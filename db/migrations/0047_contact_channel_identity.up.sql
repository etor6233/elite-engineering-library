begin;

create table communication.contact_channel_binding (
  tenant_id uuid not null references platform.tenant(tenant_id),
  channel_code text not null check (channel_code ~ '^[a-z][a-z0-9_]{0,31}$'),
  external_id_hmac text not null check (external_id_hmac ~ '^[0-9a-f]{64}$'),
  lead_id text not null,
  subject_id text not null check (length(subject_id) between 1 and 256),
  pii_allowed boolean not null,
  state text not null check (state in ('active','revoked')),
  policy_version text not null check (length(policy_version) between 1 and 128),
  evidence_sha256_hex text not null check (evidence_sha256_hex ~ '^[0-9a-f]{64}$'),
  effective_at timestamptz not null,
  version bigint not null check (version > 0),
  created_at timestamptz not null default clock_timestamp(),
  updated_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, channel_code, external_id_hmac),
  foreign key (tenant_id, lead_id) references crm.lead(tenant_id, lead_id),
  check (state <> 'revoked' or pii_allowed = false),
  check (updated_at >= created_at)
);

create table communication.contact_channel_binding_decision (
  tenant_id uuid not null,
  request_id text not null check (length(request_id) between 1 and 128),
  request_sha256_hex text not null check (request_sha256_hex ~ '^[0-9a-f]{64}$'),
  channel_code text not null,
  external_id_hmac text not null,
  binding_version bigint not null check (binding_version > 0),
  lead_id text not null,
  subject_id text not null,
  pii_allowed boolean not null,
  state text not null check (state in ('active','revoked')),
  policy_version text not null,
  evidence_sha256_hex text not null check (evidence_sha256_hex ~ '^[0-9a-f]{64}$'),
  effective_at timestamptz not null,
  recorded_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, request_id),
  unique (tenant_id, channel_code, external_id_hmac, binding_version),
  foreign key (tenant_id, channel_code, external_id_hmac) references communication.contact_channel_binding(tenant_id, channel_code, external_id_hmac),
  check (state <> 'revoked' or pii_allowed = false)
);

create function communication.reject_contact_binding_decision_mutation()
returns trigger language plpgsql as $function$
begin
  raise exception using errcode='55000', message='contact binding evidence is immutable';
end;
$function$;

create trigger contact_binding_decision_immutable before update or delete on communication.contact_channel_binding_decision
for each row execute function communication.reject_contact_binding_decision_mutation();

create index contact_binding_lead_idx on communication.contact_channel_binding(tenant_id, lead_id) where state='active';

commit;
