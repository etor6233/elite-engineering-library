package postgres

// AUTHORED exact projection/hash glue. The selected source-derived writer owns
// the gift/discount rules; all readers share this representation under order lock.
import (
	"context"
	"crypto/sha256"
	"elite.local/enterprise/internal/commerce"
	"encoding/hex"
	"encoding/json"
	"github.com/jackc/pgx/v5"
)

type orderFundingReader interface {
	QueryRow(context.Context, string, ...any) pgx.Row
}
type OrderFundingSnapshot struct {
	TenantID            string `json:"tenant_id"`
	OrganizationID      string `json:"organization_id"`
	OrderID             string `json:"order_id"`
	Currency            string `json:"currency"`
	GrossMinor          int64  `json:"gross_minor_units"`
	GiftMinor           int64  `json:"gift_minor_units"`
	DiscountMinor       int64  `json:"discount_minor_units"`
	ProviderMinor       int64  `json:"provider_minor_units"`
	ContributionsSHA256 string `json:"contributions_sha256"`
}

func fundingHash(raw []byte) string           { h := sha256.Sum256(raw); return hex.EncodeToString(h[:]) }
func (v OrderFundingSnapshot) SHA256() string { raw, _ := json.Marshal(v); return fundingHash(raw) }
func readOrderFunding(ctx context.Context, q orderFundingReader, tenant, org, order, currency string, gross int64) (OrderFundingSnapshot, error) {
	v := OrderFundingSnapshot{TenantID: tenant, OrganizationID: org, OrderID: order}
	e := q.QueryRow(ctx, `select currency,gross_minor_units,gift_minor_units,discount_minor_units,provider_due_minor_units,contributions_sha256 from payment.order_funding where tenant_id=$1 and organization_id=$2 and order_id=$3`, tenant, org, order).Scan(&v.Currency, &v.GrossMinor, &v.GiftMinor, &v.DiscountMinor, &v.ProviderMinor, &v.ContributionsSHA256)
	if e != nil {
		return v, e
	}
	hash, e := hex.DecodeString(v.ContributionsSHA256)
	if currency != v.Currency || gross != v.GrossMinor || gross <= 0 || v.GiftMinor < 0 || v.GiftMinor > gross || v.DiscountMinor < 0 || v.DiscountMinor > gross-v.GiftMinor || v.ProviderMinor < 0 || v.ProviderMinor != gross-v.GiftMinor-v.DiscountMinor || e != nil || len(hash) != 32 || hex.EncodeToString(hash) != v.ContributionsSHA256 {
		return v, commerce.ErrConflict
	}
	return v, nil
}
func orderProviderDue(ctx context.Context, q orderFundingReader, tenant, org, order, currency string, gross int64) (int64, error) {
	v, e := readOrderFunding(ctx, q, tenant, org, order, currency, gross)
	return v.ProviderMinor, e
}
