begin;

drop trigger if exists bulk_uom_conversion_line_immutable on inventory.bulk_uom_conversion_line;
drop trigger if exists bulk_uom_conversion_immutable on inventory.bulk_uom_conversion;
drop function if exists inventory.reject_bulk_uom_conversion_mutation();
drop trigger if exists bulk_uom_balance_conservation on inventory.bulk_uom_balance;
drop trigger if exists bulk_balance_uom_conservation on inventory.bulk_balance;
drop function if exists inventory.assert_bulk_uom_conservation();
drop trigger if exists bulk_balance_sync_base_uom on inventory.bulk_balance;
drop function if exists inventory.sync_bulk_base_uom();
drop trigger if exists bulk_uom_balance_validate on inventory.bulk_uom_balance;
drop function if exists inventory.validate_bulk_uom_balance();
drop table if exists inventory.bulk_uom_conversion_line;
drop table if exists inventory.bulk_uom_conversion;
drop table if exists inventory.bulk_uom_balance;
alter table inventory.item_unit_of_measure disable trigger item_unit_of_measure_immutable;
delete from inventory.item_unit_of_measure u using inventory.bulk_uom_migration_backfill b
 where u.tenant_id=b.tenant_id and u.item_id=b.item_id and u.uom_code=b.uom_code and u.is_base;
alter table inventory.item_unit_of_measure enable trigger item_unit_of_measure_immutable;
drop table if exists inventory.bulk_uom_migration_backfill;

commit;
