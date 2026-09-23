package identity

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	jose "github.com/go-jose/go-jose/v4"
)

type serviceIssuerFixture struct {
	server          *httptest.Server
	keys            []*rsa.PrivateKey
	keyIndex        int
	secret          string
	mutate          func(map[string]any, map[string]any)
	metadata        func(map[string]any)
	grants          atomic.Int32
	jwks            atomic.Int32
	invalidRequests atomic.Int32
	fail            string
	style           string
}

func newServiceIssuerFixture(t *testing.T) *serviceIssuerFixture {
	t.Helper()
	f := &serviceIssuerFixture{secret: "fixture-secret-never-live", style: "client_secret_basic"}
	for range 2 {
		key, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			t.Fatal(err)
		}
		f.keys = append(f.keys, key)
	}
	f.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/.well-known/openid-configuration":
			m := map[string]any{"issuer": f.server.URL, "authorization_endpoint": f.server.URL + "/authorize", "token_endpoint": f.server.URL + "/token", "jwks_uri": f.server.URL + "/jwks", "response_types_supported": []string{"code"}, "subject_types_supported": []string{"public"}, "id_token_signing_alg_values_supported": []string{"RS256"}, "token_endpoint_auth_methods_supported": []string{f.style}, "grant_types_supported": []string{"client_credentials"}}
			if f.metadata != nil {
				f.metadata(m)
			}
			_ = json.NewEncoder(w).Encode(m)
		case "/jwks":
			f.jwks.Add(1)
			_ = json.NewEncoder(w).Encode(jose.JSONWebKeySet{Keys: []jose.JSONWebKey{{Key: &f.keys[f.keyIndex].PublicKey, KeyID: []string{"fixture-a", "fixture-b"}[f.keyIndex], Use: "sig", Algorithm: "RS256"}}})
		case "/token":
			f.grants.Add(1)
			_ = r.ParseForm()
			id, secret, basic := r.BasicAuth()
			if f.style == "client_secret_post" {
				id = r.Form.Get("client_id")
				secret = r.Form.Get("client_secret")
				basic = r.Header.Get("Authorization") == ""
			}
			if r.Method != "POST" || !basic || id != "worker-client" || secret != f.secret || r.Form.Get("grant_type") != "client_credentials" || r.Form.Get("scope") != "message:send conversation:read" || r.Form.Get("audience") != "api" {
				f.invalidRequests.Add(1)
				w.WriteHeader(400)
				return
			}
			if f.fail == "redirect" {
				http.Redirect(w, r, f.server.URL+"/credential-exfiltration", 307)
				return
			}
			if f.fail == "large" {
				_, _ = w.Write([]byte(strings.Repeat("x", 262145)))
				return
			}
			if f.fail != "" {
				w.WriteHeader(400)
				_, _ = w.Write([]byte(`{"error":"invalid_client","error_description":"secret-canary-DO-NOT-LOG"}`))
				return
			}
			now := time.Now().UTC()
			claims := map[string]any{"iss": f.server.URL, "sub": "worker-subject", "aud": "api", "iat": now.Unix(), "exp": now.Add(120 * time.Second).Unix(), "tenant_id": "tenant", "organization_ids": []string{"org"}, "permissions": []string{"conversation:read", "message:send"}, "scope": "conversation:read message:send", "client_id": "worker-client"}
			result := map[string]any{"token_type": "Bearer", "expires_in": 120, "scope": "message:send conversation:read"}
			if f.mutate != nil {
				f.mutate(claims, result)
			}
			raw, _ := json.Marshal(claims)
			signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.RS256, Key: f.keys[f.keyIndex]}, (&jose.SignerOptions{}).WithHeader("kid", []string{"fixture-a", "fixture-b"}[f.keyIndex]))
			if err != nil {
				t.Error(err)
				return
			}
			signed, err := signer.Sign(raw)
			if err != nil {
				t.Error(err)
				return
			}
			result["access_token"], err = signed.CompactSerialize()
			if err != nil {
				t.Error(err)
				return
			}
			_ = json.NewEncoder(w).Encode(result)
		default:
			f.invalidRequests.Add(1)
			w.WriteHeader(404)
		}
	}))
	t.Cleanup(f.server.Close)
	return f
}
func (f *serviceIssuerFixture) document() serviceTokenDocument {
	return serviceTokenDocument{Schema: "elite.oidc.service-token.v1", ProfileID: "fixture-service", Revision: 1, Transport: "LOOPBACK_FIXTURE", Issuer: f.server.URL, TokenEndpoint: f.server.URL + "/token", JWKSEndpoint: f.server.URL + "/jwks", ClientID: "worker-client", ClientAuthentication: f.style, Audience: "api", Subject: "worker-subject", TenantID: "tenant", Organizations: []string{"org"}, Permissions: []string{"conversation:read", "message:send"}, Scopes: []string{"message:send", "conversation:read"}, AudienceParameter: "api", MaximumLifetimeSeconds: 180, RefreshBeforeSeconds: 10}
}
func serviceProfileForTest(t *testing.T, d serviceTokenDocument) ServiceTokenProfile {
	t.Helper()
	raw, _ := json.Marshal(d)
	sum := sha256.Sum256(raw)
	p, err := LoadServiceTokenProfile(raw, hex.EncodeToString(sum[:]))
	if err != nil {
		t.Fatal(err)
	}
	return p
}
func (f *serviceIssuerFixture) broker(t *testing.T) *ServiceTokenBroker {
	t.Helper()
	b, err := NewServiceTokenBroker(context.Background(), serviceProfileForTest(t, f.document()), func(context.Context) (string, error) { return f.secret, nil })
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func TestServiceTokenOfficialGrantCacheRenewalAndJWKSRotation(t *testing.T) {
	f := newServiceIssuerFixture(t)
	b := f.broker(t)
	var wg sync.WaitGroup
	tokens := make(chan string, 20)
	for range 20 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			token, err := b.AccessToken(context.Background())
			if err != nil {
				t.Error(err)
			}
			tokens <- token
		}()
	}
	wg.Wait()
	close(tokens)
	first := ""
	for token := range tokens {
		if first == "" {
			first = token
		}
		if token == "" || token != first {
			t.Fatal("different/missing cached token")
		}
	}
	if f.grants.Load() != 1 || f.jwks.Load() != 1 || f.invalidRequests.Load() != 0 {
		t.Fatal("grant/verification did not collapse", f.grants.Load(), f.jwks.Load())
	}
	b.token.Expiry = time.Now()
	f.keyIndex = 1
	second, err := b.AccessToken(context.Background())
	if err != nil || second == first || f.grants.Load() != 2 || f.jwks.Load() != 2 {
		t.Fatal("renewal/JWKS rotation failed", err)
	}
	principalVerifier, err := NewOIDCVerifier(context.Background(), f.server.URL, "api")
	if err != nil {
		t.Fatal(err)
	}
	principal, err := principalVerifier.Verify(context.Background(), second)
	if err != nil || principal.Subject != "worker-subject" || !principal.AllowedOrganization("org") {
		t.Fatal("backend verifier incompatible", err)
	}
	binding := b.Binding()
	binding.Permissions[0] = "*"
	if b.Binding().Permissions[0] == "*" {
		t.Fatal("profile mutable")
	}
	t.Log("SERVICE_OIDC_PASS sdk_client_credentials=true concurrent_callers=20 grants=2 jwks_fetches=2 rotation_verified=true backend_verifier=true")
}
func TestServiceTokenRejectsForeignClaimsAndInvalidLifetime(t *testing.T) {
	f := newServiceIssuerFixture(t)
	cases := []struct {
		name   string
		mutate func(map[string]any, map[string]any)
	}{
		{"issuer", func(c, r map[string]any) { c["iss"] = "https://foreign.invalid" }},
		{"audience", func(c, r map[string]any) { c["aud"] = "foreign" }},
		{"multiple_audiences", func(c, r map[string]any) { c["aud"] = []string{"api", "foreign"} }},
		{"expired", func(c, r map[string]any) { c["exp"] = time.Now().Add(-time.Second).Unix() }},
		{"not_before", func(c, r map[string]any) { c["nbf"] = time.Now().Add(time.Minute).Unix() }},
		{"issued_future", func(c, r map[string]any) { c["iat"] = time.Now().Add(time.Minute).Unix() }},
		{"issued_missing", func(c, r map[string]any) { delete(c, "iat") }},
		{"lifetime", func(c, r map[string]any) { c["exp"] = time.Now().Add(time.Hour).Unix() }},
		{"short_lifetime", func(c, r map[string]any) { c["exp"] = time.Now().Add(2 * time.Second).Unix() }},
		{"subject", func(c, r map[string]any) { c["sub"] = "foreign" }},
		{"tenant", func(c, r map[string]any) { c["tenant_id"] = "foreign" }},
		{"organization", func(c, r map[string]any) { c["organization_ids"] = []string{"foreign"} }},
		{"extra_permission", func(c, r map[string]any) {
			c["permissions"] = []string{"conversation:read", "message:send", "admin:write"}
		}},
		{"wildcard", func(c, r map[string]any) { c["permissions"] = []string{"*"} }},
		{"duplicate_scope", func(c, r map[string]any) { c["scope"] = "message:send message:send" }},
		{"missing_scope", func(c, r map[string]any) { delete(c, "scope") }},
		{"response_scope", func(c, r map[string]any) { r["scope"] = "admin" }},
		{"authorized_party", func(c, r map[string]any) { c["azp"] = "foreign" }},
		{"client_id", func(c, r map[string]any) { c["client_id"] = "foreign" }},
		{"refresh_token", func(c, r map[string]any) { r["refresh_token"] = "unsupported" }},
		{"token_type", func(c, r map[string]any) { r["token_type"] = "Other" }},
		{"missing_expiry", func(c, r map[string]any) { delete(r, "expires_in") }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f.mutate = tc.mutate
			b := f.broker(t)
			token, err := b.AccessToken(context.Background())
			if !errors.Is(err, ErrServiceTokenUnavailable) || token != "" || b.token != nil {
				t.Fatal("invalid token accepted", err)
			}
		})
	}
	f.mutate = func(c, r map[string]any) { delete(r, "scope") }
	b := f.broker(t)
	if _, err := b.AccessToken(context.Background()); err != nil {
		t.Fatal("RFC6749 unchanged response scope may be omitted", err)
	}
}
func TestServiceTokenFailureDoesNotLeakCacheOrCredentials(t *testing.T) {
	f := newServiceIssuerFixture(t)
	b := f.broker(t)
	if _, err := b.AccessToken(context.Background()); err != nil {
		t.Fatal(err)
	}
	for _, mode := range []string{"provider_error", "redirect", "large"} {
		t.Run(mode, func(t *testing.T) {
			if b.token != nil {
				b.token.Expiry = time.Now()
			}
			f.fail = mode
			token, err := b.AccessToken(context.Background())
			if token != "" || err != ErrServiceTokenUnavailable || strings.Contains(err.Error(), "canary") || b.token != nil {
				t.Fatal("unsafe failure", err)
			}
		})
	}
	if f.invalidRequests.Load() != 0 {
		t.Fatal("redirect followed")
	}
	f.fail = ""
	if _, err := b.AccessToken(context.Background()); err != nil {
		t.Fatal("failure recovery", err)
	}
	b.gate <- struct{}{}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := b.AccessToken(ctx); err != ErrServiceTokenUnavailable {
		t.Fatal("cancelled waiter")
	}
	<-b.gate
}
func TestServiceTokenExternalSecretRotationAndPostAuth(t *testing.T) {
	f := newServiceIssuerFixture(t)
	f.style = "client_secret_post"
	path := filepath.Join(t.TempDir(), "fixture-secret")
	if err := os.WriteFile(path, []byte(f.secret+"\n"), 0600); err != nil {
		t.Fatal(err)
	}
	b, err := NewServiceTokenBroker(context.Background(), serviceProfileForTest(t, f.document()), ServiceSecretFile(path))
	if err != nil {
		t.Fatal(err)
	}
	if _, err = b.AccessToken(context.Background()); err != nil {
		t.Fatal(err)
	}
	b.token.Expiry = time.Now()
	f.secret = "fixture-rotated-not-real"
	if err = os.WriteFile(path, []byte(f.secret), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err = b.AccessToken(context.Background()); err != nil {
		t.Fatal(err)
	}
	b.token.Expiry = time.Now()
	if err = os.WriteFile(path, []byte("invalid\nsecret"), 0600); err != nil {
		t.Fatal(err)
	}
	before := f.grants.Load()
	if _, err = b.AccessToken(context.Background()); err != ErrServiceTokenUnavailable || f.grants.Load() != before {
		t.Fatal("invalid secret reached provider", err)
	}
}
func TestServiceTokenProfileAndMetadataFailClosed(t *testing.T) {
	f := newServiceIssuerFixture(t)
	d := f.document()
	cases := []struct {
		name   string
		mutate func(*serviceTokenDocument)
	}{
		{"http_in_tls", func(d *serviceTokenDocument) { d.Transport = "TLS" }},
		{"foreign_loopback", func(d *serviceTokenDocument) { d.TokenEndpoint = "https://foreign.invalid/token" }},
		{"secret_in_url", func(d *serviceTokenDocument) { d.TokenEndpoint = "http://u:p@127.0.0.1/token" }},
		{"wildcard_permission", func(d *serviceTokenDocument) { d.Permissions = []string{"*"} }},
		{"duplicate_org", func(d *serviceTokenDocument) { d.Organizations = []string{"org", "org"} }},
		{"autodetect", func(d *serviceTokenDocument) { d.ClientAuthentication = "auto" }},
		{"refresh_zero", func(d *serviceTokenDocument) { d.RefreshBeforeSeconds = 0 }},
		{"revision", func(d *serviceTokenDocument) { d.Revision = 2 }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			value := f.document()
			tc.mutate(&value)
			raw, _ := json.Marshal(value)
			sum := sha256.Sum256(raw)
			if _, err := LoadServiceTokenProfile(raw, hex.EncodeToString(sum[:])); err != ErrServiceTokenConfiguration {
				t.Fatal("invalid profile", err)
			}
		})
	}
	raw, _ := json.Marshal(d)
	for _, invalid := range [][]byte{append([]byte(`{"schema":"duplicate",`), raw[1:]...), append(raw, []byte(`{}`)...), []byte(`{"unknown":1}`), []byte(strings.Repeat("x", 16385))} {
		sum := sha256.Sum256(invalid)
		if _, err := LoadServiceTokenProfile(invalid, hex.EncodeToString(sum[:])); err != ErrServiceTokenConfiguration {
			t.Fatal("malformed profile accepted")
		}
	}
	if _, err := LoadServiceTokenProfile(raw, strings.Repeat("0", 64)); err != ErrServiceTokenConfiguration {
		t.Fatal("wrong hash")
	}
	for _, field := range []string{"issuer", "token_endpoint", "jwks_uri", "token_endpoint_auth_methods_supported", "grant_types_supported"} {
		t.Run("discovery_"+field, func(t *testing.T) {
			f.metadata = func(m map[string]any) {
				if strings.HasSuffix(field, "supported") {
					m[field] = []string{"unsupported"}
				} else {
					m[field] = "https://foreign.invalid"
				}
			}
			_, err := NewServiceTokenBroker(context.Background(), serviceProfileForTest(t, d), func(context.Context) (string, error) { return f.secret, nil })
			if err != ErrServiceTokenConfiguration {
				t.Fatal("foreign metadata accepted", err)
			}
		})
	}
}
func FuzzServiceTokenProfile(f *testing.F) {
	valid, _ := json.Marshal(serviceTokenDocument{Schema: "elite.oidc.service-token.v1", ProfileID: "fuzz-valid", Revision: 1, Transport: "TLS", Issuer: "https://issuer.invalid", TokenEndpoint: "https://issuer.invalid/token", JWKSEndpoint: "https://issuer.invalid/jwks", ClientID: "worker", ClientAuthentication: "client_secret_basic", Audience: "api", Subject: "worker", TenantID: "tenant", Organizations: []string{"org"}, Permissions: []string{"message:send"}, Scopes: []string{"message:send"}, MaximumLifetimeSeconds: 300, RefreshBeforeSeconds: 30})
	f.Add(valid)
	f.Add([]byte(`{"schema":"elite.oidc.service-token.v1"}`))
	f.Add([]byte(`{"schema":"one","schema":"two"}`))
	f.Add([]byte(`[]`))
	f.Fuzz(func(t *testing.T, raw []byte) {
		sum := sha256.Sum256(raw)
		p, err := LoadServiceTokenProfile(raw, hex.EncodeToString(sum[:]))
		if err == nil && (p.hash == "" || !serviceSet(p.document.Scopes) || !serviceSet(p.document.Permissions)) {
			t.Fatal("unbound accepted profile")
		}
	})
}
