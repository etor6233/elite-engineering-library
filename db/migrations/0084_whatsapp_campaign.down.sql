begin;
do $$ begin
 if exists(select 1 from communication.whatsapp_campaign) or exists(select 1 from communication.whatsapp_schedule where campaign_id is not null)
 then raise exception 'cannot remove campaign evidence';end if;
end $$;
alter table communication.whatsapp_schedule drop constraint whatsapp_schedule_campaign_member_fk,
 drop constraint whatsapp_schedule_source_check,drop constraint whatsapp_schedule_campaign_step_unique,
 drop column campaign_id,drop column lead_id,drop column campaign_step;
alter table communication.whatsapp_schedule alter column appointment_id set not null;
drop table communication.whatsapp_campaign_stop;
drop table communication.whatsapp_campaign_member;
drop table communication.whatsapp_campaign;
commit;
