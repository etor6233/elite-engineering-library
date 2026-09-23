create table inventory.serial_reservation (
  tenant_id uuid not null,
  reservation_id text not null,
  organization_id text not null,
  stock_unit_id text not null,
  variant_id text not null,
  demand_kind text not null check (demand_kind in ('customer-order','service','transfer-outbound','manual')),
  demand_id text not null,
  demand_line_id text not null default '',
  status text not null check (status in ('reservation','released','consumed')),
  cancellation_disallowed boolean not null default false,
  expires_at timestamptz,
  version bigint not null check (version > 0),
  created_at timestamptz not null default clock_timestamp(),
  updated_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id,reservation_id),
  foreign key (tenant_id,organization_id) references org.organization(tenant_id,organization_id),
  foreign key (tenant_id,stock_unit_id) references inventory.stock_unit(tenant_id,stock_unit_id),
  foreign key (tenant_id,variant_id) references catalog.vehicle_variant(tenant_id,variant_id),
  check (expires_at is null or expires_at > created_at)
);

create unique index serial_reservation_active_stock_uidx
  on inventory.serial_reservation(tenant_id,stock_unit_id)
  where status='reservation';

create unique index serial_reservation_active_demand_uidx
  on inventory.serial_reservation(tenant_id,demand_kind,demand_id,demand_line_id)
  where status='reservation';

create index serial_reservation_scope_idx
  on inventory.serial_reservation(tenant_id,organization_id,variant_id,status,expires_at);

create table inventory.serial_transfer (
  tenant_id uuid not null,
  transfer_id text not null,
  from_organization_id text not null,
  to_organization_id text not null,
  expected_receipt_at timestamptz not null,
  state text not null check (state in ('draft','released','in-transit','received','cancelled')),
  version bigint not null check (version > 0),
  created_at timestamptz not null default clock_timestamp(),
  updated_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id,transfer_id),
  foreign key (tenant_id,from_organization_id) references org.organization(tenant_id,organization_id),
  foreign key (tenant_id,to_organization_id) references org.organization(tenant_id,organization_id),
  check (from_organization_id<>to_organization_id)
);

create table inventory.serial_transfer_unit (
  tenant_id uuid not null,
  transfer_id text not null,
  stock_unit_id text not null,
  variant_id text not null,
  primary key (tenant_id,transfer_id,stock_unit_id),
  foreign key (tenant_id,transfer_id) references inventory.serial_transfer(tenant_id,transfer_id) on delete restrict,
  foreign key (tenant_id,stock_unit_id) references inventory.stock_unit(tenant_id,stock_unit_id),
  foreign key (tenant_id,variant_id) references catalog.vehicle_variant(tenant_id,variant_id)
);

create index serial_transfer_inbound_idx
  on inventory.serial_transfer(tenant_id,to_organization_id,state,expected_receipt_at);

create function inventory.serial_atp(
  p_tenant uuid,
  p_organization text,
  p_variant text,
  p_horizon timestamptz
) returns table(
  available_inventory bigint,
  scheduled_receipt bigint,
  gross_requirement bigint,
  available_to_promise bigint
) language sql stable as $$
with available as (
  select count(*)::bigint quantity
  from inventory.stock_unit
  where tenant_id=p_tenant and organization_id=p_organization and variant_id=p_variant and state='available'
), inbound as (
  select count(*)::bigint quantity
  from inventory.serial_transfer t
  join inventory.serial_transfer_unit u using(tenant_id,transfer_id)
  where t.tenant_id=p_tenant and t.to_organization_id=p_organization and u.variant_id=p_variant
    and t.state='in-transit' and t.expected_receipt_at<=p_horizon
), demand as (
  select coalesce(sum(l.quantity-case when l.allocated_stock_unit_id is null then 0 else 1 end),0)::bigint quantity
  from sales.customer_order o
  join sales.customer_order_line l using(tenant_id,order_id)
  where o.tenant_id=p_tenant and o.organization_id=p_organization and l.variant_id=p_variant
    and o.state in ('placed','confirmed')
)
select a.quantity,i.quantity,d.quantity,greatest(0::bigint,a.quantity+i.quantity-d.quantity)
from available a cross join inbound i cross join demand d;
$$;

comment on function inventory.serial_atp(uuid,text,text,timestamptz) is
  'Portable serial-stock ATP adaptation governed by Microsoft BC Available to Promise: available inventory plus due inbound transfer receipts minus unallocated placed demand. It is not full Business Central planning/CTP.';
