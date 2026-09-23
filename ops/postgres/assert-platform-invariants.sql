begin transaction read only;

do $$
declare
  missing_relations text[];
  invalid_count bigint;
begin
  select array_agg(required_name order by required_name)
    into missing_relations
  from unnest(array[
    'platform.tenant',
    'platform.business_profile_version',
    'platform.idempotency_record',
    'platform.outbox_event',
    'platform.consumer_inbox',
    'platform.job',
    'audit.event'
  ]) as required(required_name)
  where to_regclass(required_name) is null;

  if missing_relations is not null then
    raise exception 'missing required relations: %', missing_relations;
  end if;

  select count(*) into invalid_count
  from platform.business_profile_version p
  left join platform.tenant t on t.tenant_id = p.tenant_id
  where t.tenant_id is null;
  if invalid_count <> 0 then raise exception 'orphan business profiles: %', invalid_count; end if;

  select count(*) into invalid_count
  from (
    select tenant_id
    from platform.business_profile_version
    where status = 'active'
    group by tenant_id
    having count(*) > 1
  ) duplicated_active;
  if invalid_count <> 0 then raise exception 'tenants with multiple active profiles: %', invalid_count; end if;

  select count(*) into invalid_count
  from platform.outbox_event
  where aggregate_version <= 0 or available_at < occurred_at;
  if invalid_count <> 0 then raise exception 'invalid outbox rows: %', invalid_count; end if;

  select count(*) into invalid_count
  from platform.job
  where attempts > max_attempts
     or ((claimed_by is null) <> (claimed_until is null))
     or (completed_at is not null and terminal_error_code is not null);
  if invalid_count <> 0 then raise exception 'invalid job rows: %', invalid_count; end if;
end $$;

select 'platform_restore_invariants' as check_name, 'PASS' as result;
rollback;
