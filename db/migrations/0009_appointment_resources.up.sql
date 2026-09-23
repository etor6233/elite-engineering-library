begin;

create table crm.service_resource (
  tenant_id uuid not null,
  resource_id text not null,
  organization_id text not null,
  principal_subject text,
  display_name text not null,
  resource_kind text not null check (resource_kind in ('employee','contractor','service-bay','vehicle','equipment')),
  status text not null check (status in ('active','inactive')),
  version bigint not null check (version > 0),
  created_at timestamptz not null default clock_timestamp(),
  updated_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id,resource_id),
  foreign key (tenant_id,organization_id) references org.organization(tenant_id,organization_id),
  check (length(display_name) between 1 and 160),
  check ((resource_kind in ('employee','contractor')) = (principal_subject is not null))
);

create unique index service_resource_principal_idx on crm.service_resource(tenant_id,organization_id,principal_subject) where principal_subject is not null;

create table crm.resource_skill (
  tenant_id uuid not null,
  resource_id text not null,
  appointment_kind text not null check (appointment_kind in ('consultation','test-drive','delivery','service')),
  primary key (tenant_id,resource_id,appointment_kind),
  foreign key (tenant_id,resource_id) references crm.service_resource(tenant_id,resource_id) on delete cascade
);

create table crm.appointment_resource (
  tenant_id uuid not null,
  appointment_id text not null,
  resource_id text not null,
  assigned_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id,appointment_id),
  foreign key (tenant_id,appointment_id) references crm.appointment(tenant_id,appointment_id) on delete cascade,
  foreign key (tenant_id,resource_id) references crm.service_resource(tenant_id,resource_id)
);

create or replace function crm.enforce_appointment_resource_assignment()
returns trigger language plpgsql as $$
declare candidate crm.appointment%rowtype;
begin
  perform pg_advisory_xact_lock(hashtextextended(new.tenant_id::text || '|appointment-resource|' || new.resource_id,0));
  select * into candidate from crm.appointment a where a.tenant_id=new.tenant_id and a.appointment_id=new.appointment_id;
  if candidate.appointment_id is null or candidate.state <> 'requested' then
    raise exception using errcode='23514', message='appointment is unavailable for resource assignment';
  end if;
  if not exists (
    select 1 from crm.service_resource r join crm.resource_skill s on s.tenant_id=r.tenant_id and s.resource_id=r.resource_id
    where r.tenant_id=new.tenant_id and r.resource_id=new.resource_id and r.organization_id=candidate.organization_id and r.status='active' and s.appointment_kind=candidate.appointment_kind
  ) then
    raise exception using errcode='23514', message='resource is unavailable or lacks appointment skill';
  end if;
  if exists (
    select 1 from crm.appointment_resource ar join crm.appointment a on a.tenant_id=ar.tenant_id and a.appointment_id=ar.appointment_id
    where ar.tenant_id=new.tenant_id and ar.resource_id=new.resource_id and ar.appointment_id<>new.appointment_id and a.state in ('requested','confirmed')
      and tstzrange(a.starts_at,coalesce(a.ends_at,a.starts_at+interval '1 hour'),'[)') && tstzrange(candidate.starts_at,coalesce(candidate.ends_at,candidate.starts_at+interval '1 hour'),'[)')
  ) then
    raise exception using errcode='23P01', message='resource appointment assignments overlap';
  end if;
  return new;
end;
$$;

create trigger appointment_resource_assignment_guard before insert or update of resource_id on crm.appointment_resource for each row execute function crm.enforce_appointment_resource_assignment();
create index service_resource_org_status_idx on crm.service_resource(tenant_id,organization_id,status,resource_id);
create index appointment_resource_schedule_idx on crm.appointment_resource(tenant_id,resource_id,appointment_id);

commit;
