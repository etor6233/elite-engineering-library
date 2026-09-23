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
