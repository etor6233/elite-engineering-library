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
