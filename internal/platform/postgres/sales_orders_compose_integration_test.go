package postgres

// AUTHORED fixture composes ELECTROMOBILITY-FRANCHISE-MODULES schema,
// GO-BC-SALES-CONTRACT-ADAPTER server-side pricing, and
// GO-FRANCHISE-CUSTOMER-JOURNEY-API quote issuance/acceptance into
// sales.customer_order. No payment capture or live provider effects.
import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"elite.local/enterprise/internal/franchisejourney"
	"elite.local/enterprise/internal/platform/randomid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func salesOrdersComposePool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	raw := os.Getenv("TEST_DATABASE_URL")
	if raw == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	cfg, err := pgxpool.ParseConfig(raw)
	if err != nil || cfg.ConnConfig.Host != "127.0.0.1" || !strings.HasPrefix(cfg.ConnConfig.Database, "elite_confirmation_") {
		t.Fatal("requires disposable loopback elite_confirmation_* database")
	}
	pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)
	return pool
}

func seedSalesOrdersComposeTenant(t *testing.T, pool *pgxpool.Pool) (string, *FranchiseJourney) {
	t.Helper()
	ctx := context.Background()
	tenant := randomid.Generator{}.New()
	for _, q := range []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1::uuid,'sales-compose-'||replace(($1::uuid)::text,'-',''),'Synthetic','Synthetic')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'store','store','Synthetic','store')`,
		`insert into catalog.vehicle_model(tenant_id,model_id,model_code,display_name,vehicle_class,lifecycle_state)values($1,'model','model','Synthetic','bicycle','active')`,
		`insert into catalog.vehicle_variant(tenant_id,variant_id,model_id,variant_code,display_name,battery_specification,lifecycle_state)values($1,'variant','model','variant','Synthetic','{}','active')`,
		`insert into crm.customer_profile(tenant_id,customer_principal_id,display_name,email_normalized)values($1,'customer','Synthetic','fixture@example.test')`,
		`insert into crm.lead(tenant_id,lead_id,organization_id,customer_principal_id,model_id,lifecycle_state,source_code,contact_payload)values($1,'lead','store','customer','model','new','fixture','{}')`,
		`insert into pricing.price_book(tenant_id,price_book_id,market,currency,valid_from,status)values($1,'retail','AR','ARS',clock_timestamp()-interval '1 day','active')`,
		`insert into pricing.price_book_entry(tenant_id,price_book_id,variant_id,amount_minor_units,tax_mode)values($1,'retail','variant',123456,'inclusive')`,
	} {
		if _, err := pool.Exec(ctx, q, tenant); err != nil {
			t.Fatal(err)
		}
	}
	return tenant, NewFranchiseJourney(pool)
}

func TestSalesOrdersComposeQuoteToOrder(t *testing.T) {
	pool := salesOrdersComposePool(t)
	ctx := context.Background()
	tenant, repo := seedSalesOrdersComposeTenant(t, pool)
	ids := randomid.Generator{}
	quoteHash := strings.Repeat("5", 64)
	evidence := strings.Repeat("d", 64)

	quote, replayed, err := repo.CreateQuoteAs(ctx, tenant, "sales-compose-quote-0001", franchisejourney.Quote{
		ID: "quote", OrganizationID: "store", LeadID: "lead", VariantID: "variant", PriceBookID: "retail",
		ValidUntil: time.Now().UTC().Add(24 * time.Hour), State: "issued", Version: 1,
	}, quoteHash, ids.New(), "writer")
	if err != nil || replayed || quote.Currency != "ARS" || quote.TotalMinorUnits != 123456 {
		t.Fatalf("quote=%+v replayed=%v err=%v", quote, replayed, err)
	}

	accepted, err := repo.AcceptQuote(ctx, tenant, "store", "customer", "quote", 1, evidence, "order", "line", ids.New(), ids.New())
	if err != nil || accepted.State != "accepted" || accepted.OrderID != "order" || accepted.Version != 2 {
		t.Fatalf("accepted=%+v err=%v", accepted, err)
	}

	var orderState, currency string
	var total int64
	if err = pool.QueryRow(ctx, `select state,currency,total_minor_units from sales.customer_order where tenant_id=$1 and order_id='order'`, tenant).Scan(&orderState, &currency, &total); err != nil || orderState != "placed" || currency != "ARS" || total != 123456 {
		t.Fatalf("order state=%s currency=%s total=%d err=%v", orderState, currency, total, err)
	}
	var linePrice int64
	if err = pool.QueryRow(ctx, `select unit_price_minor_units from sales.customer_order_line where tenant_id=$1 and order_id='order'`, tenant).Scan(&linePrice); err != nil || linePrice != 123456 {
		t.Fatalf("line price=%d err=%v", linePrice, err)
	}
	var acceptances, issuedEvents, acceptedEvents, placedEvents int
	if err = pool.QueryRow(ctx, `select count(*) from sales.quotation_acceptance where tenant_id=$1 and quotation_id='quote' and order_id='order'`, tenant).Scan(&acceptances); err != nil || acceptances != 1 {
		t.Fatalf("acceptances=%d err=%v", acceptances, err)
	}
	if err = pool.QueryRow(ctx, `select count(*) from platform.outbox_event where tenant_id=$1 and aggregate_id='quote' and event_type='quotation.issued'`, tenant).Scan(&issuedEvents); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select count(*) from platform.outbox_event where tenant_id=$1 and aggregate_id='quote' and event_type='quotation.accepted'`, tenant).Scan(&acceptedEvents); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select count(*) from platform.outbox_event where tenant_id=$1 and aggregate_id='order' and event_type='customer-order.placed'`, tenant).Scan(&placedEvents); err != nil || issuedEvents != 1 || acceptedEvents != 1 || placedEvents != 1 {
		t.Fatalf("outbox issued=%d accepted=%d placed=%d err=%v", issuedEvents, acceptedEvents, placedEvents, err)
	}
	if _, err = repo.AcceptQuote(ctx, tenant, "store", "customer", "quote", 1, evidence, "stale-order", "stale-line", ids.New(), ids.New()); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("stale quote acceptance succeeded")
	}

	t.Logf("SALES_ORDERS_COMPOSE_PASS tenant=%s quote=%s order=%s placed=true acceptance=true outbox=true server_price=true payment=false", tenant, quote.ID, accepted.OrderID)
}
