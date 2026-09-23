package postgres

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestEnterpriseQueriesEnforceOrganizationCustomerAndCursor(t *testing.T) {
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
	tenant := "018f4d4a-7b36-7a21-8d10-2f4c54c28901"
	cleanup := func() {
		for _, query := range []string{
			`delete from service_ops.service_case where tenant_id=$1`, `delete from service_ops.warranty where tenant_id=$1`,
			`delete from inventory.stock_unit where tenant_id=$1`, `delete from factory.production_unit where tenant_id=$1`,
			`delete from procurement.purchase_order where tenant_id=$1`, `delete from partner.supplier where tenant_id=$1`,
			`delete from sales.customer_order where tenant_id=$1`, `delete from catalog.vehicle_variant where tenant_id=$1`,
			`delete from catalog.vehicle_model where tenant_id=$1`, `delete from crm.customer_profile where tenant_id=$1`,
			`delete from org.organization where tenant_id=$1`, `delete from platform.tenant where tenant_id=$1`,
		} {
			_, _ = pool.Exec(ctx, query, tenant)
		}
	}
	cleanup()
	defer cleanup()
	fixtures := []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,'query-api','Query','Query')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,'org-a','org-a','A','store'),($1,'org-b','org-b','B','store')`,
		`insert into crm.customer_profile(tenant_id,customer_principal_id,display_name,email_normalized,status) values($1,'customer-1','Customer 1','customer1@example.test','active'),($1,'customer-2','Customer 2','customer2@example.test','active')`,
		`insert into catalog.vehicle_model(tenant_id,model_id,model_code,display_name,vehicle_class,lifecycle_state) values($1,'model','model','Model','bicycle','active')`,
		`insert into catalog.vehicle_variant(tenant_id,variant_id,model_id,variant_code,display_name,battery_specification,lifecycle_state) values($1,'variant','model','variant','Variant','{}','active')`,
		`insert into partner.supplier(tenant_id,supplier_id,supplier_code,legal_name,status) values($1,'supplier','supplier','Supplier','active')`,
		`insert into procurement.purchase_order(tenant_id,purchase_order_id,supplier_id,destination_organization_id,state,currency,total_minor_units,version) values($1,'po','supplier','org-a','accepted','USD',100,1)`,
		`insert into factory.production_unit(tenant_id,production_unit_id,purchase_order_id,variant_id,serial_number,state) values($1,'unit','po','variant','UNIT-QUERY','planned')`,
		`insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version) values($1,'order-a','org-a','customer-1','draft','USD',100,1),($1,'order-b','org-a','customer-2','draft','USD',200,1),($1,'order-c','org-b','customer-1','draft','USD',300,1)`,
		`insert into inventory.stock_unit(tenant_id,stock_unit_id,organization_id,variant_id,serial_number,vin,battery_serial_number,state,version,received_at) values($1,'stock-a','org-a','variant','STOCK-QUERY-A','VIN-QUERY-A','BATTERY-QUERY-A','available',1,clock_timestamp()),($1,'stock-b','org-b','variant','STOCK-QUERY-B','VIN-QUERY-B','BATTERY-QUERY-B','available',1,clock_timestamp())`,
		`insert into service_ops.warranty(tenant_id,warranty_id,stock_unit_id,customer_principal_id,starts_at,ends_at,terms_version,status) values($1,'warranty','stock-a','customer-1',clock_timestamp(),clock_timestamp()+interval '1 year','v1','active')`,
		`insert into service_ops.service_case(tenant_id,service_case_id,stock_unit_id,organization_id,state,severity,description,version) values($1,'case-a','stock-a','org-a','opened','medium','inspection',1)`,
	}
	for _, query := range fixtures {
		if _, err = pool.Exec(ctx, query, tenant); err != nil {
			t.Fatal(err)
		}
	}
	repository := NewEnterpriseQuery(pool)
	overview, err := repository.Overview(ctx, tenant, "org-a")
	if err != nil || overview.Orders != 2 || overview.StockAvailable != 1 || overview.OpenCases != 1 {
		t.Fatalf("overview=%+v err=%v", overview, err)
	}
	first, err := repository.Orders(ctx, tenant, "org-a", "", 1, "")
	if err != nil || len(first.Items) != 1 || first.NextCursor == "" {
		t.Fatalf("first=%+v err=%v", first, err)
	}
	second, err := repository.Orders(ctx, tenant, "org-a", "", 1, first.NextCursor)
	if err != nil || len(second.Items) != 1 || second.Items[0].ID == first.Items[0].ID {
		t.Fatalf("second=%+v err=%v", second, err)
	}
	customerOrders, err := repository.Orders(ctx, tenant, "org-a", "customer-1", 10, "")
	if err != nil || len(customerOrders.Items) != 1 || customerOrders.Items[0].CustomerSubject != "customer-1" {
		t.Fatalf("customer orders=%+v err=%v", customerOrders, err)
	}
	units, err := repository.FactoryUnits(ctx, tenant, "org-a", 10, "")
	if err != nil || len(units.Items) != 1 || units.Items[0].OrganizationID != "org-a" {
		t.Fatalf("units=%+v err=%v", units, err)
	}
	cases, err := repository.ServiceCases(ctx, tenant, "org-a", "customer-1", 10, "")
	if err != nil || len(cases.Items) != 1 {
		t.Fatalf("cases=%+v err=%v", cases, err)
	}
	foreignCases, err := repository.ServiceCases(ctx, tenant, "org-a", "customer-2", 10, "")
	if err != nil || len(foreignCases.Items) != 0 {
		t.Fatalf("foreign cases=%+v err=%v", foreignCases, err)
	}
}
