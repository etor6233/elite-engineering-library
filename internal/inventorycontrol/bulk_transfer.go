package inventorycontrol

import (
	"context"
	"fmt"
	"time"
)

type BulkTransferCommand struct {
	RequestID            string    `json:"request_id"`
	FromOrganizationID   string    `json:"from_organization_id"`
	ToOrganizationID     string    `json:"to_organization_id"`
	ReceiveBinID         string    `json:"receive_bin_id"`
	InTransitCode        string    `json:"in_transit_code"`
	ItemID               string    `json:"item_id"`
	Quantity             string    `json:"quantity"`
	SpecificReceiptEntry string    `json:"specific_receipt_entry,omitempty"`
	PostingDate          time.Time `json:"posting_date"`
	SourceKind           string    `json:"source_kind"`
	SourceID             string    `json:"source_id"`
}

type BulkTransfer struct {
	ID                   string    `json:"id"`
	LineID               string    `json:"line_id"`
	RequestID            string    `json:"request_id"`
	FromOrganizationID   string    `json:"from_organization_id"`
	ToOrganizationID     string    `json:"to_organization_id"`
	ReceiveBinID         string    `json:"receive_bin_id"`
	InTransitCode        string    `json:"in_transit_code"`
	ItemID               string    `json:"item_id"`
	Quantity             string    `json:"quantity"`
	SpecificReceiptEntry string    `json:"specific_receipt_entry,omitempty"`
	ShippedQuantity      string    `json:"shipped_quantity"`
	ReceivedQuantity     string    `json:"received_quantity"`
	CostAmount           string    `json:"cost_amount"`
	PostingID            string    `json:"posting_id,omitempty"`
	PutAwayID            string    `json:"put_away_id,omitempty"`
	Status               string    `json:"status"`
	Version              int64     `json:"version"`
	PostingDate          time.Time `json:"posting_date"`
}

type BulkTransferPostingCommand struct {
	RequestID           string    `json:"request_id"`
	Quantity            string    `json:"quantity"`
	WarehouseActivityID string    `json:"warehouse_activity_id,omitempty"`
	PostingDate         time.Time `json:"posting_date"`
}

type BulkTransferRepository interface {
	CreateBulkTransfer(context.Context, string, string, string, string, BulkTransferCommand) (BulkTransfer, error)
	ShipBulkTransfer(context.Context, string, string, string, string, int64, string, string, BulkTransferPostingCommand) (BulkTransfer, error)
	ReceiveBulkTransfer(context.Context, string, string, string, string, int64, string, string, string, BulkTransferPostingCommand) (BulkTransfer, error)
	CancelBulkTransfer(context.Context, string, string, string, string, int64, string) (BulkTransfer, error)
}

type BulkTransferService struct {
	repository BulkTransferRepository
	ids        IDGenerator
}

func NewBulkTransferService(repository BulkTransferRepository, ids IDGenerator) *BulkTransferService {
	return &BulkTransferService{repository: repository, ids: ids}
}

func (s *BulkTransferService) Create(ctx context.Context, tenant string, value BulkTransferCommand) (BulkTransfer, error) {
	if tenant == "" || value.RequestID == "" || value.FromOrganizationID == "" || value.ToOrganizationID == "" || value.FromOrganizationID == value.ToOrganizationID || value.ReceiveBinID == "" || value.InTransitCode == "" || value.ItemID == "" || !positiveDecimal(value.Quantity, quantityPattern) || value.PostingDate.IsZero() || value.SourceKind == "" || value.SourceID == "" {
		return BulkTransfer{}, fmt.Errorf("invalid bulk transfer")
	}
	value.PostingDate = value.PostingDate.UTC()
	return s.repository.CreateBulkTransfer(ctx, tenant, s.ids.New(), s.ids.New(), s.ids.New(), value)
}

func (s *BulkTransferService) Ship(ctx context.Context, tenant, transfer, fromOrganization, toOrganization string, version int64, command BulkTransferPostingCommand) (BulkTransfer, error) {
	if tenant == "" || transfer == "" || fromOrganization == "" || toOrganization == "" || fromOrganization == toOrganization || version < 1 || command.RequestID == "" || !positiveDecimal(command.Quantity, quantityPattern) || command.WarehouseActivityID == "" || command.PostingDate.IsZero() {
		return BulkTransfer{}, fmt.Errorf("invalid bulk transfer shipment")
	}
	command.PostingDate = command.PostingDate.UTC()
	return s.repository.ShipBulkTransfer(ctx, tenant, transfer, fromOrganization, toOrganization, version, s.ids.New(), s.ids.New(), command)
}

func (s *BulkTransferService) Receive(ctx context.Context, tenant, transfer, fromOrganization, toOrganization string, version int64, command BulkTransferPostingCommand) (BulkTransfer, error) {
	if tenant == "" || transfer == "" || fromOrganization == "" || toOrganization == "" || fromOrganization == toOrganization || version < 1 || command.RequestID == "" || !positiveDecimal(command.Quantity, quantityPattern) || command.WarehouseActivityID != "" || command.PostingDate.IsZero() {
		return BulkTransfer{}, fmt.Errorf("invalid bulk transfer receipt")
	}
	command.PostingDate = command.PostingDate.UTC()
	return s.repository.ReceiveBulkTransfer(ctx, tenant, transfer, fromOrganization, toOrganization, version, s.ids.New(), s.ids.New(), s.ids.New(), command)
}

func (s *BulkTransferService) Cancel(ctx context.Context, tenant, transfer, fromOrganization, toOrganization string, version int64) (BulkTransfer, error) {
	if tenant == "" || transfer == "" || fromOrganization == "" || toOrganization == "" || fromOrganization == toOrganization || version < 1 {
		return BulkTransfer{}, fmt.Errorf("invalid bulk transfer cancellation")
	}
	return s.repository.CancelBulkTransfer(ctx, tenant, transfer, fromOrganization, toOrganization, version, s.ids.New())
}
