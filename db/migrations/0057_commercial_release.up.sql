begin;

-- AUTHORED checkpoint persistence; no shipment, payment or stock journal.
create table sales.commercial_release_receipt (
 tenant_id uuid not null,
 release_id text not null,
 organization_id text not null,
 handover_id text not null,
 order_id text not null,
 payment_attempt_id text not null,
 observation_sha256_hex text not null check(observation_sha256_hex ~ '^[0-9a-f]{64}$'),
 observation_generation bigint not null check(observation_generation > 0),
 handover_version integer not null check(handover_version > 0),
 acceptance_sha256_hex text not null check(acceptance_sha256_hex ~ '^[0-9a-f]{64}$'),
 checklist_id text not null,
 checklist_version integer not null check(checklist_version > 0),
 contract_id text not null,
 contract_sha256_hex text not null check(contract_sha256_hex ~ '^[0-9a-f]{64}$'),
 effect text not null check(effect='COMMIT_COMMERCIAL_RELEASE_RECEIPT'),
 released_by_subject text not null check(length(released_by_subject) between 1 and 256),
 recorded_at timestamptz not null,
 valid_until timestamptz not null,
 primary key(tenant_id,release_id),
 unique(tenant_id,handover_id),
 unique(tenant_id,organization_id,order_id),
 foreign key(tenant_id,handover_id) references sales.delivery_handover_preparation(tenant_id,handover_id),
 foreign key(tenant_id,organization_id) references org.organization(tenant_id,organization_id),
 foreign key(tenant_id,order_id) references sales.customer_order(tenant_id,order_id),
 foreign key(tenant_id,payment_attempt_id) references payment.provider_observation(tenant_id,payment_attempt_id),
 check(valid_until > recorded_at)
);

create trigger commercial_release_receipt_immutable before update or delete on sales.commercial_release_receipt
for each row execute function communication.reject_outbound_delivery_event_mutation();

commit;
