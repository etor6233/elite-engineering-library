// AUTHORED mapping of an approved catalog publication to official Mercado Libre
// HTTP contracts. Business amounts, inventory and approval use existing owners.
package marketplacebridge

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/url"
	"regexp"
	"sort"
	"time"

	"elite.local/enterprise/internal/approval"
	cr "elite.local/enterprise/internal/catalogrelease"
	"elite.local/enterprise/internal/channels"
	sv "elite.local/enterprise/internal/storedvaluebridge"
)

var ErrBinding = errors.New("marketplace: invalid or changed binding")
var ErrUnknown = errors.New("marketplace: effect not confirmed; GET reconciliation required")
var digits = regexp.MustCompile(`^[0-9]{3,20}$`)
var itemID = regexp.MustCompile(`^ML[A-Z][0-9]{6,20}$`)
var categoryID = regexp.MustCompile(`^ML[A-Z][0-9]{2,20}$`)
var upID = regexp.MustCompile(`^ML[A-Z]U[0-9]{6,20}$`)
var siteID = regexp.MustCompile(`^ML[A-Z]$`)
var code = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_.:-]{0,63}$`)

type Attribute struct {
	ID        string `json:"id"`
	ValueID   string `json:"value_id,omitempty"`
	ValueName string `json:"value_name,omitempty"`
}
type Binding struct {
	VariantID     string      `json:"variant_id"`
	SKU           string      `json:"sku"`
	CategoryID    string      `json:"category_id"`
	ListingTypeID string      `json:"listing_type_id"`
	Attributes    []Attribute `json:"attributes"`
	StockMode     string      `json:"stock_mode"`
	StoreID       string      `json:"store_id,omitempty"`
	NetworkNodeID string      `json:"network_node_id,omitempty"`
}
type Profile struct {
	Schema         string    `json:"schema"`
	TenantID       string    `json:"tenant_id"`
	OrganizationID string    `json:"organization_id"`
	SellerID       string    `json:"seller_id"`
	SiteID         string    `json:"site_id"`
	Currency       string    `json:"currency"`
	CurrencyDigits int       `json:"currency_digits"`
	MediaOrigin    string    `json:"media_origin"`
	Bindings       []Binding `json:"bindings"`
}

func (p Profile) Validate() error {
	u, e := url.Parse(p.MediaOrigin)
	if p.Schema != "elite-marketplace-publication/v1" || !cr.ValidID(p.TenantID) || !cr.ValidID(p.OrganizationID) || !digits.MatchString(p.SellerID) || !siteID.MatchString(p.SiteID) || !regexp.MustCompile(`^[A-Z]{3}$`).MatchString(p.Currency) || p.CurrencyDigits < 0 || p.CurrencyDigits > 6 || e != nil || u.Scheme != "https" || u.Host == "" || u.User != nil || u.Path != "" || u.RawQuery != "" || u.Fragment != "" || len(p.MediaOrigin) > 512 || len(p.Bindings) < 1 || len(p.Bindings) > 32 {
		return ErrBinding
	}
	variants, skus := map[string]bool{}, map[string]bool{}
	for _, b := range p.Bindings {
		if !cr.ValidID(b.VariantID) || !code.MatchString(b.SKU) || !categoryID.MatchString(b.CategoryID) || b.CategoryID[:3] != p.SiteID || !code.MatchString(b.ListingTypeID) || variants[b.VariantID] || skus[b.SKU] || len(b.Attributes) < 1 || len(b.Attributes) > 64 {
			return ErrBinding
		}
		variants[b.VariantID] = true
		skus[b.SKU] = true
		if b.StockMode != "item" && b.StockMode != "selling_address" && b.StockMode != "seller_warehouse" {
			return ErrBinding
		}
		if b.StockMode == "seller_warehouse" {
			if !digits.MatchString(b.StoreID) || !code.MatchString(b.NetworkNodeID) {
				return ErrBinding
			}
		} else if b.StoreID != "" || b.NetworkNodeID != "" {
			return ErrBinding
		}
		found := false
		seen := map[string]bool{}
		for _, a := range b.Attributes {
			if !code.MatchString(a.ID) || seen[a.ID] || (a.ValueID == "") == (a.ValueName == "") || a.ValueID != "" && !code.MatchString(a.ValueID) || a.ValueName != "" && !cr.ValidText(a.ValueName, 256) || a.ID == "SELLER_SKU" {
				return ErrBinding
			}
			seen[a.ID] = true
			if a.ID == "ITEM_CONDITION" {
				found = true
			}
		}
		if !found {
			return ErrBinding
		}
	}
	raw, e := json.Marshal(p)
	if e != nil {
		return ErrBinding
	}
	if _, _, e = approval.CanonicalPayload(raw); e != nil {
		return ErrBinding
	}
	return nil
}
func (p Profile) SHA256() string { _, h, _ := cr.Canonical(p); return h }
func (p Profile) Binding(variant string) (Binding, error) {
	for _, b := range p.Bindings {
		if b.VariantID == variant {
			return b, nil
		}
	}
	return Binding{}, ErrBinding
}
func Decode(raw []byte, v any) error {
	if len(raw) == 0 || len(raw) > 65536 {
		return ErrBinding
	}
	if _, _, e := approval.CanonicalPayload(raw); e != nil {
		return ErrBinding
	}
	d := json.NewDecoder(bytes.NewReader(raw))
	d.DisallowUnknownFields()
	if d.Decode(v) != nil || d.Decode(new(any)) != io.EOF {
		return ErrBinding
	}
	return nil
}

type PrepareRequest struct {
	ApprovalID      string    `json:"approval_id"`
	Generation      int64     `json:"generation,string"`
	VariantID       string    `json:"variant_id"`
	Operation       string    `json:"operation"`
	ItemID          string    `json:"item_id,omitempty"`
	MediaApprovalID string    `json:"media_approval_id,omitempty"`
	ExpiresAt       time.Time `json:"expires_at"`
}

func (r PrepareRequest) Validate() error {
	if !cr.ValidID(r.ApprovalID) || len(r.ApprovalID) < 16 || r.Generation < 1 || !cr.ValidID(r.VariantID) || r.ExpiresAt.IsZero() {
		return ErrBinding
	}
	switch r.Operation {
	case "MEDIA":
		if r.ItemID != "" || r.MediaApprovalID != "" {
			return ErrBinding
		}
	case "CONTENT":
		if !itemID.MatchString(r.ItemID) || !cr.ValidID(r.MediaApprovalID) || len(r.MediaApprovalID) < 16 || r.MediaApprovalID == r.ApprovalID {
			return ErrBinding
		}
	case "CREATE":
		if r.ItemID != "" || !cr.ValidID(r.MediaApprovalID) || len(r.MediaApprovalID) < 16 || r.MediaApprovalID == r.ApprovalID {
			return ErrBinding
		}
	case "PRICE", "STOCK", "PAUSE", "RESUME":
		if r.MediaApprovalID != "" {
			return ErrBinding
		}
		if !itemID.MatchString(r.ItemID) {
			return ErrBinding
		}
	default:
		return ErrBinding
	}
	return nil
}

type Intent struct {
	PrepareRequest
	MediaSHA256     string          `json:"media_sha256,omitempty"`
	TenantID        string          `json:"tenant_id"`
	OrganizationID  string          `json:"organization_id"`
	ProfileSHA256   string          `json:"profile_sha256"`
	SourceSHA256    string          `json:"source_sha256"`
	SellerID        string          `json:"seller_id"`
	SKU             string          `json:"sku"`
	BeforeSHA256    string          `json:"before_sha256"`
	Quantity        int64           `json:"quantity"`
	PriceMinorUnits int64           `json:"price_minor_units,string"`
	StockHorizon    time.Time       `json:"stock_horizon"`
	Body            json.RawMessage `json:"body"`
	Method          string          `json:"method"`
	Path            string          `json:"path"`
	StockVersion    string          `json:"stock_version,omitempty"`
	Expected        Item            `json:"expected"`
}

func (r Intent) Message() channels.Message {
	_, h, _ := cr.Canonical(r)
	return channels.Message{ChannelCode: "mercadolibre_catalog", TenantID: r.TenantID, ExternalID: r.SellerID, ThreadID: r.SKU, Direction: channels.DirectionOut, DeliveryKey: r.ApprovalID, Text: h}
}

type Picture struct {
	Source    string `json:"source,omitempty"`
	ID        string `json:"id,omitempty"`
	SecureURL string `json:"secure_url,omitempty"`
}
type Item struct {
	SoldQuantity  *int64      `json:"sold_quantity,omitempty"`
	ID            string      `json:"id"`
	SellerID      int64       `json:"seller_id"`
	SiteID        string      `json:"site_id"`
	CategoryID    string      `json:"category_id"`
	ListingTypeID string      `json:"listing_type_id"`
	BuyingMode    string      `json:"buying_mode"`
	FamilyName    string      `json:"family_name"`
	UserProductID string      `json:"user_product_id"`
	SKU           string      `json:"seller_custom_field"`
	Currency      string      `json:"currency_id"`
	Quantity      int64       `json:"available_quantity"`
	Status        string      `json:"status"`
	Attributes    []Attribute `json:"attributes"`
	Pictures      []Picture   `json:"pictures"`
	Tags          []string    `json:"tags"`
}
type Observation struct {
	Item         Item   `json:"item"`
	PriceMinor   string `json:"price_minor"`
	Quantity     int64  `json:"quantity"`
	StockVersion string `json:"stock_version"`
	Exists       bool   `json:"exists"`
}

func (o Observation) SHA256() string { _, h, _ := cr.Canonical(o); return h }
func AttributesEqual(want, got []Attribute) bool {
	for _, a := range want {
		found := false
		for _, b := range got {
			if a.ID == b.ID && (a.ValueID != "" && a.ValueID == b.ValueID || a.ValueName != "" && a.ValueName == b.ValueName) {
				if found {
					return false
				}
				found = true
			}
		}
		if !found {
			return false
		}
	}
	return true
}
func Build(p Profile, pub cr.PublicDocument, r PrepareRequest, quantity int64, horizon time.Time, before Observation, contentPicture ...string) (Intent, error) {
	if r.Operation == "CONTENT" && (len(contentPicture) != 1 || !ValidPictureID(contentPicture[0])) || r.Operation != "CONTENT" && len(contentPicture) != 0 {
		return Intent{}, ErrBinding
	}
	if p.Validate() != nil || r.Validate() != nil || pub.Generation != r.Generation || pub.Currency != p.Currency || !cr.ValidSHA(pub.SourceSHA256) || quantity < 0 || horizon.IsZero() {
		return Intent{}, ErrBinding
	}
	b, e := p.Binding(r.VariantID)
	if e != nil {
		return Intent{}, e
	}
	var variant *cr.Variant
	var model *cr.Model
	for i := range pub.Variants {
		if pub.Variants[i].ID == r.VariantID {
			if variant != nil {
				return Intent{}, ErrBinding
			}
			variant = &pub.Variants[i]
		}
	}
	if variant == nil || variant.AmountMinorUnits <= 0 {
		return Intent{}, ErrBinding
	}
	for i := range pub.Models {
		if pub.Models[i].ID == variant.ModelID {
			if model != nil {
				return Intent{}, ErrBinding
			}
			model = &pub.Models[i]
		}
	}
	if model == nil || !cr.ValidSHA(model.Media.SHA256) || !cr.ValidText(model.DisplayName, 256) {
		return Intent{}, ErrBinding
	}
	attrs := append([]Attribute(nil), b.Attributes...)
	sort.Slice(attrs, func(i, j int) bool { return attrs[i].ID < attrs[j].ID })
	expected := Item{ID: r.ItemID, SiteID: p.SiteID, CategoryID: b.CategoryID, ListingTypeID: b.ListingTypeID, BuyingMode: "buy_it_now", FamilyName: model.DisplayName, SKU: b.SKU, Currency: p.Currency, Quantity: quantity, Attributes: attrs, Pictures: []Picture{{Source: p.MediaOrigin + "/api/public/catalog/media/" + model.Media.SHA256}}}
	expected.UserProductID = before.Item.UserProductID
	out := Intent{PrepareRequest: r, TenantID: p.TenantID, OrganizationID: p.OrganizationID, ProfileSHA256: p.SHA256(), SourceSHA256: pub.SourceSHA256, SellerID: p.SellerID, SKU: b.SKU, BeforeSHA256: before.SHA256(), Quantity: quantity, StockHorizon: horizon, Expected: expected, Method: "PUT", Path: "/items/" + r.ItemID}
	major := json.Number(sv.Major(variant.AmountMinorUnits, p.CurrencyDigits))
	var body any
	switch r.Operation {
	case "CONTENT":
		if before.Item.SoldQuantity == nil || *before.Item.SoldQuantity != 0 {
			return Intent{}, ErrBinding
		}
		out.MediaSHA256 = model.Media.SHA256
		out.Expected.Pictures = []Picture{{ID: contentPicture[0]}}
		body = contentBody(out.Expected)
	case "MEDIA", "CREATE":
		out.MediaSHA256 = model.Media.SHA256
		out.Method = "POST"
		if before.Exists {
			return Intent{}, ErrBinding
		}
		if r.Operation == "MEDIA" {
			out.Path = "/pictures/items/upload"
			body = map[string]string{"catalog_png_sha256": out.MediaSHA256}
		} else {
			if b.StockMode != "item" || len(before.Item.Pictures) != 1 || !ValidPictureID(before.Item.Pictures[0].ID) {
				return Intent{}, ErrBinding
			}
			out.Path = "/items"
			out.Expected.Pictures = []Picture{{ID: before.Item.Pictures[0].ID}}
			body = creationBody(out.Expected, major)
		}
	case "PRICE":
		body = map[string]any{"price": major}
	case "STOCK":
		switch b.StockMode {
		case "item":
			body = map[string]any{"available_quantity": quantity}
		case "selling_address":
			out.Path = "/user-products/" + before.Item.UserProductID + "/stock/type/selling_address"
			body = map[string]any{"quantity": quantity}
		case "seller_warehouse":
			out.Path = "/user-products/" + before.Item.UserProductID + "/stock/type/seller_warehouse"
			body = map[string]any{"locations": []map[string]any{{"store_id": b.StoreID, "network_node_id": b.NetworkNodeID, "quantity": quantity}}}
		}
		if b.StockMode != "item" {
			if !upID.MatchString(before.Item.UserProductID) || before.StockVersion == "" {
				return Intent{}, ErrBinding
			}
			out.StockVersion = before.StockVersion
		}
	case "PAUSE":
		body = map[string]string{"status": "paused"}
	case "RESUME":
		body = map[string]string{"status": "active"}
	}
	if r.Operation != "MEDIA" && r.Operation != "CREATE" && (!before.Exists || before.Item.ID != r.ItemID) {
		return Intent{}, ErrBinding
	}
	out.PriceMinorUnits = variant.AmountMinorUnits
	out.Body, _, e = cr.Canonical(body)
	return out, e
}
