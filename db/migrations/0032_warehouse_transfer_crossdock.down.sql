begin;

do $migration$
begin
  if exists (select 1 from inventory.warehouse_crossdock_policy)
     or exists (select 1 from inventory.warehouse_crossdock_allocation)
     or exists (select 1 from inventory.warehouse_crossdock_pick_link)
     or exists (select 1 from inventory.warehouse_bin where cross_dock) then
    raise exception using errcode='55000', message='cannot remove cross-dock schema after configuration or evidence exists';
  end if;
end;
$migration$;

drop table inventory.warehouse_crossdock_pick_link;
drop table inventory.warehouse_crossdock_allocation;
drop table inventory.warehouse_crossdock_policy;

alter table inventory.warehouse_bin
  drop constraint warehouse_bin_cross_dock_type_check,
  drop column cross_dock;

commit;
