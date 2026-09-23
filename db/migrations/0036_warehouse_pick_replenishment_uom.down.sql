begin;

do $block$
begin
  if exists (select 1 from inventory.warehouse_activity_uom_reservation) or
     exists (
       select 1 from inventory.warehouse_activity_line
        where from_uom_code<>uom_code or
              from_uom_quantity<>uom_quantity or
              from_qty_per_uom<>qty_per_uom
     ) then
    raise exception using errcode='55000',
      message='cannot remove packaging-aware warehouse schema after evidence exists';
  end if;
end;
$block$;

drop trigger warehouse_uom_reservation_delete_guard on inventory.warehouse_activity_uom_reservation;
drop function inventory.reject_warehouse_uom_reservation_delete();
drop trigger warehouse_uom_reservation_transition_guard on inventory.warehouse_activity_uom_reservation;
drop function inventory.enforce_warehouse_uom_reservation_transition();
drop table inventory.warehouse_activity_uom_reservation;

alter table inventory.warehouse_replenishment_request
  drop constraint warehouse_replenishment_target_uom_fk,
  drop column allow_breakbulk,
  drop column target_uom;

drop trigger warehouse_activity_line_from_uom_validate on inventory.warehouse_activity_line;
drop function inventory.validate_warehouse_line_from_uom();
alter table inventory.warehouse_activity_line
  drop constraint warehouse_activity_line_from_uom_fk,
  drop constraint warehouse_activity_line_from_qty_per_uom_check,
  drop constraint warehouse_activity_line_from_uom_quantity_check,
  drop column from_qty_per_uom,
  drop column from_uom_quantity,
  drop column from_uom_code;

commit;
