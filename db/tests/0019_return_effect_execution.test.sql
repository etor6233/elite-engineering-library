begin;

do $$
begin
  if to_regclass('sales.return_effect_execution') is null or to_regclass('sales.return_effect_attempt') is null or to_regclass('sales.return_effect_resume') is null then
    raise exception 'return effect execution tables missing';
  end if;
  if (select count(*) from pg_trigger where tgname in ('return_effect_execution_enqueue','return_effect_attempt_immutable','return_effect_resume_immutable') and not tgisinternal) <> 3 then
    raise exception 'return effect execution trigger contract incomplete';
  end if;
  if not exists (select 1 from pg_indexes where schemaname='sales' and indexname='return_effect_execution_claim_idx') then
    raise exception 'return effect claim index missing';
  end if;
  if exists (select 1 from sales.return_effect_request r left join sales.return_effect_execution e on e.tenant_id=r.tenant_id and e.request_id=r.request_id where e.request_id is null) then
    raise exception 'return effect request missing execution projection';
  end if;
end $$;

rollback;
