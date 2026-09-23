package postgres

import (
	"context"
	"elite.local/enterprise/internal/operations"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"testing"
)

func TestOperationsFlowAndConcurrency(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := "018f4d4a-7b36-7a21-8d10-2f4c54c28101"
	cleanup := func() {
		for _, q := range []string{`delete from platform.outbox_event where tenant_id=$1`, `delete from inventory.stock_unit where tenant_id=$1`, `delete from factory.production_unit where tenant_id=$1`, `delete from procurement.purchase_order where tenant_id=$1`, `delete from partner.supplier where tenant_id=$1`, `delete from catalog.vehicle_variant where tenant_id=$1`, `delete from catalog.vehicle_model where tenant_id=$1`, `delete from org.organization where tenant_id=$1`, `delete from platform.tenant where tenant_id=$1`} {
			_, _ = pool.Exec(ctx, q, tenant)
		}
	}
	cleanup()
	defer cleanup()
	fixtures := []struct {
		q    string
		args []any
	}{{`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1,'ops-api','Ops','Ops')`, []any{tenant}}, {`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'org','central','Central','warehouse')`, []any{tenant}}, {`insert into partner.supplier(tenant_id,supplier_id,supplier_code,legal_name,status)values($1,'supplier','supplier','Supplier','active')`, []any{tenant}}, {`insert into catalog.vehicle_model(tenant_id,model_id,model_code,display_name,vehicle_class,lifecycle_state)values($1,'model','model','Model','bicycle','active')`, []any{tenant}}, {`insert into catalog.vehicle_variant(tenant_id,variant_id,model_id,variant_code,display_name,battery_specification,lifecycle_state)values($1,'variant','model','variant','Variant','{}','active')`, []any{tenant}}}
	for _, f := range fixtures {
		if _, err := pool.Exec(ctx, f.q, f.args...); err != nil {
			t.Fatal(err)
		}
	}
	repo := NewOperations(pool)
	po := operations.PurchaseOrder{ID: "po", SupplierID: "supplier", DestinationOrganizationID: "org", State: "draft", Currency: "USD", TotalMinorUnits: 100, Version: 1}
	if err := repo.CreatePurchaseOrder(ctx, tenant, "018f4d4a-7b36-7a21-8d10-2f4c54c28102", po); err != nil {
		t.Fatal(err)
	}
	if err := repo.TransitionPurchaseOrder(ctx, tenant, "org-other", "po", "draft", 1, "submitted", "018f4d4a-7b36-7a21-8d10-2f4c54c28110"); err == nil {
		t.Fatal("purchase order crossed organization scope")
	}
	if err := repo.TransitionPurchaseOrder(ctx, tenant, "org", "po", "draft", 1, "submitted", "018f4d4a-7b36-7a21-8d10-2f4c54c28106"); err != nil {
		t.Fatal(err)
	}
	if err := repo.TransitionPurchaseOrder(ctx, tenant, "org", "po", "draft", 1, "cancelled", "018f4d4a-7b36-7a21-8d10-2f4c54c28107"); err == nil {
		t.Fatal("stale purchase order transition succeeded")
	}
	unit := operations.ProductionUnit{ID: "unit", OrganizationID: "org", PurchaseOrderID: "po", VariantID: "variant", SerialNumber: "SERIAL", State: "planned"}
	if err := repo.CreateProductionUnit(ctx, tenant, "018f4d4a-7b36-7a21-8d10-2f4c54c28103", unit); err != nil {
		t.Fatal(err)
	}
	if err := repo.TransitionProductionUnit(ctx, tenant, "org-other", "unit", "planned", "assembly", "018f4d4a-7b36-7a21-8d10-2f4c54c28111"); err == nil {
		t.Fatal("production unit crossed organization scope")
	}
	if err := repo.TransitionProductionUnit(ctx, tenant, "org", "unit", "planned", "assembly", "018f4d4a-7b36-7a21-8d10-2f4c54c28104"); err != nil {
		t.Fatal(err)
	}
	stock := operations.StockUnit{ID: "stock", OrganizationID: "org", VariantID: "variant", ProductionUnitID: "unit", SerialNumber: "SERIAL", State: "in-transit", Version: 1}
	if err := repo.CreateStockUnit(ctx, tenant, "018f4d4a-7b36-7a21-8d10-2f4c54c28105", stock); err != nil {
		t.Fatal(err)
	}
	if err := repo.TransitionStockUnit(ctx, tenant, "org-other", "stock", "in-transit", 1, "available", "018f4d4a-7b36-7a21-8d10-2f4c54c28112"); err == nil {
		t.Fatal("stock unit crossed organization scope")
	}
	if err := repo.TransitionStockUnit(ctx, tenant, "org", "stock", "in-transit", 1, "available", "018f4d4a-7b36-7a21-8d10-2f4c54c28108"); err != nil {
		t.Fatal(err)
	}
	if err := repo.TransitionStockUnit(ctx, tenant, "org", "stock", "in-transit", 1, "quarantine", "018f4d4a-7b36-7a21-8d10-2f4c54c28109"); err == nil {
		t.Fatal("stale stock transition succeeded")
	}
	var events int
	if err := pool.QueryRow(ctx, `select count(*) from platform.outbox_event where tenant_id=$1`, tenant).Scan(&events); err != nil || events != 6 {
		t.Fatalf("events=%d err=%v", events, err)
	}
}
