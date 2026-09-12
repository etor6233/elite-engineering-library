# OIDC Service Token Broker

## Metadata

```yaml
pack_id: "GO-OIDC-SERVICE-TOKEN-BROKER"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Configurable hash-bound client_credentials grant renewal using pinned oauth2/go-oidc SDKs, exact service authority validation and offline HTTP/TLS fixtures."
stacks: ["Go 1.26.8", "golang.org/x/oauth2 v0.36.0", "go-oidc/v3 v3.20.0"]
compatible_with: ["GO-ENTERPRISE-BACKEND-CORE"]
incompatible_with: ["opaque access tokens", "undeclared authority", "production loopback fixture mode"]
license_expression: "LicenseRef-Workspace-Owner AND BSD-3-Clause AND Apache-2.0"
upstream_sources: ["https://go.googlesource.com/oauth2/+/4d954e69a88d9e1ccb8439f8d5b6cbef230c4ef9", "https://github.com/coreos/go-oidc/tree/75dfa5c0626c48e0ad8b761fdd9e1dc51cb8498a"]
verified_at: "2026-09-11"
```

## Applicability

AUTHORED boundary glue only; SDK modules are DEPENDENCY_PIN. See docs/oidc-service-token-runtime.md and provenance documents for the exact supported profile. No secrets/live account/IdP implementation. Host validates profile against its selected connector before use. Global composition SCA/security/release remains a separate gate. Profile false/disabled is not a credential probe or production authorization.

## Exact file manifest

```text
CREATE config/identity/service-token.synthetic.json
CREATE docs/oidc-service-token-runtime.md
CREATE docs/provenance/oidc-service-third-party-notices.md
CREATE docs/provenance/oidc-service-token.md
CREATE internal/platform/identity/service_token_broker.go
CREATE internal/platform/identity/service_token_broker_test.go
CREATE internal/platform/identity/service_token_profile.go
CREATE internal/platform/identity/service_token_tls_test.go
CREATE licenses/coreos-go-oidc-Apache-2.0.txt
CREATE licenses/go-jose-Apache-2.0.txt
CREATE licenses/golang-oauth2-BSD-3-Clause.txt
```

## Materialization blocks

### FILE: `config/identity/service-token.synthetic.json`

```yaml
block_id: "GO-OIDC-SERVICE-TOKEN-BROKER:config/identity/service-token.synthetic.json:v1"
operation: CREATE
provenance: AUTHORED
source: "local OAuth SDK composition and scoped offline fixture"
license: "LicenseRef-Workspace-Owner"
sha256: "e9024a2da76df3366ca4caa870b64942f8ed011ecce80dc0e41d148eb574d3ae"
variables: []
secrets_allowed: false
```

````json
{
  "schema": "elite.oidc.service-token.v1",
  "profile_id": "synthetic-whatsapp-service",
  "revision": 1,
  "transport": "LOOPBACK_FIXTURE",
  "issuer": "http://127.0.0.1:9",
  "token_endpoint": "http://127.0.0.1:9/token",
  "jwks_endpoint": "http://127.0.0.1:9/jwks",
  "client_id": "fixture-worker",
  "client_authentication": "client_secret_basic",
  "audience": "fixture-api",
  "subject": "fixture-worker",
  "tenant_id": "00000000-0000-4000-8000-000000000001",
  "organization_ids": [
    "00000000-0000-4000-8000-000000000002"
  ],
  "permissions": [
    "message:send",
    "conversation:read"
  ],
  "scopes": [
    "message:send",
    "conversation:read"
  ],
  "maximum_lifetime_seconds": 300,
  "refresh_before_seconds": 30
}
````

### FILE: `docs/oidc-service-token-runtime.md`

```yaml
block_id: "GO-OIDC-SERVICE-TOKEN-BROKER:docs/oidc-service-token-runtime.md:v1"
operation: CREATE
provenance: AUTHORED
source: "local OAuth SDK composition and scoped offline fixture"
license: "LicenseRef-Workspace-Owner"
sha256: "2ddefb712ab6c892099f77f260c46c4f20e39f5224eb261b795ef97443fdbec4"
variables: []
secrets_allowed: false
```

````text
# OIDC service token runtime

This reusable component performs the OAuth2 client_credentials grant with the fixed Go oauth2 SDK, verifies the returned RS256 JWT with the fixed go-oidc SDK, and renews it before expiry. It is not an IdP implementation. No grant requests occur until AccessToken is called; construction performs discovery. Disabled host modules must not construct/call the broker.

Host integration:

1. Materialize the nonsecret profile for the chosen account and pin its exact SHA256 in deployment configuration. LoadServiceTokenProfileFile(profilePath, profileSHA256) rejects changed/unknown/duplicate/incomplete data. Configuration is not a credential probe.
2. Compare profile.Binding() with the existing host verifier issuer/audience and the connector tenant, subject, organization, permissions and scopes. Require Transport=TLS outside the explicit loopback fixture. Reject a mismatch before mounting traffic/polling; do not derive these authorities from an inbound request.
3. Construct NewServiceTokenBroker(ctx, profile, ServiceSecretFile(externalSecretPath)). That file is supplied later through the user's secret manager, outside the project. The host account owns its directory/ACL. The reader reloads it on each grant so external secret replacement takes effect on the next renewal. No secret is recorded in a profile, URL, log, receipt or command line.
4. Call AccessToken(ctx) per poll/request and continue using the existing backend verifier before effects. On ErrServiceTokenUnavailable defer the poll; never use a stale token or treat the unavailable grant as successful work. Concurrent calls collapse to one grant. Renewal failures return no cached token. There is no background goroutine and no credential/token persistence to clean up.

Suggested nonsecret environment names (the host may namespace them per connector): OIDC_SERVICE_ENABLED=false, OIDC_SERVICE_PROFILE_FILE=, OIDC_SERVICE_PROFILE_SHA256=, OIDC_SERVICE_CLIENT_SECRET_FILE=. Enabling with missing path/hash/secret reference rejects startup. Existing WHATSAPP_SERVICE_TOKEN_FILE may remain an explicitly selected static-token mode; do not silently mix two credentials. Service-token mode selection and equality with the WhatsApp profile are host-owned glue.

Profile schema elite.oidc.service-token.v1 revision1 supports one issuer, one exact token endpoint and one exact JWKS endpoint; RS256 JWT bearer tokens; client_secret_basic or client_secret_post with no auth-style probing/retry; one audience; exact subject/tenant; nonempty exact organization, permission and scope sets without wildcard; optional resource parameter or explicit audience parameter, never both. scopes in the signed token are required. A missing response scope is accepted only because OAuth2 permits omission when unchanged; a provided response scope must match. Opaque tokens, refresh_token responses, undeclared extra authority, missing expiry, future nbf/iat, excessive lifetime and ambiguous audiences fail closed. This supported contract must be provisioned by the selected real issuer later; no issuer behavior is fabricated.

The client fetches only the three fixed URLs, follows no redirects, bounds response bytes to256KiB and requests to15seconds, and emits generic errors without upstream body/secret text. TLS certificate validation uses the host trust store. LOOPBACK_FIXTURE permits only literal loopback IP endpoints and is not a production mode. The synthetic profile's127.0.0.1:9 endpoint intentionally provides no service.

This is grant renewal, not refresh_token use: RFC6749 section4.4.3 says a refresh token should not be returned for this grant. Human portal Authorization Code refresh/revocation is a separate owner and remains outside this component. JWT access-token use is the backend's explicit existing contract; this adapter does not claim all OAuth providers or the whole RFC9068 profile.

Evidence: real HTTP token/discovery/JWKS SDK tests;20 simultaneous callers; key rotation;22 wrong-claim/lifetime responses; negative metadata/profile/redirect/oversized/error paths; secret rotation; real TLS fixture; bad signature; finite profile fuzz. These are offline/sandbox results. Global composition SCA, deployment IdP configuration, target secret storage/network boundaries and production acceptance are separate gates.
````

### FILE: `docs/provenance/oidc-service-third-party-notices.md`

```yaml
block_id: "GO-OIDC-SERVICE-TOKEN-BROKER:docs/provenance/oidc-service-third-party-notices.md:v1"
operation: CREATE
provenance: AUTHORED
source: "local OAuth SDK composition and scoped offline fixture"
license: "LicenseRef-Workspace-Owner"
sha256: "4cc131838dece6f7c3d81fe65a047267c1b29283e3c6422b6fec993447d1fa98"
variables: []
secrets_allowed: false
```

````text
# Service token dependency notices

This delta uses the existing exact module graph. Go oauth2 v0.36.0 is BSD-3-Clause; coreos/go-oidc v3.20.0 and go-jose/v4 v4.1.4 are Apache-2.0. Their complete license bytes are included under licenses/. The three module versions are DEPENDENCY_PIN, not local authorship. Application/profile/transport/test glue is AUTHORED / LicenseRef-Workspace-Owner. No vendor endorsement. Transitive module notices and the complete SBOM remain composition-owned.
````

### FILE: `docs/provenance/oidc-service-token.md`

```yaml
block_id: "GO-OIDC-SERVICE-TOKEN-BROKER:docs/provenance/oidc-service-token.md:v1"
operation: CREATE
provenance: AUTHORED
source: "local OAuth SDK composition and scoped offline fixture"
license: "LicenseRef-Workspace-Owner"
sha256: "e1ea06758a8872a4db87460c5016d88ef06a45f1f577e92beffdd3679fe6be2f"
variables: []
secrets_allowed: false
```

````text
# Service token provenance and admission boundary

AUTHORED: profile/hash parsing, exact authority comparisons, bounded transport, external secret reader, lifetime/cache concurrency glue, host contract and all local fixtures. None is attributed to an external company. Existing Principal/OIDCVerifier remains owned by GO-ENTERPRISE-BACKEND-CORE and is not recopied by this pack.

DEPENDENCY_PIN: golang.org/x/oauth2 v0.36.0, clientcredentials.Config.Token executes the grant and builds its Basic/Post authentication. github.com/coreos/go-oidc/v3 v3.20.0, NewProvider and IDTokenVerifier execute discovery, signature/issuer/audience/expiry verification and RemoteKeySet key rotation. github.com/go-jose/go-jose/v4 v4.1.4 signs only ephemeral test-fixture tokens. No module/version is introduced. Exact local module info, source files, licenses, go.mod and go.sum are hashed in oidc-service-source-receipt.json.

Official method authorities: https://pkg.go.dev/golang.org/x/oauth2@v0.36.0/clientcredentials ; https://pkg.go.dev/github.com/coreos/go-oidc/v3@v3.20.0/oidc ; https://www.rfc-editor.org/rfc/rfc6749#section-4.4 . Sources are inspected; SDK integration tests execute the fixed existing modules. No upstream full-suite execution or worldwide provider compatibility is claimed.

G0/G1 source identity/license, G2 security boundary and G3 failure design are recorded before coding in oidc-service-admission.md. Local G4 unit/contract, G5 claim/transport negatives plus finite fuzz, G6 concurrency/renewal and G7 recovery/TLS checks are evidenced by exact logs. G8 exact pack reconstruction is separately recorded. This scoped evidence does not close root composition SCA/SAST/security/release or certify production.

AUTHORED glue is inevitable here: a provider SDK cannot know the library tenant, organization, connector subject or permission boundary, configured exact endpoints, external secret path or host's cancellation/poll semantics. No custom signing, OAuth token protocol or IdP framework is implemented.
````

### FILE: `internal/platform/identity/service_token_broker.go`

```yaml
block_id: "GO-OIDC-SERVICE-TOKEN-BROKER:internal/platform/identity/service_token_broker.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local OAuth SDK composition and scoped offline fixture"
license: "LicenseRef-Workspace-Owner"
sha256: "a0389d4ee9e1c87f1698a1444b194891c5cceac683a1daa5362a4406f17659bb"
variables: []
secrets_allowed: false
```

````go
package identity

// AUTHORED lifecycle/transport glue; protocol and JWT verification are DEPENDENCY_PIN.
import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

type ServiceSecretReader func(context.Context) (string, error)

// ServiceSecretFile reads on each grant, allowing externally managed secret rotation.
// Files and their directory/ACL are provisioned by the host secret manager, never the pack.
func ServiceSecretFile(path string) ServiceSecretReader {
	return func(ctx context.Context) (string, error) {
		if ctx.Err() != nil {
			return "", ErrServiceTokenUnavailable
		}
		f, err := os.Open(path)
		if err != nil {
			return "", ErrServiceTokenUnavailable
		}
		defer f.Close()
		info, err := f.Stat()
		if err != nil || !info.Mode().IsRegular() || info.Size() > 4096 {
			return "", ErrServiceTokenUnavailable
		}
		raw, err := io.ReadAll(io.LimitReader(f, 4097))
		if err != nil || len(raw) > 4096 {
			return "", ErrServiceTokenUnavailable
		}
		value := strings.TrimSuffix(strings.TrimSuffix(string(raw), "\n"), "\r")
		if len(value) < 1 || strings.ContainsAny(value, "\x00\r\n") {
			return "", ErrServiceTokenUnavailable
		}
		return value, nil
	}
}

type ServiceTokenBroker struct {
	profile  ServiceTokenProfile
	secret   ServiceSecretReader
	client   *http.Client
	verifier *oidc.IDTokenVerifier
	gate     chan struct{}
	token    *oauth2.Token
	now      func() time.Time
}

func NewServiceTokenBroker(ctx context.Context, p ServiceTokenProfile, secret ServiceSecretReader) (*ServiceTokenBroker, error) {
	return newServiceTokenBroker(ctx, p, secret, http.DefaultTransport)
}

// The transport seam is private and used for a loopback TLS fixture with its own CA.
func newServiceTokenBroker(ctx context.Context, p ServiceTokenProfile, secret ServiceSecretReader, transport http.RoundTripper) (*ServiceTokenBroker, error) {
	if p.hash == "" || secret == nil {
		return nil, ErrServiceTokenConfiguration
	}
	d := p.document
	client := &http.Client{Timeout: 15 * time.Second, Transport: serviceTokenTransport{base: transport, allowed: map[string]bool{strings.TrimSuffix(d.Issuer, "/") + "/.well-known/openid-configuration": true, d.TokenEndpoint: true, d.JWKSEndpoint: true}}, CheckRedirect: func(*http.Request, []*http.Request) error { return ErrServiceTokenUnavailable }}
	provider, err := oidc.NewProvider(oidc.ClientContext(ctx, client), d.Issuer)
	if err != nil {
		return nil, ErrServiceTokenConfiguration
	}
	var metadata struct {
		JWKSEndpoint string   `json:"jwks_uri"`
		AuthMethods  []string `json:"token_endpoint_auth_methods_supported"`
		Grants       []string `json:"grant_types_supported"`
	}
	if provider.Claims(&metadata) != nil || provider.Endpoint().TokenURL != d.TokenEndpoint || metadata.JWKSEndpoint != d.JWKSEndpoint {
		return nil, ErrServiceTokenConfiguration
	}
	contains := func(list []string, value string) bool {
		for _, item := range list {
			if item == value {
				return true
			}
		}
		return false
	}
	if len(metadata.AuthMethods) > 0 && !contains(metadata.AuthMethods, d.ClientAuthentication) || len(metadata.Grants) > 0 && !contains(metadata.Grants, "client_credentials") {
		return nil, ErrServiceTokenConfiguration
	}
	b := &ServiceTokenBroker{profile: p, secret: secret, client: client, gate: make(chan struct{}, 1), now: time.Now}
	b.verifier = provider.Verifier(&oidc.Config{ClientID: d.Audience, SupportedSigningAlgs: []string{oidc.RS256}, Now: func() time.Time { return b.now() }})
	return b, nil
}
func (b *ServiceTokenBroker) Binding() ServiceTokenBinding { return b.profile.Binding() }

// AccessToken returns only a verified short-lived bearer. Failed renewals do not
// return the previous token and do not cache malformed/foreign provider responses.
func (b *ServiceTokenBroker) AccessToken(ctx context.Context) (string, error) {
	if b == nil || ctx == nil {
		return "", ErrServiceTokenUnavailable
	}
	select {
	case b.gate <- struct{}{}:
		defer func() { <-b.gate }()
	case <-ctx.Done():
		return "", ErrServiceTokenUnavailable
	}
	if ctx.Err() != nil {
		return "", ErrServiceTokenUnavailable
	}
	d := b.profile.document
	threshold := b.now().Add(time.Duration(d.RefreshBeforeSeconds) * time.Second)
	if b.token != nil && b.token.Expiry.After(threshold) {
		return b.token.AccessToken, nil
	}
	b.token = nil
	secret, err := b.secret(ctx)
	if err != nil || len(secret) < 1 || len(secret) > 4096 || strings.ContainsAny(secret, "\x00\r\n") {
		return "", ErrServiceTokenUnavailable
	}
	style := oauth2.AuthStyleInHeader
	if d.ClientAuthentication == "client_secret_post" {
		style = oauth2.AuthStyleInParams
	}
	params := url.Values{}
	if d.Resource != "" {
		params.Set("resource", d.Resource)
	}
	if d.AudienceParameter != "" {
		params.Set("audience", d.AudienceParameter)
	}
	config := clientcredentials.Config{ClientID: d.ClientID, ClientSecret: secret, TokenURL: d.TokenEndpoint, Scopes: append([]string(nil), d.Scopes...), AuthStyle: style, EndpointParams: params}
	token, err := config.Token(context.WithValue(ctx, oauth2.HTTPClient, b.client))
	if err != nil || token == nil || !strings.EqualFold(token.TokenType, "Bearer") || len(token.AccessToken) == 0 || len(token.AccessToken) > 16384 || token.RefreshToken != "" || token.Expiry.IsZero() {
		return "", ErrServiceTokenUnavailable
	}
	jwt, err := b.verifier.Verify(ctx, token.AccessToken)
	if err != nil {
		return "", ErrServiceTokenUnavailable
	}
	var claims struct {
		TenantID        string   `json:"tenant_id"`
		Organizations   []string `json:"organization_ids"`
		Permissions     []string `json:"permissions"`
		Scope           string   `json:"scope"`
		NotBefore       int64    `json:"nbf"`
		AuthorizedParty string   `json:"azp"`
		ClientID        string   `json:"client_id"`
	}
	if jwt.Claims(&claims) != nil || jwt.Issuer != d.Issuer || jwt.Subject != d.Subject || claims.TenantID != d.TenantID || !serviceEqual(claims.Organizations, d.Organizations) || !serviceEqual(claims.Permissions, d.Permissions) || !serviceEqual(strings.Fields(claims.Scope), d.Scopes) || len(jwt.Audience) != 1 || jwt.Audience[0] != d.Audience {
		return "", ErrServiceTokenUnavailable
	}
	now := b.now()
	if claims.NotBefore > now.Unix() || jwt.IssuedAt.IsZero() || jwt.IssuedAt.After(now.Add(30*time.Second)) || jwt.Expiry.Sub(jwt.IssuedAt) > time.Duration(d.MaximumLifetimeSeconds)*time.Second || !jwt.Expiry.After(threshold) {
		return "", ErrServiceTokenUnavailable
	}
	if claims.AuthorizedParty != "" && claims.AuthorizedParty != d.ClientID || claims.ClientID != "" && claims.ClientID != d.ClientID {
		return "", ErrServiceTokenUnavailable
	}
	if raw := token.Extra("scope"); raw != nil {
		scope, ok := raw.(string)
		if !ok || !serviceEqual(strings.Fields(scope), d.Scopes) {
			return "", ErrServiceTokenUnavailable
		}
	}
	if jwt.Expiry.Before(token.Expiry) {
		token.Expiry = jwt.Expiry
	}
	if !token.Expiry.After(threshold) || token.Expiry.Sub(now) > time.Duration(d.MaximumLifetimeSeconds+30)*time.Second {
		return "", ErrServiceTokenUnavailable
	}
	b.token = token
	return token.AccessToken, nil
}

type serviceTokenTransport struct {
	base    http.RoundTripper
	allowed map[string]bool
}

func (t serviceTokenTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	if !t.allowed[request.URL.String()] || (request.Method != http.MethodGet && request.Method != http.MethodPost) {
		return nil, ErrServiceTokenUnavailable
	}
	response, err := t.base.RoundTrip(request)
	if err != nil {
		return nil, ErrServiceTokenUnavailable
	}
	defer response.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(response.Body, 262145))
	if err != nil || len(raw) > 262144 {
		return nil, ErrServiceTokenUnavailable
	}
	response.Body = io.NopCloser(bytes.NewReader(raw))
	return response, nil
}
````

### FILE: `internal/platform/identity/service_token_broker_test.go`

```yaml
block_id: "GO-OIDC-SERVICE-TOKEN-BROKER:internal/platform/identity/service_token_broker_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local OAuth SDK composition and scoped offline fixture"
license: "LicenseRef-Workspace-Owner"
sha256: "53bf3eb62131b7f256a723d2c4564fbda553c7d901eea904c57cb3d70ee8d7e2"
variables: []
secrets_allowed: false
```

````go
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
````

### FILE: `internal/platform/identity/service_token_profile.go`

```yaml
block_id: "GO-OIDC-SERVICE-TOKEN-BROKER:internal/platform/identity/service_token_profile.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local OAuth SDK composition and scoped offline fixture"
license: "LicenseRef-Workspace-Owner"
sha256: "d4cda472b005af8b484f285fdf446a8565fa006fffea6b3bb9500fc7f2b3024c"
variables: []
secrets_allowed: false
```

````go
package identity

// AUTHORED configuration and claim binding glue. OAuth/OIDC runs in pinned SDKs.
import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/netip"
	"net/url"
	"os"
	"slices"
	"strings"
)

var ErrServiceTokenConfiguration = errors.New("service token configuration rejected")
var ErrServiceTokenUnavailable = errors.New("service token unavailable")

type serviceTokenDocument struct {
	Schema                 string   `json:"schema"`
	ProfileID              string   `json:"profile_id"`
	Revision               int      `json:"revision"`
	Transport              string   `json:"transport"`
	Issuer                 string   `json:"issuer"`
	TokenEndpoint          string   `json:"token_endpoint"`
	JWKSEndpoint           string   `json:"jwks_endpoint"`
	ClientID               string   `json:"client_id"`
	ClientAuthentication   string   `json:"client_authentication"`
	Audience               string   `json:"audience"`
	Subject                string   `json:"subject"`
	TenantID               string   `json:"tenant_id"`
	Organizations          []string `json:"organization_ids"`
	Permissions            []string `json:"permissions"`
	Scopes                 []string `json:"scopes"`
	Resource               string   `json:"resource,omitempty"`
	AudienceParameter      string   `json:"audience_parameter,omitempty"`
	MaximumLifetimeSeconds int      `json:"maximum_lifetime_seconds"`
	RefreshBeforeSeconds   int      `json:"refresh_before_seconds"`
}

// ServiceTokenProfile is immutable after hash verification; zero value is invalid.
type ServiceTokenProfile struct {
	document serviceTokenDocument
	hash     string
}
type ServiceTokenBinding struct {
	ProfileID, DocumentSHA256, Transport, Issuer, Audience, ClientID, Subject, TenantID string
	Organizations, Permissions, Scopes                                                  []string
}

func (p ServiceTokenProfile) Binding() ServiceTokenBinding {
	d := p.document
	return ServiceTokenBinding{d.ProfileID, p.hash, d.Transport, d.Issuer, d.Audience, d.ClientID, d.Subject, d.TenantID, slices.Clone(d.Organizations), slices.Clone(d.Permissions), slices.Clone(d.Scopes)}
}
func LoadServiceTokenProfileFile(path, digest string) (ServiceTokenProfile, error) {
	f, err := os.Open(path)
	if err != nil {
		return ServiceTokenProfile{}, ErrServiceTokenConfiguration
	}
	defer f.Close()
	raw, err := io.ReadAll(io.LimitReader(f, 16385))
	if err != nil {
		return ServiceTokenProfile{}, ErrServiceTokenConfiguration
	}
	return LoadServiceTokenProfile(raw, digest)
}
func LoadServiceTokenProfile(raw []byte, digest string) (ServiceTokenProfile, error) {
	reject := func() (ServiceTokenProfile, error) { return ServiceTokenProfile{}, ErrServiceTokenConfiguration }
	if len(raw) == 0 || len(raw) > 16384 || len(digest) != 64 {
		return reject()
	}
	sum := sha256.Sum256(raw)
	if hex.EncodeToString(sum[:]) != digest {
		return reject()
	}
	// Reject duplicate member names before normal decoding; this schema has no objects nested in members.
	first := json.NewDecoder(bytes.NewReader(raw))
	token, err := first.Token()
	if err != nil || token != json.Delim('{') {
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
		var value json.RawMessage
		if first.Decode(&value) != nil {
			return reject()
		}
	}
	if token, err = first.Token(); err != nil || token != json.Delim('}') {
		return reject()
	}
	if _, err = first.Token(); err != io.EOF {
		return reject()
	}
	var d serviceTokenDocument
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	if dec.Decode(&d) != nil {
		return reject()
	}
	if d.Schema != "elite.oidc.service-token.v1" || d.Revision != 1 || !serviceAtom(d.ProfileID) || !serviceAtom(d.ClientID) || !serviceAtom(d.Audience) || !serviceAtom(d.Subject) || !serviceAtom(d.TenantID) {
		return reject()
	}
	if d.Transport != "TLS" && d.Transport != "LOOPBACK_FIXTURE" {
		return reject()
	}
	if !serviceURL(d.Issuer, d.Transport) || !serviceURL(d.TokenEndpoint, d.Transport) || !serviceURL(d.JWKSEndpoint, d.Transport) {
		return reject()
	}
	if d.ClientAuthentication != "client_secret_basic" && d.ClientAuthentication != "client_secret_post" {
		return reject()
	}
	if !serviceSet(d.Organizations) || !serviceSet(d.Permissions) || !serviceSet(d.Scopes) {
		return reject()
	}
	if d.Resource != "" && (!serviceURL(d.Resource, "TLS") || d.AudienceParameter != "") {
		return reject()
	}
	if d.AudienceParameter != "" && d.AudienceParameter != d.Audience {
		return reject()
	}
	if d.MaximumLifetimeSeconds < 60 || d.MaximumLifetimeSeconds > 86400 || d.RefreshBeforeSeconds < 5 || d.RefreshBeforeSeconds > 300 || d.RefreshBeforeSeconds*2 >= d.MaximumLifetimeSeconds {
		return reject()
	}
	return ServiceTokenProfile{document: d, hash: digest}, nil
}
func serviceAtom(s string) bool {
	return len(s) > 0 && len(s) <= 200 && strings.TrimSpace(s) == s && !strings.ContainsAny(s, "\x00\r\n\t ") && s != "*"
}
func serviceSet(values []string) bool {
	if len(values) == 0 || len(values) > 100 {
		return false
	}
	seen := map[string]bool{}
	for _, v := range values {
		if !serviceAtom(v) || seen[v] {
			return false
		}
		seen[v] = true
	}
	return true
}
func serviceEqual(a, b []string) bool {
	if !serviceSet(a) || len(a) != len(b) {
		return false
	}
	a = slices.Clone(a)
	b = slices.Clone(b)
	slices.Sort(a)
	slices.Sort(b)
	return slices.Equal(a, b)
}
func serviceURL(raw, mode string) bool {
	u, err := url.Parse(raw)
	if err != nil || len(raw) > 2048 || u.Host == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" || u.Opaque != "" || u.String() != raw {
		return false
	}
	if mode == "LOOPBACK_FIXTURE" {
		ip, err := netip.ParseAddr(u.Hostname())
		return err == nil && ip.IsLoopback() && (u.Scheme == "http" || u.Scheme == "https")
	}
	return u.Scheme == "https"
}
````

### FILE: `internal/platform/identity/service_token_tls_test.go`

```yaml
block_id: "GO-OIDC-SERVICE-TOKEN-BROKER:internal/platform/identity/service_token_tls_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local OAuth SDK composition and scoped offline fixture"
license: "LicenseRef-Workspace-Owner"
sha256: "bea6952f9fb2c45b13ac0a27c38041867731db047cf2c539397419999441d758"
variables: []
secrets_allowed: false
```

````go
package identity

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	jose "github.com/go-jose/go-jose/v4"
)

func TestServiceTokenTLSChainAndSignatureFailure(t *testing.T) {
	f := newServiceIssuerFixture(t)
	handler := f.server.Config.Handler
	tlsServer := httptest.NewUnstartedServer(handler)
	tlsServer.Config.ErrorLog = log.New(io.Discard, "", 0)
	tlsServer.StartTLS()
	t.Cleanup(tlsServer.Close)
	f.server = tlsServer
	p := serviceProfileForTest(t, f.document())
	secret := func(context.Context) (string, error) { return f.secret, nil }
	if _, err := NewServiceTokenBroker(context.Background(), p, secret); err != ErrServiceTokenConfiguration {
		t.Fatal("untrusted CA accepted", err)
	}
	b, err := newServiceTokenBroker(context.Background(), p, secret, tlsServer.Client().Transport)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = b.AccessToken(context.Background()); err != nil {
		t.Fatal("trusted fixture CA failed", err)
	}
	// A token signed by a different key with the same kid cannot pass a cached JWKS key.
	old := f.keys[0]
	replacement, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	f.keys[0] = replacement
	b.token.Expiry = time.Now()
	// Keep JWKS transport serving old key while token endpoint signs with replacement.
	original := tlsServer.Config.Handler
	tlsServer.Config.Handler = httpHandlerForJWKS(old, original)
	if token, err := b.AccessToken(context.Background()); err != ErrServiceTokenUnavailable || token != "" {
		t.Fatal("wrong signature accepted", err)
	}
}
func httpHandlerForJWKS(key *rsa.PrivateKey, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/jwks" {
			next.ServeHTTP(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(jose.JSONWebKeySet{Keys: []jose.JSONWebKey{{Key: &key.PublicKey, KeyID: "fixture-a", Use: "sig", Algorithm: "RS256"}}})
	})
}
````

### FILE: `licenses/coreos-go-oidc-Apache-2.0.txt`

```yaml
block_id: "GO-OIDC-SERVICE-TOKEN-BROKER:licenses/coreos-go-oidc-Apache-2.0.txt:v1"
operation: CREATE
provenance: VERBATIM
source: "exact pinned dependency LICENSE bytes; oidc-service-source-receipt.json"
license: "Apache-2.0"
sha256: "cb5e8e7e5f4a3988e1063c142c60dc2df75605f4c46515e776e3aca6df976e14"
variables: []
secrets_allowed: false
```

````text
Apache License
                           Version 2.0, January 2004
                        http://www.apache.org/licenses/

   TERMS AND CONDITIONS FOR USE, REPRODUCTION, AND DISTRIBUTION

   1. Definitions.

      "License" shall mean the terms and conditions for use, reproduction,
      and distribution as defined by Sections 1 through 9 of this document.

      "Licensor" shall mean the copyright owner or entity authorized by
      the copyright owner that is granting the License.

      "Legal Entity" shall mean the union of the acting entity and all
      other entities that control, are controlled by, or are under common
      control with that entity. For the purposes of this definition,
      "control" means (i) the power, direct or indirect, to cause the
      direction or management of such entity, whether by contract or
      otherwise, or (ii) ownership of fifty percent (50%) or more of the
      outstanding shares, or (iii) beneficial ownership of such entity.

      "You" (or "Your") shall mean an individual or Legal Entity
      exercising permissions granted by this License.

      "Source" form shall mean the preferred form for making modifications,
      including but not limited to software source code, documentation
      source, and configuration files.

      "Object" form shall mean any form resulting from mechanical
      transformation or translation of a Source form, including but
      not limited to compiled object code, generated documentation,
      and conversions to other media types.

      "Work" shall mean the work of authorship, whether in Source or
      Object form, made available under the License, as indicated by a
      copyright notice that is included in or attached to the work
      (an example is provided in the Appendix below).

      "Derivative Works" shall mean any work, whether in Source or Object
      form, that is based on (or derived from) the Work and for which the
      editorial revisions, annotations, elaborations, or other modifications
      represent, as a whole, an original work of authorship. For the purposes
      of this License, Derivative Works shall not include works that remain
      separable from, or merely link (or bind by name) to the interfaces of,
      the Work and Derivative Works thereof.

      "Contribution" shall mean any work of authorship, including
      the original version of the Work and any modifications or additions
      to that Work or Derivative Works thereof, that is intentionally
      submitted to Licensor for inclusion in the Work by the copyright owner
      or by an individual or Legal Entity authorized to submit on behalf of
      the copyright owner. For the purposes of this definition, "submitted"
      means any form of electronic, verbal, or written communication sent
      to the Licensor or its representatives, including but not limited to
      communication on electronic mailing lists, source code control systems,
      and issue tracking systems that are managed by, or on behalf of, the
      Licensor for the purpose of discussing and improving the Work, but
      excluding communication that is conspicuously marked or otherwise
      designated in writing by the copyright owner as "Not a Contribution."

      "Contributor" shall mean Licensor and any individual or Legal Entity
      on behalf of whom a Contribution has been received by Licensor and
      subsequently incorporated within the Work.

   2. Grant of Copyright License. Subject to the terms and conditions of
      this License, each Contributor hereby grants to You a perpetual,
      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
      copyright license to reproduce, prepare Derivative Works of,
      publicly display, publicly perform, sublicense, and distribute the
      Work and such Derivative Works in Source or Object form.

   3. Grant of Patent License. Subject to the terms and conditions of
      this License, each Contributor hereby grants to You a perpetual,
      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
      (except as stated in this section) patent license to make, have made,
      use, offer to sell, sell, import, and otherwise transfer the Work,
      where such license applies only to those patent claims licensable
      by such Contributor that are necessarily infringed by their
      Contribution(s) alone or by combination of their Contribution(s)
      with the Work to which such Contribution(s) was submitted. If You
      institute patent litigation against any entity (including a
      cross-claim or counterclaim in a lawsuit) alleging that the Work
      or a Contribution incorporated within the Work constitutes direct
      or contributory patent infringement, then any patent licenses
      granted to You under this License for that Work shall terminate
      as of the date such litigation is filed.

   4. Redistribution. You may reproduce and distribute copies of the
      Work or Derivative Works thereof in any medium, with or without
      modifications, and in Source or Object form, provided that You
      meet the following conditions:

      (a) You must give any other recipients of the Work or
          Derivative Works a copy of this License; and

      (b) You must cause any modified files to carry prominent notices
          stating that You changed the files; and

      (c) You must retain, in the Source form of any Derivative Works
          that You distribute, all copyright, patent, trademark, and
          attribution notices from the Source form of the Work,
          excluding those notices that do not pertain to any part of
          the Derivative Works; and

      (d) If the Work includes a "NOTICE" text file as part of its
          distribution, then any Derivative Works that You distribute must
          include a readable copy of the attribution notices contained
          within such NOTICE file, excluding those notices that do not
          pertain to any part of the Derivative Works, in at least one
          of the following places: within a NOTICE text file distributed
          as part of the Derivative Works; within the Source form or
          documentation, if provided along with the Derivative Works; or,
          within a display generated by the Derivative Works, if and
          wherever such third-party notices normally appear. The contents
          of the NOTICE file are for informational purposes only and
          do not modify the License. You may add Your own attribution
          notices within Derivative Works that You distribute, alongside
          or as an addendum to the NOTICE text from the Work, provided
          that such additional attribution notices cannot be construed
          as modifying the License.

      You may add Your own copyright statement to Your modifications and
      may provide additional or different license terms and conditions
      for use, reproduction, or distribution of Your modifications, or
      for any such Derivative Works as a whole, provided Your use,
      reproduction, and distribution of the Work otherwise complies with
      the conditions stated in this License.

   5. Submission of Contributions. Unless You explicitly state otherwise,
      any Contribution intentionally submitted for inclusion in the Work
      by You to the Licensor shall be under the terms and conditions of
      this License, without any additional terms or conditions.
      Notwithstanding the above, nothing herein shall supersede or modify
      the terms of any separate license agreement you may have executed
      with Licensor regarding such Contributions.

   6. Trademarks. This License does not grant permission to use the trade
      names, trademarks, service marks, or product names of the Licensor,
      except as required for reasonable and customary use in describing the
      origin of the Work and reproducing the content of the NOTICE file.

   7. Disclaimer of Warranty. Unless required by applicable law or
      agreed to in writing, Licensor provides the Work (and each
      Contributor provides its Contributions) on an "AS IS" BASIS,
      WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
      implied, including, without limitation, any warranties or conditions
      of TITLE, NON-INFRINGEMENT, MERCHANTABILITY, or FITNESS FOR A
      PARTICULAR PURPOSE. You are solely responsible for determining the
      appropriateness of using or redistributing the Work and assume any
      risks associated with Your exercise of permissions under this License.

   8. Limitation of Liability. In no event and under no legal theory,
      whether in tort (including negligence), contract, or otherwise,
      unless required by applicable law (such as deliberate and grossly
      negligent acts) or agreed to in writing, shall any Contributor be
      liable to You for damages, including any direct, indirect, special,
      incidental, or consequential damages of any character arising as a
      result of this License or out of the use or inability to use the
      Work (including but not limited to damages for loss of goodwill,
      work stoppage, computer failure or malfunction, or any and all
      other commercial damages or losses), even if such Contributor
      has been advised of the possibility of such damages.

   9. Accepting Warranty or Additional Liability. While redistributing
      the Work or Derivative Works thereof, You may choose to offer,
      and charge a fee for, acceptance of support, warranty, indemnity,
      or other liability obligations and/or rights consistent with this
      License. However, in accepting such obligations, You may act only
      on Your own behalf and on Your sole responsibility, not on behalf
      of any other Contributor, and only if You agree to indemnify,
      defend, and hold each Contributor harmless for any liability
      incurred by, or claims asserted against, such Contributor by reason
      of your accepting any such warranty or additional liability.

   END OF TERMS AND CONDITIONS

   APPENDIX: How to apply the Apache License to your work.

      To apply the Apache License to your work, attach the following
      boilerplate notice, with the fields enclosed by brackets "{}"
      replaced with your own identifying information. (Don't include
      the brackets!)  The text should be enclosed in the appropriate
      comment syntax for the file format. We also recommend that a
      file or class name and description of purpose be included on the
      same "printed page" as the copyright notice for easier
      identification within third-party archives.

   Copyright {yyyy} {name of copyright owner}

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

````

### FILE: `licenses/go-jose-Apache-2.0.txt`

```yaml
block_id: "GO-OIDC-SERVICE-TOKEN-BROKER:licenses/go-jose-Apache-2.0.txt:v1"
operation: CREATE
provenance: VERBATIM
source: "exact pinned dependency LICENSE bytes; oidc-service-source-receipt.json"
license: "Apache-2.0"
sha256: "cfc7749b96f63bd31c3c42b5c471bf756814053e847c10f3eb003417bc523d30"
variables: []
secrets_allowed: false
```

````text

                                 Apache License
                           Version 2.0, January 2004
                        http://www.apache.org/licenses/

   TERMS AND CONDITIONS FOR USE, REPRODUCTION, AND DISTRIBUTION

   1. Definitions.

      "License" shall mean the terms and conditions for use, reproduction,
      and distribution as defined by Sections 1 through 9 of this document.

      "Licensor" shall mean the copyright owner or entity authorized by
      the copyright owner that is granting the License.

      "Legal Entity" shall mean the union of the acting entity and all
      other entities that control, are controlled by, or are under common
      control with that entity. For the purposes of this definition,
      "control" means (i) the power, direct or indirect, to cause the
      direction or management of such entity, whether by contract or
      otherwise, or (ii) ownership of fifty percent (50%) or more of the
      outstanding shares, or (iii) beneficial ownership of such entity.

      "You" (or "Your") shall mean an individual or Legal Entity
      exercising permissions granted by this License.

      "Source" form shall mean the preferred form for making modifications,
      including but not limited to software source code, documentation
      source, and configuration files.

      "Object" form shall mean any form resulting from mechanical
      transformation or translation of a Source form, including but
      not limited to compiled object code, generated documentation,
      and conversions to other media types.

      "Work" shall mean the work of authorship, whether in Source or
      Object form, made available under the License, as indicated by a
      copyright notice that is included in or attached to the work
      (an example is provided in the Appendix below).

      "Derivative Works" shall mean any work, whether in Source or Object
      form, that is based on (or derived from) the Work and for which the
      editorial revisions, annotations, elaborations, or other modifications
      represent, as a whole, an original work of authorship. For the purposes
      of this License, Derivative Works shall not include works that remain
      separable from, or merely link (or bind by name) to the interfaces of,
      the Work and Derivative Works thereof.

      "Contribution" shall mean any work of authorship, including
      the original version of the Work and any modifications or additions
      to that Work or Derivative Works thereof, that is intentionally
      submitted to Licensor for inclusion in the Work by the copyright owner
      or by an individual or Legal Entity authorized to submit on behalf of
      the copyright owner. For the purposes of this definition, "submitted"
      means any form of electronic, verbal, or written communication sent
      to the Licensor or its representatives, including but not limited to
      communication on electronic mailing lists, source code control systems,
      and issue tracking systems that are managed by, or on behalf of, the
      Licensor for the purpose of discussing and improving the Work, but
      excluding communication that is conspicuously marked or otherwise
      designated in writing by the copyright owner as "Not a Contribution."

      "Contributor" shall mean Licensor and any individual or Legal Entity
      on behalf of whom a Contribution has been received by Licensor and
      subsequently incorporated within the Work.

   2. Grant of Copyright License. Subject to the terms and conditions of
      this License, each Contributor hereby grants to You a perpetual,
      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
      copyright license to reproduce, prepare Derivative Works of,
      publicly display, publicly perform, sublicense, and distribute the
      Work and such Derivative Works in Source or Object form.

   3. Grant of Patent License. Subject to the terms and conditions of
      this License, each Contributor hereby grants to You a perpetual,
      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
      (except as stated in this section) patent license to make, have made,
      use, offer to sell, sell, import, and otherwise transfer the Work,
      where such license applies only to those patent claims licensable
      by such Contributor that are necessarily infringed by their
      Contribution(s) alone or by combination of their Contribution(s)
      with the Work to which such Contribution(s) was submitted. If You
      institute patent litigation against any entity (including a
      cross-claim or counterclaim in a lawsuit) alleging that the Work
      or a Contribution incorporated within the Work constitutes direct
      or contributory patent infringement, then any patent licenses
      granted to You under this License for that Work shall terminate
      as of the date such litigation is filed.

   4. Redistribution. You may reproduce and distribute copies of the
      Work or Derivative Works thereof in any medium, with or without
      modifications, and in Source or Object form, provided that You
      meet the following conditions:

      (a) You must give any other recipients of the Work or
          Derivative Works a copy of this License; and

      (b) You must cause any modified files to carry prominent notices
          stating that You changed the files; and

      (c) You must retain, in the Source form of any Derivative Works
          that You distribute, all copyright, patent, trademark, and
          attribution notices from the Source form of the Work,
          excluding those notices that do not pertain to any part of
          the Derivative Works; and

      (d) If the Work includes a "NOTICE" text file as part of its
          distribution, then any Derivative Works that You distribute must
          include a readable copy of the attribution notices contained
          within such NOTICE file, excluding those notices that do not
          pertain to any part of the Derivative Works, in at least one
          of the following places: within a NOTICE text file distributed
          as part of the Derivative Works; within the Source form or
          documentation, if provided along with the Derivative Works; or,
          within a display generated by the Derivative Works, if and
          wherever such third-party notices normally appear. The contents
          of the NOTICE file are for informational purposes only and
          do not modify the License. You may add Your own attribution
          notices within Derivative Works that You distribute, alongside
          or as an addendum to the NOTICE text from the Work, provided
          that such additional attribution notices cannot be construed
          as modifying the License.

      You may add Your own copyright statement to Your modifications and
      may provide additional or different license terms and conditions
      for use, reproduction, or distribution of Your modifications, or
      for any such Derivative Works as a whole, provided Your use,
      reproduction, and distribution of the Work otherwise complies with
      the conditions stated in this License.

   5. Submission of Contributions. Unless You explicitly state otherwise,
      any Contribution intentionally submitted for inclusion in the Work
      by You to the Licensor shall be under the terms and conditions of
      this License, without any additional terms or conditions.
      Notwithstanding the above, nothing herein shall supersede or modify
      the terms of any separate license agreement you may have executed
      with Licensor regarding such Contributions.

   6. Trademarks. This License does not grant permission to use the trade
      names, trademarks, service marks, or product names of the Licensor,
      except as required for reasonable and customary use in describing the
      origin of the Work and reproducing the content of the NOTICE file.

   7. Disclaimer of Warranty. Unless required by applicable law or
      agreed to in writing, Licensor provides the Work (and each
      Contributor provides its Contributions) on an "AS IS" BASIS,
      WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
      implied, including, without limitation, any warranties or conditions
      of TITLE, NON-INFRINGEMENT, MERCHANTABILITY, or FITNESS FOR A
      PARTICULAR PURPOSE. You are solely responsible for determining the
      appropriateness of using or redistributing the Work and assume any
      risks associated with Your exercise of permissions under this License.

   8. Limitation of Liability. In no event and under no legal theory,
      whether in tort (including negligence), contract, or otherwise,
      unless required by applicable law (such as deliberate and grossly
      negligent acts) or agreed to in writing, shall any Contributor be
      liable to You for damages, including any direct, indirect, special,
      incidental, or consequential damages of any character arising as a
      result of this License or out of the use or inability to use the
      Work (including but not limited to damages for loss of goodwill,
      work stoppage, computer failure or malfunction, or any and all
      other commercial damages or losses), even if such Contributor
      has been advised of the possibility of such damages.

   9. Accepting Warranty or Additional Liability. While redistributing
      the Work or Derivative Works thereof, You may choose to offer,
      and charge a fee for, acceptance of support, warranty, indemnity,
      or other liability obligations and/or rights consistent with this
      License. However, in accepting such obligations, You may act only
      on Your own behalf and on Your sole responsibility, not on behalf
      of any other Contributor, and only if You agree to indemnify,
      defend, and hold each Contributor harmless for any liability
      incurred by, or claims asserted against, such Contributor by reason
      of your accepting any such warranty or additional liability.

   END OF TERMS AND CONDITIONS

   APPENDIX: How to apply the Apache License to your work.

      To apply the Apache License to your work, attach the following
      boilerplate notice, with the fields enclosed by brackets "[]"
      replaced with your own identifying information. (Don't include
      the brackets!)  The text should be enclosed in the appropriate
      comment syntax for the file format. We also recommend that a
      file or class name and description of purpose be included on the
      same "printed page" as the copyright notice for easier
      identification within third-party archives.

   Copyright [yyyy] [name of copyright owner]

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
````

### FILE: `licenses/golang-oauth2-BSD-3-Clause.txt`

```yaml
block_id: "GO-OIDC-SERVICE-TOKEN-BROKER:licenses/golang-oauth2-BSD-3-Clause.txt:v1"
operation: CREATE
provenance: VERBATIM
source: "exact pinned dependency LICENSE bytes; oidc-service-source-receipt.json"
license: "BSD-3-Clause"
sha256: "911f8f5782931320f5b8d1160a76365b83aea6447ee6c04fa6d5591467db9dad"
variables: []
secrets_allowed: false
```

````text
Copyright 2009 The Go Authors.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of Google LLC nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
````

