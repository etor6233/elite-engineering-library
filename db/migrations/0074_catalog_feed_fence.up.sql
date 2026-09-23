begin;
-- AUTHORED immutable admission binding; existing outbound_delivery owns states.
create table catalog.release_feed_intent(
 tenant_id uuid not null,channel_code text not null check(channel_code='catalog_feed'),
 delivery_key text not null,receiver_id text not null,receiver_profile_sha256 text not null check(receiver_profile_sha256~'^[0-9a-f]{64}$'),
 generation bigint not null,actor text not null check(length(actor) between 1 and 128),
 request_sha256 text not null check(request_sha256~'^[0-9a-f]{64}$'),
 payload_sha256 text not null,source_sha256 text not null check(source_sha256~'^[0-9a-f]{64}$'),
 payload bytea not null,created_at timestamptz not null default clock_timestamp(),
 primary key(tenant_id,channel_code,delivery_key),unique(tenant_id,generation,receiver_id),
 foreign key(tenant_id,generation) references catalog.release_publication(tenant_id,generation),
 foreign key(tenant_id,channel_code,delivery_key) references communication.outbound_delivery(tenant_id,channel_code,delivery_key),
 check(octet_length(payload)<=524288 and payload_sha256=encode(sha256(payload),'hex'))
);
create trigger catalog_release_feed_immutable before update or delete on catalog.release_feed_intent for each row execute function catalog.release_immutable();
commit;
