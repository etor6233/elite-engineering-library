package postgres

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"elite.local/enterprise/internal/accounting"
	"elite.local/enterprise/internal/returnaccounting"
	"elite.local/enterprise/internal/returneffects"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReturnAccounting struct {
	pool    *pgxpool.Pool
	effects *ReturnEffects
	ledger  *Accounting
}

func NewReturnAccounting(pool *pgxpool.Pool) *ReturnAccounting {
	return &ReturnAccounting{pool: pool, effects: NewReturnEffects(pool), ledger: NewAccounting(pool)}
}
func (r *ReturnAccounting) Claim(ctx context.Context, owner, worker, token string, lease time.Duration) (*returneffects.Work, error) {
	return r.effects.Claim(ctx, owner, worker, token, lease)
}
func (r *ReturnAccounting) Finish(ctx context.Context, work returneffects.Work, worker string, c returneffects.Completion) error {
	return r.effects.Finish(ctx, work, worker, c)
}

func stableReturnAccountingUUID(scope, value string) string {
	s := sha256.Sum256([]byte(scope + "\x00" + value))
	b := s[:16]
	b[6] = (b[6] & 0x0f) | 0x50
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

type accountingSource struct {
	disposition, order, organization, remedy, journal, status, currency string
	version, total                                                      int64
}

func (r *ReturnAccounting) PostOrReconcile(ctx context.Context, work returneffects.Work, worker string) (returnaccounting.Result, error) {
	if work.EffectKind != "accounting" || work.OwnerContext != "accounting" {
		return returnaccounting.Result{}, returnaccounting.ErrConflict
	}
	var s accountingSource
	err := r.pool.QueryRow(ctx, `
select e.disposition_id,rr.order_id,rr.organization_id,d.customer_remedy,j.journal_id,j.status,j.version,j.currency,j.total_debit_minor_units
from sales.return_effect_execution x
join sales.return_effect_request e on e.tenant_id=x.tenant_id and e.request_id=x.request_id
join sales.return_disposition d on d.tenant_id=e.tenant_id and d.disposition_id=e.disposition_id
join sales.return_receipt rr on rr.tenant_id=d.tenant_id and rr.receipt_id=d.receipt_id
join accounting.journal j on j.tenant_id=rr.tenant_id and j.organization_id=rr.organization_id and j.source_type='SALE' and j.source_id=rr.order_id and j.reversal_of is null
join sales.return_effect_request inv on inv.tenant_id=e.tenant_id and inv.disposition_id=e.disposition_id and inv.effect_kind='inventory'
join sales.return_effect_execution invx on invx.tenant_id=inv.tenant_id and invx.request_id=inv.request_id and invx.status='succeeded'
join sales.return_effect_request remedy on remedy.tenant_id=e.tenant_id and remedy.disposition_id=e.disposition_id and remedy.effect_kind=d.customer_remedy
join sales.return_effect_execution remedyx on remedyx.tenant_id=remedy.tenant_id and remedyx.request_id=remedy.request_id and remedyx.status='succeeded'
where x.tenant_id=$1 and x.request_id=$2 and x.status='claimed' and x.claimed_by=$3 and x.claim_token=$4 and x.claimed_until>=clock_timestamp()
and e.effect_kind='accounting' and e.owner_context='accounting' and j.status in('posted','reversed')`, work.TenantID, work.RequestID, worker, work.ClaimToken).Scan(&s.disposition, &s.order, &s.organization, &s.remedy, &s.journal, &s.status, &s.version, &s.currency, &s.total)
	if errors.Is(err, pgx.ErrNoRows) {
		return returnaccounting.Result{}, returnaccounting.ErrConflict
	}
	if err != nil {
		return returnaccounting.Result{}, err
	}
	reversalID := stableReturnAccountingUUID("journal", work.TenantID+"/"+work.RequestID)
	registerID := stableReturnAccountingUUID("register", work.TenantID+"/"+work.RequestID)
	eventID := stableReturnAccountingUUID("event", work.TenantID+"/"+work.RequestID)
	if s.status == "posted" {
		_, err = r.ledger.ReverseJournal(ctx, work.TenantID, s.organization, s.journal, reversalID, registerID, eventID, s.version, worker, "approved customer return "+work.RequestID)
		if err != nil && !errors.Is(err, accounting.ErrConflict) {
			return returnaccounting.Result{}, err
		}
	}
	var reversalStatus string
	var reversedTotal int64
	var exactLines int
	err = r.pool.QueryRow(ctx, `select r.status,r.total_debit_minor_units,count(*)::int from accounting.journal r join accounting.entry re on re.tenant_id=r.tenant_id and re.journal_id=r.journal_id join accounting.entry oe on oe.tenant_id=re.tenant_id and oe.entry_id=re.original_entry_id where r.tenant_id=$1 and r.journal_id=$2 and r.reversal_of=$3 and r.source_type='REVERSAL' and r.source_id=$3 and re.debit_minor_units=oe.credit_minor_units and re.credit_minor_units=oe.debit_minor_units group by r.status,r.total_debit_minor_units`, work.TenantID, reversalID, s.journal).Scan(&reversalStatus, &reversedTotal, &exactLines)
	if errors.Is(err, pgx.ErrNoRows) {
		return returnaccounting.Result{}, returnaccounting.ErrConflict
	}
	if err != nil {
		return returnaccounting.Result{}, err
	}
	if reversalStatus != "posted" || reversedTotal != s.total || exactLines < 2 {
		return returnaccounting.Result{}, returnaccounting.ErrConflict
	}

	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return returnaccounting.Result{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var attempt int
	var started time.Time
	err = tx.QueryRow(ctx, `select attempt_count,claimed_at from sales.return_effect_execution where tenant_id=$1 and request_id=$2 and status='claimed' and claimed_by=$3 and claim_token=$4 and claimed_until>=clock_timestamp() for update`, work.TenantID, work.RequestID, worker, work.ClaimToken).Scan(&attempt, &started)
	if errors.Is(err, pgx.ErrNoRows) {
		return returnaccounting.Result{}, returnaccounting.ErrConflict
	}
	if err != nil {
		return returnaccounting.Result{}, err
	}
	if attempt != work.Attempt {
		return returnaccounting.Result{}, returnaccounting.ErrConflict
	}
	command, err := json.Marshal(struct {
		RequestID, OrderID, OriginalJournalID, ReversalJournalID, Remedy, Currency string
		ReversedMinorUnits                                                         int64
	}{work.RequestID, s.order, s.journal, reversalID, s.remedy, s.currency, reversedTotal})
	if err != nil {
		return returnaccounting.Result{}, err
	}
	digest := sha256.Sum256(command)
	hash := hex.EncodeToString(digest[:])
	if _, err = tx.Exec(ctx, `insert into accounting.return_effect_posting(tenant_id,request_id,disposition_id,source_order_id,remedy,original_journal_id,reversal_journal_id,currency,reversed_minor_units,status,result_sha256_hex)values($1,$2,$3,$4,$5,$6,$7,$8,$9,'posted',$10)`, work.TenantID, work.RequestID, s.disposition, s.order, s.remedy, s.journal, reversalID, s.currency, reversedTotal, hash); err != nil {
		return returnaccounting.Result{}, err
	}
	if _, err = tx.Exec(ctx, `insert into sales.return_effect_attempt(tenant_id,attempt_id,request_id,attempt_no,worker_id,outcome,provider_reference,result_sha256_hex,started_at)values($1,$2,$3,$4,$5,'succeeded',$6,$7,$8)`, work.TenantID, work.ClaimToken, work.RequestID, attempt, worker, reversalID, hash, started); err != nil {
		return returnaccounting.Result{}, err
	}
	updated, err := tx.Exec(ctx, `update sales.return_effect_execution set status='succeeded',claimed_at=null,claimed_by=null,claim_token=null,claimed_until=null,last_error_code=null,provider_reference=$5,result_sha256_hex=$6,updated_at=clock_timestamp() where tenant_id=$1 and request_id=$2 and status='claimed' and claimed_by=$3 and claim_token=$4`, work.TenantID, work.RequestID, worker, work.ClaimToken, reversalID, hash)
	if err != nil {
		return returnaccounting.Result{}, err
	}
	if updated.RowsAffected() != 1 {
		return returnaccounting.Result{}, returnaccounting.ErrConflict
	}
	if _, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)values($1,gen_random_uuid(),'return-accounting',$2,1,'return.accounting-posted',1,clock_timestamp(),$3)`, work.TenantID, work.RequestID, command); err != nil {
		return returnaccounting.Result{}, err
	}
	if err = tx.Commit(ctx); err != nil {
		return returnaccounting.Result{}, err
	}
	return returnaccounting.Result{RequestID: work.RequestID, OriginalJournalID: s.journal, ReversalJournalID: reversalID, ReversedMinorUnits: reversedTotal, ResultSHA256: hash}, nil
}
