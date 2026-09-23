// AUTHORED read-only projection: same database, existing owner receipts, no effects.
package postgres

import (
	"context"
	"elite.local/enterprise/internal/franchisejourney"
	"errors"
	"github.com/jackc/pgx/v5"
)

func (r *FranchiseJourney) ReturnOutcome(ctx context.Context, tenant, organization, authorization string) (franchisejourney.ReturnOutcome, error) {
	var v franchisejourney.ReturnOutcome
	tx, e := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead, AccessMode: pgx.ReadOnly})
	if e != nil {
		return v, e
	}
	defer tx.Rollback(ctx)
	e = tx.QueryRow(ctx, `select a.authorization_id,a.organization_id,a.order_id,coalesce(d.disposition_id,''),coalesce(d.customer_remedy,''),transaction_timestamp()
 from sales.return_authorization a
 left join sales.return_receipt r on r.tenant_id=a.tenant_id and r.authorization_id=a.authorization_id
 left join sales.return_disposition d on d.tenant_id=r.tenant_id and d.receipt_id=r.receipt_id
 where a.tenant_id=$1 and a.organization_id=$2 and a.authorization_id=$3 and `+returnScopePredicate+`
 and (d.disposition_id is null or d.customer_remedy=case a.disposition when 'return' then 'refund' when 'exchange' then 'exchange' else '' end)
 and (r.receipt_id is null or (r.organization_id,r.order_id,r.stock_unit_id,r.customer_principal_id)=(a.organization_id,a.order_id,a.stock_unit_id,a.customer_principal_id))`, tenant, organization, authorization).Scan(&v.AuthorizationID, &v.OrganizationID, &v.OrderID, &v.DispositionID, &v.Remedy, &v.ObservedAt)
	if errors.Is(e, pgx.ErrNoRows) {
		return franchisejourney.ReturnOutcome{}, franchisejourney.ErrNotFound
	}
	if e != nil {
		return franchisejourney.ReturnOutcome{}, e
	}
	v.Stages = []franchisejourney.ReturnOutcomeStage{}
	if v.DispositionID != "" {
		rows, e := tx.Query(ctx, `select e.request_id,e.effect_kind,coalesce(x.status,'missing'),coalesce(x.attempt_count,0),coalesce(x.updated_at,e.requested_at),coalesce(x.last_error_code,''),coalesce(x.result_sha256_hex,''),
 case
 when e.effect_kind='refund' and rf.request_id is not null then jsonb_build_object('state',rf.state,'reference',coalesce(rf.provider_refund_reference,''),'amount_minor_units',rf.amount_minor_units::text,'currency',rf.currency)
 when e.effect_kind='exchange' and ex.request_id is not null then jsonb_build_object('state',ex.status,'reference',ex.replacement_handover_id,'amount_minor_units',ex.original_unit_price_minor_units::text,'currency',ex.currency,'handover_state',h.state)
 when e.effect_kind='accounting' and ac.request_id is not null then jsonb_build_object('state',ac.status,'reference',ac.reversal_journal_id,'amount_minor_units',ac.reversed_minor_units::text,'currency',ac.currency)
 when e.effect_kind='fiscal' and fc.request_id is not null then jsonb_build_object('state',fi.status,'reference',fi.invoice_id,'amount_minor_units',fc.total_minor_units::text,'currency',fc.currency)
 else null end
 from sales.return_effect_request e
 left join sales.return_effect_execution x on x.tenant_id=e.tenant_id and x.request_id=e.request_id
 left join payment.return_refund rf on rf.tenant_id=e.tenant_id and rf.request_id=e.request_id and rf.order_id=$3
 left join sales.return_exchange ex on ex.tenant_id=e.tenant_id and ex.request_id=e.request_id and ex.disposition_id=e.disposition_id and ex.original_order_id=$3
 left join sales.delivery_handover h on h.tenant_id=ex.tenant_id and h.handover_id=ex.replacement_handover_id and h.order_id=ex.replacement_order_id and h.organization_id=$4 and h.stock_unit_id=ex.replacement_stock_unit_id
 left join accounting.return_effect_posting ac on ac.tenant_id=e.tenant_id and ac.request_id=e.request_id and ac.disposition_id=e.disposition_id and ac.source_order_id=$3
 left join fiscal.return_credit_note_link fc on fc.tenant_id=e.tenant_id and fc.request_id=e.request_id and fc.disposition_id=e.disposition_id
 left join fiscal.invoice fi on fi.tenant_id=fc.tenant_id and fi.invoice_id=fc.credit_invoice_id and fi.organization_id=$4 and fi.order_id=$3
 where e.tenant_id=$1 and e.disposition_id=$2
 order by case e.effect_kind when 'inventory' then 1 when 'refund' then 2 when 'exchange' then 2 when 'accounting' then 3 else 4 end,e.request_id`, tenant, v.DispositionID, v.OrderID, organization)
		if e != nil {
			return franchisejourney.ReturnOutcome{}, e
		}
		for rows.Next() {
			var x franchisejourney.ReturnOutcomeStage
			if e = rows.Scan(&x.RequestID, &x.Kind, &x.Status, &x.Attempts, &x.UpdatedAt, &x.ErrorCode, &x.ResultSHA256, &x.Owner); e != nil {
				rows.Close()
				return franchisejourney.ReturnOutcome{}, e
			}
			v.Stages = append(v.Stages, x)
		}
		rows.Close()
		if e = rows.Err(); e != nil {
			return franchisejourney.ReturnOutcome{}, e
		}
	}
	if e = v.Validate(organization, authorization); e != nil {
		return franchisejourney.ReturnOutcome{}, e
	}
	if e = tx.Commit(ctx); e != nil {
		return franchisejourney.ReturnOutcome{}, e
	}
	return v, nil
}
