begin;
drop trigger if exists delivery_checklist_completion_guard on sales.delivery_handover;
drop function if exists sales.enforce_delivery_checklist_completion();
drop trigger if exists delivery_checklist_response_immutable on sales.delivery_checklist_response;
drop function if exists sales.prevent_delivery_checklist_response_mutation();
drop trigger if exists delivery_checklist_response_guard on sales.delivery_checklist_response;
drop function if exists sales.enforce_delivery_checklist_response();
drop table if exists sales.delivery_checklist_response;
alter table sales.delivery_handover
  drop constraint if exists delivery_handover_checklist_state_ck,
  drop constraint if exists delivery_handover_checklist_fk,
  drop column if exists checklist_completed_by_subject,
  drop column if exists checklist_completed_at,
  drop column if exists checklist_version,
  drop column if exists checklist_id;
drop trigger if exists delivery_checklist_item_immutable on sales.delivery_checklist_item;
drop table if exists sales.delivery_checklist_item;
drop trigger if exists delivery_checklist_template_immutable on sales.delivery_checklist_template;
drop function if exists sales.prevent_published_delivery_checklist_mutation();
drop table if exists sales.delivery_checklist_template;
commit;
