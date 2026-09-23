begin;
do $$ begin
 if exists(select 1 from catalog.marketplace_observation) or exists(select 1 from catalog.marketplace_effect where operation in ('MEDIA','CREATE')) then
  raise exception 'cannot remove marketplace initial publication evidence';
 end if;
end $$;
drop table catalog.marketplace_observation;
alter table catalog.marketplace_effect drop column operation;
commit;
