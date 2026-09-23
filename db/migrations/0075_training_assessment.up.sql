begin;
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

create index audit_training_attempt_idx on audit.event(tenant_id,resource_id,actor_subject,action,audit_sequence) where resource_type='training-attempt';
create index approval_training_review_idx on approval.request(tenant_id,organization_id,created_at desc,request_id) where kind='training_assessment';
create unique index training_assessment_attempt_idx on approval.request(tenant_id,organization_id,(payload->>'attempt_id')) where kind='training_assessment';
create function audit.protect_training_fact() returns trigger language plpgsql as $$
begin
 if old.resource_type='training-attempt' or (tg_op='UPDATE' and new.resource_type='training-attempt') then
  raise exception 'training participation facts are immutable';
 end if;
 if tg_op='DELETE' then return old;end if;return new;
end$$;
create trigger training_audit_immutable before update or delete on audit.event for each row execute function audit.protect_training_fact();
create function approval.require_training_facts() returns trigger language plpgsql as $$
declare r approval.request%rowtype; a jsonb;
begin
 if new.kind<>'training_assessment' then return new;end if;
 select * into strict r from approval.request where tenant_id=new.tenant_id and request_id=new.request_id;
 if r.payload->>'schema'<>'training-assessment/v1' or r.payload->>'learner_subject' is distinct from r.requester
  or r.subject_id is distinct from r.requester or r.payload->>'organization_id' is distinct from r.organization_id then
  raise exception 'training assessment identity mismatch';
 end if;
 a:=r.payload;
 if not exists(select 1 from audit.event e where e.tenant_id=r.tenant_id and e.resource_type='training-attempt'
  and e.resource_id=a->>'attempt_id' and e.actor_subject=r.requester and e.action='training.started'
  and e.evidence->>'organization_id'=r.organization_id and e.evidence->'content'=a->'content') then
  raise exception 'training assessment has no matching immutable start';
 end if;
 if not exists(select 1 from audit.event e where e.tenant_id=r.tenant_id and e.resource_type='training-attempt'
  and e.resource_id=a->>'attempt_id' and e.actor_subject=r.requester and e.action='training.submitted'
  and e.evidence->>'organization_id'=r.organization_id and e.evidence->>'payload_sha256'=r.evidence_sha) then
  raise exception 'training assessment submission not recorded';
 end if;
 if exists(select 1 from jsonb_array_elements_text(a#>'{content,course,lessons}') lesson
  where not exists(select 1 from audit.event e where e.tenant_id=r.tenant_id and e.resource_type='training-attempt'
   and e.resource_id=a->>'attempt_id' and e.actor_subject=r.requester and e.action='training.lesson-read'
   and e.evidence->>'lesson_id'=lesson and e.evidence->>'organization_id'=r.organization_id
   and e.evidence->>'profile_sha256'=a#>>'{content,profile_sha256}')) then
  raise exception 'training assessment requires recorded lesson declarations';
 end if;
 if r.state<>'pending' and not exists(select 1 from approval.decision d join audit.event e
  on e.tenant_id=d.tenant_id and e.actor_subject=d.reviewer and e.resource_type='training-attempt'
   and e.resource_id=a->>'attempt_id' and e.action='training.assessed'
  where d.tenant_id=r.tenant_id and d.request_id=r.request_id and d.reviewer<>r.requester
   and e.evidence->>'organization_id'=r.organization_id and e.evidence->>'payload_sha256'=r.evidence_sha
   and e.evidence->>'reason'=d.reason and (e.evidence->>'approved')::boolean=d.approved
   and d.approved=(r.state='approved')) then
  raise exception 'training assessment decision has no matching human review fact';
 end if;
 return new;
end$$;
create constraint trigger training_assessment_guard after insert or update on approval.request
 deferrable initially deferred for each row execute function approval.require_training_facts();
commit;
