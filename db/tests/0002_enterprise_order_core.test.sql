begin;

insert into platform.tenant (tenant_id, tenant_code, legal_name, display_name)
values
  ('018f4d4a-7b36-7a21-8d10-2f4c54c28a01', 'tenant-a', 'Tenant A S.A.', 'Tenant A'),
  ('018f4d4a-7b36-7a21-8d10-2f4c54c28a02', 'tenant-b', 'Tenant B S.A.', 'Tenant B');

insert into org.organization (
  tenant_id, organization_id, organization_code, display_name, organization_type
) values
  ('018f4d4a-7b36-7a21-8d10-2f4c54c28a01', 'org-main', 'main', 'Main A', 'enterprise'),
  ('018f4d4a-7b36-7a21-8d10-2f4c54c28a02', 'org-main', 'main', 'Main B', 'enterprise');

insert into sales.customer_order (
  tenant_id, order_id, organization_id, customer_principal_id,
  state, currency, total_minor_units, version
) values
  ('018f4d4a-7b36-7a21-8d10-2f4c54c28a01', 'order-1', 'org-main', 'customer-a', 'draft', 'USD', 100, 1),
  ('018f4d4a-7b36-7a21-8d10-2f4c54c28a02', 'order-1', 'org-main', 'customer-b', 'draft', 'USD', 200, 1);

do $test$
declare
  tenant_a_count bigint;
  tenant_b_total bigint;
begin
  select count(*) into tenant_a_count
    from sales.customer_order
   where tenant_id = '018f4d4a-7b36-7a21-8d10-2f4c54c28a01';
  select sum(total_minor_units) into tenant_b_total
    from sales.customer_order
   where tenant_id = '018f4d4a-7b36-7a21-8d10-2f4c54c28a02';
  if tenant_a_count <> 1 or tenant_b_total <> 200 then
    raise exception 'tenant isolation fixture failed';
  end if;
end;
$test$;

do $test$
begin
  begin
    insert into sales.customer_order (
      tenant_id, order_id, organization_id, customer_principal_id,
      state, currency, total_minor_units, version
    ) values (
      '018f4d4a-7b36-7a21-8d10-2f4c54c28a01', 'invalid-order', 'org-main',
      'customer-a', 'invented', 'USD', 0, 1
    );
    raise exception 'expected state check violation';
  exception
    when check_violation then null;
  end;
end;
$test$;

do $test$
declare
  changed bigint;
begin
  update sales.customer_order
     set state = 'placed', version = 2, updated_at = clock_timestamp()
   where tenant_id = '018f4d4a-7b36-7a21-8d10-2f4c54c28a01'
     and order_id = 'order-1'
     and version = 1;
  get diagnostics changed = row_count;
  if changed <> 1 then
    raise exception 'expected one optimistic update, got %', changed;
  end if;

  update sales.customer_order
     set state = 'confirmed', version = 3, updated_at = clock_timestamp()
   where tenant_id = '018f4d4a-7b36-7a21-8d10-2f4c54c28a01'
     and order_id = 'order-1'
     and version = 1;
  get diagnostics changed = row_count;
  if changed <> 0 then
    raise exception 'stale optimistic update was accepted';
  end if;
end;
$test$;

rollback;
