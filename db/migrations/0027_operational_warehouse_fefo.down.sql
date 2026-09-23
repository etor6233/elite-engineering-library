begin;
drop trigger if exists warehouse_activity_transition_guard on inventory.warehouse_activity;
drop function if exists inventory.enforce_warehouse_activity_transition();
drop trigger if exists warehouse_activity_line_immutable on inventory.warehouse_activity_line;
drop trigger if exists warehouse_receipt_line_immutable on inventory.warehouse_receipt_line;
drop trigger if exists warehouse_receipt_immutable on inventory.warehouse_receipt;
drop function if exists inventory.reject_warehouse_evidence_mutation();
drop table if exists inventory.warehouse_activity_line;
drop table if exists inventory.warehouse_activity;
drop table if exists inventory.warehouse_receipt_line;
drop table if exists inventory.warehouse_receipt;
alter table inventory.warehouse_bin drop constraint warehouse_bin_bin_type_check;
alter table inventory.warehouse_bin drop column bin_rank;
alter table inventory.warehouse_bin add constraint warehouse_bin_bin_type_check
  check (bin_type in ('receive','ship','put-away','pick','putpick'));
commit;
