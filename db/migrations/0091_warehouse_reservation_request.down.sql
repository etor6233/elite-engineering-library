begin;
do $$begin if exists(select 1 from inventory.workspace_reservation_request) then raise exception 'reservation recovery records exist; down migration refused';end if;end$$;
drop table inventory.workspace_reservation_request;
commit;
