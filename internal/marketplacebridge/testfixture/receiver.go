// Package testfixture is AUTHORED synthetic provider infrastructure. It is not a
// provider emulator certification; requests use the fixed documented paths.
package testfixture

import (
	mb "elite.local/enterprise/internal/marketplacebridge"
	sv "elite.local/enterprise/internal/storedvaluebridge"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"sync"
)

const Token = "synthetic-mercadolibre-fixture-token"

type Write struct {
	Path    string
	Body    string
	Version string
}
type Receiver struct {
	Mu                                              sync.Mutex
	Server                                          *httptest.Server
	Profile                                         mb.Profile
	Item                                            mb.Item
	Quantity, PriceMinor, Version                   int64
	DropNext, IgnorePrice, RejectNext, ConflictNext bool
	Writes                                          []Write
	Effects, Reads                                  int
}
type Transport struct {
	Target *url.URL
	Base   http.RoundTripper
}

func (t Transport) RoundTrip(r *http.Request) (*http.Response, error) {
	if r.URL.Scheme != "https" || r.URL.Host != "api.mercadolibre.com" || r.Header.Get("Authorization") != "Bearer "+Token {
		return nil, mb.ErrBinding
	}
	copy := r.Clone(r.Context())
	u := *r.URL
	u.Scheme = t.Target.Scheme
	u.Host = t.Target.Host
	copy.URL = &u
	copy.Host = t.Target.Host
	return t.Base.RoundTrip(copy)
}
func New(p mb.Profile) *Receiver {
	b := p.Bindings[0]
	seller, _ := strconv.ParseInt(p.SellerID, 10, 64)
	s := &Receiver{Profile: p, PriceMinor: 9007199254740993, Quantity: 4, Version: 7}
	s.Item = mb.Item{ID: "MLA123456789", SellerID: seller, SiteID: p.SiteID, CategoryID: b.CategoryID, ListingTypeID: b.ListingTypeID, BuyingMode: "buy_it_now", FamilyName: "Fixture bicycle", UserProductID: "MLAU123456789", SKU: b.SKU, Currency: p.Currency, Quantity: s.Quantity, Status: "active", Attributes: b.Attributes}
	s.Server = httptest.NewServer(http.HandlerFunc(s.serve))
	return s
}
func (s *Receiver) HTTP() *http.Client {
	u, _ := url.Parse(s.Server.URL)
	client := mb.NewHTTPClient()
	client.Transport = Transport{u, client.Transport}
	return client
}
func (s *Receiver) Close() { s.Server.Close() }
func (s *Receiver) serve(w http.ResponseWriter, r *http.Request) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	if r.Header.Get("Authorization") != "Bearer "+Token {
		w.WriteHeader(401)
		io.WriteString(w, `{"error":"unauthorized"}`)
		return
	}
	b := s.Profile.Bindings[0]
	itemPath := "/items/" + s.Item.ID
	stockPath := "/user-products/" + s.Item.UserProductID + "/stock"
	if r.Method == "GET" {
		s.Reads++
		switch r.URL.Path {
		case itemPath:
			_ = json.NewEncoder(w).Encode(s.Item)
		case itemPath + "/prices":
			_ = json.NewEncoder(w).Encode(map[string]any{"id": s.Item.ID, "prices": []any{map[string]any{"type": "standard", "amount": json.Number(sv.Major(s.PriceMinor, s.Profile.CurrencyDigits)), "currency_id": s.Profile.Currency, "conditions": map[string]any{"context_restrictions": []string{}}}}})
		case stockPath:
			location := map[string]any{"type": "selling_address", "quantity": s.Quantity}
			if b.StockMode == "seller_warehouse" {
				location = map[string]any{"type": b.StockMode, "quantity": s.Quantity, "store_id": b.StoreID, "network_node_id": b.NetworkNodeID}
			}
			locations := []any{location}
			if b.StockMode != "item" {
				locations = append(locations, map[string]any{"type": "meli_facility", "quantity": 99})
			}
			w.Header().Set("x-version", strconv.FormatInt(s.Version, 10))
			_ = json.NewEncoder(w).Encode(map[string]any{"id": s.Item.UserProductID, "user_id": s.Item.SellerID, "locations": locations})
		default:
			w.WriteHeader(404)
			io.WriteString(w, `{"error":"not_found"}`)
		}
		return
	}
	if r.Method != "PUT" {
		w.WriteHeader(405)
		io.WriteString(w, `{"error":"method_not_allowed"}`)
		return
	}
	raw, e := io.ReadAll(io.LimitReader(r.Body, 32769))
	if e != nil || len(raw) > 32768 {
		w.WriteHeader(413)
		return
	}
	var input map[string]json.RawMessage
	if json.Unmarshal(raw, &input) != nil {
		w.WriteHeader(400)
		return
	}
	s.Writes = append(s.Writes, Write{r.URL.Path, string(raw), r.Header.Get("x-version")})
	if s.RejectNext {
		s.RejectNext = false
		w.WriteHeader(400)
		io.WriteString(w, `{"error":"validation_error"}`)
		return
	}
	if s.ConflictNext {
		s.ConflictNext = false
		w.WriteHeader(409)
		io.WriteString(w, `{"error":"version_mismatch"}`)
		return
	}
	if r.URL.Path == itemPath {
		if len(input) != 1 {
			w.WriteHeader(400)
			io.WriteString(w, `{"error":"unexpected_fields"}`)
			return
		}
		if x, ok := input["price"]; ok {
			var n json.Number
			if json.Unmarshal(x, &n) != nil {
				w.WriteHeader(400)
				return
			}
			value, e := sv.Minor(n.String(), s.Profile.CurrencyDigits)
			if e != nil {
				w.WriteHeader(400)
				return
			}
			if !s.IgnorePrice {
				s.PriceMinor = value
				s.Effects++
			}
		} else if x, ok := input["available_quantity"]; ok {
			if b.StockMode != "item" || json.Unmarshal(x, &s.Quantity) != nil {
				w.WriteHeader(400)
				return
			}
			s.Item.Quantity = s.Quantity
			s.Effects++
		} else if x, ok := input["status"]; ok {
			var status string
			if json.Unmarshal(x, &status) != nil || status != "paused" && status != "active" {
				w.WriteHeader(400)
				return
			}
			s.Item.Status = status
			s.Effects++
		} else {
			w.WriteHeader(400)
			return
		}
	} else if r.URL.Path == stockPath+"/type/"+b.StockMode && b.StockMode != "item" {
		if r.Header.Get("x-version") != strconv.FormatInt(s.Version, 10) {
			w.WriteHeader(409)
			io.WriteString(w, `{"error":"version_mismatch"}`)
			return
		}
		if b.StockMode == "selling_address" {
			if len(input) != 1 || json.Unmarshal(input["quantity"], &s.Quantity) != nil {
				w.WriteHeader(400)
				return
			}
		} else {
			var locations []struct {
				Store    string `json:"store_id"`
				Node     string `json:"network_node_id"`
				Quantity int64  `json:"quantity"`
			}
			if len(input) != 1 || json.Unmarshal(input["locations"], &locations) != nil || len(locations) != 1 || locations[0].Store != b.StoreID || locations[0].Node != b.NetworkNodeID {
				w.WriteHeader(400)
				return
			}
			s.Quantity = locations[0].Quantity
		}
		s.Version++
		s.Effects++
	} else {
		w.WriteHeader(404)
		io.WriteString(w, `{"error":"not_found"}`)
		return
	}
	if s.DropNext {
		s.DropNext = false
		conn, _, e := w.(http.Hijacker).Hijack()
		if e == nil {
			conn.Close()
		}
		return
	}
	if strings.HasPrefix(r.URL.Path, stockPath) {
		w.WriteHeader(204)
	} else {
		io.WriteString(w, `{"accepted":true}`)
	}
}
