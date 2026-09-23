create schema if not exists fiscal;

create table fiscal.point_of_sale (
  tenant_id uuid not null,
  point_of_sale_id text not null,
  organization_id text not null,
  taxpayer_cuit text not null check (taxpayer_cuit ~ '^[0-9]{11}$'),
  environment text not null check (environment in ('homologation','production')),
  point_of_sale_number integer not null check (point_of_sale_number between 1 and 99999),
  active boolean not null default true,
  version bigint not null check (version > 0),
  active_invoice_id text,
  lease_owner text,
  lease_until timestamptz,
  created_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id,point_of_sale_id),
  foreign key (tenant_id,organization_id) references org.organization(tenant_id,organization_id),
  unique (tenant_id,taxpayer_cuit,environment,point_of_sale_number),
  check ((active_invoice_id is null and lease_owner is null and lease_until is null) or (active_invoice_id is not null and length(lease_owner) between 1 and 128 and lease_until is not null))
);

create table fiscal.invoice (
  tenant_id uuid not null,
  invoice_id text not null,
  organization_id text not null,
  order_id text not null,
  payment_attempt_id text not null,
  point_of_sale_id text not null,
  voucher_type integer not null check (voucher_type between 1 and 999),
  concept integer not null check (concept between 1 and 3),
  recipient_document_type integer not null check (recipient_document_type between 0 and 999),
  recipient_document text not null check (recipient_document ~ '^[0-9]{0,20}$'),
  recipient_vat_condition_id integer not null check (recipient_vat_condition_id between 1 and 999),
  currency text not null check (currency ~ '^[A-Z]{3}$'),
  provider_currency text not null check (provider_currency ~ '^[A-Z]{3}$'),
  total_minor_units bigint not null check (total_minor_units > 0),
  net_minor_units bigint not null check (net_minor_units >= 0),
  vat_minor_units bigint not null check (vat_minor_units >= 0),
  exempt_minor_units bigint not null check (exempt_minor_units >= 0),
  non_taxed_minor_units bigint not null check (non_taxed_minor_units >= 0),
  other_tax_minor_units bigint not null check (other_tax_minor_units >= 0),
  issued_on date not null,
  service_from date,
  service_until date,
  payment_due_on date,
  status text not null check (status in ('queued','authorizing','reconcile_required','authorized','rejected')),
  voucher_number bigint check (voucher_number > 0),
  cae text check (cae ~ '^[0-9]{8,20}$'),
  cae_expires_on date,
  request_hash text not null check (request_hash ~ '^[a-f0-9]{64}$'),
  idempotency_key text not null check (length(idempotency_key) between 16 and 128),
  last_response_hash text check (last_response_hash ~ '^[a-f0-9]{64}$'),
  provider_codes jsonb not null default '[]'::jsonb check (jsonb_typeof(provider_codes)='array'),
  attempt_count integer not null default 0 check (attempt_count >= 0),
  version bigint not null check (version > 0),
  created_at timestamptz not null default clock_timestamp(),
  updated_at timestamptz not null default clock_timestamp(),
  authorized_at timestamptz,
  primary key (tenant_id,invoice_id),
  foreign key (tenant_id,organization_id) references org.organization(tenant_id,organization_id),
  foreign key (tenant_id,order_id) references sales.customer_order(tenant_id,order_id),
  foreign key (tenant_id,payment_attempt_id) references payment.payment_attempt(tenant_id,payment_attempt_id),
  foreign key (tenant_id,point_of_sale_id) references fiscal.point_of_sale(tenant_id,point_of_sale_id),
  unique (tenant_id,idempotency_key),
  unique (tenant_id,payment_attempt_id,voucher_type),
  check (total_minor_units=net_minor_units+vat_minor_units+exempt_minor_units+non_taxed_minor_units+other_tax_minor_units),
  check ((concept=1 and service_from is null and service_until is null and payment_due_on is null) or (concept in (2,3) and service_from is not null and service_until >= service_from and payment_due_on >= service_until)),
  check ((status='authorized' and voucher_number is not null and cae is not null and cae_expires_on is not null and authorized_at is not null) or status<>'authorized'),
  check (status not in ('authorizing','reconcile_required') or voucher_number is not null)
);

create unique index fiscal_voucher_sequence_idx on fiscal.invoice(tenant_id,point_of_sale_id,voucher_type,voucher_number) where voucher_number is not null;
create index fiscal_invoice_work_idx on fiscal.invoice(status,created_at,tenant_id,invoice_id) where status in ('queued','reconcile_required');
create index fiscal_invoice_order_idx on fiscal.invoice(tenant_id,organization_id,order_id);

create table fiscal.invoice_vat (
  tenant_id uuid not null,
  invoice_id text not null,
  vat_id integer not null check (vat_id between 1 and 999),
  base_minor_units bigint not null check (base_minor_units > 0),
  amount_minor_units bigint not null check (amount_minor_units >= 0),
  primary key (tenant_id,invoice_id,vat_id),
  foreign key (tenant_id,invoice_id) references fiscal.invoice(tenant_id,invoice_id)
);

create table fiscal.invoice_other_tax (
  tenant_id uuid not null,
  invoice_id text not null,
  tax_id integer not null check (tax_id between 1 and 999),
  description text not null check (length(description) between 1 and 80),
  base_minor_units bigint not null check (base_minor_units >= 0),
  rate_basis_points integer not null check (rate_basis_points between 0 and 100000),
  amount_minor_units bigint not null check (amount_minor_units >= 0),
  primary key (tenant_id,invoice_id,tax_id),
  foreign key (tenant_id,invoice_id) references fiscal.invoice(tenant_id,invoice_id)
);

create or replace function fiscal.validate_invoice_details() returns trigger language plpgsql as $$
declare
  target_tenant uuid := case when tg_op='DELETE' then old.tenant_id else new.tenant_id end;
  target_invoice text := case when tg_op='DELETE' then old.invoice_id else new.invoice_id end;
  expected_net bigint; expected_vat bigint; expected_other bigint;
  actual_base bigint; actual_vat bigint; actual_other bigint; vat_rows bigint; tax_rows bigint;
begin
  select net_minor_units,vat_minor_units,other_tax_minor_units into expected_net,expected_vat,expected_other from fiscal.invoice where tenant_id=target_tenant and invoice_id=target_invoice;
  if not found then return null; end if;
  select coalesce(sum(base_minor_units),0),coalesce(sum(amount_minor_units),0),count(*) into actual_base,actual_vat,vat_rows from fiscal.invoice_vat where tenant_id=target_tenant and invoice_id=target_invoice;
  select coalesce(sum(amount_minor_units),0),count(*) into actual_other,tax_rows from fiscal.invoice_other_tax where tenant_id=target_tenant and invoice_id=target_invoice;
  if actual_vat<>expected_vat or actual_other<>expected_other or (expected_vat>0 and (vat_rows=0 or actual_base<>expected_net)) or (expected_vat=0 and vat_rows<>0) or (expected_other>0 and tax_rows=0) or (expected_other=0 and tax_rows<>0) then raise exception 'fiscal detail totals do not match invoice'; end if;
  return null;
end; $$;
create constraint trigger fiscal_invoice_details_valid after insert or update on fiscal.invoice deferrable initially deferred for each row execute function fiscal.validate_invoice_details();
create constraint trigger fiscal_invoice_vat_valid after insert or update or delete on fiscal.invoice_vat deferrable initially deferred for each row execute function fiscal.validate_invoice_details();
create constraint trigger fiscal_invoice_other_tax_valid after insert or update or delete on fiscal.invoice_other_tax deferrable initially deferred for each row execute function fiscal.validate_invoice_details();

create or replace function fiscal.prevent_fiscal_detail_mutation() returns trigger language plpgsql as $$ begin raise exception 'immutable fiscal detail'; end; $$;
create trigger fiscal_invoice_vat_immutable before update or delete on fiscal.invoice_vat for each row execute function fiscal.prevent_fiscal_detail_mutation();
create trigger fiscal_invoice_other_tax_immutable before update or delete on fiscal.invoice_other_tax for each row execute function fiscal.prevent_fiscal_detail_mutation();

create table fiscal.issuance_attempt (
  tenant_id uuid not null,
  attempt_id text not null,
  invoice_id text not null,
  phase text not null check (phase in ('claim','last-authorized','assign-number','consult','authorize','finish','defer')),
  voucher_number bigint,
  request_hash text not null check (request_hash ~ '^[a-f0-9]{64}$'),
  response_hash text check (response_hash ~ '^[a-f0-9]{64}$'),
  outcome text not null check (outcome in ('started','observed','authorized','rejected','deferred')),
  provider_codes jsonb not null default '[]'::jsonb check (jsonb_typeof(provider_codes)='array'),
  error_class text,
  recorded_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id,attempt_id),
  foreign key (tenant_id,invoice_id) references fiscal.invoice(tenant_id,invoice_id)
);

create index fiscal_attempt_invoice_idx on fiscal.issuance_attempt(tenant_id,invoice_id,recorded_at,attempt_id);
create or replace function fiscal.prevent_attempt_mutation() returns trigger language plpgsql as $$ begin raise exception 'immutable fiscal attempt history'; end; $$;
create trigger fiscal_attempt_immutable before update or delete on fiscal.issuance_attempt for each row execute function fiscal.prevent_attempt_mutation();
