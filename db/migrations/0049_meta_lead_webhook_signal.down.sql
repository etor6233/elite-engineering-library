begin;
drop trigger if exists meta_lead_webhook_signal_immutable on integration.meta_lead_webhook_signal;
drop trigger if exists meta_lead_webhook_batch_immutable on integration.meta_lead_webhook_batch;
drop table if exists integration.meta_lead_webhook_signal;
drop table if exists integration.meta_lead_webhook_batch;
drop function if exists integration.reject_meta_lead_webhook_mutation();
commit;
