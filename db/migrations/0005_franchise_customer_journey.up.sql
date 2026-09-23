begin;

create table org.public_location (
  tenant_id uuid not null,
  organization_id text not null,
  city text not null,
  region text not null,
  country text not null check (country ~ '^[A-Z]{2}$'),
  contact_phone text,
  contact_email text,
  published boolean not null default false,
  sort_order integer not null default 0,
  updated_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, organization_id),
  foreign key (tenant_id, organization_id) references org.organization (tenant_id, organization_id),
  check (contact_phone is null or contact_phone ~ '^\+[1-9][0-9]{7,14}$'),
  check (contact_email is null or contact_email = lower(contact_email))
);

alter table crm.lead
  add column assigned_subject text,
  add column version bigint not null default 1 check (version > 0);

create function crm.enforce_lead_transition() returns trigger language plpgsql as $$
begin
  if new.lifecycle_state <> old.lifecycle_state and not (
    (old.lifecycle_state = 'new' and new.lifecycle_state in ('contacted','lost')) or
    (old.lifecycle_state = 'contacted' and new.lifecycle_state in ('qualified','lost')) or
    (old.lifecycle_state = 'qualified' and new.lifecycle_state in ('converted','lost'))
  ) then
    raise check_violation using message = 'invalid lead lifecycle transition';
  end if;
  return new;
end $$;

create trigger enforce_lead_transition
before update of lifecycle_state on crm.lead
for each row execute function crm.enforce_lead_transition();

create table crm.appointment (
  tenant_id uuid not null,
  appointment_id text not null,
  organization_id text not null,
  lead_id text not null,
  customer_principal_id text,
  model_id text,
  appointment_kind text not null check (appointment_kind in ('consultation','test-drive','delivery','service')),
  starts_at timestamptz not null,
  state text not null check (state in ('requested','confirmed','completed','cancelled','no-show')),
  version bigint not null check (version > 0),
  created_at timestamptz not null default clock_timestamp(),
  updated_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, appointment_id),
  foreign key (tenant_id, organization_id) references org.organization (tenant_id, organization_id),
  foreign key (tenant_id, lead_id) references crm.lead (tenant_id, lead_id),
  foreign key (tenant_id, customer_principal_id) references crm.customer_profile (tenant_id, customer_principal_id),
  foreign key (tenant_id, model_id) references catalog.vehicle_model (tenant_id, model_id),
  check (starts_at > created_at),
  check (updated_at >= created_at)
);

create table sales.quotation (
  tenant_id uuid not null,
  quotation_id text not null,
  organization_id text not null,
  lead_id text not null,
  customer_principal_id text,
  variant_id text not null,
  price_book_id text not null,
  currency text not null check (currency ~ '^[A-Z]{3}$'),
  total_minor_units bigint not null check (total_minor_units > 0),
  valid_until timestamptz not null,
  state text not null check (state in ('issued','accepted','expired','withdrawn','converted')),
  version bigint not null check (version > 0),
  created_at timestamptz not null default clock_timestamp(),
  updated_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, quotation_id),
  foreign key (tenant_id, organization_id) references org.organization (tenant_id, organization_id),
  foreign key (tenant_id, lead_id) references crm.lead (tenant_id, lead_id),
  foreign key (tenant_id, customer_principal_id) references crm.customer_profile (tenant_id, customer_principal_id),
  foreign key (tenant_id, price_book_id, variant_id) references pricing.price_book_entry (tenant_id, price_book_id, variant_id),
  check (valid_until > created_at),
  check (updated_at >= created_at)
);

create table sales.delivery_handover (
  tenant_id uuid not null,
  handover_id text not null,
  organization_id text not null,
  order_id text not null,
  customer_principal_id text not null,
  stock_unit_id text not null,
  state text not null check (state in ('prepared','presented','accepted','rejected')),
  acceptance_evidence_sha256_hex text,
  customer_accepted_at timestamptz,
  version bigint not null check (version > 0),
  created_at timestamptz not null default clock_timestamp(),
  updated_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, handover_id),
  foreign key (tenant_id, organization_id) references org.organization (tenant_id, organization_id),
  foreign key (tenant_id, order_id) references sales.customer_order (tenant_id, order_id),
  foreign key (tenant_id, stock_unit_id) references inventory.stock_unit (tenant_id, stock_unit_id),
  check ((state = 'accepted') = (customer_accepted_at is not null)),
  check ((state = 'accepted') = (acceptance_evidence_sha256_hex is not null)),
  check (acceptance_evidence_sha256_hex is null or acceptance_evidence_sha256_hex ~ '^[0-9a-f]{64}$'),
  check (updated_at >= created_at)
);

create index public_location_published_idx on org.public_location (tenant_id, published, sort_order, organization_id);
create index lead_org_updated_idx on crm.lead (tenant_id, organization_id, updated_at desc, lead_id);
create index appointment_org_start_idx on crm.appointment (tenant_id, organization_id, starts_at, appointment_id);
create index appointment_customer_idx on crm.appointment (tenant_id, customer_principal_id, starts_at desc, appointment_id);
create index quotation_customer_idx on sales.quotation (tenant_id, customer_principal_id, created_at desc, quotation_id);
create index handover_customer_idx on sales.delivery_handover (tenant_id, customer_principal_id, created_at desc, handover_id);

commit;
