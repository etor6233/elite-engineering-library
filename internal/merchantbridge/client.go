package merchantbridge

import (
	"context"
	cr "elite.local/enterprise/internal/catalogrelease"
	"elite.local/enterprise/internal/outbounddelivery"
	"encoding/json"
	"strconv"
	"time"
)

type Executor interface {
	Execute(context.Context, string, Intent) (Result, error)
}
type Client struct {
	profile  Profile
	executor Executor
}

func NewClient(p Profile, x Executor) (*Client, error) {
	if p.Validate() != nil || x == nil {
		return nil, ErrBinding
	}
	raw, _ := json.Marshal(p)
	var frozen Profile
	if json.Unmarshal(raw, &frozen) != nil {
		return nil, ErrBinding
	}
	return &Client{frozen, x}, nil
}

type Attributes struct {
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
type Observation struct {
	Name          string          `json:"name"`
	Product       string          `json:"product"`
	DataSource    string          `json:"data_source"`
	OfferID       string          `json:"offer_id"`
	Language      string          `json:"content_language"`
	FeedLabel     string          `json:"feed_label"`
	Version       string          `json:"version_number"`
	Attributes    Attributes      `json:"product_attributes"`
	ProductStatus json.RawMessage `json:"product_status"`
}

func (c *Client) valid(r Intent) bool {
	b, e := c.profile.Binding(r.VariantID)
	return e == nil && r.PrepareRequest.Validate() == nil && r.TenantID == c.profile.TenantID && r.OrganizationID == c.profile.OrganizationID && r.ProfileSHA256 == c.profile.SHA256() && r.AccountID == c.profile.AccountID && r.DataSource == c.profile.DataSource() && r.Product.OfferID == b.OfferID && r.Resource == c.profile.Resource(b.OfferID) && cr.ValidSHA(r.SourceSHA256) && cr.ValidSHA(r.MediaSHA256)
}
func attributesMatch(p Product, a Attributes) bool {
	if a.Title != p.Title || a.Description != p.Description || a.Link != p.Link || a.ImageLink != p.ImageLink || a.Availability != p.Availability || a.Condition != p.Condition || a.Price != p.Price || a.Brand != p.Brand || len(a.GTINs) != len(p.GTINs) {
		return false
	}
	// Official GTIN order is not an identity guarantee; compare exact membership.
	seen := map[string]bool{}
	for _, g := range a.GTINs {
		if seen[g] {
			return false
		}
		seen[g] = true
	}
	for _, g := range p.GTINs {
		if !seen[g] {
			return false
		}
	}
	return true
}
func (c *Client) confirm(r Intent, result Result) (Observation, error) {
	var o Observation
	if result.State != "OBSERVED" || result.Resource != r.Resource || Hash(result.Response) != result.ResponseSHA256 || json.Unmarshal(result.Response, &o) != nil || o.OfferID != r.Product.OfferID || o.Language != r.Product.Language || o.FeedLabel != r.Product.FeedLabel || o.Version != strconv.FormatInt(r.Generation, 10) || !attributesMatch(r.Product, o.Attributes) {
		return o, ErrUnknown
	}
	if result.Operation == "INSERT" {
		expected := "accounts/" + r.AccountID + "/productInputs/" + r.Product.Language + "~" + r.Product.FeedLabel + "~" + r.Product.OfferID
		if o.Name != expected || o.Product != r.Resource {
			return o, ErrUnknown
		}
	} else if result.Operation == "GET" {
		if o.Name != r.Resource || o.DataSource != r.DataSource {
			return o, ErrUnknown
		}
	} else {
		return o, ErrBinding
	}
	return o, nil
}

type Record func(context.Context, Result) error

func (c *Client) Insert(ctx context.Context, r Intent, record Record) (outbounddelivery.Receipt, error) {
	if !c.valid(r) || !r.ExpiresAt.After(time.Now()) || record == nil {
		return outbounddelivery.Receipt{}, ErrBinding
	}
	result, e := c.executor.Execute(ctx, "INSERT", r)
	if e != nil {
		return outbounddelivery.Receipt{}, e
	}
	if result.State == "REJECTED" {
		_, h, _ := cr.Canonical(result)
		return outbounddelivery.Receipt{}, outbounddelivery.NewTerminalFailure("GOOGLE_MERCHANT_INSERT_REJECTED", h, nil)
	}
	if _, e = c.confirm(r, result); e != nil {
		return outbounddelivery.Receipt{}, e
	}
	if e = record(ctx, result); e != nil {
		return outbounddelivery.Receipt{}, ErrUnknown
	}
	return outbounddelivery.Receipt{ProviderMessageID: r.Resource, EvidenceSHA256: result.ResponseSHA256, AcceptedAt: time.Now().UTC()}, nil
}
func (c *Client) Observe(ctx context.Context, r Intent) (Result, error) {
	if !c.valid(r) {
		return Result{}, ErrBinding
	}
	result, e := c.executor.Execute(ctx, "GET", r)
	if e != nil {
		return result, e
	}
	// Preserve processing divergence/status separately; caller does not infer
	// provider approval or claim the input write refreshed freshness from GET.
	if _, e = c.confirm(r, result); e != nil {
		return result, e
	}
	return result, nil
}
