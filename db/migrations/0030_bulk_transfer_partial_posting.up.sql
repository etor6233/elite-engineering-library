begin;

do $migration$
declare
  constraint_name text;
begin
  for constraint_name in
    select conname
    from pg_constraint
    where conrelid='inventory.bulk_transfer'::regclass
      and contype='c'
      and pg_get_constraintdef(oid) ilike '%status%'
  loop
    execute format('alter table inventory.bulk_transfer drop constraint %I', constraint_name);
  end loop;
end;
$migration$;

alter table inventory.bulk_transfer
  add constraint bulk_transfer_status_v164_check
  check (status in ('released','partially-shipped','shipped','partially-received','received','cancelled')),
  add constraint bulk_transfer_state_v164_check
  check (
    (status='released' and shipped_at is null and received_at is null and put_away_activity_id is null) or
    (status='partially-shipped' and shipped_at is not null and received_at is null and put_away_activity_id is null) or
    (status='shipped' and shipped_at is not null and received_at is null and put_away_activity_id is null) or
    (status='partially-received' and shipped_at is not null and received_at is null and put_away_activity_id is not null) or
    (status='received' and shipped_at is not null and received_at is not null and put_away_activity_id is not null) or
    (status='cancelled' and received_at is null and put_away_activity_id is null)
  );

create table inventory.bulk_transfer_shipment (
  tenant_id uuid not null,
  shipment_id text not null,
  transfer_id text not null,
  request_id text not null,
  warehouse_activity_id text not null,
  quantity numeric(20,6) not null check (quantity > 0),
  cost_amount numeric(24,4) not null check (cost_amount >= 0),
  posting_date date not null,
  created_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, shipment_id),
  foreign key (tenant_id, transfer_id) references inventory.bulk_transfer (tenant_id, transfer_id),
  foreign key (tenant_id, warehouse_activity_id) references inventory.warehouse_activity (tenant_id, activity_id),
  unique (tenant_id, transfer_id, request_id),
  unique (tenant_id, warehouse_activity_id)
);

create table inventory.bulk_transfer_receipt (
  tenant_id uuid not null,
  receipt_id text not null,
  transfer_id text not null,
  request_id text not null,
  quantity numeric(20,6) not null check (quantity > 0),
  cost_amount numeric(24,4) not null check (cost_amount >= 0),
  posting_date date not null,
  put_away_activity_id text not null,
  created_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, receipt_id),
  foreign key (tenant_id, transfer_id) references inventory.bulk_transfer (tenant_id, transfer_id),
  foreign key (tenant_id, put_away_activity_id) references inventory.warehouse_activity (tenant_id, activity_id),
  unique (tenant_id, transfer_id, request_id),
  unique (tenant_id, put_away_activity_id)
);

alter table inventory.bulk_transfer_allocation
  add column shipment_id text,
  add constraint bulk_transfer_allocation_shipment_fk
    foreign key (tenant_id, shipment_id)
    references inventory.bulk_transfer_shipment (tenant_id, shipment_id);

alter table inventory.bulk_transfer_cost_component
  add column received_quantity numeric(20,6) not null default 0,
  add constraint bulk_transfer_cost_component_received_check
    check (received_quantity >= 0 and received_quantity <= quantity);

update inventory.bulk_transfer_cost_component
set received_quantity=quantity
where receipt_entry_id is not null;

create table inventory.bulk_transfer_receipt_component (
  tenant_id uuid not null,
  receipt_id text not null,
  transfer_id text not null,
  component_id text not null,
  receipt_entry_id text not null,
  quantity numeric(20,6) not null check (quantity > 0),
  cost_amount numeric(24,4) not null check (cost_amount >= 0),
  primary key (tenant_id, receipt_id, component_id),
  foreign key (tenant_id, receipt_id) references inventory.bulk_transfer_receipt (tenant_id, receipt_id),
  foreign key (tenant_id, transfer_id, component_id) references inventory.bulk_transfer_cost_component (tenant_id, transfer_id, component_id),
  foreign key (tenant_id, receipt_entry_id) references inventory.bulk_inventory_entry (tenant_id, entry_id)
);

drop view inventory.bulk_inventory_in_transit;

create view inventory.bulk_inventory_in_transit as
select t.tenant_id,t.transfer_id,t.from_organization_id,t.to_organization_id,t.in_transit_code,
       l.line_id,l.item_id,a.lot_id,
       sum(c.quantity-c.received_quantity) quantity,
       sum(round((c.quantity-c.received_quantity)*c.unit_cost,4)) cost_amount,
       t.shipped_at,t.version
from inventory.bulk_transfer t
join inventory.bulk_transfer_line l using(tenant_id,transfer_id)
join inventory.bulk_transfer_allocation a using(tenant_id,transfer_id,line_id)
join inventory.bulk_transfer_cost_component c using(tenant_id,transfer_id,allocation_id)
where c.received_quantity<c.quantity
group by t.tenant_id,t.transfer_id,t.from_organization_id,t.to_organization_id,t.in_transit_code,
         l.line_id,l.item_id,a.lot_id,t.shipped_at,t.version;

comment on view inventory.bulk_inventory_in_transit is
  'BC-derived transfer evidence: only shipped cost components not yet received remain visible in transit.';

drop index inventory.bulk_transfer_in_transit_idx;
create index bulk_transfer_in_transit_idx
  on inventory.bulk_transfer(tenant_id,to_organization_id,status,posting_date)
  where status in ('partially-shipped','shipped','partially-received');

commit;
