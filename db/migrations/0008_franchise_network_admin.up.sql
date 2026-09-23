begin;

alter table org.organization add column version bigint not null default 1 check (version > 0);
alter table franchise.agreement add column version bigint not null default 1 check (version > 0), add column updated_at timestamptz not null default clock_timestamp();

do $$
begin
  if exists (
    with recursive hierarchy as (
      select o.tenant_id,o.organization_id,o.parent_organization_id,array[o.organization_id]::text[] as path,false as cycle
      from org.organization o
      union all
      select p.tenant_id,p.organization_id,p.parent_organization_id,h.path || p.organization_id,p.organization_id=any(h.path)
      from hierarchy h
      join org.organization p on p.tenant_id=h.tenant_id and p.organization_id=h.parent_organization_id
      where not h.cycle
    )
    select 1 from hierarchy where cycle
  ) then
    raise exception using errcode='23514', message='existing organization hierarchy contains a cycle';
  end if;
end;
$$;

create or replace function org.enforce_organization_hierarchy()
returns trigger language plpgsql as $$
begin
  perform pg_advisory_xact_lock(hashtextextended(new.tenant_id::text || '|organization-hierarchy', 0));
  if new.parent_organization_id is not null then
    if not exists (select 1 from org.organization p where p.tenant_id=new.tenant_id and p.organization_id=new.parent_organization_id and p.status <> 'closed') then
      raise exception using errcode='23503', message='organization parent is unavailable';
    end if;
    if exists (
      with recursive ancestors as (
        select p.organization_id,p.parent_organization_id,array[p.organization_id]::text[] as path,false as cycle from org.organization p where p.tenant_id=new.tenant_id and p.organization_id=new.parent_organization_id
        union all
        select p.organization_id,p.parent_organization_id,a.path || p.organization_id,p.organization_id=any(a.path) from org.organization p join ancestors a on a.parent_organization_id=p.organization_id where p.tenant_id=new.tenant_id and not a.cycle
      ) select 1 from ancestors where organization_id=new.organization_id or cycle
    ) then
      raise exception using errcode='23514', message='organization hierarchy cycle';
    end if;
  end if;
  return new;
end;
$$;

create trigger organization_hierarchy_guard before insert or update of parent_organization_id on org.organization for each row execute function org.enforce_organization_hierarchy();

create or replace function franchise.enforce_active_territory_non_overlap()
returns trigger language plpgsql as $$
begin
  perform pg_advisory_xact_lock(hashtextextended(new.tenant_id::text || '|' || new.territory_code, 0));
  if new.status='active' and exists (
    select 1 from franchise.agreement a where a.tenant_id=new.tenant_id and a.territory_code=new.territory_code and a.agreement_id<>new.agreement_id and a.status='active'
      and daterange(a.starts_on,coalesce(a.ends_on,'infinity'::date),'[)') && daterange(new.starts_on,coalesce(new.ends_on,'infinity'::date),'[)')
  ) then
    raise exception using errcode='23P01', message='active franchise territories overlap';
  end if;
  return new;
end;
$$;

create trigger franchise_territory_non_overlap before insert or update of territory_code,starts_on,ends_on,status on franchise.agreement for each row execute function franchise.enforce_active_territory_non_overlap();

create index organization_parent_status_idx on org.organization(tenant_id,parent_organization_id,status,organization_id);
create index agreement_territory_status_idx on franchise.agreement(tenant_id,territory_code,status,starts_on,agreement_id);

commit;
