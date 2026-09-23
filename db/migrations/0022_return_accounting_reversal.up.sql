begin;

create table accounting.return_effect_posting (
  tenant_id uuid not null,
  request_id text not null,
  disposition_id text not null,
  source_order_id text not null,
  remedy text not null check (remedy in ('refund','exchange')),
  original_journal_id text not null,
  reversal_journal_id text not null,
  currency text not null check (currency ~ '^[A-Z]{3}$'),
  reversed_minor_units bigint not null check (reversed_minor_units > 0),
  status text not null check (status='posted'),
  result_sha256_hex text not null check (result_sha256_hex ~ '^[0-9a-f]{64}$'),
  posted_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id,request_id),
  unique (tenant_id,disposition_id),
  unique (tenant_id,reversal_journal_id),
  foreign key (tenant_id,request_id) references sales.return_effect_request(tenant_id,request_id),
  foreign key (tenant_id,disposition_id) references sales.return_disposition(tenant_id,disposition_id),
  foreign key (tenant_id,source_order_id) references sales.customer_order(tenant_id,order_id),
  foreign key (tenant_id,original_journal_id) references accounting.journal(tenant_id,journal_id),
  foreign key (tenant_id,reversal_journal_id) references accounting.journal(tenant_id,journal_id),
  check (original_journal_id<>reversal_journal_id)
);

create or replace function accounting.prevent_return_effect_posting_mutation() returns trigger language plpgsql as $$
begin raise exception using errcode='23514',message='return accounting posting is immutable'; end $$;
create trigger return_effect_posting_immutable before update or delete on accounting.return_effect_posting for each row execute function accounting.prevent_return_effect_posting_mutation();
create index return_effect_posting_source_idx on accounting.return_effect_posting(tenant_id,source_order_id,posted_at desc,request_id);

commit;
