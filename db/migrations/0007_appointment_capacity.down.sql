begin;
drop index if exists crm.appointment_slot_booking_idx;
alter table crm.appointment drop constraint if exists appointment_slot_time_check, drop constraint if exists appointment_slot_fk, drop column if exists ends_at, drop column if exists slot_id;
drop table if exists crm.appointment_slot;
drop function if exists crm.enforce_appointment_slot_non_overlap();
commit;
