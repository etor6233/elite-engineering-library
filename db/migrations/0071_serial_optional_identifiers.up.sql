begin;
-- AUTHORED corrective migration: the public Go contract makes VIN/battery optional.
-- PostgreSQL18 section5.5.3: NULLS DISTINCT preserves uniqueness of known values.
-- No existing row, mandatory serial, key or historical migration is rewritten.
alter table factory.production_unit
 drop constraint production_unit_tenant_id_vin_key,
 drop constraint production_unit_tenant_id_battery_serial_number_key,
 add constraint production_unit_tenant_id_vin_key unique nulls distinct (tenant_id,vin),
 add constraint production_unit_tenant_id_battery_serial_number_key unique nulls distinct (tenant_id,battery_serial_number);
alter table inventory.stock_unit
 drop constraint stock_unit_tenant_id_vin_key,
 drop constraint stock_unit_tenant_id_battery_serial_number_key,
 add constraint stock_unit_tenant_id_vin_key unique nulls distinct (tenant_id,vin),
 add constraint stock_unit_tenant_id_battery_serial_number_key unique nulls distinct (tenant_id,battery_serial_number);
commit;
