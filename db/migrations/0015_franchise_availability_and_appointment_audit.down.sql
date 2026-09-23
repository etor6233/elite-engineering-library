begin;
drop trigger if exists appointment_transition_immutable on crm.appointment_transition;
drop function if exists crm.prevent_appointment_transition_mutation();
drop table if exists crm.appointment_transition;
drop trigger if exists appointment_slot_availability_guard on crm.appointment_slot;
drop function if exists crm.enforce_appointment_slot_availability();
drop trigger if exists availability_cancellation_guard on crm.availability_entry;
drop function if exists crm.enforce_availability_cancellation();
drop trigger if exists availability_entry_guard on crm.availability_entry;
drop function if exists crm.enforce_availability_entry();
drop table if exists crm.availability_entry;
alter table crm.service_resource drop constraint if exists service_resource_org_identity_uq;
-- Restore the assignment guard owned by migration 0009.
create or replace function crm.enforce_appointment_resource_assignment() returns trigger language plpgsql as $$
declare candidate crm.appointment%rowtype;
begin
  perform pg_advisory_xact_lock(hashtextextended(new.tenant_id::text || '|appointment-resource|' || new.resource_id,0));
  select * into candidate from crm.appointment a where a.tenant_id=new.tenant_id and a.appointment_id=new.appointment_id;
  if candidate.appointment_id is null or candidate.state <> 'requested' then raise exception using errcode='23514',message='appointment is unavailable for resource assignment'; end if;
  if not exists(select 1 from crm.service_resource r join crm.resource_skill s on s.tenant_id=r.tenant_id and s.resource_id=r.resource_id where r.tenant_id=new.tenant_id and r.resource_id=new.resource_id and r.organization_id=candidate.organization_id and r.status='active' and s.appointment_kind=candidate.appointment_kind) then raise exception using errcode='23514',message='resource is unavailable or lacks appointment skill'; end if;
  if exists(select 1 from crm.appointment_resource ar join crm.appointment a on a.tenant_id=ar.tenant_id and a.appointment_id=ar.appointment_id where ar.tenant_id=new.tenant_id and ar.resource_id=new.resource_id and ar.appointment_id<>new.appointment_id and a.state in ('requested','confirmed') and tstzrange(a.starts_at,coalesce(a.ends_at,a.starts_at+interval '1 hour'),'[)') && tstzrange(candidate.starts_at,coalesce(candidate.ends_at,candidate.starts_at+interval '1 hour'),'[)')) then raise exception using errcode='23P01',message='resource appointment assignments overlap'; end if;
  return new;
end $$;
commit;
