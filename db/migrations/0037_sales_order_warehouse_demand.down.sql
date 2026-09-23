begin;

do $block$
begin
  if exists (select 1 from inventory.sales_warehouse_binding) or
     exists (select 1 from inventory.warehouse_pick_request) then
    raise exception using errcode='55000',
      message='cannot remove connected sales warehouse schema after binding or request evidence exists';
  end if;
end;
$block$;

drop trigger warehouse_pick_request_immutable on inventory.warehouse_pick_request;
drop function inventory.reject_warehouse_pick_request_mutation();
drop table inventory.warehouse_pick_request;

drop trigger sales_warehouse_binding_immutable on inventory.sales_warehouse_binding;
drop function inventory.reject_sales_warehouse_binding_mutation();
drop table inventory.sales_warehouse_binding;

commit;
