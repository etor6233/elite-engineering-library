begin;

insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)
values('49494949-4949-4949-8949-494949494949','meta-webhook-test','Meta Webhook Test','Meta Webhook Test');
insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)
values('49494949-4949-4949-8949-494949494949','store-1','store-1','Store 1','store');

insert into integration.meta_lead_webhook_batch
(tenant_id,organization_id,payload_sha256,payload,received_at,signal_count)
values('49494949-4949-4949-8949-494949494949','store-1',repeat('a',64),convert_to('{"object":"page"}','UTF8'),clock_timestamp(),1);

insert into integration.meta_lead_webhook_signal
(tenant_id,leadgen_id,organization_id,batch_payload_sha256,value_sha256,form_id,page_id,occurred_at,state)
values('49494949-4949-4949-8949-494949494949','lead-1','store-1',repeat('a',64),repeat('b',64),'form-1','page-1',clock_timestamp(),'pending_retrieval');

do $test$
begin
  begin
    update integration.meta_lead_webhook_signal set form_id='changed'
    where tenant_id='49494949-4949-4949-8949-494949494949' and leadgen_id='lead-1';
    raise exception 'signal update unexpectedly succeeded';
  exception when sqlstate '55000' then null;
  end;
  begin
    delete from integration.meta_lead_webhook_batch
    where tenant_id='49494949-4949-4949-8949-494949494949' and payload_sha256=repeat('a',64);
    raise exception 'batch delete unexpectedly succeeded';
  exception when sqlstate '55000' then null;
  end;
end;
$test$;

rollback;
