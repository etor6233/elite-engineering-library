begin;

create schema if not exists data;
create schema if not exists analytics;

-- ingest source ownership and contract
create table data.ingest_source (
  tenant_id uuid not null,
  source_id text not null,
  owner_organization_id text,
  contract_schema text not null,
  freshness_ttl interval not null,
  created_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, source_id),
  foreign key (tenant_id) references platform.tenant (tenant_id),
  check (source_id ~ '^[a-z][a-z0-9._-]{0,63}$'),
  check (length(contract_schema) between 1 and 200),
  check (freshness_ttl > interval '0')
);

-- immutable landing record with replay idempotency and event/received ordering
create table data.landing_record (
  tenant_id uuid not null,
  source_id text not null,
  external_id text not null,
  event_time timestamptz not null,
  received_at timestamptz not null default clock_timestamp(),
  payload jsonb not null,
  source_sha256 text not null,
  lineage text not null default '',
  state text not null default 'landed' check (state in ('landed', 'validated', 'rejected')),
  primary key (tenant_id, source_id, external_id, source_sha256),
  foreign key (tenant_id, source_id) references data.ingest_source (tenant_id, source_id),
  check (received_at >= event_time),
  check (jsonb_typeof(payload) = 'object'),
  check (source_sha256 ~ '^[0-9a-f]{64}$'),
  check (length(external_id) between 1 and 200)
);

create index landing_record_source_time_idx
  on data.landing_record (tenant_id, source_id, received_at);

-- versioned, immutable semantic metric definition
create table analytics.metric_definition (
  tenant_id uuid not null,
  metric_code text not null,
  version bigint not null check (version > 0),
  aggregation text not null check (aggregation in ('count', 'sum', 'min', 'max', 'avg')),
  source_ref text not null,
  dimensions jsonb not null default '[]'::jsonb,
  freshness_ttl interval not null,
  created_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, metric_code, version),
  foreign key (tenant_id) references platform.tenant (tenant_id),
  check (metric_code ~ '^[a-z][a-z0-9._-]{0,63}$'),
  check (length(source_ref) between 1 and 200),
  check (jsonb_typeof(dimensions) = 'array'),
  check (freshness_ttl > interval '0')
);

-- materialized metric snapshot
create table analytics.metric_snapshot (
  tenant_id uuid not null,
  metric_code text not null,
  version bigint not null,
  period_start timestamptz not null,
  period_end timestamptz not null,
  value numeric not null,
  dimensions jsonb not null default '{}'::jsonb,
  computed_at timestamptz not null default clock_timestamp(),
  reconciled boolean not null default false,
  primary key (tenant_id, metric_code, version, period_start, period_end),
  foreign key (tenant_id, metric_code, version)
    references analytics.metric_definition (tenant_id, metric_code, version),
  check (period_end > period_start),
  check (jsonb_typeof(dimensions) = 'object')
);

commit;
