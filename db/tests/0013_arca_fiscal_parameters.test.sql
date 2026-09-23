begin;
insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values('018f4d4a-7b36-7a21-8d10-2f4c54c29f13','fiscal-param','Fiscal Param','Fiscal Param');
insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values('018f4d4a-7b36-7a21-8d10-2f4c54c29f13','franchise','franchise-param','Franchise','franchisee');
insert into fiscal.parameter_snapshot(tenant_id,snapshot_id,organization_id,taxpayer_cuit,parameter_kind,response_hash,provider_codes,fetched_at) values('018f4d4a-7b36-7a21-8d10-2f4c54c29f13','snapshot','franchise','33693450239','vat_rate',repeat('a',64),'[]','2026-08-30T12:00:00Z');
insert into fiscal.parameter_item(tenant_id,snapshot_id,parameter_code,description,valid_from) values('018f4d4a-7b36-7a21-8d10-2f4c54c29f13','snapshot','5','21%','2009-02-20');
do $$ begin
  if exists(select 1 from fiscal.approved_parameter_snapshot where snapshot_id='snapshot') then raise exception 'unapproved snapshot exposed'; end if;
  begin update fiscal.parameter_item set description='invented' where snapshot_id='snapshot'; raise exception 'parameter history mutation accepted'; exception when raise_exception then if sqlerrm <> 'immutable fiscal parameter history' then raise; end if; end;
end $$;
insert into fiscal.parameter_decision(tenant_id,snapshot_id,decision,decided_by_subject,decision_reason) values('018f4d4a-7b36-7a21-8d10-2f4c54c29f13','snapshot','approved','tax-owner','validated against current ARCA configuration');
do $$ begin
  if (select count(*) from fiscal.approved_parameter_snapshot where snapshot_id='snapshot')<>1 then raise exception 'approved snapshot absent'; end if;
  begin insert into fiscal.parameter_decision(tenant_id,snapshot_id,decision,decided_by_subject,decision_reason) values('018f4d4a-7b36-7a21-8d10-2f4c54c29f13','snapshot','rejected','other','separate rejected decision'); raise exception 'second decision accepted'; exception when unique_violation then null; end;
end $$;
rollback;
