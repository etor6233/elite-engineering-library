package inventorycontrol

import (
	"context"
	"testing"
)

func TestUnitOfMeasureServiceFailsClosedBeforeRepository(t *testing.T) {
	repo := &bulkFake{}
	service := NewBulkService(repo, inventoryIDsForBulk{})
	invalid := []ItemUnitOfMeasure{
		{},
		{ItemID: "part", Code: "BOX", QuantityPerUnit: "12", RoundingPrecision: "1", Base: true},
		{ItemID: "part", Code: "BOX", QuantityPerUnit: "0", RoundingPrecision: "1"},
		{ItemID: "part", Code: "BOX", QuantityPerUnit: "12", RoundingPrecision: "0"},
	}
	for _, value := range invalid {
		if _, err := service.ConfigureUnitOfMeasure(context.Background(), "tenant", value); err == nil {
			t.Fatalf("accepted invalid UOM: %+v", value)
		}
	}
	if _, err := service.ConvertUnitOfMeasure(context.Background(), "tenant", "part", "BOX", "0"); err == nil {
		t.Fatal("accepted zero conversion")
	}
}
