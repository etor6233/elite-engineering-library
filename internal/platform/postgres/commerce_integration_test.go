package postgres

import (
	"context"
	"elite.local/enterprise/internal/commerce"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"testing"
	"time"
)

func TestCommercePriceOrderAllocationPaymentFlow(t *testing.T) {
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
	tenant := "018f4d4a-7b36-7a21-8d10-2f4c54c28401"
	cleanup := func() {
		for _, q := range []string{`delete from platform.outbox_event where tenant_id=$1`, `delete from payment.payment_attempt where tenant_id=$1`, `delete from inventory.serial_reservation where tenant_id=$1`, `delete from sales.customer_order_line where tenant_id=$1`, `delete from sales.customer_order where tenant_id=$1`, `delete from inventory.stock_unit where tenant_id=$1`, `delete from pricing.price_book_entry where tenant_id=$1`, `delete from pricing.price_book where tenant_id=$1`, `delete from catalog.vehicle_variant where tenant_id=$1`, `delete from catalog.vehicle_model where tenant_id=$1`, `delete from org.organization where tenant_id=$1`, `delete from platform.tenant where tenant_id=$1`} {
			_, _ = pool.Exec(ctx, q, tenant)
		}
	}
	cleanup()
	defer cleanup()
	fixtures := []string{`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1,'commerce-api','Commerce','Commerce')`, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'org','store','Store','store')`, `insert into catalog.vehicle_model(tenant_id,model_id,model_code,display_name,vehicle_class,lifecycle_state)values($1,'model','model','Model','bicycle','active')`, `insert into catalog.vehicle_variant(tenant_id,variant_id,model_id,variant_code,display_name,battery_specification,lifecycle_state)values($1,'variant','model','variant','Variant','{}','active')`, `insert into inventory.stock_unit(tenant_id,stock_unit_id,organization_id,variant_id,serial_number,state,version,received_at)values($1,'stock','org','variant','STOCK-SERIAL','available',1,clock_timestamp())`, `insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version)values($1,'order','org','customer','draft','ARS',1000,1)`}
	for _, q := range fixtures {
		if _, err := pool.Exec(ctx, q, tenant); err != nil {
			t.Fatal(err)
		}
	}
	repo := NewCommerce(pool)
	book := commerce.PriceBook{ID: "book", Market: "AR", Currency: "ARS", ValidFrom: time.Now().Add(-time.Hour), Status: "draft", Entries: []commerce.PriceEntry{{VariantID: "variant", AmountMinorUnits: 1000, TaxMode: "inclusive"}}}
	if err := repo.CreatePriceBook(ctx, tenant, "018f4d4a-7b36-7a21-8d10-2f4c54c28402", book); err != nil {
		t.Fatal(err)
	}
	if err := repo.ActivatePriceBook(ctx, tenant, "book", "018f4d4a-7b36-7a21-8d10-2f4c54c28403"); err != nil {
		t.Fatal(err)
	}
	price, err := repo.PublicPrice(ctx, "commerce-api", "AR", "variant")
	if err != nil || price.AmountMinorUnits != 1000 {
		t.Fatalf("price=%+v err=%v", price, err)
	}
	line := commerce.OrderLine{OrderID: "order", LineID: "line", OrganizationID: "org", PriceBookID: "book", VariantID: "variant", Quantity: 1}
	line, err = repo.AddOrderLine(ctx, tenant, "018f4d4a-7b36-7a21-8d10-2f4c54c28404", line, 1)
	if err != nil || line.UnitPriceMinorUnits != 1000 || line.OrderVersion != 2 {
		t.Fatalf("line=%+v err=%v", line, err)
	}
	if err := repo.PlaceOrder(ctx, tenant, "org", "order", 2, "018f4d4a-7b36-7a21-8d10-2f4c54c28405"); err != nil {
		t.Fatal(err)
	}
	if err := repo.AllocateStock(ctx, tenant, "org-other", "order", "line", "stock", 3, 1, "018f4d4a-7b36-7a21-8d10-2f4c54c28411"); err == nil {
		t.Fatal("allocation crossed organization scope")
	}
	if err := repo.AllocateStock(ctx, tenant, "org", "order", "line", "stock", 3, 1, "018f4d4a-7b36-7a21-8d10-2f4c54c28406"); err != nil {
		t.Fatal(err)
	}
	if err := repo.AllocateStock(ctx, tenant, "org", "order", "line", "stock", 3, 1, "018f4d4a-7b36-7a21-8d10-2f4c54c28407"); err == nil {
		t.Fatal("duplicate allocation succeeded")
	}
	payment := commerce.PaymentAttempt{ID: "payment", OrderID: "order", OrganizationID: "org", ProviderCode: "sandbox", Currency: "ARS", AmountMinorUnits: 1000}
	if err := repo.CreatePaymentAttempt(ctx, tenant, "018f4d4a-7b36-7a21-8d10-2f4c54c28408", "payment-key-00000001", payment); err != nil {
		t.Fatal(err)
	}
	if err := repo.TransitionPayment(ctx, tenant, "org-other", "payment", "created", "authorized", 1, "provider-ref", "018f4d4a-7b36-7a21-8d10-2f4c54c28412"); err == nil {
		t.Fatal("payment crossed organization scope")
	}
	if err := repo.TransitionPayment(ctx, tenant, "org", "payment", "created", "authorized", 1, "provider-ref", "018f4d4a-7b36-7a21-8d10-2f4c54c28409"); err != nil {
		t.Fatal(err)
	}
	if err := repo.TransitionPayment(ctx, tenant, "org", "payment", "created", "failed", 1, "", "018f4d4a-7b36-7a21-8d10-2f4c54c28410"); err == nil {
		t.Fatal("stale payment transition succeeded")
	}
	var state string
	var version int64
	if err := pool.QueryRow(ctx, `select state,version from inventory.stock_unit where tenant_id=$1 and stock_unit_id='stock'`, tenant).Scan(&state, &version); err != nil || state != "reserved" || version != 2 {
		t.Fatalf("stock=%s version=%d err=%v", state, version, err)
	}
	var reservationStatus string
	if err := pool.QueryRow(ctx, `select status from inventory.serial_reservation where tenant_id=$1 and stock_unit_id='stock' and demand_kind='customer-order' and demand_id='order' and demand_line_id='line'`, tenant).Scan(&reservationStatus); err != nil || reservationStatus != "reservation" {
		t.Fatalf("reservation=%s err=%v", reservationStatus, err)
	}
	view, err := repo.Operations(ctx, tenant, "org")
	if err != nil || view.Truncated || len(view.Orders) != 1 || len(view.Stock) != 0 {
		t.Fatalf("operation snapshot mismatch err=%v", err)
	}
	if len(view.Orders[0].Lines) != 1 || view.Orders[0].Lines[0].StockID != "stock" || len(view.Orders[0].Payments) != 1 || view.Orders[0].Payments[0].State != "authorized" || len(view.Orders[0].Handovers) != 0 {
		t.Fatal("operation projection differs from durable owners")
	}
	for _, scope := range []struct{ tenant, org string }{{tenant, "other"}, {"018f4d4a-7b36-7a21-8d10-2f4c54c28999", "org"}} {
		other, err := repo.Operations(ctx, scope.tenant, scope.org)
		if err != nil || len(other.Orders) != 0 || len(other.Stock) != 0 {
			t.Fatal("operation query crosses scope")
		}
	}
	// Every exposed truncation is explicit; incomplete data cannot authorize a UI write.
	if _, err := pool.Exec(ctx, `insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version) select $1,'limit-'||s::text,'org','customer','draft','ARS',1000,1 from generate_series(1,26)s`, tenant); err != nil {
		t.Fatal(err)
	}
	view, err = repo.Operations(ctx, tenant, "org")
	if err != nil || !view.Truncated || len(view.Orders) != 25 {
		t.Fatal("operation order limit not explicit")
	}
}
