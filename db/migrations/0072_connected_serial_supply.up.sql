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
