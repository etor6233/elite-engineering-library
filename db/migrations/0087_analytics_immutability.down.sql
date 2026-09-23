-- Owner-only rollback of the additive guards, preserving all existing data.
-- Also supports cleanup after an explicit historical 0041 DROP rollback.
begin;
do $$ begin
 if to_regclass('analytics.metric_definition') is not null then
  execute 'drop trigger if exists metric_definition_no_truncate on analytics.metric_definition';
  execute 'drop trigger if exists metric_definition_immutable on analytics.metric_definition';
 end if;
 if to_regclass('data.landing_record') is not null then
  execute 'drop trigger if exists landing_evidence_no_truncate on data.landing_record';
  execute 'drop trigger if exists landing_evidence_immutable on data.landing_record';
 end if;
end $$;
drop function if exists analytics_guard.reject_immutable_truncate();
drop function if exists analytics_guard.guard_metric_definition();
drop function if exists analytics_guard.guard_landing_evidence();
drop schema if exists analytics_guard;
commit;
