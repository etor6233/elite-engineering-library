begin;
-- AUTHORED read-model indexes; no new business state or send authority.
create index whatsapp_campaign_workspace_idx on communication.whatsapp_campaign(tenant_id,organization_id,campaign_id);
create index campaign_lead_workspace_idx on crm.lead(tenant_id,organization_id,lead_id) where lifecycle_state in('new','contacted','qualified');
commit;
