begin;
do $$ begin
  if to_regclass('accounting.return_effect_posting') is null then raise exception 'return accounting table missing'; end if;
  if (select count(*) from pg_trigger where tgname='return_effect_posting_immutable' and not tgisinternal)<>1 then raise exception 'return accounting immutability missing'; end if;
  if not exists(select 1 from pg_indexes where schemaname='accounting' and indexname='return_effect_posting_source_idx') then raise exception 'return accounting index missing'; end if;
end $$;
rollback;
