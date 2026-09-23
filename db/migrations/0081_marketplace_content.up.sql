-- AUTHORED extension of the same immutable provider effect owner.
begin;
alter table catalog.marketplace_effect drop constraint marketplace_effect_operation_check;
alter table catalog.marketplace_effect add constraint marketplace_effect_operation_check
 check(operation in ('LEGACY_MUTATION','PRICE','STOCK','PAUSE','RESUME','MEDIA','CREATE','CONTENT'));
commit;
