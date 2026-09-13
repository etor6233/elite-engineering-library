# Go Initial Handover API

## 1. Metadata

```yaml
pack_id: "GO-INITIAL-HANDOVER-API"
pack_version: "0.4.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Create and recover initial handover from existing allocation and observed payment under an explicitly materialized profile, with scoped host/API, explicit revision2 durable commercial receipt and separately revalidated current eligibility."
stacks: ["Go 1.26.8", "PostgreSQL 18.6", "Node.js 24.20.0 where frontend selected"]
compatible_with: ["GO-COMMERCE-PRICING-PAYMENT-API0.6.4", "GO-ELECTROMOBILITY-APPLICATION1.10.0", "GO-OFFICIAL-PAYMENT-WEBHOOK-ADAPTERS0.3.2", "TS-GO-API-WEB-BRIDGE0.5.17 where selected"]
incompatible_with: ["unbound tenant/provider/account", "production certification inferred from fixtures", "implicit source or business-policy attribution"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://github.com/stripe/stripe-go/tree/a2df585a800a97fe8ec4ebf551b4449bdb3d90a1", "https://github.com/mercadopago/sdk-go/tree/f910ee53fbb6819e435eaf3d0f800cb1fe74ae09", "https://docs.stripe.com/api/checkout/sessions/object"]
verified_at: "2026-09-11"
```

## 2. Applicability

Use in the current composed franchise reference with the exact declared owner revisions. Missing live credentials leave the provider lane disabled; no credentials are embedded or requested. The local fixture receipts prove the declared contract only. Select the supported materialized profile explicitly; incompatible business options are rejected, not silently invented.

## 3. Architecture contract

Existing Commerce, provider inbox, durable jobs, outbound delivery fence and PostgreSQL remain the state owners. Provider callbacks are signed resource hints; account/mode/order/attempt/currency/amount facts come from the fixed SDK GET. Send acceptance, request binding and pending state share one transaction. Unknown effects are never blindly retried. Handover uses a versioned exact-byte profile and current observations under locks; tenant/org/provider/account/connection/mode are bound. No new financial ledger, statutory decision or implicit physical shipment. Profile and customer routes use existing authorization and server-side identity. Network calls are bounded and outside database locks; queue batches are finite. Roll back the caller pack set before removing its new adapters; do not delete audit rows or provider effects.

## 4. Exact file manifest

```text
CREATE cmd/electromobility-api/commercial_release_host_test.go
CREATE cmd/electromobility-api/handover.go
CREATE cmd/electromobility-api/handover_test.go
CREATE db/migrations/0056_initial_handover.down.sql
CREATE db/migrations/0056_initial_handover.up.sql
CREATE db/migrations/0057_commercial_release.down.sql
CREATE db/migrations/0057_commercial_release.up.sql
CREATE docs/commercial-release.md
CREATE docs/handover-host-activation.md
CREATE docs/handover-profile.md
CREATE docs/initial-handover-reference.md
CREATE internal/franchisejourney/commercial_release.go
CREATE internal/franchisejourney/commercial_release_test.go
CREATE internal/franchisejourney/handover_preparation.go
CREATE internal/franchisejourney/handover_preparation_test.go
CREATE internal/franchisejourney/handover_profile.go
CREATE internal/franchisejourney/handover_profile_fuzz_test.go
CREATE internal/franchisejourney/handover_profile_test.go
CREATE internal/platform/httpapi/commercial_release.go
CREATE internal/platform/httpapi/commercial_release_test.go
CREATE internal/platform/httpapi/handover_preparation.go
CREATE internal/platform/httpapi/handover_preparation_test.go
CREATE internal/platform/postgres/commercial_release.go
CREATE internal/platform/postgres/commercial_release_connected_integration_test.go
CREATE internal/platform/postgres/commercial_release_integration_test.go
CREATE internal/platform/postgres/handover_preparation.go
CREATE internal/platform/postgres/handover_preparation_integration_test.go
CREATE tools/materialize_handover_profile.py
CREATE internal/franchisejourney/handover_context.go
CREATE internal/platform/httpapi/handover_context.go
CREATE internal/platform/httpapi/handover_context_test.go
CREATE internal/platform/postgres/handover_browser_integration_test.go
CREATE internal/platform/postgres/handover_context.go
CREATE internal/platform/postgres/handover_context_integration_test.go
CREATE db/migrations/0068_handover_funding_evidence.up.sql
CREATE db/migrations/0068_handover_funding_evidence.down.sql
```

## 5. Materialization blocks

### FILE: `cmd/electromobility-api/commercial_release_host_test.go`

```yaml
block_id: "GO-INITIAL-HANDOVER-API:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "dad9214d71365df559b359e5304f88dd56c7f0a19559d6cdf82128e980b55643"
variables: []
secrets_allowed: false
```

````go
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"elite.local/enterprise/internal/franchisejourney"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestCommercialReleaseHostRequiresSelectedEffect(t *testing.T) {
	for _, commercial := range []bool{false, true} {
		runtime, env := handoverHostFixture(t)
		if commercial {
			raw, err := os.ReadFile(env["HANDOVER_PROFILE_FILE"])
			if err != nil {
				t.Fatal(err)
			}
			var doc franchisejourney.HandoverProfileDocument
			if err = json.Unmarshal(raw, &doc); err != nil {
				t.Fatal(err)
			}
			doc.AlgorithmRevision = franchisejourney.CommercialReleaseAlgorithmRevision
			doc.Options = franchisejourney.SupportedCommercialReleaseOptions()
			raw, err = json.Marshal(doc)
			if err != nil {
				t.Fatal(err)
			}
			if err = os.WriteFile(env["HANDOVER_PROFILE_FILE"], raw, 0600); err != nil {
				t.Fatal(err)
			}
			sum := sha256.Sum256(raw)
			env["HANDOVER_PROFILE_SHA256"] = hex.EncodeToString(sum[:])
		}
		module, err := selectedInitialHandoverModule(&pgxpool.Pool{}, runtime, func(k string) string { return env[k] })
		if err != nil || module == nil {
			t.Fatal("host activation", err)
		}
		mux := http.NewServeMux()
		module.Register(mux, nil)
		// No credentials or database access: a mounted protected GET rejects
		// missing authorization, while a disabled route is absent altogether.
		for _, route := range []string{"commercial-release-result", "commercial-release-current"} {
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, httptest.NewRequest("GET", "/v1/franchise/handovers/fixture/"+route+"?organization_id=store", nil))
			want := http.StatusNotFound
			if commercial {
				want = http.StatusUnauthorized
			}
			if rec.Code != want {
				t.Fatalf("commercial=%v route=%s status=%d want=%d", commercial, route, rec.Code, want)
			}
		}
	}
}
````

### FILE: `cmd/electromobility-api/handover.go`

```yaml
block_id: "GO-INITIAL-HANDOVER-API:file2:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "3d5cacc0e1cff42be14c7d5dfd248bd511392902cf4cde59ef75af286dc19b94"
variables: []
secrets_allowed: false
```

````go
package main

// AUTHORED activation glue. Profile algorithms and transaction invariants remain
// with the existing handover owner; this host binds its activation to payments.
import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"elite.local/enterprise/internal/franchisejourney"
	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/platform/randomid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var errHandoverConfiguration = errors.New("initial handover activation or payment scope is invalid")

func loadInitialHandoverContract(paymentHost *paymentRuntime, lookup func(string) string) (*franchisejourney.HandoverReleaseContract, error) {
	if lookup == nil {
		return nil, errHandoverConfiguration
	}
	switch lookup("HANDOVER_ENABLED") {
	case "", "false":
		return nil, nil
	case "true":
	default:
		return nil, errHandoverConfiguration
	}
	if paymentHost == nil || paymentHost.worker.Scope.Validate() != nil {
		return nil, errHandoverConfiguration
	}
	scope := paymentHost.worker.Scope
	revision, err := strconv.Atoi(lookup("HANDOVER_PROFILE_REVISION"))
	if err != nil || revision < 1 {
		return nil, errHandoverConfiguration
	}
	path := lookup("HANDOVER_PROFILE_FILE")
	if path == "" || len(path) > 4096 || strings.TrimSpace(path) != path || strings.ContainsAny(path, "\x00\r\n") {
		return nil, errHandoverConfiguration
	}
	activation := franchisejourney.HandoverProfileActivation{Enabled: true, ProfileID: lookup("HANDOVER_PROFILE_ID"), ProfileRevision: revision, DocumentSHA256: lookup("HANDOVER_PROFILE_SHA256"), TenantID: scope.TenantID, OrganizationID: scope.OrganizationID, ExpectedLiveMode: scope.LiveMode, ProviderCode: scope.ProviderCode, ProviderAccountRef: scope.AccountRef, ProviderConnectionID: scope.ConnectionID}
	contract, err := franchisejourney.LoadHandoverProfileFile(path, activation)
	if err != nil || !contract.Valid() || contract.Scope != "MATERIALIZED_PROFILE" || !contract.AllowsScope(scope.TenantID, scope.OrganizationID) || contract.ExpectedLiveMode != scope.LiveMode {
		return nil, errHandoverConfiguration
	}
	return &contract, nil
}
func selectedInitialHandoverModule(pool *pgxpool.Pool, paymentHost *paymentRuntime, lookup func(string) string) (httpapi.EnterpriseModule, error) {
	contract, err := loadInitialHandoverContract(paymentHost, lookup)
	if err != nil || contract == nil {
		return nil, err
	}
	if pool == nil {
		return nil, errHandoverConfiguration
	}
	service, err := franchisejourney.NewHandoverPreparationService(postgres.NewFranchiseJourney(pool), randomid.Generator{}, *contract)
	if err != nil {
		return nil, errHandoverConfiguration
	}
	return initialHandoverHostModule{service: service, commercial: contract.AllowsCommercialRelease(paymentHost.worker.Scope.TenantID, paymentHost.worker.Scope.OrganizationID)}, nil
}

// Both routes share the exact profile-bound owner. The read-only profile never
// mounts a commercial mutation route, even if a caller guesses its URL.
type initialHandoverHostModule struct {
	service    *franchisejourney.HandoverPreparationService
	commercial bool
}

func (m initialHandoverHostModule) Register(mux *http.ServeMux, verifier identity.Verifier) {
	httpapi.InitialHandoverModule{Service: m.service}.Register(mux, verifier)
	httpapi.HandoverContextModule{Service: m.service}.Register(mux, verifier)
	if m.commercial {
		httpapi.CommercialReleaseModule{Service: m.service}.Register(mux, verifier)
	}
}
````

### FILE: `cmd/electromobility-api/handover_test.go`

```yaml
block_id: "GO-INITIAL-HANDOVER-API:file3:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "ec95742ca6b10f946562c1da44370253f1465896d1c3f83d54b396ea608813f4"
variables: []
secrets_allowed: false
```

````go
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"elite.local/enterprise/internal/franchisejourney"
	"elite.local/enterprise/internal/paymentbridge"
)

func handoverHostFixture(t *testing.T) (*paymentRuntime, map[string]string) {
	t.Helper()
	scope := paymentbridge.Scope{TenantID: "018f4d4a-7b36-7a21-8d10-000000000001", OrganizationID: "store", ConnectionID: "connection", ProviderCode: "stripe", AccountRef: "acct_fixture", Currency: "ARS", MinorUnitExponent: 2}
	mode := false
	document := franchisejourney.HandoverProfileDocument{Schema: franchisejourney.HandoverProfileSchema, ProfileID: "franchise-host-fixture", Revision: 2, Algorithm: franchisejourney.HandoverSupportedAlgorithm, AlgorithmRevision: franchisejourney.HandoverSupportedAlgorithmRevision, Scope: "MATERIALIZED_PROFILE", TenantID: scope.TenantID, OrganizationID: scope.OrganizationID, ExpectedLiveMode: &mode, ProviderCode: scope.ProviderCode, ProviderAccountRef: scope.AccountRef, ProviderConnectionID: scope.ConnectionID, MaximumObservationAgeSeconds: 120, Options: franchisejourney.SupportedHandoverProfileOptions(), AuthorityReference: "docs/handover-profile.md", DecisionReference: "PROJECT_HANDOVER_POLICY_DECISION.md"}
	raw, err := json.Marshal(document)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "profile.json")
	if err = os.WriteFile(path, raw, 0600); err != nil {
		t.Fatal(err)
	}
	digest := sha256.Sum256(raw)
	env := map[string]string{"HANDOVER_ENABLED": "true", "HANDOVER_PROFILE_FILE": path, "HANDOVER_PROFILE_ID": document.ProfileID, "HANDOVER_PROFILE_REVISION": "2", "HANDOVER_PROFILE_SHA256": hex.EncodeToString(digest[:])}
	return &paymentRuntime{worker: paymentbridge.Worker{Scope: scope}}, env
}
func TestHandoverHostDisabledDoesNotNeedPaymentOrProfile(t *testing.T) {
	for _, value := range []string{"", "false"} {
		lookup := func(key string) string {
			if key == "HANDOVER_ENABLED" {
				return value
			}
			return "missing-profile"
		}
		module, err := selectedInitialHandoverModule(nil, nil, lookup)
		if err != nil || module != nil {
			t.Fatal("disabled owner attempted activation")
		}
	}
}
func TestHandoverHostBindsMaterializedProfileToPaymentScope(t *testing.T) {
	runtime, env := handoverHostFixture(t)
	contract, err := loadInitialHandoverContract(runtime, func(k string) string { return env[k] })
	if err != nil || contract == nil || contract.Scope != "MATERIALIZED_PROFILE" || contract.ID != "franchise-host-fixture@2" || !contract.AllowsScope(runtime.worker.Scope.TenantID, runtime.worker.Scope.OrganizationID) {
		t.Fatal(contract, err)
	}
	if contract.AllowsScope(runtime.worker.Scope.TenantID, "other") {
		t.Fatal("foreign organization admitted")
	}
	// Missing payment runtime is rejected before any profile-backed route exists.
	if _, err = selectedInitialHandoverModule(nil, nil, func(k string) string { return env[k] }); err == nil {
		t.Fatal("handover activated without payment runtime")
	}
}
func TestHandoverHostRejectsIncompleteOrMismatchedActivation(t *testing.T) {
	runtime, env := handoverHostFixture(t)
	for key := range env {
		if key == "HANDOVER_ENABLED" {
			continue
		}
		t.Run("missing-"+key, func(t *testing.T) {
			lookup := func(k string) string {
				if k == key {
					return ""
				}
				return env[k]
			}
			if _, err := loadInitialHandoverContract(runtime, lookup); err == nil {
				t.Fatal("incomplete activation accepted")
			}
		})
	}
	for _, tc := range []struct{ key, value string }{
		{"HANDOVER_ENABLED", "yes"}, {"HANDOVER_PROFILE_ID", "other-profile"}, {"HANDOVER_PROFILE_REVISION", "3"}, {"HANDOVER_PROFILE_SHA256", env["HANDOVER_PROFILE_SHA256"][:63] + "x"}, {"HANDOVER_PROFILE_FILE", filepath.Join(t.TempDir(), "absent.json")},
	} {
		t.Run(tc.key+"-mismatch", func(t *testing.T) {
			lookup := func(k string) string {
				if k == tc.key {
					return tc.value
				}
				return env[k]
			}
			if _, err := loadInitialHandoverContract(runtime, lookup); err == nil {
				t.Fatal("mismatched activation accepted")
			}
		})
	}
	for _, field := range []string{"tenant", "organization", "provider", "account", "connection", "mode"} {
		t.Run("payment-"+field+"-mismatch", func(t *testing.T) {
			altered := *runtime
			switch field {
			case "tenant":
				altered.worker.Scope.TenantID = "018f4d4a-7b36-7a21-8d10-000000000002"
			case "organization":
				altered.worker.Scope.OrganizationID = "other"
			case "provider":
				altered.worker.Scope.ProviderCode = "mercadopago"
			case "account":
				altered.worker.Scope.AccountRef = "acct_other"
			case "connection":
				altered.worker.Scope.ConnectionID = "other"
			case "mode":
				altered.worker.Scope.LiveMode = true
			}
			if _, err := loadInitialHandoverContract(&altered, func(k string) string { return env[k] }); err == nil {
				t.Fatal("payment scope differs from locked profile")
			}
		})
	}
	// A valid activation never silently falls back to the old LOCAL_FIXTURES literal.
	if err := os.WriteFile(env["HANDOVER_PROFILE_FILE"], []byte(franchisejourney.ReferenceHandoverContractDocument), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := loadInitialHandoverContract(runtime, func(k string) string { return env[k] }); err == nil {
		t.Fatal("non-materialized fallback accepted")
	}
}
````

### FILE: `db/migrations/0056_initial_handover.down.sql`

```yaml
block_id: "GO-INITIAL-HANDOVER-API:file4:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "eb5f6d41bc17ea421a69077eb0f4cfef47fee3942b32c7785d98bd90a9687a6b"
variables: []
secrets_allowed: false
```

````sql
begin;
drop table sales.delivery_handover_preparation;
drop function sales.initial_handover_preparation_immutable();
commit;
````

### FILE: `db/migrations/0056_initial_handover.up.sql`

```yaml
block_id: "GO-INITIAL-HANDOVER-API:file5:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "500631f966f393a39ad8078fc9e7b5ce63af69b9c992124bb31ba15a8cdb8962"
variables: []
secrets_allowed: false
```

````sql
begin;

-- AUTHORED projection/binding glue. No price, credit, tax, stock posting or
-- shipment rules are defined by this table. Those owners remain authoritative.
create table sales.delivery_handover_preparation (
  tenant_id uuid not null,
  handover_id text not null,
  organization_id text not null,
  order_id text not null,
  order_line_id text not null,
  reservation_id text not null,
  payment_attempt_id text not null,
  observation_sha256_hex text not null check (observation_sha256_hex ~ '^[0-9a-f]{64}$'),
  contract_id text not null,
  contract_sha256_hex text not null check (contract_sha256_hex ~ '^[0-9a-f]{64}$'),
  contract_scope text not null,
  check ((contract_scope='LOCAL_FIXTURES' and contract_id='reference-single-unit-observed-payment-v1')
      or (contract_scope='MATERIALIZED_PROFILE' and contract_id ~ '^[a-z][a-z0-9._-]{0,79}@[1-9][0-9]{0,6}$')),
  maximum_observation_age_ns bigint not null check (maximum_observation_age_ns>0 and maximum_observation_age_ns<=900000000000),
  prepared_by_subject text not null check (length(prepared_by_subject) between 1 and 256),
  prepared_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id,handover_id),
  unique (tenant_id,organization_id,order_id,order_line_id),
  foreign key (tenant_id,handover_id) references sales.delivery_handover(tenant_id,handover_id),
  foreign key (tenant_id,organization_id) references org.organization(tenant_id,organization_id),
  foreign key (tenant_id,order_id,order_line_id) references sales.customer_order_line(tenant_id,order_id,line_id),
  foreign key (tenant_id,reservation_id) references inventory.serial_reservation(tenant_id,reservation_id),
  foreign key (tenant_id,payment_attempt_id) references payment.payment_attempt(tenant_id,payment_attempt_id)
);

create function sales.initial_handover_preparation_immutable() returns trigger language plpgsql as $$
begin
  raise exception 'initial handover preparation is immutable';
end;
$$;
create trigger initial_handover_preparation_immutable
  before update or delete on sales.delivery_handover_preparation
  for each row execute function sales.initial_handover_preparation_immutable();

commit;
````

### FILE: `db/migrations/0057_commercial_release.down.sql`

```yaml
block_id: "GO-INITIAL-HANDOVER-API:file6:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "2e3d61b2935649744d90c397d5fa908875ffbf1f3e25d1486563132211729fe4"
variables: []
secrets_allowed: false
```

````sql
begin;
drop table sales.commercial_release_receipt;
commit;
````

### FILE: `db/migrations/0057_commercial_release.up.sql`

```yaml
block_id: "GO-INITIAL-HANDOVER-API:file7:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "f258bb182d8f61ea636a29bc52ba0dc4ab4e3aa66043a9927bb3ce2154597baa"
variables: []
secrets_allowed: false
```

````sql
begin;

-- AUTHORED checkpoint persistence; no shipment, payment or stock journal.
create table sales.commercial_release_receipt (
 tenant_id uuid not null,
 release_id text not null,
 organization_id text not null,
 handover_id text not null,
 order_id text not null,
 payment_attempt_id text not null,
 observation_sha256_hex text not null check(observation_sha256_hex ~ '^[0-9a-f]{64}$'),
 observation_generation bigint not null check(observation_generation > 0),
 handover_version integer not null check(handover_version > 0),
 acceptance_sha256_hex text not null check(acceptance_sha256_hex ~ '^[0-9a-f]{64}$'),
 checklist_id text not null,
 checklist_version integer not null check(checklist_version > 0),
 contract_id text not null,
 contract_sha256_hex text not null check(contract_sha256_hex ~ '^[0-9a-f]{64}$'),
 effect text not null check(effect='COMMIT_COMMERCIAL_RELEASE_RECEIPT'),
 released_by_subject text not null check(length(released_by_subject) between 1 and 256),
 recorded_at timestamptz not null,
 valid_until timestamptz not null,
 primary key(tenant_id,release_id),
 unique(tenant_id,handover_id),
 unique(tenant_id,organization_id,order_id),
 foreign key(tenant_id,handover_id) references sales.delivery_handover_preparation(tenant_id,handover_id),
 foreign key(tenant_id,organization_id) references org.organization(tenant_id,organization_id),
 foreign key(tenant_id,order_id) references sales.customer_order(tenant_id,order_id),
 foreign key(tenant_id,payment_attempt_id) references payment.provider_observation(tenant_id,payment_attempt_id),
 check(valid_until > recorded_at)
);

create trigger commercial_release_receipt_immutable before update or delete on sales.commercial_release_receipt
for each row execute function communication.reject_outbound_delivery_event_mutation();

commit;
````

### FILE: `docs/commercial-release.md`

```yaml
block_id: "GO-INITIAL-HANDOVER-API:file8:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "9ebaa9b8a2f618dcb09e0771aa61eebce6f14636474b6b66d504c22e40f43239"
variables: []
secrets_allowed: false
```

````markdown
# Durable commercial release checkpoint

This owner persists a commercial release receipt after the existing handover has been accepted, its checklist completed, and its bound payment observation and reservation revalidated. The write has exactly three durable effects: one immutable receipt, one completed idempotency record, and one outbox event `commercial-release.recorded`. It does not modify money, stock, order delivery state, tax treatment, inventory costs, a carrier assignment, or a physical shipment journal.

The receipt is historical evidence. `recorded_at` records when its dependencies were checked inside the write transaction; `valid_until` is the earlier of observation expiry and reservation expiry. Neither the existence of a receipt nor an idempotent replay asserts present eligibility. A later callback may invalidate the receipt's current evidence while its history remains recoverable. `ValidateCommercialRelease` returns `current=false` after a hold, new observation generation, changed evidence/policy/acceptance, or expiry.

## Explicit profile activation

The existing `READ_ONLY_ELIGIBILITY` effect and algorithm revision 1 retain their behavior. They cannot call the durable release owner. To select the supported persistence effect, materialize the same bound handover profile with the additional option:

```text
--release-effect commercial-receipt
```

This emits algorithm revision 2 and `release_effect: COMMIT_COMMERCIAL_RELEASE_RECEIPT`. Tenant, organization, provider, account, connection, mode, profile revision and raw-document SHA-256 remain mandatory bindings. No new credentials or environment variables are introduced. The materializer remains disabled unless `--activate` is explicitly provided. A profile with revision 2 and the old read-only effect is rejected, as is the new effect under revision 1. Other business policy algorithms remain unsupported and fail closed.

Prepare the initial handover using that exact profile from the beginning: immutable preparation binds the profile hash, so changing only the release-stage profile cannot upgrade an existing preparation silently.

## HTTP composition

Mount `httpapi.CommercialReleaseModule{Service: service}` beside `InitialHandoverModule`, using the same `*franchisejourney.HandoverPreparationService` produced by the admitted loader. The optional module registers no routes when its service is absent. The existing host helper returns an `InitialHandoverModule` whose concrete service implements `httpapi.CommercialReleaseService`; the host may assert that interface and mount this module. No new database or SDK owner is constructed.

All routes require authenticated `handover:manage` permission and membership in the specified organization. Tenant and actor come from verified claims. The path identifies the handover; order, customer, payment reference, amount and stock are derived by the repository.

| Method/path | Inputs | Result |
|---|---|---|
| `POST /v1/franchise/handovers/{id}/commercial-release` | `Idempotency-Key`; JSON `organization_id`, `observation_sha256` | Immutable receipt; 201; replay header when recovering an identical command |
| `GET /v1/franchise/handovers/{id}/commercial-release-result` | `Idempotency-Key`; query `organization_id` | Historical receipt only |
| `GET /v1/franchise/handovers/{id}/commercial-release-current` | Query `organization_id` | Receipt plus current revalidation result and evaluation time |

Responses use `Cache-Control: no-store`. Unknown command fields, including forged customer, order, actor or amount, are rejected. A second idempotency key cannot create a second release for the same order; replaying a changed actor or command conflicts. A failed write leaves no receipt, idempotency record or release event.

## Transaction boundary and downstream effects

`lockCommercialRelease` reuses `lockInitialHandoverScope`: active bound provider connection, order, payment attempt, provider observation, customer, order line, serialized stock and reservation are held stable. Acceptance is checked under a handover share lock with NOWAIT to avoid inversion with the existing acceptance owner. Freshness is checked after these locks and again after writing the outbox, before committing. The receipt records the observation hash/generation, accepted handover version/evidence, checklist version and profile hash.

The `commercial-release-current` route is a read-only view. A subsequent effect cannot use that response as a durable authorization token. An effect implemented in the PostgreSQL owner must call the same locked checks and compare the receipt's generation, evidence, acceptance and deadline within its own transaction before writing its journal. No such physical consumer is claimed here. Existing Supply shipment posting requires a registered warehouse pick, UOM/FIFO cost and shipment facts; Fulfillment transport requires that actual shipment and a carrier. This receipt does not fabricate them.

## Provenance and evidence

All new code in this owner is `AUTHORED` composition, configuration, persistence, HTTP and test glue. It is not attributed to Microsoft, Stripe, Mercado Pago or any other company. No foreign source was copied and no dependency changed. Existing official SDKs remain pinned in their separate admitted owner; the integral test uses those SDKs through a local HTTP fixture, real callback authentication, durable inbox/queue and PostgreSQL reconciliation.

Existing reused owner paths are `internal/platform/postgres/handover_preparation.go`, `internal/platform/postgres/franchisejourney.go`, `internal/paymentbridge/checkout.go`, `internal/platform/postgres/payment_callback_processor.go` and the materialized profile loader. `sales_shipment.go` and `fulfillment.go` were inspected to establish their additional physical-effect prerequisites, and are not invoked by this receipt.

The candidate gates prove local infrastructure behavior with fixtures and hash-bound profiles. They do not prove production access, fiscal issuance, shipment completion or live account readiness.
````

### FILE: `docs/handover-host-activation.md`

```yaml
block_id: "GO-INITIAL-HANDOVER-API:file9:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "d452ce00db30dde395d6a7d5d2707876aff05facb679c48a8bb5bd942c454c66"
variables: []
secrets_allowed: false
```

````markdown
# Initial handover host activation

The API composes the existing `InitialHandoverModule` only after its materialized profile passes the owner's exact-byte loader and matches the active payment runtime. This wiring is AUTHORED configuration/orchestration glue. The owner retains its supported algorithm, transactional locks, stock/payment checks, idempotency and release evaluation.

`HANDOVER_ENABLED` absent or `false` omits the module and does not read a profile file or require payments. `true` requires an enabled, successfully prepared payment runtime and every `HANDOVER_*` input listed in `.env.example`. Any incomplete, mismatched or unrecognized activation stops startup before HTTP serving and before starting the payment worker. An enabled owner never falls back to the historical `LOCAL_FIXTURES` literal.

Materialize the selected profile using `tools/materialize_handover_profile.py` as described in `docs/handover-profile.md`. Preserve its `profile.json`, activation lock, decision and receipt. Set `HANDOVER_PROFILE_FILE` to that file and copy its exact profile ID, revision and SHA-256 into `HANDOVER_PROFILE_ID`, `HANDOVER_PROFILE_REVISION` and `HANDOVER_PROFILE_SHA256`. These values are configuration, not secrets. A document change requires a new reviewed revision and matching lock.

The host derives the activation's tenant, organization, provider, connection, account and expected mode directly from the prepared payment runtime. The loader must match those six fields to the hash-locked document before returning a valid contract. There are no duplicate `HANDOVER_*` scope overrides. Provider/account are connection bindings, not a new handover policy; payment storage and observations remain responsible for authenticating and reconciling provider facts. The host does not alter the contract after loading it.

After successful activation, these existing authenticated endpoints are registered:

- `POST /v1/franchise/orders/{id}/handover` prepares the initial handover using its existing command contract and idempotency key.
- `GET /v1/franchise/orders/{id}/handover-result` recovers the committed result using the same scope/key.
- `GET /v1/franchise/handovers/{id}/release-check` evaluates current eligibility against the expected observation hash.

All three require the existing `handover:manage` authorization and organization checks. Preparation and eligibility evaluation do not post a shipment or certify live production.

Selecting `--release-effect commercial-receipt` when materializing the profile
binds algorithm revision 2. The host then registers the commercial-release POST,
GET result recovery and GET current-validity endpoints documented in
`docs/commercial-release.md`, using the same profile-bound service instance.
Revision 1 keeps its existing read-only effect and does not mount these routes.
The immutable receipt records the committed commercial checkpoint; its historical
recovery is distinct from current eligibility after refunds, expiry or changed
observations. No fiscal rule, physical shipment or additional credential is
introduced by this host composition.
````

### FILE: `docs/handover-profile.md`

```yaml
block_id: "GO-INITIAL-HANDOVER-API:file10:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "e730cb3c39b23516e13229641628d39d9391b028d28be2738ddd1ec8a56319e9"
variables: []
secrets_allowed: false
```

````markdown
# Materialized initial-handover profile

The supported algorithm is `single-unit-full-observed-payment`, algorithm
revision 1. It requires one allocated serialized unit, an active matching
reservation, full observed order payment with no refunds/disputes/hold, and the
existing customer/checklist acceptance. Its final check is read-only eligibility;
no shipment, fiscal operation or irreversible release is introduced.

Generate a profile in an absent destination using the supplied standard-library
Python materializer. Tenant/organization IDs are configuration, never secrets.
The profile, exact-byte SHA-256 activation lock, documentary decision and receipt
are written together. Omission of `--activate` leaves activation disabled.

```powershell
python tools/materialize_handover_profile.py --output handover-policy --profile-id franchise-initial --revision 1 --tenant-id YOUR_TENANT_UUID --organization-id YOUR_ORGANIZATION_ID --payment-provider stripe --payment-account-ref YOUR_PROVIDER_ACCOUNT_ID --payment-connection-id YOUR_CONNECTION_ID --expected-mode sandbox --activate
```

The host loads `profile.json` with `LoadHandoverProfileFile` and an explicit
`HandoverProfileActivation`. Its tenant, organization, provider, account,
connection and expected provider mode
must also match the server/payment configuration, not request input. A future
`live` configuration uses the same algorithm and loader without Go edits. This
does not make live credentials, deployment or business acceptance PROVEN.

The loader rejects missing/altered documents, duplicate or unknown JSON fields,
unsupported schema/algorithm revisions, unbound profile revisions, scope/mode
mismatches and every option outside the supported algorithm. `policy:true` has
no meaning. A private admitted binding prevents a reconstructed JSON contract
or subsequent field edits from bypassing the loader. Both domain and PostgreSQL
owners enforce the admitted tenant/organization scope. PostgreSQL also requires
the profile's active provider connection, account and exact checkout binding.

Change profile configuration by issuing a new document revision and activation
hash. A preparation records that exact profile ID, revision and SHA-256; another
profile cannot silently reinterpret its release result. Its receipt remains
recoverable as history. Existing exact `LOCAL_FIXTURES` policy and fixtures are
retained for regression, separately from the materialized profile.

Authority/decision references are documentary selection records. They do not
assert signature verification, credentials, live production readiness, legal or
regulatory approval. Requirements beyond this supported algorithm remain a
capability decision rather than an arbitrary override flag.
````

### FILE: `docs/initial-handover-reference.md`

```yaml
block_id: "GO-INITIAL-HANDOVER-API:file11:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "f6ae20dd9b368ff87cd150f285af78bb0a17d52fcc1d3eb9be3ac54a4df3c621"
variables: []
secrets_allowed: false
```

````markdown
# Initial handover reference

This is AUTHORED composition/persistence glue over FranchiseJourney, Commerce,
the serialized reservation owner, and the official-payment observation bridge.
It does not calculate price/tax/credit, post inventory, turn a payment request
into captured money, implement fiscal policy or declare goods delivered.

The initial handover is created by `HandoverPreparationService.Prepare`, using
`postgres.FranchiseJourney`. Callers provide an authorized tenant/operator,
organization, order/line, payment attempt, observed evidence hash and idempotency
key. Customer, stock, reservation, price and currency come from durable owners.
Neither a client-supplied customer/serial/amount nor `payment_attempt=captured`
without a matching provider observation is sufficient.

The retained fixture contract is explicit: reference-single-unit-observed-payment-v1,
scope LOCAL_FIXTURES, the SHA-256 of ReferenceHandoverContractDocument, and an
explicit finite maximum observation age. This fixture selects non-live mode.
Materialized profiles use the same owner through `LoadHandoverProfileFile`, an
exact-byte activation hash, supported algorithm/revision/options and explicit
tenant/organization/mode bindings. See `docs/handover-profile.md`. Both sandbox
and future live configuration require no Go edits; neither inherits production
acceptance from the fixture tests. Missing, altered or unsupported profiles fail.

Prepare locks order, payment attempt and observation in the reconciler's order,
then fixes the customer/line/stock/reservation relationship. It requires one
quantity-one line, reserved matching serialized stock, active customer and
reservation, full matching received amount/currency, no refund/hold, the exact
observation hash and current time after lock acquisition. Unknown, old, future,
partial, disconnected or mismatched facts fail closed. The reference deliberately
does not cover split payments, multiple lines or a new commercial policy.

One transaction creates delivery_handover, its immutable preparation projection,
the completed idempotency receipt and delivery-handover.prepared outbox event.
Sixteen competing identical commands produce one creation and fifteen replays.
A changed command or a second preparation key for that initial order line is
rejected. Existing exception/reschedule owners remain responsible for successors.

`Result` recovers a committed command without retrying its mutation. It preserves
historical preparation facts and hydrates the current handover acceptance and
checklist. Its receipt is not a new release permission: a later callback can
hold payment while the receipt remains readable.

`EvaluateRelease` rechecks the current payment fence and linked owners after
customer checklist/acceptance. It returns a momentary decision in the selected
LOCAL_FIXTURES or MATERIALIZED_PROFILE scope.
It never changes order to delivered or posts a shipment. Any future effect owner
must recheck the guard inside the effect transaction; it cannot consume this
read-only response as durable authorization. A callback hold after acceptance
rejects the next evaluation.

## HTTP composition

Mount `httpapi.InitialHandoverModule{Service: preparationService}` alongside the
existing FranchiseJourney module after constructing the explicit contract.
An unconfigured module mounts no routes. These commands require the existing
handover:manage permission and organization scope from the verified principal:

- POST /v1/franchise/orders/{id}/handover with organization_id, order_line_id,
  payment_attempt_id, observation_sha256 and Idempotency-Key header.
- GET /v1/franchise/orders/{id}/handover-result?organization_id=... with the
  original Idempotency-Key header. Recovery never resends POST automatically.
- GET /v1/franchise/handovers/{id}/release-check?organization_id=...&observation_sha256=...
  checks the current reference guard without moving stock or money.

The order is taken from the route, and tenant/operator from the verified token.
Unknown customer, order or amount body overrides are rejected. Readiness/conflict
errors do not expose database details. Existing checklist and customer acceptance
routes continue unchanged. No T2804 frontend is added by this increment.

## Scope of source reuse

The underlying serial reservation and customer-shipment owners preserve their
existing Microsoft BCApps-derived admission. Bulk CustomerShipment posts costs,
quantities and inventory; an initial serialized handover is a different projection.
This glue does not claim Microsoft wrote it or that shipment semantics were copied.
No new public upstream or Business Central source group is acquired here.
Owner-domain provenance remains governed by its own admission records; this
increment does not reclassify all Commerce or FranchiseJourney as ADAPTED.

## Focused evidence

The staging runner applies all 56 migrations in an isolated loopback PostgreSQL
instance, checks migration0056 down/up, then executes quote creation/acceptance,
stock allocation, a deliberately synthetic payment observation and handover
creation/checklist/acceptance. No fixture inserts a prepared handover. It exercises
17 invalid scope/observation cases, outbox failure rollback, observation expiry
while waiting for a lock, and recovery with exact metadata after server restart.
It runs bounded domain checks and Go vet. PostgreSQL is stopped afterwards.

The original synthetic-observation suite is not an SDK/webhook end-to-end claim.
The added connected suite separately proves Stripe and Mercado Pago SDK checkout,
signed callback, real queue processor and GET reconciliation through this owner,
without fixture insertion of captured payments or prepared handovers. The profile
suite materializes the policy with Python, loads its exact file/hash in Go and
proves the same PostgreSQL flow under the admitted materialized profile. These
are local infrastructure tests; host wiring and remaining gates have their own
evidence and no production acceptance is asserted here.
````

### FILE: `internal/franchisejourney/commercial_release.go`

```yaml
block_id: "GO-INITIAL-HANDOVER-API:file12:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "a97d933c6e3d35444bd78150da2c5f051b33bf00be8372a79ecc574e4a254943"
variables: []
secrets_allowed: false
```

````go
package franchisejourney

// AUTHORED composition glue: persist the explicitly selected commercial
// checkpoint over existing acceptance/payment/stock owners. No shipment, money
// movement, tax treatment or inventory journal is created by this operation.
import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"
)

const CommercialReleaseEffect = "COMMIT_COMMERCIAL_RELEASE_RECEIPT"
const CommercialReleaseAlgorithmRevision = 2

func SupportedCommercialReleaseOptions() HandoverProfileOptions {
	opts := SupportedHandoverProfileOptions()
	opts.ReleaseEffect = CommercialReleaseEffect
	return opts
}

// Explicit additive selection; revisions1/2 retain provider-only coverage.
func SupportedStoredValueReleaseOptions() HandoverProfileOptions {
	opts := SupportedCommercialReleaseOptions()
	opts.PaymentCoverage = "FULL_ORDER_WITH_STORED_VALUE"
	return opts
}
func (p HandoverReleaseContract) AllowsStoredValueFunding(tenant, org string) bool {
	return p.AllowsScope(tenant, org) && p.profile != nil && p.profile.storedValue
}
func supportedHandoverAlgorithm(revision int, options HandoverProfileOptions) bool {
	return revision == HandoverSupportedAlgorithmRevision && options == SupportedHandoverProfileOptions() ||
		revision == CommercialReleaseAlgorithmRevision && options == SupportedCommercialReleaseOptions() || revision == 3 && options == SupportedStoredValueReleaseOptions()
}

func (p HandoverReleaseContract) AllowsCommercialRelease(tenant, organization string) bool {
	return p.AllowsScope(tenant, organization) && p.profile != nil && p.profile.releaseEffect == CommercialReleaseEffect
}

type CommitCommercialReleaseCommand struct {
	OrganizationID    string `json:"organization_id"`
	HandoverID        string `json:"handover_id"`
	ObservationSHA256 string `json:"observation_sha256"`
	IdempotencyKey    string `json:"-"`
}

// CommercialReleaseReceipt is immutable historical evidence of one committed
// checkpoint. It is deliberately not a current shipping or payment authority.
type CommercialReleaseReceipt struct {
	ID                    string    `json:"id"`
	OrganizationID        string    `json:"organization_id"`
	HandoverID            string    `json:"handover_id"`
	OrderID               string    `json:"order_id"`
	PaymentAttemptID      string    `json:"payment_attempt_id"`
	FundingReceiptID      string    `json:"funding_receipt_id,omitempty"`
	ObservationSHA256     string    `json:"observation_sha256"`
	ObservationGeneration int64     `json:"observation_generation"`
	HandoverVersion       int       `json:"handover_version"`
	AcceptanceSHA256      string    `json:"acceptance_sha256"`
	ChecklistID           string    `json:"checklist_id"`
	ChecklistVersion      int       `json:"checklist_version"`
	ContractID            string    `json:"contract_id"`
	ContractSHA256        string    `json:"contract_sha256"`
	Effect                string    `json:"effect"`
	ReleasedBy            string    `json:"released_by"`
	RecordedAt            time.Time `json:"recorded_at"`
	ValidUntil            time.Time `json:"valid_until"`
}

type CurrentCommercialRelease struct {
	Receipt     CommercialReleaseReceipt `json:"receipt"`
	EvaluatedAt time.Time                `json:"evaluated_at"`
	Current     bool                     `json:"current"`
}

type CommercialReleaseRepository interface {
	CommitInitialCommercialRelease(context.Context, string, string, string, string, CommitCommercialReleaseCommand, HandoverReleaseContract, string) (CommercialReleaseReceipt, bool, error)
	InitialCommercialReleaseResult(context.Context, string, string, string, string) (CommercialReleaseReceipt, error)
	ValidateInitialCommercialRelease(context.Context, string, string, string, HandoverReleaseContract) (CurrentCommercialRelease, error)
}

func (s *HandoverPreparationService) commercialRepository(tenant, organization string) (CommercialReleaseRepository, error) {
	if s == nil || !s.contract.AllowsCommercialRelease(tenant, organization) {
		return nil, ErrReleaseConditioned
	}
	r, ok := s.repository.(CommercialReleaseRepository)
	if !ok {
		return nil, ErrReleaseConditioned
	}
	return r, nil
}

func (s *HandoverPreparationService) CommitCommercialRelease(ctx context.Context, tenant, actor string, c CommitCommercialReleaseCommand) (CommercialReleaseReceipt, bool, error) {
	var empty CommercialReleaseReceipt
	r, err := s.commercialRepository(tenant, c.OrganizationID)
	if err != nil {
		return empty, false, err
	}
	if tenant == "" || actor == "" || len(actor) > 256 || !validPreparationID(c.OrganizationID) || !validPreparationID(c.HandoverID) || !sha256Hex(c.ObservationSHA256) || len(c.IdempotencyKey) < 8 || len(c.IdempotencyKey) > 128 {
		return empty, false, ErrInvalid
	}
	b, err := json.Marshal(struct {
		Actor       string
		Command     CommitCommercialReleaseCommand
		ContractSHA string
	}{actor, c, s.contract.DocumentSHA256})
	if err != nil {
		return empty, false, err
	}
	h := sha256.Sum256(b)
	return r.CommitInitialCommercialRelease(ctx, tenant, actor, s.ids.New(), s.ids.New(), c, s.contract, hex.EncodeToString(h[:]))
}

func (s *HandoverPreparationService) CommercialReleaseResult(ctx context.Context, tenant, organization, handover, key string) (CommercialReleaseReceipt, error) {
	r, err := s.commercialRepository(tenant, organization)
	if err != nil {
		return CommercialReleaseReceipt{}, err
	}
	if !validPreparationID(handover) || len(key) < 8 || len(key) > 128 {
		return CommercialReleaseReceipt{}, ErrInvalid
	}
	return r.InitialCommercialReleaseResult(ctx, tenant, organization, handover, key)
}

// ValidateCommercialRelease rechecks current owners under their locks, then
// matches the receipt's generation and accepted handover version. Its response
// is read-only; a later physical effect must revalidate inside its own write.
func (s *HandoverPreparationService) ValidateCommercialRelease(ctx context.Context, tenant, organization, handover string) (CurrentCommercialRelease, error) {
	r, err := s.commercialRepository(tenant, organization)
	if err != nil {
		return CurrentCommercialRelease{}, err
	}
	if !validPreparationID(handover) {
		return CurrentCommercialRelease{}, ErrInvalid
	}
	return r.ValidateInitialCommercialRelease(ctx, tenant, organization, handover, s.contract)
}
````

### FILE: `internal/franchisejourney/commercial_release_test.go`

```yaml
block_id: "GO-INITIAL-HANDOVER-API:file13:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "2c8cac2695aa38bbb7199223a7b0e0fc97497d40624833d119d19ffcf4e7dabc"
variables: []
secrets_allowed: false
```

````go
package franchisejourney

import (
	"errors"
	"testing"
)

func TestCommercialReleaseRequiresExplicitVersionedEffect(t *testing.T) {
	for _, tc := range []struct {
		name               string
		revision           int
		effect             string
		admitted, writable bool
	}{
		{"readonly", 1, "READ_ONLY_ELIGIBILITY", true, false},
		{"commercial", 2, CommercialReleaseEffect, true, true},
		{"readonly-cannot-upgrade", 1, CommercialReleaseEffect, false, false},
		{"revision-alone-cannot-upgrade", 2, "READ_ONLY_ELIGIBILITY", false, false},
		{"physical-shipment-denied", 2, "POST_SHIPMENT", false, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			d, a := profileFixture(false)
			d.AlgorithmRevision = tc.revision
			d.Options.ReleaseEffect = tc.effect
			raw, a := lockProfile(t, d, a)
			p, err := LoadHandoverProfile(raw, a)
			if tc.admitted {
				if err != nil || p.AllowsCommercialRelease(d.TenantID, d.OrganizationID) != tc.writable {
					t.Fatal("effect admission", err)
				}
			} else if !errors.Is(err, ErrReleaseConditioned) {
				t.Fatal("unsupported effect admitted", err)
			}
			if p.AllowsCommercialRelease(d.TenantID, "foreign") {
				t.Fatal("foreign organization admitted")
			}
		})
	}
}
````

### FILE: `internal/franchisejourney/handover_preparation.go`

```yaml
block_id: "GO-INITIAL-HANDOVER-API:file14:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "a76a026a7218f7dde93978e3f8239dd8340ce8dd29f06f62366d9693b409865b"
variables: []
secrets_allowed: false
```

````go
package franchisejourney

// AUTHORED composition glue. This projects existing order, allocation and
// payment observations into a handover; it neither prices an order nor posts a
// shipment. The explicitly selected reference contract is not a live policy.
import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

var ErrReleaseConditioned = errors.New("handover release contract not admitted")

// Exact reference contract, deliberately scoped to synthetic/offline execution.
// The binding conditions are selected by the reference, not a fiscal or credit rule.
const ReferenceHandoverContractDocument = "reference-single-unit-observed-payment-v1\nLOCAL_FIXTURES\nOne active customer order line, quantity one, allocated serialized stock and active matching reservation.\nThe configured provider observation must match the complete order amount and currency, with no refund, dispute, hold or unobserved state.\nPreparation does not post inventory, establish fiscal compliance, deliver goods or authorize production.\nThe observation maximum age is explicitly configured between one nanosecond and fifteen minutes and included in the command hash.\n"

type HandoverReleaseContract struct {
	ID                    string        `json:"id"`
	DocumentSHA256        string        `json:"document_sha256"`
	Scope                 string        `json:"scope"`
	ExpectedLiveMode      bool          `json:"expected_live_mode"`
	MaximumObservationAge time.Duration `json:"-"`
	profile               *handoverProfileBinding
}

func (p HandoverReleaseContract) Valid() bool {
	if p.profile != nil {
		b := p.profile
		return p.ID == b.id && p.DocumentSHA256 == b.sha && p.Scope == "MATERIALIZED_PROFILE" && p.ExpectedLiveMode == b.mode && p.MaximumObservationAge == b.age
	}
	decoded, err := hex.DecodeString(p.DocumentSHA256)
	expected := sha256.Sum256([]byte(ReferenceHandoverContractDocument))
	return p.ID == "reference-single-unit-observed-payment-v1" && p.Scope == "LOCAL_FIXTURES" && !p.ExpectedLiveMode &&
		err == nil && len(decoded) == 32 && strings.ToLower(p.DocumentSHA256) == p.DocumentSHA256 &&
		p.DocumentSHA256 == hex.EncodeToString(expected[:]) &&
		p.MaximumObservationAge > 0 && p.MaximumObservationAge <= 15*time.Minute
}

func validPreparationID(value string) bool {
	return strings.TrimSpace(value) == value && value != "" && len(value) <= 128
}

type PrepareHandoverCommand struct {
	OrganizationID    string `json:"organization_id"`
	OrderID           string `json:"order_id"`
	OrderLineID       string `json:"order_line_id"`
	PaymentAttemptID  string `json:"payment_attempt_id"`
	FundingReceiptID  string `json:"funding_receipt_id,omitempty"`
	ObservationSHA256 string `json:"observation_sha256"`
	IdempotencyKey    string `json:"-"`
}

type HandoverPreparation struct {
	Handover          Handover  `json:"handover"`
	OrderLineID       string    `json:"order_line_id"`
	ReservationID     string    `json:"reservation_id"`
	PaymentAttemptID  string    `json:"payment_attempt_id"`
	FundingReceiptID  string    `json:"funding_receipt_id,omitempty"`
	ObservationSHA256 string    `json:"observation_sha256"`
	ContractID        string    `json:"contract_id"`
	ContractSHA256    string    `json:"contract_sha256"`
	PreparedBy        string    `json:"prepared_by"`
	PreparedAt        time.Time `json:"prepared_at"`
}

type HandoverPreparationRepository interface {
	PrepareInitialHandover(context.Context, string, string, string, string, PrepareHandoverCommand, HandoverReleaseContract, string) (HandoverPreparation, bool, error)
	InitialHandoverResult(context.Context, string, string, string, string) (HandoverPreparation, error)
}

type HandoverReferenceRelease struct {
	HandoverID        string    `json:"handover_id"`
	OrganizationID    string    `json:"organization_id"`
	ObservationSHA256 string    `json:"observation_sha256"`
	ContractSHA256    string    `json:"contract_sha256"`
	EvaluatedAt       time.Time `json:"evaluated_at"`
	Scope             string    `json:"scope"`
	Eligible          bool      `json:"eligible"`
}

// EvaluateRelease reports a current reference decision. It is not a durable
// shipping authorization; the effect owner must revalidate inside its own write.
func (s *HandoverPreparationService) EvaluateRelease(ctx context.Context, tenant, organization, handover, observationSHA string) (HandoverReferenceRelease, error) {
	if s == nil || !s.contract.AllowsScope(tenant, organization) {
		return HandoverReferenceRelease{}, ErrReleaseConditioned
	}
	if tenant == "" || !validPreparationID(organization) || !validPreparationID(handover) || !sha256Hex(observationSHA) {
		return HandoverReferenceRelease{}, ErrInvalid
	}
	repo, ok := s.repository.(interface {
		EvaluateInitialHandoverRelease(context.Context, string, string, string, string, HandoverReleaseContract) (HandoverReferenceRelease, error)
	})
	if !ok {
		return HandoverReferenceRelease{}, ErrReleaseConditioned
	}
	return repo.EvaluateInitialHandoverRelease(ctx, tenant, organization, handover, observationSHA, s.contract)
}

type HandoverPreparationService struct {
	repository HandoverPreparationRepository
	ids        IDGenerator
	contract   HandoverReleaseContract
}

func NewHandoverPreparationService(repository HandoverPreparationRepository, ids IDGenerator, contract HandoverReleaseContract) (*HandoverPreparationService, error) {
	if repository == nil || ids == nil || !contract.Valid() {
		return nil, ErrReleaseConditioned
	}
	return &HandoverPreparationService{repository: repository, ids: ids, contract: contract}, nil
}

func (s *HandoverPreparationService) Prepare(ctx context.Context, tenant, actor string, command PrepareHandoverCommand) (HandoverPreparation, bool, error) {
	var empty HandoverPreparation
	if s == nil || !s.contract.AllowsScope(tenant, command.OrganizationID) {
		return empty, false, ErrReleaseConditioned
	}
	if tenant == "" || actor == "" || len(actor) > 256 || !validPreparationID(command.OrganizationID) || !validPreparationID(command.OrderID) ||
		!validPreparationID(command.OrderLineID) || !validHandoverFunding(command, s.contract, tenant) || len(command.IdempotencyKey) < 8 || len(command.IdempotencyKey) > 128 {
		return empty, false, ErrInvalid
	}
	hash, err := hex.DecodeString(command.ObservationSHA256)
	if err != nil || len(hash) != 32 || strings.ToLower(command.ObservationSHA256) != command.ObservationSHA256 {
		return empty, false, ErrInvalid
	}
	encoded, err := json.Marshal(struct {
		Actor                 string
		Command               PrepareHandoverCommand
		Contract              HandoverReleaseContract
		MaximumAgeNanoseconds int64
	}{actor, command, s.contract, int64(s.contract.MaximumObservationAge)})
	if err != nil {
		return empty, false, err
	}
	sum := sha256.Sum256(encoded)
	return s.repository.PrepareInitialHandover(ctx, tenant, actor, s.ids.New(), s.ids.New(), command, s.contract, hex.EncodeToString(sum[:]))
}

// Result recovers the committed receipt without resending the mutation. Scope
// identity is server-authenticated; no caller-supplied customer or stock is used.
func (s *HandoverPreparationService) Result(ctx context.Context, tenant, organization, order, key string) (HandoverPreparation, error) {
	if s == nil || !s.contract.AllowsScope(tenant, organization) {
		return HandoverPreparation{}, ErrReleaseConditioned
	}
	if tenant == "" || !validPreparationID(organization) || !validPreparationID(order) || len(key) < 8 || len(key) > 128 {
		return HandoverPreparation{}, ErrInvalid
	}
	return s.repository.InitialHandoverResult(ctx, tenant, organization, order, key)
}

func validHandoverFunding(c PrepareHandoverCommand, p HandoverReleaseContract, tenant string) bool {
	if c.FundingReceiptID != "" {
		return c.PaymentAttemptID == "" && validPreparationID(c.FundingReceiptID) && p.AllowsStoredValueFunding(tenant, c.OrganizationID)
	}
	return validPreparationID(c.PaymentAttemptID)
}
````

### FILE: `internal/franchisejourney/handover_preparation_test.go`

```yaml
block_id: "GO-INITIAL-HANDOVER-API:file15:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "ab1c83aa726d99cf6f1ce4a4fd57a91ddfb84491a45bb785dda1535d374b5fc9"
variables: []
secrets_allowed: false
```

````go
package franchisejourney

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"testing"
	"time"
)

type preparationRepoProbe struct {
	calls       int
	actor, hash string
	command     PrepareHandoverCommand
}

func (p *preparationRepoProbe) PrepareInitialHandover(_ context.Context, _, actor, _, _ string, command PrepareHandoverCommand, _ HandoverReleaseContract, hash string) (HandoverPreparation, bool, error) {
	p.calls++
	p.actor = actor
	p.hash = hash
	p.command = command
	return HandoverPreparation{}, false, nil
}
func (p *preparationRepoProbe) InitialHandoverResult(context.Context, string, string, string, string) (HandoverPreparation, error) {
	p.calls++
	return HandoverPreparation{}, nil
}

type preparationIDs struct{}

func (preparationIDs) New() string { return "generated-id" }
func referencePreparationContract() HandoverReleaseContract {
	hash := sha256.Sum256([]byte(ReferenceHandoverContractDocument))
	return HandoverReleaseContract{ID: "reference-single-unit-observed-payment-v1", DocumentSHA256: hex.EncodeToString(hash[:]), Scope: "LOCAL_FIXTURES", MaximumObservationAge: 5 * time.Minute}
}
func TestHandoverPreparationRequiresExplicitExactReferenceContract(t *testing.T) {
	cases := []HandoverReleaseContract{{}, referencePreparationContract(), referencePreparationContract(), referencePreparationContract(), referencePreparationContract()}
	cases[1].Scope = "PRODUCTION"
	cases[2].DocumentSHA256 = strings.Repeat("a", 64)
	cases[3].MaximumObservationAge = 0
	cases[4].MaximumObservationAge = time.Hour
	for _, contract := range cases {
		if _, err := NewHandoverPreparationService(&preparationRepoProbe{}, preparationIDs{}, contract); !errors.Is(err, ErrReleaseConditioned) {
			t.Fatalf("unsafe contract accepted: %v", err)
		}
	}
}
func TestHandoverPreparationRequestBindsActorAndObservation(t *testing.T) {
	repo := &preparationRepoProbe{}
	service, err := NewHandoverPreparationService(repo, preparationIDs{}, referencePreparationContract())
	if err != nil {
		t.Fatal(err)
	}
	command := PrepareHandoverCommand{OrganizationID: "store", OrderID: "order", OrderLineID: "line", PaymentAttemptID: "payment", ObservationSHA256: strings.Repeat("a", 64), IdempotencyKey: "prepare-key"}
	if _, _, err = service.Prepare(context.Background(), "tenant", "operator", command); err != nil {
		t.Fatal(err)
	}
	original := repo.hash
	if repo.actor != "operator" || repo.command.OrderID != "order" {
		t.Fatal("principal or scope lost")
	}
	if _, _, err = service.Prepare(context.Background(), "tenant", "other", command); err != nil || repo.hash == original {
		t.Fatal("actor not bound")
	}
	command.ObservationSHA256 = strings.Repeat("b", 64)
	if _, _, err = service.Prepare(context.Background(), "tenant", "operator", command); err != nil || repo.hash == original {
		t.Fatal("observation not bound")
	}
}
func TestHandoverPreparationRejectsInvalidInputsBeforeRepository(t *testing.T) {
	repo := &preparationRepoProbe{}
	service, _ := NewHandoverPreparationService(repo, preparationIDs{}, referencePreparationContract())
	command := PrepareHandoverCommand{OrganizationID: "store", OrderID: "order", OrderLineID: "line", PaymentAttemptID: "payment", ObservationSHA256: strings.Repeat("a", 64), IdempotencyKey: "prepare-key"}
	cases := []PrepareHandoverCommand{command, command, command, command, command}
	cases[0].OrganizationID = ""
	cases[1].OrderLineID = " "
	cases[2].PaymentAttemptID = ""
	cases[3].ObservationSHA256 = "captured"
	cases[4].IdempotencyKey = "short"
	for _, bad := range cases {
		if _, _, err := service.Prepare(context.Background(), "tenant", "operator", bad); !errors.Is(err, ErrInvalid) {
			t.Fatal(err)
		}
	}
	if repo.calls != 0 {
		t.Fatal("invalid request reached repository")
	}
}
````

### FILE: `internal/franchisejourney/handover_profile.go`

```yaml
block_id: "GO-INITIAL-HANDOVER-API:file16:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "37f8c5061e3ae2899a0137e47204c808368029a373e3b5c193c075dd239b30ee"
variables: []
secrets_allowed: false
```

````go
package franchisejourney

// AUTHORED configuration/validation glue. This selects the existing supported
// owner algorithm; it does not introduce pricing, credit, shipping or tax rules.
import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"os"
	"regexp"
	"strconv"
	"time"
)

const HandoverProfileSchema = "elite-handover-profile/v1"
const HandoverSupportedAlgorithm = "single-unit-full-observed-payment"
const HandoverSupportedAlgorithmRevision = 1

type HandoverProfileOptions struct {
	Quantity        int    `json:"quantity"`
	PaymentCoverage string `json:"payment_coverage"`
	StockSelection  string `json:"stock_selection"`
	Reservation     string `json:"reservation"`
	Refunds         string `json:"refunds"`
	Disputes        string `json:"disputes"`
	Acceptance      string `json:"acceptance"`
	ReleaseEffect   string `json:"release_effect"`
}

func SupportedHandoverProfileOptions() HandoverProfileOptions {
	return HandoverProfileOptions{Quantity: 1, PaymentCoverage: "FULL_ORDER", StockSelection: "ALLOCATED_SERIALIZED_UNIT", Reservation: "ACTIVE_MATCHING", Refunds: "ZERO", Disputes: "DENY", Acceptance: "CUSTOMER_AND_REQUIRED_CHECKLIST", ReleaseEffect: "READ_ONLY_ELIGIBILITY"}
}

type HandoverProfileDocument struct {
	Schema                       string                 `json:"schema"`
	ProfileID                    string                 `json:"profile_id"`
	Revision                     int                    `json:"revision"`
	Algorithm                    string                 `json:"algorithm"`
	AlgorithmRevision            int                    `json:"algorithm_revision"`
	Scope                        string                 `json:"scope"`
	TenantID                     string                 `json:"tenant_id"`
	OrganizationID               string                 `json:"organization_id"`
	ProviderCode                 string                 `json:"provider_code"`
	ProviderAccountRef           string                 `json:"provider_account_ref"`
	ProviderConnectionID         string                 `json:"provider_connection_id"`
	ExpectedLiveMode             *bool                  `json:"expected_live_mode"`
	MaximumObservationAgeSeconds int64                  `json:"maximum_observation_age_seconds"`
	Options                      HandoverProfileOptions `json:"options"`
	// Documentary references record who selected the supported profile. They are
	// not runtime credential proof, a signature or live readiness certification.
	AuthorityReference string `json:"authority_reference"`
	DecisionReference  string `json:"decision_reference"`
}
type HandoverProfileActivation struct {
	Enabled              bool   `json:"enabled"`
	ProfileID            string `json:"profile_id"`
	ProfileRevision      int    `json:"profile_revision"`
	DocumentSHA256       string `json:"document_sha256"`
	TenantID             string `json:"tenant_id"`
	OrganizationID       string `json:"organization_id"`
	ProviderCode         string `json:"provider_code"`
	ProviderAccountRef   string `json:"provider_account_ref"`
	ProviderConnectionID string `json:"provider_connection_id"`
	ExpectedLiveMode     bool   `json:"expected_live_mode"`
}
type handoverProfileBinding struct {
	id, sha, tenant, organization string
	provider, account, connection string
	releaseEffect                 string
	storedValue                   bool
	mode                          bool
	age                           time.Duration
}

var handoverProfileID = regexp.MustCompile(`^[a-z][a-z0-9._-]{0,79}$`)
var handoverTenantID = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

// LoadHandoverProfile requires an external exact-byte lock plus an explicit
// activation matching the deployment's identity and provider mode. A JSON
// policy:true or an unrecognized algorithm cannot admit arbitrary behavior.
func LoadHandoverProfile(raw []byte, a HandoverProfileActivation) (HandoverReleaseContract, error) {
	var empty HandoverReleaseContract
	if len(raw) == 0 || len(raw) > 32768 || !a.Enabled || !sha256Hex(a.DocumentSHA256) {
		return empty, ErrReleaseConditioned
	}
	digest := sha256.Sum256(raw)
	if hex.EncodeToString(digest[:]) != a.DocumentSHA256 || !uniqueProfileJSON(raw) {
		return empty, ErrReleaseConditioned
	}
	var d HandoverProfileDocument
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	if decoder.Decode(&d) != nil || decoder.Decode(new(any)) != io.EOF {
		return empty, ErrReleaseConditioned
	}
	if d.Schema != HandoverProfileSchema || !handoverProfileID.MatchString(d.ProfileID) || d.Revision < 1 || d.Revision > 1000000 || d.Algorithm != HandoverSupportedAlgorithm || !supportedHandoverAlgorithm(d.AlgorithmRevision, d.Options) || d.Scope != "MATERIALIZED_PROFILE" || !handoverTenantID.MatchString(d.TenantID) || !validPreparationID(d.OrganizationID) || d.ExpectedLiveMode == nil || d.MaximumObservationAgeSeconds < 1 || d.MaximumObservationAgeSeconds > 900 || !validPreparationID(d.AuthorityReference) || !validPreparationID(d.DecisionReference) {
		return empty, ErrReleaseConditioned
	}
	if (d.ProviderCode != "stripe" && d.ProviderCode != "mercadopago") || !validPreparationID(d.ProviderAccountRef) || !validPreparationID(d.ProviderConnectionID) {
		return empty, ErrReleaseConditioned
	}
	if a.ProfileID != d.ProfileID || a.ProfileRevision != d.Revision || a.TenantID != d.TenantID || a.OrganizationID != d.OrganizationID || a.ExpectedLiveMode != *d.ExpectedLiveMode || a.ProviderCode != d.ProviderCode || a.ProviderAccountRef != d.ProviderAccountRef || a.ProviderConnectionID != d.ProviderConnectionID {
		return empty, ErrReleaseConditioned
	}
	id := d.ProfileID + "@" + strconv.Itoa(d.Revision)
	age := time.Duration(d.MaximumObservationAgeSeconds) * time.Second
	bound := &handoverProfileBinding{id: id, sha: a.DocumentSHA256, tenant: d.TenantID, organization: d.OrganizationID, mode: *d.ExpectedLiveMode, age: age, provider: d.ProviderCode, account: d.ProviderAccountRef, connection: d.ProviderConnectionID, releaseEffect: d.Options.ReleaseEffect, storedValue: d.AlgorithmRevision == 3 && d.Options == SupportedStoredValueReleaseOptions()}
	return HandoverReleaseContract{ID: id, DocumentSHA256: a.DocumentSHA256, Scope: d.Scope, ExpectedLiveMode: *d.ExpectedLiveMode, MaximumObservationAge: age, profile: bound}, nil
}

func LoadHandoverProfileFile(path string, a HandoverProfileActivation) (HandoverReleaseContract, error) {
	if path == "" {
		return HandoverReleaseContract{}, ErrReleaseConditioned
	}
	f, err := os.Open(path)
	if err != nil {
		return HandoverReleaseContract{}, ErrReleaseConditioned
	}
	defer f.Close()
	raw, err := io.ReadAll(io.LimitReader(f, 32769))
	if err != nil {
		return HandoverReleaseContract{}, ErrReleaseConditioned
	}
	return LoadHandoverProfile(raw, a)
}

// Duplicate keys have ambiguous human review semantics even with a byte lock.
func uniqueProfileJSON(raw []byte) bool {
	d := json.NewDecoder(bytes.NewReader(raw))
	d.UseNumber()
	var visit func() bool
	visit = func() bool {
		token, err := d.Token()
		if err != nil {
			return false
		}
		delim, ok := token.(json.Delim)
		if !ok {
			return true
		}
		switch delim {
		case '{':
			seen := map[string]bool{}
			for d.More() {
				k, err := d.Token()
				key, ok := k.(string)
				if err != nil || !ok || seen[key] {
					return false
				}
				seen[key] = true
				if !visit() {
					return false
				}
			}
			end, err := d.Token()
			return err == nil && end == json.Delim('}')
		case '[':
			for d.More() {
				if !visit() {
					return false
				}
			}
			end, err := d.Token()
			return err == nil && end == json.Delim(']')
		default:
			return false
		}
	}
	if !visit() {
		return false
	}
	_, err := d.Token()
	return err == io.EOF
}

func (p HandoverReleaseContract) AllowsScope(tenant, organization string) bool {
	if !p.Valid() {
		return false
	}
	return p.profile == nil || p.profile.tenant == tenant && p.profile.organization == organization
}

func (p HandoverReleaseContract) PaymentBinding() (provider, account, connection string, required bool) {
	if p.profile == nil {
		return "", "", "", false
	}
	return p.profile.provider, p.profile.account, p.profile.connection, true
}
````

### FILE: `internal/franchisejourney/handover_profile_fuzz_test.go`

```yaml
block_id: "GO-INITIAL-HANDOVER-API:file17:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "6053a13ae143e627169020944aad358ba2430456591e964e7f2f60498eb3bc97"
variables: []
secrets_allowed: false
```

````go
package franchisejourney

// AUTHORED verification glue. The seed is the materialized sandbox profile
// used by the connected PostgreSQL journey, with synthetic account identifiers.
import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func fuzzHandoverProfileSeed(mode bool) ([]byte, HandoverProfileActivation) {
	d := HandoverProfileDocument{
		Schema: HandoverProfileSchema, ProfileID: "franchise-reference", Revision: 1,
		Algorithm: HandoverSupportedAlgorithm, AlgorithmRevision: 1,
		Scope: "MATERIALIZED_PROFILE", TenantID: "018f4d4a-7b36-7a21-8d10-000000000001", OrganizationID: "store",
		ProviderCode: "stripe", ProviderAccountRef: "acct_fixture", ProviderConnectionID: "checkout", ExpectedLiveMode: &mode,
		MaximumObservationAgeSeconds: 300, Options: SupportedHandoverProfileOptions(),
		AuthorityReference: "docs/initial-handover-reference.md", DecisionReference: "PROJECT_HANDOVER_POLICY_DECISION.md",
	}
	raw, _ := json.Marshal(d)
	a := HandoverProfileActivation{Enabled: true, ProfileID: d.ProfileID, ProfileRevision: d.Revision, TenantID: d.TenantID, OrganizationID: d.OrganizationID,
		ProviderCode: d.ProviderCode, ProviderAccountRef: d.ProviderAccountRef, ProviderConnectionID: d.ProviderConnectionID, ExpectedLiveMode: mode}
	return raw, a
}

func fuzzProfileHash(raw []byte) string {
	digest := sha256.Sum256(raw)
	return hex.EncodeToString(digest[:])
}

func FuzzHandoverProfileAdmission(f *testing.F) {
	base, _ := fuzzHandoverProfileSeed(false)
	live, _ := fuzzHandoverProfileSeed(true)
	f.Add(base, uint8(0))
	f.Add(live, uint8(1))
	var commercial HandoverProfileDocument
	_ = json.Unmarshal(base, &commercial)
	commercial.AlgorithmRevision = 2
	commercial.Options = SupportedCommercialReleaseOptions()
	commercialRaw, _ := json.Marshal(commercial)
	f.Add(commercialRaw, uint8(0))
	for i, seed := range [][]byte{base, live, commercialRaw} {
		_, activation := fuzzHandoverProfileSeed(i == 1)
		activation.DocumentSHA256 = fuzzProfileHash(seed)
		if contract, err := LoadHandoverProfile(seed, activation); err != nil || !contract.Valid() {
			f.Fatalf("real supported profile seed %d was rejected: %v", i, err)
		}
	}
	for variant := uint8(2); variant < 10; variant++ {
		f.Add(base, variant)
	}
	f.Add(append([]byte(`{"schema":"duplicate",`), base[1:]...), uint8(0))
	f.Add(bytes.Replace(base, []byte(`"expected_live_mode":false,`), nil, 1), uint8(0))
	f.Add(bytes.Replace(base, []byte(`"scope":"MATERIALIZED_PROFILE"`), []byte(`"scope":"LOCAL_FIXTURES"`), 1), uint8(0))
	f.Add([]byte(`{"policy":true}`), uint8(0))
	f.Add(bytes.Replace(base, []byte(`"algorithm_revision":1`), []byte(`"algorithm_revision":999`), 1), uint8(0))
	f.Fuzz(func(t *testing.T, raw []byte, variant uint8) {
		if len(raw) > 32769 {
			t.Skip()
		}
		variant %= 10
		_, a := fuzzHandoverProfileSeed(variant == 1)
		a.DocumentSHA256 = fuzzProfileHash(raw)
		switch variant {
		case 2:
			a.DocumentSHA256 = strings.Repeat("0", 64)
		case 3:
			a.TenantID = "018f4d4a-7b36-7a21-8d10-000000000099"
		case 4:
			a.OrganizationID = "another-store"
		case 5:
			a.ProviderAccountRef = "another-account"
		case 6:
			a.ProviderConnectionID = "another-connection"
		case 7:
			a.ProviderCode = "mercadopago"
		case 8:
			a.Enabled = false
		case 9:
			a.ProfileRevision = 0
		}
		contract, err := LoadHandoverProfile(raw, a)
		if err != nil {
			if contract.Valid() {
				t.Fatal("rejected profile leaked an admitted contract")
			}
			return
		}
		// Locks are external inputs. A mutated document may legitimately select
		// any supported revision, but must bind exactly to this activation.
		if !contract.Valid() || !contract.AllowsScope(a.TenantID, a.OrganizationID) || contract.AllowsScope(a.TenantID+"-foreign", a.OrganizationID) || contract.AllowsScope(a.TenantID, a.OrganizationID+"-foreign") || contract.ExpectedLiveMode != a.ExpectedLiveMode || contract.DocumentSHA256 != fuzzProfileHash(raw) {
			t.Fatal("profile escaped its explicit activation")
		}
		provider, account, connection, required := contract.PaymentBinding()
		if !required || provider != a.ProviderCode || account != a.ProviderAccountRef || connection != a.ProviderConnectionID {
			t.Fatal("profile lost its payment identity binding")
		}
		for _, mutate := range []func(*HandoverProfileActivation){
			func(p *HandoverProfileActivation) { p.Enabled = false },
			func(p *HandoverProfileActivation) { p.ExpectedLiveMode = !p.ExpectedLiveMode },
			func(p *HandoverProfileActivation) { p.TenantID += "-foreign" },
			func(p *HandoverProfileActivation) { p.OrganizationID += "-foreign" },
			func(p *HandoverProfileActivation) { p.ProviderAccountRef += "-foreign" },
			func(p *HandoverProfileActivation) { p.ProviderConnectionID += "-foreign" },
			func(p *HandoverProfileActivation) { p.ProfileID += "-foreign" },
			func(p *HandoverProfileActivation) { p.ProfileRevision++ },
		} {
			changed := a
			mutate(&changed)
			if _, err := LoadHandoverProfile(raw, changed); err == nil {
				t.Fatal("same document was admitted under another activation")
			}
		}
		// Admission is process-local: public JSON must not forge the private
		// binding, and changes to the admitted contract must invalidate it.
		wire, err := json.Marshal(contract)
		if err != nil {
			t.Fatal(err)
		}
		var forged HandoverReleaseContract
		if json.Unmarshal(wire, &forged) != nil || forged.Valid() {
			t.Fatal("serialized contract bypasses profile admission")
		}
		for _, mutate := range []func(*HandoverReleaseContract){
			func(p *HandoverReleaseContract) { p.ExpectedLiveMode = !p.ExpectedLiveMode },
			func(p *HandoverReleaseContract) { p.DocumentSHA256 = strings.Repeat("f", 64) },
			func(p *HandoverReleaseContract) { p.MaximumObservationAge += time.Second },
			func(p *HandoverReleaseContract) { p.Scope = "LOCAL_FIXTURES" },
		} {
			changed := contract
			mutate(&changed)
			if changed.Valid() {
				t.Fatal("public contract mutation retained admission")
			}
		}
		changedBytes := append(append([]byte(nil), raw...), ' ')
		if _, err := LoadHandoverProfile(changedBytes, a); err == nil {
			t.Fatal("profile lock accepted different bytes")
		}
		// A valid semantic roundtrip requires an explicitly updated byte lock.
		var document HandoverProfileDocument
		if json.Unmarshal(raw, &document) != nil {
			t.Fatal("accepted profile is not a decodable document")
		}
		normalized, err := json.Marshal(document)
		if err != nil {
			t.Fatal(err)
		}
		a.DocumentSHA256 = fuzzProfileHash(normalized)
		restored, err := LoadHandoverProfile(normalized, a)
		if err != nil || !restored.Valid() || !restored.AllowsScope(a.TenantID, a.OrganizationID) || restored.ID != contract.ID || restored.ExpectedLiveMode != contract.ExpectedLiveMode {
			t.Fatal("explicitly relocked semantic roundtrip changes admission")
		}
		// Duplicate properties must be rejected even when their bytes are
		// locked: two reviewers must not see different policy declarations.
		duplicate := append([]byte(`{"schema":"duplicate",`), normalized[1:]...)
		a.DocumentSHA256 = fuzzProfileHash(duplicate)
		if _, err := LoadHandoverProfile(duplicate, a); err == nil {
			t.Fatal("duplicate property accepted under an exact byte lock")
		}
	})
}
````

### FILE: `internal/franchisejourney/handover_profile_test.go`

```yaml
block_id: "GO-INITIAL-HANDOVER-API:file18:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "75ea6db0c29b095fb5daa810eb7595835f05e46a4fba46bc1c6b01cda91de0ce"
variables: []
secrets_allowed: false
```

````go
package franchisejourney

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func profileFixture(mode bool) (HandoverProfileDocument, HandoverProfileActivation) {
	d := HandoverProfileDocument{Schema: HandoverProfileSchema, ProfileID: "franchise-reference", Revision: 1, Algorithm: HandoverSupportedAlgorithm, AlgorithmRevision: 1, Scope: "MATERIALIZED_PROFILE", TenantID: "018f4d4a-7b36-7a21-8d10-000000000001", OrganizationID: "store", ProviderCode: "stripe", ProviderAccountRef: "acct_fixture", ProviderConnectionID: "checkout", ExpectedLiveMode: &mode, MaximumObservationAgeSeconds: 300, Options: SupportedHandoverProfileOptions(), AuthorityReference: "docs/initial-handover-reference.md", DecisionReference: "PROJECT_HANDOVER_POLICY_DECISION.md"}
	a := HandoverProfileActivation{Enabled: true, ProfileID: d.ProfileID, ProfileRevision: d.Revision, TenantID: d.TenantID, OrganizationID: d.OrganizationID, ProviderCode: d.ProviderCode, ProviderAccountRef: d.ProviderAccountRef, ProviderConnectionID: d.ProviderConnectionID, ExpectedLiveMode: mode}
	return d, a
}
func lockProfile(t *testing.T, d HandoverProfileDocument, a HandoverProfileActivation) ([]byte, HandoverProfileActivation) {
	t.Helper()
	raw, err := json.Marshal(d)
	if err != nil {
		t.Fatal(err)
	}
	hash := sha256.Sum256(raw)
	a.DocumentSHA256 = hex.EncodeToString(hash[:])
	return raw, a
}
func TestHandoverProfileBindsAlgorithmRevisionModeAndScope(t *testing.T) {
	for _, mode := range []bool{false, true} {
		d, a := profileFixture(mode)
		raw, a := lockProfile(t, d, a)
		contract, err := LoadHandoverProfile(raw, a)
		if err != nil || !contract.Valid() || contract.ExpectedLiveMode != mode || !contract.AllowsScope(d.TenantID, "store") || contract.AllowsScope(d.TenantID, "foreign") {
			t.Fatal("profile", contract, err)
		}
		path := filepath.Join(t.TempDir(), "profile.json")
		if err = os.WriteFile(path, raw, 0600); err != nil {
			t.Fatal(err)
		}
		loaded, err := LoadHandoverProfileFile(path, a)
		if err != nil || loaded.DocumentSHA256 != contract.DocumentSHA256 {
			t.Fatal("file loader", err)
		}
		forged := contract
		forged.ExpectedLiveMode = !mode
		if forged.Valid() {
			t.Fatal("mode changed after admission")
		}
		repo := &preparationRepoProbe{}
		service, err := NewHandoverPreparationService(repo, preparationIDs{}, contract)
		if err != nil {
			t.Fatal(err)
		}
		command := PrepareHandoverCommand{OrganizationID: "foreign", OrderID: "order", OrderLineID: "line", PaymentAttemptID: "payment", ObservationSHA256: strings.Repeat("a", 64), IdempotencyKey: "profile-key"}
		if _, _, err = service.Prepare(context.Background(), d.TenantID, "operator", command); !errors.Is(err, ErrReleaseConditioned) || repo.calls != 0 {
			t.Fatal("scope bypass", err)
		}
		// JSON cannot manufacture the loader's private immutable admission binding.
		encoded, _ := json.Marshal(contract)
		var roundtrip HandoverReleaseContract
		json.Unmarshal(encoded, &roundtrip)
		if roundtrip.Valid() {
			t.Fatal("unadmitted reconstructed contract accepted")
		}
	}
}
func TestHandoverProfileRejectsAlteredMissingAndUnsupportedPolicy(t *testing.T) {
	edits := map[string]func(*HandoverProfileDocument, *HandoverProfileActivation){
		"provider-binding":   func(_ *HandoverProfileDocument, a *HandoverProfileActivation) { a.ProviderCode = "mercadopago" },
		"account-binding":    func(_ *HandoverProfileDocument, a *HandoverProfileActivation) { a.ProviderAccountRef = "acct_other" },
		"connection-binding": func(_ *HandoverProfileDocument, a *HandoverProfileActivation) { a.ProviderConnectionID = "other" },
		"missing-provider":   func(d *HandoverProfileDocument, _ *HandoverProfileActivation) { d.ProviderCode = "" },
		"disabled":           func(_ *HandoverProfileDocument, a *HandoverProfileActivation) { a.Enabled = false },
		"schema":             func(d *HandoverProfileDocument, _ *HandoverProfileActivation) { d.Schema = "v2" },
		"algorithm":          func(d *HandoverProfileDocument, _ *HandoverProfileActivation) { d.Algorithm = "arbitrary-policy" },
		"algorithm-revision": func(d *HandoverProfileDocument, _ *HandoverProfileActivation) { d.AlgorithmRevision = 2 },
		"profile-revision":   func(d *HandoverProfileDocument, _ *HandoverProfileActivation) { d.Revision = 2 },
		"scope":              func(d *HandoverProfileDocument, _ *HandoverProfileActivation) { d.Scope = "PRODUCTION_PROVEN" },
		"mode-absent":        func(d *HandoverProfileDocument, _ *HandoverProfileActivation) { d.ExpectedLiveMode = nil },
		"mode-binding":       func(_ *HandoverProfileDocument, a *HandoverProfileActivation) { a.ExpectedLiveMode = true },
		"quantity":           func(d *HandoverProfileDocument, _ *HandoverProfileActivation) { d.Options.Quantity = 2 },
		"credit":             func(d *HandoverProfileDocument, _ *HandoverProfileActivation) { d.Options.PaymentCoverage = "PARTIAL" },
		"shipping": func(d *HandoverProfileDocument, _ *HandoverProfileActivation) {
			d.Options.ReleaseEffect = "POST_SHIPMENT"
		},
		"missing-decision": func(d *HandoverProfileDocument, _ *HandoverProfileActivation) { d.DecisionReference = "" },
		"missing-tenant":   func(d *HandoverProfileDocument, _ *HandoverProfileActivation) { d.TenantID = "" },
		"tenant-binding": func(_ *HandoverProfileDocument, a *HandoverProfileActivation) {
			a.TenantID = "018f4d4a-7b36-7a21-8d10-000000000002"
		},
		"stale-age": func(d *HandoverProfileDocument, _ *HandoverProfileActivation) { d.MaximumObservationAgeSeconds = 901 },
	}
	for name, edit := range edits {
		t.Run(name, func(t *testing.T) {
			d, a := profileFixture(false)
			edit(&d, &a)
			raw, a := lockProfile(t, d, a)
			if _, err := LoadHandoverProfile(raw, a); !errors.Is(err, ErrReleaseConditioned) {
				t.Fatal("unsupported policy admitted", err)
			}
		})
	}
	d, a := profileFixture(false)
	raw, a := lockProfile(t, d, a)
	if _, err := LoadHandoverProfile(append(raw, ' '), a); err == nil {
		t.Fatal("byte lock ignored")
	}
	if _, err := LoadHandoverProfile(nil, a); err == nil {
		t.Fatal("absent policy admitted")
	}
	if _, err := LoadHandoverProfileFile(filepath.Join(t.TempDir(), "absent"), a); err == nil {
		t.Fatal("absent policy file admitted")
	}
	for _, changed := range [][]byte{[]byte(`{"policy":true}`), []byte(strings.Replace(string(raw), `"revision":1`, `"revision":1,"revision":1`, 1)), append(raw, []byte(`{}`)...)} {
		hash := sha256.Sum256(changed)
		b := a
		b.DocumentSHA256 = hex.EncodeToString(hash[:])
		if _, err := LoadHandoverProfile(changed, b); err == nil {
			t.Fatal("unknown/ambiguous JSON admitted")
		}
	}
}
````

### FILE: `internal/platform/httpapi/commercial_release.go`

```yaml
block_id: "GO-INITIAL-HANDOVER-API:file19:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "fa821602e9794279a35969ca0e6960196b78027809ca3a5fd7e5c99ff6127dd0"
variables: []
secrets_allowed: false
```

````go
package httpapi

// AUTHORED HTTP composition; authenticated tenant/actor and organization
// permission bind every command and recovery. Money/stock are derived by SQL.
import (
	"context"
	"net/http"

	"elite.local/enterprise/internal/franchisejourney"
	"elite.local/enterprise/internal/platform/identity"
)

type CommercialReleaseService interface {
	CommitCommercialRelease(context.Context, string, string, franchisejourney.CommitCommercialReleaseCommand) (franchisejourney.CommercialReleaseReceipt, bool, error)
	CommercialReleaseResult(context.Context, string, string, string, string) (franchisejourney.CommercialReleaseReceipt, error)
	ValidateCommercialRelease(context.Context, string, string, string) (franchisejourney.CurrentCommercialRelease, error)
}
type CommercialReleaseModule struct{ Service CommercialReleaseService }

func (m CommercialReleaseModule) Register(mux *http.ServeMux, verifier identity.Verifier) {
	if m.Service == nil {
		return
	}
	principal := franchiseJourneyAPI{verifier: verifier}
	mux.HandleFunc("POST /v1/franchise/handovers/{id}/commercial-release", func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			OrganizationID    string `json:"organization_id"`
			ObservationSHA256 string `json:"observation_sha256"`
		}
		if !decodeStrict(w, r, &input) {
			return
		}
		actor, ok := principal.protected(w, r, "handover:manage", input.OrganizationID)
		if !ok {
			return
		}
		v, replay, err := m.Service.CommitCommercialRelease(r.Context(), actor.TenantID, actor.Subject, franchisejourney.CommitCommercialReleaseCommand{OrganizationID: input.OrganizationID, HandoverID: r.PathValue("id"), ObservationSHA256: input.ObservationSHA256, IdempotencyKey: r.Header.Get("Idempotency-Key")})
		if initialHandoverError(w, err) {
			return
		}
		w.Header().Set("Cache-Control", "no-store")
		if replay {
			w.Header().Set("Idempotency-Replayed", "true")
		}
		writeJSON(w, http.StatusCreated, v)
	})
	mux.HandleFunc("GET /v1/franchise/handovers/{id}/commercial-release-result", func(w http.ResponseWriter, r *http.Request) {
		org := r.URL.Query().Get("organization_id")
		actor, ok := principal.protected(w, r, "handover:manage", org)
		if !ok {
			return
		}
		v, err := m.Service.CommercialReleaseResult(r.Context(), actor.TenantID, org, r.PathValue("id"), r.Header.Get("Idempotency-Key"))
		if initialHandoverError(w, err) {
			return
		}
		w.Header().Set("Cache-Control", "no-store")
		writeJSON(w, http.StatusOK, v)
	})
	mux.HandleFunc("GET /v1/franchise/handovers/{id}/commercial-release-current", func(w http.ResponseWriter, r *http.Request) {
		org := r.URL.Query().Get("organization_id")
		actor, ok := principal.protected(w, r, "handover:manage", org)
		if !ok {
			return
		}
		v, err := m.Service.ValidateCommercialRelease(r.Context(), actor.TenantID, org, r.PathValue("id"))
		if initialHandoverError(w, err) {
			return
		}
		w.Header().Set("Cache-Control", "no-store")
		writeJSON(w, http.StatusOK, v)
	})
}
````

### FILE: `internal/platform/httpapi/commercial_release_test.go`

```yaml
block_id: "GO-INITIAL-HANDOVER-API:file20:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "c538bf9c6f9be6bbca4340663417ab6247d7c8a6bb7b014f4fead0a2834d56a8"
variables: []
secrets_allowed: false
```

````go
package httpapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"elite.local/enterprise/internal/franchisejourney"
)

type commercialReleaseSpy struct {
	calls         int
	tenant, actor string
	command       franchisejourney.CommitCommercialReleaseCommand
}

func (s *commercialReleaseSpy) CommitCommercialRelease(_ context.Context, tenant, actor string, c franchisejourney.CommitCommercialReleaseCommand) (franchisejourney.CommercialReleaseReceipt, bool, error) {
	s.calls++
	s.tenant = tenant
	s.actor = actor
	s.command = c
	return franchisejourney.CommercialReleaseReceipt{Effect: franchisejourney.CommercialReleaseEffect}, true, nil
}
func (s *commercialReleaseSpy) CommercialReleaseResult(_ context.Context, tenant, org, handover, key string) (franchisejourney.CommercialReleaseReceipt, error) {
	s.calls++
	s.tenant = tenant
	return franchisejourney.CommercialReleaseReceipt{}, nil
}
func (s *commercialReleaseSpy) ValidateCommercialRelease(_ context.Context, tenant, org, handover string) (franchisejourney.CurrentCommercialRelease, error) {
	s.calls++
	s.tenant = tenant
	return franchisejourney.CurrentCommercialRelease{Current: false}, nil
}
func TestCommercialReleaseHTTPAuthorizationAndHistoricalReplay(t *testing.T) {
	valid := `{"organization_id":"store","observation_sha256":"` + strings.Repeat("a", 64) + `"}`
	for _, name := range []string{"authorized", "unauthenticated", "permission", "organization", "forged-actor", "forged-money", "forged-order"} {
		t.Run(name, func(t *testing.T) {
			p := initialHandoverPrincipal()
			auth := "Bearer fixture"
			body := valid
			want := 201
			switch name {
			case "unauthenticated":
				auth = ""
				want = 401
			case "permission":
				p.Permissions = map[string]struct{}{}
				want = 403
			case "organization":
				p.Organizations = map[string]struct{}{}
				want = 403
			case "forged-actor", "forged-money", "forged-order":
				field := map[string]string{"forged-actor": "actor_subject", "forged-money": "amount_minor_units", "forged-order": "order_id"}[name]
				body = strings.TrimSuffix(valid, "}") + `,"` + field + `":"forged"}`
				want = 400
			}
			s := &commercialReleaseSpy{}
			mux := http.NewServeMux()
			CommercialReleaseModule{Service: s}.Register(mux, journeyVerifier{principal: p})
			r := httptest.NewRequest("POST", "/v1/franchise/handovers/server-handover/commercial-release", strings.NewReader(body))
			r.Header.Set("Authorization", auth)
			r.Header.Set("Idempotency-Key", "release-key")
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, r)
			if w.Code != want {
				t.Fatalf("got %d want %d", w.Code, want)
			}
			if want != 201 {
				if s.calls != 0 {
					t.Fatal("unauthorized owner call")
				}
				return
			}
			if s.tenant != "verified-tenant" || s.actor != "verified-actor" || s.command.HandoverID != "server-handover" || s.command.IdempotencyKey != "release-key" || w.Header().Get("Idempotency-Replayed") != "true" || strings.Contains(w.Body.String(), `"current":true`) {
				t.Fatal("scope or replay semantics")
			}
			for _, route := range []string{"commercial-release-result", "commercial-release-current"} {
				r = httptest.NewRequest("GET", "/v1/franchise/handovers/server-handover/"+route+"?organization_id=store", nil)
				r.Header.Set("Authorization", auth)
				r.Header.Set("Idempotency-Key", "release-key")
				w = httptest.NewRecorder()
				mux.ServeHTTP(w, r)
				if w.Code != 200 || w.Header().Get("Cache-Control") != "no-store" || strings.Contains(w.Body.String(), `"current":true`) {
					t.Fatal("recovery grants current authority", w.Code)
				}
			}
		})
	}
}
````

### FILE: `internal/platform/httpapi/handover_preparation.go`

```yaml
block_id: "GO-INITIAL-HANDOVER-API:file21:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "93d4e35f3d69a1e5b696a9bcc17f21f46344aa360e5abac85fd44257b488ea45"
variables: []
secrets_allowed: false
```

````go
package httpapi

// AUTHORED HTTP glue. The module is mounted only with an explicitly configured
// service; claims, tenant, actor and organization are verified server-side.
import (
	"context"
	"errors"
	"net/http"

	"elite.local/enterprise/internal/franchisejourney"
	"elite.local/enterprise/internal/platform/identity"
)

type InitialHandoverService interface {
	Prepare(context.Context, string, string, franchisejourney.PrepareHandoverCommand) (franchisejourney.HandoverPreparation, bool, error)
	Result(context.Context, string, string, string, string) (franchisejourney.HandoverPreparation, error)
	EvaluateRelease(context.Context, string, string, string, string) (franchisejourney.HandoverReferenceRelease, error)
}
type InitialHandoverModule struct{ Service InitialHandoverService }

func (m InitialHandoverModule) Register(mux *http.ServeMux, verifier identity.Verifier) {
	if m.Service == nil {
		return
	}
	principal := franchiseJourneyAPI{verifier: verifier}
	mux.HandleFunc("POST /v1/franchise/orders/{id}/handover", func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			OrganizationID    string `json:"organization_id"`
			OrderLineID       string `json:"order_line_id"`
			PaymentAttemptID  string `json:"payment_attempt_id"`
			FundingReceiptID  string `json:"funding_receipt_id"`
			ObservationSHA256 string `json:"observation_sha256"`
		}
		if !decodeStrict(w, r, &input) {
			return
		}
		actor, ok := principal.protected(w, r, "handover:manage", input.OrganizationID)
		if !ok {
			return
		}
		value, replay, err := m.Service.Prepare(r.Context(), actor.TenantID, actor.Subject, franchisejourney.PrepareHandoverCommand{OrganizationID: input.OrganizationID, OrderID: r.PathValue("id"), OrderLineID: input.OrderLineID, PaymentAttemptID: input.PaymentAttemptID, FundingReceiptID: input.FundingReceiptID, ObservationSHA256: input.ObservationSHA256, IdempotencyKey: r.Header.Get("Idempotency-Key")})
		if initialHandoverError(w, err) {
			return
		}
		if replay {
			w.Header().Set("Idempotency-Replayed", "true")
		}
		writeJSON(w, http.StatusCreated, value)
	})
	mux.HandleFunc("GET /v1/franchise/orders/{id}/handover-result", func(w http.ResponseWriter, r *http.Request) {
		organization := r.URL.Query().Get("organization_id")
		actor, ok := principal.protected(w, r, "handover:manage", organization)
		if !ok {
			return
		}
		value, err := m.Service.Result(r.Context(), actor.TenantID, organization, r.PathValue("id"), r.Header.Get("Idempotency-Key"))
		if initialHandoverError(w, err) {
			return
		}
		writeJSON(w, http.StatusOK, value)
	})
	mux.HandleFunc("GET /v1/franchise/handovers/{id}/release-check", func(w http.ResponseWriter, r *http.Request) {
		organization := r.URL.Query().Get("organization_id")
		actor, ok := principal.protected(w, r, "handover:manage", organization)
		if !ok {
			return
		}
		value, err := m.Service.EvaluateRelease(r.Context(), actor.TenantID, organization, r.PathValue("id"), r.URL.Query().Get("observation_sha256"))
		if initialHandoverError(w, err) {
			return
		}
		writeJSON(w, http.StatusOK, value)
	})
}
func initialHandoverError(w http.ResponseWriter, err error) bool {
	switch {
	case err == nil:
		return false
	case errors.Is(err, franchisejourney.ErrInvalid):
		writeProblem(w, 400, "INVALID_HANDOVER_PREPARATION", "handover request does not match its contract")
	case errors.Is(err, franchisejourney.ErrNotFound):
		writeProblem(w, 404, "HANDOVER_PREPARATION_NOT_FOUND", "no completed preparation exists for this scope and key")
	case errors.Is(err, franchisejourney.ErrConflict):
		writeProblem(w, 409, "HANDOVER_NOT_ELIGIBLE", "current order, stock, payment or acceptance evidence does not permit this operation")
	case errors.Is(err, franchisejourney.ErrReleaseConditioned):
		writeProblem(w, 503, "HANDOVER_CONTRACT_REQUIRED", "an admitted release contract must be configured")
	default:
		writeProblem(w, 500, "HANDOVER_PREPARATION_FAILED", "handover operation could not complete")
	}
	return true
}
````

### FILE: `internal/platform/httpapi/handover_preparation_test.go`

```yaml
block_id: "GO-INITIAL-HANDOVER-API:file22:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "461e9e3294c21972289ec59b0a74ddcb5cfb0419af8b278c2b20164282108435"
variables: []
secrets_allowed: false
```

````go
package httpapi

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"elite.local/enterprise/internal/franchisejourney"
	"elite.local/enterprise/internal/platform/identity"
)

type initialHandoverSpy struct {
	calls                     int
	tenant, actor, order, key string
	replay                    bool
	err                       error
}

func (s *initialHandoverSpy) Prepare(_ context.Context, tenant, actor string, c franchisejourney.PrepareHandoverCommand) (franchisejourney.HandoverPreparation, bool, error) {
	s.calls++
	s.tenant = tenant
	s.actor = actor
	s.order = c.OrderID
	s.key = c.IdempotencyKey
	return franchisejourney.HandoverPreparation{}, s.replay, s.err
}
func (s *initialHandoverSpy) Result(_ context.Context, tenant, organization, order, key string) (franchisejourney.HandoverPreparation, error) {
	s.calls++
	s.tenant = tenant
	s.order = order
	s.key = key
	return franchisejourney.HandoverPreparation{}, s.err
}
func (s *initialHandoverSpy) EvaluateRelease(_ context.Context, tenant, organization, handover, sha string) (franchisejourney.HandoverReferenceRelease, error) {
	s.calls++
	s.tenant = tenant
	return franchisejourney.HandoverReferenceRelease{Scope: "LOCAL_FIXTURES", Eligible: true}, s.err
}
func initialHandoverPrincipal() identity.Principal {
	return identity.Principal{TenantID: "verified-tenant", Subject: "verified-actor", Permissions: map[string]struct{}{"handover:manage": {}}, Organizations: map[string]struct{}{"store": {}}}
}
func TestInitialHandoverHTTPBindsPrincipalAndRecovery(t *testing.T) {
	spy := &initialHandoverSpy{replay: true}
	mux := http.NewServeMux()
	InitialHandoverModule{Service: spy}.Register(mux, journeyVerifier{principal: initialHandoverPrincipal()})
	request := httptest.NewRequest("POST", "/v1/franchise/orders/server-order/handover", strings.NewReader(`{"organization_id":"store","order_line_id":"line","payment_attempt_id":"payment","observation_sha256":"`+strings.Repeat("a", 64)+`"}`))
	request.Header.Set("Authorization", "Bearer fixture")
	request.Header.Set("Idempotency-Key", "recovery-key")
	recorder := httptest.NewRecorder()
	mux.ServeHTTP(recorder, request)
	if recorder.Code != 201 || recorder.Header().Get("Idempotency-Replayed") != "true" || spy.tenant != "verified-tenant" || spy.actor != "verified-actor" || spy.order != "server-order" || spy.key != "recovery-key" {
		t.Fatalf("binding failed: %d %+v", recorder.Code, spy)
	}
	request = httptest.NewRequest("GET", "/v1/franchise/orders/server-order/handover-result?organization_id=store", nil)
	request.Header.Set("Authorization", "Bearer fixture")
	request.Header.Set("Idempotency-Key", "recovery-key")
	recorder = httptest.NewRecorder()
	mux.ServeHTTP(recorder, request)
	if recorder.Code != 200 || spy.calls != 2 || spy.key != "recovery-key" {
		t.Fatalf("GET recovery failed: %d %+v", recorder.Code, spy)
	}
}
func TestInitialHandoverHTTPRejectsForgedInputsAndUnauthorizedScope(t *testing.T) {
	valid := `{"organization_id":"store","order_line_id":"line","payment_attempt_id":"payment","observation_sha256":"` + strings.Repeat("a", 64) + `"}`
	for _, name := range []string{"unauthenticated", "wrong-permission", "wrong-organization", "forged-customer", "forged-order", "forged-amount"} {
		t.Run(name, func(t *testing.T) {
			principal := initialHandoverPrincipal()
			body := valid
			authorization := "Bearer fixture"
			want := 400
			switch name {
			case "unauthenticated":
				authorization = ""
				want = 401
			case "wrong-permission":
				principal.Permissions = map[string]struct{}{}
				want = 403
			case "wrong-organization":
				principal.Organizations = map[string]struct{}{"other": {}}
				want = 403
			default:
				field := map[string]string{"forged-customer": "customer_subject", "forged-order": "order_id", "forged-amount": "amount_minor_units"}[name]
				body = strings.TrimSuffix(valid, "}") + `,"` + field + `":"forged"}`
			}
			spy := &initialHandoverSpy{}
			mux := http.NewServeMux()
			InitialHandoverModule{Service: spy}.Register(mux, journeyVerifier{principal: principal})
			request := httptest.NewRequest("POST", "/v1/franchise/orders/order/handover", strings.NewReader(body))
			request.Header.Set("Authorization", authorization)
			recorder := httptest.NewRecorder()
			mux.ServeHTTP(recorder, request)
			if recorder.Code != want || spy.calls != 0 {
				t.Fatalf("rejected request reached owner: status%d calls%d", recorder.Code, spy.calls)
			}
		})
	}
}
func TestInitialHandoverHTTPUnconfiguredModuleAndErrors(t *testing.T) {
	mux := http.NewServeMux()
	InitialHandoverModule{}.Register(mux, journeyVerifier{})
	recorder := httptest.NewRecorder()
	mux.ServeHTTP(recorder, httptest.NewRequest("GET", "/v1/franchise/orders/order/handover-result", nil))
	if recorder.Code != 404 {
		t.Fatal("unconfigured route mounted")
	}
	for _, item := range []struct {
		err    error
		status int
	}{{franchisejourney.ErrConflict, 409}, {franchisejourney.ErrNotFound, 404}, {franchisejourney.ErrReleaseConditioned, 503}, {errors.New("private database information"), 500}} {
		spy := &initialHandoverSpy{err: item.err}
		mux = http.NewServeMux()
		InitialHandoverModule{Service: spy}.Register(mux, journeyVerifier{principal: initialHandoverPrincipal()})
		request := httptest.NewRequest("GET", "/v1/franchise/handovers/handover/release-check?organization_id=store&observation_sha256="+strings.Repeat("a", 64), nil)
		request.Header.Set("Authorization", "Bearer fixture")
		recorder = httptest.NewRecorder()
		mux.ServeHTTP(recorder, request)
		if recorder.Code != item.status || strings.Contains(recorder.Body.String(), "private database information") {
			t.Fatalf("unsafe mapping %d %s", recorder.Code, recorder.Body.String())
		}
	}
}
````

### FILE: `internal/platform/postgres/commercial_release.go`

```yaml
block_id: "GO-INITIAL-HANDOVER-API:file23:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "3ccbbd76289b856f8195cc5b5fe2cdd40df48daa37076ade5dad7d31b088158c"
variables: []
secrets_allowed: false
```

````go
package postgres

// AUTHORED persistence glue over the existing handover/payment owners. The
// receipt records an atomic commercial checkpoint, never a physical shipment.
import (
	"context"
	"errors"
	"time"

	"elite.local/enterprise/internal/franchisejourney"
	"github.com/jackc/pgx/v5"
)

const commercialReleaseColumns = `r.release_id,r.organization_id,r.handover_id,r.order_id,coalesce(r.payment_attempt_id,''),coalesce(r.funding_receipt_id,''),r.observation_sha256_hex,r.observation_generation,r.handover_version,r.acceptance_sha256_hex,r.checklist_id,r.checklist_version,r.contract_id,r.contract_sha256_hex,r.effect,r.released_by_subject,r.recorded_at,r.valid_until`

func scanCommercialRelease(row pgx.Row) (franchisejourney.CommercialReleaseReceipt, error) {
	var r franchisejourney.CommercialReleaseReceipt
	err := row.Scan(&r.ID, &r.OrganizationID, &r.HandoverID, &r.OrderID, &r.PaymentAttemptID, &r.FundingReceiptID, &r.ObservationSHA256, &r.ObservationGeneration, &r.HandoverVersion, &r.AcceptanceSHA256, &r.ChecklistID, &r.ChecklistVersion, &r.ContractID, &r.ContractSHA256, &r.Effect, &r.ReleasedBy, &r.RecordedAt, &r.ValidUntil)
	if errors.Is(err, pgx.ErrNoRows) {
		err = franchisejourney.ErrNotFound
	}
	return r, err
}

func readCommercialRelease(ctx context.Context, q handoverPreparationReader, tenant, organization, handover, key string) (franchisejourney.CommercialReleaseReceipt, error) {
	return scanCommercialRelease(q.QueryRow(ctx, `select `+commercialReleaseColumns+` from sales.commercial_release_receipt r
 join platform.idempotency_record i on i.tenant_id=r.tenant_id and i.resource_id=r.release_id
 where r.tenant_id=$1 and r.organization_id=$2 and r.handover_id=$3 and i.scope='commercial-release' and i.idempotency_key=$4
 and i.status='completed' and i.resource_type='commercial-release' and i.response_code=201 and i.response_body->>'release_id'=r.release_id`, tenant, organization, handover, key))
}

func (r *FranchiseJourney) InitialCommercialReleaseResult(ctx context.Context, tenant, organization, handover, key string) (franchisejourney.CommercialReleaseReceipt, error) {
	return readCommercialRelease(ctx, r.pool, tenant, organization, handover, key)
}

// Derive the order/payment from immutable preparation. The caller cannot supply
// a different customer, stock unit, amount, payment reference or accepted state.
func lockCommercialRelease(ctx context.Context, tx pgx.Tx, tenant, organization, handover, evidence string, p franchisejourney.HandoverReleaseContract) (franchisejourney.CommercialReleaseReceipt, error) {
	var v franchisejourney.CommercialReleaseReceipt
	if !p.AllowsCommercialRelease(tenant, organization) {
		return v, franchisejourney.ErrReleaseConditioned
	}
	c := franchisejourney.PrepareHandoverCommand{OrganizationID: organization, ObservationSHA256: evidence}
	var id, sha, scope string
	var age int64
	err := tx.QueryRow(ctx, `select order_id,order_line_id,coalesce(payment_attempt_id,''),coalesce(funding_receipt_id,''),contract_id,contract_sha256_hex,contract_scope,maximum_observation_age_ns from sales.delivery_handover_preparation where tenant_id=$1 and organization_id=$2 and handover_id=$3`, tenant, organization, handover).Scan(&c.OrderID, &c.OrderLineID, &c.PaymentAttemptID, &c.FundingReceiptID, &id, &sha, &scope, &age)
	if errors.Is(err, pgx.ErrNoRows) {
		return v, franchisejourney.ErrConflict
	}
	if err != nil {
		return v, err
	}
	if id != p.ID || sha != p.DocumentSHA256 || scope != p.Scope || age != int64(p.MaximumObservationAge) {
		return v, franchisejourney.ErrConflict
	}
	facts, err := lockInitialHandoverScope(ctx, tx, tenant, c, p)
	if err != nil {
		return v, err
	}
	// Acceptance takes handover -> order. NOWAIT avoids inversion; a caller may
	// retry after the other transaction has completed without changing its key.
	err = tx.QueryRow(ctx, `select version,acceptance_evidence_sha256_hex,checklist_id,checklist_version from sales.delivery_handover
 where tenant_id=$1 and organization_id=$2 and handover_id=$3 and order_id=$4 and customer_principal_id=$5 and stock_unit_id=$6
 and state='accepted' and customer_accepted_at is not null and acceptance_evidence_sha256_hex is not null and checklist_completed_at is not null for share nowait`, tenant, organization, handover, c.OrderID, facts.customer, facts.stock).Scan(&v.HandoverVersion, &v.AcceptanceSHA256, &v.ChecklistID, &v.ChecklistVersion)
	if errors.Is(err, pgx.ErrNoRows) {
		return v, franchisejourney.ErrConflict
	}
	if err != nil {
		return v, err
	}
	observed := facts.evidenceAt
	v.ObservationGeneration = facts.generation
	var expires *time.Time
	err = tx.QueryRow(ctx, `select expires_at,clock_timestamp() from inventory.serial_reservation where tenant_id=$1 and reservation_id=$2`, tenant, facts.reservation).Scan(&expires, &v.RecordedAt)
	if err != nil {
		return v, err
	}
	v.ValidUntil = observed.Add(p.MaximumObservationAge)
	if expires != nil && expires.Before(v.ValidUntil) {
		v.ValidUntil = *expires
	}
	// Time is sampled after every relevant lock. A snapshot that became stale
	// while waiting must never be committed as a fresh release.
	if observed.After(v.RecordedAt) || !v.ValidUntil.After(v.RecordedAt) {
		return v, franchisejourney.ErrConflict
	}
	v.OrganizationID = organization
	v.HandoverID = handover
	v.OrderID = c.OrderID
	v.PaymentAttemptID = c.PaymentAttemptID
	v.FundingReceiptID = c.FundingReceiptID
	v.ObservationSHA256 = evidence
	v.ContractID = p.ID
	v.ContractSHA256 = p.DocumentSHA256
	v.Effect = franchisejourney.CommercialReleaseEffect
	return v, nil
}

func (r *FranchiseJourney) CommitInitialCommercialRelease(ctx context.Context, tenant, actor, id, event string, c franchisejourney.CommitCommercialReleaseCommand, p franchisejourney.HandoverReleaseContract, hash string) (franchisejourney.CommercialReleaseReceipt, bool, error) {
	var empty franchisejourney.CommercialReleaseReceipt
	if !p.AllowsCommercialRelease(tenant, c.OrganizationID) || tenant == "" || actor == "" || id == "" || event == "" || len(hash) != 64 {
		return empty, false, franchisejourney.ErrReleaseConditioned
	}
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return empty, false, err
	}
	defer tx.Rollback(ctx)
	insert, err := tx.Exec(ctx, `insert into platform.idempotency_record(tenant_id,scope,idempotency_key,request_sha256_hex,status,locked_until,expires_at)
 values($1,'commercial-release',$2,$3,'processing',clock_timestamp()+interval '30 seconds',clock_timestamp()+interval '24 hours') on conflict do nothing`, tenant, c.IdempotencyKey, hash)
	if err != nil {
		return empty, false, err
	}
	if insert.RowsAffected() == 0 {
		var saved string
		if err = tx.QueryRow(ctx, `select request_sha256_hex from platform.idempotency_record where tenant_id=$1 and scope='commercial-release' and idempotency_key=$2`, tenant, c.IdempotencyKey).Scan(&saved); err != nil || saved != hash {
			return empty, false, franchisejourney.ErrConflict
		}
		v, e := readCommercialRelease(ctx, tx, tenant, c.OrganizationID, c.HandoverID, c.IdempotencyKey)
		if e != nil {
			return empty, false, franchisejourney.ErrConflict
		}
		if err = tx.Commit(ctx); err != nil {
			return empty, false, err
		}
		return v, true, nil
	}
	v, err := lockCommercialRelease(ctx, tx, tenant, c.OrganizationID, c.HandoverID, c.ObservationSHA256, p)
	if err != nil {
		return empty, false, err
	}
	var exists bool
	if err = tx.QueryRow(ctx, `select exists(select 1 from sales.commercial_release_receipt where tenant_id=$1 and organization_id=$2 and order_id=$3)`, tenant, c.OrganizationID, v.OrderID).Scan(&exists); err != nil {
		return empty, false, err
	}
	if exists {
		return empty, false, franchisejourney.ErrConflict
	}
	v.ID = id
	v.ReleasedBy = actor
	_, err = tx.Exec(ctx, `insert into sales.commercial_release_receipt(tenant_id,release_id,organization_id,handover_id,order_id,payment_attempt_id,observation_sha256_hex,observation_generation,handover_version,acceptance_sha256_hex,checklist_id,checklist_version,contract_id,contract_sha256_hex,effect,released_by_subject,recorded_at,valid_until,funding_receipt_id)
 values($1,$2,$3,$4,$5,nullif($6,''),$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,nullif($19,''))`, tenant, v.ID, v.OrganizationID, v.HandoverID, v.OrderID, v.PaymentAttemptID, v.ObservationSHA256, v.ObservationGeneration, v.HandoverVersion, v.AcceptanceSHA256, v.ChecklistID, v.ChecklistVersion, v.ContractID, v.ContractSHA256, v.Effect, v.ReleasedBy, v.RecordedAt, v.ValidUntil, v.FundingReceiptID)
	if err != nil {
		return empty, false, err
	}
	if err = writeJourneyOutbox(ctx, tx, tenant, event, "commercial-release", id, 1, "commercial-release.recorded", v); err != nil {
		return empty, false, err
	}
	done, err := tx.Exec(ctx, `update platform.idempotency_record set status='completed',response_code=201,response_body=jsonb_build_object('release_id',$3::text),resource_type='commercial-release',resource_id=$3,locked_until=null where tenant_id=$1 and scope='commercial-release' and idempotency_key=$2 and status='processing'`, tenant, c.IdempotencyKey, id)
	if err != nil {
		return empty, false, err
	}
	if done.RowsAffected() != 1 {
		return empty, false, franchisejourney.ErrConflict
	}
	// An outbox/constraint wait can consume the freshness budget after the
	// initial checks. The checkpoint must still be fresh before committing.
	var commitTime time.Time
	if err = tx.QueryRow(ctx, `select clock_timestamp()`).Scan(&commitTime); err != nil {
		return empty, false, err
	}
	if !v.ValidUntil.After(commitTime) {
		return empty, false, franchisejourney.ErrConflict
	}
	if err = tx.Commit(ctx); err != nil {
		return empty, false, err
	}
	return v, false, nil
}

func (r *FranchiseJourney) ValidateInitialCommercialRelease(ctx context.Context, tenant, organization, handover string, p franchisejourney.HandoverReleaseContract) (franchisejourney.CurrentCommercialRelease, error) {
	var v franchisejourney.CurrentCommercialRelease
	if !p.AllowsCommercialRelease(tenant, organization) {
		return v, franchisejourney.ErrReleaseConditioned
	}
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return v, err
	}
	defer tx.Rollback(ctx)
	v.Receipt, err = scanCommercialRelease(tx.QueryRow(ctx, `select `+commercialReleaseColumns+` from sales.commercial_release_receipt r where r.tenant_id=$1 and r.organization_id=$2 and r.handover_id=$3`, tenant, organization, handover))
	if err != nil {
		return v, err
	}
	current, err := lockCommercialRelease(ctx, tx, tenant, organization, handover, v.Receipt.ObservationSHA256, p)
	if err != nil && !errors.Is(err, franchisejourney.ErrConflict) {
		return v, err
	}
	if e := tx.QueryRow(ctx, `select clock_timestamp()`).Scan(&v.EvaluatedAt); e != nil {
		return v, e
	}
	v.Current = err == nil && current.ObservationGeneration == v.Receipt.ObservationGeneration && current.HandoverVersion == v.Receipt.HandoverVersion && current.AcceptanceSHA256 == v.Receipt.AcceptanceSHA256 && current.ChecklistID == v.Receipt.ChecklistID && current.ChecklistVersion == v.Receipt.ChecklistVersion && current.ContractSHA256 == v.Receipt.ContractSHA256 && v.Receipt.ValidUntil.After(v.EvaluatedAt)
	if err = tx.Commit(ctx); err != nil {
		return v, err
	}
	return v, nil
}
````

### FILE: `internal/platform/postgres/commercial_release_connected_integration_test.go`

```yaml
block_id: "GO-INITIAL-HANDOVER-API:file24:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "62b6fb12c55d5d99b856f5bad78d89e27b94986ed3a89e8cbe66b49cfda63b79"
variables: []
secrets_allowed: false
```

````go
package postgres_test

// Reuses the already admitted official SDK fixture and actual callback/queue/
// reconciliation owners; adds only the durable commercial checkpoint delta.
import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"elite.local/enterprise/internal/franchisejourney"
	db "elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/platform/randomid"
)

func TestCommercialReleaseOfficialSDKConnected(t *testing.T) {
	r := newConnectedRun(t, false)
	ctx := context.Background()
	if code := r.callback(t, "evt_commercial", "checkout.session.completed", "cs_test_fixture", true); code != 200 {
		t.Fatal(code)
	}
	if n, err := r.processor.ProcessOnce(ctx); err != nil || n != 1 {
		t.Fatal("callback reconciliation", n, err)
	}
	var hash string
	if err := r.pool.QueryRow(ctx, `select evidence_sha256_hex from payment.provider_observation where tenant_id=$1`, r.tenant).Scan(&hash); err != nil {
		t.Fatal(err)
	}
	assertConnectedCommercialRelease(t, r, hash)
}

func assertConnectedCommercialRelease(t *testing.T, r *connectedRun, hash string, stored ...bool) {
	assertConnectedCommercialReleaseWithHandoverHook(t, r, hash, nil, stored...)
}

func assertConnectedCommercialReleaseWithHandoverHook(t *testing.T, r *connectedRun, hash string, hook func(*testing.T, *connectedRun, string), stored ...bool) {
	t.Helper()
	ctx := context.Background()
	policy, a := connectedCommercialProfile(t, r, len(stored) == 1 && stored[0])
	repo := db.NewFranchiseJourney(r.pool)
	ids := randomid.Generator{}
	svc, err := franchisejourney.NewHandoverPreparationService(repo, ids, policy)
	if err != nil {
		t.Fatal(err)
	}
	prepared, _, err := svc.Prepare(ctx, r.tenant, "operator", franchisejourney.PrepareHandoverCommand{OrganizationID: "store", OrderID: "order", OrderLineID: "line", PaymentAttemptID: r.payment.ID, ObservationSHA256: hash, IdempotencyKey: "commercial-integral-prepare"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = repo.PublishDeliveryChecklist(ctx, r.tenant, "operator", franchisejourney.DeliveryChecklist{ID: "commercial-checklist", OrganizationID: "store", Version: 1, Title: "Synthetic delivery", Items: []franchisejourney.ChecklistItem{{ID: "serial", Ordinal: 1, Prompt: "Read serial", ResponseType: "serial", Required: true}}}, ids.New()); err != nil {
		t.Fatal(err)
	}
	presented, err := repo.CompleteDeliveryChecklist(ctx, r.tenant, "store", "operator", prepared.Handover.ID, 1, "commercial-checklist", 1, []franchisejourney.ChecklistResponse{{ItemID: "serial", ResponseText: "SERIAL-SYNTHETIC"}}, ids.New())
	if err != nil {
		t.Fatal(err)
	}
	if _, err = repo.AcceptHandover(ctx, r.tenant, "store", "customer", prepared.Handover.ID, presented.Version, "SERIAL-SYNTHETIC", "commercial-checklist", 1, strings.Repeat("e", 64), ids.New()); err != nil {
		t.Fatal(err)
	}
	release := franchisejourney.CommitCommercialReleaseCommand{OrganizationID: "store", HandoverID: prepared.Handover.ID, ObservationSHA256: hash, IdempotencyKey: "commercial-integral-release"}
	receipt, replay, err := svc.CommitCommercialRelease(ctx, r.tenant, "operator", release)
	if err != nil || replay || receipt.ContractSHA256 != a.DocumentSHA256 {
		t.Fatal("commercial checkpoint", err)
	}
	current, err := svc.ValidateCommercialRelease(ctx, r.tenant, "store", prepared.Handover.ID)
	if err != nil || !current.Current {
		t.Fatal("checkpoint not current", err)
	}
	if hook != nil {
		hook(t, r, prepared.Handover.ID)
	}
	r.provider.mu.Lock()
	posts, sessionGets, paymentGets := r.provider.posts, r.provider.sessionGets, r.provider.paymentGets
	r.provider.refund = 1
	r.provider.mu.Unlock()
	if posts != 1 || sessionGets != 1 || paymentGets != 1 {
		t.Fatal("unexpected provider requests", posts, sessionGets, paymentGets)
	}
	if code := r.callback(t, "evt_commercial_refund", "charge.refunded", "ch_fixture", true); code != 200 {
		t.Fatal(code)
	}
	current, err = svc.ValidateCommercialRelease(ctx, r.tenant, "store", prepared.Handover.ID)
	if err != nil || current.Current {
		t.Fatal("durable callback hold did not invalidate before GET", err)
	}
	if n, e := r.processor.ProcessOnce(ctx); e != nil || n != 1 {
		t.Fatal("refund GET", n, e)
	}
	current, err = svc.ValidateCommercialRelease(ctx, r.tenant, "store", prepared.Handover.ID)
	if err != nil || current.Current {
		t.Fatal("refunded payment remained eligible", err)
	}
	recovered, err := svc.CommercialReleaseResult(ctx, r.tenant, "store", prepared.Handover.ID, release.IdempotencyKey)
	if err != nil || recovered.ID != receipt.ID {
		t.Fatal("historical recovery", err)
	}
	var receipts, events int
	if err = r.pool.QueryRow(ctx, `select (select count(*) from sales.commercial_release_receipt where tenant_id=$1),(select count(*) from platform.outbox_event where tenant_id=$1 and event_type='commercial-release.recorded')`, r.tenant).Scan(&receipts, &events); err != nil || receipts != 1 || events != 1 {
		t.Fatal("duplicate or missing durable effect", err)
	}
	t.Logf("COMMERCIAL_RELEASE_OFFICIAL_SDK_PG_PASS quote_order_allocate_create_checkout_signed_callback_durable_job_get_payment_prepare_checklist_accept_commit_release=true post_callback_current=false refund_current=false profile_materialized=true receipts=1 events=1 provider_posts=1 live_proven=false")
}

func connectedCommercialProfile(t *testing.T, r *connectedRun, storedValue bool) (franchisejourney.HandoverReleaseContract, franchisejourney.HandoverProfileActivation) {
	t.Helper()
	ctx := context.Background()
	python := os.Getenv("HANDOVER_PROFILE_PYTHON")
	if python == "" {
		t.Fatal("explicit admitted Python runtime required")
	}
	out := filepath.Join(t.TempDir(), "policy")
	cmd := exec.CommandContext(ctx, python, "-X", "utf8", "-B", filepath.Join("..", "..", "..", "tools", "materialize_handover_profile.py"), "--output", out, "--profile-id", "franchise-commercial", "--tenant-id", r.tenant, "--organization-id", "store", "--expected-mode", "sandbox", "--payment-provider", "stripe", "--payment-account-ref", "acct_fixture", "--payment-connection-id", "checkout", "--release-effect", "commercial-receipt", "--activate")
	if storedValue {
		cmd.Args = append(cmd.Args, "--funding", "stored-value")
	}
	if body, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("materializer %v %s", err, body)
	}
	raw, err := os.ReadFile(filepath.Join(out, "activation.json"))
	if err != nil {
		t.Fatal(err)
	}
	var a franchisejourney.HandoverProfileActivation
	if err = json.Unmarshal(raw, &a); err != nil {
		t.Fatal(err)
	}
	policy, err := franchisejourney.LoadHandoverProfileFile(filepath.Join(out, "profile.json"), a)
	if err != nil || !policy.AllowsCommercialRelease(r.tenant, "store") {
		t.Fatal("hash-bound commercial profile", err)
	}
	return policy, a
}
````

### FILE: `internal/platform/postgres/commercial_release_integration_test.go`

```yaml
block_id: "GO-INITIAL-HANDOVER-API:file25:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "b1a5de13ac5a2fb8e50520a12a5e9bee5b869e63dd3f05386ab6abcf2366b0b9"
variables: []
secrets_allowed: false
```

````go
package postgres

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"elite.local/enterprise/internal/franchisejourney"
	"elite.local/enterprise/internal/platform/randomid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func commercialProfile(t *testing.T, tenant string, revision int) franchisejourney.HandoverReleaseContract {
	t.Helper()
	mode := false
	d := franchisejourney.HandoverProfileDocument{Schema: franchisejourney.HandoverProfileSchema, ProfileID: "commercial-fixture", Revision: revision, Algorithm: franchisejourney.HandoverSupportedAlgorithm, AlgorithmRevision: 2, Scope: "MATERIALIZED_PROFILE", TenantID: tenant, OrganizationID: "store", ProviderCode: "stripe", ProviderAccountRef: "acct_fixture", ProviderConnectionID: "checkout", ExpectedLiveMode: &mode, MaximumObservationAgeSeconds: 300, Options: franchisejourney.SupportedCommercialReleaseOptions(), AuthorityReference: "docs/commercial-release.md", DecisionReference: "FIXTURE_ONLY"}
	raw, err := json.Marshal(d)
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(raw)
	p, err := franchisejourney.LoadHandoverProfile(raw, franchisejourney.HandoverProfileActivation{Enabled: true, ProfileID: d.ProfileID, ProfileRevision: revision, DocumentSHA256: hex.EncodeToString(sum[:]), TenantID: tenant, OrganizationID: "store", ProviderCode: "stripe", ProviderAccountRef: "acct_fixture", ProviderConnectionID: "checkout"})
	if err != nil {
		t.Fatal(err)
	}
	return p
}

func commercialFixture(t *testing.T, pool *pgxpool.Pool, accept bool) (string, *franchisejourney.HandoverPreparationService, franchisejourney.CommitCommercialReleaseCommand) {
	t.Helper()
	ctx := context.Background()
	tenant := initialHandoverFixture(t, pool)
	for _, sql := range []string{
		`insert into integration.provider_connection(tenant_id,connection_id,provider_code,secret_ref) values($1,'checkout','stripe','FIXTURE_ONLY_NOT_SECRET')`,
		`insert into payment.provider_checkout(tenant_id,payment_attempt_id,provider_code,connection_id,account_ref,request_sha256_hex,live_mode,expires_at) values($1,'payment','stripe','checkout','acct_fixture',repeat('a',64),false,clock_timestamp()+interval '1 hour')`,
	} {
		if _, err := pool.Exec(ctx, sql, tenant); err != nil {
			t.Fatal(err)
		}
	}
	repo := NewFranchiseJourney(pool)
	svc, err := franchisejourney.NewHandoverPreparationService(repo, randomid.Generator{}, commercialProfile(t, tenant, 1))
	if err != nil {
		t.Fatal(err)
	}
	prepared, _, err := svc.Prepare(ctx, tenant, "operator", initialHandoverCommand())
	if err != nil {
		t.Fatal(err)
	}
	if accept {
		_, err = repo.PublishDeliveryChecklist(ctx, tenant, "operator", franchisejourney.DeliveryChecklist{ID: "release-checklist", OrganizationID: "store", Version: 1, Title: "Fixture", Items: []franchisejourney.ChecklistItem{{ID: "serial", Ordinal: 1, Prompt: "Read serial", ResponseType: "serial", Required: true}}}, randomid.Generator{}.New())
		if err != nil {
			t.Fatal(err)
		}
		h, err := repo.CompleteDeliveryChecklist(ctx, tenant, "store", "operator", prepared.Handover.ID, 1, "release-checklist", 1, []franchisejourney.ChecklistResponse{{ItemID: "serial", ResponseText: "SERIAL-SYNTHETIC"}}, randomid.Generator{}.New())
		if err != nil {
			t.Fatal(err)
		}
		_, err = repo.AcceptHandover(ctx, tenant, "store", "customer", h.ID, h.Version, "SERIAL-SYNTHETIC", "release-checklist", 1, strings.Repeat("e", 64), randomid.Generator{}.New())
		if err != nil {
			t.Fatal(err)
		}
	}
	return tenant, svc, franchisejourney.CommitCommercialReleaseCommand{OrganizationID: "store", HandoverID: prepared.Handover.ID, ObservationSHA256: strings.Repeat("b", 64), IdempotencyKey: "commercial-release-key"}
}

func assertNoCommercialRelease(t *testing.T, pool *pgxpool.Pool, tenant string) {
	t.Helper()
	var rows, events, keys int
	err := pool.QueryRow(context.Background(), `select (select count(*) from sales.commercial_release_receipt where tenant_id=$1),(select count(*) from platform.outbox_event where tenant_id=$1 and event_type='commercial-release.recorded'),(select count(*) from platform.idempotency_record where tenant_id=$1 and scope='commercial-release')`, tenant).Scan(&rows, &events, &keys)
	if err != nil || rows != 0 || events != 0 || keys != 0 {
		t.Fatalf("partial commercial commit: rows%d events%d keys%d err%v", rows, events, keys, err)
	}
}

func TestCommercialReleaseConnectedPostgres(t *testing.T) {
	pool := initialHandoverPool(t)
	tenant, svc, c := commercialFixture(t, pool, true)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	type outcome struct {
		v      franchisejourney.CommercialReleaseReceipt
		replay bool
		err    error
	}
	out := make(chan outcome, 12)
	start := make(chan struct{})
	for range 12 {
		go func() {
			<-start
			v, r, e := svc.CommitCommercialRelease(ctx, tenant, "operator", c)
			out <- outcome{v, r, e}
		}()
	}
	close(start)
	var receipt franchisejourney.CommercialReleaseReceipt
	creates, replays := 0, 0
	for range 12 {
		o := <-out
		if o.err != nil {
			t.Fatal(o.err)
		}
		if o.replay {
			replays++
		} else {
			creates++
			receipt = o.v
		}
		if o.v.ObservationGeneration != 1 || !o.v.ValidUntil.After(o.v.RecordedAt) || o.v.AcceptanceSHA256 != strings.Repeat("e", 64) {
			t.Fatal("receipt binding")
		}
	}
	if creates != 1 || replays != 11 {
		t.Fatalf("creates%d replays%d", creates, replays)
	}
	current, err := svc.ValidateCommercialRelease(ctx, tenant, "store", c.HandoverID)
	if err != nil || !current.Current || current.Receipt.ID != receipt.ID {
		t.Fatal("current receipt", current, err)
	}
	other := c
	other.IdempotencyKey = "different-release-key"
	if _, _, err = svc.CommitCommercialRelease(ctx, tenant, "operator", other); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("second release accepted", err)
	}
	if _, _, err = svc.CommitCommercialRelease(ctx, tenant, "other-actor", c); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("divergent replay accepted", err)
	}
	if _, err = svc.CommercialReleaseResult(ctx, tenant, "other", c.HandoverID, c.IdempotencyKey); !errors.Is(err, franchisejourney.ErrReleaseConditioned) {
		t.Fatal("scope leak", err)
	}
	var rows, events, keys int
	var stock, payment string
	var amount int64
	err = pool.QueryRow(ctx, `select (select count(*) from sales.commercial_release_receipt where tenant_id=$1),(select count(*) from platform.outbox_event where tenant_id=$1 and event_type='commercial-release.recorded'),(select count(*) from platform.idempotency_record where tenant_id=$1 and scope='commercial-release'),(select state from inventory.stock_unit where tenant_id=$1 and stock_unit_id='stock'),(select state from payment.payment_attempt where tenant_id=$1 and payment_attempt_id='payment'),(select amount_minor_units from payment.payment_attempt where tenant_id=$1 and payment_attempt_id='payment')`, tenant).Scan(&rows, &events, &keys, &stock, &payment, &amount)
	if err != nil || rows != 1 || events != 1 || keys != 1 || stock != "reserved" || payment != "captured" || amount != 123456 {
		t.Fatal("unexpected effects", rows, events, keys, stock, payment, amount, err)
	}
	for _, sql := range []string{`update sales.commercial_release_receipt set released_by_subject='forged' where tenant_id=$1`, `delete from sales.commercial_release_receipt where tenant_id=$1`} {
		if _, err = pool.Exec(ctx, sql, tenant); err == nil {
			t.Fatal("immutable history changed")
		}
	}
	if _, err = pool.Exec(ctx, `update payment.provider_observation set hold=true,hold_reason='CALLBACK_PENDING',generation=generation+1 where tenant_id=$1`, tenant); err != nil {
		t.Fatal(err)
	}
	current, err = svc.ValidateCommercialRelease(ctx, tenant, "store", c.HandoverID)
	if err != nil || current.Current || current.Receipt.ID != receipt.ID {
		t.Fatal("post-callback receipt still current", err)
	}
	v, replay, err := svc.CommitCommercialRelease(ctx, tenant, "operator", c)
	if err != nil || !replay || v.ID != receipt.ID || v.ObservationGeneration != 1 {
		t.Fatal("historical recovery lost", err)
	}
	// Even an identical hash from a newer observation generation cannot reactivate
	// an old release. The current generation is part of its evidence boundary.
	if _, err = pool.Exec(ctx, `update payment.provider_observation set hold=false,hold_reason='',observed_at=clock_timestamp() where tenant_id=$1`, tenant); err != nil {
		t.Fatal(err)
	}
	current, err = svc.ValidateCommercialRelease(ctx, tenant, "store", c.HandoverID)
	if err != nil || current.Current {
		t.Fatal("old generation reactivated", err)
	}
	t.Log("COMMERCIAL_RELEASE_PG_PASS creates=1 replays=11 atomic_receipt_outbox=true historical_recovery=true post_callback_current=false no_money_or_stock_effect=true")
}

func TestCommercialReleaseRejectsUnacceptedStaleHeldAndPolicyDrift(t *testing.T) {
	pool := initialHandoverPool(t)
	ctx := context.Background()
	for _, name := range []string{"unaccepted", "stale", "held", "wrong-hash", "policy-drift", "account-drift"} {
		t.Run(name, func(t *testing.T) {
			tenant, svc, c := commercialFixture(t, pool, name != "unaccepted")
			sql := ""
			switch name {
			case "stale":
				sql = `update payment.provider_observation set observed_at=clock_timestamp()-interval '10 minutes' where tenant_id=$1`
			case "held":
				sql = `update payment.provider_observation set hold=true,hold_reason='CALLBACK_PENDING' where tenant_id=$1`
			case "account-drift":
				sql = `update payment.provider_observation set account_ref='acct_other' where tenant_id=$1`
			case "wrong-hash":
				c.ObservationSHA256 = strings.Repeat("a", 64)
			case "policy-drift":
				var err error
				svc, err = franchisejourney.NewHandoverPreparationService(NewFranchiseJourney(pool), randomid.Generator{}, commercialProfile(t, tenant, 2))
				if err != nil {
					t.Fatal(err)
				}
			}
			if sql != "" {
				if _, err := pool.Exec(ctx, sql, tenant); err != nil {
					t.Fatal(err)
				}
			}
			if _, _, err := svc.CommitCommercialRelease(ctx, tenant, "operator", c); !errors.Is(err, franchisejourney.ErrConflict) {
				t.Fatal("ineligible release accepted", err)
			}
			assertNoCommercialRelease(t, pool, tenant)
		})
	}
}

func TestCommercialReleaseWaitsForConcurrentObservation(t *testing.T) {
	pool := initialHandoverPool(t)
	for _, mode := range []string{"callback-hold", "became-stale"} {
		t.Run(mode, func(t *testing.T) {
			tenant, svc, c := commercialFixture(t, pool, true)
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			tx, err := pool.Begin(ctx)
			if err != nil {
				t.Fatal(err)
			}
			defer tx.Rollback(ctx)
			sql := `update payment.provider_observation set hold=true,hold_reason='CALLBACK_PENDING',generation=generation+1 where tenant_id=$1`
			if mode == "became-stale" {
				sql = `update payment.provider_observation set observed_at=clock_timestamp()-interval '10 minutes' where tenant_id=$1`
			}
			if _, err = tx.Exec(ctx, sql, tenant); err != nil {
				t.Fatal(err)
			}
			done := make(chan error, 1)
			go func() { _, _, e := svc.CommitCommercialRelease(ctx, tenant, "operator", c); done <- e }()
			deadline := time.Now().Add(5 * time.Second)
			for {
				var waiting bool
				err = pool.QueryRow(ctx, `select exists(select 1 from pg_stat_activity where datname=current_database() and wait_event_type='Lock' and query like '%from payment.provider_observation where tenant_id=%')`).Scan(&waiting)
				if err != nil {
					t.Fatal(err)
				}
				if waiting {
					break
				}
				if time.Now().After(deadline) {
					t.Fatal("release never reached observation lock")
				}
				time.Sleep(10 * time.Millisecond)
			}
			if err = tx.Commit(ctx); err != nil {
				t.Fatal(err)
			}
			if err = <-done; !errors.Is(err, franchisejourney.ErrConflict) {
				t.Fatal("release ignored concurrent observation", err)
			}
			assertNoCommercialRelease(t, pool, tenant)
		})
	}
}

type commercialCollisionIDs struct{ calls int }

func (i *commercialCollisionIDs) New() string {
	i.calls++
	if i.calls%2 == 0 {
		return "018f4d4a-7b36-7a21-8d10-000000000001"
	}
	return randomid.Generator{}.New()
}
func TestCommercialReleaseRollbackAndLostResponseRecovery(t *testing.T) {
	pool := initialHandoverPool(t)
	tenant, svc, c := commercialFixture(t, pool, true)
	ctx := context.Background()
	collision, err := franchisejourney.NewHandoverPreparationService(NewFranchiseJourney(pool), &commercialCollisionIDs{}, commercialProfile(t, tenant, 1))
	if err != nil {
		t.Fatal(err)
	}
	if _, _, err = collision.CommitCommercialRelease(ctx, tenant, "operator", c); err == nil {
		t.Fatal("expected existing outbox UUID collision")
	}
	assertNoCommercialRelease(t, pool, tenant)
	v, _, err := svc.CommitCommercialRelease(ctx, tenant, "operator", c)
	if err != nil {
		t.Fatal(err)
	}
	recovered, err := svc.CommercialReleaseResult(ctx, tenant, "store", c.HandoverID, c.IdempotencyKey)
	if err != nil || recovered.ID != v.ID {
		t.Fatal("lost response recovery", err)
	}
}

type commercialWaitingIDs struct {
	event string
	calls int
}

func (i *commercialWaitingIDs) New() string {
	i.calls++
	if i.calls%2 == 0 {
		return i.event
	}
	return randomid.Generator{}.New()
}

func TestCommercialReleaseRejectsExpiryDuringOutboxWait(t *testing.T) {
	pool := initialHandoverPool(t)
	tenant, _, c := commercialFixture(t, pool, true)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	event := randomid.Generator{}.New()
	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback(ctx)
	if err = writeJourneyOutbox(ctx, tx, tenant, event, "fixture", "outbox-lock", 1, "fixture.pending", map[string]any{"fixture": true}); err != nil {
		t.Fatal(err)
	}
	if _, err = pool.Exec(ctx, `update payment.provider_observation set observed_at=clock_timestamp()-interval '298 seconds' where tenant_id=$1`, tenant); err != nil {
		t.Fatal(err)
	}
	svc, err := franchisejourney.NewHandoverPreparationService(NewFranchiseJourney(pool), &commercialWaitingIDs{event: event}, commercialProfile(t, tenant, 1))
	if err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() { _, _, e := svc.CommitCommercialRelease(ctx, tenant, "operator", c); done <- e }()
	deadline := time.Now().Add(5 * time.Second)
	for {
		var waiting bool
		err = pool.QueryRow(ctx, `select exists(select 1 from pg_stat_activity where datname=current_database() and wait_event_type='Lock' and query like '%insert into platform.outbox_event%')`).Scan(&waiting)
		if err != nil {
			t.Fatal(err)
		}
		if waiting {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("release never reached outbox lock")
		}
		time.Sleep(10 * time.Millisecond)
	}
	time.Sleep(2200 * time.Millisecond)
	if err = tx.Rollback(ctx); err != nil {
		t.Fatal(err)
	}
	if err = <-done; !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("expired checkpoint committed after outbox wait", err)
	}
	assertNoCommercialRelease(t, pool, tenant)
}
````

### FILE: `internal/platform/postgres/handover_preparation.go`

```yaml
block_id: "GO-INITIAL-HANDOVER-API:file26:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "cf3ad018508e0ba2b94eef478fd42699480d0d2d8ddc536ea31b9247b5114b17"
variables: []
secrets_allowed: false
```

````go
package postgres

// AUTHORED binding/persistence glue over existing authoritative owners. This is
// an initial-handover projection, never a shipping, pricing or fiscal engine.
import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"elite.local/enterprise/internal/franchisejourney"
	"github.com/jackc/pgx/v5"
)

type handoverPreparationReader interface {
	QueryRow(context.Context, string, ...any) pgx.Row
}

func readInitialHandover(ctx context.Context, q handoverPreparationReader, tenant, organization, order, key string) (franchisejourney.HandoverPreparation, string, error) {
	var value franchisejourney.HandoverPreparation
	var requestHash string
	var items []byte
	value.Handover.ChecklistItems = []franchisejourney.ChecklistItem{}
	err := q.QueryRow(ctx, `select h.handover_id,h.organization_id,h.order_id,h.customer_principal_id,h.stock_unit_id,h.state,h.version,
      p.order_line_id,p.reservation_id,coalesce(p.payment_attempt_id,''),coalesce(p.funding_receipt_id,''),p.observation_sha256_hex,p.contract_id,p.contract_sha256_hex,p.prepared_by_subject,p.prepared_at,i.request_sha256_hex,
      h.customer_accepted_at,coalesce(h.acceptance_evidence_sha256_hex,''),coalesce(h.checklist_id,''),coalesce(h.checklist_version,0),coalesce(t.title,''),h.checklist_completed_at,coalesce(h.supersedes_handover_id,''),
      coalesce((select jsonb_agg(jsonb_build_object('id',ci.item_id,'ordinal',ci.ordinal,'prompt',ci.prompt,'response_type',ci.response_type,'required',ci.required) order by ci.ordinal,ci.item_id) from sales.delivery_checklist_item ci where ci.tenant_id=h.tenant_id and ci.organization_id=h.organization_id and ci.checklist_id=h.checklist_id and ci.checklist_version=h.checklist_version),'[]'::jsonb)
      from platform.idempotency_record i
      join sales.delivery_handover h on h.tenant_id=i.tenant_id and h.handover_id=i.resource_id
      join sales.delivery_handover_preparation p on p.tenant_id=h.tenant_id and p.handover_id=h.handover_id and p.organization_id=h.organization_id and p.order_id=h.order_id
      join sales.customer_order o on o.tenant_id=h.tenant_id and o.order_id=h.order_id and o.organization_id=h.organization_id and o.customer_principal_id=h.customer_principal_id
      join sales.customer_order_line l on l.tenant_id=p.tenant_id and l.order_id=p.order_id and l.line_id=p.order_line_id and l.allocated_stock_unit_id=h.stock_unit_id
      left join sales.delivery_checklist_template t on t.tenant_id=h.tenant_id and t.organization_id=h.organization_id and t.checklist_id=h.checklist_id and t.checklist_version=h.checklist_version
      where i.tenant_id=$1 and h.organization_id=$2 and h.order_id=$3 and i.scope='initial-handover' and i.idempotency_key=$4
      and i.status='completed' and i.resource_type='delivery-handover' and i.response_code=201
      and i.response_body->>'handover_id'=h.handover_id`, tenant, organization, order, key).Scan(
		&value.Handover.ID, &value.Handover.OrganizationID, &value.Handover.OrderID, &value.Handover.CustomerSubject, &value.Handover.StockUnitID, &value.Handover.State, &value.Handover.Version,
		&value.OrderLineID, &value.ReservationID, &value.PaymentAttemptID, &value.FundingReceiptID, &value.ObservationSHA256, &value.ContractID, &value.ContractSHA256, &value.PreparedBy, &value.PreparedAt, &requestHash,
		&value.Handover.CustomerAcceptedAt, &value.Handover.AcceptanceEvidence, &value.Handover.ChecklistID, &value.Handover.ChecklistVersion, &value.Handover.ChecklistTitle, &value.Handover.ChecklistCompletedAt, &value.Handover.SupersedesHandoverID, &items)
	if errors.Is(err, pgx.ErrNoRows) {
		return value, "", franchisejourney.ErrNotFound
	}
	if err == nil {
		err = json.Unmarshal(items, &value.Handover.ChecklistItems)
	}
	return value, requestHash, err
}

func (r *FranchiseJourney) InitialHandoverResult(ctx context.Context, tenant, organization, order, key string) (franchisejourney.HandoverPreparation, error) {
	value, _, err := readInitialHandover(ctx, r.pool, tenant, organization, order, key)
	return value, err
}

func (r *FranchiseJourney) EvaluateInitialHandoverRelease(ctx context.Context, tenant, organization, handover, observationSHA string, policy franchisejourney.HandoverReleaseContract) (franchisejourney.HandoverReferenceRelease, error) {
	var value franchisejourney.HandoverReferenceRelease
	if !policy.AllowsScope(tenant, organization) {
		return value, franchisejourney.ErrReleaseConditioned
	}
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return value, err
	}
	defer tx.Rollback(ctx)
	command := franchisejourney.PrepareHandoverCommand{OrganizationID: organization, ObservationSHA256: observationSHA}
	var contractID, contractSHA, scope string
	var age int64
	err = tx.QueryRow(ctx, `select order_id,order_line_id,coalesce(payment_attempt_id,''),coalesce(funding_receipt_id,''),contract_id,contract_sha256_hex,contract_scope,maximum_observation_age_ns from sales.delivery_handover_preparation where tenant_id=$1 and organization_id=$2 and handover_id=$3`, tenant, organization, handover).Scan(&command.OrderID, &command.OrderLineID, &command.PaymentAttemptID, &command.FundingReceiptID, &contractID, &contractSHA, &scope, &age)
	if errors.Is(err, pgx.ErrNoRows) {
		return value, franchisejourney.ErrConflict
	}
	if err != nil {
		return value, err
	}
	if contractID != policy.ID || contractSHA != policy.DocumentSHA256 || scope != policy.Scope || age != int64(policy.MaximumObservationAge) {
		return value, franchisejourney.ErrConflict
	}
	facts, err := lockInitialHandoverScope(ctx, tx, tenant, command, policy)
	if err != nil {
		return value, err
	}
	// Existing acceptance owns handover -> order locks. NOWAIT prevents lock
	// inversion while inspecting its committed result after the payment fence.
	var accepted string
	err = tx.QueryRow(ctx, `select handover_id from sales.delivery_handover where tenant_id=$1 and organization_id=$2 and handover_id=$3 and order_id=$4 and customer_principal_id=$5 and stock_unit_id=$6 and state='accepted' and customer_accepted_at is not null and acceptance_evidence_sha256_hex is not null and checklist_completed_at is not null for share nowait`, tenant, organization, handover, command.OrderID, facts.customer, facts.stock).Scan(&accepted)
	if errors.Is(err, pgx.ErrNoRows) {
		return value, franchisejourney.ErrConflict
	}
	if err != nil {
		return value, err
	}
	if err = tx.QueryRow(ctx, `select clock_timestamp()`).Scan(&value.EvaluatedAt); err != nil {
		return value, err
	}
	value.HandoverID = accepted
	value.OrganizationID = organization
	value.ObservationSHA256 = observationSHA
	value.ContractSHA256 = policy.DocumentSHA256
	value.Scope = policy.Scope
	value.Eligible = true
	if err = tx.Commit(ctx); err != nil {
		return franchisejourney.HandoverReferenceRelease{}, err
	}
	return value, nil
}

type handoverScopeFacts struct {
	customer, stock, reservation string
	evidenceAt                   time.Time
	generation                   int64
}

// Order -> payment attempt -> observation is the same lock order as the payment
// reconciler. Time/freshness is evaluated only after these locks are held.
func lockInitialHandoverScope(ctx context.Context, tx pgx.Tx, tenant string, command franchisejourney.PrepareHandoverCommand, policy franchisejourney.HandoverReleaseContract) (handoverScopeFacts, error) {
	var facts handoverScopeFacts
	profileProvider, profileAccount, profileConnection, boundProfile := policy.PaymentBinding()
	if boundProfile && command.FundingReceiptID == "" {
		var active string
		err := tx.QueryRow(ctx, `select connection_id from integration.provider_connection where tenant_id=$1 and connection_id=$2 and provider_code=$3 and state='active' and (organization_id is null or organization_id=$4) for share`, tenant, profileConnection, profileProvider, command.OrganizationID).Scan(&active)
		if errors.Is(err, pgx.ErrNoRows) {
			return facts, franchisejourney.ErrConflict
		}
		if err != nil {
			return facts, err
		}
	}
	var currency, state string
	var total int64
	err := tx.QueryRow(ctx, `select customer_principal_id,currency,total_minor_units,state from sales.customer_order where tenant_id=$1 and organization_id=$2 and order_id=$3 for update`, tenant, command.OrganizationID, command.OrderID).Scan(&facts.customer, &currency, &total, &state)
	if errors.Is(err, pgx.ErrNoRows) {
		return facts, franchisejourney.ErrConflict
	}
	if err != nil {
		return facts, err
	}
	if total <= 0 || (state != "placed" && state != "confirmed" && state != "paid" && state != "allocated") {
		return facts, franchisejourney.ErrConflict
	}
	gross := total
	total, err = orderProviderDue(ctx, tx, tenant, command.OrganizationID, command.OrderID, currency, total)
	if err != nil || total < 0 {
		return facts, franchisejourney.ErrConflict
	}
	if total != gross && !policy.AllowsStoredValueFunding(tenant, command.OrganizationID) {
		return facts, franchisejourney.ErrReleaseConditioned
	}

	var observedAt *time.Time
	if command.FundingReceiptID != "" {
		if command.PaymentAttemptID != "" || total != 0 || !policy.AllowsStoredValueFunding(tenant, command.OrganizationID) {
			return facts, franchisejourney.ErrConflict
		}
		local, e := readLocalFunding(ctx, tx, tenant, command.OrganizationID, command.OrderID, command.FundingReceiptID, command.ObservationSHA256, currency, gross)
		if e != nil {
			return facts, franchisejourney.ErrConflict
		}
		observedAt = &local.Receipt.ObservedAt
		facts.generation = 1
	} else {
		if command.PaymentAttemptID == "" || total <= 0 {
			return facts, franchisejourney.ErrConflict
		}
		var paymentProvider, paymentReference, paymentState, paymentCurrency, paymentOrder string
		var paymentAmount int64
		err = tx.QueryRow(ctx, `select provider_code,coalesce(provider_reference,''),state,currency,amount_minor_units,order_id from payment.payment_attempt where tenant_id=$1 and payment_attempt_id=$2 for update`, tenant, command.PaymentAttemptID).Scan(&paymentProvider, &paymentReference, &paymentState, &paymentCurrency, &paymentAmount, &paymentOrder)
		if errors.Is(err, pgx.ErrNoRows) {
			return facts, franchisejourney.ErrConflict
		}
		if err != nil {
			return facts, err
		}
		if paymentState != "captured" || paymentReference == "" || paymentOrder != command.OrderID || paymentCurrency != currency || paymentAmount != total {
			return facts, franchisejourney.ErrConflict
		}
		var observedOrder, observedOrganization, observedProvider, observedReference, observedCurrency, evidence, providerStatus, account string
		var observedAmount, received, refunded, generation int64
		var live, hold bool
		err = tx.QueryRow(ctx, `select order_id,organization_id,provider_code,provider_reference,currency,amount_minor_units,received_minor_units,refunded_minor_units,provider_status,live_mode,account_ref,coalesce(evidence_sha256_hex,''),observed_at,hold,generation
      from payment.provider_observation where tenant_id=$1 and payment_attempt_id=$2 for update`, tenant, command.PaymentAttemptID).Scan(&observedOrder, &observedOrganization, &observedProvider, &observedReference, &observedCurrency, &observedAmount, &received, &refunded, &providerStatus, &live, &account, &evidence, &observedAt, &hold, &generation)
		if errors.Is(err, pgx.ErrNoRows) {
			return facts, franchisejourney.ErrConflict
		}
		if err != nil {
			return facts, err
		}
		if hold || generation < 1 || observedAt == nil || evidence != command.ObservationSHA256 || observedOrder != command.OrderID || observedOrganization != command.OrganizationID || observedProvider != paymentProvider || observedReference != paymentReference || observedCurrency != currency || observedAmount != total || received != total || refunded != 0 || live != policy.ExpectedLiveMode || !((paymentProvider == "stripe" && providerStatus == "succeeded") || (paymentProvider == "mercadopago" && providerStatus == "approved")) || account == "" {
			return facts, franchisejourney.ErrConflict
		}
		if boundProfile && command.FundingReceiptID == "" {
			if observedProvider != profileProvider || account != profileAccount {
				return facts, franchisejourney.ErrConflict
			}
			var exact bool
			err = tx.QueryRow(ctx, `select exists(select 1 from payment.provider_checkout where tenant_id=$1 and payment_attempt_id=$2 and provider_code=$3 and account_ref=$4 and connection_id=$5 and live_mode=$6)`, tenant, command.PaymentAttemptID, profileProvider, profileAccount, profileConnection, policy.ExpectedLiveMode).Scan(&exact)
			if err != nil {
				return facts, err
			}
			if !exact {
				return facts, franchisejourney.ErrConflict
			}
		}
		facts.generation = generation

	}
	facts.evidenceAt = *observedAt
	var customerStatus string
	err = tx.QueryRow(ctx, `select status from crm.customer_profile where tenant_id=$1 and customer_principal_id=$2 for share`, tenant, facts.customer).Scan(&customerStatus)
	if errors.Is(err, pgx.ErrNoRows) {
		return facts, franchisejourney.ErrConflict
	}
	if err != nil {
		return facts, err
	}
	if customerStatus != "active" {
		return facts, franchisejourney.ErrConflict
	}
	var variant string
	var lineCount, paymentCount int
	err = tx.QueryRow(ctx, `select l.variant_id,l.allocated_stock_unit_id,(select count(*) from sales.customer_order_line where tenant_id=$1 and order_id=$2),(select count(*) from payment.payment_attempt where tenant_id=$1 and order_id=$2 and state not in ('failed','refunded')) from sales.customer_order_line l where l.tenant_id=$1 and l.order_id=$2 and l.line_id=$3 and l.quantity=1 and l.allocated_stock_unit_id is not null for share of l`, tenant, command.OrderID, command.OrderLineID).Scan(&variant, &facts.stock, &lineCount, &paymentCount)
	if errors.Is(err, pgx.ErrNoRows) {
		return facts, franchisejourney.ErrConflict
	}
	if err != nil {
		return facts, err
	}
	expectedPayments := 1
	if command.FundingReceiptID != "" {
		expectedPayments = 0
	}
	if lineCount != 1 || paymentCount != expectedPayments {
		return facts, franchisejourney.ErrConflict
	}
	var stock string
	err = tx.QueryRow(ctx, `select stock_unit_id from inventory.stock_unit where tenant_id=$1 and stock_unit_id=$2 and organization_id=$3 and variant_id=$4 and state='reserved' for share`, tenant, facts.stock, command.OrganizationID, variant).Scan(&stock)
	if errors.Is(err, pgx.ErrNoRows) {
		return facts, franchisejourney.ErrConflict
	}
	if err != nil {
		return facts, err
	}
	var expires *time.Time
	err = tx.QueryRow(ctx, `select reservation_id,expires_at from inventory.serial_reservation where tenant_id=$1 and organization_id=$2 and stock_unit_id=$3 and variant_id=$4 and demand_kind='customer-order' and demand_id=$5 and demand_line_id=$6 and status='reservation' for share`, tenant, command.OrganizationID, facts.stock, variant, command.OrderID, command.OrderLineID).Scan(&facts.reservation, &expires)
	if errors.Is(err, pgx.ErrNoRows) {
		return facts, franchisejourney.ErrConflict
	}
	if err != nil {
		return facts, err
	}
	var now time.Time
	if err = tx.QueryRow(ctx, `select clock_timestamp()`).Scan(&now); err != nil {
		return facts, err
	}
	if observedAt.After(now) || now.Sub(*observedAt) > policy.MaximumObservationAge || (expires != nil && !expires.After(now)) {
		return facts, franchisejourney.ErrConflict
	}
	return facts, nil
}

func (r *FranchiseJourney) PrepareInitialHandover(ctx context.Context, tenant, actor, handoverID, eventID string, command franchisejourney.PrepareHandoverCommand, policy franchisejourney.HandoverReleaseContract, requestHash string) (franchisejourney.HandoverPreparation, bool, error) {
	var empty franchisejourney.HandoverPreparation
	if !policy.AllowsScope(tenant, command.OrganizationID) || tenant == "" || actor == "" || handoverID == "" || eventID == "" || len(requestHash) != 64 {
		return empty, false, franchisejourney.ErrReleaseConditioned
	}
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return empty, false, err
	}
	defer tx.Rollback(ctx)
	inserted, err := tx.Exec(ctx, `insert into platform.idempotency_record(tenant_id,scope,idempotency_key,request_sha256_hex,status,locked_until,expires_at)
      values($1,'initial-handover',$2,$3,'processing',clock_timestamp()+interval '30 seconds',clock_timestamp()+interval '24 hours') on conflict do nothing`, tenant, command.IdempotencyKey, requestHash)
	if err != nil {
		return empty, false, err
	}
	if inserted.RowsAffected() == 0 {
		value, storedHash, readErr := readInitialHandover(ctx, tx, tenant, command.OrganizationID, command.OrderID, command.IdempotencyKey)
		if readErr != nil || storedHash != requestHash {
			return empty, false, franchisejourney.ErrConflict
		}
		if err = tx.Commit(ctx); err != nil {
			return empty, false, err
		}
		return value, true, nil
	}
	facts, err := lockInitialHandoverScope(ctx, tx, tenant, command, policy)
	if err != nil {
		return empty, false, err
	}
	var existing bool
	if err = tx.QueryRow(ctx, `select exists(select 1 from sales.delivery_handover_preparation where tenant_id=$1 and organization_id=$2 and order_id=$3 and order_line_id=$4)`, tenant, command.OrganizationID, command.OrderID, command.OrderLineID).Scan(&existing); err != nil {
		return empty, false, err
	}
	if existing {
		return empty, false, franchisejourney.ErrConflict
	}
	_, err = tx.Exec(ctx, `insert into sales.delivery_handover(tenant_id,handover_id,organization_id,order_id,customer_principal_id,stock_unit_id,state,version)values($1,$2,$3,$4,$5,$6,'prepared',1)`, tenant, handoverID, command.OrganizationID, command.OrderID, facts.customer, facts.stock)
	if err != nil {
		return empty, false, err
	}
	_, err = tx.Exec(ctx, `insert into sales.delivery_handover_preparation(tenant_id,handover_id,organization_id,order_id,order_line_id,reservation_id,payment_attempt_id,observation_sha256_hex,contract_id,contract_sha256_hex,contract_scope,maximum_observation_age_ns,prepared_by_subject,funding_receipt_id)values($1,$2,$3,$4,$5,$6,nullif($7,''),$8,$9,$10,$11,$12,$13,nullif($14,''))`, tenant, handoverID, command.OrganizationID, command.OrderID, command.OrderLineID, facts.reservation, command.PaymentAttemptID, command.ObservationSHA256, policy.ID, policy.DocumentSHA256, policy.Scope, int64(policy.MaximumObservationAge), actor, command.FundingReceiptID)
	if err != nil {
		return empty, false, err
	}
	if err = writeJourneyOutbox(ctx, tx, tenant, eventID, "delivery-handover", handoverID, 1, "delivery-handover.prepared", map[string]any{"organization_id": command.OrganizationID, "order_id": command.OrderID, "order_line_id": command.OrderLineID, "payment_attempt_id": command.PaymentAttemptID, "funding_receipt_id": command.FundingReceiptID, "observation_sha256": command.ObservationSHA256, "contract_id": policy.ID, "contract_sha256": policy.DocumentSHA256, "actor_subject": actor}); err != nil {
		return empty, false, err
	}
	result, err := tx.Exec(ctx, `update platform.idempotency_record set status='completed',response_code=201,response_body=jsonb_build_object('handover_id',$3::text),resource_type='delivery-handover',resource_id=$3,locked_until=null where tenant_id=$1 and scope='initial-handover' and idempotency_key=$2 and status='processing'`, tenant, command.IdempotencyKey, handoverID)
	if err != nil {
		return empty, false, err
	}
	if result.RowsAffected() != 1 {
		return empty, false, fmt.Errorf("handover receipt not completed: %w", franchisejourney.ErrConflict)
	}
	value, _, err := readInitialHandover(ctx, tx, tenant, command.OrganizationID, command.OrderID, command.IdempotencyKey)
	if err != nil {
		return empty, false, err
	}
	if err = tx.Commit(ctx); err != nil {
		return empty, false, err
	}
	return value, false, nil
}
````

### FILE: `internal/platform/postgres/handover_preparation_integration_test.go`

```yaml
block_id: "GO-INITIAL-HANDOVER-API:file27:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "f6e5a4c8571d66a01b14fe61189782e87a4d484b1bd535f48ac2495d99efe509"
variables: []
secrets_allowed: false
```

````go
package postgres

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"elite.local/enterprise/internal/commerce"
	"elite.local/enterprise/internal/franchisejourney"
	"elite.local/enterprise/internal/platform/randomid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func initialHandoverPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL not set")
	}
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil || cfg.ConnConfig.Host != "127.0.0.1" || !strings.HasPrefix(cfg.ConnConfig.Database, "elite_handover_") {
		t.Fatal("requires disposable loopback handover database")
	}
	pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)
	return pool
}
func initialHandoverContract() franchisejourney.HandoverReleaseContract {
	hash := sha256.Sum256([]byte(franchisejourney.ReferenceHandoverContractDocument))
	return franchisejourney.HandoverReleaseContract{ID: "reference-single-unit-observed-payment-v1", DocumentSHA256: hex.EncodeToString(hash[:]), Scope: "LOCAL_FIXTURES", MaximumObservationAge: 5 * time.Minute}
}
func initialHandoverFixture(t *testing.T, pool *pgxpool.Pool) string {
	t.Helper()
	ctx := context.Background()
	tenant := randomid.Generator{}.New()
	fixtures := []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1::uuid,'handover-'||replace(($1::uuid)::text,'-',''),'Synthetic','Synthetic')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'store','store','Synthetic','store'),($1,'other','other','Synthetic','store')`,
		`insert into catalog.vehicle_model(tenant_id,model_id,model_code,display_name,vehicle_class,lifecycle_state)values($1,'model','model','Synthetic','bicycle','active')`,
		`insert into catalog.vehicle_variant(tenant_id,variant_id,model_id,variant_code,display_name,battery_specification,lifecycle_state)values($1,'variant','model','variant','Synthetic','{}','active')`,
		`insert into crm.customer_profile(tenant_id,customer_principal_id,display_name,email_normalized)values($1,'customer','Synthetic','fixture@example.test')`,
		`insert into crm.lead(tenant_id,lead_id,organization_id,customer_principal_id,model_id,lifecycle_state,source_code,contact_payload)values($1,'lead','store','customer','model','new','fixture','{}')`,
		`insert into pricing.price_book(tenant_id,price_book_id,market,currency,valid_from,status)values($1,'retail','AR','ARS',clock_timestamp()-interval '1 day','active')`,
		`insert into pricing.price_book_entry(tenant_id,price_book_id,variant_id,amount_minor_units,tax_mode)values($1,'retail','variant',123456,'inclusive')`,
		`insert into inventory.stock_unit(tenant_id,stock_unit_id,organization_id,variant_id,serial_number,state,version,received_at)values($1,'stock','store','variant','SERIAL-SYNTHETIC','available',1,clock_timestamp())`,
	}
	for _, sql := range fixtures {
		if _, err := pool.Exec(ctx, sql, tenant); err != nil {
			t.Fatal(err)
		}
	}
	repo := NewFranchiseJourney(pool)
	quote := franchisejourney.Quote{ID: "quote", OrganizationID: "store", LeadID: "lead", VariantID: "variant", PriceBookID: "retail", ValidUntil: time.Now().UTC().Add(time.Hour)}
	if _, _, err := repo.CreateQuoteAs(ctx, tenant, "quote-key", quote, strings.Repeat("a", 64), "018f4d4a-7b36-7a21-8d10-000000000001", "operator"); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.AcceptQuote(ctx, tenant, "store", "customer", "quote", 1, strings.Repeat("a", 64), "order", "line", "018f4d4a-7b36-7a21-8d10-000000000002", "018f4d4a-7b36-7a21-8d10-000000000003"); err != nil {
		t.Fatal(err)
	}
	sales := NewCommerce(pool)
	if err := sales.AllocateStockAs(ctx, tenant, "store", "order", "line", "stock", 1, 1, "018f4d4a-7b36-7a21-8d10-000000000004", "operator"); err != nil {
		t.Fatal(err)
	}
	payment := commerce.PaymentAttempt{ID: "payment", OrganizationID: "store", OrderID: "order", ProviderCode: "stripe"}
	if _, err := sales.RecordOrderPayment(ctx, tenant, "018f4d4a-7b36-7a21-8d10-000000000005", "payment-key", payment, "operator"); err != nil {
		t.Fatal(err)
	}
	if err := sales.TransitionPayment(ctx, tenant, "store", "payment", "created", "authorized", 1, "pi_fixture", "018f4d4a-7b36-7a21-8d10-000000000006"); err != nil {
		t.Fatal(err)
	}
	if err := sales.TransitionPayment(ctx, tenant, "store", "payment", "authorized", "captured", 2, "pi_fixture", "018f4d4a-7b36-7a21-8d10-000000000007"); err != nil {
		t.Fatal(err)
	}
	// Deliberately synthetic provider observation. The combined SDK/Webhook test
	// is owned by the payment integration; this suite claims PG binding only.
	if _, err := pool.Exec(ctx, `insert into payment.provider_observation(tenant_id,payment_attempt_id,order_id,organization_id,provider_code,provider_reference,currency,amount_minor_units,received_minor_units,refunded_minor_units,provider_status,live_mode,account_ref,evidence_sha256_hex,observed_at,hold,hold_reason,generation)values($1,'payment','order','store','stripe','pi_fixture','ARS',123456,123456,0,'succeeded',false,'acct_fixture',$2,clock_timestamp(),false,'',1)`, tenant, strings.Repeat("b", 64)); err != nil {
		t.Fatal(err)
	}
	return tenant
}
func initialHandoverCommand() franchisejourney.PrepareHandoverCommand {
	return franchisejourney.PrepareHandoverCommand{OrganizationID: "store", OrderID: "order", OrderLineID: "line", PaymentAttemptID: "payment", ObservationSHA256: strings.Repeat("b", 64), IdempotencyKey: "initial-handover-key"}
}
func TestInitialHandoverConnectedPostgres(t *testing.T) {
	pool := initialHandoverPool(t)
	tenant := initialHandoverFixture(t, pool)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	repo := NewFranchiseJourney(pool)
	service, err := franchisejourney.NewHandoverPreparationService(repo, randomid.Generator{}, initialHandoverContract())
	if err != nil {
		t.Fatal(err)
	}
	command := initialHandoverCommand()
	type outcome struct {
		value  franchisejourney.HandoverPreparation
		replay bool
		err    error
	}
	start := make(chan struct{})
	results := make(chan outcome, 16)
	var wg sync.WaitGroup
	for range 16 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			v, replay, err := service.Prepare(ctx, tenant, "operator", command)
			results <- outcome{v, replay, err}
		}()
	}
	close(start)
	wg.Wait()
	close(results)
	creates, replays := 0, 0
	handover := ""
	for result := range results {
		if result.err != nil {
			t.Fatal(result.err)
		}
		if result.replay {
			replays++
		} else {
			creates++
		}
		if handover == "" {
			handover = result.value.Handover.ID
		}
		if handover != result.value.Handover.ID || result.value.Handover.CustomerSubject != "customer" || result.value.Handover.StockUnitID != "stock" || result.value.ReservationID != "018f4d4a-7b36-7a21-8d10-000000000004" {
			t.Fatal("binding or replay drift")
		}
	}
	if creates != 1 || replays != 15 {
		t.Fatalf("creates=%d replays=%d", creates, replays)
	}
	var effects, preparations int
	if err = pool.QueryRow(ctx, `select (select count(*) from platform.outbox_event where tenant_id=$1 and event_type='delivery-handover.prepared'),(select count(*) from sales.delivery_handover_preparation where tenant_id=$1)`, tenant).Scan(&effects, &preparations); err != nil || effects != 1 || preparations != 1 {
		t.Fatalf("effects=%d preparations=%d err=%v", effects, preparations, err)
	}
	recovered, err := service.Result(ctx, tenant, "store", "order", command.IdempotencyKey)
	if err != nil || recovered.Handover.ID != handover {
		t.Fatal("lost response recovery failed", err)
	}
	if _, err = service.Result(ctx, tenant, "other", "order", command.IdempotencyKey); !errors.Is(err, franchisejourney.ErrNotFound) {
		t.Fatal("scope leak", err)
	}
	changed := command
	changed.ObservationSHA256 = strings.Repeat("c", 64)
	if _, _, err = service.Prepare(ctx, tenant, "operator", changed); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("divergent replay accepted", err)
	}
	if _, err = service.EvaluateRelease(ctx, tenant, "store", handover, command.ObservationSHA256); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("release before customer acceptance", err)
	}
	_, err = repo.PublishDeliveryChecklist(ctx, tenant, "operator", franchisejourney.DeliveryChecklist{ID: "initial-checklist", OrganizationID: "store", Version: 1, Title: "Synthetic initial checklist", Items: []franchisejourney.ChecklistItem{{ID: "serial", Ordinal: 1, Prompt: "Read the serial", ResponseType: "serial", Required: true}}}, "018f4d4a-7b36-7a21-8d10-000000000008")
	if err != nil {
		t.Fatal(err)
	}
	presented, err := repo.CompleteDeliveryChecklist(ctx, tenant, "store", "operator", handover, 1, "initial-checklist", 1, []franchisejourney.ChecklistResponse{{ItemID: "serial", ResponseText: "SERIAL-SYNTHETIC"}}, "018f4d4a-7b36-7a21-8d10-000000000009")
	if err != nil {
		t.Fatal(err)
	}
	accepted, err := repo.AcceptHandover(ctx, tenant, "store", "customer", handover, presented.Version, "SERIAL-SYNTHETIC", "initial-checklist", 1, strings.Repeat("e", 64), "018f4d4a-7b36-7a21-8d10-000000000010")
	if err != nil || accepted.State != "accepted" {
		t.Fatal(err)
	}
	decision, err := service.EvaluateRelease(ctx, tenant, "store", handover, command.ObservationSHA256)
	if err != nil || !decision.Eligible || decision.Scope != "LOCAL_FIXTURES" {
		t.Fatal("reference release gate", err)
	}
	if _, err = pool.Exec(ctx, `update payment.provider_observation set hold=true,hold_reason='CALLBACK_PENDING',generation=generation+1 where tenant_id=$1`, tenant); err != nil {
		t.Fatal(err)
	}
	if _, err = service.EvaluateRelease(ctx, tenant, "store", handover, command.ObservationSHA256); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("held observation released", err)
	}
	if _, err = service.Result(ctx, tenant, "store", "order", command.IdempotencyKey); err != nil {
		t.Fatal("historical receipt cannot be recovered", err)
	}
	t.Logf("INITIAL_HANDOVER_CONNECTED_PG_PASS tenant=%s handover=%s creates=1 replays=15 quote_order_allocation=true acceptance=true no_shipping_posted=true", tenant, handover)
}

func TestInitialHandoverRejectsInvalidObservationAndScope(t *testing.T) {
	pool := initialHandoverPool(t)
	ctx := context.Background()
	cases := map[string]string{
		"missing-observation":  `delete from payment.provider_observation where tenant_id=$1`,
		"hold":                 `update payment.provider_observation set hold=true where tenant_id=$1`,
		"wrong-amount":         `update payment.provider_observation set hold=true,amount_minor_units=123457 where tenant_id=$1`,
		"partial":              `update payment.provider_observation set hold=true,received_minor_units=1 where tenant_id=$1`,
		"refunded":             `update payment.provider_observation set hold=true,refunded_minor_units=1 where tenant_id=$1`,
		"live-mode":            `update payment.provider_observation set live_mode=true where tenant_id=$1`,
		"currency":             `update payment.provider_observation set currency='USD' where tenant_id=$1`,
		"provider-reference":   `update payment.provider_observation set provider_reference='pi_other' where tenant_id=$1`,
		"organization":         `update payment.provider_observation set organization_id='other' where tenant_id=$1`,
		"unobserved":           `update payment.provider_observation set hold=true,observed_at=null,evidence_sha256_hex=null where tenant_id=$1`,
		"stale":                `update payment.provider_observation set observed_at=clock_timestamp()-interval '1 hour' where tenant_id=$1`,
		"future":               `update payment.provider_observation set observed_at=clock_timestamp()+interval '1 hour' where tenant_id=$1`,
		"payment-pending":      `update payment.payment_attempt set state='pending' where tenant_id=$1`,
		"customer-restricted":  `update crm.customer_profile set status='restricted' where tenant_id=$1`,
		"stock-scope":          `update inventory.stock_unit set organization_id='other' where tenant_id=$1`,
		"reservation-released": `update inventory.serial_reservation set status='released' where tenant_id=$1`,
		"order-cancelled":      `update sales.customer_order set state='cancelled' where tenant_id=$1`,
	}
	for name, mutation := range cases {
		t.Run(name, func(t *testing.T) {
			tenant := initialHandoverFixture(t, pool)
			if _, err := pool.Exec(ctx, mutation, tenant); err != nil {
				t.Fatal(err)
			}
			service, _ := franchisejourney.NewHandoverPreparationService(NewFranchiseJourney(pool), randomid.Generator{}, initialHandoverContract())
			if _, _, err := service.Prepare(ctx, tenant, "operator", initialHandoverCommand()); !errors.Is(err, franchisejourney.ErrConflict) {
				t.Fatalf("invalid facts accepted or wrong error: %v", err)
			}
			var count int
			if err := pool.QueryRow(ctx, `select (select count(*) from sales.delivery_handover where tenant_id=$1)+(select count(*) from sales.delivery_handover_preparation where tenant_id=$1)+(select count(*) from platform.idempotency_record where tenant_id=$1 and scope='initial-handover')`, tenant).Scan(&count); err != nil || count != 0 {
				t.Fatalf("failed attempt leaked %d rows: %v", count, err)
			}
		})
	}
	t.Log("INITIAL_HANDOVER_NEGATIVES_PASS cases=17")
}

func TestInitialHandoverOutboxAtomicityAndFreshnessAfterLock(t *testing.T) {
	pool := initialHandoverPool(t)
	ctx := context.Background()
	tenant := initialHandoverFixture(t, pool)
	repo := NewFranchiseJourney(pool)
	command := initialHandoverCommand()
	if _, _, err := repo.PrepareInitialHandover(ctx, tenant, "operator", "rollback-handover", "018f4d4a-7b36-7a21-8d10-000000000001", command, initialHandoverContract(), strings.Repeat("d", 64)); err == nil {
		t.Fatal("duplicate outbox id unexpectedly accepted")
	}
	var count int
	if err := pool.QueryRow(ctx, `select (select count(*) from sales.delivery_handover where tenant_id=$1)+(select count(*) from sales.delivery_handover_preparation where tenant_id=$1)+(select count(*) from platform.idempotency_record where tenant_id=$1 and scope='initial-handover')`, tenant).Scan(&count); err != nil || count != 0 {
		t.Fatalf("outbox failure leaked %d rows: %v", count, err)
	}
	blocker, err := pool.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer blocker.Rollback(ctx)
	if _, err = blocker.Exec(ctx, `select 1 from payment.provider_observation where tenant_id=$1 for update`, tenant); err != nil {
		t.Fatal(err)
	}
	policy := initialHandoverContract()
	policy.MaximumObservationAge = 200 * time.Millisecond
	service, _ := franchisejourney.NewHandoverPreparationService(repo, randomid.Generator{}, policy)
	done := make(chan error, 1)
	go func() { _, _, err := service.Prepare(ctx, tenant, "operator", command); done <- err }()
	time.Sleep(350 * time.Millisecond)
	if err = blocker.Commit(ctx); err != nil {
		t.Fatal(err)
	}
	select {
	case err = <-done:
		if !errors.Is(err, franchisejourney.ErrConflict) {
			t.Fatal("stale-after-lock accepted", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("lock wait did not finish")
	}
	t.Log("INITIAL_HANDOVER_ATOMICITY_AND_POSTLOCK_FRESHNESS_PASS")
}

func TestInitialHandoverRecoveryAfterRestart(t *testing.T) {
	pool := initialHandoverPool(t)
	ctx := context.Background()
	rows, err := pool.Query(ctx, `select p.tenant_id::text,p.organization_id,p.order_id,i.idempotency_key,p.handover_id from sales.delivery_handover_preparation p join platform.idempotency_record i on i.tenant_id=p.tenant_id and i.resource_id=p.handover_id and i.scope='initial-handover' order by p.tenant_id,p.handover_id`)
	if err != nil {
		t.Fatal(err)
	}
	type identity struct{ tenant, org, order, key, id string }
	var saved []identity
	for rows.Next() {
		var row identity
		if err = rows.Scan(&row.tenant, &row.org, &row.order, &row.key, &row.id); err != nil {
			t.Fatal(err)
		}
		saved = append(saved, row)
	}
	rows.Close()
	if rows.Err() != nil {
		t.Fatal(rows.Err())
	}
	if len(saved) != 1 {
		t.Fatal("expected one durable successful preparation", len(saved))
	}
	service, _ := franchisejourney.NewHandoverPreparationService(NewFranchiseJourney(pool), randomid.Generator{}, initialHandoverContract())
	for _, row := range saved {
		got, err := service.Result(ctx, row.tenant, row.org, row.order, row.key)
		if err != nil || got.Handover.ID != row.id || got.Handover.State != "accepted" || got.Handover.CustomerAcceptedAt == nil || got.Handover.ChecklistCompletedAt == nil || got.Handover.ChecklistID != "initial-checklist" || len(got.Handover.ChecklistItems) != 1 || got.Handover.AcceptanceEvidence != strings.Repeat("e", 64) {
			t.Fatal(fmt.Sprintf("durable recovery mismatch: %+v %v", got, err))
		}
	}
	t.Log("INITIAL_HANDOVER_RESTART_RECOVERY_PASS")
}
````

### FILE: `tools/materialize_handover_profile.py`

```yaml
block_id: "GO-INITIAL-HANDOVER-API:file28:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "e3ca221f68070851a618424fa4dc8b2dd22c3ec602515a7ebbb0b55e80b8ed00"
variables: []
secrets_allowed: false
```

````python
"""AUTHORED profile materialization glue; no credentials or commercial policy invention."""
import argparse
import hashlib
import json
from pathlib import Path
import re
import uuid


def main():
    p = argparse.ArgumentParser(description="Materialize the supported initial-handover algorithm and exact activation lock")
    p.add_argument("--output", required=True, type=Path)
    p.add_argument("--profile-id", required=True)
    p.add_argument("--revision", type=int, default=1)
    p.add_argument("--tenant-id", required=True)
    p.add_argument("--organization-id", required=True)
    p.add_argument("--payment-provider", choices=("stripe", "mercadopago"), required=True)
    p.add_argument("--payment-account-ref", required=True)
    p.add_argument("--payment-connection-id", required=True)
    p.add_argument("--expected-mode", choices=("sandbox", "live"), required=True)
    p.add_argument("--maximum-observation-age-seconds", type=int, default=300)
    p.add_argument("--authority-reference", default="docs/initial-handover-reference.md")
    p.add_argument("--decision-reference", default="handover-policy/DECISION.md")
    p.add_argument("--release-effect", choices=("read-only", "commercial-receipt"), default="read-only", help="Select the durable receipt explicitly; no physical shipment or money posting")
    p.add_argument("--funding", choices=("provider-only", "stored-value"), default="provider-only", help="Explicitly select revision3 with immutable approved stored-value funding")
    p.add_argument("--activate", action="store_true", help="Explicitly activate this exact supported profile; does not prove live readiness")
    a = p.parse_args()
    if not re.fullmatch(r"[a-z][a-z0-9._-]{0,79}", a.profile_id) or not 1 <= a.revision <= 1000000:
        p.error("unsupported profile id/revision")
    try:
        if str(uuid.UUID(a.tenant_id)) != a.tenant_id:
            raise ValueError()
    except ValueError:
        p.error("tenant id must be canonical UUID")
    if any(not x or len(x) > 128 or x.strip() != x for x in (a.organization_id, a.authority_reference, a.decision_reference, a.payment_account_ref, a.payment_connection_id)):
        p.error("invalid organization or documentary reference")
    if not 1 <= a.maximum_observation_age_seconds <= 900:
        p.error("maximum observation age must be 1..900 seconds")
    if a.funding == "stored-value" and a.release_effect != "commercial-receipt":
        p.error("stored-value funding requires the explicit commercial receipt profile")
    mode = a.expected_mode == "live"
    policy = {
        "schema": "elite-handover-profile/v1", "profile_id": a.profile_id,
        "revision": a.revision, "algorithm": "single-unit-full-observed-payment",
        "algorithm_revision": 3 if a.funding == "stored-value" else 2 if a.release_effect == "commercial-receipt" else 1, "scope": "MATERIALIZED_PROFILE",
        "tenant_id": a.tenant_id, "organization_id": a.organization_id,
        "provider_code": a.payment_provider, "provider_account_ref": a.payment_account_ref,
        "provider_connection_id": a.payment_connection_id,
        "expected_live_mode": mode,
        "maximum_observation_age_seconds": a.maximum_observation_age_seconds,
        "options": {"quantity": 1, "payment_coverage": "FULL_ORDER_WITH_STORED_VALUE" if a.funding == "stored-value" else "FULL_ORDER",
                    "stock_selection": "ALLOCATED_SERIALIZED_UNIT",
                    "reservation": "ACTIVE_MATCHING", "refunds": "ZERO",
                    "disputes": "DENY", "acceptance": "CUSTOMER_AND_REQUIRED_CHECKLIST",
                    "release_effect": "COMMIT_COMMERCIAL_RELEASE_RECEIPT" if a.release_effect == "commercial-receipt" else "READ_ONLY_ELIGIBILITY"},
        "authority_reference": a.authority_reference, "decision_reference": a.decision_reference,
    }
    encode = lambda value: (json.dumps(value, sort_keys=True, indent=2) + "\n").encode("utf-8")
    body = encode(policy)
    activation = {"enabled": a.activate, "profile_id": a.profile_id,
                  "profile_revision": a.revision,
                  "document_sha256": hashlib.sha256(body).hexdigest(),
                  "tenant_id": a.tenant_id, "organization_id": a.organization_id,
        "provider_code": a.payment_provider, "provider_account_ref": a.payment_account_ref,
        "provider_connection_id": a.payment_connection_id,
                  "expected_live_mode": mode}
    decision = ("# Supported handover profile selection\n\n"
                f"Profile: {a.profile_id}@{a.revision}. Expected provider mode: {a.expected_mode}.\n"
                "This selects the existing single-unit/full-observed-payment algorithm.\n"
                "The options do not grant credit, change pricing/tax rules or post shipping.\n"
                f"Funding selection: {a.funding}. Selected release effect: {a.release_effect}. A receipt records a committed checkpoint; current validity requires a new generation/freshness/acceptance check.\n"
                "The receipt never posts inventory or proves physical dispatch, taxation or production readiness.\n"
                "Authority/decision references are documentary, not secret or live-readiness proof.\n"
                f"Activation explicitly selected: {a.activate}. Production readiness remains unproven.\n")
    a.output.mkdir(parents=False, exist_ok=False)
    files = {"profile.json": body, "activation.json": encode(activation), "DECISION.md": decision.encode("utf-8")}
    for name, data in files.items():
        with (a.output / name).open("xb") as f:
            f.write(data)
    receipt = {"schema": "elite-handover-profile-materialization/v1",
               "scope": "SUPPORTED_ALGORITHM_CONFIGURATION_ONLY", "live_proven": False,
               "files": {name: hashlib.sha256(data).hexdigest() for name, data in files.items()}}
    (a.output / "receipt.json").write_bytes(encode(receipt))
    print(json.dumps({"output": str(a.output.resolve()), "document_sha256": activation["document_sha256"], "enabled": a.activate}))


if __name__ == "__main__":
    main()
````

### FILE: `internal/franchisejourney/handover_context.go`

```yaml
block_id: "HANDOVER-DELTA-GO_INITIAL_HANDOVER_API:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "a781bb90180445e92b802df10b69f1f865b8fa378bff5a687d4f02998492042d"
variables: []
secrets_allowed: false
```

````go
package franchisejourney

// AUTHORED read composition for the operator UI. Every write still revalidates
// the actual owners; this view grants no authority to deliver or move money.
import (
	"context"
	"time"
)

type HandoverOperatorContext struct {
	OrganizationID    string    `json:"organization_id"`
	OrderID           string    `json:"order_id"`
	OrderLineID       string    `json:"order_line_id"`
	PaymentAttemptID  string    `json:"payment_attempt_id"`
	FundingReceiptID  string    `json:"funding_receipt_id,omitempty"`
	ObservationSHA256 string    `json:"observation_sha256"`
	Handover          *Handover `json:"handover,omitempty"`
	CanPrepare        bool      `json:"can_prepare"`
	ReleaseEffect     string    `json:"release_effect"`
	EvaluatedAt       time.Time `json:"evaluated_at"`
}

func (s *HandoverPreparationService) OperatorContext(ctx context.Context, tenant, org, order string) (HandoverOperatorContext, error) {
	if s == nil || !s.contract.AllowsScope(tenant, org) {
		return HandoverOperatorContext{}, ErrReleaseConditioned
	}
	if !validPreparationID(order) {
		return HandoverOperatorContext{}, ErrInvalid
	}
	r, ok := s.repository.(interface {
		InitialHandoverOperatorContext(context.Context, string, string, string, HandoverReleaseContract) (HandoverOperatorContext, error)
	})
	if !ok {
		return HandoverOperatorContext{}, ErrReleaseConditioned
	}
	return r.InitialHandoverOperatorContext(ctx, tenant, org, order, s.contract)
}
````

### FILE: `internal/platform/httpapi/handover_context.go`

```yaml
block_id: "HANDOVER-DELTA-GO_INITIAL_HANDOVER_API:file2:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "6a537c1cece601c8704745df7431f53e91334498f0687dbf25499540642f1380"
variables: []
secrets_allowed: false
```

````go
package httpapi

import (
	"context"
	"elite.local/enterprise/internal/franchisejourney"
	"elite.local/enterprise/internal/platform/identity"
	"net/http"
)

type HandoverContextService interface {
	OperatorContext(context.Context, string, string, string) (franchisejourney.HandoverOperatorContext, error)
}
type HandoverContextModule struct{ Service HandoverContextService }

func (m HandoverContextModule) Register(mux *http.ServeMux, verifier identity.Verifier) {
	if m.Service == nil {
		return
	}
	principal := franchiseJourneyAPI{verifier: verifier}
	mux.HandleFunc("GET /v1/franchise/orders/{id}/handover-context", func(w http.ResponseWriter, r *http.Request) {
		org := r.URL.Query().Get("organization_id")
		actor, ok := principal.protected(w, r, "handover:manage", org)
		if !ok {
			return
		}
		v, err := m.Service.OperatorContext(r.Context(), actor.TenantID, org, r.PathValue("id"))
		if initialHandoverError(w, err) {
			return
		}
		w.Header().Set("Cache-Control", "no-store")
		writeJSON(w, 200, v)
	})
}
````

### FILE: `internal/platform/httpapi/handover_context_test.go`

```yaml
block_id: "HANDOVER-DELTA-GO_INITIAL_HANDOVER_API:file3:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "0803cb4452ba4eaaf96d57701a94f0f82a9413d36cab0adf45b2b207d39f29b3"
variables: []
secrets_allowed: false
```

````go
package httpapi

import (
	"context"
	"elite.local/enterprise/internal/franchisejourney"
	"net/http"
	"net/http/httptest"
	"testing"
)

type handoverContextSpy struct {
	calls              int
	tenant, org, order string
}

func (s *handoverContextSpy) OperatorContext(_ context.Context, t, o, id string) (franchisejourney.HandoverOperatorContext, error) {
	s.calls++
	s.tenant = t
	s.org = o
	s.order = id
	return franchisejourney.HandoverOperatorContext{}, nil
}
func TestHandoverContextAuthorization(t *testing.T) {
	for _, mode := range []string{"allowed", "unauthenticated", "permission", "organization"} {
		t.Run(mode, func(t *testing.T) {
			p := initialHandoverPrincipal()
			want := 200
			auth := "Bearer fixture"
			switch mode {
			case "unauthenticated":
				want = 401
				auth = ""
			case "permission":
				want = 403
				p.Permissions = map[string]struct{}{}
			case "organization":
				want = 403
				p.Organizations = map[string]struct{}{}
			}
			s := &handoverContextSpy{}
			mux := http.NewServeMux()
			HandoverContextModule{Service: s}.Register(mux, journeyVerifier{principal: p})
			r := httptest.NewRequest("GET", "/v1/franchise/orders/order/handover-context?organization_id=store", nil)
			r.Header.Set("Authorization", auth)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, r)
			if w.Code != want {
				t.Fatal(w.Code, want)
			}
			if want != 200 && s.calls != 0 {
				t.Fatal("unauthorized repository call")
			}
			if want == 200 && (s.tenant != "verified-tenant" || s.org != "store" || s.order != "order" || w.Header().Get("Cache-Control") != "no-store") {
				t.Fatal("unbound context")
			}
		})
	}
}
````

### FILE: `internal/platform/postgres/handover_browser_integration_test.go`

```yaml
block_id: "HANDOVER-DELTA-GO_INITIAL_HANDOVER_API:file4:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "89c29d58d19d5aef6a2013a90cbc7aaa0c3e86a473c4e1a77660ea1b7f2da47e"
variables: []
secrets_allowed: false
```

````go
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
````

### FILE: `internal/platform/postgres/handover_context.go`

```yaml
block_id: "HANDOVER-DELTA-GO_INITIAL_HANDOVER_API:file5:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "712901c92cbf83768450b7da46187c4c9f725564ad37c57333b8475d8fbfc928"
variables: []
secrets_allowed: false
```

````go
package postgres

// AUTHORED scoped projection over admitted owners; not a second eligibility algorithm.
import (
	"context"
	"elite.local/enterprise/internal/franchisejourney"
	"errors"
	"github.com/jackc/pgx/v5"
)

func (r *FranchiseJourney) InitialHandoverOperatorContext(ctx context.Context, tenant, org, order string, p franchisejourney.HandoverReleaseContract) (franchisejourney.HandoverOperatorContext, error) {
	var v franchisejourney.HandoverOperatorContext
	if !p.AllowsScope(tenant, org) {
		return v, franchisejourney.ErrReleaseConditioned
	}
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return v, err
	}
	defer tx.Rollback(ctx)
	c := franchisejourney.PrepareHandoverCommand{OrganizationID: org, OrderID: order}
	err = tx.QueryRow(ctx, `select l.line_id,coalesce(a.payment_attempt_id,''),case when a.payment_attempt_id is null then coalesce(f.funding_id,'') else '' end,coalesce(o.evidence_sha256_hex,f.receipt_sha256)
 from sales.customer_order s join sales.customer_order_line l on l.tenant_id=s.tenant_id and l.order_id=s.order_id
 left join payment.payment_attempt a on a.tenant_id=s.tenant_id and a.order_id=s.order_id and a.state='captured'
 left join payment.provider_observation o on o.tenant_id=a.tenant_id and o.payment_attempt_id=a.payment_attempt_id
 left join lateral(select funding_id,receipt_sha256 from payment.local_funding_evidence where tenant_id=s.tenant_id and order_id=s.order_id and organization_id=s.organization_id order by created_at desc,funding_id desc limit 1)f on true
 where s.tenant_id=$1 and s.organization_id=$2 and s.order_id=$3 and l.quantity=1 and l.allocated_stock_unit_id is not null and (o.evidence_sha256_hex is not null or (a.payment_attempt_id is null and f.funding_id is not null))`, tenant, org, order).Scan(&c.OrderLineID, &c.PaymentAttemptID, &c.FundingReceiptID, &c.ObservationSHA256)
	if errors.Is(err, pgx.ErrNoRows) {
		return v, franchisejourney.ErrNotFound
	}
	if err != nil {
		return v, err
	}
	if _, err = lockInitialHandoverScope(ctx, tx, tenant, c, p); err != nil {
		return v, err
	}
	var h franchisejourney.Handover
	var policyID, policySHA string
	err = tx.QueryRow(ctx, `select h.handover_id,h.organization_id,h.order_id,h.state,h.version,coalesce(h.checklist_id,''),coalesce(h.checklist_version,0),h.checklist_completed_at,h.customer_accepted_at,p.contract_id,p.contract_sha256_hex
 from sales.delivery_handover_preparation p join sales.delivery_handover h using(tenant_id,handover_id)
 where p.tenant_id=$1 and p.organization_id=$2 and p.order_id=$3 for share of h nowait`, tenant, org, order).Scan(&h.ID, &h.OrganizationID, &h.OrderID, &h.State, &h.Version, &h.ChecklistID, &h.ChecklistVersion, &h.ChecklistCompletedAt, &h.CustomerAcceptedAt, &policyID, &policySHA)
	if err == nil {
		if policyID != p.ID || policySHA != p.DocumentSHA256 {
			return v, franchisejourney.ErrConflict
		}
		h.ChecklistItems = []franchisejourney.ChecklistItem{}
		v.Handover = &h
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return v, err
	}
	v.OrganizationID = org
	v.OrderID = order
	v.OrderLineID = c.OrderLineID
	v.PaymentAttemptID = c.PaymentAttemptID
	v.FundingReceiptID = c.FundingReceiptID
	v.ObservationSHA256 = c.ObservationSHA256
	v.CanPrepare = v.Handover == nil
	v.ReleaseEffect = "READ_ONLY_ELIGIBILITY"
	if p.AllowsCommercialRelease(tenant, org) {
		v.ReleaseEffect = franchisejourney.CommercialReleaseEffect
	}
	if err = tx.QueryRow(ctx, `select clock_timestamp()`).Scan(&v.EvaluatedAt); err != nil {
		return v, err
	}
	if err = tx.Commit(ctx); err != nil {
		return v, err
	}
	return v, nil
}
````

### FILE: `internal/platform/postgres/handover_context_integration_test.go`

```yaml
block_id: "HANDOVER-DELTA-GO_INITIAL_HANDOVER_API:file6:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "b8d1ef90e0611553e10c5580024ae205224da78366abab3cb5ca68d8c9b8ccce"
variables: []
secrets_allowed: false
```

````go
package postgres

import (
	"context"
	"elite.local/enterprise/internal/franchisejourney"
	"elite.local/enterprise/internal/platform/randomid"
	"errors"
	"testing"
)

func TestHandoverOperatorContextPostgres(t *testing.T) {
	pool := initialHandoverPool(t)
	ctx := context.Background()
	tenant := initialHandoverFixture(t, pool)
	for _, sql := range []string{`insert into integration.provider_connection(tenant_id,connection_id,provider_code,secret_ref)values($1,'checkout','stripe','FIXTURE_ONLY')`, `insert into payment.provider_checkout(tenant_id,payment_attempt_id,provider_code,connection_id,account_ref,request_sha256_hex,live_mode,expires_at)values($1,'payment','stripe','checkout','acct_fixture',repeat('a',64),false,clock_timestamp()+interval '1 hour')`} {
		if _, err := pool.Exec(ctx, sql, tenant); err != nil {
			t.Fatal(err)
		}
	}
	svc, err := franchisejourney.NewHandoverPreparationService(NewFranchiseJourney(pool), randomid.Generator{}, commercialProfile(t, tenant, 1))
	if err != nil {
		t.Fatal(err)
	}
	v, err := svc.OperatorContext(ctx, tenant, "store", "order")
	if err != nil || !v.CanPrepare || v.Handover != nil || v.PaymentAttemptID != "payment" || v.OrderLineID != "line" || v.ReleaseEffect != franchisejourney.CommercialReleaseEffect {
		t.Fatal("scoped derivation", v, err)
	}
	if _, err = svc.OperatorContext(ctx, tenant, "other", "order"); !errors.Is(err, franchisejourney.ErrReleaseConditioned) {
		t.Fatal("foreign org", err)
	}
	if _, err = svc.OperatorContext(ctx, randomid.Generator{}.New(), "store", "order"); !errors.Is(err, franchisejourney.ErrReleaseConditioned) {
		t.Fatal("foreign tenant", err)
	}
	prepared, _, err := svc.Prepare(ctx, tenant, "operator", initialHandoverCommand())
	if err != nil {
		t.Fatal(err)
	}
	v, err = svc.OperatorContext(ctx, tenant, "store", "order")
	if err != nil || v.CanPrepare || v.Handover == nil || v.Handover.ID != prepared.Handover.ID || v.Handover.Version != 1 {
		t.Fatal("prepared recovery view", v, err)
	}
	if _, err = pool.Exec(ctx, `update payment.provider_observation set hold=true,hold_reason='CALLBACK_PENDING',generation=generation+1 where tenant_id=$1`, tenant); err != nil {
		t.Fatal(err)
	}
	if _, err = svc.OperatorContext(ctx, tenant, "store", "order"); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("held observation exposed as ready", err)
	}
	if _, err = svc.Result(ctx, tenant, "store", "order", initialHandoverCommand().IdempotencyKey); err != nil {
		t.Fatal("historical recovery must remain available", err)
	}
	t.Log("HANDOVER_CONTEXT_PG_PASS exact_order_line_payment_profile=true no_mutation=true foreign_scope_denied=true hold_denied=true historical_recovery_available=true")
}
````

## 6. Configuration surface

Exact future environment and profile fields are listed in the materialized docs. Defaults disable PAYMENT_CHECKOUT_ENABLED and HANDOVER_ENABLED. Enabling requires complete scope and validation before traffic. Provider secrets, webhook secret and outbound HMAC key are external runtime inputs. Hosted return URLs must be explicit trusted HTTPS configuration. Profile ID/revision/SHA and supported algorithm/options bind the same deployment scope. No secret value is present in the pack.

## 7. Dependency bill

| Component | Fixed identity | Use / license |
|---|---|---|
| Go | 1.26.8 | standard library / BSD-3-Clause |
| PostgreSQL | 18.6 | existing transaction owner / PostgreSQL |
| pgx/v5 | exact existing root go.mod/go.sum | driver / MIT |
| official payment adapter | GO-OFFICIAL-PAYMENT-WEBHOOK-ADAPTERS0.3.2 | local SDK composition / workspace terms |
| Stripe-Go | 86.3.0, a2df585a800a97fe8ec4ebf551b4449bdb3d90a1 | external SDK dependency pin / MIT |
| MercadoPago Go SDK | 1.14.0, f910ee53fbb6819e435eaf3d0f800cb1fe74ae09 | external SDK dependency pin / MIT |
| Python | CPython3.14.4 actual reference | profile materializer only / PSF |

No vendor source is relabeled as authored product code. Exact source/license/SCA and SDK reconstruction receipts remain in PAYMENT_SDK_RECONCILIATION_V402.md.

## 8. Apply order

Compose the whole selected profile into an absent external destination. Apply migrations in numeric order;0055 precedes0056 and0057. Keep the SDK module replace path, its exact module files and the existing base owner modules. For an existing target, compare baseline and preserve durable records before coordinated migrations. The down scripts are for disposable reference/reversible pre-use schema checks, not authorization to erase production payment or handover evidence.

## 9. Verification

Run the affected Go unit/HTTP/host tests, bounded native fuzz where the selected input surface supports it, and the explicit loopback PostgreSQL connected-payment/profile/read tests with skips rejected. Run vet and build on the composed revision. Frontend selection additionally uses pinned Node/Next/Vitest/TypeScript from the admitted recipe, no unreviewed package-manager invocation. Verify wrong scope, provider mode, account, amount, URL, hash, generation, claim, duplicate callback, response loss and immutable terminal-state behavior. Compare a clean canonical reconstruction byte-for-byte before promotion. Tests of fixture contracts do not authorize live accounts or production.

## 10. Reconstruction evidence

V402: the complete77pack/915file reference was reconstructed into an absent external destination; every output matched the frozen source inventory. PostgreSQL57migrations and16 connected tests passed without skips, plus vet/build and the explicitly enumerated unit/host cases. Finite fuzz and frontend contract receipts are separately bound to identical files. See reconstruction_evidence/CONNECTED_PAYMENT_HANDOVER_V402.md and CONNECTED_DELTA_ASSURANCE_V402.md. Current SAST/SCA, whole-library operations/performance/portable release and remaining capabilities retain their own gates; no live production or complete TEST02 claim.

Canonical V402 integration: selected by the current profile with exact dependencies and caller overlays. Metadata promotion records byte reconstruction, not closure of every admission/release gate. Payload provenance is unchanged.

V402 composed delta: Connected handover browser/BFF/Go/PostgreSQL gate; bounded body, stable server date formatting, Next PageProps signatures. AUTHORED integration glue; existing domain and fixed upstreams unchanged.

V402 browser evidence: reconstruction_evidence/HANDOVER_BROWSER_V402.md. Production Webpack build includes TypeScript checking;42 direct-call and23 focused BFF/date tests PASS. Chromium desktop and mobile viewport through real BFF/API/PG cover lost-response recovery and callback invalidation. Hosted IdP/live payment/physical shipment not claimed.

V402 composed delta: V402 source-backed stored-value integration: exact remaining provider due, explicit payment/funding XOR, shared approval, bounded browser transport and optional host. See STORED_VALUE_OPERATOR_FLOW_V402.md; source/pack admission successor governs final claim. Existing provider-only behavior retained.


### FILE: `db/migrations/0068_handover_funding_evidence.up.sql`

```yaml
block_id: "GO-INITIAL-HANDOVER-API:db/migrations/0068_handover_funding_evidence.up.sql:v1"
operation: CREATE
provenance: AUTHORED
source: "Local typed source/transaction/transport/UI/recovery glue; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "2d3b85bc44247780858b01ae28e6f03726d9bdfab60aeec5fe761ce42252d6fe"
variables: []
secrets_allowed: false
```

````sql
begin;
-- AUTHORED sum-type relationship: real provider payment OR local full funding.
alter table sales.delivery_handover_preparation alter column payment_attempt_id drop not null;
alter table sales.delivery_handover_preparation add column funding_receipt_id text;
alter table sales.delivery_handover_preparation add constraint initial_handover_funding_kind check ((payment_attempt_id is not null) <> (funding_receipt_id is not null));
alter table sales.delivery_handover_preparation add foreign key(tenant_id,funding_receipt_id) references payment.local_funding_receipt(tenant_id,funding_id);
alter table sales.commercial_release_receipt alter column payment_attempt_id drop not null;
alter table sales.commercial_release_receipt add column funding_receipt_id text;
alter table sales.commercial_release_receipt add constraint commercial_release_funding_kind check ((payment_attempt_id is not null) <> (funding_receipt_id is not null));
alter table sales.commercial_release_receipt add constraint commercial_local_funding_generation check(funding_receipt_id is null or observation_generation=1);
alter table sales.commercial_release_receipt add foreign key(tenant_id,funding_receipt_id) references payment.local_funding_receipt(tenant_id,funding_id);
commit;
````

### FILE: `db/migrations/0068_handover_funding_evidence.down.sql`

```yaml
block_id: "GO-INITIAL-HANDOVER-API:db/migrations/0068_handover_funding_evidence.down.sql:v1"
operation: CREATE
provenance: AUTHORED
source: "Local typed source/transaction/transport/UI/recovery glue; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "dc5a76d6d1ea2516aae385a58d516017883c135f8bcfe0ed30e96de09e952b0f"
variables: []
secrets_allowed: false
```

````sql
begin;
do $$begin
 if exists(select 1 from sales.delivery_handover_preparation where funding_receipt_id is not null) or exists(select 1 from sales.commercial_release_receipt where funding_receipt_id is not null) then raise exception 'preserve local-funded delivery history';end if;
end $$;
alter table sales.commercial_release_receipt drop constraint commercial_release_funding_kind;
alter table sales.commercial_release_receipt drop constraint commercial_local_funding_generation;
alter table sales.commercial_release_receipt drop column funding_receipt_id;
alter table sales.commercial_release_receipt alter column payment_attempt_id set not null;
alter table sales.delivery_handover_preparation drop constraint initial_handover_funding_kind;
alter table sales.delivery_handover_preparation drop column funding_receipt_id;
alter table sales.delivery_handover_preparation alter column payment_attempt_id set not null;
commit;
````

V402 composed delta: Connected warranty reuses existing transaction/approval/stock/service owners; SQL ordering and public wrapper behavior retained. Optional host factory fails closed. Exact source tested in WARRANTY_INTERFACE_AND_PORTABILITY_V402.md; no new dependency or corporate attribution.
