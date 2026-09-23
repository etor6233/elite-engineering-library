package postgres_test

// AUTHORED local composition fixture over the existing operations owner.
import (
	"context"
	"elite.local/enterprise/internal/enterprisequery"
	"elite.local/enterprise/internal/operations"
	"elite.local/enterprise/internal/platform/httpapi"
	db "elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/platform/randomid"
	"encoding/json"
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

func TestFactoryBrowserPostgres(t *testing.T) {
	if os.Getenv("ELITE_FACTORY_BROWSER") != "1" {
		t.Skip("explicit local factory browser fixture required")
	}
	rawDB := os.Getenv("FACTORY_BROWSER_DB_URL")
	parsed, err := url.Parse(rawDB)
	if err != nil || parsed.Hostname() != "127.0.0.1" || !strings.HasPrefix(parsed.Path, "/elite_payment_connected_") {
		t.Fatal("isolated loopback fixture database required")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	pool, err := pgxpool.New(ctx, rawDB)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	ids := randomid.Generator{}
	tenant := ids.New()
	otherTenant := ids.New()
	for _, q := range []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1,'factory-browser','Fixture','Fixture')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'store','store','Store','warehouse'),($1,'other','other','Other','warehouse')`,
		`insert into partner.supplier(tenant_id,supplier_id,supplier_code,legal_name,status)values($1,'supplier','supplier','Supplier','active')`,
		`insert into catalog.vehicle_model(tenant_id,model_id,model_code,display_name,vehicle_class,lifecycle_state)values($1,'model','model','Model','bicycle','active')`,
		`insert into catalog.vehicle_variant(tenant_id,variant_id,model_id,variant_code,display_name,battery_specification,lifecycle_state)values($1,'variant','model','variant','Variant','{}','active')`,
	} {
		if _, err = pool.Exec(ctx, q, tenant); err != nil {
			t.Fatal(err)
		}
	}
	repo := db.NewOperations(pool)
	service := operations.NewService(repo, ids)
	if err = repo.CreatePurchaseOrder(ctx, tenant, ids.New(), operations.PurchaseOrder{ID: "po", SupplierID: "supplier", DestinationOrganizationID: "store", State: "draft", Currency: "USD", TotalMinorUnits: 100, Version: 1}); err != nil {
		t.Fatal(err)
	}
	if err = repo.CreateProductionUnit(ctx, tenant, ids.New(), operations.ProductionUnit{ID: "unit", OrganizationID: "store", PurchaseOrderID: "po", VariantID: "variant", SerialNumber: "SERIAL-FACTORY-FIXTURE", State: "planned"}); err != nil {
		t.Fatal(err)
	}
	verifier, token := handoverBrowserIssuer(t)
	identities := map[string]any{}
	for _, name := range []string{"operator", "reader", "foreign-org", "foreign-tenant"} {
		permissions := []string{"factory:read", "factory:write"}
		orgs := []string{"store"}
		tenantID := tenant
		if name == "reader" {
			permissions = []string{"factory:read"}
		}
		if name == "foreign-org" {
			orgs = []string{"other"}
		}
		if name == "foreign-tenant" {
			tenantID = otherTenant
		}
		identities[name] = map[string]any{"subject": name, "tenantId": tenantID, "permissions": permissions, "organizations": orgs, "accessToken": token(name, tenantID, permissions, orgs)}
	}
	mux := http.NewServeMux()
	httpapi.OperationsModule{Service: service}.Register(mux, verifier)
	httpapi.EnterpriseQueryModule{Service: enterprisequery.NewService(db.NewEnterpriseQuery(pool))}.Register(mux, verifier)
	controlToken := ids.New()
	mux.HandleFunc("POST /__fixture/factory-competing-advance", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Fixture-Token") != controlToken {
			w.WriteHeader(403)
			return
		}
		if err := service.TransitionProductionUnit(r.Context(), tenant, "store", "unit", "assembly", "quality"); err != nil {
			w.WriteHeader(409)
			return
		}
		w.WriteHeader(200)
	})
	var mu sync.Mutex
	counts := map[string]int{}
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && !strings.HasPrefix(r.URL.Path, "/__fixture/") {
			mu.Lock()
			counts[r.URL.Path]++
			mu.Unlock()
		}
		mux.ServeHTTP(w, r)
	}))
	defer api.Close()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	address := listener.Addr().String()
	listener.Close()
	target, _ := url.Parse("http://" + address)
	edge := httptest.NewTLSServer(httputil.NewSingleHostReverseProxy(target))
	defer edge.Close()
	web := os.Getenv("ELITE_WEB_ROOT")
	if !filepath.IsAbs(web) {
		t.Fatal("absolute web root required")
	}
	env := []string{}
	for _, item := range os.Environ() {
		name := strings.ToUpper(strings.SplitN(item, "=", 2)[0])
		if strings.HasPrefix(name, "ELITE_") || name == "DATABASE_URL" || name == "TEST_DATABASE_URL" || name == "PAYMENT_CONNECTED_DB_URL" || name == "FACTORY_BROWSER_DB_URL" || name == "APP_BASE_URL" || name == "ENTERPRISE_API_BASE_URL" || name == "AUTH_SESSION_SECRET" {
			continue
		}
		env = append(env, item)
	}
	encoded, err := json.Marshal(identities)
	if err != nil {
		t.Fatal(err)
	}
	env = append(env, "ELITE_FACTORY_BROWSER=1", "ELITE_WORKSPACE_E2E=enabled", "ELITE_FACTORY_IDENTITIES="+string(encoded), "ELITE_WEB_ROOT="+web, "ELITE_BASE_URL="+edge.URL, "APP_BASE_URL="+edge.URL, "ENTERPRISE_API_BASE_URL="+api.URL, "AUTH_SESSION_SECRET=synthetic-"+ids.New()+ids.New(), "ELITE_FACTORY_CONTROL="+api.URL+"/__fixture/factory-competing-advance", "ELITE_FACTORY_CONTROL_TOKEN="+controlToken)
	artifacts, err := os.MkdirTemp(web, "factory-browser-artifacts-")
	if err != nil {
		t.Fatal(err)
	}
	log, err := os.Create(filepath.Join(artifacts, "next.log"))
	if err != nil {
		t.Fatal(err)
	}
	defer log.Close()
	node := os.Getenv("ELITE_NODE_BIN")
	if !filepath.IsAbs(node) {
		t.Fatal("exact Node executable required")
	}
	server := exec.CommandContext(ctx, node, filepath.Join(web, "node_modules", "next", "dist", "bin", "next"), "start", "--hostname", "127.0.0.1", "--port", strings.Split(address, ":")[1])
	server.Dir, server.Env = web, env
	server.Stdout, server.Stderr = log, log
	if err = server.Start(); err != nil {
		t.Fatal(err)
	}
	defer func() { server.Process.Kill(); server.Wait() }()
	client := &http.Client{Timeout: time.Second}
	ready := false
	for deadline := time.Now().Add(20 * time.Second); time.Now().Before(deadline); {
		resp, e := client.Get(target.String() + "/icon.svg")
		if e == nil {
			resp.Body.Close()
			if resp.StatusCode == 200 {
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
		t.Fatal("Next did not start")
	}
	gate := filepath.Join(web, "microsoft_playwright_browser_gate")
	command := exec.CommandContext(ctx, node, filepath.Join(gate, "node_modules", "@playwright", "test", "cli.js"), "test", "tests/factory-connected.spec.mjs", "--project=chromium-desktop", "--timeout=90000", "--output="+filepath.Join(artifacts, "browsers"))
	command.Dir, command.Env = gate, env
	output, err := command.CombinedOutput()
	if e := os.WriteFile(filepath.Join(artifacts, "browser.log"), output, 0600); e != nil {
		t.Fatal(e)
	}
	if err != nil || !strings.Contains(string(output), "1 passed") {
		t.Fatalf("browser %v\n%s\nartifacts=%s", err, output, artifacts)
	}

	var state string
	var assembly, quality, stock int
	err = pool.QueryRow(ctx, `select (select state from factory.production_unit where tenant_id=$1 and production_unit_id='unit'),(select count(*) from platform.outbox_event where tenant_id=$1 and event_type='production-unit.assembly'),(select count(*) from platform.outbox_event where tenant_id=$1 and event_type='production-unit.quality'),(select count(*) from inventory.stock_unit where tenant_id=$1)`, tenant).Scan(&state, &assembly, &quality, &stock)
	if err != nil || state != "quality" || assembly != 1 || quality != 1 || stock != 0 {
		t.Fatal("durable outcomes", state, assembly, quality, stock, err)
	}
	mu.Lock()
	countJSON, _ := json.Marshal(counts)
	unitPosts := counts["/v1/factory/units/unit/transitions"]
	mu.Unlock()
	if unitPosts != 1 {
		t.Fatal("unexpected backend mutation attempts", unitPosts)
	}
	if err = os.WriteFile(filepath.Join(artifacts, "api-post-counts.json"), countJSON, 0600); err != nil {
		t.Fatal(err)
	}
	t.Logf("FACTORY_BROWSER_POSTGRES_PASS actual_BFF_SQL=true JWE_RS256_JWKS=true assembly_events=1 competing_quality_events=1 backend_transition_posts=1 response_loss_observation_only=true cross_org_tenant_role_stale_rejected=true stock_mutation=false artifacts=%s", artifacts)
}
