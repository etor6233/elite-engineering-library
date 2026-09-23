begin;
create table communication.whatsapp_campaign(
 tenant_id uuid not null, campaign_id text not null, organization_id text not null,
 creator_subject text not null, campaign_sha256 text not null,
 payload jsonb not null, created_at timestamptz not null default clock_timestamp(),
 primary key(tenant_id,campaign_id),
 unique(tenant_id,campaign_id,campaign_sha256),
 foreign key(tenant_id,organization_id)references org.organization(tenant_id,organization_id),
 check(length(campaign_id)between 16 and 128),
 check(length(creator_subject)between 1 and 128),
 check(campaign_sha256~'^[0-9a-f]{64}$'),
 check(jsonb_typeof(payload)='object' and octet_length(payload::text)<=131072)
);
create trigger whatsapp_campaign_immutable before update or delete on communication.whatsapp_campaign
 for each row execute function catalog.release_immutable();
create table communication.whatsapp_campaign_member(
 tenant_id uuid not null,campaign_id text not null,lead_id text not null,
 member_sha256 text not null check(member_sha256~'^[0-9a-f]{64}$'),
 primary key(tenant_id,campaign_id,lead_id),
 foreign key(tenant_id,campaign_id)references communication.whatsapp_campaign(tenant_id,campaign_id),
 foreign key(tenant_id,lead_id)references crm.lead(tenant_id,lead_id)
);
create trigger whatsapp_campaign_member_immutable before update or delete on communication.whatsapp_campaign_member
 for each row execute function catalog.release_immutable();
create table communication.whatsapp_campaign_stop(
 tenant_id uuid not null,campaign_id text not null,campaign_sha256 text not null,
 actor_subject text not null,reason text not null,
 stopped_at timestamptz not null default clock_timestamp(),
 primary key(tenant_id,campaign_id),
 foreign key(tenant_id,campaign_id,campaign_sha256)references communication.whatsapp_campaign(tenant_id,campaign_id,campaign_sha256),
 check(length(actor_subject)between 1 and 128),check(length(reason)between 1 and 2048)
);
create trigger whatsapp_campaign_stop_immutable before update or delete on communication.whatsapp_campaign_stop
 for each row execute function catalog.release_immutable();
alter table communication.whatsapp_schedule alter column appointment_id drop not null;
alter table communication.whatsapp_schedule add column campaign_id text,add column lead_id text,add column campaign_step integer,
 add constraint whatsapp_schedule_campaign_member_fk foreign key(tenant_id,campaign_id,lead_id)
 references communication.whatsapp_campaign_member(tenant_id,campaign_id,lead_id),
 add constraint whatsapp_schedule_source_check check(
 (appointment_id is not null and appointment_version>0 and campaign_id is null and lead_id is null and campaign_step is null)
 or(appointment_id is null and appointment_version=0 and campaign_id is not null and lead_id is not null and campaign_step between 1 and 3)),
 add constraint whatsapp_schedule_campaign_step_unique unique(tenant_id,campaign_id,lead_id,campaign_step);
commit;
