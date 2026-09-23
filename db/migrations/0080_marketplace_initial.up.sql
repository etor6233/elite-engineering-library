-- AUTHORED persistence binding. Existing catalog/approval/fence remain owners.
begin;
alter table catalog.marketplace_effect add column operation text not null default 'LEGACY_MUTATION'
 check(operation in ('LEGACY_MUTATION','PRICE','STOCK','PAUSE','RESUME','MEDIA','CREATE'));
create table catalog.marketplace_observation (
 tenant_id uuid not null, approval_id text not null, channel_code text not null default 'mercadolibre_catalog' check(channel_code='mercadolibre_catalog'),
 operation text not null check(operation in ('MEDIA','CREATE')),
 provider_id text not null, media_sha256 text not null check(media_sha256~'^[0-9a-f]{64}$'),
 observation jsonb not null check(jsonb_typeof(observation)='object'),
 observation_sha256 text not null check(observation_sha256~'^[0-9a-f]{64}$'),
 created_at timestamptz not null default clock_timestamp(),
 primary key(tenant_id,approval_id),
 foreign key(tenant_id,channel_code,approval_id) references catalog.marketplace_effect(tenant_id,channel_code,delivery_key)
);
create trigger marketplace_observation_immutable before update or delete on catalog.marketplace_observation for each row execute function catalog.release_immutable();
commit;
