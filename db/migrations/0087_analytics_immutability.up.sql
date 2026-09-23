-- AUTHORED successor of GO-DATA-ANALYTICS-CORE 0.1.0 SQL contract.
-- Historical 0041 bytes remain unchanged. Runtime must not own these tables
-- or have DDL/replication privileges. Owner/superuser remains a trust boundary.
begin;

create schema analytics_guard;

create function analytics_guard.guard_landing_evidence() returns trigger
language plpgsql set search_path=pg_catalog as $$
begin
 if TG_OP='DELETE' then
  raise exception using errcode='23514',constraint='landing_evidence_immutable',message='landing evidence deletion requires owner retention procedure';
 end if;
 -- Only state is operationally mutable. All current/future other fields are
 -- immutable, including identity, payload/digest, timestamps and lineage.
 if (to_jsonb(NEW)-'state') is distinct from (to_jsonb(OLD)-'state') then
  raise exception using errcode='23514',constraint='landing_evidence_immutable',message='landing evidence fields are immutable';
 end if;
 return NEW;
end $$;
create trigger landing_evidence_immutable before update or delete on data.landing_record
 for each row execute function analytics_guard.guard_landing_evidence();

create function analytics_guard.guard_metric_definition() returns trigger
language plpgsql set search_path=pg_catalog as $$
begin
 if TG_OP='DELETE' then
  raise exception using errcode='23514',constraint='metric_definition_immutable',message='metric definition deletion requires owner retention procedure';
 end if;
 if NEW is distinct from OLD then
  raise exception using errcode='23514',constraint='metric_definition_immutable',message='metric definition is immutable; insert a new version';
 end if;
 return NEW;
end $$;
create trigger metric_definition_immutable before update or delete on analytics.metric_definition
 for each row execute function analytics_guard.guard_metric_definition();

create function analytics_guard.reject_immutable_truncate() returns trigger
language plpgsql set search_path=pg_catalog as $$
begin
 raise exception using errcode='23514',constraint='analytics_immutable_truncate',message='analytics evidence cannot be truncated by the runtime';
end $$;
create trigger landing_evidence_no_truncate before truncate on data.landing_record
 for each statement execute function analytics_guard.reject_immutable_truncate();
create trigger metric_definition_no_truncate before truncate on analytics.metric_definition
 for each statement execute function analytics_guard.reject_immutable_truncate();

comment on function analytics_guard.guard_landing_evidence() is 'AUTHORED: only landing state is mutable; controlled owner retention is documented separately';
comment on function analytics_guard.guard_metric_definition() is 'AUTHORED: definition semantics/key are immutable per version; new versions use INSERT';
commit;
