package postgres

// AUTHORED persistence/transaction glue around the existing Commerce payment,
// integration inbox and outbound fence; no new accounting or capture algorithm.
import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"elite.local/enterprise/internal/paymentbridge"
	"elite.local/enterprise/internal/providerintegration"
	officialpayments "example.com/elite/official-payment-webhooks/officialpayments"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PaymentCheckoutStore struct{ pool *pgxpool.Pool }

func NewPaymentCheckoutStore(pool *pgxpool.Pool) *PaymentCheckoutStore {
	return &PaymentCheckoutStore{pool: pool}
}

func (s *PaymentCheckoutStore) PendingRequests(ctx context.Context, c paymentbridge.Scope, limit int) ([]paymentbridge.Request, error) {
	if s == nil || s.pool == nil || c.Validate() != nil || limit < 1 || limit > 100 {
		return nil, paymentbridge.ErrCheckoutConflict
	}
	rows, err := s.pool.Query(ctx, `select p.tenant_id,o.organization_id,p.payment_attempt_id,p.order_id,o.customer_principal_id,p.provider_code,p.currency,p.amount_minor_units,
 (select date_trunc('second',min(e.occurred_at))+interval '23 hours' from platform.outbox_event e where e.tenant_id=p.tenant_id and e.aggregate_id=p.payment_attempt_id and e.event_type='payment.requested' and e.schema_version=1)
 from payment.payment_attempt p join sales.customer_order o using(tenant_id,order_id)
 join integration.provider_connection conn on conn.tenant_id=p.tenant_id and conn.connection_id=$4 and conn.provider_code=p.provider_code and conn.state='active' and (conn.organization_id is null or conn.organization_id=o.organization_id)
 left join payment.provider_checkout checkout on checkout.tenant_id=p.tenant_id and checkout.payment_attempt_id=p.payment_attempt_id
 where p.tenant_id=$1 and p.provider_code=$2 and o.organization_id=$3 and p.currency=$5 and p.state='created'
 and checkout.session_id is null and p.amount_minor_units=(select provider_due_minor_units from payment.order_funding f where f.tenant_id=o.tenant_id and f.order_id=o.order_id and f.organization_id=o.organization_id) and p.currency=o.currency
 and o.state in ('placed','confirmed','allocated') and o.customer_principal_id is not null
 and exists(select 1 from platform.outbox_event e where e.tenant_id=p.tenant_id and e.aggregate_id=p.payment_attempt_id and e.event_type='payment.requested' and e.schema_version=1)
 and (select date_trunc('second',min(e.occurred_at))+interval '23 hours' from platform.outbox_event e where e.tenant_id=p.tenant_id and e.aggregate_id=p.payment_attempt_id and e.event_type='payment.requested' and e.schema_version=1)>clock_timestamp()+interval '31 minutes'
 and not exists(select 1 from communication.outbound_delivery f where f.tenant_id=p.tenant_id and f.channel_code='payment_'||p.provider_code and f.delivery_key=p.payment_attempt_id::text and f.state in ('accepted','unknown','failed_terminal'))
 order by p.updated_at,p.payment_attempt_id limit $6`, c.TenantID, c.ProviderCode, c.OrganizationID, c.ConnectionID, c.Currency, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []paymentbridge.Request{}
	for rows.Next() {
		var r paymentbridge.Request
		if err = rows.Scan(&r.TenantID, &r.OrganizationID, &r.PaymentAttemptID, &r.OrderID, &r.CustomerSubject, &r.ProviderCode, &r.Currency, &r.AmountMinor, &r.CheckoutExpiresAt); err != nil {
			return nil, err
		}
		r.CheckoutExpiresAt = r.CheckoutExpiresAt.UTC()
		result = append(result, r)
	}
	return result, rows.Err()
}

func lockCheckoutRequest(ctx context.Context, tx pgx.Tx, c paymentbridge.Scope, attempt string) (paymentbridge.Request, error) {
	var r paymentbridge.Request
	var order string
	err := tx.QueryRow(ctx, `select order_id from payment.payment_attempt where tenant_id=$1 and payment_attempt_id=$2 and provider_code=$3`, c.TenantID, attempt, c.ProviderCode).Scan(&order)
	if errors.Is(err, pgx.ErrNoRows) {
		return r, paymentbridge.ErrCheckoutConflict
	}
	if err != nil {
		return r, err
	}
	var organization, customer, currency string
	var total int64
	err = tx.QueryRow(ctx, `select organization_id,customer_principal_id,currency,total_minor_units from sales.customer_order where tenant_id=$1 and order_id=$2 and organization_id=$3 and state in ('placed','confirmed','allocated','delivered') for update`, c.TenantID, order, c.OrganizationID).Scan(&organization, &customer, &currency, &total)
	if errors.Is(err, pgx.ErrNoRows) {
		return r, paymentbridge.ErrCheckoutConflict
	}
	if err != nil {
		return r, err
	}
	err = tx.QueryRow(ctx, `select tenant_id,payment_attempt_id,order_id,provider_code,currency,amount_minor_units from payment.payment_attempt where tenant_id=$1 and payment_attempt_id=$2 and order_id=$3 for update`, c.TenantID, attempt, order).Scan(&r.TenantID, &r.PaymentAttemptID, &r.OrderID, &r.ProviderCode, &r.Currency, &r.AmountMinor)
	if err != nil {
		return r, err
	}
	total, err = orderProviderDue(ctx, tx, c.TenantID, c.OrganizationID, order, currency, total)
	if err != nil {
		return r, err
	}
	if r.ProviderCode != c.ProviderCode || r.Currency != c.Currency || r.Currency != currency || r.AmountMinor != total || r.AmountMinor <= 0 {
		return r, paymentbridge.ErrCheckoutConflict
	}
	r.OrganizationID = organization
	r.CustomerSubject = customer
	err = tx.QueryRow(ctx, `select date_trunc('second',min(occurred_at))+interval '23 hours' from platform.outbox_event where tenant_id=$1 and aggregate_id=$2 and event_type='payment.requested' and schema_version=1`, c.TenantID, attempt).Scan(&r.CheckoutExpiresAt)
	if err != nil {
		return r, err
	}
	r.CheckoutExpiresAt = r.CheckoutExpiresAt.UTC()
	return r, nil
}

func checkoutConnection(ctx context.Context, tx pgx.Tx, c paymentbridge.Scope) error {
	var found string
	err := tx.QueryRow(ctx, `select connection_id from integration.provider_connection where tenant_id=$1 and connection_id=$2 and provider_code=$3 and state='active' and (organization_id is null or organization_id=$4) for share`, c.TenantID, c.ConnectionID, c.ProviderCode, c.OrganizationID).Scan(&found)
	if errors.Is(err, pgx.ErrNoRows) {
		return providerintegration.ErrConnection
	}
	return err
}

func (s *PaymentCheckoutStore) BindRequest(ctx context.Context, c paymentbridge.Scope, r paymentbridge.Request, hash string) error {
	if c.Validate() != nil || len(hash) != 64 {
		return paymentbridge.ErrCheckoutConflict
	}
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err = checkoutConnection(ctx, tx, c); err != nil {
		return err
	}
	actual, err := lockCheckoutRequest(ctx, tx, c, r.PaymentAttemptID)
	if err != nil {
		return err
	}
	if actual != r {
		return paymentbridge.ErrCheckoutConflict
	}
	var dispatchable bool
	err = tx.QueryRow(ctx, `select p.state='created' and o.state in ('placed','confirmed','allocated') from payment.payment_attempt p join sales.customer_order o using(tenant_id,order_id) where p.tenant_id=$1 and p.payment_attempt_id=$2`, c.TenantID, r.PaymentAttemptID).Scan(&dispatchable)
	if err != nil {
		return err
	}
	if !dispatchable || !r.CheckoutExpiresAt.After(time.Now().Add(31*time.Minute)) {
		return paymentbridge.ErrCheckoutConflict
	}
	_, err = tx.Exec(ctx, `insert into payment.provider_checkout(tenant_id,payment_attempt_id,provider_code,connection_id,account_ref,request_sha256_hex,live_mode,expires_at) values($1,$2,$3,$4,$5,$6,$7,$8) on conflict do nothing`, c.TenantID, r.PaymentAttemptID, c.ProviderCode, c.ConnectionID, c.AccountRef, hash, c.LiveMode, r.CheckoutExpiresAt)
	if err != nil {
		return err
	}
	var exact bool
	err = tx.QueryRow(ctx, `select provider_code=$3 and connection_id=$4 and account_ref=$5 and request_sha256_hex=$6 and live_mode=$7 and expires_at=$8 from payment.provider_checkout where tenant_id=$1 and payment_attempt_id=$2 for update`, c.TenantID, r.PaymentAttemptID, c.ProviderCode, c.ConnectionID, c.AccountRef, hash, c.LiveMode, r.CheckoutExpiresAt).Scan(&exact)
	if err != nil {
		return err
	}
	if !exact {
		return paymentbridge.ErrCheckoutConflict
	}
	return tx.Commit(ctx)
}

func (s *PaymentCheckoutStore) SaveCheckout(ctx context.Context, c paymentbridge.Scope, r paymentbridge.Request, hash string, result paymentbridge.Checkout) error {
	if result.OrderID != r.OrderID || result.PaymentAttemptID != r.PaymentAttemptID || result.Currency != r.Currency || result.AmountMinor != r.AmountMinor || result.LiveMode != c.LiveMode || result.AccountRef != c.AccountRef || !result.ExpiresAt.Equal(r.CheckoutExpiresAt) {
		return paymentbridge.ErrCheckoutConflict
	}
	// Stripe only returns the hosted URL while a session is active. A completed
	// GET can recover the durable identity after a lost POST response without a URL.
	urlOptional := result.ProviderCode == "stripe" && (result.Status == "complete" || result.Status == "expired")
	if c.Validate() != nil || result.ProviderCode != c.ProviderCode || result.SessionID == "" || len(result.SessionID) > 200 || (result.URL == "" && !urlOptional) || len(result.URL) > 4096 || result.Status == "" {
		return paymentbridge.ErrCheckoutConflict
	}
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err = checkoutConnection(ctx, tx, c); err != nil {
		return err
	}
	actual, err := lockCheckoutRequest(ctx, tx, c, r.PaymentAttemptID)
	if err != nil {
		return err
	}
	if actual != r {
		return paymentbridge.ErrCheckoutConflict
	}
	// A known session cannot be replaced by a later completion or a different key.
	updated, err := tx.Exec(ctx, `update payment.provider_checkout set session_id=$3,checkout_url=coalesce(nullif($4,''),checkout_url),checkout_state=$5,payment_reference=coalesce(nullif($11,''),payment_reference),updated_at=clock_timestamp() where tenant_id=$1 and payment_attempt_id=$2 and request_sha256_hex=$6 and provider_code=$7 and connection_id=$8 and account_ref=$9 and live_mode=$10 and (session_id is null or session_id=$3) and (payment_reference is null or $11='' or payment_reference=$11)`, c.TenantID, r.PaymentAttemptID, result.SessionID, result.URL, result.Status, hash, c.ProviderCode, c.ConnectionID, c.AccountRef, c.LiveMode, result.PaymentReference)
	if err != nil {
		return err
	}
	if updated.RowsAffected() != 1 {
		return paymentbridge.ErrCheckoutConflict
	}
	return tx.Commit(ctx)
}

func (s *PaymentCheckoutStore) Checkout(ctx context.Context, c paymentbridge.Scope, attempt string) (paymentbridge.Request, paymentbridge.Checkout, error) {
	var r paymentbridge.Request
	var result paymentbridge.Checkout
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return r, result, err
	}
	defer tx.Rollback(ctx)
	if err = checkoutConnection(ctx, tx, c); err != nil {
		return r, result, err
	}
	r, err = lockCheckoutRequest(ctx, tx, c, attempt)
	if err != nil {
		return r, result, err
	}
	err = tx.QueryRow(ctx, `select provider_code,session_id,coalesce(checkout_url,''),checkout_state,coalesce(payment_reference,''),expires_at from payment.provider_checkout where tenant_id=$1 and payment_attempt_id=$2 and connection_id=$3 and account_ref=$4 and live_mode=$5 and session_id is not null`, c.TenantID, attempt, c.ConnectionID, c.AccountRef, c.LiveMode).Scan(&result.ProviderCode, &result.SessionID, &result.URL, &result.Status, &result.PaymentReference, &result.ExpiresAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return r, result, paymentbridge.ErrCheckoutConflict
	}
	if err != nil {
		return r, result, err
	}
	result.OrderID = r.OrderID
	result.PaymentAttemptID = r.PaymentAttemptID
	result.Currency = r.Currency
	result.AmountMinor = r.AmountMinor
	result.AccountRef = c.AccountRef
	result.LiveMode = c.LiveMode
	return r, result, tx.Commit(ctx)
}

func (s *PaymentCheckoutStore) BeginObservation(ctx context.Context, c paymentbridge.Scope, attempt, reference string) (paymentbridge.Request, int64, error) {
	var empty paymentbridge.Request
	if c.Validate() != nil || reference == "" || len(reference) > 200 {
		return empty, 0, paymentbridge.ErrCheckoutConflict
	}
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return empty, 0, err
	}
	defer tx.Rollback(ctx)
	if err = checkoutConnection(ctx, tx, c); err != nil {
		return empty, 0, err
	}
	r, err := lockCheckoutRequest(ctx, tx, c, attempt)
	if err != nil {
		return r, 0, err
	}
	var bound bool
	err = tx.QueryRow(ctx, `select exists(select 1 from payment.provider_checkout where tenant_id=$1 and payment_attempt_id=$2 and connection_id=$3 and provider_code=$4 and account_ref=$5 and live_mode=$6)`, c.TenantID, attempt, c.ConnectionID, c.ProviderCode, c.AccountRef, c.LiveMode).Scan(&bound)
	if err != nil {
		return r, 0, err
	}
	if !bound {
		return r, 0, paymentbridge.ErrCheckoutConflict
	}
	_, err = tx.Exec(ctx, `insert into payment.provider_observation(tenant_id,payment_attempt_id,order_id,organization_id,provider_code,provider_reference,currency,amount_minor_units,live_mode,account_ref) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) on conflict(tenant_id,payment_attempt_id) do nothing`, c.TenantID, attempt, r.OrderID, r.OrganizationID, c.ProviderCode, reference, r.Currency, r.AmountMinor, c.LiveMode, c.AccountRef)
	if err != nil {
		return r, 0, err
	}
	var generation int64
	err = tx.QueryRow(ctx, `update payment.provider_observation set provider_reference=$4,hold=true,hold_reason='GET_IN_PROGRESS',generation=generation+1,updated_at=clock_timestamp() where tenant_id=$1 and payment_attempt_id=$2 and provider_code=$3 and (provider_reference=$4 or (observed_at is null and hold)) and account_ref=$5 and live_mode=$6 returning generation`, c.TenantID, attempt, c.ProviderCode, reference, c.AccountRef, c.LiveMode).Scan(&generation)
	if errors.Is(err, pgx.ErrNoRows) {
		return r, 0, paymentbridge.ErrCheckoutConflict
	}
	if err != nil {
		return r, 0, err
	}
	return r, generation, tx.Commit(ctx)
}

func (s *PaymentCheckoutStore) RecordObservation(ctx context.Context, c paymentbridge.Scope, r paymentbridge.Request, generation int64, p officialpayments.PaymentSnapshot) error {
	if c.Validate() != nil || generation < 1 {
		return paymentbridge.ErrCheckoutConflict
	}
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err = checkoutConnection(ctx, tx, c); err != nil {
		return err
	}
	actual, err := lockCheckoutRequest(ctx, tx, c, r.PaymentAttemptID)
	if err != nil {
		return err
	}
	if actual != r {
		return paymentbridge.ErrCheckoutConflict
	}
	var storedGeneration int64
	var previousRefunded int64
	var reference string
	err = tx.QueryRow(ctx, `select generation,provider_reference,refunded_minor_units from payment.provider_observation where tenant_id=$1 and payment_attempt_id=$2 for update`, c.TenantID, r.PaymentAttemptID).Scan(&storedGeneration, &reference, &previousRefunded)
	if err != nil {
		return err
	}
	if generation != storedGeneration {
		return paymentbridge.ErrObservationStale
	}
	mismatch := p.ProviderCode != c.ProviderCode || p.ProviderReference != reference || p.OrderID != r.OrderID || p.Currency != r.Currency || p.AmountMinor != r.AmountMinor || p.LiveMode != c.LiveMode || p.ReceivedMinor < 0 || p.ReceivedMinor > p.AmountMinor || p.RefundedMinor < 0 || p.RefundedMinor > p.ReceivedMinor
	if p.ProviderCode == "mercadopago" && p.CollectorID != c.AccountRef {
		mismatch = true
	}
	if p.PaymentAttemptID != "" && p.PaymentAttemptID != r.PaymentAttemptID {
		mismatch = true
	}
	if p.RefundedMinor < previousRefunded {
		mismatch = true
	}
	reason := "PROVIDER_NOT_CAPTURED"
	target := ""
	if mismatch {
		reason = "PROVIDER_BINDING_MISMATCH"
	} else if p.Disputed {
		reason = "PROVIDER_DISPUTED"
		target = "disputed"
	} else if p.RefundedMinor > 0 {
		reason = "PARTIAL_PROVIDER_REFUND"
		if p.RefundedMinor == r.AmountMinor && p.ReceivedMinor == r.AmountMinor {
			reason = "PROVIDER_REFUNDED"
			target = "refunded"
		}
	} else if (p.ProviderCode == "stripe" && p.Status == "succeeded" || p.ProviderCode == "mercadopago" && p.Status == "approved") && p.ReceivedMinor == r.AmountMinor {
		reason = "MATCHED_CAPTURE"
		target = "captured"
	}
	var financialState, boundReference string
	if err = tx.QueryRow(ctx, `select state,coalesce(provider_reference,'') from payment.payment_attempt where tenant_id=$1 and payment_attempt_id=$2`, c.TenantID, r.PaymentAttemptID).Scan(&financialState, &boundReference); err != nil {
		return err
	}
	if boundReference != "" && boundReference != reference {
		mismatch = true
		target = ""
		reason = "BOUND_PROVIDER_REFERENCE_MISMATCH"
	}
	if (financialState == "refunded" || financialState == "disputed" || financialState == "failed") && target != "" && target != financialState {
		target = ""
		reason = "TERMINAL_STATE_REQUIRES_REVIEW"
	}
	hold := reason != "MATCHED_CAPTURE"
	evidence := p.SHA256()
	if mismatch {
		_, err = tx.Exec(ctx, `update payment.provider_observation set hold=true,hold_reason=$3,updated_at=clock_timestamp() where tenant_id=$1 and payment_attempt_id=$2`, c.TenantID, r.PaymentAttemptID, reason)
	} else {
		_, err = tx.Exec(ctx, `update payment.provider_observation set received_minor_units=$3,refunded_minor_units=$4,provider_status=$5,evidence_sha256_hex=$6,observed_at=clock_timestamp(),hold=$7,hold_reason=$8,updated_at=clock_timestamp() where tenant_id=$1 and payment_attempt_id=$2`, c.TenantID, r.PaymentAttemptID, p.ReceivedMinor, p.RefundedMinor, p.Status, evidence, hold, reason)
	}
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `insert into payment.provider_observation_event(tenant_id,payment_attempt_id,generation,evidence_sha256_hex,result_code) values($1,$2,$3,$4,$5)`, c.TenantID, r.PaymentAttemptID, generation, evidence, reason)
	if err != nil {
		return err
	}
	// These are provider vocabulary projections, not simulated authorization or
	// capture operations. One observed successful payment may skip earlier notices.
	if !mismatch && target != "" {
		_, err = tx.Exec(ctx, `with changed as(update payment.payment_attempt set state=$3,provider_reference=$4,version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and payment_attempt_id=$2 and (state<>$3 or provider_reference is distinct from $4) returning version)
   insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)
   select $1,gen_random_uuid(),'payment',$2,version,'payment.provider-observed',1,clock_timestamp(),jsonb_build_object('state',$3::text,'evidence_sha256',$5::text,'provider_code',$6::text) from changed`, c.TenantID, r.PaymentAttemptID, target, reference, evidence, c.ProviderCode)
		if err != nil {
			return err
		}
	}
	reconciliationState := "matched"
	if mismatch {
		reconciliationState = "amount-mismatch"
	} else if hold {
		reconciliationState = "state-mismatch"
	}
	body, _ := json.Marshal(map[string]any{"reason": reason, "generation": generation, "evidence_sha256": evidence})
	_, err = tx.Exec(ctx, `insert into integration.reconciliation_item(tenant_id,reconciliation_id,provider_code,resource_type,internal_id,external_id,state,evidence,observed_at) values($1,gen_random_uuid(),$2,'payment',$3,$4,$5,$6,clock_timestamp())`, c.TenantID, c.ProviderCode, r.PaymentAttemptID, reference, reconciliationState, body)
	if err != nil {
		return err
	}
	if err = tx.Commit(ctx); err != nil {
		return err
	}
	if mismatch {
		return paymentbridge.ErrPaymentMismatch
	}
	return nil
}

// AcceptWebhook shares the established durable inbox transaction. Only a new
// authenticated notification advances a known observation generation/hold.
func (s *PaymentCheckoutStore) AcceptWebhook(ctx context.Context, value providerintegration.Receipt) (bool, error) {
	var body paymentbridge.InboxPayload
	if value.EventType != "payment.reconciliation-requested.v1" || json.Unmarshal(value.Payload, &body) != nil || body.Schema != "elite-payment-notification/v1" || body.Notification.ProviderCode != value.ProviderCode || body.Notification.ProviderEventID != value.ProviderEventID || body.Notification.BodySHA256 != value.BodyHash {
		return false, providerintegration.ErrInvalid
	}
	n := body.Notification
	if strings.TrimSpace(n.ProviderReference) == "" {
		return false, providerintegration.ErrInvalid
	}
	return NewProviderIntegration(s.pool).acceptWebhookWithEffect(ctx, value, "payment-provider-events", func(ctx context.Context, tx pgx.Tx) error {
		_, err := tx.Exec(ctx, `update payment.provider_observation o set hold=true,hold_reason='PROVIDER_NOTIFICATION_PENDING',generation=o.generation+1,updated_at=clock_timestamp() from payment.provider_checkout c where o.tenant_id=$1 and o.provider_code=$2 and (o.provider_reference=$3 or c.session_id=$3) and c.tenant_id=o.tenant_id and c.payment_attempt_id=o.payment_attempt_id and c.connection_id=$4`, value.TenantID, value.ProviderCode, n.ProviderReference, value.ConnectionID)
		return err
	})
}
