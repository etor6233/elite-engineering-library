begin;

do $block$
begin
  if exists (
    select 1 from inventory.warehouse_activity
     where status='open' and activity_type in ('pick','movement')
  ) then
    raise exception using errcode='55000',
      message='drain open pick and movement activities before packaging-aware migration';
  end if;
end;
$block$;

alter table inventory.warehouse_activity_line
  add column from_uom_code text,
  add column from_uom_quantity numeric(20,6),
  add column from_qty_per_uom numeric(20,6);

alter table inventory.warehouse_activity_line disable trigger warehouse_activity_line_immutable;
update inventory.warehouse_activity_line
   set from_uom_code=uom_code,
       from_uom_quantity=uom_quantity,
       from_qty_per_uom=qty_per_uom;
alter table inventory.warehouse_activity_line enable trigger warehouse_activity_line_immutable;

alter table inventory.warehouse_activity_line
  alter column from_uom_code set not null,
  alter column from_uom_quantity set not null,
  alter column from_qty_per_uom set not null,
  add constraint warehouse_activity_line_from_uom_quantity_check check (from_uom_quantity > 0),
  add constraint warehouse_activity_line_from_qty_per_uom_check check (from_qty_per_uom > 0),
  add constraint warehouse_activity_line_from_uom_fk foreign key (tenant_id,item_id,from_uom_code)
    references inventory.item_unit_of_measure(tenant_id,item_id,uom_code);

create function inventory.validate_warehouse_line_from_uom()
returns trigger language plpgsql as $function$
declare
  expected_factor numeric(20,6);
  expected_precision numeric(20,6);
begin
  select u.qty_per_uom,u.rounding_precision
    into expected_factor,expected_precision
    from inventory.item_unit_of_measure u
   where u.tenant_id=new.tenant_id and u.item_id=new.item_id and u.uom_code=new.from_uom_code;
  if not found or new.from_qty_per_uom<>expected_factor or
     mod(new.from_uom_quantity,expected_precision)<>0 or
     new.from_uom_quantity*new.from_qty_per_uom<new.quantity then
    raise exception using errcode='55000',message='warehouse line source UOM does not match the immutable item UOM contract';
  end if;
  return new;
end;
$function$;

create trigger warehouse_activity_line_from_uom_validate
before insert or update on inventory.warehouse_activity_line
for each row execute function inventory.validate_warehouse_line_from_uom();

alter table inventory.warehouse_replenishment_request
  add column target_uom text,
  add column allow_breakbulk boolean not null default false;

update inventory.warehouse_replenishment_request r
   set target_uom=u.uom_code
  from inventory.item_unit_of_measure u
 where u.tenant_id=r.tenant_id and u.item_id=r.item_id and u.is_base;

alter table inventory.warehouse_replenishment_request
  alter column target_uom set not null,
  add constraint warehouse_replenishment_target_uom_fk foreign key (tenant_id,item_id,target_uom)
    references inventory.item_unit_of_measure(tenant_id,item_id,uom_code);

create table inventory.warehouse_activity_uom_reservation (
  tenant_id uuid not null,
  activity_id text not null,
  line_id text not null,
  organization_id text not null,
  source_balance_id bigint not null,
  source_uom_code text not null,
  source_uom_quantity numeric(20,6) not null check (source_uom_quantity > 0),
  source_qty_per_uom numeric(20,6) not null check (source_qty_per_uom > 0),
  source_base_quantity numeric(20,6) not null check (source_base_quantity > 0),
  target_uom_code text not null,
  target_uom_quantity numeric(20,6) not null check (target_uom_quantity > 0),
  target_qty_per_uom numeric(20,6) not null check (target_qty_per_uom > 0),
  moved_base_quantity numeric(20,6) not null check (moved_base_quantity > 0),
  conversion_required boolean not null,
  status text not null default 'open' check (status in ('open','consumed','released')),
  version bigint not null default 1 check (version > 0),
  created_at timestamptz not null default clock_timestamp(),
  updated_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id,activity_id,line_id),
  foreign key (tenant_id,activity_id,line_id)
    references inventory.warehouse_activity_line(tenant_id,activity_id,line_id),
  foreign key (tenant_id,organization_id)
    references org.organization(tenant_id,organization_id),
  foreign key (tenant_id,activity_id,organization_id)
    references inventory.warehouse_activity(tenant_id,activity_id,organization_id),
  check (source_uom_quantity*source_qty_per_uom=source_base_quantity),
  check (target_uom_quantity*target_qty_per_uom=moved_base_quantity),
  check (source_base_quantity>=moved_base_quantity),
  check (conversion_required=(source_uom_code<>target_uom_code))
);

create function inventory.enforce_warehouse_uom_reservation_transition()
returns trigger language plpgsql as $function$
begin
  if old.status<>'open' or new.status not in ('consumed','released') or
     new.version<>old.version+1 or
     new.tenant_id<>old.tenant_id or new.activity_id<>old.activity_id or
     new.line_id<>old.line_id or new.organization_id<>old.organization_id or
     new.source_balance_id<>old.source_balance_id or
     new.source_uom_code<>old.source_uom_code or
     new.source_uom_quantity<>old.source_uom_quantity or
     new.source_qty_per_uom<>old.source_qty_per_uom or
     new.source_base_quantity<>old.source_base_quantity or
     new.target_uom_code<>old.target_uom_code or
     new.target_uom_quantity<>old.target_uom_quantity or
     new.target_qty_per_uom<>old.target_qty_per_uom or
     new.moved_base_quantity<>old.moved_base_quantity or
     new.conversion_required<>old.conversion_required or
     new.created_at<>old.created_at then
    raise exception using errcode='55000',message='invalid warehouse UOM reservation transition';
  end if;
  new.updated_at:=clock_timestamp();
  return new;
end;
$function$;

create trigger warehouse_uom_reservation_transition_guard
before update on inventory.warehouse_activity_uom_reservation
for each row execute function inventory.enforce_warehouse_uom_reservation_transition();

create function inventory.reject_warehouse_uom_reservation_delete()
returns trigger language plpgsql as $function$
begin
  raise exception using errcode='55000',message='warehouse UOM reservation evidence is immutable';
end;
$function$;

create trigger warehouse_uom_reservation_delete_guard
before delete on inventory.warehouse_activity_uom_reservation
for each row execute function inventory.reject_warehouse_uom_reservation_delete();

commit;
