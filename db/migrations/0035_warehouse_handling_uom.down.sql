begin;

drop trigger if exists warehouse_activity_line_uom_validate on inventory.warehouse_activity_line;
drop trigger if exists warehouse_receipt_line_uom_validate on inventory.warehouse_receipt_line;
drop function if exists inventory.validate_warehouse_line_uom();

alter table inventory.warehouse_activity_line
  drop constraint if exists warehouse_activity_line_uom_fk,
  drop constraint if exists warehouse_activity_line_uom_exact_check,
  drop constraint if exists warehouse_activity_line_qty_per_uom_check,
  drop constraint if exists warehouse_activity_line_uom_quantity_check,
  drop column if exists qty_per_uom,
  drop column if exists uom_quantity,
  drop column if exists uom_code;

alter table inventory.warehouse_receipt_line
  drop constraint if exists warehouse_receipt_line_uom_fk,
  drop constraint if exists warehouse_receipt_line_uom_exact_check,
  drop constraint if exists warehouse_receipt_line_qty_per_uom_check,
  drop constraint if exists warehouse_receipt_line_uom_quantity_check,
  drop column if exists qty_per_uom,
  drop column if exists uom_quantity,
  drop column if exists uom_code;

alter table inventory.warehouse_receipt drop column if exists allow_breakbulk;

alter table inventory.bulk_uom_conversion
  drop constraint if exists bulk_uom_conversion_source_id_check,
  drop constraint if exists bulk_uom_conversion_source_kind_check,
  drop column if exists source_id,
  drop column if exists source_kind;

alter table inventory.bulk_uom_balance
  drop constraint if exists bulk_uom_balance_reserved_quantity_check,
  drop column if exists reserved_quantity;

commit;
