begin;

insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)
values('44444444-4444-4444-8444-444444444444','lead-ingress-test','Lead Ingress Test','Lead Ingress Test');
insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)
values('44444444-4444-4444-8444-444444444444','store-1','store-1','Store 1','store');

insert into integration.lead_ingress_raw
(tenant_id,organization_id,provider,provider_event_id,event_type,source,schema_version,occurred_at,received_at,payload_redacted,source_payload_sha256,stored_payload_sha256,redaction_profile,state)
values('44444444-4444-4444-8444-444444444444','store-1','google_ads','lead-1','google.ads.lead.v3','https://googleads.googleapis.com/lead-form','3',clock_timestamp(),clock_timestamp(),convert_to('{"lead_id":"lead-1"}','UTF8'),repeat('a',64),repeat('b',64),'google_ads:remove-google_key:v1','normalized');

insert into integration.lead_candidate
(tenant_id,provider,provider_event_id,organization_id,provider_lead_id,submitted_at,is_test,fields,contact_eligibility)
values('44444444-4444-4444-8444-444444444444','google_ads','lead-1','store-1','lead-1',clock_timestamp(),false,'[]','pending_policy');

do $test$
begin
  begin
    update integration.lead_ingress_raw set state='rejected'
    where tenant_id='44444444-4444-4444-8444-444444444444' and provider='google_ads' and provider_event_id='lead-1';
    raise exception 'raw update unexpectedly succeeded';
  exception when sqlstate '55000' then null;
  end;
  begin
    update integration.lead_candidate set provider_lead_id='changed'
    where tenant_id='44444444-4444-4444-8444-444444444444' and provider='google_ads' and provider_event_id='lead-1';
    raise exception 'candidate update unexpectedly succeeded';
  exception when sqlstate '55000' then null;
  end;
end;
$test$;

rollback;
