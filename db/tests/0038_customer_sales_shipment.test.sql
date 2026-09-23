do $$
begin
  if to_regclass('sales.customer_shipment') is null then raise exception 'customer shipment missing'; end if;
  if to_regclass('sales.customer_shipment_line') is null then raise exception 'customer shipment line missing'; end if;
  if to_regclass('inventory.customer_shipment_allocation') is null then raise exception 'customer shipment allocation missing'; end if;
  if not exists(select 1 from information_schema.columns where table_schema='sales' and table_name='customer_order' and column_name='fulfillment_state') then raise exception 'order fulfillment state missing'; end if;
  if not exists(select 1 from information_schema.columns where table_schema='sales' and table_name='customer_order_line' and column_name='shipped_quantity') then raise exception 'order shipped quantity missing'; end if;
  if not exists(select 1 from pg_trigger where tgname='customer_shipment_immutable') then raise exception 'customer shipment immutability missing'; end if;
  if not exists(select 1 from pg_trigger where tgname='customer_shipment_line_immutable') then raise exception 'customer shipment line immutability missing'; end if;
  if not exists(select 1 from pg_trigger where tgname='customer_shipment_allocation_immutable') then raise exception 'customer shipment allocation immutability missing'; end if;
  if not exists(select 1 from pg_indexes where schemaname='sales' and indexname='customer_shipment_order_idx') then raise exception 'customer shipment order index missing'; end if;
end $$;
