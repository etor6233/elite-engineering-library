begin;
-- AUTHORED idempotency/recovery metadata at the existing royalty owner.
-- The financial owner and outbox share this exact transaction; no second ledger.
create table platform.finance_command_receipt(
 tenant_id uuid not null,
 organization_id text not null,
 requested_by_subject text not null check(length(requested_by_subject) between 1 and 256),
 request_key text not null check(request_key ~ '^[A-Za-z0-9_-]{16,128}$'),
 action text not null check(action in ('royalty_policy','royalty_accrue','royalty_open','royalty_close','royalty_reverse','royalty_reconcile','accounting_account','accounting_period','accounting_journal','accounting_post','accounting_reverse','accounting_close')),
 payload_sha256_hex text not null check(payload_sha256_hex ~ '^[a-f0-9]{64}$'),
 receipt_raw bytea not null check(octet_length(receipt_raw) between 1 and 32768),
 receipt_sha256_hex text not null check(receipt_sha256_hex ~ '^[a-f0-9]{64}$'),
 created_at timestamptz not null default clock_timestamp(),
 primary key(tenant_id,organization_id,requested_by_subject,request_key),
 foreign key(tenant_id,organization_id) references org.organization(tenant_id,organization_id)
);
create trigger royalty_command_receipt_immutable before update or delete on platform.finance_command_receipt for each row execute function royalty.prevent_financial_history_mutation();
commit;
