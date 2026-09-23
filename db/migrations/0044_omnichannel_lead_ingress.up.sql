begin;

create table integration.lead_ingress_raw (
  tenant_id uuid not null,
  organization_id text not null,
  provider text not null check (provider ~ '^[a-z][a-z0-9_]{0,31}$'),
  provider_event_id text not null check (length(provider_event_id) between 1 and 256),
  event_type text not null,
  source text not null,
  schema_version text not null,
  occurred_at timestamptz not null,
  received_at timestamptz not null,
  payload_redacted bytea not null,
  source_payload_sha256 text not null check (source_payload_sha256 ~ '^[0-9a-f]{64}$'),
  stored_payload_sha256 text not null check (stored_payload_sha256 ~ '^[0-9a-f]{64}$'),
  redaction_profile text not null check (length(redaction_profile) between 1 and 128),
  state text not null check (state in ('normalized','rejected')),
  normalization_error_code text,
  primary key (tenant_id, provider, provider_event_id),
  foreign key (tenant_id, organization_id) references org.organization (tenant_id, organization_id),
  check (octet_length(payload_redacted) between 1 and 1048576),
  check ((state='normalized' and normalization_error_code is null) or (state='rejected' and length(normalization_error_code) between 1 and 128))
);

create index lead_ingress_received_idx on integration.lead_ingress_raw
  (tenant_id, provider, received_at desc, provider_event_id);

create table integration.lead_candidate (
  tenant_id uuid not null,
  provider text not null,
  provider_event_id text not null,
  organization_id text not null,
  provider_lead_id text not null,
  form_id text,
  campaign_id text,
  ad_group_id text,
  creative_id text,
  asset_group_id text,
  click_id text,
  lead_stage text,
  source_kind text,
  submitted_at timestamptz not null,
  is_test boolean not null,
  fields jsonb not null,
  contact_eligibility text not null check (contact_eligibility='pending_policy'),
  primary key (tenant_id, provider, provider_event_id),
  foreign key (tenant_id, provider, provider_event_id)
    references integration.lead_ingress_raw (tenant_id, provider, provider_event_id),
  foreign key (tenant_id, organization_id) references org.organization (tenant_id, organization_id),
  check (jsonb_typeof(fields)='array'),
  check (length(provider_lead_id) between 1 and 256)
);

create function integration.reject_lead_ingress_mutation()
returns trigger language plpgsql as $function$
begin
  raise exception using errcode='55000', message='lead ingress evidence is append-only';
end;
$function$;

create trigger lead_ingress_raw_immutable before update or delete on integration.lead_ingress_raw
for each row execute function integration.reject_lead_ingress_mutation();

create trigger lead_candidate_immutable before update or delete on integration.lead_candidate
for each row execute function integration.reject_lead_ingress_mutation();

commit;
