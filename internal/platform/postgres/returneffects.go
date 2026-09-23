package postgres

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"elite.local/enterprise/internal/returneffects"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReturnEffects struct{ pool *pgxpool.Pool }

func NewReturnEffects(pool *pgxpool.Pool) *ReturnEffects { return &ReturnEffects{pool: pool} }

func (r *ReturnEffects) Claim(ctx context.Context, owner, worker, claimToken string, lease time.Duration) (*returneffects.Work, error) {
	if owner == "" || worker == "" || claimToken == "" || lease <= 0 {
		return nil, fmt.Errorf("invalid return effect claim")
	}
	var work returneffects.Work
	err := r.pool.QueryRow(ctx, `
with candidate as (
  select x.tenant_id,x.request_id
    from sales.return_effect_execution x
    join sales.return_effect_request e on e.tenant_id=x.tenant_id and e.request_id=x.request_id
   where e.owner_context=$1
     -- Dependency readiness belongs to claim selection: waiting must not consume
     -- a lease/attempt or turn a valid later phase into a permanent source conflict.
     and (e.owner_context not in ('fulfillment','accounting','fiscal') or (
       exists(select 1 from sales.return_effect_request dep
         join sales.return_effect_execution dx on dx.tenant_id=dep.tenant_id and dx.request_id=dep.request_id
         where dep.tenant_id=e.tenant_id and dep.disposition_id=e.disposition_id
           and dep.effect_kind='inventory' and dep.owner_context='inventory' and dx.status='succeeded')
       and (e.owner_context<>'fulfillment' or exists(
         select 1 from sales.return_effect_request dep where dep.tenant_id=e.tenant_id
           and dep.disposition_id=e.disposition_id and dep.effect_kind='accounting' and dep.owner_context='accounting'))
       and (e.owner_context not in ('accounting','fiscal') or exists(
         select 1 from sales.return_disposition d
         join sales.return_effect_request dep on dep.tenant_id=d.tenant_id and dep.disposition_id=d.disposition_id
         join sales.return_effect_execution dx on dx.tenant_id=dep.tenant_id and dx.request_id=dep.request_id
         where d.tenant_id=e.tenant_id and d.disposition_id=e.disposition_id
           and dep.effect_kind=d.customer_remedy
           and ((d.customer_remedy='refund' and dep.owner_context='payment')
             or (d.customer_remedy='exchange' and dep.owner_context='fulfillment'))
           and dx.status='succeeded'))
       and (e.owner_context<>'fiscal' or exists(
         select 1 from sales.return_effect_request dep
         join sales.return_effect_execution dx on dx.tenant_id=dep.tenant_id and dx.request_id=dep.request_id
         where dep.tenant_id=e.tenant_id and dep.disposition_id=e.disposition_id
           and dep.effect_kind='accounting' and dep.owner_context='accounting' and dx.status='succeeded'))
     ))
     and ((x.status in ('requested','retry') and x.available_at<=clock_timestamp()) or (x.status='claimed' and x.claimed_until<clock_timestamp()))
   order by x.available_at,e.requested_at,x.request_id
   for update of x skip locked
   limit 1
), claimed as (
  update sales.return_effect_execution x
     set status='claimed',attempt_count=x.attempt_count+1,claimed_at=clock_timestamp(),claimed_by=$2,claim_token=$3,claimed_until=clock_timestamp()+$4::interval,last_error_code=null,provider_reference=null,result_sha256_hex=null,updated_at=clock_timestamp()
    from candidate c
   where x.tenant_id=c.tenant_id and x.request_id=c.request_id
  returning x.tenant_id,x.request_id,x.attempt_count,x.claim_token
)
select c.tenant_id,c.request_id,e.disposition_id,e.effect_kind,e.owner_context,e.idempotency_key,c.attempt_count,c.claim_token,d.inventory_action,r.organization_id,r.stock_unit_id
  from claimed c
  join sales.return_effect_request e on e.tenant_id=c.tenant_id and e.request_id=c.request_id
  join sales.return_disposition d on d.tenant_id=e.tenant_id and d.disposition_id=e.disposition_id
  join sales.return_receipt r on r.tenant_id=d.tenant_id and r.receipt_id=d.receipt_id`, owner, worker, claimToken, lease.String()).Scan(
		&work.TenantID, &work.RequestID, &work.DispositionID, &work.EffectKind, &work.OwnerContext, &work.IdempotencyKey, &work.Attempt, &work.ClaimToken, &work.InventoryAction, &work.OrganizationID, &work.StockUnitID,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &work, nil
}

func inventoryTarget(action string) (string, bool) {
	targets := map[string]string{"quarantine": "quarantine", "restock": "available", "repair": "service", "scrap": "retired"}
	target, ok := targets[action]
	return target, ok
}

func (r *ReturnEffects) ApplyInventory(ctx context.Context, work returneffects.Work, worker string) (returneffects.Result, error) {
	target, ok := inventoryTarget(work.InventoryAction)
	if !ok || work.OwnerContext != "inventory" || work.EffectKind != "inventory" {
		return returneffects.Result{}, returneffects.ErrInventoryState
	}
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return returneffects.Result{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var attempt int
	var startedAt time.Time
	var stockState string
	var stockVersion int64
	err = tx.QueryRow(ctx, `
select x.attempt_count,x.claimed_at,s.state,s.version
  from sales.return_effect_execution x
  join sales.return_effect_request e on e.tenant_id=x.tenant_id and e.request_id=x.request_id
  join sales.return_disposition d on d.tenant_id=e.tenant_id and d.disposition_id=e.disposition_id
  join sales.return_receipt r on r.tenant_id=d.tenant_id and r.receipt_id=d.receipt_id
  join inventory.stock_unit s on s.tenant_id=r.tenant_id and s.stock_unit_id=r.stock_unit_id and s.organization_id=r.organization_id
 where x.tenant_id=$1 and x.request_id=$2 and x.status='claimed' and x.claimed_by=$3 and x.claim_token=$4 and x.claimed_until>=clock_timestamp()
   and e.effect_kind='inventory' and e.owner_context='inventory' and d.inventory_action=$5 and r.organization_id=$6 and r.stock_unit_id=$7
 for update of x,s`, work.TenantID, work.RequestID, worker, work.ClaimToken, work.InventoryAction, work.OrganizationID, work.StockUnitID).Scan(&attempt, &startedAt, &stockState, &stockVersion)
	if errors.Is(err, pgx.ErrNoRows) {
		return returneffects.Result{}, returneffects.ErrInventoryState
	}
	if err != nil {
		return returneffects.Result{}, err
	}
	if attempt != work.Attempt || stockState != "sold" {
		return returneffects.Result{}, returneffects.ErrInventoryState
	}

	newVersion := stockVersion + 1
	command, err := json.Marshal(map[string]any{"request_id": work.RequestID, "stock_unit_id": work.StockUnitID, "from_state": stockState, "to_state": target, "version": newVersion})
	if err != nil {
		return returneffects.Result{}, err
	}
	digest := sha256.Sum256(command)
	resultHash := hex.EncodeToString(digest[:])

	updated, err := tx.Exec(ctx, `update inventory.stock_unit set state=$5,version=$6,received_at=coalesce(received_at,clock_timestamp()),updated_at=clock_timestamp() where tenant_id=$1 and organization_id=$2 and stock_unit_id=$3 and state=$4 and version=$7`, work.TenantID, work.OrganizationID, work.StockUnitID, stockState, target, newVersion, stockVersion)
	if err != nil {
		return returneffects.Result{}, err
	}
	if updated.RowsAffected() != 1 {
		return returneffects.Result{}, returneffects.ErrInventoryState
	}

	if _, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload) values($1,gen_random_uuid(),'stock-unit',$2,$3,$4,1,clock_timestamp(),$5)`, work.TenantID, work.StockUnitID, newVersion, "stock-unit.return-"+target, command); err != nil {
		return returneffects.Result{}, err
	}
	if _, err = tx.Exec(ctx, `insert into sales.return_effect_attempt(tenant_id,attempt_id,request_id,attempt_no,worker_id,outcome,provider_reference,result_sha256_hex,started_at) values($1,$2,$3,$4,$5,'succeeded',$6,$7,$8)`, work.TenantID, work.ClaimToken, work.RequestID, attempt, worker, work.StockUnitID, resultHash, startedAt); err != nil {
		return returneffects.Result{}, err
	}
	completed, err := tx.Exec(ctx, `update sales.return_effect_execution set status='succeeded',claimed_at=null,claimed_by=null,claim_token=null,claimed_until=null,last_error_code=null,provider_reference=$5,result_sha256_hex=$6,updated_at=clock_timestamp() where tenant_id=$1 and request_id=$2 and status='claimed' and claimed_by=$3 and claim_token=$4`, work.TenantID, work.RequestID, worker, work.ClaimToken, work.StockUnitID, resultHash)
	if err != nil {
		return returneffects.Result{}, err
	}
	if completed.RowsAffected() != 1 {
		return returneffects.Result{}, returneffects.ErrInventoryState
	}
	if err = tx.Commit(ctx); err != nil {
		return returneffects.Result{}, err
	}
	return returneffects.Result{RequestID: work.RequestID, StockUnitID: work.StockUnitID, InventoryState: target, StockVersion: newVersion, ResultSHA256: resultHash}, nil
}

func (r *ReturnEffects) Finish(ctx context.Context, work returneffects.Work, worker string, completion returneffects.Completion) error {
	if completion.Outcome != "retry" && completion.Outcome != "blocked" && completion.Outcome != "failed" && completion.Outcome != "succeeded" {
		return fmt.Errorf("invalid completion outcome")
	}
	if completion.Outcome == "succeeded" {
		if completion.ErrorCode != "" || len(completion.ResultSHA256) != 64 {
			return fmt.Errorf("invalid successful completion")
		}
	} else if completion.ErrorCode == "" || completion.ResultSHA256 != "" {
		return fmt.Errorf("invalid unsuccessful completion")
	}
	if completion.Outcome == "retry" && completion.RetryAfter <= 0 {
		return fmt.Errorf("retry delay required")
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var attempt int
	var startedAt time.Time
	err = tx.QueryRow(ctx, `select attempt_count,claimed_at from sales.return_effect_execution where tenant_id=$1 and request_id=$2 and status='claimed' and claimed_by=$3 and claim_token=$4 and claimed_until>=clock_timestamp() for update`, work.TenantID, work.RequestID, worker, work.ClaimToken).Scan(&attempt, &startedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return returneffects.ErrInventoryState
	}
	if err != nil {
		return err
	}
	if attempt != work.Attempt {
		return returneffects.ErrInventoryState
	}
	if _, err = tx.Exec(ctx, `insert into sales.return_effect_attempt(tenant_id,attempt_id,request_id,attempt_no,worker_id,outcome,error_code,provider_reference,result_sha256_hex,started_at) values($1,$2,$3,$4,$5,$6,nullif($7,''),nullif($8,''),nullif($9,''),$10)`, work.TenantID, work.ClaimToken, work.RequestID, attempt, worker, completion.Outcome, completion.ErrorCode, completion.ProviderReference, completion.ResultSHA256, startedAt); err != nil {
		return err
	}
	available := time.Now().UTC()
	if completion.Outcome == "retry" {
		available = available.Add(completion.RetryAfter)
	}
	updated, err := tx.Exec(ctx, `update sales.return_effect_execution set status=$5,available_at=$6,claimed_at=null,claimed_by=null,claim_token=null,claimed_until=null,last_error_code=nullif($7,''),provider_reference=nullif($8,''),result_sha256_hex=nullif($9,''),updated_at=clock_timestamp() where tenant_id=$1 and request_id=$2 and status='claimed' and claimed_by=$3 and claim_token=$4`, work.TenantID, work.RequestID, worker, work.ClaimToken, completion.Outcome, available, completion.ErrorCode, completion.ProviderReference, completion.ResultSHA256)
	if err != nil {
		return err
	}
	if updated.RowsAffected() != 1 {
		return returneffects.ErrInventoryState
	}
	return tx.Commit(ctx)
}

func (r *ReturnEffects) ResumeBlocked(ctx context.Context, tenant, requestID, resumeID, actor, reason string) error {
	if tenant == "" || requestID == "" || resumeID == "" || actor == "" || reason == "" {
		return fmt.Errorf("invalid return effect resume")
	}
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	if _, err = tx.Exec(ctx, `insert into sales.return_effect_resume(tenant_id,resume_id,request_id,reason_code,requested_by_subject) values($1,$2,$3,$4,$5)`, tenant, resumeID, requestID, reason, actor); err != nil {
		return err
	}
	updated, err := tx.Exec(ctx, `update sales.return_effect_execution set status='requested',available_at=clock_timestamp(),last_error_code=null,provider_reference=null,result_sha256_hex=null,updated_at=clock_timestamp() where tenant_id=$1 and request_id=$2 and status='blocked'`, tenant, requestID)
	if err != nil {
		return err
	}
	if updated.RowsAffected() != 1 {
		return returneffects.ErrInventoryState
	}
	return tx.Commit(ctx)
}
