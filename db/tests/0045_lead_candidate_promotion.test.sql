begin;

insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)
values('45454545-4545-4545-8545-454545454545','lead-promotion-test','Lead Promotion Test','Lead Promotion Test');
insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)
values('45454545-4545-4545-8545-454545454545','store-1','store-1','Store 1','store');
insert into integration.lead_ingress_raw
(tenant_id,organization_id,provider,provider_event_id,event_type,source,schema_version,occurred_at,received_at,payload_redacted,source_payload_sha256,stored_payload_sha256,redaction_profile,state)
values('45454545-4545-4545-8545-454545454545','store-1','google_ads','lead-1','google.ads.lead.v3','https://googleads.googleapis.com/lead-form','3',clock_timestamp(),clock_timestamp(),convert_to('{"lead_id":"lead-1"}','UTF8'),repeat('a',64),repeat('b',64),'google_ads:remove-google_key:v1','normalized');
insert into integration.lead_candidate
(tenant_id,provider,provider_event_id,organization_id,provider_lead_id,submitted_at,is_test,fields,contact_eligibility)
values('45454545-4545-4545-8545-454545454545','google_ads','lead-1','store-1','lead-1',clock_timestamp(),false,'[{"id":"EMAIL","value":"person@example.test"}]','pending_policy');
insert into crm.lead(tenant_id,lead_id,organization_id,lifecycle_state,source_code,contact_payload,consent_required)
values('45454545-4545-4545-8545-454545454545','crm-lead-1','store-1','new','google_ads','{"email":"person@example.test"}',true);
insert into crm.consent_evidence(tenant_id,consent_id,lead_id,purpose_code,policy_version,decision,occurred_at,evidence_sha256_hex)
values('45454545-4545-4545-8545-454545454545','consent-1','crm-lead-1','sales-contact','policy-7','granted',clock_timestamp(),repeat('c',64));
insert into integration.lead_promotion
(tenant_id,provider,provider_event_id,lead_id,consent_id,mapping_version,field_mapping,source_payload_sha256,decision_at)
values('45454545-4545-4545-8545-454545454545','google_ads','lead-1','crm-lead-1','consent-1','google-form-v1','{"EMAIL":"email"}',repeat('a',64),clock_timestamp());

do $test$
begin
  begin
    update integration.lead_promotion set mapping_version='changed'
    where tenant_id='45454545-4545-4545-8545-454545454545' and provider='google_ads' and provider_event_id='lead-1';
    raise exception 'promotion mutation unexpectedly succeeded';
  exception when sqlstate '55000' then null;
  end;
  begin
    delete from integration.lead_promotion
    where tenant_id='45454545-4545-4545-8545-454545454545' and provider='google_ads' and provider_event_id='lead-1';
    raise exception 'promotion delete unexpectedly succeeded';
  exception when sqlstate '55000' then null;
  end;
end;
$test$;

rollback;
