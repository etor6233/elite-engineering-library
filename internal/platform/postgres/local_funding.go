package postgres

// AUTHORED immutable observation of already approved contributions. This is
// not a provider capture, cash movement, general ledger or gift-card rule.
import (
	"bytes"
	"context"
	"elite.local/enterprise/internal/commerce"
	"encoding/json"
	"io"
	"time"
)

const LocalFundingEffect = "ORDER_FULLY_FUNDED_BY_STORED_VALUE"

type LocalFundingReceipt struct {
	ID               string               `json:"funding_receipt_id"`
	RequestKey       string               `json:"request_key"`
	RequestSHA256    string               `json:"request_sha256"`
	RequestedBy      string               `json:"requested_by"`
	Allocation       OrderFundingSnapshot `json:"allocation"`
	AllocationSHA256 string               `json:"allocation_sha256"`
	ObservedAt       time.Time            `json:"observed_at"`
	Effect           string               `json:"effect"`
}
type LocalFundingResult struct {
	Receipt LocalFundingReceipt `json:"receipt"`
	SHA256  string              `json:"receipt_sha256"`
}

// Order is locked by the caller; contributions cannot change during validation.
func readLocalFunding(ctx context.Context, q orderFundingReader, tenant, org, order, id, evidence, currency string, gross int64) (LocalFundingResult, error) {
	var out LocalFundingResult
	var raw []byte
	var storedAllocation string
	var amount, gift, discount int64
	var at time.Time
	e := q.QueryRow(ctx, `select receipt_raw,receipt_sha256,allocation_sha256,gross_minor_units,gift_minor_units,discount_minor_units,created_at from payment.local_funding_evidence where tenant_id=$1 and organization_id=$2 and order_id=$3 and funding_id=$4 and currency=$5`, tenant, org, order, id, currency).Scan(&raw, &out.SHA256, &storedAllocation, &amount, &gift, &discount, &at)
	if e != nil {
		return out, e
	}
	if len(raw) > 65536 || fundingHash(raw) != out.SHA256 || out.SHA256 != evidence {
		return out, commerce.ErrConflict
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	if decoder.Decode(&out.Receipt) != nil || decoder.Decode(new(any)) != io.EOF {
		return out, commerce.ErrConflict
	}
	current, e := readOrderFunding(ctx, q, tenant, org, order, currency, gross)
	if e != nil {
		return out, e
	}
	r := out.Receipt
	if current.ProviderMinor != 0 || current.GiftMinor+current.DiscountMinor <= 0 || r.ID != id || r.Effect != LocalFundingEffect || !r.ObservedAt.Equal(at) || r.Allocation != current || r.AllocationSHA256 != current.SHA256() || r.AllocationSHA256 != storedAllocation || amount != gross || gift != current.GiftMinor || discount != current.DiscountMinor {
		return out, commerce.ErrConflict
	}
	return out, nil
}
