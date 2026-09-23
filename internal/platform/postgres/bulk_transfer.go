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

type transferCost struct {
	layerID, inboundEntryID, quantity, unitCost, costAmount string
}

func readBulkTransfer(ctx context.Context, tx pgx.Tx, tenant, transfer string) (inventorycontrol.BulkTransfer, error) {
	var value inventorycontrol.BulkTransfer
	err := tx.QueryRow(ctx, `select t.transfer_id,l.line_id,t.request_id,t.from_organization_id,t.to_organization_id,t.receive_bin_id,t.in_transit_code,l.item_id,l.quantity::text,coalesce(l.specific_receipt_entry_id,''),l.shipped_quantity::text,l.received_quantity::text,l.cost_amount::text,coalesce(t.put_away_activity_id,''),t.status,t.version,t.posting_date from inventory.bulk_transfer t join inventory.bulk_transfer_line l using(tenant_id,transfer_id) where t.tenant_id=$1 and t.transfer_id=$2`, tenant, transfer).Scan(&value.ID, &value.LineID, &value.RequestID, &value.FromOrganizationID, &value.ToOrganizationID, &value.ReceiveBinID, &value.InTransitCode, &value.ItemID, &value.Quantity, &value.SpecificReceiptEntry, &value.ShippedQuantity, &value.ReceivedQuantity, &value.CostAmount, &value.PutAwayID, &value.Status, &value.Version, &value.PostingDate)
	if errors.Is(err, pgx.ErrNoRows) {
		return value, inventorycontrol.ErrConflict
	}
	return value, err
}

func (r *InventoryControl) CreateBulkTransfer(ctx context.Context, tenant, transferID, lineID, eventID string, command inventorycontrol.BulkTransferCommand) (inventorycontrol.BulkTransfer, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return inventorycontrol.BulkTransfer{}, err
	}
	defer tx.Rollback(ctx)
	var binType, costing string
	err = tx.QueryRow(ctx, `select w.bin_type,i.costing_method from inventory.warehouse_bin w join inventory.item_bin_policy p using(tenant_id,organization_id,bin_id) join inventory.stock_item i on i.tenant_id=p.tenant_id and i.item_id=p.item_id where w.tenant_id=$1 and w.organization_id=$2 and w.bin_id=$3 and p.item_id=$4 and not w.movement_blocked for share of w,p,i`, tenant, command.ToOrganizationID, command.ReceiveBinID, command.ItemID).Scan(&binType, &costing)
	if errors.Is(err, pgx.ErrNoRows) || binType != "receive" {
		return inventorycontrol.BulkTransfer{}, inventorycontrol.ErrConflict
	}
	if err != nil {
		return inventorycontrol.BulkTransfer{}, err
	}
	if (costing == "specific") != (command.SpecificReceiptEntry != "") || (costing != "fifo" && costing != "specific") {
		return inventorycontrol.BulkTransfer{}, inventorycontrol.ErrConflict
	}
	inserted, err := tx.Exec(ctx, `insert into inventory.bulk_transfer(tenant_id,transfer_id,request_id,from_organization_id,to_organization_id,receive_bin_id,in_transit_code,status,posting_date,source_kind,source_id,version) values($1,$2,$3,$4,$5,$6,$7,'released',$8::date,$9,$10,1) on conflict (tenant_id,request_id) do nothing`, tenant, transferID, command.RequestID, command.FromOrganizationID, command.ToOrganizationID, command.ReceiveBinID, command.InTransitCode, command.PostingDate, command.SourceKind, command.SourceID)
	if err != nil {
		return inventorycontrol.BulkTransfer{}, bulkConflict(err)
	}
	if inserted.RowsAffected() == 0 {
		err = tx.QueryRow(ctx, `select transfer_id from inventory.bulk_transfer where tenant_id=$1 and request_id=$2 and from_organization_id=$3 and to_organization_id=$4 and receive_bin_id=$5 and in_transit_code=$6 and posting_date=$7::date and source_kind=$8 and source_id=$9`, tenant, command.RequestID, command.FromOrganizationID, command.ToOrganizationID, command.ReceiveBinID, command.InTransitCode, command.PostingDate, command.SourceKind, command.SourceID).Scan(&transferID)
		if errors.Is(err, pgx.ErrNoRows) {
			return inventorycontrol.BulkTransfer{}, inventorycontrol.ErrConflict
		}
		if err != nil {
			return inventorycontrol.BulkTransfer{}, err
		}
		value, readErr := readBulkTransfer(ctx, tx, tenant, transferID)
		storedQuantity, storedOK := parseRat(value.Quantity)
		requestedQuantity, requestedOK := parseRat(command.Quantity)
		if readErr != nil || !storedOK || !requestedOK || storedQuantity.Cmp(requestedQuantity) != 0 || value.ItemID != command.ItemID || value.SpecificReceiptEntry != command.SpecificReceiptEntry {
			return inventorycontrol.BulkTransfer{}, inventorycontrol.ErrConflict
		}
		return value, tx.Commit(ctx)
	}
	_, err = tx.Exec(ctx, `insert into inventory.bulk_transfer_line(tenant_id,transfer_id,line_id,item_id,quantity,specific_receipt_entry_id) values($1,$2,$3,$4,$5::numeric,nullif($6,''))`, tenant, transferID, lineID, command.ItemID, command.Quantity, command.SpecificReceiptEntry)
	if err != nil {
		return inventorycontrol.BulkTransfer{}, bulkConflict(err)
	}
	if err = recordBulkEvent(ctx, tx, tenant, eventID, "bulk-transfer", transferID, "bulk-transfer.released", 1, map[string]any{"from_organization_id": command.FromOrganizationID, "to_organization_id": command.ToOrganizationID, "item_id": command.ItemID, "quantity": command.Quantity, "in_transit_code": command.InTransitCode}); err != nil {
		return inventorycontrol.BulkTransfer{}, err
	}
	value, err := readBulkTransfer(ctx, tx, tenant, transferID)
	if err != nil {
		return value, err
	}
	return value, tx.Commit(ctx)
}

func transferCosts(ctx context.Context, tx pgx.Tx, tenant, organization, item, lot, quantity, specificReceiptEntry string) ([]transferCost, *big.Rat, error) {
	wanted, ok := parseRat(quantity)
	if !ok {
		return nil, nil, inventorycontrol.ErrConflict
	}
	query := `select layer_id,receipt_entry_id,remaining_quantity::text,unit_cost::text from inventory.bulk_cost_layer where tenant_id=$1 and organization_id=$2 and item_id=$3 and lot_id is not distinct from nullif($4,'') and remaining_quantity>0`
	args := []any{tenant, organization, item, lot}
	if specificReceiptEntry != "" {
		query += ` and receipt_entry_id=$5`
		args = append(args, specificReceiptEntry)
	}
	query += ` order by posting_date,receipt_entry_id for update`
	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	remaining, total := new(big.Rat).Set(wanted), new(big.Rat)
	values := []transferCost{}
	for rows.Next() && remaining.Sign() > 0 {
		var value transferCost
		var availableText string
		if err = rows.Scan(&value.layerID, &value.inboundEntryID, &availableText, &value.unitCost); err != nil {
			return nil, nil, err
		}
		available, validAvailable := parseRat(availableText)
		unit, validUnit := parseRat(value.unitCost)
		if !validAvailable || !validUnit {
			return nil, nil, inventorycontrol.ErrConflict
		}
		take := minRat(available, remaining)
		value.quantity = formatRat(take, 6)
		value.costAmount = formatRat(new(big.Rat).Mul(take, unit), 4)
		rounded, _ := parseRat(value.costAmount)
		total.Add(total, rounded)
		values = append(values, value)
		remaining.Sub(remaining, take)
	}
	if err = rows.Err(); err != nil {
		return nil, nil, err
	}
	if remaining.Sign() != 0 || len(values) == 0 {
		return nil, nil, inventorycontrol.ErrConflict
	}
	return values, total, nil
}

func (r *InventoryControl) ShipBulkTransfer(ctx context.Context, tenant, transferID, fromOrganization, toOrganization string, version int64, shipmentID, eventID string, command inventorycontrol.BulkTransferPostingCommand) (inventorycontrol.BulkTransfer, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return inventorycontrol.BulkTransfer{}, err
	}
	defer tx.Rollback(ctx)
	var replayID, replayActivity, replayQuantity string
	var replayDate time.Time
	err = tx.QueryRow(ctx, `select s.shipment_id,s.warehouse_activity_id,s.quantity::text,s.posting_date from inventory.bulk_transfer_shipment s join inventory.bulk_transfer t using(tenant_id,transfer_id) where s.tenant_id=$1 and s.transfer_id=$2 and s.request_id=$3 and t.from_organization_id=$4 and t.to_organization_id=$5`, tenant, transferID, command.RequestID, fromOrganization, toOrganization).Scan(&replayID, &replayActivity, &replayQuantity, &replayDate)
	if err == nil {
		stored, storedOK := parseRat(replayQuantity)
		requested, requestedOK := parseRat(command.Quantity)
		if !storedOK || !requestedOK || stored.Cmp(requested) != 0 || replayActivity != command.WarehouseActivityID || !replayDate.Equal(command.PostingDate) {
			return inventorycontrol.BulkTransfer{}, inventorycontrol.ErrConflict
		}
		value, readErr := readBulkTransfer(ctx, tx, tenant, transferID)
		if readErr != nil {
			return value, readErr
		}
		value.PostingID = replayID
		return value, tx.Commit(ctx)
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return inventorycontrol.BulkTransfer{}, err
	}
	var lineID, organization, itemID, totalText, shippedText, receivedText, costing, specificReceiptEntry string
	err = tx.QueryRow(ctx, `select l.line_id,t.from_organization_id,l.item_id,l.quantity::text,l.shipped_quantity::text,l.received_quantity::text,i.costing_method,coalesce(l.specific_receipt_entry_id,'') from inventory.bulk_transfer t join inventory.bulk_transfer_line l using(tenant_id,transfer_id) join inventory.stock_item i on i.tenant_id=t.tenant_id and i.item_id=l.item_id where t.tenant_id=$1 and t.transfer_id=$2 and t.from_organization_id=$3 and t.to_organization_id=$4 and t.status in ('released','partially-shipped','partially-received') and t.version=$5 for update of t,l`, tenant, transferID, fromOrganization, toOrganization, version).Scan(&lineID, &organization, &itemID, &totalText, &shippedText, &receivedText, &costing, &specificReceiptEntry)
	if errors.Is(err, pgx.ErrNoRows) || (costing == "specific") != (specificReceiptEntry != "") || (costing != "fifo" && costing != "specific") {
		return inventorycontrol.BulkTransfer{}, inventorycontrol.ErrConflict
	}
	if err != nil {
		return inventorycontrol.BulkTransfer{}, err
	}
	if err = warehouseLock(ctx, tx, tenant, organization, itemID); err != nil {
		return inventorycontrol.BulkTransfer{}, err
	}
	total, totalOK := parseRat(totalText)
	shippedBefore, shippedOK := parseRat(shippedText)
	receivedBefore, receivedOK := parseRat(receivedText)
	postingQuantity, postingOK := parseRat(command.Quantity)
	if !totalOK || !shippedOK || !receivedOK || !postingOK || postingQuantity.Cmp(new(big.Rat).Sub(total, shippedBefore)) > 0 {
		return inventorycontrol.BulkTransfer{}, inventorycontrol.ErrConflict
	}
	rows, err := tx.Query(ctx, `select r.reservation_id,r.bin_id,coalesce(r.lot_id,''),r.quantity::text,r.version from inventory.warehouse_activity a join inventory.warehouse_activity_line al using(tenant_id,activity_id) join inventory.bulk_reservation r on r.tenant_id=al.tenant_id and r.reservation_id=al.reservation_id join inventory.warehouse_bin w on w.tenant_id=r.tenant_id and w.organization_id=r.organization_id and w.bin_id=r.bin_id where a.tenant_id=$1 and a.organization_id=$2 and a.activity_id=$3 and a.activity_type='pick' and a.status='registered' and a.source_kind='transfer-outbound' and a.source_id=$4 and a.source_line_id=$5 and r.item_id=$6 and r.demand_kind='transfer-outbound' and r.demand_id=$4 and r.demand_line_id=a.source_line_id||':'||a.activity_id||':'||lpad(al.sequence_no::text,6,'0') and r.status='reservation' and (r.expires_at is null or r.expires_at>clock_timestamp()) and w.bin_type='ship' and not w.movement_blocked order by al.sequence_no for update of a,r`, tenant, organization, command.WarehouseActivityID, transferID, lineID, itemID)
	if err != nil {
		return inventorycontrol.BulkTransfer{}, err
	}
	type allocation struct {
		reservation, bin, lot, quantity string
		reservationVersion              int64
	}
	allocations, shipped := []allocation{}, new(big.Rat)
	for rows.Next() {
		var a allocation
		if err = rows.Scan(&a.reservation, &a.bin, &a.lot, &a.quantity, &a.reservationVersion); err != nil {
			rows.Close()
			return inventorycontrol.BulkTransfer{}, err
		}
		q, ok := parseRat(a.quantity)
		if !ok {
			rows.Close()
			return inventorycontrol.BulkTransfer{}, inventorycontrol.ErrConflict
		}
		shipped.Add(shipped, q)
		allocations = append(allocations, a)
	}
	rows.Close()
	if err = rows.Err(); err != nil {
		return inventorycontrol.BulkTransfer{}, err
	}
	if shipped.Cmp(postingQuantity) != 0 || len(allocations) == 0 {
		return inventorycontrol.BulkTransfer{}, inventorycontrol.ErrConflict
	}
	_, err = tx.Exec(ctx, `insert into inventory.bulk_transfer_shipment(tenant_id,shipment_id,transfer_id,request_id,warehouse_activity_id,quantity,cost_amount,posting_date) values($1,$2,$3,$4,$5,$6::numeric,0,$7::date)`, tenant, shipmentID, transferID, command.RequestID, command.WarehouseActivityID, command.Quantity, command.PostingDate)
	if err != nil {
		return inventorycontrol.BulkTransfer{}, bulkConflict(err)
	}
	totalCost := new(big.Rat)
	for index, a := range allocations {
		outboundID := fmt.Sprintf("%s-entry-%06d", shipmentID, index+1)
		allocationID := fmt.Sprintf("%s-allocation-%06d", shipmentID, index+1)
		costs, allocationCost, costErr := transferCosts(ctx, tx, tenant, organization, itemID, a.lot, a.quantity, specificReceiptEntry)
		if costErr != nil {
			return inventorycontrol.BulkTransfer{}, costErr
		}
		totalCost.Add(totalCost, allocationCost)
		unitCost := formatRat(new(big.Rat).Quo(allocationCost, mustRat(a.quantity)), 4)
		_, err = tx.Exec(ctx, `insert into inventory.bulk_inventory_entry(tenant_id,entry_id,organization_id,bin_id,item_id,lot_id,entry_type,quantity,unit_cost,cost_amount,posting_date,source_kind,source_id) values($1,$2,$3,$4,$5,nullif($6,''),'issue',-$7::numeric,$8::numeric,-$9::numeric,$10::date,'bulk-transfer-shipment',$11)`, tenant, outboundID, organization, a.bin, itemID, a.lot, a.quantity, unitCost, formatRat(allocationCost, 4), command.PostingDate, fmt.Sprintf("%s:%06d", shipmentID, index+1))
		if err != nil {
			return inventorycontrol.BulkTransfer{}, bulkConflict(err)
		}
		_, err = tx.Exec(ctx, `insert into inventory.bulk_transfer_allocation(tenant_id,transfer_id,line_id,allocation_id,reservation_id,from_bin_id,item_id,lot_id,quantity,outbound_entry_id,shipment_id) values($1,$2,$3,$4,$5,$6,$7,nullif($8,''),$9::numeric,$10,$11)`, tenant, transferID, lineID, allocationID, a.reservation, a.bin, itemID, a.lot, a.quantity, outboundID, shipmentID)
		if err != nil {
			return inventorycontrol.BulkTransfer{}, bulkConflict(err)
		}
		for componentIndex, cost := range costs {
			updated, updateErr := tx.Exec(ctx, `update inventory.bulk_cost_layer set remaining_quantity=remaining_quantity-$3::numeric where tenant_id=$1 and layer_id=$2 and remaining_quantity >= $3::numeric`, tenant, cost.layerID, cost.quantity)
			if updateErr != nil || updated.RowsAffected() != 1 {
				if updateErr != nil {
					return inventorycontrol.BulkTransfer{}, updateErr
				}
				return inventorycontrol.BulkTransfer{}, inventorycontrol.ErrConflict
			}
			_, err = tx.Exec(ctx, `insert into inventory.bulk_cost_application(tenant_id,outbound_entry_id,inbound_entry_id,quantity,cost_amount) values($1,$2,$3,$4::numeric,$5::numeric)`, tenant, outboundID, cost.inboundEntryID, cost.quantity, cost.costAmount)
			if err != nil {
				return inventorycontrol.BulkTransfer{}, bulkConflict(err)
			}
			componentID := fmt.Sprintf("%s-cost-%06d-%06d", shipmentID, index+1, componentIndex+1)
			_, err = tx.Exec(ctx, `insert into inventory.bulk_transfer_cost_component(tenant_id,transfer_id,allocation_id,component_id,source_receipt_entry_id,quantity,unit_cost,cost_amount) values($1,$2,$3,$4,$5,$6::numeric,$7::numeric,$8::numeric)`, tenant, transferID, allocationID, componentID, cost.inboundEntryID, cost.quantity, cost.unitCost, cost.costAmount)
			if err != nil {
				return inventorycontrol.BulkTransfer{}, bulkConflict(err)
			}
		}
		if err = consumeRegisteredPickPackaging(ctx, tx, tenant, organization, command.WarehouseActivityID, a.reservation); err != nil {
			return inventorycontrol.BulkTransfer{}, err
		}
		updated, updateErr := tx.Exec(ctx, `update inventory.bulk_balance set quantity=quantity-$6::numeric,reserved_quantity=reserved_quantity-$6::numeric,version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and organization_id=$2 and bin_id=$3 and item_id=$4 and lot_id is not distinct from nullif($5,'') and quantity >= $6::numeric and reserved_quantity >= $6::numeric`, tenant, organization, a.bin, itemID, a.lot, a.quantity)
		if updateErr != nil || updated.RowsAffected() != 1 {
			if updateErr != nil {
				return inventorycontrol.BulkTransfer{}, updateErr
			}
			return inventorycontrol.BulkTransfer{}, inventorycontrol.ErrConflict
		}
		updated, err = tx.Exec(ctx, `update inventory.bulk_reservation set status='consumed',version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and reservation_id=$2 and status='reservation' and version=$3`, tenant, a.reservation, a.reservationVersion)
		if err != nil || updated.RowsAffected() != 1 {
			if err != nil {
				return inventorycontrol.BulkTransfer{}, err
			}
			return inventorycontrol.BulkTransfer{}, inventorycontrol.ErrConflict
		}
	}
	shippedAfter := new(big.Rat).Add(shippedBefore, postingQuantity)
	status := "partially-shipped"
	if receivedBefore.Sign() > 0 {
		status = "partially-received"
	} else if shippedAfter.Cmp(total) == 0 {
		status = "shipped"
	}
	updated, err := tx.Exec(ctx, `update inventory.bulk_transfer set status=$4,version=version+1,shipped_at=coalesce(shipped_at,clock_timestamp()),updated_at=clock_timestamp() where tenant_id=$1 and transfer_id=$2 and status in ('released','partially-shipped','partially-received') and version=$3`, tenant, transferID, version, status)
	if err != nil || updated.RowsAffected() != 1 {
		if err != nil {
			return inventorycontrol.BulkTransfer{}, err
		}
		return inventorycontrol.BulkTransfer{}, inventorycontrol.ErrConflict
	}
	_, err = tx.Exec(ctx, `update inventory.bulk_transfer_line set shipped_quantity=shipped_quantity+$3::numeric,cost_amount=cost_amount+$4::numeric where tenant_id=$1 and transfer_id=$2`, tenant, transferID, command.Quantity, formatRat(totalCost, 4))
	if err != nil {
		return inventorycontrol.BulkTransfer{}, err
	}
	_, err = tx.Exec(ctx, `update inventory.bulk_transfer_shipment set cost_amount=$3::numeric where tenant_id=$1 and shipment_id=$2`, tenant, shipmentID, formatRat(totalCost, 4))
	if err != nil {
		return inventorycontrol.BulkTransfer{}, err
	}
	if err = recordBulkEvent(ctx, tx, tenant, eventID, "bulk-transfer", transferID, "bulk-transfer.shipped", version+1, map[string]any{"shipment_id": shipmentID, "request_id": command.RequestID, "organization_id": organization, "item_id": itemID, "quantity": command.Quantity, "cost_amount": formatRat(totalCost, 4), "allocations": len(allocations)}); err != nil {
		return inventorycontrol.BulkTransfer{}, err
	}
	value, err := readBulkTransfer(ctx, tx, tenant, transferID)
	if err != nil {
		return value, err
	}
	value.PostingID = shipmentID
	return value, tx.Commit(ctx)
}

func mustRat(value string) *big.Rat {
	parsed, ok := parseRat(value)
	if !ok {
		panic("validated decimal became invalid")
	}
	return parsed
}

func (r *InventoryControl) ReceiveBulkTransfer(ctx context.Context, tenant, transferID, fromOrganization, toOrganization string, version int64, receiptID, putAwayID, eventID string, command inventorycontrol.BulkTransferPostingCommand) (inventorycontrol.BulkTransfer, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return inventorycontrol.BulkTransfer{}, err
	}
	defer tx.Rollback(ctx)
	var replayID, replayPutAway, replayQuantity string
	var replayDate time.Time
	err = tx.QueryRow(ctx, `select r.receipt_id,r.put_away_activity_id,r.quantity::text,r.posting_date from inventory.bulk_transfer_receipt r join inventory.bulk_transfer t using(tenant_id,transfer_id) where r.tenant_id=$1 and r.transfer_id=$2 and r.request_id=$3 and t.from_organization_id=$4 and t.to_organization_id=$5`, tenant, transferID, command.RequestID, fromOrganization, toOrganization).Scan(&replayID, &replayPutAway, &replayQuantity, &replayDate)
	if err == nil {
		stored, storedOK := parseRat(replayQuantity)
		requested, requestedOK := parseRat(command.Quantity)
		if !storedOK || !requestedOK || stored.Cmp(requested) != 0 || !replayDate.Equal(command.PostingDate) {
			return inventorycontrol.BulkTransfer{}, inventorycontrol.ErrConflict
		}
		value, readErr := readBulkTransfer(ctx, tx, tenant, transferID)
		if readErr != nil {
			return value, readErr
		}
		value.PostingID, value.PutAwayID = replayID, replayPutAway
		return value, tx.Commit(ctx)
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return inventorycontrol.BulkTransfer{}, err
	}
	var lineID, organization, receiveBin, itemID, totalText, shippedText, receivedText string
	err = tx.QueryRow(ctx, `select l.line_id,t.to_organization_id,t.receive_bin_id,l.item_id,l.quantity::text,l.shipped_quantity::text,l.received_quantity::text from inventory.bulk_transfer t join inventory.bulk_transfer_line l using(tenant_id,transfer_id) join inventory.warehouse_bin w on w.tenant_id=t.tenant_id and w.organization_id=t.to_organization_id and w.bin_id=t.receive_bin_id where t.tenant_id=$1 and t.transfer_id=$2 and t.from_organization_id=$3 and t.to_organization_id=$4 and t.status in ('partially-shipped','shipped','partially-received') and t.version=$5 and w.bin_type='receive' and not w.movement_blocked for update of t,l,w`, tenant, transferID, fromOrganization, toOrganization, version).Scan(&lineID, &organization, &receiveBin, &itemID, &totalText, &shippedText, &receivedText)
	if errors.Is(err, pgx.ErrNoRows) {
		return inventorycontrol.BulkTransfer{}, inventorycontrol.ErrConflict
	}
	if err != nil {
		return inventorycontrol.BulkTransfer{}, err
	}
	if err = warehouseLock(ctx, tx, tenant, organization, itemID); err != nil {
		return inventorycontrol.BulkTransfer{}, err
	}
	total, totalOK := parseRat(totalText)
	shipped, shippedOK := parseRat(shippedText)
	receivedBefore, receivedOK := parseRat(receivedText)
	postingQuantity, postingOK := parseRat(command.Quantity)
	if !totalOK || !shippedOK || !receivedOK || !postingOK || postingQuantity.Cmp(new(big.Rat).Sub(shipped, receivedBefore)) > 0 {
		return inventorycontrol.BulkTransfer{}, inventorycontrol.ErrConflict
	}
	rows, err := tx.Query(ctx, `select c.component_id,coalesce(a.lot_id,''),(c.quantity-c.received_quantity)::text,c.unit_cost::text from inventory.bulk_transfer_cost_component c join inventory.bulk_transfer_allocation a using(tenant_id,transfer_id,allocation_id) join inventory.bulk_transfer_shipment s on s.tenant_id=a.tenant_id and s.shipment_id=a.shipment_id where c.tenant_id=$1 and c.transfer_id=$2 and c.received_quantity<c.quantity order by s.created_at,c.component_id for update of c`, tenant, transferID)
	if err != nil {
		return inventorycontrol.BulkTransfer{}, err
	}
	type component struct{ id, lot, quantity, unit, cost string }
	components := []component{}
	remaining := new(big.Rat).Set(postingQuantity)
	totalCost := new(big.Rat)
	for rows.Next() && remaining.Sign() > 0 {
		var c component
		var availableText string
		if err = rows.Scan(&c.id, &c.lot, &availableText, &c.unit); err != nil {
			rows.Close()
			return inventorycontrol.BulkTransfer{}, err
		}
		available, availableOK := parseRat(availableText)
		unit, unitOK := parseRat(c.unit)
		if !availableOK || !unitOK {
			rows.Close()
			return inventorycontrol.BulkTransfer{}, inventorycontrol.ErrConflict
		}
		take := minRat(available, remaining)
		c.quantity = formatRat(take, 6)
		c.cost = formatRat(new(big.Rat).Mul(take, unit), 4)
		rounded, _ := parseRat(c.cost)
		totalCost.Add(totalCost, rounded)
		components = append(components, c)
		remaining.Sub(remaining, take)
	}
	rows.Close()
	if err = rows.Err(); err != nil {
		return inventorycontrol.BulkTransfer{}, err
	}
	if len(components) == 0 || remaining.Sign() != 0 {
		return inventorycontrol.BulkTransfer{}, inventorycontrol.ErrConflict
	}
	_, err = tx.Exec(ctx, `insert into inventory.warehouse_activity(tenant_id,activity_id,organization_id,activity_type,status,source_kind,source_id,source_line_id,request_id,version) values($1,$2,$3,'put-away','open','bulk-transfer-receipt',$4,$5,$6,1)`, tenant, putAwayID, organization, transferID, lineID, command.RequestID)
	if err != nil {
		return inventorycontrol.BulkTransfer{}, bulkConflict(err)
	}
	_, err = tx.Exec(ctx, `insert into inventory.bulk_transfer_receipt(tenant_id,receipt_id,transfer_id,request_id,quantity,cost_amount,posting_date,put_away_activity_id) values($1,$2,$3,$4,$5::numeric,$6::numeric,$7::date,$8)`, tenant, receiptID, transferID, command.RequestID, command.Quantity, formatRat(totalCost, 4), command.PostingDate, putAwayID)
	if err != nil {
		return inventorycontrol.BulkTransfer{}, bulkConflict(err)
	}
	for index, c := range components {
		receiptEntryID := fmt.Sprintf("%s-entry-%06d", receiptID, index+1)
		layerID := fmt.Sprintf("%s-layer-%06d", receiptID, index+1)
		var balanceVersion int64
		err = tx.QueryRow(ctx, `insert into inventory.bulk_balance(tenant_id,organization_id,bin_id,item_id,lot_id,quantity,reserved_quantity,version) values($1,$2,$3,$4,nullif($5,''),$6::numeric,0,1) on conflict (tenant_id,organization_id,bin_id,item_id,lot_id) do update set quantity=inventory.bulk_balance.quantity+excluded.quantity,version=inventory.bulk_balance.version+1,updated_at=clock_timestamp() where (select max_quantity is null or inventory.bulk_balance.quantity+excluded.quantity<=max_quantity from inventory.item_bin_policy where tenant_id=$1 and organization_id=$2 and item_id=$4 and bin_id=$3) returning version`, tenant, organization, receiveBin, itemID, c.lot, c.quantity).Scan(&balanceVersion)
		if errors.Is(err, pgx.ErrNoRows) {
			return inventorycontrol.BulkTransfer{}, inventorycontrol.ErrConflict
		}
		if err != nil {
			return inventorycontrol.BulkTransfer{}, bulkConflict(err)
		}
		_, err = tx.Exec(ctx, `insert into inventory.bulk_inventory_entry(tenant_id,entry_id,organization_id,bin_id,item_id,lot_id,entry_type,quantity,unit_cost,cost_amount,posting_date,source_kind,source_id) values($1,$2,$3,$4,$5,nullif($6,''),'receipt',$7::numeric,$8::numeric,$9::numeric,$10::date,'bulk-transfer-receipt',$11)`, tenant, receiptEntryID, organization, receiveBin, itemID, c.lot, c.quantity, c.unit, c.cost, command.PostingDate, receiptID+":"+c.id)
		if err != nil {
			return inventorycontrol.BulkTransfer{}, bulkConflict(err)
		}
		_, err = tx.Exec(ctx, `insert into inventory.bulk_cost_layer(tenant_id,layer_id,receipt_entry_id,organization_id,item_id,lot_id,posting_date,original_quantity,remaining_quantity,unit_cost) values($1,$2,$3,$4,$5,nullif($6,''),$7::date,$8::numeric,$8::numeric,$9::numeric)`, tenant, layerID, receiptEntryID, organization, itemID, c.lot, command.PostingDate, c.quantity, c.unit)
		if err != nil {
			return inventorycontrol.BulkTransfer{}, bulkConflict(err)
		}
		_, err = tx.Exec(ctx, `insert into inventory.bulk_transfer_receipt_component(tenant_id,receipt_id,transfer_id,component_id,receipt_entry_id,quantity,cost_amount) values($1,$2,$3,$4,$5,$6::numeric,$7::numeric)`, tenant, receiptID, transferID, c.id, receiptEntryID, c.quantity, c.cost)
		if err != nil {
			return inventorycontrol.BulkTransfer{}, bulkConflict(err)
		}
		updated, updateErr := tx.Exec(ctx, `update inventory.bulk_transfer_cost_component set receipt_entry_id=coalesce(receipt_entry_id,$4),received_quantity=received_quantity+$5::numeric where tenant_id=$1 and transfer_id=$2 and component_id=$3 and received_quantity+$5::numeric<=quantity`, tenant, transferID, c.id, receiptEntryID, c.quantity)
		if updateErr != nil || updated.RowsAffected() != 1 {
			if updateErr != nil {
				return inventorycontrol.BulkTransfer{}, updateErr
			}
			return inventorycontrol.BulkTransfer{}, inventorycontrol.ErrConflict
		}
	}
	type lotTotal struct{ lot, quantity, sourceEntry string }
	lotOrder := []string{}
	lotTotals := map[string]*lotTotal{}
	for index, c := range components {
		current, exists := lotTotals[c.lot]
		if !exists {
			current = &lotTotal{lot: c.lot, quantity: "0", sourceEntry: fmt.Sprintf("%s-entry-%06d", receiptID, index+1)}
			lotTotals[c.lot] = current
			lotOrder = append(lotOrder, c.lot)
		}
		total, _ := parseRat(current.quantity)
		part, _ := parseRat(c.quantity)
		current.quantity = formatRat(total.Add(total, part), 6)
	}
	sequence := 0
	for _, lotID := range lotOrder {
		lot := lotTotals[lotID]
		rows, queryErr := tx.Query(ctx, `select p.bin_id,coalesce(p.max_quantity::text,''),coalesce(b.quantity::text,'0'),coalesce((select sum(al.quantity)::text from inventory.warehouse_activity_line al join inventory.warehouse_activity a using(tenant_id,activity_id) where al.tenant_id=p.tenant_id and a.organization_id=p.organization_id and a.activity_type='put-away' and a.status='open' and al.to_bin_id=p.bin_id and al.item_id=p.item_id),'0') from inventory.item_bin_policy p join inventory.warehouse_bin w using(tenant_id,organization_id,bin_id) left join inventory.bulk_balance b on b.tenant_id=p.tenant_id and b.organization_id=p.organization_id and b.bin_id=p.bin_id and b.item_id=p.item_id and b.lot_id is not distinct from nullif($4,'') where p.tenant_id=$1 and p.organization_id=$2 and p.item_id=$3 and p.bin_id<>$5 and w.bin_type in ('put-away','putpick') and not w.movement_blocked order by p.is_default desc,w.bin_rank desc,p.fixed desc,w.bin_code for share of p,w`, tenant, organization, itemID, lotID, receiveBin)
		if queryErr != nil {
			return inventorycontrol.BulkTransfer{}, queryErr
		}
		remaining, ok := parseRat(lot.quantity)
		if !ok {
			rows.Close()
			return inventorycontrol.BulkTransfer{}, inventorycontrol.ErrConflict
		}
		placements := []warehousePlacement{}
		for rows.Next() && remaining.Sign() > 0 {
			var binID, maximum, onHand, planned string
			if err = rows.Scan(&binID, &maximum, &onHand, &planned); err != nil {
				rows.Close()
				return inventorycontrol.BulkTransfer{}, err
			}
			capacity, valid := decimalCapacity(maximum, onHand, planned)
			if !valid {
				rows.Close()
				return inventorycontrol.BulkTransfer{}, inventorycontrol.ErrConflict
			}
			take := new(big.Rat).Set(remaining)
			if capacity != nil {
				take = minRat(take, capacity)
			}
			if take.Sign() > 0 {
				placements = append(placements, warehousePlacement{binID: binID, quantity: formatRat(take, 6)})
				remaining.Sub(remaining, take)
			}
		}
		rows.Close()
		if err = rows.Err(); err != nil {
			return inventorycontrol.BulkTransfer{}, err
		}
		if remaining.Sign() != 0 || len(placements) == 0 {
			return inventorycontrol.BulkTransfer{}, inventorycontrol.ErrConflict
		}
		var expiration *time.Time
		if lotID != "" {
			if err = tx.QueryRow(ctx, `select expiration_date from inventory.inventory_lot where tenant_id=$1 and lot_id=$2`, tenant, lotID).Scan(&expiration); err != nil {
				return inventorycontrol.BulkTransfer{}, err
			}
		}
		for _, placement := range placements {
			sequence++
			activityLineID := fmt.Sprintf("%s-%06d", putAwayID, sequence)
			_, err = tx.Exec(ctx, `insert into inventory.warehouse_activity_line(tenant_id,activity_id,organization_id,line_id,sequence_no,from_bin_id,to_bin_id,item_id,lot_id,quantity,expiration_date,source_entry_id,uom_code,uom_quantity,qty_per_uom,from_uom_code,from_uom_quantity,from_qty_per_uom) select $1,$2,$3,$4,$5,$6,$7,$8,nullif($9,''),$10::numeric,$11,$12,u.uom_code,$10::numeric/u.qty_per_uom,u.qty_per_uom,u.uom_code,$10::numeric/u.qty_per_uom,u.qty_per_uom from inventory.item_unit_of_measure u where u.tenant_id=$1 and u.item_id=$8 and u.is_base`, tenant, putAwayID, organization, activityLineID, sequence, receiveBin, placement.binID, itemID, lotID, placement.quantity, expiration, lot.sourceEntry)
			if err != nil {
				return inventorycontrol.BulkTransfer{}, bulkConflict(err)
			}
			var receiveBalanceID int64
			err = tx.QueryRow(ctx, `update inventory.bulk_balance set reserved_quantity=reserved_quantity+$6::numeric,version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and organization_id=$2 and bin_id=$3 and item_id=$4 and lot_id is not distinct from nullif($5,'') and quantity-reserved_quantity >= $6::numeric returning balance_id`, tenant, organization, receiveBin, itemID, lotID, placement.quantity).Scan(&receiveBalanceID)
			if errors.Is(err, pgx.ErrNoRows) {
				return inventorycontrol.BulkTransfer{}, inventorycontrol.ErrConflict
			}
			if err != nil {
				return inventorycontrol.BulkTransfer{}, err
			}
			reserved, reserveErr := tx.Exec(ctx, `update inventory.bulk_uom_balance composition set reserved_quantity=composition.reserved_quantity+($3::numeric/u.qty_per_uom),version=composition.version+1,updated_at=clock_timestamp() from inventory.item_unit_of_measure u where composition.balance_id=$1 and u.tenant_id=$2 and u.item_id=$4 and u.is_base and composition.uom_code=u.uom_code and composition.quantity-composition.reserved_quantity >= ($3::numeric/u.qty_per_uom)`, receiveBalanceID, tenant, placement.quantity, itemID)
			if reserveErr != nil {
				return inventorycontrol.BulkTransfer{}, bulkConflict(reserveErr)
			}
			if reserved.RowsAffected() != 1 {
				return inventorycontrol.BulkTransfer{}, inventorycontrol.ErrConflict
			}
		}
	}
	receivedAfter := new(big.Rat).Add(receivedBefore, postingQuantity)
	status := "partially-received"
	finalReceipt := receivedAfter.Cmp(total) == 0
	if finalReceipt {
		status = "received"
	}
	updated, err := tx.Exec(ctx, `update inventory.bulk_transfer set status=$4,version=version+1,received_at=case when $5 then clock_timestamp() else null end,put_away_activity_id=$6,updated_at=clock_timestamp() where tenant_id=$1 and transfer_id=$2 and status in ('partially-shipped','shipped','partially-received') and version=$3`, tenant, transferID, version, status, finalReceipt, putAwayID)
	if err != nil || updated.RowsAffected() != 1 {
		if err != nil {
			return inventorycontrol.BulkTransfer{}, err
		}
		return inventorycontrol.BulkTransfer{}, inventorycontrol.ErrConflict
	}
	_, err = tx.Exec(ctx, `update inventory.bulk_transfer_line set received_quantity=received_quantity+$3::numeric where tenant_id=$1 and transfer_id=$2`, tenant, transferID, command.Quantity)
	if err != nil {
		return inventorycontrol.BulkTransfer{}, err
	}
	if err = recordBulkEvent(ctx, tx, tenant, eventID, "bulk-transfer", transferID, "bulk-transfer.received", version+1, map[string]any{"receipt_id": receiptID, "request_id": command.RequestID, "organization_id": organization, "item_id": itemID, "quantity": command.Quantity, "cost_amount": formatRat(totalCost, 4), "cost_components": len(components), "receive_bin_id": receiveBin, "put_away_id": putAwayID, "put_away_lines": sequence}); err != nil {
		return inventorycontrol.BulkTransfer{}, err
	}
	value, err := readBulkTransfer(ctx, tx, tenant, transferID)
	if err != nil {
		return value, err
	}
	value.PostingID, value.PutAwayID = receiptID, putAwayID
	return value, tx.Commit(ctx)
}

func (r *InventoryControl) CancelBulkTransfer(ctx context.Context, tenant, transferID, fromOrganization, toOrganization string, version int64, eventID string) (inventorycontrol.BulkTransfer, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return inventorycontrol.BulkTransfer{}, err
	}
	defer tx.Rollback(ctx)
	var organization, itemID string
	err = tx.QueryRow(ctx, `select t.from_organization_id,l.item_id from inventory.bulk_transfer t join inventory.bulk_transfer_line l using(tenant_id,transfer_id) where t.tenant_id=$1 and t.transfer_id=$2 and t.from_organization_id=$3 and t.to_organization_id=$4 and t.status='released' and t.version=$5 for update of t,l`, tenant, transferID, fromOrganization, toOrganization, version).Scan(&organization, &itemID)
	if errors.Is(err, pgx.ErrNoRows) {
		return inventorycontrol.BulkTransfer{}, inventorycontrol.ErrConflict
	}
	if err != nil {
		return inventorycontrol.BulkTransfer{}, err
	}
	if err = warehouseLock(ctx, tx, tenant, organization, itemID); err != nil {
		return inventorycontrol.BulkTransfer{}, err
	}
	rows, err := tx.Query(ctx, `select a.activity_id,al.reservation_id,al.from_bin_id,coalesce(al.lot_id,''),al.quantity::text,r.version from inventory.warehouse_activity a join inventory.warehouse_activity_line al using(tenant_id,activity_id) join inventory.bulk_reservation r on r.tenant_id=al.tenant_id and r.reservation_id=al.reservation_id where a.tenant_id=$1 and a.organization_id=$2 and a.activity_type='pick' and a.status='open' and a.source_kind='transfer-outbound' and a.source_id=$3 and r.status='reservation' order by a.activity_id,al.sequence_no for update of a,r`, tenant, organization, transferID)
	if err != nil {
		return inventorycontrol.BulkTransfer{}, err
	}
	activities := map[string]bool{}
	type transferRelease struct {
		reservation, bin, lot, quantity string
		reservationVersion              int64
	}
	releases := []transferRelease{}
	for rows.Next() {
		var activity, reservation, bin, lot, quantity string
		var reservationVersion int64
		if err = rows.Scan(&activity, &reservation, &bin, &lot, &quantity, &reservationVersion); err != nil {
			rows.Close()
			return inventorycontrol.BulkTransfer{}, err
		}
		activities[activity] = true
		releases = append(releases, transferRelease{reservation: reservation, bin: bin, lot: lot, quantity: quantity, reservationVersion: reservationVersion})
	}
	rows.Close()
	if err = rows.Err(); err != nil {
		return inventorycontrol.BulkTransfer{}, err
	}
	for _, release := range releases {
		updated, updateErr := tx.Exec(ctx, `update inventory.bulk_balance set reserved_quantity=reserved_quantity-$6::numeric,version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and organization_id=$2 and bin_id=$3 and item_id=$4 and lot_id is not distinct from nullif($5,'') and reserved_quantity >= $6::numeric`, tenant, organization, release.bin, itemID, release.lot, release.quantity)
		if updateErr != nil || updated.RowsAffected() != 1 {
			if updateErr != nil {
				return inventorycontrol.BulkTransfer{}, updateErr
			}
			return inventorycontrol.BulkTransfer{}, inventorycontrol.ErrConflict
		}
		updated, updateErr = tx.Exec(ctx, `update inventory.bulk_reservation set status='released',cancellation_disallowed=false,version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and reservation_id=$2 and status='reservation' and version=$3`, tenant, release.reservation, release.reservationVersion)
		if updateErr != nil || updated.RowsAffected() != 1 {
			if updateErr != nil {
				return inventorycontrol.BulkTransfer{}, updateErr
			}
			return inventorycontrol.BulkTransfer{}, inventorycontrol.ErrConflict
		}
	}
	for activity := range activities {
		updated, updateErr := tx.Exec(ctx, `update inventory.warehouse_activity set status='cancelled',version=version+1 where tenant_id=$1 and activity_id=$2 and status='open'`, tenant, activity)
		if updateErr != nil || updated.RowsAffected() != 1 {
			if updateErr != nil {
				return inventorycontrol.BulkTransfer{}, updateErr
			}
			return inventorycontrol.BulkTransfer{}, inventorycontrol.ErrConflict
		}
	}
	updated, err := tx.Exec(ctx, `update inventory.bulk_transfer set status='cancelled',version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and transfer_id=$2 and status='released' and version=$3`, tenant, transferID, version)
	if err != nil || updated.RowsAffected() != 1 {
		if err != nil {
			return inventorycontrol.BulkTransfer{}, err
		}
		return inventorycontrol.BulkTransfer{}, inventorycontrol.ErrConflict
	}
	if err = recordBulkEvent(ctx, tx, tenant, eventID, "bulk-transfer", transferID, "bulk-transfer.cancelled", version+1, map[string]any{"organization_id": organization, "cancelled_pick_activities": len(activities)}); err != nil {
		return inventorycontrol.BulkTransfer{}, err
	}
	value, err := readBulkTransfer(ctx, tx, tenant, transferID)
	if err != nil {
		return value, err
	}
	return value, tx.Commit(ctx)
}
