begin;

create schema catalog;
create schema crm;
create schema partner;
create schema procurement;
create schema factory;
create schema inventory;
create schema pricing;
create schema payment;
create schema logistics;
create schema service_ops;
create schema franchise;
create schema integration;
create schema marketing;
create schema communication;

create table catalog.vehicle_model (
  tenant_id uuid not null,
  model_id text not null,
  model_code text not null,
  display_name text not null,
  vehicle_class text not null check (vehicle_class in ('motorcycle', 'bicycle', 'scooter', 'utility', 'other')),
  lifecycle_state text not null check (lifecycle_state in ('draft', 'active', 'retired')),
  publicly_visible boolean not null default false,
  specification jsonb not null default '{}'::jsonb,
  created_at timestamptz not null default clock_timestamp(),
  updated_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, model_id),
  foreign key (tenant_id) references platform.tenant (tenant_id),
  unique (tenant_id, model_code),
  check (jsonb_typeof(specification) = 'object'),
  check (updated_at >= created_at)
);

create table catalog.vehicle_variant (
  tenant_id uuid not null,
  variant_id text not null,
  model_id text not null,
  variant_code text not null,
  display_name text not null,
  battery_specification jsonb not null,
  homologation_state text not null default 'unknown'
    check (homologation_state in ('unknown', 'pending', 'approved', 'rejected', 'expired')),
  lifecycle_state text not null check (lifecycle_state in ('draft', 'active', 'retired')),
  primary key (tenant_id, variant_id),
  foreign key (tenant_id, model_id) references catalog.vehicle_model (tenant_id, model_id),
  unique (tenant_id, variant_code),
  check (jsonb_typeof(battery_specification) = 'object')
);

create table crm.customer_profile (
  tenant_id uuid not null,
  customer_principal_id text not null,
  display_name text not null,
  email_normalized text,
  phone_e164 text,
  status text not null default 'active' check (status in ('active', 'restricted', 'deleted')),
  created_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, customer_principal_id),
  foreign key (tenant_id) references platform.tenant (tenant_id),
  unique nulls not distinct (tenant_id, email_normalized),
  check (email_normalized is null or email_normalized = lower(email_normalized)),
  check (phone_e164 is null or phone_e164 ~ '^\+[1-9][0-9]{7,14}$')
);

create table crm.lead (
  tenant_id uuid not null,
  lead_id text not null,
  organization_id text not null,
  customer_principal_id text,
  model_id text,
  lifecycle_state text not null check (lifecycle_state in ('new', 'qualified', 'contacted', 'converted', 'lost')),
  source_code text not null,
  contact_payload jsonb not null,
  consent_required boolean not null default true,
  created_at timestamptz not null default clock_timestamp(),
  updated_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, lead_id),
  foreign key (tenant_id, organization_id) references org.organization (tenant_id, organization_id),
  foreign key (tenant_id, customer_principal_id) references crm.customer_profile (tenant_id, customer_principal_id),
  foreign key (tenant_id, model_id) references catalog.vehicle_model (tenant_id, model_id),
  check (jsonb_typeof(contact_payload) = 'object'),
  check (updated_at >= created_at)
);

create table crm.consent_evidence (
  tenant_id uuid not null,
  consent_id text not null,
  lead_id text,
  customer_principal_id text,
  purpose_code text not null,
  policy_version text not null,
  decision text not null check (decision in ('granted', 'denied', 'withdrawn')),
  occurred_at timestamptz not null,
  evidence_sha256_hex text not null check (evidence_sha256_hex ~ '^[0-9a-f]{64}$'),
  primary key (tenant_id, consent_id),
  foreign key (tenant_id, lead_id) references crm.lead (tenant_id, lead_id),
  foreign key (tenant_id, customer_principal_id) references crm.customer_profile (tenant_id, customer_principal_id),
  check (lead_id is not null or customer_principal_id is not null)
);

create table partner.supplier (
  tenant_id uuid not null,
  supplier_id text not null,
  supplier_code text not null,
  legal_name text not null,
  status text not null check (status in ('onboarding', 'active', 'suspended', 'terminated')),
  qualification jsonb not null default '{}'::jsonb,
  primary key (tenant_id, supplier_id),
  foreign key (tenant_id) references platform.tenant (tenant_id),
  unique (tenant_id, supplier_code),
  check (jsonb_typeof(qualification) = 'object')
);

create table procurement.purchase_order (
  tenant_id uuid not null,
  purchase_order_id text not null,
  supplier_id text not null,
  destination_organization_id text not null,
  state text not null check (state in ('draft', 'submitted', 'accepted', 'in-production', 'shipped', 'received', 'cancelled')),
  currency text not null check (currency ~ '^[A-Z]{3}$'),
  total_minor_units bigint not null check (total_minor_units >= 0),
  version bigint not null check (version > 0),
  created_at timestamptz not null default clock_timestamp(),
  updated_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, purchase_order_id),
  foreign key (tenant_id, supplier_id) references partner.supplier (tenant_id, supplier_id),
  foreign key (tenant_id, destination_organization_id) references org.organization (tenant_id, organization_id),
  check (updated_at >= created_at)
);

create table factory.production_unit (
  tenant_id uuid not null,
  production_unit_id text not null,
  purchase_order_id text not null,
  variant_id text not null,
  serial_number text not null,
  vin text,
  battery_serial_number text,
  state text not null check (state in ('planned', 'assembly', 'quality', 'released', 'shipped', 'received', 'rejected')),
  milestone_payload jsonb not null default '{}'::jsonb,
  updated_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, production_unit_id),
  foreign key (tenant_id, purchase_order_id) references procurement.purchase_order (tenant_id, purchase_order_id),
  foreign key (tenant_id, variant_id) references catalog.vehicle_variant (tenant_id, variant_id),
  unique (tenant_id, serial_number),
  unique nulls not distinct (tenant_id, vin),
  unique nulls not distinct (tenant_id, battery_serial_number),
  check (jsonb_typeof(milestone_payload) = 'object')
);

create table inventory.stock_unit (
  tenant_id uuid not null,
  stock_unit_id text not null,
  organization_id text not null,
  variant_id text not null,
  production_unit_id text,
  serial_number text not null,
  vin text,
  battery_serial_number text,
  state text not null check (state in ('in-transit', 'available', 'reserved', 'sold', 'service', 'quarantine', 'retired')),
  version bigint not null check (version > 0),
  received_at timestamptz,
  updated_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, stock_unit_id),
  foreign key (tenant_id, organization_id) references org.organization (tenant_id, organization_id),
  foreign key (tenant_id, variant_id) references catalog.vehicle_variant (tenant_id, variant_id),
  foreign key (tenant_id, production_unit_id) references factory.production_unit (tenant_id, production_unit_id),
  unique (tenant_id, serial_number),
  unique nulls not distinct (tenant_id, vin),
  unique nulls not distinct (tenant_id, battery_serial_number),
  check ((state = 'in-transit') or received_at is not null)
);

create table pricing.price_book (
  tenant_id uuid not null,
  price_book_id text not null,
  market text not null check (market ~ '^[A-Z]{2}$'),
  currency text not null check (currency ~ '^[A-Z]{3}$'),
  valid_from timestamptz not null,
  valid_until timestamptz,
  status text not null check (status in ('draft', 'active', 'retired')),
  primary key (tenant_id, price_book_id),
  foreign key (tenant_id) references platform.tenant (tenant_id),
  check (valid_until is null or valid_until > valid_from)
);

create table pricing.price_book_entry (
  tenant_id uuid not null,
  price_book_id text not null,
  variant_id text not null,
  amount_minor_units bigint not null check (amount_minor_units >= 0),
  tax_mode text not null check (tax_mode in ('exclusive', 'inclusive', 'not-applicable')),
  primary key (tenant_id, price_book_id, variant_id),
  foreign key (tenant_id, price_book_id) references pricing.price_book (tenant_id, price_book_id),
  foreign key (tenant_id, variant_id) references catalog.vehicle_variant (tenant_id, variant_id)
);

create table sales.customer_order_line (
  tenant_id uuid not null,
  order_id text not null,
  line_id text not null,
  variant_id text not null,
  quantity integer not null check (quantity > 0),
  unit_price_minor_units bigint not null check (unit_price_minor_units >= 0),
  allocated_stock_unit_id text,
  primary key (tenant_id, order_id, line_id),
  foreign key (tenant_id, order_id) references sales.customer_order (tenant_id, order_id),
  foreign key (tenant_id, variant_id) references catalog.vehicle_variant (tenant_id, variant_id),
  foreign key (tenant_id, allocated_stock_unit_id) references inventory.stock_unit (tenant_id, stock_unit_id),
  check (quantity = 1 or allocated_stock_unit_id is null)
);

create table payment.payment_attempt (
  tenant_id uuid not null,
  payment_attempt_id text not null,
  order_id text not null,
  provider_code text not null,
  provider_reference text,
  idempotency_key text not null,
  state text not null check (state in ('created', 'pending', 'authorized', 'captured', 'failed', 'refunded', 'disputed')),
  currency text not null check (currency ~ '^[A-Z]{3}$'),
  amount_minor_units bigint not null check (amount_minor_units > 0),
  version bigint not null check (version > 0),
  updated_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, payment_attempt_id),
  foreign key (tenant_id, order_id) references sales.customer_order (tenant_id, order_id),
  unique (tenant_id, provider_code, idempotency_key),
  unique nulls not distinct (tenant_id, provider_code, provider_reference)
);

create table logistics.shipment (
  tenant_id uuid not null,
  shipment_id text not null,
  provider_code text not null,
  provider_reference text,
  origin_organization_id text not null,
  destination_organization_id text not null,
  state text not null check (state in ('planned', 'booked', 'picked-up', 'in-transit', 'delivered', 'exception', 'cancelled')),
  last_event_at timestamptz,
  primary key (tenant_id, shipment_id),
  foreign key (tenant_id, origin_organization_id) references org.organization (tenant_id, organization_id),
  foreign key (tenant_id, destination_organization_id) references org.organization (tenant_id, organization_id),
  unique nulls not distinct (tenant_id, provider_code, provider_reference)
);

create table service_ops.warranty (
  tenant_id uuid not null,
  warranty_id text not null,
  stock_unit_id text not null,
  customer_principal_id text not null,
  starts_at timestamptz not null,
  ends_at timestamptz not null,
  terms_version text not null,
  status text not null check (status in ('active', 'expired', 'void')),
  primary key (tenant_id, warranty_id),
  foreign key (tenant_id, stock_unit_id) references inventory.stock_unit (tenant_id, stock_unit_id),
  foreign key (tenant_id, customer_principal_id) references crm.customer_profile (tenant_id, customer_principal_id),
  unique (tenant_id, stock_unit_id),
  check (ends_at > starts_at)
);

create table service_ops.service_case (
  tenant_id uuid not null,
  service_case_id text not null,
  stock_unit_id text not null,
  organization_id text not null,
  state text not null check (state in ('opened', 'diagnosis', 'awaiting-parts', 'repair', 'quality', 'closed', 'cancelled')),
  severity text not null check (severity in ('low', 'medium', 'high', 'safety')),
  description text not null,
  version bigint not null check (version > 0),
  opened_at timestamptz not null default clock_timestamp(),
  closed_at timestamptz,
  primary key (tenant_id, service_case_id),
  foreign key (tenant_id, stock_unit_id) references inventory.stock_unit (tenant_id, stock_unit_id),
  foreign key (tenant_id, organization_id) references org.organization (tenant_id, organization_id),
  check ((state = 'closed') = (closed_at is not null))
);

create table service_ops.recall (
  tenant_id uuid not null,
  recall_id text not null,
  recall_code text not null,
  title text not null,
  severity text not null check (severity in ('service', 'safety', 'regulatory')),
  status text not null check (status in ('draft', 'active', 'closed')),
  published_at timestamptz,
  primary key (tenant_id, recall_id),
  foreign key (tenant_id) references platform.tenant (tenant_id),
  unique (tenant_id, recall_code),
  check ((status = 'draft') = (published_at is null))
);

create table service_ops.recall_unit (
  tenant_id uuid not null,
  recall_id text not null,
  stock_unit_id text not null,
  state text not null check (state in ('identified', 'notified', 'scheduled', 'remediated', 'exempt')),
  updated_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, recall_id, stock_unit_id),
  foreign key (tenant_id, recall_id) references service_ops.recall (tenant_id, recall_id),
  foreign key (tenant_id, stock_unit_id) references inventory.stock_unit (tenant_id, stock_unit_id)
);

create table franchise.agreement (
  tenant_id uuid not null,
  agreement_id text not null,
  franchise_organization_id text not null,
  territory_code text not null,
  terms_version text not null,
  starts_on date not null,
  ends_on date,
  status text not null check (status in ('draft', 'active', 'suspended', 'terminated', 'expired')),
  primary key (tenant_id, agreement_id),
  foreign key (tenant_id, franchise_organization_id) references org.organization (tenant_id, organization_id),
  unique (tenant_id, franchise_organization_id, territory_code, starts_on),
  check (ends_on is null or ends_on > starts_on)
);

create table integration.external_mapping (
  tenant_id uuid not null,
  provider_code text not null,
  resource_type text not null,
  internal_id text not null,
  external_id text not null,
  version_token text,
  updated_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, provider_code, resource_type, internal_id),
  foreign key (tenant_id) references platform.tenant (tenant_id),
  unique (tenant_id, provider_code, resource_type, external_id)
);

create table integration.reconciliation_item (
  tenant_id uuid not null,
  reconciliation_id text not null,
  provider_code text not null,
  resource_type text not null,
  internal_id text,
  external_id text,
  state text not null check (state in ('matched', 'missing-internal', 'missing-external', 'amount-mismatch', 'state-mismatch', 'resolved')),
  evidence jsonb not null,
  observed_at timestamptz not null,
  resolved_at timestamptz,
  primary key (tenant_id, reconciliation_id),
  foreign key (tenant_id) references platform.tenant (tenant_id),
  check (jsonb_typeof(evidence) = 'object'),
  check ((state = 'resolved') = (resolved_at is not null))
);

create table marketing.touchpoint (
  tenant_id uuid not null,
  touchpoint_id text not null,
  lead_id text,
  channel text not null,
  campaign_external_id text,
  consent_basis text not null,
  occurred_at timestamptz not null,
  attributes jsonb not null default '{}'::jsonb,
  primary key (tenant_id, touchpoint_id),
  foreign key (tenant_id, lead_id) references crm.lead (tenant_id, lead_id),
  check (jsonb_typeof(attributes) = 'object')
);

create table communication.message (
  tenant_id uuid not null,
  message_id text not null,
  recipient_principal_id text not null,
  channel text not null check (channel in ('email', 'sms', 'push', 'chat')),
  template_code text not null,
  template_version text not null,
  state text not null check (state in ('queued', 'sent', 'delivered', 'failed', 'suppressed')),
  provider_reference text,
  last_error_code text,
  created_at timestamptz not null default clock_timestamp(),
  updated_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, message_id),
  foreign key (tenant_id) references platform.tenant (tenant_id),
  check (updated_at >= created_at)
);

create index lead_org_state_idx on crm.lead (tenant_id, organization_id, lifecycle_state, created_at desc);
create index production_po_state_idx on factory.production_unit (tenant_id, purchase_order_id, state);
create index stock_org_state_variant_idx on inventory.stock_unit (tenant_id, organization_id, state, variant_id);
create index payment_order_state_idx on payment.payment_attempt (tenant_id, order_id, state);
create index service_org_state_idx on service_ops.service_case (tenant_id, organization_id, state, opened_at);
create index reconciliation_open_idx on integration.reconciliation_item (tenant_id, provider_code, observed_at) where state <> 'resolved';
create index message_ready_idx on communication.message (tenant_id, created_at, message_id) where state = 'queued';

commit;
