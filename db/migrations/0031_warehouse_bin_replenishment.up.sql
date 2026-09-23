begin;

alter table inventory.warehouse_activity
  drop constraint warehouse_activity_activity_type_check,
  add constraint warehouse_activity_activity_type_v165_check
    check (activity_type in ('put-away','pick','movement'));

alter table inventory.bulk_reservation
  drop constraint bulk_reservation_demand_kind_check,
  add constraint bulk_reservation_demand_kind_v165_check
    check (demand_kind in ('customer-order','service','transfer-outbound','manual','warehouse-replenishment'));

create table inventory.warehouse_replenishment_request (
  tenant_id uuid not null,
  request_id text not null,
  activity_id text not null,
  organization_id text not null,
  to_bin_id text not null,
  item_id text not null,
  use_fefo boolean not null,
  planned_quantity numeric(20,6) not null check (planned_quantity > 0),
  created_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, request_id),
  foreign key (tenant_id, activity_id)
    references inventory.warehouse_activity (tenant_id, activity_id),
  foreign key (tenant_id, organization_id)
    references org.organization (tenant_id, organization_id),
  foreign key (tenant_id, organization_id, to_bin_id)
    references inventory.warehouse_bin (tenant_id, organization_id, bin_id),
  foreign key (tenant_id, item_id)
    references inventory.stock_item (tenant_id, item_id),
  unique (tenant_id, activity_id),
  check (length(request_id) between 1 and 128)
);

commit;
