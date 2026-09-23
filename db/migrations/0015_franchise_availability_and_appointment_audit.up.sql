begin;

alter table crm.service_resource add constraint service_resource_org_identity_uq unique (tenant_id,organization_id,resource_id);

create table crm.availability_entry (
  tenant_id uuid not null,
  availability_id text not null,
  organization_id text not null,
  resource_id text,
  entry_type text not null check (entry_type in ('working','unavailable')),
  reason_code text,
  starts_at timestamptz not null,
  ends_at timestamptz not null,
  state text not null check (state in ('active','cancelled')),
  version bigint not null check (version>0),
  created_by_subject text not null check (length(created_by_subject) between 1 and 255),
  created_at timestamptz not null default clock_timestamp(),
  cancelled_by_subject text,
  cancellation_reason_code text,
  cancelled_at timestamptz,
  primary key (tenant_id,availability_id),
  foreign key (tenant_id,organization_id) references org.organization(tenant_id,organization_id),
  foreign key (tenant_id,organization_id,resource_id) references crm.service_resource(tenant_id,organization_id,resource_id),
  check (ends_at>starts_at and ends_at-starts_at<=interval '366 days'),
  check ((entry_type='working' and reason_code is null) or (entry_type='unavailable' and reason_code ~ '^[a-z][a-z0-9]*(-[a-z0-9]+)*$')),
  check ((state='active' and cancelled_by_subject is null and cancellation_reason_code is null and cancelled_at is null) or (state='cancelled' and cancelled_by_subject is not null and cancellation_reason_code is not null and cancelled_at is not null))
);

create index availability_scope_time_idx on crm.availability_entry(tenant_id,organization_id,resource_id,starts_at,ends_at) where state='active';

create or replace function crm.enforce_availability_entry() returns trigger language plpgsql as $$
begin
  perform pg_advisory_xact_lock(hashtextextended(new.tenant_id::text||'|availability|'||new.organization_id||'|'||coalesce(new.resource_id,'organization'),0));
  if exists (
    select 1 from crm.availability_entry e where e.tenant_id=new.tenant_id and e.organization_id=new.organization_id
      and e.resource_id is not distinct from new.resource_id and e.entry_type=new.entry_type and e.state='active'
      and e.availability_id<>new.availability_id and tstzrange(e.starts_at,e.ends_at,'[)') && tstzrange(new.starts_at,new.ends_at,'[)')
  ) then raise exception using errcode='23P01',message='availability entries overlap'; end if;
  if new.entry_type='unavailable' and exists (
    select 1 from crm.appointment a left join crm.appointment_resource ar on ar.tenant_id=a.tenant_id and ar.appointment_id=a.appointment_id
    where a.tenant_id=new.tenant_id and a.organization_id=new.organization_id and a.state in ('requested','confirmed')
      and (new.resource_id is null or ar.resource_id=new.resource_id)
      and tstzrange(a.starts_at,coalesce(a.ends_at,a.starts_at+interval '1 hour'),'[)') && tstzrange(new.starts_at,new.ends_at,'[)')
  ) then raise exception using errcode='23P01',message='unavailability conflicts with active appointments'; end if;
  return new;
end $$;
create trigger availability_entry_guard before insert or update of organization_id,resource_id,entry_type,starts_at,ends_at,state on crm.availability_entry for each row execute function crm.enforce_availability_entry();

create or replace function crm.enforce_availability_cancellation() returns trigger language plpgsql as $$
begin
  if old.state='active' and new.state='cancelled' and old.entry_type='working' and exists (
    select 1 from crm.appointment a left join crm.appointment_resource ar on ar.tenant_id=a.tenant_id and ar.appointment_id=a.appointment_id
    where a.tenant_id=old.tenant_id and a.organization_id=old.organization_id and a.state in ('requested','confirmed')
      and (old.resource_id is null or ar.resource_id=old.resource_id)
      and tstzrange(a.starts_at,coalesce(a.ends_at,a.starts_at+interval '1 hour'),'[)') && tstzrange(old.starts_at,old.ends_at,'[)')
  ) then raise exception using errcode='23P01',message='working availability has active appointments'; end if;
  return new;
end $$;
create trigger availability_cancellation_guard before update of state on crm.availability_entry for each row execute function crm.enforce_availability_cancellation();

create or replace function crm.enforce_appointment_slot_availability() returns trigger language plpgsql as $$
begin
  if new.state='open' and (not exists(select 1 from crm.availability_entry e where e.tenant_id=new.tenant_id and e.organization_id=new.organization_id and e.resource_id is null and e.entry_type='working' and e.state='active' and e.starts_at<=new.starts_at and e.ends_at>=new.ends_at)
    or exists(select 1 from crm.availability_entry e where e.tenant_id=new.tenant_id and e.organization_id=new.organization_id and e.resource_id is null and e.entry_type='unavailable' and e.state='active' and tstzrange(e.starts_at,e.ends_at,'[)') && tstzrange(new.starts_at,new.ends_at,'[)'))) then
    raise exception using errcode='23514',message='appointment slot is outside organization availability';
  end if;
  return new;
end $$;
create trigger appointment_slot_availability_guard before insert or update of organization_id,starts_at,ends_at,state on crm.appointment_slot for each row execute function crm.enforce_appointment_slot_availability();

create or replace function crm.enforce_appointment_resource_assignment() returns trigger language plpgsql as $$
declare candidate crm.appointment%rowtype;
begin
  perform pg_advisory_xact_lock(hashtextextended(new.tenant_id::text || '|appointment-resource|' || new.resource_id,0));
  select * into candidate from crm.appointment a where a.tenant_id=new.tenant_id and a.appointment_id=new.appointment_id;
  if candidate.appointment_id is null or candidate.state <> 'requested' then raise exception using errcode='23514',message='appointment is unavailable for resource assignment'; end if;
  if not exists(select 1 from crm.service_resource r join crm.resource_skill s on s.tenant_id=r.tenant_id and s.resource_id=r.resource_id where r.tenant_id=new.tenant_id and r.resource_id=new.resource_id and r.organization_id=candidate.organization_id and r.status='active' and s.appointment_kind=candidate.appointment_kind) then raise exception using errcode='23514',message='resource is unavailable or lacks appointment skill'; end if;
  if not exists(select 1 from crm.availability_entry e where e.tenant_id=new.tenant_id and e.organization_id=candidate.organization_id and e.resource_id=new.resource_id and e.entry_type='working' and e.state='active' and e.starts_at<=candidate.starts_at and e.ends_at>=coalesce(candidate.ends_at,candidate.starts_at+interval '1 hour'))
    or exists(select 1 from crm.availability_entry e where e.tenant_id=new.tenant_id and e.resource_id=new.resource_id and e.entry_type='unavailable' and e.state='active' and tstzrange(e.starts_at,e.ends_at,'[)') && tstzrange(candidate.starts_at,coalesce(candidate.ends_at,candidate.starts_at+interval '1 hour'),'[)')) then raise exception using errcode='23514',message='resource is outside working availability'; end if;
  if exists(select 1 from crm.appointment_resource ar join crm.appointment a on a.tenant_id=ar.tenant_id and a.appointment_id=ar.appointment_id where ar.tenant_id=new.tenant_id and ar.resource_id=new.resource_id and ar.appointment_id<>new.appointment_id and a.state in ('requested','confirmed') and tstzrange(a.starts_at,coalesce(a.ends_at,a.starts_at+interval '1 hour'),'[)') && tstzrange(candidate.starts_at,coalesce(candidate.ends_at,candidate.starts_at+interval '1 hour'),'[)')) then raise exception using errcode='23P01',message='resource appointment assignments overlap'; end if;
  return new;
end $$;

create table crm.appointment_transition (
  tenant_id uuid not null,
  transition_id uuid not null,
  appointment_id text not null,
  from_state text not null,
  to_state text not null,
  actor_subject text not null check (length(actor_subject) between 1 and 255),
  reason_code text,
  occurred_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id,transition_id),
  foreign key (tenant_id,appointment_id) references crm.appointment(tenant_id,appointment_id),
  check ((to_state in ('cancelled','no-show') and reason_code ~ '^[a-z][a-z0-9]*(-[a-z0-9]+)*$') or (to_state not in ('cancelled','no-show') and reason_code is null))
);
create index appointment_transition_history_idx on crm.appointment_transition(tenant_id,appointment_id,occurred_at,transition_id);
create or replace function crm.prevent_appointment_transition_mutation() returns trigger language plpgsql as $$ begin raise exception 'immutable appointment transition'; end $$;
create trigger appointment_transition_immutable before update or delete on crm.appointment_transition for each row execute function crm.prevent_appointment_transition_mutation();

commit;
