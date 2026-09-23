package httpapi

import (
	"context"
	"elite.local/enterprise/internal/electromobility"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/platform/randomid"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// AUTHORED integration harness: real admitted handler/repository and Microsoft
// Playwright runtime. No successful authentication double, provider or production claim.
type browserDenyVerifier struct{}

func (browserDenyVerifier) Verify(context.Context, string) (identity.Principal, error) {
	return identity.Principal{}, identity.ErrUnauthenticated
}

type publicBrowserFixture struct {
	testName      string
	setup         func(context.Context, *testing.T, *pgxpool.Pool, string) []EnterpriseModule
	verify        func(context.Context, *testing.T, *pgxpool.Pool, string)
	cleanupTables []string
}

func TestPublicLeadBrowserPostgres(t *testing.T) {
	runPublicBrowserPostgres(t, publicBrowserFixture{testName: "public lead crosses"})
}

func runPublicBrowserPostgres(t *testing.T, fixture publicBrowserFixture) {
	t.Helper()
	if os.Getenv("ELITE_PUBLIC_LEAD_E2E") != "1" {
		t.Skip("explicit disposable browser/Go/PostgreSQL gate not requested")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Minute)
	defer cancel()
	config, err := pgxpool.ParseConfig(os.Getenv("TEST_DATABASE_URL"))
	if err != nil {
		t.Fatal("invalid disposable database configuration")
	}
	if config.ConnConfig.Host != "127.0.0.1" || !strings.HasPrefix(config.ConnConfig.Database, "elite_browser_") {
		t.Fatal("requires explicit loopback database named elite_browser_*; never use a project database")
	}
	webRoot := os.Getenv("ELITE_WEB_ROOT")
	if !filepath.IsAbs(webRoot) {
		t.Fatal("ELITE_WEB_ROOT must be an absolute, freshly built enterprise-web composition")
	}
	gate := filepath.Join(webRoot, "microsoft_playwright_browser_gate")
	cli := filepath.Join(gate, "node_modules", "@playwright", "test", "cli.js")
	for _, path := range []string{cli, filepath.Join(webRoot, ".next", "BUILD_ID")} {
		if _, err := os.Stat(path); err != nil {
			t.Fatal("missing built web or installed pinned Playwright runtime")
		}
	}
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		t.Fatal("cannot initialize disposable database pool")
	}
	defer pool.Close()
	generator := randomid.Generator{}
	tenant := generator.New()
	code := "browser-" + strings.ReplaceAll(tenant, "-", "")
	_, err = pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1,$2,'Browser Fixture','Browser Fixture')`, tenant, code)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		cleanup, done := context.WithTimeout(context.Background(), 10*time.Second)
		defer done()
		// Only rows belonging to the random tenant created above are removed.
		for _, table := range append(fixture.cleanupTables, "platform.outbox_event", "platform.idempotency_record", "crm.consent_evidence", "crm.lead", "catalog.vehicle_model", "org.organization", "platform.tenant") {
			if _, err := pool.Exec(cleanup, "delete from "+table+" where tenant_id=$1", tenant); err != nil {
				t.Errorf("fixture cleanup failed for %s: %v", table, err)
			}
		}
	}()
	if _, err = pool.Exec(ctx, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'browser-org','browser-store','Browser Store','store')`, tenant); err != nil {
		t.Fatal(err)
	}
	if _, err = pool.Exec(ctx, `insert into catalog.vehicle_model(tenant_id,model_id,model_code,display_name,vehicle_class,lifecycle_state,publicly_visible,specification)values($1,'browser-model','browser-model','Browser Fixture Model','bicycle','active',true,'{}')`, tenant); err != nil {
		t.Fatal(err)
	}
	var modules []EnterpriseModule
	if fixture.setup != nil {
		modules = fixture.setup(ctx, t, pool, tenant)
	}
	server := httptest.NewServer(NewEnterprise(nil, browserDenyVerifier{}, electromobility.NewService(postgres.NewElectromobility(pool), generator), modules...))
	defer server.Close()
	command := exec.CommandContext(ctx, "node", cli, "test", "tests/enterprise-web.spec.mjs", "--grep", fixture.testName, "--workers=1", "--retries=0", "--max-failures=1")
	command.Dir = gate
	for _, entry := range os.Environ() {
		name, _, _ := strings.Cut(entry, "=")
		switch strings.ToUpper(name) {
		case "DATABASE_URL", "TEST_DATABASE_URL", "ELITE_BASE_URL", "ELITE_RUNTIME_ONLY", "ENTERPRISE_API_BASE_URL", "ENTERPRISE_TENANT_CODE", "ENTERPRISE_ORGANIZATION_CODE":
			continue
		}
		command.Env = append(command.Env, entry)
	}
	command.Env = append(command.Env, "ENTERPRISE_API_BASE_URL="+server.URL, "ENTERPRISE_TENANT_CODE="+code, "ENTERPRISE_ORGANIZATION_CODE=browser-store")
	output, err := command.CombinedOutput()
	t.Log(string(output))
	if err != nil {
		t.Fatalf("connected browser gate failed: %v", err)
	}
	var leads, consents, events, receipts int
	for _, check := range []struct {
		sql   string
		count *int
	}{
		{`select count(*) from crm.lead where tenant_id=$1 and organization_id='browser-org' and model_id='browser-model' and lifecycle_state='new' and source_code='public-web' and contact_payload->>'email' like '%@example.invalid'`, &leads},
		{`select count(*) from crm.consent_evidence where tenant_id=$1 and decision='granted'`, &consents},
		{`select count(*) from platform.outbox_event where tenant_id=$1 and event_type='lead.captured' and not (payload ? 'email') and not (payload ? 'contact_payload')`, &events},
		{`select count(*) from platform.idempotency_record where tenant_id=$1 and scope='public-lead' and status='completed'`, &receipts},
	} {
		if err := pool.QueryRow(ctx, check.sql, tenant).Scan(check.count); err != nil {
			t.Fatal(err)
		}
	}
	if leads != 4 || consents != 4 || events != 4 || receipts != 4 {
		t.Fatalf("durable invariant failed: leads=%d consents=%d events=%d receipts=%d; want four, no duplicates", leads, consents, events, receipts)
	}
	t.Log("PUBLIC_LEAD_BROWSER_POSTGRES_PASS browsers=4 durable_leads=4 consent=4 outbox=4 idempotency=4")
	if fixture.verify != nil {
		fixture.verify(ctx, t, pool, tenant)
	}
}

type emRepo struct {
	models  []electromobility.Model
	created int
	lead    int
}

func (r *emRepo) ListPublicModels(context.Context, string) ([]electromobility.Model, error) {
	return r.models, nil
}
func (r *emRepo) CreateModel(_ context.Context, _ string, _ string, _ electromobility.Model) error {
	r.created++
	return nil
}
func (r *emRepo) ResolvePublicOrganization(context.Context, string, string) (string, string, error) {
	return "018f4d4a-7b36-7a21-8d10-2f4c54c28f01", "org", nil
}
func (r *emRepo) CreateLead(_ context.Context, lead electromobility.Lead, _, _ string) (string, bool, error) {
	r.lead++
	return lead.ID, false, nil
}

type emIDs struct{ n int }

func (i *emIDs) New() string {
	i.n++
	return "018f4d4a-7b36-7a21-8d10-" + []string{"2f4c54c28f10", "2f4c54c28f11", "2f4c54c28f12", "2f4c54c28f13", "2f4c54c28f14", "2f4c54c28f15"}[i.n-1]
}

type emVerifier struct{ principal identity.Principal }

func (v emVerifier) Verify(context.Context, string) (identity.Principal, error) {
	return v.principal, nil
}
func TestElectromobilityHTTPBoundaries(t *testing.T) {
	repo := &emRepo{models: []electromobility.Model{{ID: "m", Code: "urban", DisplayName: "Urban", VehicleClass: "bicycle", Specification: json.RawMessage(`{}`)}}}
	ids := &emIDs{}
	service := electromobility.NewService(repo, ids)
	verifier := emVerifier{principal: identity.Principal{Subject: "admin", TenantID: "018f4d4a-7b36-7a21-8d10-2f4c54c28f01", Permissions: map[string]struct{}{"catalog:write": {}}}}
	handler := NewEnterprise(nil, verifier, service)
	request := httptest.NewRequest("GET", "/v1/public/acme/models", nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != 200 {
		t.Fatalf("list status=%d", response.Code)
	}
	request = httptest.NewRequest("POST", "/v1/catalog/models", strings.NewReader(`{"code":"urban-one","displayName":"Urban One","vehicleClass":"bicycle","specification":{}}`))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer valid")
	response = httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != 201 || repo.created != 1 {
		t.Fatalf("create status=%d created=%d body=%s", response.Code, repo.created, response.Body.String())
	}
	request = httptest.NewRequest("POST", "/v1/public/acme/store/leads", strings.NewReader(`{"model_id":"m","source_code":"public-web","contact":{"email":"person@example.test"},"consent_granted":true}`))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Idempotency-Key", "lead-request-00000001")
	response = httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != 202 || repo.lead != 1 {
		t.Fatalf("lead status=%d leads=%d body=%s", response.Code, repo.lead, response.Body.String())
	}
	request = httptest.NewRequest("POST", "/v1/public/acme/store/leads", strings.NewReader(`{"source_code":"public-web","contact":{},"consent_granted":false}`))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Idempotency-Key", "lead-request-00000002")
	response = httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != 400 || repo.lead != 1 {
		t.Fatalf("consent status=%d leads=%d", response.Code, repo.lead)
	}
}
