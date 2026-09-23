begin;

insert into platform.tenant (
  tenant_id,
  tenant_code,
  legal_name,
  display_name
) values (
  '018f4d4a-7b36-7a21-8d10-2f4c54c28a01',
  'test-tenant',
  'Test Tenant S.A.',
  'Test Tenant'
);

insert into platform.business_profile_version (
  tenant_id,
  profile_version,
  schema_version,
  profile,
  sha256_hex,
  status,
  created_by_subject,
  activated_at
) values (
  '018f4d4a-7b36-7a21-8d10-2f4c54c28a01',
  1,
  '1.0.0',
  '{"schemaVersion":"1.0.0"}'::jsonb,
  repeat('a', 64),
  'active',
  'test:operator',
  clock_timestamp()
);

do $test$
begin
  begin
    insert into platform.business_profile_version (
      tenant_id, profile_version, schema_version, profile,
      sha256_hex, status, created_by_subject, activated_at
    ) values (
      '018f4d4a-7b36-7a21-8d10-2f4c54c28a01', 2, '1.0.0',
      '{"schemaVersion":"1.0.0"}'::jsonb, repeat('b', 64),
      'active', 'test:operator', clock_timestamp()
    );
    raise exception 'expected one-active-profile unique violation';
  exception
    when unique_violation then null;
  end;
end;
$test$;

insert into platform.idempotency_record (
  tenant_id,
  scope,
  idempotency_key,
  request_sha256_hex,
  status,
  locked_until,
  expires_at
) values (
  '018f4d4a-7b36-7a21-8d10-2f4c54c28a01',
  'order:create',
  'idem-test-0001',
  repeat('c', 64),
  'processing',
  clock_timestamp() + interval '30 seconds',
  clock_timestamp() + interval '24 hours'
);

insert into platform.outbox_event (
  tenant_id,
  event_id,
  aggregate_type,
  aggregate_id,
  aggregate_version,
  event_type,
  schema_version,
  occurred_at,
  payload
) values (
  '018f4d4a-7b36-7a21-8d10-2f4c54c28a01',
  '018f4d4a-7b36-7a21-8d10-2f4c54c28a10',
  'order',
  'order-test-1',
  1,
  'order.created',
  1,
  clock_timestamp(),
  '{"orderId":"order-test-1"}'::jsonb
);

insert into audit.event (
  tenant_id,
  event_id,
  actor_subject,
  action,
  resource_type,
  resource_id,
  decision
) values (
  '018f4d4a-7b36-7a21-8d10-2f4c54c28a01',
  '018f4d4a-7b36-7a21-8d10-2f4c54c28a20',
  'test:operator',
  'business-profile.activate',
  'business-profile',
  '1',
  'allowed'
);

do $test$
begin
  begin
    update audit.event
       set action = 'tampered'
     where event_id = '018f4d4a-7b36-7a21-8d10-2f4c54c28a20';
    raise exception 'expected append-only mutation rejection';
  exception
    when object_not_in_prerequisite_state then null;
  end;
end;
$test$;

do $test$
declare
  actual_count bigint;
begin
  select count(*) into actual_count
    from platform.outbox_event
   where tenant_id = '018f4d4a-7b36-7a21-8d10-2f4c54c28a01';

  if actual_count <> 1 then
    raise exception 'expected exactly one outbox event, got %', actual_count;
  end if;
end;
$test$;

rollback;
