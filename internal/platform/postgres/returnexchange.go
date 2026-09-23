package postgres

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"time"

	"elite.local/enterprise/internal/returneffects"
	"elite.local/enterprise/internal/returnexchange"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReturnExchange struct {
	pool    *pgxpool.Pool
	effects *ReturnEffects
}

func NewReturnExchange(pool *pgxpool.Pool) *ReturnExchange {
	return &ReturnExchange{pool: pool, effects: NewReturnEffects(pool)}
}

func (r *ReturnExchange) Claim(ctx context.Context, owner, worker, token string, lease time.Duration) (*returneffects.Work, error) {
	return r.effects.Claim(ctx, owner, worker, token, lease)
}

func (r *ReturnExchange) Finish(ctx context.Context, work returneffects.Work, worker string, completion returneffects.Completion) error {
	return r.effects.Finish(ctx, work, worker, completion)
}

type exchangeSource struct {
	attempt           int
	startedAt         time.Time
	dispositionID     string
	organizationID    string
	originalOrderID   string
	originalLineID    string
	originalStockID   string
	customerID        string
	variantID         string
	currency          string
	unitPrice         int64
	accountingRequest string
}

func (r *ReturnExchange) PrepareExchange(ctx context.Context, work returneffects.Work, worker string) (returnexchange.Result, error) {
	if work.EffectKind != "exchange" || work.OwnerContext != "fulfillment" {
		return returnexchange.Result{}, returnexchange.ErrConflict
	}
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return returnexchange.Result{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var source exchangeSource
	err = tx.QueryRow(ctx, `
select x.attempt_count,x.claimed_at,e.disposition_id,r.organization_id,r.order_id,l.line_id,r.stock_unit_id,r.customer_principal_id,l.variant_id,o.currency,l.unit_price_minor_units,accounting.request_id
  from sales.return_effect_execution x
  join sales.return_effect_request e on e.tenant_id=x.tenant_id and e.request_id=x.request_id
  join sales.return_disposition d on d.tenant_id=e.tenant_id and d.disposition_id=e.disposition_id
  join sales.return_receipt r on r.tenant_id=d.tenant_id and r.receipt_id=d.receipt_id
  join sales.return_authorization a on a.tenant_id=r.tenant_id and a.authorization_id=r.authorization_id
  join sales.customer_order o on o.tenant_id=r.tenant_id and o.order_id=r.order_id
  join sales.customer_order_line l on l.tenant_id=o.tenant_id and l.order_id=o.order_id and l.allocated_stock_unit_id=r.stock_unit_id
  join sales.return_effect_request inventory on inventory.tenant_id=e.tenant_id and inventory.disposition_id=e.disposition_id and inventory.effect_kind='inventory' and inventory.owner_context='inventory'
  join sales.return_effect_execution inventory_x on inventory_x.tenant_id=inventory.tenant_id and inventory_x.request_id=inventory.request_id and inventory_x.status='succeeded'
  join sales.return_effect_request accounting on accounting.tenant_id=e.tenant_id and accounting.disposition_id=e.disposition_id and accounting.effect_kind='accounting' and accounting.owner_context='accounting'
 where x.tenant_id=$1 and x.request_id=$2 and x.status='claimed' and x.claimed_by=$3 and x.claim_token=$4 and x.claimed_until>=clock_timestamp()
   and e.effect_kind='exchange' and e.owner_context='fulfillment' and d.customer_remedy='exchange' and a.disposition='exchange'
   and o.state='delivered' and l.quantity=1 and l.unit_price_minor_units>0 and o.total_minor_units=l.unit_price_minor_units
   and (select count(*) from sales.customer_order_line all_lines where all_lines.tenant_id=o.tenant_id and all_lines.order_id=o.order_id)=1
 for update of x,o,l`, work.TenantID, work.RequestID, worker, work.ClaimToken).Scan(
		&source.attempt, &source.startedAt, &source.dispositionID, &source.organizationID, &source.originalOrderID, &source.originalLineID, &source.originalStockID, &source.customerID, &source.variantID, &source.currency, &source.unitPrice, &source.accountingRequest,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return returnexchange.Result{}, returnexchange.ErrConflict
	}
	if err != nil {
		return returnexchange.Result{}, err
	}
	if source.attempt != work.Attempt {
		return returnexchange.Result{}, returnexchange.ErrConflict
	}

	var replacementStock string
	var stockVersion int64
	err = tx.QueryRow(ctx, `select stock_unit_id,version from inventory.stock_unit where tenant_id=$1 and organization_id=$2 and variant_id=$3 and state='available' and stock_unit_id<>$4 order by received_at nulls last,stock_unit_id for update skip locked limit 1`, work.TenantID, source.organizationID, source.variantID, source.originalStockID).Scan(&replacementStock, &stockVersion)
	if errors.Is(err, pgx.ErrNoRows) {
		return returnexchange.Result{}, returnexchange.ErrStockUnavailable
	}
	if err != nil {
		return returnexchange.Result{}, err
	}

	var replacementOrder, replacementLine, replacementHandover string
	if err = tx.QueryRow(ctx, `select gen_random_uuid()::text,gen_random_uuid()::text,gen_random_uuid()::text`).Scan(&replacementOrder, &replacementLine, &replacementHandover); err != nil {
		return returnexchange.Result{}, err
	}
	if _, err = tx.Exec(ctx, `insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version) values($1,$2,$3,$4,'allocated',$5,0,1)`, work.TenantID, replacementOrder, source.organizationID, source.customerID, source.currency); err != nil {
		return returnexchange.Result{}, err
	}
	if _, err = tx.Exec(ctx, `insert into sales.customer_order_line(tenant_id,order_id,line_id,variant_id,quantity,unit_price_minor_units,allocated_stock_unit_id) values($1,$2,$3,$4,1,0,$5)`, work.TenantID, replacementOrder, replacementLine, source.variantID, replacementStock); err != nil {
		return returnexchange.Result{}, err
	}
	updated, err := tx.Exec(ctx, `update inventory.stock_unit set state='reserved',version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and organization_id=$2 and stock_unit_id=$3 and variant_id=$4 and state='available' and version=$5`, work.TenantID, source.organizationID, replacementStock, source.variantID, stockVersion)
	if err != nil {
		return returnexchange.Result{}, err
	}
	if updated.RowsAffected() != 1 {
		return returnexchange.Result{}, returnexchange.ErrStockUnavailable
	}
	if _, err = tx.Exec(ctx, `insert into sales.delivery_handover(tenant_id,handover_id,organization_id,order_id,customer_principal_id,stock_unit_id,state,version) values($1,$2,$3,$4,$5,$6,'prepared',1)`, work.TenantID, replacementHandover, source.organizationID, replacementOrder, source.customerID, replacementStock); err != nil {
		return returnexchange.Result{}, err
	}

	command, err := json.Marshal(struct {
		RequestID             string `json:"request_id"`
		OriginalOrderID       string `json:"original_order_id"`
		ReplacementOrderID    string `json:"replacement_order_id"`
		ReplacementStockID    string `json:"replacement_stock_unit_id"`
		ReplacementHandoverID string `json:"replacement_handover_id"`
	}{work.RequestID, source.originalOrderID, replacementOrder, replacementStock, replacementHandover})
	if err != nil {
		return returnexchange.Result{}, err
	}
	digest := sha256.Sum256(command)
	resultHash := hex.EncodeToString(digest[:])

	if _, err = tx.Exec(ctx, `insert into sales.return_exchange(tenant_id,request_id,disposition_id,original_order_id,original_line_id,original_stock_unit_id,replacement_order_id,replacement_line_id,replacement_stock_unit_id,replacement_handover_id,accounting_request_id,currency,original_unit_price_minor_units,settlement_mode,status,result_sha256_hex) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,'even-exchange-zero-balance','prepared',$14)`, work.TenantID, work.RequestID, source.dispositionID, source.originalOrderID, source.originalLineID, source.originalStockID, replacementOrder, replacementLine, replacementStock, replacementHandover, source.accountingRequest, source.currency, source.unitPrice, resultHash); err != nil {
		return returnexchange.Result{}, err
	}
	if _, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload) values($1,gen_random_uuid(),'return-exchange',$2,1,'return.exchange-prepared',1,clock_timestamp(),$3),($1,gen_random_uuid(),'order',$4,1,'order.exchange-replacement-created',1,clock_timestamp(),$3),($1,gen_random_uuid(),'stock-unit',$5,$6,'stock-unit.exchange-reserved',1,clock_timestamp(),$3)`, work.TenantID, work.RequestID, command, replacementOrder, replacementStock, stockVersion+1); err != nil {
		return returnexchange.Result{}, err
	}
	if _, err = tx.Exec(ctx, `insert into sales.return_effect_attempt(tenant_id,attempt_id,request_id,attempt_no,worker_id,outcome,provider_reference,result_sha256_hex,started_at) values($1,$2,$3,$4,$5,'succeeded',$6,$7,$8)`, work.TenantID, work.ClaimToken, work.RequestID, source.attempt, worker, replacementOrder, resultHash, source.startedAt); err != nil {
		return returnexchange.Result{}, err
	}
	completed, err := tx.Exec(ctx, `update sales.return_effect_execution set status='succeeded',claimed_at=null,claimed_by=null,claim_token=null,claimed_until=null,last_error_code=null,provider_reference=$5,result_sha256_hex=$6,updated_at=clock_timestamp() where tenant_id=$1 and request_id=$2 and status='claimed' and claimed_by=$3 and claim_token=$4`, work.TenantID, work.RequestID, worker, work.ClaimToken, replacementOrder, resultHash)
	if err != nil {
		return returnexchange.Result{}, err
	}
	if completed.RowsAffected() != 1 {
		return returnexchange.Result{}, returnexchange.ErrConflict
	}
	if err = tx.Commit(ctx); err != nil {
		return returnexchange.Result{}, err
	}
	return returnexchange.Result{RequestID: work.RequestID, ReplacementOrderID: replacementOrder, ReplacementStockID: replacementStock, ReplacementHandoverID: replacementHandover, ResultSHA256: resultHash}, nil
}
