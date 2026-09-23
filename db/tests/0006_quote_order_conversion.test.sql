begin;

do $$
begin
  if not exists (select 1 from information_schema.columns where table_schema='sales' and table_name='quotation' and column_name='order_id') then
    raise exception 'quotation.order_id missing';
  end if;
  if not exists (select 1 from information_schema.tables where table_schema='sales' and table_name='quotation_acceptance') then
    raise exception 'quotation_acceptance missing';
  end if;
  if not exists (select 1 from pg_constraint where conname='quotation_acceptance_order_check') then
    raise exception 'quotation acceptance invariant missing';
  end if;
end $$;

rollback;
