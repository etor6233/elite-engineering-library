package merchantbridge

// AUTHORED provider representation and source binding. Catalog content, prices,
// inventory ATP, approvals and the outbound fence remain their original owners.
import (
	"bytes"
	"elite.local/enterprise/internal/approval"
	cr "elite.local/enterprise/internal/catalogrelease"
	"elite.local/enterprise/internal/channels"
	"encoding/json"
	"errors"
	"io"
	"net/url"
	"regexp"
	"strconv"
	"time"
)

var ErrBinding = errors.New("merchant: invalid or changed source binding")
var ErrUnknown = errors.New("merchant: provider result not confirmed")
var number = regexp.MustCompile(`^[0-9]{3,20}$`)
var offer = regexp.MustCompile(`^[A-Za-z0-9._-]{1,50}$`)

type Binding struct {
	VariantID string   `json:"variant_id"`
	OfferID   string   `json:"offer_id"`
	Condition string   `json:"condition"`
	Brand     string   `json:"brand"`
	GTINs     []string `json:"gtins"`
}
type Profile struct {
	Schema             string    `json:"schema"`
	TenantID           string    `json:"tenant_id"`
	OrganizationID     string    `json:"organization_id"`
	AccountID          string    `json:"account_id"`
	DataSourceID       string    `json:"data_source_id"`
	Language           string    `json:"content_language"`
	FeedLabel          string    `json:"feed_label"`
	Currency           string    `json:"currency"`
	CurrencyDigits     int       `json:"currency_digits"`
	Origin             string    `json:"origin"`
	RefreshCadenceDays int       `json:"refresh_cadence_days"`
	Bindings           []Binding `json:"bindings"`
}

func (p Profile) Validate() error {
	if p.RefreshCadenceDays < 1 || p.RefreshCadenceDays > 30 {
		return ErrBinding
	}
	u, e := url.Parse(p.Origin)
	if p.Schema != "elite-connected-merchant/v1" || !cr.ValidID(p.TenantID) || !cr.ValidID(p.OrganizationID) || !number.MatchString(p.AccountID) || !number.MatchString(p.DataSourceID) || !regexp.MustCompile(`^[a-z]{2}$`).MatchString(p.Language) || !regexp.MustCompile(`^[A-Z0-9-]{2,20}$`).MatchString(p.FeedLabel) || !regexp.MustCompile(`^[A-Z]{3}$`).MatchString(p.Currency) || p.CurrencyDigits < 0 || p.CurrencyDigits > 6 || e != nil || u.Scheme != "https" || u.Host == "" || u.User != nil || u.Path != "" || u.RawQuery != "" || u.Fragment != "" || len(p.Origin) > 512 || len(p.Bindings) < 1 || len(p.Bindings) > 32 {
		return ErrBinding
	}
	variants, offers := map[string]bool{}, map[string]bool{}
	for _, b := range p.Bindings {
		if !cr.ValidID(b.VariantID) || !offer.MatchString(b.OfferID) || variants[b.VariantID] || offers[b.OfferID] || b.Condition != "NEW" && b.Condition != "USED" && b.Condition != "REFURBISHED" || len(b.Brand) > 70 || b.Brand != "" && !cr.ValidText(b.Brand, 70) || len(b.GTINs) > 10 {
			return ErrBinding
		}
		variants[b.VariantID] = true
		offers[b.OfferID] = true
		seen := map[string]bool{}
		for _, g := range b.GTINs {
			if !regexp.MustCompile(`^(?:[0-9]{8}|[0-9]{12,14})$`).MatchString(g) || seen[g] {
				return ErrBinding
			}
			seen[g] = true
		}
	}
	if _, _, e = cr.Canonical(p); e != nil {
		return ErrBinding
	}
	return nil
}
func (p Profile) SHA256() string { _, h, _ := cr.Canonical(p); return h }
func (p Profile) DataSource() string {
	return "accounts/" + p.AccountID + "/dataSources/" + p.DataSourceID
}
func (p Profile) Binding(id string) (Binding, error) {
	for _, b := range p.Bindings {
		if b.VariantID == id {
			return b, nil
		}
	}
	return Binding{}, ErrBinding
}
func (p Profile) Resource(offerID string) string {
	return "accounts/" + p.AccountID + "/products/" + p.Language + "~" + p.FeedLabel + "~" + offerID
}

type PrepareRequest struct {
	ApprovalID string    `json:"approval_id"`
	Generation int64     `json:"generation,string"`
	VariantID  string    `json:"variant_id"`
	ExpiresAt  time.Time `json:"expires_at"`
}

func (r PrepareRequest) Validate() error {
	if !cr.ValidID(r.ApprovalID) || len(r.ApprovalID) < 16 || r.Generation < 1 || !cr.ValidID(r.VariantID) || r.ExpiresAt.IsZero() {
		return ErrBinding
	}
	return nil
}

type Price struct {
	AmountMicros string `json:"amount_micros"`
	Currency     string `json:"currency_code"`
}
type Product struct {
	OfferID      string   `json:"offer_id"`
	Language     string   `json:"content_language"`
	FeedLabel    string   `json:"feed_label"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Link         string   `json:"link"`
	ImageLink    string   `json:"image_link"`
	Availability string   `json:"availability"`
	Condition    string   `json:"condition"`
	Price        Price    `json:"price"`
	Brand        string   `json:"brand"`
	GTINs        []string `json:"gtins"`
}
type Intent struct {
	PrepareRequest
	TenantID       string    `json:"tenant_id"`
	OrganizationID string    `json:"organization_id"`
	ProfileSHA256  string    `json:"profile_sha256"`
	SourceSHA256   string    `json:"source_sha256"`
	MediaSHA256    string    `json:"media_sha256"`
	AccountID      string    `json:"account_id"`
	DataSource     string    `json:"data_source"`
	Resource       string    `json:"resource"`
	StockHorizon   time.Time `json:"stock_horizon"`
	Quantity       int64     `json:"quantity,string"`
	Product        Product   `json:"product"`
}

func (r Intent) Message() channels.Message {
	_, h, _ := cr.Canonical(r)
	return channels.Message{ChannelCode: "google_merchant", TenantID: r.TenantID, ExternalID: r.AccountID, ThreadID: r.Product.OfferID, Direction: channels.DirectionOut, DeliveryKey: r.ApprovalID, Text: h}
}
func Decode(raw []byte, v any) error {
	if len(raw) == 0 || len(raw) > 32768 {
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
func Build(p Profile, pub cr.PublicDocument, r PrepareRequest, quantity int64, horizon time.Time) (Intent, error) {
	if p.Validate() != nil || r.Validate() != nil || pub.Generation != r.Generation || pub.Currency != p.Currency || !cr.ValidSHA(pub.SourceSHA256) || quantity < 0 || horizon.IsZero() {
		return Intent{}, ErrBinding
	}
	b, e := p.Binding(r.VariantID)
	if e != nil {
		return Intent{}, e
	}
	var v *cr.Variant
	var m *cr.Model
	for i := range pub.Variants {
		if pub.Variants[i].ID == r.VariantID {
			if v != nil {
				return Intent{}, ErrBinding
			}
			v = &pub.Variants[i]
		}
	}
	if v == nil || v.AmountMinorUnits <= 0 {
		return Intent{}, ErrBinding
	}
	for i := range pub.Models {
		if pub.Models[i].ID == v.ModelID {
			if m != nil {
				return Intent{}, ErrBinding
			}
			m = &pub.Models[i]
		}
	}
	if m == nil || !cr.ValidSHA(m.Media.SHA256) || !cr.ValidText(m.DisplayName, 150) {
		return Intent{}, ErrBinding
	}
	var specification map[string]json.RawMessage
	var description string
	if json.Unmarshal(m.Specification, &specification) != nil || json.Unmarshal(specification["merchant_description"], &description) != nil || !cr.ValidText(description, 5000) {
		return Intent{}, ErrBinding
	}
	// Representation only: convert original integer minor units into Google
	// integer micros without floating point, rounding, tax or exchange rules.
	scale := int64(1)
	for i := p.CurrencyDigits; i < 6; i++ {
		scale *= 10
	}
	if v.AmountMinorUnits > 1000000000000000000/scale {
		return Intent{}, ErrBinding
	}
	amount := v.AmountMinorUnits * scale
	availability := "OUT_OF_STOCK"
	if quantity > 0 {
		availability = "IN_STOCK"
	}
	gtins := append([]string{}, b.GTINs...)
	product := Product{b.OfferID, p.Language, p.FeedLabel, m.DisplayName, description, m.CanonicalURL, p.Origin + "/api/public/catalog/media/" + m.Media.SHA256, availability, b.Condition, Price{strconv.FormatInt(amount, 10), p.Currency}, b.Brand, gtins}
	link, e := url.Parse(product.Link)
	origin, _ := url.Parse(p.Origin)
	if e != nil || link.Scheme != "https" || link.Host != origin.Host || link.User != nil || link.RawQuery != "" || link.Fragment != "" {
		return Intent{}, ErrBinding
	}
	out := Intent{PrepareRequest: r, TenantID: p.TenantID, OrganizationID: p.OrganizationID, ProfileSHA256: p.SHA256(), SourceSHA256: pub.SourceSHA256, MediaSHA256: m.Media.SHA256, AccountID: p.AccountID, DataSource: p.DataSource(), Resource: p.Resource(b.OfferID), StockHorizon: horizon, Quantity: quantity, Product: product}
	if _, _, e = cr.Canonical(out); e != nil {
		return Intent{}, ErrBinding
	}
	return out, nil
}
