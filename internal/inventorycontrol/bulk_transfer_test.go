package inventorycontrol

import (
	"context"
	"testing"
	"time"
)

type transferRepo struct{ calls int }

func (r *transferRepo) CreateBulkTransfer(context.Context, string, string, string, string, BulkTransferCommand) (BulkTransfer, error) {
	r.calls++
	return BulkTransfer{Status: "released"}, nil
}
func (r *transferRepo) ShipBulkTransfer(context.Context, string, string, string, string, int64, string, string, BulkTransferPostingCommand) (BulkTransfer, error) {
	r.calls++
	return BulkTransfer{Status: "shipped"}, nil
}
func (r *transferRepo) ReceiveBulkTransfer(context.Context, string, string, string, string, int64, string, string, string, BulkTransferPostingCommand) (BulkTransfer, error) {
	r.calls++
	return BulkTransfer{Status: "received"}, nil
}
func (r *transferRepo) CancelBulkTransfer(context.Context, string, string, string, string, int64, string) (BulkTransfer, error) {
	r.calls++
	return BulkTransfer{Status: "cancelled"}, nil
}

type transferIDs struct{ n int }

func (g *transferIDs) New() string { g.n++; return "id" }

func TestBulkTransferServiceRejectsAmbiguityBeforeRepository(t *testing.T) {
	repo := &transferRepo{}
	service := NewBulkTransferService(repo, &transferIDs{})
	_, err := service.Create(context.Background(), "tenant", BulkTransferCommand{RequestID: "request", FromOrganizationID: "same", ToOrganizationID: "same", ReceiveBinID: "receive", InTransitCode: "TRANSIT", ItemID: "item", Quantity: "1", PostingDate: time.Now(), SourceKind: "plan", SourceID: "source"})
	if err == nil || repo.calls != 0 {
		t.Fatalf("expected fail-closed validation, err=%v calls=%d", err, repo.calls)
	}
	_, err = service.Create(context.Background(), "tenant", BulkTransferCommand{RequestID: "request", FromOrganizationID: "from", ToOrganizationID: "to", ReceiveBinID: "receive", InTransitCode: "TRANSIT", ItemID: "item", Quantity: "1.250000", PostingDate: time.Now(), SourceKind: "plan", SourceID: "source"})
	if err != nil || repo.calls != 1 {
		t.Fatalf("expected valid transfer, err=%v calls=%d", err, repo.calls)
	}
	_, err = service.Ship(context.Background(), "tenant", "transfer", "from", "to", 1, BulkTransferPostingCommand{RequestID: "ship", Quantity: "1", PostingDate: time.Now()})
	if err == nil || repo.calls != 1 {
		t.Fatalf("shipment without exact pick must fail, err=%v calls=%d", err, repo.calls)
	}
	_, err = service.Ship(context.Background(), "tenant", "transfer", "from", "to", 1, BulkTransferPostingCommand{RequestID: "ship", Quantity: "1", WarehouseActivityID: "pick", PostingDate: time.Now()})
	if err != nil || repo.calls != 2 {
		t.Fatalf("valid shipment command, err=%v calls=%d", err, repo.calls)
	}
	_, err = service.Receive(context.Background(), "tenant", "transfer", "from", "to", 2, BulkTransferPostingCommand{RequestID: "receipt", Quantity: "1", WarehouseActivityID: "unexpected", PostingDate: time.Now()})
	if err == nil || repo.calls != 2 {
		t.Fatalf("receipt with warehouse activity must fail, err=%v calls=%d", err, repo.calls)
	}
}
