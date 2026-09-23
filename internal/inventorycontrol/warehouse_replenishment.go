package inventorycontrol

import (
	"context"
	"fmt"
)

type WarehouseReplenishmentCommand struct {
	RequestID      string `json:"request_id"`
	OrganizationID string `json:"organization_id"`
	ToBinID        string `json:"to_bin_id"`
	ItemID         string `json:"item_id"`
	TargetUOM      string `json:"target_uom,omitempty"`
	AllowBreakbulk bool   `json:"allow_breakbulk"`
	UseFEFO        bool   `json:"use_fefo"`
	AssignedTo     string `json:"assigned_to,omitempty"`
}

func (s *WarehouseService) CreateReplenishment(ctx context.Context, tenant string, value WarehouseReplenishmentCommand) (WarehouseActivity, error) {
	if tenant == "" || value.RequestID == "" || value.OrganizationID == "" || value.ToBinID == "" || value.ItemID == "" {
		return WarehouseActivity{}, fmt.Errorf("invalid warehouse replenishment")
	}
	return s.repository.CreateWarehouseReplenishment(ctx, tenant, s.newIDs(), value)
}

func (s *WarehouseService) CancelReplenishment(ctx context.Context, tenant, organization, activity string, version int64) (WarehouseActivity, error) {
	if tenant == "" || organization == "" || activity == "" || version < 1 {
		return WarehouseActivity{}, fmt.Errorf("invalid warehouse replenishment cancellation")
	}
	return s.repository.CancelWarehouseReplenishment(ctx, tenant, organization, activity, version, s.ids.New())
}
