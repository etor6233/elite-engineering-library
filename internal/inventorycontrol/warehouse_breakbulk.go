package inventorycontrol

import (
	"context"
	"fmt"
)

type PackagingConversionCommand struct {
	RequestID      string `json:"request_id"`
	OrganizationID string `json:"organization_id"`
	BinID          string `json:"bin_id"`
	ItemID         string `json:"item_id"`
	LotID          string `json:"lot_id,omitempty"`
	Operation      string `json:"operation"`
	FromUOM        string `json:"from_uom"`
	ToUOM          string `json:"to_uom"`
	FromQuantity   string `json:"from_quantity"`
}

type PackagingConversionIDs struct {
	ConversionID string
	TakeLineID   string
	PlaceLineID  string
	EventID      string
}

type PackagingConversionLine struct {
	Sequence     int    `json:"sequence"`
	Action       string `json:"action"`
	UOM          string `json:"uom"`
	Quantity     string `json:"quantity"`
	BaseQuantity string `json:"base_quantity"`
}

type PackagingConversionResult struct {
	ID              string                    `json:"id"`
	RequestID       string                    `json:"request_id"`
	OrganizationID  string                    `json:"organization_id"`
	BinID           string                    `json:"bin_id"`
	ItemID          string                    `json:"item_id"`
	LotID           string                    `json:"lot_id,omitempty"`
	Operation       string                    `json:"operation"`
	FromUOM         string                    `json:"from_uom"`
	ToUOM           string                    `json:"to_uom"`
	FromQuantity    string                    `json:"from_quantity"`
	ToQuantity      string                    `json:"to_quantity"`
	BaseQuantity    string                    `json:"base_quantity"`
	FromQuantityPer string                    `json:"from_quantity_per_uom"`
	ToQuantityPer   string                    `json:"to_quantity_per_uom"`
	Version         int64                     `json:"version"`
	Lines           []PackagingConversionLine `json:"lines"`
}

func (s *BulkService) ConvertHandlingUnits(ctx context.Context, tenant string, value PackagingConversionCommand) (PackagingConversionResult, error) {
	operations := map[string]bool{"breakbulk": true, "gather": true}
	if tenant == "" || value.RequestID == "" || value.OrganizationID == "" || value.BinID == "" || value.ItemID == "" || !operations[value.Operation] || value.FromUOM == "" || value.ToUOM == "" || value.FromUOM == value.ToUOM || !positiveDecimal(value.FromQuantity, quantityPattern) {
		return PackagingConversionResult{}, fmt.Errorf("invalid packaging conversion")
	}
	ids := PackagingConversionIDs{ConversionID: s.ids.New(), TakeLineID: s.ids.New(), PlaceLineID: s.ids.New(), EventID: s.ids.New()}
	return s.repository.ConvertBulkHandlingUnits(ctx, tenant, ids, value)
}
