begin;
-- AUTHORED relationship/receipt at the existing accounting owner, not a ledger.
create table accounting.fx_journal_receipt(
 tenant_id uuid not null,
 organization_id text not null,
 conversion_id text not null,
 journal_id text not null,
 requested_by_subject text not null check(length(requested_by_subject) between 1 and 256),
 request_key text not null check(request_key ~ '^[A-Za-z0-9_-]{16,128}$'),
 request_sha256_hex text not null check(request_sha256_hex ~ '^[0-9a-f]{64}$'),
 receipt_raw bytea not null check(octet_length(receipt_raw) between 1 and 16384),
 receipt_sha256_hex text not null check(receipt_sha256_hex ~ '^[0-9a-f]{64}$'),
 recorded_at timestamptz not null,
 primary key(tenant_id,conversion_id),
 unique(tenant_id,request_key),
 unique(tenant_id,journal_id),
 foreign key(tenant_id,organization_id) references org.organization(tenant_id,organization_id),
 foreign key(tenant_id,conversion_id) references accounting.fx_conversion_receipt(tenant_id,conversion_id),
 foreign key(tenant_id,journal_id) references accounting.journal(tenant_id,journal_id)
);
create trigger fx_journal_receipt_immutable before update or delete on accounting.fx_journal_receipt for each row execute function accounting.prevent_posted_history_mutation();
commit;
