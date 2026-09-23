begin;

create table crm.appointment_slot (
  tenant_id uuid not null,
  slot_id text not null,
  organization_id text not null,
  appointment_kind text not null check (appointment_kind in ('consultation','test-drive','delivery','service')),
  starts_at timestamptz not null,
  ends_at timestamptz not null,
  capacity integer not null check (capacity between 1 and 100),
  state text not null check (state in ('open','closed')),
  version bigint not null check (version > 0),
  created_at timestamptz not null default clock_timestamp(),
  updated_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, slot_id),
  unique (tenant_id, organization_id, appointment_kind, starts_at),
  foreign key (tenant_id, organization_id) references org.organization (tenant_id, organization_id),
  check (ends_at > starts_at),
  check (ends_at - starts_at <= interval '8 hours'),
  check (updated_at >= created_at)
);

alter table crm.appointment
  add column slot_id text,
  add column ends_at timestamptz,
  add constraint appointment_slot_fk foreign key (tenant_id, slot_id) references crm.appointment_slot (tenant_id, slot_id),
  add constraint appointment_slot_time_check check ((slot_id is null and ends_at is null) or (slot_id is not null and ends_at > starts_at));

create index appointment_slot_public_idx on crm.appointment_slot (tenant_id, organization_id, appointment_kind, state, starts_at);
create index appointment_slot_booking_idx on crm.appointment (tenant_id, slot_id, state) where slot_id is not null;

create or replace function crm.enforce_appointment_slot_non_overlap()
returns trigger
language plpgsql
as $$
begin
  perform pg_advisory_xact_lock(hashtextextended(new.tenant_id::text || '|' || new.organization_id || '|' || new.appointment_kind, 0));
  if new.state = 'open' and exists (
    select 1
    from crm.appointment_slot existing
    where existing.tenant_id = new.tenant_id
      and existing.organization_id = new.organization_id
      and existing.appointment_kind = new.appointment_kind
      and existing.state = 'open'
      and existing.slot_id <> new.slot_id
      and tstzrange(existing.starts_at, existing.ends_at, '[)') && tstzrange(new.starts_at, new.ends_at, '[)')
  ) then
    raise exception using errcode = '23P01', message = 'appointment slots overlap';
  end if;
  return new;
end;
$$;

create trigger appointment_slot_non_overlap
before insert or update of organization_id,appointment_kind,starts_at,ends_at,state
on crm.appointment_slot
for each row execute function crm.enforce_appointment_slot_non_overlap();

commit;
