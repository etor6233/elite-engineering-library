begin;
insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name,status) values('018f4d4a-7b36-7a21-8d10-2f4c54c29f12','fiscal-v133','Fiscal V133 SA','Fiscal V133','active');
insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type,status) values('018f4d4a-7b36-7a21-8d10-2f4c54c29f12','org-v133','org-v133','Organization V133','franchisee','active');
insert into fiscal.parameter_refresh_schedule(tenant_id,schedule_id,organization_id,taxpayer_cuit,parameter_kind,interval_seconds,next_run_at,active,version,configured_by_subject) values('018f4d4a-7b36-7a21-8d10-2f4c54c29f12','schedule-v133','org-v133','33693450239','vat_rate',900,clock_timestamp(),true,1,'controller');
do $$ begin
  if (select count(*) from fiscal.parameter_refresh_schedule where tenant_id='018f4d4a-7b36-7a21-8d10-2f4c54c29f12' and active and next_run_at<=clock_timestamp())<>1 then raise exception 'due schedule missing'; end if;
  begin
    insert into fiscal.parameter_refresh_schedule(tenant_id,schedule_id,organization_id,taxpayer_cuit,parameter_kind,interval_seconds,next_run_at,active,version,configured_by_subject) values('018f4d4a-7b36-7a21-8d10-2f4c54c29f12','bad-v133','org-v133','33693450239','vat_rate',899,clock_timestamp(),true,1,'controller');
    raise exception 'unsafe cadence accepted';
  exception when check_violation then null; end;
end $$;
rollback;
