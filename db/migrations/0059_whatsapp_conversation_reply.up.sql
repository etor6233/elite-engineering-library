begin;
-- AUTHORED extension of the existing status owner: every status batch anchors
-- an already accepted outbound row, whether its approval was appointment or
-- an exact human-reviewed conversational response.
alter table communication.whatsapp_status_batch drop constraint whatsapp_status_batch_tenant_id_delivery_key_fkey;
alter table communication.whatsapp_status_batch add column channel_code text not null default 'whatsapp' check(channel_code='whatsapp');
alter table communication.whatsapp_status_batch add constraint whatsapp_status_batch_outbound_fkey foreign key(tenant_id,channel_code,delivery_key) references communication.outbound_delivery(tenant_id,channel_code,delivery_key);

create view communication.whatsapp_delivery_approval as
 select g.tenant_id,g.channel_code,g.delivery_key,g.organization_id,g.appointment_id,g.confirmation_event_id::text as approval_event_id,g.lead_id,g.external_id_hmac,g.profile_sha256,g.message_sha256,g.expires_at
 from communication.whatsapp_appointment_approval g
 union all
 select r.tenant_id,'whatsapp',r.request_id,r.organization_id,''::text,r.request_id,
 r.payload->>'lead_id',r.payload->>'external_id_hmac',r.payload->>'profile_sha256',r.payload->>'message_sha256',(r.payload->>'expires_at')::timestamptz
 from approval.request r
 where r.kind='whatsapp_reply' and r.state='approved' and r.payload->>'schema'='elite-whatsapp-reply-approval/v1'
 and r.payload->>'organization_id'=r.organization_id;
commit;
