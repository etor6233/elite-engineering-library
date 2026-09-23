package postgres

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"math/big"

	"elite.local/enterprise/internal/inventorycontrol"
	"github.com/jackc/pgx/v5"
)

type warehouseUOMContract struct {
	code      string
	factor    string
	precision string
}

type warehouseUOMPlan struct {
	balanceID          int64
	sourceUOM          string
	sourceQuantity     string
	sourceFactor       string
	sourceBase         string
	targetUOM          string
	targetQuantity     string
	targetFactor       string
	movedBase          string
	requiresConversion bool
}

func resolveWarehouseUOM(ctx context.Context, tx pgx.Tx, tenant, item, requested string) (warehouseUOMContract, error) {
	var value warehouseUOMContract
	err := tx.QueryRow(ctx, `select u.uom_code,u.qty_per_uom::text,u.rounding_precision::text from inventory.item_unit_of_measure u where u.tenant_id=$1 and u.item_id=$2 and u.uom_code=coalesce(nullif($3,''),(select b.uom_code from inventory.item_unit_of_measure b where b.tenant_id=$1 and b.item_id=$2 and b.is_base)) for share of u`, tenant, item, requested).Scan(&value.code, &value.factor, &value.precision)
	if errors.Is(err, pgx.ErrNoRows) {
		return value, inventorycontrol.ErrConflict
	}
	return value, err
}

func floorToQuantum(value, quantum *big.Rat) *big.Rat {
	if value.Sign() <= 0 || quantum.Sign() <= 0 {
		return new(big.Rat)
	}
	steps := new(big.Rat).Quo(value, quantum)
	whole := new(big.Int).Quo(steps.Num(), steps.Denom())
	return new(big.Rat).Mul(new(big.Rat).SetInt(whole), quantum)
}

func ceilToQuantum(value, quantum *big.Rat) *big.Rat {
	if value.Sign() <= 0 || quantum.Sign() <= 0 {
		return new(big.Rat)
	}
	steps := new(big.Rat).Quo(value, quantum)
	whole, remainder := new(big.Int), new(big.Int)
	whole.QuoRem(steps.Num(), steps.Denom(), remainder)
	if remainder.Sign() > 0 {
		whole.Add(whole, big.NewInt(1))
	}
	return new(big.Rat).Mul(new(big.Rat).SetInt(whole), quantum)
}

func planWarehousePackaging(ctx context.Context, tx pgx.Tx, balanceID int64, target warehouseUOMContract, allowBreakbulk bool, maximumBase *big.Rat) ([]warehouseUOMPlan, *big.Rat, error) {
	targetFactor, targetFactorOK := parseRat(target.factor)
	targetPrecision, targetPrecisionOK := parseRat(target.precision)
	if !targetFactorOK || !targetPrecisionOK {
		return nil, new(big.Rat), inventorycontrol.ErrConflict
	}
	targetQuantum := new(big.Rat).Mul(targetFactor, targetPrecision)
	remaining := new(big.Rat).Set(maximumBase)
	rows, err := tx.Query(ctx, `select c.uom_code,c.qty_per_uom::text,u.rounding_precision::text,(c.quantity-c.reserved_quantity)::text from inventory.bulk_uom_balance c join inventory.item_unit_of_measure u on u.tenant_id=(select b.tenant_id from inventory.bulk_balance b where b.balance_id=c.balance_id) and u.item_id=(select b.item_id from inventory.bulk_balance b where b.balance_id=c.balance_id) and u.uom_code=c.uom_code where c.balance_id=$1 and c.quantity>c.reserved_quantity and (c.uom_code=$2 or ($3 and c.qty_per_uom>$4::numeric)) order by (c.uom_code=$2) desc,c.qty_per_uom,c.uom_code for update of c`, balanceID, target.code, allowBreakbulk, target.factor)
	if err != nil {
		return nil, remaining, err
	}
	type composition struct{ code, factor, precision, available string }
	available := []composition{}
	for rows.Next() {
		var value composition
		if err = rows.Scan(&value.code, &value.factor, &value.precision, &value.available); err != nil {
			rows.Close()
			return nil, remaining, err
		}
		available = append(available, value)
	}
	rows.Close()
	if err = rows.Err(); err != nil {
		return nil, remaining, err
	}
	plans := []warehouseUOMPlan{}
	for _, source := range available {
		if remaining.Sign() <= 0 {
			break
		}
		sourceFactor, factorOK := parseRat(source.factor)
		sourcePrecision, precisionOK := parseRat(source.precision)
		sourceAvailable, availableOK := parseRat(source.available)
		if !factorOK || !precisionOK || !availableOK {
			return nil, remaining, inventorycontrol.ErrConflict
		}
		sourceAvailableBase := new(big.Rat).Mul(sourceAvailable, sourceFactor)
		move := floorToQuantum(minRat(remaining, sourceAvailableBase), targetQuantum)
		if move.Sign() <= 0 {
			continue
		}
		sourceQuantity := new(big.Rat).Quo(move, sourceFactor)
		if source.code != target.code {
			sourceQuantity = ceilToQuantum(sourceQuantity, sourcePrecision)
		}
		if sourceQuantity.Cmp(sourceAvailable) > 0 {
			continue
		}
		sourceBase := new(big.Rat).Mul(sourceQuantity, sourceFactor)
		targetQuantity := new(big.Rat).Quo(move, targetFactor)
		plans = append(plans, warehouseUOMPlan{
			balanceID: balanceID, sourceUOM: source.code, sourceQuantity: formatRat(sourceQuantity, 6),
			sourceFactor: source.factor, sourceBase: formatRat(sourceBase, 6), targetUOM: target.code,
			targetQuantity: formatRat(targetQuantity, 6), targetFactor: target.factor,
			movedBase: formatRat(move, 6), requiresConversion: source.code != target.code,
		})
		remaining.Sub(remaining, move)
	}
	return plans, remaining, nil
}

func reserveWarehousePackaging(ctx context.Context, tx pgx.Tx, tenant, activityID, organization string, plan warehouseUOMPlan, line inventorycontrol.WarehouseActivityLine) error {
	updated, err := tx.Exec(ctx, `update inventory.bulk_balance set reserved_quantity=reserved_quantity+$2::numeric,version=version+1,updated_at=clock_timestamp() where balance_id=$1 and quantity-reserved_quantity >= $2::numeric`, plan.balanceID, plan.movedBase)
	if err != nil || updated.RowsAffected() != 1 {
		if err != nil {
			return bulkConflict(err)
		}
		return inventorycontrol.ErrConflict
	}
	updated, err = tx.Exec(ctx, `update inventory.bulk_uom_balance set reserved_quantity=reserved_quantity+$3::numeric,version=version+1,updated_at=clock_timestamp() where balance_id=$1 and uom_code=$2 and quantity-reserved_quantity >= $3::numeric`, plan.balanceID, plan.sourceUOM, plan.sourceQuantity)
	if err != nil || updated.RowsAffected() != 1 {
		if err != nil {
			return bulkConflict(err)
		}
		return inventorycontrol.ErrConflict
	}
	_, err = tx.Exec(ctx, `insert into inventory.warehouse_activity_uom_reservation(tenant_id,activity_id,line_id,organization_id,source_balance_id,source_uom_code,source_uom_quantity,source_qty_per_uom,source_base_quantity,target_uom_code,target_uom_quantity,target_qty_per_uom,moved_base_quantity,conversion_required) values($1,$2,$3,$4,$5,$6,$7::numeric,$8::numeric,$9::numeric,$10,$11::numeric,$12::numeric,$13::numeric,$14)`, tenant, activityID, line.ID, organization, plan.balanceID, plan.sourceUOM, plan.sourceQuantity, plan.sourceFactor, plan.sourceBase, plan.targetUOM, plan.targetQuantity, plan.targetFactor, plan.movedBase, plan.requiresConversion)
	return bulkConflict(err)
}

func deterministicWarehouseUUID(parts ...string) string {
	sum := sha256.Sum256([]byte(fmt.Sprint(parts)))
	sum[6] = (sum[6] & 0x0f) | 0x50
	sum[8] = (sum[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", sum[0:4], sum[4:6], sum[6:8], sum[8:10], sum[10:16])
}

func moveWarehousePackagingSource(ctx context.Context, tx pgx.Tx, tenant, organization, activityID, sourceKind string, line inventorycontrol.WarehouseActivityLine) error {
	var balanceID int64
	var sourceUOM, sourceQuantity, sourceFactor, sourceBase, targetUOM, targetQuantity, targetFactor, movedBase string
	var conversion bool
	err := tx.QueryRow(ctx, `select source_balance_id,source_uom_code,source_uom_quantity::text,source_qty_per_uom::text,source_base_quantity::text,target_uom_code,target_uom_quantity::text,target_qty_per_uom::text,moved_base_quantity::text,conversion_required from inventory.warehouse_activity_uom_reservation where tenant_id=$1 and activity_id=$2 and line_id=$3 and organization_id=$4 and status='open' and version=1 for update`, tenant, activityID, line.ID, organization).Scan(&balanceID, &sourceUOM, &sourceQuantity, &sourceFactor, &sourceBase, &targetUOM, &targetQuantity, &targetFactor, &movedBase, &conversion)
	if errors.Is(err, pgx.ErrNoRows) {
		return inventorycontrol.ErrConflict
	}
	if err != nil {
		return err
	}
	if movedBase != line.Quantity || sourceUOM != line.FromUOMCode || sourceQuantity != line.FromUOMQuantity || sourceFactor != line.FromQuantityPerUOM || targetUOM != line.UOMCode || targetQuantity != line.UOMQuantity || targetFactor != line.QuantityPerUOM {
		return inventorycontrol.ErrConflict
	}
	var sourceIsBase, targetIsBase bool
	var targetPrecision string
	err = tx.QueryRow(ctx, `select s.is_base,t.is_base,t.rounding_precision::text from inventory.bulk_balance b join inventory.item_unit_of_measure s on s.tenant_id=b.tenant_id and s.item_id=b.item_id and s.uom_code=$2 join inventory.item_unit_of_measure t on t.tenant_id=b.tenant_id and t.item_id=b.item_id and t.uom_code=$3 where b.balance_id=$1 for share of b,s,t`, balanceID, sourceUOM, targetUOM).Scan(&sourceIsBase, &targetIsBase, &targetPrecision)
	if err != nil {
		return err
	}
	if sourceIsBase {
		updated, updateErr := tx.Exec(ctx, `update inventory.bulk_uom_balance set reserved_quantity=reserved_quantity-$3::numeric,version=version+1,updated_at=clock_timestamp() where balance_id=$1 and uom_code=$2 and reserved_quantity >= $3::numeric`, balanceID, sourceUOM, sourceQuantity)
		if updateErr != nil || updated.RowsAffected() != 1 {
			if updateErr != nil {
				return bulkConflict(updateErr)
			}
			return inventorycontrol.ErrConflict
		}
	} else {
		updated, updateErr := tx.Exec(ctx, `update inventory.bulk_uom_balance set quantity=quantity-$3::numeric,reserved_quantity=reserved_quantity-$3::numeric,quantity_base=quantity_base-$4::numeric,version=version+1,updated_at=clock_timestamp() where balance_id=$1 and uom_code=$2 and quantity >= $3::numeric and reserved_quantity >= $3::numeric and quantity_base >= $4::numeric`, balanceID, sourceUOM, sourceQuantity, sourceBase)
		if updateErr != nil || updated.RowsAffected() != 1 {
			if updateErr != nil {
				return bulkConflict(updateErr)
			}
			return inventorycontrol.ErrConflict
		}
		if _, err = tx.Exec(ctx, `delete from inventory.bulk_uom_balance where balance_id=$1 and uom_code=$2 and quantity_base=0 and reserved_quantity=0`, balanceID, sourceUOM); err != nil {
			return err
		}
		_, err = tx.Exec(ctx, `insert into inventory.bulk_uom_balance(balance_id,uom_code,quantity,qty_per_uom,quantity_base) select $1,u.uom_code,$2::numeric/u.qty_per_uom,u.qty_per_uom,$2::numeric from inventory.bulk_balance b join inventory.item_unit_of_measure u on u.tenant_id=b.tenant_id and u.item_id=b.item_id and u.is_base where b.balance_id=$1 on conflict(balance_id,uom_code) do update set quantity=inventory.bulk_uom_balance.quantity+excluded.quantity,quantity_base=inventory.bulk_uom_balance.quantity_base+excluded.quantity_base,version=inventory.bulk_uom_balance.version+1,updated_at=clock_timestamp()`, balanceID, sourceBase)
		if err != nil {
			return bulkConflict(err)
		}
	}
	updated, err := tx.Exec(ctx, `update inventory.bulk_balance set quantity=quantity-$2::numeric,reserved_quantity=reserved_quantity-$2::numeric,version=version+1,updated_at=clock_timestamp() where balance_id=$1 and quantity >= $2::numeric and reserved_quantity >= $2::numeric`, balanceID, movedBase)
	if err != nil || updated.RowsAffected() != 1 {
		if err != nil {
			return err
		}
		return inventorycontrol.ErrConflict
	}
	if conversion {
		toTotal, exact := exactUOMQuantity(sourceBase, targetFactor, targetPrecision)
		if !exact {
			return inventorycontrol.ErrConflict
		}
		conversionID := activityID + ":" + line.ID + ":breakbulk"
		requestID := activityID + ":" + line.ID + ":automatic-breakbulk"
		_, err = tx.Exec(ctx, `insert into inventory.bulk_uom_conversion(tenant_id,conversion_id,request_id,organization_id,bin_id,item_id,lot_id,operation,from_uom,to_uom,from_quantity,to_quantity,base_quantity,from_qty_per_uom,to_qty_per_uom,source_kind,source_id) values($1,$2,$3,$4,$5,$6,nullif($7,''),'breakbulk',$8,$9,$10::numeric,$11::numeric,$12::numeric,$13::numeric,$14::numeric,$15,$16)`, tenant, conversionID, requestID, organization, line.FromBinID, line.ItemID, line.LotID, sourceUOM, targetUOM, sourceQuantity, toTotal, sourceBase, sourceFactor, targetFactor, sourceKind, activityID)
		if err != nil {
			return bulkConflict(err)
		}
		_, err = tx.Exec(ctx, `insert into inventory.bulk_uom_conversion_line(tenant_id,conversion_id,line_id,sequence_no,action_type,uom_code,quantity,quantity_base) values($1,$2,$3,1,'take',$5,$6::numeric,$7::numeric),($1,$2,$4,2,'place',$8,$9::numeric,$7::numeric)`, tenant, conversionID, conversionID+":take", conversionID+":place", sourceUOM, sourceQuantity, sourceBase, targetUOM, toTotal)
		if err != nil {
			return bulkConflict(err)
		}
		if err = recordBulkEvent(ctx, tx, tenant, deterministicWarehouseUUID(tenant, activityID, line.ID, "breakbulk"), "bulk-uom-conversion", conversionID, "bulk-uom-conversion.completed", 1, map[string]any{"organization_id": organization, "activity_id": activityID, "line_id": line.ID, "from_uom": sourceUOM, "to_uom": targetUOM, "source_base_quantity": sourceBase, "moved_base_quantity": movedBase}); err != nil {
			return err
		}
	}
	leftover, leftOK := parseRat(sourceBase)
	moved, movedOK := parseRat(movedBase)
	if !leftOK || !movedOK {
		return inventorycontrol.ErrConflict
	}
	leftover.Sub(leftover, moved)
	if leftover.Sign() > 0 && !targetIsBase {
		leftQuantity, exact := exactUOMQuantity(formatRat(leftover, 6), targetFactor, targetPrecision)
		if !exact {
			return inventorycontrol.ErrConflict
		}
		var baseUOM string
		if err = tx.QueryRow(ctx, `select uom_code from inventory.item_unit_of_measure u join inventory.bulk_balance b on b.tenant_id=u.tenant_id and b.item_id=u.item_id where b.balance_id=$1 and u.is_base`, balanceID).Scan(&baseUOM); err != nil {
			return err
		}
		updated, err = tx.Exec(ctx, `update inventory.bulk_uom_balance set quantity=quantity-($2::numeric/qty_per_uom),quantity_base=quantity_base-$2::numeric,version=version+1,updated_at=clock_timestamp() where balance_id=$1 and uom_code=$3 and quantity_base >= $2::numeric`, balanceID, formatRat(leftover, 6), baseUOM)
		if err != nil || updated.RowsAffected() != 1 {
			if err != nil {
				return bulkConflict(err)
			}
			return inventorycontrol.ErrConflict
		}
		_, err = tx.Exec(ctx, `insert into inventory.bulk_uom_balance(balance_id,uom_code,quantity,qty_per_uom,quantity_base) values($1,$2,$3::numeric,$4::numeric,$5::numeric) on conflict(balance_id,uom_code) do update set quantity=inventory.bulk_uom_balance.quantity+excluded.quantity,quantity_base=inventory.bulk_uom_balance.quantity_base+excluded.quantity_base,version=inventory.bulk_uom_balance.version+1,updated_at=clock_timestamp()`, balanceID, targetUOM, leftQuantity, targetFactor, formatRat(leftover, 6))
		if err != nil {
			return bulkConflict(err)
		}
	}
	updated, err = tx.Exec(ctx, `update inventory.warehouse_activity_uom_reservation set status='consumed',version=version+1 where tenant_id=$1 and activity_id=$2 and line_id=$3 and status='open' and version=1`, tenant, activityID, line.ID)
	if err != nil || updated.RowsAffected() != 1 {
		if err != nil {
			return err
		}
		return inventorycontrol.ErrConflict
	}
	return nil
}

func releaseWarehousePackagingReservation(ctx context.Context, tx pgx.Tx, tenant, organization, activityID string, line inventorycontrol.WarehouseActivityLine) error {
	var balanceID int64
	var sourceUOM, sourceQuantity, movedBase string
	err := tx.QueryRow(ctx, `select source_balance_id,source_uom_code,source_uom_quantity::text,moved_base_quantity::text from inventory.warehouse_activity_uom_reservation where tenant_id=$1 and activity_id=$2 and line_id=$3 and organization_id=$4 and status='open' and version=1 for update`, tenant, activityID, line.ID, organization).Scan(&balanceID, &sourceUOM, &sourceQuantity, &movedBase)
	if errors.Is(err, pgx.ErrNoRows) {
		return inventorycontrol.ErrConflict
	}
	if err != nil {
		return err
	}
	updated, err := tx.Exec(ctx, `update inventory.bulk_balance set reserved_quantity=reserved_quantity-$2::numeric,version=version+1,updated_at=clock_timestamp() where balance_id=$1 and reserved_quantity >= $2::numeric`, balanceID, movedBase)
	if err != nil || updated.RowsAffected() != 1 {
		if err != nil {
			return err
		}
		return inventorycontrol.ErrConflict
	}
	updated, err = tx.Exec(ctx, `update inventory.bulk_uom_balance set reserved_quantity=reserved_quantity-$3::numeric,version=version+1,updated_at=clock_timestamp() where balance_id=$1 and uom_code=$2 and reserved_quantity >= $3::numeric`, balanceID, sourceUOM, sourceQuantity)
	if err != nil || updated.RowsAffected() != 1 {
		if err != nil {
			return err
		}
		return inventorycontrol.ErrConflict
	}
	updated, err = tx.Exec(ctx, `update inventory.warehouse_activity_uom_reservation set status='released',version=version+1 where tenant_id=$1 and activity_id=$2 and line_id=$3 and status='open' and version=1`, tenant, activityID, line.ID)
	if err != nil || updated.RowsAffected() != 1 {
		if err != nil {
			return err
		}
		return inventorycontrol.ErrConflict
	}
	return nil
}

func consumeRegisteredPickPackaging(ctx context.Context, tx pgx.Tx, tenant, organization, activityID, reservationID string) error {
	var balanceID int64
	var uomCode, uomQuantity, baseQuantity, baseUOM, baseFactor string
	var isBase bool
	err := tx.QueryRow(ctx, `select b.balance_id,al.uom_code,al.uom_quantity::text,al.quantity::text,u.is_base,base.uom_code,base.qty_per_uom::text from inventory.warehouse_activity a join inventory.warehouse_activity_line al using(tenant_id,activity_id) join inventory.bulk_reservation r on r.tenant_id=al.tenant_id and r.reservation_id=al.reservation_id join inventory.bulk_balance b on b.tenant_id=r.tenant_id and b.organization_id=r.organization_id and b.bin_id=r.bin_id and b.item_id=r.item_id and b.lot_id is not distinct from r.lot_id join inventory.item_unit_of_measure u on u.tenant_id=b.tenant_id and u.item_id=b.item_id and u.uom_code=al.uom_code join inventory.item_unit_of_measure base on base.tenant_id=b.tenant_id and base.item_id=b.item_id and base.is_base where a.tenant_id=$1 and a.organization_id=$2 and a.activity_id=$3 and a.activity_type='pick' and a.status='registered' and r.reservation_id=$4 and r.status='reservation' and b.quantity>=al.quantity and b.reserved_quantity>=al.quantity for update of b,u,base`, tenant, organization, activityID, reservationID).Scan(&balanceID, &uomCode, &uomQuantity, &baseQuantity, &isBase, &baseUOM, &baseFactor)
	if errors.Is(err, pgx.ErrNoRows) {
		return inventorycontrol.ErrConflict
	}
	if err != nil {
		return err
	}
	if isBase {
		return nil
	}
	updated, err := tx.Exec(ctx, `update inventory.bulk_uom_balance set quantity=quantity-$3::numeric,quantity_base=quantity_base-$4::numeric,version=version+1,updated_at=clock_timestamp() where balance_id=$1 and uom_code=$2 and quantity-reserved_quantity >= $3::numeric and quantity_base >= $4::numeric`, balanceID, uomCode, uomQuantity, baseQuantity)
	if err != nil || updated.RowsAffected() != 1 {
		if err != nil {
			return bulkConflict(err)
		}
		return inventorycontrol.ErrConflict
	}
	if _, err = tx.Exec(ctx, `delete from inventory.bulk_uom_balance where balance_id=$1 and uom_code=$2 and quantity_base=0 and reserved_quantity=0`, balanceID, uomCode); err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `insert into inventory.bulk_uom_balance(balance_id,uom_code,quantity,qty_per_uom,quantity_base) values($1,$2,$3::numeric/$4::numeric,$4::numeric,$3::numeric) on conflict(balance_id,uom_code) do update set quantity=inventory.bulk_uom_balance.quantity+excluded.quantity,quantity_base=inventory.bulk_uom_balance.quantity_base+excluded.quantity_base,version=inventory.bulk_uom_balance.version+1,updated_at=clock_timestamp()`, balanceID, baseUOM, baseQuantity, baseFactor)
	return bulkConflict(err)
}
