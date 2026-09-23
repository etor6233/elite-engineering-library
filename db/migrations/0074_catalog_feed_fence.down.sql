begin;
do $$begin
 if exists(select 1 from catalog.release_feed_intent) then raise exception 'catalog feed history exists; downgrade refused';end if;
end$$;
drop table catalog.release_feed_intent;
commit;
