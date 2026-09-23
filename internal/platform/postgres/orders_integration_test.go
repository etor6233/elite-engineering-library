package postgres

import (
	"context"
	"errors"
	"os"
	"testing"

	"elite.local/enterprise/internal/order"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestOrdersRepositoryIntegration(t *testing.T) {
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	const tenantID = "018f4d4a-7b36-7a21-8d10-2f4c54c28b01"
	const organizationID = "integration-org"
	const orderID = "integration-order"
	cleanup := func() {
		_, _ = pool.Exec(ctx, `delete from platform.outbox_event where tenant_id=$1`, tenantID)
		_, _ = pool.Exec(ctx, `delete from platform.idempotency_record where tenant_id=$1`, tenantID)
		_, _ = pool.Exec(ctx, `delete from sales.customer_order where tenant_id=$1`, tenantID)
		_, _ = pool.Exec(ctx, `delete from org.organization where tenant_id=$1`, tenantID)
		_, _ = pool.Exec(ctx, `delete from platform.tenant where tenant_id=$1`, tenantID)
	}
	cleanup()
	defer cleanup()
	if _, err := pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,'integration-tenant','Integration Tenant S.A.','Integration Tenant')`, tenantID); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,$2,'integration-org','Integration Org','enterprise')`, tenantID, organizationID); err != nil {
		t.Fatal(err)
	}

	repository := NewOrders(pool)
	created := order.Order{ID: orderID, TenantID: tenantID, OrganizationID: organizationID, CustomerPrincipal: "customer-1", State: order.Draft, Currency: "USD", TotalMinorUnits: 100, Version: 1}
	actual, replayed, err := repository.Create(ctx, created, "integration-idem-0001", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	if err != nil || replayed || actual.ID != orderID {
		t.Fatalf("create failed: order=%+v replayed=%v err=%v", actual, replayed, err)
	}
	replay, replayed, err := repository.Create(ctx, created, "integration-idem-0001", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	if err != nil || !replayed || replay.ID != orderID {
		t.Fatalf("replay failed: order=%+v replayed=%v err=%v", replay, replayed, err)
	}
	if _, _, err := repository.Create(ctx, created, "integration-idem-0001", "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"); !errors.Is(err, order.ErrConflict) {
		t.Fatalf("expected hash conflict, got %v", err)
	}
	loaded, err := repository.Get(ctx, tenantID, organizationID, orderID)
	if err != nil || loaded.TenantID != tenantID {
		t.Fatalf("tenant-scoped get failed: order=%+v err=%v", loaded, err)
	}
	next, err := loaded.Transition(order.Placed, map[string]struct{}{"order:create": {}})
	if err != nil {
		t.Fatal(err)
	}
	if err := repository.SaveTransition(ctx, next, loaded.Version); err != nil {
		t.Fatal(err)
	}
	loaded, err = repository.Get(ctx, tenantID, organizationID, orderID)
	if err != nil || loaded.State != order.Placed || loaded.Version != 2 {
		t.Fatalf("transition persistence failed: order=%+v err=%v", loaded, err)
	}
}
