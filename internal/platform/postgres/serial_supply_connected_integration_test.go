// AUTHORED connected J2 reference: real source owners, scoped fixture principals.
// No provider accounts, legal quality certification or production acceptance.
package postgres_test

import (
	"context"
	"elite.local/enterprise/internal/operations"
	"elite.local/enterprise/internal/platform/identity"
	db "elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/platform/randomid"
	sc "elite.local/enterprise/internal/serialsupply"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

type serialSupplyReference interface {
	BindPlan(context.Context, identity.Principal, sc.PlanRequest) (sc.Receipt, error)
	Apply(context.Context, identity.Principal, sc.Command) (sc.Receipt, error)
	Plan(context.Context, identity.Principal, string, string) (sc.Plan, error)
	CommandReceipt(context.Context, identity.Principal, string, string) (sc.Receipt, error)
}
type serialSupplyHook func(*testing.T, *pgxpool.Pool, *db.SerialSupply) serialSupplyReference

func TestSerialSupplyConnectedReference(t *testing.T) { testSerialSupplyConnected(t, nil) }
func testSerialSupplyConnected(t *testing.T, hook serialSupplyHook) {
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
	tenant := "ca6c9cb7-e548-4f19-a2e6-27d2aa2d4892"
	fixtures := []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,'supply-j2','Fixture','Fixture')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,'j2-store','j2-store','Fixture','store'),($1,'j2-factory','j2-factory','Fixture','factory')`,
		`insert into partner.supplier(tenant_id,supplier_id,supplier_code,legal_name,status) values($1,'j2-supplier','j2-supplier','Fixture','active')`,
		`insert into catalog.vehicle_model(tenant_id,model_id,model_code,display_name,vehicle_class,lifecycle_state) values($1,'j2-model','j2-model','Fixture','bicycle','active')`,
		`insert into catalog.vehicle_variant(tenant_id,variant_id,model_id,variant_code,display_name,battery_specification,lifecycle_state) values($1,'j2-a','j2-model','j2-a','Fixture','{}','active'),($1,'j2-b','j2-model','j2-b','Fixture','{}','active')`,
	}
	for _, q := range fixtures {
		if _, err = pool.Exec(ctx, q, tenant); err != nil {
			t.Fatal(err)
		}
	}
	person := func(subject, org string, permissions ...string) identity.Principal {
		p := identity.Principal{Subject: subject, TenantID: tenant, Permissions: map[string]struct{}{}, Organizations: map[string]struct{}{org: {}}}
		for _, v := range permissions {
			p.Permissions[v] = struct{}{}
		}
		return p
	}
	buyer := person("j2-buyer", "j2-store", "supply:plan", "supply:read")
	maker := person("j2-maker", "j2-factory", "supply:factory", "supply:factory-read")
	factoryQA := person("j2-factory-qa", "j2-factory", "supply:factory", "supply:factory-read")
	receiver := person("j2-receiver", "j2-store", "supply:receive", "supply:inspect", "supply:read")
	qa := person("j2-receipt-qa", "j2-store", "supply:release", "supply:read")
	store, err := db.NewSerialSupply(pool)
	if err != nil {
		t.Fatal(err)
	}
	var supply serialSupplyReference = store
	if hook != nil {
		supply = hook(t, pool, store)
	}
	ops := operations.NewService(db.NewOperations(pool), randomid.Generator{})
	po, err := ops.CreatePurchaseOrder(ctx, tenant, operations.PurchaseOrder{SupplierID: "j2-supplier", DestinationOrganizationID: "j2-store", Currency: "ARS", TotalMinorUnits: 90000})
	if err != nil {
		t.Fatal(err)
	}
	planRequest := sc.PlanRequest{PurchaseOrderID: po.ID, CommandID: "plan", FactoryOrganizationID: "j2-factory", DemandReference: "fixture-demand-3-serials", PolicyCode: sc.PolicyCode, EvidenceSHA256: strings.Repeat("a", 64), Lines: []sc.Line{{ID: "a", VariantID: "j2-a", Quantity: 2}, {ID: "b", VariantID: "j2-b", Quantity: 1}}}
	first, err := supply.BindPlan(ctx, buyer, planRequest)
	if err != nil {
		t.Fatal("bind plan", err)
	}
	replay, err := supply.BindPlan(ctx, buyer, planRequest)
	if err != nil || !replay.Replay || replay.RequestSHA256 != first.RequestSHA256 {
		t.Fatal("plan replay", replay, err)
	}
	wrong := person("j2-other", "other", "supply:plan", "supply:read", "supply:receive")
	if _, err = supply.Plan(ctx, wrong, po.ID, ""); err == nil {
		t.Fatal("cross organization plan read")
	}
	version := first.Version
	sequence := 0
	command := func(kind string) sc.Command {
		sequence++
		return sc.Command{PurchaseOrderID: po.ID, CommandID: fmt.Sprintf("j2-%03d", sequence), ExpectedVersion: version, Kind: kind, EvidenceSHA256: strings.Repeat("a", 64)}
	}
	apply := func(p identity.Principal, r sc.Command) sc.Receipt {
		t.Helper()
		out, err := supply.Apply(ctx, p, r)
		if err != nil {
			t.Fatalf("%s/%s: %v", r.Kind, r.CommandID, err)
		}
		if out.Replay || out.Version != version+1 || out.Actor != p.Subject {
			t.Fatalf("unexpected receipt: %+v", out)
		}
		version = out.Version
		return out
	}
	negative := func(p identity.Principal, r sc.Command, label string) {
		t.Helper()
		if _, err := supply.Apply(ctx, p, r); err == nil {
			t.Fatal(label + " accepted")
		}
	}
	snapshot := func() string {
		t.Helper()
		out := ""
		for _, table := range []string{"procurement.purchase_order", "factory.production_unit", "inventory.stock_unit", "procurement.serial_supply_plan", "procurement.serial_supply_line", "procurement.serial_supply_unit", "procurement.serial_supply_shipment", "procurement.serial_supply_manifest", "procurement.serial_supply_receipt", "procurement.serial_supply_step", "procurement.serial_supply_quality", "procurement.serial_supply_effect", "approval.request", "approval.decision", "platform.outbox_event"} {
			var raw string
			query := fmt.Sprintf("select coalesce(jsonb_agg(to_jsonb(x) order by to_jsonb(x)::text),'[]'::jsonb)::text from %s x where tenant_id=$1", table)
			if err := pool.QueryRow(ctx, query, tenant).Scan(&raw); err != nil {
				t.Fatal(table, err)
			}
			out += table + raw
		}
		return out
	}
	apply(buyer, command("submit"))
	negative(buyer, command("confirm"), "buyer acting as factory")
	apply(maker, command("confirm"))
	apply(maker, command("start"))
	register := func(line, serial string) string {
		t.Helper()
		r := command("register")
		r.LineID = line
		r.SerialNumber = serial
		out := apply(maker, r)
		var payload struct {
			Effect struct {
				Unit operations.ProductionUnit `json:"unit"`
			} `json:"effect"`
		}
		if err := json.Unmarshal(out.Payload, &payload); err != nil || payload.Effect.Unit.ID == "" {
			t.Fatal(err, string(out.Payload))
		}
		return payload.Effect.Unit.ID
	}
	milestone := func(actor identity.Principal, unit, target string) sc.Receipt {
		t.Helper()
		r := command("milestone")
		r.UnitID = unit
		r.TargetState = target
		if target == "released" || target == "rejected" {
			r.Reason = "fixture inspection and recorded evidence"
		}
		return apply(actor, r)
	}
	// A rejected factory serial is retained and a replacement consumes its line slot.
	rejected := register("a", "J2-REJECTED")
	milestone(maker, rejected, "assembly")
	milestone(maker, rejected, "quality")
	before := snapshot()
	r := command("milestone")
	r.UnitID = rejected
	r.TargetState = "released"
	r.Reason = "self review forbidden"
	negative(maker, r, "factory self approval")
	if snapshot() != before {
		t.Fatal("self approval changed durable state")
	}
	milestone(factoryQA, rejected, "rejected")
	a1 := register("a", "J2-A1")
	a2 := register("a", "J2-A2")
	b1 := register("b", "J2-B1")
	r = command("register")
	r.LineID = "a"
	r.SerialNumber = "J2-OVER"
	negative(maker, r, "overproduction")
	r = command("register")
	r.LineID = "b"
	r.SerialNumber = "J2-A1"
	negative(maker, r, "reused serial/quantity")
	r = command("ship")
	r.ShipmentID = "premature"
	r.Units = []string{a1}
	negative(maker, r, "premature shipment")
	// Legacy operations may not bypass the connected plan even with the right org.
	before = snapshot()
	if err = ops.TransitionProductionUnit(ctx, tenant, "j2-store", a1, "planned", "assembly"); err == nil {
		t.Fatal("generic factory bypass")
	}
	if snapshot() != before {
		t.Fatal("rejected generic command leaked state/outbox")
	}
	for _, id := range []string{a1, a2, b1} {
		milestone(maker, id, "assembly")
		milestone(maker, id, "quality")
		milestone(factoryQA, id, "released")
	}
	r = command("ship")
	r.ShipmentID = "asn-one"
	r.Units = []string{a1, a2}
	apply(maker, r)
	plan, err := supply.Plan(ctx, buyer, po.ID, "")
	if err != nil || plan.State != "in-production" {
		t.Fatal("partial shipment order state", plan.State, err)
	}
	r = command("receive")
	r.ShipmentID = "asn-one"
	r.Units = []string{a1}
	negative(maker, r, "factory receiving as destination")
	negative(wrong, r, "wrong receiving organization")
	apply(receiver, r)
	if hook != nil {
		before := snapshot()
		var approvalID, approvalHash, stockID string
		var stockVersion int64
		err = pool.QueryRow(ctx, `select q.approval_id,q.payload_sha256,m.stock_unit_id,i.version from procurement.serial_supply_quality q
   join procurement.serial_supply_manifest m using(tenant_id,production_unit_id)
   join inventory.stock_unit i on i.tenant_id=m.tenant_id and i.stock_unit_id=m.stock_unit_id
   where q.tenant_id=$1 and q.production_unit_id=$2 and q.stage='receipt' order by q.attempt desc limit 1`, tenant, a1).Scan(&approvalID, &approvalHash, &stockID, &stockVersion)
		if err != nil {
			t.Fatal(err)
		}
		if _, err = db.NewHumanApprovals(pool).Decide(ctx, qa, tenant, approvalID, "j2-store", approvalHash, true, "generic decision bypass", "supply:release", nil); err == nil {
			t.Fatal("generic approval bypass")
		}
		if err = ops.TransitionStockUnit(ctx, tenant, "j2-store", stockID, "quarantine", "available", stockVersion); err == nil {
			t.Fatal("generic stock release bypass")
		}
		if _, err = pool.Exec(ctx, `update procurement.purchase_order set total_minor_units=total_minor_units+1 where tenant_id=$1 and purchase_order_id=$2`, tenant, po.ID); err == nil {
			t.Fatal("bound purchase amount changed")
		}
		if snapshot() != before {
			t.Fatal("approval/stock/purchase bypass changed durable state")
		}
	}
	// Newly received stock is quarantine, not ATP.
	atp, err := db.NewInventoryControl(pool).AvailableToPromise(ctx, tenant, "j2-store", "j2-a", time.Now().UTC().Add(time.Hour))
	if err != nil || atp.AvailableInventory != 0 {
		t.Fatal("quarantine ATP", atp, err)
	}
	r = command("quality")
	r.UnitID = a1
	r.Reason = "fixture receipt accepted"
	self := receiver
	self.Permissions = map[string]struct{}{"supply:release": {}}
	before = snapshot()
	negative(self, r, "receiver self quality approval")
	if snapshot() != before {
		t.Fatal("receiver self decision leaked")
	}
	apply(qa, r)
	atp, err = db.NewInventoryControl(pool).AvailableToPromise(ctx, tenant, "j2-store", "j2-a", time.Now().UTC().Add(time.Hour))
	if err != nil || atp.AvailableInventory != 1 {
		t.Fatal("released ATP", atp, err)
	}
	r = command("ship")
	r.ShipmentID = "asn-two"
	r.Units = []string{b1}
	apply(maker, r)
	plan, err = supply.Plan(ctx, buyer, po.ID, "")
	if err != nil || plan.State != "shipped" {
		t.Fatal("all shipped order state", plan.State, err)
	}
	// Fail after state changes, pending approval and connected receipt were written.
	_, err = pool.Exec(ctx, `create function public.j2_outbox_fail() returns trigger language plpgsql as $$
 begin if new.aggregate_type='serial-supply' and new.payload->>'command_id'='receipt-rollback' then raise exception 'fixture outbox rollback';end if;return new;end $$;
 create trigger j2_outbox_fail before insert on platform.outbox_event for each row execute function public.j2_outbox_fail();`)
	if err != nil {
		t.Fatal(err)
	}
	r = command("receive")
	r.CommandID = "receipt-rollback"
	r.ShipmentID = "asn-one"
	r.Units = []string{a2}
	before = snapshot()
	negative(receiver, r, "outbox failure")
	if snapshot() != before {
		t.Fatal("outbox rollback changed supply/stock/approval history")
	}
	if _, err = pool.Exec(ctx, `drop trigger j2_outbox_fail on platform.outbox_event;drop function public.j2_outbox_fail();`); err != nil {
		t.Fatal(err)
	}
	// Concurrent same command: one new effect; other calls either replay or conflict.
	r = command("receive")
	r.ShipmentID = "asn-one"
	r.Units = []string{a2}
	parallel := r
	var wg sync.WaitGroup
	var mu sync.Mutex
	fresh, replays, conflicts := 0, 0, 0
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			out, err := supply.Apply(ctx, receiver, parallel)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				conflicts++
				return
			}
			if out.Replay {
				replays++
			} else {
				fresh++
			}
			if out.Actor != receiver.Subject || out.Version != version+1 {
				t.Errorf("bad concurrent receipt %+v", out)
			}
		}()
	}
	wg.Wait()
	if fresh != 1 || fresh+replays+conflicts != 8 {
		t.Fatalf("concurrency fresh/replay/conflict=%d/%d/%d", fresh, replays, conflicts)
	}
	version++
	// Recover committed evidence from a new pool; no hidden POST retry.
	otherPool, err := pgxpool.New(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	restarted, err := db.NewSerialSupply(otherPool)
	if err != nil {
		t.Fatal(err)
	}
	recovered, err := restarted.CommandReceipt(ctx, receiver, po.ID, parallel.CommandID)
	otherPool.Close()
	if err != nil || recovered.Actor != receiver.Subject || recovered.Version != version {
		t.Fatal("recovery", recovered, err)
	}
	_, wantHash, err := sc.Canonical(parallel)
	if err != nil || wantHash != recovered.RequestSHA256 {
		t.Fatal("recovery request hash", err)
	}
	again, err := supply.Apply(ctx, receiver, parallel)
	if err != nil || !again.Replay {
		t.Fatal("exact replay", again, err)
	}
	changed := parallel
	changed.EvidenceSHA256 = strings.Repeat("c", 64)
	negative(receiver, changed, "changed replay payload")
	secondReceiver := person("j2-second-receiver", "j2-store", "supply:receive")
	negative(secondReceiver, parallel, "changed replay actor")
	r = command("receive")
	r.ShipmentID = "asn-one"
	r.Units = []string{b1}
	negative(receiver, r, "wrong ASN membership")
	r = command("receive")
	r.ShipmentID = "asn-two"
	r.Units = []string{b1}
	apply(receiver, r)
	plan, err = supply.Plan(ctx, buyer, po.ID, "")
	if err != nil || plan.State != "received" || len(plan.Units) != 4 || len(plan.Lines) != 2 {
		t.Fatal("final physical receipt", plan, err)
	}
	r = command("quality-reject")
	r.UnitID = a2
	r.Reason = "fixture inspection failed"
	apply(qa, r)
	r = command("reinspect")
	r.UnitID = a2
	r.Reason = "new inspection"
	negative(receiver, r, "unchanged rejected evidence")
	r.EvidenceSHA256 = strings.Repeat("b", 64)
	apply(receiver, r)
	r = command("quality")
	r.UnitID = a2
	r.Reason = "corrected fixture evidence accepted"
	apply(qa, r)
	r = command("quality")
	r.UnitID = b1
	r.Reason = "fixture evidence accepted"
	apply(qa, r)
	atp, err = db.NewInventoryControl(pool).AvailableToPromise(ctx, tenant, "j2-store", "j2-a", time.Now().UTC().Add(time.Hour))
	if err != nil || atp.AvailableInventory != 2 {
		t.Fatal("final A ATP", atp, err)
	}
	bATP, err := db.NewInventoryControl(pool).AvailableToPromise(ctx, tenant, "j2-store", "j2-b", time.Now().UTC().Add(time.Hour))
	if err != nil || bATP.AvailableInventory != 1 {
		t.Fatal("final B ATP", bATP, err)
	}
	var manifested, receivedCount, available, poReceived int
	err = pool.QueryRow(ctx, `select (select count(*) from procurement.serial_supply_manifest where tenant_id=$1),(select count(*) from procurement.serial_supply_receipt where tenant_id=$1),(select count(*) from inventory.stock_unit where tenant_id=$1 and state='available'),(select count(*) from procurement.purchase_order where tenant_id=$1 and state='received')`, tenant).Scan(&manifested, &receivedCount, &available, &poReceived)
	if err != nil || manifested != 3 || receivedCount != 3 || available != 3 || poReceived != 1 {
		t.Fatal("counts", manifested, receivedCount, available, poReceived, err)
	}
	before = snapshot()
	if _, err = pool.Exec(ctx, `update procurement.serial_supply_line set quantity=4 where tenant_id=$1`, tenant); err == nil {
		t.Fatal("line evidence mutated")
	}
	if snapshot() != before {
		t.Fatal("immutability changed history")
	}
	t.Logf("SERIAL_SUPPLY_CONNECTED_PASS planned3/registered4/rejected1/shipped3/received3/available3; splitASN2; factory/receipt review; rollback/recovery; concurrency%dnew/%dreplay/%dconflict; version%d", fresh, replays, conflicts, version)
}
