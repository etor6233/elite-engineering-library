begin;

create table communication.outbound_delivery (
  tenant_id uuid not null references platform.tenant(tenant_id),
  channel_code text not null check (channel_code ~ '^[a-z][a-z0-9_]{0,31}$'),
  delivery_key text not null check (length(delivery_key) between 16 and 128),
  request_sha256_hex text not null check (request_sha256_hex ~ '^[0-9a-f]{64}$'),
  recipient_hmac text not null check (recipient_hmac ~ '^[0-9a-f]{64}$'),
  state text not null check (state in ('sending','accepted','unknown','failed_terminal')),
  attempt_count integer not null check (attempt_count=1),
  locked_until timestamptz,
  provider_message_hmac text check (provider_message_hmac is null or provider_message_hmac ~ '^[0-9a-f]{64}$'),
  evidence_sha256_hex text check (evidence_sha256_hex is null or evidence_sha256_hex ~ '^[0-9a-f]{64}$'),
  failure_code text,
  accepted_at timestamptz,
  created_at timestamptz not null default clock_timestamp(),
  updated_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, channel_code, delivery_key),
  check ((state='sending' and locked_until is not null) or state<>'sending'),
  check ((state='accepted' and provider_message_hmac is not null and evidence_sha256_hex is not null and accepted_at is not null) or state<>'accepted'),
  check (failure_code is null or failure_code ~ '^[A-Z][A-Z0-9_]{0,127}$')
);

create table communication.outbound_delivery_event (
  tenant_id uuid not null,
  channel_code text not null,
  delivery_key text not null,
  sequence bigint not null check (sequence > 0),
  state text not null check (state in ('sending','accepted','unknown','failed_terminal')),
  evidence_sha256_hex text check (evidence_sha256_hex is null or evidence_sha256_hex ~ '^[0-9a-f]{64}$'),
  failure_code text,
  recorded_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id,channel_code,delivery_key,sequence),
  foreign key (tenant_id,channel_code,delivery_key) references communication.outbound_delivery(tenant_id,channel_code,delivery_key),
  check (failure_code is null or failure_code ~ '^[A-Z][A-Z0-9_]{0,127}$')
);

create function communication.reject_outbound_delivery_event_mutation()
returns trigger language plpgsql as $function$
begin
  raise exception using errcode='55000', message='outbound delivery event is immutable';
end;
$function$;

create trigger outbound_delivery_event_immutable before update or delete on communication.outbound_delivery_event
for each row execute function communication.reject_outbound_delivery_event_mutation();

create index outbound_delivery_unknown_idx on communication.outbound_delivery(updated_at) where state='unknown';

commit;
