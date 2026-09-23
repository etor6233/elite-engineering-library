package postgres

// AUTHORED binding/persistence glue over existing authoritative owners. This is
// an initial-handover projection, never a shipping, pricing or fiscal engine.
import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"elite.local/enterprise/internal/franchisejourney"
	"github.com/jackc/pgx/v5"
)

type handoverPreparationReader interface {
	QueryRow(context.Context, string, ...any) pgx.Row
}

func readInitialHandover(ctx context.Context, q handoverPreparationReader, tenant, organization, order, key string) (franchisejourney.HandoverPreparation, string, error) {
	var value franchisejourney.HandoverPreparation
	var requestHash string
	var items []byte
	value.Handover.ChecklistItems = []franchisejourney.ChecklistItem{}
	err := q.QueryRow(ctx, `select h.handover_id,h.organization_id,h.order_id,h.customer_principal_id,h.stock_unit_id,h.state,h.version,
      p.order_line_id,p.reservation_id,coalesce(p.payment_attempt_id,''),coalesce(p.funding_receipt_id,''),p.observation_sha256_hex,p.contract_id,p.contract_sha256_hex,p.prepared_by_subject,p.prepared_at,i.request_sha256_hex,
      h.customer_accepted_at,coalesce(h.acceptance_evidence_sha256_hex,''),coalesce(h.checklist_id,''),coalesce(h.checklist_version,0),coalesce(t.title,''),h.checklist_completed_at,coalesce(h.supersedes_handover_id,''),
      coalesce((select jsonb_agg(jsonb_build_object('id',ci.item_id,'ordinal',ci.ordinal,'prompt',ci.prompt,'response_type',ci.response_type,'required',ci.required) order by ci.ordinal,ci.item_id) from sales.delivery_checklist_item ci where ci.tenant_id=h.tenant_id and ci.organization_id=h.organization_id and ci.checklist_id=h.checklist_id and ci.checklist_version=h.checklist_version),'[]'::jsonb)
      from platform.idempotency_record i
      join sales.delivery_handover h on h.tenant_id=i.tenant_id and h.handover_id=i.resource_id
      join sales.delivery_handover_preparation p on p.tenant_id=h.tenant_id and p.handover_id=h.handover_id and p.organization_id=h.organization_id and p.order_id=h.order_id
      join sales.customer_order o on o.tenant_id=h.tenant_id and o.order_id=h.order_id and o.organization_id=h.organization_id and o.customer_principal_id=h.customer_principal_id
      join sales.customer_order_line l on l.tenant_id=p.tenant_id and l.order_id=p.order_id and l.line_id=p.order_line_id and l.allocated_stock_unit_id=h.stock_unit_id
      left join sales.delivery_checklist_template t on t.tenant_id=h.tenant_id and t.organization_id=h.organization_id and t.checklist_id=h.checklist_id and t.checklist_version=h.checklist_version
      where i.tenant_id=$1 and h.organization_id=$2 and h.order_id=$3 and i.scope='initial-handover' and i.idempotency_key=$4
      and i.status='completed' and i.resource_type='delivery-handover' and i.response_code=201
      and i.response_body->>'handover_id'=h.handover_id`, tenant, organization, order, key).Scan(
		&value.Handover.ID, &value.Handover.OrganizationID, &value.Handover.OrderID, &value.Handover.CustomerSubject, &value.Handover.StockUnitID, &value.Handover.State, &value.Handover.Version,
		&value.OrderLineID, &value.ReservationID, &value.PaymentAttemptID, &value.FundingReceiptID, &value.ObservationSHA256, &value.ContractID, &value.ContractSHA256, &value.PreparedBy, &value.PreparedAt, &requestHash,
		&value.Handover.CustomerAcceptedAt, &value.Handover.AcceptanceEvidence, &value.Handover.ChecklistID, &value.Handover.ChecklistVersion, &value.Handover.ChecklistTitle, &value.Handover.ChecklistCompletedAt, &value.Handover.SupersedesHandoverID, &items)
	if errors.Is(err, pgx.ErrNoRows) {
		return value, "", franchisejourney.ErrNotFound
	}
	if err == nil {
		err = json.Unmarshal(items, &value.Handover.ChecklistItems)
	}
	return value, requestHash, err
}

func (r *FranchiseJourney) InitialHandoverResult(ctx context.Context, tenant, organization, order, key string) (franchisejourney.HandoverPreparation, error) {
	value, _, err := readInitialHandover(ctx, r.pool, tenant, organization, order, key)
	return value, err
}

func (r *FranchiseJourney) EvaluateInitialHandoverRelease(ctx context.Context, tenant, organization, handover, observationSHA string, policy franchisejourney.HandoverReleaseContract) (franchisejourney.HandoverReferenceRelease, error) {
	var value franchisejourney.HandoverReferenceRelease
	if !policy.AllowsScope(tenant, organization) {
		return value, franchisejourney.ErrReleaseConditioned
	}
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return value, err
	}
	defer tx.Rollback(ctx)
	command := franchisejourney.PrepareHandoverCommand{OrganizationID: organization, ObservationSHA256: observationSHA}
	var contractID, contractSHA, scope string
	var age int64
	err = tx.QueryRow(ctx, `select order_id,order_line_id,coalesce(payment_attempt_id,''),coalesce(funding_receipt_id,''),contract_id,contract_sha256_hex,contract_scope,maximum_observation_age_ns from sales.delivery_handover_preparation where tenant_id=$1 and organization_id=$2 and handover_id=$3`, tenant, organization, handover).Scan(&command.OrderID, &command.OrderLineID, &command.PaymentAttemptID, &command.FundingReceiptID, &contractID, &contractSHA, &scope, &age)
	if errors.Is(err, pgx.ErrNoRows) {
		return value, franchisejourney.ErrConflict
	}
	if err != nil {
		return value, err
	}
	if contractID != policy.ID || contractSHA != policy.DocumentSHA256 || scope != policy.Scope || age != int64(policy.MaximumObservationAge) {
		return value, franchisejourney.ErrConflict
	}
	facts, err := lockInitialHandoverScope(ctx, tx, tenant, command, policy)
	if err != nil {
		return value, err
	}
	// Existing acceptance owns handover -> order locks. NOWAIT prevents lock
	// inversion while inspecting its committed result after the payment fence.
	var accepted string
	err = tx.QueryRow(ctx, `select handover_id from sales.delivery_handover where tenant_id=$1 and organization_id=$2 and handover_id=$3 and order_id=$4 and customer_principal_id=$5 and stock_unit_id=$6 and state='accepted' and customer_accepted_at is not null and acceptance_evidence_sha256_hex is not null and checklist_completed_at is not null for share nowait`, tenant, organization, handover, command.OrderID, facts.customer, facts.stock).Scan(&accepted)
	if errors.Is(err, pgx.ErrNoRows) {
		return value, franchisejourney.ErrConflict
	}
	if err != nil {
		return value, err
	}
	if err = tx.QueryRow(ctx, `select clock_timestamp()`).Scan(&value.EvaluatedAt); err != nil {
		return value, err
	}
	value.HandoverID = accepted
	value.OrganizationID = organization
	value.ObservationSHA256 = observationSHA
	value.ContractSHA256 = policy.DocumentSHA256
	value.Scope = policy.Scope
	value.Eligible = true
	if err = tx.Commit(ctx); err != nil {
		return franchisejourney.HandoverReferenceRelease{}, err
	}
	return value, nil
}

type handoverScopeFacts struct {
	customer, stock, reservation string
	evidenceAt                   time.Time
	generation                   int64
}

// Order -> payment attempt -> observation is the same lock order as the payment
// reconciler. Time/freshness is evaluated only after these locks are held.
func lockInitialHandoverScope(ctx context.Context, tx pgx.Tx, tenant string, command franchisejourney.PrepareHandoverCommand, policy franchisejourney.HandoverReleaseContract) (handoverScopeFacts, error) {
	var facts handoverScopeFacts
	profileProvider, profileAccount, profileConnection, boundProfile := policy.PaymentBinding()
	if boundProfile && command.FundingReceiptID == "" {
		var active string
		err := tx.QueryRow(ctx, `select connection_id from integration.provider_connection where tenant_id=$1 and connection_id=$2 and provider_code=$3 and state='active' and (organization_id is null or organization_id=$4) for share`, tenant, profileConnection, profileProvider, command.OrganizationID).Scan(&active)
		if errors.Is(err, pgx.ErrNoRows) {
			return facts, franchisejourney.ErrConflict
		}
		if err != nil {
			return facts, err
		}
	}
	var currency, state string
	var total int64
	err := tx.QueryRow(ctx, `select customer_principal_id,currency,total_minor_units,state from sales.customer_order where tenant_id=$1 and organization_id=$2 and order_id=$3 for update`, tenant, command.OrganizationID, command.OrderID).Scan(&facts.customer, &currency, &total, &state)
	if errors.Is(err, pgx.ErrNoRows) {
		return facts, franchisejourney.ErrConflict
	}
	if err != nil {
		return facts, err
	}
	if total <= 0 || (state != "placed" && state != "confirmed" && state != "paid" && state != "allocated") {
		return facts, franchisejourney.ErrConflict
	}
	gross := total
	total, err = orderProviderDue(ctx, tx, tenant, command.OrganizationID, command.OrderID, currency, total)
	if err != nil || total < 0 {
		return facts, franchisejourney.ErrConflict
	}
	if total != gross && !policy.AllowsStoredValueFunding(tenant, command.OrganizationID) {
		return facts, franchisejourney.ErrReleaseConditioned
	}

	var observedAt *time.Time
	if command.FundingReceiptID != "" {
		if command.PaymentAttemptID != "" || total != 0 || !policy.AllowsStoredValueFunding(tenant, command.OrganizationID) {
			return facts, franchisejourney.ErrConflict
		}
		local, e := readLocalFunding(ctx, tx, tenant, command.OrganizationID, command.OrderID, command.FundingReceiptID, command.ObservationSHA256, currency, gross)
		if e != nil {
			return facts, franchisejourney.ErrConflict
		}
		observedAt = &local.Receipt.ObservedAt
		facts.generation = 1
	} else {
		if command.PaymentAttemptID == "" || total <= 0 {
			return facts, franchisejourney.ErrConflict
		}
		var paymentProvider, paymentReference, paymentState, paymentCurrency, paymentOrder string
		var paymentAmount int64
		err = tx.QueryRow(ctx, `select provider_code,coalesce(provider_reference,''),state,currency,amount_minor_units,order_id from payment.payment_attempt where tenant_id=$1 and payment_attempt_id=$2 for update`, tenant, command.PaymentAttemptID).Scan(&paymentProvider, &paymentReference, &paymentState, &paymentCurrency, &paymentAmount, &paymentOrder)
		if errors.Is(err, pgx.ErrNoRows) {
			return facts, franchisejourney.ErrConflict
		}
		if err != nil {
			return facts, err
		}
		if paymentState != "captured" || paymentReference == "" || paymentOrder != command.OrderID || paymentCurrency != currency || paymentAmount != total {
			return facts, franchisejourney.ErrConflict
		}
		var observedOrder, observedOrganization, observedProvider, observedReference, observedCurrency, evidence, providerStatus, account string
		var observedAmount, received, refunded, generation int64
		var live, hold bool
		err = tx.QueryRow(ctx, `select order_id,organization_id,provider_code,provider_reference,currency,amount_minor_units,received_minor_units,refunded_minor_units,provider_status,live_mode,account_ref,coalesce(evidence_sha256_hex,''),observed_at,hold,generation
      from payment.provider_observation where tenant_id=$1 and payment_attempt_id=$2 for update`, tenant, command.PaymentAttemptID).Scan(&observedOrder, &observedOrganization, &observedProvider, &observedReference, &observedCurrency, &observedAmount, &received, &refunded, &providerStatus, &live, &account, &evidence, &observedAt, &hold, &generation)
		if errors.Is(err, pgx.ErrNoRows) {
			return facts, franchisejourney.ErrConflict
		}
		if err != nil {
			return facts, err
		}
		if hold || generation < 1 || observedAt == nil || evidence != command.ObservationSHA256 || observedOrder != command.OrderID || observedOrganization != command.OrganizationID || observedProvider != paymentProvider || observedReference != paymentReference || observedCurrency != currency || observedAmount != total || received != total || refunded != 0 || live != policy.ExpectedLiveMode || !((paymentProvider == "stripe" && providerStatus == "succeeded") || (paymentProvider == "mercadopago" && providerStatus == "approved")) || account == "" {
			return facts, franchisejourney.ErrConflict
		}
		if boundProfile && command.FundingReceiptID == "" {
			if observedProvider != profileProvider || account != profileAccount {
				return facts, franchisejourney.ErrConflict
			}
			var exact bool
			err = tx.QueryRow(ctx, `select exists(select 1 from payment.provider_checkout where tenant_id=$1 and payment_attempt_id=$2 and provider_code=$3 and account_ref=$4 and connection_id=$5 and live_mode=$6)`, tenant, command.PaymentAttemptID, profileProvider, profileAccount, profileConnection, policy.ExpectedLiveMode).Scan(&exact)
			if err != nil {
				return facts, err
			}
			if !exact {
				return facts, franchisejourney.ErrConflict
			}
		}
		facts.generation = generation

	}
	facts.evidenceAt = *observedAt
	var customerStatus string
	err = tx.QueryRow(ctx, `select status from crm.customer_profile where tenant_id=$1 and customer_principal_id=$2 for share`, tenant, facts.customer).Scan(&customerStatus)
	if errors.Is(err, pgx.ErrNoRows) {
		return facts, franchisejourney.ErrConflict
	}
	if err != nil {
		return facts, err
	}
	if customerStatus != "active" {
		return facts, franchisejourney.ErrConflict
	}
	var variant string
	var lineCount, paymentCount int
	err = tx.QueryRow(ctx, `select l.variant_id,l.allocated_stock_unit_id,(select count(*) from sales.customer_order_line where tenant_id=$1 and order_id=$2),(select count(*) from payment.payment_attempt where tenant_id=$1 and order_id=$2 and state not in ('failed','refunded')) from sales.customer_order_line l where l.tenant_id=$1 and l.order_id=$2 and l.line_id=$3 and l.quantity=1 and l.allocated_stock_unit_id is not null for share of l`, tenant, command.OrderID, command.OrderLineID).Scan(&variant, &facts.stock, &lineCount, &paymentCount)
	if errors.Is(err, pgx.ErrNoRows) {
		return facts, franchisejourney.ErrConflict
	}
	if err != nil {
		return facts, err
	}
	expectedPayments := 1
	if command.FundingReceiptID != "" {
		expectedPayments = 0
	}
	if lineCount != 1 || paymentCount != expectedPayments {
		return facts, franchisejourney.ErrConflict
	}
	var stock string
	err = tx.QueryRow(ctx, `select stock_unit_id from inventory.stock_unit where tenant_id=$1 and stock_unit_id=$2 and organization_id=$3 and variant_id=$4 and state='reserved' for share`, tenant, facts.stock, command.OrganizationID, variant).Scan(&stock)
	if errors.Is(err, pgx.ErrNoRows) {
		return facts, franchisejourney.ErrConflict
	}
	if err != nil {
		return facts, err
	}
	var expires *time.Time
	err = tx.QueryRow(ctx, `select reservation_id,expires_at from inventory.serial_reservation where tenant_id=$1 and organization_id=$2 and stock_unit_id=$3 and variant_id=$4 and demand_kind='customer-order' and demand_id=$5 and demand_line_id=$6 and status='reservation' for share`, tenant, command.OrganizationID, facts.stock, variant, command.OrderID, command.OrderLineID).Scan(&facts.reservation, &expires)
	if errors.Is(err, pgx.ErrNoRows) {
		return facts, franchisejourney.ErrConflict
	}
	if err != nil {
		return facts, err
	}
	var now time.Time
	if err = tx.QueryRow(ctx, `select clock_timestamp()`).Scan(&now); err != nil {
		return facts, err
	}
	if observedAt.After(now) || now.Sub(*observedAt) > policy.MaximumObservationAge || (expires != nil && !expires.After(now)) {
		return facts, franchisejourney.ErrConflict
	}
	return facts, nil
}

func (r *FranchiseJourney) PrepareInitialHandover(ctx context.Context, tenant, actor, handoverID, eventID string, command franchisejourney.PrepareHandoverCommand, policy franchisejourney.HandoverReleaseContract, requestHash string) (franchisejourney.HandoverPreparation, bool, error) {
	var empty franchisejourney.HandoverPreparation
	if !policy.AllowsScope(tenant, command.OrganizationID) || tenant == "" || actor == "" || handoverID == "" || eventID == "" || len(requestHash) != 64 {
		return empty, false, franchisejourney.ErrReleaseConditioned
	}
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return empty, false, err
	}
	defer tx.Rollback(ctx)
	inserted, err := tx.Exec(ctx, `insert into platform.idempotency_record(tenant_id,scope,idempotency_key,request_sha256_hex,status,locked_until,expires_at)
      values($1,'initial-handover',$2,$3,'processing',clock_timestamp()+interval '30 seconds',clock_timestamp()+interval '24 hours') on conflict do nothing`, tenant, command.IdempotencyKey, requestHash)
	if err != nil {
		return empty, false, err
	}
	if inserted.RowsAffected() == 0 {
		value, storedHash, readErr := readInitialHandover(ctx, tx, tenant, command.OrganizationID, command.OrderID, command.IdempotencyKey)
		if readErr != nil || storedHash != requestHash {
			return empty, false, franchisejourney.ErrConflict
		}
		if err = tx.Commit(ctx); err != nil {
			return empty, false, err
		}
		return value, true, nil
	}
	facts, err := lockInitialHandoverScope(ctx, tx, tenant, command, policy)
	if err != nil {
		return empty, false, err
	}
	var existing bool
	if err = tx.QueryRow(ctx, `select exists(select 1 from sales.delivery_handover_preparation where tenant_id=$1 and organization_id=$2 and order_id=$3 and order_line_id=$4)`, tenant, command.OrganizationID, command.OrderID, command.OrderLineID).Scan(&existing); err != nil {
		return empty, false, err
	}
	if existing {
		return empty, false, franchisejourney.ErrConflict
	}
	_, err = tx.Exec(ctx, `insert into sales.delivery_handover(tenant_id,handover_id,organization_id,order_id,customer_principal_id,stock_unit_id,state,version)values($1,$2,$3,$4,$5,$6,'prepared',1)`, tenant, handoverID, command.OrganizationID, command.OrderID, facts.customer, facts.stock)
	if err != nil {
		return empty, false, err
	}
	_, err = tx.Exec(ctx, `insert into sales.delivery_handover_preparation(tenant_id,handover_id,organization_id,order_id,order_line_id,reservation_id,payment_attempt_id,observation_sha256_hex,contract_id,contract_sha256_hex,contract_scope,maximum_observation_age_ns,prepared_by_subject,funding_receipt_id)values($1,$2,$3,$4,$5,$6,nullif($7,''),$8,$9,$10,$11,$12,$13,nullif($14,''))`, tenant, handoverID, command.OrganizationID, command.OrderID, command.OrderLineID, facts.reservation, command.PaymentAttemptID, command.ObservationSHA256, policy.ID, policy.DocumentSHA256, policy.Scope, int64(policy.MaximumObservationAge), actor, command.FundingReceiptID)
	if err != nil {
		return empty, false, err
	}
	if err = writeJourneyOutbox(ctx, tx, tenant, eventID, "delivery-handover", handoverID, 1, "delivery-handover.prepared", map[string]any{"organization_id": command.OrganizationID, "order_id": command.OrderID, "order_line_id": command.OrderLineID, "payment_attempt_id": command.PaymentAttemptID, "funding_receipt_id": command.FundingReceiptID, "observation_sha256": command.ObservationSHA256, "contract_id": policy.ID, "contract_sha256": policy.DocumentSHA256, "actor_subject": actor}); err != nil {
		return empty, false, err
	}
	result, err := tx.Exec(ctx, `update platform.idempotency_record set status='completed',response_code=201,response_body=jsonb_build_object('handover_id',$3::text),resource_type='delivery-handover',resource_id=$3,locked_until=null where tenant_id=$1 and scope='initial-handover' and idempotency_key=$2 and status='processing'`, tenant, command.IdempotencyKey, handoverID)
	if err != nil {
		return empty, false, err
	}
	if result.RowsAffected() != 1 {
		return empty, false, fmt.Errorf("handover receipt not completed: %w", franchisejourney.ErrConflict)
	}
	value, _, err := readInitialHandover(ctx, tx, tenant, command.OrganizationID, command.OrderID, command.IdempotencyKey)
	if err != nil {
		return empty, false, err
	}
	if err = tx.Commit(ctx); err != nil {
		return empty, false, err
	}
	return value, false, nil
}
