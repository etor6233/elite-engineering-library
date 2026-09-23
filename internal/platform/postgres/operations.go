package postgres

import (
	"context"
	"elite.local/enterprise/internal/operations"
	"encoding/json"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Operations struct{ pool *pgxpool.Pool }

func NewOperations(pool *pgxpool.Pool) *Operations { return &Operations{pool: pool} }
func outboxPayload(value any) json.RawMessage      { payload, _ := json.Marshal(value); return payload }
func (r *Operations) CreatePurchaseOrder(ctx context.Context, tenant, eventID string, value operations.PurchaseOrder) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err := createPurchaseOrderInTx(ctx, tx, tenant, eventID, value); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// AUTHORED extraction: existing SQL/order retained for a shared transaction.
func createPurchaseOrderInTx(ctx context.Context, tx pgx.Tx, tenant, eventID string, value operations.PurchaseOrder) error {
	var err error
	_, err = tx.Exec(ctx, `insert into procurement.purchase_order(tenant_id,purchase_order_id,supplier_id,destination_organization_id,state,currency,total_minor_units,version)values($1,$2,$3,$4,$5,$6,$7,$8)`, tenant, value.ID, value.SupplierID, value.DestinationOrganizationID, value.State, value.Currency, value.TotalMinorUnits, value.Version)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)values($1,$2,'purchase-order',$3,$4,'purchase-order.created',1,clock_timestamp(),$5)`, tenant, eventID, value.ID, value.Version, outboxPayload(value))
	if err != nil {
		return err
	}
	return nil
}
func (r *Operations) TransitionPurchaseOrder(ctx context.Context, tenant, organization, id, current string, version int64, target, eventID string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err := transitionPurchaseOrderInTx(ctx, tx, tenant, organization, id, current, version, target, eventID); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// AUTHORED extraction: existing SQL/order retained for a shared transaction.
func transitionPurchaseOrderInTx(ctx context.Context, tx pgx.Tx, tenant, organization, id, current string, version int64, target, eventID string) error {
	result, err := tx.Exec(ctx, `update procurement.purchase_order set state=$6,version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and destination_organization_id=$2 and purchase_order_id=$3 and state=$4 and version=$5`, tenant, organization, id, current, version, target)
	if err != nil {
		return err
	}
	if result.RowsAffected() != 1 {
		return operations.ErrConflict
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)values($1,$2,'purchase-order',$3,$4,$5,1,clock_timestamp(),$6)`, tenant, eventID, id, version+1, "purchase-order."+target, outboxPayload(map[string]string{"state": target}))
	if err != nil {
		return err
	}
	return nil
}
func (r *Operations) CreateProductionUnit(ctx context.Context, tenant, eventID string, value operations.ProductionUnit) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err := createProductionUnitInTx(ctx, tx, tenant, eventID, value); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// AUTHORED extraction: existing SQL/order retained for a shared transaction.
func createProductionUnitInTx(ctx context.Context, tx pgx.Tx, tenant, eventID string, value operations.ProductionUnit) error {
	result, err := tx.Exec(ctx, `insert into factory.production_unit(tenant_id,production_unit_id,purchase_order_id,variant_id,serial_number,vin,battery_serial_number,state) select $1,$2,p.purchase_order_id,$4,$5,nullif($6,''),nullif($7,''),$8 from procurement.purchase_order p where p.tenant_id=$1 and p.purchase_order_id=$3 and p.destination_organization_id=$9`, tenant, value.ID, value.PurchaseOrderID, value.VariantID, value.SerialNumber, value.VIN, value.BatterySerialNumber, value.State, value.OrganizationID)
	if err != nil {
		return err
	}
	if result.RowsAffected() != 1 {
		return operations.ErrConflict
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)values($1,$2,'production-unit',$3,1,'production-unit.registered',1,clock_timestamp(),$4)`, tenant, eventID, value.ID, outboxPayload(value))
	if err != nil {
		return err
	}
	return nil
}
func (r *Operations) TransitionProductionUnit(ctx context.Context, tenant, organization, id, current, target, eventID string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err := transitionProductionUnitInTx(ctx, tx, tenant, organization, id, current, target, eventID); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// AUTHORED extraction: existing SQL/order retained for a shared transaction.
func transitionProductionUnitInTx(ctx context.Context, tx pgx.Tx, tenant, organization, id, current, target, eventID string) error {
	result, err := tx.Exec(ctx, `update factory.production_unit u set state=$5,updated_at=clock_timestamp() from procurement.purchase_order p where u.tenant_id=$1 and u.production_unit_id=$3 and u.state=$4 and p.tenant_id=u.tenant_id and p.purchase_order_id=u.purchase_order_id and p.destination_organization_id=$2`, tenant, organization, id, current, target)
	if err != nil {
		return err
	}
	if result.RowsAffected() != 1 {
		return operations.ErrConflict
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)values($1,$2,'production-unit',$3,extract(epoch from clock_timestamp())::bigint,$4,1,clock_timestamp(),$5)`, tenant, eventID, id, "production-unit."+target, outboxPayload(map[string]string{"state": target}))
	if err != nil {
		return err
	}
	return nil
}
func (r *Operations) CreateStockUnit(ctx context.Context, tenant, eventID string, value operations.StockUnit) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err := createStockUnitInTx(ctx, tx, tenant, eventID, value); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// AUTHORED extraction: existing SQL/order retained for a shared transaction.
func createStockUnitInTx(ctx context.Context, tx pgx.Tx, tenant, eventID string, value operations.StockUnit) error {
	var err error
	_, err = tx.Exec(ctx, `insert into inventory.stock_unit(tenant_id,stock_unit_id,organization_id,variant_id,production_unit_id,serial_number,vin,battery_serial_number,state,version)values($1,$2,$3,$4,nullif($5,''),$6,nullif($7,''),nullif($8,''),$9,$10)`, tenant, value.ID, value.OrganizationID, value.VariantID, value.ProductionUnitID, value.SerialNumber, value.VIN, value.BatterySerialNumber, value.State, value.Version)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)values($1,$2,'stock-unit',$3,$4,'stock-unit.received-in-transit',1,clock_timestamp(),$5)`, tenant, eventID, value.ID, value.Version, outboxPayload(value))
	if err != nil {
		return err
	}
	return nil
}
func (r *Operations) TransitionStockUnit(ctx context.Context, tenant, organization, id, current string, version int64, target, eventID string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err := transitionStockUnitInTx(ctx, tx, tenant, organization, id, current, version, target, eventID); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// AUTHORED extraction: existing SQL/order retained for a shared transaction.
func transitionStockUnitInTx(ctx context.Context, tx pgx.Tx, tenant, organization, id, current string, version int64, target, eventID string) error {
	result, err := tx.Exec(ctx, `update inventory.stock_unit set state=$6,version=version+1,received_at=case when $6 in ('available','quarantine') then coalesce(received_at,clock_timestamp()) else received_at end,updated_at=clock_timestamp() where tenant_id=$1 and organization_id=$2 and stock_unit_id=$3 and state=$4 and version=$5`, tenant, organization, id, current, version, target)
	if err != nil {
		return err
	}
	if result.RowsAffected() != 1 {
		return operations.ErrConflict
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)values($1,$2,'stock-unit',$3,$4,$5,1,clock_timestamp(),$6)`, tenant, eventID, id, version+1, "stock-unit."+target, outboxPayload(map[string]string{"state": target}))
	if err != nil {
		return err
	}
	return nil
}
