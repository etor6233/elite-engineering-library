begin;
-- AUTHORED records at the existing accounting owner. No new ledger or posting.
create table accounting.fx_rate_snapshot(
 tenant_id uuid not null,
 profile_id text not null,
 profile_revision integer not null check(profile_revision>0),
 organization_id text not null,
 profile_sha256_hex text not null check(profile_sha256_hex ~ '^[0-9a-f]{64}$'),
 source_sha256_hex text not null check(source_sha256_hex ~ '^[0-9a-f]{64}$'),
 profile_raw bytea not null check(octet_length(profile_raw) between 1 and 1048576),
 source_raw bytea not null check(octet_length(source_raw) between 1 and 1048576),
 valid_from timestamptz not null,
 valid_until timestamptz not null check(valid_until>valid_from),
 primary key(tenant_id,profile_id,profile_revision),
 unique(tenant_id,profile_id,profile_revision,profile_sha256_hex),
 foreign key(tenant_id,organization_id) references org.organization(tenant_id,organization_id)
);
create trigger fx_rate_snapshot_immutable before update or delete on accounting.fx_rate_snapshot for each row execute function accounting.prevent_posted_history_mutation();
create table accounting.fx_conversion_receipt(
 tenant_id uuid not null,
 conversion_id text not null,
 organization_id text not null,
 requested_by_subject text not null check(length(requested_by_subject) between 1 and 256),
 request_key text not null check(request_key ~ '^[A-Za-z0-9_-]{16,128}$'),
 request_sha256_hex text not null check(request_sha256_hex ~ '^[0-9a-f]{64}$'),
 profile_id text not null,
 profile_revision integer not null,
 profile_sha256_hex text not null,
 receipt_raw bytea not null check(octet_length(receipt_raw) between 1 and 16384),
 receipt_sha256_hex text not null check(receipt_sha256_hex ~ '^[0-9a-f]{64}$'),
 recorded_at timestamptz not null,
 primary key(tenant_id,conversion_id),
 unique(tenant_id,request_key),
 foreign key(tenant_id,organization_id) references org.organization(tenant_id,organization_id),
 foreign key(tenant_id,profile_id,profile_revision,profile_sha256_hex) references accounting.fx_rate_snapshot(tenant_id,profile_id,profile_revision,profile_sha256_hex)
);
create trigger fx_conversion_receipt_immutable before update or delete on accounting.fx_conversion_receipt for each row execute function accounting.prevent_posted_history_mutation();
commit;
