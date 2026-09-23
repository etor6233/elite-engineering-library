begin;
do $$ begin
  if to_regclass('fiscal.return_credit_note_link') is null then raise exception 'return credit note link missing'; end if;
  if (select count(*) from information_schema.table_constraints where table_schema='fiscal' and table_name='return_credit_note_link' and constraint_type='FOREIGN KEY')<>4 then raise exception 'return credit note foreign keys missing'; end if;
  if not exists(select 1 from pg_trigger where tgrelid='fiscal.return_credit_note_link'::regclass and tgname='return_credit_note_link_immutable' and not tgisinternal) then raise exception 'return credit note immutability missing'; end if;
end $$;
rollback;
