begin;
create table franchise.network_command_receipt(
 tenant_id uuid not null references platform.tenant(tenant_id),
 command_id text not null check(length(command_id)between 1 and 128),
 action text not null check(action in('create-organization','transition-organization','create-agreement','transition-agreement')),
 scope_organization_id text not null check(length(scope_organization_id)<=128),
 actor_subject text not null check(length(actor_subject)between 1 and 256),
 request_sha256 text not null check(request_sha256~'^[0-9a-f]{64}$'),
 entity_payload jsonb not null check(jsonb_typeof(entity_payload)='object'),
 entity_sha256 text not null check(entity_sha256~'^[0-9a-f]{64}$'),
 recorded_at timestamptz not null default clock_timestamp(),
 primary key(tenant_id,command_id)
);
create function franchise.reject_network_receipt_mutation()returns trigger language plpgsql as $$begin raise exception 'network receipt is immutable';end$$;
create trigger network_receipt_immutable before update or delete on franchise.network_command_receipt for each row execute function franchise.reject_network_receipt_mutation();
create index network_receipt_actor_idx on franchise.network_command_receipt(tenant_id,actor_subject,scope_organization_id,recorded_at,command_id);
commit;
