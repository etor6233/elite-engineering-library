package postgres_test

import (
	"context"
	"encoding/json"
	"math/big"
	"strings"
	"testing"
	"time"

	"elite.local/enterprise/internal/businesspolicy"
	"elite.local/enterprise/internal/franchisejourney"
	"elite.local/enterprise/internal/inventorycontrol"
	"elite.local/enterprise/internal/platform/identity"
	db "elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/platform/randomid"
	wc "elite.local/enterprise/internal/warrantyclaim"
)

type warrantyFixtureClock struct{}

func (warrantyFixtureClock) Now() time.Time { return time.Now().UTC() }

func assertWarrantyClaimLifecycle(t *testing.T, r *connectedRun, store warrantyClaimCommands, operator, customer identity.Principal, activation wc.Activation) {
	t.Helper()
	ctx := context.Background()
	ids := randomid.Generator{}
	must := func(q string, args ...any) {
		t.Helper()
		if _, err := r.pool.Exec(ctx, q, args...); err != nil {
			t.Fatal(err)
		}
	}
	must(`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'factory','factory','Synthetic factory','factory')`, r.tenant)
	must(`insert into org.public_location(tenant_id,organization_id,city,region,country,published)values($1,'store','Synthetic','Synthetic','AR',true)`, r.tenant)
	// Explicit synthetic appointment profile, admitted by the existing policy
	// contract: zero lead time and one-second slots avoid wall-clock simulation.
	// This operational profile does not amend any warranty terms already sold.
	var doc map[string]any
	if json.Unmarshal(businesspolicy.ReferenceJSON(), &doc) != nil {
		t.Fatal("reference profile")
	}
	doc["profile_id"] = "warranty-fixture"
	doc["appointments"].(map[string]any)["lead_time_seconds"] = 0
	raw, err := json.Marshal(doc)
	if err != nil {
		t.Fatal(err)
	}
	policy, err := businesspolicy.Load(raw, wc.SHA(raw))
	if err != nil {
		t.Fatal(err)
	}
	journey, err := db.NewFranchiseJourneyWithProfile(r.pool, policy)
	if err != nil {
		t.Fatal(err)
	}
	booking, err := franchisejourney.NewServiceWithProfile(journey, ids, warrantyFixtureClock{}, policy)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = journey.CreateServiceResource(ctx, r.tenant, franchisejourney.ServiceResource{ID: "warranty-bay", OrganizationID: "store", Kind: "service-bay", DisplayName: "Synthetic service bay", Skills: []string{"service"}}, ids.New()); err != nil {
		t.Fatal(err)
	}
	for _, resource := range []string{"", "warranty-bay"} {
		must(`insert into crm.availability_entry(tenant_id,availability_id,organization_id,resource_id,entry_type,starts_at,ends_at,state,version,created_by_subject)
 values($1,$2,'store',nullif($3,''),'working',clock_timestamp()-interval '1 hour',clock_timestamp()+interval '6 hours','active',1,'fixture-calendar')`, r.tenant, ids.New(), resource)
	}
	var tenantCode string
	if err = r.pool.QueryRow(ctx, `select tenant_code from platform.tenant where tenant_id=$1`, r.tenant).Scan(&tenantCode); err != nil {
		t.Fatal(err)
	}
	appointment := func(id string) string {
		t.Helper()
		var start time.Time
		if err := r.pool.QueryRow(ctx, `select clock_timestamp()+interval '1 second'`).Scan(&start); err != nil {
			t.Fatal(err)
		}
		if _, err := booking.CreateAppointmentSlot(ctx, r.tenant, franchisejourney.AppointmentSlot{OrganizationID: "store", Kind: "service", StartsAt: start, EndsAt: start.Add(time.Second), Capacity: 1}); err != nil {
			t.Fatal(err)
		}
		a, _, err := booking.RequestAppointment(ctx, tenantCode, "store", id+"-appointment-booking-key", wc.SHA([]byte(id)), franchisejourney.Appointment{LeadID: "lead", ModelID: "model", Kind: "service", StartsAt: start})
		if err != nil {
			t.Fatal("request appointment", err)
		}
		a, err = journey.AssignAppointmentResource(ctx, r.tenant, "store", a.ID, "warranty-bay", a.Version, ids.New())
		if err != nil {
			t.Fatal("assign service bay", err)
		}
		a, err = journey.TransitionAppointment(ctx, r.tenant, "store", a.ID, "requested", "confirmed", a.Version, "operator", "", ids.New())
		if err != nil {
			t.Fatal("confirm appointment", err)
		}
		wait := time.Until(start.Add(time.Second + 30*time.Millisecond))
		if wait > 0 {
			if wait > 3*time.Second {
				t.Fatal("fixture wait exceeds bound")
			}
			time.Sleep(wait)
		}
		a, err = journey.TransitionAppointment(ctx, r.tenant, "store", a.ID, "confirmed", "completed", a.Version, "operator", "", ids.New())
		if err != nil {
			t.Fatal("complete attended appointment", err)
		}
		return a.ID
	}
	inv := db.NewInventoryControl(r.pool)
	if _, err = inv.CreateBulkItem(ctx, r.tenant, ids.New(), inventorycontrol.BulkItem{ID: "warranty-part", Code: "WARRANTY-PART", Description: "Synthetic replacement part", BaseUOM: "EA", BaseRoundingPrecision: "1", TrackingMode: "lot", CostingMethod: "fifo", Version: 1}); err != nil {
		t.Fatal(err)
	}
	if _, err = inv.CreateWarehouseBin(ctx, r.tenant, ids.New(), inventorycontrol.WarehouseBin{ID: "warranty-bin", OrganizationID: "store", Code: "SERVICE", Type: "pick", Version: 1}); err != nil {
		t.Fatal(err)
	}
	if _, err = inv.ConfigureItemBin(ctx, r.tenant, ids.New(), inventorycontrol.ItemBinPolicy{OrganizationID: "store", ItemID: "warranty-part", BinID: "warranty-bin", Fixed: true, Default: true, MinQuantity: "0", MaxQuantity: "100", Version: 1}); err != nil {
		t.Fatal(err)
	}
	expires := time.Now().UTC().AddDate(2, 0, 0)
	for i, cost := range []string{"100", "120"} {
		_, err = inv.ReceiveBulk(ctx, r.tenant, ids.New(), "warranty-receipt-"+cost, "warranty-layer-"+cost, "warranty-lot", inventorycontrol.BulkReceipt{OrganizationID: "store", BinID: "warranty-bin", ItemID: "warranty-part", LotNo: "SERVICE-FIXTURE", ExpirationDate: &expires, Quantity: "5", UnitCost: cost, PostingDate: time.Now().UTC().AddDate(0, 0, -2+i), SourceKind: "fixture-receipt", SourceID: "fixture-" + cost})
		if err != nil {
			t.Fatal(err)
		}
	}
	for _, permission := range []string{"warranty:request", "warranty:diagnose", "warranty:plan", "warranty:work", "warranty:approve", "warranty:quality", "warranty:cancel"} {
		operator.Permissions[permission] = struct{}{}
	}
	reviewer := operator
	reviewer.Subject = "warranty-reviewer"
	qualityActor := operator
	qualityActor.Subject = "warranty-quality"
	factory := identity.Principal{TenantID: r.tenant, Subject: "factory-reviewer", Organizations: map[string]struct{}{"factory": {}}, Permissions: map[string]struct{}{"warranty:factory-read": {}, "warranty:reconcile": {}}}
	part := func(qty string) []wc.Part {
		return []wc.Part{{LineID: "brake", ItemID: "warranty-part", BinID: "warranty-bin", LotID: "warranty-lot", Quantity: qty}}
	}
	open := func(id, fault string) (wc.Step, wc.Diagnose) {
		t.Helper()
		request := wc.OpenClaim{CaseID: id, HandoverID: activation.HandoverID, AppointmentID: appointment(id), Severity: "medium", Description: "Synthetic brake diagnosis request", EvidenceSHA256: strings.Repeat("a", 64)}
		opened, err := store.OpenClaim(ctx, customer, request)
		if err != nil {
			t.Fatal("open claim", err)
		}
		again, err := store.OpenClaim(ctx, customer, request)
		if err != nil || !again.Replay || again.Version != opened.Version {
			t.Fatal("open recovery", again, err)
		}
		diagnose := wc.Diagnose{Command: wc.Command{CaseID: id, CommandID: "diagnose", ExpectedVersion: 1}, FaultCode: fault, Description: "Synthetic inspection of brake", EvidenceSHA256: strings.Repeat("b", 64)}
		step, err := store.Diagnose(ctx, operator, diagnose)
		if err != nil {
			t.Fatal("diagnosis", err)
		}
		return step, diagnose
	}
	planned := func(id, qty string, version int64) (wc.Step, string) {
		t.Helper()
		step, err := store.PlanRepair(ctx, operator, wc.Plan{Command: wc.Command{CaseID: id, CommandID: "plan", ExpectedVersion: version}, LaborWork: "Synthetic brake adjustment", Parts: part(qty)})
		if err != nil {
			t.Fatal("plan", err)
		}
		var bound struct {
			SHA string `json:"payload_sha256"`
		}
		if json.Unmarshal(step.Payload, &bound) != nil || !wc.ValidSHA(bound.SHA) {
			t.Fatal("plan hash")
		}
		return step, bound.SHA
	}
	step, _ := open("claim-main", "fixture-wear")
	step, planSHA := planned("claim-main", "7", step.Version)
	decision := wc.DecideRepair{Command: wc.Command{CaseID: "claim-main", CommandID: "approve", ExpectedVersion: step.Version}, PayloadSHA256: planSHA, Approved: true, Reason: "Fixture covered repair"}
	if _, err = store.DecideRepair(ctx, operator, decision); err == nil {
		t.Fatal("requester approved own repair")
	}
	// The old generic case endpoint cannot bypass the connected approval writer.
	if err = db.NewFulfillment(r.pool).TransitionServiceCase(ctx, r.tenant, "store", "claim-main", "diagnosis", "repair", step.Version, ids.New()); err == nil {
		t.Fatal("generic service path bypassed approval")
	}
	step, err = store.DecideRepair(ctx, reviewer, decision)
	if err != nil || step.State != "repair" {
		t.Fatal("human approval", step, err)
	}
	again, err := store.DecideRepair(ctx, reviewer, decision)
	if err != nil || !again.Replay {
		t.Fatal("decision recovery", again, err)
	}
	work := wc.CompleteWork{Command: wc.Command{CaseID: "claim-main", CommandID: "work", ExpectedVersion: step.Version}, EvidenceSHA256: strings.Repeat("c", 64)}
	// Failure after real FIFO stock writes must roll back every stock/case effect.
	must(`create function service_ops.fixture_work_failure() returns trigger language plpgsql as $$begin if new.event_type='warranty.claim-work' then raise exception 'fixture work commit failure';end if;return new;end $$;
 create trigger fixture_work_failure before insert on platform.outbox_event for each row execute function service_ops.fixture_work_failure()`)
	if _, err = store.CompleteRepairWork(ctx, operator, work); err == nil {
		t.Fatal("work fault injection did not fail")
	}
	var quantity, reserved string
	var issues int
	if err = r.pool.QueryRow(ctx, `select quantity::text,reserved_quantity::text,(select count(*) from inventory.bulk_inventory_entry where tenant_id=$1 and entry_type='issue') from inventory.bulk_balance where tenant_id=$1 and item_id='warranty-part'`, r.tenant).Scan(&quantity, &reserved, &issues); err != nil || quantity != "10.000000" || reserved != "7.000000" || issues != 0 {
		t.Fatal("FIFO rollback", quantity, reserved, issues, err)
	}
	must(`drop trigger fixture_work_failure on platform.outbox_event;drop function service_ops.fixture_work_failure()`)
	step, err = store.CompleteRepairWork(ctx, operator, work)
	if err != nil || step.State != "quality" {
		t.Fatal("complete repair", step, err)
	}
	var performed wc.WorkReceipt
	if json.Unmarshal(step.Payload, &performed) != nil || len(performed.Parts) != 1 || !warrantyExactCost(performed.Parts[0].Issue.CostAmount, "740") || performed.Parts[0].Issue.Applications != 2 {
		t.Fatal("source FIFO receipt", string(step.Payload))
	}
	workSHA := step.PayloadSHA256
	again, err = store.CompleteRepairWork(ctx, operator, work)
	if err != nil || !again.Replay || again.PayloadSHA256 != workSHA {
		t.Fatal("work recovery", again, err)
	}
	failedQuality := wc.Quality{Command: wc.Command{CaseID: "claim-main", CommandID: "quality-failed", ExpectedVersion: step.Version}, Passed: false, WorkEvidenceSHA256: workSHA, EvidenceSHA256: strings.Repeat("d", 64)}
	if _, err = store.RecordRepairQuality(ctx, operator, failedQuality); err == nil {
		t.Fatal("technician self-verified work")
	}
	step, err = store.RecordRepairQuality(ctx, qualityActor, failedQuality)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = store.AcceptRepair(ctx, customer, wc.AcceptRepair{Command: wc.Command{CaseID: "claim-main", CommandID: "accept-failed", ExpectedVersion: step.Version}, QualitySHA256: step.PayloadSHA256, EvidenceSHA256: strings.Repeat("e", 64)}); err == nil {
		t.Fatal("failed quality accepted")
	}
	corrected := wc.Quality{Command: wc.Command{CaseID: "claim-main", CommandID: "quality-passed", ExpectedVersion: step.Version}, Passed: true, WorkEvidenceSHA256: workSHA, EvidenceSHA256: strings.Repeat("f", 64), CorrectionEvidenceSHA256: strings.Repeat("1", 64)}
	step, err = store.RecordRepairQuality(ctx, qualityActor, corrected)
	if err != nil {
		t.Fatal("corrected quality", err)
	}
	accept := wc.AcceptRepair{Command: wc.Command{CaseID: "claim-main", CommandID: "accept", ExpectedVersion: step.Version}, QualitySHA256: step.PayloadSHA256, EvidenceSHA256: strings.Repeat("2", 64)}
	wrong := customer
	wrong.Subject = "different-customer"
	if _, err = store.AcceptRepair(ctx, wrong, accept); err == nil {
		t.Fatal("wrong customer accepted repair")
	}
	step, err = store.AcceptRepair(ctx, customer, accept)
	if err != nil {
		t.Fatal("customer acceptance", err)
	}
	reconcile := wc.ReconcileRepair{Command: wc.Command{CaseID: "claim-main", CommandID: "reconcile", ExpectedVersion: step.Version}, AcceptanceSHA256: step.PayloadSHA256, EvidenceSHA256: strings.Repeat("3", 64)}
	if _, err = store.ReconcileRepair(ctx, operator, reconcile); err == nil {
		t.Fatal("store performed factory acknowledgement")
	}
	step, err = store.ReconcileRepair(ctx, factory, reconcile)
	if err != nil || step.State != "closed" {
		t.Fatal("factory reconciliation", step, err)
	}
	var receipt struct {
		Cost string `json:"recorded_inventory_cost"`
		Paid bool   `json:"payment_created"`
	}
	if json.Unmarshal(step.Payload, &receipt) != nil || receipt.Cost != "740.0000" || receipt.Paid {
		t.Fatal("invented settlement", string(step.Payload))
	}
	again, err = store.ReconcileRepair(ctx, factory, reconcile)
	if err != nil || !again.Replay {
		t.Fatal("reconciliation recovery", again, err)
	}
	// Declined covered service and cancelled approved service both release the
	// original reservations without writing an issue or erasing human decisions.
	for _, scenario := range []string{"declined", "cancelled", "excluded"} {
		fault := "fixture-wear"
		if scenario == "excluded" {
			fault = "fixture-exclusion"
		}
		st, _ := open("claim-"+scenario, fault)
		st, hash := planned("claim-"+scenario, "1", st.Version)
		d := wc.DecideRepair{Command: wc.Command{CaseID: "claim-" + scenario, CommandID: "decision", ExpectedVersion: st.Version}, PayloadSHA256: hash, Approved: scenario == "cancelled", Reason: "Synthetic decision"}
		if scenario == "excluded" {
			bad := d
			bad.Approved = true
			if _, err = store.DecideRepair(ctx, reviewer, bad); err == nil {
				t.Fatal("excluded fault approved")
			}
		}
		st, err = store.DecideRepair(ctx, reviewer, d)
		if err != nil {
			t.Fatal(scenario, err)
		}
		if scenario == "cancelled" {
			st, err = store.CancelRepair(ctx, reviewer, wc.CancelRepair{Command: wc.Command{CaseID: st.CaseID, CommandID: "cancel", ExpectedVersion: st.Version}, Reason: "Customer withdrew before stock issue", EvidenceSHA256: strings.Repeat("4", 64)})
			if err != nil {
				t.Fatal("cancel approved unperformed work", err)
			}
		}
		if st.State != "cancelled" {
			t.Fatal(scenario, st)
		}
	}
	if err = r.pool.QueryRow(ctx, `select quantity::text,reserved_quantity::text,(select count(*) from inventory.bulk_inventory_entry where tenant_id=$1 and entry_type='issue') from inventory.bulk_balance where tenant_id=$1 and item_id='warranty-part'`, r.tenant).Scan(&quantity, &reserved, &issues); err != nil || quantity != "3.000000" || reserved != "0.000000" || issues != 1 {
		t.Fatal("final stock", quantity, reserved, issues, err)
	}
	t.Logf("WARRANTY_CLAIM_J4_PASS real_appointment_request_assignment_confirmation_completion=true appointment_profile_sha256=%s diagnosis_approval_parts_fifo_quality_customer_factory=true claims=4 stock_issues=1 fifo_applications=2 source_cost=740.0000 atomic_rollback=true rejected_and_cancelled_holds_released=true live_proven=false", policy.SHA256())
}

func warrantyExactCost(actual, want string) bool {
	a, ok := new(big.Rat).SetString(actual)
	if !ok {
		return false
	}
	b, ok := new(big.Rat).SetString(want)
	return ok && a.Cmp(b) == 0
}
