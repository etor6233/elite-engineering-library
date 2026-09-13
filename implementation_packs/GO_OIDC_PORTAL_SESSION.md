# OIDC portal durable session composition

## 1. Metadata

```yaml
pack_id: "GO-OIDC-PORTAL-SESSION"
pack_version: "0.1.1"
status:
  authority: SUPPORTED_REFERENCE
  implementation: RECONSTRUCTIBLE
  admission: CONDITIONED
claim: "Hash-bound portal lifecycle: official OIDC code/PKCE/refresh/revocation, encrypted random browser handle and bound token vault, PostgreSQL CAS/logout/revocation sweep and optional joined host. AUTHORED glue; no IdP/password implementation or live certification."
stacks: ["Go 1.26.8", "PostgreSQL 18.6", "Node.js 24.20.0 where frontend selected"]
compatible_with: ["GO-ELECTROMOBILITY-APPLICATION", "GO-OIDC-SERVICE-TOKEN-BROKER", "TS-OIDC-PORTAL-ADAPTER"]
incompatible_with: ["unbound tenant/provider/account", "production certification inferred from fixtures", "implicit source or business-policy attribution"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://github.com/panva/openid-client", "https://github.com/panva/jose", "https://github.com/coreos/go-oidc", "https://github.com/jackc/pgx"]
verified_at: "2026-09-13"
```

## 2. Applicability

Selected alongside the existing TypeScript OIDC portal, original Go application and official service-token broker; only the exact configured provider/tenant/org/service identity can access storage.

## 3. Architecture contract

Purpose-separated jose vault/handle; session/profile/ciphertext binding; one durable refresh owner, no ambiguous grant replay, local logout before provider revocation; finite worker sweep/retention. No role grants or IdP administration invented.

## 4. Exact file manifest

```text
CREATE docs/J5_IDP_ACCESS_CONTRACT.md
CREATE internal/platform/postgres/portal_access_review_integration_test.go
CREATE cmd/electromobility-api/portal_host_integration_test.go
CREATE cmd/electromobility-api/portal_maintenance.go
CREATE cmd/electromobility-api/portal_maintenance_test.go
CREATE cmd/electromobility-api/portal_session.go
CREATE config/identity/portal-lifecycle.example.json
CREATE db/migrations/0062_portal_session.down.sql
CREATE db/migrations/0062_portal_session.up.sql
CREATE docs/PORTAL_IDENTITY_LIFECYCLE.md
CREATE internal/platform/httpapi/portal_session.go
CREATE internal/platform/identity/portal_profile.go
CREATE internal/platform/identity/portal_profile_test.go
CREATE internal/platform/identity/portal_service_broker.go
CREATE internal/platform/identity/portal_session.go
CREATE internal/platform/postgres/portal_session.go
CREATE internal/platform/postgres/portal_session_integration_test.go
CREATE internal/platform/postgres/portal_session_sweep.go
```

## 5. Materialization blocks

### FILE: `cmd/electromobility-api/portal_host_integration_test.go`

```yaml
block_id: "GO-OIDC-PORTAL-SESSION:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "a81b7a90da10976df22e2c4a212c4ce041423c44882c7748d46d50298b468cff"
variables: []
secrets_allowed: false
```

````go
package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"elite.local/enterprise/internal/platform/identity"
	jose "github.com/go-jose/go-jose/v4"
)

// Exercise the host selector and domain/worker identity adapter with the real
// pinned client-credentials and RS256/JWKS verifiers, not a token provider stub.
func TestPortalHostUsesOIDCServiceBroker(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	var requests, grants atomic.Int32
	var server *httptest.Server
	var verifier *identity.OIDCVerifier
	var polls atomic.Int32
	observed := make(chan struct{}, 1)
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/.well-known/openid-configuration":
			_ = json.NewEncoder(w).Encode(map[string]any{"issuer": server.URL, "authorization_endpoint": server.URL + "/authorize", "token_endpoint": server.URL + "/token", "jwks_uri": server.URL + "/jwks", "response_types_supported": []string{"code"}, "subject_types_supported": []string{"public"}, "id_token_signing_alg_values_supported": []string{"RS256"}, "grant_types_supported": []string{"client_credentials"}, "token_endpoint_auth_methods_supported": []string{"client_secret_basic"}})
		case "/api/internal/oidc/session-maintenance":
			principal, e := verifier.Verify(r.Context(), strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer "))
			if e != nil || principal.Subject != "fixture-worker" || r.Method != "POST" {
				t.Error("unauthenticated maintenance")
				w.WriteHeader(401)
				return
			}
			polls.Add(1)
			_ = json.NewEncoder(w).Encode(map[string]int{"claimed": 0, "confirmed": 0, "pending": 0, "purged": 0, "unconfirmed_purged": 0})
			select {
			case observed <- struct{}{}:
			default:
			}
		case "/jwks":
			_ = json.NewEncoder(w).Encode(jose.JSONWebKeySet{Keys: []jose.JSONWebKey{{Key: &key.PublicKey, KeyID: "fixture-key", Use: "sig", Algorithm: "RS256"}}})
		case "/token":
			grants.Add(1)
			_ = r.ParseForm()
			id, secret, ok := r.BasicAuth()
			if r.Method != "POST" || !ok || id != "portal-worker" || secret != "synthetic-secret" || r.Form.Get("grant_type") != "client_credentials" || r.Form.Get("scope") != "portal-session:manage" {
				t.Error("invalid grant contract")
				w.WriteHeader(400)
				return
			}
			now := time.Now()
			raw, _ := json.Marshal(map[string]any{"iss": server.URL, "sub": "fixture-worker", "aud": "enterprise-api", "iat": now.Unix(), "exp": now.Add(120 * time.Second).Unix(), "tenant_id": "tenant", "organization_ids": []string{"org"}, "permissions": []string{"portal-session:manage"}, "scope": "portal-session:manage", "client_id": "portal-worker"})
			signer, e := jose.NewSigner(jose.SigningKey{Algorithm: jose.RS256, Key: key}, (&jose.SignerOptions{}).WithHeader("kid", "fixture-key"))
			if e != nil {
				t.Error(e)
				w.WriteHeader(500)
				return
			}
			signed, e := signer.Sign(raw)
			if e != nil {
				t.Error(e)
				w.WriteHeader(500)
				return
			}
			token, e := signed.CompactSerialize()
			if e != nil {
				t.Error(e)
				w.WriteHeader(500)
				return
			}
			_ = json.NewEncoder(w).Encode(map[string]any{"token_type": "Bearer", "expires_in": 120, "access_token": token})
		default:
			w.WriteHeader(404)
		}
	}))
	defer server.Close()
	ctx := context.Background()
	verifier, err = identity.NewOIDCVerifier(ctx, server.URL, "enterprise-api")
	if err != nil {
		t.Fatal(err)
	}
	directory := t.TempDir()
	profileFile := filepath.Join(directory, "profile.json")
	secretFile := filepath.Join(directory, "secret")
	var document map[string]any
	if json.Unmarshal([]byte(`{"schema": "elite.oidc.portal-lifecycle.v1", "profile_id": "portal-host-fixture", "revision": 1, "transport": "LOOPBACK_FIXTURE", "issuer": "https://issuer.invalid", "authorization_endpoint": "https://issuer.invalid/authorize", "token_endpoint": "https://issuer.invalid/token", "jwks_endpoint": "https://issuer.invalid/jwks", "revocation_endpoint": "https://issuer.invalid/revoke", "end_session_endpoint": "https://issuer.invalid/logout", "client_id": "reference-portal", "callback_url": "https://portal.invalid/api/auth/callback", "post_logout_url": "https://portal.invalid/", "bridge_url": "https://backend.invalid", "backend_audience": "enterprise-api", "service_client_id": "portal-worker", "service_subject": "fixture-worker", "service_scopes": ["portal-session:manage"], "service_audience_parameter": "", "tenant_id": "tenant", "organization_ids": ["org"], "allowed_permissions": ["customer:read"], "scopes": ["openid", "offline_access"], "maximum_session_seconds": 3600, "refresh_before_seconds": 30, "retention_seconds": 86400}`), &document) != nil {
		t.Fatal("fixture")
	}
	for field, path := range map[string]string{"issuer": "", "authorization_endpoint": "/authorize", "token_endpoint": "/token", "jwks_endpoint": "/jwks", "revocation_endpoint": "/revoke", "end_session_endpoint": "/logout", "callback_url": "/api/auth/callback", "post_logout_url": "/", "bridge_url": ""} {
		document[field] = server.URL + path
	}
	raw, _ := json.Marshal(document)
	hash := sha256.Sum256(raw)
	if os.WriteFile(profileFile, raw, 0600) != nil || os.WriteFile(secretFile, []byte("synthetic-secret"), 0600) != nil {
		t.Fatal("fixture files")
	}
	values := map[string]string{"OIDC_PORTAL_LIFECYCLE_ENABLED": "true", "OIDC_PORTAL_PROFILE_FILE": profileFile, "OIDC_PORTAL_PROFILE_SHA256": hex.EncodeToString(hash[:]), "OIDC_PORTAL_SERVICE_CLIENT_SECRET_FILE": secretFile, "OIDC_ISSUER": server.URL, "OIDC_AUDIENCE": "enterprise-api"}
	lookup := func(k string) string { return values[k] }
	pool, e := pgxpool.New(ctx, "postgres://fixture:fixture@127.0.0.1:1/constructor_only?sslmode=disable")
	if e != nil {
		t.Fatal(e)
	}
	defer pool.Close()
	before := requests.Load()
	if _, e = selectedPortalHost(ctx, pool, lookup); e == nil {
		t.Fatal("loopback activated without explicit fixture optin")
	}
	if requests.Load() != before {
		t.Fatal("rejected config contacted issuer")
	}
	values["OIDC_PORTAL_ALLOW_LOOPBACK_FIXTURE"] = "true"
	runtime, e := selectedPortalHost(ctx, pool, lookup)
	if e != nil {
		t.Fatal(e)
	}
	if runtime == nil {
		t.Fatal("missing runtime")
	}
	worker, cancel := context.WithCancel(ctx)
	defer cancel()
	done := make(chan struct{})
	go func() { defer close(done); runtime.run(worker) }()
	select {
	case <-observed:
	case <-time.After(3 * time.Second):
		t.Fatal("host maintenance missing")
	}
	cancel()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("host shutdown not joined")
	}
	if grants.Load() != 1 || polls.Load() != 1 {
		t.Fatal("unexpected grant/poll counts", grants.Load(), polls.Load())
	}
	t.Log("PORTAL_HOST_PASS exact_profile=true real_client_credentials_and_JWKS=true authenticated_maintenance=1 grants=1 joined_shutdown=true; PG behavior proven by separate connected fixture")
}
````

### FILE: `cmd/electromobility-api/portal_maintenance.go`

```yaml
block_id: "GO-OIDC-PORTAL-SESSION:file2:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "c237376b5ee8c1676a4ef95398a4dc0f4249877ed6f495065eb1a22f8c3168fa"
variables: []
secrets_allowed: false
```

````go
package main

import (
	"context"
	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/platform/identity"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"io"
	"log/slog"
	"mime"
	"net/http"
	"net/url"
	"time"
)

type portalSessionRuntime struct {
	*httpapi.PortalSessionModule
	broker interface {
		AccessToken(context.Context) (string, error)
	}
	target string
	client *http.Client
}

func init() { portalRuntimeFactory = preparePortalSessionRuntime }
func preparePortalSessionRuntime(ctx context.Context, pool *pgxpool.Pool, getenv func(string) string) (portalRuntime, error) {
	module, err := selectedPortalSessionModule(pool, getenv)
	if err != nil || module == nil {
		return nil, errors.New("portal session runtime unavailable")
	}
	secret := getenv("OIDC_PORTAL_SERVICE_CLIENT_SECRET_FILE")
	if secret == "" {
		return nil, errors.New("portal service credential reference required")
	}
	broker, err := identity.NewPortalServiceTokenBroker(ctx, module.Profile, identity.ServiceSecretFile(secret))
	if err != nil {
		return nil, err
	}
	base, err := url.Parse(module.Profile.Document().PostLogoutURL)
	if err != nil {
		return nil, err
	}
	target := base.ResolveReference(&url.URL{Path: "/api/internal/oidc/session-maintenance"}).String()
	return &portalSessionRuntime{PortalSessionModule: module, broker: broker, target: target, client: &http.Client{Timeout: 45 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return errors.New("portal maintenance redirect rejected") }}}, nil
}
func (p *portalSessionRuntime) poll(ctx context.Context) error {
	token, err := p.broker.AccessToken(ctx)
	if err != nil {
		return errors.New("portal service grant unavailable")
	}
	request, err := http.NewRequestWithContext(ctx, "POST", p.target, nil)
	if err != nil {
		return errors.New("portal maintenance request unavailable")
	}
	request.Header.Set("Authorization", "Bearer "+token)
	response, err := p.client.Do(request)
	if err != nil {
		return errors.New("portal maintenance transport unavailable")
	}
	defer response.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(response.Body, 4097))
	if err != nil || len(raw) > 4096 {
		return errors.New("portal maintenance response unavailable")
	}
	var result struct {
		Claimed           int `json:"claimed"`
		Confirmed         int `json:"confirmed"`
		Pending           int `json:"pending"`
		Purged            int `json:"purged"`
		UnconfirmedPurged int `json:"unconfirmed_purged"`
	}
	var fields map[string]json.RawMessage
	media, _, mediaErr := mime.ParseMediaType(response.Header.Get("Content-Type"))
	if mediaErr != nil || media != "application/json" || json.Unmarshal(raw, &fields) != nil || len(fields) != 5 || json.Unmarshal(raw, &result) != nil || response.StatusCode != 200 && response.StatusCode != 503 {
		return errors.New("portal maintenance result unavailable")
	}
	for _, name := range []string{"claimed", "confirmed", "pending", "purged", "unconfirmed_purged"} {
		if _, ok := fields[name]; !ok {
			return errors.New("portal maintenance result incomplete")
		}
	}
	if result.Claimed < 0 || result.Claimed > 2 || result.Confirmed < 0 || result.Pending < 0 || result.Confirmed+result.Pending != result.Claimed || result.Purged < 0 || result.Purged > 100 || result.UnconfirmedPurged < 0 || result.UnconfirmedPurged > result.Purged {
		return errors.New("portal maintenance counters invalid")
	}
	if result.UnconfirmedPurged > 0 {
		slog.Error("portal credentials removed at retention deadline without confirmed provider revocation", "count", result.UnconfirmedPurged)
	}
	if response.StatusCode != 200 || result.Pending > 0 {
		return errors.New("portal provider revocation pending")
	}
	return nil
}
func (p *portalSessionRuntime) run(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for {
		if ctx.Err() != nil {
			return
		}
		if err := p.poll(ctx); err != nil && ctx.Err() == nil {
			slog.Warn("portal session maintenance deferred")
		}
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}
````

### FILE: `cmd/electromobility-api/portal_maintenance_test.go`

```yaml
block_id: "GO-OIDC-PORTAL-SESSION:file3:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "62422442768b18763d7cb32cfa9275f0a989d175c8ad4a5a4d4a8e4fe28647d5"
variables: []
secrets_allowed: false
```

````go
package main

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type portalFixtureBroker struct{ failure bool }

func (b portalFixtureBroker) AccessToken(context.Context) (string, error) {
	if b.failure {
		return "", errors.New("fixture grant unavailable")
	}
	return "synthetic-service", nil
}
func TestPortalMaintenanceResponseBoundary(t *testing.T) {
	cases := []struct {
		name, body, media string
		status            int
		ok                bool
	}{
		{"empty", `{"claimed":0,"confirmed":0,"pending":0,"purged":0,"unconfirmed_purged":0}`, "application/json", 200, true},
		{"confirmed", `{"claimed":2,"confirmed":2,"pending":0,"purged":100,"unconfirmed_purged":1}`, "application/json", 200, true},
		{"pending", `{"claimed":1,"confirmed":0,"pending":1,"purged":0,"unconfirmed_purged":0}`, "application/json", 503, false},
		{"missing", `{"claimed":0,"confirmed":0,"pending":0,"purged":0}`, "application/json", 200, false},
		{"unknown", `{"claimed":0,"confirmed":0,"pending":0,"purged":0,"other":0}`, "application/json", 200, false},
		{"inconsistent", `{"claimed":2,"confirmed":1,"pending":0,"purged":0,"unconfirmed_purged":0}`, "application/json", 200, false},
		{"negative", `{"claimed":-1,"confirmed":0,"pending":0,"purged":0,"unconfirmed_purged":0}`, "application/json", 200, false},
		{"purge_bound", `{"claimed":0,"confirmed":0,"pending":0,"purged":101,"unconfirmed_purged":0}`, "application/json", 200, false},
		{"revocation_bound", `{"claimed":0,"confirmed":0,"pending":0,"purged":1,"unconfirmed_purged":2}`, "application/json", 200, false},
		{"claim_bound", `{"claimed":3,"confirmed":3,"pending":0,"purged":0,"unconfirmed_purged":0}`, "application/json", 200, false},
		{"wrong_media", `{"claimed":0,"confirmed":0,"pending":0,"purged":0,"unconfirmed_purged":0}`, "text/plain", 200, false},
		{"json_prefix", `{"claimed":0,"confirmed":0,"pending":0,"purged":0,"unconfirmed_purged":0}`, "application/jsonjunk", 200, false},
		{"wrong_status", `{"claimed":0,"confirmed":0,"pending":0,"purged":0,"unconfirmed_purged":0}`, "application/json", 500, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls++
				if r.Method != "POST" || r.Header.Get("Authorization") != "Bearer synthetic-service" {
					t.Error("request contract")
				}
				w.Header().Set("Content-Type", tc.media)
				w.WriteHeader(tc.status)
				_, _ = w.Write([]byte(tc.body))
			}))
			defer server.Close()
			p := portalSessionRuntime{broker: portalFixtureBroker{}, target: server.URL, client: server.Client()}
			if err := p.poll(context.Background()); (err == nil) != tc.ok {
				t.Fatalf("unexpected result %v", err)
			}
			if calls != 1 {
				t.Fatal("retried poll")
			}
		})
	}
}
func TestPortalMaintenanceCancellationAndGrant(t *testing.T) {
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls++; w.WriteHeader(500) }))
	defer server.Close()
	p := portalSessionRuntime{broker: portalFixtureBroker{failure: true}, target: server.URL, client: server.Client()}
	if p.poll(context.Background()) == nil || calls != 0 {
		t.Fatal("grant failure must not reach target")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	done := make(chan struct{})
	go func() { defer close(done); p.run(ctx) }()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("shutdown not joined")
	}
	if calls != 0 {
		t.Fatal("cancelled worker ran")
	}
}
func TestPortalActivationRequiresSelectedPack(t *testing.T) {
	old := portalRuntimeFactory
	defer func() { portalRuntimeFactory = old }()
	portalRuntimeFactory = nil
	for _, enabled := range []string{"true", "TRUE", "1"} {
		_, err := selectedPortalHost(context.Background(), nil, func(string) string { return enabled })
		if err == nil {
			t.Fatal("activation accepted", enabled)
		}
	}
	for _, enabled := range []string{"", "false"} {
		runtime, err := selectedPortalHost(context.Background(), nil, func(string) string { return enabled })
		if err != nil || runtime != nil {
			t.Fatal("disabled lifecycle not inert")
		}
	}
}
````

### FILE: `cmd/electromobility-api/portal_session.go`

```yaml
block_id: "GO-OIDC-PORTAL-SESSION:file4:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "cb7707bd3092ea577ef29f5c90838e1660f41e550e8839d663040c27bc32177c"
variables: []
secrets_allowed: false
```

````go
package main

import (
	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/postgres"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
)

func selectedPortalSessionModule(pool *pgxpool.Pool, getenv func(string) string) (*httpapi.PortalSessionModule, error) {
	enabled := getenv("OIDC_PORTAL_LIFECYCLE_ENABLED")
	if enabled == "" || enabled == "false" {
		return nil, nil
	}
	bad := errors.New("portal session configuration rejected")
	if enabled != "true" || pool == nil {
		return nil, bad
	}
	profile, err := identity.LoadPortalProfileFile(getenv("OIDC_PORTAL_PROFILE_FILE"), getenv("OIDC_PORTAL_PROFILE_SHA256"))
	if err != nil {
		return nil, bad
	}
	d := profile.Document()
	if d.Issuer != getenv("OIDC_ISSUER") || d.BackendAudience != getenv("OIDC_AUDIENCE") || (d.Transport != "TLS" && getenv("OIDC_PORTAL_ALLOW_LOOPBACK_FIXTURE") != "true") {
		return nil, bad
	}
	return &httpapi.PortalSessionModule{Profile: profile, Store: &postgres.PortalSessions{Pool: pool}}, nil
}
````

### FILE: `config/identity/portal-lifecycle.example.json`

```yaml
block_id: "GO-OIDC-PORTAL-SESSION:file5:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "6a57d5cb1c78b52f6ce641838ee22018d4a882853a283c6a6679699912f71835"
variables: []
secrets_allowed: false
```

````json
{
  "schema": "elite.oidc.portal-lifecycle.v1",
  "profile_id": "reference-portal",
  "revision": 1,
  "transport": "TLS",
  "issuer": "https://issuer.invalid",
  "authorization_endpoint": "https://issuer.invalid/authorize",
  "token_endpoint": "https://issuer.invalid/token",
  "jwks_endpoint": "https://issuer.invalid/jwks",
  "revocation_endpoint": "https://issuer.invalid/revoke",
  "end_session_endpoint": "https://issuer.invalid/logout",
  "client_id": "reference-portal",
  "callback_url": "https://portal.invalid/api/auth/callback",
  "post_logout_url": "https://portal.invalid/",
  "bridge_url": "https://backend.invalid",
  "backend_audience": "reference-api",
  "service_client_id": "reference-portal-maintenance",
  "service_subject": "reference-portal-maintenance",
  "service_scopes": [
    "portal-session:manage"
  ],
  "service_audience_parameter": "reference-api",
  "tenant_id": "00000000-0000-4000-8000-000000000001",
  "organization_ids": [
    "00000000-0000-4000-8000-000000000002"
  ],
  "allowed_permissions": [
    "customer:read"
  ],
  "scopes": [
    "openid",
    "offline_access"
  ],
  "maximum_session_seconds": 3600,
  "refresh_before_seconds": 30,
  "retention_seconds": 86400
}
````

### FILE: `db/migrations/0062_portal_session.down.sql`

```yaml
block_id: "GO-OIDC-PORTAL-SESSION:file6:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "50a6356e78b0557e2370198360f033d4fb6f6892bfc37c4f845f6ad8a2921387"
variables: []
secrets_allowed: false
```

````sql
begin;
do $$ begin
 if exists(select 1 from platform.portal_session) or exists(select 1 from platform.portal_session_retention_summary)
 then raise exception 'cannot remove portal session or revocation evidence';end if;
end $$;
drop table platform.portal_session;
drop table platform.portal_session_retention_summary;
commit;
````

### FILE: `db/migrations/0062_portal_session.up.sql`

```yaml
block_id: "GO-OIDC-PORTAL-SESSION:file7:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "49d7cdc600f6cb6915374368542fb63ff54fc713b40018b932c0cd97b8a33b1c"
variables: []
secrets_allowed: false
```

````sql
-- AUTHORED opaque session storage/CAS glue. No plaintext OAuth credentials.
create table platform.portal_session (
 profile_sha256 text not null check(profile_sha256 ~ '^[0-9a-f]{64}$'),
 session_id_sha256 bytea not null check(octet_length(session_id_sha256)=32),
 version bigint not null default 1 check(version>0),
 state text not null check(state in ('active','refreshing','reauth_required','revoked')),
 operation_id text not null default '',
 ciphertext text not null check(length(ciphertext)<=32768),
 ciphertext_sha256 bytea not null check(octet_length(ciphertext_sha256)=32),
 access_expires_at timestamptz not null,
 absolute_expires_at timestamptz not null,
 revocation_pending boolean not null default false,
 revocation_lease_until timestamptz not null default '-infinity',
 revocation_operation_id text not null default '',
 created_at timestamptz not null default clock_timestamp(),
 updated_at timestamptz not null default clock_timestamp(),
 primary key(profile_sha256,session_id_sha256),
 check(access_expires_at<=absolute_expires_at),
 check((state='active' and not revocation_pending) or state<>'active')
);
create index portal_session_expiry on platform.portal_session(absolute_expires_at);
create table platform.portal_session_retention_summary (
 profile_sha256 text primary key check(profile_sha256 ~ '^[0-9a-f]{64}$'),
 purged_count bigint not null default 0,
 unconfirmed_provider_revocation_count bigint not null default 0,
 updated_at timestamptz not null default clock_timestamp()
);
````

### FILE: `docs/PORTAL_IDENTITY_LIFECYCLE.md`

```yaml
block_id: "GO-OIDC-PORTAL-SESSION:file8:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "ea423631af4030bb4b83f762a3c4e085ba1a4364cd29ab2065ff060ab46daf61"
variables: []
secrets_allowed: false
```

````markdown
# Portal identity lifecycle — library reference

AUTHORED composition glue around the pinned openid-client/jose and original Go OIDC verifier/pgx; no password provider, corporate authorship or new dependency. Original IdP is the authority for bootstrap admins, user lifecycle, MFA and signed permissions. This adapter verifies the fixed tenant, issuer, audience, organization and permission allowlists and stores encrypted refresh credentials; it never grants a role or manufactures an admin.

Enable OIDC_PORTAL_LIFECYCLE_ENABLED=true on BFF and Go together after applying migration0062 (all selected migrations, in name order). OIDC_PORTAL_PROFILE_FILE names config/identity/portal-lifecycle.example.json adapted to the user's own target; OIDC_PORTAL_PROFILE_SHA256 is the SHA256 of those exact UTF8 bytes. Example .invalid hosts are deliberately unusable. BFF APP_BASE_URL must agree with callback_url; Go OIDC_ISSUER/OIDC_AUDIENCE must agree with issuer/backend_audience. Production requires TLS. LOOPBACK_FIXTURE is limited to numeric loopback and needs OIDC_PORTAL_ALLOW_LOOPBACK_FIXTURE=true on Go; the BFF rejects it in production.

Future secret file references: OIDC_PORTAL_CLIENT_SECRET_FILE (BFF client), OIDC_PORTAL_SERVICE_CLIENT_SECRET_FILE (BFF and Go maintenance client), OIDC_PORTAL_SESSION_KEY_FILE (BFF only, at least32characters). Files stay outside product/source and are never requested during library preparation. Register exact callback/logout URIs, code+S256 PKCE, offline_access, rotated refresh tokens, RS256 JWKS and client_secret_basic. The service principal needs exactly portal-session:manage with the configured tenant/organization/service subject; user tokens cannot use the private bridge.

Browser receives only a purpose-separated encrypted random session handle. The vault binds its authenticated session ID and profile hash, and PostgreSQL binds the ciphertext hash and fixed absolute lifetime. One durable CAS owner may send a refresh grant. Ambiguous grants require reauthentication; a consumed credential is never reused. Logout revokes locally before provider calls; failures stay pending and are visible. The Go host joins the maintenance worker at shutdown. Every30seconds it uses the official service grant to request at most2 revocations and100 expired-row removals. The BFF decrypts only leased, hash-matching rows and uses the official revocation endpoint. Unconfirmed provider revocation at retention expiry is counted durably, never called confirmed.

Rotation: rotate provider client credentials atomically in the external files; configuration/grant caches are bounded and backend access tokens expire. Replacing the session key or profile hash invalidates existing browser handles; revoke old provider sessions before removing the prior key/profile. No automatic key migration or silent role widening. Role withdrawal in the IdP is observed on issuance/refresh and the Go API independently validates each access token; already-issued bearer tokens remain bounded by their expiry unless the selected IdP/API policy adds immediate revocation. Provider account setup and target acceptance remain CONDITIONED. A library fixture does not prove a live IdP's compatibility.

Evidence: IDENTITY_PORTAL_RELEASE_V402.md/json. Integration uses the real pinned SDK, signed synthetic issuer, authenticated Go HTTP bridge and PostgreSQL. It covers code/PKCE/state/nonce, tenant/permission rejections, JWKS rotation, refresh contention, lost commit, unknown grant, ciphertext substitution, durable logout, revocation recovery and bounded retention. This is one T2803 claim; full composition SCA/source and J5 access administration remain separate gates until their receipts close.
````

### FILE: `internal/platform/httpapi/portal_session.go`

```yaml
block_id: "GO-OIDC-PORTAL-SESSION:file9:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "35fa22d7d394d7c6a907ec91aa54a8f93b3679a03bfd5d0e4e32afabeed8d105"
variables: []
secrets_allowed: false
```

````go
package httpapi

import (
	"bytes"
	"elite.local/enterprise/internal/platform/identity"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

type PortalSessionModule struct {
	Profile identity.PortalProfile
	Store   identity.PortalSessionStore
}

func (m PortalSessionModule) Register(mux *http.ServeMux, verifier identity.Verifier) {
	if m.Profile.SHA256() == "" || m.Store == nil {
		return
	}
	mux.HandleFunc("POST /v1/private/portal-sessions/sweep", func(w http.ResponseWriter, r *http.Request) {
		principal, err := authenticate(r.Context(), r.Header.Get("Authorization"), verifier)
		if err != nil {
			writeProblem(w, 401, "UNAUTHENTICATED", "service identity required")
			return
		}
		if !m.Profile.ServiceAllowed(principal) || r.Header.Get("X-Portal-Profile-SHA256") != m.Profile.SHA256() {
			writeProblem(w, 403, "FORBIDDEN", "portal service profile required")
			return
		}
		if r.Header.Get("Content-Type") != "application/json" {
			writeProblem(w, 415, "UNSUPPORTED_MEDIA_TYPE", "JSON required")
			return
		}
		var input struct {
			OperationID string `json:"operation_id"`
		}
		dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 256))
		dec.DisallowUnknownFields()
		if dec.Decode(&input) != nil || dec.Decode(new(any)) != io.EOF {
			writeProblem(w, 400, "INVALID_BODY", "bounded operation required")
			return
		}
		result, err := m.Store.Sweep(r.Context(), m.Profile, input.OperationID)
		if err != nil {
			writeProblem(w, 503, "SESSION_UNAVAILABLE", "session maintenance unavailable")
			return
		}
		writeJSON(w, 200, result)
	})
	for _, action := range []string{"create", "read", "claim", "commit", "abort", "revoke", "ack-revocation"} {
		mux.HandleFunc("POST /v1/private/portal-sessions/"+action, func(w http.ResponseWriter, r *http.Request) {
			principal, err := authenticate(r.Context(), r.Header.Get("Authorization"), verifier)
			if err != nil {
				writeProblem(w, 401, "UNAUTHENTICATED", "service identity required")
				return
			}
			if !m.Profile.ServiceAllowed(principal) || r.Header.Get("X-Portal-Profile-SHA256") != m.Profile.SHA256() {
				writeProblem(w, 403, "FORBIDDEN", "portal service profile required")
				return
			}
			if r.Header.Get("Content-Type") != "application/json" {
				writeProblem(w, 415, "UNSUPPORTED_MEDIA_TYPE", "JSON required")
				return
			}
			raw, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 40000))
			if err != nil {
				writeProblem(w, 413, "BODY_TOO_LARGE", "bounded JSON required")
				return
			}
			first := json.NewDecoder(bytes.NewReader(raw))
			token, err := first.Token()
			if err != nil || token != json.Delim('{') {
				writeProblem(w, 400, "INVALID_BODY", "object required")
				return
			}
			seen := map[string]bool{}
			for first.More() {
				key, e := first.Token()
				name, ok := key.(string)
				if e != nil || !ok || seen[name] {
					writeProblem(w, 400, "INVALID_BODY", "duplicate field")
					return
				}
				seen[name] = true
				var value json.RawMessage
				if first.Decode(&value) != nil {
					writeProblem(w, 400, "INVALID_BODY", "invalid field")
					return
				}
			}
			var c identity.PortalSessionCommand
			dec := json.NewDecoder(strings.NewReader(string(raw)))
			dec.DisallowUnknownFields()
			if dec.Decode(&c) != nil || dec.Decode(new(any)) != io.EOF {
				writeProblem(w, 400, "INVALID_BODY", "invalid session command")
				return
			}
			record, err := m.Store.Execute(r.Context(), m.Profile, action, c)
			if errors.Is(err, identity.ErrPortalSessionNotFound) {
				writeProblem(w, 404, "SESSION_NOT_FOUND", "session unavailable")
				return
			}
			if errors.Is(err, identity.ErrPortalSessionConflict) {
				writeProblem(w, 409, "SESSION_CONFLICT", "session state changed")
				return
			}
			if err != nil {
				writeProblem(w, 503, "SESSION_UNAVAILABLE", "session store unavailable")
				return
			}
			writeJSON(w, 200, record)
		})
	}
}
````

### FILE: `internal/platform/identity/portal_profile.go`

```yaml
block_id: "GO-OIDC-PORTAL-SESSION:file10:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "48c10ff1ffe20ca40336a47e2da828500fbbf3e19a60a3fed087d06313845388"
variables: []
secrets_allowed: false
```

````go
package identity

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/url"
	"os"
	"slices"
	"strings"
)

// PortalProfile is hash-bound application configuration, not issuer credentials.
type PortalProfile struct {
	document PortalProfileDocument
	digest   string
}
type PortalProfileDocument struct {
	Schema                   string   `json:"schema"`
	ProfileID                string   `json:"profile_id"`
	Revision                 int      `json:"revision"`
	Transport                string   `json:"transport"`
	Issuer                   string   `json:"issuer"`
	AuthorizationEndpoint    string   `json:"authorization_endpoint"`
	TokenEndpoint            string   `json:"token_endpoint"`
	JWKSEndpoint             string   `json:"jwks_endpoint"`
	RevocationEndpoint       string   `json:"revocation_endpoint"`
	EndSessionEndpoint       string   `json:"end_session_endpoint"`
	ClientID                 string   `json:"client_id"`
	CallbackURL              string   `json:"callback_url"`
	PostLogoutURL            string   `json:"post_logout_url"`
	BridgeURL                string   `json:"bridge_url"`
	BackendAudience          string   `json:"backend_audience"`
	ServiceClientID          string   `json:"service_client_id"`
	ServiceSubject           string   `json:"service_subject"`
	ServiceScopes            []string `json:"service_scopes"`
	ServiceAudienceParameter string   `json:"service_audience_parameter"`
	TenantID                 string   `json:"tenant_id"`
	Organizations            []string `json:"organization_ids"`
	AllowedPermissions       []string `json:"allowed_permissions"`
	Scopes                   []string `json:"scopes"`
	MaximumSessionSeconds    int      `json:"maximum_session_seconds"`
	RefreshBeforeSeconds     int      `json:"refresh_before_seconds"`
	RetentionSeconds         int      `json:"retention_seconds"`
}

func (p PortalProfile) SHA256() string { return p.digest }
func (p PortalProfile) Document() PortalProfileDocument {
	d := p.document
	d.Organizations = slices.Clone(d.Organizations)
	d.AllowedPermissions = slices.Clone(d.AllowedPermissions)
	d.Scopes = slices.Clone(d.Scopes)
	d.ServiceScopes = slices.Clone(d.ServiceScopes)
	return d
}
func LoadPortalProfileFile(path, digest string) (PortalProfile, error) {
	f, err := os.Open(path)
	if err != nil {
		return PortalProfile{}, ErrServiceTokenConfiguration
	}
	defer f.Close()
	raw, err := io.ReadAll(io.LimitReader(f, 16385))
	if err != nil {
		return PortalProfile{}, ErrServiceTokenConfiguration
	}
	return LoadPortalProfile(raw, digest)
}
func LoadPortalProfile(raw []byte, digest string) (PortalProfile, error) {
	reject := func() (PortalProfile, error) { return PortalProfile{}, ErrServiceTokenConfiguration }
	if len(raw) == 0 || len(raw) > 16384 || len(digest) != 64 {
		return reject()
	}
	sum := sha256.Sum256(raw)
	if hex.EncodeToString(sum[:]) != digest {
		return reject()
	}
	first := json.NewDecoder(bytes.NewReader(raw))
	tok, err := first.Token()
	if err != nil || tok != json.Delim('{') {
		return reject()
	}
	seen := map[string]bool{}
	for first.More() {
		key, err := first.Token()
		name, ok := key.(string)
		if err != nil || !ok || seen[name] {
			return reject()
		}
		seen[name] = true
		var item json.RawMessage
		if first.Decode(&item) != nil {
			return reject()
		}
	}
	if _, err = first.Token(); err != nil {
		return reject()
	}
	if _, err = first.Token(); err != io.EOF {
		return reject()
	}
	var d PortalProfileDocument
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	if dec.Decode(&d) != nil {
		return reject()
	}
	if d.Schema != "elite.oidc.portal-lifecycle.v1" || d.Revision != 1 || !serviceAtom(d.ProfileID) || !serviceAtom(d.ClientID) || !serviceAtom(d.ServiceClientID) || !serviceAtom(d.ServiceSubject) || !serviceAtom(d.TenantID) || !serviceAtom(d.BackendAudience) {
		return reject()
	}
	if d.Transport != "TLS" && d.Transport != "LOOPBACK_FIXTURE" {
		return reject()
	}
	for _, u := range []string{d.Issuer, d.AuthorizationEndpoint, d.TokenEndpoint, d.JWKSEndpoint, d.RevocationEndpoint, d.EndSessionEndpoint, d.CallbackURL, d.PostLogoutURL, d.BridgeURL} {
		if !serviceURL(u, d.Transport) {
			return reject()
		}
	}
	if !strings.HasSuffix(d.CallbackURL, "/api/auth/callback") || strings.TrimSuffix(d.CallbackURL, "/api/auth/callback")+"/" != d.PostLogoutURL || strings.HasSuffix(d.BridgeURL, "/") {
		return reject()
	}
	bridge, err := url.Parse(d.BridgeURL)
	if err != nil || bridge.Path != "" && bridge.Path != "/" {
		return reject()
	}
	if !serviceSet(d.Organizations) || !serviceSet(d.AllowedPermissions) || !serviceSet(d.ServiceScopes) || !serviceSet(d.Scopes) || !slices.Contains(d.Scopes, "openid") || !slices.Contains(d.Scopes, "offline_access") {
		return reject()
	}
	if d.ServiceAudienceParameter != "" && d.ServiceAudienceParameter != d.BackendAudience {
		return reject()
	}
	if d.MaximumSessionSeconds < 300 || d.MaximumSessionSeconds > 86400 || d.RefreshBeforeSeconds < 5 || d.RefreshBeforeSeconds > 300 || d.RefreshBeforeSeconds*2 >= d.MaximumSessionSeconds || d.RetentionSeconds < 3600 || d.RetentionSeconds > 604800 {
		return reject()
	}
	return PortalProfile{d, digest}, nil
}

// ServiceAllowed grants only this private storage capability under a fixed profile.
func (p PortalProfile) ServiceAllowed(principal Principal) bool {
	if p.digest == "" || principal.Subject != p.document.ServiceSubject || principal.TenantID != p.document.TenantID || len(principal.Permissions) != 1 {
		return false
	}
	if _, ok := principal.Permissions["portal-session:manage"]; !ok {
		return false
	}
	organizations := make([]string, 0, len(principal.Organizations))
	for id := range principal.Organizations {
		organizations = append(organizations, id)
	}
	return serviceEqual(organizations, p.document.Organizations)
}
````

### FILE: `internal/platform/identity/portal_profile_test.go`

```yaml
block_id: "GO-OIDC-PORTAL-SESSION:file11:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "c4224d1e5ae66f6cd79ac3a3d59cb401cc8201958b01f9f1e29c8458962efa45"
variables: []
secrets_allowed: false
```

````go
package identity

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"testing"
)

var portalProfileFixture = []byte(`{
  "schema": "elite.oidc.portal-lifecycle.v1",
  "profile_id": "reference-portal",
  "revision": 1,
  "transport": "TLS",
  "issuer": "https://issuer.invalid",
  "authorization_endpoint": "https://issuer.invalid/authorize",
  "token_endpoint": "https://issuer.invalid/token",
  "jwks_endpoint": "https://issuer.invalid/jwks",
  "revocation_endpoint": "https://issuer.invalid/revoke",
  "end_session_endpoint": "https://issuer.invalid/logout",
  "client_id": "reference-portal",
  "callback_url": "https://portal.invalid/api/auth/callback",
  "post_logout_url": "https://portal.invalid/",
  "bridge_url": "https://backend.invalid",
  "backend_audience": "reference-api",
  "service_client_id": "reference-portal-maintenance",
  "service_subject": "reference-portal-maintenance",
  "service_scopes": [
    "portal-session:manage"
  ],
  "service_audience_parameter": "reference-api",
  "tenant_id": "00000000-0000-4000-8000-000000000001",
  "organization_ids": [
    "00000000-0000-4000-8000-000000000002"
  ],
  "allowed_permissions": [
    "customer:read"
  ],
  "scopes": [
    "openid",
    "offline_access"
  ],
  "maximum_session_seconds": 3600,
  "refresh_before_seconds": 30,
  "retention_seconds": 86400
}`)

func portalProfileFixtureHash(raw []byte) string {
	sum := sha256.Sum256(raw)
	return hex.EncodeToString(sum[:])
}
func TestPortalProfileContract(t *testing.T) {
	p, err := LoadPortalProfile(portalProfileFixture, portalProfileFixtureHash(portalProfileFixture))
	if err != nil {
		t.Fatal(err)
	}
	doc := p.Document()
	doc.Organizations[0] = "foreign"
	if p.Document().Organizations[0] == "foreign" {
		t.Fatal("mutable profile")
	}
	for _, change := range []map[string]any{{"bridge_url": "https://backend.invalid/path"}, {"bridge_url": "https://backend.invalid/"}, {"allowed_permissions": []string{"*"}}, {"transport": "LOOPBACK_FIXTURE"}, {"post_logout_url": "https://foreign.invalid/"}} {
		var value map[string]any
		if json.Unmarshal(portalProfileFixture, &value) != nil {
			t.Fatal("fixture")
		}
		for k, v := range change {
			value[k] = v
		}
		raw, _ := json.Marshal(value)
		if _, err := LoadPortalProfile(raw, portalProfileFixtureHash(raw)); err == nil {
			t.Fatal("invalid profile accepted", change)
		}
	}
	duplicate := append([]byte(`{"profile_id":"foreign",`), portalProfileFixture[1:]...)
	if _, err := LoadPortalProfile(duplicate, portalProfileFixtureHash(duplicate)); err == nil {
		t.Fatal("duplicate")
	}
	if _, err := LoadPortalProfile(portalProfileFixture, "0000000000000000000000000000000000000000000000000000000000000000"); err == nil {
		t.Fatal("wrong hash")
	}
}
func FuzzPortalProfile(f *testing.F) {
	f.Add(portalProfileFixture)
	f.Add([]byte(`{"profile_id":"foreign","profile_id":"fixture"}`))
	f.Add([]byte("null"))
	f.Fuzz(func(t *testing.T, raw []byte) {
		p, err := LoadPortalProfile(raw, portalProfileFixtureHash(raw))
		if err == nil {
			if p.SHA256() != portalProfileFixtureHash(raw) || len(raw) > 16384 {
				t.Fatal("unbound profile")
			}
			d := p.Document()
			if d.MaximumSessionSeconds < 300 || d.MaximumSessionSeconds > 86400 || len(d.Organizations) == 0 {
				t.Fatal("unbounded accepted profile")
			}
		}
	})
}
````

### FILE: `internal/platform/identity/portal_service_broker.go`

```yaml
block_id: "GO-OIDC-PORTAL-SESSION:file12:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "1b3eb390f3097a64b5b1d717a1ce59b85d427e7bf490b9886919fa739dbc56b1"
variables: []
secrets_allowed: false
```

````go
package identity

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

// Derived strictly from the hash-bound portal service identity; no body tenant.
func NewPortalServiceTokenBroker(ctx context.Context, p PortalProfile, reader ServiceSecretReader) (*ServiceTokenBroker, error) {
	if p.SHA256() == "" {
		return nil, ErrServiceTokenConfiguration
	}
	d := p.Document()
	document := serviceTokenDocument{Schema: "elite.oidc.service-token.v1", ProfileID: d.ProfileID, Revision: 1, Transport: d.Transport, Issuer: d.Issuer, TokenEndpoint: d.TokenEndpoint, JWKSEndpoint: d.JWKSEndpoint, ClientID: d.ServiceClientID, ClientAuthentication: "client_secret_basic", Audience: d.BackendAudience, Subject: d.ServiceSubject, TenantID: d.TenantID, Organizations: d.Organizations, Permissions: []string{"portal-session:manage"}, Scopes: d.ServiceScopes, AudienceParameter: d.ServiceAudienceParameter, MaximumLifetimeSeconds: 3600, RefreshBeforeSeconds: 10}
	raw, err := json.Marshal(document)
	if err != nil {
		return nil, ErrServiceTokenConfiguration
	}
	sum := sha256.Sum256(raw)
	profile, err := LoadServiceTokenProfile(raw, hex.EncodeToString(sum[:]))
	if err != nil {
		return nil, err
	}
	return NewServiceTokenBroker(ctx, profile, reader)
}
````

### FILE: `internal/platform/identity/portal_session.go`

```yaml
block_id: "GO-OIDC-PORTAL-SESSION:file13:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "168f2379d594770b0c2fbb5b0c67db5208a37a5010534d2b38aec3b7b46ce4a1"
variables: []
secrets_allowed: false
```

````go
package identity

import (
	"context"
	"errors"
	"time"
)

var ErrPortalSessionConflict = errors.New("portal session conflict")
var ErrPortalSessionNotFound = errors.New("portal session not found")

type PortalSessionCommand struct {
	SessionID         string    `json:"session_id"`
	OperationID       string    `json:"operation_id,omitempty"`
	Version           int64     `json:"version,omitempty"`
	Ciphertext        string    `json:"ciphertext,omitempty"`
	AccessExpiresAt   time.Time `json:"access_expires_at,omitempty"`
	AbsoluteExpiresAt time.Time `json:"absolute_expires_at,omitempty"`
}
type PortalSessionRecord struct {
	Version           int64     `json:"version"`
	State             string    `json:"state"`
	OperationID       string    `json:"operation_id"`
	Ciphertext        string    `json:"ciphertext"`
	AccessExpiresAt   time.Time `json:"access_expires_at"`
	AbsoluteExpiresAt time.Time `json:"absolute_expires_at"`
	RevocationPending bool      `json:"revocation_pending"`
}
type PortalSessionStore interface {
	Execute(context.Context, PortalProfile, string, PortalSessionCommand) (PortalSessionRecord, error)
	Sweep(context.Context, PortalProfile, string) (PortalSessionSweep, error)
}

type PortalSessionSweepItem struct {
	SessionIDSHA256 string `json:"session_id_sha256"`
	PortalSessionRecord
}
type PortalSessionSweep struct {
	Items             []PortalSessionSweepItem `json:"items"`
	Purged            int64                    `json:"purged"`
	UnconfirmedPurged int64                    `json:"unconfirmed_purged"`
}
````

### FILE: `internal/platform/postgres/portal_session.go`

```yaml
block_id: "GO-OIDC-PORTAL-SESSION:file14:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "f2656cd45e101aecb27fba2761015be83545275d6cc1766f11f69f3d5759858d"
variables: []
secrets_allowed: false
```

````go
package postgres

// AUTHORED PostgreSQL CAS adapter; OAuth algorithms remain in official SDKs.
import (
	"bytes"
	"context"
	"crypto/sha256"
	"elite.local/enterprise/internal/platform/identity"
	"encoding/base64"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"strings"
	"time"
)

type PortalSessions struct{ Pool *pgxpool.Pool }

func portalOpaque(s string) bool {
	raw, err := base64.RawURLEncoding.DecodeString(s)
	return err == nil && len(raw) == 32 && base64.RawURLEncoding.EncodeToString(raw) == s
}
func (s *PortalSessions) Execute(ctx context.Context, p identity.PortalProfile, action string, c identity.PortalSessionCommand) (identity.PortalSessionRecord, error) {
	var out identity.PortalSessionRecord
	bad := identity.ErrPortalSessionConflict
	if s == nil || s.Pool == nil || p.SHA256() == "" || !portalOpaque(c.SessionID) || len(c.Ciphertext) > 32768 {
		return out, bad
	}
	key := sha256.Sum256([]byte(c.SessionID))
	cipherHash := sha256.Sum256([]byte(c.Ciphertext))
	tx, err := s.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return out, err
	}
	defer tx.Rollback(ctx)
	var now time.Time
	if err = tx.QueryRow(ctx, `select clock_timestamp()`).Scan(&now); err != nil {
		return out, err
	}
	if action == "create" {
		if c.Version != 0 || c.OperationID != "" || len(c.Ciphertext) < 32 || strings.Count(c.Ciphertext, ".") != 4 || !c.AccessExpiresAt.After(now) || c.AccessExpiresAt.After(c.AbsoluteExpiresAt) || !c.AbsoluteExpiresAt.After(now) || c.AbsoluteExpiresAt.After(now.Add(time.Duration(p.Document().MaximumSessionSeconds)*time.Second)) {
			return out, bad
		}
		_, err = tx.Exec(ctx, `insert into platform.portal_session(profile_sha256,session_id_sha256,state,ciphertext,ciphertext_sha256,access_expires_at,absolute_expires_at) values($1,$2,'active',$3,$4,$5,$6) on conflict do nothing`, p.SHA256(), key[:], c.Ciphertext, cipherHash[:], c.AccessExpiresAt, c.AbsoluteExpiresAt)
		if err != nil {
			return out, err
		}
	}
	var storedHash []byte
	err = tx.QueryRow(ctx, `select version,state,operation_id,ciphertext,ciphertext_sha256,access_expires_at,absolute_expires_at,revocation_pending from platform.portal_session where profile_sha256=$1 and session_id_sha256=$2 for update`, p.SHA256(), key[:]).Scan(&out.Version, &out.State, &out.OperationID, &out.Ciphertext, &storedHash, &out.AccessExpiresAt, &out.AbsoluteExpiresAt, &out.RevocationPending)
	if errors.Is(err, pgx.ErrNoRows) {
		return out, identity.ErrPortalSessionNotFound
	}
	if err != nil {
		return out, err
	}
	currentHash := sha256.Sum256([]byte(out.Ciphertext))
	if !bytes.Equal(currentHash[:], storedHash) {
		return out, bad
	}
	// Re-read time after acquiring the row lock: waiting must not extend a lease/session.
	if err = tx.QueryRow(ctx, `select clock_timestamp()`).Scan(&now); err != nil {
		return out, err
	}
	if !out.AbsoluteExpiresAt.After(now) && out.State != "revoked" {
		_, err = tx.Exec(ctx, `update platform.portal_session set state='reauth_required',revocation_pending=true,updated_at=clock_timestamp() where profile_sha256=$1 and session_id_sha256=$2`, p.SHA256(), key[:])
		if err != nil {
			return out, err
		}
		out.State = "reauth_required"
		out.RevocationPending = true
	}
	switch action {
	case "create":
		if out.Version != 1 || out.State != "active" || out.Ciphertext != c.Ciphertext || !out.AccessExpiresAt.Equal(c.AccessExpiresAt) || !out.AbsoluteExpiresAt.Equal(c.AbsoluteExpiresAt) {
			return out, bad
		}
	case "read":
	case "claim":
		if !portalOpaque(c.OperationID) || c.Version != out.Version || !out.AbsoluteExpiresAt.After(now) {
			return out, bad
		}
		if out.State == "refreshing" && out.OperationID == c.OperationID {
			break
		}
		if out.State != "active" {
			return out, bad
		}
		out.State = "refreshing"
		out.OperationID = c.OperationID
	case "commit":
		if !portalOpaque(c.OperationID) || out.OperationID != c.OperationID || len(c.Ciphertext) < 32 || strings.Count(c.Ciphertext, ".") != 4 || !c.AccessExpiresAt.After(now) || c.AccessExpiresAt.After(out.AbsoluteExpiresAt) || !out.AbsoluteExpiresAt.After(now) {
			return out, bad
		}
		if out.Version == c.Version+1 && out.Ciphertext == c.Ciphertext && out.AccessExpiresAt.Equal(c.AccessExpiresAt) {
			break
		}
		if out.Version != c.Version || (out.State != "refreshing" && out.State != "revoked") {
			return out, bad
		}
		out.Version++
		out.Ciphertext = c.Ciphertext
		out.AccessExpiresAt = c.AccessExpiresAt
		if out.State == "refreshing" {
			out.State = "active"
		} else {
			out.RevocationPending = true
		}
	case "abort":
		if out.OperationID != c.OperationID || out.Version != c.Version {
			return out, bad
		}
		if out.State == "refreshing" {
			out.State = "reauth_required"
			out.RevocationPending = true
		}
	case "revoke":
		out.State = "revoked"
		out.RevocationPending = out.Ciphertext != ""
	case "ack-revocation":
		if out.State != "revoked" || out.Version != c.Version || out.OperationID != c.OperationID {
			return out, bad
		}
		out.RevocationPending = false
		out.Ciphertext = ""
	default:
		return out, bad
	}
	finalHash := sha256.Sum256([]byte(out.Ciphertext))
	_, err = tx.Exec(ctx, `update platform.portal_session set version=$3,state=$4,operation_id=$5,ciphertext=$6,ciphertext_sha256=$7,access_expires_at=$8,revocation_pending=$9,updated_at=clock_timestamp() where profile_sha256=$1 and session_id_sha256=$2`, p.SHA256(), key[:], out.Version, out.State, out.OperationID, out.Ciphertext, finalHash[:], out.AccessExpiresAt, out.RevocationPending)
	if err != nil {
		return out, err
	}
	if err = tx.Commit(ctx); err != nil {
		return identity.PortalSessionRecord{}, err
	}
	return out, nil
}
````

### FILE: `internal/platform/postgres/portal_session_integration_test.go`

```yaml
block_id: "GO-OIDC-PORTAL-SESSION:file15:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "88fa21d191f768c7362e204989dea75ac29f9ed2b271e0badb3906f58714a716"
variables: []
secrets_allowed: false
```

````go
package postgres_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/platform/identity"
	db "elite.local/enterprise/internal/platform/postgres"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	jose "github.com/go-jose/go-jose/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

type portalFixture struct {
	mu                                                      sync.Mutex
	issuer                                                  *httptest.Server
	api                                                     *httptest.Server
	keys                                                    []*rsa.PrivateKey
	keyIndex                                                int
	mode                                                    string
	codes                                                   map[string]url.Values
	refresh                                                 map[string]bool
	counter                                                 int
	grants, refreshGrants, serviceGrants, jwks, revocations int
	tokenAuthValid, tokenClientValid, tokenBodyGrantValid   bool
	pool                                                    *pgxpool.Pool
	profile                                                 identity.PortalProfile
}

func (f *portalFixture) jwt(claims map[string]any) string {
	raw, _ := json.Marshal(claims)
	signer, _ := jose.NewSigner(jose.SigningKey{Algorithm: jose.RS256, Key: f.keys[f.keyIndex]}, (&jose.SignerOptions{}).WithHeader("kid", fmt.Sprintf("fixture-%d", f.keyIndex)))
	signed, _ := signer.Sign(raw)
	value, _ := signed.CompactSerialize()
	return value
}
func (f *portalFixture) issue(w http.ResponseWriter, nonce string, service bool) {
	now := time.Now()
	f.counter++
	claims := map[string]any{"iss": f.issuer.URL, "sub": "person", "aud": "portal-client", "iat": now.Unix(), "exp": now.Add(120 * time.Second).Unix(), "tenant_id": "tenant", "organization_ids": []string{"org"}, "permissions": []string{"customer:read"}}
	if nonce != "" {
		claims["nonce"] = nonce
	}
	if service {
		claims["sub"] = "portal-service"
		claims["aud"] = "api"
		claims["permissions"] = []string{"portal-session:manage"}
		claims["scope"] = "portal-session:manage"
	}
	if !service {
		switch f.mode {
		case "permission_withdrawn":
			claims["permissions"] = []string{}
		case "subject_changed":
			claims["sub"] = "different-person"
		case "wrong_issuer":
			claims["iss"] = "https://foreign.invalid"
		case "wrong_audience":
			claims["aud"] = "foreign"
		case "expired":
			claims["exp"] = now.Add(-time.Hour).Unix()
		case "wrong_nonce":
			claims["nonce"] = "foreign"
		case "foreign_tenant":
			claims["tenant_id"] = "foreign"
		case "extra_permission":
			claims["permissions"] = []string{"admin:write"}
		}
	}
	id := f.jwt(claims)
	claims["aud"] = "api"
	access := f.jwt(claims)
	response := map[string]any{"access_token": access, "token_type": "Bearer", "expires_in": 120}
	if !service {
		refresh := fmt.Sprintf("fixture-refresh-%d", f.counter)
		f.refresh[refresh] = true
		response["id_token"] = id
		response["refresh_token"] = refresh
		response["scope"] = "openid profile offline_access"
	}
	if f.mode == "grant_lost" && !service {
		_, _ = w.Write([]byte(`{"truncated":`))
		f.mode = ""
		return
	}
	_ = json.NewEncoder(w).Encode(response)
}
func (f *portalFixture) serveIssuer(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	switch r.URL.Path {
	case "/.well-known/openid-configuration":
		_ = json.NewEncoder(w).Encode(map[string]any{"issuer": f.issuer.URL, "authorization_endpoint": f.issuer.URL + "/authorize", "token_endpoint": f.issuer.URL + "/token", "jwks_uri": f.issuer.URL + "/jwks", "revocation_endpoint": f.issuer.URL + "/revoke", "end_session_endpoint": f.issuer.URL + "/logout", "response_types_supported": []string{"code"}, "subject_types_supported": []string{"public"}, "id_token_signing_alg_values_supported": []string{"RS256"}, "code_challenge_methods_supported": []string{"S256"}, "token_endpoint_auth_methods_supported": []string{"client_secret_basic"}, "grant_types_supported": []string{"authorization_code", "refresh_token", "client_credentials"}})
	case "/jwks":
		f.jwks++
		_ = json.NewEncoder(w).Encode(jose.JSONWebKeySet{Keys: []jose.JSONWebKey{{Key: &f.keys[f.keyIndex].PublicKey, KeyID: fmt.Sprintf("fixture-%d", f.keyIndex), Algorithm: "RS256", Use: "sig"}}})
	case "/authorize":
		q := r.URL.Query()
		if q.Get("client_id") != "portal-client" || q.Get("code_challenge_method") != "S256" || q.Get("code_challenge") == "" || q.Get("state") == "" || q.Get("nonce") == "" || q.Get("redirect_uri") != f.profile.Document().CallbackURL {
			w.WriteHeader(400)
			return
		}
		f.counter++
		code := fmt.Sprintf("fixture-code-%d", f.counter)
		f.codes[code] = q
		callback, _ := url.Parse(q.Get("redirect_uri"))
		query := callback.Query()
		query.Set("code", code)
		query.Set("state", q.Get("state"))
		callback.RawQuery = query.Encode()
		http.Redirect(w, r, callback.String(), 303)
	case "/token":
		_ = r.ParseForm()
		client, secret, ok := r.BasicAuth()
		client, _ = url.QueryUnescape(client)
		secret, _ = url.QueryUnescape(secret)
		f.tokenAuthValid = ok && secret == "synthetic-oidc-only"
		f.tokenClientValid = client == "portal-client" || client == "portal-service-client"
		f.tokenBodyGrantValid = r.Form.Get("grant_type") == "authorization_code" || r.Form.Get("grant_type") == "refresh_token" || r.Form.Get("grant_type") == "client_credentials"
		if !ok || secret != "synthetic-oidc-only" {
			w.WriteHeader(401)
			return
		}
		if f.mode == "account_disabled" && client == "portal-client" {
			w.WriteHeader(400)
			_, _ = w.Write([]byte(`{"error":"invalid_grant"}`))
			return
		}
		switch r.Form.Get("grant_type") {
		case "client_credentials":
			if client != "portal-service-client" || r.Form.Get("scope") != "portal-session:manage" || r.Form.Get("audience") != "api" {
				w.WriteHeader(400)
				return
			}
			f.serviceGrants++
			f.issue(w, "", true)
		case "authorization_code":
			q, exists := f.codes[r.Form.Get("code")]
			sum := sha256.Sum256([]byte(r.Form.Get("code_verifier")))
			if client != "portal-client" || !exists || base64.RawURLEncoding.EncodeToString(sum[:]) != q.Get("code_challenge") || r.Form.Get("redirect_uri") != q.Get("redirect_uri") {
				w.WriteHeader(400)
				_, _ = w.Write([]byte(`{"error":"invalid_grant"}`))
				return
			}
			delete(f.codes, r.Form.Get("code"))
			f.grants++
			f.issue(w, q.Get("nonce"), false)
		case "refresh_token":
			if client != "portal-client" || !f.refresh[r.Form.Get("refresh_token")] {
				w.WriteHeader(400)
				_, _ = w.Write([]byte(`{"error":"invalid_grant"}`))
				return
			}
			f.refresh[r.Form.Get("refresh_token")] = false
			f.refreshGrants++
			f.issue(w, "", false)
		default:
			w.WriteHeader(400)
		}
	case "/revoke":
		_ = r.ParseForm()
		client, secret, ok := r.BasicAuth()
		client, _ = url.QueryUnescape(client)
		secret, _ = url.QueryUnescape(secret)
		if !ok || client != "portal-client" || secret != "synthetic-oidc-only" {
			w.WriteHeader(401)
			return
		}
		if f.mode == "revoke_unavailable" {
			w.WriteHeader(503)
			_, _ = w.Write([]byte(`{"error":"temporarily_unavailable"}`))
			return
		}
		f.revocations++
		if r.Form.Get("token_type_hint") == "refresh_token" {
			f.refresh[r.Form.Get("token")] = false
		}
		_, _ = w.Write([]byte(`{}`))
	case "/logout":
		http.Redirect(w, r, f.profile.Document().PostLogoutURL, 303)
	default:
		w.WriteHeader(404)
	}
}
func newPortalFixture(t *testing.T, pool *pgxpool.Pool) *portalFixture {
	f := &portalFixture{pool: pool, codes: map[string]url.Values{}, refresh: map[string]bool{}}
	for range 2 {
		key, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			t.Fatal(err)
		}
		f.keys = append(f.keys, key)
	}
	f.issuer = httptest.NewServer(http.HandlerFunc(f.serveIssuer))
	t.Cleanup(f.issuer.Close)
	var handler http.Handler
	f.api = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/fixture/") {
			f.control(w, r)
			return
		}
		f.mu.Lock()
		lost := f.mode == "commit_lost" && strings.HasSuffix(r.URL.Path, "/commit")
		if lost {
			f.mode = ""
		}
		f.mu.Unlock()
		if lost {
			record := httptest.NewRecorder()
			handler.ServeHTTP(record, r)
			if record.Code != 200 {
				t.Error("commit fixture failed", record.Code)
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"truncated":`))
			return
		}
		handler.ServeHTTP(w, r)
	}))
	t.Cleanup(f.api.Close)
	d := identity.PortalProfileDocument{Schema: "elite.oidc.portal-lifecycle.v1", ProfileID: "fixture-portal", Revision: 1, Transport: "LOOPBACK_FIXTURE", Issuer: f.issuer.URL, AuthorizationEndpoint: f.issuer.URL + "/authorize", TokenEndpoint: f.issuer.URL + "/token", JWKSEndpoint: f.issuer.URL + "/jwks", RevocationEndpoint: f.issuer.URL + "/revoke", EndSessionEndpoint: f.issuer.URL + "/logout", ClientID: "portal-client", CallbackURL: "http://127.0.0.1:4567/api/auth/callback", PostLogoutURL: "http://127.0.0.1:4567/", BridgeURL: f.api.URL, BackendAudience: "api", ServiceClientID: "portal-service-client", ServiceSubject: "portal-service", ServiceScopes: []string{"portal-session:manage"}, ServiceAudienceParameter: "api", TenantID: "tenant", Organizations: []string{"org"}, AllowedPermissions: []string{"customer:read"}, Scopes: []string{"openid", "profile", "offline_access"}, MaximumSessionSeconds: 3600, RefreshBeforeSeconds: 15, RetentionSeconds: 86400}
	raw, _ := json.Marshal(d)
	sum := sha256.Sum256(raw)
	p, err := identity.LoadPortalProfile(raw, hex.EncodeToString(sum[:]))
	if err != nil {
		t.Fatal(err)
	}
	f.profile = p
	verifier, err := identity.NewOIDCVerifier(context.Background(), f.issuer.URL, "api")
	if err != nil {
		t.Fatal(err)
	}
	mux := http.NewServeMux()
	httpapi.PortalSessionModule{Profile: p, Store: &db.PortalSessions{Pool: pool}}.Register(mux, verifier)
	handler = mux
	return f
}
func (f *portalFixture) control(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	if r.Method == "POST" && r.URL.Path == "/fixture/control" {
		var input struct {
			Mode           string `json:"mode"`
			Expire         bool   `json:"expire"`
			Rotate         bool   `json:"rotate"`
			Swap           bool   `json:"swap"`
			ExpireAbsolute bool   `json:"expire_absolute"`
			Retention      bool   `json:"retention"`
		}
		if json.NewDecoder(r.Body).Decode(&input) != nil {
			w.WriteHeader(400)
			return
		}
		f.mode = input.Mode
		if input.Rotate {
			f.keyIndex = 1
		}
		if input.Expire {
			_, err := f.pool.Exec(r.Context(), `update platform.portal_session set access_expires_at=clock_timestamp()-interval '1 second' where profile_sha256=$1 and state='active'`, f.profile.SHA256())
			if err != nil {
				w.WriteHeader(500)
				return
			}
		}
		if input.Swap {
			tx, err := f.pool.Begin(r.Context())
			if err != nil {
				w.WriteHeader(500)
				return
			}
			defer tx.Rollback(r.Context())
			rows, err := tx.Query(r.Context(), `select session_id_sha256,ciphertext,ciphertext_sha256 from platform.portal_session where profile_sha256=$1 order by created_at desc limit 2 for update`, f.profile.SHA256())
			if err != nil {
				w.WriteHeader(500)
				return
			}
			var ids, hashes [][]byte
			var payloads []string
			for rows.Next() {
				var id, hash []byte
				var body string
				if rows.Scan(&id, &body, &hash) != nil {
					w.WriteHeader(500)
					return
				}
				ids = append(ids, id)
				hashes = append(hashes, hash)
				payloads = append(payloads, body)
			}
			rows.Close()
			if len(ids) != 2 {
				w.WriteHeader(500)
				return
			}
			for i := 0; i < 2; i++ {
				_, err = tx.Exec(r.Context(), `update platform.portal_session set ciphertext=$3,ciphertext_sha256=$4 where profile_sha256=$1 and session_id_sha256=$2`, f.profile.SHA256(), ids[i], payloads[1-i], hashes[1-i])
				if err != nil {
					w.WriteHeader(500)
					return
				}
			}
			if tx.Commit(r.Context()) != nil {
				w.WriteHeader(500)
				return
			}
		}
		if input.ExpireAbsolute || input.Retention {
			interval := "1 second"
			if input.Retention {
				interval = "2 days"
			}
			_, err := f.pool.Exec(r.Context(), `update platform.portal_session set absolute_expires_at=clock_timestamp()-$2::interval,access_expires_at=least(access_expires_at,clock_timestamp()-$2::interval) where profile_sha256=$1`, f.profile.SHA256(), interval)
			if err != nil {
				w.WriteHeader(500)
				return
			}
		}
		_, _ = w.Write([]byte(`{}`))
		return
	}
	if r.URL.Path == "/fixture/stats" {
		var active, revoked, reauth int
		_ = f.pool.QueryRow(r.Context(), `select count(*) filter(where state='active'),count(*) filter(where state='revoked'),count(*) filter(where state='reauth_required') from platform.portal_session where profile_sha256=$1`, f.profile.SHA256()).Scan(&active, &revoked, &reauth)
		_ = json.NewEncoder(w).Encode(map[string]any{"code_grants": f.grants, "refresh_grants": f.refreshGrants, "service_grants": f.serviceGrants, "jwks": f.jwks, "revocations": f.revocations, "active": active, "revoked": revoked, "reauth_required": reauth, "token_auth_valid": f.tokenAuthValid, "token_client_valid": f.tokenClientValid, "token_body_grant_valid": f.tokenBodyGrantValid})
		return
	}
	w.WriteHeader(404)
}
func TestPortalLifecycleOfficialSDKPostgres(t *testing.T) {
	dbURL := os.Getenv("PORTAL_SESSION_DB_URL")
	if dbURL == "" {
		t.Skip("explicit disposable portal PG fixture required")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	f := newPortalFixture(t, pool)
	directory := t.TempDir()
	doc, _ := json.Marshal(f.profile.Document())
	profile := filepath.Join(directory, "profile.json")
	if err = os.WriteFile(profile, doc, 0600); err != nil {
		t.Fatal(err)
	}
	secret := filepath.Join(directory, "client-secret")
	key := filepath.Join(directory, "session-key")
	if err = os.WriteFile(secret, []byte("synthetic-oidc-only"), 0600); err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(key, []byte("synthetic-portal-key-not-for-production-000000"), 0600); err != nil {
		t.Fatal(err)
	}
	node := os.Getenv("ELITE_NODE_EXACT")
	if node == "" {
		t.Fatal("pinned Node path required")
	}
	root, err := filepath.Abs("../../..")
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.CommandContext(ctx, node, filepath.Join(root, "node_modules/vitest/vitest.mjs"), "run", "src/platform/auth/portal-lifecycle.connected.test.ts", "--reporter=verbose")
	cmd.Dir = root
	cmd.Env = append(os.Environ(), "OIDC_PORTAL_LIFECYCLE_ENABLED=true", "OIDC_PORTAL_PROFILE_FILE="+profile, "OIDC_PORTAL_PROFILE_SHA256="+f.profile.SHA256(), "OIDC_PORTAL_CLIENT_SECRET_FILE="+secret, "OIDC_PORTAL_SERVICE_CLIENT_SECRET_FILE="+secret, "OIDC_PORTAL_SESSION_KEY_FILE="+key, "PORTAL_FIXTURE_CONTROL_URL="+f.api.URL, "APP_BASE_URL=http://127.0.0.1:4567")
	output, err := cmd.CombinedOutput()
	t.Log(string(output))
	if err != nil {
		t.Fatal("real SDK BFF/PG fixture failed", err)
	}
	var plain int
	if err = pool.QueryRow(ctx, `select count(*) from platform.portal_session where profile_sha256=$1 and (ciphertext like '%fixture-refresh%' or ciphertext like '%accessToken%')`, f.profile.SHA256()).Scan(&plain); err != nil || plain != 0 {
		t.Fatal("plaintext token persisted", err)
	}
	t.Log("PORTAL_LIFECYCLE_SDK_PG_PASS actual_code_pkce=true actual_refresh_rotation=true durable_cas=true logout_replay_rejected=true plaintext_tokens=0")
}

func TestPortalSessionCASAndScope(t *testing.T) {
	dbURL := os.Getenv("PORTAL_SESSION_DB_URL")
	if dbURL == "" {
		t.Skip("explicit disposable portal PG fixture required")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	f := newPortalFixture(t, pool)
	store := &db.PortalSessions{Pool: pool}
	id := func() string {
		bytes := make([]byte, 32)
		_, _ = rand.Read(bytes)
		return base64.RawURLEncoding.EncodeToString(bytes)
	}
	sid := id()
	now := time.Now().UTC().Truncate(time.Millisecond)
	cipher := "header..iv.ciphertext-payload-fixture.tag"
	create := identity.PortalSessionCommand{SessionID: sid, Ciphertext: cipher, AccessExpiresAt: now.Add(time.Minute), AbsoluteExpiresAt: now.Add(time.Hour - time.Second)}
	row, err := store.Execute(ctx, f.profile, "create", create)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = store.Execute(ctx, f.profile, "create", create); err != nil {
		t.Fatal("create replay", err)
	}
	var wg sync.WaitGroup
	wins := make(chan string, 16)
	for range 16 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			op := id()
			if _, err := store.Execute(ctx, f.profile, "claim", identity.PortalSessionCommand{SessionID: sid, Version: row.Version, OperationID: op}); err == nil {
				wins <- op
			}
		}()
	}
	wg.Wait()
	close(wins)
	winner := ""
	count := 0
	for op := range wins {
		winner = op
		count++
	}
	if count != 1 {
		t.Fatal("multiple refresh owners", count)
	}
	if _, err = store.Execute(ctx, f.profile, "commit", identity.PortalSessionCommand{SessionID: sid, Version: 1, OperationID: id(), Ciphertext: cipher, AccessExpiresAt: now.Add(2 * time.Minute)}); err == nil {
		t.Fatal("foreign operation committed")
	}
	if _, err = store.Execute(ctx, f.profile, "revoke", identity.PortalSessionCommand{SessionID: sid}); err != nil {
		t.Fatal(err)
	}
	changed := cipher + "-rotated"
	commit := identity.PortalSessionCommand{SessionID: sid, Version: 1, OperationID: winner, Ciphertext: changed, AccessExpiresAt: now.Add(2 * time.Minute)}
	row, err = store.Execute(ctx, f.profile, "commit", commit)
	if err != nil || row.State != "revoked" || !row.RevocationPending || row.Version != 2 {
		t.Fatal("refresh reactivated logged out session", row.State, err)
	}
	if _, err = store.Execute(ctx, f.profile, "ack-revocation", identity.PortalSessionCommand{SessionID: sid, Version: 1, OperationID: winner}); err == nil {
		t.Fatal("stale revocation acknowledgement deleted new token")
	}
	if _, err = store.Execute(ctx, f.profile, "ack-revocation", identity.PortalSessionCommand{SessionID: sid, Version: 2, OperationID: winner}); err != nil {
		t.Fatal(err)
	}
	row, err = store.Execute(ctx, f.profile, "read", identity.PortalSessionCommand{SessionID: sid})
	if err != nil || row.Ciphertext != "" || row.State != "revoked" {
		t.Fatal("revocation cleanup", err)
	}
	d := f.profile.Document()
	d.ProfileID = "foreign-profile"
	raw, _ := json.Marshal(d)
	sum := sha256.Sum256(raw)
	foreign, err := identity.LoadPortalProfile(raw, hex.EncodeToString(sum[:]))
	if err != nil {
		t.Fatal(err)
	}
	if _, err = store.Execute(ctx, foreign, "read", identity.PortalSessionCommand{SessionID: sid}); err != identity.ErrPortalSessionNotFound {
		t.Fatal("profile isolation", err)
	}
	t.Log("PORTAL_CAS_PASS concurrent_claims=16 grant_owner=1 logout_wins=true stale_ack_rejected=true profile_isolation=true")

	validClaims := func() map[string]any {
		return map[string]any{"iss": f.issuer.URL, "sub": "portal-service", "aud": "api", "iat": time.Now().Unix(), "exp": time.Now().Add(time.Minute).Unix(), "tenant_id": "tenant", "organization_ids": []string{"org"}, "permissions": []string{"portal-session:manage"}}
	}
	for _, tc := range []struct {
		name   string
		field  string
		value  any
		status int
	}{{"user", "sub", "person", 403}, {"tenant", "tenant_id", "foreign", 403}, {"organization", "organization_ids", []string{"foreign"}, 403}, {"wildcard", "permissions", []string{"*"}, 403}, {"extra_permission", "permissions", []string{"portal-session:manage", "admin:write"}, 403}, {"issuer", "iss", "https://foreign.invalid", 401}, {"audience", "aud", "foreign", 401}, {"expiry", "exp", time.Now().Add(-time.Minute).Unix(), 401}} {
		t.Run("http_"+tc.name, func(t *testing.T) {
			claims := validClaims()
			claims[tc.field] = tc.value
			request, _ := http.NewRequest("POST", f.api.URL+"/v1/private/portal-sessions/read", strings.NewReader(`{"session_id":"`+sid+`"}`))
			request.Header.Set("Authorization", "Bearer "+f.jwt(claims))
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("X-Portal-Profile-SHA256", f.profile.SHA256())
			response, err := http.DefaultClient.Do(request)
			if err != nil {
				t.Fatal(err)
			}
			defer response.Body.Close()
			if response.StatusCode != tc.status {
				t.Fatal("authorization negative", response.StatusCode)
			}
		})
	}
	for _, tc := range []struct {
		name, body, profile string
		status              int
	}{{"profile", `{"session_id":"` + sid + `"}`, strings.Repeat("0", 64), 403}, {"duplicate", `{"session_id":"` + sid + `","session_id":"` + sid + `"}`, f.profile.SHA256(), 400}, {"unknown", `{"session_id":"` + sid + `","tenant_id":"foreign"}`, f.profile.SHA256(), 400}, {"too_large", strings.Repeat("x", 40001), f.profile.SHA256(), 413}} {
		t.Run("http_"+tc.name, func(t *testing.T) {
			request, _ := http.NewRequest("POST", f.api.URL+"/v1/private/portal-sessions/read", strings.NewReader(tc.body))
			request.Header.Set("Authorization", "Bearer "+f.jwt(validClaims()))
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("X-Portal-Profile-SHA256", tc.profile)
			response, err := http.DefaultClient.Do(request)
			if err != nil {
				t.Fatal(err)
			}
			defer response.Body.Close()
			_, _ = io.Copy(io.Discard, response.Body)
			if response.StatusCode != tc.status {
				t.Fatal("transport negative", response.StatusCode)
			}
		})
	}
	// Expiry while waiting on a row lock cannot be evaluated using pre-wait time.
	shortID := id()
	at := time.Now().UTC().Truncate(time.Millisecond)
	_, err = store.Execute(ctx, f.profile, "create", identity.PortalSessionCommand{SessionID: shortID, Ciphertext: cipher, AccessExpiresAt: at.Add(200 * time.Millisecond), AbsoluteExpiresAt: at.Add(350 * time.Millisecond)})
	if err != nil {
		t.Fatal(err)
	}
	lock, err := pool.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	shortHash := sha256.Sum256([]byte(shortID))
	if _, err = lock.Exec(ctx, `select 1 from platform.portal_session where profile_sha256=$1 and session_id_sha256=$2 for update`, f.profile.SHA256(), shortHash[:]); err != nil {
		t.Fatal(err)
	}
	result := make(chan error, 1)
	go func() {
		_, e := store.Execute(ctx, f.profile, "claim", identity.PortalSessionCommand{SessionID: shortID, Version: 1, OperationID: id()})
		result <- e
	}()
	time.Sleep(450 * time.Millisecond)
	if err = lock.Commit(ctx); err != nil {
		t.Fatal(err)
	}
	if err = <-result; err != identity.ErrPortalSessionConflict {
		t.Fatal("expired lock-wait accepted", err)
	}
}
````

### FILE: `internal/platform/postgres/portal_session_sweep.go`

```yaml
block_id: "GO-OIDC-PORTAL-SESSION:file16:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "9d01c482a814551bb19f60f463b2b2b885eb614a4139c63945399f536462b0ea"
variables: []
secrets_allowed: false
```

````go
package postgres

import (
	"bytes"
	"context"
	"crypto/sha256"
	"elite.local/enterprise/internal/platform/identity"
	"github.com/jackc/pgx/v5"
)

// Sweep leases a bounded batch without needing a browser cookie. Retention
// counts unconfirmed provider revocation durably; it never labels it success.
func (s *PortalSessions) Sweep(ctx context.Context, p identity.PortalProfile, operation string) (identity.PortalSessionSweep, error) {
	result := identity.PortalSessionSweep{Items: []identity.PortalSessionSweepItem{}}
	if p.SHA256() == "" || !portalOpaque(operation) || s == nil || s.Pool == nil {
		return result, identity.ErrPortalSessionConflict
	}
	tx, err := s.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return result, err
	}
	defer tx.Rollback(ctx)
	err = tx.QueryRow(ctx, `with picked as (select profile_sha256,session_id_sha256 from platform.portal_session where profile_sha256=$1 and absolute_expires_at<clock_timestamp()-make_interval(secs=>$2) order by absolute_expires_at for update skip locked limit 100), removed as (delete from platform.portal_session s using picked p where s.profile_sha256=p.profile_sha256 and s.session_id_sha256=p.session_id_sha256 returning s.revocation_pending,s.state) select count(*),count(*) filter(where revocation_pending or state<>'revoked') from removed`, p.SHA256(), p.Document().RetentionSeconds).Scan(&result.Purged, &result.UnconfirmedPurged)
	if err != nil {
		return result, err
	}
	if result.Purged > 0 {
		_, err = tx.Exec(ctx, `insert into platform.portal_session_retention_summary(profile_sha256,purged_count,unconfirmed_provider_revocation_count) values($1,$2,$3) on conflict(profile_sha256) do update set purged_count=platform.portal_session_retention_summary.purged_count+excluded.purged_count,unconfirmed_provider_revocation_count=platform.portal_session_retention_summary.unconfirmed_provider_revocation_count+excluded.unconfirmed_provider_revocation_count,updated_at=clock_timestamp()`, p.SHA256(), result.Purged, result.UnconfirmedPurged)
		if err != nil {
			return result, err
		}
	}
	rows, err := tx.Query(ctx, `with picked as (select profile_sha256,session_id_sha256 from platform.portal_session where profile_sha256=$1 and ciphertext<>'' and revocation_lease_until<clock_timestamp() and (state in ('reauth_required','revoked') or absolute_expires_at<=clock_timestamp() or (state='refreshing' and updated_at<clock_timestamp()-interval '120 seconds')) order by updated_at for update skip locked limit 2) update platform.portal_session s set state='revoked',revocation_pending=true,revocation_lease_until=clock_timestamp()+interval '120 seconds',revocation_operation_id=$2,updated_at=clock_timestamp() from picked p where s.profile_sha256=p.profile_sha256 and s.session_id_sha256=p.session_id_sha256 returning encode(s.session_id_sha256,'hex'),s.version,s.state,s.operation_id,s.ciphertext,s.ciphertext_sha256,s.access_expires_at,s.absolute_expires_at,s.revocation_pending`, p.SHA256(), operation)
	if err != nil {
		return result, err
	}
	for rows.Next() {
		var item identity.PortalSessionSweepItem
		var hash []byte
		if err = rows.Scan(&item.SessionIDSHA256, &item.Version, &item.State, &item.OperationID, &item.Ciphertext, &hash, &item.AccessExpiresAt, &item.AbsoluteExpiresAt, &item.RevocationPending); err != nil {
			rows.Close()
			return result, err
		}
		current := sha256.Sum256([]byte(item.Ciphertext))
		if !bytes.Equal(current[:], hash) {
			rows.Close()
			return result, identity.ErrPortalSessionConflict
		}
		result.Items = append(result.Items, item)
	}
	err = rows.Err()
	rows.Close()
	if err != nil {
		return result, err
	}
	if err = tx.Commit(ctx); err != nil {
		return result, err
	}
	return result, nil
}
````

## 6. Configuration surface

docs/PORTAL_IDENTITY_LIFECYCLE.md and config/identity/portal-lifecycle.example.json. Future external credential file paths only; production requires TLS.

## 7. Dependency bill

No dependency/runtime upgrade. Original openid-client6.8.5, jose6.2.10, zod4.4.3, Go OIDC/service broker and pgx pins remain unchanged; all added files are AUTHORED composition glue.

## 8. Apply order

Apply selected migrations by name, including0062. Enable BFF and Go with exact same profile hash; optional Go host hook preserves independent profiles. BFF disabled mode preserves prior nonrefresh login.

## 9. Verification

Real SDK/signed fixture/Go bridge/PG7lifecycle cases and16CAS contenders plus12scope negatives;18TS regressions, host real service grant/JWKS/poll/shutdown, bounded responses, populated rollback refusal/emptydownup, Go profile/finite2sfuzz/vet/build.

## 10. Reconstruction evidence

reconstruction_evidence/IDENTITY_PORTAL_RELEASE_V402.md/json. Closes this session lifecycle claim within T2803; J5 administration and whole-composition source/SCA remain separate required gates.


### FILE: `docs/J5_IDP_ACCESS_CONTRACT.md`

```yaml
block_id: "J5-EXTENSION-GO_OIDC_PORTAL_SESSION:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "7b66a1b133df8e53d6cccd7c43ffdf2d9ef66b89d1b9654728b40e18c950d651"
variables: []
secrets_allowed: false
```

````markdown
# J5 identity administration contract — library reference

The selected integration is standards-based OIDC with an external IdP as the user, administrator, MFA and signed-role authority. No Keycloak/Entra/other proprietary admin API is selected or claimed. Creating IdP accounts, assigning roles and revoking them are user-side IdP administration; the complete code interface consumed here is code/PKCE, signed claims, refresh, service grant and revocation. This is not a locally authored password or identity-provider stack.

Bootstrap is explicit: a signed principal requires both network:admin and network:bootstrap to create an initial root organization. No global wildcard is needed or added. Ordinary administration additionally requires the exact organization ID in organization_ids. The initial administrator loses network:bootstrap after setup through the IdP's own role administration and authenticates again/refreshes. New organizations do not automatically add organization access to any principal. A provider administrator must explicitly grant the new scope. The durable original network command receipt records the authenticated subject, tenant, command hash and effect; another actor cannot recover it as their own.

For the future target, register the exact issuer/audience/client/redirects and bounded token lifetimes from PORTAL_IDENTITY_LIFECYCLE.md. Populate only the allowlisted permissions and organization IDs required by each role; no automatic default admin or first-login privilege. Keep the service principal separate with exactly portal-session:manage. IdP MFA, break-glass policy and its audit retention are account configuration/acceptance responsibilities, never inferred from library fixtures. Replacing a user or linking a different subject requires a new authorized identity flow: refresh cannot switch subjects in place.

Access review and deprovisioning: remove the permission/scope or disable the subject at the IdP, revoke its sessions there, and log out the local portal. On refresh, signed reduced permissions replace the old set; invalid_grant places the durable session into required reauthentication. Already issued access tokens are accepted only until their signed expiry under the selected issuer policy. The library does not claim immediate universal bearer revocation or a vendor admin API. Targets requiring immediate revocation must select and admit the corresponding IdP/introspection policy before live acceptance. No accounts or secrets are required for these local fixtures.

Local evidence: actual official SDK refresh and PostgreSQL enforce permission withdrawal, changed-subject rejection and disabled-account reauthentication. Signed RS256/JWKS tokens drive the original network HTTP/PG owner through narrow bootstrap, scoped operations and rejected post-withdrawal/expired/foreign-tenant/foreign-actor requests;2organizations/3immutable receipts/3outbox events. Earlier network agreement/branch and training proofs remain linked, without repeating their unchanged journeys. Library readiness/target go-live remains T2801/T2808; whole composition source/SCA remains T2803 until its own receipt is complete.
````

### FILE: `internal/platform/postgres/portal_access_review_integration_test.go`

```yaml
block_id: "J5-EXTENSION-GO_OIDC_PORTAL_SESSION:file2:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "98040fbbc6e6bad3d7382f06212b33a03e28a2804bb5ede47b6e890ee1042af8"
variables: []
secrets_allowed: false
```

````go
package postgres_test

import (
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestPortalAccessReviewOfficialSDKPostgres(t *testing.T) {
	dbURL := os.Getenv("PORTAL_SESSION_DB_URL")
	if dbURL == "" {
		t.Skip("explicit disposable portal PG fixture required")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	f := newPortalFixture(t, pool)
	directory := t.TempDir()
	doc, _ := json.Marshal(f.profile.Document())
	profile := filepath.Join(directory, "profile.json")
	if err = os.WriteFile(profile, doc, 0600); err != nil {
		t.Fatal(err)
	}
	secret := filepath.Join(directory, "client-secret")
	key := filepath.Join(directory, "session-key")
	if err = os.WriteFile(secret, []byte("synthetic-oidc-only"), 0600); err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(key, []byte("synthetic-portal-key-not-for-production-000000"), 0600); err != nil {
		t.Fatal(err)
	}
	node := os.Getenv("ELITE_NODE_EXACT")
	if node == "" {
		t.Fatal("pinned Node path required")
	}
	root, err := filepath.Abs("../../..")
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.CommandContext(ctx, node, filepath.Join(root, "node_modules/vitest/vitest.mjs"), "run", "src/platform/auth/portal-access-review.connected.test.ts", "--reporter=verbose")
	cmd.Dir = root
	cmd.Env = append(os.Environ(), "OIDC_PORTAL_LIFECYCLE_ENABLED=true", "OIDC_PORTAL_PROFILE_FILE="+profile, "OIDC_PORTAL_PROFILE_SHA256="+f.profile.SHA256(), "OIDC_PORTAL_CLIENT_SECRET_FILE="+secret, "OIDC_PORTAL_SERVICE_CLIENT_SECRET_FILE="+secret, "OIDC_PORTAL_SESSION_KEY_FILE="+key, "PORTAL_FIXTURE_CONTROL_URL="+f.api.URL, "APP_BASE_URL=http://127.0.0.1:4567")
	output, err := cmd.CombinedOutput()
	t.Log(string(output))
	if err != nil {
		t.Fatal("real SDK BFF/PG fixture failed", err)
	}
	var plain int
	if err = pool.QueryRow(ctx, `select count(*) from platform.portal_session where profile_sha256=$1 and (ciphertext like '%fixture-refresh%' or ciphertext like '%accessToken%')`, f.profile.SHA256()).Scan(&plain); err != nil || plain != 0 {
		t.Fatal("plaintext token persisted", err)
	}
	t.Log("PORTAL_ACCESS_REVIEW_PASS official_sdk_refresh=true permission_withdrawal=true changed_subject_rejected=true account_disabled_reauth=true plaintext_tokens=0")
}
````


V402 composed delta: Explicit network:admin + network:bootstrap, original scoped ongoing access; IDENTITY_J5_RELEASE_V402.md/json. IdP owns role grants and deprovisioning.
