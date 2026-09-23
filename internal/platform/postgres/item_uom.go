package postgres

import (
	"context"
	"errors"

	"elite.local/enterprise/internal/inventorycontrol"
	"github.com/jackc/pgx/v5"
)

func (r *InventoryControl) ConfigureItemUnitOfMeasure(ctx context.Context, tenant, eventID string, value inventorycontrol.ItemUnitOfMeasure) (inventorycontrol.ItemUnitOfMeasure, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return value, err
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, `insert into inventory.item_unit_of_measure(tenant_id,item_id,uom_code,qty_per_uom,rounding_precision,is_base,version) values($1,$2,$3,$4::numeric,$5::numeric,false,1)`, tenant, value.ItemID, value.Code, value.QuantityPerUnit, value.RoundingPrecision)
	if err != nil {
		return value, bulkConflict(err)
	}
	if err = recordBulkEvent(ctx, tx, tenant, eventID, "item-unit-of-measure", value.ItemID+":"+value.Code, "item-unit-of-measure.created", 1, map[string]string{"item_id": value.ItemID, "code": value.Code, "quantity_per_unit": value.QuantityPerUnit, "rounding_precision": value.RoundingPrecision}); err != nil {
		return value, err
	}
	return value, tx.Commit(ctx)
}

func (r *InventoryControl) ConvertItemUnitOfMeasure(ctx context.Context, tenant, itemID, code, quantity string) (inventorycontrol.UnitOfMeasureConversion, error) {
	result := inventorycontrol.UnitOfMeasureConversion{ItemID: itemID, Code: code, Quantity: quantity}
	err := r.pool.QueryRow(ctx, `select u.qty_per_uom::text,i.base_uom,(($4::numeric)*u.qty_per_uom)::numeric(20,6)::text,b.rounding_precision::text from inventory.item_unit_of_measure u join inventory.stock_item i using(tenant_id,item_id) join inventory.item_unit_of_measure b on b.tenant_id=u.tenant_id and b.item_id=u.item_id and b.is_base where u.tenant_id=$1 and u.item_id=$2 and u.uom_code=$3 and mod(($4::numeric)*u.qty_per_uom,b.rounding_precision)=0`, tenant, itemID, code, quantity).Scan(&result.QuantityPerUnit, &result.BaseCode, &result.BaseQuantity, &result.RoundingPrecision)
	if errors.Is(err, pgx.ErrNoRows) {
		return inventorycontrol.UnitOfMeasureConversion{}, inventorycontrol.ErrConflict
	}
	if err != nil {
		return inventorycontrol.UnitOfMeasureConversion{}, err
	}
	return result, nil
}
