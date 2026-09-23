package postgres

// AUTHORED final observation and request recovery; all applied contributions
// already exist under immutable source-derived, separately approved operations.
import (
	"context"
	"elite.local/enterprise/internal/platform/identity"
	sv "elite.local/enterprise/internal/storedvaluebridge"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5"
	"time"
)

type FinalizeStoredValueFunding struct {
	OrganizationID       string `json:"organization_id"`
	OrderID              string `json:"order_id"`
	ExpectedOrderVersion int64  `json:"expected_order_version"`
	RequestKey           string `json:"request_key"`
}

func storedFundingRequestHash(p identity.Principal, c FinalizeStoredValueFunding) string {
	raw, _ := json.Marshal(struct {
		Tenant, Actor string
		Command       FinalizeStoredValueFunding
	}{p.TenantID, p.Subject, c})
	return fundingHash(raw)
}
func readStoredFundingRequest(ctx context.Context, q orderFundingReader, tenant, org, actor, key, hash string) (LocalFundingResult, error) {
	var out LocalFundingResult
	var raw []byte
	var saved string
	e := q.QueryRow(ctx, `select receipt_raw,receipt_sha256,request_sha256 from payment.local_funding_receipt where tenant_id=$1 and organization_id=$2 and requested_by_subject=$3 and request_key=$4`, tenant, org, actor, key).Scan(&raw, &out.SHA256, &saved)
	if e != nil {
		return out, e
	}
	if hash != "" && saved != hash {
		return out, sv.ErrBinding
	}
	if fundingHash(raw) != out.SHA256 || json.Unmarshal(raw, &out.Receipt) != nil || out.Receipt.RequestedBy != actor || out.Receipt.RequestKey != key || out.Receipt.RequestSHA256 != saved || out.Receipt.Effect != LocalFundingEffect || out.Receipt.Allocation.TenantID != tenant || out.Receipt.Allocation.OrganizationID != org || out.Receipt.AllocationSHA256 != out.Receipt.Allocation.SHA256() {
		return out, sv.ErrBinding
	}
	return out, nil
}
func (s *StoredValue) FundingResult(ctx context.Context, p identity.Principal, key string) (LocalFundingResult, error) {
	if !s.allowed(p, "stored_value:read") || len(key) < 16 || !sv.ValidID(key) {
		return LocalFundingResult{}, sv.ErrBinding
	}
	tenant, org := s.profile.Scope()
	return readStoredFundingRequest(ctx, s.pool, tenant, org, p.Subject, key, "")
}
func (s *StoredValue) FinalizeFunding(ctx context.Context, p identity.Principal, c FinalizeStoredValueFunding) (LocalFundingResult, bool, error) {
	var empty LocalFundingResult
	if !s.allowed(p, "stored_value:fund") || !sv.ValidID(c.OrderID) || !sv.ValidID(c.RequestKey) || len(c.RequestKey) < 16 || c.ExpectedOrderVersion < 1 {
		return empty, false, sv.ErrBinding
	}
	tenant, org := s.profile.Scope()
	if org != c.OrganizationID {
		return empty, false, sv.ErrBinding
	}
	hash := storedFundingRequestHash(p, c)
	tx, e := s.pool.Begin(ctx)
	if e != nil {
		return empty, false, e
	}
	defer tx.Rollback(ctx)
	if old, e := readStoredFundingRequest(ctx, tx, tenant, org, p.Subject, c.RequestKey, hash); e == nil {
		return old, true, tx.Commit(ctx)
	} else if !errors.Is(e, pgx.ErrNoRows) {
		return empty, false, e
	}
	var currency, state string
	var gross, version int64
	e = tx.QueryRow(ctx, `select currency,total_minor_units,state,version from sales.customer_order where tenant_id=$1 and organization_id=$2 and order_id=$3 for update`, tenant, org, c.OrderID).Scan(&currency, &gross, &state, &version)
	if e != nil {
		return empty, false, e
	}
	if version != c.ExpectedOrderVersion || (state != "placed" && state != "confirmed" && state != "allocated") {
		return empty, false, sv.ErrBinding
	}
	a, e := readOrderFunding(ctx, tx, tenant, org, c.OrderID, currency, gross)
	if e != nil {
		return empty, false, e
	}
	if a.ProviderMinor != 0 || a.GiftMinor+a.DiscountMinor <= 0 {
		return empty, false, sv.ErrBinding
	}
	var active bool
	e = tx.QueryRow(ctx, `select exists(select 1 from payment.payment_attempt where tenant_id=$1 and order_id=$2 and state not in ('failed','refunded'))`, tenant, c.OrderID).Scan(&active)
	if e != nil {
		return empty, false, e
	}
	if active {
		return empty, false, sv.ErrBinding
	}
	var now time.Time
	if e = tx.QueryRow(ctx, `select clock_timestamp()`).Scan(&now); e != nil {
		return empty, false, e
	}
	rec := LocalFundingReceipt{ID: sv.StableID("funding", tenant, c.RequestKey), RequestKey: c.RequestKey, RequestSHA256: hash, RequestedBy: p.Subject, Allocation: a, AllocationSHA256: a.SHA256(), ObservedAt: now.UTC(), Effect: LocalFundingEffect}
	raw, e := json.Marshal(rec)
	if e != nil {
		return empty, false, e
	}
	receiptHash := fundingHash(raw)
	inserted, e := tx.Exec(ctx, `insert into payment.local_funding_receipt(tenant_id,funding_id,request_key,request_sha256,requested_by_subject,receipt_raw,order_id,organization_id,currency,gross_minor_units,gift_minor_units,discount_minor_units,provider_minor_units,allocation_sha256,receipt_sha256,receipt,created_at)values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,0,$13,$14,$15,$16) on conflict(tenant_id,request_key)do nothing`, tenant, rec.ID, c.RequestKey, hash, p.Subject, raw, c.OrderID, org, currency, gross, a.GiftMinor, a.DiscountMinor, a.SHA256(), receiptHash, raw, now)
	if e != nil {
		return empty, false, e
	}
	if inserted.RowsAffected() == 0 {
		old, e := readStoredFundingRequest(ctx, tx, tenant, org, p.Subject, c.RequestKey, hash)
		if e != nil {
			return empty, false, e
		}
		return old, true, tx.Commit(ctx)
	}
	if _, e = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)values($1,gen_random_uuid(),'order-funding',$2,1,'order.local-funding-observed',1,$3,$4)`, tenant, rec.ID, now, raw); e != nil {
		return empty, false, e
	}
	if e = tx.Commit(ctx); e != nil {
		return empty, false, e
	}
	return LocalFundingResult{Receipt: rec, SHA256: receiptHash}, false, nil
}
