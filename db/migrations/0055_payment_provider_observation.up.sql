begin;

-- AUTHORED persistence glue for the existing payment owner and official SDKs.
-- It does not create a second balance, capture or money-movement ledger.
create table payment.provider_checkout (
  tenant_id uuid not null,
  payment_attempt_id text not null,
  provider_code text not null check (provider_code in ('stripe','mercadopago')),
  connection_id text not null,
  account_ref text not null check (length(account_ref) between 1 and 200),
  request_sha256_hex text not null check (request_sha256_hex ~ '^[0-9a-f]{64}$'),
  session_id text,
  payment_reference text,
  checkout_url text,
  checkout_state text,
  live_mode boolean not null,
  expires_at timestamptz not null,
  created_at timestamptz not null default clock_timestamp(),
  updated_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id,payment_attempt_id),
  foreign key (tenant_id,payment_attempt_id) references payment.payment_attempt(tenant_id,payment_attempt_id),
  foreign key (tenant_id,connection_id) references integration.provider_connection(tenant_id,connection_id),
  unique (tenant_id,provider_code,session_id),
  unique (tenant_id,provider_code,payment_reference),
  check (session_id is null or length(session_id) between 1 and 200),
  check (payment_reference is null or length(payment_reference) between 1 and 200),
  check (checkout_url is null or length(checkout_url) between 1 and 4096),
  check ((session_id is null and checkout_url is null and checkout_state is null)
      or (session_id is not null and checkout_state is not null and
          (checkout_url is not null or (provider_code='stripe' and checkout_state in ('complete','expired')))))
);

create table payment.provider_observation (
  tenant_id uuid not null,
  payment_attempt_id text not null,
  order_id text not null,
  organization_id text not null,
  provider_code text not null check (provider_code in ('stripe','mercadopago')),
  provider_reference text not null check (length(provider_reference) between 1 and 200),
  currency text not null check (currency ~ '^[A-Z]{3}$'),
  amount_minor_units bigint not null check (amount_minor_units > 0),
  received_minor_units bigint not null default 0 check (received_minor_units >= 0),
  refunded_minor_units bigint not null default 0 check (refunded_minor_units >= 0),
  provider_status text not null default 'unobserved',
  live_mode boolean not null,
  account_ref text not null check (length(account_ref) between 1 and 200),
  evidence_sha256_hex text check (evidence_sha256_hex ~ '^[0-9a-f]{64}$'),
  observed_at timestamptz,
  hold boolean not null default true,
  hold_reason text not null default 'OBSERVATION_PENDING',
  generation bigint not null default 1 check (generation > 0),
  updated_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id,payment_attempt_id),
  foreign key (tenant_id,payment_attempt_id) references payment.payment_attempt(tenant_id,payment_attempt_id),
  foreign key (tenant_id,order_id) references sales.customer_order(tenant_id,order_id),
  unique (tenant_id,provider_code,provider_reference),
  check (hold or (observed_at is not null and evidence_sha256_hex is not null and received_minor_units=amount_minor_units and refunded_minor_units=0))
);

create table payment.provider_observation_event (
  tenant_id uuid not null,
  payment_attempt_id text not null,
  generation bigint not null,
  evidence_sha256_hex text not null check (evidence_sha256_hex ~ '^[0-9a-f]{64}$'),
  result_code text not null,
  observed_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id,payment_attempt_id,generation),
  foreign key (tenant_id,payment_attempt_id) references payment.provider_observation(tenant_id,payment_attempt_id)
);
create trigger payment_observation_event_immutable before update or delete on payment.provider_observation_event
for each row execute function communication.reject_outbound_delivery_event_mutation();
create index payment_observation_hold_idx on payment.provider_observation(tenant_id,provider_code,updated_at) where hold;

commit;
