package postgres

import (
	"context"
	"errors"

	"elite.local/enterprise/internal/inventorycontrol"
	"github.com/jackc/pgx/v5"
)

func loadPackagingConversion(ctx context.Context, tx pgx.Tx, tenant, request string) (inventorycontrol.PackagingConversionResult, bool, error) {
	var value inventorycontrol.PackagingConversionResult
	err := tx.QueryRow(ctx, `select conversion_id,request_id,organization_id,bin_id,item_id,coalesce(lot_id,''),operation,from_uom,to_uom,from_quantity::text,to_quantity::text,base_quantity::text,from_qty_per_uom::text,to_qty_per_uom::text,version from inventory.bulk_uom_conversion where tenant_id=$1 and request_id=$2`, tenant, request).Scan(&value.ID, &value.RequestID, &value.OrganizationID, &value.BinID, &value.ItemID, &value.LotID, &value.Operation, &value.FromUOM, &value.ToUOM, &value.FromQuantity, &value.ToQuantity, &value.BaseQuantity, &value.FromQuantityPer, &value.ToQuantityPer, &value.Version)
	if errors.Is(err, pgx.ErrNoRows) {
		return value, false, nil
	}
	if err != nil {
		return value, false, err
	}
	rows, err := tx.Query(ctx, `select sequence_no,action_type,uom_code,quantity::text,quantity_base::text from inventory.bulk_uom_conversion_line where tenant_id=$1 and conversion_id=$2 order by sequence_no`, tenant, value.ID)
	if err != nil {
		return value, false, err
	}
	defer rows.Close()
	for rows.Next() {
		var line inventorycontrol.PackagingConversionLine
		if err = rows.Scan(&line.Sequence, &line.Action, &line.UOM, &line.Quantity, &line.BaseQuantity); err != nil {
			return value, false, err
		}
		value.Lines = append(value.Lines, line)
	}
	return value, true, rows.Err()
}

func samePackagingRequest(value inventorycontrol.PackagingConversionResult, command inventorycontrol.PackagingConversionCommand) bool {
	return value.RequestID == command.RequestID && value.OrganizationID == command.OrganizationID && value.BinID == command.BinID && value.ItemID == command.ItemID && value.LotID == command.LotID && value.Operation == command.Operation && value.FromUOM == command.FromUOM && value.ToUOM == command.ToUOM
}

func (r *InventoryControl) ConvertBulkHandlingUnits(ctx context.Context, tenant string, ids inventorycontrol.PackagingConversionIDs, command inventorycontrol.PackagingConversionCommand) (inventorycontrol.PackagingConversionResult, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return inventorycontrol.PackagingConversionResult{}, err
	}
	defer tx.Rollback(ctx)
	if _, err = tx.Exec(ctx, `select pg_advisory_xact_lock(hashtextextended($1,0))`, tenant+":"+command.RequestID); err != nil {
		return inventorycontrol.PackagingConversionResult{}, err
	}
	if existing, found, loadErr := loadPackagingConversion(ctx, tx, tenant, command.RequestID); loadErr != nil {
		return existing, loadErr
	} else if found {
		var quantityMatches bool
		if err = tx.QueryRow(ctx, `select from_quantity=$3::numeric from inventory.bulk_uom_conversion where tenant_id=$1 and request_id=$2`, tenant, command.RequestID, command.FromQuantity).Scan(&quantityMatches); err != nil {
			return existing, err
		}
		if !samePackagingRequest(existing, command) || !quantityMatches {
			return inventorycontrol.PackagingConversionResult{}, inventorycontrol.ErrConflict
		}
		return existing, tx.Commit(ctx)
	}

	var balanceID int64
	var fromFactor, toFactor, baseQuantity, toQuantity string
	err = tx.QueryRow(ctx, `select b.balance_id,fu.qty_per_uom::text,tu.qty_per_uom::text,($9::numeric*fu.qty_per_uom)::numeric(20,6)::text,(($9::numeric*fu.qty_per_uom)/tu.qty_per_uom)::numeric(20,6)::text
from inventory.bulk_balance b
join inventory.item_unit_of_measure fu on fu.tenant_id=b.tenant_id and fu.item_id=b.item_id and fu.uom_code=$7
join inventory.item_unit_of_measure tu on tu.tenant_id=b.tenant_id and tu.item_id=b.item_id and tu.uom_code=$8
where b.tenant_id=$1 and b.organization_id=$2 and b.bin_id=$3 and b.item_id=$4 and b.lot_id is not distinct from nullif($5,'')
  and (($6='breakbulk' and fu.qty_per_uom>tu.qty_per_uom) or ($6='gather' and fu.qty_per_uom<tu.qty_per_uom))
  and mod($9::numeric,fu.rounding_precision)=0
  and mod(($9::numeric*fu.qty_per_uom)/tu.qty_per_uom,tu.rounding_precision)=0
  and b.quantity-b.reserved_quantity >= $9::numeric*fu.qty_per_uom
for update of b`, tenant, command.OrganizationID, command.BinID, command.ItemID, command.LotID, command.Operation, command.FromUOM, command.ToUOM, command.FromQuantity).Scan(&balanceID, &fromFactor, &toFactor, &baseQuantity, &toQuantity)
	if errors.Is(err, pgx.ErrNoRows) {
		return inventorycontrol.PackagingConversionResult{}, inventorycontrol.ErrConflict
	}
	if err != nil {
		return inventorycontrol.PackagingConversionResult{}, bulkConflict(err)
	}

	result, err := tx.Exec(ctx, `update inventory.bulk_uom_balance set quantity=quantity-$3::numeric,quantity_base=quantity_base-$4::numeric,version=version+1,updated_at=clock_timestamp() where balance_id=$1 and uom_code=$2 and quantity-reserved_quantity >= $3::numeric and quantity_base>=$4::numeric`, balanceID, command.FromUOM, command.FromQuantity, baseQuantity)
	if err != nil {
		return inventorycontrol.PackagingConversionResult{}, bulkConflict(err)
	}
	if result.RowsAffected() != 1 {
		return inventorycontrol.PackagingConversionResult{}, inventorycontrol.ErrConflict
	}
	if _, err = tx.Exec(ctx, `delete from inventory.bulk_uom_balance where balance_id=$1 and uom_code=$2 and quantity_base=0 and reserved_quantity=0`, balanceID, command.FromUOM); err != nil {
		return inventorycontrol.PackagingConversionResult{}, err
	}
	_, err = tx.Exec(ctx, `insert into inventory.bulk_uom_balance(balance_id,uom_code,quantity,qty_per_uom,quantity_base) values($1,$2,$3::numeric,$4::numeric,$5::numeric) on conflict(balance_id,uom_code) do update set quantity=inventory.bulk_uom_balance.quantity+excluded.quantity,quantity_base=inventory.bulk_uom_balance.quantity_base+excluded.quantity_base,version=inventory.bulk_uom_balance.version+1,updated_at=clock_timestamp()`, balanceID, command.ToUOM, toQuantity, toFactor, baseQuantity)
	if err != nil {
		return inventorycontrol.PackagingConversionResult{}, bulkConflict(err)
	}
	_, err = tx.Exec(ctx, `insert into inventory.bulk_uom_conversion(tenant_id,conversion_id,request_id,organization_id,bin_id,item_id,lot_id,operation,from_uom,to_uom,from_quantity,to_quantity,base_quantity,from_qty_per_uom,to_qty_per_uom) values($1,$2,$3,$4,$5,$6,nullif($7,''),$8,$9,$10,$11::numeric,$12::numeric,$13::numeric,$14::numeric,$15::numeric)`, tenant, ids.ConversionID, command.RequestID, command.OrganizationID, command.BinID, command.ItemID, command.LotID, command.Operation, command.FromUOM, command.ToUOM, command.FromQuantity, toQuantity, baseQuantity, fromFactor, toFactor)
	if err != nil {
		return inventorycontrol.PackagingConversionResult{}, bulkConflict(err)
	}
	_, err = tx.Exec(ctx, `insert into inventory.bulk_uom_conversion_line(tenant_id,conversion_id,line_id,sequence_no,action_type,uom_code,quantity,quantity_base) values($1,$2,$3,1,'take',$5,$6::numeric,$7::numeric),($1,$2,$4,2,'place',$8,$9::numeric,$7::numeric)`, tenant, ids.ConversionID, ids.TakeLineID, ids.PlaceLineID, command.FromUOM, command.FromQuantity, baseQuantity, command.ToUOM, toQuantity)
	if err != nil {
		return inventorycontrol.PackagingConversionResult{}, bulkConflict(err)
	}
	if err = recordBulkEvent(ctx, tx, tenant, ids.EventID, "bulk-uom-conversion", ids.ConversionID, "bulk-uom-conversion.completed", 1, map[string]string{"request_id": command.RequestID, "organization_id": command.OrganizationID, "bin_id": command.BinID, "item_id": command.ItemID, "lot_id": command.LotID, "operation": command.Operation, "from_uom": command.FromUOM, "to_uom": command.ToUOM, "from_quantity": command.FromQuantity, "to_quantity": toQuantity, "base_quantity": baseQuantity}); err != nil {
		return inventorycontrol.PackagingConversionResult{}, err
	}
	value, found, err := loadPackagingConversion(ctx, tx, tenant, command.RequestID)
	if err != nil || !found {
		return value, err
	}
	return value, tx.Commit(ctx)
}
