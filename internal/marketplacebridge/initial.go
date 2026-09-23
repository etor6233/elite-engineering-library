package marketplacebridge

// AUTHORED official-contract mapping. PNG normalization, monetary representation,
// catalog authority, approval and durable delivery remain their existing owners.
import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"elite.local/enterprise/internal/approval"
	cr "elite.local/enterprise/internal/catalogrelease"
	"elite.local/enterprise/internal/outbounddelivery"
	sv "elite.local/enterprise/internal/storedvaluebridge"
)

var pictureID = regexp.MustCompile(`^[0-9]{1,20}-ML[A-Z][0-9]{1,24}_[0-9]{6}$`)

func ValidPictureID(s string) bool { return pictureID.MatchString(s) }
func creationBody(i Item, price json.Number) any {
	return struct {
		Family     string      `json:"family_name"`
		Category   string      `json:"category_id"`
		Price      json.Number `json:"price"`
		Currency   string      `json:"currency_id"`
		Quantity   int64       `json:"available_quantity"`
		Mode       string      `json:"buying_mode"`
		Listing    string      `json:"listing_type_id"`
		SKU        string      `json:"seller_custom_field"`
		Attributes []Attribute `json:"attributes"`
		Pictures   []Picture   `json:"pictures"`
	}{i.FamilyName, i.CategoryID, price, i.Currency, i.Quantity, i.BuyingMode, i.ListingTypeID, i.SKU, i.Attributes, i.Pictures}
}
func (c *Client) validateInitial(r Intent) error {
	if r.PrepareRequest.Validate() != nil || r.Operation != "MEDIA" && r.Operation != "CREATE" || r.TenantID != c.profile.TenantID || r.OrganizationID != c.profile.OrganizationID || r.ProfileSHA256 != c.profile.SHA256() || r.SellerID != c.profile.SellerID || !cr.ValidSHA(r.MediaSHA256) || !cr.ValidSHA(r.SourceSHA256) || !cr.ValidSHA(r.BeforeSHA256) || r.Method != "POST" || r.StockVersion != "" || r.PriceMinorUnits <= 0 || r.Quantity < 0 {
		return ErrBinding
	}
	b, e := c.profile.Binding(r.VariantID)
	if e != nil || b.SKU != r.SKU {
		return ErrBinding
	}
	var want any
	if r.Operation == "MEDIA" {
		if r.Path != "/pictures/items/upload" {
			return ErrBinding
		}
		want = map[string]string{"catalog_png_sha256": r.MediaSHA256}
	} else {
		x := r.Expected
		if r.Path != "/items" || b.StockMode != "item" || x.ID != "" || x.UserProductID != "" || x.SKU != b.SKU || x.CategoryID != b.CategoryID || x.ListingTypeID != b.ListingTypeID || x.BuyingMode != "buy_it_now" || x.Currency != c.profile.Currency || x.SiteID != c.profile.SiteID || x.Quantity != r.Quantity || !cr.ValidText(x.FamilyName, 256) || !AttributesEqual(b.Attributes, x.Attributes) || len(x.Attributes) != len(b.Attributes) || len(x.Pictures) != 1 || !ValidPictureID(x.Pictures[0].ID) || x.Pictures[0].Source != "" || x.Pictures[0].SecureURL != "" {
			return ErrBinding
		}
		want = creationBody(x, json.Number(sv.Major(r.PriceMinorUnits, c.profile.CurrencyDigits)))
	}
	raw, _, e := cr.Canonical(want)
	actual, _, a := approval.CanonicalPayload(r.Body)
	if e != nil || a != nil || !bytes.Equal(raw, actual) {
		return ErrBinding
	}
	return nil
}

// Search is bounded and authenticated. Empty search is only a precondition,
// never proof that a previous ambiguous write had no effect.
func (c *Client) SearchSKU(ctx context.Context, variant string) ([]string, error) {
	b, e := c.profile.Binding(variant)
	if e != nil {
		return nil, e
	}
	status, raw, _, e := c.call(ctx, "GET", "/users/"+c.profile.SellerID+"/items/search?sku="+url.QueryEscape(b.SKU)+"&limit=2", nil, "")
	var result struct {
		SellerID string   `json:"seller_id"`
		Results  []string `json:"results"`
		Paging   struct {
			Total  int `json:"total"`
			Offset int `json:"offset"`
		} `json:"paging"`
	}
	if e != nil || status != 200 || json.Unmarshal(raw, &result) != nil || result.SellerID != c.profile.SellerID || result.Paging.Offset != 0 || result.Paging.Total < 0 || result.Paging.Total > 1 || len(result.Results) != result.Paging.Total {
		return nil, ErrUnknown
	}
	for _, id := range result.Results {
		if !itemID.MatchString(id) || id[:3] != c.profile.SiteID {
			return nil, ErrBinding
		}
	}
	return result.Results, nil
}
func (c *Client) CreationPreflight(ctx context.Context, variant string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	status, raw, _, e := c.call(ctx, "GET", "/users/me", nil, "")
	var user struct {
		ID   int64    `json:"id"`
		Site string   `json:"site_id"`
		Tags []string `json:"tags"`
	}
	if e != nil || status != 200 || json.Unmarshal(raw, &user) != nil || strconv.FormatInt(user.ID, 10) != c.profile.SellerID || user.Site != c.profile.SiteID || !hasTag(user.Tags, "user_product_seller") {
		return ErrBinding
	}
	ids, e := c.SearchSKU(ctx, variant)
	if e != nil {
		return e
	}
	if len(ids) != 0 {
		return ErrBinding
	}
	return nil
}

type InitialObservation struct {
	ProviderID  string          `json:"provider_id"`
	Operation   string          `json:"operation"`
	MediaSHA256 string          `json:"media_sha256"`
	Response    json.RawMessage `json:"response"`
}
type RecordInitial func(context.Context, InitialObservation) error

func initialReceipt(o InitialObservation) (outbounddelivery.Receipt, error) {
	_, h, e := cr.Canonical(o)
	return outbounddelivery.Receipt{ProviderMessageID: o.ProviderID, EvidenceSHA256: h, AcceptedAt: time.Now().UTC()}, e
}
func ValidateUploadObservation(r Intent, o InitialObservation) error {
	if r.Operation != "MEDIA" || o.Operation != "MEDIA" || o.MediaSHA256 != r.MediaSHA256 || !ValidPictureID(o.ProviderID) {
		return ErrBinding
	}
	var response struct {
		ID         string `json:"id"`
		Variations []struct {
			Size string `json:"size"`
			URL  string `json:"secure_url"`
		} `json:"variations"`
	}
	if _, _, e := approval.CanonicalPayload(o.Response); e != nil {
		return ErrBinding
	}
	if json.Unmarshal(o.Response, &response) != nil || response.ID != o.ProviderID || len(response.Variations) < 1 || len(response.Variations) > 16 {
		return ErrBinding
	}
	for _, v := range response.Variations {
		u, e := url.Parse(v.URL)
		if e != nil || u.Scheme != "https" || u.Host != "http2.mlstatic.com" || u.User != nil || u.RawQuery != "" || u.Fragment != "" || !strings.HasPrefix(u.Path, "/D_NQ_NP_"+o.ProviderID+"-") || !strings.HasSuffix(u.Path, ".jpg") || !regexp.MustCompile(`^[1-9][0-9]{0,4}x[1-9][0-9]{0,4}$`).MatchString(v.Size) {
			return ErrBinding
		}
	}
	return nil
}
func (c *Client) AcceptedUpload(r Intent, o InitialObservation) (outbounddelivery.Receipt, error) {
	if c.validateInitial(r) != nil || ValidateUploadObservation(r, o) != nil {
		return outbounddelivery.Receipt{}, ErrBinding
	}
	return initialReceipt(o)
}
func (c *Client) Upload(ctx context.Context, r Intent, png []byte, record RecordInitial) (outbounddelivery.Receipt, error) {
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	sum := sha256.Sum256(png)
	if c.validateInitial(r) != nil || r.Operation != "MEDIA" || !r.ExpiresAt.After(time.Now()) || len(png) < 8 || len(png) > 4*1024*1024 || !bytes.Equal(png[:8], []byte{137, 80, 78, 71, 13, 10, 26, 10}) || hex.EncodeToString(sum[:]) != r.MediaSHA256 || record == nil {
		return outbounddelivery.Receipt{}, ErrBinding
	}
	// Only already-normalized catalog PNG reaches this method. No new image codec.
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	if e := writer.SetBoundary("elite-" + r.MediaSHA256); e != nil {
		return outbounddelivery.Receipt{}, ErrBinding
	}
	part, e := writer.CreateFormFile("file", r.MediaSHA256+".png")
	if e != nil {
		return outbounddelivery.Receipt{}, ErrBinding
	}
	if _, e = part.Write(png); e != nil {
		return outbounddelivery.Receipt{}, ErrBinding
	}
	if writer.Close() != nil {
		return outbounddelivery.Receipt{}, ErrBinding
	}
	token, e := c.tokens.MercadoLibreAccessToken(ctx)
	if e != nil || len(token) < 20 || len(token) > 4096 || strings.ContainsAny(token, "\r\n") {
		return outbounddelivery.Receipt{}, ErrBinding
	}
	request, e := http.NewRequestWithContext(ctx, "POST", APIBaseURL+"/pictures/items/upload", &body)
	if e != nil {
		return outbounddelivery.Receipt{}, ErrBinding
	}
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", writer.FormDataContentType())
	response, e := c.http.Do(request)
	if e != nil {
		return outbounddelivery.Receipt{}, ErrUnknown
	}
	defer response.Body.Close()
	raw, e := io.ReadAll(io.LimitReader(response.Body, 32769))
	if e != nil || len(raw) > 32768 {
		return outbounddelivery.Receipt{}, ErrUnknown
	}
	if response.StatusCode != 200 && response.StatusCode != 201 {
		return outbounddelivery.Receipt{}, initialRejection(response.StatusCode, raw)
	}
	var upload struct {
		ID string `json:"id"`
	}
	if json.Unmarshal(raw, &upload) != nil {
		return outbounddelivery.Receipt{}, ErrUnknown
	}
	o := InitialObservation{upload.ID, "MEDIA", r.MediaSHA256, raw}
	if ValidateUploadObservation(r, o) != nil {
		return outbounddelivery.Receipt{}, ErrUnknown
	}
	if e = record(ctx, o); e != nil {
		return outbounddelivery.Receipt{}, ErrUnknown
	}
	return initialReceipt(o)
}
func initialRejection(status int, raw []byte) error {
	if status == 400 || status == 401 || status == 403 || status == 409 || status == 422 {
		if _, h, e := approval.CanonicalPayload(raw); e == nil {
			return outbounddelivery.NewTerminalFailure("MERCADOLIBRE_INITIAL_REJECTED", h, nil)
		}
	}
	return ErrUnknown
}
func (c *Client) Create(ctx context.Context, r Intent, record RecordInitial) (outbounddelivery.Receipt, error) {
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	if c.validateInitial(r) != nil || r.Operation != "CREATE" || !r.ExpiresAt.After(time.Now()) || record == nil {
		return outbounddelivery.Receipt{}, ErrBinding
	}
	if e := c.CreationPreflight(ctx, r.VariantID); e != nil {
		return outbounddelivery.Receipt{}, e
	}
	wire, _, e := approval.CanonicalPayload(r.Body)
	if e != nil {
		return outbounddelivery.Receipt{}, ErrBinding
	}
	status, raw, _, e := c.call(ctx, "POST", "/items", wire, "")
	if e != nil {
		return outbounddelivery.Receipt{}, ErrUnknown
	}
	if status != 201 {
		return outbounddelivery.Receipt{}, initialRejection(status, raw)
	}
	var created struct {
		ID string `json:"id"`
	}
	if json.Unmarshal(raw, &created) != nil || !itemID.MatchString(created.ID) || created.ID[:3] != c.profile.SiteID {
		return outbounddelivery.Receipt{}, ErrUnknown
	}
	if e = record(ctx, InitialObservation{created.ID, "CREATE", r.MediaSHA256, raw}); e != nil {
		return outbounddelivery.Receipt{}, ErrUnknown
	}
	return c.ReconcileCreation(ctx, r, created.ID)
}
func (c *Client) ReconcileCreation(ctx context.Context, r Intent, knownID string) (outbounddelivery.Receipt, error) {
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	if c.validateInitial(r) != nil || r.Operation != "CREATE" {
		return outbounddelivery.Receipt{}, ErrBinding
	}
	if knownID == "" {
		ids, e := c.SearchSKU(ctx, r.VariantID)
		if e != nil || len(ids) != 1 {
			return outbounddelivery.Receipt{}, ErrUnknown
		}
		knownID = ids[0]
	}
	observed, e := c.Observe(ctx, r.VariantID, knownID)
	if e != nil {
		return outbounddelivery.Receipt{}, ErrUnknown
	}
	x := observed.Item
	w := r.Expected
	if observed.PriceMinor != strconv.FormatInt(r.PriceMinorUnits, 10) || observed.Quantity != r.Quantity || x.FamilyName != w.FamilyName || x.ListingTypeID != w.ListingTypeID || x.BuyingMode != w.BuyingMode || !AttributesEqual(w.Attributes, x.Attributes) || len(x.Pictures) != 1 || x.Pictures[0].ID != w.Pictures[0].ID || x.Status != "active" && x.Status != "paused" && x.Status != "under_review" {
		return outbounddelivery.Receipt{}, ErrUnknown
	}
	_, h, e := cr.Canonical(map[string]any{"request": r, "observation": observed})
	return outbounddelivery.Receipt{ProviderMessageID: knownID, EvidenceSHA256: h, AcceptedAt: time.Now().UTC()}, e
}
