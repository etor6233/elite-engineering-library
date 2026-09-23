begin;

do $$
begin
  if not exists(select 1 from information_schema.columns where table_schema='org' and table_name='organization' and column_name='version') then raise exception 'organization version missing'; end if;
  if not exists(select 1 from information_schema.columns where table_schema='franchise' and table_name='agreement' and column_name='version') then raise exception 'agreement version missing'; end if;
  if to_regprocedure('org.enforce_organization_hierarchy()') is null then raise exception 'hierarchy guard missing'; end if;
  if to_regprocedure('franchise.enforce_active_territory_non_overlap()') is null then raise exception 'territory guard missing'; end if;
end;
$$;

insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)
values('018f4d4a-7b36-7a21-8d10-2f4c54c28900','network-test','Network Test','Network Test');

insert into org.organization(tenant_id,organization_id,parent_organization_id,organization_code,display_name,organization_type)
values
('018f4d4a-7b36-7a21-8d10-2f4c54c28900','root',null,'root','Root','enterprise'),
('018f4d4a-7b36-7a21-8d10-2f4c54c28900','franchisee','root','franchisee','Franchisee','franchisee'),
('018f4d4a-7b36-7a21-8d10-2f4c54c28900','store','franchisee','store','Store','store');

do $$
begin
  begin
    update org.organization set parent_organization_id='store'
    where tenant_id='018f4d4a-7b36-7a21-8d10-2f4c54c28900' and organization_id='root';
    raise exception 'organization hierarchy accepted a cycle';
  exception when check_violation then
    null;
  end;
end;
$$;

insert into franchise.agreement(tenant_id,agreement_id,franchise_organization_id,territory_code,terms_version,starts_on,status)
values
('018f4d4a-7b36-7a21-8d10-2f4c54c28900','active-agreement','franchisee','AR-CBA','v1','2026-01-01','active'),
('018f4d4a-7b36-7a21-8d10-2f4c54c28900','draft-agreement','franchisee','AR-CBA','v2','2026-06-01','draft');

do $$
begin
  begin
    update franchise.agreement set status='active'
    where tenant_id='018f4d4a-7b36-7a21-8d10-2f4c54c28900' and agreement_id='draft-agreement';
    raise exception 'overlapping active territory accepted';
  exception when exclusion_violation then
    null;
  end;
end;
$$;

rollback;
