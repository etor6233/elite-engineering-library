package inventorycontrol

import (
	"context"
	"fmt"
)

type ItemUnitOfMeasure struct {
	ItemID            string `json:"item_id"`
	Code              string `json:"code"`
	QuantityPerUnit   string `json:"quantity_per_unit"`
	RoundingPrecision string `json:"rounding_precision"`
	Base              bool   `json:"base"`
	Version           int64  `json:"version"`
}

type UnitOfMeasureConversion struct {
	ItemID            string `json:"item_id"`
	Code              string `json:"code"`
	Quantity          string `json:"quantity"`
	QuantityPerUnit   string `json:"quantity_per_unit"`
	BaseCode          string `json:"base_code"`
	BaseQuantity      string `json:"base_quantity"`
	RoundingPrecision string `json:"rounding_precision"`
}

func (s *BulkService) ConfigureUnitOfMeasure(ctx context.Context, tenant string, value ItemUnitOfMeasure) (ItemUnitOfMeasure, error) {
	if tenant == "" || value.ItemID == "" || value.Code == "" || value.Base || !positiveDecimal(value.QuantityPerUnit, quantityPattern) || !positiveDecimal(value.RoundingPrecision, quantityPattern) {
		return ItemUnitOfMeasure{}, fmt.Errorf("invalid item unit of measure")
	}
	value.Version = 1
	return s.repository.ConfigureItemUnitOfMeasure(ctx, tenant, s.ids.New(), value)
}

func (s *BulkService) ConvertUnitOfMeasure(ctx context.Context, tenant, itemID, code, quantity string) (UnitOfMeasureConversion, error) {
	if tenant == "" || itemID == "" || code == "" || !positiveDecimal(quantity, quantityPattern) {
		return UnitOfMeasureConversion{}, fmt.Errorf("invalid unit of measure conversion")
	}
	return s.repository.ConvertItemUnitOfMeasure(ctx, tenant, itemID, code, quantity)
}
