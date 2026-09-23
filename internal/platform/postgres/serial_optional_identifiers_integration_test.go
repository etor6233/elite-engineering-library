// AUTHORED regression for optional identifiers in the existing serial owner.
// Synthetic fixtures; no connected J2 journey is claimed by this focused test.
package postgres

import (
	"context"
	"elite.local/enterprise/internal/operations"
	"elite.local/enterprise/internal/platform/randomid"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"testing"
)

func TestSerialSupplyOptionalIdentifiers(t *testing.T) {
	url := os.Getenv("PAYMENT_CONNECTED_DB_URL")
	if url == "" {
		t.Skip("owned fixture DB not configured")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := "ca6c9cb7-e548-4f19-a2e6-27d2aa2d4890"
	fixtures := []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,'supply-optionals','Fixture','Fixture')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,'supply-org','supply-org','Fixture','warehouse')`,
		`insert into partner.supplier(tenant_id,supplier_id,supplier_code,legal_name,status) values($1,'supplier','supplier','Fixture','active')`,
		`insert into catalog.vehicle_model(tenant_id,model_id,model_code,display_name,vehicle_class,lifecycle_state) values($1,'model','model','Fixture','bicycle','active')`,
		`insert into catalog.vehicle_variant(tenant_id,variant_id,model_id,variant_code,display_name,battery_specification,lifecycle_state) values($1,'variant','model','variant','Fixture','{}','active')`,
	}
	for _, q := range fixtures {
		if _, err = pool.Exec(ctx, q, tenant); err != nil {
			t.Fatal(err)
		}
	}
	svc := operations.NewService(NewOperations(pool), randomid.Generator{})
	po, err := svc.CreatePurchaseOrder(ctx, tenant, operations.PurchaseOrder{SupplierID: "supplier", DestinationOrganizationID: "supply-org", Currency: "ARS", TotalMinorUnits: 100})
	if err != nil {
		t.Fatal(err)
	}
	for _, serial := range []string{"FRAME-1", "FRAME-2"} {
		unit, err := svc.RegisterProductionUnit(ctx, tenant, operations.ProductionUnit{OrganizationID: "supply-org", PurchaseOrderID: po.ID, VariantID: "variant", SerialNumber: serial})
		if err != nil {
			t.Fatalf("optional VIN/battery cannot exclude second valid unit %s: %v", serial, err)
		}
		if _, err = svc.ReceiveStockUnit(ctx, tenant, operations.StockUnit{OrganizationID: "supply-org", VariantID: "variant", ProductionUnitID: unit.ID, SerialNumber: serial}); err != nil {
			t.Fatal(err)
		}
	}
	input := operations.ProductionUnit{OrganizationID: "supply-org", PurchaseOrderID: po.ID, VariantID: "variant", SerialNumber: "FRAME-3", VIN: "VIN-REAL", BatterySerialNumber: "BATTERY-REAL"}
	specific, err := svc.RegisterProductionUnit(ctx, tenant, input)
	if err != nil {
		t.Fatal(err)
	}
	specificStock := operations.StockUnit{OrganizationID: "supply-org", VariantID: "variant", ProductionUnitID: specific.ID, SerialNumber: input.SerialNumber, VIN: input.VIN, BatterySerialNumber: input.BatterySerialNumber}
	if _, err = svc.ReceiveStockUnit(ctx, tenant, specificStock); err != nil {
		t.Fatal(err)
	}
	specificStock.ProductionUnitID = ""
	specificStock.SerialNumber = "FRAME-STOCK-4"
	specificStock.BatterySerialNumber = "BATTERY-STOCK-OTHER"
	if _, err = svc.ReceiveStockUnit(ctx, tenant, specificStock); err == nil {
		t.Fatal("duplicate stock VIN admitted")
	}
	specificStock.VIN = "VIN-STOCK-OTHER"
	specificStock.BatterySerialNumber = "BATTERY-REAL"
	if _, err = svc.ReceiveStockUnit(ctx, tenant, specificStock); err == nil {
		t.Fatal("duplicate stock battery admitted")
	}
	specificStock.SerialNumber = "FRAME-1"
	specificStock.VIN = ""
	specificStock.BatterySerialNumber = ""
	if _, err = svc.ReceiveStockUnit(ctx, tenant, specificStock); err == nil {
		t.Fatal("duplicate stock mandatory serial admitted")
	}
	input.SerialNumber = "FRAME-4"
	input.BatterySerialNumber = "BATTERY-OTHER"
	if _, err = svc.RegisterProductionUnit(ctx, tenant, input); err == nil {
		t.Fatal("duplicate nonnull VIN admitted")
	}
	input.VIN = "VIN-OTHER"
	input.BatterySerialNumber = "BATTERY-REAL"
	if _, err = svc.RegisterProductionUnit(ctx, tenant, input); err == nil {
		t.Fatal("duplicate nonnull battery admitted")
	}
	input.SerialNumber = "FRAME-1"
	input.VIN = ""
	input.BatterySerialNumber = ""
	if _, err = svc.RegisterProductionUnit(ctx, tenant, input); err == nil {
		t.Fatal("duplicate mandatory serial admitted")
	}
	var units, stock, events int
	err = pool.QueryRow(ctx, `select (select count(*) from factory.production_unit where tenant_id=$1),(select count(*) from inventory.stock_unit where tenant_id=$1),(select count(*) from platform.outbox_event where tenant_id=$1)`, tenant).Scan(&units, &stock, &events)
	if err != nil || units != 3 || stock != 3 || events != 7 {
		t.Fatalf("unit/stock/events=%d/%d/%d err=%v", units, stock, events, err)
	}
	t.Log("SERIAL_OPTIONAL_IDENTIFIERS_PASS missing optional identifiers allowed; present identifiers and mandatory serial unique; rejected writes leak no event")
}
