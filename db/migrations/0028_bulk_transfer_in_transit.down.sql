begin;
drop view if exists inventory.bulk_inventory_in_transit;
drop table if exists inventory.bulk_transfer_cost_component;
drop table if exists inventory.bulk_transfer_allocation;
drop table if exists inventory.bulk_transfer_line;
drop table if exists inventory.bulk_transfer;
commit;
