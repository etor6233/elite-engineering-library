package postgres

// AUTHORED exact read of the existing FactoryUnits projection; no state mutation.
import (
	"context"
	"elite.local/enterprise/internal/enterprisequery"
	"errors"
	"github.com/jackc/pgx/v5"
)

func (r *EnterpriseQuery) FactoryUnit(ctx context.Context, tenant, organization, id string) (enterprisequery.FactoryUnit, error) {
	var v enterprisequery.FactoryUnit
	err := r.pool.QueryRow(ctx, `select u.production_unit_id,p.destination_organization_id,u.purchase_order_id,u.variant_id,u.serial_number,u.state
 from factory.production_unit u join procurement.purchase_order p on p.tenant_id=u.tenant_id and p.purchase_order_id=u.purchase_order_id
 where u.tenant_id=$1 and p.destination_organization_id=$2 and u.production_unit_id=$3`, tenant, organization, id).Scan(&v.ID, &v.OrganizationID, &v.PurchaseOrderID, &v.VariantID, &v.SerialNumber, &v.State)
	if errors.Is(err, pgx.ErrNoRows) {
		return enterprisequery.FactoryUnit{}, enterprisequery.ErrFactoryUnitNotFound
	}
	return v, err
}
