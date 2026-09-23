package testfixture

// AUTHORED loopback fixture for the locked multipart and User Products contracts.
import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"

	mb "elite.local/enterprise/internal/marketplacebridge"
	sv "elite.local/enterprise/internal/storedvaluebridge"
)

const PictureID = "123456-MLA123456789_092026"

type InitialReceiver struct {
	*Receiver
	Exists, DropUpload, DropCreate, BadPicture, DuplicateSearch bool
	Uploads, Creates                                            int
	UploadedSHA                                                 string
}

func NewInitial(p mb.Profile) *InitialReceiver {
	base := New(p)
	base.Server.Close()
	s := &InitialReceiver{Receiver: base}
	base.Server = httptest.NewServer(http.HandlerFunc(s.serveInitial))
	return s
}
func (s *InitialReceiver) serveInitial(w http.ResponseWriter, r *http.Request) {
	s.Mu.Lock()
	handled := true
	defer func() {
		s.Mu.Unlock()
		if !handled {
			s.Receiver.serve(w, r)
		}
	}()
	w.Header().Set("Content-Type", "application/json")
	if r.Header.Get("Authorization") != "Bearer "+Token {
		w.WriteHeader(401)
		io.WriteString(w, `{"error":"unauthorized"}`)
		return
	}
	if r.Method == "GET" && r.URL.Path == "/users/me" {
		s.Reads++
		json.NewEncoder(w).Encode(map[string]any{"id": s.Item.SellerID, "site_id": s.Profile.SiteID, "tags": []string{"user_product_seller"}})
		return
	}
	if r.Method == "GET" && r.URL.Path == "/users/"+s.Profile.SellerID+"/items/search" {
		s.Reads++
		if r.URL.Query().Get("sku") != s.Profile.Bindings[0].SKU || r.URL.Query().Get("limit") != "2" {
			w.WriteHeader(400)
			return
		}
		ids := []string{}
		if s.Exists {
			ids = append(ids, s.Item.ID)
		}
		if s.DuplicateSearch {
			ids = append(ids, "MLA987654321")
		}
		json.NewEncoder(w).Encode(map[string]any{"seller_id": s.Profile.SellerID, "results": ids, "paging": map[string]any{"total": len(ids), "offset": 0, "limit": 2}})
		return
	}
	if r.Method == "POST" && r.URL.Path == "/pictures/items/upload" {
		s.Uploads++
		reader, e := r.MultipartReader()
		if e != nil {
			w.WriteHeader(400)
			return
		}
		part, e := reader.NextPart()
		if e != nil || part.FormName() != "file" {
			w.WriteHeader(400)
			return
		}
		raw, e := io.ReadAll(io.LimitReader(part, 4*1024*1024+1))
		if e != nil || len(raw) > 4*1024*1024 {
			w.WriteHeader(413)
			return
		}
		sum := sha256.Sum256(raw)
		s.UploadedSHA = hex.EncodeToString(sum[:])
		if part.FileName() != s.UploadedSHA+".png" {
			w.WriteHeader(400)
			return
		}
		if _, e = reader.NextPart(); e != io.EOF {
			w.WriteHeader(400)
			return
		}
		if s.DropUpload {
			s.DropUpload = false
			conn, _, _ := w.(http.Hijacker).Hijack()
			conn.Close()
			return
		}
		w.WriteHeader(201)
		json.NewEncoder(w).Encode(map[string]any{"id": PictureID, "variations": []any{map[string]string{"size": "2x2", "secure_url": "https://http2.mlstatic.com/D_NQ_NP_" + PictureID + "-F.jpg"}}})
		return
	}
	if r.Method == "POST" && r.URL.Path == "/items" {
		s.Creates++
		raw, e := io.ReadAll(io.LimitReader(r.Body, 32769))
		if e != nil || len(raw) > 32768 {
			w.WriteHeader(413)
			return
		}
		var input map[string]json.RawMessage
		var next mb.Item
		var price json.Number
		if json.Unmarshal(raw, &input) != nil || len(input) != 10 || json.Unmarshal(raw, &next) != nil || json.Unmarshal(input["price"], &price) != nil {
			w.WriteHeader(400)
			io.WriteString(w, `{"error":"invalid_item"}`)
			return
		}
		amount, e := sv.Minor(price.String(), s.Profile.CurrencyDigits)
		if e != nil || next.SKU != s.Profile.Bindings[0].SKU || len(next.Pictures) != 1 || next.Pictures[0].ID != PictureID || s.UploadedSHA == "" || next.FamilyName == "" {
			w.WriteHeader(400)
			io.WriteString(w, `{"error":"invalid_binding"}`)
			return
		}
		if s.RejectNext {
			s.RejectNext = false
			w.WriteHeader(400)
			io.WriteString(w, `{"error":"category_validation"}`)
			return
		}
		if s.Exists {
			w.WriteHeader(409)
			io.WriteString(w, `{"error":"fixture_existing_sku"}`)
			return
		}
		next.ID = s.Item.ID
		next.UserProductID = s.Item.UserProductID
		next.SellerID = s.Item.SellerID
		next.SiteID = s.Profile.SiteID
		next.Status = "paused"
		if next.Quantity > 0 {
			next.Status = "active"
		}
		if s.BadPicture {
			next.Pictures[0].ID = "654321-MLA123456789_092026"
		}
		s.Item = next
		s.Quantity = next.Quantity
		s.PriceMinor = amount
		s.Exists = true
		s.Effects++
		s.Writes = append(s.Writes, Write{Path: r.URL.Path, Body: string(raw)})
		if s.DropCreate {
			s.DropCreate = false
			conn, _, _ := w.(http.Hijacker).Hijack()
			conn.Close()
			return
		}
		w.WriteHeader(201)
		io.WriteString(w, `{"id":`+strconv.Quote(next.ID)+`}`)
		return
	}
	if !s.Exists && r.Method == "GET" {
		w.WriteHeader(404)
		io.WriteString(w, `{"error":"not_found"}`)
		return
	}
	handled = false
}
