begin;

create table inventory.sales_warehouse_binding (
  tenant_id uuid not null,
  request_id text not null,
  variant_id text not null,
  item_id text not null,
  sales_uom_code text not null,
  version bigint not null default 1 check (version = 1),
  created_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, variant_id),
  foreign key (tenant_id, variant_id)
    references catalog.vehicle_variant(tenant_id, variant_id),
  foreign key (tenant_id, item_id, sales_uom_code)
    references inventory.item_unit_of_measure(tenant_id, item_id, uom_code),
  unique (tenant_id, request_id),
  check (length(request_id) between 1 and 128)
);

create index sales_warehouse_binding_item_idx
  on inventory.sales_warehouse_binding(tenant_id, item_id, sales_uom_code);

create function inventory.reject_sales_warehouse_binding_mutation()
returns trigger language plpgsql as $function$
begin
  raise exception using errcode='55000',
    message='sales warehouse binding is immutable; create a new sellable variant for a changed item or unit contract';
end;
$function$;

create trigger sales_warehouse_binding_immutable
before update or delete on inventory.sales_warehouse_binding
for each row execute function inventory.reject_sales_warehouse_binding_mutation();

create table inventory.warehouse_pick_request (
  tenant_id uuid not null,
  request_id text not null,
  activity_id text not null,
  organization_id text not null,
  ship_bin_id text not null,
  item_id text not null,
  demand_kind text not null,
  demand_id text not null,
  demand_line_id text not null,
  quantity numeric(20,6) not null check (quantity > 0),
  requested_handling_uom text not null default '',
  resolved_handling_uom text not null,
  allow_breakbulk boolean not null,
  use_fefo boolean not null,
  allow_dedicated boolean not null,
  assigned_to text not null default '',
  created_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, request_id),
  foreign key (tenant_id, activity_id, organization_id)
    references inventory.warehouse_activity(tenant_id, activity_id, organization_id),
  foreign key (tenant_id, organization_id, ship_bin_id)
    references inventory.warehouse_bin(tenant_id, organization_id, bin_id),
  foreign key (tenant_id, item_id, resolved_handling_uom)
    references inventory.item_unit_of_measure(tenant_id, item_id, uom_code),
  unique (tenant_id, activity_id),
  check (length(request_id) between 1 and 128)
);

create function inventory.reject_warehouse_pick_request_mutation()
returns trigger language plpgsql as $function$
begin
  raise exception using errcode='55000', message='warehouse pick request evidence is immutable';
end;
$function$;

create trigger warehouse_pick_request_immutable
before update or delete on inventory.warehouse_pick_request
for each row execute function inventory.reject_warehouse_pick_request_mutation();

commit;
