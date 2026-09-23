begin;
-- AUTHORED lookup support; not UNIQUE: ambiguity must be rejected, not erased.
create index whatsapp_status_message_route_idx
on communication.outbound_delivery(tenant_id,provider_message_hmac,recipient_hmac)
where channel_code='whatsapp' and state='accepted';
commit;
