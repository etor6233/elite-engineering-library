begin;

alter table inventory.bulk_uom_balance
  add column reserved_quantity numeric(20,6) not null default 0,
  add constraint bulk_uom_balance_reserved_quantity_check
    check (reserved_quantity >= 0 and reserved_quantity <= quantity);

alter table inventory.bulk_uom_conversion
  add column source_kind text not null default 'manual',
  add column source_id text not null default '',
  add constraint bulk_uom_conversion_source_kind_check
    check (length(source_kind) between 1 and 64),
  add constraint bulk_uom_conversion_source_id_check
    check (length(source_id) <= 128);

alter table inventory.warehouse_receipt
  add column allow_breakbulk boolean not null default false;

alter table inventory.warehouse_receipt_line
  add column uom_code text,
  add column uom_quantity numeric(20,6),
  add column qty_per_uom numeric(20,6);

alter table inventory.warehouse_activity_line
  add column uom_code text,
  add column uom_quantity numeric(20,6),
  add column qty_per_uom numeric(20,6);

alter table inventory.warehouse_receipt_line disable trigger warehouse_receipt_line_immutable;
update inventory.warehouse_receipt_line l
   set uom_code=u.uom_code,
       uom_quantity=l.quantity/u.qty_per_uom,
       qty_per_uom=u.qty_per_uom
  from inventory.item_unit_of_measure u
 where u.tenant_id=l.tenant_id and u.item_id=l.item_id and u.is_base;
alter table inventory.warehouse_receipt_line enable trigger warehouse_receipt_line_immutable;

alter table inventory.warehouse_activity_line disable trigger warehouse_activity_line_immutable;
update inventory.warehouse_activity_line l
   set uom_code=u.uom_code,
       uom_quantity=l.quantity/u.qty_per_uom,
       qty_per_uom=u.qty_per_uom
  from inventory.item_unit_of_measure u
 where u.tenant_id=l.tenant_id and u.item_id=l.item_id and u.is_base;
alter table inventory.warehouse_activity_line enable trigger warehouse_activity_line_immutable;

alter table inventory.warehouse_receipt_line
  alter column uom_code set not null,
  alter column uom_quantity set not null,
  alter column qty_per_uom set not null,
  add constraint warehouse_receipt_line_uom_quantity_check check (uom_quantity > 0),
  add constraint warehouse_receipt_line_qty_per_uom_check check (qty_per_uom > 0),
  add constraint warehouse_receipt_line_uom_exact_check check (uom_quantity * qty_per_uom = quantity),
  add constraint warehouse_receipt_line_uom_fk foreign key (tenant_id,item_id,uom_code)
    references inventory.item_unit_of_measure(tenant_id,item_id,uom_code);

alter table inventory.warehouse_activity_line
  alter column uom_code set not null,
  alter column uom_quantity set not null,
  alter column qty_per_uom set not null,
  add constraint warehouse_activity_line_uom_quantity_check check (uom_quantity > 0),
  add constraint warehouse_activity_line_qty_per_uom_check check (qty_per_uom > 0),
  add constraint warehouse_activity_line_uom_exact_check check (uom_quantity * qty_per_uom = quantity),
  add constraint warehouse_activity_line_uom_fk foreign key (tenant_id,item_id,uom_code)
    references inventory.item_unit_of_measure(tenant_id,item_id,uom_code);

create function inventory.validate_warehouse_line_uom()
returns trigger language plpgsql as $function$
declare
  expected_factor numeric(20,6);
  expected_precision numeric(20,6);
begin
  select u.qty_per_uom,u.rounding_precision
    into expected_factor,expected_precision
    from inventory.item_unit_of_measure u
   where u.tenant_id=new.tenant_id and u.item_id=new.item_id and u.uom_code=new.uom_code;
  if not found or new.qty_per_uom<>expected_factor or mod(new.uom_quantity,expected_precision)<>0 then
    raise exception using errcode='55000',message='warehouse line UOM does not match the immutable item UOM contract';
  end if;
  return new;
end;
$function$;

create trigger warehouse_receipt_line_uom_validate
before insert or update on inventory.warehouse_receipt_line
for each row execute function inventory.validate_warehouse_line_uom();

create trigger warehouse_activity_line_uom_validate
before insert or update on inventory.warehouse_activity_line
for each row execute function inventory.validate_warehouse_line_uom();

commit;
