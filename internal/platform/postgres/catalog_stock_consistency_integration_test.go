package postgres

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

const catalogStockSliceTenant = "018f4d4a-7b36-7a21-8d10-2f4c54c2888"

func requireCatalogStockDB(t *testing.T) *pgxpool.Pool {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	return pool
}

func cleanupCatalogStockSlice(t *testing.T, pool *pgxpool.Pool) {
	ctx := context.Background()
	tenant := catalogStockSliceTenant
	for _, q := range []string{
		`delete from pricing.price_book_entry where tenant_id=$1`,
		`delete from pricing.price_book where tenant_id=$1`,
		`delete from inventory.stock_unit where tenant_id=$1`,
		`delete from catalog.vehicle_variant where tenant_id=$1`,
		`delete from catalog.vehicle_model where tenant_id=$1`,
		`delete from org.organization where tenant_id=$1`,
		`delete from platform.tenant where tenant_id=$1`,
	} {
		if _, err := pool.Exec(ctx, q, tenant); err != nil {
			t.Fatal(err)
		}
	}
}

func seedCatalogStockSlice(t *testing.T, pool *pgxpool.Pool) {
	ctx := context.Background()
	tenant := catalogStockSliceTenant
	fixtures := []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1,'catalog-stock-go','Catalog Stock Go','Catalog Stock Go')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'store','store','Store','store')`,
		`insert into catalog.vehicle_model(tenant_id,model_id,model_code,display_name,vehicle_class,lifecycle_state)values($1,'model','model','Go Model','bicycle','active')`,
		`insert into catalog.vehicle_variant(tenant_id,variant_id,model_id,variant_code,display_name,battery_specification,lifecycle_state)values($1,'variant','model','variant','Go Variant','{}','active')`,
		`insert into inventory.stock_unit(tenant_id,stock_unit_id,organization_id,variant_id,serial_number,state,version,received_at)values($1,'stock','store','variant','GO-SLICE-SERIAL','available',1,clock_timestamp())`,
	}
	for _, q := range fixtures {
		if _, err := pool.Exec(ctx, q, tenant); err != nil {
			t.Fatal(err)
		}
	}
}

func TestCatalogStockConsistencyForeignKeys(t *testing.T) {
	pool := requireCatalogStockDB(t)
	defer pool.Close()
	cleanupCatalogStockSlice(t, pool)
	defer cleanupCatalogStockSlice(t, pool)
	seedCatalogStockSlice(t, pool)

	ctx := context.Background()
	tenant := catalogStockSliceTenant
	_, err := pool.Exec(ctx, `insert into inventory.stock_unit(tenant_id,stock_unit_id,organization_id,variant_id,serial_number,state,version,received_at)values($1,'bad-stock','store','ghost-variant','GHOST-SERIAL','available',1,clock_timestamp())`, tenant)
	if err == nil {
		t.Fatal("expected FK violation inserting stock for unknown variant")
	}
}

func TestCatalogStockConsistencyJoin(t *testing.T) {
	pool := requireCatalogStockDB(t)
	defer pool.Close()
	cleanupCatalogStockSlice(t, pool)
	defer cleanupCatalogStockSlice(t, pool)
	seedCatalogStockSlice(t, pool)

	ctx := context.Background()
	tenant := catalogStockSliceTenant
	var label string
	err := pool.QueryRow(ctx, `select v.display_name from inventory.stock_unit s join catalog.vehicle_variant v using(tenant_id,variant_id) where s.tenant_id=$1 and s.stock_unit_id='stock'`, tenant).Scan(&label)
	if err != nil || label != "Go Variant" {
		t.Fatalf("join label=%q err=%v", label, err)
	}
}

func TestCatalogStockConsistencyPriceBook(t *testing.T) {
	pool := requireCatalogStockDB(t)
	defer pool.Close()
	cleanupCatalogStockSlice(t, pool)
	defer cleanupCatalogStockSlice(t, pool)
	seedCatalogStockSlice(t, pool)

	ctx := context.Background()
	tenant := catalogStockSliceTenant
	if _, err := pool.Exec(ctx, `insert into pricing.price_book(tenant_id,price_book_id,market,currency,valid_from,status)values($1,'retail','AR','ARS',clock_timestamp()-interval '1 day','active')`, tenant); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `insert into pricing.price_book_entry(tenant_id,price_book_id,variant_id,amount_minor_units,tax_mode)values($1,'retail','variant',99900,'inclusive')`, tenant); err != nil {
		t.Fatal(err)
	}
	_, err := pool.Exec(ctx, `insert into pricing.price_book_entry(tenant_id,price_book_id,variant_id,amount_minor_units,tax_mode)values($1,'retail','ghost-variant',100,'inclusive')`, tenant)
	if err == nil {
		t.Fatal("expected FK violation for price book entry on unknown variant")
	}
	var amount int64
	if err := pool.QueryRow(ctx, `select amount_minor_units from pricing.price_book_entry where tenant_id=$1 and price_book_id='retail' and variant_id='variant'`, tenant).Scan(&amount); err != nil || amount != 99900 {
		t.Fatalf("amount=%d err=%v", amount, err)
	}
}
