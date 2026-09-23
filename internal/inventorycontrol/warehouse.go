package inventorycontrol

import (
	"context"
	"fmt"
	"time"
)

type WarehouseReceiptCommand struct {
	RequestID        string     `json:"request_id"`
	OrganizationID   string     `json:"organization_id"`
	ReceiveBinID     string     `json:"receive_bin_id"`
	ItemID           string     `json:"item_id"`
	LotNo            string     `json:"lot_no,omitempty"`
	ExpirationDate   *time.Time `json:"expiration_date,omitempty"`
	WarrantyDate     *time.Time `json:"warranty_date,omitempty"`
	Quantity         string     `json:"quantity,omitempty"`
	HandlingUOM      string     `json:"handling_uom,omitempty"`
	HandlingQuantity string     `json:"handling_quantity,omitempty"`
	AllowBreakbulk   bool       `json:"allow_breakbulk"`
	UnitCost         string     `json:"unit_cost"`
	PostingDate      time.Time  `json:"posting_date"`
	SourceKind       string     `json:"source_kind"`
	SourceID         string     `json:"source_id"`
}

type WarehousePickCommand struct {
	RequestID      string `json:"request_id"`
	OrganizationID string `json:"organization_id"`
	ShipBinID      string `json:"ship_bin_id"`
	ItemID         string `json:"item_id"`
	DemandKind     string `json:"demand_kind"`
	DemandID       string `json:"demand_id"`
	DemandLineID   string `json:"demand_line_id"`
	Quantity       string `json:"quantity"`
	HandlingUOM    string `json:"handling_uom,omitempty"`
	AllowBreakbulk bool   `json:"allow_breakbulk"`
	UseFEFO        bool   `json:"use_fefo"`
	AllowDedicated bool   `json:"allow_dedicated"`
	AssignedTo     string `json:"assigned_to,omitempty"`
}

type SalesWarehouseBinding struct {
	RequestID    string `json:"request_id"`
	VariantID    string `json:"variant_id"`
	ItemID       string `json:"item_id"`
	SalesUOMCode string `json:"sales_uom_code"`
	Version      int64  `json:"version"`
}

type CustomerShipmentCommand struct {
	RequestID           string    `json:"request_id"`
	OrganizationID      string    `json:"organization_id"`
	OrderID             string    `json:"order_id"`
	WarehouseActivityID string    `json:"warehouse_activity_id"`
	PostingDate         time.Time `json:"posting_date"`
}

type CustomerShipmentLine struct {
	ID              string `json:"id"`
	OrderLineID     string `json:"order_line_id"`
	VariantID       string `json:"variant_id"`
	ItemID          string `json:"item_id"`
	SalesUOMCode    string `json:"sales_uom_code"`
	Quantity        string `json:"quantity"`
	QuantityBase    string `json:"quantity_base"`
	CostAmount      string `json:"cost_amount"`
	AllocationCount int    `json:"allocation_count"`
}

type CustomerShipment struct {
	ID                  string               `json:"id"`
	RequestID           string               `json:"request_id"`
	OrganizationID      string               `json:"organization_id"`
	OrderID             string               `json:"order_id"`
	WarehouseActivityID string               `json:"warehouse_activity_id"`
	PostingDate         time.Time            `json:"posting_date"`
	OrderVersion        int64                `json:"order_version"`
	FulfillmentState    string               `json:"fulfillment_state"`
	Line                CustomerShipmentLine `json:"line"`
}

type CustomerShipmentIDs struct {
	ShipmentID string
	LineID     string
	EventID    string
}

type WarehouseActivityLine struct {
	ID                 string     `json:"id"`
	Sequence           int        `json:"sequence"`
	FromBinID          string     `json:"from_bin_id"`
	ToBinID            string     `json:"to_bin_id"`
	ItemID             string     `json:"item_id"`
	LotID              string     `json:"lot_id,omitempty"`
	LotNo              string     `json:"lot_no,omitempty"`
	ReservationID      string     `json:"reservation_id,omitempty"`
	ReservationVersion int64      `json:"reservation_version,omitempty"`
	Quantity           string     `json:"quantity"`
	FromUOMCode        string     `json:"from_uom_code"`
	FromUOMQuantity    string     `json:"from_uom_quantity"`
	FromQuantityPerUOM string     `json:"from_quantity_per_uom"`
	UOMCode            string     `json:"uom_code"`
	UOMQuantity        string     `json:"uom_quantity"`
	QuantityPerUOM     string     `json:"quantity_per_uom"`
	ExpirationDate     *time.Time `json:"expiration_date,omitempty"`
	SourceEntryID      string     `json:"source_entry_id,omitempty"`
}

type WarehouseActivity struct {
	ID             string                  `json:"id"`
	OrganizationID string                  `json:"organization_id"`
	Type           string                  `json:"type"`
	Status         string                  `json:"status"`
	SourceKind     string                  `json:"source_kind"`
	SourceID       string                  `json:"source_id"`
	SourceLineID   string                  `json:"source_line_id,omitempty"`
	RequestID      string                  `json:"request_id"`
	AssignedTo     string                  `json:"assigned_to,omitempty"`
	Version        int64                   `json:"version"`
	Lines          []WarehouseActivityLine `json:"lines"`
}

type WarehouseReceiptResult struct {
	ReceiptID        string            `json:"receipt_id"`
	EntryID          string            `json:"entry_id"`
	LotID            string            `json:"lot_id,omitempty"`
	Quantity         string            `json:"quantity"`
	HandlingUOM      string            `json:"handling_uom"`
	HandlingQuantity string            `json:"handling_quantity"`
	CostAmount       string            `json:"cost_amount"`
	PutAway          WarehouseActivity `json:"put_away"`
}

type WarehouseIDs struct {
	ActivityID        string
	ReceiptID         string
	LineID            string
	EntryID           string
	LayerID           string
	LotID             string
	EventID           string
	ConversionID      string
	TakeLineID        string
	PlaceLineID       string
	ConversionEventID string
}

type WarehouseRepository interface {
	ConfigureSalesWarehouseBinding(context.Context, string, string, SalesWarehouseBinding) (SalesWarehouseBinding, error)
	PostCustomerShipment(context.Context, string, CustomerShipmentIDs, CustomerShipmentCommand) (CustomerShipment, error)
	PostWarehouseReceipt(context.Context, string, WarehouseIDs, WarehouseReceiptCommand) (WarehouseReceiptResult, error)
	CreateWarehousePick(context.Context, string, WarehouseIDs, WarehousePickCommand) (WarehouseActivity, error)
	CreateWarehouseReplenishment(context.Context, string, WarehouseIDs, WarehouseReplenishmentCommand) (WarehouseActivity, error)
	ConfigureWarehouseCrossDock(context.Context, string, string, WarehouseCrossDockPolicy) (WarehouseCrossDockPolicy, error)
	RegisterWarehouseActivity(context.Context, string, string, string, int64, time.Time, string) (WarehouseActivity, error)
	CancelWarehousePick(context.Context, string, string, string, int64, string) (WarehouseActivity, error)
	CancelWarehousePutAway(context.Context, string, string, string, int64, string) (WarehouseActivity, error)
	CancelWarehouseReplenishment(context.Context, string, string, string, int64, string) (WarehouseActivity, error)
}

type WarehouseService struct {
	repository WarehouseRepository
	ids        IDGenerator
}

func NewWarehouseService(repository WarehouseRepository, ids IDGenerator) *WarehouseService {
	return &WarehouseService{repository: repository, ids: ids}
}

func (s *WarehouseService) newIDs() WarehouseIDs {
	return WarehouseIDs{ActivityID: s.ids.New(), ReceiptID: s.ids.New(), LineID: s.ids.New(), EntryID: s.ids.New(), LayerID: s.ids.New(), LotID: s.ids.New(), EventID: s.ids.New(), ConversionID: s.ids.New(), TakeLineID: s.ids.New(), PlaceLineID: s.ids.New(), ConversionEventID: s.ids.New()}
}

func (s *WarehouseService) ConfigureSalesBinding(ctx context.Context, tenant string, value SalesWarehouseBinding) (SalesWarehouseBinding, error) {
	if tenant == "" || value.RequestID == "" || value.VariantID == "" || value.ItemID == "" || value.SalesUOMCode == "" {
		return SalesWarehouseBinding{}, fmt.Errorf("invalid sales warehouse binding")
	}
	value.Version = 1
	return s.repository.ConfigureSalesWarehouseBinding(ctx, tenant, s.ids.New(), value)
}

func (s *WarehouseService) PostCustomerShipment(ctx context.Context, tenant string, value CustomerShipmentCommand) (CustomerShipment, error) {
	if tenant == "" || value.RequestID == "" || value.OrganizationID == "" || value.OrderID == "" || value.WarehouseActivityID == "" || value.PostingDate.IsZero() {
		return CustomerShipment{}, fmt.Errorf("invalid customer shipment")
	}
	value.PostingDate = value.PostingDate.UTC()
	ids := CustomerShipmentIDs{ShipmentID: s.ids.New(), LineID: s.ids.New(), EventID: s.ids.New()}
	return s.repository.PostCustomerShipment(ctx, tenant, ids, value)
}

func (s *WarehouseService) PostReceipt(ctx context.Context, tenant string, value WarehouseReceiptCommand) (WarehouseReceiptResult, error) {
	legacyQuantity := positiveDecimal(value.Quantity, quantityPattern) && value.HandlingUOM == "" && value.HandlingQuantity == ""
	handlingQuantity := value.Quantity == "" && value.HandlingUOM != "" && positiveDecimal(value.HandlingQuantity, quantityPattern)
	if tenant == "" || value.RequestID == "" || value.OrganizationID == "" || value.ReceiveBinID == "" || value.ItemID == "" || (!legacyQuantity && !handlingQuantity) || !positiveDecimal(value.UnitCost, costPattern) || value.PostingDate.IsZero() || value.SourceKind == "" || value.SourceID == "" {
		return WarehouseReceiptResult{}, fmt.Errorf("invalid warehouse receipt")
	}
	value.PostingDate = value.PostingDate.UTC()
	return s.repository.PostWarehouseReceipt(ctx, tenant, s.newIDs(), value)
}

func (s *WarehouseService) CancelPutAway(ctx context.Context, tenant, organization, activity string, version int64) (WarehouseActivity, error) {
	if tenant == "" || organization == "" || activity == "" || version < 1 {
		return WarehouseActivity{}, fmt.Errorf("invalid warehouse put-away cancellation")
	}
	return s.repository.CancelWarehousePutAway(ctx, tenant, organization, activity, version, s.ids.New())
}

func (s *WarehouseService) CreatePick(ctx context.Context, tenant string, value WarehousePickCommand) (WarehouseActivity, error) {
	allowed := map[string]bool{"customer-order": true, "service": true, "transfer-outbound": true, "manual": true}
	if tenant == "" || value.RequestID == "" || value.OrganizationID == "" || value.ShipBinID == "" || value.ItemID == "" || !allowed[value.DemandKind] || value.DemandID == "" || value.DemandLineID == "" || !positiveDecimal(value.Quantity, quantityPattern) {
		return WarehouseActivity{}, fmt.Errorf("invalid warehouse pick")
	}
	return s.repository.CreateWarehousePick(ctx, tenant, s.newIDs(), value)
}

func (s *WarehouseService) Register(ctx context.Context, tenant, organization, activity string, version int64, postingDate time.Time) (WarehouseActivity, error) {
	if tenant == "" || organization == "" || activity == "" || version < 1 || postingDate.IsZero() {
		return WarehouseActivity{}, fmt.Errorf("invalid warehouse registration")
	}
	return s.repository.RegisterWarehouseActivity(ctx, tenant, organization, activity, version, postingDate.UTC(), s.ids.New())
}

func (s *WarehouseService) CancelPick(ctx context.Context, tenant, organization, activity string, version int64) (WarehouseActivity, error) {
	if tenant == "" || organization == "" || activity == "" || version < 1 {
		return WarehouseActivity{}, fmt.Errorf("invalid warehouse pick cancellation")
	}
	return s.repository.CancelWarehousePick(ctx, tenant, organization, activity, version, s.ids.New())
}
