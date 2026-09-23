begin;

-- AUTHORED glue over the existing PostgreSQL turn and budget owners.
-- This is a lifetime token reservation cap, not a provider invoice estimate.
create table communication.conversation_budget (
 tenant_id uuid primary key references platform.tenant(tenant_id),
 cap bigint not null check(cap between 1 and 1000000000000),
 reserved bigint not null default 0 check(reserved>=0 and reserved<=cap)
);
create table communication.conversation_reservation (
 tenant_id uuid not null references communication.conversation_budget(tenant_id),
 channel_code text not null,
 provider_message_id text not null,
 generation integer not null check(generation between 1 and 100),
 tokens bigint not null check(tokens>0),
 created_at timestamptz not null default clock_timestamp(),
 primary key(tenant_id,channel_code,provider_message_id,generation)
);
create table communication.conversation_intent (
 tenant_id uuid not null,
 channel_code text not null,
 provider_message_id text not null,
 intent_sha256 text not null check(intent_sha256 ~ '^[a-f0-9]{64}$'),
 created_at timestamptz not null default clock_timestamp(),
 primary key(tenant_id,channel_code,provider_message_id),
 foreign key(tenant_id,channel_code,provider_message_id) references communication.conversation_turn(tenant_id,channel_code,provider_message_id)
);
create function communication.reject_conversation_governance_update()
returns trigger language plpgsql as $f$
begin
 raise exception using errcode='55000',message='conversation reservation and intent are immutable';
end;
$f$;
create trigger conversation_reservation_immutable before update on communication.conversation_reservation
for each row execute function communication.reject_conversation_governance_update();
create trigger conversation_intent_immutable before update on communication.conversation_intent
for each row execute function communication.reject_conversation_governance_update();
commit;
