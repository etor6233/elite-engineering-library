package trainingbridge

// AUTHORED real-browser fixture over the actual module and durable owners.
import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
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

func TestTrainingHumanBrowserPostgres(t *testing.T) {
	if os.Getenv("ELITE_TRAINING_BROWSER") != "1" {
		t.Skip("explicit local training browser fixture required")
	}
	f := connectedFixture(t)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	verifier, token := trainingBrowserIssuer(t)
	identities := map[string]any{}
	for _, name := range []string{"learner", "reviewer", "other-learner", "foreign-org", "reader"} {
		permissions := []string{"training:learn"}
		orgs := []string{"store-1"}
		if name == "reviewer" {
			permissions = []string{"training:review"}
		}
		if name == "foreign-org" {
			orgs = []string{"other"}
		}
		if name == "reader" {
			permissions = []string{}
		}
		identities[name] = map[string]any{"subject": name, "tenantId": f.learner.TenantID, "permissions": permissions, "organizations": orgs, "accessToken": token(name, f.learner.TenantID, permissions, orgs)}
	}
	mux := http.NewServeMux()
	Module{Store: f.store}.Register(mux, verifier)
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
	defer api.Close()
	controlToken := uuid.NewString()
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
	node := os.Getenv("ELITE_NODE_BIN")
	if !filepath.IsAbs(web) || !filepath.IsAbs(node) {
		t.Fatal("fixed web/runtime paths required")
	}
	env := []string{}
	for _, item := range os.Environ() {
		name := strings.ToUpper(strings.SplitN(item, "=", 2)[0])
		if strings.HasPrefix(name, "ELITE_") || strings.Contains(name, "DATABASE_URL") || name == "APP_BASE_URL" || name == "ENTERPRISE_API_BASE_URL" || name == "AUTH_SESSION_SECRET" {
			continue
		}
		env = append(env, item)
	}
	encoded, _ := json.Marshal(identities)
	env = append(env, "ELITE_TRAINING_BROWSER=1", "ELITE_WORKSPACE_E2E=enabled", "ELITE_TRAINING_IDENTITIES="+string(encoded), "ELITE_WEB_ROOT="+web, "ELITE_BASE_URL="+edge.URL, "APP_BASE_URL="+edge.URL, "ENTERPRISE_API_BASE_URL="+api.URL, "AUTH_SESSION_SECRET=synthetic-browser-only-"+controlToken, "ELITE_TRAINING_FIXTURE=1")
	artifacts, err := os.MkdirTemp(web, "training-browser-artifacts-")
	if err != nil {
		t.Fatal(err)
	}
	log, err := os.Create(filepath.Join(artifacts, "next.log"))
	if err != nil {
		t.Fatal(err)
	}
	defer log.Close()
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
		t.Fatal("Next did not start", artifacts)
	}
	gate := filepath.Join(web, "microsoft_playwright_browser_gate")
	command := exec.CommandContext(ctx, node, filepath.Join(gate, "node_modules", "@playwright", "test", "cli.js"), "test", "tests/training-connected.spec.mjs", "--project=chromium-desktop", "--timeout=90000", "--output="+filepath.Join(artifacts, "browsers"))
	command.Dir, command.Env = gate, env
	output, err := command.CombinedOutput()
	if e := os.WriteFile(filepath.Join(artifacts, "browser.log"), output, 0600); e != nil {
		t.Fatal(e)
	}
	if err != nil || !strings.Contains(string(output), "1 passed") {
		t.Fatalf("browser %v\n%s\nartifacts=%s", err, output, artifacts)
	}

	var facts, events, requests, decisions, resources int
	err = f.pool.QueryRow(ctx, `select (select count(*) from audit.event where tenant_id=$1 and resource_type='training-attempt'),(select count(*) from platform.outbox_event where tenant_id=$1 and aggregate_type='training-evidence'),(select count(*) from approval.request where tenant_id=$1 and kind='training_assessment'),(select count(*) from approval.decision where tenant_id=$1),(select count(*) from crm.service_resource where tenant_id=$1)`, f.learner.TenantID).Scan(&facts, &events, &requests, &decisions, &resources)
	if err != nil || facts != 4 || events != 4 || requests != 1 || decisions != 1 || resources != 0 {
		t.Fatal("unexpected durable outcomes", facts, events, requests, decisions, resources, err)
	}
	mu.Lock()
	encodedCounts, _ := json.Marshal(counts)
	total := 0
	for _, count := range counts {
		total += count
	}
	mu.Unlock()
	if total != 4 {
		t.Fatal("duplicate or unauthorized backend writes", total)
	}
	if err = os.WriteFile(filepath.Join(artifacts, "api-post-counts.json"), encodedCounts, 0600); err != nil {
		t.Fatal(err)
	}
	t.Logf("TRAINING_BROWSER_POSTGRES_PASS actual_BFF_SQL=true JWE_RS256_JWKS=true learner_start_read_response=true human_assessment=true response_loss_read_only_recovery=true actor_org_isolation=true durable_facts=4 outbox=4 backend_posts=4 permission_grants=0 artifacts=%s", artifacts)
}
