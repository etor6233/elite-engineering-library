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

type warehousePlacement struct {
	binID                 string
	quantity              string
	uomCode               string
	uomQuantity           string
	quantityPerUOM        string
	crossDockAllocationID string
	transferID            string
	transferLineID        string
	dueDate               time.Time
	lotID                 string
}

func warehouseLock(ctx context.Context, tx pgx.Tx, tenant, organization, item string) error {
	_, err := tx.Exec(ctx, `select pg_advisory_xact_lock(hashtextextended($1 || ':' || $2 || ':' || $3,0))`, tenant, organization, item)
	return err
}

func decimalCapacity(maximum, onHand, planned string) (*big.Rat, bool) {
	if maximum == "" {
		return nil, true
	}
	max, okMax := parseRat(maximum)
	stock, okStock := parseRat(onHand)
	pending, okPending := parseRat(planned)
	if !okMax || !okStock || !okPending {
		return nil, false
	}
	available := new(big.Rat).Sub(max, stock)
	available.Sub(available, pending)
	if available.Sign() < 0 {
		available.SetInt64(0)
	}
	return available, true
}

func exactUOMQuantity(baseQuantity, quantityPerUOM, roundingPrecision string) (string, bool) {
	base, baseOK := parseRat(baseQuantity)
	factor, factorOK := parseRat(quantityPerUOM)
	precision, precisionOK := parseRat(roundingPrecision)
	if !baseOK || !factorOK || !precisionOK || factor.Sign() <= 0 || precision.Sign() <= 0 {
		return "", false
	}
	quantity := new(big.Rat).Quo(base, factor)
	steps := new(big.Rat).Quo(quantity, precision)
	if !steps.IsInt() {
		return "", false
	}
	return formatRat(quantity, 6), true
}

func (r *InventoryControl) PostWarehouseReceipt(ctx context.Context, tenant string, ids inventorycontrol.WarehouseIDs, value inventorycontrol.WarehouseReceiptCommand) (inventorycontrol.WarehouseReceiptResult, error) {
	result := inventorycontrol.WarehouseReceiptResult{ReceiptID: ids.ReceiptID, EntryID: ids.EntryID}
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return result, err
	}
	defer tx.Rollback(ctx)
	if err = warehouseLock(ctx, tx, tenant, value.OrganizationID, value.ItemID); err != nil {
		return result, err
	}
	var tracking, receiveType string
	err = tx.QueryRow(ctx, `select i.tracking_mode,w.bin_type from inventory.stock_item i join inventory.item_bin_policy p on p.tenant_id=i.tenant_id and p.item_id=i.item_id join inventory.warehouse_bin w on w.tenant_id=p.tenant_id and w.organization_id=p.organization_id and w.bin_id=p.bin_id where i.tenant_id=$1 and i.item_id=$2 and p.organization_id=$3 and p.bin_id=$4 and not w.movement_blocked for share of i,p,w`, tenant, value.ItemID, value.OrganizationID, value.ReceiveBinID).Scan(&tracking, &receiveType)
	if errors.Is(err, pgx.ErrNoRows) {
		return result, fmt.Errorf("warehouse receipt receive-bin contract: %w", inventorycontrol.ErrConflict)
	}
	if err != nil {
		return result, err
	}
	if receiveType != "receive" {
		return result, fmt.Errorf("warehouse receipt receive-bin contract: %w", inventorycontrol.ErrConflict)
	}
	handlingUOM, handlingQuantity := value.HandlingUOM, value.HandlingQuantity
	if handlingUOM == "" {
		handlingQuantity = value.Quantity
	}
	var baseUOM, baseQuantityPerUOM, basePrecision, quantityPerUOM, handlingPrecision, baseQuantity string
	err = tx.QueryRow(ctx, `select u.uom_code,u.qty_per_uom::text,u.rounding_precision::text,b.uom_code,b.qty_per_uom::text,b.rounding_precision::text,($4::numeric*u.qty_per_uom)::numeric(20,6)::text
from inventory.item_unit_of_measure u
join inventory.item_unit_of_measure b on b.tenant_id=u.tenant_id and b.item_id=u.item_id and b.is_base
where u.tenant_id=$1 and u.item_id=$2 and u.uom_code=coalesce(nullif($3,''),b.uom_code)
  and mod($4::numeric,u.rounding_precision)=0
  and mod($4::numeric*u.qty_per_uom,b.rounding_precision)=0
for share of u,b`, tenant, value.ItemID, handlingUOM, handlingQuantity).Scan(&handlingUOM, &quantityPerUOM, &handlingPrecision, &baseUOM, &baseQuantityPerUOM, &basePrecision, &baseQuantity)
	if errors.Is(err, pgx.ErrNoRows) {
		return result, fmt.Errorf("warehouse receipt handling UOM contract: %w", inventorycontrol.ErrConflict)
	}
	if err != nil {
		return result, bulkConflict(err)
	}
	result.Quantity, result.HandlingUOM, result.HandlingQuantity = baseQuantity, handlingUOM, handlingQuantity
	lotID := ""
	if tracking == "lot" {
		if value.LotNo == "" {
			return result, fmt.Errorf("warehouse receipt lot required: %w", inventorycontrol.ErrConflict)
		}
		var existingExpiration, existingWarranty string
		err = tx.QueryRow(ctx, `select lot_id,coalesce(expiration_date::text,''),coalesce(warranty_date::text,'') from inventory.inventory_lot where tenant_id=$1 and item_id=$2 and lot_no=$3 for update`, tenant, value.ItemID, value.LotNo).Scan(&lotID, &existingExpiration, &existingWarranty)
		expiration, warranty := "", ""
		if value.ExpirationDate != nil {
			expiration = value.ExpirationDate.UTC().Format("2006-01-02")
		}
		if value.WarrantyDate != nil {
			warranty = value.WarrantyDate.UTC().Format("2006-01-02")
		}
		if errors.Is(err, pgx.ErrNoRows) {
			lotID = ids.LotID
			_, err = tx.Exec(ctx, `insert into inventory.inventory_lot(tenant_id,lot_id,item_id,lot_no,expiration_date,warranty_date) values($1,$2,$3,$4,nullif($5,'')::date,nullif($6,'')::date)`, tenant, lotID, value.ItemID, value.LotNo, expiration, warranty)
		} else if err == nil && (existingExpiration != expiration || existingWarranty != warranty) {
			return result, fmt.Errorf("warehouse receipt lot metadata mismatch: %w", inventorycontrol.ErrConflict)
		}
		if err != nil {
			return result, fmt.Errorf("warehouse receipt lot: %w", bulkConflict(err))
		}
	} else if value.LotNo != "" || value.ExpirationDate != nil || value.WarrantyDate != nil {
		return result, fmt.Errorf("warehouse receipt unexpected tracking: %w", inventorycontrol.ErrConflict)
	}
	result.LotID = lotID
	var balanceID, balanceVersion int64
	err = tx.QueryRow(ctx, `insert into inventory.bulk_balance(tenant_id,organization_id,bin_id,item_id,lot_id,quantity,reserved_quantity,version) values($1,$2,$3,$4,nullif($5,''),$6::numeric,0,1) on conflict (tenant_id,organization_id,bin_id,item_id,lot_id) do update set quantity=inventory.bulk_balance.quantity+excluded.quantity,version=inventory.bulk_balance.version+1,updated_at=clock_timestamp() where (select max_quantity is null or inventory.bulk_balance.quantity+excluded.quantity<=max_quantity from inventory.item_bin_policy where tenant_id=$1 and organization_id=$2 and item_id=$4 and bin_id=$3) returning balance_id,version`, tenant, value.OrganizationID, value.ReceiveBinID, value.ItemID, lotID, baseQuantity).Scan(&balanceID, &balanceVersion)
	if errors.Is(err, pgx.ErrNoRows) {
		return result, fmt.Errorf("warehouse receipt receive capacity: %w", inventorycontrol.ErrConflict)
	}
	if err != nil {
		return result, fmt.Errorf("warehouse receipt balance: %w", bulkConflict(err))
	}
	err = tx.QueryRow(ctx, `insert into inventory.bulk_inventory_entry(tenant_id,entry_id,organization_id,bin_id,item_id,lot_id,entry_type,quantity,unit_cost,cost_amount,posting_date,source_kind,source_id) values($1,$2,$3,$4,$5,nullif($6,''),'receipt',$7::numeric,$8::numeric,round($7::numeric*$8::numeric,4),$9::date,$10,$11) returning cost_amount::text`, tenant, ids.EntryID, value.OrganizationID, value.ReceiveBinID, value.ItemID, lotID, baseQuantity, value.UnitCost, value.PostingDate, value.SourceKind, value.SourceID).Scan(&result.CostAmount)
	if err != nil {
		return result, fmt.Errorf("warehouse receipt entry: %w", bulkConflict(err))
	}
	_, err = tx.Exec(ctx, `insert into inventory.bulk_cost_layer(tenant_id,layer_id,receipt_entry_id,organization_id,item_id,lot_id,posting_date,original_quantity,remaining_quantity,unit_cost) values($1,$2,$3,$4,$5,nullif($6,''),$7::date,$8::numeric,$8::numeric,$9::numeric)`, tenant, ids.LayerID, ids.EntryID, value.OrganizationID, value.ItemID, lotID, value.PostingDate, baseQuantity, value.UnitCost)
	if err != nil {
		return result, fmt.Errorf("warehouse receipt cost layer: %w", bulkConflict(err))
	}
	_, err = tx.Exec(ctx, `insert into inventory.warehouse_receipt(tenant_id,receipt_id,organization_id,receive_bin_id,source_kind,source_id,posting_date,status,allow_breakbulk) values($1,$2,$3,$4,$5,$6,$7::date,'posted',$8)`, tenant, ids.ReceiptID, value.OrganizationID, value.ReceiveBinID, value.SourceKind, value.SourceID, value.PostingDate, value.AllowBreakbulk)
	if err != nil {
		return result, fmt.Errorf("warehouse receipt header: %w", bulkConflict(err))
	}
	_, err = tx.Exec(ctx, `insert into inventory.warehouse_receipt_line(tenant_id,receipt_id,line_id,item_id,lot_id,quantity,unit_cost,inventory_entry_id,uom_code,uom_quantity,qty_per_uom) values($1,$2,$3,$4,nullif($5,''),$6::numeric,$7::numeric,$8,$9,$10::numeric,$11::numeric)`, tenant, ids.ReceiptID, ids.LineID, value.ItemID, lotID, baseQuantity, value.UnitCost, ids.EntryID, handlingUOM, handlingQuantity, quantityPerUOM)
	if err != nil {
		return result, fmt.Errorf("warehouse receipt line: %w", bulkConflict(err))
	}
	planningValue := value
	planningValue.Quantity = baseQuantity
	crossDockPlacements, remaining, err := planWarehouseCrossDock(ctx, tx, tenant, ids, planningValue, lotID)
	if err != nil {
		return result, fmt.Errorf("warehouse receipt cross-dock plan: %w", err)
	}
	rows, err := tx.Query(ctx, `select p.bin_id,coalesce(p.max_quantity::text,''),coalesce(b.quantity::text,'0'),coalesce((select sum(al.quantity)::text from inventory.warehouse_activity_line al join inventory.warehouse_activity a using(tenant_id,activity_id) where al.tenant_id=p.tenant_id and a.organization_id=p.organization_id and a.activity_type='put-away' and a.status='open' and al.to_bin_id=p.bin_id and al.item_id=p.item_id),'0') from inventory.item_bin_policy p join inventory.warehouse_bin w using(tenant_id,organization_id,bin_id) left join inventory.bulk_balance b on b.tenant_id=p.tenant_id and b.organization_id=p.organization_id and b.bin_id=p.bin_id and b.item_id=p.item_id and b.lot_id is not distinct from nullif($4,'') where p.tenant_id=$1 and p.organization_id=$2 and p.item_id=$3 and p.bin_id<>$5 and w.bin_type in ('put-away','putpick') and not w.cross_dock and not w.movement_blocked order by p.is_default desc,w.bin_rank desc,p.fixed desc,w.bin_code for share of p,w`, tenant, value.OrganizationID, value.ItemID, lotID, value.ReceiveBinID)
	if err != nil {
		return result, err
	}
	placements := crossDockPlacements
	for rows.Next() && remaining.Sign() > 0 {
		var binID, maximum, onHand, planned string
		if err = rows.Scan(&binID, &maximum, &onHand, &planned); err != nil {
			rows.Close()
			return result, err
		}
		capacity, valid := decimalCapacity(maximum, onHand, planned)
		if !valid {
			rows.Close()
			return result, fmt.Errorf("warehouse receipt target capacity parse: %w", inventorycontrol.ErrConflict)
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
		return result, err
	}
	if remaining.Sign() != 0 || len(placements) == 0 {
		return result, fmt.Errorf("warehouse receipt put-away capacity: %w", inventorycontrol.ErrConflict)
	}
	requiresBreakbulk := false
	for index := range placements {
		quantity, exact := exactUOMQuantity(placements[index].quantity, quantityPerUOM, handlingPrecision)
		if !exact {
			requiresBreakbulk = handlingUOM != baseUOM
			break
		}
		placements[index].uomCode = handlingUOM
		placements[index].uomQuantity = quantity
		placements[index].quantityPerUOM = quantityPerUOM
	}
	if requiresBreakbulk && !value.AllowBreakbulk {
		return result, fmt.Errorf("warehouse receipt requires explicit breakbulk authorization: %w", inventorycontrol.ErrConflict)
	}
	if requiresBreakbulk {
		for index := range placements {
			quantity, exact := exactUOMQuantity(placements[index].quantity, baseQuantityPerUOM, basePrecision)
			if !exact {
				return result, fmt.Errorf("warehouse receipt base UOM placement: %w", inventorycontrol.ErrConflict)
			}
			placements[index].uomCode = baseUOM
			placements[index].uomQuantity = quantity
			placements[index].quantityPerUOM = baseQuantityPerUOM
		}
	}
	if handlingUOM != baseUOM && !requiresBreakbulk {
		updated, updateErr := tx.Exec(ctx, `update inventory.bulk_uom_balance set quantity=quantity-($3::numeric/qty_per_uom),quantity_base=quantity_base-$3::numeric,version=version+1,updated_at=clock_timestamp() where balance_id=$1 and uom_code=$2 and quantity-reserved_quantity >= ($3::numeric/qty_per_uom) and quantity_base >= $3::numeric`, balanceID, baseUOM, baseQuantity)
		if updateErr != nil || updated.RowsAffected() != 1 {
			if updateErr != nil {
				return result, bulkConflict(updateErr)
			}
			return result, inventorycontrol.ErrConflict
		}
		_, err = tx.Exec(ctx, `delete from inventory.bulk_uom_balance where balance_id=$1 and uom_code=$2 and quantity_base=0 and reserved_quantity=0`, balanceID, baseUOM)
		if err != nil {
			return result, err
		}
		_, err = tx.Exec(ctx, `insert into inventory.bulk_uom_balance(balance_id,uom_code,quantity,qty_per_uom,quantity_base) values($1,$2,$3::numeric,$4::numeric,$5::numeric) on conflict(balance_id,uom_code) do update set quantity=inventory.bulk_uom_balance.quantity+excluded.quantity,quantity_base=inventory.bulk_uom_balance.quantity_base+excluded.quantity_base,version=inventory.bulk_uom_balance.version+1,updated_at=clock_timestamp()`, balanceID, handlingUOM, handlingQuantity, quantityPerUOM, baseQuantity)
		if err != nil {
			return result, bulkConflict(err)
		}
	}
	if requiresBreakbulk {
		baseHandlingQuantity, exact := exactUOMQuantity(baseQuantity, baseQuantityPerUOM, basePrecision)
		if !exact {
			return result, inventorycontrol.ErrConflict
		}
		_, err = tx.Exec(ctx, `insert into inventory.bulk_uom_conversion(tenant_id,conversion_id,request_id,organization_id,bin_id,item_id,lot_id,operation,from_uom,to_uom,from_quantity,to_quantity,base_quantity,from_qty_per_uom,to_qty_per_uom,source_kind,source_id) values($1,$2,$3,$4,$5,$6,nullif($7,''),'breakbulk',$8,$9,$10::numeric,$11::numeric,$12::numeric,$13::numeric,$14::numeric,'warehouse-receipt',$15)`, tenant, ids.ConversionID, value.RequestID+"/automatic-breakbulk", value.OrganizationID, value.ReceiveBinID, value.ItemID, lotID, handlingUOM, baseUOM, handlingQuantity, baseHandlingQuantity, baseQuantity, quantityPerUOM, baseQuantityPerUOM, ids.ReceiptID)
		if err != nil {
			return result, fmt.Errorf("warehouse receipt automatic breakbulk: %w", bulkConflict(err))
		}
		_, err = tx.Exec(ctx, `insert into inventory.bulk_uom_conversion_line(tenant_id,conversion_id,line_id,sequence_no,action_type,uom_code,quantity,quantity_base) values($1,$2,$3,1,'take',$5,$6::numeric,$7::numeric),($1,$2,$4,2,'place',$8,$9::numeric,$7::numeric)`, tenant, ids.ConversionID, ids.TakeLineID, ids.PlaceLineID, handlingUOM, handlingQuantity, baseQuantity, baseUOM, baseHandlingQuantity)
		if err != nil {
			return result, fmt.Errorf("warehouse receipt automatic breakbulk lines: %w", bulkConflict(err))
		}
		if err = recordBulkEvent(ctx, tx, tenant, ids.ConversionEventID, "bulk-uom-conversion", ids.ConversionID, "bulk-uom-conversion.completed", 1, map[string]string{"request_id": value.RequestID + "/automatic-breakbulk", "organization_id": value.OrganizationID, "bin_id": value.ReceiveBinID, "item_id": value.ItemID, "lot_id": lotID, "operation": "breakbulk", "from_uom": handlingUOM, "to_uom": baseUOM, "from_quantity": handlingQuantity, "to_quantity": baseHandlingQuantity, "base_quantity": baseQuantity, "source_kind": "warehouse-receipt", "source_id": ids.ReceiptID}); err != nil {
			return result, err
		}
	}
	_, err = tx.Exec(ctx, `insert into inventory.warehouse_activity(tenant_id,activity_id,organization_id,activity_type,status,source_kind,source_id,source_line_id,request_id,version) values($1,$2,$3,'put-away','open','warehouse-receipt',$4,$5,$6,1)`, tenant, ids.ActivityID, value.OrganizationID, ids.ReceiptID, ids.LineID, value.RequestID)
	if err != nil {
		return result, fmt.Errorf("warehouse put-away header: %w", bulkConflict(err))
	}
	activity := inventorycontrol.WarehouseActivity{ID: ids.ActivityID, OrganizationID: value.OrganizationID, Type: "put-away", Status: "open", SourceKind: "warehouse-receipt", SourceID: ids.ReceiptID, SourceLineID: ids.LineID, RequestID: value.RequestID, Version: 1, Lines: []inventorycontrol.WarehouseActivityLine{}}
	for index, placement := range placements {
		lineID := fmt.Sprintf("%s-%06d", ids.ActivityID, index+1)
		_, err = tx.Exec(ctx, `insert into inventory.warehouse_activity_line(tenant_id,activity_id,organization_id,line_id,sequence_no,from_bin_id,to_bin_id,item_id,lot_id,quantity,expiration_date,source_entry_id,uom_code,uom_quantity,qty_per_uom,from_uom_code,from_uom_quantity,from_qty_per_uom) values($1,$2,$3,$4,$5,$6,$7,$8,nullif($9,''),$10::numeric,$11,$12,$13,$14::numeric,$15::numeric,$13,$14::numeric,$15::numeric)`, tenant, ids.ActivityID, value.OrganizationID, lineID, index+1, value.ReceiveBinID, placement.binID, value.ItemID, lotID, placement.quantity, value.ExpirationDate, ids.EntryID, placement.uomCode, placement.uomQuantity, placement.quantityPerUOM)
		if err != nil {
			return result, fmt.Errorf("warehouse put-away line: %w", bulkConflict(err))
		}
		if placement.crossDockAllocationID != "" {
			_, err = tx.Exec(ctx, `insert into inventory.warehouse_crossdock_allocation(tenant_id,allocation_id,receipt_id,receipt_line_id,put_away_activity_id,put_away_line_id,organization_id,bin_id,item_id,lot_id,transfer_id,transfer_line_id,quantity,due_date,available) values($1,$2,$3,$4,$5,$6,$7,$8,$9,nullif($10,''),$11,$12,$13::numeric,$14::date,false)`, tenant, placement.crossDockAllocationID, ids.ReceiptID, ids.LineID, ids.ActivityID, lineID, value.OrganizationID, placement.binID, value.ItemID, lotID, placement.transferID, placement.transferLineID, placement.quantity, placement.dueDate)
			if err != nil {
				return result, fmt.Errorf("warehouse cross-dock allocation: %w", bulkConflict(err))
			}
		}
		activity.Lines = append(activity.Lines, inventorycontrol.WarehouseActivityLine{ID: lineID, Sequence: index + 1, FromBinID: value.ReceiveBinID, ToBinID: placement.binID, ItemID: value.ItemID, LotID: lotID, LotNo: value.LotNo, Quantity: placement.quantity, FromUOMCode: placement.uomCode, FromUOMQuantity: placement.uomQuantity, FromQuantityPerUOM: placement.quantityPerUOM, UOMCode: placement.uomCode, UOMQuantity: placement.uomQuantity, QuantityPerUOM: placement.quantityPerUOM, ExpirationDate: value.ExpirationDate, SourceEntryID: ids.EntryID})
	}
	activityUOM := placements[0].uomCode
	activityUOMQuantity := baseQuantity
	if activityUOM == handlingUOM {
		activityUOMQuantity = handlingQuantity
	} else {
		activityUOMQuantity, _ = exactUOMQuantity(baseQuantity, baseQuantityPerUOM, basePrecision)
	}
	updated, updateErr := tx.Exec(ctx, `update inventory.bulk_balance set reserved_quantity=reserved_quantity+$2::numeric,version=version+1,updated_at=clock_timestamp() where balance_id=$1 and quantity-reserved_quantity >= $2::numeric`, balanceID, baseQuantity)
	if updateErr != nil || updated.RowsAffected() != 1 {
		if updateErr != nil {
			return result, updateErr
		}
		return result, inventorycontrol.ErrConflict
	}
	updated, updateErr = tx.Exec(ctx, `update inventory.bulk_uom_balance set reserved_quantity=reserved_quantity+$3::numeric,version=version+1,updated_at=clock_timestamp() where balance_id=$1 and uom_code=$2 and quantity-reserved_quantity >= $3::numeric`, balanceID, activityUOM, activityUOMQuantity)
	if updateErr != nil || updated.RowsAffected() != 1 {
		if updateErr != nil {
			return result, bulkConflict(updateErr)
		}
		return result, inventorycontrol.ErrConflict
	}
	result.PutAway = activity
	if err = recordBulkEvent(ctx, tx, tenant, ids.EventID, "warehouse-receipt", ids.ReceiptID, "warehouse-receipt.posted", 1, map[string]any{"organization_id": value.OrganizationID, "item_id": value.ItemID, "lot_id": lotID, "quantity": baseQuantity, "handling_uom": handlingUOM, "handling_quantity": handlingQuantity, "allow_breakbulk": value.AllowBreakbulk, "automatic_breakbulk": requiresBreakbulk, "put_away_id": ids.ActivityID, "put_away_lines": len(activity.Lines), "cross_dock_allocations": len(crossDockPlacements)}); err != nil {
		return result, err
	}
	return result, tx.Commit(ctx)
}

func (r *InventoryControl) CreateWarehousePick(ctx context.Context, tenant string, ids inventorycontrol.WarehouseIDs, value inventorycontrol.WarehousePickCommand) (inventorycontrol.WarehouseActivity, error) {
	activity := inventorycontrol.WarehouseActivity{ID: ids.ActivityID, OrganizationID: value.OrganizationID, Type: "pick", Status: "open", SourceKind: value.DemandKind, SourceID: value.DemandID, SourceLineID: value.DemandLineID, RequestID: value.RequestID, AssignedTo: value.AssignedTo, Version: 1, Lines: []inventorycontrol.WarehouseActivityLine{}}
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return activity, err
	}
	defer tx.Rollback(ctx)
	if _, err = tx.Exec(ctx, `select pg_advisory_xact_lock(hashtextextended($1,0))`, tenant+"\x1fwarehouse-pick-request\x1f"+value.RequestID); err != nil {
		return activity, err
	}
	if replay, found, replayErr := readWarehousePickReplay(ctx, tx, tenant, value); replayErr != nil {
		return activity, replayErr
	} else if found {
		if err = tx.Commit(ctx); err != nil {
			return activity, err
		}
		return replay, nil
	}
	if err = warehouseLock(ctx, tx, tenant, value.OrganizationID, value.ItemID); err != nil {
		return activity, err
	}
	var tracking, targetType string
	err = tx.QueryRow(ctx, `select i.tracking_mode,w.bin_type from inventory.stock_item i join inventory.item_bin_policy p on p.tenant_id=i.tenant_id and p.item_id=i.item_id join inventory.warehouse_bin w on w.tenant_id=p.tenant_id and w.organization_id=p.organization_id and w.bin_id=p.bin_id where i.tenant_id=$1 and i.item_id=$2 and p.organization_id=$3 and p.bin_id=$4 and not w.movement_blocked for share of i,p,w`, tenant, value.ItemID, value.OrganizationID, value.ShipBinID).Scan(&tracking, &targetType)
	if errors.Is(err, pgx.ErrNoRows) {
		return activity, fmt.Errorf("warehouse pick ship-bin contract: %w", inventorycontrol.ErrConflict)
	}
	if err != nil {
		return activity, err
	}
	if targetType != "ship" {
		return activity, fmt.Errorf("warehouse pick ship-bin type: %w", inventorycontrol.ErrConflict)
	}
	demandSalesUOM, demandTotalBase, demandOutstandingBase := "", "", ""
	if value.DemandKind == "customer-order" {
		demandSalesUOM, demandTotalBase, demandOutstandingBase, err = validateCustomerOrderWarehouseDemand(ctx, tx, tenant, value)
		if err != nil {
			return activity, err
		}
	}
	if value.DemandKind == "transfer-outbound" {
		var outstandingText string
		err = tx.QueryRow(ctx, `select (l.quantity-l.shipped_quantity-coalesce((select sum(al.quantity) from inventory.warehouse_activity a join inventory.warehouse_activity_line al using(tenant_id,activity_id) where a.tenant_id=t.tenant_id and a.organization_id=t.from_organization_id and a.activity_type='pick' and a.source_kind='transfer-outbound' and a.source_id=t.transfer_id and a.source_line_id=l.line_id and (a.status='open' or (a.status='registered' and not exists(select 1 from inventory.bulk_transfer_shipment s where s.tenant_id=a.tenant_id and s.warehouse_activity_id=a.activity_id)))),0))::text from inventory.bulk_transfer t join inventory.bulk_transfer_line l using(tenant_id,transfer_id) where t.tenant_id=$1 and t.transfer_id=$2 and l.line_id=$3 and t.from_organization_id=$4 and l.item_id=$5 and t.status in ('released','partially-shipped','partially-received') for share of t,l`, tenant, value.DemandID, value.DemandLineID, value.OrganizationID, value.ItemID).Scan(&outstandingText)
		if errors.Is(err, pgx.ErrNoRows) {
			return activity, fmt.Errorf("warehouse pick transfer demand: %w", inventorycontrol.ErrConflict)
		}
		if err != nil {
			return activity, err
		}
		outstanding, validOutstanding := parseRat(outstandingText)
		requested, validRequested := parseRat(value.Quantity)
		if !validOutstanding || !validRequested || requested.Cmp(outstanding) > 0 {
			return activity, fmt.Errorf("warehouse pick transfer outstanding: %w", inventorycontrol.ErrConflict)
		}
	}
	remaining, ok := parseRat(value.Quantity)
	if !ok {
		return activity, fmt.Errorf("warehouse pick quantity parse: %w", inventorycontrol.ErrConflict)
	}
	targetUOM, err := resolveWarehouseUOM(ctx, tx, tenant, value.ItemID, value.HandlingUOM)
	if err != nil {
		return activity, fmt.Errorf("warehouse pick handling UOM: %w", err)
	}
	if _, exact := exactUOMQuantity(value.Quantity, targetUOM.factor, targetUOM.precision); !exact {
		return activity, fmt.Errorf("warehouse pick handling UOM quantity: %w", inventorycontrol.ErrConflict)
	}
	type candidate struct {
		binID, lotID, lotNo, available, allocationID string
		balanceID                                    int64
		expiration                                   *time.Time
		plan                                         warehouseUOMPlan
	}
	candidates := []candidate{}
	appendCandidates := func(rows pgx.Rows) error {
		raw := []candidate{}
		for rows.Next() {
			var c candidate
			if scanErr := rows.Scan(&c.balanceID, &c.binID, &c.lotID, &c.lotNo, &c.expiration, &c.available, &c.allocationID); scanErr != nil {
				rows.Close()
				return scanErr
			}
			raw = append(raw, c)
		}
		rows.Close()
		if rows.Err() != nil {
			return rows.Err()
		}
		for _, c := range raw {
			if remaining.Sign() <= 0 {
				break
			}
			available, valid := parseRat(c.available)
			if !valid || available.Sign() <= 0 {
				return fmt.Errorf("warehouse pick availability parse: %w", inventorycontrol.ErrConflict)
			}
			limit := minRat(available, remaining)
			plans, leftover, planErr := planWarehousePackaging(ctx, tx, c.balanceID, targetUOM, value.AllowBreakbulk, limit)
			if planErr != nil {
				return planErr
			}
			moved := new(big.Rat).Sub(limit, leftover)
			for _, plan := range plans {
				planned := c
				planned.plan = plan
				candidates = append(candidates, planned)
			}
			remaining.Sub(remaining, moved)
		}
		return nil
	}
	if value.DemandKind == "transfer-outbound" {
		rows, queryErr := tx.Query(ctx, `select b.balance_id,x.bin_id,coalesce(x.lot_id,''),coalesce(l.lot_no,''),l.expiration_date,least(x.quantity-x.reserved_quantity-x.picked_quantity,b.quantity-b.reserved_quantity)::text,x.allocation_id from inventory.warehouse_crossdock_allocation x join inventory.bulk_balance b on b.tenant_id=x.tenant_id and b.organization_id=x.organization_id and b.bin_id=x.bin_id and b.item_id=x.item_id and b.lot_id is not distinct from x.lot_id join inventory.warehouse_bin w on w.tenant_id=x.tenant_id and w.organization_id=x.organization_id and w.bin_id=x.bin_id left join inventory.inventory_lot l on l.tenant_id=x.tenant_id and l.lot_id=x.lot_id where x.tenant_id=$1 and x.organization_id=$2 and x.item_id=$3 and x.transfer_id=$4 and x.transfer_line_id=$5 and x.available and x.quantity-x.reserved_quantity-x.picked_quantity>0 and b.quantity-b.reserved_quantity>0 and w.cross_dock and not w.movement_blocked and coalesce(l.blocked,false)=false and (l.expiration_date is null or l.expiration_date>=current_date) order by x.due_date,w.bin_rank desc,w.bin_code,l.expiration_date nulls last,l.lot_no for update of x,b`, tenant, value.OrganizationID, value.ItemID, value.DemandID, value.DemandLineID)
		if queryErr != nil {
			return activity, queryErr
		}
		if err = appendCandidates(rows); err != nil {
			return activity, err
		}
	}
	if remaining.Sign() > 0 {
		query := `select b.balance_id,b.bin_id,coalesce(b.lot_id,''),coalesce(l.lot_no,''),l.expiration_date,(b.quantity-b.reserved_quantity)::text,'' from inventory.bulk_balance b join inventory.item_bin_policy p using(tenant_id,organization_id,item_id,bin_id) join inventory.warehouse_bin w using(tenant_id,organization_id,bin_id) left join inventory.inventory_lot l on l.tenant_id=b.tenant_id and l.lot_id=b.lot_id where b.tenant_id=$1 and b.organization_id=$2 and b.item_id=$3 and b.quantity-b.reserved_quantity>0 and w.bin_type in ('pick','putpick') and not w.cross_dock and not w.movement_blocked and ($4 or not p.dedicated) and coalesce(l.blocked,false)=false and (l.expiration_date is null or l.expiration_date>=current_date)`
		if tracking == "lot" && value.UseFEFO {
			query += ` order by l.expiration_date nulls last,l.lot_no,l.created_at,w.bin_rank desc,w.bin_code`
		} else {
			query += ` order by w.bin_rank desc,w.bin_code,l.expiration_date nulls last,l.lot_no`
		}
		query += ` for update of b`
		rows, queryErr := tx.Query(ctx, query, tenant, value.OrganizationID, value.ItemID, value.AllowDedicated)
		if queryErr != nil {
			return activity, queryErr
		}
		if err = appendCandidates(rows); err != nil {
			return activity, err
		}
	}
	if remaining.Sign() != 0 || len(candidates) == 0 {
		return activity, fmt.Errorf("warehouse pick insufficient FEFO availability: %w", inventorycontrol.ErrConflict)
	}
	_, err = tx.Exec(ctx, `insert into inventory.warehouse_activity(tenant_id,activity_id,organization_id,activity_type,status,source_kind,source_id,source_line_id,request_id,assigned_to,version) values($1,$2,$3,'pick','open',$4,$5,$6,$7,$8,1)`, tenant, ids.ActivityID, value.OrganizationID, value.DemandKind, value.DemandID, value.DemandLineID, value.RequestID, value.AssignedTo)
	if err != nil {
		return activity, bulkConflict(err)
	}
	_, err = tx.Exec(ctx, `insert into inventory.warehouse_pick_request(tenant_id,request_id,activity_id,organization_id,ship_bin_id,item_id,demand_kind,demand_id,demand_line_id,quantity,requested_handling_uom,resolved_handling_uom,allow_breakbulk,use_fefo,allow_dedicated,assigned_to) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10::numeric,$11,$12,$13,$14,$15,$16)`, tenant, value.RequestID, ids.ActivityID, value.OrganizationID, value.ShipBinID, value.ItemID, value.DemandKind, value.DemandID, value.DemandLineID, value.Quantity, value.HandlingUOM, targetUOM.code, value.AllowBreakbulk, value.UseFEFO, value.AllowDedicated, value.AssignedTo)
	if err != nil {
		return activity, bulkConflict(err)
	}
	for index, c := range candidates {
		lineID := fmt.Sprintf("%s-%06d", ids.ActivityID, index+1)
		reservationID := fmt.Sprintf("%s-r-%06d", ids.ActivityID, index+1)
		_, err = tx.Exec(ctx, `insert into inventory.bulk_reservation(tenant_id,reservation_id,organization_id,bin_id,item_id,lot_id,demand_kind,demand_id,demand_line_id,quantity,status,cancellation_disallowed,version) values($1,$2,$3,$4,$5,nullif($6,''),$7,$8,$9,$10::numeric,'reservation',true,1)`, tenant, reservationID, value.OrganizationID, c.binID, value.ItemID, c.lotID, value.DemandKind, value.DemandID, fmt.Sprintf("%s:%s:%06d", value.DemandLineID, ids.ActivityID, index+1), c.plan.movedBase)
		if err != nil {
			return activity, bulkConflict(err)
		}
		_, err = tx.Exec(ctx, `insert into inventory.warehouse_activity_line(tenant_id,activity_id,organization_id,line_id,sequence_no,from_bin_id,to_bin_id,item_id,lot_id,reservation_id,quantity,expiration_date,uom_code,uom_quantity,qty_per_uom,from_uom_code,from_uom_quantity,from_qty_per_uom) values($1,$2,$3,$4,$5,$6,$7,$8,nullif($9,''),$10,$11::numeric,$12,$13,$14::numeric,$15::numeric,$16,$17::numeric,$18::numeric)`, tenant, ids.ActivityID, value.OrganizationID, lineID, index+1, c.binID, value.ShipBinID, value.ItemID, c.lotID, reservationID, c.plan.movedBase, c.expiration, c.plan.targetUOM, c.plan.targetQuantity, c.plan.targetFactor, c.plan.sourceUOM, c.plan.sourceQuantity, c.plan.sourceFactor)
		if err != nil {
			return activity, bulkConflict(err)
		}
		line := inventorycontrol.WarehouseActivityLine{ID: lineID, Sequence: index + 1, FromBinID: c.binID, ToBinID: value.ShipBinID, ItemID: value.ItemID, LotID: c.lotID, LotNo: c.lotNo, ReservationID: reservationID, ReservationVersion: 1, Quantity: c.plan.movedBase, FromUOMCode: c.plan.sourceUOM, FromUOMQuantity: c.plan.sourceQuantity, FromQuantityPerUOM: c.plan.sourceFactor, UOMCode: c.plan.targetUOM, UOMQuantity: c.plan.targetQuantity, QuantityPerUOM: c.plan.targetFactor, ExpirationDate: c.expiration}
		if err = reserveWarehousePackaging(ctx, tx, tenant, ids.ActivityID, value.OrganizationID, c.plan, line); err != nil {
			return activity, err
		}
		if c.allocationID != "" {
			updated, err := tx.Exec(ctx, `update inventory.warehouse_crossdock_allocation set reserved_quantity=reserved_quantity+$4::numeric where tenant_id=$1 and allocation_id=$2 and available and transfer_id=$3 and quantity-reserved_quantity-picked_quantity >= $4::numeric`, tenant, c.allocationID, value.DemandID, c.plan.movedBase)
			if err != nil || updated.RowsAffected() != 1 {
				if err != nil {
					return activity, err
				}
				return activity, inventorycontrol.ErrConflict
			}
			_, err = tx.Exec(ctx, `insert into inventory.warehouse_crossdock_pick_link(tenant_id,pick_activity_id,pick_line_id,allocation_id,quantity,status) values($1,$2,$3,$4,$5::numeric,'open')`, tenant, ids.ActivityID, lineID, c.allocationID, c.plan.movedBase)
			if err != nil {
				return activity, bulkConflict(err)
			}
		}
		activity.Lines = append(activity.Lines, line)
	}
	if err = recordBulkEvent(ctx, tx, tenant, ids.EventID, "warehouse-pick", ids.ActivityID, "warehouse-pick.created", 1, map[string]any{"organization_id": value.OrganizationID, "item_id": value.ItemID, "quantity": value.Quantity, "handling_uom": targetUOM.code, "allow_breakbulk": value.AllowBreakbulk, "use_fefo": value.UseFEFO, "lines": len(activity.Lines), "demand_sales_uom": demandSalesUOM, "demand_total_base": demandTotalBase, "demand_outstanding_before_base": demandOutstandingBase}); err != nil {
		return activity, err
	}
	return activity, tx.Commit(ctx)
}

func readWarehouseActivity(ctx context.Context, tx pgx.Tx, tenant, organization, activityID string) (inventorycontrol.WarehouseActivity, error) {
	var value inventorycontrol.WarehouseActivity
	err := tx.QueryRow(ctx, `select activity_id,organization_id,activity_type,status,source_kind,source_id,source_line_id,request_id,assigned_to,version from inventory.warehouse_activity where tenant_id=$1 and organization_id=$2 and activity_id=$3`, tenant, organization, activityID).Scan(&value.ID, &value.OrganizationID, &value.Type, &value.Status, &value.SourceKind, &value.SourceID, &value.SourceLineID, &value.RequestID, &value.AssignedTo, &value.Version)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return value, inventorycontrol.ErrConflict
		}
		return value, err
	}
	rows, err := tx.Query(ctx, `select al.line_id,al.sequence_no,al.from_bin_id,al.to_bin_id,al.item_id,coalesce(al.lot_id,''),coalesce(l.lot_no,''),coalesce(al.reservation_id,''),al.quantity::text,al.from_uom_code,al.from_uom_quantity::text,al.from_qty_per_uom::text,al.uom_code,al.uom_quantity::text,al.qty_per_uom::text,al.expiration_date,coalesce(al.source_entry_id,'') from inventory.warehouse_activity_line al left join inventory.inventory_lot l on l.tenant_id=al.tenant_id and l.lot_id=al.lot_id where al.tenant_id=$1 and al.activity_id=$2 order by al.sequence_no`, tenant, activityID)
	if err != nil {
		return value, err
	}
	defer rows.Close()
	value.Lines = []inventorycontrol.WarehouseActivityLine{}
	for rows.Next() {
		var line inventorycontrol.WarehouseActivityLine
		if err = rows.Scan(&line.ID, &line.Sequence, &line.FromBinID, &line.ToBinID, &line.ItemID, &line.LotID, &line.LotNo, &line.ReservationID, &line.Quantity, &line.FromUOMCode, &line.FromUOMQuantity, &line.FromQuantityPerUOM, &line.UOMCode, &line.UOMQuantity, &line.QuantityPerUOM, &line.ExpirationDate, &line.SourceEntryID); err != nil {
			return value, err
		}
		value.Lines = append(value.Lines, line)
	}
	return value, rows.Err()
}

func releasePutAwaySourceComposition(ctx context.Context, tx pgx.Tx, tenant, organization string, line inventorycontrol.WarehouseActivityLine) error {
	var balanceID int64
	var isBase bool
	err := tx.QueryRow(ctx, `select b.balance_id,u.is_base from inventory.bulk_balance b join inventory.item_unit_of_measure u on u.tenant_id=b.tenant_id and u.item_id=b.item_id and u.uom_code=$7 where b.tenant_id=$1 and b.organization_id=$2 and b.bin_id=$3 and b.item_id=$4 and b.lot_id is not distinct from nullif($5,'') and b.quantity >= $6::numeric and b.reserved_quantity >= $6::numeric for update of b,u`, tenant, organization, line.FromBinID, line.ItemID, line.LotID, line.Quantity, line.UOMCode).Scan(&balanceID, &isBase)
	if errors.Is(err, pgx.ErrNoRows) {
		return inventorycontrol.ErrConflict
	}
	if err != nil {
		return err
	}
	if isBase {
		updated, updateErr := tx.Exec(ctx, `update inventory.bulk_uom_balance set reserved_quantity=reserved_quantity-$3::numeric,version=version+1,updated_at=clock_timestamp() where balance_id=$1 and uom_code=$2 and reserved_quantity >= $3::numeric`, balanceID, line.UOMCode, line.UOMQuantity)
		if updateErr != nil {
			return bulkConflict(updateErr)
		}
		if updated.RowsAffected() != 1 {
			return inventorycontrol.ErrConflict
		}
		return nil
	}
	updated, err := tx.Exec(ctx, `update inventory.bulk_uom_balance set quantity=quantity-$3::numeric,reserved_quantity=reserved_quantity-$3::numeric,quantity_base=quantity_base-$4::numeric,version=version+1,updated_at=clock_timestamp() where balance_id=$1 and uom_code=$2 and quantity >= $3::numeric and reserved_quantity >= $3::numeric and quantity_base >= $4::numeric`, balanceID, line.UOMCode, line.UOMQuantity, line.Quantity)
	if err != nil {
		return bulkConflict(err)
	}
	if updated.RowsAffected() != 1 {
		return inventorycontrol.ErrConflict
	}
	if _, err = tx.Exec(ctx, `delete from inventory.bulk_uom_balance where balance_id=$1 and uom_code=$2 and quantity_base=0 and reserved_quantity=0`, balanceID, line.UOMCode); err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `insert into inventory.bulk_uom_balance(balance_id,uom_code,quantity,qty_per_uom,quantity_base) select $1,u.uom_code,$2::numeric/u.qty_per_uom,u.qty_per_uom,$2::numeric from inventory.bulk_balance b join inventory.item_unit_of_measure u on u.tenant_id=b.tenant_id and u.item_id=b.item_id and u.is_base where b.balance_id=$1 on conflict(balance_id,uom_code) do update set quantity=inventory.bulk_uom_balance.quantity+excluded.quantity,quantity_base=inventory.bulk_uom_balance.quantity_base+excluded.quantity_base,version=inventory.bulk_uom_balance.version+1,updated_at=clock_timestamp()`, balanceID, line.Quantity)
	return bulkConflict(err)
}

func preservePutAwayTargetComposition(ctx context.Context, tx pgx.Tx, balanceID int64, line inventorycontrol.WarehouseActivityLine) error {
	var isBase bool
	err := tx.QueryRow(ctx, `select is_base from inventory.item_unit_of_measure u join inventory.bulk_balance b on b.tenant_id=u.tenant_id and b.item_id=u.item_id where b.balance_id=$1 and u.uom_code=$2`, balanceID, line.UOMCode).Scan(&isBase)
	if err != nil || isBase {
		return err
	}
	updated, err := tx.Exec(ctx, `update inventory.bulk_uom_balance target set quantity=target.quantity-($2::numeric/target.qty_per_uom),quantity_base=target.quantity_base-$2::numeric,version=target.version+1,updated_at=clock_timestamp() where target.balance_id=$1 and target.uom_code=(select u.uom_code from inventory.bulk_balance b join inventory.item_unit_of_measure u on u.tenant_id=b.tenant_id and u.item_id=b.item_id and u.is_base where b.balance_id=$1) and target.quantity-target.reserved_quantity >= ($2::numeric/target.qty_per_uom) and target.quantity_base >= $2::numeric`, balanceID, line.Quantity)
	if err != nil {
		return bulkConflict(err)
	}
	if updated.RowsAffected() != 1 {
		return inventorycontrol.ErrConflict
	}
	_, err = tx.Exec(ctx, `delete from inventory.bulk_uom_balance target where target.balance_id=$1 and target.quantity_base=0 and target.reserved_quantity=0 and target.uom_code=(select u.uom_code from inventory.bulk_balance b join inventory.item_unit_of_measure u on u.tenant_id=b.tenant_id and u.item_id=b.item_id and u.is_base where b.balance_id=$1)`, balanceID)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `insert into inventory.bulk_uom_balance(balance_id,uom_code,quantity,qty_per_uom,quantity_base) values($1,$2,$3::numeric,$4::numeric,$5::numeric) on conflict(balance_id,uom_code) do update set quantity=inventory.bulk_uom_balance.quantity+excluded.quantity,quantity_base=inventory.bulk_uom_balance.quantity_base+excluded.quantity_base,version=inventory.bulk_uom_balance.version+1,updated_at=clock_timestamp()`, balanceID, line.UOMCode, line.UOMQuantity, line.QuantityPerUOM, line.Quantity)
	return bulkConflict(err)
}

func (r *InventoryControl) RegisterWarehouseActivity(ctx context.Context, tenant, organization, activityID string, version int64, postingDate time.Time, eventID string) (inventorycontrol.WarehouseActivity, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return inventorycontrol.WarehouseActivity{}, err
	}
	defer tx.Rollback(ctx)
	var activityType, itemID string
	err = tx.QueryRow(ctx, `select a.activity_type,al.item_id from inventory.warehouse_activity a join inventory.warehouse_activity_line al using(tenant_id,activity_id) where a.tenant_id=$1 and a.organization_id=$2 and a.activity_id=$3 and a.status='open' and a.version=$4 order by al.sequence_no limit 1 for update of a`, tenant, organization, activityID, version).Scan(&activityType, &itemID)
	if errors.Is(err, pgx.ErrNoRows) {
		return inventorycontrol.WarehouseActivity{}, inventorycontrol.ErrConflict
	}
	if err != nil {
		return inventorycontrol.WarehouseActivity{}, err
	}
	if err = warehouseLock(ctx, tx, tenant, organization, itemID); err != nil {
		return inventorycontrol.WarehouseActivity{}, err
	}
	activity, err := readWarehouseActivity(ctx, tx, tenant, organization, activityID)
	if err != nil {
		return activity, err
	}
	for index, line := range activity.Lines {
		if activityType == "put-away" {
			if err = releasePutAwaySourceComposition(ctx, tx, tenant, organization, line); err != nil {
				return activity, err
			}
			updated, updateErr := tx.Exec(ctx, `update inventory.bulk_balance set quantity=quantity-$6::numeric,reserved_quantity=reserved_quantity-$6::numeric,version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and organization_id=$2 and bin_id=$3 and item_id=$4 and lot_id is not distinct from nullif($5,'') and quantity >= $6::numeric and reserved_quantity >= $6::numeric`, tenant, organization, line.FromBinID, line.ItemID, line.LotID, line.Quantity)
			if updateErr != nil {
				return activity, updateErr
			}
			if updated.RowsAffected() != 1 {
				return activity, inventorycontrol.ErrConflict
			}
			var targetBalanceID, targetVersion int64
			err = tx.QueryRow(ctx, `insert into inventory.bulk_balance(tenant_id,organization_id,bin_id,item_id,lot_id,quantity,reserved_quantity,version) values($1,$2,$3,$4,nullif($5,''),$6::numeric,0,1) on conflict (tenant_id,organization_id,bin_id,item_id,lot_id) do update set quantity=inventory.bulk_balance.quantity+excluded.quantity,version=inventory.bulk_balance.version+1,updated_at=clock_timestamp() where exists(select 1 from inventory.item_bin_policy p join inventory.warehouse_bin w using(tenant_id,organization_id,bin_id) where p.tenant_id=$1 and p.organization_id=$2 and p.item_id=$4 and p.bin_id=$3 and w.bin_type in ('put-away','putpick') and not w.movement_blocked and (p.max_quantity is null or inventory.bulk_balance.quantity+excluded.quantity<=p.max_quantity)) returning balance_id,version`, tenant, organization, line.ToBinID, line.ItemID, line.LotID, line.Quantity).Scan(&targetBalanceID, &targetVersion)
			if errors.Is(err, pgx.ErrNoRows) {
				return activity, inventorycontrol.ErrConflict
			}
			if err != nil {
				return activity, bulkConflict(err)
			}
			if err = preservePutAwayTargetComposition(ctx, tx, targetBalanceID, line); err != nil {
				return activity, err
			}
			_, err = tx.Exec(ctx, `update inventory.warehouse_crossdock_allocation set available=true where tenant_id=$1 and put_away_activity_id=$2 and put_away_line_id=$3 and not available`, tenant, activityID, line.ID)
			if err != nil {
				return activity, err
			}
		} else if activityType == "pick" {
			if err = moveWarehousePackagingSource(ctx, tx, tenant, organization, activityID, "warehouse-pick", line); err != nil {
				return activity, err
			}
			activity.Lines[index].ReservationVersion = 2
			var targetBalanceID, targetVersion int64
			err = tx.QueryRow(ctx, `insert into inventory.bulk_balance(tenant_id,organization_id,bin_id,item_id,lot_id,quantity,reserved_quantity,version) values($1,$2,$3,$4,nullif($5,''),$6::numeric,$6::numeric,1) on conflict (tenant_id,organization_id,bin_id,item_id,lot_id) do update set quantity=inventory.bulk_balance.quantity+excluded.quantity,reserved_quantity=inventory.bulk_balance.reserved_quantity+excluded.reserved_quantity,version=inventory.bulk_balance.version+1,updated_at=clock_timestamp() where exists(select 1 from inventory.item_bin_policy p join inventory.warehouse_bin w using(tenant_id,organization_id,bin_id) where p.tenant_id=$1 and p.organization_id=$2 and p.item_id=$4 and p.bin_id=$3 and w.bin_type='ship' and not w.movement_blocked and (p.max_quantity is null or inventory.bulk_balance.quantity+excluded.quantity<=p.max_quantity)) returning balance_id,version`, tenant, organization, line.ToBinID, line.ItemID, line.LotID, line.Quantity).Scan(&targetBalanceID, &targetVersion)
			if errors.Is(err, pgx.ErrNoRows) {
				return activity, inventorycontrol.ErrConflict
			}
			if err != nil {
				return activity, bulkConflict(err)
			}
			if err = preservePutAwayTargetComposition(ctx, tx, targetBalanceID, line); err != nil {
				return activity, err
			}
			updated, err := tx.Exec(ctx, `update inventory.bulk_reservation set bin_id=$4,version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and organization_id=$2 and reservation_id=$3 and status='reservation' and version=1`, tenant, organization, line.ReservationID, line.ToBinID)
			if err != nil {
				return activity, err
			}
			if updated.RowsAffected() != 1 {
				return activity, inventorycontrol.ErrConflict
			}
			var allocationID string
			linkErr := tx.QueryRow(ctx, `select allocation_id from inventory.warehouse_crossdock_pick_link where tenant_id=$1 and pick_activity_id=$2 and pick_line_id=$3 and status='open' for update`, tenant, activityID, line.ID).Scan(&allocationID)
			if linkErr == nil {
				updated, err = tx.Exec(ctx, `update inventory.warehouse_crossdock_allocation set reserved_quantity=reserved_quantity-$3::numeric,picked_quantity=picked_quantity+$3::numeric where tenant_id=$1 and allocation_id=$2 and reserved_quantity >= $3::numeric and reserved_quantity+picked_quantity <= quantity`, tenant, allocationID, line.Quantity)
				if err != nil || updated.RowsAffected() != 1 {
					if err != nil {
						return activity, err
					}
					return activity, inventorycontrol.ErrConflict
				}
				updated, err = tx.Exec(ctx, `update inventory.warehouse_crossdock_pick_link set status='picked' where tenant_id=$1 and pick_activity_id=$2 and pick_line_id=$3 and status='open'`, tenant, activityID, line.ID)
				if err != nil || updated.RowsAffected() != 1 {
					if err != nil {
						return activity, err
					}
					return activity, inventorycontrol.ErrConflict
				}
			} else if !errors.Is(linkErr, pgx.ErrNoRows) {
				return activity, linkErr
			}
		} else if activityType == "movement" {
			if err = moveWarehousePackagingSource(ctx, tx, tenant, organization, activityID, "warehouse-replenishment", line); err != nil {
				return activity, err
			}
			var targetBalanceID, targetVersion int64
			err = tx.QueryRow(ctx, `insert into inventory.bulk_balance(tenant_id,organization_id,bin_id,item_id,lot_id,quantity,reserved_quantity,version) values($1,$2,$3,$4,nullif($5,''),$6::numeric,0,1) on conflict (tenant_id,organization_id,bin_id,item_id,lot_id) do update set quantity=inventory.bulk_balance.quantity+excluded.quantity,version=inventory.bulk_balance.version+1,updated_at=clock_timestamp() where exists(select 1 from inventory.item_bin_policy p join inventory.warehouse_bin w using(tenant_id,organization_id,bin_id) where p.tenant_id=$1 and p.organization_id=$2 and p.item_id=$4 and p.bin_id=$3 and p.fixed and w.bin_type in ('pick','putpick') and not w.movement_blocked and p.max_quantity is not null and (select coalesce(sum(quantity),0) from inventory.bulk_balance total where total.tenant_id=$1 and total.organization_id=$2 and total.bin_id=$3 and total.item_id=$4)+excluded.quantity<=p.max_quantity) returning balance_id,version`, tenant, organization, line.ToBinID, line.ItemID, line.LotID, line.Quantity).Scan(&targetBalanceID, &targetVersion)
			if errors.Is(err, pgx.ErrNoRows) {
				return activity, inventorycontrol.ErrConflict
			}
			if err != nil {
				return activity, bulkConflict(err)
			}
			if err = preservePutAwayTargetComposition(ctx, tx, targetBalanceID, line); err != nil {
				return activity, err
			}
			updated, err := tx.Exec(ctx, `update inventory.bulk_reservation set status='consumed',version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and organization_id=$2 and reservation_id=$3 and status='reservation' and version=1`, tenant, organization, line.ReservationID)
			if err != nil || updated.RowsAffected() != 1 {
				if err != nil {
					return activity, err
				}
				return activity, inventorycontrol.ErrConflict
			}
		} else {
			return activity, inventorycontrol.ErrConflict
		}
		outID := fmt.Sprintf("%s-out-%06d", eventID, index+1)
		inID := fmt.Sprintf("%s-in-%06d", eventID, index+1)
		_, err = tx.Exec(ctx, `insert into inventory.bulk_inventory_entry(tenant_id,entry_id,organization_id,bin_id,item_id,lot_id,entry_type,quantity,unit_cost,cost_amount,posting_date,source_kind,source_id) values($1,$2,$3,$4,$5,nullif($6,''),'movement-out',-$7::numeric,0,0,$8::date,$9,$10),($1,$11,$3,$12,$5,nullif($6,''),'movement-in',$7::numeric,0,0,$8::date,$9,$10)`, tenant, outID, organization, line.FromBinID, line.ItemID, line.LotID, line.Quantity, postingDate, "warehouse-"+activityType, line.ID, inID, line.ToBinID)
		if err != nil {
			return activity, bulkConflict(err)
		}
	}
	updated, err := tx.Exec(ctx, `update inventory.warehouse_activity set status='registered',version=version+1,registered_at=clock_timestamp() where tenant_id=$1 and organization_id=$2 and activity_id=$3 and status='open' and version=$4`, tenant, organization, activityID, version)
	if err != nil {
		return activity, err
	}
	if updated.RowsAffected() != 1 {
		return activity, inventorycontrol.ErrConflict
	}
	activity.Status, activity.Version = "registered", version+1
	if err = recordBulkEvent(ctx, tx, tenant, eventID, "warehouse-activity", activityID, "warehouse-activity.registered", activity.Version, map[string]any{"organization_id": organization, "activity_type": activityType, "lines": len(activity.Lines)}); err != nil {
		return activity, err
	}
	return activity, tx.Commit(ctx)
}

func (r *InventoryControl) CancelWarehousePick(ctx context.Context, tenant, organization, activityID string, version int64, eventID string) (inventorycontrol.WarehouseActivity, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return inventorycontrol.WarehouseActivity{}, err
	}
	defer tx.Rollback(ctx)
	var itemID string
	err = tx.QueryRow(ctx, `select al.item_id from inventory.warehouse_activity a join inventory.warehouse_activity_line al using(tenant_id,activity_id) where a.tenant_id=$1 and a.organization_id=$2 and a.activity_id=$3 and a.activity_type='pick' and a.status='open' and a.version=$4 order by al.sequence_no limit 1 for update of a`, tenant, organization, activityID, version).Scan(&itemID)
	if errors.Is(err, pgx.ErrNoRows) {
		return inventorycontrol.WarehouseActivity{}, inventorycontrol.ErrConflict
	}
	if err != nil {
		return inventorycontrol.WarehouseActivity{}, err
	}
	if err = warehouseLock(ctx, tx, tenant, organization, itemID); err != nil {
		return inventorycontrol.WarehouseActivity{}, err
	}
	activity, err := readWarehouseActivity(ctx, tx, tenant, organization, activityID)
	if err != nil {
		return activity, err
	}
	for _, line := range activity.Lines {
		if err = releaseWarehousePackagingReservation(ctx, tx, tenant, organization, activityID, line); err != nil {
			return activity, err
		}
		updated, err := tx.Exec(ctx, `update inventory.bulk_reservation set status='released',cancellation_disallowed=false,version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and organization_id=$2 and reservation_id=$3 and status='reservation' and version=1`, tenant, organization, line.ReservationID)
		if err != nil {
			return activity, err
		}
		if updated.RowsAffected() != 1 {
			return activity, inventorycontrol.ErrConflict
		}
		var allocationID string
		linkErr := tx.QueryRow(ctx, `select allocation_id from inventory.warehouse_crossdock_pick_link where tenant_id=$1 and pick_activity_id=$2 and pick_line_id=$3 and status='open' for update`, tenant, activityID, line.ID).Scan(&allocationID)
		if linkErr == nil {
			updated, err = tx.Exec(ctx, `update inventory.warehouse_crossdock_allocation set reserved_quantity=reserved_quantity-$3::numeric where tenant_id=$1 and allocation_id=$2 and reserved_quantity >= $3::numeric`, tenant, allocationID, line.Quantity)
			if err != nil || updated.RowsAffected() != 1 {
				if err != nil {
					return activity, err
				}
				return activity, inventorycontrol.ErrConflict
			}
			updated, err = tx.Exec(ctx, `update inventory.warehouse_crossdock_pick_link set status='released' where tenant_id=$1 and pick_activity_id=$2 and pick_line_id=$3 and status='open'`, tenant, activityID, line.ID)
			if err != nil || updated.RowsAffected() != 1 {
				if err != nil {
					return activity, err
				}
				return activity, inventorycontrol.ErrConflict
			}
		} else if !errors.Is(linkErr, pgx.ErrNoRows) {
			return activity, linkErr
		}
	}
	updated, err := tx.Exec(ctx, `update inventory.warehouse_activity set status='cancelled',version=version+1 where tenant_id=$1 and organization_id=$2 and activity_id=$3 and status='open' and version=$4`, tenant, organization, activityID, version)
	if err != nil {
		return activity, err
	}
	if updated.RowsAffected() != 1 {
		return activity, inventorycontrol.ErrConflict
	}
	activity.Status, activity.Version = "cancelled", version+1
	if err = recordBulkEvent(ctx, tx, tenant, eventID, "warehouse-pick", activityID, "warehouse-pick.cancelled", activity.Version, map[string]any{"organization_id": organization, "lines": len(activity.Lines)}); err != nil {
		return activity, err
	}
	return activity, tx.Commit(ctx)
}

func (r *InventoryControl) CancelWarehousePutAway(ctx context.Context, tenant, organization, activityID string, version int64, eventID string) (inventorycontrol.WarehouseActivity, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return inventorycontrol.WarehouseActivity{}, err
	}
	defer tx.Rollback(ctx)
	var itemID string
	err = tx.QueryRow(ctx, `select al.item_id from inventory.warehouse_activity a join inventory.warehouse_activity_line al using(tenant_id,activity_id) where a.tenant_id=$1 and a.organization_id=$2 and a.activity_id=$3 and a.activity_type='put-away' and a.status='open' and a.version=$4 order by al.sequence_no limit 1 for update of a`, tenant, organization, activityID, version).Scan(&itemID)
	if errors.Is(err, pgx.ErrNoRows) {
		return inventorycontrol.WarehouseActivity{}, inventorycontrol.ErrConflict
	}
	if err != nil {
		return inventorycontrol.WarehouseActivity{}, err
	}
	if err = warehouseLock(ctx, tx, tenant, organization, itemID); err != nil {
		return inventorycontrol.WarehouseActivity{}, err
	}
	activity, err := readWarehouseActivity(ctx, tx, tenant, organization, activityID)
	if err != nil {
		return activity, err
	}
	for _, line := range activity.Lines {
		var balanceID int64
		err = tx.QueryRow(ctx, `update inventory.bulk_balance set reserved_quantity=reserved_quantity-$6::numeric,version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and organization_id=$2 and bin_id=$3 and item_id=$4 and lot_id is not distinct from nullif($5,'') and reserved_quantity >= $6::numeric returning balance_id`, tenant, organization, line.FromBinID, line.ItemID, line.LotID, line.Quantity).Scan(&balanceID)
		if errors.Is(err, pgx.ErrNoRows) {
			return activity, inventorycontrol.ErrConflict
		}
		if err != nil {
			return activity, err
		}
		updated, updateErr := tx.Exec(ctx, `update inventory.bulk_uom_balance set reserved_quantity=reserved_quantity-$3::numeric,version=version+1,updated_at=clock_timestamp() where balance_id=$1 and uom_code=$2 and reserved_quantity >= $3::numeric`, balanceID, line.UOMCode, line.UOMQuantity)
		if updateErr != nil {
			return activity, bulkConflict(updateErr)
		}
		if updated.RowsAffected() != 1 {
			return activity, inventorycontrol.ErrConflict
		}
	}
	updated, err := tx.Exec(ctx, `update inventory.warehouse_activity set status='cancelled',version=version+1 where tenant_id=$1 and organization_id=$2 and activity_id=$3 and activity_type='put-away' and status='open' and version=$4`, tenant, organization, activityID, version)
	if err != nil {
		return activity, err
	}
	if updated.RowsAffected() != 1 {
		return activity, inventorycontrol.ErrConflict
	}
	activity.Status, activity.Version = "cancelled", version+1
	if err = recordBulkEvent(ctx, tx, tenant, eventID, "warehouse-put-away", activityID, "warehouse-put-away.cancelled", activity.Version, map[string]any{"organization_id": organization, "lines": len(activity.Lines)}); err != nil {
		return activity, err
	}
	return activity, tx.Commit(ctx)
}
