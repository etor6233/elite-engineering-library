package operations

import (
	"context"
	"errors"
	"fmt"
	"regexp"
)

var ErrConflict = errors.New("operation conflict")
var code = regexp.MustCompile(`^[A-Z]{3}$`)

type PurchaseOrder struct {
	ID                        string `json:"id"`
	SupplierID                string `json:"supplier_id"`
	DestinationOrganizationID string `json:"destination_organization_id"`
	State                     string `json:"state"`
	Currency                  string `json:"currency"`
	TotalMinorUnits           int64  `json:"total_minor_units"`
	Version                   int64  `json:"version"`
}
type ProductionUnit struct {
	ID                  string `json:"id"`
	OrganizationID      string `json:"organization_id"`
	PurchaseOrderID     string `json:"purchase_order_id"`
	VariantID           string `json:"variant_id"`
	SerialNumber        string `json:"serial_number"`
	VIN                 string `json:"vin,omitempty"`
	BatterySerialNumber string `json:"battery_serial_number,omitempty"`
	State               string `json:"state"`
}
type StockUnit struct {
	ID                  string `json:"id"`
	OrganizationID      string `json:"organization_id"`
	VariantID           string `json:"variant_id"`
	ProductionUnitID    string `json:"production_unit_id,omitempty"`
	SerialNumber        string `json:"serial_number"`
	VIN                 string `json:"vin,omitempty"`
	BatterySerialNumber string `json:"battery_serial_number,omitempty"`
	State               string `json:"state"`
	Version             int64  `json:"version"`
}
type Repository interface {
	CreatePurchaseOrder(context.Context, string, string, PurchaseOrder) error
	TransitionPurchaseOrder(context.Context, string, string, string, string, int64, string, string) error
	CreateProductionUnit(context.Context, string, string, ProductionUnit) error
	TransitionProductionUnit(context.Context, string, string, string, string, string, string) error
	CreateStockUnit(context.Context, string, string, StockUnit) error
	TransitionStockUnit(context.Context, string, string, string, string, int64, string, string) error
}
type IDGenerator interface{ New() string }
type Service struct {
	repository Repository
	ids        IDGenerator
}

func NewService(r Repository, ids IDGenerator) *Service { return &Service{repository: r, ids: ids} }

var poTransitions = map[string]map[string]bool{"draft": {"submitted": true, "cancelled": true}, "submitted": {"accepted": true, "cancelled": true}, "accepted": {"in-production": true, "cancelled": true}, "in-production": {"shipped": true}, "shipped": {"received": true}}
var factoryTransitions = map[string]map[string]bool{"planned": {"assembly": true, "rejected": true}, "assembly": {"quality": true, "rejected": true}, "quality": {"released": true, "rejected": true}, "released": {"shipped": true}, "shipped": {"received": true}}
var stockTransitions = map[string]map[string]bool{"in-transit": {"available": true, "quarantine": true}, "available": {"reserved": true, "service": true, "quarantine": true}, "reserved": {"available": true, "sold": true}, "sold": {"service": true}, "service": {"available": true, "retired": true}, "quarantine": {"available": true, "retired": true}}

func (s *Service) CreatePurchaseOrder(ctx context.Context, tenant string, input PurchaseOrder) (PurchaseOrder, error) {
	if tenant == "" || input.SupplierID == "" || input.DestinationOrganizationID == "" || !code.MatchString(input.Currency) || input.TotalMinorUnits < 0 {
		return PurchaseOrder{}, fmt.Errorf("invalid purchase order")
	}
	input.ID = s.ids.New()
	input.State = "draft"
	input.Version = 1
	if err := s.repository.CreatePurchaseOrder(ctx, tenant, s.ids.New(), input); err != nil {
		return PurchaseOrder{}, err
	}
	return input, nil
}
func (s *Service) TransitionPurchaseOrder(ctx context.Context, tenant, organization, id, current, target string, version int64) error {
	if organization == "" || !poTransitions[current][target] || version < 1 {
		return fmt.Errorf("%w: invalid purchase order transition", ErrConflict)
	}
	return s.repository.TransitionPurchaseOrder(ctx, tenant, organization, id, current, version, target, s.ids.New())
}
func (s *Service) RegisterProductionUnit(ctx context.Context, tenant string, input ProductionUnit) (ProductionUnit, error) {
	if tenant == "" || input.OrganizationID == "" || input.PurchaseOrderID == "" || input.VariantID == "" || input.SerialNumber == "" {
		return ProductionUnit{}, fmt.Errorf("invalid production unit")
	}
	input.ID = s.ids.New()
	input.State = "planned"
	if err := s.repository.CreateProductionUnit(ctx, tenant, s.ids.New(), input); err != nil {
		return ProductionUnit{}, err
	}
	return input, nil
}
func (s *Service) TransitionProductionUnit(ctx context.Context, tenant, organization, id, current, target string) error {
	if organization == "" || !factoryTransitions[current][target] {
		return fmt.Errorf("%w: invalid production transition", ErrConflict)
	}
	return s.repository.TransitionProductionUnit(ctx, tenant, organization, id, current, target, s.ids.New())
}
func (s *Service) ReceiveStockUnit(ctx context.Context, tenant string, input StockUnit) (StockUnit, error) {
	if tenant == "" || input.OrganizationID == "" || input.VariantID == "" || input.SerialNumber == "" {
		return StockUnit{}, fmt.Errorf("invalid stock unit")
	}
	input.ID = s.ids.New()
	input.State = "in-transit"
	input.Version = 1
	if err := s.repository.CreateStockUnit(ctx, tenant, s.ids.New(), input); err != nil {
		return StockUnit{}, err
	}
	return input, nil
}
func (s *Service) TransitionStockUnit(ctx context.Context, tenant, organization, id, current, target string, version int64) error {
	if organization == "" || !stockTransitions[current][target] || version < 1 {
		return fmt.Errorf("%w: invalid stock transition", ErrConflict)
	}
	return s.repository.TransitionStockUnit(ctx, tenant, organization, id, current, version, target, s.ids.New())
}
