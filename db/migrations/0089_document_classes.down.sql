begin;
do $$begin if exists(select 1 from document.original where mode='TYPED_FIXTURE') then raise exception 'typed document evidence present; rollback refused';end if;end$$;
create or replace function document.guard_review()returns trigger language plpgsql as $$
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
drop trigger typed_original_guard on document.original;
drop function document.guard_typed_original();
drop index document.document_inbox_scope;
drop index document.document_inbox_uploader;
alter table document.original drop constraint original_class_fk;
alter table document.original drop constraint original_mode_check;
alter table document.original add constraint original_mode_check check(mode in('FIXTURE','PROVIDER'));
alter table document.original drop column class_id,drop column schema_version;
drop function document.valid_typed_fields(text,text,jsonb);
drop table document.class_schema;
commit;
