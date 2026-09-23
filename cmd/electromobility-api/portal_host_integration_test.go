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
