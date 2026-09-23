begin;
-- Restoring the historical restriction is refused when valid new data needs NULLS DISTINCT.
-- The entire downgrade is atomic; never delete, rewrite or synthesize identifiers.
do $$
begin
 if exists(select 1 from factory.production_unit where vin is null group by tenant_id having count(*)>1)
 or exists(select 1 from factory.production_unit where battery_serial_number is null group by tenant_id having count(*)>1)
 or exists(select 1 from inventory.stock_unit where vin is null group by tenant_id having count(*)>1)
 or exists(select 1 from inventory.stock_unit where battery_serial_number is null group by tenant_id having count(*)>1)
 then raise exception 'optional serial identifiers require NULLS DISTINCT; downgrade refused';end if;
end $$;
alter table factory.production_unit
 drop constraint production_unit_tenant_id_vin_key,
 drop constraint production_unit_tenant_id_battery_serial_number_key,
 add constraint production_unit_tenant_id_vin_key unique nulls not distinct (tenant_id,vin),
 add constraint production_unit_tenant_id_battery_serial_number_key unique nulls not distinct (tenant_id,battery_serial_number);
alter table inventory.stock_unit
 drop constraint stock_unit_tenant_id_vin_key,
 drop constraint stock_unit_tenant_id_battery_serial_number_key,
 add constraint stock_unit_tenant_id_vin_key unique nulls not distinct (tenant_id,vin),
 add constraint stock_unit_tenant_id_battery_serial_number_key unique nulls not distinct (tenant_id,battery_serial_number);
commit;
