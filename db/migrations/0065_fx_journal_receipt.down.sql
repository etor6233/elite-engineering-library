begin;
do $$ begin
 if exists(select 1 from accounting.fx_journal_receipt) then
  raise exception 'FX journal receipts exist; preserve accounting history';
 end if;
end $$;
drop table accounting.fx_journal_receipt;
commit;
