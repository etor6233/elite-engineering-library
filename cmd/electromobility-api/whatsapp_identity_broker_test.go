package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"elite.local/enterprise/internal/platform/identity"
	jose "github.com/go-jose/go-jose/v4"
)

// Exercise the host selector and domain/worker identity adapter with the real
// pinned client-credentials and RS256/JWKS verifiers, not a token provider stub.
func TestWhatsAppHostUsesOIDCServiceBroker(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	var requests, grants atomic.Int32
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/.well-known/openid-configuration":
			_ = json.NewEncoder(w).Encode(map[string]any{"issuer": server.URL, "authorization_endpoint": server.URL + "/authorize", "token_endpoint": server.URL + "/token", "jwks_uri": server.URL + "/jwks", "response_types_supported": []string{"code"}, "subject_types_supported": []string{"public"}, "id_token_signing_alg_values_supported": []string{"RS256"}, "grant_types_supported": []string{"client_credentials"}, "token_endpoint_auth_methods_supported": []string{"client_secret_basic"}})
		case "/jwks":
			_ = json.NewEncoder(w).Encode(jose.JSONWebKeySet{Keys: []jose.JSONWebKey{{Key: &key.PublicKey, KeyID: "fixture-key", Use: "sig", Algorithm: "RS256"}}})
		case "/token":
			grants.Add(1)
			_ = r.ParseForm()
			id, secret, ok := r.BasicAuth()
			if r.Method != "POST" || !ok || id != "whatsapp-worker" || secret != "synthetic-secret" || r.Form.Get("grant_type") != "client_credentials" || r.Form.Get("scope") != "whatsapp:process appointment:manage" {
				t.Error("invalid grant contract")
				w.WriteHeader(400)
				return
			}
			now := time.Now()
			raw, _ := json.Marshal(map[string]any{"iss": server.URL, "sub": "fixture-worker", "aud": "enterprise-api", "iat": now.Unix(), "exp": now.Add(120 * time.Second).Unix(), "tenant_id": "tenant", "organization_ids": []string{"org"}, "permissions": []string{"whatsapp:process", "appointment:manage"}, "scope": "whatsapp:process appointment:manage", "client_id": "whatsapp-worker"})
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
	verifier, err := identity.NewOIDCVerifier(ctx, server.URL, "enterprise-api")
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	profileFile := filepath.Join(dir, "identity.json")
	secretFile := filepath.Join(dir, "secret")
	raw, _ := json.Marshal(map[string]any{"schema": "elite.oidc.service-token.v1", "profile_id": "wa-fixture", "revision": 1, "transport": "LOOPBACK_FIXTURE", "issuer": server.URL, "token_endpoint": server.URL + "/token", "jwks_endpoint": server.URL + "/jwks", "client_id": "whatsapp-worker", "client_authentication": "client_secret_basic", "audience": "enterprise-api", "subject": "fixture-worker", "tenant_id": "tenant", "organization_ids": []string{"org"}, "permissions": []string{"whatsapp:process", "appointment:manage"}, "scopes": []string{"whatsapp:process", "appointment:manage"}, "maximum_lifetime_seconds": 120, "refresh_before_seconds": 5})
	if os.WriteFile(profileFile, raw, 0600) != nil || os.WriteFile(secretFile, []byte("synthetic-secret\n"), 0600) != nil {
		t.Fatal("fixture files")
	}
	hash := sha256.Sum256(raw)
	env := map[string]string{"WHATSAPP_SERVICE_IDENTITY_PROFILE_FILE": profileFile, "WHATSAPP_SERVICE_IDENTITY_PROFILE_SHA256": hex.EncodeToString(hash[:]), "WHATSAPP_SERVICE_CLIENT_SECRET_FILE": secretFile}
	lookup := func(k string) string { return env[k] }
	cfg := whatsappHostConfig{TenantID: "tenant", OrganizationID: "org"}
	before := requests.Load()
	if _, e := selectWhatsAppIdentity(ctx, cfg, verifier, lookup); e == nil {
		t.Fatal("loopback profile accepted without reference opt-in")
	}
	if requests.Load() != before {
		t.Fatal("rejected fixture contacted issuer")
	}
	cfg.ReferenceFixture = true
	cfg.OrganizationID = "foreign"
	if _, e := selectWhatsAppIdentity(ctx, cfg, verifier, lookup); e == nil {
		t.Fatal("foreign organization accepted")
	}
	if requests.Load() != before {
		t.Fatal("scope rejection contacted issuer")
	}
	cfg.OrganizationID = "org"
	env["WHATSAPP_SERVICE_TOKEN_FILE"] = filepath.Join(dir, "unused")
	if _, e := selectWhatsAppIdentity(ctx, cfg, verifier, lookup); e == nil {
		t.Fatal("competing token owners accepted")
	}
	if requests.Load() != before {
		t.Fatal("mixed identity contacted issuer")
	}
	delete(env, "WHATSAPP_SERVICE_TOKEN_FILE")
	source, e := selectWhatsAppIdentity(ctx, cfg, verifier, lookup)
	if e != nil {
		t.Fatal(e)
	}
	var group sync.WaitGroup
	for i := 0; i < 12; i++ {
		group.Add(1)
		go func() {
			defer group.Done()
			if token, e := source.AccessToken(ctx); e != nil || token == "" {
				t.Error("domain bearer unavailable", e)
			}
			if principal, e := source.Resolve(ctx); e != nil || principal.Subject != "fixture-worker" {
				t.Error("worker principal unavailable", e)
			}
		}()
	}
	group.Wait()
	if grants.Load() != 1 {
		t.Fatalf("concurrent host calls performed %d grants", grants.Load())
	}
	cancelled, cancel := context.WithCancel(ctx)
	cancel()
	if _, e := source.AccessToken(cancelled); e == nil {
		t.Fatal("cancelled host got cached token")
	}
	if grants.Load() != 1 {
		t.Fatal("cancelled host performed another grant")
	}
}
