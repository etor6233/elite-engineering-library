begin;

create table integration.provider_connection (
  tenant_id uuid not null,
  connection_id text not null,
  provider_code text not null,
  organization_id text,
  secret_ref text not null,
  state text not null default 'active' check (state in ('active', 'disabled')),
  created_at timestamptz not null default clock_timestamp(),
  updated_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, connection_id),
  foreign key (tenant_id) references platform.tenant (tenant_id),
  foreign key (tenant_id, organization_id) references org.organization (tenant_id, organization_id),
  unique (tenant_id, connection_id, provider_code),
  check (connection_id ~ '^[a-z0-9][a-z0-9._-]{0,127}$'),
  check (provider_code ~ '^[a-z0-9][a-z0-9._-]{0,63}$'),
  check (secret_ref ~ '^[A-Z][A-Z0-9_]{2,127}$'),
  check (updated_at >= created_at)
);

create table integration.webhook_event (
  tenant_id uuid not null,
  connection_id text not null,
  provider_code text not null,
  provider_event_id text not null,
  event_type text not null,
  body_sha256_hex text not null,
  payload jsonb not null,
  received_at timestamptz not null default clock_timestamp(),
  processed_at timestamptz,
  state text not null default 'received' check (state in ('received', 'processed', 'failed')),
  last_error_code text,
  primary key (tenant_id, connection_id, provider_event_id),
  foreign key (tenant_id, connection_id, provider_code)
    references integration.provider_connection (tenant_id, connection_id, provider_code),
  check (provider_event_id <> '' and length(provider_event_id) <= 200),
  check (event_type <> '' and length(event_type) <= 200),
  check (body_sha256_hex ~ '^[0-9a-f]{64}$'),
  check (jsonb_typeof(payload) = 'object'),
  check ((state = 'processed') = (processed_at is not null))
);

create index provider_webhook_ready_idx
  on integration.webhook_event (tenant_id, provider_code, received_at, provider_event_id)
  where state = 'received';

commit;
