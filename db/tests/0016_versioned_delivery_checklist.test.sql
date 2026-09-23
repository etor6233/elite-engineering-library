begin;
do $$ begin
  if to_regclass('sales.delivery_checklist_template') is null then raise exception 'delivery checklist template missing'; end if;
  if to_regclass('sales.delivery_checklist_item') is null then raise exception 'delivery checklist item missing'; end if;
  if to_regclass('sales.delivery_checklist_response') is null then raise exception 'delivery checklist response missing'; end if;
  if not exists(select 1 from pg_trigger where tgname='delivery_checklist_template_immutable') then raise exception 'published checklist immutability missing'; end if;
  if not exists(select 1 from pg_trigger where tgname='delivery_checklist_response_immutable') then raise exception 'response immutability missing'; end if;
  if not exists(select 1 from pg_trigger where tgname='delivery_checklist_completion_guard') then raise exception 'completion guard missing'; end if;
  if not exists(select 1 from information_schema.columns where table_schema='sales' and table_name='delivery_handover' and column_name='checklist_version') then raise exception 'handover checklist binding missing'; end if;
end $$;
rollback;
