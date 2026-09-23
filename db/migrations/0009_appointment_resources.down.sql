begin;
drop index if exists crm.appointment_resource_schedule_idx;
drop index if exists crm.service_resource_org_status_idx;
drop trigger if exists appointment_resource_assignment_guard on crm.appointment_resource;
drop function if exists crm.enforce_appointment_resource_assignment();
drop table if exists crm.appointment_resource;
drop table if exists crm.resource_skill;
drop table if exists crm.service_resource;
commit;
