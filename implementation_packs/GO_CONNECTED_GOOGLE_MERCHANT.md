# Connected Google Merchant catalog publication

## 1. Metadata

```yaml
pack_id: "GO-CONNECTED-GOOGLE-MERCHANT"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: RECONSTRUCTIBLE
  admission: CONDITIONED
claim: "Local connected Google Merchant catalog upsert/refresh: approved current source, original serial ATP/exact integer price, distinct manual approval, actual fixed SDK1.8.0, one-attempt fence, durable raw receipts and GET-only processing/recovery. AUTHORED glue around unchanged admitted owners; no Google approval or whole T2805/live claim."
stacks: ["Go 1.26.8", "PostgreSQL 18.6", "CPython 3.14 / Windows x64", "Google Merchant SDK 1.8.0"]
compatible_with: ["GO-CONNECTED-CATALOG-PUBLICATION", "GO-SUPPLY-FACTORY-INVENTORY-API", "GO-HUMAN-APPROVAL-CORE", "GO-PG-OUTBOUND-DELIVERY-FENCE", "PYTHON-GOOGLE-MERCHANT-PRODUCT-SYNC-ADAPTER", "GO-ELECTROMOBILITY-APPLICATION"]
incompatible_with: ["unbound tenant/provider/account", "production certification inferred from fixtures", "implicit source or business-policy attribution"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["google_merchant_product_sync/sdk-artifact.lock.json", "docs/merchant/official-contract.lock.json"]
verified_at: "2026-09-13"
```

## 2. Applicability

Dedicated API primary data source, current reviewed offer namespace; Windows x64/CPython3.14 fixed SDK lane; selected manual upsert/refresh and processing observation.

## 3. Architecture contract

Current source/ATP/integer micros/profile/version binding. Manual distinct review, shared immutable approval/fence, one provider attempt, GET-only recovery, raw SDK processing status and acknowledgement-based refresh queue.

## 4. Exact file manifest

```text
CREATE cmd/electromobility-api/merchant.go
CREATE cmd/electromobility-api/merchant_test.go
CREATE config/merchant.reference.json
CREATE db/migrations/0082_merchant_feed.down.sql
CREATE db/migrations/0082_merchant_feed.up.sql
CREATE docs/MERCHANT_CONNECTED_CONTRACT_PLAN.md
CREATE docs/MERCHANT_CONNECTED_REFERENCE.md
CREATE docs/merchant/official-contract.lock.json
CREATE google_merchant_product_sync/connected-runtime.lock.json
CREATE google_merchant_product_sync/connected_worker.py
CREATE google_merchant_product_sync/install_connected_runtime.py
CREATE google_merchant_product_sync/test_connected_worker.py
CREATE internal/merchantbridge/client.go
CREATE internal/merchantbridge/contract.go
CREATE internal/merchantbridge/fuzz_test.go
CREATE internal/merchantbridge/process.go
CREATE internal/merchantbridge/queue.go
CREATE internal/merchantbridge/queue_test.go
CREATE internal/merchantbridge/testfixture/server.go
CREATE internal/platform/httpapi/merchant.go
CREATE internal/platform/postgres/merchant_observation.go
CREATE internal/platform/postgres/merchant_publication.go
CREATE internal/platform/postgres/merchant_publication_integration_test.go
CREATE internal/platform/postgres/merchant_queue.go
```

## 5. Materialization blocks

### FILE: `cmd/electromobility-api/merchant.go`

```yaml
block_id: "GO-CONNECTED-GOOGLE-MERCHANT:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "21a73b8cea76f2615d69c86515c88eb3b4c177f78a4f88b35a5dc9fc4d29b749"
variables: []
secrets_allowed: false
```

````go
package main

// AUTHORED activation binding. Account credentials are supplied only in the
// future process environment; no provider call is made during module creation.
import (
	"context"
	cr "elite.local/enterprise/internal/catalogrelease"
	mb "elite.local/enterprise/internal/merchantbridge"
	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/platform/postgres"
	"encoding/base64"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"io"
	"os"
)

var errMerchantConfiguration = errors.New("merchant activation invalid")

func init() { merchantModuleFactory = selectedMerchantModule }

type merchantRuntime struct {
	Schema            string `json:"schema"`
	Python            string `json:"python"`
	PythonSHA256      string `json:"python_sha256"`
	Script            string `json:"script"`
	ScriptSHA256      string `json:"script_sha256"`
	OwnerSHA256       string `json:"owner_sha256"`
	RuntimeLockSHA256 string `json:"runtime_lock_sha256"`
	Mode              string `json:"mode"`
	FixtureOrigin     string `json:"fixture_origin"`
}

func merchantConfig(path, hash string, v any) error {
	f, e := os.Open(path)
	if e != nil {
		return errMerchantConfiguration
	}
	defer f.Close()
	raw, e := io.ReadAll(io.LimitReader(f, 32769))
	if e != nil || len(raw) > 32768 || mb.Hash(raw) != hash || mb.Decode(raw, v) != nil {
		return errMerchantConfiguration
	}
	return nil
}
func selectedMerchantModule(ctx context.Context, pool *pgxpool.Pool, price *postgres.Commerce, lookup func(string) string) (httpapi.EnterpriseModule, error) {
	if lookup == nil {
		return nil, errMerchantConfiguration
	}
	switch lookup("MERCHANT_ENABLED") {
	case "", "false":
		return nil, nil
	case "true":
	default:
		return nil, errMerchantConfiguration
	}
	if lookup("CATALOG_RELEASE_ENABLED") != "true" || pool == nil || price == nil {
		return nil, errMerchantConfiguration
	}
	var profile mb.Profile
	var runtime merchantRuntime
	if merchantConfig(lookup("MERCHANT_PROFILE_FILE"), lookup("MERCHANT_PROFILE_SHA256"), &profile) != nil || profile.Validate() != nil || merchantConfig(lookup("MERCHANT_RUNTIME_FILE"), lookup("MERCHANT_RUNTIME_SHA256"), &runtime) != nil || runtime.Schema != "elite-merchant-runtime/v1" {
		return nil, errMerchantConfiguration
	}
	cp, e := cr.LoadProfile(lookup("CATALOG_RELEASE_PROFILE_FILE"), lookup("CATALOG_RELEASE_PROFILE_SHA256"))
	if e != nil || cp.Origin != profile.Origin {
		return nil, errMerchantConfiguration
	}
	catalog, e := postgres.NewCatalogRelease(pool, cp, price)
	if e != nil {
		return nil, errMerchantConfiguration
	}
	key, e := base64.StdEncoding.DecodeString(lookup("MERCHANT_HMAC_KEY"))
	if e != nil || len(key) != 32 {
		return nil, errMerchantConfiguration
	}
	process := mb.Process{Python: runtime.Python, PythonSHA256: runtime.PythonSHA256, Script: runtime.Script, ScriptSHA256: runtime.ScriptSHA256, OwnerSHA256: runtime.OwnerSHA256, RuntimeLockSHA256: runtime.RuntimeLockSHA256, Mode: runtime.Mode, FixtureOrigin: runtime.FixtureOrigin, CredentialsFile: lookup("GOOGLE_APPLICATION_CREDENTIALS")}
	if process.Validate() != nil {
		return nil, errMerchantConfiguration
	}
	var ready bool
	e = pool.QueryRow(ctx, `select
 exists(select 1 from pg_trigger where tgrelid=to_regclass('catalog.merchant_effect') and tgname='merchant_effect_immutable' and tgenabled in('O','A'))
 and exists(select 1 from pg_trigger where tgrelid=to_regclass('catalog.merchant_observation') and tgname='merchant_observation_immutable' and tgenabled in('O','A'))
 and (select count(*) from pg_trigger where not tgisinternal and tgenabled in('O','A') and
 (tgrelid=to_regclass('approval.request') and tgname='bound_request_immutable' or tgrelid=to_regclass('approval.decision') and tgname='bound_decision_immutable'))=2
 and pg_get_functiondef(to_regprocedure('approval.guard_bound_request()')) like '%marketplace_mutation%'
 and pg_get_functiondef(to_regprocedure('approval.guard_bound_decision()')) like '%marketplace_mutation%'`).Scan(&ready)
	if e != nil || !ready {
		return nil, errMerchantConfiguration
	}
	service, e := postgres.NewMerchantPublication(catalog, profile, key, process)
	if e != nil {
		return nil, errMerchantConfiguration
	}
	return httpapi.MerchantModule{Service: service, TenantID: profile.TenantID, OrganizationID: profile.OrganizationID}, nil
}
````

### FILE: `cmd/electromobility-api/merchant_test.go`

```yaml
block_id: "GO-CONNECTED-GOOGLE-MERCHANT:file2:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "cd7c15d1ddb2e40647d1c6a39aaf54aeda6782e0535ad11358b418bfe6b1b631"
variables: []
secrets_allowed: false
```

````go
package main

import (
	"context"
	"crypto/sha256"
	cr "elite.local/enterprise/internal/catalogrelease"
	mb "elite.local/enterprise/internal/merchantbridge"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/postgres"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMerchantHostGuards(t *testing.T) {
	database := os.Getenv("PAYMENT_CONNECTED_DB_URL")
	if database == "" {
		t.Skip("owned fixture required")
	}
	ctx := context.Background()
	pool, e := pgxpool.New(ctx, database)
	if e != nil {
		t.Fatal(e)
	}
	defer pool.Close()
	cp := cr.Profile{Schema: "elite-catalog-publication/v1", TenantID: "d84c64f9-1252-4db5-bf5a-3af8026190f1", OrganizationID: "marketplace-store", Market: "AR", Currency: "ARS", Origin: "https://catalog.example.invalid", PolicyCode: cr.PolicyCode}

	var variant string
	if e = pool.QueryRow(ctx, `select payload->>'variant_id' from approval.request where tenant_id=$1 and payload->>'account_id'='123456' order by created_at desc limit 1`, cp.TenantID).Scan(&variant); e != nil {
		t.Fatal(e)
	}
	p := mb.Profile{Schema: "elite-connected-merchant/v1", TenantID: cp.TenantID, OrganizationID: cp.OrganizationID, AccountID: "123456", DataSourceID: "789012", Language: "es", FeedLabel: "AR", Currency: "ARS", CurrencyDigits: 2, Origin: cp.Origin, RefreshCadenceDays: 30, Bindings: []mb.Binding{{VariantID: variant, OfferID: "fixture-bicycle-standard", Condition: "NEW", GTINs: []string{}}}}
	fileHash := func(path string) string {
		t.Helper()
		raw, e := os.ReadFile(path)
		if e != nil {
			t.Fatal(e)
		}
		return mb.Hash(raw)
	}
	script := filepath.Join(os.Getenv("MERCHANT_TEST_ROOT"), "google_merchant_product_sync", "connected_worker.py")
	python := os.Getenv("MERCHANT_TEST_PYTHON")
	runtime := merchantRuntime{Schema: "elite-merchant-runtime/v1", Python: python, PythonSHA256: fileHash(python), Script: script, ScriptSHA256: fileHash(script), OwnerSHA256: fileHash(filepath.Join(filepath.Dir(script), "sync_product.py")), RuntimeLockSHA256: fileHash(filepath.Join(filepath.Dir(script), "connected-runtime.lock.json")), Mode: "LOCAL_FIXTURES", FixtureOrigin: "http://127.0.0.1:1"}
	config := map[string]string{"MERCHANT_ENABLED": "true", "CATALOG_RELEASE_ENABLED": "true", "MERCHANT_HMAC_KEY": base64.StdEncoding.EncodeToString([]byte(strings.Repeat("f", 32)))}
	raw, _ := json.Marshal(runtime)
	f := filepath.Join(t.TempDir(), "runtime.json")
	if e = os.WriteFile(f, raw, 0600); e != nil {
		t.Fatal(e)
	}
	config["MERCHANT_RUNTIME_FILE"] = f
	config["MERCHANT_RUNTIME_SHA256"] = mb.Hash(raw)
	for _, v := range []struct {
		prefix string
		value  any
	}{{"MERCHANT", p}, {"CATALOG_RELEASE", cp}} {
		raw, _ := json.Marshal(v.value)
		sum := sha256.Sum256(raw)
		file := filepath.Join(t.TempDir(), v.prefix+".json")
		if e = os.WriteFile(file, raw, 0600); e != nil {
			t.Fatal(e)
		}
		config[v.prefix+"_PROFILE_FILE"] = file
		config[v.prefix+"_PROFILE_SHA256"] = hex.EncodeToString(sum[:])
	}
	lookup := func(k string) string { return config[k] }
	price := postgres.NewCommerce(pool)
	if module, e := selectedMerchantModule(ctx, pool, price, lookup); e != nil || module == nil {
		t.Fatal("current merchant host", e)
	}

	catalog, e := postgres.NewCatalogRelease(pool, cp, price)
	if e != nil {
		t.Fatal(e)
	}
	process := mb.Process{Python: runtime.Python, PythonSHA256: runtime.PythonSHA256, Script: runtime.Script, ScriptSHA256: runtime.ScriptSHA256, OwnerSHA256: runtime.OwnerSHA256, RuntimeLockSHA256: runtime.RuntimeLockSHA256, Mode: runtime.Mode, FixtureOrigin: runtime.FixtureOrigin}
	service, e := postgres.NewMerchantPublication(catalog, p, []byte(strings.Repeat("f", 32)), process)
	if e != nil {
		t.Fatal(e)
	}
	principal := identity.Principal{TenantID: p.TenantID, Subject: "queue-reader", Permissions: map[string]struct{}{"merchant:read": {}}, Organizations: map[string]struct{}{p.OrganizationID: {}}}
	queue, e := service.Queue(ctx, principal)
	if e != nil || len(queue.Items) != 1 || queue.Items[0].Action != "PROVIDER_REJECTED" || queue.Items[0].RefreshDueAt == nil || queue.Items[0].Generation != 1 || queue.Items[0].DesiredProductSHA256 == "" {
		t.Fatalf("queue: %+v %v", queue, e)
	}
	principal.Organizations = map[string]struct{}{"foreign": {}}
	if _, e = service.Queue(ctx, principal); e == nil {
		t.Fatal("foreign queue admitted")
	}
	// The queue and disabled/unconfigured host must never contact the fixture
	// origin (port1 deliberately has no provider) or mutate the proved journey.
	for _, key := range []string{"MERCHANT_PROFILE_SHA256", "CATALOG_RELEASE_ENABLED", "MERCHANT_RUNTIME_SHA256", "MERCHANT_HMAC_KEY"} {
		old := config[key]
		config[key] = ""
		if _, e := selectedMerchantModule(ctx, pool, price, lookup); e == nil {
			t.Fatal("missing binding admitted", key)
		}
		config[key] = old
	}
	if _, e = pool.Exec(ctx, `alter table catalog.merchant_observation disable trigger merchant_observation_immutable`); e != nil {
		t.Fatal(e)
	}
	_, rejected := selectedMerchantModule(ctx, pool, price, lookup)
	if _, e = pool.Exec(ctx, `alter table catalog.merchant_observation enable trigger merchant_observation_immutable`); e != nil {
		t.Fatal(e)
	}
	if rejected == nil {
		t.Fatal("missing observation guard admitted")
	}
	if _, e = pool.Exec(ctx, `alter table catalog.merchant_effect disable trigger merchant_effect_immutable`); e != nil {
		t.Fatal(e)
	}
	defer func() {
		if _, e := pool.Exec(context.Background(), `alter table catalog.merchant_effect enable trigger merchant_effect_immutable`); e != nil {
			t.Error(e)
		}
	}()
	if _, e = selectedMerchantModule(ctx, pool, price, lookup); e == nil {
		t.Fatal("missing immutable effect guard admitted")
	}
	t.Log("MERCHANT_HOST_PASS exact profile/hash/catalog activation/runtime/HMAC/migration required; queue exact and scoped; no provider network call")
}
````

### FILE: `config/merchant.reference.json`

```yaml
block_id: "GO-CONNECTED-GOOGLE-MERCHANT:file3:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "cf18bea8a5e8ea831a8cfbb4ae0230af5fb7cc4c797b6c2fd6159a764f076162"
variables: []
secrets_allowed: false
```

````json
{
  "schema": "elite-connected-merchant/v1",
  "tenant_id": "d84c64f9-1252-4db5-bf5a-3af8026190f1",
  "organization_id": "marketplace-store",
  "account_id": "123456",
  "data_source_id": "789012",
  "content_language": "es",
  "feed_label": "AR",
  "currency": "ARS",
  "currency_digits": 2,
  "origin": "https://catalog.example.invalid",
  "refresh_cadence_days": 28,
  "bindings": [
    {
      "variant_id": "replace-with-approved-catalog-variant",
      "offer_id": "reference-bicycle",
      "condition": "NEW",
      "brand": "",
      "gtins": []
    }
  ]
}
````

### FILE: `db/migrations/0082_merchant_feed.down.sql`

```yaml
block_id: "GO-CONNECTED-GOOGLE-MERCHANT:file4:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "74891d832aed9f3aea611d117ec33b76b2e067f84b67db91d44aa8abbaa36855"
variables: []
secrets_allowed: false
```

````sql
begin;
do $$ begin
 if exists(select 1 from catalog.merchant_effect) then raise exception 'cannot remove merchant feed evidence';end if;
end $$;
drop table catalog.merchant_observation;
drop table catalog.merchant_effect;
commit;
````

### FILE: `db/migrations/0082_merchant_feed.up.sql`

```yaml
block_id: "GO-CONNECTED-GOOGLE-MERCHANT:file5:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "6342dbe8894edc320fc31884768049c7f04b2324eb7440b2136fa55a5d2076e2"
variables: []
secrets_allowed: false
```

````sql
-- AUTHORED binding to the existing catalog, manual approval and delivery fence.
begin;
create table catalog.merchant_effect(
 tenant_id uuid not null,channel_code text not null check(channel_code='google_merchant'),
 delivery_key text not null,account_id text not null,offer_id text not null,
 generation bigint not null,source_sha256 text not null check(source_sha256~'^[0-9a-f]{64}$'),
 profile_sha256 text not null check(profile_sha256~'^[0-9a-f]{64}$'),
 request_sha256 text not null check(request_sha256~'^[0-9a-f]{64}$'),
 created_at timestamptz not null default clock_timestamp(),
 primary key(tenant_id,channel_code,delivery_key),
 foreign key(tenant_id,generation) references catalog.release_publication(tenant_id,generation),
 foreign key(tenant_id,delivery_key) references approval.request(tenant_id,request_id),
 foreign key(tenant_id,channel_code,delivery_key) references communication.outbound_delivery(tenant_id,channel_code,delivery_key)
);
create index merchant_effect_offer_idx on catalog.merchant_effect(tenant_id,account_id,offer_id);
create trigger merchant_effect_immutable before update or delete on catalog.merchant_effect for each row execute function catalog.release_immutable();
create table catalog.merchant_observation(
 tenant_id uuid not null,channel_code text not null default 'google_merchant' check(channel_code='google_merchant'),
 approval_id text not null,operation text not null check(operation in('INSERT','GET')),
 result_state text not null check(result_state in('OBSERVED','NOT_FOUND','REJECTED','UNKNOWN')),
 confirmed boolean not null,
 response bytea not null check(octet_length(response)<=32768),
 response_sha256 text not null check(response_sha256~'^[0-9a-f]{64}$'),
 observed_at timestamptz not null default clock_timestamp(),
 primary key(tenant_id,approval_id,operation,result_state,response_sha256),
 foreign key(tenant_id,channel_code,approval_id) references catalog.merchant_effect(tenant_id,channel_code,delivery_key),
 check(not confirmed or result_state='OBSERVED'),
 check(operation<>'INSERT' or confirmed)
);
create trigger merchant_observation_immutable before update or delete on catalog.merchant_observation for each row execute function catalog.release_immutable();
commit;
````

### FILE: `docs/MERCHANT_CONNECTED_CONTRACT_PLAN.md`

```yaml
block_id: "GO-CONNECTED-GOOGLE-MERCHANT:file6:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "aaf7bfa52ba77ea93730647d2a62baec0bd356df7c7cb67cd94c3acc0ca4469c"
variables: []
secrets_allowed: false
```

````markdown
# Connected Merchant contract decisions

Contract decisions for the connected local claim; final receipts and remaining T2805 scope are in MERCHANT_CONNECTED_REFERENCE.md.

Reuse PYTHON_GOOGLE_MERCHANT_PRODUCT_SYNC_ADAPTER1.8.0 SDK pin, its exact
build_insert_request and official ProductInput/ProductAttributes/Price models.
Reuse the admitted catalog source, serial ATP, human approval and outbound fence.
The Go mapping, bounded Python IPC, PostgreSQL binding and fixture transport are
AUTHORED glue. No copied tutorial or new pricing/inventory algorithm.

Selected reference: one offer in a dedicated API primary data source. Offer,
language/feed label, account/data source and product metadata are explicitly
profile/source-bound. Source model specification must contain reviewed
merchant_description; never manufacture a description, brand or GTIN.
The product link and immutable PNG URL use the current catalog origin/routes.
Price minor units convert exactly to integer micros; overflow or unrepresentable
currency precision is rejected. Availability follows original current serial ATP
with a recorded horizon, without claiming physical receipt or reserving stock.

Use official ProductInput.version_number for primary-source insertion and compare
Product.version_number when observing the processed product. Bind it to source
generation for this dedicated reference namespace; lower versions are rejected
by the provider. Refresh at the same generation is supported by the documented
contract. Existing foreign version namespaces require target admission; no
automatic source takeover or shared-account migration is implied.

One manual approval binds the exact provider request, source generation, profile,
offer and ATP. One provider attempt is admitted by the original durable fence.
Persist the exact SDK response before subsequent processing observation. Missing
acknowledgements remain UNKNOWN; recovery is GET-only by the deterministic
account/language/feed-label/offer resource name. Confirm identity/data source,
generation and selected fields. NotFound is pending, never proof of no effect.

An accepted ProductInput is not Google approval. Processed attributes can differ
because of rules/supplemental sources. Preserve product_status/issues separately,
report mismatch, and never mark ads/product approval from an insert200.
An exact GET match after an ambiguous insert proves observed desired state, not
that the provider's30day freshness deadline was refreshed by that write.

The explicit refresh-due/manual approval queue and operator reconciliation path are connected. Queue reads never re-send expired manual requests. Automated communications remain the separate T2805 job/notification scope.

Source acquisition: merchant-connected-authority/acquisition.json reacquires the
same20official wheels with exact hashes and fixes three current official docs;
runtime-receipt.json restores the pinned Windows/CPython3.14 environment offline,
passes pip check and loads actual SDK1.8.0 insert/get/product_status interfaces.
No accounts, billing action or live product write occurred.
````

### FILE: `docs/MERCHANT_CONNECTED_REFERENCE.md`

```yaml
block_id: "GO-CONNECTED-GOOGLE-MERCHANT:file7:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "e414eb5852200306fc10e0ab5bfc0efb7c6d0bc142546fe53015af80402eff83"
variables: []
secrets_allowed: false
```

````markdown
# Connected Google Merchant reference

Library fixture claim: current reviewed catalog, original serial ATP and exact
integer minor-unit price -> manual distinct approval -> official SDK1.8.0
ProductInput insertion/upsert -> immutable PostgreSQL receipt -> processing GET.
The Go/Python process, SQL and HTTP binding is AUTHORED glue. Google owns the
unchanged generated SDK; original source/SDK wheel/20 dependencies are fixed in
the existing sdk-artifact.lock.json and requirements-windows-py314.lock.
docs/merchant/official-contract.lock.json fixes current official API docs.

Select a dedicated API primary data source and a stable offer namespace in
config/merchant.reference.json. It contains synthetic examples, no credential.
Bind its variant to the reviewed catalog and set model.specification's
merchant_description explicitly. Source generation is the official positive
int64 version_number. Updates and refreshes use insert/upsert at that version;
source rollback is a new monotonically increasing local publication generation.
No takeover of a foreign version namespace is inferred. Supplemental/rule
transformations remain observable divergence, not approval.

## Runtime and host

Use the fixed Windows x64 / CPython3.14 lane. Acquire the exact twenty wheels
from their URLs/hashes in connected-runtime.lock.json (the original hash-locked
requirements file is the same graph). Run install_connected_runtime.py with
--wheels ABS_CACHE --destination ABS_ABSENT_RUNTIME --python-sha256 BASE_SHA
--lock-sha256 LOCK_SHA using the admitted base Python. It validates all wheels
before creating a venv, installs offline without dependency resolution, runs
pip check, verifies1041 installed payload files and emits merchant-runtime.json
plus its hash in setup-receipt.json. The script never reads an account/secret.
Source-specific fixtures run with test_connected_worker.py --python VENV_PYTHON
--receipt ABS_ABSENT_RECEIPT. This uses the actual fixed REST SDK on loopback.

At future activation supply MERCHANT_ENABLED=true,
MERCHANT_PROFILE_FILE / MERCHANT_PROFILE_SHA256,
MERCHANT_RUNTIME_FILE / MERCHANT_RUNTIME_SHA256,
MERCHANT_HMAC_KEY (32 bytes, base64) and GOOGLE_APPLICATION_CREDENTIALS through
the runtime secret environment. Catalog activation/profile/origin must match;
migration0082 and original immutable approval/fence guards must be active.
The emitted runtime defaults to CREDENTIALS; only fixture harnesses use
LOCAL_FIXTURES plus an explicit127.0.0.1 origin and no credential file.
Enabled invalid configurations fail startup; disabled modules make no calls.
The old standalone sync_product.py remains an admitted lower-level utility;
the connected host uses its exact request builder/models, with PostgreSQL
receipts and the original outbound fence replacing its standalone file runner.

## Operator workflow

Authenticated tenant/organization role endpoints under /v1/admin/merchant:
- GET /queue: bounded per-profile view (maximum32 offers). Read only.
- POST /prepare: approval_id, generation as string, variant_id, expires_at.
- POST /requests: submit the exact returned intent.
- POST /requests/{id}/decision: request_sha256, explicit approve, reason;
  reviewer must differ from requester.
- POST /requests/{id}/send with {}: one provider attempt for that approval.
- GET /requests/{id}: durable approval, fence, input receipt and processing.
- POST /requests/{id}/reconcile with {}: GET only, never resends the insertion.

Permissions are merchant:read/request/approve/send/reconcile, with verifier
tenant and organization scope. All request bodies/time/process/response sizes
are bounded; duplicate/case-folded/unknown JSON fields fail closed.

CURRENT_INPUT means an acknowledged input with a future refresh deadline,
not Google product approval. A GET match after a lost insertion proves desired
state but leaves freshness unconfirmed. Input acknowledgements alone anchor
refresh_due_at. At REFRESH_APPROVAL_REQUIRED or changed source, prepare a new
request and review it; expired/manual approvals are never auto-renewed.
UNKNOWN/SENDING blocks new attempts for that offer across profile changes.
Missing resources, incorrect data source/version/fields and rejected writes
stay explicit. Product status/issues are preserved as exact SDK JSON bytes.
Polling /queue does not create jobs or send products; its operator read path is
the selected manual feed-renewal workflow. Automated campaigns/reminders are
the separate T2805 notification owner. This reference selects upsert/refresh
of its bound offers; bulk feeds, multi-source takeover and automated removal
are not claimed as equivalent provider operations.

## Evidence and limits

MERCHANT_CONNECTED_RELEASE_V402.md/json binds the current pack, SDK/runtime
restore, connected API/PG proof, host/migration and finite fuzz results.
Input processing can be pending or altered by Google; no insert200 or local
PASS demonstrates product approval, live account acceptance or advertising.
Future account credentials/access are the connected provider activation
condition. Local fixtures require none. T2805 remains open until its other
selected mappings and communications workflow are connected.
````

### FILE: `docs/merchant/official-contract.lock.json`

```yaml
block_id: "GO-CONNECTED-GOOGLE-MERCHANT:file8:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "5c4c6cf948b20e25e619f6170fa6d1cb76acc1f2467015a294f8ad28c04bd7b7"
variables: []
secrets_allowed: false
```

````json
{
  "schema": "elite-official-contract-snapshot/v1",
  "verified_at": "2026-09-13",
  "sources": [
    {
      "id": "products",
      "url": "https://developers.google.com/merchant/api/guides/products/add-manage",
      "sha256": "87770da1baad98d13df098cf0e0a800dff1a1ecbd6334ff2586faca9a0f52443",
      "bytes": 375991
    },
    {
      "id": "insert",
      "url": "https://developers.google.com/merchant/api/reference/rest/products_v1/accounts.productInputs/insert",
      "sha256": "5900b8cb15b5fbb463dd6b882a9939863cbef11540d9caada7826c6ef37b7f4a",
      "bytes": 485464
    },
    {
      "id": "get",
      "url": "https://developers.google.com/merchant/api/reference/rest/products_v1/accounts.products/get",
      "sha256": "01f659734f4a952a011addd88383550702df706b3876718e2cb6a0265b9fa3ca",
      "bytes": 486194
    }
  ],
  "sdk_artifact_lock": "google_merchant_product_sync/sdk-artifact.lock.json",
  "sdk_artifact_lock_sha256": "b070f37394487a1fd674c9e492543275b232255379e03095ea52008ca536f0ee",
  "installed_source_lock": "google_merchant_product_sync/connected-runtime.lock.json",
  "installed_source_lock_sha256": "4b8b7a660fb91db30f5bca4caf863f1fdd9ecd1f6676293cd9d12eed1fc11043",
  "mapping_provenance": "AUTHORED; original Google SDK is DEPENDENCY_PIN, not rewritten or attributed to this bridge"
}
````

### FILE: `google_merchant_product_sync/connected-runtime.lock.json`

```yaml
block_id: "GO-CONNECTED-GOOGLE-MERCHANT:file9:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "4b8b7a660fb91db30f5bca4caf863f1fdd9ecd1f6676293cd9d12eed1fc11043"
variables: []
secrets_allowed: false
```

````json
{
  "abi": "Windows x86-64",
  "excluded": "Installed RECORD/generated installer files and wheel data outside purelib/platlib; original dependency license obligations remain unchanged.",
  "files": {
    "_cffi_backend.cp314-win_amd64.pyd": "9e99b2184c60f21a88cafc374386ded3559e8694ed7cf111f2869e066d303334",
    "certifi-2026.7.22.dist-info/METADATA": "ef5af1638fbb23676ac3c5777dfcfc2cd9c348fe4172ed5ba3d277655b248090",
    "certifi-2026.7.22.dist-info/WHEEL": "2b6eb4118ce7cd7b09601406aa623c553c4476265836f0d9c16f5c061f7efcc0",
    "certifi-2026.7.22.dist-info/licenses/LICENSE": "e93716da6b9c0d5a4a1df60fe695b370f0695603d21f6f83f053e42cfc10caf7",
    "certifi-2026.7.22.dist-info/top_level.txt": "28cbb8bd409fb232eb90f6d235d81d7a44bea552730402453bffe723c345ebe5",
    "certifi/__init__.py": "7a52e5f4205f9b3d12a31898ceaa7e3f6837d19cc230700bcae1a9d91d1486c1",
    "certifi/__main__.py": "c410688fdd394d45812d118034e71fee88ba7beddd30fe1c1281bd3b232cd758",
    "certifi/cacert.pem": "9cc2a774b5198dcff14d9be1e66091f538975d867ce029a96bce15a55dfd730f",
    "certifi/core.py": "5c55f2727746e697f7edac9e17c377d8752e0da7ecca191531b3b80403d61dad",
    "certifi/py.typed": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
    "certifi/tests/__init__.py": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
    "certifi/tests/test_certify.py": "7c179d1edf37f7e912b133663735fb46c6338d5d0e5f8915d11c3bc287d227b8",
    "cffi-2.1.1.dist-info/METADATA": "67b6a7595df7ede2195016a05d285faa7d362a8f57c238984edb8bb674b5770d",
    "cffi-2.1.1.dist-info/WHEEL": "a48533112e6f0981ca051e807f15f10c5caa8d9f5c90285ebac78f7ead45612a",
    "cffi-2.1.1.dist-info/entry_points.txt": "a6b6108011f8e27050fe5269bba88c8051b414697880d3b0b6fbe3a03a0b001c",
    "cffi-2.1.1.dist-info/licenses/LICENSE": "5ba24ddc57067f9249add644c3afc41a5d6dc37e23433ef759d95df370b0af63",
    "cffi-2.1.1.dist-info/top_level.txt": "ac4ed6477ad97cd2b1588f7e8e7ea1b0708097b303901f859ae41bc568c57a14",
    "cffi/__init__.py": "c47b2db28c5f197e242d4dd4af95b7ebb0da65763fc76972cafd06f93e942cda",
    "cffi/_cffi_errors.h": "5b383ed01bb5fcad8610ad29494e5be90d86000b9d71a17cec7346f0ba05996c",
    "cffi/_cffi_gen_src.py": "c94479de91526ec6b7b2c17bfd53aa5d368a78be402531ccfddc99f48198f6cc",
    "cffi/_cffi_include.h": "cd8e8a168a3f6e96b4d8f6a042cd61cde5bd3b9eac9e2a25a32c6abd1280a2df",
    "cffi/_embedding.h": "3e101d00c70a0a5013559fa505c65e0c4735d85e85a9575a1ed689d8f920e84f",
    "cffi/_imp_emulation.py": "4714441bccc06c8d913c6070c3dd2eff97e2f2c59d6a1a5d8a93a83f3929ec2d",
    "cffi/_shimmed_dist_utils.py": "0638f6c26f3265bbc5bd6131e4011f9aa6aa6726458587e8c8b2d01e45d9b9aa",
    "cffi/api.py": "7165f5665b1c98ec9392fb61881a5458b712441b6a9699fe9a806ee1eb019b32",
    "cffi/backend_ctypes.py": "4fb092735e0fe6fca5a22e61ca8c7af2a841597d3bf68c075da1836d8fd4edba",
    "cffi/cffi_opcode.py": "e7ed5c53aa691c6cecb2d17f7a88250fff2d08393bebf6619338f77d6e1289cd",
    "cffi/commontypes.py": "ecdeb33ed08596fc57316847574f29b148dd6082b65e0b0ddf2a39760b9afefe",
    "cffi/cparser.py": "598fa8c785064ae552613b609b4d8598681b99cb9a6812b79a816591d1529ac5",
    "cffi/error.py": "bfac53892e14d24bc3732e21fc10d1a39bf7f5942e8fe20c4582efe444dd759b",
    "cffi/ffiplatform.py": "6afc458dd8a460626812d9893bb7b0566c06fd511597a119fd668d859602aafe",
    "cffi/gen_src.py": "4d4fd9543c282e71cf95a97cfe63e6a71ab2589ed4ef9da34dcdd4b13984d502",
    "cffi/lock.py": "97d4d37703083298ba8c39091a742013d72f4c847b0809ed209afc1061edde96",
    "cffi/model.py": "700f9b10089b70751bb5ba9f80cab408b42314f02357425247333229588d805e",
    "cffi/parse_c_type.h": "39dc107f033d92dababe5081e377b11509b10c1b63d8c04d74af0b625d79b63c",
    "cffi/pkgconfig.py": "bb2b75d7b4d079e6b77636a03355b9828f11a67c921938280b337c407c997657",
    "cffi/recompiler.py": "4c9545d3b985cd819357ab33a76ccfac977674928737e781fd2ba22621a5f4ac",
    "cffi/setuptools_ext.py": "a59ae74ea54290632abbc37123191d6c8b65eb981e84c3f52e8ada51d4fbfb85",
    "cffi/vengine_cpy.py": "ad1a2e610e819650085cd8d91933d45b9e2b07de76aa2de0726d7a4965d628b7",
    "cffi/vengine_gen.py": "67f2aa685eabe789f98b886341385975bfb50c14c48a69d4e32e04fcac520838",
    "cffi/verifier.py": "c48cd5f220cc859d7eaa15f74515f80d6a36934ed6a5c271761840157023ae66",
    "charset_normalizer-3.5.1.dist-info/METADATA": "f00192aad9979ce50bdab2bfd471ad0705aee5a33f4bf6d27c1bd24417c873e5",
    "charset_normalizer-3.5.1.dist-info/WHEEL": "f8868e546779e082b2f4ef1acdce705bd9cc24829ee8add213946b7906f8242d",
    "charset_normalizer-3.5.1.dist-info/entry_points.txt": "0034932ab91767786174e5458ba0dc5041d0452d317f10c813fa44cf8c0b2170",
    "charset_normalizer-3.5.1.dist-info/licenses/LICENSE": "18577485d3704f1a479ded8e573c0976cfed315fd2fd17983fa988da4c2f70d1",
    "charset_normalizer-3.5.1.dist-info/top_level.txt": "ec04b2cde3ebf3fc6e65626c9ea263201b7257cbe1128d30042bf530f4518b74",
    "charset_normalizer/__init__.py": "c9c9fcadec07e5fbc1c07b9aaae8b54c238f49524b40629f7ddfd44719ce770a",
    "charset_normalizer/__main__.py": "dac8ff052e87d2c536e42d5b32acfd0d5c1aea438af657214846d253efe2bbb3",
    "charset_normalizer/api.py": "9032c5c8a60b2c19b5101d1543e4aa8b0ded26618d28a3d45c3372454cf20a6a",
    "charset_normalizer/cd.cp314-win_amd64.pyd": "f234a1dbd19363a3260d6a1e2d82ec1c6e30b2e23d6057796021231d668d978a",
    "charset_normalizer/cd.py": "f5b50bb7ecb028ed2f769ac02757b6320abdd4aa30e155125e480f15c28db8db",
    "charset_normalizer/cli/__init__.py": "77d314c7ed55fea0f7c7d8a0232e094f8a02e42534ca3ba593b4326567911618",
    "charset_normalizer/cli/__main__.py": "1a29e6b4a0d8d75b16f4fe7dda781cfe0f8c159e76008481655348f693551887",
    "charset_normalizer/constant.py": "9b4e0a49913370b4283d26466bec73af7dd1cdf13bd70d9325ac453f82e8628f",
    "charset_normalizer/legacy.py": "8d5a1916b9fc14b890f4b68db9a715e72189a854aa067430a65f035252649bdb",
    "charset_normalizer/md.cp314-win_amd64.pyd": "faf5cb006492b3526abcbfdac9a0d39cd52d32425292db9d2a210d8b44d9d925",
    "charset_normalizer/md.py": "e2ba134fe7bcfdd6bcd2d0cbccc71e38cc7dc5a25b8118b74b4d4f1e1f4f9c2d",
    "charset_normalizer/models.py": "01189116f8501996fd87b57f12639d763df21bcd801b9fc768db3a5c9519354c",
    "charset_normalizer/py.typed": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
    "charset_normalizer/utils.py": "4dc210b6670b40d9332925b893dd792b4437365f147abe00251a08606b35a853",
    "charset_normalizer/version.py": "c6e9896047e40b6a966dae4565e75bf5b42790a9c01a53046358bbac1723a110",
    "cryptography-50.0.1.dist-info/METADATA": "8e9b5570da4686848ee5c6bb83fffc224bcddf9793136731dd050ce0c05fffe2",
    "cryptography-50.0.1.dist-info/WHEEL": "9efe56059dbb731c438c822ac4132a2b670e174db90f09059aa246a8bdd2a575",
    "cryptography-50.0.1.dist-info/licenses/LICENSE": "3e0c7c091a948b82533ba98fd7cbb40432d6f1a9acbf85f5922d2f99a93ae6bb",
    "cryptography-50.0.1.dist-info/licenses/LICENSE.APACHE": "aac73b3148f6d1d7111dbca32099f68d26c644c6813ae1e4f05f6579aa2663fe",
    "cryptography-50.0.1.dist-info/licenses/LICENSE.BSD": "602c4c7482de6479dd2e9793cda275e5e63d773dacd1eca689232ab7008fb4fb",
    "cryptography-50.0.1.dist-info/sboms/cryptography-rust.cyclonedx.json": "ca81fd6162fbb9c39d1c4530c4beb49dd3f938a77cb867f82ce788def949e3ce",
    "cryptography-50.0.1.dist-info/sboms/sbom.json": "3375a1384a4782d4a871ffe6d00d22926c140ab2db7c3671087e15a1ff239234",
    "cryptography/__about__.py": "49425154c5ad66852ef967592df63ba7fb36c66e5842334b163f3462b50ff45d",
    "cryptography/__init__.py": "9ad86e52b4dde0544e0a9518ad322a863cfab3a4fd763019ad5ee7675a0c9b6a",
    "cryptography/cobblestone.py": "4fff85d3cf8051fde58daca965e92535357585639b1f9ee8a38e66552ce4f707",
    "cryptography/exceptions.py": "f37e445882dcd9fc31c3e83214cae27220b64aa8558844c1d742c14b2b670c87",
    "cryptography/fernet.py": "dac2c23ce99e8e13956ee75c4b742c3353c8b876f109e45c46d5a49dc437982a",
    "cryptography/hazmat/__init__.py": "e48c2b2d6ad5a7402312bff815d586fd5d39ecd489198fd6e1e80d36cb9cb748",
    "cryptography/hazmat/_oid.py": "3a17b69d413a5e32b1d5417a4c0b06532daff06509f958bd0d2e32c52ba4b795",
    "cryptography/hazmat/asn1/__init__.py": "62a613a071858cb38a8a0442d341c68fa86b8c4c1bbb548fa408aeebebbe7daa",
    "cryptography/hazmat/asn1/asn1.py": "f6671fae676f2bb2ad8e1458a1f3330235b4bbb9efa9b7ecbcf7db84d40ad414",
    "cryptography/hazmat/backends/__init__.py": "3b98ef28541d6675e129ea89f87b6e95a10bf4d8bb9abd660f3658e641e56212",
    "cryptography/hazmat/backends/openssl/__init__.py": "a778e625f9c26a0f62139b1d32b37a56f544bb9e6ee3ac5a4bf223a08d12ae60",
    "cryptography/hazmat/backends/openssl/backend.py": "d0d73315b8b01b17eaf6a8b6f281d8af4f2bfebed5420e005c78dbd2814fc2a0",
    "cryptography/hazmat/bindings/__init__.py": "c0709b59f69e7daf9f93a4c74b0f6d87d7c952c4ad268ef6e39c1f141aa676e0",
    "cryptography/hazmat/bindings/_rust.pyd": "6a42262974f0e086c1f2defaff4ccd9fb25046cf193422135e2dda5b4b29b2a2",
    "cryptography/hazmat/bindings/_rust/__init__.pyi": "b949e7bceabfb9a99e62704da66b2fdd244054ed16406a63239f1a47594be0ac",
    "cryptography/hazmat/bindings/_rust/_openssl.pyi": "4ffe53a583bc6e08483e3010dc1b40001ab0f82718b241e9f76c26e4f358e9d7",
    "cryptography/hazmat/bindings/_rust/asn1.pyi": "06b1a30bc27a9f0b92fabdc455c757241f27768b5f63d99b41839fa5b3c6d070",
    "cryptography/hazmat/bindings/_rust/cobblestone.pyi": "f6551fd56d614ebb33109e18a527093483d20fd04e33ba3e0bcdbe9ad4831bda",
    "cryptography/hazmat/bindings/_rust/declarative_asn1.pyi": "30198e6732d37663c9ebed4ae80ab73d3b872753068117c8bb484881a83ba1c1",
    "cryptography/hazmat/bindings/_rust/exceptions.pyi": "7b15ebdb1c3fd29075924f777186ccdcca216f3a149233a6b3564c522d2e4191",
    "cryptography/hazmat/bindings/_rust/ocsp.pyi": "54f556b8a1c8f4432cd3d64b458006bd1d08cf49823261335c0920247a2fa683",
    "cryptography/hazmat/bindings/_rust/openssl/__init__.pyi": "3ba18ecaf0f1b7aca589c0766851aebb945c199cbd5232a762a191f2bab3bddc",
    "cryptography/hazmat/bindings/_rust/openssl/aead.pyi": "eddfbec5d734beecd17b84b6bae41a1a6fbd1862afe0d7122d070bc9671a0fab",
    "cryptography/hazmat/bindings/_rust/openssl/ciphers.pyi": "2e13f31d649726ae20ac02579facd2bd249d57e698208b1c1c3c083e524618fb",
    "cryptography/hazmat/bindings/_rust/openssl/cmac.pyi": "9cf1f45f9ed1629b00911a305698d0887139eba4e1513c7b617aec69d9ab9879",
    "cryptography/hazmat/bindings/_rust/openssl/dh.pyi": "6774c2f86d38f931ed49d00e3cb3358761bb9a5e5b75ad4495951c9f9c29ba19",
    "cryptography/hazmat/bindings/_rust/openssl/dsa.pyi": "a81b64823d9a95bb76a8572767d503ae1ce83610953bb1d36f2e5549fd445cc2",
    "cryptography/hazmat/bindings/_rust/openssl/ec.pyi": "cc9cb4a516b99fefe9d9d9b8e4fc44081ffe07a49567234a7e3c450e93ea4f7f",
    "cryptography/hazmat/bindings/_rust/openssl/ed25519.pyi": "5577d77791ba8548af837f7d4750d87665b77936f411e6f30d3aa3442da0691c",
    "cryptography/hazmat/bindings/_rust/openssl/ed448.pyi": "631e3d96a7678ec0fb6f188357591c68cac392dba0e5ebe2e5aeb37ab322cb6b",
    "cryptography/hazmat/bindings/_rust/openssl/hashes.pyi": "6c57b74f5dd47b24489d424c0fe47902b19318ba960e25a8d8bbf6286cd46dfe",
    "cryptography/hazmat/bindings/_rust/openssl/hmac.pyi": "057667ecd0e32f72406d85b4490f298358b20b90db4175615008b0b22f03151f",
    "cryptography/hazmat/bindings/_rust/openssl/hpke.pyi": "a7f8462e7e981fe11aac91755796d4b14b638a9be2100a5c4793b4b141c92ed7",
    "cryptography/hazmat/bindings/_rust/openssl/kdf.pyi": "1e288b74407d9ea32c4ef330885a694297cfcfa7c7533f7f8298c407f5c03970",
    "cryptography/hazmat/bindings/_rust/openssl/keys.pyi": "b5e22df0ce9910c26b9f8b375b45275b40d9fb7d0977af169d2b0a286d76e25d",
    "cryptography/hazmat/bindings/_rust/openssl/keywrap.pyi": "9860891284389718b8091393044b94cdbf54e24cbc9f9d8bde20ef08e5cd0512",
    "cryptography/hazmat/bindings/_rust/openssl/mldsa.pyi": "1fce170d6d134372c785d1890f9b8f6834f98a31222ba319cad5d873cac135fe",
    "cryptography/hazmat/bindings/_rust/openssl/mlkem.pyi": "6f2a1ebace3e9483a55a5122450742734d2a0a45e9381e6978163aca8f6293cb",
    "cryptography/hazmat/bindings/_rust/openssl/poly1305.pyi": "fd25bd36d4391439406dd72516d5a94f89469b1288287a4dfb88fc27607361f4",
    "cryptography/hazmat/bindings/_rust/openssl/rsa.pyi": "d8e4023525e4c6073edeec35c620829682104d5ea9f7f90aefd62ed2b85980f7",
    "cryptography/hazmat/bindings/_rust/openssl/x25519.pyi": "7b09f81a94326fbccfc04fa78bb1adc9080c0b403598bbaa62c4b2aafea723fb",
    "cryptography/hazmat/bindings/_rust/openssl/x448.pyi": "8ee4d94e6962f233bfe5572e7e0faf1efc7fb42c9ece648b221ffd3d4dd37332",
    "cryptography/hazmat/bindings/_rust/pkcs12.pyi": "bc411de700e266f6fc64615ace22c26962f3030a06fedbcf5312d0c39fee3a5f",
    "cryptography/hazmat/bindings/_rust/pkcs7.pyi": "b711812628ea66c84472aadae9bc8f35b9e2b08765c733922073f6865f5aacfb",
    "cryptography/hazmat/bindings/_rust/test_support.pyi": "3cf86577e5a43bbe3789714f79b786d0bb602b4693cc6763708b1acb51a6e461",
    "cryptography/hazmat/bindings/_rust/x509.pyi": "2d54f886d04db1ee608a3421e0685535811b00ed04d40360136f2d5aac378382",
    "cryptography/hazmat/bindings/openssl/__init__.py": "c0709b59f69e7daf9f93a4c74b0f6d87d7c952c4ad268ef6e39c1f141aa676e0",
    "cryptography/hazmat/bindings/openssl/_conditional.py": "bdfce2fb11dd4ad01c2508cbaaa3a54a8b3de0bf03feda85b417bb3afc3e2ad6",
    "cryptography/hazmat/bindings/openssl/binding.py": "16cb8bc5b4e7aea5bad3e7405ad6a482e7deabb4e7a9a41a8b53dcec3a30ae3a",
    "cryptography/hazmat/decrepit/__init__.py": "c0709b59f69e7daf9f93a4c74b0f6d87d7c952c4ad268ef6e39c1f141aa676e0",
    "cryptography/hazmat/decrepit/ciphers/__init__.py": "c0709b59f69e7daf9f93a4c74b0f6d87d7c952c4ad268ef6e39c1f141aa676e0",
    "cryptography/hazmat/decrepit/ciphers/algorithms.py": "5923e3336b177eb2c35ebb4de51ecdeab0355c5edbf772357bf96f0f8cb1da80",
    "cryptography/hazmat/decrepit/ciphers/modes.py": "3aafcf1300a47b938b73339faffbf300e27ac0fc7eae5b2f1c05cf06e8796fda",
    "cryptography/hazmat/primitives/__init__.py": "c0709b59f69e7daf9f93a4c74b0f6d87d7c952c4ad268ef6e39c1f141aa676e0",
    "cryptography/hazmat/primitives/_asymmetric.py": "46181ca2e501e874e214306b4752f1aa4323a54c4888dbd0d6bff3263446eaa4",
    "cryptography/hazmat/primitives/_cipheralgorithm.py": "121de2ee5c1e7477e2d1e2d2b07f773d9c502b363d23a9642baeef2f8579b4e7",
    "cryptography/hazmat/primitives/_modes.py": "fc4880383f096fa5a915796ec247144f7102cff4c7f318c6bdb4251a4366e7d4",
    "cryptography/hazmat/primitives/_serialization.py": "862d3125b941019f29999c769566be96bba9871befbc6de0f4334ad06f0e75fb",
    "cryptography/hazmat/primitives/asymmetric/__init__.py": "c0709b59f69e7daf9f93a4c74b0f6d87d7c952c4ad268ef6e39c1f141aa676e0",
    "cryptography/hazmat/primitives/asymmetric/dh.py": "467b767cb683a501234125fc7cdb659935ce2478e5b7f0ca419d5595f85d8dae",
    "cryptography/hazmat/primitives/asymmetric/dsa.py": "910cd3562a8023d3e210b34e05ca3f9fa31f3f737d7a120f91276e3a9ba9b743",
    "cryptography/hazmat/primitives/asymmetric/ec.py": "2e495508b2f16db433aa213a1d16984d16ae2cca61614c37b44e9cf28048c584",
    "cryptography/hazmat/primitives/asymmetric/ed25519.py": "99b7df2efb3edcf8020de835ad05d55e1ec62e7b6554be004e390811ce29b793",
    "cryptography/hazmat/primitives/asymmetric/ed448.py": "9f474aeb2f952f1bb0fb3beb0f1b8f4f8e7df21ca6817fb5472fe3abd6944211",
    "cryptography/hazmat/primitives/asymmetric/mldsa.py": "2c731a694d26eba12c4ab4a87567f87b691d99d0d14c41910b5cbe7c22b118ad",
    "cryptography/hazmat/primitives/asymmetric/mlkem.py": "c2f5647854fee287fe55c500111dbecc8ff98574faaf04c93cc090f2017d1930",
    "cryptography/hazmat/primitives/asymmetric/padding.py": "8db81a8ef18763090f8cc0e58c7bb99713019437e60b15947b6005d92e8a135d",
    "cryptography/hazmat/primitives/asymmetric/rsa.py": "94813b9467b8e3d5b95112d3282a6f1467ab213a2e2a8655f46c2e4f349cb12a",
    "cryptography/hazmat/primitives/asymmetric/types.py": "4a32614eb9c544b9533dd62b043d69a0c8ba723606f9f8747d221a427dd0522c",
    "cryptography/hazmat/primitives/asymmetric/utils.py": "42cf117bd1853e35bfb4da7fef721e2448cf09fd2c94e22c58ec68eaaca64fa9",
    "cryptography/hazmat/primitives/asymmetric/x25519.py": "950720524f8f8ae6e3f329c535641e6267b7b5864c80742a9112a42033872f3a",
    "cryptography/hazmat/primitives/asymmetric/x448.py": "6f7ee283b93ba681ba4890ebcb8da3f849428f9ab8790d17e421663744a7f33f",
    "cryptography/hazmat/primitives/ciphers/__init__.py": "7b21179a393afc265768e3d80ebeef018197af6f50bf38162f6fb8092a252c5b",
    "cryptography/hazmat/primitives/ciphers/aead.py": "173972c7bc3c29841a9330e9d735a026722beb6ce066b815875bb68787b1079e",
    "cryptography/hazmat/primitives/ciphers/algorithms.py": "2067c2c9c6093a1d53be8237e2b7f0f0a7c9ff19add69fc3a4186dd8582171ea",
    "cryptography/hazmat/primitives/ciphers/base.py": "6810bb1c7041a22c5e6e6a5aad5af4525383b37543d00ec1ea8cff01a4630eff",
    "cryptography/hazmat/primitives/ciphers/modes.py": "1fb0194a02624930428b482b2fdecdc5a912f79d09c852982ecadaf2d8c73f14",
    "cryptography/hazmat/primitives/cmac.py": "b33fece87fdc6273afc7e54d59d20a85185edd89a9f33f09d03dc206a397de08",
    "cryptography/hazmat/primitives/constant_time.py": "c5dba7593d277fc3af29d72a52186514a6b21a9e3f3e0549454d96d702d2aff0",
    "cryptography/hazmat/primitives/hashes.py": "4abb7081cdffa45bd97c88518fd7b5a7859d8be17f3e995e18796a240460b4e8",
    "cryptography/hazmat/primitives/hmac.py": "469077cfdcf9b248ab090ae6ef341bb67a7da4b327023ae54d4bcaa85e5a0c37",
    "cryptography/hazmat/primitives/hpke.py": "46c1e2b2c0b997e7533e7d69d896dcac8ac09c306b2ea339ba0045486a662aee",
    "cryptography/hazmat/primitives/kdf/__init__.py": "bf7ca2601194dbbd84a233576f07d865d6db7c8f5c54e0821b79d78530dd6b79",
    "cryptography/hazmat/primitives/kdf/argon2.py": "649c7e7a7525000e28edbe3a0b5c25cc6fd7010416d3022416050a92b33fc00f",
    "cryptography/hazmat/primitives/kdf/concatkdf.py": "0459ce292ef6b71285d3242f012af435aeee873bcb1e95027242b988d002fda4",
    "cryptography/hazmat/primitives/kdf/hkdf.py": "33494011f468738929a78fa7c038fdc81faf359ba420e60442b5255ac04d9fda",
    "cryptography/hazmat/primitives/kdf/kbkdf.py": "0b7c3df6f8c5aa1ab3031cc5387e33cef5681e46f9002e8bcfa00440657b77f8",
    "cryptography/hazmat/primitives/kdf/pbkdf2.py": "5dc86493df6cf8c93a6fa2920a8baa764f0557288d54e48521e9e69d6ded2e04",
    "cryptography/hazmat/primitives/kdf/scrypt.py": "5f259475452686e23d57a4ea00f3afba3092306bf55d0760d1adb52160a63be5",
    "cryptography/hazmat/primitives/kdf/x963kdf.py": "40a8e145e86e4ec02cb4bdfce1b29fbe3bafd727dd6a6007561b014567a7d95c",
    "cryptography/hazmat/primitives/keywrap.py": "b7db7c299bb367852cdf4012521236ecc64f3cf2b2e041114439002be94ced34",
    "cryptography/hazmat/primitives/padding.py": "413f94f8dbd5d9e4063b5c153db0e234635273d91e4434be8a0e5c40eacbcf41",
    "cryptography/hazmat/primitives/poly1305.py": "3f910f415f9107f1493da869834d6ed13b384bf3e7026b2ba312065db19e451a",
    "cryptography/hazmat/primitives/serialization/__init__.py": "7d772ee3dba02fdaf1869465ff06639cca84f430a484231cb11c5f8ffa93bde8",
    "cryptography/hazmat/primitives/serialization/base.py": "48006bbe8d12f3a41ca01963ec3a6378b5b2da5cd032bca28a827a0eb1184150",
    "cryptography/hazmat/primitives/serialization/pkcs12.py": "992f5c14d1b869fcefb1ea1ce5ed4c5a863656c91f2fc37c63f3858e5ebb96e6",
    "cryptography/hazmat/primitives/serialization/pkcs7.py": "98533b385b99f1c0b1506528a380a87fa06708773f21bc750108e37f19cf971a",
    "cryptography/hazmat/primitives/serialization/ssh.py": "162b177bf9d429d3c67ea10d5612a99a86b399a23ca87f067be5466dcd1dca4c",
    "cryptography/hazmat/primitives/twofactor/__init__.py": "b66319181fa0e08535afb94816a012534d7dcebd2e3e9ff010161cc1d0c22820",
    "cryptography/hazmat/primitives/twofactor/hotp.py": "8af668e41adc08658bb2aaa5e27655d175c28f218f8bf8877c3165b4694e2709",
    "cryptography/hazmat/primitives/twofactor/totp.py": "9b92cfa512f4d24a78cd8f204e3af9e477f3f5a3253d2e7790799591240299d6",
    "cryptography/py.typed": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
    "cryptography/utils.py": "21aed2d106734b47b3c2920741f63c1b57bea381a3ac78a31175e3c00af58ae2",
    "cryptography/x509/__init__.py": "d0a54d4e646d758357fe2fdec2c4cdcdb4e6d9ed6ac9fa0e85ce1e9a883984be",
    "cryptography/x509/base.py": "50830a9309acc75d01477fe09c14410814f0429d7e3bcc247d59bb6b31fefe61",
    "cryptography/x509/certificate_transparency.py": "26aa0e203865c089eb60c156e88167efb589d2f885f8f07fae5655defb3d3187",
    "cryptography/x509/extensions.py": "fa25dcc169c82f5ce5e4c977055245904a3fef59dc5def387a1beacd6d6c2b83",
    "cryptography/x509/general_name.py": "b0ffeb575d50969b24e31dd75c6258fccbf443fb3d76d8de2dc907b23a4b4284",
    "cryptography/x509/name.py": "a08db0e95758f33830eeaa4d865c44c79d765b40215835d0ac361fe44b0cd8f6",
    "cryptography/x509/ocsp.py": "61ecba35d155d4c3e3a29db8323fd57a78c4a60de451a328a5258e2b401b781b",
    "cryptography/x509/oid.py": "054ce05d71955a296406474a3d39bd478a84944f600061e07583cc640a7b3c9a",
    "cryptography/x509/verification.py": "811d82d9cf97650b5b959853e53e6f8d228eb426fbe1e7f66a53d59847301653",
    "google/_upb/_message.pyd": "e05ec368a9fe4d2b723128c8a6ba7b68b7cbbfb140ecf1ab0f68ba6b11cc60ca",
    "google/api/annotations.proto": "e79ea741cb605a65e78ca322174764a4af9fde1962c1631e12b84c4934ba9a6c",
    "google/api/annotations_pb2.py": "d0fd397fbcb2132243c9eb8417d6cf733abf1678bfb60e6568698cbbf2b3e9fb",
    "google/api/annotations_pb2.pyi": "c80b1a4674059631b94d97854f5cb781cbd885e96b7607f4081550ec5ca6c448",
    "google/api/auth.proto": "01264e4e5b8ee20388dc350066704cad3816bf51430b580f2ccdbdbab36ac336",
    "google/api/auth_pb2.py": "8bf777f4253d0a7b14908cf519dc4cb8ceeba86f236d29537355bd460c7e2014",
    "google/api/auth_pb2.pyi": "9e38f4c1d4911ef3a96cfde078dda4fb70e9fa30692a50a64d6ff28aea64532f",
    "google/api/backend.proto": "8d56380cc90a7cbe87766e53bb228baf90539058421a78a487b81eb6fbc5a5cc",
    "google/api/backend_pb2.py": "82da69e52f3bb4d577c5582f0e8a51870fb3b9212eef32f78e4b131c8ee6f960",
    "google/api/backend_pb2.pyi": "b4bb8fc696705afc62a2403bd12dfb167062ff11e241d47b97e5bb96a1ab6829",
    "google/api/billing.proto": "09e4ad9b42f918160bc82f14de80ab93603158600246e4cd5aa88c7bf7156f26",
    "google/api/billing_pb2.py": "d150935f133f969c96530fc8fcef1989ec6378731ffbbf5ddd4f81cbe25f1263",
    "google/api/billing_pb2.pyi": "21b66a891fe2440281996ccd1c97f2fa3586921c20f1ff90f0a0b4c78b15b8ab",
    "google/api/client.proto": "5e9fc117290fe7d3472791543c979bec7cf4497f3fe6ffaf9de9ae99cc4f1804",
    "google/api/client_pb2.py": "3a0e05b3879929b2cdf4e3a5d05e3d065ac35de30387b8d5ce3cf458711135a1",
    "google/api/client_pb2.pyi": "48c44f7e24b08b9a7a88a389ab2dbcb3e9d92f1248663c8806c85a2ffaf1c628",
    "google/api/config_change.proto": "569e013add9082005f06ca659738cf29c06b86a03d6afb6cd95449603bc5a2d5",
    "google/api/config_change_pb2.py": "38e96f65f31b1706b1e0c96e9d6bc9d535ebb25064dedcaf0359856c966f816f",
    "google/api/config_change_pb2.pyi": "1a94a919baa0757e1e01cfe1e5f05161a386dcffc0ce32836c81ed3a75e23302",
    "google/api/consumer.proto": "999acb311628d9c2941d0cfdaabdf18047be2ee061dda2d9543d79aa6defec84",
    "google/api/consumer_pb2.py": "233ea7b55a1bb540b7e5edddd3bb207ce802f20f41194f73727443efbd73c970",
    "google/api/consumer_pb2.pyi": "281c8d2ea76854f929ef70518829e66ebd4d9c61278ea025c4188804e9b08a31",
    "google/api/context.proto": "9a66fb2aa3ba46df37e3bd377ba8b842025a4418af78f38beb998af5aaf97185",
    "google/api/context_pb2.py": "4398a401461baf24a6820f5d389637e213e7692cf77e11320638348a821a9993",
    "google/api/context_pb2.pyi": "ffdabbcdd94528c2b3aa58be38bebdbc54984cee09e396361d890b8ec5a5e41f",
    "google/api/control.proto": "581efb54aea34a7cd12cc06e0dd2e3155a713645269ecef215629723338965fe",
    "google/api/control_pb2.py": "c4caa461e5224dacfa057613a36c3564e8995947a6636d60135209e57313988d",
    "google/api/control_pb2.pyi": "d70642612d11f8b48eda054f3aabdb55badd4fea37fd53d049dc024664e27072",
    "google/api/distribution.proto": "94bd1a17dbe20bcfb9ad14e160647b3eb2e09a95edaba3437b83b19d118fc286",
    "google/api/distribution_pb2.py": "cccc6ca5e834cec97c50d9e7c572858fcd9badd6e2b30bc713b788a5fc02cb51",
    "google/api/distribution_pb2.pyi": "a9fde1acd7b9e0e352afa390f3e213c2d5bf56ae153a8f635c00a37c33d6408a",
    "google/api/documentation.proto": "a632846b2a6a06de4ab3880d842f3142a385105553a17ffc9ffe3718b8c58747",
    "google/api/documentation_pb2.py": "80a256587bb8d53d2ad2608cf68e3b847aae365f6d5002f8b8e37ca150743a4a",
    "google/api/documentation_pb2.pyi": "219c794c7ca8aeff97ec049e910c93bb4f3d4f959a10931ba208c2c58d67b3fc",
    "google/api/endpoint.proto": "a57c5c460f6c8ec3dca195c6bd04a24bea67aea379ac95acb00973de668dc38d",
    "google/api/endpoint_pb2.py": "4bc22ea330e92608eef7cfea174819cc02a3122d0c1318604e1ebdf5f5adbeb8",
    "google/api/endpoint_pb2.pyi": "2e9f4056b68c8b20feeec141f6ea4d9708187c2e51fd3f177e7fb42c21f0acea",
    "google/api/error_reason.proto": "c69dbc781135d7ca94485ec674e82cce9458b01ae706906a5b73eb275cb5130e",
    "google/api/error_reason_pb2.py": "2cd4837bda6544d063863398cfe8e9b7df0b8f8cb621ed15b382ec07f6f889c2",
    "google/api/error_reason_pb2.pyi": "6d8edc78dea110b6243a4692bee6f6f08cb4670abc78c1788dfefb64c112445d",
    "google/api/field_behavior.proto": "bb7ed77725c95f751ec2c47cd92359d0137c39383b8a62aef41bab7aa7478da5",
    "google/api/field_behavior_pb2.py": "9d3b745733ce8b1823eed691459e54684229a59c326dc34c0f3018cc5db52e30",
    "google/api/field_behavior_pb2.pyi": "1ca5532291dfbe1d8a25a6899f2aaacd506098dffca32a6eb495487229799d20",
    "google/api/field_info.proto": "53bb0a673642b4efff2cb48f6eb9dbb0693b198ac837da684ae54f751e8cf39e",
    "google/api/field_info_pb2.py": "58dd3cbcc0c7e02884267fb81a62c23de3a0ba0e8920d2c618004eace488c492",
    "google/api/field_info_pb2.pyi": "7225adf1dee76f393227292177b2c3f81914cab21ba4b0daf41013b677b00e7c",
    "google/api/http.proto": "ead99129aa15dd5f6233942030c72eec33bc0d7b1c7c260dc143293ef66c5b78",
    "google/api/http_pb2.py": "a174349b065eb8da9214dc8c7d117f89c60d03ff5d23d43f42b96387afaab839",
    "google/api/http_pb2.pyi": "4ce2870969d6582d245ae53e5a2d0a5dcfde5e75f10c4692879342c35ef46d1a",
    "google/api/httpbody.proto": "454a5102396e030edb0dae5c09056c8953041a27f665e4b1e3c12280f6ed2421",
    "google/api/httpbody_pb2.py": "9b97a8aef3ebb549e35bc1ceb4e20a9b4e5afa6e836ffd372a17df8729ac4d07",
    "google/api/httpbody_pb2.pyi": "4635efaa21a0b4322c5b4bb146316946882283a8ef7b6041fa431670b4265af4",
    "google/api/label.proto": "dadfde6d5ecce36b7e3232e2439a0d5ac83357500732a2b7760e89cc66111806",
    "google/api/label_pb2.py": "5aed3c3e2e321e9503f8d42ad131274bfb469c50312388536fa51cc7838e94d5",
    "google/api/label_pb2.pyi": "7bdc5b89373b596c17a3d9ee764aff620a7288db9978e0974e35968645b345e4",
    "google/api/launch_stage.proto": "68e22c9821227b4e84c7029569fce1055c21e6ccab991e0986bc75f12ef4dd82",
    "google/api/launch_stage_pb2.py": "931cfaf85389c2a06ef1a6fafe4856b132b463d4a9f3dcdb3949283cc95c5b66",
    "google/api/launch_stage_pb2.pyi": "12d535eae3b5a506128f8726d42240ce003e3530452e7776ab6602d68b46a5f7",
    "google/api/log.proto": "4989e9f19f6dba74de6ae2915a046fbbec8f7aa61b55552f45317df7ac8fec27",
    "google/api/log_pb2.py": "70995696d347784974c85be2b7dde3ed13d9449ee6c1b000aab316ca8ea3ae8a",
    "google/api/log_pb2.pyi": "1caa4dbe7e867865c66f94a5b0679e2c7db14b87d57886aee2a87394f769b934",
    "google/api/logging.proto": "803958363b6412ed765ecae4e2aab2986b411637c5cb49158196575e61f9eb58",
    "google/api/logging_pb2.py": "cb310725f86c3173d0a6cf7858551e15dbf30713c14b2cf22515e8297be9b0b4",
    "google/api/logging_pb2.pyi": "607d699e4fc440fb1db2d049cd3515cdffc3faedaaf55a41047bbc92ecab60c9",
    "google/api/metric.proto": "021f0ef550db23d98f8bafb74aca4189bed165911c8890bf8769ad052c7ad6cd",
    "google/api/metric_pb2.py": "50b29c00a99d8191c4548a4a1bd834a3aef456c5580352e47d9eeeb072c9ffe8",
    "google/api/metric_pb2.pyi": "a66c38c4d665faaf83de991c1d41024be46d10b8b4e87026bf88d09e5fd9fcf1",
    "google/api/monitored_resource.proto": "cf7d64333f5b6287adb0eb8bdbe552bc38813a85ce8134c39cbe1c56be7e9cd2",
    "google/api/monitored_resource_pb2.py": "8ddaf6c06f86116afc1e10ee2d082256a81c2ccd132f4210a651973139462feb",
    "google/api/monitored_resource_pb2.pyi": "137d9ec248e0f451edf6019c5b7afa5596cc8edbbc7299918a0e801e0eb3abb9",
    "google/api/monitoring.proto": "06563dc549e756dc559ba73136ca7c5e94ae52f9a88f44fd33d5d99e83fccecf",
    "google/api/monitoring_pb2.py": "925d0e8f664ccee7a94293303ab27f0042bc7af8f3fcb060f89afd8b5c5d32c0",
    "google/api/monitoring_pb2.pyi": "1c987fe2b3e1d235af9f8b2d58da711405dab409cf906984b6e7d627c314d10b",
    "google/api/policy.proto": "695cc0f9b612e43c9585aec68c70fe7a741026ade82ab288cd7bfb3905c7968a",
    "google/api/policy_pb2.py": "0bbee8a4cf3eb5bebc4caf45ed6b5864c80f5c778ff0da8878d2807b576ea797",
    "google/api/policy_pb2.pyi": "1551c3d4f80317028b3cc6dd458e1068c4c89593dcd3ae39bb478a5b476fee02",
    "google/api/quota.proto": "ad5e8a53037fd180bf646894f4e36705ac984f204dc2d5d98d8fe99dedb725d3",
    "google/api/quota_pb2.py": "6154dac337855ba1f95eec1370155212038b29d2dae7ee9847de6b4a56550d3a",
    "google/api/quota_pb2.pyi": "00a9b3e8cdebbf0d19c4c2c5468a4d7acb6c51e97adc3892074f6642d71b9d30",
    "google/api/resource.proto": "815d51eada4341bf01672b474b5ec0abbea190e8df9f4aa26a0c96dbdad2993b",
    "google/api/resource_pb2.py": "5347c627b6274116d6904f364f5033cd961bd15df62672d01a7cf8ee9d394f26",
    "google/api/resource_pb2.pyi": "5fa34f3a83374ef64a3a14380312b4bb959dd95fc32288dd199155d2d70b1b60",
    "google/api/routing.proto": "f31cb43c3efe07a80cd814eede587f1a38b4dac88053069dcfeaac86700511be",
    "google/api/routing_pb2.py": "f18bf2401c5378ebf20aed208b92255082ec297270b2b1c37156822c9fa34ca7",
    "google/api/routing_pb2.pyi": "da2cb3d48159433c9ec9eb83b8f4b216ad382dc96484ad8776100c319305a362",
    "google/api/service.proto": "91c47190669fb339437f09f911d2ea4307c6e99f5d66ab504c48d1a3d2274da3",
    "google/api/service_pb2.py": "6cab0441a013d6d8428165f0f956a031d9a5bb777eb1aa4a0dce851c94f89dfc",
    "google/api/service_pb2.pyi": "0f21f7895baf50eea89e38e1a1508a32610d8fff0f24d7b10d49e406b0e0598d",
    "google/api/source_info.proto": "8698efc898d714b206edcc8af0bcbcf80d538e7dd005fe206d2b922c314bab63",
    "google/api/source_info_pb2.py": "5ea7519cc90ff3c96c9af69112009136122425b4901d639cd8eb1dc5c49cf05e",
    "google/api/source_info_pb2.pyi": "a7aac39e942f4f4378218e20763e7624029db995e6ed1d8bcb2a3dea34252a1c",
    "google/api/system_parameter.proto": "b45151987677ac2789912bc8bbb270709b4f64a853383c65a133013e0974f793",
    "google/api/system_parameter_pb2.py": "8af7e4927b1a8fc2daa7498a2534055573a56bee70bc6dbdf1ad39090047a70c",
    "google/api/system_parameter_pb2.pyi": "ae40eb3727893d1b1f358e1c237755622f9d9b18dd62034af70787ff767958c6",
    "google/api/usage.proto": "26c382037e3e7c6caeba6d18329b6d7ddf37c1632ec5bcd57c5854d5be007bad",
    "google/api/usage_pb2.py": "ba36eb2aa424769964885c250309363804d01eb3e4a6050477f5994ef7373b77",
    "google/api/usage_pb2.pyi": "3807b8db25d595bc8ca344f7026b14a5125f894a1482e906c87606f3e6d1e379",
    "google/api/visibility.proto": "80f0b5c5174f0c73380e0e91e865e142133930c5f35b3c07c2a52ec6610a1810",
    "google/api/visibility_pb2.py": "7a727a43789ff1c598c32f713cf98ef413ef7a6383653acdd421f4dc8a344a2d",
    "google/api/visibility_pb2.pyi": "0a791cfae99885bbec3a3d53d6c736559d136c837a9c3e148204bdd4ff9b93e7",
    "google/api_core/__init__.py": "4d8d15b64318487e926fa0046ff21a855850f56ba1049db52158046ef1962757",
    "google/api_core/_feature_gating_helpers.py": "5de3b9b34323d94d63432953bb837a9504b43391ac29577a752ef44b2dbbf048",
    "google/api_core/_python_package_support.py": "7020e9eb6faa5fd7d3dcf6cc521bf24ab3c286a663ba3cd7898c61a7e207df19",
    "google/api_core/_python_version_support.py": "448c858012d3754d75012fc94d57fccf725e66c726761125a1202d03f67b0d4a",
    "google/api_core/_rest_streaming_base.py": "9d1c23dffb2e99be25fb9f7711c8f3c996396e30cfa85f3d4aff85446392c5bd",
    "google/api_core/bidi.py": "3b1074808d1fb633ad99318fb6937d883be26d73745fd75d4db19e76680c5b8a",
    "google/api_core/bidi_async.py": "f70134bdd6fa4078d431fcfe3fbb398b3f1346ec816cd223266bb86cfbffc190",
    "google/api_core/bidi_base.py": "0b9592b619c8def2b92be2edafafd86afc86622d4c0228d3c3e6589e2e552255",
    "google/api_core/client_info.py": "b0fb9a9ef7a25c403cd29953a350858eb938cb76f0568831561da78658ec7288",
    "google/api_core/client_logging.py": "0f02d0d9a0b6739d52e1ac68ba12d8b25d617ff7bb88dc1503ec801b1d89c79e",
    "google/api_core/client_options.py": "3ce6d6826af2fedff991f2611937ff4689ec6acd1ff52939a8a2be8d4fb49142",
    "google/api_core/datetime_helpers.py": "0cf57cce1ee3f920de6f90c5a61637bdf71501f13bcc3d6efd724f5cf608d3f9",
    "google/api_core/exceptions.py": "c61b9562784692668b7d9c51ef75fdc9bacbb4b74984bbecccddc26a789ee06c",
    "google/api_core/extended_operation.py": "afdc5239b94d177e65c27da1aeb8d443e7f7243a28d1ae19f31c0ecbf564bcbd",
    "google/api_core/future/__init__.py": "eec4e8c4d36ef5cff1a9ca663bc75badc48b38ec6996760e3925e3397f50217c",
    "google/api_core/future/_helpers.py": "4de7a7fbb512bf8a6ba22f5e825aa8d580400dcba0e1c23d233a96d29317d57c",
    "google/api_core/future/async_future.py": "73c57d039d3f5873a728db3c74fa5e4780283d32b1fb4540609dcc9fa22a5653",
    "google/api_core/future/base.py": "487cae75a992591ec4c94b18690f97ac61a42de6029527d77ec1c81d2a832182",
    "google/api_core/future/polling.py": "dbd5c7a4a547b16cebc33334d99f4451415eaf002e5d169e2c4ab8bc98ba591b",
    "google/api_core/gapic_v1/__init__.py": "30ea2dbca6fef8c6cf36762edc0794080ff70a4fc983d652170de8728b0fad66",
    "google/api_core/gapic_v1/client_info.py": "c4e22adb43872bc1013084ed19ff04440e6a59b7f6bd52030133c0e735a2505f",
    "google/api_core/gapic_v1/config.py": "125c1af35b432ee7cc79e3a9806e74ac939acbcc298302d47da2caea2172f951",
    "google/api_core/gapic_v1/config_async.py": "fe3ac1e58bfaaf1c5253a2b0ccec5643e1bfc799978a5a5214082743fea4b6b5",
    "google/api_core/gapic_v1/method.py": "cc63c8e8ea316efafa12ee08a05a3191aca9dc2a5c24c15b515a524efd298223",
    "google/api_core/gapic_v1/method_async.py": "41d6ad30583baefd922348ce7cbe29bc06a47e98caf55e337beec085b9d3fdd0",
    "google/api_core/gapic_v1/requests.py": "618895f3aa0e3c95ec8ba0acc67c700618bf22689178d7944739e240a3021dea",
    "google/api_core/gapic_v1/routing_header.py": "90928e629352da681265ae10b7c21bd90e4e35f370a49c1b365a1527c7b6c0cb",
    "google/api_core/general_helpers.py": "186830bdd60c2c8891bd4aa88dddee28130e3c82e92bbf19539c1ef95447b701",
    "google/api_core/grpc_helpers.py": "0bd33cbe79ab1cf016633bc2fc6f40fc35eb1e9a867c2853860b9ac6392f8b75",
    "google/api_core/grpc_helpers_async.py": "97641af5bdad4e32275b264033d1df304cdf20a24988a52f72c076a590d6ce7f",
    "google/api_core/iam.py": "046cfadc7b4e3f9ff9a07f59b3ddd13f463a42c86dcb678e8451188ff0a813ae",
    "google/api_core/operation.py": "50ff5aa0dd03acb2c55fe73b13f8a60d02d76149a82cf52fb5ef889938640041",
    "google/api_core/operation_async.py": "8d4aaa7c9113c8d2ed2e92b9ae25f25312debd0f5b6d97b690dea51a58e77436",
    "google/api_core/operations_v1/__init__.py": "ac25397df339c5b93e53574ba688e90937dfb6da24ee1ed9c4078e4d643d4a6b",
    "google/api_core/operations_v1/abstract_operations_base_client.py": "ef7a5da8f6d1e938b336e40ae709573ac90faf4a2e6dcee3b860c88ff3ab8c70",
    "google/api_core/operations_v1/abstract_operations_client.py": "798c643bc4203b2d00c465bdf69384ecab6125aea24ba08e06eb3ef7a9175b71",
    "google/api_core/operations_v1/operations_async_client.py": "6fabef8ec93c9d27e700892c970c3a256430e0540f8443ddc10ab0547ed7b750",
    "google/api_core/operations_v1/operations_client.py": "92c420e2ce6f7b9ae4b53b89e9533c79c9f0705ca9d0f4866569ea559b35d3eb",
    "google/api_core/operations_v1/operations_client_config.py": "bfb07416255ce69f47867a4f635ff7148a2615d03e27ee258a5a267a80bd4a44",
    "google/api_core/operations_v1/operations_rest_client_async.py": "f02731d40a3c6250bbc3faa8d2b3dcf10b577c3f37d34a66fb5fac63d5ea4f25",
    "google/api_core/operations_v1/pagers.py": "ba630f02580741905b4511c8c6f9c3ff7ef4b83790af7f7d1dcb4e92d0e7f65b",
    "google/api_core/operations_v1/pagers_async.py": "38503378c850c9190086c3c50aa75f2d8dab558bf105cf9d45cf122834fb2e12",
    "google/api_core/operations_v1/pagers_base.py": "c2e171dd4b48f9da6dca32a7b768d2fcbbe1c8b4a33b21167f194104856e8e2a",
    "google/api_core/operations_v1/transports/__init__.py": "8f0db214df3f9072cf5d9a7abfd4d31c0e7b7f39586b257a3942895c1b9a74d6",
    "google/api_core/operations_v1/transports/base.py": "7eddc783d258170c0e3b95f7d57c1c8741cfed39a3098ef762f4557717cc334a",
    "google/api_core/operations_v1/transports/rest.py": "2755d9793fc50a9b68fd8a7d1d2f0d151b1db2637e9d570cab0d4583bf871bf8",
    "google/api_core/operations_v1/transports/rest_asyncio.py": "3809b8367d470af657e1e0e7055d187df452ab3513fa26ab284b53ee9f1eca40",
    "google/api_core/page_iterator.py": "15731fa9b865558004563a688f2b58022525b95358015482e3531d7c0847017e",
    "google/api_core/page_iterator_async.py": "4dbb97a2b4613f5c1c4130f7ada04985681224fd49c093bf9c2289a6109b55dc",
    "google/api_core/path_template.py": "b3a0e009444f47e93dbeba61535e15f5ee84f0e548c71c5d6c25ab8539e696fd",
    "google/api_core/protobuf_helpers.py": "d3d0a5c97035bcb394f19d062c8303fd0a85c5c48f7fb476afc8be0a81ffcd8c",
    "google/api_core/py.typed": "abc7601fd9759a85178ae7c70558ea23432e272e017bd6b7acd1fc665fec2020",
    "google/api_core/rest_helpers.py": "9abdecef4d6d252db5636c4393866970130dcc5792c326779d2d53e5e08e01e9",
    "google/api_core/rest_streaming.py": "2920ff6e370c2c6421d3c7bde7f28cca078765d2923c1fece8ab12ddeb7d7ba5",
    "google/api_core/rest_streaming_async.py": "c5c18a1b37d385ffb39e4b9fb0d8335cb40e27e76f0234e0e86494c7885d31cf",
    "google/api_core/retry/__init__.py": "98f1d210239ad36f08332e60778d28f358059c1aaf155b6017ea79876b259675",
    "google/api_core/retry/retry_base.py": "b4859c1fa97416d79c874ef2b8cd90984a782605e91ee0ff0a17c466681783c9",
    "google/api_core/retry/retry_streaming.py": "a3219d4579f6d2f168a7f6036d2599a46565fb6bf6afe77e0596f514eb6d3430",
    "google/api_core/retry/retry_streaming_async.py": "6c1db01ed5fa56697f755431d637dc3e93695b2bab3d4573dc50d85f1ff55bbe",
    "google/api_core/retry/retry_unary.py": "13eb82696a60daff4b82e1e61da48b66bcd13690f068825c6307265c93355eb0",
    "google/api_core/retry/retry_unary_async.py": "f36bf9de68e2d63cea37485fdfbfd32bb9b2351687c2e35384db359b6ce75c36",
    "google/api_core/retry_async.py": "51378632c934c03041c030f1a989437b05b6f94df7bcd4532e4db7abc1690cb0",
    "google/api_core/timeout.py": "85e8a5d04eac72ec8590cbf299b4766c0df76662526af1ff4a644d2bd929a9c3",
    "google/api_core/universe.py": "ab374611acd8bf114eb1c0293cac42f1763dbe7bf0919762a3fd730cf44ce126",
    "google/api_core/version.py": "9ea4ffe489d0f6b03b7b684265e3b4aa4c9e5344ee2f3e3cee8807b214234b9f",
    "google/api_core/version_header.py": "b84157a2c0a9f141fb5e1ce71b9190b1e4d8b563682475d13c0e79ec35ac5310",
    "google/auth/__init__.py": "02ad48cacbf1be3e6ca9bb0d1cf151e4231b049b3a2422820a69f926d35a2906",
    "google/auth/_agent_identity_utils.py": "e595025bc2b1654ecff485034a6d8bc95222fb6f87e3b17e30ad7168ae28f4c0",
    "google/auth/_cache.py": "f716ea764c056bef56e6dcdbd9633f6cdcc5c19ea129accc72c3ee82a1ee0320",
    "google/auth/_cloud_sdk.py": "a705fabd9858a84a6edef912371276308e70313e6dfb15fbdbbdd0a82b561599",
    "google/auth/_credentials_async.py": "37d555d9ca55ec21cda62e86b143411ab66e5687879b1355d1bdbc6609eee3c9",
    "google/auth/_credentials_base.py": "2b1742672a05cafadf5a125b367b98929521c6cd1b99bc58bdc1fcfb1a79854b",
    "google/auth/_default.py": "16ac217251c280a359cf5a7e94c444bfa998f652939f695ee823f17241a2b719",
    "google/auth/_default_async.py": "66453c8e283f9446ff229b184fc0c1ccd3cb44681f60d6ff41b2f44b79698e4b",
    "google/auth/_exponential_backoff.py": "ab103d7a4f34ac1928011c741209766f5325614d03fa940d56a80299487e9605",
    "google/auth/_helpers.py": "63b585e8538af7632baf2bb46f781245ea86b5708582bb71b755c97b81075319",
    "google/auth/_jwt_async.py": "865886f2a62f5881d174686fcaf12bf2cbbb20a7e10d1b1c0e0406df852822d1",
    "google/auth/_oauth2client.py": "f69052864d5e93090f3680f95a84093501c70f796e35fce9055a10268f27d6b7",
    "google/auth/_refresh_worker.py": "cd3da38fe3e33ed8c5dc1aa0fe7e2dd4274bdf6c815445213ca650c2836ef978",
    "google/auth/_regional_access_boundary_utils.py": "6510d3566b1fce95811be2cb63a936433fa7f449bb6ce366d381cdfa6bfe868a",
    "google/auth/_service_account_info.py": "11d38466a97dfe61b6709e782a62f74398cedc2df071a77094239b83880a8eeb",
    "google/auth/aio/__init__.py": "7b74e80315cd1e1a892c1816f01ebaeb9d3176aad30990cb7303f6a790e108f3",
    "google/auth/aio/_helpers.py": "466cc1671ea7e2f5a706ad6985529022152a0692206049f63108c442c3e81373",
    "google/auth/aio/credentials.py": "957634fd227d737e8cce9e3b7cde1af094f91df373150b48230b19130ce80d1e",
    "google/auth/aio/transport/__init__.py": "2f71118c9a0c2af3cce3c0ee2ee10ca4c731ee1c2ee63250b9865a9529a8fb9d",
    "google/auth/aio/transport/aiohttp.py": "b9e76955b0ced502eccb0180c9ed783da73ef6fd230add79d54b03662a83678e",
    "google/auth/aio/transport/mtls.py": "ccffcc5e98f0c4034b0bfc20b0cc84055a696d1f25136e6d28910ef975f3812d",
    "google/auth/aio/transport/sessions.py": "d467113a58f8f8074117dfc92f7e0204d15439bc0276deaa0567c3836d04a56f",
    "google/auth/api_key.py": "3de89e4d871e1c92050a8d33428d4403d3440cbfc87ba4bbf2a983fba220d7bb",
    "google/auth/app_engine.py": "990ecdaca4187772e267b17fb6d96ab3d4e37f8104b493474cf7e98b10b2716d",
    "google/auth/aws.py": "3e8fbc759e67342f9f5d507f67bbeaa27f63f367799f645c05ea86edab912be8",
    "google/auth/compute_engine/__init__.py": "06a79391afa8c871404e4cacdd218a4653b2590f26555d2f55a3f684ec15e10c",
    "google/auth/compute_engine/_metadata.py": "4f638d21d33aeb3ecbbd5becf7601113bea115870010a0030ce70bb8639a4917",
    "google/auth/compute_engine/_mtls.py": "ab6d34c553c3a3aad1dc646774e1ec6623421ca95f82b18b91b3211c933e29a9",
    "google/auth/compute_engine/credentials.py": "0b7efb91f6731b78080b4c14d7d4433ba423fe356ebcdf24c0106be792318035",
    "google/auth/credentials.py": "53a1ff1fdb2877110ea7e7182398d2b63a5151a6c5c3fa489dae2ca385d2ad9f",
    "google/auth/crypt/__init__.py": "ba40e62ab754987a86739ac32dca1e97f4ca5dd579c17d39af91aa8bd8675811",
    "google/auth/crypt/_cryptography_rsa.py": "a364134647c344bb44062abe7dba5b4d6ca9bf16b1503c33971d8da571bda34c",
    "google/auth/crypt/_helpers.py": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
    "google/auth/crypt/_python_rsa.py": "131ab3f4e8039bedc91c777adfa82701d5dd707ea98c3581349cf6d6e0d19d41",
    "google/auth/crypt/base.py": "dc226b9d0b29a51ea1f969d346a4abb7584f1105641dcc0d21527f079334d216",
    "google/auth/crypt/es.py": "dde1d20d4fe587de69318137788184ac61cf113aac4421de2b656d427a953e4f",
    "google/auth/crypt/es256.py": "705683b9045d5f1ddffda572f871f4550e86e3d54d480841d388e76e3cb11cd2",
    "google/auth/crypt/rsa.py": "7282ebd05172b830260daa1779fa6b887fe768025b9a0c0539150355819d0585",
    "google/auth/downscoped.py": "2a6982f256c1594154b0862dd55a5c6d35ecdf2609617b5e1f5f6a759420a1b0",
    "google/auth/environment_vars.py": "e75987f6c33fc4ab6bee0980ce32323d67a8ce237cfab7642dc981cf38cfd328",
    "google/auth/exceptions.py": "f2811e07fd548ab248c81a209504690858c239cc5c4da2ebe84a41aebc63784b",
    "google/auth/external_account.py": "da394a2a19b515bcaacb2f3a7d4cbccf14619e4655ed3327e1bfe4d4b8776c28",
    "google/auth/external_account_authorized_user.py": "91f0b45190694ccd4d4321b75903839dbf95edeaa77c7caa11e15e0126dd76fc",
    "google/auth/iam.py": "6fca0be6b5a78676fc4b86ce0b1916d3ebe576c28ba285e2315a4d7281ef6a15",
    "google/auth/identity_pool.py": "f4db4a5d4c69a9f69a43b84520f0ac6cd48deb59581e238c8fcff1e774605051",
    "google/auth/impersonated_credentials.py": "01dcbf7ec35aac9368687022e48e29001baddfefde04d160b927c59b99ad4e83",
    "google/auth/jwt.py": "0a73deda85e3fd4a6a2eceed146d924eec03b69a5fd26c4435f154a32e842097",
    "google/auth/metrics.py": "71a6c4364840d30e1309363881a6b702968661711bef68ca30c2f1b655705754",
    "google/auth/pluggable.py": "37d97122289e813e8322fde9a1ce41dc1fba3a556615166bf8098d3ba5f389de",
    "google/auth/py.typed": "974e7f2d3822de8cbe7a29ca06bc3aeace9a6af8331b6a3e49e90a3ce6376b23",
    "google/auth/transport/__init__.py": "bd37565030425d617ac24dd7ceeecbe48d656ed1984f138d8be20a60e9837453",
    "google/auth/transport/_aiohttp_requests.py": "d64d05e9f906f83fe71d4d06015a3f6abc2bd60539de079ec35535ca0deb1e40",
    "google/auth/transport/_custom_tls_signer.py": "8d65ea367d9cb62ba5bc428d553d1cbe29c14fec175fe9b65adf009bc61f1a27",
    "google/auth/transport/_http_client.py": "fe51364da40e6f6859ab93ab6e534061ba59ccb5f5f0b1974b268d26e5c7276c",
    "google/auth/transport/_mtls_helper.py": "2bae0ebaa9ec2b31c86a9030b1930a66d56dc39dc8e20dcc0f4b9798dad13d01",
    "google/auth/transport/_requests_base.py": "e32d2d4cc47f84f19e026052c8253afe6a21f7d671259e0217bf44dacfaded30",
    "google/auth/transport/grpc.py": "1111107f11d9fd15e697678b5170fd4f4cc4e8dce65ef3e0ee208797767d7445",
    "google/auth/transport/mtls.py": "e1e6c1ddc77de2001b1a8f481ccaca1ca47a540b812c4a954a60d0dc012f0d34",
    "google/auth/transport/requests.py": "501eed10c1c1aedd4d76e9189eb67703aa1c68c231c781cd92e9178f23936e77",
    "google/auth/transport/urllib3.py": "2d7f08139e0442e8c230e266e4fa93bc721f205b3f9dda441497720f94e2aa77",
    "google/auth/version.py": "ef21f7c7a3d5dc3cfe4d4241f0a335ed21e98573db0231dde2b728d3df51187a",
    "google/cloud/common_resources.proto": "ad34ac9832a84e18a8b293ad4472d95f6ed6edc04f084252e6ba552b56cff58d",
    "google/cloud/common_resources_pb2.py": "99623d7e5e94aa99dd973bff7ae8e39ac2c2ecc16bcb33e96b9a0ed0998af627",
    "google/cloud/common_resources_pb2.pyi": "6ae380c84f3e016bec63daaee71857d1371a1967cff6f3ce51cbc45a49278017",
    "google/cloud/extended_operations.proto": "808e252f9f10670e0938bd35fe3e518ab939034d1fe8c1b43f49b9797f72b523",
    "google/cloud/extended_operations_pb2.py": "3cf4773fe3ecd577071ef20f9afa059fa75b455f3a70a785ea91bf9da7c377f9",
    "google/cloud/extended_operations_pb2.pyi": "ece2c6223aeceefc9dd186516db1756e6978245f6149da253021274775d213dd",
    "google/cloud/location/locations.proto": "40fa03c073fe86d9577ce97d9ad9b65f358e997307cdbc57d5f608bdf01e02ce",
    "google/cloud/location/locations_pb2.py": "52cb6acdfd1f8e673a9607389c65b37c6f012e42b2b4110287d909f71b0c01a6",
    "google/cloud/location/locations_pb2.pyi": "9e032ec307ffd03bd272928cd36693a02b144ad0e7cf88068b16101063acb2e0",
    "google/gapic/metadata/gapic_metadata.proto": "fd5c638ea297a0be176053810500b021b210ef6e200fa48e47e009053c0f9386",
    "google/gapic/metadata/gapic_metadata_pb2.py": "e6d58190e22214b3eb53202f232a117e6859ddba42fab551764fbdeec0a5ff43",
    "google/gapic/metadata/gapic_metadata_pb2.pyi": "5e95872e1f74708d2d72886268a8f580087e1707e76078453cd821c2e180b135",
    "google/logging/type/http_request.proto": "cf16052273c8732143ad6418c0de059aa49a3b2f402f62e6a7c63491436cc465",
    "google/logging/type/http_request_pb2.py": "20177c1873bc9b5850f0060af8af47aea949660d4634e7a8c1c13764d1740c2c",
    "google/logging/type/http_request_pb2.pyi": "5762071412444bcc54f1faa2bb9c5b193055dfc2c7bb1b260d36e062adda4b66",
    "google/logging/type/log_severity.proto": "cea9c15d29125773c61b4b018ddbc9b3d7f088a2f737d6a74e38f1b607dda51e",
    "google/logging/type/log_severity_pb2.py": "14591b8e633baa54cfe8fc5712870d1e25d7f61ab957b3b1648623afaafcdf2b",
    "google/logging/type/log_severity_pb2.pyi": "e3b8b4f354af5feb7b86a959ce6cf2f0d663ece2b79b27a835eba95421d68bfd",
    "google/longrunning/operations_grpc.py": "66d4deffb1c9b03586e0e66897c422d2aa15ec790ae01031b7c561d1587ffb1b",
    "google/longrunning/operations_grpc_pb2.py": "8a0960dfb4532c85ea9324287cb882e3226bcf2c1185794f1346a9bddd234ca3",
    "google/longrunning/operations_pb2.py": "211129cd3556f74a12786896ab4b11a20b5782207310bf2432ec0bf170b40aa7",
    "google/longrunning/operations_pb2_grpc.py": "0c543bc76e713c4afd5454fdcd255c67682a3ecbc6fb9c614067cb9b8991bb7b",
    "google/longrunning/operations_proto.proto": "fcf560f3b8b9d1bdda82803026ac8caa616b0ee19d0f554e8eb9ef7b2f201fc2",
    "google/longrunning/operations_proto.py": "bc87bb4685597e3b12ecc499078f3b18630807fe7fefefac4ff732b7c7281829",
    "google/longrunning/operations_proto_pb2.py": "f164f088ffcc97deacae01dbde445bf068aa85d05adb77624c2c7bf58978f638",
    "google/longrunning/operations_proto_pb2.pyi": "e7f57dbef6dda525d7744ca7d59a013693e359ad1b157dc257ac2d19e543deef",
    "google/oauth2/__init__.py": "883c93a71ba1f3ae2b2dfe18c8836078ee303092d0e7711430f2698b0d460375",
    "google/oauth2/_client.py": "39230518c764aad2d8d2d178be6b49571b81c2ef386543f310869ad017071417",
    "google/oauth2/_client_async.py": "e17b568ee00abbdd9cfb86bdc7e5503eecae639728e820c6be655dca7e8ea09d",
    "google/oauth2/_credentials_async.py": "854aee71091c62e625c8274731c8bcb736959dc9e342515cdac01f36ee43b7c9",
    "google/oauth2/_id_token_async.py": "a2de3fbbeae1036d93b9f58e011e4ca5ab9aed0d9a02d378abf9099a3f989c82",
    "google/oauth2/_reauth_async.py": "5caeb7c75082593057e131bf9b07355a7466077bc6d09459e45f9d35fc6875cb",
    "google/oauth2/_service_account_async.py": "8f497447100637f5b5a6042409898bc18f0282ecb0ccffc81743d009daad9484",
    "google/oauth2/challenges.py": "3ddcbc229c0626cb73f3dd02b73ff166f854c08953bca89b4a798735b44165dd",
    "google/oauth2/credentials.py": "2d08c8e5eed77bcadbbcb3086e491e1a2fa2962cb66ee32832d50461c25d8e00",
    "google/oauth2/gdch_credentials.py": "098ea23e73ee73638221ed59ba3c20d4cbbd41297588626a18ecba4e4525787c",
    "google/oauth2/id_token.py": "744523f1e16f1de9cc0b530bf76aa851e61b5c15bd08a659f872e57fad98eb9c",
    "google/oauth2/py.typed": "2349ae5d145d6dda49a19fd5ca17a2b224d361c9a34c08ad42936fbae87a7ccc",
    "google/oauth2/reauth.py": "0a4903a6fefdffc4ac403ae8980c88dec6b6952dd5fccedb499a961598f7d27c",
    "google/oauth2/service_account.py": "4f0604894c3f938319baafadf262d30294d54c9038d152f97f8aecea4c3e19c2",
    "google/oauth2/sts.py": "75ffa5ce17bc7012abd8efad5994da59d818a272a6d49f472390b3fc07d6d615",
    "google/oauth2/utils.py": "e1cac076929b0eda1ba507d725cdee17a666e85dc8cdf7ef452a3ef61ff9d79c",
    "google/oauth2/webauthn_handler.py": "0efe2eff917e786599997273e4faa05cd8a715db7b8e746a4413447caedd95ed",
    "google/oauth2/webauthn_handler_factory.py": "b28139728919de92cd6c1a351c2e85d4df8df88fa2afe0c6943d2dac63c2a41b",
    "google/oauth2/webauthn_types.py": "20776a51ef8458e53e0a6c8b9c72e78aa73bbc61d9f07a5abc1d1d7239fb665f",
    "google/protobuf/__init__.py": "5edecc0aade5b8251832a6160b0a7301aef4d8314bbe805c19e32be69a484d5b",
    "google/protobuf/any.py": "d563db4061f17713e176318dd400f0033ae39267c4703ce9900b51e7b98ed380",
    "google/protobuf/any_pb2.py": "5190875f855fb1118190d01dd49e8429065aab3bf3cc59601ead67ca91b4f139",
    "google/protobuf/api_pb2.py": "b1d2f62f9432532531b9257648dd5173dc5aa914e0329328719f2d1d88ede712",
    "google/protobuf/compiler/__init__.py": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
    "google/protobuf/compiler/plugin_pb2.py": "60376be07b552f125f08417d17b78f62771fffaebfb283c8b323c1ac2520fd8a",
    "google/protobuf/descriptor.py": "9b1239b4ecb00eac2a12b4a1d15d90a72aa3856420e1357ee64a16465074b894",
    "google/protobuf/descriptor_database.py": "d27fb0f8d420add42231e67a2e70b02c08b0a281e3170c0370956916655fa2e7",
    "google/protobuf/descriptor_pb2.py": "a3a6d0cfdfeec9cc74c979bf2ea3c39fee83894d3e26737afc842f7bafaf2434",
    "google/protobuf/descriptor_pool.py": "319322d01b4bde165221d4323f3fe356bc50e4081adae8c873ec3c835bbc7e9a",
    "google/protobuf/duration.py": "bd04f05728a2c869b75b2dcb5bca21037b641a4ad42a84e7fe9e12744054f1b3",
    "google/protobuf/duration_pb2.py": "7ae170e8a2c5f49bf3d296dd608662d2f89163cecfb12be978bf63210ccd55eb",
    "google/protobuf/empty_pb2.py": "e47a0eb864009d011dbc2154a4835b341bbbb70f3357e13b754a1bf31218d1ad",
    "google/protobuf/field_mask_pb2.py": "6bcfca2dc5152956759dfa3a2cf0e2a4037c232af63ed6935d3daa92bc7dc5a7",
    "google/protobuf/internal/__init__.py": "b93e0182c8f2283152dcbaba56bb4000dfa54ac7a17c9c60d600a55f94e164d1",
    "google/protobuf/internal/api_implementation.py": "e651f709c17bc034344726a3a6863bd404bf8fb254c4c558a50a1ab4e2e1be2a",
    "google/protobuf/internal/builder.py": "a6dd21583d41378032592f7c51ffa6da6a8eee25e81254cc7ecf6b821d53bff0",
    "google/protobuf/internal/containers.py": "327ba2e1cea0304321ee7436f2f9196e9893d26646219d5bc22fd11883df25c8",
    "google/protobuf/internal/decoder.py": "16cd39157aebf77f020fb635a2f65bbb059aa3e07e0ad60aebb0d20553eb2806",
    "google/protobuf/internal/encoder.py": "54ab81264e52c8c779c601238936c6a0414e739c8c1a37c597bfcb7a18dbdac2",
    "google/protobuf/internal/enum_type_wrapper.py": "d5cab3075b586cd12cd2436a3869a6bc220002f16ad35749dab85ea578c075ef",
    "google/protobuf/internal/extension_dict.py": "82d7c196c76fbad10a83c2c317a04dd7b103e8b0299e1f7793ed58ba1219aa88",
    "google/protobuf/internal/field_mask.py": "651c7fe7460883c914cc3e83a91b99690dea90d2e6474c11702964f28ea34788",
    "google/protobuf/internal/message_benchmark.py": "c19b8e9a048b5c30c8c37a15497a81b12660d7e6bf8be20ea619f8ccd61871dd",
    "google/protobuf/internal/message_listener.py": "8404f4120a76ab5092d2e8a62831d6dfc3525109aa3c4d8e876ee08d6ca1bff4",
    "google/protobuf/internal/python_edition_defaults.py": "7fb7071b6bc4845347c1f3312dcb8de0451949a476cbff94af7a793fa3097118",
    "google/protobuf/internal/python_message.py": "7e882bd51ddaf133e1d6fd23a63bf485feac3ee28464dbd6bc50f37b23cc2fd1",
    "google/protobuf/internal/testing_refleaks.py": "39a36642b2e695cc1979104af5fe1718fe2daf4816faf849ed644b15ef4f5666",
    "google/protobuf/internal/type_checkers.py": "4c6a3fa576c81227095156f88da878f8c406b29d40164d13abff0e16bcabd3e4",
    "google/protobuf/internal/well_known_types.py": "a0fdac8e76f32657f6c414c51b2904e062bfab8d98c6e68911c51b86492d769c",
    "google/protobuf/internal/wire_format.py": "3aa64c99d0cea9579aae6088b18b5515f7511006f33f0574e6956456804cc167",
    "google/protobuf/json_enumvalue_options_pb2.py": "6dffbc209c663bdaac1437bcedf9223605db71977f520a350af921dc0f2e760b",
    "google/protobuf/json_format.py": "5c5a6efc7d58494c2b6e5ac437bc18ef3bf92bf06e84d7b50bfcab9ac42c2f88",
    "google/protobuf/json_options_pb2.py": "d337328744d4ccb91c04a04b3ffc2f5b4d7b6a305b46f83bb9d83cd318fc9a2f",
    "google/protobuf/message.py": "65e00c107f9422abb2df305e3dce4770b0b0234ac8f3f55b56921b925f57f700",
    "google/protobuf/message_factory.py": "5b8a8cf7d0d9c07ebd4fb7b4f6b00c24445239ead246583feef17285deb22c6f",
    "google/protobuf/proto.py": "720fe4f668f52787069a7b7df9ee2b167937c6a89a66a43327c0bc0a8c32bd57",
    "google/protobuf/proto_builder.py": "ee9417fdcc246cf33e63b3fa5ea1f4390cf896a54ad135f303272b061ccbe05a",
    "google/protobuf/proto_json.py": "cffa128b5ca72ff42f71c20aff6c8f117327409b04c088e8b4a40212be33ffc0",
    "google/protobuf/proto_text.py": "81a9191ecd194ac078f5fdd8c8d16a844ae3f3466842cbd680c9e8bfa62c5486",
    "google/protobuf/pyext/__init__.py": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
    "google/protobuf/pyext/cpp_message.py": "33215a6e51a61422474d204ea2ecda38a7c62509ce8de18e249fcafe8ef78230",
    "google/protobuf/reflection.py": "74d1990d195b4fa12edbb081ec1bd7c5756678044fee74482caae242c036acd6",
    "google/protobuf/runtime_version.py": "83c2ce094f32e0b39db99ab3d139b5c8fa75dc676a7e54d5fb9218f16b018498",
    "google/protobuf/service_reflection.py": "a98d7c8a8ac56e7cdc61a6d0b9f7ef57876c0f16a120d40fafdb55daee46266f",
    "google/protobuf/source_context_pb2.py": "ca9f75fec9f7ef27908b785cdd0b354e236c714a187ccfd76d3f7568b35d9113",
    "google/protobuf/struct_pb2.py": "740b29b21517cceed46564e3a4bf244d2044d13b59d1a96907dbf0eaf4410618",
    "google/protobuf/symbol_database.py": "61da9188d8be686f547289473b39aec144fff21ad69b785f324ef7adb6176928",
    "google/protobuf/testdata/__init__.py": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
    "google/protobuf/text_encoding.py": "9a76280ab63e5cdac1d8d21d56e1c9261633835501bdc5e70f25f556ac5f3825",
    "google/protobuf/text_format.py": "7ac692a079a16b502cc17519591e234c5fa4b66f1ba6b17778100834e20ffcf8",
    "google/protobuf/timestamp.py": "b36dcb5aaea10e2148780b55527f0bc1f11ce5a44ced60304f3fe109a3959dd9",
    "google/protobuf/timestamp_pb2.py": "96793f6be608ed35e2414e49bf51655c4c53661082518e8ddd762bb9bdf3e0b8",
    "google/protobuf/type_pb2.py": "1606556d6519d54e25f8f5d76e2208570011fb94f315ff6d68b4ae58919496c7",
    "google/protobuf/unknown_fields.py": "c13630d365876920b97c37e3c394e5f5cda9a325903d0a81151620169303542b",
    "google/protobuf/util/__init__.py": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
    "google/protobuf/wrappers_pb2.py": "699e7d3e601fc32dc5810384793257d6cd52e15ff1044ef9a707067daf7f3f33",
    "google/rpc/code.proto": "83a63185b83cf7332196fc5ead0ec91b93ac207fbe91b3d5bc3d72f461b0bb31",
    "google/rpc/code_pb2.py": "2ed7928536fe2e411e29627ef3712bf43687e593f510fed14c909c5900f92c1f",
    "google/rpc/code_pb2.pyi": "4b4b55b0566beff417a4b521464442a838f712868d00ed119f9155ed25c86aab",
    "google/rpc/context/attribute_context.proto": "8aa3f0fe9f4872be0a32a1ad4f6203b18a67b7a4bb69fc00a4652550b455a93e",
    "google/rpc/context/attribute_context_pb2.py": "91f53c65c405fd890ffd538c214ad556b698e2fab1ca0d73307e6bbcf4bcb8b9",
    "google/rpc/context/attribute_context_pb2.pyi": "a9c2c565ebb3fb734e01870abc7275e7ddc6e8b18b9d5d8c7e23ada9c17c5495",
    "google/rpc/context/audit_context.proto": "cd1a53c653cbd626aeb60a226f19913bc975277153e1b5dbae5a3ba43bda045e",
    "google/rpc/context/audit_context_pb2.py": "6aeab8b3dc19c8eec4113a5e7876e5320479b554f7855ba196b8d7733bcdc948",
    "google/rpc/context/audit_context_pb2.pyi": "19f26d0d04dd47aa9385db38d63560823a628af658e881666daed265c16f52fb",
    "google/rpc/error_details.proto": "5c18d42e6924f447622534b4b7cdb8b53c0606c8ce729b912ee5edf88d1c5520",
    "google/rpc/error_details_pb2.py": "03566335de9750b768b8a5a73194d6cd0ab1914debbe9a021f2528a63465cc76",
    "google/rpc/error_details_pb2.pyi": "49eee5127b1e6af5fd4df9f6b97578a04d2473b5501913ecdefd58d969bf9b85",
    "google/rpc/http.proto": "f924995faae52eb92c09d7ea25d9b217fe9a662954c06398e75d5cb4e5ac0386",
    "google/rpc/http_pb2.py": "4734e24fc3f95d64899bcb600e686c7f53ecbbcc2eb05a76ead98b903f299a00",
    "google/rpc/http_pb2.pyi": "efa82022801b41adbacc533e3288cc1b0e97aec051a7a32664bf9636be38bb79",
    "google/rpc/status.proto": "f5bfd262e6705c7ae73f32e0ad8ee20ce8c0a2578df8c4f76ebf76b572f295ed",
    "google/rpc/status_pb2.py": "0da1265e75ec9c52d995e5c4893a5f43f30e7a233f77a559687aef0a29318cef",
    "google/rpc/status_pb2.pyi": "9091947bee80d6719244bcf326daa35171a947e7a8df10fbd0054011191b8e02",
    "google/shopping/merchant_products/__init__.py": "f48f938229d0fd967fb44715b4d8fdcfb0b30723d8b90514b3409f709c43617a",
    "google/shopping/merchant_products/gapic_version.py": "080ea6686bc6d05ada44661c235064ba39e282deef5be8b7b85822be4f28bbc9",
    "google/shopping/merchant_products/py.typed": "5848b8880c6ab8acf26d4234723be90642d4b640dbcba98b9344649c7d02b948",
    "google/shopping/merchant_products_v1/__init__.py": "f4e0126e2cfee4bcf9a10f19c26931c62d56924b94db7f1156b81342635d2de0",
    "google/shopping/merchant_products_v1/gapic_metadata.json": "4e30a3f1e387facdbb487d36a1553a346692d2c64b48f3b13fad5c72fa53e8e2",
    "google/shopping/merchant_products_v1/gapic_version.py": "080ea6686bc6d05ada44661c235064ba39e282deef5be8b7b85822be4f28bbc9",
    "google/shopping/merchant_products_v1/py.typed": "5848b8880c6ab8acf26d4234723be90642d4b640dbcba98b9344649c7d02b948",
    "google/shopping/merchant_products_v1/services/__init__.py": "02edd8492b7d3d907692bbc187e7775581a62d1e0422577e67132ea1dbedf119",
    "google/shopping/merchant_products_v1/services/product_inputs_service/__init__.py": "71e02bc273aee73bfc21e8ab73d14479fb9255577e96697060ee5854ad787c09",
    "google/shopping/merchant_products_v1/services/product_inputs_service/async_client.py": "e3e482adfae67f86e94cf7a27d9db508c56ad4834896f2d1706895f7a80af948",
    "google/shopping/merchant_products_v1/services/product_inputs_service/client.py": "fa0372845ab2ef2edf48b70a947c91a6df00cf70b8571afb068ad0bd1b1ab025",
    "google/shopping/merchant_products_v1/services/product_inputs_service/transports/__init__.py": "83ecb126f12b8c0b0d4c985bfa3cd24b67a1a23d1a602d45c71a101249c4a5c0",
    "google/shopping/merchant_products_v1/services/product_inputs_service/transports/base.py": "3b5abb530c2b79e63282df77ed2d4395024d45b69124620a3fb9c4ff3415cf2d",
    "google/shopping/merchant_products_v1/services/product_inputs_service/transports/grpc.py": "187b8d97ec83d1cf2e996a46dd6d8df2c68014d3acf42492093626027f4abdbb",
    "google/shopping/merchant_products_v1/services/product_inputs_service/transports/grpc_asyncio.py": "508d2b62001a70960fc4ce120499ac81c43bc016209d77d3e6881e0379bfc687",
    "google/shopping/merchant_products_v1/services/product_inputs_service/transports/rest.py": "ae6156cabd8783524533ab89aa51d264ffb8c298c632fca9fbff8fe76bcdbd70",
    "google/shopping/merchant_products_v1/services/product_inputs_service/transports/rest_base.py": "4bfec917d573605f96da5f256e5255f5bec09187a103ca759909ae3a9461d8de",
    "google/shopping/merchant_products_v1/services/products_service/__init__.py": "5453a2835946c86ea7fb3e2949571c237ca198b484c7ef1bc098acfd65e499f6",
    "google/shopping/merchant_products_v1/services/products_service/async_client.py": "7787c3516671c42e10b65af5f59b66292fe69244ad336d1afa04be113c90a621",
    "google/shopping/merchant_products_v1/services/products_service/client.py": "9267cb76c919ba8541f941be5800837c8f94d4a7fff7b8faf719a4a86c169216",
    "google/shopping/merchant_products_v1/services/products_service/pagers.py": "2e8febe145a1843b88baef0ae2196cddd947bd1e6d14bc1022d8527ac0d58a43",
    "google/shopping/merchant_products_v1/services/products_service/transports/__init__.py": "d458f904826a6d2b81031a2448cfb45eafed72a332adfe2a32979193b16d7c0d",
    "google/shopping/merchant_products_v1/services/products_service/transports/base.py": "e631be028d1300d505f177de187a191541f7d58da2c4a7179f27a7d87cf9ef9d",
    "google/shopping/merchant_products_v1/services/products_service/transports/grpc.py": "11044e2aa763cbfb477bb907663766bce4455b150ccef59c0d6c913635a0bfa5",
    "google/shopping/merchant_products_v1/services/products_service/transports/grpc_asyncio.py": "bd032a9aa7f0df54db07b0dc1192be90202fa1573378c08f8413dc4bf0a91c43",
    "google/shopping/merchant_products_v1/services/products_service/transports/rest.py": "f1a583382ca4618fcf2624c80111e5c2b04e0fab8e4a0ca376b8c881ac40eb33",
    "google/shopping/merchant_products_v1/services/products_service/transports/rest_base.py": "d38e1843fed748ba2f7db2c2ca377064a47da1050af78e4ab321df3320956a43",
    "google/shopping/merchant_products_v1/types/__init__.py": "49dae89f2ecfadf0b00300b1336c876cc5432675124efc04a5e7d2b5e3ce3570",
    "google/shopping/merchant_products_v1/types/productinputs.py": "e044c37a8077c4cd6a8f80f1ed3754ab96a91b5527a02d7131fd409ea8f57192",
    "google/shopping/merchant_products_v1/types/products.py": "e87a5bc309fca54521fdbfb75eac3366c9196ec86f7efaf7dc232065d87fbaa1",
    "google/shopping/merchant_products_v1/types/products_common.py": "bae83969749e8ad881107deeb54dbb20c124dc499f76a2ad72325fbe19a0c778",
    "google/shopping/merchant_products_v1beta/__init__.py": "dcbed57bec95edeeff04a375820d74333ca38a354949a914ebdc6a2b77eaaf25",
    "google/shopping/merchant_products_v1beta/gapic_metadata.json": "f04782679308c0a8be96eb40c79fd5295cbff61841dcee7073d91af5d3e864ee",
    "google/shopping/merchant_products_v1beta/gapic_version.py": "080ea6686bc6d05ada44661c235064ba39e282deef5be8b7b85822be4f28bbc9",
    "google/shopping/merchant_products_v1beta/py.typed": "5848b8880c6ab8acf26d4234723be90642d4b640dbcba98b9344649c7d02b948",
    "google/shopping/merchant_products_v1beta/services/__init__.py": "02edd8492b7d3d907692bbc187e7775581a62d1e0422577e67132ea1dbedf119",
    "google/shopping/merchant_products_v1beta/services/product_inputs_service/__init__.py": "71e02bc273aee73bfc21e8ab73d14479fb9255577e96697060ee5854ad787c09",
    "google/shopping/merchant_products_v1beta/services/product_inputs_service/async_client.py": "465e3db26639a14248248cfb866e399980113455cd3466df268d4b58ac9bd575",
    "google/shopping/merchant_products_v1beta/services/product_inputs_service/client.py": "d142b3673aee9f96dd047e802f61480664813610c3e3a780bde45c31774afc0b",
    "google/shopping/merchant_products_v1beta/services/product_inputs_service/transports/__init__.py": "83ecb126f12b8c0b0d4c985bfa3cd24b67a1a23d1a602d45c71a101249c4a5c0",
    "google/shopping/merchant_products_v1beta/services/product_inputs_service/transports/base.py": "a8a15e865007f5413c94c74648cb520a143cf5fa1ee56bd8b593b2e50254c46f",
    "google/shopping/merchant_products_v1beta/services/product_inputs_service/transports/grpc.py": "cfcffeae931fc7c8d5238d50752283c1f45b6a8af39af521b606c26288f1de86",
    "google/shopping/merchant_products_v1beta/services/product_inputs_service/transports/grpc_asyncio.py": "7a28529efcad781af3b4f584effc841c9034ac2df4ec688fff07a03e895e2cc6",
    "google/shopping/merchant_products_v1beta/services/product_inputs_service/transports/rest.py": "346edf43a8921b715db47ab4806739474adb96ea1e3c5492652b0d7127b3d15d",
    "google/shopping/merchant_products_v1beta/services/product_inputs_service/transports/rest_base.py": "c07b8ec26f2e7fc4e410cfd9dc5deb5fcfec3bf4861234a631792dac39c55f63",
    "google/shopping/merchant_products_v1beta/services/products_service/__init__.py": "5453a2835946c86ea7fb3e2949571c237ca198b484c7ef1bc098acfd65e499f6",
    "google/shopping/merchant_products_v1beta/services/products_service/async_client.py": "5b31f485e024f313c7f9dc43e54ccbda8103233d0b858c3e5585b5e7d62b9d02",
    "google/shopping/merchant_products_v1beta/services/products_service/client.py": "fbc29073450deebfb7d2a916d974c979b8a4d99d755a1a1e9a1408e5205686d0",
    "google/shopping/merchant_products_v1beta/services/products_service/pagers.py": "dff6a9009c4fea42a5322dbe80eb4b9033deb450e3406deb02c74d6bc543df19",
    "google/shopping/merchant_products_v1beta/services/products_service/transports/__init__.py": "d458f904826a6d2b81031a2448cfb45eafed72a332adfe2a32979193b16d7c0d",
    "google/shopping/merchant_products_v1beta/services/products_service/transports/base.py": "a289dfbfdf39a0ac6db33385b9a33b8d6a9215ef1ff19b885f6d72e4edf9c54e",
    "google/shopping/merchant_products_v1beta/services/products_service/transports/grpc.py": "3d1fb130c23d728382ca4cd8a33e5699efcea9bd4e34355aadd8a594df5dfa1d",
    "google/shopping/merchant_products_v1beta/services/products_service/transports/grpc_asyncio.py": "f8556508678360397019ba4fbec2f12ecd1214e30324ff1dc91aff295fc357ca",
    "google/shopping/merchant_products_v1beta/services/products_service/transports/rest.py": "7ac04fecca75bbb67cdfcda7c2b54eadffe984111e37d57a9859a94ceb8b181b",
    "google/shopping/merchant_products_v1beta/services/products_service/transports/rest_base.py": "8ed35ddb4588d6eb061615f1a75353584906de69a152d76343683c14870b910b",
    "google/shopping/merchant_products_v1beta/types/__init__.py": "9fc91c4aa7a25df6d6a2e4c9a4fd40524e8047698fa27d88bd96e02b0ae19929",
    "google/shopping/merchant_products_v1beta/types/productinputs.py": "cd6e63c8fa1c720b13b1b65cb980591720553af00c6e7f5dd855871da94c4da5",
    "google/shopping/merchant_products_v1beta/types/products.py": "4d162bfe61dec2ad6c60bcd9a2a2cf4c126c1b99bec1e8e5ec1afef087dc3f39",
    "google/shopping/merchant_products_v1beta/types/products_common.py": "b2c548dc5175784f3f6230168c8c91d636dc120ac091b1ea4aa418abd9fa2edc",
    "google/shopping/type/__init__.py": "2b71005e6cc68025caf866b7657689ea869dd58b337cd9cef01735b42e6044da",
    "google/shopping/type/gapic_metadata.json": "c0ca4f3ce207fd9caaa57a9559bb9d34f38b7ae194cda2cdbf64bb220f06de3d",
    "google/shopping/type/gapic_version.py": "fe61fc1cefc6f13e04f2bb60a9767bed648914c57369be63ee27b701503724da",
    "google/shopping/type/py.typed": "ec3ef2cbeaa327dcbc5fd8d11dfc96211167ec34bc0ab06cd9f6c080217b927e",
    "google/shopping/type/services/__init__.py": "02edd8492b7d3d907692bbc187e7775581a62d1e0422577e67132ea1dbedf119",
    "google/shopping/type/types/__init__.py": "379240332047e5498100ac525934b5ab42b859e178d6268bc778d6d674e4c766",
    "google/shopping/type/types/types.py": "1842413f181f14a8e8a6962de6a0024e88f9e6370680cccaf5f313810d670c5b",
    "google/type/calendar_period.proto": "f7217988c4b2a5e5b6179c2083f521bd9c6709e764b021bea8b7854516b8eec3",
    "google/type/calendar_period_pb2.py": "91e801054ae0531fb9f5d4100a4045f8ec303f8bc4c4dc4bb722896872b46b97",
    "google/type/calendar_period_pb2.pyi": "b15a3cd76c0b4a0dbc55e0bfa191cd48d6efc5681b6f0f855b258303520aff18",
    "google/type/color.proto": "56e75860f7b31538aa69971d59701ef042624d7f1a5b9c78e01e64dd5249d681",
    "google/type/color_pb2.py": "d97366645aec2a8d893868cdc0dab73005d32d15e112a2e5d84d671db930a868",
    "google/type/color_pb2.pyi": "7cc1d3c61705a4dd68426efbbb2978aa6a48991915a1555d04e472d92fe2d1e9",
    "google/type/date.proto": "965aed94dc0496471176bd607f7dc9413854f07c81253933da075dee2669a9da",
    "google/type/date_pb2.py": "df532075efb9694808da59116f7e4742b3687ec60d2522ae3ded175a9065bb76",
    "google/type/date_pb2.pyi": "6e1728feedea4ec9294b13ed871acd62869fbeea67a373fb0f43c339f8b0f0a5",
    "google/type/datetime.proto": "8d3e9db9a18d579f0b18577446c6eb89f6272ff8d297e9e1642c8354534b649f",
    "google/type/datetime_pb2.py": "90ad8b71ef848f386a527ece5f6ce0d99672a9822658e10a92f3e62f92e6084a",
    "google/type/datetime_pb2.pyi": "d90daad0afe09082e4e8a6166e3aad83c25b07f4e1fea40aec9ac793bf425074",
    "google/type/dayofweek.proto": "ca681f90502ad818207105f07e50c0f5e2676c0af0a87feefa071b299c43205b",
    "google/type/dayofweek_pb2.py": "4dd71d2a2dc9523df45cec231bbef7e5ebaaf1cbd4f8305bf357e9523327afce",
    "google/type/dayofweek_pb2.pyi": "45b17589bae2f3196848f3c7476e7336ecc14f3a706f39fabbf146b13fcfe221",
    "google/type/decimal.proto": "07b1b8b884b250ec5d7a247b0361054265ef9b4b686a13fcbd6e93679e2d7660",
    "google/type/decimal_pb2.py": "797e692c4438c24c634e137da207a80db31a8fb756171bbb13de67b1665ab89e",
    "google/type/decimal_pb2.pyi": "11330d81b0d2324c1bc15cc5d5dcb5898bf9175bdf49956b65369c062ab0e0c4",
    "google/type/expr.proto": "d1d3dc4fb64016ea18d21d52ed97ff67e6556d69bffc94ff4b02f6d80fbd0a63",
    "google/type/expr_pb2.py": "adaca5cd67653d4fae028ab3d3194e2fe1cd23b5079ab2af063886fc15135216",
    "google/type/expr_pb2.pyi": "f2947104dcbafa5eb469b38d81aaf7c2c445c6134ed7b0175d76d992fd7e6877",
    "google/type/fraction.proto": "92689823752af9e8941677757dcae82bcabc2b4da702b2e7d7a203354f0e8e9e",
    "google/type/fraction_pb2.py": "967634ce0860759ad02d781f63a6f97ce55083e95988dde517a72297a787a7ee",
    "google/type/fraction_pb2.pyi": "c71a3095698a8bd8e2aa94f124f255b945b8f8eaa9993ffb269f83d658866926",
    "google/type/interval.proto": "29aeb806c14aa5382f1bcbd95bd73235dea1ecae360d5757ba5c9e1bf16e666a",
    "google/type/interval_pb2.py": "743ccf1b6a31c9cb4fee391c6e199444aeb585b3181b0dab8e2fef697fefeeb0",
    "google/type/interval_pb2.pyi": "4e1339d796e0cc661f117592d53cb757ca99b2f11378604a12c2f4e4002cb64c",
    "google/type/latlng.proto": "39666561b6c2cb45d488a56ad909b8ec9b5a086ccf62e28c662059170f3d1cbd",
    "google/type/latlng_pb2.py": "d0d0abf15a3500da6347b84f8b4a08701512d79ae39068721be40dda146370d6",
    "google/type/latlng_pb2.pyi": "1f5fdd7771947f0f8cdc2ae8d95a325104cde385c35ca4d6698ae9fcab7f9de9",
    "google/type/localized_text.proto": "4225588ea1a34c13416636ffff2a08c849b2fab4b82cbc2de105d424abf17106",
    "google/type/localized_text_pb2.py": "b768a6f09498dd9eb00f6542902198cb627c4f437a260fed64138c6299afdc11",
    "google/type/localized_text_pb2.pyi": "94751db288f04960d830d043cd2cb8ffd38ac9f94205d03dd706744c047977e0",
    "google/type/money.proto": "4a139ef51e243525668ac1ddd0591225452d89633efb0ed42637e87b7c18f349",
    "google/type/money_pb2.py": "c94c7ec6d341d8133439ce1c81408afd1a75caa1db7f6c8b3dc64a3c1a491333",
    "google/type/money_pb2.pyi": "7535257aeac26011a5dd4b058042fe00819eaeb3ddebe1d4017d3690b9976786",
    "google/type/month.proto": "7db9afc53a778a06409c4a86c45d1e73ff18931b64344e8e706363058bc49974",
    "google/type/month_pb2.py": "f83140b3e5336d9e0354bc00ae19ee250e9e38183aee9397b6eb971a24aaa59f",
    "google/type/month_pb2.pyi": "08ffc4040deb0a485f1d1a8e40340c2a5f2cbf782f423a56ebd910ed3e4c76cb",
    "google/type/phone_number.proto": "c159b838206a754d75ccb9f687ce021351b56dbdfeef6da8fec5c7c8a20d486d",
    "google/type/phone_number_pb2.py": "dab4d2fc4f7cf3611d087adbad9c45656b5eff30b170bfc019a239be2ccc6ced",
    "google/type/phone_number_pb2.pyi": "f2b41fbc7b24eae61e8857c732ebe39195ebd8a21944d71c201b8ca3333a1ac5",
    "google/type/postal_address.proto": "6420a116343be5fdf1afc571ed9d2a7f3972723d66e0d26338a0cab03a8003fa",
    "google/type/postal_address_pb2.py": "6156815efa530d21311f90bca54cbbdff1902ee5931b9ab28059b21cc3932330",
    "google/type/postal_address_pb2.pyi": "90cf4790c2949b92a2a63f3d851d43b727be91ee9a3d6352bd82b9d08975afe8",
    "google/type/quaternion.proto": "e33bcce0b925e547024448a2a7d832167f6b28c6ac72192791db2363b438a375",
    "google/type/quaternion_pb2.py": "087f9245b72ce516928d817aeea2eb2963c70c884d53f450974c4b282aed1dfc",
    "google/type/quaternion_pb2.pyi": "efc923ff3d065ebba56f2a0271ba27dadce55ac607c30a133d59d639697568e9",
    "google/type/timeofday.proto": "95e0e5b05489f3c7e17b115ff8db9fe5b6b9f5162e83381e34d8208cf659728e",
    "google/type/timeofday_pb2.py": "7c1d47469ea578f3212adef01cc0ca47d1a5f95abf6157a0d45bf377a37cc931",
    "google/type/timeofday_pb2.pyi": "07d9b33024bf2c482a506aef6ffd1f34363402729c76e02e3799897bbc639f5b",
    "google_api_core-2.34.0.dist-info/METADATA": "3c256848c45d6f20c5e03bbf034a57ee4caeedcb327c1b74f01f960eccbc1660",
    "google_api_core-2.34.0.dist-info/WHEEL": "2b6eb4118ce7cd7b09601406aa623c553c4476265836f0d9c16f5c061f7efcc0",
    "google_api_core-2.34.0.dist-info/licenses/LICENSE": "cfc7749b96f63bd31c3c42b5c471bf756814053e847c10f3eb003417bc523d30",
    "google_api_core-2.34.0.dist-info/top_level.txt": "ff542f48922114019fc5befd0fa0e107b494c365fa4f8af09f3fcb2eb6dc0f77",
    "google_auth-2.57.0.dist-info/METADATA": "43c435aa5c37f392f12a90ab0dd67541643b7d4bea94633c332e7e52ab7146bc",
    "google_auth-2.57.0.dist-info/WHEEL": "61532836a2b3111b7ec23519c09df7c41180c27165fb872a6d8913b566b88ad1",
    "google_auth-2.57.0.dist-info/licenses/LICENSE": "c71d239df91726fc519c6eb72d318ec65820627232b2f796219e87dcf35d0ab4",
    "google_auth-2.57.0.dist-info/top_level.txt": "ff542f48922114019fc5befd0fa0e107b494c365fa4f8af09f3fcb2eb6dc0f77",
    "google_shopping_merchant_products-1.8.0.dist-info/METADATA": "53270f8a249ed2858c9cf5a82abb7d1c5405486f14e47ed4c284ff54463c9d69",
    "google_shopping_merchant_products-1.8.0.dist-info/WHEEL": "2b6eb4118ce7cd7b09601406aa623c553c4476265836f0d9c16f5c061f7efcc0",
    "google_shopping_merchant_products-1.8.0.dist-info/licenses/LICENSE": "cfc7749b96f63bd31c3c42b5c471bf756814053e847c10f3eb003417bc523d30",
    "google_shopping_merchant_products-1.8.0.dist-info/top_level.txt": "ff542f48922114019fc5befd0fa0e107b494c365fa4f8af09f3fcb2eb6dc0f77",
    "google_shopping_type-1.5.0.dist-info/METADATA": "ca7f2d2da14ed1acf23311a89516a08472c46531587934e161049d3323bf06bb",
    "google_shopping_type-1.5.0.dist-info/WHEEL": "69e6228a0d35958183cc1812f07c565ce837b95eb51bd8a33acba9fa4f68d6c9",
    "google_shopping_type-1.5.0.dist-info/licenses/LICENSE": "cfc7749b96f63bd31c3c42b5c471bf756814053e847c10f3eb003417bc523d30",
    "google_shopping_type-1.5.0.dist-info/top_level.txt": "ff542f48922114019fc5befd0fa0e107b494c365fa4f8af09f3fcb2eb6dc0f77",
    "googleapis_common_protos-1.75.2.dist-info/METADATA": "6cca679d083b78f033e50dd9249696af34007bae9f7a6ab6d05ed7f066271655",
    "googleapis_common_protos-1.75.2.dist-info/WHEEL": "61532836a2b3111b7ec23519c09df7c41180c27165fb872a6d8913b566b88ad1",
    "googleapis_common_protos-1.75.2.dist-info/licenses/LICENSE": "cfc7749b96f63bd31c3c42b5c471bf756814053e847c10f3eb003417bc523d30",
    "googleapis_common_protos-1.75.2.dist-info/top_level.txt": "ff542f48922114019fc5befd0fa0e107b494c365fa4f8af09f3fcb2eb6dc0f77",
    "grpc/__init__.py": "c3a37159b3817d71ce8d6b7a986555e4724e35787d518f3656410be19dd89730",
    "grpc/_auth.py": "b866c512875fc634d4f134ad70e55895031ad1f07061ed267fba49537728e546",
    "grpc/_channel.py": "4723bd60f3a41153db1c43efbdd0a4b7e0ba7fc6ff9de75160d61e3054c5ec00",
    "grpc/_common.py": "56001444100eaaab8a03c26a9ac8744346d630e1e740cbf983096e762a61f1a1",
    "grpc/_compression.py": "928ca76fd19f196711b5accaccb5c7c3f84e7fb9576582ef8dfaf617eedaaea3",
    "grpc/_cython/__init__.py": "bfe6cc9a17e79183fb911e8cc3dc0f2c6fdc09a824f992a25198baa726a9d865",
    "grpc/_cython/_credentials/roots.pem": "8e73da9610986216bde16ec8953c1bcc5bc2ca9081f3ae322518809ef7126c64",
    "grpc/_cython/_cygrpc/__init__.py": "bfe6cc9a17e79183fb911e8cc3dc0f2c6fdc09a824f992a25198baa726a9d865",
    "grpc/_cython/_cygrpc/private_key_signing/private_key_signer_py_wrapper.cc": "87e42340a2902d397058e3e3eb0a8c5c704b3165bc1d5a265ab3a07562495baa",
    "grpc/_cython/_cygrpc/private_key_signing/private_key_signer_py_wrapper.h": "149ee27ddb90ed32d6b5be28c2557bf0cfd6a294664c21ea056ee324727ecb90",
    "grpc/_cython/cygrpc.cp314-win_amd64.pyd": "e6e08f7e985b92f086a552b92c6a8ae164ac467dd85c2b0d50fd87b115c20620",
    "grpc/_cython/cygrpc.pyi": "45954c70d1fb5132e062a9a3df0752db826038a79ca6096867556951a0dfba9b",
    "grpc/_grpcio_metadata.py": "a1618f150eaad7afa7b7d0b7f5bf2ae68cc00221c0a36af265f12081b5ea1ad9",
    "grpc/_interceptor.py": "37d90a853a9255d1a0b036d9ac9a1829a43b7e0ad4e06e6682dcea309f5f667d",
    "grpc/_observability.py": "161845b6c0230b7098a2eb2c3da7c05cc509e7ba1b578c28f6c5cf231337c7df",
    "grpc/_plugin_wrapping.py": "b643cbe065bf695c6f09c0f925e58928fe30d140142f1e3e4ee821144cfa0db4",
    "grpc/_runtime_protos.py": "0d4d4813ef481dcc6100fdc6346fcf8b1ddf245d2b60f912c08ca7921b21a696",
    "grpc/_server.py": "677eb160954e1b764220e02adf4c3f55e9cd95954ca1812dbfe84e037fc49761",
    "grpc/_simple_stubs.py": "280cc83f3f0542bf84a77556df3454ca47cfcb5b764346c3e3470995acf5c6eb",
    "grpc/_typing.py": "66123f381ee85c8342f32f2624b3ee1298263a6101e6a490327812d3904823b2",
    "grpc/_utilities.py": "4a4677f1d1718c98d7ae3b11adb6f1960ec1915be8084912368f79a8b2d1b4f4",
    "grpc/aio/__init__.py": "2c307f85b770cf1149dc748a9a5be773fac8d031a27b477716e0e3412dfe0f15",
    "grpc/aio/_base_call.py": "688743719377f102807cc9090f9a93b2efd6b2c6d242d2d9723f26060a8eb0a0",
    "grpc/aio/_base_channel.py": "a740347cb36c59186cfea7f047095404de7af439f357ce0a5d81757eace40cc0",
    "grpc/aio/_base_server.py": "9a1bf146a7baded42500ad09539b7fd8c5279123ac1c2e331f99945d7695ef24",
    "grpc/aio/_call.py": "889558a666fdf7dad46ba9698f3f61e1ab72c6a25f3a5957107136f47a9c29cc",
    "grpc/aio/_channel.py": "901d9c0a6ddaa1a48bb97d9de53368740dc7131801fb7942832615f42f03242a",
    "grpc/aio/_interceptor.py": "c434fc4e7f34c2ea89f9c0b4bc0c1d0a1cf0468ec6f32342849ed6dd79216414",
    "grpc/aio/_metadata.py": "9d4e9a8f31ac637d6a6cca9e24bd8f3186e97a46e597d98222fad927676a6df5",
    "grpc/aio/_server.py": "bc9f9d49c2ee38cf3d2a6e6a38be975702bb764dc3a34fc7b452591f034fcbc9",
    "grpc/aio/_typing.py": "f08d97104285f8f8cc12ebab0cf3bdb4d11042360f42e276878e1b5d5a4ba9ef",
    "grpc/aio/_utils.py": "170f05a4ab6e704d9ad9544d17727200cce1d01b6035064aa44e635a51e5c10c",
    "grpc/beta/__init__.py": "bfe6cc9a17e79183fb911e8cc3dc0f2c6fdc09a824f992a25198baa726a9d865",
    "grpc/beta/_client_adaptations.py": "47c4570b225c488af25280303621ea88796a6cc1b59dd635bd81b656b1571377",
    "grpc/beta/_metadata.py": "d10a0f04be0bcc8f40c28c3630f42779c793f77705f7c89247ff2bcaf26d340f",
    "grpc/beta/_server_adaptations.py": "6d58c4e1c15081d519c797ca5af2cfad2001d300f54eb98e5cc4a62b80fa9cba",
    "grpc/beta/implementations.py": "4e12525b6f2dceb4fa891505944d2736b5a527bd8d75179828b1f77e4d9cd85a",
    "grpc/beta/interfaces.py": "27ec9cf28226fae6a9a2165f0486252abf8e4abbcb9e3e1f3c6e6435258506fa",
    "grpc/beta/utilities.py": "8b035724dd06a7e8cf77b5aacab70ce3884cb9c39c04f304c63cc03c2b9390f8",
    "grpc/experimental/__init__.py": "7c7c1913598d2ebaf81750637f7db407e496847cef5259a76196ae61ab146de1",
    "grpc/experimental/aio/__init__.py": "42812d69ae42e34f0879a69230abe8a183497560aac97e97f54605fa66897086",
    "grpc/experimental/gevent.py": "66614bd222bb8ab842f49b530b62493f6dfe2111b7fd90a8851401861556bb23",
    "grpc/experimental/session_cache.py": "ec606c31a1669962dd7b4151f783e989369f3d9881cc28979047fd8cbeb44ddc",
    "grpc/framework/__init__.py": "bfe6cc9a17e79183fb911e8cc3dc0f2c6fdc09a824f992a25198baa726a9d865",
    "grpc/framework/common/__init__.py": "bfe6cc9a17e79183fb911e8cc3dc0f2c6fdc09a824f992a25198baa726a9d865",
    "grpc/framework/common/cardinality.py": "961c92d0285cda45021fce0ad3af45d8f80d2c5247a637c18d274defda252fc1",
    "grpc/framework/common/style.py": "785117e66036ca70731c57fc1f678a780169da69be7d1df8dcc886abd5eafec2",
    "grpc/framework/foundation/__init__.py": "bfe6cc9a17e79183fb911e8cc3dc0f2c6fdc09a824f992a25198baa726a9d865",
    "grpc/framework/foundation/abandonment.py": "948dc54ae8b3750f39cfcb3110a96e86f623f48eb861e226bea583b96bba257b",
    "grpc/framework/foundation/callable_util.py": "5469e66ec33fdf52a4e3a7e3ca8b412aa4b2ac23522e97bb8f5339d4c1d1c852",
    "grpc/framework/foundation/future.py": "58692a5affb85b46a1804d2246c82c29d22e54d2019f2db63ab01141be0031f7",
    "grpc/framework/foundation/logging_pool.py": "c6943d905ea144521e373de11fa7dbbc8ef21b5f8481faae6e28d5a50a693112",
    "grpc/framework/foundation/stream.py": "97654c8052b9aa13cc053d328ccef9669a473087f92f3cfa2309e9e9677dc00c",
    "grpc/framework/foundation/stream_util.py": "1769d1955b101092a20ff881153d85b6691bca740a10f8cdd0286c50d63b16f5",
    "grpc/framework/interfaces/__init__.py": "bfe6cc9a17e79183fb911e8cc3dc0f2c6fdc09a824f992a25198baa726a9d865",
    "grpc/framework/interfaces/base/__init__.py": "bfe6cc9a17e79183fb911e8cc3dc0f2c6fdc09a824f992a25198baa726a9d865",
    "grpc/framework/interfaces/base/base.py": "df9a358487d7a3b89e3e48a9abbdb691ec0311725f9a4f613f4f26cc80b14e6e",
    "grpc/framework/interfaces/base/utilities.py": "7cdfbf0d418665d6dbf07c0798e99660771244765b5a23f2d5271139bf730b37",
    "grpc/framework/interfaces/face/__init__.py": "bfe6cc9a17e79183fb911e8cc3dc0f2c6fdc09a824f992a25198baa726a9d865",
    "grpc/framework/interfaces/face/face.py": "325de5a0344bc2bba73240acfc1f6213ab8911b57741209c896c4b5b2906698a",
    "grpc/framework/interfaces/face/utilities.py": "8b685c532daad5b19f676606833f2f4751e18f1777566a964003791ba4de97ed",
    "grpc_status/__init__.py": "6389c23a1f650e352d38f2e83bb5561e407ef7a951354a3493a0acec7fc34cad",
    "grpc_status/_async.py": "5556277edce0ab487925688f95a9e5717148d6c251cac7e480a3666ec6d2b0e7",
    "grpc_status/_common.py": "7e7b65fb74997fb3342618e389c1ce6cd02370abf8a92f108d2259aacf2418d6",
    "grpc_status/rpc_status.py": "9e2d16c2ac8e4ff1a3efa654402cbaf392b24a02d51a106facd691c1ff20fd6e",
    "grpcio-1.83.0.dist-info/METADATA": "dc5f36717badd3622862331bc5d0b6ac8704579b60f0408ea0b147b6f87e48fd",
    "grpcio-1.83.0.dist-info/WHEEL": "9d93da95b41de81179804e86a7f02d126aa61a37349beedfc1022b424b259946",
    "grpcio-1.83.0.dist-info/licenses/LICENSE": "ed9abc49f1fa3e3ce962d8ac42e176125c0ab0051ec5783dfcc807dd2a4c1f0f",
    "grpcio-1.83.0.dist-info/top_level.txt": "78477626afda550169dfc6d65bc3dfc23cffe6289ba9e3854fedb35e53c0abff",
    "grpcio_status-1.83.0.dist-info/METADATA": "55158f0bb6c7cbf69cb94d4f906e9cffb44adb089d6f27171500ad06b1f3c805",
    "grpcio_status-1.83.0.dist-info/WHEEL": "2b6eb4118ce7cd7b09601406aa623c553c4476265836f0d9c16f5c061f7efcc0",
    "grpcio_status-1.83.0.dist-info/licenses/LICENSE": "58c86a3f2b5cca0ebdf9ed3e7f4e6dc38e3b38796d8c3dc92064c62a4e39d0dd",
    "grpcio_status-1.83.0.dist-info/top_level.txt": "c9093bcd7af761798ab60ece75d22fea39574d31453139c3a574a56b6c6d579f",
    "idna-3.19.dist-info/METADATA": "4d113161aca8582e8d28fddf3ea50f19b607209c2b3e4364379a163454680084",
    "idna-3.19.dist-info/WHEEL": "ff23073da055face9f40cc14a28f3735275308cddcf315d6157461a9a7cdfbb6",
    "idna-3.19.dist-info/entry_points.txt": "ec7de718e1daa778e72c4e5eeeaed8c2bf55abc6b107b5888f9f82ff8376be1c",
    "idna-3.19.dist-info/licenses/LICENSE.md": "1a9a4f0e3d479a27240ddd59a9137a66ab4a0f9dfdc8ca6188cc0bfd85187f04",
    "idna/__init__.py": "8514c3ed53136a3596ebdf512fa487bbdd7da5a99adcaed82e0363d2c306d3af",
    "idna/__main__.py": "e0930aeba5a3e2e2d94ca6c5fac4f72c0c4eb2be9bba283becf98e909091471c",
    "idna/cli.py": "2e01972a25e67afdfba5d824013c01e808893916539c0db33f72f37c42b8896a",
    "idna/codec.py": "71ae48bb9eb70f678b195ea6e297d06f1900378cecc50c2bd1b705a12cf35267",
    "idna/compat.py": "4ff7995ce52dd1be5381ae7a70fc9b42ef420982b41532b599d88f6fdea76bc0",
    "idna/core.py": "d9f2cb3e2f9da9c494ffb1e27e501441d395cf9c12aa8321c9656fda9bb0da0d",
    "idna/idnadata.py": "01ffa6a3c5819a48402ba4f25ca3901fcf0e5f49903432ad74bed45ad429a699",
    "idna/intranges.py": "838f6c70b4a4a89b4084b98e3836bb8551eb8928e66fad2d893b04a287097412",
    "idna/package_data.py": "335f9830c0f69a399ae75464766817ee8bed94d3e9944b286e2629b13213cd1b",
    "idna/py.typed": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
    "idna/uts46data.py": "a8d8adbfa87a4f72a8a9f9e386d8cc83eefd4c88f8d364c4826eaba258481567",
    "proto/__init__.py": "52daea47d1870b13d17ff5bac292d8ba230ecbbe17c98c85e8bb9a7c03967b90",
    "proto/_file_info.py": "e1a47b159ca7674a1851ee7cead6b4e0ceff790f177ee97ef93f80094b40177d",
    "proto/_package_info.py": "51d17e579135f7792f3adfc8a94a9eca527377163507f45161cefe68230c97ce",
    "proto/datetime_helpers.py": "1c9688b86eb739432a3d0b636c2daab6b37018ade6f59fb5928756118d07d587",
    "proto/enums.py": "73de586daf6799c01d8d4d9efe76e6bf8673a4688825c941608708ac6bb7d1e8",
    "proto/fields.py": "3fdbf3036049cc8733aeb33266fdc5c4593d3fc5a75c018012d00374133c70ad",
    "proto/marshal/__init__.py": "ba085dc41813a59fc6b4e5a74b659acf5c9511b0a7a3eeb59b97915266316d2a",
    "proto/marshal/collections/__init__.py": "997d9b3be1d3be9e2f716c350ee7bd58c7fbd7bf903fded984d787be45c626cc",
    "proto/marshal/collections/maps.py": "f9736e49cf540836ce3f7bba786ac9632e8fc3df9d81a760d7f56abc918768ee",
    "proto/marshal/collections/repeated.py": "8f4160bbffe021736403fa6b74a38f124c752f6b71cb9850f8695da2291b843d",
    "proto/marshal/compat.py": "30b03f0b7cd73bc414240dd578c5c18324a66aaae81aa8cdf136b9db44a4665c",
    "proto/marshal/marshal.py": "cf92977d4b01537c3081fbd0c5ac6552111b70e78afea8698ecc8c3a17e26b6b",
    "proto/marshal/rules/__init__.py": "addb6d8ff8b47b8229459617f686524ece409e3016e1997ffc1a0d1df5141e95",
    "proto/marshal/rules/bytes.py": "01fdba986e5ed47aa7d2dcb93710b186ad26ec2c5a7985936ad390f214cdd233",
    "proto/marshal/rules/dates.py": "f281a996942d8b16fc8e71ca988ed28e0328414912f0a361a727f7c2af52892b",
    "proto/marshal/rules/enums.py": "f4d648436d78b2063745cc91eb8d134569ca2c7e2b977ea64561b2539e8403ee",
    "proto/marshal/rules/field_mask.py": "fc21bd92f95ae3edf6a6a1d5633a172433df0bdcb2489cc424200a6500bed88a",
    "proto/marshal/rules/message.py": "6ef001127ac17bff866a3ed3be2aa911fca02e80bac519ef05cbf8c85330f245",
    "proto/marshal/rules/stringy_numbers.py": "b0eb02a25f9d07b7a96d50a3fc066ed59118a39b9a966839ce6e219345f59255",
    "proto/marshal/rules/struct.py": "918da21993b473c231a23805f0e9657ee0dbeb79667698c7fb9b41a1de732e84",
    "proto/marshal/rules/wrappers.py": "14356ec0bc570a7866efb8f4919f85d615466383e596681085afd31fdbe31377",
    "proto/message.py": "65f6ab7190fa3d9446b2306a8c357381be02ba46f58636f27b32c4137e547e69",
    "proto/modules.py": "d60dc70195b4256071ecb0ab94c3da99438c1f169a7fe405c1b61aeb3822a658",
    "proto/primitives.py": "170bba290767d62eaf54f3c3dccc505c72af6f3ef3e41e364253b049391c89d2",
    "proto/utils.py": "5f8d0d16dfb17b8fcbdba0aa2ac6813f88583d4abddb0eb636262aaeb9d39ae8",
    "proto/version.py": "db2cca5ceb9848022629444d5ea7e6d55f5a5dcfa5c8981fb02e41352bf00491",
    "proto_plus-1.28.4.dist-info/METADATA": "8c802fadb1eee700e21d04a90ecad523f2c1c92259299e134a81dcf65965a9c4",
    "proto_plus-1.28.4.dist-info/WHEEL": "61532836a2b3111b7ec23519c09df7c41180c27165fb872a6d8913b566b88ad1",
    "proto_plus-1.28.4.dist-info/licenses/LICENSE": "cfc7749b96f63bd31c3c42b5c471bf756814053e847c10f3eb003417bc523d30",
    "proto_plus-1.28.4.dist-info/top_level.txt": "1c583f356f55c610f24b3a860e6531a9487ac3c4199552e1e190c8c5892c4c23",
    "protobuf-7.36.0.dist-info/LICENSE": "6e5e117324afd944dcf67f36cf329843bc1a92229a8cd9bb573d7a83130fea7d",
    "protobuf-7.36.0.dist-info/METADATA": "a318aecdc5dbf73e6fa1377b1c3b7b8e3e439e17ebeb4cb37357c2eb6a2ef182",
    "protobuf-7.36.0.dist-info/WHEEL": "30d18a8aa5737049bf47ec7833f079f591ab68a9aef577a21773d55a595d7ece",
    "pyasn1-0.6.4.dist-info/METADATA": "a092d99fad1e3744d6aa26eef2598e8ee2343ea403137ca73d06dccef19d70e2",
    "pyasn1-0.6.4.dist-info/WHEEL": "2b6eb4118ce7cd7b09601406aa623c553c4476265836f0d9c16f5c061f7efcc0",
    "pyasn1-0.6.4.dist-info/licenses/LICENSE.rst": "2aad5fc00f705c4a1addb83eed10a6a75d286a3779f0cf8519d87e62bc4735fd",
    "pyasn1-0.6.4.dist-info/top_level.txt": "76734442dde720320ee6648208e079a1b407ae30ce52c47271d06e8dcdafad61",
    "pyasn1-0.6.4.dist-info/zip-safe": "01ba4719c80b6fe911b091a7c05124b64eeece964e09c058ef8f9805daca546b",
    "pyasn1/__init__.py": "a31c29fa7f4793674abced83c239362fe89c6e9a748068cf202d9d3b28a1b5d3",
    "pyasn1/codec/__init__.py": "1040e52584b5ef6107dfd19489d37ff056e435c598f4e555f1edf4015e7ca67d",
    "pyasn1/codec/ber/__init__.py": "1040e52584b5ef6107dfd19489d37ff056e435c598f4e555f1edf4015e7ca67d",
    "pyasn1/codec/ber/decoder.py": "3016055c7113a5d33e6308d5f5a698629042679b8126fb7b3367af3c199ec8cb",
    "pyasn1/codec/ber/encoder.py": "964b55bfd1489151e1f4f45dcc8f5b44c5f08900919e5690fb9db2fd953815ab",
    "pyasn1/codec/ber/eoo.py": "76ca4b29cdb1aff5b94db72bd9671f2ddfdb24b84e8e8b6ad58c4a9f70c240d2",
    "pyasn1/codec/cer/__init__.py": "1040e52584b5ef6107dfd19489d37ff056e435c598f4e555f1edf4015e7ca67d",
    "pyasn1/codec/cer/decoder.py": "fa65fd94f7afb7912b3c7485fce59d24393d81cb4a53214f9ed68505e7aced98",
    "pyasn1/codec/cer/encoder.py": "47feb32d6543a9041d491e83ae55d07296fb5556b4b7c8365ea603c93f51bd71",
    "pyasn1/codec/der/__init__.py": "1040e52584b5ef6107dfd19489d37ff056e435c598f4e555f1edf4015e7ca67d",
    "pyasn1/codec/der/decoder.py": "86a03127ab48e1ff3407aeceab998ca14a6dc831ff8abc95419f0c10657f15f5",
    "pyasn1/codec/der/encoder.py": "32cda1aa485ba4f3e3f822ba5dd0f4669657d59252b00fb418f2971e8808df38",
    "pyasn1/codec/native/__init__.py": "1040e52584b5ef6107dfd19489d37ff056e435c598f4e555f1edf4015e7ca67d",
    "pyasn1/codec/native/decoder.py": "6505c75eedf46ac6ae7850bded2dfd86fc109c99e98ab4effbf8a77df0f1d23e",
    "pyasn1/codec/native/encoder.py": "603d94ce455550d0b08b2e1e334e9d7a17b7f37c66f6132a4c5a58ffcdefa1c8",
    "pyasn1/codec/streaming.py": "569f950e1d12940e61ed3d77deb9def54365265aafda88695335654821a3ab6e",
    "pyasn1/compat/__init__.py": "fbd14e255d524c505ab5fda955188e627d781a608a0bc458dd3602c4ea9f4576",
    "pyasn1/compat/integer.py": "94c5ea6c9053ca3837e11871e89945717ca84310da7971b185a20869bf3a857f",
    "pyasn1/debug.py": "bbe5a62057dec2aa74d38d5ecefb538ef859714f4ad78388ea9d3402b5d9eb78",
    "pyasn1/error.py": "7b7e76a2a5b7dec79e87631b205dbbb054a0a627a08ecb5a6c2305c76a624743",
    "pyasn1/type/__init__.py": "1040e52584b5ef6107dfd19489d37ff056e435c598f4e555f1edf4015e7ca67d",
    "pyasn1/type/base.py": "b63051bd72104a21c44b9f9ee6b05bb279f90ad22f0600ae7e5ba30db76bb643",
    "pyasn1/type/char.py": "46f8f9ca940b3cd5dc74791f515f27ba5d575fae91fc0927d20d875322e3d6a6",
    "pyasn1/type/constraint.py": "8e6aede5eb0b6b4f795dd7d2d1b7aa6a846e5239ee1e24ca7644dd09c2b1d452",
    "pyasn1/type/error.py": "da4c186246ddda35c8544139e9384b46604438665f69fc288043a8fbd455fc66",
    "pyasn1/type/namedtype.py": "8e74c29485284598b4db919363d1a5325308fa3e5da8472ffe297367b8b48544",
    "pyasn1/type/namedval.py": "f38bbac0a39fb5eed4e3b696ac5a88651337b4edabca2be9b01a956e53decee7",
    "pyasn1/type/opentype.py": "8e3a926d3800682c6548749feba61c2dbaf1b5f87ff7c9c0c76bfcc335b7e4c5",
    "pyasn1/type/tag.py": "4bc9adfaf8b2ba6ce082707d3b5f3cd1187c5d61945c1293ff28006b35130641",
    "pyasn1/type/tagmap.py": "6a527d65f0c64c0b0f7b28074fac8e3536a05240a39608a3f36617a4f690ffef",
    "pyasn1/type/univ.py": "9b025b1fbe2f713f274887da492ff8684244cddd8ca95ea4a87628152a7298f3",
    "pyasn1/type/useful.py": "c1aaab6322af1f7d6b75162f3e748bc40377474de716ac5c0f964e99b4578db2",
    "pyasn1_modules-0.4.2.dist-info/METADATA": "162666d75297dfcdd0116cbeab134d78fab808b8a66bfa5c1d74832d9b754771",
    "pyasn1_modules-0.4.2.dist-info/WHEEL": "0a6c85234931e5c7443132e238d4116c64308c8a11d5a21807b782010e0a3c9d",
    "pyasn1_modules-0.4.2.dist-info/licenses/LICENSE.txt": "2aad5fc00f705c4a1addb83eed10a6a75d286a3779f0cf8519d87e62bc4735fd",
    "pyasn1_modules-0.4.2.dist-info/top_level.txt": "7bf0288df1350cd63833c3fd2c04bbaa1f05c7778e9a8be86ea92becd123960e",
    "pyasn1_modules-0.4.2.dist-info/zip-safe": "01ba4719c80b6fe911b091a7c05124b64eeece964e09c058ef8f9805daca546b",
    "pyasn1_modules/__init__.py": "24a02816448325bb78ae0feb0c36938a7f8957028d288c0c8373923e194d0223",
    "pyasn1_modules/pem.py": "27dc7cac4c7f8e175704a78ce17dc94b9b9c44db1d43268d26bb380dfa20d0c3",
    "pyasn1_modules/rfc1155.py": "c86dc0d029700b776209d331132c2f2ef75cf1e06d2ff7e1f7221d2266228f61",
    "pyasn1_modules/rfc1157.py": "24191b243ecb732187e2fcf16085e0030e7c5f5e56c7a46b2f9f4d3012ea6115",
    "pyasn1_modules/rfc1901.py": "aa9bfb3fc702d24732df7ad9e7f12dd005e0265381e91b2fea0a4ca569d90469",
    "pyasn1_modules/rfc1902.py": "55f07acd3873988cc13dc024a2c3d57cd407d03f8cb4a61aacbec1be03094b45",
    "pyasn1_modules/rfc1905.py": "abefbeaf9da3d80c90e072f78516d605c34634b4da843b00180262098cec5b44",
    "pyasn1_modules/rfc2251.py": "8c71fd0341530f0369e2524b0abdb8fed152a1f8bfe6bdababc04b101f847dc2",
    "pyasn1_modules/rfc2314.py": "87b1be01d2fcf57f7edeac6d8dc418022915b6f2be978c66c43fca6724831f9d",
    "pyasn1_modules/rfc2315.py": "ff4b4c1da7e338ce255183a0f8273ffbf15dfc11a53d3df67aba58b309ed44b3",
    "pyasn1_modules/rfc2437.py": "d6589ed4c20c74e1419e50044c5d938692ecd2bcfa695f0025a4a5c1632c4546",
    "pyasn1_modules/rfc2459.py": "64729abd9ba8c0a6c46f34f1a63b52362de70544b782c8a04998f03283f1d84d",
    "pyasn1_modules/rfc2511.py": "22e25e924c8380e0f999837d6b6cc77f8869292ba13c1827c1b3f1e4450e1f49",
    "pyasn1_modules/rfc2560.py": "ced7562d13e0f7d6e21863939a47d3209199ac21305cd1a050f4bef5b3b36814",
    "pyasn1_modules/rfc2631.py": "1deb789c73d5163ea8125a4400d624290ba29dc51ad26b3948eb7de0f87c8e1b",
    "pyasn1_modules/rfc2634.py": "eec4eedd8cac6c72269272e4ec26dd4088c98ede9c0bce3df97aa4b840e014f9",
    "pyasn1_modules/rfc2876.py": "c96c7ed12fcb9bea86684989e1edc3672ec62a9bf11537f4ad055f9f3fe47bef",
    "pyasn1_modules/rfc2985.py": "f062fc8e4586a4dd6deec55a12dca15607c233cd1786560e5048bd8e17005f41",
    "pyasn1_modules/rfc2986.py": "7cbc4ca50101182d9eebc012ad058252f473a6775fcb5b8e892bdfd696b457e3",
    "pyasn1_modules/rfc3058.py": "74a023335494224f87a23e69b39303cdb0eff1f1d099ad89e81f12a19aa5adbf",
    "pyasn1_modules/rfc3114.py": "d3678308ad9b9543726d36865fce6fc467c24e7cc75d76bd04ff4869556870af",
    "pyasn1_modules/rfc3125.py": "6ddf341be16bceedc0feb7ea974e162c971b2a28108ec0d8ffa9050629cc1509",
    "pyasn1_modules/rfc3161.py": "f64cff4ef439fcea413ee1d00c08765b2a8a78e4e18313eaf04e62041fac369f",
    "pyasn1_modules/rfc3274.py": "6542db30ddf092cbeffdf5af4ff0b5bec931ba1fc8cd10800fd403d61764fa5a",
    "pyasn1_modules/rfc3279.py": "b916967ef230e165c1a0937d81c021b16f0c4c3ca61a86be16bac2da4d37dd32",
    "pyasn1_modules/rfc3280.py": "2712032a0aab8449119abe78c8d9a836c1424482bc2f3f17f39c8b1e23135dc8",
    "pyasn1_modules/rfc3281.py": "20188454184aa1695bd0a270570a803d6f753cf5bebfc0e56bca79da7d868176",
    "pyasn1_modules/rfc3370.py": "ab7753a093a862f6b114ef3d66a6bf9d88bd977ad0bc6c027cb17ba857e16dfb",
    "pyasn1_modules/rfc3412.py": "32aad7c18df914f119a9735d7b2e0b678e3004eea7ac84ea0f9fad844fed2ddf",
    "pyasn1_modules/rfc3414.py": "366199ddd5996377bf3afbf43f21c374af213aa75b69aaf602e5128bb60a9037",
    "pyasn1_modules/rfc3447.py": "201c47706d1783c917a26301ea116ff8231afd8f574106be6f56ececbd1b86c3",
    "pyasn1_modules/rfc3537.py": "2342d54a7b32c8c393349c8b75afbc4fdf094a162d0e4d77c8f77f6e6ad9dce4",
    "pyasn1_modules/rfc3560.py": "dd477bb18ece015ff8286367fe1839c596e5124c44fc82c7d643f64c8f8a6d9c",
    "pyasn1_modules/rfc3565.py": "9d17a985c5d8ee2a06e48e22693ea6490606c1aa2e4505e8327a76905410384d",
    "pyasn1_modules/rfc3657.py": "1f8f386d513b549d39ad56e4e817656c86c058277b085aade5ee41a6d36cdf06",
    "pyasn1_modules/rfc3709.py": "280686ed22934fd11ff8acdae599ffa97919a69ba5f16b7c30f4a4cd2e2ab39a",
    "pyasn1_modules/rfc3739.py": "a7e0fcf7ba9162de904747d69a0a0f56607a5579f3ff598c6351d12671fd79e1",
    "pyasn1_modules/rfc3770.py": "b9ed106a2cacf09f3a33ef04b4b36b71fb979bceccaf6190e1fdf4952b34bd71",
    "pyasn1_modules/rfc3779.py": "c7c1d828609a18edc1a2109110741411ad6860602b5880b62743dfb7188fae32",
    "pyasn1_modules/rfc3820.py": "855ef452989b99c4b31f8fd39dc1cf3c7b30e0eaed16331b7a94b812f89e40b9",
    "pyasn1_modules/rfc3852.py": "55c59c6d300062d67414361cba87c3c5555231d2a08ad85d4c18b143cf7c4b3b",
    "pyasn1_modules/rfc4010.py": "276169b4942fd37a3c37e9a8eef27dab9fc0f6faf0124d1cae9bc09179951fbe",
    "pyasn1_modules/rfc4043.py": "3963e05737cadc7b39b0d4094aa5019089e182292ffb1d79fb12d2837d31c0d1",
    "pyasn1_modules/rfc4055.py": "7f6ae5c9a05e361976f3b6ffa8b2ec3698e0c316115730606ce1f6f149c96704",
    "pyasn1_modules/rfc4073.py": "6c756cb10137c97c1eb5eb354cf5804f7d1484e1229c78fcbc405a623582db88",
    "pyasn1_modules/rfc4108.py": "f88eb767472b9ff125bebf399d26bd06a025471edc2097c46fd22d3c392af096",
    "pyasn1_modules/rfc4210.py": "de7761b09e72171dd952f50ff040de26152295214075e2fc6da4f9b4a3f2222d",
    "pyasn1_modules/rfc4211.py": "4c7aabf67e0f1a0d2354c9a3f24f5b84acdbe63ec86a9c2d0629e15900272fd9",
    "pyasn1_modules/rfc4334.py": "43e7dc624b2bba7028d6dd3b1c4da39b95a54201407f9a37f6eb697a5d198dc2",
    "pyasn1_modules/rfc4357.py": "b69659fbac43b7fbbc698af8a9bc1ee542714fe25518e8a1bb13da7a60750173",
    "pyasn1_modules/rfc4387.py": "0b296f11046957f53dbd58afcd32783d50e3644040292fba4eb555a1ea519361",
    "pyasn1_modules/rfc4476.py": "9255a6305660fda32a9951ae6048d741adefde22a5a2992132e981a0b766613e",
    "pyasn1_modules/rfc4490.py": "67d9e49ed49273d9e958706357d73ce90196fe2b6f1984389798c93b46b61820",
    "pyasn1_modules/rfc4491.py": "2e97a3d7b4f931f176e459df5a8d4216a981616413252e164121a48b984282c3",
    "pyasn1_modules/rfc4683.py": "5f43fa967df866c62409ea601891fc031521a06f394ee2e87f9e632080cc35e2",
    "pyasn1_modules/rfc4985.py": "a160811b7b649c52d424e786e1a285ec9ae4038a8c8cfc891464e77fbc663ddf",
    "pyasn1_modules/rfc5035.py": "c60c3dced00cfdb24a948502ce78b6cdc13fcf712b12e5e95913c9a57235284c",
    "pyasn1_modules/rfc5083.py": "10d5c810bd0262b4eabdffe2c29bc0901049a62da93853410c5118737ef2a85f",
    "pyasn1_modules/rfc5084.py": "8bdb057549256dd4d0a1d4f26b804d167a5e171188076b92d5a3647c5759a6ee",
    "pyasn1_modules/rfc5126.py": "6a97bccbe85cb25a235343110faf890a8425271b1a4b487fd2c5acd859541aa2",
    "pyasn1_modules/rfc5208.py": "69dc22c1aeb7f5576168d3114ca5df10f97e0246c1e60aeb2c930937a15956c8",
    "pyasn1_modules/rfc5275.py": "9a08ab484be5dce267d42e3474790d4c7d384a6e1a12f4fc73af36144fdc390f",
    "pyasn1_modules/rfc5280.py": "185c12725b2fa5351cdab8c41767f000ae3cc072373e45f11c86d53247c3e2c4",
    "pyasn1_modules/rfc5480.py": "1b305380a43af15f8bf90cb4481ac240c82a47998617b53bdee5e50737d5d899",
    "pyasn1_modules/rfc5636.py": "db3d0da31236b8c9531c7180d744345bfdf6f4fb2664e8d21677222f60a2cb01",
    "pyasn1_modules/rfc5639.py": "dbc6256d4e3d8f896bfc19707e39f6f96272cca421a493cf64c356d255979529",
    "pyasn1_modules/rfc5649.py": "dc0fbe2d02fb8b0f03197483ca24947a1eb0c053ca4101b2558f7898dcd8d049",
    "pyasn1_modules/rfc5652.py": "eb903f7636064ae121b858f4f269fbcb692f56b6ebe00af1716d1b5b9255ceb1",
    "pyasn1_modules/rfc5697.py": "696aa2f90059ac4afa2396c8a6945eb7875112a145f04d6d3b606f69904638e1",
    "pyasn1_modules/rfc5751.py": "33c9132c046176a87752a9a566ffc55897ee25bfa987b3fa3151b148fed0e304",
    "pyasn1_modules/rfc5752.py": "4a38a3a9b8b90bf75e567249301b4e1f3ef005614cecf093c03599da5a26593d",
    "pyasn1_modules/rfc5753.py": "d8dc2bf06b15d938234c32ed13942222e9cf380eb5ec5dc25c1fed3d87bb0477",
    "pyasn1_modules/rfc5755.py": "459dbc35e0a71001abda92d148d170d0145bfdbfdeba59b16a0fa5453994793a",
    "pyasn1_modules/rfc5913.py": "39ac8c9a98b6f5995023512ccc8c57694f0c8708b8d41ac7e5eb28c92f3479f4",
    "pyasn1_modules/rfc5914.py": "9d739be12bc449b105608f21d2711891102b359f70e59ab1bdaff8b8074c5cd6",
    "pyasn1_modules/rfc5915.py": "56a31177f2ac9b42c5bc4e575f8fcc3ba05d146ec287b35d41cc13fc331600ae",
    "pyasn1_modules/rfc5916.py": "807ac53bd957db587a59adc99c4ab28eea97425713134968508bbdd774a2b741",
    "pyasn1_modules/rfc5917.py": "9ccd3cac69bd0f73bcbaa49b9ac86fa7bfdf1e5d9d61a4dd8541952501ded317",
    "pyasn1_modules/rfc5924.py": "ffc4ea109f50edc152776bb765aeabce534fa862e5ec8038a07b58569a0985d0",
    "pyasn1_modules/rfc5934.py": "efbcfde9278fe2233647a465e7e571ece684340f1942fceb7e10d5651cbd96a9",
    "pyasn1_modules/rfc5940.py": "ebaacc9a0c8a0616b2f9166c59a2b3ecf506c29d1ba8401550b3dbe04764d6f9",
    "pyasn1_modules/rfc5958.py": "3593f1fbb16f8f3804af3da54d446246af26dd7099ec3f506c60e1b4817ecc21",
    "pyasn1_modules/rfc5990.py": "f9bd12b7ae1b6b72d544649e3666c6a0c21b914f1cf050e9a38cc5585d0f0853",
    "pyasn1_modules/rfc6010.py": "178dc0615154c2efb6ff18c91365a6c3559db742bb977be0d3f7c26bf407a815",
    "pyasn1_modules/rfc6019.py": "bf38f9b5f1b8ebde3eb9ca44ae902d1350d5384e3ebf476437cf7866bf719b8a",
    "pyasn1_modules/rfc6031.py": "5f6723372567ad7dc6db31bb903e11abffe417af9fb669a7a879424090c2b8c5",
    "pyasn1_modules/rfc6032.py": "b8d02ee732c78346f9f37c71cc53546710a725a0b3330d62a1bcd112e7a33283",
    "pyasn1_modules/rfc6120.py": "25e846643f18d01761affa23a4c4a31e09d11c47545dab99c6a2f14709ece827",
    "pyasn1_modules/rfc6170.py": "b0bdb23d9ccefbe308e13a1e03094510ffb1e88d3e7adb89c53da68003e310ee",
    "pyasn1_modules/rfc6187.py": "8ce322221c381c0527ee18f7ee028898d53f84af136a601f77457426bc21fd85",
    "pyasn1_modules/rfc6210.py": "c0b89f2bf11286fd5ae133a1189fa4f730359155585433634be461d288660a1d",
    "pyasn1_modules/rfc6211.py": "5e8b5305056c78aef2d27241e05c7e9e9761447787e7720cf38906ba9588a6b9",
    "pyasn1_modules/rfc6402.py": "92c83a62c69c4bdbc903339ba8538f442a1ef26a72bd10dbe37674bfe40a8673",
    "pyasn1_modules/rfc6482.py": "d74fd7c9bd9368f171ef6214099b6ef35687e6b9988a185d2f43fe3d56f2d72b",
    "pyasn1_modules/rfc6486.py": "6b7ff9389be4cf61bbc56382d1da9b36a2500ec1d000e53ad805a29f5d3b7389",
    "pyasn1_modules/rfc6487.py": "813515905609c9472bd44e2ea1e3767173cd697ca361b8b1ba96c1285400e234",
    "pyasn1_modules/rfc6664.py": "9eaf05e700de3b8f45a01195403c7c8afbe0fc6b2e6dd59ad5ba5933fe344e6b",
    "pyasn1_modules/rfc6955.py": "14155bf0ba4728c663477c0e26d9be04f6e2e443224681ae516878d6bd6ca025",
    "pyasn1_modules/rfc6960.py": "06110308b2eb69ee11682a4cb8a25cd24c356c6b39e95d3f17e37188ef5cb6ec",
    "pyasn1_modules/rfc7030.py": "b7eb36043c97dd9936b32fe331097e3f623635714e9fb86ebb4c0570cfb6b2ab",
    "pyasn1_modules/rfc7191.py": "b8cb01cc9f75ebbc31b223d80d052765015514d7e05319c2c1135e2a75f13463",
    "pyasn1_modules/rfc7229.py": "192894cf842460e0df9c8bcb4572a269bca7a3d1a677a097d33591ec7a08a429",
    "pyasn1_modules/rfc7292.py": "c0e4630c60ff6aa1e8ba3076c2eea736b123613c3754eff10e9f90c745562db7",
    "pyasn1_modules/rfc7296.py": "780669676760521c5b26b2cb1ad0df7f852ca5ab86ed3af9763f1610b6079d43",
    "pyasn1_modules/rfc7508.py": "6662456d03bddf816cf30c5ca4ed208397d4d1df321059649050f728c5106c01",
    "pyasn1_modules/rfc7585.py": "4f4fac7733c9a28a758db076449fb0cd49dfeade82783d9e30c5e9733e792448",
    "pyasn1_modules/rfc7633.py": "f0ffdf0569281a4debb24ed21009ba41970f8e84464d1190b9ab1630b3ab2ca6",
    "pyasn1_modules/rfc7773.py": "e9418f5b2558ba270a7bab276429c3d70b80bb530e560ccfa1200b2f6baf4eb2",
    "pyasn1_modules/rfc7894.py": "1cb69206839407efdc484e7ddf94d7027b85055a5906fea30672ce3e9ffe2cd9",
    "pyasn1_modules/rfc7906.py": "9837f5a56c1536509c1107ecc1486d4034ad027c12fb9c5b66d8cc95f9d62dd2",
    "pyasn1_modules/rfc7914.py": "2715869d757e575df1cce9fb73bfbfdefc4336990fb5d648614e0a1769055d1e",
    "pyasn1_modules/rfc8017.py": "a703d149c86f32d5ee6ad7022d42c7baf48bf2400f12a902e1a2098dde6f100a",
    "pyasn1_modules/rfc8018.py": "f3fe3dc40def10e76519485ab30db14d4bf84e91c1be346ea289cc4ff9354d39",
    "pyasn1_modules/rfc8103.py": "a4d60015f28236b83d66644ab0d370af6a28a6d100045de031a3dba82ae84674",
    "pyasn1_modules/rfc8209.py": "f44434efbae30fdba84d958838699e39a1cb0c3ab42115e1dd1b7479807e62c7",
    "pyasn1_modules/rfc8226.py": "9ae765560aec27a5de1e7166c41356fcd81c61c16c1d4bcad38fcc4ebdd49113",
    "pyasn1_modules/rfc8358.py": "6a21da5d001a68ffaae5cf74c7fb991d2a5045307ec9e93085ee95f7e12dac58",
    "pyasn1_modules/rfc8360.py": "4f8b18ea8d952d53e767db38c89f0fcdf540f18eb49defb529c54db70e72b7eb",
    "pyasn1_modules/rfc8398.py": "8b797081fffff6827339a68724a5a60c0c7777f75e28d082baf2035aa416889e",
    "pyasn1_modules/rfc8410.py": "9ed78ac9329c230565821d6a525fbc904eb790a1be2a05ad2eb7c5f764d6cb24",
    "pyasn1_modules/rfc8418.py": "79308f4ce9bab7e4721dde8f968c0ba200f350eef695175d11260b89288ea42d",
    "pyasn1_modules/rfc8419.py": "a9cbc1957c6abec0af1bf17a00a2a3a8175eaea6d6c01cbcce364e8c03dd614e",
    "pyasn1_modules/rfc8479.py": "ac32b3ae9f8c984174b7713b96a297860c1c820bf1f0da1655bb491c62f10d83",
    "pyasn1_modules/rfc8494.py": "18c86dd517406e31cbb52a87749d9c2cef075d1cfa48b20f136e784f8a32d12e",
    "pyasn1_modules/rfc8520.py": "fe8d3496fd8c61c88eaa8d142a3959050358fcccf3810b75495f555c22344fd0",
    "pyasn1_modules/rfc8619.py": "a9262205e7cb485ba42e0e9522047a767857fae07024c22dc6a1e33579c180cd",
    "pyasn1_modules/rfc8649.py": "a070902bb838bcab3507420ef4682289d4f23ce938a73e5b6245d246604e007a",
    "pyasn1_modules/rfc8692.py": "78b2dba71ec275c0922f7320a20135ed6bcfa15dd55be2b56e8e9dab8991a1e8",
    "pyasn1_modules/rfc8696.py": "8aa1a1eee859f703b6514e93936b48053c3fbb03aeb6882492c9c3760ffe4238",
    "pyasn1_modules/rfc8702.py": "b7b27d82f299de7d365025234c778965675c0bc9a805ed90980fcfe8af0b570c",
    "pyasn1_modules/rfc8708.py": "88736c1e49351b8d7d97b28e400670ab6e722fd9bac3200298cd254b87f95c6c",
    "pyasn1_modules/rfc8769.py": "bf044909f7666abe05750f888421c9472ca7e1f566ffcd96b1532ea81783cb9a",
    "pycparser-3.0.dist-info/METADATA": "f547a60b0ab54cc2f242213d1f87960ad80df7dd5dfeb98d17a444d48bfb6b93",
    "pycparser-3.0.dist-info/WHEEL": "a842dba36b35633977f599ab0226d70368e33cb5187e756150d4e5c85d6bab46",
    "pycparser-3.0.dist-info/licenses/LICENSE": "0c846399369ea76ddd7b5c44fe6d16497415fcf015f5cbb508c24bf98b81c5b1",
    "pycparser-3.0.dist-info/top_level.txt": "73e94f712ef82fff0aa07ec813a3d0179a1fca2ad140d57856191b48520f7963",
    "pycparser/__init__.py": "a6156247202e52682a13890d99a0aaa66e56544048bf351215aa41bf85d7df1a",
    "pycparser/_ast_gen.py": "1311f9626e2993b7503c422442bf512629b97deced7414304813efa6f13ee413",
    "pycparser/_c_ast.cfg": "95de5ecc4f72cc82452150147f0edecc94a5322e275ca342cdf9aa8cec904cda",
    "pycparser/ast_transforms.py": "5f032c6ab73968375d35680888a9b8fa33963116226fdeb2350528d3fbb6f160",
    "pycparser/c_ast.py": "bb091c6561df5c3408c3a583bc22f5ee258cffed11f94443a843263e35cb3807",
    "pycparser/c_generator.py": "4552893e0bafd82be8bc71d27d092a8a672452bf70194f4fa1fc421a3db9d500",
    "pycparser/c_lexer.py": "075568a9b6213d690e2495827a6e0e638ce3d08c6b3e4b4165915ed98f3b9148",
    "pycparser/c_parser.py": "dc504a18b8e3942deff1a7f00ff727311ebb226a18232ef76b858a411ea15fb8",
    "requests-2.34.2.dist-info/METADATA": "8c384ba3e979480faae2859d3c5e6c1276dd2c3616e322e124d52c8cfc556f27",
    "requests-2.34.2.dist-info/WHEEL": "69e6228a0d35958183cc1812f07c565ce837b95eb51bd8a33acba9fa4f68d6c9",
    "requests-2.34.2.dist-info/licenses/LICENSE": "09e8a9bcec8067104652c168685ab0931e7868f9c8284b66f5ae6edae5f1130b",
    "requests-2.34.2.dist-info/licenses/NOTICE": "f5110972dedad2b4e9d314518daf3b7d72d6e02e499acd802181de6f74571dcc",
    "requests-2.34.2.dist-info/top_level.txt": "7cc4959877dbe6b6c63a8eb1bfe3bfb545fa8fe5b28b1b2c13e4a7c1c0d1c4d4",
    "requests/__init__.py": "311157e7fa4aa9166c827d5237e6fd694eeb5ca549ee40e2e90f89f21abcd56e",
    "requests/__version__.py": "abba90670b0370a3dc40d97c2dea46a35d4bec941ff10d7c3cf3098ef04a17f9",
    "requests/_internal_utils.py": "4c7d8d132c9898fc7d715e473f3ac74785ddc4ab96d2c9240f87835dc6d981ff",
    "requests/_types.py": "d06df79f54a279c0cc3faec596c446246d5ad14edfd529eb328defbd62a44a89",
    "requests/adapters.py": "647eabd4453de232aa40b55ff033d6f19778910b7e7069c51fdf63076fc77493",
    "requests/api.py": "4d15480ac046f089209798e8650476ef4a28ebe6f81b400758f8ef42ec6b5509",
    "requests/auth.py": "e93951a552d4c3c5fd9b84a5e1d49b03caaba4723c4ff112da6053600ec85a65",
    "requests/certs.py": "fd9c6b83359cef90ff6c4eeeab8dcc2388da382ebca7d00a499b3c0b434a87e4",
    "requests/compat.py": "5d79ba289690b5015044341a1f7bc7d24134fc13a2db3efcf852515bf9aa18b4",
    "requests/cookies.py": "0a1ce6140f2b9420522c5080e356e689fa796c2522f55ca10826f5ed8e0a723c",
    "requests/exceptions.py": "c5e18f3352454968e337bb3056f0867df2f37c0a0d2c1a77dc05e7e0b3c0411b",
    "requests/help.py": "723519bb1884da18d84f6b2fb78f7ebf7fb57f732070b6025ae071dac6a2d179",
    "requests/hooks.py": "ebd8a02475d31a0e473a8f553e9501ff43645b9563885ad52844e7a63f0d76ab",
    "requests/models.py": "d1bc0d990abf5d5ebee05f890911b4363fadf2d5264b686a963df47c529b6ace",
    "requests/packages.py": "fe0d2067af355320252874631fa91a9db6a8c71d9e01beaacdc5e2383c932287",
    "requests/py.typed": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
    "requests/sessions.py": "96fbb30bbbf06a59a5268d13b57885149756aa3f31695b5c15e41dd7bb2f67a6",
    "requests/status_codes.py": "1950f47c89cf18019787e07d8ce48f66d6b38e622f7d94d70a90247fee6c040e",
    "requests/structures.py": "ba9460c39078f25e6f1d2a24ac941ac6f8d2ee97197fa8c8d0c262d8a1e67a02",
    "requests/utils.py": "657fd02343b2586bbd138a4ec0f03ceeb7610c2712f96ac54f28a708bf7d079c",
    "typing_extensions-4.16.0.dist-info/METADATA": "b05084ca1d50879865178d9fff9fabeab61bdfb1f361bfbde95421ffc8f9be46",
    "typing_extensions-4.16.0.dist-info/WHEEL": "1b68144734c4b66791f27add5d425f3620775585718a03d0f9b110ba3a4d88db",
    "typing_extensions-4.16.0.dist-info/licenses/LICENSE": "3b2f81fe21d181c499c59a256c8e1968455d6689d269aa85373bfb6af41da3bf",
    "typing_extensions.py": "4040ca1a1ecbee00d1385c12a93084d1c5bd46f0b774f07e5ae7e91c4f55e696",
    "urllib3-2.7.0.dist-info/METADATA": "66147856d30bbb6bed01c6ce3326d0f48d84ed5f286ac177633a1486b82eace5",
    "urllib3-2.7.0.dist-info/WHEEL": "41c708c5adba6e097513ab8ccb9f1d7865a2fb469e22491a9e01dcc64da459fc",
    "urllib3-2.7.0.dist-info/licenses/LICENSE.txt": "130e3a64d5fdd5d096a752694634a7d9df284469de86e5732100268041e3d686",
    "urllib3/__init__.py": "24ca35b60d67215d40789daf10d0bf4f17e5d1ee61e86ce5f43195935ad645ba",
    "urllib3/_base_connection.py": "1f37121077b1803ad152beb48cd9e20762908ddc3ded279d0f996e45f1ed0978",
    "urllib3/_collections.py": "68e566da62a296fbaf4f579f006b6403ba62eb90b5342dc725fa45c6f6014bc0",
    "urllib3/_request_methods.py": "802785f3948efd45385a83f0607228cffb70f9e33f1153a42c5a7c385b02ec30",
    "urllib3/_version.py": "fa032b2a1b0f88e8f45f351ecd515ae7d7754ba4204222a0413e601a0c8cf3ec",
    "urllib3/connection.py": "668b37ab12835bdf8643a695a9f050551339b280742231f163fc575a925a1532",
    "urllib3/connectionpool.py": "b0616775d5d8c25c7b282e0908fd602af74d18b34af984c22437460021a3dd8f",
    "urllib3/contrib/__init__.py": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
    "urllib3/contrib/emscripten/__init__.py": "c325ef7bcae6a97eecd8aa91401c43e56978f23cd63e7e7ed6efd7a100442d57",
    "urllib3/contrib/emscripten/connection.py": "822125b01a14b0a55445b6736fc182ac99aa5b6dd79ef8f668dc90545e36b258",
    "urllib3/contrib/emscripten/emscripten_fetch_worker.js": "cf5937cd9e3f84329ddfeb4def0cf3f0b1e31c2da9c4dfeebbc07793d0fd0377",
    "urllib3/contrib/emscripten/fetch.py": "e7171dfbebe217165dda7072d1a2bbdddb49f53b21d7261953f4945f09f089b9",
    "urllib3/contrib/emscripten/request.py": "98bdbcb33cb52af137349856a2be633666aba7c830a650d4fbb8301996398344",
    "urllib3/contrib/emscripten/response.py": "083a58d0616896e47723110290cda01f9b837d80ccd042f5156afb07beb0b603",
    "urllib3/contrib/pyopenssl.py": "5d963e433c93eecf1de37aa2b00e28f9e49960178890fc8131986d8f69242d9b",
    "urllib3/contrib/socks.py": "ad2ba495f85ae3ebeda3ef2a48b8523e67f996644f218baa740c4a7bd6e0d910",
    "urllib3/exceptions.py": "8503c7a2aa38cb0e3a7ac3525647225505d5520031580825b21e8c6f242d6be4",
    "urllib3/fields.py": "6862c5015669554f856c996596f7b8021834fa91fd659050d88cfb95fa648e43",
    "urllib3/filepost.py": "53c78d67e9a928a1e1ae56c7104893c7180ad7a21e8e111aeeecf8db2a80fdd2",
    "urllib3/http2/__init__.py": "c73ac0487ed1e4035190f24ea2de651a70133aadca2aec97cc8e36adc9f09aab",
    "urllib3/http2/connection.py": "6c7307e9f36f6adc173eb2aaadc9fbe32037a5459ca8f0e9a672b52dc2826cff",
    "urllib3/http2/probe.py": "9e7024a9b8406a43a217be6bcfb5b4b9d677f047a1fee0fc7e357be0def71442",
    "urllib3/poolmanager.py": "734ae1d2b7140b5b79b431ae51e188800ce512b6ac3cec16c0b881ebb957f293",
    "urllib3/py.typed": "51a0ae3c56b71fc5006a46edfb91bc48f69c95d4ce1af26fd7ca4f8d42798036",
    "urllib3/response.py": "f525f806491da0b82cc759deea1a6dfb93e772b9c3ccfcdea82d4369285273e3",
    "urllib3/util/__init__.py": "faa792d1071e8af6b3bc110a0cd142008fba00271d0ce1384ccbe8ed22cd9404",
    "urllib3/util/connection.py": "2633bbdb69731e5ccb5cf4e4afd65605d86c7979cc5633126f50c92d5ad74a74",
    "urllib3/util/proxy.py": "b1e3fcf90e41e9b07474cb703e3f98719650df4bc7b8ba91bbeb48d096767f3b",
    "urllib3/util/request.py": "8ada670bcba0ec3e2755f0e6194091325824011510d77aff5ccc529f34f09a91",
    "urllib3/util/response.py": "bd013adfdba81218f5be98c4771bb994d22124249466477ba6a965508d0164e0",
    "urllib3/util/retry.py": "d989d25fefc579c312843eb5331e6cebc274fdbb541d9aeb73f064155452d4fe",
    "urllib3/util/ssl_.py": "3aa7b7ac88545377b7195819a1bda1c6647c43464386b84448f94f3ff18b6954",
    "urllib3/util/ssl_match_hostname.py": "16de38289cd3cc632629f7ff6573fdd658b65785a19af132d7a3ca05cbf8bd99",
    "urllib3/util/ssltransport.py": "133e0ef2947fbd3f1d6a7fc5bea0584ba7600df05710c7d57ebcdc754a167e2e",
    "urllib3/util/timeout.py": "e1e4f5155799654ee1ee6603d49ab639735ee1fc5e91d36f868594919bac4690",
    "urllib3/util/url.py": "59187e4cc617a2c9a0a7c9bc953e07e6ca681f0e7252395c3027d4e77024a00b",
    "urllib3/util/util.py": "8f795b64ad633f28b00f7e13f08809cdd5846554fee04fb4bd82098bd52378d0",
    "urllib3/util/wait.py": "fe987c22b511deca8faa2d0ea29420254947e30ce419e3390a2c80ed7186b662"
  },
  "python": "3.14",
  "schema": "elite-merchant-installed-source-lock/v1",
  "sdk": "1.8.0",
  "wheels": [
    {
      "bytes": 136983,
      "distribution": "certifi",
      "sha256": "62f22742b58a1a33014a2b6b706588a8d7e2a88ae7bd1a6ebe8c992928483775",
      "url": "https://files.pythonhosted.org/packages/0b/a7/71ac2cff56fec219ed242bb11b8efb69fcc4bec75db06fb7bfe35de520e6/certifi-2026.7.22-py3-none-any.whl"
    },
    {
      "bytes": 187949,
      "distribution": "cffi",
      "sha256": "3222ba5d678f80a030e6afbcc33dc1ae5cb45facabb61cee2c7016b8432fde48",
      "url": "https://files.pythonhosted.org/packages/a7/06/1c3e01e3ba14c39f6d10bfbac52753b7e22259e38088e5cfe1d704918690/cffi-2.1.1-cp314-cp314-win_amd64.whl"
    },
    {
      "bytes": 204175,
      "distribution": "charset-normalizer",
      "sha256": "c658c50ac0c98cd755a2dd50b7977d3bca7df401dcc47fbdfa87db53ef7d4e8b",
      "url": "https://files.pythonhosted.org/packages/7a/7c/4938c329b6a9d446f6a59aa2092ff7118f274209b5ed0e26893d1d30a63c/charset_normalizer-3.5.1-cp314-cp314-win_amd64.whl"
    },
    {
      "bytes": 3842826,
      "distribution": "cryptography",
      "sha256": "aed8db4f6d71c51efb89530e12d9464e7bf2923d46c3205dc794a2a93f8c0648",
      "url": "https://files.pythonhosted.org/packages/42/8b/cb12b1b60c91b074ca6bf0fdd59aa8f10d8bc5f73af8faece86ef0421b37/cryptography-50.0.1-cp311-abi3-win_amd64.whl"
    },
    {
      "bytes": 180545,
      "distribution": "google-api-core",
      "sha256": "cdf9c67e7ca2402d86ccbfde5f2503fc83e3cc3f58cc78456ae96cad24a6d2de",
      "url": "https://files.pythonhosted.org/packages/bc/c1/a8a92ae1bc4b1a8f804c776d7d3f0c771b78a62c3ad4df1be41b3fd8c767/google_api_core-2.34.0-py3-none-any.whl"
    },
    {
      "bytes": 259728,
      "distribution": "google-auth",
      "sha256": "180dafe015cfb62193bea26b677500fab5b9fd51a1e825ebf3ad9b182047ae59",
      "url": "https://files.pythonhosted.org/packages/00/f3/8508a702c094af5f6e89773f4dfdeee74913df0f41a02c21b5e7dc3d75cd/google_auth-2.57.0-py3-none-any.whl"
    },
    {
      "bytes": 244077,
      "distribution": "google-shopping-merchant-products",
      "sha256": "722ef095eca35126255c129964a61818f0d7bc9eeef7f1da12732d81c1dcd505",
      "url": "https://files.pythonhosted.org/packages/ea/7f/4de247c1900a4bb9a2e6bec51128c907545980133f1f39293c21f7c17527/google_shopping_merchant_products-1.8.0-py3-none-any.whl"
    },
    {
      "bytes": 15913,
      "distribution": "google-shopping-type",
      "sha256": "afc1180f2d068713bc1bee4b1c6882876023be6ef42e1af020991dfab2595072",
      "url": "https://files.pythonhosted.org/packages/37/b0/2b6287659916fd70564951e5d5797ca987da3cc67897fd7a1d1724075401/google_shopping_type-1.5.0-py3-none-any.whl"
    },
    {
      "bytes": 307002,
      "distribution": "googleapis-common-protos",
      "sha256": "6b83302f554ea93a0f48409c7fc2050f954bcbcddb7e3a9c76d4a823cb22920e",
      "url": "https://files.pythonhosted.org/packages/47/5b/1c9e55363c3b1890a98cae813de5b4ea327845756cd8fb7ee690140c7eac/googleapis_common_protos-1.75.2-py3-none-any.whl"
    },
    {
      "bytes": 5298932,
      "distribution": "grpcio",
      "sha256": "2bb48cb5e6dd005ca12b89ce4b6ac0b48ff3112c747542ee7986ef611a8ca6d9",
      "url": "https://files.pythonhosted.org/packages/a1/00/b1b26431c9d54eee11724fd6e5585473a2ed47fbc1fb95e5204906a642ce/grpcio-1.83.0-cp314-cp314-win_amd64.whl"
    },
    {
      "bytes": 14636,
      "distribution": "grpcio-status",
      "sha256": "f6a838a7c5fb84ae98833ec0ef81ed438c26e11e54b2ddb8e92ad328c861de69",
      "url": "https://files.pythonhosted.org/packages/d6/00/73204406228cf989bea6b0fd9fe4702fab49a8a152a0c6f90856dadb6ac7/grpcio_status-1.83.0-py3-none-any.whl"
    },
    {
      "bytes": 68550,
      "distribution": "idna",
      "sha256": "815e7be7a7806d54abb586dc943addc79e8b2ee16915059658cbeff4b1b43bf4",
      "url": "https://files.pythonhosted.org/packages/57/b0/0e52c878c53f245edd3a11020f20979b3f490f245af532c7cae3027754b5/idna-3.19-py3-none-any.whl"
    },
    {
      "bytes": 50797,
      "distribution": "proto-plus",
      "sha256": "4b01341272f8a348db3f003b6143109f83ab43091019d5181b3fcdf500ab32aa",
      "url": "https://files.pythonhosted.org/packages/41/5d/0f04b85dafdc3250ced7f2592efc17dce7f40712e941e9632202481e600d/proto_plus-1.28.4-py3-none-any.whl"
    },
    {
      "bytes": 453731,
      "distribution": "protobuf",
      "sha256": "1781cc1de61249b750848029bca452c0a8b7e990080316b9bbc2518b2117b488",
      "url": "https://files.pythonhosted.org/packages/0e/4e/12cb93270967a2affff5b3f720694700d4d87712a67afd05c8cb3f6fa52c/protobuf-7.36.0-cp310-abi3-win_amd64.whl"
    },
    {
      "bytes": 84410,
      "distribution": "pyasn1",
      "sha256": "deda9277cfd454080ec40b207fb6df82206a3a2688735233cdcd8d3d565f088b",
      "url": "https://files.pythonhosted.org/packages/9a/3b/6163796d69c3977d1e4287bea4a6979161cbbdd170ebb430511e8e1999ce/pyasn1-0.6.4-py3-none-any.whl"
    },
    {
      "bytes": 181259,
      "distribution": "pyasn1-modules",
      "sha256": "29253a9207ce32b64c3ac6600edc75368f98473906e8fd1043bd6b5b1de2c14a",
      "url": "https://files.pythonhosted.org/packages/47/8d/d529b5d697919ba8c11ad626e835d4039be708a35b0d22de83a269a6682c/pyasn1_modules-0.4.2-py3-none-any.whl"
    },
    {
      "bytes": 48172,
      "distribution": "pycparser",
      "sha256": "b727414169a36b7d524c1c3e31839a521725078d7b2ff038656844266160a992",
      "url": "https://files.pythonhosted.org/packages/0c/c3/44f3fbbfa403ea2a7c779186dc20772604442dde72947e7d01069cbe98e3/pycparser-3.0-py3-none-any.whl"
    },
    {
      "bytes": 73075,
      "distribution": "requests",
      "sha256": "2a0d60c172f83ac6ab31e4554906c0f3b3588d37b5cb939b1c061f4907e278e0",
      "url": "https://files.pythonhosted.org/packages/a0/f4/c67b0b3f1b9245e8d266f0f112c500d50e5b4e83cb6f3b71b6528104182a/requests-2.34.2-py3-none-any.whl"
    },
    {
      "bytes": 45571,
      "distribution": "typing-extensions",
      "sha256": "481caa481374e813c1b176ada14e97f1f67a4539ce9cfeb3f350d78d6370c2e8",
      "url": "https://files.pythonhosted.org/packages/49/d3/b8441a820a491ddfc024b0b0cf0393375b75ea13866d9c66727e54c2fc80/typing_extensions-4.16.0-py3-none-any.whl"
    },
    {
      "bytes": 131087,
      "distribution": "urllib3",
      "sha256": "9fb4c81ebbb1ce9531cce37674bbc6f1360472bc18ca9a553ede278ef7276897",
      "url": "https://files.pythonhosted.org/packages/7f/3e/5db95bcf282c52709639744ca2a8b149baccf648e39c8cc87553df9eae0c/urllib3-2.7.0-py3-none-any.whl"
    }
  ]
}
````

### FILE: `google_merchant_product_sync/connected_worker.py`

```yaml
block_id: "GO-CONNECTED-GOOGLE-MERCHANT:file10:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "bd2c5934ec5e32a9d06ddec91204f66775c18576fd5bbc2d5fc9b1aeff121d25"
variables: []
secrets_allowed: false
```

````python
"""AUTHORED bounded IPC around the exact admitted Google models/client.

Candidate implementation. The original build_insert_request stays unchanged.
No provider response, secret or traceback is printed outside the typed result.
"""
from __future__ import annotations
import hashlib
import importlib.util
import importlib.metadata
import importlib.machinery
import json
from pathlib import Path
from types import ModuleType
import re
import sys
from typing import Any
from urllib.parse import urlsplit

LIMIT = 32768
SDK_VERSION = "1.8.0"


def admit_installed_sources(expected: str) -> None:
    """Same verified-byte principle as the existing isolated source launcher.

    -I -S prevents site/.pth startup. Source imports compile the verified bytes;
    injected package initializers or bytecode cannot replace those bytes.
    """
    if not sys.flags.isolated or not sys.flags.no_site or sys.version_info[:2] != (3, 14) or sys.platform != "win32":
        raise ValueError("runtime lane")
    path = Path(__file__).resolve().with_name("connected-runtime.lock.json")
    raw = path.read_bytes()
    if len(raw) > 1048576 or not re.fullmatch("[0-9a-f]{64}", expected) or hashlib.sha256(raw).hexdigest() != expected:
        raise ValueError("runtime source lock")
    lock = json.loads(raw)
    if lock["schema"] != "elite-merchant-installed-source-lock/v1" or lock["sdk"] != SDK_VERSION:
        raise ValueError("runtime lock schema")
    site = (Path(sys.executable).resolve().parent.parent / "Lib" / "site-packages").resolve()
    sources = {}
    for relative, digest in lock["files"].items():
        target = (site / relative).resolve()
        if not target.is_relative_to(site) or target.is_symlink():
            raise ValueError("installed source path")
        data = target.read_bytes()
        if hashlib.sha256(data).hexdigest() != digest:
            raise ValueError("installed source mismatch")
        sources[str(target)] = data

    class VerifiedSource(importlib.machinery.SourceFileLoader):
        def get_code(self, fullname):
            data = sources[str(Path(self.path).resolve())]
            return compile(data, self.path, "exec", dont_inherit=True)

    class VerifiedFinder:
        @staticmethod
        def find_spec(fullname, path=None, target=None):
            spec = importlib.machinery.PathFinder.find_spec(fullname, path, target)
            if spec is None or spec.origin in (None, "built-in", "frozen"):
                return None
            origin = Path(spec.origin).resolve()
            if not origin.is_relative_to(site):
                return None
            if str(origin) not in sources:
                raise ImportError("unlocked installed module")
            if isinstance(spec.loader, importlib.machinery.SourceFileLoader):
                spec.loader = VerifiedSource(fullname, str(origin))
            return spec

    sys.meta_path.insert(0, VerifiedFinder())
    sys.path.append(str(site))


def strict_object(pairs: list[tuple[str, Any]]) -> dict[str, Any]:
    result = {}
    seen = set()
    for key, value in pairs:
        folded = key.casefold()
        if folded in seen:
            raise ValueError("duplicate key")
        seen.add(folded)
        result[key] = value
    return result


def decode(raw: bytes) -> dict[str, Any]:
    if not 0 < len(raw) <= LIMIT:
        raise ValueError("input bound")
    value = json.loads(raw.decode("utf-8"), object_pairs_hook=strict_object,
                       parse_constant=lambda _: (_ for _ in ()).throw(ValueError("nonfinite")))
    if not isinstance(value, dict):
        raise ValueError("object required")
    return value


def load_owner(expected: str):
    path = Path(__file__).resolve().with_name("sync_product.py")
    raw = path.read_bytes()
    if len(raw) > 65536 or not re.fullmatch("[0-9a-f]{64}", expected) or hashlib.sha256(raw).hexdigest() != expected:
        raise ValueError("source lock")
    module = ModuleType("elite_locked_merchant_owner")
    module.__file__ = str(path)
    sys.modules[module.__name__] = module
    exec(compile(raw, str(path), "exec", dont_inherit=True), module.__dict__)
    if importlib.metadata.version("google-shopping-merchant-products") != SDK_VERSION:
        raise ValueError("SDK pin")
    return module


def request_from_intent(owner, intent: dict[str, Any]):
    product = intent["product"]
    account, datasource = intent["account_id"], intent["data_source"]
    generation = intent["generation"]
    if not isinstance(account, str) or not re.fullmatch("[0-9]{3,20}", account):
        raise ValueError("account")
    if not isinstance(datasource, str) or not re.fullmatch("accounts/" + account + "/dataSources/[0-9]{3,20}", datasource):
        raise ValueError("data source")
    if not isinstance(generation, str) or not re.fullmatch("[1-9][0-9]{0,18}", generation) or int(generation) > 9223372036854775807:
        raise ValueError("generation")
    if not isinstance(product, dict) or set(product) != owner.PRODUCT_KEYS:
        raise ValueError("product shape")
    product = dict(product)
    if not isinstance(product["price"], dict) or set(product["price"]) != {"amount_micros", "currency_code"}:
        raise ValueError("price")
    micros = product["price"]["amount_micros"]
    if not isinstance(micros, str) or not re.fullmatch("[1-9][0-9]{0,18}", micros):
        raise ValueError("micros")
    product["price"] = dict(product["price"], amount_micros=int(micros))
    profile = {"merchant_account_id": account, "data_source_name": datasource,
               "approved_offer_ids": [product["offer_id"]]}
    request = owner.build_insert_request(profile, product)
    request.product_input.version_number = int(generation)
    resource = f"accounts/{account}/products/{product['content_language']}~{product['feed_label']}~{product['offer_id']}"
    if intent["resource"] != resource:
        raise ValueError("resource binding")
    return request, resource


def bounded_adapter(fixture_origin: str = ""):
    from requests.adapters import HTTPAdapter

    if fixture_origin:
        target = urlsplit(fixture_origin)
        if target.scheme != "http" or target.hostname != "127.0.0.1" or not target.port or target.path or target.query or target.fragment or target.username:
            raise ValueError("fixture origin")

    class Adapter(HTTPAdapter):
        def send(self, request, **kwargs):
            origin = urlsplit(request.url)
            if origin.scheme != "https" or origin.netloc != "merchantapi.googleapis.com" or not origin.path.startswith("/products/v1/accounts/") or request.method not in {"GET", "POST"}:
                raise ValueError("provider origin/path")
            if fixture_origin:
                request.url = fixture_origin + origin.path + ("?" + origin.query if origin.query else "")
                request.headers["Authorization"] = "Bearer synthetic-merchant-sdk-fixture"
            kwargs.update(stream=True, timeout=8, proxies={})
            response = super().send(request, **kwargs)
            if 300 <= response.status_code < 400:
                response.close()
                raise ValueError("redirect denied")
            body = bytearray()
            try:
                for part in response.iter_content(chunk_size=4096):
                    if len(body) + len(part) > LIMIT:
                        raise ValueError("response bound")
                    body.extend(part)
            finally:
                response.close()
            response._content = bytes(body)
            response._content_consumed = True
            return response
    return Adapter(max_retries=0)


def clients(owner, mode: str, fixture_origin: str):
    if mode not in {"CREDENTIALS", "LOCAL_FIXTURES"} or (mode == "LOCAL_FIXTURES") != bool(fixture_origin):
        raise ValueError("runtime scope")
    kwargs = {"transport": "rest", "client_options": {"api_endpoint": "merchantapi.googleapis.com"}}
    if mode == "LOCAL_FIXTURES":
        from google.auth.credentials import AnonymousCredentials
        kwargs["credentials"] = AnonymousCredentials()
    inputs = owner.merchant.ProductInputsServiceClient(**kwargs)
    products = owner.merchant.ProductsServiceClient(**kwargs)
    for client in (inputs, products):
        # This is a pinned transport wrapper, not a modification to Google source.
        session = client.transport._session
        session.trust_env = False
        session.max_redirects = 0
        session._max_refresh_attempts = 0
        session.mount("https://", bounded_adapter(fixture_origin))
    return inputs, products


def execute(owner, command: dict[str, Any], inputs, products) -> dict[str, Any]:
    from google.api_core import exceptions
    request, resource = request_from_intent(owner, command["intent"])
    operation = command["operation"]
    if operation not in {"INSERT", "GET"}:
        raise ValueError("operation")
    try:
        if operation == "INSERT":
            result = inputs.insert_product_input(request, retry=None, timeout=8)
        else:
            result = products.get_product(owner.merchant.GetProductRequest(name=resource), retry=None, timeout=8)
        value = owner.message_dict(result)
        raw = json.dumps(value, ensure_ascii=False, sort_keys=True, separators=(",", ":")).encode("utf-8")
        if len(raw) > LIMIT:
            return {"state": "UNKNOWN", "operation": operation}
        return {"state": "OBSERVED", "operation": operation, "resource": resource,
                "response": value, "response_sha256": hashlib.sha256(raw).hexdigest()}
    except exceptions.NotFound:
        return {"state": "NOT_FOUND" if operation == "GET" else "UNKNOWN", "operation": operation}
    except (exceptions.BadRequest, exceptions.Unauthorized, exceptions.Forbidden,
            exceptions.Conflict, exceptions.Aborted):
        return {"state": "REJECTED" if operation == "INSERT" else "UNKNOWN", "operation": operation}
    except Exception:
        return {"state": "UNKNOWN", "operation": operation}


def main() -> int:
    result = {"state": "UNKNOWN"}
    try:
        if len(sys.argv) != 3:
            raise ValueError("source hash argument")
        admit_installed_sources(sys.argv[2])
        owner = load_owner(sys.argv[1])
        command = decode(sys.stdin.buffer.read(LIMIT + 1))
        if set(command) != {"operation", "intent", "mode", "fixture_origin"}:
            raise ValueError("command keys")
        # Validate the request before credential lookup or SDK client construction.
        request_from_intent(owner, command["intent"])
        inputs, products = clients(owner, command["mode"], command["fixture_origin"])
        try:
            result = execute(owner, command, inputs, products)
        finally:
            inputs.transport.close()
            products.transport.close()
    except Exception:
        result = {"state": "UNKNOWN"}
    raw = json.dumps(result, ensure_ascii=False, sort_keys=True, separators=(",", ":")).encode("utf-8")
    if len(raw) > LIMIT:
        raw = b'{"state":"UNKNOWN"}'
    sys.stdout.buffer.write(raw)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
````

### FILE: `google_merchant_product_sync/install_connected_runtime.py`

```yaml
block_id: "GO-CONNECTED-GOOGLE-MERCHANT:file11:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "ce27b54a17927d34d8683bb233e5e7c57d6c04a77d9c72d4b065e52af02a679a"
variables: []
secrets_allowed: false
```

````python
"""AUTHORED offline setup for the admitted Windows CPython3.14 wheel graph.

Run from the composed project with an already acquired exact wheel directory.
The absent destination is created once; failures preserve logs, never delete it.
No account or credential is read, and no provider API is contacted.
"""
from __future__ import annotations

import argparse
import hashlib
import json
import os
from pathlib import Path, PurePosixPath
import platform
import subprocess
import sys
from urllib.parse import urlsplit


def require(ok: bool, message: str) -> None:
    if not ok:
        raise ValueError(message)


def digest(path: Path) -> str:
    return hashlib.sha256(path.read_bytes()).hexdigest()


def setup(args: argparse.Namespace) -> dict:
    root = Path(__file__).resolve().parent
    require(sys.platform == "win32" and sys.version_info[:2] == (3, 14)
            and platform.machine().lower() in {"amd64", "x86_64"}, "fixed ABI required")
    require(digest(Path(sys.executable)) == args.python_sha256, "base Python pin mismatch")
    lock_path = root / "connected-runtime.lock.json"
    require(digest(lock_path) == args.lock_sha256, "runtime lock mismatch")
    lock = json.loads(lock_path.read_bytes())
    require(lock["schema"] == "elite-merchant-installed-source-lock/v1"
            and len(lock["wheels"]) == 20, "exact wheel graph required")
    require(args.destination.is_absolute() and not args.destination.exists(),
            "absolute absent runtime destination required")
    require(args.wheels.is_absolute() and args.wheels.is_dir(), "absolute wheel cache required")
    # Validate every artifact before any environment is created or pip is run.
    lines = []
    for row in lock["wheels"]:
        url = urlsplit(row["url"])
        require(url.scheme == "https" and url.netloc == "files.pythonhosted.org"
                and not url.query and not url.fragment, "official wheel origin required")
        name = PurePosixPath(url.path).name
        require(name.endswith(".whl") and "\\" not in name, "wheel filename required")
        wheel = args.wheels / name
        require(wheel.is_file() and wheel.stat().st_size == row["bytes"]
                and digest(wheel) == row["sha256"], "wheel pin mismatch: " + name)
        lines.append(row["distribution"] + " @ " + wheel.as_uri() + "#sha256=" + row["sha256"])
    args.destination.mkdir()
    env = {k: v for k, v in os.environ.items()
           if not k.startswith(("PIP_", "PYTHON", "GOOGLE_", "GCLOUD_", "CLOUDSDK_"))}
    env.update(PIP_CONFIG_FILE=os.devnull, PIP_DISABLE_PIP_VERSION_CHECK="1",
               PYTHONDONTWRITEBYTECODE="1")
    requirements = args.destination / "requirements-offline.lock"
    requirements.write_text("\n".join(lines) + "\n", encoding="utf-8", newline="\n")
    rows = []
    result = {"state": "FAIL", "live_account_access": False, "network_required": False,
              "base_python_sha256": args.python_sha256, "lock_sha256": args.lock_sha256}

    def run(name: str, command: list) -> None:
        path = args.destination / (name + ".log")
        with path.open("xb") as output:
            process = subprocess.run(list(map(str, command)), env=env, stdin=subprocess.DEVNULL,
                                     stdout=output, stderr=subprocess.STDOUT, timeout=180,
                                     creationflags=subprocess.CREATE_NO_WINDOW)
        rows.append({"step": name, "exit_code": process.returncode, "sha256": digest(path)})
        require(process.returncode == 0, "setup failed: " + name)

    try:
        runtime = args.destination / "venv"
        run("venv", [sys.executable, "-I", "-B", "-m", "venv", "--copies", runtime])
        python = runtime / "Scripts/python.exe"
        run("install", [python, "-I", "-B", "-m", "pip", "install", "--no-index", "--no-deps",
                        "--require-hashes", "-r", requirements])
        run("pip-check", [python, "-I", "-B", "-m", "pip", "check"])
        site = runtime / "Lib/site-packages"
        for name, sha in lock["files"].items():
            parts = PurePosixPath(name)
            require(not parts.is_absolute() and ".." not in parts.parts
                    and "\\" not in name and ":" not in name, "invalid installed path")
            require(digest(site / name) == sha, "installed source differs: " + name)
        script = root / "connected_worker.py"
        config = {"schema": "elite-merchant-runtime/v1", "python": str(python),
                  "python_sha256": digest(python), "script": str(script),
                  "script_sha256": digest(script), "owner_sha256": digest(root / "sync_product.py"),
                  "runtime_lock_sha256": args.lock_sha256,
                  "mode": "CREDENTIALS", "fixture_origin": ""}
        config_path = args.destination / "merchant-runtime.json"
        config_path.write_text(json.dumps(config, indent=2) + "\n", encoding="utf-8", newline="\n")
        result.update(state="PASS", installed_files=len(lock["files"]), wheels=len(lock["wheels"]),
                      runtime_file=str(config_path), runtime_sha256=digest(config_path),
                      condition="future user credential supplied through GOOGLE_APPLICATION_CREDENTIALS")
        return result
    finally:
        result["receipts"] = rows
        (args.destination / "setup-receipt.json").write_text(
            json.dumps(result, indent=2) + "\n", encoding="utf-8", newline="\n")


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--destination", type=Path, required=True)
    parser.add_argument("--wheels", type=Path, required=True)
    parser.add_argument("--python-sha256", required=True)
    parser.add_argument("--lock-sha256", required=True)
    args = parser.parse_args()
    print(json.dumps(setup(args), sort_keys=True))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
````

### FILE: `google_merchant_product_sync/test_connected_worker.py`

```yaml
block_id: "GO-CONNECTED-GOOGLE-MERCHANT:file12:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "b9922064a3449519122995c7f6ca20573b81ad6ad48de4d9e71619ab3a9c8da8"
variables: []
secrets_allowed: false
```

````python
from pathlib import Path
import json,subprocess,threading,http.server,urllib.parse,hashlib,os,sys
import argparse
parser=argparse.ArgumentParser();parser.add_argument("--python",type=Path,required=True);parser.add_argument("--receipt",type=Path,required=True)
args=parser.parse_args();require_receipt=args.receipt
if require_receipt.exists():raise ValueError("absent receipt required")
sha=lambda p:hashlib.sha256(p.read_bytes()).hexdigest()
root=Path(__file__).resolve().parent;python=args.python;worker=root/'connected_worker.py';owner=root/'sync_product.py';lock=root/'connected-runtime.lock.json'
state={'posts':0,'gets':0,'item':None,'mode':'normal','requests':[]}
class Handler(http.server.BaseHTTPRequestHandler):
 def log_message(self,*args):pass
 def respond(self,status,data):
  raw=json.dumps(data).encode();self.send_response(status);self.send_header('Content-Type','application/json');self.send_header('Content-Length',str(len(raw)));self.end_headers();self.wfile.write(raw)
 def do_POST(self):
  state['posts']+=1
  u=urllib.parse.urlsplit(self.path);q=urllib.parse.parse_qs(u.query)
  assert u.path=='/products/v1/accounts/123456/productInputs:insert'and q.get('dataSource')==['accounts/123456/dataSources/789012'],self.path
  assert self.headers.get('Authorization')=='Bearer synthetic-merchant-sdk-fixture'
  raw=self.rfile.read(int(self.headers['Content-Length']));product=json.loads(raw);state['requests'].append(product)
  assert product['offerId']=='reference-bicycle'and product['versionNumber']=='3'
  assert product['productAttributes']['price']['amountMicros']=='12345678901230000'
  if state['mode']=='unauthorized':return self.respond(401,{'error':{'code':401,'message':'synthetic authorization rejection','status':'UNAUTHENTICATED'}})
  if state['mode']=='redirect':
   self.send_response(307);self.send_header('Location',f'http://127.0.0.1:{self.server.server_port}/must-not-follow');self.end_headers();return
  if state['mode']=='oversize':return self.respond(200,{'unknown':'x'*32769})
  if state['mode']=='reject':return self.respond(400,{'error':{'code':400,'message':'synthetic rejection','status':'INVALID_ARGUMENT'}})
  item=dict(product,name='accounts/123456/productInputs/es~AR~reference-bicycle',product='accounts/123456/products/es~AR~reference-bicycle')
  state['item']=item
  if state['mode']=='lost':
   self.close_connection=True;self.connection.close();return
  return self.respond(200,item)
 def do_GET(self):
  state['gets']+=1
  u=urllib.parse.urlsplit(self.path)
  assert u.path=='/products/v1/accounts/123456/products/es~AR~reference-bicycle' and urllib.parse.parse_qs(u.query)=={'$alt':['json;enum-encoding=int']},self.path
  if state['item']is None:return self.respond(404,{'error':{'code':404,'message':'pending','status':'NOT_FOUND'}})
  product=dict(state['item']);product['name']=product.pop('product');product['dataSource']='accounts/123456/dataSources/789012'
  product['productStatus']={'destinationStatuses':[{'reportingContext':'SHOPPING_ADS','pendingCountries':['AR']}],'lastUpdateDate':'2026-09-13T00:00:00Z'}
  return self.respond(200,product)
server=http.server.ThreadingHTTPServer(('127.0.0.1',0),Handler);thread=threading.Thread(target=server.serve_forever,daemon=True);thread.start()
intent={'account_id':'123456','data_source':'accounts/123456/dataSources/789012','generation':'3','resource':'accounts/123456/products/es~AR~reference-bicycle','product':{'offer_id':'reference-bicycle','content_language':'es','feed_label':'AR','title':'Bicicleta de referencia','description':'Descripción sintética revisada del catálogo.','link':'https://catalog.example.invalid/catalog/bicycle','image_link':'https://catalog.example.invalid/api/public/catalog/media/'+'a'*64,'availability':'OUT_OF_STOCK','condition':'NEW','price':{'amount_micros':'12345678901230000','currency_code':'ARS'},'brand':'','gtins':[]}}
env={k:v for k,v in os.environ.items()if k in ['SYSTEMROOT','WINDIR','TEMP','TMP']};receipts=[]
def call(operation,value=None):
 command={'operation':operation,'intent':value or intent,'mode':'LOCAL_FIXTURES','fixture_origin':f'http://127.0.0.1:{server.server_port}'}
 q=subprocess.run([python,'-I','-S','-B',worker,sha(owner),sha(lock)],input=json.dumps(command).encode(),capture_output=True,env=env,timeout=25,creationflags=subprocess.CREATE_NO_WINDOW)
 assert q.returncode==0 and not q.stderr,q.stderr.decode()
 out=json.loads(q.stdout);receipts.append({'operation':operation,'state':out.get('state'),'result':out})
 return out
try:
 assert call('GET')['state']=='NOT_FOUND'
 out=call('INSERT');assert out['state']=='OBSERVED',out
 assert out['response']['version_number']=='3'
 assert out['response']['product_attributes']['price']['amount_micros']=='12345678901230000'
 out=call('GET');assert out['state']=='OBSERVED'and out['response']['data_source']==intent['data_source'],out
 assert out['response']['product_status']['destination_statuses'][0]['pending_countries']==['AR']
 state['mode']='lost';assert call('INSERT')['state']=='UNKNOWN'
 assert state['posts']==2
 assert call('GET')['state']=='OBSERVED'and state['posts']==2
 state['mode']='reject';assert call('INSERT')['state']=='REJECTED'
 tampered=json.loads(json.dumps(intent));tampered['resource']='accounts/999999/products/es~AR~reference-bicycle'
 assert call('INSERT',tampered)['state']=='UNKNOWN'and state['posts']==3
 for mode,want in [('unauthorized','REJECTED'),('redirect','UNKNOWN'),('oversize','UNKNOWN')]:
  before=state['posts'];state['mode']=mode;out=call('INSERT')
  assert out['state']==want and state['posts']==before+1,(mode,out,state['posts'])
 result={'state':'PASS','real_sdk':'1.8.0','posts':state['posts'],'gets':state['gets'],'requests':state['requests'],'receipts':receipts,'sources':{str(p):sha(p)for p in [worker,owner,lock]},'live_effects':False}
except BaseException as e:
 result={'state':'FAIL','failure':str(e),'type':type(e).__name__,'posts':state['posts'],'gets':state['gets'],'receipts':receipts};raise
finally:
 server.shutdown();server.server_close();thread.join(3)
 require_receipt.write_text(json.dumps(result,indent=2)+'\n',encoding='utf-8',newline='\n')
 print(json.dumps({k:result[k]for k in ['state','posts','gets']}))
````

### FILE: `internal/merchantbridge/client.go`

```yaml
block_id: "GO-CONNECTED-GOOGLE-MERCHANT:file13:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "6bbaed1739e223cf5baeddda6e9adc71e49430a5ed9e50f5ffadd9a498f8f566"
variables: []
secrets_allowed: false
```

````go
package merchantbridge

import (
	"context"
	cr "elite.local/enterprise/internal/catalogrelease"
	"elite.local/enterprise/internal/outbounddelivery"
	"encoding/json"
	"strconv"
	"time"
)

type Executor interface {
	Execute(context.Context, string, Intent) (Result, error)
}
type Client struct {
	profile  Profile
	executor Executor
}

func NewClient(p Profile, x Executor) (*Client, error) {
	if p.Validate() != nil || x == nil {
		return nil, ErrBinding
	}
	raw, _ := json.Marshal(p)
	var frozen Profile
	if json.Unmarshal(raw, &frozen) != nil {
		return nil, ErrBinding
	}
	return &Client{frozen, x}, nil
}

type Attributes struct {
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Link         string   `json:"link"`
	ImageLink    string   `json:"image_link"`
	Availability string   `json:"availability"`
	Condition    string   `json:"condition"`
	Price        Price    `json:"price"`
	Brand        string   `json:"brand"`
	GTINs        []string `json:"gtins"`
}
type Observation struct {
	Name          string          `json:"name"`
	Product       string          `json:"product"`
	DataSource    string          `json:"data_source"`
	OfferID       string          `json:"offer_id"`
	Language      string          `json:"content_language"`
	FeedLabel     string          `json:"feed_label"`
	Version       string          `json:"version_number"`
	Attributes    Attributes      `json:"product_attributes"`
	ProductStatus json.RawMessage `json:"product_status"`
}

func (c *Client) valid(r Intent) bool {
	b, e := c.profile.Binding(r.VariantID)
	return e == nil && r.PrepareRequest.Validate() == nil && r.TenantID == c.profile.TenantID && r.OrganizationID == c.profile.OrganizationID && r.ProfileSHA256 == c.profile.SHA256() && r.AccountID == c.profile.AccountID && r.DataSource == c.profile.DataSource() && r.Product.OfferID == b.OfferID && r.Resource == c.profile.Resource(b.OfferID) && cr.ValidSHA(r.SourceSHA256) && cr.ValidSHA(r.MediaSHA256)
}
func attributesMatch(p Product, a Attributes) bool {
	if a.Title != p.Title || a.Description != p.Description || a.Link != p.Link || a.ImageLink != p.ImageLink || a.Availability != p.Availability || a.Condition != p.Condition || a.Price != p.Price || a.Brand != p.Brand || len(a.GTINs) != len(p.GTINs) {
		return false
	}
	// Official GTIN order is not an identity guarantee; compare exact membership.
	seen := map[string]bool{}
	for _, g := range a.GTINs {
		if seen[g] {
			return false
		}
		seen[g] = true
	}
	for _, g := range p.GTINs {
		if !seen[g] {
			return false
		}
	}
	return true
}
func (c *Client) confirm(r Intent, result Result) (Observation, error) {
	var o Observation
	if result.State != "OBSERVED" || result.Resource != r.Resource || Hash(result.Response) != result.ResponseSHA256 || json.Unmarshal(result.Response, &o) != nil || o.OfferID != r.Product.OfferID || o.Language != r.Product.Language || o.FeedLabel != r.Product.FeedLabel || o.Version != strconv.FormatInt(r.Generation, 10) || !attributesMatch(r.Product, o.Attributes) {
		return o, ErrUnknown
	}
	if result.Operation == "INSERT" {
		expected := "accounts/" + r.AccountID + "/productInputs/" + r.Product.Language + "~" + r.Product.FeedLabel + "~" + r.Product.OfferID
		if o.Name != expected || o.Product != r.Resource {
			return o, ErrUnknown
		}
	} else if result.Operation == "GET" {
		if o.Name != r.Resource || o.DataSource != r.DataSource {
			return o, ErrUnknown
		}
	} else {
		return o, ErrBinding
	}
	return o, nil
}

type Record func(context.Context, Result) error

func (c *Client) Insert(ctx context.Context, r Intent, record Record) (outbounddelivery.Receipt, error) {
	if !c.valid(r) || !r.ExpiresAt.After(time.Now()) || record == nil {
		return outbounddelivery.Receipt{}, ErrBinding
	}
	result, e := c.executor.Execute(ctx, "INSERT", r)
	if e != nil {
		return outbounddelivery.Receipt{}, e
	}
	if result.State == "REJECTED" {
		_, h, _ := cr.Canonical(result)
		return outbounddelivery.Receipt{}, outbounddelivery.NewTerminalFailure("GOOGLE_MERCHANT_INSERT_REJECTED", h, nil)
	}
	if _, e = c.confirm(r, result); e != nil {
		return outbounddelivery.Receipt{}, e
	}
	if e = record(ctx, result); e != nil {
		return outbounddelivery.Receipt{}, ErrUnknown
	}
	return outbounddelivery.Receipt{ProviderMessageID: r.Resource, EvidenceSHA256: result.ResponseSHA256, AcceptedAt: time.Now().UTC()}, nil
}
func (c *Client) Observe(ctx context.Context, r Intent) (Result, error) {
	if !c.valid(r) {
		return Result{}, ErrBinding
	}
	result, e := c.executor.Execute(ctx, "GET", r)
	if e != nil {
		return result, e
	}
	// Preserve processing divergence/status separately; caller does not infer
	// provider approval or claim the input write refreshed freshness from GET.
	if _, e = c.confirm(r, result); e != nil {
		return result, e
	}
	return result, nil
}
````

### FILE: `internal/merchantbridge/contract.go`

```yaml
block_id: "GO-CONNECTED-GOOGLE-MERCHANT:file14:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "3c42c3c74fdc8ab964396bd91c514840c91d2954d538e6fc48350028adbacef0"
variables: []
secrets_allowed: false
```

````go
package merchantbridge

// AUTHORED provider representation and source binding. Catalog content, prices,
// inventory ATP, approvals and the outbound fence remain their original owners.
import (
	"bytes"
	"elite.local/enterprise/internal/approval"
	cr "elite.local/enterprise/internal/catalogrelease"
	"elite.local/enterprise/internal/channels"
	"encoding/json"
	"errors"
	"io"
	"net/url"
	"regexp"
	"strconv"
	"time"
)

var ErrBinding = errors.New("merchant: invalid or changed source binding")
var ErrUnknown = errors.New("merchant: provider result not confirmed")
var number = regexp.MustCompile(`^[0-9]{3,20}$`)
var offer = regexp.MustCompile(`^[A-Za-z0-9._-]{1,50}$`)

type Binding struct {
	VariantID string   `json:"variant_id"`
	OfferID   string   `json:"offer_id"`
	Condition string   `json:"condition"`
	Brand     string   `json:"brand"`
	GTINs     []string `json:"gtins"`
}
type Profile struct {
	Schema             string    `json:"schema"`
	TenantID           string    `json:"tenant_id"`
	OrganizationID     string    `json:"organization_id"`
	AccountID          string    `json:"account_id"`
	DataSourceID       string    `json:"data_source_id"`
	Language           string    `json:"content_language"`
	FeedLabel          string    `json:"feed_label"`
	Currency           string    `json:"currency"`
	CurrencyDigits     int       `json:"currency_digits"`
	Origin             string    `json:"origin"`
	RefreshCadenceDays int       `json:"refresh_cadence_days"`
	Bindings           []Binding `json:"bindings"`
}

func (p Profile) Validate() error {
	if p.RefreshCadenceDays < 1 || p.RefreshCadenceDays > 30 {
		return ErrBinding
	}
	u, e := url.Parse(p.Origin)
	if p.Schema != "elite-connected-merchant/v1" || !cr.ValidID(p.TenantID) || !cr.ValidID(p.OrganizationID) || !number.MatchString(p.AccountID) || !number.MatchString(p.DataSourceID) || !regexp.MustCompile(`^[a-z]{2}$`).MatchString(p.Language) || !regexp.MustCompile(`^[A-Z0-9-]{2,20}$`).MatchString(p.FeedLabel) || !regexp.MustCompile(`^[A-Z]{3}$`).MatchString(p.Currency) || p.CurrencyDigits < 0 || p.CurrencyDigits > 6 || e != nil || u.Scheme != "https" || u.Host == "" || u.User != nil || u.Path != "" || u.RawQuery != "" || u.Fragment != "" || len(p.Origin) > 512 || len(p.Bindings) < 1 || len(p.Bindings) > 32 {
		return ErrBinding
	}
	variants, offers := map[string]bool{}, map[string]bool{}
	for _, b := range p.Bindings {
		if !cr.ValidID(b.VariantID) || !offer.MatchString(b.OfferID) || variants[b.VariantID] || offers[b.OfferID] || b.Condition != "NEW" && b.Condition != "USED" && b.Condition != "REFURBISHED" || len(b.Brand) > 70 || b.Brand != "" && !cr.ValidText(b.Brand, 70) || len(b.GTINs) > 10 {
			return ErrBinding
		}
		variants[b.VariantID] = true
		offers[b.OfferID] = true
		seen := map[string]bool{}
		for _, g := range b.GTINs {
			if !regexp.MustCompile(`^(?:[0-9]{8}|[0-9]{12,14})$`).MatchString(g) || seen[g] {
				return ErrBinding
			}
			seen[g] = true
		}
	}
	if _, _, e = cr.Canonical(p); e != nil {
		return ErrBinding
	}
	return nil
}
func (p Profile) SHA256() string { _, h, _ := cr.Canonical(p); return h }
func (p Profile) DataSource() string {
	return "accounts/" + p.AccountID + "/dataSources/" + p.DataSourceID
}
func (p Profile) Binding(id string) (Binding, error) {
	for _, b := range p.Bindings {
		if b.VariantID == id {
			return b, nil
		}
	}
	return Binding{}, ErrBinding
}
func (p Profile) Resource(offerID string) string {
	return "accounts/" + p.AccountID + "/products/" + p.Language + "~" + p.FeedLabel + "~" + offerID
}

type PrepareRequest struct {
	ApprovalID string    `json:"approval_id"`
	Generation int64     `json:"generation,string"`
	VariantID  string    `json:"variant_id"`
	ExpiresAt  time.Time `json:"expires_at"`
}

func (r PrepareRequest) Validate() error {
	if !cr.ValidID(r.ApprovalID) || len(r.ApprovalID) < 16 || r.Generation < 1 || !cr.ValidID(r.VariantID) || r.ExpiresAt.IsZero() {
		return ErrBinding
	}
	return nil
}

type Price struct {
	AmountMicros string `json:"amount_micros"`
	Currency     string `json:"currency_code"`
}
type Product struct {
	OfferID      string   `json:"offer_id"`
	Language     string   `json:"content_language"`
	FeedLabel    string   `json:"feed_label"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Link         string   `json:"link"`
	ImageLink    string   `json:"image_link"`
	Availability string   `json:"availability"`
	Condition    string   `json:"condition"`
	Price        Price    `json:"price"`
	Brand        string   `json:"brand"`
	GTINs        []string `json:"gtins"`
}
type Intent struct {
	PrepareRequest
	TenantID       string    `json:"tenant_id"`
	OrganizationID string    `json:"organization_id"`
	ProfileSHA256  string    `json:"profile_sha256"`
	SourceSHA256   string    `json:"source_sha256"`
	MediaSHA256    string    `json:"media_sha256"`
	AccountID      string    `json:"account_id"`
	DataSource     string    `json:"data_source"`
	Resource       string    `json:"resource"`
	StockHorizon   time.Time `json:"stock_horizon"`
	Quantity       int64     `json:"quantity,string"`
	Product        Product   `json:"product"`
}

func (r Intent) Message() channels.Message {
	_, h, _ := cr.Canonical(r)
	return channels.Message{ChannelCode: "google_merchant", TenantID: r.TenantID, ExternalID: r.AccountID, ThreadID: r.Product.OfferID, Direction: channels.DirectionOut, DeliveryKey: r.ApprovalID, Text: h}
}
func Decode(raw []byte, v any) error {
	if len(raw) == 0 || len(raw) > 32768 {
		return ErrBinding
	}
	if _, _, e := approval.CanonicalPayload(raw); e != nil {
		return ErrBinding
	}
	d := json.NewDecoder(bytes.NewReader(raw))
	d.DisallowUnknownFields()
	if d.Decode(v) != nil || d.Decode(new(any)) != io.EOF {
		return ErrBinding
	}
	return nil
}
func Build(p Profile, pub cr.PublicDocument, r PrepareRequest, quantity int64, horizon time.Time) (Intent, error) {
	if p.Validate() != nil || r.Validate() != nil || pub.Generation != r.Generation || pub.Currency != p.Currency || !cr.ValidSHA(pub.SourceSHA256) || quantity < 0 || horizon.IsZero() {
		return Intent{}, ErrBinding
	}
	b, e := p.Binding(r.VariantID)
	if e != nil {
		return Intent{}, e
	}
	var v *cr.Variant
	var m *cr.Model
	for i := range pub.Variants {
		if pub.Variants[i].ID == r.VariantID {
			if v != nil {
				return Intent{}, ErrBinding
			}
			v = &pub.Variants[i]
		}
	}
	if v == nil || v.AmountMinorUnits <= 0 {
		return Intent{}, ErrBinding
	}
	for i := range pub.Models {
		if pub.Models[i].ID == v.ModelID {
			if m != nil {
				return Intent{}, ErrBinding
			}
			m = &pub.Models[i]
		}
	}
	if m == nil || !cr.ValidSHA(m.Media.SHA256) || !cr.ValidText(m.DisplayName, 150) {
		return Intent{}, ErrBinding
	}
	var specification map[string]json.RawMessage
	var description string
	if json.Unmarshal(m.Specification, &specification) != nil || json.Unmarshal(specification["merchant_description"], &description) != nil || !cr.ValidText(description, 5000) {
		return Intent{}, ErrBinding
	}
	// Representation only: convert original integer minor units into Google
	// integer micros without floating point, rounding, tax or exchange rules.
	scale := int64(1)
	for i := p.CurrencyDigits; i < 6; i++ {
		scale *= 10
	}
	if v.AmountMinorUnits > 1000000000000000000/scale {
		return Intent{}, ErrBinding
	}
	amount := v.AmountMinorUnits * scale
	availability := "OUT_OF_STOCK"
	if quantity > 0 {
		availability = "IN_STOCK"
	}
	gtins := append([]string{}, b.GTINs...)
	product := Product{b.OfferID, p.Language, p.FeedLabel, m.DisplayName, description, m.CanonicalURL, p.Origin + "/api/public/catalog/media/" + m.Media.SHA256, availability, b.Condition, Price{strconv.FormatInt(amount, 10), p.Currency}, b.Brand, gtins}
	link, e := url.Parse(product.Link)
	origin, _ := url.Parse(p.Origin)
	if e != nil || link.Scheme != "https" || link.Host != origin.Host || link.User != nil || link.RawQuery != "" || link.Fragment != "" {
		return Intent{}, ErrBinding
	}
	out := Intent{PrepareRequest: r, TenantID: p.TenantID, OrganizationID: p.OrganizationID, ProfileSHA256: p.SHA256(), SourceSHA256: pub.SourceSHA256, MediaSHA256: m.Media.SHA256, AccountID: p.AccountID, DataSource: p.DataSource(), Resource: p.Resource(b.OfferID), StockHorizon: horizon, Quantity: quantity, Product: product}
	if _, _, e = cr.Canonical(out); e != nil {
		return Intent{}, ErrBinding
	}
	return out, nil
}
````

### FILE: `internal/merchantbridge/fuzz_test.go`

```yaml
block_id: "GO-CONNECTED-GOOGLE-MERCHANT:file15:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "f910e24824dca117a014c4edce0f0175ea2ff3ffc59db37990c3058eb2c34fc7"
variables: []
secrets_allowed: false
```

````go
package merchantbridge

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func FuzzMerchantBoundedCommand(f *testing.F) {
	f.Add([]byte(`{"approval_id":"fixture-approval-0001","generation":"1","variant_id":"bicycle","expires_at":"2026-09-13T00:00:00Z"}`))
	f.Add([]byte(`{"generation":"1","Generation":"2"}`))
	f.Add([]byte(`{"price":1e9999}`))
	f.Add([]byte(`{"x":null,"x":true}`))
	f.Fuzz(func(t *testing.T, raw []byte) {
		var r PrepareRequest
		e := Decode(raw, &r)
		if len(raw) > 32768 && e == nil {
			t.Fatal("unbounded command")
		}
		if e == nil {
			b, e := json.Marshal(r)
			if e != nil {
				t.Fatal(e)
			}
			var again PrepareRequest
			if Decode(b, &again) != nil || again != r {
				t.Fatal("canonical roundtrip")
			}
			if r.Validate() == nil && (r.Generation < 1 || len(r.ApprovalID) < 16 || r.ExpiresAt == (time.Time{})) {
				t.Fatal("invalid admitted intent")
			}
		}
	})
}
func TestMerchantStrictCommandBoundary(t *testing.T) {
	for _, raw := range []string{`{"generation":"1","Generation":"2"}`, `{"generation":"1","generation":"2"}`, `{"generation":1}`, `{"generation":"1"} {}`, strings.Repeat(" ", 32769)} {
		if Decode([]byte(raw), new(PrepareRequest)) == nil {
			t.Fatal("invalid admitted", raw[:min(60, len(raw))])
		}
	}
}
````

### FILE: `internal/merchantbridge/process.go`

```yaml
block_id: "GO-CONNECTED-GOOGLE-MERCHANT:file16:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "1eac36b67462bdd2f849b1b39fef5f459fc3c252fb2ab4ded2cd5e64902c1050"
variables: []
secrets_allowed: false
```

````go
package merchantbridge

// AUTHORED bounded process binding following the admitted social/stored-value
// launchers. No shell, inherited proxies, Python paths or secrets in arguments.
import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Process struct {
	Python, PythonSHA256           string
	Script, ScriptSHA256           string
	OwnerSHA256, RuntimeLockSHA256 string
	Mode, FixtureOrigin            string
	CredentialsFile                string
}

func Hash(b []byte) string { s := sha256.Sum256(b); return hex.EncodeToString(s[:]) }
func exactFile(path, expected string, maximum int64) error {
	f, e := os.Open(path)
	if e != nil {
		return ErrBinding
	}
	defer f.Close()
	b, e := io.ReadAll(io.LimitReader(f, maximum+1))
	if e != nil || int64(len(b)) > maximum || len(expected) != 64 || Hash(b) != expected {
		return ErrBinding
	}
	return nil
}
func (p Process) Validate() error {
	if !filepath.IsAbs(p.Python) || !filepath.IsAbs(p.Script) || exactFile(p.Python, p.PythonSHA256, 4*1024*1024) != nil || exactFile(p.Script, p.ScriptSHA256, 65536) != nil || exactFile(filepath.Join(filepath.Dir(p.Script), "sync_product.py"), p.OwnerSHA256, 65536) != nil || exactFile(filepath.Join(filepath.Dir(p.Script), "connected-runtime.lock.json"), p.RuntimeLockSHA256, 1048576) != nil {
		return ErrBinding
	}
	if p.Mode == "LOCAL_FIXTURES" {
		u, e := url.Parse(p.FixtureOrigin)
		if e != nil || u.Scheme != "http" || u.Hostname() != "127.0.0.1" || u.Port() == "" || u.User != nil || u.Path != "" || u.RawQuery != "" || u.Fragment != "" || p.CredentialsFile != "" {
			return ErrBinding
		}
	} else if p.Mode == "CREDENTIALS" {
		if p.FixtureOrigin != "" || !filepath.IsAbs(p.CredentialsFile) {
			return ErrBinding
		}
		if info, e := os.Stat(p.CredentialsFile); e != nil || info.IsDir() {
			return ErrBinding
		}
	} else {
		return ErrBinding
	}
	return nil
}

type Result struct {
	State          string          `json:"state"`
	Operation      string          `json:"operation,omitempty"`
	Resource       string          `json:"resource,omitempty"`
	Response       json.RawMessage `json:"response,omitempty"`
	ResponseSHA256 string          `json:"response_sha256,omitempty"`
}
type output struct {
	bytes.Buffer
	exceeded bool
}

func (o *output) Write(b []byte) (int, error) {
	if o.Len()+len(b) > 32768 {
		o.exceeded = true
		return 0, io.ErrShortBuffer
	}
	return o.Buffer.Write(b)
}
func (p Process) Execute(ctx context.Context, operation string, r Intent) (Result, error) {
	if p.Validate() != nil || operation != "INSERT" && operation != "GET" {
		return Result{}, ErrBinding
	}
	input, e := json.Marshal(struct {
		Operation string `json:"operation"`
		Intent    Intent `json:"intent"`
		Mode      string `json:"mode"`
		Origin    string `json:"fixture_origin"`
	}{operation, r, p.Mode, p.FixtureOrigin})
	if e != nil || len(input) > 32768 {
		return Result{}, ErrBinding
	}
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	c := exec.CommandContext(ctx, p.Python, "-I", "-S", "-B", p.Script, p.OwnerSHA256, p.RuntimeLockSHA256)
	c.Stdin = bytes.NewReader(input)
	c.Dir = filepath.Dir(p.Script)
	for _, key := range []string{"SYSTEMROOT", "WINDIR", "TEMP", "TMP"} {
		if v := os.Getenv(key); v != "" {
			c.Env = append(c.Env, key+"="+v)
		}
	}
	if p.Mode == "CREDENTIALS" {
		if strings.ContainsAny(p.CredentialsFile, "\r\n\x00") {
			return Result{}, ErrBinding
		}
		c.Env = append(c.Env, "GOOGLE_APPLICATION_CREDENTIALS="+p.CredentialsFile)
	}
	var buffer output
	c.Stdout = &buffer
	c.Stderr = io.Discard
	if c.Run() != nil || buffer.exceeded {
		return Result{State: "UNKNOWN", Operation: operation}, nil
	}
	var result Result
	if Decode(buffer.Bytes(), &result) != nil {
		return Result{State: "UNKNOWN", Operation: operation}, nil
	}
	if result.State == "UNKNOWN" {
		return Result{State: "UNKNOWN", Operation: operation}, nil
	}
	if result.Operation != operation {
		return Result{}, ErrBinding
	}
	if result.State == "OBSERVED" {
		if result.Resource != r.Resource || len(result.Response) == 0 || Hash(result.Response) != result.ResponseSHA256 {
			return Result{}, ErrBinding
		}
	} else if result.State != "REJECTED" && result.State != "NOT_FOUND" {
		return Result{}, ErrBinding
	}
	return result, nil
}
````

### FILE: `internal/merchantbridge/queue.go`

```yaml
block_id: "GO-CONNECTED-GOOGLE-MERCHANT:file17:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "d384c717fbefc3d79cc878459a1dabc958175d9fcba0ed62be50c8189dd942bb"
variables: []
secrets_allowed: false
```

````go
package merchantbridge

// AUTHORED read-model classification over original approval/fence receipts.
import "time"

type QueueItem struct {
	VariantID            string     `json:"variant_id"`
	OfferID              string     `json:"offer_id"`
	Generation           int64      `json:"generation,string"`
	SourceSHA256         string     `json:"source_sha256"`
	DesiredProductSHA256 string     `json:"desired_product_sha256"`
	ApprovalID           string     `json:"approval_id"`
	ApprovalState        string     `json:"approval_state"`
	DeliveryState        string     `json:"delivery_state"`
	UnresolvedApprovalID string     `json:"unresolved_approval_id"`
	RefreshDueAt         *time.Time `json:"refresh_due_at"`
	Action               string     `json:"action"`
}

func QueueAction(now time.Time, item QueueItem, expires time.Time, changed bool) string {
	if item.UnresolvedApprovalID != "" {
		return "RECONCILE"
	}
	if item.ApprovalID == "" {
		return "INITIAL_APPROVAL_REQUIRED"
	}
	if item.DeliveryState == "failed_terminal" {
		return "PROVIDER_REJECTED"
	}
	if item.ApprovalState == "rejected" {
		return "REVIEW_REJECTED"
	}
	if changed {
		return "CHANGED_SOURCE_APPROVAL_REQUIRED"
	}
	if item.DeliveryState == "" {
		if !expires.After(now) {
			return "EXPIRED_APPROVAL_REQUIRED"
		}
		if item.ApprovalState == "approved" {
			return "SEND_APPROVED"
		}
		return "AWAITING_REVIEW"
	}
	if item.DeliveryState != "accepted" {
		return "RECONCILE"
	}
	if item.RefreshDueAt == nil {
		return "FRESHNESS_UNCONFIRMED"
	}
	if !item.RefreshDueAt.After(now) {
		return "REFRESH_APPROVAL_REQUIRED"
	}
	return "CURRENT_INPUT"
}
````

### FILE: `internal/merchantbridge/queue_test.go`

```yaml
block_id: "GO-CONNECTED-GOOGLE-MERCHANT:file18:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "d369bb2f19a62e92afe4f706d653be74d2f1b370ad6272fa4fbb43d938475c1b"
variables: []
secrets_allowed: false
```

````go
package merchantbridge

import (
	"testing"
	"time"
)

func TestMerchantQueueFreshnessAndApproval(t *testing.T) {
	now := time.Date(2026, 9, 13, 0, 0, 0, 0, time.UTC)
	past := now.Add(-time.Second)
	future := now.Add(time.Second)
	cases := []struct {
		name    string
		item    QueueItem
		expires time.Time
		changed bool
		want    string
	}{
		{"initial", QueueItem{}, future, false, "INITIAL_APPROVAL_REQUIRED"},
		{"ambiguous-prior-profile", QueueItem{UnresolvedApprovalID: "old"}, past, true, "RECONCILE"},
		{"review", QueueItem{ApprovalID: "a", ApprovalState: "pending"}, future, false, "AWAITING_REVIEW"},
		{"expired", QueueItem{ApprovalID: "a", ApprovalState: "approved"}, now, false, "EXPIRED_APPROVAL_REQUIRED"},
		{"approved", QueueItem{ApprovalID: "a", ApprovalState: "approved"}, future, false, "SEND_APPROVED"},
		{"source-changed", QueueItem{ApprovalID: "a", DeliveryState: "accepted"}, future, true, "CHANGED_SOURCE_APPROVAL_REQUIRED"},
		{"get-is-not-refresh", QueueItem{ApprovalID: "a", DeliveryState: "accepted"}, future, false, "FRESHNESS_UNCONFIRMED"},
		{"refresh-deadline", QueueItem{ApprovalID: "a", DeliveryState: "accepted", RefreshDueAt: &now}, past, false, "REFRESH_APPROVAL_REQUIRED"},
		{"current", QueueItem{ApprovalID: "a", DeliveryState: "accepted", RefreshDueAt: &future}, past, false, "CURRENT_INPUT"},
		{"provider-reject", QueueItem{ApprovalID: "a", DeliveryState: "failed_terminal", RefreshDueAt: &future}, future, false, "PROVIDER_REJECTED"},
		{"review-reject", QueueItem{ApprovalID: "a", ApprovalState: "rejected"}, future, false, "REVIEW_REJECTED"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := QueueAction(now, c.item, c.expires, c.changed); got != c.want {
				t.Fatalf("%s != %s", got, c.want)
			}
		})
	}
}
````

### FILE: `internal/merchantbridge/testfixture/server.go`

```yaml
block_id: "GO-CONNECTED-GOOGLE-MERCHANT:file19:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "0acc8cd56b0a37fae2fd9b997708da0d8432beef826f656874b25a2705b3ff4e"
variables: []
secrets_allowed: false
```

````go
package testfixture

// AUTHORED loopback transport fixture. The actual pinned Google SDK builds and
// sends these requests; this receiver never claims Google account approval.
import (
	mb "elite.local/enterprise/internal/merchantbridge"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
)

type Server struct {
	Mu                                           sync.Mutex
	HTTP                                         *httptest.Server
	Profile                                      mb.Profile
	Input                                        map[string]any
	Inserts, Gets                                int
	DropNext, RejectNext, Pending, ForeignSource bool
}

func New(p mb.Profile) *Server {
	s := &Server{Profile: p}
	s.HTTP = httptest.NewServer(http.HandlerFunc(s.serve))
	return s
}
func (s *Server) Close() { s.HTTP.Close() }
func (s *Server) serve(w http.ResponseWriter, r *http.Request) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	bad := func(status int, message string) {
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(map[string]any{"error": map[string]any{"code": status, "message": message}})
	}
	if r.Header.Get("Authorization") != "Bearer synthetic-merchant-sdk-fixture" || r.URL.Query().Get("$alt") != "json;enum-encoding=int" {
		bad(400, "fixture auth/query")
		return
	}
	if r.Method == "POST" && r.URL.Path == "/products/v1/accounts/"+s.Profile.AccountID+"/productInputs:insert" {
		s.Inserts++
		if r.URL.Query().Get("dataSource") != s.Profile.DataSource() {
			bad(400, "fixture data source")
			return
		}
		raw, e := io.ReadAll(io.LimitReader(r.Body, 32769))
		if e != nil || len(raw) > 32768 {
			bad(413, "fixture body")
			return
		}
		var input map[string]any
		if json.Unmarshal(raw, &input) != nil {
			bad(400, "fixture JSON")
			return
		}
		if input["offerId"] != s.Profile.Bindings[0].OfferID || input["contentLanguage"] != s.Profile.Language || input["feedLabel"] != s.Profile.FeedLabel || input["versionNumber"] != "1" {
			bad(400, "fixture key/version")
			return
		}
		if s.RejectNext {
			s.RejectNext = false
			bad(400, "fixture provider rejection")
			return
		}
		input["name"] = "accounts/" + s.Profile.AccountID + "/productInputs/" + s.Profile.Language + "~" + s.Profile.FeedLabel + "~" + s.Profile.Bindings[0].OfferID
		input["product"] = s.Profile.Resource(s.Profile.Bindings[0].OfferID)
		s.Input = input
		if s.DropNext {
			s.DropNext = false
			conn, _, _ := w.(http.Hijacker).Hijack()
			conn.Close()
			return
		}
		json.NewEncoder(w).Encode(input)
		return
	}
	if r.Method == "GET" && r.URL.Path == "/products/v1/"+s.Profile.Resource(s.Profile.Bindings[0].OfferID) {
		s.Gets++
		if s.Pending || s.Input == nil {
			bad(404, "fixture processing pending")
			return
		}
		data := map[string]any{}
		for k, v := range s.Input {
			data[k] = v
		}
		delete(data, "product")
		data["name"] = s.Profile.Resource(s.Profile.Bindings[0].OfferID)
		data["dataSource"] = s.Profile.DataSource()
		if s.ForeignSource {
			data["dataSource"] = "accounts/999999/dataSources/999999"
		}
		data["productStatus"] = map[string]any{"destinationStatuses": []any{map[string]any{"reportingContext": "SHOPPING_ADS", "pendingCountries": []string{"AR"}}}, "lastUpdateDate": "2026-09-13T00:00:00Z"}
		json.NewEncoder(w).Encode(data)
		return
	}
	bad(404, "fixture path")
}
````

### FILE: `internal/platform/httpapi/merchant.go`

```yaml
block_id: "GO-CONNECTED-GOOGLE-MERCHANT:file20:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "5448981334301f9016674a47d5f057a6a8c94675b6644f1f01b3a0d70fd9b76e"
variables: []
secrets_allowed: false
```

````go
package httpapi

// AUTHORED bounded role HTTP adapter. All authority comes from the verifier.
import (
	"context"
	"io"
	"net/http"
	"time"

	"elite.local/enterprise/internal/approval"
	cr "elite.local/enterprise/internal/catalogrelease"
	mb "elite.local/enterprise/internal/merchantbridge"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/postgres"
)

type MerchantPublicationService interface {
	Queue(context.Context, identity.Principal) (postgres.MerchantQueue, error)
	Prepare(context.Context, identity.Principal, mb.PrepareRequest) (mb.Intent, error)
	Submit(context.Context, identity.Principal, mb.Intent) (string, bool, error)
	Decide(context.Context, identity.Principal, string, string, bool, string) (approval.State, error)
	Status(context.Context, identity.Principal, string) (postgres.MerchantStatus, error)
	Send(context.Context, identity.Principal, string) error
	Reconcile(context.Context, identity.Principal, string) error
}
type MerchantModule struct {
	Service                  MerchantPublicationService
	TenantID, OrganizationID string
}

func (m MerchantModule) Register(mux *http.ServeMux, verifier identity.Verifier) {
	if m.Service == nil {
		return
	}
	route := func(method, path, permission string, action func(http.ResponseWriter, *http.Request, identity.Principal)) {
		mux.HandleFunc(method+" /v1/admin/merchant"+path, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", "no-store")
			p, e := authenticate(r.Context(), r.Header.Get("Authorization"), verifier)
			if e != nil {
				writeProblem(w, 401, "UNAUTHENTICATED", "valid bearer required")
				return
			}
			if p.TenantID != m.TenantID || !p.Allowed(permission) || !p.AllowedOrganization(m.OrganizationID) {
				writeProblem(w, 403, "FORBIDDEN", "scope required")
				return
			}
			if r.URL.RawQuery != "" || r.PathValue("id") != "" && !cr.ValidID(r.PathValue("id")) {
				writeProblem(w, 400, "INVALID_TARGET", "exact target required")
				return
			}
			ctx, cancel := context.WithTimeout(r.Context(), 25*time.Second)
			defer cancel()
			action(w, r.WithContext(ctx), p)
		})
	}
	decode := func(w http.ResponseWriter, r *http.Request, v any) bool {
		if r.Header.Get("Content-Type") != "application/json" {
			writeProblem(w, 415, "UNSUPPORTED_MEDIA_TYPE", "application/json required")
			return false
		}
		if e := http.NewResponseController(w).SetReadDeadline(time.Now().Add(5 * time.Second)); e != nil {
			writeProblem(w, 503, "READ_BOUNDARY", "bounded read unavailable")
			return false
		}
		raw, e := io.ReadAll(http.MaxBytesReader(w, r.Body, 32768))
		if e != nil || mb.Decode(raw, v) != nil {
			writeProblem(w, 400, "INVALID_BODY", "exact bounded JSON required")
			return false
		}
		return true
	}
	failed := func(w http.ResponseWriter, e error) bool {
		if e == nil {
			return false
		}
		writeProblem(w, 409, "MERCHANT_UNCONFIRMED", "consult immutable request and provider reconciliation")
		return true
	}
	route("GET", "/queue", "merchant:read", func(w http.ResponseWriter, r *http.Request, p identity.Principal) {
		out, e := m.Service.Queue(r.Context(), p)
		if !failed(w, e) {
			writeJSON(w, 200, out)
		}
	})
	route("POST", "/prepare", "merchant:request", func(w http.ResponseWriter, r *http.Request, p identity.Principal) {
		var in mb.PrepareRequest
		if !decode(w, r, &in) {
			return
		}
		out, e := m.Service.Prepare(r.Context(), p, in)
		if !failed(w, e) {
			writeJSON(w, 200, out)
		}
	})
	route("POST", "/requests", "merchant:request", func(w http.ResponseWriter, r *http.Request, p identity.Principal) {
		var in mb.Intent
		if !decode(w, r, &in) {
			return
		}
		hash, replay, e := m.Service.Submit(r.Context(), p, in)
		if !failed(w, e) {
			writeJSON(w, 200, map[string]any{"approval_id": in.ApprovalID, "request_sha256": hash, "replay": replay})
		}
	})
	route("POST", "/requests/{id}/decision", "merchant:approve", func(w http.ResponseWriter, r *http.Request, p identity.Principal) {
		var in struct {
			Hash    string `json:"request_sha256"`
			Approve *bool  `json:"approve"`
			Reason  string `json:"reason"`
		}
		if !decode(w, r, &in) {
			return
		}
		if in.Approve == nil || !cr.ValidSHA(in.Hash) {
			writeProblem(w, 400, "INVALID_DECISION", "explicit bound decision required")
			return
		}
		out, e := m.Service.Decide(r.Context(), p, r.PathValue("id"), in.Hash, *in.Approve, in.Reason)
		if !failed(w, e) {
			writeJSON(w, 200, map[string]any{"state": out})
		}
	})
	route("GET", "/requests/{id}", "merchant:read", func(w http.ResponseWriter, r *http.Request, p identity.Principal) {
		out, e := m.Service.Status(r.Context(), p, r.PathValue("id"))
		if !failed(w, e) {
			writeJSON(w, 200, out)
		}
	})
	for _, v := range []struct {
		name, permission string
		call             func(context.Context, identity.Principal, string) error
	}{{"send", "merchant:send", m.Service.Send}, {"reconcile", "merchant:reconcile", m.Service.Reconcile}} {
		route("POST", "/requests/{id}/"+v.name, v.permission, func(w http.ResponseWriter, r *http.Request, p identity.Principal) {
			var empty struct{}
			if !decode(w, r, &empty) {
				return
			}
			if !failed(w, v.call(r.Context(), p, r.PathValue("id"))) {
				writeJSON(w, 200, map[string]string{"outcome": "confirmed"})
			}
		})
	}
}
````

### FILE: `internal/platform/postgres/merchant_observation.go`

```yaml
block_id: "GO-CONNECTED-GOOGLE-MERCHANT:file21:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "68349f6d1f6ac0a7301b389b0300369dd8fc0181ff2480df376b1e9646ebb9c2"
variables: []
secrets_allowed: false
```

````go
package postgres

// AUTHORED provider evidence binding. Raw SDK bytes survive PostgreSQL JSONB
// whitespace changes; observed processing is separate from input acknowledgement.
import (
	"context"
	mb "elite.local/enterprise/internal/merchantbridge"
)

func (s *MerchantPublication) recordMerchantObservation(ctx context.Context, r mb.Intent, result mb.Result, confirmed bool) error {
	if result.Operation != "INSERT" && result.Operation != "GET" || result.State != "OBSERVED" && result.State != "NOT_FOUND" && result.State != "REJECTED" && result.State != "UNKNOWN" || confirmed && result.State != "OBSERVED" || result.Operation == "INSERT" && !confirmed {
		return mb.ErrBinding
	}
	raw := []byte("{}")
	if result.State == "OBSERVED" {
		if result.Resource != r.Resource || len(result.Response) == 0 || mb.Hash(result.Response) != result.ResponseSHA256 {
			return mb.ErrBinding
		}
		raw = result.Response
	}
	command, e := s.catalog.pool.Exec(ctx, `insert into catalog.merchant_observation(tenant_id,approval_id,operation,result_state,confirmed,response,response_sha256)
 select i.tenant_id,i.delivery_key,$3,$4,$5,$6,$7 from catalog.merchant_effect i
 join communication.outbound_delivery d using(tenant_id,channel_code,delivery_key)
 where i.tenant_id=$1 and i.delivery_key=$2 and i.profile_sha256=$8 and i.source_sha256=$9
 and d.state in ('sending','unknown','accepted')
 on conflict(tenant_id,approval_id,operation,result_state,response_sha256) do nothing`, r.TenantID, r.ApprovalID, result.Operation, result.State, confirmed, raw, mb.Hash(raw), r.ProfileSHA256, r.SourceSHA256)
	if e != nil {
		return e
	}
	if command.RowsAffected() == 1 {
		return nil
	}
	var exists bool
	e = s.catalog.pool.QueryRow(ctx, `select exists(select 1 from catalog.merchant_observation o join catalog.merchant_effect i
 on i.tenant_id=o.tenant_id and i.channel_code=o.channel_code and i.delivery_key=o.approval_id
 where o.tenant_id=$1 and o.approval_id=$2 and o.operation=$3 and o.result_state=$4 and o.confirmed=$5
 and o.response_sha256=$6 and i.profile_sha256=$7 and i.source_sha256=$8)`, r.TenantID, r.ApprovalID, result.Operation, result.State, confirmed, mb.Hash(raw), r.ProfileSHA256, r.SourceSHA256).Scan(&exists)
	if e != nil {
		return e
	}
	if !exists {
		return mb.ErrBinding
	}
	return nil
}
````

### FILE: `internal/platform/postgres/merchant_publication.go`

```yaml
block_id: "GO-CONNECTED-GOOGLE-MERCHANT:file22:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "ae93b4fa50113da85bdb740f93a4e68c960e889462fc30c2067e74b7658bcf7f"
variables: []
secrets_allowed: false
```

````go
package postgres

// AUTHORED composition of immutable catalog, existing ATP, human approvals and
// the shared outbound fence. No new monetary/inventory algorithm or ledger.
import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"elite.local/enterprise/internal/approval"
	cr "elite.local/enterprise/internal/catalogrelease"
	"elite.local/enterprise/internal/channels"
	mb "elite.local/enterprise/internal/merchantbridge"
	"elite.local/enterprise/internal/outbounddelivery"
	"elite.local/enterprise/internal/platform/identity"
	"github.com/jackc/pgx/v5"
)

type MerchantPublication struct {
	catalog   *CatalogRelease
	profile   mb.Profile
	approvals *HumanApprovals
	fence     *OutboundDeliveryStore
	client    *mb.Client
}

func NewMerchantPublication(catalog *CatalogRelease, p mb.Profile, key []byte, executor mb.Executor) (*MerchantPublication, error) {
	if catalog == nil || p.Validate() != nil || p.TenantID != catalog.profile.TenantID || p.OrganizationID != catalog.profile.OrganizationID || p.Currency != catalog.profile.Currency {
		return nil, mb.ErrBinding
	}
	raw, _ := json.Marshal(p)
	var frozen mb.Profile
	if json.Unmarshal(raw, &frozen) != nil {
		return nil, mb.ErrBinding
	}
	client, e := mb.NewClient(frozen, executor)
	if e != nil {
		return nil, e
	}
	fence, e := NewOutboundDeliveryStore(catalog.pool, key, 30*time.Second)
	if e != nil {
		return nil, e
	}
	return &MerchantPublication{catalog, frozen, NewHumanApprovals(catalog.pool), fence, client}, nil
}
func (s *MerchantPublication) authorized(p identity.Principal, permission string) bool {
	return s != nil && approvalPrincipal(p, s.profile.TenantID, s.profile.OrganizationID, permission)
}
func (s *MerchantPublication) source(ctx context.Context, tx pgx.Tx, generation int64, variant string, horizon time.Time) (cr.PublicDocument, int64, error) {
	var publication cr.Publication
	var quantity int64
	if _, e := tx.Exec(ctx, `select pg_advisory_xact_lock(hashtextextended($1,0))`, "catalog-release:"+s.profile.TenantID); e != nil {
		return cr.PublicDocument{}, 0, e
	}
	e := tx.QueryRow(ctx, `select p.generation,d.snapshot_sha256,d.snapshot_canonical,p.effective_price_book_id from catalog.release_publication p
 join catalog.release_draft d using(tenant_id,draft_id)
 join catalog.release_public_current c on c.tenant_id=p.tenant_id and c.generation=p.generation
 join org.organization o on o.tenant_id=p.tenant_id and o.organization_id=$3
 join platform.tenant t on t.tenant_id=p.tenant_id
 where p.tenant_id=$1 and p.generation=$2 and o.status='active' and t.status='active'
 for share of o,t`, s.profile.TenantID, generation, s.profile.OrganizationID).Scan(&publication.Generation, &publication.SHA256, &publication.Snapshot, &publication.EffectivePriceBookID)
	if e != nil {
		return cr.PublicDocument{}, 0, mb.ErrBinding
	}
	e = tx.QueryRow(ctx, `select available_to_promise from inventory.serial_atp($1,$2,$3,$4)`, s.profile.TenantID, s.profile.OrganizationID, variant, horizon).Scan(&quantity)
	if e != nil || quantity < 0 {
		return cr.PublicDocument{}, 0, mb.ErrBinding
	}
	doc, e := cr.ProjectPublic(publication)
	return doc, quantity, e
}

func (s *MerchantPublication) Prepare(ctx context.Context, p identity.Principal, r mb.PrepareRequest) (mb.Intent, error) {
	if !s.authorized(p, "merchant:request") || r.Validate() != nil {
		return mb.Intent{}, mb.ErrBinding
	}
	tx, e := s.catalog.pool.Begin(ctx)
	if e != nil {
		return mb.Intent{}, e
	}
	defer tx.Rollback(ctx)
	var now time.Time
	if e = tx.QueryRow(ctx, `select clock_timestamp()`).Scan(&now); e != nil {
		return mb.Intent{}, e
	}
	if !r.ExpiresAt.After(now) || r.ExpiresAt.After(now.Add(15*time.Minute)) {
		return mb.Intent{}, mb.ErrBinding
	}
	doc, q, e := s.source(ctx, tx, r.Generation, r.VariantID, now)
	if e != nil {
		return mb.Intent{}, e
	}
	intent, e := mb.Build(s.profile, doc, r, q, now)
	if e != nil {
		return mb.Intent{}, e
	}
	return intent, tx.Commit(ctx)
}
func (s *MerchantPublication) guard(ctx context.Context, tx pgx.Tx, r mb.Intent, requireCurrent bool) error {
	if r.TenantID != s.profile.TenantID || r.OrganizationID != s.profile.OrganizationID || r.ProfileSHA256 != s.profile.SHA256() {
		return mb.ErrBinding
	}
	if !requireCurrent {
		return nil
	}
	var now time.Time
	if e := tx.QueryRow(ctx, `select clock_timestamp()`).Scan(&now); e != nil {
		return e
	}
	if !r.ExpiresAt.After(now) || r.ExpiresAt.After(r.StockHorizon.Add(15*time.Minute)) || r.StockHorizon.After(now) {
		return mb.ErrBinding
	}
	doc, q, e := s.source(ctx, tx, r.Generation, r.VariantID, r.StockHorizon)
	if e != nil {
		return e
	}
	rebuilt, e := mb.Build(s.profile, doc, r.PrepareRequest, q, r.StockHorizon)
	if e != nil {
		return e
	}
	_, a, e := cr.Canonical(rebuilt)
	if e != nil {
		return e
	}
	_, b, e := cr.Canonical(r)
	if e != nil || a != b {
		return mb.ErrBinding
	}
	return nil
}
func (s *MerchantPublication) Submit(ctx context.Context, p identity.Principal, r mb.Intent) (string, bool, error) {
	if !s.authorized(p, "merchant:request") {
		return "", false, mb.ErrBinding
	}
	raw, hash, e := cr.Canonical(r)
	if e != nil {
		return "", false, e
	}
	spec := HumanApprovalSpec{Request: approval.Request{TenantID: p.TenantID, ID: r.ApprovalID, Kind: approval.KindMarketplaceMutation, SubjectID: r.Product.OfferID, Requester: p.Subject, EvidenceSHA: hash}, OrganizationID: s.profile.OrganizationID, Payload: raw}
	replay, e := s.approvals.Submit(ctx, p, spec, "merchant:request", func(ctx context.Context, tx pgx.Tx) error { return s.guard(ctx, tx, r, true) })
	return hash, replay, e
}
func (s *MerchantPublication) request(ctx context.Context, id string) (mb.Intent, string, error) {
	var r mb.Intent
	var raw []byte
	var hash string
	e := s.catalog.pool.QueryRow(ctx, `select payload,evidence_sha from approval.request where tenant_id=$1 and request_id=$2 and organization_id=$3 and kind='marketplace_mutation'`, s.profile.TenantID, id, s.profile.OrganizationID).Scan(&raw, &hash)
	if e != nil || mb.Decode(raw, &r) != nil || r.ApprovalID != id || r.ProfileSHA256 != s.profile.SHA256() {
		return r, "", mb.ErrBinding
	}
	_, h, e := cr.Canonical(r)
	if e != nil || h != hash {
		return r, "", mb.ErrBinding
	}
	return r, hash, nil
}
func (s *MerchantPublication) Decide(ctx context.Context, p identity.Principal, id, hash string, approve bool, reason string) (approval.State, error) {
	if !s.authorized(p, "merchant:approve") {
		return "", mb.ErrBinding
	}
	r, stored, e := s.request(ctx, id)
	if e != nil || stored != hash {
		return "", mb.ErrBinding
	}
	return s.approvals.Decide(ctx, p, s.profile.TenantID, id, s.profile.OrganizationID, hash, approve, reason, "merchant:approve", func(ctx context.Context, tx pgx.Tx) error {
		if !approve {
			return nil
		}
		return s.guard(ctx, tx, r, true)
	})
}

type MerchantStatus struct {
	InputAcknowledged   bool            `json:"input_acknowledged"`
	RefreshDueAt        *time.Time      `json:"refresh_due_at,omitempty"`
	ProcessingState     string          `json:"processing_state"`
	ProcessingConfirmed bool            `json:"processing_confirmed"`
	ProcessingResponse  json.RawMessage `json:"processing_response,omitempty"`
	ApprovalID          string          `json:"approval_id"`
	ApprovalState       string          `json:"approval_state"`
	RequestSHA256       string          `json:"request_sha256"`
	Payload             json.RawMessage `json:"payload"`
	DeliveryState       string          `json:"delivery_state"`
	FailureCode         string          `json:"failure_code"`
	EvidenceSHA256      string          `json:"evidence_sha256"`
}

func (s *MerchantPublication) Status(ctx context.Context, p identity.Principal, id string) (MerchantStatus, error) {
	var out MerchantStatus
	if !s.authorized(p, "merchant:read") || !cr.ValidID(id) {
		return out, mb.ErrBinding
	}
	e := s.catalog.pool.QueryRow(ctx, `select a.request_id,a.state,a.evidence_sha,a.payload,coalesce(d.state,''),coalesce(d.failure_code,''),coalesce(d.evidence_sha256_hex,'')
 from approval.request a left join communication.outbound_delivery d on d.tenant_id=a.tenant_id and d.channel_code='google_merchant' and d.delivery_key=a.request_id
 where a.tenant_id=$1 and a.request_id=$2 and a.organization_id=$3 and a.kind='marketplace_mutation' and a.payload->>'profile_sha256'=$4`, s.profile.TenantID, id, s.profile.OrganizationID, s.profile.SHA256()).Scan(&out.ApprovalID, &out.ApprovalState, &out.RequestSHA256, &out.Payload, &out.DeliveryState, &out.FailureCode, &out.EvidenceSHA256)
	if e != nil {
		return out, e
	}
	e = s.catalog.pool.QueryRow(ctx, `select exists(select 1 from catalog.merchant_observation where tenant_id=$1 and approval_id=$2 and operation='INSERT' and confirmed),
 (select max(observed_at)+make_interval(days=>$3) from catalog.merchant_observation where tenant_id=$1 and approval_id=$2 and operation='INSERT' and confirmed)`, s.profile.TenantID, id, s.profile.RefreshCadenceDays).Scan(&out.InputAcknowledged, &out.RefreshDueAt)
	if e != nil {
		return out, e
	}
	var raw []byte
	var hash string
	e = s.catalog.pool.QueryRow(ctx, `select result_state,confirmed,response,response_sha256 from catalog.merchant_observation
 where tenant_id=$1 and approval_id=$2 and operation='GET' order by observed_at desc,response_sha256 limit 1`, s.profile.TenantID, id).Scan(&out.ProcessingState, &out.ProcessingConfirmed, &raw, &hash)
	if errors.Is(e, pgx.ErrNoRows) {
		out.ProcessingState = "UNOBSERVED"
		return out, nil
	}
	if e != nil {
		return out, e
	}
	if mb.Hash(raw) != hash {
		return out, mb.ErrBinding
	}
	out.ProcessingResponse = raw
	return out, nil
}

type merchantBound struct {
	service *MerchantPublication
	intent  mb.Intent
	hash    string
}

func (b *merchantBound) Claim(ctx context.Context, m channels.Message, h string) (outbounddelivery.Claim, error) {
	return b.service.fence.claimWithAdmission(ctx, m, h, func(ctx context.Context, tx pgx.Tx) error {
		r := b.intent
		s := b.service
		// One active effect per tenant/provider SKU, including another approval id.
		if _, e := tx.Exec(ctx, `select pg_advisory_xact_lock(hashtextextended($1,0))`, "merchant:"+r.TenantID+":"+r.AccountID+":"+r.Product.OfferID); e != nil {
			return e
		}
		var busy bool

		e := tx.QueryRow(ctx, `select exists(select 1 from catalog.merchant_effect i
  join communication.outbound_delivery d using(tenant_id,channel_code,delivery_key)
  where i.tenant_id=$1 and i.account_id=$2 and i.offer_id=$3 and d.state in ('sending','unknown'))`, r.TenantID, r.AccountID, r.Product.OfferID).Scan(&busy)
		if e != nil {
			return e
		}
		if busy {
			return mb.ErrUnknown
		}
		if e = s.guard(ctx, tx, r, true); e != nil {
			return e
		}
		if _, e = LookupHumanApproval(ctx, tx, r.TenantID, r.ApprovalID, r.OrganizationID, approval.KindMarketplaceMutation, b.hash); e != nil {
			return e
		}
		_, e = tx.Exec(ctx, `insert into catalog.merchant_effect(tenant_id,channel_code,delivery_key,account_id,offer_id,generation,source_sha256,profile_sha256,request_sha256) values($1,'google_merchant',$2,$3,$4,$5,$6,$7,$8)`, r.TenantID, r.ApprovalID, r.AccountID, r.Product.OfferID, r.Generation, r.SourceSHA256, r.ProfileSHA256, b.hash)
		return e
	})
}
func (b *merchantBound) Complete(c context.Context, m channels.Message, h string, r outbounddelivery.Receipt) error {
	return b.service.fence.Complete(c, m, h, r)
}
func (b *merchantBound) MarkUnknown(c context.Context, m channels.Message, h, code string) error {
	return b.service.fence.MarkUnknown(c, m, h, code)
}
func (b *merchantBound) MarkFailed(c context.Context, m channels.Message, h, sha, code string) error {
	return b.service.fence.MarkFailed(c, m, h, sha, code)
}

type merchantSender struct {
	client  *mb.Client
	intent  mb.Intent
	service *MerchantPublication
}

func (s merchantSender) SendWithReceipt(ctx context.Context, m channels.Message) (outbounddelivery.Receipt, error) {
	want, _ := outbounddelivery.MessageSHA256(s.intent.Message())
	got, e := outbounddelivery.MessageSHA256(m)
	if e != nil || got != want {
		return outbounddelivery.Receipt{}, mb.ErrBinding
	}

	return s.client.Insert(ctx, s.intent, func(ctx context.Context, result mb.Result) error {
		return s.service.recordMerchantObservation(ctx, s.intent, result, true)
	})
}

type merchantOutbound struct{}

func (merchantOutbound) Code() string { return "google_merchant" }
func (merchantOutbound) Receive(context.Context) ([]channels.Message, error) {
	return nil, mb.ErrBinding
}
func (merchantOutbound) Send(context.Context, channels.Message) error { return mb.ErrBinding }
func (s *MerchantPublication) Send(ctx context.Context, p identity.Principal, id string) error {
	if !s.authorized(p, "merchant:send") {
		return mb.ErrBinding
	}
	r, h, e := s.request(ctx, id)
	if e != nil {
		return e
	}
	channel := outbounddelivery.Channel{CodeValue: "google_merchant", Receiver: merchantOutbound{}, Sender: merchantSender{s.client, r, s}, Store: &merchantBound{s, r, h}}
	return channel.Send(ctx, r.Message())
}

func (s *MerchantPublication) Reconcile(ctx context.Context, p identity.Principal, id string) error {
	if !s.authorized(p, "merchant:reconcile") {
		return mb.ErrBinding
	}
	r, _, e := s.request(ctx, id)
	if e != nil {
		return e
	}
	status, e := s.Status(ctx, p, id)
	if e != nil {
		return e
	}
	if status.DeliveryState != "accepted" && status.DeliveryState != "unknown" && status.DeliveryState != "sending" {
		return mb.ErrBinding
	}
	message := r.Message()
	hash, e := outbounddelivery.MessageSHA256(message)
	if e != nil {
		return e
	}
	if status.DeliveryState != "accepted" {
		if _, e = s.fence.Claim(ctx, message, hash); !errors.Is(e, outbounddelivery.ErrUnknown) {
			return mb.ErrUnknown
		}
	}
	observation, observeErr := s.client.Observe(ctx, r)
	if e = s.recordMerchantObservation(ctx, r, observation, observeErr == nil); e != nil {
		return e
	}
	if observeErr != nil {
		return observeErr
	}
	if status.DeliveryState == "accepted" {
		return nil
	}
	receipt := outbounddelivery.Receipt{ProviderMessageID: r.Resource, EvidenceSHA256: observation.ResponseSHA256, AcceptedAt: time.Now().UTC()}
	return s.fence.ReconcileAccepted(ctx, message, hash, receipt)
}
````

### FILE: `internal/platform/postgres/merchant_publication_integration_test.go`

```yaml
block_id: "GO-CONNECTED-GOOGLE-MERCHANT:file23:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "2505cdf6b62cdefa04801eb596d3f7fe71902f67fed9d30f50d02dbc972883aa"
variables: []
secrets_allowed: false
```

````go
package postgres_test

// AUTHORED HTTP/PG/provider fixture for initial catalog publication only.
import (
	"bytes"
	"context"
	cr "elite.local/enterprise/internal/catalogrelease"
	"elite.local/enterprise/internal/electromobility"
	mb "elite.local/enterprise/internal/merchantbridge"
	fixture "elite.local/enterprise/internal/merchantbridge/testfixture"
	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/platform/identity"
	db "elite.local/enterprise/internal/platform/postgres"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"image"
	"image/png"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestMerchantCatalogConnected(t *testing.T) {
	database := os.Getenv("PAYMENT_CONNECTED_DB_URL")
	if database == "" {
		t.Skip("owned fixture DB required")
	}
	ctx := context.Background()
	pool, e := pgxpool.New(ctx, database)
	if e != nil {
		t.Fatal(e)
	}
	defer pool.Close()
	tenant := "d84c64f9-1252-4db5-bf5a-3af8026190f1"
	org := "marketplace-store"
	for _, q := range []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1,'marketplace-fixture','Fixture','Fixture')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'marketplace-store','marketplace-store','Fixture','store')`,
	} {
		if _, e = pool.Exec(ctx, q, tenant); e != nil {
			t.Fatal(e)
		}
	}
	principal := func(subject string) identity.Principal {
		return identity.Principal{TenantID: tenant, Subject: subject, Permissions: map[string]struct{}{"catalog:draft": {}, "catalog:read": {}, "catalog:publish": {}, "catalog:review:legal": {}, "catalog:review:technical": {}, "catalog:review:media": {}, "catalog:review:publication": {}, "merchant:request": {}, "merchant:approve": {}, "merchant:read": {}, "merchant:send": {}, "merchant:reconcile": {}}, Organizations: map[string]struct{}{org: {}}}
	}
	maker, reviewer := principal("marketplace-maker"), principal("marketplace-reviewer")
	cp := cr.Profile{Schema: "elite-catalog-publication/v1", TenantID: tenant, OrganizationID: org, Market: "AR", Currency: "ARS", Origin: "https://catalog.example.invalid", PolicyCode: cr.PolicyCode}
	catalog, e := db.NewCatalogRelease(pool, cp, db.NewCommerce(pool))
	if e != nil {
		t.Fatal(e)
	}
	source, e := catalog.CreateSource(ctx, maker, cr.SourceRequest{CommandID: "marketplace-source", Model: electromobility.Model{Code: "fixture-bicycle", DisplayName: "Fixture bicycle", VehicleClass: "bicycle", Specification: json.RawMessage(`{"merchant_description":"Descripción sintética revisada del catálogo."}`)}, Variants: []cr.SourceVariant{{Code: "fixture-bicycle-standard", DisplayName: "Standard", BatterySpecification: json.RawMessage("{}"), AmountMinorUnits: 1234567890123, TaxMode: "not-applicable"}}, ValidFrom: time.Now().Add(-time.Hour)})
	if e != nil {
		t.Fatal("source", e)
	}
	var pngRaw bytes.Buffer
	if e = png.Encode(&pngRaw, image.NewNRGBA(image.Rect(0, 0, 2, 2))); e != nil {
		t.Fatal(e)
	}
	media, e := catalog.UploadPNG(ctx, maker, "marketplace-media", pngRaw.Bytes())
	if e != nil {
		t.Fatal(e)
	}
	draft, e := catalog.CreateDraft(ctx, maker, cr.DraftRequest{CommandID: "marketplace-draft", PriceBookID: source.SourcePriceBookID, Models: []cr.ModelReference{{ModelID: source.ResourceID, MediaID: media.ResourceID}}})
	if e != nil {
		t.Fatal("draft", e)
	}
	for _, stage := range []string{"legal", "technical", "media", "publication"} {
		_, e = catalog.Review(ctx, reviewer, cr.ReviewRequest{CommandID: "marketplace-review-" + stage, DraftID: draft.ResourceID, Stage: stage, SnapshotSHA256: draft.SnapshotSHA256, Approved: true, Reason: "synthetic source review", EvidenceSHA256: strings.Repeat("a", 64)})
		if e != nil {
			t.Fatal("review", e)
		}
	}
	published, e := catalog.Publish(ctx, reviewer, cr.PublishRequest{CommandID: "marketplace-publish", DraftID: draft.ResourceID, SnapshotSHA256: draft.SnapshotSHA256, ExpectedGeneration: 0, Reason: "fixture publication"})
	if e != nil {
		t.Fatal("publish", e)
	}

	profile := mb.Profile{Schema: "elite-connected-merchant/v1", TenantID: tenant, OrganizationID: org, AccountID: "123456", DataSourceID: "789012", Language: "es", FeedLabel: "AR", Currency: "ARS", CurrencyDigits: 2, Origin: cp.Origin, RefreshCadenceDays: 30, Bindings: []mb.Binding{{VariantID: source.SourceVariantIDs[0], OfferID: "fixture-bicycle-standard", Condition: "NEW", GTINs: []string{}}}}
	receiver := fixture.New(profile)
	defer receiver.Close()
	root := os.Getenv("MERCHANT_TEST_ROOT")
	python := os.Getenv("MERCHANT_TEST_PYTHON")
	fileHash := func(path string) string {
		t.Helper()
		raw, e := os.ReadFile(path)
		if e != nil {
			t.Fatal(e)
		}
		return mb.Hash(raw)
	}
	script := filepath.Join(root, "google_merchant_product_sync", "connected_worker.py")
	process := mb.Process{Python: python, PythonSHA256: fileHash(python), Script: script, ScriptSHA256: fileHash(script), OwnerSHA256: fileHash(filepath.Join(filepath.Dir(script), "sync_product.py")), RuntimeLockSHA256: fileHash(filepath.Join(filepath.Dir(script), "connected-runtime.lock.json")), Mode: "LOCAL_FIXTURES", FixtureOrigin: receiver.HTTP.URL}
	service, e := db.NewMerchantPublication(catalog, profile, bytes.Repeat([]byte("f"), 32), process)
	if e != nil {
		t.Fatal(e)
	}

	forbidden := principal("unprivileged")
	forbidden.Permissions = map[string]struct{}{}
	wrongOrg := principal("wrong-org")
	wrongOrg.Organizations = map[string]struct{}{"other": {}}
	wrongTenant := principal("wrong-tenant")
	wrongTenant.TenantID = "33b9e628-aa94-4c70-9c24-6ac4a3fb805e"
	verifier := marketplaceFixtureVerifier{"maker": maker, "reviewer": reviewer, "forbidden": forbidden, "wrong-org": wrongOrg, "wrong-tenant": wrongTenant}
	mux := http.NewServeMux()
	httpapi.MerchantModule{Service: service, TenantID: tenant, OrganizationID: org}.Register(mux, verifier)
	var apiMu sync.Mutex
	dropAPI := false
	apiPosts := 0
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recorder := httptest.NewRecorder()
		mux.ServeHTTP(marketplaceRecorder{recorder, w}, r)
		apiMu.Lock()
		defer apiMu.Unlock()
		if r.Method == "POST" {
			apiPosts++
		}
		if dropAPI && strings.HasSuffix(r.URL.Path, "/send") && recorder.Code == 200 {
			dropAPI = false
			conn, _, e := w.(http.Hijacker).Hijack()
			if e == nil {
				conn.Close()
			}
			return
		}
		for k, values := range recorder.Header() {
			for _, v := range values {
				w.Header().Add(k, v)
			}
		}
		w.WriteHeader(recorder.Code)
		_, _ = w.Write(recorder.Body.Bytes())
	}))
	defer api.Close()
	call := func(token, method, path string, in, out any) (int, error) {
		var raw []byte
		if in != nil {
			raw, _ = json.Marshal(in)
		}
		request, _ := http.NewRequest(method, api.URL+"/v1/admin/merchant"+path, bytes.NewReader(raw))
		request.Header.Set("Authorization", "Bearer "+token)
		request.Header.Set("Content-Type", "application/json")
		response, e := api.Client().Do(request)
		if e != nil {
			return 0, e
		}
		defer response.Body.Close()
		data, e := io.ReadAll(io.LimitReader(response.Body, 65537))
		if e != nil {
			return response.StatusCode, e
		}
		if response.Header.Get("Cache-Control") != "no-store" {
			return response.StatusCode, errors.New("private response cacheable")
		}
		if response.StatusCode >= 400 {
			return response.StatusCode, fmt.Errorf("HTTP %d %s", response.StatusCode, data)
		}
		if out != nil {
			e = json.Unmarshal(data, out)
		}
		return response.StatusCode, e
	}

	prepare := func(id string) mb.Intent {
		t.Helper()
		var r mb.Intent
		_, e := call("maker", "POST", "/prepare", mb.PrepareRequest{ApprovalID: id, Generation: published.Generation, VariantID: profile.Bindings[0].VariantID, ExpiresAt: time.Now().Add(5 * time.Minute)}, &r)
		if e != nil {
			t.Fatal("prepare", e)
		}
		if r.SourceSHA256 != published.SnapshotSHA256 || r.Product.Price.AmountMicros != "12345678901230000" || r.Quantity != 0 || r.Product.Availability != "OUT_OF_STOCK" {
			t.Fatal("source/amount/ATP", r)
		}
		return r
	}

	submit := func(r mb.Intent) string {
		t.Helper()
		var result struct {
			Hash string `json:"request_sha256"`
		}
		if _, e := call("maker", "POST", "/requests", r, &result); e != nil {
			t.Fatal("submit", e)
		}
		return result.Hash
	}
	approve := func(r mb.Intent, hash string) {
		t.Helper()
		body := map[string]any{"request_sha256": hash, "approve": true, "reason": "exact synthetic provider mutation"}
		if status, e := call("maker", "POST", "/requests/"+r.ApprovalID+"/decision", body, nil); e == nil || status != 409 {
			t.Fatal("self approval")
		}
		if _, e := call("reviewer", "POST", "/requests/"+r.ApprovalID+"/decision", body, nil); e != nil {
			t.Fatal("review", e)
		}
	}
	status := func(id, want string) {
		t.Helper()
		var result db.MerchantStatus
		if _, e := call("maker", "GET", "/requests/"+id, nil, &result); e != nil || result.DeliveryState != want {
			t.Fatal("durable state", result.ApprovalID, result.ApprovalState, result.DeliveryState, result.FailureCode, "provider writes", receiver.Inserts, e)
		}
	}

	first := prepare("merchant-lost-input")
	for _, token := range []string{"forbidden", "wrong-org", "wrong-tenant"} {
		if code, e := call(token, "POST", "/requests", first, nil); e == nil || code != 403 {
			t.Fatal("scope", token, e)
		}
	}
	tampered := first
	tampered.Product.Price.AmountMicros = "1"
	if _, e := call("maker", "POST", "/requests", tampered, nil); e == nil {
		t.Fatal("client price override")
	}
	approve(first, submit(first))
	receiver.Mu.Lock()
	receiver.DropNext = true
	receiver.Mu.Unlock()
	if _, e := call("maker", "POST", "/requests/"+first.ApprovalID+"/send", struct{}{}, nil); e == nil {
		t.Fatal("lost SDK acknowledgement")
	}
	status(first.ApprovalID, "unknown")
	if receiver.Inserts != 1 {
		t.Fatal("SDK retries", receiver.Inserts)
	}
	receiver.Mu.Lock()
	receiver.Pending = true
	receiver.Mu.Unlock()
	if _, e := call("maker", "POST", "/requests/"+first.ApprovalID+"/reconcile", struct{}{}, nil); e == nil {
		t.Fatal("pending inferred success")
	}
	receiver.Mu.Lock()
	receiver.Pending = false
	receiver.ForeignSource = true
	receiver.Mu.Unlock()
	if _, e := call("maker", "POST", "/requests/"+first.ApprovalID+"/reconcile", struct{}{}, nil); e == nil {
		t.Fatal("foreign data source accepted")
	}
	receiver.Mu.Lock()
	receiver.ForeignSource = false
	receiver.Mu.Unlock()
	if _, e := call("maker", "POST", "/requests/"+first.ApprovalID+"/reconcile", struct{}{}, nil); e != nil {
		t.Fatal("GET source recovery", e)
	}
	status(first.ApprovalID, "accepted")
	var observed db.MerchantStatus
	if _, e := call("maker", "GET", "/requests/"+first.ApprovalID, nil, &observed); e != nil {
		t.Fatal(e)
	}
	if observed.InputAcknowledged || observed.RefreshDueAt != nil || !observed.ProcessingConfirmed || !bytes.Contains(observed.ProcessingResponse, []byte("pending_countries")) {
		t.Fatal("GET invented input refresh or Google approval", observed)
	}
	// A separately approved refresh retains the official source generation.
	refresh := prepare("merchant-approved-refresh")
	approve(refresh, submit(refresh))
	apiMu.Lock()
	dropAPI = true
	apiMu.Unlock()
	if _, e := call("maker", "POST", "/requests/"+refresh.ApprovalID+"/send", struct{}{}, nil); e == nil {
		t.Fatal("lost local API acknowledgement")
	}
	status(refresh.ApprovalID, "accepted")
	var wg sync.WaitGroup
	for range 12 {
		wg.Go(func() { call("maker", "POST", "/requests/"+refresh.ApprovalID+"/send", struct{}{}, nil) })
	}
	wg.Wait()
	if receiver.Inserts != 2 {
		t.Fatal("refresh replay", receiver.Inserts)
	}
	if _, e := call("maker", "POST", "/requests/"+refresh.ApprovalID+"/reconcile", struct{}{}, nil); e != nil {
		t.Fatal("processing status", e)
	}
	if _, e := call("maker", "GET", "/requests/"+refresh.ApprovalID, nil, &observed); e != nil || !observed.InputAcknowledged || observed.RefreshDueAt == nil || !observed.ProcessingConfirmed {
		t.Fatal("input/status separation", e)
	}
	rejected := prepare("merchant-provider-rejected")
	approve(rejected, submit(rejected))
	receiver.Mu.Lock()
	receiver.RejectNext = true
	receiver.Mu.Unlock()
	if _, e := call("maker", "POST", "/requests/"+rejected.ApprovalID+"/send", struct{}{}, nil); e == nil {
		t.Fatal("provider rejection")
	}
	status(rejected.ApprovalID, "failed_terminal")
	if _, e := call("maker", "POST", "/requests/"+rejected.ApprovalID+"/send", struct{}{}, nil); e == nil || receiver.Inserts != 3 {
		t.Fatal("rejected retry")
	}
	for _, q := range []string{
		`update catalog.merchant_observation set response='{}'::bytea where tenant_id=$1`,
		`delete from catalog.merchant_observation where tenant_id=$1`,
		`update catalog.merchant_effect set offer_id='forged' where tenant_id=$1`,
	} {
		if _, e := pool.Exec(ctx, q, tenant); e == nil {
			t.Fatal("immutable provider evidence")
		}
	}
	var approvals, effects, observations, attempts int
	for _, v := range []struct {
		sql string
		out *int
	}{
		{`select count(*) from approval.request where tenant_id=$1 and kind='marketplace_mutation'`, &approvals},
		{`select count(*) from catalog.merchant_effect where tenant_id=$1`, &effects},
		{`select count(*) from catalog.merchant_observation where tenant_id=$1`, &observations},
		{`select sum(attempt_count) from communication.outbound_delivery where tenant_id=$1 and channel_code='google_merchant'`, &attempts},
	} {
		if e := pool.QueryRow(ctx, v.sql, tenant).Scan(v.out); e != nil {
			t.Fatal(e)
		}
	}
	if approvals != 3 || effects != 3 || observations != 5 || attempts != 3 {
		t.Fatal("durable counts", approvals, effects, observations, attempts)
	}
	t.Logf("MERCHANT_CONNECTED_PASS approvals=%d effects=%d observations=%d attempts=%d actual_sdk_inserts=%d actual_sdk_gets=%d api_posts=%d source=%s", approvals, effects, observations, attempts, receiver.Inserts, receiver.Gets, apiPosts, published.SnapshotSHA256)
}
````

### FILE: `internal/platform/postgres/merchant_queue.go`

```yaml
block_id: "GO-CONNECTED-GOOGLE-MERCHANT:file24:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "226cde2e36c3e517a3fda33ff4b1acc71fb46e8eeb6696beab78d4aa0b019dbb"
variables: []
secrets_allowed: false
```

````go
package postgres

// AUTHORED bounded read-only projection. No queue read creates approvals,
// provider writes or jobs; each refresh follows the existing manual workflow.
import (
	"context"
	cr "elite.local/enterprise/internal/catalogrelease"
	mb "elite.local/enterprise/internal/merchantbridge"
	"elite.local/enterprise/internal/platform/identity"
	"errors"
	"github.com/jackc/pgx/v5"
	"time"
)

type MerchantQueue struct {
	ObservedAt    time.Time      `json:"observed_at"`
	ProfileSHA256 string         `json:"profile_sha256"`
	Items         []mb.QueueItem `json:"items"`
}

func (s *MerchantPublication) Queue(ctx context.Context, p identity.Principal) (MerchantQueue, error) {
	var out MerchantQueue
	if !s.authorized(p, "merchant:read") {
		return out, mb.ErrBinding
	}
	tx, e := s.catalog.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead})
	if e != nil {
		return out, e
	}
	defer tx.Rollback(ctx)
	var generation int64
	if e = tx.QueryRow(ctx, `select clock_timestamp(),generation from catalog.release_public_current where tenant_id=$1`, s.profile.TenantID).Scan(&out.ObservedAt, &generation); e != nil {
		return out, mb.ErrBinding
	}
	out.ProfileSHA256 = s.profile.SHA256()
	out.Items = make([]mb.QueueItem, 0, len(s.profile.Bindings))
	for _, b := range s.profile.Bindings {
		horizon := out.ObservedAt.Add(24 * time.Hour)
		doc, q, e := s.source(ctx, tx, generation, b.VariantID, horizon)
		if e != nil {
			return out, e
		}
		desired, e := mb.Build(s.profile, doc, mb.PrepareRequest{ApprovalID: "queue-projection-only", Generation: generation, VariantID: b.VariantID, ExpiresAt: out.ObservedAt.Add(time.Minute)}, q, horizon)
		if e != nil {
			return out, e
		}
		_, desiredSHA, e := cr.Canonical(desired.Product)
		if e != nil {
			return out, e
		}
		item := mb.QueueItem{VariantID: b.VariantID, OfferID: b.OfferID, Generation: generation, SourceSHA256: doc.SourceSHA256, DesiredProductSHA256: desiredSHA}
		var raw []byte
		var latest mb.Intent
		e = tx.QueryRow(ctx, `select a.request_id,a.state,coalesce(d.state,''),a.payload from approval.request a
 left join communication.outbound_delivery d on d.tenant_id=a.tenant_id and d.delivery_key=a.request_id and d.channel_code='google_merchant'
 where a.tenant_id=$1 and a.organization_id=$2 and a.kind='marketplace_mutation' and a.payload->>'profile_sha256'=$3
 and a.payload->'product'->>'offer_id'=$4 order by a.created_at desc,a.request_id desc limit 1`,
			s.profile.TenantID, s.profile.OrganizationID, s.profile.SHA256(), b.OfferID).Scan(&item.ApprovalID, &item.ApprovalState, &item.DeliveryState, &raw)
		if e != nil && !errors.Is(e, pgx.ErrNoRows) {
			return out, e
		}
		changed := false
		if e == nil {
			if mb.Decode(raw, &latest) != nil {
				return out, mb.ErrBinding
			}
			_, h, e := cr.Canonical(latest.Product)
			if e != nil {
				return out, e
			}
			changed = latest.SourceSHA256 != doc.SourceSHA256 || latest.Generation != generation || h != desiredSHA
		}
		e = tx.QueryRow(ctx, `select coalesce((select i.delivery_key from catalog.merchant_effect i
 join communication.outbound_delivery d using(tenant_id,channel_code,delivery_key)
 where i.tenant_id=$1 and i.account_id=$2 and i.offer_id=$3 and d.state in('sending','unknown')
 order by i.delivery_key limit 1),''),
 (select max(o.observed_at)+make_interval(days=>$5) from catalog.merchant_observation o
 join catalog.merchant_effect i on i.tenant_id=o.tenant_id and i.delivery_key=o.approval_id
 where i.tenant_id=$1 and i.account_id=$2 and i.offer_id=$3 and i.profile_sha256=$4
 and o.operation='INSERT' and o.confirmed)`,
			s.profile.TenantID, s.profile.AccountID, b.OfferID, s.profile.SHA256(), s.profile.RefreshCadenceDays).Scan(&item.UnresolvedApprovalID, &item.RefreshDueAt)
		if e != nil {
			return out, e
		}
		item.Action = mb.QueueAction(out.ObservedAt, item, latest.ExpiresAt, changed)
		out.Items = append(out.Items, item)
	}
	return out, tx.Commit(ctx)
}
````

## 6. Configuration surface

docs/MERCHANT_CONNECTED_REFERENCE.md, config/merchant.reference.json and offline install_connected_runtime.py; no credential in a materialized file.

## 7. Dependency bill

Unchanged PYTHON-GOOGLE-MERCHANT-PRODUCT-SYNC-ADAPTER SDK/artifact/20-wheel lock. New code and source-lock metadata are AUTHORED glue; Google dependency remains DEPENDENCY_PIN Apache-2.0, with all original dependency notices preserved.

## 8. Apply order

Original catalog, supply/ATP, human approval, outbound fence, Google SDK and application owners. Apply migration0082; enabled host fails closed without exact configuration/guards.

## 9. Verification

Actual REST SDK/API/PG proof3approvals/3effects/5observations, lost provider/API replies,12 replays, source/scope/tamper/pending/datasource negatives; queue/host; populated downgrade refusal/empty reversal; independent offline runtime1041files;6SDKPOST/3GET bounded transport;finite2s fuzz/unit/vet/build; OSV2.5.1 exact20/20 zero findings.

## 10. Reconstruction evidence

reconstruction_evidence/MERCHANT_CONNECTED_RELEASE_V402.md/json; fixed composition rebuild receipts. Other selected T2805 mappings/comms and later ordered controls remain open.

