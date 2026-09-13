# Connected Mercado Libre catalog mutations

## 1. Metadata

```yaml
pack_id: "GO-CONNECTED-MARKETPLACE-MUTATION"
pack_version: "0.3.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: RECONSTRUCTIBLE
  admission: CONDITIONED
claim: "T2805 CONTENT from approved catalog and accepted media: exact name/attributes/picture, single unsold User Product item guard, original human approval/fence and GET recovery. Initial publication and existing mutations retained. Local fixtures proven; remaining feeds/comms open. AUTHORED glue, no corporate authorship."
stacks: ["Go 1.26.8", "PostgreSQL 18.6", "Node.js 24.20.0 where frontend selected"]
compatible_with: ["GO-CONNECTED-CATALOG-PUBLICATION", "GO-SUPPLY-FACTORY-INVENTORY-API", "GO-HUMAN-APPROVAL-CORE", "GO-PG-OUTBOUND-DELIVERY-FENCE", "GO-APPROVED-STORED-VALUE-TENDER", "GO-ELECTROMOBILITY-APPLICATION"]
incompatible_with: ["unbound tenant/provider/account", "production certification inferred from fixtures", "implicit source or business-policy attribution"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["docs/marketplace/official-contract.lock.json"]
verified_at: "2026-09-13"
```

## 2. Applicability

Library local fixture claim for initial approved-PNG User Products publication and existing seller-owned item mutations. No production-live or whole-T2805 closure.

## 3. Architecture contract

Exact current source amount and original ATP snapshot; no floating point or tax/stock algorithm. Provider-specific location/version, human separation, immutable request and one-attempt fence.

## 4. Exact file manifest

```text
CREATE cmd/electromobility-api/marketplace.go
CREATE cmd/electromobility-api/marketplace_test.go
CREATE config/marketplace.reference.json
CREATE db/migrations/0079_marketplace_mutation.down.sql
CREATE db/migrations/0079_marketplace_mutation.up.sql
CREATE docs/MARKETPLACE_MUTATION_REFERENCE.md
CREATE docs/marketplace/official-contract.lock.json
CREATE internal/marketplacebridge/approval_test.go
CREATE internal/marketplacebridge/contract.go
CREATE internal/marketplacebridge/fuzz_test.go
CREATE internal/marketplacebridge/http.go
CREATE internal/marketplacebridge/http_test.go
CREATE internal/marketplacebridge/testfixture/receiver.go
CREATE internal/platform/httpapi/marketplace.go
CREATE internal/platform/postgres/marketplace_publication.go
CREATE internal/platform/postgres/marketplace_publication_integration_test.go
CREATE db/migrations/0080_marketplace_initial.down.sql
CREATE db/migrations/0080_marketplace_initial.up.sql
CREATE docs/MARKETPLACE_INITIAL_REFERENCE.md
CREATE internal/marketplacebridge/initial.go
CREATE internal/marketplacebridge/initial_test.go
CREATE internal/marketplacebridge/testfixture/initial.go
CREATE internal/platform/postgres/marketplace_initial.go
CREATE internal/platform/postgres/marketplace_initial_integration_test.go
CREATE db/migrations/0081_marketplace_content.down.sql
CREATE db/migrations/0081_marketplace_content.up.sql
CREATE docs/MARKETPLACE_CONTENT_REFERENCE.md
CREATE internal/marketplacebridge/content.go
CREATE internal/marketplacebridge/content_test.go
CREATE internal/marketplacebridge/testfixture/content.go
CREATE internal/platform/postgres/marketplace_content_integration_test.go
```

## 5. Materialization blocks

### FILE: `cmd/electromobility-api/marketplace.go`

```yaml
block_id: "GO-CONNECTED-MARKETPLACE-MUTATION:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "9371743a2a5ef1ed21d405c43204a949df1abf8e676ed50d27fcbdedd9bedd93"
variables: []
secrets_allowed: false
```

````go
package main

// AUTHORED activation glue. Credentials stay in the process environment and
// provider calls use the fixed official origin with no redirect or proxy.
import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"strings"

	cr "elite.local/enterprise/internal/catalogrelease"
	mb "elite.local/enterprise/internal/marketplacebridge"
	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/platform/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

var errMarketplaceConfiguration = errors.New("marketplace activation invalid")

type marketplaceToken struct{ lookup func(string) string }

func (t marketplaceToken) MercadoLibreAccessToken(context.Context) (string, error) {
	return t.lookup("MERCADOLIBRE_ACCESS_TOKEN"), nil
}
func init() { marketplaceModuleFactory = selectedMarketplaceModule }
func selectedMarketplaceModule(ctx context.Context, pool *pgxpool.Pool, price *postgres.Commerce, lookup func(string) string) (httpapi.EnterpriseModule, error) {
	if lookup == nil {
		return nil, errMarketplaceConfiguration
	}
	switch lookup("MARKETPLACE_ENABLED") {
	case "", "false":
		return nil, nil
	case "true":
	default:
		return nil, errMarketplaceConfiguration
	}
	if lookup("CATALOG_RELEASE_ENABLED") != "true" || pool == nil || price == nil {
		return nil, errMarketplaceConfiguration
	}
	f, e := os.Open(lookup("MARKETPLACE_PROFILE_FILE"))
	if e != nil {
		return nil, errMarketplaceConfiguration
	}
	defer f.Close()
	raw, e := io.ReadAll(io.LimitReader(f, 32769))
	if e != nil || len(raw) > 32768 {
		return nil, errMarketplaceConfiguration
	}
	sum := sha256.Sum256(raw)
	if hex.EncodeToString(sum[:]) != lookup("MARKETPLACE_PROFILE_SHA256") {
		return nil, errMarketplaceConfiguration
	}
	var profile mb.Profile
	if mb.Decode(raw, &profile) != nil || profile.Validate() != nil {
		return nil, errMarketplaceConfiguration
	}
	cp, e := cr.LoadProfile(lookup("CATALOG_RELEASE_PROFILE_FILE"), lookup("CATALOG_RELEASE_PROFILE_SHA256"))
	if e != nil {
		return nil, errMarketplaceConfiguration
	}
	catalog, e := postgres.NewCatalogRelease(pool, cp, price)
	if e != nil {
		return nil, errMarketplaceConfiguration
	}
	key, e := base64.StdEncoding.DecodeString(lookup("MARKETPLACE_HMAC_KEY"))
	if e != nil || len(key) != 32 {
		return nil, errMarketplaceConfiguration
	}
	token := lookup("MERCADOLIBRE_ACCESS_TOKEN")
	if len(token) < 20 || len(token) > 4096 || strings.ContainsAny(token, "\r\n") {
		return nil, errMarketplaceConfiguration
	}
	var ready bool
	e = pool.QueryRow(ctx, `select exists(select 1 from pg_trigger where tgrelid=to_regclass('catalog.marketplace_effect') and tgname='marketplace_effect_immutable' and tgenabled in('O','A'))
 and exists(select 1 from pg_trigger where tgrelid=to_regclass('catalog.marketplace_observation') and tgname='marketplace_observation_immutable' and tgenabled in('O','A'))
 and exists(select 1 from pg_attribute where attrelid=to_regclass('catalog.marketplace_effect') and attname='operation' and not attisdropped)
 and exists(select 1 from pg_constraint where conrelid=to_regclass('catalog.marketplace_effect') and conname='marketplace_effect_operation_check' and pg_get_constraintdef(oid) like '%CONTENT%')
 and exists(select 1 from pg_constraint where conrelid=to_regclass('approval.request') and conname='request_kind_check' and pg_get_constraintdef(oid) like '%marketplace_mutation%')
 and (select count(*) from pg_trigger where not tgisinternal and tgenabled in('O','A') and
 (tgrelid=to_regclass('approval.request') and tgname='bound_request_immutable' or tgrelid=to_regclass('approval.decision') and tgname='bound_decision_immutable'))=2
 and pg_get_functiondef(to_regprocedure('approval.guard_bound_request()')) like '%marketplace_mutation%'
 and pg_get_functiondef(to_regprocedure('approval.guard_bound_decision()')) like '%marketplace_mutation%'`).Scan(&ready)
	if e != nil || !ready {
		return nil, errMarketplaceConfiguration
	}
	service, e := postgres.NewMarketplacePublication(catalog, profile, key, marketplaceToken{lookup}, mb.NewHTTPClient())
	if e != nil {
		return nil, errMarketplaceConfiguration
	}
	return httpapi.MarketplaceModule{Service: service, TenantID: profile.TenantID, OrganizationID: profile.OrganizationID}, nil
}
````

### FILE: `cmd/electromobility-api/marketplace_test.go`

```yaml
block_id: "GO-CONNECTED-MARKETPLACE-MUTATION:file2:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "208f70d935a38c91f949b9c7e5c194e607298a87490953af7ba96e29c108470f"
variables: []
secrets_allowed: false
```

````go
package main

import (
	"context"
	"crypto/sha256"
	cr "elite.local/enterprise/internal/catalogrelease"
	mb "elite.local/enterprise/internal/marketplacebridge"
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

func TestMarketplaceHostGuards(t *testing.T) {
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
	p := mb.Profile{Schema: "elite-marketplace-publication/v1", TenantID: cp.TenantID, OrganizationID: cp.OrganizationID, SellerID: "123456789", SiteID: "MLA", Currency: "ARS", CurrencyDigits: 2, MediaOrigin: cp.Origin, Bindings: []mb.Binding{{VariantID: "fixture-variant", SKU: "fixture-bicycle", CategoryID: "MLA3530", ListingTypeID: "gold_special", StockMode: "item", Attributes: []mb.Attribute{{ID: "ITEM_CONDITION", ValueID: "2230284"}}}}}
	config := map[string]string{"MARKETPLACE_ENABLED": "true", "CATALOG_RELEASE_ENABLED": "true", "MARKETPLACE_HMAC_KEY": base64.StdEncoding.EncodeToString([]byte(strings.Repeat("f", 32))), "MERCADOLIBRE_ACCESS_TOKEN": "synthetic-fixture-token-no-live"}
	for _, v := range []struct {
		prefix string
		value  any
	}{{"MARKETPLACE", p}, {"CATALOG_RELEASE", cp}} {
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
	if module, e := selectedMarketplaceModule(ctx, pool, price, lookup); e != nil || module == nil {
		t.Fatal("current host", e)
	}
	for _, key := range []string{"MARKETPLACE_PROFILE_SHA256", "CATALOG_RELEASE_ENABLED", "MERCADOLIBRE_ACCESS_TOKEN", "MARKETPLACE_HMAC_KEY"} {
		old := config[key]
		config[key] = ""
		if _, e := selectedMarketplaceModule(ctx, pool, price, lookup); e == nil {
			t.Fatal("missing binding admitted", key)
		}
		config[key] = old
	}
	if _, e = pool.Exec(ctx, `alter table catalog.marketplace_observation disable trigger marketplace_observation_immutable`); e != nil {
		t.Fatal(e)
	}
	_, rejected := selectedMarketplaceModule(ctx, pool, price, lookup)
	if _, e = pool.Exec(ctx, `alter table catalog.marketplace_observation enable trigger marketplace_observation_immutable`); e != nil {
		t.Fatal(e)
	}
	if rejected == nil {
		t.Fatal("missing observation guard admitted")
	}
	if _, e = pool.Exec(ctx, `alter table catalog.marketplace_effect disable trigger marketplace_effect_immutable`); e != nil {
		t.Fatal(e)
	}
	defer func() {
		if _, e := pool.Exec(context.Background(), `alter table catalog.marketplace_effect enable trigger marketplace_effect_immutable`); e != nil {
			t.Error(e)
		}
	}()
	if _, e = selectedMarketplaceModule(ctx, pool, price, lookup); e == nil {
		t.Fatal("missing immutable effect guard admitted")
	}
	t.Log("MARKETPLACE_HOST_PASS exact profile/hash/catalog activation/token/HMAC/migration required; no provider network call")
}
````

### FILE: `config/marketplace.reference.json`

```yaml
block_id: "GO-CONNECTED-MARKETPLACE-MUTATION:file3:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "27f93422cfe44a500f30419ead672c0e43a009ea93c237409cc93bb2ba851c56"
variables: []
secrets_allowed: false
```

````json
{
  "schema": "elite-marketplace-publication/v1",
  "tenant_id": "d84c64f9-1252-4db5-bf5a-3af8026190f1",
  "organization_id": "marketplace-store",
  "seller_id": "123456789",
  "site_id": "MLA",
  "currency": "ARS",
  "currency_digits": 2,
  "media_origin": "https://catalog.example.invalid",
  "bindings": [
    {
      "variant_id": "fixture-variant",
      "sku": "fixture-bicycle",
      "category_id": "MLA3530",
      "listing_type_id": "gold_special",
      "attributes": [
        {
          "id": "ITEM_CONDITION",
          "value_id": "2230284"
        }
      ],
      "stock_mode": "item"
    }
  ]
}
````

### FILE: `db/migrations/0079_marketplace_mutation.down.sql`

```yaml
block_id: "GO-CONNECTED-MARKETPLACE-MUTATION:file4:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "95ccd8a66fe26f7ba92455ff8759845616edc7d3bcdd87af7c7fcec3669fb889"
variables: []
secrets_allowed: false
```

````sql
begin;
drop table catalog.marketplace_effect;
-- AUTHORED service/approval/inventory binding; upstream owners remain authoritative.
alter table approval.request drop constraint request_kind_check;
alter table approval.request add constraint request_kind_check check(kind in
 ('reservation','sale','refund','payment','whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review','training_assessment'));
alter table approval.request drop constraint request_payload_binding_check;
alter table approval.request add constraint request_payload_binding_check check
 (kind not in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review','training_assessment') or
  (organization_id is not null and length(organization_id) between 1 and 128 and
   payload is not null and jsonb_typeof(payload)='object' and pg_column_size(payload)<=65536 and amount_minor_units=0));
create or replace function approval.guard_bound_request() returns trigger language plpgsql as $$
begin
 if tg_op='DELETE' then
  if old.kind in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review','training_assessment') then raise exception 'bound approval request is immutable'; end if;
  return old;
 end if;
 if old.kind in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review','training_assessment') or new.kind in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review','training_assessment') then
  if row(new.tenant_id,new.request_id,new.kind,new.subject_id,new.amount_minor_units,new.requester,new.evidence_sha,new.organization_id,new.payload,new.created_at)
   is distinct from row(old.tenant_id,old.request_id,old.kind,old.subject_id,old.amount_minor_units,old.requester,old.evidence_sha,old.organization_id,old.payload,old.created_at)
   or (old.state<>'pending' and row(new.state,new.decided_at) is distinct from row(old.state,old.decided_at))
  then raise exception 'bound approval request is immutable'; end if;
 end if;
 return new;
end $$;
create or replace function approval.guard_bound_decision() returns trigger language plpgsql as $$
begin
 if exists(select 1 from approval.request where tenant_id=old.tenant_id and request_id=old.request_id and kind in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review','training_assessment')) then
  raise exception 'bound approval decision is immutable';
 end if;
 if tg_op='DELETE' then return old;end if;return new;
end $$;

commit;
````

### FILE: `db/migrations/0079_marketplace_mutation.up.sql`

```yaml
block_id: "GO-CONNECTED-MARKETPLACE-MUTATION:file5:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "2cba50e0cf726b06d21075a4064f6ce74a52147502367ab9ab38fbe1a81c2564"
variables: []
secrets_allowed: false
```

````sql
begin;
-- AUTHORED service/approval/inventory binding; upstream owners remain authoritative.
alter table approval.request drop constraint request_kind_check;
alter table approval.request add constraint request_kind_check check(kind in
 ('reservation','sale','refund','payment','whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review','training_assessment','marketplace_mutation'));
alter table approval.request drop constraint request_payload_binding_check;
alter table approval.request add constraint request_payload_binding_check check
 (kind not in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review','training_assessment','marketplace_mutation') or
  (organization_id is not null and length(organization_id) between 1 and 128 and
   payload is not null and jsonb_typeof(payload)='object' and pg_column_size(payload)<=65536 and amount_minor_units=0));
create or replace function approval.guard_bound_request() returns trigger language plpgsql as $$
begin
 if tg_op='DELETE' then
  if old.kind in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review','training_assessment','marketplace_mutation') then raise exception 'bound approval request is immutable'; end if;
  return old;
 end if;
 if old.kind in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review','training_assessment','marketplace_mutation') or new.kind in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review','training_assessment','marketplace_mutation') then
  if row(new.tenant_id,new.request_id,new.kind,new.subject_id,new.amount_minor_units,new.requester,new.evidence_sha,new.organization_id,new.payload,new.created_at)
   is distinct from row(old.tenant_id,old.request_id,old.kind,old.subject_id,old.amount_minor_units,old.requester,old.evidence_sha,old.organization_id,old.payload,old.created_at)
   or (old.state<>'pending' and row(new.state,new.decided_at) is distinct from row(old.state,old.decided_at))
  then raise exception 'bound approval request is immutable'; end if;
 end if;
 return new;
end $$;
create or replace function approval.guard_bound_decision() returns trigger language plpgsql as $$
begin
 if exists(select 1 from approval.request where tenant_id=old.tenant_id and request_id=old.request_id and kind in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review','training_assessment','marketplace_mutation')) then
  raise exception 'bound approval decision is immutable';
 end if;
 if tg_op='DELETE' then return old;end if;return new;
end $$;

create table catalog.marketplace_effect (
 tenant_id uuid not null, channel_code text not null check(channel_code='mercadolibre_catalog'),
 delivery_key text not null, seller_id text not null, sku text not null,
 generation bigint not null, source_sha256 text not null check(source_sha256~'^[0-9a-f]{64}$'),
 profile_sha256 text not null check(profile_sha256~'^[0-9a-f]{64}$'),
 request_sha256 text not null check(request_sha256~'^[0-9a-f]{64}$'),
 created_at timestamptz not null default clock_timestamp(),
 primary key(tenant_id,channel_code,delivery_key),
 foreign key(tenant_id,generation) references catalog.release_publication(tenant_id,generation),
 foreign key(tenant_id,delivery_key) references approval.request(tenant_id,request_id),
 foreign key(tenant_id,channel_code,delivery_key) references communication.outbound_delivery(tenant_id,channel_code,delivery_key)
);
create index marketplace_effect_sku_idx on catalog.marketplace_effect(tenant_id,seller_id,sku);
create trigger marketplace_effect_immutable before update or delete on catalog.marketplace_effect for each row execute function catalog.release_immutable();
commit;
````

### FILE: `docs/MARKETPLACE_MUTATION_REFERENCE.md`

```yaml
block_id: "GO-CONNECTED-MARKETPLACE-MUTATION:file6:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "10313c9dec375752e7e3aae4deead024fff512e9b9df9f2991e1d62951c0bfce"
variables: []
secrets_allowed: false
```

````markdown
# Mercado Libre: approved catalog mutations

This pack supplies AUTHORED HTTP, SQL and host composition around the existing
catalog, inventory ATP, human approval and outbound-delivery owners. It contains
no Mercado Libre SDK source. Exact official documentation snapshots are listed
in docs/marketplace/official-contract.lock.json. The archived Go SDK is not used.

The closed local claim is an existing, seller-owned User Products item:
PRICE, STOCK, PAUSE and RESUME. MEDIA upload and CREATE initial publication
are connected in MARKETPLACE_INITIAL_REFERENCE.md; existing-item CONTENT
updates are connected in MARKETPLACE_CONTENT_REFERENCE.md. Google Merchant
mapping and remaining T2805 orchestration remain tracked in the library gap map. These endpoints do not
silently approximate those operations.

Profile fixes tenant/organization, seller/site/currency and each variant's SKU,
category and inventory location mode. The reference uses explicit synthetic
values. Listing metadata in the intent records the approved source context;
these four operations never change title, attributes or pictures.

The price is the exact current published minor-unit amount, serialized using
the existing decimal representation owner. No floats, currency conversion,
tax calculation or pricing policy is introduced. The stock quantity comes
from the existing inventory.serial_atp function with a captured horizon.
This is its ATP policy, not a new promise that inbound stock is physically
received. Neither provider write nor reconciliation reserves local inventory.
Live channel allocation and provider orders retain their existing owners.

Supported stock modes are the normal item quantity, Full/Flex selling-address
quantity and one explicitly bound seller warehouse. Full inventory is never
written. The adapter reads item ownership, standard prices and stock locations
before preparing a request. Location writes retain the returned x-version.
Ambiguous/multiple standard prices and unexpected location topology fail closed.
Dynamic pricing blocks PRICE before the PUT. The not-yet-available
/prices/standard write contract is not selected.

Workflow through the same Go application and identity verifier:

1. POST /v1/admin/marketplace/prepare with approval_id (16–64 characters),
   generation (decimal string), variant_id, operation, item_id and expires_at.
   This only reads the provider and the current catalog/ATP source.
2. POST the returned exact intent to /v1/admin/marketplace/requests.
3. A distinct reviewer POSTs request_sha256, approve and reason to
   /v1/admin/marketplace/requests/{id}/decision.
4. POST {} to /v1/admin/marketplace/requests/{id}/send.
5. GET /v1/admin/marketplace/requests/{id} after any missing response.
   If unknown, POST {} to /reconcile, which only reads provider state.

Each permission is explicit: marketplace:request, marketplace:approve,
marketplace:send, marketplace:read and marketplace:reconcile. The original
global administrator permission retains its documented organization semantics.
Unprivileged users are rejected before reading the body. JSON is bounded to
32 KiB with a five-second read deadline; duplicate/extra fields fail. No private
response is cacheable. Preparation expires within fifteen minutes.

The shared PostgreSQL fence permits one provider attempt per approved id.
Another unresolved command for the seller/SKU cannot create a second effect.
Lost acknowledgements, expired sender leases and readback mismatch remain
unknown; no automatic PUT retry exists. A documented rejection or a changed
precondition becomes terminal and requires a new reviewed command.
HTTP200 alone does not mean the requested price/stock/state was applied.
Reconciliation confirms observed state, not a distributed transaction or
exclusive ownership of concurrent edits made outside this application.
Catalog changes after claim may require a subsequent approved synchronization;
the receipt always identifies its original generation and source hash.

Activation uses the existing CATALOG_RELEASE_* settings plus:

    MARKETPLACE_ENABLED=false
    MARKETPLACE_PROFILE_FILE=
    MARKETPLACE_PROFILE_SHA256=
    MARKETPLACE_HMAC_KEY=
    MERCADOLIBRE_ACCESS_TOKEN=

When enabled, the application requires the exact profile/hash, catalog
activation, HMAC material and token. Missing migration/immutable approval
guards prevent startup. The HTTP origin is fixed to api.mercadolibre.com,
TLS >=1.2, no redirects or environment proxy, bounded response/timeouts.
Provider token renewal is supplied by the existing activation environment
owner; credentials never enter browser requests, receipts or logs.

The local reference uses a loopback fixture transport only in tests, synthetic
principals and explicit human decisions. Tests cover exact large amounts,
three stock modes, provider/API acknowledgement loss, twelve concurrent sends,
JSONB canonicalization, cross-product/seller binding, ignored price, version
drift, immutable evidence, startup guard and populated/empty migration rollback.
This is library infrastructure evidence, not live provider or production approval.
````

### FILE: `docs/marketplace/official-contract.lock.json`

```yaml
block_id: "GO-CONNECTED-MARKETPLACE-MUTATION:file7:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "0ec5cacd075cff804daec0ade51fc743010eb27d95cf199c38dcacc9a25eecb1"
variables: []
secrets_allowed: false
```

````json
{
  "schema": "elite-marketplace-contract-lock/v1",
  "observed_at": "2026-09-13T03:58:03.128909+00:00",
  "provenance": "AUTHORED HTTP contract mapping; no Mercado Libre source code copied",
  "sources": [
    {
      "id": "publish",
      "url": "https://developers.mercadolibre.com.ar/es_ar/publica-productos",
      "status": 200,
      "bytes": 337519,
      "sha256": "81afd38a0b2129bdbee94a94b7a4c5c70d8d655e0f30064d6c31681959b22237"
    },
    {
      "id": "updates",
      "url": "https://developers.mercadolibre.com.ar/es_ar/producto-sincroniza-modifica-publicaciones",
      "status": 200,
      "bytes": 284334,
      "sha256": "fd65167d83faf26d90e3b716dc59e97be9d8b651e2a8c355836046333d0b83c0"
    },
    {
      "id": "prices",
      "url": "https://developers.mercadolibre.com.ar/es_ar/api-de-precios",
      "status": 200,
      "bytes": 209009,
      "sha256": "a55b9cc9640770e3aaea3c4455e1bdb587fe8cdf22ba5b9af4c60fc535d7c764"
    },
    {
      "id": "stock",
      "url": "https://developers.mercadolibre.com.ar/en_us/distributed-stock",
      "status": 200,
      "bytes": 163126,
      "sha256": "ccc83823912518a919c48f02bb507e5a0eb44408f302d3fdaf51fe36848c9770"
    },
    {
      "id": "search",
      "url": "https://developers.mercadolibre.com.ar/items-y-busquedas",
      "status": 200,
      "bytes": 282492,
      "sha256": "8c469657bc8c01e210f51c132d19ac13e7d662fcc31ec41ad4082233975b4272"
    },
    {
      "id": "user-products",
      "url": "https://developers.mercadolibre.com.ar/es_ar/user-products",
      "status": 200,
      "bytes": 209787,
      "sha256": "465aba2731d7d7e4ac08cf368d881364c72a8aa90877421621666f2d8136d6b6"
    },
    {
      "id": "automation",
      "url": "https://developers.mercadolibre.com.ar/automatizaciones-de-precios",
      "status": 200,
      "bytes": 363840,
      "sha256": "59ae4836ea4eecb5fa632a2a9808a864960d6eeeab30685bbec6187f78381caa",
      "observed_at": "2026-09-13T04:23:23.921123+00:00"
    },
    {
      "id": "pictures",
      "url": "https://developers.mercadolibre.com.ar/es_ar/saldo-de-la-cuenta/trabajar-con-imagenes",
      "sha256": "2e85a485a52fa1311b7e8aaa1391fb18df140dc16b880dd723b5705f43fbe695",
      "bytes": 207091,
      "status": 200,
      "observed_at": "2026-09-13T05:00:47.432677+00:00"
    }
  ],
  "revalidation": "At provider activation or official contract/advisory change. A mismatch fails closed and reopens provider admission."
}
````

### FILE: `internal/marketplacebridge/approval_test.go`

```yaml
block_id: "GO-CONNECTED-MARKETPLACE-MUTATION:file8:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "0372853930a3bf62d3c2689f9bf6aa764503272524c76f1ea6c001c8b0bb99ad"
variables: []
secrets_allowed: false
```

````go
package marketplacebridge_test

import (
	"elite.local/enterprise/internal/approval"
	"strings"
	"testing"
)

func TestMarketplaceApprovalAlwaysManual(t *testing.T) {
	r := approval.NewRegistry(approval.Policy{AutoApproveMinorUnits: 1000})
	state, e := r.Submit(approval.Request{TenantID: "fixture", ID: "marketplace-approval-0001", Kind: approval.KindMarketplaceMutation, SubjectID: "fixture-sku", Requester: "maker", EvidenceSHA: strings.Repeat("a", 64)})
	if e != nil || state != approval.StatePending {
		t.Fatal("provider mutation automatically approved", state, e)
	}
	if _, e = r.Approve("fixture", "marketplace-approval-0001", "maker", "self"); e == nil {
		t.Fatal("separation of duties")
	}
	if state, e = r.Approve("fixture", "marketplace-approval-0001", "reviewer", "exact reviewed payload"); e != nil || state != approval.StateApproved {
		t.Fatal(state, e)
	}
}
````

### FILE: `internal/marketplacebridge/contract.go`

```yaml
block_id: "GO-CONNECTED-MARKETPLACE-MUTATION:file9:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "45fd9aafa4ab471fe995aae9b1c0e4f3de68ddfa10cce2a500af1890a139642f"
variables: []
secrets_allowed: false
```

````go
// AUTHORED mapping of an approved catalog publication to official Mercado Libre
// HTTP contracts. Business amounts, inventory and approval use existing owners.
package marketplacebridge

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/url"
	"regexp"
	"sort"
	"time"

	"elite.local/enterprise/internal/approval"
	cr "elite.local/enterprise/internal/catalogrelease"
	"elite.local/enterprise/internal/channels"
	sv "elite.local/enterprise/internal/storedvaluebridge"
)

var ErrBinding = errors.New("marketplace: invalid or changed binding")
var ErrUnknown = errors.New("marketplace: effect not confirmed; GET reconciliation required")
var digits = regexp.MustCompile(`^[0-9]{3,20}$`)
var itemID = regexp.MustCompile(`^ML[A-Z][0-9]{6,20}$`)
var categoryID = regexp.MustCompile(`^ML[A-Z][0-9]{2,20}$`)
var upID = regexp.MustCompile(`^ML[A-Z]U[0-9]{6,20}$`)
var siteID = regexp.MustCompile(`^ML[A-Z]$`)
var code = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_.:-]{0,63}$`)

type Attribute struct {
	ID        string `json:"id"`
	ValueID   string `json:"value_id,omitempty"`
	ValueName string `json:"value_name,omitempty"`
}
type Binding struct {
	VariantID     string      `json:"variant_id"`
	SKU           string      `json:"sku"`
	CategoryID    string      `json:"category_id"`
	ListingTypeID string      `json:"listing_type_id"`
	Attributes    []Attribute `json:"attributes"`
	StockMode     string      `json:"stock_mode"`
	StoreID       string      `json:"store_id,omitempty"`
	NetworkNodeID string      `json:"network_node_id,omitempty"`
}
type Profile struct {
	Schema         string    `json:"schema"`
	TenantID       string    `json:"tenant_id"`
	OrganizationID string    `json:"organization_id"`
	SellerID       string    `json:"seller_id"`
	SiteID         string    `json:"site_id"`
	Currency       string    `json:"currency"`
	CurrencyDigits int       `json:"currency_digits"`
	MediaOrigin    string    `json:"media_origin"`
	Bindings       []Binding `json:"bindings"`
}

func (p Profile) Validate() error {
	u, e := url.Parse(p.MediaOrigin)
	if p.Schema != "elite-marketplace-publication/v1" || !cr.ValidID(p.TenantID) || !cr.ValidID(p.OrganizationID) || !digits.MatchString(p.SellerID) || !siteID.MatchString(p.SiteID) || !regexp.MustCompile(`^[A-Z]{3}$`).MatchString(p.Currency) || p.CurrencyDigits < 0 || p.CurrencyDigits > 6 || e != nil || u.Scheme != "https" || u.Host == "" || u.User != nil || u.Path != "" || u.RawQuery != "" || u.Fragment != "" || len(p.MediaOrigin) > 512 || len(p.Bindings) < 1 || len(p.Bindings) > 32 {
		return ErrBinding
	}
	variants, skus := map[string]bool{}, map[string]bool{}
	for _, b := range p.Bindings {
		if !cr.ValidID(b.VariantID) || !code.MatchString(b.SKU) || !categoryID.MatchString(b.CategoryID) || b.CategoryID[:3] != p.SiteID || !code.MatchString(b.ListingTypeID) || variants[b.VariantID] || skus[b.SKU] || len(b.Attributes) < 1 || len(b.Attributes) > 64 {
			return ErrBinding
		}
		variants[b.VariantID] = true
		skus[b.SKU] = true
		if b.StockMode != "item" && b.StockMode != "selling_address" && b.StockMode != "seller_warehouse" {
			return ErrBinding
		}
		if b.StockMode == "seller_warehouse" {
			if !digits.MatchString(b.StoreID) || !code.MatchString(b.NetworkNodeID) {
				return ErrBinding
			}
		} else if b.StoreID != "" || b.NetworkNodeID != "" {
			return ErrBinding
		}
		found := false
		seen := map[string]bool{}
		for _, a := range b.Attributes {
			if !code.MatchString(a.ID) || seen[a.ID] || (a.ValueID == "") == (a.ValueName == "") || a.ValueID != "" && !code.MatchString(a.ValueID) || a.ValueName != "" && !cr.ValidText(a.ValueName, 256) || a.ID == "SELLER_SKU" {
				return ErrBinding
			}
			seen[a.ID] = true
			if a.ID == "ITEM_CONDITION" {
				found = true
			}
		}
		if !found {
			return ErrBinding
		}
	}
	raw, e := json.Marshal(p)
	if e != nil {
		return ErrBinding
	}
	if _, _, e = approval.CanonicalPayload(raw); e != nil {
		return ErrBinding
	}
	return nil
}
func (p Profile) SHA256() string { _, h, _ := cr.Canonical(p); return h }
func (p Profile) Binding(variant string) (Binding, error) {
	for _, b := range p.Bindings {
		if b.VariantID == variant {
			return b, nil
		}
	}
	return Binding{}, ErrBinding
}
func Decode(raw []byte, v any) error {
	if len(raw) == 0 || len(raw) > 65536 {
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

type PrepareRequest struct {
	ApprovalID      string    `json:"approval_id"`
	Generation      int64     `json:"generation,string"`
	VariantID       string    `json:"variant_id"`
	Operation       string    `json:"operation"`
	ItemID          string    `json:"item_id,omitempty"`
	MediaApprovalID string    `json:"media_approval_id,omitempty"`
	ExpiresAt       time.Time `json:"expires_at"`
}

func (r PrepareRequest) Validate() error {
	if !cr.ValidID(r.ApprovalID) || len(r.ApprovalID) < 16 || r.Generation < 1 || !cr.ValidID(r.VariantID) || r.ExpiresAt.IsZero() {
		return ErrBinding
	}
	switch r.Operation {
	case "MEDIA":
		if r.ItemID != "" || r.MediaApprovalID != "" {
			return ErrBinding
		}
	case "CONTENT":
		if !itemID.MatchString(r.ItemID) || !cr.ValidID(r.MediaApprovalID) || len(r.MediaApprovalID) < 16 || r.MediaApprovalID == r.ApprovalID {
			return ErrBinding
		}
	case "CREATE":
		if r.ItemID != "" || !cr.ValidID(r.MediaApprovalID) || len(r.MediaApprovalID) < 16 || r.MediaApprovalID == r.ApprovalID {
			return ErrBinding
		}
	case "PRICE", "STOCK", "PAUSE", "RESUME":
		if r.MediaApprovalID != "" {
			return ErrBinding
		}
		if !itemID.MatchString(r.ItemID) {
			return ErrBinding
		}
	default:
		return ErrBinding
	}
	return nil
}

type Intent struct {
	PrepareRequest
	MediaSHA256     string          `json:"media_sha256,omitempty"`
	TenantID        string          `json:"tenant_id"`
	OrganizationID  string          `json:"organization_id"`
	ProfileSHA256   string          `json:"profile_sha256"`
	SourceSHA256    string          `json:"source_sha256"`
	SellerID        string          `json:"seller_id"`
	SKU             string          `json:"sku"`
	BeforeSHA256    string          `json:"before_sha256"`
	Quantity        int64           `json:"quantity"`
	PriceMinorUnits int64           `json:"price_minor_units,string"`
	StockHorizon    time.Time       `json:"stock_horizon"`
	Body            json.RawMessage `json:"body"`
	Method          string          `json:"method"`
	Path            string          `json:"path"`
	StockVersion    string          `json:"stock_version,omitempty"`
	Expected        Item            `json:"expected"`
}

func (r Intent) Message() channels.Message {
	_, h, _ := cr.Canonical(r)
	return channels.Message{ChannelCode: "mercadolibre_catalog", TenantID: r.TenantID, ExternalID: r.SellerID, ThreadID: r.SKU, Direction: channels.DirectionOut, DeliveryKey: r.ApprovalID, Text: h}
}

type Picture struct {
	Source    string `json:"source,omitempty"`
	ID        string `json:"id,omitempty"`
	SecureURL string `json:"secure_url,omitempty"`
}
type Item struct {
	SoldQuantity  *int64      `json:"sold_quantity,omitempty"`
	ID            string      `json:"id"`
	SellerID      int64       `json:"seller_id"`
	SiteID        string      `json:"site_id"`
	CategoryID    string      `json:"category_id"`
	ListingTypeID string      `json:"listing_type_id"`
	BuyingMode    string      `json:"buying_mode"`
	FamilyName    string      `json:"family_name"`
	UserProductID string      `json:"user_product_id"`
	SKU           string      `json:"seller_custom_field"`
	Currency      string      `json:"currency_id"`
	Quantity      int64       `json:"available_quantity"`
	Status        string      `json:"status"`
	Attributes    []Attribute `json:"attributes"`
	Pictures      []Picture   `json:"pictures"`
	Tags          []string    `json:"tags"`
}
type Observation struct {
	Item         Item   `json:"item"`
	PriceMinor   string `json:"price_minor"`
	Quantity     int64  `json:"quantity"`
	StockVersion string `json:"stock_version"`
	Exists       bool   `json:"exists"`
}

func (o Observation) SHA256() string { _, h, _ := cr.Canonical(o); return h }
func AttributesEqual(want, got []Attribute) bool {
	for _, a := range want {
		found := false
		for _, b := range got {
			if a.ID == b.ID && (a.ValueID != "" && a.ValueID == b.ValueID || a.ValueName != "" && a.ValueName == b.ValueName) {
				if found {
					return false
				}
				found = true
			}
		}
		if !found {
			return false
		}
	}
	return true
}
func Build(p Profile, pub cr.PublicDocument, r PrepareRequest, quantity int64, horizon time.Time, before Observation, contentPicture ...string) (Intent, error) {
	if r.Operation == "CONTENT" && (len(contentPicture) != 1 || !ValidPictureID(contentPicture[0])) || r.Operation != "CONTENT" && len(contentPicture) != 0 {
		return Intent{}, ErrBinding
	}
	if p.Validate() != nil || r.Validate() != nil || pub.Generation != r.Generation || pub.Currency != p.Currency || !cr.ValidSHA(pub.SourceSHA256) || quantity < 0 || horizon.IsZero() {
		return Intent{}, ErrBinding
	}
	b, e := p.Binding(r.VariantID)
	if e != nil {
		return Intent{}, e
	}
	var variant *cr.Variant
	var model *cr.Model
	for i := range pub.Variants {
		if pub.Variants[i].ID == r.VariantID {
			if variant != nil {
				return Intent{}, ErrBinding
			}
			variant = &pub.Variants[i]
		}
	}
	if variant == nil || variant.AmountMinorUnits <= 0 {
		return Intent{}, ErrBinding
	}
	for i := range pub.Models {
		if pub.Models[i].ID == variant.ModelID {
			if model != nil {
				return Intent{}, ErrBinding
			}
			model = &pub.Models[i]
		}
	}
	if model == nil || !cr.ValidSHA(model.Media.SHA256) || !cr.ValidText(model.DisplayName, 256) {
		return Intent{}, ErrBinding
	}
	attrs := append([]Attribute(nil), b.Attributes...)
	sort.Slice(attrs, func(i, j int) bool { return attrs[i].ID < attrs[j].ID })
	expected := Item{ID: r.ItemID, SiteID: p.SiteID, CategoryID: b.CategoryID, ListingTypeID: b.ListingTypeID, BuyingMode: "buy_it_now", FamilyName: model.DisplayName, SKU: b.SKU, Currency: p.Currency, Quantity: quantity, Attributes: attrs, Pictures: []Picture{{Source: p.MediaOrigin + "/api/public/catalog/media/" + model.Media.SHA256}}}
	expected.UserProductID = before.Item.UserProductID
	out := Intent{PrepareRequest: r, TenantID: p.TenantID, OrganizationID: p.OrganizationID, ProfileSHA256: p.SHA256(), SourceSHA256: pub.SourceSHA256, SellerID: p.SellerID, SKU: b.SKU, BeforeSHA256: before.SHA256(), Quantity: quantity, StockHorizon: horizon, Expected: expected, Method: "PUT", Path: "/items/" + r.ItemID}
	major := json.Number(sv.Major(variant.AmountMinorUnits, p.CurrencyDigits))
	var body any
	switch r.Operation {
	case "CONTENT":
		if before.Item.SoldQuantity == nil || *before.Item.SoldQuantity != 0 {
			return Intent{}, ErrBinding
		}
		out.MediaSHA256 = model.Media.SHA256
		out.Expected.Pictures = []Picture{{ID: contentPicture[0]}}
		body = contentBody(out.Expected)
	case "MEDIA", "CREATE":
		out.MediaSHA256 = model.Media.SHA256
		out.Method = "POST"
		if before.Exists {
			return Intent{}, ErrBinding
		}
		if r.Operation == "MEDIA" {
			out.Path = "/pictures/items/upload"
			body = map[string]string{"catalog_png_sha256": out.MediaSHA256}
		} else {
			if b.StockMode != "item" || len(before.Item.Pictures) != 1 || !ValidPictureID(before.Item.Pictures[0].ID) {
				return Intent{}, ErrBinding
			}
			out.Path = "/items"
			out.Expected.Pictures = []Picture{{ID: before.Item.Pictures[0].ID}}
			body = creationBody(out.Expected, major)
		}
	case "PRICE":
		body = map[string]any{"price": major}
	case "STOCK":
		switch b.StockMode {
		case "item":
			body = map[string]any{"available_quantity": quantity}
		case "selling_address":
			out.Path = "/user-products/" + before.Item.UserProductID + "/stock/type/selling_address"
			body = map[string]any{"quantity": quantity}
		case "seller_warehouse":
			out.Path = "/user-products/" + before.Item.UserProductID + "/stock/type/seller_warehouse"
			body = map[string]any{"locations": []map[string]any{{"store_id": b.StoreID, "network_node_id": b.NetworkNodeID, "quantity": quantity}}}
		}
		if b.StockMode != "item" {
			if !upID.MatchString(before.Item.UserProductID) || before.StockVersion == "" {
				return Intent{}, ErrBinding
			}
			out.StockVersion = before.StockVersion
		}
	case "PAUSE":
		body = map[string]string{"status": "paused"}
	case "RESUME":
		body = map[string]string{"status": "active"}
	}
	if r.Operation != "MEDIA" && r.Operation != "CREATE" && (!before.Exists || before.Item.ID != r.ItemID) {
		return Intent{}, ErrBinding
	}
	out.PriceMinorUnits = variant.AmountMinorUnits
	out.Body, _, e = cr.Canonical(body)
	return out, e
}
````

### FILE: `internal/marketplacebridge/fuzz_test.go`

```yaml
block_id: "GO-CONNECTED-MARKETPLACE-MUTATION:file10:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "3ea58e21c051fc5d9d672228679be5eb204d9459c2cdbabce28eaaf3916976cc"
variables: []
secrets_allowed: false
```

````go
package marketplacebridge

import (
	"bytes"
	"encoding/json"
	"testing"
)

func FuzzMarketplaceCommandBoundary(f *testing.F) {
	f.Add([]byte(`{"approval_id":"marketplace-content-01","media_approval_id":"marketplace-media-0001","generation":"1","variant_id":"bike-standard","operation":"CONTENT","item_id":"MLA123456789","expires_at":"2026-09-13T00:00:00Z"}`))
	f.Add([]byte(`{"approval_id":"marketplace-media-0001","generation":"1","variant_id":"bike-standard","operation":"MEDIA","expires_at":"2026-09-13T00:00:00Z"}`))
	f.Add([]byte(`{"approval_id":"marketplace-create-001","media_approval_id":"marketplace-media-0001","generation":"1","variant_id":"bike-standard","operation":"CREATE","expires_at":"2026-09-13T00:00:00Z"}`))
	f.Add([]byte(`{"approval_id":"marketplace-price-0001","generation":"1","variant_id":"bike-standard","operation":"PRICE","item_id":"MLA123456789","expires_at":"2026-09-13T00:00:00Z"}`))
	f.Add([]byte(`{"approval_id":"one","Approval_ID":"two","operation":"STOCK"}`))
	f.Add([]byte(`{"operation":"PRICE","item_id":"../../users"}`))
	f.Fuzz(func(t *testing.T, raw []byte) {
		var r PrepareRequest
		if Decode(raw, &r) != nil {
			return
		}
		before, _ := json.Marshal(r)
		err := r.Validate()
		after, _ := json.Marshal(r)
		if !bytes.Equal(before, after) {
			t.Fatal("validation mutated command")
		}
		var round PrepareRequest
		if Decode(after, &round) != nil {
			t.Fatal("accepted typed serialization cannot decode")
		}
		if err == nil {
			if len(r.ApprovalID) < 16 || r.Generation < 1 || r.ExpiresAt.IsZero() {
				t.Fatal("command boundary lost")
			}
			switch r.Operation {
			case "CONTENT":
				if !itemID.MatchString(r.ItemID) || len(r.MediaApprovalID) < 16 || r.MediaApprovalID == r.ApprovalID {
					t.Fatal("unbound content")
				}
			case "MEDIA":
				if r.ItemID != "" || r.MediaApprovalID != "" {
					t.Fatal("media target injection")
				}
			case "CREATE":
				if r.ItemID != "" || len(r.MediaApprovalID) < 16 || r.MediaApprovalID == r.ApprovalID {
					t.Fatal("unbound create")
				}
			case "PRICE", "STOCK", "PAUSE", "RESUME":
				if !itemID.MatchString(r.ItemID) || r.MediaApprovalID != "" {
					t.Fatal("mutation target injection")
				}
			default:
				t.Fatal("unadmitted write")
			}
		}
	})
}
````

### FILE: `internal/marketplacebridge/http.go`

```yaml
block_id: "GO-CONNECTED-MARKETPLACE-MUTATION:file11:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "f0ad3742eb047af2e995efdc19a17a003aaf4172f18c8cb3790791b6f4e86edf"
variables: []
secrets_allowed: false
```

````go
package marketplacebridge

// AUTHORED bounded HTTP/serialization glue for official provider contracts.
// A 2xx response alone does not confirm the requested effect.
import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"elite.local/enterprise/internal/approval"
	"elite.local/enterprise/internal/outbounddelivery"
	sv "elite.local/enterprise/internal/storedvaluebridge"
)

const APIBaseURL = "https://api.mercadolibre.com"

type Doer interface {
	Do(*http.Request) (*http.Response, error)
}
type TokenSource interface {
	MercadoLibreAccessToken(context.Context) (string, error)
}
type Client struct {
	profile Profile
	tokens  TokenSource
	http    Doer
}

func NewClient(p Profile, t TokenSource, h Doer) (*Client, error) {
	if p.Validate() != nil || t == nil || h == nil {
		return nil, ErrBinding
	}
	raw, _ := json.Marshal(p)
	var cloned Profile
	if json.Unmarshal(raw, &cloned) != nil {
		return nil, ErrBinding
	}
	return &Client{cloned, t, h}, nil
}
func NewHTTPClient() *http.Client {
	return &http.Client{Timeout: 8 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
		Transport: &http.Transport{Proxy: nil, DialContext: (&net.Dialer{Timeout: 3 * time.Second, KeepAlive: 30 * time.Second}).DialContext, TLSClientConfig: &tls.Config{MinVersion: tls.VersionTLS12}, TLSHandshakeTimeout: 3 * time.Second, ResponseHeaderTimeout: 5 * time.Second, MaxResponseHeaderBytes: 32768}}
}
func (c *Client) call(ctx context.Context, method, path string, body []byte, version string) (int, []byte, http.Header, error) {
	if !strings.HasPrefix(path, "/") || strings.HasPrefix(path, "//") || strings.ContainsAny(path, "\r\n\\") || len(path) > 512 || len(body) > 32768 {
		return 0, nil, nil, ErrBinding
	}
	token, e := c.tokens.MercadoLibreAccessToken(ctx)
	if e != nil || len(token) < 20 || len(token) > 4096 || strings.ContainsAny(token, "\r\n") {
		return 0, nil, nil, ErrBinding
	}
	request, e := http.NewRequestWithContext(ctx, method, APIBaseURL+path, bytes.NewReader(body))
	if e != nil {
		return 0, nil, nil, ErrBinding
	}
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("Accept", "application/json")
	if len(body) > 0 {
		request.Header.Set("Content-Type", "application/json")
	}
	if version != "" {
		request.Header.Set("x-version", version)
	}
	response, e := c.http.Do(request)
	if e != nil {
		return 0, nil, nil, ErrUnknown
	}
	defer response.Body.Close()
	raw, e := io.ReadAll(io.LimitReader(response.Body, 32769))
	if e != nil || len(raw) > 32768 {
		return response.StatusCode, nil, nil, ErrUnknown
	}
	if len(raw) > 0 {
		if _, _, e = approval.CanonicalPayload(raw); e != nil {
			return response.StatusCode, nil, nil, ErrUnknown
		}
	}
	return response.StatusCode, raw, response.Header.Clone(), nil
}
func hasTag(tags []string, want string) bool {
	for _, s := range tags {
		if s == want {
			return true
		}
	}
	return false
}
func (c *Client) Observe(ctx context.Context, variant, id string) (Observation, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	var out Observation
	b, e := c.profile.Binding(variant)
	if e != nil || !itemID.MatchString(id) || id[:3] != c.profile.SiteID {
		return out, ErrBinding
	}
	status, raw, _, e := c.call(ctx, "GET", "/items/"+id, nil, "")
	if e != nil || status != 200 || json.Unmarshal(raw, &out.Item) != nil {
		return out, ErrUnknown
	}
	v := out.Item
	if v.ID != id || strconv.FormatInt(v.SellerID, 10) != c.profile.SellerID || v.SiteID != c.profile.SiteID || v.CategoryID != b.CategoryID || v.SKU != b.SKU || !upID.MatchString(v.UserProductID) || v.FamilyName == "" || v.Currency != c.profile.Currency {
		return out, ErrBinding
	}
	status, raw, _, e = c.call(ctx, "GET", "/items/"+id+"/prices", nil, "")
	var prices struct {
		ID     string `json:"id"`
		Prices []struct {
			Type       string      `json:"type"`
			Amount     json.Number `json:"amount"`
			Currency   string      `json:"currency_id"`
			Conditions struct {
				Context []string `json:"context_restrictions"`
				Start   *string  `json:"start_time"`
				End     *string  `json:"end_time"`
				Min     *int64   `json:"min_purchase_unit"`
			} `json:"conditions"`
		} `json:"prices"`
	}
	if e != nil || status != 200 || json.Unmarshal(raw, &prices) != nil || prices.ID != id {
		return out, ErrUnknown
	}
	matched := 0
	for _, p := range prices.Prices {
		if p.Type != "standard" || p.Conditions.Start != nil || p.Conditions.End != nil || p.Conditions.Min != nil {
			continue
		}
		if len(p.Conditions.Context) != 0 && !(len(p.Conditions.Context) == 1 && p.Conditions.Context[0] == "channel_marketplace") {
			continue
		}
		if p.Currency != c.profile.Currency {
			return out, ErrBinding
		}
		n, e := sv.Minor(p.Amount.String(), c.profile.CurrencyDigits)
		if e != nil || n <= 0 {
			return out, ErrBinding
		}
		out.PriceMinor = strconv.FormatInt(n, 10)
		matched++
	}
	if matched != 1 {
		return out, ErrBinding
	}
	status, raw, headers, e := c.call(ctx, "GET", "/user-products/"+v.UserProductID+"/stock", nil, "")
	var stock struct {
		ID        string `json:"id"`
		UserID    int64  `json:"user_id"`
		Locations []struct {
			Type     string `json:"type"`
			StoreID  string `json:"store_id"`
			Node     string `json:"network_node_id"`
			Quantity int64  `json:"quantity"`
		} `json:"locations"`
	}
	if e != nil || status != 200 || json.Unmarshal(raw, &stock) != nil || stock.ID != v.UserProductID || strconv.FormatInt(stock.UserID, 10) != c.profile.SellerID {
		return out, ErrUnknown
	}
	matched = 0
	for _, loc := range stock.Locations {
		if loc.Quantity < 0 {
			return out, ErrBinding
		}
		switch b.StockMode {
		case "item":
			// Legacy available_quantity write is selected only for this normal,
			// non-multi-origin, non-Full inventory configuration.
			if len(stock.Locations) != 1 || loc.Type != "selling_address" || loc.StoreID != "" || loc.Node != "" {
				return out, ErrBinding
			}
			out.Quantity = v.Quantity
			matched++
		case "selling_address":
			if loc.Type == "seller_warehouse" {
				return out, ErrBinding
			}
			if loc.Type == "selling_address" {
				out.Quantity = loc.Quantity
				matched++
			}
		case "seller_warehouse":
			if loc.Type == "selling_address" {
				return out, ErrBinding
			}
			if loc.Type == "seller_warehouse" && loc.StoreID == b.StoreID && loc.Node == b.NetworkNodeID {
				out.Quantity = loc.Quantity
				matched++
			}
		}
	}
	if matched != 1 || out.Quantity < 0 {
		return out, ErrBinding
	}
	if b.StockMode != "item" {
		version := headers.Get("x-version")
		n, e := strconv.ParseInt(version, 10, 64)
		if e != nil || n < 0 || strconv.FormatInt(n, 10) != version {
			return out, ErrBinding
		}
		out.StockVersion = version
	}
	out.Exists = true
	return out, nil
}
func (c *Client) validate(r Intent) error {
	if r.PrepareRequest.Validate() != nil || r.TenantID != c.profile.TenantID || r.OrganizationID != c.profile.OrganizationID || r.ProfileSHA256 != c.profile.SHA256() || r.SellerID != c.profile.SellerID || r.Method != "PUT" || r.PriceMinorUnits <= 0 || r.Quantity < 0 {
		return ErrBinding
	}
	b, e := c.profile.Binding(r.VariantID)
	if e != nil || r.SKU != b.SKU {
		return ErrBinding
	}
	path := "/items/" + r.ItemID
	version := ""
	var body any
	switch r.Operation {
	case "CONTENT":
		if validateContent(c.profile, r) != nil {
			return ErrBinding
		}
		body = contentBody(r.Expected)
	case "PRICE":
		body = map[string]any{"price": json.Number(sv.Major(r.PriceMinorUnits, c.profile.CurrencyDigits))}
	case "STOCK":
		if b.StockMode == "item" {
			body = map[string]any{"available_quantity": r.Quantity}
		} else {
			if !upID.MatchString(r.Expected.UserProductID) {
				return ErrBinding
			}
			path = "/user-products/" + r.Expected.UserProductID + "/stock/type/" + b.StockMode
			version = r.StockVersion
			n, e := strconv.ParseInt(version, 10, 64)
			if e != nil || n < 0 || strconv.FormatInt(n, 10) != version {
				return ErrBinding
			}
			body = map[string]any{"quantity": r.Quantity}
			if b.StockMode == "seller_warehouse" {
				body = map[string]any{"locations": []map[string]any{{"store_id": b.StoreID, "network_node_id": b.NetworkNodeID, "quantity": r.Quantity}}}
			}
		}
	case "PAUSE":
		body = map[string]string{"status": "paused"}
	case "RESUME":
		body = map[string]string{"status": "active"}
	}
	raw, _ := json.Marshal(body)
	canonical, _, e := approval.CanonicalPayload(raw)
	actual, _, parseErr := approval.CanonicalPayload(r.Body)
	if e != nil || parseErr != nil || r.Path != path || r.StockVersion != version || !bytes.Equal(canonical, actual) {
		return ErrBinding
	}
	return nil
}
func (c *Client) confirms(r Intent, o Observation) bool {
	if !o.Exists || o.Item.ID != r.ItemID || o.Item.UserProductID != r.Expected.UserProductID {
		return false
	}
	switch r.Operation {
	case "CONTENT":
		return o.Item.FamilyName == r.Expected.FamilyName && AttributesEqual(r.Expected.Attributes, o.Item.Attributes) && len(o.Item.Pictures) == 1 && o.Item.Pictures[0].ID == r.Expected.Pictures[0].ID
	case "PRICE":
		return o.PriceMinor == strconv.FormatInt(r.PriceMinorUnits, 10)
	case "STOCK":
		return o.Quantity == r.Quantity
	case "PAUSE":
		return o.Item.Status == "paused"
	case "RESUME":
		return o.Item.Status == "active"
	}
	return false
}
func (c *Client) Reconcile(ctx context.Context, r Intent) (outbounddelivery.Receipt, error) {
	if r.Operation == "CREATE" {
		return c.ReconcileCreation(ctx, r, "")
	}
	if c.validate(r) != nil {
		return outbounddelivery.Receipt{}, ErrBinding
	}
	o, e := c.Observe(ctx, r.VariantID, r.ItemID)
	if e != nil || !c.confirms(r, o) {
		return outbounddelivery.Receipt{}, ErrUnknown
	}
	return outbounddelivery.Receipt{ProviderMessageID: r.ItemID, EvidenceSHA256: o.SHA256(), AcceptedAt: time.Now().UTC()}, nil
}
func (c *Client) Write(ctx context.Context, r Intent) (outbounddelivery.Receipt, error) {
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	if c.validate(r) != nil || !r.ExpiresAt.After(time.Now()) {
		return outbounddelivery.Receipt{}, ErrBinding
	}
	o, e := c.Observe(ctx, r.VariantID, r.ItemID)
	if e != nil {
		return outbounddelivery.Receipt{}, e
	}
	if o.SHA256() != r.BeforeSHA256 || o.Item.UserProductID != r.Expected.UserProductID || r.Operation == "PRICE" && hasTag(o.Item.Tags, "dynamic_standard_price") {
		return outbounddelivery.Receipt{}, outbounddelivery.NewTerminalFailure("MERCADOLIBRE_PRECONDITION_CHANGED", o.SHA256(), nil)
	}
	if r.Operation == "CONTENT" {
		if e := c.SingleUnsoldItem(ctx, o.Item); e != nil {
			return outbounddelivery.Receipt{}, outbounddelivery.NewTerminalFailure("MERCADOLIBRE_CONTENT_SCOPE_CHANGED", o.SHA256(), nil)
		}
	}
	wire, _, e := approval.CanonicalPayload(r.Body)
	if e != nil {
		return outbounddelivery.Receipt{}, ErrBinding
	}
	status, raw, _, e := c.call(ctx, r.Method, r.Path, wire, r.StockVersion)
	if e != nil {
		return outbounddelivery.Receipt{}, ErrUnknown
	}
	if status == 400 || status == 401 || status == 403 || status == 409 || status == 422 {
		_, h, e := approval.CanonicalPayload(raw)
		if e != nil {
			return outbounddelivery.Receipt{}, ErrUnknown
		}
		return outbounddelivery.Receipt{}, outbounddelivery.NewTerminalFailure("MERCADOLIBRE_MUTATION_REJECTED", h, nil)
	}
	if status != 200 && status != 204 {
		return outbounddelivery.Receipt{}, ErrUnknown
	}
	receipt, e := c.Reconcile(ctx, r)
	if e != nil {
		return receipt, ErrUnknown
	}
	// Preserve the write and GET evidence together, without raw tokens or PII.
	_, h, e := approval.CanonicalPayload(mustJSON(map[string]any{"write_status": status, "write_body": json.RawMessage(nonempty(raw)), "observation_sha256": receipt.EvidenceSHA256, "request_body": r.Body}))
	if e != nil {
		return outbounddelivery.Receipt{}, ErrUnknown
	}
	receipt.EvidenceSHA256 = h
	return receipt, nil
}
func nonempty(b []byte) []byte {
	if len(b) == 0 {
		return []byte("{}")
	}
	return b
}
func mustJSON(v any) []byte {
	b, e := json.Marshal(v)
	if e != nil {
		panic(errors.New("bounded typed serialization failed"))
	}
	return b
}
````

### FILE: `internal/marketplacebridge/http_test.go`

```yaml
block_id: "GO-CONNECTED-MARKETPLACE-MUTATION:file12:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "d218e8205971a9ac577631aa3d5d6a8a9a2b3da374495f884094799669adc2a7"
variables: []
secrets_allowed: false
```

````go
package marketplacebridge_test

import (
	"bytes"
	"context"
	cr "elite.local/enterprise/internal/catalogrelease"
	"elite.local/enterprise/internal/electromobility"
	mb "elite.local/enterprise/internal/marketplacebridge"
	fixture "elite.local/enterprise/internal/marketplacebridge/testfixture"
	"elite.local/enterprise/internal/outbounddelivery"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestMarketplaceJSONBWhitespaceAndCrossProductBinding(t *testing.T) {
	p := profile("seller_warehouse")
	receiver := fixture.New(p)
	defer receiver.Close()
	client, _ := mb.NewClient(p, tokens{}, receiver.HTTP())
	intent := prepare(t, p, client, "PRICE")
	var spaced bytes.Buffer
	if e := json.Indent(&spaced, intent.Body, "", "  "); e != nil {
		t.Fatal(e)
	}
	intent.Body = spaced.Bytes()
	if _, e := client.Write(context.Background(), intent); e != nil {
		t.Fatal("semantic JSONB representation changed wire authorization", e)
	}
	if receiver.Writes[0].Body != `{"price":90071992547409.91}` {
		t.Fatal("wire is not canonical", receiver.Writes[0].Body)
	}
	intent = prepare(t, p, client, "STOCK")
	original := intent.Expected.UserProductID
	intent.Expected.UserProductID = "MLAU987654321"
	intent.Path = strings.ReplaceAll(intent.Path, original, intent.Expected.UserProductID)
	var terminal *outbounddelivery.TerminalFailure
	if _, e := client.Write(context.Background(), intent); !errors.As(e, &terminal) || len(receiver.Writes) != 1 {
		t.Fatal("observed item does not own write UP", e)
	}
}

type tokens struct{}

func (tokens) MercadoLibreAccessToken(context.Context) (string, error) { return fixture.Token, nil }
func profile(mode string) mb.Profile {
	b := mb.Binding{VariantID: "fixture-variant", SKU: "fixture-bicycle", CategoryID: "MLA3530", ListingTypeID: "gold_special", StockMode: mode, Attributes: []mb.Attribute{{ID: "ITEM_CONDITION", ValueID: "2230284"}}}
	if mode == "seller_warehouse" {
		b.StoreID = "123456"
		b.NetworkNodeID = "ARP12345"
	}
	return mb.Profile{Schema: "elite-marketplace-publication/v1", TenantID: "ca6c9cb7-e548-4f19-a2e6-27d2aa2d4893", OrganizationID: "j3-store", SellerID: "123456789", SiteID: "MLA", Currency: "ARS", CurrencyDigits: 2, MediaOrigin: "https://catalog.example.invalid", Bindings: []mb.Binding{b}}
}
func pub() cr.PublicDocument {
	return cr.PublicDocument{Generation: 1, SourceSHA256: strings.Repeat("a", 64), Currency: "ARS", Models: []cr.Model{{Model: electromobility.Model{ID: "model", DisplayName: "Fixture bicycle"}, Media: cr.Media{SHA256: strings.Repeat("b", 64)}}}, Variants: []cr.Variant{{ID: "fixture-variant", ModelID: "model", AmountMinorUnits: 9007199254740991}}}
}
func prepare(t *testing.T, p mb.Profile, c *mb.Client, op string) mb.Intent {
	t.Helper()
	before, e := c.Observe(context.Background(), p.Bindings[0].VariantID, "MLA123456789")
	if e != nil {
		t.Fatal(e)
	}
	in := mb.PrepareRequest{ApprovalID: "marketplace-unit-0001", Generation: 1, VariantID: p.Bindings[0].VariantID, Operation: op, ItemID: "MLA123456789", ExpiresAt: time.Now().Add(time.Minute)}
	out, e := mb.Build(p, pub(), in, 2, time.Now(), before)
	if e != nil {
		t.Fatal(e)
	}
	return out
}
func TestMarketplaceMutationContracts(t *testing.T) {
	for _, mode := range []string{"item", "selling_address", "seller_warehouse"} {
		t.Run(mode, func(t *testing.T) {
			p := profile(mode)
			receiver := fixture.New(p)
			defer receiver.Close()
			client, e := mb.NewClient(p, tokens{}, receiver.HTTP())
			if e != nil {
				t.Fatal(e)
			}
			for _, op := range []string{"PRICE", "STOCK", "PAUSE", "RESUME"} {
				intent := prepare(t, p, client, op)
				receipt, e := client.Write(context.Background(), intent)
				if e != nil || receipt.Validate() != nil {
					t.Fatal(op, e)
				}
			}
			if len(receiver.Writes) != 4 || receiver.Effects != 4 || receiver.PriceMinor != 9007199254740991 || receiver.Quantity != 2 || receiver.Item.Status != "active" {
				t.Fatal("request/effect mismatch", receiver.Writes)
			}
			if !strings.Contains(receiver.Writes[0].Body, "90071992547409.91") {
				t.Fatal("minor units lost precision")
			}
			if mode != "item" && receiver.Writes[1].Version != "7" {
				t.Fatal("stock version not preserved")
			}
			t.Log("one exact field per write; standard price exact decimal; configured stock location; pause/resume read back")
		})
	}
}
func TestMarketplaceFalseSuccessAndUnknownAreReadOnly(t *testing.T) {
	p := profile("seller_warehouse")
	receiver := fixture.New(p)
	defer receiver.Close()
	client, _ := mb.NewClient(p, tokens{}, receiver.HTTP())
	intent := prepare(t, p, client, "PRICE")
	receiver.IgnorePrice = true
	if _, e := client.Write(context.Background(), intent); !errors.Is(e, mb.ErrUnknown) {
		t.Fatal("200 with ignored price accepted", e)
	}
	if _, e := client.Reconcile(context.Background(), intent); !errors.Is(e, mb.ErrUnknown) {
		t.Fatal("wrong price reconciled", e)
	}
	receiver.IgnorePrice = false
	receiver.PriceMinor = intent.PriceMinorUnits
	if _, e := client.Reconcile(context.Background(), intent); e != nil {
		t.Fatal(e)
	}
	if len(receiver.Writes) != 1 {
		t.Fatal("recovery wrote")
	}
	intent = prepare(t, p, client, "STOCK")
	receiver.DropNext = true
	if _, e := client.Write(context.Background(), intent); !errors.Is(e, mb.ErrUnknown) {
		t.Fatal("lost response accepted", e)
	}
	if _, e := client.Reconcile(context.Background(), intent); e != nil {
		t.Fatal(e)
	}
	if len(receiver.Writes) != 2 || receiver.Quantity != 2 {
		t.Fatal("lost acknowledgement resent")
	}
}
func TestMarketplaceDriftApprovalBindingAndTerminal(t *testing.T) {
	p := profile("item")
	receiver := fixture.New(p)
	defer receiver.Close()
	client, _ := mb.NewClient(p, tokens{}, receiver.HTTP())
	intent := prepare(t, p, client, "PRICE")
	receiver.Item.Tags = []string{"dynamic_standard_price"}
	var terminal *outbounddelivery.TerminalFailure
	if _, e := client.Write(context.Background(), intent); !errors.As(e, &terminal) || len(receiver.Writes) != 0 {
		t.Fatal("automation guard", e)
	}
	receiver.Item.Tags = nil
	intent = prepare(t, p, client, "PRICE")
	tampered := intent
	tampered.Path = "/users/123456789"
	if _, e := client.Write(context.Background(), tampered); !errors.Is(e, mb.ErrBinding) || len(receiver.Writes) != 0 {
		t.Fatal("path injection")
	}
	tampered = intent
	tampered.Body = json.RawMessage(`{"price":1}`)
	if _, e := client.Write(context.Background(), tampered); !errors.Is(e, mb.ErrBinding) {
		t.Fatal("body injection")
	}
	receiver.RejectNext = true
	if _, e := client.Write(context.Background(), intent); !errors.As(e, &terminal) {
		t.Fatal("documented rejection", e)
	}
	receiver.Item.SellerID = 987654321
	if _, e := client.Observe(context.Background(), p.Bindings[0].VariantID, "MLA123456789"); !errors.Is(e, mb.ErrBinding) {
		t.Fatal("foreign seller")
	}
}
````

### FILE: `internal/marketplacebridge/testfixture/receiver.go`

```yaml
block_id: "GO-CONNECTED-MARKETPLACE-MUTATION:file13:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "8e6965b6114e033fcfa04cf9c7ec9c79ea9fa1508e9f736f5f1f6084f1cd93f7"
variables: []
secrets_allowed: false
```

````go
// Package testfixture is AUTHORED synthetic provider infrastructure. It is not a
// provider emulator certification; requests use the fixed documented paths.
package testfixture

import (
	mb "elite.local/enterprise/internal/marketplacebridge"
	sv "elite.local/enterprise/internal/storedvaluebridge"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"sync"
)

const Token = "synthetic-mercadolibre-fixture-token"

type Write struct {
	Path    string
	Body    string
	Version string
}
type Receiver struct {
	Mu                                              sync.Mutex
	Server                                          *httptest.Server
	Profile                                         mb.Profile
	Item                                            mb.Item
	Quantity, PriceMinor, Version                   int64
	DropNext, IgnorePrice, RejectNext, ConflictNext bool
	Writes                                          []Write
	Effects, Reads                                  int
}
type Transport struct {
	Target *url.URL
	Base   http.RoundTripper
}

func (t Transport) RoundTrip(r *http.Request) (*http.Response, error) {
	if r.URL.Scheme != "https" || r.URL.Host != "api.mercadolibre.com" || r.Header.Get("Authorization") != "Bearer "+Token {
		return nil, mb.ErrBinding
	}
	copy := r.Clone(r.Context())
	u := *r.URL
	u.Scheme = t.Target.Scheme
	u.Host = t.Target.Host
	copy.URL = &u
	copy.Host = t.Target.Host
	return t.Base.RoundTrip(copy)
}
func New(p mb.Profile) *Receiver {
	b := p.Bindings[0]
	seller, _ := strconv.ParseInt(p.SellerID, 10, 64)
	s := &Receiver{Profile: p, PriceMinor: 9007199254740993, Quantity: 4, Version: 7}
	s.Item = mb.Item{ID: "MLA123456789", SellerID: seller, SiteID: p.SiteID, CategoryID: b.CategoryID, ListingTypeID: b.ListingTypeID, BuyingMode: "buy_it_now", FamilyName: "Fixture bicycle", UserProductID: "MLAU123456789", SKU: b.SKU, Currency: p.Currency, Quantity: s.Quantity, Status: "active", Attributes: b.Attributes}
	s.Server = httptest.NewServer(http.HandlerFunc(s.serve))
	return s
}
func (s *Receiver) HTTP() *http.Client {
	u, _ := url.Parse(s.Server.URL)
	client := mb.NewHTTPClient()
	client.Transport = Transport{u, client.Transport}
	return client
}
func (s *Receiver) Close() { s.Server.Close() }
func (s *Receiver) serve(w http.ResponseWriter, r *http.Request) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	if r.Header.Get("Authorization") != "Bearer "+Token {
		w.WriteHeader(401)
		io.WriteString(w, `{"error":"unauthorized"}`)
		return
	}
	b := s.Profile.Bindings[0]
	itemPath := "/items/" + s.Item.ID
	stockPath := "/user-products/" + s.Item.UserProductID + "/stock"
	if r.Method == "GET" {
		s.Reads++
		switch r.URL.Path {
		case itemPath:
			_ = json.NewEncoder(w).Encode(s.Item)
		case itemPath + "/prices":
			_ = json.NewEncoder(w).Encode(map[string]any{"id": s.Item.ID, "prices": []any{map[string]any{"type": "standard", "amount": json.Number(sv.Major(s.PriceMinor, s.Profile.CurrencyDigits)), "currency_id": s.Profile.Currency, "conditions": map[string]any{"context_restrictions": []string{}}}}})
		case stockPath:
			location := map[string]any{"type": "selling_address", "quantity": s.Quantity}
			if b.StockMode == "seller_warehouse" {
				location = map[string]any{"type": b.StockMode, "quantity": s.Quantity, "store_id": b.StoreID, "network_node_id": b.NetworkNodeID}
			}
			locations := []any{location}
			if b.StockMode != "item" {
				locations = append(locations, map[string]any{"type": "meli_facility", "quantity": 99})
			}
			w.Header().Set("x-version", strconv.FormatInt(s.Version, 10))
			_ = json.NewEncoder(w).Encode(map[string]any{"id": s.Item.UserProductID, "user_id": s.Item.SellerID, "locations": locations})
		default:
			w.WriteHeader(404)
			io.WriteString(w, `{"error":"not_found"}`)
		}
		return
	}
	if r.Method != "PUT" {
		w.WriteHeader(405)
		io.WriteString(w, `{"error":"method_not_allowed"}`)
		return
	}
	raw, e := io.ReadAll(io.LimitReader(r.Body, 32769))
	if e != nil || len(raw) > 32768 {
		w.WriteHeader(413)
		return
	}
	var input map[string]json.RawMessage
	if json.Unmarshal(raw, &input) != nil {
		w.WriteHeader(400)
		return
	}
	s.Writes = append(s.Writes, Write{r.URL.Path, string(raw), r.Header.Get("x-version")})
	if s.RejectNext {
		s.RejectNext = false
		w.WriteHeader(400)
		io.WriteString(w, `{"error":"validation_error"}`)
		return
	}
	if s.ConflictNext {
		s.ConflictNext = false
		w.WriteHeader(409)
		io.WriteString(w, `{"error":"version_mismatch"}`)
		return
	}
	if r.URL.Path == itemPath {
		if len(input) != 1 {
			w.WriteHeader(400)
			io.WriteString(w, `{"error":"unexpected_fields"}`)
			return
		}
		if x, ok := input["price"]; ok {
			var n json.Number
			if json.Unmarshal(x, &n) != nil {
				w.WriteHeader(400)
				return
			}
			value, e := sv.Minor(n.String(), s.Profile.CurrencyDigits)
			if e != nil {
				w.WriteHeader(400)
				return
			}
			if !s.IgnorePrice {
				s.PriceMinor = value
				s.Effects++
			}
		} else if x, ok := input["available_quantity"]; ok {
			if b.StockMode != "item" || json.Unmarshal(x, &s.Quantity) != nil {
				w.WriteHeader(400)
				return
			}
			s.Item.Quantity = s.Quantity
			s.Effects++
		} else if x, ok := input["status"]; ok {
			var status string
			if json.Unmarshal(x, &status) != nil || status != "paused" && status != "active" {
				w.WriteHeader(400)
				return
			}
			s.Item.Status = status
			s.Effects++
		} else {
			w.WriteHeader(400)
			return
		}
	} else if r.URL.Path == stockPath+"/type/"+b.StockMode && b.StockMode != "item" {
		if r.Header.Get("x-version") != strconv.FormatInt(s.Version, 10) {
			w.WriteHeader(409)
			io.WriteString(w, `{"error":"version_mismatch"}`)
			return
		}
		if b.StockMode == "selling_address" {
			if len(input) != 1 || json.Unmarshal(input["quantity"], &s.Quantity) != nil {
				w.WriteHeader(400)
				return
			}
		} else {
			var locations []struct {
				Store    string `json:"store_id"`
				Node     string `json:"network_node_id"`
				Quantity int64  `json:"quantity"`
			}
			if len(input) != 1 || json.Unmarshal(input["locations"], &locations) != nil || len(locations) != 1 || locations[0].Store != b.StoreID || locations[0].Node != b.NetworkNodeID {
				w.WriteHeader(400)
				return
			}
			s.Quantity = locations[0].Quantity
		}
		s.Version++
		s.Effects++
	} else {
		w.WriteHeader(404)
		io.WriteString(w, `{"error":"not_found"}`)
		return
	}
	if s.DropNext {
		s.DropNext = false
		conn, _, e := w.(http.Hijacker).Hijack()
		if e == nil {
			conn.Close()
		}
		return
	}
	if strings.HasPrefix(r.URL.Path, stockPath) {
		w.WriteHeader(204)
	} else {
		io.WriteString(w, `{"accepted":true}`)
	}
}
````

### FILE: `internal/platform/httpapi/marketplace.go`

```yaml
block_id: "GO-CONNECTED-MARKETPLACE-MUTATION:file14:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "c3ad5e3a719f3f7d38e701b7d832638ef33ac7965a898de125defeb0e139344f"
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
	mb "elite.local/enterprise/internal/marketplacebridge"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/postgres"
)

type MarketplacePublicationService interface {
	Prepare(context.Context, identity.Principal, mb.PrepareRequest) (mb.Intent, error)
	Submit(context.Context, identity.Principal, mb.Intent) (string, bool, error)
	Decide(context.Context, identity.Principal, string, string, bool, string) (approval.State, error)
	Status(context.Context, identity.Principal, string) (postgres.MarketplaceStatus, error)
	Send(context.Context, identity.Principal, string) error
	Reconcile(context.Context, identity.Principal, string) error
}
type MarketplaceModule struct {
	Service                  MarketplacePublicationService
	TenantID, OrganizationID string
}

func (m MarketplaceModule) Register(mux *http.ServeMux, verifier identity.Verifier) {
	if m.Service == nil {
		return
	}
	route := func(method, path, permission string, action func(http.ResponseWriter, *http.Request, identity.Principal)) {
		mux.HandleFunc(method+" /v1/admin/marketplace"+path, func(w http.ResponseWriter, r *http.Request) {
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
		writeProblem(w, 409, "MARKETPLACE_UNCONFIRMED", "consult immutable request and provider reconciliation")
		return true
	}
	route("POST", "/prepare", "marketplace:request", func(w http.ResponseWriter, r *http.Request, p identity.Principal) {
		var in mb.PrepareRequest
		if !decode(w, r, &in) {
			return
		}
		out, e := m.Service.Prepare(r.Context(), p, in)
		if !failed(w, e) {
			writeJSON(w, 200, out)
		}
	})
	route("POST", "/requests", "marketplace:request", func(w http.ResponseWriter, r *http.Request, p identity.Principal) {
		var in mb.Intent
		if !decode(w, r, &in) {
			return
		}
		hash, replay, e := m.Service.Submit(r.Context(), p, in)
		if !failed(w, e) {
			writeJSON(w, 200, map[string]any{"approval_id": in.ApprovalID, "request_sha256": hash, "replay": replay})
		}
	})
	route("POST", "/requests/{id}/decision", "marketplace:approve", func(w http.ResponseWriter, r *http.Request, p identity.Principal) {
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
	route("GET", "/requests/{id}", "marketplace:read", func(w http.ResponseWriter, r *http.Request, p identity.Principal) {
		out, e := m.Service.Status(r.Context(), p, r.PathValue("id"))
		if !failed(w, e) {
			writeJSON(w, 200, out)
		}
	})
	for _, v := range []struct {
		name, permission string
		call             func(context.Context, identity.Principal, string) error
	}{{"send", "marketplace:send", m.Service.Send}, {"reconcile", "marketplace:reconcile", m.Service.Reconcile}} {
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

### FILE: `internal/platform/postgres/marketplace_publication.go`

```yaml
block_id: "GO-CONNECTED-MARKETPLACE-MUTATION:file15:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "075188764397aa0acdb3bd21a291136de5418d92870d951d3fbf59db0fd613bf"
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
	mb "elite.local/enterprise/internal/marketplacebridge"
	"elite.local/enterprise/internal/outbounddelivery"
	"elite.local/enterprise/internal/platform/identity"
	"github.com/jackc/pgx/v5"
)

type MarketplacePublication struct {
	catalog   *CatalogRelease
	profile   mb.Profile
	approvals *HumanApprovals
	fence     *OutboundDeliveryStore
	client    *mb.Client
}

func NewMarketplacePublication(catalog *CatalogRelease, p mb.Profile, key []byte, tokens mb.TokenSource, http mb.Doer) (*MarketplacePublication, error) {
	if catalog == nil || p.Validate() != nil || p.TenantID != catalog.profile.TenantID || p.OrganizationID != catalog.profile.OrganizationID || p.Currency != catalog.profile.Currency {
		return nil, mb.ErrBinding
	}
	raw, _ := json.Marshal(p)
	var frozen mb.Profile
	if json.Unmarshal(raw, &frozen) != nil {
		return nil, mb.ErrBinding
	}
	client, e := mb.NewClient(frozen, tokens, http)
	if e != nil {
		return nil, e
	}
	fence, e := NewOutboundDeliveryStore(catalog.pool, key, 30*time.Second)
	if e != nil {
		return nil, e
	}
	return &MarketplacePublication{catalog, frozen, NewHumanApprovals(catalog.pool), fence, client}, nil
}
func (s *MarketplacePublication) authorized(p identity.Principal, permission string) bool {
	return s != nil && approvalPrincipal(p, s.profile.TenantID, s.profile.OrganizationID, permission)
}
func (s *MarketplacePublication) source(ctx context.Context, tx pgx.Tx, generation int64, variant string, horizon time.Time) (cr.PublicDocument, int64, error) {
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
func (s *MarketplacePublication) Prepare(ctx context.Context, p identity.Principal, r mb.PrepareRequest) (mb.Intent, error) {
	if !s.authorized(p, "marketplace:request") || r.Validate() != nil {
		return mb.Intent{}, mb.ErrBinding
	}
	// Remote preflight is read-only. Source publication and ATP remain database owners.
	before, e := s.initialObservation(ctx, r)
	if r.Operation != "MEDIA" && r.Operation != "CREATE" {
		before, e = s.client.Observe(ctx, r.VariantID, r.ItemID)
	}
	if e != nil {
		return mb.Intent{}, e
	}
	if r.Operation == "CONTENT" {
		if e = s.client.SingleUnsoldItem(ctx, before.Item); e != nil {
			return mb.Intent{}, e
		}
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
	doc, quantity, e := s.source(ctx, tx, r.Generation, r.VariantID, now)
	if e != nil {
		return mb.Intent{}, e
	}
	var pictures []string
	if r.Operation == "CONTENT" {
		pic, e := s.acceptedMedia(ctx, tx, r)
		if e != nil {
			return mb.Intent{}, e
		}
		pictures = []string{pic}
	}
	out, e := mb.Build(s.profile, doc, r, quantity, now, before, pictures...)
	if e != nil {
		return out, e
	}
	return out, tx.Commit(ctx)
}
func (s *MarketplacePublication) guard(ctx context.Context, tx pgx.Tx, r mb.Intent, requireCurrent bool) error {
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
	// Rebuild all source-derived fields; use the recorded observation only as a
	// before-state binding. The live client separately compares that observation.
	before := mb.Observation{Item: r.Expected, Exists: true, StockVersion: r.StockVersion}
	if r.Operation == "MEDIA" {
		before = mb.Observation{}
	}
	if r.Operation == "CREATE" {
		pic, e := s.acceptedMedia(ctx, tx, r.PrepareRequest)
		if e != nil {
			return e
		}
		before = mb.Observation{Item: mb.Item{Pictures: []mb.Picture{{ID: pic}}}}
	}
	var pictures []string
	if r.Operation == "CONTENT" {
		pic, e := s.acceptedMedia(ctx, tx, r.PrepareRequest)
		if e != nil {
			return e
		}
		pictures = []string{pic}
		zero := int64(0)
		before.Item.SoldQuantity = &zero
	}
	rebuilt, e := mb.Build(s.profile, doc, r.PrepareRequest, q, r.StockHorizon, before, pictures...)
	if e != nil {
		return e
	}
	rebuilt.BeforeSHA256 = r.BeforeSHA256
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
func (s *MarketplacePublication) Submit(ctx context.Context, p identity.Principal, r mb.Intent) (string, bool, error) {
	if !s.authorized(p, "marketplace:request") || !cr.ValidSHA(r.BeforeSHA256) {
		return "", false, mb.ErrBinding
	}
	raw, hash, e := cr.Canonical(r)
	if e != nil {
		return "", false, e
	}
	spec := HumanApprovalSpec{Request: approval.Request{TenantID: p.TenantID, ID: r.ApprovalID, Kind: approval.KindMarketplaceMutation, SubjectID: r.SKU, Requester: p.Subject, EvidenceSHA: hash}, OrganizationID: s.profile.OrganizationID, Payload: raw}
	replay, e := s.approvals.Submit(ctx, p, spec, "marketplace:request", func(ctx context.Context, tx pgx.Tx) error { return s.guard(ctx, tx, r, true) })
	return hash, replay, e
}
func (s *MarketplacePublication) request(ctx context.Context, id string) (mb.Intent, string, error) {
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
func (s *MarketplacePublication) Decide(ctx context.Context, p identity.Principal, id, hash string, approve bool, reason string) (approval.State, error) {
	if !s.authorized(p, "marketplace:approve") {
		return "", mb.ErrBinding
	}
	r, stored, e := s.request(ctx, id)
	if e != nil || stored != hash {
		return "", mb.ErrBinding
	}
	return s.approvals.Decide(ctx, p, s.profile.TenantID, id, s.profile.OrganizationID, hash, approve, reason, "marketplace:approve", func(ctx context.Context, tx pgx.Tx) error {
		if !approve {
			return nil
		}
		return s.guard(ctx, tx, r, true)
	})
}

type MarketplaceStatus struct {
	ApprovalID     string          `json:"approval_id"`
	ApprovalState  string          `json:"approval_state"`
	RequestSHA256  string          `json:"request_sha256"`
	Payload        json.RawMessage `json:"payload"`
	DeliveryState  string          `json:"delivery_state"`
	FailureCode    string          `json:"failure_code"`
	EvidenceSHA256 string          `json:"evidence_sha256"`
}

func (s *MarketplacePublication) Status(ctx context.Context, p identity.Principal, id string) (MarketplaceStatus, error) {
	var out MarketplaceStatus
	if !s.authorized(p, "marketplace:read") || !cr.ValidID(id) {
		return out, mb.ErrBinding
	}
	e := s.catalog.pool.QueryRow(ctx, `select a.request_id,a.state,a.evidence_sha,a.payload,coalesce(d.state,''),coalesce(d.failure_code,''),coalesce(d.evidence_sha256_hex,'')
 from approval.request a left join communication.outbound_delivery d on d.tenant_id=a.tenant_id and d.channel_code='mercadolibre_catalog' and d.delivery_key=a.request_id
 where a.tenant_id=$1 and a.request_id=$2 and a.organization_id=$3 and a.kind='marketplace_mutation' and a.payload->>'profile_sha256'=$4`, s.profile.TenantID, id, s.profile.OrganizationID, s.profile.SHA256()).Scan(&out.ApprovalID, &out.ApprovalState, &out.RequestSHA256, &out.Payload, &out.DeliveryState, &out.FailureCode, &out.EvidenceSHA256)
	return out, e
}

type marketplaceBound struct {
	service *MarketplacePublication
	intent  mb.Intent
	hash    string
}

func (b *marketplaceBound) Claim(ctx context.Context, m channels.Message, h string) (outbounddelivery.Claim, error) {
	return b.service.fence.claimWithAdmission(ctx, m, h, func(ctx context.Context, tx pgx.Tx) error {
		r := b.intent
		s := b.service
		// One active effect per tenant/provider SKU, including another approval id.
		if _, e := tx.Exec(ctx, `select pg_advisory_xact_lock(hashtextextended($1,0))`, "marketplace:"+r.TenantID+":"+r.SellerID+":"+r.SKU); e != nil {
			return e
		}
		var busy bool
		e := tx.QueryRow(ctx, `select exists(select 1 from catalog.marketplace_effect i join communication.outbound_delivery d using(tenant_id,channel_code,delivery_key)
    where i.tenant_id=$1 and i.seller_id=$2 and i.sku=$3 and (d.state in ('sending','unknown') or (d.state='accepted' and i.operation='CREATE' and $4='CREATE')) and i.operation<>'MEDIA' and $4<>'MEDIA')`, r.TenantID, r.SellerID, r.SKU, r.Operation).Scan(&busy)
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
		_, e = tx.Exec(ctx, `insert into catalog.marketplace_effect(tenant_id,channel_code,delivery_key,seller_id,sku,generation,source_sha256,profile_sha256,request_sha256,operation) values($1,'mercadolibre_catalog',$2,$3,$4,$5,$6,$7,$8,$9)`, r.TenantID, r.ApprovalID, r.SellerID, r.SKU, r.Generation, r.SourceSHA256, r.ProfileSHA256, b.hash, r.Operation)
		return e
	})
}
func (b *marketplaceBound) Complete(c context.Context, m channels.Message, h string, r outbounddelivery.Receipt) error {
	return b.service.fence.Complete(c, m, h, r)
}
func (b *marketplaceBound) MarkUnknown(c context.Context, m channels.Message, h, code string) error {
	return b.service.fence.MarkUnknown(c, m, h, code)
}
func (b *marketplaceBound) MarkFailed(c context.Context, m channels.Message, h, sha, code string) error {
	return b.service.fence.MarkFailed(c, m, h, sha, code)
}

type marketplaceSender struct {
	client  *mb.Client
	intent  mb.Intent
	service *MarketplacePublication
}

func (s marketplaceSender) SendWithReceipt(ctx context.Context, m channels.Message) (outbounddelivery.Receipt, error) {
	want, _ := outbounddelivery.MessageSHA256(s.intent.Message())
	got, e := outbounddelivery.MessageSHA256(m)
	if e != nil || got != want {
		return outbounddelivery.Receipt{}, mb.ErrBinding
	}
	if s.intent.Operation == "MEDIA" || s.intent.Operation == "CREATE" {
		return s.service.sendInitial(ctx, s.intent)
	}
	return s.client.Write(ctx, s.intent)
}

type marketplaceOutbound struct{}

func (marketplaceOutbound) Code() string { return "mercadolibre_catalog" }
func (marketplaceOutbound) Receive(context.Context) ([]channels.Message, error) {
	return nil, mb.ErrBinding
}
func (marketplaceOutbound) Send(context.Context, channels.Message) error { return mb.ErrBinding }
func (s *MarketplacePublication) Send(ctx context.Context, p identity.Principal, id string) error {
	if !s.authorized(p, "marketplace:send") {
		return mb.ErrBinding
	}
	r, h, e := s.request(ctx, id)
	if e != nil {
		return e
	}
	channel := outbounddelivery.Channel{CodeValue: "mercadolibre_catalog", Receiver: marketplaceOutbound{}, Sender: marketplaceSender{s.client, r, s}, Store: &marketplaceBound{s, r, h}}
	return channel.Send(ctx, r.Message())
}
func (s *MarketplacePublication) Reconcile(ctx context.Context, p identity.Principal, id string) error {
	if !s.authorized(p, "marketplace:reconcile") {
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
	if status.DeliveryState == "accepted" {
		return nil
	}
	if status.DeliveryState != "unknown" && status.DeliveryState != "sending" {
		return mb.ErrBinding
	}
	message := r.Message()
	hash, e := outbounddelivery.MessageSHA256(message)
	if e != nil {
		return e
	}
	// Existing lease expiry makes a crashed sender unknown. No new claim is made.
	if _, e = s.fence.Claim(ctx, message, hash); !errors.Is(e, outbounddelivery.ErrUnknown) {
		return mb.ErrUnknown
	}
	receipt, e := s.reconcileInitial(ctx, r)
	if r.Operation != "MEDIA" && r.Operation != "CREATE" {
		receipt, e = s.client.Reconcile(ctx, r)
	}
	if e != nil {
		return e
	}
	return s.fence.ReconcileAccepted(ctx, message, hash, receipt)
}
````

### FILE: `internal/platform/postgres/marketplace_publication_integration_test.go`

```yaml
block_id: "GO-CONNECTED-MARKETPLACE-MUTATION:file16:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "c089a9b49f5c183ec78f606f8319705bc6d4bd6d089dd1023bb6ac19af496182"
variables: []
secrets_allowed: false
```

````go
package postgres_test

// AUTHORED local HTTP/PG/provider fixture. No live account or real inventory.
import (
	"bytes"
	"context"
	cr "elite.local/enterprise/internal/catalogrelease"
	"elite.local/enterprise/internal/electromobility"
	mb "elite.local/enterprise/internal/marketplacebridge"
	fixture "elite.local/enterprise/internal/marketplacebridge/testfixture"
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
	"strings"
	"sync"
	"testing"
	"time"
)

type marketplaceFixtureToken struct{}

func (marketplaceFixtureToken) MercadoLibreAccessToken(context.Context) (string, error) {
	return fixture.Token, nil
}

type marketplaceFixtureVerifier map[string]identity.Principal
type marketplaceRecorder struct {
	*httptest.ResponseRecorder
	target http.ResponseWriter
}

func (w marketplaceRecorder) Unwrap() http.ResponseWriter { return w.target }

func (v marketplaceFixtureVerifier) Verify(_ context.Context, token string) (identity.Principal, error) {
	p, ok := v[token]
	if !ok {
		return p, identity.ErrUnauthenticated
	}
	return p, nil
}

func TestMarketplacePublicationConnected(t *testing.T) {
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
		return identity.Principal{TenantID: tenant, Subject: subject, Permissions: map[string]struct{}{"catalog:draft": {}, "catalog:read": {}, "catalog:publish": {}, "catalog:review:legal": {}, "catalog:review:technical": {}, "catalog:review:media": {}, "catalog:review:publication": {}, "marketplace:request": {}, "marketplace:approve": {}, "marketplace:read": {}, "marketplace:send": {}, "marketplace:reconcile": {}}, Organizations: map[string]struct{}{org: {}}}
	}
	maker, reviewer := principal("marketplace-maker"), principal("marketplace-reviewer")
	cp := cr.Profile{Schema: "elite-catalog-publication/v1", TenantID: tenant, OrganizationID: org, Market: "AR", Currency: "ARS", Origin: "https://catalog.example.invalid", PolicyCode: cr.PolicyCode}
	catalog, e := db.NewCatalogRelease(pool, cp, db.NewCommerce(pool))
	if e != nil {
		t.Fatal(e)
	}
	source, e := catalog.CreateSource(ctx, maker, cr.SourceRequest{CommandID: "marketplace-source", Model: electromobility.Model{Code: "fixture-bicycle", DisplayName: "Fixture bicycle", VehicleClass: "bicycle", Specification: json.RawMessage("{}")}, Variants: []cr.SourceVariant{{Code: "fixture-bicycle-standard", DisplayName: "Standard", BatterySpecification: json.RawMessage("{}"), AmountMinorUnits: 9007199254740991, TaxMode: "not-applicable"}}, ValidFrom: time.Now().Add(-time.Hour)})
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
	profile := mb.Profile{Schema: "elite-marketplace-publication/v1", TenantID: tenant, OrganizationID: org, SellerID: "123456789", SiteID: "MLA", Currency: "ARS", CurrencyDigits: 2, MediaOrigin: cp.Origin, Bindings: []mb.Binding{{VariantID: source.SourceVariantIDs[0], SKU: "fixture-bicycle-standard", CategoryID: "MLA3530", ListingTypeID: "gold_special", StockMode: "seller_warehouse", StoreID: "123456", NetworkNodeID: "ARP12345", Attributes: []mb.Attribute{{ID: "ITEM_CONDITION", ValueID: "2230284"}}}}}
	receiver := fixture.New(profile)
	defer receiver.Close()
	service, e := db.NewMarketplacePublication(catalog, profile, bytes.Repeat([]byte("f"), 32), marketplaceFixtureToken{}, receiver.HTTP())
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
	httpapi.MarketplaceModule{Service: service, TenantID: tenant, OrganizationID: org}.Register(mux, verifier)
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
		request, _ := http.NewRequest(method, api.URL+"/v1/admin/marketplace"+path, bytes.NewReader(raw))
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
	prepare := func(op, id string) mb.Intent {
		t.Helper()
		var r mb.Intent
		_, e := call("maker", "POST", "/prepare", mb.PrepareRequest{ApprovalID: id, Generation: published.Generation, VariantID: profile.Bindings[0].VariantID, Operation: op, ItemID: receiver.Item.ID, ExpiresAt: time.Now().Add(5 * time.Minute)}, &r)
		if e != nil {
			t.Fatal("prepare", e)
		}
		if r.SourceSHA256 != published.SnapshotSHA256 || r.PriceMinorUnits != 9007199254740991 || r.Quantity != 0 {
			t.Fatal("source/ATP binding", r)
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
		var result db.MarketplaceStatus
		if _, e := call("maker", "GET", "/requests/"+id, nil, &result); e != nil || result.DeliveryState != want {
			t.Fatal("durable state", result.ApprovalID, result.ApprovalState, result.DeliveryState, result.FailureCode, "provider writes", len(receiver.Writes), e)
		}
	}
	price := prepare("PRICE", "marketplace-price-0001")
	for _, token := range []string{"forbidden", "wrong-org", "wrong-tenant"} {
		if code, e := call(token, "POST", "/requests", price, nil); e == nil || code != 403 {
			t.Fatal("scope boundary", token)
		}
	}
	tampered := price
	tampered.PriceMinorUnits = 1
	if _, e := call("maker", "POST", "/requests", tampered, nil); e == nil {
		t.Fatal("source amount replaced")
	}
	hash := submit(price)
	if _, e := call("maker", "POST", "/requests/"+price.ApprovalID+"/send", struct{}{}, nil); e == nil {
		t.Fatal("unapproved effect")
	}
	approve(price, hash)
	var wg sync.WaitGroup
	successes := 0
	var mu sync.Mutex
	for i := 0; i < 12; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, e := call("maker", "POST", "/requests/"+price.ApprovalID+"/send", struct{}{}, nil)
			if e == nil {
				mu.Lock()
				successes++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	status(price.ApprovalID, "accepted")
	if successes == 0 || len(receiver.Writes) != 1 || receiver.PriceMinor != 9007199254740991 {
		t.Fatal("concurrent send", successes, len(receiver.Writes))
	}
	stock := prepare("STOCK", "marketplace-stock-0001")
	approve(stock, submit(stock))
	receiver.DropNext = true
	if _, e := call("maker", "POST", "/requests/"+stock.ApprovalID+"/send", struct{}{}, nil); e == nil {
		t.Fatal("lost provider response accepted")
	}
	status(stock.ApprovalID, "unknown")
	if _, e := call("maker", "POST", "/requests/"+stock.ApprovalID+"/send", struct{}{}, nil); e == nil {
		t.Fatal("unknown resent")
	}
	if _, e := call("maker", "POST", "/requests/"+stock.ApprovalID+"/reconcile", struct{}{}, nil); e != nil {
		t.Fatal("GET recovery", e)
	}
	status(stock.ApprovalID, "accepted")
	if len(receiver.Writes) != 2 || receiver.Quantity != 0 || receiver.Writes[1].Version != "7" {
		t.Fatal("stock recovery changed write/version")
	}
	pause := prepare("PAUSE", "marketplace-pause-0001")
	approve(pause, submit(pause))
	apiMu.Lock()
	dropAPI = true
	apiMu.Unlock()
	if _, e := call("maker", "POST", "/requests/"+pause.ApprovalID+"/send", struct{}{}, nil); e == nil {
		t.Fatal("API loss fixture")
	}
	status(pause.ApprovalID, "accepted")
	if _, e := call("maker", "POST", "/requests/"+pause.ApprovalID+"/send", struct{}{}, nil); e != nil {
		t.Fatal("accepted replay", e)
	}
	resume := prepare("RESUME", "marketplace-resume-0001")
	approve(resume, submit(resume))
	if _, e := call("maker", "POST", "/requests/"+resume.ApprovalID+"/send", struct{}{}, nil); e != nil {
		t.Fatal("resume", e)
	}
	rejection := prepare("PRICE", "marketplace-reject-0001")
	approve(rejection, submit(rejection))
	receiver.RejectNext = true
	if _, e := call("maker", "POST", "/requests/"+rejection.ApprovalID+"/send", struct{}{}, nil); e == nil {
		t.Fatal("rejection accepted")
	}
	status(rejection.ApprovalID, "failed_terminal")
	if _, e := call("maker", "POST", "/requests/"+rejection.ApprovalID+"/send", struct{}{}, nil); e == nil {
		t.Fatal("terminal retry")
	}
	drift := prepare("STOCK", "marketplace-drift-0001")
	approve(drift, submit(drift))
	receiver.Version++
	if _, e := call("maker", "POST", "/requests/"+drift.ApprovalID+"/send", struct{}{}, nil); e == nil {
		t.Fatal("remote version drift accepted")
	}
	status(drift.ApprovalID, "failed_terminal")
	var approvals, effects, attempts int
	if e = pool.QueryRow(ctx, `select count(*) from approval.request where tenant_id=$1 and kind='marketplace_mutation'`, tenant).Scan(&approvals); e != nil {
		t.Fatal(e)
	}
	if e = pool.QueryRow(ctx, `select count(*),coalesce(sum(attempt_count),0) from communication.outbound_delivery where tenant_id=$1 and channel_code='mercadolibre_catalog'`, tenant).Scan(&effects, &attempts); e != nil {
		t.Fatal(e)
	}
	if approvals != 6 || effects != 6 || attempts != 6 || len(receiver.Writes) != 5 || receiver.Effects != 4 {
		t.Fatal("durable counts", approvals, effects, attempts, len(receiver.Writes), receiver.Effects)
	}
	for _, q := range []string{
		`update approval.request set payload=payload||'{"quantity":999}'::jsonb where tenant_id=$1 and kind='marketplace_mutation'`,
		`delete from approval.decision where tenant_id=$1 and request_id='marketplace-price-0001'`,
		`delete from catalog.marketplace_effect where tenant_id=$1`,
	} {
		if _, e = pool.Exec(ctx, q, tenant); e == nil {
			t.Fatal("evidence mutable")
		}
	}
	t.Logf("MARKETPLACE_CONNECTED_PASS source publication=%s; exact price=9007199254740991; existing ATP=0; 6 bound approvals/fences; 5 provider PUT attempts/4 effects; 12 concurrent sends -> 1 PUT; provider/API loss recovered; rejection/drift terminal; API POSTs=%d", published.SnapshotSHA256, apiPosts)
}
````

## 6. Configuration surface

docs/MARKETPLACE_MUTATION_REFERENCE.md and config/marketplace.reference.json. Fixed official API origin; future credentials supplied only at activation.

## 7. Dependency bill

AUTHORED inevitable composition/HTTP/serialization/proof. No Mercado Libre SDK source copied. Existing Go/PG/BC/Odoo-derived owners retain their own narrow provenance and notices.

## 8. Apply order

Requires connected catalog, inventory, stored-value decimal utility, human approval, outbound fence and application owners. Includes finite local proof; no hosted CI or live entitlement required.

## 9. Verification

Six approvals/fences,12 concurrent sends->1PUT,5provider attempts/4effects, provider/API loss recovery, JSONB and cross-UP binding,3stockmodes,exact9007199254740991amount,ignoredprice/automation/rejection/versiondrift,host+rollback,2s fuzz,vet/build.

## 10. Reconstruction evidence

MARKETPLACE_MUTATION_RELEASE_V402.md/json. TEST02/03/07 and overall library READY remain open pending the ordered controls.


V402 composed delta: T2805 approved immutable catalog PNG to separately approved MEDIA upload and CREATE User Products item, existing mutations, same host/API/PG, durable observation/fence and bounded GET-only creation recovery. Local fixtures proven; remaining content/feed/comms owners stay open. AUTHORED glue, no provider-source attribution.

### FILE: `db/migrations/0080_marketplace_initial.down.sql`

```yaml
block_id: "GO-CONNECTED-MARKETPLACE-MUTATION-INITIAL-DELTA:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "3be818cc502284bdd5d91748497da44074319dca45c08caa5a28d81a86a81db8"
variables: []
secrets_allowed: false
```

````sql
begin;
do $$ begin
 if exists(select 1 from catalog.marketplace_observation) or exists(select 1 from catalog.marketplace_effect where operation in ('MEDIA','CREATE')) then
  raise exception 'cannot remove marketplace initial publication evidence';
 end if;
end $$;
drop table catalog.marketplace_observation;
alter table catalog.marketplace_effect drop column operation;
commit;
````

### FILE: `db/migrations/0080_marketplace_initial.up.sql`

```yaml
block_id: "GO-CONNECTED-MARKETPLACE-MUTATION-INITIAL-DELTA:file2:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "2e219f33af0151045f0b7c53ee6eb1f607806f32e6404f9485d00d594a975969"
variables: []
secrets_allowed: false
```

````sql
-- AUTHORED persistence binding. Existing catalog/approval/fence remain owners.
begin;
alter table catalog.marketplace_effect add column operation text not null default 'LEGACY_MUTATION'
 check(operation in ('LEGACY_MUTATION','PRICE','STOCK','PAUSE','RESUME','MEDIA','CREATE'));
create table catalog.marketplace_observation (
 tenant_id uuid not null, approval_id text not null, channel_code text not null default 'mercadolibre_catalog' check(channel_code='mercadolibre_catalog'),
 operation text not null check(operation in ('MEDIA','CREATE')),
 provider_id text not null, media_sha256 text not null check(media_sha256~'^[0-9a-f]{64}$'),
 observation jsonb not null check(jsonb_typeof(observation)='object'),
 observation_sha256 text not null check(observation_sha256~'^[0-9a-f]{64}$'),
 created_at timestamptz not null default clock_timestamp(),
 primary key(tenant_id,approval_id),
 foreign key(tenant_id,channel_code,approval_id) references catalog.marketplace_effect(tenant_id,channel_code,delivery_key)
);
create trigger marketplace_observation_immutable before update or delete on catalog.marketplace_observation for each row execute function catalog.release_immutable();
commit;
````

### FILE: `docs/MARKETPLACE_INITIAL_REFERENCE.md`

```yaml
block_id: "GO-CONNECTED-MARKETPLACE-MUTATION-INITIAL-DELTA:file3:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "826a9fa07a57c40d620bbe781e92b4ea7f30ceb2888de3fa22faff52f8a087f8"
variables: []
secrets_allowed: false
```

````markdown
# Approved catalog to initial Mercado Libre publication

Scope: local/fixtures, User Products buy_it_now, one selected physical variant,
one immutable normalized PNG, and initial item stock mode. Existing location
mutation modes retain their separate contract. No classified vehicle, Full-stock
write, automatic price rule or whole-provider certification is implied.

## Run the connected owner

Apply migrations through0080 and enable the existing catalog and marketplace
modules in the same Go host. Use the existing MARKETPLACE_* profile/hash/HMAC
settings and future MERCADOLIBRE_ACCESS_TOKEN. The reference profile selects
item stock mode for initial publication. Provider category, listing type,
condition/attributes and seller account must match the target product; provider
validation rejects incompatible settings without choosing commercial policy.

The existing role API is /v1/admin/marketplace. Every request is authenticated,
tenant/organization scoped, private/no-store and bounded. No secrets enter JSON.

1. POST /prepare with operation MEDIA, a new approval_id, current generation,
   variant_id and expires_at (within15minutes). No item_id or media_approval_id.
2. POST the returned exact intent to /requests; a distinct authorized human
   approves its request_sha256 at /requests/{id}/decision.
3. POST {} to /requests/{id}/send. The original durable fence admits one
   multipart POST. Only the current catalog owner's normalized PNG bytes,
   verified against their SHA256, can be uploaded. GET request status after
   any lost API acknowledgement; do not generate another identifier to retry.
4. Prepare operation CREATE with a new approval_id and the accepted
   media_approval_id, same generation/variant. No item_id. The server binds the
   uploaded provider picture ID, current source price/attributes/family name and
   original serial ATP. Submit, obtain distinct human approval and send.
5. POST {} to /requests/{id}/reconcile when the original send is unknown.
   Creation recovery uses a persisted provider ID or bounded seller-SKU search,
   then original item/prices/stock GETs. It compares seller/site/category/SKU,
   User Product, family, listing type, mode, attributes, exact decimal amount,
   quantity and the accepted picture ID. It never issues another item POST.

available_to_promise keeps its original owner and horizon semantics; it can
include inbound supply. Publication neither reserves stock nor invents physical
receipt. An initial quantity zero can produce a paused listing. A confirmed
publication is not a claim that Mercado Libre approved or activated the listing.

## Durable ambiguity and recovery

Each media upload and creation has a separate immutable human request, original
outbound fence and immutable catalog effect. Parsed provider upload/create
acknowledgements persist before subsequent confirmation. A lost local API reply
is recovered by GET. A lost item POST reply is recovered by exact readback.
Existing accepted CREATE evidence prevents another local creation of that SKU,
even if provider search later returns empty. A replacement selling condition or
deleted listing requires explicit target policy and a new admission; it is not
an automatic duplicate.

An upload whose provider response is completely lost and whose picture ID was
never durably observed remains UNKNOWN. There is no documented source-hash
lookup contract and no invented image equivalence. That original operation
cannot resend or be marked accepted from a guessed ID. A separately reviewed
MEDIA request may upload a replacement PNG; it does not clear the unknown
orphan. Unassociated provider media has no product publication effect.

Mercado Libre transforms images. Evidence binds the exact catalog PNG hash to
the multipart request, authenticated provider response/picture ID and listing
readback. It does not claim that provider JPEG bytes equal the original PNG.

## Provenance and verification

All new HTTP/JSON/PG/host/test code is AUTHORED composition, with no corporate
authorship claim. Existing catalog normalization, amounts, serial ATP, manual
approval and fence algorithms retain their original owners and notices.
Official contracts are fixed in docs/marketplace/official-contract.lock.json.
No archived Mercado Libre SDK or unpinned source code is included.

Run TestMarketplaceInitialContracts for seven transport/recovery scenarios;
TestMarketplaceInitialConnected against an owned database with all migrations
for source→API→PG→provider effects,12callers→1upload, human separation, tampering,
unknown/no-retry, accepted upload binding, duplicate CREATE denial and immutable
evidence. Host guards require both effect and observation triggers. Populated
0080down fails and rolls back; empty down/up passes. The exact command-boundary
fuzz campaign includes MEDIA/CREATE plus the prior mutation invariants.

This closes initial publication/media in local fixtures. Existing-item CONTENT
updates are connected in MARKETPLACE_CONTENT_REFERENCE.md. Remaining feeds/comms
and the global TEST02/03/07/T2805 closure remain separately tracked. Provider credentials are supplied later by the user.
````

### FILE: `internal/marketplacebridge/initial.go`

```yaml
block_id: "GO-CONNECTED-MARKETPLACE-MUTATION-INITIAL-DELTA:file4:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "b0ff0816686b6cd95cfab9e67c9a9ab7a07b7901f93bb61a921c892713ec22b6"
variables: []
secrets_allowed: false
```

````go
package marketplacebridge

// AUTHORED official-contract mapping. PNG normalization, monetary representation,
// catalog authority, approval and durable delivery remain their existing owners.
import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"elite.local/enterprise/internal/approval"
	cr "elite.local/enterprise/internal/catalogrelease"
	"elite.local/enterprise/internal/outbounddelivery"
	sv "elite.local/enterprise/internal/storedvaluebridge"
)

var pictureID = regexp.MustCompile(`^[0-9]{1,20}-ML[A-Z][0-9]{1,24}_[0-9]{6}$`)

func ValidPictureID(s string) bool { return pictureID.MatchString(s) }
func creationBody(i Item, price json.Number) any {
	return struct {
		Family     string      `json:"family_name"`
		Category   string      `json:"category_id"`
		Price      json.Number `json:"price"`
		Currency   string      `json:"currency_id"`
		Quantity   int64       `json:"available_quantity"`
		Mode       string      `json:"buying_mode"`
		Listing    string      `json:"listing_type_id"`
		SKU        string      `json:"seller_custom_field"`
		Attributes []Attribute `json:"attributes"`
		Pictures   []Picture   `json:"pictures"`
	}{i.FamilyName, i.CategoryID, price, i.Currency, i.Quantity, i.BuyingMode, i.ListingTypeID, i.SKU, i.Attributes, i.Pictures}
}
func (c *Client) validateInitial(r Intent) error {
	if r.PrepareRequest.Validate() != nil || r.Operation != "MEDIA" && r.Operation != "CREATE" || r.TenantID != c.profile.TenantID || r.OrganizationID != c.profile.OrganizationID || r.ProfileSHA256 != c.profile.SHA256() || r.SellerID != c.profile.SellerID || !cr.ValidSHA(r.MediaSHA256) || !cr.ValidSHA(r.SourceSHA256) || !cr.ValidSHA(r.BeforeSHA256) || r.Method != "POST" || r.StockVersion != "" || r.PriceMinorUnits <= 0 || r.Quantity < 0 {
		return ErrBinding
	}
	b, e := c.profile.Binding(r.VariantID)
	if e != nil || b.SKU != r.SKU {
		return ErrBinding
	}
	var want any
	if r.Operation == "MEDIA" {
		if r.Path != "/pictures/items/upload" {
			return ErrBinding
		}
		want = map[string]string{"catalog_png_sha256": r.MediaSHA256}
	} else {
		x := r.Expected
		if r.Path != "/items" || b.StockMode != "item" || x.ID != "" || x.UserProductID != "" || x.SKU != b.SKU || x.CategoryID != b.CategoryID || x.ListingTypeID != b.ListingTypeID || x.BuyingMode != "buy_it_now" || x.Currency != c.profile.Currency || x.SiteID != c.profile.SiteID || x.Quantity != r.Quantity || !cr.ValidText(x.FamilyName, 256) || !AttributesEqual(b.Attributes, x.Attributes) || len(x.Attributes) != len(b.Attributes) || len(x.Pictures) != 1 || !ValidPictureID(x.Pictures[0].ID) || x.Pictures[0].Source != "" || x.Pictures[0].SecureURL != "" {
			return ErrBinding
		}
		want = creationBody(x, json.Number(sv.Major(r.PriceMinorUnits, c.profile.CurrencyDigits)))
	}
	raw, _, e := cr.Canonical(want)
	actual, _, a := approval.CanonicalPayload(r.Body)
	if e != nil || a != nil || !bytes.Equal(raw, actual) {
		return ErrBinding
	}
	return nil
}

// Search is bounded and authenticated. Empty search is only a precondition,
// never proof that a previous ambiguous write had no effect.
func (c *Client) SearchSKU(ctx context.Context, variant string) ([]string, error) {
	b, e := c.profile.Binding(variant)
	if e != nil {
		return nil, e
	}
	status, raw, _, e := c.call(ctx, "GET", "/users/"+c.profile.SellerID+"/items/search?sku="+url.QueryEscape(b.SKU)+"&limit=2", nil, "")
	var result struct {
		SellerID string   `json:"seller_id"`
		Results  []string `json:"results"`
		Paging   struct {
			Total  int `json:"total"`
			Offset int `json:"offset"`
		} `json:"paging"`
	}
	if e != nil || status != 200 || json.Unmarshal(raw, &result) != nil || result.SellerID != c.profile.SellerID || result.Paging.Offset != 0 || result.Paging.Total < 0 || result.Paging.Total > 1 || len(result.Results) != result.Paging.Total {
		return nil, ErrUnknown
	}
	for _, id := range result.Results {
		if !itemID.MatchString(id) || id[:3] != c.profile.SiteID {
			return nil, ErrBinding
		}
	}
	return result.Results, nil
}
func (c *Client) CreationPreflight(ctx context.Context, variant string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	status, raw, _, e := c.call(ctx, "GET", "/users/me", nil, "")
	var user struct {
		ID   int64    `json:"id"`
		Site string   `json:"site_id"`
		Tags []string `json:"tags"`
	}
	if e != nil || status != 200 || json.Unmarshal(raw, &user) != nil || strconv.FormatInt(user.ID, 10) != c.profile.SellerID || user.Site != c.profile.SiteID || !hasTag(user.Tags, "user_product_seller") {
		return ErrBinding
	}
	ids, e := c.SearchSKU(ctx, variant)
	if e != nil {
		return e
	}
	if len(ids) != 0 {
		return ErrBinding
	}
	return nil
}

type InitialObservation struct {
	ProviderID  string          `json:"provider_id"`
	Operation   string          `json:"operation"`
	MediaSHA256 string          `json:"media_sha256"`
	Response    json.RawMessage `json:"response"`
}
type RecordInitial func(context.Context, InitialObservation) error

func initialReceipt(o InitialObservation) (outbounddelivery.Receipt, error) {
	_, h, e := cr.Canonical(o)
	return outbounddelivery.Receipt{ProviderMessageID: o.ProviderID, EvidenceSHA256: h, AcceptedAt: time.Now().UTC()}, e
}
func ValidateUploadObservation(r Intent, o InitialObservation) error {
	if r.Operation != "MEDIA" || o.Operation != "MEDIA" || o.MediaSHA256 != r.MediaSHA256 || !ValidPictureID(o.ProviderID) {
		return ErrBinding
	}
	var response struct {
		ID         string `json:"id"`
		Variations []struct {
			Size string `json:"size"`
			URL  string `json:"secure_url"`
		} `json:"variations"`
	}
	if _, _, e := approval.CanonicalPayload(o.Response); e != nil {
		return ErrBinding
	}
	if json.Unmarshal(o.Response, &response) != nil || response.ID != o.ProviderID || len(response.Variations) < 1 || len(response.Variations) > 16 {
		return ErrBinding
	}
	for _, v := range response.Variations {
		u, e := url.Parse(v.URL)
		if e != nil || u.Scheme != "https" || u.Host != "http2.mlstatic.com" || u.User != nil || u.RawQuery != "" || u.Fragment != "" || !strings.HasPrefix(u.Path, "/D_NQ_NP_"+o.ProviderID+"-") || !strings.HasSuffix(u.Path, ".jpg") || !regexp.MustCompile(`^[1-9][0-9]{0,4}x[1-9][0-9]{0,4}$`).MatchString(v.Size) {
			return ErrBinding
		}
	}
	return nil
}
func (c *Client) AcceptedUpload(r Intent, o InitialObservation) (outbounddelivery.Receipt, error) {
	if c.validateInitial(r) != nil || ValidateUploadObservation(r, o) != nil {
		return outbounddelivery.Receipt{}, ErrBinding
	}
	return initialReceipt(o)
}
func (c *Client) Upload(ctx context.Context, r Intent, png []byte, record RecordInitial) (outbounddelivery.Receipt, error) {
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	sum := sha256.Sum256(png)
	if c.validateInitial(r) != nil || r.Operation != "MEDIA" || !r.ExpiresAt.After(time.Now()) || len(png) < 8 || len(png) > 4*1024*1024 || !bytes.Equal(png[:8], []byte{137, 80, 78, 71, 13, 10, 26, 10}) || hex.EncodeToString(sum[:]) != r.MediaSHA256 || record == nil {
		return outbounddelivery.Receipt{}, ErrBinding
	}
	// Only already-normalized catalog PNG reaches this method. No new image codec.
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	if e := writer.SetBoundary("elite-" + r.MediaSHA256); e != nil {
		return outbounddelivery.Receipt{}, ErrBinding
	}
	part, e := writer.CreateFormFile("file", r.MediaSHA256+".png")
	if e != nil {
		return outbounddelivery.Receipt{}, ErrBinding
	}
	if _, e = part.Write(png); e != nil {
		return outbounddelivery.Receipt{}, ErrBinding
	}
	if writer.Close() != nil {
		return outbounddelivery.Receipt{}, ErrBinding
	}
	token, e := c.tokens.MercadoLibreAccessToken(ctx)
	if e != nil || len(token) < 20 || len(token) > 4096 || strings.ContainsAny(token, "\r\n") {
		return outbounddelivery.Receipt{}, ErrBinding
	}
	request, e := http.NewRequestWithContext(ctx, "POST", APIBaseURL+"/pictures/items/upload", &body)
	if e != nil {
		return outbounddelivery.Receipt{}, ErrBinding
	}
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", writer.FormDataContentType())
	response, e := c.http.Do(request)
	if e != nil {
		return outbounddelivery.Receipt{}, ErrUnknown
	}
	defer response.Body.Close()
	raw, e := io.ReadAll(io.LimitReader(response.Body, 32769))
	if e != nil || len(raw) > 32768 {
		return outbounddelivery.Receipt{}, ErrUnknown
	}
	if response.StatusCode != 200 && response.StatusCode != 201 {
		return outbounddelivery.Receipt{}, initialRejection(response.StatusCode, raw)
	}
	var upload struct {
		ID string `json:"id"`
	}
	if json.Unmarshal(raw, &upload) != nil {
		return outbounddelivery.Receipt{}, ErrUnknown
	}
	o := InitialObservation{upload.ID, "MEDIA", r.MediaSHA256, raw}
	if ValidateUploadObservation(r, o) != nil {
		return outbounddelivery.Receipt{}, ErrUnknown
	}
	if e = record(ctx, o); e != nil {
		return outbounddelivery.Receipt{}, ErrUnknown
	}
	return initialReceipt(o)
}
func initialRejection(status int, raw []byte) error {
	if status == 400 || status == 401 || status == 403 || status == 409 || status == 422 {
		if _, h, e := approval.CanonicalPayload(raw); e == nil {
			return outbounddelivery.NewTerminalFailure("MERCADOLIBRE_INITIAL_REJECTED", h, nil)
		}
	}
	return ErrUnknown
}
func (c *Client) Create(ctx context.Context, r Intent, record RecordInitial) (outbounddelivery.Receipt, error) {
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	if c.validateInitial(r) != nil || r.Operation != "CREATE" || !r.ExpiresAt.After(time.Now()) || record == nil {
		return outbounddelivery.Receipt{}, ErrBinding
	}
	if e := c.CreationPreflight(ctx, r.VariantID); e != nil {
		return outbounddelivery.Receipt{}, e
	}
	wire, _, e := approval.CanonicalPayload(r.Body)
	if e != nil {
		return outbounddelivery.Receipt{}, ErrBinding
	}
	status, raw, _, e := c.call(ctx, "POST", "/items", wire, "")
	if e != nil {
		return outbounddelivery.Receipt{}, ErrUnknown
	}
	if status != 201 {
		return outbounddelivery.Receipt{}, initialRejection(status, raw)
	}
	var created struct {
		ID string `json:"id"`
	}
	if json.Unmarshal(raw, &created) != nil || !itemID.MatchString(created.ID) || created.ID[:3] != c.profile.SiteID {
		return outbounddelivery.Receipt{}, ErrUnknown
	}
	if e = record(ctx, InitialObservation{created.ID, "CREATE", r.MediaSHA256, raw}); e != nil {
		return outbounddelivery.Receipt{}, ErrUnknown
	}
	return c.ReconcileCreation(ctx, r, created.ID)
}
func (c *Client) ReconcileCreation(ctx context.Context, r Intent, knownID string) (outbounddelivery.Receipt, error) {
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	if c.validateInitial(r) != nil || r.Operation != "CREATE" {
		return outbounddelivery.Receipt{}, ErrBinding
	}
	if knownID == "" {
		ids, e := c.SearchSKU(ctx, r.VariantID)
		if e != nil || len(ids) != 1 {
			return outbounddelivery.Receipt{}, ErrUnknown
		}
		knownID = ids[0]
	}
	observed, e := c.Observe(ctx, r.VariantID, knownID)
	if e != nil {
		return outbounddelivery.Receipt{}, ErrUnknown
	}
	x := observed.Item
	w := r.Expected
	if observed.PriceMinor != strconv.FormatInt(r.PriceMinorUnits, 10) || observed.Quantity != r.Quantity || x.FamilyName != w.FamilyName || x.ListingTypeID != w.ListingTypeID || x.BuyingMode != w.BuyingMode || !AttributesEqual(w.Attributes, x.Attributes) || len(x.Pictures) != 1 || x.Pictures[0].ID != w.Pictures[0].ID || x.Status != "active" && x.Status != "paused" && x.Status != "under_review" {
		return outbounddelivery.Receipt{}, ErrUnknown
	}
	_, h, e := cr.Canonical(map[string]any{"request": r, "observation": observed})
	return outbounddelivery.Receipt{ProviderMessageID: knownID, EvidenceSHA256: h, AcceptedAt: time.Now().UTC()}, e
}
````

### FILE: `internal/marketplacebridge/initial_test.go`

```yaml
block_id: "GO-CONNECTED-MARKETPLACE-MUTATION-INITIAL-DELTA:file5:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "9a6d42f67671fdd97b753496e67dc94eb17d99795b3031df30638c0a3bac2b69"
variables: []
secrets_allowed: false
```

````go
package marketplacebridge_test

import (
	"bytes"
	"context"
	"crypto/sha256"
	mb "elite.local/enterprise/internal/marketplacebridge"
	fixture "elite.local/enterprise/internal/marketplacebridge/testfixture"
	"elite.local/enterprise/internal/outbounddelivery"
	"encoding/hex"
	"encoding/json"
	"errors"
	"image"
	"image/png"
	"testing"
	"time"
)

func TestMarketplaceInitialContracts(t *testing.T) {
	for _, scenario := range []string{"accepted", "lost-create", "lost-upload", "wrong-picture", "duplicate-search", "rejected", "persist-failed"} {
		t.Run(scenario, func(t *testing.T) {
			p := profile("item")
			server := fixture.NewInitial(p)
			defer server.Close()
			client, _ := mb.NewClient(p, tokens{}, server.HTTP())
			document := pub()
			var data bytes.Buffer
			if e := png.Encode(&data, image.NewNRGBA(image.Rect(0, 0, 2, 2))); e != nil {
				t.Fatal(e)
			}
			sum := sha256.Sum256(data.Bytes())
			document.Models[0].Media.SHA256 = hex.EncodeToString(sum[:])
			request := mb.PrepareRequest{ApprovalID: "initial-media-unit", Generation: 1, VariantID: p.Bindings[0].VariantID, Operation: "MEDIA", ExpiresAt: time.Now().Add(time.Minute)}
			media, e := mb.Build(p, document, request, 2, time.Now(), mb.Observation{})
			if e != nil {
				t.Fatal(e)
			}
			var stored mb.InitialObservation
			record := func(_ context.Context, o mb.InitialObservation) error { stored = o; return nil }
			if scenario == "lost-upload" {
				server.DropUpload = true
			}
			receipt, e := client.Upload(context.Background(), media, data.Bytes(), record)
			if scenario == "lost-upload" {
				if !errors.Is(e, mb.ErrUnknown) || stored.ProviderID != "" || server.Uploads != 1 {
					t.Fatal("unknown upload", e)
				}
				return
			}
			if e != nil || receipt.ProviderMessageID != fixture.PictureID || server.UploadedSHA != media.MediaSHA256 {
				t.Fatal("upload binding", e)
			}
			if _, e = client.AcceptedUpload(media, stored); e != nil {
				t.Fatal("durable response recovery", e)
			}
			tampered := append([]byte(nil), data.Bytes()...)
			tampered[len(tampered)-1] ^= 1
			if _, e = client.Upload(context.Background(), media, tampered, record); e == nil || server.Uploads != 1 {
				t.Fatal("source hash bypass")
			}
			request = mb.PrepareRequest{ApprovalID: "initial-create-unit", MediaApprovalID: media.ApprovalID, Generation: 1, VariantID: p.Bindings[0].VariantID, Operation: "CREATE", ExpiresAt: time.Now().Add(time.Minute)}
			creation, e := mb.Build(p, document, request, 2, time.Now(), mb.Observation{Item: mb.Item{Pictures: []mb.Picture{{ID: fixture.PictureID}}}})
			if e != nil {
				t.Fatal(e)
			}
			var pretty bytes.Buffer
			json.Indent(&pretty, creation.Body, "", " ")
			creation.Body = pretty.Bytes()
			if scenario == "lost-create" {
				server.DropCreate = true
			}
			if scenario == "wrong-picture" {
				server.BadPicture = true
			}
			if scenario == "duplicate-search" {
				server.Exists = true
				server.DuplicateSearch = true
			}
			if scenario == "rejected" {
				server.RejectNext = true
			}
			if scenario == "persist-failed" {
				record = func(context.Context, mb.InitialObservation) error { return errors.New("fixture store unavailable") }
			}
			receipt, e = client.Create(context.Background(), creation, record)
			switch scenario {
			case "accepted":
				if e != nil || receipt.ProviderMessageID != server.Item.ID {
					t.Fatal("create", e)
				}
			case "lost-create", "persist-failed":
				if !errors.Is(e, mb.ErrUnknown) {
					t.Fatal("must be unknown", e)
				}
				if _, e = client.ReconcileCreation(context.Background(), creation, ""); e != nil {
					t.Fatal("GET-only recovery", e)
				}
			case "wrong-picture":
				if !errors.Is(e, mb.ErrUnknown) {
					t.Fatal("false success", e)
				}
				if _, e = client.ReconcileCreation(context.Background(), creation, server.Item.ID); e == nil {
					t.Fatal("wrong provider image accepted")
				}
			case "duplicate-search":
				if e == nil || server.Creates != 0 {
					t.Fatal("duplicate precondition")
				}
			case "rejected":
				var terminal *outbounddelivery.TerminalFailure
				if !errors.As(e, &terminal) {
					t.Fatal("documented rejection", e)
				}
			}
			if scenario != "duplicate-search" && server.Creates != 1 {
				t.Fatal("provider POST count", server.Creates)
			}
		})
	}
}
````

### FILE: `internal/marketplacebridge/testfixture/initial.go`

```yaml
block_id: "GO-CONNECTED-MARKETPLACE-MUTATION-INITIAL-DELTA:file6:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "d7eeb6c4bcc02db4b2a933622054b52270bc6ced680b9e63f842b9445a4ea139"
variables: []
secrets_allowed: false
```

````go
package testfixture

// AUTHORED loopback fixture for the locked multipart and User Products contracts.
import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"

	mb "elite.local/enterprise/internal/marketplacebridge"
	sv "elite.local/enterprise/internal/storedvaluebridge"
)

const PictureID = "123456-MLA123456789_092026"

type InitialReceiver struct {
	*Receiver
	Exists, DropUpload, DropCreate, BadPicture, DuplicateSearch bool
	Uploads, Creates                                            int
	UploadedSHA                                                 string
}

func NewInitial(p mb.Profile) *InitialReceiver {
	base := New(p)
	base.Server.Close()
	s := &InitialReceiver{Receiver: base}
	base.Server = httptest.NewServer(http.HandlerFunc(s.serveInitial))
	return s
}
func (s *InitialReceiver) serveInitial(w http.ResponseWriter, r *http.Request) {
	s.Mu.Lock()
	handled := true
	defer func() {
		s.Mu.Unlock()
		if !handled {
			s.Receiver.serve(w, r)
		}
	}()
	w.Header().Set("Content-Type", "application/json")
	if r.Header.Get("Authorization") != "Bearer "+Token {
		w.WriteHeader(401)
		io.WriteString(w, `{"error":"unauthorized"}`)
		return
	}
	if r.Method == "GET" && r.URL.Path == "/users/me" {
		s.Reads++
		json.NewEncoder(w).Encode(map[string]any{"id": s.Item.SellerID, "site_id": s.Profile.SiteID, "tags": []string{"user_product_seller"}})
		return
	}
	if r.Method == "GET" && r.URL.Path == "/users/"+s.Profile.SellerID+"/items/search" {
		s.Reads++
		if r.URL.Query().Get("sku") != s.Profile.Bindings[0].SKU || r.URL.Query().Get("limit") != "2" {
			w.WriteHeader(400)
			return
		}
		ids := []string{}
		if s.Exists {
			ids = append(ids, s.Item.ID)
		}
		if s.DuplicateSearch {
			ids = append(ids, "MLA987654321")
		}
		json.NewEncoder(w).Encode(map[string]any{"seller_id": s.Profile.SellerID, "results": ids, "paging": map[string]any{"total": len(ids), "offset": 0, "limit": 2}})
		return
	}
	if r.Method == "POST" && r.URL.Path == "/pictures/items/upload" {
		s.Uploads++
		reader, e := r.MultipartReader()
		if e != nil {
			w.WriteHeader(400)
			return
		}
		part, e := reader.NextPart()
		if e != nil || part.FormName() != "file" {
			w.WriteHeader(400)
			return
		}
		raw, e := io.ReadAll(io.LimitReader(part, 4*1024*1024+1))
		if e != nil || len(raw) > 4*1024*1024 {
			w.WriteHeader(413)
			return
		}
		sum := sha256.Sum256(raw)
		s.UploadedSHA = hex.EncodeToString(sum[:])
		if part.FileName() != s.UploadedSHA+".png" {
			w.WriteHeader(400)
			return
		}
		if _, e = reader.NextPart(); e != io.EOF {
			w.WriteHeader(400)
			return
		}
		if s.DropUpload {
			s.DropUpload = false
			conn, _, _ := w.(http.Hijacker).Hijack()
			conn.Close()
			return
		}
		w.WriteHeader(201)
		json.NewEncoder(w).Encode(map[string]any{"id": PictureID, "variations": []any{map[string]string{"size": "2x2", "secure_url": "https://http2.mlstatic.com/D_NQ_NP_" + PictureID + "-F.jpg"}}})
		return
	}
	if r.Method == "POST" && r.URL.Path == "/items" {
		s.Creates++
		raw, e := io.ReadAll(io.LimitReader(r.Body, 32769))
		if e != nil || len(raw) > 32768 {
			w.WriteHeader(413)
			return
		}
		var input map[string]json.RawMessage
		var next mb.Item
		var price json.Number
		if json.Unmarshal(raw, &input) != nil || len(input) != 10 || json.Unmarshal(raw, &next) != nil || json.Unmarshal(input["price"], &price) != nil {
			w.WriteHeader(400)
			io.WriteString(w, `{"error":"invalid_item"}`)
			return
		}
		amount, e := sv.Minor(price.String(), s.Profile.CurrencyDigits)
		if e != nil || next.SKU != s.Profile.Bindings[0].SKU || len(next.Pictures) != 1 || next.Pictures[0].ID != PictureID || s.UploadedSHA == "" || next.FamilyName == "" {
			w.WriteHeader(400)
			io.WriteString(w, `{"error":"invalid_binding"}`)
			return
		}
		if s.RejectNext {
			s.RejectNext = false
			w.WriteHeader(400)
			io.WriteString(w, `{"error":"category_validation"}`)
			return
		}
		if s.Exists {
			w.WriteHeader(409)
			io.WriteString(w, `{"error":"fixture_existing_sku"}`)
			return
		}
		next.ID = s.Item.ID
		next.UserProductID = s.Item.UserProductID
		next.SellerID = s.Item.SellerID
		next.SiteID = s.Profile.SiteID
		next.Status = "paused"
		if next.Quantity > 0 {
			next.Status = "active"
		}
		if s.BadPicture {
			next.Pictures[0].ID = "654321-MLA123456789_092026"
		}
		s.Item = next
		s.Quantity = next.Quantity
		s.PriceMinor = amount
		s.Exists = true
		s.Effects++
		s.Writes = append(s.Writes, Write{Path: r.URL.Path, Body: string(raw)})
		if s.DropCreate {
			s.DropCreate = false
			conn, _, _ := w.(http.Hijacker).Hijack()
			conn.Close()
			return
		}
		w.WriteHeader(201)
		io.WriteString(w, `{"id":`+strconv.Quote(next.ID)+`}`)
		return
	}
	if !s.Exists && r.Method == "GET" {
		w.WriteHeader(404)
		io.WriteString(w, `{"error":"not_found"}`)
		return
	}
	handled = false
}
````

### FILE: `internal/platform/postgres/marketplace_initial.go`

```yaml
block_id: "GO-CONNECTED-MARKETPLACE-MUTATION-INITIAL-DELTA:file7:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "e5daf497ce05d6088d788bff421d08247233bc7be9dff21d8a8bdbf926dbe836"
variables: []
secrets_allowed: false
```

````go
package postgres

// AUTHORED binding of the admitted catalog PNG, human approvals and outbound
// fence to provider observations. No credentials or raw bearer tokens persist.
import (
	"context"

	cr "elite.local/enterprise/internal/catalogrelease"
	mb "elite.local/enterprise/internal/marketplacebridge"
	"elite.local/enterprise/internal/outbounddelivery"
	"github.com/jackc/pgx/v5"
)

type marketplaceQuery interface {
	QueryRow(context.Context, string, ...any) pgx.Row
}

func (s *MarketplacePublication) acceptedMedia(ctx context.Context, q marketplaceQuery, r mb.PrepareRequest) (string, error) {
	var payload, raw []byte
	var sha string
	e := q.QueryRow(ctx, `select a.payload,o.observation,o.observation_sha256
 from approval.request a join communication.outbound_delivery d on d.tenant_id=a.tenant_id and d.channel_code='mercadolibre_catalog' and d.delivery_key=a.request_id
 join catalog.marketplace_observation o on o.tenant_id=a.tenant_id and o.approval_id=a.request_id
 where a.tenant_id=$1 and a.request_id=$2 and a.organization_id=$3 and a.kind='marketplace_mutation' and a.state='approved' and d.state='accepted' and o.operation='MEDIA'`, s.profile.TenantID, r.MediaApprovalID, s.profile.OrganizationID).Scan(&payload, &raw, &sha)
	var media mb.Intent
	var observation mb.InitialObservation
	if e != nil || mb.Decode(payload, &media) != nil || mb.Decode(raw, &observation) != nil || media.Operation != "MEDIA" || media.ApprovalID != r.MediaApprovalID || media.Generation != r.Generation || media.VariantID != r.VariantID || media.ProfileSHA256 != s.profile.SHA256() || media.TenantID != s.profile.TenantID || media.OrganizationID != s.profile.OrganizationID || mb.ValidateUploadObservation(media, observation) != nil {
		return "", mb.ErrBinding
	}
	_, h, e := cr.Canonical(observation)
	if e != nil || h != sha {
		return "", mb.ErrBinding
	}
	return observation.ProviderID, nil
}
func (s *MarketplacePublication) initialObservation(ctx context.Context, r mb.PrepareRequest) (mb.Observation, error) {
	if r.Operation != "CREATE" {
		return mb.Observation{}, nil
	}
	pic, e := s.acceptedMedia(ctx, s.catalog.pool, r)
	if e != nil {
		return mb.Observation{}, e
	}
	if e = s.client.CreationPreflight(ctx, r.VariantID); e != nil {
		return mb.Observation{}, e
	}
	return mb.Observation{Item: mb.Item{Pictures: []mb.Picture{{ID: pic}}}}, nil
}
func (s *MarketplacePublication) recordInitial(ctx context.Context, r mb.Intent, o mb.InitialObservation) error {
	if o.Operation != r.Operation || o.MediaSHA256 != r.MediaSHA256 {
		return mb.ErrBinding
	}
	raw, h, e := cr.Canonical(o)
	if e != nil {
		return e
	}
	result, e := s.catalog.pool.Exec(ctx, `insert into catalog.marketplace_observation(tenant_id,approval_id,operation,provider_id,media_sha256,observation,observation_sha256)
 select i.tenant_id,i.delivery_key,$3,$4,$5,$6,$7 from catalog.marketplace_effect i
 join communication.outbound_delivery d using(tenant_id,channel_code,delivery_key)
 where i.tenant_id=$1 and i.delivery_key=$2 and i.channel_code='mercadolibre_catalog' and i.operation=$3
 and i.profile_sha256=$8 and i.source_sha256=$9 and d.state in ('sending','unknown')
 on conflict(tenant_id,approval_id) do nothing`, r.TenantID, r.ApprovalID, r.Operation, o.ProviderID, r.MediaSHA256, raw, h, r.ProfileSHA256, r.SourceSHA256)
	if e != nil {
		return e
	}
	if result.RowsAffected() == 1 {
		return nil
	}
	var stored string
	if e = s.catalog.pool.QueryRow(ctx, `select observation_sha256 from catalog.marketplace_observation where tenant_id=$1 and approval_id=$2`, r.TenantID, r.ApprovalID).Scan(&stored); e != nil || stored != h {
		return mb.ErrBinding
	}
	return nil
}
func (s *MarketplacePublication) sendInitial(ctx context.Context, r mb.Intent) (outbounddelivery.Receipt, error) {
	record := func(ctx context.Context, o mb.InitialObservation) error { return s.recordInitial(ctx, r, o) }
	if r.Operation == "CREATE" {
		return s.client.Create(ctx, r, record)
	}
	png, e := s.catalog.PublicMedia(ctx, r.MediaSHA256)
	if e != nil {
		return outbounddelivery.Receipt{}, e
	}
	return s.client.Upload(ctx, r, png, record)
}
func (s *MarketplacePublication) reconcileInitial(ctx context.Context, r mb.Intent) (outbounddelivery.Receipt, error) {
	var raw []byte
	var h string
	var o mb.InitialObservation
	e := s.catalog.pool.QueryRow(ctx, `select observation,observation_sha256 from catalog.marketplace_observation where tenant_id=$1 and approval_id=$2`, r.TenantID, r.ApprovalID).Scan(&raw, &h)
	if e == nil {
		if mb.Decode(raw, &o) != nil || o.Operation != r.Operation || o.MediaSHA256 != r.MediaSHA256 {
			return outbounddelivery.Receipt{}, mb.ErrBinding
		}
		_, hash, e := cr.Canonical(o)
		if e != nil || hash != h {
			return outbounddelivery.Receipt{}, mb.ErrBinding
		}
	} else if e != pgx.ErrNoRows {
		return outbounddelivery.Receipt{}, e
	}
	if r.Operation == "MEDIA" {
		if e != nil {
			return outbounddelivery.Receipt{}, mb.ErrUnknown
		}
		return s.client.AcceptedUpload(r, o)
	}
	if r.Operation == "CREATE" {
		return s.client.ReconcileCreation(ctx, r, o.ProviderID)
	}
	return outbounddelivery.Receipt{}, mb.ErrBinding
}
````

### FILE: `internal/platform/postgres/marketplace_initial_integration_test.go`

```yaml
block_id: "GO-CONNECTED-MARKETPLACE-MUTATION-INITIAL-DELTA:file8:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "069a5dbffb0e13db0932ea08aad33f9c7b2b7f97666e236c8b62529eee9e326f"
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
	mb "elite.local/enterprise/internal/marketplacebridge"
	fixture "elite.local/enterprise/internal/marketplacebridge/testfixture"
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
	"strings"
	"sync"
	"testing"
	"time"
)

func TestMarketplaceInitialConnected(t *testing.T) {
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
		return identity.Principal{TenantID: tenant, Subject: subject, Permissions: map[string]struct{}{"catalog:draft": {}, "catalog:read": {}, "catalog:publish": {}, "catalog:review:legal": {}, "catalog:review:technical": {}, "catalog:review:media": {}, "catalog:review:publication": {}, "marketplace:request": {}, "marketplace:approve": {}, "marketplace:read": {}, "marketplace:send": {}, "marketplace:reconcile": {}}, Organizations: map[string]struct{}{org: {}}}
	}
	maker, reviewer := principal("marketplace-maker"), principal("marketplace-reviewer")
	cp := cr.Profile{Schema: "elite-catalog-publication/v1", TenantID: tenant, OrganizationID: org, Market: "AR", Currency: "ARS", Origin: "https://catalog.example.invalid", PolicyCode: cr.PolicyCode}
	catalog, e := db.NewCatalogRelease(pool, cp, db.NewCommerce(pool))
	if e != nil {
		t.Fatal(e)
	}
	source, e := catalog.CreateSource(ctx, maker, cr.SourceRequest{CommandID: "marketplace-source", Model: electromobility.Model{Code: "fixture-bicycle", DisplayName: "Fixture bicycle", VehicleClass: "bicycle", Specification: json.RawMessage("{}")}, Variants: []cr.SourceVariant{{Code: "fixture-bicycle-standard", DisplayName: "Standard", BatterySpecification: json.RawMessage("{}"), AmountMinorUnits: 9007199254740991, TaxMode: "not-applicable"}}, ValidFrom: time.Now().Add(-time.Hour)})
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
	profile := mb.Profile{Schema: "elite-marketplace-publication/v1", TenantID: tenant, OrganizationID: org, SellerID: "123456789", SiteID: "MLA", Currency: "ARS", CurrencyDigits: 2, MediaOrigin: cp.Origin, Bindings: []mb.Binding{{VariantID: source.SourceVariantIDs[0], SKU: "fixture-bicycle-standard", CategoryID: "MLA3530", ListingTypeID: "gold_special", StockMode: "item", Attributes: []mb.Attribute{{ID: "ITEM_CONDITION", ValueID: "2230284"}}}}}
	receiver := fixture.NewInitial(profile)
	defer receiver.Close()
	service, e := db.NewMarketplacePublication(catalog, profile, bytes.Repeat([]byte("f"), 32), marketplaceFixtureToken{}, receiver.HTTP())
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
	httpapi.MarketplaceModule{Service: service, TenantID: tenant, OrganizationID: org}.Register(mux, verifier)
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
		request, _ := http.NewRequest(method, api.URL+"/v1/admin/marketplace"+path, bytes.NewReader(raw))
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
	prepare := func(op, id, mediaID string) mb.Intent {
		t.Helper()
		var r mb.Intent
		_, e := call("maker", "POST", "/prepare", mb.PrepareRequest{ApprovalID: id, Generation: published.Generation, VariantID: profile.Bindings[0].VariantID, Operation: op, MediaApprovalID: mediaID, ExpiresAt: time.Now().Add(5 * time.Minute)}, &r)
		if e != nil {
			t.Fatal("prepare", e)
		}
		if r.SourceSHA256 != published.SnapshotSHA256 || r.PriceMinorUnits != 9007199254740991 || r.Quantity != 0 {
			t.Fatal("source/ATP binding", r)
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
		var result db.MarketplaceStatus
		if _, e := call("maker", "GET", "/requests/"+id, nil, &result); e != nil || result.DeliveryState != want {
			t.Fatal("durable state", result.ApprovalID, result.ApprovalState, result.DeliveryState, result.FailureCode, "provider writes", len(receiver.Writes), e)
		}
	}

	mediaRequest := prepare("MEDIA", "initial-media-approved", "")
	for _, token := range []string{"forbidden", "wrong-org", "wrong-tenant"} {
		if code, e := call(token, "POST", "/requests", mediaRequest, nil); e == nil || code != 403 {
			t.Fatal("scope", token, e)
		}
	}
	tampered := mediaRequest
	tampered.MediaSHA256 = strings.Repeat("f", 64)
	if _, e := call("maker", "POST", "/requests", tampered, nil); e == nil {
		t.Fatal("foreign source hash admitted")
	}
	mediaHash := submit(mediaRequest)
	if _, e := call("maker", "POST", "/requests/"+mediaRequest.ApprovalID+"/send", struct{}{}, nil); e == nil {
		t.Fatal("unapproved upload")
	}
	approve(mediaRequest, mediaHash)
	// Concurrent clicks share the original fence; only one multipart upload.
	var wg sync.WaitGroup
	for range 12 {
		wg.Go(func() { call("maker", "POST", "/requests/"+mediaRequest.ApprovalID+"/send", struct{}{}, nil) })
	}
	wg.Wait()
	status(mediaRequest.ApprovalID, "accepted")
	if receiver.Uploads != 1 || receiver.UploadedSHA != mediaRequest.MediaSHA256 {
		t.Fatal("upload/fence/source", receiver.Uploads)
	}
	creation := prepare("CREATE", "initial-create-approved", mediaRequest.ApprovalID)
	if creation.Expected.Pictures[0].ID != fixture.PictureID || creation.SourceSHA256 != mediaRequest.SourceSHA256 {
		t.Fatal("accepted upload binding")
	}
	other := creation
	other.ApprovalID = "initial-create-invalid"
	other.MediaApprovalID = "missing-media-receipt"
	if _, e := call("maker", "POST", "/requests", other, nil); e == nil {
		t.Fatal("unaccepted media")
	}
	approve(creation, submit(creation))
	receiver.Mu.Lock()
	receiver.DropCreate = true
	receiver.Mu.Unlock()
	if _, e := call("maker", "POST", "/requests/"+creation.ApprovalID+"/send", struct{}{}, nil); e == nil {
		t.Fatal("lost provider acknowledgement")
	}
	status(creation.ApprovalID, "unknown")
	if receiver.Creates != 1 || receiver.Effects != 1 {
		t.Fatal("initial write count")
	}
	if _, e := call("maker", "POST", "/requests/"+creation.ApprovalID+"/send", struct{}{}, nil); e == nil || receiver.Creates != 1 {
		t.Fatal("unknown retried")
	}
	if _, e := call("maker", "POST", "/requests/"+creation.ApprovalID+"/reconcile", struct{}{}, nil); e != nil {
		t.Fatal("read-only reconcile", e)
	}
	status(creation.ApprovalID, "accepted")
	if _, e := call("maker", "POST", "/requests/"+creation.ApprovalID+"/send", struct{}{}, nil); e != nil || receiver.Creates != 1 {
		t.Fatal("accepted replay", e)
	}
	if receiver.Item.Quantity != 0 || receiver.Item.Status != "paused" || receiver.PriceMinor != 9007199254740991 {
		t.Fatal("ATP/price readback")
	}
	// No known provider image after a lost multipart acknowledgement: explicit
	// unknown, no forged SHA equivalence, no upload retry, new approval required.
	lost := prepare("MEDIA", "initial-media-lost-ack", "")
	approve(lost, submit(lost))
	receiver.Mu.Lock()
	receiver.DropUpload = true
	receiver.Mu.Unlock()
	if _, e := call("maker", "POST", "/requests/"+lost.ApprovalID+"/send", struct{}{}, nil); e == nil {
		t.Fatal("lost upload")
	}
	status(lost.ApprovalID, "unknown")
	if _, e := call("maker", "POST", "/requests/"+lost.ApprovalID+"/reconcile", struct{}{}, nil); e == nil {
		t.Fatal("invented image recovery")
	}
	if _, e := call("maker", "POST", "/requests/"+lost.ApprovalID+"/send", struct{}{}, nil); e == nil || receiver.Uploads != 2 {
		t.Fatal("upload replay")
	}
	// Another separately approved upload can replace an unassociated orphan,
	// without releasing or resetting the unknown operation.
	replacement := prepare("MEDIA", "initial-media-replacement", "")
	approve(replacement, submit(replacement))
	apiMu.Lock()
	dropAPI = true
	apiMu.Unlock()
	if _, e := call("maker", "POST", "/requests/"+replacement.ApprovalID+"/send", struct{}{}, nil); e == nil {
		t.Fatal("lost local API acknowledgement")
	}
	status(replacement.ApprovalID, "accepted")
	if _, e := call("maker", "POST", "/requests/"+replacement.ApprovalID+"/send", struct{}{}, nil); e != nil || receiver.Uploads != 3 {
		t.Fatal("API recovery replay")
	}
	// A second CREATE stays rejected even if provider search is temporarily empty.
	receiver.Mu.Lock()
	receiver.Exists = false
	receiver.Mu.Unlock()
	duplicate := prepare("CREATE", "initial-create-duplicate", replacement.ApprovalID)
	approve(duplicate, submit(duplicate))
	if _, e := call("maker", "POST", "/requests/"+duplicate.ApprovalID+"/send", struct{}{}, nil); e == nil || receiver.Creates != 1 {
		t.Fatal("historical local SKU exclusion")
	}
	for _, q := range []string{
		`update catalog.marketplace_observation set provider_id='fake' where tenant_id=$1`,
		`delete from catalog.marketplace_observation where tenant_id=$1`,
		`update catalog.marketplace_effect set operation='MEDIA' where tenant_id=$1`,
	} {
		if _, e := pool.Exec(ctx, q, tenant); e == nil {
			t.Fatal("immutable observation/effect")
		}
	}
	var approvals, effects, observations, attempts int
	for _, v := range []struct {
		sql string
		out *int
	}{
		{`select count(*) from approval.request where tenant_id=$1 and kind='marketplace_mutation'`, &approvals},
		{`select count(*) from catalog.marketplace_effect where tenant_id=$1`, &effects},
		{`select count(*) from catalog.marketplace_observation where tenant_id=$1`, &observations},
		{`select coalesce(sum(attempt_count),0) from communication.outbound_delivery where tenant_id=$1 and channel_code='mercadolibre_catalog'`, &attempts},
	} {
		if e := pool.QueryRow(ctx, v.sql, tenant).Scan(v.out); e != nil {
			t.Fatal(e)
		}
	}
	if approvals != 5 || effects != 4 || observations != 2 || attempts != 4 {
		t.Fatal("durable counts", approvals, effects, observations, attempts)
	}
	t.Logf("MARKETPLACE_INITIAL_PASS approvals=%d effects=%d observations=%d attempts=%d multipart=%d item_posts=%d api_posts=%d source=%s image=%s", approvals, effects, observations, attempts, receiver.Uploads, receiver.Creates, apiPosts, published.SnapshotSHA256, mediaRequest.MediaSHA256)
}
````


V402 composed delta: T2805 CONTENT from approved catalog and accepted media: exact name/attributes/picture, single unsold User Product item guard, original human approval/fence and GET recovery. Initial publication and existing mutations retained. Local fixtures proven; remaining feeds/comms open. AUTHORED glue, no corporate authorship.

### FILE: `db/migrations/0081_marketplace_content.down.sql`

```yaml
block_id: "GO-CONNECTED-MARKETPLACE-MUTATION-CONTENT-DELTA:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "6fcf2cdffb6002fc754ebd75c6d6c2f63db71fa2b389faef9b607b90ff56098e"
variables: []
secrets_allowed: false
```

````sql
begin;
-- Restoring the previous constraint rejects CONTENT history atomically.
alter table catalog.marketplace_effect drop constraint marketplace_effect_operation_check;
alter table catalog.marketplace_effect add constraint marketplace_effect_operation_check
 check(operation in ('LEGACY_MUTATION','PRICE','STOCK','PAUSE','RESUME','MEDIA','CREATE'));
commit;
````

### FILE: `db/migrations/0081_marketplace_content.up.sql`

```yaml
block_id: "GO-CONNECTED-MARKETPLACE-MUTATION-CONTENT-DELTA:file2:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "d320d043117d7882f5f63cde065d958edb7580edf61b6355ec75577c1405ff3a"
variables: []
secrets_allowed: false
```

````sql
-- AUTHORED extension of the same immutable provider effect owner.
begin;
alter table catalog.marketplace_effect drop constraint marketplace_effect_operation_check;
alter table catalog.marketplace_effect add constraint marketplace_effect_operation_check
 check(operation in ('LEGACY_MUTATION','PRICE','STOCK','PAUSE','RESUME','MEDIA','CREATE','CONTENT'));
commit;
````

### FILE: `docs/MARKETPLACE_CONTENT_REFERENCE.md`

```yaml
block_id: "GO-CONNECTED-MARKETPLACE-MUTATION-CONTENT-DELTA:file3:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "33557b851bc3a87e6619635d7b8ae603f191282cf048d9fb75c6d5ffc284ab8d"
variables: []
secrets_allowed: false
```

````markdown
# Existing item content from the approved catalog

CONTENT extends the same Mercado Libre owner, Go host, role API, approval kind
and outbound fence. It sends one PUT/items/{id} containing exactly family_name,
attributes and pictures. It does not include price, quantity, status or a legacy
title/variations field. Family name comes from the current immutable catalog
model; attributes are the explicitly bound provider profile; the picture ID
comes from an accepted MEDIA upload for the same generation and variant.

Use the workflow in MARKETPLACE_INITIAL_REFERENCE.md to upload the source PNG.
Then POST /v1/admin/marketplace/prepare with operation CONTENT, a new approval_id,
accepted media_approval_id, current generation, variant_id, existing item_id and
expires_at within15minutes. Submit the returned exact intent, obtain approval
from a distinct authorized reviewer, and send through the existing request path.
GET status after a missing API response. Reconciliation is read-only.

Official User Products documentation states that family_name propagates to all
selling conditions of a User Product, and cannot be changed after an associated
selling condition has sales. This reference therefore requires a bounded
authenticated search returning exactly the selected item and an explicit
sold_quantity zero. Missing sales metadata, more than one associated item,
unexpected owner/UP or changed preconditions fail closed before PUT. No inferred
permission to rename other listings is granted.

Readback confirms the same owner/item/UP and exact requested name, attributes
and picture ID. Provider image bytes remain transformed, not SHA-equivalent to
the normalized PNG. Concurrent changes outside this application cannot be
locked by a local transaction; the receipt identifies the observed scope.

Migration0081 extends the immutable effect's operation constraint. Populated
rollback refuses to erase CONTENT evidence; empty down/up is supported.
Host activation requires that constraint plus the existing immutable guards.

TestMarketplaceContentContracts checks eight contract scenarios. The connected
source/API/PG proof uses3approvals/3fences/3attempts,1multipart,1contentPUT,
lost-response recovery and12replays without another effect. Unsafe scope,
client title tampering and tenant/org/permission boundaries are rejected.
Price and stock remain byte-semantically unchanged by the three-field write.

All new implementation/proof is AUTHORED glue around the existing catalog,
media, human approval, ATP, decimal and fence owners. Eight fixed official
contracts remain in docs/marketplace/official-contract.lock.json. No dependency,
SDK or upstream algorithm changed. Other provider feeds and communications
automation remain separate T2805 work.
````

### FILE: `internal/marketplacebridge/content.go`

```yaml
block_id: "GO-CONNECTED-MARKETPLACE-MUTATION-CONTENT-DELTA:file4:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "8cf471c86a1431f0f2d4cfff91e61963ddd625e5646da19fc3c50dfa28a33638"
variables: []
secrets_allowed: false
```

````go
package marketplacebridge

// AUTHORED field/scope mapping for the fixed official User Products contract.
import (
	"context"
	cr "elite.local/enterprise/internal/catalogrelease"
	"encoding/json"
	"net/url"
	"time"
)

func contentBody(i Item) any {
	return struct {
		Family     string      `json:"family_name"`
		Attributes []Attribute `json:"attributes"`
		Pictures   []Picture   `json:"pictures"`
	}{i.FamilyName, i.Attributes, i.Pictures}
}
func validateContent(p Profile, r Intent) error {
	b, e := p.Binding(r.VariantID)
	x := r.Expected
	if e != nil || r.Operation != "CONTENT" || !cr.ValidSHA(r.MediaSHA256) || !upID.MatchString(x.UserProductID) || x.ID != r.ItemID || x.SKU != b.SKU || !cr.ValidText(x.FamilyName, 256) || x.CategoryID != b.CategoryID || !AttributesEqual(b.Attributes, x.Attributes) || len(b.Attributes) != len(x.Attributes) || len(x.Pictures) != 1 || !ValidPictureID(x.Pictures[0].ID) || x.Pictures[0].Source != "" || x.Pictures[0].SecureURL != "" {
		return ErrBinding
	}
	return nil
}
func (c *Client) SingleUnsoldItem(ctx context.Context, item Item) error {
	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()
	if !itemID.MatchString(item.ID) || !upID.MatchString(item.UserProductID) || item.SoldQuantity == nil || *item.SoldQuantity != 0 {
		return ErrBinding
	}
	status, raw, _, e := c.call(ctx, "GET", "/users/"+c.profile.SellerID+"/items/search?user_product_id="+url.QueryEscape(item.UserProductID)+"&limit=2", nil, "")
	var result struct {
		Seller  string   `json:"seller_id"`
		Results []string `json:"results"`
		Paging  struct {
			Total  int `json:"total"`
			Offset int `json:"offset"`
		} `json:"paging"`
	}
	if e != nil || status != 200 || json.Unmarshal(raw, &result) != nil || result.Seller != c.profile.SellerID || result.Paging.Total != 1 || result.Paging.Offset != 0 || len(result.Results) != 1 || result.Results[0] != item.ID {
		return ErrBinding
	}
	return nil
}
````

### FILE: `internal/marketplacebridge/content_test.go`

```yaml
block_id: "GO-CONNECTED-MARKETPLACE-MUTATION-CONTENT-DELTA:file5:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "76fe24aec2d11b7b3e10e54522b58c90e9a0f64b3b8b2a413166ec868a9dd283"
variables: []
secrets_allowed: false
```

````go
package marketplacebridge_test

import (
	"context"
	mb "elite.local/enterprise/internal/marketplacebridge"
	fixture "elite.local/enterprise/internal/marketplacebridge/testfixture"
	"elite.local/enterprise/internal/outbounddelivery"
	"encoding/json"
	"errors"
	"testing"
	"time"
)

func TestMarketplaceContentContracts(t *testing.T) {
	for _, scenario := range []string{"accepted", "lost", "multi", "sold", "missing-sales", "wrong-picture", "rejected", "tampered"} {
		t.Run(scenario, func(t *testing.T) {
			p := profile("item")
			server := fixture.NewContent(p)
			defer server.Close()
			client, _ := mb.NewClient(p, tokens{}, server.HTTP())
			before, e := client.Observe(context.Background(), p.Bindings[0].VariantID, server.Item.ID)
			if e != nil {
				t.Fatal(e)
			}
			request := mb.PrepareRequest{ApprovalID: "content-unit-approval", MediaApprovalID: "content-unit-media-id", Generation: 1, VariantID: p.Bindings[0].VariantID, Operation: "CONTENT", ItemID: server.Item.ID, ExpiresAt: time.Now().Add(time.Minute)}
			intent, e := mb.Build(p, pub(), request, 2, time.Now(), before, fixture.PictureID)
			if e != nil {
				t.Fatal(e)
			}
			if scenario == "multi" {
				server.Multiple = true
			}
			if scenario == "sold" {
				server.Sales = 1
			}
			if scenario == "missing-sales" {
				server.MissingSales = true
			}
			if scenario == "lost" {
				server.DropNext = true
			}
			if scenario == "rejected" {
				server.RejectNext = true
			}
			if scenario == "tampered" {
				var body map[string]any
				json.Unmarshal(intent.Body, &body)
				body["price"] = 1
				intent.Body, _ = json.Marshal(body)
			}
			_, e = client.Write(context.Background(), intent)
			switch scenario {
			case "accepted", "wrong-picture":
				if e != nil {
					t.Fatal("content effect", e)
				}
				if scenario == "wrong-picture" {
					server.Item.Pictures[0].ID = "654321-MLA123456789_092026"
					if _, e = client.Reconcile(context.Background(), intent); e == nil {
						t.Fatal("wrong picture accepted")
					}
				}
			case "lost":
				if !errors.Is(e, mb.ErrUnknown) {
					t.Fatal("lost write", e)
				}
				if _, e = client.Reconcile(context.Background(), intent); e != nil {
					t.Fatal("read-only recovery", e)
				}
			case "multi", "sold", "missing-sales", "rejected":
				var terminal *outbounddelivery.TerminalFailure
				if !errors.As(e, &terminal) {
					t.Fatal("unsafe scope/rejection", e)
				}
			case "tampered":
				if e == nil {
					t.Fatal("extra monetary field")
				}
			}
			expected := 1
			if scenario == "multi" || scenario == "sold" || scenario == "missing-sales" || scenario == "tampered" {
				expected = 0
			}
			if server.ContentWrites != expected {
				t.Fatal("write count", server.ContentWrites)
			}
			if server.PriceMinor != 9007199254740993 || server.Quantity != 4 {
				t.Fatal("content changed commerce fields")
			}
		})
	}
}
````

### FILE: `internal/marketplacebridge/testfixture/content.go`

```yaml
block_id: "GO-CONNECTED-MARKETPLACE-MUTATION-CONTENT-DELTA:file6:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "f3b506b90f167244f6e32311abfb101e55446091899b46a986168f80ed31e44e"
variables: []
secrets_allowed: false
```

````go
package testfixture

import (
	mb "elite.local/enterprise/internal/marketplacebridge"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
)

// ContentReceiver adds only the new bounded association and three-field update
// contracts. Existing upload/initial/mutation fixture handlers remain unchanged.
type ContentReceiver struct {
	*InitialReceiver
	Multiple, MissingSales bool
	Sales                  int64
	ContentWrites          int
}

func NewContent(p mb.Profile) *ContentReceiver {
	base := NewInitial(p)
	base.Server.Close()
	s := &ContentReceiver{InitialReceiver: base}
	s.Exists = true
	s.Item.FamilyName = "Old supplier name"
	s.Item.Pictures = []mb.Picture{{ID: "654321-MLA123456789_092026"}}
	base.Server = httptest.NewServer(http.HandlerFunc(s.serveContent))
	return s
}
func (s *ContentReceiver) serveContent(w http.ResponseWriter, r *http.Request) {
	s.Mu.Lock()
	handled := true
	defer func() {
		s.Mu.Unlock()
		if !handled {
			s.InitialReceiver.serveInitial(w, r)
		}
	}()
	w.Header().Set("Content-Type", "application/json")
	if r.Header.Get("Authorization") != "Bearer "+Token {
		w.WriteHeader(401)
		return
	}
	if r.Method == "GET" && r.URL.Path == "/users/"+s.Profile.SellerID+"/items/search" && r.URL.Query().Get("user_product_id") != "" {
		s.Reads++
		if r.URL.Query().Get("user_product_id") != s.Item.UserProductID || r.URL.Query().Get("limit") != "2" {
			w.WriteHeader(400)
			return
		}
		ids := []string{s.Item.ID}
		if s.Multiple {
			ids = append(ids, "MLA987654321")
		}
		json.NewEncoder(w).Encode(map[string]any{"seller_id": s.Profile.SellerID, "results": ids, "paging": map[string]any{"total": len(ids), "offset": 0}})
		return
	}
	if r.Method == "PUT" && r.URL.Path == "/items/"+s.Item.ID {
		raw, e := io.ReadAll(io.LimitReader(r.Body, 32769))
		if e != nil || len(raw) > 32768 {
			w.WriteHeader(413)
			return
		}
		var input map[string]json.RawMessage
		var wanted mb.Item
		if json.Unmarshal(raw, &input) != nil || len(input) != 3 || json.Unmarshal(raw, &wanted) != nil || wanted.FamilyName == "" || len(wanted.Pictures) != 1 || wanted.Pictures[0].ID != PictureID {
			w.WriteHeader(400)
			io.WriteString(w, `{"error":"invalid_content"}`)
			return
		}
		s.ContentWrites++
		s.Writes = append(s.Writes, Write{Path: r.URL.Path, Body: string(raw)})
		if s.RejectNext {
			s.RejectNext = false
			w.WriteHeader(400)
			io.WriteString(w, `{"error":"content_validation"}`)
			return
		}
		s.Item.FamilyName = wanted.FamilyName
		s.Item.Attributes = wanted.Attributes
		s.Item.Pictures = wanted.Pictures
		s.Effects++
		if s.DropNext {
			s.DropNext = false
			conn, _, _ := w.(http.Hijacker).Hijack()
			conn.Close()
			return
		}
		json.NewEncoder(w).Encode(s.Item)
		return
	}
	if r.Method == "GET" && r.URL.Path == "/items/"+s.Item.ID {
		s.Item.SoldQuantity = nil
		if !s.MissingSales {
			n := s.Sales
			s.Item.SoldQuantity = &n
		}
	}
	handled = false
}
````

### FILE: `internal/platform/postgres/marketplace_content_integration_test.go`

```yaml
block_id: "GO-CONNECTED-MARKETPLACE-MUTATION-CONTENT-DELTA:file7:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "b7ca82154a178fcd967190e12cc1c34009ed1ea87b54d3a019aeb58b7806ba77"
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
	mb "elite.local/enterprise/internal/marketplacebridge"
	fixture "elite.local/enterprise/internal/marketplacebridge/testfixture"
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
	"strings"
	"sync"
	"testing"
	"time"
)

func TestMarketplaceContentConnected(t *testing.T) {
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
		return identity.Principal{TenantID: tenant, Subject: subject, Permissions: map[string]struct{}{"catalog:draft": {}, "catalog:read": {}, "catalog:publish": {}, "catalog:review:legal": {}, "catalog:review:technical": {}, "catalog:review:media": {}, "catalog:review:publication": {}, "marketplace:request": {}, "marketplace:approve": {}, "marketplace:read": {}, "marketplace:send": {}, "marketplace:reconcile": {}}, Organizations: map[string]struct{}{org: {}}}
	}
	maker, reviewer := principal("marketplace-maker"), principal("marketplace-reviewer")
	cp := cr.Profile{Schema: "elite-catalog-publication/v1", TenantID: tenant, OrganizationID: org, Market: "AR", Currency: "ARS", Origin: "https://catalog.example.invalid", PolicyCode: cr.PolicyCode}
	catalog, e := db.NewCatalogRelease(pool, cp, db.NewCommerce(pool))
	if e != nil {
		t.Fatal(e)
	}
	source, e := catalog.CreateSource(ctx, maker, cr.SourceRequest{CommandID: "marketplace-source", Model: electromobility.Model{Code: "fixture-bicycle", DisplayName: "Fixture bicycle", VehicleClass: "bicycle", Specification: json.RawMessage("{}")}, Variants: []cr.SourceVariant{{Code: "fixture-bicycle-standard", DisplayName: "Standard", BatterySpecification: json.RawMessage("{}"), AmountMinorUnits: 9007199254740991, TaxMode: "not-applicable"}}, ValidFrom: time.Now().Add(-time.Hour)})
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
	profile := mb.Profile{Schema: "elite-marketplace-publication/v1", TenantID: tenant, OrganizationID: org, SellerID: "123456789", SiteID: "MLA", Currency: "ARS", CurrencyDigits: 2, MediaOrigin: cp.Origin, Bindings: []mb.Binding{{VariantID: source.SourceVariantIDs[0], SKU: "fixture-bicycle-standard", CategoryID: "MLA3530", ListingTypeID: "gold_special", StockMode: "item", Attributes: []mb.Attribute{{ID: "ITEM_CONDITION", ValueID: "2230284"}}}}}
	receiver := fixture.NewContent(profile)
	defer receiver.Close()
	service, e := db.NewMarketplacePublication(catalog, profile, bytes.Repeat([]byte("f"), 32), marketplaceFixtureToken{}, receiver.HTTP())
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
	httpapi.MarketplaceModule{Service: service, TenantID: tenant, OrganizationID: org}.Register(mux, verifier)
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
		request, _ := http.NewRequest(method, api.URL+"/v1/admin/marketplace"+path, bytes.NewReader(raw))
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
	prepare := func(op, id, mediaID string) mb.Intent {
		t.Helper()
		var r mb.Intent
		item := ""
		if op == "CONTENT" {
			item = receiver.Item.ID
		}
		_, e := call("maker", "POST", "/prepare", mb.PrepareRequest{ApprovalID: id, Generation: published.Generation, VariantID: profile.Bindings[0].VariantID, Operation: op, MediaApprovalID: mediaID, ItemID: item, ExpiresAt: time.Now().Add(5 * time.Minute)}, &r)
		if e != nil {
			t.Fatal("prepare", e)
		}
		if r.SourceSHA256 != published.SnapshotSHA256 || r.PriceMinorUnits != 9007199254740991 || r.Quantity != 0 {
			t.Fatal("source/ATP binding", r)
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
		var result db.MarketplaceStatus
		if _, e := call("maker", "GET", "/requests/"+id, nil, &result); e != nil || result.DeliveryState != want {
			t.Fatal("durable state", result.ApprovalID, result.ApprovalState, result.DeliveryState, result.FailureCode, "provider writes", len(receiver.Writes), e)
		}
	}

	mediaRequest := prepare("MEDIA", "content-approved-media", "")
	approve(mediaRequest, submit(mediaRequest))
	if _, e := call("maker", "POST", "/requests/"+mediaRequest.ApprovalID+"/send", struct{}{}, nil); e != nil {
		t.Fatal("media", e)
	}
	status(mediaRequest.ApprovalID, "accepted")
	rejected := prepare("CONTENT", "content-scope-reject", mediaRequest.ApprovalID)
	approve(rejected, submit(rejected))
	receiver.Mu.Lock()
	receiver.Multiple = true
	receiver.Mu.Unlock()
	if _, e := call("maker", "POST", "/requests/"+rejected.ApprovalID+"/send", struct{}{}, nil); e == nil {
		t.Fatal("multi-item update")
	}
	status(rejected.ApprovalID, "failed_terminal")
	if receiver.ContentWrites != 0 {
		t.Fatal("unapproved family affected")
	}
	receiver.Mu.Lock()
	receiver.Multiple = false
	receiver.Mu.Unlock()
	content := prepare("CONTENT", "content-approved-loss", mediaRequest.ApprovalID)
	if content.Expected.FamilyName != "Fixture bicycle" || content.Expected.Pictures[0].ID != fixture.PictureID {
		t.Fatal("source mapping")
	}
	tampered := content
	tampered.Expected.FamilyName = "Client invented title"
	if _, e := call("maker", "POST", "/requests", tampered, nil); e == nil {
		t.Fatal("source tamper")
	}
	for _, token := range []string{"forbidden", "wrong-org", "wrong-tenant"} {
		if code, e := call(token, "POST", "/requests", content, nil); e == nil || code != 403 {
			t.Fatal("scope", token, e)
		}
	}
	approve(content, submit(content))
	receiver.Mu.Lock()
	receiver.DropNext = true
	receiver.Mu.Unlock()
	if _, e := call("maker", "POST", "/requests/"+content.ApprovalID+"/send", struct{}{}, nil); e == nil {
		t.Fatal("missing write ack")
	}
	status(content.ApprovalID, "unknown")
	if _, e := call("maker", "POST", "/requests/"+content.ApprovalID+"/send", struct{}{}, nil); e == nil || receiver.ContentWrites != 1 {
		t.Fatal("unknown retried")
	}
	if _, e := call("maker", "POST", "/requests/"+content.ApprovalID+"/reconcile", struct{}{}, nil); e != nil {
		t.Fatal("GET reconcile", e)
	}
	status(content.ApprovalID, "accepted")
	var wg sync.WaitGroup
	for range 12 {
		wg.Go(func() { call("maker", "POST", "/requests/"+content.ApprovalID+"/send", struct{}{}, nil) })
	}
	wg.Wait()
	if receiver.ContentWrites != 1 || receiver.Item.FamilyName != "Fixture bicycle" || receiver.Item.Pictures[0].ID != fixture.PictureID || receiver.PriceMinor != 9007199254740993 || receiver.Quantity != 4 {
		t.Fatal("content changed price/stock or duplicate", receiver.ContentWrites)
	}
	for _, state := range []string{"sold", "missing", "multiple"} {
		receiver.Mu.Lock()
		receiver.Sales = 0
		receiver.MissingSales = false
		receiver.Multiple = false
		if state == "sold" {
			receiver.Sales = 1
		}
		if state == "missing" {
			receiver.MissingSales = true
		}
		if state == "multiple" {
			receiver.Multiple = true
		}
		receiver.Mu.Unlock()
		input := mb.PrepareRequest{ApprovalID: "content-denied-" + state, MediaApprovalID: mediaRequest.ApprovalID, Generation: published.Generation, VariantID: profile.Bindings[0].VariantID, Operation: "CONTENT", ItemID: receiver.Item.ID, ExpiresAt: time.Now().Add(time.Minute)}
		if _, e := call("maker", "POST", "/prepare", input, nil); e == nil {
			t.Fatal("unsafe preparation", state)
		}
	}
	var approvals, effects, attempts int
	for _, v := range []struct {
		sql string
		out *int
	}{
		{`select count(*) from approval.request where tenant_id=$1 and kind='marketplace_mutation'`, &approvals},
		{`select count(*) from catalog.marketplace_effect where tenant_id=$1`, &effects},
		{`select coalesce(sum(attempt_count),0) from communication.outbound_delivery where tenant_id=$1 and channel_code='mercadolibre_catalog'`, &attempts},
	} {
		if e := pool.QueryRow(ctx, v.sql, tenant).Scan(v.out); e != nil {
			t.Fatal(e)
		}
	}
	if approvals != 3 || effects != 3 || attempts != 3 || receiver.Uploads != 1 || receiver.ContentWrites != 1 {
		t.Fatal("durable counts", approvals, effects, attempts)
	}
	t.Logf("MARKETPLACE_CONTENT_PASS approvals=%d effects=%d attempts=%d multipart=%d content_puts=%d api_posts=%d source=%s", approvals, effects, attempts, receiver.Uploads, receiver.ContentWrites, apiPosts, published.SnapshotSHA256)
}
````


V402 composed delta: Current scope correction: existing-item CONTENT now proven; remaining Google Merchant/feeds/comms tracked separately. No additional runtime delta.
