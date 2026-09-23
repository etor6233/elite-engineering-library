-- 0043_finops.test.sql — verifica tags obligatorios (nada sin atribuir),
-- budget cap (spent <= cap) y cost entry por recurso.
-- Precondición: migración 0001 (platform.tenant) y 0043 aplicadas.

begin;

insert into platform.tenant (tenant_id, tenant_code, legal_name, display_name) values
  ('11111111-1111-1111-1111-111111111111', 'tenant-a', 'A', 'Tenant A');

insert into finops.resource
  (tenant_id, provider, resource_id, kind, region, tags, estimated_monthly_minor_units)
values
  ('11111111-1111-1111-1111-111111111111', 'gcp', 'r1', 'cloud_run', 'us-central1',
   '{"tenant":"t","environment":"prod","owner":"platform"}', 100);

-- recurso sin environment/owner → check_violation
do $$
begin
  begin
    insert into finops.resource
      (tenant_id, provider, resource_id, kind, region, tags, estimated_monthly_minor_units)
    values
      ('11111111-1111-1111-1111-111111111111', 'aws', 'r2', 'compute', 'us-east-1',
       '{"tenant":"t"}', 100);
    raise exception 'untagged resource accepted';
  exception when check_violation then
    null; -- expected
  end;
end $$;

-- budget: spent <= cap
insert into finops.budget (tenant_id, period, cap_minor_units, spent_minor_units) values
  ('11111111-1111-1111-1111-111111111111', '2026-09', 100, 60);

do $$
begin
  begin
    update finops.budget set spent_minor_units = 101
     where tenant_id = '11111111-1111-1111-1111-111111111111' and period = '2026-09';
    raise exception 'over-cap accepted';
  exception when check_violation then
    null; -- expected
  end;
end $$;

-- cost entry por recurso
insert into finops.cost_entry (tenant_id, provider, resource_id, period, amount_minor_units) values
  ('11111111-1111-1111-1111-111111111111', 'gcp', 'r1', '2026-09', 40);

do $$
declare n int;
begin
  select count(*) into n from finops.resource
   where tenant_id = '11111111-1111-1111-1111-111111111111';
  if n <> 1 then raise exception 'expected 1 resource, got %', n; end if;
end $$;

rollback;
