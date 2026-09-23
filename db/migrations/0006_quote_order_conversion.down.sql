begin;

drop table if exists sales.quotation_acceptance;
drop index if exists sales.quotation_order_unique_idx;
alter table sales.quotation
  drop constraint if exists quotation_acceptance_order_check,
  drop constraint if exists quotation_order_fk,
  drop column if exists order_id;

commit;
