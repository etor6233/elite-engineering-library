package operations

import (
	"context"
	"testing"
)

type fakeRepo struct{ po, factory, stock int }

func (f *fakeRepo) CreatePurchaseOrder(context.Context, string, string, PurchaseOrder) error {
	f.po++
	return nil
}
func (f *fakeRepo) TransitionPurchaseOrder(context.Context, string, string, string, string, int64, string, string) error {
	f.po++
	return nil
}
func (f *fakeRepo) CreateProductionUnit(context.Context, string, string, ProductionUnit) error {
	f.factory++
	return nil
}
func (f *fakeRepo) TransitionProductionUnit(context.Context, string, string, string, string, string, string) error {
	f.factory++
	return nil
}
func (f *fakeRepo) CreateStockUnit(context.Context, string, string, StockUnit) error {
	f.stock++
	return nil
}
func (f *fakeRepo) TransitionStockUnit(context.Context, string, string, string, string, int64, string, string) error {
	f.stock++
	return nil
}

type ids struct{ n int }

func (i *ids) New() string {
	i.n++
	return []string{"018f4d4a-7b36-7a21-8d10-2f4c54c28001", "018f4d4a-7b36-7a21-8d10-2f4c54c28002", "018f4d4a-7b36-7a21-8d10-2f4c54c28003", "018f4d4a-7b36-7a21-8d10-2f4c54c28004", "018f4d4a-7b36-7a21-8d10-2f4c54c28005", "018f4d4a-7b36-7a21-8d10-2f4c54c28006"}[i.n-1]
}
func TestStateMachinesRejectSkippedTransitions(t *testing.T) {
	repo := &fakeRepo{}
	service := NewService(repo, &ids{})
	ctx := context.Background()
	if _, err := service.CreatePurchaseOrder(ctx, "tenant", PurchaseOrder{SupplierID: "s", DestinationOrganizationID: "o", Currency: "USD", TotalMinorUnits: 1}); err != nil {
		t.Fatal(err)
	}
	if err := service.TransitionPurchaseOrder(ctx, "tenant", "o", "po", "draft", "received", 1); err == nil {
		t.Fatal("purchase order skipped states")
	}
	if err := service.TransitionProductionUnit(ctx, "tenant", "o", "u", "planned", "released"); err == nil {
		t.Fatal("factory skipped states")
	}
	if err := service.TransitionStockUnit(ctx, "tenant", "o", "s", "available", "sold", 1); err == nil {
		t.Fatal("stock skipped reservation")
	}
	if repo.po != 1 || repo.factory != 0 || repo.stock != 0 {
		t.Fatalf("unexpected repository calls %+v", repo)
	}
}
