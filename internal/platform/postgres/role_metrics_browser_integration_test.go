package postgres_test

// AUTHORED browser composition of scoped read models and existing survey writer.
import (
	"context"
	"elite.local/enterprise/internal/customerfeedback"
	"elite.local/enterprise/internal/enterprisequery"
	"elite.local/enterprise/internal/platform/httpapi"
	db "elite.local/enterprise/internal/platform/postgres"
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

func TestRoleMetricsBrowser(t *testing.T) {
	if os.Getenv("ELITE_ROLE_METRICS_BROWSER") != "1" {
		t.Skip("explicit local browser fixture required")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	pool, e := pgxpool.New(ctx, os.Getenv("PAYMENT_CONNECTED_DB_URL"))
	if e != nil {
		t.Fatal(e)
	}
	defer pool.Close()
	tenant := roleMetricSeed(t, pool)
	must := func(q string, args ...any) {
		t.Helper()
		if _, e := pool.Exec(ctx, q, args...); e != nil {
			t.Fatal(e)
		}
	}
	must(`insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version)values($1,'large','store','customer','confirmed','ARS',9007199254740993,1)`, tenant)
	stored, p, program := fixtureStoredValueSetup(t, pool, tenant, "gift_card")
	_ = p
	_ = program
	must(`insert into crm.survey_definition(tenant_id,organization_id,survey_id,prompt,consent_version,consent_notice,opens_at,closes_at,retain_until,minimum_responses,active)values($1,'store','survey','Fixture','v1','Fixture',clock_timestamp()-interval '1 hour',clock_timestamp()+interval '1 hour',clock_timestamp()+interval '1 day',2,true)`, tenant)
	metrics, e := db.NewRoleMetrics(pool)
	if e != nil {
		t.Fatal(e)
	}
	verifier, token := catalogRoleIssuer(t)
	identities := map[string]any{}
	for _, name := range []string{"admin", "customer", "customer-2", "none", "foreign"} {
		perms := []string{"customer:self"}
		orgs := []string{"store"}
		if name == "admin" {
			perms = []string{"admin:read", "surveys:read", "stored_value:read"}
			orgs = []string{"store", "other"}
		}
		if name == "none" {
			perms = []string{}
		}
		if name == "foreign" {
			perms = []string{"admin:read"}
			orgs = []string{"other"}
		}
		identities[name] = map[string]any{"subject": name, "tenantId": tenant, "permissions": perms, "organizations": orgs, "accessToken": token(name, tenant, perms, orgs)}
	}
	mux := http.NewServeMux()
	httpapi.EnterpriseQueryModule{Service: enterprisequery.NewService(db.NewEnterpriseQuery(pool)), Metrics: metrics}.Register(mux, verifier)
	httpapi.StoredValueModule{Service: stored}.Register(mux, verifier)
	httpapi.CustomerFeedbackModule{Service: customerfeedback.NewService(db.NewCustomerFeedback(pool))}.Register(mux, verifier)
	var mu sync.Mutex
	counts := map[string]int{}
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		counts[r.Method+" "+r.URL.Path]++
		mu.Unlock()
		mux.ServeHTTP(w, r)
	}))
	defer api.Close()
	web, node := os.Getenv("ELITE_WEB_ROOT"), os.Getenv("ELITE_NODE_BIN")
	if !filepath.IsAbs(web) || !filepath.IsAbs(node) {
		t.Fatal("fixed runtime/root required")
	}
	selection:=map[string]any{"schema_version":"1","tenants":map[string]any{tenant:map[string]any{"organizations":map[string]any{"store":map[string]any{"metricReferences":[]map[string]any{
	 {"value":program,"label":"Programa de ensayo","allowed_subjects":[]string{"admin"}},
	 {"value":"survey","label":"Encuesta de ensayo","allowed_subjects":[]string{"admin"}},
	}}}}}}
	selectionBytes,_:=json.Marshal(selection)
	if e=os.WriteFile(filepath.Join(web,"config","experience-options.json"),selectionBytes,0600);e!=nil{t.Fatal(e)}
	listener, e := net.Listen("tcp", "127.0.0.1:0")
	if e != nil {
		t.Fatal(e)
	}
	address := listener.Addr().String()
	listener.Close()
	target, _ := url.Parse("http://" + address)
	edge := httptest.NewTLSServer(httputil.NewSingleHostReverseProxy(target))
	defer edge.Close()
	env := []string{}
	for _, v := range os.Environ() {
		k := strings.ToUpper(strings.SplitN(v, "=", 2)[0])
		if strings.HasPrefix(k, "ELITE_") || strings.HasPrefix(k, "CATALOG_") || strings.HasPrefix(k, "PUBLIC_") || strings.HasPrefix(k, "ENTERPRISE_") || strings.Contains(k, "DATABASE_URL") || k == "APP_BASE_URL" || k == "AUTH_SESSION_SECRET" || k == "BUSINESS_CONFIG_FILE" {
			continue
		}
		env = append(env, v)
	}
	identitiesJSON, _ := json.Marshal(identities)
	env = append(env, "ELITE_ROLE_METRICS_BROWSER=1", "ELITE_WORKSPACE_E2E=enabled", "ELITE_METRIC_IDENTITIES="+string(identitiesJSON), "ELITE_WEB_ROOT="+web, "ELITE_BASE_URL="+edge.URL,
		"APP_BASE_URL="+edge.URL, "ENTERPRISE_API_BASE_URL="+api.URL, "AUTH_SESSION_SECRET=synthetic-role-metrics-"+uuid.NewString(), "BUSINESS_CONFIG_FILE=role.metrics.fixture.json",
		"CATALOG_RELEASE_ENABLED=false", "PUBLIC_SITE_ORIGIN=https://catalog.example.invalid", "PUBLIC_INDEXING_ENABLED=0", "ENTERPRISE_TENANT_CODE=role-catalog", "ENTERPRISE_ORGANIZATION_CODE=role-org", "NEXT_TELEMETRY_DISABLED=1")
	artifacts, e := os.MkdirTemp(web, "role-metrics-artifacts-")
	if e != nil {
		t.Fatal(e)
	}
	log, e := os.Create(filepath.Join(artifacts, "next.log"))
	if e != nil {
		t.Fatal(e)
	}
	defer log.Close()
	server := exec.CommandContext(ctx, node, filepath.Join(web, "node_modules", "next", "dist", "bin", "next"), "start", "--hostname", "127.0.0.1", "--port", strings.Split(address, ":")[1])
	server.Dir, server.Env, server.Stdout, server.Stderr = web, env, log, log
	if e = server.Start(); e != nil {
		t.Fatal(e)
	}
	defer func() { server.Process.Kill(); server.Wait() }()
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
	command := exec.CommandContext(ctx, node, filepath.Join(gate, "node_modules", "@playwright", "test", "cli.js"), "test", "tests/role-metrics-connected.spec.mjs", "--project=chromium-desktop", "--timeout=90000", "--output="+filepath.Join(artifacts, "browsers"))
	command.Dir, command.Env = gate, env
	output, e := command.CombinedOutput()
	if x := os.WriteFile(filepath.Join(artifacts, "browser.log"), output, 0600); x != nil {
		t.Fatal(x)
	}
	if e != nil || !strings.Contains(string(output), "1 passed") {
		t.Fatalf("role metrics browser %v\n%s\nartifacts=%s", e, output, artifacts)
	}

	var n int
	if e = pool.QueryRow(ctx, `select count(*)from crm.survey_response where tenant_id=$1`, tenant).Scan(&n); e != nil || n != 2 {
		t.Fatal("durable survey writes", n, e)
	}
	mu.Lock()
	rawCounts, _ := json.Marshal(counts)
	posts := 0
	for k, n := range counts {
		if strings.HasPrefix(k, "POST ") {
			posts += n
		}
	}
	mu.Unlock()
	if posts != 2 {
		t.Fatal("unexpected write", posts)
	}
	if e = os.WriteFile(filepath.Join(artifacts, "api-counts.json"), rawCounts, 0600); e != nil {
		t.Fatal(e)
	}
	t.Logf("ROLE_METRICS_BROWSER_PASS JWE_RS256_JWKS=true PG_BFF_exact_amounts=true scope_and_unavailable=true real_survey_deltas=2 metric_POSTs=0 artifacts=%s", artifacts)
}
