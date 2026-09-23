begin;
do $$begin
 if exists(select 1 from franchise.network_command_receipt)then raise exception 'cannot discard durable network command evidence';end if;
end$$;
drop table franchise.network_command_receipt;
drop function franchise.reject_network_receipt_mutation();
commit;
