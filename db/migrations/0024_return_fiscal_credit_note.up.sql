begin;
create table fiscal.return_credit_note_link (
  tenant_id uuid not null,
  request_id text not null,
  disposition_id text not null,
  original_invoice_id text not null,
  credit_invoice_id text not null,
  currency text not null check(currency ~ '^[A-Z]{3}$'),
  total_minor_units bigint not null check(total_minor_units>0),
  request_sha256_hex text not null check(request_sha256_hex ~ '^[0-9a-f]{64}$'),
  created_at timestamptz not null default clock_timestamp(),
  primary key(tenant_id,request_id),
  unique(tenant_id,disposition_id),
  unique(tenant_id,credit_invoice_id),
  foreign key(tenant_id,request_id) references sales.return_effect_request(tenant_id,request_id),
  foreign key(tenant_id,disposition_id) references sales.return_disposition(tenant_id,disposition_id),
  foreign key(tenant_id,original_invoice_id) references fiscal.invoice(tenant_id,invoice_id),
  foreign key(tenant_id,credit_invoice_id) references fiscal.invoice(tenant_id,invoice_id),
  check(original_invoice_id<>credit_invoice_id)
);
create or replace function fiscal.prevent_return_credit_link_mutation() returns trigger language plpgsql as $$ begin raise exception using errcode='23514',message='return credit note link is immutable'; end $$;
create trigger return_credit_note_link_immutable before update or delete on fiscal.return_credit_note_link for each row execute function fiscal.prevent_return_credit_link_mutation();
commit;
