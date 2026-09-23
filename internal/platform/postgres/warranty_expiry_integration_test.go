package postgres_test

// Narrow commit-time lease regression. Accepted handover and attended service
// are explicit historical fixtures here; the separate HTTP J1/J4 test proves
// their full creation journey. The approval/reservation transaction is real.
import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"elite.local/enterprise/internal/inventorycontrol"
	"elite.local/enterprise/internal/platform/identity"
	db "elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/platform/randomid"
	wc "elite.local/enterprise/internal/warrantyclaim"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestWarrantyPendingLeaseExpiresAtCommit(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 85*time.Second)
	defer cancel()
	pool := connectedPool(t)
	ids := randomid.Generator{}
	var store *db.Warranty
	var actor identity.Principal
	tenant := connectedSeedOrder(t, pool, "stripe", func(t *testing.T, pool *pgxpool.Pool, tenant string) int64 {
		var profile wc.Profile
		store, actor, profile = fixtureWarranty(t, pool, tenant, 365, 60)
		offered, err := store.BindOffer(ctx, actor, wc.OfferRequest{QuoteID: "quote", QuoteVersion: 1, ProfileSHA256: profile.Hash()})
		if err != nil {
			t.Fatal(err)
		}
		customer := identity.Principal{TenantID: tenant, Subject: "customer", Organizations: map[string]struct{}{"store": {}}, Permissions: map[string]struct{}{"warranty:self": {}}}
		if _, err = store.Acknowledge(ctx, customer, wc.AcknowledgeRequest{QuoteID: "quote", QuoteVersion: offered.QuoteVersion, ProfileSHA256: profile.Hash(), EvidenceSHA256: strings.Repeat("a", 64)}); err != nil {
			t.Fatal(err)
		}
		return offered.QuoteVersion
	})
	must := func(sql string, args ...any) {
		t.Helper()
		if _, err := pool.Exec(ctx, sql, args...); err != nil {
			t.Fatal(err)
		}
	}
	must(`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'factory','factory','Fixture factory','factory')`, tenant)
	must(`insert into sales.delivery_handover(tenant_id,handover_id,organization_id,order_id,customer_principal_id,stock_unit_id,state,acceptance_evidence_sha256_hex,customer_accepted_at,version)
 values($1,'historical-handover','store','order','customer','stock','accepted',$2,clock_timestamp(),3)`, tenant, strings.Repeat("b", 64))
	if _, err := store.Activate(ctx, actor, "historical-handover"); err != nil {
		t.Fatal(err)
	}
	must(`insert into crm.appointment(tenant_id,appointment_id,organization_id,lead_id,customer_principal_id,model_id,appointment_kind,starts_at,state,version,created_at,updated_at)
 values($1,'historical-service','store','lead','customer','model','service',clock_timestamp()-interval '30 minutes','completed',3,clock_timestamp()-interval '1 hour',clock_timestamp())`, tenant)
	for _, p := range []string{"warranty:request", "warranty:diagnose", "warranty:plan", "warranty:approve"} {
		actor.Permissions[p] = struct{}{}
	}
	review := actor
	review.Subject = "reviewer"
	inv := db.NewInventoryControl(pool)
	if _, err := inv.CreateBulkItem(ctx, tenant, ids.New(), inventorycontrol.BulkItem{ID: "part", Code: "PART", Description: "Fixture part", BaseUOM: "EA", BaseRoundingPrecision: "1", TrackingMode: "none", CostingMethod: "fifo", Version: 1}); err != nil {
		t.Fatal(err)
	}
	if _, err := inv.CreateWarehouseBin(ctx, tenant, ids.New(), inventorycontrol.WarehouseBin{ID: "bin", OrganizationID: "store", Code: "BIN", Type: "pick", Version: 1}); err != nil {
		t.Fatal(err)
	}
	if _, err := inv.ConfigureItemBin(ctx, tenant, ids.New(), inventorycontrol.ItemBinPolicy{OrganizationID: "store", ItemID: "part", BinID: "bin", MinQuantity: "0", Version: 1}); err != nil {
		t.Fatal(err)
	}
	if _, err := inv.ReceiveBulk(ctx, tenant, ids.New(), ids.New(), ids.New(), ids.New(), inventorycontrol.BulkReceipt{OrganizationID: "store", BinID: "bin", ItemID: "part", Quantity: "1", UnitCost: "10", PostingDate: time.Now().UTC(), SourceKind: "fixture", SourceID: "expiry-stock"}); err != nil {
		t.Fatal(err)
	}
	if _, err := store.OpenClaim(ctx, actor, wc.OpenClaim{CaseID: "expiry-case", HandoverID: "historical-handover", AppointmentID: "historical-service", Severity: "low", Description: "Fixture lease inspection", EvidenceSHA256: strings.Repeat("c", 64)}); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Diagnose(ctx, actor, wc.Diagnose{Command: wc.Command{CaseID: "expiry-case", CommandID: "diagnosis", ExpectedVersion: 1}, FaultCode: "fixture-wear", Description: "Fixture covered fault", EvidenceSHA256: strings.Repeat("d", 64)}); err != nil {
		t.Fatal(err)
	}
	plan, err := store.PlanRepair(ctx, actor, wc.Plan{Command: wc.Command{CaseID: "expiry-case", CommandID: "plan", ExpectedVersion: 2}, Parts: []wc.Part{{LineID: "line", ItemID: "part", BinID: "bin", Quantity: "1"}}, LaborWork: ""})
	if err != nil {
		t.Fatal(err)
	}
	var bound struct {
		ID  string `json:"approval_id"`
		SHA string `json:"payload_sha256"`
	}
	if json.Unmarshal(plan.Payload, &bound) != nil {
		t.Fatal("plan payload")
	}
	must(`create function service_ops.fixture_delay_warranty_commit() returns trigger language plpgsql as $$declare expiry timestamptz;begin
 if new.event_type='warranty.claim-decision' then
 select expires_at into expiry from service_ops.warranty_work_plan where tenant_id=new.tenant_id and service_case_id=new.aggregate_id;
 perform pg_sleep(greatest(0,extract(epoch from(expiry-clock_timestamp())))+0.1);
 end if;return new;end $$;
 create trigger fixture_delay_warranty_commit before insert on platform.outbox_event for each row execute function service_ops.fixture_delay_warranty_commit()`)
	t.Log("WARRANTY_EXPIRY_PROBE waiting only for the explicit60second pending lease to cross COMMIT; no shortened production policy")
	_, err = store.DecideRepair(ctx, review, wc.DecideRepair{Command: wc.Command{CaseID: "expiry-case", CommandID: "approve", ExpectedVersion: 3}, PayloadSHA256: bound.SHA, Approved: true, Reason: "Fixture approval before lease expiry"})
	if err == nil || ctx.Err() != nil || !strings.Contains(err.Error(), "warranty pending reservation expired before commit") {
		t.Fatal("commit did not reject exact lease expiry", err, ctx.Err())
	}
	must(`drop trigger fixture_delay_warranty_commit on platform.outbox_event;drop function service_ops.fixture_delay_warranty_commit()`)
	var decisions, steps, events, issues int
	var state, reservationState string
	var version int64
	var held string
	err = pool.QueryRow(ctx, `select (select count(*) from approval.decision where tenant_id=$1),(select count(*) from service_ops.warranty_claim_step where tenant_id=$1 and kind='decision'),(select count(*) from platform.outbox_event where tenant_id=$1 and event_type='warranty.claim-decision'),(select count(*) from inventory.bulk_inventory_entry where tenant_id=$1 and entry_type='issue'),(select state from approval.request where tenant_id=$1 and request_id=$2),(select status from inventory.bulk_reservation where tenant_id=$1),(select version from inventory.bulk_reservation where tenant_id=$1),(select reserved_quantity::text from inventory.bulk_balance where tenant_id=$1)`, tenant, bound.ID).Scan(&decisions, &steps, &events, &issues, &state, &reservationState, &version, &held)
	if err != nil || decisions != 0 || steps != 0 || events != 0 || issues != 0 || state != "pending" || reservationState != "reservation" || version != 1 || held != "1.000000" {
		t.Fatal("expired approval leaked effect", decisions, steps, events, issues, state, reservationState, version, held, err)
	}
	rejected, err := store.DecideRepair(ctx, review, wc.DecideRepair{Command: wc.Command{CaseID: "expiry-case", CommandID: "reject-expired", ExpectedVersion: 3}, PayloadSHA256: bound.SHA, Approved: false, Reason: "Expired pending lease released"})
	if err != nil || rejected.State != "cancelled" {
		t.Fatal("expired hold recovery", rejected, err)
	}
	if err = pool.QueryRow(ctx, `select reserved_quantity::text from inventory.bulk_balance where tenant_id=$1`, tenant).Scan(&held); err != nil || held != "0.000000" {
		t.Fatal("hold not released", held, err)
	}
	t.Log("WARRANTY_EXPIRY_COMMIT_PASS approved_before_wait_commit_after_expiry=true leaked_decisions_steps_events_issues=0 original_pending_reservation_preserved=true explicit_rejection_releases_hold=true")
}
