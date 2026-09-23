package httpapi

import (
	"context"
	"elite.local/enterprise/internal/customerfeedback"
	"elite.local/enterprise/internal/platform/postgres"
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
	"sync/atomic"
	"testing"
	"time"
)

func TestFeedbackBrowserPostgres(t *testing.T) {
	if os.Getenv("ELITE_FEEDBACK_E2E") != "1" {
		t.Skip("explicit disposable survey browser gate not requested")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	cfg, err := pgxpool.ParseConfig(os.Getenv("TEST_DATABASE_URL"))
	if err != nil || cfg.ConnConfig.Host != "127.0.0.1" || !strings.HasPrefix(cfg.ConnConfig.Database, "elite_feedback_") {
		t.Fatal("disposable loopback elite_feedback_* database required")
	}
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := randomid.Generator{}.New()
	for _, q := range []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,'feedback-browser','Synthetic','Synthetic')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,'org-a','org-a','A','store'),($1,'org-b','org-b','B','store')`,
		`insert into crm.customer_profile(tenant_id,customer_principal_id,display_name,email_normalized,status) values($1,'alice','Synthetic Alice','a@example.test','active'),($1,'bob','Synthetic Bob','b@example.test','active')`,
	} {
		if _, err = pool.Exec(ctx, q, tenant); err != nil {
			t.Fatal(err)
		}
	}
	for _, id := range []string{"chromium-desktop", "chromium-mobile", "firefox-desktop", "webkit-desktop"} {
		_, err = pool.Exec(ctx, `insert into crm.survey_definition(tenant_id,organization_id,survey_id,prompt,consent_version,consent_notice,opens_at,closes_at,retain_until,minimum_responses,active) values($1,'org-a',$2,'Synthetic survey <img src=x onerror=alert(1)>','fixture-v1','Synthetic fixture only',clock_timestamp()-interval '1 hour',clock_timestamp()+interval '1 hour',clock_timestamp()+interval '1 day',2,true)`, tenant, id)
		if err != nil {
			t.Fatal(err)
		}
	}
	verifier, token := confirmationTestIssuer(t)
	identities := map[string]any{}
	for _, name := range []string{"alice", "bob", "admin"} {
		permissions := []string{"customer:self"}
		if name == "admin" {
			permissions = []string{"surveys:read"}
		}
		identities[name] = map[string]any{"subject": name, "tenantId": tenant, "permissions": permissions, "organizations": []string{"org-a"}, "accessToken": token(name, tenant, permissions, []string{"org-a"})}
	}
	mux := http.NewServeMux()
	repository := postgres.NewCustomerFeedback(pool)
	CustomerFeedbackModule{Service: customerfeedback.NewService(repository)}.Register(mux, verifier)
	var writes atomic.Int64
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			writes.Add(1)
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
		if strings.HasPrefix(name, "ELITE_") || name == "DATABASE_URL" || name == "TEST_DATABASE_URL" || name == "APP_BASE_URL" || name == "ENTERPRISE_API_BASE_URL" || name == "AUTH_SESSION_SECRET" {
			continue
		}
		env = append(env, item)
	}
	encoded, err := json.Marshal(identities)
	if err != nil {
		t.Fatal(err)
	}
	env = append(env, "ELITE_FEEDBACK_E2E=1", "ELITE_WORKSPACE_E2E=enabled", "ELITE_FEEDBACK_IDENTITIES="+string(encoded), "ELITE_WEB_ROOT="+web, "ELITE_BASE_URL="+edge.URL, "APP_BASE_URL="+edge.URL, "ENTERPRISE_API_BASE_URL="+api.URL, "AUTH_SESSION_SECRET=synthetic-"+randomid.Generator{}.New()+randomid.Generator{}.New())
	artifacts, err := os.MkdirTemp(web, "survey-browser-artifacts-")
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
	defer func() { server.Process.Kill(); server.Wait() }()
	ready := false
	client := &http.Client{Timeout: time.Second}
	for deadline := time.Now().Add(20 * time.Second); time.Now().Before(deadline); {
		response, e := client.Get(target.String() + "/icon.svg")
		if e == nil {
			response.Body.Close()
			if response.StatusCode == 200 {
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
	command := exec.CommandContext(ctx, "node", filepath.Join(gate, "node_modules", "@playwright", "test", "cli.js"), "test", "tests/customer-survey.spec.mjs", "--timeout=35000", "--output="+filepath.Join(artifacts, "browsers"))
	command.Dir, command.Env = gate, env
	output, err := command.CombinedOutput()
	if e := os.WriteFile(filepath.Join(artifacts, "browser.log"), output, 0600); e != nil {
		t.Fatal(e)
	}
	if err != nil || !strings.Contains(string(output), "4 passed") {
		t.Fatalf("browser: %v\n%s", err, output)
	}
	var rows int
	if err = pool.QueryRow(ctx, `select count(*) from crm.survey_response where tenant_id=$1`, tenant).Scan(&rows); err != nil || rows != 8 || writes.Load() != 8 {
		t.Fatalf("responses=%d writes=%d err=%v", rows, writes.Load(), err)
	}
	t.Logf("FEEDBACK_BROWSER_POSTGRES_PASS browsers=4 durable_responses=8 writes=8 artifacts=%s", artifacts)
	t.Run("retention_scope_and_limit", func(t *testing.T) {
		for _, org := range []string{"org-a", "org-b"} {
			if _, e := pool.Exec(ctx, `insert into crm.survey_definition(tenant_id,organization_id,survey_id,prompt,consent_version,consent_notice,opens_at,closes_at,retain_until,minimum_responses)values($1,$2,'expired','Synthetic expired','v1','Synthetic',clock_timestamp()-interval '3 day',clock_timestamp()-interval '2 day',clock_timestamp()-interval '1 day',2)`, tenant, org); e != nil {
				t.Fatal(e)
			}
			if _, e := pool.Exec(ctx, `insert into crm.survey_response(tenant_id,organization_id,survey_id,customer_principal_id,score,consent_version)values($1,$2,'expired','alice',9,'v1'),($1,$2,'expired','bob',0,'v1')`, tenant, org); e != nil {
				t.Fatal(e)
			}
		}
		scope := customerfeedback.Scope{Tenant: tenant, Organization: "org-a", Customer: "alice"}
		if _, e := repository.OwnAnswer(ctx, scope, "expired"); e != customerfeedback.ErrNotFound {
			t.Fatalf("expired answer=%v", e)
		}
		if _, e := repository.Definition(ctx, scope, "expired"); e != customerfeedback.ErrUnavailable {
			t.Fatalf("expired definition=%v", e)
		}
		if _, e := repository.Summary(ctx, scope, "expired"); e != customerfeedback.ErrUnavailable {
			t.Fatalf("expired summary=%v", e)
		}
		for _, limit := range []int{0, 1001} {
			if _, e := repository.PurgeExpired(ctx, tenant, "org-a", limit); e != customerfeedback.ErrInvalid {
				t.Fatalf("limit: %v", e)
			}
		}
		for _, want := range []int64{1, 1, 0} {
			n, e := repository.PurgeExpired(ctx, tenant, "org-a", 1)
			if e != nil || n != want {
				t.Fatalf("purge=%d want=%d err=%v", n, want, e)
			}
		}
		var other, live int
		if e := pool.QueryRow(ctx, `select count(*) filter(where organization_id='org-b'),count(*) filter(where survey_id<>'expired') from crm.survey_response where tenant_id=$1`, tenant).Scan(&other, &live); e != nil || other != 2 || live != 8 {
			t.Fatalf("retention crossed boundary other=%d live=%d err=%v", other, live, e)
		}
		t.Log("FEEDBACK_RETENTION_PASS expired_reads_blocked=true bounded=1 foreign_org_unchanged=true future_unchanged=true")
	})
}
