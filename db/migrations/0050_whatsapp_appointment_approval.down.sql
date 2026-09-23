begin;
-- Destructive schema rollback. Export/retain evidence and stop the consumer
-- before using this on an approved target; not a retry mechanism.
drop trigger whatsapp_approval_immutable on communication.whatsapp_appointment_approval;
drop function communication.reject_whatsapp_approval_mutation();
drop table communication.whatsapp_appointment_approval;
drop index crm.whatsapp_notification_consent_lookup_idx;
commit;
