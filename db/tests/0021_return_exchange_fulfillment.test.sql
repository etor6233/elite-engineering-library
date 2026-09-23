begin;

do $$
begin
  if to_regclass('sales.return_exchange') is null then
    raise exception 'return exchange table missing';
  end if;
  if (select count(*) from pg_trigger where tgname='return_exchange_immutable' and not tgisinternal) <> 1 then
    raise exception 'return exchange immutability trigger missing';
  end if;
  if not exists (select 1 from pg_indexes where schemaname='sales' and indexname='return_exchange_original_order_idx') then
    raise exception 'return exchange lookup index missing';
  end if;
  if exists (select 1 from sales.return_exchange where original_order_id=replacement_order_id or original_stock_unit_id=replacement_stock_unit_id) then
    raise exception 'return exchange reused original order or stock';
  end if;
end $$;

rollback;
