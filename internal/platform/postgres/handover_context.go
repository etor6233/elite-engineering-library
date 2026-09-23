package postgres

// AUTHORED scoped projection over admitted owners; not a second eligibility algorithm.
import (
	"context"
	"elite.local/enterprise/internal/franchisejourney"
	"errors"
	"github.com/jackc/pgx/v5"
)

func (r *FranchiseJourney) InitialHandoverOperatorContext(ctx context.Context, tenant, org, order string, p franchisejourney.HandoverReleaseContract) (franchisejourney.HandoverOperatorContext, error) {
	var v franchisejourney.HandoverOperatorContext
	if !p.AllowsScope(tenant, org) {
		return v, franchisejourney.ErrReleaseConditioned
	}
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return v, err
	}
	defer tx.Rollback(ctx)
	c := franchisejourney.PrepareHandoverCommand{OrganizationID: org, OrderID: order}
	err = tx.QueryRow(ctx, `select l.line_id,coalesce(a.payment_attempt_id,''),case when a.payment_attempt_id is null then coalesce(f.funding_id,'') else '' end,coalesce(o.evidence_sha256_hex,f.receipt_sha256)
 from sales.customer_order s join sales.customer_order_line l on l.tenant_id=s.tenant_id and l.order_id=s.order_id
 left join payment.payment_attempt a on a.tenant_id=s.tenant_id and a.order_id=s.order_id and a.state='captured'
 left join payment.provider_observation o on o.tenant_id=a.tenant_id and o.payment_attempt_id=a.payment_attempt_id
 left join lateral(select funding_id,receipt_sha256 from payment.local_funding_evidence where tenant_id=s.tenant_id and order_id=s.order_id and organization_id=s.organization_id order by created_at desc,funding_id desc limit 1)f on true
 where s.tenant_id=$1 and s.organization_id=$2 and s.order_id=$3 and l.quantity=1 and l.allocated_stock_unit_id is not null and (o.evidence_sha256_hex is not null or (a.payment_attempt_id is null and f.funding_id is not null))`, tenant, org, order).Scan(&c.OrderLineID, &c.PaymentAttemptID, &c.FundingReceiptID, &c.ObservationSHA256)
	if errors.Is(err, pgx.ErrNoRows) {
		return v, franchisejourney.ErrNotFound
	}
	if err != nil {
		return v, err
	}
	if _, err = lockInitialHandoverScope(ctx, tx, tenant, c, p); err != nil {
		return v, err
	}
	var h franchisejourney.Handover
	var policyID, policySHA string
	err = tx.QueryRow(ctx, `select h.handover_id,h.organization_id,h.order_id,h.state,h.version,coalesce(h.checklist_id,''),coalesce(h.checklist_version,0),h.checklist_completed_at,h.customer_accepted_at,p.contract_id,p.contract_sha256_hex
 from sales.delivery_handover_preparation p join sales.delivery_handover h using(tenant_id,handover_id)
 where p.tenant_id=$1 and p.organization_id=$2 and p.order_id=$3 for share of h nowait`, tenant, org, order).Scan(&h.ID, &h.OrganizationID, &h.OrderID, &h.State, &h.Version, &h.ChecklistID, &h.ChecklistVersion, &h.ChecklistCompletedAt, &h.CustomerAcceptedAt, &policyID, &policySHA)
	if err == nil {
		if policyID != p.ID || policySHA != p.DocumentSHA256 {
			return v, franchisejourney.ErrConflict
		}
		h.ChecklistItems = []franchisejourney.ChecklistItem{}
		v.Handover = &h
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return v, err
	}
	v.OrganizationID = org
	v.OrderID = order
	v.OrderLineID = c.OrderLineID
	v.PaymentAttemptID = c.PaymentAttemptID
	v.FundingReceiptID = c.FundingReceiptID
	v.ObservationSHA256 = c.ObservationSHA256
	v.CanPrepare = v.Handover == nil
	v.ReleaseEffect = "READ_ONLY_ELIGIBILITY"
	if p.AllowsCommercialRelease(tenant, org) {
		v.ReleaseEffect = franchisejourney.CommercialReleaseEffect
	}
	if err = tx.QueryRow(ctx, `select clock_timestamp()`).Scan(&v.EvaluatedAt); err != nil {
		return v, err
	}
	if err = tx.Commit(ctx); err != nil {
		return v, err
	}
	return v, nil
}
