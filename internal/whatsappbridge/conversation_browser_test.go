package whatsappbridge

// AUTHORED browser composition fixture. Existing conversation fixture supplies
// inbox/runtime/proposal/transport; real OIDC verifies the browser's backend token.
import (
	"context"
	"encoding/json"
	"fmt"
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

func TestWhatsAppHumanBrowserPostgres(t *testing.T) {
	if os.Getenv("ELITE_WHATSAPP_BROWSER") != "1" {
		t.Fatal("explicit WhatsApp browser fixture required")
	}
	f := newConnectedReplyFixture(t)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	for _, id := range []string{"wamid.browser-normal", "wamid.lost", "wamid.browser-reject"} {
		f.ingest(t, f.inbound(t, id, time.Now().Add(-time.Second), "5491112345678"))
		f.process(t)
	}
	if f.metaCalls.Load() != 0 || f.llmCalls.Load() != 6 || f.domainCalls.Load() != 3 {
		t.Fatal("proposal effects", f.metaCalls.Load(), f.llmCalls.Load(), f.domainCalls.Load())
	}
	verifier, token := whatsappBrowserIssuer(t)
	identities := map[string]any{}
	for _, name := range []string{"human", "review-only", "reader", "foreign-org"} {
		permissions := []string{"whatsapp:approve", "whatsapp:send"}
		orgs := []string{"store-1"}
		subject := "fixture-human"
		if name == "review-only" {
			permissions = []string{"whatsapp:approve"}
			subject = "review-only"
		}
		if name == "reader" {
			permissions = []string{"lead:read"}
			subject = "reader"
		}
		if name == "foreign-org" {
			orgs = []string{"store-2"}
			subject = "foreign-org"
		}
		identities[name] = map[string]any{"subject": subject, "tenantId": f.human.TenantID, "permissions": permissions, "organizations": orgs, "accessToken": token(subject, f.human.TenantID, permissions, orgs)}
	}
	mux := http.NewServeMux()
	f.module.Register(mux, verifier)
	controlToken := fmt.Sprintf("fixture-%d", time.Now().UnixNano())
	mux.HandleFunc("POST /__fixture/delivered", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Fixture-Token") != controlToken {
			w.WriteHeader(403)
			return
		}
		var body struct {
			Index int `json:"index"`
		}
		if json.NewDecoder(r.Body).Decode(&body) != nil || body.Index < 1 || body.Index > 2 {
			w.WriteHeader(400)
			return
		}
		f.ingest(t, statusBatch("delivered", fmt.Sprint(time.Now().Unix()), map[string]any{"id": fmt.Sprintf("wamid.reply.%d", body.Index)}))
		f.process(t)
		w.WriteHeader(200)
	})
	var mu sync.Mutex
	counts := map[string]int{}
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" && !strings.HasPrefix(r.URL.Path, "/__fixture/") {
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
	env = append(env, "ELITE_WHATSAPP_BROWSER=1", "ELITE_WORKSPACE_E2E=enabled", "ELITE_WHATSAPP_IDENTITIES="+string(encoded), "ELITE_WEB_ROOT="+web, "ELITE_BASE_URL="+edge.URL, "APP_BASE_URL="+edge.URL, "ENTERPRISE_API_BASE_URL="+api.URL, "AUTH_SESSION_SECRET=synthetic-browser-only-"+controlToken, "ELITE_WHATSAPP_CONTROL="+api.URL+"/__fixture/delivered", "ELITE_WHATSAPP_CONTROL_TOKEN="+controlToken)
	artifacts, err := os.MkdirTemp(web, "whatsapp-browser-artifacts-")
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
	command := exec.CommandContext(ctx, node, filepath.Join(gate, "node_modules", "@playwright", "test", "cli.js"), "test", "tests/whatsapp-connected.spec.mjs", "--project=chromium-desktop", "--timeout=90000", "--output="+filepath.Join(artifacts, "browsers"))
	command.Dir, command.Env = gate, env
	output, err := command.CombinedOutput()
	if e := os.WriteFile(filepath.Join(artifacts, "browser.log"), output, 0600); e != nil {
		t.Fatal(e)
	}
	if err != nil || !strings.Contains(string(output), "1 passed") {
		t.Fatalf("browser %v\n%s\nartifacts=%s", err, output, artifacts)
	}
	var accepted, attempts, observations, approved, rejected int
	err = f.pool.QueryRow(ctx, `select (select count(*) from communication.outbound_delivery where tenant_id=$1 and state='accepted'),(select coalesce(sum(attempt_count),0) from communication.outbound_delivery where tenant_id=$1),(select count(*) from communication.whatsapp_status_observation where tenant_id=$1),(select count(*) from approval.request where tenant_id=$1 and state='approved'),(select count(*) from approval.request where tenant_id=$1 and state='rejected')`, f.human.TenantID).Scan(&accepted, &attempts, &observations, &approved, &rejected)
	if err != nil || accepted != 2 || attempts != 2 || observations != 2 || approved != 2 || rejected != 1 || f.metaCalls.Load() != 2 {
		t.Fatal("durable outcomes", accepted, attempts, observations, approved, rejected, f.metaCalls.Load(), err)
	}
	mu.Lock()
	countsJSON, _ := json.Marshal(counts)
	mu.Unlock()
	if err = os.WriteFile(filepath.Join(artifacts, "api-post-counts.json"), countsJSON, 0600); err != nil {
		t.Fatal(err)
	}
	t.Logf("WHATSAPP_BROWSER_POSTGRES_PASS actual_BFF_SQL=true JWE_RS256_JWKS=true runtime_proposals=3 approvals=2 rejected=1 provider_POSTs=2 accepted=2 delivered=2 lost_receipt_recovered_without_POST=true reader_navigation_hidden=true send_permission_enforced=true artifacts=%s", artifacts)
}
