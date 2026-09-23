begin;

create table inventory.bulk_uom_migration_backfill (
  tenant_id uuid not null,
  item_id text not null,
  uom_code text not null,
  primary key (tenant_id,item_id),
  foreign key (tenant_id,item_id) references inventory.stock_item(tenant_id,item_id)
);

insert into inventory.bulk_uom_migration_backfill(tenant_id,item_id,uom_code)
select i.tenant_id,i.item_id,i.base_uom
  from inventory.stock_item i
 where not exists (select 1 from inventory.item_unit_of_measure u where u.tenant_id=i.tenant_id and u.item_id=i.item_id and u.is_base);

insert into inventory.item_unit_of_measure(tenant_id,item_id,uom_code,qty_per_uom,rounding_precision,is_base,version)
select tenant_id,item_id,uom_code,1,0.000001,true,1
  from inventory.bulk_uom_migration_backfill;

create table inventory.bulk_uom_balance (
  balance_id bigint not null references inventory.bulk_balance(balance_id) on delete cascade,
  uom_code text not null,
  quantity numeric(20,6) not null check (quantity >= 0),
  qty_per_uom numeric(20,6) not null check (qty_per_uom > 0),
  quantity_base numeric(20,6) not null check (quantity_base >= 0),
  version bigint not null default 1 check (version > 0),
  updated_at timestamptz not null default clock_timestamp(),
  primary key (balance_id, uom_code),
  check (quantity * qty_per_uom = quantity_base)
);

create table inventory.bulk_uom_conversion (
  tenant_id uuid not null,
  conversion_id text not null,
  request_id text not null,
  organization_id text not null,
  bin_id text not null,
  item_id text not null,
  lot_id text,
  operation text not null check (operation in ('breakbulk','gather')),
  from_uom text not null,
  to_uom text not null,
  from_quantity numeric(20,6) not null check (from_quantity > 0),
  to_quantity numeric(20,6) not null check (to_quantity > 0),
  base_quantity numeric(20,6) not null check (base_quantity > 0),
  from_qty_per_uom numeric(20,6) not null check (from_qty_per_uom > 0),
  to_qty_per_uom numeric(20,6) not null check (to_qty_per_uom > 0),
  version bigint not null default 1 check (version > 0),
  created_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, conversion_id),
  unique (tenant_id, request_id),
  foreign key (tenant_id, item_id, from_uom) references inventory.item_unit_of_measure(tenant_id,item_id,uom_code),
  foreign key (tenant_id, item_id, to_uom) references inventory.item_unit_of_measure(tenant_id,item_id,uom_code),
  check (from_uom <> to_uom),
  check (from_quantity * from_qty_per_uom = base_quantity),
  check (to_quantity * to_qty_per_uom = base_quantity),
  check ((operation = 'breakbulk' and from_qty_per_uom > to_qty_per_uom) or (operation = 'gather' and from_qty_per_uom < to_qty_per_uom))
);

create table inventory.bulk_uom_conversion_line (
  tenant_id uuid not null,
  conversion_id text not null,
  line_id text not null,
  sequence_no smallint not null check (sequence_no in (1,2)),
  action_type text not null check (action_type in ('take','place')),
  uom_code text not null,
  quantity numeric(20,6) not null check (quantity > 0),
  quantity_base numeric(20,6) not null check (quantity_base > 0),
  created_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, conversion_id, line_id),
  unique (tenant_id, conversion_id, sequence_no),
  foreign key (tenant_id, conversion_id) references inventory.bulk_uom_conversion(tenant_id,conversion_id),
  check ((sequence_no=1 and action_type='take') or (sequence_no=2 and action_type='place'))
);

create function inventory.validate_bulk_uom_balance()
returns trigger language plpgsql as $$
declare
  expected_factor numeric(20,6);
  expected_precision numeric(20,6);
begin
  select u.qty_per_uom,u.rounding_precision into expected_factor,expected_precision
    from inventory.bulk_balance b
    join inventory.item_unit_of_measure u on u.tenant_id=b.tenant_id and u.item_id=b.item_id and u.uom_code=new.uom_code
   where b.balance_id=new.balance_id;
  if expected_factor is null or new.qty_per_uom<>expected_factor or mod(new.quantity,expected_precision)<>0 then
    raise exception using errcode='55000',message='bulk UOM composition does not match the immutable item UOM contract';
  end if;
  return new;
end;
$$;

create trigger bulk_uom_balance_validate
before insert or update on inventory.bulk_uom_balance
for each row execute function inventory.validate_bulk_uom_balance();

create function inventory.sync_bulk_base_uom()
returns trigger language plpgsql as $$
declare
  delta numeric(20,6);
  base_code text;
  base_factor numeric(20,6);
  base_precision numeric(20,6);
begin
  delta := new.quantity - case when tg_op='INSERT' then 0 else old.quantity end;
  if delta=0 then return new; end if;
  select u.uom_code,u.qty_per_uom,u.rounding_precision into base_code,base_factor,base_precision
    from inventory.item_unit_of_measure u
   where u.tenant_id=new.tenant_id and u.item_id=new.item_id and u.is_base;
  if base_code is null then
    raise exception using errcode='55000',message='base UOM missing for bulk balance';
  end if;
  if delta>0 then
    insert into inventory.bulk_uom_balance(balance_id,uom_code,quantity,qty_per_uom,quantity_base)
    values(new.balance_id,base_code,delta/base_factor,base_factor,delta)
    on conflict(balance_id,uom_code) do update
      set quantity=inventory.bulk_uom_balance.quantity+excluded.quantity,
          quantity_base=inventory.bulk_uom_balance.quantity_base+excluded.quantity_base,
          version=inventory.bulk_uom_balance.version+1,updated_at=clock_timestamp();
  else
    update inventory.bulk_uom_balance
       set quantity=quantity-((-delta)/qty_per_uom),quantity_base=quantity_base+delta,
           version=version+1,updated_at=clock_timestamp()
     where balance_id=new.balance_id and uom_code=base_code and quantity_base>=-delta
       and mod((-delta)/qty_per_uom,base_precision)=0;
    if not found then
      raise exception using errcode='55000',message='bulk operation requires explicit handling-unit conversion';
    end if;
    delete from inventory.bulk_uom_balance where balance_id=new.balance_id and uom_code=base_code and quantity_base=0;
  end if;
  return new;
end;
$$;

create trigger bulk_balance_sync_base_uom
after insert or update of quantity on inventory.bulk_balance
for each row execute function inventory.sync_bulk_base_uom();

insert into inventory.bulk_uom_balance(balance_id,uom_code,quantity,qty_per_uom,quantity_base)
select b.balance_id,u.uom_code,b.quantity/u.qty_per_uom,u.qty_per_uom,b.quantity
  from inventory.bulk_balance b
  join inventory.item_unit_of_measure u on u.tenant_id=b.tenant_id and u.item_id=b.item_id and u.is_base
 where b.quantity>0;

create function inventory.assert_bulk_uom_conservation()
returns trigger language plpgsql as $$
declare
  affected bigint;
  physical numeric(20,6);
  composed numeric(20,6);
begin
  affected := case when tg_table_name='bulk_balance' then coalesce(new.balance_id,old.balance_id) else coalesce(new.balance_id,old.balance_id) end;
  select quantity into physical from inventory.bulk_balance where balance_id=affected;
  if physical is null then return null; end if;
  select coalesce(sum(quantity_base),0) into composed from inventory.bulk_uom_balance where balance_id=affected;
  if physical<>composed then
    raise exception using errcode='55000',message='bulk physical quantity and UOM composition diverged';
  end if;
  return null;
end;
$$;

create constraint trigger bulk_balance_uom_conservation
after insert or update of quantity or delete on inventory.bulk_balance
deferrable initially deferred for each row execute function inventory.assert_bulk_uom_conservation();

create constraint trigger bulk_uom_balance_conservation
after insert or update or delete on inventory.bulk_uom_balance
deferrable initially deferred for each row execute function inventory.assert_bulk_uom_conservation();

create function inventory.reject_bulk_uom_conversion_mutation()
returns trigger language plpgsql as $$
begin
  raise exception using errcode='55000',message='bulk UOM conversion evidence is immutable';
end;
$$;

create trigger bulk_uom_conversion_immutable
before update or delete on inventory.bulk_uom_conversion
for each row execute function inventory.reject_bulk_uom_conversion_mutation();

create trigger bulk_uom_conversion_line_immutable
before update or delete on inventory.bulk_uom_conversion_line
for each row execute function inventory.reject_bulk_uom_conversion_mutation();

commit;
