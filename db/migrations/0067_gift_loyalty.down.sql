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
