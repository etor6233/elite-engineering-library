begin;

drop table if exists analytics.metric_snapshot;
drop table if exists analytics.metric_definition;
drop table if exists data.landing_record;
drop table if exists data.ingest_source;
drop schema if exists analytics;
drop schema if exists data;

commit;
