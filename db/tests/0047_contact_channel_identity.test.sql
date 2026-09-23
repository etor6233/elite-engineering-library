begin;

insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)
values('47474747-4747-4747-8747-474747474747','contact-id-test','Contact Identity Test','Contact Identity Test');
insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)
values('47474747-4747-4747-8747-474747474747','store-1','store-1','Store 1','store');
insert into crm.lead(tenant_id,lead_id,organization_id,lifecycle_state,source_code,contact_payload,consent_required)
values('47474747-4747-4747-8747-474747474747','lead-1','store-1','new','meta','{}',true);
insert into communication.contact_channel_binding
(tenant_id,channel_code,external_id_hmac,lead_id,subject_id,pii_allowed,state,policy_version,evidence_sha256_hex,effective_at,version)
values('47474747-4747-4747-8747-474747474747','whatsapp',repeat('a',64),'lead-1','lead:lead-1',true,'active','policy-1',repeat('b',64),clock_timestamp(),1);
insert into communication.contact_channel_binding_decision
(tenant_id,request_id,request_sha256_hex,channel_code,external_id_hmac,binding_version,lead_id,subject_id,pii_allowed,state,policy_version,evidence_sha256_hex,effective_at)
values('47474747-4747-4747-8747-474747474747','request-1',repeat('c',64),'whatsapp',repeat('a',64),1,'lead-1','lead:lead-1',true,'active','policy-1',repeat('b',64),clock_timestamp());

do $test$
begin
  begin
    update communication.contact_channel_binding_decision set policy_version='changed'
    where tenant_id='47474747-4747-4747-8747-474747474747' and request_id='request-1';
    raise exception 'decision mutation unexpectedly succeeded';
  exception when sqlstate '55000' then null;
  end;
  if exists(select 1 from communication.contact_channel_binding where external_id_hmac like '%54911%') then
    raise exception 'raw channel identity leaked';
  end if;
end;
$test$;

rollback;
