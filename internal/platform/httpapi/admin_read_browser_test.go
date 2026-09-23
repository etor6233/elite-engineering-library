package httpapi

import (
	"context"
	"elite.local/enterprise/internal/enterprisequery"
	"elite.local/enterprise/internal/franchisejourney"
	"elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/platform/randomid"
	"encoding/json"
	"fmt"
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
	"sync/atomic"
	"testing"
	"time"
)

// AUTHORED read-only browser fixture. Requires the selected journey test issuer;
// it exercises actual RS256/JWKS, HTTP scopes, PostgreSQL queries and Next JWE.
func TestAdministrativeReadBrowserPostgres(t *testing.T) { runPrivateBrowserPostgres(t, false) }

func TestCustomerCancellationBrowserPostgres(t *testing.T) {
 if os.Getenv("ELITE_CUSTOMER_CANCEL_E2E") != "1" { t.Skip("explicit disposable customer cancellation gate not requested") }
 runPrivateBrowserPostgres(t, true)
}

func runPrivateBrowserPostgres(t *testing.T, cancellation bool) {
	if os.Getenv("ELITE_ADMIN_READ_E2E") != "1" {
		t.Skip("explicit disposable admin-read gate not requested")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	cfg, err := pgxpool.ParseConfig(os.Getenv("TEST_DATABASE_URL"))
	if err != nil || cfg.ConnConfig.Host != "127.0.0.1" || !strings.HasPrefix(cfg.ConnConfig.Database, "elite_confirmation_") {
		t.Fatal("requires disposable loopback elite_confirmation_* database")
	}
	web := os.Getenv("ELITE_WEB_ROOT")
	if !filepath.IsAbs(web) {
		t.Fatal("absolute composed web root required")
	}
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := randomid.Generator{}.New()
	fixtures := []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,'query-api','Query','Query')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,'org-a','org-a','A','store'),($1,'org-b','org-b','B','store')`,
		`insert into crm.customer_profile(tenant_id,customer_principal_id,display_name,email_normalized,status) values($1,'customer-1','Customer 1','customer1@example.test','active'),($1,'customer-2','Customer 2','customer2@example.test','active')`,
		`insert into catalog.vehicle_model(tenant_id,model_id,model_code,display_name,vehicle_class,lifecycle_state) values($1,'model','model','Model','bicycle','active')`,
		`insert into catalog.vehicle_variant(tenant_id,variant_id,model_id,variant_code,display_name,battery_specification,lifecycle_state) values($1,'variant','model','variant','Variant','{}','active')`,
		`insert into partner.supplier(tenant_id,supplier_id,supplier_code,legal_name,status) values($1,'supplier','supplier','Supplier','active')`,
		`insert into procurement.purchase_order(tenant_id,purchase_order_id,supplier_id,destination_organization_id,state,currency,total_minor_units,version) values($1,'po','supplier','org-a','accepted','USD',100,1)`,
		`insert into factory.production_unit(tenant_id,production_unit_id,purchase_order_id,variant_id,serial_number,state) values($1,'unit','po','variant','UNIT-QUERY','planned')`,
		`insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version) values($1,'order-a','org-a','customer-1','draft','USD',100,1),($1,'order-b','org-a','customer-2','draft','USD',200,1),($1,'order-c','org-b','customer-1','draft','USD',300,1)`,
		`insert into inventory.stock_unit(tenant_id,stock_unit_id,organization_id,variant_id,serial_number,vin,battery_serial_number,state,version,received_at) values($1,'stock-a','org-a','variant','STOCK-QUERY-A','VIN-QUERY-A','BATTERY-QUERY-A','available',1,clock_timestamp()),($1,'stock-b','org-b','variant','STOCK-QUERY-B','VIN-QUERY-B','BATTERY-QUERY-B','available',1,clock_timestamp())`,
		`insert into service_ops.warranty(tenant_id,warranty_id,stock_unit_id,customer_principal_id,starts_at,ends_at,terms_version,status) values($1,'warranty','stock-a','customer-1',clock_timestamp(),clock_timestamp()+interval '1 year','v1','active')`,
		`insert into service_ops.service_case(tenant_id,service_case_id,stock_unit_id,organization_id,state,severity,description,version) values($1,'case-a','stock-a','org-a','opened','medium','inspection',1)`,
		`insert into crm.lead(tenant_id,lead_id,organization_id,customer_principal_id,lifecycle_state,source_code,contact_payload) values($1,'lead-visible','org-a','customer-1','new','fixture','{}'),($1,'lead-foreign-org','org-b','customer-2','new','fixture','{}')`,
	}
	fixtures = append(fixtures,
		`insert into pricing.price_book(tenant_id,price_book_id,market,currency,valid_from,status) values($1,'book','AR','USD',clock_timestamp()-interval '1 day','active')`,
		`insert into pricing.price_book_entry(tenant_id,price_book_id,variant_id,amount_minor_units,tax_mode) values($1,'book','variant',100,'inclusive')`,
		`insert into sales.quotation(tenant_id,quotation_id,organization_id,lead_id,customer_principal_id,variant_id,price_book_id,currency,total_minor_units,valid_until,state,version) values($1,'quote-visible','org-a','lead-visible','customer-1','variant','book','USD',100,clock_timestamp()+interval '1 day','issued',1)`,
	)

	fixtures = append(fixtures,
		`insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version) select $1,'zz-order-'||lpad(i::text,3,'0'),'org-a','customer-1','draft','USD',100,1 from generate_series(1,26) i`,
		`insert into service_ops.service_case(tenant_id,service_case_id,stock_unit_id,organization_id,state,severity,description,version) select $1,'zz-case-'||lpad(i::text,3,'0'),'stock-a','org-a','opened','medium','pagination fixture',1 from generate_series(1,26) i`,
		`insert into factory.production_unit(tenant_id,production_unit_id,purchase_order_id,variant_id,serial_number,vin,battery_serial_number,state) select $1,'zz-unit-'||lpad(i::text,3,'0'),'po','variant','PAGE-UNIT-'||lpad(i::text,3,'0'),'PAGE-VIN-'||i::text,'PAGE-BATTERY-'||i::text,'planned' from generate_series(1,26) i`,
	)

	fixtures=append(fixtures,`insert into crm.lead(tenant_id,lead_id,organization_id,customer_principal_id,lifecycle_state,source_code,contact_payload) select $1,'zz-lead-'||lpad(i::text,3,'0'),'org-a','customer-1','new','fixture','{}' from generate_series(1,26) i`)

	fixtures = append(fixtures,
        `insert into crm.appointment(tenant_id,appointment_id,organization_id,lead_id,customer_principal_id,appointment_kind,starts_at,state,version) values($1,'time-winter','org-a','lead-visible','customer-1','consultation','2035-01-02T01:30:00Z','requested',1),($1,'time-summer','org-a','lead-visible','customer-1','service','2035-07-02T01:30:00Z','requested',1),($1,'time-foreign','org-b','lead-foreign-org','customer-2','delivery','2035-01-03T01:30:00Z','requested',1)`,
    )

	for _, q := range fixtures {
		if _, err = pool.Exec(ctx, q, tenant); err != nil {
			t.Fatal(err)
		}
	}

 cancelFixtures := map[string]map[string]map[string]string{}
 if cancellation {
  for projectIndex, project := range []string{"chromium-desktop","chromium-mobile","firefox-desktop","webkit-desktop"} {
   cancelFixtures[project]=map[string]map[string]string{}
   for phaseIndex, phase := range []string{"normal","lost","malformed","stale"} {
    id:="cc-"+project+"-"+phase
    at:=time.Date(2035,1,10+projectIndex*4+phaseIndex,1,30,0,0,time.UTC).Format(time.RFC3339)
    cancelFixtures[project][phase]=map[string]string{"id":id,"instant":at}
    if _,err:=pool.Exec(ctx,`insert into crm.appointment(tenant_id,appointment_id,organization_id,lead_id,customer_principal_id,appointment_kind,starts_at,state,version)values($1,$2,'org-a','lead-visible','customer-1','consultation',$3,'requested',1)`,tenant,id,at);err!=nil{t.Fatal(err)}
   }
  }
 }

	snapshot := func() string {
		var digest string
		q := `select md5(string_agg(payload,'|' order by payload)) from (
   select row_to_json(x)::text payload from sales.customer_order x where tenant_id=$1
   union all select row_to_json(x)::text from service_ops.service_case x where tenant_id=$1
   union all select row_to_json(x)::text from factory.production_unit x where tenant_id=$1
   union all select row_to_json(x)::text from crm.lead x where tenant_id=$1
   union all select row_to_json(x)::text from inventory.stock_unit x where tenant_id=$1
   union all select row_to_json(x)::text from sales.quotation x where tenant_id=$1
   union all select row_to_json(x)::text from crm.appointment x where tenant_id=$1
  )q`
        if cancellation {q=strings.Replace(q,"select row_to_json(x)::text from crm.appointment x","select (to_jsonb(x)-'state'-'version'-'updated_at')::text from crm.appointment x",1)}
		if err := pool.QueryRow(ctx, q, tenant).Scan(&digest); err != nil {
			t.Fatal(err)
		}
		return digest
	}
	before := snapshot()
	verifier, token := confirmationTestIssuer(t)
	identities := map[string]any{}
	tokens := map[string]string{}
	add := func(name string, permissions, organizations []string, tenantID string) {
		subject := name
		if name == "customer" {
			subject = "customer-1"
		}
		bearer := token(subject, tenantID, permissions, organizations)
		tokens[name] = bearer
		identities[name] = map[string]any{"subject": subject, "tenantId": tenantID, "permissions": permissions, "organizations": organizations, "accessToken": bearer}
	}
	for name, permissions := range map[string][]string{"admin": {"admin:read"}, "admin-lead": {"admin:read", "lead:read"}, "admin-case": {"admin:read", "Lead:Read"}, "wildcard": {"*"}, "lead": {"lead:read"}, "factory": {"factory:read"}, "forged-admin": {"admin:read"}} {
		add(name, permissions, []string{"org-a"}, tenant)
	}
	add("customer", []string{"customer:self"}, []string{"org-a"}, tenant)
	add("other-org", []string{"admin:read", "lead:read"}, []string{"org-b"}, tenant)
	add("other-tenant", []string{"admin:read", "lead:read"}, []string{"org-a"}, randomid.Generator{}.New())
	mux := http.NewServeMux()
	EnterpriseQueryModule{Service: enterprisequery.NewService(postgres.NewEnterpriseQuery(pool))}.Register(mux, verifier)
	FranchiseJourneyModule{Service: franchisejourney.NewService(postgres.NewFranchiseJourney(pool), randomid.Generator{}, browserAppointmentClock{})}.Register(mux, verifier)
	var writes, adminLeadReads, leadAdminReads atomic.Int64
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
            writes.Add(1)
            if !cancellation || r.Method!=http.MethodPost || !strings.HasPrefix(r.URL.Path,"/v1/customer/appointments/") || !strings.HasSuffix(r.URL.Path,"/cancel") {w.WriteHeader(405);return}
        }
		if r.Header.Get("Authorization") == "Bearer "+tokens["admin"] && r.URL.Path == "/v1/franchise/leads" {
			adminLeadReads.Add(1)
		}
		if r.Header.Get("Authorization") == "Bearer "+tokens["lead"] && strings.HasPrefix(r.URL.Path, "/v1/admin/") {
			leadAdminReads.Add(1)
		}
		mux.ServeHTTP(w, r)
	}))
	defer api.Close()
	// Independent real HTTP negatives before browser counters are recorded.
	for _, probe := range []struct {
		path, bearer string
		status       int
	}{
		{"/v1/admin/overview?organization_id=org-a", "invalid", 401},
		{"/v1/admin/overview?organization_id=org-a", tokens["lead"], 403},
		{"/v1/franchise/leads?organization_id=org-a", tokens["forged-admin"], 403},
		{"/v1/admin/overview?organization_id=org-b", tokens["admin"], 403},
		{"/v1/factory/units?organization_id=org-a", tokens["admin"], 403},
	} {
		req, _ := http.NewRequestWithContext(ctx, "GET", api.URL+probe.path, nil)
		req.Header.Set("Authorization", "Bearer "+probe.bearer)
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		res.Body.Close()
		if res.StatusCode != probe.status {
			t.Fatalf("probe %s status=%d want=%d", probe.path, res.StatusCode, probe.status)
		}
	}
	baselineLeadAdmin := leadAdminReads.Load()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	address := listener.Addr().String()
	listener.Close()
	target, _ := url.Parse("http://" + address)
	edge := httptest.NewTLSServer(httputil.NewSingleHostReverseProxy(target))
	defer edge.Close()
	env := []string{}
	for _, item := range os.Environ() {
		name := strings.ToUpper(strings.SplitN(item, "=", 2)[0])
		if name == "TEST_DATABASE_URL" || name == "DATABASE_URL" || strings.HasPrefix(name, "ELITE_") || name == "ENTERPRISE_API_BASE_URL" || name == "AUTH_SESSION_SECRET" || name == "APP_BASE_URL" {
			continue
		}
		env = append(env, item)
	}
	encoded, err := json.Marshal(identities)
	if err != nil {
		t.Fatal(err)
	}
	env = append(env, "ELITE_APPOINTMENT_TIME_EXPECTATIONS="+os.Getenv("ELITE_APPOINTMENT_TIME_EXPECTATIONS"), "ELITE_ADMIN_READ_E2E=1", "ELITE_WORKSPACE_E2E=enabled", "ELITE_ADMIN_IDENTITIES="+string(encoded), "ELITE_WEB_ROOT="+web, "ENTERPRISE_API_BASE_URL="+api.URL, "AUTH_SESSION_SECRET=synthetic-"+randomid.Generator{}.New()+randomid.Generator{}.New(), "APP_BASE_URL="+edge.URL, "ELITE_BASE_URL="+edge.URL)
    if cancellation {fixtureJSON,err:=json.Marshal(cancelFixtures);if err!=nil{t.Fatal(err)};env=append(env,"ELITE_CUSTOMER_CANCEL_E2E=1","ELITE_CUSTOMER_CANCEL_FIXTURES="+string(fixtureJSON))}

	artifacts, err := os.MkdirTemp(web, "admin-read-artifacts-")
	if err != nil {
		t.Fatal(err)
	}
	nextlog, err := os.Create(filepath.Join(artifacts, "next.log"))
	if err != nil {
		t.Fatal(err)
	}
	defer nextlog.Close()
	server := exec.CommandContext(ctx, "node", filepath.Join(web, "node_modules", "next", "dist", "bin", "next"), "start", "--hostname", "127.0.0.1", "--port", strings.Split(address, ":")[1])
	server.Dir, server.Env = web, env
	server.Stdout, server.Stderr = nextlog, nextlog
	if err = server.Start(); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = server.Process.Kill(); _ = server.Wait() }()
	ready := false
	client := &http.Client{Timeout: time.Second}
	for deadline := time.Now().Add(20 * time.Second); time.Now().Before(deadline); {
		res, err := client.Get(target.String() + "/icon.svg")
		if err == nil {
			res.Body.Close()
			if res.StatusCode == 200 {
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
		t.Fatal("web did not start")
	}
	gate := filepath.Join(web, "microsoft_playwright_browser_gate")
    selectedTest:="admin reads preserve independent grants";if cancellation {selectedTest="customer cancels with durable result recovery"}
	cmd := exec.CommandContext(ctx, "node", filepath.Join(gate, "node_modules", "@playwright", "test", "cli.js"), "test", "tests/role-workspace.spec.mjs", "-g", selectedTest, "--timeout=60000", "--output="+filepath.Join(artifacts, "browsers"))
	cmd.Dir, cmd.Env = gate, env
	output, err := cmd.CombinedOutput()
	if e := os.WriteFile(filepath.Join(artifacts, "browser.log"), output, 0600); e != nil {
		t.Fatal(e)
	}
	if err != nil {
		t.Fatalf("browser: %v\n%s", err, output)
	}
	if !strings.Contains(string(output), "4 passed") {
		t.Fatalf("missing four browser successes: %s", output)
	}

 if cancellation {
  var appointments,transitions,events,unchanged int
  checks:=[]struct{q string;dst *int;want int}{
   {`select count(*) from crm.appointment where tenant_id=$1 and appointment_id like 'cc-%' and state='cancelled' and version=2`,&appointments,16},
   {`select count(*) from crm.appointment_transition where tenant_id=$1 and appointment_id like 'cc-%' and from_state='requested' and to_state='cancelled' and actor_subject='customer-1' and reason_code='customer-request'`,&transitions,16},
   {`select count(*) from platform.outbox_event where tenant_id=$1 and aggregate_type='appointment' and aggregate_id like 'cc-%' and aggregate_version=2 and event_type='appointment.cancelled' and payload->>'actor_subject'='customer-1' and payload->>'channel'='customer' and payload->>'reason_code'='customer-request'`,&events,16},
   {`select count(*) from crm.appointment where tenant_id=$1 and appointment_id not like 'cc-%' and state='requested' and version=1`,&unchanged,3},
  }
  for _,check:=range checks {if err:=pool.QueryRow(ctx,check.q,tenant).Scan(check.dst);err!=nil || *check.dst!=check.want {t.Fatalf("cancellation durable check got=%d want=%d err=%v",*check.dst,check.want,err)}}
  var allTransitions,allEvents int
  if err:=pool.QueryRow(ctx,`select (select count(*) from crm.appointment_transition where tenant_id=$1),(select count(*) from platform.outbox_event where tenant_id=$1)`,tenant).Scan(&allTransitions,&allEvents);err!=nil || allTransitions!=16 || allEvents!=16 {t.Fatalf("unexpected durable effects transitions=%d events=%d err=%v",allTransitions,allEvents,err)}
  if writes.Load()!=28 || snapshot()!=before {t.Fatalf("cancellation writes=%d or unrelated/immutable data changed",writes.Load())}
  t.Logf("CUSTOMER_CANCELLATION_BROWSER_PASS browsers=4 cancelled=16 transitions=16 outbox=16 writes=28 immutable_snapshot=%s artifacts=%s",before,artifacts)
  return
 }

	if writes.Load() != 0 || adminLeadReads.Load() != 0 || leadAdminReads.Load() != baselineLeadAdmin {
		t.Fatalf("unexpected requests writes=%d admin_leads=%d lead_admin=%d", writes.Load(), adminLeadReads.Load(), leadAdminReads.Load()-baselineLeadAdmin)
	}
	if snapshot() != before {
		t.Fatal("read-only browser modified durable records")
	}
	t.Logf("ADMIN_READ_BROWSER_POSTGRES_PASS browsers=4 http_negatives=5 admin_only_lead_queries=0 non_admin_queries=0 writes=0 paging_records=81 admin_paging_records=82 appointment_times=2 snapshot_tables=7 durable_snapshot=%s artifacts=%s", before, artifacts)
	fmt.Fprintln(os.Stdout, "ADMIN_READ_OWNED_NEXT_STOP_ON_RETURN")
}
