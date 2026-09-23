begin;

create schema if not exists finops;

create table finops.resource (
  tenant_id uuid not null,
  provider text not null,
  resource_id text not null,
  kind text not null,
  region text not null default '',
  tags jsonb not null,
  estimated_monthly_minor_units bigint not null check (estimated_monthly_minor_units >= 0),
  created_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, provider, resource_id),
  foreign key (tenant_id) references platform.tenant (tenant_id),
  check (provider ~ '^[a-z][a-z0-9_-]{0,31}$'),
  check (kind ~ '^[a-z][a-z0-9_-]{0,63}$'),
  check (jsonb_typeof(tags) = 'object'),
  check (tags ? 'tenant' and tags ? 'environment' and tags ? 'owner')
);

create table finops.budget (
  tenant_id uuid not null,
  period text not null check (period ~ '^\d{4}-\d{2}$'),
  cap_minor_units bigint not null check (cap_minor_units >= 0),
  spent_minor_units bigint not null default 0 check (spent_minor_units >= 0),
  primary key (tenant_id, period),
  foreign key (tenant_id) references platform.tenant (tenant_id),
  check (spent_minor_units <= cap_minor_units)
);

create table finops.cost_entry (
  tenant_id uuid not null,
  provider text not null,
  resource_id text not null,
  period text not null check (period ~ '^\d{4}-\d{2}$'),
  amount_minor_units bigint not null check (amount_minor_units >= 0),
  recorded_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, provider, resource_id, period),
  foreign key (tenant_id) references platform.tenant (tenant_id)
);

commit;
