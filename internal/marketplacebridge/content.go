package marketplacebridge

// AUTHORED field/scope mapping for the fixed official User Products contract.
import (
	"context"
	cr "elite.local/enterprise/internal/catalogrelease"
	"encoding/json"
	"net/url"
	"time"
)

func contentBody(i Item) any {
	return struct {
		Family     string      `json:"family_name"`
		Attributes []Attribute `json:"attributes"`
		Pictures   []Picture   `json:"pictures"`
	}{i.FamilyName, i.Attributes, i.Pictures}
}
func validateContent(p Profile, r Intent) error {
	b, e := p.Binding(r.VariantID)
	x := r.Expected
	if e != nil || r.Operation != "CONTENT" || !cr.ValidSHA(r.MediaSHA256) || !upID.MatchString(x.UserProductID) || x.ID != r.ItemID || x.SKU != b.SKU || !cr.ValidText(x.FamilyName, 256) || x.CategoryID != b.CategoryID || !AttributesEqual(b.Attributes, x.Attributes) || len(b.Attributes) != len(x.Attributes) || len(x.Pictures) != 1 || !ValidPictureID(x.Pictures[0].ID) || x.Pictures[0].Source != "" || x.Pictures[0].SecureURL != "" {
		return ErrBinding
	}
	return nil
}
func (c *Client) SingleUnsoldItem(ctx context.Context, item Item) error {
	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()
	if !itemID.MatchString(item.ID) || !upID.MatchString(item.UserProductID) || item.SoldQuantity == nil || *item.SoldQuantity != 0 {
		return ErrBinding
	}
	status, raw, _, e := c.call(ctx, "GET", "/users/"+c.profile.SellerID+"/items/search?user_product_id="+url.QueryEscape(item.UserProductID)+"&limit=2", nil, "")
	var result struct {
		Seller  string   `json:"seller_id"`
		Results []string `json:"results"`
		Paging  struct {
			Total  int `json:"total"`
			Offset int `json:"offset"`
		} `json:"paging"`
	}
	if e != nil || status != 200 || json.Unmarshal(raw, &result) != nil || result.Seller != c.profile.SellerID || result.Paging.Total != 1 || result.Paging.Offset != 0 || len(result.Results) != 1 || result.Results[0] != item.ID {
		return ErrBinding
	}
	return nil
}
