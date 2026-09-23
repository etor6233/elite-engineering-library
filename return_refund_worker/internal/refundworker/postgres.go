package refundworker

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStore struct{ pool *pgxpool.Pool }

func NewPostgresStore(pool *pgxpool.Pool) *PostgresStore { return &PostgresStore{pool: pool} }

func (s *PostgresStore) Claim(ctx context.Context, workerID, claimToken string, lease time.Duration) (*Work, error) {
	if s == nil || s.pool == nil || workerID == "" || claimToken == "" || lease <= 0 {
		return nil, fmt.Errorf("invalid refund claim")
	}
	var work Work
	err := s.pool.QueryRow(ctx, `
with candidate as (
  select x.tenant_id,x.request_id
    from sales.return_effect_execution x
    join sales.return_effect_request e on e.tenant_id=x.tenant_id and e.request_id=x.request_id
   where e.effect_kind='refund' and e.owner_context='payment'
     and ((x.status in ('requested','retry') and x.available_at<=clock_timestamp()) or (x.status='claimed' and x.claimed_until<clock_timestamp()))
   order by x.available_at,e.requested_at,x.request_id
   for update of x skip locked limit 1
), claimed as (
  update sales.return_effect_execution x
     set status='claimed',attempt_count=x.attempt_count+1,claimed_at=clock_timestamp(),claimed_by=$1,claim_token=$2,claimed_until=clock_timestamp()+$3::interval,last_error_code=null,result_sha256_hex=null,updated_at=clock_timestamp()
    from candidate c where x.tenant_id=c.tenant_id and x.request_id=c.request_id
  returning x.tenant_id,x.request_id,x.attempt_count,x.claim_token
)
select c.tenant_id,c.request_id,e.idempotency_key,c.attempt_count,c.claim_token
  from claimed c join sales.return_effect_request e on e.tenant_id=c.tenant_id and e.request_id=c.request_id`, workerID, claimToken, lease.String()).Scan(&work.TenantID, &work.RequestID, &work.IdempotencyKey, &work.Attempt, &work.ClaimToken)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &work, nil
}

func scanRefund(row pgx.Row) (Refund, error) {
	var value Refund
	err := row.Scan(&value.TenantID, &value.RequestID, &value.PaymentAttemptID, &value.OrderID, &value.LineID, &value.StockUnitID, &value.Provider, &value.ProviderPaymentReference, &value.IdempotencyKey, &value.Currency, &value.AmountMinorUnits, &value.ProviderRefundReference, &value.ProviderStatus)
	return value, err
}

func (s *PostgresStore) Prepare(ctx context.Context, work Work, workerID string) (Refund, error) {
	if s == nil || s.pool == nil {
		return Refund{}, fmt.Errorf("invalid refund store")
	}
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return Refund{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var attempt int
	err = tx.QueryRow(ctx, `select attempt_count from sales.return_effect_execution where tenant_id=$1 and request_id=$2 and status='claimed' and claimed_by=$3 and claim_token=$4 and claimed_until>=clock_timestamp() for update`, work.TenantID, work.RequestID, workerID, work.ClaimToken).Scan(&attempt)
	if errors.Is(err, pgx.ErrNoRows) {
		return Refund{}, ErrStaleClaim
	}
	if err != nil {
		return Refund{}, err
	}
	if attempt != work.Attempt {
		return Refund{}, ErrStaleClaim
	}

	existing, existingErr := scanRefund(tx.QueryRow(ctx, `select tenant_id,request_id,payment_attempt_id,order_id,line_id,stock_unit_id,provider_code,provider_payment_reference,idempotency_key,currency,amount_minor_units,coalesce(provider_refund_reference,''),coalesce(provider_status,'') from payment.return_refund where tenant_id=$1 and request_id=$2 for update`, work.TenantID, work.RequestID))
	if existingErr == nil {
		if err = tx.Commit(ctx); err != nil {
			return Refund{}, err
		}
		return existing, nil
	}
	if !errors.Is(existingErr, pgx.ErrNoRows) {
		return Refund{}, existingErr
	}

	type lineSource struct {
		orderID, lineID, stockID, currency string
		amount                             int64
	}
	rows, err := tx.Query(ctx, `
select r.order_id,l.line_id,r.stock_unit_id,o.currency,l.unit_price_minor_units
  from sales.return_effect_request e
  join sales.return_disposition d on d.tenant_id=e.tenant_id and d.disposition_id=e.disposition_id
  join sales.return_receipt r on r.tenant_id=d.tenant_id and r.receipt_id=d.receipt_id
  join sales.customer_order o on o.tenant_id=r.tenant_id and o.order_id=r.order_id and o.organization_id=r.organization_id
  join sales.customer_order_line l on l.tenant_id=r.tenant_id and l.order_id=r.order_id and l.allocated_stock_unit_id=r.stock_unit_id and l.quantity=1
 where e.tenant_id=$1 and e.request_id=$2 and e.effect_kind='refund' and e.owner_context='payment' and e.state='requested' and d.customer_remedy='refund'
 for update of l,o`, work.TenantID, work.RequestID)
	if err != nil {
		return Refund{}, err
	}
	var lines []lineSource
	for rows.Next() {
		var line lineSource
		if err = rows.Scan(&line.orderID, &line.lineID, &line.stockID, &line.currency, &line.amount); err != nil {
			rows.Close()
			return Refund{}, err
		}
		lines = append(lines, line)
	}
	rows.Close()
	if err = rows.Err(); err != nil {
		return Refund{}, err
	}
	if len(lines) != 1 || lines[0].amount <= 0 {
		return Refund{}, ErrMappingConflict
	}

	type paymentSource struct {
		id, provider, reference, currency string
		amount, refunded                  int64
	}
	rows, err = tx.Query(ctx, `
-- AUTHORED provider-name compatibility at the Commerce -> refund-owner boundary.
-- Preserve payment/refund IDs and idempotency keys; existing refund rows retain
-- their historical mercado_pago provider code and are returned above unchanged.
select p.payment_attempt_id,case p.provider_code when 'mercadopago' then 'mercado_pago' else p.provider_code end,p.provider_reference,p.currency,p.amount_minor_units,
       coalesce((select sum(rr.amount_minor_units) from payment.return_refund rr where rr.tenant_id=p.tenant_id and rr.payment_attempt_id=p.payment_attempt_id and rr.state='succeeded'),0)
  from payment.payment_attempt p
 where p.tenant_id=$1 and p.order_id=$2 and p.state in ('captured','disputed') and p.provider_reference is not null and p.provider_code in ('stripe','mercado_pago','mercadopago')
 for update of p`, work.TenantID, lines[0].orderID)
	if err != nil {
		return Refund{}, err
	}
	var payments []paymentSource
	for rows.Next() {
		var value paymentSource
		if err = rows.Scan(&value.id, &value.provider, &value.reference, &value.currency, &value.amount, &value.refunded); err != nil {
			rows.Close()
			return Refund{}, err
		}
		if value.currency == lines[0].currency && value.amount-value.refunded >= lines[0].amount {
			payments = append(payments, value)
		}
	}
	rows.Close()
	if err = rows.Err(); err != nil {
		return Refund{}, err
	}
	if len(payments) != 1 {
		return Refund{}, ErrMappingConflict
	}

	value := Refund{TenantID: work.TenantID, RequestID: work.RequestID, PaymentAttemptID: payments[0].id, OrderID: lines[0].orderID, LineID: lines[0].lineID, StockUnitID: lines[0].stockID, Provider: payments[0].provider, ProviderPaymentReference: payments[0].reference, IdempotencyKey: work.IdempotencyKey, Currency: lines[0].currency, AmountMinorUnits: lines[0].amount}
	_, err = tx.Exec(ctx, `insert into payment.return_refund(tenant_id,request_id,payment_attempt_id,order_id,line_id,stock_unit_id,provider_code,provider_payment_reference,idempotency_key,currency,amount_minor_units,state) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,'prepared')`, value.TenantID, value.RequestID, value.PaymentAttemptID, value.OrderID, value.LineID, value.StockUnitID, value.Provider, value.ProviderPaymentReference, value.IdempotencyKey, value.Currency, value.AmountMinorUnits)
	if err != nil {
		return Refund{}, err
	}
	if err = tx.Commit(ctx); err != nil {
		return Refund{}, err
	}
	return value, nil
}

func canonicalResult(result ProviderResult) (string, error) {
	encoded, err := json.Marshal(result)
	if err != nil {
		return "", err
	}
	digest := sha256.Sum256(encoded)
	return hex.EncodeToString(digest[:]), nil
}

func refundState(outcome, provider, status string) string {
	if outcome == "succeeded" {
		return "succeeded"
	}
	if outcome == "blocked" {
		return "blocked"
	}
	if outcome == "failed" {
		if (provider == "stripe" && status == "canceled") || (provider == "mercado_pago" && (status == "cancelled" || status == "canceled")) {
			return "canceled"
		}
		return "failed"
	}
	if status == "requires_action" {
		return "requires_action"
	}
	return "pending"
}

func validOutcome(outcome, code string, retry time.Duration) bool {
	if outcome == "succeeded" {
		return code == "" && retry == 0
	}
	if outcome == "retry" {
		return code != "" && retry > 0
	}
	return (outcome == "blocked" || outcome == "failed") && code != "" && retry == 0
}

func (s *PostgresStore) Complete(ctx context.Context, work Work, workerID string, refund Refund, result ProviderResult, outcome, code string, retryAfter time.Duration) error {
	if s == nil || s.pool == nil || !validOutcome(outcome, code, retryAfter) {
		return fmt.Errorf("invalid refund completion")
	}
	resultHash, err := canonicalResult(result)
	if err != nil {
		return err
	}
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var attempt int
	var startedAt time.Time
	err = tx.QueryRow(ctx, `select attempt_count,claimed_at from sales.return_effect_execution where tenant_id=$1 and request_id=$2 and status='claimed' and claimed_by=$3 and claim_token=$4 and claimed_until>=clock_timestamp() for update`, work.TenantID, work.RequestID, workerID, work.ClaimToken).Scan(&attempt, &startedAt)
	if errors.Is(err, pgx.ErrNoRows) || attempt != work.Attempt {
		return ErrStaleClaim
	}
	if err != nil {
		return err
	}
	var stored Refund
	stored, err = scanRefund(tx.QueryRow(ctx, `select tenant_id,request_id,payment_attempt_id,order_id,line_id,stock_unit_id,provider_code,provider_payment_reference,idempotency_key,currency,amount_minor_units,coalesce(provider_refund_reference,''),coalesce(provider_status,'') from payment.return_refund where tenant_id=$1 and request_id=$2 for update`, work.TenantID, work.RequestID))
	if err != nil || stored.PaymentAttemptID != refund.PaymentAttemptID || stored.AmountMinorUnits != refund.AmountMinorUnits || stored.Currency != refund.Currency {
		if err != nil {
			return err
		}
		return ErrMappingConflict
	}
	state := refundState(outcome, refund.Provider, result.ProviderStatus)
	updated, err := tx.Exec(ctx, `update payment.return_refund set provider_refund_reference=$3,provider_status=$4,state=$5,response_sha256_hex=$6,version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and request_id=$2`, work.TenantID, work.RequestID, result.ProviderRefundReference, result.ProviderStatus, state, resultHash)
	if err != nil || updated.RowsAffected() != 1 {
		if err != nil {
			return err
		}
		return ErrMappingConflict
	}
	_, err = tx.Exec(ctx, `insert into payment.return_refund_observation(tenant_id,observation_id,request_id,provider_refund_reference,provider_status,amount_minor_units,currency,response_sha256_hex) values($1,$2,$3,$4,$5,$6,$7,$8)`, work.TenantID, work.ClaimToken, work.RequestID, result.ProviderRefundReference, result.ProviderStatus, result.AmountMinorUnits, result.Currency, resultHash)
	if err != nil {
		return err
	}
	if outcome == "succeeded" {
		var paymentAmount, paymentVersion int64
		var paymentState string
		err = tx.QueryRow(ctx, `select amount_minor_units,state,version from payment.payment_attempt where tenant_id=$1 and payment_attempt_id=$2 for update`, work.TenantID, refund.PaymentAttemptID).Scan(&paymentAmount, &paymentState, &paymentVersion)
		if err != nil {
			return err
		}
		var refunded int64
		if err = tx.QueryRow(ctx, `select coalesce(sum(amount_minor_units),0) from payment.return_refund where tenant_id=$1 and payment_attempt_id=$2 and state='succeeded'`, work.TenantID, refund.PaymentAttemptID).Scan(&refunded); err != nil {
			return err
		}
		if refunded > paymentAmount {
			return ErrMappingConflict
		}
		if refunded == paymentAmount && (paymentState == "captured" || paymentState == "disputed") {
			if _, err = tx.Exec(ctx, `update payment.payment_attempt set state='refunded',version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and payment_attempt_id=$2 and version=$3`, work.TenantID, refund.PaymentAttemptID, paymentVersion); err != nil {
				return err
			}
			if _, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload) values($1,gen_random_uuid(),'payment',$2,$3,'payment.refunded',1,clock_timestamp(),jsonb_build_object('return_request_id',$4::text,'refund_reference',$5::text,'amount_minor_units',$6::bigint,'currency',$7::text))`, work.TenantID, refund.PaymentAttemptID, paymentVersion+1, work.RequestID, result.ProviderRefundReference, refund.AmountMinorUnits, refund.Currency); err != nil {
				return err
			}
		}
		if _, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload) values($1,gen_random_uuid(),'return-refund',$2,1,'return-refund.succeeded',1,clock_timestamp(),jsonb_build_object('payment_attempt_id',$3::text,'line_id',$4::text,'stock_unit_id',$5::text,'provider',$6::text,'refund_reference',$7::text,'amount_minor_units',$8::bigint,'currency',$9::text))`, work.TenantID, work.RequestID, refund.PaymentAttemptID, refund.LineID, refund.StockUnitID, refund.Provider, result.ProviderRefundReference, refund.AmountMinorUnits, refund.Currency); err != nil {
			return err
		}
	}
	if _, err = tx.Exec(ctx, `insert into sales.return_effect_attempt(tenant_id,attempt_id,request_id,attempt_no,worker_id,outcome,error_code,provider_reference,result_sha256_hex,started_at) values($1,$2,$3,$4,$5,$6,nullif($7,''),$8,case when $6='succeeded' then $9 else null end,$10)`, work.TenantID, work.ClaimToken, work.RequestID, attempt, workerID, outcome, code, result.ProviderRefundReference, resultHash, startedAt); err != nil {
		return err
	}
	available := time.Now().UTC()
	if outcome == "retry" {
		available = available.Add(retryAfter)
	}
	updated, err = tx.Exec(ctx, `update sales.return_effect_execution set status=$5,available_at=$6,claimed_at=null,claimed_by=null,claim_token=null,claimed_until=null,last_error_code=nullif($7,''),provider_reference=$8,result_sha256_hex=case when $5='succeeded' then $9 else null end,updated_at=clock_timestamp() where tenant_id=$1 and request_id=$2 and status='claimed' and claimed_by=$3 and claim_token=$4`, work.TenantID, work.RequestID, workerID, work.ClaimToken, outcome, available, code, result.ProviderRefundReference, resultHash)
	if err != nil || updated.RowsAffected() != 1 {
		if err != nil {
			return err
		}
		return ErrStaleClaim
	}
	return tx.Commit(ctx)
}

func (s *PostgresStore) Finish(ctx context.Context, work Work, workerID, outcome, code string, retryAfter time.Duration) error {
	if s == nil || s.pool == nil || !validOutcome(outcome, code, retryAfter) || outcome == "succeeded" {
		return fmt.Errorf("invalid refund finish")
	}
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var attempt int
	var startedAt time.Time
	err = tx.QueryRow(ctx, `select attempt_count,claimed_at from sales.return_effect_execution where tenant_id=$1 and request_id=$2 and status='claimed' and claimed_by=$3 and claim_token=$4 and claimed_until>=clock_timestamp() for update`, work.TenantID, work.RequestID, workerID, work.ClaimToken).Scan(&attempt, &startedAt)
	if errors.Is(err, pgx.ErrNoRows) || attempt != work.Attempt {
		return ErrStaleClaim
	}
	if err != nil {
		return err
	}
	if _, err = tx.Exec(ctx, `insert into sales.return_effect_attempt(tenant_id,attempt_id,request_id,attempt_no,worker_id,outcome,error_code,started_at) values($1,$2,$3,$4,$5,$6,$7,$8)`, work.TenantID, work.ClaimToken, work.RequestID, attempt, workerID, outcome, code, startedAt); err != nil {
		return err
	}
	available := time.Now().UTC()
	if outcome == "retry" {
		available = available.Add(retryAfter)
	}
	updated, err := tx.Exec(ctx, `update sales.return_effect_execution set status=$5,available_at=$6,claimed_at=null,claimed_by=null,claim_token=null,claimed_until=null,last_error_code=$7,result_sha256_hex=null,updated_at=clock_timestamp() where tenant_id=$1 and request_id=$2 and status='claimed' and claimed_by=$3 and claim_token=$4`, work.TenantID, work.RequestID, workerID, work.ClaimToken, outcome, available, code)
	if err != nil || updated.RowsAffected() != 1 {
		if err != nil {
			return err
		}
		return ErrStaleClaim
	}
	return tx.Commit(ctx)
}
