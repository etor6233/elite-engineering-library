begin;

alter table inventory.warehouse_bin
  add column cross_dock boolean not null default false,
  add constraint warehouse_bin_cross_dock_type_check
    check (not cross_dock or bin_type in ('pick','putpick'));

create table inventory.warehouse_crossdock_policy (
  tenant_id uuid not null,
  organization_id text not null,
  item_id text not null,
  bin_id text not null,
  due_date_days integer not null check (due_date_days between 0 and 365),
  enabled boolean not null,
  version bigint not null check (version > 0),
  created_at timestamptz not null default clock_timestamp(),
  updated_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, organization_id, item_id),
  foreign key (tenant_id, organization_id, bin_id)
    references inventory.warehouse_bin (tenant_id, organization_id, bin_id),
  foreign key (tenant_id, item_id)
    references inventory.stock_item (tenant_id, item_id)
);

create table inventory.warehouse_crossdock_allocation (
  tenant_id uuid not null,
  allocation_id text not null,
  receipt_id text not null,
  receipt_line_id text not null,
  put_away_activity_id text not null,
  put_away_line_id text not null,
  organization_id text not null,
  bin_id text not null,
  item_id text not null,
  lot_id text,
  transfer_id text not null,
  transfer_line_id text not null,
  quantity numeric(20,6) not null check (quantity > 0),
  reserved_quantity numeric(20,6) not null default 0,
  picked_quantity numeric(20,6) not null default 0,
  due_date date not null,
  available boolean not null default false,
  created_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, allocation_id),
  foreign key (tenant_id, receipt_id, receipt_line_id)
    references inventory.warehouse_receipt_line (tenant_id, receipt_id, line_id),
  foreign key (tenant_id, put_away_activity_id, put_away_line_id)
    references inventory.warehouse_activity_line (tenant_id, activity_id, line_id),
  foreign key (tenant_id, organization_id, bin_id)
    references inventory.warehouse_bin (tenant_id, organization_id, bin_id),
  foreign key (tenant_id, item_id)
    references inventory.stock_item (tenant_id, item_id),
  foreign key (tenant_id, lot_id)
    references inventory.inventory_lot (tenant_id, lot_id),
  foreign key (tenant_id, transfer_id, transfer_line_id)
    references inventory.bulk_transfer_line (tenant_id, transfer_id, line_id),
  unique (tenant_id, put_away_activity_id, put_away_line_id),
  check (reserved_quantity >= 0 and picked_quantity >= 0 and reserved_quantity + picked_quantity <= quantity)
);

create index warehouse_crossdock_allocation_demand_idx
  on inventory.warehouse_crossdock_allocation(tenant_id,organization_id,transfer_id,transfer_line_id,item_id,available,due_date);

create table inventory.warehouse_crossdock_pick_link (
  tenant_id uuid not null,
  pick_activity_id text not null,
  pick_line_id text not null,
  allocation_id text not null,
  quantity numeric(20,6) not null check (quantity > 0),
  status text not null check (status in ('open','picked','released')),
  created_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, pick_activity_id, pick_line_id),
  foreign key (tenant_id, pick_activity_id, pick_line_id)
    references inventory.warehouse_activity_line (tenant_id, activity_id, line_id),
  foreign key (tenant_id, allocation_id)
    references inventory.warehouse_crossdock_allocation (tenant_id, allocation_id)
);

commit;
