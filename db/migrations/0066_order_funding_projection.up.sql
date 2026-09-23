begin;
-- AUTHORED base projection. Gross order pricing remains with sales; this view
-- introduces no discount rule, alternate order, payment or accounting entry.
create view payment.order_funding as
 select tenant_id,order_id,organization_id,currency,version as order_version,
 total_minor_units as gross_minor_units,0::bigint as gift_minor_units,
 0::bigint as discount_minor_units,total_minor_units as provider_due_minor_units,
 encode(sha256(convert_to('[]','UTF8')),'hex') as contributions_sha256
 from sales.customer_order;
create table payment.local_funding_receipt (
 tenant_id uuid not null,
 funding_id text not null,
 request_key text not null check(length(request_key) between 16 and 128),
 request_sha256 text not null check(request_sha256 ~ '^[0-9a-f]{64}$'),
 requested_by_subject text not null check(length(requested_by_subject) between 1 and 256),
 receipt_raw bytea not null check(octet_length(receipt_raw) between 1 and 65536),
 order_id text not null,
 organization_id text not null,
 currency text not null,
 gross_minor_units bigint not null check(gross_minor_units>0),
 gift_minor_units bigint not null check(gift_minor_units>=0),
 discount_minor_units bigint not null check(discount_minor_units>=0),
 provider_minor_units bigint not null check(provider_minor_units>=0),
 payment_attempt_id text,
 provider_evidence_sha256 text,
 allocation_sha256 text not null check(allocation_sha256 ~ '^[0-9a-f]{64}$'),
 receipt_sha256 text not null check(receipt_sha256 ~ '^[0-9a-f]{64}$'),
 receipt jsonb not null check(jsonb_typeof(receipt)='object' and pg_column_size(receipt)<=65536),
 created_at timestamptz not null default clock_timestamp(),
 primary key(tenant_id,funding_id),
 unique(tenant_id,request_key),
 foreign key(tenant_id,order_id) references sales.customer_order(tenant_id,order_id),
 foreign key(tenant_id,organization_id) references org.organization(tenant_id,organization_id),
 foreign key(tenant_id,payment_attempt_id) references payment.payment_attempt(tenant_id,payment_attempt_id),
 check(gross_minor_units::numeric=gift_minor_units::numeric+discount_minor_units::numeric+provider_minor_units::numeric),
 check((provider_minor_units=0 and payment_attempt_id is null and provider_evidence_sha256 is null) or
       (provider_minor_units>0 and payment_attempt_id is not null and provider_evidence_sha256 is not null and provider_evidence_sha256 ~ '^[0-9a-f]{64}$'))
);
create view payment.local_funding_evidence as select tenant_id,funding_id,order_id,organization_id,currency,gross_minor_units,gift_minor_units,discount_minor_units,allocation_sha256,receipt_raw,receipt_sha256,created_at from payment.local_funding_receipt where provider_minor_units=0 and payment_attempt_id is null and provider_evidence_sha256 is null;
create index local_funding_order on payment.local_funding_receipt(tenant_id,order_id,created_at desc);
create function payment.local_funding_immutable() returns trigger language plpgsql as $$
begin raise exception 'local funding observation is immutable';end $$;
create trigger local_funding_immutable before update or delete on payment.local_funding_receipt for each row execute function payment.local_funding_immutable();
commit;
