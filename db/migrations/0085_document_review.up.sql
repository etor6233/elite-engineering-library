begin;
-- AUTHORED service/approval/inventory binding; upstream owners remain authoritative.
alter table approval.request drop constraint request_kind_check;
alter table approval.request add constraint request_kind_check check(kind in
 ('reservation','sale','refund','payment','whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review','training_assessment','marketplace_mutation','whatsapp_schedule','document_review'));
alter table approval.request drop constraint request_payload_binding_check;
alter table approval.request add constraint request_payload_binding_check check
 (kind not in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review','training_assessment','marketplace_mutation','whatsapp_schedule','document_review') or
  (organization_id is not null and length(organization_id) between 1 and 128 and
   payload is not null and jsonb_typeof(payload)='object' and pg_column_size(payload)<=65536 and amount_minor_units=0));
create or replace function approval.guard_bound_request() returns trigger language plpgsql as $$
begin
 if tg_op='DELETE' then
  if old.kind in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review','training_assessment','marketplace_mutation','whatsapp_schedule','document_review') then raise exception 'bound approval request is immutable'; end if;
  return old;
 end if;
 if old.kind in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review','training_assessment','marketplace_mutation','whatsapp_schedule','document_review') or new.kind in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review','training_assessment','marketplace_mutation','whatsapp_schedule','document_review') then
  if row(new.tenant_id,new.request_id,new.kind,new.subject_id,new.amount_minor_units,new.requester,new.evidence_sha,new.organization_id,new.payload,new.created_at)
   is distinct from row(old.tenant_id,old.request_id,old.kind,old.subject_id,old.amount_minor_units,old.requester,old.evidence_sha,old.organization_id,old.payload,old.created_at)
   or (old.state<>'pending' and row(new.state,new.decided_at) is distinct from row(old.state,old.decided_at))
  then raise exception 'bound approval request is immutable'; end if;
 end if;
 return new;
end $$;
create or replace function approval.guard_bound_decision() returns trigger language plpgsql as $$
begin
 if exists(select 1 from approval.request where tenant_id=old.tenant_id and request_id=old.request_id and kind in ('whatsapp_reply','social_publish','social_revoke','stored_value_operation','warranty_repair','serial_quality','catalog_review','training_assessment','marketplace_mutation','whatsapp_schedule','document_review')) then
  raise exception 'bound approval decision is immutable';
 end if;
 if tg_op='DELETE' then return old;end if;return new;
end $$;



create schema document;
create table document.original(
 tenant_id uuid not null,document_id uuid not null,organization_id text not null,
 uploader text not null,name text not null,original_sha256 text not null,
 profile_sha256 text not null,mode text not null,content bytea not null,
 received_at timestamptz not null default clock_timestamp(),
 primary key(tenant_id,document_id),
 foreign key(tenant_id,organization_id)references org.organization(tenant_id,organization_id),
 check(length(uploader)between 1 and 128),check(length(name)between 1 and 128),
 check(original_sha256~'^[0-9a-f]{64}$'and profile_sha256~'^[0-9a-f]{64}$'),
 check(encode(sha256(content),'hex')=original_sha256),
 check(octet_length(content)between 1 and 2097152),check(mode in('FIXTURE','PROVIDER')),
 constraint original_fixture_identity check(mode<>'FIXTURE'or original_sha256='489f0c63b6a05e0ae0fcfbf29121299fc2ccacab35d5cd40a5d63d42764078cb')
);
create table document.extraction(
 tenant_id uuid not null,document_id uuid not null,job_id uuid not null,attempt integer not null,
 security_receipt bytea not null,provider_response bytea not null,analysis_receipt bytea not null,
 evidence_sha256 text not null,suggested jsonb not null,
 extracted_at timestamptz not null default clock_timestamp(),
 primary key(tenant_id,document_id),
 foreign key(tenant_id,document_id)references document.original(tenant_id,document_id),
 foreign key(tenant_id,job_id)references platform.job(tenant_id,job_id),
 check(attempt>0),check(evidence_sha256~'^[0-9a-f]{64}$'),
 check(encode(sha256(analysis_receipt),'hex')=evidence_sha256),
 check(octet_length(security_receipt)between 1 and 65536),
 check(octet_length(provider_response)between 1 and 4194304),
 check(octet_length(analysis_receipt)between 1 and 32768),
 check(jsonb_typeof(suggested)='object')
);
create table document.attempt_failure(
 tenant_id uuid not null,document_id uuid not null,attempt integer not null,code text not null,
 failed_at timestamptz not null default clock_timestamp(),
 primary key(tenant_id,document_id,attempt),
 foreign key(tenant_id,document_id)references document.original(tenant_id,document_id),
 check(code in('SECURITY_REJECTED','EXTRACTION_UNAVAILABLE','INVALID_CONTRACT'))
);
create table document.committed(
 tenant_id uuid not null,document_id uuid not null,request_id text not null,
 payload_sha256 text not null,fields jsonb not null,reviewer text not null,
 committed_at timestamptz not null default clock_timestamp(),
 primary key(tenant_id,document_id),unique(tenant_id,request_id),
 foreign key(tenant_id,document_id)references document.extraction(tenant_id,document_id),
 foreign key(tenant_id,request_id)references approval.request(tenant_id,request_id),
 check(payload_sha256~'^[0-9a-f]{64}$'),check(jsonb_typeof(fields)='object'),check(length(reviewer)between 1 and 128)
);
create unique index document_review_once on approval.request(tenant_id,organization_id,(payload->>'document_id'))where kind='document_review';
create trigger original_immutable before update or delete on document.original for each row execute function catalog.release_immutable();
create trigger extraction_immutable before update or delete on document.extraction for each row execute function catalog.release_immutable();
create trigger attempt_failure_immutable before update or delete on document.attempt_failure for each row execute function catalog.release_immutable();
create trigger committed_immutable before update or delete on document.committed for each row execute function catalog.release_immutable();
create function document.guard_review()returns trigger language plpgsql as $$
begin
 if new.kind='document_review' and not exists(
 select 1 from document.original o join document.extraction e using(tenant_id,document_id)
 where o.tenant_id=new.tenant_id and o.document_id::text=new.payload->>'document_id'
 and o.organization_id=new.organization_id and o.original_sha256=new.payload->>'original_sha256'
 and o.profile_sha256=new.payload->>'profile_sha256'and o.mode=new.payload->>'mode'
 and e.evidence_sha256=new.payload->>'evidence_sha256'
 and new.payload->>'schema'='document-review/v1'
 and jsonb_typeof(new.payload->'fields')='object'
 and(select count(*)from jsonb_object_keys(new.payload->'fields'))=4
 and jsonb_typeof(new.payload->'fields'->'invoice_number')='string'and octet_length(new.payload->'fields'->>'invoice_number')between 1 and 128
 and jsonb_typeof(new.payload->'fields'->'vendor')='string'and octet_length(new.payload->'fields'->>'vendor')between 1 and 512
 and jsonb_typeof(new.payload->'fields'->'total')='string'and octet_length(new.payload->'fields'->>'total')between 1 and 64
 and jsonb_typeof(new.payload->'fields'->'currency')='string'and(new.payload->'fields'->>'currency')~'^[A-Z]{3}$'
 and not exists(select 1 from jsonb_each_text(new.payload->'fields')f where f.value<>btrim(f.value)or f.value='')
 )then raise exception 'document review evidence not bound';end if;
 return new;
end$$;
create trigger document_review_guard before insert on approval.request for each row execute function document.guard_review();
create function document.guard_commit()returns trigger language plpgsql as $$
begin
 if not exists(select 1 from approval.request r join approval.decision d using(tenant_id,request_id)
 where r.tenant_id=new.tenant_id and r.request_id=new.request_id and r.kind='document_review'
 and r.state='approved'and d.approved and d.reviewer<>r.requester and d.reviewer=new.reviewer
 and r.evidence_sha=new.payload_sha256 and r.payload->>'document_id'=new.document_id::text
 and r.payload->'fields'=new.fields
 and(select count(*)from approval.decision x where x.tenant_id=r.tenant_id and x.request_id=r.request_id)=1)
 then raise exception 'document commit requires exact human approval';end if;
 return new;
end$$;
create constraint trigger document_commit_guard after insert on document.committed deferrable initially deferred for each row execute function document.guard_commit();
commit;
