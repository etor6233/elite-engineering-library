begin;
do $$ begin
 if exists(select 1 from catalog.merchant_effect) then raise exception 'cannot remove merchant feed evidence';end if;
end $$;
drop table catalog.merchant_observation;
drop table catalog.merchant_effect;
commit;
