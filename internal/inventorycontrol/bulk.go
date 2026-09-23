package inventorycontrol

import (
	"context"
	"fmt"
	"math/big"
	"regexp"
	"time"
)

var (
	quantityPattern = regexp.MustCompile(`^(0|[1-9][0-9]{0,13})(\.[0-9]{1,6})?$`)
	costPattern     = regexp.MustCompile(`^(0|[1-9][0-9]{0,15})(\.[0-9]{1,4})?$`)
)

type BulkItem struct {
	ID                    string `json:"id"`
	Code                  string `json:"code"`
	Description           string `json:"description"`
	BaseUOM               string `json:"base_uom"`
	BaseRoundingPrecision string `json:"base_rounding_precision"`
	TrackingMode          string `json:"tracking_mode"`
	CostingMethod         string `json:"costing_method"`
	Version               int64  `json:"version"`
}

type WarehouseBin struct {
	ID              string `json:"id"`
	OrganizationID  string `json:"organization_id"`
	Code            string `json:"code"`
	Type            string `json:"type"`
	Ranking         int    `json:"ranking"`
	MovementBlocked bool   `json:"movement_blocked"`
	CrossDock       bool   `json:"cross_dock"`
	Version         int64  `json:"version"`
}

type ItemBinPolicy struct {
	OrganizationID string `json:"organization_id"`
	ItemID         string `json:"item_id"`
	BinID          string `json:"bin_id"`
	Fixed          bool   `json:"fixed"`
	Dedicated      bool   `json:"dedicated"`
	Default        bool   `json:"default"`
	MinQuantity    string `json:"min_quantity"`
	MaxQuantity    string `json:"max_quantity,omitempty"`
	Version        int64  `json:"version"`
}

type BulkReceipt struct {
	OrganizationID string     `json:"organization_id"`
	BinID          string     `json:"bin_id"`
	ItemID         string     `json:"item_id"`
	LotNo          string     `json:"lot_no,omitempty"`
	ExpirationDate *time.Time `json:"expiration_date,omitempty"`
	WarrantyDate   *time.Time `json:"warranty_date,omitempty"`
	Quantity       string     `json:"quantity"`
	UnitCost       string     `json:"unit_cost"`
	PostingDate    time.Time  `json:"posting_date"`
	SourceKind     string     `json:"source_kind"`
	SourceID       string     `json:"source_id"`
}

type BulkReceiptResult struct {
	EntryID    string `json:"entry_id"`
	LotID      string `json:"lot_id,omitempty"`
	Quantity   string `json:"quantity"`
	UnitCost   string `json:"unit_cost"`
	CostAmount string `json:"cost_amount"`
}

type BulkAvailability struct {
	OrganizationID    string     `json:"organization_id"`
	BinID             string     `json:"bin_id"`
	BinCode           string     `json:"bin_code"`
	ItemID            string     `json:"item_id"`
	LotID             string     `json:"lot_id,omitempty"`
	LotNo             string     `json:"lot_no,omitempty"`
	ExpirationDate    *time.Time `json:"expiration_date,omitempty"`
	Quantity          string     `json:"quantity"`
	ReservedQuantity  string     `json:"reserved_quantity"`
	AvailableQuantity string     `json:"available_quantity"`
}

type BulkReservation struct {
	ID                     string     `json:"id"`
	OrganizationID         string     `json:"organization_id"`
	BinID                  string     `json:"bin_id"`
	ItemID                 string     `json:"item_id"`
	LotID                  string     `json:"lot_id,omitempty"`
	DemandKind             string     `json:"demand_kind"`
	DemandID               string     `json:"demand_id"`
	DemandLineID           string     `json:"demand_line_id"`
	Quantity               string     `json:"quantity"`
	CancellationDisallowed bool       `json:"cancellation_disallowed"`
	ExpiresAt              *time.Time `json:"expires_at,omitempty"`
	Status                 string     `json:"status"`
	Version                int64      `json:"version"`
}

type BulkMovement struct {
	OrganizationID string    `json:"organization_id"`
	FromBinID      string    `json:"from_bin_id"`
	ToBinID        string    `json:"to_bin_id"`
	ItemID         string    `json:"item_id"`
	LotID          string    `json:"lot_id,omitempty"`
	Quantity       string    `json:"quantity"`
	PostingDate    time.Time `json:"posting_date"`
	SourceID       string    `json:"source_id"`
}

type BulkIssue struct {
	OrganizationID       string    `json:"organization_id"`
	ReservationID        string    `json:"reservation_id"`
	ReservationVersion   int64     `json:"reservation_version"`
	SpecificReceiptEntry string    `json:"specific_receipt_entry,omitempty"`
	PostingDate          time.Time `json:"posting_date"`
	SourceKind           string    `json:"source_kind"`
	SourceID             string    `json:"source_id"`
}

type BulkIssueResult struct {
	EntryID      string `json:"entry_id"`
	Quantity     string `json:"quantity"`
	CostAmount   string `json:"cost_amount"`
	Applications int    `json:"applications"`
}

type BulkRepository interface {
	CreateBulkItem(context.Context, string, string, BulkItem) (BulkItem, error)
	ConfigureItemUnitOfMeasure(context.Context, string, string, ItemUnitOfMeasure) (ItemUnitOfMeasure, error)
	ConvertItemUnitOfMeasure(context.Context, string, string, string, string) (UnitOfMeasureConversion, error)
	ConvertBulkHandlingUnits(context.Context, string, PackagingConversionIDs, PackagingConversionCommand) (PackagingConversionResult, error)
	CreateWarehouseBin(context.Context, string, string, WarehouseBin) (WarehouseBin, error)
	ConfigureItemBin(context.Context, string, string, ItemBinPolicy) (ItemBinPolicy, error)
	ReceiveBulk(context.Context, string, string, string, string, string, BulkReceipt) (BulkReceiptResult, error)
	BulkAvailability(context.Context, string, string, string) ([]BulkAvailability, error)
	ReserveBulk(context.Context, string, string, BulkReservation) (BulkReservation, error)
	ReleaseBulk(context.Context, string, string, string, int64, string) error
	MoveBulk(context.Context, string, string, string, BulkMovement) error
	IssueBulk(context.Context, string, string, string, BulkIssue) (BulkIssueResult, error)
}

type BulkService struct {
	repository BulkRepository
	ids        IDGenerator
}

func NewBulkService(repository BulkRepository, ids IDGenerator) *BulkService {
	return &BulkService{repository: repository, ids: ids}
}

func positiveDecimal(value string, pattern *regexp.Regexp) bool {
	if !pattern.MatchString(value) {
		return false
	}
	n, ok := new(big.Rat).SetString(value)
	return ok && n.Sign() > 0
}

func nonNegativeDecimal(value string) bool {
	if !quantityPattern.MatchString(value) {
		return false
	}
	n, ok := new(big.Rat).SetString(value)
	return ok && n.Sign() >= 0
}

func (s *BulkService) CreateItem(ctx context.Context, tenant string, value BulkItem) (BulkItem, error) {
	tracking := map[string]bool{"none": true, "lot": true}
	costing := map[string]bool{"fifo": true, "specific": true}
	if tenant == "" || value.Code == "" || value.Description == "" || value.BaseUOM == "" || !positiveDecimal(value.BaseRoundingPrecision, quantityPattern) || !tracking[value.TrackingMode] || !costing[value.CostingMethod] {
		return BulkItem{}, fmt.Errorf("invalid bulk item")
	}
	value.ID, value.Version = s.ids.New(), 1
	return s.repository.CreateBulkItem(ctx, tenant, s.ids.New(), value)
}

func (s *BulkService) CreateBin(ctx context.Context, tenant string, value WarehouseBin) (WarehouseBin, error) {
	types := map[string]bool{"receive": true, "ship": true, "put-away": true, "pick": true, "putpick": true, "qc": true}
	if tenant == "" || value.OrganizationID == "" || value.Code == "" || !types[value.Type] || value.Ranking < 0 || value.Ranking > 1_000_000 {
		return WarehouseBin{}, fmt.Errorf("invalid warehouse bin")
	}
	value.ID, value.Version = s.ids.New(), 1
	return s.repository.CreateWarehouseBin(ctx, tenant, s.ids.New(), value)
}

func (s *BulkService) ConfigureBin(ctx context.Context, tenant string, value ItemBinPolicy) (ItemBinPolicy, error) {
	if tenant == "" || value.OrganizationID == "" || value.ItemID == "" || value.BinID == "" || !nonNegativeDecimal(value.MinQuantity) || (value.MaxQuantity != "" && !positiveDecimal(value.MaxQuantity, quantityPattern)) {
		return ItemBinPolicy{}, fmt.Errorf("invalid item bin policy")
	}
	value.Version = 1
	return s.repository.ConfigureItemBin(ctx, tenant, s.ids.New(), value)
}

func (s *BulkService) Receive(ctx context.Context, tenant string, value BulkReceipt) (BulkReceiptResult, error) {
	if tenant == "" || value.OrganizationID == "" || value.BinID == "" || value.ItemID == "" || !positiveDecimal(value.Quantity, quantityPattern) || !positiveDecimal(value.UnitCost, costPattern) || value.PostingDate.IsZero() || value.SourceKind == "" || value.SourceID == "" {
		return BulkReceiptResult{}, fmt.Errorf("invalid bulk receipt")
	}
	value.PostingDate = value.PostingDate.UTC()
	return s.repository.ReceiveBulk(ctx, tenant, s.ids.New(), s.ids.New(), s.ids.New(), s.ids.New(), value)
}

func (s *BulkService) Availability(ctx context.Context, tenant, organization, item string) ([]BulkAvailability, error) {
	if tenant == "" || organization == "" || item == "" {
		return nil, fmt.Errorf("invalid bulk availability query")
	}
	return s.repository.BulkAvailability(ctx, tenant, organization, item)
}

func (s *BulkService) Reserve(ctx context.Context, tenant string, value BulkReservation) (BulkReservation, error) {
	allowed := map[string]bool{"service": true, "manual": true, "transfer-outbound": true}
	if tenant == "" || value.OrganizationID == "" || value.BinID == "" || value.ItemID == "" || !allowed[value.DemandKind] || value.DemandID == "" || value.DemandLineID == "" || !positiveDecimal(value.Quantity, quantityPattern) {
		return BulkReservation{}, fmt.Errorf("invalid bulk reservation")
	}
	value.ID, value.Status, value.Version = s.ids.New(), "reservation", 1
	return s.repository.ReserveBulk(ctx, tenant, s.ids.New(), value)
}

func (s *BulkService) Release(ctx context.Context, tenant, organization, reservation string, version int64) error {
	if tenant == "" || organization == "" || reservation == "" || version < 1 {
		return fmt.Errorf("invalid bulk release")
	}
	return s.repository.ReleaseBulk(ctx, tenant, organization, reservation, version, s.ids.New())
}

func (s *BulkService) Move(ctx context.Context, tenant string, value BulkMovement) error {
	if tenant == "" || value.OrganizationID == "" || value.FromBinID == "" || value.ToBinID == "" || value.FromBinID == value.ToBinID || value.ItemID == "" || !positiveDecimal(value.Quantity, quantityPattern) || value.PostingDate.IsZero() || value.SourceID == "" {
		return fmt.Errorf("invalid bulk movement")
	}
	value.PostingDate = value.PostingDate.UTC()
	return s.repository.MoveBulk(ctx, tenant, s.ids.New(), s.ids.New(), value)
}

func (s *BulkService) Issue(ctx context.Context, tenant string, value BulkIssue) (BulkIssueResult, error) {
	if tenant == "" || value.OrganizationID == "" || value.ReservationID == "" || value.ReservationVersion < 1 || value.PostingDate.IsZero() || value.SourceKind == "" || value.SourceID == "" {
		return BulkIssueResult{}, fmt.Errorf("invalid bulk issue")
	}
	value.PostingDate = value.PostingDate.UTC()
	return s.repository.IssueBulk(ctx, tenant, s.ids.New(), s.ids.New(), value)
}
