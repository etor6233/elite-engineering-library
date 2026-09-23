begin;
-- AUTHORED service/approval/inventory binding; upstream owners remain authoritative.
alter table approval.request drop constraint request_kind_check;
alter table approval.request add constraint request_kind_check check(kind in
 ('reservation','sale','refund','payment','whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review','training_assessment','marketplace_mutation','whatsapp_schedule'));
alter table approval.request drop constraint request_payload_binding_check;
alter table approval.request add constraint request_payload_binding_check check
 (kind not in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review','training_assessment','marketplace_mutation','whatsapp_schedule') or
  (organization_id is not null and length(organization_id) between 1 and 128 and
   payload is not null and jsonb_typeof(payload)='object' and pg_column_size(payload)<=65536 and amount_minor_units=0));
create or replace function approval.guard_bound_request() returns trigger language plpgsql as $$
begin
 if tg_op='DELETE' then
  if old.kind in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review','training_assessment','marketplace_mutation','whatsapp_schedule') then raise exception 'bound approval request is immutable'; end if;
  return old;
 end if;
 if old.kind in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review','training_assessment','marketplace_mutation','whatsapp_schedule') or new.kind in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review','training_assessment','marketplace_mutation','whatsapp_schedule') then
  if row(new.tenant_id,new.request_id,new.kind,new.subject_id,new.amount_minor_units,new.requester,new.evidence_sha,new.organization_id,new.payload,new.created_at)
   is distinct from row(old.tenant_id,old.request_id,old.kind,old.subject_id,old.amount_minor_units,old.requester,old.evidence_sha,old.organization_id,old.payload,old.created_at)
   or (old.state<>'pending' and row(new.state,new.decided_at) is distinct from row(old.state,old.decided_at))
  then raise exception 'bound approval request is immutable'; end if;
 end if;
 return new;
end $$;
create or replace function approval.guard_bound_decision() returns trigger language plpgsql as $$
begin
 if exists(select 1 from approval.request where tenant_id=old.tenant_id and request_id=old.request_id and kind in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review','training_assessment','marketplace_mutation','whatsapp_schedule')) then
  raise exception 'bound approval decision is immutable';
 end if;
 if tg_op='DELETE' then return old;end if;return new;
end $$;


create table communication.whatsapp_schedule(
 tenant_id uuid not null,delivery_key text not null,job_id uuid not null,
 organization_id text not null,appointment_id text not null,appointment_version bigint not null,
 not_before timestamptz not null,expires_at timestamptz not null,
 request_sha256 text not null check(request_sha256~'^[0-9a-f]{64}$'),
 primary key(tenant_id,delivery_key),unique(tenant_id,job_id),
 unique(tenant_id,appointment_id,appointment_version,not_before),
 foreign key(tenant_id,delivery_key) references approval.request(tenant_id,request_id),
 foreign key(tenant_id,appointment_id) references crm.appointment(tenant_id,appointment_id),
 check(expires_at>not_before)
);
create trigger whatsapp_schedule_immutable before update or delete on communication.whatsapp_schedule
 for each row execute function catalog.release_immutable();
create table communication.whatsapp_schedule_cancellation(
 tenant_id uuid not null,delivery_key text not null,actor_subject text not null,reason text not null,
 request_sha256 text not null check(request_sha256~'^[0-9a-f]{64}$'),
 cancelled_at timestamptz not null default clock_timestamp(),
 primary key(tenant_id,delivery_key),
 foreign key(tenant_id,delivery_key) references communication.whatsapp_schedule(tenant_id,delivery_key),
 check(length(actor_subject)between 1 and 128),check(length(reason)between 1 and 2048)
);
create trigger whatsapp_schedule_cancellation_immutable before update or delete on communication.whatsapp_schedule_cancellation
 for each row execute function catalog.release_immutable();
create table communication.whatsapp_schedule_result(
 tenant_id uuid not null,delivery_key text not null,job_id uuid not null,
 outcome text not null check(outcome in('ACCEPTED','SUPPRESSED','RECONCILIATION_REQUIRED','FAILED_TERMINAL')),
 completed_at timestamptz not null default clock_timestamp(),
 primary key(tenant_id,delivery_key),
 foreign key(tenant_id,delivery_key) references communication.whatsapp_schedule(tenant_id,delivery_key),
 foreign key(tenant_id,job_id) references platform.job(tenant_id,job_id)
);
create trigger whatsapp_schedule_result_immutable before update or delete on communication.whatsapp_schedule_result
 for each row execute function catalog.release_immutable();
create or replace view communication.whatsapp_delivery_approval as
 select g.tenant_id,g.channel_code,g.delivery_key,g.organization_id,g.appointment_id,g.confirmation_event_id::text as approval_event_id,g.lead_id,g.external_id_hmac,g.profile_sha256,g.message_sha256,g.expires_at
 from communication.whatsapp_appointment_approval g
 union all
 select r.tenant_id,'whatsapp',r.request_id,r.organization_id,''::text,r.request_id,
 r.payload->>'lead_id',r.payload->>'external_id_hmac',r.payload->>'profile_sha256',r.payload->>'message_sha256',(r.payload->>'expires_at')::timestamptz
 from approval.request r
 where r.kind='whatsapp_reply' and r.state='approved' and r.payload->>'schema'='elite-whatsapp-reply-approval/v1'
 and r.payload->>'organization_id'=r.organization_id
 union all
 select r.tenant_id,'whatsapp',r.request_id,r.organization_id,''::text,r.request_id,
 r.payload->>'lead_id',r.payload->>'external_id_hmac',r.payload->>'profile_sha256',
 r.payload->>'message_sha256',(r.payload->>'expires_at')::timestamptz
 from approval.request r join communication.whatsapp_schedule s on s.tenant_id=r.tenant_id and s.delivery_key=r.request_id
 where r.kind='whatsapp_schedule' and r.state='approved' and r.payload->>'schema'='elite-whatsapp-schedule/v1'
 and r.payload->>'organization_id'=r.organization_id;
commit;
