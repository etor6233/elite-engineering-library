package postgres_test

// AUTHORED focused regression for the newly materializable source write only.
// Publication, four reviews and storefront will be exercised by the role browser.
import (
	"context"
	cr "elite.local/enterprise/internal/catalogrelease"
	"elite.local/enterprise/internal/electromobility"
	"elite.local/enterprise/internal/platform/identity"
	db "elite.local/enterprise/internal/platform/postgres"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"sync"
	"testing"
	"time"
)

func TestCatalogSourceAuthoringAtomic(t *testing.T) {
	raw := os.Getenv("PAYMENT_CONNECTED_DB_URL")
	if raw == "" {
		t.Skip("owned fixture required")
	}
	ctx := context.Background()
	pool, e := pgxpool.New(ctx, raw)
	if e != nil {
		t.Fatal(e)
	}
	defer pool.Close()
	tenant := uuid.NewString()
	for _, q := range []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1,'source-'||($1::uuid)::text,'Fixture','Fixture')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'source-org','source-org','Fixture','store')`,
	} {
		if _, e = pool.Exec(ctx, q, tenant); e != nil {
			t.Fatal(e)
		}
	}
	profile := cr.Profile{Schema: "elite-catalog-publication/v1", TenantID: tenant, OrganizationID: "source-org", Market: "AR", Currency: "ARS", Origin: "https://catalog.example.invalid", PolicyCode: cr.PolicyCode}
	store, e := db.NewCatalogRelease(pool, profile, db.NewCommerce(pool))
	if e != nil {
		t.Fatal(e)
	}
	actor := identity.Principal{TenantID: tenant, Subject: "source-maker", Organizations: map[string]struct{}{"source-org": {}}, Permissions: map[string]struct{}{"catalog:draft": {}, "catalog:read": {}}}
	in := cr.SourceRequest{CommandID: "source-create", Model: electromobility.Model{Code: "created-through-source", DisplayName: "Draft authored model", VehicleClass: "bicycle", Specification: json.RawMessage(`{"description":"Fixture"}`)}, Variants: []cr.SourceVariant{{Code: "source-variant", DisplayName: "Reference variant", BatterySpecification: json.RawMessage(`{}`), AmountMinorUnits: 123456, TaxMode: "not-applicable"}}, ValidFrom: time.Now().UTC().Truncate(time.Second).Add(-time.Hour)}
	const n = 4
	var wg sync.WaitGroup
	got := make(chan cr.Receipt, n)
	errs := make(chan error, n)
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); r, e := store.CreateSource(ctx, actor, in); got <- r; errs <- e }()
	}
	wg.Wait()
	close(got)
	close(errs)
	for e := range errs {
		if e != nil {
			t.Fatal(e)
		}
	}
	var result cr.Receipt
	newCount := 0
	for r := range got {
		if !r.Replay {
			newCount++
		}
		if result.ResourceID != "" && result.ResourceID != r.ResourceID {
			t.Fatal("different model")
		}
		result = r
	}
	if newCount != 1 || len(result.SourceVariantIDs) != 1 || result.SourcePriceBookID == "" {
		t.Fatal("duplicate or missing source receipt", newCount, result)
	}
	snapshot := func() string {
		var raw string
		e := pool.QueryRow(ctx, `select jsonb_build_array(
  (select count(*) from catalog.vehicle_model where tenant_id=$1),
  (select count(*) from catalog.vehicle_variant where tenant_id=$1),
  (select count(*) from pricing.price_book where tenant_id=$1),
  (select count(*) from pricing.price_book_entry where tenant_id=$1),
  (select count(*) from catalog.release_command where tenant_id=$1),
  (select count(*) from platform.outbox_event where tenant_id=$1))::text`, tenant).Scan(&raw)
		if e != nil {
			t.Fatal(e)
		}
		return raw
	}
	before := snapshot()
	if before != "[1, 1, 1, 1, 1, 3]" {
		t.Fatal(before)
	}
	recovered, e := store.CommandReceipt(ctx, actor, in.CommandID)
	if e != nil || recovered.ResourceID != result.ResourceID || recovered.RequestSHA256 != result.RequestSHA256 {
		t.Fatal("read recovery", e)
	}
	other := actor
	other.Subject = "other"
	if _, e = store.CreateSource(ctx, other, in); e == nil {
		t.Fatal("actor replay")
	}
	if _, e = store.CommandReceipt(ctx, other, in.CommandID); e == nil {
		t.Fatal("receipt leak")
	}
	edited := in
	edited.Model.DisplayName = "Changed"
	if _, e = store.CreateSource(ctx, actor, edited); e == nil {
		t.Fatal("hash conflict")
	}
	foreign := actor
	foreign.Organizations = map[string]struct{}{"foreign": {}}
	if _, e = store.CreateSource(ctx, foreign, in); e == nil {
		t.Fatal("org bypass")
	}
	if snapshot() != before {
		t.Fatal("rejections wrote")
	}
	var hidden bool
	if e = pool.QueryRow(ctx, `select m.lifecycle_state='draft' and not m.publicly_visible and v.lifecycle_state='draft' and v.homologation_state='unknown' and b.status='draft'
 from catalog.vehicle_model m join catalog.vehicle_variant v using(tenant_id,model_id) join pricing.price_book_entry x using(tenant_id,variant_id) join pricing.price_book b using(tenant_id,price_book_id) where m.tenant_id=$1`, tenant).Scan(&hidden); e != nil || !hidden {
		t.Fatal("source promoted without review", e)
	}
	_, e = pool.Exec(ctx, `create function catalog.test_source_outbox_fail()returns trigger language plpgsql as $$begin if new.aggregate_type='catalog-release-command' and new.event_type='catalog-release.source' then raise exception 'fixture outbox failure';end if;return new;end$$;
 create trigger test_source_outbox_fail before insert on platform.outbox_event for each row execute function catalog.test_source_outbox_fail()`)
	if e != nil {
		t.Fatal(e)
	}
	broken := in
	broken.CommandID = "source-rollback"
	broken.Model.Code = "rolled-back-model"
	broken.Variants = append([]cr.SourceVariant(nil), in.Variants...)
	broken.Variants[0].Code = "rolled-back-variant"
	if _, e = store.CreateSource(ctx, actor, broken); e == nil {
		t.Fatal("outbox accepted")
	}
	if snapshot() != before {
		t.Fatal("source/model/price/outbox partial commit")
	}
	if _, e = pool.Exec(ctx, `drop trigger test_source_outbox_fail on platform.outbox_event;drop function catalog.test_source_outbox_fail()`); e != nil {
		t.Fatal(e)
	}
	reopened, e := db.NewCatalogRelease(pool, profile, db.NewCommerce(pool))
	if e != nil {
		t.Fatal(e)
	}
	same, e := reopened.CommandReceipt(ctx, actor, in.CommandID)
	if e != nil || same.SourcePriceBookID != result.SourcePriceBookID {
		t.Fatal("restart receipt", e)
	}
	t.Log("CATALOG_SOURCE_AUTHORING_PASS model_variant_price_original_owners=true concurrent_one_new_three_replay=true command_actor_hash_org_bound=true all6tables_atomic=true source_draft_no_publication=true durable_recovery=true")
}
