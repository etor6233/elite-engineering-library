// AUTHORED focused proof for the new order/plan transaction composition.
package postgres_test

import (
	"context"
	"elite.local/enterprise/internal/platform/identity"
	db "elite.local/enterprise/internal/platform/postgres"
	sc "elite.local/enterprise/internal/serialsupply"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"strings"
	"sync"
	"testing"
)

func TestSerialSupplyCreateAtomic(t *testing.T) {
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
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1,'supply-create','Fixture','Fixture')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'store','store','Fixture','store'),($1,'factory','factory','Fixture','factory')`,
		`insert into partner.supplier(tenant_id,supplier_id,supplier_code,legal_name,status)values($1,'supplier','supplier','Fixture','active')`,
		`insert into catalog.vehicle_model(tenant_id,model_id,model_code,display_name,vehicle_class,lifecycle_state)values($1,'model','model','Fixture','bicycle','active')`,
		`insert into catalog.vehicle_variant(tenant_id,variant_id,model_id,variant_code,display_name,battery_specification,lifecycle_state)values($1,'variant','model','variant','Fixture','{}','active')`,
	} {
		if _, e = pool.Exec(ctx, q, tenant); e != nil {
			t.Fatal(e)
		}
	}
	store, e := db.NewSerialSupply(pool)
	if e != nil {
		t.Fatal(e)
	}
	p := identity.Principal{TenantID: tenant, Subject: "buyer", Organizations: map[string]struct{}{"store": {}}, Permissions: map[string]struct{}{"supply:plan": {}, "supply:read": {}}}
	r := sc.CreateRequest{PlanRequest: sc.PlanRequest{PurchaseOrderID: uuid.NewString(), CommandID: uuid.NewString(), FactoryOrganizationID: "factory", DemandReference: "Demand <á> & fixture", PolicyCode: sc.PolicyCode, EvidenceSHA256: strings.Repeat("a", 64), Lines: []sc.Line{{ID: "line-1", VariantID: "variant", Quantity: 2}}}, DestinationOrganizationID: "store", SupplierID: "supplier", Currency: "ARS", TotalMinorUnits: 9007199254740993}
	snapshot := func() string {
		t.Helper()
		out := ""
		for _, table := range []string{"procurement.purchase_order", "procurement.serial_supply_plan", "procurement.serial_supply_line", "procurement.serial_supply_step", "platform.outbox_event"} {
			var text string
			if e = pool.QueryRow(ctx, "select coalesce(jsonb_agg(to_jsonb(x) order by to_jsonb(x)::text),'[]'::jsonb)::text from "+table+" x where tenant_id=$1", tenant).Scan(&text); e != nil {
				t.Fatal(e)
			}
			out += text
		}
		return out
	}
	before := snapshot()
	bad := r
	bad.Lines = []sc.Line{{ID: "line-1", VariantID: "missing", Quantity: 2}}
	if _, e = store.Create(ctx, p, bad); e == nil {
		t.Fatal("invalid variant accepted")
	}
	if snapshot() != before {
		t.Fatal("order/outbox escaped failed plan")
	}
	// Failure after the plan receipt also rolls back its purchase owner/outbox.
	if _, e = pool.Exec(ctx, `create function public.supply_role_fail()returns trigger language plpgsql as $$begin if new.aggregate_type='serial-supply' and new.payload->>'command_id'='forced-rollback' then raise exception 'fixture';end if;return new;end$$;create trigger supply_role_fail before insert on platform.outbox_event for each row execute function public.supply_role_fail()`); e != nil {
		t.Fatal(e)
	}
	bad = r
	bad.CommandID = "forced-rollback"
	if _, e = store.Create(ctx, p, bad); e == nil {
		t.Fatal("forced outbox accepted")
	}
	if snapshot() != before {
		t.Fatal("receipt rollback leaked")
	}
	if _, e = pool.Exec(ctx, `drop trigger supply_role_fail on platform.outbox_event;drop function public.supply_role_fail()`); e != nil {
		t.Fatal(e)
	}
	var wg sync.WaitGroup
	out := make(chan sc.Receipt, 4)
	errs := make(chan error, 4)
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); v, e := store.Create(ctx, p, r); out <- v; errs <- e }()
	}
	wg.Wait()
	close(out)
	close(errs)
	for e := range errs {
		if e != nil {
			t.Fatal(e)
		}
	}
	fresh, replay := 0, 0
	for v := range out {
		if v.Replay {
			replay++
		} else {
			fresh++
		}
		if v.PurchaseOrderID != r.PurchaseOrderID || v.Version != 1 || v.Actor != p.Subject || v.Kind != "planned" {
			t.Fatal(v)
		}
	}
	if fresh != 1 || replay != 3 {
		t.Fatal(fresh, replay)
	}
	durable := snapshot()
	for _, variation := range []string{"actor", "hash", "org"} {
		other := p
		bad = r
		switch variation {
		case "actor":
			other.Subject = "other"
		case "hash":
			bad.TotalMinorUnits++
		case "org":
			bad.DestinationOrganizationID = "other"
		}
		if _, e = store.Create(ctx, other, bad); e == nil {
			t.Fatal("replay mismatch", variation)
		}
	}
	if snapshot() != durable {
		t.Fatal("negative replay changed history")
	}
	plan, e := store.Plan(ctx, p, r.PurchaseOrderID, "")
	if e != nil || plan.TotalMinorUnits != r.TotalMinorUnits || plan.Currency != "ARS" || len(plan.Lines) != 1 || plan.Version != 1 {
		t.Fatal(plan, e)
	}
	rebuilt, e := db.NewSerialSupply(pool)
	if e != nil {
		t.Fatal(e)
	}
	receipt, e := rebuilt.CommandReceipt(ctx, p, r.PurchaseOrderID, r.CommandID)
	_, hash, _ := sc.Canonical(r)
	if e != nil || receipt.RequestSHA256 != hash {
		t.Fatal("durable recovery", receipt, e)
	}
	var counts string
	if e = pool.QueryRow(ctx, `select jsonb_build_array((select count(*)from procurement.purchase_order where tenant_id=$1),(select count(*)from procurement.serial_supply_plan where tenant_id=$1),(select count(*)from procurement.serial_supply_step where tenant_id=$1),(select count(*)from platform.outbox_event where tenant_id=$1))::text`, tenant).Scan(&counts); e != nil || counts != "[1, 1, 1, 2]" {
		t.Fatal(counts, e)
	}
	rawPlan, _ := json.Marshal(plan)
	if !strings.Contains(string(rawPlan), `"total_minor_units":"9007199254740993"`) {
		t.Fatal("unsafe numeric projection", string(rawPlan))
	}
	t.Log("SUPPLY_CREATE_ATOMIC_PASS concurrency=1new3replay tables5_rollback=true actor_hash_scope_bound=true durable_recovery=true exact_int64=true")
}
