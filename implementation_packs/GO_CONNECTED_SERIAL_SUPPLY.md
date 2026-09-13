# Connected serialized supply reference

## 1. Metadata

```yaml
pack_id: "GO-CONNECTED-SERIAL-SUPPLY"
pack_version: "0.2.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: RECONSTRUCTIBLE
  admission: CONDITIONED
claim: "Connected demand/PO quantities, factory serial quality, split ASN shipments, receiving, quarantine and distinct-human stock release using existing owners; local reference."
stacks: ["Go 1.26.8", "PostgreSQL 18.6", "Node.js 24.20.0 where frontend selected"]
compatible_with: ["MARKDOWN-COMPOSITOR0.3.0", "V402 complete franchise closure"]
incompatible_with: ["unbound tenant/provider/account", "production certification inferred from fixtures", "implicit source or business-policy attribution"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: []
verified_at: "2026-09-12"
```

## 2. Applicability

Explicit J2 of ENTERPRISE_FULL_STACK_BLUEPRINT. Select complete reference closure. No implicit pricing, legal quality standard, overreceipt, unscoped actor or stock ledger.

## 3. Architecture contract

Bound immutable plan and versioned commands compose existing Operations SQL transaction writers, shared human approval and stock. One purchase lock; actor/hash replay. ASN membership and rejected replacement quantities are enforced. Deferred physical/approval guards preserve immutable history and reject bypass. Receipt quarantine excludes ATP until review. Existing post-release commerce/service remains authoritative.

## 4. Exact file manifest

```text
CREATE cmd/electromobility-api/serial_supply.go
CREATE cmd/electromobility-api/serial_supply_test.go
CREATE db/migrations/0072_connected_serial_supply.down.sql
CREATE db/migrations/0072_connected_serial_supply.up.sql
CREATE deploy/serial-supply/plan.reference.json
CREATE docs/SERIAL_SUPPLY_REFERENCE.md
CREATE internal/approval/serial_quality_test.go
CREATE internal/platform/httpapi/serial_supply.go
CREATE internal/platform/httpapi/serial_supply_test.go
CREATE internal/platform/postgres/serial_supply.go
CREATE internal/platform/postgres/serial_supply_commands.go
CREATE internal/platform/postgres/serial_supply_connected_integration_test.go
CREATE internal/platform/postgres/serial_supply_http_integration_test.go
CREATE internal/platform/postgres/serial_supply_owner_tx.go
CREATE internal/platform/postgres/serial_supply_quality.go
CREATE internal/platform/postgres/serial_supply_shipping.go
CREATE internal/serialsupply/contract.go
CREATE internal/serialsupply/contract_fuzz_test.go
```

## 5. Materialization blocks

### FILE: `cmd/electromobility-api/serial_supply.go`

```yaml
block_id: "GO-CONNECTED-SERIAL-SUPPLY:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "19a931e0968ba78dabb3840f03148a325315b2fc0b814ccc2668d01bb32a77dc"
variables: []
secrets_allowed: false
```

````go
package main

// AUTHORED opt-in activation of the explicit serialized reference policy.
import (
	"context"
	"crypto/sha256"
	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/platform/postgres"
	sc "elite.local/enterprise/internal/serialsupply"
	"encoding/hex"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
)

var errSerialSupplyConfiguration = errors.New("serial supply activation is invalid")

func init() { serialSupplyModuleFactory = selectedSerialSupplyModule }
func selectedSerialSupplyModule(ctx context.Context, pool *pgxpool.Pool, lookup func(string) string) (httpapi.EnterpriseModule, error) {
	if lookup == nil {
		return nil, errSerialSupplyConfiguration
	}
	switch lookup("SERIAL_SUPPLY_ENABLED") {
	case "", "false":
		return nil, nil
	case "true":
	default:
		return nil, errSerialSupplyConfiguration
	}
	hash := sha256.Sum256([]byte(sc.PolicyJSON))
	if pool == nil || lookup("SERIAL_SUPPLY_POLICY_SHA256") != hex.EncodeToString(hash[:]) {
		return nil, errSerialSupplyConfiguration
	}
	var ready bool
	err := pool.QueryRow(ctx, `select
 (select count(*) from pg_trigger where not tgisinternal and tgenabled in ('O','A') and (
  tgrelid=to_regclass('procurement.purchase_order') and tgname='serial_supply_purchase_guard'
  or tgrelid=to_regclass('factory.production_unit') and tgname='serial_supply_factory_guard'
  or tgrelid=to_regclass('inventory.stock_unit') and tgname='serial_supply_stock_guard'
  or tgrelid=to_regclass('approval.request') and tgname='serial_supply_approval_guard'
  or tgrelid=to_regclass('procurement.serial_supply_plan') and tgname in ('serial_supply_plan_guard','serial_supply_plan_immutable')
  or tgrelid=to_regclass('procurement.serial_supply_line') and tgname in ('serial_supply_line_guard','serial_supply_line_immutable')
  or tgrelid=to_regclass('procurement.serial_supply_step') and tgname='serial_supply_step_immutable'
  or tgrelid=to_regclass('procurement.serial_supply_unit') and tgname='serial_supply_unit_immutable'
  or tgrelid=to_regclass('procurement.serial_supply_shipment') and tgname='serial_supply_shipment_immutable'
  or tgrelid=to_regclass('procurement.serial_supply_manifest') and tgname='serial_supply_manifest_immutable'
  or tgrelid=to_regclass('procurement.serial_supply_receipt') and tgname='serial_supply_receipt_immutable'
  or tgrelid=to_regclass('procurement.serial_supply_quality') and tgname='serial_supply_quality_immutable'
  or tgrelid=to_regclass('procurement.serial_supply_effect') and tgname='serial_supply_effect_immutable'
 ))=15 and exists(select 1 from pg_constraint where conrelid=to_regclass('approval.request') and conname='request_kind_check'
 and pg_get_constraintdef(oid) like '%serial_quality%')`).Scan(&ready)
	if err != nil || !ready {
		return nil, errSerialSupplyConfiguration
	}
	store, err := postgres.NewSerialSupply(pool)
	if err != nil {
		return nil, errSerialSupplyConfiguration
	}
	return httpapi.SerialSupplyModule{Service: store}, nil
}
````

### FILE: `cmd/electromobility-api/serial_supply_test.go`

```yaml
block_id: "GO-CONNECTED-SERIAL-SUPPLY:file2:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "055dd15e32cc1d85e3ef7f6d1271fd438ec368fc56dda17d259cdbdc289e16a2"
variables: []
secrets_allowed: false
```

````go
package main

// AUTHORED activation and active-guard boundary; no credential access.
import (
	"context"
	"crypto/sha256"
	sc "elite.local/enterprise/internal/serialsupply"
	"encoding/hex"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"testing"
)

func TestSerialSupplyHostDisabledDoesNotReadPolicy(t *testing.T) {
	calls := 0
	module, err := selectedSerialSupplyModule(context.Background(), nil, func(key string) string {
		calls++
		if key != "SERIAL_SUPPLY_ENABLED" {
			t.Fatal("disabled module read", key)
		}
		return "false"
	})
	if err != nil || module != nil || calls != 1 {
		t.Fatal(module, err, calls)
	}
	if _, err = selectedSerialSupplyModule(context.Background(), nil, func(string) string { return "TRUE" }); err == nil {
		t.Fatal("invalid flag enabled")
	}
}
func TestSerialSupplyHostRequiresPolicyAndGuards(t *testing.T) {
	url := os.Getenv("PAYMENT_CONNECTED_DB_URL")
	if url == "" {
		t.Skip("owned fixture DB not configured")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	hash := sha256.Sum256([]byte(sc.PolicyJSON))
	want := hex.EncodeToString(hash[:])
	config := map[string]string{"SERIAL_SUPPLY_ENABLED": "true", "SERIAL_SUPPLY_POLICY_SHA256": want}
	lookup := func(key string) string { return config[key] }
	if _, err = selectedSerialSupplyModule(ctx, nil, lookup); err == nil {
		t.Fatal("missing pool admitted")
	}
	config["SERIAL_SUPPLY_POLICY_SHA256"] = ""
	if _, err = selectedSerialSupplyModule(ctx, pool, lookup); err == nil {
		t.Fatal("missing policy hash admitted")
	}
	config["SERIAL_SUPPLY_POLICY_SHA256"] = want
	module, err := selectedSerialSupplyModule(ctx, pool, lookup)
	if err != nil || module == nil {
		t.Fatal("valid schema/profile rejected", err)
	}
	if _, err = pool.Exec(ctx, `alter table approval.request disable trigger serial_supply_approval_guard`); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if _, err := pool.Exec(context.Background(), `alter table approval.request enable trigger serial_supply_approval_guard`); err != nil {
			t.Error(err)
		}
	}()
	if _, err = selectedSerialSupplyModule(ctx, pool, lookup); err == nil {
		t.Fatal("disabled approval binding guard admitted")
	}
	t.Log("SERIAL_SUPPLY_HOST_PASS exact policy hash, pool and15active guards required; disabled mode reads no policy")
}
````

### FILE: `db/migrations/0072_connected_serial_supply.down.sql`

```yaml
block_id: "GO-CONNECTED-SERIAL-SUPPLY:file3:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "3b41494faff932f51879a7a38dee54336d947932b35f45e06dd0aff01ac6ef4b"
variables: []
secrets_allowed: false
```

````sql
begin;
do $$ begin
 if exists(select 1 from procurement.serial_supply_plan) or exists(select 1 from approval.request where kind='serial_quality')
 then raise exception 'serial supply immutable history exists; downgrade refused';end if;
end $$;
drop trigger serial_supply_approval_guard on approval.request;
drop function procurement.serial_supply_guard_approval();
drop trigger serial_supply_stock_guard on inventory.stock_unit;
drop trigger serial_supply_factory_guard on factory.production_unit;
drop trigger serial_supply_purchase_guard on procurement.purchase_order;
drop function procurement.serial_supply_guard_stock();
drop function procurement.serial_supply_guard_factory();
drop function procurement.serial_supply_guard_purchase();
drop table procurement.serial_supply_effect;
drop table procurement.serial_supply_quality;
drop table procurement.serial_supply_receipt;
drop table procurement.serial_supply_manifest;
drop table procurement.serial_supply_shipment;
drop table procurement.serial_supply_unit;
drop table procurement.serial_supply_step;
drop table procurement.serial_supply_line;
drop table procurement.serial_supply_plan;
drop function procurement.serial_supply_guard_plan();
drop function procurement.serial_supply_fact_immutable();
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

commit;
````

### FILE: `db/migrations/0072_connected_serial_supply.up.sql`

```yaml
block_id: "GO-CONNECTED-SERIAL-SUPPLY:file4:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "809aafb6ae1ec9968cb7a200a4dd667a1157132605f436c5728c97353f6fe999"
variables: []
secrets_allowed: false
```

````sql
begin;
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

-- AUTHORED physical bindings around the existing purchase/factory/stock owners.
-- The selected reference policy admits zero overreceipt and separate human QC.
-- This is orchestration, not copied Business Central or a new stock/cost ledger.
create table procurement.serial_supply_plan (
 tenant_id uuid not null, purchase_order_id text not null,
 destination_organization_id text not null, factory_organization_id text not null,
 supplier_id text not null, demand_reference text not null check(length(demand_reference) between 1 and 1024),
 policy_code text not null check(policy_code='strict-serial-reference/v1'),
 policy_sha256 text not null check(policy_sha256 ~ '^[0-9a-f]{64}$'),
 created_by text not null check(length(created_by) between 1 and 128),
 created_at timestamptz not null default clock_timestamp(),
 primary key(tenant_id,purchase_order_id),
 foreign key(tenant_id,purchase_order_id) references procurement.purchase_order(tenant_id,purchase_order_id),
 foreign key(tenant_id,destination_organization_id) references org.organization(tenant_id,organization_id),
 foreign key(tenant_id,factory_organization_id) references org.organization(tenant_id,organization_id),
 foreign key(tenant_id,supplier_id) references partner.supplier(tenant_id,supplier_id)
);
create table procurement.serial_supply_line (
 tenant_id uuid not null,purchase_order_id text not null,line_id text not null,
 variant_id text not null,quantity integer not null check(quantity between 1 and 1000),
 primary key(tenant_id,purchase_order_id,line_id),
 unique(tenant_id,purchase_order_id,variant_id),
 foreign key(tenant_id,purchase_order_id) references procurement.serial_supply_plan(tenant_id,purchase_order_id),
 foreign key(tenant_id,variant_id) references catalog.vehicle_variant(tenant_id,variant_id)
);
create table procurement.serial_supply_step (
 tenant_id uuid not null,purchase_order_id text not null,command_id text not null,
 version bigint not null check(version>0),
 kind text not null check(kind in ('planned','submit','confirm','start','cancel','register','milestone','ship','receive','quality','quality-reject','reinspect')),
 actor text not null check(length(actor) between 1 and 128),
 request_sha256 text not null check(request_sha256 ~ '^[0-9a-f]{64}$'),
 evidence_sha256 text not null check(evidence_sha256 ~ '^[0-9a-f]{64}$'),
 payload_sha256 text not null check(payload_sha256 ~ '^[0-9a-f]{64}$'),
 payload jsonb not null check(jsonb_typeof(payload)='object' and pg_column_size(payload)<=131072),
 recorded_at timestamptz not null default clock_timestamp(),
 transaction_id xid8 not null default pg_current_xact_id(),
 primary key(tenant_id,purchase_order_id,command_id),
 unique(tenant_id,purchase_order_id,version),
 foreign key(tenant_id,purchase_order_id) references procurement.serial_supply_plan(tenant_id,purchase_order_id)
);
create table procurement.serial_supply_unit (
 tenant_id uuid not null,purchase_order_id text not null,line_id text not null,
 production_unit_id text not null,command_id text not null,
 primary key(tenant_id,production_unit_id),
 unique(tenant_id,purchase_order_id,production_unit_id),
 foreign key(tenant_id,purchase_order_id,line_id) references procurement.serial_supply_line(tenant_id,purchase_order_id,line_id),
 foreign key(tenant_id,production_unit_id) references factory.production_unit(tenant_id,production_unit_id),
 foreign key(tenant_id,purchase_order_id,command_id) references procurement.serial_supply_step(tenant_id,purchase_order_id,command_id) deferrable initially deferred
);
create table procurement.serial_supply_shipment (
 tenant_id uuid not null,purchase_order_id text not null,shipment_id text not null,command_id text not null,
 primary key(tenant_id,purchase_order_id,shipment_id),
 foreign key(tenant_id,purchase_order_id) references procurement.serial_supply_plan(tenant_id,purchase_order_id),
 foreign key(tenant_id,purchase_order_id,command_id) references procurement.serial_supply_step(tenant_id,purchase_order_id,command_id) deferrable initially deferred
);
create table procurement.serial_supply_manifest (
 tenant_id uuid not null,purchase_order_id text not null,shipment_id text not null,
 production_unit_id text not null,stock_unit_id text not null,
 primary key(tenant_id,production_unit_id),unique(tenant_id,stock_unit_id),
 unique(tenant_id,purchase_order_id,shipment_id,production_unit_id),
 foreign key(tenant_id,purchase_order_id,shipment_id) references procurement.serial_supply_shipment(tenant_id,purchase_order_id,shipment_id),
 foreign key(tenant_id,purchase_order_id,production_unit_id) references procurement.serial_supply_unit(tenant_id,purchase_order_id,production_unit_id),
 foreign key(tenant_id,stock_unit_id) references inventory.stock_unit(tenant_id,stock_unit_id)
);
create table procurement.serial_supply_receipt (
 tenant_id uuid not null,purchase_order_id text not null,shipment_id text not null,
 production_unit_id text not null,command_id text not null,
 primary key(tenant_id,production_unit_id),
 foreign key(tenant_id,purchase_order_id,shipment_id,production_unit_id) references procurement.serial_supply_manifest(tenant_id,purchase_order_id,shipment_id,production_unit_id),
 foreign key(tenant_id,purchase_order_id,command_id) references procurement.serial_supply_step(tenant_id,purchase_order_id,command_id) deferrable initially deferred
);
create table procurement.serial_supply_quality (
 tenant_id uuid not null,purchase_order_id text not null,production_unit_id text not null,
 stage text not null check(stage in ('factory','receipt')),attempt integer not null check(attempt>0),
 approval_id text not null,command_id text not null,source_state text not null,
 source_version bigint not null check(source_version>=0),
 payload_sha256 text not null check(payload_sha256 ~ '^[0-9a-f]{64}$'),
 primary key(tenant_id,production_unit_id,stage,attempt),unique(tenant_id,approval_id),
 foreign key(tenant_id,purchase_order_id,production_unit_id) references procurement.serial_supply_unit(tenant_id,purchase_order_id,production_unit_id),
 foreign key(tenant_id,approval_id) references approval.request(tenant_id,request_id),
 foreign key(tenant_id,purchase_order_id,command_id) references procurement.serial_supply_step(tenant_id,purchase_order_id,command_id) deferrable initially deferred
);
create table procurement.serial_supply_effect (
 tenant_id uuid not null,purchase_order_id text not null,command_id text not null,
 entity text not null check(entity in ('purchase','factory','stock')),entity_id text not null,
 from_state text not null,to_state text not null,from_version bigint not null,to_version bigint not null,
 transaction_id xid8 not null default pg_current_xact_id(),
 primary key(tenant_id,purchase_order_id,command_id,entity,entity_id,to_state),
 foreign key(tenant_id,purchase_order_id,command_id) references procurement.serial_supply_step(tenant_id,purchase_order_id,command_id) deferrable initially deferred
);
create function procurement.serial_supply_fact_immutable() returns trigger language plpgsql as $$
begin raise exception 'serial supply evidence is immutable';end $$;
create function procurement.serial_supply_guard_purchase() returns trigger language plpgsql as $$
begin
 if not exists(select 1 from procurement.serial_supply_plan where tenant_id=old.tenant_id and purchase_order_id=old.purchase_order_id) then return null;end if;
 if tg_op='DELETE' then raise exception 'bound supply order cannot be deleted';end if;
 if row(new.tenant_id,new.purchase_order_id,new.supplier_id,new.destination_organization_id,new.currency,new.total_minor_units,new.created_at)
 is distinct from row(old.tenant_id,old.purchase_order_id,old.supplier_id,old.destination_organization_id,old.currency,old.total_minor_units,old.created_at)
 then raise exception 'bound supply order terms are immutable';end if;
 if row(new.state,new.version) is distinct from row(old.state,old.version) and not exists(
  select 1 from procurement.serial_supply_effect e where e.tenant_id=new.tenant_id and e.purchase_order_id=new.purchase_order_id
  and e.entity='purchase' and e.entity_id=new.purchase_order_id and e.from_state=old.state and e.to_state=new.state
  and e.from_version=old.version and e.to_version=new.version and e.transaction_id=pg_current_xact_id())
 then raise exception 'bound supply order requires connected command';end if;
 return null;
end $$;
create constraint trigger serial_supply_purchase_guard after update or delete on procurement.purchase_order
 deferrable initially deferred for each row execute function procurement.serial_supply_guard_purchase();
create function procurement.serial_supply_guard_factory() returns trigger language plpgsql as $$
declare po text;ten uuid;line text;expected_variant text;maximum integer;actual bigint;
begin
 if tg_op='DELETE' then po:=old.purchase_order_id;ten:=old.tenant_id;else po:=new.purchase_order_id;ten:=new.tenant_id;end if;
 if not exists(select 1 from procurement.serial_supply_plan where tenant_id=ten and purchase_order_id=po) then
  if tg_op='UPDATE' and exists(select 1 from procurement.serial_supply_unit where tenant_id=old.tenant_id and production_unit_id=old.production_unit_id)
   then raise exception 'bound factory identity is immutable';end if;return null;
 end if;
 if tg_op='DELETE' then raise exception 'bound factory unit cannot be deleted';end if;
 select l.line_id,l.variant_id,l.quantity into line,expected_variant,maximum
 from procurement.serial_supply_unit b join procurement.serial_supply_line l using(tenant_id,purchase_order_id,line_id)
 where b.tenant_id=ten and b.purchase_order_id=po and b.production_unit_id=new.production_unit_id;
 if not found or expected_variant<>new.variant_id then raise exception 'factory unit requires matching supply line';end if;
 if tg_op='UPDATE' and row(new.tenant_id,new.production_unit_id,new.purchase_order_id,new.variant_id,new.serial_number,new.vin,new.battery_serial_number)
 is distinct from row(old.tenant_id,old.production_unit_id,old.purchase_order_id,old.variant_id,old.serial_number,old.vin,old.battery_serial_number)
 then raise exception 'bound factory identity is immutable';end if;
 if (tg_op='INSERT' or new.state is distinct from old.state) and not exists(
  select 1 from procurement.serial_supply_effect e where e.tenant_id=ten and e.purchase_order_id=po and e.entity='factory'
  and e.entity_id=new.production_unit_id and e.from_state=case when tg_op='INSERT' then '' else old.state end
  and e.to_state=new.state and e.transaction_id=pg_current_xact_id())
 then raise exception 'bound factory state requires connected command';end if;
 select count(*) into actual from procurement.serial_supply_unit b join factory.production_unit u using(tenant_id,production_unit_id)
 where b.tenant_id=ten and b.purchase_order_id=po and b.line_id=line and u.state<>'rejected';
 if actual>maximum then raise exception 'serial supply quantity exceeded';end if;
 return null;
end $$;
create constraint trigger serial_supply_factory_guard after insert or update or delete on factory.production_unit
 deferrable initially deferred for each row execute function procurement.serial_supply_guard_factory();
create function procurement.serial_supply_guard_stock() returns trigger language plpgsql as $$
declare ten uuid;unit text;po text;bound_stock text;released boolean;
begin
 if tg_op='DELETE' then ten:=old.tenant_id;unit:=old.production_unit_id;else ten:=new.tenant_id;unit:=new.production_unit_id;end if;
 select purchase_order_id into po from procurement.serial_supply_unit where tenant_id=ten and production_unit_id=unit;
 if not found then
  if tg_op='UPDATE' and exists(select 1 from procurement.serial_supply_manifest where tenant_id=old.tenant_id and stock_unit_id=old.stock_unit_id)
   then raise exception 'bound stock identity is immutable';end if;return null;
 end if;
 if tg_op='DELETE' then raise exception 'bound supply stock cannot be deleted';end if;
 select stock_unit_id into bound_stock from procurement.serial_supply_manifest where tenant_id=ten and production_unit_id=unit;
 if not found or bound_stock<>new.stock_unit_id then raise exception 'stock requires exact supply manifest';end if;
 if not exists(select 1 from factory.production_unit u where u.tenant_id=ten and u.production_unit_id=unit
  and row(u.variant_id,u.serial_number,u.vin,u.battery_serial_number) is not distinct from row(new.variant_id,new.serial_number,new.vin,new.battery_serial_number))
 then raise exception 'stock serial binding differs from factory';end if;
 select exists(select 1 from procurement.serial_supply_quality q join approval.request a on a.tenant_id=q.tenant_id and a.request_id=q.approval_id
  where q.tenant_id=ten and q.production_unit_id=unit and q.stage='receipt' and a.kind='serial_quality'
  and a.state='approved' and a.evidence_sha=q.payload_sha256) into released;
 -- After first approved receipt, legitimate commerce/transfer/service owners keep
 -- their existing state/version contract. Identity stays tied to the manifest.
 if released and tg_op='UPDATE' then return null;end if;
 if (tg_op='INSERT' or row(new.state,new.version) is distinct from row(old.state,old.version)) and not exists(
  select 1 from procurement.serial_supply_effect e where e.tenant_id=ten and e.purchase_order_id=po and e.entity='stock'
  and e.entity_id=new.stock_unit_id and e.from_state=case when tg_op='INSERT' then '' else old.state end
  and e.to_state=new.state and e.from_version=case when tg_op='INSERT' then 0 else old.version end
  and e.to_version=new.version and e.transaction_id=pg_current_xact_id())
 then raise exception 'bound supply stock requires connected command';end if;
 if new.state='available' and not released then raise exception 'supply stock requires human quality release';end if;
 if new.state='quarantine' and not exists(select 1 from procurement.serial_supply_receipt where tenant_id=ten and production_unit_id=unit)
 then raise exception 'supply quarantine requires receipt';end if;
 return null;
end $$;
create constraint trigger serial_supply_stock_guard after insert or update or delete on inventory.stock_unit
 deferrable initially deferred for each row execute function procurement.serial_supply_guard_stock();


create trigger serial_supply_plan_immutable before update or delete on procurement.serial_supply_plan for each row execute function procurement.serial_supply_fact_immutable();

create trigger serial_supply_line_immutable before update or delete on procurement.serial_supply_line for each row execute function procurement.serial_supply_fact_immutable();

create trigger serial_supply_step_immutable before update or delete on procurement.serial_supply_step for each row execute function procurement.serial_supply_fact_immutable();

create trigger serial_supply_unit_immutable before update or delete on procurement.serial_supply_unit for each row execute function procurement.serial_supply_fact_immutable();

create trigger serial_supply_shipment_immutable before update or delete on procurement.serial_supply_shipment for each row execute function procurement.serial_supply_fact_immutable();

create trigger serial_supply_manifest_immutable before update or delete on procurement.serial_supply_manifest for each row execute function procurement.serial_supply_fact_immutable();

create trigger serial_supply_receipt_immutable before update or delete on procurement.serial_supply_receipt for each row execute function procurement.serial_supply_fact_immutable();

create trigger serial_supply_quality_immutable before update or delete on procurement.serial_supply_quality for each row execute function procurement.serial_supply_fact_immutable();

create trigger serial_supply_effect_immutable before update or delete on procurement.serial_supply_effect for each row execute function procurement.serial_supply_fact_immutable();

create function procurement.serial_supply_guard_plan() returns trigger language plpgsql as $$
declare total bigint;lines bigint;
begin
 if not exists(select 1 from procurement.serial_supply_plan s join procurement.purchase_order p using(tenant_id,purchase_order_id)
  where s.tenant_id=new.tenant_id and s.purchase_order_id=new.purchase_order_id and p.state='draft'
  and s.supplier_id=p.supplier_id and s.destination_organization_id=p.destination_organization_id)
 then raise exception 'serial plan must bind an unchanged draft order';end if;
 select count(*),coalesce(sum(quantity),0) into lines,total from procurement.serial_supply_line
 where tenant_id=new.tenant_id and purchase_order_id=new.purchase_order_id;
 if lines<1 or lines>32 or total<1 or total>1000 then raise exception 'serial plan quantity bounds';end if;
 if not exists(select 1 from procurement.serial_supply_step st join procurement.serial_supply_plan p using(tenant_id,purchase_order_id)
  where st.tenant_id=new.tenant_id and st.purchase_order_id=new.purchase_order_id
  and st.kind='planned' and st.version=1 and st.actor=p.created_by and st.transaction_id=pg_current_xact_id())
 then raise exception 'serial plan lines require creation evidence in the same transaction';end if;
 return null;
end $$;
create constraint trigger serial_supply_plan_guard after insert on procurement.serial_supply_plan
 deferrable initially deferred for each row execute function procurement.serial_supply_guard_plan();
create constraint trigger serial_supply_line_guard after insert on procurement.serial_supply_line
 deferrable initially deferred for each row execute function procurement.serial_supply_guard_plan();


create function procurement.serial_supply_guard_approval() returns trigger language plpgsql as $$
begin
 if new.kind<>'serial_quality' then return null;end if;
 if not exists(select 1 from procurement.serial_supply_quality q join procurement.serial_supply_plan p using(tenant_id,purchase_order_id)
  where q.tenant_id=new.tenant_id and q.approval_id=new.request_id and q.payload_sha256=new.evidence_sha
  and new.subject_id=q.production_unit_id and new.amount_minor_units=0
  and new.organization_id=case q.stage when 'factory' then p.factory_organization_id else p.destination_organization_id end)
 then raise exception 'serial quality approval requires exact physical binding';end if;
 if tg_op='INSERT' and not exists(
  select 1 from procurement.serial_supply_quality q join procurement.serial_supply_step st using(tenant_id,purchase_order_id,command_id)
  where q.tenant_id=new.tenant_id and q.approval_id=new.request_id and st.actor=new.requester and st.transaction_id=pg_current_xact_id())
 then raise exception 'serial quality request requires connected evidence';end if;
 if tg_op='UPDATE' and new.state is distinct from old.state and not exists(
  select 1 from procurement.serial_supply_quality q join procurement.serial_supply_step st using(tenant_id,purchase_order_id)
  where q.tenant_id=new.tenant_id and q.approval_id=new.request_id and exists(select 1 from approval.decision d where d.tenant_id=new.tenant_id and d.request_id=new.request_id and d.reviewer=st.actor and d.approved=(new.state='approved'))
  and st.kind in ('milestone','quality','quality-reject') and st.payload->'effect'->>'approval_id'=new.request_id
  and st.transaction_id=pg_current_xact_id())
 then raise exception 'serial quality decision requires connected command';end if;
 return null;
end $$;
create constraint trigger serial_supply_approval_guard after insert or update on approval.request
 deferrable initially deferred for each row execute function procurement.serial_supply_guard_approval();

commit;
````

### FILE: `deploy/serial-supply/plan.reference.json`

```yaml
block_id: "GO-CONNECTED-SERIAL-SUPPLY:file5:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "16e2abb8c5482e6dbd47414bbf1a28c4fbae1c7e984fcbc08233bc183b2d9b5d"
variables: []
secrets_allowed: false
```

````json
{
  "purchase_order_id": "replace-with-existing-draft-po",
  "command_id": "plan-1",
  "factory_organization_id": "replace-with-factory-org",
  "demand_reference": "synthetic-reference-demand",
  "policy_code": "strict-serial-reference/v1",
  "evidence_sha256": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
  "lines": [
    {
      "id": "line-a",
      "variant_id": "replace-with-active-variant",
      "quantity": 3
    }
  ]
}
````

### FILE: `docs/SERIAL_SUPPLY_REFERENCE.md`

```yaml
block_id: "GO-CONNECTED-SERIAL-SUPPLY:file6:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "2e5d1afdbae29b8ba73df83e23ef5325d154b0eee97b97aa94402420c3d3b456"
variables: []
secrets_allowed: false
```

````markdown
# Connected serial supply reference — J2

Scope: LIBRARY_INFRASTRUCTURE, strict-serial-reference/v1. AUTHORED orchestration
of existing procurement, factory, shared human approval and inventory owners.
No new pricing or stock ledger. No third-party authorship is claimed for this
binding. BC-derived inventory owners retain their own provenance and notices.

The existing draft purchase order must already identify its active supplier,
destination, currency and total. Bind a demand reference, active factory and
1–32 unique variant lines, total1–1000 whole serialized units. The reference
permits no overreceipt. Its purchase total is unchanged; no unit price is inferred.
Missing VIN/battery is permitted; all present identifiers and serials remain unique.

## Activation

Apply all migrations in order through0072. Enable SERIAL_SUPPLY_ENABLED=true
and set SERIAL_SUPPLY_POLICY_SHA256=02767ecb3efb11f8ddaed3bdf237f3111024ad175e6001526f8f028c3171b8fa.
The SHA binds exact UTF-8 PolicyJSON bytes in internal/serialsupply/contract.go,
without an added newline. This explicitly selects the built-in synthetic policy:
{"schema":"elite-serial-supply-policy/v1","scope":"LIBRARY_INFRASTRUCTURE_REFERENCE","quantity":"whole-serialized-units","overreceipt":0,"receipt_state":"quarantine","release":"distinct-human-with-evidence","pricing":"existing-purchase-order-total-unchanged","production_authorized":false}

Disabled mode reads no policy. Enabled mode requires a pool, exact policy hash
and15enabled database guards. No credential, remote endpoint or account is used.
The existing host supplies database configuration and the selected identity
verifier. A target later assigns real organization-scoped permissions explicitly.

## Connected flow and roles

Every route includes organization_id once. Identity comes from the verifier;
request bodies cannot choose actor or tenant. Each route narrows a broad principal
to its one permission and selected organization. Versions are JSON decimal strings.

POST /v1/franchise/supply/orders/{id}/plan binds the immutable plan (supply:plan).
Franchise submit/cancel use supply:plan. Factory confirm/start/register/milestone/
ship use supply:factory under the bound factory organization. A released or rejected
quality milestone requires a distinct human and evidence. Rejected serial history
is retained; a replacement may consume the released plan slot.

Ship writes the immutable ASN/manifest and existing stock in-transit entries.
Ship/receive each take1–100 unique units, allowing split shipments and receipts.
Franchise receive requires supply:receive and manifest membership; stock enters
quarantine, excluded from existing available-to-promise. quality/quality-reject
require supply:release and a distinct reviewer. A rejected receipt stays quarantined.
reinspect requires supply:inspect, a prior rejection and a new evidence hash.
Only approved receipt quality makes the existing stock available. After approval,
ordinary existing commerce/transfer/service transitions remain authoritative.

GET /v1/{franchise|factory}/supply/orders/{id} uses supply:read or
supply:factory-read. Units are bounded100/page, with next_unit_id/after_unit.
GET with command_id recovers an immutable command receipt. Following an uncertain
POST, compare actor, purchase order, command ID and request SHA before deciding
whether an explicit identical retry is appropriate. Never generate a new command
ID merely because a response was lost. Actor or payload changes cannot replay.

All effects, approval decisions and outbox rows commit together. Deferred database
guards reject generic API/SQL bypass before release. Down0072 refuses populated
history; on an empty database down/up is supported. It never erases audit facts.

## Evidence and limits

SERIAL_SUPPLY_CONNECTED_RELEASE_V402.md/json in the canonical library binds exact
source, G0–G8 and reconstruction. Actual PG18.6 fixture: planned3, registered4,
one rejected/replaced, two ASNs, three received/quarantined/reviewed/available;
30versions. Eight concurrent identical commands: one new effect and seven replays.
Outbox failure rolls back15table snapshots. Two committed HTTP responses lost and
recovered by GET without another accepted POST. Generic approval/stock/PO bypass
rejected. Host policy/guard validation,13transport boundaries, populated/empty
downgrade and finite identity fuzz pass. JWT/provider/live factory certification
and role UI are not inferred; UI is the separate T2804 work item.
````

### FILE: `internal/approval/serial_quality_test.go`

```yaml
block_id: "GO-CONNECTED-SERIAL-SUPPLY:file7:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "81a49a0423c6494fe96190f93888e92eaa084770aa9490d6ac652ddc8c1da42d"
variables: []
secrets_allowed: false
```

````go
package approval

import (
	"strings"
	"testing"
)

func TestSerialQualityAlwaysNeedsDistinctHuman(t *testing.T) {
	r := NewRegistry(Policy{AutoApproveMinorUnits: 1000000})
	request := Request{TenantID: "fixture", ID: "serial-quality", Kind: KindSerialQuality, SubjectID: "unit", Requester: "receiver", EvidenceSHA: strings.Repeat("a", 64)}
	state, err := r.Submit(request)
	if err != nil || state != StatePending {
		t.Fatal("serial quality auto-approved", state, err)
	}
	if _, err = r.Approve("fixture", "serial-quality", "receiver", "self"); err != ErrSeparation {
		t.Fatal("self approval", err)
	}
	state, err = r.Approve("fixture", "serial-quality", "reviewer", "observed evidence")
	if err != nil || state != StateApproved {
		t.Fatal(state, err)
	}
}
````

### FILE: `internal/platform/httpapi/serial_supply.go`

```yaml
block_id: "GO-CONNECTED-SERIAL-SUPPLY:file8:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "7b484246331012a03a5d82b610b45dfef0962c094ed8f1b5a364f2d3c9504f8b"
variables: []
secrets_allowed: false
```

````go
// AUTHORED protected transport for the connected serial supply owner.
package httpapi

import (
	"bytes"
	"context"
	"elite.local/enterprise/internal/approval"
	"elite.local/enterprise/internal/operations"
	"elite.local/enterprise/internal/platform/identity"
	sc "elite.local/enterprise/internal/serialsupply"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"io"
	"net/http"
)

type SerialSupplyService interface {
	Plan(context.Context, identity.Principal, string, string) (sc.Plan, error)
	CommandReceipt(context.Context, identity.Principal, string, string) (sc.Receipt, error)
	BindPlan(context.Context, identity.Principal, sc.PlanRequest) (sc.Receipt, error)
	Apply(context.Context, identity.Principal, sc.Command) (sc.Receipt, error)
}
type SerialSupplyCreator interface {
	Create(context.Context, identity.Principal, sc.CreateRequest) (sc.Receipt, error)
}
type SerialSupplyModule struct{ Service SerialSupplyService }

func supplyHTTPError(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}
	var pg *pgconn.PgError
	switch {
	case errors.Is(err, sc.ErrInvalid):
		writeProblem(w, 400, "SUPPLY_INVALID", "invalid serial supply command")
	case errors.Is(err, sc.ErrNotFound), errors.Is(err, approval.ErrNotFound):
		writeProblem(w, 404, "SUPPLY_NOT_FOUND", "resource is outside the authorized scope")
	case errors.Is(err, sc.ErrConflict), errors.Is(err, operations.ErrConflict), errors.Is(err, approval.ErrSeparation), errors.Is(err, approval.ErrNotPending), errors.Is(err, approval.ErrInvalidRequest), errors.Is(err, approval.ErrDuplicate):
		writeProblem(w, 409, "SUPPLY_CONFLICT", "consult the saved command and current order version")
	case errors.As(err, &pg) && (pg.Code == "23505" || pg.Code == "23514" || pg.Code == "P0001" || pg.Code == "40001" || pg.Code == "40P01"):
		writeProblem(w, 409, "SUPPLY_CONFLICT", "consult the saved command and current order version")
	default:
		writeProblem(w, 503, "SUPPLY_UNCONFIRMED", "result unconfirmed; consult the saved command reference")
	}
	return true
}
func supplyDecode(w http.ResponseWriter, r *http.Request, out any) bool {
	raw, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 32768))
	if err != nil {
		writeProblem(w, 413, "SUPPLY_BODY_TOO_LARGE", "command exceeds its bounded size")
		return false
	}
	if _, _, err = approval.CanonicalPayload(raw); err != nil {
		writeProblem(w, 400, "SUPPLY_INVALID_JSON", "unambiguous JSON object required")
		return false
	}
	d := json.NewDecoder(bytes.NewReader(raw))
	d.DisallowUnknownFields()
	if d.Decode(out) != nil || d.Decode(new(any)) != io.EOF {
		writeProblem(w, 400, "SUPPLY_INVALID_JSON", "unknown fields or invalid typed values")
		return false
	}
	return true
}
func supplyReply(w http.ResponseWriter, v sc.Receipt) {
	status := http.StatusCreated
	if v.Replay {
		status = http.StatusOK
	}
	writeJSON(w, status, v)
}
func (m SerialSupplyModule) Register(mux *http.ServeMux, verifier identity.Verifier) {
	if m.Service == nil {
		return
	}
	api := franchiseJourneyAPI{verifier: verifier}
	protect := func(w http.ResponseWriter, r *http.Request, permission string) (identity.Principal, bool) {
		w.Header().Set("Cache-Control", "no-store")
		query := r.URL.Query()
		for key, values := range query {
			if len(values) != 1 || values[0] == "" || (key != "organization_id" && (r.Method != "GET" || (key != "command_id" && key != "after_unit"))) {
				writeProblem(w, 400, "SUPPLY_INVALID_QUERY", "one explicit organization and bounded query parameters required")
				return identity.Principal{}, false
			}
		}
		if query.Get("organization_id") == "" || (query.Get("command_id") != "" && query.Get("after_unit") != "") {
			writeProblem(w, 400, "SUPPLY_SCOPE_REQUIRED", "one organization_id and one recovery or page cursor required")
			return identity.Principal{}, false
		}
		org := query.Get("organization_id")
		p, ok := api.protected(w, r, permission, org)
		if !ok {
			return p, false
		}
		p.Permissions = map[string]struct{}{permission: {}}
		p.Organizations = map[string]struct{}{org: {}}
		return p, true
	}
	for _, surface := range []struct{ path, permission string }{{"franchise", "supply:read"}, {"factory", "supply:factory-read"}} {
		mux.HandleFunc("GET /v1/"+surface.path+"/supply/orders/{id}", func(w http.ResponseWriter, r *http.Request) {
			p, ok := protect(w, r, surface.permission)
			if !ok {
				return
			}
			if key := r.URL.Query().Get("command_id"); key != "" {
				v, err := m.Service.CommandReceipt(r.Context(), p, r.PathValue("id"), key)
				if !supplyHTTPError(w, err) {
					writeJSON(w, 200, v)
				}
				return
			}
			v, err := m.Service.Plan(r.Context(), p, r.PathValue("id"), r.URL.Query().Get("after_unit"))
			if !supplyHTTPError(w, err) {
				writeJSON(w, 200, v)
			}
		})
	}
	mux.HandleFunc("POST /v1/franchise/supply/orders/{id}/plan", func(w http.ResponseWriter, r *http.Request) {
		p, ok := protect(w, r, "supply:plan")
		if !ok {
			return
		}
		var c sc.PlanRequest
		if !supplyDecode(w, r, &c) {
			return
		}
		if c.PurchaseOrderID != r.PathValue("id") {
			supplyHTTPError(w, sc.ErrInvalid)
			return
		}
		v, err := m.Service.BindPlan(r.Context(), p, c)
		if !supplyHTTPError(w, err) {
			supplyReply(w, v)
		}
	})

	if creator, ok := m.Service.(SerialSupplyCreator); ok {
		mux.HandleFunc("POST /v1/franchise/supply/orders/{id}/create", func(w http.ResponseWriter, r *http.Request) {
			p, ok := protect(w, r, "supply:plan")
			if !ok {
				return
			}
			var c sc.CreateRequest
			if !supplyDecode(w, r, &c) {
				return
			}
			if c.PurchaseOrderID != r.PathValue("id") || c.DestinationOrganizationID != r.URL.Query().Get("organization_id") {
				supplyHTTPError(w, sc.ErrInvalid)
				return
			}
			v, err := creator.Create(r.Context(), p, c)
			if !supplyHTTPError(w, err) {
				supplyReply(w, v)
			}
		})
	}
	routes := []struct{ surface, kind, permission string }{
		{"franchise", "submit", "supply:plan"}, {"franchise", "cancel", "supply:plan"},
		{"factory", "confirm", "supply:factory"}, {"factory", "start", "supply:factory"},
		{"factory", "register", "supply:factory"}, {"factory", "milestone", "supply:factory"}, {"factory", "ship", "supply:factory"},
		{"franchise", "receive", "supply:receive"}, {"franchise", "quality", "supply:release"}, {"franchise", "quality-reject", "supply:release"}, {"franchise", "reinspect", "supply:inspect"},
	}
	for _, route := range routes {
		mux.HandleFunc("POST /v1/"+route.surface+"/supply/orders/{id}/"+route.kind, func(w http.ResponseWriter, r *http.Request) {
			p, ok := protect(w, r, route.permission)
			if !ok {
				return
			}
			var c sc.Command
			if !supplyDecode(w, r, &c) {
				return
			}
			if c.PurchaseOrderID != r.PathValue("id") || c.Kind != route.kind {
				supplyHTTPError(w, sc.ErrInvalid)
				return
			}
			v, err := m.Service.Apply(r.Context(), p, c)
			if !supplyHTTPError(w, err) {
				supplyReply(w, v)
			}
		})
	}
}
````

### FILE: `internal/platform/httpapi/serial_supply_test.go`

```yaml
block_id: "GO-CONNECTED-SERIAL-SUPPLY:file9:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "bf0ec794d8863de59d3de28f16e556aa7e6d3fa55eb6b38409096375cd6bfb26"
variables: []
secrets_allowed: false
```

````go
package httpapi

// AUTHORED trust-boundary regressions. The spy asserts no repository invocation
// for malformed commands and that broad identity grants are narrowed per route.
import (
	"context"
	"elite.local/enterprise/internal/platform/identity"
	sc "elite.local/enterprise/internal/serialsupply"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type supplyBoundaryVerifier struct{}

func (supplyBoundaryVerifier) Verify(context.Context, string) (identity.Principal, error) {
	return identity.Principal{Subject: "operator", TenantID: "tenant", Permissions: map[string]struct{}{"*": {}}, Organizations: map[string]struct{}{"*": {}}}, nil
}

type supplyBoundaryStore struct {
	calls   int
	p       identity.Principal
	command sc.Command
}

func (s *supplyBoundaryStore) Plan(context.Context, identity.Principal, string, string) (sc.Plan, error) {
	s.calls++
	return sc.Plan{}, nil
}
func (s *supplyBoundaryStore) CommandReceipt(context.Context, identity.Principal, string, string) (sc.Receipt, error) {
	s.calls++
	return sc.Receipt{}, nil
}
func (s *supplyBoundaryStore) BindPlan(context.Context, identity.Principal, sc.PlanRequest) (sc.Receipt, error) {
	s.calls++
	return sc.Receipt{}, nil
}
func (s *supplyBoundaryStore) Apply(_ context.Context, p identity.Principal, r sc.Command) (sc.Receipt, error) {
	s.calls++
	s.p = p
	s.command = r
	return sc.Receipt{Version: r.ExpectedVersion}, nil
}
func TestSerialSupplyHTTPBoundaries(t *testing.T) {
	raw := `{"purchase_order_id":"po","command_id":"receive-1","kind":"receive","expected_version":"9007199254740993","evidence_sha256":"` + strings.Repeat("a", 64) + `","shipment_id":"asn","units":["u1"]}`
	for _, c := range []struct {
		name, body, query string
		status            int
	}{
		{"exact-int64", raw, "organization_id=store", 201},
		{"number-not-string", strings.Replace(raw, `"9007199254740993"`, `9007199254740993`, 1), "organization_id=store", 400},
		{"duplicate-key", strings.Replace(raw, `"command_id":"receive-1"`, `"command_id":"receive-1","command_id":"changed"`, 1), "organization_id=store", 400},
		{"case-collision", strings.Replace(raw, `"command_id":"receive-1"`, `"command_id":"receive-1","COMMAND_ID":"changed"`, 1), "organization_id=store", 400},
		{"actor-injection", strings.TrimSuffix(raw, "}") + `,"actor":"other"}`, "organization_id=store", 400},
		{"path-body-mismatch", strings.Replace(raw, `"po"`, `"other"`, 1), "organization_id=store", 400},
		{"kind-route-mismatch", strings.Replace(raw, `"receive"`, `"ship"`, 1), "organization_id=store", 400},
		{"trailing-object", raw + "{}", "organization_id=store", 400},
		{"oversize", raw + strings.Repeat(" ", 32769), "organization_id=store", 413},
		{"missing-scope", raw, "", 400},
		{"repeated-scope", raw, "organization_id=store&organization_id=other", 400},
		{"write-cursor", raw, "organization_id=store&after_unit=u", 400},
		{"unknown-query", raw, "organization_id=store&actor=other", 400},
	} {
		t.Run(c.name, func(t *testing.T) {
			store := &supplyBoundaryStore{}
			mux := http.NewServeMux()
			SerialSupplyModule{Service: store}.Register(mux, supplyBoundaryVerifier{})
			req := httptest.NewRequest("POST", "/v1/franchise/supply/orders/po/receive?"+c.query, strings.NewReader(c.body))
			req.Header.Set("Authorization", "Bearer fixture")
			out := httptest.NewRecorder()
			mux.ServeHTTP(out, req)
			if out.Code != c.status || out.Header().Get("Cache-Control") != "no-store" {
				t.Fatal(out.Code, out.Body.String())
			}
			if c.status != 201 {
				if store.calls != 0 {
					t.Fatal("invalid input invoked repository")
				}
				return
			}
			if store.calls != 1 || store.command.ExpectedVersion != 9007199254740993 || len(store.p.Permissions) != 1 || !store.p.Allowed("supply:receive") || store.p.Allowed("supply:plan") || len(store.p.Organizations) != 1 {
				t.Fatal("authority/version changed", store)
			}
			if _, ok := store.p.Organizations["store"]; !ok {
				t.Fatal("route organization not bound")
			}
			var value map[string]any
			if err := json.Unmarshal(out.Body.Bytes(), &value); err != nil || value["version"] != "9007199254740993" {
				t.Fatal("response lost int64 precision", value, err)
			}
		})
	}
}
````

### FILE: `internal/platform/postgres/serial_supply.go`

```yaml
block_id: "GO-CONNECTED-SERIAL-SUPPLY:file10:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "5b204f8f46d74b58652c6f1644ba3fe4f31c4a054110785733cde4bb590b6930"
variables: []
secrets_allowed: false
```

````go
// AUTHORED physical/actor/evidence binding and recovery around existing owners.
package postgres

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"

	"elite.local/enterprise/internal/platform/identity"
	sc "elite.local/enterprise/internal/serialsupply"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SerialSupply struct{ pool *pgxpool.Pool }

func NewSerialSupply(pool *pgxpool.Pool) (*SerialSupply, error) {
	if pool == nil {
		return nil, sc.ErrInvalid
	}
	return &SerialSupply{pool: pool}, nil
}
func supplyError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return sc.ErrNotFound
	}
	var p *pgconn.PgError
	if postgresConflict(err) || errors.As(err, &p) && p.Code == "P0001" {
		return sc.ErrConflict
	}
	return err
}

type supplyQuery interface {
	QueryRow(context.Context, string, ...any) pgx.Row
	Query(context.Context, string, ...any) (pgx.Rows, error)
}

func readSupplyPlan(ctx context.Context, q supplyQuery, tenant, po string, lock bool) (sc.Plan, error) {
	var p sc.Plan
	sql := `select s.purchase_order_id,s.destination_organization_id,s.factory_organization_id,s.supplier_id,p.state,p.version,
 coalesce((select max(version) from procurement.serial_supply_step where tenant_id=s.tenant_id and purchase_order_id=s.purchase_order_id),0),
 s.policy_code,s.demand_reference,p.currency,p.total_minor_units from procurement.serial_supply_plan s join procurement.purchase_order p using(tenant_id,purchase_order_id)
 where s.tenant_id=$1 and s.purchase_order_id=$2`
	if lock {
		sql += " for update of p"
	}
	err := q.QueryRow(ctx, sql, tenant, po).Scan(&p.PurchaseOrderID, &p.DestinationOrganizationID, &p.FactoryOrganizationID, &p.SupplierID, &p.State, &p.PurchaseVersion, &p.Version, &p.PolicyCode, &p.DemandReference, &p.Currency, &p.TotalMinorUnits)
	return p, supplyError(err)
}
func readSupplyReceipt(ctx context.Context, q supplyQuery, tenant, po, command string) (sc.Receipt, error) {
	var r sc.Receipt
	err := q.QueryRow(ctx, `select purchase_order_id,command_id,version,kind,actor,request_sha256,payload_sha256,payload,recorded_at
 from procurement.serial_supply_step where tenant_id=$1 and purchase_order_id=$2 and ($3::text='' or command_id=$3) order by version desc limit 1`, tenant, po, command).Scan(&r.PurchaseOrderID, &r.CommandID, &r.Version, &r.Kind, &r.Actor, &r.RequestSHA256, &r.PayloadSHA256, &r.Payload, &r.RecordedAt)
	if err != nil {
		return r, supplyError(err)
	}
	// CanonicalPayload already enforces exact JSON numbers; decode through that
	// shared owner directly, not via float64 in an interface.
	raw, hash, err := sc.Canonical(r.Payload)
	if err != nil || hash != r.PayloadSHA256 {
		return sc.Receipt{}, sc.ErrConflict
	}
	r.Payload = raw
	return r, nil
}
func supplyReadAllowed(p identity.Principal, plan sc.Plan) bool {
	return approvalPrincipal(p, p.TenantID, plan.DestinationOrganizationID, "supply:read") ||
		approvalPrincipal(p, p.TenantID, plan.FactoryOrganizationID, "supply:factory-read")
}
func (s *SerialSupply) Plan(ctx context.Context, p identity.Principal, id, afterUnit string) (sc.Plan, error) {
	if s == nil || !sc.ValidID(id) || afterUnit != "" && !sc.ValidID(afterUnit) || p.TenantID == "" {
		return sc.Plan{}, sc.ErrInvalid
	}
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead, AccessMode: pgx.ReadOnly})
	if err != nil {
		return sc.Plan{}, err
	}
	defer tx.Rollback(ctx)
	plan, err := readSupplyPlan(ctx, tx, p.TenantID, id, false)
	if err != nil {
		return sc.Plan{}, err
	}
	if !supplyReadAllowed(p, plan) {
		return sc.Plan{}, sc.ErrNotFound
	}
	rows, err := tx.Query(ctx, `select line_id,variant_id,quantity from procurement.serial_supply_line where tenant_id=$1 and purchase_order_id=$2 order by line_id`, p.TenantID, id)
	if err != nil {
		return sc.Plan{}, err
	}
	plan.Lines = []sc.Line{}
	for rows.Next() {
		var l sc.Line
		if err = rows.Scan(&l.ID, &l.VariantID, &l.Quantity); err != nil {
			rows.Close()
			return sc.Plan{}, err
		}
		plan.Lines = append(plan.Lines, l)
	}
	rows.Close()
	if err = rows.Err(); err != nil {
		return sc.Plan{}, err
	}
	rows, err = tx.Query(ctx, `select u.production_unit_id,b.line_id,u.state,u.serial_number,coalesce(m.stock_unit_id,''),coalesce(i.state,''),coalesce(m.shipment_id,''),
 coalesce((select a.state from procurement.serial_supply_quality q join approval.request a on a.tenant_id=q.tenant_id and a.request_id=q.approval_id where q.tenant_id=u.tenant_id and q.production_unit_id=u.production_unit_id and q.stage='factory' order by q.attempt desc limit 1),''),
 coalesce((select a.requester from procurement.serial_supply_quality q join approval.request a on a.tenant_id=q.tenant_id and a.request_id=q.approval_id where q.tenant_id=u.tenant_id and q.production_unit_id=u.production_unit_id and q.stage='factory' order by q.attempt desc limit 1),''),
 coalesce((select a.state from procurement.serial_supply_quality q join approval.request a on a.tenant_id=q.tenant_id and a.request_id=q.approval_id where q.tenant_id=u.tenant_id and q.production_unit_id=u.production_unit_id and q.stage='receipt' order by q.attempt desc limit 1),''),
 coalesce((select a.requester from procurement.serial_supply_quality q join approval.request a on a.tenant_id=q.tenant_id and a.request_id=q.approval_id where q.tenant_id=u.tenant_id and q.production_unit_id=u.production_unit_id and q.stage='receipt' order by q.attempt desc limit 1),'')
 from procurement.serial_supply_unit b join factory.production_unit u using(tenant_id,production_unit_id)
 left join procurement.serial_supply_manifest m using(tenant_id,production_unit_id)
 left join inventory.stock_unit i on i.tenant_id=m.tenant_id and i.stock_unit_id=m.stock_unit_id
 where b.tenant_id=$1 and b.purchase_order_id=$2 and u.production_unit_id>$3 order by u.production_unit_id limit 101`, p.TenantID, id, afterUnit)
	if err != nil {
		return sc.Plan{}, err
	}
	plan.Units = []sc.Unit{}
	for rows.Next() {
		var u sc.Unit
		if err = rows.Scan(&u.ID, &u.LineID, &u.State, &u.SerialNumber, &u.StockUnitID, &u.StockState, &u.ShipmentID, &u.FactoryReviewState, &u.FactoryRequester, &u.ReceiptReviewState, &u.ReceiptRequester); err != nil {
			rows.Close()
			return sc.Plan{}, err
		}
		plan.Units = append(plan.Units, u)
	}
	rows.Close()
	if err = rows.Err(); err != nil {
		return sc.Plan{}, err
	}
	if len(plan.Units) > 100 {
		plan.Units = plan.Units[:100]
		plan.NextUnitID = plan.Units[99].ID
	}
	latest, err := readSupplyReceipt(ctx, tx, p.TenantID, id, "")
	if err != nil {
		return sc.Plan{}, err
	}
	plan.Latest = &latest
	if err = tx.Commit(ctx); err != nil {
		return sc.Plan{}, err
	}
	return plan, nil
}
func (s *SerialSupply) CommandReceipt(ctx context.Context, p identity.Principal, id, command string) (sc.Receipt, error) {
	if s == nil || !sc.ValidID(id) || !sc.ValidID(command) || p.TenantID == "" {
		return sc.Receipt{}, sc.ErrInvalid
	}
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead, AccessMode: pgx.ReadOnly})
	if err != nil {
		return sc.Receipt{}, err
	}
	defer tx.Rollback(ctx)
	plan, err := readSupplyPlan(ctx, tx, p.TenantID, id, false)
	if err != nil {
		return sc.Receipt{}, err
	}
	if !supplyReadAllowed(p, plan) {
		return sc.Receipt{}, sc.ErrNotFound
	}
	r, err := readSupplyReceipt(ctx, tx, p.TenantID, id, command)
	if err != nil {
		return sc.Receipt{}, err
	}
	if err = tx.Commit(ctx); err != nil {
		return sc.Receipt{}, err
	}
	return r, nil
}
func writeSupplyReceipt(ctx context.Context, tx pgx.Tx, tenant, po, command, kind, actor, requestSHA, evidence string, version int64, payload any) (sc.Receipt, error) {
	raw, hash, err := sc.Canonical(payload)
	if err != nil {
		return sc.Receipt{}, err
	}
	_, err = tx.Exec(ctx, `insert into procurement.serial_supply_step(tenant_id,purchase_order_id,command_id,version,kind,actor,request_sha256,evidence_sha256,payload_sha256,payload) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`, tenant, po, command, version, kind, actor, requestSHA, evidence, hash, raw)
	if err != nil {
		return sc.Receipt{}, supplyError(err)
	}
	r, err := readSupplyReceipt(ctx, tx, tenant, po, command)
	if err != nil {
		return sc.Receipt{}, err
	}
	event, err := json.Marshal(r)
	if err != nil {
		return sc.Receipt{}, err
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)
 values($1,gen_random_uuid(),'serial-supply',$2,$3,$4,1,clock_timestamp(),$5)`, tenant, po, version, "serial-supply."+kind, event)
	return r, err
}
func supplyEffect(ctx context.Context, tx pgx.Tx, tenant, po, command, entity, id, from, to string, fromVersion, toVersion int64) error {
	_, err := tx.Exec(ctx, `insert into procurement.serial_supply_effect(tenant_id,purchase_order_id,command_id,entity,entity_id,from_state,to_state,from_version,to_version)
 values($1,$2,$3,$4,$5,$6,$7,$8,$9)`, tenant, po, command, entity, id, from, to, fromVersion, toVersion)
	return err
}
func (s *SerialSupply) BindPlan(ctx context.Context, p identity.Principal, r sc.PlanRequest) (sc.Receipt, error) {
	if s == nil || !r.Valid() || p.TenantID == "" {
		return sc.Receipt{}, sc.ErrInvalid
	}
	_, requestSHA, err := sc.Canonical(r)
	if err != nil {
		return sc.Receipt{}, err
	}
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return sc.Receipt{}, err
	}
	defer tx.Rollback(ctx)
	receipt, err := bindSupplyPlanInTx(ctx, tx, p, r, requestSHA)
	if err != nil {
		return sc.Receipt{}, err
	}
	return receipt, tx.Commit(ctx)
}

// AUTHORED extraction: original BindPlan SQL/authorization/replay are retained.
func bindSupplyPlanInTx(ctx context.Context, tx pgx.Tx, p identity.Principal, r sc.PlanRequest, requestSHA string) (sc.Receipt, error) {
	var org, supplier, state string
	err := tx.QueryRow(ctx, `select destination_organization_id,supplier_id,state from procurement.purchase_order where tenant_id=$1 and purchase_order_id=$2 for update`, p.TenantID, r.PurchaseOrderID).Scan(&org, &supplier, &state)
	if err != nil {
		return sc.Receipt{}, supplyError(err)
	}
	if !approvalPrincipal(p, p.TenantID, org, "supply:plan") {
		return sc.Receipt{}, sc.ErrNotFound
	}
	existing, err := readSupplyReceipt(ctx, tx, p.TenantID, r.PurchaseOrderID, r.CommandID)
	if err == nil {
		if existing.Kind != "planned" || existing.Actor != p.Subject || existing.RequestSHA256 != requestSHA {
			return sc.Receipt{}, sc.ErrConflict
		}
		existing.Replay = true
		return existing, nil
	}
	if !errors.Is(err, sc.ErrNotFound) {
		return sc.Receipt{}, err
	}
	if state != "draft" {
		return sc.Receipt{}, sc.ErrConflict
	}
	var active bool
	err = tx.QueryRow(ctx, `select exists(select 1 from partner.supplier where tenant_id=$1 and supplier_id=$2 and status='active')
 and exists(select 1 from org.organization where tenant_id=$1 and organization_id=$3 and status='active')
 and exists(select 1 from org.organization where tenant_id=$1 and organization_id=$4 and status='active' and organization_type='factory') and not exists(select 1 from factory.production_unit where tenant_id=$1 and purchase_order_id=$5)`, p.TenantID, supplier, org, r.FactoryOrganizationID, r.PurchaseOrderID).Scan(&active)
	if err != nil {
		return sc.Receipt{}, err
	}
	if !active {
		return sc.Receipt{}, sc.ErrConflict
	}
	policyHash := sha256.Sum256([]byte(sc.PolicyJSON))
	_, err = tx.Exec(ctx, `insert into procurement.serial_supply_plan(tenant_id,purchase_order_id,destination_organization_id,factory_organization_id,supplier_id,demand_reference,policy_code,policy_sha256,created_by)
 values($1,$2,$3,$4,$5,$6,$7,$8,$9)`, p.TenantID, r.PurchaseOrderID, org, r.FactoryOrganizationID, supplier, r.DemandReference, sc.PolicyCode, hex.EncodeToString(policyHash[:]), p.Subject)
	if err != nil {
		return sc.Receipt{}, supplyError(err)
	}
	for _, line := range r.Lines {
		tag, err := tx.Exec(ctx, `insert into procurement.serial_supply_line(tenant_id,purchase_order_id,line_id,variant_id,quantity)
 select $1,$2,$3,variant_id,$5 from catalog.vehicle_variant where tenant_id=$1 and variant_id=$4 and lifecycle_state='active'`, p.TenantID, r.PurchaseOrderID, line.ID, line.VariantID, line.Quantity)
		if err != nil {
			return sc.Receipt{}, supplyError(err)
		}
		if tag.RowsAffected() != 1 {
			return sc.Receipt{}, sc.ErrConflict
		}
	}
	receipt, err := writeSupplyReceipt(ctx, tx, p.TenantID, r.PurchaseOrderID, r.CommandID, "planned", p.Subject, requestSHA, r.EvidenceSHA256, 1, map[string]any{"plan": r, "destination_organization_id": org, "supplier_id": supplier, "policy_sha256": hex.EncodeToString(policyHash[:])})
	if err != nil {
		return sc.Receipt{}, err
	}
	return receipt, nil
}
````

### FILE: `internal/platform/postgres/serial_supply_commands.go`

```yaml
block_id: "GO-CONNECTED-SERIAL-SUPPLY:file11:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "1755c91f9fec3f4a61f5c03eae9e2f0af3722c175a538cc0bc289708caa07bf4"
variables: []
secrets_allowed: false
```

````go
// AUTHORED transaction composition. Domain preparation/transitions remain in
// operations.Service and its original PostgreSQL writers.
package postgres

import (
	"context"
	"elite.local/enterprise/internal/operations"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/randomid"
	sc "elite.local/enterprise/internal/serialsupply"
	"errors"
	"github.com/jackc/pgx/v5"
)

func supplyCommandPermission(r sc.Command, plan sc.Plan) (string, string) {
	switch r.Kind {
	case "confirm", "start", "register", "milestone", "ship":
		return plan.FactoryOrganizationID, "supply:factory"
	case "receive":
		return plan.DestinationOrganizationID, "supply:receive"
	case "quality", "quality-reject":
		return plan.DestinationOrganizationID, "supply:release"
	case "reinspect":
		return plan.DestinationOrganizationID, "supply:inspect"
	default:
		return plan.DestinationOrganizationID, "supply:plan"
	}
}
func moveSupplyPurchase(ctx context.Context, tx pgx.Tx, tenant string, p sc.Plan, r sc.Command, target string) error {
	svc := operations.NewService(operationsTxRepository{tx: tx}, randomid.Generator{})
	if err := svc.TransitionPurchaseOrder(ctx, tenant, p.DestinationOrganizationID, p.PurchaseOrderID, p.State, target, p.PurchaseVersion); err != nil {
		return supplyError(err)
	}
	return supplyEffect(ctx, tx, tenant, p.PurchaseOrderID, r.CommandID, "purchase", p.PurchaseOrderID, p.State, target, p.PurchaseVersion, p.PurchaseVersion+1)
}
func moveSupplyFactory(ctx context.Context, tx pgx.Tx, tenant string, p sc.Plan, r sc.Command, id, from, to string) error {
	svc := operations.NewService(operationsTxRepository{tx: tx}, randomid.Generator{})
	if err := svc.TransitionProductionUnit(ctx, tenant, p.DestinationOrganizationID, id, from, to); err != nil {
		return supplyError(err)
	}
	return supplyEffect(ctx, tx, tenant, p.PurchaseOrderID, r.CommandID, "factory", id, from, to, 0, 0)
}
func moveSupplyStock(ctx context.Context, tx pgx.Tx, tenant string, p sc.Plan, r sc.Command, id, from, to string, version int64) error {
	svc := operations.NewService(operationsTxRepository{tx: tx}, randomid.Generator{})
	if err := svc.TransitionStockUnit(ctx, tenant, p.DestinationOrganizationID, id, from, to, version); err != nil {
		return supplyError(err)
	}
	return supplyEffect(ctx, tx, tenant, p.PurchaseOrderID, r.CommandID, "stock", id, from, to, version, version+1)
}

type supplyUnitState struct {
	ID, LineID, VariantID, State, Serial, VIN, Battery string
	StockID, StockState, ShipmentID                    string
	StockVersion                                       int64
}

func readSupplyUnit(ctx context.Context, tx pgx.Tx, tenant, po, id string) (supplyUnitState, error) {
	var u supplyUnitState
	err := tx.QueryRow(ctx, `select u.production_unit_id,b.line_id,u.variant_id,u.state,u.serial_number,coalesce(u.vin,''),coalesce(u.battery_serial_number,''),
 coalesce(m.stock_unit_id,''),coalesce(i.state,''),coalesce(m.shipment_id,''),coalesce(i.version,0)
 from procurement.serial_supply_unit b join factory.production_unit u using(tenant_id,production_unit_id)
 left join procurement.serial_supply_manifest m using(tenant_id,production_unit_id)
 left join inventory.stock_unit i on i.tenant_id=m.tenant_id and i.stock_unit_id=m.stock_unit_id
 where b.tenant_id=$1 and b.purchase_order_id=$2 and b.production_unit_id=$3 for update of u`, tenant, po, id).Scan(&u.ID, &u.LineID, &u.VariantID, &u.State, &u.Serial, &u.VIN, &u.Battery, &u.StockID, &u.StockState, &u.ShipmentID, &u.StockVersion)
	if err != nil {
		return u, supplyError(err)
	}
	if u.StockID != "" {
		// Lock the stock owner too, rather than relying on a stale outer-join snapshot.
		var org, variant, unit, serial, vin, battery string
		err = tx.QueryRow(ctx, `select organization_id,variant_id,coalesce(production_unit_id,''),serial_number,coalesce(vin,''),coalesce(battery_serial_number,''),state,version
  from inventory.stock_unit where tenant_id=$1 and stock_unit_id=$2 for update`, tenant, u.StockID).Scan(&org, &variant, &unit, &serial, &vin, &battery, &u.StockState, &u.StockVersion)
		if err != nil {
			return u, supplyError(err)
		}
		var destination string
		if err = tx.QueryRow(ctx, `select destination_organization_id from procurement.serial_supply_plan where tenant_id=$1 and purchase_order_id=$2`, tenant, po).Scan(&destination); err != nil {
			return u, err
		}
		if org != destination || variant != u.VariantID || unit != u.ID || serial != u.Serial || vin != u.VIN || battery != u.Battery {
			return u, sc.ErrConflict
		}
	}
	return u, nil
}
func registerSupplyUnit(ctx context.Context, tx pgx.Tx, tenant string, p sc.Plan, r sc.Command) (any, error) {
	if p.State != "in-production" {
		return nil, sc.ErrConflict
	}
	var variant string
	var maximum, active int
	err := tx.QueryRow(ctx, `select l.variant_id,l.quantity,
 (select count(*) from procurement.serial_supply_unit b join factory.production_unit u using(tenant_id,production_unit_id)
 where b.tenant_id=l.tenant_id and b.purchase_order_id=l.purchase_order_id and b.line_id=l.line_id and u.state<>'rejected')
 from procurement.serial_supply_line l where l.tenant_id=$1 and l.purchase_order_id=$2 and l.line_id=$3`, tenant, p.PurchaseOrderID, r.LineID).Scan(&variant, &maximum, &active)
	if err != nil {
		return nil, supplyError(err)
	}
	if active >= maximum {
		return nil, sc.ErrConflict
	}
	svc := operations.NewService(operationsTxRepository{tx: tx}, randomid.Generator{})
	unit, err := svc.RegisterProductionUnit(ctx, tenant, operations.ProductionUnit{OrganizationID: p.DestinationOrganizationID, PurchaseOrderID: p.PurchaseOrderID, VariantID: variant, SerialNumber: r.SerialNumber, VIN: r.VIN, BatterySerialNumber: r.BatterySerialNumber})
	if err != nil {
		return nil, supplyError(err)
	}
	_, err = tx.Exec(ctx, `insert into procurement.serial_supply_unit(tenant_id,purchase_order_id,line_id,production_unit_id,command_id) values($1,$2,$3,$4,$5)`, tenant, p.PurchaseOrderID, r.LineID, unit.ID, r.CommandID)
	if err != nil {
		return nil, supplyError(err)
	}
	if err = supplyEffect(ctx, tx, tenant, p.PurchaseOrderID, r.CommandID, "factory", unit.ID, "", "planned", 0, 0); err != nil {
		return nil, err
	}
	return map[string]any{"unit": unit, "line_id": r.LineID}, nil
}
func supplyLifecycle(ctx context.Context, tx pgx.Tx, tenant string, p sc.Plan, r sc.Command) (any, error) {
	target := map[string]string{"submit": "submitted", "confirm": "accepted", "start": "in-production", "cancel": "cancelled"}[r.Kind]
	if target == "" {
		return nil, sc.ErrInvalid
	}
	if r.Kind == "confirm" {
		var active bool
		if err := tx.QueryRow(ctx, `select exists(select 1 from partner.supplier where tenant_id=$1 and supplier_id=$2 and status='active')`, tenant, p.SupplierID).Scan(&active); err != nil {
			return nil, err
		}
		if !active {
			return nil, sc.ErrConflict
		}
	}
	if err := moveSupplyPurchase(ctx, tx, tenant, p, r, target); err != nil {
		return nil, err
	}
	payload := map[string]any{"purchase_order_id": p.PurchaseOrderID, "from_state": p.State, "to_state": target, "supplier_id": p.SupplierID}
	if r.Kind == "confirm" {
		payload["confirmation_kind"] = "authorized-recorded-evidence"
	}
	return payload, nil
}

type supplyAction func(context.Context, pgx.Tx, identity.Principal, sc.Plan, sc.Command) (any, error)

// executeSupplyCommand supplies the single lock/replay/receipt boundary for all
// commands. The action is selected internally by the typed public dispatcher.
func (s *SerialSupply) executeSupplyCommand(ctx context.Context, p identity.Principal, r sc.Command, action supplyAction) (sc.Receipt, error) {
	if s == nil || !r.Valid() || p.TenantID == "" || action == nil {
		return sc.Receipt{}, sc.ErrInvalid
	}
	_, hash, err := sc.Canonical(r)
	if err != nil {
		return sc.Receipt{}, err
	}
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return sc.Receipt{}, err
	}
	defer tx.Rollback(ctx)
	plan, err := readSupplyPlan(ctx, tx, p.TenantID, r.PurchaseOrderID, true)
	if err != nil {
		return sc.Receipt{}, err
	}
	org, permission := supplyCommandPermission(r, plan)
	if !approvalPrincipal(p, p.TenantID, org, permission) {
		return sc.Receipt{}, sc.ErrNotFound
	}
	old, err := readSupplyReceipt(ctx, tx, p.TenantID, r.PurchaseOrderID, r.CommandID)
	if err == nil {
		if old.Kind != r.Kind || old.Actor != p.Subject || old.RequestSHA256 != hash {
			return sc.Receipt{}, sc.ErrConflict
		}
		old.Replay = true
		return old, supplyError(tx.Commit(ctx))
	}
	if !errors.Is(err, sc.ErrNotFound) {
		return sc.Receipt{}, err
	}
	if plan.Version != r.ExpectedVersion {
		return sc.Receipt{}, sc.ErrConflict
	}
	payload, err := action(ctx, tx, p, plan, r)
	if err != nil {
		return sc.Receipt{}, err
	}
	receipt, err := writeSupplyReceipt(ctx, tx, p.TenantID, r.PurchaseOrderID, r.CommandID, r.Kind, p.Subject, hash, r.EvidenceSHA256, plan.Version+1, map[string]any{"effect": payload, "evidence_sha256": r.EvidenceSHA256})
	if err != nil {
		return sc.Receipt{}, err
	}
	return receipt, supplyError(tx.Commit(ctx))
}
````

### FILE: `internal/platform/postgres/serial_supply_connected_integration_test.go`

```yaml
block_id: "GO-CONNECTED-SERIAL-SUPPLY:file12:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "0536c8b909faa89e38558a1f35e1555874aef8ad99dc467ddd7550644d574d11"
variables: []
secrets_allowed: false
```

````go
// AUTHORED connected J2 reference: real source owners, scoped fixture principals.
// No provider accounts, legal quality certification or production acceptance.
package postgres_test

import (
	"context"
	"elite.local/enterprise/internal/operations"
	"elite.local/enterprise/internal/platform/identity"
	db "elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/platform/randomid"
	sc "elite.local/enterprise/internal/serialsupply"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

type serialSupplyReference interface {
	BindPlan(context.Context, identity.Principal, sc.PlanRequest) (sc.Receipt, error)
	Apply(context.Context, identity.Principal, sc.Command) (sc.Receipt, error)
	Plan(context.Context, identity.Principal, string, string) (sc.Plan, error)
	CommandReceipt(context.Context, identity.Principal, string, string) (sc.Receipt, error)
}
type serialSupplyHook func(*testing.T, *pgxpool.Pool, *db.SerialSupply) serialSupplyReference

func TestSerialSupplyConnectedReference(t *testing.T) { testSerialSupplyConnected(t, nil) }
func testSerialSupplyConnected(t *testing.T, hook serialSupplyHook) {
	url := os.Getenv("PAYMENT_CONNECTED_DB_URL")
	if url == "" {
		t.Skip("owned fixture DB not configured")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := "ca6c9cb7-e548-4f19-a2e6-27d2aa2d4892"
	fixtures := []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,'supply-j2','Fixture','Fixture')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,'j2-store','j2-store','Fixture','store'),($1,'j2-factory','j2-factory','Fixture','factory')`,
		`insert into partner.supplier(tenant_id,supplier_id,supplier_code,legal_name,status) values($1,'j2-supplier','j2-supplier','Fixture','active')`,
		`insert into catalog.vehicle_model(tenant_id,model_id,model_code,display_name,vehicle_class,lifecycle_state) values($1,'j2-model','j2-model','Fixture','bicycle','active')`,
		`insert into catalog.vehicle_variant(tenant_id,variant_id,model_id,variant_code,display_name,battery_specification,lifecycle_state) values($1,'j2-a','j2-model','j2-a','Fixture','{}','active'),($1,'j2-b','j2-model','j2-b','Fixture','{}','active')`,
	}
	for _, q := range fixtures {
		if _, err = pool.Exec(ctx, q, tenant); err != nil {
			t.Fatal(err)
		}
	}
	person := func(subject, org string, permissions ...string) identity.Principal {
		p := identity.Principal{Subject: subject, TenantID: tenant, Permissions: map[string]struct{}{}, Organizations: map[string]struct{}{org: {}}}
		for _, v := range permissions {
			p.Permissions[v] = struct{}{}
		}
		return p
	}
	buyer := person("j2-buyer", "j2-store", "supply:plan", "supply:read")
	maker := person("j2-maker", "j2-factory", "supply:factory", "supply:factory-read")
	factoryQA := person("j2-factory-qa", "j2-factory", "supply:factory", "supply:factory-read")
	receiver := person("j2-receiver", "j2-store", "supply:receive", "supply:inspect", "supply:read")
	qa := person("j2-receipt-qa", "j2-store", "supply:release", "supply:read")
	store, err := db.NewSerialSupply(pool)
	if err != nil {
		t.Fatal(err)
	}
	var supply serialSupplyReference = store
	if hook != nil {
		supply = hook(t, pool, store)
	}
	ops := operations.NewService(db.NewOperations(pool), randomid.Generator{})
	po, err := ops.CreatePurchaseOrder(ctx, tenant, operations.PurchaseOrder{SupplierID: "j2-supplier", DestinationOrganizationID: "j2-store", Currency: "ARS", TotalMinorUnits: 90000})
	if err != nil {
		t.Fatal(err)
	}
	planRequest := sc.PlanRequest{PurchaseOrderID: po.ID, CommandID: "plan", FactoryOrganizationID: "j2-factory", DemandReference: "fixture-demand-3-serials", PolicyCode: sc.PolicyCode, EvidenceSHA256: strings.Repeat("a", 64), Lines: []sc.Line{{ID: "a", VariantID: "j2-a", Quantity: 2}, {ID: "b", VariantID: "j2-b", Quantity: 1}}}
	first, err := supply.BindPlan(ctx, buyer, planRequest)
	if err != nil {
		t.Fatal("bind plan", err)
	}
	replay, err := supply.BindPlan(ctx, buyer, planRequest)
	if err != nil || !replay.Replay || replay.RequestSHA256 != first.RequestSHA256 {
		t.Fatal("plan replay", replay, err)
	}
	wrong := person("j2-other", "other", "supply:plan", "supply:read", "supply:receive")
	if _, err = supply.Plan(ctx, wrong, po.ID, ""); err == nil {
		t.Fatal("cross organization plan read")
	}
	version := first.Version
	sequence := 0
	command := func(kind string) sc.Command {
		sequence++
		return sc.Command{PurchaseOrderID: po.ID, CommandID: fmt.Sprintf("j2-%03d", sequence), ExpectedVersion: version, Kind: kind, EvidenceSHA256: strings.Repeat("a", 64)}
	}
	apply := func(p identity.Principal, r sc.Command) sc.Receipt {
		t.Helper()
		out, err := supply.Apply(ctx, p, r)
		if err != nil {
			t.Fatalf("%s/%s: %v", r.Kind, r.CommandID, err)
		}
		if out.Replay || out.Version != version+1 || out.Actor != p.Subject {
			t.Fatalf("unexpected receipt: %+v", out)
		}
		version = out.Version
		return out
	}
	negative := func(p identity.Principal, r sc.Command, label string) {
		t.Helper()
		if _, err := supply.Apply(ctx, p, r); err == nil {
			t.Fatal(label + " accepted")
		}
	}
	snapshot := func() string {
		t.Helper()
		out := ""
		for _, table := range []string{"procurement.purchase_order", "factory.production_unit", "inventory.stock_unit", "procurement.serial_supply_plan", "procurement.serial_supply_line", "procurement.serial_supply_unit", "procurement.serial_supply_shipment", "procurement.serial_supply_manifest", "procurement.serial_supply_receipt", "procurement.serial_supply_step", "procurement.serial_supply_quality", "procurement.serial_supply_effect", "approval.request", "approval.decision", "platform.outbox_event"} {
			var raw string
			query := fmt.Sprintf("select coalesce(jsonb_agg(to_jsonb(x) order by to_jsonb(x)::text),'[]'::jsonb)::text from %s x where tenant_id=$1", table)
			if err := pool.QueryRow(ctx, query, tenant).Scan(&raw); err != nil {
				t.Fatal(table, err)
			}
			out += table + raw
		}
		return out
	}
	apply(buyer, command("submit"))
	negative(buyer, command("confirm"), "buyer acting as factory")
	apply(maker, command("confirm"))
	apply(maker, command("start"))
	register := func(line, serial string) string {
		t.Helper()
		r := command("register")
		r.LineID = line
		r.SerialNumber = serial
		out := apply(maker, r)
		var payload struct {
			Effect struct {
				Unit operations.ProductionUnit `json:"unit"`
			} `json:"effect"`
		}
		if err := json.Unmarshal(out.Payload, &payload); err != nil || payload.Effect.Unit.ID == "" {
			t.Fatal(err, string(out.Payload))
		}
		return payload.Effect.Unit.ID
	}
	milestone := func(actor identity.Principal, unit, target string) sc.Receipt {
		t.Helper()
		r := command("milestone")
		r.UnitID = unit
		r.TargetState = target
		if target == "released" || target == "rejected" {
			r.Reason = "fixture inspection and recorded evidence"
		}
		return apply(actor, r)
	}
	// A rejected factory serial is retained and a replacement consumes its line slot.
	rejected := register("a", "J2-REJECTED")
	milestone(maker, rejected, "assembly")
	milestone(maker, rejected, "quality")
	before := snapshot()
	r := command("milestone")
	r.UnitID = rejected
	r.TargetState = "released"
	r.Reason = "self review forbidden"
	negative(maker, r, "factory self approval")
	if snapshot() != before {
		t.Fatal("self approval changed durable state")
	}
	milestone(factoryQA, rejected, "rejected")
	a1 := register("a", "J2-A1")
	a2 := register("a", "J2-A2")
	b1 := register("b", "J2-B1")
	r = command("register")
	r.LineID = "a"
	r.SerialNumber = "J2-OVER"
	negative(maker, r, "overproduction")
	r = command("register")
	r.LineID = "b"
	r.SerialNumber = "J2-A1"
	negative(maker, r, "reused serial/quantity")
	r = command("ship")
	r.ShipmentID = "premature"
	r.Units = []string{a1}
	negative(maker, r, "premature shipment")
	// Legacy operations may not bypass the connected plan even with the right org.
	before = snapshot()
	if err = ops.TransitionProductionUnit(ctx, tenant, "j2-store", a1, "planned", "assembly"); err == nil {
		t.Fatal("generic factory bypass")
	}
	if snapshot() != before {
		t.Fatal("rejected generic command leaked state/outbox")
	}
	for _, id := range []string{a1, a2, b1} {
		milestone(maker, id, "assembly")
		milestone(maker, id, "quality")
		milestone(factoryQA, id, "released")
	}
	r = command("ship")
	r.ShipmentID = "asn-one"
	r.Units = []string{a1, a2}
	apply(maker, r)
	plan, err := supply.Plan(ctx, buyer, po.ID, "")
	if err != nil || plan.State != "in-production" {
		t.Fatal("partial shipment order state", plan.State, err)
	}
	r = command("receive")
	r.ShipmentID = "asn-one"
	r.Units = []string{a1}
	negative(maker, r, "factory receiving as destination")
	negative(wrong, r, "wrong receiving organization")
	apply(receiver, r)
	if hook != nil {
		before := snapshot()
		var approvalID, approvalHash, stockID string
		var stockVersion int64
		err = pool.QueryRow(ctx, `select q.approval_id,q.payload_sha256,m.stock_unit_id,i.version from procurement.serial_supply_quality q
   join procurement.serial_supply_manifest m using(tenant_id,production_unit_id)
   join inventory.stock_unit i on i.tenant_id=m.tenant_id and i.stock_unit_id=m.stock_unit_id
   where q.tenant_id=$1 and q.production_unit_id=$2 and q.stage='receipt' order by q.attempt desc limit 1`, tenant, a1).Scan(&approvalID, &approvalHash, &stockID, &stockVersion)
		if err != nil {
			t.Fatal(err)
		}
		if _, err = db.NewHumanApprovals(pool).Decide(ctx, qa, tenant, approvalID, "j2-store", approvalHash, true, "generic decision bypass", "supply:release", nil); err == nil {
			t.Fatal("generic approval bypass")
		}
		if err = ops.TransitionStockUnit(ctx, tenant, "j2-store", stockID, "quarantine", "available", stockVersion); err == nil {
			t.Fatal("generic stock release bypass")
		}
		if _, err = pool.Exec(ctx, `update procurement.purchase_order set total_minor_units=total_minor_units+1 where tenant_id=$1 and purchase_order_id=$2`, tenant, po.ID); err == nil {
			t.Fatal("bound purchase amount changed")
		}
		if snapshot() != before {
			t.Fatal("approval/stock/purchase bypass changed durable state")
		}
	}
	// Newly received stock is quarantine, not ATP.
	atp, err := db.NewInventoryControl(pool).AvailableToPromise(ctx, tenant, "j2-store", "j2-a", time.Now().UTC().Add(time.Hour))
	if err != nil || atp.AvailableInventory != 0 {
		t.Fatal("quarantine ATP", atp, err)
	}
	r = command("quality")
	r.UnitID = a1
	r.Reason = "fixture receipt accepted"
	self := receiver
	self.Permissions = map[string]struct{}{"supply:release": {}}
	before = snapshot()
	negative(self, r, "receiver self quality approval")
	if snapshot() != before {
		t.Fatal("receiver self decision leaked")
	}
	apply(qa, r)
	atp, err = db.NewInventoryControl(pool).AvailableToPromise(ctx, tenant, "j2-store", "j2-a", time.Now().UTC().Add(time.Hour))
	if err != nil || atp.AvailableInventory != 1 {
		t.Fatal("released ATP", atp, err)
	}
	r = command("ship")
	r.ShipmentID = "asn-two"
	r.Units = []string{b1}
	apply(maker, r)
	plan, err = supply.Plan(ctx, buyer, po.ID, "")
	if err != nil || plan.State != "shipped" {
		t.Fatal("all shipped order state", plan.State, err)
	}
	// Fail after state changes, pending approval and connected receipt were written.
	_, err = pool.Exec(ctx, `create function public.j2_outbox_fail() returns trigger language plpgsql as $$
 begin if new.aggregate_type='serial-supply' and new.payload->>'command_id'='receipt-rollback' then raise exception 'fixture outbox rollback';end if;return new;end $$;
 create trigger j2_outbox_fail before insert on platform.outbox_event for each row execute function public.j2_outbox_fail();`)
	if err != nil {
		t.Fatal(err)
	}
	r = command("receive")
	r.CommandID = "receipt-rollback"
	r.ShipmentID = "asn-one"
	r.Units = []string{a2}
	before = snapshot()
	negative(receiver, r, "outbox failure")
	if snapshot() != before {
		t.Fatal("outbox rollback changed supply/stock/approval history")
	}
	if _, err = pool.Exec(ctx, `drop trigger j2_outbox_fail on platform.outbox_event;drop function public.j2_outbox_fail();`); err != nil {
		t.Fatal(err)
	}
	// Concurrent same command: one new effect; other calls either replay or conflict.
	r = command("receive")
	r.ShipmentID = "asn-one"
	r.Units = []string{a2}
	parallel := r
	var wg sync.WaitGroup
	var mu sync.Mutex
	fresh, replays, conflicts := 0, 0, 0
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			out, err := supply.Apply(ctx, receiver, parallel)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				conflicts++
				return
			}
			if out.Replay {
				replays++
			} else {
				fresh++
			}
			if out.Actor != receiver.Subject || out.Version != version+1 {
				t.Errorf("bad concurrent receipt %+v", out)
			}
		}()
	}
	wg.Wait()
	if fresh != 1 || fresh+replays+conflicts != 8 {
		t.Fatalf("concurrency fresh/replay/conflict=%d/%d/%d", fresh, replays, conflicts)
	}
	version++
	// Recover committed evidence from a new pool; no hidden POST retry.
	otherPool, err := pgxpool.New(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	restarted, err := db.NewSerialSupply(otherPool)
	if err != nil {
		t.Fatal(err)
	}
	recovered, err := restarted.CommandReceipt(ctx, receiver, po.ID, parallel.CommandID)
	otherPool.Close()
	if err != nil || recovered.Actor != receiver.Subject || recovered.Version != version {
		t.Fatal("recovery", recovered, err)
	}
	_, wantHash, err := sc.Canonical(parallel)
	if err != nil || wantHash != recovered.RequestSHA256 {
		t.Fatal("recovery request hash", err)
	}
	again, err := supply.Apply(ctx, receiver, parallel)
	if err != nil || !again.Replay {
		t.Fatal("exact replay", again, err)
	}
	changed := parallel
	changed.EvidenceSHA256 = strings.Repeat("c", 64)
	negative(receiver, changed, "changed replay payload")
	secondReceiver := person("j2-second-receiver", "j2-store", "supply:receive")
	negative(secondReceiver, parallel, "changed replay actor")
	r = command("receive")
	r.ShipmentID = "asn-one"
	r.Units = []string{b1}
	negative(receiver, r, "wrong ASN membership")
	r = command("receive")
	r.ShipmentID = "asn-two"
	r.Units = []string{b1}
	apply(receiver, r)
	plan, err = supply.Plan(ctx, buyer, po.ID, "")
	if err != nil || plan.State != "received" || len(plan.Units) != 4 || len(plan.Lines) != 2 {
		t.Fatal("final physical receipt", plan, err)
	}
	r = command("quality-reject")
	r.UnitID = a2
	r.Reason = "fixture inspection failed"
	apply(qa, r)
	r = command("reinspect")
	r.UnitID = a2
	r.Reason = "new inspection"
	negative(receiver, r, "unchanged rejected evidence")
	r.EvidenceSHA256 = strings.Repeat("b", 64)
	apply(receiver, r)
	r = command("quality")
	r.UnitID = a2
	r.Reason = "corrected fixture evidence accepted"
	apply(qa, r)
	r = command("quality")
	r.UnitID = b1
	r.Reason = "fixture evidence accepted"
	apply(qa, r)
	atp, err = db.NewInventoryControl(pool).AvailableToPromise(ctx, tenant, "j2-store", "j2-a", time.Now().UTC().Add(time.Hour))
	if err != nil || atp.AvailableInventory != 2 {
		t.Fatal("final A ATP", atp, err)
	}
	bATP, err := db.NewInventoryControl(pool).AvailableToPromise(ctx, tenant, "j2-store", "j2-b", time.Now().UTC().Add(time.Hour))
	if err != nil || bATP.AvailableInventory != 1 {
		t.Fatal("final B ATP", bATP, err)
	}
	var manifested, receivedCount, available, poReceived int
	err = pool.QueryRow(ctx, `select (select count(*) from procurement.serial_supply_manifest where tenant_id=$1),(select count(*) from procurement.serial_supply_receipt where tenant_id=$1),(select count(*) from inventory.stock_unit where tenant_id=$1 and state='available'),(select count(*) from procurement.purchase_order where tenant_id=$1 and state='received')`, tenant).Scan(&manifested, &receivedCount, &available, &poReceived)
	if err != nil || manifested != 3 || receivedCount != 3 || available != 3 || poReceived != 1 {
		t.Fatal("counts", manifested, receivedCount, available, poReceived, err)
	}
	before = snapshot()
	if _, err = pool.Exec(ctx, `update procurement.serial_supply_line set quantity=4 where tenant_id=$1`, tenant); err == nil {
		t.Fatal("line evidence mutated")
	}
	if snapshot() != before {
		t.Fatal("immutability changed history")
	}
	t.Logf("SERIAL_SUPPLY_CONNECTED_PASS planned3/registered4/rejected1/shipped3/received3/available3; splitASN2; factory/receipt review; rollback/recovery; concurrency%dnew/%dreplay/%dconflict; version%d", fresh, replays, conflicts, version)
}
````

### FILE: `internal/platform/postgres/serial_supply_http_integration_test.go`

```yaml
block_id: "GO-CONNECTED-SERIAL-SUPPLY:file13:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "8e3b09156a55f54abd87e302274fc90fd71c23bc3ad0caee37627d53bb12b72b"
variables: []
secrets_allowed: false
```

````go
// AUTHORED loopback HTTP integration fixture; bearer identities are explicit
// fixture principals, not proof of JWT/IdP cryptography.
package postgres_test

import (
	"bytes"
	"context"
	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/platform/identity"
	db "elite.local/enterprise/internal/platform/postgres"
	sc "elite.local/enterprise/internal/serialsupply"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"
)

type supplyFixtureVerifier struct {
	mu         sync.Mutex
	next       int
	principals map[string]identity.Principal
}

func (v *supplyFixtureVerifier) token(p identity.Principal) string {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.next++
	token := fmt.Sprintf("supply-fixture-%d", v.next)
	v.principals[token] = p
	return token
}
func (v *supplyFixtureVerifier) Verify(_ context.Context, token string) (identity.Principal, error) {
	v.mu.Lock()
	defer v.mu.Unlock()
	p, ok := v.principals[token]
	if !ok {
		return p, identity.ErrUnauthenticated
	}
	return p, nil
}

type supplyHTTPReference struct {
	client        *http.Client
	base          string
	verifier      *supplyFixtureVerifier
	mu            sync.Mutex
	posts         map[string]int
	acceptedPosts map[string]int
	dropped       map[string]bool
	recoveries    int
}

func supplyFixtureOrg(p identity.Principal) string {
	values := []string{}
	for o := range p.Organizations {
		values = append(values, o)
	}
	sort.Strings(values)
	if len(values) == 0 {
		return ""
	}
	return values[0]
}
func supplyFixtureSurface(p identity.Principal) string {
	if p.Allowed("supply:factory-read") {
		return "factory"
	}
	return "franchise"
}
func (c *supplyHTTPReference) call(ctx context.Context, p identity.Principal, method, path string, query url.Values, input, out any) error {
	if query == nil {
		query = url.Values{}
	}
	query.Set("organization_id", supplyFixtureOrg(p))
	var body []byte
	var err error
	if input != nil {
		body, err = json.Marshal(input)
		if err != nil {
			return err
		}
	}
	req, err := http.NewRequestWithContext(ctx, method, c.base+path+"?"+query.Encode(), bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.verifier.token(p))
	req.Header.Set("Content-Type", "application/json")
	response, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(response.Body, 1048577))
	if err != nil {
		return err
	}
	if len(raw) > 1048576 {
		return errors.New("oversized fixture response")
	}
	if response.Header.Get("Cache-Control") != "no-store" {
		return errors.New("private response not no-store")
	}
	if response.StatusCode >= 400 {
		return fmt.Errorf("HTTP %d: %s", response.StatusCode, raw)
	}
	return json.Unmarshal(raw, out)
}
func (c *supplyHTTPReference) BindPlan(ctx context.Context, p identity.Principal, r sc.PlanRequest) (sc.Receipt, error) {
	var out sc.Receipt
	err := c.call(ctx, p, "POST", "/v1/franchise/supply/orders/"+url.PathEscape(r.PurchaseOrderID)+"/plan", nil, r, &out)
	return out, err
}
func (c *supplyHTTPReference) Plan(ctx context.Context, p identity.Principal, id, after string) (sc.Plan, error) {
	var out sc.Plan
	q := url.Values{}
	if after != "" {
		q.Set("after_unit", after)
	}
	err := c.call(ctx, p, "GET", "/v1/"+supplyFixtureSurface(p)+"/supply/orders/"+url.PathEscape(id), q, nil, &out)
	return out, err
}
func (c *supplyHTTPReference) CommandReceipt(ctx context.Context, p identity.Principal, id, key string) (sc.Receipt, error) {
	var out sc.Receipt
	err := c.call(ctx, p, "GET", "/v1/"+supplyFixtureSurface(p)+"/supply/orders/"+url.PathEscape(id), url.Values{"command_id": {key}}, nil, &out)
	return out, err
}
func (c *supplyHTTPReference) Apply(ctx context.Context, p identity.Principal, r sc.Command) (sc.Receipt, error) {
	surface := "franchise"
	switch r.Kind {
	case "confirm", "start", "register", "milestone", "ship":
		surface = "factory"
	}
	var out sc.Receipt
	err := c.call(ctx, p, "POST", "/v1/"+surface+"/supply/orders/"+url.PathEscape(r.PurchaseOrderID)+"/"+r.Kind, nil, r, &out)
	if err == nil {
		return out, nil
	}
	// Only transport loss triggers fixture recovery. A rejected HTTP command is
	// never silently retried, and no second POST is made here.
	var network *url.Error
	if !errors.As(err, &network) {
		return out, err
	}
	recovered, getErr := c.CommandReceipt(ctx, p, r.PurchaseOrderID, r.CommandID)
	if getErr != nil {
		return out, fmt.Errorf("unconfirmed POST %v; recovery %v", err, getErr)
	}
	_, hash, hashErr := sc.Canonical(r)
	if hashErr != nil || recovered.PurchaseOrderID != r.PurchaseOrderID || recovered.CommandID != r.CommandID || recovered.Actor != p.Subject || recovered.RequestSHA256 != hash {
		return out, errors.New("recovery evidence mismatch")
	}
	c.mu.Lock()
	c.recoveries++
	c.mu.Unlock()
	return recovered, nil
}
func TestSerialSupplyHTTPConnectedReference(t *testing.T) {
	var reference *supplyHTTPReference
	testSerialSupplyConnected(t, func(t *testing.T, _ *pgxpool.Pool, store *db.SerialSupply) serialSupplyReference {
		verifier := &supplyFixtureVerifier{principals: map[string]identity.Principal{}}
		mux := http.NewServeMux()
		httpapi.SerialSupplyModule{Service: store}.Register(mux, verifier)
		reference = &supplyHTTPReference{verifier: verifier, posts: map[string]int{}, acceptedPosts: map[string]int{}, dropped: map[string]bool{}}
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var command sc.Command
			if r.Method == "POST" {
				raw, err := io.ReadAll(io.LimitReader(r.Body, 32769))
				if err != nil {
					t.Error(err)
				}
				r.Body.Close()
				r.Body = io.NopCloser(bytes.NewReader(raw))
				_ = json.Unmarshal(raw, &command)
				reference.mu.Lock()
				reference.posts[command.CommandID]++
				reference.mu.Unlock()
			}
			result := httptest.NewRecorder()
			mux.ServeHTTP(result, r)
			if r.Method == "POST" && (result.Code == 200 || result.Code == 201) {
				reference.mu.Lock()
				reference.acceptedPosts[command.CommandID]++
				reference.mu.Unlock()
			}
			drop := ""
			if r.Method == "POST" && result.Code == 201 {
				if command.Kind == "receive" && command.ShipmentID == "asn-one" {
					drop = "first-receipt"
				}
				if command.Kind == "ship" && command.ShipmentID == "asn-two" {
					drop = "last-shipment"
				}
			}
			reference.mu.Lock()
			lose := drop != "" && !reference.dropped[drop]
			if lose {
				reference.dropped[drop] = true
				reference.dropped["command:"+command.CommandID] = true
			}
			reference.mu.Unlock()
			if lose {
				connection, _, err := w.(http.Hijacker).Hijack()
				if err != nil {
					t.Error(err)
					return
				}
				connection.Close()
				return
			}
			for key, values := range result.Header() {
				for _, value := range values {
					w.Header().Add(key, value)
				}
			}
			w.WriteHeader(result.Code)
			_, _ = w.Write(result.Body.Bytes())
		}))
		transport := &http.Transport{DisableKeepAlives: true}
		reference.client = &http.Client{Transport: transport, Timeout: 5 * time.Second}
		reference.base = server.URL
		t.Cleanup(func() { transport.CloseIdleConnections(); server.Close() })
		return reference
	})
	reference.mu.Lock()
	defer reference.mu.Unlock()
	if reference.recoveries != 2 || !reference.dropped["first-receipt"] || !reference.dropped["last-shipment"] {
		t.Fatal("lost-response coverage", reference.recoveries, reference.dropped)
	}
	for key := range reference.dropped {
		if strings.HasPrefix(key, "command:") && reference.acceptedPosts[strings.TrimPrefix(key, "command:")] != 1 {
			t.Fatal("hidden POST retry", key, reference.posts)
		}
	}
	t.Log("SERIAL_SUPPLY_HTTP_PASS 30-version connected journey; two committed responses lost and recovered by actor/hash GET; generic approval/stock/purchase bypass refused")
}
````

### FILE: `internal/platform/postgres/serial_supply_owner_tx.go`

```yaml
block_id: "GO-CONNECTED-SERIAL-SUPPLY:file14:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "463cc14be88f99ceee80d1366da10e4b7e1300d581a9947d29c58753bc345960"
variables: []
secrets_allowed: false
```

````go
// AUTHORED transaction adapter. Existing operations.Service performs all preparation
// and transition validation; existing PostgreSQL owner helpers perform the writes.
package postgres

import (
	"context"
	"elite.local/enterprise/internal/operations"
	"github.com/jackc/pgx/v5"
)

type operationsTxRepository struct{ tx pgx.Tx }

var _ operations.Repository = operationsTxRepository{}

func (r operationsTxRepository) CreatePurchaseOrder(ctx context.Context, tenant, eventID string, value operations.PurchaseOrder) error {
	return createPurchaseOrderInTx(ctx, r.tx, tenant, eventID, value)
}
func (r operationsTxRepository) TransitionPurchaseOrder(ctx context.Context, tenant, organization, id, current string, version int64, target, eventID string) error {
	return transitionPurchaseOrderInTx(ctx, r.tx, tenant, organization, id, current, version, target, eventID)
}
func (r operationsTxRepository) CreateProductionUnit(ctx context.Context, tenant, eventID string, value operations.ProductionUnit) error {
	return createProductionUnitInTx(ctx, r.tx, tenant, eventID, value)
}
func (r operationsTxRepository) TransitionProductionUnit(ctx context.Context, tenant, organization, id, current, target, eventID string) error {
	return transitionProductionUnitInTx(ctx, r.tx, tenant, organization, id, current, target, eventID)
}
func (r operationsTxRepository) CreateStockUnit(ctx context.Context, tenant, eventID string, value operations.StockUnit) error {
	return createStockUnitInTx(ctx, r.tx, tenant, eventID, value)
}
func (r operationsTxRepository) TransitionStockUnit(ctx context.Context, tenant, organization, id, current string, version int64, target, eventID string) error {
	return transitionStockUnitInTx(ctx, r.tx, tenant, organization, id, current, version, target, eventID)
}
````

### FILE: `internal/platform/postgres/serial_supply_quality.go`

```yaml
block_id: "GO-CONNECTED-SERIAL-SUPPLY:file15:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "dd2dab39216499780c69c9442ff7f80cbe48d4161f097c422ba2c5ac3afdb62c"
variables: []
secrets_allowed: false
```

````go
// AUTHORED quality payload binding. Decisions use the existing shared registry.
package postgres

import (
	"context"
	"elite.local/enterprise/internal/approval"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/randomid"
	sc "elite.local/enterprise/internal/serialsupply"
	"github.com/jackc/pgx/v5"
	"strconv"
)

type supplyQualityState struct {
	ID, Hash, State, SourceState, Evidence string
	Version                                int64
	Attempt                                int
}

func readSupplyQuality(ctx context.Context, tx pgx.Tx, tenant, unit, stage string) (supplyQualityState, error) {
	var q supplyQualityState
	err := tx.QueryRow(ctx, `select q.approval_id,q.payload_sha256,a.state,q.source_state,q.source_version,q.attempt,coalesce(a.payload->>'basis_evidence_sha256','')
 from procurement.serial_supply_quality q join approval.request a on a.tenant_id=q.tenant_id and a.request_id=q.approval_id
 where q.tenant_id=$1 and q.production_unit_id=$2 and q.stage=$3 and a.kind='serial_quality' and a.evidence_sha=q.payload_sha256
 order by q.attempt desc limit 1 for update of a`, tenant, unit, stage).Scan(&q.ID, &q.Hash, &q.State, &q.SourceState, &q.Version, &q.Attempt, &q.Evidence)
	return q, supplyError(err)
}
func (s *SerialSupply) proposeSupplyQuality(ctx context.Context, tx pgx.Tx, p identity.Principal, plan sc.Plan, r sc.Command, u supplyUnitState, stage string) (string, error) {
	org, permission := supplyCommandPermission(r, plan)
	state, version := u.State, int64(0)
	if stage == "receipt" {
		state = u.StockState
		version = u.StockVersion
	}
	raw, hash, err := sc.Canonical(map[string]any{"purchase_order_id": plan.PurchaseOrderID, "unit_id": u.ID, "stock_unit_id": u.StockID, "serial_number": u.Serial, "variant_id": u.VariantID, "stage": stage, "source_state": state, "source_version": strconv.FormatInt(version, 10), "policy_code": sc.PolicyCode, "organization_id": org, "basis_evidence_sha256": r.EvidenceSHA256, "requester": p.Subject})
	if err != nil {
		return "", err
	}
	id := randomid.Generator{}.New()
	spec := HumanApprovalSpec{Request: approval.Request{TenantID: p.TenantID, ID: id, Kind: approval.KindSerialQuality, SubjectID: u.ID, Requester: p.Subject, EvidenceSHA: hash}, OrganizationID: org, Payload: raw}
	if _, err = NewHumanApprovals(s.pool).submitTx(ctx, tx, p, spec, permission, nil); err != nil {
		return "", err
	}
	_, err = tx.Exec(ctx, `insert into procurement.serial_supply_quality(tenant_id,purchase_order_id,production_unit_id,stage,attempt,approval_id,command_id,source_state,source_version,payload_sha256)
 select $1,$2,$3,$4,coalesce(max(attempt),0)+1,$5,$6,$7,$8,$9 from procurement.serial_supply_quality where tenant_id=$1 and production_unit_id=$3 and stage=$4`, p.TenantID, plan.PurchaseOrderID, u.ID, stage, id, r.CommandID, state, version, hash)
	return id, supplyError(err)
}
func (s *SerialSupply) decideSupplyQuality(ctx context.Context, tx pgx.Tx, p identity.Principal, plan sc.Plan, r sc.Command, u supplyUnitState, stage string, approved bool) (string, error) {
	q, err := readSupplyQuality(ctx, tx, p.TenantID, u.ID, stage)
	if err != nil {
		return "", err
	}
	state, version := u.State, int64(0)
	if stage == "receipt" {
		state = u.StockState
		version = u.StockVersion
	}
	if q.State != "pending" || q.SourceState != state || q.Version != version {
		return "", sc.ErrConflict
	}
	org, permission := supplyCommandPermission(r, plan)
	if _, err = NewHumanApprovals(s.pool).decideTx(ctx, tx, p, p.TenantID, q.ID, org, q.Hash, approved, r.Reason, permission, nil); err != nil {
		return "", err
	}
	return q.ID, nil
}
func (s *SerialSupply) supplyMilestone(ctx context.Context, tx pgx.Tx, p identity.Principal, plan sc.Plan, r sc.Command) (any, error) {
	if plan.State != "in-production" {
		return nil, sc.ErrConflict
	}
	u, err := readSupplyUnit(ctx, tx, p.TenantID, plan.PurchaseOrderID, r.UnitID)
	if err != nil {
		return nil, err
	}
	var approvalID string
	if u.State == "quality" && (r.TargetState == "released" || r.TargetState == "rejected") {
		approvalID, err = s.decideSupplyQuality(ctx, tx, p, plan, r, u, "factory", r.TargetState == "released")
		if err != nil {
			return nil, err
		}
	}
	if err = moveSupplyFactory(ctx, tx, p.TenantID, plan, r, u.ID, u.State, r.TargetState); err != nil {
		return nil, err
	}
	if r.TargetState == "quality" {
		u.State = "quality"
		approvalID, err = s.proposeSupplyQuality(ctx, tx, p, plan, r, u, "factory")
		if err != nil {
			return nil, err
		}
	}
	return map[string]any{"unit_id": u.ID, "to_state": r.TargetState, "approval_id": approvalID}, nil
}
func (s *SerialSupply) supplyReceiptQuality(ctx context.Context, tx pgx.Tx, p identity.Principal, plan sc.Plan, r sc.Command) (any, error) {
	u, err := readSupplyUnit(ctx, tx, p.TenantID, plan.PurchaseOrderID, r.UnitID)
	if err != nil {
		return nil, err
	}
	if u.State != "received" || u.StockState != "quarantine" {
		return nil, sc.ErrConflict
	}
	if r.Kind == "reinspect" {
		q, err := readSupplyQuality(ctx, tx, p.TenantID, u.ID, "receipt")
		if err != nil {
			return nil, err
		}
		if q.State != "rejected" || q.Evidence == r.EvidenceSHA256 {
			return nil, sc.ErrConflict
		}
		id, err := s.proposeSupplyQuality(ctx, tx, p, plan, r, u, "receipt")
		return map[string]any{"unit_id": u.ID, "stock_unit_id": u.StockID, "approval_id": id, "stock_state": "quarantine", "inspection_reason": r.Reason}, err
	}
	approved := r.Kind == "quality"
	id, err := s.decideSupplyQuality(ctx, tx, p, plan, r, u, "receipt", approved)
	if err != nil {
		return nil, err
	}
	target := "quarantine"
	if approved {
		target = "available"
		if err = moveSupplyStock(ctx, tx, p.TenantID, plan, r, u.StockID, u.StockState, target, u.StockVersion); err != nil {
			return nil, err
		}
	}
	return map[string]any{"unit_id": u.ID, "stock_unit_id": u.StockID, "approval_id": id, "approved": approved, "stock_state": target}, nil
}
````

### FILE: `internal/platform/postgres/serial_supply_shipping.go`

```yaml
block_id: "GO-CONNECTED-SERIAL-SUPPLY:file16:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "ed82d857d6a7a0807de6c20f630db2a94e8838b980f03e795cb85e33bf50611d"
variables: []
secrets_allowed: false
```

````go
// AUTHORED serial ASN/partial receiving composition over existing state writers.
package postgres

import (
	"context"
	"elite.local/enterprise/internal/operations"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/randomid"
	sc "elite.local/enterprise/internal/serialsupply"
	"github.com/jackc/pgx/v5"
	"sort"
)

func (s *SerialSupply) supplyShip(ctx context.Context, tx pgx.Tx, p identity.Principal, plan sc.Plan, r sc.Command) (any, error) {
	if plan.State != "in-production" {
		return nil, sc.ErrConflict
	}
	_, err := tx.Exec(ctx, `insert into procurement.serial_supply_shipment(tenant_id,purchase_order_id,shipment_id,command_id) values($1,$2,$3,$4)`, p.TenantID, plan.PurchaseOrderID, r.ShipmentID, r.CommandID)
	if err != nil {
		return nil, supplyError(err)
	}
	ids := append([]string(nil), r.Units...)
	sort.Strings(ids)
	lines := []map[string]string{}
	svc := operations.NewService(operationsTxRepository{tx: tx}, randomid.Generator{})
	for _, id := range ids {
		u, err := readSupplyUnit(ctx, tx, p.TenantID, plan.PurchaseOrderID, id)
		if err != nil {
			return nil, err
		}
		if u.State != "released" || u.StockID != "" || u.ShipmentID != "" {
			return nil, sc.ErrConflict
		}
		q, err := readSupplyQuality(ctx, tx, p.TenantID, u.ID, "factory")
		if err != nil {
			return nil, err
		}
		if q.State != "approved" || q.SourceState != "quality" {
			return nil, sc.ErrConflict
		}
		if err = moveSupplyFactory(ctx, tx, p.TenantID, plan, r, u.ID, u.State, "shipped"); err != nil {
			return nil, err
		}
		stock, err := svc.ReceiveStockUnit(ctx, p.TenantID, operations.StockUnit{OrganizationID: plan.DestinationOrganizationID, VariantID: u.VariantID, ProductionUnitID: u.ID, SerialNumber: u.Serial, VIN: u.VIN, BatterySerialNumber: u.Battery})
		if err != nil {
			return nil, supplyError(err)
		}
		_, err = tx.Exec(ctx, `insert into procurement.serial_supply_manifest(tenant_id,purchase_order_id,shipment_id,production_unit_id,stock_unit_id) values($1,$2,$3,$4,$5)`, p.TenantID, plan.PurchaseOrderID, r.ShipmentID, u.ID, stock.ID)
		if err != nil {
			return nil, supplyError(err)
		}
		if err = supplyEffect(ctx, tx, p.TenantID, plan.PurchaseOrderID, r.CommandID, "stock", stock.ID, "", "in-transit", 0, 1); err != nil {
			return nil, err
		}
		lines = append(lines, map[string]string{"unit_id": u.ID, "stock_unit_id": stock.ID, "serial_number": u.Serial, "variant_id": u.VariantID, "line_id": u.LineID})
	}
	var complete bool
	err = tx.QueryRow(ctx, `select not exists(select 1 from procurement.serial_supply_line l where tenant_id=$1 and purchase_order_id=$2
 and l.quantity<>(select count(*) from procurement.serial_supply_unit b join procurement.serial_supply_manifest m using(tenant_id,production_unit_id)
 where b.tenant_id=l.tenant_id and b.purchase_order_id=l.purchase_order_id and b.line_id=l.line_id))`, p.TenantID, plan.PurchaseOrderID).Scan(&complete)
	if err != nil {
		return nil, err
	}
	if complete {
		if err = moveSupplyPurchase(ctx, tx, p.TenantID, plan, r, "shipped"); err != nil {
			return nil, err
		}
	}
	return map[string]any{"shipment_id": r.ShipmentID, "manifest": lines, "all_declared_units_shipped": complete}, nil
}
func (s *SerialSupply) supplyReceive(ctx context.Context, tx pgx.Tx, p identity.Principal, plan sc.Plan, r sc.Command) (any, error) {
	if plan.State != "in-production" && plan.State != "shipped" {
		return nil, sc.ErrConflict
	}
	ids := append([]string(nil), r.Units...)
	sort.Strings(ids)
	received := []map[string]string{}
	for _, id := range ids {
		u, err := readSupplyUnit(ctx, tx, p.TenantID, plan.PurchaseOrderID, id)
		if err != nil {
			return nil, err
		}
		if u.State != "shipped" || u.StockState != "in-transit" || u.ShipmentID != r.ShipmentID || u.StockID == "" {
			return nil, sc.ErrConflict
		}
		_, err = tx.Exec(ctx, `insert into procurement.serial_supply_receipt(tenant_id,purchase_order_id,shipment_id,production_unit_id,command_id) values($1,$2,$3,$4,$5)`, p.TenantID, plan.PurchaseOrderID, r.ShipmentID, id, r.CommandID)
		if err != nil {
			return nil, supplyError(err)
		}
		if err = moveSupplyFactory(ctx, tx, p.TenantID, plan, r, id, u.State, "received"); err != nil {
			return nil, err
		}
		if err = moveSupplyStock(ctx, tx, p.TenantID, plan, r, u.StockID, u.StockState, "quarantine", u.StockVersion); err != nil {
			return nil, err
		}
		u.State = "received"
		u.StockState = "quarantine"
		u.StockVersion++
		approvalID, err := s.proposeSupplyQuality(ctx, tx, p, plan, r, u, "receipt")
		if err != nil {
			return nil, err
		}
		received = append(received, map[string]string{"unit_id": id, "stock_unit_id": u.StockID, "approval_id": approvalID, "state": "quarantine"})
	}
	var complete bool
	err := tx.QueryRow(ctx, `select not exists(select 1 from procurement.serial_supply_line l where tenant_id=$1 and purchase_order_id=$2
 and l.quantity<>(select count(*) from procurement.serial_supply_unit b join procurement.serial_supply_receipt r using(tenant_id,production_unit_id)
 where b.tenant_id=l.tenant_id and b.purchase_order_id=l.purchase_order_id and b.line_id=l.line_id))`, p.TenantID, plan.PurchaseOrderID).Scan(&complete)
	if err != nil {
		return nil, err
	}
	if complete {
		if plan.State != "shipped" {
			return nil, sc.ErrConflict
		}
		if err = moveSupplyPurchase(ctx, tx, p.TenantID, plan, r, "received"); err != nil {
			return nil, err
		}
	}
	return map[string]any{"shipment_id": r.ShipmentID, "received": received, "all_declared_units_received": complete, "payment_created": false}, nil
}
func (s *SerialSupply) Apply(ctx context.Context, p identity.Principal, r sc.Command) (sc.Receipt, error) {
	if !r.Valid() {
		return sc.Receipt{}, sc.ErrInvalid
	}
	var action supplyAction
	switch r.Kind {
	case "submit", "confirm", "start", "cancel":
		action = func(ctx context.Context, tx pgx.Tx, p identity.Principal, plan sc.Plan, r sc.Command) (any, error) {
			return supplyLifecycle(ctx, tx, p.TenantID, plan, r)
		}
	case "register":
		action = func(ctx context.Context, tx pgx.Tx, p identity.Principal, plan sc.Plan, r sc.Command) (any, error) {
			return registerSupplyUnit(ctx, tx, p.TenantID, plan, r)
		}
	case "milestone":
		action = s.supplyMilestone
	case "ship":
		action = s.supplyShip
	case "receive":
		action = s.supplyReceive
	case "quality", "quality-reject", "reinspect":
		action = s.supplyReceiptQuality
	default:
		return sc.Receipt{}, sc.ErrInvalid
	}
	return s.executeSupplyCommand(ctx, p, r, action)
}
````

### FILE: `internal/serialsupply/contract.go`

```yaml
block_id: "GO-CONNECTED-SERIAL-SUPPLY:file17:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "455d496ef0c4861e94f93e3ad7de18a29cc6b5aded7aefe25fe2aa9be8a770de"
variables: []
secrets_allowed: false
```

````go
// AUTHORED bounded contracts for connecting existing procurement, factory and
// inventory owners. No new pricing, stock ledger or corporate implementation.
package serialsupply

import (
	"encoding/json"
	"errors"
	"regexp"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"elite.local/enterprise/internal/approval"
)

var ErrInvalid = errors.New("invalid serial supply request")
var ErrNotFound = errors.New("serial supply not found")
var ErrConflict = errors.New("serial supply conflict")
var idPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:-]{0,127}$`)
var shaPattern = regexp.MustCompile(`^[0-9a-f]{64}$`)

const PolicyCode = "strict-serial-reference/v1"
const PolicyJSON = `{"schema":"elite-serial-supply-policy/v1","scope":"LIBRARY_INFRASTRUCTURE_REFERENCE","quantity":"whole-serialized-units","overreceipt":0,"receipt_state":"quarantine","release":"distinct-human-with-evidence","pricing":"existing-purchase-order-total-unchanged","production_authorized":false}`

func ValidID(s string) bool  { return idPattern.MatchString(s) }
func ValidSHA(s string) bool { return shaPattern.MatchString(s) }
func ValidText(s string, max int) bool {
	if s == "" || len(s) > max || !utf8.ValidString(s) || strings.TrimSpace(s) != s {
		return false
	}
	for _, r := range s {
		if unicode.IsControl(r) {
			return false
		}
	}
	return true
}
func Canonical(value any) (json.RawMessage, string, error) {
	raw, err := json.Marshal(value)
	if err != nil {
		return nil, "", err
	}
	out, hash, err := approval.CanonicalPayload(raw)
	return out, hash, err
}

type Line struct {
	ID        string `json:"id"`
	VariantID string `json:"variant_id"`
	Quantity  int    `json:"quantity"`
}
type PlanRequest struct {
	PurchaseOrderID       string `json:"purchase_order_id"`
	CommandID             string `json:"command_id"`
	FactoryOrganizationID string `json:"factory_organization_id"`
	DemandReference       string `json:"demand_reference"`
	PolicyCode            string `json:"policy_code"`
	EvidenceSHA256        string `json:"evidence_sha256"`
	Lines                 []Line `json:"lines"`
}

func (r PlanRequest) Valid() bool {
	if !ValidID(r.PurchaseOrderID) || !ValidID(r.CommandID) || !ValidID(r.FactoryOrganizationID) || !ValidText(r.DemandReference, 1024) || r.PolicyCode != PolicyCode || !ValidSHA(r.EvidenceSHA256) || len(r.Lines) < 1 || len(r.Lines) > 32 {
		return false
	}
	ids := map[string]bool{}
	variants := map[string]bool{}
	total := 0
	for _, l := range r.Lines {
		if !ValidID(l.ID) || !ValidID(l.VariantID) || l.Quantity < 1 || l.Quantity > 1000 || ids[l.ID] || variants[l.VariantID] {
			return false
		}
		ids[l.ID] = true
		variants[l.VariantID] = true
		total += l.Quantity
	}
	return total <= 1000
}

type Command struct {
	PurchaseOrderID     string   `json:"purchase_order_id"`
	CommandID           string   `json:"command_id"`
	ExpectedVersion     int64    `json:"expected_version,string"`
	Kind                string   `json:"kind"`
	EvidenceSHA256      string   `json:"evidence_sha256"`
	LineID              string   `json:"line_id,omitempty"`
	UnitID              string   `json:"unit_id,omitempty"`
	SerialNumber        string   `json:"serial_number,omitempty"`
	VIN                 string   `json:"vin,omitempty"`
	BatterySerialNumber string   `json:"battery_serial_number,omitempty"`
	TargetState         string   `json:"target_state,omitempty"`
	ShipmentID          string   `json:"shipment_id,omitempty"`
	Units               []string `json:"units,omitempty"`
	Reason              string   `json:"reason,omitempty"`
}

func (r Command) Valid() bool {
	if !ValidID(r.PurchaseOrderID) || !ValidID(r.CommandID) || r.ExpectedVersion < 1 || !ValidSHA(r.EvidenceSHA256) {
		return false
	}
	allowed := map[string]bool{}
	switch r.Kind {
	case "submit", "confirm", "start", "cancel":
	case "register":
		if !ValidID(r.LineID) || !ValidText(r.SerialNumber, 128) || r.VIN != "" && !ValidText(r.VIN, 128) || r.BatterySerialNumber != "" && !ValidText(r.BatterySerialNumber, 128) {
			return false
		}
		allowed["line_id"] = true
		allowed["serial_number"] = true
		allowed["vin"] = true
		allowed["battery_serial_number"] = true
	case "milestone":
		if !ValidID(r.UnitID) || !map[string]bool{"assembly": true, "quality": true, "released": true, "rejected": true}[r.TargetState] {
			return false
		}
		allowed["unit_id"] = true
		allowed["target_state"] = true
		if r.TargetState == "released" || r.TargetState == "rejected" {
			allowed["reason"] = true
			if !ValidText(r.Reason, 2000) {
				return false
			}
		}
	case "ship", "receive":
		if !ValidID(r.ShipmentID) || len(r.Units) < 1 || len(r.Units) > 100 {
			return false
		}
		seen := map[string]bool{}
		for _, u := range r.Units {
			if !ValidID(u) || seen[u] {
				return false
			}
			seen[u] = true
		}
		allowed["shipment_id"] = true
		allowed["units"] = true
	case "quality", "quality-reject", "reinspect":
		if !ValidID(r.UnitID) {
			return false
		}
		allowed["unit_id"] = true
		allowed["reason"] = true
		if !ValidText(r.Reason, 2000) {
			return false
		}
	default:
		return false
	}
	values := map[string]string{"line_id": r.LineID, "unit_id": r.UnitID, "serial_number": r.SerialNumber, "vin": r.VIN, "battery_serial_number": r.BatterySerialNumber, "target_state": r.TargetState, "shipment_id": r.ShipmentID, "reason": r.Reason}
	for k, v := range values {
		if v != "" && !allowed[k] {
			return false
		}
	}
	return len(r.Units) == 0 || allowed["units"]
}

type Receipt struct {
	PurchaseOrderID string          `json:"purchase_order_id"`
	CommandID       string          `json:"command_id"`
	Version         int64           `json:"version,string"`
	Kind            string          `json:"kind"`
	Actor           string          `json:"actor"`
	RequestSHA256   string          `json:"request_sha256"`
	PayloadSHA256   string          `json:"payload_sha256"`
	Payload         json.RawMessage `json:"payload"`
	RecordedAt      time.Time       `json:"recorded_at"`
	Replay          bool            `json:"replay"`
}
type Plan struct {
	PurchaseOrderID           string   `json:"purchase_order_id"`
	DestinationOrganizationID string   `json:"destination_organization_id"`
	FactoryOrganizationID     string   `json:"factory_organization_id"`
	SupplierID                string   `json:"supplier_id"`
	State                     string   `json:"state"`
	PurchaseVersion           int64    `json:"purchase_version,string"`
	Currency                  string   `json:"currency"`
	TotalMinorUnits           int64    `json:"total_minor_units,string"`
	Version                   int64    `json:"version,string"`
	PolicyCode                string   `json:"policy_code"`
	DemandReference           string   `json:"demand_reference"`
	Lines                     []Line   `json:"lines"`
	Units                     []Unit   `json:"units"`
	Latest                    *Receipt `json:"latest,omitempty"`
	NextUnitID                string   `json:"next_unit_id,omitempty"`
}
type Unit struct {
	FactoryReviewState string `json:"factory_review_state,omitempty"`
	FactoryRequester   string `json:"factory_requester,omitempty"`
	ReceiptReviewState string `json:"receipt_review_state,omitempty"`
	ReceiptRequester   string `json:"receipt_requester,omitempty"`
	ID                 string `json:"id"`
	LineID             string `json:"line_id"`
	State              string `json:"state"`
	SerialNumber       string `json:"serial_number"`
	StockUnitID        string `json:"stock_unit_id,omitempty"`
	StockState         string `json:"stock_state,omitempty"`
	ShipmentID         string `json:"shipment_id,omitempty"`
}
````

### FILE: `internal/serialsupply/contract_fuzz_test.go`

```yaml
block_id: "GO-CONNECTED-SERIAL-SUPPLY:file18:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "2f5767c74fc5fb05d625c952718508a3069e2fd45d15236c10d6db6322f82042"
variables: []
secrets_allowed: false
```

````go
package serialsupply

// AUTHORED finite-fuzz invariants for physical serial identity, receipt hashes
// and immutable command versions; not a deployment security certification.
import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func FuzzSerialSupplyCommandIdentity(f *testing.F) {
	for _, s := range []string{"FRAME-0001", "BATERÍA-ñ-2", " serial ", "a\x00b", "\xff", "電池-7", strings.Repeat("A", 129)} {
		f.Add(s, int64(9007199254740993))
	}
	f.Fuzz(func(t *testing.T, serial string, version int64) {
		r := Command{PurchaseOrderID: "po", CommandID: "register", Kind: "register", LineID: "line", SerialNumber: serial, ExpectedVersion: version, EvidenceSHA256: strings.Repeat("a", 64)}
		if !r.Valid() {
			return
		}
		raw, hash, err := Canonical(r)
		if err != nil {
			t.Fatal(err)
		}
		var decoded Command
		if err = json.Unmarshal(raw, &decoded); err != nil || decoded.SerialNumber != serial || decoded.ExpectedVersion != version || !decoded.Valid() {
			t.Fatal("serial/version identity lost", err)
		}
		again, other, err := Canonical(decoded)
		if err != nil || hash != other || !bytes.Equal(raw, again) {
			t.Fatal("unstable durable command")
		}
		decoded.CommandID += "-changed"
		_, different, err := Canonical(decoded)
		if err != nil || different == hash {
			t.Fatal("command identity not hash-bound")
		}
		shipment := Command{PurchaseOrderID: "po", CommandID: "receive", Kind: "receive", ShipmentID: "asn", Units: []string{"unit", "unit"}, ExpectedVersion: version, EvidenceSHA256: r.EvidenceSHA256}
		if shipment.Valid() {
			t.Fatal("duplicate physical unit admitted")
		}
	})
}
````

## 6. Configuration surface

docs/SERIAL_SUPPLY_REFERENCE.md documents built-in policy hash, optional host activation, permissions, paths, bounded bodies, exact quoted int64 and recovery. deploy/serial-supply/plan.reference.json is an explicit synthetic template. Disabled host reads no policy; enabled host requires exact hash and15active guards. No accounts requested.

## 7. Dependency bill

Zero new dependencies or upstream code. Existing admitted Operations, shared Approval and BC-derived Inventory remain selected with notices/locks. This pack is AUTHORED configuration and transactional/interface glue; source-derived stock algorithms are not reimplemented.

## 8. Apply order

MARKDOWN-COMPOSITOR0.3.0 with complete dependency closure, absent destination. Apply migration72 after71. Existing public owner wrappers remain. Populated down refuses immutable history; empty down/up tested. Do not rewrite old migrations.

## 9. Verification

Actual PG69migrations/Go1.26.8: three planned, one reject/replacement, two ASNs, three received and approved available, version30; concurrency1new/7replay;15table rollback snapshots; two lost responses recovered by GET; generic bypass rejected;13transport boundaries; activation hash/guards; empty/populated down and finite domain fuzz. Separate T2803 composition SCA, T2804 role UI and T2810 release remain open.

## 10. Reconstruction evidence

Four exact affected profile reconstructions required before canonical publication. SERIAL_SUPPLY_CONNECTED_RELEASE_V402.md/json binds current source and receipts. Backend-only earlier evidence is preserved, not rewritten.


V402 composed delta: T2804 connected supply role: original Operations/BindPlan composed atomically; bounded forms, latest review/requester projection, nav and GET recovery. No migration/dependency/state-machine change. SUPPLY_ROLE_RELEASE_V402.md.
