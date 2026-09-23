package postgres

import (
	"context"

	"elite.local/enterprise/internal/enterprisequery"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EnterpriseQuery struct{ pool *pgxpool.Pool }

func NewEnterpriseQuery(pool *pgxpool.Pool) *EnterpriseQuery { return &EnterpriseQuery{pool: pool} }

func (r *EnterpriseQuery) Overview(ctx context.Context, tenant, organization string) (enterprisequery.Overview, error) {
	value := enterprisequery.Overview{OrganizationID: organization}
	err := r.pool.QueryRow(ctx, `select
  (select count(*) from sales.customer_order where tenant_id=$1 and organization_id=$2),
  (select count(*) from crm.lead where tenant_id=$1 and organization_id=$2 and lifecycle_state not in ('converted','lost')),
  (select count(*) from inventory.stock_unit where tenant_id=$1 and organization_id=$2 and state='available'),
  (select count(*) from service_ops.service_case where tenant_id=$1 and organization_id=$2 and state not in ('closed','cancelled')),
  (select count(*) from logistics.shipment where tenant_id=$1 and (origin_organization_id=$2 or destination_organization_id=$2) and state not in ('delivered','cancelled'))`, tenant, organization).Scan(&value.Orders, &value.OpenLeads, &value.StockAvailable, &value.OpenCases, &value.ActiveShipments)
	return value, err
}

func (r *EnterpriseQuery) Orders(ctx context.Context, tenant, organization, customer string, limit int, after string) (enterprisequery.Page[enterprisequery.Order], error) {
	rows, err := r.pool.Query(ctx, `select order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version
from sales.customer_order where tenant_id=$1 and organization_id=$2 and ($3='' or customer_principal_id=$3) and ($4='' or order_id>$4)
order by order_id limit $5`, tenant, organization, customer, after, limit+1)
	if err != nil {
		return enterprisequery.Page[enterprisequery.Order]{}, err
	}
	defer rows.Close()
	items := make([]enterprisequery.Order, 0, limit)
	for rows.Next() {
		var value enterprisequery.Order
		if err = rows.Scan(&value.ID, &value.OrganizationID, &value.CustomerSubject, &value.State, &value.Currency, &value.TotalMinorUnits, &value.Version); err != nil {
			return enterprisequery.Page[enterprisequery.Order]{}, err
		}
		items = append(items, value)
	}
	if err = rows.Err(); err != nil {
		return enterprisequery.Page[enterprisequery.Order]{}, err
	}
	return orderPage(items, limit), nil
}

func orderPage(items []enterprisequery.Order, limit int) enterprisequery.Page[enterprisequery.Order] {
	page := enterprisequery.Page[enterprisequery.Order]{Items: items}
	if len(items) > limit {
		page.NextCursor = items[limit-1].ID
		page.Items = items[:limit]
	}
	return page
}

func (r *EnterpriseQuery) FactoryUnits(ctx context.Context, tenant, organization string, limit int, after string) (enterprisequery.Page[enterprisequery.FactoryUnit], error) {
	rows, err := r.pool.Query(ctx, `select u.production_unit_id,p.destination_organization_id,u.purchase_order_id,u.variant_id,u.serial_number,u.state
from factory.production_unit u join procurement.purchase_order p on p.tenant_id=u.tenant_id and p.purchase_order_id=u.purchase_order_id
where u.tenant_id=$1 and p.destination_organization_id=$2 and ($3='' or u.production_unit_id>$3)
order by u.production_unit_id limit $4`, tenant, organization, after, limit+1)
	if err != nil {
		return enterprisequery.Page[enterprisequery.FactoryUnit]{}, err
	}
	defer rows.Close()
	items := make([]enterprisequery.FactoryUnit, 0, limit)
	for rows.Next() {
		var value enterprisequery.FactoryUnit
		if err = rows.Scan(&value.ID, &value.OrganizationID, &value.PurchaseOrderID, &value.VariantID, &value.SerialNumber, &value.State); err != nil {
			return enterprisequery.Page[enterprisequery.FactoryUnit]{}, err
		}
		items = append(items, value)
	}
	if err = rows.Err(); err != nil {
		return enterprisequery.Page[enterprisequery.FactoryUnit]{}, err
	}
	page := enterprisequery.Page[enterprisequery.FactoryUnit]{Items: items}
	if len(items) > limit {
		page.NextCursor = items[limit-1].ID
		page.Items = items[:limit]
	}
	return page, nil
}

func (r *EnterpriseQuery) ServiceCases(ctx context.Context, tenant, organization, customer string, limit int, after string) (enterprisequery.Page[enterprisequery.ServiceCase], error) {
	rows, err := r.pool.Query(ctx, `select c.service_case_id,c.organization_id,c.stock_unit_id,c.state,c.severity,c.description,c.version
from service_ops.service_case c where c.tenant_id=$1 and c.organization_id=$2 and ($3='' or exists(select 1 from service_ops.warranty w where w.tenant_id=c.tenant_id and w.stock_unit_id=c.stock_unit_id and w.customer_principal_id=$3)) and ($4='' or c.service_case_id>$4)
order by c.service_case_id limit $5`, tenant, organization, customer, after, limit+1)
	if err != nil {
		return enterprisequery.Page[enterprisequery.ServiceCase]{}, err
	}
	defer rows.Close()
	items := make([]enterprisequery.ServiceCase, 0, limit)
	for rows.Next() {
		var value enterprisequery.ServiceCase
		if err = rows.Scan(&value.ID, &value.OrganizationID, &value.StockUnitID, &value.State, &value.Severity, &value.Description, &value.Version); err != nil {
			return enterprisequery.Page[enterprisequery.ServiceCase]{}, err
		}
		items = append(items, value)
	}
	if err = rows.Err(); err != nil {
		return enterprisequery.Page[enterprisequery.ServiceCase]{}, err
	}
	page := enterprisequery.Page[enterprisequery.ServiceCase]{Items: items}
	if len(items) > limit {
		page.NextCursor = items[limit-1].ID
		page.Items = items[:limit]
	}
	return page, nil
}
