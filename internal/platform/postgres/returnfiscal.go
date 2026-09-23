package postgres

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"elite.local/enterprise/internal/fiscal"
	"elite.local/enterprise/internal/returneffects"
	"elite.local/enterprise/internal/returnfiscal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReturnFiscal struct {
	pool    *pgxpool.Pool
	effects *ReturnEffects
	fiscal  *Fiscal
}

func NewReturnFiscal(pool *pgxpool.Pool) *ReturnFiscal {
	return &ReturnFiscal{pool: pool, effects: NewReturnEffects(pool), fiscal: NewFiscal(pool)}
}
func (r *ReturnFiscal) Claim(ctx context.Context, owner, worker, token string, lease time.Duration) (*returneffects.Work, error) {
	return r.effects.Claim(ctx, owner, worker, token, lease)
}
func (r *ReturnFiscal) Finish(ctx context.Context, work returneffects.Work, worker string, c returneffects.Completion) error {
	return r.effects.Finish(ctx, work, worker, c)
}

func stableReturnFiscalUUID(scope, value string) string {
	sum := sha256.Sum256([]byte(scope + "\x00" + value))
	b := sum[:16]
	b[6] = (b[6] & 0x0f) | 0x50
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
func stableReturnFiscalHash(scope, value string) string {
	sum := sha256.Sum256([]byte(scope + "\x00" + value))
	return hex.EncodeToString(sum[:])
}

func (r *ReturnFiscal) RequestOrObserve(ctx context.Context, work returneffects.Work, worker string) (returnfiscal.Result, error) {
	if work.EffectKind != "fiscal" || work.OwnerContext != "fiscal" {
		return returnfiscal.Result{}, returnfiscal.ErrConflict
	}
	var originalID string
	var decidedAt time.Time
	err := r.pool.QueryRow(ctx, `
select original.invoice_id,d.decided_at
from sales.return_effect_execution x
join sales.return_effect_request e on e.tenant_id=x.tenant_id and e.request_id=x.request_id
join sales.return_disposition d on d.tenant_id=e.tenant_id and d.disposition_id=e.disposition_id and d.customer_remedy='refund'
join sales.return_receipt receipt on receipt.tenant_id=d.tenant_id and receipt.receipt_id=d.receipt_id
join sales.return_effect_request inventory on inventory.tenant_id=e.tenant_id and inventory.disposition_id=e.disposition_id and inventory.effect_kind='inventory'
join sales.return_effect_execution inventory_x on inventory_x.tenant_id=inventory.tenant_id and inventory_x.request_id=inventory.request_id and inventory_x.status='succeeded'
join sales.return_effect_request remedy on remedy.tenant_id=e.tenant_id and remedy.disposition_id=e.disposition_id and remedy.effect_kind='refund'
join sales.return_effect_execution remedy_x on remedy_x.tenant_id=remedy.tenant_id and remedy_x.request_id=remedy.request_id and remedy_x.status='succeeded'
join payment.return_refund refund on refund.tenant_id=remedy.tenant_id and refund.request_id=remedy.request_id and refund.state='succeeded'
join sales.return_effect_request accounting on accounting.tenant_id=e.tenant_id and accounting.disposition_id=e.disposition_id and accounting.effect_kind='accounting'
join sales.return_effect_execution accounting_x on accounting_x.tenant_id=accounting.tenant_id and accounting_x.request_id=accounting.request_id and accounting_x.status='succeeded'
join accounting.return_effect_posting posting on posting.tenant_id=accounting.tenant_id and posting.request_id=accounting.request_id and posting.status='posted'
join fiscal.invoice original on original.tenant_id=receipt.tenant_id and original.organization_id=receipt.organization_id and original.order_id=receipt.order_id and original.status='authorized' and original.voucher_type in(1,6,11)
where x.tenant_id=$1 and x.request_id=$2 and x.status='claimed' and x.claimed_by=$3 and x.claim_token=$4 and x.claimed_until>=clock_timestamp()
and e.effect_kind='fiscal' and e.owner_context='fiscal' and refund.amount_minor_units=original.total_minor_units and refund.currency=original.currency and posting.reversed_minor_units=original.total_minor_units
and 1=(select count(*) from fiscal.invoice candidate where candidate.tenant_id=original.tenant_id and candidate.organization_id=original.organization_id and candidate.order_id=original.order_id and candidate.status='authorized' and candidate.voucher_type in(1,6,11))`, work.TenantID, work.RequestID, worker, work.ClaimToken).Scan(&originalID, &decidedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return returnfiscal.Result{}, returnfiscal.ErrConflict
	}
	if err != nil {
		return returnfiscal.Result{}, err
	}
	original, err := r.fiscal.GetInvoice(ctx, work.TenantID, work.OrganizationID, originalID)
	if err != nil {
		return returnfiscal.Result{}, err
	}
	creditType := map[int]int{1: 3, 6: 8, 11: 13}[original.VoucherType]
	if creditType == 0 {
		return returnfiscal.Result{}, returnfiscal.ErrConflict
	}
	creditID := stableReturnFiscalUUID("credit-invoice", work.TenantID+"/"+work.RequestID)
	credit := original
	credit.TenantID = ""
	credit.ID = creditID
	credit.VoucherType = creditType
	credit.IssuedOn = decidedAt.UTC()
	credit.Status = "queued"
	credit.VoucherNumber = 0
	credit.CAE = ""
	credit.CAEExpiresOn = nil
	credit.Version = 1
	credit.AssociatedVouchers = []fiscal.AssociatedVoucher{{InvoiceID: original.ID}}
	payload, err := json.Marshal(struct {
		RequestID, OriginalInvoiceID, CreditInvoiceID string
		VoucherType                                   int
		IssuedOn                                      string
		Total                                         int64
	}{work.RequestID, original.ID, creditID, creditType, credit.IssuedOn.Format("2006-01-02"), credit.TotalMinorUnits})
	if err != nil {
		return returnfiscal.Result{}, err
	}
	digest := sha256.Sum256(payload)
	requestHash := hex.EncodeToString(digest[:])
	idempotency := stableReturnFiscalHash("credit-idempotency", work.TenantID+"/"+work.RequestID)
	created, _, err := r.fiscal.RequestInvoice(ctx, work.TenantID, idempotency, requestHash, stableReturnFiscalUUID("credit-event", work.TenantID+"/"+work.RequestID), credit)
	if err != nil {
		return returnfiscal.Result{}, err
	}
	inserted, err := r.pool.Exec(ctx, `insert into fiscal.return_credit_note_link(tenant_id,request_id,disposition_id,original_invoice_id,credit_invoice_id,currency,total_minor_units,request_sha256_hex)values($1,$2,$3,$4,$5,$6,$7,$8) on conflict do nothing`, work.TenantID, work.RequestID, work.DispositionID, original.ID, created.ID, created.Currency, created.TotalMinorUnits, requestHash)
	if err != nil {
		return returnfiscal.Result{}, err
	}
	_ = inserted
	var exact int
	err = r.pool.QueryRow(ctx, `select count(*) from fiscal.return_credit_note_link where tenant_id=$1 and request_id=$2 and disposition_id=$3 and original_invoice_id=$4 and credit_invoice_id=$5 and currency=$6 and total_minor_units=$7 and request_sha256_hex=$8`, work.TenantID, work.RequestID, work.DispositionID, original.ID, created.ID, created.Currency, created.TotalMinorUnits, requestHash).Scan(&exact)
	if err != nil {
		return returnfiscal.Result{}, err
	}
	if exact != 1 {
		return returnfiscal.Result{}, returnfiscal.ErrConflict
	}
	result := returnfiscal.Result{RequestID: work.RequestID, OriginalInvoiceID: original.ID, CreditInvoiceID: created.ID, Status: created.Status}
	if created.Status == "rejected" {
		return result, returnfiscal.ErrRejected
	}
	if created.Status != "authorized" {
		return result, returnfiscal.ErrPending
	}
	resultBytes, _ := json.Marshal(struct {
		RequestID, CreditInvoiceID, CAE string
		VoucherNumber                   int64
	}{work.RequestID, created.ID, created.CAE, created.VoucherNumber})
	resultDigest := sha256.Sum256(resultBytes)
	result.ResultSHA256 = hex.EncodeToString(resultDigest[:])
	return result, nil
}
