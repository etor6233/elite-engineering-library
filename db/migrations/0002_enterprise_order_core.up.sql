begin;

create schema org;
create schema sales;

create table org.organization (
  tenant_id uuid not null,
  organization_id text not null,
  parent_organization_id text,
  organization_code text not null,
  display_name text not null,
  organization_type text not null
    check (organization_type in ('enterprise', 'franchisor', 'franchisee', 'factory', 'warehouse', 'store', 'service_center')),
  status text not null default 'active'
    check (status in ('provisioning', 'active', 'suspended', 'closed')),
  created_at timestamptz not null default clock_timestamp(),
  updated_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, organization_id),
  foreign key (tenant_id) references platform.tenant (tenant_id),
  foreign key (tenant_id, parent_organization_id)
    references org.organization (tenant_id, organization_id),
  unique (tenant_id, organization_code),
  check (length(organization_id) between 1 and 128),
  check (organization_code ~ '^[a-z][a-z0-9]*(-[a-z0-9]+)*$'),
  check (parent_organization_id is null or parent_organization_id <> organization_id)
);

create table sales.customer_order (
  tenant_id uuid not null,
  order_id text not null,
  organization_id text not null,
  customer_principal_id text not null,
  state text not null
    check (state in ('draft', 'placed', 'confirmed', 'paid', 'allocated', 'delivered', 'cancelled')),
  currency text not null check (currency ~ '^[A-Z]{3}$'),
  total_minor_units bigint not null check (total_minor_units >= 0),
  version bigint not null check (version > 0),
  created_at timestamptz not null default clock_timestamp(),
  updated_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, order_id),
  foreign key (tenant_id, organization_id)
    references org.organization (tenant_id, organization_id),
  check (length(order_id) between 1 and 128),
  check (length(customer_principal_id) between 1 and 256),
  check (updated_at >= created_at)
);

create index customer_order_org_state_idx
  on sales.customer_order (tenant_id, organization_id, state, created_at desc, order_id);

create index customer_order_customer_idx
  on sales.customer_order (tenant_id, customer_principal_id, created_at desc, order_id);

commit;
