begin;

create table inventory.stock_item (
  tenant_id uuid not null,
  item_id text not null,
  item_code text not null,
  description text not null,
  base_uom text not null,
  tracking_mode text not null check (tracking_mode in ('none', 'lot')),
  costing_method text not null check (costing_method in ('fifo', 'specific')),
  version bigint not null default 1 check (version > 0),
  created_at timestamptz not null default clock_timestamp(),
  updated_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, item_id),
  foreign key (tenant_id) references platform.tenant (tenant_id),
  unique (tenant_id, item_code),
  check (length(item_id) between 1 and 128),
  check (item_code ~ '^[A-Za-z0-9][A-Za-z0-9._/-]{0,63}$'),
  check (base_uom ~ '^[A-Z0-9][A-Z0-9._/-]{0,15}$')
);

create table inventory.warehouse_bin (
  tenant_id uuid not null,
  organization_id text not null,
  bin_id text not null,
  bin_code text not null,
  bin_type text not null check (bin_type in ('receive', 'ship', 'put-away', 'pick', 'putpick')),
  movement_blocked boolean not null default false,
  version bigint not null default 1 check (version > 0),
  created_at timestamptz not null default clock_timestamp(),
  updated_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, organization_id, bin_id),
  foreign key (tenant_id, organization_id) references org.organization (tenant_id, organization_id),
  unique (tenant_id, organization_id, bin_code),
  check (length(bin_id) between 1 and 128),
  check (bin_code ~ '^[A-Za-z0-9][A-Za-z0-9._/-]{0,63}$')
);

create table inventory.item_bin_policy (
  tenant_id uuid not null,
  organization_id text not null,
  item_id text not null,
  bin_id text not null,
  fixed boolean not null default false,
  dedicated boolean not null default false,
  is_default boolean not null default false,
  min_quantity numeric(20,6) not null default 0 check (min_quantity >= 0),
  max_quantity numeric(20,6),
  version bigint not null default 1 check (version > 0),
  created_at timestamptz not null default clock_timestamp(),
  updated_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, organization_id, item_id, bin_id),
  foreign key (tenant_id, organization_id, bin_id) references inventory.warehouse_bin (tenant_id, organization_id, bin_id),
  foreign key (tenant_id, item_id) references inventory.stock_item (tenant_id, item_id),
  check (max_quantity is null or max_quantity > min_quantity)
);

create unique index item_bin_policy_one_default_uidx
  on inventory.item_bin_policy(tenant_id, organization_id, item_id)
  where is_default;

create table inventory.inventory_lot (
  tenant_id uuid not null,
  lot_id text not null,
  item_id text not null,
  lot_no text not null,
  expiration_date date,
  warranty_date date,
  blocked boolean not null default false,
  version bigint not null default 1 check (version > 0),
  created_at timestamptz not null default clock_timestamp(),
  updated_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, lot_id),
  foreign key (tenant_id, item_id) references inventory.stock_item (tenant_id, item_id),
  unique (tenant_id, item_id, lot_no),
  check (length(lot_id) between 1 and 128),
  check (length(lot_no) between 1 and 50)
);

create table inventory.bulk_balance (
  tenant_id uuid not null,
  organization_id text not null,
  bin_id text not null,
  item_id text not null,
  lot_id text,
  quantity numeric(20,6) not null default 0 check (quantity >= 0),
  reserved_quantity numeric(20,6) not null default 0 check (reserved_quantity >= 0 and reserved_quantity <= quantity),
  version bigint not null default 1 check (version > 0),
  updated_at timestamptz not null default clock_timestamp(),
  balance_id bigint generated always as identity primary key,
  foreign key (tenant_id, organization_id, item_id, bin_id) references inventory.item_bin_policy (tenant_id, organization_id, item_id, bin_id),
  foreign key (tenant_id, lot_id) references inventory.inventory_lot (tenant_id, lot_id),
  unique nulls not distinct (tenant_id, organization_id, bin_id, item_id, lot_id)
);

create table inventory.bulk_reservation (
  tenant_id uuid not null,
  reservation_id text not null,
  organization_id text not null,
  bin_id text not null,
  item_id text not null,
  lot_id text,
  demand_kind text not null check (demand_kind in ('customer-order', 'service', 'transfer-outbound', 'manual')),
  demand_id text not null,
  demand_line_id text not null,
  quantity numeric(20,6) not null check (quantity > 0),
  status text not null check (status in ('reservation', 'released', 'consumed', 'expired')),
  cancellation_disallowed boolean not null default false,
  expires_at timestamptz,
  version bigint not null check (version > 0),
  created_at timestamptz not null default clock_timestamp(),
  updated_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, reservation_id),
  foreign key (tenant_id, organization_id, item_id, bin_id) references inventory.item_bin_policy (tenant_id, organization_id, item_id, bin_id),
  foreign key (tenant_id, lot_id) references inventory.inventory_lot (tenant_id, lot_id),
  check (expires_at is null or expires_at > created_at)
);

create unique index bulk_reservation_active_demand_uidx
  on inventory.bulk_reservation(tenant_id, demand_kind, demand_id, demand_line_id)
  where status='reservation';

create index bulk_reservation_scope_idx
  on inventory.bulk_reservation(tenant_id, organization_id, item_id, status, expires_at);

create table inventory.bulk_inventory_entry (
  tenant_id uuid not null,
  entry_id text not null,
  organization_id text not null,
  bin_id text not null,
  item_id text not null,
  lot_id text,
  entry_type text not null check (entry_type in ('receipt', 'issue', 'movement-in', 'movement-out', 'adjustment')),
  quantity numeric(20,6) not null check (quantity <> 0),
  unit_cost numeric(20,4) not null check (unit_cost >= 0),
  cost_amount numeric(24,4) not null,
  posting_date date not null,
  source_kind text not null,
  source_id text not null,
  created_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, entry_id),
  foreign key (tenant_id, organization_id, item_id, bin_id) references inventory.item_bin_policy (tenant_id, organization_id, item_id, bin_id),
  foreign key (tenant_id, lot_id) references inventory.inventory_lot (tenant_id, lot_id),
  unique (tenant_id, source_kind, source_id, entry_type, organization_id, bin_id),
  check ((quantity > 0 and entry_type in ('receipt','movement-in','adjustment')) or (quantity < 0 and entry_type in ('issue','movement-out','adjustment')))
);

create table inventory.bulk_cost_layer (
  tenant_id uuid not null,
  layer_id text not null,
  receipt_entry_id text not null,
  organization_id text not null,
  item_id text not null,
  lot_id text,
  posting_date date not null,
  original_quantity numeric(20,6) not null check (original_quantity > 0),
  remaining_quantity numeric(20,6) not null check (remaining_quantity >= 0 and remaining_quantity <= original_quantity),
  unit_cost numeric(20,4) not null check (unit_cost >= 0),
  primary key (tenant_id, layer_id),
  foreign key (tenant_id, receipt_entry_id) references inventory.bulk_inventory_entry (tenant_id, entry_id),
  foreign key (tenant_id, item_id) references inventory.stock_item (tenant_id, item_id),
  foreign key (tenant_id, lot_id) references inventory.inventory_lot (tenant_id, lot_id),
  unique (tenant_id, receipt_entry_id)
);

create index bulk_cost_layer_fifo_idx
  on inventory.bulk_cost_layer(tenant_id, organization_id, item_id, posting_date, receipt_entry_id)
  where remaining_quantity > 0;

create table inventory.bulk_cost_application (
  tenant_id uuid not null,
  outbound_entry_id text not null,
  inbound_entry_id text not null,
  quantity numeric(20,6) not null check (quantity > 0),
  cost_amount numeric(24,4) not null check (cost_amount >= 0),
  created_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, outbound_entry_id, inbound_entry_id),
  foreign key (tenant_id, outbound_entry_id) references inventory.bulk_inventory_entry (tenant_id, entry_id),
  foreign key (tenant_id, inbound_entry_id) references inventory.bulk_inventory_entry (tenant_id, entry_id)
);

create function inventory.reject_bulk_ledger_mutation()
returns trigger language plpgsql as $function$
begin
  raise exception using errcode='55000', message='bulk inventory ledger is append-only';
end;
$function$;

create trigger bulk_inventory_entry_immutable
before update or delete on inventory.bulk_inventory_entry
for each row execute function inventory.reject_bulk_ledger_mutation();

create trigger bulk_cost_application_immutable
before update or delete on inventory.bulk_cost_application
for each row execute function inventory.reject_bulk_ledger_mutation();

create view inventory.bulk_available as
select b.tenant_id,b.organization_id,b.bin_id,b.item_id,b.lot_id,
       b.quantity,b.reserved_quantity,(b.quantity-b.reserved_quantity) available_quantity,
       l.lot_no,l.expiration_date,p.fixed,p.dedicated,p.is_default,w.bin_code,w.bin_type
from inventory.bulk_balance b
join inventory.item_bin_policy p using(tenant_id,organization_id,item_id,bin_id)
join inventory.warehouse_bin w using(tenant_id,organization_id,bin_id)
left join inventory.inventory_lot l using(tenant_id,lot_id)
where not w.movement_blocked
  and not p.dedicated
  and coalesce(l.blocked,false)=false
  and (l.expiration_date is null or l.expiration_date >= current_date);

comment on view inventory.bulk_available is
  'Portable BC-derived warehouse availability: physical bin quantity minus reserved quantity, excluding movement-blocked/dedicated bins and blocked/expired lots.';

commit;
