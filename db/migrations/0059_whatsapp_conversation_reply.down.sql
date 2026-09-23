begin;
-- Removing this owner while conversational status evidence exists would sever
-- its FK. Fail closed; export/retain evidence before a deliberate downgrade.
do $$ begin if exists(select 1 from communication.whatsapp_status_batch b where not exists(select 1 from communication.whatsapp_appointment_approval a where a.tenant_id=b.tenant_id and a.delivery_key=b.delivery_key)) then raise exception 'conversational status evidence prevents downgrade'; end if; end $$;
drop view communication.whatsapp_delivery_approval;
alter table communication.whatsapp_status_batch drop constraint whatsapp_status_batch_outbound_fkey;
alter table communication.whatsapp_status_batch drop column channel_code;
alter table communication.whatsapp_status_batch add foreign key(tenant_id,delivery_key) references communication.whatsapp_appointment_approval(tenant_id,delivery_key);
commit;
