package postgres

// AUTHORED persistence glue over the existing handover/payment owners. The
// receipt records an atomic commercial checkpoint, never a physical shipment.
import (
	"context"
	"errors"
	"time"

	"elite.local/enterprise/internal/franchisejourney"
	"github.com/jackc/pgx/v5"
)

const commercialReleaseColumns = `r.release_id,r.organization_id,r.handover_id,r.order_id,coalesce(r.payment_attempt_id,''),coalesce(r.funding_receipt_id,''),r.observation_sha256_hex,r.observation_generation,r.handover_version,r.acceptance_sha256_hex,r.checklist_id,r.checklist_version,r.contract_id,r.contract_sha256_hex,r.effect,r.released_by_subject,r.recorded_at,r.valid_until`

func scanCommercialRelease(row pgx.Row) (franchisejourney.CommercialReleaseReceipt, error) {
	var r franchisejourney.CommercialReleaseReceipt
	err := row.Scan(&r.ID, &r.OrganizationID, &r.HandoverID, &r.OrderID, &r.PaymentAttemptID, &r.FundingReceiptID, &r.ObservationSHA256, &r.ObservationGeneration, &r.HandoverVersion, &r.AcceptanceSHA256, &r.ChecklistID, &r.ChecklistVersion, &r.ContractID, &r.ContractSHA256, &r.Effect, &r.ReleasedBy, &r.RecordedAt, &r.ValidUntil)
	if errors.Is(err, pgx.ErrNoRows) {
		err = franchisejourney.ErrNotFound
	}
	return r, err
}

func readCommercialRelease(ctx context.Context, q handoverPreparationReader, tenant, organization, handover, key string) (franchisejourney.CommercialReleaseReceipt, error) {
	return scanCommercialRelease(q.QueryRow(ctx, `select `+commercialReleaseColumns+` from sales.commercial_release_receipt r
 join platform.idempotency_record i on i.tenant_id=r.tenant_id and i.resource_id=r.release_id
 where r.tenant_id=$1 and r.organization_id=$2 and r.handover_id=$3 and i.scope='commercial-release' and i.idempotency_key=$4
 and i.status='completed' and i.resource_type='commercial-release' and i.response_code=201 and i.response_body->>'release_id'=r.release_id`, tenant, organization, handover, key))
}

func (r *FranchiseJourney) InitialCommercialReleaseResult(ctx context.Context, tenant, organization, handover, key string) (franchisejourney.CommercialReleaseReceipt, error) {
	return readCommercialRelease(ctx, r.pool, tenant, organization, handover, key)
}

// Derive the order/payment from immutable preparation. The caller cannot supply
// a different customer, stock unit, amount, payment reference or accepted state.
func lockCommercialRelease(ctx context.Context, tx pgx.Tx, tenant, organization, handover, evidence string, p franchisejourney.HandoverReleaseContract) (franchisejourney.CommercialReleaseReceipt, error) {
	var v franchisejourney.CommercialReleaseReceipt
	if !p.AllowsCommercialRelease(tenant, organization) {
		return v, franchisejourney.ErrReleaseConditioned
	}
	c := franchisejourney.PrepareHandoverCommand{OrganizationID: organization, ObservationSHA256: evidence}
	var id, sha, scope string
	var age int64
	err := tx.QueryRow(ctx, `select order_id,order_line_id,coalesce(payment_attempt_id,''),coalesce(funding_receipt_id,''),contract_id,contract_sha256_hex,contract_scope,maximum_observation_age_ns from sales.delivery_handover_preparation where tenant_id=$1 and organization_id=$2 and handover_id=$3`, tenant, organization, handover).Scan(&c.OrderID, &c.OrderLineID, &c.PaymentAttemptID, &c.FundingReceiptID, &id, &sha, &scope, &age)
	if errors.Is(err, pgx.ErrNoRows) {
		return v, franchisejourney.ErrConflict
	}
	if err != nil {
		return v, err
	}
	if id != p.ID || sha != p.DocumentSHA256 || scope != p.Scope || age != int64(p.MaximumObservationAge) {
		return v, franchisejourney.ErrConflict
	}
	facts, err := lockInitialHandoverScope(ctx, tx, tenant, c, p)
	if err != nil {
		return v, err
	}
	// Acceptance takes handover -> order. NOWAIT avoids inversion; a caller may
	// retry after the other transaction has completed without changing its key.
	err = tx.QueryRow(ctx, `select version,acceptance_evidence_sha256_hex,checklist_id,checklist_version from sales.delivery_handover
 where tenant_id=$1 and organization_id=$2 and handover_id=$3 and order_id=$4 and customer_principal_id=$5 and stock_unit_id=$6
 and state='accepted' and customer_accepted_at is not null and acceptance_evidence_sha256_hex is not null and checklist_completed_at is not null for share nowait`, tenant, organization, handover, c.OrderID, facts.customer, facts.stock).Scan(&v.HandoverVersion, &v.AcceptanceSHA256, &v.ChecklistID, &v.ChecklistVersion)
	if errors.Is(err, pgx.ErrNoRows) {
		return v, franchisejourney.ErrConflict
	}
	if err != nil {
		return v, err
	}
	observed := facts.evidenceAt
	v.ObservationGeneration = facts.generation
	var expires *time.Time
	err = tx.QueryRow(ctx, `select expires_at,clock_timestamp() from inventory.serial_reservation where tenant_id=$1 and reservation_id=$2`, tenant, facts.reservation).Scan(&expires, &v.RecordedAt)
	if err != nil {
		return v, err
	}
	v.ValidUntil = observed.Add(p.MaximumObservationAge)
	if expires != nil && expires.Before(v.ValidUntil) {
		v.ValidUntil = *expires
	}
	// Time is sampled after every relevant lock. A snapshot that became stale
	// while waiting must never be committed as a fresh release.
	if observed.After(v.RecordedAt) || !v.ValidUntil.After(v.RecordedAt) {
		return v, franchisejourney.ErrConflict
	}
	v.OrganizationID = organization
	v.HandoverID = handover
	v.OrderID = c.OrderID
	v.PaymentAttemptID = c.PaymentAttemptID
	v.FundingReceiptID = c.FundingReceiptID
	v.ObservationSHA256 = evidence
	v.ContractID = p.ID
	v.ContractSHA256 = p.DocumentSHA256
	v.Effect = franchisejourney.CommercialReleaseEffect
	return v, nil
}

func (r *FranchiseJourney) CommitInitialCommercialRelease(ctx context.Context, tenant, actor, id, event string, c franchisejourney.CommitCommercialReleaseCommand, p franchisejourney.HandoverReleaseContract, hash string) (franchisejourney.CommercialReleaseReceipt, bool, error) {
	var empty franchisejourney.CommercialReleaseReceipt
	if !p.AllowsCommercialRelease(tenant, c.OrganizationID) || tenant == "" || actor == "" || id == "" || event == "" || len(hash) != 64 {
		return empty, false, franchisejourney.ErrReleaseConditioned
	}
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return empty, false, err
	}
	defer tx.Rollback(ctx)
	insert, err := tx.Exec(ctx, `insert into platform.idempotency_record(tenant_id,scope,idempotency_key,request_sha256_hex,status,locked_until,expires_at)
 values($1,'commercial-release',$2,$3,'processing',clock_timestamp()+interval '30 seconds',clock_timestamp()+interval '24 hours') on conflict do nothing`, tenant, c.IdempotencyKey, hash)
	if err != nil {
		return empty, false, err
	}
	if insert.RowsAffected() == 0 {
		var saved string
		if err = tx.QueryRow(ctx, `select request_sha256_hex from platform.idempotency_record where tenant_id=$1 and scope='commercial-release' and idempotency_key=$2`, tenant, c.IdempotencyKey).Scan(&saved); err != nil || saved != hash {
			return empty, false, franchisejourney.ErrConflict
		}
		v, e := readCommercialRelease(ctx, tx, tenant, c.OrganizationID, c.HandoverID, c.IdempotencyKey)
		if e != nil {
			return empty, false, franchisejourney.ErrConflict
		}
		if err = tx.Commit(ctx); err != nil {
			return empty, false, err
		}
		return v, true, nil
	}
	v, err := lockCommercialRelease(ctx, tx, tenant, c.OrganizationID, c.HandoverID, c.ObservationSHA256, p)
	if err != nil {
		return empty, false, err
	}
	var exists bool
	if err = tx.QueryRow(ctx, `select exists(select 1 from sales.commercial_release_receipt where tenant_id=$1 and organization_id=$2 and order_id=$3)`, tenant, c.OrganizationID, v.OrderID).Scan(&exists); err != nil {
		return empty, false, err
	}
	if exists {
		return empty, false, franchisejourney.ErrConflict
	}
	v.ID = id
	v.ReleasedBy = actor
	_, err = tx.Exec(ctx, `insert into sales.commercial_release_receipt(tenant_id,release_id,organization_id,handover_id,order_id,payment_attempt_id,observation_sha256_hex,observation_generation,handover_version,acceptance_sha256_hex,checklist_id,checklist_version,contract_id,contract_sha256_hex,effect,released_by_subject,recorded_at,valid_until,funding_receipt_id)
 values($1,$2,$3,$4,$5,nullif($6,''),$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,nullif($19,''))`, tenant, v.ID, v.OrganizationID, v.HandoverID, v.OrderID, v.PaymentAttemptID, v.ObservationSHA256, v.ObservationGeneration, v.HandoverVersion, v.AcceptanceSHA256, v.ChecklistID, v.ChecklistVersion, v.ContractID, v.ContractSHA256, v.Effect, v.ReleasedBy, v.RecordedAt, v.ValidUntil, v.FundingReceiptID)
	if err != nil {
		return empty, false, err
	}
	if err = writeJourneyOutbox(ctx, tx, tenant, event, "commercial-release", id, 1, "commercial-release.recorded", v); err != nil {
		return empty, false, err
	}
	done, err := tx.Exec(ctx, `update platform.idempotency_record set status='completed',response_code=201,response_body=jsonb_build_object('release_id',$3::text),resource_type='commercial-release',resource_id=$3,locked_until=null where tenant_id=$1 and scope='commercial-release' and idempotency_key=$2 and status='processing'`, tenant, c.IdempotencyKey, id)
	if err != nil {
		return empty, false, err
	}
	if done.RowsAffected() != 1 {
		return empty, false, franchisejourney.ErrConflict
	}
	// An outbox/constraint wait can consume the freshness budget after the
	// initial checks. The checkpoint must still be fresh before committing.
	var commitTime time.Time
	if err = tx.QueryRow(ctx, `select clock_timestamp()`).Scan(&commitTime); err != nil {
		return empty, false, err
	}
	if !v.ValidUntil.After(commitTime) {
		return empty, false, franchisejourney.ErrConflict
	}
	if err = tx.Commit(ctx); err != nil {
		return empty, false, err
	}
	return v, false, nil
}

func (r *FranchiseJourney) ValidateInitialCommercialRelease(ctx context.Context, tenant, organization, handover string, p franchisejourney.HandoverReleaseContract) (franchisejourney.CurrentCommercialRelease, error) {
	var v franchisejourney.CurrentCommercialRelease
	if !p.AllowsCommercialRelease(tenant, organization) {
		return v, franchisejourney.ErrReleaseConditioned
	}
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return v, err
	}
	defer tx.Rollback(ctx)
	v.Receipt, err = scanCommercialRelease(tx.QueryRow(ctx, `select `+commercialReleaseColumns+` from sales.commercial_release_receipt r where r.tenant_id=$1 and r.organization_id=$2 and r.handover_id=$3`, tenant, organization, handover))
	if err != nil {
		return v, err
	}
	current, err := lockCommercialRelease(ctx, tx, tenant, organization, handover, v.Receipt.ObservationSHA256, p)
	if err != nil && !errors.Is(err, franchisejourney.ErrConflict) {
		return v, err
	}
	if e := tx.QueryRow(ctx, `select clock_timestamp()`).Scan(&v.EvaluatedAt); e != nil {
		return v, e
	}
	v.Current = err == nil && current.ObservationGeneration == v.Receipt.ObservationGeneration && current.HandoverVersion == v.Receipt.HandoverVersion && current.AcceptanceSHA256 == v.Receipt.AcceptanceSHA256 && current.ChecklistID == v.Receipt.ChecklistID && current.ChecklistVersion == v.Receipt.ChecklistVersion && current.ContractSHA256 == v.Receipt.ContractSHA256 && v.Receipt.ValidUntil.After(v.EvaluatedAt)
	if err = tx.Commit(ctx); err != nil {
		return v, err
	}
	return v, nil
}
