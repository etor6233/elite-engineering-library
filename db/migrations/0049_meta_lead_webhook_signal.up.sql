begin;

create table integration.meta_lead_webhook_batch (
  tenant_id uuid not null references platform.tenant(tenant_id),
  organization_id text not null,
  payload_sha256 text not null check (payload_sha256 ~ '^[0-9a-f]{64}$'),
  payload bytea not null check (octet_length(payload) between 1 and 1048576),
  received_at timestamptz not null,
  signal_count integer not null check (signal_count > 0),
  primary key (tenant_id,payload_sha256),
  foreign key (tenant_id,organization_id) references org.organization(tenant_id,organization_id)
);

create table integration.meta_lead_webhook_signal (
  tenant_id uuid not null references platform.tenant(tenant_id),
  leadgen_id text not null check (length(leadgen_id) between 1 and 256),
  organization_id text not null,
  batch_payload_sha256 text not null,
  value_sha256 text not null check (value_sha256 ~ '^[0-9a-f]{64}$'),
  form_id text not null check (length(form_id) between 1 and 256),
  page_id text not null check (length(page_id) between 1 and 256),
  ad_id text,
  ad_group_id text,
  occurred_at timestamptz not null,
  state text not null check (state='pending_retrieval'),
  primary key (tenant_id,leadgen_id),
  foreign key (tenant_id,organization_id) references org.organization(tenant_id,organization_id),
  foreign key (tenant_id,batch_payload_sha256) references integration.meta_lead_webhook_batch(tenant_id,payload_sha256)
);

create index meta_lead_webhook_pending_idx on integration.meta_lead_webhook_signal
  (tenant_id,state,occurred_at,leadgen_id);

create function integration.reject_meta_lead_webhook_mutation()
returns trigger language plpgsql as $function$
begin
  raise exception using errcode='55000', message='Meta lead webhook evidence is append-only';
end;
$function$;

create trigger meta_lead_webhook_batch_immutable before update or delete on integration.meta_lead_webhook_batch
for each row execute function integration.reject_meta_lead_webhook_mutation();

create trigger meta_lead_webhook_signal_immutable before update or delete on integration.meta_lead_webhook_signal
for each row execute function integration.reject_meta_lead_webhook_mutation();

commit;
