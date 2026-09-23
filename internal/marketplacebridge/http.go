package marketplacebridge

// AUTHORED bounded HTTP/serialization glue for official provider contracts.
// A 2xx response alone does not confirm the requested effect.
import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"elite.local/enterprise/internal/approval"
	"elite.local/enterprise/internal/outbounddelivery"
	sv "elite.local/enterprise/internal/storedvaluebridge"
)

const APIBaseURL = "https://api.mercadolibre.com"

type Doer interface {
	Do(*http.Request) (*http.Response, error)
}
type TokenSource interface {
	MercadoLibreAccessToken(context.Context) (string, error)
}
type Client struct {
	profile Profile
	tokens  TokenSource
	http    Doer
}

func NewClient(p Profile, t TokenSource, h Doer) (*Client, error) {
	if p.Validate() != nil || t == nil || h == nil {
		return nil, ErrBinding
	}
	raw, _ := json.Marshal(p)
	var cloned Profile
	if json.Unmarshal(raw, &cloned) != nil {
		return nil, ErrBinding
	}
	return &Client{cloned, t, h}, nil
}
func NewHTTPClient() *http.Client {
	return &http.Client{Timeout: 8 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
		Transport: &http.Transport{Proxy: nil, DialContext: (&net.Dialer{Timeout: 3 * time.Second, KeepAlive: 30 * time.Second}).DialContext, TLSClientConfig: &tls.Config{MinVersion: tls.VersionTLS12}, TLSHandshakeTimeout: 3 * time.Second, ResponseHeaderTimeout: 5 * time.Second, MaxResponseHeaderBytes: 32768}}
}
func (c *Client) call(ctx context.Context, method, path string, body []byte, version string) (int, []byte, http.Header, error) {
	if !strings.HasPrefix(path, "/") || strings.HasPrefix(path, "//") || strings.ContainsAny(path, "\r\n\\") || len(path) > 512 || len(body) > 32768 {
		return 0, nil, nil, ErrBinding
	}
	token, e := c.tokens.MercadoLibreAccessToken(ctx)
	if e != nil || len(token) < 20 || len(token) > 4096 || strings.ContainsAny(token, "\r\n") {
		return 0, nil, nil, ErrBinding
	}
	request, e := http.NewRequestWithContext(ctx, method, APIBaseURL+path, bytes.NewReader(body))
	if e != nil {
		return 0, nil, nil, ErrBinding
	}
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("Accept", "application/json")
	if len(body) > 0 {
		request.Header.Set("Content-Type", "application/json")
	}
	if version != "" {
		request.Header.Set("x-version", version)
	}
	response, e := c.http.Do(request)
	if e != nil {
		return 0, nil, nil, ErrUnknown
	}
	defer response.Body.Close()
	raw, e := io.ReadAll(io.LimitReader(response.Body, 32769))
	if e != nil || len(raw) > 32768 {
		return response.StatusCode, nil, nil, ErrUnknown
	}
	if len(raw) > 0 {
		if _, _, e = approval.CanonicalPayload(raw); e != nil {
			return response.StatusCode, nil, nil, ErrUnknown
		}
	}
	return response.StatusCode, raw, response.Header.Clone(), nil
}
func hasTag(tags []string, want string) bool {
	for _, s := range tags {
		if s == want {
			return true
		}
	}
	return false
}
func (c *Client) Observe(ctx context.Context, variant, id string) (Observation, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	var out Observation
	b, e := c.profile.Binding(variant)
	if e != nil || !itemID.MatchString(id) || id[:3] != c.profile.SiteID {
		return out, ErrBinding
	}
	status, raw, _, e := c.call(ctx, "GET", "/items/"+id, nil, "")
	if e != nil || status != 200 || json.Unmarshal(raw, &out.Item) != nil {
		return out, ErrUnknown
	}
	v := out.Item
	if v.ID != id || strconv.FormatInt(v.SellerID, 10) != c.profile.SellerID || v.SiteID != c.profile.SiteID || v.CategoryID != b.CategoryID || v.SKU != b.SKU || !upID.MatchString(v.UserProductID) || v.FamilyName == "" || v.Currency != c.profile.Currency {
		return out, ErrBinding
	}
	status, raw, _, e = c.call(ctx, "GET", "/items/"+id+"/prices", nil, "")
	var prices struct {
		ID     string `json:"id"`
		Prices []struct {
			Type       string      `json:"type"`
			Amount     json.Number `json:"amount"`
			Currency   string      `json:"currency_id"`
			Conditions struct {
				Context []string `json:"context_restrictions"`
				Start   *string  `json:"start_time"`
				End     *string  `json:"end_time"`
				Min     *int64   `json:"min_purchase_unit"`
			} `json:"conditions"`
		} `json:"prices"`
	}
	if e != nil || status != 200 || json.Unmarshal(raw, &prices) != nil || prices.ID != id {
		return out, ErrUnknown
	}
	matched := 0
	for _, p := range prices.Prices {
		if p.Type != "standard" || p.Conditions.Start != nil || p.Conditions.End != nil || p.Conditions.Min != nil {
			continue
		}
		if len(p.Conditions.Context) != 0 && !(len(p.Conditions.Context) == 1 && p.Conditions.Context[0] == "channel_marketplace") {
			continue
		}
		if p.Currency != c.profile.Currency {
			return out, ErrBinding
		}
		n, e := sv.Minor(p.Amount.String(), c.profile.CurrencyDigits)
		if e != nil || n <= 0 {
			return out, ErrBinding
		}
		out.PriceMinor = strconv.FormatInt(n, 10)
		matched++
	}
	if matched != 1 {
		return out, ErrBinding
	}
	status, raw, headers, e := c.call(ctx, "GET", "/user-products/"+v.UserProductID+"/stock", nil, "")
	var stock struct {
		ID        string `json:"id"`
		UserID    int64  `json:"user_id"`
		Locations []struct {
			Type     string `json:"type"`
			StoreID  string `json:"store_id"`
			Node     string `json:"network_node_id"`
			Quantity int64  `json:"quantity"`
		} `json:"locations"`
	}
	if e != nil || status != 200 || json.Unmarshal(raw, &stock) != nil || stock.ID != v.UserProductID || strconv.FormatInt(stock.UserID, 10) != c.profile.SellerID {
		return out, ErrUnknown
	}
	matched = 0
	for _, loc := range stock.Locations {
		if loc.Quantity < 0 {
			return out, ErrBinding
		}
		switch b.StockMode {
		case "item":
			// Legacy available_quantity write is selected only for this normal,
			// non-multi-origin, non-Full inventory configuration.
			if len(stock.Locations) != 1 || loc.Type != "selling_address" || loc.StoreID != "" || loc.Node != "" {
				return out, ErrBinding
			}
			out.Quantity = v.Quantity
			matched++
		case "selling_address":
			if loc.Type == "seller_warehouse" {
				return out, ErrBinding
			}
			if loc.Type == "selling_address" {
				out.Quantity = loc.Quantity
				matched++
			}
		case "seller_warehouse":
			if loc.Type == "selling_address" {
				return out, ErrBinding
			}
			if loc.Type == "seller_warehouse" && loc.StoreID == b.StoreID && loc.Node == b.NetworkNodeID {
				out.Quantity = loc.Quantity
				matched++
			}
		}
	}
	if matched != 1 || out.Quantity < 0 {
		return out, ErrBinding
	}
	if b.StockMode != "item" {
		version := headers.Get("x-version")
		n, e := strconv.ParseInt(version, 10, 64)
		if e != nil || n < 0 || strconv.FormatInt(n, 10) != version {
			return out, ErrBinding
		}
		out.StockVersion = version
	}
	out.Exists = true
	return out, nil
}
func (c *Client) validate(r Intent) error {
	if r.PrepareRequest.Validate() != nil || r.TenantID != c.profile.TenantID || r.OrganizationID != c.profile.OrganizationID || r.ProfileSHA256 != c.profile.SHA256() || r.SellerID != c.profile.SellerID || r.Method != "PUT" || r.PriceMinorUnits <= 0 || r.Quantity < 0 {
		return ErrBinding
	}
	b, e := c.profile.Binding(r.VariantID)
	if e != nil || r.SKU != b.SKU {
		return ErrBinding
	}
	path := "/items/" + r.ItemID
	version := ""
	var body any
	switch r.Operation {
	case "CONTENT":
		if validateContent(c.profile, r) != nil {
			return ErrBinding
		}
		body = contentBody(r.Expected)
	case "PRICE":
		body = map[string]any{"price": json.Number(sv.Major(r.PriceMinorUnits, c.profile.CurrencyDigits))}
	case "STOCK":
		if b.StockMode == "item" {
			body = map[string]any{"available_quantity": r.Quantity}
		} else {
			if !upID.MatchString(r.Expected.UserProductID) {
				return ErrBinding
			}
			path = "/user-products/" + r.Expected.UserProductID + "/stock/type/" + b.StockMode
			version = r.StockVersion
			n, e := strconv.ParseInt(version, 10, 64)
			if e != nil || n < 0 || strconv.FormatInt(n, 10) != version {
				return ErrBinding
			}
			body = map[string]any{"quantity": r.Quantity}
			if b.StockMode == "seller_warehouse" {
				body = map[string]any{"locations": []map[string]any{{"store_id": b.StoreID, "network_node_id": b.NetworkNodeID, "quantity": r.Quantity}}}
			}
		}
	case "PAUSE":
		body = map[string]string{"status": "paused"}
	case "RESUME":
		body = map[string]string{"status": "active"}
	}
	raw, _ := json.Marshal(body)
	canonical, _, e := approval.CanonicalPayload(raw)
	actual, _, parseErr := approval.CanonicalPayload(r.Body)
	if e != nil || parseErr != nil || r.Path != path || r.StockVersion != version || !bytes.Equal(canonical, actual) {
		return ErrBinding
	}
	return nil
}
func (c *Client) confirms(r Intent, o Observation) bool {
	if !o.Exists || o.Item.ID != r.ItemID || o.Item.UserProductID != r.Expected.UserProductID {
		return false
	}
	switch r.Operation {
	case "CONTENT":
		return o.Item.FamilyName == r.Expected.FamilyName && AttributesEqual(r.Expected.Attributes, o.Item.Attributes) && len(o.Item.Pictures) == 1 && o.Item.Pictures[0].ID == r.Expected.Pictures[0].ID
	case "PRICE":
		return o.PriceMinor == strconv.FormatInt(r.PriceMinorUnits, 10)
	case "STOCK":
		return o.Quantity == r.Quantity
	case "PAUSE":
		return o.Item.Status == "paused"
	case "RESUME":
		return o.Item.Status == "active"
	}
	return false
}
func (c *Client) Reconcile(ctx context.Context, r Intent) (outbounddelivery.Receipt, error) {
	if r.Operation == "CREATE" {
		return c.ReconcileCreation(ctx, r, "")
	}
	if c.validate(r) != nil {
		return outbounddelivery.Receipt{}, ErrBinding
	}
	o, e := c.Observe(ctx, r.VariantID, r.ItemID)
	if e != nil || !c.confirms(r, o) {
		return outbounddelivery.Receipt{}, ErrUnknown
	}
	return outbounddelivery.Receipt{ProviderMessageID: r.ItemID, EvidenceSHA256: o.SHA256(), AcceptedAt: time.Now().UTC()}, nil
}
func (c *Client) Write(ctx context.Context, r Intent) (outbounddelivery.Receipt, error) {
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	if c.validate(r) != nil || !r.ExpiresAt.After(time.Now()) {
		return outbounddelivery.Receipt{}, ErrBinding
	}
	o, e := c.Observe(ctx, r.VariantID, r.ItemID)
	if e != nil {
		return outbounddelivery.Receipt{}, e
	}
	if o.SHA256() != r.BeforeSHA256 || o.Item.UserProductID != r.Expected.UserProductID || r.Operation == "PRICE" && hasTag(o.Item.Tags, "dynamic_standard_price") {
		return outbounddelivery.Receipt{}, outbounddelivery.NewTerminalFailure("MERCADOLIBRE_PRECONDITION_CHANGED", o.SHA256(), nil)
	}
	if r.Operation == "CONTENT" {
		if e := c.SingleUnsoldItem(ctx, o.Item); e != nil {
			return outbounddelivery.Receipt{}, outbounddelivery.NewTerminalFailure("MERCADOLIBRE_CONTENT_SCOPE_CHANGED", o.SHA256(), nil)
		}
	}
	wire, _, e := approval.CanonicalPayload(r.Body)
	if e != nil {
		return outbounddelivery.Receipt{}, ErrBinding
	}
	status, raw, _, e := c.call(ctx, r.Method, r.Path, wire, r.StockVersion)
	if e != nil {
		return outbounddelivery.Receipt{}, ErrUnknown
	}
	if status == 400 || status == 401 || status == 403 || status == 409 || status == 422 {
		_, h, e := approval.CanonicalPayload(raw)
		if e != nil {
			return outbounddelivery.Receipt{}, ErrUnknown
		}
		return outbounddelivery.Receipt{}, outbounddelivery.NewTerminalFailure("MERCADOLIBRE_MUTATION_REJECTED", h, nil)
	}
	if status != 200 && status != 204 {
		return outbounddelivery.Receipt{}, ErrUnknown
	}
	receipt, e := c.Reconcile(ctx, r)
	if e != nil {
		return receipt, ErrUnknown
	}
	// Preserve the write and GET evidence together, without raw tokens or PII.
	_, h, e := approval.CanonicalPayload(mustJSON(map[string]any{"write_status": status, "write_body": json.RawMessage(nonempty(raw)), "observation_sha256": receipt.EvidenceSHA256, "request_body": r.Body}))
	if e != nil {
		return outbounddelivery.Receipt{}, ErrUnknown
	}
	receipt.EvidenceSHA256 = h
	return receipt, nil
}
func nonempty(b []byte) []byte {
	if len(b) == 0 {
		return []byte("{}")
	}
	return b
}
func mustJSON(v any) []byte {
	b, e := json.Marshal(v)
	if e != nil {
		panic(errors.New("bounded typed serialization failed"))
	}
	return b
}
