begin;

do $$
begin
  if to_regclass('sales.return_receipt') is null or to_regclass('sales.return_disposition') is null or to_regclass('sales.return_effect_request') is null then
    raise exception 'return processing tables missing';
  end if;
  if (select count(*) from pg_trigger where tgname in ('return_receipt_validate','return_disposition_validate','return_receipt_immutable','return_disposition_immutable','return_effect_request_immutable') and not tgisinternal) <> 5 then
    raise exception 'return processing trigger contract incomplete';
  end if;
  if not exists (select 1 from pg_indexes where schemaname='sales' and indexname='return_effect_request_state_idx') then
    raise exception 'return effect owner queue index missing';
  end if;
end $$;

rollback;
