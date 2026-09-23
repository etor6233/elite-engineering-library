begin;
do $$ begin
 if exists(select 1 from approval.request where kind='whatsapp_schedule') or exists(select 1 from communication.whatsapp_schedule)
 then raise exception 'cannot remove scheduled WhatsApp evidence';end if;
end $$;
create or replace view communication.whatsapp_delivery_approval as
 select g.tenant_id,g.channel_code,g.delivery_key,g.organization_id,g.appointment_id,g.confirmation_event_id::text as approval_event_id,g.lead_id,g.external_id_hmac,g.profile_sha256,g.message_sha256,g.expires_at
 from communication.whatsapp_appointment_approval g
 union all
 select r.tenant_id,'whatsapp',r.request_id,r.organization_id,''::text,r.request_id,
 r.payload->>'lead_id',r.payload->>'external_id_hmac',r.payload->>'profile_sha256',r.payload->>'message_sha256',(r.payload->>'expires_at')::timestamptz
 from approval.request r
 where r.kind='whatsapp_reply' and r.state='approved' and r.payload->>'schema'='elite-whatsapp-reply-approval/v1'
 and r.payload->>'organization_id'=r.organization_id;
drop table communication.whatsapp_schedule_result;
drop table communication.whatsapp_schedule_cancellation;
drop table communication.whatsapp_schedule;
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

commit;
