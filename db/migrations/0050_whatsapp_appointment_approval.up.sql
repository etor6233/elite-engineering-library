begin;

create table communication.whatsapp_appointment_approval (
  tenant_id uuid not null,
  delivery_key text not null check(length(delivery_key) between 16 and 128),
  channel_code text not null default 'whatsapp' check(channel_code='whatsapp'),
  appointment_id text not null,
  appointment_version bigint not null check(appointment_version>0),
  organization_id text not null,
  confirmation_event_id uuid not null,
  external_id_hmac text not null check(external_id_hmac ~ '^[0-9a-f]{64}$'),
  binding_version bigint not null check(binding_version>0),
  lead_id text not null,
  subject_id text not null,
  policy_version text not null check(length(policy_version) between 1 and 128),
  consent_id text not null,
  consent_purpose text not null check(length(consent_purpose) between 1 and 128),
  consent_evidence_sha256 text not null check(consent_evidence_sha256 ~ '^[0-9a-f]{64}$'),
  message_sha256 text not null check(message_sha256 ~ '^[0-9a-f]{64}$'),
  profile_sha256 text not null check(profile_sha256 ~ '^[0-9a-f]{64}$'),
  approval_request_sha256 text not null check(approval_request_sha256 ~ '^[0-9a-f]{64}$'),
  evidence_sha256 text not null check(evidence_sha256 ~ '^[0-9a-f]{64}$'),
  approved_by text not null check(length(approved_by) between 1 and 255),
  approved_at timestamptz not null default statement_timestamp(),
  expires_at timestamptz not null,
  primary key(tenant_id,delivery_key),
  unique(tenant_id,confirmation_event_id,channel_code),
  foreign key(tenant_id,appointment_id) references crm.appointment(tenant_id,appointment_id),
  foreign key(tenant_id,confirmation_event_id) references crm.appointment_transition(tenant_id,transition_id),
  foreign key(tenant_id,consent_id) references crm.consent_evidence(tenant_id,consent_id),
  foreign key(tenant_id,channel_code,external_id_hmac) references communication.contact_channel_binding(tenant_id,channel_code,external_id_hmac),
  check(expires_at>approved_at)
);

create function communication.reject_whatsapp_approval_mutation() returns trigger
language plpgsql as $$ begin raise exception using errcode='55000',message='WhatsApp approval evidence is immutable'; end $$;
create trigger whatsapp_approval_immutable before update or delete
on communication.whatsapp_appointment_approval for each row
execute function communication.reject_whatsapp_approval_mutation();

create index whatsapp_notification_consent_lookup_idx
on crm.consent_evidence(tenant_id,lead_id,purpose_code,occurred_at desc);

commit;
