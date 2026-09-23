-- AUTHORED binding to the existing catalog, manual approval and delivery fence.
begin;
create table catalog.merchant_effect(
 tenant_id uuid not null,channel_code text not null check(channel_code='google_merchant'),
 delivery_key text not null,account_id text not null,offer_id text not null,
 generation bigint not null,source_sha256 text not null check(source_sha256~'^[0-9a-f]{64}$'),
 profile_sha256 text not null check(profile_sha256~'^[0-9a-f]{64}$'),
 request_sha256 text not null check(request_sha256~'^[0-9a-f]{64}$'),
 created_at timestamptz not null default clock_timestamp(),
 primary key(tenant_id,channel_code,delivery_key),
 foreign key(tenant_id,generation) references catalog.release_publication(tenant_id,generation),
 foreign key(tenant_id,delivery_key) references approval.request(tenant_id,request_id),
 foreign key(tenant_id,channel_code,delivery_key) references communication.outbound_delivery(tenant_id,channel_code,delivery_key)
);
create index merchant_effect_offer_idx on catalog.merchant_effect(tenant_id,account_id,offer_id);
create trigger merchant_effect_immutable before update or delete on catalog.merchant_effect for each row execute function catalog.release_immutable();
create table catalog.merchant_observation(
 tenant_id uuid not null,channel_code text not null default 'google_merchant' check(channel_code='google_merchant'),
 approval_id text not null,operation text not null check(operation in('INSERT','GET')),
 result_state text not null check(result_state in('OBSERVED','NOT_FOUND','REJECTED','UNKNOWN')),
 confirmed boolean not null,
 response bytea not null check(octet_length(response)<=32768),
 response_sha256 text not null check(response_sha256~'^[0-9a-f]{64}$'),
 observed_at timestamptz not null default clock_timestamp(),
 primary key(tenant_id,approval_id,operation,result_state,response_sha256),
 foreign key(tenant_id,channel_code,approval_id) references catalog.merchant_effect(tenant_id,channel_code,delivery_key),
 check(not confirmed or result_state='OBSERVED'),
 check(operation<>'INSERT' or confirmed)
);
create trigger merchant_observation_immutable before update or delete on catalog.merchant_observation for each row execute function catalog.release_immutable();
commit;
