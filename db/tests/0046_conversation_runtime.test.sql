begin;

insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)
values('018f4d4a-7b36-7a21-8d10-2f4c54c20003','conversation-sql','Conversation SQL','Conversation');

insert into communication.conversation_turn
(tenant_id,channel_code,provider_message_id,external_id,thread_id,occurred_at,request_sha256_hex,state,attempt_count,locked_until,expires_at)
values('018f4d4a-7b36-7a21-8d10-2f4c54c20003','whatsapp','wamid.sql.1','contact-a','thread-a',clock_timestamp(),repeat('a',64),'processing',1,clock_timestamp()+interval '30 seconds',clock_timestamp()+interval '1 day');

update communication.conversation_turn set state='completed',user_text='consulta',assistant_text='respuesta',response_sha256_hex=repeat('b',64),locked_until=null,completed_at=clock_timestamp(),updated_at=clock_timestamp()
where tenant_id='018f4d4a-7b36-7a21-8d10-2f4c54c20003' and channel_code='whatsapp' and provider_message_id='wamid.sql.1';

do $test$
begin
  if (select count(*) from communication.conversation_turn where tenant_id='018f4d4a-7b36-7a21-8d10-2f4c54c20003' and state='completed' and assistant_text='respuesta')<>1 then
    raise exception 'completed conversation was not persisted';
  end if;
  begin
    update communication.conversation_turn set external_id='attacker'
    where tenant_id='018f4d4a-7b36-7a21-8d10-2f4c54c20003' and channel_code='whatsapp' and provider_message_id='wamid.sql.1';
    raise exception 'terminal identity mutation was accepted';
  exception when sqlstate '55000' then
    null;
  end;
end;
$test$;

rollback;
