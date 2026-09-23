begin;

create table inventory.bulk_transfer (
  tenant_id uuid not null,
  transfer_id text not null,
  request_id text not null,
  from_organization_id text not null,
  to_organization_id text not null,
  receive_bin_id text not null,
  in_transit_code text not null,
  status text not null check (status in ('released','shipped','received','cancelled')),
  posting_date date not null,
  source_kind text not null,
  source_id text not null,
  put_away_activity_id text,
  version bigint not null check (version > 0),
  shipped_at timestamptz,
  received_at timestamptz,
  created_at timestamptz not null default clock_timestamp(),
  updated_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, transfer_id),
  foreign key (tenant_id, from_organization_id) references org.organization (tenant_id, organization_id),
  foreign key (tenant_id, to_organization_id, receive_bin_id) references inventory.warehouse_bin (tenant_id, organization_id, bin_id),
  foreign key (tenant_id, put_away_activity_id) references inventory.warehouse_activity (tenant_id, activity_id),
  unique (tenant_id, request_id),
  check (from_organization_id <> to_organization_id),
  check (length(in_transit_code) between 1 and 50),
  check ((status='released' and shipped_at is null and received_at is null and put_away_activity_id is null) or
         (status='shipped' and shipped_at is not null and received_at is null and put_away_activity_id is null) or
         (status='received' and shipped_at is not null and received_at is not null and put_away_activity_id is not null) or
         (status='cancelled' and received_at is null and put_away_activity_id is null))
);

create table inventory.bulk_transfer_line (
  tenant_id uuid not null,
  transfer_id text not null,
  line_id text not null,
  item_id text not null,
  quantity numeric(20,6) not null check (quantity > 0),
  shipped_quantity numeric(20,6) not null default 0 check (shipped_quantity >= 0 and shipped_quantity <= quantity),
  received_quantity numeric(20,6) not null default 0 check (received_quantity >= 0 and received_quantity <= shipped_quantity),
  cost_amount numeric(24,4) not null default 0 check (cost_amount >= 0),
  primary key (tenant_id, transfer_id, line_id),
  foreign key (tenant_id, transfer_id) references inventory.bulk_transfer (tenant_id, transfer_id),
  foreign key (tenant_id, item_id) references inventory.stock_item (tenant_id, item_id),
  unique (tenant_id, transfer_id, item_id)
);

create table inventory.bulk_transfer_allocation (
  tenant_id uuid not null,
  transfer_id text not null,
  line_id text not null,
  allocation_id text not null,
  reservation_id text not null,
  from_bin_id text not null,
  item_id text not null,
  lot_id text,
  quantity numeric(20,6) not null check (quantity > 0),
  outbound_entry_id text not null,
  primary key (tenant_id, transfer_id, allocation_id),
  foreign key (tenant_id, transfer_id, line_id) references inventory.bulk_transfer_line (tenant_id, transfer_id, line_id),
  foreign key (tenant_id, reservation_id) references inventory.bulk_reservation (tenant_id, reservation_id),
  foreign key (tenant_id, outbound_entry_id) references inventory.bulk_inventory_entry (tenant_id, entry_id),
  foreign key (tenant_id, lot_id) references inventory.inventory_lot (tenant_id, lot_id),
  unique (tenant_id, reservation_id)
);

create table inventory.bulk_transfer_cost_component (
  tenant_id uuid not null,
  transfer_id text not null,
  allocation_id text not null,
  component_id text not null,
  source_receipt_entry_id text not null,
  quantity numeric(20,6) not null check (quantity > 0),
  unit_cost numeric(20,4) not null check (unit_cost >= 0),
  cost_amount numeric(24,4) not null check (cost_amount >= 0),
  receipt_entry_id text,
  primary key (tenant_id, transfer_id, component_id),
  foreign key (tenant_id, transfer_id, allocation_id) references inventory.bulk_transfer_allocation (tenant_id, transfer_id, allocation_id),
  foreign key (tenant_id, source_receipt_entry_id) references inventory.bulk_inventory_entry (tenant_id, entry_id),
  foreign key (tenant_id, receipt_entry_id) references inventory.bulk_inventory_entry (tenant_id, entry_id),
  check (round(quantity * unit_cost,4)=cost_amount)
);

create index bulk_transfer_in_transit_idx
  on inventory.bulk_transfer(tenant_id,to_organization_id,status,posting_date)
  where status='shipped';

create view inventory.bulk_inventory_in_transit as
select t.tenant_id,t.transfer_id,t.from_organization_id,t.to_organization_id,t.in_transit_code,
       l.line_id,l.item_id,a.lot_id,sum(c.quantity) quantity,sum(c.cost_amount) cost_amount,
       t.shipped_at,t.version
from inventory.bulk_transfer t
join inventory.bulk_transfer_line l using(tenant_id,transfer_id)
join inventory.bulk_transfer_allocation a using(tenant_id,transfer_id,line_id)
join inventory.bulk_transfer_cost_component c using(tenant_id,transfer_id,allocation_id)
where t.status='shipped'
group by t.tenant_id,t.transfer_id,t.from_organization_id,t.to_organization_id,t.in_transit_code,
         l.line_id,l.item_id,a.lot_id,t.shipped_at,t.version;

comment on view inventory.bulk_inventory_in_transit is
  'BC-derived transfer order evidence: shipped demand remains unavailable at source and visible in transit until destination receipt.';

commit;
