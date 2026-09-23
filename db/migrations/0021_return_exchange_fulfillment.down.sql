begin;

drop index if exists sales.return_exchange_original_order_idx;
drop trigger if exists return_exchange_immutable on sales.return_exchange;
drop function if exists sales.prevent_return_exchange_mutation();
drop table if exists sales.return_exchange;

commit;
