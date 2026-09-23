package inventorycontrol

import (
	"context"
	"fmt"
)

type WarehouseCrossDockPolicy struct {
	OrganizationID string `json:"organization_id"`
	ItemID         string `json:"item_id"`
	BinID          string `json:"bin_id"`
	DueDateDays    int    `json:"due_date_days"`
	Enabled        bool   `json:"enabled"`
	Version        int64  `json:"version"`
}

func (s *WarehouseService) ConfigureCrossDock(ctx context.Context, tenant string, value WarehouseCrossDockPolicy) (WarehouseCrossDockPolicy, error) {
	if tenant == "" || value.OrganizationID == "" || value.ItemID == "" || value.BinID == "" || value.DueDateDays < 0 || value.DueDateDays > 365 || value.Version != 1 {
		return value, fmt.Errorf("invalid warehouse cross-dock policy")
	}
	return s.repository.ConfigureWarehouseCrossDock(ctx, tenant, s.ids.New(), value)
}
