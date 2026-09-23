package postgres

// AUTHORED bounded read models of existing query/domain owners. No business writes.
import (
	"context"
	"elite.local/enterprise/internal/platform/identity"
	rm "elite.local/enterprise/internal/rolemetrics"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RoleMetrics struct{ pool *pgxpool.Pool }

func NewRoleMetrics(pool *pgxpool.Pool) (*RoleMetrics, error) {
	if pool == nil {
		return nil, rm.ErrInvalid
	}
	return &RoleMetrics{pool}, nil
}

type metricProjection struct{ source, scope, sql string }

func metricQuery(kind string) (metricProjection, bool) {
	switch kind {
	case "orders":
		return metricProjection{"sales.customer_order", "organization", `select state,currency,count(*)::text,sum(total_minor_units::numeric)::text from sales.customer_order where tenant_id=$1 and organization_id=$2 group by state,currency order by state,currency limit 101`}, true
	case "own-orders":
		return metricProjection{"sales.customer_order", "customer", `select state,currency,count(*)::text,sum(total_minor_units::numeric)::text from sales.customer_order where tenant_id=$1 and organization_id=$2 and customer_principal_id=$3 group by state,currency order by state,currency limit 101`}, true
	case "leads":
		return metricProjection{"crm.lead", "organization", `select lifecycle_state,'',count(*)::text,'' from crm.lead where tenant_id=$1 and organization_id=$2 group by lifecycle_state order by lifecycle_state limit 101`}, true
	case "stock":
		return metricProjection{"inventory.stock_unit", "organization", `select state,'',count(*)::text,'' from inventory.stock_unit where tenant_id=$1 and organization_id=$2 group by state order by state limit 101`}, true
	case "cases":
		return metricProjection{"service_ops.service_case", "organization", `select state,'',count(*)::text,'' from service_ops.service_case where tenant_id=$1 and organization_id=$2 group by state order by state limit 101`}, true
	case "own-cases":
		return metricProjection{"service_ops.service_case + service_ops.warranty", "customer", `select c.state,'',count(*)::text,'' from service_ops.service_case c where c.tenant_id=$1 and c.organization_id=$2 and exists(select 1 from service_ops.warranty w where w.tenant_id=c.tenant_id and w.stock_unit_id=c.stock_unit_id and w.customer_principal_id=$3)group by c.state order by c.state limit 101`}, true
	case "shipments":
		return metricProjection{"logistics.shipment", "origin_or_destination", `select state,'',count(*)::text,'' from logistics.shipment where tenant_id=$1 and(origin_organization_id=$2 or destination_organization_id=$2)group by state order by state limit 101`}, true
	case "appointments":
		return metricProjection{"crm.appointment", "organization", `select state,'',count(*)::text,'' from crm.appointment where tenant_id=$1 and organization_id=$2 group by state order by state limit 101`}, true
	case "own-appointments":
		return metricProjection{"crm.appointment", "customer", `select state,'',count(*)::text,'' from crm.appointment where tenant_id=$1 and organization_id=$2 and customer_principal_id=$3 group by state order by state limit 101`}, true
	case "factory-destination":
		return metricProjection{"factory.production_unit + procurement.purchase_order", "destination", `select u.state,'',count(*)::text,'' from factory.production_unit u join procurement.purchase_order p on p.tenant_id=u.tenant_id and p.purchase_order_id=u.purchase_order_id where u.tenant_id=$1 and p.destination_organization_id=$2 group by u.state order by u.state limit 101`}, true
	case "factory-owned":
		return metricProjection{"factory.production_unit + procurement.serial_supply_plan", "factory", `select u.state,'',count(*)::text,'' from factory.production_unit u join procurement.serial_supply_plan p on p.tenant_id=u.tenant_id and p.purchase_order_id=u.purchase_order_id where u.tenant_id=$1 and p.factory_organization_id=$2 group by u.state order by u.state limit 101`}, true
	case "supply":
		return metricProjection{"procurement.purchase_order + procurement.serial_supply_plan", "destination", `select o.state,'',count(*)::text,'' from procurement.purchase_order o join procurement.serial_supply_plan p on p.tenant_id=o.tenant_id and p.purchase_order_id=o.purchase_order_id where o.tenant_id=$1 and p.destination_organization_id=$2 group by o.state order by o.state limit 101`}, true
	}
	return metricProjection{}, false
}
func (s *RoleMetrics) Read(ctx context.Context, p identity.Principal, kind, org string) (rm.Snapshot, error) {
	var empty rm.Snapshot
	if !rm.Authorized(p, kind, org) {
		return empty, rm.ErrUnavailable
	}
	spec, ok := metricQuery(kind)
	if !ok {
		return empty, rm.ErrInvalid
	}
	tx, e := s.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead, AccessMode: pgx.ReadOnly})
	if e != nil {
		return empty, e
	}
	defer tx.Rollback(ctx)
	if _, e = tx.Exec(ctx, "set local statement_timeout='3s'"); e != nil {
		return empty, e
	}
	result := rm.Snapshot{Kind: kind, OrganizationID: org, Scope: spec.scope, Source: spec.source, Basis: "current_registered_records_by_state", Rows: []rm.Row{}}
	if e = tx.QueryRow(ctx, "select clock_timestamp()").Scan(&result.ObservedAt); e != nil {
		return empty, e
	}
	args := []any{p.TenantID, org}
	if spec.scope == "customer" {
		args = append(args, p.Subject)
	}
	rows, e := tx.Query(ctx, spec.sql, args...)
	if e != nil {
		return empty, e
	}
	for rows.Next() {
		var row rm.Row
		if e = rows.Scan(&row.State, &row.Currency, &row.Count, &row.TotalMinor); e != nil {
			rows.Close()
			return empty, e
		}
		result.Rows = append(result.Rows, row)
	}
	e = rows.Err()
	rows.Close()
	if e != nil {
		return empty, e
	}
	if len(result.Rows) > 100 {
		return empty, rm.ErrCardinality
	}
	return result, tx.Commit(ctx)
}
