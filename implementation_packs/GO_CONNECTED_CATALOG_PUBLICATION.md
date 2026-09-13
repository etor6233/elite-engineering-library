# Connected approved catalog publication

## 1. Metadata

```yaml
pack_id: "GO-CONNECTED-CATALOG-PUBLICATION"
pack_version: "0.2.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: RECONSTRUCTIBLE
  admission: CONDITIONED
claim: "J3 immutable catalog/price/media publication, distinct shared reviews, rollback, search and durable reference feed connected to existing owners; local reference."
stacks: ["Go 1.26.8", "PostgreSQL 18.6", "Node.js 24.20.0 where frontend selected"]
compatible_with: ["MARKDOWN-COMPOSITOR0.3.0", "V402 existing owner closure"]
incompatible_with: ["unbound tenant/provider/account", "production certification inferred from fixtures", "implicit source or business-policy attribution"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: []
verified_at: "2026-09-12"
```

## 2. Applicability

Blueprint J3 in LIBRARY_INFRASTRUCTURE. One explicit tenant/organization/market/currency/origin stream. No inferred tax/legal rules. Full Go closure required.

## 3. Architecture contract

Original Catalog/Commerce/BC eligibility/Approval/Search/OutboundDelivery owners remain authoritative. Snapshot immutable; every publication creates a fresh effective book via source writers. One transaction for price/search/receipt/outbox. Four distinct-maker shared human reviews; command identity/actor/hash and generations; existing send fence plus immutable feed intent.

## 4. Exact file manifest

```text
CREATE cmd/electromobility-api/catalog_release.go
CREATE cmd/electromobility-api/catalog_release_test.go
CREATE db/migrations/0073_connected_catalog_publication.down.sql
CREATE db/migrations/0073_connected_catalog_publication.up.sql
CREATE db/migrations/0074_catalog_feed_fence.down.sql
CREATE db/migrations/0074_catalog_feed_fence.up.sql
CREATE deploy/catalog/feed.reference.json
CREATE deploy/catalog/publication.reference.json
CREATE docs/CATALOG_PUBLICATION_REFERENCE.md
CREATE internal/catalogrelease/catalog_fuzz_test.go
CREATE internal/catalogrelease/contract.go
CREATE internal/catalogrelease/feed.go
CREATE internal/catalogrelease/media.go
CREATE internal/catalogrelease/media_test.go
CREATE internal/catalogrelease/profile.go
CREATE internal/catalogrelease/profile_test.go
CREATE internal/platform/httpapi/catalog_release.go
CREATE internal/platform/httpapi/catalog_release_test.go
CREATE internal/platform/postgres/catalog_feed.go
CREATE internal/platform/postgres/catalog_feed_integration_test.go
CREATE internal/platform/postgres/catalog_public_owner.go
CREATE internal/platform/postgres/catalog_release.go
CREATE internal/platform/postgres/catalog_release_draft.go
CREATE internal/platform/postgres/catalog_release_http_integration_test.go
CREATE internal/platform/postgres/catalog_release_integration_test.go
CREATE internal/platform/postgres/catalog_release_public.go
CREATE internal/platform/postgres/catalog_release_publish.go
CREATE internal/platform/postgres/catalog_release_storefront_integration_test.go
```

## 5. Materialization blocks

### FILE: `cmd/electromobility-api/catalog_release.go`

```yaml
block_id: "GO-CONNECTED-CATALOG-PUBLICATION:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "17d5bac229040510a38aa795dd5ab2293d67dc8dcc22979ab6a3aa44dd45ea30"
variables: []
secrets_allowed: false
```

````go
package main

// AUTHORED optional host: the same selected Commerce owner supplies pricing.
import (
	"context"
	cr "elite.local/enterprise/internal/catalogrelease"
	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/platform/postgres"
	"encoding/base64"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
)

var errCatalogReleaseConfiguration = errors.New("catalog publication activation invalid")

func init() { catalogReleaseModuleFactory = selectedCatalogReleaseModule }
func selectedCatalogReleaseModule(ctx context.Context, pool *pgxpool.Pool, price *postgres.Commerce, lookup func(string) string) (httpapi.EnterpriseModule, error) {
	if lookup == nil {
		return nil, errCatalogReleaseConfiguration
	}
	switch lookup("CATALOG_RELEASE_ENABLED") {
	case "", "false":
		return nil, nil
	case "true":
	default:
		return nil, errCatalogReleaseConfiguration
	}
	if pool == nil || price == nil {
		return nil, errCatalogReleaseConfiguration
	}
	profile, e := cr.LoadProfile(lookup("CATALOG_RELEASE_PROFILE_FILE"), lookup("CATALOG_RELEASE_PROFILE_SHA256"))
	if e != nil {
		return nil, errCatalogReleaseConfiguration
	}
	var ready bool
	e = pool.QueryRow(ctx, `select (select count(*) from pg_trigger where not tgisinternal and tgenabled in('O','A') and
 (tgname='catalog_release_immutable' and tgrelid in(to_regclass('catalog.release_stream'),to_regclass('catalog.release_media'),
 to_regclass('catalog.release_draft'),to_regclass('catalog.release_review'),to_regclass('catalog.release_command'),
 to_regclass('catalog.release_review_action'),to_regclass('catalog.release_publication'))
 or tgname='catalog_release_approval_guard' and tgrelid=to_regclass('approval.request')
 or tgname='catalog_release_publication_guard' and tgrelid=to_regclass('catalog.release_publication')
 or tgname='catalog_release_price_guard' and tgrelid=to_regclass('pricing.price_book')
 or tgname='catalog_release_entry_guard' and tgrelid=to_regclass('pricing.price_book_entry')
 or tgname='catalog_release_search_guard' and tgrelid=to_regclass('search.document')))=12
 and to_regclass('catalog.release_public_models') is not null
 and exists(select 1 from pg_constraint where conrelid=to_regclass('catalog.release_command') and conname='release_command_kind_check' and pg_get_constraintdef(oid) like '%source%')
 and exists(select 1 from pg_constraint where conrelid=to_regclass('approval.request') and conname='request_kind_check' and pg_get_constraintdef(oid) like '%catalog_review%')`).Scan(&ready)
	if e != nil || !ready {
		return nil, errCatalogReleaseConfiguration
	}
	store, e := postgres.NewCatalogRelease(pool, profile, price)
	if e != nil {
		return nil, errCatalogReleaseConfiguration
	}
	module := httpapi.CatalogReleaseModule{Service: store}
	switch lookup("CATALOG_FEED_ENABLED") {
	case "", "false":
		return module, nil
	case "true":
	default:
		return nil, errCatalogReleaseConfiguration
	}
	feedProfile, e := cr.LoadFeedProfile(lookup("CATALOG_FEED_PROFILE_FILE"), lookup("CATALOG_FEED_PROFILE_SHA256"))
	if e != nil {
		return nil, errCatalogReleaseConfiguration
	}
	key, e := base64.StdEncoding.DecodeString(lookup("CATALOG_FEED_HMAC_KEY"))
	if e != nil || len(key) != 32 {
		return nil, errCatalogReleaseConfiguration
	}
	e = pool.QueryRow(ctx, `select exists(select 1 from pg_trigger where tgrelid=to_regclass('catalog.release_feed_intent') and tgname='catalog_release_feed_immutable' and tgenabled in('O','A'))`).Scan(&ready)
	if e != nil || !ready {
		return nil, errCatalogReleaseConfiguration
	}
	receiver, e := cr.NewFeedHTTP(feedProfile, lookup("CATALOG_FEED_TOKEN"))
	if e != nil {
		return nil, errCatalogReleaseConfiguration
	}
	feeder, e := postgres.NewCatalogFeed(store, key, receiver)
	if e != nil {
		receiver.Close()
		return nil, errCatalogReleaseConfiguration
	}
	module.Feed = feeder
	return module, nil
}
````

### FILE: `cmd/electromobility-api/catalog_release_test.go`

```yaml
block_id: "GO-CONNECTED-CATALOG-PUBLICATION:file2:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "a4063c97a7ffacde9707efca29d84d8779433e2f2dd49ddf144d4ca3669318f3"
variables: []
secrets_allowed: false
```

````go
package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	cr "elite.local/enterprise/internal/catalogrelease"
	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/platform/postgres"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestCatalogReleaseHost(t *testing.T) {
	ctx := context.Background()
	calls := 0
	off := func(k string) string {
		calls++
		if k != "CATALOG_RELEASE_ENABLED" {
			t.Fatal("disabled module read", k)
		}
		return "false"
	}
	if module, e := selectedCatalogReleaseModule(ctx, nil, nil, off); e != nil || module != nil || calls != 1 {
		t.Fatal(module, e, calls)
	}
	rawURL := os.Getenv("PAYMENT_CONNECTED_DB_URL")
	if rawURL == "" {
		t.Skip("owned fixture DB not configured")
	}
	pool, e := pgxpool.New(ctx, rawURL)
	if e != nil {
		t.Fatal(e)
	}
	defer pool.Close()
	profile := cr.Profile{Schema: "elite-catalog-publication/v1", TenantID: "ca6c9cb7-e548-4f19-a2e6-27d2aa2d4893", OrganizationID: "j3-store", Market: "AR", Currency: "ARS", Origin: "https://catalog.example.invalid", PolicyCode: cr.PolicyCode}
	raw, _ := json.Marshal(profile)
	sum := sha256.Sum256(raw)
	hash := hex.EncodeToString(sum[:])
	file := filepath.Join(t.TempDir(), "profile.json")
	if e = os.WriteFile(file, raw, 0600); e != nil {
		t.Fatal(e)
	}
	config := map[string]string{"CATALOG_RELEASE_ENABLED": "true", "CATALOG_RELEASE_PROFILE_FILE": file, "CATALOG_RELEASE_PROFILE_SHA256": hash}
	lookup := func(k string) string { return config[k] }
	price := postgres.NewCommerce(pool)
	if _, e = selectedCatalogReleaseModule(ctx, pool, nil, lookup); e == nil {
		t.Fatal("missing selected price owner admitted")
	}
	config["CATALOG_RELEASE_PROFILE_SHA256"] = ""
	if _, e = selectedCatalogReleaseModule(ctx, pool, price, lookup); e == nil {
		t.Fatal("unbound profile admitted")
	}
	config["CATALOG_RELEASE_PROFILE_SHA256"] = hash
	if module, e := selectedCatalogReleaseModule(ctx, pool, price, lookup); e != nil || module == nil {
		t.Fatal("valid activation", e)
	}
	receiver := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("activation made a provider request")
		w.WriteHeader(500)
	}))
	defer receiver.Close()
	feed := cr.FeedProfile{Schema: "elite-catalog-feed/v1", ReceiverID: "fixture-receiver", BaseURL: receiver.URL, Mode: "fixture"}
	feedRaw, _ := json.Marshal(feed)
	feedSum := sha256.Sum256(feedRaw)
	feedFile := filepath.Join(t.TempDir(), "feed.json")
	if e = os.WriteFile(feedFile, feedRaw, 0600); e != nil {
		t.Fatal(e)
	}
	config["CATALOG_FEED_ENABLED"] = "true"
	config["CATALOG_FEED_PROFILE_FILE"] = feedFile
	config["CATALOG_FEED_PROFILE_SHA256"] = hex.EncodeToString(feedSum[:])
	config["CATALOG_FEED_TOKEN"] = "synthetic-fixture-token"
	if _, e = selectedCatalogReleaseModule(ctx, pool, price, lookup); e == nil {
		t.Fatal("missing feed fence key admitted")
	}
	config["CATALOG_FEED_HMAC_KEY"] = base64.StdEncoding.EncodeToString(bytes.Repeat([]byte{0x42}, 32))
	if module, e := selectedCatalogReleaseModule(ctx, pool, price, lookup); e != nil || module == nil {
		t.Fatal("feed activation", e)
	} else if module.(httpapi.CatalogReleaseModule).Feed == nil {
		t.Fatal("feed absent")
	}
	delete(config, "CATALOG_FEED_TOKEN")
	if _, e = selectedCatalogReleaseModule(ctx, pool, price, lookup); e == nil {
		t.Fatal("missing receiver token admitted")
	}
	config["CATALOG_FEED_TOKEN"] = "synthetic-fixture-token"
	config["CATALOG_FEED_PROFILE_SHA256"] = hash
	if _, e = selectedCatalogReleaseModule(ctx, pool, price, lookup); e == nil {
		t.Fatal("feed hash mismatch admitted")
	}
	config["CATALOG_FEED_PROFILE_SHA256"] = hex.EncodeToString(feedSum[:])
	if _, e = pool.Exec(ctx, `alter table catalog.release_feed_intent disable trigger catalog_release_feed_immutable`); e != nil {
		t.Fatal(e)
	}
	func() {
		defer func() {
			if _, e := pool.Exec(context.Background(), `alter table catalog.release_feed_intent enable trigger catalog_release_feed_immutable`); e != nil {
				t.Error(e)
			}
		}()
		if _, e = selectedCatalogReleaseModule(ctx, pool, price, lookup); e == nil {
			t.Fatal("disabled feed guard admitted")
		}
	}()
	if _, e = pool.Exec(ctx, `alter table approval.request disable trigger catalog_release_approval_guard`); e != nil {
		t.Fatal(e)
	}
	defer func() {
		if _, e := pool.Exec(context.Background(), `alter table approval.request enable trigger catalog_release_approval_guard`); e != nil {
			t.Error(e)
		}
	}()
	if _, e = selectedCatalogReleaseModule(ctx, pool, price, lookup); e == nil {
		t.Fatal("disabled approval guard admitted")
	}
	t.Log("CATALOG_RELEASE_HOST_PASS disabled no-read; selected price owner, exact profile hashes,12base guards, feed immutable guard/key/token required; zero provider requests at startup")
}
````

### FILE: `db/migrations/0073_connected_catalog_publication.down.sql`

```yaml
block_id: "GO-CONNECTED-CATALOG-PUBLICATION:file3:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "45dac70a21dd41f4de640bc93228735692240b88bd3323b95089f510ee5a892b"
variables: []
secrets_allowed: false
```

````sql
begin;
do $$begin
 if exists(select 1 from catalog.release_stream) or exists(select 1 from approval.request where kind='catalog_review')
 then raise exception 'catalog release history exists; downgrade refused';end if;
end $$;
drop trigger catalog_release_search_guard on search.document;
drop function catalog.release_guard_search();
drop trigger catalog_release_approval_guard on approval.request;
drop function catalog.release_guard_approval();
drop trigger catalog_release_price_guard on pricing.price_book;
drop trigger catalog_release_entry_guard on pricing.price_book_entry;
drop function catalog.release_guard_price();
drop view catalog.release_public_models;
drop view catalog.release_public_current;
drop view catalog.release_current;
drop table catalog.release_publication;
drop function catalog.release_guard_publication();
drop table catalog.release_review_action;
drop table catalog.release_command;
drop table catalog.release_review;
drop table catalog.release_draft;
drop table catalog.release_media;
drop table catalog.release_stream;
drop function catalog.release_immutable();
-- AUTHORED service/approval/inventory binding; upstream owners remain authoritative.
alter table approval.request drop constraint request_kind_check;
alter table approval.request add constraint request_kind_check check(kind in
 ('reservation','sale','refund','payment','whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality'));
alter table approval.request drop constraint request_payload_binding_check;
alter table approval.request add constraint request_payload_binding_check check
 (kind not in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality') or
  (organization_id is not null and length(organization_id) between 1 and 128 and
   payload is not null and jsonb_typeof(payload)='object' and pg_column_size(payload)<=65536 and amount_minor_units=0));
create or replace function approval.guard_bound_request() returns trigger language plpgsql as $$
begin
 if tg_op='DELETE' then
  if old.kind in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality') then raise exception 'bound approval request is immutable'; end if;
  return old;
 end if;
 if old.kind in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality') or new.kind in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality') then
  if row(new.tenant_id,new.request_id,new.kind,new.subject_id,new.amount_minor_units,new.requester,new.evidence_sha,new.organization_id,new.payload,new.created_at)
   is distinct from row(old.tenant_id,old.request_id,old.kind,old.subject_id,old.amount_minor_units,old.requester,old.evidence_sha,old.organization_id,old.payload,old.created_at)
   or (old.state<>'pending' and row(new.state,new.decided_at) is distinct from row(old.state,old.decided_at))
  then raise exception 'bound approval request is immutable'; end if;
 end if;
 return new;
end $$;
create or replace function approval.guard_bound_decision() returns trigger language plpgsql as $$
begin
 if exists(select 1 from approval.request where tenant_id=old.tenant_id and request_id=old.request_id and kind in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality')) then
  raise exception 'bound approval decision is immutable';
 end if;
 if tg_op='DELETE' then return old;end if;return new;
end $$;


commit;
````

### FILE: `db/migrations/0073_connected_catalog_publication.up.sql`

```yaml
block_id: "GO-CONNECTED-CATALOG-PUBLICATION:file4:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "21b9be2ddbfb43459eb9648f9aa62e5257a1b9d4bcec197282c53e3d3ee96152"
variables: []
secrets_allowed: false
```

````sql
begin;
-- AUTHORED service/approval/inventory binding; upstream owners remain authoritative.
alter table approval.request drop constraint request_kind_check;
alter table approval.request add constraint request_kind_check check(kind in
 ('reservation','sale','refund','payment','whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review'));
alter table approval.request drop constraint request_payload_binding_check;
alter table approval.request add constraint request_payload_binding_check check
 (kind not in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review') or
  (organization_id is not null and length(organization_id) between 1 and 128 and
   payload is not null and jsonb_typeof(payload)='object' and pg_column_size(payload)<=65536 and amount_minor_units=0));
create or replace function approval.guard_bound_request() returns trigger language plpgsql as $$
begin
 if tg_op='DELETE' then
  if old.kind in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review') then raise exception 'bound approval request is immutable'; end if;
  return old;
 end if;
 if old.kind in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review') or new.kind in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review') then
  if row(new.tenant_id,new.request_id,new.kind,new.subject_id,new.amount_minor_units,new.requester,new.evidence_sha,new.organization_id,new.payload,new.created_at)
   is distinct from row(old.tenant_id,old.request_id,old.kind,old.subject_id,old.amount_minor_units,old.requester,old.evidence_sha,old.organization_id,old.payload,old.created_at)
   or (old.state<>'pending' and row(new.state,new.decided_at) is distinct from row(old.state,old.decided_at))
  then raise exception 'bound approval request is immutable'; end if;
 end if;
 return new;
end $$;
create or replace function approval.guard_bound_decision() returns trigger language plpgsql as $$
begin
 if exists(select 1 from approval.request where tenant_id=old.tenant_id and request_id=old.request_id and kind in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review')) then
  raise exception 'bound approval decision is immutable';
 end if;
 if tg_op='DELETE' then return old;end if;return new;
end $$;

create table catalog.release_stream(
 tenant_id uuid primary key references platform.tenant(tenant_id),
 organization_id text not null,
 market text not null check(market~'^[A-Z]{2}$'),
 currency text not null check(currency~'^[A-Z]{3}$'),
 profile jsonb not null,
 profile_canonical text not null,
 profile_sha256 text not null,
 foreign key(tenant_id,organization_id) references org.organization(tenant_id,organization_id),
 check(profile=profile_canonical::jsonb and profile_sha256=encode(sha256(convert_to(profile_canonical,'UTF8')),'hex'))
);
create table catalog.release_media(
 tenant_id uuid not null references catalog.release_stream(tenant_id),
 media_id text not null, maker text not null, original_sha256 text not null check(original_sha256~'^[0-9a-f]{64}$'),
 sha256 text not null, width integer not null, height integer not null, png bytea not null,
 primary key(tenant_id,media_id),check(width between 1 and 2048 and height between 1 and 2048 and width*height<=1048576),
 check(octet_length(png)<=4194304 and sha256=encode(sha256(png),'hex'))
);
create table catalog.release_draft(
 tenant_id uuid not null references catalog.release_stream(tenant_id),
 draft_id text not null, maker text not null, price_book_id text not null,
 snapshot jsonb not null, snapshot_canonical text not null, snapshot_sha256 text not null,
 created_at timestamptz not null default clock_timestamp(),transaction_id xid8 not null default pg_current_xact_id(),
 primary key(tenant_id,draft_id),foreign key(tenant_id,price_book_id) references pricing.price_book(tenant_id,price_book_id),
 check(octet_length(snapshot_canonical)<=262144 and snapshot=snapshot_canonical::jsonb and snapshot_sha256=encode(sha256(convert_to(snapshot_canonical,'UTF8')),'hex')),
 check(jsonb_typeof(snapshot->'models')='array' and jsonb_array_length(snapshot->'models') between 1 and 32),
 check(jsonb_typeof(snapshot->'variants')='array' and jsonb_array_length(snapshot->'variants') between 1 and 128)
);
create table catalog.release_review(
 tenant_id uuid not null,draft_id text not null,stage text not null check(stage in('legal','technical','media','publication')),
 approval_id text not null,payload_sha256 text not null,
 primary key(tenant_id,draft_id,stage),unique(tenant_id,approval_id),
 foreign key(tenant_id,draft_id) references catalog.release_draft(tenant_id,draft_id),
 foreign key(tenant_id,approval_id) references approval.request(tenant_id,request_id) deferrable initially deferred
);
create table catalog.release_command(
 tenant_id uuid not null references catalog.release_stream(tenant_id),
 command_id text not null,actor text not null,kind text not null check(kind in('media','draft','review','publish')),
 resource_id text not null,request_sha256 text not null check(request_sha256~'^[0-9a-f]{64}$'),
 receipt jsonb not null,transaction_id xid8 not null default pg_current_xact_id(),
 primary key(tenant_id,command_id),check(jsonb_typeof(receipt)='object')
);
create table catalog.release_review_action(
 tenant_id uuid not null,command_id text not null,approval_id text not null,actor text not null,
 approved boolean not null,evidence_sha256 text not null check(evidence_sha256~'^[0-9a-f]{64}$'),
 transaction_id xid8 not null default pg_current_xact_id(),
 primary key(tenant_id,approval_id),unique(tenant_id,command_id),
 foreign key(tenant_id,command_id) references catalog.release_command(tenant_id,command_id) deferrable initially deferred,
 foreign key(tenant_id,approval_id) references catalog.release_review(tenant_id,approval_id)
);
create table catalog.release_publication(
 tenant_id uuid not null,generation bigint not null check(generation>0),draft_id text not null,command_id text not null,
 actor text not null,reason text not null,effective_price_book_id text not null,created_at timestamptz not null default clock_timestamp(),
 transaction_id xid8 not null default pg_current_xact_id(),
 primary key(tenant_id,generation),unique(tenant_id,command_id),
 foreign key(tenant_id,effective_price_book_id) references pricing.price_book(tenant_id,price_book_id),
 foreign key(tenant_id,draft_id) references catalog.release_draft(tenant_id,draft_id),
 foreign key(tenant_id,command_id) references catalog.release_command(tenant_id,command_id) deferrable initially deferred
);
create index catalog_release_draft_book on catalog.release_draft(tenant_id,price_book_id);
create unique index catalog_release_effective_book on catalog.release_publication(tenant_id,effective_price_book_id);
create index catalog_release_media_digest on catalog.release_media(tenant_id,sha256);
create function catalog.release_immutable() returns trigger language plpgsql as $$
begin raise exception 'catalog release facts are immutable';end $$;
do $$ declare n text;begin
 foreach n in array array['release_stream','release_media','release_draft','release_review','release_command','release_review_action','release_publication'] loop
 execute format('create trigger catalog_release_immutable before update or delete on catalog.%I for each row execute function catalog.release_immutable()',n);
 end loop;
end $$;
create view catalog.release_current as
 select p.*,d.snapshot,d.snapshot_canonical,d.snapshot_sha256,d.price_book_id as source_price_book_id
 from catalog.release_publication p join catalog.release_draft d using(tenant_id,draft_id)
 where p.generation=(select max(x.generation) from catalog.release_publication x where x.tenant_id=p.tenant_id);
create view catalog.release_public_current as
 select c.* from catalog.release_current c join pricing.price_book b on b.tenant_id=c.tenant_id and b.price_book_id=c.effective_price_book_id
 join catalog.release_stream st on st.tenant_id=c.tenant_id join platform.tenant tenant on tenant.tenant_id=c.tenant_id join org.organization org on org.tenant_id=st.tenant_id and org.organization_id=st.organization_id
 where tenant.status='active' and org.status='active' and b.status='active' and b.valid_from<=statement_timestamp() and (b.valid_until is null or b.valid_until>statement_timestamp());
create view catalog.release_public_models as
 select c.tenant_id,m->>'id' as model_id,m->>'code' as model_code,m->>'displayName' as display_name,
 m->>'vehicleClass' as vehicle_class,m->'specification' as specification
 from catalog.release_public_current c cross join lateral jsonb_array_elements(c.snapshot->'models') m
 union all
 select tenant_id,model_id,model_code,display_name,vehicle_class,specification from catalog.vehicle_model m
 where m.lifecycle_state='active' and m.publicly_visible and not exists(select 1 from catalog.release_stream s where s.tenant_id=m.tenant_id);
create function catalog.release_guard_price() returns trigger language plpgsql as $$
declare t uuid;bid text;market_ text;currency_ text;begin
 if tg_op='DELETE' then t:=old.tenant_id;bid:=old.price_book_id;else t:=new.tenant_id;bid:=new.price_book_id;end if;
 if tg_table_name='price_book_entry' then
  if (exists(select 1 from catalog.release_draft where tenant_id=t and price_book_id=bid) or exists(select 1 from catalog.release_publication where tenant_id=t and effective_price_book_id=bid)) then raise exception 'captured price entries immutable';end if;
 else
  if tg_op='DELETE' then
   if (exists(select 1 from catalog.release_draft where tenant_id=t and price_book_id=bid) or exists(select 1 from catalog.release_publication where tenant_id=t and effective_price_book_id=bid)) then raise exception 'captured price book immutable';end if;
  elsif tg_op='UPDATE' then
   if (exists(select 1 from catalog.release_draft where tenant_id=old.tenant_id and price_book_id=old.price_book_id) or exists(select 1 from catalog.release_publication where tenant_id=old.tenant_id and effective_price_book_id=old.price_book_id))
   and row(new.tenant_id,new.price_book_id,new.market,new.currency,new.valid_from,new.valid_until)
    is distinct from row(old.tenant_id,old.price_book_id,old.market,old.currency,old.valid_from,old.valid_until)
   then raise exception 'captured price header immutable';end if;
  end if;
  if tg_op<>'DELETE' then
   market_:=new.market;currency_:=new.currency;
   if exists(select 1 from catalog.release_stream where tenant_id=t and market=market_ and currency=currency_)
   and ((tg_op='INSERT' and new.status='active') or (tg_op='UPDATE' and new.status is distinct from old.status))
   and not exists(select 1 from catalog.release_publication where tenant_id=t and transaction_id=pg_current_xact_id())
   then raise exception 'price publication command required';end if;
  end if;
 end if;
 if tg_op='DELETE' then return old;end if;return new;
end $$;
create trigger catalog_release_price_guard before insert or update or delete on pricing.price_book for each row execute function catalog.release_guard_price();
create trigger catalog_release_entry_guard before insert or update or delete on pricing.price_book_entry for each row execute function catalog.release_guard_price();
create function catalog.release_guard_approval() returns trigger language plpgsql as $$
declare binding record;action record;begin
 if new.kind<>'catalog_review' then return new;end if;
 select r.*,d.maker,d.snapshot_sha256,d.transaction_id,s.organization_id into binding
 from catalog.release_review r join catalog.release_draft d using(tenant_id,draft_id)
 join catalog.release_stream s using(tenant_id)
 where r.tenant_id=new.tenant_id and r.approval_id=new.request_id;
 if not found or new.requester<>binding.maker or new.organization_id<>binding.organization_id
  or new.evidence_sha<>binding.payload_sha256 or new.amount_minor_units<>0
  or new.payload->>'draft_id' is distinct from binding.draft_id or new.payload->>'stage' is distinct from binding.stage
  or new.payload->>'snapshot_sha256' is distinct from binding.snapshot_sha256
 then raise exception 'catalog approval binding mismatch';end if;
 if tg_op='INSERT' then
  if binding.transaction_id<>pg_current_xact_id() then raise exception 'catalog draft command required';end if;
 elsif new.state<>old.state then
  select a.* into action from catalog.release_review_action a
  join catalog.release_command c using(tenant_id,command_id)
  join approval.decision v on v.tenant_id=a.tenant_id and v.request_id=a.approval_id
  where a.tenant_id=new.tenant_id and a.approval_id=new.request_id and a.transaction_id=pg_current_xact_id()
   and c.transaction_id=pg_current_xact_id() and c.kind='review' and c.actor=a.actor
   and a.actor=v.reviewer and a.actor<>binding.maker and a.approved=v.approved
   and new.state=case when a.approved then 'approved' else 'rejected' end;
  if not found then raise exception 'catalog review command required';end if;
  if binding.stage='publication' and action.approved and
   (select count(*) from catalog.release_review r join approval.request q on q.tenant_id=r.tenant_id and q.request_id=r.approval_id
    where r.tenant_id=new.tenant_id and r.draft_id=binding.draft_id and r.stage<>'publication' and q.state='approved')<>3
  then raise exception 'prior catalog reviews required';end if;
 end if;
 return new;
end $$;
create constraint trigger catalog_release_approval_guard after insert or update on approval.request deferrable initially deferred for each row execute function catalog.release_guard_approval();
create function catalog.release_guard_publication() returns trigger language plpgsql as $$
declare d record;s record;c record;expected bigint;begin
 select * into d from catalog.release_draft where tenant_id=new.tenant_id and draft_id=new.draft_id;
 select * into s from catalog.release_stream where tenant_id=new.tenant_id;
 select * into c from catalog.release_command where tenant_id=new.tenant_id and command_id=new.command_id;
 select coalesce(max(generation),0)+1 into expected from catalog.release_publication where tenant_id=new.tenant_id and generation<new.generation;
 if new.transaction_id<>pg_current_xact_id() or c.transaction_id<>pg_current_xact_id() or c.kind<>'publish' or c.actor<>new.actor
  or c.resource_id<>new.draft_id or new.generation<>expected or new.actor=d.maker
  or (c.receipt->>'generation')::bigint is distinct from new.generation or c.receipt->>'snapshot_sha256' is distinct from d.snapshot_sha256
 then raise exception 'catalog publication receipt mismatch';end if;
 if (select count(*) from catalog.release_review r join approval.request q on q.tenant_id=r.tenant_id and q.request_id=r.approval_id
  where r.tenant_id=new.tenant_id and r.draft_id=new.draft_id and q.kind='catalog_review' and q.state='approved')<>4
 then raise exception 'catalog publication needs four approved reviews';end if;
 if not exists(select 1 from pricing.price_book b where b.tenant_id=new.tenant_id and b.price_book_id=new.effective_price_book_id
  and b.market=d.snapshot->'profile'->>'market' and b.currency=d.snapshot->'profile'->>'currency'
  and b.valid_from=(d.snapshot->>'valid_from')::timestamptz and b.valid_until is not distinct from (d.snapshot->>'valid_until')::timestamptz)
 or c.receipt->>'effective_price_book_id' is distinct from new.effective_price_book_id
 or (select count(*) from pricing.price_book_entry e where e.tenant_id=new.tenant_id and e.price_book_id=new.effective_price_book_id)<>jsonb_array_length(d.snapshot->'variants')
 or exists(select 1 from jsonb_array_elements(d.snapshot->'variants') v where not exists(
  select 1 from pricing.price_book_entry e where e.tenant_id=new.tenant_id and e.price_book_id=new.effective_price_book_id and e.variant_id=v->>'id'
   and e.amount_minor_units=(v->>'amount_minor_units')::bigint and e.tax_mode=v->>'tax_mode'))
 then raise exception 'effective price book differs from approved snapshot';end if;

 if not exists(select 1 from pricing.price_book b where b.tenant_id=new.tenant_id and b.price_book_id=new.effective_price_book_id and b.status='active'
  and b.valid_from<=clock_timestamp() and (b.valid_until is null or b.valid_until>clock_timestamp()))
 then raise exception 'catalog price expired or not active at commit';end if;
 if exists(select 1 from pricing.price_book b where b.tenant_id=new.tenant_id and b.market=s.market and b.currency=s.currency and b.status='active' and b.price_book_id<>new.effective_price_book_id)
 then raise exception 'catalog competing active price';end if;
 if (select count(*) from search.document x where x.tenant_id=new.tenant_id and x.kind='catalog-model')<>jsonb_array_length(d.snapshot->'models')
 or exists(select 1 from jsonb_array_elements(d.snapshot->'models') m where not exists(
  select 1 from search.document x where x.tenant_id=new.tenant_id and x.document_id='catalog:'||(m->>'id') and x.kind='catalog-model'
   and x.title=m->>'displayName' and x.body=m->>'canonical_url' and x.external_id=m->>'id' and x.content_sha256=d.snapshot_sha256))
 then raise exception 'catalog search projection mismatch';end if;
 return new;
end $$;
create constraint trigger catalog_release_publication_guard after insert on catalog.release_publication deferrable initially deferred for each row execute function catalog.release_guard_publication();
create function catalog.release_guard_search() returns trigger language plpgsql as $$
declare t uuid;begin
 if tg_op='DELETE' then t:=old.tenant_id;else t:=new.tenant_id;end if;
 if ((tg_op<>'INSERT' and old.kind='catalog-model') or (tg_op<>'DELETE' and new.kind='catalog-model'))
 and exists(select 1 from catalog.release_stream where tenant_id=t)
 and not exists(select 1 from catalog.release_publication where tenant_id=t and transaction_id=pg_current_xact_id())
 then raise exception 'catalog projection requires publication';end if;
 if tg_op='DELETE' then return old;end if;return new;
end $$;
create trigger catalog_release_search_guard before insert or update or delete on search.document for each row execute function catalog.release_guard_search();

commit;
````

### FILE: `db/migrations/0074_catalog_feed_fence.down.sql`

```yaml
block_id: "GO-CONNECTED-CATALOG-PUBLICATION:file5:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "6224fa9975d454243562fcf89ff25cf14e50a13e5803ff4973608a7a6dfb53df"
variables: []
secrets_allowed: false
```

````sql
begin;
do $$begin
 if exists(select 1 from catalog.release_feed_intent) then raise exception 'catalog feed history exists; downgrade refused';end if;
end$$;
drop table catalog.release_feed_intent;
commit;
````

### FILE: `db/migrations/0074_catalog_feed_fence.up.sql`

```yaml
block_id: "GO-CONNECTED-CATALOG-PUBLICATION:file6:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "44066d5beac0f81ff58bab52ee5b6f6a181c76223faf5db2d31056616d964a6d"
variables: []
secrets_allowed: false
```

````sql
begin;
-- AUTHORED immutable admission binding; existing outbound_delivery owns states.
create table catalog.release_feed_intent(
 tenant_id uuid not null,channel_code text not null check(channel_code='catalog_feed'),
 delivery_key text not null,receiver_id text not null,receiver_profile_sha256 text not null check(receiver_profile_sha256~'^[0-9a-f]{64}$'),
 generation bigint not null,actor text not null check(length(actor) between 1 and 128),
 request_sha256 text not null check(request_sha256~'^[0-9a-f]{64}$'),
 payload_sha256 text not null,source_sha256 text not null check(source_sha256~'^[0-9a-f]{64}$'),
 payload bytea not null,created_at timestamptz not null default clock_timestamp(),
 primary key(tenant_id,channel_code,delivery_key),unique(tenant_id,generation,receiver_id),
 foreign key(tenant_id,generation) references catalog.release_publication(tenant_id,generation),
 foreign key(tenant_id,channel_code,delivery_key) references communication.outbound_delivery(tenant_id,channel_code,delivery_key),
 check(octet_length(payload)<=524288 and payload_sha256=encode(sha256(payload),'hex'))
);
create trigger catalog_release_feed_immutable before update or delete on catalog.release_feed_intent for each row execute function catalog.release_immutable();
commit;
````

### FILE: `deploy/catalog/feed.reference.json`

```yaml
block_id: "GO-CONNECTED-CATALOG-PUBLICATION:file7:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "aada0362d91babaccd41583b81f03b44c4842b9405ffb51cc672ee54ace12cf9"
variables: []
secrets_allowed: false
```

````json
{
  "schema": "elite-catalog-feed/v1",
  "receiver_id": "fixture-receiver",
  "base_url": "http://127.0.0.1:4199",
  "mode": "fixture"
}
````

### FILE: `deploy/catalog/publication.reference.json`

```yaml
block_id: "GO-CONNECTED-CATALOG-PUBLICATION:file8:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "ec8c9a87108d782f7f7ceefa95fbb1b643310b7d05be9d1c002ac475ac885d27"
variables: []
secrets_allowed: false
```

````json
{
  "schema": "elite-catalog-publication/v1",
  "tenant_id": "ca6c9cb7-e548-4f19-a2e6-27d2aa2d4893",
  "organization_id": "j3-store",
  "market": "AR",
  "currency": "ARS",
  "origin": "https://catalog.example.invalid",
  "policy_code": "catalog-reference-v1"
}
````

### FILE: `docs/CATALOG_PUBLICATION_REFERENCE.md`

```yaml
block_id: "GO-CONNECTED-CATALOG-PUBLICATION:file9:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "87574e00cda2ac33f62b252080368ba2816618c623f6cb20c3db74ed49e7da8f"
variables: []
secrets_allowed: false
```

````markdown
# Connected catalog publication — J3

LIBRARY_INFRASTRUCTURE / local reference. AUTHORED orchestration reuses the
existing Catalog, Commerce, BC eligibility, shared Approval, Search and
OutboundDelivery owners. No new price, tax, legal, ranking or business policy.

Apply migrations through0074 in order. Select the complete franchise composition.
CATALOG_RELEASE_ENABLED=true requires CATALOG_RELEASE_PROFILE_FILE (absolute
regular file) and CATALOG_RELEASE_PROFILE_SHA256 (SHA256 of exact UTF-8 bytes).
deploy/catalog/publication.reference.json is an explicit synthetic profile.
Choose an existing active tenant/organization, market/currency and HTTPS origin;
never infer fiscal settings. Host requires the same selected Commerce repository
and12enabled database guards. Disabled mode reads no profile or provider keys.

Draft creation captures existing validated models/variants, a bounded existing
price book and normalized PNG assets.1–32models/1–128variants; snapshot at most
256KiB. PNG input at most1MiB, dimensions at most2048each/1048576totalpixels,
normalized output at most4MiB. Go1.26.8 image/png supplies decoding/reencoding.
This removes uninterpreted trailers; it is not antivirus or libxml2 assurance.

Four shared human review requests bind the exact draft SHA: legal, technical,
media and publication. Each reviewer differs from the maker; the final stage
requires three preceding approvals. These are recorded decisions, not automatic
legal/technical certification. Publish requires all four approvals and current
database-clock price validity. Generic approval bypass is rejected.

Each publication/rollback copies the exact approved price header and entries
into a fresh effective book via the original Commerce transaction writer, then
activates it once. This preserves created1/activated2 outbox semantics, old quote
references and immutable history. One transaction publishes effective pricing,
search projection, receipt and outbox. Rollback appends a generation and refuses
expired prices. Source model/price edits never mutate a captured snapshot.

Private routes require a single organization_id query and verifier-supplied
tenant/actor. Bodies cannot choose identity. POST /v1/admin/catalog/media/{command}
uses image/png and catalog:draft; POST /drafts and GET /drafts/{id} use catalog:draft
and catalog:read respectively. POST /drafts/{id}/review/{stage} uses
catalog:review:{stage}; POST /drafts/{id}/publish uses catalog:publish.
GET /commands/{command} recovers actor/request-hash-bound receipts after a lost
response. Do not issue another command ID merely because the transport failed.
Versions and minor amounts are quoted int64 decimal strings. JSON32KiB boundary,
duplicate/unknown fields and wrong organization/permissions are rejected.

GET /v1/public/catalog and /feed return identical typed current publication bytes.
They omit private profile/reviewer configuration, include the public tenant code,
and emit ETag generation/sourceSHA with max-age=0,must-revalidate. Current-only
normalized media and scoped PG full-text search use the same approved source.
The existing public model and price owners read this publication when selected;
unconfigured tenants and narrower Go profiles preserve their existing behavior.

Optional reference feed: CATALOG_FEED_ENABLED=true additionally requires exact
CATALOG_FEED_PROFILE_FILE/SHA256, CATALOG_FEED_HMAC_KEY (base6432bytes) and
CATALOG_FEED_TOKEN. The fixture profile uses literal loopback HTTP; other profiles
require HTTPS. No provider call occurs during activation. Real values are supplied
later from target secret management, never in the profile/receipt. Enabling feed
without these inputs or its immutable database guard fails closed.

POST /v1/admin/catalog/feed/{generation} and POST the same path/reconcile require
catalog:feed; GET uses catalog:read. The existing OutboundDelivery state machine
and fence own the send. An immutable intent binds actor/generation/source/payload/
receiver profile hashes in the same claim transaction. Receiver protocol:
POST and GET /publications/{idempotency-key}; Authorization Bearer,
Idempotency-Key, X-Catalog-Generation and X-Content-SHA256 headers. ACK JSON binds
status accepted/rejected, quoted generation, payload_sha256, remote_id and
accepted_at. Accepted200/201 must match the request; terminal rejected422 remains
rejected. An unknown send is reconciled by GET, never automatically re-POSTed.
The receiver preserves per-key receipts and rejects older generations.
This is an explicit reference receiver contract. Mapping to the actual admitted
Google/Meta/ML SDKs belongs to T2805, not a claim of credential-only readiness yet.

Actual PG18.6/71migrations, Go1.26.8 and Next16.3.4/Node24.20.0/Chromium:
3snapshots/12reviews/5publications;1new/7replays;13table rollback; source isolation;
lost publication response recovered by GET;3stale ETags; four reference feed sends,
one terminal rejection and one lost accepted response reconciled once. Actual
storefront follows generations1/2/3/5 with canonical/sitemap/robots/PNG/404.
73/74down refuses populated history intact; empty74→73down then73→74up passes.
Fuzz finite2s/7seeds/121059executions. Narrow compile is not a runtime claim.
Canonical CATALOG_CONNECTED_RELEASE_V402.md/json binds exact source and receipts.
Role editor UI T2804, provider mappings T2805, identity/SCA T2803, operations and
signed release remain their own gates. No production or corporate glue attribution.
````

### FILE: `internal/catalogrelease/catalog_fuzz_test.go`

```yaml
block_id: "GO-CONNECTED-CATALOG-PUBLICATION:file10:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "597beb2f9543f37a3f60203c76e9b72abd68bad0293b4a6b5102b4f0bff0be2a"
variables: []
secrets_allowed: false
```

````go
package catalogrelease

import (
	"bytes"
	"crypto/sha256"
	"elite.local/enterprise/internal/approval"
	"encoding/hex"
	"image"
	"image/png"
	"testing"
)

func FuzzCatalogPublicationBoundaries(f *testing.F) {
	var seed bytes.Buffer
	if e := png.Encode(&seed, image.NewNRGBA(image.Rect(0, 0, 2, 2))); e != nil {
		f.Fatal(e)
	}
	for _, raw := range [][]byte{seed.Bytes(), append(append([]byte(nil), seed.Bytes()...), []byte("uninterpreted trailer")...),
		[]byte(`{"command_id":"publish-one","draft_id":"draft-one","expected_generation":"2","reason":"approved snapshot"}`),
		[]byte(`{"command_id":"one","command_id":"two"}`), []byte(`{"amount_minor_units":"9223372036854775807"}`), []byte{}, []byte("<svg/>")} {
		f.Add(raw)
	}
	f.Fuzz(func(t *testing.T, raw []byte) {
		if len(raw) > 1<<20 {
			return
		}
		if canonical, hash, e := approval.CanonicalPayload(raw); e == nil {
			again, againHash, e := approval.CanonicalPayload(canonical)
			if e != nil || !bytes.Equal(canonical, again) || hash != againHash {
				t.Fatal("canonical command identity unstable")
			}
		}
		media, clean, e := NormalizePNG(raw)
		if e != nil {
			return
		}
		originalHash := sha256.Sum256(raw)
		cleanHash := sha256.Sum256(clean)
		if media.OriginalSHA256 != hex.EncodeToString(originalHash[:]) || media.SHA256 != hex.EncodeToString(cleanHash[:]) ||
			media.Width < 1 || media.Height < 1 || media.Width > 2048 || media.Height > 2048 || media.Width*media.Height > 1048576 || len(clean) > 4<<20 {
			t.Fatal("normalized media escaped declared boundary")
		}
		next, twice, e := NormalizePNG(clean)
		if e != nil || next.SHA256 != media.SHA256 || !bytes.Equal(clean, twice) {
			t.Fatal("normalized media is not stable")
		}
	})
}
````

### FILE: `internal/catalogrelease/contract.go`

```yaml
block_id: "GO-CONNECTED-CATALOG-PUBLICATION:file11:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "21e17f7e93ee19200ecc45e8272ac82ea5641cc8833d7a6f194e729a8dc6f227"
variables: []
secrets_allowed: false
```

````go
// AUTHORED publication binding contracts for existing catalog/price/approval
// owners. No legal, tax, pricing or ranking rules are inferred.
package catalogrelease

import (
	"elite.local/enterprise/internal/approval"
	"elite.local/enterprise/internal/electromobility"
	"encoding/json"
	"errors"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

var ErrInvalid = errors.New("invalid catalog publication request")
var ErrNotFound = errors.New("catalog publication outside scope")
var ErrConflict = errors.New("catalog publication conflict")
var identifier = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:-]{0,63}$`)
var digestPattern = regexp.MustCompile(`^[0-9a-f]{64}$`)

const PolicyCode = "catalog-reference-v1"

func ValidID(s string) bool  { return identifier.MatchString(s) }
func ValidSHA(s string) bool { return digestPattern.MatchString(s) }
func ValidText(s string, n int) bool {
	if s == "" || len(s) > n || !utf8.ValidString(s) || strings.TrimSpace(s) != s {
		return false
	}
	for _, r := range s {
		if unicode.IsControl(r) {
			return false
		}
	}
	return true
}
func Canonical(v any) (json.RawMessage, string, error) {
	raw, err := json.Marshal(v)
	if err != nil {
		return nil, "", err
	}
	return approval.CanonicalPayload(raw)
}

type Profile struct {
	Schema         string `json:"schema"`
	TenantID       string `json:"tenant_id"`
	OrganizationID string `json:"organization_id"`
	Market         string `json:"market"`
	Currency       string `json:"currency"`
	Origin         string `json:"origin"`
	PolicyCode     string `json:"policy_code"`
}

func (p Profile) Valid() bool {
	if p.Schema != "elite-catalog-publication/v1" || !regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`).MatchString(p.TenantID) || !ValidID(p.OrganizationID) || !regexp.MustCompile(`^[A-Z]{2}$`).MatchString(p.Market) || !regexp.MustCompile(`^[A-Z]{3}$`).MatchString(p.Currency) || p.PolicyCode != PolicyCode {
		return false
	}
	u, e := url.Parse(p.Origin)
	return e == nil && u.Scheme == "https" && u.Host != "" && u.User == nil && u.RawQuery == "" && u.Fragment == "" && u.Path == "" && u.Opaque == "" && len(p.Origin) <= 512
}

type DraftRequest struct {
	CommandID   string           `json:"command_id"`
	PriceBookID string           `json:"price_book_id"`
	Models      []ModelReference `json:"models"`
}
type ModelReference struct {
	ModelID string `json:"model_id"`
	MediaID string `json:"media_id"`
}

func (r DraftRequest) Valid() bool {
	if !ValidID(r.CommandID) || !ValidID(r.PriceBookID) || len(r.Models) == 0 || len(r.Models) > 32 {
		return false
	}
	seen := map[string]bool{}
	for _, m := range r.Models {
		if !ValidID(m.ModelID) || !ValidID(m.MediaID) || seen[m.ModelID] {
			return false
		}
		seen[m.ModelID] = true
	}
	return true
}

type ReviewRequest struct {
	CommandID      string `json:"command_id"`
	DraftID        string `json:"draft_id"`
	Stage          string `json:"stage"`
	SnapshotSHA256 string `json:"snapshot_sha256"`
	Approved       bool   `json:"approved"`
	Reason         string `json:"reason"`
	EvidenceSHA256 string `json:"evidence_sha256"`
}

func (r ReviewRequest) Valid() bool {
	return ValidID(r.CommandID) && ValidID(r.DraftID) && ValidSHA(r.SnapshotSHA256) && ValidSHA(r.EvidenceSHA256) && ValidText(r.Reason, 2000) && map[string]bool{"legal": true, "technical": true, "media": true, "publication": true}[r.Stage]
}

type PublishRequest struct {
	CommandID          string `json:"command_id"`
	DraftID            string `json:"draft_id"`
	SnapshotSHA256     string `json:"snapshot_sha256"`
	ExpectedGeneration int64  `json:"expected_generation,string"`
	Reason             string `json:"reason"`
}

func (r PublishRequest) Valid() bool {
	return ValidID(r.CommandID) && ValidID(r.DraftID) && ValidSHA(r.SnapshotSHA256) && r.ExpectedGeneration >= 0 && ValidText(r.Reason, 2000)
}

type Receipt struct {
	SourcePriceBookID    string   `json:"source_price_book_id,omitempty"`
	SourceVariantIDs     []string `json:"source_variant_ids,omitempty"`
	EffectivePriceBookID string   `json:"effective_price_book_id,omitempty"`
	CommandID            string   `json:"command_id"`
	Actor                string   `json:"actor"`
	Kind                 string   `json:"kind"`
	ResourceID           string   `json:"resource_id"`
	RequestSHA256        string   `json:"request_sha256"`
	SnapshotSHA256       string   `json:"snapshot_sha256,omitempty"`
	Generation           int64    `json:"generation,string"`
	Replay               bool     `json:"replay"`
}
type Media struct {
	ID             string `json:"id"`
	OriginalSHA256 string `json:"original_sha256"`
	SHA256         string `json:"sha256"`
	Width          int    `json:"width"`
	Height         int    `json:"height"`
}
type Variant struct {
	ID                   string          `json:"id"`
	ModelID              string          `json:"model_id"`
	Code                 string          `json:"code"`
	DisplayName          string          `json:"display_name"`
	BatterySpecification json.RawMessage `json:"battery_specification"`
	HomologationState    string          `json:"homologation_state"`
	AmountMinorUnits     int64           `json:"amount_minor_units,string"`
	TaxMode              string          `json:"tax_mode"`
}
type Model struct {
	electromobility.Model
	Media        Media  `json:"media"`
	CanonicalURL string `json:"canonical_url"`
}
type Snapshot struct {
	TenantCode  string     `json:"tenant_code"`
	Schema      string     `json:"schema"`
	Profile     Profile    `json:"profile"`
	PriceBookID string     `json:"price_book_id"`
	ValidFrom   time.Time  `json:"valid_from"`
	ValidUntil  *time.Time `json:"valid_until"`
	Models      []Model    `json:"models"`
	Variants    []Variant  `json:"variants"`
}
type Draft struct {
	ID       string            `json:"id"`
	Maker    string            `json:"maker"`
	SHA256   string            `json:"sha256"`
	Snapshot json.RawMessage   `json:"snapshot"`
	Reviews  map[string]string `json:"reviews"`
}
type Publication struct {
	EffectivePriceBookID string          `json:"effective_price_book_id,omitempty"`
	Generation           int64           `json:"generation,string"`
	DraftID              string          `json:"draft_id"`
	SHA256               string          `json:"sha256"`
	Snapshot             json.RawMessage `json:"snapshot"`
}

type SearchHit struct {
	ModelID      string `json:"model_id"`
	Title        string `json:"title"`
	URL          string `json:"url"`
	Generation   int64  `json:"generation,string"`
	SourceSHA256 string `json:"source_sha256"`
}

// PublicDocument omits private stream ownership and review configuration.
type PublicDocument struct {
	TenantCode           string    `json:"tenant_code"`
	EffectivePriceBookID string    `json:"effective_price_book_id,omitempty"`
	Generation           int64     `json:"generation,string"`
	SourceSHA256         string    `json:"source_sha256"`
	Market               string    `json:"market"`
	Currency             string    `json:"currency"`
	Models               []Model   `json:"models"`
	Variants             []Variant `json:"variants"`
}

func ProjectPublic(p Publication) (PublicDocument, error) {
	var s Snapshot
	if err := json.Unmarshal(p.Snapshot, &s); err != nil {
		return PublicDocument{}, err
	}
	return PublicDocument{TenantCode: s.TenantCode, EffectivePriceBookID: p.EffectivePriceBookID, Generation: p.Generation, SourceSHA256: p.SHA256, Market: s.Profile.Market, Currency: s.Profile.Currency, Models: s.Models, Variants: s.Variants}, nil
}
````

### FILE: `internal/catalogrelease/feed.go`

```yaml
block_id: "GO-CONNECTED-CATALOG-PUBLICATION:file12:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "b9a4e54ca3b0e04473f9df61e4c0df9b3f075bfe30c873e9a7fa3c20df6c94ca"
variables: []
secrets_allowed: false
```

````go
package catalogrelease

// AUTHORED reference feed protocol over pinned Go net/http. Provider-specific
// marketplace SDKs remain separate. The receiver must preserve per-key receipts
// and reject an older generation after a newer one; it is not a blind webhook.
import (
	"bytes"
	"context"
	"crypto/sha256"
	"crypto/tls"
	"elite.local/enterprise/internal/approval"
	"elite.local/enterprise/internal/outbounddelivery"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type FeedProfile struct {
	Schema     string `json:"schema"`
	ReceiverID string `json:"receiver_id"`
	BaseURL    string `json:"base_url"`
	Mode       string `json:"mode"`
}

func (p FeedProfile) Valid() bool {
	if p.Schema != "elite-catalog-feed/v1" || !ValidID(p.ReceiverID) || len(p.BaseURL) > 1024 {
		return false
	}
	u, e := url.Parse(p.BaseURL)
	if e != nil || u.Host == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" || u.Path != "" || u.Opaque != "" {
		return false
	}
	if p.Mode == "fixture" {
		return u.Scheme == "http" && (u.Hostname() == "127.0.0.1" || u.Hostname() == "::1") && u.Port() != ""
	}
	return p.Mode == "https" && u.Scheme == "https"
}

type FeedAck struct {
	Status        string    `json:"status"`
	Generation    int64     `json:"generation,string"`
	PayloadSHA256 string    `json:"payload_sha256"`
	RemoteID      string    `json:"remote_id"`
	AcceptedAt    time.Time `json:"accepted_at"`
}
type FeedHTTP struct {
	profile    FeedProfile
	profileSHA string
	token      string
	client     *http.Client
	transport  *http.Transport
}

func NewFeedHTTP(p FeedProfile, token string) (*FeedHTTP, error) {
	if !p.Valid() || !ValidText(token, 4096) {
		return nil, ErrInvalid
	}
	_, sha, e := Canonical(p)
	if e != nil {
		return nil, e
	}
	transport := &http.Transport{TLSClientConfig: &tls.Config{MinVersion: tls.VersionTLS12}, ResponseHeaderTimeout: 5 * time.Second, MaxResponseHeaderBytes: 16384, DisableKeepAlives: true}
	client := &http.Client{Transport: transport, Timeout: 8 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	return &FeedHTTP{p, sha, token, client, transport}, nil
}
func (c *FeedHTTP) Profile() FeedProfile  { return c.profile }
func (c *FeedHTTP) ProfileSHA256() string { return c.profileSHA }
func (c *FeedHTTP) Close()                { c.transport.CloseIdleConnections() }
func (c *FeedHTTP) exchange(ctx context.Context, method, key string, generation int64, payload []byte, payloadSHA string) (outbounddelivery.Receipt, error) {
	if !ValidText(key, 128) || generation < 1 || !ValidSHA(payloadSHA) {
		return outbounddelivery.Receipt{}, ErrInvalid
	}
	req, e := http.NewRequestWithContext(ctx, method, c.profile.BaseURL+"/publications/"+url.PathEscape(key), bytes.NewReader(payload))
	if e != nil {
		return outbounddelivery.Receipt{}, ErrInvalid
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Idempotency-Key", key)
	req.Header.Set("X-Catalog-Generation", fmt.Sprint(generation))
	req.Header.Set("X-Content-SHA256", payloadSHA)
	response, e := c.client.Do(req)
	if e != nil {
		return outbounddelivery.Receipt{}, errors.New("feed response unconfirmed")
	}
	defer response.Body.Close()
	raw, e := io.ReadAll(io.LimitReader(response.Body, 16385))
	if e != nil || len(raw) > 16384 {
		return outbounddelivery.Receipt{}, errors.New("feed response exceeds boundary")
	}
	if _, _, e = approval.CanonicalPayload(raw); e != nil {
		return outbounddelivery.Receipt{}, errors.New("feed receipt malformed")
	}
	var ack FeedAck
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	if dec.Decode(&ack) != nil || dec.Decode(new(any)) != io.EOF || ack.Generation != generation || ack.PayloadSHA256 != payloadSHA {
		return outbounddelivery.Receipt{}, errors.New("feed receipt binding mismatch")
	}
	digest := sha256.Sum256(raw)
	evidence := hex.EncodeToString(digest[:])
	if ack.Status == "rejected" && (response.StatusCode == 422 || method == "GET" && response.StatusCode == 200) {
		return outbounddelivery.Receipt{}, outbounddelivery.NewTerminalFailure("FEED_REJECTED", evidence, nil)
	}
	if ack.Status != "accepted" || (response.StatusCode != 200 && response.StatusCode != 201) || strings.TrimSpace(ack.RemoteID) == "" {
		return outbounddelivery.Receipt{}, errors.New("feed receipt not accepted")
	}
	receipt := outbounddelivery.Receipt{ProviderMessageID: ack.RemoteID, EvidenceSHA256: evidence, AcceptedAt: ack.AcceptedAt}
	if e = receipt.Validate(); e != nil {
		return outbounddelivery.Receipt{}, e
	}
	return receipt, nil
}
func (c *FeedHTTP) Send(ctx context.Context, key string, generation int64, payload []byte, payloadSHA string) (outbounddelivery.Receipt, error) {
	sum := sha256.Sum256(payload)
	if hex.EncodeToString(sum[:]) != payloadSHA || len(payload) > 524288 {
		return outbounddelivery.Receipt{}, ErrInvalid
	}
	return c.exchange(ctx, "POST", key, generation, payload, payloadSHA)
}
func (c *FeedHTTP) Reconcile(ctx context.Context, key string, generation int64, payloadSHA string) (outbounddelivery.Receipt, error) {
	return c.exchange(ctx, "GET", key, generation, nil, payloadSHA)
}

type FeedStatus struct {
	DeliveryKey   string `json:"delivery_key"`
	Generation    int64  `json:"generation,string"`
	State         string `json:"state"`
	PayloadSHA256 string `json:"payload_sha256"`
	SourceSHA256  string `json:"source_sha256"`
}
````

### FILE: `internal/catalogrelease/media.go`

```yaml
block_id: "GO-CONNECTED-CATALOG-PUBLICATION:file13:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "8440834c8cc68d963525cfef5c396d563e17a13a4dc5b912c7fa826c6d7180d2"
variables: []
secrets_allowed: false
```

````go
package catalogrelease

// DEPENDENCY_PIN: existing Go1.26.8 standard library image/png performs PNG
// decoding/reencoding. This AUTHORED wrapper bounds resources and records hashes;
// it does not claim malware scanning or native-image-library repair.
import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"image/png"
)

func NormalizePNG(raw []byte) (Media, []byte, error) {
	if len(raw) == 0 || len(raw) > 1<<20 {
		return Media{}, nil, ErrInvalid
	}
	cfg, err := png.DecodeConfig(bytes.NewReader(raw))
	if err != nil || cfg.Width < 1 || cfg.Height < 1 || cfg.Width > 2048 || cfg.Height > 2048 || cfg.Width*cfg.Height > 1048576 {
		return Media{}, nil, ErrInvalid
	}
	im, err := png.Decode(bytes.NewReader(raw))
	if err != nil {
		return Media{}, nil, ErrInvalid
	}
	var out bytes.Buffer
	if err = png.Encode(&out, im); err != nil || out.Len() > 4<<20 {
		return Media{}, nil, ErrInvalid
	}
	original := sha256.Sum256(raw)
	clean := sha256.Sum256(out.Bytes())
	return Media{OriginalSHA256: hex.EncodeToString(original[:]), SHA256: hex.EncodeToString(clean[:]), Width: cfg.Width, Height: cfg.Height}, out.Bytes(), nil
}
````

### FILE: `internal/catalogrelease/media_test.go`

```yaml
block_id: "GO-CONNECTED-CATALOG-PUBLICATION:file14:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "9bd3a2c5dc637f3a7d9ecadc4f866adfc08411380671614465a38aa50fd8c1cc"
variables: []
secrets_allowed: false
```

````go
package catalogrelease

import (
	"bytes"
	"image"
	"image/png"
	"testing"
)

func TestCatalogMediaStructuralBoundary(t *testing.T) {
	var b bytes.Buffer
	if e := png.Encode(&b, image.NewNRGBA(image.Rect(0, 0, 4, 3))); e != nil {
		t.Fatal(e)
	}
	original := append(append([]byte(nil), b.Bytes()...), []byte("<script>trailer</script>")...)
	m, clean, e := NormalizePNG(original)
	if e != nil || m.Width != 4 || m.Height != 3 || bytes.Contains(clean, []byte("script")) || m.SHA256 == m.OriginalSHA256 {
		t.Fatal("PNG normalization", m, e)
	}
	again, out, e := NormalizePNG(clean)
	if e != nil || !bytes.Equal(out, clean) || again.SHA256 != m.SHA256 {
		t.Fatal("non-deterministic PNG")
	}
	for _, bad := range [][]byte{[]byte("<svg onload='x'/>"), b.Bytes()[:20], make([]byte, (1<<20)+1)} {
		if _, _, e := NormalizePNG(bad); e == nil {
			t.Fatal("invalid input accepted")
		}
	}
	b.Reset()
	if e = png.Encode(&b, image.NewNRGBA(image.Rect(0, 0, 2049, 1))); e != nil {
		t.Fatal(e)
	}
	if _, _, e = NormalizePNG(b.Bytes()); e == nil {
		t.Fatal("dimension budget ignored")
	}
}
````

### FILE: `internal/catalogrelease/profile.go`

```yaml
block_id: "GO-CONNECTED-CATALOG-PUBLICATION:file15:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "f9f7387bc8f01a87efe409d4ffa6090bd3feda95de3de52851a6d5545549ca86"
variables: []
secrets_allowed: false
```

````go
package catalogrelease

// AUTHORED bounded profile loading; activation never accepts an unpinned file.
import (
	"bytes"
	"crypto/sha256"
	"elite.local/enterprise/internal/approval"
	"encoding/hex"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

func LoadProfile(path, expected string) (Profile, error) {
	var p Profile
	if !filepath.IsAbs(path) || !ValidSHA(expected) {
		return p, ErrInvalid
	}
	info, e := os.Lstat(path)
	if e != nil || !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 || info.Size() < 2 || info.Size() > 16384 {
		return p, ErrInvalid
	}
	f, e := os.Open(path)
	if e != nil {
		return p, ErrInvalid
	}
	defer f.Close()
	actual, e := f.Stat()
	if e != nil || !os.SameFile(info, actual) {
		return p, ErrInvalid
	}
	raw, e := io.ReadAll(io.LimitReader(f, 16385))
	if e != nil || len(raw) > 16384 {
		return p, ErrInvalid
	}
	digest := sha256.Sum256(raw)
	if hex.EncodeToString(digest[:]) != expected {
		return p, ErrInvalid
	}
	if _, _, e = approval.CanonicalPayload(raw); e != nil {
		return p, ErrInvalid
	}
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	if dec.Decode(&p) != nil || dec.Decode(new(any)) != io.EOF || !p.Valid() {
		return Profile{}, ErrInvalid
	}
	return p, nil
}

func LoadFeedProfile(path, expected string) (FeedProfile, error) {
	var p FeedProfile
	if !filepath.IsAbs(path) || !ValidSHA(expected) {
		return p, ErrInvalid
	}
	info, e := os.Lstat(path)
	if e != nil || !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 || info.Size() < 2 || info.Size() > 16384 {
		return p, ErrInvalid
	}
	f, e := os.Open(path)
	if e != nil {
		return p, ErrInvalid
	}
	defer f.Close()
	actual, e := f.Stat()
	if e != nil || !os.SameFile(info, actual) {
		return p, ErrInvalid
	}
	raw, e := io.ReadAll(io.LimitReader(f, 16385))
	if e != nil || len(raw) > 16384 {
		return p, ErrInvalid
	}
	digest := sha256.Sum256(raw)
	if hex.EncodeToString(digest[:]) != expected {
		return p, ErrInvalid
	}
	if _, _, e = approval.CanonicalPayload(raw); e != nil {
		return p, ErrInvalid
	}
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	if dec.Decode(&p) != nil || dec.Decode(new(any)) != io.EOF || !p.Valid() {
		return FeedProfile{}, ErrInvalid
	}
	return p, nil
}
````

### FILE: `internal/catalogrelease/profile_test.go`

```yaml
block_id: "GO-CONNECTED-CATALOG-PUBLICATION:file16:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "d6c792e77f699892b937fa09af6b72551f47d8896c10b5a965372ab272e3430c"
variables: []
secrets_allowed: false
```

````go
package catalogrelease

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestCatalogProfilePinnedBoundary(t *testing.T) {
	p := Profile{Schema: "elite-catalog-publication/v1", TenantID: "ca6c9cb7-e548-4f19-a2e6-27d2aa2d4893", OrganizationID: "store", Market: "AR", Currency: "ARS", Origin: "https://catalog.example.invalid", PolicyCode: PolicyCode}
	raw, _ := json.Marshal(p)
	digest := sha256.Sum256(raw)
	hash := hex.EncodeToString(digest[:])
	path := filepath.Join(t.TempDir(), "profile.json")
	if e := os.WriteFile(path, raw, 0600); e != nil {
		t.Fatal(e)
	}
	if got, e := LoadProfile(path, hash); e != nil || got != p {
		t.Fatal("valid profile", got, e)
	}
	if _, e := LoadProfile("relative.json", hash); e == nil {
		t.Fatal("relative path")
	}
	if _, e := LoadProfile(filepath.Dir(path), hash); e == nil {
		t.Fatal("directory path")
	}
	changed := append(append([]byte(nil), raw...), '\n')
	if e := os.WriteFile(path, changed, 0600); e != nil {
		t.Fatal(e)
	}
	if _, e := LoadProfile(path, hash); e == nil {
		t.Fatal("changed bytes under old hash")
	}
	bad := append(append([]byte(nil), raw[:len(raw)-1]...), []byte(`,"secret":"not-an-input"}`)...)
	digest = sha256.Sum256(bad)
	if e := os.WriteFile(path, bad, 0600); e != nil {
		t.Fatal(e)
	}
	if _, e := LoadProfile(path, hex.EncodeToString(digest[:])); e == nil {
		t.Fatal("unknown field admitted")
	}
	p.Origin = "https://user:password@example.invalid"
	if p.Valid() {
		t.Fatal("credential URL")
	}
	p.Origin = "https://example.invalid/path"
	if p.Valid() {
		t.Fatal("origin path")
	}
}
````

### FILE: `internal/platform/httpapi/catalog_release.go`

```yaml
block_id: "GO-CONNECTED-CATALOG-PUBLICATION:file17:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "adb4c2a1f24fc886ac0d569cca88326822cc3efd44b5d0bc3b32924e78c9c3ee"
variables: []
secrets_allowed: false
```

````go
// AUTHORED HTTP binding to the existing authenticated role transport and
// connected catalog publication owner. No provider credentials are read.
package httpapi

import (
	"bytes"
	"context"
	"elite.local/enterprise/internal/approval"
	cr "elite.local/enterprise/internal/catalogrelease"
	"elite.local/enterprise/internal/platform/identity"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"io"
	"net/http"
	"strconv"
)

type CatalogReleaseService interface {
	UploadPNG(context.Context, identity.Principal, string, []byte) (cr.Receipt, error)
	CreateDraft(context.Context, identity.Principal, cr.DraftRequest) (cr.Receipt, error)
	Draft(context.Context, identity.Principal, string) (cr.Draft, error)
	Review(context.Context, identity.Principal, cr.ReviewRequest) (cr.Receipt, error)
	Publish(context.Context, identity.Principal, cr.PublishRequest) (cr.Receipt, error)
	CommandReceipt(context.Context, identity.Principal, string) (cr.Receipt, error)
	Public(context.Context) (cr.Publication, error)
	PublicMedia(context.Context, string) ([]byte, error)
	PublicSearch(context.Context, string) ([]cr.SearchHit, error)
}
type CatalogCurrentService interface {
	Current(context.Context, identity.Principal) (cr.Publication, error)
}
type CatalogSourceService interface {
	CreateSource(context.Context, identity.Principal, cr.SourceRequest) (cr.Receipt, error)
}
type CatalogFeedService interface {
	Send(context.Context, identity.Principal, int64) (cr.FeedStatus, error)
	Status(context.Context, identity.Principal, int64) (cr.FeedStatus, error)
	Reconcile(context.Context, identity.Principal, int64) (cr.FeedStatus, error)
}
type CatalogReleaseModule struct {
	Service CatalogReleaseService
	Feed    CatalogFeedService
}

func catalogError(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}
	var pg *pgconn.PgError
	switch {
	case errors.Is(err, cr.ErrInvalid):
		writeProblem(w, 400, "CATALOG_INVALID", "invalid bounded catalog request")
	case errors.Is(err, cr.ErrNotFound):
		writeProblem(w, 404, "CATALOG_NOT_FOUND", "catalog resource not available in this scope")
	case errors.Is(err, cr.ErrConflict), errors.Is(err, approval.ErrNotPending), errors.Is(err, approval.ErrSeparation):
		writeProblem(w, 409, "CATALOG_CONFLICT", "consult the immutable command and current publication")
	case errors.As(err, &pg) && (pg.Code == "23505" || pg.Code == "23514" || pg.Code == "P0001" || pg.Code == "40001" || pg.Code == "40P01"):
		writeProblem(w, 409, "CATALOG_CONFLICT", "consult the immutable command and current publication")
	default:
		writeProblem(w, 503, "CATALOG_UNCONFIRMED", "result unconfirmed; consult the immutable command receipt")
	}
	return true
}
func catalogDecode(w http.ResponseWriter, r *http.Request, v any) bool {
	raw, e := io.ReadAll(http.MaxBytesReader(w, r.Body, 32768))
	if e != nil {
		writeProblem(w, 413, "CATALOG_TOO_LARGE", "bounded command required")
		return false
	}
	if _, _, e = approval.CanonicalPayload(raw); e != nil {
		writeProblem(w, 400, "CATALOG_JSON", "unambiguous JSON required")
		return false
	}
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	if dec.Decode(v) != nil || dec.Decode(new(any)) != io.EOF {
		writeProblem(w, 400, "CATALOG_JSON", "typed command required without extra fields")
		return false
	}
	return true
}
func (m CatalogReleaseModule) Register(mux *http.ServeMux, verifier identity.Verifier) {
	if m.Service == nil {
		return
	}
	api := franchiseJourneyAPI{verifier: verifier}
	protect := func(w http.ResponseWriter, r *http.Request, permission string) (identity.Principal, bool) {
		w.Header().Set("Cache-Control", "no-store")
		q := r.URL.Query()
		org := q.Get("organization_id")
		if len(q) != 1 || len(q["organization_id"]) != 1 || !cr.ValidID(org) {
			writeProblem(w, 400, "CATALOG_SCOPE", "one explicit organization required")
			return identity.Principal{}, false
		}
		p, ok := api.protected(w, r, permission, org)
		if !ok {
			return p, false
		}
		p.Permissions = map[string]struct{}{permission: {}}
		p.Organizations = map[string]struct{}{org: {}}
		return p, true
	}
	reply := func(w http.ResponseWriter, v cr.Receipt, e error) {
		if catalogError(w, e) {
			return
		}
		status := 201
		if v.Replay {
			status = 200
		}
		writeJSON(w, status, v)
	}
	mux.HandleFunc("POST /v1/admin/catalog/media/{command}", func(w http.ResponseWriter, r *http.Request) {
		p, ok := protect(w, r, "catalog:draft")
		if !ok {
			return
		}
		if r.Header.Get("Content-Type") != "image/png" {
			writeProblem(w, 415, "CATALOG_PNG", "PNG required")
			return
		}
		raw, e := io.ReadAll(http.MaxBytesReader(w, r.Body, 1<<20))
		if e != nil {
			writeProblem(w, 413, "CATALOG_TOO_LARGE", "PNG exceeds limit")
			return
		}
		v, e := m.Service.UploadPNG(r.Context(), p, r.PathValue("command"), raw)
		reply(w, v, e)
	})
	if current, ok := m.Service.(CatalogCurrentService); ok {
		mux.HandleFunc("GET /v1/admin/catalog/current", func(w http.ResponseWriter, r *http.Request) {
			p, ok := protect(w, r, "catalog:read")
			if !ok {
				return
			}
			v, e := current.Current(r.Context(), p)
			if catalogError(w, e) {
				return
			}
			value, e := cr.ProjectPublic(v)
			if !catalogError(w, e) {
				writeJSON(w, 200, value)
			}
		})
	}
	if source, ok := m.Service.(CatalogSourceService); ok {
		mux.HandleFunc("POST /v1/admin/catalog/sources", func(w http.ResponseWriter, r *http.Request) {
			p, ok := protect(w, r, "catalog:draft")
			if !ok {
				return
			}
			var in cr.SourceRequest
			if !catalogDecode(w, r, &in) {
				return
			}
			v, e := source.CreateSource(r.Context(), p, in)
			reply(w, v, e)
		})
	}
	mux.HandleFunc("POST /v1/admin/catalog/drafts", func(w http.ResponseWriter, r *http.Request) {
		p, ok := protect(w, r, "catalog:draft")
		if !ok {
			return
		}
		var in cr.DraftRequest
		if !catalogDecode(w, r, &in) {
			return
		}
		v, e := m.Service.CreateDraft(r.Context(), p, in)
		reply(w, v, e)
	})
	mux.HandleFunc("GET /v1/admin/catalog/drafts/{id}", func(w http.ResponseWriter, r *http.Request) {
		p, ok := protect(w, r, "catalog:read")
		if !ok {
			return
		}
		v, e := m.Service.Draft(r.Context(), p, r.PathValue("id"))
		if !catalogError(w, e) {
			writeJSON(w, 200, v)
		}
	})
	mux.HandleFunc("GET /v1/admin/catalog/commands/{id}", func(w http.ResponseWriter, r *http.Request) {
		p, ok := protect(w, r, "catalog:read")
		if !ok {
			return
		}
		v, e := m.Service.CommandReceipt(r.Context(), p, r.PathValue("id"))
		if !catalogError(w, e) {
			writeJSON(w, 200, v)
		}
	})
	for _, stage := range []string{"legal", "technical", "media", "publication"} {
		mux.HandleFunc("POST /v1/admin/catalog/drafts/{id}/review/"+stage, func(w http.ResponseWriter, r *http.Request) {
			p, ok := protect(w, r, "catalog:review:"+stage)
			if !ok {
				return
			}
			var in cr.ReviewRequest
			if !catalogDecode(w, r, &in) {
				return
			}
			if in.DraftID != r.PathValue("id") || in.Stage != stage {
				catalogError(w, cr.ErrInvalid)
				return
			}
			v, e := m.Service.Review(r.Context(), p, in)
			reply(w, v, e)
		})
	}
	mux.HandleFunc("POST /v1/admin/catalog/drafts/{id}/publish", func(w http.ResponseWriter, r *http.Request) {
		p, ok := protect(w, r, "catalog:publish")
		if !ok {
			return
		}
		var in cr.PublishRequest
		if !catalogDecode(w, r, &in) {
			return
		}
		if in.DraftID != r.PathValue("id") {
			catalogError(w, cr.ErrInvalid)
			return
		}
		v, e := m.Service.Publish(r.Context(), p, in)
		reply(w, v, e)
	})
	public := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=0, must-revalidate")
		if len(r.URL.Query()) != 0 {
			catalogError(w, cr.ErrInvalid)
			return
		}
		v, e := m.Service.Public(r.Context())
		if catalogError(w, e) {
			return
		}
		etag := fmt.Sprintf(`"catalog-%d-%s"`, v.Generation, v.SHA256)
		w.Header().Set("ETag", etag)
		if r.Header.Get("If-None-Match") == etag {
			w.WriteHeader(304)
			return
		}
		value, e := cr.ProjectPublic(v)
		if !catalogError(w, e) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(value)
		}
	}
	mux.HandleFunc("GET /v1/public/catalog", public)
	mux.HandleFunc("GET /v1/public/catalog/feed", public)
	mux.HandleFunc("GET /v1/public/catalog/media/{sha}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=0, must-revalidate")
		if len(r.URL.Query()) != 0 {
			catalogError(w, cr.ErrInvalid)
			return
		}
		raw, e := m.Service.PublicMedia(r.Context(), r.PathValue("sha"))
		if catalogError(w, e) {
			return
		}
		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		etag := `"` + r.PathValue("sha") + `"`
		w.Header().Set("ETag", etag)
		if r.Header.Get("If-None-Match") == etag {
			w.WriteHeader(304)
			return
		}
		_, _ = w.Write(raw)
	})
	mux.HandleFunc("GET /v1/public/catalog/search", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		q := r.URL.Query()
		if len(q) != 1 || len(q["q"]) != 1 {
			catalogError(w, cr.ErrInvalid)
			return
		}
		rows, e := m.Service.PublicSearch(r.Context(), q.Get("q"))
		if !catalogError(w, e) {
			writeJSON(w, 200, rows)
		}
	})
	if m.Feed != nil {
		for _, route := range []struct{ method, suffix, permission string }{{"POST", "", "catalog:feed"}, {"GET", "", "catalog:read"}, {"POST", "/reconcile", "catalog:feed"}} {
			mux.HandleFunc(route.method+" /v1/admin/catalog/feed/{generation}"+route.suffix, func(w http.ResponseWriter, r *http.Request) {
				p, ok := protect(w, r, route.permission)
				if !ok {
					return
				}
				raw := r.PathValue("generation")
				generation, e := strconv.ParseInt(raw, 10, 64)
				if e != nil || generation < 1 || strconv.FormatInt(generation, 10) != raw {
					catalogError(w, cr.ErrInvalid)
					return
				}
				if r.Body != nil {
					value, e := io.ReadAll(http.MaxBytesReader(w, r.Body, 1))
					if e != nil || len(value) != 0 {
						catalogError(w, cr.ErrInvalid)
						return
					}
				}
				var status cr.FeedStatus
				switch {
				case route.method == "GET":
					status, e = m.Feed.Status(r.Context(), p, generation)
				case route.suffix == "/reconcile":
					status, e = m.Feed.Reconcile(r.Context(), p, generation)
				default:
					status, e = m.Feed.Send(r.Context(), p, generation)
				}
				if !catalogError(w, e) {
					writeJSON(w, 200, status)
				}
			})
		}
	}
}
````

### FILE: `internal/platform/httpapi/catalog_release_test.go`

```yaml
block_id: "GO-CONNECTED-CATALOG-PUBLICATION:file18:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "ad30b8cf5415640c48f90d9436aec0fbd9bc1b2777966c9bb1d9ab8b0b0dbf9e"
variables: []
secrets_allowed: false
```

````go
package httpapi
import(
 "context"
 "encoding/json"
 "net/http"
 "net/http/httptest"
 "testing"
 cr "elite.local/enterprise/internal/catalogrelease"
)
type catalogCacheProbe struct{CatalogReleaseService}
func(catalogCacheProbe)Public(context.Context)(cr.Publication,error){
 raw,_:=json.Marshal(cr.Snapshot{Profile:cr.Profile{Market:"AR",Currency:"ARS"}})
 return cr.Publication{Generation:1,SHA256:"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",Snapshot:raw},nil
}
func TestCatalogPublicCachePolicy(t *testing.T){
 mux:=http.NewServeMux();CatalogReleaseModule{Service:catalogCacheProbe{}}.Register(mux,nil)
 first:=httptest.NewRecorder();mux.ServeHTTP(first,httptest.NewRequest("GET","/v1/public/catalog",nil))
 if first.Code!=200||first.Header().Get("Cache-Control")!="public, max-age=0, must-revalidate"{t.Fatal(first.Code,first.Header())}
 etag:=first.Header().Get("ETag");if etag==""{t.Fatal("no ETag")}
 req:=httptest.NewRequest("GET","/v1/public/catalog",nil);req.Header.Set("If-None-Match",etag)
 next:=httptest.NewRecorder();mux.ServeHTTP(next,req)
 if next.Code!=304||next.Body.Len()!=0||next.Header().Get("Cache-Control")!="public, max-age=0, must-revalidate"{t.Fatal(next.Code,next.Header())}
}
````

### FILE: `internal/platform/postgres/catalog_feed.go`

```yaml
block_id: "GO-CONNECTED-CATALOG-PUBLICATION:file19:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "d81857dd9c5ca99e3210866392a4b330cd8d68f086663253da23476a35f57a81"
variables: []
secrets_allowed: false
```

````go
package postgres

// AUTHORED feed binding to the existing outbound delivery state machine/fence.
import (
	"context"
	"crypto/sha256"
	cr "elite.local/enterprise/internal/catalogrelease"
	"elite.local/enterprise/internal/channels"
	"elite.local/enterprise/internal/outbounddelivery"
	"elite.local/enterprise/internal/platform/identity"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"time"
)

type CatalogFeed struct {
	catalog  *CatalogRelease
	fence    *OutboundDeliveryStore
	receiver *cr.FeedHTTP
}

func NewCatalogFeed(catalog *CatalogRelease, hmacKey []byte, receiver *cr.FeedHTTP) (*CatalogFeed, error) {
	if catalog == nil || receiver == nil {
		return nil, cr.ErrInvalid
	}
	fence, e := NewOutboundDeliveryStore(catalog.pool, hmacKey, 15*time.Second)
	if e != nil {
		return nil, e
	}
	return &CatalogFeed{catalog, fence, receiver}, nil
}

type catalogFeedPlan struct {
	message               channels.Message
	payload               []byte
	payloadSHA, sourceSHA string
	generation            int64
}

func (s *CatalogFeed) plan(ctx context.Context, p identity.Principal, generation int64, permission string) (catalogFeedPlan, error) {
	var plan catalogFeedPlan
	if generation < 1 || !s.catalog.authorized(p, permission) {
		return plan, cr.ErrNotFound
	}
	var publication cr.Publication
	e := s.catalog.pool.QueryRow(ctx, `select p.generation,p.draft_id,d.snapshot_sha256,d.snapshot_canonical,p.effective_price_book_id
 from catalog.release_publication p join catalog.release_draft d using(tenant_id,draft_id)
 where p.tenant_id=$1 and p.generation=$2`, p.TenantID, generation).Scan(&publication.Generation, &publication.DraftID, &publication.SHA256, &publication.Snapshot, &publication.EffectivePriceBookID)
	if errors.Is(e, pgx.ErrNoRows) {
		return plan, cr.ErrNotFound
	}
	if e != nil {
		return plan, e
	}
	doc, e := cr.ProjectPublic(publication)
	if e != nil {
		return plan, e
	}
	raw, e := json.Marshal(doc)
	if e != nil {
		return plan, e
	}
	raw = append(raw, '\n')
	digest := sha256.Sum256(raw)
	plan = catalogFeedPlan{payload: raw, payloadSHA: hex.EncodeToString(digest[:]), sourceSHA: publication.SHA256, generation: generation}
	receiver := s.receiver.Profile()
	metadata, _, e := cr.Canonical(map[string]string{"source_sha256": plan.sourceSHA, "payload_sha256": plan.payloadSHA, "receiver_profile_sha256": s.receiver.ProfileSHA256()})
	if e != nil {
		return plan, e
	}
	plan.message = channels.Message{ChannelCode: "catalog_feed", TenantID: p.TenantID, ExternalID: receiver.ReceiverID, ThreadID: fmt.Sprint(generation), DeliveryKey: fmt.Sprintf("catalog-feed:%d:%s", generation, receiver.ReceiverID), Direction: channels.DirectionOut, Text: string(metadata)}
	if e = plan.message.Validate(); e != nil {
		return plan, e
	}
	return plan, nil
}

type catalogFeedBoundStore struct {
	service *CatalogFeed
	actor   identity.Principal
	plan    catalogFeedPlan
}

func (s *catalogFeedBoundStore) checkIntent(ctx context.Context) error {
	var actor, profile string
	e := s.service.catalog.pool.QueryRow(ctx, `select actor,receiver_profile_sha256 from catalog.release_feed_intent where tenant_id=$1 and channel_code='catalog_feed' and delivery_key=$2`, s.actor.TenantID, s.plan.message.DeliveryKey).Scan(&actor, &profile)
	if errors.Is(e, pgx.ErrNoRows) {
		return nil
	}
	if e != nil {
		return e
	}
	if actor != s.actor.Subject || profile != s.service.receiver.ProfileSHA256() {
		return cr.ErrConflict
	}
	return nil
}
func (s *catalogFeedBoundStore) Claim(ctx context.Context, m channels.Message, h string) (outbounddelivery.Claim, error) {
	expected, e := outbounddelivery.MessageSHA256(s.plan.message)
	if e != nil || expected != h {
		return outbounddelivery.Claim{}, cr.ErrInvalid
	}
	if e = s.checkIntent(ctx); e != nil {
		return outbounddelivery.Claim{}, e
	}
	claim, e := s.service.fence.claimWithAdmission(ctx, m, h, func(ctx context.Context, tx pgx.Tx) error {
		var current bool
		e := tx.QueryRow(ctx, `select exists(select 1 from catalog.release_public_current where tenant_id=$1 and generation=$2 and snapshot_sha256=$3)`, m.TenantID, s.plan.generation, s.plan.sourceSHA).Scan(&current)
		if e != nil {
			return e
		}
		if !current {
			return cr.ErrConflict
		}
		_, e = tx.Exec(ctx, `insert into catalog.release_feed_intent(tenant_id,channel_code,delivery_key,receiver_id,receiver_profile_sha256,generation,actor,request_sha256,payload_sha256,source_sha256,payload)
   values($1,'catalog_feed',$2,$3,$4,$5,$6,$7,$8,$9,$10)`, m.TenantID, m.DeliveryKey, m.ExternalID, s.service.receiver.ProfileSHA256(), s.plan.generation, s.actor.Subject, h, s.plan.payloadSHA, s.plan.sourceSHA, s.plan.payload)
		return e
	})
	if e != nil {
		return claim, e
	}
	return claim, s.checkIntent(ctx)
}
func (s *catalogFeedBoundStore) Complete(c context.Context, m channels.Message, h string, r outbounddelivery.Receipt) error {
	return s.service.fence.Complete(c, m, h, r)
}
func (s *catalogFeedBoundStore) MarkUnknown(c context.Context, m channels.Message, h, code string) error {
	return s.service.fence.MarkUnknown(c, m, h, code)
}
func (s *catalogFeedBoundStore) MarkFailed(c context.Context, m channels.Message, h, sha, code string) error {
	return s.service.fence.MarkFailed(c, m, h, sha, code)
}

type catalogFeedSender struct {
	receiver *cr.FeedHTTP
	plan     catalogFeedPlan
}

func (s catalogFeedSender) SendWithReceipt(ctx context.Context, m channels.Message) (outbounddelivery.Receipt, error) {
	got, e := outbounddelivery.MessageSHA256(m)
	expected, x := outbounddelivery.MessageSHA256(s.plan.message)
	if e != nil || x != nil || got != expected {
		return outbounddelivery.Receipt{}, cr.ErrInvalid
	}
	return s.receiver.Send(ctx, m.DeliveryKey, s.plan.generation, s.plan.payload, s.plan.payloadSHA)
}

type catalogOutboundOnly struct{}

func (catalogOutboundOnly) Code() string { return "catalog_feed" }
func (catalogOutboundOnly) Receive(context.Context) ([]channels.Message, error) {
	return nil, outbounddelivery.ErrInvalid
}
func (catalogOutboundOnly) Send(context.Context, channels.Message) error {
	return outbounddelivery.ErrInvalid
}
func (s *CatalogFeed) Send(ctx context.Context, p identity.Principal, generation int64) (cr.FeedStatus, error) {
	plan, e := s.plan(ctx, p, generation, "catalog:feed")
	if e != nil {
		return cr.FeedStatus{}, e
	}
	bound := &catalogFeedBoundStore{s, p, plan}
	channel := outbounddelivery.Channel{CodeValue: "catalog_feed", Receiver: catalogOutboundOnly{}, Sender: catalogFeedSender{s.receiver, plan}, Store: bound}
	e = channel.Send(ctx, plan.message)
	status, readErr := s.readStatus(ctx, p, plan)
	if readErr != nil {
		return status, errors.Join(e, readErr)
	}
	return status, e
}
func (s *CatalogFeed) readStatus(ctx context.Context, p identity.Principal, plan catalogFeedPlan) (cr.FeedStatus, error) {
	var out cr.FeedStatus
	e := s.catalog.pool.QueryRow(ctx, `select i.delivery_key,i.generation,d.state,i.payload_sha256,i.source_sha256
 from catalog.release_feed_intent i join communication.outbound_delivery d using(tenant_id,channel_code,delivery_key)
 where i.tenant_id=$1 and i.delivery_key=$2 and i.channel_code='catalog_feed' and i.actor=$3 and i.receiver_profile_sha256=$4`,
		p.TenantID, plan.message.DeliveryKey, p.Subject, s.receiver.ProfileSHA256()).Scan(&out.DeliveryKey, &out.Generation, &out.State, &out.PayloadSHA256, &out.SourceSHA256)
	if errors.Is(e, pgx.ErrNoRows) {
		return out, cr.ErrNotFound
	}
	return out, e
}
func (s *CatalogFeed) Status(ctx context.Context, p identity.Principal, generation int64) (cr.FeedStatus, error) {
	plan, e := s.plan(ctx, p, generation, "catalog:read")
	if e != nil {
		return cr.FeedStatus{}, e
	}
	return s.readStatus(ctx, p, plan)
}
func (s *CatalogFeed) Reconcile(ctx context.Context, p identity.Principal, generation int64) (cr.FeedStatus, error) {
	plan, e := s.plan(ctx, p, generation, "catalog:feed")
	if e != nil {
		return cr.FeedStatus{}, e
	}
	status, e := s.readStatus(ctx, p, plan)
	if e != nil {
		return status, e
	}
	if status.State == "accepted" || status.State == "failed_terminal" {
		return status, nil
	}
	if status.State != "unknown" {
		return status, outbounddelivery.ErrInProgress
	}
	receipt, e := s.receiver.Reconcile(ctx, plan.message.DeliveryKey, generation, plan.payloadSHA)
	hash, hashErr := outbounddelivery.MessageSHA256(plan.message)
	if hashErr != nil {
		return status, hashErr
	}
	if e != nil {
		var terminal *outbounddelivery.TerminalFailure
		if !errors.As(e, &terminal) {
			return status, e
		}
		if e = s.fence.ReconcileFailed(ctx, plan.message, hash, terminal.EvidenceSHA256, terminal.Code); e != nil {
			return status, e
		}
	} else if e = s.fence.ReconcileAccepted(ctx, plan.message, hash, receipt); e != nil {
		return status, e
	}
	return s.readStatus(ctx, p, plan)
}
````

### FILE: `internal/platform/postgres/catalog_feed_integration_test.go`

```yaml
block_id: "GO-CONNECTED-CATALOG-PUBLICATION:file20:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "1cea884eb30c70e7891ea3b0e35d4b7d3200abaa6fe8c00c7122ee65ea706ab8"
variables: []
secrets_allowed: false
```

````go
package postgres_test

// AUTHORED receiver protocol fixture, not a live Google/Meta/ML API claim.
// Provider-specific SDK mapping remains the explicit T2805 adapter owner.
import (
	"bytes"
	"context"
	"crypto/sha256"
	cr "elite.local/enterprise/internal/catalogrelease"
	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/platform/identity"
	db "elite.local/enterprise/internal/platform/postgres"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

func attachCatalogFeedFixture(t *testing.T, client *catalogHTTPFixture, store *db.CatalogRelease) httpapi.CatalogFeedService {
	t.Helper()
	var mu sync.Mutex
	posts := map[int64]int{}
	gets := map[int64]int{}
	expected := map[int64][]byte{}
	acks := map[string]cr.FeedAck{}
	last := int64(0)
	receiver := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer fixture-feed-token" {
			w.WriteHeader(401)
			return
		}
		key := strings.TrimPrefix(r.URL.Path, "/publications/")
		generation, e := strconv.ParseInt(r.Header.Get("X-Catalog-Generation"), 10, 64)
		if e != nil || key == "" || key == r.URL.Path {
			w.WriteHeader(400)
			return
		}
		mu.Lock()
		defer mu.Unlock()
		if r.Method == "GET" {
			gets[generation]++
			ack, ok := acks[key]
			if !ok {
				w.WriteHeader(404)
				return
			}
			_ = json.NewEncoder(w).Encode(ack)
			return
		}
		if r.Method != "POST" {
			w.WriteHeader(405)
			return
		}
		raw, e := io.ReadAll(io.LimitReader(r.Body, 524289))
		if e != nil || len(raw) > 524288 {
			w.WriteHeader(413)
			return
		}
		sum := sha256.Sum256(raw)
		hash := hex.EncodeToString(sum[:])
		if !bytes.Equal(raw, expected[generation]) || r.Header.Get("X-Content-SHA256") != hash || r.Header.Get("Idempotency-Key") != key {
			t.Error("feed bytes/header differ from published source")
			w.WriteHeader(400)
			return
		}
		posts[generation]++
		if posts[generation] != 1 {
			t.Error("duplicate provider POST", generation)
		}
		ack := cr.FeedAck{Generation: generation, PayloadSHA256: hash, RemoteID: fmt.Sprintf("fixture-catalog-%d", generation), AcceptedAt: time.Now().UTC()}
		if generation == 1 || generation < last {
			ack.Status = "rejected"
			acks[key] = ack
			w.WriteHeader(422)
			_ = json.NewEncoder(w).Encode(ack)
			return
		}
		ack.Status = "accepted"
		last = generation
		acks[key] = ack
		if generation == 2 {
			conn, _, e := w.(http.Hijacker).Hijack()
			if e != nil {
				t.Error(e)
				return
			}
			conn.Close()
			return
		}
		w.WriteHeader(201)
		_ = json.NewEncoder(w).Encode(ack)
	}))
	t.Cleanup(receiver.Close)
	sender, e := cr.NewFeedHTTP(cr.FeedProfile{Schema: "elite-catalog-feed/v1", ReceiverID: "fixture-receiver", BaseURL: receiver.URL, Mode: "fixture"}, "fixture-feed-token")
	if e != nil {
		t.Fatal(e)
	}
	t.Cleanup(sender.Close)
	feeder, e := db.NewCatalogFeed(store, bytes.Repeat([]byte{'f'}, 32), sender)
	if e != nil {
		t.Fatal(e)
	}
	client.feedProbe = func(ctx context.Context, p identity.Principal, generation int64, raw []byte) {
		t.Helper()
		mu.Lock()
		expected[generation] = append([]byte(nil), raw...)
		mu.Unlock()
		path := fmt.Sprintf("/v1/admin/catalog/feed/%d", generation)
		invoke := func(method, suffix string) (cr.FeedStatus, error) {
			var out cr.FeedStatus
			data, _, _, e := client.request(ctx, &p, method, path+suffix, "", nil, "")
			if e == nil {
				e = json.Unmarshal(data, &out)
			}
			return out, e
		}
		_, e := invoke("POST", "")
		if generation == 1 || generation == 2 {
			if e == nil {
				t.Fatal("rejected/unknown feed reported accepted", generation)
			}
		} else if e != nil {
			t.Fatal("feed", e)
		}
		status, e := invoke("GET", "")
		if e != nil {
			t.Fatal(e)
		}
		wanted := "accepted"
		if generation == 1 {
			wanted = "failed_terminal"
		}
		if generation == 2 {
			wanted = "unknown"
		}
		digest := sha256.Sum256(raw)
		if status.State != wanted || status.Generation != generation || status.PayloadSHA256 != hex.EncodeToString(digest[:]) {
			t.Fatal("feed durable status", status, wanted)
		}
		// Repeating a known terminal or uncertain send cannot issue another POST.
		_, again := invoke("POST", "")
		if (generation == 1 || generation == 2) && again == nil {
			t.Fatal("terminal/unknown retry not fenced")
		}
		if generation == 2 {
			status, e = invoke("POST", "/reconcile")
			if e != nil || status.State != "accepted" {
				t.Fatal("GET reconciliation", status, e)
			}
			if status, e = invoke("POST", "/reconcile"); e != nil || status.State != "accepted" {
				t.Fatal("reconciled replay", e)
			}
		}
	}
	t.Cleanup(func() {
		mu.Lock()
		defer mu.Unlock()
		if len(posts) != 4 || posts[1] != 1 || posts[2] != 1 || posts[3] != 1 || posts[5] != 1 || gets[2] != 1 || last != 5 {
			t.Error("feed coverage", posts, gets, last)
		}
		t.Log("CATALOG_FEED_FENCE_PASS four current publications; exact storefront bytes; rejection terminal; accepted response lost then GET reconciled; no repeated provider POST; retained PG intent/fence")
	})
	return feeder
}
````

### FILE: `internal/platform/postgres/catalog_public_owner.go`

```yaml
block_id: "GO-CONNECTED-CATALOG-PUBLICATION:file21:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "c7a65ece2104306300c59b07968e7f7f015850b88698530b6948f501b0e6c53d"
variables: []
secrets_allowed: false
```

````go
// AUTHORED optional projection binding. Full composition requires migration73;
// narrow profiles retain the unchanged original public model query.
package postgres

func init() {
	electromobilityPublicModelsSQL = `select m.model_id,m.model_code,m.display_name,m.vehicle_class,m.specification from catalog.release_public_models m join platform.tenant t on t.tenant_id=m.tenant_id where t.tenant_code=$1 and t.status='active' order by m.display_name,m.model_id`
}
````

### FILE: `internal/platform/postgres/catalog_release.go`

```yaml
block_id: "GO-CONNECTED-CATALOG-PUBLICATION:file22:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "89eab5ebde780e4a5f535922e019eed5b7dd11c8d891633748e38478b634e2bb"
variables: []
secrets_allowed: false
```

````go
// AUTHORED transaction/identity binding for existing catalog/commerce/approval.
package postgres

import (
	"context"
	cr "elite.local/enterprise/internal/catalogrelease"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/randomid"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CatalogRelease struct {
	pool    *pgxpool.Pool
	profile cr.Profile
	price   *Commerce
}

func NewCatalogRelease(pool *pgxpool.Pool, profile cr.Profile, price *Commerce) (*CatalogRelease, error) {
	if pool == nil || !profile.Valid() || price == nil || price.pool != pool {
		return nil, cr.ErrInvalid
	}
	return &CatalogRelease{pool, profile, price}, nil
}
func (s *CatalogRelease) authorized(p identity.Principal, permission string) bool {
	return s != nil && approvalPrincipal(p, s.profile.TenantID, s.profile.OrganizationID, permission)
}
func (s *CatalogRelease) begin(ctx context.Context, p identity.Principal, permission string) (pgx.Tx, error) {
	if !s.authorized(p, permission) {
		return nil, cr.ErrNotFound
	}
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	ok := false
	defer func() {
		if !ok {
			_ = tx.Rollback(ctx)
		}
	}()
	if _, err = tx.Exec(ctx, `select pg_advisory_xact_lock(hashtextextended($1,0))`, "catalog-release:"+s.profile.TenantID); err != nil {
		return nil, err
	}
	raw, hash, err := cr.Canonical(s.profile)
	if err != nil {
		return nil, err
	}
	var exists bool
	err = tx.QueryRow(ctx, `select exists(select 1 from catalog.release_stream where tenant_id=$1)`, s.profile.TenantID).Scan(&exists)
	if err != nil {
		return nil, err
	}
	if !exists {
		var active bool
		err = tx.QueryRow(ctx, `select exists(select 1 from pricing.price_book where tenant_id=$1 and market=$2 and currency=$3 and status='active')`, s.profile.TenantID, s.profile.Market, s.profile.Currency).Scan(&active)
		if err != nil {
			return nil, err
		}
		if active {
			return nil, cr.ErrConflict
		}
		_, err = tx.Exec(ctx, `insert into catalog.release_stream(tenant_id,organization_id,market,currency,profile,profile_canonical,profile_sha256)
   select $1,$2,$3,$4,$5,$6,$7 where exists(select 1 from org.organization where tenant_id=$1 and organization_id=$2 and status='active')`,
			s.profile.TenantID, s.profile.OrganizationID, s.profile.Market, s.profile.Currency, raw, string(raw), hash)
		if err != nil {
			return nil, err
		}
	}
	var same bool
	err = tx.QueryRow(ctx, `select s.profile_sha256=$2 and o.status='active' and t.status='active' from catalog.release_stream s
  join org.organization o using(tenant_id,organization_id) join platform.tenant t using(tenant_id) where s.tenant_id=$1`, s.profile.TenantID, hash).Scan(&same)
	if errors.Is(err, pgx.ErrNoRows) || err == nil && !same {
		return nil, cr.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	ok = true
	return tx, nil
}
func (s *CatalogRelease) replay(ctx context.Context, tx pgx.Tx, p identity.Principal, key, kind, hash string) (cr.Receipt, bool, error) {
	var raw []byte
	var actor, oldKind, oldHash string
	err := tx.QueryRow(ctx, `select actor,kind,request_sha256,receipt from catalog.release_command where tenant_id=$1 and command_id=$2`, p.TenantID, key).Scan(&actor, &oldKind, &oldHash, &raw)
	if errors.Is(err, pgx.ErrNoRows) {
		return cr.Receipt{}, false, nil
	}
	if err != nil {
		return cr.Receipt{}, false, err
	}
	if actor != p.Subject || oldKind != kind || oldHash != hash {
		return cr.Receipt{}, false, cr.ErrConflict
	}
	var r cr.Receipt
	err = json.Unmarshal(raw, &r)
	r.Replay = true
	return r, true, err
}
func (s *CatalogRelease) save(ctx context.Context, tx pgx.Tx, p identity.Principal, r cr.Receipt) error {
	raw, _, err := cr.Canonical(r)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `insert into catalog.release_command(tenant_id,command_id,actor,kind,resource_id,request_sha256,receipt) values($1,$2,$3,$4,$5,$6,$7)`,
		p.TenantID, r.CommandID, p.Subject, r.Kind, r.ResourceID, r.RequestSHA256, raw)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)
 values($1,$2,'catalog-release-command',$3,$4,$5,1,clock_timestamp(),$6)`, p.TenantID, (randomid.Generator{}).New(), r.CommandID, 1, "catalog-release."+r.Kind, raw)
	return err
}
func (s *CatalogRelease) CommandReceipt(ctx context.Context, p identity.Principal, key string) (cr.Receipt, error) {
	if !s.authorized(p, "catalog:read") || !cr.ValidID(key) {
		return cr.Receipt{}, cr.ErrNotFound
	}
	var raw []byte
	err := s.pool.QueryRow(ctx, `select receipt from catalog.release_command where tenant_id=$1 and command_id=$2 and actor=$3`, p.TenantID, key, p.Subject).Scan(&raw)
	if errors.Is(err, pgx.ErrNoRows) {
		return cr.Receipt{}, cr.ErrNotFound
	}
	if err != nil {
		return cr.Receipt{}, err
	}
	var r cr.Receipt
	err = json.Unmarshal(raw, &r)
	return r, err
}
````

### FILE: `internal/platform/postgres/catalog_release_draft.go`

```yaml
block_id: "GO-CONNECTED-CATALOG-PUBLICATION:file23:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "393c9892f594859d481446e77d44a1fa766984359526748f2516e2d9b05a5002"
variables: []
secrets_allowed: false
```

````go
package postgres

import (
	"context"
	"elite.local/enterprise/internal/approval"
	cr "elite.local/enterprise/internal/catalogrelease"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/randomid"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5"
	"net/url"
	"sort"
)

func (s *CatalogRelease) UploadPNG(ctx context.Context, p identity.Principal, key string, raw []byte) (cr.Receipt, error) {
	if !cr.ValidID(key) || !s.authorized(p, "catalog:draft") {
		return cr.Receipt{}, cr.ErrInvalid
	}
	media, clean, err := cr.NormalizePNG(raw)
	if err != nil {
		return cr.Receipt{}, err
	}
	_, hash, err := cr.Canonical(map[string]string{"command_id": key, "original_sha256": media.OriginalSHA256})
	if err != nil {
		return cr.Receipt{}, err
	}
	tx, err := s.begin(ctx, p, "catalog:draft")
	if err != nil {
		return cr.Receipt{}, err
	}
	defer tx.Rollback(ctx)
	if r, yes, err := s.replay(ctx, tx, p, key, "media", hash); err != nil || yes {
		return r, err
	}
	media.ID = (randomid.Generator{}).New()
	_, err = tx.Exec(ctx, `insert into catalog.release_media(tenant_id,media_id,maker,original_sha256,sha256,width,height,png) values($1,$2,$3,$4,$5,$6,$7,$8)`,
		p.TenantID, media.ID, p.Subject, media.OriginalSHA256, media.SHA256, media.Width, media.Height, clean)
	if err != nil {
		return cr.Receipt{}, err
	}
	r := cr.Receipt{CommandID: key, Actor: p.Subject, Kind: "media", ResourceID: media.ID, RequestSHA256: hash}
	if err = s.save(ctx, tx, p, r); err != nil {
		return cr.Receipt{}, err
	}
	return r, tx.Commit(ctx)
}
func (s *CatalogRelease) CreateDraft(ctx context.Context, p identity.Principal, input cr.DraftRequest) (cr.Receipt, error) {
	if !input.Valid() {
		return cr.Receipt{}, cr.ErrInvalid
	}
	_, hash, err := cr.Canonical(input)
	if err != nil {
		return cr.Receipt{}, err
	}
	tx, err := s.begin(ctx, p, "catalog:draft")
	if err != nil {
		return cr.Receipt{}, err
	}
	defer tx.Rollback(ctx)
	if r, yes, err := s.replay(ctx, tx, p, input.CommandID, "draft", hash); err != nil || yes {
		return r, err
	}
	snap := cr.Snapshot{Schema: "elite-catalog-snapshot/v1", Profile: s.profile, PriceBookID: input.PriceBookID, Models: []cr.Model{}, Variants: []cr.Variant{}}
	if err = tx.QueryRow(ctx, `select tenant_code from platform.tenant where tenant_id=$1`, p.TenantID).Scan(&snap.TenantCode); err != nil {
		return cr.Receipt{}, err
	}
	err = tx.QueryRow(ctx, `select valid_from,valid_until from pricing.price_book where tenant_id=$1 and price_book_id=$2 and market=$3 and currency=$4 for update`, p.TenantID, input.PriceBookID, s.profile.Market, s.profile.Currency).Scan(&snap.ValidFrom, &snap.ValidUntil)
	if errors.Is(err, pgx.ErrNoRows) {
		return cr.Receipt{}, cr.ErrNotFound
	}
	if err != nil {
		return cr.Receipt{}, err
	}
	selected := append([]cr.ModelReference(nil), input.Models...)
	sort.Slice(selected, func(i, j int) bool { return selected[i].ModelID < selected[j].ModelID })
	modelIDs := map[string]bool{}
	for _, ref := range selected {
		var m cr.Model
		err = tx.QueryRow(ctx, `select model_id,model_code,display_name,vehicle_class,specification from catalog.vehicle_model
    where tenant_id=$1 and model_id=$2 and lifecycle_state in('draft','active') and octet_length(display_name)<=160 and pg_column_size(specification)<=32768 for share`, p.TenantID, ref.ModelID).Scan(&m.ID, &m.Code, &m.DisplayName, &m.VehicleClass, &m.Specification)
		if errors.Is(err, pgx.ErrNoRows) {
			return cr.Receipt{}, cr.ErrNotFound
		}
		if err != nil {
			return cr.Receipt{}, err
		}
		err = tx.QueryRow(ctx, `select media_id,original_sha256,sha256,width,height from catalog.release_media where tenant_id=$1 and media_id=$2`, p.TenantID, ref.MediaID).Scan(&m.Media.ID, &m.Media.OriginalSHA256, &m.Media.SHA256, &m.Media.Width, &m.Media.Height)
		if errors.Is(err, pgx.ErrNoRows) {
			return cr.Receipt{}, cr.ErrNotFound
		}
		if err != nil {
			return cr.Receipt{}, err
		}
		m.CanonicalURL = s.profile.Origin + "/models/" + url.PathEscape(m.Code)
		snap.Models = append(snap.Models, m)
		modelIDs[m.ID] = false
	}
	rows, err := tx.Query(ctx, `select v.variant_id,v.model_id,v.variant_code,v.display_name,v.battery_specification,v.homologation_state,e.amount_minor_units,e.tax_mode
 from pricing.price_book_entry e join catalog.vehicle_variant v using(tenant_id,variant_id)
 where e.tenant_id=$1 and e.price_book_id=$2 and v.lifecycle_state in('draft','active') and pg_column_size(v.battery_specification)<=32768
 order by v.variant_id limit 129 for share of e,v`, p.TenantID, input.PriceBookID)
	if err != nil {
		return cr.Receipt{}, err
	}
	for rows.Next() {
		var v cr.Variant
		if err = rows.Scan(&v.ID, &v.ModelID, &v.Code, &v.DisplayName, &v.BatterySpecification, &v.HomologationState, &v.AmountMinorUnits, &v.TaxMode); err != nil {
			rows.Close()
			return cr.Receipt{}, err
		}
		if _, ok := modelIDs[v.ModelID]; !ok {
			rows.Close()
			return cr.Receipt{}, cr.ErrInvalid
		}
		modelIDs[v.ModelID] = true
		snap.Variants = append(snap.Variants, v)
	}
	err = rows.Err()
	rows.Close()
	if err != nil {
		return cr.Receipt{}, err
	}
	if len(snap.Variants) == 0 || len(snap.Variants) > 128 {
		return cr.Receipt{}, cr.ErrInvalid
	}
	for _, present := range modelIDs {
		if !present {
			return cr.Receipt{}, cr.ErrInvalid
		}
	}
	var entryCount int
	if err = tx.QueryRow(ctx, `select count(*) from pricing.price_book_entry where tenant_id=$1 and price_book_id=$2`, p.TenantID, input.PriceBookID).Scan(&entryCount); err != nil {
		return cr.Receipt{}, err
	}
	if entryCount != len(snap.Variants) {
		return cr.Receipt{}, cr.ErrInvalid
	}
	raw, sha, err := cr.Canonical(snap)
	if err != nil || len(raw) > 262144 {
		return cr.Receipt{}, cr.ErrInvalid
	}
	id := (randomid.Generator{}).New()
	_, err = tx.Exec(ctx, `insert into catalog.release_draft(tenant_id,draft_id,maker,price_book_id,snapshot,snapshot_canonical,snapshot_sha256) values($1,$2,$3,$4,$5,$6,$7)`, p.TenantID, id, p.Subject, input.PriceBookID, raw, string(raw), sha)
	if err != nil {
		return cr.Receipt{}, err
	}
	approvals := NewHumanApprovals(s.pool)
	for _, stage := range []string{"legal", "technical", "media", "publication"} {
		payload, h, err := approval.CanonicalPayload([]byte(`{"draft_id":"` + id + `","snapshot_sha256":"` + sha + `","stage":"` + stage + `"}`))
		if err != nil {
			return cr.Receipt{}, err
		}
		aid := (randomid.Generator{}).New()
		_, err = tx.Exec(ctx, `insert into catalog.release_review(tenant_id,draft_id,stage,approval_id,payload_sha256) values($1,$2,$3,$4,$5)`, p.TenantID, id, stage, aid, h)
		if err != nil {
			return cr.Receipt{}, err
		}
		spec := HumanApprovalSpec{Request: approval.Request{TenantID: p.TenantID, ID: aid, Kind: approval.KindCatalogReview, SubjectID: id, Requester: p.Subject, EvidenceSHA: h}, OrganizationID: s.profile.OrganizationID, Payload: payload}
		if _, err = approvals.submitTx(ctx, tx, p, spec, "catalog:draft", nil); err != nil {
			return cr.Receipt{}, err
		}
	}
	r := cr.Receipt{CommandID: input.CommandID, Actor: p.Subject, Kind: "draft", ResourceID: id, RequestSHA256: hash, SnapshotSHA256: sha}
	if err = s.save(ctx, tx, p, r); err != nil {
		return cr.Receipt{}, err
	}
	return r, tx.Commit(ctx)
}
func (s *CatalogRelease) Draft(ctx context.Context, p identity.Principal, id string) (cr.Draft, error) {
	if !s.authorized(p, "catalog:read") || !cr.ValidID(id) {
		return cr.Draft{}, cr.ErrNotFound
	}
	var d cr.Draft
	err := s.pool.QueryRow(ctx, `select draft_id,maker,snapshot_sha256,snapshot_canonical from catalog.release_draft where tenant_id=$1 and draft_id=$2`, p.TenantID, id).Scan(&d.ID, &d.Maker, &d.SHA256, &d.Snapshot)
	if errors.Is(err, pgx.ErrNoRows) {
		return d, cr.ErrNotFound
	}
	if err != nil {
		return d, err
	}
	var raw []byte
	err = s.pool.QueryRow(ctx, `select jsonb_object_agg(r.stage,q.state) from catalog.release_review r join approval.request q on q.tenant_id=r.tenant_id and q.request_id=r.approval_id where r.tenant_id=$1 and r.draft_id=$2`, p.TenantID, id).Scan(&raw)
	if err != nil {
		return d, err
	}
	err = json.Unmarshal(raw, &d.Reviews)
	return d, err
}
````

### FILE: `internal/platform/postgres/catalog_release_http_integration_test.go`

```yaml
block_id: "GO-CONNECTED-CATALOG-PUBLICATION:file24:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "57ad341ab57642171a5ce6452f9240ea51cc4b6e40f8c6eba1c7b32faea311f0"
variables: []
secrets_allowed: false
```

````go
package postgres_test

// AUTHORED real HTTP/PG fixture. Bearer verifier assigns explicit principals;
// this does not substitute T2803 JWT/IdP evidence.
import (
	"bytes"
	"context"
	cr "elite.local/enterprise/internal/catalogrelease"
	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/platform/identity"
	db "elite.local/enterprise/internal/platform/postgres"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"
)

type catalogHTTPFixture struct {
	publisher identity.Principal
	feedProbe func(context.Context, identity.Principal, int64, []byte)
	webProbe  func(cr.PublicDocument)

	t            *testing.T
	client       *http.Client
	base         string
	mu           sync.Mutex
	counter      int
	principals   map[string]identity.Principal
	accepted     map[string]int
	dropped      bool
	recoveries   int
	etag         string
	cached       []byte
	cacheChanges int
}

func (c *catalogHTTPFixture) Verify(_ context.Context, token string) (identity.Principal, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	p, ok := c.principals[token]
	if !ok {
		return p, identity.ErrUnauthenticated
	}
	return p, nil
}
func (c *catalogHTTPFixture) request(ctx context.Context, p *identity.Principal, method, path, contentType string, body []byte, etag string) ([]byte, int, string, error) {
	if p != nil {
		org := ""
		for v := range p.Organizations {
			org = v
			break
		}
		path += "?organization_id=" + url.QueryEscape(org)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.base+path, bytes.NewReader(body))
	if err != nil {
		return nil, 0, "", err
	}
	req.Header.Set("Content-Type", contentType)
	if etag != "" {
		req.Header.Set("If-None-Match", etag)
	}
	if p != nil {
		c.mu.Lock()
		c.counter++
		token := fmt.Sprintf("catalog-fixture-%d", c.counter)
		c.principals[token] = *p
		c.mu.Unlock()
		req.Header.Set("Authorization", "Bearer "+token)
	}
	response, err := c.client.Do(req)
	if err != nil {
		return nil, 0, "", err
	}
	defer response.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(response.Body, (4<<20)+1))
	if err != nil {
		return nil, 0, "", err
	}
	if len(raw) > 4<<20 {
		return nil, 0, "", errors.New("oversize fixture response")
	}
	if p != nil && response.Header.Get("Cache-Control") != "no-store" {
		return nil, 0, "", errors.New("private response cacheable")
	}
	if response.StatusCode >= 400 {
		return nil, response.StatusCode, "", fmt.Errorf("HTTP %d: %s", response.StatusCode, raw)
	}
	return raw, response.StatusCode, response.Header.Get("ETag"), nil
}
func (c *catalogHTTPFixture) call(ctx context.Context, p identity.Principal, method, path string, in, out any) error {
	var raw []byte
	var err error
	if in != nil {
		raw, err = json.Marshal(in)
		if err != nil {
			return err
		}
	}
	body, _, _, err := c.request(ctx, &p, method, path, "application/json", raw, "")
	if err != nil {
		return err
	}
	return json.Unmarshal(body, out)
}
func (c *catalogHTTPFixture) UploadPNG(ctx context.Context, p identity.Principal, key string, raw []byte) (cr.Receipt, error) {
	var out cr.Receipt
	body, _, _, err := c.request(ctx, &p, "POST", "/v1/admin/catalog/media/"+url.PathEscape(key), "image/png", raw, "")
	if err == nil {
		err = json.Unmarshal(body, &out)
	}
	return out, err
}
func (c *catalogHTTPFixture) CreateDraft(ctx context.Context, p identity.Principal, in cr.DraftRequest) (cr.Receipt, error) {
	var out cr.Receipt
	err := c.call(ctx, p, "POST", "/v1/admin/catalog/drafts", in, &out)
	return out, err
}
func (c *catalogHTTPFixture) Draft(ctx context.Context, p identity.Principal, id string) (cr.Draft, error) {
	var out cr.Draft
	err := c.call(ctx, p, "GET", "/v1/admin/catalog/drafts/"+url.PathEscape(id), nil, &out)
	return out, err
}
func (c *catalogHTTPFixture) Review(ctx context.Context, p identity.Principal, in cr.ReviewRequest) (cr.Receipt, error) {
	var out cr.Receipt
	err := c.call(ctx, p, "POST", "/v1/admin/catalog/drafts/"+url.PathEscape(in.DraftID)+"/review/"+in.Stage, in, &out)
	return out, err
}
func (c *catalogHTTPFixture) CommandReceipt(ctx context.Context, p identity.Principal, id string) (cr.Receipt, error) {
	var out cr.Receipt
	err := c.call(ctx, p, "GET", "/v1/admin/catalog/commands/"+url.PathEscape(id), nil, &out)
	return out, err
}
func (c *catalogHTTPFixture) Publish(ctx context.Context, p identity.Principal, in cr.PublishRequest) (cr.Receipt, error) {
	c.mu.Lock()
	c.publisher = p
	c.mu.Unlock()
	var out cr.Receipt
	err := c.call(ctx, p, "POST", "/v1/admin/catalog/drafts/"+url.PathEscape(in.DraftID)+"/publish", in, &out)
	if err == nil {
		return out, nil
	}
	var network *url.Error
	if !errors.As(err, &network) {
		return out, err
	}
	got, getErr := c.CommandReceipt(ctx, p, in.CommandID)
	if getErr != nil {
		return out, getErr
	}
	_, hash, hashErr := cr.Canonical(in)
	if hashErr != nil || got.RequestSHA256 != hash || got.Actor != p.Subject || got.CommandID != in.CommandID || got.ResourceID != in.DraftID || got.SnapshotSHA256 != in.SnapshotSHA256 {
		return out, errors.New("unconfirmed publication identity")
	}
	c.mu.Lock()
	c.recoveries++
	c.mu.Unlock()
	return got, nil
}
func (c *catalogHTTPFixture) Public(ctx context.Context) (cr.Publication, error) {
	raw, status, etag, err := c.request(ctx, nil, "GET", "/v1/public/catalog", "", nil, c.etag)
	if err != nil {
		return cr.Publication{}, err
	}
	if status == 304 {
		raw = c.cached
	} else {
		if c.etag != "" {
			if etag == c.etag {
				return cr.Publication{}, errors.New("new version did not invalidate cache")
			}
			c.cacheChanges++
		}
		c.etag = etag
		c.cached = raw
	}
	_, sameStatus, _, err := c.request(ctx, nil, "GET", "/v1/public/catalog", "", nil, etag)
	if err != nil || sameStatus != 304 {
		return cr.Publication{}, fmt.Errorf("unchanged ETag %d %v", sameStatus, err)
	}
	feed, _, _, err := c.request(ctx, nil, "GET", "/v1/public/catalog/feed", "", nil, "")
	if err != nil || !bytes.Equal(raw, feed) {
		return cr.Publication{}, fmt.Errorf("feed differs from storefront: %v", err)
	}
	if bytes.Contains(raw, []byte("tenant_id")) || bytes.Contains(raw, []byte("organization_id")) {
		return cr.Publication{}, errors.New("private stream profile leaked")
	}
	var document cr.PublicDocument
	if err = json.Unmarshal(raw, &document); err != nil {
		return cr.Publication{}, err
	}
	if c.feedProbe != nil {
		c.feedProbe(ctx, c.publisher, document.Generation, raw)
	}
	if c.webProbe != nil {
		c.webProbe(document)
	}
	return cr.Publication{Generation: document.Generation, SHA256: document.SourceSHA256, EffectivePriceBookID: document.EffectivePriceBookID}, nil
}
func (c *catalogHTTPFixture) PublicMedia(ctx context.Context, sha string) ([]byte, error) {
	raw, _, _, err := c.request(ctx, nil, "GET", "/v1/public/catalog/media/"+sha, "", nil, "")
	return raw, err
}
func (c *catalogHTTPFixture) PublicSearch(ctx context.Context, q string) ([]cr.SearchHit, error) {
	raw, _, _, err := c.request(ctx, nil, "GET", "/v1/public/catalog/search?q="+url.QueryEscape(q), "", nil, "")
	if err != nil {
		return nil, err
	}
	var out []cr.SearchHit
	err = json.Unmarshal(raw, &out)
	return out, err
}
func TestCatalogReleaseHTTPReference(t *testing.T) { testCatalogReleaseHTTPReference(t, false) }
func TestCatalogReleaseFeedReference(t *testing.T) { testCatalogReleaseHTTPReference(t, true) }
func testCatalogReleaseHTTPReference(t *testing.T, withFeed bool) {
	var client *catalogHTTPFixture
	testCatalogReleaseConnected(t, func(t *testing.T, store *db.CatalogRelease) catalogReleaseReference {
		client = &catalogHTTPFixture{t: t, principals: map[string]identity.Principal{}, accepted: map[string]int{}}
		mux := http.NewServeMux()
		var feeder httpapi.CatalogFeedService
		if withFeed {
			feeder = attachCatalogFeedFixture(t, client, store)
		}
		httpapi.CatalogReleaseModule{Service: store, Feed: feeder}.Register(mux, client)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var cmd cr.PublishRequest
			if r.Method == "POST" && strings.HasSuffix(r.URL.Path, "/publish") {
				raw, _ := io.ReadAll(io.LimitReader(r.Body, 32769))
				r.Body.Close()
				r.Body = io.NopCloser(bytes.NewReader(raw))
				_ = json.Unmarshal(raw, &cmd)
			}
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, r)
			drop := false
			if cmd.CommandID != "" && (rec.Code == 200 || rec.Code == 201) {
				client.mu.Lock()
				client.accepted[cmd.CommandID]++
				if cmd.CommandID == "publish-two" && rec.Code == 201 && !client.dropped {
					client.dropped = true
					drop = true
				}
				client.mu.Unlock()
			}
			if drop {
				conn, _, err := w.(http.Hijacker).Hijack()
				if err != nil {
					t.Error(err)
					return
				}
				conn.Close()
				return
			}
			for k, values := range rec.Header() {
				for _, value := range values {
					w.Header().Add(k, value)
				}
			}
			w.WriteHeader(rec.Code)
			_, _ = w.Write(rec.Body.Bytes())
		}))
		transport := &http.Transport{DisableKeepAlives: true}
		client.client = &http.Client{Transport: transport, Timeout: 5 * time.Second}
		client.base = server.URL
		t.Cleanup(func() { transport.CloseIdleConnections(); server.Close() })
		if catalogStorefrontFixtureRequested() {
			attachCatalogStorefrontFixture(t, client)
		}
		return client
	})
	if client.recoveries != 1 || !client.dropped || client.accepted["publish-two"] != 1 || client.cacheChanges != 3 {
		t.Fatal("HTTP recovery/cache coverage", client.recoveries, client.accepted, client.cacheChanges)
	}
	t.Log("CATALOG_RELEASE_HTTP_PASS connected role commands; lost committed response recovered by GET without repeat POST; storefront/feed byte-identical; three stale ETags invalidated; normalized media and scoped FTS match current publication")
}
````

### FILE: `internal/platform/postgres/catalog_release_integration_test.go`

```yaml
block_id: "GO-CONNECTED-CATALOG-PUBLICATION:file25:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "87e2fdc7e2e789877540ef19908b67e514d9f8066301228286def345165e76c0"
variables: []
secrets_allowed: false
```

````go
package postgres_test

// AUTHORED actual PostgreSQL owner composition; synthetic principals and media.
import (
	"bytes"
	"context"
	cr "elite.local/enterprise/internal/catalogrelease"
	"elite.local/enterprise/internal/commerce"
	"elite.local/enterprise/internal/electromobility"
	"elite.local/enterprise/internal/platform/identity"
	db "elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/platform/randomid"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"image"
	"image/png"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

type catalogReleaseReference interface {
	UploadPNG(context.Context, identity.Principal, string, []byte) (cr.Receipt, error)
	CreateDraft(context.Context, identity.Principal, cr.DraftRequest) (cr.Receipt, error)
	Draft(context.Context, identity.Principal, string) (cr.Draft, error)
	Review(context.Context, identity.Principal, cr.ReviewRequest) (cr.Receipt, error)
	Publish(context.Context, identity.Principal, cr.PublishRequest) (cr.Receipt, error)
	CommandReceipt(context.Context, identity.Principal, string) (cr.Receipt, error)
	Public(context.Context) (cr.Publication, error)
	PublicMedia(context.Context, string) ([]byte, error)
	PublicSearch(context.Context, string) ([]cr.SearchHit, error)
}

func TestCatalogReleaseConnectedReference(t *testing.T) { testCatalogReleaseConnected(t, nil) }
func testCatalogReleaseConnected(t *testing.T, hook func(*testing.T, *db.CatalogRelease) catalogReleaseReference) {

	rawURL := os.Getenv("PAYMENT_CONNECTED_DB_URL")
	if rawURL == "" {
		t.Skip("owned fixture DB not configured")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, rawURL)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := "ca6c9cb7-e548-4f19-a2e6-27d2aa2d4893"
	for _, q := range []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,'catalog-j3','Fixture','Fixture')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,'j3-store','j3-store','Fixture','store')`,
	} {
		if _, err = pool.Exec(ctx, q, tenant); err != nil {
			t.Fatal(err)
		}
	}
	profile := cr.Profile{Schema: "elite-catalog-publication/v1", TenantID: tenant, OrganizationID: "j3-store", Market: "AR", Currency: "ARS", Origin: "https://catalog.example.invalid", PolicyCode: cr.PolicyCode}
	price := db.NewCommerce(pool)
	backend, err := db.NewCatalogRelease(pool, profile, price)
	if err != nil {
		t.Fatal(err)
	}
	var store catalogReleaseReference = backend
	if hook != nil {
		store = hook(t, backend)
	}
	person := func(subject string) identity.Principal {
		return identity.Principal{TenantID: tenant, Subject: subject, Permissions: map[string]struct{}{"catalog:feed": {}, "catalog:draft": {}, "catalog:read": {}, "catalog:publish": {}, "catalog:review:legal": {}, "catalog:review:technical": {}, "catalog:review:media": {}, "catalog:review:publication": {}}, Organizations: map[string]struct{}{"j3-store": {}}}
	}
	maker, reviewer := person("j3-maker"), person("j3-reviewer")
	var pngBytes bytes.Buffer
	if err = png.Encode(&pngBytes, image.NewNRGBA(image.Rect(0, 0, 2, 2))); err != nil {
		t.Fatal(err)
	}
	media, err := store.UploadPNG(ctx, maker, "upload-png", append(pngBytes.Bytes(), []byte("<script>untrusted trailer</script>")...))
	if err != nil {
		t.Fatal("upload", err)
	}
	model, err := electromobility.NewService(db.NewElectromobility(pool), randomid.Generator{}).CreateModel(ctx, tenant, electromobility.Model{Code: "bicycle-one", DisplayName: "Version One", VehicleClass: "bicycle", Specification: json.RawMessage(`{"range_km":20}`)})
	if err != nil {
		t.Fatal(err)
	}
	_, err = pool.Exec(ctx, `insert into catalog.vehicle_variant(tenant_id,variant_id,model_id,variant_code,display_name,battery_specification,lifecycle_state) values($1,'j3-variant',$2,'j3-variant','Fixture','{}','active')`, tenant, model.ID)
	if err != nil {
		t.Fatal(err)
	}
	sales := commerce.NewService(price, randomid.Generator{})
	book := func(amount int64, until *time.Time) commerce.PriceBook {
		t.Helper()
		var now time.Time
		if err := pool.QueryRow(ctx, `select clock_timestamp()`).Scan(&now); err != nil {
			t.Fatal(err)
		}
		out, err := sales.CreatePriceBook(ctx, tenant, commerce.PriceBook{Market: "AR", Currency: "ARS", ValidFrom: now.Add(-time.Hour), ValidUntil: until, Entries: []commerce.PriceEntry{{VariantID: "j3-variant", AmountMinorUnits: amount, TaxMode: "not-applicable"}}})
		if err != nil {
			t.Fatal("price book", err)
		}
		return out
	}
	firstBook := book(120000, nil)
	draft := func(key, bookID string) cr.Receipt {
		t.Helper()
		r, err := store.CreateDraft(ctx, maker, cr.DraftRequest{CommandID: key, PriceBookID: bookID, Models: []cr.ModelReference{{ModelID: model.ID, MediaID: media.ResourceID}}})
		if err != nil {
			t.Fatal("draft", err)
		}
		return r
	}
	one := draft("draft-one", firstBook.ID)
	snapshot := func() string {
		t.Helper()
		var out string
		for _, table := range []string{"catalog.release_stream", "catalog.release_media", "catalog.release_draft", "catalog.release_review", "catalog.release_review_action", "catalog.release_command", "catalog.release_publication", "pricing.price_book", "pricing.price_book_entry", "search.document", "approval.request", "approval.decision", "platform.outbox_event"} {
			var text string
			if err := pool.QueryRow(ctx, fmt.Sprintf("select coalesce(jsonb_agg(to_jsonb(x) order by to_jsonb(x)::text),'[]'::jsonb)::text from %s x where tenant_id=$1", table), tenant).Scan(&text); err != nil {
				t.Fatal(table, err)
			}
			out += text
		}
		return out
	}
	expectDenied := func(action func() error, name string) {
		t.Helper()
		before := snapshot()
		if err := action(); err == nil {
			t.Fatal(name + " accepted")
		}
		if snapshot() != before {
			t.Fatal(name + " leaked effects")
		}
	}
	requestReview := func(d cr.Receipt, stage, key string) cr.ReviewRequest {
		return cr.ReviewRequest{CommandID: key, DraftID: d.ResourceID, Stage: stage, SnapshotSHA256: d.SnapshotSHA256, Approved: true, Reason: "synthetic reviewer evidence; no legal certification", EvidenceSHA256: strings.Repeat("a", 64)}
	}
	expectDenied(func() error { _, e := store.Review(ctx, maker, requestReview(one, "legal", "self-review")); return e }, "self review")
	expectDenied(func() error {
		_, e := store.Review(ctx, reviewer, requestReview(one, "publication", "early-final"))
		return e
	}, "early final review")
	var aid, approvalHash string
	if err = pool.QueryRow(ctx, `select approval_id,payload_sha256 from catalog.release_review where tenant_id=$1 and draft_id=$2 and stage='legal'`, tenant, one.ResourceID).Scan(&aid, &approvalHash); err != nil {
		t.Fatal(err)
	}
	expectDenied(func() error {
		_, e := db.NewHumanApprovals(pool).Decide(ctx, reviewer, tenant, aid, "j3-store", approvalHash, true, "bypass", "catalog:review:legal", nil)
		return e
	}, "generic approval bypass")
	approve := func(d cr.Receipt) {
		t.Helper()
		for _, stage := range []string{"legal", "technical", "media", "publication"} {
			if _, err := store.Review(ctx, reviewer, requestReview(d, stage, d.CommandID+"-"+stage)); err != nil {
				t.Fatal(stage, err)
			}
		}
	}
	if public, err := db.NewElectromobility(pool).ListPublicModels(ctx, "catalog-j3"); err != nil || len(public) != 0 {
		t.Fatal("draft visible", public, err)
	}
	if _, err = pool.Exec(ctx, `update catalog.vehicle_model set display_name='Version Two' where tenant_id=$1 and model_id=$2`, tenant, model.ID); err != nil {
		t.Fatal(err)
	}
	captured, err := store.Draft(ctx, maker, one.ResourceID)
	if err != nil || !bytes.Contains(captured.Snapshot, []byte("Version One")) || bytes.Contains(captured.Snapshot, []byte("Version Two")) {
		t.Fatal("captured source mutated", err)
	}
	expectDenied(func() error { return sales.ActivatePriceBook(ctx, tenant, firstBook.ID) }, "generic price activation")
	expectDenied(func() error {
		_, e := pool.Exec(ctx, `update pricing.price_book_entry set amount_minor_units=1 where tenant_id=$1 and price_book_id=$2`, tenant, firstBook.ID)
		return e
	}, "captured price update")
	approve(one)
	pub := cr.PublishRequest{CommandID: "publish-one", DraftID: one.ResourceID, SnapshotSHA256: one.SnapshotSHA256, ExpectedGeneration: 0, Reason: "approved synthetic release"}
	expectDenied(func() error { _, e := store.Publish(ctx, maker, pub); return e }, "maker publish")
	var wg sync.WaitGroup
	var mu sync.Mutex
	newCount, replays := 0, 0
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r, e := store.Publish(ctx, reviewer, pub)
			mu.Lock()
			defer mu.Unlock()
			if e != nil {
				t.Error(e)
				return
			}
			if r.Replay {
				replays++
			} else {
				newCount++
			}
		}()
	}
	wg.Wait()
	if newCount != 1 || replays != 7 {
		t.Fatal("concurrency", newCount, replays)
	}
	checkPublic := func(name string, amount int64, generation int64, sha string) {
		t.Helper()
		models, e := db.NewElectromobility(pool).ListPublicModels(ctx, "catalog-j3")
		if e != nil || len(models) != 1 || models[0].DisplayName != name {
			t.Fatal("public owner", models, e)
		}
		priced, e := sales.PublicPrice(ctx, "catalog-j3", "AR", "j3-variant")
		if e != nil || priced.AmountMinorUnits != amount {
			t.Fatal("price projection", priced, e)
		}
		release, e := store.Public(ctx)
		if e != nil || release.Generation != generation || release.SHA256 != sha {
			t.Fatal("public release", release, e)
		}
		hits, e := store.PublicSearch(ctx, "Version")
		if e != nil || len(hits) != 1 || hits[0].Title != name || hits[0].Generation != generation || hits[0].SourceSHA256 != sha {
			t.Fatal("public search projection", hits, e)
		}
		var assetHash string
		if e = pool.QueryRow(ctx, `select sha256 from catalog.release_media where tenant_id=$1 and media_id=$2`, tenant, media.ResourceID).Scan(&assetHash); e != nil {
			t.Fatal(e)
		}
		asset, e := store.PublicMedia(ctx, assetHash)
		if e != nil || bytes.Contains(asset, []byte("script")) {
			t.Fatal("public normalized media", e)
		}
		var count int
		e = pool.QueryRow(ctx, `select count(*) from search.document where tenant_id=$1 and kind='catalog-model' and content_sha256=$2 and title=$3`, tenant, sha, name).Scan(&count)
		if e != nil || count != 1 {
			t.Fatal("search source mismatch", e, count)
		}
	}
	checkPublic("Version One", 120000, 1, one.SnapshotSHA256)
	secondBook := book(130000, nil)
	two := draft("draft-two", secondBook.ID)
	approve(two)
	next := cr.PublishRequest{CommandID: "publish-two", DraftID: two.ResourceID, SnapshotSHA256: two.SnapshotSHA256, ExpectedGeneration: 1, Reason: "approved second version"}
	_, err = pool.Exec(ctx, `create function catalog.reject_fixture_outbox() returns trigger language plpgsql as $$begin if new.event_type='catalog-release.publish' then raise exception 'fixture outbox unavailable';end if;return new;end$$;create trigger reject_catalog_outbox before insert on platform.outbox_event for each row execute function catalog.reject_fixture_outbox()`)
	if err != nil {
		t.Fatal(err)
	}
	expectDenied(func() error { _, e := store.Publish(ctx, reviewer, next); return e }, "publication outbox rollback")
	if _, err = pool.Exec(ctx, `drop trigger reject_catalog_outbox on platform.outbox_event;drop function catalog.reject_fixture_outbox()`); err != nil {
		t.Fatal(err)
	}
	if _, err = store.Publish(ctx, reviewer, next); err != nil {
		t.Fatal("second publish", err)
	}
	checkPublic("Version Two", 130000, 2, two.SnapshotSHA256)
	rollback := cr.PublishRequest{CommandID: "rollback-one", DraftID: one.ResourceID, SnapshotSHA256: one.SnapshotSHA256, ExpectedGeneration: 2, Reason: "explicit rollback to prior approved snapshot"}
	if _, err = store.Publish(ctx, reviewer, rollback); err != nil {
		t.Fatal("rollback", err)
	}
	checkPublic("Version One", 120000, 3, one.SnapshotSHA256)
	expectDenied(func() error {
		_, e := pool.Exec(ctx, `update search.document set title='bypass' where tenant_id=$1 and kind='catalog-model'`, tenant)
		return e
	}, "search bypass")
	// Publish a briefly valid version, leave it, then reject rollback after expiry.
	var now time.Time
	if err = pool.QueryRow(ctx, `select clock_timestamp()`).Scan(&now); err != nil {
		t.Fatal(err)
	}
	until := now.Add(2 * time.Second)
	shortBook := book(140000, &until)
	three := draft("draft-short", shortBook.ID)
	approve(three)
	short := cr.PublishRequest{CommandID: "publish-short", DraftID: three.ResourceID, SnapshotSHA256: three.SnapshotSHA256, ExpectedGeneration: 3, Reason: "synthetic expiry boundary"}
	if _, err = store.Publish(ctx, reviewer, short); err != nil {
		t.Fatal(err)
	}
	next.CommandID = "return-two"
	next.ExpectedGeneration = 4
	if _, err = store.Publish(ctx, reviewer, next); err != nil {
		t.Fatal(err)
	}
	if _, err = pool.Exec(ctx, `select pg_sleep(greatest(0,extract(epoch from $1::timestamptz-clock_timestamp()))+0.02)`, until); err != nil {
		t.Fatal(err)
	}
	short.CommandID = "expired-rollback"
	short.ExpectedGeneration = 5
	expectDenied(func() error { _, e := store.Publish(ctx, reviewer, short); return e }, "expired rollback")
	checkPublic("Version Two", 130000, 5, two.SnapshotSHA256)
	otherPool, err := pgxpool.New(ctx, rawURL)
	if err != nil {
		t.Fatal(err)
	}
	defer otherPool.Close()
	reopened, err := db.NewCatalogRelease(otherPool, profile, db.NewCommerce(otherPool))
	if err != nil {
		t.Fatal(err)
	}
	recovered, err := reopened.CommandReceipt(ctx, reviewer, "rollback-one")
	if err != nil || recovered.Generation != 3 || recovered.SnapshotSHA256 != one.SnapshotSHA256 {
		t.Fatal("durable recovery", recovered, err)
	}
	t.Log("CATALOG_RELEASE_CONNECTED_PASS three immutable snapshots/four shared reviews each; five publications including rollback; source edit isolated; storefront/price/search match; concurrency1new7replay;13table rollback; expired rollback denied; durable recovery")
}
````

### FILE: `internal/platform/postgres/catalog_release_public.go`

```yaml
block_id: "GO-CONNECTED-CATALOG-PUBLICATION:file26:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "f22ef6787cec2cc7d7e4ab0816ece6a2b9d81ce70971108d130fc6194170df5f"
variables: []
secrets_allowed: false
```

````go
package postgres

import (
	"context"
	cr "elite.local/enterprise/internal/catalogrelease"
	"errors"
	"github.com/jackc/pgx/v5"
)

// PublicMedia only resolves an immutable normalized asset in the currently
// eligible publication. Knowing another draft's hash does not publish its media.
func (s *CatalogRelease) PublicMedia(ctx context.Context, hash string) ([]byte, error) {
	if !cr.ValidSHA(hash) {
		return nil, cr.ErrNotFound
	}
	var raw []byte
	err := s.pool.QueryRow(ctx, `select a.png from catalog.release_media a
 join catalog.release_public_current c using(tenant_id)
 where a.tenant_id=$1 and a.sha256=$2 and exists(select 1 from jsonb_array_elements(c.snapshot->'models') m
 where m->'media'->>'id'=a.media_id and m->'media'->>'sha256'=a.sha256) limit 1`, s.profile.TenantID, hash).Scan(&raw)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, cr.ErrNotFound
	}
	return raw, err
}
func (s *CatalogRelease) PublicSearch(ctx context.Context, q string) ([]cr.SearchHit, error) {
	if !cr.ValidText(q, 200) {
		return nil, cr.ErrInvalid
	}
	rows, err := s.pool.Query(ctx, `select x.external_id,x.title,x.body,c.generation,c.snapshot_sha256
 from search.document x join catalog.release_public_current c on c.tenant_id=x.tenant_id and c.snapshot_sha256=x.content_sha256
 cross join websearch_to_tsquery('simple',$2) query
 where x.tenant_id=$1 and x.kind='catalog-model' and x.tsv@@query
 order by ts_rank(x.tsv,query) desc,x.external_id limit 20`, s.profile.TenantID, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []cr.SearchHit{}
	for rows.Next() {
		var hit cr.SearchHit
		if err = rows.Scan(&hit.ModelID, &hit.Title, &hit.URL, &hit.Generation, &hit.SourceSHA256); err != nil {
			return nil, err
		}
		out = append(out, hit)
	}
	return out, rows.Err()
}
````

### FILE: `internal/platform/postgres/catalog_release_publish.go`

```yaml
block_id: "GO-CONNECTED-CATALOG-PUBLICATION:file27:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "6e6247d574eb2630a87d34981a40438395eef0d17f7c4d60bb03ace1848c499b"
variables: []
secrets_allowed: false
```

````go
package postgres

import (
	"context"
	cr "elite.local/enterprise/internal/catalogrelease"
	"elite.local/enterprise/internal/commerce"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/randomid"
	"elite.local/enterprise/internal/search"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5"
	"math"
)

func (s *CatalogRelease) Review(ctx context.Context, p identity.Principal, input cr.ReviewRequest) (cr.Receipt, error) {
	if !input.Valid() {
		return cr.Receipt{}, cr.ErrInvalid
	}
	permission := "catalog:review:" + input.Stage
	tx, err := s.begin(ctx, p, permission)
	if err != nil {
		return cr.Receipt{}, err
	}
	defer tx.Rollback(ctx)
	_, hash, err := cr.Canonical(input)
	if err != nil {
		return cr.Receipt{}, err
	}
	if r, yes, err := s.replay(ctx, tx, p, input.CommandID, "review", hash); err != nil || yes {
		return r, err
	}
	var aid, payloadSHA, maker string
	err = tx.QueryRow(ctx, `select r.approval_id,r.payload_sha256,d.maker from catalog.release_review r join catalog.release_draft d using(tenant_id,draft_id)
 where r.tenant_id=$1 and r.draft_id=$2 and r.stage=$3 and d.snapshot_sha256=$4`, p.TenantID, input.DraftID, input.Stage, input.SnapshotSHA256).Scan(&aid, &payloadSHA, &maker)
	if errors.Is(err, pgx.ErrNoRows) {
		return cr.Receipt{}, cr.ErrNotFound
	}
	if err != nil {
		return cr.Receipt{}, err
	}
	if maker == p.Subject {
		return cr.Receipt{}, cr.ErrConflict
	}
	if input.Stage == "publication" && input.Approved {
		var count int
		err = tx.QueryRow(ctx, `select count(*) from catalog.release_review r join approval.request q on q.tenant_id=r.tenant_id and q.request_id=r.approval_id
   where r.tenant_id=$1 and r.draft_id=$2 and r.stage<>'publication' and q.state='approved'`, p.TenantID, input.DraftID).Scan(&count)
		if err != nil {
			return cr.Receipt{}, err
		}
		if count != 3 {
			return cr.Receipt{}, cr.ErrConflict
		}
	}
	_, err = NewHumanApprovals(s.pool).decideTx(ctx, tx, p, p.TenantID, aid, s.profile.OrganizationID, payloadSHA, input.Approved, input.Reason, permission, nil)
	if err != nil {
		return cr.Receipt{}, err
	}
	_, err = tx.Exec(ctx, `insert into catalog.release_review_action(tenant_id,command_id,approval_id,actor,approved,evidence_sha256) values($1,$2,$3,$4,$5,$6)`,
		p.TenantID, input.CommandID, aid, p.Subject, input.Approved, input.EvidenceSHA256)
	if err != nil {
		return cr.Receipt{}, err
	}
	r := cr.Receipt{CommandID: input.CommandID, Actor: p.Subject, Kind: "review", ResourceID: input.DraftID, RequestSHA256: hash, SnapshotSHA256: input.SnapshotSHA256}
	if err = s.save(ctx, tx, p, r); err != nil {
		return cr.Receipt{}, err
	}
	return r, tx.Commit(ctx)
}
func (s *CatalogRelease) Publish(ctx context.Context, p identity.Principal, input cr.PublishRequest) (cr.Receipt, error) {
	if !input.Valid() {
		return cr.Receipt{}, cr.ErrInvalid
	}
	tx, err := s.begin(ctx, p, "catalog:publish")
	if err != nil {
		return cr.Receipt{}, err
	}
	defer tx.Rollback(ctx)
	_, hash, err := cr.Canonical(input)
	if err != nil {
		return cr.Receipt{}, err
	}
	if r, yes, err := s.replay(ctx, tx, p, input.CommandID, "publish", hash); err != nil || yes {
		return r, err
	}
	var current int64
	if err = tx.QueryRow(ctx, `select coalesce(max(generation),0) from catalog.release_publication where tenant_id=$1`, p.TenantID).Scan(&current); err != nil {
		return cr.Receipt{}, err
	}
	if current != input.ExpectedGeneration || current == math.MaxInt64 {
		return cr.Receipt{}, cr.ErrConflict
	}
	var raw []byte
	var maker, book string
	err = tx.QueryRow(ctx, `select snapshot_canonical,maker,price_book_id from catalog.release_draft where tenant_id=$1 and draft_id=$2 and snapshot_sha256=$3`, p.TenantID, input.DraftID, input.SnapshotSHA256).Scan(&raw, &maker, &book)
	if errors.Is(err, pgx.ErrNoRows) {
		return cr.Receipt{}, cr.ErrNotFound
	}
	if err != nil {
		return cr.Receipt{}, err
	}
	if maker == p.Subject {
		return cr.Receipt{}, cr.ErrConflict
	}
	var approved int
	err = tx.QueryRow(ctx, `select count(*) from catalog.release_review r join approval.request q on q.tenant_id=r.tenant_id and q.request_id=r.approval_id
 where r.tenant_id=$1 and r.draft_id=$2 and q.state='approved'`, p.TenantID, input.DraftID).Scan(&approved)
	if err != nil {
		return cr.Receipt{}, err
	}
	if approved != 4 {
		return cr.Receipt{}, cr.ErrConflict
	}
	var eligible bool
	err = tx.QueryRow(ctx, `select valid_from<=clock_timestamp() and (valid_until is null or valid_until>clock_timestamp())
 from pricing.price_book where tenant_id=$1 and price_book_id=$2 for update`, p.TenantID, book).Scan(&eligible)
	if err != nil {
		return cr.Receipt{}, err
	}
	if !eligible {
		return cr.Receipt{}, cr.ErrConflict
	}
	var snapshot cr.Snapshot
	if err = json.Unmarshal(raw, &snapshot); err != nil {
		return cr.Receipt{}, err
	}
	// Every publication gets a new immutable effective book. The source owner
	// permits one activation per book; rollback clones the approved historical
	// amounts and validity without rewriting either history or source events.
	effective := commerce.PriceBook{ID: (randomid.Generator{}).New(), Market: snapshot.Profile.Market, Currency: snapshot.Profile.Currency, ValidFrom: snapshot.ValidFrom, ValidUntil: snapshot.ValidUntil}
	for _, v := range snapshot.Variants {
		effective.Entries = append(effective.Entries, commerce.PriceEntry{VariantID: v.ID, AmountMinorUnits: v.AmountMinorUnits, TaxMode: v.TaxMode})
	}
	if err = s.price.createPriceBookTx(ctx, tx, p.TenantID, (randomid.Generator{}).New(), effective); err != nil {
		return cr.Receipt{}, err
	}
	r := cr.Receipt{CommandID: input.CommandID, Actor: p.Subject, Kind: "publish", EffectivePriceBookID: effective.ID, ResourceID: input.DraftID, RequestSHA256: hash, SnapshotSHA256: input.SnapshotSHA256, Generation: current + 1}
	_, err = tx.Exec(ctx, `insert into catalog.release_publication(tenant_id,generation,draft_id,command_id,actor,reason,effective_price_book_id) values($1,$2,$3,$4,$5,$6,$7)`, p.TenantID, r.Generation, input.DraftID, input.CommandID, p.Subject, input.Reason, effective.ID)
	if err != nil {
		return cr.Receipt{}, err
	}
	// Status selection is publication glue; the original Commerce activation
	// still owns market/currency overlap, activation SQL and its outbox event.
	_, err = tx.Exec(ctx, `update pricing.price_book set status='retired' where tenant_id=$1 and market=$2 and currency=$3 and status='active'`, p.TenantID, s.profile.Market, s.profile.Currency)
	if err != nil {
		return cr.Receipt{}, err
	}
	if err = s.price.activatePriceBookTx(ctx, tx, p.TenantID, effective.ID, (randomid.Generator{}).New()); err != nil {
		return cr.Receipt{}, err
	}
	if _, err = tx.Exec(ctx, `delete from search.document where tenant_id=$1 and kind='catalog-model'`, p.TenantID); err != nil {
		return cr.Receipt{}, err
	}
	for _, m := range snapshot.Models {
		doc := search.Document{TenantID: p.TenantID, ID: "catalog:" + m.ID, Kind: "catalog-model", ExternalID: m.ID, Title: m.DisplayName, Body: m.CanonicalURL, ContentSHA256: input.SnapshotSHA256}
		if err = doc.Validate(); err != nil {
			return cr.Receipt{}, err
		}
		_, err = tx.Exec(ctx, `insert into search.document(tenant_id,document_id,kind,external_id,title,body,content_sha256) values($1,$2,$3,$4,$5,$6,$7)`,
			doc.TenantID, doc.ID, doc.Kind, doc.ExternalID, doc.Title, doc.Body, doc.ContentSHA256)
		if err != nil {
			return cr.Receipt{}, err
		}
	}
	if err = s.save(ctx, tx, p, r); err != nil {
		return cr.Receipt{}, err
	}
	return r, tx.Commit(ctx)
}
func (s *CatalogRelease) Public(ctx context.Context) (cr.Publication, error) {
	var out cr.Publication
	err := s.pool.QueryRow(ctx, `select generation,draft_id,snapshot_sha256,snapshot_canonical,effective_price_book_id from catalog.release_public_current where tenant_id=$1`, s.profile.TenantID).Scan(&out.Generation, &out.DraftID, &out.SHA256, &out.Snapshot, &out.EffectivePriceBookID)
	if errors.Is(err, pgx.ErrNoRows) {
		return out, cr.ErrNotFound
	}
	return out, err
}
````

### FILE: `internal/platform/postgres/catalog_release_storefront_integration_test.go`

```yaml
block_id: "GO-CONNECTED-CATALOG-PUBLICATION:file28:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "a33af8957f0b18358967bf6d494ee1e620a8751eada601f7e3cece4c60ac3292"
variables: []
secrets_allowed: false
```

````go
package postgres_test

// AUTHORED connected test: current PG publication -> actual Next -> Chromium.
// Synthetic loopback only; the existing HTTP journey owns all catalog writes.
import (
	"context"
	cr "elite.local/enterprise/internal/catalogrelease"
	"encoding/json"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func catalogStorefrontFixtureRequested() bool { return os.Getenv("ELITE_CATALOG_STOREFRONT") == "1" }
func TestCatalogReleaseStorefrontReference(t *testing.T) {
	if !catalogStorefrontFixtureRequested() {
		t.Skip("explicit catalog storefront fixture required")
	}
	testCatalogReleaseHTTPReference(t, true)
}
func attachCatalogStorefrontFixture(t *testing.T, c *catalogHTTPFixture) {
	t.Helper()
	web, node := os.Getenv("ELITE_WEB_ROOT"), os.Getenv("ELITE_NODE_BIN")
	if !filepath.IsAbs(web) || !filepath.IsAbs(node) {
		t.Fatal("absolute exact web and Node paths required")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	t.Cleanup(cancel)
	listener, e := net.Listen("tcp", "127.0.0.1:0")
	if e != nil {
		t.Fatal(e)
	}
	address := listener.Addr().String()
	listener.Close()
	base := "http://" + address
	env := []string{}
	for _, v := range os.Environ() {
		k := strings.ToUpper(strings.SplitN(v, "=", 2)[0])
		if strings.HasPrefix(k, "ELITE_") || strings.HasPrefix(k, "CATALOG_") || strings.HasPrefix(k, "PUBLIC_") ||
			strings.HasPrefix(k, "ENTERPRISE_") || k == "DATABASE_URL" || k == "TEST_DATABASE_URL" || k == "PAYMENT_CONNECTED_DB_URL" ||
			k == "AUTH_SESSION_SECRET" || k == "APP_BASE_URL" {
			continue
		}
		env = append(env, v)
	}
	env = append(env, "CATALOG_RELEASE_ENABLED=true", "PUBLIC_SITE_ORIGIN=https://catalog.example.invalid", "PUBLIC_INDEXING_ENABLED=1",
		"ENTERPRISE_TENANT_CODE=catalog-j3", "ENTERPRISE_ORGANIZATION_CODE=j3-store", "ENTERPRISE_API_BASE_URL="+c.base,
		"APP_BASE_URL="+base, "AUTH_SESSION_SECRET=synthetic-catalog-next-reference-000000000000000",
		"ELITE_WEB_ROOT="+web, "ELITE_BASE_URL="+base, "ELITE_CATALOG_STOREFRONT=1", "NEXT_TELEMETRY_DISABLED=1", "BUSINESS_CONFIG_FILE=business.example.json")
	artifacts, e := os.MkdirTemp(web, "catalog-storefront-artifacts-")
	if e != nil {
		t.Fatal(e)
	}
	log, e := os.Create(filepath.Join(artifacts, "next.log"))
	if e != nil {
		t.Fatal(e)
	}
	t.Cleanup(func() { log.Close() })
	server := exec.CommandContext(ctx, node, filepath.Join(web, "node_modules", "next", "dist", "bin", "next"), "start", "--hostname", "127.0.0.1", "--port", strings.Split(address, ":")[1])
	server.Dir, server.Env, server.Stdout, server.Stderr = web, env, log, log
	if e = server.Start(); e != nil {
		t.Fatal(e)
	}
	t.Cleanup(func() { server.Process.Kill(); server.Wait() })
	client := &http.Client{Timeout: time.Second}
	ready := false
	for deadline := time.Now().Add(25 * time.Second); time.Now().Before(deadline); {
		r, e := client.Get(base + "/icon.svg")
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
		t.Fatal("Next startup failed; see ", artifacts)
	}
	c.webProbe = func(document cr.PublicDocument) {
		raw, e := json.Marshal(document)
		if e != nil {
			t.Fatal(e)
		}
		generation := filepath.Join(artifacts, "generation-"+jsonNumber(document.Generation))
		if e = os.Mkdir(generation, 0700); e != nil {
			t.Fatal(e)
		}
		if e = os.WriteFile(filepath.Join(generation, "publication.json"), raw, 0600); e != nil {
			t.Fatal(e)
		}
		gate := filepath.Join(web, "microsoft_playwright_browser_gate")
		command := exec.CommandContext(ctx, node, filepath.Join(gate, "node_modules", "@playwright", "test", "cli.js"), "test", "tests/catalog-publication-connected.spec.mjs", "--project=chromium-desktop", "--timeout=30000", "--output="+filepath.Join(generation, "browser"))
		command.Dir, command.Env = gate, append(append([]string(nil), env...), "ELITE_CATALOG_EXPECTED="+string(raw))
		output, e := command.CombinedOutput()
		if x := os.WriteFile(filepath.Join(generation, "browser.log"), output, 0600); x != nil {
			t.Fatal(x)
		}
		if e != nil || !strings.Contains(string(output), "1 passed") {
			t.Fatalf("storefront generation%d: %v\n%s\nartifacts=%s", document.Generation, e, output, artifacts)
		}
	}
	t.Cleanup(func() {
		t.Log("CATALOG_STOREFRONT_PASS actual Next/Chromium/PG; published versions and rollback, canonical, sitemap, robots, normalized PNG and absent-model404; artifacts=", artifacts)
	})
}
func jsonNumber(value int64) string { raw, _ := json.Marshal(value); return string(raw) }
````

## 6. Configuration surface

docs/CATALOG_PUBLICATION_REFERENCE.md and synthetic deploy/catalog profiles document exact hash-bound optional activation, permissions and protocol. Missing selected owners/guards/configuration fail closed. No live accounts requested.

## 7. Dependency bill

No new upstream or dependency. AUTHORED configuration, SQL/interface/orchestration/test glue around existing admitted owners. Go1.26.8 PNG/HTTP, PG18.6 and existing source adapters retain exact provenance.

## 8. Apply order

Materialize with MARKDOWN-COMPOSITOR0.3.0 and complete closure into an absent target. Apply migrations73/74 in order. Populated downgrade refuses history; empty74/73down and73/74up proven. Narrow original APIs preserved.

## 9. Verification

Actual PG71migrations,3snapshots/12reviews/5publications,1new7replay,13table rollback, source isolation, expiry/receipt recovery, current search/media, HTTP lost-response GET, feed rejection/unknown/GET reconciliation. Host12base+feed guard/hash/key/token gates; finite2s fuzz121059execs. Next/Chromium same four publication versions. Provider SDK mappings T2805 remain explicit.

## 10. Reconstruction evidence

CATALOG_CONNECTED_RELEASE_V402.md/json binds source, failed and successful receipts, exact affected profile reconstruction and notices. No global readiness or production claim from this one journey.


V402 composed delta: T2804 catalog role source/edit/review/publication transport reuses original model/price writers and catalog owner. Bounded PNG and text defaults retained. No new dependency. CATALOG_ROLE_AUTHORING_RELEASE_V402.md.
