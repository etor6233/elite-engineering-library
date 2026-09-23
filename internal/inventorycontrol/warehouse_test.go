package inventorycontrol

import (
	"context"
	"testing"
	"time"
)

type warehouseIDs struct{ n int }

func (g *warehouseIDs) New() string {
	g.n++
	return "id-" + time.Unix(int64(g.n), 0).UTC().Format("150405")
}

type warehouseRepo struct {
	bindings             int
	shipments            int
	receipts             int
	picks                int
	replenishments       int
	cancellations        int
	putAwayCancellations int
	crossDock            int
}

func (r *warehouseRepo) ConfigureSalesWarehouseBinding(_ context.Context, _ string, _ string, value SalesWarehouseBinding) (SalesWarehouseBinding, error) {
	r.bindings++
	return value, nil
}
func (r *warehouseRepo) PostCustomerShipment(_ context.Context, _ string, ids CustomerShipmentIDs, value CustomerShipmentCommand) (CustomerShipment, error) {
	r.shipments++
	return CustomerShipment{ID: ids.ShipmentID, RequestID: value.RequestID, OrganizationID: value.OrganizationID, OrderID: value.OrderID, WarehouseActivityID: value.WarehouseActivityID, PostingDate: value.PostingDate, Line: CustomerShipmentLine{ID: ids.LineID}}, nil
}

func (r *warehouseRepo) PostWarehouseReceipt(_ context.Context, _ string, ids WarehouseIDs, value WarehouseReceiptCommand) (WarehouseReceiptResult, error) {
	r.receipts++
	return WarehouseReceiptResult{ReceiptID: ids.ReceiptID, Quantity: value.Quantity}, nil
}
func (r *warehouseRepo) CreateWarehousePick(_ context.Context, _ string, ids WarehouseIDs, value WarehousePickCommand) (WarehouseActivity, error) {
	r.picks++
	return WarehouseActivity{ID: ids.ActivityID, OrganizationID: value.OrganizationID, Type: "pick", Version: 1}, nil
}
func (r *warehouseRepo) CreateWarehouseReplenishment(_ context.Context, _ string, ids WarehouseIDs, value WarehouseReplenishmentCommand) (WarehouseActivity, error) {
	r.replenishments++
	return WarehouseActivity{ID: ids.ActivityID, OrganizationID: value.OrganizationID, Type: "movement", Version: 1}, nil
}
func (r *warehouseRepo) ConfigureWarehouseCrossDock(_ context.Context, _ string, _ string, value WarehouseCrossDockPolicy) (WarehouseCrossDockPolicy, error) {
	r.crossDock++
	return value, nil
}
func (*warehouseRepo) RegisterWarehouseActivity(context.Context, string, string, string, int64, time.Time, string) (WarehouseActivity, error) {
	return WarehouseActivity{}, nil
}
func (*warehouseRepo) CancelWarehousePick(context.Context, string, string, string, int64, string) (WarehouseActivity, error) {
	return WarehouseActivity{}, nil
}
func (r *warehouseRepo) CancelWarehousePutAway(context.Context, string, string, string, int64, string) (WarehouseActivity, error) {
	r.putAwayCancellations++
	return WarehouseActivity{}, nil
}
func (r *warehouseRepo) CancelWarehouseReplenishment(context.Context, string, string, string, int64, string) (WarehouseActivity, error) {
	r.cancellations++
	return WarehouseActivity{}, nil
}

func TestWarehouseServiceRejectsAmbiguityBeforeRepository(t *testing.T) {
	repo := &warehouseRepo{}
	service := NewWarehouseService(repo, &warehouseIDs{})
	if _, err := service.ConfigureSalesBinding(context.Background(), "tenant", SalesWarehouseBinding{RequestID: "binding-request", VariantID: "variant", ItemID: "part"}); err == nil || repo.bindings != 0 {
		t.Fatalf("invalid sales binding reached repository: err=%v calls=%d", err, repo.bindings)
	}
	binding, err := service.ConfigureSalesBinding(context.Background(), "tenant", SalesWarehouseBinding{RequestID: "binding-request", VariantID: "variant", ItemID: "part", SalesUOMCode: "BOX"})
	if err != nil || repo.bindings != 1 || binding.Version != 1 {
		t.Fatalf("valid sales binding failed: value=%+v err=%v calls=%d", binding, err, repo.bindings)
	}
	shipment := CustomerShipmentCommand{RequestID: "shipment-request", OrganizationID: "warehouse", OrderID: "order", WarehouseActivityID: "pick", PostingDate: time.Now()}
	invalidShipment := shipment
	invalidShipment.WarehouseActivityID = ""
	if _, err := service.PostCustomerShipment(context.Background(), "tenant", invalidShipment); err == nil || repo.shipments != 0 {
		t.Fatalf("invalid customer shipment reached repository: err=%v calls=%d", err, repo.shipments)
	}
	posted, err := service.PostCustomerShipment(context.Background(), "tenant", shipment)
	if err != nil || repo.shipments != 1 || posted.ID == "" || posted.Line.ID == "" || !posted.PostingDate.Equal(shipment.PostingDate.UTC()) {
		t.Fatalf("valid customer shipment failed: value=%+v err=%v calls=%d", posted, err, repo.shipments)
	}
	validReceipt := WarehouseReceiptCommand{RequestID: "receipt-request", OrganizationID: "warehouse", ReceiveBinID: "receive", ItemID: "part", Quantity: "1", UnitCost: "10", PostingDate: time.Now(), SourceKind: "purchase", SourceID: "po-1"}
	invalidReceipt := validReceipt
	invalidReceipt.Quantity = "1e3"
	if _, err := service.PostReceipt(context.Background(), "tenant", invalidReceipt); err == nil || repo.receipts != 0 {
		t.Fatalf("invalid receipt reached repository: err=%v calls=%d", err, repo.receipts)
	}
	if _, err := service.PostReceipt(context.Background(), "tenant", validReceipt); err != nil || repo.receipts != 1 {
		t.Fatalf("valid receipt failed: err=%v calls=%d", err, repo.receipts)
	}
	packagedReceipt := validReceipt
	packagedReceipt.Quantity = ""
	packagedReceipt.HandlingUOM = "BOX"
	packagedReceipt.HandlingQuantity = "2"
	if _, err := service.PostReceipt(context.Background(), "tenant", packagedReceipt); err != nil || repo.receipts != 2 {
		t.Fatalf("valid packaged receipt failed: err=%v calls=%d", err, repo.receipts)
	}
	ambiguousReceipt := packagedReceipt
	ambiguousReceipt.Quantity = "24"
	if _, err := service.PostReceipt(context.Background(), "tenant", ambiguousReceipt); err == nil || repo.receipts != 2 {
		t.Fatalf("ambiguous receipt reached repository: err=%v calls=%d", err, repo.receipts)
	}
	if _, err := service.CancelPutAway(context.Background(), "tenant", "warehouse", "put-away", 0); err == nil || repo.putAwayCancellations != 0 {
		t.Fatalf("invalid put-away cancellation reached repository: err=%v calls=%d", err, repo.putAwayCancellations)
	}
	if _, err := service.CancelPutAway(context.Background(), "tenant", "warehouse", "put-away", 1); err != nil || repo.putAwayCancellations != 1 {
		t.Fatalf("valid put-away cancellation failed: err=%v calls=%d", err, repo.putAwayCancellations)
	}
	invalidPick := WarehousePickCommand{RequestID: "pick-request", OrganizationID: "warehouse", ShipBinID: "ship", ItemID: "part", DemandKind: "invented", DemandID: "order", DemandLineID: "line", Quantity: "1", UseFEFO: true}
	if _, err := service.CreatePick(context.Background(), "tenant", invalidPick); err == nil || repo.picks != 0 {
		t.Fatalf("invalid pick reached repository: err=%v calls=%d", err, repo.picks)
	}
	invalidPick.DemandKind = "customer-order"
	if _, err := service.CreatePick(context.Background(), "tenant", invalidPick); err != nil || repo.picks != 1 {
		t.Fatalf("valid pick failed: err=%v calls=%d", err, repo.picks)
	}
}

func TestWarehouseReplenishmentRejectsAmbiguityBeforeRepository(t *testing.T) {
	repo := &warehouseRepo{}
	service := NewWarehouseService(repo, &warehouseIDs{})
	value := WarehouseReplenishmentCommand{RequestID: "replenish-request", OrganizationID: "warehouse", ToBinID: "pick", ItemID: "part", UseFEFO: true}
	invalid := value
	invalid.ToBinID = ""
	if _, err := service.CreateReplenishment(context.Background(), "tenant", invalid); err == nil || repo.replenishments != 0 {
		t.Fatalf("invalid replenishment reached repository: err=%v calls=%d", err, repo.replenishments)
	}
	if _, err := service.CreateReplenishment(context.Background(), "tenant", value); err != nil || repo.replenishments != 1 {
		t.Fatalf("valid replenishment failed: err=%v calls=%d", err, repo.replenishments)
	}
	if _, err := service.CancelReplenishment(context.Background(), "tenant", "warehouse", "activity", 0); err == nil || repo.cancellations != 0 {
		t.Fatalf("invalid cancellation reached repository: err=%v calls=%d", err, repo.cancellations)
	}
	if _, err := service.CancelReplenishment(context.Background(), "tenant", "warehouse", "activity", 1); err != nil || repo.cancellations != 1 {
		t.Fatalf("valid cancellation failed: err=%v calls=%d", err, repo.cancellations)
	}
}

func TestWarehouseCrossDockRejectsAmbiguityBeforeRepository(t *testing.T) {
	repo := &warehouseRepo{}
	service := NewWarehouseService(repo, &warehouseIDs{})
	value := WarehouseCrossDockPolicy{OrganizationID: "warehouse", ItemID: "part", BinID: "cross-dock", DueDateDays: 3, Enabled: true, Version: 1}
	invalid := value
	invalid.DueDateDays = 366
	if _, err := service.ConfigureCrossDock(context.Background(), "tenant", invalid); err == nil || repo.crossDock != 0 {
		t.Fatalf("invalid cross-dock policy reached repository: err=%v calls=%d", err, repo.crossDock)
	}
	if _, err := service.ConfigureCrossDock(context.Background(), "tenant", value); err != nil || repo.crossDock != 1 {
		t.Fatalf("valid cross-dock policy failed: err=%v calls=%d", err, repo.crossDock)
	}
}
