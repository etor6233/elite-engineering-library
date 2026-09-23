begin;
create table communication.whatsapp_status_batch (
 tenant_id uuid not null,
 delivery_key text not null,
 observations_sha256 text not null check(observations_sha256 ~ '^[0-9a-f]{64}$'),
 send_receipt_sha256 text not null check(send_receipt_sha256 ~ '^[0-9a-f]{64}$'),
 verification_receipt jsonb not null check(jsonb_typeof(verification_receipt)='object' and octet_length(verification_receipt::text)<=65536),
 observed_by text not null check(length(observed_by) between 1 and 255),
 recorded_at timestamptz not null default statement_timestamp(),
 primary key(tenant_id,delivery_key,observations_sha256),
 foreign key(tenant_id,delivery_key) references communication.whatsapp_appointment_approval(tenant_id,delivery_key)
);
create table communication.whatsapp_status_observation (
 tenant_id uuid not null,
 delivery_key text not null,
 event_key_sha256 text not null check(event_key_sha256 ~ '^[0-9a-f]{64}$'),
 event_sha256 text not null check(event_sha256 ~ '^[0-9a-f]{64}$'),
 provider_timestamp bigint not null check(provider_timestamp between 1 and 999999999999),
 provider_status text not null check(provider_status in ('sent','delivered','read','failed','deleted')),
 observations_sha256 text not null,
 primary key(tenant_id,delivery_key,event_key_sha256),
 foreign key(tenant_id,delivery_key,observations_sha256) references communication.whatsapp_status_batch(tenant_id,delivery_key,observations_sha256)
);
create index whatsapp_status_latest_idx on communication.whatsapp_status_observation(tenant_id,delivery_key,provider_timestamp desc);
create trigger whatsapp_status_batch_immutable before update or delete on communication.whatsapp_status_batch
 for each row execute function communication.reject_whatsapp_approval_mutation();
create trigger whatsapp_status_observation_immutable before update or delete on communication.whatsapp_status_observation
 for each row execute function communication.reject_whatsapp_approval_mutation();
commit;
