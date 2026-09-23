begin;
insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)
values('48484848-4848-4848-8848-484848484848','outbound-test','Outbound Test','Outbound Test');
insert into communication.outbound_delivery
(tenant_id,channel_code,delivery_key,request_sha256_hex,recipient_hmac,state,attempt_count,locked_until)
values('48484848-4848-4848-8848-484848484848','whatsapp',repeat('d',64),repeat('a',64),repeat('b',64),'sending',1,clock_timestamp()+interval '1 minute');
insert into communication.outbound_delivery_event(tenant_id,channel_code,delivery_key,sequence,state)
values('48484848-4848-4848-8848-484848484848','whatsapp',repeat('d',64),1,'sending');
do $test$
begin
  begin
    update communication.outbound_delivery_event set state='accepted'
    where tenant_id='48484848-4848-4848-8848-484848484848';
    raise exception 'delivery event mutation unexpectedly succeeded';
  exception when sqlstate '55000' then null;
  end;
  if exists(select 1 from communication.outbound_delivery where attempt_count<>1) then
    raise exception 'automatic retry became possible';
  end if;
end;
$test$;
rollback;
