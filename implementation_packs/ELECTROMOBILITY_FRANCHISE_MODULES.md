# Electromobility Franchise Modules

## 1. Metadata

```yaml
pack_id: "ELECTROMOBILITY-FRANCHISE-MODULES"
pack_version: "0.1.2"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa el modelo transaccional multi-tenant para una red de franquicias de vehículos eléctricos: catálogo, CRM, procurement, fábrica, stock serializado, pricing, venta, pago, logística, garantía, service, recall, integraciones, marketing y comunicación."
stacks: ["PostgreSQL 18.6"]
compatible_with: ["PG-TX-FOUNDATION 0.1.x", "GO-ENTERPRISE-BACKEND >=0.3.0 <1.0.0"]
incompatible_with: ["instalación sin schemas org y sales", "contabilidad/fiscalidad no definida por jurisdicción"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: []
verified_at: "2026-08-25"
```

## 2. Applicability

Este pack expresa sources of truth e invariantes técnicas comunes. No implementa impuestos, facturación legal, homologación, contratos de financiación, normas de baterías, garantías legales ni APIs de proveedores. Esas extensiones se activan por mercado mediante código, tests y asesoramiento competente.

Se adopta para una red multi-tenant que necesita el modelo transaccional completo de catálogo a posventa. Se rechaza si faltan los schemas foundation, si otro componente ya es source of truth de estas entidades o si el proyecto pretende tratar reglas jurisdiccionales como defaults universales.

## 3. Architecture contract

PostgreSQL es el source of truth. Cada agregado conserva `tenant_id`; los límites organizacionales se completan en las APIs autorizadas. Catálogo, CRM, procurement, fábrica, stock serializado, pricing, venta, pago, logística, servicio, franquicia e integraciones viven en schemas separados y se relacionan mediante claves explícitas. Identidades serializadas y eventos externos son únicos; cambios de estado inválidos fallan en transacción. La migración `down` sólo se usa en entornos seguros sin datos que deban preservarse; producción evoluciona expand/migrate/contract.

## 4. Exact file manifest

```text
CREATE db/migrations/0003_electromobility_franchise_modules.up.sql
CREATE db/migrations/0003_electromobility_franchise_modules.down.sql
CREATE db/tests/0003_electromobility_franchise_modules.test.sql
```

## 5. Materialization blocks

### FILE: `db/migrations/0003_electromobility_franchise_modules.up.sql`

```yaml
block_id: "ELECTROMOBILITY-FRANCHISE-MODULES:migration-up:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "208f268916400f4f6c8be3c86be2d923c5065a51f5402a115b35d43868550225"
variables: []
secrets_allowed: false
```

````sql
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
````

### FILE: `db/migrations/0003_electromobility_franchise_modules.down.sql`

```yaml
block_id: "ELECTROMOBILITY-FRANCHISE-MODULES:migration-down:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "a1a0afcb554726f48920777c90fdb613ce58f242d10c9bf93c146af8c64ff6ab"
variables: []
secrets_allowed: false
```

````sql
begin;
drop table sales.customer_order_line;
drop schema communication cascade;
drop schema marketing cascade;
drop schema integration cascade;
drop schema franchise cascade;
drop schema service_ops cascade;
drop schema logistics cascade;
drop schema payment cascade;
drop schema pricing cascade;
drop schema inventory cascade;
drop schema factory cascade;
drop schema procurement cascade;
drop schema partner cascade;
drop schema crm cascade;
drop schema catalog cascade;
commit;
````

### FILE: `db/tests/0003_electromobility_franchise_modules.test.sql`

```yaml
block_id: "ELECTROMOBILITY-FRANCHISE-MODULES:migration-test:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "651e30a1ae88f4b6b59b39c34b431fbea7d5d1c154d3568b5bc6ab88700d8b37"
variables: []
secrets_allowed: false
```

````sql
begin;

insert into platform.tenant (tenant_id, tenant_code, legal_name, display_name)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c28c01', 'mobility-test', 'Mobility Test S.A.', 'Mobility Test');
insert into org.organization (tenant_id, organization_id, organization_code, display_name, organization_type)
values
 ('018f4d4a-7b36-7a21-8d10-2f4c54c28c01', 'hq', 'hq', 'Headquarters', 'franchisor'),
 ('018f4d4a-7b36-7a21-8d10-2f4c54c28c01', 'franchise-1', 'franchise-1', 'Franchise 1', 'franchisee');
insert into catalog.vehicle_model values (
 '018f4d4a-7b36-7a21-8d10-2f4c54c28c01','model-1','moto-one','Moto One','motorcycle','active',true,
 '{"rangeKm":120}',clock_timestamp(),clock_timestamp()
);
insert into catalog.vehicle_variant values (
 '018f4d4a-7b36-7a21-8d10-2f4c54c28c01','variant-1','model-1','moto-one-lfp','Moto One LFP',
 '{"chemistry":"LFP","capacityWh":4000}','approved','active'
);
insert into crm.customer_profile (tenant_id,customer_principal_id,display_name,email_normalized)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c28c01','customer-1','Customer One','customer@example.com');
insert into crm.lead (tenant_id,lead_id,organization_id,customer_principal_id,model_id,lifecycle_state,source_code,contact_payload)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c28c01','lead-1','franchise-1','customer-1','model-1','qualified','google-ads','{"email":"customer@example.com"}');
insert into crm.consent_evidence values (
 '018f4d4a-7b36-7a21-8d10-2f4c54c28c01','consent-1','lead-1','customer-1','sales-contact','v1','granted',clock_timestamp(),repeat('a',64)
);
insert into partner.supplier (tenant_id,supplier_id,supplier_code,legal_name,status)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c28c01','supplier-1','factory-one','Factory One','active');
insert into procurement.purchase_order (tenant_id,purchase_order_id,supplier_id,destination_organization_id,state,currency,total_minor_units,version)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c28c01','po-1','supplier-1','hq','received','USD',100000,1);
insert into factory.production_unit (tenant_id,production_unit_id,purchase_order_id,variant_id,serial_number,vin,battery_serial_number,state)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c28c01','production-1','po-1','variant-1','SERIAL-1','VIN-1','BATTERY-1','received');
insert into inventory.stock_unit (tenant_id,stock_unit_id,organization_id,variant_id,production_unit_id,serial_number,vin,battery_serial_number,state,version,received_at)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c28c01','stock-1','franchise-1','variant-1','production-1','SERIAL-1','VIN-1','BATTERY-1','available',1,clock_timestamp());
insert into sales.customer_order (tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c28c01','order-1','franchise-1','customer-1','allocated','USD',120000,1);
insert into sales.customer_order_line values (
 '018f4d4a-7b36-7a21-8d10-2f4c54c28c01','order-1','line-1','variant-1',1,120000,'stock-1'
);
insert into payment.payment_attempt (tenant_id,payment_attempt_id,order_id,provider_code,provider_reference,idempotency_key,state,currency,amount_minor_units,version)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c28c01','payment-1','order-1','sandbox','provider-payment-1','idem-payment-1','captured','USD',120000,1);
insert into service_ops.warranty values (
 '018f4d4a-7b36-7a21-8d10-2f4c54c28c01','warranty-1','stock-1','customer-1',clock_timestamp(),clock_timestamp()+interval '2 years','terms-v1','active'
);
insert into franchise.agreement (tenant_id,agreement_id,franchise_organization_id,territory_code,terms_version,starts_on,status)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c28c01','agreement-1','franchise-1','AR-CBA','v1',current_date,'active');

do $test$
begin
  begin
    insert into inventory.stock_unit (tenant_id,stock_unit_id,organization_id,variant_id,serial_number,state,version,received_at)
    values ('018f4d4a-7b36-7a21-8d10-2f4c54c28c01','stock-duplicate','franchise-1','variant-1','SERIAL-1','available',1,clock_timestamp());
    raise exception 'expected duplicate serial rejection';
  exception when unique_violation then null;
  end;
end;
$test$;

do $test$
declare
  model_count bigint;
  paid_total bigint;
begin
  select count(*) into model_count from catalog.vehicle_model where tenant_id='018f4d4a-7b36-7a21-8d10-2f4c54c28c01' and publicly_visible;
  select sum(amount_minor_units) into paid_total from payment.payment_attempt where tenant_id='018f4d4a-7b36-7a21-8d10-2f4c54c28c01' and state='captured';
  if model_count <> 1 or paid_total <> 120000 then
    raise exception 'vertical fixture mismatch';
  end if;
end;
$test$;

rollback;
````

## 6. Configuration surface

Este pack no lee variables runtime: materializa schema SQL determinista. La configuración de módulos, mercado, permisos y proveedores pertenece al blueprint y a las capas API/worker; una capability desactivada no debe recibir rutas ni jobs. Nombres de schema, unidades monetarias menores e IDs no son configurables sin una migración/versionado explícito.

## 7. Dependency bill

| Tool/runtime | Pin | Use | License | Scope | Official source |
|---|---|---|---|---|---|
| PostgreSQL | `18.6` verified baseline | transactional schema and tests | PostgreSQL | runtime/test | `postgresql.org` |

No hay extensiones de terceros ni código upstream embebido.

## 8. Apply order

Aplicar después de `PG-TX-FOUNDATION` y su migration/test; ejecutar `0003 ... up`, luego su fixture dentro de una base descartable. Componer APIs/workers sólo después. En un workspace existente, auditar colisiones de schemas/tablas y ownership antes de ejecutar. Rollback de desarrollo usa `down`; producción requiere backup probado y plan de migración forward-compatible.

## 9. Verification

Ejecutar después de migrations 0001 y 0002. Estado `REBUILD_VERIFIED / CONDITIONED`, evidencia `ELECTROMOBILITY-MODULES-20260824-V1`: 23 tablas, 14 schemas y 7 índices se ejecutaron en PostgreSQL 18.6; fixture vertical, duplicado serial y `up-test-down-up-test` pasaron.

Condiciones: APIs/UI/workers, tenant negatives adicionales, asignación de stock concurrente, crashes order/payment/reconciliation, migration locks, datos representativos, roles mínimos, reglas fiscales/homologación/garantía por jurisdicción y restore. No es production-ready ni `REUSABLE_PACK`.

## 10. Reconstruction evidence

La reconstrucción limpia, hashes, migration up/down/up y fixtures sobre PostgreSQL 18.6 están registrados en `reconstruction_evidence/ELECTROMOBILITY_FRANCHISE_MODULES_2026-08-24_V1.md`; la auditoría final de biblioteca vuelve a validar el pack vigente.
