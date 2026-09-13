# Approved stored-value tender connected to checkout and delivery

## 1. Metadata

```yaml
pack_id: "GO-APPROVED-STORED-VALUE-TENDER"
pack_version: "0.2.1"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Typed Go/PG approval, durable funding projection, authenticated HTTP/BFF/operator UI and zero-provider delivery receipt around the source calculator and existing commerce/handover owners"
stacks: ["CPython 3.14.4", "Go 1.26.8", "PostgreSQL 18.6", "Next.js 16.3.4 for selected portal"]
compatible_with: ["MARKDOWN-COMPOSITOR 0.3.0", "V402 connected franchise profile", "explicit hash-bound reference program"]
incompatible_with: ["public bearer-wallet inferred from assisted internal account IDs", "arbitrary financial or tax policy", "production certification inferred from fixtures"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://github.com/odoo/odoo/tree/99edb6dd82b7b560930c00b03b694ba700785370"]
verified_at: "2026-09-12"
```

## 2. Applicability

Select both stored-value packs only with the connected commerce/payment/handover,
shared human-approval and portal owners in FRANCHISE_COMPLETE_PACK_PLAN. A named
reference profile supplies rules explicitly; it does not decide the target program,
tax, wage or accounting policy. No public gift-card code service is claimed.

## 3. Architecture contract

The LGPL calculator owns only the exact mapped commercial calculations. CPython
runs -I -S -B and executes verified source snapshots with no filesystem package
imports. AUTHORED Go glue owns scoped principals, locked order/account snapshots,
immutable approvals, idempotency, outbox and PostgreSQL receipts. Separate reviewers
approve bound effects. Deferred commit constraints reject expired operations after
waits; reservations and order locks prevent double spend. Receipt history and
current eligibility are separate. Full coverage creates an immutable local funding
receipt, never a zero-value provider charge; handover explicitly selects profile3.
BFF transports int64 and exact payload/receipt bytes as strings, bounds body/time,
and recovers persisted request keys after uncertain responses.

## 4. Exact file manifest

```text
CREATE cmd/electromobility-api/stored_value.go
CREATE cmd/electromobility-api/stored_value_test.go
CREATE db/migrations/0067_gift_loyalty.down.sql
CREATE db/migrations/0067_gift_loyalty.up.sql
CREATE deploy/stored-value/profile.reference.json
CREATE docs/stored-value-runtime.md
CREATE internal/platform/httpapi/stored_value.go
CREATE internal/platform/postgres/stored_value.go
CREATE internal/platform/postgres/stored_value_browser_integration_test.go
CREATE internal/platform/postgres/stored_value_funding.go
CREATE internal/platform/postgres/stored_value_http_integration_test.go
CREATE internal/platform/postgres/stored_value_integration_test.go
CREATE internal/platform/postgres/stored_value_read.go
CREATE internal/platform/postgres/stored_value_tender_integration_test.go
CREATE internal/platform/postgres/zero_funding_integration_test.go
CREATE internal/storedvaluebridge/contract.go
CREATE internal/storedvaluebridge/contract_test.go
CREATE internal/storedvaluebridge/fuzz_test.go
CREATE internal/storedvaluebridge/persistence.go
CREATE internal/storedvaluebridge/process.go
CREATE microsoft_playwright_browser_gate/tests/stored-value-connected.spec.mjs
CREATE src/app/api/enterprise/franchise/stored-value/route.test.ts
CREATE src/app/api/enterprise/franchise/stored-value/route.ts
CREATE src/app/franchise/stored-value/page.tsx
CREATE src/components/stored-value-operations.tsx
CREATE src/platform/stored-value/contracts.ts
CREATE internal/platform/httpapi/stored_value_metrics.go
CREATE internal/platform/postgres/stored_value_metrics.go
```

## 5. Materialization blocks

### FILE: `cmd/electromobility-api/stored_value.go`

```yaml
block_id: "GO-APPROVED-STORED-VALUE-TENDER:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration/configuration/authorization/persistence/transport/UI glue; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "4b2aa935ef9449a61946cbf83619039ab469458b4c15349db7788125d50572f6"
variables: []
secrets_allowed: false
```

````go
package main

// AUTHORED optional activation. No route is mounted until profile, source
// manifest and isolated calculation process pass their exact local bindings.
import (
	"context"
	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/platform/postgres"
	sv "elite.local/enterprise/internal/storedvaluebridge"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"path/filepath"
)

var errStoredValueConfiguration = errors.New("stored-value activation is invalid")

func init() { storedValueModuleFactory = selectedStoredValueModule }
func selectedStoredValueModule(ctx context.Context, pool *pgxpool.Pool, lookup func(string) string) (httpapi.EnterpriseModule, error) {
	if lookup == nil {
		return nil, errStoredValueConfiguration
	}
	switch lookup("STORED_VALUE_ENABLED") {
	case "", "false":
		return nil, nil
	case "true":
	default:
		return nil, errStoredValueConfiguration
	}
	if pool == nil {
		return nil, errStoredValueConfiguration
	}
	path := lookup("STORED_VALUE_PROFILE_FILE")
	if !filepath.IsAbs(path) {
		return nil, errStoredValueConfiguration
	}
	f, e := os.Open(path)
	if e != nil {
		return nil, errStoredValueConfiguration
	}
	info, e := f.Stat()
	f.Close()
	if e != nil || !info.Mode().IsRegular() || info.Size() > 65536 {
		return nil, errStoredValueConfiguration
	}
	raw, e := os.ReadFile(path)
	if e != nil {
		return nil, errStoredValueConfiguration
	}
	profile, e := sv.Load(raw, lookup("STORED_VALUE_PROFILE_SHA256"))
	if e != nil {
		return nil, errStoredValueConfiguration
	}
	tenant, org := profile.Scope()
	if tenant != lookup("STORED_VALUE_TENANT_ID") || org != lookup("STORED_VALUE_ORGANIZATION_ID") {
		return nil, errStoredValueConfiguration
	}
	process := sv.Process{Python: lookup("STORED_VALUE_PYTHON"), Script: lookup("STORED_VALUE_SCRIPT"), ScriptSHA256: lookup("STORED_VALUE_SCRIPT_SHA256"), Manifest: lookup("STORED_VALUE_MANIFEST"), ManifestSHA256: lookup("STORED_VALUE_MANIFEST_SHA256")}
	if process.Preflight(ctx, profile) != nil {
		return nil, errStoredValueConfiguration
	}
	store, e := postgres.NewStoredValue(pool, profile, process)
	if e != nil {
		return nil, errStoredValueConfiguration
	}
	return httpapi.StoredValueModule{Service: store}, nil
}
````

### FILE: `cmd/electromobility-api/stored_value_test.go`

```yaml
block_id: "GO-APPROVED-STORED-VALUE-TENDER:file2:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration/configuration/authorization/persistence/transport/UI glue; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "6895ec73166950f4b95ffd199b1060332e089682b29fb91dde35d422d6a3e9a8"
variables: []
secrets_allowed: false
```

````go
package main

// AUTHORED activation contract. Pool is lazy and no database is contacted here;
// the connected PG/browser suite separately verifies every mounted operation.
import (
	"context"
	sv "elite.local/enterprise/internal/storedvaluebridge"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestStoredValueOptionalActivation(t *testing.T) {
	ctx := context.Background()
	root, e := filepath.Abs("../..")
	if e != nil {
		t.Fatal(e)
	}
	read := func(rel string) []byte {
		t.Helper()
		b, e := os.ReadFile(filepath.Join(root, rel))
		if e != nil {
			t.Fatal(e)
		}
		return b
	}
	config := map[string]string{"STORED_VALUE_ENABLED": "true", "STORED_VALUE_PROFILE_FILE": filepath.Join(root, "deploy/stored-value/profile.reference.json"), "STORED_VALUE_PROFILE_SHA256": sv.Hash(read("deploy/stored-value/profile.reference.json")), "STORED_VALUE_TENANT_ID": "11111111-1111-4111-8111-111111111111", "STORED_VALUE_ORGANIZATION_ID": "franchise-1", "STORED_VALUE_PYTHON": os.Getenv("HANDOVER_PROFILE_PYTHON"), "STORED_VALUE_SCRIPT": filepath.Join(root, "odoo_loyalty/run.py"), "STORED_VALUE_SCRIPT_SHA256": sv.Hash(read("odoo_loyalty/run.py")), "STORED_VALUE_MANIFEST": filepath.Join(root, "odoo_loyalty/engine-lock.json"), "STORED_VALUE_MANIFEST_SHA256": sv.Hash(read("odoo_loyalty/engine-lock.json"))}
	if config["STORED_VALUE_PYTHON"] == "" {
		t.Fatal("explicit fixed Python required for source preflight")
	}
	pool, e := pgxpool.New(ctx, "postgres://fixture@127.0.0.1:1/fixture?sslmode=disable&connect_timeout=1")
	if e != nil {
		t.Fatal(e)
	}
	defer pool.Close()
	for _, which := range []string{"valid", "disabled", "invalid-flag", "profile-sha", "organization", "script-sha", "manifest-sha", "relative-profile"} {
		t.Run(which, func(t *testing.T) {
			v := map[string]string{}
			for k, val := range config {
				v[k] = val
			}
			switch which {
			case "disabled":
				v["STORED_VALUE_ENABLED"] = "false"
			case "invalid-flag":
				v["STORED_VALUE_ENABLED"] = "1"
			case "profile-sha":
				v["STORED_VALUE_PROFILE_SHA256"] = "bad"
			case "organization":
				v["STORED_VALUE_ORGANIZATION_ID"] = "other"
			case "script-sha":
				v["STORED_VALUE_SCRIPT_SHA256"] = "bad"
			case "manifest-sha":
				v["STORED_VALUE_MANIFEST_SHA256"] = "bad"
			case "relative-profile":
				v["STORED_VALUE_PROFILE_FILE"] = "profile.json"
			}
			module, e := storedValueModuleFactory(ctx, pool, func(k string) string { return v[k] })
			if which == "disabled" {
				if module != nil || e != nil {
					t.Fatal("disabled module mounted", e)
				}
				return
			}
			if which != "valid" {
				if e == nil || module != nil {
					t.Fatal("invalid activation admitted")
				}
				return
			}
			if e != nil || module == nil {
				t.Fatal("valid activation", e)
			}
			mux := http.NewServeMux()
			module.Register(mux, nil)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, httptest.NewRequest("GET", "/v1/franchise/stored-value/profile?organization_id=franchise-1", nil))
			if w.Code != 401 {
				t.Fatal("route was not mounted behind authentication", w.Code)
			}
		})
	}
}
````

### FILE: `db/migrations/0067_gift_loyalty.down.sql`

```yaml
block_id: "GO-APPROVED-STORED-VALUE-TENDER:file3:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration/configuration/authorization/persistence/transport/UI glue; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "fa74b9016e5570428a76c5662611b4ae540b32d83547520537695bcd4d280748"
variables: []
secrets_allowed: false
```

````sql
begin;
do $$begin
 if exists(select 1 from payment.local_funding_receipt) or exists(select 1 from stored_value.operation) or exists(select 1 from stored_value.account) or exists(select 1 from stored_value.reservation) or exists(select 1 from approval.request where kind='stored_value_operation') then
  raise exception 'cannot remove retained stored value evidence';
 end if;
end $$;
create or replace view payment.order_funding as select tenant_id,order_id,organization_id,currency,version as order_version,total_minor_units as gross_minor_units,0::bigint as gift_minor_units,0::bigint as discount_minor_units,total_minor_units as provider_due_minor_units,encode(sha256(convert_to('[]','UTF8')),'hex') as contributions_sha256 from sales.customer_order;
-- Remove only owned objects. Unknown dependencies fail the transaction rather
-- than being erased by CASCADE, even when these owned ledgers are empty.
drop trigger stored_value_proposal_lease on approval.request;
drop view stored_value.order_allocation;
drop table stored_value.reservation;
drop table stored_value.entry;
drop table stored_value.operation;
drop table stored_value.account;
drop function stored_value.guard_reservation();
drop function stored_value.immutable_record();
drop function stored_value.require_live_proposal();
drop function stored_value.require_live_operation();
drop schema stored_value;
alter table approval.request drop constraint request_kind_check;
alter table approval.request add constraint request_kind_check check(kind in ('reservation','sale','refund','payment','whatsapp_reply','social_publish','social_revoke'));
alter table approval.request drop constraint request_payload_binding_check;
alter table approval.request add constraint request_payload_binding_check check(kind not in ('whatsapp_reply','social_publish','social_revoke') or (organization_id is not null and length(organization_id) between 1 and 128 and payload is not null and jsonb_typeof(payload)='object' and pg_column_size(payload)<=65536 and amount_minor_units=0));
create or replace function approval.guard_bound_request() returns trigger language plpgsql as $$
begin
 if tg_op='DELETE' then
  if old.kind in ('whatsapp_reply','social_publish','social_revoke') then raise exception 'bound approval request is immutable'; end if;
  return old;
 end if;
 if old.kind in ('whatsapp_reply','social_publish','social_revoke') or new.kind in ('whatsapp_reply','social_publish','social_revoke') then
  if row(new.tenant_id,new.request_id,new.kind,new.subject_id,new.amount_minor_units,new.requester,new.evidence_sha,new.organization_id,new.payload,new.created_at)
   is distinct from row(old.tenant_id,old.request_id,old.kind,old.subject_id,old.amount_minor_units,old.requester,old.evidence_sha,old.organization_id,old.payload,old.created_at)
   or (old.state<>'pending' and row(new.state,new.decided_at) is distinct from row(old.state,old.decided_at))
  then raise exception 'bound approval request is immutable'; end if;
 end if;
 return new;
end $$;
create or replace function approval.guard_bound_decision() returns trigger language plpgsql as $$
begin
 if exists(select 1 from approval.request where tenant_id=old.tenant_id and request_id=old.request_id and kind in ('whatsapp_reply','social_publish','social_revoke')) then
  raise exception 'bound approval decision is immutable';
 end if;
 if tg_op='DELETE' then return old;end if;return new;
end $$;
commit;
````

### FILE: `db/migrations/0067_gift_loyalty.up.sql`

```yaml
block_id: "GO-APPROVED-STORED-VALUE-TENDER:file4:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration/configuration/authorization/persistence/transport/UI glue; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "659fdf7ad9785d515c610a4d360bc5c28e2e292dd61282e2cc17cd4dd8b639e0"
variables: []
secrets_allowed: false
```

````sql
begin;
-- AUTHORED durable representation, reservation and receipt glue. No financial
-- general ledger or external-provider payment is synthesized here.
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

create schema stored_value;
create table stored_value.account (
 tenant_id uuid not null,
 account_id text not null check(length(account_id) between 1 and 128),
 organization_id text not null,
 program_id text not null,
 kind text not null check(kind in ('gift_card','loyalty')),
 customer_subject text,
 currency text not null check(currency ~ '^[A-Z]{3}$'),
 created_at timestamptz not null default clock_timestamp(),
 primary key(tenant_id,account_id),
 foreign key(tenant_id,organization_id) references org.organization(tenant_id,organization_id),
 check((kind='loyalty' and customer_subject is not null) or kind='gift_card')
);
create table stored_value.operation (
 tenant_id uuid not null,
 operation_id text not null check(length(operation_id) between 1 and 128),
 organization_id text not null,
 order_id text not null,
 order_version bigint not null check(order_version>0),
 program_id text not null,
 operation text not null check(operation in ('issue','accrue','redeem','reverse')),
 original_operation_id text,
 approval_id text not null,
 profile_sha256 text not null check(profile_sha256 ~ '^[0-9a-f]{64}$'),
 request_sha256 text not null check(request_sha256 ~ '^[0-9a-f]{64}$'),
 calculation_sha256 text not null check(calculation_sha256 ~ '^[0-9a-f]{64}$'),
 receipt_sha256 text not null check(receipt_sha256 ~ '^[0-9a-f]{64}$'),
 receipt jsonb not null check(jsonb_typeof(receipt)='object' and pg_column_size(receipt)<=65536),
 created_at timestamptz not null default clock_timestamp(),
 primary key(tenant_id,operation_id),
 foreign key(tenant_id,organization_id) references org.organization(tenant_id,organization_id),
 foreign key(tenant_id,order_id) references sales.customer_order(tenant_id,order_id),
 foreign key(tenant_id,approval_id) references approval.request(tenant_id,request_id),
 foreign key(tenant_id,original_operation_id) references stored_value.operation(tenant_id,operation_id),
 unique(tenant_id,approval_id),
 check((operation='reverse')=(original_operation_id is not null))
);
create unique index one_stored_value_reversal on stored_value.operation(tenant_id,original_operation_id) where original_operation_id is not null;
create index stored_value_order_operations on stored_value.operation(tenant_id,order_id,program_id,created_at);
create table stored_value.entry (
 tenant_id uuid not null,
 operation_id text not null,
 account_id text not null,
 points_delta numeric(38,6) not null check(points_delta::text not in ('NaN','Infinity','-Infinity')),
 applied_minor_units bigint not null default 0,
 primary key(tenant_id,operation_id,account_id),
 foreign key(tenant_id,operation_id) references stored_value.operation(tenant_id,operation_id),
 foreign key(tenant_id,account_id) references stored_value.account(tenant_id,account_id),
 check(points_delta<>0 or applied_minor_units<>0)
);
create index stored_value_account_entries on stored_value.entry(tenant_id,account_id,operation_id);
create table stored_value.reservation (
 tenant_id uuid not null,
 operation_id text not null,
 account_id text not null,
 order_id text not null,
 approval_id text not null,
 points numeric(38,6) not null check(points>0 and points::text not in ('NaN','Infinity','-Infinity')),
 state text not null check(state in ('reserved','consumed','rejected','expired')),
 expires_at timestamptz not null,
 created_at timestamptz not null default clock_timestamp(),
 primary key(tenant_id,operation_id,account_id),
 foreign key(tenant_id,account_id) references stored_value.account(tenant_id,account_id),
 foreign key(tenant_id,order_id) references sales.customer_order(tenant_id,order_id),
 foreign key(tenant_id,approval_id) references approval.request(tenant_id,request_id) deferrable initially deferred,
 check(expires_at>created_at)
);
create index stored_value_live_reservations on stored_value.reservation(tenant_id,account_id,expires_at) where state='reserved';
create view stored_value.order_allocation as
 select o.tenant_id,o.order_id,o.organization_id,o.currency,o.version as order_version,o.total_minor_units as gross_minor_units,
 coalesce(a.gift_minor,0)::bigint as gift_minor_units,coalesce(a.discount_minor,0)::bigint as discount_minor_units,
 (o.total_minor_units-coalesce(a.gift_minor,0)-coalesce(a.discount_minor,0))::bigint as provider_due_minor_units
 from sales.customer_order o left join (
  select op.tenant_id,op.order_id,
  sum(e.applied_minor_units) filter(where a.kind='gift_card') as gift_minor,
  sum(e.applied_minor_units) filter(where a.kind='loyalty') as discount_minor
  from stored_value.operation op join stored_value.entry e using(tenant_id,operation_id)
  join stored_value.account a using(tenant_id,account_id)
  group by op.tenant_id,op.order_id
 ) a using(tenant_id,order_id);
create or replace view payment.order_funding as
 select a.tenant_id,a.order_id,a.organization_id,a.currency,a.order_version,a.gross_minor_units,a.gift_minor_units,a.discount_minor_units,a.provider_due_minor_units,
 encode(sha256(convert_to(coalesce((select jsonb_agg(jsonb_build_array(o.operation_id,o.receipt_sha256) order by o.operation_id)
 from stored_value.operation o where o.tenant_id=a.tenant_id and o.order_id=a.order_id
 and exists(select 1 from stored_value.entry e where e.tenant_id=o.tenant_id and e.operation_id=o.operation_id and e.applied_minor_units<>0)),'[]'::jsonb)::text,'UTF8')),'hex') as contributions_sha256
 from stored_value.order_allocation a;
create function stored_value.immutable_record() returns trigger language plpgsql as $$
begin raise exception 'stored value evidence is immutable; record a bound reversal';end $$;
create trigger stored_value_account_immutable before update or delete on stored_value.account for each row execute function stored_value.immutable_record();
create trigger stored_value_operation_immutable before update or delete on stored_value.operation for each row execute function stored_value.immutable_record();
create trigger stored_value_entry_immutable before update or delete on stored_value.entry for each row execute function stored_value.immutable_record();
create function stored_value.guard_reservation() returns trigger language plpgsql as $$
begin
 if tg_op='DELETE' then raise exception 'reservation evidence is immutable';end if;
 if row(new.tenant_id,new.operation_id,new.account_id,new.order_id,new.approval_id,new.points,new.expires_at,new.created_at)
  is distinct from row(old.tenant_id,old.operation_id,old.account_id,old.order_id,old.approval_id,old.points,old.expires_at,old.created_at)
  or (old.state<>'reserved' and new.state<>old.state)
 then raise exception 'reservation identity or terminal state is immutable';end if;
 return new;
end $$;
create trigger stored_value_reservation_immutable before update or delete on stored_value.reservation for each row execute function stored_value.guard_reservation();
-- A lock wait, slow IPC or outbox write must not turn an expired reviewed
-- snapshot into a committed effect. Deferred constraints execute at COMMIT,
-- after every statement, using the database clock rather than transaction time.
create function stored_value.require_live_proposal() returns trigger language plpgsql as $$
begin
 if new.kind='stored_value_operation' and
  not coalesce((new.payload->>'expires_at')::timestamptz>clock_timestamp(),false) then
  raise exception using errcode='23514', constraint='stored_value_lease_valid', message='stored-value proposal expired before commit';
 end if;
 return null;
end $$;
create constraint trigger stored_value_proposal_lease after insert on approval.request
 deferrable initially deferred for each row execute function stored_value.require_live_proposal();
create function stored_value.require_live_operation() returns trigger language plpgsql as $$
begin
 if not coalesce((new.receipt->'operation'->>'expires_at')::timestamptz>clock_timestamp(),false) then
  raise exception using errcode='23514', constraint='stored_value_lease_valid', message='stored-value operation expired before commit';
 end if;
 return null;
end $$;
create constraint trigger stored_value_operation_lease after insert on stored_value.operation
 deferrable initially deferred for each row execute function stored_value.require_live_operation();
commit;
````

### FILE: `deploy/stored-value/profile.reference.json`

```yaml
block_id: "GO-APPROVED-STORED-VALUE-TENDER:file5:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration/configuration/authorization/persistence/transport/UI glue; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "db96a3558474e1dd57448d4fd25af9753cdef7238faabac13b1f17727994e727"
variables: []
secrets_allowed: false
```

````json
{
  "schema": "elite.stored-value-profile.v1",
  "tenant_id": "11111111-1111-4111-8111-111111111111",
  "organization_id": "franchise-1",
  "amounts_mode": "bound_gross_inclusive",
  "review": "one_distinct_human",
  "issuance_states": [
    "confirmed",
    "paid",
    "allocated",
    "delivered"
  ],
  "redemption_states": [
    "placed",
    "confirmed",
    "allocated"
  ],
  "programs": [
    {
      "calculation": {
        "id": "reference-gift_card",
        "kind": "gift_card",
        "currency": "ARS",
        "currency_digits": 2,
        "applies_on": "future",
        "nominative": false,
        "trigger": "auto",
        "trigger_product_ids": [
          "gift-50"
        ],
        "rules": [
          {
            "id": "reference-gift_card-rule",
            "mode": "money",
            "points": "1",
            "split": true,
            "minimum_qty": "1",
            "minimum_amount": "0",
            "tax_mode": "incl",
            "product_ids": [
              "gift-50"
            ],
            "code_required": false
          }
        ]
      },
      "reward": {
        "mode": "per_point",
        "discount": "1",
        "required_points": "1",
        "max_amount": "0",
        "clear_wallet": false
      }
    },
    {
      "calculation": {
        "id": "reference-loyalty",
        "kind": "loyalty",
        "currency": "ARS",
        "currency_digits": 2,
        "applies_on": "both",
        "nominative": true,
        "trigger": "auto",
        "trigger_product_ids": [],
        "rules": [
          {
            "id": "reference-loyalty-rule",
            "mode": "money",
            "points": "1",
            "split": false,
            "minimum_qty": "1",
            "minimum_amount": "0",
            "tax_mode": "incl",
            "product_ids": [],
            "code_required": false
          }
        ]
      },
      "reward": {
        "mode": "percent",
        "discount": "5",
        "required_points": "200",
        "max_amount": "0",
        "clear_wallet": false
      }
    }
  ],
  "reservation_seconds": 900
}
````

### FILE: `docs/stored-value-runtime.md`

```yaml
block_id: "GO-APPROVED-STORED-VALUE-TENDER:file6:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration/configuration/authorization/persistence/transport/UI glue; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "50e8c50d6abc27fca122130a75c98728c4b3045dfd5ec2123d1b7488b5b21b26"
variables: []
secrets_allowed: false
```

````markdown
# Gift cards and loyalty in the local franchise reference

This optional composition joins the separately source-derived calculation module to existing order, human approval, payment, delivery and commercial-release owners. It is an assisted operator workflow with a distinct reviewer, not a public bearer-card wallet or a full ERP. The reference policy is explicitly selected, immutable and hash-bound: ARS, two monetary digits, inclusive bound order totals, gift issuance at selected confirmed-sale states, nominated loyalty and one distinct reviewer. Different fiscal/currency/program policies require a reviewed profile or a source-admitted extension. The reference profile is a fixture choice, not an invented live rule for a future business.

## Startup and configuration

Select both the Python calculator and Go approved stored-value packs in the complete franchise composition. Apply migrations0066,0067 and0068 after the existing payment/approval/handover migrations; never edit historical migrations. Verify the composition and docs/provenance/ODOO_STORED_VALUE_RUNTIME_LOCK.json before activation. The runtime requires an admitted absolute Python executable and starts it with -I -S -B, a5second deadline, bounded IPC and a cleared environment. Source bytes are checked and compiled from their verified snapshots, without filesystem package or cached-bytecode imports.

The host's disabled default mounts no stored-value routes. To activate, configure STORED_VALUE_ENABLED=true, STORED_VALUE_PROFILE_FILE and PROFILE_SHA256, STORED_VALUE_TENANT_ID and ORGANIZATION_ID, STORED_VALUE_PYTHON, STORED_VALUE_SCRIPT and SCRIPT_SHA256, and STORED_VALUE_MANIFEST and MANIFEST_SHA256; all names use the STORED_VALUE_ prefix. The selected profile and source lock supply expected hashes. Paths must be absolute. The host runs the actual source preflight and refuses mismatches. No account or credential is needed for this local module. The existing OIDC host verifies identity; its credentials remain separate from this module's source/profile configuration.

Roles need only their assigned stored_value:read, stored_value:request, stored_value:approve or stored_value:fund permissions and the selected organization. Browser session identity never becomes a request field. The same subject cannot approve its own request. Portal /franchise/stored-value provides order lookup, selected program, proposal, exact amounts/points/accounts, review/reject, receipt and recovery. Account identifiers are internal operator references; they are not public redemption secrets. Approval lists are limited50per page and scoped to an order, with a stable next cursor.

## Connected transitions

Issue/accrue derives points from the existing confirmed source order. Propose redeems against the current order version and holds a bounded reservation. Decision re-reads the locked order/accounts, recomputes the exact approved payload and atomically records immutable operation/entries, reservation consumption and outbox. Deferred constraints require the proposal/operation lease to remain valid at COMMIT, including time spent waiting for the outbox. Expired or changed snapshots apply nothing. Rejected/expired reservations cease to consume availability without deleting evidence.

payment.order_funding preserves gross order value and derives gift/discount/provider remainder. Existing payment creation and dispatch serialize on the same order and validate the remainder after waiting for the lock. Partially funded purchases use the existing official payment SDK, callback and reconciliation. A fully funded purchase creates its own payment.local_funding_receipt observation and outbox, not a zero-amount or fabricated provider capture. Handover algorithm revision3/FULL_ORDER_WITH_STORED_VALUE explicitly admits that funding type. Revisions1/2 retain their original provider-only meaning. Initial handover and commercial receipt preserve payment XOR local funding, real FKs and current-evidence checks.

A cancellation reversal requires the authoritative order/payment states and a separate reviewed inverse. It preserves source-calculated point deltas and immutable history, including negative balances if an already spent issuance is reversed. A canceled/held/refunded order invalidates present eligibility; it does not delete historical funding or release receipts. This does not add a post-delivery cancellation endpoint or an automatic refund/credit policy.

## Recovery, operations and verification

Retain the request before sending. Propose recovers by stable approval ID; decision uses that ID and the reviewed hash; funding uses its original request key. An uncertain response means consult GET or retry exactly that saved command, never create a new request identity to guess the result. The UI persists this marker under tenant/actor/org scope and removes it only after a verified result. Current funding and release eligibility remain distinct from historical receipts. A fresh funding observation may be requested after a prior observation ages, using unchanged contributions and a new explicitly recorded key.

The durable sources of truth are approval.request/decision, stored_value.account/operation/entry/reservation, payment.local_funding_receipt, existing delivery/commercial records and platform.outbox_event. Audit includes requester, reviewer, profile/source/payload hashes and exact result. Do not log document payloads or put arbitrary customer PII in operation references. Operation/account history is immutable; retention and any permitted pseudonymization belong to the target's configured financial/privacy owners. No automatic deletion schedule is invented here. Back up/restore these tables together using the composition's PostgreSQL procedures.

Rollback application deployment without erasing evidence. All three down migrations refuse retained funding/operation/delivery rows;0067 removes only known objects and refuses unknown dependencies. Empty down/up is tested. Preserve source/profile snapshots for historical receipts; a changed source result cannot silently approve an old hash. Source failure should leave the saved proposal available for inspection, with no partial entry/event.

Reproducible evidence in the canonical library: STORED_VALUE_TENDER_V402.md, STORED_VALUE_OPERATOR_FLOW_V402.md and the source-admission successor. It includes source-oracle comparisons, tenant/order/account isolation, concurrent reservation/funding, lease/outbox rollback, payment SDK/reconciliation, real Chromium/Next/RS256/JWKS/PG with separate reviewer, lost-response recovery and zero-provider delivery. The source loader and replacement tests add manifest/cached-code closure. Native Go fuzz covers exact decimal representation and profile immutability with a finite budget; it is not SAST, DAST or a target-load claim. Global composition SCA/identity, deployment, load/recovery and live acceptance remain their separate library/target gates.
````

### FILE: `internal/platform/httpapi/stored_value.go`

```yaml
block_id: "GO-APPROVED-STORED-VALUE-TENDER:file7:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration/configuration/authorization/persistence/transport/UI glue; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "b8833060dbc82f1b7a3f8f24fbdab0075afe29bad271907ca240815e05515b4d"
variables: []
secrets_allowed: false
```

````go
package httpapi

// AUTHORED monetary-safe transport only. Identity and scope come from verified
// bearer claims; exact source calculations and durable effects stay with owners.
import (
	"context"
	"elite.local/enterprise/internal/approval"
	"elite.local/enterprise/internal/platform/identity"
	db "elite.local/enterprise/internal/platform/postgres"
	sv "elite.local/enterprise/internal/storedvaluebridge"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"net/http"
	"strconv"
)

type StoredValueService interface {
	Metrics(context.Context, identity.Principal, string, string) (db.StoredValueMetrics, error)
	Scope() (string, string)
	ProfileDocument(context.Context, identity.Principal, string) (json.RawMessage, error)
	ProfileSHA256() string
	Approvals(context.Context, identity.Principal, string, string) (db.StoredValueApprovalPage, error)
	Propose(context.Context, identity.Principal, sv.Request) (sv.ApprovalView, error)
	Read(context.Context, identity.Principal, string) (sv.ApprovalView, error)
	Decide(context.Context, identity.Principal, sv.Decision) (sv.ApprovalView, error)
	Allocation(context.Context, identity.Principal, string) (sv.Allocation, error)
	FinalizeFunding(context.Context, identity.Principal, db.FinalizeStoredValueFunding) (db.LocalFundingResult, bool, error)
	FundingResult(context.Context, identity.Principal, string) (db.LocalFundingResult, error)
}
type StoredValueModule struct{ Service StoredValueService }

func storedValueVersion(raw string) (int64, error) {
	n, e := strconv.ParseInt(raw, 10, 64)
	if e != nil || n < 1 || strconv.FormatInt(n, 10) != raw {
		return 0, sv.ErrBinding
	}
	return n, nil
}
func storedValueError(w http.ResponseWriter, e error) bool {
	if e == nil {
		return false
	}
	var pg *pgconn.PgError
	switch {
	case errors.Is(e, approval.ErrNotFound) || errors.Is(e, pgx.ErrNoRows):
		writeProblem(w, 404, "STORED_VALUE_NOT_FOUND", "no result exists in this authorized scope")
	case errors.Is(e, sv.ErrBinding) || errors.Is(e, approval.ErrDuplicate) || errors.Is(e, approval.ErrNotPending) || errors.Is(e, approval.ErrSeparation) || errors.Is(e, approval.ErrInvalidRequest) || errors.As(e, &pg) && (pg.Code == "23514" || pg.Code == "23505"):
		writeProblem(w, 409, "STORED_VALUE_BINDING_REJECTED", "consult current order and approval state before retrying")
	default:
		writeProblem(w, 503, "STORED_VALUE_UNCONFIRMED", "consult the saved request reference; the result is not confirmed")
	}
	return true
}
func storedValueApprovalJSON(w http.ResponseWriter, status int, v sv.ApprovalView) {
	raw, hash, e := sv.Canonical(v.Payload)
	if e != nil || hash != v.PayloadSHA256 {
		storedValueError(w, sv.ErrBinding)
		return
	}
	entries := make([]map[string]string, 0, len(v.Payload.Entries))
	for _, e := range v.Payload.Entries {
		entries = append(entries, map[string]string{"account_id": e.AccountID, "points_delta": e.PointsDelta, "applied_minor_units": sv.Number(e.AppliedMinor)})
	}
	review := map[string]any{"operation_id": v.Payload.Request.OperationID, "program_id": v.Payload.Request.ProgramID, "original_operation_id": v.Payload.Request.OriginalOperationID, "currency": v.Payload.Order.Currency, "gross_minor_units": sv.Number(v.Payload.GrossMinor), "gift_minor_units": sv.Number(v.Payload.GiftMinor), "discount_minor_units": sv.Number(v.Payload.DiscountMinor), "entries": entries}
	// JSON document strings preserve signed int64 amounts across JavaScript.
	writeJSON(w, status, map[string]any{"review": review, "approval_id": v.ApprovalID, "payload_sha256": hash, "state": v.State, "replay": v.Replay, "payload_json": string(raw), "receipt_json": string(v.Receipt), "order_id": v.Payload.Request.OrderID, "organization_id": v.Payload.Request.OrganizationID, "operation": v.Payload.Request.Operation, "requester": v.Payload.Requester, "expires_at": v.Payload.ExpiresAt})
}
func storedValueFundingJSON(w http.ResponseWriter, status int, v db.LocalFundingResult) {
	raw, e := json.Marshal(v.Receipt)
	if e != nil {
		storedValueError(w, e)
		return
	}
	writeJSON(w, status, map[string]any{"funding_receipt_id": v.Receipt.ID, "receipt_sha256": v.SHA256, "receipt_json": string(raw), "organization_id": v.Receipt.Allocation.OrganizationID, "order_id": v.Receipt.Allocation.OrderID, "request_key": v.Receipt.RequestKey})
}
func (m StoredValueModule) Register(mux *http.ServeMux, verifier identity.Verifier) {
	if m.Service == nil {
		return
	}
	api := franchiseJourneyAPI{verifier: verifier}
	protected := func(w http.ResponseWriter, r *http.Request, permission, org string) (identity.Principal, bool) {
		p, ok := api.protected(w, r, permission, org)
		if !ok {
			return p, false
		}
		tenant, bound := m.Service.Scope()
		if p.TenantID != tenant || org != bound {
			writeProblem(w, 403, "FORBIDDEN", "the selected program is outside this request scope")
			return p, false
		}
		return p, true
	}

	registerStoredValueMetrics(m, mux, protected)
	mux.HandleFunc("GET /v1/franchise/stored-value/profile", func(w http.ResponseWriter, r *http.Request) {
		org := r.URL.Query().Get("organization_id")
		p, ok := protected(w, r, "stored_value:read", org)
		if !ok {
			return
		}
		raw, e := m.Service.ProfileDocument(r.Context(), p, org)
		if storedValueError(w, e) {
			return
		}
		writeJSON(w, 200, map[string]string{"profile_sha256": m.Service.ProfileSHA256(), "profile_json": string(raw), "organization_id": org})
	})
	mux.HandleFunc("GET /v1/franchise/stored-value/orders/{id}", func(w http.ResponseWriter, r *http.Request) {
		org := r.URL.Query().Get("organization_id")
		p, ok := protected(w, r, "stored_value:read", org)
		if !ok {
			return
		}
		a, e := m.Service.Allocation(r.Context(), p, r.PathValue("id"))
		if storedValueError(w, e) {
			return
		}
		page, e := m.Service.Approvals(r.Context(), p, r.PathValue("id"), r.URL.Query().Get("after"))
		if storedValueError(w, e) {
			return
		}
		writeJSON(w, 200, map[string]any{"organization_id": org, "order_id": a.OrderID, "order_version": sv.Number(a.OrderVersion), "currency": a.Currency, "gross_minor_units": sv.Number(a.GrossMinor), "gift_minor_units": sv.Number(a.GiftMinor), "discount_minor_units": sv.Number(a.DiscountMinor), "provider_due_minor_units": sv.Number(a.ProviderDueMinor), "approvals": page})
	})
	mux.HandleFunc("POST /v1/franchise/stored-value/proposals", func(w http.ResponseWriter, r *http.Request) {
		var c struct {
			OperationID          string `json:"operation_id"`
			Operation            string `json:"operation"`
			OrganizationID       string `json:"organization_id"`
			ProgramID            string `json:"program_id"`
			OrderID              string `json:"order_id"`
			ExpectedOrderVersion string `json:"expected_order_version"`
			AccountID            string `json:"account_id"`
			OriginalOperationID  string `json:"original_operation_id"`
			ProfileSHA256        string `json:"profile_sha256"`
		}
		if !decodeStrict(w, r, &c) {
			return
		}
		p, ok := protected(w, r, "stored_value:request", c.OrganizationID)
		if !ok {
			return
		}
		n, e := storedValueVersion(c.ExpectedOrderVersion)
		if storedValueError(w, e) {
			return
		}
		v, e := m.Service.Propose(r.Context(), p, sv.Request{OperationID: c.OperationID, Operation: c.Operation, OrganizationID: c.OrganizationID, ProgramID: c.ProgramID, OrderID: c.OrderID, ExpectedOrderVersion: n, AccountID: c.AccountID, OriginalOperationID: c.OriginalOperationID, ProfileSHA256: c.ProfileSHA256})
		if storedValueError(w, e) {
			return
		}
		storedValueApprovalJSON(w, 201, v)
	})
	mux.HandleFunc("GET /v1/franchise/stored-value/approvals/{id}", func(w http.ResponseWriter, r *http.Request) {
		p, ok := protected(w, r, "stored_value:read", r.URL.Query().Get("organization_id"))
		if !ok {
			return
		}
		v, e := m.Service.Read(r.Context(), p, r.PathValue("id"))
		if storedValueError(w, e) {
			return
		}
		storedValueApprovalJSON(w, 200, v)
	})
	mux.HandleFunc("POST /v1/franchise/stored-value/approvals/{id}/decision", func(w http.ResponseWriter, r *http.Request) {
		var c struct {
			OrganizationID string `json:"organization_id"`
			PayloadSHA256  string `json:"payload_sha256"`
			Approved       *bool  `json:"approved"`
			Reason         string `json:"reason"`
		}
		if !decodeStrict(w, r, &c) {
			return
		}
		p, ok := protected(w, r, "stored_value:approve", c.OrganizationID)
		if !ok {
			return
		}
		if c.Approved == nil {
			storedValueError(w, sv.ErrBinding)
			return
		}
		v, e := m.Service.Decide(r.Context(), p, sv.Decision{OrganizationID: c.OrganizationID, ApprovalID: r.PathValue("id"), PayloadSHA256: c.PayloadSHA256, Approved: *c.Approved, Reason: c.Reason})
		if storedValueError(w, e) {
			return
		}
		storedValueApprovalJSON(w, 200, v)
	})
	mux.HandleFunc("POST /v1/franchise/stored-value/orders/{id}/funding", func(w http.ResponseWriter, r *http.Request) {
		var c struct {
			OrganizationID       string `json:"organization_id"`
			ExpectedOrderVersion string `json:"expected_order_version"`
		}
		if !decodeStrict(w, r, &c) {
			return
		}
		p, ok := protected(w, r, "stored_value:fund", c.OrganizationID)
		if !ok {
			return
		}
		n, e := storedValueVersion(c.ExpectedOrderVersion)
		if storedValueError(w, e) {
			return
		}
		v, replay, e := m.Service.FinalizeFunding(r.Context(), p, db.FinalizeStoredValueFunding{OrganizationID: c.OrganizationID, OrderID: r.PathValue("id"), ExpectedOrderVersion: n, RequestKey: r.Header.Get("Idempotency-Key")})
		if storedValueError(w, e) {
			return
		}
		if replay {
			w.Header().Set("Idempotency-Replayed", "true")
		}
		storedValueFundingJSON(w, 201, v)
	})
	mux.HandleFunc("GET /v1/franchise/stored-value/funding-result", func(w http.ResponseWriter, r *http.Request) {
		p, ok := protected(w, r, "stored_value:read", r.URL.Query().Get("organization_id"))
		if !ok {
			return
		}
		v, e := m.Service.FundingResult(r.Context(), p, r.Header.Get("Idempotency-Key"))
		if storedValueError(w, e) {
			return
		}
		storedValueFundingJSON(w, 200, v)
	})
}
````

### FILE: `internal/platform/postgres/stored_value.go`

```yaml
block_id: "GO-APPROVED-STORED-VALUE-TENDER:file8:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration/configuration/authorization/persistence/transport/UI glue; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "aa2e4b56ecfcf187522bb2a9f4455bf5400939bb07a31e6114bb905482666c3c"
variables: []
secrets_allowed: false
```

````go
package postgres

// AUTHORED transaction/binding glue. Commercial points and reward arithmetic
// are executed by the separately licensed, hash-locked Odoo derivation.
import (
	"context"
	"encoding/json"
	"errors"
	"math/big"
	"sort"
	"time"

	"elite.local/enterprise/internal/approval"
	"elite.local/enterprise/internal/bcamounts"
	"elite.local/enterprise/internal/platform/identity"
	sv "elite.local/enterprise/internal/storedvaluebridge"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StoredValue struct {
	pool      *pgxpool.Pool
	profile   *sv.Profile
	process   sv.Process
	approvals *HumanApprovals
}

func NewStoredValue(pool *pgxpool.Pool, profile *sv.Profile, process sv.Process) (*StoredValue, error) {
	if pool == nil || !profile.Valid() || process.Validate() != nil {
		return nil, sv.ErrBinding
	}
	return &StoredValue{pool, profile, process, NewHumanApprovals(pool)}, nil
}
func (s *StoredValue) allowed(p identity.Principal, permission string) bool {
	if s == nil || s.pool == nil || !s.profile.Valid() {
		return false
	}
	tenant, org := s.profile.Scope()
	return approvalPrincipal(p, tenant, org, permission)
}
func (s *StoredValue) ProfileSHA256() string {
	if s == nil {
		return ""
	}
	return s.profile.SHA256()
}
func points(v string) (*big.Rat, error) {
	n, e := sv.Minor(v, 6)
	if e != nil {
		return nil, e
	}
	return new(big.Rat).SetFrac(big.NewInt(n), big.NewInt(1000000)), nil
}
func decimalPoints(v *big.Rat) string { return v.FloatString(6) }

// ApprovalView checks tenant, organization and the complete approved payload.
func (s *StoredValue) Read(ctx context.Context, p identity.Principal, id string) (sv.ApprovalView, error) {
	var out sv.ApprovalView
	if !s.allowed(p, "stored_value:read") || !sv.ValidID(id) {
		return out, sv.ErrBinding
	}
	tx, e := s.pool.Begin(ctx)
	if e != nil {
		return out, e
	}
	defer tx.Rollback(ctx)
	out, e = s.readTx(ctx, tx, p, id)
	if e != nil {
		return out, e
	}
	return out, tx.Commit(ctx)
}
func (s *StoredValue) readTx(ctx context.Context, tx pgx.Tx, p identity.Principal, id string) (sv.ApprovalView, error) {
	var out sv.ApprovalView
	tenant, org := s.profile.Scope()
	a, state, e := readHumanApproval(ctx, tx, tenant, id, false)
	if e != nil {
		return out, e
	}
	if a.OrganizationID != org || a.Request.Kind != approval.KindStoredValueOperation || sv.Decode(a.Payload, &out.Payload) != nil {
		return out, sv.ErrBinding
	}
	_, hash, e := sv.Canonical(out.Payload)
	if e != nil || hash != a.Request.EvidenceSHA {
		return out, sv.ErrBinding
	}
	out.ApprovalID = id
	out.PayloadSHA256 = hash
	out.State = string(state)
	out.Receipt = json.RawMessage("null")
	var storedSHA string
	e = tx.QueryRow(ctx, `select receipt,receipt_sha256 from stored_value.operation where tenant_id=$1 and approval_id=$2`, tenant, id).Scan(&out.Receipt, &storedSHA)
	if e == nil {
		raw, hash, err := approval.CanonicalPayload(out.Receipt)
		if err != nil || hash != storedSHA {
			return out, sv.ErrBinding
		}
		out.Receipt = raw
	}
	if errors.Is(e, pgx.ErrNoRows) {
		e = nil
	}
	return out, e
}
func (s *StoredValue) Propose(ctx context.Context, p identity.Principal, r sv.Request) (sv.ApprovalView, error) {
	var out sv.ApprovalView
	if !s.allowed(p, "stored_value:request") || r.Validate(s.profile) != nil {
		return out, sv.ErrBinding
	}
	tenant, _ := s.profile.Scope()
	id := sv.StableID("svapproval", r.OperationID)
	tx, e := s.pool.Begin(ctx)
	if e != nil {
		return out, e
	}
	defer tx.Rollback(ctx)
	// Serializes same request before order/account locks; an exact replay does not
	// consult a later order revision or regenerate the immutable lease/payload.
	if _, e = tx.Exec(ctx, `select pg_advisory_xact_lock(hashtextextended($1,0))`, tenant+":stored-value-request:"+r.OperationID); e != nil {
		return out, e
	}
	existing, e := s.readTx(ctx, tx, p, id)
	if e == nil {
		a, _, _ := sv.Canonical(r)
		b, _, _ := sv.Canonical(existing.Payload.Request)
		if string(a) != string(b) || existing.Payload.Requester != p.Subject {
			return out, approval.ErrDuplicate
		}
		existing.Replay = true
		return existing, tx.Commit(ctx)
	}
	if !errors.Is(e, approval.ErrNotFound) {
		return out, e
	}
	var expires time.Time
	if e = tx.QueryRow(ctx, `select date_trunc('microseconds',clock_timestamp())+make_interval(secs=>$1)`, s.profile.ReservationSeconds()).Scan(&expires); e != nil {
		return out, e
	}
	bound, e := s.calculate(ctx, tx, r, p.Subject, expires.UTC().Format(time.RFC3339Nano), false)
	if e != nil {
		return out, e
	}
	raw, hash, e := sv.Canonical(bound)
	if e != nil {
		return out, e
	}
	spec := HumanApprovalSpec{Request: approval.Request{TenantID: tenant, ID: id, Kind: approval.KindStoredValueOperation, SubjectID: r.OrderID, Requester: p.Subject, EvidenceSHA: hash}, OrganizationID: r.OrganizationID, Payload: raw}
	if _, e = s.approvals.submitTx(ctx, tx, p, spec, "stored_value:request", nil); e != nil {
		return out, e
	}
	if r.Operation == "redeem" {
		for _, entry := range bound.Entries {
			delta, e := points(entry.PointsDelta)
			if e != nil {
				return out, e
			}
			delta.Neg(delta)
			_, e = tx.Exec(ctx, `insert into stored_value.reservation(tenant_id,operation_id,account_id,order_id,approval_id,points,state,expires_at) values($1,$2,$3,$4,$5,$6::numeric,'reserved',$7)`, tenant, r.OperationID, entry.AccountID, r.OrderID, id, decimalPoints(delta), expires)
			if e != nil {
				return out, e
			}
		}
	}
	out = sv.ApprovalView{ApprovalID: id, PayloadSHA256: hash, State: string(approval.StatePending), Payload: bound, Receipt: json.RawMessage("null")}
	return out, tx.Commit(ctx)
}

func (s *StoredValue) Decide(ctx context.Context, p identity.Principal, d sv.Decision) (sv.ApprovalView, error) {
	var out sv.ApprovalView
	if !s.allowed(p, "stored_value:approve") || !sv.ValidID(d.ApprovalID) {
		return out, sv.ErrBinding
	}
	tenant, org := s.profile.Scope()
	if d.OrganizationID != org {
		return out, sv.ErrBinding
	}
	_, e := s.approvals.Decide(ctx, p, tenant, d.ApprovalID, org, d.PayloadSHA256, d.Approved, d.Reason, "stored_value:approve", func(ctx context.Context, tx pgx.Tx) error {
		a, _, e := readHumanApproval(ctx, tx, tenant, d.ApprovalID, false)
		if e != nil {
			return e
		}
		var approved sv.BoundOperation
		if a.Request.Kind != approval.KindStoredValueOperation || sv.Decode(a.Payload, &approved) != nil {
			return sv.ErrBinding
		}
		if !d.Approved {
			_, e = tx.Exec(ctx, `update stored_value.reservation set state='rejected' where tenant_id=$1 and approval_id=$2 and state='reserved'`, tenant, d.ApprovalID)
			return e
		}
		actual, e := s.calculate(ctx, tx, approved.Request, approved.Requester, approved.ExpiresAt, true)
		if e != nil {
			return e
		}
		_, actualSHA, e := sv.Canonical(actual)
		if e != nil || actualSHA != d.PayloadSHA256 {
			return sv.ErrBinding
		}
		return s.commit(ctx, tx, d.ApprovalID, d.PayloadSHA256, p.Subject, actual)
	})
	if e != nil {
		return out, e
	}
	// The authorization for this response is the decision permission, not a
	// separate read grant that could hide an already committed local result.
	tx, e := s.pool.Begin(ctx)
	if e != nil {
		return out, e
	}
	defer tx.Rollback(ctx)
	out, e = s.readTx(ctx, tx, p, d.ApprovalID)
	if e != nil {
		return out, e
	}
	return out, tx.Commit(ctx)
}

func (s *StoredValue) order(ctx context.Context, tx pgx.Tx, r sv.Request) (sv.Order, string, sv.Allocation, error) {
	var o sv.Order
	var state string
	var a sv.Allocation
	tenant, _ := s.profile.Scope()
	var subject *string
	e := tx.QueryRow(ctx, `select o.order_id,o.organization_id,o.customer_principal_id,o.state,o.currency,o.version,o.total_minor_units from sales.customer_order o where tenant_id=$1 and order_id=$2 and organization_id=$3 for update`, tenant, r.OrderID, r.OrganizationID).Scan(&o.OrderID, &o.OrganizationID, &subject, &state, &o.Currency, &a.OrderVersion, &a.GrossMinor)
	if e != nil {
		return o, state, a, e
	}
	if a.OrderVersion != r.ExpectedOrderVersion || subject == nil || !sv.ValidID(*subject) {
		return o, state, a, sv.ErrBinding
	}
	o.SubjectID = *subject
	o.State = "draft"
	o.EnabledRuleIDs = []string{}
	o.Lines = []sv.Line{}
	e = tx.QueryRow(ctx, `select gift_minor_units,discount_minor_units,provider_due_minor_units from stored_value.order_allocation where tenant_id=$1 and order_id=$2`, tenant, r.OrderID).Scan(&a.GiftMinor, &a.DiscountMinor, &a.ProviderDueMinor)
	if e != nil {
		return o, state, a, e
	}
	if a.GiftMinor < 0 || a.DiscountMinor < 0 || a.ProviderDueMinor < 0 {
		return o, state, a, sv.ErrBinding
	}
	a.OrderID = o.OrderID
	a.Currency = o.Currency
	rows, e := tx.Query(ctx, `select line_id,variant_id,quantity,unit_price_minor_units from sales.customer_order_line where tenant_id=$1 and order_id=$2 order by line_id`, tenant, r.OrderID)
	if e != nil {
		return o, state, a, e
	}
	total := new(big.Int)
	for rows.Next() {
		var l sv.Line
		var qty, price int64
		if e = rows.Scan(&l.ID, &l.ProductID, &qty, &price); e != nil {
			rows.Close()
			return o, state, a, e
		}
		amount, e := bcamounts.LineAmount(qty, price, 0)
		if e != nil {
			rows.Close()
			return o, state, a, e
		}
		total.Add(total, big.NewInt(amount))
		l.Quantity = sv.Number(qty)
		l.Total = sv.Major(amount, 2)
		l.Subtotal = l.Total
		l.Tax = "0"
		o.Lines = append(o.Lines, l)
	}
	e = rows.Err()
	rows.Close()
	if e != nil {
		return o, state, a, e
	}
	if len(o.Lines) < 1 || len(o.Lines) > 40 || !total.IsInt64() || total.Int64() != a.GrossMinor {
		return o, state, a, sv.ErrBinding
	}
	// Upstream reward lines are represented separately from immutable gross
	// sales lines. These typed negative lines preserve source points filtering.
	rows, e = tx.Query(ctx, `select op.operation_id,op.program_id,ac.kind,sum(en.applied_minor_units)::bigint from stored_value.operation op join stored_value.entry en using(tenant_id,operation_id) join stored_value.account ac using(tenant_id,account_id) where op.tenant_id=$1 and op.order_id=$2 group by op.operation_id,op.program_id,ac.kind order by op.operation_id`, tenant, r.OrderID)
	if e != nil {
		return o, state, a, e
	}
	for rows.Next() {
		var id, program, kind string
		var amount int64
		if e = rows.Scan(&id, &program, &kind, &amount); e != nil {
			rows.Close()
			return o, state, a, e
		}
		if amount == 0 {
			continue
		}
		val := sv.Major(-amount, 2)
		o.Lines = append(o.Lines, sv.Line{ID: id, ProductID: "stored-value-reward", Quantity: "1", Subtotal: val, Tax: "0", Total: val, RewardProgramType: kind, RewardProgramID: program, RewardTrigger: "auto"})
	}
	e = rows.Err()
	rows.Close()
	o.Total = sv.Major(a.ProviderDueMinor, 2)
	return o, state, a, e
}

func (s *StoredValue) account(ctx context.Context, tx pgx.Tx, r sv.Request, subject, id string) (string, error) {
	tenant, _ := s.profile.Scope()
	if _, e := tx.Exec(ctx, `select pg_advisory_xact_lock(hashtextextended($1,0))`, tenant+":stored-value-account:"+id); e != nil {
		return "", e
	}
	var kind, currency, program, org string
	var owner *string
	e := tx.QueryRow(ctx, `select kind,currency,program_id,organization_id,customer_subject from stored_value.account where tenant_id=$1 and account_id=$2 for update`, tenant, id).Scan(&kind, &currency, &program, &org, &owner)
	if e != nil {
		return "", e
	}
	pp, e := s.profile.Program(r.ProgramID)
	if e != nil || kind != pp.Calculation.Kind || currency != pp.Calculation.Currency || program != r.ProgramID || org != r.OrganizationID || (owner != nil && *owner != subject) {
		return "", sv.ErrBinding
	}
	var balance string
	e = tx.QueryRow(ctx, `select (coalesce((select sum(points_delta) from stored_value.entry where tenant_id=$1 and account_id=$2),0)-coalesce((select sum(points) from stored_value.reservation where tenant_id=$1 and account_id=$2 and operation_id<>$3 and state='reserved' and expires_at>clock_timestamp()),0))::numeric(38,6)::text`, tenant, id, r.OperationID).Scan(&balance)
	return balance, e
}

func (s *StoredValue) calculate(ctx context.Context, tx pgx.Tx, r sv.Request, requester, expiry string, deciding bool) (sv.BoundOperation, error) {
	out := sv.BoundOperation{Schema: "elite.stored-value-operation.v1", Request: r, Requester: requester, EngineSHA256: s.process.ManifestSHA256, ExpiresAt: expiry, Entries: []sv.Entry{}}
	if r.Validate(s.profile) != nil {
		return out, sv.ErrBinding
	}
	tenant, _ := s.profile.Scope()
	pp, e := s.profile.Program(r.ProgramID)
	if e != nil {
		return out, e
	}
	o, state, a, e := s.order(ctx, tx, r)
	if e != nil {
		return out, e
	}
	out.Order = o
	out.OrderState = state
	out.GrossMinor = a.GrossMinor
	out.GiftMinor = a.GiftMinor
	out.DiscountMinor = a.DiscountMinor
	if o.Currency != pp.Calculation.Currency {
		return out, sv.ErrBinding
	}
	var validTime bool
	if e = tx.QueryRow(ctx, `select $1::timestamptz>clock_timestamp()`, expiry).Scan(&validTime); e != nil {
		return out, e
	}
	if !validTime {
		return out, sv.ErrBinding
	}
	if r.Operation != "reverse" && !s.profile.StateAllowed(r.Operation, state) {
		return out, sv.ErrBinding
	}
	if r.Operation == "redeem" || r.Operation == "reverse" {
		var frozen bool
		// Unknown provider effects remain fenced. Source cancellation may reopen
		// stored value only after the authoritative payment owner reaches refunded.
		e = tx.QueryRow(ctx, `select exists(select 1 from payment.payment_attempt where tenant_id=$1 and order_id=$2 and state not in ('failed','refunded')) or ($3::text<>'cancelled' and exists(select 1 from payment.local_funding_receipt where tenant_id=$1 and order_id=$2))`, tenant, r.OrderID, state).Scan(&frozen)
		if e != nil {
			return out, e
		}
		if frozen {
			return out, sv.ErrBinding
		}
	}
	data := map[string]any{}
	calcOp := "evaluate"
	switch r.Operation {
	case "issue", "accrue":
		if (r.Operation == "issue") != (pp.Calculation.Kind == "gift_card") {
			return out, sv.ErrBinding
		}
		var repeated bool
		e = tx.QueryRow(ctx, `select exists(select 1 from stored_value.operation op where tenant_id=$1 and order_id=$2 and program_id=$3 and operation in ('issue','accrue') and not exists(select 1 from stored_value.operation rev where rev.tenant_id=op.tenant_id and rev.original_operation_id=op.operation_id))`, tenant, r.OrderID, r.ProgramID).Scan(&repeated)
		if e != nil {
			return out, e
		}
		if repeated {
			return out, sv.ErrBinding
		}
	case "redeem":
		if pp.Calculation.Kind == "loyalty" && (a.DiscountMinor != 0 || a.GiftMinor != 0) {
			return out, sv.ErrBinding
		}
		balance, e := s.account(ctx, tx, r, o.SubjectID, r.AccountID)
		if e != nil {
			return out, e
		}
		var futureFromThis bool
		e = tx.QueryRow(ctx, `select exists(select 1 from stored_value.entry en join stored_value.operation op using(tenant_id,operation_id) where en.tenant_id=$1 and en.account_id=$2 and op.order_id=$3 and op.operation='issue')`, tenant, r.AccountID, r.OrderID).Scan(&futureFromThis)
		if e != nil {
			return out, e
		}
		if futureFromThis {
			return out, sv.ErrBinding
		}
		if deciding {
			var reserved bool
			e = tx.QueryRow(ctx, `select exists(select 1 from stored_value.reservation where tenant_id=$1 and operation_id=$2 and account_id=$3 and state='reserved' and expires_at=$4::timestamptz and expires_at>clock_timestamp())`, tenant, r.OperationID, r.AccountID, expiry).Scan(&reserved)
			if e != nil {
				return out, e
			}
			if !reserved {
				return out, sv.ErrBinding
			}
		}
		calcOp = "reward"
		data = map[string]any{"coupon_id": r.AccountID, "balance": balance, "pending_earned": "0", "pending_cost": "0", "discountable": o.Total, "reward": pp.Reward}
	case "reverse":
		if state != "cancelled" {
			return out, sv.ErrBinding
		}
		var originalOrder, originalProgram, operation string
		e = tx.QueryRow(ctx, `select order_id,program_id,operation from stored_value.operation where tenant_id=$1 and operation_id=$2`, tenant, r.OriginalOperationID).Scan(&originalOrder, &originalProgram, &operation)
		if e != nil {
			return out, e
		}
		if originalOrder != r.OrderID || originalProgram != r.ProgramID || operation == "reverse" {
			return out, sv.ErrBinding
		}
		rows, e := tx.Query(ctx, `select account_id,points_delta::text,applied_minor_units from stored_value.entry where tenant_id=$1 and operation_id=$2 order by account_id`, tenant, r.OriginalOperationID)
		if e != nil {
			return out, e
		}
		originals := []sv.Entry{}
		for rows.Next() {
			var entry sv.Entry
			if e = rows.Scan(&entry.AccountID, &entry.PointsDelta, &entry.AppliedMinor); e != nil {
				rows.Close()
				return out, e
			}
			originals = append(originals, entry)
		}
		e = rows.Err()
		rows.Close()
		if e != nil {
			return out, e
		}
		for _, entry := range originals {
			if _, e = s.account(ctx, tx, r, o.SubjectID, entry.AccountID); e != nil {
				return out, e
			}
		}
		calcOp = "reverse"
		data = map[string]any{"entries": originals}
	}
	request, e := sv.NewCalculation(calcOp, pp.Calculation, o, data)
	if e != nil {
		return out, e
	}
	result, e := s.process.Execute(ctx, request)
	if e != nil {
		return out, e
	}
	out.Calculation = request
	out.Result = result
	switch r.Operation {
	case "issue", "accrue":
		var resultData struct {
			Points []string `json:"points"`
			Error  string   `json:"error"`
		}
		if json.Unmarshal(result.Result, &resultData) != nil || resultData.Error != "" {
			return out, sv.ErrBinding
		}
		for i, value := range resultData.Points {
			n, e := points(value)
			if e != nil || n.Sign() < 0 {
				return out, sv.ErrBinding
			}
			if n.Sign() == 0 {
				continue
			}
			id := sv.StableID("gift", tenant, r.OperationID, sv.Number(int64(i)))
			if r.Operation == "accrue" {
				id = sv.StableID("loyalty", tenant, r.OrganizationID, r.ProgramID, o.SubjectID)
			}
			out.Entries = append(out.Entries, sv.Entry{AccountID: id, PointsDelta: decimalPoints(n)})
		}
		if len(out.Entries) == 0 {
			return out, sv.ErrBinding
		}
	case "redeem":
		var v struct {
			Value      string `json:"value"`
			PointsCost string `json:"points_cost"`
		}
		if sv.Decode(result.Result, &v) != nil {
			return out, sv.ErrBinding
		}
		amount, e := sv.Minor(v.Value, 2)
		if e != nil || amount <= 0 || amount > a.ProviderDueMinor {
			return out, sv.ErrBinding
		}
		n, e := points(v.PointsCost)
		if e != nil || n.Sign() <= 0 {
			return out, sv.ErrBinding
		}
		n.Neg(n)
		out.Entries = []sv.Entry{{AccountID: r.AccountID, PointsDelta: decimalPoints(n), AppliedMinor: amount}}
	case "reverse":
		if json.Unmarshal(result.Result, &out.Entries) != nil || len(out.Entries) == 0 {
			return out, sv.ErrBinding
		}
	}
	sort.Slice(out.Entries, func(i, j int) bool { return out.Entries[i].AccountID < out.Entries[j].AccountID })
	return out, nil
}

func (s *StoredValue) commit(ctx context.Context, tx pgx.Tx, approvalID, approvedHash, reviewer string, b sv.BoundOperation) error {
	tenant, _ := s.profile.Scope()
	r := b.Request
	pp, e := s.profile.Program(r.ProgramID)
	if e != nil {
		return e
	}
	receipt := map[string]any{"schema": "elite.stored-value-receipt.v1", "approval_id": approvalID, "approved_sha256": approvedHash, "reviewer": reviewer, "operation": b}
	raw, receiptSHA, e := sv.Canonical(receipt)
	if e != nil {
		return e
	}
	_, requestSHA, e := sv.Canonical(r)
	if e != nil {
		return e
	}
	_, calcSHA, e := sv.Canonical(b.Result)
	if e != nil {
		return e
	}
	for _, entry := range b.Entries {
		if r.Operation == "issue" || r.Operation == "accrue" {
			if _, e = tx.Exec(ctx, `select pg_advisory_xact_lock(hashtextextended($1,0))`, tenant+":stored-value-account:"+entry.AccountID); e != nil {
				return e
			}
			var owner any
			if pp.Calculation.Nominative {
				owner = b.Order.SubjectID
			}
			_, e = tx.Exec(ctx, `insert into stored_value.account(tenant_id,account_id,organization_id,program_id,kind,customer_subject,currency) values($1,$2,$3,$4,$5,$6,$7) on conflict do nothing`, tenant, entry.AccountID, r.OrganizationID, r.ProgramID, pp.Calculation.Kind, owner, pp.Calculation.Currency)
			if e != nil {
				return e
			}
			if _, e = s.account(ctx, tx, r, b.Order.SubjectID, entry.AccountID); e != nil {
				return e
			}
		}
	}
	_, e = tx.Exec(ctx, `insert into stored_value.operation(tenant_id,operation_id,organization_id,order_id,order_version,program_id,operation,original_operation_id,approval_id,profile_sha256,request_sha256,calculation_sha256,receipt_sha256,receipt) values($1,$2,$3,$4,$5,$6,$7,nullif($8,''),$9,$10,$11,$12,$13,$14)`, tenant, r.OperationID, r.OrganizationID, r.OrderID, r.ExpectedOrderVersion, r.ProgramID, r.Operation, r.OriginalOperationID, approvalID, r.ProfileSHA256, requestSHA, calcSHA, receiptSHA, raw)
	if e != nil {
		return e
	}
	allocated := false
	for _, entry := range b.Entries {
		_, e = tx.Exec(ctx, `insert into stored_value.entry(tenant_id,operation_id,account_id,points_delta,applied_minor_units) values($1,$2,$3,$4::numeric,$5)`, tenant, r.OperationID, entry.AccountID, entry.PointsDelta, entry.AppliedMinor)
		if e != nil {
			return e
		}
		allocated = allocated || entry.AppliedMinor != 0
	}
	if _, e = tx.Exec(ctx, `update stored_value.reservation set state='consumed' where tenant_id=$1 and operation_id=$2 and state='reserved'`, tenant, r.OperationID); e != nil {
		return e
	}
	if allocated {
		if _, e = tx.Exec(ctx, `update sales.customer_order set version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and order_id=$2`, tenant, r.OrderID); e != nil {
			return e
		}
	}
	_, e = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload) values($1,gen_random_uuid(),'stored-value',$2,1,'stored-value.committed',1,clock_timestamp(),$3)`, tenant, r.OperationID, raw)
	return e
}

func (s *StoredValue) Allocation(ctx context.Context, p identity.Principal, order string) (sv.Allocation, error) {
	var out sv.Allocation
	if !s.allowed(p, "stored_value:read") || !sv.ValidID(order) {
		return out, sv.ErrBinding
	}
	tenant, org := s.profile.Scope()
	e := s.pool.QueryRow(ctx, `select order_id,order_version,currency,gross_minor_units,gift_minor_units,discount_minor_units,provider_due_minor_units from stored_value.order_allocation where tenant_id=$1 and organization_id=$2 and order_id=$3`, tenant, org, order).Scan(&out.OrderID, &out.OrderVersion, &out.Currency, &out.GrossMinor, &out.GiftMinor, &out.DiscountMinor, &out.ProviderDueMinor)
	return out, e
}
````

### FILE: `internal/platform/postgres/stored_value_browser_integration_test.go`

```yaml
block_id: "GO-APPROVED-STORED-VALUE-TENDER:file9:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration/configuration/authorization/persistence/transport/UI glue; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "627e7ad4c0b2ae4bb0a92cd7f8badbb8ba0801e3f95cb976b4f53f2c0486aa5a"
variables: []
secrets_allowed: false
```

````go
package postgres_test

// AUTHORED connected fixture: distinct reviewers, exact source calculation,
// JWE role sessions, real OIDC/JWKS, Next BFF and durable zero-provider delivery.
import (
	"context"
	"crypto/sha256"
	"elite.local/enterprise/internal/commerce"
	"elite.local/enterprise/internal/franchisejourney"
	"elite.local/enterprise/internal/platform/httpapi"
	db "elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/platform/randomid"
	"encoding/hex"
	"encoding/json"
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

func TestStoredValueBrowserPostgres(t *testing.T) {
	if os.Getenv("ELITE_STORED_VALUE_BROWSER") != "1" {
		t.Skip("explicit stored value browser fixture required")
	}
	pool := connectedPool(t)
	tenant := connectedSeedOrder(t, pool, "stripe")
	store, _, _ := fixtureStoredValueSetup(t, pool, tenant, "gift_card", 200000)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	mode := false
	d := franchisejourney.HandoverProfileDocument{Schema: franchisejourney.HandoverProfileSchema, ProfileID: "browser-franchise", Revision: 1, Algorithm: franchisejourney.HandoverSupportedAlgorithm, AlgorithmRevision: 3, Scope: "MATERIALIZED_PROFILE", TenantID: tenant, OrganizationID: "store", ProviderCode: "stripe", ProviderAccountRef: "acct_fixture", ProviderConnectionID: "checkout", ExpectedLiveMode: &mode, MaximumObservationAgeSeconds: 900, Options: franchisejourney.SupportedStoredValueReleaseOptions(), AuthorityReference: "docs/handover-operator-flow.md", DecisionReference: "LOCAL_BROWSER_FIXTURE"}
	raw, err := json.Marshal(d)
	if err != nil {
		t.Fatal(err)
	}
	digest := sha256.Sum256(raw)
	policy, err := franchisejourney.LoadHandoverProfile(raw, franchisejourney.HandoverProfileActivation{Enabled: true, ProfileID: d.ProfileID, ProfileRevision: 1, DocumentSHA256: hex.EncodeToString(digest[:]), TenantID: tenant, OrganizationID: "store", ProviderCode: "stripe", ProviderAccountRef: "acct_fixture", ProviderConnectionID: "checkout"})
	if err != nil {
		t.Fatal(err)
	}
	repo := db.NewFranchiseJourney(pool)
	ids := randomid.Generator{}
	service, err := franchisejourney.NewHandoverPreparationService(repo, ids, policy)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = repo.PublishDeliveryChecklist(ctx, tenant, "fixture-operator", franchisejourney.DeliveryChecklist{ID: "browser-checklist", OrganizationID: "store", Version: 1, Title: "Preparación de referencia", Items: []franchisejourney.ChecklistItem{{ID: "serial", Ordinal: 1, Prompt: "Verificar serie", ResponseType: "serial", Required: true}}}, ids.New()); err != nil {
		t.Fatal(err)
	}
	verifier, token := handoverBrowserIssuer(t)
	identities := map[string]any{}
	for _, name := range []string{"operator", "reviewer", "customer", "stranger", "foreign-org"} {
		permissions := []string{"customer:self"}
		organizations := []string{"store"}
		if name == "operator" || name == "reviewer" || name == "foreign-org" {
			permissions = []string{"handover:manage", "stored_value:read", "stored_value:request", "stored_value:approve", "stored_value:fund"}
		}
		if name == "foreign-org" {
			organizations = []string{"other"}
		}
		identities[name] = map[string]any{"subject": name, "tenantId": tenant, "permissions": permissions, "organizations": organizations, "accessToken": token(name, tenant, permissions, organizations)}
	}
	mux := http.NewServeMux()
	httpapi.FranchiseJourneyModule{Service: franchisejourney.NewService(repo, ids, handoverBrowserClock{})}.Register(mux, verifier)
	httpapi.CommerceModule{Service: commerce.NewService(db.NewCommerce(pool), ids), PaymentProvider: "stripe", PaymentTenantID: tenant, PaymentOrganizationID: "store", ProviderObservedPayments: true}.Register(mux, verifier)
	httpapi.InitialHandoverModule{Service: service}.Register(mux, verifier)
	httpapi.HandoverContextModule{Service: service}.Register(mux, verifier)
	httpapi.CommercialReleaseModule{Service: service}.Register(mux, verifier)
	httpapi.StoredValueModule{Service: store}.Register(mux, verifier)
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
	env = append(env, "ELITE_STORED_VALUE_BROWSER=1", "ELITE_WORKSPACE_E2E=enabled", "ELITE_HANDOVER_IDENTITIES="+string(encoded), "ELITE_WEB_ROOT="+web, "ELITE_BASE_URL="+edge.URL, "APP_BASE_URL="+edge.URL, "ENTERPRISE_API_BASE_URL="+api.URL, "AUTH_SESSION_SECRET=synthetic-"+ids.New()+ids.New())
	artifacts, err := os.MkdirTemp(web, "stored-value-browser-artifacts-")
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
	command := exec.CommandContext(ctx, node, filepath.Join(gate, "node_modules", "@playwright", "test", "cli.js"), "test", "tests/stored-value-connected.spec.mjs", "--project=chromium-desktop", "--timeout=90000", "--output="+filepath.Join(artifacts, "browsers"))
	command.Dir, command.Env = gate, env
	output, err := command.CombinedOutput()
	if e := os.WriteFile(filepath.Join(artifacts, "browser.log"), output, 0600); e != nil {
		t.Fatal(e)
	}
	if err != nil || !strings.Contains(string(output), "1 passed") {
		t.Fatalf("browser %v\n%s\nartifacts=%s", err, output, artifacts)
	}

	var prepared, releases, accepted, operations, funding, payments int
	err = pool.QueryRow(ctx, `select (select count(*) from sales.delivery_handover_preparation where tenant_id=$1),(select count(*) from sales.commercial_release_receipt where tenant_id=$1),(select count(*) from sales.delivery_handover where tenant_id=$1 and state='accepted'),(select count(*) from stored_value.operation where tenant_id=$1),(select count(*) from payment.local_funding_receipt where tenant_id=$1),(select count(*) from payment.payment_attempt where tenant_id=$1)`, tenant).Scan(&prepared, &releases, &accepted, &operations, &funding, &payments)
	if err != nil || prepared != 1 || releases != 1 || accepted != 1 || operations != 2 || funding != 1 || payments != 0 {
		t.Fatal("durable outcomes", prepared, releases, accepted, operations, funding, payments, err)
	}
	mu.Lock()
	countsJSON, _ := json.Marshal(counts)
	mu.Unlock()
	if err = os.WriteFile(filepath.Join(artifacts, "api-post-counts.json"), countsJSON, 0600); err != nil {
		t.Fatal(err)
	}
	t.Logf("STORED_VALUE_BROWSER_POSTGRES_PASS actual_BFF_SQL=true JWE_RS256_JWKS=true source_operations=2 funding=1 payments=0 handover=1 commercial_receipt=1 artifacts=%s", artifacts)
}
````

### FILE: `internal/platform/postgres/stored_value_funding.go`

```yaml
block_id: "GO-APPROVED-STORED-VALUE-TENDER:file10:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration/configuration/authorization/persistence/transport/UI glue; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "bdffe3f7b971cfeb9b5239c6e44833a90b80c1b7c705825cfe2cfc101bd1966d"
variables: []
secrets_allowed: false
```

````go
package postgres

// AUTHORED final observation and request recovery; all applied contributions
// already exist under immutable source-derived, separately approved operations.
import (
	"context"
	"elite.local/enterprise/internal/platform/identity"
	sv "elite.local/enterprise/internal/storedvaluebridge"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5"
	"time"
)

type FinalizeStoredValueFunding struct {
	OrganizationID       string `json:"organization_id"`
	OrderID              string `json:"order_id"`
	ExpectedOrderVersion int64  `json:"expected_order_version"`
	RequestKey           string `json:"request_key"`
}

func storedFundingRequestHash(p identity.Principal, c FinalizeStoredValueFunding) string {
	raw, _ := json.Marshal(struct {
		Tenant, Actor string
		Command       FinalizeStoredValueFunding
	}{p.TenantID, p.Subject, c})
	return fundingHash(raw)
}
func readStoredFundingRequest(ctx context.Context, q orderFundingReader, tenant, org, actor, key, hash string) (LocalFundingResult, error) {
	var out LocalFundingResult
	var raw []byte
	var saved string
	e := q.QueryRow(ctx, `select receipt_raw,receipt_sha256,request_sha256 from payment.local_funding_receipt where tenant_id=$1 and organization_id=$2 and requested_by_subject=$3 and request_key=$4`, tenant, org, actor, key).Scan(&raw, &out.SHA256, &saved)
	if e != nil {
		return out, e
	}
	if hash != "" && saved != hash {
		return out, sv.ErrBinding
	}
	if fundingHash(raw) != out.SHA256 || json.Unmarshal(raw, &out.Receipt) != nil || out.Receipt.RequestedBy != actor || out.Receipt.RequestKey != key || out.Receipt.RequestSHA256 != saved || out.Receipt.Effect != LocalFundingEffect || out.Receipt.Allocation.TenantID != tenant || out.Receipt.Allocation.OrganizationID != org || out.Receipt.AllocationSHA256 != out.Receipt.Allocation.SHA256() {
		return out, sv.ErrBinding
	}
	return out, nil
}
func (s *StoredValue) FundingResult(ctx context.Context, p identity.Principal, key string) (LocalFundingResult, error) {
	if !s.allowed(p, "stored_value:read") || len(key) < 16 || !sv.ValidID(key) {
		return LocalFundingResult{}, sv.ErrBinding
	}
	tenant, org := s.profile.Scope()
	return readStoredFundingRequest(ctx, s.pool, tenant, org, p.Subject, key, "")
}
func (s *StoredValue) FinalizeFunding(ctx context.Context, p identity.Principal, c FinalizeStoredValueFunding) (LocalFundingResult, bool, error) {
	var empty LocalFundingResult
	if !s.allowed(p, "stored_value:fund") || !sv.ValidID(c.OrderID) || !sv.ValidID(c.RequestKey) || len(c.RequestKey) < 16 || c.ExpectedOrderVersion < 1 {
		return empty, false, sv.ErrBinding
	}
	tenant, org := s.profile.Scope()
	if org != c.OrganizationID {
		return empty, false, sv.ErrBinding
	}
	hash := storedFundingRequestHash(p, c)
	tx, e := s.pool.Begin(ctx)
	if e != nil {
		return empty, false, e
	}
	defer tx.Rollback(ctx)
	if old, e := readStoredFundingRequest(ctx, tx, tenant, org, p.Subject, c.RequestKey, hash); e == nil {
		return old, true, tx.Commit(ctx)
	} else if !errors.Is(e, pgx.ErrNoRows) {
		return empty, false, e
	}
	var currency, state string
	var gross, version int64
	e = tx.QueryRow(ctx, `select currency,total_minor_units,state,version from sales.customer_order where tenant_id=$1 and organization_id=$2 and order_id=$3 for update`, tenant, org, c.OrderID).Scan(&currency, &gross, &state, &version)
	if e != nil {
		return empty, false, e
	}
	if version != c.ExpectedOrderVersion || (state != "placed" && state != "confirmed" && state != "allocated") {
		return empty, false, sv.ErrBinding
	}
	a, e := readOrderFunding(ctx, tx, tenant, org, c.OrderID, currency, gross)
	if e != nil {
		return empty, false, e
	}
	if a.ProviderMinor != 0 || a.GiftMinor+a.DiscountMinor <= 0 {
		return empty, false, sv.ErrBinding
	}
	var active bool
	e = tx.QueryRow(ctx, `select exists(select 1 from payment.payment_attempt where tenant_id=$1 and order_id=$2 and state not in ('failed','refunded'))`, tenant, c.OrderID).Scan(&active)
	if e != nil {
		return empty, false, e
	}
	if active {
		return empty, false, sv.ErrBinding
	}
	var now time.Time
	if e = tx.QueryRow(ctx, `select clock_timestamp()`).Scan(&now); e != nil {
		return empty, false, e
	}
	rec := LocalFundingReceipt{ID: sv.StableID("funding", tenant, c.RequestKey), RequestKey: c.RequestKey, RequestSHA256: hash, RequestedBy: p.Subject, Allocation: a, AllocationSHA256: a.SHA256(), ObservedAt: now.UTC(), Effect: LocalFundingEffect}
	raw, e := json.Marshal(rec)
	if e != nil {
		return empty, false, e
	}
	receiptHash := fundingHash(raw)
	inserted, e := tx.Exec(ctx, `insert into payment.local_funding_receipt(tenant_id,funding_id,request_key,request_sha256,requested_by_subject,receipt_raw,order_id,organization_id,currency,gross_minor_units,gift_minor_units,discount_minor_units,provider_minor_units,allocation_sha256,receipt_sha256,receipt,created_at)values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,0,$13,$14,$15,$16) on conflict(tenant_id,request_key)do nothing`, tenant, rec.ID, c.RequestKey, hash, p.Subject, raw, c.OrderID, org, currency, gross, a.GiftMinor, a.DiscountMinor, a.SHA256(), receiptHash, raw, now)
	if e != nil {
		return empty, false, e
	}
	if inserted.RowsAffected() == 0 {
		old, e := readStoredFundingRequest(ctx, tx, tenant, org, p.Subject, c.RequestKey, hash)
		if e != nil {
			return empty, false, e
		}
		return old, true, tx.Commit(ctx)
	}
	if _, e = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)values($1,gen_random_uuid(),'order-funding',$2,1,'order.local-funding-observed',1,$3,$4)`, tenant, rec.ID, now, raw); e != nil {
		return empty, false, e
	}
	if e = tx.Commit(ctx); e != nil {
		return empty, false, e
	}
	return LocalFundingResult{Receipt: rec, SHA256: receiptHash}, false, nil
}
````

### FILE: `internal/platform/postgres/stored_value_http_integration_test.go`

```yaml
block_id: "GO-APPROVED-STORED-VALUE-TENDER:file11:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration/configuration/authorization/persistence/transport/UI glue; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "e8d42e88b12c039b2bdc768165a7180bb93224753a5435e3a0580a5dab4b6341"
variables: []
secrets_allowed: false
```

````go
package postgres_test

// AUTHORED HTTP composition fixture. Only bearer verification is a principal
// fixture; real HTTP decoding, permission checks, PG owners and source IPC run.
import (
	"context"
	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/platform/identity"
	db "elite.local/enterprise/internal/platform/postgres"
	sv "elite.local/enterprise/internal/storedvaluebridge"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type storedValueFixtureVerifier struct{ p identity.Principal }

func (v storedValueFixtureVerifier) Verify(ctx context.Context, raw string) (identity.Principal, error) {
	if raw != "local-bearer-fixture" {
		return identity.Principal{}, fmt.Errorf("invalid fixture bearer")
	}
	return v.p, nil
}
func storedValueHTTPCall(t *testing.T, store *db.StoredValue, p identity.Principal, method, path, key, body string, want int) []byte {
	t.Helper()
	mux := http.NewServeMux()
	httpapi.StoredValueModule{Service: store}.Register(mux, storedValueFixtureVerifier{p})
	server := httptest.NewServer(mux)
	defer server.Close()
	request, e := http.NewRequest(method, server.URL+path, strings.NewReader(body))
	if e != nil {
		t.Fatal(e)
	}
	request.Header.Set("Authorization", "Bearer local-bearer-fixture")
	request.Header.Set("Content-Type", "application/json")
	if key != "" {
		request.Header.Set("Idempotency-Key", key)
	}
	client := http.Client{Timeout: 8 * time.Second}
	response, e := client.Do(request)
	if e != nil {
		t.Fatal(e)
	}
	defer response.Body.Close()
	raw, e := io.ReadAll(io.LimitReader(response.Body, 262145))
	if e != nil || len(raw) > 262144 {
		t.Fatal(e)
	}
	if response.StatusCode != want {
		t.Fatalf("%s %s status%d want%d: %s", method, path, response.StatusCode, want, raw)
	}
	return raw
}
func storedValueHTTPView(t *testing.T, raw []byte) sv.ApprovalView {
	t.Helper()
	var v struct {
		ApprovalID string `json:"approval_id"`
		SHA        string `json:"payload_sha256"`
		State      string `json:"state"`
		Replay     bool   `json:"replay"`
		Payload    string `json:"payload_json"`
		Receipt    string `json:"receipt_json"`
	}
	if json.Unmarshal(raw, &v) != nil {
		t.Fatal("wire view")
	}
	var payload sv.BoundOperation
	if sv.Decode([]byte(v.Payload), &payload) != nil {
		t.Fatal("exact payload")
	}
	_, hash, e := sv.Canonical(payload)
	if e != nil || hash != v.SHA {
		t.Fatal("wire hash")
	}
	return sv.ApprovalView{ApprovalID: v.ApprovalID, PayloadSHA256: v.SHA, State: v.State, Replay: v.Replay, Payload: payload, Receipt: json.RawMessage(v.Receipt)}
}
func storedValueHTTPPropose(t *testing.T, store *db.StoredValue, p identity.Principal, r sv.Request) sv.ApprovalView {
	t.Helper()
	raw, e := json.Marshal(map[string]string{"operation_id": r.OperationID, "operation": r.Operation, "organization_id": r.OrganizationID, "program_id": r.ProgramID, "order_id": r.OrderID, "expected_order_version": sv.Number(r.ExpectedOrderVersion), "account_id": r.AccountID, "original_operation_id": r.OriginalOperationID, "profile_sha256": r.ProfileSHA256})
	if e != nil {
		t.Fatal(e)
	}
	return storedValueHTTPView(t, storedValueHTTPCall(t, store, p, "POST", "/v1/franchise/stored-value/proposals", "", string(raw), 201))
}
func storedValueHTTPApprove(t *testing.T, store *db.StoredValue, p identity.Principal, v sv.ApprovalView) sv.ApprovalView {
	t.Helper()
	raw, _ := json.Marshal(map[string]any{"organization_id": v.Payload.Request.OrganizationID, "payload_sha256": v.PayloadSHA256, "approved": true, "reason": "Exact source result reviewed through HTTP fixture"})
	out := storedValueHTTPView(t, storedValueHTTPCall(t, store, p, "POST", "/v1/franchise/stored-value/approvals/"+v.ApprovalID+"/decision", "", string(raw), 200))
	recovered := storedValueHTTPView(t, storedValueHTTPCall(t, store, p, "GET", "/v1/franchise/stored-value/approvals/"+v.ApprovalID+"?organization_id="+v.Payload.Request.OrganizationID, "", "", 200))
	if string(recovered.Receipt) != string(out.Receipt) || recovered.PayloadSHA256 != out.PayloadSHA256 {
		t.Fatal("HTTP decision recovery")
	}
	return out
}
func TestStoredValueHTTPRejectsScopeAndUnsafeNumbers(t *testing.T) {
	ctx := context.Background()
	pool := connectedPool(t)
	tenant := connectedSeedOrder(t, pool, "stripe")
	store, p, _ := fixtureStoredValue(t, pool, tenant, "gift_card")
	p.Organizations = map[string]struct{}{"store": {}, "other": {}}
	storedValueHTTPCall(t, store, p, "GET", "/v1/franchise/stored-value/orders/order?organization_id=other", "", "", 403)
	body := `{"operation_id":"invalid-number","operation":"issue","organization_id":"store","program_id":"reference-gift_card","order_id":"source","expected_order_version":9007199254740993,"account_id":"","original_operation_id":"","profile_sha256":"` + store.ProfileSHA256() + `"}`
	storedValueHTTPCall(t, store, p, "POST", "/v1/franchise/stored-value/proposals", "", body, 400)
	if _, e := pool.Exec(ctx, `insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version)values($1,'exact-large','store','customer','placed','ARS',9007199254740993,9007199254740993)`, tenant); e != nil {
		t.Fatal(e)
	}
	raw := storedValueHTTPCall(t, store, p, "GET", "/v1/franchise/stored-value/orders/exact-large?organization_id=store", "", "", 200)
	var out map[string]any
	if json.Unmarshal(raw, &out) != nil || out["gross_minor_units"] != "9007199254740993" || out["order_version"] != "9007199254740993" {
		t.Fatal("large amount/version lost", string(raw))
	}
	denied := p
	denied.Permissions = map[string]struct{}{"stored_value:read": {}}
	storedValueHTTPCall(t, store, denied, "POST", "/v1/franchise/stored-value/orders/order/funding", "permission-negative-key", `{"organization_id":"store","expected_order_version":"2"}`, 403)
	t.Log("actual HTTP/PG: large int64 strings exact; wrong configured scope and numeric JSON version rejected")
}
````

### FILE: `internal/platform/postgres/stored_value_integration_test.go`

```yaml
block_id: "GO-APPROVED-STORED-VALUE-TENDER:file12:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration/configuration/authorization/persistence/transport/UI glue; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "e2bd4b4eff2776f6f3b5c6d98cbb1058407a372ca3d7f3a1c42e8efd3992a278"
variables: []
secrets_allowed: false
```

````go
package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"elite.local/enterprise/internal/approval"
	"elite.local/enterprise/internal/platform/identity"
	sv "elite.local/enterprise/internal/storedvaluebridge"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestStoredValueConnected(t *testing.T) {
	raw := os.Getenv("STORED_VALUE_DB_URL")
	if raw == "" {
		t.Skip("STORED_VALUE_DB_URL not set")
	}
	u, e := url.Parse(raw)
	if e != nil || u.Scheme != "postgres" || u.Hostname() != "127.0.0.1" || !strings.HasPrefix(u.Path, "/elite_stored_value_") {
		t.Fatal("dedicated loopback fixture required")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 70*time.Second)
	defer cancel()
	pool, e := pgxpool.New(ctx, raw)
	if e != nil {
		t.Fatal(e)
	}
	defer pool.Close()
	tenant := "11111111-1111-4111-8111-111111111111"
	org := "franchise-1"
	must := func(q string, args ...any) {
		t.Helper()
		if _, e := pool.Exec(ctx, q, args...); e != nil {
			t.Fatal(e)
		}
	}
	root, e := filepath.Abs("../../..")
	if e != nil {
		t.Fatal(e)
	}
	read := func(p string) []byte {
		t.Helper()
		b, e := os.ReadFile(filepath.Join(root, p))
		if e != nil {
			t.Fatal(e)
		}
		return b
	}
	profile, e := sv.Load(read("deploy/stored-value/profile.reference.json"), sv.Hash(read("deploy/stored-value/profile.reference.json")))
	if e != nil {
		t.Fatal(e)
	}
	process := sv.Process{Python: `C:\Python314\python.exe`, Script: filepath.Join(root, "odoo_loyalty/run.py"), ScriptSHA256: sv.Hash(read("odoo_loyalty/run.py")), Manifest: filepath.Join(root, "odoo_loyalty/engine-lock.json"), ManifestSHA256: sv.Hash(read("odoo_loyalty/engine-lock.json"))}
	store, e := NewStoredValue(pool, profile, process)
	if e != nil {
		t.Fatal(e)
	}
	maker := identity.Principal{TenantID: tenant, Subject: "operator", Permissions: map[string]struct{}{"stored_value:request": {}, "stored_value:read": {}, "stored_value:approve": {}}, Organizations: map[string]struct{}{org: {}}}
	reviewer := maker
	reviewer.Subject = "reviewer"
	must(`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,'stored-value','Fixture','Fixture')`, tenant)
	must(`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,$2,$2,'Fixture','store')`, tenant, org)
	must(`insert into catalog.vehicle_model(tenant_id,model_id,model_code,display_name,vehicle_class,lifecycle_state) values($1,'m','m','Fixture','other','active')`, tenant)
	for _, v := range []string{"gift-50", "product-a"} {
		must(`insert into catalog.vehicle_variant(tenant_id,variant_id,model_id,variant_code,display_name,battery_specification,lifecycle_state) values($1,$2,'m',$2,'Fixture','{}','active')`, tenant, v)
	}
	order := func(id, state, product string, qty, price int64) {
		t.Helper()
		must(`insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version)values($1,$2,$3,'customer',$4,'ARS',$5,1)`, tenant, id, org, state, qty*price)
		must(`insert into sales.customer_order_line(tenant_id,order_id,line_id,variant_id,quantity,unit_price_minor_units)values($1,$2,'line',$3,$4,$5)`, tenant, id, product, qty, price)
	}
	request := func(id, op, order, program, account string, version int64) sv.Request {
		return sv.Request{OperationID: id, Operation: op, OrganizationID: org, OrderID: order, ProgramID: program, ExpectedOrderVersion: version, AccountID: account, ProfileSHA256: profile.SHA256()}
	}
	approve := func(v sv.ApprovalView) sv.ApprovalView {
		t.Helper()
		out, e := store.Decide(ctx, reviewer, sv.Decision{OrganizationID: org, ApprovalID: v.ApprovalID, PayloadSHA256: v.PayloadSHA256, Approved: true, Reason: "Fixture reviewed exact snapshot"})
		if e != nil {
			t.Fatal(e)
		}
		return out
	}
	order("gift-source", "confirmed", "gift-50", 2, 5000)
	giftRequest := request("gift-issue", "issue", "gift-source", "reference-gift_card", "", 1)
	proposal, e := store.Propose(ctx, maker, giftRequest)
	if e != nil {

		tx, ee := pool.Begin(ctx)
		if ee == nil {
			oo, st, aa, oe := store.order(ctx, tx, giftRequest)
			t.Logf("fixture diagnostics: valid=%v allowed=%v order=%+v state=%s allocation=%+v ordererr=%v", giftRequest.Validate(profile), store.allowed(maker, "stored_value:request"), oo, st, aa, oe)
			bb, be := store.calculate(ctx, tx, giftRequest, maker.Subject, time.Now().UTC().Add(time.Hour).Format(time.RFC3339Nano), false)
			t.Logf("calculation result=%+v entries=%+v err=%v", bb.Result, bb.Entries, be)
			tx.Rollback(ctx)
		}
		t.Fatal("issue proposal", e)
	}
	if len(proposal.Payload.Entries) != 2 || proposal.Payload.Entries[0].PointsDelta != "50.000000" {
		t.Fatal(proposal)
	}
	wrong := reviewer
	wrong.TenantID = "22222222-2222-4222-8222-222222222222"
	if _, e = store.Decide(ctx, wrong, sv.Decision{OrganizationID: org, ApprovalID: proposal.ApprovalID, PayloadSHA256: proposal.PayloadSHA256, Approved: true, Reason: "wrong"}); e == nil {
		t.Fatal("cross tenant approved")
	}
	if _, e = store.Decide(ctx, maker, sv.Decision{OrganizationID: org, ApprovalID: proposal.ApprovalID, PayloadSHA256: proposal.PayloadSHA256, Approved: true, Reason: "self"}); !errors.Is(e, approval.ErrSeparation) {
		t.Fatal("separation", e)
	}
	approved := approve(proposal)
	if approved.State != "approved" || string(approved.Receipt) == "null" {
		t.Fatal(approved)
	}
	replay, e := store.Propose(ctx, maker, giftRequest)
	if e != nil || !replay.Replay || replay.PayloadSHA256 != proposal.PayloadSHA256 {
		t.Fatal("replay", replay, e)
	}
	different := giftRequest
	different.OperationID = "gift-again"
	if _, e = store.Propose(ctx, maker, different); e == nil {
		t.Fatal("duplicate issuance")
	}
	account := proposal.Payload.Entries[0].AccountID
	other := proposal.Payload.Entries[1].AccountID
	if _, e = store.Propose(ctx, maker, request("same-order", "redeem", "gift-source", "reference-gift_card", account, 1)); e == nil {
		t.Fatal("future card on issuing order")
	}
	order("purchase-a", "placed", "product-a", 1, 11500)
	order("purchase-b", "placed", "product-a", 1, 11500)
	var wg sync.WaitGroup
	result := make(chan sv.ApprovalView, 2)
	errs := make(chan error, 2)
	for _, id := range []string{"a", "b"} {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			v, e := store.Propose(ctx, maker, request("spend-"+id, "redeem", "purchase-"+id, "reference-gift_card", account, 1))
			if e == nil {
				result <- v
			} else {
				errs <- e
			}
		}(id)
	}
	wg.Wait()
	close(result)
	close(errs)
	if len(result) != 1 || len(errs) != 1 {
		t.Fatal("double reservation", len(result), len(errs))
	}
	redeemed := approve(<-result)
	a, e := store.Allocation(ctx, maker, redeemed.Payload.Request.OrderID)
	if e != nil || a.GrossMinor != 11500 || a.GiftMinor != 5000 || a.ProviderDueMinor != 6500 || a.OrderVersion != 2 {
		t.Fatal(a, e)
	}
	// Rejecting a separate reservation returns availability without deleting history.
	order("small", "placed", "product-a", 1, 575)
	v, e := store.Propose(ctx, maker, request("spend-small-reject", "redeem", "small", "reference-gift_card", other, 1))
	if e != nil {
		t.Fatal(e)
	}
	_, e = store.Decide(ctx, reviewer, sv.Decision{OrganizationID: org, ApprovalID: v.ApprovalID, PayloadSHA256: v.PayloadSHA256, Approved: false, Reason: "Fixture reject"})
	if e != nil {
		t.Fatal(e)
	}
	v, e = store.Propose(ctx, maker, request("spend-small", "redeem", "small", "reference-gift_card", other, 1))
	if e != nil {
		t.Fatal(e)
	}
	approve(v)
	a, e = store.Allocation(ctx, maker, "small")
	if e != nil || a.ProviderDueMinor != 0 || a.GiftMinor != 575 {
		t.Fatal(a, e)
	}
	// Source cancellation is an immutable inverse, including negative balance
	// when a previously issued card was already consumed elsewhere.
	must(`update sales.customer_order set state='cancelled',version=version+1 where tenant_id=$1 and order_id='gift-source'`, tenant)
	reverse := request("gift-reverse", "reverse", "gift-source", "reference-gift_card", "", 2)
	reverse.OriginalOperationID = "gift-issue"
	v, e = store.Propose(ctx, maker, reverse)
	if e != nil {
		t.Fatal("reverse proposal", e)
	}
	approve(v)
	var balance string
	if e = pool.QueryRow(ctx, `select sum(points_delta)::numeric(38,6)::text from stored_value.entry where tenant_id=$1 and account_id=$2`, tenant, account).Scan(&balance); e != nil || balance != "-50.000000" {
		t.Fatal(balance, e)
	}
	reverse.OperationID = "gift-reverse-again"
	if v, e = store.Propose(ctx, maker, reverse); e == nil {
		if _, e = store.Decide(ctx, reviewer, sv.Decision{OrganizationID: org, ApprovalID: v.ApprovalID, PayloadSHA256: v.PayloadSHA256, Approved: true, Reason: "duplicate"}); e == nil {
			t.Fatal("double reversal")
		}
	}
	// Loyalty is bound to the customer and uses the admitted source default
	// money earning plus 5-percent/200-point reward, including minor rounding.
	order("earn", "confirmed", "product-a", 1, 20000)
	v, e = store.Propose(ctx, maker, request("loyalty-earn", "accrue", "earn", "reference-loyalty", "", 1))
	if e != nil {
		t.Fatal(e)
	}
	approve(v)
	loyalty := v.Payload.Entries[0].AccountID
	order("discount", "placed", "product-a", 1, 1001)
	v, e = store.Propose(ctx, maker, request("loyalty-spend", "redeem", "discount", "reference-loyalty", loyalty, 1))
	if e != nil {
		t.Fatal(e)
	}
	approve(v)
	a, e = store.Allocation(ctx, maker, "discount")
	if e != nil || a.DiscountMinor != 50 || a.ProviderDueMinor != 951 {
		t.Fatal(a, e)
	}

	// Seed a short-lived, correctly hash-bound proposal through the same typed
	// owner. Only the fixture lease is shorter than the deployment profile's
	// 60-second minimum; no production clock or validation is replaced.
	pending := func(id string, lifetime time.Duration) (sv.ApprovalView, time.Time) {
		t.Helper()
		order(id, "confirmed", "gift-50", 1, 5000)
		r := request(id, "issue", id, "reference-gift_card", "", 1)
		tx, e := pool.Begin(ctx)
		if e != nil {
			t.Fatal(e)
		}
		defer tx.Rollback(ctx)
		var expiry time.Time
		if e = tx.QueryRow(ctx, `select clock_timestamp()+($1::bigint*interval '1 millisecond')`, lifetime.Milliseconds()).Scan(&expiry); e != nil {
			t.Fatal(e)
		}
		bound, e := store.calculate(ctx, tx, r, maker.Subject, expiry.UTC().Format(time.RFC3339Nano), false)
		if e != nil {
			t.Fatal(e)
		}
		raw, hash, e := sv.Canonical(bound)
		if e != nil {
			t.Fatal(e)
		}
		id = sv.StableID("svapproval", id)
		_, e = store.approvals.submitTx(ctx, tx, maker, HumanApprovalSpec{Request: approval.Request{TenantID: tenant, ID: id, Kind: approval.KindStoredValueOperation, SubjectID: r.OrderID, Requester: maker.Subject, EvidenceSHA: hash}, OrganizationID: org, Payload: raw}, "stored_value:request", nil)
		if e != nil {
			t.Fatal(e)
		}
		if e = tx.Commit(ctx); e != nil {
			t.Fatal(e)
		}
		return sv.ApprovalView{ApprovalID: id, PayloadSHA256: hash, Payload: bound}, expiry
	}
	t.Run("expiry_during_outbox_wait", func(t *testing.T) {
		v, expiry := pending("expiry-outbox", 2*time.Second)
		blocker, e := pool.Begin(ctx)
		if e != nil {
			t.Fatal(e)
		}
		defer blocker.Rollback(ctx)
		if _, e = blocker.Exec(ctx, `lock table platform.outbox_event in access exclusive mode`); e != nil {
			t.Fatal(e)
		}
		done := make(chan error, 1)
		go func() {
			_, e := store.Decide(ctx, reviewer, sv.Decision{OrganizationID: org, ApprovalID: v.ApprovalID, PayloadSHA256: v.PayloadSHA256, Approved: true, Reason: "Exact short-lease fixture"})
			done <- e
		}()
		deadline := time.Now().Add(time.Second)
		for {
			var waiting bool
			if e = pool.QueryRow(ctx, `select exists(select 1 from pg_stat_activity where datname=current_database() and wait_event_type='Lock' and query like 'insert into platform.outbox_event%')`).Scan(&waiting); e != nil {
				t.Fatal(e)
			}
			if waiting {
				break
			}
			if time.Now().After(deadline) {
				t.Fatal("decision never reached blocked outbox")
			}
			time.Sleep(10 * time.Millisecond)
		}
		must(`select pg_sleep(greatest(0,extract(epoch from ($1::timestamptz-clock_timestamp())))+0.05)`, expiry)
		if e = blocker.Rollback(ctx); e != nil {
			t.Fatal(e)
		}
		if e = <-done; e == nil {
			t.Fatal("expired operation committed after outbox wait")
		}
		var operations, decisions, events int
		if e = pool.QueryRow(ctx, `select (select count(*) from stored_value.operation where tenant_id=$1 and approval_id=$2),(select count(*) from approval.decision where tenant_id=$1 and request_id=$2),(select count(*) from platform.outbox_event where tenant_id=$1 and aggregate_id='expiry-outbox')`, tenant, v.ApprovalID).Scan(&operations, &decisions, &events); e != nil || operations+decisions+events != 0 {
			t.Fatal("expired transaction left effects", operations, decisions, events, e)
		}
		got, e := store.Read(ctx, maker, v.ApprovalID)
		if e != nil || got.State != "pending" || string(got.Receipt) != "null" {
			t.Fatal(got, e)
		}
	})

	t.Run("proposal_expired_before_commit", func(t *testing.T) {
		order("expiry-proposal", "confirmed", "gift-50", 1, 5000)
		r := request("expiry-proposal", "issue", "expiry-proposal", "reference-gift_card", "", 1)
		tx, e := pool.Begin(ctx)
		if e != nil {
			t.Fatal(e)
		}
		defer tx.Rollback(ctx)
		var expiry time.Time
		if e = tx.QueryRow(ctx, `select clock_timestamp()+interval '1 second'`).Scan(&expiry); e != nil {
			t.Fatal(e)
		}
		bound, e := store.calculate(ctx, tx, r, maker.Subject, expiry.UTC().Format(time.RFC3339Nano), false)
		if e != nil {
			t.Fatal(e)
		}
		raw, hash, e := sv.Canonical(bound)
		if e != nil {
			t.Fatal(e)
		}
		id := sv.StableID("svapproval", r.OperationID)
		_, e = store.approvals.submitTx(ctx, tx, maker, HumanApprovalSpec{Request: approval.Request{TenantID: tenant, ID: id, Kind: approval.KindStoredValueOperation, SubjectID: r.OrderID, Requester: maker.Subject, EvidenceSHA: hash}, OrganizationID: org, Payload: raw}, "stored_value:request", nil)
		if e != nil {
			t.Fatal(e)
		}
		if _, e = tx.Exec(ctx, `select pg_sleep(greatest(0,extract(epoch from ($1::timestamptz-clock_timestamp())))+0.05)`, expiry); e != nil {
			t.Fatal(e)
		}
		if e = tx.Commit(ctx); e == nil {
			t.Fatal("expired proposal committed")
		}
		if _, e = store.Read(ctx, maker, id); !errors.Is(e, approval.ErrNotFound) {
			t.Fatal("expired proposal retained", e)
		}
	})
	// Immutable register, one event per applied operation, no second ledger.
	if _, e = pool.Exec(ctx, `update stored_value.entry set points_delta=0 where tenant_id=$1`, tenant); e == nil {
		t.Fatal("mutable history")
	}
	var counts []byte
	e = pool.QueryRow(ctx, `select jsonb_build_object('operations',(select count(*) from stored_value.operation),'events',(select count(*) from platform.outbox_event where event_type='stored-value.committed'),'approvals',(select count(*) from approval.request where kind='stored_value_operation'),'reservations',(select count(*) from stored_value.reservation))`).Scan(&counts)
	if e != nil {
		t.Fatal(e)
	}
	var numbers map[string]int
	if json.Unmarshal(counts, &numbers) != nil || numbers["operations"] != numbers["events"] {
		t.Fatal(string(counts))
	}
	t.Log(string(counts))
}
````

### FILE: `internal/platform/postgres/stored_value_read.go`

```yaml
block_id: "GO-APPROVED-STORED-VALUE-TENDER:file13:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration/configuration/authorization/persistence/transport/UI glue; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "bec65bf912d8e9d5a05f0541261f9df378d31db08ec27b95118933bf1f26d515"
variables: []
secrets_allowed: false
```

````go
package postgres

// AUTHORED scoped, bounded operator reads; monetary calculations remain in the
// source-derived process and allocations remain in the existing SQL projection.
import (
	"context"
	"elite.local/enterprise/internal/platform/identity"
	sv "elite.local/enterprise/internal/storedvaluebridge"
	"encoding/json"
)

func (s *StoredValue) ProfileDocument(ctx context.Context, p identity.Principal, org string) (json.RawMessage, error) {
	if !s.allowed(p, "stored_value:read") {
		return nil, sv.ErrBinding
	}
	_, bound := s.profile.Scope()
	if org != bound {
		return nil, sv.ErrBinding
	}
	return s.profile.Document(), nil
}

type StoredValueApprovalPage struct {
	IDs  []string `json:"ids"`
	Next string   `json:"next_cursor,omitempty"`
}

func (s *StoredValue) Approvals(ctx context.Context, p identity.Principal, order, after string) (StoredValueApprovalPage, error) {
	out := StoredValueApprovalPage{IDs: []string{}}
	if !s.allowed(p, "stored_value:read") || !sv.ValidID(order) || after != "" && !sv.ValidID(after) {
		return out, sv.ErrBinding
	}
	tenant, org := s.profile.Scope()
	rows, e := s.pool.Query(ctx, `select request_id from approval.request where tenant_id=$1 and organization_id=$2 and kind='stored_value_operation' and subject_id=$3 and request_id>$4 order by request_id limit 51`, tenant, org, order, after)
	if e != nil {
		return out, e
	}
	defer rows.Close()
	for rows.Next() {
		var id string
		if e = rows.Scan(&id); e != nil {
			return out, e
		}
		out.IDs = append(out.IDs, id)
	}
	if e = rows.Err(); e != nil {
		return out, e
	}
	if len(out.IDs) > 50 {
		out.IDs = out.IDs[:50]
		out.Next = out.IDs[49]
	}
	return out, nil
}

func (s *StoredValue) Scope() (string, string) {
	if s == nil {
		return "", ""
	}
	return s.profile.Scope()
}
````

### FILE: `internal/platform/postgres/stored_value_tender_integration_test.go`

```yaml
block_id: "GO-APPROVED-STORED-VALUE-TENDER:file14:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration/configuration/authorization/persistence/transport/UI glue; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "61061059ed4a9518273d5c79b872d48df69cfc2bb76c58dfbf2d563ce9614af3"
variables: []
secrets_allowed: false
```

````go
package postgres_test

// AUTHORED integration fixture: source-derived approved stored value changes
// only provider due, then the real SDK/callback/receipt journey completes.
import (
	"context"
	"elite.local/enterprise/internal/commerce"
	"elite.local/enterprise/internal/franchisejourney"
	"elite.local/enterprise/internal/platform/identity"
	db "elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/platform/randomid"
	sv "elite.local/enterprise/internal/storedvaluebridge"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"path/filepath"
	"testing"
)

func fixtureStoredValueSetup(t *testing.T, pool *pgxpool.Pool, tenant, kind string, giftValue ...int64) (*db.StoredValue, identity.Principal, string) {
	t.Helper()
	ctx := context.Background()
	root, e := filepath.Abs("../../..")
	if e != nil {
		t.Fatal(e)
	}
	read := func(rel string) []byte {
		b, e := os.ReadFile(filepath.Join(root, rel))
		if e != nil {
			t.Fatal(e)
		}
		return b
	}
	var c map[string]any
	if json.Unmarshal(read("deploy/stored-value/profile.reference.json"), &c) != nil {
		t.Fatal("profile JSON")
	}
	c["tenant_id"] = tenant
	c["organization_id"] = "store"
	raw, e := json.Marshal(c)
	if e != nil {
		t.Fatal(e)
	}
	profile, e := sv.Load(raw, sv.Hash(raw))
	if e != nil {
		t.Fatal(e)
	}
	python := os.Getenv("HANDOVER_PROFILE_PYTHON")
	if python == "" {
		t.Fatal("explicit fixed Python required")
	}
	process := sv.Process{Python: python, Script: filepath.Join(root, "odoo_loyalty/run.py"), ScriptSHA256: sv.Hash(read("odoo_loyalty/run.py")), Manifest: filepath.Join(root, "odoo_loyalty/engine-lock.json"), ManifestSHA256: sv.Hash(read("odoo_loyalty/engine-lock.json"))}
	store, e := db.NewStoredValue(pool, profile, process)
	if e != nil {
		t.Fatal(e)
	}
	actor := identity.Principal{TenantID: tenant, Subject: "operator", Organizations: map[string]struct{}{"store": {}}, Permissions: map[string]struct{}{"stored_value:read": {}, "stored_value:request": {}, "stored_value:approve": {}, "stored_value:fund": {}}}
	must := func(q string, args ...any) {
		t.Helper()
		if _, e := pool.Exec(ctx, q, args...); e != nil {
			t.Fatal(e)
		}
	}
	must(`insert into catalog.vehicle_variant(tenant_id,variant_id,model_id,variant_code,display_name,battery_specification,lifecycle_state)values($1,'gift-50','model','gift-50','Synthetic gift','{}','active')`, tenant)
	product, price := "gift-50", int64(5000)
	if len(giftValue) == 1 {
		price = giftValue[0]
	}
	if kind == "loyalty" {
		product, price = "variant", 20000
	}
	// The upstream issues/accrues at confirmed sale. This source-order fixture
	// is explicit; the target purchase uses actual quote/accept/stock writers.
	must(`insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version)values($1,'source','store','customer','confirmed','ARS',$2,1)`, tenant, price)
	must(`insert into sales.customer_order_line(tenant_id,order_id,line_id,variant_id,quantity,unit_price_minor_units)values($1,'source','source-line',$2,1,$3)`, tenant, product, price)
	return store, actor, "reference-" + kind
}

func fixtureStoredValue(t *testing.T, pool *pgxpool.Pool, tenant, kind string, giftValue ...int64) (*db.StoredValue, identity.Principal, string) {
	t.Helper()
	store, actor, program := fixtureStoredValueSetup(t, pool, tenant, kind, giftValue...)
	ctx := context.Background()
	review := actor
	review.Subject = "reviewer"
	op := "issue"
	if kind == "loyalty" {
		op = "accrue"
	}
	req := sv.Request{OperationID: "source-operation", Operation: op, OrganizationID: "store", OrderID: "source", ProgramID: program, ExpectedOrderVersion: 1, ProfileSHA256: store.ProfileSHA256()}
	v := storedValueHTTPPropose(t, store, actor, req)
	v = storedValueHTTPApprove(t, store, review, v)
	account := v.Payload.Entries[0].AccountID
	var version int64
	if e := pool.QueryRow(ctx, `select version from sales.customer_order where tenant_id=$1 and order_id='order'`, tenant).Scan(&version); e != nil {
		t.Fatal(e)
	}
	req = sv.Request{OperationID: "redemption", Operation: "redeem", OrganizationID: "store", OrderID: "order", ProgramID: program, ExpectedOrderVersion: version, AccountID: account, ProfileSHA256: store.ProfileSHA256()}
	v = storedValueHTTPPropose(t, store, actor, req)
	v = storedValueHTTPApprove(t, store, review, v)
	if v.State != "approved" {
		t.Fatal(v)
	}
	return store, actor, account
}

func TestStoredValueOfficialSDKTenderAndRelease(t *testing.T) {
	for _, kind := range []string{"gift_card", "loyalty"} {
		t.Run(kind, func(t *testing.T) {
			ctx := context.Background()
			var store *db.StoredValue
			var actor identity.Principal
			var account string
			var allocation sv.Allocation
			r := newConnectedRunWithAllocation(t, kind == "gift_card", func(t *testing.T, pool *pgxpool.Pool, tenant string) int64 {
				store, actor, account = fixtureStoredValue(t, pool, tenant, kind)
				var e error
				allocation, e = store.Allocation(ctx, actor, "order")
				if e != nil {
					t.Fatal(e)
				}
				expected := int64(118456)
				if kind == "loyalty" {
					expected = 117283
				}
				if allocation.GrossMinor != 123456 || allocation.ProviderDueMinor != expected {
					t.Fatal("gross/due", allocation)
				}
				// Legacy caller-supplied gross may not bypass the derived remaining due.
				_, e = db.NewCommerce(pool).RecordPaymentIntent(ctx, tenant, randomid.Generator{}.New(), "forged-gross-00001", commerce.PaymentAttempt{ID: "forged-gross", OrderID: "order", OrganizationID: "store", ProviderCode: "stripe", Currency: "ARS", AmountMinorUnits: 123456}, "operator")
				if e == nil {
					t.Fatal("gross charged after tender")
				}
				return expected
			})
			if code := r.callback(t, "evt_tender", "checkout.session.completed", "cs_test_fixture", true); code != 200 {
				t.Fatal(code)
			}
			if n, e := r.processor.ProcessOnce(ctx); e != nil || n != 1 {
				t.Fatal("reconciliation", n, e)
			}
			var amount, gross int64
			var state, hash string
			e := r.pool.QueryRow(ctx, `select p.amount_minor_units,o.total_minor_units,p.state,ob.evidence_sha256_hex from payment.payment_attempt p join sales.customer_order o using(tenant_id,order_id) join payment.provider_observation ob using(tenant_id,payment_attempt_id) where p.tenant_id=$1`, r.tenant).Scan(&amount, &gross, &state, &hash)
			if e != nil || amount != allocation.ProviderDueMinor || gross != 123456 || state != "captured" {
				t.Fatal("captured tender", amount, gross, state, e)
			}
			// Once provider intent exists, another allocation cannot change its amount.
			_, e = store.Propose(ctx, actor, sv.Request{OperationID: "after-payment", Operation: "redeem", OrganizationID: "store", OrderID: "order", ProgramID: "reference-" + kind, ExpectedOrderVersion: allocation.OrderVersion, AccountID: account, ProfileSHA256: store.ProfileSHA256()})
			if e == nil {
				t.Fatal("allocation changed after dispatch")
			}
			assertConnectedCommercialRelease(t, r, hash, true)
			t.Logf("gross=%d provider=%d gift=%d discount=%d; SDK callback, recovery, handover/commercial receipt and partial-refund invalidation passed", gross, amount, allocation.GiftMinor, allocation.DiscountMinor)
		})
	}
}

func TestStoredValueHandoverRejectsProviderOnlyProfile(t *testing.T) {
	ctx := context.Background()
	r := newConnectedRunWithAllocation(t, false, func(t *testing.T, pool *pgxpool.Pool, tenant string) int64 {
		store, actor, _ := fixtureStoredValue(t, pool, tenant, "gift_card")
		a, e := store.Allocation(ctx, actor, "order")
		if e != nil {
			t.Fatal(e)
		}
		return a.ProviderDueMinor
	})
	if code := r.callback(t, "evt_tender_policy", "checkout.session.completed", "cs_test_fixture", true); code != 200 {
		t.Fatal(code)
	}
	if n, e := r.processor.ProcessOnce(ctx); e != nil || n != 1 {
		t.Fatal(n, e)
	}
	var hash string
	if e := r.pool.QueryRow(ctx, `select evidence_sha256_hex from payment.provider_observation where tenant_id=$1`, r.tenant).Scan(&hash); e != nil {
		t.Fatal(e)
	}
	policy, _ := connectedCommercialProfile(t, r, false)
	svc, e := franchisejourney.NewHandoverPreparationService(db.NewFranchiseJourney(r.pool), randomid.Generator{}, policy)
	if e != nil {
		t.Fatal(e)
	}
	_, _, e = svc.Prepare(ctx, r.tenant, "operator", franchisejourney.PrepareHandoverCommand{OrganizationID: "store", OrderID: "order", OrderLineID: "line", PaymentAttemptID: r.payment.ID, ObservationSHA256: hash, IdempotencyKey: "provider-only-must-reject"})
	if e == nil {
		t.Fatal("provider-only profile accepted stored-value funding without explicit algorithm selection")
	}
	var n int
	if e = r.pool.QueryRow(ctx, `select count(*) from sales.delivery_handover_preparation where tenant_id=$1`, r.tenant).Scan(&n); e != nil || n != 0 {
		t.Fatal("unselected policy left an effect", n, e)
	}
}
````

### FILE: `internal/platform/postgres/zero_funding_integration_test.go`

```yaml
block_id: "GO-APPROVED-STORED-VALUE-TENDER:file15:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration/configuration/authorization/persistence/transport/UI glue; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "2e6342abe489765a28fb8ea9ab3fc29dd6bfa15fb089a31fb183bb5630b5f705"
variables: []
secrets_allowed: false
```

````go
package postgres_test

import (
	"context"
	"elite.local/enterprise/internal/commerce"
	"elite.local/enterprise/internal/franchisejourney"
	db "elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/platform/randomid"
	sv "elite.local/enterprise/internal/storedvaluebridge"
	"encoding/json"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestZeroProviderFundingReceiptThroughDelivery(t *testing.T) {
	ctx := context.Background()
	pool := connectedPool(t)
	tenant := connectedSeedOrder(t, pool, "stripe")
	store, actor, account := fixtureStoredValue(t, pool, tenant, "gift_card", 200000)
	allocation, e := store.Allocation(ctx, actor, "order")
	if e != nil || allocation.ProviderDueMinor != 0 || allocation.GiftMinor != 123456 {
		t.Fatal(allocation, e)
	}
	sales := commerce.NewService(db.NewCommerce(pool), randomid.Generator{})
	if _, e = sales.RequestOrderPayment(ctx, tenant, "store", "order", "stripe", "no-zero-provider-payment", "operator"); e == nil {
		t.Fatal("zero provider payment created")
	}
	// No provider setup can be used by this path, even though the profile records
	// the provider selection for other orders in the same local deployment.
	if _, e = pool.Exec(ctx, `update integration.provider_connection set state='disabled' where tenant_id=$1`, tenant); e != nil {
		t.Fatal(e)
	}
	command := db.FinalizeStoredValueFunding{OrganizationID: "store", OrderID: "order", ExpectedOrderVersion: allocation.OrderVersion, RequestKey: "full-funding-observation-0001"}
	denied := actor
	denied.Permissions = map[string]struct{}{"stored_value:read": {}}
	if _, _, e = store.FinalizeFunding(ctx, denied, command); e == nil {
		t.Fatal("missing funding permission accepted")
	}
	var wg sync.WaitGroup
	results := make(chan db.LocalFundingResult, 12)
	errors := make(chan error, 12)
	for range 12 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r, _, e := store.FinalizeFunding(ctx, actor, command)
			if e != nil {
				errors <- e
			} else {
				results <- r
			}
		}()
	}
	wg.Wait()
	close(results)
	close(errors)
	if len(errors) != 0 || len(results) != 12 {
		t.Fatal("concurrent funding", len(results), len(errors), <-errors)
	}
	var funding db.LocalFundingResult
	for r := range results {
		if funding.SHA256 != "" && r.SHA256 != funding.SHA256 {
			t.Fatal("duplicate observations")
		}
		funding = r
	}
	if funding.Receipt.Allocation.ProviderMinor != 0 || funding.Receipt.Effect != db.LocalFundingEffect {
		t.Fatal(funding)
	}
	got, e := store.FundingResult(ctx, actor, command.RequestKey)
	if e != nil || got.SHA256 != funding.SHA256 {
		t.Fatal("lost response", got, e)
	}
	changed := command
	changed.ExpectedOrderVersion++
	if _, _, e = store.FinalizeFunding(ctx, actor, changed); e == nil {
		t.Fatal("changed request replay")
	}
	stranger := actor
	stranger.Subject = "other-actor"
	if _, _, e = store.FinalizeFunding(ctx, stranger, command); e == nil {
		t.Fatal("other actor adopted receipt")
	}
	// A new observation stalled on outbox must roll back its receipt atomically.
	blocker, e := pool.Begin(ctx)
	if e != nil {
		t.Fatal(e)
	}
	if _, e = blocker.Exec(ctx, `lock table platform.outbox_event in access exclusive mode`); e != nil {
		t.Fatal(e)
	}
	short, cancel := context.WithTimeout(ctx, 150*time.Millisecond)
	blocked := command
	blocked.RequestKey = "full-funding-cancelled-0001"
	_, _, e = store.FinalizeFunding(short, actor, blocked)
	cancel()
	if e == nil {
		t.Fatal("blocked funding committed")
	}
	if e = blocker.Rollback(ctx); e != nil {
		t.Fatal(e)
	}
	var count int
	if e = pool.QueryRow(ctx, `select count(*) from payment.local_funding_receipt where tenant_id=$1`, tenant).Scan(&count); e != nil || count != 1 {
		t.Fatal("orphan funding", count, e)
	}
	r := &connectedRun{pool: pool, tenant: tenant}
	old, _ := connectedCommercialProfile(t, r, false)
	policy, _ := connectedCommercialProfile(t, r, true)
	repo := db.NewFranchiseJourney(pool)
	ids := randomid.Generator{}
	svc, e := franchisejourney.NewHandoverPreparationService(repo, ids, policy)
	if e != nil {
		t.Fatal(e)
	}
	prepare := franchisejourney.PrepareHandoverCommand{OrganizationID: "store", OrderID: "order", OrderLineID: "line", FundingReceiptID: funding.Receipt.ID, ObservationSHA256: funding.SHA256, IdempotencyKey: "local-funding-prepare-0001"}
	oldSvc, e := franchisejourney.NewHandoverPreparationService(repo, ids, old)
	if e != nil {
		t.Fatal(e)
	}
	if _, _, e = oldSvc.Prepare(ctx, tenant, "operator", prepare); e == nil {
		t.Fatal("legacy profile accepted local funding")
	}
	wrong := prepare
	wrong.PaymentAttemptID = "invented-provider-payment"
	if _, _, e = svc.Prepare(ctx, tenant, "operator", wrong); e == nil {
		t.Fatal("mixed identity accepted")
	}
	wrong = prepare
	wrong.ObservationSHA256 = strings.Repeat("0", 64)
	if _, _, e = svc.Prepare(ctx, tenant, "operator", wrong); e == nil {
		t.Fatal("forged receipt hash")
	}
	prepared, replay, e := svc.Prepare(ctx, tenant, "operator", prepare)
	if e != nil || replay || prepared.PaymentAttemptID != "" || prepared.FundingReceiptID != funding.Receipt.ID {
		t.Fatal("full funding preparation", prepared, e)
	}
	recovered, e := svc.Result(ctx, tenant, "store", "order", prepare.IdempotencyKey)
	if e != nil || recovered.Handover.ID != prepared.Handover.ID {
		t.Fatal("prepare recovery", e)
	}
	view, e := svc.OperatorContext(ctx, tenant, "store", "order")
	if e != nil || view.PaymentAttemptID != "" || view.FundingReceiptID != funding.Receipt.ID {
		t.Fatal("operator funding context", view, e)
	}
	if _, e = repo.PublishDeliveryChecklist(ctx, tenant, "operator", franchisejourney.DeliveryChecklist{ID: "local-funded-checklist", OrganizationID: "store", Version: 1, Title: "Synthetic delivery", Items: []franchisejourney.ChecklistItem{{ID: "serial", Ordinal: 1, Prompt: "Read serial", ResponseType: "serial", Required: true}}}, ids.New()); e != nil {
		t.Fatal(e)
	}
	presented, e := repo.CompleteDeliveryChecklist(ctx, tenant, "store", "operator", prepared.Handover.ID, 1, "local-funded-checklist", 1, []franchisejourney.ChecklistResponse{{ItemID: "serial", ResponseText: "SERIAL-SYNTHETIC"}}, ids.New())
	if e != nil {
		t.Fatal(e)
	}
	if _, e = repo.AcceptHandover(ctx, tenant, "store", "customer", prepared.Handover.ID, presented.Version, "SERIAL-SYNTHETIC", "local-funded-checklist", 1, strings.Repeat("e", 64), ids.New()); e != nil {
		t.Fatal(e)
	}
	commit := franchisejourney.CommitCommercialReleaseCommand{OrganizationID: "store", HandoverID: prepared.Handover.ID, ObservationSHA256: funding.SHA256, IdempotencyKey: "local-funding-release-0001"}
	receipt, replay, e := svc.CommitCommercialRelease(ctx, tenant, "operator", commit)
	if e != nil || replay || receipt.PaymentAttemptID != "" || receipt.FundingReceiptID != funding.Receipt.ID {
		t.Fatal("local funding release", receipt, e)
	}
	current, e := svc.ValidateCommercialRelease(ctx, tenant, "store", prepared.Handover.ID)
	if e != nil || !current.Current {
		t.Fatal("current receipt", current, e)
	}
	recoveredRelease, e := svc.CommercialReleaseResult(ctx, tenant, "store", prepared.Handover.ID, commit.IdempotencyKey)
	if e != nil || recoveredRelease.ID != receipt.ID {
		t.Fatal("release recovery", e)
	}
	// Negative source-state fixture, not a claim that a customer cancellation
	// endpoint is permitted after physical delivery. Existing truth invalidates.
	if _, e = pool.Exec(ctx, `update sales.customer_order set state='cancelled',version=version+1 where tenant_id=$1 and order_id='order'`, tenant); e != nil {
		t.Fatal(e)
	}
	current, e = svc.ValidateCommercialRelease(ctx, tenant, "store", prepared.Handover.ID)
	if e != nil || current.Current {
		t.Fatal("cancelled order remained releasable", e)
	}
	var version int64
	if e = pool.QueryRow(ctx, `select version from sales.customer_order where tenant_id=$1 and order_id='order'`, tenant).Scan(&version); e != nil {
		t.Fatal(e)
	}
	v, e := store.Propose(ctx, actor, sv.Request{OperationID: "reverse-funded-redemption", Operation: "reverse", OrganizationID: "store", OrderID: "order", ProgramID: "reference-gift_card", ExpectedOrderVersion: version, OriginalOperationID: "redemption", ProfileSHA256: store.ProfileSHA256()})
	if e != nil {
		t.Fatal("approved cancellation inverse", e)
	}
	reviewer := actor
	reviewer.Subject = "reviewer"
	if _, e = store.Decide(ctx, reviewer, sv.Decision{OrganizationID: "store", ApprovalID: v.ApprovalID, PayloadSHA256: v.PayloadSHA256, Approved: true, Reason: "Exact source inverse after fixture cancellation"}); e != nil {
		t.Fatal(e)
	}
	var balance string
	if e = pool.QueryRow(ctx, `select sum(points_delta)::numeric(38,6)::text from stored_value.entry where tenant_id=$1 and account_id=$2`, tenant, account).Scan(&balance); e != nil || balance != "2000.000000" {
		t.Fatal("gift conservation", balance, e)
	}
	history, e := store.FundingResult(ctx, actor, command.RequestKey)
	if e != nil || history.SHA256 != funding.SHA256 {
		t.Fatal("historical receipt changed", e)
	}
	var snapshot []byte
	e = pool.QueryRow(ctx, `select jsonb_build_object('payments',(select count(*) from payment.payment_attempt where tenant_id=$1),'observations',(select count(*) from payment.provider_observation where tenant_id=$1),'funding',(select count(*) from payment.local_funding_receipt where tenant_id=$1),'funding_events',(select count(*) from platform.outbox_event where tenant_id=$1 and event_type='order.local-funding-observed'),'releases',(select count(*) from sales.commercial_release_receipt where tenant_id=$1))`, tenant).Scan(&snapshot)
	var totals map[string]int
	if e != nil || json.Unmarshal(snapshot, &totals) != nil || totals["payments"] != 0 || totals["observations"] != 0 || totals["funding"] != 1 || totals["funding_events"] != 1 || totals["releases"] != 1 {
		t.Fatal("no fake payment", string(snapshot), e)
	}
	t.Log(string(snapshot))
}
````

### FILE: `internal/storedvaluebridge/contract.go`

```yaml
block_id: "GO-APPROVED-STORED-VALUE-TENDER:file16:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration/configuration/authorization/persistence/transport/UI glue; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "9b06c217fe2a5baef3da4dcee1d418a387bf97d3b516cde16559e42badc0840b"
variables: []
secrets_allowed: false
```

````go
package storedvaluebridge

// AUTHORED contracts and exact serialization glue. Commercial calculations are
// provided by the separately licensed source-derived Odoo process.
import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"math/big"
	"regexp"
	"strconv"
	"strings"

	"elite.local/enterprise/internal/approval"
)

var ErrBinding = errors.New("stored value binding is invalid")
var identityRE = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_.:-]{0,127}$`)
var hashRE = regexp.MustCompile(`^[0-9a-f]{64}$`)
var decimalRE = regexp.MustCompile(`^-?(0|[1-9][0-9]{0,17})(\.[0-9]{1,6})?$`)

func Hash(b []byte) string  { h := sha256.Sum256(b); return hex.EncodeToString(h[:]) }
func ValidID(v string) bool { return identityRE.MatchString(v) }
func Decode(raw []byte, target any) error {
	canonical, _, e := approval.CanonicalPayload(raw)
	if e != nil {
		return ErrBinding
	}
	d := json.NewDecoder(bytes.NewReader(raw))
	d.DisallowUnknownFields()
	if e := d.Decode(target); e != nil {
		return ErrBinding
	}
	normalized, _, e := Canonical(target)
	if e != nil || !bytes.Equal(canonical, normalized) {
		return ErrBinding
	}
	return nil
}
func Canonical(v any) (json.RawMessage, string, error) {
	b, e := json.Marshal(v)
	if e != nil {
		return nil, "", e
	}
	return approval.CanonicalPayload(b)
}
func Minor(value string, digits int) (int64, error) {
	if !decimalRE.MatchString(value) || digits < 0 || digits > 6 {
		return 0, ErrBinding
	}
	r, ok := new(big.Rat).SetString(value)
	if !ok {
		return 0, ErrBinding
	}
	scale := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(digits)), nil)
	r.Mul(r, new(big.Rat).SetInt(scale))
	if !r.IsInt() || !r.Num().IsInt64() {
		return 0, ErrBinding
	}
	return r.Num().Int64(), nil
}
func Major(value int64, digits int) string {
	n := big.NewInt(value)
	negative := n.Sign() < 0
	n.Abs(n)
	v := n.String()
	if digits > 0 {
		for len(v) <= digits {
			v = "0" + v
		}
		v = v[:len(v)-digits] + "." + v[len(v)-digits:]
	}
	if negative {
		v = "-" + v
	}
	return v
}

type Rule struct {
	ID            string   `json:"id"`
	Mode          string   `json:"mode"`
	Points        string   `json:"points"`
	Split         bool     `json:"split"`
	MinimumQty    string   `json:"minimum_qty"`
	MinimumAmount string   `json:"minimum_amount"`
	TaxMode       string   `json:"tax_mode"`
	ProductIDs    []string `json:"product_ids"`
	CodeRequired  bool     `json:"code_required"`
}
type Program struct {
	ID                string   `json:"id"`
	Kind              string   `json:"kind"`
	Currency          string   `json:"currency"`
	CurrencyDigits    int      `json:"currency_digits"`
	AppliesOn         string   `json:"applies_on"`
	Nominative        bool     `json:"nominative"`
	Trigger           string   `json:"trigger"`
	TriggerProductIDs []string `json:"trigger_product_ids"`
	Rules             []Rule   `json:"rules"`
}
type Reward struct {
	Mode           string `json:"mode"`
	Discount       string `json:"discount"`
	RequiredPoints string `json:"required_points"`
	MaxAmount      string `json:"max_amount"`
	ClearWallet    bool   `json:"clear_wallet"`
}
type PolicyProgram struct {
	Calculation Program `json:"calculation"`
	Reward      Reward  `json:"reward"`
}
type ProfileConfig struct {
	Schema             string          `json:"schema"`
	TenantID           string          `json:"tenant_id"`
	OrganizationID     string          `json:"organization_id"`
	AmountsMode        string          `json:"amounts_mode"`
	Review             string          `json:"review"`
	ReservationSeconds int             `json:"reservation_seconds"`
	IssuanceStates     []string        `json:"issuance_states"`
	RedemptionStates   []string        `json:"redemption_states"`
	Programs           []PolicyProgram `json:"programs"`
}
type Profile struct {
	config ProfileConfig
	raw    json.RawMessage
	hash   string
}

func Load(raw []byte, expected string) (*Profile, error) {
	var c ProfileConfig
	if Decode(raw, &c) != nil || !hashRE.MatchString(expected) || Hash(raw) != expected {
		return nil, ErrBinding
	}
	if c.Schema != "elite.stored-value-profile.v1" || !ValidID(c.TenantID) || !ValidID(c.OrganizationID) || c.AmountsMode != "bound_gross_inclusive" || c.Review != "one_distinct_human" || c.ReservationSeconds < 60 || c.ReservationSeconds > 86400 || len(c.Programs) == 0 || len(c.Programs) > 8 {
		return nil, ErrBinding
	}
	validStates := func(states []string, allowed string) bool {
		if len(states) == 0 {
			return false
		}
		seen := map[string]bool{}
		for _, s := range states {
			if seen[s] || !strings.Contains("|"+allowed+"|", "|"+s+"|") {
				return false
			}
			seen[s] = true
		}
		return true
	}
	if !validStates(c.IssuanceStates, "confirmed|paid|allocated|delivered") || !validStates(c.RedemptionStates, "placed|confirmed|allocated") {
		return nil, ErrBinding
	}
	seen := map[string]bool{}
	kinds := map[string]bool{}
	for _, p := range c.Programs {
		v := p.Calculation
		if !ValidID(v.ID) || seen[v.ID] || kinds[v.Kind] || (v.Kind != "gift_card" && v.Kind != "loyalty") || v.Currency != "ARS" || v.CurrencyDigits != 2 || len(v.Rules) == 0 || len(v.Rules) > 20 {
			return nil, ErrBinding
		}
		seen[v.ID] = true
		kinds[v.Kind] = true
		if (v.Kind == "gift_card" && (v.AppliesOn != "future" || v.Nominative)) || (v.Kind == "loyalty" && (v.AppliesOn != "both" || !v.Nominative)) {
			return nil, ErrBinding
		}
		if v.Trigger != "auto" || v.TriggerProductIDs == nil {
			return nil, ErrBinding
		}
		positive := func(v string, zero bool) bool {
			r, ok := new(big.Rat).SetString(v)
			return ok && decimalRE.MatchString(v) && (r.Sign() > 0 || (zero && r.Sign() == 0))
		}
		ruleIDs := map[string]bool{}
		for _, r := range v.Rules {
			if !ValidID(r.ID) || ruleIDs[r.ID] || r.TaxMode != "incl" || r.CodeRequired || r.ProductIDs == nil || !positive(r.Points, false) || !positive(r.MinimumQty, true) || !positive(r.MinimumAmount, true) || (r.Mode != "order" && r.Mode != "money" && r.Mode != "unit") || (r.Split && v.AppliesOn == "both") {
				return nil, ErrBinding
			}
			ruleIDs[r.ID] = true
		}
		if p.Reward.Mode != "per_point" && p.Reward.Mode != "per_order" && p.Reward.Mode != "percent" {
			return nil, ErrBinding
		}
		if !positive(p.Reward.Discount, false) || !positive(p.Reward.RequiredPoints, false) || !positive(p.Reward.MaxAmount, true) {
			return nil, ErrBinding
		}
	}
	return &Profile{c, append(json.RawMessage(nil), raw...), expected}, nil
}
func (p *Profile) Valid() bool {
	return p != nil && hashRE.MatchString(p.hash) && Hash(p.raw) == p.hash
}
func (p *Profile) SHA256() string {
	if p == nil {
		return ""
	}
	return p.hash
}
func (p *Profile) Scope() (string, string) {
	if !p.Valid() {
		return "", ""
	}
	return p.config.TenantID, p.config.OrganizationID
}
func (p *Profile) Program(id string) (PolicyProgram, error) {
	if !p.Valid() {
		return PolicyProgram{}, ErrBinding
	}
	for _, v := range p.config.Programs {
		if v.Calculation.ID == id {
			// Keep the validated hash-locked profile immutable across caller DTO changes.
			b, err := json.Marshal(v)
			if err != nil {
				return PolicyProgram{}, ErrBinding
			}
			var copied PolicyProgram
			if json.Unmarshal(b, &copied) != nil {
				return PolicyProgram{}, ErrBinding
			}
			return copied, nil
		}
	}
	return PolicyProgram{}, ErrBinding
}
func (p *Profile) ReservationSeconds() int {
	if !p.Valid() {
		return 0
	}
	return p.config.ReservationSeconds
}
func (p *Profile) StateAllowed(operation, state string) bool {
	if !p.Valid() {
		return false
	}
	states := p.config.RedemptionStates
	if operation == "issue" || operation == "accrue" {
		states = p.config.IssuanceStates
	}
	for _, v := range states {
		if v == state {
			return true
		}
	}
	return false
}

type Line struct {
	ID                string `json:"id"`
	ProductID         string `json:"product_id"`
	Quantity          string `json:"quantity"`
	Subtotal          string `json:"subtotal"`
	Tax               string `json:"tax"`
	Total             string `json:"total"`
	RewardProgramType string `json:"reward_program_type"`
	RewardProgramID   string `json:"reward_program_id"`
	RewardTrigger     string `json:"reward_trigger"`
	ThresholdExcluded bool   `json:"threshold_excluded"`
}
type Order struct {
	OrderID        string   `json:"order_id"`
	OrganizationID string   `json:"organization_id"`
	SubjectID      string   `json:"subject_id"`
	State          string   `json:"state"`
	PublicSubject  bool     `json:"public_subject"`
	Currency       string   `json:"currency"`
	Total          string   `json:"total"`
	EnabledRuleIDs []string `json:"enabled_rule_ids"`
	Lines          []Line   `json:"lines"`
}
type Calculation struct {
	Schema        string          `json:"schema"`
	Operation     string          `json:"operation"`
	Program       Program         `json:"program"`
	ProgramSHA256 string          `json:"program_sha256"`
	Order         Order           `json:"order"`
	Data          json.RawMessage `json:"data"`
}

func NewCalculation(operation string, p Program, o Order, data any) (Calculation, error) {
	_, hash, e := Canonical(p)
	if e != nil {
		return Calculation{}, e
	}
	raw, e := json.Marshal(data)
	return Calculation{"elite.odoo-loyalty-calc.v1", operation, p, hash, o, raw}, e
}

type Result struct {
	Schema         string          `json:"schema"`
	SourceRevision string          `json:"source_revision"`
	ProgramSHA256  string          `json:"program_sha256"`
	RequestSHA256  string          `json:"request_sha256"`
	Result         json.RawMessage `json:"result"`
}
type Request struct {
	OperationID          string `json:"operation_id"`
	Operation            string `json:"operation"`
	OrganizationID       string `json:"organization_id"`
	ProgramID            string `json:"program_id"`
	OrderID              string `json:"order_id"`
	ExpectedOrderVersion int64  `json:"expected_order_version"`
	AccountID            string `json:"account_id"`
	OriginalOperationID  string `json:"original_operation_id"`
	ProfileSHA256        string `json:"profile_sha256"`
}

func (r Request) Validate(p *Profile) error {
	if !p.Valid() || !ValidID(r.OperationID) || !ValidID(r.OrderID) || r.ExpectedOrderVersion < 1 || r.OrganizationID != p.config.OrganizationID || r.ProfileSHA256 != p.hash {
		return ErrBinding
	}
	if _, e := p.Program(r.ProgramID); e != nil {
		return e
	}
	switch r.Operation {
	case "issue", "accrue":
		if r.AccountID != "" || r.OriginalOperationID != "" {
			return ErrBinding
		}
	case "redeem":
		if !ValidID(r.AccountID) || r.OriginalOperationID != "" {
			return ErrBinding
		}
	case "reverse":
		if !ValidID(r.OriginalOperationID) || r.AccountID != "" {
			return ErrBinding
		}
	default:
		return ErrBinding
	}
	return nil
}
func StableID(kind string, parts ...string) string {
	return kind + "_" + Hash([]byte(strings.Join(parts, "\x00")))
}
func Number(v int64) string { return strconv.FormatInt(v, 10) }

// Document returns a copy of the selected immutable policy for authorized review.
func (p *Profile) Document() json.RawMessage {
	if !p.Valid() {
		return nil
	}
	return append(json.RawMessage(nil), p.raw...)
}
````

### FILE: `internal/storedvaluebridge/contract_test.go`

```yaml
block_id: "GO-APPROVED-STORED-VALUE-TENDER:file17:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration/configuration/authorization/persistence/transport/UI glue; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "edc15b637120670783a8ba5671eca88390f17f78f0968bc82361e54a492712ae"
variables: []
secrets_allowed: false
```

````go
package storedvaluebridge

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func fixture(t *testing.T) (*Profile, Process) {
	t.Helper()
	_, file, _, _ := runtime.Caller(0)
	root := filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
	raw, e := os.ReadFile(filepath.Join(root, "deploy", "stored-value", "profile.reference.json"))
	if e != nil {
		t.Fatal(e)
	}
	p, e := Load(raw, Hash(raw))
	if e != nil {
		t.Fatal(e)
	}
	script := filepath.Join(root, "odoo_loyalty", "run.py")
	manifest := filepath.Join(root, "odoo_loyalty", "engine-lock.json")
	scriptBytes, _ := os.ReadFile(script)
	manifestBytes, _ := os.ReadFile(manifest)
	return p, Process{Python: `C:\Python314\python.exe`, Script: script, ScriptSHA256: Hash(scriptBytes), Manifest: manifest, ManifestSHA256: Hash(manifestBytes)}
}
func TestExactAmountsAndProfileContracts(t *testing.T) {
	for _, n := range []int64{0, 1, 100, -25, 9223372036854775807, -9223372036854775808} {
		v, e := Minor(Major(n, 2), 2)
		if e != nil || v != n {
			t.Fatalf("roundtrip %d: %d %v", n, v, e)
		}
	}
	for _, v := range []string{"0.001", "NaN", "1e2", "92233720368547758.08"} {
		if _, e := Minor(v, 2); e == nil {
			t.Fatalf("accepted %s", v)
		}
	}
	p, _ := fixture(t)
	raw := p.raw
	for _, bad := range [][]byte{[]byte(strings.Replace(string(raw), `"currency_digits": 2`, `"currency_digits": null`, 1)), []byte(strings.Replace(string(raw), `"points": "1"`, `"points": "0"`, 1)), []byte(strings.Replace(string(raw), `"schema"`, `"Schema"`, 1)), []byte(strings.Replace(string(raw), `"tax_mode": "incl"`, `"tax_mode": "excl"`, 1))} {
		if _, e := Load(bad, Hash(bad)); e == nil {
			t.Fatal("accepted ambiguous or unsupported profile")
		}
	}
	if _, e := Load(raw, strings.Repeat("0", 64)); e == nil {
		t.Fatal("profile hash ignored")
	}
}
func TestActualIsolatedProcessAndSourceBinding(t *testing.T) {
	p, process := fixture(t)
	if _, e := os.Stat(process.Python); e != nil {
		t.Fatal("qualified Python runtime missing")
	}
	if e := process.Preflight(context.Background(), p); e != nil {
		t.Fatal(e)
	}
	entry, _ := p.Program("reference-gift_card")
	order := Order{OrderID: "order-1", OrganizationID: "franchise-1", SubjectID: "customer-1", State: "draft", Currency: "ARS", Total: "100", EnabledRuleIDs: []string{}, Lines: []Line{{ID: "line-1", ProductID: "gift-50", Quantity: "2", Subtotal: "100", Tax: "0", Total: "100"}}}
	request, e := NewCalculation("evaluate", entry.Calculation, order, map[string]any{})
	if e != nil {
		t.Fatal(e)
	}
	result, e := process.Execute(context.Background(), request)
	if e != nil {
		t.Fatal(e)
	}
	points, e := ReadPoints(result)
	if e != nil || strings.Join(points, ",") != "50.00,50.00" {
		t.Fatalf("%s %v", result.Result, e)
	}
	bad := process
	bad.ManifestSHA256 = strings.Repeat("0", 64)
	if _, e = bad.Execute(context.Background(), request); e == nil {
		t.Fatal("changed manifest binding accepted")
	}
	// Same manifest bytes with a changed module in an isolated copied payload must fail before calculation.
	temp := t.TempDir()
	base := filepath.Dir(process.Script)
	e = filepath.WalkDir(base, func(path string, d os.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if d.IsDir() {
			return nil
		}
		rel, _ := filepath.Rel(base, path)
		dest := filepath.Join(temp, rel)
		if e = os.MkdirAll(filepath.Dir(dest), 0700); e != nil {
			return e
		}
		body, e := os.ReadFile(path)
		if e != nil {
			return e
		}
		return os.WriteFile(dest, body, 0600)
	})
	if e != nil {
		t.Fatal(e)
	}
	if e = os.WriteFile(filepath.Join(temp, "engine.py"), []byte("raise RuntimeError('tampered')\n"), 0600); e != nil {
		t.Fatal(e)
	}
	bad = process
	bad.Script = filepath.Join(temp, "run.py")
	bad.Manifest = filepath.Join(temp, "engine-lock.json")
	if _, e = bad.Execute(context.Background(), request); e == nil {
		t.Fatal("changed source bytes accepted")
	}
}
func TestSourceProfileChangeChangesCalculatedResult(t *testing.T) {
	p, process := fixture(t)
	entry, _ := p.Program("reference-loyalty")
	o := Order{OrderID: "order-1", OrganizationID: "franchise-1", SubjectID: "customer-1", State: "draft", Currency: "ARS", Total: "100", EnabledRuleIDs: []string{}, Lines: []Line{{ID: "line-1", ProductID: "sku", Quantity: "2", Subtotal: "100", Tax: "0", Total: "100"}}}
	first, _ := NewCalculation("evaluate", entry.Calculation, o, map[string]any{})
	a, e := process.Execute(context.Background(), first)
	if e != nil {
		t.Fatal(e)
	}
	entry.Calculation.Rules[0].Points = "2"
	unchanged, _ := p.Program("reference-loyalty")
	if unchanged.Calculation.Rules[0].Points != "1" {
		t.Fatal("mutated immutable profile")
	}
	second, _ := NewCalculation("evaluate", entry.Calculation, o, map[string]any{})
	b, e := process.Execute(context.Background(), second)
	if e != nil {
		t.Fatal(e)
	}
	if string(a.Result) == string(b.Result) || a.ProgramSHA256 == b.ProgramSHA256 {
		t.Fatal("configuration change not reflected")
	}
	var result struct{ Points []string }
	if json.Unmarshal(b.Result, &result) != nil || len(result.Points) != 1 || result.Points[0] != "200.00" {
		t.Fatalf("%s", b.Result)
	}
}

func TestPointsRepresentationAndReceiptIntegerBinding(t *testing.T) {
	for _, v := range []string{"50.00", "-0.250000", "200.000001"} {
		if _, e := Minor(v, 6); e != nil {
			t.Fatal(v, e)
		}
	}
	var c Calculation
	raw := []byte(`{"schema":"x","operation":"reverse","program":{"id":"","kind":"","currency":"","currency_digits":0,"applies_on":"","nominative":false,"trigger":"","trigger_product_ids":null,"rules":null},"program_sha256":"","order":{"order_id":"","organization_id":"","subject_id":"","state":"","public_subject":false,"currency":"","total":"","enabled_rule_ids":null,"lines":null},"data":{"applied_minor_units":9007199254740993}}`)
	if Decode(raw, &c) != nil || !strings.Contains(string(c.Data), "9007199254740993") {
		t.Fatal("approval data lost integer precision")
	}
}
````

### FILE: `internal/storedvaluebridge/fuzz_test.go`

```yaml
block_id: "GO-APPROVED-STORED-VALUE-TENDER:file18:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration/configuration/authorization/persistence/transport/UI glue; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "dba7a092cd680a244196d4b56488fde83fa1cd02ef931548dd476d39bf7193a7"
variables: []
secrets_allowed: false
```

````go
package storedvaluebridge

// AUTHORED domain-bound representation invariants, no external effects.
import (
	"math"
	"os"
	"path/filepath"
	"testing"
)

func FuzzStoredValueMinorRoundTrip(f *testing.F) {
	for _, v := range []int64{0, 1, -1, 50, 5000, 123456, 9007199254740993, math.MaxInt64, math.MinInt64} {
		f.Add(v, uint8(2))
	}
	f.Fuzz(func(t *testing.T, v int64, d uint8) {
		digits := int(d % 7)
		raw := Major(v, digits)
		got, e := Minor(raw, digits)
		// The wire decimal contract admits at most18whole digits, so int64 extrema
		// at scale0 are correctly rejected; all representable values round-trip.
		if digits == 0 && (v > 999999999999999999 || v < -999999999999999999) {
			if e == nil {
				t.Fatal("unbounded decimal")
			}
			return
		}
		if e != nil || got != v {
			t.Fatal("amount representation changed", raw, digits, got, e)
		}
	})
}
func FuzzStoredValueProfileBinding(f *testing.F) {
	raw, e := os.ReadFile(filepath.Join("..", "..", "deploy", "stored-value", "profile.reference.json"))
	if e != nil {
		f.Fatal(e)
	}
	f.Add(raw)
	f.Add([]byte(`{"schema":"elite.stored-value-profile.v1","schema":"duplicate"}`))
	f.Add([]byte("null"))
	f.Fuzz(func(t *testing.T, raw []byte) {
		if len(raw) > 8192 {
			return
		}
		p, e := Load(raw, Hash(raw))
		if e != nil {
			return
		}
		if !p.Valid() || Hash(p.Document()) != p.SHA256() {
			t.Fatal("mutable or unbound profile")
		}
		copy := p.Document()
		if len(copy) > 0 {
			copy[0] ^= 1
		}
		if !p.Valid() {
			t.Fatal("document aliases live policy")
		}
	})
}
````

### FILE: `internal/storedvaluebridge/persistence.go`

```yaml
block_id: "GO-APPROVED-STORED-VALUE-TENDER:file19:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration/configuration/authorization/persistence/transport/UI glue; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "68e6cfa6f5a8926a30c0e54f2a981f216f735583e9593d7d65fa5794d98b59b4"
variables: []
secrets_allowed: false
```

````go
package storedvaluebridge

// AUTHORED typed transaction and receipt contracts. No pricing or points rule
// is implemented here; Result and Calculation identify the source-derived IPC.
import "encoding/json"

type Entry struct {
	AccountID    string `json:"account_id"`
	PointsDelta  string `json:"points_delta"`
	AppliedMinor int64  `json:"applied_minor_units"`
}
type BoundOperation struct {
	Schema        string      `json:"schema"`
	Request       Request     `json:"request"`
	Requester     string      `json:"requester"`
	Order         Order       `json:"order"`
	OrderState    string      `json:"order_state"`
	GrossMinor    int64       `json:"gross_minor_units"`
	GiftMinor     int64       `json:"gift_minor_units"`
	DiscountMinor int64       `json:"discount_minor_units"`
	EngineSHA256  string      `json:"engine_sha256"`
	ExpiresAt     string      `json:"expires_at"`
	Calculation   Calculation `json:"calculation"`
	Result        Result      `json:"result"`
	Entries       []Entry     `json:"entries"`
}
type ApprovalView struct {
	ApprovalID    string          `json:"approval_id"`
	PayloadSHA256 string          `json:"payload_sha256"`
	State         string          `json:"state"`
	Replay        bool            `json:"replay"`
	Payload       BoundOperation  `json:"payload"`
	Receipt       json.RawMessage `json:"receipt"`
}
type Decision struct {
	OrganizationID string `json:"organization_id"`
	ApprovalID     string `json:"approval_id"`
	PayloadSHA256  string `json:"payload_sha256"`
	Approved       bool   `json:"approved"`
	Reason         string `json:"reason"`
}
type Allocation struct {
	OrderID          string `json:"order_id"`
	OrderVersion     int64  `json:"order_version"`
	Currency         string `json:"currency"`
	GrossMinor       int64  `json:"gross_minor_units"`
	GiftMinor        int64  `json:"gift_minor_units"`
	DiscountMinor    int64  `json:"discount_minor_units"`
	ProviderDueMinor int64  `json:"provider_due_minor_units"`
}
````

### FILE: `internal/storedvaluebridge/process.go`

```yaml
block_id: "GO-APPROVED-STORED-VALUE-TENDER:file20:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration/configuration/authorization/persistence/transport/UI glue; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "be467b3d2cc6baf380e7bc52bcb7498b773d444dfbdaa0c95800075949fe1996"
variables: []
secrets_allowed: false
```

````go
package storedvaluebridge

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// Process verifies the materialized source manifest before each bounded pure
// calculation. No credentials, inherited Python path, HTTP or DB are exposed.
type Process struct{ Python, Script, ScriptSHA256, Manifest, ManifestSHA256 string }

func (p Process) Validate() error {
	if !filepath.IsAbs(p.Python) || !filepath.IsAbs(p.Script) || !filepath.IsAbs(p.Manifest) || !hashRE.MatchString(p.ScriptSHA256) || !hashRE.MatchString(p.ManifestSHA256) {
		return ErrBinding
	}
	if info, e := os.Stat(p.Python); e != nil || info.IsDir() {
		return ErrBinding
	}
	for path, expected := range map[string]string{p.Script: p.ScriptSHA256, p.Manifest: p.ManifestSHA256} {
		b, e := os.ReadFile(path)
		if e != nil || len(b) > 65536 || Hash(b) != expected {
			return ErrBinding
		}
	}
	return nil
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
func (p Process) Execute(ctx context.Context, request Calculation) (Result, error) {
	var result Result
	if p.Validate() != nil {
		return result, ErrBinding
	}
	input, expected, e := Canonical(request)
	if e != nil {
		return result, e
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	c := exec.CommandContext(ctx, p.Python, "-I", "-S", "-B", p.Script, p.ManifestSHA256)
	c.Env = []string{}
	for _, key := range []string{"SYSTEMROOT", "WINDIR", "TEMP", "TMP"} {
		if v := os.Getenv(key); v != "" {
			c.Env = append(c.Env, key+"="+v)
		}
	}
	c.Stdin = bytes.NewReader(input)
	var out output
	c.Stdout = &out
	c.Stderr = io.Discard
	if c.Run() != nil || out.exceeded || Decode(out.Bytes(), &result) != nil {
		return result, ErrBinding
	}
	if result.Schema != "elite.odoo-loyalty-result.v1" || result.SourceRevision != "99edb6dd82b7b560930c00b03b694ba700785370" || result.ProgramSHA256 != request.ProgramSHA256 || result.RequestSHA256 != expected {
		return Result{}, ErrBinding
	}
	return result, nil
}
func (p Process) Preflight(ctx context.Context, profile *Profile) error {
	if p.Validate() != nil || !profile.Valid() {
		return ErrBinding
	}
	for _, v := range profile.config.Programs {
		o := Order{OrderID: "preflight", OrganizationID: profile.config.OrganizationID, SubjectID: "preflight", State: "draft", Currency: v.Calculation.Currency, Total: "0", EnabledRuleIDs: []string{}, Lines: []Line{{ID: "preflight", ProductID: "preflight", Quantity: "1", Subtotal: "0", Tax: "0", Total: "0"}}}
		r, e := NewCalculation("evaluate", v.Calculation, o, map[string]any{})
		if e != nil {
			return e
		}
		if _, e = p.Execute(ctx, r); e != nil {
			return e
		}
	}
	return nil
}
func ReadPoints(result Result) ([]string, error) {
	var v struct {
		Points []string `json:"points"`
		Error  string   `json:"error"`
	}
	if json.Unmarshal(result.Result, &v) != nil || v.Error != "" || v.Points == nil {
		return nil, ErrBinding
	}
	return v.Points, nil
}
````

### FILE: `microsoft_playwright_browser_gate/tests/stored-value-connected.spec.mjs`

```yaml
block_id: "GO-APPROVED-STORED-VALUE-TENDER:file21:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration/configuration/authorization/persistence/transport/UI glue; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "894d783942d8034f28d808ba58783c6da65332e6e9f17eb764ae3fb9d49fb14f"
variables: []
secrets_allowed: false
```

````text
import {test,expect} from '@playwright/test';
import {createHash} from 'node:crypto';
import {createRequire} from 'node:module';
import {resolve} from 'node:path';
import {pathToFileURL} from 'node:url';

test('source-backed gift issue, distinct review and redemption fund actual delivery without a provider',async({page,context},info)=>{
 if(process.env.ELITE_STORED_VALUE_BROWSER!=='1')throw new Error('explicit local handover fixture required');
 const base=process.env.ELITE_BASE_URL;expect(new URL(base).protocol).toBe('https:');expect(new URL(base).hostname).toBe('127.0.0.1');
 const identities=JSON.parse(process.env.ELITE_HANDOVER_IDENTITIES);
 const require=createRequire(resolve(process.env.ELITE_WEB_ROOT,'package.json'));
 const {EncryptJWT}=await import(pathToFileURL(require.resolve('jose')).href);
 async function identity(name){
  const jwt=await new EncryptJWT(identities[name]).setProtectedHeader({alg:'dir',enc:'A256GCM',typ:'JWT'}).setIssuedAt().setExpirationTime('300s').encrypt(createHash('sha256').update(process.env.AUTH_SESSION_SECRET).digest());
  await context.clearCookies();await context.addCookies([{name:'__Host-elite_session',value:jwt,url:base+'/',secure:true,httpOnly:true,sameSite:'Lax'}]);
 }
 page.setDefaultTimeout(12000);
 const errors=[];page.on('pageerror',error=>errors.push(error.message));
 const writes=[];page.on('request',request=>{if(request.method()==='POST'&&request.url().includes('/api/enterprise/'))writes.push(request.postDataJSON()?.action)});

 async function consult(id){await page.getByLabel('Referencia del pedido',{exact:true}).fill(id);await page.getByRole('button',{name:'Consultar pedido',exact:true}).click();await expect(page.getByText('Total del pedido',{exact:true})).toBeVisible()}
 async function propose(op,order,account){
  await identity('operator');await page.goto('/franchise/stored-value');await consult(order);
  await page.getByRole('combobox',{name:'Operación',exact:true}).selectOption(op);
  if(account)await page.getByLabel('Cuenta de gift card o fidelidad',{exact:true}).fill(account);
  await page.getByRole('button',{name:'Guardar solicitud para revisión',exact:true}).click();
  await expect(page.getByText('Solicitud guardada.',{exact:false})).toBeVisible();
  await expect(page.getByRole('button',{name:'Aprobar contenido revisado',exact:true})).toHaveCount(0);
 }
 async function approve(order){
  await identity('reviewer');await page.goto('/franchise/stored-value');await consult(order);
  await page.getByRole('button',{name:/^Revisar svapproval_/}).click();
  await page.getByLabel('Motivo de la decisión',{exact:true}).fill('Pedido, cuenta y valores exactos revisados en la composición de referencia');
  await page.getByRole('button',{name:'Aprobar contenido revisado',exact:true}).click();
  await expect(page.getByText('Decisión guardada.',{exact:false})).toBeVisible();
 }
 await propose('issue','source');await approve('source');
 const account=(await page.getByRole('heading',{name:/^Cuenta gift_/}).innerText()).replace('Cuenta ','');
 await propose('redeem','order',account);await approve('order');
 await identity('operator');await page.goto('/franchise/stored-value');await consult('order');
 await expect(page.getByRole('button',{name:'Registrar cobertura total para entrega',exact:true})).toBeEnabled();
 await page.route('**/api/enterprise/franchise/stored-value',async route=>{
  if(route.request().method()!=='POST'||route.request().postDataJSON()?.action!=='fund'){await route.continue();return}
  const response=await route.fetch();expect(response.status()).toBe(201);await route.abort('failed');
 });
 await page.getByRole('button',{name:'Registrar cobertura total para entrega',exact:true}).click();
 await expect(page.getByText('Resultado sin confirmar.',{exact:false})).toBeVisible();
 await page.unrouteAll();await page.reload();
 await page.getByRole('button',{name:'Consultar resultado pendiente',exact:true}).click();
 await expect(page.getByRole('heading',{name:'Cobertura total registrada',exact:true})).toBeVisible();
 expect(writes.filter(x=>x==='fund')).toHaveLength(1);
 await page.setViewportSize({width:390,height:844});expect(await page.evaluate(()=>document.documentElement.scrollWidth<=window.innerWidth+1)).toBe(true);
 await page.screenshot({path:info.outputPath('stored-value-recovered-mobile.png'),fullPage:true});await page.setViewportSize({width:1280,height:900});
 await page.getByRole('link',{name:'Continuar con la entrega',exact:true}).click();
 const panel=page.getByRole('article',{name:'Entrega del pedido order',exact:true});
 await expect(panel.getByRole('button',{name:'Preparar entrega',exact:true})).toBeEnabled();
 // Discard only the browser response after the actual Go/PG transaction.
 await page.route('**/api/enterprise/handovers',async route=>{
  const action=route.request().postDataJSON()?.action;
  if(route.request().method()!=='POST'||action!=='prepare'){await route.continue();return}
  const response=await route.fetch();expect(response.status()).toBe(201);await route.abort('failed');
 });
 await panel.getByRole('button',{name:'Preparar entrega',exact:true}).evaluate(button=>{button.click();button.click()});
 await expect(panel.getByText('El resultado quedó sin comprobar.',{exact:false})).toBeVisible();
 expect(writes.filter(x=>x==='prepare')).toHaveLength(1);
 await panel.getByRole('button',{name:'Consultar resultado de preparación',exact:true}).click();
 await expect(panel.getByText('Preparación recuperada:',{exact:false})).toBeVisible();
 await page.unrouteAll();await page.reload();
 const presentation=panel.getByRole('region',{name:'Presentación de entrega',exact:true});
 await expect(presentation.getByLabel('Entrega',{exact:true})).toHaveAttribute('readonly','');
 const handoverID=await presentation.getByLabel('Entrega',{exact:true}).inputValue();
 await expect(presentation.getByLabel('Versión actual de entrega')).toHaveValue('1');
 await presentation.getByLabel('ID del checklist',{exact:true}).fill('browser-checklist');
 await presentation.getByLabel('Versión del checklist',{exact:true}).fill('1');
 await presentation.getByLabel('Respuestas, una por línea').fill('serial=SERIAL-SYNTHETIC');
 await presentation.getByRole('button',{name:'Completar y presentar',exact:true}).click();
 await expect(presentation.getByText('Respuesta recibida.',{exact:false})).toBeVisible();
 await presentation.getByRole('button',{name:'Consultar presentación',exact:true}).click();
 await expect(presentation.getByText('Presentación recuperada:',{exact:false})).toBeVisible();
 await identity('stranger');await page.goto('/customer/handovers');
 await expect(page.getByRole('button',{name:'Registrar recepción',exact:true})).toHaveCount(0);
 const unauthorized=await context.request.post(base+'/api/enterprise/franchise/commands',{headers:{origin:base,'content-type':'application/json'},data:{action:'accept-handover',organizationId:'store',handoverId:handoverID,version:2,confirmedReceived:true,serialNumber:'SERIAL-SYNTHETIC',checklistId:'browser-checklist',checklistVersion:1}});
 expect([404,409]).toContain(unauthorized.status());
 await identity('customer');await page.goto('/customer/handovers');
 expect((await context.request.get(base+'/api/enterprise/handovers?kind=context&organizationId=store&orderId=order')).status()).toBe(403);
 await page.getByLabel('Número de serie observado').fill('SERIAL-SYNTHETIC');
 await page.getByRole('checkbox',{name:'Confirmo que recibí el activo identificado y revisé la preparación indicada'}).check();
 await page.getByRole('button',{name:'Registrar recepción',exact:true}).click();
 await expect(page.getByText('Aceptado:',{exact:false})).toBeVisible();
 await expect(page.getByRole('button',{name:'Registrar recepción',exact:true})).toHaveCount(0);
 await identity('operator');await page.goto('/franchise');
 await expect(panel.getByRole('button',{name:'Registrar cierre comercial',exact:true})).toBeEnabled();
 await page.route('**/api/enterprise/handovers',async route=>{
  if(route.request().method()!=='POST'||route.request().postDataJSON()?.action!=='release'){await route.continue();return}
  const response=await route.fetch();expect(response.status()).toBe(201);await route.abort('failed');
 });
 await panel.getByRole('button',{name:'Registrar cierre comercial',exact:true}).click();
 await expect(panel.getByText('El resultado quedó sin comprobar.',{exact:false})).toBeVisible();
 await page.unrouteAll();

 await panel.getByRole('button',{name:'Consultar resultado comercial',exact:true}).click();
 await expect(panel.getByText('Recibo histórico recuperado.',{exact:false})).toBeVisible();
 await panel.getByRole('button',{name:'Consultar vigencia del recibo',exact:true}).click();
 expect(writes.filter(x=>x==='propose')).toHaveLength(2);expect(writes.filter(x=>x==='decision')).toHaveLength(2);expect(writes.filter(x=>x==='prepare')).toHaveLength(1);expect(writes.filter(x=>x==='release')).toHaveLength(1);
 await identity('customer');expect((await context.request.get(base+'/api/enterprise/franchise/stored-value?kind=order&organization_id=store&order_id=order')).status()).toBe(403);
 await identity('foreign-org');expect((await context.request.get(base+'/api/enterprise/franchise/stored-value?kind=order&organization_id=store&order_id=order')).status()).toBe(403);
 expect(errors).toEqual([]);
 console.log('STORED_VALUE_BROWSER_PASS issue=1 redeem=1 distinct_decisions=2 funding=1 prepare=1 acceptance=1 release=1 recovered_lost_responses=true zero_provider=true');
});
````

### FILE: `src/app/api/enterprise/franchise/stored-value/route.test.ts`

```yaml
block_id: "GO-APPROVED-STORED-VALUE-TENDER:file22:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration/configuration/authorization/persistence/transport/UI glue; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "47e0e3506e980213ea619ed1de23ca5f1ec6b2e53ce4d3b23246c46098f38ca1"
variables: []
secrets_allowed: false
```

````typescript
import {beforeEach,expect,it,vi} from "vitest";
import {createHash} from "node:crypto";
const m=vi.hoisted(()=>({session:vi.fn(),get:vi.fn(),post:vi.fn()}));
vi.mock("@/platform/auth/session",()=>({readSession:m.session,allowed:(s:{permissions:string[]},p:string)=>s.permissions.includes(p)}));
vi.mock("@/platform/auth/oidc-client",()=>({applicationBaseUrl:()=>new URL("https://portal.example.test")}));
vi.mock("@/platform/backend/protected-client",()=>({protectedGet:m.get,protectedPost:m.post}));
import {GET,POST} from "./route";
const hash=(v:string)=>createHash("sha256").update(v).digest("hex");const key="a03c7850-13ad-4b38-81d0-120000000001";
const session={subject:"operator",tenantId:"tenant",organizations:["store"],permissions:["stored_value:read","stored_value:request","stored_value:approve","stored_value:fund"],accessToken:"never-browser-token"};
const proposal={action:"propose",organization_id:"store",operation_id:key,operation:"redeem",program_id:"gift",order_id:"order",expected_order_version:"9007199254740993",account_id:"account",original_operation_id:"",profile_sha256:"a".repeat(64)};
const payload='{"gross_minor_units":9007199254740993}';
const view=()=>({approval_id:"svapproval_"+hash(key),payload_sha256:hash(payload),payload_json:payload,receipt_json:"null",state:"pending",replay:false,order_id:"order",organization_id:"store",operation:"redeem",requester:"operator",expires_at:"2026-09-12T18:00:00Z",review:{operation_id:key,program_id:"gift",original_operation_id:"",currency:"ARS",gross_minor_units:"9007199254740993",gift_minor_units:"0",discount_minor_units:"0",entries:[{account_id:"account",points_delta:"-1.000000",applied_minor_units:"100"}]}});
const url="https://portal.example.test/api/enterprise/franchise/stored-value";
const post=(body:unknown=proposal)=>new Request(url,{method:"POST",headers:{origin:"https://portal.example.test","content-type":"application/json"},body:JSON.stringify(body)});
beforeEach(()=>{vi.clearAllMocks();m.session.mockResolvedValue(session);m.post.mockResolvedValue(view());m.get.mockResolvedValue(view())});
it("keeps int64 and hash-bound approval bytes exact",async()=>{const r=await POST(post());expect(r.status).toBe(201);const v=await r.json();expect(v.payload_json).toBe(payload);expect(v.review.gross_minor_units).toBe("9007199254740993");expect(JSON.stringify(m.post.mock.calls)).toContain('"expected_order_version":"9007199254740993"');expect(JSON.stringify(v)).not.toContain("never-browser-token")});
it.each([9007199254740992,"01","1e2","9223372036854775808"])("rejects invalid version %s",async value=>{expect((await POST(post({...proposal,expected_order_version:value}))).status).toBe(400);expect(m.post).not.toHaveBeenCalled()});
it.each(["amount_minor_units","requester","tenant_id"])("rejects forged field %s",async field=>{expect((await POST(post({...proposal,[field]:"forged"}))).status).toBe(400);expect(m.post).not.toHaveBeenCalled()});
it.each(["origin","permission","organization","unauthenticated"])("rejects %s before backend",async mode=>{let r=post();if(mode==="origin")r=new Request(r,{headers:{origin:"https://evil.test","content-type":"application/json"}});if(mode==="permission")m.session.mockResolvedValue({...session,permissions:[]});if(mode==="organization")r=post({...proposal,organization_id:"foreign"});if(mode==="unauthenticated")m.session.mockResolvedValue(null);expect([401,403]).toContain((await POST(r)).status);expect(m.post).not.toHaveBeenCalled()});
it.each(["digest","reference","organization"])("rejects mismatched %s in backend response",async mode=>{m.post.mockResolvedValue({...view(),...(mode==="digest"?{payload_json:payload+" "}:mode==="reference"?{approval_id:"other"}:{organization_id:"other"})});expect((await POST(post())).status).toBe(503);expect(m.post).toHaveBeenCalledTimes(1)});
it("recovers a proposal by its stable reference without a second mutation",async()=>{const r=await GET(new Request(url+"?kind=proposal-result&organization_id=store&operation_id="+key));expect(r.status).toBe(200);expect(m.get).toHaveBeenCalledWith(session,"/v1/franchise/stored-value/approvals/svapproval_"+hash(key),{organization_id:"store"});expect(m.post).not.toHaveBeenCalled()});
it("recovers exact full funding evidence with the original key",async()=>{const raw='{"amount":9007199254740993}';m.get.mockResolvedValue({funding_receipt_id:"funding",receipt_sha256:hash(raw),receipt_json:raw,organization_id:"store",order_id:"order",request_key:key});const r=await GET(new Request(url+"?kind=funding-result&organization_id=store&request_key="+key));expect(r.status).toBe(200);expect((await r.json()).receipt_json).toBe(raw);expect(m.get).toHaveBeenCalledWith(session,"/v1/franchise/stored-value/funding-result",{organization_id:"store"},key);expect(m.post).not.toHaveBeenCalled()});
````

### FILE: `src/app/api/enterprise/franchise/stored-value/route.ts`

```yaml
block_id: "GO-APPROVED-STORED-VALUE-TENDER:file23:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration/configuration/authorization/persistence/transport/UI glue; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "a320f3d5ac3509795ce5f088948297aca3a8f91af32985a32475b7c389a40872"
variables: []
secrets_allowed: false
```

````typescript
// AUTHORED scoped BFF. Exact JSON payload strings are hashed and forwarded as
// text; no amount is converted to Number or sent back as a business decision.
import { NextResponse } from "next/server";
import { createHash } from "node:crypto";
import { allowed,readSession } from "@/platform/auth/session";
import { applicationBaseUrl } from "@/platform/auth/oidc-client";
import { protectedGet,protectedPost } from "@/platform/backend/protected-client";
import { BackendProblem } from "@/platform/backend/public-client";
import { boundedCommandBody } from "@/platform/backend/bounded-command";
import { storedCommand,storedQuery,storedOrder,storedApproval,storedFunding } from "@/platform/stored-value/contracts";
const base="/v1/franchise/stored-value";
const reply=(value:unknown,status=200)=>NextResponse.json(value,{status,headers:{"cache-control":"no-store","x-content-type-options":"nosniff"}});
const digest=(s:string)=>createHash("sha256").update(s,"utf8").digest("hex");
function failure(e:unknown){return reply({code:"STORED_VALUE_UNCONFIRMED"},e instanceof BackendProblem&&[400,401,403,404,409,503].includes(e.status)?e.status:503)}
function approval(raw:unknown,org:string){const v=storedApproval.parse(raw);if(v.organization_id!==org||digest(v.payload_json)!==v.payload_sha256)throw new Error("exact approval mismatch");return v}
function funding(raw:unknown,org:string){const v=storedFunding.parse(raw);if(v.organization_id!==org||digest(v.receipt_json)!==v.receipt_sha256)throw new Error("exact receipt mismatch");return v}
export async function POST(request:Request){
 try{if(request.headers.get("origin")!==applicationBaseUrl().origin)return reply({code:"CROSS_ORIGIN_REJECTED"},403)}catch{return reply({code:"CROSS_ORIGIN_REJECTED"},403)}
 if(request.headers.get("content-type")!=="application/json")return reply({code:"UNSUPPORTED_MEDIA_TYPE"},415);
 const session=await readSession();if(!session)return reply({code:"UNAUTHENTICATED"},401);
 let parsed;try{parsed=storedCommand.safeParse(JSON.parse(await boundedCommandBody(request)))}catch{return reply({code:"INVALID_COMMAND"},400)}
 if(!parsed.success)return reply({code:"INVALID_COMMAND"},400);const c=parsed.data;
 const permission=c.action==="propose"?"stored_value:request":c.action==="decision"?"stored_value:approve":"stored_value:fund";
 if(!session.organizations.includes(c.organization_id)||!allowed(session,permission))return reply({code:"FORBIDDEN"},403);
 try{
  if(c.action==="propose"){
   const {action,...body}=c;void action;
   const v=approval(await protectedPost<unknown>(session,`${base}/proposals`,body),c.organization_id);
   if(v.order_id!==c.order_id||v.approval_id!=="svapproval_"+digest(c.operation_id))throw new Error("proposal reference mismatch");return reply(v,201);
  }
  if(c.action==="decision"){
   const v=approval(await protectedPost<unknown>(session,`${base}/approvals/${encodeURIComponent(c.approval_id)}/decision`,{organization_id:c.organization_id,payload_sha256:c.payload_sha256,approved:c.approved,reason:c.reason}),c.organization_id);
   if(v.approval_id!==c.approval_id||v.payload_sha256!==c.payload_sha256)throw new Error("decision reference mismatch");return reply(v);
  }
  const v=funding(await protectedPost<unknown>(session,`${base}/orders/${encodeURIComponent(c.order_id)}/funding`,{organization_id:c.organization_id,expected_order_version:c.expected_order_version},c.request_key),c.organization_id);
  if(v.order_id!==c.order_id||v.request_key!==c.request_key)throw new Error("funding reference mismatch");return reply(v,201);
 }catch(e){return failure(e)}
}
export async function GET(request:Request){
 if(["cross-site","none"].includes(request.headers.get("sec-fetch-site")??""))return reply({code:"CROSS_ORIGIN_REJECTED"},403);
 const session=await readSession();if(!session)return reply({code:"UNAUTHENTICATED"},401);if(!allowed(session,"stored_value:read"))return reply({code:"FORBIDDEN"},403);
 const search=new URL(request.url).searchParams;if(request.url.length>2048||[...search.keys()].some(k=>search.getAll(k).length!==1))return reply({code:"INVALID_QUERY"},400);
 const parsed=storedQuery.safeParse(Object.fromEntries(search));if(!parsed.success)return reply({code:"INVALID_QUERY"},400);const q=parsed.data;
 if(!session.organizations.includes(q.organization_id))return reply({code:"FORBIDDEN"},403);
 try{
  if(q.kind==="order"){
   const v=storedOrder.parse(await protectedGet<unknown>(session,`${base}/orders/${encodeURIComponent(q.order_id)}`,{organization_id:q.organization_id,after:q.after}));
   if(v.organization_id!==q.organization_id||v.order_id!==q.order_id)throw new Error("order scope mismatch");return reply(v);
  }
  if(q.kind==="funding-result"){
   const v=funding(await protectedGet<unknown>(session,`${base}/funding-result`,{organization_id:q.organization_id},q.request_key),q.organization_id);if(v.request_key!==q.request_key)throw new Error("funding reference mismatch");return reply(v);
  }
  const id=q.kind==="approval"?q.approval_id:"svapproval_"+digest(q.operation_id);
  const v=approval(await protectedGet<unknown>(session,`${base}/approvals/${encodeURIComponent(id)}`,{organization_id:q.organization_id}),q.organization_id);if(v.approval_id!==id)throw new Error("approval reference mismatch");return reply(v);
 }catch(e){return failure(e)}
}
````

### FILE: `src/app/franchise/stored-value/page.tsx`

```yaml
block_id: "GO-APPROVED-STORED-VALUE-TENDER:file24:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration/configuration/authorization/persistence/transport/UI glue; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "a9402448ba08c13074e68dc9364809dbcc4f7aa9943280da7dbd8481536051c6"
variables: []
secrets_allowed: false
```

````tsx
import {loadPrivateI18n} from "@/platform/i18n/load-private-i18n";
import { redirect } from "next/navigation";
import type { Route } from "next";
import { createHash } from "node:crypto";
import { z } from "zod";
import { readSession,allowed } from "@/platform/auth/session";
import { protectedGet } from "@/platform/backend/protected-client";
import { storedProfile,storedID } from "@/platform/stored-value/contracts";
import { StoredValueOperations } from "@/components/stored-value-operations";
const configuration=z.object({programs:z.array(z.object({calculation:z.object({id:storedID,kind:z.enum(["gift_card","loyalty"])})})).min(1).max(8)});
export default async function StoredValuePage({searchParams}:{searchParams:Promise<{organization?:string}>}){
 const {locale:privateLocale,t,controlled}=await loadPrivateI18n();

 const session=await readSession();if(!session)redirect("/api/auth/login?return_to=/franchise/stored-value" as Route);
 if(!allowed(session,"stored_value:read"))return <><h1>{t("p0088")}</h1><p>{t("p0089")}</p></>;
 const organization=(await searchParams).organization??session.organizations[0];
 if(!organization||!session.organizations.includes(organization))return <><h1>{t("p0088")}</h1><p>{t("p0090")}</p></>;
 const choose=<form action="/franchise/stored-value" method="get"><label>{t("p0091")}<select name="organization" defaultValue={organization}>{session.organizations.map(o=><option key={o}>{o}</option>)}</select></label><button className="button">{t("p0092")}</button></form>;
 try{
  const p=storedProfile.parse(await protectedGet<unknown>(session,"/v1/franchise/stored-value/profile",{organization_id:organization}));
  if(p.organization_id!==organization||createHash("sha256").update(p.profile_json,"utf8").digest("hex")!==p.profile_sha256)throw new Error("profile mismatch");
  const programs=configuration.parse(JSON.parse(p.profile_json)).programs.map(p=>p.calculation);
  return <><h1 className="pageTitle">{t("p0088")}</h1><p>{t("p0093")}</p>{choose}<StoredValueOperations organization={organization} profileHash={p.profile_sha256} programs={programs} permissions={session.permissions} subject={session.subject} scope={createHash("sha256").update(JSON.stringify([session.tenantId,session.subject,organization])).digest("hex")}/></>;
 }catch{return <><h1>{t("p0088")}</h1>{choose}<p role="alert">{t("p0094")}</p><a className="button" href={(`/franchise/stored-value?organization=${encodeURIComponent(organization)}`) as Route}>{t("p0095")}</a></>}
}
````

### FILE: `src/components/stored-value-operations.tsx`

```yaml
block_id: "GO-APPROVED-STORED-VALUE-TENDER:file25:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration/configuration/authorization/persistence/transport/UI glue; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "f58c81fb59930aafb071b8b936cbc5b291db6d0af15679618ab57d9475d2d4ca"
variables: []
secrets_allowed: false
```

````tsx
"use client";
import {usePrivateI18n} from "@/platform/i18n/private-provider";

// AUTHORED operator interaction and recovery glue. Browser submits references
// and selected operations; it never supplies an amount or approves itself.
import { useEffect,useState } from "react";
import { storedCommand,storedOrder,storedApproval,storedFunding,formatStoredMoney,formatStoredAdjustment,type StoredCommand,type StoredOrder,type StoredApproval,type StoredFunding } from "@/platform/stored-value/contracts";
type Program={id:string;kind:"gift_card"|"loyalty"};
export function StoredValueOperations({organization,profileHash,programs,permissions,subject,scope}:{organization:string;profileHash:string;programs:Program[];permissions:string[];subject:string;scope:string}){
 const {locale:privateLocale,t,controlled}=usePrivateI18n();

 const [orderID,setOrderID]=useState("");const [order,setOrder]=useState<StoredOrder|null>(null);const [view,setView]=useState<StoredApproval|null>(null);const [funding,setFunding]=useState<StoredFunding|null>(null);
 const [program,setProgram]=useState(programs[0]?.id??"");const [operation,setOperation]=useState<"issue"|"accrue"|"redeem"|"reverse">("redeem");const [account,setAccount]=useState("");const [original,setOriginal]=useState("");const [reason,setReason]=useState("");
 const [busy,setBusy]=useState(false);const [notice,setNotice]=useState("");const [pending,setPending]=useState<StoredCommand|null>(null);const [storageReady,setStorageReady]=useState(false);
 const storageKey=`stored-value-request:${scope}`;
 useEffect(()=>{try{const raw=localStorage.getItem(storageKey);if(raw!==null){const v=storedCommand.parse(JSON.parse(raw));if(v.organization_id!==organization)throw new Error("scope");setPending(v)}setStorageReady(true)}catch{setNotice(t("p0757"))}},[storageKey,organization]);
 const can=(p:string)=>permissions.includes(p);
 async function read(query:Record<string,string>){const r=await fetch("/api/enterprise/franchise/stored-value?"+new URLSearchParams({organization_id:organization,...query}),{cache:"no-store"});if(!r.ok)throw new Error("unconfirmed");return r.json() as Promise<unknown>}
 async function loadOrder(after?:string){if(busy)return;setBusy(true);setNotice("");setOrder(null);try{setOrder(storedOrder.parse(await read({kind:"order",order_id:orderID,...(after?{after}:{})})))}catch{setNotice(t("p0758"))}finally{setBusy(false)}}
 async function loadApproval(id:string){if(busy)return;setBusy(true);setNotice("");setView(null);try{setView(storedApproval.parse(await read({kind:"approval",approval_id:id})))}catch{setNotice(t("p0759"))}finally{setBusy(false)}}
 function keep(c:StoredCommand){const v=storedCommand.parse(c);const old=localStorage.getItem(storageKey);const raw=JSON.stringify(v);if(old!==null&&old!==raw)throw new Error("pending request");localStorage.setItem(storageKey,raw);if(localStorage.getItem(storageKey)!==raw)throw new Error("not retained");setPending(v)}
 function clear(){localStorage.removeItem(storageKey);setPending(null)}
 async function send(c:StoredCommand){
  if(busy||!storageReady)return;setBusy(true);setNotice("");
  try{keep(c);const r=await fetch("/api/enterprise/franchise/stored-value",{method:"POST",headers:{"content-type":"application/json"},body:JSON.stringify(c)});if(!r.ok)throw new Error("unconfirmed");const raw:unknown=await r.json();if(c.action==="fund")setFunding(storedFunding.parse(raw));else setView(storedApproval.parse(raw));clear();setOrder(null);setNotice(c.action==="propose"?t("p0760"):c.action==="fund"?t("p0761"):t("p0762"))}
  catch{setNotice(t("p0763"))}
  finally{setBusy(false)}
 }
 async function recover(){
  if(!pending||busy)return;setBusy(true);setNotice("");
  try{
   if(pending.action==="fund")setFunding(storedFunding.parse(await read({kind:"funding-result",request_key:pending.request_key})));
   else {const v=storedApproval.parse(await read(pending.action==="propose"?{kind:"proposal-result",operation_id:pending.operation_id}:{kind:"approval",approval_id:pending.approval_id}));setView(v);if(pending.action==="decision"&&v.state==="pending"){setNotice(t("p0764"));return}}
   clear();setOrder(null);setNotice(t("p0765"));
  }catch{setNotice(t("p0766"))}finally{setBusy(false)}
 }
 const blocked=busy||!storageReady||pending!==null;
 const submit=()=>{if(!order)return;void send({action:"propose",organization_id:organization,operation_id:crypto.randomUUID(),operation,program_id:program,order_id:order.order_id,expected_order_version:order.order_version,account_id:operation==="redeem"?account:"",original_operation_id:operation==="reverse"?original:"",profile_sha256:profileHash})};
 return <>
  <p role="status" aria-live="polite">{notice}</p>
  {pending&&<section className="card" aria-label={t("p0767")}><h2>{t("p0767")}</h2><p>{t("p0768")}</p><div style={{display:"flex",gap:12,flexWrap:"wrap"}}><button className="button" disabled={busy} onClick={()=>void recover()}>{t("p0168")}</button><button className="button" disabled={busy} onClick={()=>void send(pending)}>{t("p0769")}</button></div></section>}
  <section className="card"><h2>{t("p0770")}</h2><form onSubmit={e=>{e.preventDefault();void loadOrder()}}><label>{t("p0771")}<input name="order" value={orderID} onChange={e=>{setOrderID(e.target.value);setOrder(null);setView(null);setFunding(null)}} maxLength={128} required/></label><button className="button" disabled={busy}>{t("p0770")}</button></form>
   {order&&<><dl><dt>{t("p0772")}</dt><dd>{formatStoredMoney(order.gross_minor_units)}</dd><dt>{t("p0773")}</dt><dd>{formatStoredMoney(order.gift_minor_units)}</dd><dt>{t("p0774")}</dt><dd>{formatStoredMoney(order.discount_minor_units)}</dd><dt>{t("p0775")}</dt><dd>{formatStoredMoney(order.provider_due_minor_units)}</dd></dl>
    <h3>{t("p0776")}</h3>{order.approvals.ids.length===0?<p>{t("p0777")}</p>:<ul>{order.approvals.ids.map(id=><li key={id}><button className="button" disabled={busy} onClick={()=>void loadApproval(id)}>{t("p0778")} {id}</button></li>)}</ul>}
    {order.approvals.next_cursor&&<button className="button" disabled={busy} onClick={()=>void loadOrder(order.approvals.next_cursor)}>{t("p0779")}</button>}
    {can("stored_value:fund")&&order.provider_due_minor_units==="0"&&<button className="button" disabled={blocked} onClick={()=>void send({action:"fund",organization_id:organization,order_id:order.order_id,expected_order_version:order.order_version,request_key:crypto.randomUUID()})}>{t("p0780")}</button>}
   </>}
  </section>
  {can("stored_value:request")&&<section className="card"><h2>{t("p0781")}</h2><p>{t("p0782")}</p><form onSubmit={e=>{e.preventDefault();submit()}}>
   <label>{t("p0783")}<select value={program} onChange={e=>setProgram(e.target.value)}>{programs.map(p=><option key={p.id} value={p.id}>{p.kind==="gift_card"?t("p0784"):t("p0785")} {t("p0016")} {p.id}</option>)}</select></label>
   <label>{t("p0010")}<select value={operation} onChange={e=>setOperation(e.target.value as typeof operation)}><option value="redeem">{t("p0786")}</option><option value="issue">{t("p0787")}</option><option value="accrue">{t("p0788")}</option><option value="reverse">{t("p0789")}</option></select></label>
   {operation==="redeem"&&<label>{t("p0790")}<input value={account} onChange={e=>setAccount(e.target.value)} maxLength={128} required/></label>}
   {operation==="reverse"&&<label>{t("p0791")}<input value={original} onChange={e=>setOriginal(e.target.value)} maxLength={128} required/></label>}
   <button className="button" disabled={blocked||!order}>{t("p0792")}</button>
  </form></section>}
  {view&&<section className="card" style={{overflowWrap:"anywhere"}}><h2>{t("p0793")}</h2><p>{t("p0794")} {view.order_id}{t("p0795")} {view.requester}{t("p0796")} {view.state==="pending"?t("p0797"):view.state==="approved"?t("p0798"):t("p0337")}{t("p0008")}</p><p>{t("p0799")} {view.expires_at}</p><p>{t("p0800")}</p><dl><dt>{t("p0783")}</dt><dd>{view.review.program_id}</dd><dt>{t("p0010")}</dt><dd>{{issue:t("p0801"),accrue:t("p0802"),redeem:t("p0803"),reverse:t("p0804")}[view.operation]}</dd><dt>{t("p0805")}</dt><dd>{view.review.operation_id}</dd><dt>{t("p0806")}</dt><dd>{formatStoredMoney(view.review.gross_minor_units)}</dd>{view.review.original_operation_id&&<><dt>{t("p0807")}</dt><dd>{view.review.original_operation_id}</dd></>}</dl>
   {view.review.entries.map(entry=><article key={entry.account_id}><h3>{t("p0808")} {entry.account_id}</h3><p>{t("p0809")} {entry.points_delta}</p><p>{t("p0810")} {formatStoredAdjustment(entry.applied_minor_units)}</p></article>)}
   <details><summary>{t("p0811")}</summary><pre style={{whiteSpace:"pre-wrap",overflowWrap:"anywhere",maxHeight:400,overflow:"auto"}} aria-label={t("p0812")}>{view.payload_json}</pre></details>
   {can("stored_value:approve")&&view.state==="pending"&&view.requester!==subject&&<><label>{t("p0220")}<textarea value={reason} onChange={e=>setReason(e.target.value)} maxLength={2048} required/></label><div style={{display:"flex",gap:12,flexWrap:"wrap"}}>{[true,false].map(approved=><button className="button" key={String(approved)} disabled={blocked||!reason.trim()} onClick={()=>void send({action:"decision",organization_id:organization,approval_id:view.approval_id,payload_sha256:view.payload_sha256,approved,reason})}>{approved?t("p0813"):t("p0814")}</button>)}</div></>}
   {view.requester===subject&&view.state==="pending"&&<p>{t("p0815")}</p>}
   <button className="button" disabled={busy} onClick={()=>void loadApproval(view.approval_id)}>{t("p0816")}</button>
   {view.receipt_json!=="null"&&<details><summary>{t("p0817")}</summary><pre style={{whiteSpace:"pre-wrap",overflowWrap:"anywhere"}}>{view.receipt_json}</pre></details>}
  </section>}
  {funding&&<section className="card" style={{overflowWrap:"anywhere"}}><h2>{t("p0818")}</h2><p>{t("p0251")} {funding.order_id}{t("p0819")} {funding.funding_receipt_id}</p><a className="button" href="/franchise">{t("p0820")}</a><details><summary>{t("p0821")}</summary><pre style={{whiteSpace:"pre-wrap",overflowWrap:"anywhere"}}>{funding.receipt_json}</pre></details></section>}
 </>;
}
````

### FILE: `src/platform/stored-value/contracts.ts`

```yaml
block_id: "GO-APPROVED-STORED-VALUE-TENDER:file26:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration/configuration/authorization/persistence/transport/UI glue; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "88bb4e889e28de0aa00c80671270157a18a3af63d62a49c9105f100e338edb87"
variables: []
secrets_allowed: false
```

````typescript
// AUTHORED transport validation. Amounts remain decimal strings; the selected
// backend/source profile owns every points, discount and eligibility decision.
import { z } from "zod";
export const storedID=z.string().regex(/^[A-Za-z0-9][A-Za-z0-9_.:-]{0,127}$/);
export const storedHash=z.string().regex(/^[0-9a-f]{64}$/);
const integer=z.string().regex(/^(0|[1-9][0-9]{0,18})$/).refine(v=>BigInt(v)<=9223372036854775807n);
const version=integer.refine(v=>v!=="0");
const signedInteger=z.string().regex(/^-?(0|[1-9][0-9]{0,18})$/).refine(v=>BigInt(v)>=-9223372036854775808n&&BigInt(v)<=9223372036854775807n);
const document=z.string().min(1).max(65536);
export const storedProfile=z.object({organization_id:storedID,profile_sha256:storedHash,profile_json:document}).strict();
export const storedOrder=z.object({organization_id:storedID,order_id:storedID,order_version:version,currency:z.literal("ARS"),gross_minor_units:integer,gift_minor_units:integer,discount_minor_units:integer,provider_due_minor_units:integer,approvals:z.object({ids:z.array(storedID).max(50),next_cursor:storedID.optional()}).strict()}).strict().refine(v=>BigInt(v.gross_minor_units)===BigInt(v.gift_minor_units)+BigInt(v.discount_minor_units)+BigInt(v.provider_due_minor_units));
export const storedApproval=z.object({review:z.object({operation_id:storedID,program_id:storedID,original_operation_id:z.union([storedID,z.literal("")]),currency:z.literal("ARS"),gross_minor_units:integer,gift_minor_units:integer,discount_minor_units:integer,entries:z.array(z.object({account_id:storedID,points_delta:z.string().regex(/^-?(0|[1-9][0-9]*)\.[0-9]{6}$/),applied_minor_units:signedInteger}).strict()).min(1).max(40)}).strict(),approval_id:storedID,payload_sha256:storedHash,state:z.enum(["pending","approved","rejected"]),replay:z.boolean(),payload_json:document,receipt_json:document,order_id:storedID,organization_id:storedID,operation:z.enum(["issue","accrue","redeem","reverse"]),requester:z.string().min(1).max(128),expires_at:z.iso.datetime({offset:true})}).strict();
export const storedFunding=z.object({funding_receipt_id:storedID,receipt_sha256:storedHash,receipt_json:document,organization_id:storedID,order_id:storedID,request_key:storedID}).strict();
export type StoredOrder=z.infer<typeof storedOrder>;export type StoredApproval=z.infer<typeof storedApproval>;export type StoredFunding=z.infer<typeof storedFunding>;
const base={organization_id:storedID};
export const storedCommand=z.discriminatedUnion("action",[
 z.object({...base,action:z.literal("propose"),operation_id:z.uuid(),operation:z.enum(["issue","accrue","redeem","reverse"]),program_id:storedID,order_id:storedID,expected_order_version:version,account_id:z.union([storedID,z.literal("")]),original_operation_id:z.union([storedID,z.literal("")]),profile_sha256:storedHash}).strict(),
 z.object({...base,action:z.literal("decision"),approval_id:storedID,payload_sha256:storedHash,approved:z.boolean(),reason:z.string().trim().min(1).max(2048)}).strict(),
 z.object({...base,action:z.literal("fund"),order_id:storedID,expected_order_version:version,request_key:z.uuid()}).strict(),
]);
export type StoredCommand=z.infer<typeof storedCommand>;
export const storedQuery=z.discriminatedUnion("kind",[
 z.object({...base,kind:z.literal("order"),order_id:storedID,after:storedID.optional()}).strict(),
 z.object({...base,kind:z.literal("approval"),approval_id:storedID}).strict(),
 z.object({...base,kind:z.literal("proposal-result"),operation_id:z.uuid()}).strict(),
 z.object({...base,kind:z.literal("funding-result"),request_key:z.uuid()}).strict(),
]);
export function formatStoredMoney(value:string){const v=integer.parse(value);const cents=v.padStart(3,"0");return `${cents.slice(0,-2)},${cents.slice(-2)} ARS`;}

export function formatStoredAdjustment(value:string){const v=signedInteger.parse(value);return v.startsWith("-")?"−"+formatStoredMoney(v.slice(1)):formatStoredMoney(v)}
````

## 6. Configuration surface

See deploy/stored-value/profile.reference.json and docs/stored-value-runtime.md in
the connected pair. Enable only a hash-bound tenant/org/profile and exact Python
source lock. Empty credentials do not authorize a live payment. Source overrides
write a new tree/receipt and are not automatically admitted.

## 7. Dependency bill

Official Odoo19.0 commit99edb6dd82b7b560930c00b03b694ba700785370, selected source only,
LGPL-3.0-only. No Odoo application installation or added pip dependency. Existing
locked CPython3.14.4/Go1.26.8/PG18.6/Next16.3.4 runtime owners are retained. Detailed
source/Gitblob/SHA256 and runtime manifests are under docs/provenance/ODOO_*.
Complete modified module source and exact LICENSE/COPYRIGHT ship with the Python
pack; IPC alone does not determine the legal boundary.

## 8. Apply order

Compose into an absent/empty destination using the canonical full profile and
MARKDOWN-COMPOSITOR0.3.0. Apply selected migrations in numeric order; shared0066,
stored-value0067 and handover0068 retain distinct owners. Never replace existing
history. Down migrations refuse populated immutable history; explicit dependency
drops avoid CASCADE. Enable the optional host module only after exact profile/source
validation. The updater preserves source LICENSE without adding a final LF.

## 9. Verification

Source oracle/expectation comparisons, loader/source-replacement regressions and
Go fuzz receipts are recorded in STORED_VALUE_SOURCE_CLOSURE_V402.md. Real PG/HTTP
partial gift and loyalty SDK payments, fully funded delivery, separate reviewer,
commit expiry rollback, concurrency, body/permission negatives, mobile browser and
lost-response recovery are in STORED_VALUE_TENDER_V402.md and
STORED_VALUE_OPERATOR_FLOW_V402.md. Changed-source checks must rerun their affected
gates. These receipts do not claim the full Odoo suite, SAST/DAST, target load or
production acceptance. Full composition SCA remains the separate T2803 owner.

## 10. Reconstruction evidence

Five profiles reconstruct exactly from source packs: full franchise84/1112,
backend40/567, web8/163, serverless57/793 and HTTP metrics85/1120. The two narrower
Go compositions compile changed owners and test packages with zero tests executed.
Runtime receipts from the identical candidate source are retained without replaying
unchanged payment/browser/domain suites. Canonical admission evidence is
reconstruction_evidence/STORED_VALUE_CONNECTED_RELEASE_V402.md and its JSON.

CONDITIONED selection is limited to the declared assisted reference program,
exact source/runtime/profile hashes, corresponding LGPL source/notices and the
existing scoped approval/payment/handover owners. The module's local functionality
is proven. Whole-composition SCA/T2803, operations/T2809 and target business/live
acceptance remain their explicit owners; these are not asserted by this pack.
The metadata-only change to REBUILD_VERIFIED preserves every materialized byte;
a final reference reconstruction binds the admitted pack bytes before publication.

V402 composed delta: T2804 role metrics use existing domain read models with exact strings, organization/customer/factory/program permissions, NPS minimum/retention and no zero on unavailable. FAIL868 converted lead and FAIL457 bounded generic body corrected. ROLE_METRICS_RELEASE_V402.md/json; no new dependency or corporate authorship.

### FILE: `internal/platform/httpapi/stored_value_metrics.go`

```yaml
block_id: "GO-APPROVED-STORED-VALUE-TENDER-METRIC-DELTA:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "01ea702292998669fea8d729ddc0ab5810e77edbf87e26f47f3ea8dc53d33aad"
variables: []
secrets_allowed: false
```

````go
package httpapi

import (
	"context"
	"elite.local/enterprise/internal/platform/identity"
	sv "elite.local/enterprise/internal/storedvaluebridge"
	"net/http"
	"net/url"
	"time"
)

func registerStoredValueMetrics(m StoredValueModule, mux *http.ServeMux, protected func(http.ResponseWriter, *http.Request, string, string) (identity.Principal, bool)) {
	mux.HandleFunc("GET /v1/franchise/stored-value/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "private, no-store")
		w.Header().Set("Vary", "Authorization")
		q, e := url.ParseQuery(r.URL.RawQuery)
		if e != nil || len(r.URL.RawQuery) > 512 || len(q) != 2 || len(q["organization_id"]) != 1 || len(q["program_id"]) != 1 || !sv.ValidID(q.Get("organization_id")) || !sv.ValidID(q.Get("program_id")) {
			writeProblem(w, 400, "METRIC_INVALID", "one organization and selected program required")
			return
		}
		p, ok := protected(w, r, "stored_value:read", q.Get("organization_id"))
		if !ok {
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 4*time.Second)
		defer cancel()
		v, e := m.Service.Metrics(ctx, p, q.Get("organization_id"), q.Get("program_id"))
		if storedValueError(w, e) {
			return
		}
		writeJSON(w, 200, v)
	})
}
````

### FILE: `internal/platform/postgres/stored_value_metrics.go`

```yaml
block_id: "GO-APPROVED-STORED-VALUE-TENDER-METRIC-DELTA:file2:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "6422c65c90aea5b29d24e0b2c6bdeaac7ce949fe5069b85076e212cd0316be65"
variables: []
secrets_allowed: false
```

````go
package postgres

// AUTHORED reporting projection. Existing profile/auth/immutable ledgers govern data.
import (
	"context"
	"elite.local/enterprise/internal/platform/identity"
	sv "elite.local/enterprise/internal/storedvaluebridge"
	"github.com/jackc/pgx/v5"
	"time"
)

type StoredValueMetricRow struct {
	Operation    string `json:"operation"`
	Operations   string `json:"operations"`
	PointsDelta  string `json:"points_delta"`
	AppliedMinor string `json:"applied_minor_units"`
}
type StoredValueMetrics struct {
	OrganizationID string                 `json:"organization_id"`
	ProgramID      string                 `json:"program_id"`
	ProgramKind    string                 `json:"program_kind"`
	Currency       string                 `json:"currency"`
	ProfileSHA256  string                 `json:"profile_sha256"`
	Source         string                 `json:"source"`
	Basis          string                 `json:"basis"`
	ObservedAt     time.Time              `json:"observed_at"`
	Rows           []StoredValueMetricRow `json:"rows"`
}

func (s *StoredValue) Metrics(ctx context.Context, p identity.Principal, org, programID string) (StoredValueMetrics, error) {
	var empty StoredValueMetrics
	if !s.allowed(p, "stored_value:read") || !sv.ValidID(programID) {
		return empty, sv.ErrBinding
	}
	tenant, bound := s.profile.Scope()
	if org != bound {
		return empty, sv.ErrBinding
	}
	program, e := s.profile.Program(programID)
	if e != nil {
		return empty, e
	}
	tx, e := s.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead, AccessMode: pgx.ReadOnly})
	if e != nil {
		return empty, e
	}
	defer tx.Rollback(ctx)
	if _, e = tx.Exec(ctx, "set local statement_timeout='3s'"); e != nil {
		return empty, e
	}
	out := StoredValueMetrics{OrganizationID: org, ProgramID: programID, ProgramKind: program.Calculation.Kind, Currency: program.Calculation.Currency, ProfileSHA256: s.profile.SHA256(), Source: "stored_value.operation + entry + account", Basis: "committed_ledger_by_operation", Rows: []StoredValueMetricRow{}}
	if e = tx.QueryRow(ctx, "select clock_timestamp()").Scan(&out.ObservedAt); e != nil {
		return empty, e
	}
	rows, e := tx.Query(ctx, `select op.operation,count(distinct op.operation_id)::text,sum(e.points_delta)::text,sum(e.applied_minor_units::numeric)::text
 from stored_value.operation op join stored_value.entry e on e.tenant_id=op.tenant_id and e.operation_id=op.operation_id
 join stored_value.account a on a.tenant_id=e.tenant_id and a.account_id=e.account_id
 where op.tenant_id=$1 and op.organization_id=$2 and op.program_id=$3 and a.organization_id=$2 and a.program_id=$3 and a.currency=$4 and a.kind=$5
 group by op.operation order by op.operation`, tenant, org, programID, out.Currency, out.ProgramKind)
	if e != nil {
		return empty, e
	}
	for rows.Next() {
		var v StoredValueMetricRow
		if e = rows.Scan(&v.Operation, &v.Operations, &v.PointsDelta, &v.AppliedMinor); e != nil {
			rows.Close()
			return empty, e
		}
		out.Rows = append(out.Rows, v)
	}
	e = rows.Err()
	rows.Close()
	if e != nil {
		return empty, e
	}
	return out, tx.Commit(ctx)
}
````


V402 composed delta: T2804 private es/en display, per-user/tenant preference and browser negotiation;1198messages,21exact guides,5hash-bound curricula; original commands/content/policies preserved. PRIVATE_LOCALE_RELEASE_V402.md/json. AUTHORED glue; no new dependency or corporate attribution.
