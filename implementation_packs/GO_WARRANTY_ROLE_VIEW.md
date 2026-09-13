# Authorized warranty role projections

## 1. Metadata

```yaml
pack_id: "GO-WARRANTY-ROLE-VIEW"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: RECONSTRUCTIBLE
  admission: CONDITIONED
claim: "Sold terms, explicit customer acknowledgement, activation and complete warranty repair workflow through authorized role forms and durable recovery; local reference."
stacks: ["Go 1.26.8", "PostgreSQL 18.6", "Node.js 24.20.0 where frontend selected"]
compatible_with: ["MARKDOWN-COMPOSITOR0.3.0", "V402 existing owner closure"]
incompatible_with: ["unbound tenant/provider/account", "production certification inferred from fixtures", "implicit source or business-policy attribution"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: []
verified_at: "2026-09-12"
```

## 2. Applicability

LIBRARY_INFRASTRUCTURE T2804/J1/J4: existing sold warranty/coverage/repair, approval/inventory, Next/OIDC/BFF/PG owners. Optional feature enabled only with configured warranty profile.

## 3. Architecture contract

Original writers, domain guards and hash verification retained. RoleContext projected only after existing Claim authorization in its repeatable-read transaction. Quote source read requires warranty:offer; monetary/version int64 as strings. No new state machine.

## 4. Exact file manifest

```text
CREATE internal/platform/postgres/warranty_role_browser_integration_test.go
CREATE internal/platform/postgres/warranty_role_view.go
CREATE internal/warrantyclaim/role_golden_test.go
CREATE internal/warrantyclaim/role_view.go
```

## 5. Materialization blocks

### FILE: `internal/platform/postgres/warranty_role_browser_integration_test.go`

```yaml
block_id: "GO-WARRANTY-ROLE-VIEW:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "dcce9f39987377e2e47814c51c693b03f00abb1d3f3a3a68a3044702bda17547"
variables: []
secrets_allowed: false
```

````go
package postgres_test

// AUTHORED fixture: existing writers prepare payment/handover/service/parts;
// terms, consent, activation and every repair step are performed in the UI.
import (
	"context"
	"elite.local/enterprise/internal/businesspolicy"
	"elite.local/enterprise/internal/franchisejourney"
	"elite.local/enterprise/internal/inventorycontrol"
	"elite.local/enterprise/internal/platform/httpapi"
	db "elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/platform/randomid"
	wc "elite.local/enterprise/internal/warrantyclaim"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"net"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

func startWarrantyRoleBrowser(t *testing.T, ctx context.Context, pool *pgxpool.Pool, tenant string, store *db.Warranty) (func(string, map[string]string), func()) {
	t.Helper()
	verifier, token := catalogRoleIssuer(t)
	identities := map[string]any{}

	for _, name := range []string{"operator", "reviewer", "quality", "customer", "factory", "reader", "foreign", "other-customer", "unprivileged"} {
		permissions := []string{"warranty:read"}
		orgs := []string{"store"}
		switch name {
		case "operator":
			permissions = append(permissions, "warranty:offer", "warranty:activate", "warranty:request", "warranty:diagnose", "warranty:plan", "warranty:approve", "warranty:work", "warranty:quality", "warranty:cancel")
		case "reviewer":
			permissions = append(permissions, "warranty:approve")
		case "quality":
			permissions = append(permissions, "warranty:quality")
		case "customer", "other-customer":
			permissions = []string{"warranty:self"}
		case "factory":
			permissions = []string{"warranty:factory-read", "warranty:reconcile"}
			orgs = []string{"factory"}
		case "foreign":
			orgs = []string{"foreign"}
		case "unprivileged":
			permissions = []string{}
		}
		identities[name] = map[string]any{"subject": name, "tenantId": tenant, "permissions": permissions, "organizations": orgs, "accessToken": token(name, tenant, permissions, orgs)}
	}
	mux := http.NewServeMux()
	httpapi.WarrantyModule{Service: store}.Register(mux, verifier)
	var mu sync.Mutex
	counts := map[string]int{}
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			mu.Lock()
			counts[r.URL.Path]++
			mu.Unlock()
		}
		mux.ServeHTTP(w, r)
	}))
	t.Cleanup(api.Close)
	web, node := os.Getenv("ELITE_WEB_ROOT"), os.Getenv("ELITE_NODE_BIN")
	if !filepath.IsAbs(web) || !filepath.IsAbs(node) {
		t.Fatal("fixed runtime/root required")
	}
	listener, e := net.Listen("tcp", "127.0.0.1:0")
	if e != nil {
		t.Fatal(e)
	}
	address := listener.Addr().String()
	listener.Close()
	target, _ := url.Parse("http://" + address)
	edge := httptest.NewTLSServer(httputil.NewSingleHostReverseProxy(target))
	t.Cleanup(edge.Close)
	env := []string{}
	for _, v := range os.Environ() {
		k := strings.ToUpper(strings.SplitN(v, "=", 2)[0])
		if strings.HasPrefix(k, "ELITE_") || strings.HasPrefix(k, "CATALOG_") || strings.HasPrefix(k, "PUBLIC_") || strings.HasPrefix(k, "ENTERPRISE_") || strings.Contains(k, "DATABASE_URL") || k == "APP_BASE_URL" || k == "AUTH_SESSION_SECRET" || k == "BUSINESS_CONFIG_FILE" {
			continue
		}
		env = append(env, v)
	}
	identitiesJSON, _ := json.Marshal(identities)
	env = append(env, "ELITE_WARRANTY_ROLE_BROWSER=1", "ELITE_WORKSPACE_E2E=enabled", "ELITE_WARRANTY_IDENTITIES="+string(identitiesJSON), "ELITE_WEB_ROOT="+web, "ELITE_BASE_URL="+edge.URL,
		"APP_BASE_URL="+edge.URL, "ENTERPRISE_API_BASE_URL="+api.URL, "AUTH_SESSION_SECRET=synthetic-warranty-role-"+uuid.NewString(), "BUSINESS_CONFIG_FILE=warranty.role.fixture.json",
		"CATALOG_RELEASE_ENABLED=false", "PUBLIC_SITE_ORIGIN=https://catalog.example.invalid", "PUBLIC_INDEXING_ENABLED=0", "ENTERPRISE_TENANT_CODE=role-catalog", "ENTERPRISE_ORGANIZATION_CODE=role-org", "NEXT_TELEMETRY_DISABLED=1")
	artifacts, e := os.MkdirTemp(web, "warranty-role-artifacts-")
	if e != nil {
		t.Fatal(e)
	}
	log, e := os.Create(filepath.Join(artifacts, "next.log"))
	if e != nil {
		t.Fatal(e)
	}
	t.Cleanup(func() { log.Close() })
	server := exec.CommandContext(ctx, node, filepath.Join(web, "node_modules", "next", "dist", "bin", "next"), "start", "--hostname", "127.0.0.1", "--port", strings.Split(address, ":")[1])
	server.Dir, server.Env, server.Stdout, server.Stderr = web, env, log, log
	if e = server.Start(); e != nil {
		t.Fatal(e)
	}
	t.Cleanup(func() { server.Process.Kill(); server.Wait() })
	client := &http.Client{Timeout: time.Second}
	ready := false
	for deadline := time.Now().Add(20 * time.Second); time.Now().Before(deadline); {
		r, e := client.Get(target.String() + "/icon.svg")
		if e == nil {
			r.Body.Close()
			if r.StatusCode == 200 {
				ready = true
				break
			}
		}
		select {
		case <-ctx.Done():
			t.Fatal("startup cancelled")
		case <-time.After(100 * time.Millisecond):
		}
	}
	if !ready {
		t.Fatal("Next startup", artifacts)
	}
	gate := filepath.Join(web, "microsoft_playwright_browser_gate")
	run := func(phase string, inputs map[string]string) {
		inputsJSON, _ := json.Marshal(inputs)
		command := exec.CommandContext(ctx, node, filepath.Join(gate, "node_modules", "@playwright", "test", "cli.js"), "test", "tests/warranty-role-connected.spec.mjs", "--project=chromium-desktop", "--timeout=90000", "--output="+filepath.Join(artifacts, "browsers-"+phase))
		command.Dir, command.Env = gate, append(append([]string{}, env...), "ELITE_WARRANTY_PHASE="+phase, "ELITE_WARRANTY_INPUTS="+string(inputsJSON))
		output, e := command.CombinedOutput()
		if x := os.WriteFile(filepath.Join(artifacts, "browser-"+phase+".log"), output, 0600); x != nil {
			t.Fatal(x)
		}
		if e != nil || !strings.Contains(string(output), "1 passed") {
			t.Fatalf("warranty role browser %v\n%s\nartifacts=%s", e, output, artifacts)
		}
	}
	finish := func() {
		mu.Lock()
		defer mu.Unlock()
		total := 0
		for _, n := range counts {
			total += n
		}
		if total != 18 {
			t.Fatal("backend writes", total, counts)
		}
		raw, _ := json.Marshal(counts)
		if e := os.WriteFile(filepath.Join(artifacts, "api-post-counts.json"), raw, 0600); e != nil {
			t.Fatal(e)
		}
		t.Logf("WARRANTY_ROLE_BROWSER_PASS phases=terms,claims backend_POSTs=18 response_loss_GET_only=true artifacts=%s", artifacts)
	}
	return run, finish
}

func prepareWarrantyRoleInputs(t *testing.T, r *connectedRun) []string {
	ctx := context.Background()
	ids := randomid.Generator{}
	must := func(q string, args ...any) {
		t.Helper()
		if _, err := r.pool.Exec(ctx, q, args...); err != nil {
			t.Fatal(err)
		}
	}
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

	return []string{appointment("role-main"), appointment("role-cancel"), appointment("role-excluded")}
}
func TestWarrantyRoleBrowser(t *testing.T) {
	if os.Getenv("ELITE_WARRANTY_ROLE_BROWSER") != "1" {
		t.Skip("explicit fixture required")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	var run func(string, map[string]string)
	var finish func()
	r := newConnectedRunWithAllocation(t, false, nil, func(t *testing.T, pool *pgxpool.Pool, tenant string) int64 {
		if _, e := pool.Exec(ctx, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'factory','factory','Synthetic factory','factory')`, tenant); e != nil {
			t.Fatal(e)
		}
		store, operator, _ := fixtureWarranty(t, pool, tenant, 365)
		run, finish = startWarrantyRoleBrowser(t, ctx, pool, tenant, store)
		run("terms", map[string]string{"quote": "quote"})
		offer, e := store.Offer(ctx, operator, "quote")
		if e != nil || !offer.Acknowledged || offer.QuoteVersion != 2 {
			t.Fatal("UI terms", offer, e)
		}
		return offer.QuoteVersion
	})
	if code := r.callback(t, "evt_warranty_role", "checkout.session.completed", "cs_test_fixture", true); code != 200 {
		t.Fatal(code)
	}
	if n, e := r.processor.ProcessOnce(ctx); e != nil || n != 1 {
		t.Fatal(n, e)
	}
	var hash string
	if e := r.pool.QueryRow(ctx, `select evidence_sha256_hex from payment.provider_observation where tenant_id=$1`, r.tenant).Scan(&hash); e != nil {
		t.Fatal(e)
	}
	policy, a := connectedCommercialProfile(t, r, false)
	repo := db.NewFranchiseJourney(r.pool)
	ids := randomid.Generator{}
	svc, err := franchisejourney.NewHandoverPreparationService(repo, ids, policy)
	if err != nil {
		t.Fatal(err)
	}
	prepared, _, err := svc.Prepare(ctx, r.tenant, "operator", franchisejourney.PrepareHandoverCommand{OrganizationID: "store", OrderID: "order", OrderLineID: "line", PaymentAttemptID: r.payment.ID, ObservationSHA256: hash, IdempotencyKey: "commercial-integral-prepare"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = repo.PublishDeliveryChecklist(ctx, r.tenant, "operator", franchisejourney.DeliveryChecklist{ID: "commercial-checklist", OrganizationID: "store", Version: 1, Title: "Synthetic delivery", Items: []franchisejourney.ChecklistItem{{ID: "serial", Ordinal: 1, Prompt: "Read serial", ResponseType: "serial", Required: true}}}, ids.New()); err != nil {
		t.Fatal(err)
	}
	presented, err := repo.CompleteDeliveryChecklist(ctx, r.tenant, "store", "operator", prepared.Handover.ID, 1, "commercial-checklist", 1, []franchisejourney.ChecklistResponse{{ItemID: "serial", ResponseText: "SERIAL-SYNTHETIC"}}, ids.New())
	if err != nil {
		t.Fatal(err)
	}
	if _, err = repo.AcceptHandover(ctx, r.tenant, "store", "customer", prepared.Handover.ID, presented.Version, "SERIAL-SYNTHETIC", "commercial-checklist", 1, strings.Repeat("e", 64), ids.New()); err != nil {
		t.Fatal(err)
	}
	release := franchisejourney.CommitCommercialReleaseCommand{OrganizationID: "store", HandoverID: prepared.Handover.ID, ObservationSHA256: hash, IdempotencyKey: "commercial-integral-release"}
	receipt, replay, err := svc.CommitCommercialRelease(ctx, r.tenant, "operator", release)
	if err != nil || replay || receipt.ContractSHA256 != a.DocumentSHA256 {
		t.Fatal("commercial checkpoint", err)
	}
	current, err := svc.ValidateCommercialRelease(ctx, r.tenant, "store", prepared.Handover.ID)
	if err != nil || !current.Current {
		t.Fatal("checkpoint not current", err)
	}

	appointments := prepareWarrantyRoleInputs(t, r)
	run("claims", map[string]string{"handover": prepared.Handover.ID, "main": appointments[0], "cancel": appointments[1], "excluded": appointments[2]})
	var actual string
	err = r.pool.QueryRow(ctx, `select jsonb_build_array(
 (select count(*)from service_ops.service_case where tenant_id=$1 and state='closed'),
 (select count(*)from service_ops.service_case where tenant_id=$1 and state='cancelled'),
 (select count(*)from service_ops.warranty_claim_step where tenant_id=$1),
 (select quantity::text from inventory.bulk_balance where tenant_id=$1 and item_id='warranty-part'),
 (select sum(cost_amount)::text from inventory.bulk_inventory_entry where tenant_id=$1 and entry_type='issue'))::text`, r.tenant).Scan(&actual)
	if err != nil || actual != `[1, 2, 15, "3.000000", "-740.0000"]` {
		t.Fatal("durable effects", actual, err)
	}
	finish()
	t.Log("WARRANTY_ROLE_DURABLE_EFFECTS", actual)
}
````

### FILE: `internal/platform/postgres/warranty_role_view.go`

```yaml
block_id: "GO-WARRANTY-ROLE-VIEW:file2:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "f8dc71ac98e87b759d445300e69363c74034cd81944dd5adb1fb7cf1f8b88ce5"
variables: []
secrets_allowed: false
```

````go
// AUTHORED bounded projection of already authorized immutable evidence, read in
// the Claim repeatable-read transaction. No write or new domain algorithm.
package postgres

import (
	"context"
	"elite.local/enterprise/internal/platform/identity"
	wc "elite.local/enterprise/internal/warrantyclaim"
	"encoding/json"
	"errors"
)

func readWarrantyRoleContext(ctx context.Context, q warrantyRowReader, tenant, id string) (*wc.RoleContext, error) {
	out := &wc.RoleContext{}
	opened, e := readWarrantyStep(ctx, q, tenant, id, "", "opened")
	if e != nil {
		return nil, e
	}
	var opening struct {
		Request wc.OpenClaim `json:"request"`
	}
	if json.Unmarshal(opened.Payload, &opening) != nil || opening.Request.CaseID != id {
		return nil, wc.ErrConflict
	}
	out.Description = opening.Request.Description
	out.Severity = opening.Request.Severity
	for _, kind := range []string{"diagnosed", "work", "quality", "accepted", "reconciled"} {
		step, e := readWarrantyStep(ctx, q, tenant, id, "", kind)
		if errors.Is(e, wc.ErrNotFound) {
			continue
		}
		if e != nil {
			return nil, e
		}
		switch kind {
		case "diagnosed":
			var v struct {
				Diagnosis wc.Diagnose `json:"diagnosis"`
				Excluded  bool        `json:"fault_excluded"`
			}
			if json.Unmarshal(step.Payload, &v) != nil || v.Diagnosis.CaseID != id {
				return nil, wc.ErrConflict
			}
			out.Diagnosis = &wc.RoleDiagnosis{FaultCode: v.Diagnosis.FaultCode, Description: v.Diagnosis.Description, Excluded: v.Excluded, EvidenceSHA256: v.Diagnosis.EvidenceSHA256}
		case "work":
			var v wc.WorkReceipt
			if json.Unmarshal(step.Payload, &v) != nil || v.Request.CaseID != id {
				return nil, wc.ErrConflict
			}
			out.Work = &wc.RoleWork{Actor: step.Actor, PayloadSHA256: step.PayloadSHA256, EvidenceSHA256: v.Request.EvidenceSHA256}
		case "quality":
			var v wc.Quality
			if json.Unmarshal(step.Payload, &v) != nil || v.CaseID != id {
				return nil, wc.ErrConflict
			}
			out.Quality = &wc.RoleQuality{Passed: v.Passed, PayloadSHA256: step.PayloadSHA256, EvidenceSHA256: v.EvidenceSHA256}
		case "accepted":
			var v wc.AcceptRepair
			if json.Unmarshal(step.Payload, &v) != nil || v.CaseID != id {
				return nil, wc.ErrConflict
			}
			out.Acceptance = &wc.RoleAcceptance{Actor: step.Actor, PayloadSHA256: step.PayloadSHA256}
		case "reconciled":
			var v wc.RoleReconciliation
			if json.Unmarshal(step.Payload, &v) != nil || v.Settlement != "INTERNAL_WARRANTY_SERVICE_ACKNOWLEDGEMENT" || v.RecordedInventoryCost == "" {
				return nil, wc.ErrConflict
			}
			out.Reconciliation = &v
		}
	}
	plan, approvalID, hash, e := readWarrantyPlan(ctx, q, tenant, id)
	if e == nil {
		var state string
		if e = q.QueryRow(ctx, `select state from approval.request where tenant_id=$1 and request_id=$2 and evidence_sha=$3`, tenant, approvalID, hash).Scan(&state); e != nil {
			return nil, e
		}
		out.Plan = &wc.RolePlan{Request: plan.Request, Requester: plan.Requester, PayloadSHA256: hash, Excluded: plan.Excluded, ExpiresAt: plan.ExpiresAt, ApprovalState: state}
	} else if !errors.Is(e, wc.ErrNotFound) {
		return nil, e
	}
	return out, nil
}

func (s *Warranty) QuoteForOffer(ctx context.Context, p identity.Principal, id string) (wc.QuoteForOffer, error) {
	var out wc.QuoteForOffer
	if !s.allowed(p, "warranty:offer") || !wc.ValidID(id) {
		return out, wc.ErrNotFound
	}
	tenant, org := s.profile.Scope()
	e := s.pool.QueryRow(ctx, `select quotation_id,organization_id,version,coalesce(customer_principal_id,''),state,valid_until>clock_timestamp(),currency,total_minor_units from sales.quotation where tenant_id=$1 and organization_id=$2 and quotation_id=$3`, tenant, org, id).Scan(&out.QuoteID, &out.OrganizationID, &out.QuoteVersion, &out.CustomerSubject, &out.State, &out.Current, &out.Currency, &out.TotalMinorUnits)
	return out, warrantyError(e)
}
````

### FILE: `internal/warrantyclaim/role_golden_test.go`

```yaml
block_id: "GO-WARRANTY-ROLE-VIEW:file3:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "acf6e0e2ba5d49cba39f6da78e3991c0260c281df627a00af5fbb817ed5b26a0"
variables: []
secrets_allowed: false
```

````go
package warrantyclaim

import (
	"encoding/json"
	"os"
	"testing"
)

func TestWarrantyRoleWireGoldens(t *testing.T) {
	raw, e := os.ReadFile("../../deploy/warranty/role-goldens.json")
	if e != nil {
		t.Fatal(e)
	}
	var cases []struct {
		Name      string
		Payload   json.RawMessage
		Canonical string
		SHA256    string
	}
	if e = json.Unmarshal(raw, &cases); e != nil {
		t.Fatal(e)
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			var v any
			switch c.Name {
			case "offer":
				v = new(OfferRequest)
			case "acknowledge":
				v = new(AcknowledgeRequest)
			case "activate":
				v = new(struct {
					HandoverID string `json:"handover_id"`
				})
			case "open":
				v = new(OpenClaim)
			case "diagnose":
				v = new(Diagnose)
			case "plan":
				v = new(Plan)
			case "decide":
				v = new(DecideRepair)
			case "work":
				v = new(CompleteWork)
			case "quality":
				v = new(Quality)
			case "accept":
				v = new(AcceptRepair)
			case "reconcile":
				v = new(ReconcileRepair)
			case "cancel":
				v = new(CancelRepair)
			default:
				t.Fatal("unknown")
			}
			if e = json.Unmarshal(c.Payload, v); e != nil {
				t.Fatal(e)
			}
			raw, h, e := Canonical(v)
			if e != nil || string(raw) != c.Canonical || h != c.SHA256 {
				t.Fatal(string(raw), h, e)
			}
		})
	}
}
````

### FILE: `internal/warrantyclaim/role_view.go`

```yaml
block_id: "GO-WARRANTY-ROLE-VIEW:file4:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "470a5720b0c445f6ed96ec9dfa97e16e205692f8565c9e96b9fc182ac88ead1f"
variables: []
secrets_allowed: false
```

````go
// AUTHORED read projection: all decisions remain in the existing owners.
package warrantyclaim

import "time"

type RoleContext struct {
	Description    string              `json:"description"`
	Severity       string              `json:"severity"`
	Diagnosis      *RoleDiagnosis      `json:"diagnosis,omitempty"`
	Plan           *RolePlan           `json:"plan,omitempty"`
	Work           *RoleWork           `json:"work,omitempty"`
	Quality        *RoleQuality        `json:"quality,omitempty"`
	Acceptance     *RoleAcceptance     `json:"acceptance,omitempty"`
	Reconciliation *RoleReconciliation `json:"reconciliation,omitempty"`
}
type RoleDiagnosis struct {
	FaultCode      string `json:"fault_code"`
	Description    string `json:"description"`
	Excluded       bool   `json:"excluded"`
	EvidenceSHA256 string `json:"evidence_sha256"`
}
type RolePlan struct {
	Request       Plan      `json:"request"`
	Requester     string    `json:"requester"`
	PayloadSHA256 string    `json:"payload_sha256"`
	Excluded      bool      `json:"excluded"`
	ExpiresAt     time.Time `json:"expires_at"`
	ApprovalState string    `json:"approval_state"`
}
type RoleWork struct {
	Actor          string `json:"actor"`
	PayloadSHA256  string `json:"payload_sha256"`
	EvidenceSHA256 string `json:"evidence_sha256"`
}
type RoleQuality struct {
	Passed         bool   `json:"passed"`
	PayloadSHA256  string `json:"payload_sha256"`
	EvidenceSHA256 string `json:"evidence_sha256"`
}
type RoleAcceptance struct {
	Actor         string `json:"actor"`
	PayloadSHA256 string `json:"payload_sha256"`
}
type RoleReconciliation struct {
	Settlement            string `json:"settlement"`
	RecordedInventoryCost string `json:"recorded_inventory_cost"`
}

// QuoteForOffer avoids operator-supplied version numbers; BindOffer still locks
// and validates the real quote when a write is requested.
type QuoteForOffer struct {
	QuoteID         string `json:"quote_id"`
	OrganizationID  string `json:"organization_id"`
	QuoteVersion    int64  `json:"quote_version,string"`
	CustomerSubject string `json:"customer_subject"`
	State           string `json:"state"`
	Current         bool   `json:"current"`
	Currency        string `json:"currency"`
	TotalMinorUnits int64  `json:"total_minor_units,string"`
}
````

## 6. Configuration surface

docs/WARRANTY_ROLE_REFERENCE.md. features.warranty_portal opt-in. Role read and explicit action permissions. Evidence hash and read-only predecessor context; up to16parts and32KiB BFF. Pending result reference scoped to tenant/org/actor/view, GET only recovery.

## 7. Dependency bill

14new and7changed AUTHORED unavoidable projection/transport/form/config/fixture glue. No new dependency/runtime; retained BC date adaptation and other official source licenses, no corporate attribution.

## 8. Apply order

Select GO-CONNECTED-WARRANTY-CLAIM0.2.0 with this role-view companion and TS-WARRANTY-ROLE-PORTAL. Browser fixture uses existing GO-CONNECTED-CATALOG-AUTHORING issuer and payment/handover fixture closure. Domain golden fixture is owned by the TS companion. Full reference required for connected tests.

## 9. Verification

18web tests,12typed Go/TS goldens, full Next webpack/types. Both actual browser phases pass:18POST,15steps,3claims,1closed/2cancelled,3remaining parts, signed ledger -740.0000 and positive service cost740.0000. Two response losses recovered GET after reload with0extraPOST; review/quality separation, exclusion/rejection/cancellation and foreign/other-customer checks. Screenshot desktop/390px inspected.

## 10. Reconstruction evidence

WARRANTY_ROLE_RELEASE_V402.md/json records exact source/rebuilds and original RED wrapper FAIL859. Only final ledger-sign fixture assertion corrected; existing owned DB verified read-only without replay. Original domain fuzz/host/transaction guards retained, no new algorithm/migration. Does not close wholeT2804 or production.

