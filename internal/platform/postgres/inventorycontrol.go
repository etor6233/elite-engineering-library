// Inventory control persistence implements the Microsoft BC-derived portable
// contract recorded in docs/inventory/MICROSOFT_BC_DERIVATION.md.
package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"elite.local/enterprise/internal/inventorycontrol"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type InventoryControl struct{ pool *pgxpool.Pool }

func NewInventoryControl(pool *pgxpool.Pool) *InventoryControl { return &InventoryControl{pool: pool} }

func (r *InventoryControl) AvailableToPromise(ctx context.Context, tenant, organization, variant string, horizon time.Time) (inventorycontrol.ATP, error) {
	value := inventorycontrol.ATP{OrganizationID: organization, VariantID: variant, Horizon: horizon}
	err := r.pool.QueryRow(ctx, `select available_inventory,scheduled_receipt,gross_requirement,available_to_promise from inventory.serial_atp($1,$2,$3,$4)`, tenant, organization, variant, horizon).Scan(&value.AvailableInventory, &value.ScheduledReceipt, &value.GrossRequirement, &value.AvailableToPromise)
	return value, err
}

func (r *InventoryControl) Reserve(ctx context.Context, tenant, eventID, reservationID string, value inventorycontrol.Reservation, stockVersion int64) (inventorycontrol.Reservation, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return value, err
	}
	defer tx.Rollback(ctx)
	result, err := tx.Exec(ctx, `update inventory.stock_unit set state='reserved',version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and organization_id=$2 and stock_unit_id=$3 and variant_id=$4 and state='available' and version=$5`, tenant, value.OrganizationID, value.StockUnitID, value.VariantID, stockVersion)
	if err != nil {
		return value, err
	}
	if result.RowsAffected() != 1 {
		return value, inventorycontrol.ErrConflict
	}
	_, err = tx.Exec(ctx, `insert into inventory.serial_reservation(tenant_id,reservation_id,organization_id,stock_unit_id,variant_id,demand_kind,demand_id,demand_line_id,status,cancellation_disallowed,expires_at,version) values($1,$2,$3,$4,$5,$6,$7,$8,'reservation',$9,$10,1)`, tenant, reservationID, value.OrganizationID, value.StockUnitID, value.VariantID, value.DemandKind, value.DemandID, value.DemandLineID, value.CancellationDisallowed, value.ExpiresAt)
	if err != nil {
		return value, err
	}
	payload, _ := json.Marshal(map[string]string{"reservation_id": reservationID, "stock_unit_id": value.StockUnitID, "demand_kind": value.DemandKind, "demand_id": value.DemandID})
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload) values($1,$2,'serial-reservation',$3,1,'serial-reservation.created',1,clock_timestamp(),$4)`, tenant, eventID, reservationID, payload)
	if err != nil {
		return value, err
	}
	if err = tx.Commit(ctx); err != nil {
		return value, err
	}
	return value, nil
}

func (r *InventoryControl) Release(ctx context.Context, tenant, organization, reservationID string, version int64, eventID string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	var stockID string
	err = tx.QueryRow(ctx, `update inventory.serial_reservation set status='released',version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and organization_id=$2 and reservation_id=$3 and status='reservation' and version=$4 and not cancellation_disallowed returning stock_unit_id`, tenant, organization, reservationID, version).Scan(&stockID)
	if errors.Is(err, pgx.ErrNoRows) {
		return inventorycontrol.ErrConflict
	}
	if err != nil {
		return err
	}
	result, err := tx.Exec(ctx, `update inventory.stock_unit set state='available',version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and organization_id=$2 and stock_unit_id=$3 and state='reserved'`, tenant, organization, stockID)
	if err != nil {
		return err
	}
	if result.RowsAffected() != 1 {
		return inventorycontrol.ErrConflict
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload) values($1,$2,'serial-reservation',$3,$4,'serial-reservation.released',1,clock_timestamp(),jsonb_build_object('stock_unit_id',$5::text))`, tenant, eventID, reservationID, version+1, stockID)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *InventoryControl) CreateTransfer(ctx context.Context, tenant, transferID string, value inventorycontrol.Transfer, stockVersions map[string]int64, eventID string) (inventorycontrol.Transfer, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return value, err
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, `insert into inventory.serial_transfer(tenant_id,transfer_id,from_organization_id,to_organization_id,expected_receipt_at,state,version) values($1,$2,$3,$4,$5,'draft',1)`, tenant, transferID, value.FromOrganizationID, value.ToOrganizationID, value.ExpectedReceiptAt)
	if err != nil {
		return value, err
	}
	for _, stockID := range value.StockUnitIDs {
		var variant string
		err = tx.QueryRow(ctx, `update inventory.stock_unit set state='reserved',version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and organization_id=$2 and stock_unit_id=$3 and state='available' and version=$4 returning variant_id`, tenant, value.FromOrganizationID, stockID, stockVersions[stockID]).Scan(&variant)
		if errors.Is(err, pgx.ErrNoRows) {
			return value, inventorycontrol.ErrConflict
		}
		if err != nil {
			return value, err
		}
		_, err = tx.Exec(ctx, `insert into inventory.serial_transfer_unit(tenant_id,transfer_id,stock_unit_id,variant_id) values($1,$2,$3,$4)`, tenant, transferID, stockID, variant)
		if err != nil {
			return value, err
		}
		_, err = tx.Exec(ctx, `insert into inventory.serial_reservation(tenant_id,reservation_id,organization_id,stock_unit_id,variant_id,demand_kind,demand_id,demand_line_id,status,version) values($1,$2,$3,$4,$5,'transfer-outbound',$6,$4,'reservation',1)`, tenant, transferID+":"+stockID, value.FromOrganizationID, stockID, variant, transferID)
		if err != nil {
			return value, err
		}
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload) values($1,$2,'serial-transfer',$3,1,'serial-transfer.created',1,clock_timestamp(),jsonb_build_object('from_organization_id',$4::text,'to_organization_id',$5::text,'unit_count',$6::int))`, tenant, eventID, transferID, value.FromOrganizationID, value.ToOrganizationID, len(value.StockUnitIDs))
	if err != nil {
		return value, err
	}
	if err = tx.Commit(ctx); err != nil {
		return value, err
	}
	return value, nil
}

func (r *InventoryControl) TransitionTransfer(ctx context.Context, tenant, authorizedOrganization, transferID, current string, version int64, target, eventID string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	var fromOrganization, toOrganization string
	err = tx.QueryRow(ctx, `select from_organization_id,to_organization_id from inventory.serial_transfer where tenant_id=$1 and transfer_id=$2 and state=$3 and version=$4 for update`, tenant, transferID, current, version).Scan(&fromOrganization, &toOrganization)
	if errors.Is(err, pgx.ErrNoRows) {
		return inventorycontrol.ErrConflict
	}
	if err != nil {
		return err
	}
	if (target == "received" && authorizedOrganization != toOrganization) || (target != "received" && authorizedOrganization != fromOrganization) {
		return inventorycontrol.ErrConflict
	}
	rows, err := tx.Query(ctx, `select stock_unit_id from inventory.serial_transfer_unit where tenant_id=$1 and transfer_id=$2 order by stock_unit_id for update`, tenant, transferID)
	if err != nil {
		return err
	}
	var stockIDs []string
	for rows.Next() {
		var stockID string
		if err = rows.Scan(&stockID); err != nil {
			rows.Close()
			return err
		}
		stockIDs = append(stockIDs, stockID)
	}
	rows.Close()
	if err = rows.Err(); err != nil || len(stockIDs) == 0 {
		if err != nil {
			return err
		}
		return inventorycontrol.ErrConflict
	}
	for _, stockID := range stockIDs {
		var result pgconn.CommandTag
		switch target {
		case "in-transit":
			result, err = tx.Exec(ctx, `update inventory.stock_unit set state='in-transit',version=version+1,received_at=null,updated_at=clock_timestamp() where tenant_id=$1 and organization_id=$2 and stock_unit_id=$3 and state='reserved'`, tenant, fromOrganization, stockID)
		case "received":
			result, err = tx.Exec(ctx, `update inventory.stock_unit set organization_id=$3,state='available',version=version+1,received_at=clock_timestamp(),updated_at=clock_timestamp() where tenant_id=$1 and organization_id=$2 and stock_unit_id=$4 and state='in-transit'`, tenant, fromOrganization, toOrganization, stockID)
		case "cancelled":
			result, err = tx.Exec(ctx, `update inventory.stock_unit set state='available',version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and organization_id=$2 and stock_unit_id=$3 and state='reserved'`, tenant, fromOrganization, stockID)
		default:
			continue
		}
		if err != nil {
			return err
		}
		if result.RowsAffected() != 1 {
			return inventorycontrol.ErrConflict
		}
	}
	if target == "in-transit" {
		_, err = tx.Exec(ctx, `update inventory.serial_reservation set status='consumed',version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and demand_kind='transfer-outbound' and demand_id=$2 and status='reservation'`, tenant, transferID)
	} else if target == "cancelled" {
		_, err = tx.Exec(ctx, `update inventory.serial_reservation set status='released',version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and demand_kind='transfer-outbound' and demand_id=$2 and status='reservation'`, tenant, transferID)
	}
	if err != nil {
		return err
	}
	result, err := tx.Exec(ctx, `update inventory.serial_transfer set state=$5,version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and transfer_id=$2 and state=$3 and version=$4`, tenant, transferID, current, version, target)
	if err != nil {
		return err
	}
	if result.RowsAffected() != 1 {
		return inventorycontrol.ErrConflict
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload) values($1,$2,'serial-transfer',$3,$4,$5,1,clock_timestamp(),jsonb_build_object('unit_count',$6::int))`, tenant, eventID, transferID, version+1, "serial-transfer."+target, len(stockIDs))
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}
