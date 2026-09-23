package postgres

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"elite.local/enterprise/internal/inventorycontrol"
	"github.com/jackc/pgx/v5"
)

type customerShipmentAllocation struct {
	reservationID      string
	binID              string
	lotID              string
	quantityBase       string
	reservationVersion int64
	costs              []transferCost
	costAmount         string
}

func readCustomerShipmentReplay(ctx context.Context, tx pgx.Tx, tenant string, command inventorycontrol.CustomerShipmentCommand) (inventorycontrol.CustomerShipment, bool, error) {
	var value inventorycontrol.CustomerShipment
	var postingDate time.Time
	err := tx.QueryRow(ctx, `select s.shipment_id,s.request_id,s.organization_id,s.order_id,s.warehouse_activity_id,s.posting_date,s.order_version,s.fulfillment_state,l.shipment_line_id,l.order_line_id,l.variant_id,l.item_id,l.sales_uom_code,l.quantity::text,l.quantity_base::text,l.cost_amount::text,(select count(*) from inventory.customer_shipment_allocation a where a.tenant_id=s.tenant_id and a.shipment_id=s.shipment_id) from sales.customer_shipment s join sales.customer_shipment_line l using(tenant_id,shipment_id) where s.tenant_id=$1 and s.request_id=$2 for share of s,l`, tenant, command.RequestID).Scan(&value.ID, &value.RequestID, &value.OrganizationID, &value.OrderID, &value.WarehouseActivityID, &postingDate, &value.OrderVersion, &value.FulfillmentState, &value.Line.ID, &value.Line.OrderLineID, &value.Line.VariantID, &value.Line.ItemID, &value.Line.SalesUOMCode, &value.Line.Quantity, &value.Line.QuantityBase, &value.Line.CostAmount, &value.Line.AllocationCount)
	if errors.Is(err, pgx.ErrNoRows) {
		return value, false, nil
	}
	if err != nil {
		return value, false, err
	}
	if value.OrganizationID != command.OrganizationID || value.OrderID != command.OrderID || value.WarehouseActivityID != command.WarehouseActivityID || !postingDate.Equal(command.PostingDate) {
		return inventorycontrol.CustomerShipment{}, false, inventorycontrol.ErrConflict
	}
	value.PostingDate = postingDate
	return value, true, nil
}

func (r *InventoryControl) PostCustomerShipment(ctx context.Context, tenant string, ids inventorycontrol.CustomerShipmentIDs, command inventorycontrol.CustomerShipmentCommand) (inventorycontrol.CustomerShipment, error) {
	value := inventorycontrol.CustomerShipment{ID: ids.ShipmentID, RequestID: command.RequestID, OrganizationID: command.OrganizationID, OrderID: command.OrderID, WarehouseActivityID: command.WarehouseActivityID, PostingDate: command.PostingDate, Line: inventorycontrol.CustomerShipmentLine{ID: ids.LineID}}
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return value, err
	}
	defer tx.Rollback(ctx)
	if _, err = tx.Exec(ctx, `select pg_advisory_xact_lock(hashtextextended($1,0))`, tenant+"\x1fcustomer-shipment-request\x1f"+command.RequestID); err != nil {
		return value, err
	}
	if replay, found, replayErr := readCustomerShipmentReplay(ctx, tx, tenant, command); replayErr != nil {
		return value, replayErr
	} else if found {
		if err = tx.Commit(ctx); err != nil {
			return value, err
		}
		return replay, nil
	}
	var orderLineID, variantID, itemID, salesUOM, salesFactorText, orderedText, shippedBeforeText, costing string
	var orderVersion int64
	err = tx.QueryRow(ctx, `select l.line_id,l.variant_id,b.item_id,b.sales_uom_code,u.qty_per_uom::text,l.quantity::numeric::text,l.shipped_quantity::text,i.costing_method,o.version from sales.customer_order o join inventory.warehouse_activity a on a.tenant_id=o.tenant_id and a.organization_id=o.organization_id and a.source_kind='customer-order' and a.source_id=o.order_id join sales.customer_order_line l on l.tenant_id=o.tenant_id and l.order_id=o.order_id and l.line_id=a.source_line_id join inventory.sales_warehouse_binding b on b.tenant_id=l.tenant_id and b.variant_id=l.variant_id join inventory.item_unit_of_measure u on u.tenant_id=b.tenant_id and u.item_id=b.item_id and u.uom_code=b.sales_uom_code join inventory.stock_item i on i.tenant_id=b.tenant_id and i.item_id=b.item_id where o.tenant_id=$1 and o.organization_id=$2 and o.order_id=$3 and a.activity_id=$4 and a.activity_type='pick' and a.status='registered' and o.state in ('placed','confirmed','paid','allocated') and o.fulfillment_state in ('unfulfilled','partially-shipped') for update of o,l,a`, tenant, command.OrganizationID, command.OrderID, command.WarehouseActivityID).Scan(&orderLineID, &variantID, &itemID, &salesUOM, &salesFactorText, &orderedText, &shippedBeforeText, &costing, &orderVersion)
	if errors.Is(err, pgx.ErrNoRows) || costing != "fifo" {
		return value, inventorycontrol.ErrConflict
	}
	if err != nil {
		return value, err
	}
	if err = warehouseLock(ctx, tx, tenant, command.OrganizationID, itemID); err != nil {
		return value, err
	}
	rows, err := tx.Query(ctx, `select r.reservation_id,r.bin_id,coalesce(r.lot_id,''),r.quantity::text,r.version from inventory.warehouse_activity a join inventory.warehouse_activity_line al using(tenant_id,activity_id) join inventory.bulk_reservation r on r.tenant_id=al.tenant_id and r.reservation_id=al.reservation_id join inventory.warehouse_bin w on w.tenant_id=r.tenant_id and w.organization_id=r.organization_id and w.bin_id=r.bin_id where a.tenant_id=$1 and a.organization_id=$2 and a.activity_id=$3 and a.activity_type='pick' and a.status='registered' and a.source_kind='customer-order' and a.source_id=$4 and a.source_line_id=$5 and r.item_id=$6 and r.demand_kind='customer-order' and r.demand_id=$4 and r.demand_line_id=a.source_line_id||':'||a.activity_id||':'||lpad(al.sequence_no::text,6,'0') and r.status='reservation' and w.bin_type='ship' and not w.movement_blocked order by al.sequence_no for update of r`, tenant, command.OrganizationID, command.WarehouseActivityID, command.OrderID, orderLineID, itemID)
	if err != nil {
		return value, err
	}
	allocations := []customerShipmentAllocation{}
	shipmentBase := new(big.Rat)
	for rows.Next() {
		var allocation customerShipmentAllocation
		if err = rows.Scan(&allocation.reservationID, &allocation.binID, &allocation.lotID, &allocation.quantityBase, &allocation.reservationVersion); err != nil {
			rows.Close()
			return value, err
		}
		quantity, ok := parseRat(allocation.quantityBase)
		if !ok {
			rows.Close()
			return value, inventorycontrol.ErrConflict
		}
		shipmentBase.Add(shipmentBase, quantity)
		allocations = append(allocations, allocation)
	}
	rows.Close()
	if err = rows.Err(); err != nil {
		return value, err
	}
	if len(allocations) == 0 || shipmentBase.Sign() <= 0 {
		return value, inventorycontrol.ErrConflict
	}
	salesFactor, factorOK := parseRat(salesFactorText)
	ordered, orderedOK := parseRat(orderedText)
	shippedBefore, shippedOK := parseRat(shippedBeforeText)
	if !factorOK || !orderedOK || !shippedOK || salesFactor.Sign() <= 0 {
		return value, inventorycontrol.ErrConflict
	}
	shipmentSalesRaw := new(big.Rat).Quo(shipmentBase, salesFactor)
	shipmentSalesText := formatRat(shipmentSalesRaw, 6)
	shipmentSales, salesOK := parseRat(shipmentSalesText)
	if !salesOK || new(big.Rat).Mul(shipmentSales, salesFactor).Cmp(shipmentBase) != 0 || new(big.Rat).Add(shippedBefore, shipmentSales).Cmp(ordered) > 0 {
		return value, inventorycontrol.ErrConflict
	}
	totalCost := new(big.Rat)
	for index := range allocations {
		costs, cost, costErr := transferCosts(ctx, tx, tenant, command.OrganizationID, itemID, allocations[index].lotID, allocations[index].quantityBase, "")
		if costErr != nil {
			return value, costErr
		}
		allocations[index].costs = costs
		allocations[index].costAmount = formatRat(cost, 4)
		totalCost.Add(totalCost, cost)
		for _, component := range costs {
			updated, updateErr := tx.Exec(ctx, `update inventory.bulk_cost_layer set remaining_quantity=remaining_quantity-$3::numeric where tenant_id=$1 and layer_id=$2 and remaining_quantity >= $3::numeric`, tenant, component.layerID, component.quantity)
			if updateErr != nil || updated.RowsAffected() != 1 {
				if updateErr != nil {
					return value, updateErr
				}
				return value, inventorycontrol.ErrConflict
			}
		}
	}
	for index, allocation := range allocations {
		outboundID := fmt.Sprintf("%s-entry-%06d", ids.ShipmentID, index+1)
		unitCost := formatRat(new(big.Rat).Quo(mustRat(allocation.costAmount), mustRat(allocation.quantityBase)), 4)
		_, err = tx.Exec(ctx, `insert into inventory.bulk_inventory_entry(tenant_id,entry_id,organization_id,bin_id,item_id,lot_id,entry_type,quantity,unit_cost,cost_amount,posting_date,source_kind,source_id) values($1,$2,$3,$4,$5,nullif($6,''),'issue',-$7::numeric,$8::numeric,-$9::numeric,$10::date,'customer-shipment',$11)`, tenant, outboundID, command.OrganizationID, allocation.binID, itemID, allocation.lotID, allocation.quantityBase, unitCost, allocation.costAmount, command.PostingDate, fmt.Sprintf("%s:%06d", ids.ShipmentID, index+1))
		if err != nil {
			return value, bulkConflict(err)
		}
		for _, component := range allocation.costs {
			if _, err = tx.Exec(ctx, `insert into inventory.bulk_cost_application(tenant_id,outbound_entry_id,inbound_entry_id,quantity,cost_amount) values($1,$2,$3,$4::numeric,$5::numeric)`, tenant, outboundID, component.inboundEntryID, component.quantity, component.costAmount); err != nil {
				return value, bulkConflict(err)
			}
		}
		if err = consumeRegisteredPickPackaging(ctx, tx, tenant, command.OrganizationID, command.WarehouseActivityID, allocation.reservationID); err != nil {
			return value, err
		}
		updated, updateErr := tx.Exec(ctx, `update inventory.bulk_balance set quantity=quantity-$6::numeric,reserved_quantity=reserved_quantity-$6::numeric,version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and organization_id=$2 and bin_id=$3 and item_id=$4 and lot_id is not distinct from nullif($5,'') and quantity >= $6::numeric and reserved_quantity >= $6::numeric`, tenant, command.OrganizationID, allocation.binID, itemID, allocation.lotID, allocation.quantityBase)
		if updateErr != nil || updated.RowsAffected() != 1 {
			if updateErr != nil {
				return value, updateErr
			}
			return value, inventorycontrol.ErrConflict
		}
		updated, updateErr = tx.Exec(ctx, `update inventory.bulk_reservation set status='consumed',version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and reservation_id=$2 and status='reservation' and version=$3`, tenant, allocation.reservationID, allocation.reservationVersion)
		if updateErr != nil || updated.RowsAffected() != 1 {
			if updateErr != nil {
				return value, updateErr
			}
			return value, inventorycontrol.ErrConflict
		}
	}
	shippedAfter := new(big.Rat).Add(shippedBefore, shipmentSales)
	updated, err := tx.Exec(ctx, `update sales.customer_order_line set shipped_quantity=shipped_quantity+$4::numeric where tenant_id=$1 and order_id=$2 and line_id=$3 and shipped_quantity+$4::numeric<=quantity`, tenant, command.OrderID, orderLineID, shipmentSalesText)
	if err != nil || updated.RowsAffected() != 1 {
		if err != nil {
			return value, err
		}
		return value, inventorycontrol.ErrConflict
	}
	fulfillmentState := "partially-shipped"
	var incomplete bool
	if err = tx.QueryRow(ctx, `select exists(select 1 from sales.customer_order_line where tenant_id=$1 and order_id=$2 and shipped_quantity<quantity)`, tenant, command.OrderID).Scan(&incomplete); err != nil {
		return value, err
	}
	if !incomplete && shippedAfter.Cmp(ordered) == 0 {
		fulfillmentState = "shipped"
	}
	updated, err = tx.Exec(ctx, `update sales.customer_order set fulfillment_state=$4,version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and order_id=$2 and version=$3 and fulfillment_state in ('unfulfilled','partially-shipped')`, tenant, command.OrderID, orderVersion, fulfillmentState)
	if err != nil || updated.RowsAffected() != 1 {
		if err != nil {
			return value, err
		}
		return value, inventorycontrol.ErrConflict
	}
	value.OrderVersion = orderVersion + 1
	value.FulfillmentState = fulfillmentState
	value.Line = inventorycontrol.CustomerShipmentLine{ID: ids.LineID, OrderLineID: orderLineID, VariantID: variantID, ItemID: itemID, SalesUOMCode: salesUOM, Quantity: shipmentSalesText, QuantityBase: formatRat(shipmentBase, 6), CostAmount: formatRat(totalCost, 4), AllocationCount: len(allocations)}
	_, err = tx.Exec(ctx, `insert into sales.customer_shipment(tenant_id,shipment_id,request_id,organization_id,order_id,warehouse_activity_id,posting_date,total_cost_amount,order_version,fulfillment_state) values($1,$2,$3,$4,$5,$6,$7::date,$8::numeric,$9,$10)`, tenant, value.ID, value.RequestID, value.OrganizationID, value.OrderID, value.WarehouseActivityID, value.PostingDate, value.Line.CostAmount, value.OrderVersion, value.FulfillmentState)
	if err != nil {
		return value, bulkConflict(err)
	}
	_, err = tx.Exec(ctx, `insert into sales.customer_shipment_line(tenant_id,shipment_id,shipment_line_id,order_id,order_line_id,variant_id,item_id,sales_uom_code,quantity,quantity_base,cost_amount) values($1,$2,$3,$4,$5,$6,$7,$8,$9::numeric,$10::numeric,$11::numeric)`, tenant, value.ID, value.Line.ID, value.OrderID, value.Line.OrderLineID, value.Line.VariantID, value.Line.ItemID, value.Line.SalesUOMCode, value.Line.Quantity, value.Line.QuantityBase, value.Line.CostAmount)
	if err != nil {
		return value, bulkConflict(err)
	}
	for index, allocation := range allocations {
		outboundID := fmt.Sprintf("%s-entry-%06d", ids.ShipmentID, index+1)
		allocationID := fmt.Sprintf("%s-allocation-%06d", ids.ShipmentID, index+1)
		_, err = tx.Exec(ctx, `insert into inventory.customer_shipment_allocation(tenant_id,shipment_id,shipment_line_id,allocation_id,reservation_id,outbound_entry_id,organization_id,from_bin_id,item_id,lot_id,quantity_base,cost_amount) values($1,$2,$3,$4,$5,$6,$7,$8,$9,nullif($10,''),$11::numeric,$12::numeric)`, tenant, value.ID, value.Line.ID, allocationID, allocation.reservationID, outboundID, value.OrganizationID, allocation.binID, itemID, allocation.lotID, allocation.quantityBase, allocation.costAmount)
		if err != nil {
			return value, bulkConflict(err)
		}
	}
	if err = recordBulkEvent(ctx, tx, tenant, ids.EventID, "customer-shipment", value.ID, "customer-shipment.posted", value.OrderVersion, map[string]any{"request_id": value.RequestID, "organization_id": value.OrganizationID, "order_id": value.OrderID, "order_line_id": value.Line.OrderLineID, "warehouse_activity_id": value.WarehouseActivityID, "quantity": value.Line.Quantity, "quantity_base": value.Line.QuantityBase, "sales_uom_code": value.Line.SalesUOMCode, "cost_amount": value.Line.CostAmount, "fulfillment_state": value.FulfillmentState, "allocations": value.Line.AllocationCount}); err != nil {
		return value, err
	}
	return value, tx.Commit(ctx)
}
