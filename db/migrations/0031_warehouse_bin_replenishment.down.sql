begin;

do $migration$
begin
  if exists (select 1 from inventory.warehouse_replenishment_request)
     or exists (select 1 from inventory.warehouse_activity where activity_type='movement') then
    raise exception using errcode='55000', message='cannot remove replenishment schema after movement evidence exists';
  end if;
end;
$migration$;

drop table inventory.warehouse_replenishment_request;

alter table inventory.bulk_reservation
  drop constraint bulk_reservation_demand_kind_v165_check,
  add constraint bulk_reservation_demand_kind_check
    check (demand_kind in ('customer-order','service','transfer-outbound','manual'));

alter table inventory.warehouse_activity
  drop constraint warehouse_activity_activity_type_v165_check,
  add constraint warehouse_activity_activity_type_check
    check (activity_type in ('put-away','pick'));

commit;
