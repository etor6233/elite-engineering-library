-- 0041_data_analytics.test.sql — verifica replay idempotente, rechazo de
-- received-before-event, definición versionada de métrica y snapshot.
-- Precondición: migración 0001 (platform.tenant) y 0041 aplicadas.

begin;

insert into platform.tenant (tenant_id, tenant_code, legal_name, display_name) values
  ('11111111-1111-1111-1111-111111111111', 'tenant-a', 'A', 'Tenant A');

insert into data.ingest_source (tenant_id, source_id, contract_schema, freshness_ttl) values
  ('11111111-1111-1111-1111-111111111111', 'pos', 'pos-events/v1', interval '1 hour');

insert into data.landing_record
  (tenant_id, source_id, external_id, event_time, received_at, payload, source_sha256)
values
  ('11111111-1111-1111-1111-111111111111', 'pos', 'rec-1', now() - interval '1 minute', now(),
   '{"total":10}', repeat('a',64));

insert into data.landing_record
  (tenant_id, source_id, external_id, event_time, received_at, payload, source_sha256)
values
  ('11111111-1111-1111-1111-111111111111', 'pos', 'rec-1', now() - interval '1 minute', now(),
   '{"total":10}', repeat('a',64))
on conflict do nothing;

-- replay: exactly one row despite the duplicate insert
do $$
declare n int;
begin
  select count(*) into n from data.landing_record
   where tenant_id = '11111111-1111-1111-1111-111111111111' and external_id = 'rec-1';
  if n <> 1 then raise exception 'replay dedup expected 1, got %', n; end if;
end $$;

-- ordering: received_at < event_time must violate the check constraint
do $$
begin
  begin
    insert into data.landing_record
      (tenant_id, source_id, external_id, event_time, received_at, payload, source_sha256)
    values
      ('11111111-1111-1111-1111-111111111111', 'pos', 'rec-2', now(), now() - interval '1 minute',
       '{"x":1}', repeat('b',64));
    raise exception 'received-before-event was accepted';
  exception when check_violation then
    null; -- expected
  end;
end $$;

-- versioned metric definition + snapshot
insert into analytics.metric_definition
  (tenant_id, metric_code, version, aggregation, source_ref, dimensions, freshness_ttl)
values
  ('11111111-1111-1111-1111-111111111111', 'sales.revenue', 1, 'sum', 'sales.customer_order',
   '["organization","month"]', interval '1 hour');

insert into analytics.metric_snapshot
  (tenant_id, metric_code, version, period_start, period_end, value, dimensions, reconciled)
values
  ('11111111-1111-1111-1111-111111111111', 'sales.revenue', 1,
   '2026-09-01 00:00:00+00', '2026-09-02 00:00:00+00', 1234.56,
   '{"organization":"org-1"}', true);

do $$
declare n int;
begin
  select count(*) into n from analytics.metric_snapshot
   where tenant_id = '11111111-1111-1111-1111-111111111111';
  if n <> 1 then raise exception 'snapshot expected 1, got %', n; end if;
end $$;

rollback;
