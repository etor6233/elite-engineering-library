# Connected catalog source authoring

## 1. Metadata

```yaml
pack_id: "GO-CONNECTED-CATALOG-AUTHORING"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: RECONSTRUCTIBLE
  admission: CONDITIONED
claim: "Role catalog authoring from initial model/variants/prices and PNG through existing approval/publication/storefront; local reference."
stacks: ["Go 1.26.8", "PostgreSQL 18.6", "Node.js 24.20.0 where frontend selected"]
compatible_with: ["MARKDOWN-COMPOSITOR0.3.0", "V402 existing owner closure"]
incompatible_with: ["unbound tenant/provider/account", "production certification inferred from fixtures", "implicit source or business-policy attribution"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: []
verified_at: "2026-09-12"
```

## 2. Applicability

LIBRARY_INFRASTRUCTURE T2804/J3. Existing catalog release, model/price services, shared approvals, OIDC/BFF and Next are required. Source form creates one model with1–16variants; frozen drafts are immutable, publication supports rollback to a reviewed version.

## 3. Architecture contract

Original model SQL/event extracted into its owner transaction method; original model and Commerce services are bound to one transaction through narrow repositories. Existing variant table, catalog command ledger and outbox; source remains draft/private/unknown homologation. Existing publication and review rules unchanged.

## 4. Exact file manifest

```text
CREATE cmd/electromobility-api/catalog_source_test.go
CREATE db/migrations/0076_catalog_source_authoring.down.sql
CREATE db/migrations/0076_catalog_source_authoring.up.sql
CREATE internal/catalogrelease/authoring_golden_test.go
CREATE internal/catalogrelease/source.go
CREATE internal/catalogrelease/source_fuzz_test.go
CREATE internal/platform/postgres/catalog_role_browser_integration_test.go
CREATE internal/platform/postgres/catalog_role_issuer_test.go
CREATE internal/platform/postgres/catalog_source.go
CREATE internal/platform/postgres/catalog_source_integration_test.go
```

## 5. Materialization blocks

### FILE: `cmd/electromobility-api/catalog_source_test.go`

```yaml
block_id: "GO-CONNECTED-CATALOG-AUTHORING:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "ef4ca0c90655a351d4fa044f4636a32e8051e115e3779f1420638b37316ad0e6"
variables: []
secrets_allowed: false
```

````go
package main

import (
	"context"
	"crypto/sha256"
	cr "elite.local/enterprise/internal/catalogrelease"
	"elite.local/enterprise/internal/platform/postgres"
	"encoding/hex"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"path/filepath"
	"testing"
)

func TestCatalogSourceHostMigration(t *testing.T) {
	raw := os.Getenv("PAYMENT_CONNECTED_DB_URL")
	if raw == "" {
		t.Skip("owned fixture required")
	}
	ctx := context.Background()
	pool, e := pgxpool.New(ctx, raw)
	if e != nil {
		t.Fatal(e)
	}
	defer pool.Close()
	profile := cr.Profile{Schema: "elite-catalog-publication/v1", TenantID: "ca6c9cb7-e548-4f19-a2e6-27d2aa2d4893", OrganizationID: "j3-store", Market: "AR", Currency: "ARS", Origin: "https://catalog.example.invalid", PolicyCode: cr.PolicyCode}
	body, _ := json.Marshal(profile)
	sum := sha256.Sum256(body)
	file := filepath.Join(t.TempDir(), "profile.json")
	if e = os.WriteFile(file, body, 0600); e != nil {
		t.Fatal(e)
	}
	config := map[string]string{"CATALOG_RELEASE_ENABLED": "true", "CATALOG_RELEASE_PROFILE_FILE": file, "CATALOG_RELEASE_PROFILE_SHA256": hex.EncodeToString(sum[:])}
	lookup := func(k string) string { return config[k] }
	price := postgres.NewCommerce(pool)
	if module, e := selectedCatalogReleaseModule(ctx, pool, price, lookup); e != nil || module == nil {
		t.Fatal("current source host rejected", e)
	}
	if _, e = pool.Exec(ctx, `alter table catalog.release_command rename constraint release_command_kind_check to source_kind_saved_fixture`); e != nil {
		t.Fatal(e)
	}
	defer func() {
		if _, e := pool.Exec(context.Background(), `alter table catalog.release_command rename constraint source_kind_saved_fixture to release_command_kind_check`); e != nil {
			t.Error(e)
		}
	}()
	if _, e = selectedCatalogReleaseModule(ctx, pool, price, lookup); e == nil {
		t.Fatal("missing source migration admitted")
	}
	t.Log("CATALOG_SOURCE_HOST_PASS exact source migration required; earlier host guards retained; no provider call")
}
````

### FILE: `db/migrations/0076_catalog_source_authoring.down.sql`

```yaml
block_id: "GO-CONNECTED-CATALOG-AUTHORING:file2:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "9948b7846dc68fabfa813ee8b34c220d60c31ed93c78adb8964e5488ff5b6c03"
variables: []
secrets_allowed: false
```

````sql
begin;
do $$begin
 if exists(select 1 from catalog.release_command where kind='source') then
  raise exception 'catalog source authoring history exists; downgrade refused';
 end if;
end$$;
alter table catalog.release_command drop constraint release_command_kind_check;
alter table catalog.release_command add constraint release_command_kind_check
 check(kind in('media','draft','review','publish'));
commit;
````

### FILE: `db/migrations/0076_catalog_source_authoring.up.sql`

```yaml
block_id: "GO-CONNECTED-CATALOG-AUTHORING:file3:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "556a8945db541bddaae34c25a95969000d42ff29c161567165d9a64fc03cf73c"
variables: []
secrets_allowed: false
```

````sql
begin;
alter table catalog.release_command drop constraint release_command_kind_check;
alter table catalog.release_command add constraint release_command_kind_check
 check(kind in('media','draft','review','publish','source'));
commit;
````

### FILE: `internal/catalogrelease/authoring_golden_test.go`

```yaml
block_id: "GO-CONNECTED-CATALOG-AUTHORING:file4:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "32dbbb75297c910bcbd3aadbc01e53264cc1ea6c8d352b783f01a5d95e8084a0"
variables: []
secrets_allowed: false
```

````go
package catalogrelease

import (
	"encoding/json"
	"os"
	"testing"
)

func TestCatalogRoleCommandGolden(t *testing.T) {
	raw, e := os.ReadFile("../../deploy/catalog/authoring-goldens.json")
	if e != nil {
		t.Fatal(e)
	}
	var cases []struct {
		Name      string          `json:"name"`
		Payload   json.RawMessage `json:"payload"`
		Canonical string          `json:"canonical"`
		SHA256    string          `json:"sha256"`
	}
	if e = json.Unmarshal(raw, &cases); e != nil {
		t.Fatal(e)
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			var value any
			switch c.Name {
			case "source":
				v := new(SourceRequest)
				e = json.Unmarshal(c.Payload, v)
				value = v
			case "review":
				v := new(ReviewRequest)
				e = json.Unmarshal(c.Payload, v)
				value = v
			case "publish":
				v := new(PublishRequest)
				e = json.Unmarshal(c.Payload, v)
				value = v
			case "media":
				v := map[string]string{}
				e = json.Unmarshal(c.Payload, &v)
				value = v
			default:
				t.Fatal("unknown golden")
			}
			if e != nil {
				t.Fatal(e)
			}
			got, h, e := Canonical(value)
			if e != nil || string(got) != c.Canonical || h != c.SHA256 {
				t.Fatalf("typed Go command differs: %v %s %s", e, h, got)
			}
		})
	}
}
````

### FILE: `internal/catalogrelease/source.go`

```yaml
block_id: "GO-CONNECTED-CATALOG-AUTHORING:file5:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "ae4ca94274ec0d92fca4c2a881f77cf4db4763995b5ef69b1be7c2beae7921d7"
variables: []
secrets_allowed: false
```

````go
// AUTHORED typed binding to the existing model and price services.
package catalogrelease

import (
	"elite.local/enterprise/internal/approval"
	"elite.local/enterprise/internal/electromobility"
	"encoding/json"
	"regexp"
	"time"
)

type SourceVariant struct {
	Code                 string          `json:"code"`
	DisplayName          string          `json:"display_name"`
	BatterySpecification json.RawMessage `json:"battery_specification"`
	AmountMinorUnits     int64           `json:"amount_minor_units,string"`
	TaxMode              string          `json:"tax_mode"`
}
type SourceRequest struct {
	CommandID  string                `json:"command_id"`
	Model      electromobility.Model `json:"model"`
	Variants   []SourceVariant       `json:"variants"`
	ValidFrom  time.Time             `json:"valid_from"`
	ValidUntil *time.Time            `json:"valid_until"`
}

func (s SourceRequest) Valid() bool {
	if !ValidID(s.CommandID) || s.Model.ID != "" || !ValidText(s.Model.DisplayName, 160) || len(s.Model.Code) > 64 || !sourceObject(s.Model.Specification) || s.ValidFrom.IsZero() || s.ValidUntil != nil && !s.ValidUntil.After(s.ValidFrom) || len(s.Variants) < 1 || len(s.Variants) > 16 {
		return false
	}
	seen := map[string]bool{}
	for _, v := range s.Variants {
		if len(v.Code) > 64 || !regexp.MustCompile(`^[a-z][a-z0-9]*(?:-[a-z0-9]+)*$`).MatchString(v.Code) || seen[v.Code] || !ValidText(v.DisplayName, 160) || !sourceObject(v.BatterySpecification) || v.AmountMinorUnits < 0 {
			return false
		}
		seen[v.Code] = true
	}
	return true
}
func sourceObject(raw json.RawMessage) bool {
	if len(raw) > 8192 {
		return false
	}
	var v map[string]any
	_, _, e := approval.CanonicalPayload(raw)
	return e == nil && json.Unmarshal(raw, &v) == nil && v != nil
}
````

### FILE: `internal/catalogrelease/source_fuzz_test.go`

```yaml
block_id: "GO-CONNECTED-CATALOG-AUTHORING:file6:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "75956799e0345b7ba89c8dc57c37c47f0c4196d70c87712310056472b1bb893b"
variables: []
secrets_allowed: false
```

````go
package catalogrelease

import (
	"elite.local/enterprise/internal/approval"
	"encoding/json"
	"testing"
)

func FuzzCatalogSourceBounds(f *testing.F) {
	f.Add([]byte(`{"command_id":"source","model":{"id":"","code":"bike","displayName":"Fixture","vehicleClass":"bicycle","specification":{"description":"bounded"}},"variants":[{"code":"bike-one","display_name":"Fixture","battery_specification":{},"amount_minor_units":"123456","tax_mode":"not-applicable"}],"valid_from":"2026-09-12T00:00:00Z","valid_until":null}`))
	f.Add([]byte(`{"command_id":"a","Command_ID":"b"}`))
	f.Add([]byte(`{}`))
	f.Fuzz(func(t *testing.T, raw []byte) {
		if len(raw) > 32768 {
			return
		}
		canon, _, e := approval.CanonicalPayload(raw)
		if e != nil {
			return
		}
		var value SourceRequest
		if json.Unmarshal(canon, &value) != nil {
			return
		}
		before, _ := json.Marshal(value)
		valid := value.Valid()
		after, _ := json.Marshal(value)
		if string(before) != string(after) {
			t.Fatal("validation mutated source")
		}
		if valid && (value.Model.ID != "" || len(value.Variants) < 1 || len(value.Variants) > 16 || value.ValidFrom.IsZero() || value.ValidUntil != nil && !value.ValidUntil.After(value.ValidFrom)) {
			t.Fatal("unbounded source admitted")
		}
		for _, v := range value.Variants {
			if valid && v.AmountMinorUnits < 0 {
				t.Fatal("negative amount")
			}
		}
	})
}
````

### FILE: `internal/platform/postgres/catalog_role_browser_integration_test.go`

```yaml
block_id: "GO-CONNECTED-CATALOG-AUTHORING:file7:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "9ea7ef9cffc958059b67741753f5534e8bca8e7ee31938f09cc40eeaea7bb9c3"
variables: []
secrets_allowed: false
```

````go
package postgres_test

// AUTHORED actual role browser on an empty catalog; only tenant/org are seeded.
import (
	"bytes"
	"context"
	cr "elite.local/enterprise/internal/catalogrelease"
	"elite.local/enterprise/internal/platform/httpapi"
	db "elite.local/enterprise/internal/platform/postgres"
	"encoding/base64"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"image"
	"image/png"
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

func TestCatalogAuthoringRoleBrowser(t *testing.T) {
	if os.Getenv("ELITE_CATALOG_ROLE_BROWSER") != "1" {
		t.Skip("explicit local browser fixture required")
	}
	raw := os.Getenv("PAYMENT_CONNECTED_DB_URL")
	if raw == "" {
		t.Fatal("owned database required")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	pool, e := pgxpool.New(ctx, raw)
	if e != nil {
		t.Fatal(e)
	}
	defer pool.Close()
	tenant := uuid.NewString()
	for _, q := range []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1,'role-catalog','Fixture','Fixture')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'role-org','role-org','Fixture','store')`,
	} {
		if _, e = pool.Exec(ctx, q, tenant); e != nil {
			t.Fatal(e)
		}
	}
	profile := cr.Profile{Schema: "elite-catalog-publication/v1", TenantID: tenant, OrganizationID: "role-org", Market: "AR", Currency: "ARS", Origin: "https://catalog.example.invalid", PolicyCode: cr.PolicyCode}
	store, e := db.NewCatalogRelease(pool, profile, db.NewCommerce(pool))
	if e != nil {
		t.Fatal(e)
	}
	verifier, token := catalogRoleIssuer(t)
	identities := map[string]any{}
	for _, name := range []string{"maker", "reviewer", "reader", "foreign-org", "unprivileged"} {
		permissions := []string{"catalog:read"}
		orgs := []string{"role-org"}
		if name == "maker" {
			permissions = append(permissions, "catalog:draft")
		}
		if name == "reviewer" {
			permissions = append(permissions, "catalog:publish", "catalog:review:legal", "catalog:review:technical", "catalog:review:media", "catalog:review:publication")
		}
		if name == "foreign-org" {
			orgs = []string{"foreign"}
		}
		if name == "unprivileged" {
			permissions = []string{}
		}
		identities[name] = map[string]any{"subject": name, "tenantId": tenant, "permissions": permissions, "organizations": orgs, "accessToken": token(name, tenant, permissions, orgs)}
	}
	mux := http.NewServeMux()
	httpapi.CatalogReleaseModule{Service: store}.Register(mux, verifier)
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
	web, node := os.Getenv("ELITE_WEB_ROOT"), os.Getenv("ELITE_NODE_BIN")
	if !filepath.IsAbs(web) || !filepath.IsAbs(node) {
		t.Fatal("fixed runtime/root required")
	}
	listener, e := net.Listen("tcp", "127.0.0.1:0")
	if e != nil {
		t.Fatal(e)
	}
	address := listener.Addr().String()
	listener.Close()
	target, _ := url.Parse("http://" + address)
	edge := httptest.NewTLSServer(httputil.NewSingleHostReverseProxy(target))
	defer edge.Close()
	env := []string{}
	for _, v := range os.Environ() {
		k := strings.ToUpper(strings.SplitN(v, "=", 2)[0])
		if strings.HasPrefix(k, "ELITE_") || strings.HasPrefix(k, "CATALOG_") || strings.HasPrefix(k, "PUBLIC_") || strings.HasPrefix(k, "ENTERPRISE_") || strings.Contains(k, "DATABASE_URL") || k == "APP_BASE_URL" || k == "AUTH_SESSION_SECRET" || k == "BUSINESS_CONFIG_FILE" {
			continue
		}
		env = append(env, v)
	}
	identitiesJSON, _ := json.Marshal(identities)
	var picture bytes.Buffer
	if e = png.Encode(&picture, image.NewNRGBA(image.Rect(0, 0, 2, 2))); e != nil {
		t.Fatal(e)
	}
	env = append(env, "ELITE_CATALOG_ROLE_BROWSER=1", "ELITE_WORKSPACE_E2E=enabled", "ELITE_CATALOG_IDENTITIES="+string(identitiesJSON), "ELITE_WEB_ROOT="+web, "ELITE_BASE_URL="+edge.URL,
		"ELITE_CATALOG_PNG="+base64.StdEncoding.EncodeToString(picture.Bytes()), "ELITE_CATALOG_VALID_FROM="+time.Now().UTC().Add(-time.Hour).Format("2006-01-02T15:04"),
		"APP_BASE_URL="+edge.URL, "ENTERPRISE_API_BASE_URL="+api.URL, "AUTH_SESSION_SECRET=synthetic-catalog-role-"+uuid.NewString(), "BUSINESS_CONFIG_FILE=catalog.authoring.fixture.json",
		"CATALOG_RELEASE_ENABLED=true", "PUBLIC_SITE_ORIGIN=https://catalog.example.invalid", "PUBLIC_INDEXING_ENABLED=1", "ENTERPRISE_TENANT_CODE=role-catalog", "ENTERPRISE_ORGANIZATION_CODE=role-org", "NEXT_TELEMETRY_DISABLED=1")
	artifacts, e := os.MkdirTemp(web, "catalog-role-artifacts-")
	if e != nil {
		t.Fatal(e)
	}
	log, e := os.Create(filepath.Join(artifacts, "next.log"))
	if e != nil {
		t.Fatal(e)
	}
	defer log.Close()
	server := exec.CommandContext(ctx, node, filepath.Join(web, "node_modules", "next", "dist", "bin", "next"), "start", "--hostname", "127.0.0.1", "--port", strings.Split(address, ":")[1])
	server.Dir, server.Env, server.Stdout, server.Stderr = web, env, log, log
	if e = server.Start(); e != nil {
		t.Fatal(e)
	}
	defer func() { server.Process.Kill(); server.Wait() }()
	client := &http.Client{Timeout: time.Second}
	ready := false
	for deadline := time.Now().Add(20 * time.Second); time.Now().Before(deadline); {
		r, e := client.Get(target.String() + "/icon.svg")
		if e == nil {
			r.Body.Close()
			if r.StatusCode == 200 {
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
		t.Fatal("Next startup", artifacts)
	}
	gate := filepath.Join(web, "microsoft_playwright_browser_gate")
	command := exec.CommandContext(ctx, node, filepath.Join(gate, "node_modules", "@playwright", "test", "cli.js"), "test", "tests/catalog-authoring-connected.spec.mjs", "--project=chromium-desktop", "--timeout=90000", "--output="+filepath.Join(artifacts, "browsers"))
	command.Dir, command.Env = gate, env
	output, e := command.CombinedOutput()
	if x := os.WriteFile(filepath.Join(artifacts, "browser.log"), output, 0600); x != nil {
		t.Fatal(x)
	}
	if e != nil || !strings.Contains(string(output), "1 passed") {
		t.Fatalf("catalog role browser %v\n%s\nartifacts=%s", e, output, artifacts)
	}
	var actual string
	e = pool.QueryRow(ctx, `select jsonb_build_array(
 (select count(*) from catalog.vehicle_model where tenant_id=$1),
 (select count(*) from catalog.vehicle_variant where tenant_id=$1),
 (select count(*) from pricing.price_book where tenant_id=$1),
 (select count(*) from catalog.release_media where tenant_id=$1),
 (select count(*) from catalog.release_draft where tenant_id=$1),
 (select count(*) from catalog.release_review_action where tenant_id=$1),
 (select count(*) from catalog.release_publication where tenant_id=$1),
 (select count(*) from catalog.release_command where tenant_id=$1))::text`, tenant).Scan(&actual)
	if e != nil || actual != "[2, 2, 5, 1, 2, 8, 3, 16]" {
		t.Fatal("durable effects", actual, e)
	}
	mu.Lock()
	countJSON, _ := json.Marshal(counts)
	total := 0
	for _, n := range counts {
		total += n
	}
	mu.Unlock()
	if total != 16 {
		t.Fatal("unexpected backend writes", total)
	}
	if e = os.WriteFile(filepath.Join(artifacts, "api-post-counts.json"), countJSON, 0600); e != nil {
		t.Fatal(e)
	}
	t.Logf("CATALOG_ROLE_BROWSER_PASS source_to_publication_actual=true JWE_RS256_JWKS=true source_and_publish_response_loss_GET_only=true two_models_two_snapshots_eight_reviews_three_publications=true backend_POSTs=16 artifacts=%s", artifacts)
}
````

### FILE: `internal/platform/postgres/catalog_role_issuer_test.go`

```yaml
block_id: "GO-CONNECTED-CATALOG-AUTHORING:file8:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "98315a2e7285174099462bbefa33c1cca713a2bec213d27df9ca9e8b083f3ccb"
variables: []
secrets_allowed: false
```

````go
package postgres_test

// AUTHORED synthetic RS256/JWKS fixture reused from admitted local browser tests.
import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"elite.local/enterprise/internal/platform/identity"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func catalogRoleIssuer(t *testing.T) (identity.Verifier, func(string, string, []string, []string) string) {
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
````

### FILE: `internal/platform/postgres/catalog_source.go`

```yaml
block_id: "GO-CONNECTED-CATALOG-AUTHORING:file9:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "64aaf4125c7928b98e7ac4c5dab2925778aea493a56ef703416b14bfc9e898a0"
variables: []
secrets_allowed: false
```

````go
package postgres

// AUTHORED transaction adapters; model and pricing validation/SQL remain in
// their original services. No alternate pricing or publication algorithm.
import (
	"context"
	cr "elite.local/enterprise/internal/catalogrelease"
	"elite.local/enterprise/internal/commerce"
	"elite.local/enterprise/internal/electromobility"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/randomid"
	"github.com/jackc/pgx/v5"
)

type catalogModelTx struct {
	*Electromobility
	tx pgx.Tx
}

func (r catalogModelTx) CreateModel(ctx context.Context, tenant, event string, m electromobility.Model) error {
	return r.createModelTx(ctx, r.tx, tenant, event, m)
}

type catalogPriceTx struct {
	*Commerce
	tx pgx.Tx
}

func (r catalogPriceTx) CreatePriceBook(ctx context.Context, tenant, event string, b commerce.PriceBook) error {
	return r.createPriceBookTx(ctx, r.tx, tenant, event, b)
}

func (s *CatalogRelease) CreateSource(ctx context.Context, p identity.Principal, in cr.SourceRequest) (cr.Receipt, error) {
	if !in.Valid() {
		return cr.Receipt{}, cr.ErrInvalid
	}
	_, hash, e := cr.Canonical(in)
	if e != nil {
		return cr.Receipt{}, e
	}
	tx, e := s.begin(ctx, p, "catalog:draft")
	if e != nil {
		return cr.Receipt{}, e
	}
	defer tx.Rollback(ctx)
	if r, yes, e := s.replay(ctx, tx, p, in.CommandID, "source", hash); e != nil || yes {
		return r, e
	}
	ids := randomid.Generator{}
	model, e := electromobility.NewService(catalogModelTx{NewElectromobility(s.pool), tx}, ids).CreateModel(ctx, p.TenantID, in.Model)
	if e != nil {
		return cr.Receipt{}, e
	}
	entries := []commerce.PriceEntry{}
	variantIDs := []string{}
	for _, v := range in.Variants {
		id := ids.New()
		_, e = tx.Exec(ctx, `insert into catalog.vehicle_variant(tenant_id,variant_id,model_id,variant_code,display_name,battery_specification,lifecycle_state)
  values($1,$2,$3,$4,$5,$6,'draft')`, p.TenantID, id, model.ID, v.Code, v.DisplayName, v.BatterySpecification)
		if e != nil {
			return cr.Receipt{}, e
		}
		variantIDs = append(variantIDs, id)
		entries = append(entries, commerce.PriceEntry{VariantID: id, AmountMinorUnits: v.AmountMinorUnits, TaxMode: v.TaxMode})
	}
	book, e := commerce.NewService(catalogPriceTx{s.price, tx}, ids).CreatePriceBook(ctx, p.TenantID, commerce.PriceBook{Market: s.profile.Market, Currency: s.profile.Currency, ValidFrom: in.ValidFrom, ValidUntil: in.ValidUntil, Entries: entries})
	if e != nil {
		return cr.Receipt{}, e
	}
	r := cr.Receipt{CommandID: in.CommandID, Actor: p.Subject, Kind: "source", ResourceID: model.ID, RequestSHA256: hash, SourcePriceBookID: book.ID, SourceVariantIDs: variantIDs}
	if e = s.save(ctx, tx, p, r); e != nil {
		return cr.Receipt{}, e
	}
	return r, tx.Commit(ctx)
}

func (s *CatalogRelease) Current(ctx context.Context, p identity.Principal) (cr.Publication, error) {
	if !s.authorized(p, "catalog:read") {
		return cr.Publication{}, cr.ErrInvalid
	}
	return s.Public(ctx)
}
````

### FILE: `internal/platform/postgres/catalog_source_integration_test.go`

```yaml
block_id: "GO-CONNECTED-CATALOG-AUTHORING:file10:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "b10f16c90289f44aca867224825057dd01e9fdf75e7e9e79ef939b4978fe4fe1"
variables: []
secrets_allowed: false
```

````go
package postgres_test

// AUTHORED focused regression for the newly materializable source write only.
// Publication, four reviews and storefront will be exercised by the role browser.
import (
	"context"
	cr "elite.local/enterprise/internal/catalogrelease"
	"elite.local/enterprise/internal/electromobility"
	"elite.local/enterprise/internal/platform/identity"
	db "elite.local/enterprise/internal/platform/postgres"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"sync"
	"testing"
	"time"
)

func TestCatalogSourceAuthoringAtomic(t *testing.T) {
	raw := os.Getenv("PAYMENT_CONNECTED_DB_URL")
	if raw == "" {
		t.Skip("owned fixture required")
	}
	ctx := context.Background()
	pool, e := pgxpool.New(ctx, raw)
	if e != nil {
		t.Fatal(e)
	}
	defer pool.Close()
	tenant := uuid.NewString()
	for _, q := range []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1,'source-'||($1::uuid)::text,'Fixture','Fixture')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'source-org','source-org','Fixture','store')`,
	} {
		if _, e = pool.Exec(ctx, q, tenant); e != nil {
			t.Fatal(e)
		}
	}
	profile := cr.Profile{Schema: "elite-catalog-publication/v1", TenantID: tenant, OrganizationID: "source-org", Market: "AR", Currency: "ARS", Origin: "https://catalog.example.invalid", PolicyCode: cr.PolicyCode}
	store, e := db.NewCatalogRelease(pool, profile, db.NewCommerce(pool))
	if e != nil {
		t.Fatal(e)
	}
	actor := identity.Principal{TenantID: tenant, Subject: "source-maker", Organizations: map[string]struct{}{"source-org": {}}, Permissions: map[string]struct{}{"catalog:draft": {}, "catalog:read": {}}}
	in := cr.SourceRequest{CommandID: "source-create", Model: electromobility.Model{Code: "created-through-source", DisplayName: "Draft authored model", VehicleClass: "bicycle", Specification: json.RawMessage(`{"description":"Fixture"}`)}, Variants: []cr.SourceVariant{{Code: "source-variant", DisplayName: "Reference variant", BatterySpecification: json.RawMessage(`{}`), AmountMinorUnits: 123456, TaxMode: "not-applicable"}}, ValidFrom: time.Now().UTC().Truncate(time.Second).Add(-time.Hour)}
	const n = 4
	var wg sync.WaitGroup
	got := make(chan cr.Receipt, n)
	errs := make(chan error, n)
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); r, e := store.CreateSource(ctx, actor, in); got <- r; errs <- e }()
	}
	wg.Wait()
	close(got)
	close(errs)
	for e := range errs {
		if e != nil {
			t.Fatal(e)
		}
	}
	var result cr.Receipt
	newCount := 0
	for r := range got {
		if !r.Replay {
			newCount++
		}
		if result.ResourceID != "" && result.ResourceID != r.ResourceID {
			t.Fatal("different model")
		}
		result = r
	}
	if newCount != 1 || len(result.SourceVariantIDs) != 1 || result.SourcePriceBookID == "" {
		t.Fatal("duplicate or missing source receipt", newCount, result)
	}
	snapshot := func() string {
		var raw string
		e := pool.QueryRow(ctx, `select jsonb_build_array(
  (select count(*) from catalog.vehicle_model where tenant_id=$1),
  (select count(*) from catalog.vehicle_variant where tenant_id=$1),
  (select count(*) from pricing.price_book where tenant_id=$1),
  (select count(*) from pricing.price_book_entry where tenant_id=$1),
  (select count(*) from catalog.release_command where tenant_id=$1),
  (select count(*) from platform.outbox_event where tenant_id=$1))::text`, tenant).Scan(&raw)
		if e != nil {
			t.Fatal(e)
		}
		return raw
	}
	before := snapshot()
	if before != "[1, 1, 1, 1, 1, 3]" {
		t.Fatal(before)
	}
	recovered, e := store.CommandReceipt(ctx, actor, in.CommandID)
	if e != nil || recovered.ResourceID != result.ResourceID || recovered.RequestSHA256 != result.RequestSHA256 {
		t.Fatal("read recovery", e)
	}
	other := actor
	other.Subject = "other"
	if _, e = store.CreateSource(ctx, other, in); e == nil {
		t.Fatal("actor replay")
	}
	if _, e = store.CommandReceipt(ctx, other, in.CommandID); e == nil {
		t.Fatal("receipt leak")
	}
	edited := in
	edited.Model.DisplayName = "Changed"
	if _, e = store.CreateSource(ctx, actor, edited); e == nil {
		t.Fatal("hash conflict")
	}
	foreign := actor
	foreign.Organizations = map[string]struct{}{"foreign": {}}
	if _, e = store.CreateSource(ctx, foreign, in); e == nil {
		t.Fatal("org bypass")
	}
	if snapshot() != before {
		t.Fatal("rejections wrote")
	}
	var hidden bool
	if e = pool.QueryRow(ctx, `select m.lifecycle_state='draft' and not m.publicly_visible and v.lifecycle_state='draft' and v.homologation_state='unknown' and b.status='draft'
 from catalog.vehicle_model m join catalog.vehicle_variant v using(tenant_id,model_id) join pricing.price_book_entry x using(tenant_id,variant_id) join pricing.price_book b using(tenant_id,price_book_id) where m.tenant_id=$1`, tenant).Scan(&hidden); e != nil || !hidden {
		t.Fatal("source promoted without review", e)
	}
	_, e = pool.Exec(ctx, `create function catalog.test_source_outbox_fail()returns trigger language plpgsql as $$begin if new.aggregate_type='catalog-release-command' and new.event_type='catalog-release.source' then raise exception 'fixture outbox failure';end if;return new;end$$;
 create trigger test_source_outbox_fail before insert on platform.outbox_event for each row execute function catalog.test_source_outbox_fail()`)
	if e != nil {
		t.Fatal(e)
	}
	broken := in
	broken.CommandID = "source-rollback"
	broken.Model.Code = "rolled-back-model"
	broken.Variants = append([]cr.SourceVariant(nil), in.Variants...)
	broken.Variants[0].Code = "rolled-back-variant"
	if _, e = store.CreateSource(ctx, actor, broken); e == nil {
		t.Fatal("outbox accepted")
	}
	if snapshot() != before {
		t.Fatal("source/model/price/outbox partial commit")
	}
	if _, e = pool.Exec(ctx, `drop trigger test_source_outbox_fail on platform.outbox_event;drop function catalog.test_source_outbox_fail()`); e != nil {
		t.Fatal(e)
	}
	reopened, e := db.NewCatalogRelease(pool, profile, db.NewCommerce(pool))
	if e != nil {
		t.Fatal(e)
	}
	same, e := reopened.CommandReceipt(ctx, actor, in.CommandID)
	if e != nil || same.SourcePriceBookID != result.SourcePriceBookID {
		t.Fatal("restart receipt", e)
	}
	t.Log("CATALOG_SOURCE_AUTHORING_PASS model_variant_price_original_owners=true concurrent_one_new_three_replay=true command_actor_hash_org_bound=true all6tables_atomic=true source_draft_no_publication=true durable_recovery=true")
}
````

## 6. Configuration surface

docs/CATALOG_ROLE_AUTHORING_REFERENCE.md; opt-in features.catalog_editor with catalog:read and independent draft/review/publish permissions. Go host also requires migration0076. Exact existing catalog profile and provider conditions retained.

## 7. Dependency bill

20new/10changed AUTHORED unavoidable transport/form/configuration/transaction/test glue; no new upstream or dependency. Existing ADAPTED business sources and pinned runtimes/licenses retained; no corporate attribution.

## 8. Apply order

Select both authoring packs with GO-CONNECTED-CATALOG-PUBLICATION and the full reference owner closure. Migration0076 follows0075 and preserves old kinds. New source kind makes authoring companion required with the updated publication pack. Web pack owns shared command golden fixture used by the Go test.

## 9. Verification

Actual empty-catalog role browser:2models/2variants/5pricebooks/1media/2snapshots/8reviews/3publications/16commands. Source and publication response losses recover via GET with zero extra POST. Four concurrent source requests create one aggregate;6table rollback and durable receipt recovery.14focused web tests, strict types/build,4Go/TS command goldens, source-only finite fuzz. Host requires migration; populated down refuses16commands/3publications/2models, empty73migrations/down/up PASS.

## 10. Reconstruction evidence

CATALOG_ROLE_AUTHORING_RELEASE_V402.md/json binds exact source/delta, all RED/PASS receipts, reconstruction and notices. Supply/warranty/J5/CMS/KPIs/private locale remain explicit T2804 owners. No whole T2804, TEST02 or overall readiness promotion.

