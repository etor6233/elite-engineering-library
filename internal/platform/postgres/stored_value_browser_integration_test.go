package postgres_test

// AUTHORED connected fixture: distinct reviewers, exact source calculation,
// JWE role sessions, real OIDC/JWKS, Next BFF and durable zero-provider delivery.
import (
	"context"
	"crypto/sha256"
	"elite.local/enterprise/internal/commerce"
	"elite.local/enterprise/internal/franchisejourney"
	"elite.local/enterprise/internal/platform/httpapi"
	db "elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/platform/randomid"
	"encoding/hex"
	"encoding/json"
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

func TestStoredValueBrowserPostgres(t *testing.T) {
	if os.Getenv("ELITE_STORED_VALUE_BROWSER") != "1" {
		t.Skip("explicit stored value browser fixture required")
	}
	pool := connectedPool(t)
	tenant := connectedSeedOrder(t, pool, "stripe")
	store, _, _ := fixtureStoredValueSetup(t, pool, tenant, "gift_card", 200000)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	mode := false
	d := franchisejourney.HandoverProfileDocument{Schema: franchisejourney.HandoverProfileSchema, ProfileID: "browser-franchise", Revision: 1, Algorithm: franchisejourney.HandoverSupportedAlgorithm, AlgorithmRevision: 3, Scope: "MATERIALIZED_PROFILE", TenantID: tenant, OrganizationID: "store", ProviderCode: "stripe", ProviderAccountRef: "acct_fixture", ProviderConnectionID: "checkout", ExpectedLiveMode: &mode, MaximumObservationAgeSeconds: 900, Options: franchisejourney.SupportedStoredValueReleaseOptions(), AuthorityReference: "docs/handover-operator-flow.md", DecisionReference: "LOCAL_BROWSER_FIXTURE"}
	raw, err := json.Marshal(d)
	if err != nil {
		t.Fatal(err)
	}
	digest := sha256.Sum256(raw)
	policy, err := franchisejourney.LoadHandoverProfile(raw, franchisejourney.HandoverProfileActivation{Enabled: true, ProfileID: d.ProfileID, ProfileRevision: 1, DocumentSHA256: hex.EncodeToString(digest[:]), TenantID: tenant, OrganizationID: "store", ProviderCode: "stripe", ProviderAccountRef: "acct_fixture", ProviderConnectionID: "checkout"})
	if err != nil {
		t.Fatal(err)
	}
	repo := db.NewFranchiseJourney(pool)
	ids := randomid.Generator{}
	service, err := franchisejourney.NewHandoverPreparationService(repo, ids, policy)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = repo.PublishDeliveryChecklist(ctx, tenant, "fixture-operator", franchisejourney.DeliveryChecklist{ID: "browser-checklist", OrganizationID: "store", Version: 1, Title: "Preparación de referencia", Items: []franchisejourney.ChecklistItem{{ID: "serial", Ordinal: 1, Prompt: "Verificar serie", ResponseType: "serial", Required: true}}}, ids.New()); err != nil {
		t.Fatal(err)
	}
	verifier, token := handoverBrowserIssuer(t)
	identities := map[string]any{}
	for _, name := range []string{"operator", "reviewer", "customer", "stranger", "foreign-org"} {
		permissions := []string{"customer:self"}
		organizations := []string{"store"}
		if name == "operator" || name == "reviewer" || name == "foreign-org" {
			permissions = []string{"handover:manage", "stored_value:read", "stored_value:request", "stored_value:approve", "stored_value:fund"}
		}
		if name == "foreign-org" {
			organizations = []string{"other"}
		}
		identities[name] = map[string]any{"subject": name, "tenantId": tenant, "permissions": permissions, "organizations": organizations, "accessToken": token(name, tenant, permissions, organizations)}
	}
	mux := http.NewServeMux()
	httpapi.FranchiseJourneyModule{Service: franchisejourney.NewService(repo, ids, handoverBrowserClock{})}.Register(mux, verifier)
	httpapi.CommerceModule{Service: commerce.NewService(db.NewCommerce(pool), ids), PaymentProvider: "stripe", PaymentTenantID: tenant, PaymentOrganizationID: "store", ProviderObservedPayments: true}.Register(mux, verifier)
	httpapi.InitialHandoverModule{Service: service}.Register(mux, verifier)
	httpapi.HandoverContextModule{Service: service}.Register(mux, verifier)
	httpapi.CommercialReleaseModule{Service: service}.Register(mux, verifier)
	httpapi.StoredValueModule{Service: store}.Register(mux, verifier)
	var mu sync.Mutex
	counts := map[string]int{}
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		if request.Method == "POST" && !strings.HasPrefix(request.URL.Path, "/__fixture/") {
			mu.Lock()
			counts[request.URL.Path]++
			mu.Unlock()
		}
		mux.ServeHTTP(w, request)
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
		if strings.HasPrefix(name, "ELITE_") || name == "DATABASE_URL" || name == "TEST_DATABASE_URL" || name == "PAYMENT_CONNECTED_DB_URL" || name == "APP_BASE_URL" || name == "ENTERPRISE_API_BASE_URL" || name == "AUTH_SESSION_SECRET" {
			continue
		}
		env = append(env, item)
	}
	encoded, err := json.Marshal(identities)
	if err != nil {
		t.Fatal(err)
	}
	env = append(env, "ELITE_STORED_VALUE_BROWSER=1", "ELITE_WORKSPACE_E2E=enabled", "ELITE_HANDOVER_IDENTITIES="+string(encoded), "ELITE_WEB_ROOT="+web, "ELITE_BASE_URL="+edge.URL, "APP_BASE_URL="+edge.URL, "ENTERPRISE_API_BASE_URL="+api.URL, "AUTH_SESSION_SECRET=synthetic-"+ids.New()+ids.New())
	artifacts, err := os.MkdirTemp(web, "stored-value-browser-artifacts-")
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(filepath.Join(artifacts, "profile.json"), raw, 0600); err != nil {
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
	command := exec.CommandContext(ctx, node, filepath.Join(gate, "node_modules", "@playwright", "test", "cli.js"), "test", "tests/stored-value-connected.spec.mjs", "--project=chromium-desktop", "--timeout=90000", "--output="+filepath.Join(artifacts, "browsers"))
	command.Dir, command.Env = gate, env
	output, err := command.CombinedOutput()
	if e := os.WriteFile(filepath.Join(artifacts, "browser.log"), output, 0600); e != nil {
		t.Fatal(e)
	}
	if err != nil || !strings.Contains(string(output), "1 passed") {
		t.Fatalf("browser %v\n%s\nartifacts=%s", err, output, artifacts)
	}

	var prepared, releases, accepted, operations, funding, payments int
	err = pool.QueryRow(ctx, `select (select count(*) from sales.delivery_handover_preparation where tenant_id=$1),(select count(*) from sales.commercial_release_receipt where tenant_id=$1),(select count(*) from sales.delivery_handover where tenant_id=$1 and state='accepted'),(select count(*) from stored_value.operation where tenant_id=$1),(select count(*) from payment.local_funding_receipt where tenant_id=$1),(select count(*) from payment.payment_attempt where tenant_id=$1)`, tenant).Scan(&prepared, &releases, &accepted, &operations, &funding, &payments)
	if err != nil || prepared != 1 || releases != 1 || accepted != 1 || operations != 2 || funding != 1 || payments != 0 {
		t.Fatal("durable outcomes", prepared, releases, accepted, operations, funding, payments, err)
	}
	mu.Lock()
	countsJSON, _ := json.Marshal(counts)
	mu.Unlock()
	if err = os.WriteFile(filepath.Join(artifacts, "api-post-counts.json"), countsJSON, 0600); err != nil {
		t.Fatal(err)
	}
	t.Logf("STORED_VALUE_BROWSER_POSTGRES_PASS actual_BFF_SQL=true JWE_RS256_JWKS=true source_operations=2 funding=1 payments=0 handover=1 commercial_receipt=1 artifacts=%s", artifacts)
}
