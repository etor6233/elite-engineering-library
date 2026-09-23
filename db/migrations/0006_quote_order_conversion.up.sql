begin;

alter table sales.quotation
  add column order_id text,
  add constraint quotation_order_fk foreign key (tenant_id, order_id)
    references sales.customer_order (tenant_id, order_id),
  add constraint quotation_acceptance_order_check check (
    (state = 'accepted' and order_id is not null) or
    (state <> 'accepted' and order_id is null)
  );

create unique index quotation_order_unique_idx
  on sales.quotation (tenant_id, order_id)
  where order_id is not null;

create table sales.quotation_acceptance (
  tenant_id uuid not null,
  quotation_id text not null,
  order_id text not null,
  customer_principal_id text not null,
  evidence_sha256_hex text not null check (evidence_sha256_hex ~ '^[0-9a-f]{64}$'),
  accepted_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, quotation_id),
  unique (tenant_id, order_id),
  foreign key (tenant_id, quotation_id) references sales.quotation (tenant_id, quotation_id),
  foreign key (tenant_id, order_id) references sales.customer_order (tenant_id, order_id),
  foreign key (tenant_id, customer_principal_id) references crm.customer_profile (tenant_id, customer_principal_id)
);

commit;
