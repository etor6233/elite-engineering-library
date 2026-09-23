begin;

alter table sales.customer_order
  add column fulfillment_state text not null default 'unfulfilled'
    check (fulfillment_state in ('unfulfilled','partially-shipped','shipped','delivered'));

alter table sales.customer_order_line
  add column shipped_quantity numeric(20,6) not null default 0
    check (shipped_quantity >= 0 and shipped_quantity <= quantity);

create table sales.customer_shipment (
  tenant_id uuid not null,
  shipment_id text not null,
  request_id text not null,
  organization_id text not null,
  order_id text not null,
  warehouse_activity_id text not null,
  posting_date date not null,
  total_cost_amount numeric(20,4) not null check (total_cost_amount >= 0),
  order_version bigint not null check (order_version > 0),
  fulfillment_state text not null check (fulfillment_state in ('partially-shipped','shipped')),
  created_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, shipment_id),
  unique (tenant_id, request_id),
  unique (tenant_id, warehouse_activity_id),
  foreign key (tenant_id, organization_id) references org.organization(tenant_id,organization_id),
  foreign key (tenant_id, order_id) references sales.customer_order(tenant_id,order_id),
  foreign key (tenant_id, warehouse_activity_id, organization_id)
    references inventory.warehouse_activity(tenant_id,activity_id,organization_id),
  check (length(request_id) between 1 and 128)
);

create table sales.customer_shipment_line (
  tenant_id uuid not null,
  shipment_id text not null,
  shipment_line_id text not null,
  order_id text not null,
  order_line_id text not null,
  variant_id text not null,
  item_id text not null,
  sales_uom_code text not null,
  quantity numeric(20,6) not null check (quantity > 0),
  quantity_base numeric(20,6) not null check (quantity_base > 0),
  cost_amount numeric(20,4) not null check (cost_amount >= 0),
  primary key (tenant_id, shipment_id, shipment_line_id),
  foreign key (tenant_id, shipment_id) references sales.customer_shipment(tenant_id,shipment_id),
  foreign key (tenant_id, order_id, order_line_id) references sales.customer_order_line(tenant_id,order_id,line_id),
  foreign key (tenant_id, variant_id) references catalog.vehicle_variant(tenant_id,variant_id),
  foreign key (tenant_id, item_id, sales_uom_code) references inventory.item_unit_of_measure(tenant_id,item_id,uom_code),
  check (quantity_base > 0)
);

create table inventory.customer_shipment_allocation (
  tenant_id uuid not null,
  shipment_id text not null,
  shipment_line_id text not null,
  allocation_id text not null,
  reservation_id text not null,
  outbound_entry_id text not null,
  organization_id text not null,
  from_bin_id text not null,
  item_id text not null,
  lot_id text,
  quantity_base numeric(20,6) not null check (quantity_base > 0),
  cost_amount numeric(20,4) not null check (cost_amount >= 0),
  primary key (tenant_id, shipment_id, allocation_id),
  unique (tenant_id, reservation_id),
  unique (tenant_id, outbound_entry_id),
  foreign key (tenant_id, shipment_id, shipment_line_id)
    references sales.customer_shipment_line(tenant_id,shipment_id,shipment_line_id),
  foreign key (tenant_id, reservation_id) references inventory.bulk_reservation(tenant_id,reservation_id),
  foreign key (tenant_id, outbound_entry_id) references inventory.bulk_inventory_entry(tenant_id,entry_id),
  foreign key (tenant_id, organization_id, from_bin_id)
    references inventory.warehouse_bin(tenant_id,organization_id,bin_id),
  foreign key (tenant_id, lot_id) references inventory.inventory_lot(tenant_id,lot_id)
);

create function sales.reject_customer_shipment_mutation()
returns trigger language plpgsql as $function$
begin
  raise exception using errcode='55000', message='posted customer shipment evidence is immutable';
end;
$function$;

create trigger customer_shipment_immutable
before update or delete on sales.customer_shipment
for each row execute function sales.reject_customer_shipment_mutation();

create trigger customer_shipment_line_immutable
before update or delete on sales.customer_shipment_line
for each row execute function sales.reject_customer_shipment_mutation();

create trigger customer_shipment_allocation_immutable
before update or delete on inventory.customer_shipment_allocation
for each row execute function sales.reject_customer_shipment_mutation();

create index customer_shipment_order_idx
  on sales.customer_shipment(tenant_id,order_id,posting_date,shipment_id);

commit;
