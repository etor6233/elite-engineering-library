begin;

create schema platform;
create schema audit;

create table platform.tenant (
  tenant_id uuid primary key,
  tenant_code text not null unique,
  legal_name text not null,
  display_name text not null,
  status text not null default 'active'
    check (status in ('provisioning', 'active', 'suspended', 'closed')),
  created_at timestamptz not null default clock_timestamp(),
  closed_at timestamptz,
  check ((status = 'closed') = (closed_at is not null)),
  check (tenant_code ~ '^[a-z][a-z0-9]*(-[a-z0-9]+)*$')
);

create table platform.business_profile_version (
  tenant_id uuid not null,
  profile_version bigint not null check (profile_version > 0),
  schema_version text not null,
  profile jsonb not null,
  sha256_hex text not null,
  status text not null
    check (status in ('draft', 'admitted', 'active', 'retired', 'rejected')),
  created_by_subject text not null,
  created_at timestamptz not null default clock_timestamp(),
  activated_at timestamptz,
  primary key (tenant_id, profile_version),
  foreign key (tenant_id) references platform.tenant (tenant_id),
  check (jsonb_typeof(profile) = 'object'),
  check (sha256_hex ~ '^[0-9a-f]{64}$'),
  check ((status = 'active') = (activated_at is not null))
);

create unique index business_profile_one_active_per_tenant_uq
  on platform.business_profile_version (tenant_id)
  where status = 'active';

create table platform.idempotency_record (
  tenant_id uuid not null,
  scope text not null,
  idempotency_key text not null,
  request_sha256_hex text not null,
  status text not null
    check (status in ('processing', 'completed', 'failed_terminal')),
  response_code integer,
  response_body jsonb,
  resource_type text,
  resource_id text,
  locked_until timestamptz,
  created_at timestamptz not null default clock_timestamp(),
  expires_at timestamptz not null,
  primary key (tenant_id, scope, idempotency_key),
  foreign key (tenant_id) references platform.tenant (tenant_id),
  check (length(scope) between 1 and 128),
  check (length(idempotency_key) between 8 and 200),
  check (request_sha256_hex ~ '^[0-9a-f]{64}$'),
  check (expires_at > created_at),
  check (
    (status = 'processing' and response_code is null)
    or
    (status in ('completed', 'failed_terminal') and response_code between 100 and 599)
  )
);

create index idempotency_expiry_idx
  on platform.idempotency_record (expires_at);

create table platform.outbox_event (
  tenant_id uuid not null,
  event_id uuid not null,
  aggregate_type text not null,
  aggregate_id text not null,
  aggregate_version bigint not null check (aggregate_version > 0),
  event_type text not null,
  schema_version integer not null check (schema_version > 0),
  occurred_at timestamptz not null,
  trace_id text,
  payload jsonb not null,
  headers jsonb not null default '{}'::jsonb,
  available_at timestamptz not null default clock_timestamp(),
  claimed_by text,
  claimed_until timestamptz,
  attempts integer not null default 0 check (attempts >= 0),
  published_at timestamptz,
  last_error_code text,
  primary key (tenant_id, event_id),
  foreign key (tenant_id) references platform.tenant (tenant_id),
  unique (tenant_id, aggregate_type, aggregate_id, aggregate_version, event_type),
  check (jsonb_typeof(payload) = 'object'),
  check (jsonb_typeof(headers) = 'object'),
  check ((claimed_by is null) = (claimed_until is null)),
  check (published_at is null or published_at >= occurred_at)
);

create index outbox_ready_idx
  on platform.outbox_event (available_at, occurred_at, event_id)
  where published_at is null;

create table platform.consumer_inbox (
  tenant_id uuid not null,
  consumer_name text not null,
  event_id uuid not null,
  received_at timestamptz not null default clock_timestamp(),
  completed_at timestamptz,
  result_sha256_hex text,
  primary key (tenant_id, consumer_name, event_id),
  foreign key (tenant_id) references platform.tenant (tenant_id),
  check (length(consumer_name) between 1 and 128),
  check (completed_at is null or completed_at >= received_at),
  check (result_sha256_hex is null or result_sha256_hex ~ '^[0-9a-f]{64}$')
);

create index inbox_incomplete_idx
  on platform.consumer_inbox (received_at)
  where completed_at is null;

create table platform.job (
  tenant_id uuid not null,
  job_id uuid not null,
  queue text not null,
  job_type text not null,
  schema_version integer not null check (schema_version > 0),
  payload jsonb not null,
  priority smallint not null default 0,
  available_at timestamptz not null default clock_timestamp(),
  attempts integer not null default 0 check (attempts >= 0),
  max_attempts integer not null check (max_attempts between 1 and 100),
  claimed_by text,
  claimed_until timestamptz,
  completed_at timestamptz,
  terminal_error_code text,
  created_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, job_id),
  foreign key (tenant_id) references platform.tenant (tenant_id),
  check (jsonb_typeof(payload) = 'object'),
  check ((claimed_by is null) = (claimed_until is null)),
  check (not (completed_at is not null and terminal_error_code is not null)),
  check (completed_at is null or completed_at >= created_at)
);

create index job_ready_idx
  on platform.job (queue, priority desc, available_at, job_id)
  where completed_at is null and terminal_error_code is null;

create table audit.event (
  audit_sequence bigint generated always as identity primary key,
  tenant_id uuid not null,
  event_id uuid not null,
  occurred_at timestamptz not null default clock_timestamp(),
  actor_subject text not null,
  action text not null,
  resource_type text not null,
  resource_id text not null,
  trace_id text,
  source_ip inet,
  decision text not null check (decision in ('allowed', 'denied', 'system')),
  reason_code text,
  evidence jsonb not null default '{}'::jsonb,
  foreign key (tenant_id) references platform.tenant (tenant_id),
  unique (tenant_id, event_id),
  check (jsonb_typeof(evidence) = 'object')
);

create index audit_tenant_time_idx
  on audit.event (tenant_id, occurred_at desc, audit_sequence desc);

create function audit.reject_mutation()
returns trigger
language plpgsql
as $function$
begin
  raise exception using
    errcode = '55000',
    message = 'audit.event is append-only';
end;
$function$;

create trigger audit_event_reject_update_or_delete
before update or delete on audit.event
for each row execute function audit.reject_mutation();

commit;
