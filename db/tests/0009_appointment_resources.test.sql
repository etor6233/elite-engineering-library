begin;
do $$
begin
  if to_regclass('crm.service_resource') is null then raise exception 'service_resource missing'; end if;
  if to_regclass('crm.resource_skill') is null then raise exception 'resource_skill missing'; end if;
  if to_regclass('crm.appointment_resource') is null then raise exception 'appointment_resource missing'; end if;
  if not exists(select 1 from pg_indexes where schemaname='crm' and indexname='service_resource_principal_idx' and indexdef like '%WHERE (principal_subject IS NOT NULL)%') then raise exception 'partial principal uniqueness missing'; end if;
  if to_regprocedure('crm.enforce_appointment_resource_assignment()') is null then raise exception 'assignment guard missing'; end if;
end;
$$;
rollback;
