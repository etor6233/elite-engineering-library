begin;

drop trigger if exists item_unit_of_measure_immutable on inventory.item_unit_of_measure;
drop function if exists inventory.reject_item_unit_of_measure_mutation();
drop table if exists inventory.item_unit_of_measure;

commit;
