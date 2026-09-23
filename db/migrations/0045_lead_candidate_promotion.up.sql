begin;

create table integration.lead_promotion (
  tenant_id uuid not null,
  provider text not null,
  provider_event_id text not null,
  lead_id text not null,
  consent_id text not null,
  mapping_version text not null,
  field_mapping jsonb not null,
  source_payload_sha256 text not null check (source_payload_sha256 ~ '^[0-9a-f]{64}$'),
  decision_at timestamptz not null,
  created_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, provider, provider_event_id),
  unique (tenant_id, lead_id),
  unique (tenant_id, consent_id),
  foreign key (tenant_id, provider, provider_event_id) references integration.lead_candidate (tenant_id, provider, provider_event_id),
  foreign key (tenant_id, lead_id) references crm.lead (tenant_id, lead_id),
  foreign key (tenant_id, consent_id) references crm.consent_evidence (tenant_id, consent_id),
  check (length(mapping_version) between 1 and 128),
  check (jsonb_typeof(field_mapping)='object')
);

create trigger lead_promotion_immutable before update or delete on integration.lead_promotion
for each row execute function integration.reject_lead_ingress_mutation();

commit;
