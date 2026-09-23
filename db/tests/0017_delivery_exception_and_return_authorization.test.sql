begin;
do $$ begin
  if to_regclass('sales.delivery_exception') is null then raise exception 'delivery exception missing'; end if;
  if to_regclass('sales.delivery_exception_resolution') is null then raise exception 'delivery exception resolution missing'; end if;
  if to_regclass('sales.return_authorization') is null then raise exception 'return authorization missing'; end if;
  if not exists(select 1 from pg_trigger where tgname='delivery_exception_lifecycle') then raise exception 'delivery exception lifecycle guard missing'; end if;
  if not exists(select 1 from pg_trigger where tgname='delivery_exception_resolution_immutable') then raise exception 'resolution immutability missing'; end if;
  if not exists(select 1 from pg_indexes where schemaname='sales' and indexname='delivery_handover_one_successor_uq') then raise exception 'one successor invariant missing'; end if;
end $$;
rollback;
