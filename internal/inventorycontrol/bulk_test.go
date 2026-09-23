package inventorycontrol

import (
	"context"
	"testing"
	"time"
)

type bulkFake struct{ created int }

func (*bulkFake) ConfigureItemUnitOfMeasure(context.Context, string, string, ItemUnitOfMeasure) (ItemUnitOfMeasure, error) {
	return ItemUnitOfMeasure{}, nil
}
func (*bulkFake) ConvertItemUnitOfMeasure(context.Context, string, string, string, string) (UnitOfMeasureConversion, error) {
	return UnitOfMeasureConversion{}, nil
}
func (*bulkFake) ConvertBulkHandlingUnits(context.Context, string, PackagingConversionIDs, PackagingConversionCommand) (PackagingConversionResult, error) {
	return PackagingConversionResult{}, nil
}

func (f *bulkFake) CreateBulkItem(_ context.Context, _ string, _ string, v BulkItem) (BulkItem, error) {
	f.created++
	return v, nil
}
func (*bulkFake) CreateWarehouseBin(context.Context, string, string, WarehouseBin) (WarehouseBin, error) {
	return WarehouseBin{}, nil
}
func (*bulkFake) ConfigureItemBin(context.Context, string, string, ItemBinPolicy) (ItemBinPolicy, error) {
	return ItemBinPolicy{}, nil
}
func (*bulkFake) ReceiveBulk(context.Context, string, string, string, string, string, BulkReceipt) (BulkReceiptResult, error) {
	return BulkReceiptResult{}, nil
}
func (*bulkFake) BulkAvailability(context.Context, string, string, string) ([]BulkAvailability, error) {
	return nil, nil
}
func (*bulkFake) ReserveBulk(context.Context, string, string, BulkReservation) (BulkReservation, error) {
	return BulkReservation{}, nil
}
func (*bulkFake) ReleaseBulk(context.Context, string, string, string, int64, string) error {
	return nil
}
func (*bulkFake) MoveBulk(context.Context, string, string, string, BulkMovement) error { return nil }
func (*bulkFake) IssueBulk(context.Context, string, string, string, BulkIssue) (BulkIssueResult, error) {
	return BulkIssueResult{}, nil
}

func TestBulkServiceRejectsUnsupportedCostAndSerialDuplication(t *testing.T) {
	fake := &bulkFake{}
	service := NewBulkService(fake, inventoryIDsForBulk{})
	for _, value := range []BulkItem{
		{Code: "PART", Description: "Part", BaseUOM: "EA", TrackingMode: "serial", CostingMethod: "specific"},
		{Code: "PART", Description: "Part", BaseUOM: "EA", TrackingMode: "lot", CostingMethod: "average"},
		{Code: "PART", Description: "Part", BaseUOM: "EA", TrackingMode: "lot", CostingMethod: "lifo"},
	} {
		if _, err := service.CreateItem(context.Background(), "tenant", value); err == nil {
			t.Fatalf("unsupported item admitted: %+v", value)
		}
	}
	if fake.created != 0 {
		t.Fatalf("repository called %d times", fake.created)
	}
	if _, err := service.Receive(context.Background(), "tenant", BulkReceipt{OrganizationID: "o", BinID: "b", ItemID: "i", Quantity: "1.0000001", UnitCost: "1", PostingDate: time.Now(), SourceKind: "purchase", SourceID: "p"}); err == nil {
		t.Fatal("over-precision quantity admitted")
	}
}

func TestBulkServiceValidatesPackagingConversionBeforeRepository(t *testing.T) {
	service := NewBulkService(&bulkFake{}, inventoryIDsForBulk{})
	valid := PackagingConversionCommand{RequestID: "request-1", OrganizationID: "warehouse", BinID: "pick", ItemID: "part", Operation: "breakbulk", FromUOM: "BOX", ToUOM: "EA", FromQuantity: "1"}
	if _, err := service.ConvertHandlingUnits(context.Background(), "tenant", valid); err != nil {
		t.Fatal(err)
	}
	for _, invalid := range []PackagingConversionCommand{
		{},
		{RequestID: "r", OrganizationID: "o", BinID: "b", ItemID: "i", Operation: "invented", FromUOM: "BOX", ToUOM: "EA", FromQuantity: "1"},
		{RequestID: "r", OrganizationID: "o", BinID: "b", ItemID: "i", Operation: "gather", FromUOM: "EA", ToUOM: "EA", FromQuantity: "1"},
		{RequestID: "r", OrganizationID: "o", BinID: "b", ItemID: "i", Operation: "gather", FromUOM: "EA", ToUOM: "BOX", FromQuantity: "0"},
	} {
		if _, err := service.ConvertHandlingUnits(context.Background(), "tenant", invalid); err == nil {
			t.Fatalf("invalid packaging conversion admitted: %+v", invalid)
		}
	}
}

type inventoryIDsForBulk struct{ n int }

func (g inventoryIDsForBulk) New() string { return "018f4d4a-7b36-7a21-8d10-2f4c54c28888" }
