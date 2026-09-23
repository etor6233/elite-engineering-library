begin;

create table inventory.item_unit_of_measure (
  tenant_id uuid not null,
  item_id text not null,
  uom_code text not null,
  qty_per_uom numeric(20,6) not null check (qty_per_uom > 0),
  rounding_precision numeric(20,6) not null check (rounding_precision > 0 and rounding_precision <= 1),
  is_base boolean not null default false,
  version bigint not null default 1 check (version > 0),
  created_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, item_id, uom_code),
  foreign key (tenant_id, item_id) references inventory.stock_item (tenant_id, item_id),
  check (uom_code ~ '^[A-Z0-9][A-Z0-9._/-]{0,15}$'),
  check (not is_base or qty_per_uom = 1),
  check (mod(qty_per_uom, rounding_precision) = 0)
);

create unique index item_unit_of_measure_one_base_uidx
  on inventory.item_unit_of_measure(tenant_id, item_id)
  where is_base;

create function inventory.reject_item_unit_of_measure_mutation()
returns trigger language plpgsql as $$
begin
  raise exception using errcode='55000', message='item unit of measure is immutable; create a new item or approved code instead';
end;
$$;

create trigger item_unit_of_measure_immutable
before update or delete on inventory.item_unit_of_measure
for each row execute function inventory.reject_item_unit_of_measure_mutation();

commit;
