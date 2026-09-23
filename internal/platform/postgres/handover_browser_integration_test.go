package postgres_test

// AUTHORED local composition fixture. Real encrypted role sessions, RS256/JWKS
// bearer verification, Next BFF, admitted SDKs and PostgreSQL; no live login.
import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"math/big"
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

	"elite.local/enterprise/internal/commerce"
	"elite.local/enterprise/internal/franchisejourney"
	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/platform/identity"
	db "elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/platform/randomid"
)

type handoverBrowserClock struct{}

func (handoverBrowserClock) Now() time.Time { return time.Now() }

func TestHandoverBrowserPostgres(t *testing.T) {
	if os.Getenv("ELITE_HANDOVER_BROWSER") != "1" {
		t.Skip("explicit local handover browser fixture required")
	}
	r := newConnectedRun(t, false)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	if code := r.callback(t, "evt_browser_initial", "checkout.session.completed", "cs_test_fixture", true); code != 200 {
		t.Fatal(code)
	}
	if n, err := r.processor.ProcessOnce(ctx); err != nil || n != 1 {
		t.Fatal(n, err)
	}
	mode := false
	d := franchisejourney.HandoverProfileDocument{Schema: franchisejourney.HandoverProfileSchema, ProfileID: "browser-franchise", Revision: 1, Algorithm: franchisejourney.HandoverSupportedAlgorithm, AlgorithmRevision: 2, Scope: "MATERIALIZED_PROFILE", TenantID: r.tenant, OrganizationID: "store", ProviderCode: "stripe", ProviderAccountRef: "acct_fixture", ProviderConnectionID: "checkout", ExpectedLiveMode: &mode, MaximumObservationAgeSeconds: 900, Options: franchisejourney.SupportedCommercialReleaseOptions(), AuthorityReference: "docs/handover-operator-flow.md", DecisionReference: "LOCAL_BROWSER_FIXTURE"}
	raw, err := json.Marshal(d)
	if err != nil {
		t.Fatal(err)
	}
	digest := sha256.Sum256(raw)
	policy, err := franchisejourney.LoadHandoverProfile(raw, franchisejourney.HandoverProfileActivation{Enabled: true, ProfileID: d.ProfileID, ProfileRevision: 1, DocumentSHA256: hex.EncodeToString(digest[:]), TenantID: r.tenant, OrganizationID: "store", ProviderCode: "stripe", ProviderAccountRef: "acct_fixture", ProviderConnectionID: "checkout"})
	if err != nil {
		t.Fatal(err)
	}
	repo := db.NewFranchiseJourney(r.pool)
	ids := randomid.Generator{}
	service, err := franchisejourney.NewHandoverPreparationService(repo, ids, policy)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = repo.PublishDeliveryChecklist(ctx, r.tenant, "fixture-operator", franchisejourney.DeliveryChecklist{ID: "browser-checklist", OrganizationID: "store", Version: 1, Title: "Preparación de referencia", Items: []franchisejourney.ChecklistItem{{ID: "serial", Ordinal: 1, Prompt: "Verificar serie", ResponseType: "serial", Required: true}}}, ids.New()); err != nil {
		t.Fatal(err)
	}
	verifier, token := handoverBrowserIssuer(t)
	identities := map[string]any{}
	for _, name := range []string{"operator", "customer", "stranger", "foreign-org"} {
		permissions := []string{"customer:self"}
		organizations := []string{"store"}
		if name == "operator" || name == "foreign-org" {
			permissions = []string{"handover:manage"}
		}
		if name == "foreign-org" {
			organizations = []string{"other"}
		}
		identities[name] = map[string]any{"subject": name, "tenantId": r.tenant, "permissions": permissions, "organizations": organizations, "accessToken": token(name, r.tenant, permissions, organizations)}
	}
	mux := http.NewServeMux()
	httpapi.FranchiseJourneyModule{Service: franchisejourney.NewService(repo, ids, handoverBrowserClock{})}.Register(mux, verifier)
	httpapi.CommerceModule{Service: commerce.NewService(db.NewCommerce(r.pool), ids), PaymentProvider: "stripe", PaymentTenantID: r.tenant, PaymentOrganizationID: "store", ProviderObservedPayments: true}.Register(mux, verifier)
	httpapi.InitialHandoverModule{Service: service}.Register(mux, verifier)
	httpapi.HandoverContextModule{Service: service}.Register(mux, verifier)
	httpapi.CommercialReleaseModule{Service: service}.Register(mux, verifier)
	controlToken := ids.New()
	mux.HandleFunc("POST /__fixture/handover-observation", func(w http.ResponseWriter, request *http.Request) {
		if request.Header.Get("X-Fixture-Token") != controlToken {
			w.WriteHeader(403)
			return
		}
		var body struct {
			Action string `json:"action"`
		}
		if json.NewDecoder(http.MaxBytesReader(w, request.Body, 1024)).Decode(&body) != nil {
			w.WriteHeader(400)
			return
		}
		switch body.Action {
		case "hold":
			r.provider.mu.Lock()
			r.provider.refund = 1
			r.provider.mu.Unlock()
			if r.callback(t, "evt_browser_refund", "charge.refunded", "ch_fixture", true) != 200 {
				w.WriteHeader(500)
				return
			}
		case "reconcile":
			if n, e := r.processor.ProcessOnce(request.Context()); e != nil || n != 1 {
				w.WriteHeader(500)
				return
			}
		default:
			w.WriteHeader(400)
			return
		}
		w.WriteHeader(200)
	})
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
	env = append(env, "ELITE_HANDOVER_BROWSER=1", "ELITE_WORKSPACE_E2E=enabled", "ELITE_HANDOVER_IDENTITIES="+string(encoded), "ELITE_WEB_ROOT="+web, "ELITE_BASE_URL="+edge.URL, "APP_BASE_URL="+edge.URL, "ENTERPRISE_API_BASE_URL="+api.URL, "AUTH_SESSION_SECRET=synthetic-"+ids.New()+ids.New(), "ELITE_HANDOVER_CONTROL="+api.URL+"/__fixture/handover-observation", "ELITE_HANDOVER_CONTROL_TOKEN="+controlToken)
	artifacts, err := os.MkdirTemp(web, "handover-browser-artifacts-")
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
	command := exec.CommandContext(ctx, node, filepath.Join(gate, "node_modules", "@playwright", "test", "cli.js"), "test", "tests/handover-connected.spec.mjs", "--project=chromium-desktop", "--timeout=90000", "--output="+filepath.Join(artifacts, "browsers"))
	command.Dir, command.Env = gate, env
	output, err := command.CombinedOutput()
	if e := os.WriteFile(filepath.Join(artifacts, "browser.log"), output, 0600); e != nil {
		t.Fatal(e)
	}
	if err != nil || !strings.Contains(string(output), "1 passed") {
		t.Fatalf("browser %v\n%s\nartifacts=%s", err, output, artifacts)
	}
	var prepared, releases, releaseEvents, accepted int
	var stock string
	var hold bool
	err = r.pool.QueryRow(ctx, `select (select count(*) from sales.delivery_handover_preparation where tenant_id=$1),(select count(*) from sales.commercial_release_receipt where tenant_id=$1),(select count(*) from platform.outbox_event where tenant_id=$1 and event_type='commercial-release.recorded'),(select count(*) from sales.delivery_handover where tenant_id=$1 and state='accepted'),(select state from inventory.stock_unit where tenant_id=$1 and stock_unit_id='stock'),(select hold from payment.provider_observation where tenant_id=$1)`, r.tenant).Scan(&prepared, &releases, &releaseEvents, &accepted, &stock, &hold)
	if err != nil || prepared != 1 || releases != 1 || releaseEvents != 1 || accepted != 1 || stock != "reserved" || !hold {
		t.Fatal("durable outcomes", prepared, releases, releaseEvents, accepted, stock, hold, err)
	}
	mu.Lock()
	countsJSON, _ := json.Marshal(counts)
	mu.Unlock()
	if err = os.WriteFile(filepath.Join(artifacts, "api-post-counts.json"), countsJSON, 0600); err != nil {
		t.Fatal(err)
	}
	t.Logf("HANDOVER_BROWSER_POSTGRES_PASS browser=chromium-desktop mobile_viewport=true actual_BFF_SQL=true role_sessions_JWE_RS256_JWKS=true provider_SDK_callback=true preparations=1 acceptances=1 commercial_receipts=1 release_events=1 refunded_hold=true artifacts=%s", artifacts)
}

func handoverBrowserIssuer(t *testing.T) (identity.Verifier, func(string, string, []string, []string) string) {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	encode := base64.RawURLEncoding.EncodeToString
	var issuer *httptest.Server
	issuer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/.well-known/openid-configuration":
			_ = json.NewEncoder(w).Encode(map[string]any{"issuer": issuer.URL, "jwks_uri": issuer.URL + "/keys", "id_token_signing_alg_values_supported": []string{"RS256"}})
		case "/keys":
			_ = json.NewEncoder(w).Encode(map[string]any{"keys": []any{map[string]any{"kty": "RSA", "use": "sig", "alg": "RS256", "kid": "local-confirmation", "n": encode(key.N.Bytes()), "e": encode(big.NewInt(int64(key.E)).Bytes())}}})
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(issuer.Close)
	verifier, err := identity.NewOIDCVerifier(context.Background(), issuer.URL, "confirmation-api")
	if err != nil {
		t.Fatal(err)
	}
	token := func(subject, tenant string, permissions, organizations []string) string {
		header, err := json.Marshal(map[string]string{"alg": "RS256", "kid": "local-confirmation", "typ": "JWT"})
		if err != nil {
			t.Fatal(err)
		}
		claims, err := json.Marshal(map[string]any{"iss": issuer.URL, "aud": "confirmation-api", "sub": subject, "iat": time.Now().Unix(), "exp": time.Now().Add(5 * time.Minute).Unix(), "tenant_id": tenant, "permissions": permissions, "organization_ids": organizations})
		if err != nil {
			t.Fatal(err)
		}
		unsigned := encode(header) + "." + encode(claims)
		digest := sha256.Sum256([]byte(unsigned))
		signature, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, digest[:])
		if err != nil {
			t.Fatal(err)
		}
		return unsigned + "." + encode(signature)
	}
	return verifier, token
}
