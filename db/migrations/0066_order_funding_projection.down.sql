begin;
do $$begin if exists(select 1 from payment.local_funding_receipt) then raise exception 'preserve local funding evidence';end if;end $$;
-- No CASCADE: dependent handover columns must be rolled back explicitly first.
drop view payment.local_funding_evidence;
drop view payment.order_funding;
drop table payment.local_funding_receipt;
drop function payment.local_funding_immutable();
commit;
