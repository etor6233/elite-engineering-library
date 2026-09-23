begin;
create table inventory.workspace_reservation_request(
 tenant_id uuid not null, organization_id text not null, request_id text not null,
 actor_id text not null, payload_sha256 text not null check(payload_sha256 ~ '^[a-f0-9]{64}$'),
 reservation_id text not null, created_at timestamptz not null default clock_timestamp(),
 primary key(tenant_id,organization_id,request_id),
 foreign key(tenant_id,reservation_id) references inventory.bulk_reservation(tenant_id,reservation_id),
 foreign key(tenant_id,organization_id) references org.organization(tenant_id,organization_id),
 check(length(request_id) between 1 and 128),check(length(actor_id) between 1 and 256)
);
create trigger workspace_reservation_request_immutable before update or delete on inventory.workspace_reservation_request for each row execute function inventory.reject_bulk_ledger_mutation();
commit;
