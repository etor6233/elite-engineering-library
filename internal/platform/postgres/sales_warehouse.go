package postgres

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"elite.local/enterprise/internal/inventorycontrol"
	"github.com/jackc/pgx/v5"
)

func (r *InventoryControl) ConfigureSalesWarehouseBinding(ctx context.Context, tenant, eventID string, value inventorycontrol.SalesWarehouseBinding) (inventorycontrol.SalesWarehouseBinding, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return value, err
	}
	defer tx.Rollback(ctx)
	var existing inventorycontrol.SalesWarehouseBinding
	err = tx.QueryRow(ctx, `select request_id,variant_id,item_id,sales_uom_code,version from inventory.sales_warehouse_binding where tenant_id=$1 and request_id=$2 for update`, tenant, value.RequestID).Scan(&existing.RequestID, &existing.VariantID, &existing.ItemID, &existing.SalesUOMCode, &existing.Version)
	if err == nil {
		if existing.VariantID != value.VariantID || existing.ItemID != value.ItemID || existing.SalesUOMCode != value.SalesUOMCode {
			return value, inventorycontrol.ErrConflict
		}
		if err = tx.Commit(ctx); err != nil {
			return value, err
		}
		return existing, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return value, err
	}
	_, err = tx.Exec(ctx, `insert into inventory.sales_warehouse_binding(tenant_id,request_id,variant_id,item_id,sales_uom_code,version) values($1,$2,$3,$4,$5,1)`, tenant, value.RequestID, value.VariantID, value.ItemID, value.SalesUOMCode)
	if err != nil {
		return value, bulkConflict(err)
	}
	if err = recordBulkEvent(ctx, tx, tenant, eventID, "sales-warehouse-binding", value.VariantID, "sales-warehouse-binding.configured", 1, map[string]any{"item_id": value.ItemID, "sales_uom_code": value.SalesUOMCode, "request_id": value.RequestID}); err != nil {
		return value, err
	}
	value.Version = 1
	return value, tx.Commit(ctx)
}

func readWarehousePickReplay(ctx context.Context, tx pgx.Tx, tenant string, value inventorycontrol.WarehousePickCommand) (inventorycontrol.WarehouseActivity, bool, error) {
	var activityID, organization, shipBin, itemID, demandKind, demandID, demandLineID, quantity, requestedUOM, assignedTo string
	var allowBreakbulk, useFEFO, allowDedicated bool
	err := tx.QueryRow(ctx, `select activity_id,organization_id,ship_bin_id,item_id,demand_kind,demand_id,demand_line_id,quantity::text,requested_handling_uom,allow_breakbulk,use_fefo,allow_dedicated,assigned_to from inventory.warehouse_pick_request where tenant_id=$1 and request_id=$2 for update`, tenant, value.RequestID).Scan(&activityID, &organization, &shipBin, &itemID, &demandKind, &demandID, &demandLineID, &quantity, &requestedUOM, &allowBreakbulk, &useFEFO, &allowDedicated, &assignedTo)
	if errors.Is(err, pgx.ErrNoRows) {
		return inventorycontrol.WarehouseActivity{}, false, nil
	}
	if err != nil {
		return inventorycontrol.WarehouseActivity{}, false, err
	}
	storedQuantity, storedOK := parseRat(quantity)
	requestedQuantity, requestedOK := parseRat(value.Quantity)
	if !storedOK || !requestedOK || storedQuantity.Cmp(requestedQuantity) != 0 || organization != value.OrganizationID || shipBin != value.ShipBinID || itemID != value.ItemID || demandKind != value.DemandKind || demandID != value.DemandID || demandLineID != value.DemandLineID || requestedUOM != value.HandlingUOM || allowBreakbulk != value.AllowBreakbulk || useFEFO != value.UseFEFO || allowDedicated != value.AllowDedicated || assignedTo != value.AssignedTo {
		return inventorycontrol.WarehouseActivity{}, false, fmt.Errorf("warehouse pick divergent request replay: %w", inventorycontrol.ErrConflict)
	}
	activity, err := readWarehouseActivity(ctx, tx, tenant, organization, activityID)
	return activity, true, err
}

func validateCustomerOrderWarehouseDemand(ctx context.Context, tx pgx.Tx, tenant string, value inventorycontrol.WarehousePickCommand) (string, string, string, error) {
	var salesUOM, totalBase string
	err := tx.QueryRow(ctx, `select b.sales_uom_code,(l.quantity::numeric*u.qty_per_uom)::text from sales.customer_order o join sales.customer_order_line l using(tenant_id,order_id) join inventory.sales_warehouse_binding b on b.tenant_id=l.tenant_id and b.variant_id=l.variant_id join inventory.item_unit_of_measure u on u.tenant_id=b.tenant_id and u.item_id=b.item_id and u.uom_code=b.sales_uom_code where o.tenant_id=$1 and o.organization_id=$2 and o.order_id=$3 and l.line_id=$4 and b.item_id=$5 and o.state in ('placed','confirmed','paid','allocated') for share of o,l,b,u`, tenant, value.OrganizationID, value.DemandID, value.DemandLineID, value.ItemID).Scan(&salesUOM, &totalBase)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", "", "", fmt.Errorf("warehouse pick customer-order source contract: %w", inventorycontrol.ErrConflict)
	}
	if err != nil {
		return "", "", "", err
	}
	var handledBase string
	if err = tx.QueryRow(ctx, `select coalesce(sum(al.quantity),0)::text from inventory.warehouse_activity a join inventory.warehouse_activity_line al using(tenant_id,activity_id) where a.tenant_id=$1 and a.organization_id=$2 and a.activity_type='pick' and a.source_kind='customer-order' and a.source_id=$3 and a.source_line_id=$4 and a.status in ('open','registered')`, tenant, value.OrganizationID, value.DemandID, value.DemandLineID).Scan(&handledBase); err != nil {
		return "", "", "", err
	}
	total, totalOK := parseRat(totalBase)
	handled, handledOK := parseRat(handledBase)
	requested, requestedOK := parseRat(value.Quantity)
	if !totalOK || !handledOK || !requestedOK {
		return "", "", "", fmt.Errorf("warehouse pick customer-order quantity parse: %w", inventorycontrol.ErrConflict)
	}
	outstanding := new(big.Rat).Sub(total, handled)
	if outstanding.Sign() <= 0 || requested.Cmp(outstanding) > 0 {
		return "", "", "", fmt.Errorf("warehouse pick customer-order outstanding: %w", inventorycontrol.ErrConflict)
	}
	return salesUOM, formatRat(total, 6), formatRat(outstanding, 6), nil
}
