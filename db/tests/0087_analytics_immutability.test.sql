-- AUTHORED, self-contained disposable SQL regression. Requires fixture database owner.
-- Role/grants/seed are transaction-local effects rolled back at the end.
begin;
create role elite_analytics_sql_runtime nologin nosuperuser nocreatedb nocreaterole noinherit noreplication;
grant usage on schema data,analytics to elite_analytics_sql_runtime;
grant select,insert,update,delete,truncate on all tables in schema data,analytics to elite_analytics_sql_runtime;
insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values
 ('018f4d4a-7b36-7a21-8d10-2f4c54c2e411','analytics-sql-a','Fixture','Fixture'),
 ('018f4d4a-7b36-7a21-8d10-2f4c54c2e412','analytics-sql-b','Fixture','Fixture');
insert into data.ingest_source(tenant_id,source_id,contract_schema,freshness_ttl)
 select tenant_id,s,'fixture/v1',interval '1 hour' from platform.tenant cross join unnest(array['fixture','second'])s where tenant_id in ('018f4d4a-7b36-7a21-8d10-2f4c54c2e411','018f4d4a-7b36-7a21-8d10-2f4c54c2e412');
insert into data.landing_record(tenant_id,source_id,external_id,event_time,received_at,payload,source_sha256,lineage)
 values('018f4d4a-7b36-7a21-8d10-2f4c54c2e411','fixture','record','2026-09-01','2026-09-02','{"value":1}',repeat('a',64),'fixture'),
 ('018f4d4a-7b36-7a21-8d10-2f4c54c2e411','fixture','expired','2025-01-01','2025-01-02','{"value":2}',repeat('c',64),'fixture');
insert into analytics.metric_definition(tenant_id,metric_code,version,aggregation,source_ref,freshness_ttl)
 values('018f4d4a-7b36-7a21-8d10-2f4c54c2e411','sales.total',1,'sum','fixture.owner',interval '1 hour');

set local role elite_analytics_sql_runtime;
do $$ begin begin update data.landing_record set payload='{"value":9}' where tenant_id='018f4d4a-7b36-7a21-8d10-2f4c54c2e411' and source_id='fixture' and external_id='record'; raise exception using errcode='P0001',message='guard missed landing-payload'; exception when check_violation then null; end; end $$;
do $$ begin begin update data.landing_record set source_sha256=repeat('b',64) where tenant_id='018f4d4a-7b36-7a21-8d10-2f4c54c2e411' and source_id='fixture' and external_id='record'; raise exception using errcode='P0001',message='guard missed landing-digest'; exception when check_violation then null; end; end $$;
do $$ begin begin update data.landing_record set external_id='replacement' where tenant_id='018f4d4a-7b36-7a21-8d10-2f4c54c2e411' and source_id='fixture' and external_id='record'; raise exception using errcode='P0001',message='guard missed landing-external-id'; exception when check_violation then null; end; end $$;
do $$ begin begin update data.landing_record set source_id='second' where tenant_id='018f4d4a-7b36-7a21-8d10-2f4c54c2e411' and source_id='fixture' and external_id='record'; raise exception using errcode='P0001',message='guard missed landing-source'; exception when check_violation then null; end; end $$;
do $$ begin begin update data.landing_record set tenant_id='018f4d4a-7b36-7a21-8d10-2f4c54c2e412' where tenant_id='018f4d4a-7b36-7a21-8d10-2f4c54c2e411' and source_id='fixture' and external_id='record'; raise exception using errcode='P0001',message='guard missed landing-tenant'; exception when check_violation then null; end; end $$;
do $$ begin begin update data.landing_record set event_time=event_time-interval '1 hour' where tenant_id='018f4d4a-7b36-7a21-8d10-2f4c54c2e411' and source_id='fixture' and external_id='record'; raise exception using errcode='P0001',message='guard missed landing-event-time'; exception when check_violation then null; end; end $$;
do $$ begin begin update data.landing_record set received_at=received_at+interval '1 hour' where tenant_id='018f4d4a-7b36-7a21-8d10-2f4c54c2e411' and source_id='fixture' and external_id='record'; raise exception using errcode='P0001',message='guard missed landing-received-at'; exception when check_violation then null; end; end $$;
do $$ begin begin update data.landing_record set lineage='replacement' where tenant_id='018f4d4a-7b36-7a21-8d10-2f4c54c2e411' and source_id='fixture' and external_id='record'; raise exception using errcode='P0001',message='guard missed landing-lineage'; exception when check_violation then null; end; end $$;
do $$ begin begin update analytics.metric_definition set aggregation='count' where tenant_id='018f4d4a-7b36-7a21-8d10-2f4c54c2e411' and metric_code='sales.total' and version=1; raise exception using errcode='P0001',message='guard missed metric-aggregation'; exception when check_violation then null; end; end $$;
do $$ begin begin update analytics.metric_definition set source_ref='another.owner' where tenant_id='018f4d4a-7b36-7a21-8d10-2f4c54c2e411' and metric_code='sales.total' and version=1; raise exception using errcode='P0001',message='guard missed metric-source'; exception when check_violation then null; end; end $$;
do $$ begin begin update analytics.metric_definition set dimensions='["organization"]' where tenant_id='018f4d4a-7b36-7a21-8d10-2f4c54c2e411' and metric_code='sales.total' and version=1; raise exception using errcode='P0001',message='guard missed metric-dimensions'; exception when check_violation then null; end; end $$;
do $$ begin begin update analytics.metric_definition set freshness_ttl=interval '2 hours' where tenant_id='018f4d4a-7b36-7a21-8d10-2f4c54c2e411' and metric_code='sales.total' and version=1; raise exception using errcode='P0001',message='guard missed metric-ttl'; exception when check_violation then null; end; end $$;
do $$ begin begin update analytics.metric_definition set version=9 where tenant_id='018f4d4a-7b36-7a21-8d10-2f4c54c2e411' and metric_code='sales.total' and version=1; raise exception using errcode='P0001',message='guard missed metric-version'; exception when check_violation then null; end; end $$;
do $$ begin begin update analytics.metric_definition set metric_code='sales.other' where tenant_id='018f4d4a-7b36-7a21-8d10-2f4c54c2e411' and metric_code='sales.total' and version=1; raise exception using errcode='P0001',message='guard missed metric-code'; exception when check_violation then null; end; end $$;
do $$ begin begin update analytics.metric_definition set tenant_id='018f4d4a-7b36-7a21-8d10-2f4c54c2e412' where tenant_id='018f4d4a-7b36-7a21-8d10-2f4c54c2e411' and metric_code='sales.total' and version=1; raise exception using errcode='P0001',message='guard missed metric-tenant'; exception when check_violation then null; end; end $$;
do $$ begin begin update analytics.metric_definition set created_at=created_at+interval '1 hour' where tenant_id='018f4d4a-7b36-7a21-8d10-2f4c54c2e411' and metric_code='sales.total' and version=1; raise exception using errcode='P0001',message='guard missed metric-created'; exception when check_violation then null; end; end $$;
do $$ begin begin delete from data.landing_record where tenant_id='018f4d4a-7b36-7a21-8d10-2f4c54c2e411' and source_id='fixture' and external_id='record'; raise exception using errcode='P0001',message='guard missed landing-delete'; exception when check_violation then null; end; end $$;
do $$ begin begin delete from analytics.metric_definition where tenant_id='018f4d4a-7b36-7a21-8d10-2f4c54c2e411' and metric_code='sales.total' and version=1; raise exception using errcode='P0001',message='guard missed metric-delete'; exception when check_violation then null; end; end $$;
do $$ begin begin truncate data.landing_record; raise exception using errcode='P0001',message='guard missed landing-truncate'; exception when check_violation then null; end; end $$;
do $$ begin begin truncate analytics.metric_definition cascade; raise exception using errcode='P0001',message='guard missed metric-truncate'; exception when check_violation then null; end; end $$;
rollback;
