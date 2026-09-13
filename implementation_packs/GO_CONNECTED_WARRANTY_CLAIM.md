# Connected sold warranty and repair claim reference

## 1. Metadata

```yaml
pack_id: "GO-CONNECTED-WARRANTY-CLAIM"
pack_version: "0.2.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Immutable sold warranty terms and source-derived coverage connected to existing appointment/service/approval/FIFO stock owners, repair quality/customer acceptance and factory acknowledgement; local reference profile only."
stacks: ["Go 1.26.8", "PostgreSQL 18.6", "Node.js 24.20.0 where frontend selected"]
compatible_with: ["MARKDOWN-COMPOSITOR0.3.0", "V402 complete franchise dependency closure", "GO-BC-WARRANTY-COVERAGE-ADAPTER0.1.0"]
incompatible_with: ["unbound tenant/provider/account", "production certification inferred from fixtures", "implicit source or business-policy attribution"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: []
verified_at: "2026-09-12"
```

## 2. Applicability

Select the admitted BC date adapter and existing franchise commerce/payment/handover, service, inventory and approval owners. This is the J1 warranty activation and J4 service reference. The reference policy is explicit and hash-bound; no law, tax, cash settlement or business duration is inferred. Role UI is a separate T2804 owner.

## 3. Architecture contract

One transaction composes existing service, approval and stock writers with immutable warranty bindings and outbox. Sold terms/quote acknowledgement precede order acceptance. Activation reads the sold profile, not later config. Distinct humans approve bound work; existing reservations hold parts and existing FIFO issues them. Deferred constraints reject lease expiry at COMMIT. Work, quality, customer acceptance and factory reconciliation reference the same case versions/receipts. Recovery is actor/request-bound. Internal settlement creates no payment. One bounded work plan per case; additional parts use a follow-up case.

## 4. Exact file manifest

```text
CREATE cmd/electromobility-api/warranty.go
CREATE cmd/electromobility-api/warranty_test.go
CREATE cmd/warranty-profile/main.go
CREATE cmd/warranty-profile/main_test.go
CREATE db/migrations/0069_warranty_terms_activation.down.sql
CREATE db/migrations/0069_warranty_terms_activation.up.sql
CREATE db/migrations/0070_warranty_claim.down.sql
CREATE db/migrations/0070_warranty_claim.up.sql
CREATE deploy/warranty/profile.reference.json
CREATE docs/WARRANTY_REFERENCE.md
CREATE docs/provenance/BC_WARRANTY_NOTICES.md
CREATE internal/approval/warranty_test.go
CREATE internal/platform/httpapi/warranty.go
CREATE internal/platform/httpapi/warranty_test.go
CREATE internal/platform/postgres/warranty_claim.go
CREATE internal/platform/postgres/warranty_claim_integration_test.go
CREATE internal/platform/postgres/warranty_claim_work.go
CREATE internal/platform/postgres/warranty_connected_integration_test.go
CREATE internal/platform/postgres/warranty_expiry_integration_test.go
CREATE internal/platform/postgres/warranty_http_integration_test.go
CREATE internal/platform/postgres/warranty_terms.go
CREATE internal/warrantyclaim/claim.go
CREATE internal/warrantyclaim/contract.go
CREATE internal/warrantyclaim/contract_test.go
CREATE internal/warrantyclaim/profile_fuzz_test.go
CREATE internal/warrantyclaim/terms.go
```

## 5. Materialization blocks

### FILE: `cmd/electromobility-api/warranty.go`

```yaml
block_id: "GO-CONNECTED-WARRANTY-CLAIM:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "87ed190936a4627f958e488851f566f618e151a7bd4e5f2ce5fd22600498aef4"
variables: []
secrets_allowed: false
```

````go
package main

// AUTHORED optional local activation; credentials are neither read nor needed.
import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"

	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/platform/postgres"
	wc "elite.local/enterprise/internal/warrantyclaim"
	"github.com/jackc/pgx/v5/pgxpool"
)

var errWarrantyConfiguration = errors.New("warranty activation is invalid")

func init() { warrantyModuleFactory = selectedWarrantyModule }
func selectedWarrantyModule(ctx context.Context, pool *pgxpool.Pool, lookup func(string) string) (httpapi.EnterpriseModule, error) {
	if lookup == nil {
		return nil, errWarrantyConfiguration
	}
	switch lookup("WARRANTY_ENABLED") {
	case "", "false":
		return nil, nil
	case "true":
	default:
		return nil, errWarrantyConfiguration
	}
	if pool == nil {
		return nil, errWarrantyConfiguration
	}
	path := lookup("WARRANTY_PROFILE_FILE")
	if !filepath.IsAbs(path) {
		return nil, errWarrantyConfiguration
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, errWarrantyConfiguration
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil || !info.Mode().IsRegular() || info.Size() < 1 || info.Size() > 32768 {
		return nil, errWarrantyConfiguration
	}
	raw, err := io.ReadAll(io.LimitReader(file, 32769))
	if err != nil {
		return nil, errWarrantyConfiguration
	}
	profile, err := wc.LoadProfile(raw, lookup("WARRANTY_PROFILE_SHA256"))
	if err != nil {
		return nil, errWarrantyConfiguration
	}
	tenant, org := profile.Scope()
	if tenant != lookup("WARRANTY_TENANT_ID") || org != lookup("WARRANTY_ORGANIZATION_ID") {
		return nil, errWarrantyConfiguration
	}
	store, err := postgres.NewWarranty(pool, profile)
	if err != nil {
		return nil, errWarrantyConfiguration
	}
	// Fail before mounting routes if the selected composition lacks the schema
	// or its immutable-fact/connected-command guards.
	var ready bool
	err = pool.QueryRow(ctx, `select to_regclass('service_ops.warranty_offer') is not null and to_regclass('service_ops.warranty_activation') is not null
 and to_regclass('service_ops.warranty_claim_step') is not null and to_regclass('service_ops.warranty_part_reservation') is not null
 and (select count(*) from pg_trigger where not tgisinternal and tgenabled in ('O','A') and
 (tgrelid=to_regclass('service_ops.service_case') and tgname='warranty_claim_case_guard'
 or tgrelid=to_regclass('sales.quotation') and tgname='warranty_quote_acceptance_guard'
 or tgrelid=to_regclass('service_ops.warranty_work_plan') and tgname='warranty_plan_commit_expiry'
 or tgrelid=to_regclass('service_ops.warranty_claim_step') and tgname='warranty_approval_commit_expiry'))=4`).Scan(&ready)
	if err != nil || !ready {
		return nil, errWarrantyConfiguration
	}
	return httpapi.WarrantyModule{Service: store}, nil
}
````

### FILE: `cmd/electromobility-api/warranty_test.go`

```yaml
block_id: "GO-CONNECTED-WARRANTY-CLAIM:file2:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "7991b8874de9b6a442ecc77a8d29336a67ff7b31f7eda8ce4a4e8094efaa8fc3"
variables: []
secrets_allowed: false
```

````go
package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	wc "elite.local/enterprise/internal/warrantyclaim"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestWarrantyHostDisabledReadsNoProfile(t *testing.T) {
	for _, enabled := range []string{"", "false"} {
		m, err := selectedWarrantyModule(context.Background(), nil, func(key string) string {
			if key != "WARRANTY_ENABLED" {
				t.Fatal("disabled module read configuration", key)
			}
			return enabled
		})
		if err != nil || m != nil {
			t.Fatal(m, err)
		}
	}
	if _, err := selectedWarrantyModule(context.Background(), nil, func(string) string { return "TRUE" }); err == nil {
		t.Fatal("untyped activation admitted")
	}
}
func TestWarrantyHostProfileAndSchemaActivation(t *testing.T) {
	rawURL := os.Getenv("PAYMENT_CONNECTED_DB_URL")
	if rawURL == "" {
		t.Skip("owned reference PostgreSQL URL is required")
	}
	parsed, err := url.Parse(rawURL)
	if err != nil || parsed.Hostname() != "127.0.0.1" || !strings.HasPrefix(parsed.Path, "/elite_payment_connected_") {
		t.Fatal("fixture database outside owned scope")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, rawURL)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	d := wc.ProfileDocument{Schema: "elite-warranty-profile/v1", Scope: "MATERIALIZED_PROFILE", Algorithm: "bc-inclusive-fixed-terms", AlgorithmRevision: 1, TenantID: "aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa", OrganizationID: "store", FactoryOrganizationID: "factory", PolicyID: "synthetic", TermsVersion: "fixture-v1", TermsText: "Fixture profile only; no legal terms inferred.", BusinessTimeZone: "UTC", PartsDurationDays: 365, LaborDurationDays: 365, WorkReservationSeconds: 3600, FaultExclusions: []string{}, Settlement: "INTERNAL_WARRANTY_SERVICE_ACKNOWLEDGEMENT", AuthorityReference: "fixture-authority", DecisionReference: "fixture-decision"}
	raw, err := json.Marshal(d)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "profile.json")
	if err = os.WriteFile(path, raw, 0600); err != nil {
		t.Fatal(err)
	}
	values := map[string]string{"WARRANTY_ENABLED": "true", "WARRANTY_PROFILE_FILE": path, "WARRANTY_PROFILE_SHA256": wc.SHA(raw), "WARRANTY_TENANT_ID": d.TenantID, "WARRANTY_ORGANIZATION_ID": d.OrganizationID}
	lookup := func(key string) string { return values[key] }
	module, err := selectedWarrantyModule(ctx, pool, lookup)
	if err != nil || module == nil {
		t.Fatal("selected module", err)
	}
	mux := http.NewServeMux()
	module.Register(mux, nil)
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, httptest.NewRequest("GET", "/v1/franchise/warranty/profile?organization_id=store", nil))
	if response.Code != 401 {
		t.Fatal("route not protected/mounted", response.Code)
	}
	for _, v := range []struct{ key, value string }{{"WARRANTY_PROFILE_SHA256", strings.Repeat("f", 64)}, {"WARRANTY_PROFILE_FILE", "relative.json"}, {"WARRANTY_TENANT_ID", "bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbbbb"}, {"WARRANTY_ORGANIZATION_ID", "other"}} {
		previous := values[v.key]
		values[v.key] = v.value
		_, err := selectedWarrantyModule(ctx, pool, lookup)
		values[v.key] = previous
		if err == nil {
			t.Fatal("mismatched configuration admitted", v.key)
		}
	}
	if err = os.WriteFile(path, append(raw, []byte(strings.Repeat(" ", 32769))...), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err = selectedWarrantyModule(ctx, pool, lookup); err == nil {
		t.Fatal("oversized profile admitted")
	}
	if err = os.WriteFile(path, raw, 0600); err != nil {
		t.Fatal(err)
	}
	// This is the owned disposable fixture database; restore the guard even if
	// the negative activation assertion fails.
	if _, err = pool.Exec(ctx, `alter table service_ops.service_case disable trigger warranty_claim_case_guard`); err != nil {
		t.Fatal(err)
	}
	func() {
		defer func() {
			if _, err := pool.Exec(ctx, `alter table service_ops.service_case enable trigger warranty_claim_case_guard`); err != nil {
				t.Fatal("restore fixture guard", err)
			}
		}()
		if _, err := selectedWarrantyModule(ctx, pool, lookup); err == nil {
			t.Fatal("disabled command guard admitted")
		}
	}()
	t.Log("WARRANTY_HOST_PG_PASS exact_profile_scope_size_schema_guard=true routes_protected=true no_live_credentials=true")
}
````

### FILE: `cmd/warranty-profile/main.go`

```yaml
block_id: "GO-CONNECTED-WARRANTY-CLAIM:file3:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "dba3f721bc0594dff844153f22d7cb5407ecd399b4546f9fb28a913c13af66f1"
variables: []
secrets_allowed: false
```

````go
// AUTHORED profile materialization glue. The same Go loader used by the host
// validates all source/override data before creating an absent destination.
package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	wc "elite.local/enterprise/internal/warrantyclaim"
)

type activation struct {
	Schema                string            `json:"schema"`
	Scope                 string            `json:"scope"`
	ProfileFile           string            `json:"profile_file"`
	ProfileSHA256         string            `json:"profile_sha256"`
	SourceSHA256          string            `json:"source_sha256"`
	TenantID              string            `json:"tenant_id"`
	OrganizationID        string            `json:"organization_id"`
	FactoryOrganizationID string            `json:"factory_organization_id"`
	Overrides             map[string]string `json:"overrides"`
	ProductionAuthorized  bool              `json:"production_authorized"`
}

func materialize(args []string) (activation, error) {
	var empty activation
	flags := flag.NewFlagSet("warranty-profile", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	input := flags.String("input", "", "existing profile JSON")
	expected := flags.String("sha256", "", "exact source profile SHA256")
	output := flags.String("output", "", "absent output directory")
	tenant := flags.String("tenant-id", "", "explicit tenant override")
	org := flags.String("organization-id", "", "explicit store override")
	factory := flags.String("factory-organization-id", "", "explicit factory override")
	if err := flags.Parse(args); err != nil {
		return empty, err
	}
	if flags.NArg() != 0 || *input == "" || *output == "" || !wc.ValidSHA(*expected) {
		return empty, errors.New("input, sha256 and absent output are required")
	}
	f, err := os.Open(*input)
	if err != nil {
		return empty, err
	}
	info, err := f.Stat()
	if err != nil || !info.Mode().IsRegular() || info.Size() < 1 || info.Size() > 32768 {
		f.Close()
		return empty, errors.New("source profile must be a bounded regular file")
	}
	raw, err := io.ReadAll(io.LimitReader(f, 32769))
	closeErr := f.Close()
	if err != nil {
		return empty, err
	}
	if closeErr != nil {
		return empty, closeErr
	}
	source, err := wc.LoadProfile(raw, *expected)
	if err != nil {
		return empty, err
	}
	doc := source.Document()
	overrides := map[string]string{}
	if *tenant != "" {
		doc.TenantID = *tenant
		overrides["tenant_id"] = *tenant
	}
	if *org != "" {
		doc.OrganizationID = *org
		overrides["organization_id"] = *org
	}
	if *factory != "" {
		doc.FactoryOrganizationID = *factory
		overrides["factory_organization_id"] = *factory
	}
	profileBytes, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return empty, err
	}
	profileBytes = append(profileBytes, '\n')
	profileHash := wc.SHA(profileBytes)
	selected, err := wc.LoadProfile(profileBytes, profileHash)
	if err != nil {
		return empty, err
	}
	boundTenant, boundOrg := selected.Scope()
	value := activation{Schema: "elite-warranty-activation/v1", Scope: "LIBRARY_INFRASTRUCTURE", ProfileFile: "profile.json", ProfileSHA256: profileHash, SourceSHA256: *expected, TenantID: boundTenant, OrganizationID: boundOrg, FactoryOrganizationID: doc.FactoryOrganizationID, Overrides: overrides, ProductionAuthorized: false}
	activationBytes, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return empty, err
	}
	activationBytes = append(activationBytes, '\n')
	dest, err := filepath.Abs(*output)
	if err != nil {
		return empty, err
	}
	// Mkdir is exclusive; existing directories/files are never reused. A partial
	// failure is preserved for inspection and cannot be mistaken for a new run.
	if err = os.Mkdir(dest, 0700); err != nil {
		return empty, err
	}
	for _, artifact := range []struct {
		name string
		data []byte
	}{{"profile.json", profileBytes}, {"activation.json", activationBytes}} {
		file, err := os.OpenFile(filepath.Join(dest, artifact.name), os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
		if err != nil {
			return empty, err
		}
		_, writeErr := file.Write(artifact.data)
		syncErr := file.Sync()
		closeErr := file.Close()
		if writeErr != nil {
			return empty, writeErr
		}
		if syncErr != nil {
			return empty, syncErr
		}
		if closeErr != nil {
			return empty, closeErr
		}
	}
	return value, nil
}
func main() {
	result, err := materialize(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, "warranty profile materialization failed:", err)
		os.Exit(1)
	}
	if err = json.NewEncoder(os.Stdout).Encode(result); err != nil {
		os.Exit(1)
	}
}
````

### FILE: `cmd/warranty-profile/main_test.go`

```yaml
block_id: "GO-CONNECTED-WARRANTY-CLAIM:file4:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "002fae89a9c2ab590deb4c30142f923d32ba2518887b3bb7762c13cd4cd434c3"
variables: []
secrets_allowed: false
```

````go
package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	wc "elite.local/enterprise/internal/warrantyclaim"
)

func TestWarrantyProfileMaterializesIdenticallyAndRefusesOverwrite(t *testing.T) {
	source := filepath.Join("..", "..", "deploy", "warranty", "profile.reference.json")
	raw, err := os.ReadFile(source)
	if err != nil {
		t.Fatal(err)
	}
	sha := wc.SHA(raw)
	root := t.TempDir()
	args := func(dest, hash string) []string {
		return []string{"--input", source, "--sha256", hash, "--output", dest, "--tenant-id", "bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbbbb", "--organization-id", "branch", "--factory-organization-id", "maker"}
	}
	first := filepath.Join(root, "first")
	second := filepath.Join(root, "second")
	a, err := materialize(args(first, sha))
	if err != nil {
		t.Fatal(err)
	}
	b, err := materialize(args(second, sha))
	if err != nil {
		t.Fatal(err)
	}
	if a.ProfileSHA256 != b.ProfileSHA256 || a.SourceSHA256 != sha || a.TenantID != "bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbbbb" || a.OrganizationID != "branch" || a.FactoryOrganizationID != "maker" || a.ProductionAuthorized {
		t.Fatal("activation binding", a, b)
	}
	for _, name := range []string{"profile.json", "activation.json"} {
		one, err := os.ReadFile(filepath.Join(first, name))
		if err != nil {
			t.Fatal(err)
		}
		two, err := os.ReadFile(filepath.Join(second, name))
		if err != nil || !bytes.Equal(one, two) {
			t.Fatal("nonportable rebuild", name, err)
		}
	}
	if _, err = materialize(args(first, sha)); err == nil {
		t.Fatal("existing destination overwritten")
	}
	bad := filepath.Join(root, "bad-hash")
	if _, err = materialize(args(bad, strings.Repeat("f", 64))); err == nil {
		t.Fatal("unlocked source admitted")
	}
	if _, err = os.Lstat(bad); !os.IsNotExist(err) {
		t.Fatal("invalid source created output")
	}
	invalid := filepath.Join(root, "invalid-override")
	wrong := args(invalid, sha)
	wrong[7] = "not-a-tenant"
	if _, err = materialize(wrong); err == nil {
		t.Fatal("invalid override admitted")
	}
	if _, err = os.Lstat(invalid); !os.IsNotExist(err) {
		t.Fatal("invalid scope created output")
	}
}
````

### FILE: `db/migrations/0069_warranty_terms_activation.down.sql`

```yaml
block_id: "GO-CONNECTED-WARRANTY-CLAIM:file5:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "2ad82e8dea6acfefb30f82572020b8dd406623a61403f739d4353f8b5562b5f3"
variables: []
secrets_allowed: false
```

````sql
begin;
do $$ begin
 if exists(select 1 from service_ops.warranty_offer) or exists(select 1 from service_ops.warranty_activation)
 then raise exception 'populated warranty terms cannot be downgraded without an admitted archival plan'; end if;
end $$;
drop trigger warranty_quote_acceptance_guard on sales.quotation;
drop function service_ops.warranty_quote_acceptance_guard();
drop table service_ops.warranty_activation;
drop function service_ops.warranty_activation_binding();
drop table service_ops.warranty_acknowledgement;
drop function service_ops.warranty_acknowledgement_binding();
drop table service_ops.warranty_offer;
drop function service_ops.warranty_offer_binding();
drop function service_ops.warranty_fact_immutable();
commit;
````

### FILE: `db/migrations/0069_warranty_terms_activation.up.sql`

```yaml
block_id: "GO-CONNECTED-WARRANTY-CLAIM:file6:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "d4ab0dd7dfa5ffe1e5b77e6f8a8b14079086809f00f6de21cf2b94e004065628"
variables: []
secrets_allowed: false
```

````sql
begin;
-- AUTHORED immutable consent and source-date persistence glue. No statutory
-- duration, fault exclusion, reimbursement or coverage percentage is inferred.
create table service_ops.warranty_offer (
 tenant_id uuid not null,
 quotation_id text not null,
 organization_id text not null,
 customer_principal_id text not null,
 quotation_version bigint not null check(quotation_version>1),
 profile_sha256 text not null check(profile_sha256 ~ '^[0-9a-f]{64}$'),
 profile_bytes bytea not null check(octet_length(profile_bytes) between 1 and 32768),
 offered_by text not null check(length(offered_by) between 1 and 128),
 offered_at timestamptz not null default clock_timestamp(),
 primary key(tenant_id,quotation_id),
 foreign key(tenant_id,quotation_id) references sales.quotation(tenant_id,quotation_id),
 foreign key(tenant_id,organization_id) references org.organization(tenant_id,organization_id),
 foreign key(tenant_id,customer_principal_id) references crm.customer_profile(tenant_id,customer_principal_id)
);
create table service_ops.warranty_acknowledgement (
 tenant_id uuid not null,
 quotation_id text not null,
 quotation_version bigint not null,
 customer_principal_id text not null,
 profile_sha256 text not null check(profile_sha256 ~ '^[0-9a-f]{64}$'),
 evidence_sha256 text not null check(evidence_sha256 ~ '^[0-9a-f]{64}$'),
 acknowledged_at timestamptz not null default clock_timestamp(),
 primary key(tenant_id,quotation_id),
 foreign key(tenant_id,quotation_id) references service_ops.warranty_offer(tenant_id,quotation_id)
);
create table service_ops.warranty_activation (
 tenant_id uuid not null,
 warranty_id text not null,
 handover_id text not null,
 quotation_id text not null,
 order_id text not null,
 organization_id text not null,
 customer_principal_id text not null,
 stock_unit_id text not null,
 profile_sha256 text not null check(profile_sha256 ~ '^[0-9a-f]{64}$'),
 terms_version text not null,
 parts_start date not null, parts_end date not null,
 labor_start date not null, labor_end date not null,
 handover_accepted_at timestamptz not null,
 activated_by text not null check(length(activated_by) between 1 and 128),
 activated_at timestamptz not null default clock_timestamp(),
 primary key(tenant_id,warranty_id),
 unique(tenant_id,handover_id),
 unique(tenant_id,stock_unit_id),
 foreign key(tenant_id,warranty_id) references service_ops.warranty(tenant_id,warranty_id),
 foreign key(tenant_id,quotation_id) references service_ops.warranty_acknowledgement(tenant_id,quotation_id),
 foreign key(tenant_id,handover_id) references sales.delivery_handover(tenant_id,handover_id),
 foreign key(tenant_id,order_id) references sales.customer_order(tenant_id,order_id),
 check(parts_start<=parts_end and labor_start<=labor_end)
);
create function service_ops.warranty_fact_immutable() returns trigger language plpgsql as $$
begin raise exception 'warranty fact is immutable'; end;
$$;
create trigger warranty_offer_immutable before update or delete on service_ops.warranty_offer
 for each row execute function service_ops.warranty_fact_immutable();
create trigger warranty_acknowledgement_immutable before update or delete on service_ops.warranty_acknowledgement
 for each row execute function service_ops.warranty_fact_immutable();
create trigger warranty_activation_immutable before update or delete on service_ops.warranty_activation
 for each row execute function service_ops.warranty_fact_immutable();

create function service_ops.warranty_offer_binding() returns trigger language plpgsql as $$
begin
 perform 1 from sales.quotation q where q.tenant_id=new.tenant_id and q.quotation_id=new.quotation_id
 and q.organization_id=new.organization_id and q.customer_principal_id=new.customer_principal_id
 and q.version=new.quotation_version and q.state='issued' and q.valid_until>clock_timestamp() for update;
 if not found then raise exception 'warranty offer requires current issued quotation'; end if;
 return new;
end;
$$;
create trigger warranty_offer_binding before insert on service_ops.warranty_offer
 for each row execute function service_ops.warranty_offer_binding();
create function service_ops.warranty_acknowledgement_binding() returns trigger language plpgsql as $$
begin
 perform 1 from sales.quotation q join service_ops.warranty_offer o using(tenant_id,quotation_id)
 where q.tenant_id=new.tenant_id and q.quotation_id=new.quotation_id and q.state='issued'
 and q.valid_until>clock_timestamp() and q.version=new.quotation_version and o.quotation_version=q.version
 and q.customer_principal_id=new.customer_principal_id and o.customer_principal_id=new.customer_principal_id
 and o.profile_sha256=new.profile_sha256 for update of q;
 if not found then raise exception 'warranty acknowledgement does not bind offered quotation'; end if;
 return new;
end;
$$;
create trigger warranty_acknowledgement_binding before insert on service_ops.warranty_acknowledgement
 for each row execute function service_ops.warranty_acknowledgement_binding();
create function service_ops.warranty_quote_acceptance_guard() returns trigger language plpgsql as $$
declare o service_ops.warranty_offer%rowtype;
begin
 select * into o from service_ops.warranty_offer where tenant_id=old.tenant_id and quotation_id=old.quotation_id;
 if not found then return new; end if;
 if row(new.organization_id,new.customer_principal_id,new.variant_id,new.price_book_id,new.currency,new.total_minor_units,new.valid_until)
 is distinct from row(old.organization_id,old.customer_principal_id,old.variant_id,old.price_book_id,old.currency,old.total_minor_units,old.valid_until)
 then raise exception 'offered warranty quote terms are immutable; issue a new quotation'; end if;
 if old.state='issued' and new.state='accepted' then
  perform 1 from service_ops.warranty_acknowledgement a where a.tenant_id=old.tenant_id and a.quotation_id=old.quotation_id
  and a.quotation_version=old.version and o.quotation_version=old.version and a.profile_sha256=o.profile_sha256
  and a.customer_principal_id=old.customer_principal_id;
  if not found then raise exception 'warranty terms acknowledgement required before quote acceptance'; end if;
 end if;
 return new;
end;
$$;
create trigger warranty_quote_acceptance_guard before update on sales.quotation
 for each row execute function service_ops.warranty_quote_acceptance_guard();

create function service_ops.warranty_activation_binding() returns trigger language plpgsql as $$
begin
 perform 1 from sales.delivery_handover h
 join sales.quotation_acceptance qa on qa.tenant_id=h.tenant_id and qa.order_id=h.order_id
 join service_ops.warranty_offer o on o.tenant_id=qa.tenant_id and o.quotation_id=qa.quotation_id
 join service_ops.warranty_acknowledgement a on a.tenant_id=o.tenant_id and a.quotation_id=o.quotation_id
 join service_ops.warranty w on w.tenant_id=h.tenant_id and w.warranty_id=new.warranty_id
 where h.tenant_id=new.tenant_id and h.handover_id=new.handover_id and h.state='accepted'
 and h.customer_accepted_at=new.handover_accepted_at and h.organization_id=new.organization_id
 and h.order_id=new.order_id and h.stock_unit_id=new.stock_unit_id and h.customer_principal_id=new.customer_principal_id
 and qa.quotation_id=new.quotation_id and qa.customer_principal_id=h.customer_principal_id
 and o.organization_id=h.organization_id and o.customer_principal_id=h.customer_principal_id
 and o.profile_sha256=new.profile_sha256 and a.profile_sha256=o.profile_sha256
 and a.quotation_version=o.quotation_version and a.customer_principal_id=h.customer_principal_id
 and w.stock_unit_id=h.stock_unit_id and w.customer_principal_id=h.customer_principal_id and w.terms_version=new.terms_version;
 if not found then raise exception 'warranty activation requires sold terms and accepted handover'; end if;
 return new;
end;
$$;
create trigger warranty_activation_binding before insert on service_ops.warranty_activation
 for each row execute function service_ops.warranty_activation_binding();
commit;
````

### FILE: `db/migrations/0070_warranty_claim.down.sql`

```yaml
block_id: "GO-CONNECTED-WARRANTY-CLAIM:file7:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "f2853fc70779b57f8caf5737d7ace3d72bd16b6c71ca371de78014a4788341d5"
variables: []
secrets_allowed: false
```

````sql
begin;
do $$ begin
 if exists(select 1 from service_ops.warranty_claim) or exists(select 1 from approval.request where kind='warranty_repair')
 then raise exception 'populated warranty claims cannot be downgraded';end if;
end $$;
drop trigger warranty_claim_case_guard on service_ops.service_case;
drop function service_ops.warranty_claim_case_guard();
drop table service_ops.warranty_part_reservation;
drop table service_ops.warranty_work_plan;
drop table service_ops.warranty_claim_step;
drop function service_ops.warranty_pending_commit_guard();
drop table service_ops.warranty_claim;
alter table approval.request drop constraint request_kind_check;
alter table approval.request add constraint request_kind_check check(kind in
 ('reservation','sale','refund','payment','whatsapp_reply','social_publish','social_revoke','stored_value_operation'));
alter table approval.request drop constraint request_payload_binding_check;
alter table approval.request add constraint request_payload_binding_check check
 (kind not in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation') or
  (organization_id is not null and length(organization_id) between 1 and 128 and
   payload is not null and jsonb_typeof(payload)='object' and pg_column_size(payload)<=65536 and amount_minor_units=0));
create or replace function approval.guard_bound_request() returns trigger language plpgsql as $$
begin
 if tg_op='DELETE' then
  if old.kind in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation') then raise exception 'bound approval request is immutable'; end if;
  return old;
 end if;
 if old.kind in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation') or new.kind in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation') then
  if row(new.tenant_id,new.request_id,new.kind,new.subject_id,new.amount_minor_units,new.requester,new.evidence_sha,new.organization_id,new.payload,new.created_at)
   is distinct from row(old.tenant_id,old.request_id,old.kind,old.subject_id,old.amount_minor_units,old.requester,old.evidence_sha,old.organization_id,old.payload,old.created_at)
   or (old.state<>'pending' and row(new.state,new.decided_at) is distinct from row(old.state,old.decided_at))
  then raise exception 'bound approval request is immutable'; end if;
 end if;
 return new;
end $$;
create or replace function approval.guard_bound_decision() returns trigger language plpgsql as $$
begin
 if exists(select 1 from approval.request where tenant_id=old.tenant_id and request_id=old.request_id and kind in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation')) then
  raise exception 'bound approval decision is immutable';
 end if;
 if tg_op='DELETE' then return old;end if;return new;
end $$;

commit;
````

### FILE: `db/migrations/0070_warranty_claim.up.sql`

```yaml
block_id: "GO-CONNECTED-WARRANTY-CLAIM:file8:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "374bd24812e0c3bb7c6578e4cc1f1486865855fdd5249a8c6619f0bf0681b3e5"
variables: []
secrets_allowed: false
```

````sql
begin;
-- AUTHORED service/approval/inventory binding; upstream owners remain authoritative.
alter table approval.request drop constraint request_kind_check;
alter table approval.request add constraint request_kind_check check(kind in
 ('reservation','sale','refund','payment','whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair'));
alter table approval.request drop constraint request_payload_binding_check;
alter table approval.request add constraint request_payload_binding_check check
 (kind not in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair') or
  (organization_id is not null and length(organization_id) between 1 and 128 and
   payload is not null and jsonb_typeof(payload)='object' and pg_column_size(payload)<=65536 and amount_minor_units=0));
create or replace function approval.guard_bound_request() returns trigger language plpgsql as $$
begin
 if tg_op='DELETE' then
  if old.kind in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair') then raise exception 'bound approval request is immutable'; end if;
  return old;
 end if;
 if old.kind in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair') or new.kind in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair') then
  if row(new.tenant_id,new.request_id,new.kind,new.subject_id,new.amount_minor_units,new.requester,new.evidence_sha,new.organization_id,new.payload,new.created_at)
   is distinct from row(old.tenant_id,old.request_id,old.kind,old.subject_id,old.amount_minor_units,old.requester,old.evidence_sha,old.organization_id,old.payload,old.created_at)
   or (old.state<>'pending' and row(new.state,new.decided_at) is distinct from row(old.state,old.decided_at))
  then raise exception 'bound approval request is immutable'; end if;
 end if;
 return new;
end $$;
create or replace function approval.guard_bound_decision() returns trigger language plpgsql as $$
begin
 if exists(select 1 from approval.request where tenant_id=old.tenant_id and request_id=old.request_id and kind in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair')) then
  raise exception 'bound approval decision is immutable';
 end if;
 if tg_op='DELETE' then return old;end if;return new;
end $$;

create table service_ops.warranty_claim (
 tenant_id uuid not null, service_case_id text not null, warranty_id text not null,
 handover_id text not null, appointment_id text not null,
 organization_id text not null, factory_organization_id text not null,
 customer_principal_id text not null, stock_unit_id text not null,
 service_date date not null, profile_sha256 text not null check(profile_sha256 ~ '^[0-9a-f]{64}$'),
 parts_covered boolean not null, labor_covered boolean not null,
 opened_by text not null, request_sha256 text not null check(request_sha256 ~ '^[0-9a-f]{64}$'),
 opened_at timestamptz not null default clock_timestamp(),
 primary key(tenant_id,service_case_id),
 unique(tenant_id,appointment_id,stock_unit_id),
 foreign key(tenant_id,service_case_id) references service_ops.service_case(tenant_id,service_case_id),
 foreign key(tenant_id,warranty_id) references service_ops.warranty_activation(tenant_id,warranty_id),
 foreign key(tenant_id,appointment_id) references crm.appointment(tenant_id,appointment_id),
 foreign key(tenant_id,factory_organization_id) references org.organization(tenant_id,organization_id)
);
create trigger warranty_claim_immutable before update or delete on service_ops.warranty_claim
 for each row execute function service_ops.warranty_fact_immutable();
create table service_ops.warranty_claim_step (
 tenant_id uuid not null, service_case_id text not null, command_id text not null,
 version bigint not null check(version>0), kind text not null check(kind in ('opened','diagnosed','planned','decision','work','quality','accepted','reconciled','cancelled')),
 from_state text not null, to_state text not null, actor text not null check(length(actor) between 1 and 128),
 request_sha256 text not null check(request_sha256 ~ '^[0-9a-f]{64}$'),
 payload_sha256 text not null check(payload_sha256 ~ '^[0-9a-f]{64}$'),
 payload jsonb not null check(jsonb_typeof(payload)='object' and pg_column_size(payload)<=65536),
 recorded_at timestamptz not null default clock_timestamp(),
 primary key(tenant_id,service_case_id,version),
 unique(tenant_id,service_case_id,command_id),
 foreign key(tenant_id,service_case_id) references service_ops.warranty_claim(tenant_id,service_case_id)
);
create unique index warranty_single_step_kind on service_ops.warranty_claim_step(tenant_id,service_case_id,kind) where kind<>'quality';
create trigger warranty_claim_step_immutable before update or delete on service_ops.warranty_claim_step
 for each row execute function service_ops.warranty_fact_immutable();
create table service_ops.warranty_work_plan (
 tenant_id uuid not null, service_case_id text not null, approval_id text not null,
 payload_sha256 text not null check(payload_sha256 ~ '^[0-9a-f]{64}$'),
 payload jsonb not null check(jsonb_typeof(payload)='object' and pg_column_size(payload)<=65536),
 expires_at timestamptz not null, created_at timestamptz not null default clock_timestamp(),
 primary key(tenant_id,service_case_id),unique(tenant_id,approval_id),
 foreign key(tenant_id,service_case_id) references service_ops.warranty_claim(tenant_id,service_case_id),
 foreign key(tenant_id,approval_id) references approval.request(tenant_id,request_id),
 check(expires_at>created_at)
);
create trigger warranty_work_plan_immutable before update or delete on service_ops.warranty_work_plan
 for each row execute function service_ops.warranty_fact_immutable();
create table service_ops.warranty_part_reservation (
 tenant_id uuid not null, service_case_id text not null, line_id text not null, reservation_id text not null,
 primary key(tenant_id,service_case_id,line_id),unique(tenant_id,reservation_id),
 foreign key(tenant_id,service_case_id) references service_ops.warranty_work_plan(tenant_id,service_case_id),
 foreign key(tenant_id,reservation_id) references inventory.bulk_reservation(tenant_id,reservation_id)
);
create trigger warranty_part_reservation_immutable before update or delete on service_ops.warranty_part_reservation
 for each row execute function service_ops.warranty_fact_immutable();

create function service_ops.warranty_claim_case_guard() returns trigger language plpgsql as $$
begin
 if not exists(select 1 from service_ops.warranty_claim c where c.tenant_id=old.tenant_id and c.service_case_id=old.service_case_id) then
  if tg_op='DELETE' then return old;end if;return new;
 end if;
 if tg_op='DELETE' then raise exception 'warranty claim case cannot be deleted';end if;
 if row(new.tenant_id,new.service_case_id,new.organization_id,new.stock_unit_id,new.severity,new.description,new.opened_at)
 is distinct from row(old.tenant_id,old.service_case_id,old.organization_id,old.stock_unit_id,old.severity,old.description,old.opened_at)
 then raise exception 'warranty case identity is immutable';end if;
 if new.version<>old.version+1 or not exists(select 1 from service_ops.warranty_claim_step s
 where s.tenant_id=old.tenant_id and s.service_case_id=old.service_case_id and s.version=new.version
 and s.from_state=old.state and s.to_state=new.state)
 then raise exception 'warranty case requires its connected command receipt';end if;
 return new;
end $$;
create trigger warranty_claim_case_guard before update or delete on service_ops.service_case
 for each row execute function service_ops.warranty_claim_case_guard();

-- A pending reservation may expire while a transaction is waiting. Approval
-- binds it durably to the repair; it is no longer a pending reservation lease.
create function service_ops.warranty_pending_commit_guard() returns trigger language plpgsql as $$
declare expiry timestamptz;
begin
 if tg_table_name='warranty_work_plan' then expiry:=new.expires_at;
 elsif new.kind='decision' and new.to_state='repair' then
  select expires_at into expiry from service_ops.warranty_work_plan where tenant_id=new.tenant_id and service_case_id=new.service_case_id;
 else return null;end if;
 if expiry is null or expiry<=clock_timestamp() then raise exception 'warranty pending reservation expired before commit';end if;
 return null;
end $$;
create constraint trigger warranty_plan_commit_expiry after insert on service_ops.warranty_work_plan
 deferrable initially deferred for each row execute function service_ops.warranty_pending_commit_guard();
create constraint trigger warranty_approval_commit_expiry after insert on service_ops.warranty_claim_step
 deferrable initially deferred for each row execute function service_ops.warranty_pending_commit_guard();
commit;
````

### FILE: `deploy/warranty/profile.reference.json`

```yaml
block_id: "GO-CONNECTED-WARRANTY-CLAIM:file9:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "ed494deafc6fc950d7edb4ac1b8809897762a794ab9baed539427c64de61d8d0"
variables: []
secrets_allowed: false
```

````json
{
  "schema": "elite-warranty-profile/v1",
  "scope": "MATERIALIZED_PROFILE",
  "algorithm": "bc-inclusive-fixed-terms",
  "algorithm_revision": 1,
  "tenant_id": "aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa",
  "organization_id": "store",
  "factory_organization_id": "factory",
  "policy_id": "reference-service-warranty",
  "terms_version": "reference-fixture-v1",
  "terms_text": "Fixture local de servicio: intervalos configurados de piezas y mano de obra, revisión humana de causa, calidad, aceptación del cliente y comprobante interno. Estos datos de referencia no certifican términos legales de un país.",
  "business_time_zone": "America/Argentina/Buenos_Aires",
  "parts_duration_days": 365,
  "labor_duration_days": 365,
  "work_reservation_seconds": 3600,
  "fault_exclusions": ["fixture-exclusion"],
  "settlement": "INTERNAL_WARRANTY_SERVICE_ACKNOWLEDGEMENT",
  "authority_reference": "reference-fixture-terms",
  "decision_reference": "reference-fixture-selection"
}
````

### FILE: `docs/WARRANTY_REFERENCE.md`

```yaml
block_id: "GO-CONNECTED-WARRANTY-CLAIM:file10:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "340a32a0317216cfe66586b16726e99913e7024ed2387196d033ee201f8b1158"
variables: []
secrets_allowed: false
```

````markdown
# Garantía de servicio en la composición de referencia

La garantía conserva los términos de la cotización aceptada. El pedido, pago,
entrega, caso de servicio, aprobación y stock mantienen sus escritores existentes.
Las comparaciones inclusivas de fechas de piezas/mano de obra son la adaptación
de `ServiceItemLine.CheckWarranty` identificada en
`docs/provenance/BC_WARRANTY_DERIVATION.json`. Configuración, consentimiento,
calendario, transporte y coordinación transaccional son AUTHORED; este módulo no
reproduce todo Business Central ni se atribuye ese trabajo a Microsoft.

El perfil configura explícitamente términos, zona horaria, días, exclusiones,
organizaciones y duración de la reserva pendiente. El último día incluido se calcula como fecha inicial + `*_duration_days`; el predicado incluye ambos extremos. El ejemplo es un fixture de
infraestructura. No contiene credenciales ni afirma obligaciones legales de un
país. La conciliación seleccionada es un comprobante interno de servicio, con
referencias al costo de inventario. No crea un reembolso, nómina o documento fiscal.

Desde la raíz materializada, un perfil se genera en un destino ausente:

```powershell
$source = 'deploy/warranty/profile.reference.json'
$sourceHash = (Get-FileHash -LiteralPath $source -Algorithm SHA256).Hash.ToLowerInvariant()
go run ./cmd/warranty-profile --input $source --sha256 $sourceHash --output deploy/warranty/active
$activation = Get-Content -LiteralPath deploy/warranty/active/activation.json -Raw | ConvertFrom-Json
$env:WARRANTY_ENABLED = 'true'
$env:WARRANTY_TENANT_ID = $activation.tenant_id
$env:WARRANTY_ORGANIZATION_ID = $activation.organization_id
$env:WARRANTY_PROFILE_SHA256 = $activation.profile_sha256
$env:WARRANTY_PROFILE_FILE = (Resolve-Path -LiteralPath deploy/warranty/active/profile.json).Path
```

El CLI admite overrides explícitos `--tenant-id`, `--organization-id` y
`--factory-organization-id`; vuelve a validar mediante el mismo loader del host.
Conserva el hash de origen y los cambios en `activation.json`. Ambos archivos se
reconstruyen con bytes idénticos en directorios distintos; la activación contiene
una ruta relativa. No reutiliza ni sobrescribe destinos existentes. Un fallo de
escritura conserva el destino parcial para inspección.

Con `WARRANTY_ENABLED=false` no se leen perfiles. Activado exige ruta absoluta,
archivo regular de hasta32768bytes, hash y scope exactos, migraciones0069/0070 y
los guards habilitados. La política de términos ofrecidos se vuelve inmutable;
un cambio de precio/cliente/variante requiere otra cotización. El cliente reconoce
el hash y versión antes de aceptar. Una configuración posterior no modifica el
perfil vendido ni inventa una anulación por reembolso.

Todos los endpoints requieren `organization_id` y bearer verificado. La API
reduce la autoridad al permiso y organización de esa ruta, incluso con roles
que poseen varias organizaciones. Las versiones viajan como strings decimales;
se rechazan claves duplicadas, campos desconocidos y objetos mayores de32768bytes.

| Operación | Ruta POST | Permiso |
| --- | --- | --- |
| Ofrecer términos | `/v1/franchise/warranty/quotes/{id}/terms` | `warranty:offer` |
| Consentimiento | `/v1/customer/warranty/quotes/{id}/acknowledgements` | `warranty:self` |
| Activación | `/v1/franchise/warranty/handovers/{id}/activation` | `warranty:activate` |
| Abrir caso | `/v1/customer/warranty/claims` o `/v1/franchise/warranty/claims` | `warranty:self` o `warranty:request` |
| Diagnóstico | `/v1/franchise/warranty/claims/{id}/diagnosis` | `warranty:diagnose` |
| Plan/reservas | `/v1/franchise/warranty/claims/{id}/plan` | `warranty:plan` |
| Decisión | `/v1/franchise/warranty/claims/{id}/decision` | `warranty:approve` |
| Trabajo/repuestos | `/v1/franchise/warranty/claims/{id}/work` | `warranty:work` |
| Calidad | `/v1/franchise/warranty/claims/{id}/quality` | `warranty:quality` |
| Aceptación | `/v1/customer/warranty/claims/{id}/acceptance` | `warranty:self` |
| Conciliación | `/v1/factory/warranty/claims/{id}/reconciliation` | `warranty:reconcile` |
| Cancelación previa al consumo | `/v1/franchise/warranty/claims/{id}/cancellation` | `warranty:cancel` |

Las lecturas de cotización/activación tienen rutas GET de cliente y franquicia.
El perfil actual se consulta en `/v1/franchise/warranty/profile`. El caso y su
último recibo se consultan en `/v1/{customer|franchise|factory}/warranty/claims/{id}`;
`command_id` recupera el recibo exacto de un comando. Permisos de lectura:
`warranty:self`, `warranty:read` o `warranty:factory-read`. La fábrica se toma del
perfil vendido; el cliente debe coincidir con el propietario del caso.

La apertura liga unidad, cliente y cita de servicio confirmada/completada. El
diagnóstico conserva causa y evidencia. Un plan por caso admite hasta16líneas
de repuestos enteros y trabajo de servicio; usa reservas `demand_kind=service`.
La decisión exige un sujeto diferente y payload exacto. No se autoaprueba. Un
plan fuera de cobertura puede rechazarse; no reserva piezas ni permite aprobación.
La aprobación convierte la reserva pendiente en hold durable de reparación.

El trabajo consume el stock con el writer FIFO/specific existente en la misma
transacción Serializable del recibo. La calidad la registra un sujeto diferente
del técnico. Un resultado fallido permite otra revisión con evidencia de
corrección; no permite aceptación. El cliente acepta el resultado aprobado y la
fábrica confirma el comprobante enlazado antes de cerrar. Nuevas necesidades de
repuestos después de un plan o consumos ya terminados requieren un caso de
seguimiento; este perfil no agrega líneas silenciosamente a un plan aprobado.

Ante respuesta perdida, consultar el comando guardado antes de decidir un retry.
El replay exige actor y payload originales. El rechazo o la cancelación anterior
al consumo liberan las reservas y conservan la decisión histórica. No revierten
stock ya consumido ni alteran dinero. El caso genérico no puede omitir los recibos
conectados. Los downgrades de garantía rechazan datos existentes; una base vacía
puede bajar/subir ambas migraciones. El frontend por rol se compone por su owner.
````

### FILE: `docs/provenance/BC_WARRANTY_NOTICES.md`

```yaml
block_id: "GO-CONNECTED-WARRANTY-CLAIM:file11:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "f59682d140362666dfdc39aea014d3a6389f9ef1183199261e993de2696ddcb8"
variables: []
secrets_allowed: false
```

````markdown
# Warranty source and integration notice

Microsoft BCApps snapshot 2eae56d704a1fd035d104f333602aea7091b7749, MIT.
The complete selected ServiceItemLine.Table.al is retained VERBATIM under
source/ServiceItemLine.Table.al (SHA256
32b4fcf60b25d363d31ef6a1927b1f051eb4ef4ef48d9ca6b2ecc04844a854ff).
The exact upstream MIT notice is licenses/Microsoft-BCApps-MIT.txt (SHA256
c2cfccb812fe482101a8f04597dfc5a9991a6b2748266c47ac91b6a5aae15383).
Source: https://github.com/microsoft/BCApps/blob/2eae56d704a1fd035d104f333602aea7091b7749/src/Layers/W1/BaseApp/Service/Document/ServiceItemLine.Table.al

Only the inclusive parts/labor date predicates and invalid-parts-period result
from CheckWarranty (lines1884–1953, segment SHA256
e46ec6bb958d2d3b946ef59a3833664dc2c8d13351dd6c876493a3980c5b95e0)
are ADAPTED in internal/warrantycoverage/coverage.go. Ordinal dates and Go errors
are explicit representation adaptations. Source-projected vectors and generator
ship with the source. No Microsoft AL runtime, event subscribers or full service
suite is executed or claimed.

BC_WARRANTY_DERIVATION.json preserves the earlier source investigation state;
its CANDIDATE/SOURCE_IDENTIFIED fields are historical, not release status.
Current narrow admission and exact connected reconstruction are recorded in
WARRANTY_SOURCE_SCOPE_V402.md and WARRANTY_CONNECTED_RELEASE_V402.md in the library.

All sold-term binding, configuration, authorization, PostgreSQL transactions,
HTTP/host/materializer and tests around that predicate are AUTHORED local glue.
Existing appointment/service/stock/FIFO/approval owners remain authoritative;
there is no second stock or money ledger. No corporate authorship is attributed
to that glue. Workspace-owner licensing applies as declared per pack block.

The named reference policy is a fixture. Days and exclusions are explicit;
start plus configured days is inclusive at both endpoints. Factory reconciliation
is an internal service acknowledgement and creates no provider payment or fiscal
invoice. A target supplies its approved policy and identity grants explicitly;
this is not a statement of legal warranty rights or production acceptance.
````

### FILE: `internal/approval/warranty_test.go`

```yaml
block_id: "GO-CONNECTED-WARRANTY-CLAIM:file12:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "dcfcfe58ecebd6e8f6780986ea4cd0b25c3a84f86988c52f8b4c207f66d29077"
variables: []
secrets_allowed: false
```

````go
package approval

import (
	"strings"
	"testing"
)

func TestWarrantyRepairAlwaysNeedsDistinctHuman(t *testing.T) {
	r := NewRegistry(Policy{AutoApproveMinorUnits: 1000000})
	request := Request{TenantID: "fixture", ID: "warranty", Kind: KindWarrantyRepair, SubjectID: "case", Requester: "technician", EvidenceSHA: strings.Repeat("a", 64)}
	state, err := r.Submit(request)
	if err != nil || state != StatePending {
		t.Fatal("warranty auto-approved", state, err)
	}
	if _, err = r.Approve("fixture", "warranty", "technician", "self"); err != ErrSeparation {
		t.Fatal("self approval", err)
	}
	state, err = r.Approve("fixture", "warranty", "reviewer", "observed evidence")
	if err != nil || state != StateApproved {
		t.Fatal(state, err)
	}
}
````

### FILE: `internal/platform/httpapi/warranty.go`

```yaml
block_id: "GO-CONNECTED-WARRANTY-CLAIM:file13:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "359321f457fb422cd1101f7044c6cf135b2f2ef9fc45d664e7ec62075adf33ff"
variables: []
secrets_allowed: false
```

````go
// AUTHORED bounded transport. Verified principal/scope and existing domain
// owners determine effects; request JSON cannot choose tenant or actor.
package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"elite.local/enterprise/internal/approval"
	"elite.local/enterprise/internal/platform/identity"
	wc "elite.local/enterprise/internal/warrantyclaim"
	"github.com/jackc/pgx/v5/pgconn"
)

type WarrantyService interface {
	Scope() (string, string)
	ProfileSHA256() string
	ProfileDocument(context.Context, identity.Principal) (json.RawMessage, error)
	Offer(context.Context, identity.Principal, string) (wc.Offer, error)
	BindOffer(context.Context, identity.Principal, wc.OfferRequest) (wc.Offer, error)
	Acknowledge(context.Context, identity.Principal, wc.AcknowledgeRequest) (wc.Offer, error)
	Activation(context.Context, identity.Principal, string) (wc.Activation, error)
	Activate(context.Context, identity.Principal, string) (wc.Activation, error)
	Claim(context.Context, identity.Principal, string) (wc.Claim, error)
	ClaimCommand(context.Context, identity.Principal, string, string) (wc.Step, error)
	OpenClaim(context.Context, identity.Principal, wc.OpenClaim) (wc.Step, error)
	Diagnose(context.Context, identity.Principal, wc.Diagnose) (wc.Step, error)
	PlanRepair(context.Context, identity.Principal, wc.Plan) (wc.Step, error)
	DecideRepair(context.Context, identity.Principal, wc.DecideRepair) (wc.Step, error)
	CompleteRepairWork(context.Context, identity.Principal, wc.CompleteWork) (wc.Step, error)
	RecordRepairQuality(context.Context, identity.Principal, wc.Quality) (wc.Step, error)
	AcceptRepair(context.Context, identity.Principal, wc.AcceptRepair) (wc.Step, error)
	ReconcileRepair(context.Context, identity.Principal, wc.ReconcileRepair) (wc.Step, error)
	CancelRepair(context.Context, identity.Principal, wc.CancelRepair) (wc.Step, error)
}
type WarrantyQuoteReader interface {
	QuoteForOffer(context.Context, identity.Principal, string) (wc.QuoteForOffer, error)
}
type WarrantyModule struct{ Service WarrantyService }
type warrantyProtection func(http.ResponseWriter, *http.Request, string, bool) (identity.Principal, bool)

func warrantyHTTPError(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}
	var pg *pgconn.PgError
	switch {
	case errors.Is(err, wc.ErrNotFound), errors.Is(err, approval.ErrNotFound):
		writeProblem(w, 404, "WARRANTY_NOT_FOUND", "no warranty resource exists in this authorized scope")
	case errors.Is(err, wc.ErrInvalid):
		writeProblem(w, 400, "WARRANTY_INVALID", "the warranty command is invalid")
	case errors.Is(err, wc.ErrConflict), errors.Is(err, approval.ErrDuplicate), errors.Is(err, approval.ErrNotPending), errors.Is(err, approval.ErrSeparation), errors.Is(err, approval.ErrInvalidRequest):
		writeProblem(w, 409, "WARRANTY_CONFLICT", "consult the saved command and current case version before retrying")
	case errors.As(err, &pg) && (pg.Code == "23505" || pg.Code == "23514" || pg.Code == "P0001" || pg.Code == "40001" || pg.Code == "40P01"):
		writeProblem(w, 409, "WARRANTY_CONFLICT", "consult the saved command and current case version before retrying")
	default:
		writeProblem(w, 503, "WARRANTY_UNCONFIRMED", "the result is unconfirmed; consult its saved command reference")
	}
	return true
}
func warrantyDecode(w http.ResponseWriter, r *http.Request, out any) bool {
	raw, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 32768))
	if err != nil {
		writeProblem(w, 413, "WARRANTY_BODY_TOO_LARGE", "the command exceeds its bounded size")
		return false
	}
	// CanonicalPayload rejects duplicate keys, including case-insensitive ones,
	// before typed decoding could discard an ambiguous field.
	if _, _, err = approval.CanonicalPayload(raw); err != nil {
		writeProblem(w, 400, "WARRANTY_INVALID_JSON", "a bounded unambiguous JSON object is required")
		return false
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	if decoder.Decode(out) != nil || decoder.Decode(new(any)) != io.EOF {
		writeProblem(w, 400, "WARRANTY_INVALID_JSON", "unknown fields or invalid typed values")
		return false
	}
	return true
}
func warrantyReply(w http.ResponseWriter, status int, value any) {
	replay := false
	switch v := value.(type) {
	case wc.Step:
		replay = v.Replay
	case wc.Offer:
		replay = v.Replay
	case wc.Activation:
		replay = v.Replay
	}
	if replay {
		status = 200
	}
	writeJSON(w, status, value)
}
func warrantyPost[T any](mux *http.ServeMux, route string, protect warrantyProtection, permission string, factory bool, id func(T) string, action func(context.Context, identity.Principal, T) (any, error)) {
	mux.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		p, ok := protect(w, r, permission, factory)
		if !ok {
			return
		}
		var c T
		if !warrantyDecode(w, r, &c) {
			return
		}
		if path := r.PathValue("id"); path != "" && id(c) != path {
			warrantyHTTPError(w, wc.ErrInvalid)
			return
		}
		value, err := action(r.Context(), p, c)
		if warrantyHTTPError(w, err) {
			return
		}
		warrantyReply(w, 201, value)
	})
}
func (m WarrantyModule) Register(mux *http.ServeMux, verifier identity.Verifier) {
	if m.Service == nil {
		return
	}
	api := franchiseJourneyAPI{verifier: verifier}
	protect := func(w http.ResponseWriter, r *http.Request, permission string, factory bool) (identity.Principal, bool) {
		query := r.URL.Query()
		if len(query["organization_id"]) != 1 || query.Get("organization_id") == "" {
			writeProblem(w, 400, "WARRANTY_SCOPE_REQUIRED", "one organization_id is required")
			return identity.Principal{}, false
		}
		org := query.Get("organization_id")
		p, ok := api.protected(w, r, permission, org)
		if !ok {
			return p, false
		}
		tenant, bound := m.Service.Scope()
		if p.TenantID != tenant || !factory && org != bound {
			writeProblem(w, 403, "FORBIDDEN", "the module is outside this request scope")
			return p, false
		}
		// Pass only the authority verified for this route/organization. A principal
		// with several organizations or '*' cannot silently switch request scope.
		p.Permissions = map[string]struct{}{permission: {}}
		p.Organizations = map[string]struct{}{org: {}}
		return p, true
	}
	mux.HandleFunc("GET /v1/franchise/warranty/profile", func(w http.ResponseWriter, r *http.Request) {
		p, ok := protect(w, r, "warranty:read", false)
		if !ok {
			return
		}
		raw, err := m.Service.ProfileDocument(r.Context(), p)
		if warrantyHTTPError(w, err) {
			return
		}
		writeJSON(w, 200, map[string]any{"profile_sha256": m.Service.ProfileSHA256(), "profile": raw})
	})
	if reader, ok := m.Service.(WarrantyQuoteReader); ok {
		mux.HandleFunc("GET /v1/franchise/warranty/quotes/{id}/source", func(w http.ResponseWriter, r *http.Request) {
			p, ok := protect(w, r, "warranty:offer", false)
			if !ok {
				return
			}
			v, e := reader.QuoteForOffer(r.Context(), p, r.PathValue("id"))
			if !warrantyHTTPError(w, e) {
				warrantyReply(w, 200, v)
			}
		})
	}
	for _, surface := range []struct{ path, permission string }{{"franchise", "warranty:read"}, {"customer", "warranty:self"}} {
		mux.HandleFunc("GET /v1/"+surface.path+"/warranty/quotes/{id}", func(w http.ResponseWriter, r *http.Request) {
			p, ok := protect(w, r, surface.permission, false)
			if !ok {
				return
			}
			v, err := m.Service.Offer(r.Context(), p, r.PathValue("id"))
			if !warrantyHTTPError(w, err) {
				warrantyReply(w, 200, v)
			}
		})
		mux.HandleFunc("GET /v1/"+surface.path+"/warranty/handovers/{id}/activation", func(w http.ResponseWriter, r *http.Request) {
			p, ok := protect(w, r, surface.permission, false)
			if !ok {
				return
			}
			v, err := m.Service.Activation(r.Context(), p, r.PathValue("id"))
			if !warrantyHTTPError(w, err) {
				warrantyReply(w, 200, v)
			}
		})
	}
	for _, surface := range []struct {
		path, permission string
		factory          bool
	}{{"franchise", "warranty:read", false}, {"customer", "warranty:self", false}, {"factory", "warranty:factory-read", true}} {
		mux.HandleFunc("GET /v1/"+surface.path+"/warranty/claims/{id}", func(w http.ResponseWriter, r *http.Request) {
			p, ok := protect(w, r, surface.permission, surface.factory)
			if !ok {
				return
			}
			if command := r.URL.Query().Get("command_id"); command != "" {
				v, err := m.Service.ClaimCommand(r.Context(), p, r.PathValue("id"), command)
				if !warrantyHTTPError(w, err) {
					warrantyReply(w, 200, v)
				}
				return
			}
			v, err := m.Service.Claim(r.Context(), p, r.PathValue("id"))
			if !warrantyHTTPError(w, err) {
				warrantyReply(w, 200, v)
			}
		})
	}
	warrantyPost(mux, "POST /v1/franchise/warranty/quotes/{id}/terms", protect, "warranty:offer", false, func(c wc.OfferRequest) string { return c.QuoteID }, func(ctx context.Context, p identity.Principal, c wc.OfferRequest) (any, error) {
		return m.Service.BindOffer(ctx, p, c)
	})
	warrantyPost(mux, "POST /v1/customer/warranty/quotes/{id}/acknowledgements", protect, "warranty:self", false, func(c wc.AcknowledgeRequest) string { return c.QuoteID }, func(ctx context.Context, p identity.Principal, c wc.AcknowledgeRequest) (any, error) {
		return m.Service.Acknowledge(ctx, p, c)
	})
	type activate struct {
		HandoverID string `json:"handover_id"`
	}
	warrantyPost(mux, "POST /v1/franchise/warranty/handovers/{id}/activation", protect, "warranty:activate", false, func(c activate) string { return c.HandoverID }, func(ctx context.Context, p identity.Principal, c activate) (any, error) {
		return m.Service.Activate(ctx, p, c.HandoverID)
	})
	for _, surface := range []struct{ path, permission string }{{"franchise", "warranty:request"}, {"customer", "warranty:self"}} {
		warrantyPost(mux, "POST /v1/"+surface.path+"/warranty/claims", protect, surface.permission, false, func(c wc.OpenClaim) string { return c.CaseID }, func(ctx context.Context, p identity.Principal, c wc.OpenClaim) (any, error) {
			return m.Service.OpenClaim(ctx, p, c)
		})
	}
	warrantyPost(mux, "POST /v1/franchise/warranty/claims/{id}/diagnosis", protect, "warranty:diagnose", false, func(c wc.Diagnose) string { return c.CaseID }, func(ctx context.Context, p identity.Principal, c wc.Diagnose) (any, error) {
		return m.Service.Diagnose(ctx, p, c)
	})
	warrantyPost(mux, "POST /v1/franchise/warranty/claims/{id}/plan", protect, "warranty:plan", false, func(c wc.Plan) string { return c.CaseID }, func(ctx context.Context, p identity.Principal, c wc.Plan) (any, error) {
		return m.Service.PlanRepair(ctx, p, c)
	})
	warrantyPost(mux, "POST /v1/franchise/warranty/claims/{id}/decision", protect, "warranty:approve", false, func(c wc.DecideRepair) string { return c.CaseID }, func(ctx context.Context, p identity.Principal, c wc.DecideRepair) (any, error) {
		return m.Service.DecideRepair(ctx, p, c)
	})
	warrantyPost(mux, "POST /v1/franchise/warranty/claims/{id}/work", protect, "warranty:work", false, func(c wc.CompleteWork) string { return c.CaseID }, func(ctx context.Context, p identity.Principal, c wc.CompleteWork) (any, error) {
		return m.Service.CompleteRepairWork(ctx, p, c)
	})
	warrantyPost(mux, "POST /v1/franchise/warranty/claims/{id}/quality", protect, "warranty:quality", false, func(c wc.Quality) string { return c.CaseID }, func(ctx context.Context, p identity.Principal, c wc.Quality) (any, error) {
		return m.Service.RecordRepairQuality(ctx, p, c)
	})
	warrantyPost(mux, "POST /v1/customer/warranty/claims/{id}/acceptance", protect, "warranty:self", false, func(c wc.AcceptRepair) string { return c.CaseID }, func(ctx context.Context, p identity.Principal, c wc.AcceptRepair) (any, error) {
		return m.Service.AcceptRepair(ctx, p, c)
	})
	warrantyPost(mux, "POST /v1/factory/warranty/claims/{id}/reconciliation", protect, "warranty:reconcile", true, func(c wc.ReconcileRepair) string { return c.CaseID }, func(ctx context.Context, p identity.Principal, c wc.ReconcileRepair) (any, error) {
		return m.Service.ReconcileRepair(ctx, p, c)
	})
	warrantyPost(mux, "POST /v1/franchise/warranty/claims/{id}/cancellation", protect, "warranty:cancel", false, func(c wc.CancelRepair) string { return c.CaseID }, func(ctx context.Context, p identity.Principal, c wc.CancelRepair) (any, error) {
		return m.Service.CancelRepair(ctx, p, c)
	})
}
````

### FILE: `internal/platform/httpapi/warranty_test.go`

```yaml
block_id: "GO-CONNECTED-WARRANTY-CLAIM:file14:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "7ad62ae0c3a0845ec70e1dc3cf869d005cb5dd602eccd1cb8f47f522a4b2230c"
variables: []
secrets_allowed: false
```

````go
package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"elite.local/enterprise/internal/platform/identity"
	wc "elite.local/enterprise/internal/warrantyclaim"
)

type warrantyTestVerifier struct{ p identity.Principal }

func (v warrantyTestVerifier) Verify(_ context.Context, token string) (identity.Principal, error) {
	if token != "fixture" {
		return identity.Principal{}, identity.ErrUnauthenticated
	}
	return v.p, nil
}

type warrantyTestService struct {
	WarrantyService
	calls     int
	principal identity.Principal
	request   wc.OfferRequest
}

func (*warrantyTestService) Scope() (string, string) {
	return "aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa", "store"
}
func (s *warrantyTestService) BindOffer(_ context.Context, p identity.Principal, r wc.OfferRequest) (wc.Offer, error) {
	s.calls++
	s.principal = p
	s.request = r
	return wc.Offer{QuoteID: r.QuoteID, QuoteVersion: r.QuoteVersion, ProfileSHA256: r.ProfileSHA256}, nil
}
func TestWarrantyTransportPreservesExactVersionsAndRejectsAmbiguity(t *testing.T) {
	profile := strings.Repeat("a", 64)
	valid := `{"quote_id":"quote","quote_version":"9007199254740993","profile_sha256":"` + profile + `"}`
	for _, test := range []struct {
		name, body, query, token string
		code                     int
	}{
		{"exact-int64", valid, "organization_id=store", "fixture", 201},
		{"number-not-string", strings.Replace(valid, `"9007199254740993"`, `9007199254740993`, 1), "organization_id=store", "fixture", 400},
		{"overflow", strings.Replace(valid, "9007199254740993", "9223372036854775808", 1), "organization_id=store", "fixture", 400},
		{"duplicate-case", strings.Replace(valid, `"quote_id":"quote"`, `"quote_id":"quote","Quote_ID":"other"`, 1), "organization_id=store", "fixture", 400},
		{"unknown-actor", strings.Replace(valid, `{`, `{"actor":"forged",`, 1), "organization_id=store", "fixture", 400},
		{"path-body-mismatch", strings.Replace(valid, `"quote"`, `"other"`, 1), "organization_id=store", "fixture", 400},
		{"duplicate-scope", valid, "organization_id=store&organization_id=other", "fixture", 400},
		{"other-scope", valid, "organization_id=other", "fixture", 403},
		{"oversized", `{"padding":"` + strings.Repeat("x", 32768) + `"}`, "organization_id=store", "fixture", 413},
		{"unauthenticated", valid, "organization_id=store", "invalid", 401},
	} {
		t.Run(test.name, func(t *testing.T) {
			service := &warrantyTestService{}
			tenant, _ := service.Scope()
			principal := identity.Principal{TenantID: tenant, Subject: "operator", Organizations: map[string]struct{}{"store": {}, "other": {}}, Permissions: map[string]struct{}{"*": {}}}
			mux := http.NewServeMux()
			WarrantyModule{Service: service}.Register(mux, warrantyTestVerifier{p: principal})
			req := httptest.NewRequest("POST", "/v1/franchise/warranty/quotes/quote/terms?"+test.query, strings.NewReader(test.body))
			req.Header.Set("Authorization", "Bearer "+test.token)
			res := httptest.NewRecorder()
			mux.ServeHTTP(res, req)
			if res.Code != test.code {
				t.Fatal(res.Code, res.Body.String())
			}
			if test.code != 201 {
				if service.calls != 0 {
					t.Fatal("rejected body reached owner")
				}
				return
			}
			if service.calls != 1 || service.request.QuoteVersion != 9007199254740993 || service.principal.Subject != "operator" || service.principal.AllowedOrganization("other") || service.principal.Allowed("warranty:approve") {
				t.Fatal("precision or authority changed", service)
			}
			var body map[string]any
			if json.Unmarshal(res.Body.Bytes(), &body) != nil || body["quote_version"] != "9007199254740993" {
				t.Fatal("response lost int64", res.Body.String())
			}
		})
	}
}
````

### FILE: `internal/platform/postgres/warranty_claim.go`

```yaml
block_id: "GO-CONNECTED-WARRANTY-CLAIM:file15:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "c615835afa1744d5dbf1f4ba457533495f9a0bedd0457b3e76936d786df44a76"
variables: []
secrets_allowed: false
```

````go
// AUTHORED transaction and evidence composition. Existing service cases, human
// approval registry and inventory costing writers retain their ownership.
package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"elite.local/enterprise/internal/approval"
	"elite.local/enterprise/internal/fulfillment"
	"elite.local/enterprise/internal/inventorycontrol"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/randomid"
	wc "elite.local/enterprise/internal/warrantyclaim"
	"github.com/jackc/pgx/v5"
)

type warrantyRowReader interface {
	QueryRow(context.Context, string, ...any) pgx.Row
}

func readWarrantyClaim(ctx context.Context, q warrantyRowReader, tenant, org, id string, lock bool) (wc.Claim, error) {
	var c wc.Claim
	sql := `select c.service_case_id,c.warranty_id,c.handover_id,c.appointment_id,c.organization_id,c.factory_organization_id,c.customer_principal_id,c.stock_unit_id,c.service_date::text,s.state,s.version,c.profile_sha256,c.parts_covered,c.labor_covered
 from service_ops.warranty_claim c join service_ops.service_case s using(tenant_id,service_case_id)
 where c.tenant_id=$1 and c.organization_id=$2 and c.service_case_id=$3`
	if lock {
		sql += " for update of s"
	}
	err := q.QueryRow(ctx, sql, tenant, org, id).Scan(&c.CaseID, &c.WarrantyID, &c.HandoverID, &c.AppointmentID, &c.OrganizationID, &c.FactoryOrganizationID, &c.CustomerSubject, &c.StockUnitID, &c.ServiceDate, &c.State, &c.Version, &c.ProfileSHA256, &c.PartsCovered, &c.LaborCovered)
	return c, warrantyError(err)
}
func readWarrantyStep(ctx context.Context, q warrantyRowReader, tenant, id, command, kind string) (wc.Step, error) {
	var s wc.Step
	err := q.QueryRow(ctx, `select service_case_id,command_id,version,kind,to_state,actor,request_sha256,payload_sha256,payload,recorded_at
 from service_ops.warranty_claim_step where tenant_id=$1 and service_case_id=$2 and ($3::text='' or command_id=$3) and ($4::text='' or kind=$4) order by version desc limit 1`, tenant, id, command, kind).Scan(&s.CaseID, &s.CommandID, &s.Version, &s.Kind, &s.State, &s.Actor, &s.RequestSHA256, &s.PayloadSHA256, &s.Payload, &s.RecordedAt)
	if err != nil {
		return s, warrantyError(err)
	}
	raw, hash, err := approval.CanonicalPayload(s.Payload)
	if err != nil || hash != s.PayloadSHA256 {
		return wc.Step{}, wc.ErrConflict
	}
	s.Payload = raw
	return s, nil
}
func (s *Warranty) claimAllowed(p identity.Principal, c wc.Claim, permission string) bool {
	tenant, org := s.profile.Scope()
	if c.OrganizationID != org {
		return false
	}
	switch permission {
	case "warranty:factory-read", "warranty:reconcile":
		return approvalPrincipal(p, tenant, c.FactoryOrganizationID, permission)
	case "warranty:self":
		return p.Subject == c.CustomerSubject && approvalPrincipal(p, tenant, org, permission)
	default:
		return approvalPrincipal(p, tenant, org, permission)
	}
}
func (s *Warranty) Claim(ctx context.Context, p identity.Principal, id string) (wc.Claim, error) {
	if s == nil || !s.profile.Valid() || !wc.ValidID(id) {
		return wc.Claim{}, wc.ErrInvalid
	}
	tenant, org := s.profile.Scope()
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead, AccessMode: pgx.ReadOnly})
	if err != nil {
		return wc.Claim{}, err
	}
	defer tx.Rollback(ctx)
	c, err := readWarrantyClaim(ctx, tx, tenant, org, id, false)
	if err != nil {
		return wc.Claim{}, err
	}
	if !(s.claimAllowed(p, c, "warranty:read") || s.claimAllowed(p, c, "warranty:self") || s.claimAllowed(p, c, "warranty:factory-read")) {
		return wc.Claim{}, wc.ErrNotFound
	}
	latest, err := readWarrantyStep(ctx, tx, tenant, id, "", "")
	if err != nil {
		return wc.Claim{}, err
	}
	c.Latest = &latest
	c.Context, err = readWarrantyRoleContext(ctx, tx, tenant, id)
	if err != nil {
		return wc.Claim{}, err
	}
	return c, tx.Commit(ctx)
}
func (s *Warranty) ClaimCommand(ctx context.Context, p identity.Principal, id, command string) (wc.Step, error) {
	if !wc.ValidID(command) {
		return wc.Step{}, wc.ErrInvalid
	}
	if _, err := s.Claim(ctx, p, id); err != nil {
		return wc.Step{}, err
	}
	tenant, _ := s.profile.Scope()
	return readWarrantyStep(ctx, s.pool, tenant, id, command, "")
}
func claimSoldProfile(ctx context.Context, q warrantyRowReader, tenant string, c wc.Claim) (wc.Profile, error) {
	var raw []byte
	var hash string
	err := q.QueryRow(ctx, `select o.profile_bytes,o.profile_sha256 from service_ops.warranty_activation a
 join service_ops.warranty_offer o on o.tenant_id=a.tenant_id and o.quotation_id=a.quotation_id
 where a.tenant_id=$1 and a.warranty_id=$2 and a.organization_id=$3 and a.profile_sha256=$4`, tenant, c.WarrantyID, c.OrganizationID, c.ProfileSHA256).Scan(&raw, &hash)
	if err != nil {
		return wc.Profile{}, warrantyError(err)
	}
	p, err := wc.LoadProfile(raw, hash)
	if err != nil {
		return wc.Profile{}, err
	}
	pt, po := p.Scope()
	if pt != tenant || po != c.OrganizationID || p.Document().FactoryOrganizationID != c.FactoryOrganizationID {
		return wc.Profile{}, wc.ErrConflict
	}
	return p, nil
}
func writeWarrantyStep(ctx context.Context, tx pgx.Tx, tenant string, c wc.Claim, command, kind, target, actor, requestSHA string, payload any, opening bool) (wc.Step, error) {
	raw, hash, err := wc.Canonical(payload)
	if err != nil {
		return wc.Step{}, err
	}
	version := c.Version + 1
	if opening {
		version = 1
	}
	_, err = tx.Exec(ctx, `insert into service_ops.warranty_claim_step(tenant_id,service_case_id,command_id,version,kind,from_state,to_state,actor,request_sha256,payload_sha256,payload)
 values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`, tenant, c.CaseID, command, version, kind, c.State, target, actor, requestSHA, hash, raw)
	if err != nil {
		return wc.Step{}, warrantyError(err)
	}
	if !opening {
		tag, err := tx.Exec(ctx, `update service_ops.service_case set state=$4,version=version+1,closed_at=case when $4='closed' then clock_timestamp() else null end
 where tenant_id=$1 and service_case_id=$2 and version=$3`, tenant, c.CaseID, c.Version, target)
		if err != nil {
			return wc.Step{}, err
		}
		if tag.RowsAffected() != 1 {
			return wc.Step{}, wc.ErrConflict
		}
	}
	result, err := readWarrantyStep(ctx, tx, tenant, c.CaseID, command, "")
	if err != nil {
		return wc.Step{}, err
	}
	if err = warrantyEvent(ctx, tx, tenant, c.CaseID, "warranty.claim-"+kind, version, result); err != nil {
		return wc.Step{}, err
	}
	return result, nil
}

type warrantyClaimEffect func(context.Context, pgx.Tx, string, wc.Claim) (kind, target string, payload any, err error)

func (s *Warranty) claimCommand(ctx context.Context, p identity.Principal, c wc.Command, permission string, request any, effect warrantyClaimEffect) (wc.Step, error) {
	if s == nil || !s.profile.Valid() || !c.Valid() {
		return wc.Step{}, wc.ErrInvalid
	}
	_, requestSHA, err := wc.Canonical(request)
	if err != nil {
		return wc.Step{}, err
	}
	tenant, org := s.profile.Scope()
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return wc.Step{}, err
	}
	defer tx.Rollback(ctx)
	claim, err := readWarrantyClaim(ctx, tx, tenant, org, c.CaseID, true)
	if err != nil {
		return wc.Step{}, err
	}
	if !s.claimAllowed(p, claim, permission) {
		return wc.Step{}, wc.ErrNotFound
	}
	existing, err := readWarrantyStep(ctx, tx, tenant, c.CaseID, c.CommandID, "")
	if err == nil {
		if existing.Actor != p.Subject || existing.RequestSHA256 != requestSHA {
			return wc.Step{}, wc.ErrConflict
		}
		existing.Replay = true
		return existing, tx.Commit(ctx)
	}
	if !errors.Is(err, wc.ErrNotFound) {
		return wc.Step{}, err
	}
	if claim.Version != c.ExpectedVersion {
		return wc.Step{}, wc.ErrConflict
	}
	kind, target, payload, err := effect(ctx, tx, tenant, claim)
	if err != nil {
		return wc.Step{}, err
	}
	result, err := writeWarrantyStep(ctx, tx, tenant, claim, c.CommandID, kind, target, p.Subject, requestSHA, payload, false)
	if err != nil {
		return wc.Step{}, err
	}
	return result, tx.Commit(ctx)
}
func (s *Warranty) OpenClaim(ctx context.Context, p identity.Principal, r wc.OpenClaim) (wc.Step, error) {
	if !(s.allowed(p, "warranty:request") || s.allowed(p, "warranty:self")) || !wc.ValidID(r.CaseID) || !wc.ValidID(r.HandoverID) || !wc.ValidID(r.AppointmentID) || !wc.ValidSHA(r.EvidenceSHA256) || len(strings.TrimSpace(r.Description)) < 3 || len(r.Description) > 4000 || !map[string]bool{"low": true, "medium": true, "high": true, "safety": true}[r.Severity] {
		return wc.Step{}, wc.ErrInvalid
	}
	_, hash, err := wc.Canonical(r)
	if err != nil {
		return wc.Step{}, err
	}
	tenant, org := s.profile.Scope()
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return wc.Step{}, err
	}
	defer tx.Rollback(ctx)
	if _, err = tx.Exec(ctx, `select pg_advisory_xact_lock(hashtextextended($1,0))`, tenant+":warranty-case:"+r.CaseID); err != nil {
		return wc.Step{}, err
	}
	existing, err := readWarrantyStep(ctx, tx, tenant, r.CaseID, "open", "")
	if err == nil {
		if _, scopeErr := readWarrantyClaim(ctx, tx, tenant, org, r.CaseID, false); scopeErr != nil {
			return wc.Step{}, wc.ErrNotFound
		}
		if existing.Actor != p.Subject || existing.RequestSHA256 != hash {
			return wc.Step{}, wc.ErrConflict
		}
		existing.Replay = true
		return existing, tx.Commit(ctx)
	}
	if !errors.Is(err, wc.ErrNotFound) {
		return wc.Step{}, err
	}
	a, err := readWarrantyActivation(ctx, tx, tenant, org, r.HandoverID)
	if err != nil {
		return wc.Step{}, err
	}
	if !s.allowed(p, "warranty:request") && p.Subject != a.CustomerSubject {
		return wc.Step{}, wc.ErrNotFound
	}
	var warrantyStatus string
	if err = tx.QueryRow(ctx, `select status from service_ops.warranty where tenant_id=$1 and warranty_id=$2 for share`, tenant, a.WarrantyID).Scan(&warrantyStatus); err != nil || warrantyStatus == "void" {
		return wc.Step{}, wc.ErrConflict
	}
	// The appointment and unit must belong to the same actual customer/model.
	var valid bool
	err = tx.QueryRow(ctx, `select true from crm.appointment ap join inventory.stock_unit u on u.tenant_id=ap.tenant_id
 join catalog.vehicle_variant v on v.tenant_id=u.tenant_id and v.variant_id=u.variant_id
 where ap.tenant_id=$1 and ap.organization_id=$2 and ap.appointment_id=$3 and ap.customer_principal_id=$4
 and ap.appointment_kind='service' and ap.state in ('confirmed','completed') and u.stock_unit_id=$5
 and (ap.model_id is null or ap.model_id=v.model_id) for share of ap`, tenant, org, r.AppointmentID, a.CustomerSubject, a.StockUnitID).Scan(&valid)
	if err != nil || !valid {
		return wc.Step{}, wc.ErrConflict
	}
	var raw []byte
	err = tx.QueryRow(ctx, `select profile_bytes from service_ops.warranty_offer where tenant_id=$1 and quotation_id=$2`, tenant, a.QuoteID).Scan(&raw)
	if err != nil {
		return wc.Step{}, err
	}
	sold, err := wc.LoadProfile(raw, a.ProfileSHA256)
	if err != nil {
		return wc.Step{}, err
	}
	doc := sold.Document()
	zone, err := time.LoadLocation(doc.BusinessTimeZone)
	if err != nil {
		return wc.Step{}, err
	}
	var now time.Time
	if err = tx.QueryRow(ctx, `select clock_timestamp()`).Scan(&now); err != nil {
		return wc.Step{}, err
	}
	date := now.In(zone).Format(time.DateOnly)
	coverage, err := wc.CheckCoverage(date, a.Dates)
	if err != nil {
		return wc.Step{}, err
	}
	id := randomid.Generator{}
	if err = openServiceCaseInTx(ctx, tx, tenant, id.New(), fulfillment.ServiceCase{ID: r.CaseID, OrganizationID: org, StockUnitID: a.StockUnitID, Severity: r.Severity, Description: r.Description}); err != nil {
		return wc.Step{}, err
	}
	_, err = tx.Exec(ctx, `insert into service_ops.warranty_claim(tenant_id,service_case_id,warranty_id,handover_id,appointment_id,organization_id,factory_organization_id,customer_principal_id,stock_unit_id,service_date,profile_sha256,parts_covered,labor_covered,opened_by,request_sha256)
 values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10::date,$11,$12,$13,$14,$15)`, tenant, r.CaseID, a.WarrantyID, r.HandoverID, r.AppointmentID, org, doc.FactoryOrganizationID, a.CustomerSubject, a.StockUnitID, date, sold.Hash(), coverage.Parts, coverage.Labor, p.Subject, hash)
	if err != nil {
		return wc.Step{}, warrantyError(err)
	}
	claim, err := readWarrantyClaim(ctx, tx, tenant, org, r.CaseID, false)
	if err != nil {
		return wc.Step{}, err
	}
	step, err := writeWarrantyStep(ctx, tx, tenant, claim, "open", "opened", "opened", p.Subject, hash, map[string]any{"request": r, "coverage": coverage, "profile_sha256": sold.Hash(), "service_date": date}, true)
	if err != nil {
		return wc.Step{}, err
	}
	return step, tx.Commit(ctx)
}
func (s *Warranty) Diagnose(ctx context.Context, p identity.Principal, r wc.Diagnose) (wc.Step, error) {
	if !wc.ValidID(r.FaultCode) || !wc.ValidSHA(r.EvidenceSHA256) || len(strings.TrimSpace(r.Description)) < 3 || len(r.Description) > 4000 {
		return wc.Step{}, wc.ErrInvalid
	}
	return s.claimCommand(ctx, p, r.Command, "warranty:diagnose", r, func(ctx context.Context, tx pgx.Tx, tenant string, c wc.Claim) (string, string, any, error) {
		if c.State != "opened" {
			return "", "", nil, wc.ErrConflict
		}
		sold, err := claimSoldProfile(ctx, tx, tenant, c)
		if err != nil {
			return "", "", nil, err
		}
		return "diagnosed", "diagnosis", map[string]any{"diagnosis": r, "fault_excluded": sold.FaultExcluded(r.FaultCode)}, nil
	})
}
func readWarrantyPlan(ctx context.Context, q warrantyRowReader, tenant, id string) (wc.RepairPlan, string, string, error) {
	var plan wc.RepairPlan
	var raw []byte
	var approvalID, hash string
	err := q.QueryRow(ctx, `select approval_id,payload_sha256,payload from service_ops.warranty_work_plan where tenant_id=$1 and service_case_id=$2`, tenant, id).Scan(&approvalID, &hash, &raw)
	if err != nil {
		return plan, "", "", warrantyError(err)
	}
	_, got, err := approval.CanonicalPayload(raw)
	if err != nil || got != hash || json.Unmarshal(raw, &plan) != nil || plan.ClaimID != id {
		return plan, "", "", wc.ErrConflict
	}
	return plan, approvalID, hash, nil
}
func (s *Warranty) PlanRepair(ctx context.Context, p identity.Principal, r wc.Plan) (wc.Step, error) {
	if !r.Valid() {
		return wc.Step{}, wc.ErrInvalid
	}
	return s.claimCommand(ctx, p, r.Command, "warranty:plan", r, func(ctx context.Context, tx pgx.Tx, tenant string, c wc.Claim) (string, string, any, error) {
		if c.State != "diagnosis" {
			return "", "", nil, wc.ErrConflict
		}
		diagnosis, err := readWarrantyStep(ctx, tx, tenant, c.CaseID, "", "diagnosed")
		if err != nil {
			return "", "", nil, err
		}
		var d struct {
			FaultExcluded bool `json:"fault_excluded"`
		}
		if json.Unmarshal(diagnosis.Payload, &d) != nil {
			return "", "", nil, wc.ErrConflict
		}
		sold, err := claimSoldProfile(ctx, tx, tenant, c)
		if err != nil {
			return "", "", nil, err
		}
		// An uncovered diagnosis can be submitted for an explicit human rejection;
		// it never reserves parts or permits an approved covered-service claim.
		covered := !d.FaultExcluded && (c.PartsCovered || c.LaborCovered)
		if covered && (len(r.Parts) > 0 && !c.PartsCovered || r.LaborWork != "" && !c.LaborCovered) {
			return "", "", nil, wc.ErrConflict
		}
		var expires time.Time
		if err = tx.QueryRow(ctx, `select clock_timestamp()+make_interval(secs=>$1)`, sold.Document().WorkReservationSeconds).Scan(&expires); err != nil {
			return "", "", nil, err
		}
		plan := wc.RepairPlan{Schema: "elite-warranty-repair-plan/v1", ClaimID: c.CaseID, ProfileSHA256: c.ProfileSHA256, DiagnosisSHA256: diagnosis.PayloadSHA256, Requester: p.Subject, Request: r, PartsCovered: c.PartsCovered, LaborCovered: c.LaborCovered, Excluded: d.FaultExcluded, ExpiresAt: expires, Reservations: []inventorycontrol.BulkReservation{}}
		ids := randomid.Generator{}
		if covered {
			for _, part := range r.Parts {
				reservation, err := reserveBulkInTx(ctx, tx, tenant, ids.New(), inventorycontrol.BulkReservation{ID: wc.StableID("warranty-part", tenant, c.CaseID, part.LineID), OrganizationID: c.OrganizationID, BinID: part.BinID, ItemID: part.ItemID, LotID: part.LotID, DemandKind: "service", DemandID: c.CaseID, DemandLineID: part.LineID, Quantity: part.Quantity, Status: "reservation", Version: 1, ExpiresAt: &expires})
				if err != nil {
					return "", "", nil, err
				}
				plan.Reservations = append(plan.Reservations, reservation)
			}
		}
		raw, hash, err := wc.Canonical(plan)
		if err != nil {
			return "", "", nil, err
		}
		approvalID := wc.StableID("warranty-approval", tenant, c.CaseID)
		spec := HumanApprovalSpec{Request: approval.Request{TenantID: tenant, ID: approvalID, Kind: approval.KindWarrantyRepair, SubjectID: c.CaseID, Requester: p.Subject, EvidenceSHA: hash}, OrganizationID: c.OrganizationID, Payload: raw}
		if _, err = NewHumanApprovals(s.pool).submitTx(ctx, tx, p, spec, "warranty:plan", nil); err != nil {
			return "", "", nil, err
		}
		if _, err = tx.Exec(ctx, `insert into service_ops.warranty_work_plan(tenant_id,service_case_id,approval_id,payload_sha256,payload,expires_at) values($1,$2,$3,$4,$5,$6)`, tenant, c.CaseID, approvalID, hash, raw, expires); err != nil {
			return "", "", nil, err
		}
		for _, reservation := range plan.Reservations {
			if _, err = tx.Exec(ctx, `insert into service_ops.warranty_part_reservation(tenant_id,service_case_id,line_id,reservation_id) values($1,$2,$3,$4)`, tenant, c.CaseID, reservation.DemandLineID, reservation.ID); err != nil {
				return "", "", nil, err
			}
		}
		return "planned", "diagnosis", map[string]any{"approval_id": approvalID, "payload_sha256": hash, "plan": plan}, nil
	})
}
````

### FILE: `internal/platform/postgres/warranty_claim_integration_test.go`

```yaml
block_id: "GO-CONNECTED-WARRANTY-CLAIM:file16:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "539c5e2aee6f7b950d2641726f87b22bfa674554b931d7eaba5dd424d6ec2e98"
variables: []
secrets_allowed: false
```

````go
package postgres_test

import (
	"context"
	"encoding/json"
	"math/big"
	"strings"
	"testing"
	"time"

	"elite.local/enterprise/internal/businesspolicy"
	"elite.local/enterprise/internal/franchisejourney"
	"elite.local/enterprise/internal/inventorycontrol"
	"elite.local/enterprise/internal/platform/identity"
	db "elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/platform/randomid"
	wc "elite.local/enterprise/internal/warrantyclaim"
)

type warrantyFixtureClock struct{}

func (warrantyFixtureClock) Now() time.Time { return time.Now().UTC() }

func assertWarrantyClaimLifecycle(t *testing.T, r *connectedRun, store warrantyClaimCommands, operator, customer identity.Principal, activation wc.Activation) {
	t.Helper()
	ctx := context.Background()
	ids := randomid.Generator{}
	must := func(q string, args ...any) {
		t.Helper()
		if _, err := r.pool.Exec(ctx, q, args...); err != nil {
			t.Fatal(err)
		}
	}
	must(`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'factory','factory','Synthetic factory','factory')`, r.tenant)
	must(`insert into org.public_location(tenant_id,organization_id,city,region,country,published)values($1,'store','Synthetic','Synthetic','AR',true)`, r.tenant)
	// Explicit synthetic appointment profile, admitted by the existing policy
	// contract: zero lead time and one-second slots avoid wall-clock simulation.
	// This operational profile does not amend any warranty terms already sold.
	var doc map[string]any
	if json.Unmarshal(businesspolicy.ReferenceJSON(), &doc) != nil {
		t.Fatal("reference profile")
	}
	doc["profile_id"] = "warranty-fixture"
	doc["appointments"].(map[string]any)["lead_time_seconds"] = 0
	raw, err := json.Marshal(doc)
	if err != nil {
		t.Fatal(err)
	}
	policy, err := businesspolicy.Load(raw, wc.SHA(raw))
	if err != nil {
		t.Fatal(err)
	}
	journey, err := db.NewFranchiseJourneyWithProfile(r.pool, policy)
	if err != nil {
		t.Fatal(err)
	}
	booking, err := franchisejourney.NewServiceWithProfile(journey, ids, warrantyFixtureClock{}, policy)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = journey.CreateServiceResource(ctx, r.tenant, franchisejourney.ServiceResource{ID: "warranty-bay", OrganizationID: "store", Kind: "service-bay", DisplayName: "Synthetic service bay", Skills: []string{"service"}}, ids.New()); err != nil {
		t.Fatal(err)
	}
	for _, resource := range []string{"", "warranty-bay"} {
		must(`insert into crm.availability_entry(tenant_id,availability_id,organization_id,resource_id,entry_type,starts_at,ends_at,state,version,created_by_subject)
 values($1,$2,'store',nullif($3,''),'working',clock_timestamp()-interval '1 hour',clock_timestamp()+interval '6 hours','active',1,'fixture-calendar')`, r.tenant, ids.New(), resource)
	}
	var tenantCode string
	if err = r.pool.QueryRow(ctx, `select tenant_code from platform.tenant where tenant_id=$1`, r.tenant).Scan(&tenantCode); err != nil {
		t.Fatal(err)
	}
	appointment := func(id string) string {
		t.Helper()
		var start time.Time
		if err := r.pool.QueryRow(ctx, `select clock_timestamp()+interval '1 second'`).Scan(&start); err != nil {
			t.Fatal(err)
		}
		if _, err := booking.CreateAppointmentSlot(ctx, r.tenant, franchisejourney.AppointmentSlot{OrganizationID: "store", Kind: "service", StartsAt: start, EndsAt: start.Add(time.Second), Capacity: 1}); err != nil {
			t.Fatal(err)
		}
		a, _, err := booking.RequestAppointment(ctx, tenantCode, "store", id+"-appointment-booking-key", wc.SHA([]byte(id)), franchisejourney.Appointment{LeadID: "lead", ModelID: "model", Kind: "service", StartsAt: start})
		if err != nil {
			t.Fatal("request appointment", err)
		}
		a, err = journey.AssignAppointmentResource(ctx, r.tenant, "store", a.ID, "warranty-bay", a.Version, ids.New())
		if err != nil {
			t.Fatal("assign service bay", err)
		}
		a, err = journey.TransitionAppointment(ctx, r.tenant, "store", a.ID, "requested", "confirmed", a.Version, "operator", "", ids.New())
		if err != nil {
			t.Fatal("confirm appointment", err)
		}
		wait := time.Until(start.Add(time.Second + 30*time.Millisecond))
		if wait > 0 {
			if wait > 3*time.Second {
				t.Fatal("fixture wait exceeds bound")
			}
			time.Sleep(wait)
		}
		a, err = journey.TransitionAppointment(ctx, r.tenant, "store", a.ID, "confirmed", "completed", a.Version, "operator", "", ids.New())
		if err != nil {
			t.Fatal("complete attended appointment", err)
		}
		return a.ID
	}
	inv := db.NewInventoryControl(r.pool)
	if _, err = inv.CreateBulkItem(ctx, r.tenant, ids.New(), inventorycontrol.BulkItem{ID: "warranty-part", Code: "WARRANTY-PART", Description: "Synthetic replacement part", BaseUOM: "EA", BaseRoundingPrecision: "1", TrackingMode: "lot", CostingMethod: "fifo", Version: 1}); err != nil {
		t.Fatal(err)
	}
	if _, err = inv.CreateWarehouseBin(ctx, r.tenant, ids.New(), inventorycontrol.WarehouseBin{ID: "warranty-bin", OrganizationID: "store", Code: "SERVICE", Type: "pick", Version: 1}); err != nil {
		t.Fatal(err)
	}
	if _, err = inv.ConfigureItemBin(ctx, r.tenant, ids.New(), inventorycontrol.ItemBinPolicy{OrganizationID: "store", ItemID: "warranty-part", BinID: "warranty-bin", Fixed: true, Default: true, MinQuantity: "0", MaxQuantity: "100", Version: 1}); err != nil {
		t.Fatal(err)
	}
	expires := time.Now().UTC().AddDate(2, 0, 0)
	for i, cost := range []string{"100", "120"} {
		_, err = inv.ReceiveBulk(ctx, r.tenant, ids.New(), "warranty-receipt-"+cost, "warranty-layer-"+cost, "warranty-lot", inventorycontrol.BulkReceipt{OrganizationID: "store", BinID: "warranty-bin", ItemID: "warranty-part", LotNo: "SERVICE-FIXTURE", ExpirationDate: &expires, Quantity: "5", UnitCost: cost, PostingDate: time.Now().UTC().AddDate(0, 0, -2+i), SourceKind: "fixture-receipt", SourceID: "fixture-" + cost})
		if err != nil {
			t.Fatal(err)
		}
	}
	for _, permission := range []string{"warranty:request", "warranty:diagnose", "warranty:plan", "warranty:work", "warranty:approve", "warranty:quality", "warranty:cancel"} {
		operator.Permissions[permission] = struct{}{}
	}
	reviewer := operator
	reviewer.Subject = "warranty-reviewer"
	qualityActor := operator
	qualityActor.Subject = "warranty-quality"
	factory := identity.Principal{TenantID: r.tenant, Subject: "factory-reviewer", Organizations: map[string]struct{}{"factory": {}}, Permissions: map[string]struct{}{"warranty:factory-read": {}, "warranty:reconcile": {}}}
	part := func(qty string) []wc.Part {
		return []wc.Part{{LineID: "brake", ItemID: "warranty-part", BinID: "warranty-bin", LotID: "warranty-lot", Quantity: qty}}
	}
	open := func(id, fault string) (wc.Step, wc.Diagnose) {
		t.Helper()
		request := wc.OpenClaim{CaseID: id, HandoverID: activation.HandoverID, AppointmentID: appointment(id), Severity: "medium", Description: "Synthetic brake diagnosis request", EvidenceSHA256: strings.Repeat("a", 64)}
		opened, err := store.OpenClaim(ctx, customer, request)
		if err != nil {
			t.Fatal("open claim", err)
		}
		again, err := store.OpenClaim(ctx, customer, request)
		if err != nil || !again.Replay || again.Version != opened.Version {
			t.Fatal("open recovery", again, err)
		}
		diagnose := wc.Diagnose{Command: wc.Command{CaseID: id, CommandID: "diagnose", ExpectedVersion: 1}, FaultCode: fault, Description: "Synthetic inspection of brake", EvidenceSHA256: strings.Repeat("b", 64)}
		step, err := store.Diagnose(ctx, operator, diagnose)
		if err != nil {
			t.Fatal("diagnosis", err)
		}
		return step, diagnose
	}
	planned := func(id, qty string, version int64) (wc.Step, string) {
		t.Helper()
		step, err := store.PlanRepair(ctx, operator, wc.Plan{Command: wc.Command{CaseID: id, CommandID: "plan", ExpectedVersion: version}, LaborWork: "Synthetic brake adjustment", Parts: part(qty)})
		if err != nil {
			t.Fatal("plan", err)
		}
		var bound struct {
			SHA string `json:"payload_sha256"`
		}
		if json.Unmarshal(step.Payload, &bound) != nil || !wc.ValidSHA(bound.SHA) {
			t.Fatal("plan hash")
		}
		return step, bound.SHA
	}
	step, _ := open("claim-main", "fixture-wear")
	step, planSHA := planned("claim-main", "7", step.Version)
	decision := wc.DecideRepair{Command: wc.Command{CaseID: "claim-main", CommandID: "approve", ExpectedVersion: step.Version}, PayloadSHA256: planSHA, Approved: true, Reason: "Fixture covered repair"}
	if _, err = store.DecideRepair(ctx, operator, decision); err == nil {
		t.Fatal("requester approved own repair")
	}
	// The old generic case endpoint cannot bypass the connected approval writer.
	if err = db.NewFulfillment(r.pool).TransitionServiceCase(ctx, r.tenant, "store", "claim-main", "diagnosis", "repair", step.Version, ids.New()); err == nil {
		t.Fatal("generic service path bypassed approval")
	}
	step, err = store.DecideRepair(ctx, reviewer, decision)
	if err != nil || step.State != "repair" {
		t.Fatal("human approval", step, err)
	}
	again, err := store.DecideRepair(ctx, reviewer, decision)
	if err != nil || !again.Replay {
		t.Fatal("decision recovery", again, err)
	}
	work := wc.CompleteWork{Command: wc.Command{CaseID: "claim-main", CommandID: "work", ExpectedVersion: step.Version}, EvidenceSHA256: strings.Repeat("c", 64)}
	// Failure after real FIFO stock writes must roll back every stock/case effect.
	must(`create function service_ops.fixture_work_failure() returns trigger language plpgsql as $$begin if new.event_type='warranty.claim-work' then raise exception 'fixture work commit failure';end if;return new;end $$;
 create trigger fixture_work_failure before insert on platform.outbox_event for each row execute function service_ops.fixture_work_failure()`)
	if _, err = store.CompleteRepairWork(ctx, operator, work); err == nil {
		t.Fatal("work fault injection did not fail")
	}
	var quantity, reserved string
	var issues int
	if err = r.pool.QueryRow(ctx, `select quantity::text,reserved_quantity::text,(select count(*) from inventory.bulk_inventory_entry where tenant_id=$1 and entry_type='issue') from inventory.bulk_balance where tenant_id=$1 and item_id='warranty-part'`, r.tenant).Scan(&quantity, &reserved, &issues); err != nil || quantity != "10.000000" || reserved != "7.000000" || issues != 0 {
		t.Fatal("FIFO rollback", quantity, reserved, issues, err)
	}
	must(`drop trigger fixture_work_failure on platform.outbox_event;drop function service_ops.fixture_work_failure()`)
	step, err = store.CompleteRepairWork(ctx, operator, work)
	if err != nil || step.State != "quality" {
		t.Fatal("complete repair", step, err)
	}
	var performed wc.WorkReceipt
	if json.Unmarshal(step.Payload, &performed) != nil || len(performed.Parts) != 1 || !warrantyExactCost(performed.Parts[0].Issue.CostAmount, "740") || performed.Parts[0].Issue.Applications != 2 {
		t.Fatal("source FIFO receipt", string(step.Payload))
	}
	workSHA := step.PayloadSHA256
	again, err = store.CompleteRepairWork(ctx, operator, work)
	if err != nil || !again.Replay || again.PayloadSHA256 != workSHA {
		t.Fatal("work recovery", again, err)
	}
	failedQuality := wc.Quality{Command: wc.Command{CaseID: "claim-main", CommandID: "quality-failed", ExpectedVersion: step.Version}, Passed: false, WorkEvidenceSHA256: workSHA, EvidenceSHA256: strings.Repeat("d", 64)}
	if _, err = store.RecordRepairQuality(ctx, operator, failedQuality); err == nil {
		t.Fatal("technician self-verified work")
	}
	step, err = store.RecordRepairQuality(ctx, qualityActor, failedQuality)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = store.AcceptRepair(ctx, customer, wc.AcceptRepair{Command: wc.Command{CaseID: "claim-main", CommandID: "accept-failed", ExpectedVersion: step.Version}, QualitySHA256: step.PayloadSHA256, EvidenceSHA256: strings.Repeat("e", 64)}); err == nil {
		t.Fatal("failed quality accepted")
	}
	corrected := wc.Quality{Command: wc.Command{CaseID: "claim-main", CommandID: "quality-passed", ExpectedVersion: step.Version}, Passed: true, WorkEvidenceSHA256: workSHA, EvidenceSHA256: strings.Repeat("f", 64), CorrectionEvidenceSHA256: strings.Repeat("1", 64)}
	step, err = store.RecordRepairQuality(ctx, qualityActor, corrected)
	if err != nil {
		t.Fatal("corrected quality", err)
	}
	accept := wc.AcceptRepair{Command: wc.Command{CaseID: "claim-main", CommandID: "accept", ExpectedVersion: step.Version}, QualitySHA256: step.PayloadSHA256, EvidenceSHA256: strings.Repeat("2", 64)}
	wrong := customer
	wrong.Subject = "different-customer"
	if _, err = store.AcceptRepair(ctx, wrong, accept); err == nil {
		t.Fatal("wrong customer accepted repair")
	}
	step, err = store.AcceptRepair(ctx, customer, accept)
	if err != nil {
		t.Fatal("customer acceptance", err)
	}
	reconcile := wc.ReconcileRepair{Command: wc.Command{CaseID: "claim-main", CommandID: "reconcile", ExpectedVersion: step.Version}, AcceptanceSHA256: step.PayloadSHA256, EvidenceSHA256: strings.Repeat("3", 64)}
	if _, err = store.ReconcileRepair(ctx, operator, reconcile); err == nil {
		t.Fatal("store performed factory acknowledgement")
	}
	step, err = store.ReconcileRepair(ctx, factory, reconcile)
	if err != nil || step.State != "closed" {
		t.Fatal("factory reconciliation", step, err)
	}
	var receipt struct {
		Cost string `json:"recorded_inventory_cost"`
		Paid bool   `json:"payment_created"`
	}
	if json.Unmarshal(step.Payload, &receipt) != nil || receipt.Cost != "740.0000" || receipt.Paid {
		t.Fatal("invented settlement", string(step.Payload))
	}
	again, err = store.ReconcileRepair(ctx, factory, reconcile)
	if err != nil || !again.Replay {
		t.Fatal("reconciliation recovery", again, err)
	}
	// Declined covered service and cancelled approved service both release the
	// original reservations without writing an issue or erasing human decisions.
	for _, scenario := range []string{"declined", "cancelled", "excluded"} {
		fault := "fixture-wear"
		if scenario == "excluded" {
			fault = "fixture-exclusion"
		}
		st, _ := open("claim-"+scenario, fault)
		st, hash := planned("claim-"+scenario, "1", st.Version)
		d := wc.DecideRepair{Command: wc.Command{CaseID: "claim-" + scenario, CommandID: "decision", ExpectedVersion: st.Version}, PayloadSHA256: hash, Approved: scenario == "cancelled", Reason: "Synthetic decision"}
		if scenario == "excluded" {
			bad := d
			bad.Approved = true
			if _, err = store.DecideRepair(ctx, reviewer, bad); err == nil {
				t.Fatal("excluded fault approved")
			}
		}
		st, err = store.DecideRepair(ctx, reviewer, d)
		if err != nil {
			t.Fatal(scenario, err)
		}
		if scenario == "cancelled" {
			st, err = store.CancelRepair(ctx, reviewer, wc.CancelRepair{Command: wc.Command{CaseID: st.CaseID, CommandID: "cancel", ExpectedVersion: st.Version}, Reason: "Customer withdrew before stock issue", EvidenceSHA256: strings.Repeat("4", 64)})
			if err != nil {
				t.Fatal("cancel approved unperformed work", err)
			}
		}
		if st.State != "cancelled" {
			t.Fatal(scenario, st)
		}
	}
	if err = r.pool.QueryRow(ctx, `select quantity::text,reserved_quantity::text,(select count(*) from inventory.bulk_inventory_entry where tenant_id=$1 and entry_type='issue') from inventory.bulk_balance where tenant_id=$1 and item_id='warranty-part'`, r.tenant).Scan(&quantity, &reserved, &issues); err != nil || quantity != "3.000000" || reserved != "0.000000" || issues != 1 {
		t.Fatal("final stock", quantity, reserved, issues, err)
	}
	t.Logf("WARRANTY_CLAIM_J4_PASS real_appointment_request_assignment_confirmation_completion=true appointment_profile_sha256=%s diagnosis_approval_parts_fifo_quality_customer_factory=true claims=4 stock_issues=1 fifo_applications=2 source_cost=740.0000 atomic_rollback=true rejected_and_cancelled_holds_released=true live_proven=false", policy.SHA256())
}

func warrantyExactCost(actual, want string) bool {
	a, ok := new(big.Rat).SetString(actual)
	if !ok {
		return false
	}
	b, ok := new(big.Rat).SetString(want)
	return ok && a.Cmp(b) == 0
}
````

### FILE: `internal/platform/postgres/warranty_claim_work.go`

```yaml
block_id: "GO-CONNECTED-WARRANTY-CLAIM:file17:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "1452f707dbe48845e0fd9fcfefb3136585c0ed3459924c7553977aabd970cfb8"
variables: []
secrets_allowed: false
```

````go
package postgres

// AUTHORED same-transaction service/approval/inventory composition. Repair
// acknowledgements never represent a paid reimbursement or a tax document.
import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"elite.local/enterprise/internal/approval"
	"elite.local/enterprise/internal/inventorycontrol"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/randomid"
	wc "elite.local/enterprise/internal/warrantyclaim"
	"github.com/jackc/pgx/v5"
)

func (s *Warranty) DecideRepair(ctx context.Context, p identity.Principal, r wc.DecideRepair) (wc.Step, error) {
	if !wc.ValidSHA(r.PayloadSHA256) || len(strings.TrimSpace(r.Reason)) < 3 || len(r.Reason) > 2048 {
		return wc.Step{}, wc.ErrInvalid
	}
	return s.claimCommand(ctx, p, r.Command, "warranty:approve", r, func(ctx context.Context, tx pgx.Tx, tenant string, c wc.Claim) (string, string, any, error) {
		if c.State != "diagnosis" {
			return "", "", nil, wc.ErrConflict
		}
		plan, id, hash, err := readWarrantyPlan(ctx, tx, tenant, c.CaseID)
		if err != nil {
			return "", "", nil, err
		}
		if hash != r.PayloadSHA256 || plan.ProfileSHA256 != c.ProfileSHA256 {
			return "", "", nil, wc.ErrConflict
		}
		if r.Approved && (plan.Excluded || !plan.PartsCovered && !plan.LaborCovered) {
			return "", "", nil, wc.ErrConflict
		}
		if _, err = NewHumanApprovals(s.pool).decideTx(ctx, tx, p, tenant, id, c.OrganizationID, hash, r.Approved, r.Reason, "warranty:approve", nil); err != nil {
			return "", "", nil, err
		}
		ids := randomid.Generator{}
		for _, reservation := range plan.Reservations {
			if r.Approved {
				// Approval converts the pending lease to an explicit durable repair hold.
				// Only metadata/version changes; the existing stock reservation remains.
				tag, err := tx.Exec(ctx, `update inventory.bulk_reservation set expires_at=null,version=version+1,updated_at=clock_timestamp()
 where tenant_id=$1 and organization_id=$2 and reservation_id=$3 and version=$4 and status='reservation' and expires_at=$5 and expires_at>clock_timestamp()`, tenant, c.OrganizationID, reservation.ID, reservation.Version, plan.ExpiresAt)
				if err != nil {
					return "", "", nil, err
				}
				if tag.RowsAffected() != 1 {
					return "", "", nil, wc.ErrConflict
				}
			} else if err = releaseBulkInTx(ctx, tx, tenant, c.OrganizationID, reservation.ID, reservation.Version, ids.New()); err != nil {
				return "", "", nil, err
			}
		}
		target := "cancelled"
		if r.Approved {
			target = "repair"
		}
		return "decision", target, map[string]any{"decision": r, "approval_id": id, "approved_sha256": hash, "requester": plan.Requester, "reviewer": p.Subject}, nil
	})
}
func (s *Warranty) CompleteRepairWork(ctx context.Context, p identity.Principal, r wc.CompleteWork) (wc.Step, error) {
	if !wc.ValidSHA(r.EvidenceSHA256) {
		return wc.Step{}, wc.ErrInvalid
	}
	return s.claimCommand(ctx, p, r.Command, "warranty:work", r, func(ctx context.Context, tx pgx.Tx, tenant string, c wc.Claim) (string, string, any, error) {
		if c.State != "repair" {
			return "", "", nil, wc.ErrConflict
		}
		plan, id, hash, err := readWarrantyPlan(ctx, tx, tenant, c.CaseID)
		if err != nil {
			return "", "", nil, err
		}
		if _, err = LookupHumanApproval(ctx, tx, tenant, id, c.OrganizationID, approval.KindWarrantyRepair, hash); err != nil {
			return "", "", nil, err
		}
		if plan.Excluded || !plan.PartsCovered && !plan.LaborCovered || len(plan.Reservations) != len(plan.Request.Parts) {
			return "", "", nil, wc.ErrConflict
		}
		receipt := wc.WorkReceipt{Request: r, ApprovalID: id, ApprovedSHA256: hash, Parts: []wc.IssuedPart{}}
		sold, err := claimSoldProfile(ctx, tx, tenant, c)
		if err != nil {
			return "", "", nil, err
		}
		zone, err := time.LoadLocation(sold.Document().BusinessTimeZone)
		if err != nil {
			return "", "", nil, err
		}
		var now time.Time
		if err = tx.QueryRow(ctx, `select clock_timestamp()`).Scan(&now); err != nil {
			return "", "", nil, err
		}
		local := now.In(zone)
		posting := time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, time.UTC)
		ids := randomid.Generator{}
		for i, part := range plan.Request.Parts {
			reservation := plan.Reservations[i]
			if reservation.DemandLineID != part.LineID || reservation.ItemID != part.ItemID || reservation.BinID != part.BinID {
				return "", "", nil, wc.ErrConflict
			}
			result, err := issueBulkInTx(ctx, tx, tenant, ids.New(), ids.New(), inventorycontrol.BulkIssue{OrganizationID: c.OrganizationID, ReservationID: reservation.ID, ReservationVersion: reservation.Version + 1, SpecificReceiptEntry: part.SpecificReceiptEntry, PostingDate: posting, SourceKind: "warranty-repair", SourceID: c.CaseID})
			if err != nil {
				return "", "", nil, err
			}
			receipt.Parts = append(receipt.Parts, wc.IssuedPart{LineID: part.LineID, ReservationID: reservation.ID, Issue: result})
		}
		return "work", "quality", receipt, nil
	})
}
func (s *Warranty) RecordRepairQuality(ctx context.Context, p identity.Principal, r wc.Quality) (wc.Step, error) {
	if !wc.ValidSHA(r.EvidenceSHA256) || !wc.ValidSHA(r.WorkEvidenceSHA256) || r.CorrectionEvidenceSHA256 != "" && !wc.ValidSHA(r.CorrectionEvidenceSHA256) {
		return wc.Step{}, wc.ErrInvalid
	}
	return s.claimCommand(ctx, p, r.Command, "warranty:quality", r, func(ctx context.Context, tx pgx.Tx, tenant string, c wc.Claim) (string, string, any, error) {
		if c.State != "quality" {
			return "", "", nil, wc.ErrConflict
		}
		work, err := readWarrantyStep(ctx, tx, tenant, c.CaseID, "", "work")
		if err != nil {
			return "", "", nil, err
		}
		if work.PayloadSHA256 != r.WorkEvidenceSHA256 || work.Actor == p.Subject {
			return "", "", nil, wc.ErrConflict
		}
		if _, err = readWarrantyStep(ctx, tx, tenant, c.CaseID, "", "accepted"); err == nil {
			return "", "", nil, wc.ErrConflict
		} else if !errors.Is(err, wc.ErrNotFound) {
			return "", "", nil, err
		}
		previous, err := readWarrantyStep(ctx, tx, tenant, c.CaseID, "", "quality")
		if err == nil {
			var q wc.Quality
			if json.Unmarshal(previous.Payload, &q) != nil || q.Passed || !wc.ValidSHA(r.CorrectionEvidenceSHA256) {
				return "", "", nil, wc.ErrConflict
			}
		} else if !errors.Is(err, wc.ErrNotFound) {
			return "", "", nil, err
		}
		return "quality", "quality", r, nil
	})
}
func (s *Warranty) AcceptRepair(ctx context.Context, p identity.Principal, r wc.AcceptRepair) (wc.Step, error) {
	if !wc.ValidSHA(r.QualitySHA256) || !wc.ValidSHA(r.EvidenceSHA256) {
		return wc.Step{}, wc.ErrInvalid
	}
	return s.claimCommand(ctx, p, r.Command, "warranty:self", r, func(ctx context.Context, tx pgx.Tx, tenant string, c wc.Claim) (string, string, any, error) {
		if c.State != "quality" {
			return "", "", nil, wc.ErrConflict
		}
		quality, err := readWarrantyStep(ctx, tx, tenant, c.CaseID, "", "quality")
		if err != nil {
			return "", "", nil, err
		}
		var q wc.Quality
		if json.Unmarshal(quality.Payload, &q) != nil || !q.Passed || quality.PayloadSHA256 != r.QualitySHA256 {
			return "", "", nil, wc.ErrConflict
		}
		return "accepted", "quality", r, nil
	})
}
func (s *Warranty) ReconcileRepair(ctx context.Context, p identity.Principal, r wc.ReconcileRepair) (wc.Step, error) {
	if !wc.ValidSHA(r.AcceptanceSHA256) || !wc.ValidSHA(r.EvidenceSHA256) {
		return wc.Step{}, wc.ErrInvalid
	}
	return s.claimCommand(ctx, p, r.Command, "warranty:reconcile", r, func(ctx context.Context, tx pgx.Tx, tenant string, c wc.Claim) (string, string, any, error) {
		if c.State != "quality" {
			return "", "", nil, wc.ErrConflict
		}
		accepted, err := readWarrantyStep(ctx, tx, tenant, c.CaseID, "", "accepted")
		if err != nil {
			return "", "", nil, err
		}
		if accepted.PayloadSHA256 != r.AcceptanceSHA256 || accepted.Actor != c.CustomerSubject {
			return "", "", nil, wc.ErrConflict
		}
		quality, err := readWarrantyStep(ctx, tx, tenant, c.CaseID, "", "quality")
		if err != nil {
			return "", "", nil, err
		}
		var ack wc.AcceptRepair
		var q wc.Quality
		if json.Unmarshal(accepted.Payload, &ack) != nil || json.Unmarshal(quality.Payload, &q) != nil || !q.Passed || ack.QualitySHA256 != quality.PayloadSHA256 {
			return "", "", nil, wc.ErrConflict
		}
		work, err := readWarrantyStep(ctx, tx, tenant, c.CaseID, "", "work")
		if err != nil {
			return "", "", nil, err
		}
		var performed wc.WorkReceipt
		if json.Unmarshal(work.Payload, &performed) != nil || q.WorkEvidenceSHA256 != work.PayloadSHA256 {
			return "", "", nil, wc.ErrConflict
		}
		plan, id, hash, err := readWarrantyPlan(ctx, tx, tenant, c.CaseID)
		if err != nil {
			return "", "", nil, err
		}
		if _, err = LookupHumanApproval(ctx, tx, tenant, id, c.OrganizationID, approval.KindWarrantyRepair, hash); err != nil {
			return "", "", nil, err
		}
		if performed.ApprovalID != id || performed.ApprovedSHA256 != hash || len(performed.Parts) != len(plan.Reservations) {
			return "", "", nil, wc.ErrConflict
		}
		entries := []string{}
		for i, part := range performed.Parts {
			if part.LineID != plan.Request.Parts[i].LineID || part.ReservationID != plan.Reservations[i].ID {
				return "", "", nil, wc.ErrConflict
			}
			var valid bool
			err = tx.QueryRow(ctx, `select exists(select 1 from inventory.bulk_inventory_entry e join inventory.bulk_reservation b on b.tenant_id=e.tenant_id and b.reservation_id=$3
 where e.tenant_id=$1 and e.entry_id=$2 and e.organization_id=$4 and e.entry_type='issue' and e.source_kind='warranty-repair' and e.source_id=$5
 and e.item_id=b.item_id and e.bin_id=b.bin_id and e.lot_id is not distinct from b.lot_id and e.quantity=-b.quantity
 and -e.cost_amount=$6::numeric and b.status='consumed' and b.demand_kind='service' and b.demand_id=$5 and b.demand_line_id=$7)`, tenant, part.Issue.EntryID, part.ReservationID, c.OrganizationID, c.CaseID, part.Issue.CostAmount, part.LineID).Scan(&valid)
			if err != nil {
				return "", "", nil, err
			}
			if !valid {
				return "", "", nil, wc.ErrConflict
			}
			entries = append(entries, part.Issue.EntryID)
		}
		var recordedCost string
		err = tx.QueryRow(ctx, `select coalesce(sum(-cost_amount),0)::numeric(38,4)::text from inventory.bulk_inventory_entry where tenant_id=$1 and entry_id=any($2::text[])`, tenant, entries).Scan(&recordedCost)
		if err != nil {
			return "", "", nil, err
		}
		sold, err := claimSoldProfile(ctx, tx, tenant, c)
		if err != nil {
			return "", "", nil, err
		}
		return "reconciled", "closed", map[string]any{"request": r, "settlement": sold.Document().Settlement, "factory_organization_id": c.FactoryOrganizationID, "work_sha256": work.PayloadSHA256, "quality_sha256": quality.PayloadSHA256, "acceptance_sha256": accepted.PayloadSHA256, "approved_sha256": hash, "recorded_inventory_cost": recordedCost, "inventory_entry_ids": entries, "payment_created": false}, nil
	})
}
func (s *Warranty) CancelRepair(ctx context.Context, p identity.Principal, r wc.CancelRepair) (wc.Step, error) {
	if !wc.ValidSHA(r.EvidenceSHA256) || len(strings.TrimSpace(r.Reason)) < 3 || len(r.Reason) > 2048 {
		return wc.Step{}, wc.ErrInvalid
	}
	return s.claimCommand(ctx, p, r.Command, "warranty:cancel", r, func(ctx context.Context, tx pgx.Tx, tenant string, c wc.Claim) (string, string, any, error) {
		if c.State == "closed" || c.State == "cancelled" || c.State == "quality" {
			return "", "", nil, wc.ErrConflict
		}
		if _, err := readWarrantyStep(ctx, tx, tenant, c.CaseID, "", "work"); err == nil {
			return "", "", nil, wc.ErrConflict
		} else if !errors.Is(err, wc.ErrNotFound) {
			return "", "", nil, err
		}
		plan, id, hash, err := readWarrantyPlan(ctx, tx, tenant, c.CaseID)
		if err == nil {
			_, state, err := readHumanApproval(ctx, tx, tenant, id, true)
			if err != nil {
				return "", "", nil, err
			}
			versionDelta := int64(0)
			if state == approval.StatePending {
				if _, err = NewHumanApprovals(s.pool).decideTx(ctx, tx, p, tenant, id, c.OrganizationID, hash, false, r.Reason, "warranty:cancel", nil); err != nil {
					return "", "", nil, err
				}
			} else if state == approval.StateApproved {
				versionDelta = 1
			} else {
				return "", "", nil, wc.ErrConflict
			}
			ids := randomid.Generator{}
			for _, reservation := range plan.Reservations {
				if err = releaseBulkInTx(ctx, tx, tenant, c.OrganizationID, reservation.ID, reservation.Version+versionDelta, ids.New()); err != nil {
					return "", "", nil, err
				}
			}
		} else if !errors.Is(err, wc.ErrNotFound) {
			return "", "", nil, err
		}
		return "cancelled", "cancelled", map[string]any{"request": r, "approval_id": id, "approved_sha256": hash, "inventory_issue_reversed": false}, nil
	})
}
````

### FILE: `internal/platform/postgres/warranty_connected_integration_test.go`

```yaml
block_id: "GO-CONNECTED-WARRANTY-CLAIM:file18:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "a540ac6f2ddf0cfcaff1af20908d856a6cf4edc8bf94868a5bbff429a70cc769"
variables: []
secrets_allowed: false
```

````go
package postgres_test

// AUTHORED fixture composes existing quote/order, official SDK, callback,
// durable reconciliation, checklist and handover with sold warranty terms.
import (
	"context"
	"encoding/json"
	"strings"
	"sync"
	"testing"

	"elite.local/enterprise/internal/platform/identity"
	db "elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/platform/randomid"
	wc "elite.local/enterprise/internal/warrantyclaim"
	"github.com/jackc/pgx/v5/pgxpool"
)

func fixtureWarranty(t *testing.T, pool *pgxpool.Pool, tenant string, days int, leaseSeconds ...int) (*db.Warranty, identity.Principal, wc.Profile) {
	t.Helper()
	document := wc.ProfileDocument{Schema: "elite-warranty-profile/v1", Scope: "MATERIALIZED_PROFILE", Algorithm: "bc-inclusive-fixed-terms", AlgorithmRevision: 1, TenantID: tenant, OrganizationID: "store", FactoryOrganizationID: "factory", PolicyID: "synthetic", TermsVersion: "fixture-v1", TermsText: "Synthetic fixture terms; no jurisdiction or legal compliance claim.", BusinessTimeZone: "America/Argentina/Buenos_Aires", PartsDurationDays: days, LaborDurationDays: days, FaultExclusions: []string{"fixture-exclusion"}, WorkReservationSeconds: 3600, Settlement: "INTERNAL_WARRANTY_SERVICE_ACKNOWLEDGEMENT", AuthorityReference: "synthetic-authority", DecisionReference: "synthetic-acceptance"}
	if len(leaseSeconds) > 1 {
		t.Fatal("one fixture lease maximum")
	}
	if len(leaseSeconds) == 1 {
		document.WorkReservationSeconds = leaseSeconds[0]
	}
	raw, err := json.Marshal(document)
	if err != nil {
		t.Fatal(err)
	}
	p, err := wc.LoadProfile(raw, wc.SHA(raw))
	if err != nil {
		t.Fatal(err)
	}
	store, err := db.NewWarranty(pool, p)
	if err != nil {
		t.Fatal(err)
	}
	actor := identity.Principal{TenantID: tenant, Subject: "operator", Organizations: map[string]struct{}{"store": {}}, Permissions: map[string]struct{}{"warranty:offer": {}, "warranty:read": {}, "warranty:activate": {}}}
	return store, actor, p
}
func TestWarrantySoldTermsThroughConnectedHandover(t *testing.T) {
	ctx := context.Background()
	ids := randomid.Generator{}
	var store *db.Warranty
	var operator, customer identity.Principal
	var sold wc.Profile
	var activated wc.Activation
	r := newConnectedRunWithAllocation(t, false, nil, func(t *testing.T, pool *pgxpool.Pool, tenant string) int64 {
		store, operator, sold = fixtureWarranty(t, pool, tenant, 365)
		customer = identity.Principal{TenantID: tenant, Subject: "customer", Organizations: map[string]struct{}{"store": {}}, Permissions: map[string]struct{}{"warranty:self": {}}}
		request := wc.OfferRequest{QuoteID: "quote", QuoteVersion: 1, ProfileSHA256: sold.Hash()}
		offer, err := store.BindOffer(ctx, operator, request)
		if err != nil || offer.QuoteVersion != 2 {
			t.Fatal("bind", offer, err)
		}
		replay, err := store.BindOffer(ctx, operator, request)
		if err != nil || !replay.Replay || replay.OfferedAt != offer.OfferedAt {
			t.Fatal("offer recovery", replay, err)
		}
		repo := db.NewFranchiseJourney(pool)
		if _, err = repo.AcceptQuote(ctx, tenant, "store", "customer", "quote", 2, strings.Repeat("a", 64), "order", "line", ids.New(), ids.New()); err == nil {
			t.Fatal("accepted without warranty consent")
		}
		var orders int
		if err = pool.QueryRow(ctx, `select count(*) from sales.customer_order where tenant_id=$1`, tenant).Scan(&orders); err != nil || orders != 0 {
			t.Fatal("failed acceptance leaked order", orders, err)
		}
		wrong := customer
		wrong.Subject = "different-customer"
		ack := wc.AcknowledgeRequest{QuoteID: "quote", QuoteVersion: 2, ProfileSHA256: sold.Hash(), EvidenceSHA256: strings.Repeat("d", 64)}
		if _, err = store.Acknowledge(ctx, wrong, ack); err == nil {
			t.Fatal("other customer acknowledged")
		}
		if _, err = store.Offer(ctx, wrong, "quote"); err == nil {
			t.Fatal("other customer read terms")
		}
		wrong = operator
		wrong.TenantID = ids.New()
		if _, err = store.BindOffer(ctx, wrong, request); err == nil {
			t.Fatal("other tenant bound terms")
		}
		wrong = operator
		wrong.Organizations = map[string]struct{}{"other": {}}
		if _, err = store.Offer(ctx, wrong, "quote"); err == nil {
			t.Fatal("other organization read terms")
		}
		bad := ack
		bad.ProfileSHA256 = strings.Repeat("f", 64)
		if _, err = store.Acknowledge(ctx, customer, bad); err == nil {
			t.Fatal("different terms acknowledged")
		}
		bad = ack
		bad.QuoteVersion = 3
		if _, err = store.Acknowledge(ctx, customer, bad); err == nil {
			t.Fatal("stale version acknowledged")
		}
		if _, err = pool.Exec(ctx, `update sales.quotation set total_minor_units=1 where tenant_id=$1 and quotation_id='quote'`, tenant); err == nil {
			t.Fatal("offered price mutated")
		}
		// Two uncertain duplicate submissions serialize to one immutable consent.
		var wg sync.WaitGroup
		results := make(chan error, 2)
		for i := 0; i < 2; i++ {
			wg.Add(1)
			go func() { defer wg.Done(); _, err := store.Acknowledge(ctx, customer, ack); results <- err }()
		}
		wg.Wait()
		close(results)
		for err := range results {
			if err != nil {
				t.Fatal("concurrent consent", err)
			}
		}
		offer, err = store.Offer(ctx, customer, "quote")
		if err != nil || !offer.Acknowledged || offer.EvidenceSHA256 != ack.EvidenceSHA256 {
			t.Fatal("consent recovery", offer, err)
		}
		bad = ack
		bad.EvidenceSHA256 = strings.Repeat("e", 64)
		if _, err = store.Acknowledge(ctx, customer, bad); err == nil {
			t.Fatal("changed replay accepted")
		}
		return offer.QuoteVersion
	})
	if code := r.callback(t, "evt_warranty", "checkout.session.completed", "cs_test_fixture", true); code != 200 {
		t.Fatal(code)
	}
	if n, err := r.processor.ProcessOnce(ctx); err != nil || n != 1 {
		t.Fatal(n, err)
	}
	var hash string
	if err := r.pool.QueryRow(ctx, `select evidence_sha256_hex from payment.provider_observation where tenant_id=$1`, r.tenant).Scan(&hash); err != nil {
		t.Fatal(err)
	}
	assertConnectedCommercialReleaseWithHandoverHook(t, r, hash, func(t *testing.T, r *connectedRun, handover string) {
		// New configuration must not replace the profile the customer actually bought.
		changed, _, current := fixtureWarranty(t, r.pool, r.tenant, 7)
		if current.Hash() == sold.Hash() {
			t.Fatal("fixture requires distinct policy revisions")
		}
		// Fail after warranty/activation writes but before transaction commit.
		_, err := r.pool.Exec(ctx, `create function service_ops.fixture_warranty_outbox_failure() returns trigger language plpgsql as $$ begin if new.event_type='warranty.activated' then raise exception 'fixture commit failure'; end if;return new;end $$;
 create trigger fixture_warranty_outbox_failure before insert on platform.outbox_event for each row execute function service_ops.fixture_warranty_outbox_failure()`)
		if err != nil {
			t.Fatal(err)
		}
		if _, err = changed.Activate(ctx, operator, handover); err == nil {
			t.Fatal("fault injection did not fail")
		}
		var count int
		if err = r.pool.QueryRow(ctx, `select count(*) from service_ops.warranty where tenant_id=$1`, r.tenant).Scan(&count); err != nil || count != 0 {
			t.Fatal("activation rollback leaked warranty", count, err)
		}
		if _, err = r.pool.Exec(ctx, `drop trigger fixture_warranty_outbox_failure on platform.outbox_event;drop function service_ops.fixture_warranty_outbox_failure()`); err != nil {
			t.Fatal(err)
		}
		activated, err = changed.Activate(ctx, operator, handover)
		if err != nil {
			t.Fatal("activate sold terms", err)
		}
		dates, err := sold.DatesAt(activated.AcceptedAt)
		if err != nil || activated.ProfileSHA256 != sold.Hash() || activated.Dates != dates {
			t.Fatal("current configuration replaced sold terms", activated, dates, err)
		}
		again, err := changed.Activate(ctx, operator, handover)
		if err != nil || !again.Replay || again.WarrantyID != activated.WarrantyID || !again.ActivatedAt.Equal(activated.ActivatedAt) {
			t.Fatal("activation recovery", again, err)
		}
		read, err := changed.Activation(ctx, customer, handover)
		if err != nil || read.WarrantyID != activated.WarrantyID {
			t.Fatal("customer receipt", read, err)
		}
		wrong := operator
		wrong.Subject = "another-operator"
		if _, err = changed.Activate(ctx, wrong, handover); err == nil {
			t.Fatal("different actor replay")
		}
		if _, err = r.pool.Exec(ctx, `update service_ops.warranty_activation set parts_end=parts_end+1 where tenant_id=$1`, r.tenant); err == nil {
			t.Fatal("activation rewritten")
		}
		wire := newWarrantyHTTPFixture(t, changed, true)
		assertWarrantyClaimLifecycle(t, r, wire, operator, customer, activated)
		wire.mu.Lock()
		drops, recovered := len(wire.drops), wire.recoveries
		workPosts := wire.posts["/v1/franchise/warranty/claims/claim-main/work"]
		reconciliationPosts := wire.posts["/v1/factory/warranty/claims/claim-main/reconciliation"]
		wire.mu.Unlock()
		if drops != 2 || recovered < 2 || workPosts != 3 || reconciliationPosts != 3 {
			t.Fatal("HTTP loss recovery accounting", drops, recovered, workPosts, reconciliationPosts)
		}
		t.Log("WARRANTY_HTTP_PG_PASS claims=4 dropped_commit_responses=2 recovery_gets_verified=true hidden_post_retry=false denied_factory_post=1 verifier_fixture_only=true")
	})
	var warranties, activations, events int
	var state string
	err := r.pool.QueryRow(ctx, `select (select count(*) from service_ops.warranty where tenant_id=$1),(select count(*) from service_ops.warranty_activation where tenant_id=$1),(select count(*) from platform.outbox_event where tenant_id=$1 and event_type='warranty.activated'),(select status from service_ops.warranty where tenant_id=$1)`, r.tenant).Scan(&warranties, &activations, &events, &state)
	if err != nil || warranties != 1 || activations != 1 || events != 1 || state != "active" {
		t.Fatal("durable warranty or invented refund voiding", warranties, activations, events, state, err)
	}
	t.Log("WARRANTY_CONNECTED_PASS quote_terms_customer_consent_order_official_sdk_callback_reconcile_handover_commercial_receipt_activation=true sold_profile_preserved=true rollback=true scope=true exact_replay=true claim_workflow_connected=true live_proven=false")
}
````

### FILE: `internal/platform/postgres/warranty_expiry_integration_test.go`

```yaml
block_id: "GO-CONNECTED-WARRANTY-CLAIM:file19:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "468b44425a4d793f0dfc9d9b79e2317d99c337572bc9373049f335391cb02478"
variables: []
secrets_allowed: false
```

````go
package postgres_test

// Narrow commit-time lease regression. Accepted handover and attended service
// are explicit historical fixtures here; the separate HTTP J1/J4 test proves
// their full creation journey. The approval/reservation transaction is real.
import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"elite.local/enterprise/internal/inventorycontrol"
	"elite.local/enterprise/internal/platform/identity"
	db "elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/platform/randomid"
	wc "elite.local/enterprise/internal/warrantyclaim"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestWarrantyPendingLeaseExpiresAtCommit(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 85*time.Second)
	defer cancel()
	pool := connectedPool(t)
	ids := randomid.Generator{}
	var store *db.Warranty
	var actor identity.Principal
	tenant := connectedSeedOrder(t, pool, "stripe", func(t *testing.T, pool *pgxpool.Pool, tenant string) int64 {
		var profile wc.Profile
		store, actor, profile = fixtureWarranty(t, pool, tenant, 365, 60)
		offered, err := store.BindOffer(ctx, actor, wc.OfferRequest{QuoteID: "quote", QuoteVersion: 1, ProfileSHA256: profile.Hash()})
		if err != nil {
			t.Fatal(err)
		}
		customer := identity.Principal{TenantID: tenant, Subject: "customer", Organizations: map[string]struct{}{"store": {}}, Permissions: map[string]struct{}{"warranty:self": {}}}
		if _, err = store.Acknowledge(ctx, customer, wc.AcknowledgeRequest{QuoteID: "quote", QuoteVersion: offered.QuoteVersion, ProfileSHA256: profile.Hash(), EvidenceSHA256: strings.Repeat("a", 64)}); err != nil {
			t.Fatal(err)
		}
		return offered.QuoteVersion
	})
	must := func(sql string, args ...any) {
		t.Helper()
		if _, err := pool.Exec(ctx, sql, args...); err != nil {
			t.Fatal(err)
		}
	}
	must(`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'factory','factory','Fixture factory','factory')`, tenant)
	must(`insert into sales.delivery_handover(tenant_id,handover_id,organization_id,order_id,customer_principal_id,stock_unit_id,state,acceptance_evidence_sha256_hex,customer_accepted_at,version)
 values($1,'historical-handover','store','order','customer','stock','accepted',$2,clock_timestamp(),3)`, tenant, strings.Repeat("b", 64))
	if _, err := store.Activate(ctx, actor, "historical-handover"); err != nil {
		t.Fatal(err)
	}
	must(`insert into crm.appointment(tenant_id,appointment_id,organization_id,lead_id,customer_principal_id,model_id,appointment_kind,starts_at,state,version,created_at,updated_at)
 values($1,'historical-service','store','lead','customer','model','service',clock_timestamp()-interval '30 minutes','completed',3,clock_timestamp()-interval '1 hour',clock_timestamp())`, tenant)
	for _, p := range []string{"warranty:request", "warranty:diagnose", "warranty:plan", "warranty:approve"} {
		actor.Permissions[p] = struct{}{}
	}
	review := actor
	review.Subject = "reviewer"
	inv := db.NewInventoryControl(pool)
	if _, err := inv.CreateBulkItem(ctx, tenant, ids.New(), inventorycontrol.BulkItem{ID: "part", Code: "PART", Description: "Fixture part", BaseUOM: "EA", BaseRoundingPrecision: "1", TrackingMode: "none", CostingMethod: "fifo", Version: 1}); err != nil {
		t.Fatal(err)
	}
	if _, err := inv.CreateWarehouseBin(ctx, tenant, ids.New(), inventorycontrol.WarehouseBin{ID: "bin", OrganizationID: "store", Code: "BIN", Type: "pick", Version: 1}); err != nil {
		t.Fatal(err)
	}
	if _, err := inv.ConfigureItemBin(ctx, tenant, ids.New(), inventorycontrol.ItemBinPolicy{OrganizationID: "store", ItemID: "part", BinID: "bin", MinQuantity: "0", Version: 1}); err != nil {
		t.Fatal(err)
	}
	if _, err := inv.ReceiveBulk(ctx, tenant, ids.New(), ids.New(), ids.New(), ids.New(), inventorycontrol.BulkReceipt{OrganizationID: "store", BinID: "bin", ItemID: "part", Quantity: "1", UnitCost: "10", PostingDate: time.Now().UTC(), SourceKind: "fixture", SourceID: "expiry-stock"}); err != nil {
		t.Fatal(err)
	}
	if _, err := store.OpenClaim(ctx, actor, wc.OpenClaim{CaseID: "expiry-case", HandoverID: "historical-handover", AppointmentID: "historical-service", Severity: "low", Description: "Fixture lease inspection", EvidenceSHA256: strings.Repeat("c", 64)}); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Diagnose(ctx, actor, wc.Diagnose{Command: wc.Command{CaseID: "expiry-case", CommandID: "diagnosis", ExpectedVersion: 1}, FaultCode: "fixture-wear", Description: "Fixture covered fault", EvidenceSHA256: strings.Repeat("d", 64)}); err != nil {
		t.Fatal(err)
	}
	plan, err := store.PlanRepair(ctx, actor, wc.Plan{Command: wc.Command{CaseID: "expiry-case", CommandID: "plan", ExpectedVersion: 2}, Parts: []wc.Part{{LineID: "line", ItemID: "part", BinID: "bin", Quantity: "1"}}, LaborWork: ""})
	if err != nil {
		t.Fatal(err)
	}
	var bound struct {
		ID  string `json:"approval_id"`
		SHA string `json:"payload_sha256"`
	}
	if json.Unmarshal(plan.Payload, &bound) != nil {
		t.Fatal("plan payload")
	}
	must(`create function service_ops.fixture_delay_warranty_commit() returns trigger language plpgsql as $$declare expiry timestamptz;begin
 if new.event_type='warranty.claim-decision' then
 select expires_at into expiry from service_ops.warranty_work_plan where tenant_id=new.tenant_id and service_case_id=new.aggregate_id;
 perform pg_sleep(greatest(0,extract(epoch from(expiry-clock_timestamp())))+0.1);
 end if;return new;end $$;
 create trigger fixture_delay_warranty_commit before insert on platform.outbox_event for each row execute function service_ops.fixture_delay_warranty_commit()`)
	t.Log("WARRANTY_EXPIRY_PROBE waiting only for the explicit60second pending lease to cross COMMIT; no shortened production policy")
	_, err = store.DecideRepair(ctx, review, wc.DecideRepair{Command: wc.Command{CaseID: "expiry-case", CommandID: "approve", ExpectedVersion: 3}, PayloadSHA256: bound.SHA, Approved: true, Reason: "Fixture approval before lease expiry"})
	if err == nil || ctx.Err() != nil || !strings.Contains(err.Error(), "warranty pending reservation expired before commit") {
		t.Fatal("commit did not reject exact lease expiry", err, ctx.Err())
	}
	must(`drop trigger fixture_delay_warranty_commit on platform.outbox_event;drop function service_ops.fixture_delay_warranty_commit()`)
	var decisions, steps, events, issues int
	var state, reservationState string
	var version int64
	var held string
	err = pool.QueryRow(ctx, `select (select count(*) from approval.decision where tenant_id=$1),(select count(*) from service_ops.warranty_claim_step where tenant_id=$1 and kind='decision'),(select count(*) from platform.outbox_event where tenant_id=$1 and event_type='warranty.claim-decision'),(select count(*) from inventory.bulk_inventory_entry where tenant_id=$1 and entry_type='issue'),(select state from approval.request where tenant_id=$1 and request_id=$2),(select status from inventory.bulk_reservation where tenant_id=$1),(select version from inventory.bulk_reservation where tenant_id=$1),(select reserved_quantity::text from inventory.bulk_balance where tenant_id=$1)`, tenant, bound.ID).Scan(&decisions, &steps, &events, &issues, &state, &reservationState, &version, &held)
	if err != nil || decisions != 0 || steps != 0 || events != 0 || issues != 0 || state != "pending" || reservationState != "reservation" || version != 1 || held != "1.000000" {
		t.Fatal("expired approval leaked effect", decisions, steps, events, issues, state, reservationState, version, held, err)
	}
	rejected, err := store.DecideRepair(ctx, review, wc.DecideRepair{Command: wc.Command{CaseID: "expiry-case", CommandID: "reject-expired", ExpectedVersion: 3}, PayloadSHA256: bound.SHA, Approved: false, Reason: "Expired pending lease released"})
	if err != nil || rejected.State != "cancelled" {
		t.Fatal("expired hold recovery", rejected, err)
	}
	if err = pool.QueryRow(ctx, `select reserved_quantity::text from inventory.bulk_balance where tenant_id=$1`, tenant).Scan(&held); err != nil || held != "0.000000" {
		t.Fatal("hold not released", held, err)
	}
	t.Log("WARRANTY_EXPIRY_COMMIT_PASS approved_before_wait_commit_after_expiry=true leaked_decisions_steps_events_issues=0 original_pending_reservation_preserved=true explicit_rejection_releases_hold=true")
}
````

### FILE: `internal/platform/postgres/warranty_http_integration_test.go`

```yaml
block_id: "GO-CONNECTED-WARRANTY-CLAIM:file20:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "c4b93c70df57023e17c7233bc84f869f84fee49b46c0b090a70b4a72a2487ffe"
variables: []
secrets_allowed: false
```

````go
package postgres_test

// AUTHORED transport fixture. The verifier returns explicitly assigned test
// principals; JWT cryptography is not claimed by this fixture. Real HTTP,
// handlers, repository authorization, durable commands and PG effects execute.
import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/randomid"
	wc "elite.local/enterprise/internal/warrantyclaim"
)

type warrantyClaimCommands interface {
	OpenClaim(context.Context, identity.Principal, wc.OpenClaim) (wc.Step, error)
	Diagnose(context.Context, identity.Principal, wc.Diagnose) (wc.Step, error)
	PlanRepair(context.Context, identity.Principal, wc.Plan) (wc.Step, error)
	DecideRepair(context.Context, identity.Principal, wc.DecideRepair) (wc.Step, error)
	CompleteRepairWork(context.Context, identity.Principal, wc.CompleteWork) (wc.Step, error)
	RecordRepairQuality(context.Context, identity.Principal, wc.Quality) (wc.Step, error)
	AcceptRepair(context.Context, identity.Principal, wc.AcceptRepair) (wc.Step, error)
	ReconcileRepair(context.Context, identity.Principal, wc.ReconcileRepair) (wc.Step, error)
	CancelRepair(context.Context, identity.Principal, wc.CancelRepair) (wc.Step, error)
}
type warrantyHTTPFixture struct {
	mu         sync.Mutex
	server     *httptest.Server
	principals map[string]identity.Principal
	drops      map[string]int
	recoveries int
	posts      map[string]int
}

func (h *warrantyHTTPFixture) Verify(_ context.Context, token string) (identity.Principal, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	p, ok := h.principals[token]
	if !ok {
		return p, identity.ErrUnauthenticated
	}
	return p, nil
}
func newWarrantyHTTPFixture(t *testing.T, service httpapi.WarrantyService, drop bool) *warrantyHTTPFixture {
	t.Helper()
	h := &warrantyHTTPFixture{principals: map[string]identity.Principal{}, drops: map[string]int{}, posts: map[string]int{}}
	mux := http.NewServeMux()
	httpapi.WarrantyModule{Service: service}.Register(mux, h)
	h.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, r)
		shouldDrop := false
		if r.Method == "POST" {
			h.mu.Lock()
			h.posts[r.URL.Path]++
			if drop && recorder.Code >= 200 && recorder.Code < 300 && (strings.HasSuffix(r.URL.Path, "/work") || strings.HasSuffix(r.URL.Path, "/reconciliation")) && h.drops[r.URL.Path] == 0 {
				h.drops[r.URL.Path]++
				shouldDrop = true
			}
			h.mu.Unlock()
		}
		if shouldDrop {
			conn, _, err := w.(http.Hijacker).Hijack()
			if err != nil {
				t.Error(err)
				return
			}
			conn.Close()
			return
		}
		for key, values := range recorder.Header() {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
		w.WriteHeader(recorder.Code)
		_, _ = w.Write(recorder.Body.Bytes())
	}))
	h.server.Client().Timeout = 10 * time.Second
	t.Cleanup(h.server.Close)
	return h
}
func (h *warrantyHTTPFixture) send(ctx context.Context, p identity.Principal, method, path string, body any, out any) error {
	raw, err := json.Marshal(body)
	if err != nil {
		return err
	}
	id := randomid.Generator{}.New()
	copy := identity.Principal{TenantID: p.TenantID, Subject: p.Subject, Permissions: map[string]struct{}{}, Organizations: map[string]struct{}{}}
	for k := range p.Permissions {
		copy.Permissions[k] = struct{}{}
	}
	for k := range p.Organizations {
		copy.Organizations[k] = struct{}{}
	}
	h.mu.Lock()
	h.principals[id] = copy
	h.mu.Unlock()
	req, err := http.NewRequestWithContext(ctx, method, h.server.URL+path, bytes.NewReader(raw))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+id)
	req.Header.Set("Content-Type", "application/json")
	response, err := h.server.Client().Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	data, err := io.ReadAll(io.LimitReader(response.Body, 131073))
	if err != nil {
		return err
	}
	if len(data) > 131072 {
		return fmt.Errorf("fixture response exceeded bound")
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return fmt.Errorf("HTTP%d: %s", response.StatusCode, data)
	}
	return json.Unmarshal(data, out)
}
func (h *warrantyHTTPFixture) step(ctx context.Context, p identity.Principal, c wc.Command, surface, action string, body any) (wc.Step, error) {
	org := "store"
	if surface == "factory" {
		org = "factory"
	}
	path := "/v1/" + surface + "/warranty/claims/" + url.PathEscape(c.CaseID) + "/" + action + "?organization_id=" + org
	var out wc.Step
	err := h.send(ctx, p, "POST", path, body, &out)
	if err != nil {
		var recovered wc.Step
		lookup := "/v1/" + surface + "/warranty/claims/" + url.PathEscape(c.CaseID) + "?organization_id=" + org + "&command_id=" + url.QueryEscape(c.CommandID)
		if recovery := h.send(ctx, p, "GET", lookup, nil, &recovered); recovery == nil {
			_, expected, e := wc.Canonical(body)
			if e != nil || recovered.RequestSHA256 != expected || recovered.Actor != p.Subject {
				return wc.Step{}, fmt.Errorf("recovery did not bind original request")
			}
			h.mu.Lock()
			h.recoveries++
			h.mu.Unlock()
			return recovered, nil
		}
	}
	return out, err
}
func (h *warrantyHTTPFixture) OpenClaim(ctx context.Context, p identity.Principal, r wc.OpenClaim) (wc.Step, error) {
	surface := "franchise"
	if p.Subject == "customer" {
		surface = "customer"
	}
	var out wc.Step
	err := h.send(ctx, p, "POST", "/v1/"+surface+"/warranty/claims?organization_id=store", r, &out)
	return out, err
}

func (h *warrantyHTTPFixture) Diagnose(ctx context.Context, p identity.Principal, r wc.Diagnose) (wc.Step, error) {
	return h.step(ctx, p, r.Command, "franchise", "diagnosis", r)
}

func (h *warrantyHTTPFixture) PlanRepair(ctx context.Context, p identity.Principal, r wc.Plan) (wc.Step, error) {
	return h.step(ctx, p, r.Command, "franchise", "plan", r)
}

func (h *warrantyHTTPFixture) DecideRepair(ctx context.Context, p identity.Principal, r wc.DecideRepair) (wc.Step, error) {
	return h.step(ctx, p, r.Command, "franchise", "decision", r)
}

func (h *warrantyHTTPFixture) CompleteRepairWork(ctx context.Context, p identity.Principal, r wc.CompleteWork) (wc.Step, error) {
	return h.step(ctx, p, r.Command, "franchise", "work", r)
}

func (h *warrantyHTTPFixture) RecordRepairQuality(ctx context.Context, p identity.Principal, r wc.Quality) (wc.Step, error) {
	return h.step(ctx, p, r.Command, "franchise", "quality", r)
}

func (h *warrantyHTTPFixture) AcceptRepair(ctx context.Context, p identity.Principal, r wc.AcceptRepair) (wc.Step, error) {
	return h.step(ctx, p, r.Command, "customer", "acceptance", r)
}

func (h *warrantyHTTPFixture) ReconcileRepair(ctx context.Context, p identity.Principal, r wc.ReconcileRepair) (wc.Step, error) {
	return h.step(ctx, p, r.Command, "factory", "reconciliation", r)
}

func (h *warrantyHTTPFixture) CancelRepair(ctx context.Context, p identity.Principal, r wc.CancelRepair) (wc.Step, error) {
	return h.step(ctx, p, r.Command, "franchise", "cancellation", r)
}
````

### FILE: `internal/platform/postgres/warranty_terms.go`

```yaml
block_id: "GO-CONNECTED-WARRANTY-CLAIM:file21:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "238efbd6fc8fc58820909c50f7edc02aa7a81528ad33756202f87c43778b609a"
variables: []
secrets_allowed: false
```

````go
// AUTHORED atomic consent, historical scope and outbox composition. Coverage
// predicates come from the separately declared BC adaptation, not this glue.
package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"elite.local/enterprise/internal/platform/identity"
	wc "elite.local/enterprise/internal/warrantyclaim"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Warranty struct {
	pool    *pgxpool.Pool
	profile wc.Profile
}

func (s *Warranty) Scope() (string, string) {
	if s == nil {
		return "", ""
	}
	return s.profile.Scope()
}
func (s *Warranty) ProfileDocument(ctx context.Context, p identity.Principal) (json.RawMessage, error) {
	if !s.allowed(p, "warranty:read") {
		return nil, wc.ErrNotFound
	}
	return s.profile.Bytes(), nil
}

func NewWarranty(pool *pgxpool.Pool, profile wc.Profile) (*Warranty, error) {
	if pool == nil || !profile.Valid() {
		return nil, wc.ErrInvalid
	}
	return &Warranty{pool: pool, profile: profile}, nil
}
func (s *Warranty) ProfileSHA256() string {
	if s == nil {
		return ""
	}
	return s.profile.Hash()
}
func (s *Warranty) allowed(p identity.Principal, permission string) bool {
	if s == nil || s.pool == nil || !s.profile.Valid() {
		return false
	}
	tenant, org := s.profile.Scope()
	return approvalPrincipal(p, tenant, org, permission)
}
func warrantyError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return wc.ErrNotFound
	}
	if postgresConflict(err) {
		return wc.ErrConflict
	}
	return err
}
func warrantyEvent(ctx context.Context, tx pgx.Tx, tenant, id, event string, version int64, payload any) error {
	raw, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)
 values($1,gen_random_uuid(),'warranty',$2,$3,$4,1,clock_timestamp(),$5)`, tenant, id, version, event, raw)
	return err
}
func readWarrantyOffer(ctx context.Context, q interface {
	QueryRow(context.Context, string, ...any) pgx.Row
}, tenant, org, id string) (wc.Offer, error) {
	var out wc.Offer
	err := q.QueryRow(ctx, `select o.quotation_id,o.quotation_version,o.customer_principal_id,o.profile_sha256,o.profile_bytes,o.offered_by,o.offered_at,
 a.quotation_id is not null,coalesce(a.evidence_sha256,'') from service_ops.warranty_offer o
 left join service_ops.warranty_acknowledgement a on a.tenant_id=o.tenant_id and a.quotation_id=o.quotation_id
 where o.tenant_id=$1 and o.organization_id=$2 and o.quotation_id=$3`, tenant, org, id).Scan(&out.QuoteID, &out.QuoteVersion, &out.CustomerSubject, &out.ProfileSHA256, &out.Profile, &out.OfferedBy, &out.OfferedAt, &out.Acknowledged, &out.EvidenceSHA256)
	if err != nil {
		return out, warrantyError(err)
	}
	profile, err := wc.LoadProfile(out.Profile, out.ProfileSHA256)
	if err != nil {
		return wc.Offer{}, err
	}
	pt, po := profile.Scope()
	if pt != tenant || po != org {
		return wc.Offer{}, wc.ErrConflict
	}
	return out, nil
}
func (s *Warranty) Offer(ctx context.Context, p identity.Principal, id string) (wc.Offer, error) {
	if !wc.ValidID(id) || !(s.allowed(p, "warranty:read") || s.allowed(p, "warranty:self")) {
		return wc.Offer{}, wc.ErrInvalid
	}
	tenant, org := s.profile.Scope()
	out, err := readWarrantyOffer(ctx, s.pool, tenant, org, id)
	if err == nil && !s.allowed(p, "warranty:read") && out.CustomerSubject != p.Subject {
		return wc.Offer{}, wc.ErrNotFound
	}
	return out, err
}
func (s *Warranty) BindOffer(ctx context.Context, p identity.Principal, r wc.OfferRequest) (wc.Offer, error) {
	if !s.allowed(p, "warranty:offer") || !wc.ValidID(r.QuoteID) || r.QuoteVersion < 1 || r.ProfileSHA256 != s.profile.Hash() {
		return wc.Offer{}, wc.ErrInvalid
	}
	tenant, org := s.profile.Scope()
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return wc.Offer{}, err
	}
	defer tx.Rollback(ctx)
	var customer, state string
	var version int64
	var current bool
	err = tx.QueryRow(ctx, `select coalesce(customer_principal_id,''),state,version,valid_until>clock_timestamp() from sales.quotation
 where tenant_id=$1 and organization_id=$2 and quotation_id=$3 for update`, tenant, org, r.QuoteID).Scan(&customer, &state, &version, &current)
	if err != nil {
		return wc.Offer{}, warrantyError(err)
	}
	existing, err := readWarrantyOffer(ctx, tx, tenant, org, r.QuoteID)
	if err == nil {
		if existing.OfferedBy != p.Subject || existing.ProfileSHA256 != r.ProfileSHA256 || existing.QuoteVersion-1 != r.QuoteVersion {
			return wc.Offer{}, wc.ErrConflict
		}
		existing.Replay = true
		return existing, tx.Commit(ctx)
	}
	if !errors.Is(err, wc.ErrNotFound) {
		return wc.Offer{}, err
	}
	if customer == "" || state != "issued" || version != r.QuoteVersion || !current {
		return wc.Offer{}, wc.ErrConflict
	}
	if _, err = tx.Exec(ctx, `update sales.quotation set version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and quotation_id=$2`, tenant, r.QuoteID); err != nil {
		return wc.Offer{}, err
	}
	_, err = tx.Exec(ctx, `insert into service_ops.warranty_offer(tenant_id,quotation_id,organization_id,customer_principal_id,quotation_version,profile_sha256,profile_bytes,offered_by)
 values($1,$2,$3,$4,$5,$6,$7,$8)`, tenant, r.QuoteID, org, customer, version+1, r.ProfileSHA256, s.profile.Bytes(), p.Subject)
	if err != nil {
		return wc.Offer{}, warrantyError(err)
	}
	out, err := readWarrantyOffer(ctx, tx, tenant, org, r.QuoteID)
	if err != nil {
		return wc.Offer{}, err
	}
	if err = warrantyEvent(ctx, tx, tenant, r.QuoteID, "warranty.terms-offered", version+1, map[string]any{"organization_id": org, "quote_id": r.QuoteID, "quote_version": out.QuoteVersion, "profile_sha256": r.ProfileSHA256, "actor": p.Subject}); err != nil {
		return wc.Offer{}, err
	}
	return out, tx.Commit(ctx)
}
func (s *Warranty) Acknowledge(ctx context.Context, p identity.Principal, r wc.AcknowledgeRequest) (wc.Offer, error) {
	if !s.allowed(p, "warranty:self") || !wc.ValidID(r.QuoteID) || r.QuoteVersion < 2 || !wc.ValidSHA(r.ProfileSHA256) || !wc.ValidSHA(r.EvidenceSHA256) {
		return wc.Offer{}, wc.ErrInvalid
	}
	tenant, org := s.profile.Scope()
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return wc.Offer{}, err
	}
	defer tx.Rollback(ctx)
	var customer, state string
	var version int64
	var current bool
	err = tx.QueryRow(ctx, `select coalesce(customer_principal_id,''),state,version,valid_until>clock_timestamp() from sales.quotation
 where tenant_id=$1 and organization_id=$2 and quotation_id=$3 and customer_principal_id=$4 for update`, tenant, org, r.QuoteID, p.Subject).Scan(&customer, &state, &version, &current)
	if err != nil {
		return wc.Offer{}, warrantyError(err)
	}
	offer, err := readWarrantyOffer(ctx, tx, tenant, org, r.QuoteID)
	if err != nil {
		return wc.Offer{}, err
	}
	if offer.CustomerSubject != p.Subject || offer.ProfileSHA256 != r.ProfileSHA256 || offer.QuoteVersion != r.QuoteVersion {
		return wc.Offer{}, wc.ErrConflict
	}
	if offer.Acknowledged {
		if offer.EvidenceSHA256 != r.EvidenceSHA256 {
			return wc.Offer{}, wc.ErrConflict
		}
		offer.Replay = true
		return offer, tx.Commit(ctx)
	}
	if state != "issued" || !current || version != r.QuoteVersion {
		return wc.Offer{}, wc.ErrConflict
	}
	_, err = tx.Exec(ctx, `insert into service_ops.warranty_acknowledgement(tenant_id,quotation_id,quotation_version,customer_principal_id,profile_sha256,evidence_sha256)
 values($1,$2,$3,$4,$5,$6)`, tenant, r.QuoteID, r.QuoteVersion, p.Subject, r.ProfileSHA256, r.EvidenceSHA256)
	if err != nil {
		return wc.Offer{}, warrantyError(err)
	}
	offer.Acknowledged = true
	offer.EvidenceSHA256 = r.EvidenceSHA256
	if err = warrantyEvent(ctx, tx, tenant, r.QuoteID, "warranty.terms-acknowledged", version, map[string]any{"organization_id": org, "quote_id": r.QuoteID, "quote_version": version, "profile_sha256": r.ProfileSHA256, "evidence_sha256": r.EvidenceSHA256, "actor": p.Subject}); err != nil {
		return wc.Offer{}, err
	}
	return offer, tx.Commit(ctx)
}
func readWarrantyActivation(ctx context.Context, q interface {
	QueryRow(context.Context, string, ...any) pgx.Row
}, tenant, org, handover string) (wc.Activation, error) {
	var a wc.Activation
	err := q.QueryRow(ctx, `select warranty_id,handover_id,quotation_id,order_id,stock_unit_id,organization_id,customer_principal_id,profile_sha256,terms_version,
 parts_start::text,parts_end::text,labor_start::text,labor_end::text,handover_accepted_at,activated_by,activated_at from service_ops.warranty_activation
 where tenant_id=$1 and organization_id=$2 and handover_id=$3`, tenant, org, handover).Scan(&a.WarrantyID, &a.HandoverID, &a.QuoteID, &a.OrderID, &a.StockUnitID, &a.OrganizationID, &a.CustomerSubject, &a.ProfileSHA256, &a.TermsVersion, &a.Dates.PartsStart, &a.Dates.PartsEnd, &a.Dates.LaborStart, &a.Dates.LaborEnd, &a.AcceptedAt, &a.ActivatedBy, &a.ActivatedAt)
	return a, warrantyError(err)
}
func (s *Warranty) Activation(ctx context.Context, p identity.Principal, handover string) (wc.Activation, error) {
	if !wc.ValidID(handover) || !(s.allowed(p, "warranty:read") || s.allowed(p, "warranty:self")) {
		return wc.Activation{}, wc.ErrInvalid
	}
	tenant, org := s.profile.Scope()
	out, err := readWarrantyActivation(ctx, s.pool, tenant, org, handover)
	if err == nil && !s.allowed(p, "warranty:read") && out.CustomerSubject != p.Subject {
		return wc.Activation{}, wc.ErrNotFound
	}
	return out, err
}
func (s *Warranty) Activate(ctx context.Context, p identity.Principal, handover string) (wc.Activation, error) {
	var out wc.Activation
	if !wc.ValidID(handover) || !s.allowed(p, "warranty:activate") {
		return out, wc.ErrInvalid
	}
	tenant, org := s.profile.Scope()
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return out, err
	}
	defer tx.Rollback(ctx)
	var state string
	err = tx.QueryRow(ctx, `select handover_id,order_id,stock_unit_id,customer_principal_id,state,customer_accepted_at from sales.delivery_handover
 where tenant_id=$1 and organization_id=$2 and handover_id=$3 and state='accepted' for update`, tenant, org, handover).Scan(&out.HandoverID, &out.OrderID, &out.StockUnitID, &out.CustomerSubject, &state, &out.AcceptedAt)
	if err != nil {
		return out, warrantyError(err)
	}
	existing, err := readWarrantyActivation(ctx, tx, tenant, org, handover)
	if err == nil {
		if existing.ActivatedBy != p.Subject {
			return out, wc.ErrConflict
		}
		existing.Replay = true
		return existing, tx.Commit(ctx)
	}
	if !errors.Is(err, wc.ErrNotFound) {
		return out, err
	}
	err = tx.QueryRow(ctx, `select qa.quotation_id from sales.quotation_acceptance qa where qa.tenant_id=$1 and qa.order_id=$2 and qa.customer_principal_id=$3`, tenant, out.OrderID, out.CustomerSubject).Scan(&out.QuoteID)
	if err != nil {
		return out, warrantyError(err)
	}
	offer, err := readWarrantyOffer(ctx, tx, tenant, org, out.QuoteID)
	if err != nil {
		return out, err
	}
	if !offer.Acknowledged || offer.CustomerSubject != out.CustomerSubject {
		return out, wc.ErrConflict
	}
	sold, err := wc.LoadProfile(offer.Profile, offer.ProfileSHA256)
	if err != nil {
		return out, err
	}
	dates, err := sold.DatesAt(out.AcceptedAt)
	if err != nil {
		return out, err
	}
	// service_ops.warranty uses instants; the declared date intervals are inclusive.
	// Store the following local midnight as the exclusive outer bound.
	doc := sold.Document()
	zone, err := time.LoadLocation(doc.BusinessTimeZone)
	if err != nil {
		return out, err
	}
	start, err := time.ParseInLocation(time.DateOnly, dates.PartsStart, zone)
	if err != nil {
		return out, err
	}
	last := dates.PartsEnd
	if dates.LaborEnd > last {
		last = dates.LaborEnd
	}
	end, err := time.ParseInLocation(time.DateOnly, last, zone)
	if err != nil {
		return out, err
	}
	end = end.AddDate(0, 0, 1)
	if end.Year() > 9999 {
		return out, wc.ErrInvalid
	}
	out.WarrantyID = wc.StableID("warranty", tenant, handover)
	_, err = tx.Exec(ctx, `insert into service_ops.warranty(tenant_id,warranty_id,stock_unit_id,customer_principal_id,starts_at,ends_at,terms_version,status)
 values($1,$2,$3,$4,$5,$6,$7,'active')`, tenant, out.WarrantyID, out.StockUnitID, out.CustomerSubject, start, end, doc.TermsVersion)
	if err != nil {
		return out, warrantyError(err)
	}
	_, err = tx.Exec(ctx, `insert into service_ops.warranty_activation(tenant_id,warranty_id,handover_id,quotation_id,order_id,organization_id,customer_principal_id,stock_unit_id,profile_sha256,terms_version,
 parts_start,parts_end,labor_start,labor_end,handover_accepted_at,activated_by)
 values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11::date,$12::date,$13::date,$14::date,$15,$16)`, tenant, out.WarrantyID, handover, out.QuoteID, out.OrderID, org, out.CustomerSubject, out.StockUnitID, sold.Hash(), doc.TermsVersion, dates.PartsStart, dates.PartsEnd, dates.LaborStart, dates.LaborEnd, out.AcceptedAt, p.Subject)
	if err != nil {
		return out, warrantyError(err)
	}
	out, err = readWarrantyActivation(ctx, tx, tenant, org, handover)
	if err != nil {
		return out, err
	}
	if err = warrantyEvent(ctx, tx, tenant, out.WarrantyID, "warranty.activated", 1, out); err != nil {
		return out, err
	}
	return out, tx.Commit(ctx)
}
````

### FILE: `internal/warrantyclaim/claim.go`

```yaml
block_id: "GO-CONNECTED-WARRANTY-CLAIM:file22:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "d227722f741ee708ca0d5f50e642b1d73795ec3dc329755b9381690af9d2d9b5"
variables: []
secrets_allowed: false
```

````go
// AUTHORED typed service/approval/stock orchestration contracts. No new stock
// costing algorithm, reimbursement, payroll or financial ledger is defined.
package warrantyclaim

import (
	"elite.local/enterprise/internal/approval"
	"elite.local/enterprise/internal/inventorycontrol"
	"encoding/json"
	"regexp"
	"time"
)

type OpenClaim struct {
	CaseID         string `json:"case_id"`
	HandoverID     string `json:"handover_id"`
	AppointmentID  string `json:"appointment_id"`
	Severity       string `json:"severity"`
	Description    string `json:"description"`
	EvidenceSHA256 string `json:"evidence_sha256"`
}
type Command struct {
	CaseID          string `json:"case_id"`
	CommandID       string `json:"command_id"`
	ExpectedVersion int64  `json:"expected_version,string"`
}

func (r Command) Valid() bool {
	return ValidID(r.CaseID) && ValidID(r.CommandID) && r.ExpectedVersion > 0
}

type Diagnose struct {
	Command
	FaultCode      string `json:"fault_code"`
	Description    string `json:"description"`
	EvidenceSHA256 string `json:"evidence_sha256"`
}
type Part struct {
	LineID               string `json:"line_id"`
	ItemID               string `json:"item_id"`
	BinID                string `json:"bin_id"`
	LotID                string `json:"lot_id,omitempty"`
	Quantity             string `json:"quantity"`
	SpecificReceiptEntry string `json:"specific_receipt_entry,omitempty"`
}
type Plan struct {
	Command
	LaborWork string `json:"labor_work"`
	Parts     []Part `json:"parts"`
}
type DecideRepair struct {
	Command
	PayloadSHA256 string `json:"payload_sha256"`
	Approved      bool   `json:"approved"`
	Reason        string `json:"reason"`
}
type CompleteWork struct {
	Command
	EvidenceSHA256 string `json:"evidence_sha256"`
}
type Quality struct {
	Command
	Passed                   bool   `json:"passed"`
	WorkEvidenceSHA256       string `json:"work_evidence_sha256"`
	EvidenceSHA256           string `json:"evidence_sha256"`
	CorrectionEvidenceSHA256 string `json:"correction_evidence_sha256,omitempty"`
}
type AcceptRepair struct {
	Command
	QualitySHA256  string `json:"quality_sha256"`
	EvidenceSHA256 string `json:"evidence_sha256"`
}
type ReconcileRepair struct {
	Command
	AcceptanceSHA256 string `json:"acceptance_sha256"`
	EvidenceSHA256   string `json:"evidence_sha256"`
}
type CancelRepair struct {
	Command
	Reason         string `json:"reason"`
	EvidenceSHA256 string `json:"evidence_sha256"`
}
type IssuedPart struct {
	LineID        string                           `json:"line_id"`
	ReservationID string                           `json:"reservation_id"`
	Issue         inventorycontrol.BulkIssueResult `json:"issue"`
}
type WorkReceipt struct {
	Request        CompleteWork `json:"request"`
	ApprovalID     string       `json:"approval_id"`
	ApprovedSHA256 string       `json:"approved_sha256"`
	Parts          []IssuedPart `json:"parts"`
}
type Claim struct {
	Context               *RoleContext `json:"context,omitempty"`
	CaseID                string       `json:"case_id"`
	WarrantyID            string       `json:"warranty_id"`
	HandoverID            string       `json:"handover_id"`
	AppointmentID         string       `json:"appointment_id"`
	OrganizationID        string       `json:"organization_id"`
	FactoryOrganizationID string       `json:"factory_organization_id"`
	CustomerSubject       string       `json:"customer_subject"`
	StockUnitID           string       `json:"stock_unit_id"`
	ServiceDate           string       `json:"service_date"`
	State                 string       `json:"state"`
	Version               int64        `json:"version,string"`
	ProfileSHA256         string       `json:"profile_sha256"`
	PartsCovered          bool         `json:"parts_covered"`
	LaborCovered          bool         `json:"labor_covered"`
	Latest                *Step        `json:"latest,omitempty"`
}
type Step struct {
	CaseID        string          `json:"case_id"`
	CommandID     string          `json:"command_id"`
	Version       int64           `json:"version,string"`
	Kind          string          `json:"kind"`
	State         string          `json:"state"`
	Actor         string          `json:"actor"`
	RequestSHA256 string          `json:"request_sha256"`
	PayloadSHA256 string          `json:"payload_sha256"`
	Payload       json.RawMessage `json:"payload"`
	RecordedAt    time.Time       `json:"recorded_at"`
	Replay        bool            `json:"replay"`
}
type RepairPlan struct {
	Schema          string                             `json:"schema"`
	ClaimID         string                             `json:"claim_id"`
	ProfileSHA256   string                             `json:"profile_sha256"`
	DiagnosisSHA256 string                             `json:"diagnosis_sha256"`
	Requester       string                             `json:"requester"`
	Request         Plan                               `json:"request"`
	PartsCovered    bool                               `json:"parts_covered"`
	LaborCovered    bool                               `json:"labor_covered"`
	Excluded        bool                               `json:"excluded"`
	ExpiresAt       time.Time                          `json:"expires_at"`
	Reservations    []inventorycontrol.BulkReservation `json:"reservations"`
}

func Canonical(value any) ([]byte, string, error) {
	raw, err := json.Marshal(value)
	if err != nil {
		return nil, "", err
	}
	return approval.CanonicalPayload(raw)
}

var partQuantity = regexp.MustCompile(`^[1-9][0-9]{0,5}$`)

func (r Plan) Valid() bool {
	if !r.Command.Valid() || len(r.LaborWork) > 4000 || len(r.Parts) > 16 || len(r.Parts) == 0 && len(r.LaborWork) == 0 {
		return false
	}
	seen := map[string]bool{}
	for _, p := range r.Parts {
		if !ValidID(p.LineID) || seen[p.LineID] || !ValidID(p.ItemID) || !ValidID(p.BinID) || p.LotID != "" && !ValidID(p.LotID) || !partQuantity.MatchString(p.Quantity) || p.SpecificReceiptEntry != "" && !ValidID(p.SpecificReceiptEntry) {
			return false
		}
		seen[p.LineID] = true
	}
	return true
}
````

### FILE: `internal/warrantyclaim/contract.go`

```yaml
block_id: "GO-CONNECTED-WARRANTY-CLAIM:file23:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "96a22374a8395b44d6f4ca7c00e32b39baff94169ce34a458ecbf246cad180ad"
variables: []
secrets_allowed: false
```

````go
// AUTHORED typed configuration and calendar representation glue around the
// admitted BC CheckWarranty predicate. No legal period or coverage is inferred.
package warrantyclaim

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"regexp"
	"strings"
	"time"
	_ "time/tzdata"

	"elite.local/enterprise/internal/approval"
	"elite.local/enterprise/internal/warrantycoverage"
)

var (
	ErrInvalid  = errors.New("invalid warranty command")
	ErrConflict = errors.New("warranty conflict")
	ErrNotFound = errors.New("warranty resource unavailable")
	identifier  = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:-]{0,119}$`)
	tenantID    = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
)

type ProfileDocument struct {
	Schema                 string   `json:"schema"`
	Scope                  string   `json:"scope"`
	Algorithm              string   `json:"algorithm"`
	AlgorithmRevision      int      `json:"algorithm_revision"`
	TenantID               string   `json:"tenant_id"`
	OrganizationID         string   `json:"organization_id"`
	FactoryOrganizationID  string   `json:"factory_organization_id"`
	PolicyID               string   `json:"policy_id"`
	TermsVersion           string   `json:"terms_version"`
	TermsText              string   `json:"terms_text"`
	BusinessTimeZone       string   `json:"business_time_zone"`
	PartsDurationDays      int      `json:"parts_duration_days"`
	LaborDurationDays      int      `json:"labor_duration_days"`
	WorkReservationSeconds int      `json:"work_reservation_seconds"`
	FaultExclusions        []string `json:"fault_exclusions"`
	Settlement             string   `json:"settlement"`
	AuthorityReference     string   `json:"authority_reference"`
	DecisionReference      string   `json:"decision_reference"`
}

// Profile is immutable after validation. Document returns a copy, including slices.
type Profile struct {
	document ProfileDocument
	raw      []byte
	hash     string
	zone     *time.Location
}

func (p Profile) Document() ProfileDocument {
	d := p.document
	d.FaultExclusions = append([]string(nil), d.FaultExclusions...)
	return d
}
func (p Profile) Bytes() []byte           { return append([]byte(nil), p.raw...) }
func (p Profile) Hash() string            { return p.hash }
func (p Profile) Scope() (string, string) { return p.document.TenantID, p.document.OrganizationID }
func (p Profile) Valid() bool             { return p.zone != nil && len(p.hash) == 64 }
func ValidID(value string) bool           { return identifier.MatchString(value) }
func SHA(raw []byte) string               { h := sha256.Sum256(raw); return hex.EncodeToString(h[:]) }

func LoadProfile(raw []byte, expectedSHA string) (Profile, error) {
	var empty Profile
	if len(raw) == 0 || len(raw) > 32768 || len(expectedSHA) != 64 || SHA(raw) != expectedSHA {
		return empty, ErrInvalid
	}
	if _, _, err := approval.CanonicalPayload(raw); err != nil {
		return empty, ErrInvalid
	}
	var d ProfileDocument
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	if decoder.Decode(&d) != nil || decoder.Decode(new(any)) != io.EOF {
		return empty, ErrInvalid
	}
	if d.Schema != "elite-warranty-profile/v1" || d.Scope != "MATERIALIZED_PROFILE" || d.Algorithm != "bc-inclusive-fixed-terms" || d.AlgorithmRevision != 1 || !tenantID.MatchString(d.TenantID) || !ValidID(d.OrganizationID) || !ValidID(d.FactoryOrganizationID) || !ValidID(d.PolicyID) || !ValidID(d.TermsVersion) || !ValidID(d.AuthorityReference) || !ValidID(d.DecisionReference) {
		return empty, ErrInvalid
	}
	if len(strings.TrimSpace(d.TermsText)) < 1 || len(d.TermsText) > 16000 || d.PartsDurationDays < 1 || d.LaborDurationDays < 1 || d.PartsDurationDays > 3652500 || d.LaborDurationDays > 3652500 || len(d.FaultExclusions) > 64 || d.WorkReservationSeconds < 60 || d.WorkReservationSeconds > 86400 || d.Settlement != "INTERNAL_WARRANTY_SERVICE_ACKNOWLEDGEMENT" {
		return empty, ErrInvalid
	}
	seen := map[string]bool{}
	for _, reason := range d.FaultExclusions {
		if !ValidID(reason) || seen[reason] {
			return empty, ErrInvalid
		}
		seen[reason] = true
	}
	zone, err := time.LoadLocation(d.BusinessTimeZone)
	if err != nil || d.BusinessTimeZone == "" || d.BusinessTimeZone == "Local" {
		return empty, ErrInvalid
	}
	return Profile{document: d, raw: append([]byte(nil), raw...), hash: expectedSHA, zone: zone}, nil
}

type Dates struct {
	PartsStart string `json:"parts_start"`
	PartsEnd   string `json:"parts_end"`
	LaborStart string `json:"labor_start"`
	LaborEnd   string `json:"labor_end"`
}

// Integer-day formulas are explicit profile data. Go AddDate represents the
// selected calendar-day operation; it does not choose or infer a warranty term.
// Date evaluation itself remains the separately admitted BC source predicate.
func (p Profile) DatesAt(acceptedAt time.Time) (Dates, error) {
	if !p.Valid() || acceptedAt.IsZero() {
		return Dates{}, ErrInvalid
	}
	local := acceptedAt.In(p.zone)
	start := time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, time.UTC)
	partsEnd := start.AddDate(0, 0, p.document.PartsDurationDays)
	laborEnd := start.AddDate(0, 0, p.document.LaborDurationDays)
	if start.Year() < 1753 || start.Year() > 9999 || partsEnd.Year() > 9999 || laborEnd.Year() > 9999 {
		return Dates{}, ErrInvalid
	}
	return Dates{PartsStart: start.Format(time.DateOnly), PartsEnd: partsEnd.Format(time.DateOnly), LaborStart: start.Format(time.DateOnly), LaborEnd: laborEnd.Format(time.DateOnly)}, nil
}

func ordinal(value string) (int32, error) {
	day, err := time.Parse(time.DateOnly, value)
	if err != nil || len(value) != 10 || day.Format(time.DateOnly) != value || day.Year() < 1753 || day.Year() > 9999 {
		return 0, ErrInvalid
	}
	// Positive Gregorian ordinal, independent of local UTC offsets/DST.
	return int32(day.Unix()/86400 + 719163), nil
}
func CheckCoverage(date string, dates Dates) (warrantycoverage.Coverage, error) {
	values := []string{date, dates.PartsStart, dates.PartsEnd, dates.LaborStart, dates.LaborEnd}
	days := make([]int32, 5)
	for i, value := range values {
		parsed, err := ordinal(value)
		if err != nil {
			return warrantycoverage.Coverage{}, err
		}
		days[i] = parsed
	}
	return warrantycoverage.Evaluate(days[0], warrantycoverage.Period{Start: days[1], End: days[2]}, warrantycoverage.Period{Start: days[3], End: days[4]})
}
func (p Profile) FaultExcluded(reason string) bool {
	for _, code := range p.document.FaultExclusions {
		if code == reason {
			return true
		}
	}
	return false
}
````

### FILE: `internal/warrantyclaim/contract_test.go`

```yaml
block_id: "GO-CONNECTED-WARRANTY-CLAIM:file24:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "6e3d531f3f28a06042efa3fc4d95a6ba8776a66b813ba8bc0c027d86d2128fc0"
variables: []
secrets_allowed: false
```

````go
package warrantyclaim

import (
	"encoding/json"
	"testing"
	"time"
)

func fixtureProfileDocument() ProfileDocument {
	return ProfileDocument{
		Schema: "elite-warranty-profile/v1", Scope: "MATERIALIZED_PROFILE", Algorithm: "bc-inclusive-fixed-terms", AlgorithmRevision: 1,
		TenantID: "aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa", OrganizationID: "store", FactoryOrganizationID: "factory",
		PolicyID: "synthetic", TermsVersion: "fixture-v1", TermsText: "Synthetic fixture only; no legal warranty assertion.",
		BusinessTimeZone: "America/New_York", PartsDurationDays: 2, LaborDurationDays: 1, FaultExclusions: []string{"fixture-exclusion"},
		WorkReservationSeconds: 3600, Settlement: "INTERNAL_WARRANTY_SERVICE_ACKNOWLEDGEMENT", AuthorityReference: "fixture-authority", DecisionReference: "fixture-decision",
	}
}
func TestProfileRejectsImplicitOrMutablePolicy(t *testing.T) {
	d := fixtureProfileDocument()
	raw, _ := json.Marshal(d)
	p, err := LoadProfile(raw, SHA(raw))
	if err != nil {
		t.Fatal(err)
	}
	raw[0] = 'X'
	copy := p.Document()
	copy.FaultExclusions[0] = "changed"
	bytes := p.Bytes()
	bytes[0] = 'X'
	if !p.FaultExcluded("fixture-exclusion") || p.FaultExcluded("changed") || p.Bytes()[0] != '{' {
		t.Fatal("profile mutation escaped")
	}
	for _, mutate := range []func(*ProfileDocument){func(d *ProfileDocument) { d.BusinessTimeZone = "Local" }, func(d *ProfileDocument) { d.PartsDurationDays = 0 }, func(d *ProfileDocument) { d.Settlement = "automatic-refund" }, func(d *ProfileDocument) { d.FaultExclusions = []string{"same", "same"} }, func(d *ProfileDocument) { d.AuthorityReference = "" }} {
		d := fixtureProfileDocument()
		mutate(&d)
		raw, _ := json.Marshal(d)
		if _, err := LoadProfile(raw, SHA(raw)); err == nil {
			t.Fatal("invalid profile admitted", string(raw))
		}
	}
	d = fixtureProfileDocument()
	raw, _ = json.Marshal(d)
	duplicate := append([]byte(`{"Schema":"elite-warranty-profile/v1",`), raw[1:]...)
	if _, err := LoadProfile(duplicate, SHA(duplicate)); err == nil {
		t.Fatal("case-insensitive duplicate admitted")
	}
	if _, err := LoadProfile(raw, SHA([]byte("wrong"))); err == nil {
		t.Fatal("unlocked profile admitted")
	}
}
func TestSoldDatesUseBusinessCalendarAndSourceBoundaries(t *testing.T) {
	d := fixtureProfileDocument()
	raw, _ := json.Marshal(d)
	p, err := LoadProfile(raw, SHA(raw))
	if err != nil {
		t.Fatal(err)
	}
	// DST changes at 02:00 local: two calendar days remain March 10, not a
	// floating 48-hour instant. The source predicate includes the ending date.
	dates, err := p.DatesAt(time.Date(2026, 3, 8, 6, 30, 0, 0, time.UTC))
	if err != nil || dates != (Dates{"2026-03-08", "2026-03-10", "2026-03-08", "2026-03-09"}) {
		t.Fatal(dates, err)
	}
	for _, v := range []struct {
		day               string
		any, parts, labor bool
	}{{"2026-03-07", false, false, false}, {"2026-03-08", true, true, true}, {"2026-03-09", true, true, true}, {"2026-03-10", true, true, false}, {"2026-03-11", false, false, false}} {
		c, err := CheckCoverage(v.day, dates)
		if err != nil || c.Any != v.any || c.Parts != v.parts || c.Labor != v.labor {
			t.Fatal(v, c, err)
		}
	}
	if _, err = p.DatesAt(time.Date(9999, 12, 31, 23, 0, 0, 0, time.UTC)); err == nil {
		t.Fatal("out-of-range term admitted")
	}
	if _, err = CheckCoverage("2026-02-30", dates); err == nil {
		t.Fatal("impossible date admitted")
	}
	if _, err = (Profile{}).DatesAt(time.Now()); err == nil {
		t.Fatal("unconfigured profile admitted")
	}
}
````

### FILE: `internal/warrantyclaim/profile_fuzz_test.go`

```yaml
block_id: "GO-CONNECTED-WARRANTY-CLAIM:file25:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "ebdda2da89db3957f59727e12b1107af82959e16399ffd4157447f15fff8d76f"
variables: []
secrets_allowed: false
```

````go
package warrantyclaim

import (
	"encoding/json"
	"testing"
	"time"
)

func FuzzWarrantyProfileCalendar(f *testing.F) {
	raw, _ := json.Marshal(fixtureProfileDocument())
	f.Add(raw)
	f.Add([]byte(`{"schema":"elite-warranty-profile/v1","Schema":"ambiguous"}`))
	f.Add([]byte{0xff, '{', '}'})
	f.Fuzz(func(t *testing.T, raw []byte) {
		if len(raw) > 32769 {
			return
		}
		p, err := LoadProfile(raw, SHA(raw))
		if err != nil {
			return
		}
		if !p.Valid() || p.Hash() != SHA(raw) || SHA(p.Bytes()) != SHA(raw) {
			t.Fatal("accepted profile changed identity")
		}
		doc := p.Document()
		again, _ := json.Marshal(doc)
		round, err := LoadProfile(again, SHA(again))
		if err != nil {
			t.Fatal("accepted semantic profile cannot reload", err)
		}
		at := time.Date(2026, 3, 8, 6, 30, 0, 0, time.UTC)
		dates, e1 := p.DatesAt(at)
		reconstructed, e2 := round.DatesAt(at)
		if (e1 == nil) != (e2 == nil) || e1 == nil && dates != reconstructed {
			t.Fatal("JSON representation changed calendar")
		}
		if e1 == nil {
			first, err := CheckCoverage(dates.PartsStart, dates)
			if err != nil || !first.Parts || !first.Labor {
				t.Fatal("accepted term excludes first date", dates, first, err)
			}
			last, err := CheckCoverage(dates.PartsEnd, dates)
			if err != nil || !last.Parts {
				t.Fatal("accepted parts end is not inclusive", dates, last, err)
			}
		}
		copy := p.Bytes()
		if len(copy) > 0 {
			copy[0] ^= 255
		}
		if SHA(p.Bytes()) != p.Hash() {
			t.Fatal("mutable profile bytes escaped")
		}
	})
}
````

### FILE: `internal/warrantyclaim/terms.go`

```yaml
block_id: "GO-CONNECTED-WARRANTY-CLAIM:file26:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "df054beb8e04032f68871ecaeda88f6e3710d793fe9090e32f0b8a897afecdee"
variables: []
secrets_allowed: false
```

````go
// AUTHORED wire and receipt binding around the admitted coverage predicate.
package warrantyclaim

import (
	"encoding/json"
	"strings"
	"time"
)

type OfferRequest struct {
	QuoteID       string `json:"quote_id"`
	QuoteVersion  int64  `json:"quote_version,string"`
	ProfileSHA256 string `json:"profile_sha256"`
}
type AcknowledgeRequest struct {
	QuoteID        string `json:"quote_id"`
	QuoteVersion   int64  `json:"quote_version,string"`
	ProfileSHA256  string `json:"profile_sha256"`
	EvidenceSHA256 string `json:"evidence_sha256"`
}
type Offer struct {
	QuoteID         string          `json:"quote_id"`
	QuoteVersion    int64           `json:"quote_version,string"`
	CustomerSubject string          `json:"customer_subject"`
	ProfileSHA256   string          `json:"profile_sha256"`
	Profile         json.RawMessage `json:"profile"`
	OfferedBy       string          `json:"offered_by"`
	OfferedAt       time.Time       `json:"offered_at"`
	Acknowledged    bool            `json:"acknowledged"`
	EvidenceSHA256  string          `json:"evidence_sha256,omitempty"`
	Replay          bool            `json:"replay"`
}
type Activation struct {
	WarrantyID      string    `json:"warranty_id"`
	HandoverID      string    `json:"handover_id"`
	QuoteID         string    `json:"quote_id"`
	OrderID         string    `json:"order_id"`
	StockUnitID     string    `json:"stock_unit_id"`
	OrganizationID  string    `json:"organization_id"`
	CustomerSubject string    `json:"customer_subject"`
	ProfileSHA256   string    `json:"profile_sha256"`
	TermsVersion    string    `json:"terms_version"`
	Dates           Dates     `json:"dates"`
	AcceptedAt      time.Time `json:"accepted_at"`
	ActivatedBy     string    `json:"activated_by"`
	ActivatedAt     time.Time `json:"activated_at"`
	Replay          bool      `json:"replay"`
}

func ValidSHA(value string) bool {
	if len(value) != 64 {
		return false
	}
	for _, c := range value {
		if !(c >= '0' && c <= '9' || c >= 'a' && c <= 'f') {
			return false
		}
	}
	return true
}
func StableID(kind string, values ...string) string {
	return kind + "-" + SHA([]byte(strings.Join(values, "\x00")))[:40]
}
````

## 6. Configuration surface

deploy/warranty/profile.reference.json and docs/WARRANTY_REFERENCE.md define the schema and CLI. Disabled feature reads no profile. Enabled host requires exact hash/tenant/org, bounded regular file and migrated enabled guards. No secret is stored or requested. Profile builder refuses existing output and emits deterministic profile/activation manifests.

## 7. Dependency bill

No added Go, Python or Node dependency. Existing fixed Go1.26.8/PG18.6, pgx, payment and service owners remain selected. GO-BC-WARRANTY-COVERAGE-ADAPTER0.1.0 supplies the adapted predicate and MIT source; this pack is AUTHORED orchestration, not BC implementation.

## 8. Apply order

Compose only into an absent/empty destination with MARKDOWN-COMPOSITOR0.3.0 and complete selected dependency closure. Migrations0069/0070 extend existing quote/service/approval history; apply in numeric order. Down refuses populated immutable history, never CASCADE. Preserve existing migrations. The single BC MIT license owner is selected once.

## 9. Verification

WARRANTY_CONNECTED_TERMS_V402.md, WARRANTY_CONNECTED_CLAIM_V402.md and WARRANTY_INTERFACE_AND_PORTABILITY_V402.md bind the unchanged runtime source: connected quote/payment/handover activation, four J4 cases, exact cost740 through two FIFO layers, rollback and COMMIT expiry, human/object negatives, HTTP lost-response recovery, host guard validation, empty/populated downgrade, deterministic CLI and bounded profile fuzz. Full composition SCA/ops/release remain T2803/T2809/T2810; no compilation-as-business-test claim.

## 10. Reconstruction evidence

Exact reconstruction pending at initial emission. Historical runtime receipts are reused only after per-file identity to the frozen tested delta. Publication admission is recorded separately in WARRANTY_CONNECTED_RELEASE_V402.md/json; historical candidate receipts are not rewritten.


V402 narrow admission: four profiles reconstructed exactly; connected source equals the tested frozen40file delta. G0–G8 are recorded in WARRANTY_CONNECTED_RELEASE_V402.md/json. CONDITIONED on selecting the exact sold-policy and complete compatible owners; whole-composition SCA/ops/release remain separate gates, not hidden credential gaps. No frontend claim in this pack.

V402 composed delta: T2804 warranty roles: same authorized transaction projects existing immutable diagnosis/plan/work/quality/acceptance; quote version read, bounded forms and durable GET recovery. No migration/dependency/domain-rule change. WARRANTY_ROLE_RELEASE_V402.md.
