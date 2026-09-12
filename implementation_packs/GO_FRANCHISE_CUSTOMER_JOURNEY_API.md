# Go Franchise Customer Journey API

V312: corrección AUTHORED de recepción y decisión de devoluciones.
Consulta exacta por autorización, scope de grafo bloqueado en escritura,
lectura atómica, actor y evidencia inmutables. UI conserva huella y consulta
sin reenviar aun con lista inaccesible. Refund4/exchange3 solicitudes.
Evidencia reconstruction_evidence/RETURN_OPERATIONS_RECOVERY_V312.md.
Sin efectos downstream, cambio comercial ni promoción integral.

V311: corrección AUTHORED de recuperación de checklist completada por
identidad natural de entrega, actor y respuestas inmutables. Consulta scoped
preserva estado vigente accepted/rejected/presented; no reenvía ni reabre.
Evidencia reconstruction_evidence/CHECKLIST_COMPLETION_RECOVERY_V311.md.
Sin cambio de reglas comerciales, upstreams ni promoción integral.

V310: corrección AUTHORED de publicación de checklist con consulta exacta
por identidad natural org/id/versión, actor atómico y recuperación sin POST.
Evidencia reconstruction_evidence/CHECKLIST_PUBLICATION_RECOVERY_V310.md.
No altera inmutabilidad, reglas de entrega, upstreams ni promoción integral.

V309: corrección AUTHORED de publicación/consulta de turnos con identidad
retenida, Idempotency-Key y recuperación scoped sin reenvío. Pruebas y límites
en reconstruction_evidence/SLOT_CREATION_RECOVERY_V309.md. No cambio de reglas de agenda,
upstreams ni promoción integral; FAIL457 conserva los otros comandos pendientes.

V308: corrección AUTHORED de alta/consulta de recursos con identidad
retenida, Idempotency-Key y recuperación scoped sin reenvío. Pruebas y límites
en reconstruction_evidence/RESOURCE_CREATION_RECOVERY_V308.md. No cambio de reglas de agenda,
upstreams ni promoción integral; FAIL457 conserva los otros comandos pendientes.

V307: corrección AUTHORED de alta/consulta de disponibilidad con identidad
retenida, Idempotency-Key y recuperación scoped sin reenvío. Pruebas y límites
en reconstruction_evidence/AVAILABILITY_CREATION_RECOVERY_V307.md. No cambio de reglas de agenda,
upstreams ni promoción integral; FAIL457 conserva los otros comandos pendientes.

V306: corrección AUTHORED de emisión/consulta de cotización con identidad
retenida, Idempotency-Key y recuperación scoped sin reenvío. Pruebas y límites
en reconstruction_evidence/QUOTE_CREATION_RECOVERY_V306.md. No cambio de pricing,
upstreams ni promoción integral; FAIL457 conserva los otros comandos pendientes.

V303 / 0.10.8 amplía únicamente el harness: ELITE_RESOLUTION_RECOVERY_E2E=1
requiere la entrega rechazada real del fixture ELITE_DELIVERY_ACTION_E2E.
Tres decisiones existentes por UI/BFF/HTTP/PG; pérdida de respuesta, JSON
incompleto, consulta503, storage ausente, replay y carrera sin efectos duplicados.
Los dos preparados adicionales son fixtures, no el creador inicial de producto.
Evidencia: reconstruction_evidence/DELIVERY_EXCEPTION_RECOVERY_V303.md.

V302 / 0.10.7 cambia sólo el harness HTTP/browser existente: lane opt-in de
secciones por rol, fallo503 de lectura y recuperación, JWT verificado y cero
efecto para escritura no autorizada. No cambia dominio, SQL de producto ni
política comercial. Ver reconstruction_evidence/OPERATOR_SECTIONS_RECOVERY_V302.md.

V300 valida también las lecturas de actas/excepciones; CustomerJourney usa un
snapshot read-only común y no devuelve resultados parciales al detectar vínculos
incoherentes. Items sólo para actas seleccionadas. Cambios AUTHORED; evidencia
en reconstruction_evidence/DELIVERY_READ_RECOVERY_V300.md. No habilita entrega inicial.

V299 endurece las relaciones de entrega existentes: checklist, aceptación,
rechazo y resolución contrastan organización/pedido/cliente/stock dentro de la
misma transacción. Cambio AUTHORED, no código Microsoft. No crea la preparación
inicial ni prueba pago o liberación financiera. Evidencia y límites en
reconstruction_evidence/DELIVERY_RELATIONAL_SCOPE_V299.md.

V294 añade el gate opt-in TestQuoteAcceptanceBrowserPostgres en el owner de
tests existente. Exige base descartable loopback elite_confirmation_* con todas
las migraciones, ELITE_QUOTE_E2E=1 y ELITE_WEB_ROOT absoluto al web construido.
Conecta portal/sesión RS256/JWKS/BFF/Go/PostgreSQL, respuesta perdida después de
commit, dos aceptaciones concurrentes, permisos y recuperación. No cambia el
dominio: el segundo POST sigue 409 y el cliente recupera el pedido mediante GET.
ELITE_QUOTE_PROJECT permite diagnóstico de un solo navegador, no cierre global.
Fixtures/harness son AUTHORED, no implementación copiada de Microsoft.
Evidencia y límites: reconstruction_evidence/QUOTE_ACCEPTANCE_CONNECTED_V294.md.

Continuación V294: después de la recuperación, los mismos pedidos de navegador
compiten por una unidad serial mediante Commerce existente; sólo uno reserva.
Una solicitud de pago created queda durable y rechaza clave duplicada, moneda,
importe u organización incorrectos. No dispara worker ni proveedor; no prueba
UI operativa, autorización HTTP de Commerce, cobro, confirmación ni entrega.

## 1. Metadata

```yaml
pack_id: "GO-FRANCHISE-CUSTOMER-JOURNEY-API"
pack_version: "0.10.21"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Extiende el backend canónico sin duplicar dominios: ubicación, agenda/recursos, leads, cotización idempotente y quote-to-order, y entrega con checklist, excepción, autorización, recepción/inspección física y efectos separados."
stacks: ["Go 1.26.7", "PostgreSQL 18.6", "pgx 5.10.0"]
compatible_with: ["GO-ELECTROMOBILITY-PUBLIC-CRM-API 0.2.3 (shared browser harness)", "GO-ENTERPRISE-BACKEND 0.4.x", "ELECTROMOBILITY-FRANCHISE-MODULES 0.1.x", "GO-ENTERPRISE-QUERY-API 0.1.x", "GO-ELECTROMOBILITY-APPLICATION 1.5.x", "GO-BC-SALES-CONTRACT-ADAPTER 0.1.0"]
incompatible_with: ["alternate owner for CRM, quotation, appointment or delivery handover without synchronization contract"]
license_expression: "LicenseRef-Workspace-Owner AND MIT dependency"
upstream_sources: ["https://github.com/microsoft/BCApps/blob/31a860b527f0dc72c7a44a255d7e7d403cfa4789/src/Layers/W1/DemoTool/CreateResourceCapacityEntry.Codeunit.al", "https://github.com/microsoft/BCApps/blob/31a860b527f0dc72c7a44a255d7e7d403cfa4789/src/Layers/W1/Tests/Resource/ResourceMatrixManagement.Codeunit.al", "https://github.com/microsoft/BCApps/blob/31a860b527f0dc72c7a44a255d7e7d403cfa4789/src/Layers/W1/BaseApp/Projects/Resources/Resource/Resource.Table.al", "https://github.com/microsoft/BCApps/blob/31a860b527f0dc72c7a44a255d7e7d403cfa4789/src/Layers/W1/BaseApp/Projects/Resources/Resource/ResCapacityEntry.Table.al", "https://github.com/microsoft/BCApps/blob/31a860b527f0dc72c7a44a255d7e7d403cfa4789/src/Layers/W1/BaseApp/Service/Resources/ResourceSkill.Table.al", "https://github.com/microsoft/BCApps/blob/31a860b527f0dc72c7a44a255d7e7d403cfa4789/src/Layers/W1/Tests/SCM-Service/ServiceResourceSkill.Codeunit.al", "https://learn.microsoft.com/en-us/dynamics365/business-central/hr-how-manage-absence", "https://learn.microsoft.com/en-us/dynamics365/business-central/production-how-to-create-work-center-calendars", "https://learn.microsoft.com/en-gb/dynamics365/business-central/across-how-to-assign-base-calendars", "https://learn.microsoft.com/en-us/dynamics365/business-central/application/base-application/table/microsoft.humanresources.absence.employee-absence", "https://learn.microsoft.com/en-us/dynamics365/business-central/ui-post-sales", "https://learn.microsoft.com/en-us/dynamics365/field-service/inspections-overview", "https://www.postgresql.org/docs/18/explicit-locking.html", "https://github.com/GoogleCloudPlatform/microservices-demo/tree/5b3a712ab85ccb8f6f7cd5b720d36ba9a8d041eb", "https://sre.google/sre-book/reliable-product-launches/", "https://sre.google/sre-book/release-engineering/", "https://www.nasa.gov/reference/system-engineering-handbook-appendix/"]
verified_at: "2026-09-08"
```

Todo bloque es `AUTHORED`: Microsoft BCApps y Google Online Boutique gobiernan límites, flujos y contratos de referencia, pero ninguna línea local se atribuye a Microsoft o Google. El precio se resuelve desde `pricing.price_book_entry`; la UI y el caller no son autoridad monetaria.

V402 staging candidate / 0.10.20: delegates the precisely mapped price/quote transformations to GO-BC-SALES-CONTRACT-ADAPTER0.1.0. Existing caller provenance stays AUTHORED; scheduling, pricebook exclusivity, consent and financial state policy are not reclassified. No canonical publication is asserted by this candidate.

## 2. Applicability

V298 amplía únicamente el harness conectado existente: ELITE_PAYMENT_E2E=1 añade
request/desactivado/reinicio sobre pedidos creados en el portal y owner Commerce.
RS256/JWKS, dos intents/outbox con actor, respuesta perdida y carrera8 claves.
La selección stripe se inyecta sólo en fixture sintético, nunca invoca proveedor.
Budget finito5 minutos por navegador para las nueve fases; código del dominio
FranchiseJourney no se cambia. AUTHORED; evidencia PAYMENT_REQUEST_PORTAL_V298.md.

V261 agrega verificación de solicitud pública de turno mediante el harness compartido de `GO-ELECTROMOBILITY-PUBLIC-CRM-API 0.2.3` y `MICROSOFT-PLAYWRIGHT-BROWSER-GATE 0.1.4`. Go propaga `Idempotency-Replayed` sólo para el replay demostrado. El gate atraviesa catálogo/lead/BFF/API/PG, pierde una respuesta después del commit y recupera el mismo turno sin duplicación; prueba capacidad agotada, divergencia y payload inválido. El estado demostrado es `requested`, no confirmación por operador. Tests/glue `AUTHORED`, sin migraciones o dependencias nuevas. Evidencia: `reconstruction_evidence/PUBLIC_APPOINTMENT_RECOVERY_V261.md`.

Use when the existing Go/PostgreSQL franchise composition needs a single visitor-to-customer operational slice without a second CRM, catalog, order or inventory owner. Reject it when those records belong to an external ERP/CRM unless bidirectional ownership and reconciliation are designed. It assumes migrations 0001-0004, verified OIDC principals and the existing outbox.

## 3. Architecture contract

The pack extends existing bounded contexts instead of copying them. Public locations are explicit publications. Appointment creation resolves tenant, organization, lead and customer inside one transaction and uses durable request-hash idempotency. Lead assignment and lifecycle changes are tenant/organization scoped, versioned and protected both by repository predicates and a database transition trigger. Quote creation now uses the same durable platform idempotency owner: the tenant-scoped key and exact body hash serialize concurrent retries, same-key/same-hash returns the original quote, same-key/different-hash conflicts, and quote plus one outbox event commit atomically. Quotes use active server-side price-book values. Quote acceptance locks an issued, unexpired, customer-owned quote and creates the existing canonical order, its line, immutable acceptance evidence and two outbox events in one transaction; price and currency never come from the browser. Customer reads and handover acceptance bind ownership to the verified token subject. Cross-scope, stale-version, unknown-price, expired-quote and invalid-transition paths fail closed. Rollback removes only this extension and preserves earlier schemas.

Capacity is an organization/kind calendar owned by PostgreSQL. Explicit absolute working windows and unavailable intervals model calculated availability without inventing recurrence, timezone, employment or statutory rules. An open slot must be fully covered by organization working time and must not overlap organization unavailability. Resource assignment additionally requires a covering resource working window, no resource absence, the correct skill and no overlapping active appointment. Advisory transaction locks serialize competing availability, assignment and capacity decisions. Cancelling a working window that supports an active appointment fails closed.

Every staff appointment transition binds the actor to the verified token subject. Cancellation and no-show require a governed reason code; unrelated states reject one. The append-only transition ledger is immutable. Customer cancellation is restricted to the verified customer, the organization scope, a future requested/confirmed appointment and the expected version; state, transition evidence and outbox event commit together. Delivery acceptance also binds the verified customer and version, requires explicit receipt confirmation, compares the observed serial with the durable stock unit and stores a server-calculated SHA-256 of the exact accepted command; browser-authored evidence digests are rejected. This is an authenticated audit record, not a qualified electronic signature. Confirmation requires an assignment; completion/no-show cannot precede the appointment. Public queries return only published, future, non-full slots over at most 31 days. Performance budget: bounded pages are at most 100 records; public slot results are capped at 500; availability queries cover at most 366 days and return at most 1000 rows; customer timelines are capped at 100 rows per collection; HTTP bodies are capped at 1 MiB. Project-specific employment/privacy rules, recurrence/timezone calendars, retention, legal acceptance/signature, quote terms, taxes and financing remain explicit conditions.

Delivery checklist definitions are explicit organization-scoped versions. Publication and all items commit atomically; a published definition cannot be edited. Staff responses are append-only, bind the verified operator and exact handover/version, and the database refuses presentation while any required item is unanswered. Customer acceptance then requires that same completed checklist ID/version in addition to the durable asset serial and expected handover version. Microsoft Field Service inspections govern the published/versioned required-response model and Business Central posting governs the durable operational boundary; all local implementation remains declared `AUTHORED`.

A customer can reject only the currently presented handover in its own organization and identity scope. Rejection, exact-command evidence, durable exception and outbox event commit atomically. A verified operator resolves an open exception exactly once. Correction never rewrites the rejected handover: it creates one prepared successor linked by `supersedes_handover_id`, which must complete a checklist and be presented again. Return/exchange creates an immutable authorization linked to the original order and stock unit with `exact_cost_source_order_id`; it does not mutate inventory, issue money or post accounting. Microsoft Business Central sales-return documentation and public BCApps return-order/receipt code govern this separation and exact-cost link.

Physical receipt binds the authorization, organization, original order, durable stock unit, customer and exact serial; condition, notes, evidence and verified receiver are append-only. A single disposition derives refund versus exchange from the authorization rather than accepting that authority from the browser. It creates immutable, idempotent requests for the existing inventory, payment or fulfillment, accounting and fiscal owners and one outbox event atomically. A request is not a completed effect: each selected project must connect, reconcile and prove those owners before production.

## 4. Exact file manifest

V263 adds `GET /v1/franchise/agenda` under `appointment:manage` and verified tenant/organization scope. Explicit RFC3339 `from` inclusive / `to` exclusive, maximum 31 days; read-only repeatable-read snapshot, maximum 200 appointments and 200 active resources, `truncated=true` on overflow. The read model omits resource principal subjects, does not determine guaranteed availability and never replaces assignment/transition constraints. The web gates select a UTC day, block actions on incomplete results, remove manual appointment IDs/versions and recover ambiguous writes by reading durable state. APP_BASE_URL is the exact trusted public origin, never inferred from forwarded headers. The opt-in browser harness uses real loopback TLS with a synthetic certificate and JWE/RS256 fixtures, not a live IdP login. New read model/tests remain AUTHORED under the existing owner. Full evidence: `reconstruction_evidence/OPERATOR_AGENDA_BROWSER_POSTGRES_V263.md`.

V262 adds an opt-in confirmation integration gate in the existing HTTP test file: materialize this backend plus TS-OIDC-PORTAL-ADAPTER 0.2.2, install the exact web lock, apply all migrations to a disposable loopback database named `elite_confirmation_*`, and run `go test ./internal/platform/httpapi -run '^TestAppointmentConfirmationBFFPostgres$' -v -count=1 -timeout=3m` with `ELITE_CONFIRMATION_E2E=1`, `TEST_DATABASE_URL` and absolute `ELITE_WEB_ROOT`. The existing BFF client calls the real Go API with synthetic RS256 tokens validated through local OIDC discovery/JWKS. Two confirmations produce one committed state, immutable actor audit and outbox event; the customer read is ownership-scoped. AUTHORED fixture, not official copied code, live IdP login, browser E2E or delivery notification evidence. Audit fixtures deliberately remain until the caller discards that dedicated test database; never disable immutable-audit triggers to clean them. Evidence and official method references: `reconstruction_evidence/APPOINTMENT_CONFIRMATION_BFF_POSTGRES_V262.md`.

```text
CREATE db/migrations/0005_franchise_customer_journey.up.sql
CREATE db/migrations/0005_franchise_customer_journey.down.sql
CREATE db/tests/0005_franchise_customer_journey.test.sql
CREATE db/migrations/0006_quote_order_conversion.up.sql
CREATE db/migrations/0006_quote_order_conversion.down.sql
CREATE db/tests/0006_quote_order_conversion.test.sql
CREATE db/migrations/0007_appointment_capacity.up.sql
CREATE db/migrations/0007_appointment_capacity.down.sql
CREATE db/tests/0007_appointment_capacity.test.sql
CREATE internal/franchisejourney/service.go
CREATE internal/franchisejourney/service_test.go
CREATE internal/platform/postgres/franchisejourney.go
CREATE internal/platform/postgres/franchisejourney_integration_test.go
CREATE internal/platform/httpapi/franchisejourney.go
CREATE internal/platform/httpapi/franchisejourney_test.go
CREATE db/migrations/0009_appointment_resources.up.sql
CREATE db/migrations/0009_appointment_resources.down.sql
CREATE db/tests/0009_appointment_resources.test.sql
CREATE internal/franchisejourney/availability.go
CREATE internal/franchisejourney/availability_test.go
CREATE internal/platform/postgres/franchisejourney_availability.go
CREATE db/migrations/0015_franchise_availability_and_appointment_audit.up.sql
CREATE db/migrations/0015_franchise_availability_and_appointment_audit.down.sql
CREATE db/tests/0015_franchise_availability_and_appointment_audit.test.sql
CREATE db/migrations/0016_versioned_delivery_checklist.up.sql
CREATE db/migrations/0016_versioned_delivery_checklist.down.sql
CREATE db/tests/0016_versioned_delivery_checklist.test.sql
CREATE db/migrations/0017_delivery_exception_and_return_authorization.up.sql
CREATE db/migrations/0017_delivery_exception_and_return_authorization.down.sql
CREATE db/tests/0017_delivery_exception_and_return_authorization.test.sql
CREATE db/migrations/0018_return_receipt_disposition_effects.up.sql
CREATE db/migrations/0018_return_receipt_disposition_effects.down.sql
CREATE db/tests/0018_return_receipt_disposition_effects.test.sql
```

## 5. Materialization blocks

### FILE: `db/migrations/0005_franchise_customer_journey.up.sql`

```yaml
block_id: "GO-FRANCHISE-JOURNEY:file:01"
operation: CREATE
provenance: AUTHORED
source: "local implementation governed by the upstream references in metadata"
license: "LicenseRef-Workspace-Owner"
sha256: "aa750ffe7a56cbdbbbc33dee4aabd0860b18268227e31dee141914b8198bf958"
variables: []
secrets_allowed: false
```

````sql
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
````

### FILE: `db/migrations/0005_franchise_customer_journey.down.sql`

```yaml
block_id: "GO-FRANCHISE-JOURNEY:file:02"
operation: CREATE
provenance: AUTHORED
source: "local implementation governed by the upstream references in metadata"
license: "LicenseRef-Workspace-Owner"
sha256: "11cf0e32119a91607aa0f2033d91c6b342360027fe8b0ecc1606ce813aec1152"
variables: []
secrets_allowed: false
```

````sql
begin;
drop table sales.delivery_handover;
drop table sales.quotation;
drop table crm.appointment;
drop trigger enforce_lead_transition on crm.lead;
drop function crm.enforce_lead_transition();
drop index crm.lead_org_updated_idx;
alter table crm.lead drop column version, drop column assigned_subject;
drop table org.public_location;
commit;
````

### FILE: `db/tests/0005_franchise_customer_journey.test.sql`

```yaml
block_id: "GO-FRANCHISE-JOURNEY:file:03"
operation: CREATE
provenance: AUTHORED
source: "local implementation governed by the upstream references in metadata"
license: "LicenseRef-Workspace-Owner"
sha256: "51877e18ae70930131bee2ae6d9292839d9ad5ad0cc11f4f38a474593eb65032"
variables: []
secrets_allowed: false
```

````sql
begin;

insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c29a01','journey-test','Journey Test','Journey Test');
insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c29a01','store','store','Store','store');
insert into org.public_location(tenant_id,organization_id,city,region,country,published)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c29a01','store','Cordoba','Cordoba','AR',true);
insert into catalog.vehicle_model(tenant_id,model_id,model_code,display_name,vehicle_class,lifecycle_state)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c29a01','model','model','Model','bicycle','active');
insert into catalog.vehicle_variant(tenant_id,variant_id,model_id,variant_code,display_name,battery_specification,lifecycle_state)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c29a01','variant','model','variant','Variant','{}','active');
insert into crm.customer_profile(tenant_id,customer_principal_id,display_name,email_normalized)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c29a01','customer','Customer','customer@example.test');
insert into crm.lead(tenant_id,lead_id,organization_id,customer_principal_id,model_id,lifecycle_state,source_code,contact_payload)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c29a01','lead','store','customer','model','new','public-web','{}');
insert into crm.appointment(tenant_id,appointment_id,organization_id,lead_id,customer_principal_id,model_id,appointment_kind,starts_at,state,version)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c29a01','appointment','store','lead','customer','model','test-drive',clock_timestamp()+interval '1 hour','requested',1);
insert into pricing.price_book(tenant_id,price_book_id,market,currency,valid_from,status)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c29a01','retail','AR','ARS',clock_timestamp()-interval '1 day','active');
insert into pricing.price_book_entry(tenant_id,price_book_id,variant_id,amount_minor_units,tax_mode)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c29a01','retail','variant',100000,'inclusive');
insert into sales.quotation(tenant_id,quotation_id,organization_id,lead_id,customer_principal_id,variant_id,price_book_id,currency,total_minor_units,valid_until,state,version)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c29a01','quote','store','lead','customer','variant','retail','ARS',100000,clock_timestamp()+interval '7 days','issued',1);
insert into inventory.stock_unit(tenant_id,stock_unit_id,organization_id,variant_id,serial_number,state,version,received_at)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c29a01','stock','store','variant','SERIAL','sold',1,clock_timestamp());
insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c29a01','order','store','customer','delivered','ARS',100000,1);
insert into sales.delivery_handover(tenant_id,handover_id,organization_id,order_id,customer_principal_id,stock_unit_id,state,version)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c29a01','handover','store','order','customer','stock','prepared',1);

do $$
begin
  begin
    update crm.lead set lifecycle_state='converted',version=version+1 where tenant_id='018f4d4a-7b36-7a21-8d10-2f4c54c29a01' and lead_id='lead';
    raise exception 'test failed: database alone allowed a skipped lead state';
  exception when check_violation then null;
  end;
end $$;

select 1 / case when (select count(*) from org.public_location where published)=1 then 1 else 0 end;
select 1 / case when (select count(*) from crm.appointment)=1 then 1 else 0 end;
select 1 / case when (select count(*) from sales.quotation)=1 then 1 else 0 end;
select 1 / case when (select count(*) from sales.delivery_handover)=1 then 1 else 0 end;

rollback;
````

### FILE: `db/migrations/0006_quote_order_conversion.up.sql`

```yaml
block_id: "GO-FRANCHISE-JOURNEY:file:10"
operation: CREATE
provenance: AUTHORED
source: "local implementation governed by the upstream references in metadata"
license: "LicenseRef-Workspace-Owner"
sha256: "8f21e222c2f95ff1a677ef0f54d715e91889ec4aa11393c7ef938464fea0fb49"
variables: []
secrets_allowed: false
```

````sql
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
````

### FILE: `db/migrations/0006_quote_order_conversion.down.sql`

```yaml
block_id: "GO-FRANCHISE-JOURNEY:file:11"
operation: CREATE
provenance: AUTHORED
source: "local implementation governed by the upstream references in metadata"
license: "LicenseRef-Workspace-Owner"
sha256: "3255f45a760b1b177ea94998f0cc22c84c4f2ec6bcc32b8f5433d2465995224a"
variables: []
secrets_allowed: false
```

````sql
begin;

drop table if exists sales.quotation_acceptance;
drop index if exists sales.quotation_order_unique_idx;
alter table sales.quotation
  drop constraint if exists quotation_acceptance_order_check,
  drop constraint if exists quotation_order_fk,
  drop column if exists order_id;

commit;
````

### FILE: `db/tests/0006_quote_order_conversion.test.sql`

```yaml
block_id: "GO-FRANCHISE-JOURNEY:file:12"
operation: CREATE
provenance: AUTHORED
source: "local implementation governed by the upstream references in metadata"
license: "LicenseRef-Workspace-Owner"
sha256: "cc3d045843af5daee3e3dbfcdee31adc5e38361505a2852b7e4177f9f5ee27e3"
variables: []
secrets_allowed: false
```

````sql
begin;

do $$
begin
  if not exists (select 1 from information_schema.columns where table_schema='sales' and table_name='quotation' and column_name='order_id') then
    raise exception 'quotation.order_id missing';
  end if;
  if not exists (select 1 from information_schema.tables where table_schema='sales' and table_name='quotation_acceptance') then
    raise exception 'quotation_acceptance missing';
  end if;
  if not exists (select 1 from pg_constraint where conname='quotation_acceptance_order_check') then
    raise exception 'quotation acceptance invariant missing';
  end if;
end $$;

rollback;
````

### FILE: `internal/franchisejourney/service.go`

```yaml
block_id: "GO-FRANCHISE-JOURNEY:file:04"
operation: CREATE
provenance: AUTHORED
source: "local implementation governed by the upstream references in metadata"
license: "LicenseRef-Workspace-Owner"
sha256: "c1d1e2ba2644a2fe26f3483c4387ac13a3143f024c13011187bbccb15102c28e"
variables: []
secrets_allowed: false
```

````go
package franchisejourney

import (
	"context"
	"elite.local/enterprise/internal/businesspolicy"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

var (
	ErrInvalid  = errors.New("invalid franchise journey input")
	ErrConflict = errors.New("franchise journey conflict")
	ErrNotFound = errors.New("franchise journey resource not found")
)

type Location struct {
	OrganizationID string `json:"organization_id"`
	Code           string `json:"code"`
	Name           string `json:"name"`
	City           string `json:"city"`
	Region         string `json:"region"`
	Country        string `json:"country"`
	ContactPhone   string `json:"contact_phone,omitempty"`
	ContactEmail   string `json:"contact_email,omitempty"`
}

type Appointment struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organization_id"`
	LeadID         string    `json:"lead_id"`
	ModelID        string    `json:"model_id,omitempty"`
	Kind           string    `json:"kind"`
	StartsAt       time.Time `json:"starts_at"`
	State          string    `json:"state"`
	Version        int64     `json:"version"`
	SlotID         string    `json:"slot_id,omitempty"`
	EndsAt         time.Time `json:"ends_at,omitempty"`
	ResourceID     string    `json:"resource_id,omitempty"`
}

type ServiceResource struct {
	ID               string   `json:"id"`
	OrganizationID   string   `json:"organization_id"`
	PrincipalSubject string   `json:"principal_subject,omitempty"`
	DisplayName      string   `json:"display_name"`
	Kind             string   `json:"kind"`
	Status           string   `json:"status"`
	Version          int64    `json:"version"`
	Skills           []string `json:"skills"`
}

type AppointmentAgenda struct {
	Appointments []Appointment     `json:"appointments"`
	Resources    []ServiceResource `json:"resources"`
	Truncated    bool              `json:"truncated"`
}

type AppointmentSlot struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organization_id"`
	Kind           string    `json:"kind"`
	StartsAt       time.Time `json:"starts_at"`
	EndsAt         time.Time `json:"ends_at"`
	Capacity       int       `json:"capacity"`
	Booked         int64     `json:"booked"`
	State          string    `json:"state"`
	Version        int64     `json:"version"`
}

type Lead struct {
	ID              string    `json:"id"`
	OrganizationID  string    `json:"organization_id"`
	ModelID         string    `json:"model_id,omitempty"`
	State           string    `json:"state"`
	SourceCode      string    `json:"source_code"`
	AssignedSubject string    `json:"assigned_subject,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	Version         int64     `json:"version"`
}

type Quote struct {
	ID              string    `json:"id"`
	OrganizationID  string    `json:"organization_id"`
	LeadID          string    `json:"lead_id"`
	CustomerSubject string    `json:"customer_subject,omitempty"`
	VariantID       string    `json:"variant_id"`
	Currency        string    `json:"currency"`
	TotalMinorUnits int64     `json:"total_minor_units"`
	ValidUntil      time.Time `json:"valid_until"`
	State           string    `json:"state"`
	Version         int64     `json:"version"`
	PriceBookID     string    `json:"price_book_id"`
	OrderID         string    `json:"order_id,omitempty"`
}

type Handover struct {
	ID                   string          `json:"id"`
	OrganizationID       string          `json:"organization_id"`
	OrderID              string          `json:"order_id"`
	CustomerSubject      string          `json:"customer_subject"`
	StockUnitID          string          `json:"stock_unit_id"`
	State                string          `json:"state"`
	Version              int64           `json:"version"`
	CustomerAcceptedAt   *time.Time      `json:"customer_accepted_at,omitempty"`
	AcceptanceEvidence   string          `json:"acceptance_evidence_sha256,omitempty"`
	ChecklistID          string          `json:"checklist_id,omitempty"`
	ChecklistVersion     int64           `json:"checklist_version,omitempty"`
	ChecklistTitle       string          `json:"checklist_title,omitempty"`
	ChecklistItems       []ChecklistItem `json:"checklist_items"`
	ChecklistCompletedAt *time.Time      `json:"checklist_completed_at,omitempty"`
	SupersedesHandoverID string          `json:"supersedes_handover_id,omitempty"`
}

type ChecklistItem struct {
	ID           string `json:"id"`
	Ordinal      int    `json:"ordinal"`
	Prompt       string `json:"prompt"`
	ResponseType string `json:"response_type"`
	Required     bool   `json:"required"`
}

type DeliveryChecklist struct {
	ID             string          `json:"id"`
	OrganizationID string          `json:"organization_id"`
	Version        int64           `json:"version"`
	Title          string          `json:"title"`
	State          string          `json:"state"`
	Items          []ChecklistItem `json:"items"`
}

type ChecklistResponse struct {
	ItemID         string `json:"item_id"`
	ResponseText   string `json:"response_text"`
	EvidenceSHA256 string `json:"evidence_sha256,omitempty"`
}

type DeliveryException struct {
	ID                    string     `json:"id"`
	OrganizationID        string     `json:"organization_id"`
	HandoverID            string     `json:"handover_id"`
	CustomerSubject       string     `json:"customer_subject"`
	ReasonCode            string     `json:"reason_code"`
	Details               string     `json:"details"`
	State                 string     `json:"state"`
	Version               int64      `json:"version"`
	CreatedAt             time.Time  `json:"created_at"`
	ResolvedAt            *time.Time `json:"resolved_at,omitempty"`
	ResolutionAction      string     `json:"resolution_action,omitempty"`
	SuccessorHandoverID   string     `json:"successor_handover_id,omitempty"`
	ReturnAuthorizationID string     `json:"return_authorization_id,omitempty"`
}

type DeliveryResolution struct {
	Exception             DeliveryException `json:"exception"`
	SuccessorHandover     *Handover         `json:"successor_handover,omitempty"`
	ReturnAuthorizationID string            `json:"return_authorization_id,omitempty"`
	Disposition           string            `json:"disposition,omitempty"`
}

type ReturnReceipt struct {
	ID                   string    `json:"id"`
	AuthorizationID      string    `json:"authorization_id"`
	OrganizationID       string    `json:"organization_id"`
	OrderID              string    `json:"order_id"`
	StockUnitID          string    `json:"stock_unit_id"`
	CustomerSubject      string    `json:"customer_subject"`
	ReceivedSerialNumber string    `json:"received_serial_number"`
	ConditionCode        string    `json:"condition_code"`
	Notes                string    `json:"notes"`
	EvidenceSHA256       string    `json:"evidence_sha256"`
	ReceivedBySubject    string    `json:"received_by_subject"`
	ReceivedAt           time.Time `json:"received_at"`
}

type ReturnEffectRequest struct {
	ID             string    `json:"id"`
	EffectKind     string    `json:"effect_kind"`
	OwnerContext   string    `json:"owner_context"`
	State          string    `json:"state"`
	IdempotencyKey string    `json:"idempotency_key"`
	RequestedAt    time.Time `json:"requested_at"`
}

type ReturnDisposition struct {
	ID               string                `json:"id"`
	ReceiptID        string                `json:"receipt_id"`
	InventoryAction  string                `json:"inventory_action"`
	CustomerRemedy   string                `json:"customer_remedy"`
	Notes            string                `json:"notes"`
	DecidedBySubject string                `json:"decided_by_subject"`
	DecidedAt        time.Time             `json:"decided_at"`
	Effects          []ReturnEffectRequest `json:"effect_requests"`
}

type ReturnCase struct {
	AuthorizationID  string             `json:"authorization_id"`
	OrganizationID   string             `json:"organization_id"`
	OrderID          string             `json:"order_id"`
	StockUnitID      string             `json:"stock_unit_id"`
	CustomerSubject  string             `json:"customer_subject"`
	AuthorizedAction string             `json:"authorized_action"`
	AuthorizedAt     time.Time          `json:"authorized_at"`
	Receipt          *ReturnReceipt     `json:"receipt,omitempty"`
	Disposition      *ReturnDisposition `json:"disposition,omitempty"`
}

type CustomerJourney struct {
	Appointments []Appointment       `json:"appointments"`
	Quotes       []Quote             `json:"quotes"`
	Handovers    []Handover          `json:"handovers"`
	Exceptions   []DeliveryException `json:"delivery_exceptions"`
}

type Page[T any] struct {
	Items      []T    `json:"items"`
	NextCursor string `json:"next_cursor,omitempty"`
}

type Repository interface {
	AppointmentAgenda(context.Context, string, string, time.Time, time.Time) (AppointmentAgenda, error)
	PublicLocations(context.Context, string) ([]Location, error)
	PublicAppointmentSlots(context.Context, string, string, string, time.Time, time.Time) ([]AppointmentSlot, error)
	CreateAppointmentSlot(context.Context, string, AppointmentSlot, string) (AppointmentSlot, error)
	RequestAppointment(context.Context, string, string, string, Appointment, string, string) (Appointment, bool, error)
	CreateServiceResource(context.Context, string, ServiceResource, string) (ServiceResource, error)
	AssignAppointmentResource(context.Context, string, string, string, string, int64, string) (Appointment, error)
	CreateAvailability(context.Context, string, string, AvailabilityEntry, string) (AvailabilityEntry, error)
	CancelAvailability(context.Context, string, string, string, int64, string, string, string) (AvailabilityEntry, error)
	Availability(context.Context, string, string, string, time.Time, time.Time) ([]AvailabilityEntry, error)
	TransitionAppointment(context.Context, string, string, string, string, string, int64, string, string, string) (Appointment, error)
	CancelCustomerAppointment(context.Context, string, string, string, string, int64, string, string, string) (Appointment, error)
	Leads(context.Context, string, string, int, string) (Page[Lead], error)
	AssignLead(context.Context, string, string, string, string, int64, string) (Lead, error)
	TransitionLead(context.Context, string, string, string, string, string, int64, string) (Lead, error)
	CreateQuote(context.Context, string, string, Quote, string, string) (Quote, bool, error)
	PublishDeliveryChecklist(context.Context, string, string, DeliveryChecklist, string) (DeliveryChecklist, error)
	CompleteDeliveryChecklist(context.Context, string, string, string, string, int64, string, int64, []ChecklistResponse, string) (Handover, error)
	RejectHandover(context.Context, string, string, string, string, int64, string, string, string, string, string) (DeliveryException, error)
	DeliveryExceptions(context.Context, string, string, int) ([]DeliveryException, error)
	ResolveDeliveryException(context.Context, string, string, string, string, int64, string, string, string, string, string, string) (DeliveryResolution, error)
	ReturnCases(context.Context, string, string, int) ([]ReturnCase, error)
	ReceiveReturn(context.Context, string, string, string, string, string, string, string, string, string, string) (ReturnReceipt, error)
	DecideReturn(context.Context, string, string, string, string, string, string, string, string, string, string, string, string) (ReturnDisposition, error)
	AcceptQuote(context.Context, string, string, string, string, int64, string, string, string, string, string) (Quote, error)
	CustomerJourney(context.Context, string, string, string) (CustomerJourney, error)
	AcceptHandover(context.Context, string, string, string, string, int64, string, string, int64, string, string) (Handover, error)
}

type IDGenerator interface{ New() string }
type Clock interface{ Now() time.Time }

type Service struct {
	repository Repository
	ids        IDGenerator
	clock      Clock
	policy     *businesspolicy.Profile
}

func NewService(repository Repository, ids IDGenerator, clock Clock) *Service {
	policy := businesspolicy.Reference()
	if bound, ok := repository.(interface{ BusinessPolicySHA256() string }); ok && bound.BusinessPolicySHA256() != policy.SHA256() {
		panic("custom repository policy requires NewServiceWithProfile")
	}
	return &Service{repository: repository, ids: ids, clock: clock, policy: policy}
}

// NewServiceWithProfile binds the same validated policy used by persistence.
// Compatibility constructors retain the reference profile and historical keys.
func NewServiceWithProfile(repository Repository, ids IDGenerator, clock Clock, policy *businesspolicy.Profile) (*Service, error) {
	bound, ok := repository.(interface{ BusinessPolicySHA256() string })
	if !policy.Valid() || !ok || bound.BusinessPolicySHA256() != policy.SHA256() || ids == nil || clock == nil {
		return nil, businesspolicy.ErrProfile
	}
	return &Service{repository: repository, ids: ids, clock: clock, policy: policy}, nil
}

func (s *Service) PublicLocations(ctx context.Context, tenantCode string) ([]Location, error) {
	if !code(tenantCode) {
		return nil, ErrInvalid
	}
	return s.repository.PublicLocations(ctx, tenantCode)
}

func (s *Service) PublicAppointmentSlots(ctx context.Context, tenantCode, organizationCode, kind string, from, to time.Time) ([]AppointmentSlot, error) {
	if !code(tenantCode) || !code(organizationCode) || !appointmentKinds[kind] || from.Before(s.clock.Now()) || !to.After(from) || to.Sub(from) > 31*24*time.Hour {
		return nil, ErrInvalid
	}
	return s.repository.PublicAppointmentSlots(ctx, tenantCode, organizationCode, kind, from, to)
}

func (s *Service) CreateAppointmentSlot(ctx context.Context, tenant string, value AppointmentSlot) (AppointmentSlot, error) {
	value, err := s.prepareAppointmentSlot(tenant, value)
	if err != nil {
		return AppointmentSlot{}, err
	}
	return s.repository.CreateAppointmentSlot(ctx, tenant, value, s.ids.New())
}

func (s *Service) prepareAppointmentSlot(tenant string, value AppointmentSlot) (AppointmentSlot, error) {
	value.ID = s.ids.New()
	value.State = "open"
	value.Version = 1
	if tenant == "" || value.OrganizationID == "" || !appointmentKinds[value.Kind] || !s.policy.AllowsSlot(value.StartsAt, value.EndsAt, s.clock.Now(), value.Capacity) {
		return AppointmentSlot{}, ErrInvalid
	}
	return value, nil
}

func (s *Service) RequestAppointment(ctx context.Context, tenantCode, organizationCode, idempotencyKey, requestHash string, value Appointment) (Appointment, bool, error) {
	value.ID = s.ids.New()
	value.OrganizationID = ""
	value.State = "requested"
	value.Version = 1
	if !code(tenantCode) || !code(organizationCode) || len(idempotencyKey) < 16 || !sha256Hex(requestHash) || value.LeadID == "" || !appointmentKinds[value.Kind] || value.StartsAt.Before(s.clock.Now().Add(s.policy.LeadTime())) {
		return Appointment{}, false, ErrInvalid
	}
	return s.repository.RequestAppointment(ctx, tenantCode, organizationCode, idempotencyKey, value, requestHash, s.ids.New())
}

func (s *Service) CreateServiceResource(ctx context.Context, tenant string, value ServiceResource) (ServiceResource, error) {
	value, err := s.prepareServiceResource(tenant, value)
	if err != nil {
		return ServiceResource{}, err
	}
	return s.repository.CreateServiceResource(ctx, tenant, value, s.ids.New())
}

func (s *Service) prepareServiceResource(tenant string, value ServiceResource) (ServiceResource, error) {
	value.ID = s.ids.New()
	value.Status = "active"
	value.Version = 1
	if tenant == "" || value.OrganizationID == "" || len(strings.TrimSpace(value.DisplayName)) < 1 || len(value.DisplayName) > 160 || !resourceKinds[value.Kind] || len(value.Skills) < 1 || len(value.Skills) > 4 || ((value.Kind == "employee" || value.Kind == "contractor") != (value.PrincipalSubject != "")) {
		return ServiceResource{}, ErrInvalid
	}
	seen := map[string]bool{}
	for _, skill := range value.Skills {
		if !appointmentKinds[skill] || seen[skill] {
			return ServiceResource{}, ErrInvalid
		}
		seen[skill] = true
	}
	return value, nil
}

func (s *Service) AssignAppointmentResource(ctx context.Context, tenant, organization, appointment, resource string, version int64) (Appointment, error) {
	if tenant == "" || organization == "" || appointment == "" || resource == "" || version < 1 {
		return Appointment{}, ErrInvalid
	}
	return s.repository.AssignAppointmentResource(ctx, tenant, organization, appointment, resource, version, s.ids.New())
}

var appointmentTransitions = map[string]map[string]bool{"requested": {"confirmed": true, "cancelled": true}, "confirmed": {"completed": true, "cancelled": true, "no-show": true}}

func (s *Service) TransitionAppointment(ctx context.Context, tenant, organization, appointment, current, target string, version int64, subject, reason string) (Appointment, error) {
	if tenant == "" || organization == "" || appointment == "" || version < 1 || subject == "" || !appointmentTransitions[current][target] || ((target == "cancelled" || target == "no-show") && !code(reason)) || ((target != "cancelled" && target != "no-show") && reason != "") {
		return Appointment{}, ErrInvalid
	}
	return s.repository.TransitionAppointment(ctx, tenant, organization, appointment, current, target, version, subject, reason, s.ids.New())
}

func (s *Service) Leads(ctx context.Context, tenant, organization string, limit int, after string) (Page[Lead], error) {
	if tenant == "" || organization == "" || limit < 1 || limit > 100 {
		return Page[Lead]{}, ErrInvalid
	}
	return s.repository.Leads(ctx, tenant, organization, limit, after)
}

func (s *Service) AssignLead(ctx context.Context, tenant, organization, lead, subject string, version int64) (Lead, error) {
	if tenant == "" || organization == "" || lead == "" || subject == "" || version < 1 {
		return Lead{}, ErrInvalid
	}
	return s.repository.AssignLead(ctx, tenant, organization, lead, subject, version, s.ids.New())
}

// HTTP commands require actor-aware persistence; legacy callers retain their API.
type auditedLeadRepository interface {
	AssignLeadAs(context.Context, string, string, string, string, int64, string, string) (Lead, error)
	TransitionLeadAs(context.Context, string, string, string, string, string, int64, string, string) (Lead, error)
}

func (s *Service) AssignLeadAs(ctx context.Context, tenant, organization, lead, subject string, version int64, actor string) (Lead, error) {
	repo, ok := s.repository.(auditedLeadRepository)
	if !ok || actor == "" || tenant == "" || organization == "" || lead == "" || subject == "" || version < 1 {
		return Lead{}, ErrInvalid
	}
	return repo.AssignLeadAs(ctx, tenant, organization, lead, subject, version, s.ids.New(), actor)
}

func (s *Service) TransitionLeadAs(ctx context.Context, tenant, organization, lead, current, target string, version int64, actor string) (Lead, error) {
	repo, ok := s.repository.(auditedLeadRepository)
	if !ok || actor == "" || tenant == "" || organization == "" || lead == "" || version < 1 || !leadTransitions[current][target] {
		return Lead{}, ErrInvalid
	}
	return repo.TransitionLeadAs(ctx, tenant, organization, lead, current, target, version, s.ids.New(), actor)
}

var leadTransitions = map[string]map[string]bool{
	"new":       {"contacted": true, "lost": true},
	"contacted": {"qualified": true, "lost": true},
	"qualified": {"converted": true, "lost": true},
}

func (s *Service) TransitionLead(ctx context.Context, tenant, organization, lead, current, target string, version int64) (Lead, error) {
	if tenant == "" || organization == "" || lead == "" || version < 1 || !leadTransitions[current][target] {
		return Lead{}, ErrInvalid
	}
	return s.repository.TransitionLead(ctx, tenant, organization, lead, current, target, version, s.ids.New())
}

func (s *Service) CreateQuote(ctx context.Context, tenant, idempotencyKey, requestHash string, value Quote) (Quote, bool, error) {
	return s.createQuote(ctx, tenant, idempotencyKey, requestHash, value, "")
}
func (s *Service) CreateQuoteAs(ctx context.Context, tenant, subject, key, hash string, value Quote) (Quote, bool, error) {
	if subject == "" {
		return Quote{}, false, ErrInvalid
	}
	return s.createQuote(ctx, tenant, key, hash, value, subject)
}
func (s *Service) createQuote(ctx context.Context, tenant, idempotencyKey, requestHash string, value Quote, actor string) (Quote, bool, error) {
	value.ID = s.ids.New()
	value.State = "issued"
	value.Version = 1
	now := s.clock.Now()
	if tenant == "" || len(idempotencyKey) < 16 || !sha256Hex(requestHash) || value.OrganizationID == "" || value.LeadID == "" || value.VariantID == "" || value.PriceBookID == "" || !value.ValidUntil.After(now) {
		return Quote{}, false, ErrInvalid
	}
	if actor != "" {
		repository, ok := s.repository.(interface {
			CreateQuoteAs(context.Context, string, string, Quote, string, string, string) (Quote, bool, error)
		})
		if !ok {
			return Quote{}, false, ErrConflict
		}
		return repository.CreateQuoteAs(ctx, tenant, idempotencyKey, value, requestHash, s.ids.New(), actor)
	}
	return s.repository.CreateQuote(ctx, tenant, idempotencyKey, value, requestHash, s.ids.New())
}

func (s *Service) CustomerJourney(ctx context.Context, tenant, organization, customer string) (CustomerJourney, error) {
	if tenant == "" || organization == "" || customer == "" {
		return CustomerJourney{}, ErrInvalid
	}
	return s.repository.CustomerJourney(ctx, tenant, organization, customer)
}

func (s *Service) AppointmentAgenda(ctx context.Context, tenant, organization string, from, to time.Time) (AppointmentAgenda, error) {
	if tenant == "" || organization == "" || from.IsZero() || !to.After(from) || to.Sub(from) > 31*24*time.Hour {
		return AppointmentAgenda{}, ErrInvalid
	}
	return s.repository.AppointmentAgenda(ctx, tenant, organization, from, to)
}

var checklistResponseTypes = map[string]bool{"confirmation": true, "text": true, "serial": true, "evidence": true}

func (s *Service) PublishDeliveryChecklist(ctx context.Context, tenant, subject string, value DeliveryChecklist) (DeliveryChecklist, error) {
	value.State = "published"
	if tenant == "" || subject == "" || value.OrganizationID == "" || !code(value.ID) || value.Version < 1 || len(strings.TrimSpace(value.Title)) < 1 || len(value.Title) > 160 || len(value.Items) < 1 || len(value.Items) > 64 {
		return DeliveryChecklist{}, ErrInvalid
	}
	seen := map[string]bool{}
	for index := range value.Items {
		item := &value.Items[index]
		item.Ordinal = index + 1
		if !code(item.ID) || seen[item.ID] || len(strings.TrimSpace(item.Prompt)) < 1 || len(item.Prompt) > 500 || !checklistResponseTypes[item.ResponseType] {
			return DeliveryChecklist{}, ErrInvalid
		}
		seen[item.ID] = true
	}
	return s.repository.PublishDeliveryChecklist(ctx, tenant, subject, value, s.ids.New())
}

func (s *Service) CompleteDeliveryChecklist(ctx context.Context, tenant, organization, subject, handover string, version int64, checklistID string, checklistVersion int64, responses []ChecklistResponse) (Handover, error) {
	if tenant == "" || organization == "" || subject == "" || handover == "" || version < 1 || !code(checklistID) || checklistVersion < 1 || len(responses) < 1 || len(responses) > 64 {
		return Handover{}, ErrInvalid
	}
	seen := map[string]bool{}
	for _, response := range responses {
		if !code(response.ItemID) || seen[response.ItemID] || len(strings.TrimSpace(response.ResponseText)) < 1 || len(response.ResponseText) > 2048 || (response.EvidenceSHA256 != "" && !sha256Hex(response.EvidenceSHA256)) {
			return Handover{}, ErrInvalid
		}
		seen[response.ItemID] = true
	}
	return s.repository.CompleteDeliveryChecklist(ctx, tenant, organization, subject, handover, version, checklistID, checklistVersion, responses, s.ids.New())
}

func (s *Service) RejectHandover(ctx context.Context, tenant, organization, customer, handover string, version int64, reasonCode, details, evidence string) (DeliveryException, error) {
	if tenant == "" || organization == "" || customer == "" || handover == "" || version < 1 || !code(reasonCode) || len(strings.TrimSpace(details)) < 1 || len(details) > 1000 || !sha256Hex(evidence) {
		return DeliveryException{}, ErrInvalid
	}
	return s.repository.RejectHandover(ctx, tenant, organization, customer, handover, version, reasonCode, details, evidence, s.ids.New(), s.ids.New())
}

func (s *Service) DeliveryExceptions(ctx context.Context, tenant, organization string, limit int) ([]DeliveryException, error) {
	if tenant == "" || organization == "" || limit < 1 || limit > 100 {
		return nil, ErrInvalid
	}
	return s.repository.DeliveryExceptions(ctx, tenant, organization, limit)
}

var deliveryResolutionActions = map[string]bool{"correct-and-represent": true, "return": true, "exchange": true}

func (s *Service) ResolveDeliveryException(ctx context.Context, tenant, organization, subject, exceptionID string, version int64, action, notes string) (DeliveryResolution, error) {
	if tenant == "" || organization == "" || subject == "" || exceptionID == "" || version < 1 || !deliveryResolutionActions[action] || len(strings.TrimSpace(notes)) < 1 || len(notes) > 1000 {
		return DeliveryResolution{}, ErrInvalid
	}
	return s.repository.ResolveDeliveryException(ctx, tenant, organization, subject, exceptionID, version, action, notes, s.ids.New(), s.ids.New(), s.ids.New(), s.ids.New())
}

func (s *Service) ReturnCases(ctx context.Context, tenant, organization string, limit int) ([]ReturnCase, error) {
	if tenant == "" || organization == "" || limit < 1 || limit > 100 {
		return nil, ErrInvalid
	}
	return s.repository.ReturnCases(ctx, tenant, organization, limit)
}

var returnConditionCodes = map[string]bool{"sealed": true, "opened": true, "damaged": true, "incomplete": true}

func (s *Service) ReceiveReturn(ctx context.Context, tenant, organization, subject, authorizationID, serialNumber, conditionCode, notes, evidence string) (ReturnReceipt, error) {
	if tenant == "" || organization == "" || subject == "" || authorizationID == "" || len(serialNumber) < 1 || len(serialNumber) > 128 || !returnConditionCodes[conditionCode] || len(strings.TrimSpace(notes)) < 1 || len(notes) > 1000 || !sha256Hex(evidence) {
		return ReturnReceipt{}, ErrInvalid
	}
	return s.repository.ReceiveReturn(ctx, tenant, organization, subject, authorizationID, serialNumber, conditionCode, notes, evidence, s.ids.New(), s.ids.New())
}

var returnInventoryActions = map[string]bool{"quarantine": true, "restock": true, "repair": true, "scrap": true}

func (s *Service) DecideReturn(ctx context.Context, tenant, organization, subject, receiptID, inventoryAction, notes string) (ReturnDisposition, error) {
	if tenant == "" || organization == "" || subject == "" || receiptID == "" || !returnInventoryActions[inventoryAction] || len(strings.TrimSpace(notes)) < 1 || len(notes) > 1000 {
		return ReturnDisposition{}, ErrInvalid
	}
	return s.repository.DecideReturn(ctx, tenant, organization, subject, receiptID, inventoryAction, notes, s.ids.New(), s.ids.New(), s.ids.New(), s.ids.New(), s.ids.New(), s.ids.New())
}

func (s *Service) AcceptQuote(ctx context.Context, tenant, organization, customer, quote string, version int64, evidence string) (Quote, error) {
	if tenant == "" || organization == "" || customer == "" || quote == "" || version < 1 || !sha256Hex(evidence) {
		return Quote{}, ErrInvalid
	}
	return s.repository.AcceptQuote(ctx, tenant, organization, customer, quote, version, evidence, s.ids.New(), s.ids.New(), s.ids.New(), s.ids.New())
}

func (s *Service) AcceptHandover(ctx context.Context, tenant, organization, customer, handover string, version int64, serialNumber, checklistID string, checklistVersion int64, evidence string) (Handover, error) {
	if tenant == "" || organization == "" || customer == "" || handover == "" || version < 1 || serialNumber == "" || len(serialNumber) > 128 || !code(checklistID) || checklistVersion < 1 || !sha256Hex(evidence) {
		return Handover{}, ErrInvalid
	}
	return s.repository.AcceptHandover(ctx, tenant, organization, customer, handover, version, serialNumber, checklistID, checklistVersion, evidence, s.ids.New())
}

func sha256Hex(value string) bool {
	decoded, err := hex.DecodeString(value)
	return err == nil && len(decoded) == 32 && value == strings.ToLower(value)
}

var appointmentKinds = map[string]bool{"consultation": true, "test-drive": true, "delivery": true, "service": true}
var resourceKinds = map[string]bool{"employee": true, "contractor": true, "service-bay": true, "vehicle": true, "equipment": true}

func code(value string) bool {
	if len(value) < 1 || len(value) > 64 || value[0] < 'a' || value[0] > 'z' {
		return false
	}
	for _, r := range value {
		if (r < 'a' || r > 'z') && (r < '0' || r > '9') && r != '-' {
			return false
		}
	}
	return !strings.Contains(value, "--") && !strings.HasSuffix(value, "-")
}

func WrapConflict(operation string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%w: %s", ErrConflict, operation)
}

// QuoteResult only reads the durable result; absence never authorizes a new command.
func (s *Service) QuoteResult(ctx context.Context, tenant, organization, lead, key string) (Quote, error) {
	if tenant == "" || organization == "" || lead == "" || len(key) < 16 || len(key) > 128 {
		return Quote{}, ErrInvalid
	}
	repository, ok := s.repository.(interface {
		QuoteResult(context.Context, string, string, string, string) (Quote, error)
	})
	if !ok {
		return Quote{}, ErrConflict
	}
	return repository.QuoteResult(ctx, tenant, organization, lead, key)
}

var resourceCreationKeyPattern = regexp.MustCompile(`^[A-Za-z0-9_-]{16,128}$`)

func (s *Service) CreateServiceResourceOnce(ctx context.Context, tenant, subject, key, hash string, value ServiceResource) (ServiceResource, bool, error) {
	if subject == "" || !resourceCreationKeyPattern.MatchString(key) || !sha256Hex(hash) {
		return ServiceResource{}, false, ErrInvalid
	}
	value, err := s.prepareServiceResource(tenant, value)
	if err != nil {
		return ServiceResource{}, false, err
	}
	r, ok := s.repository.(interface {
		CreateServiceResourceOnce(context.Context, string, string, string, string, ServiceResource, string) (ServiceResource, bool, error)
	})
	if !ok {
		return ServiceResource{}, false, ErrConflict
	}
	return r.CreateServiceResourceOnce(ctx, tenant, subject, key, hash, value, s.ids.New())
}
func (s *Service) ServiceResourceCreationResult(ctx context.Context, tenant, organization, key string) (ServiceResource, error) {
	if tenant == "" || organization == "" || !resourceCreationKeyPattern.MatchString(key) {
		return ServiceResource{}, ErrInvalid
	}
	r, ok := s.repository.(interface {
		ServiceResourceCreationResult(context.Context, string, string, string) (ServiceResource, error)
	})
	if !ok {
		return ServiceResource{}, ErrConflict
	}
	return r.ServiceResourceCreationResult(ctx, tenant, organization, key)
}

var slotCreationKeyPattern = regexp.MustCompile(`^[A-Za-z0-9_-]{16,128}$`)

func (s *Service) CreateAppointmentSlotOnce(ctx context.Context, tenant, subject, key, hash string, value AppointmentSlot) (AppointmentSlot, bool, error) {
	if subject == "" || !slotCreationKeyPattern.MatchString(key) || !sha256Hex(hash) {
		return AppointmentSlot{}, false, ErrInvalid
	}
	value, err := s.prepareAppointmentSlot(tenant, value)
	if err != nil {
		return AppointmentSlot{}, false, err
	}
	r, ok := s.repository.(interface {
		CreateAppointmentSlotOnce(context.Context, string, string, string, string, AppointmentSlot, string) (AppointmentSlot, bool, error)
	})
	if !ok {
		return AppointmentSlot{}, false, ErrConflict
	}
	return r.CreateAppointmentSlotOnce(ctx, tenant, subject, key, hash, value, s.ids.New())
}
func (s *Service) AppointmentSlotCreationResult(ctx context.Context, tenant, organization, key string) (AppointmentSlot, error) {
	if tenant == "" || organization == "" || !slotCreationKeyPattern.MatchString(key) {
		return AppointmentSlot{}, ErrInvalid
	}
	r, ok := s.repository.(interface {
		AppointmentSlotCreationResult(context.Context, string, string, string) (AppointmentSlot, error)
	})
	if !ok {
		return AppointmentSlot{}, ErrConflict
	}
	return r.AppointmentSlotCreationResult(ctx, tenant, organization, key)
}

func (s *Service) PublishedDeliveryChecklist(ctx context.Context, tenant, organization, id string, version int64) (DeliveryChecklist, error) {
	if tenant == "" || organization == "" || !code(id) || version < 1 {
		return DeliveryChecklist{}, ErrInvalid
	}
	repo, ok := s.repository.(interface {
		PublishedDeliveryChecklist(context.Context, string, string, string, int64) (DeliveryChecklist, error)
	})
	if !ok {
		return DeliveryChecklist{}, ErrConflict
	}
	v, err := repo.PublishedDeliveryChecklist(ctx, tenant, organization, id, version)
	if err != nil {
		return DeliveryChecklist{}, err
	}
	if v.ID != id || v.OrganizationID != organization || v.Version != version || v.State != "published" || len(strings.TrimSpace(v.Title)) < 1 || len(v.Title) > 160 || len(v.Items) < 1 || len(v.Items) > 64 {
		return DeliveryChecklist{}, ErrConflict
	}
	seen := map[string]bool{}
	for index, item := range v.Items {
		if item.Ordinal != index+1 || !code(item.ID) || seen[item.ID] || len(strings.TrimSpace(item.Prompt)) < 1 || len(item.Prompt) > 500 || !checklistResponseTypes[item.ResponseType] {
			return DeliveryChecklist{}, ErrConflict
		}
		seen[item.ID] = true
	}
	return v, nil
}

type ChecklistCompletion struct {
	HandoverID       string              `json:"handover_id"`
	OrganizationID   string              `json:"organization_id"`
	State            string              `json:"state"`
	Version          int64               `json:"version"`
	ChecklistID      string              `json:"checklist_id"`
	ChecklistVersion int64               `json:"checklist_version"`
	CompletedAt      time.Time           `json:"completed_at"`
	ActorSubject     string              `json:"actor_subject"`
	Responses        []ChecklistResponse `json:"responses"`
}

func (s *Service) DeliveryChecklistCompletion(ctx context.Context, tenant, organization, handover string) (ChecklistCompletion, error) {
	if tenant == "" || organization == "" || handover == "" || len(handover) > 128 {
		return ChecklistCompletion{}, ErrInvalid
	}
	repo, ok := s.repository.(interface {
		DeliveryChecklistCompletion(context.Context, string, string, string) (ChecklistCompletion, error)
	})
	if !ok {
		return ChecklistCompletion{}, ErrConflict
	}
	v, err := repo.DeliveryChecklistCompletion(ctx, tenant, organization, handover)
	if err != nil {
		return ChecklistCompletion{}, err
	}
	if v.HandoverID != handover || v.OrganizationID != organization || v.Version < 2 || (v.State != "presented" && v.State != "accepted" && v.State != "rejected") || !code(v.ChecklistID) || v.ChecklistVersion < 1 || v.CompletedAt.IsZero() || strings.TrimSpace(v.ActorSubject) == "" || len(v.ActorSubject) > 255 || len(v.Responses) < 1 || len(v.Responses) > 64 {
		return ChecklistCompletion{}, ErrConflict
	}
	previous := ""
	for _, response := range v.Responses {
		if !code(response.ItemID) || response.ItemID <= previous || len(strings.TrimSpace(response.ResponseText)) < 1 || len(response.ResponseText) > 2048 || (response.EvidenceSHA256 != "" && !sha256Hex(response.EvidenceSHA256)) {
			return ChecklistCompletion{}, ErrConflict
		}
		previous = response.ItemID
	}
	return v, nil
}

func (s *Service) ReturnCaseResult(ctx context.Context, tenant, organization, authorization string) (ReturnCase, error) {
	if tenant == "" || organization == "" || authorization == "" || len(authorization) > 128 {
		return ReturnCase{}, ErrInvalid
	}
	repo, ok := s.repository.(interface {
		ReturnCaseResult(context.Context, string, string, string) (ReturnCase, error)
	})
	if !ok {
		return ReturnCase{}, ErrConflict
	}
	v, err := repo.ReturnCaseResult(ctx, tenant, organization, authorization)
	if err != nil {
		return ReturnCase{}, err
	}
	bounded := func(v string, max int) bool { return len(strings.TrimSpace(v)) > 0 && len(v) <= max }
	if v.AuthorizationID != authorization || v.OrganizationID != organization || !bounded(v.OrderID, 128) || !bounded(v.StockUnitID, 128) || !bounded(v.CustomerSubject, 255) || (v.AuthorizedAction != "return" && v.AuthorizedAction != "exchange") || v.AuthorizedAt.IsZero() {
		return ReturnCase{}, ErrConflict
	}
	if v.Receipt == nil {
		if v.Disposition != nil {
			return ReturnCase{}, ErrConflict
		}
		return v, nil
	}
	receipt := v.Receipt
	if !bounded(receipt.ID, 128) || receipt.AuthorizationID != authorization || receipt.OrganizationID != organization || receipt.OrderID != v.OrderID || receipt.StockUnitID != v.StockUnitID || receipt.CustomerSubject != v.CustomerSubject || (len(receipt.ReceivedSerialNumber) < 1 || len(receipt.ReceivedSerialNumber) > 128) || !returnConditionCodes[receipt.ConditionCode] || !bounded(receipt.Notes, 1000) || !sha256Hex(receipt.EvidenceSHA256) || !bounded(receipt.ReceivedBySubject, 255) || receipt.ReceivedAt.IsZero() {
		return ReturnCase{}, ErrConflict
	}
	if v.Disposition == nil {
		return v, nil
	}
	d := v.Disposition
	remedy := "refund"
	if v.AuthorizedAction == "exchange" {
		remedy = "exchange"
	}
	if !bounded(d.ID, 128) || d.ReceiptID != receipt.ID || !returnInventoryActions[d.InventoryAction] || d.CustomerRemedy != remedy || !bounded(d.Notes, 1000) || !bounded(d.DecidedBySubject, 255) || d.DecidedAt.IsZero() {
		return ReturnCase{}, ErrConflict
	}
	owners := map[string]string{"inventory": "inventory", "accounting": "accounting", remedy: map[string]string{"refund": "payment", "exchange": "fulfillment"}[remedy]}
	if remedy == "refund" {
		owners["fiscal"] = "fiscal"
	}
	if len(d.Effects) != len(owners) {
		return ReturnCase{}, ErrConflict
	}
	seenKinds, seenIDs := map[string]bool{}, map[string]bool{}
	for _, e := range d.Effects {
		owner, ok := owners[e.EffectKind]
		if !ok || owner != e.OwnerContext || seenKinds[e.EffectKind] || seenIDs[e.ID] || !bounded(e.ID, 128) || e.State != "requested" || len(e.IdempotencyKey) < 16 || len(e.IdempotencyKey) > 128 || e.RequestedAt.IsZero() {
			return ReturnCase{}, ErrConflict
		}
		seenKinds[e.EffectKind] = true
		seenIDs[e.ID] = true
	}
	return v, nil
}
````

### FILE: `internal/franchisejourney/service_test.go`

```yaml
block_id: "GO-FRANCHISE-JOURNEY:file:05"
operation: CREATE
provenance: AUTHORED
source: "local implementation governed by the upstream references in metadata"
license: "LicenseRef-Workspace-Owner"
sha256: "6fc88efaef5ae994bb1b786a94e8d86325cd348ff900602dc3e7d835eaee159d"
variables: []
secrets_allowed: false
```

````go
package franchisejourney

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"
)

type fixedIDs struct{ n int }

type auditedLeadFake struct {
	fakeRepository
	actor string
	calls int
}

func (f *auditedLeadFake) AssignLeadAs(_ context.Context, _, organization, lead, subject string, version int64, _, actor string) (Lead, error) {
	f.actor = actor
	f.calls++
	return Lead{ID: lead, OrganizationID: organization, AssignedSubject: subject, Version: version + 1}, nil
}
func (f *auditedLeadFake) TransitionLeadAs(_ context.Context, _, organization, lead, _, target string, version int64, _, actor string) (Lead, error) {
	f.actor = actor
	f.calls++
	return Lead{ID: lead, OrganizationID: organization, State: target, Version: version + 1}, nil
}

func TestLeadAuditedRepositoryRequired(t *testing.T) {
	ids := &fixedIDs{}
	s := NewService(&fakeRepository{}, ids, fixedClock{})
	if _, err := s.AssignLeadAs(context.Background(), "tenant", "org", "lead", "assignee", 1, "actor"); !errors.Is(err, ErrInvalid) {
		t.Fatal("unadmitted assignment fallback", err)
	}
	if _, err := s.TransitionLeadAs(context.Background(), "tenant", "org", "lead", "new", "contacted", 1, "actor"); !errors.Is(err, ErrInvalid) {
		t.Fatal("unadmitted transition fallback", err)
	}
	if ids.n != 0 {
		t.Fatal("generated IDs before audited admission")
	}
}

func TestLeadActorAndValidation(t *testing.T) {
	repo := &auditedLeadFake{}
	s := NewService(repo, &fixedIDs{}, fixedClock{})
	ctx := context.Background()
	if _, err := s.AssignLeadAs(ctx, "tenant", "org", "lead", "assignee", 1, ""); !errors.Is(err, ErrInvalid) {
		t.Fatal("missing actor", err)
	}
	if _, err := s.TransitionLeadAs(ctx, "tenant", "org", "lead", "new", "contacted", 1, ""); !errors.Is(err, ErrInvalid) {
		t.Fatal("missing transition actor", err)
	}
	if _, err := s.TransitionLeadAs(ctx, "tenant", "org", "lead", "lost", "new", 1, "actor"); !errors.Is(err, ErrInvalid) {
		t.Fatal("illegal transition", err)
	}
	if repo.calls != 0 {
		t.Fatal("invalid command reached persistence")
	}
	v, err := s.AssignLeadAs(ctx, "tenant", "org", "lead", "assignee", 1, "actual-actor")
	if err != nil || v.AssignedSubject != "assignee" || v.Version != 2 || repo.actor != "actual-actor" {
		t.Fatal(v, err, repo.actor)
	}
	v, err = s.TransitionLeadAs(ctx, "tenant", "org", "lead", "new", "contacted", 2, "second-actor")
	if err != nil || v.State != "contacted" || v.Version != 3 || repo.actor != "second-actor" || repo.calls != 2 {
		t.Fatal(v, err, repo.actor)
	}
}

func (i *fixedIDs) New() string { i.n++; return "id-" + string(rune('a'+i.n)) }

type fixedClock struct{ now time.Time }

func (c fixedClock) Now() time.Time { return c.now }

type fakeRepository struct {
	appointment           Appointment
	slot                  AppointmentSlot
	assigned              Lead
	transition            Lead
	quote                 Quote
	handover              Handover
	accepted              Quote
	resource              ServiceResource
	resourced             Appointment
	appointmentTransition Appointment
	availability          AvailabilityEntry
	checklist             DeliveryChecklist
	completedChecklist    Handover
	deliveryException     DeliveryException
	deliveryResolution    DeliveryResolution
	returnReceipt         ReturnReceipt
	returnDisposition     ReturnDisposition
}

func (f *fakeRepository) PublicLocations(context.Context, string) ([]Location, error) {
	return []Location{{Code: "cordoba"}}, nil
}
func (f *fakeRepository) PublicAppointmentSlots(context.Context, string, string, string, time.Time, time.Time) ([]AppointmentSlot, error) {
	return []AppointmentSlot{{ID: "slot", Capacity: 2, Booked: 1}}, nil
}
func (f *fakeRepository) CreateAppointmentSlot(_ context.Context, _ string, value AppointmentSlot, _ string) (AppointmentSlot, error) {
	f.slot = value
	return value, nil
}
func (f *fakeRepository) RequestAppointment(_ context.Context, _, _, _ string, value Appointment, _, _ string) (Appointment, bool, error) {
	f.appointment = value
	return value, false, nil
}
func (f *fakeRepository) CreateServiceResource(_ context.Context, _ string, value ServiceResource, _ string) (ServiceResource, error) {
	f.resource = value
	return value, nil
}
func (f *fakeRepository) AssignAppointmentResource(_ context.Context, _, _, appointment, resource string, version int64, _ string) (Appointment, error) {
	f.resourced = Appointment{ID: appointment, ResourceID: resource, State: "requested", Version: version + 1}
	return f.resourced, nil
}
func (f *fakeRepository) CreateAvailability(_ context.Context, _, _ string, value AvailabilityEntry, _ string) (AvailabilityEntry, error) {
	f.availability = value
	return value, nil
}
func (f *fakeRepository) CancelAvailability(_ context.Context, _, _, _ string, version int64, _, _, _ string) (AvailabilityEntry, error) {
	return AvailabilityEntry{State: "cancelled", Version: version + 1}, nil
}
func (f *fakeRepository) Availability(context.Context, string, string, string, time.Time, time.Time) ([]AvailabilityEntry, error) {
	return []AvailabilityEntry{{ID: "availability"}}, nil
}
func (f *fakeRepository) TransitionAppointment(_ context.Context, _, _, appointment, _, target string, version int64, _, _, _ string) (Appointment, error) {
	f.appointmentTransition = Appointment{ID: appointment, State: target, Version: version + 1}
	return f.appointmentTransition, nil
}
func (f *fakeRepository) CancelCustomerAppointment(_ context.Context, _, _, _, appointment string, version int64, _, _, _ string) (Appointment, error) {
	return Appointment{ID: appointment, State: "cancelled", Version: version + 1}, nil
}
func (f *fakeRepository) Leads(context.Context, string, string, int, string) (Page[Lead], error) {
	return Page[Lead]{Items: []Lead{{ID: "lead"}}}, nil
}
func (f *fakeRepository) AssignLead(_ context.Context, _, _, _, subject string, version int64, _ string) (Lead, error) {
	f.assigned = Lead{AssignedSubject: subject, Version: version + 1}
	return f.assigned, nil
}
func (f *fakeRepository) TransitionLead(_ context.Context, _, _, _, _, target string, version int64, _ string) (Lead, error) {
	f.transition = Lead{State: target, Version: version + 1}
	return f.transition, nil
}
func (f *fakeRepository) CreateQuote(_ context.Context, _, _ string, value Quote, _, _ string) (Quote, bool, error) {
	f.quote = value
	return value, false, nil
}
func (f *fakeRepository) PublishDeliveryChecklist(_ context.Context, _, _ string, value DeliveryChecklist, _ string) (DeliveryChecklist, error) {
	f.checklist = value
	return value, nil
}
func (f *fakeRepository) CompleteDeliveryChecklist(_ context.Context, _, _, _, _ string, version int64, checklistID string, checklistVersion int64, _ []ChecklistResponse, _ string) (Handover, error) {
	f.completedChecklist = Handover{State: "presented", Version: version + 1, ChecklistID: checklistID, ChecklistVersion: checklistVersion}
	return f.completedChecklist, nil
}
func (f *fakeRepository) RejectHandover(_ context.Context, _, organization, customer, handover string, _ int64, reason, details, _ string, exceptionID, _ string) (DeliveryException, error) {
	f.deliveryException = DeliveryException{ID: exceptionID, OrganizationID: organization, HandoverID: handover, CustomerSubject: customer, ReasonCode: reason, Details: details, State: "open", Version: 1}
	return f.deliveryException, nil
}
func (f *fakeRepository) DeliveryExceptions(context.Context, string, string, int) ([]DeliveryException, error) {
	return []DeliveryException{{ID: "exception", State: "open"}}, nil
}
func (f *fakeRepository) ResolveDeliveryException(_ context.Context, _, _, _ string, exceptionID string, version int64, action, _ string, successorID, authorizationID, _, _ string) (DeliveryResolution, error) {
	f.deliveryResolution = DeliveryResolution{Exception: DeliveryException{ID: exceptionID, State: "resolved", Version: version + 1, ResolutionAction: action}}
	if action == "correct-and-represent" {
		f.deliveryResolution.SuccessorHandover = &Handover{ID: successorID, State: "prepared", Version: 1}
	} else {
		f.deliveryResolution.ReturnAuthorizationID = authorizationID
		f.deliveryResolution.Disposition = action
	}
	return f.deliveryResolution, nil
}
func (f *fakeRepository) ReturnCases(context.Context, string, string, int) ([]ReturnCase, error) {
	return []ReturnCase{{AuthorizationID: "authorization", AuthorizedAction: "return"}}, nil
}
func (f *fakeRepository) ReceiveReturn(_ context.Context, _, organization, subject, authorizationID, serial, condition, notes, evidence, receiptID, _ string) (ReturnReceipt, error) {
	f.returnReceipt = ReturnReceipt{ID: receiptID, AuthorizationID: authorizationID, OrganizationID: organization, ReceivedSerialNumber: serial, ConditionCode: condition, Notes: notes, EvidenceSHA256: evidence, ReceivedBySubject: subject}
	return f.returnReceipt, nil
}
func (f *fakeRepository) DecideReturn(_ context.Context, _, _, subject, receiptID, inventoryAction, notes, dispositionID, inventoryID, remedyID, accountingID, _, _ string) (ReturnDisposition, error) {
	f.returnDisposition = ReturnDisposition{ID: dispositionID, ReceiptID: receiptID, InventoryAction: inventoryAction, CustomerRemedy: "refund", Notes: notes, DecidedBySubject: subject, Effects: []ReturnEffectRequest{{ID: inventoryID, EffectKind: "inventory"}, {ID: remedyID, EffectKind: "refund"}, {ID: accountingID, EffectKind: "accounting"}}}
	return f.returnDisposition, nil
}
func (f *fakeRepository) AcceptQuote(_ context.Context, _, _, _, _ string, version int64, _ string, orderID, _, _, _ string) (Quote, error) {
	f.accepted = Quote{State: "accepted", Version: version + 1, OrderID: orderID}
	return f.accepted, nil
}
func (f *fakeRepository) CustomerJourney(context.Context, string, string, string) (CustomerJourney, error) {
	return CustomerJourney{Quotes: []Quote{{ID: "quote"}}}, nil
}

func (f *fakeRepository) AppointmentAgenda(context.Context, string, string, time.Time, time.Time) (AppointmentAgenda, error) {
	return AppointmentAgenda{Appointments: []Appointment{}, Resources: []ServiceResource{}}, nil
}
func (f *fakeRepository) AcceptHandover(_ context.Context, _, _, _, _ string, version int64, _, checklistID string, checklistVersion int64, evidence, _ string) (Handover, error) {
	f.handover = Handover{State: "accepted", Version: version + 1, AcceptanceEvidence: evidence, ChecklistID: checklistID, ChecklistVersion: checklistVersion}
	return f.handover, nil
}

func newService() (*Service, *fakeRepository, time.Time) {
	now := time.Date(2026, 8, 29, 20, 0, 0, 0, time.UTC)
	repository := &fakeRepository{}
	return NewService(repository, &fixedIDs{}, fixedClock{now}), repository, now
}

func TestPublicAndAppointmentContracts(t *testing.T) {
	service, repository, now := newService()
	if _, err := service.PublicLocations(context.Background(), "Bad_Code"); !errors.Is(err, ErrInvalid) {
		t.Fatal("unsafe tenant code accepted")
	}
	value := Appointment{LeadID: "lead", ModelID: "model", Kind: "test-drive", StartsAt: now.Add(time.Hour)}
	validHash := strings.Repeat("0", 64)
	created, replayed, err := service.RequestAppointment(context.Background(), "tenant", "store", "appointment-key-1", validHash, value)
	if err != nil || replayed || created.State != "requested" || repository.appointment.Version != 1 {
		t.Fatalf("appointment=%+v replayed=%v err=%v", created, replayed, err)
	}
	value.StartsAt = now
	if _, _, err = service.RequestAppointment(context.Background(), "tenant", "store", "appointment-key-2", validHash, value); !errors.Is(err, ErrInvalid) {
		t.Fatal("past appointment accepted")
	}
	value.StartsAt = now.Add(time.Hour)
	if _, _, err = service.RequestAppointment(context.Background(), "tenant", "store", "appointment-key-3", strings.Repeat("z", 64), value); !errors.Is(err, ErrInvalid) {
		t.Fatal("non-hex request digest accepted")
	}
}

func TestAppointmentCapacityContracts(t *testing.T) {
	service, repository, now := newService()
	created, err := service.CreateAppointmentSlot(context.Background(), "tenant", AppointmentSlot{OrganizationID: "store", Kind: "service", StartsAt: now.Add(time.Hour), EndsAt: now.Add(2 * time.Hour), Capacity: 3})
	if err != nil || created.State != "open" || created.Version != 1 || repository.slot.ID == "" {
		t.Fatalf("slot=%+v err=%v", created, err)
	}
	items, err := service.PublicAppointmentSlots(context.Background(), "tenant", "store", "service", now.Add(time.Minute), now.Add(24*time.Hour))
	if err != nil || len(items) != 1 || items[0].Booked != 1 {
		t.Fatalf("slots=%+v err=%v", items, err)
	}
	invalid := []AppointmentSlot{
		{OrganizationID: "store", Kind: "unknown", StartsAt: now.Add(time.Hour), EndsAt: now.Add(2 * time.Hour), Capacity: 1},
		{OrganizationID: "store", Kind: "service", StartsAt: now.Add(time.Hour), EndsAt: now.Add(10 * time.Hour), Capacity: 1},
		{OrganizationID: "store", Kind: "service", StartsAt: now.Add(time.Hour), EndsAt: now.Add(2 * time.Hour), Capacity: 0},
	}
	for _, value := range invalid {
		if _, err = service.CreateAppointmentSlot(context.Background(), "tenant", value); !errors.Is(err, ErrInvalid) {
			t.Fatalf("invalid slot accepted: %+v err=%v", value, err)
		}
	}
	if _, err = service.PublicAppointmentSlots(context.Background(), "tenant", "store", "service", now.Add(time.Minute), now.Add(32*24*time.Hour)); !errors.Is(err, ErrInvalid) {
		t.Fatal("unbounded public slot range accepted")
	}
}

func TestAppointmentResourceAndLifecycleContracts(t *testing.T) {
	service, _, _ := newService()
	resource, err := service.CreateServiceResource(context.Background(), "tenant", ServiceResource{OrganizationID: "store", PrincipalSubject: "technician", DisplayName: "Technician", Kind: "employee", Skills: []string{"service", "delivery"}})
	if err != nil || resource.Status != "active" || resource.Version != 1 || resource.ID == "" {
		t.Fatalf("resource=%+v err=%v", resource, err)
	}
	if _, err = service.CreateServiceResource(context.Background(), "tenant", ServiceResource{OrganizationID: "store", DisplayName: "Invalid employee", Kind: "employee", Skills: []string{"service"}}); !errors.Is(err, ErrInvalid) {
		t.Fatal("employee without principal accepted")
	}
	if _, err = service.CreateServiceResource(context.Background(), "tenant", ServiceResource{OrganizationID: "store", DisplayName: "Bay", Kind: "service-bay", Skills: []string{"service", "service"}}); !errors.Is(err, ErrInvalid) {
		t.Fatal("duplicate resource skill accepted")
	}
	assigned, err := service.AssignAppointmentResource(context.Background(), "tenant", "store", "appointment", resource.ID, 1)
	if err != nil || assigned.ResourceID != resource.ID || assigned.Version != 2 {
		t.Fatalf("assigned=%+v err=%v", assigned, err)
	}
	confirmed, err := service.TransitionAppointment(context.Background(), "tenant", "store", "appointment", "requested", "confirmed", 2, "operator", "")
	if err != nil || confirmed.State != "confirmed" || confirmed.Version != 3 {
		t.Fatalf("confirmed=%+v err=%v", confirmed, err)
	}
	if _, err = service.TransitionAppointment(context.Background(), "tenant", "store", "appointment", "requested", "completed", 2, "operator", ""); !errors.Is(err, ErrInvalid) {
		t.Fatal("appointment skipped confirmation")
	}
}

func TestLeadStateMachineAndAssignment(t *testing.T) {
	service, _, _ := newService()
	assigned, err := service.AssignLead(context.Background(), "tenant", "store", "lead", "sales-person", 1)
	if err != nil || assigned.AssignedSubject != "sales-person" || assigned.Version != 2 {
		t.Fatalf("assignment=%+v err=%v", assigned, err)
	}
	value, err := service.TransitionLead(context.Background(), "tenant", "store", "lead", "new", "contacted", 2)
	if err != nil || value.State != "contacted" || value.Version != 3 {
		t.Fatalf("transition=%+v err=%v", value, err)
	}
	if _, err = service.TransitionLead(context.Background(), "tenant", "store", "lead", "new", "converted", 2); !errors.Is(err, ErrInvalid) {
		t.Fatal("lead skipped states")
	}
}

func TestQuoteAndCustomerOwnershipInputs(t *testing.T) {
	service, _, now := newService()
	quote, replayed, err := service.CreateQuote(context.Background(), "tenant", "quote-request-0001", strings.Repeat("1", 64), Quote{OrganizationID: "store", LeadID: "lead", VariantID: "variant", PriceBookID: "retail", Currency: "ARS", TotalMinorUnits: 1000, ValidUntil: now.Add(24 * time.Hour)})
	if err != nil || replayed || quote.State != "issued" || quote.Version != 1 {
		t.Fatalf("quote=%+v replayed=%v err=%v", quote, replayed, err)
	}
	if _, _, err = service.CreateQuote(context.Background(), "tenant", "quote-request-0002", strings.Repeat("2", 64), Quote{OrganizationID: "store", LeadID: "lead", VariantID: "variant", ValidUntil: now.Add(time.Hour)}); !errors.Is(err, ErrInvalid) {
		t.Fatal("missing price book accepted")
	}
	if _, _, err = service.CreateQuote(context.Background(), "tenant", "short", strings.Repeat("2", 64), Quote{OrganizationID: "store", LeadID: "lead", VariantID: "variant", PriceBookID: "retail", ValidUntil: now.Add(time.Hour)}); !errors.Is(err, ErrInvalid) {
		t.Fatal("short quote idempotency key accepted")
	}
	if _, err = service.CustomerJourney(context.Background(), "tenant", "store", ""); !errors.Is(err, ErrInvalid) {
		t.Fatal("empty customer accepted")
	}
	evidence := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	checklist, err := service.PublishDeliveryChecklist(context.Background(), "tenant", "operator", DeliveryChecklist{ID: "standard-delivery", OrganizationID: "store", Version: 2, Title: "Entrega estándar", Items: []ChecklistItem{{ID: "serial-observed", Prompt: "Verificar serie", ResponseType: "serial", Required: true}, {ID: "asset-condition", Prompt: "Confirmar condición", ResponseType: "confirmation", Required: true}}})
	if err != nil || checklist.State != "published" || len(checklist.Items) != 2 || checklist.Items[1].Ordinal != 2 {
		t.Fatalf("checklist=%+v err=%v", checklist, err)
	}
	completed, err := service.CompleteDeliveryChecklist(context.Background(), "tenant", "store", "operator", "handover", 1, checklist.ID, checklist.Version, []ChecklistResponse{{ItemID: "serial-observed", ResponseText: "SERIAL"}, {ItemID: "asset-condition", ResponseText: "confirmed"}})
	if err != nil || completed.State != "presented" || completed.Version != 2 {
		t.Fatalf("completed=%+v err=%v", completed, err)
	}
	handover, err := service.AcceptHandover(context.Background(), "tenant", "store", "customer", "handover", 2, "SERIAL", checklist.ID, checklist.Version, evidence)
	if err != nil || handover.State != "accepted" {
		t.Fatalf("handover=%+v err=%v", handover, err)
	}
	if _, err = service.AcceptHandover(context.Background(), "tenant", "store", "customer", "handover", 2, "", checklist.ID, checklist.Version, evidence); !errors.Is(err, ErrInvalid) {
		t.Fatalf("handover without serial accepted: %v", err)
	}
	if _, err = service.PublishDeliveryChecklist(context.Background(), "tenant", "operator", DeliveryChecklist{ID: "duplicate", OrganizationID: "store", Version: 1, Title: "Duplicate", Items: []ChecklistItem{{ID: "same", Prompt: "A", ResponseType: "text"}, {ID: "same", Prompt: "B", ResponseType: "text"}}}); !errors.Is(err, ErrInvalid) {
		t.Fatal("duplicate checklist item accepted")
	}
	rejected, err := service.RejectHandover(context.Background(), "tenant", "store", "customer", "handover", 2, "visible-damage", "Rayón en cuadro", evidence)
	if err != nil || rejected.State != "open" || rejected.ReasonCode != "visible-damage" {
		t.Fatalf("rejected=%+v err=%v", rejected, err)
	}
	if _, err = service.RejectHandover(context.Background(), "tenant", "store", "customer", "handover", 2, "Bad_Code", "detail", evidence); !errors.Is(err, ErrInvalid) {
		t.Fatal("invalid rejection reason accepted")
	}
	resolved, err := service.ResolveDeliveryException(context.Background(), "tenant", "store", "operator", rejected.ID, 1, "correct-and-represent", "Corregir preparación")
	if err != nil || resolved.SuccessorHandover == nil || resolved.SuccessorHandover.State != "prepared" {
		t.Fatalf("resolved=%+v err=%v", resolved, err)
	}
	if _, err = service.ResolveDeliveryException(context.Background(), "tenant", "store", "operator", rejected.ID, 1, "refund-now", "invalid"); !errors.Is(err, ErrInvalid) {
		t.Fatal("unknown resolution action accepted")
	}
	receipt, err := service.ReceiveReturn(context.Background(), "tenant", "store", "operator", "authorization", "SERIAL", "damaged", "Recepción con daño", evidence)
	if err != nil || receipt.AuthorizationID != "authorization" || receipt.ConditionCode != "damaged" || receipt.ReceivedBySubject != "operator" {
		t.Fatalf("receipt=%+v err=%v", receipt, err)
	}
	if _, err = service.ReceiveReturn(context.Background(), "tenant", "store", "operator", "authorization", "SERIAL", "unknown", "invalid", evidence); !errors.Is(err, ErrInvalid) {
		t.Fatal("unknown return condition accepted")
	}
	disposition, err := service.DecideReturn(context.Background(), "tenant", "store", "operator", receipt.ID, "quarantine", "Separar hasta efectos downstream")
	if err != nil || disposition.InventoryAction != "quarantine" || disposition.CustomerRemedy != "refund" || len(disposition.Effects) != 3 {
		t.Fatalf("disposition=%+v err=%v", disposition, err)
	}
	if _, err = service.DecideReturn(context.Background(), "tenant", "store", "operator", receipt.ID, "sell-as-new", "invalid"); !errors.Is(err, ErrInvalid) {
		t.Fatal("unknown inventory action accepted")
	}
	accepted, err := service.AcceptQuote(context.Background(), "tenant", "store", "customer", "quote", 1, evidence)
	if err != nil || accepted.State != "accepted" || accepted.OrderID == "" || accepted.Version != 2 {
		t.Fatalf("accepted quote=%+v err=%v", accepted, err)
	}
	if _, err = service.AcceptQuote(context.Background(), "tenant", "store", "customer", "quote", 1, "not-a-digest"); !errors.Is(err, ErrInvalid) {
		t.Fatal("invalid quote acceptance evidence accepted")
	}
}

type checklistReadFake struct {
	fakeRepository
	value DeliveryChecklist
}

func (f *checklistReadFake) PublishedDeliveryChecklist(context.Context, string, string, string, int64) (DeliveryChecklist, error) {
	return f.value, nil
}
func TestPublishedChecklistRejectsIncoherentProjection(t *testing.T) {
	valid := func() DeliveryChecklist {
		return DeliveryChecklist{ID: "checklist", OrganizationID: "store", Version: 1, State: "published", Title: "Synthetic", Items: []ChecklistItem{{ID: "serial", Ordinal: 1, Prompt: "Serie", ResponseType: "serial", Required: true}}}
	}
	cases := map[string]func(*DeliveryChecklist){"id": func(v *DeliveryChecklist) { v.ID = "other" }, "organization": func(v *DeliveryChecklist) { v.OrganizationID = "other" }, "version": func(v *DeliveryChecklist) { v.Version = 2 }, "draft": func(v *DeliveryChecklist) { v.State = "draft" }, "blank title": func(v *DeliveryChecklist) { v.Title = " " }, "long title": func(v *DeliveryChecklist) { v.Title = strings.Repeat("á", 81) }, "no items": func(v *DeliveryChecklist) { v.Items = nil }, "ordinal": func(v *DeliveryChecklist) { v.Items[0].Ordinal = 2 }, "prompt": func(v *DeliveryChecklist) { v.Items[0].Prompt = " " }, "long prompt": func(v *DeliveryChecklist) { v.Items[0].Prompt = strings.Repeat("á", 251) }, "response type": func(v *DeliveryChecklist) { v.Items[0].ResponseType = "unknown" }, "duplicate id": func(v *DeliveryChecklist) { v.Items = append(v.Items, v.Items[0]); v.Items[1].Ordinal = 2 }}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			value := valid()
			mutate(&value)
			s := NewService(&checklistReadFake{value: value}, &fixedIDs{}, fixedClock{})
			got, e := s.PublishedDeliveryChecklist(context.Background(), "tenant", "store", "checklist", 1)
			if !errors.Is(e, ErrConflict) || got.ID != "" {
				t.Fatal(got, e)
			}
		})
	}
	s := NewService(&checklistReadFake{value: valid()}, &fixedIDs{}, fixedClock{})
	if got, e := s.PublishedDeliveryChecklist(context.Background(), "tenant", "store", "checklist", 1); e != nil || got.ID != "checklist" {
		t.Fatal(got, e)
	}
}

type completionReadFake struct {
	fakeRepository
	value ChecklistCompletion
}

func (f *completionReadFake) DeliveryChecklistCompletion(context.Context, string, string, string) (ChecklistCompletion, error) {
	return f.value, nil
}
func TestChecklistCompletionRejectsIncoherentProjection(t *testing.T) {
	valid := func() ChecklistCompletion {
		return ChecklistCompletion{HandoverID: "handover", OrganizationID: "store", State: "presented", Version: 2, ChecklistID: "checklist", ChecklistVersion: 1, CompletedAt: time.Now().UTC(), ActorSubject: "operator", Responses: []ChecklistResponse{{ItemID: "confirmed", ResponseText: "confirmed"}, {ItemID: "serial", ResponseText: "Synthetic"}}}
	}
	cases := map[string]func(*ChecklistCompletion){"handover": func(v *ChecklistCompletion) { v.HandoverID = "other" }, "organization": func(v *ChecklistCompletion) { v.OrganizationID = "other" }, "version": func(v *ChecklistCompletion) { v.Version = 1 }, "prepared": func(v *ChecklistCompletion) { v.State = "prepared" }, "checklist": func(v *ChecklistCompletion) { v.ChecklistID = "" }, "checklist version": func(v *ChecklistCompletion) { v.ChecklistVersion = 0 }, "time": func(v *ChecklistCompletion) { v.CompletedAt = time.Time{} }, "actor": func(v *ChecklistCompletion) { v.ActorSubject = " " }, "no responses": func(v *ChecklistCompletion) { v.Responses = nil }, "order": func(v *ChecklistCompletion) { v.Responses[0], v.Responses[1] = v.Responses[1], v.Responses[0] }, "duplicate": func(v *ChecklistCompletion) { v.Responses[1].ItemID = v.Responses[0].ItemID }, "blank": func(v *ChecklistCompletion) { v.Responses[0].ResponseText = " " }, "too long": func(v *ChecklistCompletion) { v.Responses[0].ResponseText = strings.Repeat("á", 1025) }, "evidence": func(v *ChecklistCompletion) { v.Responses[0].EvidenceSHA256 = "invalid" }}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			v := valid()
			mutate(&v)
			s := NewService(&completionReadFake{value: v}, &fixedIDs{}, fixedClock{})
			got, e := s.DeliveryChecklistCompletion(context.Background(), "tenant", "store", "handover")
			if !errors.Is(e, ErrConflict) || got.HandoverID != "" {
				t.Fatal(got, e)
			}
		})
	}
	for _, state := range []string{"presented", "accepted", "rejected"} {
		v := valid()
		v.State = state
		s := NewService(&completionReadFake{value: v}, &fixedIDs{}, fixedClock{})
		if got, e := s.DeliveryChecklistCompletion(context.Background(), "tenant", "store", "handover"); e != nil || got.State != state {
			t.Fatal(got, e)
		}
	}
}

type returnCaseReadFake struct {
	fakeRepository
	value ReturnCase
}

func (f *returnCaseReadFake) ReturnCaseResult(context.Context, string, string, string) (ReturnCase, error) {
	return f.value, nil
}
func returnResultFixture(action string) ReturnCase {
	now := time.Date(2026, 9, 8, 0, 0, 0, 0, time.UTC)
	remedy := "refund"
	if action == "exchange" {
		remedy = "exchange"
	}
	v := ReturnCase{AuthorizationID: "authorization", OrganizationID: "store", OrderID: "order", StockUnitID: "stock", CustomerSubject: "customer", AuthorizedAction: action, AuthorizedAt: now, Receipt: &ReturnReceipt{ID: "receipt", AuthorizationID: "authorization", OrganizationID: "store", OrderID: "order", StockUnitID: "stock", CustomerSubject: "customer", ReceivedSerialNumber: "SYNTHETIC-SERIAL", ConditionCode: "sealed", Notes: "Synthetic receipt", EvidenceSHA256: strings.Repeat("a", 64), ReceivedBySubject: "receiver", ReceivedAt: now}, Disposition: &ReturnDisposition{ID: "disposition", ReceiptID: "receipt", InventoryAction: "quarantine", CustomerRemedy: remedy, Notes: "Synthetic decision", DecidedBySubject: "decider", DecidedAt: now}}
	requested := []struct{ kind, owner string }{{"inventory", "inventory"}, {remedy, map[string]string{"refund": "payment", "exchange": "fulfillment"}[remedy]}, {"accounting", "accounting"}}
	if action == "return" {
		requested = append(requested, struct{ kind, owner string }{"fiscal", "fiscal"})
	}
	for _, item := range requested {
		id := "request-" + item.kind + "-0001"
		v.Disposition.Effects = append(v.Disposition.Effects, ReturnEffectRequest{ID: id, EffectKind: item.kind, OwnerContext: item.owner, State: "requested", IdempotencyKey: id, RequestedAt: now})
	}
	return v
}
func TestReturnCaseResultRejectsIncoherentProjection(t *testing.T) {
	cases := map[string]func(*ReturnCase){
		"authorization": func(v *ReturnCase) { v.AuthorizationID = "other" }, "organization": func(v *ReturnCase) { v.OrganizationID = "other" }, "action": func(v *ReturnCase) { v.AuthorizedAction = "unknown" }, "authorized time": func(v *ReturnCase) { v.AuthorizedAt = time.Time{} },
		"receipt authorization": func(v *ReturnCase) { v.Receipt.AuthorizationID = "other" }, "receipt organization": func(v *ReturnCase) { v.Receipt.OrganizationID = "other" }, "receipt order": func(v *ReturnCase) { v.Receipt.OrderID = "other" }, "receipt stock": func(v *ReturnCase) { v.Receipt.StockUnitID = "other" }, "receipt customer": func(v *ReturnCase) { v.Receipt.CustomerSubject = "other" }, "serial": func(v *ReturnCase) { v.Receipt.ReceivedSerialNumber = "" }, "condition": func(v *ReturnCase) { v.Receipt.ConditionCode = "unknown" }, "receipt notes": func(v *ReturnCase) { v.Receipt.Notes = " " }, "evidence": func(v *ReturnCase) { v.Receipt.EvidenceSHA256 = "invalid" }, "receiver": func(v *ReturnCase) { v.Receipt.ReceivedBySubject = "" }, "received time": func(v *ReturnCase) { v.Receipt.ReceivedAt = time.Time{} },
		"orphan decision": func(v *ReturnCase) { v.Receipt = nil }, "decision receipt": func(v *ReturnCase) { v.Disposition.ReceiptID = "other" }, "inventory": func(v *ReturnCase) { v.Disposition.InventoryAction = "unknown" }, "remedy": func(v *ReturnCase) { v.Disposition.CustomerRemedy = "exchange" }, "decision notes": func(v *ReturnCase) { v.Disposition.Notes = strings.Repeat("á", 501) }, "decider": func(v *ReturnCase) { v.Disposition.DecidedBySubject = "" }, "decided time": func(v *ReturnCase) { v.Disposition.DecidedAt = time.Time{} },
		"missing effect": func(v *ReturnCase) { v.Disposition.Effects = v.Disposition.Effects[:3] }, "effect owner": func(v *ReturnCase) { v.Disposition.Effects[0].OwnerContext = "payment" }, "effect kind": func(v *ReturnCase) { v.Disposition.Effects[0].EffectKind = "unknown" }, "duplicate id": func(v *ReturnCase) { v.Disposition.Effects[1].ID = v.Disposition.Effects[0].ID }, "duplicate kind": func(v *ReturnCase) { v.Disposition.Effects[1].EffectKind = v.Disposition.Effects[0].EffectKind }, "state": func(v *ReturnCase) { v.Disposition.Effects[0].State = "completed" }, "key": func(v *ReturnCase) { v.Disposition.Effects[0].IdempotencyKey = "short" }, "requested time": func(v *ReturnCase) { v.Disposition.Effects[0].RequestedAt = time.Time{} },
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			v := returnResultFixture("return")
			mutate(&v)
			s := NewService(&returnCaseReadFake{value: v}, &fixedIDs{}, fixedClock{})
			got, e := s.ReturnCaseResult(context.Background(), "tenant", "store", "authorization")
			if !errors.Is(e, ErrConflict) || got.AuthorizationID != "" {
				t.Fatal(got, e)
			}
		})
	}
	for _, action := range []string{"return", "exchange"} {
		for _, stage := range []string{"authorized", "received", "decided"} {
			v := returnResultFixture(action)
			if stage == "authorized" {
				v.Receipt = nil
				v.Disposition = nil
			}
			if stage == "received" {
				v.Disposition = nil
			}
			s := NewService(&returnCaseReadFake{value: v}, &fixedIDs{}, fixedClock{})
			if got, e := s.ReturnCaseResult(context.Background(), "tenant", "store", "authorization"); e != nil || got.AuthorizedAction != action {
				t.Fatal(stage, got, e)
			}
		}
	}
}

// Domain seeds model the authorized, received and decided fixtures exercised against PostgreSQL.
// Invariants: reads mint no IDs, never mutate input, expose no partial failure result,
// preserve authorization/receipt scope and retain the existing refund-four/exchange-three request rule.
func FuzzReturnCaseResult(f *testing.F) {
	for _, action := range []string{"return", "exchange"} {
		for _, stage := range []string{"authorized", "received", "decided", "cross-stock"} {
			v := returnResultFixture(action)
			if stage == "authorized" {
				v.Receipt = nil
				v.Disposition = nil
			}
			if stage == "received" {
				v.Disposition = nil
			}
			if stage == "cross-stock" {
				v.Receipt.StockUnitID = "foreign-stock"
			}
			raw, e := json.Marshal(v)
			if e != nil {
				f.Fatal(e)
			}
			f.Add(raw)
		}
	}
	f.Fuzz(func(t *testing.T, raw []byte) {
		if len(raw) > 65536 {
			return
		}
		var v ReturnCase
		if json.Unmarshal(raw, &v) != nil {
			return
		}
		before, e := json.Marshal(v)
		if e != nil {
			return
		}
		ids := &fixedIDs{}
		s := NewService(&returnCaseReadFake{value: v}, ids, fixedClock{})
		got, err := s.ReturnCaseResult(context.Background(), "tenant", "store", "authorization")
		after, e := json.Marshal(v)
		if e != nil || string(before) != string(after) {
			t.Fatal("read mutated its source")
		}
		if ids.n != 0 {
			t.Fatal("read generated mutation IDs")
		}
		if err != nil {
			if !reflect.DeepEqual(got, ReturnCase{}) {
				t.Fatal("error exposed a partial result")
			}
			return
		}
		if got.AuthorizationID != "authorization" || got.OrganizationID != "store" {
			t.Fatal("authorization scope escaped")
		}
		if !reflect.DeepEqual(got, v) {
			t.Fatal("read rewrote source evidence")
		}
		if got.Receipt == nil {
			if got.Disposition != nil {
				t.Fatal("decision without receipt")
			}
			return
		}
		r := got.Receipt
		if r.AuthorizationID != got.AuthorizationID || r.OrganizationID != got.OrganizationID || r.OrderID != got.OrderID || r.StockUnitID != got.StockUnitID || r.CustomerSubject != got.CustomerSubject {
			t.Fatal("receipt binding escaped")
		}
		if got.Disposition == nil {
			return
		}
		d := got.Disposition
		if d.ReceiptID != r.ID {
			t.Fatal("decision binding escaped")
		}
		expected := map[string]string{"inventory": "inventory", "accounting": "accounting"}
		switch got.AuthorizedAction {
		case "return":
			if d.CustomerRemedy != "refund" {
				t.Fatal("refund contract")
			}
			expected["refund"] = "payment"
			expected["fiscal"] = "fiscal"
		case "exchange":
			if d.CustomerRemedy != "exchange" {
				t.Fatal("exchange contract")
			}
			expected["exchange"] = "fulfillment"
		default:
			t.Fatal("unknown authorization action")
		}
		if len(d.Effects) != len(expected) {
			t.Fatal("request count contradicts authorization")
		}
		for _, effect := range d.Effects {
			owner, exists := expected[effect.EffectKind]
			if !exists || effect.OwnerContext != owner || effect.State != "requested" {
				t.Fatal("request contract escaped")
			}
			delete(expected, effect.EffectKind)
		}
		if len(expected) != 0 {
			t.Fatal("missing request owner")
		}
	})
}
````

### FILE: `internal/platform/postgres/franchisejourney.go`

```yaml
block_id: "GO-FRANCHISE-JOURNEY:file:06"
operation: CREATE
provenance: AUTHORED
source: "local implementation governed by the upstream references in metadata"
license: "LicenseRef-Workspace-Owner"
sha256: "04ccb4100d045e30830604d4165f44df2ac8e7882691f61eab82f14a3c971910"
variables: []
secrets_allowed: false
```

````go
package postgres

import (
	"context"
	"elite.local/enterprise/internal/businesspolicy"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"elite.local/enterprise/internal/bcsales"
	"elite.local/enterprise/internal/franchisejourney"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FranchiseJourney struct {
	pool   *pgxpool.Pool
	policy *businesspolicy.Profile
}

func NewFranchiseJourney(pool *pgxpool.Pool) *FranchiseJourney {
	return &FranchiseJourney{pool: pool, policy: businesspolicy.Reference()}
}

func NewFranchiseJourneyWithProfile(pool *pgxpool.Pool, policy *businesspolicy.Profile) (*FranchiseJourney, error) {
	if pool == nil || !policy.Valid() {
		return nil, businesspolicy.ErrProfile
	}
	return &FranchiseJourney{pool: pool, policy: policy}, nil
}
func (r *FranchiseJourney) BusinessPolicySHA256() string {
	if r == nil {
		return ""
	}
	return r.policy.SHA256()
}

// Read model only; assignment/transition owners remain authoritative for writes.
func (r *FranchiseJourney) AppointmentAgenda(ctx context.Context, tenant, organization string, from, to time.Time) (franchisejourney.AppointmentAgenda, error) {
	result := franchisejourney.AppointmentAgenda{Appointments: []franchisejourney.Appointment{}, Resources: []franchisejourney.ServiceResource{}}
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead, AccessMode: pgx.ReadOnly})
	if err != nil {
		return result, err
	}
	defer tx.Rollback(ctx)
	rows, err := tx.Query(ctx, `select a.appointment_id,a.organization_id,a.lead_id,coalesce(a.model_id,''),a.appointment_kind,a.starts_at,a.state,a.version,coalesce(a.slot_id,''),coalesce(a.ends_at,a.starts_at),coalesce(ar.resource_id,'')
 from crm.appointment a left join crm.appointment_resource ar on ar.tenant_id=a.tenant_id and ar.appointment_id=a.appointment_id
 where a.tenant_id=$1 and a.organization_id=$2 and a.starts_at >= $3 and a.starts_at < $4
 order by a.starts_at,a.appointment_id limit 201`, tenant, organization, from, to)
	if err != nil {
		return result, err
	}
	for rows.Next() {
		var a franchisejourney.Appointment
		if err := rows.Scan(&a.ID, &a.OrganizationID, &a.LeadID, &a.ModelID, &a.Kind, &a.StartsAt, &a.State, &a.Version, &a.SlotID, &a.EndsAt, &a.ResourceID); err != nil {
			rows.Close()
			return result, err
		}
		result.Appointments = append(result.Appointments, a)
	}
	err = rows.Err()
	rows.Close()
	if err != nil {
		return result, err
	}
	rows, err = tx.Query(ctx, `select r.resource_id,r.organization_id,r.display_name,r.resource_kind,r.status,r.version,
 coalesce(array(select s.appointment_kind from crm.resource_skill s where s.tenant_id=r.tenant_id and s.resource_id=r.resource_id order by s.appointment_kind),'{}')
 from crm.service_resource r where r.tenant_id=$1 and r.organization_id=$2 and r.status='active'
 order by r.display_name,r.resource_id limit 201`, tenant, organization)
	if err != nil {
		return result, err
	}
	for rows.Next() {
		var item franchisejourney.ServiceResource
		if err := rows.Scan(&item.ID, &item.OrganizationID, &item.DisplayName, &item.Kind, &item.Status, &item.Version, &item.Skills); err != nil {
			rows.Close()
			return result, err
		}
		result.Resources = append(result.Resources, item)
	}
	err = rows.Err()
	rows.Close()
	if err != nil {
		return result, err
	}
	if len(result.Appointments) > 200 {
		result.Truncated = true
		result.Appointments = result.Appointments[:200]
	}
	if len(result.Resources) > 200 {
		result.Truncated = true
		result.Resources = result.Resources[:200]
	}
	return result, tx.Commit(ctx)
}

func (r *FranchiseJourney) PublicLocations(ctx context.Context, tenantCode string) ([]franchisejourney.Location, error) {
	rows, err := r.pool.Query(ctx, `select o.organization_id,o.organization_code,o.display_name,l.city,l.region,l.country,coalesce(l.contact_phone,''),coalesce(l.contact_email,'') from platform.tenant t join org.organization o on o.tenant_id=t.tenant_id join org.public_location l on l.tenant_id=o.tenant_id and l.organization_id=o.organization_id where t.tenant_code=$1 and t.status='active' and o.status='active' and l.published=true order by l.sort_order,o.organization_id`, tenantCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []franchisejourney.Location{}
	for rows.Next() {
		var item franchisejourney.Location
		if err := rows.Scan(&item.OrganizationID, &item.Code, &item.Name, &item.City, &item.Region, &item.Country, &item.ContactPhone, &item.ContactEmail); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *FranchiseJourney) PublicAppointmentSlots(ctx context.Context, tenantCode, organizationCode, kind string, from, to time.Time) ([]franchisejourney.AppointmentSlot, error) {
	rows, err := r.pool.Query(ctx, `select s.slot_id,s.organization_id,s.appointment_kind,s.starts_at,s.ends_at,s.capacity,count(a.appointment_id),s.state,s.version
		from platform.tenant t
		join org.organization o on o.tenant_id=t.tenant_id
		join org.public_location l on l.tenant_id=o.tenant_id and l.organization_id=o.organization_id and l.published=true
		join crm.appointment_slot s on s.tenant_id=o.tenant_id and s.organization_id=o.organization_id
		left join crm.appointment a on a.tenant_id=s.tenant_id and a.slot_id=s.slot_id and a.state in ('requested','confirmed')
		where t.tenant_code=$1 and t.status='active' and o.organization_code=$2 and o.status='active'
		  and s.appointment_kind=$3 and s.state='open' and s.starts_at >= $4 and s.starts_at < $5
		  and s.starts_at >= clock_timestamp()+make_interval(secs=>$6::double precision)
		  and exists(select 1 from crm.availability_entry e where e.tenant_id=s.tenant_id and e.organization_id=s.organization_id and e.resource_id is null and e.entry_type='working' and e.state='active' and e.starts_at<=s.starts_at and e.ends_at>=s.ends_at)
		  and not exists(select 1 from crm.availability_entry e where e.tenant_id=s.tenant_id and e.organization_id=s.organization_id and e.resource_id is null and e.entry_type='unavailable' and e.state='active' and tstzrange(e.starts_at,e.ends_at,'[)') && tstzrange(s.starts_at,s.ends_at,'[)'))
		group by s.tenant_id,s.slot_id
		having count(a.appointment_id) < s.capacity
		order by s.starts_at,s.slot_id limit 500`, tenantCode, organizationCode, kind, from, to, r.policy.LeadTime().Seconds())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []franchisejourney.AppointmentSlot{}
	for rows.Next() {
		var item franchisejourney.AppointmentSlot
		if err := rows.Scan(&item.ID, &item.OrganizationID, &item.Kind, &item.StartsAt, &item.EndsAt, &item.Capacity, &item.Booked, &item.State, &item.Version); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *FranchiseJourney) CreateAppointmentSlot(ctx context.Context, tenant string, value franchisejourney.AppointmentSlot, eventID string) (franchisejourney.AppointmentSlot, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return franchisejourney.AppointmentSlot{}, err
	}
	defer tx.Rollback(ctx)
	value, err = r.createAppointmentSlotTx(ctx, tx, tenant, "", value, eventID)
	if err != nil {
		return franchisejourney.AppointmentSlot{}, err
	}
	return value, tx.Commit(ctx)
}

func (r *FranchiseJourney) createAppointmentSlotTx(ctx context.Context, tx pgx.Tx, tenant, subject string, value franchisejourney.AppointmentSlot, eventID string) (franchisejourney.AppointmentSlot, error) {
	// Direct repository callers receive the same profile bounds as the service.
	// The persisted slot retains its own capacity; later profiles are prospective.
	var now time.Time
	if err := tx.QueryRow(ctx, `select clock_timestamp()`).Scan(&now); err != nil {
		return franchisejourney.AppointmentSlot{}, err
	}
	if !r.policy.AllowsSlot(value.StartsAt, value.EndsAt, now, value.Capacity) {
		return franchisejourney.AppointmentSlot{}, franchisejourney.ErrConflict
	}
	err := tx.QueryRow(ctx, `insert into crm.appointment_slot(tenant_id,slot_id,organization_id,appointment_kind,starts_at,ends_at,capacity,state,version)
		select $1,$2,o.organization_id,$4,$5,$6,$7,'open',1 from org.organization o
		where o.tenant_id=$1 and o.organization_id=$3 and o.status='active'
		returning slot_id,organization_id,appointment_kind,starts_at,ends_at,capacity,0,state,version`, tenant, value.ID, value.OrganizationID, value.Kind, value.StartsAt, value.EndsAt, value.Capacity).Scan(&value.ID, &value.OrganizationID, &value.Kind, &value.StartsAt, &value.EndsAt, &value.Capacity, &value.Booked, &value.State, &value.Version)
	if errors.Is(err, pgx.ErrNoRows) || postgresConflict(err) {
		return franchisejourney.AppointmentSlot{}, franchisejourney.ErrConflict
	}
	if err != nil {
		return franchisejourney.AppointmentSlot{}, err
	}
	payload := map[string]any{"organization_id": value.OrganizationID, "appointment_kind": value.Kind, "starts_at": value.StartsAt, "capacity": value.Capacity, "policy_sha256": r.policy.SHA256()}
	if subject != "" {
		payload["actor_subject"] = subject
	}
	if err = writeJourneyOutbox(ctx, tx, tenant, eventID, "appointment-slot", value.ID, 1, "appointment-slot.created", payload); err != nil {
		return franchisejourney.AppointmentSlot{}, err
	}
	return value, nil
}

func (r *FranchiseJourney) CreateServiceResource(ctx context.Context, tenant string, value franchisejourney.ServiceResource, eventID string) (franchisejourney.ServiceResource, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return franchisejourney.ServiceResource{}, err
	}
	defer tx.Rollback(ctx)
	value, err = r.createServiceResourceTx(ctx, tx, tenant, "", value, eventID)
	if err != nil {
		return franchisejourney.ServiceResource{}, err
	}
	return value, tx.Commit(ctx)
}

func (r *FranchiseJourney) createServiceResourceTx(ctx context.Context, tx pgx.Tx, tenant, subject string, value franchisejourney.ServiceResource, eventID string) (franchisejourney.ServiceResource, error) {
	result, err := tx.Exec(ctx, `insert into crm.service_resource(tenant_id,resource_id,organization_id,principal_subject,display_name,resource_kind,status,version)
		select $1,$2,o.organization_id,nullif($4,''),$5,$6,'active',1 from org.organization o where o.tenant_id=$1 and o.organization_id=$3 and o.status='active'`, tenant, value.ID, value.OrganizationID, value.PrincipalSubject, value.DisplayName, value.Kind)
	if err != nil {
		if postgresConflict(err) {
			return franchisejourney.ServiceResource{}, franchisejourney.ErrConflict
		}
		return franchisejourney.ServiceResource{}, err
	}
	if result.RowsAffected() != 1 {
		return franchisejourney.ServiceResource{}, franchisejourney.ErrConflict
	}
	if _, err = tx.Exec(ctx, `insert into crm.resource_skill(tenant_id,resource_id,appointment_kind) select $1,$2,unnest($3::text[])`, tenant, value.ID, value.Skills); err != nil {
		if postgresConflict(err) {
			return franchisejourney.ServiceResource{}, franchisejourney.ErrConflict
		}
		return franchisejourney.ServiceResource{}, err
	}
	payload := map[string]any{"organization_id": value.OrganizationID, "resource_kind": value.Kind, "skills": value.Skills}
	if subject != "" {
		payload["actor_subject"] = subject
	}
	if err = writeJourneyOutbox(ctx, tx, tenant, eventID, "service-resource", value.ID, 1, "service-resource.created", payload); err != nil {
		return franchisejourney.ServiceResource{}, err
	}
	return value, nil
}

func (r *FranchiseJourney) AssignAppointmentResource(ctx context.Context, tenant, organization, appointment, resource string, version int64, eventID string) (franchisejourney.Appointment, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return franchisejourney.Appointment{}, err
	}
	defer tx.Rollback(ctx)
	if _, err = tx.Exec(ctx, `insert into crm.appointment_resource(tenant_id,appointment_id,resource_id) values($1,$2,$3)`, tenant, appointment, resource); err != nil {
		if postgresConflict(err) {
			return franchisejourney.Appointment{}, franchisejourney.ErrConflict
		}
		return franchisejourney.Appointment{}, err
	}
	var value franchisejourney.Appointment
	err = tx.QueryRow(ctx, `update crm.appointment set version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and organization_id=$2 and appointment_id=$3 and state='requested' and version=$4 returning appointment_id,organization_id,lead_id,coalesce(model_id,''),appointment_kind,starts_at,state,version,coalesce(slot_id,''),coalesce(ends_at,starts_at)`, tenant, organization, appointment, version).Scan(&value.ID, &value.OrganizationID, &value.LeadID, &value.ModelID, &value.Kind, &value.StartsAt, &value.State, &value.Version, &value.SlotID, &value.EndsAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return franchisejourney.Appointment{}, franchisejourney.ErrConflict
	}
	if err != nil {
		return franchisejourney.Appointment{}, err
	}
	value.ResourceID = resource
	if err = writeJourneyOutbox(ctx, tx, tenant, eventID, "appointment", appointment, value.Version, "appointment.resource-assigned", map[string]any{"organization_id": organization, "resource_id": resource}); err != nil {
		return franchisejourney.Appointment{}, err
	}
	return value, tx.Commit(ctx)
}

func (r *FranchiseJourney) TransitionAppointment(ctx context.Context, tenant, organization, appointment, current, target string, version int64, subject, reason, eventID string) (franchisejourney.Appointment, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return franchisejourney.Appointment{}, err
	}
	defer tx.Rollback(ctx)
	var value franchisejourney.Appointment
	err = tx.QueryRow(ctx, `update crm.appointment a set state=$5,version=version+1,updated_at=clock_timestamp()
		where tenant_id=$1 and organization_id=$2 and appointment_id=$3 and state=$4 and version=$6
		and ($5<>'confirmed' or exists(select 1 from crm.appointment_resource ar where ar.tenant_id=a.tenant_id and ar.appointment_id=a.appointment_id))
		and ($5 not in ('completed','no-show') or a.starts_at<=clock_timestamp())
		returning appointment_id,organization_id,lead_id,coalesce(model_id,''),appointment_kind,starts_at,state,version,coalesce(slot_id,''),coalesce(ends_at,starts_at),coalesce((select ar.resource_id from crm.appointment_resource ar where ar.tenant_id=a.tenant_id and ar.appointment_id=a.appointment_id),'')`, tenant, organization, appointment, current, target, version).Scan(&value.ID, &value.OrganizationID, &value.LeadID, &value.ModelID, &value.Kind, &value.StartsAt, &value.State, &value.Version, &value.SlotID, &value.EndsAt, &value.ResourceID)
	if errors.Is(err, pgx.ErrNoRows) {
		return franchisejourney.Appointment{}, franchisejourney.ErrConflict
	}
	if err != nil {
		return franchisejourney.Appointment{}, err
	}
	if _, err = tx.Exec(ctx, `insert into crm.appointment_transition(tenant_id,transition_id,appointment_id,from_state,to_state,actor_subject,reason_code) values($1,$2,$3,$4,$5,$6,nullif($7,''))`, tenant, eventID, appointment, current, target, subject, reason); err != nil {
		return franchisejourney.Appointment{}, err
	}
	if err = writeJourneyOutbox(ctx, tx, tenant, eventID, "appointment", appointment, value.Version, "appointment."+target, map[string]any{"organization_id": organization, "resource_id": value.ResourceID, "actor_subject": subject, "reason_code": reason}); err != nil {
		return franchisejourney.Appointment{}, err
	}
	return value, tx.Commit(ctx)
}

func (r *FranchiseJourney) RequestAppointment(ctx context.Context, tenantCode, organizationCode, idempotencyKey string, value franchisejourney.Appointment, requestHash, eventID string) (franchisejourney.Appointment, bool, error) {
	boundHash, bindErr := r.policy.BindRequestHash(requestHash)
	if bindErr != nil {
		return franchisejourney.Appointment{}, false, franchisejourney.ErrInvalid
	}
	requestHash = boundHash
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return franchisejourney.Appointment{}, false, err
	}
	defer tx.Rollback(ctx)
	var tenant, organization string
	err = tx.QueryRow(ctx, `select t.tenant_id,o.organization_id from platform.tenant t join org.organization o on o.tenant_id=t.tenant_id join org.public_location l on l.tenant_id=o.tenant_id and l.organization_id=o.organization_id where t.tenant_code=$1 and t.status='active' and o.organization_code=$2 and o.status='active' and l.published=true`, tenantCode, organizationCode).Scan(&tenant, &organization)
	if errors.Is(err, pgx.ErrNoRows) {
		return franchisejourney.Appointment{}, false, franchisejourney.ErrNotFound
	}
	if err != nil {
		return franchisejourney.Appointment{}, false, err
	}
	result, err := tx.Exec(ctx, `insert into platform.idempotency_record(tenant_id,scope,idempotency_key,request_sha256_hex,status,locked_until,expires_at) values($1,'public-appointment',$2,$3,'processing',clock_timestamp()+interval '30 seconds',clock_timestamp()+interval '24 hours') on conflict do nothing`, tenant, idempotencyKey, requestHash)
	if err != nil {
		return franchisejourney.Appointment{}, false, err
	}
	if result.RowsAffected() == 0 {
		var storedHash, status, resourceID string
		if err = tx.QueryRow(ctx, `select request_sha256_hex,status,coalesce(resource_id,'') from platform.idempotency_record where tenant_id=$1 and scope='public-appointment' and idempotency_key=$2`, tenant, idempotencyKey).Scan(&storedHash, &status, &resourceID); err != nil {
			return franchisejourney.Appointment{}, false, err
		}
		if storedHash != requestHash || status != "completed" || resourceID == "" {
			return franchisejourney.Appointment{}, false, franchisejourney.ErrConflict
		}
		var replay franchisejourney.Appointment
		err = tx.QueryRow(ctx, `select appointment_id,organization_id,lead_id,coalesce(model_id,''),appointment_kind,starts_at,state,version,coalesce(slot_id,''),coalesce(ends_at,starts_at) from crm.appointment where tenant_id=$1 and appointment_id=$2`, tenant, resourceID).Scan(&replay.ID, &replay.OrganizationID, &replay.LeadID, &replay.ModelID, &replay.Kind, &replay.StartsAt, &replay.State, &replay.Version, &replay.SlotID, &replay.EndsAt)
		if err != nil {
			return franchisejourney.Appointment{}, false, err
		}
		if err = tx.Commit(ctx); err != nil {
			return franchisejourney.Appointment{}, false, err
		}
		return replay, true, nil
	}
	value.OrganizationID = organization
	if err = lockOrganizationSchedule(ctx, tx, tenant, organization); err != nil {
		return franchisejourney.Appointment{}, false, err
	}
	var capacity int
	err = tx.QueryRow(ctx, `select s.slot_id,s.ends_at,s.capacity from crm.appointment_slot s where s.tenant_id=$1 and s.organization_id=$2 and s.appointment_kind=$3 and s.starts_at=$4 and s.state='open' and s.starts_at>=clock_timestamp()+make_interval(secs=>$5::double precision)
		and exists(select 1 from crm.availability_entry e where e.tenant_id=s.tenant_id and e.organization_id=s.organization_id and e.resource_id is null and e.entry_type='working' and e.state='active' and e.starts_at<=s.starts_at and e.ends_at>=s.ends_at)
		and not exists(select 1 from crm.availability_entry e where e.tenant_id=s.tenant_id and e.organization_id=s.organization_id and e.resource_id is null and e.entry_type='unavailable' and e.state='active' and tstzrange(e.starts_at,e.ends_at,'[)') && tstzrange(s.starts_at,s.ends_at,'[)')) for update of s`, tenant, organization, value.Kind, value.StartsAt, r.policy.LeadTime().Seconds()).Scan(&value.SlotID, &value.EndsAt, &capacity)
	if errors.Is(err, pgx.ErrNoRows) {
		return franchisejourney.Appointment{}, false, franchisejourney.ErrConflict
	}
	if err != nil {
		return franchisejourney.Appointment{}, false, err
	}
	var booked int
	if err = tx.QueryRow(ctx, `select count(*) from crm.appointment where tenant_id=$1 and slot_id=$2 and state in ('requested','confirmed')`, tenant, value.SlotID).Scan(&booked); err != nil {
		return franchisejourney.Appointment{}, false, err
	}
	if booked >= capacity {
		return franchisejourney.Appointment{}, false, franchisejourney.ErrConflict
	}
	result, err = tx.Exec(ctx, `insert into crm.appointment(tenant_id,appointment_id,organization_id,lead_id,customer_principal_id,model_id,appointment_kind,starts_at,ends_at,slot_id,state,version) select $1,$2,$3,l.lead_id,l.customer_principal_id,nullif($5,''),$6,$7,$8,$9,'requested',1 from crm.lead l where l.tenant_id=$1 and l.organization_id=$3 and l.lead_id=$4`, tenant, value.ID, organization, value.LeadID, value.ModelID, value.Kind, value.StartsAt, value.EndsAt, value.SlotID)
	if err != nil {
		return franchisejourney.Appointment{}, false, err
	}
	if result.RowsAffected() != 1 {
		return franchisejourney.Appointment{}, false, franchisejourney.ErrConflict
	}
	if err = writeJourneyOutbox(ctx, tx, tenant, eventID, "appointment", value.ID, 1, "appointment.requested", map[string]any{"organization_id": organization, "lead_id": value.LeadID, "policy_sha256": r.policy.SHA256()}); err != nil {
		return franchisejourney.Appointment{}, false, err
	}
	_, err = tx.Exec(ctx, `update platform.idempotency_record set status='completed',response_code=202,response_body=jsonb_build_object('appointment_id',$3::text),resource_type='appointment',resource_id=$3,locked_until=null where tenant_id=$1 and scope='public-appointment' and idempotency_key=$2 and status='processing'`, tenant, idempotencyKey, value.ID)
	if err != nil {
		return franchisejourney.Appointment{}, false, err
	}
	if err = tx.Commit(ctx); err != nil {
		return franchisejourney.Appointment{}, false, err
	}
	return value, false, nil
}

func (r *FranchiseJourney) Leads(ctx context.Context, tenant, organization string, limit int, after string) (franchisejourney.Page[franchisejourney.Lead], error) {
	rows, err := r.pool.Query(ctx, `select lead_id,organization_id,coalesce(model_id,''),lifecycle_state,source_code,coalesce(assigned_subject,''),created_at,version from crm.lead where tenant_id=$1 and organization_id=$2 and ($3='' or lead_id>$3) order by lead_id limit $4`, tenant, organization, after, limit+1)
	if err != nil {
		return franchisejourney.Page[franchisejourney.Lead]{}, err
	}
	defer rows.Close()
	items := []franchisejourney.Lead{}
	for rows.Next() {
		var item franchisejourney.Lead
		if err := rows.Scan(&item.ID, &item.OrganizationID, &item.ModelID, &item.State, &item.SourceCode, &item.AssignedSubject, &item.CreatedAt, &item.Version); err != nil {
			return franchisejourney.Page[franchisejourney.Lead]{}, err
		}
		items = append(items, item)
	}
	page := franchisejourney.Page[franchisejourney.Lead]{Items: items}
	if len(items) > limit {
		page.NextCursor = items[limit-1].ID
		page.Items = items[:limit]
	}
	return page, rows.Err()
}

func (r *FranchiseJourney) AssignLead(ctx context.Context, tenant, organization, lead, subject string, version int64, eventID string) (franchisejourney.Lead, error) {
	return r.assignLead(ctx, tenant, organization, lead, subject, version, eventID, "")
}

func (r *FranchiseJourney) AssignLeadAs(ctx context.Context, tenant, organization, lead, subject string, version int64, eventID, actor string) (franchisejourney.Lead, error) {
	if actor == "" {
		return franchisejourney.Lead{}, franchisejourney.ErrInvalid
	}
	return r.assignLead(ctx, tenant, organization, lead, subject, version, eventID, actor)
}

func (r *FranchiseJourney) assignLead(ctx context.Context, tenant, organization, lead, subject string, version int64, eventID, actor string) (franchisejourney.Lead, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return franchisejourney.Lead{}, err
	}
	defer tx.Rollback(ctx)
	var value franchisejourney.Lead
	err = tx.QueryRow(ctx, `update crm.lead set assigned_subject=$5,version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and organization_id=$2 and lead_id=$3 and version=$4 and lifecycle_state not in ('converted','lost') returning lead_id,organization_id,coalesce(model_id,''),lifecycle_state,source_code,assigned_subject,created_at,version`, tenant, organization, lead, version, subject).Scan(&value.ID, &value.OrganizationID, &value.ModelID, &value.State, &value.SourceCode, &value.AssignedSubject, &value.CreatedAt, &value.Version)
	if errors.Is(err, pgx.ErrNoRows) {
		return franchisejourney.Lead{}, franchisejourney.ErrConflict
	}
	if err != nil {
		return franchisejourney.Lead{}, err
	}
	payload := map[string]any{"organization_id": organization, "assigned_subject": subject}
	if actor != "" {
		payload["actor_subject"] = actor
	}
	if err = writeJourneyOutbox(ctx, tx, tenant, eventID, "lead", lead, value.Version, "lead.assigned", payload); err != nil {
		return franchisejourney.Lead{}, err
	}
	return value, tx.Commit(ctx)
}

func (r *FranchiseJourney) TransitionLead(ctx context.Context, tenant, organization, lead, current, target string, version int64, eventID string) (franchisejourney.Lead, error) {
	return r.transitionLead(ctx, tenant, organization, lead, current, target, version, eventID, "")
}

func (r *FranchiseJourney) TransitionLeadAs(ctx context.Context, tenant, organization, lead, current, target string, version int64, eventID, actor string) (franchisejourney.Lead, error) {
	if actor == "" {
		return franchisejourney.Lead{}, franchisejourney.ErrInvalid
	}
	return r.transitionLead(ctx, tenant, organization, lead, current, target, version, eventID, actor)
}

func (r *FranchiseJourney) transitionLead(ctx context.Context, tenant, organization, lead, current, target string, version int64, eventID, actor string) (franchisejourney.Lead, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return franchisejourney.Lead{}, err
	}
	defer tx.Rollback(ctx)
	var value franchisejourney.Lead
	err = tx.QueryRow(ctx, `update crm.lead set lifecycle_state=$6,version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and organization_id=$2 and lead_id=$3 and lifecycle_state=$4 and version=$5 returning lead_id,organization_id,coalesce(model_id,''),lifecycle_state,source_code,coalesce(assigned_subject,''),created_at,version`, tenant, organization, lead, current, version, target).Scan(&value.ID, &value.OrganizationID, &value.ModelID, &value.State, &value.SourceCode, &value.AssignedSubject, &value.CreatedAt, &value.Version)
	if errors.Is(err, pgx.ErrNoRows) {
		return franchisejourney.Lead{}, franchisejourney.ErrConflict
	}
	if err != nil {
		return franchisejourney.Lead{}, err
	}
	payload := map[string]any{"organization_id": organization}
	if actor != "" {
		payload["actor_subject"] = actor
	}
	if err = writeJourneyOutbox(ctx, tx, tenant, eventID, "lead", lead, value.Version, "lead."+target, payload); err != nil {
		return franchisejourney.Lead{}, err
	}
	return value, tx.Commit(ctx)
}

func (r *FranchiseJourney) CreateQuote(ctx context.Context, tenant, idempotencyKey string, value franchisejourney.Quote, requestHash, eventID string) (franchisejourney.Quote, bool, error) {
	return r.createQuote(ctx, tenant, idempotencyKey, value, requestHash, eventID, "")
}
func (r *FranchiseJourney) CreateQuoteAs(ctx context.Context, tenant, key string, value franchisejourney.Quote, hash, event, actor string) (franchisejourney.Quote, bool, error) {
	if actor == "" {
		return franchisejourney.Quote{}, false, franchisejourney.ErrInvalid
	}
	return r.createQuote(ctx, tenant, key, value, hash, event, actor)
}
func (r *FranchiseJourney) createQuote(ctx context.Context, tenant, idempotencyKey string, value franchisejourney.Quote, requestHash, eventID, actor string) (franchisejourney.Quote, bool, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return franchisejourney.Quote{}, false, err
	}
	defer tx.Rollback(ctx)
	result, err := tx.Exec(ctx, `insert into platform.idempotency_record(tenant_id,scope,idempotency_key,request_sha256_hex,status,locked_until,expires_at) values($1,'franchise-quote',$2,$3,'processing',clock_timestamp()+interval '30 seconds',clock_timestamp()+interval '24 hours') on conflict do nothing`, tenant, idempotencyKey, requestHash)
	if err != nil {
		return franchisejourney.Quote{}, false, err
	}
	if result.RowsAffected() == 0 {
		var storedHash, status, resourceID string
		if err = tx.QueryRow(ctx, `select request_sha256_hex,status,coalesce(resource_id,'') from platform.idempotency_record where tenant_id=$1 and scope='franchise-quote' and idempotency_key=$2`, tenant, idempotencyKey).Scan(&storedHash, &status, &resourceID); err != nil {
			return franchisejourney.Quote{}, false, err
		}
		if storedHash != requestHash || status != "completed" || resourceID == "" {
			return franchisejourney.Quote{}, false, franchisejourney.ErrConflict
		}
		var replay franchisejourney.Quote
		err = tx.QueryRow(ctx, `select quotation_id,organization_id,lead_id,coalesce(customer_principal_id,''),variant_id,price_book_id,currency,total_minor_units,valid_until,state,version,coalesce(order_id,'') from sales.quotation where tenant_id=$1 and quotation_id=$2`, tenant, resourceID).Scan(&replay.ID, &replay.OrganizationID, &replay.LeadID, &replay.CustomerSubject, &replay.VariantID, &replay.PriceBookID, &replay.Currency, &replay.TotalMinorUnits, &replay.ValidUntil, &replay.State, &replay.Version, &replay.OrderID)
		if err != nil {
			return franchisejourney.Quote{}, false, err
		}
		if err = tx.Commit(ctx); err != nil {
			return franchisejourney.Quote{}, false, err
		}
		return replay, true, nil
	}
	err = tx.QueryRow(ctx, `select b.currency,e.amount_minor_units from pricing.price_book b join pricing.price_book_entry e on e.tenant_id=b.tenant_id and e.price_book_id=b.price_book_id `+bcPriceTimeContextSQL+` where b.tenant_id=$1 and b.price_book_id=$2 and e.variant_id=$3 and `+bcPriceEligibilitySQL+``, tenant, value.PriceBookID, value.VariantID).Scan(&value.Currency, &value.TotalMinorUnits)
	if errors.Is(err, pgx.ErrNoRows) {
		return franchisejourney.Quote{}, false, franchisejourney.ErrConflict
	}
	if err != nil {
		return franchisejourney.Quote{}, false, err
	}
	result, err = tx.Exec(ctx, `insert into sales.quotation(tenant_id,quotation_id,organization_id,lead_id,customer_principal_id,variant_id,price_book_id,currency,total_minor_units,valid_until,state,version) select $1,$2,$3,$4,customer_principal_id,$5,$6,$7,$8,$9,'issued',1 from crm.lead where tenant_id=$1 and organization_id=$3 and lead_id=$4`, tenant, value.ID, value.OrganizationID, value.LeadID, value.VariantID, value.PriceBookID, value.Currency, value.TotalMinorUnits, value.ValidUntil)
	if err != nil {
		return franchisejourney.Quote{}, false, err
	}
	if result.RowsAffected() != 1 {
		return franchisejourney.Quote{}, false, franchisejourney.ErrConflict
	}
	if err = writeJourneyOutbox(ctx, tx, tenant, eventID, "quotation", value.ID, 1, "quotation.issued", map[string]any{"organization_id": value.OrganizationID, "lead_id": value.LeadID, "actor_subject": actor}); err != nil {
		return franchisejourney.Quote{}, false, err
	}
	_, err = tx.Exec(ctx, `update platform.idempotency_record set status='completed',response_code=201,response_body=jsonb_build_object('quotation_id',$3::text),resource_type='quotation',resource_id=$3,locked_until=null where tenant_id=$1 and scope='franchise-quote' and idempotency_key=$2 and status='processing'`, tenant, idempotencyKey, value.ID)
	if err != nil {
		return franchisejourney.Quote{}, false, err
	}
	if err = tx.Commit(ctx); err != nil {
		return franchisejourney.Quote{}, false, err
	}
	return value, false, nil
}

func (r *FranchiseJourney) AcceptQuote(ctx context.Context, tenant, organization, customer, quote string, version int64, evidence, orderID, lineID, quoteEventID, orderEventID string) (franchisejourney.Quote, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return franchisejourney.Quote{}, err
	}
	defer tx.Rollback(ctx)
	var value franchisejourney.Quote
	err = tx.QueryRow(ctx, `select quotation_id,organization_id,lead_id,customer_principal_id,variant_id,price_book_id,currency,total_minor_units,valid_until,state,version from sales.quotation where tenant_id=$1 and organization_id=$2 and customer_principal_id=$3 and quotation_id=$4 and version=$5 and state='issued' and valid_until>clock_timestamp() for update`, tenant, organization, customer, quote, version).Scan(&value.ID, &value.OrganizationID, &value.LeadID, &value.CustomerSubject, &value.VariantID, &value.PriceBookID, &value.Currency, &value.TotalMinorUnits, &value.ValidUntil, &value.State, &value.Version)
	if errors.Is(err, pgx.ErrNoRows) {
		return franchisejourney.Quote{}, franchisejourney.ErrConflict
	}
	if err != nil {
		return franchisejourney.Quote{}, err
	}
	// AUTHORED mapping from the already authorized one-line quote contract.
	// BC supplies document/line copying; consent, states and retention stay here.
	orderHeader, orderLines, err := bcsales.TransferQuoteToOrder(
		bcsales.Header{DocumentType: bcsales.Quote, ID: value.ID, CustomerID: value.CustomerSubject, Currency: value.Currency},
		[]bcsales.Line{{DocumentType: bcsales.Quote, DocumentID: value.ID, VariantID: value.VariantID, Quantity: 1, UnitPriceMinor: value.TotalMinorUnits}}, orderID)
	if err != nil || len(orderLines) != 1 {
		return franchisejourney.Quote{}, franchisejourney.ErrConflict
	}
	orderLine := orderLines[0]
	if _, err = tx.Exec(ctx, `insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version) values($1,$2,$3,$4,'placed',$5,$6,1)`, tenant, orderHeader.ID, organization, orderHeader.CustomerID, orderHeader.Currency, value.TotalMinorUnits); err != nil {
		return franchisejourney.Quote{}, err
	}
	if _, err = tx.Exec(ctx, `insert into sales.customer_order_line(tenant_id,order_id,line_id,variant_id,quantity,unit_price_minor_units) values($1,$2,$3,$4,$5,$6)`, tenant, orderLine.DocumentID, lineID, orderLine.VariantID, orderLine.Quantity, orderLine.UnitPriceMinor); err != nil {
		return franchisejourney.Quote{}, err
	}
	if _, err = tx.Exec(ctx, `update sales.quotation set state='accepted',order_id=$6,version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and organization_id=$2 and customer_principal_id=$3 and quotation_id=$4 and version=$5`, tenant, organization, customer, quote, version, orderID); err != nil {
		return franchisejourney.Quote{}, err
	}
	if _, err = tx.Exec(ctx, `insert into sales.quotation_acceptance(tenant_id,quotation_id,order_id,customer_principal_id,evidence_sha256_hex) values($1,$2,$3,$4,$5)`, tenant, quote, orderID, customer, evidence); err != nil {
		return franchisejourney.Quote{}, err
	}
	value.State = "accepted"
	value.Version++
	value.OrderID = orderID
	if err = writeJourneyOutbox(ctx, tx, tenant, quoteEventID, "quotation", quote, value.Version, "quotation.accepted", map[string]any{"organization_id": organization, "customer_subject": customer, "order_id": orderID}); err != nil {
		return franchisejourney.Quote{}, err
	}
	if err = writeJourneyOutbox(ctx, tx, tenant, orderEventID, "customer-order", orderID, 1, "customer-order.placed", map[string]any{"organization_id": organization, "customer_subject": customer, "quotation_id": quote, "currency": value.Currency, "total_minor_units": value.TotalMinorUnits}); err != nil {
		return franchisejourney.Quote{}, err
	}
	return value, tx.Commit(ctx)
}

func (r *FranchiseJourney) PublishDeliveryChecklist(ctx context.Context, tenant, subject string, value franchisejourney.DeliveryChecklist, eventID string) (franchisejourney.DeliveryChecklist, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return franchisejourney.DeliveryChecklist{}, err
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, `insert into sales.delivery_checklist_template(tenant_id,organization_id,checklist_id,checklist_version,title,state,created_by_subject) values($1,$2,$3,$4,$5,'draft',$6)`, tenant, value.OrganizationID, value.ID, value.Version, value.Title, subject)
	if postgresConflict(err) {
		return franchisejourney.DeliveryChecklist{}, franchisejourney.ErrConflict
	}
	if err != nil {
		return franchisejourney.DeliveryChecklist{}, err
	}
	for _, item := range value.Items {
		_, err = tx.Exec(ctx, `insert into sales.delivery_checklist_item(tenant_id,organization_id,checklist_id,checklist_version,item_id,ordinal,prompt,response_type,required) values($1,$2,$3,$4,$5,$6,$7,$8,$9)`, tenant, value.OrganizationID, value.ID, value.Version, item.ID, item.Ordinal, item.Prompt, item.ResponseType, item.Required)
		if postgresConflict(err) {
			return franchisejourney.DeliveryChecklist{}, franchisejourney.ErrConflict
		}
		if err != nil {
			return franchisejourney.DeliveryChecklist{}, err
		}
	}
	if err = tx.QueryRow(ctx, `update sales.delivery_checklist_template set state='published',published_at=clock_timestamp() where tenant_id=$1 and organization_id=$2 and checklist_id=$3 and checklist_version=$4 and state='draft' returning state`, tenant, value.OrganizationID, value.ID, value.Version).Scan(&value.State); err != nil {
		return franchisejourney.DeliveryChecklist{}, err
	}
	if err = writeJourneyOutbox(ctx, tx, tenant, eventID, "delivery-checklist", value.ID, value.Version, "delivery-checklist.published", map[string]any{"organization_id": value.OrganizationID, "checklist_version": value.Version, "item_count": len(value.Items), "actor_subject": subject}); err != nil {
		return franchisejourney.DeliveryChecklist{}, err
	}
	return value, tx.Commit(ctx)
}

func (r *FranchiseJourney) CompleteDeliveryChecklist(ctx context.Context, tenant, organization, subject, handover string, version int64, checklistID string, checklistVersion int64, responses []franchisejourney.ChecklistResponse, eventID string) (franchisejourney.Handover, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return franchisejourney.Handover{}, err
	}
	defer tx.Rollback(ctx)
	if err = lockDeliveryScope(ctx, tx, tenant, organization, handover); err != nil {
		return franchisejourney.Handover{}, err
	}
	for _, response := range responses {
		var evidence any
		if response.EvidenceSHA256 != "" {
			evidence = response.EvidenceSHA256
		}
		_, err = tx.Exec(ctx, `insert into sales.delivery_checklist_response(tenant_id,organization_id,handover_id,checklist_id,checklist_version,item_id,response_text,evidence_sha256_hex,answered_by_subject) values($1,$2,$3,$4,$5,$6,$7,$8,$9)`, tenant, organization, handover, checklistID, checklistVersion, response.ItemID, response.ResponseText, evidence, subject)
		if postgresConflict(err) {
			return franchisejourney.Handover{}, franchisejourney.ErrConflict
		}
		if err != nil {
			return franchisejourney.Handover{}, err
		}
	}
	var value franchisejourney.Handover
	value.ChecklistItems = []franchisejourney.ChecklistItem{}
	err = tx.QueryRow(ctx, `update sales.delivery_handover h set state='presented',checklist_id=$5,checklist_version=$6,checklist_completed_at=clock_timestamp(),checklist_completed_by_subject=$7,updated_at=clock_timestamp(),version=h.version+1 from sales.delivery_checklist_template t where h.tenant_id=$1 and h.organization_id=$2 and h.handover_id=$3 and h.version=$4 and h.state='prepared' and t.tenant_id=h.tenant_id and t.organization_id=h.organization_id and t.checklist_id=$5 and t.checklist_version=$6 and t.state='published' returning h.handover_id,h.organization_id,h.order_id,h.customer_principal_id,h.stock_unit_id,h.state,h.version,h.checklist_id,h.checklist_version,t.title,h.checklist_completed_at`, tenant, organization, handover, version, checklistID, checklistVersion, subject).Scan(&value.ID, &value.OrganizationID, &value.OrderID, &value.CustomerSubject, &value.StockUnitID, &value.State, &value.Version, &value.ChecklistID, &value.ChecklistVersion, &value.ChecklistTitle, &value.ChecklistCompletedAt)
	if errors.Is(err, pgx.ErrNoRows) || postgresConflict(err) {
		return franchisejourney.Handover{}, franchisejourney.ErrConflict
	}
	if err != nil {
		return franchisejourney.Handover{}, err
	}
	if err = writeJourneyOutbox(ctx, tx, tenant, eventID, "delivery-handover", handover, value.Version, "delivery-handover.checklist-completed", map[string]any{"organization_id": organization, "checklist_id": checklistID, "checklist_version": checklistVersion, "response_count": len(responses), "actor_subject": subject}); err != nil {
		return franchisejourney.Handover{}, err
	}
	return value, tx.Commit(ctx)
}

func (r *FranchiseJourney) RejectHandover(ctx context.Context, tenant, organization, customer, handover string, version int64, reasonCode, details, evidence, exceptionID, eventID string) (franchisejourney.DeliveryException, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return franchisejourney.DeliveryException{}, err
	}
	defer tx.Rollback(ctx)
	if err = lockDeliveryScope(ctx, tx, tenant, organization, handover); err != nil {
		return franchisejourney.DeliveryException{}, err
	}
	var nextVersion int64
	err = tx.QueryRow(ctx, `update sales.delivery_handover set state='rejected',updated_at=clock_timestamp(),version=version+1 where tenant_id=$1 and organization_id=$2 and customer_principal_id=$3 and handover_id=$4 and version=$5 and state='presented' and checklist_completed_at is not null returning version`, tenant, organization, customer, handover, version).Scan(&nextVersion)
	if errors.Is(err, pgx.ErrNoRows) {
		return franchisejourney.DeliveryException{}, franchisejourney.ErrConflict
	}
	if err != nil {
		return franchisejourney.DeliveryException{}, err
	}
	value := franchisejourney.DeliveryException{ID: exceptionID, OrganizationID: organization, HandoverID: handover, CustomerSubject: customer, ReasonCode: reasonCode, Details: details, State: "open", Version: 1}
	err = tx.QueryRow(ctx, `insert into sales.delivery_exception(tenant_id,exception_id,organization_id,handover_id,customer_principal_id,reason_code,details,rejection_evidence_sha256_hex,state,version) values($1,$2,$3,$4,$5,$6,$7,$8,'open',1) returning created_at`, tenant, exceptionID, organization, handover, customer, reasonCode, details, evidence).Scan(&value.CreatedAt)
	if postgresConflict(err) {
		return franchisejourney.DeliveryException{}, franchisejourney.ErrConflict
	}
	if err != nil {
		return franchisejourney.DeliveryException{}, err
	}
	if err = writeJourneyOutbox(ctx, tx, tenant, eventID, "delivery-handover", handover, nextVersion, "delivery-handover.rejected", map[string]any{"organization_id": organization, "customer_subject": customer, "exception_id": exceptionID, "reason_code": reasonCode, "rejection_evidence_sha256": evidence}); err != nil {
		return franchisejourney.DeliveryException{}, err
	}
	return value, tx.Commit(ctx)
}

func (r *FranchiseJourney) DeliveryExceptions(ctx context.Context, tenant, organization string, limit int) ([]franchisejourney.DeliveryException, error) {
	rows, err := r.pool.Query(ctx, `select e.exception_id,e.organization_id,e.handover_id,e.customer_principal_id,e.reason_code,e.details,e.state,e.version,e.created_at,e.resolved_at,coalesce(x.action,''),coalesce(x.successor_handover_id,''),coalesce(x.return_authorization_id,''),exists(select 1 from sales.delivery_handover h join sales.customer_order o on o.tenant_id=h.tenant_id and o.order_id=h.order_id and o.organization_id=h.organization_id and o.customer_principal_id=h.customer_principal_id join inventory.stock_unit su on su.tenant_id=h.tenant_id and su.stock_unit_id=h.stock_unit_id and su.organization_id=h.organization_id where h.tenant_id=e.tenant_id and h.handover_id=e.handover_id and h.organization_id=e.organization_id and h.customer_principal_id=e.customer_principal_id) from sales.delivery_exception e left join sales.delivery_exception_resolution x on x.tenant_id=e.tenant_id and x.exception_id=e.exception_id where e.tenant_id=$1 and e.organization_id=$2 order by e.created_at desc,e.exception_id limit $3`, tenant, organization, limit)
	if err != nil {
		return nil, err
	}
	return scanDeliveryExceptions(rows)
}

func scanDeliveryExceptions(rows pgx.Rows) ([]franchisejourney.DeliveryException, error) {
	defer rows.Close()
	items := []franchisejourney.DeliveryException{}
	for rows.Next() {
		var item franchisejourney.DeliveryException
		var coherent bool
		if err := rows.Scan(&item.ID, &item.OrganizationID, &item.HandoverID, &item.CustomerSubject, &item.ReasonCode, &item.Details, &item.State, &item.Version, &item.CreatedAt, &item.ResolvedAt, &item.ResolutionAction, &item.SuccessorHandoverID, &item.ReturnAuthorizationID, &coherent); err != nil {
			return nil, err
		}
		if !coherent {
			return nil, franchisejourney.ErrConflict
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *FranchiseJourney) ResolveDeliveryException(ctx context.Context, tenant, organization, subject, exceptionID string, version int64, action, notes, successorID, authorizationID, resolutionID, eventID string) (franchisejourney.DeliveryResolution, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return franchisejourney.DeliveryResolution{}, err
	}
	defer tx.Rollback(ctx)
	var handover franchisejourney.Handover
	var customer, orderID, stockUnitID string
	if err = tx.QueryRow(ctx, `select h.handover_id,h.customer_principal_id,h.order_id,h.stock_unit_id from sales.delivery_exception e join sales.delivery_handover h on h.tenant_id=e.tenant_id and h.handover_id=e.handover_id where e.tenant_id=$1 and e.organization_id=$2 and e.exception_id=$3 and e.state='open' and e.version=$4 and h.state='rejected' for update of e,h`, tenant, organization, exceptionID, version).Scan(&handover.ID, &customer, &orderID, &stockUnitID); errors.Is(err, pgx.ErrNoRows) {
		return franchisejourney.DeliveryResolution{}, franchisejourney.ErrConflict
	}
	if err != nil {
		return franchisejourney.DeliveryResolution{}, err
	}
	result := franchisejourney.DeliveryResolution{}
	if err = lockDeliveryScope(ctx, tx, tenant, organization, handover.ID); err != nil {
		return franchisejourney.DeliveryResolution{}, err
	}
	originalHandoverID := handover.ID
	var successor any
	var authorization any
	if action == "correct-and-represent" {
		successor = successorID
		handover = franchisejourney.Handover{ID: successorID, OrganizationID: organization, OrderID: orderID, CustomerSubject: customer, StockUnitID: stockUnitID, State: "prepared", Version: 1, ChecklistItems: []franchisejourney.ChecklistItem{}, SupersedesHandoverID: originalHandoverID}
		_, err = tx.Exec(ctx, `insert into sales.delivery_handover(tenant_id,handover_id,organization_id,order_id,customer_principal_id,stock_unit_id,state,version,supersedes_handover_id) values($1,$2,$3,$4,$5,$6,'prepared',1,$7)`, tenant, successorID, organization, orderID, customer, stockUnitID, originalHandoverID)
		result.SuccessorHandover = &handover
	} else {
		authorization = authorizationID
		_, err = tx.Exec(ctx, `insert into sales.return_authorization(tenant_id,authorization_id,organization_id,exception_id,handover_id,order_id,stock_unit_id,customer_principal_id,disposition,exact_cost_source_order_id,state,authorized_by_subject) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$6,'authorized',$10)`, tenant, authorizationID, organization, exceptionID, originalHandoverID, orderID, stockUnitID, customer, action, subject)
		result.ReturnAuthorizationID = authorizationID
		result.Disposition = action
	}
	if postgresConflict(err) {
		return franchisejourney.DeliveryResolution{}, franchisejourney.ErrConflict
	}
	if err != nil {
		return franchisejourney.DeliveryResolution{}, err
	}
	_, err = tx.Exec(ctx, `insert into sales.delivery_exception_resolution(tenant_id,resolution_id,exception_id,action,notes,resolved_by_subject,successor_handover_id,return_authorization_id) values($1,$2,$3,$4,$5,$6,$7,$8)`, tenant, resolutionID, exceptionID, action, notes, subject, successor, authorization)
	if postgresConflict(err) {
		return franchisejourney.DeliveryResolution{}, franchisejourney.ErrConflict
	}
	if err != nil {
		return franchisejourney.DeliveryResolution{}, err
	}
	resolved := franchisejourney.DeliveryException{ID: exceptionID, OrganizationID: organization, HandoverID: originalHandoverID, CustomerSubject: customer, State: "resolved", Version: version + 1, ResolutionAction: action}
	if result.SuccessorHandover != nil {
		resolved.SuccessorHandoverID = successorID
	} else {
		resolved.ReturnAuthorizationID = authorizationID
	}
	err = tx.QueryRow(ctx, `update sales.delivery_exception set state='resolved',resolved_at=clock_timestamp(),resolved_by_subject=$5,version=version+1 where tenant_id=$1 and organization_id=$2 and exception_id=$3 and version=$4 and state='open' returning reason_code,details,created_at,resolved_at`, tenant, organization, exceptionID, version, subject).Scan(&resolved.ReasonCode, &resolved.Details, &resolved.CreatedAt, &resolved.ResolvedAt)
	if errors.Is(err, pgx.ErrNoRows) || postgresConflict(err) {
		return franchisejourney.DeliveryResolution{}, franchisejourney.ErrConflict
	}
	if err != nil {
		return franchisejourney.DeliveryResolution{}, err
	}
	result.Exception = resolved
	if err = writeJourneyOutbox(ctx, tx, tenant, eventID, "delivery-exception", exceptionID, resolved.Version, "delivery-exception.resolved", map[string]any{"organization_id": organization, "action": action, "successor_handover_id": resolved.SuccessorHandoverID, "return_authorization_id": resolved.ReturnAuthorizationID, "actor_subject": subject}); err != nil {
		return franchisejourney.DeliveryResolution{}, err
	}
	return result, tx.Commit(ctx)
}

func (r *FranchiseJourney) ReturnCases(ctx context.Context, tenant, organization string, limit int) ([]franchisejourney.ReturnCase, error) {
	if limit < 1 || limit > 100 {
		return nil, franchisejourney.ErrInvalid
	}
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead, AccessMode: pgx.ReadOnly})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	rows, err := tx.Query(ctx, `select a.authorization_id,a.organization_id,a.order_id,a.stock_unit_id,a.customer_principal_id,a.disposition,a.authorized_at,r.receipt_id,r.received_serial_number,r.condition_code,r.notes,r.evidence_sha256_hex,r.received_by_subject,r.received_at,d.disposition_id,d.inventory_action,d.customer_remedy,d.notes,d.decided_by_subject,d.decided_at, `+returnScopePredicate+` and (r.receipt_id is null or (r.organization_id,r.order_id,r.stock_unit_id,r.customer_principal_id)=(a.organization_id,a.order_id,a.stock_unit_id,a.customer_principal_id)) from sales.return_authorization a left join sales.return_receipt r on r.tenant_id=a.tenant_id and r.authorization_id=a.authorization_id left join sales.return_disposition d on d.tenant_id=r.tenant_id and d.receipt_id=r.receipt_id where a.tenant_id=$1 and a.organization_id=$2 order by a.authorized_at desc,a.authorization_id limit $3`, tenant, organization, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []franchisejourney.ReturnCase{}
	dispositions := map[string]*franchisejourney.ReturnDisposition{}
	for rows.Next() {
		var item franchisejourney.ReturnCase
		var receiptID, serial, condition, receiptNotes, evidence, receiver *string
		var receivedAt *time.Time
		var dispositionID, inventoryAction, remedy, dispositionNotes, decider *string
		var decidedAt *time.Time
		var scopeValid bool
		if err = rows.Scan(&item.AuthorizationID, &item.OrganizationID, &item.OrderID, &item.StockUnitID, &item.CustomerSubject, &item.AuthorizedAction, &item.AuthorizedAt, &receiptID, &serial, &condition, &receiptNotes, &evidence, &receiver, &receivedAt, &dispositionID, &inventoryAction, &remedy, &dispositionNotes, &decider, &decidedAt, &scopeValid); err != nil {
			return nil, err
		}
		if !scopeValid {
			return nil, franchisejourney.ErrConflict
		}
		if receiptID != nil {
			item.Receipt = &franchisejourney.ReturnReceipt{ID: *receiptID, AuthorizationID: item.AuthorizationID, OrganizationID: item.OrganizationID, OrderID: item.OrderID, StockUnitID: item.StockUnitID, CustomerSubject: item.CustomerSubject, ReceivedSerialNumber: *serial, ConditionCode: *condition, Notes: *receiptNotes, EvidenceSHA256: *evidence, ReceivedBySubject: *receiver, ReceivedAt: *receivedAt}
		}
		if dispositionID != nil {
			item.Disposition = &franchisejourney.ReturnDisposition{ID: *dispositionID, ReceiptID: *receiptID, InventoryAction: *inventoryAction, CustomerRemedy: *remedy, Notes: *dispositionNotes, DecidedBySubject: *decider, DecidedAt: *decidedAt, Effects: []franchisejourney.ReturnEffectRequest{}}
			dispositions[*dispositionID] = item.Disposition
		}
		items = append(items, item)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	rows.Close()
	if len(dispositions) == 0 {
		return items, tx.Commit(ctx)
	}
	dispositionIDs := make([]string, 0, len(dispositions))
	for id := range dispositions {
		dispositionIDs = append(dispositionIDs, id)
	}
	effectRows, err := tx.Query(ctx, `select e.disposition_id,e.request_id,e.effect_kind,e.owner_context,e.state,e.idempotency_key,e.requested_at from sales.return_effect_request e join sales.return_disposition d on d.tenant_id=e.tenant_id and d.disposition_id=e.disposition_id join sales.return_receipt r on r.tenant_id=d.tenant_id and r.receipt_id=d.receipt_id where e.tenant_id=$1 and r.organization_id=$2 and e.disposition_id=any($3::text[]) order by e.requested_at,e.request_id`, tenant, organization, dispositionIDs)
	if err != nil {
		return nil, err
	}
	defer effectRows.Close()
	for effectRows.Next() {
		var dispositionID string
		var effect franchisejourney.ReturnEffectRequest
		if err = effectRows.Scan(&dispositionID, &effect.ID, &effect.EffectKind, &effect.OwnerContext, &effect.State, &effect.IdempotencyKey, &effect.RequestedAt); err != nil {
			return nil, err
		}
		if disposition := dispositions[dispositionID]; disposition != nil {
			disposition.Effects = append(disposition.Effects, effect)
		}
	}
	effectRows.Close()
	if err = effectRows.Err(); err != nil {
		return nil, err
	}
	return items, tx.Commit(ctx)
}

func (r *FranchiseJourney) ReceiveReturn(ctx context.Context, tenant, organization, subject, authorizationID, serialNumber, conditionCode, notes, evidence, receiptID, eventID string) (franchisejourney.ReturnReceipt, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return franchisejourney.ReturnReceipt{}, err
	}
	defer tx.Rollback(ctx)
	if err = lockReturnAuthorizationScope(ctx, tx, tenant, organization, authorizationID); err != nil {
		return franchisejourney.ReturnReceipt{}, err
	}
	value := franchisejourney.ReturnReceipt{ID: receiptID, AuthorizationID: authorizationID, OrganizationID: organization, ReceivedSerialNumber: serialNumber, ConditionCode: conditionCode, Notes: notes, EvidenceSHA256: evidence, ReceivedBySubject: subject}
	err = tx.QueryRow(ctx, `insert into sales.return_receipt(tenant_id,receipt_id,authorization_id,organization_id,order_id,stock_unit_id,customer_principal_id,received_serial_number,condition_code,notes,evidence_sha256_hex,received_by_subject) select a.tenant_id,$4,a.authorization_id,a.organization_id,a.order_id,a.stock_unit_id,a.customer_principal_id,$5,$6,$7,$8,$9 from sales.return_authorization a where a.tenant_id=$1 and a.organization_id=$2 and a.authorization_id=$3 and a.state='authorized' returning order_id,stock_unit_id,customer_principal_id,received_at`, tenant, organization, authorizationID, receiptID, serialNumber, conditionCode, notes, evidence, subject).Scan(&value.OrderID, &value.StockUnitID, &value.CustomerSubject, &value.ReceivedAt)
	if errors.Is(err, pgx.ErrNoRows) || postgresConflict(err) {
		return franchisejourney.ReturnReceipt{}, franchisejourney.ErrConflict
	}
	if err != nil {
		return franchisejourney.ReturnReceipt{}, err
	}
	if err = writeJourneyOutbox(ctx, tx, tenant, eventID, "return-authorization", authorizationID, 1, "return.received", map[string]any{"organization_id": organization, "receipt_id": receiptID, "order_id": value.OrderID, "stock_unit_id": value.StockUnitID, "condition_code": conditionCode, "evidence_sha256": evidence, "actor_subject": subject}); err != nil {
		return franchisejourney.ReturnReceipt{}, err
	}
	return value, tx.Commit(ctx)
}

func (r *FranchiseJourney) DecideReturn(ctx context.Context, tenant, organization, subject, receiptID, inventoryAction, notes, dispositionID, inventoryRequestID, remedyRequestID, accountingRequestID, fiscalRequestID, eventID string) (franchisejourney.ReturnDisposition, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return franchisejourney.ReturnDisposition{}, err
	}
	defer tx.Rollback(ctx)
	var authorizationID string
	err = tx.QueryRow(ctx, `select authorization_id from sales.return_receipt where tenant_id=$1 and organization_id=$2 and receipt_id=$3`, tenant, organization, receiptID).Scan(&authorizationID)
	if errors.Is(err, pgx.ErrNoRows) {
		return franchisejourney.ReturnDisposition{}, franchisejourney.ErrConflict
	}
	if err != nil {
		return franchisejourney.ReturnDisposition{}, err
	}
	if err = lockReturnAuthorizationScope(ctx, tx, tenant, organization, authorizationID); err != nil {
		return franchisejourney.ReturnDisposition{}, err
	}

	var remedy string
	err = tx.QueryRow(ctx, `select case a.disposition when 'return' then 'refund' else 'exchange' end from sales.return_receipt r join sales.return_authorization a on a.tenant_id=r.tenant_id and a.authorization_id=r.authorization_id where r.tenant_id=$1 and r.organization_id=$2 and r.receipt_id=$3 and (r.organization_id,r.order_id,r.stock_unit_id,r.customer_principal_id)=(a.organization_id,a.order_id,a.stock_unit_id,a.customer_principal_id) for update of r`, tenant, organization, receiptID).Scan(&remedy)
	if errors.Is(err, pgx.ErrNoRows) {
		return franchisejourney.ReturnDisposition{}, franchisejourney.ErrConflict
	}
	if err != nil {
		return franchisejourney.ReturnDisposition{}, err
	}
	value := franchisejourney.ReturnDisposition{ID: dispositionID, ReceiptID: receiptID, InventoryAction: inventoryAction, CustomerRemedy: remedy, Notes: notes, DecidedBySubject: subject, Effects: []franchisejourney.ReturnEffectRequest{}}
	err = tx.QueryRow(ctx, `insert into sales.return_disposition(tenant_id,disposition_id,receipt_id,inventory_action,customer_remedy,notes,decided_by_subject) values($1,$2,$3,$4,$5,$6,$7) returning decided_at`, tenant, dispositionID, receiptID, inventoryAction, remedy, notes, subject).Scan(&value.DecidedAt)
	if postgresConflict(err) {
		return franchisejourney.ReturnDisposition{}, franchisejourney.ErrConflict
	}
	if err != nil {
		return franchisejourney.ReturnDisposition{}, err
	}
	effects := []struct{ id, kind, owner string }{{inventoryRequestID, "inventory", "inventory"}, {remedyRequestID, remedy, map[string]string{"refund": "payment", "exchange": "fulfillment"}[remedy]}, {accountingRequestID, "accounting", "accounting"}}
	if remedy == "refund" {
		effects = append(effects, struct{ id, kind, owner string }{fiscalRequestID, "fiscal", "fiscal"})
	}
	for _, requested := range effects {
		var effect franchisejourney.ReturnEffectRequest
		err = tx.QueryRow(ctx, `insert into sales.return_effect_request(tenant_id,request_id,disposition_id,effect_kind,owner_context,state,idempotency_key) values($1,$2,$3,$4,$5,'requested',$2) returning request_id,effect_kind,owner_context,state,idempotency_key,requested_at`, tenant, requested.id, dispositionID, requested.kind, requested.owner).Scan(&effect.ID, &effect.EffectKind, &effect.OwnerContext, &effect.State, &effect.IdempotencyKey, &effect.RequestedAt)
		if postgresConflict(err) {
			return franchisejourney.ReturnDisposition{}, franchisejourney.ErrConflict
		}
		if err != nil {
			return franchisejourney.ReturnDisposition{}, err
		}
		value.Effects = append(value.Effects, effect)
	}
	if err = writeJourneyOutbox(ctx, tx, tenant, eventID, "return-disposition", dispositionID, 1, "return.disposition-requested", map[string]any{"organization_id": organization, "receipt_id": receiptID, "inventory_action": inventoryAction, "customer_remedy": remedy, "effect_count": len(value.Effects), "actor_subject": subject}); err != nil {
		return franchisejourney.ReturnDisposition{}, err
	}
	return value, tx.Commit(ctx)
}

func (r *FranchiseJourney) CustomerJourney(ctx context.Context, tenant, organization, customer string) (franchisejourney.CustomerJourney, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead, AccessMode: pgx.ReadOnly})
	if err != nil {
		return franchisejourney.CustomerJourney{}, err
	}
	defer tx.Rollback(ctx)
	result, err := readCustomerJourney(ctx, tx, tenant, organization, customer)
	if err != nil {
		return franchisejourney.CustomerJourney{}, err
	}
	if err = tx.Commit(ctx); err != nil {
		return franchisejourney.CustomerJourney{}, err
	}
	return result, nil
}

func readCustomerJourney(ctx context.Context, tx pgx.Tx, tenant, organization, customer string) (franchisejourney.CustomerJourney, error) {
	result := franchisejourney.CustomerJourney{Appointments: []franchisejourney.Appointment{}, Quotes: []franchisejourney.Quote{}, Handovers: []franchisejourney.Handover{}, Exceptions: []franchisejourney.DeliveryException{}}
	rows, err := tx.Query(ctx, `select appointment_id,organization_id,lead_id,coalesce(model_id,''),appointment_kind,starts_at,state,version,coalesce(slot_id,''),coalesce(ends_at,starts_at) from crm.appointment where tenant_id=$1 and organization_id=$2 and customer_principal_id=$3 order by starts_at desc limit 100`, tenant, organization, customer)
	if err != nil {
		return result, err
	}
	for rows.Next() {
		var v franchisejourney.Appointment
		if err = rows.Scan(&v.ID, &v.OrganizationID, &v.LeadID, &v.ModelID, &v.Kind, &v.StartsAt, &v.State, &v.Version, &v.SlotID, &v.EndsAt); err != nil {
			rows.Close()
			return result, err
		}
		result.Appointments = append(result.Appointments, v)
	}
	rows.Close()
	if err = rows.Err(); err != nil {
		return result, err
	}
	rows, err = tx.Query(ctx, `select quotation_id,organization_id,lead_id,coalesce(customer_principal_id,''),variant_id,price_book_id,currency,total_minor_units,valid_until,state,version,coalesce(order_id,'') from sales.quotation where tenant_id=$1 and organization_id=$2 and customer_principal_id=$3 order by created_at desc limit 100`, tenant, organization, customer)
	if err != nil {
		return result, err
	}
	for rows.Next() {
		var v franchisejourney.Quote
		if err = rows.Scan(&v.ID, &v.OrganizationID, &v.LeadID, &v.CustomerSubject, &v.VariantID, &v.PriceBookID, &v.Currency, &v.TotalMinorUnits, &v.ValidUntil, &v.State, &v.Version, &v.OrderID); err != nil {
			rows.Close()
			return result, err
		}
		result.Quotes = append(result.Quotes, v)
	}
	rows.Close()
	if err = rows.Err(); err != nil {
		return result, err
	}
	rows, err = tx.Query(ctx, `select h.handover_id,h.organization_id,h.order_id,h.customer_principal_id,h.stock_unit_id,h.state,h.version,h.customer_accepted_at,coalesce(h.acceptance_evidence_sha256_hex,''),coalesce(h.checklist_id,''),coalesce(h.checklist_version,0),coalesce(t.title,''),h.checklist_completed_at,coalesce(h.supersedes_handover_id,''),exists(select 1 from sales.customer_order o join inventory.stock_unit su on su.tenant_id=o.tenant_id and su.stock_unit_id=h.stock_unit_id and su.organization_id=h.organization_id where o.tenant_id=h.tenant_id and o.order_id=h.order_id and o.organization_id=h.organization_id and o.customer_principal_id=h.customer_principal_id) from sales.delivery_handover h left join sales.delivery_checklist_template t on t.tenant_id=h.tenant_id and t.organization_id=h.organization_id and t.checklist_id=h.checklist_id and t.checklist_version=h.checklist_version where h.tenant_id=$1 and h.organization_id=$2 and h.customer_principal_id=$3 order by h.created_at desc limit 100`, tenant, organization, customer)
	if err != nil {
		return result, err
	}
	for rows.Next() {
		var v franchisejourney.Handover
		var coherent bool
		v.ChecklistItems = []franchisejourney.ChecklistItem{}
		if err = rows.Scan(&v.ID, &v.OrganizationID, &v.OrderID, &v.CustomerSubject, &v.StockUnitID, &v.State, &v.Version, &v.CustomerAcceptedAt, &v.AcceptanceEvidence, &v.ChecklistID, &v.ChecklistVersion, &v.ChecklistTitle, &v.ChecklistCompletedAt, &v.SupersedesHandoverID, &coherent); err != nil {
			rows.Close()
			return result, err
		}
		if !coherent {
			rows.Close()
			return result, franchisejourney.ErrConflict
		}
		result.Handovers = append(result.Handovers, v)
	}
	rows.Close()
	if err = rows.Err(); err != nil {
		return result, err
	}
	indexes := map[string]int{}
	handoverIDs := []string{}
	for index := range result.Handovers {
		indexes[result.Handovers[index].ID] = index
		handoverIDs = append(handoverIDs, result.Handovers[index].ID)
	}
	rows, err = tx.Query(ctx, `select h.handover_id,i.item_id,i.ordinal,i.prompt,i.response_type,i.required from sales.delivery_handover h join sales.delivery_checklist_item i on i.tenant_id=h.tenant_id and i.organization_id=h.organization_id and i.checklist_id=h.checklist_id and i.checklist_version=h.checklist_version where h.tenant_id=$1 and h.organization_id=$2 and h.customer_principal_id=$3 and h.handover_id=any($4::text[]) order by h.created_at desc,h.handover_id,i.ordinal`, tenant, organization, customer, handoverIDs)
	if err != nil {
		return result, err
	}
	for rows.Next() {
		var handoverID string
		var item franchisejourney.ChecklistItem
		if err = rows.Scan(&handoverID, &item.ID, &item.Ordinal, &item.Prompt, &item.ResponseType, &item.Required); err != nil {
			rows.Close()
			return result, err
		}
		if index, ok := indexes[handoverID]; ok {
			result.Handovers[index].ChecklistItems = append(result.Handovers[index].ChecklistItems, item)
		}
	}
	rows.Close()
	if err = rows.Err(); err != nil {
		return result, err
	}
	rows, err = tx.Query(ctx, `select e.exception_id,e.organization_id,e.handover_id,e.customer_principal_id,e.reason_code,e.details,e.state,e.version,e.created_at,e.resolved_at,coalesce(x.action,''),coalesce(x.successor_handover_id,''),coalesce(x.return_authorization_id,''),exists(select 1 from sales.delivery_handover h join sales.customer_order o on o.tenant_id=h.tenant_id and o.order_id=h.order_id and o.organization_id=h.organization_id and o.customer_principal_id=h.customer_principal_id join inventory.stock_unit su on su.tenant_id=h.tenant_id and su.stock_unit_id=h.stock_unit_id and su.organization_id=h.organization_id where h.tenant_id=e.tenant_id and h.handover_id=e.handover_id and h.organization_id=e.organization_id and h.customer_principal_id=e.customer_principal_id) from sales.delivery_exception e left join sales.delivery_exception_resolution x on x.tenant_id=e.tenant_id and x.exception_id=e.exception_id where e.tenant_id=$1 and e.organization_id=$2 and e.customer_principal_id=$3 order by e.created_at desc,e.exception_id limit 100`, tenant, organization, customer)
	if err != nil {
		return result, err
	}
	result.Exceptions, err = scanDeliveryExceptions(rows)
	return result, err
}

func postgresConflict(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && (pgErr.Code == "23503" || pgErr.Code == "23505" || pgErr.Code == "23514" || pgErr.Code == "23P01")
}

func (r *FranchiseJourney) AcceptHandover(ctx context.Context, tenant, organization, customer, handover string, version int64, serialNumber, checklistID string, checklistVersion int64, evidence, eventID string) (franchisejourney.Handover, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return franchisejourney.Handover{}, err
	}
	defer tx.Rollback(ctx)
	if err = lockDeliveryScope(ctx, tx, tenant, organization, handover); err != nil {
		return franchisejourney.Handover{}, err
	}
	var value franchisejourney.Handover
	value.ChecklistItems = []franchisejourney.ChecklistItem{}
	err = tx.QueryRow(ctx, `update sales.delivery_handover h set state='accepted',acceptance_evidence_sha256_hex=$6,customer_accepted_at=clock_timestamp(),updated_at=clock_timestamp(),version=h.version+1 from inventory.stock_unit i,sales.delivery_checklist_template t where h.tenant_id=$1 and h.organization_id=$2 and h.customer_principal_id=$3 and h.handover_id=$4 and h.version=$5 and h.state='presented' and h.checklist_completed_at is not null and h.checklist_id=$8 and h.checklist_version=$9 and i.tenant_id=h.tenant_id and i.stock_unit_id=h.stock_unit_id and i.serial_number=$7 and t.tenant_id=h.tenant_id and t.organization_id=h.organization_id and t.checklist_id=h.checklist_id and t.checklist_version=h.checklist_version returning h.handover_id,h.organization_id,h.order_id,h.customer_principal_id,h.stock_unit_id,h.state,h.version,h.customer_accepted_at,h.acceptance_evidence_sha256_hex,h.checklist_id,h.checklist_version,t.title,h.checklist_completed_at`, tenant, organization, customer, handover, version, evidence, serialNumber, checklistID, checklistVersion).Scan(&value.ID, &value.OrganizationID, &value.OrderID, &value.CustomerSubject, &value.StockUnitID, &value.State, &value.Version, &value.CustomerAcceptedAt, &value.AcceptanceEvidence, &value.ChecklistID, &value.ChecklistVersion, &value.ChecklistTitle, &value.ChecklistCompletedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return franchisejourney.Handover{}, franchisejourney.ErrConflict
	}
	if err != nil {
		return franchisejourney.Handover{}, err
	}
	if err = writeJourneyOutbox(ctx, tx, tenant, eventID, "delivery-handover", handover, value.Version, "delivery-handover.accepted", map[string]any{"organization_id": organization, "customer_subject": customer, "acceptance_evidence_sha256": evidence, "checklist_id": checklistID, "checklist_version": checklistVersion}); err != nil {
		return franchisejourney.Handover{}, err
	}
	return value, tx.Commit(ctx)
}

// Tenant-scoped foreign keys alone do not authorize a relationship. Keep the
// linked rows stable until the command and its audit commit together. This is
// an integrity check, NOT proof of payment, financial release or delivery.
func lockDeliveryScope(ctx context.Context, tx pgx.Tx, tenant, organization, handover string) error {
	var id string
	err := tx.QueryRow(ctx, `select h.handover_id
		from sales.delivery_handover h
		join sales.customer_order o on o.tenant_id=h.tenant_id and o.order_id=h.order_id
			and o.organization_id=h.organization_id and o.customer_principal_id=h.customer_principal_id
		join inventory.stock_unit i on i.tenant_id=h.tenant_id and i.stock_unit_id=h.stock_unit_id
			and i.organization_id=h.organization_id
		where h.tenant_id=$1 and h.organization_id=$2 and h.handover_id=$3
		for update of h for share of o,i`, tenant, organization, handover).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return franchisejourney.ErrConflict
	}
	return err
}

func writeJourneyOutbox(ctx context.Context, tx pgx.Tx, tenant, eventID, aggregateType, aggregateID string, version int64, eventType string, payload any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload) values($1,$2,$3,$4,$5,$6,1,clock_timestamp(),$7)`, tenant, eventID, aggregateType, aggregateID, version, eventType, body)
	if err != nil {
		return fmt.Errorf("write journey outbox: %w", err)
	}
	return nil
}

func (r *FranchiseJourney) QuoteResult(ctx context.Context, tenant, organization, lead, key string) (franchisejourney.Quote, error) {
	var v franchisejourney.Quote
	err := r.pool.QueryRow(ctx, `select q.quotation_id,q.organization_id,q.lead_id,q.variant_id,q.price_book_id,q.currency,q.total_minor_units,q.valid_until,q.state,q.version,coalesce(q.order_id,'')
 from platform.idempotency_record i join sales.quotation q on q.tenant_id=i.tenant_id and q.quotation_id=i.resource_id
 join crm.lead l on l.tenant_id=q.tenant_id and l.organization_id=q.organization_id and l.lead_id=q.lead_id
 where i.tenant_id=$1 and i.scope='franchise-quote' and i.idempotency_key=$4 and i.status='completed' and i.resource_type='quotation' and i.response_code=201
 and i.response_body->>'quotation_id'=q.quotation_id and q.organization_id=$2 and q.lead_id=$3
 and q.customer_principal_id is not distinct from l.customer_principal_id`, tenant, organization, lead, key).Scan(&v.ID, &v.OrganizationID, &v.LeadID, &v.VariantID, &v.PriceBookID, &v.Currency, &v.TotalMinorUnits, &v.ValidUntil, &v.State, &v.Version, &v.OrderID)
	if errors.Is(err, pgx.ErrNoRows) {
		return franchisejourney.Quote{}, franchisejourney.ErrConflict
	}
	if err != nil {
		return franchisejourney.Quote{}, err
	}
	return v, nil
}

func (r *FranchiseJourney) CreateServiceResourceOnce(ctx context.Context, tenant, subject, key, hash string, value franchisejourney.ServiceResource, event string) (franchisejourney.ServiceResource, bool, error) {
	var empty franchisejourney.ServiceResource
	if subject == "" {
		return empty, false, franchisejourney.ErrInvalid
	}
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return empty, false, err
	}
	defer tx.Rollback(ctx)
	result, err := tx.Exec(ctx, `insert into platform.idempotency_record(tenant_id,scope,idempotency_key,request_sha256_hex,status,locked_until,expires_at)values($1,'franchise-resource',$2,$3,'processing',clock_timestamp()+interval '30 seconds',clock_timestamp()+interval '24 hours')on conflict do nothing`, tenant, key, hash)
	if err != nil {
		return empty, false, err
	}
	if result.RowsAffected() == 0 {
		replay, storedHash, err := readServiceResourceCreation(ctx, tx, tenant, value.OrganizationID, key)
		if err != nil {
			return empty, false, err
		}
		if hash != storedHash {
			return empty, false, franchisejourney.ErrConflict
		}
		if err = tx.Commit(ctx); err != nil {
			return empty, false, err
		}
		return replay, true, nil
	}
	value, err = r.createServiceResourceTx(ctx, tx, tenant, subject, value, event)
	if err != nil {
		return empty, false, err
	}
	result, err = tx.Exec(ctx, `update platform.idempotency_record set status='completed',response_code=201,response_body=jsonb_build_object('resource_id',$3::text),resource_type='service-resource',resource_id=$3,locked_until=null where tenant_id=$1 and scope='franchise-resource' and idempotency_key=$2 and status='processing'`, tenant, key, value.ID)
	if err != nil {
		return empty, false, err
	}
	if result.RowsAffected() != 1 {
		return empty, false, franchisejourney.ErrConflict
	}
	if err = tx.Commit(ctx); err != nil {
		return empty, false, err
	}
	return value, false, nil
}

type serviceResourceRowReader interface {
	QueryRow(context.Context, string, ...any) pgx.Row
}

func readServiceResourceCreation(ctx context.Context, q serviceResourceRowReader, tenant, organization, key string) (franchisejourney.ServiceResource, string, error) {
	var v franchisejourney.ServiceResource
	var hash string
	err := q.QueryRow(ctx, `select i.request_sha256_hex,a.resource_id,a.organization_id,coalesce(a.principal_subject,''),a.display_name,a.resource_kind,a.status,a.version,
array(select rs.appointment_kind from crm.resource_skill rs where rs.tenant_id=a.tenant_id and rs.resource_id=a.resource_id order by rs.appointment_kind)
from platform.idempotency_record i join crm.service_resource a on a.tenant_id=i.tenant_id and a.resource_id=i.resource_id
where i.tenant_id=$1 and i.scope='franchise-resource' and i.idempotency_key=$3 and i.status='completed' and i.resource_type='service-resource' and i.response_code=201 and i.response_body->>'resource_id'=a.resource_id and a.organization_id=$2`, tenant, organization, key).Scan(&hash, &v.ID, &v.OrganizationID, &v.PrincipalSubject, &v.DisplayName, &v.Kind, &v.Status, &v.Version, &v.Skills)
	if errors.Is(err, pgx.ErrNoRows) {
		return franchisejourney.ServiceResource{}, "", franchisejourney.ErrConflict
	}
	if err != nil {
		return franchisejourney.ServiceResource{}, "", err
	}
	return v, hash, nil
}
func (r *FranchiseJourney) ServiceResourceCreationResult(ctx context.Context, tenant, organization, key string) (franchisejourney.ServiceResource, error) {
	v, _, err := readServiceResourceCreation(ctx, r.pool, tenant, organization, key)
	return v, err
}

func (r *FranchiseJourney) CreateAppointmentSlotOnce(ctx context.Context, tenant, subject, key, hash string, value franchisejourney.AppointmentSlot, event string) (franchisejourney.AppointmentSlot, bool, error) {
	boundHash, bindErr := r.policy.BindRequestHash(hash)
	if bindErr != nil {
		return franchisejourney.AppointmentSlot{}, false, franchisejourney.ErrInvalid
	}
	hash = boundHash
	var empty franchisejourney.AppointmentSlot
	if subject == "" {
		return empty, false, franchisejourney.ErrInvalid
	}
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return empty, false, err
	}
	defer tx.Rollback(ctx)
	result, err := tx.Exec(ctx, `insert into platform.idempotency_record(tenant_id,scope,idempotency_key,request_sha256_hex,status,locked_until,expires_at)values($1,'franchise-slot',$2,$3,'processing',clock_timestamp()+interval '30 seconds',clock_timestamp()+interval '24 hours')on conflict do nothing`, tenant, key, hash)
	if err != nil {
		return empty, false, err
	}
	if result.RowsAffected() == 0 {
		replay, storedHash, err := readAppointmentSlotCreation(ctx, tx, tenant, value.OrganizationID, key)
		if err != nil {
			return empty, false, err
		}
		if hash != storedHash {
			return empty, false, franchisejourney.ErrConflict
		}
		if err = tx.Commit(ctx); err != nil {
			return empty, false, err
		}
		return replay, true, nil
	}
	value, err = r.createAppointmentSlotTx(ctx, tx, tenant, subject, value, event)
	if err != nil {
		return empty, false, err
	}
	result, err = tx.Exec(ctx, `update platform.idempotency_record set status='completed',response_code=201,response_body=jsonb_build_object('slot_id',$3::text),resource_type='appointment-slot',resource_id=$3,locked_until=null where tenant_id=$1 and scope='franchise-slot' and idempotency_key=$2 and status='processing'`, tenant, key, value.ID)
	if err != nil {
		return empty, false, err
	}
	if result.RowsAffected() != 1 {
		return empty, false, franchisejourney.ErrConflict
	}
	if err = tx.Commit(ctx); err != nil {
		return empty, false, err
	}
	return value, false, nil
}

type appointmentSlotRowReader interface {
	QueryRow(context.Context, string, ...any) pgx.Row
}

func readAppointmentSlotCreation(ctx context.Context, q appointmentSlotRowReader, tenant, organization, key string) (franchisejourney.AppointmentSlot, string, error) {
	var v franchisejourney.AppointmentSlot
	var hash string
	err := q.QueryRow(ctx, `select i.request_sha256_hex,a.slot_id,a.organization_id,a.appointment_kind,a.starts_at,a.ends_at,a.capacity,
(select count(*) from crm.appointment p where p.tenant_id=a.tenant_id and p.slot_id=a.slot_id and p.state in ('requested','confirmed')),a.state,a.version
from platform.idempotency_record i join crm.appointment_slot a on a.tenant_id=i.tenant_id and a.slot_id=i.resource_id
where i.tenant_id=$1 and i.scope='franchise-slot' and i.idempotency_key=$3 and i.status='completed' and i.resource_type='appointment-slot' and i.response_code=201 and i.response_body->>'slot_id'=a.slot_id and a.organization_id=$2`, tenant, organization, key).Scan(&hash, &v.ID, &v.OrganizationID, &v.Kind, &v.StartsAt, &v.EndsAt, &v.Capacity, &v.Booked, &v.State, &v.Version)
	if errors.Is(err, pgx.ErrNoRows) {
		return franchisejourney.AppointmentSlot{}, "", franchisejourney.ErrConflict
	}
	if err != nil {
		return franchisejourney.AppointmentSlot{}, "", err
	}
	return v, hash, nil
}
func (r *FranchiseJourney) AppointmentSlotCreationResult(ctx context.Context, tenant, organization, key string) (franchisejourney.AppointmentSlot, error) {
	v, _, err := readAppointmentSlotCreation(ctx, r.pool, tenant, organization, key)
	return v, err
}

func (r *FranchiseJourney) PublishedDeliveryChecklist(ctx context.Context, tenant, organization, id string, version int64) (franchisejourney.DeliveryChecklist, error) {
	var value franchisejourney.DeliveryChecklist
	err := r.pool.QueryRow(ctx, `select t.checklist_id,t.organization_id,t.checklist_version,t.title,t.state,
 coalesce((select jsonb_agg(jsonb_build_object('id',i.item_id,'ordinal',i.ordinal,'prompt',i.prompt,'response_type',i.response_type,'required',i.required) order by i.ordinal)
 from sales.delivery_checklist_item i where i.tenant_id=t.tenant_id and i.organization_id=t.organization_id and i.checklist_id=t.checklist_id and i.checklist_version=t.checklist_version),'[]'::jsonb)
 from sales.delivery_checklist_template t where t.tenant_id=$1 and t.organization_id=$2 and t.checklist_id=$3 and t.checklist_version=$4 and t.state='published'`, tenant, organization, id, version).Scan(&value.ID, &value.OrganizationID, &value.Version, &value.Title, &value.State, &value.Items)
	if errors.Is(err, pgx.ErrNoRows) {
		return franchisejourney.DeliveryChecklist{}, franchisejourney.ErrConflict
	}
	if err != nil {
		return franchisejourney.DeliveryChecklist{}, err
	}
	return value, nil
}

func (r *FranchiseJourney) DeliveryChecklistCompletion(ctx context.Context, tenant, organization, handover string) (franchisejourney.ChecklistCompletion, error) {
	var v franchisejourney.ChecklistCompletion
	err := r.pool.QueryRow(ctx, `select h.handover_id,h.organization_id,h.state,h.version,h.checklist_id,h.checklist_version,h.checklist_completed_at,h.checklist_completed_by_subject,
 coalesce((select jsonb_agg(jsonb_build_object('item_id',r.item_id,'response_text',r.response_text,'evidence_sha256',coalesce(r.evidence_sha256_hex,'')) order by r.item_id collate "C") from sales.delivery_checklist_response r where r.tenant_id=h.tenant_id and r.handover_id=h.handover_id),'[]'::jsonb)
 from sales.delivery_handover h
 join sales.customer_order o on o.tenant_id=h.tenant_id and o.order_id=h.order_id and o.organization_id=h.organization_id and o.customer_principal_id=h.customer_principal_id
 join inventory.stock_unit i on i.tenant_id=h.tenant_id and i.stock_unit_id=h.stock_unit_id and i.organization_id=h.organization_id
 join sales.delivery_checklist_template t on t.tenant_id=h.tenant_id and t.organization_id=h.organization_id and t.checklist_id=h.checklist_id and t.checklist_version=h.checklist_version and t.state='published'
 where h.tenant_id=$1 and h.organization_id=$2 and h.handover_id=$3 and h.checklist_completed_at is not null
 and not exists(select 1 from sales.delivery_checklist_response r where r.tenant_id=h.tenant_id and r.handover_id=h.handover_id and (r.organization_id,r.checklist_id,r.checklist_version,r.answered_by_subject) is distinct from (h.organization_id,h.checklist_id,h.checklist_version,h.checklist_completed_by_subject))`, tenant, organization, handover).Scan(&v.HandoverID, &v.OrganizationID, &v.State, &v.Version, &v.ChecklistID, &v.ChecklistVersion, &v.CompletedAt, &v.ActorSubject, &v.Responses)
	if errors.Is(err, pgx.ErrNoRows) {
		return franchisejourney.ChecklistCompletion{}, franchisejourney.ErrConflict
	}
	if err != nil {
		return franchisejourney.ChecklistCompletion{}, err
	}
	return v, nil
}

// Shared read invariant; the authorization alias is always scoped by tenant and organization.
const returnScopePredicate = `exists(select 1 from sales.delivery_handover h
 join sales.customer_order o on o.tenant_id=h.tenant_id and o.order_id=h.order_id and o.organization_id=h.organization_id and o.customer_principal_id=h.customer_principal_id
 join inventory.stock_unit i on i.tenant_id=h.tenant_id and i.stock_unit_id=h.stock_unit_id and i.organization_id=h.organization_id
 join sales.delivery_exception x on x.tenant_id=h.tenant_id and x.handover_id=h.handover_id and x.organization_id=h.organization_id and x.customer_principal_id=h.customer_principal_id
 where h.tenant_id=a.tenant_id and h.handover_id=a.handover_id and x.exception_id=a.exception_id
 and (h.organization_id,h.order_id,h.stock_unit_id,h.customer_principal_id)=(a.organization_id,a.order_id,a.stock_unit_id,a.customer_principal_id))`

func lockReturnAuthorizationScope(ctx context.Context, tx pgx.Tx, tenant, organization, authorization string) error {
	var id string
	err := tx.QueryRow(ctx, `select a.authorization_id from sales.return_authorization a
 join sales.delivery_handover h on h.tenant_id=a.tenant_id and h.handover_id=a.handover_id and (h.organization_id,h.order_id,h.stock_unit_id,h.customer_principal_id)=(a.organization_id,a.order_id,a.stock_unit_id,a.customer_principal_id)
 join sales.customer_order o on o.tenant_id=h.tenant_id and o.order_id=h.order_id and o.organization_id=h.organization_id and o.customer_principal_id=h.customer_principal_id
 join inventory.stock_unit i on i.tenant_id=h.tenant_id and i.stock_unit_id=h.stock_unit_id and i.organization_id=h.organization_id
 join sales.delivery_exception x on x.tenant_id=a.tenant_id and x.exception_id=a.exception_id and x.handover_id=h.handover_id and x.organization_id=h.organization_id and x.customer_principal_id=h.customer_principal_id
 where a.tenant_id=$1 and a.organization_id=$2 and a.authorization_id=$3 for update of h for share of a,x,o,i`, tenant, organization, authorization).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return franchisejourney.ErrConflict
	}
	return err
}

func (r *FranchiseJourney) ReturnCaseResult(ctx context.Context, tenant, organization, authorization string) (franchisejourney.ReturnCase, error) {
	var v franchisejourney.ReturnCase
	var scopeValid bool
	err := r.pool.QueryRow(ctx, `select a.authorization_id,a.organization_id,a.order_id,a.stock_unit_id,a.customer_principal_id,a.disposition,a.authorized_at,
 case when r.receipt_id is null then null else jsonb_build_object('id',r.receipt_id,'authorization_id',r.authorization_id,'organization_id',r.organization_id,'order_id',r.order_id,'stock_unit_id',r.stock_unit_id,'customer_subject',r.customer_principal_id,'received_serial_number',r.received_serial_number,'condition_code',r.condition_code,'notes',r.notes,'evidence_sha256',r.evidence_sha256_hex,'received_by_subject',r.received_by_subject,'received_at',r.received_at) end,
 case when d.disposition_id is null then null else jsonb_build_object('id',d.disposition_id,'receipt_id',d.receipt_id,'inventory_action',d.inventory_action,'customer_remedy',d.customer_remedy,'notes',d.notes,'decided_by_subject',d.decided_by_subject,'decided_at',d.decided_at,'effect_requests',coalesce((select jsonb_agg(jsonb_build_object('id',e.request_id,'effect_kind',e.effect_kind,'owner_context',e.owner_context,'state',e.state,'idempotency_key',e.idempotency_key,'requested_at',e.requested_at) order by e.effect_kind) from sales.return_effect_request e where e.tenant_id=d.tenant_id and e.disposition_id=d.disposition_id),'[]'::jsonb)) end,
 `+returnScopePredicate+` and (r.receipt_id is null or (r.organization_id,r.order_id,r.stock_unit_id,r.customer_principal_id)=(a.organization_id,a.order_id,a.stock_unit_id,a.customer_principal_id))
 from sales.return_authorization a left join sales.return_receipt r on r.tenant_id=a.tenant_id and r.authorization_id=a.authorization_id left join sales.return_disposition d on d.tenant_id=r.tenant_id and d.receipt_id=r.receipt_id
 where a.tenant_id=$1 and a.organization_id=$2 and a.authorization_id=$3`, tenant, organization, authorization).Scan(&v.AuthorizationID, &v.OrganizationID, &v.OrderID, &v.StockUnitID, &v.CustomerSubject, &v.AuthorizedAction, &v.AuthorizedAt, &v.Receipt, &v.Disposition, &scopeValid)
	if errors.Is(err, pgx.ErrNoRows) {
		return franchisejourney.ReturnCase{}, franchisejourney.ErrConflict
	}
	if err != nil {
		return franchisejourney.ReturnCase{}, err
	}
	if !scopeValid {
		return franchisejourney.ReturnCase{}, franchisejourney.ErrConflict
	}
	return v, nil
}
````

### FILE: `internal/platform/postgres/franchisejourney_integration_test.go`

```yaml
block_id: "GO-FRANCHISE-JOURNEY:file:07"
operation: CREATE
provenance: AUTHORED
source: "local implementation governed by the upstream references in metadata"
license: "LicenseRef-Workspace-Owner"
sha256: "af93b84ec8f18f2b7afad402f399f23fde3865d8ab22d874f42a8da74767afe4"
variables: []
secrets_allowed: false
```

````go
package postgres

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"elite.local/enterprise/internal/franchisejourney"
	"elite.local/enterprise/internal/platform/randomid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type journeyReadProbeKey struct{}

func TestBookingRequiresCurrentAvailability(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil || cfg.ConnConfig.Host != "127.0.0.1" || !strings.HasPrefix(cfg.ConnConfig.Database, "elite_confirmation_") {
		t.Fatal("requires disposable loopback database")
	}
	for _, mode := range []string{"cancelled", "unavailable", "race-cancel", "race-unavailable"} {
		t.Run(mode, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			pool, err := pgxpool.NewWithConfig(ctx, cfg.Copy())
			if err != nil {
				t.Fatal(err)
			}
			defer pool.Close()
			tenant := randomid.Generator{}.New()
			tenantCode := "booking-" + strings.ReplaceAll(tenant, "-", "")
			for _, q := range []string{
				`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1::uuid,'booking-'||replace(($1::uuid)::text,'-',''),'Synthetic','Synthetic')`,
				`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'store','store','Synthetic','store')`,
				`insert into org.public_location(tenant_id,organization_id,city,region,country,published)values($1,'store','Synthetic','Synthetic','AR',true)`,
				`insert into crm.lead(tenant_id,lead_id,organization_id,lifecycle_state,source_code,contact_payload)values($1,'lead','store','new','fixture','{}')`,
			} {
				if _, err = pool.Exec(ctx, q, tenant); err != nil {
					t.Fatal(err)
				}
			}
			repo := NewFranchiseJourney(pool)
			start := time.Now().UTC().Add(72 * time.Hour).Truncate(time.Second)
			end := start.Add(time.Hour)
			window := franchisejourney.AvailabilityEntry{ID: "working", OrganizationID: "store", EntryType: "working", StartsAt: start, EndsAt: end}
			if _, err = repo.CreateAvailability(ctx, tenant, "scheduler", window, randomid.Generator{}.New()); err != nil {
				t.Fatal(err)
			}
			if _, err = repo.CreateAppointmentSlot(ctx, tenant, franchisejourney.AppointmentSlot{ID: "slot", OrganizationID: "store", Kind: "consultation", StartsAt: start, EndsAt: end, Capacity: 1}, randomid.Generator{}.New()); err != nil {
				t.Fatal(err)
			}
			appointment := franchisejourney.Appointment{ID: "booking", LeadID: "lead", Kind: "consultation", StartsAt: start}
			book := func() error {
				_, _, e := repo.RequestAppointment(ctx, tenantCode, "store", "booking-key-v305", appointment, strings.Repeat("a", 64), randomid.Generator{}.New())
				return e
			}
			change := func() error {
				if mode == "cancelled" || mode == "race-cancel" {
					_, e := repo.CancelAvailability(ctx, tenant, "store", "working", 1, "scheduler", "schedule-correction", randomid.Generator{}.New())
					return e
				}
				window.ID, window.EntryType, window.ReasonCode = "absence", "unavailable", "synthetic-absence"
				_, e := repo.CreateAvailability(ctx, tenant, "scheduler", window, randomid.Generator{}.New())
				return e
			}
			if !strings.HasPrefix(mode, "race-") {
				if err = change(); err != nil {
					t.Fatal(err)
				}
				items, e := repo.PublicAppointmentSlots(ctx, tenantCode, "store", "consultation", start.Add(-time.Minute), end.Add(time.Minute))
				if e != nil || len(items) != 0 {
					t.Errorf("unavailable slot published count=%d err=%v", len(items), e)
				}
				if e = book(); !errors.Is(e, franchisejourney.ErrConflict) {
					t.Errorf("booking outside availability accepted: %v", e)
				}
				var effects int
				if e = pool.QueryRow(ctx, `select (select count(*) from crm.appointment where tenant_id=$1)+(select count(*) from platform.idempotency_record where tenant_id=$1 and scope='public-appointment')+(select count(*) from platform.outbox_event where tenant_id=$1 and event_type='appointment.requested')`, tenant).Scan(&effects); e != nil || effects != 0 {
					t.Errorf("unexpected booking effects=%d err=%v", effects, e)
				}
			} else {
				startTogether := make(chan struct{})
				results := make(chan error, 2)
				go func() { <-startTogether; results <- book() }()
				go func() { <-startTogether; results <- change() }()
				close(startTogether)
				success, conflicts := 0, 0
				for i := 0; i < 2; i++ {
					e := <-results
					if e == nil {
						success++
					} else if errors.Is(e, franchisejourney.ErrConflict) {
						conflicts++
					} else {
						t.Fatal(e)
					}
				}
				if success != 1 || conflicts != 1 {
					t.Fatalf("race success=%d conflicts=%d", success, conflicts)
				}
				var booked, eligible int
				if err = pool.QueryRow(ctx, `select (select count(*) from crm.appointment where tenant_id=$1),(select count(*) from crm.availability_entry where tenant_id=$1 and availability_id='working' and state='active')-(select count(*) from crm.availability_entry where tenant_id=$1 and entry_type='unavailable' and state='active')`, tenant).Scan(&booked, &eligible); err != nil || booked != eligible {
					t.Fatal("contradictory schedule", booked, eligible, err)
				}
				if booked == 1 {
					replayed, wasReplay, e := repo.RequestAppointment(ctx, tenantCode, "store", "booking-key-v305", appointment, strings.Repeat("a", 64), randomid.Generator{}.New())
					if e != nil || !wasReplay || replayed.ID != "booking" {
						t.Fatal("replay changed", replayed, wasReplay, e)
					}
				}
			}
		})
	}
}

func TestAvailabilityCancellationAtomicity(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil || cfg.ConnConfig.Host != "127.0.0.1" || !strings.HasPrefix(cfg.ConnConfig.Database, "elite_confirmation_") {
		t.Fatal("requires disposable loopback database")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := randomid.Generator{}.New()
	for _, q := range []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1::uuid,'avail-'||replace(($1::uuid)::text,'-',''),'Synthetic','Synthetic')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'store','store','Synthetic','store')`,
		`insert into crm.lead(tenant_id,lead_id,organization_id,lifecycle_state,source_code,contact_payload)values($1,'lead','store','new','fixture','{}')`,
	} {
		if _, err = pool.Exec(ctx, q, tenant); err != nil {
			t.Fatal(err)
		}
	}
	repo := NewFranchiseJourney(pool)
	start := time.Now().UTC().Add(72 * time.Hour).Truncate(time.Second)
	creationEvent := randomid.Generator{}.New()
	if _, err = repo.CreateAvailability(ctx, tenant, "creator", franchisejourney.AvailabilityEntry{ID: "atomic-window", OrganizationID: "store", EntryType: "working", StartsAt: start, EndsAt: start.Add(time.Hour)}, creationEvent); err != nil {
		t.Fatal(err)
	}
	if _, err = repo.CancelAvailability(ctx, tenant, "other", "atomic-window", 1, "operator", "schedule-correction", randomid.Generator{}.New()); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("foreign organization", err)
	}
	if _, err = repo.CancelAvailability(ctx, randomid.Generator{}.New(), "store", "atomic-window", 1, "operator", "schedule-correction", randomid.Generator{}.New()); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("foreign tenant", err)
	}
	if _, err = repo.CancelAvailability(ctx, tenant, "store", "atomic-window", 2, "operator", "schedule-correction", randomid.Generator{}.New()); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("stale version", err)
	}
	_, err = repo.CancelAvailability(ctx, tenant, "store", "atomic-window", 1, "operator", "schedule-correction", creationEvent)
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) || pgErr.Code != "23505" {
		t.Fatal("expected duplicate audit event", err)
	}
	var state, actor string
	var version, events int
	if err = pool.QueryRow(ctx, `select state,version,coalesce(cancelled_by_subject,'') from crm.availability_entry where tenant_id=$1 and availability_id='atomic-window'`, tenant).Scan(&state, &version, &actor); err != nil || state != "active" || version != 1 || actor != "" {
		t.Fatal("partial cancellation", state, version, actor, err)
	}
	results := make(chan error, 2)
	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, e := repo.CancelAvailability(ctx, tenant, "store", "atomic-window", 1, "operator", "schedule-correction", randomid.Generator{}.New())
			results <- e
		}()
	}
	wg.Wait()
	close(results)
	success, conflicts := 0, 0
	for e := range results {
		if e == nil {
			success++
		} else if errors.Is(e, franchisejourney.ErrConflict) {
			conflicts++
		} else {
			t.Fatal(e)
		}
	}
	if success != 1 || conflicts != 1 {
		t.Fatal("race", success, conflicts)
	}
	pool.Close()
	pool, err = pgxpool.NewWithConfig(ctx, cfg.Copy())
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	repo = NewFranchiseJourney(pool)
	items, err := repo.Availability(ctx, tenant, "store", "", start.Add(-time.Minute), start.Add(2*time.Hour))
	if err != nil || len(items) != 1 || items[0].ID != "atomic-window" || items[0].State != "cancelled" || items[0].Version != 2 {
		t.Fatal("reconnected read", items, err)
	}
	if err = pool.QueryRow(ctx, `select count(*) from platform.outbox_event where tenant_id=$1 and aggregate_id='atomic-window' and event_type='availability-entry.cancelled' and payload->>'actor_subject'='operator'`, tenant).Scan(&events); err != nil || events != 1 {
		t.Fatal("audit", events, err)
	}
	// Existing cancellation policy: an active appointment prevents removing its working interval.
	start = start.Add(24 * time.Hour)
	if _, err = repo.CreateAvailability(ctx, tenant, "creator", franchisejourney.AvailabilityEntry{ID: "booked-window", OrganizationID: "store", EntryType: "working", StartsAt: start, EndsAt: start.Add(2 * time.Hour)}, randomid.Generator{}.New()); err != nil {
		t.Fatal(err)
	}
	if _, err = pool.Exec(ctx, `insert into crm.appointment(tenant_id,appointment_id,organization_id,lead_id,appointment_kind,starts_at,state,version)values($1,'booked','store','lead','consultation',$2,'requested',1)`, tenant, start.Add(30*time.Minute)); err != nil {
		t.Fatal(err)
	}
	if _, err = repo.CancelAvailability(ctx, tenant, "store", "booked-window", 1, "operator", "schedule-correction", randomid.Generator{}.New()); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("active appointment guard", err)
	}
	if err = pool.QueryRow(ctx, `select state,version,(select count(*) from platform.outbox_event where tenant_id=$1 and aggregate_id='booked-window' and event_type='availability-entry.cancelled') from crm.availability_entry where tenant_id=$1 and availability_id='booked-window'`, tenant).Scan(&state, &version, &events); err != nil || state != "active" || version != 1 || events != 0 {
		t.Fatal("booked interval changed", state, version, events, err)
	}
}

func TestLeadAuditedAtomicity(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil || cfg.ConnConfig.Host != "127.0.0.1" || !strings.HasPrefix(cfg.ConnConfig.Database, "elite_confirmation_") {
		t.Fatal("requires disposable loopback database")
	}
	ctx := context.Background()
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := randomid.Generator{}.New()
	eventOne, eventTwo := randomid.Generator{}.New(), randomid.Generator{}.New()
	for _, q := range []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1::uuid,'lead-'||replace(($1::uuid)::text,'-',''),'Synthetic','Synthetic')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'store','store','Synthetic','store')`,
		`insert into crm.customer_profile(tenant_id,customer_principal_id,display_name,email_normalized)values($1,'customer','Synthetic','fixture@example.invalid')`,
		`insert into crm.lead(tenant_id,lead_id,organization_id,customer_principal_id,lifecycle_state,source_code,contact_payload)values($1,'atomic-lead','store','customer','new','fixture','{}')`,
	} {
		if _, err := pool.Exec(ctx, q, tenant); err != nil {
			t.Fatal(err)
		}
	}
	repo := NewFranchiseJourney(pool)
	if _, err := repo.AssignLeadAs(ctx, tenant, "store", "atomic-lead", "assignee", 1, "empty", ""); !errors.Is(err, franchisejourney.ErrInvalid) {
		t.Fatal("missing actor accepted", err)
	}
	if _, err := repo.TransitionLeadAs(ctx, tenant, "store", "atomic-lead", "new", "contacted", 1, "empty", ""); !errors.Is(err, franchisejourney.ErrInvalid) {
		t.Fatal("missing actor transition accepted", err)
	}
	if _, err := repo.AssignLeadAs(ctx, tenant, "store", "atomic-lead", "assignee", 1, eventOne, "actor-one"); err != nil {
		t.Fatal(err)
	}
	// Duplicate event ID must roll back the preceding update, not merely its audit.
	if _, err := repo.AssignLeadAs(ctx, tenant, "store", "atomic-lead", "other", 2, eventOne, "actor-two"); err == nil {
		t.Fatal("duplicate event assignment accepted")
	}
	if _, err := repo.TransitionLeadAs(ctx, tenant, "store", "atomic-lead", "new", "contacted", 2, eventOne, "actor-two"); err == nil {
		t.Fatal("duplicate event transition accepted")
	}
	var state, assignee string
	var version, events int
	if err := pool.QueryRow(ctx, `select lifecycle_state,assigned_subject,version from crm.lead where tenant_id=$1 and lead_id='atomic-lead'`, tenant).Scan(&state, &assignee, &version); err != nil || state != "new" || assignee != "assignee" || version != 2 {
		t.Fatal("partial update", state, assignee, version, err)
	}
	if _, err := repo.TransitionLeadAs(ctx, tenant, "store", "atomic-lead", "new", "contacted", 2, eventTwo, "actor-two"); err != nil {
		t.Fatal(err)
	}
	pool.Close()
	pool, err = pgxpool.NewWithConfig(ctx, cfg.Copy())
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	if err := pool.QueryRow(ctx, `select lifecycle_state,assigned_subject,version from crm.lead where tenant_id=$1 and lead_id='atomic-lead'`, tenant).Scan(&state, &assignee, &version); err != nil || state != "contacted" || assignee != "assignee" || version != 3 {
		t.Fatal("recovery", state, assignee, version, err)
	}
	if err := pool.QueryRow(ctx, `select count(*) from platform.outbox_event where tenant_id=$1 and aggregate_id='atomic-lead'`, tenant).Scan(&events); err != nil || events != 2 {
		t.Fatal("events", events, err)
	}
	for _, pair := range [][2]string{{eventOne, "actor-one"}, {eventTwo, "actor-two"}} {
		var actor string
		if err := pool.QueryRow(ctx, `select payload->>'actor_subject' from platform.outbox_event where tenant_id=$1 and event_id=$2`, tenant, pair[0]).Scan(&actor); err != nil || actor != pair[1] {
			t.Fatal("actor", actor, err)
		}
	}
}

type journeyReadProbe struct {
	afterFirstRead func() error
	once           sync.Once
	err            error
	inTransaction  bool
	boundedItems   int
}

func (p *journeyReadProbe) TraceQueryStart(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	if strings.Contains(data.SQL, "h.handover_id=any($4::text[])") && len(data.Args) == 4 {
		if ids, ok := data.Args[3].([]string); ok {
			p.boundedItems = len(ids)
		}
	}
	return context.WithValue(ctx, journeyReadProbeKey{}, data.SQL)
}
func (p *journeyReadProbe) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, _ pgx.TraceQueryEndData) {
	query, _ := ctx.Value(journeyReadProbeKey{}).(string)
	if strings.HasPrefix(query, "select appointment_id,organization_id") {
		p.once.Do(func() { p.inTransaction = conn.PgConn().TxStatus() == 'T'; p.err = p.afterFirstRead() })
	}
}

func TestFranchiseJourneyPersistenceIsolationAndReplay(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := "018f4d4a-7b36-7a21-8d10-2f4c54c29b01"
	cleanup := func() {
		tx, cleanupErr := pool.Begin(ctx)
		if cleanupErr != nil {
			return
		}
		defer tx.Rollback(ctx)
		_, _ = tx.Exec(ctx, `set local session_replication_role=replica`)
		for _, q := range []string{`delete from sales.return_effect_attempt where tenant_id=$1`, `delete from sales.return_effect_resume where tenant_id=$1`, `delete from sales.return_effect_execution where tenant_id=$1`, `delete from sales.return_effect_request where tenant_id=$1`, `delete from sales.return_disposition where tenant_id=$1`, `delete from sales.return_receipt where tenant_id=$1`} {
			_, _ = tx.Exec(ctx, q, tenant)
		}
		for _, q := range []string{`delete from sales.delivery_exception_resolution where tenant_id=$1`, `delete from sales.return_authorization where tenant_id=$1`, `delete from sales.delivery_exception where tenant_id=$1`} {
			_, _ = tx.Exec(ctx, q, tenant)
		}
		for _, q := range []string{`delete from platform.outbox_event where tenant_id=$1`, `delete from platform.idempotency_record where tenant_id=$1`, `delete from sales.delivery_checklist_response where tenant_id=$1`, `delete from sales.delivery_handover where tenant_id=$1`, `delete from sales.delivery_checklist_item where tenant_id=$1`, `delete from sales.delivery_checklist_template where tenant_id=$1`, `delete from sales.quotation_acceptance where tenant_id=$1`, `delete from sales.quotation where tenant_id=$1`, `delete from sales.customer_order_line where tenant_id=$1`, `delete from sales.customer_order where tenant_id=$1`, `delete from crm.appointment_transition where tenant_id=$1`, `delete from crm.appointment_resource where tenant_id=$1`, `delete from crm.appointment where tenant_id=$1`, `delete from crm.appointment_slot where tenant_id=$1`, `delete from crm.availability_entry where tenant_id=$1`, `delete from crm.resource_skill where tenant_id=$1`, `delete from crm.service_resource where tenant_id=$1`, `delete from crm.consent_evidence where tenant_id=$1`, `delete from crm.lead where tenant_id=$1`, `delete from crm.customer_profile where tenant_id=$1`, `delete from org.public_location where tenant_id=$1`, `delete from inventory.stock_unit where tenant_id=$1`, `delete from pricing.price_book_entry where tenant_id=$1`, `delete from pricing.price_book where tenant_id=$1`, `delete from catalog.vehicle_variant where tenant_id=$1`, `delete from catalog.vehicle_model where tenant_id=$1`, `delete from org.organization where tenant_id=$1`, `delete from platform.tenant where tenant_id=$1`} {
			_, _ = tx.Exec(ctx, q, tenant)
		}
		_ = tx.Commit(ctx)
	}
	cleanup()
	defer cleanup()
	fixtures := []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1,'journey-api','Journey','Journey')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'store','store','Store','store'),($1,'other','other','Other','store')`,
		`insert into org.public_location(tenant_id,organization_id,city,region,country,published)values($1,'store','Cordoba','Cordoba','AR',true)`,
		`insert into catalog.vehicle_model(tenant_id,model_id,model_code,display_name,vehicle_class,lifecycle_state)values($1,'model','model','Model','bicycle','active')`,
		`insert into catalog.vehicle_variant(tenant_id,variant_id,model_id,variant_code,display_name,battery_specification,lifecycle_state)values($1,'variant','model','variant','Variant','{}','active')`,
		`insert into crm.customer_profile(tenant_id,customer_principal_id,display_name,email_normalized)values($1,'customer','Customer','customer@example.test')`,
		`insert into crm.lead(tenant_id,lead_id,organization_id,customer_principal_id,model_id,lifecycle_state,source_code,contact_payload)values($1,'lead','store','customer','model','new','public-web','{}')`,
		`insert into pricing.price_book(tenant_id,price_book_id,market,currency,valid_from,status)values($1,'retail','AR','ARS',clock_timestamp()-interval '1 day','active')`,
		`insert into pricing.price_book_entry(tenant_id,price_book_id,variant_id,amount_minor_units,tax_mode)values($1,'retail','variant',123456,'inclusive')`,
		`insert into inventory.stock_unit(tenant_id,stock_unit_id,organization_id,variant_id,serial_number,state,version,received_at)values($1,'stock','store','variant','SERIAL-JOURNEY','sold',1,clock_timestamp())`,
		`insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version)values($1,'order','store','customer','delivered','ARS',123456,1)`,
		`insert into sales.delivery_handover(tenant_id,handover_id,organization_id,order_id,customer_principal_id,stock_unit_id,state,version)values($1,'handover','store','order','customer','stock','prepared',1)`,
		`insert into sales.delivery_handover(tenant_id,handover_id,organization_id,order_id,customer_principal_id,stock_unit_id,state,version)values($1,'handover-reject','store','order','customer','stock','prepared',1)`,
		`insert into sales.delivery_handover(tenant_id,handover_id,organization_id,order_id,customer_principal_id,stock_unit_id,state,version)values($1,'handover-return','store','order','customer','stock','prepared',1)`,
	}
	for _, q := range fixtures {
		if _, err = pool.Exec(ctx, q, tenant); err != nil {
			t.Fatal(err)
		}
	}
	repo := NewFranchiseJourney(pool)
	locations, err := repo.PublicLocations(ctx, "journey-api")
	if err != nil || len(locations) != 1 || locations[0].OrganizationID != "store" {
		t.Fatalf("locations=%+v err=%v", locations, err)
	}
	slotStart := time.Now().Add(time.Hour).Truncate(time.Second)
	slotEnd := slotStart.Add(45 * time.Minute)
	workingEnd := slotStart.Add(8 * time.Hour)
	organizationWorking, err := repo.CreateAvailability(ctx, tenant, "scheduler", franchisejourney.AvailabilityEntry{ID: "store-working", OrganizationID: "store", EntryType: "working", StartsAt: slotStart.Add(-time.Hour), EndsAt: workingEnd, State: "active", Version: 1}, "018f4d4a-7b36-7a21-8d10-2f4c54c29b20")
	if err != nil || organizationWorking.State != "active" {
		t.Fatalf("organization working window=%+v err=%v", organizationWorking, err)
	}
	slot, err := repo.CreateAppointmentSlot(ctx, tenant, franchisejourney.AppointmentSlot{ID: "test-drive-slot", OrganizationID: "store", Kind: "test-drive", StartsAt: slotStart, EndsAt: slotEnd, Capacity: 1, State: "open", Version: 1}, "018f4d4a-7b36-7a21-8d10-2f4c54c29b00")
	if err != nil || slot.Booked != 0 || slot.Capacity != 1 {
		t.Fatalf("slot=%+v err=%v", slot, err)
	}
	available, err := repo.PublicAppointmentSlots(ctx, "journey-api", "store", "test-drive", slotStart.Add(-time.Minute), slotEnd.Add(time.Minute))
	if err != nil || len(available) != 1 || available[0].ID != "test-drive-slot" {
		t.Fatalf("available=%+v err=%v", available, err)
	}
	if _, err = repo.CreateAppointmentSlot(ctx, tenant, franchisejourney.AppointmentSlot{ID: "overlap-slot", OrganizationID: "store", Kind: "test-drive", StartsAt: slotStart.Add(15 * time.Minute), EndsAt: slotEnd.Add(15 * time.Minute), Capacity: 1, State: "open", Version: 1}, "unused"); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatalf("overlapping slot accepted: %v", err)
	}
	appointment := franchisejourney.Appointment{ID: "appointment", LeadID: "lead", ModelID: "model", Kind: "test-drive", StartsAt: slotStart, State: "requested", Version: 1}
	created, replayed, err := repo.RequestAppointment(ctx, "journey-api", "store", "appointment-key-1", appointment, "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "018f4d4a-7b36-7a21-8d10-2f4c54c29b02")
	if err != nil || replayed || created.OrganizationID != "store" {
		t.Fatalf("appointment=%+v replay=%v err=%v", created, replayed, err)
	}
	if created.SlotID != "test-drive-slot" || !created.EndsAt.Equal(slotEnd) {
		t.Fatalf("appointment was not server-bound to slot: %+v", created)
	}
	resource, err := repo.CreateServiceResource(ctx, tenant, franchisejourney.ServiceResource{ID: "technician", OrganizationID: "store", PrincipalSubject: "technician-subject", DisplayName: "Technician", Kind: "employee", Status: "active", Version: 1, Skills: []string{"test-drive", "service"}}, "018f4d4a-7b36-7a21-8d10-2f4c54c29b16")
	if err != nil || resource.Version != 1 {
		t.Fatalf("resource=%+v err=%v", resource, err)
	}
	if _, err = repo.CreateServiceResource(ctx, tenant, franchisejourney.ServiceResource{ID: "wrong-skill", OrganizationID: "store", DisplayName: "Service Bay", Kind: "service-bay", Status: "active", Version: 1, Skills: []string{"service"}}, "018f4d4a-7b36-7a21-8d10-2f4c54c29b17"); err != nil {
		t.Fatal(err)
	}
	resourceWorking, err := repo.CreateAvailability(ctx, tenant, "scheduler", franchisejourney.AvailabilityEntry{ID: "technician-working", OrganizationID: "store", ResourceID: "technician", EntryType: "working", StartsAt: slotStart.Add(-time.Hour), EndsAt: workingEnd, State: "active", Version: 1}, "018f4d4a-7b36-7a21-8d10-2f4c54c29b21")
	if err != nil || resourceWorking.ResourceID != "technician" {
		t.Fatalf("resource working window=%+v err=%v", resourceWorking, err)
	}
	if _, err = repo.AssignAppointmentResource(ctx, tenant, "store", "appointment", "wrong-skill", 1, "unused"); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatalf("wrong skill assignment accepted: %v", err)
	}
	assignedAppointment, err := repo.AssignAppointmentResource(ctx, tenant, "store", "appointment", "technician", 1, "018f4d4a-7b36-7a21-8d10-2f4c54c29b18")
	if err != nil || assignedAppointment.ResourceID != "technician" || assignedAppointment.Version != 2 {
		t.Fatalf("assigned appointment=%+v err=%v", assignedAppointment, err)
	}
	if _, err = pool.Exec(ctx, `insert into crm.appointment(tenant_id,appointment_id,organization_id,lead_id,customer_principal_id,model_id,appointment_kind,starts_at,ends_at,slot_id,state,version) values($1,'manual-overlap','store','lead','customer','model','test-drive',$2,$3,'test-drive-slot','requested',1)`, tenant, slotStart, slotEnd); err != nil {
		t.Fatal(err)
	}
	if _, err = repo.AssignAppointmentResource(ctx, tenant, "store", "manual-overlap", "technician", 1, "unused"); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatalf("overlapping resource assignment accepted: %v", err)
	}
	confirmed, err := repo.TransitionAppointment(ctx, tenant, "store", "appointment", "requested", "confirmed", 2, "operator", "", "018f4d4a-7b36-7a21-8d10-2f4c54c29b19")
	if err != nil || confirmed.State != "confirmed" || confirmed.ResourceID != "technician" || confirmed.Version != 3 {
		t.Fatalf("confirmed appointment=%+v err=%v", confirmed, err)
	}
	if _, err = repo.TransitionAppointment(ctx, tenant, "store", "appointment", "requested", "cancelled", 2, "operator", "customer-request", "unused"); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("stale appointment transition accepted")
	}
	if _, err = repo.TransitionAppointment(ctx, tenant, "store", "appointment", "confirmed", "completed", 3, "operator", "", "unused"); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("future appointment completion accepted")
	}
	if _, err = repo.CreateAvailability(ctx, tenant, "scheduler", franchisejourney.AvailabilityEntry{ID: "conflicting-leave", OrganizationID: "store", ResourceID: "technician", EntryType: "unavailable", ReasonCode: "annual-leave", StartsAt: slotStart, EndsAt: slotEnd, State: "active", Version: 1}, "unused"); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatalf("absence overlapping active appointment accepted: %v", err)
	}
	_, replayed, err = repo.RequestAppointment(ctx, "journey-api", "store", "appointment-key-1", appointment, "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "unused")
	if err != nil || !replayed {
		t.Fatalf("replay=%v err=%v", replayed, err)
	}
	if _, _, err = repo.RequestAppointment(ctx, "journey-api", "store", "appointment-key-1", appointment, "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", "unused"); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("changed replay accepted")
	}
	available, err = repo.PublicAppointmentSlots(ctx, "journey-api", "store", "test-drive", slotStart.Add(-time.Minute), slotEnd.Add(time.Minute))
	if err != nil || len(available) != 0 {
		t.Fatalf("full slot remained public: %+v err=%v", available, err)
	}
	second := appointment
	second.ID = "appointment-over-capacity"
	if _, _, err = repo.RequestAppointment(ctx, "journey-api", "store", "appointment-key-over-capacity", second, strings.Repeat("c", 64), "unused"); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatalf("over-capacity booking accepted: %v", err)
	}
	var partialIdempotency int
	if err = pool.QueryRow(ctx, `select count(*) from platform.idempotency_record where tenant_id=$1 and scope='public-appointment' and idempotency_key='appointment-key-over-capacity'`, tenant).Scan(&partialIdempotency); err != nil || partialIdempotency != 0 {
		t.Fatalf("capacity conflict left idempotency residue=%d err=%v", partialIdempotency, err)
	}
	cancelled, err := repo.TransitionAppointment(ctx, tenant, "store", "appointment", "confirmed", "cancelled", 3, "operator", "customer-request", "018f4d4a-7b36-7a21-8d10-2f4c54c29b1a")
	if err != nil || cancelled.State != "cancelled" || cancelled.Version != 4 {
		t.Fatalf("cancelled appointment=%+v err=%v", cancelled, err)
	}
	var transitionActor, transitionReason string
	if err = pool.QueryRow(ctx, `select actor_subject,reason_code from crm.appointment_transition where tenant_id=$1 and appointment_id='appointment' and to_state='cancelled'`, tenant).Scan(&transitionActor, &transitionReason); err != nil || transitionActor != "operator" || transitionReason != "customer-request" {
		t.Fatalf("transition actor=%s reason=%s err=%v", transitionActor, transitionReason, err)
	}
	concurrentStart := slotStart.Add(3 * time.Hour)
	if _, err = repo.CreateAppointmentSlot(ctx, tenant, franchisejourney.AppointmentSlot{ID: "concurrent-slot", OrganizationID: "store", Kind: "service", StartsAt: concurrentStart, EndsAt: concurrentStart.Add(time.Hour), Capacity: 1, State: "open", Version: 1}, "018f4d4a-7b36-7a21-8d10-2f4c54c29b13"); err != nil {
		t.Fatal(err)
	}
	type bookingResult struct {
		created franchisejourney.Appointment
		err     error
	}
	startBookings := make(chan struct{})
	results := make(chan bookingResult, 2)
	for _, input := range []struct{ appointment, key, hash, event string }{
		{"concurrent-a", "appointment-concurrent-a", strings.Repeat("1", 64), "018f4d4a-7b36-7a21-8d10-2f4c54c29b14"},
		{"concurrent-b", "appointment-concurrent-b", strings.Repeat("2", 64), "018f4d4a-7b36-7a21-8d10-2f4c54c29b15"},
	} {
		go func(input struct{ appointment, key, hash, event string }) {
			<-startBookings
			created, _, bookingErr := repo.RequestAppointment(ctx, "journey-api", "store", input.key, franchisejourney.Appointment{ID: input.appointment, LeadID: "lead", ModelID: "model", Kind: "service", StartsAt: concurrentStart, State: "requested", Version: 1}, input.hash, input.event)
			results <- bookingResult{created: created, err: bookingErr}
		}(input)
	}
	close(startBookings)
	successes, conflicts := 0, 0
	for range 2 {
		result := <-results
		if result.err == nil && result.created.SlotID == "concurrent-slot" {
			successes++
		} else if errors.Is(result.err, franchisejourney.ErrConflict) {
			conflicts++
		} else {
			t.Fatalf("unexpected concurrent booking result: %+v", result)
		}
	}
	if successes != 1 || conflicts != 1 {
		t.Fatalf("concurrent capacity successes=%d conflicts=%d", successes, conflicts)
	}
	customerCancelStart := slotStart.Add(6 * time.Hour)
	if _, err = repo.CreateAppointmentSlot(ctx, tenant, franchisejourney.AppointmentSlot{ID: "customer-cancel-slot", OrganizationID: "store", Kind: "service", StartsAt: customerCancelStart, EndsAt: customerCancelStart.Add(time.Hour), Capacity: 1, State: "open", Version: 1}, "018f4d4a-7b36-7a21-8d10-2f4c54c29b22"); err != nil {
		t.Fatal(err)
	}
	customerAppointment := franchisejourney.Appointment{ID: "customer-cancel-appointment", LeadID: "lead", ModelID: "model", Kind: "service", StartsAt: customerCancelStart, State: "requested", Version: 1}
	createdForCustomer, _, err := repo.RequestAppointment(ctx, "journey-api", "store", "appointment-customer-cancel", customerAppointment, strings.Repeat("4", 64), "018f4d4a-7b36-7a21-8d10-2f4c54c29b23")
	if err != nil || createdForCustomer.Version != 1 {
		t.Fatalf("customer appointment=%+v err=%v", createdForCustomer, err)
	}
	customerCancelled, err := repo.CancelCustomerAppointment(ctx, tenant, "store", "customer", "customer-cancel-appointment", 1, "customer-request", "018f4d4a-7b36-7a21-8d10-2f4c54c29b24", "018f4d4a-7b36-7a21-8d10-2f4c54c29b25")
	if err != nil || customerCancelled.State != "cancelled" || customerCancelled.Version != 2 {
		t.Fatalf("customer cancellation=%+v err=%v", customerCancelled, err)
	}
	if _, err = repo.CancelCustomerAppointment(ctx, tenant, "store", "other-customer", "customer-cancel-appointment", 2, "customer-request", "unused", "unused"); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatalf("cross-customer cancellation accepted: %v", err)
	}
	var concurrentBookings int
	if err = pool.QueryRow(ctx, `select count(*) from crm.appointment where tenant_id=$1 and slot_id='concurrent-slot' and state in ('requested','confirmed')`, tenant).Scan(&concurrentBookings); err != nil || concurrentBookings != 1 {
		t.Fatalf("concurrent persisted bookings=%d err=%v", concurrentBookings, err)
	}
	assigned, err := repo.AssignLead(ctx, tenant, "store", "lead", "sales", 1, "018f4d4a-7b36-7a21-8d10-2f4c54c29b03")
	if err != nil || assigned.Version != 2 {
		t.Fatalf("assigned=%+v err=%v", assigned, err)
	}
	if _, err = repo.AssignLead(ctx, tenant, "other", "lead", "intruder", 2, "unused"); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("cross organization assignment succeeded")
	}
	contacted, err := repo.TransitionLead(ctx, tenant, "store", "lead", "new", "contacted", 2, "018f4d4a-7b36-7a21-8d10-2f4c54c29b04")
	if err != nil || contacted.Version != 3 {
		t.Fatalf("contacted=%+v err=%v", contacted, err)
	}
	if _, err = repo.TransitionLead(ctx, tenant, "store", "lead", "new", "lost", 2, "unused"); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("stale transition succeeded")
	}
	quoteHash := strings.Repeat("5", 64)
	quote, replayedQuote, err := repo.CreateQuoteAs(ctx, tenant, "quote-request-0001", franchisejourney.Quote{ID: "quote", OrganizationID: "store", LeadID: "lead", VariantID: "variant", PriceBookID: "retail", ValidUntil: time.Now().Add(24 * time.Hour), State: "issued", Version: 1}, quoteHash, "018f4d4a-7b36-7a21-8d10-2f4c54c29b05", "writer")
	if err != nil || replayedQuote || quote.TotalMinorUnits != 123456 || quote.Currency != "ARS" {
		t.Fatalf("quote=%+v replayed=%v err=%v", quote, replayedQuote, err)
	}

	recoveredQuote, err := repo.QuoteResult(ctx, tenant, "store", "lead", "quote-request-0001")
	if err != nil || recoveredQuote.ID != "quote" || recoveredQuote.TotalMinorUnits != quote.TotalMinorUnits {
		t.Fatalf("quote result=%+v err=%v", recoveredQuote, err)
	}
	for _, scope := range [][3]string{{tenant, "other", "lead"}, {tenant, "store", "other"}, {randomid.Generator{}.New(), "store", "lead"}} {
		if v, e := repo.QuoteResult(ctx, scope[0], scope[1], scope[2], "quote-request-0001"); !errors.Is(e, franchisejourney.ErrConflict) || v.ID != "" {
			t.Fatalf("lookup scope leak=%+v err=%v", v, e)
		}
	}
	if v, e := repo.QuoteResult(ctx, tenant, "store", "lead", "unknown-key-00001"); !errors.Is(e, franchisejourney.ErrConflict) || v.ID != "" {
		t.Fatal("unknown quote identity accepted")
	}
	for _, mutation := range []string{"resource_type='other'", "response_code=200", "response_body='{\"quotation_id\":\"other\"}'"} {
		if _, e := pool.Exec(ctx, "update platform.idempotency_record set "+mutation+" where tenant_id=$1 and scope='franchise-quote' and idempotency_key='quote-request-0001'", tenant); e != nil {
			t.Fatal(e)
		}
		if v, e := repo.QuoteResult(ctx, tenant, "store", "lead", "quote-request-0001"); !errors.Is(e, franchisejourney.ErrConflict) || v.ID != "" {
			t.Fatalf("incoherent result=%+v err=%v", v, e)
		}
		if _, e := pool.Exec(ctx, `update platform.idempotency_record set resource_type='quotation',response_code=201,response_body='{"quotation_id":"quote"}' where tenant_id=$1 and scope='franchise-quote' and idempotency_key='quote-request-0001'`, tenant); e != nil {
			t.Fatal(e)
		}
	}
	replayedValue, replayedQuote, err := repo.CreateQuote(ctx, tenant, "quote-request-0001", franchisejourney.Quote{ID: "unused-replay-id", OrganizationID: "store", LeadID: "lead", VariantID: "variant", PriceBookID: "retail", ValidUntil: quote.ValidUntil, State: "issued", Version: 1}, quoteHash, "unused-replay-event")
	if err != nil || !replayedQuote || replayedValue.ID != quote.ID {
		t.Fatalf("quote replay=%+v replayed=%v err=%v", replayedValue, replayedQuote, err)
	}
	if _, _, err = repo.CreateQuote(ctx, tenant, "quote-request-0001", franchisejourney.Quote{ID: "conflicting-quote", OrganizationID: "store", LeadID: "lead", VariantID: "variant", PriceBookID: "retail", ValidUntil: quote.ValidUntil, State: "issued", Version: 1}, strings.Repeat("6", 64), "unused-conflict-event"); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatalf("same key with different quote request accepted: %v", err)
	}

	var originalActor string
	if _, _, e := repo.CreateQuoteAs(ctx, tenant, "quote-request-0001", quote, quoteHash, "unused", "second-operator"); e != nil {
		t.Fatal(e)
	}
	if e := pool.QueryRow(ctx, `select payload->>'actor_subject' from platform.outbox_event where tenant_id=$1 and aggregate_id='quote' and event_type='quotation.issued'`, tenant).Scan(&originalActor); e != nil || originalActor != "writer" {
		t.Fatalf("replay changed actor=%s err=%v", originalActor, e)
	}
	attempted := quote
	attempted.ID = "quote-audit-rollback"
	if _, _, e := repo.CreateQuoteAs(ctx, tenant, "quote-atomicity-key", attempted, strings.Repeat("a", 64), "018f4d4a-7b36-7a21-8d10-2f4c54c29b05", "writer"); e == nil {
		t.Fatal("duplicate audit event accepted")
	}
	var residue int
	if e := pool.QueryRow(ctx, `select (select count(*) from sales.quotation where tenant_id=$1 and quotation_id='quote-audit-rollback')+(select count(*) from platform.idempotency_record where tenant_id=$1 and scope='franchise-quote' and idempotency_key='quote-atomicity-key')`, tenant).Scan(&residue); e != nil || residue != 0 {
		t.Fatalf("audit rollback residue=%d err=%v", residue, e)
	}
	var quoteRows, quoteEvents int
	if err = pool.QueryRow(ctx, `select count(*) from sales.quotation where tenant_id=$1 and quotation_id=$2`, tenant, quote.ID).Scan(&quoteRows); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select count(*) from platform.outbox_event where tenant_id=$1 and aggregate_id=$2 and event_type='quotation.issued'`, tenant, quote.ID).Scan(&quoteEvents); err != nil || quoteRows != 1 || quoteEvents != 1 {
		t.Fatalf("quote rows=%d events=%d err=%v", quoteRows, quoteEvents, err)
	}
	type quoteResult struct {
		value    franchisejourney.Quote
		replayed bool
		err      error
	}
	startQuotes := make(chan struct{})
	quoteResults := make(chan quoteResult, 2)
	concurrentValidUntil := time.Now().UTC().Add(48 * time.Hour)
	for _, input := range []struct{ id, event string }{
		{"quote-concurrent-a", "018f4d4a-7b36-7a21-8d10-2f4c54c29b40"},
		{"quote-concurrent-b", "018f4d4a-7b36-7a21-8d10-2f4c54c29b41"},
	} {
		go func(input struct{ id, event string }) {
			<-startQuotes
			value, replayed, quoteErr := repo.CreateQuote(ctx, tenant, "quote-concurrent-key", franchisejourney.Quote{ID: input.id, OrganizationID: "store", LeadID: "lead", VariantID: "variant", PriceBookID: "retail", ValidUntil: concurrentValidUntil, State: "issued", Version: 1}, strings.Repeat("9", 64), input.event)
			quoteResults <- quoteResult{value: value, replayed: replayed, err: quoteErr}
		}(input)
	}
	close(startQuotes)
	createdQuotes, replayedQuotes := 0, 0
	var concurrentQuoteID string
	for range 2 {
		result := <-quoteResults
		if result.err != nil {
			t.Fatalf("concurrent quote failed: %v", result.err)
		}
		if concurrentQuoteID == "" {
			concurrentQuoteID = result.value.ID
		} else if result.value.ID != concurrentQuoteID {
			t.Fatalf("concurrent replay returned different ids: %q vs %q", concurrentQuoteID, result.value.ID)
		}
		if result.replayed {
			replayedQuotes++
		} else {
			createdQuotes++
		}
	}
	if createdQuotes != 1 || replayedQuotes != 1 {
		t.Fatalf("concurrent quote created=%d replayed=%d", createdQuotes, replayedQuotes)
	}
	if err = pool.QueryRow(ctx, `select count(*) from sales.quotation where tenant_id=$1 and quotation_id=$2`, tenant, concurrentQuoteID).Scan(&quoteRows); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select count(*) from platform.outbox_event where tenant_id=$1 and aggregate_id=$2 and event_type='quotation.issued'`, tenant, concurrentQuoteID).Scan(&quoteEvents); err != nil || quoteRows != 1 || quoteEvents != 1 {
		t.Fatalf("concurrent quote rows=%d events=%d err=%v", quoteRows, quoteEvents, err)
	}
	accepted, err := repo.AcceptQuote(ctx, tenant, "store", "customer", "quote", 1, "dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd", "accepted-order", "accepted-line", "018f4d4a-7b36-7a21-8d10-2f4c54c29b09", "018f4d4a-7b36-7a21-8d10-2f4c54c29b0a")
	if err != nil || accepted.State != "accepted" || accepted.OrderID != "accepted-order" || accepted.Version != 2 {
		t.Fatalf("accepted=%+v err=%v", accepted, err)
	}
	var orderState, currency string
	var total int64
	if err = pool.QueryRow(ctx, `select state,currency,total_minor_units from sales.customer_order where tenant_id=$1 and order_id='accepted-order'`, tenant).Scan(&orderState, &currency, &total); err != nil || orderState != "placed" || currency != "ARS" || total != 123456 {
		t.Fatalf("order state=%s currency=%s total=%d err=%v", orderState, currency, total, err)
	}
	var linePrice int64
	if err = pool.QueryRow(ctx, `select unit_price_minor_units from sales.customer_order_line where tenant_id=$1 and order_id='accepted-order'`, tenant).Scan(&linePrice); err != nil || linePrice != 123456 {
		t.Fatalf("line price=%d err=%v", linePrice, err)
	}
	if _, err = repo.AcceptQuote(ctx, tenant, "store", "customer", "quote", 1, strings.Repeat("d", 64), "stale-order", "stale-line", "018f4d4a-7b36-7a21-8d10-2f4c54c29b0b", "018f4d4a-7b36-7a21-8d10-2f4c54c29b0c"); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("stale quote acceptance succeeded")
	}
	if _, err = repo.AcceptQuote(ctx, tenant, "store", "other-customer", "quote", 2, strings.Repeat("d", 64), "cross-order", "cross-line", "018f4d4a-7b36-7a21-8d10-2f4c54c29b0d", "018f4d4a-7b36-7a21-8d10-2f4c54c29b0e"); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("cross-customer quote acceptance succeeded")
	}
	if _, err = pool.Exec(ctx, `insert into sales.quotation(tenant_id,quotation_id,organization_id,lead_id,customer_principal_id,variant_id,price_book_id,currency,total_minor_units,valid_until,state,version,created_at) values($1,'expired-quote','store','lead','customer','variant','retail','ARS',123456,clock_timestamp()-interval '1 day','issued',1,clock_timestamp()-interval '2 days')`, tenant); err != nil {
		t.Fatal(err)
	}
	if _, err = repo.AcceptQuote(ctx, tenant, "store", "customer", "expired-quote", 1, strings.Repeat("e", 64), "expired-order", "expired-line", "018f4d4a-7b36-7a21-8d10-2f4c54c29b0f", "018f4d4a-7b36-7a21-8d10-2f4c54c29b10"); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("expired quote acceptance succeeded")
	}
	failedQuote, _, err := repo.CreateQuote(ctx, tenant, "quote-request-0002", franchisejourney.Quote{ID: "rollback-quote", OrganizationID: "store", LeadID: "lead", VariantID: "variant", PriceBookID: "retail", ValidUntil: time.Now().Add(24 * time.Hour), State: "issued", Version: 1}, strings.Repeat("7", 64), "018f4d4a-7b36-7a21-8d10-2f4c54c29b11")
	if err != nil || failedQuote.ID == "" {
		t.Fatal(err)
	}
	_, err = repo.AcceptQuote(ctx, tenant, "store", "customer", "rollback-quote", 1, strings.Repeat("f", 64), "rollback-order", "rollback-line", "018f4d4a-7b36-7a21-8d10-2f4c54c29b12", "018f4d4a-7b36-7a21-8d10-2f4c54c29b11")
	if err == nil {
		t.Fatal("duplicate outbox event did not abort acceptance")
	}
	var rollbackOrders, rollbackAcceptances, rollbackAcceptedEvents int
	if err = pool.QueryRow(ctx, `select count(*) from sales.customer_order where tenant_id=$1 and order_id='rollback-order'`, tenant).Scan(&rollbackOrders); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select count(*) from sales.quotation_acceptance where tenant_id=$1 and quotation_id='rollback-quote'`, tenant).Scan(&rollbackAcceptances); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select count(*) from platform.outbox_event where tenant_id=$1 and aggregate_id='rollback-quote' and event_type='quotation.accepted'`, tenant).Scan(&rollbackAcceptedEvents); err != nil {
		t.Fatal(err)
	}
	var rollbackState string
	if err = pool.QueryRow(ctx, `select state from sales.quotation where tenant_id=$1 and quotation_id='rollback-quote'`, tenant).Scan(&rollbackState); err != nil || rollbackOrders != 0 || rollbackAcceptances != 0 || rollbackAcceptedEvents != 0 || rollbackState != "issued" {
		t.Fatalf("rollback state=%s orders=%d acceptances=%d events=%d err=%v", rollbackState, rollbackOrders, rollbackAcceptances, rollbackAcceptedEvents, err)
	}
	if _, _, err = repo.CreateQuote(ctx, tenant, "quote-request-0003", franchisejourney.Quote{ID: "orphan-quote", OrganizationID: "store", LeadID: "missing", VariantID: "variant", PriceBookID: "retail", ValidUntil: time.Now().Add(24 * time.Hour), State: "issued", Version: 1}, strings.Repeat("8", 64), "018f4d4a-7b36-7a21-8d10-2f4c54c29b08"); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatalf("quote without scoped lead was accepted: %v", err)
	}
	var orphanEvents int
	if err = pool.QueryRow(ctx, `select count(*) from platform.outbox_event where tenant_id=$1 and aggregate_id='orphan-quote'`, tenant).Scan(&orphanEvents); err != nil || orphanEvents != 0 {
		t.Fatalf("orphan quote outbox events=%d err=%v", orphanEvents, err)
	}
	evidence := "cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc"
	checklist, err := repo.PublishDeliveryChecklist(ctx, tenant, "operator", franchisejourney.DeliveryChecklist{ID: "standard-delivery", OrganizationID: "store", Version: 1, Title: "Entrega estándar", State: "published", Items: []franchisejourney.ChecklistItem{{ID: "serial-observed", Ordinal: 1, Prompt: "Verificar serie", ResponseType: "serial", Required: true}, {ID: "asset-condition", Ordinal: 2, Prompt: "Confirmar condición", ResponseType: "confirmation", Required: true}}}, "018f4d4a-7b36-7a21-8d10-2f4c54c29b30")
	if err != nil || checklist.State != "published" {
		t.Fatalf("checklist=%+v err=%v", checklist, err)
	}
	// The existing foreign keys are tenant-scoped, not relational authorization.
	// Corrupt only this synthetic tenant to exercise a stale or inconsistent link.
	for phaseNo, phase := range []string{"complete", "accept", "reject"} {
		for caseNo, mismatch := range []struct{ name, corrupt, restore string }{
			{"order-organization", `update sales.customer_order set organization_id='other' where tenant_id=$1 and order_id='order'`, `update sales.customer_order set organization_id='store' where tenant_id=$1 and order_id='order'`},
			{"order-customer", `update sales.customer_order set customer_principal_id='other-customer' where tenant_id=$1 and order_id='order'`, `update sales.customer_order set customer_principal_id='customer' where tenant_id=$1 and order_id='order'`},
			{"stock-organization", `update inventory.stock_unit set organization_id='other' where tenant_id=$1 and stock_unit_id='stock'`, `update inventory.stock_unit set organization_id='store' where tenant_id=$1 and stock_unit_id='stock'`},
		} {
			t.Run(phase+"-"+mismatch.name, func(t *testing.T) {
				id := "scope-" + phase + "-" + mismatch.name
				if _, err := pool.Exec(ctx, `insert into sales.delivery_handover(tenant_id,handover_id,organization_id,order_id,customer_principal_id,stock_unit_id,state,version) values($1,$2,'store','order','customer','stock','prepared',1)`, tenant, id); err != nil {
					t.Fatal(err)
				}
				responses := []franchisejourney.ChecklistResponse{{ItemID: "serial-observed", ResponseText: "SERIAL-JOURNEY"}, {ItemID: "asset-condition", ResponseText: "confirmed"}}
				event := fmt.Sprintf("018f4d4a-7b36-7a21-8d10-2f4c54c2%04d", 9100+phaseNo*10+caseNo)
				if phase != "complete" {
					if _, err := repo.CompleteDeliveryChecklist(ctx, tenant, "store", "operator", id, 1, checklist.ID, 1, responses, event); err != nil {
						t.Fatal(err)
					}
				}
				if _, err := pool.Exec(ctx, mismatch.corrupt, tenant); err != nil {
					t.Fatal(err)
				}
				defer func() {
					if _, err := pool.Exec(ctx, mismatch.restore, tenant); err != nil {
						t.Fatal(err)
					}
				}()
				event = fmt.Sprintf("018f4d4a-7b36-7a21-8d10-2f4c54c2%04d", 9200+phaseNo*10+caseNo)
				if phase == "complete" {
					read, readErr := repo.CustomerJourney(ctx, tenant, "store", "customer")
					if !errors.Is(readErr, franchisejourney.ErrConflict) || len(read.Appointments)+len(read.Quotes)+len(read.Handovers)+len(read.Exceptions) != 0 {
						t.Errorf("inconsistent customer projection returned data: err=%v handovers=%d", readErr, len(read.Handovers))
					}
					unrelated, unrelatedErr := repo.CustomerJourney(ctx, tenant, "store", "unrelated-customer")
					if unrelatedErr != nil || len(unrelated.Handovers) != 0 {
						t.Errorf("unrelated customer affected: %v", unrelatedErr)
					}
					ops, opsErr := NewCommerce(pool).Operations(ctx, tenant, "store")
					if mismatch.name != "order-organization" {
						if opsErr == nil || len(ops.Orders)+len(ops.Stock) != 0 {
							t.Errorf("inconsistent operation projection returned data: err=%v orders=%d", opsErr, len(ops.Orders))
						}
					} else {
						for _, order := range ops.Orders {
							if order.ID == "order" {
								t.Error("foreign organization order leaked")
							}
						}
					}
				}
				var actionErr error
				switch phase {
				case "complete":
					_, actionErr = repo.CompleteDeliveryChecklist(ctx, tenant, "store", "operator", id, 1, checklist.ID, 1, responses, event)
				case "accept":
					_, actionErr = repo.AcceptHandover(ctx, tenant, "store", "customer", id, 2, "SERIAL-JOURNEY", checklist.ID, 1, evidence, event)
				case "reject":
					_, actionErr = repo.RejectHandover(ctx, tenant, "store", "customer", id, 2, "visible-damage", "synthetic", evidence, id+"-exception", event)
				}
				if !errors.Is(actionErr, franchisejourney.ErrConflict) {
					t.Errorf("inconsistent %s allowed %s: %v", mismatch.name, phase, actionErr)
				}
				var state string
				var events, responseCount int
				if err := pool.QueryRow(ctx, `select state,(select count(*) from platform.outbox_event where tenant_id=$1 and event_id=$3),(select count(*) from sales.delivery_checklist_response where tenant_id=$1 and handover_id=$2) from sales.delivery_handover where tenant_id=$1 and handover_id=$2`, tenant, id, event).Scan(&state, &events, &responseCount); err != nil {
					t.Fatal(err)
				}
				wantState, wantResponses := "presented", 2
				if phase == "complete" {
					wantState, wantResponses = "prepared", 0
				}
				if state != wantState || events != 0 || responseCount != wantResponses {
					t.Errorf("partial effect state=%s events=%d responses=%d", state, events, responseCount)
				}
			})
		}
	}
	if _, err = repo.AcceptHandover(ctx, tenant, "store", "customer", "handover", 1, "SERIAL-JOURNEY", checklist.ID, checklist.Version, evidence, "unused"); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatalf("handover accepted before checklist completion: %v", err)
	}
	if _, err = repo.CompleteDeliveryChecklist(ctx, tenant, "store", "operator", "handover", 1, checklist.ID, checklist.Version, []franchisejourney.ChecklistResponse{{ItemID: "serial-observed", ResponseText: "SERIAL-JOURNEY"}}, "unused"); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatalf("incomplete required checklist accepted: %v", err)
	}
	completed, err := repo.CompleteDeliveryChecklist(ctx, tenant, "store", "operator", "handover", 1, checklist.ID, checklist.Version, []franchisejourney.ChecklistResponse{{ItemID: "serial-observed", ResponseText: "SERIAL-JOURNEY"}, {ItemID: "asset-condition", ResponseText: "confirmed"}}, "018f4d4a-7b36-7a21-8d10-2f4c54c29b31")
	if err != nil || completed.State != "presented" || completed.Version != 2 {
		t.Fatalf("completed=%+v err=%v", completed, err)
	}
	if _, err = repo.AcceptHandover(ctx, tenant, "store", "customer", "handover", 2, "WRONG-SERIAL", checklist.ID, checklist.Version, evidence, "unused"); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatalf("handover accepted with wrong serial: %v", err)
	}
	if _, err = repo.AcceptHandover(ctx, tenant, "store", "customer", "handover", 2, "SERIAL-JOURNEY", checklist.ID, 2, evidence, "unused"); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatalf("handover accepted against wrong checklist version: %v", err)
	}
	var acceptanceWG sync.WaitGroup
	// Audit failure must roll back acceptance, leaving the same version usable.
	if _, err := repo.AcceptHandover(ctx, tenant, "store", "customer", "handover", 2, "SERIAL-JOURNEY", checklist.ID, checklist.Version, evidence, "018f4d4a-7b36-7a21-8d10-2f4c54c29b31"); err == nil {
		t.Fatal("duplicate outbox event allowed acceptance")
	}
	var recoverState string
	var recoverVersion int64
	if err := pool.QueryRow(ctx, `select state,version from sales.delivery_handover where tenant_id=$1 and handover_id='handover'`, tenant).Scan(&recoverState, &recoverVersion); err != nil || recoverState != "presented" || recoverVersion != 2 {
		t.Fatalf("acceptance rollback state=%s version=%d err=%v", recoverState, recoverVersion, err)
	}
	// Check the actual row locks, not timing of an assumed concurrent writer.
	for _, mutation := range []string{`update sales.customer_order set customer_principal_id='other-customer' where tenant_id=$1 and order_id='order'`, `update inventory.stock_unit set organization_id='other' where tenant_id=$1 and stock_unit_id='stock'`} {
		locked, err := pool.Begin(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if err = lockDeliveryScope(ctx, locked, tenant, "store", "handover"); err != nil {
			_ = locked.Rollback(ctx)
			t.Fatal(err)
		}
		writer, err := pool.Begin(ctx)
		if err != nil {
			_ = locked.Rollback(ctx)
			t.Fatal(err)
		}
		_, err = writer.Exec(ctx, `set local lock_timeout='50ms'`)
		if err != nil {
			_ = writer.Rollback(ctx)
			_ = locked.Rollback(ctx)
			t.Fatal(err)
		}
		_, err = writer.Exec(ctx, mutation, tenant)
		_ = writer.Rollback(ctx)
		_ = locked.Rollback(ctx)
		var pgErr *pgconn.PgError
		if !errors.As(err, &pgErr) || pgErr.Code != "55P03" {
			t.Fatalf("linked row not fenced by transaction: %v", err)
		}
	}
	acceptanceResults := make(chan error, 16)
	for range 16 {
		acceptanceWG.Add(1)
		go func() {
			defer acceptanceWG.Done()
			value, err := repo.AcceptHandover(ctx, tenant, "store", "customer", "handover", 2, "SERIAL-JOURNEY", checklist.ID, checklist.Version, evidence, "018f4d4a-7b36-7a21-8d10-2f4c54c29b06")
			if err == nil && (value.State != "accepted" || value.Version != 3) {
				err = fmt.Errorf("invalid receipt state/version")
			}
			acceptanceResults <- err
		}()
	}
	acceptanceWG.Wait()
	close(acceptanceResults)
	acceptedCount, conflictCount := 0, 0
	for err := range acceptanceResults {
		if err == nil {
			acceptedCount++
		} else if errors.Is(err, franchisejourney.ErrConflict) {
			conflictCount++
		} else {
			t.Fatal(err)
		}
	}
	if acceptedCount != 1 || conflictCount != 15 {
		t.Fatalf("concurrent acceptance successes=%d conflicts=%d", acceptedCount, conflictCount)
	}
	var acceptedEvents int
	if err := pool.QueryRow(ctx, `select count(*) from platform.outbox_event where tenant_id=$1 and aggregate_id='handover' and event_type='delivery-handover.accepted'`, tenant).Scan(&acceptedEvents); err != nil || acceptedEvents != 1 {
		t.Fatalf("acceptance events=%d err=%v", acceptedEvents, err)
	}
	if _, err = repo.AcceptHandover(ctx, tenant, "store", "other-customer", "handover", 3, "SERIAL-JOURNEY", checklist.ID, checklist.Version, evidence, "unused"); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("handover crossed customer scope")
	}
	rejectedReady, err := repo.CompleteDeliveryChecklist(ctx, tenant, "store", "operator", "handover-reject", 1, checklist.ID, checklist.Version, []franchisejourney.ChecklistResponse{{ItemID: "serial-observed", ResponseText: "SERIAL-JOURNEY"}, {ItemID: "asset-condition", ResponseText: "confirmed"}}, "018f4d4a-7b36-7a21-8d10-2f4c54c29b32")
	if err != nil || rejectedReady.State != "presented" {
		t.Fatalf("rejected-ready=%+v err=%v", rejectedReady, err)
	}
	if _, err = repo.RejectHandover(ctx, tenant, "store", "other-customer", "handover-reject", 2, "visible-damage", "Rayón visible", evidence, "unused", "unused"); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("handover rejection crossed customer scope")
	}
	exception, err := repo.RejectHandover(ctx, tenant, "store", "customer", "handover-reject", 2, "visible-damage", "Rayón visible", evidence, "delivery-exception", "018f4d4a-7b36-7a21-8d10-2f4c54c29b33")
	if err != nil || exception.State != "open" || exception.Version != 1 {
		t.Fatalf("exception=%+v err=%v", exception, err)
	}
	openExceptions, err := repo.DeliveryExceptions(ctx, tenant, "store", 100)
	if err != nil || len(openExceptions) != 1 || openExceptions[0].ID != exception.ID {
		t.Fatalf("open exceptions=%+v err=%v", openExceptions, err)
	}
	for caseNo, mismatch := range []struct{ name, corrupt, restore string }{
		{"order-organization", `update sales.customer_order set organization_id='other' where tenant_id=$1 and order_id='order'`, `update sales.customer_order set organization_id='store' where tenant_id=$1 and order_id='order'`},
		{"order-customer", `update sales.customer_order set customer_principal_id='other-customer' where tenant_id=$1 and order_id='order'`, `update sales.customer_order set customer_principal_id='customer' where tenant_id=$1 and order_id='order'`},
		{"stock-organization", `update inventory.stock_unit set organization_id='other' where tenant_id=$1 and stock_unit_id='stock'`, `update inventory.stock_unit set organization_id='store' where tenant_id=$1 and stock_unit_id='stock'`},
	} {
		t.Run("resolve-"+mismatch.name, func(t *testing.T) {
			if _, err := pool.Exec(ctx, mismatch.corrupt, tenant); err != nil {
				t.Fatal(err)
			}
			defer func() {
				if _, err := pool.Exec(ctx, mismatch.restore, tenant); err != nil {
					t.Fatal(err)
				}
			}()
			_, err := repo.ResolveDeliveryException(ctx, tenant, "store", "operator", exception.ID, 1, "correct-and-represent", "synthetic", "bad-successor", "unused", "bad-resolution", fmt.Sprintf("018f4d4a-7b36-7a21-8d10-2f4c54c2%04d", 9300+caseNo))
			if !errors.Is(err, franchisejourney.ErrConflict) {
				t.Errorf("inconsistent resolution allowed: %v", err)
			}
			var state string
			var effects int
			if err := pool.QueryRow(ctx, `select state,(select count(*) from sales.delivery_exception_resolution where tenant_id=$1 and exception_id=$2) from sales.delivery_exception where tenant_id=$1 and exception_id=$2`, tenant, exception.ID).Scan(&state, &effects); err != nil || state != "open" || effects != 0 {
				t.Errorf("resolution state=%s effects=%d err=%v", state, effects, err)
			}
		})
	}
	resolution, err := repo.ResolveDeliveryException(ctx, tenant, "store", "operator", exception.ID, 1, "correct-and-represent", "Corregir preparación", "handover-successor", "unused-authorization", "delivery-resolution", "018f4d4a-7b36-7a21-8d10-2f4c54c29b34")
	if err != nil || resolution.Exception.State != "resolved" || resolution.SuccessorHandover == nil || resolution.SuccessorHandover.State != "prepared" || resolution.SuccessorHandover.SupersedesHandoverID != "handover-reject" {
		t.Fatalf("resolution=%+v err=%v", resolution, err)
	}
	if _, err = repo.ResolveDeliveryException(ctx, tenant, "store", "operator", exception.ID, 1, "return", "retry", "unused", "unused", "unused", "unused"); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("resolved delivery exception was processed twice")
	}
	returnReady, err := repo.CompleteDeliveryChecklist(ctx, tenant, "store", "operator", "handover-return", 1, checklist.ID, checklist.Version, []franchisejourney.ChecklistResponse{{ItemID: "serial-observed", ResponseText: "SERIAL-JOURNEY"}, {ItemID: "asset-condition", ResponseText: "confirmed"}}, "018f4d4a-7b36-7a21-8d10-2f4c54c29b35")
	if err != nil || returnReady.State != "presented" {
		t.Fatalf("return-ready=%+v err=%v", returnReady, err)
	}
	returnException, err := repo.RejectHandover(ctx, tenant, "store", "customer", "handover-return", 2, "customer-return", "Cliente solicita devolución", evidence, "delivery-return-exception", "018f4d4a-7b36-7a21-8d10-2f4c54c29b36")
	if err != nil || returnException.State != "open" {
		t.Fatalf("return exception=%+v err=%v", returnException, err)
	}
	returnResolution, err := repo.ResolveDeliveryException(ctx, tenant, "store", "operator", returnException.ID, 1, "return", "Autorizar recepción y revisión", "unused-successor", "return-authorization", "delivery-return-resolution", "018f4d4a-7b36-7a21-8d10-2f4c54c29b37")
	if err != nil || returnResolution.Exception.State != "resolved" || returnResolution.ReturnAuthorizationID != "return-authorization" || returnResolution.Disposition != "return" || returnResolution.SuccessorHandover != nil {
		t.Fatalf("return resolution=%+v err=%v", returnResolution, err)
	}
	var authorizationOrder, authorizationStock, exactCostOrder, authorizationState string
	if err = pool.QueryRow(ctx, `select order_id,stock_unit_id,exact_cost_source_order_id,state from sales.return_authorization where tenant_id=$1 and authorization_id='return-authorization'`, tenant).Scan(&authorizationOrder, &authorizationStock, &exactCostOrder, &authorizationState); err != nil {
		t.Fatal(err)
	}
	if authorizationOrder != "order" || authorizationStock != "stock" || exactCostOrder != "order" || authorizationState != "authorized" {
		t.Fatalf("return authorization order=%q stock=%q exact_cost_order=%q state=%q", authorizationOrder, authorizationStock, exactCostOrder, authorizationState)
	}
	if _, err = pool.Exec(ctx, `update sales.return_authorization set state='received' where tenant_id=$1 and authorization_id='return-authorization'`, tenant); err == nil {
		t.Fatal("return authorization accepted an in-place mutation")
	}
	if _, err = repo.ReceiveReturn(ctx, tenant, "store", "operator", "return-authorization", "WRONG-SERIAL", "damaged", "Daño confirmado", evidence, "return-receipt-invalid", "unused"); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatalf("return receipt accepted wrong durable serial: %v", err)
	}
	receipt, err := repo.ReceiveReturn(ctx, tenant, "store", "operator", "return-authorization", "SERIAL-JOURNEY", "damaged", "Daño confirmado", evidence, "return-receipt", "018f4d4a-7b36-7a21-8d10-2f4c54c29b38")
	if err != nil || receipt.AuthorizationID != "return-authorization" || receipt.OrderID != "order" || receipt.StockUnitID != "stock" || receipt.ConditionCode != "damaged" {
		t.Fatalf("return receipt=%+v err=%v", receipt, err)
	}
	if _, err = repo.ReceiveReturn(ctx, tenant, "store", "operator", "return-authorization", "SERIAL-JOURNEY", "damaged", "duplicate", evidence, "return-receipt-duplicate", "unused"); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("return authorization was received twice")
	}
	disposition, err := repo.DecideReturn(ctx, tenant, "store", "operator", receipt.ID, "quarantine", "Separar y solicitar efectos", "return-disposition", "return-effect-inventory", "return-effect-remedy", "return-effect-accounting", "return-effect-fiscal", "018f4d4a-7b36-7a21-8d10-2f4c54c29b39")
	if err != nil || disposition.CustomerRemedy != "refund" || disposition.InventoryAction != "quarantine" || len(disposition.Effects) != 4 {
		t.Fatalf("return disposition=%+v err=%v", disposition, err)
	}
	owners := map[string]string{}
	for _, effect := range disposition.Effects {
		owners[effect.EffectKind] = effect.OwnerContext
	}
	if owners["inventory"] != "inventory" || owners["refund"] != "payment" || owners["accounting"] != "accounting" || owners["fiscal"] != "fiscal" {
		t.Fatalf("return effect owners=%+v", owners)
	}
	if _, err = repo.DecideReturn(ctx, tenant, "store", "operator", receipt.ID, "restock", "duplicate", "unused", "unused", "unused", "unused", "unused", "unused"); !errors.Is(err, franchisejourney.ErrConflict) {
		t.Fatal("return receipt received two dispositions")
	}
	returnCases, err := repo.ReturnCases(ctx, tenant, "store", 100)
	if err != nil || len(returnCases) != 1 || returnCases[0].Receipt == nil || returnCases[0].Disposition == nil || len(returnCases[0].Disposition.Effects) != 4 {
		t.Fatalf("return cases=%+v err=%v", returnCases, err)
	}
	if _, err = pool.Exec(ctx, `update sales.return_receipt set condition_code='sealed' where tenant_id=$1 and receipt_id='return-receipt'`, tenant); err == nil {
		t.Fatal("return receipt accepted in-place mutation")
	}
	journey, err := repo.CustomerJourney(ctx, tenant, "store", "customer")
	acceptedQuotes := 0
	for _, item := range journey.Quotes {
		if item.ID == "quote" && item.State == "accepted" && item.OrderID == "accepted-order" {
			acceptedQuotes++
		}
	}
	successors := 0
	for _, item := range journey.Handovers {
		if item.SupersedesHandoverID == "handover-reject" && item.State == "prepared" {
			successors++
		}
	}
	negativeHandovers := 0
	for _, h := range journey.Handovers {
		if strings.HasPrefix(h.ID, "scope-") {
			negativeHandovers++
		}
	}
	if err != nil || len(journey.Appointments) != 4 || len(journey.Quotes) != 4 || acceptedQuotes != 1 || len(journey.Handovers) != 13 || negativeHandovers != 9 || successors != 1 || len(journey.Exceptions) != 2 || journey.Exceptions[0].State != "resolved" || journey.Exceptions[1].State != "resolved" {
		t.Fatalf("journey=%+v err=%v", journey, err)
	}
	fresh, err := pgxpool.New(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer fresh.Close()
	var durableState string
	var durableVersion int64
	var durableEvents int
	if err := fresh.QueryRow(ctx, `select state,version,(select count(*) from platform.outbox_event where tenant_id=$1 and aggregate_id='handover' and event_type='delivery-handover.accepted') from sales.delivery_handover where tenant_id=$1 and handover_id='handover'`, tenant).Scan(&durableState, &durableVersion, &durableEvents); err != nil || durableState != "accepted" || durableVersion != 3 || durableEvents != 1 {
		t.Fatalf("fresh pool state=%s version=%d events=%d err=%v", durableState, durableVersion, durableEvents, err)
	}
	// Keep immutable history intact; this is a separate, FK-valid bad import.
	probe := &journeyReadProbe{afterFirstRead: func() error {
		_, err := pool.Exec(ctx, `update sales.customer_order set customer_principal_id='other-customer' where tenant_id=$1 and order_id='order'`, tenant)
		return err
	}}
	probeConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		t.Fatal(err)
	}
	probeConfig.ConnConfig.Tracer = probe
	probePool, err := pgxpool.NewWithConfig(ctx, probeConfig)
	if err != nil {
		t.Fatal(err)
	}
	defer probePool.Close()
	consistent, consistentErr := NewFranchiseJourney(probePool).CustomerJourney(ctx, tenant, "store", "customer")
	changed, changedErr := repo.CustomerJourney(ctx, tenant, "store", "customer")
	if _, err := pool.Exec(ctx, `update sales.customer_order set customer_principal_id='customer' where tenant_id=$1 and order_id='order'`, tenant); err != nil {
		t.Fatal(err)
	}
	if consistentErr != nil || len(consistent.Handovers) != 13 || probe.err != nil || !probe.inTransaction || probe.boundedItems != 13 {
		t.Fatalf("snapshot/bound failed read=%v writer=%v handovers=%d tx=%v items=%d", consistentErr, probe.err, len(consistent.Handovers), probe.inTransaction, probe.boundedItems)
	}
	if !errors.Is(changedErr, franchisejourney.ErrConflict) || len(changed.Handovers) != 0 {
		t.Fatalf("new snapshot did not reject changed relation: %v", changedErr)
	}
	recovered, recoveredErr := repo.CustomerJourney(ctx, tenant, "store", "customer")
	if recoveredErr != nil || len(recovered.Handovers) != 13 {
		t.Fatalf("read recovery failed: %v", recoveredErr)
	}
	for _, q := range []string{
		`insert into crm.customer_profile(tenant_id,customer_principal_id,display_name,email_normalized)values($1,'unrelated-customer','Synthetic','unrelated@example.test')`,
		`insert into sales.delivery_exception(tenant_id,exception_id,organization_id,handover_id,customer_principal_id,reason_code,details,rejection_evidence_sha256_hex,state,version) values($1,'bad-import-exception','store','handover-successor','unrelated-customer','synthetic-case','PRIVATE-SYNTHETIC-NOTE',repeat('a',64),'open',1)`,
	} {
		if _, err := pool.Exec(ctx, q, tenant); err != nil {
			t.Fatal(err)
		}
	}
	wrongJourney, wrongJourneyErr := repo.CustomerJourney(ctx, tenant, "store", "unrelated-customer")
	wrongExceptions, wrongExceptionsErr := repo.DeliveryExceptions(ctx, tenant, "store", 100)
	if !errors.Is(wrongJourneyErr, franchisejourney.ErrConflict) || len(wrongJourney.Exceptions) != 0 {
		t.Errorf("foreign exception exposed to customer: %v", wrongJourneyErr)
	}
	if !errors.Is(wrongExceptionsErr, franchisejourney.ErrConflict) || len(wrongExceptions) != 0 {
		t.Errorf("incoherent exception exposed to operator: %v", wrongExceptionsErr)
	}
}

func TestAvailabilityCreationIdentityAndRecovery(t *testing.T) {
	raw := os.Getenv("TEST_DATABASE_URL")
	if raw == "" {
		t.Skip("disposable database not selected")
	}
	cfg, err := pgxpool.ParseConfig(raw)
	if err != nil || cfg.ConnConfig.Host != "127.0.0.1" || !strings.HasPrefix(cfg.ConnConfig.Database, "elite_confirmation_") {
		t.Fatal("requires disposable loopback database")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := (randomid.Generator{}).New()
	for _, sql := range []string{`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1::uuid,'create-'||replace(($1::uuid)::text,'-',''),'Synthetic','Synthetic')`, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'store','store','Synthetic','store')`} {
		if _, err = pool.Exec(ctx, sql, tenant); err != nil {
			t.Fatal(err)
		}
	}
	repo := NewFranchiseJourney(pool)
	start := time.Now().UTC().Add(72 * time.Hour).Truncate(time.Second)
	input := franchisejourney.AvailabilityEntry{ID: "created-interval", OrganizationID: "store", EntryType: "working", StartsAt: start, EndsAt: start.Add(time.Hour)}
	const key = "availability-concurrent-key"
	hash := strings.Repeat("a", 64)
	type outcome struct {
		v      franchisejourney.AvailabilityEntry
		replay bool
		err    error
	}
	results := make(chan outcome, 8)
	begin := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-begin
			v := input
			v.ID = fmt.Sprintf("created-%d", i)
			got, replay, e := repo.CreateAvailabilityOnce(ctx, tenant, "scheduler", key, hash, v, (randomid.Generator{}).New())
			results <- outcome{got, replay, e}
		}(i)
	}
	close(begin)
	wg.Wait()
	close(results)
	created, replays := 0, 0
	var id string
	for r := range results {
		if r.err != nil {
			t.Fatal(r.err)
		}
		if id != "" && id != r.v.ID {
			t.Fatal("duplicate identity")
		}
		id = r.v.ID
		if r.replay {
			replays++
		} else {
			created++
		}
	}
	if created != 1 || replays != 7 {
		t.Fatalf("created=%d replay=%d", created, replays)
	}
	recovered, e := repo.AvailabilityCreationResult(ctx, tenant, "store", key)
	if e != nil || recovered.ID != id {
		t.Fatalf("recover=%+v err=%v", recovered, e)
	}
	for _, scope := range [][3]string{{tenant, "other", key}, {tenant, "store", "unknown-availability-key"}, {(randomid.Generator{}).New(), "store", key}} {
		if v, e := repo.AvailabilityCreationResult(ctx, scope[0], scope[1], scope[2]); !errors.Is(e, franchisejourney.ErrConflict) || v.ID != "" {
			t.Fatalf("scope leak=%+v %v", v, e)
		}
	}
	if _, _, e := repo.CreateAvailabilityOnce(ctx, tenant, "scheduler", key, strings.Repeat("b", 64), input, (randomid.Generator{}).New()); !errors.Is(e, franchisejourney.ErrConflict) {
		t.Fatal("divergent intent accepted")
	}
	var event, actor string
	var rows, keys, events int
	if e := pool.QueryRow(ctx, `select event_id::text,payload->>'actor_subject' from platform.outbox_event where tenant_id=$1 and aggregate_id=$2 and event_type='availability-entry.created'`, tenant, id).Scan(&event, &actor); e != nil || actor != "scheduler" {
		t.Fatalf("actor=%s err=%v", actor, e)
	}
	if _, _, e := repo.CreateAvailabilityOnce(ctx, tenant, "another-operator", key, hash, input, (randomid.Generator{}).New()); e != nil {
		t.Fatal(e)
	}
	next := input
	next.ID = "rolled-back"
	next.StartsAt = start.Add(48 * time.Hour)
	next.EndsAt = next.StartsAt.Add(time.Hour)
	if v, _, e := repo.CreateAvailabilityOnce(ctx, tenant, "scheduler", "availability-rollback-key", hash, next, event); e == nil || v.ID != "" {
		t.Fatalf("audit failure not atomic=%+v %v", v, e)
	}
	if e := pool.QueryRow(ctx, `select (select count(*) from crm.availability_entry where tenant_id=$1),(select count(*) from platform.idempotency_record where tenant_id=$1 and scope='franchise-availability'),(select count(*) from platform.outbox_event where tenant_id=$1 and event_type='availability-entry.created')`, tenant).Scan(&rows, &keys, &events); e != nil || rows != 1 || keys != 1 || events != 1 {
		t.Fatalf("effects=%d/%d/%d err=%v", rows, keys, events, e)
	}
	if _, e := repo.CancelAvailability(ctx, tenant, "store", id, 1, "scheduler", "schedule-correction", (randomid.Generator{}).New()); e != nil {
		t.Fatal(e)
	}
	after, replay, e := repo.CreateAvailabilityOnce(ctx, tenant, "another-operator", key, hash, input, (randomid.Generator{}).New())
	if e != nil || !replay || after.ID != id || after.State != "cancelled" || after.Version != 2 {
		t.Fatalf("replay resurrected interval=%+v replay=%v err=%v", after, replay, e)
	}
	fresh, e := pgxpool.NewWithConfig(ctx, cfg.Copy())
	if e != nil {
		t.Fatal(e)
	}
	defer fresh.Close()
	after, e = NewFranchiseJourney(fresh).AvailabilityCreationResult(ctx, tenant, "store", key)
	if e != nil || after.State != "cancelled" {
		t.Fatalf("fresh pool recovery=%+v err=%v", after, e)
	}
	for _, mutation := range []string{"resource_type='other'", "response_code=200", "response_body='{\"availability_id\":\"other\"}'"} {
		if _, e := pool.Exec(ctx, "update platform.idempotency_record set "+mutation+" where tenant_id=$1 and scope='franchise-availability'", tenant); e != nil {
			t.Fatal(e)
		}
		if v, e := repo.AvailabilityCreationResult(ctx, tenant, "store", key); !errors.Is(e, franchisejourney.ErrConflict) || v.ID != "" {
			t.Fatal("incoherent receipt accepted")
		}
		if _, e := pool.Exec(ctx, `update platform.idempotency_record set resource_type='availability-entry',response_code=201,response_body=jsonb_build_object('availability_id',$2::text) where tenant_id=$1 and scope='franchise-availability'`, tenant, id); e != nil {
			t.Fatal(e)
		}
	}
	if e := pool.QueryRow(ctx, `select payload->>'actor_subject' from platform.outbox_event where tenant_id=$1 and event_id=$2::uuid`, tenant, event).Scan(&actor); e != nil || actor != "scheduler" {
		t.Fatal("replay overwrote creator")
	}
	t.Log("AVAILABILITY_IDENTITY_PASS create=1 replay=7 actor_preserved cancelled_replay rollback scoped_lookup corruption fresh_pool")
}

func TestResourceCreationIdentityAndRecovery(t *testing.T) {
	raw := os.Getenv("TEST_DATABASE_URL")
	if raw == "" {
		t.Skip("disposable database not selected")
	}
	cfg, err := pgxpool.ParseConfig(raw)
	if err != nil || cfg.ConnConfig.Host != "127.0.0.1" || !strings.HasPrefix(cfg.ConnConfig.Database, "elite_confirmation_") {
		t.Fatal("requires disposable loopback database")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := (randomid.Generator{}).New()
	for _, sql := range []string{`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1::uuid,'create-'||replace(($1::uuid)::text,'-',''),'Synthetic','Synthetic')`, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'store','store','Synthetic','store')`} {
		if _, err = pool.Exec(ctx, sql, tenant); err != nil {
			t.Fatal(err)
		}
	}
	repo := NewFranchiseJourney(pool)
	input := franchisejourney.ServiceResource{ID: "created-resource", OrganizationID: "store", DisplayName: "Synthetic bay", Kind: "service-bay", Status: "active", Version: 1, Skills: []string{"service", "delivery"}}
	const key = "resource-concurrent-key"
	hash := strings.Repeat("a", 64)
	type outcome struct {
		v      franchisejourney.ServiceResource
		replay bool
		err    error
	}
	results := make(chan outcome, 8)
	begin := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-begin
			v := input
			v.ID = fmt.Sprintf("created-%d", i)
			got, replay, e := repo.CreateServiceResourceOnce(ctx, tenant, "resource-manager", key, hash, v, (randomid.Generator{}).New())
			results <- outcome{got, replay, e}
		}(i)
	}
	close(begin)
	wg.Wait()
	close(results)
	created, replays := 0, 0
	var id string
	for r := range results {
		if r.err != nil {
			t.Fatal(r.err)
		}
		if id != "" && id != r.v.ID {
			t.Fatal("duplicate identity")
		}
		id = r.v.ID
		if r.replay {
			replays++
		} else {
			created++
		}
	}
	if created != 1 || replays != 7 {
		t.Fatalf("created=%d replay=%d", created, replays)
	}
	recovered, e := repo.ServiceResourceCreationResult(ctx, tenant, "store", key)
	if e != nil || recovered.ID != id {
		t.Fatalf("recover=%+v err=%v", recovered, e)
	}
	for _, scope := range [][3]string{{tenant, "other", key}, {tenant, "store", "unknown-resource-key"}, {(randomid.Generator{}).New(), "store", key}} {
		if v, e := repo.ServiceResourceCreationResult(ctx, scope[0], scope[1], scope[2]); !errors.Is(e, franchisejourney.ErrConflict) || v.ID != "" {
			t.Fatalf("scope leak=%+v %v", v, e)
		}
	}
	if _, _, e := repo.CreateServiceResourceOnce(ctx, tenant, "resource-manager", key, strings.Repeat("b", 64), input, (randomid.Generator{}).New()); !errors.Is(e, franchisejourney.ErrConflict) {
		t.Fatal("divergent intent accepted")
	}
	var event, actor string
	var rows, keys, events int
	if e := pool.QueryRow(ctx, `select event_id::text,payload->>'actor_subject' from platform.outbox_event where tenant_id=$1 and aggregate_id=$2 and event_type='service-resource.created'`, tenant, id).Scan(&event, &actor); e != nil || actor != "resource-manager" {
		t.Fatalf("actor=%s err=%v", actor, e)
	}
	if _, _, e := repo.CreateServiceResourceOnce(ctx, tenant, "another-operator", key, hash, input, (randomid.Generator{}).New()); e != nil {
		t.Fatal(e)
	}
	next := input
	next.ID = "rolled-back"
	if v, _, e := repo.CreateServiceResourceOnce(ctx, tenant, "resource-manager", "resource-rollback-key", hash, next, event); e == nil || v.ID != "" {
		t.Fatalf("audit failure not atomic=%+v %v", v, e)
	}
	if e := pool.QueryRow(ctx, `select (select count(*) from crm.service_resource where tenant_id=$1),(select count(*) from platform.idempotency_record where tenant_id=$1 and scope='franchise-resource'),(select count(*) from platform.outbox_event where tenant_id=$1 and event_type='service-resource.created')`, tenant).Scan(&rows, &keys, &events); e != nil || rows != 1 || keys != 1 || events != 1 {
		t.Fatalf("effects=%d/%d/%d err=%v", rows, keys, events, e)
	}
	if _, e := pool.Exec(ctx, `update crm.service_resource set status='inactive',version=2 where tenant_id=$1 and resource_id=$2`, tenant, id); e != nil {
		t.Fatal(e)
	}
	after, replay, e := repo.CreateServiceResourceOnce(ctx, tenant, "another-operator", key, hash, input, (randomid.Generator{}).New())
	if e != nil || !replay || after.ID != id || after.Status != "inactive" || after.Version != 2 {
		t.Fatalf("replay reactivated resource=%+v replay=%v err=%v", after, replay, e)
	}
	fresh, e := pgxpool.NewWithConfig(ctx, cfg.Copy())
	if e != nil {
		t.Fatal(e)
	}
	defer fresh.Close()
	after, e = NewFranchiseJourney(fresh).ServiceResourceCreationResult(ctx, tenant, "store", key)
	if e != nil || after.Status != "inactive" {
		t.Fatalf("fresh pool recovery=%+v err=%v", after, e)
	}
	for _, mutation := range []string{"resource_type='other'", "response_code=200", "response_body='{\"resource_id\":\"other\"}'"} {
		if _, e := pool.Exec(ctx, "update platform.idempotency_record set "+mutation+" where tenant_id=$1 and scope='franchise-resource'", tenant); e != nil {
			t.Fatal(e)
		}
		if v, e := repo.ServiceResourceCreationResult(ctx, tenant, "store", key); !errors.Is(e, franchisejourney.ErrConflict) || v.ID != "" {
			t.Fatal("incoherent receipt accepted")
		}
		if _, e := pool.Exec(ctx, `update platform.idempotency_record set resource_type='service-resource',response_code=201,response_body=jsonb_build_object('resource_id',$2::text) where tenant_id=$1 and scope='franchise-resource'`, tenant, id); e != nil {
			t.Fatal(e)
		}
	}
	if e := pool.QueryRow(ctx, `select payload->>'actor_subject' from platform.outbox_event where tenant_id=$1 and event_id=$2::uuid`, tenant, event).Scan(&actor); e != nil || actor != "resource-manager" {
		t.Fatal("replay overwrote creator")
	}
	var skills int
	if err := pool.QueryRow(ctx, `select count(*) from crm.resource_skill where tenant_id=$1`, tenant).Scan(&skills); err != nil || skills != 2 {
		t.Fatalf("skill transaction leaked rows=%d error=%v", skills, err)
	}
	if len(after.Skills) != 2 || after.Skills[0] != "delivery" || after.Skills[1] != "service" {
		t.Fatal("skills not recovered coherently")
	}

	t.Log("RESOURCE_IDENTITY_PASS create=1 replay=7 actor_preserved inactive_replay rollback scoped_lookup corruption fresh_pool")
}

func TestSlotCreationIdentityAndRecovery(t *testing.T) {
	raw := os.Getenv("TEST_DATABASE_URL")
	if raw == "" {
		t.Skip("disposable database not selected")
	}
	cfg, err := pgxpool.ParseConfig(raw)
	if err != nil || cfg.ConnConfig.Host != "127.0.0.1" || !strings.HasPrefix(cfg.ConnConfig.Database, "elite_confirmation_") {
		t.Fatal("requires disposable loopback database")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := (randomid.Generator{}).New()
	for _, sql := range []string{`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1::uuid,'create-'||replace(($1::uuid)::text,'-',''),'Synthetic','Synthetic')`, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'store','store','Synthetic','store')`} {
		if _, err = pool.Exec(ctx, sql, tenant); err != nil {
			t.Fatal(err)
		}
	}
	repo := NewFranchiseJourney(pool)
	start := time.Now().UTC().Add(72 * time.Hour).Truncate(time.Second)
	input := franchisejourney.AppointmentSlot{ID: "created-slot", OrganizationID: "store", Kind: "service", StartsAt: start, EndsAt: start.Add(time.Hour), Capacity: 2, State: "open", Version: 1}
	if _, err := repo.CreateAvailability(ctx, tenant, "slot-manager", franchisejourney.AvailabilityEntry{ID: "working-window", OrganizationID: "store", EntryType: "working", StartsAt: start.Add(-time.Hour), EndsAt: start.Add(72 * time.Hour)}, randomid.Generator{}.New()); err != nil {
		t.Fatal(err)
	}

	const key = "slot-concurrent-key"
	hash := strings.Repeat("a", 64)
	type outcome struct {
		v      franchisejourney.AppointmentSlot
		replay bool
		err    error
	}
	results := make(chan outcome, 8)
	begin := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-begin
			v := input
			v.ID = fmt.Sprintf("created-%d", i)
			got, replay, e := repo.CreateAppointmentSlotOnce(ctx, tenant, "slot-manager", key, hash, v, (randomid.Generator{}).New())
			results <- outcome{got, replay, e}
		}(i)
	}
	close(begin)
	wg.Wait()
	close(results)
	created, replays := 0, 0
	var id string
	for r := range results {
		if r.err != nil {
			t.Fatal(r.err)
		}
		if id != "" && id != r.v.ID {
			t.Fatal("duplicate identity")
		}
		id = r.v.ID
		if r.replay {
			replays++
		} else {
			created++
		}
	}
	if created != 1 || replays != 7 {
		t.Fatalf("created=%d replay=%d", created, replays)
	}
	recovered, e := repo.AppointmentSlotCreationResult(ctx, tenant, "store", key)
	if e != nil || recovered.ID != id {
		t.Fatalf("recover=%+v err=%v", recovered, e)
	}
	for _, scope := range [][3]string{{tenant, "other", key}, {tenant, "store", "unknown-slot-key"}, {(randomid.Generator{}).New(), "store", key}} {
		if v, e := repo.AppointmentSlotCreationResult(ctx, scope[0], scope[1], scope[2]); !errors.Is(e, franchisejourney.ErrConflict) || v.ID != "" {
			t.Fatalf("scope leak=%+v %v", v, e)
		}
	}
	if _, _, e := repo.CreateAppointmentSlotOnce(ctx, tenant, "slot-manager", key, strings.Repeat("b", 64), input, (randomid.Generator{}).New()); !errors.Is(e, franchisejourney.ErrConflict) {
		t.Fatal("divergent intent accepted")
	}
	var event, actor string
	var rows, keys, events int
	if e := pool.QueryRow(ctx, `select event_id::text,payload->>'actor_subject' from platform.outbox_event where tenant_id=$1 and aggregate_id=$2 and event_type='appointment-slot.created'`, tenant, id).Scan(&event, &actor); e != nil || actor != "slot-manager" {
		t.Fatalf("actor=%s err=%v", actor, e)
	}
	if _, _, e := repo.CreateAppointmentSlotOnce(ctx, tenant, "another-operator", key, hash, input, (randomid.Generator{}).New()); e != nil {
		t.Fatal(e)
	}
	next := input
	next.ID = "rolled-back"
	next.StartsAt = start.Add(48 * time.Hour)
	next.EndsAt = next.StartsAt.Add(time.Hour)
	if v, _, e := repo.CreateAppointmentSlotOnce(ctx, tenant, "slot-manager", "slot-rollback-key", hash, next, event); e == nil || v.ID != "" {
		t.Fatalf("audit failure not atomic=%+v %v", v, e)
	}
	if e := pool.QueryRow(ctx, `select (select count(*) from crm.appointment_slot where tenant_id=$1),(select count(*) from platform.idempotency_record where tenant_id=$1 and scope='franchise-slot'),(select count(*) from platform.outbox_event where tenant_id=$1 and event_type='appointment-slot.created')`, tenant).Scan(&rows, &keys, &events); e != nil || rows != 1 || keys != 1 || events != 1 {
		t.Fatalf("effects=%d/%d/%d err=%v", rows, keys, events, e)
	}
	for _, q := range []string{
		`insert into org.public_location(tenant_id,organization_id,city,region,country,published)values($1,'store','Synthetic','Synthetic','AR',true)`,
		`insert into crm.lead(tenant_id,lead_id,organization_id,lifecycle_state,source_code,contact_payload)values($1,'lead','store','new','fixture','{}')`,
	} {
		if _, e := pool.Exec(ctx, q, tenant); e != nil {
			t.Fatal(e)
		}
	}
	tenantCode := "create-" + strings.ReplaceAll(tenant, "-", "")
	for i := 0; i < 2; i++ {
		appointment := franchisejourney.Appointment{ID: fmt.Sprintf("booking-%d", i), LeadID: "lead", Kind: "service", StartsAt: start}
		if _, _, e := repo.RequestAppointment(ctx, tenantCode, "store", fmt.Sprintf("slot-booking-fixture-%d", i), appointment, hash, randomid.Generator{}.New()); e != nil {
			t.Fatal(e)
		}
	}
	full, e := repo.AppointmentSlotCreationResult(ctx, tenant, "store", key)
	if e != nil || full.Booked != 2 || full.Capacity != 2 {
		t.Fatalf("full slot not recovered: %+v %v", full, e)
	}
	public, e := repo.PublicAppointmentSlots(ctx, tenantCode, "store", "service", start.Add(-time.Minute), start.Add(2*time.Hour))
	if e != nil || len(public) != 0 {
		t.Fatal("full slot publicly bookable")
	}
	if _, e := pool.Exec(ctx, `update crm.appointment_slot set state='closed',version=2 where tenant_id=$1 and slot_id=$2`, tenant, id); e != nil {
		t.Fatal(e)
	}

	after, replay, e := repo.CreateAppointmentSlotOnce(ctx, tenant, "another-operator", key, hash, input, (randomid.Generator{}).New())
	if e != nil || !replay || after.ID != id || after.State != "closed" || after.Version != 2 {
		t.Fatalf("replay reopened slot=%+v replay=%v err=%v", after, replay, e)
	}
	fresh, e := pgxpool.NewWithConfig(ctx, cfg.Copy())
	if e != nil {
		t.Fatal(e)
	}
	defer fresh.Close()
	after, e = NewFranchiseJourney(fresh).AppointmentSlotCreationResult(ctx, tenant, "store", key)
	if e != nil || after.State != "closed" {
		t.Fatalf("fresh pool recovery=%+v err=%v", after, e)
	}
	for _, mutation := range []string{"resource_type='other'", "response_code=200", "response_body='{\"slot_id\":\"other\"}'"} {
		if _, e := pool.Exec(ctx, "update platform.idempotency_record set "+mutation+" where tenant_id=$1 and scope='franchise-slot'", tenant); e != nil {
			t.Fatal(e)
		}
		if v, e := repo.AppointmentSlotCreationResult(ctx, tenant, "store", key); !errors.Is(e, franchisejourney.ErrConflict) || v.ID != "" {
			t.Fatal("incoherent receipt accepted")
		}
		if _, e := pool.Exec(ctx, `update platform.idempotency_record set resource_type='appointment-slot',response_code=201,response_body=jsonb_build_object('slot_id',$2::text) where tenant_id=$1 and scope='franchise-slot'`, tenant, id); e != nil {
			t.Fatal(e)
		}
	}
	if e := pool.QueryRow(ctx, `select payload->>'actor_subject' from platform.outbox_event where tenant_id=$1 and event_id=$2::uuid`, tenant, event).Scan(&actor); e != nil || actor != "slot-manager" {
		t.Fatal("replay overwrote creator")
	}
	t.Log("SLOT_IDENTITY_PASS create=1 replay=7 actor_preserved full_closed_replay rollback scoped_lookup corruption fresh_pool")
}

func TestChecklistPublicationIdentityAndRecovery(t *testing.T) {
	raw := os.Getenv("TEST_DATABASE_URL")
	if raw == "" {
		t.Skip("disposable database not selected")
	}
	cfg, err := pgxpool.ParseConfig(raw)
	if err != nil || cfg.ConnConfig.Host != "127.0.0.1" || !strings.HasPrefix(cfg.ConnConfig.Database, "elite_confirmation_") {
		t.Fatal("requires disposable loopback database")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := (randomid.Generator{}).New()
	for _, sql := range []string{`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1::uuid,'publish-'||replace(($1::uuid)::text,'-',''),'Synthetic','Synthetic')`, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'store','store','Synthetic','store')`} {
		if _, err = pool.Exec(ctx, sql, tenant); err != nil {
			t.Fatal(err)
		}
	}
	repo := NewFranchiseJourney(pool)
	input := franchisejourney.DeliveryChecklist{ID: "published-checklist", OrganizationID: "store", Version: 1, Title: "Synthetic checklist", Items: []franchisejourney.ChecklistItem{{ID: "serial", Ordinal: 1, Prompt: "Serie", ResponseType: "serial", Required: true}, {ID: "accept", Ordinal: 2, Prompt: "Confirmación", ResponseType: "confirmation", Required: true}}}
	type outcome struct {
		actor string
		err   error
	}
	results := make(chan outcome, 8)
	begin := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-begin
			actor := fmt.Sprintf("publisher-%d", i)
			_, e := repo.PublishDeliveryChecklist(ctx, tenant, actor, input, randomid.Generator{}.New())
			results <- outcome{actor, e}
		}(i)
	}
	close(begin)
	wg.Wait()
	close(results)
	created, conflicts := 0, 0
	winner := ""
	for result := range results {
		if result.err == nil {
			created++
			winner = result.actor
		} else if errors.Is(result.err, franchisejourney.ErrConflict) {
			conflicts++
		} else {
			t.Fatal(result.err)
		}
	}
	if created != 1 || conflicts != 7 {
		t.Fatal(created, conflicts)
	}
	got, err := repo.PublishedDeliveryChecklist(ctx, tenant, "store", input.ID, 1)
	if err != nil || got.State != "published" || got.Title != input.Title || len(got.Items) != 2 || got.Items[0].ID != "serial" || got.Items[1].Ordinal != 2 {
		t.Fatal(got, err)
	}
	var rows, items, events int
	var creator, actor, eventID string
	err = pool.QueryRow(ctx, `select (select count(*) from sales.delivery_checklist_template where tenant_id=$1),(select count(*) from sales.delivery_checklist_item where tenant_id=$1),(select count(*) from platform.outbox_event where tenant_id=$1),t.created_by_subject,e.payload->>'actor_subject',e.event_id from sales.delivery_checklist_template t join platform.outbox_event e on e.tenant_id=t.tenant_id and e.aggregate_id=t.checklist_id where t.tenant_id=$1`, tenant).Scan(&rows, &items, &events, &creator, &actor, &eventID)
	if err != nil || rows != 1 || items != 2 || events != 1 || creator != winner || actor != winner {
		t.Fatal(rows, items, events, creator, actor, winner, err)
	}
	for _, scope := range [][3]string{{randomid.Generator{}.New(), "store", input.ID}, {tenant, "other", input.ID}, {tenant, "store", "unknown"}} {
		if v, e := repo.PublishedDeliveryChecklist(ctx, scope[0], scope[1], scope[2], 1); !errors.Is(e, franchisejourney.ErrConflict) || v.ID != "" {
			t.Fatal("scope leak", v, e)
		}
	}
	next := input
	next.Version = 2
	next.Title = "Second version"
	if _, e := repo.PublishDeliveryChecklist(ctx, tenant, "later-actor", next, eventID); e == nil {
		t.Fatal("duplicate outbox ID must roll back")
	}
	if v, e := repo.PublishedDeliveryChecklist(ctx, tenant, "store", input.ID, 2); !errors.Is(e, franchisejourney.ErrConflict) || v.ID != "" {
		t.Fatal("rollback lookup", v, e)
	}
	var remnants int
	if e := pool.QueryRow(ctx, `select (select count(*) from sales.delivery_checklist_template where tenant_id=$1 and checklist_version=2)+(select count(*) from sales.delivery_checklist_item where tenant_id=$1 and checklist_version=2)`, tenant).Scan(&remnants); e != nil || remnants != 0 {
		t.Fatal("partial publication", remnants, e)
	}
	if _, e := repo.PublishDeliveryChecklist(ctx, tenant, "later-actor", next, randomid.Generator{}.New()); e != nil {
		t.Fatal(e)
	}
	for _, sql := range []string{
		`update sales.delivery_checklist_template set title='mutated' where tenant_id=$1 and checklist_version=1`,
		`delete from sales.delivery_checklist_template where tenant_id=$1 and checklist_version=1`,
		`update sales.delivery_checklist_item set prompt='mutated' where tenant_id=$1 and checklist_version=1`,
		`delete from sales.delivery_checklist_item where tenant_id=$1 and checklist_version=1`,
		`insert into sales.delivery_checklist_item(tenant_id,organization_id,checklist_id,checklist_version,item_id,ordinal,prompt,response_type,required) values($1,'store','published-checklist',1,'extra',3,'Extra','text',false)`,
	} {
		if _, e := pool.Exec(ctx, sql, tenant); e == nil {
			t.Fatal("published version accepted mutation")
		}
	}
	if _, e := pool.Exec(ctx, `insert into sales.delivery_checklist_template(tenant_id,organization_id,checklist_id,checklist_version,title,state,created_by_subject) values($1,'store','draft-checklist',1,'Draft','draft','draft-author')`, tenant); e != nil {
		t.Fatal(e)
	}
	if v, e := repo.PublishedDeliveryChecklist(ctx, tenant, "store", "draft-checklist", 1); !errors.Is(e, franchisejourney.ErrConflict) || v.ID != "" {
		t.Fatal("draft exposed", v, e)
	}
	fresh, e := pgxpool.NewWithConfig(ctx, cfg.Copy())
	if e != nil {
		t.Fatal(e)
	}
	defer fresh.Close()
	for version, title := range map[int64]string{1: input.Title, 2: next.Title} {
		v, e := NewFranchiseJourney(fresh).PublishedDeliveryChecklist(ctx, tenant, "store", input.ID, version)
		if e != nil || v.Version != version || v.Title != title || len(v.Items) != 2 || v.Items[0].Prompt != "Serie" {
			t.Fatal("fresh lookup", v, e)
		}
	}
	t.Log("CHECKLIST_PUBLICATION_POSTGRES_PASS created=1 conflicts=7 versions=2 rollback=1 immutable=5 scope=3 draft=blocked fresh_pool=pass")
}

func TestChecklistCompletionIdentityAndRecovery(t *testing.T) {
	raw := os.Getenv("TEST_DATABASE_URL")
	if raw == "" {
		t.Skip("disposable database required")
	}
	cfg, err := pgxpool.ParseConfig(raw)
	if err != nil || cfg.ConnConfig.Host != "127.0.0.1" || !strings.HasPrefix(cfg.ConnConfig.Database, "elite_confirmation_") {
		t.Fatal("requires disposable loopback database")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := randomid.Generator{}.New()
	for _, q := range []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1::uuid,'complete-'||replace(($1::uuid)::text,'-',''),'Synthetic','Synthetic')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'store','store','Synthetic','store')`,
		`insert into catalog.vehicle_model(tenant_id,model_id,model_code,display_name,vehicle_class,lifecycle_state)values($1,'model','model','Synthetic','bicycle','active')`,
		`insert into catalog.vehicle_variant(tenant_id,variant_id,model_id,variant_code,display_name,battery_specification,lifecycle_state)values($1,'variant','model','variant','Synthetic','{}','active')`,
		`insert into crm.customer_profile(tenant_id,customer_principal_id,display_name,email_normalized)values($1,'customer','Synthetic','synthetic@example.test')`,
		`insert into inventory.stock_unit(tenant_id,stock_unit_id,organization_id,variant_id,serial_number,state,version,received_at)values($1,'stock','store','variant','SYNTHETIC-SERIAL','sold',1,clock_timestamp())`,
		`insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version)values($1,'order','store','customer','delivered','ARS',123456,1)`,
		`insert into sales.delivery_handover(tenant_id,handover_id,organization_id,order_id,customer_principal_id,stock_unit_id,state,version)values($1,'completion-fixture','store','order','customer','stock','prepared',1)`,
	} {
		if _, err = pool.Exec(ctx, q, tenant); err != nil {
			t.Fatal(err)
		}
	}
	repo := NewFranchiseJourney(pool)
	checklist := franchisejourney.DeliveryChecklist{ID: "completion-checklist", OrganizationID: "store", Version: 1, Title: "Synthetic", Items: []franchisejourney.ChecklistItem{{ID: "serial", Ordinal: 1, Prompt: "Serie", ResponseType: "serial", Required: true}, {ID: "confirmed", Ordinal: 2, Prompt: "Confirmación", ResponseType: "confirmation", Required: true}}}
	if _, err = repo.PublishDeliveryChecklist(ctx, tenant, "fixture-author", checklist, randomid.Generator{}.New()); err != nil {
		t.Fatal(err)
	}
	responses := []franchisejourney.ChecklistResponse{{ItemID: "serial", ResponseText: "SYNTHETIC-SERIAL"}, {ItemID: "confirmed", ResponseText: "confirmed"}}
	type outcome struct {
		actor string
		err   error
	}
	results := make(chan outcome, 8)
	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-start
			actor := fmt.Sprintf("completer-%d", i)
			_, e := repo.CompleteDeliveryChecklist(ctx, tenant, "store", actor, "completion-fixture", 1, checklist.ID, 1, responses, randomid.Generator{}.New())
			results <- outcome{actor, e}
		}(i)
	}
	close(start)
	wg.Wait()
	close(results)
	created, conflicts := 0, 0
	winner := ""
	for result := range results {
		if result.err == nil {
			created++
			winner = result.actor
		} else if errors.Is(result.err, franchisejourney.ErrConflict) {
			conflicts++
		} else {
			t.Fatal(result.err)
		}
	}
	if created != 1 || conflicts != 7 {
		t.Fatal(created, conflicts)
	}
	got, err := repo.DeliveryChecklistCompletion(ctx, tenant, "store", "completion-fixture")
	if err != nil || got.Version != 2 || got.State != "presented" || got.ActorSubject != winner || len(got.Responses) != 2 || got.Responses[0].ItemID != "confirmed" || got.Responses[1].ResponseText != "SYNTHETIC-SERIAL" {
		t.Fatal(got, err)
	}
	var events, answers int
	var eventID, actor string
	err = pool.QueryRow(ctx, `select (select count(*) from platform.outbox_event where tenant_id=$1 and event_type='delivery-handover.checklist-completed'),(select count(*) from sales.delivery_checklist_response where tenant_id=$1),event_id,payload->>'actor_subject' from platform.outbox_event where tenant_id=$1 and event_type='delivery-handover.checklist-completed'`, tenant).Scan(&events, &answers, &eventID, &actor)
	if err != nil || events != 1 || answers != 2 || actor != winner {
		t.Fatal(events, answers, actor, err)
	}
	for _, scope := range [][3]string{{randomid.Generator{}.New(), "store", "completion-fixture"}, {tenant, "other", "completion-fixture"}, {tenant, "store", "unknown"}} {
		if v, e := repo.DeliveryChecklistCompletion(ctx, scope[0], scope[1], scope[2]); !errors.Is(e, franchisejourney.ErrConflict) || v.HandoverID != "" {
			t.Fatal("scope", v, e)
		}
	}
	if _, e := pool.Exec(ctx, `insert into sales.delivery_handover(tenant_id,handover_id,organization_id,order_id,customer_principal_id,stock_unit_id,state,version) values($1,'rollback-fixture','store','order','customer','stock','prepared',1)`, tenant); e != nil {
		t.Fatal(e)
	}
	if _, e := repo.CompleteDeliveryChecklist(ctx, tenant, "store", "later-actor", "rollback-fixture", 1, checklist.ID, 1, responses, eventID); e == nil {
		t.Fatal("duplicate event must roll back")
	}
	var state string
	var version int64
	if e := pool.QueryRow(ctx, `select state,version,(select count(*) from sales.delivery_checklist_response where tenant_id=$1 and handover_id='rollback-fixture') from sales.delivery_handover where tenant_id=$1 and handover_id='rollback-fixture'`, tenant).Scan(&state, &version, &answers); e != nil || state != "prepared" || version != 1 || answers != 0 {
		t.Fatal("partial completion", state, version, answers, e)
	}
	if v, e := repo.DeliveryChecklistCompletion(ctx, tenant, "store", "rollback-fixture"); !errors.Is(e, franchisejourney.ErrConflict) || v.HandoverID != "" {
		t.Fatal("prepared exposed", v, e)
	}
	if _, e := repo.AcceptHandover(ctx, tenant, "store", "customer", "completion-fixture", 2, "SYNTHETIC-SERIAL", checklist.ID, 1, strings.Repeat("a", 64), randomid.Generator{}.New()); e != nil {
		t.Fatal(e)
	}
	if _, e := repo.CompleteDeliveryChecklist(ctx, tenant, "store", "later-actor", "rollback-fixture", 1, checklist.ID, 1, responses, randomid.Generator{}.New()); e != nil {
		t.Fatal(e)
	}
	if _, e := repo.RejectHandover(ctx, tenant, "store", "customer", "rollback-fixture", 2, "serial-mismatch", "Synthetic rejection", strings.Repeat("b", 64), randomid.Generator{}.New(), randomid.Generator{}.New()); e != nil {
		t.Fatal(e)
	}
	for _, sql := range []string{`update sales.delivery_checklist_response set response_text='changed' where tenant_id=$1 and handover_id='completion-fixture'`, `delete from sales.delivery_checklist_response where tenant_id=$1 and handover_id='completion-fixture'`, `update sales.delivery_handover set checklist_completed_by_subject='forged' where tenant_id=$1 and handover_id='completion-fixture'`} {
		if _, e := pool.Exec(ctx, sql, tenant); e == nil {
			t.Fatal("completed binding accepted mutation")
		}
	}
	for _, sql := range []string{`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'other','other','Other','store')`, `insert into inventory.stock_unit(tenant_id,stock_unit_id,organization_id,variant_id,serial_number,vin,battery_serial_number,state,version,received_at)values($1,'foreign-stock','other','variant','FOREIGN-SYNTHETIC','FOREIGN-SYNTHETIC','FOREIGN-SYNTHETIC','sold',1,clock_timestamp())`, `update sales.delivery_handover set stock_unit_id='foreign-stock' where tenant_id=$1 and handover_id='completion-fixture'`} {
		if _, e := pool.Exec(ctx, sql, tenant); e != nil {
			t.Fatal(e)
		}
	}
	if v, e := repo.DeliveryChecklistCompletion(ctx, tenant, "store", "completion-fixture"); !errors.Is(e, franchisejourney.ErrConflict) || v.HandoverID != "" {
		t.Fatal("foreign stock leaked", v, e)
	}
	if _, e := pool.Exec(ctx, `update sales.delivery_handover set stock_unit_id='stock' where tenant_id=$1 and handover_id='completion-fixture'`, tenant); e != nil {
		t.Fatal(e)
	}
	fresh, e := pgxpool.NewWithConfig(ctx, cfg.Copy())
	if e != nil {
		t.Fatal(e)
	}
	defer fresh.Close()
	for id, want := range map[string]string{"completion-fixture": "accepted", "rollback-fixture": "rejected"} {
		v, e := NewFranchiseJourney(fresh).DeliveryChecklistCompletion(ctx, tenant, "store", id)
		if e != nil || v.State != want || v.Version != 3 || len(v.Responses) != 2 || v.Responses[0].ResponseText != "confirmed" || v.Responses[1].ResponseText != "SYNTHETIC-SERIAL" {
			t.Fatal("current state recovery", v, e)
		}
		if id == "completion-fixture" && (v.ActorSubject != winner || !v.CompletedAt.Equal(got.CompletedAt)) {
			t.Fatal("completion audit changed", v)
		}
	}
	t.Log("CHECKLIST_COMPLETION_POSTGRES_PASS completed=1 conflicts=7 rollback=1 immutable=3 scope=3 foreign_stock=blocked accepted=preserved rejected=preserved actor=preserved fresh_pool=pass")
}

func TestReturnScopeBoundaries(t *testing.T) {
	raw := os.Getenv("TEST_DATABASE_URL")
	if raw == "" {
		t.Skip("disposable database required")
	}
	cfg, err := pgxpool.ParseConfig(raw)
	if err != nil || cfg.ConnConfig.Host != "127.0.0.1" || !strings.HasPrefix(cfg.ConnConfig.Database, "elite_confirmation_") {
		t.Fatal("requires disposable loopback database")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := randomid.Generator{}.New()
	for _, q := range []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1::uuid,'complete-'||replace(($1::uuid)::text,'-',''),'Synthetic','Synthetic')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'store','store','Synthetic','store')`,
		`insert into catalog.vehicle_model(tenant_id,model_id,model_code,display_name,vehicle_class,lifecycle_state)values($1,'model','model','Synthetic','bicycle','active')`,
		`insert into catalog.vehicle_variant(tenant_id,variant_id,model_id,variant_code,display_name,battery_specification,lifecycle_state)values($1,'variant','model','variant','Synthetic','{}','active')`,
		`insert into crm.customer_profile(tenant_id,customer_principal_id,display_name,email_normalized)values($1,'customer','Synthetic','synthetic@example.test')`,
		`insert into inventory.stock_unit(tenant_id,stock_unit_id,organization_id,variant_id,serial_number,state,version,received_at)values($1,'stock','store','variant','SYNTHETIC-SERIAL','sold',1,clock_timestamp())`,
		`insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version)values($1,'order','store','customer','delivered','ARS',123456,1)`,
		`insert into sales.delivery_handover(tenant_id,handover_id,organization_id,order_id,customer_principal_id,stock_unit_id,state,version)values($1,'completion-fixture','store','order','customer','stock','prepared',1)`,
	} {
		if _, err = pool.Exec(ctx, q, tenant); err != nil {
			t.Fatal(err)
		}
	}
	repo := NewFranchiseJourney(pool)
	checklist := franchisejourney.DeliveryChecklist{ID: "completion-checklist", OrganizationID: "store", Version: 1, Title: "Synthetic", Items: []franchisejourney.ChecklistItem{{ID: "serial", Ordinal: 1, Prompt: "Serie", ResponseType: "serial", Required: true}, {ID: "confirmed", Ordinal: 2, Prompt: "Confirmación", ResponseType: "confirmation", Required: true}}}
	if _, err = repo.PublishDeliveryChecklist(ctx, tenant, "fixture-author", checklist, randomid.Generator{}.New()); err != nil {
		t.Fatal(err)
	}

	responses := []franchisejourney.ChecklistResponse{{ItemID: "serial", ResponseText: "SYNTHETIC-SERIAL"}, {ItemID: "confirmed", ResponseText: "confirmed"}}
	if _, err = repo.CompleteDeliveryChecklist(ctx, tenant, "store", "fixture-operator", "completion-fixture", 1, checklist.ID, 1, responses, randomid.Generator{}.New()); err != nil {
		t.Fatal(err)
	}
	exception, err := repo.RejectHandover(ctx, tenant, "store", "customer", "completion-fixture", 2, "serial-mismatch", "Synthetic rejection", strings.Repeat("a", 64), randomid.Generator{}.New(), randomid.Generator{}.New())
	if err != nil {
		t.Fatal(err)
	}
	resolution, err := repo.ResolveDeliveryException(ctx, tenant, "store", "fixture-manager", exception.ID, 1, "return", "Synthetic authorization", "", "return-authorization", randomid.Generator{}.New(), randomid.Generator{}.New())
	if err != nil || resolution.ReturnAuthorizationID != "return-authorization" {
		t.Fatal(resolution, err)
	}

	if _, err = pool.Exec(ctx, `insert into sales.delivery_handover(tenant_id,handover_id,organization_id,order_id,customer_principal_id,stock_unit_id,state,version) select tenant_id,'scope-second',organization_id,order_id,customer_principal_id,stock_unit_id,'prepared',1 from sales.delivery_handover where tenant_id=$1 and handover_id='completion-fixture'`, tenant); err != nil {
		t.Fatal(err)
	}
	if _, err = repo.CompleteDeliveryChecklist(ctx, tenant, "store", "fixture-operator", "scope-second", 1, checklist.ID, 1, responses, randomid.Generator{}.New()); err != nil {
		t.Fatal(err)
	}
	second, err := repo.RejectHandover(ctx, tenant, "store", "customer", "scope-second", 2, "serial-mismatch", "Synthetic rejection", strings.Repeat("a", 64), randomid.Generator{}.New(), randomid.Generator{}.New())
	if err != nil {
		t.Fatal(err)
	}
	if _, err = repo.ResolveDeliveryException(ctx, tenant, "store", "fixture-manager", second.ID, 1, "return", "Synthetic authorization", "", "scope-second-authorization", randomid.Generator{}.New(), randomid.Generator{}.New()); err != nil {
		t.Fatal(err)
	}
	received, err := repo.ReceiveReturn(ctx, tenant, "store", "valid-receiver", "scope-second-authorization", "SYNTHETIC-SERIAL", "sealed", "Synthetic receipt", strings.Repeat("b", 64), randomid.Generator{}.New(), randomid.Generator{}.New())
	if err != nil {
		t.Fatal(err)
	}
	if _, err = pool.Exec(ctx, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'other','other','Other','store')`, tenant); err != nil {
		t.Fatal(err)
	}
	if _, err = pool.Exec(ctx, `update inventory.stock_unit set organization_id='other' where tenant_id=$1 and stock_unit_id='stock'`, tenant); err != nil {
		t.Fatal(err)
	}
	if values, e := repo.ReturnCases(ctx, tenant, "store", 100); !errors.Is(e, franchisejourney.ErrConflict) || len(values) != 0 {
		t.Errorf("foreign graph listed cases=%d error=%v", len(values), e)
	}
	if value, e := repo.ReceiveReturn(ctx, tenant, "store", "scope-receiver", "return-authorization", "SYNTHETIC-SERIAL", "sealed", "Synthetic receipt", strings.Repeat("b", 64), randomid.Generator{}.New(), randomid.Generator{}.New()); !errors.Is(e, franchisejourney.ErrConflict) || value.ID != "" {
		t.Errorf("foreign graph received id=%s error=%v", value.ID, e)
	}
	if value, e := repo.DecideReturn(ctx, tenant, "store", "scope-decider", received.ID, "quarantine", "Synthetic decision", randomid.Generator{}.New(), randomid.Generator{}.New(), randomid.Generator{}.New(), randomid.Generator{}.New(), randomid.Generator{}.New(), randomid.Generator{}.New()); !errors.Is(e, franchisejourney.ErrConflict) || value.ID != "" {
		t.Errorf("foreign graph decided id=%s effects=%d error=%v", value.ID, len(value.Effects), e)
	}
	var receipts, decisions, effects, events int
	err = pool.QueryRow(ctx, `select (select count(*) from sales.return_receipt where tenant_id=$1 and authorization_id='return-authorization'),(select count(*) from sales.return_disposition where tenant_id=$1),(select count(*) from sales.return_effect_request where tenant_id=$1),(select count(*) from platform.outbox_event where tenant_id=$1 and payload->>'actor_subject' in ('scope-receiver','scope-decider'))`, tenant).Scan(&receipts, &decisions, &effects, &events)
	if err != nil || receipts != 0 || decisions != 0 || effects != 0 || events != 0 {
		t.Errorf("foreign graph effects=%d/%d/%d/%d error=%v", receipts, decisions, effects, events, err)
	}
	if !t.Failed() {
		t.Log("RETURN_SCOPE_PASS foreign_stock_blocks_list_receive_decide effects=0")
	}
}

func TestReturnOperationsIdentityAndRecovery(t *testing.T) {
	raw := os.Getenv("TEST_DATABASE_URL")
	if raw == "" {
		t.Skip("disposable database required")
	}
	cfg, err := pgxpool.ParseConfig(raw)
	if err != nil || cfg.ConnConfig.Host != "127.0.0.1" || !strings.HasPrefix(cfg.ConnConfig.Database, "elite_confirmation_") {
		t.Fatal("requires disposable loopback database")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := randomid.Generator{}.New()
	for _, q := range []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1::uuid,'complete-'||replace(($1::uuid)::text,'-',''),'Synthetic','Synthetic')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'store','store','Synthetic','store')`,
		`insert into catalog.vehicle_model(tenant_id,model_id,model_code,display_name,vehicle_class,lifecycle_state)values($1,'model','model','Synthetic','bicycle','active')`,
		`insert into catalog.vehicle_variant(tenant_id,variant_id,model_id,variant_code,display_name,battery_specification,lifecycle_state)values($1,'variant','model','variant','Synthetic','{}','active')`,
		`insert into crm.customer_profile(tenant_id,customer_principal_id,display_name,email_normalized)values($1,'customer','Synthetic','synthetic@example.test')`,
		`insert into inventory.stock_unit(tenant_id,stock_unit_id,organization_id,variant_id,serial_number,state,version,received_at)values($1,'stock','store','variant','SYNTHETIC-SERIAL','sold',1,clock_timestamp())`,
		`insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version)values($1,'order','store','customer','delivered','ARS',123456,1)`,
		`insert into sales.delivery_handover(tenant_id,handover_id,organization_id,order_id,customer_principal_id,stock_unit_id,state,version)values($1,'completion-fixture','store','order','customer','stock','prepared',1)`,
	} {
		if _, err = pool.Exec(ctx, q, tenant); err != nil {
			t.Fatal(err)
		}
	}
	repo := NewFranchiseJourney(pool)
	checklist := franchisejourney.DeliveryChecklist{ID: "completion-checklist", OrganizationID: "store", Version: 1, Title: "Synthetic", Items: []franchisejourney.ChecklistItem{{ID: "serial", Ordinal: 1, Prompt: "Serie", ResponseType: "serial", Required: true}, {ID: "confirmed", Ordinal: 2, Prompt: "Confirmación", ResponseType: "confirmation", Required: true}}}
	if _, err = repo.PublishDeliveryChecklist(ctx, tenant, "fixture-author", checklist, randomid.Generator{}.New()); err != nil {
		t.Fatal(err)
	}

	responses := []franchisejourney.ChecklistResponse{{ItemID: "serial", ResponseText: "SYNTHETIC-SERIAL"}, {ItemID: "confirmed", ResponseText: "confirmed"}}
	if _, err = repo.CompleteDeliveryChecklist(ctx, tenant, "store", "fixture-operator", "completion-fixture", 1, checklist.ID, 1, responses, randomid.Generator{}.New()); err != nil {
		t.Fatal(err)
	}
	exception, err := repo.RejectHandover(ctx, tenant, "store", "customer", "completion-fixture", 2, "serial-mismatch", "Synthetic rejection", strings.Repeat("a", 64), randomid.Generator{}.New(), randomid.Generator{}.New())
	if err != nil {
		t.Fatal(err)
	}
	resolution, err := repo.ResolveDeliveryException(ctx, tenant, "store", "fixture-manager", exception.ID, 1, "return", "Synthetic authorization", "", "return-authorization", randomid.Generator{}.New(), randomid.Generator{}.New())
	if err != nil || resolution.ReturnAuthorizationID != "return-authorization" {
		t.Fatal(resolution, err)
	}

	newID := func() string { return randomid.Generator{}.New() }
	authorize := func(id, action string) {
		t.Helper()
		if _, e := pool.Exec(ctx, `insert into sales.delivery_handover(tenant_id,handover_id,organization_id,order_id,customer_principal_id,stock_unit_id,state,version)values($1,$2,'store','order','customer','stock','prepared',1)`, tenant, id); e != nil {
			t.Fatal(e)
		}
		if _, e := repo.CompleteDeliveryChecklist(ctx, tenant, "store", "fixture", id, 1, checklist.ID, 1, responses, newID()); e != nil {
			t.Fatal(e)
		}
		x, e := repo.RejectHandover(ctx, tenant, "store", "customer", id, 2, "serial-mismatch", "Synthetic", strings.Repeat("a", 64), newID(), newID())
		if e != nil {
			t.Fatal(e)
		}
		if _, e = repo.ResolveDeliveryException(ctx, tenant, "store", "fixture", x.ID, 1, action, "Synthetic", "", id, newID(), newID()); e != nil {
			t.Fatal(e)
		}
	}
	receive := func(auth, actor, event string) (franchisejourney.ReturnReceipt, error) {
		return repo.ReceiveReturn(ctx, tenant, "store", actor, auth, "SYNTHETIC-SERIAL", "sealed", "Synthetic receipt", strings.Repeat("b", 64), newID(), event)
	}
	decide := func(receipt, actor, event string) (franchisejourney.ReturnDisposition, error) {
		return repo.DecideReturn(ctx, tenant, "store", actor, receipt, "quarantine", "Synthetic decision", newID(), newID(), newID(), newID(), newID(), event)
	}
	lookup := func(auth string) franchisejourney.ReturnCase {
		t.Helper()
		v, e := repo.ReturnCaseResult(ctx, tenant, "store", auth)
		if e != nil {
			t.Fatal(e)
		}
		return v
	}
	if v := lookup("return-authorization"); v.Receipt != nil || v.Disposition != nil {
		t.Fatal("uncreated operations exposed")
	}
	if v, e := repo.ReceiveReturn(ctx, tenant, "store", "wrong-serial", "return-authorization", "WRONG", "sealed", "Synthetic", strings.Repeat("b", 64), newID(), newID()); !errors.Is(e, franchisejourney.ErrConflict) || v.ID != "" {
		t.Fatal(v, e)
	}
	type outcome struct {
		id  string
		err error
	}
	race := func(fn func(int) (string, error)) string {
		t.Helper()
		ch := make(chan outcome, 8)
		start := make(chan struct{})
		var wg sync.WaitGroup
		for i := 0; i < 8; i++ {
			wg.Add(1)
			go func(i int) { defer wg.Done(); <-start; id, e := fn(i); ch <- outcome{id, e} }(i)
		}
		close(start)
		wg.Wait()
		close(ch)
		wins, conflicts, id := 0, 0, ""
		for v := range ch {
			if v.err == nil {
				wins++
				id = v.id
			} else if errors.Is(v.err, franchisejourney.ErrConflict) {
				conflicts++
			} else {
				t.Fatal(v.err)
			}
		}
		if wins != 1 || conflicts != 7 || id == "" {
			t.Fatal(wins, conflicts, id)
		}
		return id
	}
	receiptID := race(func(i int) (string, error) {
		v, e := receive("return-authorization", fmt.Sprintf("receiver-%d", i), newID())
		return v.ID, e
	})
	decisionID := race(func(i int) (string, error) {
		v, e := decide(receiptID, fmt.Sprintf("decider-%d", i), newID())
		return v.ID, e
	})
	first := lookup("return-authorization")
	if first.Receipt.ID != receiptID || first.Disposition.ID != decisionID || len(first.Disposition.Effects) != 4 {
		t.Fatal(first)
	}
	for _, actor := range []string{first.Receipt.ReceivedBySubject, first.Disposition.DecidedBySubject} {
		var n int
		if e := pool.QueryRow(ctx, `select count(*) from platform.outbox_event where tenant_id=$1 and payload->>'actor_subject'=$2`, tenant, actor).Scan(&n); e != nil || n != 1 {
			t.Fatal(actor, n, e)
		}
	}
	authorize("exchange-authorization", "exchange")
	er, e := receive("exchange-authorization", "exchange-receiver", newID())
	if e != nil {
		t.Fatal(e)
	}
	ed, e := decide(er.ID, "exchange-decider", newID())
	if e != nil || len(ed.Effects) != 3 || ed.CustomerRemedy != "exchange" {
		t.Fatal(ed, e)
	}
	for _, effect := range ed.Effects {
		if effect.OwnerContext == "fiscal" {
			t.Fatal("exchange created fiscal request")
		}
	}
	authorize("atomic-authorization", "return")
	var occupiedEvent string
	if e = pool.QueryRow(ctx, `select event_id from platform.outbox_event where tenant_id=$1 limit 1`, tenant).Scan(&occupiedEvent); e != nil {
		t.Fatal(e)
	}
	if v, e := receive("atomic-authorization", "atomic-receiver", occupiedEvent); e == nil || v.ID != "" {
		t.Fatal("outbox collision committed receipt", v, e)
	}
	if v := lookup("atomic-authorization"); v.Receipt != nil {
		t.Fatal("receipt survived rollback")
	}
	ar, e := receive("atomic-authorization", "atomic-receiver", newID())
	if e != nil {
		t.Fatal(e)
	}
	if v, e := decide(ar.ID, "atomic-decider", occupiedEvent); e == nil || v.ID != "" {
		t.Fatal("outbox collision committed decision", v, e)
	}
	if v := lookup("atomic-authorization"); v.Disposition != nil {
		t.Fatal("decision survived rollback")
	}
	var requests int
	if e = pool.QueryRow(ctx, `select count(*) from sales.return_effect_request where tenant_id=$1`, tenant).Scan(&requests); e != nil || requests != 7 {
		t.Fatal("requests survived rollback", requests, e)
	}
	if _, e = decide(ar.ID, "atomic-decider", newID()); e != nil {
		t.Fatal(e)
	}
	for _, q := range []string{`update sales.return_receipt set notes='changed' where tenant_id=$1`, `delete from sales.return_receipt where tenant_id=$1`, `update sales.return_disposition set notes='changed' where tenant_id=$1`, `delete from sales.return_disposition where tenant_id=$1`, `update sales.return_effect_request set state='changed' where tenant_id=$1`, `delete from sales.return_effect_request where tenant_id=$1`} {
		if _, e = pool.Exec(ctx, q, tenant); e == nil {
			t.Fatal("immutable evidence changed", q)
		}
	}
	// More than the historical list limit: absence from a page is not absence of a committed operation.
	for i := 0; i < 101; i++ {
		authorize(fmt.Sprintf("newer-%03d", i), "return")
	}
	list, e := repo.ReturnCases(ctx, tenant, "store", 100)
	if e != nil || len(list) != 100 {
		t.Fatal(len(list), e)
	}
	for _, v := range list {
		if v.AuthorizationID == "return-authorization" {
			t.Fatal("old case unexpectedly in newest page")
		}
	}
	fresh, e := pgxpool.NewWithConfig(ctx, cfg.Copy())
	if e != nil {
		t.Fatal(e)
	}
	defer fresh.Close()
	again, e := NewFranchiseJourney(fresh).ReturnCaseResult(ctx, tenant, "store", "return-authorization")
	if e != nil || again.Receipt.ID != receiptID || again.Disposition.ID != decisionID || again.Receipt.EvidenceSHA256 != first.Receipt.EvidenceSHA256 || again.Receipt.ReceivedBySubject != first.Receipt.ReceivedBySubject {
		t.Fatal("fresh pool lost evidence", again, e)
	}
	for _, scope := range [][3]string{{tenant, "other", "return-authorization"}, {newID(), "store", "return-authorization"}, {tenant, "store", "unknown"}} {
		v, e := repo.ReturnCaseResult(ctx, scope[0], scope[1], scope[2])
		if !errors.Is(e, franchisejourney.ErrConflict) || v.AuthorizationID != "" {
			t.Fatal(v, e)
		}
	}
	// The scope lock blocks concurrent movement of the authorization's handover, order and stock.
	lock, e := pool.Begin(ctx)
	if e != nil {
		t.Fatal(e)
	}
	defer lock.Rollback(ctx)
	if e = lockReturnAuthorizationScope(ctx, lock, tenant, "store", "return-authorization"); e != nil {
		t.Fatal(e)
	}
	for _, q := range []string{`update sales.delivery_handover set version=version+1 where tenant_id=$1 and handover_id='completion-fixture'`, `update sales.customer_order set version=version+1 where tenant_id=$1 and order_id='order'`, `update inventory.stock_unit set version=version+1 where tenant_id=$1 and stock_unit_id='stock'`} {
		writer, e := fresh.Begin(ctx)
		if e != nil {
			t.Fatal(e)
		}
		if _, e = writer.Exec(ctx, `set local lock_timeout='100ms'`); e != nil {
			t.Fatal(e)
		}
		_, e = writer.Exec(ctx, q, tenant)
		var pe *pgconn.PgError
		if !errors.As(e, &pe) || pe.Code != "55P03" {
			t.Fatal("scope lock failed", e)
		}
		writer.Rollback(ctx)
	}
	lock.Rollback(ctx)
	t.Log("RETURN_IDENTITY_RECOVERY_PASS receipt_race=1/8 decision_race=1/8 requests=4+3+4 atomic_rollbacks=2 immutable=6 outside_limit=100 fresh_pool=PASS scope_locks=3")
}
````

### FILE: `internal/platform/httpapi/franchisejourney.go`

```yaml
block_id: "GO-FRANCHISE-JOURNEY:file:08"
operation: CREATE
provenance: AUTHORED
source: "local implementation governed by the upstream references in metadata"
license: "LicenseRef-Workspace-Owner"
sha256: "2439a78ae334f0a62376963c8a627e52a4627c39e7926b3927fda361230964fe"
variables: []
secrets_allowed: false
```

````go
package httpapi

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"elite.local/enterprise/internal/franchisejourney"
	"elite.local/enterprise/internal/platform/identity"
)

type FranchiseJourneyModule struct{ Service *franchisejourney.Service }

func (m FranchiseJourneyModule) Register(mux *http.ServeMux, verifier identity.Verifier) {
	a := franchiseJourneyAPI{service: m.Service, verifier: verifier}
	mux.HandleFunc("GET /v1/public/{tenantCode}/locations", a.publicLocations)
	mux.HandleFunc("GET /v1/public/{tenantCode}/{organizationCode}/appointment-slots", a.publicAppointmentSlots)
	mux.HandleFunc("POST /v1/public/{tenantCode}/{organizationCode}/appointments", a.requestAppointment)
	mux.HandleFunc("POST /v1/franchise/appointment-slots", a.createAppointmentSlot)
	mux.HandleFunc("GET /v1/franchise/appointment-slots/result", a.appointmentSlotCreationResult)
	mux.HandleFunc("POST /v1/franchise/resources", a.createServiceResource)
	mux.HandleFunc("GET /v1/franchise/resources/result", a.serviceResourceCreationResult)
	mux.HandleFunc("GET /v1/franchise/availability", a.availability)
	mux.HandleFunc("POST /v1/franchise/availability", a.createAvailability)
	mux.HandleFunc("GET /v1/franchise/availability/result", a.availabilityCreationResult)
	mux.HandleFunc("POST /v1/franchise/availability/{id}/cancel", a.cancelAvailability)
	mux.HandleFunc("POST /v1/franchise/appointments/{id}/resources", a.assignAppointmentResource)
	mux.HandleFunc("POST /v1/franchise/appointments/{id}/transitions", a.transitionAppointment)
	mux.HandleFunc("GET /v1/franchise/leads", a.leads)
	mux.HandleFunc("GET /v1/franchise/agenda", a.appointmentAgenda)
	mux.HandleFunc("POST /v1/franchise/leads/{id}/assign", a.assignLead)
	mux.HandleFunc("POST /v1/franchise/leads/{id}/transitions", a.transitionLead)
	mux.HandleFunc("POST /v1/franchise/quotes", a.createQuote)
	mux.HandleFunc("GET /v1/franchise/quotes/result", a.quoteResult)
	mux.HandleFunc("POST /v1/franchise/delivery-checklists", a.publishDeliveryChecklist)
	mux.HandleFunc("GET /v1/franchise/delivery-checklists/result", a.publishedDeliveryChecklist)
	mux.HandleFunc("POST /v1/franchise/handovers/{id}/complete-checklist", a.completeDeliveryChecklist)
	mux.HandleFunc("GET /v1/franchise/handovers/{id}/checklist-result", a.deliveryChecklistCompletion)
	mux.HandleFunc("GET /v1/franchise/delivery-exceptions", a.deliveryExceptions)
	mux.HandleFunc("POST /v1/franchise/delivery-exceptions/{id}/resolve", a.resolveDeliveryException)
	mux.HandleFunc("GET /v1/franchise/returns", a.returnCases)
	mux.HandleFunc("GET /v1/franchise/returns/result", a.returnCaseResult)
	mux.HandleFunc("POST /v1/franchise/return-authorizations/{id}/receive", a.receiveReturn)
	mux.HandleFunc("POST /v1/franchise/return-receipts/{id}/decide", a.decideReturn)
	mux.HandleFunc("GET /v1/customer/journey", a.customerJourney)
	mux.HandleFunc("POST /v1/customer/appointments/{id}/cancel", a.cancelCustomerAppointment)
	mux.HandleFunc("POST /v1/customer/quotes/{id}/accept", a.acceptQuote)
	mux.HandleFunc("POST /v1/customer/handovers/{id}/accept", a.acceptHandover)
	mux.HandleFunc("POST /v1/customer/handovers/{id}/reject", a.rejectHandover)
}

func (a franchiseJourneyAPI) publishDeliveryChecklist(w http.ResponseWriter, r *http.Request) {
	var input struct {
		OrganizationID string                           `json:"organization_id"`
		ChecklistID    string                           `json:"checklist_id"`
		Version        int64                            `json:"version"`
		Title          string                           `json:"title"`
		Items          []franchisejourney.ChecklistItem `json:"items"`
	}
	if !decodeStrict(w, r, &input) {
		return
	}
	p, ok := a.protected(w, r, "handover:manage", input.OrganizationID)
	if !ok {
		return
	}
	value, err := a.service.PublishDeliveryChecklist(r.Context(), p.TenantID, p.Subject, franchisejourney.DeliveryChecklist{ID: input.ChecklistID, OrganizationID: input.OrganizationID, Version: input.Version, Title: input.Title, Items: input.Items})
	if errors.Is(err, franchisejourney.ErrInvalid) {
		writeProblem(w, 400, "INVALID_DELIVERY_CHECKLIST", "delivery checklist does not match contract")
		return
	}
	if errors.Is(err, franchisejourney.ErrConflict) {
		writeProblem(w, 409, "DELIVERY_CHECKLIST_CONFLICT", "delivery checklist version already exists or violates scope")
		return
	}
	if err != nil {
		writeProblem(w, 500, "DELIVERY_CHECKLIST_FAILED", "delivery checklist publication failed")
		return
	}
	writeJSON(w, 201, value)
}

func (a franchiseJourneyAPI) completeDeliveryChecklist(w http.ResponseWriter, r *http.Request) {
	var input struct {
		OrganizationID   string                               `json:"organization_id"`
		Version          int64                                `json:"version"`
		ChecklistID      string                               `json:"checklist_id"`
		ChecklistVersion int64                                `json:"checklist_version"`
		Responses        []franchisejourney.ChecklistResponse `json:"responses"`
	}
	if !decodeStrict(w, r, &input) {
		return
	}
	p, ok := a.protected(w, r, "handover:manage", input.OrganizationID)
	if !ok {
		return
	}
	value, err := a.service.CompleteDeliveryChecklist(r.Context(), p.TenantID, input.OrganizationID, p.Subject, r.PathValue("id"), input.Version, input.ChecklistID, input.ChecklistVersion, input.Responses)
	writeJourneyResult(w, value, err)
}

func (a franchiseJourneyAPI) deliveryExceptions(w http.ResponseWriter, r *http.Request) {
	organization := r.URL.Query().Get("organization_id")
	p, ok := a.protected(w, r, "handover:manage", organization)
	if !ok {
		return
	}
	items, err := a.service.DeliveryExceptions(r.Context(), p.TenantID, organization, 100)
	if err != nil {
		writeProblem(w, 500, "QUERY_FAILED", "delivery exceptions query failed")
		return
	}
	writeJSON(w, 200, map[string]any{"items": items})
}

func (a franchiseJourneyAPI) resolveDeliveryException(w http.ResponseWriter, r *http.Request) {
	var input struct {
		OrganizationID string `json:"organization_id"`
		Version        int64  `json:"version"`
		Action         string `json:"action"`
		Notes          string `json:"notes"`
	}
	if !decodeStrict(w, r, &input) {
		return
	}
	p, ok := a.protected(w, r, "handover:manage", input.OrganizationID)
	if !ok {
		return
	}
	value, err := a.service.ResolveDeliveryException(r.Context(), p.TenantID, input.OrganizationID, p.Subject, r.PathValue("id"), input.Version, input.Action, input.Notes)
	writeJourneyResult(w, value, err)
}

func (a franchiseJourneyAPI) returnCases(w http.ResponseWriter, r *http.Request) {
	organization := r.URL.Query().Get("organization_id")
	p, ok := a.protected(w, r, "handover:manage", organization)
	if !ok {
		return
	}
	items, err := a.service.ReturnCases(r.Context(), p.TenantID, organization, 100)
	if err != nil {
		writeProblem(w, 500, "QUERY_FAILED", "return cases query failed")
		return
	}
	writeJSON(w, 200, map[string]any{"items": items})
}

func (a franchiseJourneyAPI) receiveReturn(w http.ResponseWriter, r *http.Request) {
	var input struct {
		OrganizationID string `json:"organization_id"`
		SerialNumber   string `json:"serial_number"`
		ConditionCode  string `json:"condition_code"`
		Notes          string `json:"notes"`
		EvidenceSHA256 string `json:"evidence_sha256"`
	}
	if !decodeStrict(w, r, &input) {
		return
	}
	p, ok := a.protected(w, r, "handover:manage", input.OrganizationID)
	if !ok {
		return
	}
	value, err := a.service.ReceiveReturn(r.Context(), p.TenantID, input.OrganizationID, p.Subject, r.PathValue("id"), input.SerialNumber, input.ConditionCode, input.Notes, input.EvidenceSHA256)
	writeJourneyResult(w, value, err)
}

func (a franchiseJourneyAPI) decideReturn(w http.ResponseWriter, r *http.Request) {
	var input struct {
		OrganizationID  string `json:"organization_id"`
		InventoryAction string `json:"inventory_action"`
		Notes           string `json:"notes"`
	}
	if !decodeStrict(w, r, &input) {
		return
	}
	p, ok := a.protected(w, r, "handover:manage", input.OrganizationID)
	if !ok {
		return
	}
	value, err := a.service.DecideReturn(r.Context(), p.TenantID, input.OrganizationID, p.Subject, r.PathValue("id"), input.InventoryAction, input.Notes)
	writeJourneyResult(w, value, err)
}

func (a franchiseJourneyAPI) publicAppointmentSlots(w http.ResponseWriter, r *http.Request) {
	from, fromErr := time.Parse(time.RFC3339, r.URL.Query().Get("from"))
	to, toErr := time.Parse(time.RFC3339, r.URL.Query().Get("to"))
	if fromErr != nil || toErr != nil {
		writeProblem(w, 400, "INVALID_SLOT_RANGE", "from and to must be RFC3339 timestamps")
		return
	}
	items, err := a.service.PublicAppointmentSlots(r.Context(), r.PathValue("tenantCode"), r.PathValue("organizationCode"), r.URL.Query().Get("kind"), from, to)
	if errors.Is(err, franchisejourney.ErrInvalid) {
		writeProblem(w, 400, "INVALID_SLOT_QUERY", "slot query does not match contract")
		return
	}
	if err != nil {
		writeProblem(w, 500, "QUERY_FAILED", "appointment slot query failed")
		return
	}
	writeJSON(w, 200, map[string]any{"items": items})
}

func (a franchiseJourneyAPI) createAppointmentSlot(w http.ResponseWriter, r *http.Request) {
	var input struct {
		OrganizationID string    `json:"organization_id"`
		Kind           string    `json:"kind"`
		StartsAt       time.Time `json:"starts_at"`
		EndsAt         time.Time `json:"ends_at"`
		Capacity       int       `json:"capacity"`
	}
	hash, decoded := decodeHashedJSON(w, r, &input)
	if !decoded {
		return
	}
	p, ok := a.protected(w, r, "appointment:manage", input.OrganizationID)
	if !ok {
		return
	}
	value, replay, err := a.service.CreateAppointmentSlotOnce(r.Context(), p.TenantID, p.Subject, r.Header.Get("Idempotency-Key"), hash, franchisejourney.AppointmentSlot{OrganizationID: input.OrganizationID, Kind: input.Kind, StartsAt: input.StartsAt, EndsAt: input.EndsAt, Capacity: input.Capacity})
	if errors.Is(err, franchisejourney.ErrConflict) {
		writeProblem(w, 409, "SLOT_CONFLICT", "slot overlaps existing capacity or organization is unavailable")
		return
	}
	if errors.Is(err, franchisejourney.ErrInvalid) {
		writeProblem(w, 400, "INVALID_SLOT", "slot does not match contract")
		return
	}
	if err != nil {
		writeProblem(w, 500, "SLOT_FAILED", "appointment slot could not be persisted")
		return
	}
	status := 201
	if replay {
		status = 200
		w.Header().Set("Idempotency-Replayed", "true")
	}
	writeJSON(w, status, value)
}

func (a franchiseJourneyAPI) createServiceResource(w http.ResponseWriter, r *http.Request) {
	var input struct {
		OrganizationID   string   `json:"organization_id"`
		PrincipalSubject string   `json:"principal_subject"`
		DisplayName      string   `json:"display_name"`
		Kind             string   `json:"kind"`
		Skills           []string `json:"skills"`
	}
	hash, decoded := decodeHashedJSON(w, r, &input)
	if !decoded {
		return
	}
	p, ok := a.protected(w, r, "resource:manage", input.OrganizationID)
	if !ok {
		return
	}
	value, replay, err := a.service.CreateServiceResourceOnce(r.Context(), p.TenantID, p.Subject, r.Header.Get("Idempotency-Key"), hash, franchisejourney.ServiceResource{OrganizationID: input.OrganizationID, PrincipalSubject: input.PrincipalSubject, DisplayName: input.DisplayName, Kind: input.Kind, Skills: input.Skills})
	if errors.Is(err, franchisejourney.ErrConflict) {
		writeProblem(w, 409, "RESOURCE_CONFLICT", "resource identity or organization conflicts with this request")
		return
	}
	if errors.Is(err, franchisejourney.ErrInvalid) {
		writeProblem(w, 400, "INVALID_RESOURCE", "resource does not match contract")
		return
	}
	if err != nil {
		writeProblem(w, 500, "RESOURCE_FAILED", "resource could not be persisted")
		return
	}
	status := 201
	if replay {
		status = 200
		w.Header().Set("Idempotency-Replayed", "true")
	}
	writeJSON(w, status, value)
}

func (a franchiseJourneyAPI) assignAppointmentResource(w http.ResponseWriter, r *http.Request) {
	var input struct {
		OrganizationID string `json:"organization_id"`
		ResourceID     string `json:"resource_id"`
		Version        int64  `json:"version"`
	}
	if !decodeStrict(w, r, &input) {
		return
	}
	p, ok := a.protected(w, r, "appointment:manage", input.OrganizationID)
	if !ok {
		return
	}
	value, err := a.service.AssignAppointmentResource(r.Context(), p.TenantID, input.OrganizationID, r.PathValue("id"), input.ResourceID, input.Version)
	writeJourneyResult(w, value, err)
}

func (a franchiseJourneyAPI) transitionAppointment(w http.ResponseWriter, r *http.Request) {
	var input struct {
		OrganizationID string `json:"organization_id"`
		Current        string `json:"current"`
		Target         string `json:"target"`
		Version        int64  `json:"version"`
		ReasonCode     string `json:"reason_code"`
	}
	if !decodeStrict(w, r, &input) {
		return
	}
	p, ok := a.protected(w, r, "appointment:manage", input.OrganizationID)
	if !ok {
		return
	}
	value, err := a.service.TransitionAppointment(r.Context(), p.TenantID, input.OrganizationID, r.PathValue("id"), input.Current, input.Target, input.Version, p.Subject, input.ReasonCode)
	writeJourneyResult(w, value, err)
}

func (a franchiseJourneyAPI) createAvailability(w http.ResponseWriter, r *http.Request) {
	var input franchisejourney.AvailabilityEntry
	hash, ok := decodeHashedJSON(w, r, &input)
	if !ok {
		return
	}
	p, ok := a.protected(w, r, "availability:manage", input.OrganizationID)
	if !ok {
		return
	}
	value, replayed, err := a.service.CreateAvailabilityOnce(r.Context(), p.TenantID, p.Subject, r.Header.Get("Idempotency-Key"), hash, input)
	if errors.Is(err, franchisejourney.ErrConflict) {
		writeProblem(w, 409, "AVAILABILITY_CONFLICT", "availability overlaps or conflicts with active appointments")
		return
	}
	if errors.Is(err, franchisejourney.ErrInvalid) {
		writeProblem(w, 400, "INVALID_AVAILABILITY", "availability does not match contract")
		return
	}
	if err != nil {
		writeProblem(w, 500, "AVAILABILITY_FAILED", "availability could not be persisted")
		return
	}
	status := 201
	if replayed {
		status = 200
		w.Header().Set("Idempotency-Replayed", "true")
	}
	writeJSON(w, status, value)
}

func (a franchiseJourneyAPI) cancelAvailability(w http.ResponseWriter, r *http.Request) {
	var input struct {
		OrganizationID string `json:"organization_id"`
		Version        int64  `json:"version"`
		ReasonCode     string `json:"reason_code"`
	}
	if !decodeStrict(w, r, &input) {
		return
	}
	p, ok := a.protected(w, r, "availability:manage", input.OrganizationID)
	if !ok {
		return
	}
	value, err := a.service.CancelAvailability(r.Context(), p.TenantID, input.OrganizationID, r.PathValue("id"), input.Version, p.Subject, input.ReasonCode)
	writeJourneyResult(w, value, err)
}

func (a franchiseJourneyAPI) availability(w http.ResponseWriter, r *http.Request) {
	organization := r.URL.Query().Get("organization_id")
	p, ok := a.protected(w, r, "availability:read", organization)
	if !ok {
		return
	}
	from, fromErr := time.Parse(time.RFC3339, r.URL.Query().Get("from"))
	to, toErr := time.Parse(time.RFC3339, r.URL.Query().Get("to"))
	if fromErr != nil || toErr != nil {
		writeProblem(w, 400, "INVALID_AVAILABILITY_RANGE", "from and to must be RFC3339 timestamps")
		return
	}
	items, err := a.service.Availability(r.Context(), p.TenantID, organization, r.URL.Query().Get("resource_id"), from, to)
	if errors.Is(err, franchisejourney.ErrInvalid) {
		writeProblem(w, 400, "INVALID_AVAILABILITY_RANGE", "availability range does not match contract")
		return
	}
	if err != nil {
		writeProblem(w, 500, "AVAILABILITY_QUERY_FAILED", "availability query failed")
		return
	}
	writeJSON(w, 200, map[string]any{"items": items})
}

func (a franchiseJourneyAPI) cancelCustomerAppointment(w http.ResponseWriter, r *http.Request) {
	var input struct {
		OrganizationID string `json:"organization_id"`
		Version        int64  `json:"version"`
		ReasonCode     string `json:"reason_code"`
	}
	if !decodeStrict(w, r, &input) {
		return
	}
	p, ok := a.protected(w, r, "customer:self", input.OrganizationID)
	if !ok {
		return
	}
	value, err := a.service.CancelCustomerAppointment(r.Context(), p.TenantID, input.OrganizationID, p.Subject, r.PathValue("id"), input.Version, input.ReasonCode)
	writeJourneyResult(w, value, err)
}

type franchiseJourneyAPI struct {
	service  *franchisejourney.Service
	verifier identity.Verifier
}

func (a franchiseJourneyAPI) protected(w http.ResponseWriter, r *http.Request, permission, organization string) (identity.Principal, bool) {
	p, err := authenticate(r.Context(), r.Header.Get("Authorization"), a.verifier)
	if err != nil {
		writeProblem(w, 401, "UNAUTHENTICATED", "a valid bearer token is required")
		return p, false
	}
	if !p.Allowed(permission) {
		writeProblem(w, 403, "FORBIDDEN", permission+" permission is required")
		return p, false
	}
	if !p.AllowedOrganization(organization) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return p, false
	}
	return p, true
}

func (a franchiseJourneyAPI) publicLocations(w http.ResponseWriter, r *http.Request) {
	items, err := a.service.PublicLocations(r.Context(), r.PathValue("tenantCode"))
	if errors.Is(err, franchisejourney.ErrInvalid) {
		writeProblem(w, 400, "INVALID_TENANT", "tenant code is invalid")
		return
	}
	if err != nil {
		writeProblem(w, 500, "QUERY_FAILED", "locations query failed")
		return
	}
	writeJSON(w, 200, map[string]any{"items": items})
}

func (a franchiseJourneyAPI) appointmentAgenda(w http.ResponseWriter, r *http.Request) {
	organization := r.URL.Query().Get("organization_id")
	p, ok := a.protected(w, r, "appointment:manage", organization)
	if !ok {
		return
	}
	from, fromErr := time.Parse(time.RFC3339, r.URL.Query().Get("from"))
	to, toErr := time.Parse(time.RFC3339, r.URL.Query().Get("to"))
	if fromErr != nil || toErr != nil {
		writeProblem(w, 400, "INVALID_AGENDA_RANGE", "explicit RFC3339 range required")
		return
	}
	value, err := a.service.AppointmentAgenda(r.Context(), p.TenantID, organization, from, to)
	w.Header().Set("Cache-Control", "no-store")
	writeJourneyResult(w, value, err)
}

func (a franchiseJourneyAPI) requestAppointment(w http.ResponseWriter, r *http.Request) {
	var input franchisejourney.Appointment
	hash, ok := decodeHashedJSON(w, r, &input)
	if !ok {
		return
	}
	key := r.Header.Get("Idempotency-Key")
	value, replayed, err := a.service.RequestAppointment(r.Context(), r.PathValue("tenantCode"), r.PathValue("organizationCode"), key, hash, input)
	if errors.Is(err, franchisejourney.ErrNotFound) {
		writeProblem(w, 404, "LOCATION_NOT_FOUND", "public location not found")
		return
	}
	if errors.Is(err, franchisejourney.ErrConflict) {
		writeProblem(w, 409, "APPOINTMENT_CONFLICT", "idempotency key or appointment capacity conflicts with this request")
		return
	}
	if errors.Is(err, franchisejourney.ErrInvalid) {
		writeProblem(w, 400, "INVALID_APPOINTMENT", "appointment does not match contract")
		return
	}
	if err != nil {
		writeProblem(w, 500, "APPOINTMENT_FAILED", "appointment could not be persisted")
		return
	}
	status := 202
	if replayed {
		status = 200
		w.Header().Set("Idempotency-Replayed", "true")
	}
	writeJSON(w, status, value)
}

func (a franchiseJourneyAPI) leads(w http.ResponseWriter, r *http.Request) {
	organization := r.URL.Query().Get("organization_id")
	p, ok := a.protected(w, r, "lead:read", organization)
	if !ok {
		return
	}
	limit, after, ok := pageInput(w, r)
	if !ok {
		return
	}
	value, err := a.service.Leads(r.Context(), p.TenantID, organization, limit, after)
	if err != nil {
		writeProblem(w, 500, "QUERY_FAILED", "leads query failed")
		return
	}
	writeJSON(w, 200, value)
}

func (a franchiseJourneyAPI) assignLead(w http.ResponseWriter, r *http.Request) {
	var input struct {
		OrganizationID  string `json:"organization_id"`
		AssignedSubject string `json:"assigned_subject"`
		Version         int64  `json:"version"`
	}
	if !decodeStrict(w, r, &input) {
		return
	}
	p, ok := a.protected(w, r, "lead:assign", input.OrganizationID)
	if !ok {
		return
	}
	value, err := a.service.AssignLeadAs(r.Context(), p.TenantID, input.OrganizationID, r.PathValue("id"), input.AssignedSubject, input.Version, p.Subject)
	writeJourneyResult(w, value, err)
}

func (a franchiseJourneyAPI) transitionLead(w http.ResponseWriter, r *http.Request) {
	var input struct {
		OrganizationID string `json:"organization_id"`
		Current        string `json:"current"`
		Target         string `json:"target"`
		Version        int64  `json:"version"`
	}
	if !decodeStrict(w, r, &input) {
		return
	}
	p, ok := a.protected(w, r, "lead:update", input.OrganizationID)
	if !ok {
		return
	}
	value, err := a.service.TransitionLeadAs(r.Context(), p.TenantID, input.OrganizationID, r.PathValue("id"), input.Current, input.Target, input.Version, p.Subject)
	writeJourneyResult(w, value, err)
}

func (a franchiseJourneyAPI) createQuote(w http.ResponseWriter, r *http.Request) {
	var input franchisejourney.Quote
	hash, ok := decodeHashedJSON(w, r, &input)
	if !ok {
		return
	}
	p, ok := a.protected(w, r, "quote:write", input.OrganizationID)
	if !ok {
		return
	}
	value, replayed, err := a.service.CreateQuoteAs(r.Context(), p.TenantID, p.Subject, r.Header.Get("Idempotency-Key"), hash, input)
	if errors.Is(err, franchisejourney.ErrConflict) {
		writeProblem(w, 409, "PRICE_OR_LEAD_CONFLICT", "active server price or lead not found")
		return
	}
	if errors.Is(err, franchisejourney.ErrInvalid) {
		writeProblem(w, 400, "INVALID_QUOTE", "quote does not match contract")
		return
	}
	if err != nil {
		writeProblem(w, 500, "QUOTE_FAILED", "quote could not be persisted")
		return
	}
	status := 201
	if replayed {
		status = 200
	}
	writeJSON(w, status, value)
}

func (a franchiseJourneyAPI) customerJourney(w http.ResponseWriter, r *http.Request) {
	organization := r.URL.Query().Get("organization_id")
	p, ok := a.protected(w, r, "customer:self", organization)
	if !ok {
		return
	}
	value, err := a.service.CustomerJourney(r.Context(), p.TenantID, organization, p.Subject)
	if err != nil {
		writeProblem(w, 500, "QUERY_FAILED", "customer journey query failed")
		return
	}
	writeJSON(w, 200, value)
}

func (a franchiseJourneyAPI) acceptHandover(w http.ResponseWriter, r *http.Request) {
	var input struct {
		OrganizationID    string `json:"organization_id"`
		Version           int64  `json:"version"`
		ConfirmedReceived bool   `json:"confirmed_received"`
		SerialNumber      string `json:"serial_number"`
		ChecklistID       string `json:"checklist_id"`
		ChecklistVersion  int64  `json:"checklist_version"`
	}
	evidence, decoded := decodeHashedJSON(w, r, &input)
	if !decoded {
		return
	}
	p, ok := a.protected(w, r, "customer:self", input.OrganizationID)
	if !ok {
		return
	}
	if !input.ConfirmedReceived {
		writeProblem(w, 400, "HANDOVER_CONFIRMATION_REQUIRED", "confirmed_received must be true")
		return
	}
	value, err := a.service.AcceptHandover(r.Context(), p.TenantID, input.OrganizationID, p.Subject, r.PathValue("id"), input.Version, input.SerialNumber, input.ChecklistID, input.ChecklistVersion, evidence)
	writeJourneyResult(w, value, err)
}

func (a franchiseJourneyAPI) rejectHandover(w http.ResponseWriter, r *http.Request) {
	var input struct {
		OrganizationID string `json:"organization_id"`
		Version        int64  `json:"version"`
		ReasonCode     string `json:"reason_code"`
		Details        string `json:"details"`
	}
	evidence, decoded := decodeHashedJSON(w, r, &input)
	if !decoded {
		return
	}
	p, ok := a.protected(w, r, "customer:self", input.OrganizationID)
	if !ok {
		return
	}
	value, err := a.service.RejectHandover(r.Context(), p.TenantID, input.OrganizationID, p.Subject, r.PathValue("id"), input.Version, input.ReasonCode, input.Details, evidence)
	writeJourneyResult(w, value, err)
}

func (a franchiseJourneyAPI) acceptQuote(w http.ResponseWriter, r *http.Request) {
	var input struct {
		OrganizationID string `json:"organization_id"`
		Version        int64  `json:"version"`
	}
	evidence, ok := decodeHashedJSON(w, r, &input)
	if !ok {
		return
	}
	p, ok := a.protected(w, r, "customer:self", input.OrganizationID)
	if !ok {
		return
	}
	value, err := a.service.AcceptQuote(r.Context(), p.TenantID, input.OrganizationID, p.Subject, r.PathValue("id"), input.Version, evidence)
	writeJourneyResult(w, value, err)
}

func decodeHashedJSON(w http.ResponseWriter, r *http.Request, destination any) (string, bool) {
	if r.Header.Get("Content-Type") != "application/json" {
		writeProblem(w, 415, "UNSUPPORTED_MEDIA_TYPE", "Content-Type must be application/json")
		return "", false
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeProblem(w, 400, "INVALID_JSON", "body is invalid or too large")
		return "", false
	}
	r.Body = io.NopCloser(bytes.NewReader(body))
	if !decodeStrict(w, r, destination) {
		return "", false
	}
	sum := sha256.Sum256(body)
	return hex.EncodeToString(sum[:]), true
}

func writeJourneyResult(w http.ResponseWriter, value any, err error) {
	if errors.Is(err, franchisejourney.ErrConflict) {
		writeProblem(w, 409, "JOURNEY_CONFLICT", "state or version conflict")
		return
	}
	if errors.Is(err, franchisejourney.ErrInvalid) {
		writeProblem(w, 400, "INVALID_JOURNEY_COMMAND", "command does not match contract")
		return
	}
	if err != nil {
		writeProblem(w, 500, "JOURNEY_FAILED", "command failed")
		return
	}
	writeJSON(w, 200, value)
}

func (a franchiseJourneyAPI) quoteResult(w http.ResponseWriter, r *http.Request) {
	organization := r.URL.Query().Get("organization_id")
	p, ok := a.protected(w, r, "quote:write", organization)
	if !ok {
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	value, err := a.service.QuoteResult(r.Context(), p.TenantID, organization, r.URL.Query().Get("lead_id"), r.URL.Query().Get("request_key"))
	writeJourneyResult(w, value, err)
}

func (a franchiseJourneyAPI) availabilityCreationResult(w http.ResponseWriter, r *http.Request) {
	organization := r.URL.Query().Get("organization_id")
	p, ok := a.protected(w, r, "availability:manage", organization)
	if !ok {
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	value, err := a.service.AvailabilityCreationResult(r.Context(), p.TenantID, organization, r.URL.Query().Get("request_key"))
	writeJourneyResult(w, value, err)
}

func (a franchiseJourneyAPI) serviceResourceCreationResult(w http.ResponseWriter, r *http.Request) {
	organization := r.URL.Query().Get("organization_id")
	p, ok := a.protected(w, r, "resource:manage", organization)
	if !ok {
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	value, err := a.service.ServiceResourceCreationResult(r.Context(), p.TenantID, organization, r.URL.Query().Get("request_key"))
	writeJourneyResult(w, value, err)
}

func (a franchiseJourneyAPI) appointmentSlotCreationResult(w http.ResponseWriter, r *http.Request) {
	organization := r.URL.Query().Get("organization_id")
	p, ok := a.protected(w, r, "appointment:manage", organization)
	if !ok {
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	value, err := a.service.AppointmentSlotCreationResult(r.Context(), p.TenantID, organization, r.URL.Query().Get("request_key"))
	writeJourneyResult(w, value, err)
}

func (a franchiseJourneyAPI) publishedDeliveryChecklist(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	organization := r.URL.Query().Get("organization_id")
	p, ok := a.protected(w, r, "handover:manage", organization)
	if !ok {
		return
	}
	version, err := strconv.ParseInt(r.URL.Query().Get("version"), 10, 64)
	if err != nil || version < 1 {
		writeProblem(w, 400, "INVALID_DELIVERY_CHECKLIST_QUERY", "checklist version is invalid")
		return
	}
	value, err := a.service.PublishedDeliveryChecklist(r.Context(), p.TenantID, organization, r.URL.Query().Get("checklist_id"), version)
	writeJourneyResult(w, value, err)
}

func (a franchiseJourneyAPI) deliveryChecklistCompletion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	organization := r.URL.Query().Get("organization_id")
	p, ok := a.protected(w, r, "handover:manage", organization)
	if !ok {
		return
	}
	value, err := a.service.DeliveryChecklistCompletion(r.Context(), p.TenantID, organization, r.PathValue("id"))
	writeJourneyResult(w, value, err)
}

func (a franchiseJourneyAPI) returnCaseResult(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	organization := r.URL.Query().Get("organization_id")
	p, ok := a.protected(w, r, "handover:manage", organization)
	if !ok {
		return
	}
	value, err := a.service.ReturnCaseResult(r.Context(), p.TenantID, organization, r.URL.Query().Get("authorization_id"))
	writeJourneyResult(w, value, err)
}
````

### FILE: `internal/platform/httpapi/franchisejourney_test.go`

```yaml
block_id: "GO-FRANCHISE-JOURNEY:file:09"
operation: CREATE
provenance: AUTHORED
source: "local implementation governed by the upstream references in metadata"
license: "LicenseRef-Workspace-Owner"
sha256: "83b304013de77fbc41097704a4e59967c1b7b84a48825facba198c563cc39cef"
variables: []
secrets_allowed: false
```

````go
package httpapi

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"elite.local/enterprise/internal/commerce"
	"elite.local/enterprise/internal/franchisejourney"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/platform/randomid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type browserAppointmentClock struct{}

// AUTHORED local integration fixture. Real RS256/JWKS verification, not a
// successful verifier double. No live IdP login or production identity claim.
func confirmationTestIssuer(t *testing.T) (identity.Verifier, func(string, string, []string, []string) string) {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	encode := base64.RawURLEncoding.EncodeToString
	var issuer *httptest.Server
	issuer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/.well-known/openid-configuration":
			_ = json.NewEncoder(w).Encode(map[string]any{"issuer": issuer.URL, "jwks_uri": issuer.URL + "/keys", "id_token_signing_alg_values_supported": []string{"RS256"}})
		case "/keys":
			_ = json.NewEncoder(w).Encode(map[string]any{"keys": []any{map[string]any{"kty": "RSA", "use": "sig", "alg": "RS256", "kid": "local-confirmation", "n": encode(key.N.Bytes()), "e": encode(big.NewInt(int64(key.E)).Bytes())}}})
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(issuer.Close)
	verifier, err := identity.NewOIDCVerifier(context.Background(), issuer.URL, "confirmation-api")
	if err != nil {
		t.Fatal(err)
	}
	token := func(subject, tenant string, permissions, organizations []string) string {
		header, err := json.Marshal(map[string]string{"alg": "RS256", "kid": "local-confirmation", "typ": "JWT"})
		if err != nil {
			t.Fatal(err)
		}
		claims, err := json.Marshal(map[string]any{"iss": issuer.URL, "aud": "confirmation-api", "sub": subject, "iat": time.Now().Unix(), "exp": time.Now().Add(5 * time.Minute).Unix(), "tenant_id": tenant, "permissions": permissions, "organization_ids": organizations})
		if err != nil {
			t.Fatal(err)
		}
		unsigned := encode(header) + "." + encode(claims)
		digest := sha256.Sum256([]byte(unsigned))
		signature, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, digest[:])
		if err != nil {
			t.Fatal(err)
		}
		return unsigned + "." + encode(signature)
	}
	return verifier, token
}

func TestAppointmentConfirmationBFFPostgres(t *testing.T) {
	runAppointmentConfirmationPostgres(t, false)
}

func TestAppointmentAgendaBrowserPostgres(t *testing.T) {
	for _, project := range []string{"chromium-desktop", "chromium-mobile", "firefox-desktop", "webkit-desktop"} {
		t.Run(project, func(t *testing.T) {
			t.Setenv("ELITE_AGENDA_BROWSER_PROJECT", project)
			runAppointmentConfirmationPostgres(t, true)
		})
	}
}

func runAppointmentConfirmationPostgres(t *testing.T, browser bool) {
	t.Helper()
	if os.Getenv("ELITE_CONFIRMATION_E2E") != "1" {
		t.Skip("explicit disposable confirmation gate not requested")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	config, err := pgxpool.ParseConfig(os.Getenv("TEST_DATABASE_URL"))
	if err != nil {
		t.Fatal("invalid test database configuration")
	}
	// Immutable audit rows deliberately remain for inspection until the caller
	// discards this dedicated database. Never disable their protective triggers.
	if config.ConnConfig.Host != "127.0.0.1" || !strings.HasPrefix(config.ConnConfig.Database, "elite_confirmation_") {
		t.Fatal("requires a disposable loopback elite_confirmation_* database, never a project database")
	}
	web := os.Getenv("ELITE_WEB_ROOT")
	if !filepath.IsAbs(web) {
		t.Fatal("ELITE_WEB_ROOT must be absolute")
	}
	cli := filepath.Join(web, "node_modules", "vitest", "vitest.mjs")
	if _, err := os.Stat(cli); err != nil {
		t.Fatal("materialize web and install its exact frozen lock first")
	}
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := (randomid.Generator{}).New()
	code := "confirmation-" + strings.ReplaceAll(tenant, "-", "")
	if _, err := pool.Exec(ctx, "insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1,$2,'Synthetic confirmation','Synthetic confirmation')", tenant, code); err != nil {
		t.Fatal(err)
	}
	for _, sql := range []string{
		"insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'store','store','Synthetic Store','store')",
		"insert into org.public_location(tenant_id,organization_id,city,region,country,published)values($1,'store','Fixture City','Fixture Region','AR',true)",
		"insert into crm.customer_profile(tenant_id,customer_principal_id,display_name,email_normalized)values($1,'customer','Synthetic Customer','confirmation@example.invalid')",
		"insert into crm.lead(tenant_id,lead_id,organization_id,customer_principal_id,lifecycle_state,source_code,contact_payload)values($1,'lead','store','customer','new','public-web','{}')",
	} {
		if _, err := pool.Exec(ctx, sql, tenant); err != nil {
			t.Fatal(err)
		}
	}
	verifier, token := confirmationTestIssuer(t)
	mux := http.NewServeMux()
	FranchiseJourneyModule{Service: franchisejourney.NewService(postgres.NewFranchiseJourney(pool), randomid.Generator{}, browserAppointmentClock{})}.Register(mux, verifier)
	api := httptest.NewServer(mux)
	defer api.Close()
	operator := token("operator", tenant, []string{"appointment:manage", "resource:manage", "availability:manage"}, []string{"store"})
	post := func(path, bearer string, payload any, expected int) map[string]any {
		t.Helper()
		body, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequestWithContext(ctx, "POST", api.URL+path, strings.NewReader(string(body)))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+bearer)
		req.Header.Set("Idempotency-Key", "confirmation-request-key")
		if path == "/v1/franchise/availability" {
			req.Header.Set("Idempotency-Key", "availability-"+(randomid.Generator{}).New())
		}
		res, err := api.Client().Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		var value map[string]any
		if err := json.NewDecoder(res.Body).Decode(&value); err != nil {
			t.Fatal(err)
		}
		if res.StatusCode != expected {
			t.Fatalf("%s expected=%d actual=%d response=%v", path, expected, res.StatusCode, value)
		}
		return value
	}
	start := time.Now().UTC().Add(24 * time.Hour).Truncate(time.Second)
	end := start.Add(30 * time.Minute)
	working := map[string]any{"organization_id": "store", "entry_type": "working", "starts_at": start.Add(-time.Hour), "ends_at": end.Add(time.Hour)}
	post("/v1/franchise/availability", operator, working, 201)
	post("/v1/franchise/appointment-slots", operator, map[string]any{"organization_id": "store", "kind": "consultation", "starts_at": start, "ends_at": end, "capacity": 1}, 201)
	appointment := post("/v1/public/"+code+"/store/appointments", "", map[string]any{"lead_id": "lead", "kind": "consultation", "starts_at": start}, 202)
	id, ok := appointment["id"].(string)
	if !ok || id == "" {
		t.Fatal("missing appointment receipt")
	}
	transitionPath := "/v1/franchise/appointments/" + id + "/transitions"
	confirmation := map[string]any{"organization_id": "store", "current": "requested", "target": "confirmed", "version": 1}
	for _, negative := range []struct {
		name, bearer string
		status       int
	}{
		{"invalid-signature", operator[:strings.LastIndex(operator, ".")+1] + "AAAA", 401},
		{"no-permission", token("observer", tenant, []string{"lead:read"}, []string{"store"}), 403},
		{"other-organization", token("operator", tenant, []string{"appointment:manage"}, []string{"other"}), 403},
		{"other-tenant", token("operator", (randomid.Generator{}).New(), []string{"appointment:manage"}, []string{"store"}), 409},
		{"no-resource", operator, 409},
	} {
		t.Run(negative.name, func(t *testing.T) { post(transitionPath, negative.bearer, confirmation, negative.status) })
	}
	resource := post("/v1/franchise/resources", operator, map[string]any{"organization_id": "store", "principal_subject": "technician", "display_name": "Synthetic technician", "kind": "employee", "skills": []string{"consultation"}}, 201)
	assignment := map[string]any{"organization_id": "store", "resource_id": resource["id"], "version": 1}
	assignmentPath := "/v1/franchise/appointments/" + id + "/resources"
	post(assignmentPath, operator, assignment, 409) // no resource working window
	working["resource_id"] = resource["id"]
	post("/v1/franchise/availability", operator, working, 201)
	assignment["version"] = 2
	post(assignmentPath, operator, assignment, 409) // stale version, no effect
	assignment["version"] = 1
	if !browser {
		assigned := post(assignmentPath, operator, assignment, 200)
		if assigned["version"] != float64(2) || assigned["state"] != "requested" {
			t.Fatalf("invalid assignment receipt: %v", assigned)
		}
	}
	// Real server-side BFF client forwards signed tokens to this API. The test
	// races confirmations then reads the customer's durable, scoped timeline.
	cmd := exec.CommandContext(ctx, "node", cli, "run", "src/platform/backend/protected-client.test.ts", "-t", "connected appointment confirmation")
	cmd.Dir = web
	if browser {
		cmd = exec.CommandContext(ctx, "node", filepath.Join(web, "microsoft_playwright_browser_gate", "node_modules", "@playwright", "test", "cli.js"), "test", "tests/enterprise-web.spec.mjs", "--grep", "operator agenda confirms", "--project", os.Getenv("ELITE_AGENDA_BROWSER_PROJECT"), "--workers=1", "--retries=0", "--max-failures=1")
		cmd.Dir = filepath.Join(web, "microsoft_playwright_browser_gate")
	}
	for _, item := range os.Environ() {
		name := strings.ToUpper(strings.SplitN(item, "=", 2)[0])
		if name == "TEST_DATABASE_URL" || name == "DATABASE_URL" || strings.HasPrefix(name, "ELITE_CONFIRMATION_") || name == "ENTERPRISE_API_BASE_URL" || name == "AUTH_SESSION_SECRET" || name == "ELITE_BASE_URL" || name == "ELITE_RUNTIME_ONLY" || name == "ELITE_OPERATOR_AGENDA_E2E" || name == "ELITE_WEB_ROOT" {
			continue
		}
		cmd.Env = append(cmd.Env, item)
	}
	cmd.Env = append(cmd.Env, "ENTERPRISE_API_BASE_URL="+api.URL, "ELITE_CONFIRMATION_E2E=1", "ELITE_CONFIRMATION_ID="+id, "ELITE_CONFIRMATION_TENANT="+tenant, "ELITE_CONFIRMATION_OPERATOR="+operator, "ELITE_CONFIRMATION_CUSTOMER="+token("customer", tenant, []string{"customer:self"}, []string{"store"}), "ELITE_CONFIRMATION_STRANGER="+token("stranger", tenant, []string{"customer:self"}, []string{"store"}))
	if browser {
		cmd.Env = append(cmd.Env, "ELITE_OPERATOR_AGENDA_E2E=1", "ELITE_WEB_ROOT="+web, "APP_BASE_URL=https://127.0.0.1:4173", "ELITE_CONFIRMATION_DAY="+start.Format("2006-01-02"), "AUTH_SESSION_SECRET=synthetic-"+(randomid.Generator{}).New()+(randomid.Generator{}).New())
	}
	if browser {
		target, _ := url.Parse("http://127.0.0.1:4173")
		edge := httptest.NewTLSServer(httputil.NewSingleHostReverseProxy(target))
		defer edge.Close()
		filtered := cmd.Env[:0]
		for _, item := range cmd.Env {
			if !strings.HasPrefix(item, "APP_BASE_URL=") && !strings.HasPrefix(item, "ELITE_BASE_URL=") {
				filtered = append(filtered, item)
			}
		}
		cmd.Env = append(filtered, "APP_BASE_URL="+edge.URL, "ELITE_BASE_URL="+edge.URL)
		server := exec.CommandContext(ctx, "node", filepath.Join(web, "node_modules", "next", "dist", "bin", "next"), "start", "--hostname", "127.0.0.1", "--port", "4173")
		server.Dir = web
		server.Env = cmd.Env
		if err := server.Start(); err != nil {
			t.Fatal(err)
		}
		defer func() { _ = server.Process.Kill(); _ = server.Wait() }()
		client := &http.Client{Timeout: time.Second}
		ready := false
		for deadline := time.Now().Add(15 * time.Second); time.Now().Before(deadline); {
			response, err := client.Get(target.String() + "/icon.svg")
			if err == nil {
				response.Body.Close()
				if response.StatusCode == 200 {
					ready = true
					break
				}
			}
			select {
			case <-ctx.Done():
				t.Fatal("web fixture startup cancelled")
			case <-time.After(100 * time.Millisecond):
			}
		}
		if !ready {
			t.Fatal("web fixture did not become ready")
		}
	}
	output, err := cmd.CombinedOutput()
	// Never log signed fixture tokens, even on assertion failures.
	clean := string(output)
	for _, item := range cmd.Env {
		parts := strings.SplitN(item, "=", 2)
		if strings.HasPrefix(parts[0], "ELITE_CONFIRMATION_") && len(parts) == 2 && strings.Count(parts[1], ".") == 2 {
			clean = strings.ReplaceAll(clean, parts[1], "[REDACTED_TEST_TOKEN]")
		}
	}
	t.Log(clean)
	if err != nil {
		t.Fatalf("BFF confirmation gate failed: %v", err)
	}
	var state, actor string
	var version int
	if err := pool.QueryRow(ctx, "select a.state,a.version,x.actor_subject from crm.appointment a join crm.appointment_transition x on x.tenant_id=a.tenant_id and x.appointment_id=a.appointment_id where a.tenant_id=$1 and a.appointment_id=$2", tenant, id).Scan(&state, &version, &actor); err != nil {
		t.Fatal(err)
	}
	if state != "confirmed" || version != 3 || actor != "operator" {
		t.Fatalf("durable confirmation differs: %s %d %s", state, version, actor)
	}
	for _, sql := range []string{
		"select count(*) from crm.appointment_transition where tenant_id=$1 and appointment_id=$2 and from_state='requested' and to_state='confirmed'",
		"select count(*) from platform.outbox_event where tenant_id=$1 and aggregate_id=$2 and event_type='appointment.confirmed'",
		"select count(*) from crm.appointment_resource where tenant_id=$1 and appointment_id=$2",
	} {
		var count int
		if err := pool.QueryRow(ctx, sql, tenant, id).Scan(&count); err != nil {
			t.Fatal(err)
		}
		if count != 1 {
			t.Fatalf("expected exactly one durable effect, got %d", count)
		}
	}
	if _, err := pool.Exec(ctx, "delete from crm.appointment_transition where tenant_id=$1 and appointment_id=$2", tenant, id); err == nil {
		t.Fatal("immutable audit accepted deletion")
	}
	readRepo := postgres.NewFranchiseJourney(pool)
	agenda, err := readRepo.AppointmentAgenda(ctx, tenant, "store", start.Add(-time.Hour), end.Add(time.Hour))
	if err != nil || len(agenda.Appointments) != 1 || len(agenda.Resources) != 1 || agenda.Appointments[0].State != "confirmed" || agenda.Resources[0].PrincipalSubject != "" || agenda.Truncated {
		t.Fatalf("agenda snapshot mismatch: %v", err)
	}
	for _, scope := range [][2]string{{tenant, "other"}, {(randomid.Generator{}).New(), "store"}} {
		isolated, err := readRepo.AppointmentAgenda(ctx, scope[0], scope[1], start.Add(-time.Hour), end.Add(time.Hour))
		if err != nil || len(isolated.Appointments) != 0 || len(isolated.Resources) != 0 {
			t.Fatal("agenda scope leaked or query failed")
		}
	}
	if _, err := pool.Exec(ctx, `insert into crm.service_resource(tenant_id,resource_id,organization_id,display_name,resource_kind,status,version) select $1,'overflow-'||n,'store','Synthetic overflow '||n,'service-bay','active',1 from generate_series(1,201) n`, tenant); err != nil {
		t.Fatal(err)
	}
	bounded, err := readRepo.AppointmentAgenda(ctx, tenant, "store", start.Add(-time.Hour), end.Add(time.Hour))
	if err != nil || !bounded.Truncated || len(bounded.Resources) != 200 {
		t.Fatal("agenda silently truncated or exceeded resource budget")
	}
	if browser {
		t.Log("APPOINTMENT_AGENDA_BROWSER_POSTGRES_PASS signature=RS256 state=confirmed version=3 actor=operator audit=1 outbox=1 scoped_read=true bounded_read=true")
	} else {
		t.Log("APPOINTMENT_CONFIRMATION_BFF_POSTGRES_PASS signature=RS256 state=confirmed version=3 actor=operator confirmations=1 audit=1 outbox=1 duplicate=409 customer_scoped=true")
	}
}

func (browserAppointmentClock) Now() time.Time { return time.Now() }

func TestPublicAppointmentBrowserPostgres(t *testing.T) {
	runPublicBrowserPostgres(t, publicBrowserFixture{
		testName:      "public appointment follows captured lead",
		cleanupTables: []string{"crm.appointment", "crm.appointment_slot", "crm.availability_entry", "org.public_location"},
		setup: func(ctx context.Context, t *testing.T, pool *pgxpool.Pool, tenant string) []EnterpriseModule {
			if _, err := pool.Exec(ctx, `insert into org.public_location(tenant_id,organization_id,city,region,country,published)values($1,'browser-org','Fixture City','Fixture Region','AR',true)`, tenant); err != nil {
				t.Fatal(err)
			}
			start := time.Now().UTC().Add(24 * time.Hour).Truncate(time.Minute)
			if _, err := pool.Exec(ctx, `insert into crm.availability_entry(tenant_id,availability_id,organization_id,entry_type,starts_at,ends_at,state,version,created_by_subject)values($1,'browser-working','browser-org','working',$2,$3,'active',1,'synthetic-fixture')`, tenant, start.Add(-time.Hour), start.Add(5*time.Hour)); err != nil {
				t.Fatal(err)
			}
			for i, project := range []string{"chromium-desktop", "chromium-mobile", "firefox-desktop", "webkit-desktop"} {
				at := start.Add(time.Duration(i) * time.Hour)
				if _, err := pool.Exec(ctx, `insert into crm.appointment_slot(tenant_id,slot_id,organization_id,appointment_kind,starts_at,ends_at,capacity,state,version)values($1,$2,'browser-org','consultation',$3,$4,1,'open',1)`, tenant, "browser-slot-"+project, at, at.Add(30*time.Minute)); err != nil {
					t.Fatal(err)
				}
			}
			return []EnterpriseModule{FranchiseJourneyModule{Service: franchisejourney.NewService(postgres.NewFranchiseJourney(pool), randomid.Generator{}, browserAppointmentClock{})}}
		},
		verify: func(ctx context.Context, t *testing.T, pool *pgxpool.Pool, tenant string) {
			var appointments, events, receipts, violations int
			for _, check := range []struct {
				sql    string
				target *int
			}{
				{`select count(*) from crm.appointment a join crm.lead l on l.tenant_id=a.tenant_id and l.lead_id=a.lead_id join crm.appointment_slot s on s.tenant_id=a.tenant_id and s.slot_id=a.slot_id where a.tenant_id=$1 and a.organization_id=l.organization_id and a.model_id=l.model_id and a.starts_at=s.starts_at and a.ends_at=s.ends_at and a.state='requested' and a.version=1`, &appointments},
				{`select count(*) from platform.outbox_event where tenant_id=$1 and event_type='appointment.requested'`, &events},
				{`select count(*) from platform.idempotency_record where tenant_id=$1 and scope='public-appointment' and status='completed'`, &receipts},
				{`select count(*) from (select s.slot_id from crm.appointment_slot s left join crm.appointment a on a.tenant_id=s.tenant_id and a.slot_id=s.slot_id and a.state in ('requested','confirmed') where s.tenant_id=$1 group by s.slot_id,s.capacity having count(a.appointment_id)<>s.capacity) q`, &violations},
			} {
				if err := pool.QueryRow(ctx, check.sql, tenant).Scan(check.target); err != nil {
					t.Fatal(err)
				}
			}
			if appointments != 4 || events != 4 || receipts != 4 || violations != 0 {
				t.Fatalf("appointment invariants: appointments=%d events=%d receipts=%d capacity_violations=%d", appointments, events, receipts, violations)
			}
			t.Log("PUBLIC_APPOINTMENT_BROWSER_POSTGRES_PASS browsers=4 appointments=4 outbox=4 receipts=4 capacity_violations=0")
		},
	})
}

type journeyIDs struct{ n int }

func (i *journeyIDs) New() string { i.n++; return "journey-id" }

type journeyClock struct{}

func (journeyClock) Now() time.Time { return time.Date(2026, 8, 29, 20, 0, 0, 0, time.UTC) }

type journeyVerifier struct {
	principal identity.Principal
	err       error
}

func (v journeyVerifier) Verify(context.Context, string) (identity.Principal, error) {
	return v.principal, v.err
}

type journeyRepo struct {
	readFailure  bool
	appointment  franchisejourney.Appointment
	slot         franchisejourney.AppointmentSlot
	customer     string
	conflict     bool
	resource     franchisejourney.ServiceResource
	resourced    franchisejourney.Appointment
	availability franchisejourney.AvailabilityEntry
	quoteActor   string
	quoteKey     string
	quoteHash    string
	quoteReplay  bool
}

func (r *journeyRepo) PublicLocations(context.Context, string) ([]franchisejourney.Location, error) {
	return []franchisejourney.Location{{Code: "store", Name: "Store"}}, nil
}
func (r *journeyRepo) PublicAppointmentSlots(context.Context, string, string, string, time.Time, time.Time) ([]franchisejourney.AppointmentSlot, error) {
	return []franchisejourney.AppointmentSlot{{ID: "slot", Capacity: 2}}, nil
}
func (r *journeyRepo) CreateAppointmentSlotOnce(ctx context.Context, tenant, subject, key, hash string, v franchisejourney.AppointmentSlot, event string) (franchisejourney.AppointmentSlot, bool, error) {
	value, err := r.CreateAppointmentSlot(ctx, tenant, v, event)
	return value, false, err
}
func (r *journeyRepo) CreateAppointmentSlot(_ context.Context, _ string, value franchisejourney.AppointmentSlot, _ string) (franchisejourney.AppointmentSlot, error) {
	r.slot = value
	return value, nil
}
func (r *journeyRepo) RequestAppointment(_ context.Context, _, _, _ string, v franchisejourney.Appointment, _, _ string) (franchisejourney.Appointment, bool, error) {
	r.appointment = v
	return v, false, nil
}
func (r *journeyRepo) CreateServiceResourceOnce(ctx context.Context, tenant, subject, key, hash string, v franchisejourney.ServiceResource, event string) (franchisejourney.ServiceResource, bool, error) {
	value, err := r.CreateServiceResource(ctx, tenant, v, event)
	return value, false, err
}
func (r *journeyRepo) CreateServiceResource(_ context.Context, _ string, v franchisejourney.ServiceResource, _ string) (franchisejourney.ServiceResource, error) {
	r.resource = v
	return v, nil
}
func (r *journeyRepo) AssignAppointmentResource(_ context.Context, _, _, appointment, resource string, version int64, _ string) (franchisejourney.Appointment, error) {
	r.resourced = franchisejourney.Appointment{ID: appointment, ResourceID: resource, State: "requested", Version: version + 1}
	return r.resourced, nil
}
func (r *journeyRepo) CreateAvailability(_ context.Context, _, _ string, value franchisejourney.AvailabilityEntry, _ string) (franchisejourney.AvailabilityEntry, error) {
	r.availability = value
	return value, nil
}
func (r *journeyRepo) CancelAvailability(_ context.Context, _, _, entry string, version int64, _, _, _ string) (franchisejourney.AvailabilityEntry, error) {
	return franchisejourney.AvailabilityEntry{ID: entry, State: "cancelled", Version: version + 1}, nil
}
func (r *journeyRepo) Availability(context.Context, string, string, string, time.Time, time.Time) ([]franchisejourney.AvailabilityEntry, error) {
	return []franchisejourney.AvailabilityEntry{{ID: "working-window", EntryType: "working", State: "active", Version: 1}}, nil
}
func (r *journeyRepo) TransitionAppointment(_ context.Context, _, _, appointment, _, target string, version int64, _, _, _ string) (franchisejourney.Appointment, error) {
	return franchisejourney.Appointment{ID: appointment, State: target, Version: version + 1}, nil
}
func (r *journeyRepo) CancelCustomerAppointment(_ context.Context, _, _, customer, appointment string, version int64, _, _, _ string) (franchisejourney.Appointment, error) {
	r.customer = customer
	return franchisejourney.Appointment{ID: appointment, State: "cancelled", Version: version + 1}, nil
}
func (r *journeyRepo) Leads(context.Context, string, string, int, string) (franchisejourney.Page[franchisejourney.Lead], error) {
	return franchisejourney.Page[franchisejourney.Lead]{Items: []franchisejourney.Lead{{ID: "lead", Version: 1}}}, nil
}
func (r *journeyRepo) AssignLead(context.Context, string, string, string, string, int64, string) (franchisejourney.Lead, error) {
	if r.conflict {
		return franchisejourney.Lead{}, franchisejourney.ErrConflict
	}
	return franchisejourney.Lead{ID: "lead", Version: 2}, nil
}

func (r *journeyRepo) AssignLeadAs(ctx context.Context, tenant, organization, lead, subject string, version int64, eventID, actor string) (franchisejourney.Lead, error) {
	if actor == "" {
		return franchisejourney.Lead{}, franchisejourney.ErrInvalid
	}
	return r.AssignLead(ctx, tenant, organization, lead, subject, version, eventID)
}
func (r *journeyRepo) TransitionLeadAs(ctx context.Context, tenant, organization, lead, current, target string, version int64, eventID, actor string) (franchisejourney.Lead, error) {
	if actor == "" {
		return franchisejourney.Lead{}, franchisejourney.ErrInvalid
	}
	return r.TransitionLead(ctx, tenant, organization, lead, current, target, version, eventID)
}
func (r *journeyRepo) TransitionLead(context.Context, string, string, string, string, string, int64, string) (franchisejourney.Lead, error) {
	return franchisejourney.Lead{ID: "lead", State: "contacted", Version: 2}, nil
}
func (r *journeyRepo) CreateQuote(_ context.Context, _, key string, v franchisejourney.Quote, hash, _ string) (franchisejourney.Quote, bool, error) {
	r.quoteKey = key
	r.quoteHash = hash
	v.Currency = "ARS"
	v.TotalMinorUnits = 1000
	return v, r.quoteReplay, nil
}
func (r *journeyRepo) AcceptQuote(_ context.Context, _, _, customer, _ string, version int64, _ string, orderID, _, _, _ string) (franchisejourney.Quote, error) {
	r.customer = customer
	return franchisejourney.Quote{State: "accepted", Version: version + 1, OrderID: orderID}, nil
}
func (r *journeyRepo) PublishDeliveryChecklist(_ context.Context, _, subject string, value franchisejourney.DeliveryChecklist, _ string) (franchisejourney.DeliveryChecklist, error) {
	r.customer = subject
	value.State = "published"
	return value, nil
}
func (r *journeyRepo) CompleteDeliveryChecklist(_ context.Context, _, _, subject, _ string, version int64, checklistID string, checklistVersion int64, _ []franchisejourney.ChecklistResponse, _ string) (franchisejourney.Handover, error) {
	r.customer = subject
	return franchisejourney.Handover{State: "presented", Version: version + 1, ChecklistID: checklistID, ChecklistVersion: checklistVersion}, nil
}
func (r *journeyRepo) RejectHandover(_ context.Context, _, organization, customer, handover string, _ int64, reason, details, _ string, exceptionID, _ string) (franchisejourney.DeliveryException, error) {
	r.customer = customer
	return franchisejourney.DeliveryException{ID: exceptionID, OrganizationID: organization, HandoverID: handover, CustomerSubject: customer, ReasonCode: reason, Details: details, State: "open", Version: 1}, nil
}
func (r *journeyRepo) DeliveryExceptions(context.Context, string, string, int) ([]franchisejourney.DeliveryException, error) {
	return []franchisejourney.DeliveryException{{ID: "exception", State: "open", Version: 1}}, nil
}
func (r *journeyRepo) ResolveDeliveryException(_ context.Context, _, _, subject, exceptionID string, version int64, action, _ string, successorID, authorizationID, _, _ string) (franchisejourney.DeliveryResolution, error) {
	r.customer = subject
	result := franchisejourney.DeliveryResolution{Exception: franchisejourney.DeliveryException{ID: exceptionID, State: "resolved", Version: version + 1, ResolutionAction: action}}
	if action == "correct-and-represent" {
		result.SuccessorHandover = &franchisejourney.Handover{ID: successorID, State: "prepared", Version: 1}
	} else {
		result.ReturnAuthorizationID = authorizationID
		result.Disposition = action
	}
	return result, nil
}
func (r *journeyRepo) ReturnCases(context.Context, string, string, int) ([]franchisejourney.ReturnCase, error) {
	return []franchisejourney.ReturnCase{{AuthorizationID: "authorization", AuthorizedAction: "return"}}, nil
}
func (r *journeyRepo) ReceiveReturn(_ context.Context, _, organization, subject, authorizationID, serial, condition, notes, evidence, receiptID, _ string) (franchisejourney.ReturnReceipt, error) {
	r.customer = subject
	return franchisejourney.ReturnReceipt{ID: receiptID, AuthorizationID: authorizationID, OrganizationID: organization, ReceivedSerialNumber: serial, ConditionCode: condition, Notes: notes, EvidenceSHA256: evidence, ReceivedBySubject: subject}, nil
}
func (r *journeyRepo) DecideReturn(_ context.Context, _, _, subject, receiptID, inventoryAction, notes, dispositionID, inventoryID, remedyID, accountingID, _, _ string) (franchisejourney.ReturnDisposition, error) {
	r.customer = subject
	return franchisejourney.ReturnDisposition{ID: dispositionID, ReceiptID: receiptID, InventoryAction: inventoryAction, CustomerRemedy: "refund", Notes: notes, Effects: []franchisejourney.ReturnEffectRequest{{ID: inventoryID}, {ID: remedyID}, {ID: accountingID}}}, nil
}
func (r *journeyRepo) CustomerJourney(_ context.Context, _, _, customer string) (franchisejourney.CustomerJourney, error) {
	r.customer = customer
	if r.readFailure {
		return franchisejourney.CustomerJourney{Handovers: []franchisejourney.Handover{{ID: "PRIVATE-PARTIAL-RESULT"}}}, errors.New("PRIVATE-INTERNAL-DETAIL")
	}
	return franchisejourney.CustomerJourney{}, nil
}

func TestCustomerJourneyReadFailureIsRedactedAndRecoverable(t *testing.T) {
	verifier, token := confirmationTestIssuer(t)
	repo := &journeyRepo{readFailure: true}
	mux := http.NewServeMux()
	FranchiseJourneyModule{Service: franchisejourney.NewService(repo, randomid.Generator{}, browserAppointmentClock{})}.Register(mux, verifier)
	for _, tc := range []struct {
		name, bearer string
		status       int
	}{
		{"anonymous", "", 401},
		{"wrong-permission", token("customer", "018f4d4a-7b36-7a21-8d10-2f4c54c29b01", []string{"admin:read"}, []string{"store"}), 403},
		{"wrong-organization", token("customer", "018f4d4a-7b36-7a21-8d10-2f4c54c29b01", []string{"customer:self"}, []string{"other"}), 403},
		{"read-failure", token("customer", "018f4d4a-7b36-7a21-8d10-2f4c54c29b01", []string{"customer:self"}, []string{"store"}), 500},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/v1/customer/journey?organization_id=store", nil)
			if tc.bearer != "" {
				req.Header.Set("Authorization", "Bearer "+tc.bearer)
			}
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)
			if w.Code != tc.status || strings.Contains(w.Body.String(), "PRIVATE-") {
				t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
			}
			if tc.status == 500 && !strings.Contains(w.Body.String(), "QUERY_FAILED") {
				t.Fatal("missing controlled error")
			}
		})
	}
	repo.readFailure = false
	req := httptest.NewRequest("GET", "/v1/customer/journey?organization_id=store", nil)
	req.Header.Set("Authorization", "Bearer "+token("customer", "018f4d4a-7b36-7a21-8d10-2f4c54c29b01", []string{"customer:self"}, []string{"store"}))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != 200 || repo.customer != "customer" {
		t.Fatalf("recovery status=%d", w.Code)
	}
}

func (r *journeyRepo) AppointmentAgenda(context.Context, string, string, time.Time, time.Time) (franchisejourney.AppointmentAgenda, error) {
	return franchisejourney.AppointmentAgenda{Appointments: []franchisejourney.Appointment{}, Resources: []franchisejourney.ServiceResource{}}, nil
}

func TestAppointmentAgendaAuthorizationAndRange(t *testing.T) {
	principal := identity.Principal{TenantID: "tenant", Subject: "operator", Permissions: map[string]struct{}{"appointment:manage": {}}, Organizations: map[string]struct{}{"store": {}}}
	valid := "/v1/franchise/agenda?organization_id=store&from=2026-09-06T00:00:00Z&to=2026-09-07T00:00:00Z"
	for _, tc := range []struct {
		name, path string
		principal  identity.Principal
		err        error
		status     int
	}{
		{"authorized", valid, principal, nil, 200},
		{"unauthenticated", valid, principal, identity.ErrUnauthenticated, 401},
		{"wrong-role", valid, identity.Principal{TenantID: "tenant", Organizations: principal.Organizations}, nil, 403},
		{"wrong-organization", strings.Replace(valid, "organization_id=store", "organization_id=other", 1), principal, nil, 403},
		{"missing-dates", "/v1/franchise/agenda?organization_id=store", principal, nil, 400},
		{"unbounded", strings.Replace(valid, "2026-09-07", "2026-12-07", 1), principal, nil, 400},
		{"reversed", strings.Replace(valid, "2026-09-07", "2026-09-05", 1), principal, nil, 400},
	} {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			journeyHandler(&journeyRepo{}, tc.principal, tc.err).ServeHTTP(w, request("GET", tc.path, ""))
			if w.Code != tc.status {
				t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
			}
			if tc.status == 200 && (w.Header().Get("Cache-Control") != "no-store" || !strings.Contains(w.Body.String(), "\"appointments\":[]")) {
				t.Fatal("read contract/cache mismatch")
			}
		})
	}
}
func (r *journeyRepo) AcceptHandover(_ context.Context, _, _, customer, _ string, _ int64, _, checklistID string, checklistVersion int64, evidence, _ string) (franchisejourney.Handover, error) {
	r.customer = customer
	return franchisejourney.Handover{State: "accepted", AcceptanceEvidence: evidence, ChecklistID: checklistID, ChecklistVersion: checklistVersion}, nil
}

func journeyHandler(repository *journeyRepo, principal identity.Principal, verifyErr error) http.Handler {
	service := franchisejourney.NewService(repository, &journeyIDs{}, journeyClock{})
	mux := http.NewServeMux()
	FranchiseJourneyModule{Service: service}.Register(mux, journeyVerifier{principal: principal, err: verifyErr})
	return mux
}
func request(method, path, body string) *http.Request {
	r := httptest.NewRequest(method, path, strings.NewReader(body))
	if body != "" {
		r.Header.Set("Content-Type", "application/json")
	}
	r.Header.Set("Authorization", "Bearer token")
	return r
}

func TestPublicLocationsAndAppointment(t *testing.T) {
	repository := &journeyRepo{}
	handler := journeyHandler(repository, identity.Principal{}, nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, request("GET", "/v1/public/tenant/locations", ""))
	if w.Code != 200 || !strings.Contains(w.Body.String(), "Store") {
		t.Fatalf("locations status=%d body=%s", w.Code, w.Body.String())
	}
	body := `{"lead_id":"lead","model_id":"model","kind":"test-drive","starts_at":"2026-08-29T21:00:00Z"}`
	r := request("POST", "/v1/public/tenant/store/appointments", body)
	r.Header.Set("Idempotency-Key", "appointment-key-1")
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, r)
	if w.Code != 202 || repository.appointment.State != "requested" {
		t.Fatalf("appointment status=%d body=%s value=%+v", w.Code, w.Body.String(), repository.appointment)
	}
}

func TestCreateQuoteCarriesIdempotencyAndReportsReplay(t *testing.T) {
	repository := &journeyRepo{}
	principal := identity.Principal{Subject: "writer", TenantID: "tenant", Permissions: map[string]struct{}{"quote:write": {}}, Organizations: map[string]struct{}{"store": {}}}
	handler := journeyHandler(repository, principal, nil)
	body := `{"organization_id":"store","lead_id":"lead","variant_id":"variant","price_book_id":"retail","valid_until":"2026-08-30T20:00:00Z"}`
	r := request("POST", "/v1/franchise/quotes", body)
	r.Header.Set("Idempotency-Key", "quote-request-0001")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)
	if w.Code != 201 || repository.quoteActor != "writer" || repository.quoteKey != "quote-request-0001" || len(repository.quoteHash) != 64 {
		t.Fatalf("quote status=%d key=%q hash=%q body=%s", w.Code, repository.quoteKey, repository.quoteHash, w.Body.String())
	}
	noKey := httptest.NewRecorder()
	handler.ServeHTTP(noKey, request("POST", "/v1/franchise/quotes", body))
	if noKey.Code != 400 {
		t.Fatalf("missing key status=%d", noKey.Code)
	}
	repository.quoteReplay = true
	r = request("POST", "/v1/franchise/quotes", body)
	r.Header.Set("Idempotency-Key", "quote-request-0001")
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, r)
	if w.Code != 200 {
		t.Fatalf("quote replay status=%d body=%s", w.Code, w.Body.String())
	}
	r = request("POST", "/v1/franchise/quotes", body)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, r)
	if w.Code != 400 {
		t.Fatalf("quote without key status=%d body=%s", w.Code, w.Body.String())
	}
}

func TestPublicAndProtectedAppointmentCapacity(t *testing.T) {
	repository := &journeyRepo{}
	principal := identity.Principal{Subject: "slot-manager", TenantID: "tenant", Permissions: map[string]struct{}{"appointment:manage": {}}, Organizations: map[string]struct{}{"store": {}}}
	handler := journeyHandler(repository, principal, nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, request("GET", "/v1/public/tenant/store/appointment-slots?kind=service&from=2026-08-29T20:01:00Z&to=2026-08-30T20:00:00Z", ""))
	if w.Code != 200 || !strings.Contains(w.Body.String(), "slot") {
		t.Fatalf("public slots status=%d body=%s", w.Code, w.Body.String())
	}
	w = httptest.NewRecorder()
	slotRequest := request("POST", "/v1/franchise/appointment-slots", `{"organization_id":"store","kind":"service","starts_at":"2026-08-29T21:00:00Z","ends_at":"2026-08-29T22:00:00Z","capacity":2}`)
	slotRequest.Header.Set("Idempotency-Key", "slot-unit-fixture-key")
	handler.ServeHTTP(w, slotRequest)
	if w.Code != 201 || repository.slot.Capacity != 2 || repository.slot.State != "open" {
		t.Fatalf("create slot status=%d body=%s slot=%+v", w.Code, w.Body.String(), repository.slot)
	}
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, request("POST", "/v1/franchise/appointment-slots", `{"organization_id":"store","kind":"service","starts_at":"2026-08-29T21:00:00Z","ends_at":"2026-08-29T22:00:00Z","capacity":2,"state":"open"}`))
	if w.Code != 400 {
		t.Fatalf("server-managed slot field accepted status=%d", w.Code)
	}
}

func TestProtectedAppointmentResourceOperations(t *testing.T) {
	repository := &journeyRepo{}
	principal := identity.Principal{Subject: "resource-manager", TenantID: "tenant", Permissions: map[string]struct{}{"resource:manage": {}, "appointment:manage": {}}, Organizations: map[string]struct{}{"store": {}}}
	handler := journeyHandler(repository, principal, nil)
	w := httptest.NewRecorder()
	resourceRequest := request("POST", "/v1/franchise/resources", `{"organization_id":"store","principal_subject":"technician","display_name":"Technician","kind":"employee","skills":["service"]}`)
	resourceRequest.Header.Set("Idempotency-Key", "resource-unit-fixture-key")
	handler.ServeHTTP(w, resourceRequest)
	if w.Code != 201 || repository.resource.Status != "active" || repository.resource.ID == "" {
		t.Fatalf("resource status=%d body=%s value=%+v", w.Code, w.Body.String(), repository.resource)
	}
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, request("POST", "/v1/franchise/appointments/appointment/resources", `{"organization_id":"store","resource_id":"resource","version":1}`))
	if w.Code != 200 || repository.resourced.ResourceID != "resource" {
		t.Fatalf("assignment status=%d body=%s value=%+v", w.Code, w.Body.String(), repository.resourced)
	}
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, request("POST", "/v1/franchise/resources", `{"organization_id":"store","display_name":"Bay","kind":"service-bay","skills":["service"],"status":"active"}`))
	if w.Code != 400 {
		t.Fatalf("server-managed resource state accepted status=%d", w.Code)
	}
}

func TestProtectedScopeAndCustomerSubject(t *testing.T) {
	repository := &journeyRepo{}
	principal := identity.Principal{Subject: "customer-1", TenantID: "tenant", Permissions: map[string]struct{}{"customer:self": {}}, Organizations: map[string]struct{}{"store": {}}}
	handler := journeyHandler(repository, principal, nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, request("GET", "/v1/customer/journey?organization_id=other", ""))
	if w.Code != 403 {
		t.Fatalf("cross organization status=%d", w.Code)
	}
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, request("GET", "/v1/customer/journey?organization_id=store", ""))
	if w.Code != 200 || repository.customer != "customer-1" {
		t.Fatalf("customer status=%d bound=%s", w.Code, repository.customer)
	}
}

func TestAvailabilityAndCustomerCancellationBindActorAndScope(t *testing.T) {
	repository := &journeyRepo{}
	principal := identity.Principal{Subject: "scheduler", TenantID: "tenant", Permissions: map[string]struct{}{"availability:manage": {}, "availability:read": {}}, Organizations: map[string]struct{}{"store": {}}}
	handler := journeyHandler(repository, principal, nil)
	w := httptest.NewRecorder()
	creationRequest := request("POST", "/v1/franchise/availability", `{"organization_id":"store","resource_id":"technician","entry_type":"unavailable","reason_code":"annual-leave","starts_at":"2026-08-30T21:00:00Z","ends_at":"2026-08-30T22:00:00Z"}`)
	creationRequest.Header.Set("Idempotency-Key", "availability-test-key")
	handler.ServeHTTP(w, creationRequest)
	if w.Code != 201 || repository.availability.State != "active" || repository.availability.Version != 1 {
		t.Fatalf("create availability status=%d body=%s value=%+v", w.Code, w.Body.String(), repository.availability)
	}
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, request("GET", "/v1/franchise/availability?organization_id=store&resource_id=technician&from=2026-08-30T20:00:00Z&to=2026-08-31T20:00:00Z", ""))
	if w.Code != 200 || !strings.Contains(w.Body.String(), "working-window") {
		t.Fatalf("availability status=%d body=%s", w.Code, w.Body.String())
	}

	principal = identity.Principal{Subject: "customer-1", TenantID: "tenant", Permissions: map[string]struct{}{"customer:self": {}}, Organizations: map[string]struct{}{"store": {}}}
	handler = journeyHandler(repository, principal, nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, request("POST", "/v1/customer/appointments/appointment/cancel", `{"organization_id":"store","version":2,"reason_code":"customer-request"}`))
	if w.Code != 200 || repository.customer != "customer-1" || !strings.Contains(w.Body.String(), `"state":"cancelled"`) {
		t.Fatalf("customer cancel status=%d customer=%s body=%s", w.Code, repository.customer, w.Body.String())
	}
}

func TestLeadPermissionsConflictAndStrictJSON(t *testing.T) {
	repository := &journeyRepo{conflict: true}
	principal := identity.Principal{Subject: "sales", TenantID: "tenant", Permissions: map[string]struct{}{"lead:assign": {}}, Organizations: map[string]struct{}{"store": {}}}
	handler := journeyHandler(repository, principal, nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, request("POST", "/v1/franchise/leads/lead/assign", `{"organization_id":"store","assigned_subject":"sales","version":1}`))
	if w.Code != 409 {
		t.Fatalf("conflict status=%d body=%s", w.Code, w.Body.String())
	}
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, request("POST", "/v1/franchise/leads/lead/assign", `{"organization_id":"store","assigned_subject":"sales","version":1,"unknown":true}`))
	if w.Code != 400 {
		t.Fatalf("strict json status=%d body=%s", w.Code, w.Body.String())
	}
}

func TestAuthenticationFailure(t *testing.T) {
	handler := journeyHandler(&journeyRepo{}, identity.Principal{}, errors.New("bad token"))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, request("GET", "/v1/franchise/leads?organization_id=store", ""))
	if w.Code != 401 {
		t.Fatalf("status=%d", w.Code)
	}
}

func TestAcceptQuoteBindsIdentityAndRejectsCrossOrganization(t *testing.T) {
	repository := &journeyRepo{}
	principal := identity.Principal{Subject: "customer-1", TenantID: "tenant", Permissions: map[string]struct{}{"customer:self": {}}, Organizations: map[string]struct{}{"store": {}}}
	handler := journeyHandler(repository, principal, nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, request("POST", "/v1/customer/quotes/quote/accept", `{"organization_id":"other","version":1}`))
	if w.Code != 403 {
		t.Fatalf("cross organization status=%d body=%s", w.Code, w.Body.String())
	}
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, request("POST", "/v1/customer/quotes/quote/accept", `{"organization_id":"store","version":1}`))
	if w.Code != 200 || repository.customer != "customer-1" || !strings.Contains(w.Body.String(), `"state":"accepted"`) {
		t.Fatalf("accept status=%d customer=%s body=%s", w.Code, repository.customer, w.Body.String())
	}
}

func TestAcceptHandoverHashesExactCommandAndBindsIdentity(t *testing.T) {
	repository := &journeyRepo{}
	principal := identity.Principal{Subject: "customer-1", TenantID: "tenant", Permissions: map[string]struct{}{"customer:self": {}}, Organizations: map[string]struct{}{"store": {}}}
	handler := journeyHandler(repository, principal, nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, request("POST", "/v1/customer/handovers/handover/accept", `{"organization_id":"store","version":2,"confirmed_received":true,"serial_number":"SERIAL","checklist_id":"standard-delivery","checklist_version":1}`))
	if w.Code != 200 || repository.customer != "customer-1" || !strings.Contains(w.Body.String(), `"state":"accepted"`) {
		t.Fatalf("accept status=%d customer=%s body=%s", w.Code, repository.customer, w.Body.String())
	}
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, request("POST", "/v1/customer/handovers/handover/accept", `{"organization_id":"store","version":2,"confirmed_received":false,"serial_number":"SERIAL","checklist_id":"standard-delivery","checklist_version":1}`))
	if w.Code != 400 {
		t.Fatalf("unconfirmed handover accepted status=%d body=%s", w.Code, w.Body.String())
	}
}

func TestDeliveryChecklistPublicationAndCompletionBindOperator(t *testing.T) {
	repository := &journeyRepo{}
	principal := identity.Principal{Subject: "operator-1", TenantID: "tenant", Permissions: map[string]struct{}{"handover:manage": {}}, Organizations: map[string]struct{}{"store": {}}}
	handler := journeyHandler(repository, principal, nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, request("POST", "/v1/franchise/delivery-checklists", `{"organization_id":"store","checklist_id":"standard-delivery","version":1,"title":"Entrega estándar","items":[{"id":"serial-observed","prompt":"Verificar serie","response_type":"serial","required":true}]}`))
	if w.Code != 201 || repository.customer != "operator-1" || !strings.Contains(w.Body.String(), `"state":"published"`) {
		t.Fatalf("publish status=%d subject=%s body=%s", w.Code, repository.customer, w.Body.String())
	}
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, request("POST", "/v1/franchise/handovers/handover/complete-checklist", `{"organization_id":"store","version":1,"checklist_id":"standard-delivery","checklist_version":1,"responses":[{"item_id":"serial-observed","response_text":"SERIAL"}]}`))
	if w.Code != 200 || repository.customer != "operator-1" || !strings.Contains(w.Body.String(), `"state":"presented"`) {
		t.Fatalf("complete status=%d subject=%s body=%s", w.Code, repository.customer, w.Body.String())
	}
}

func TestDeliveryRejectionAndResolutionBindVerifiedActors(t *testing.T) {
	repository := &journeyRepo{}
	customer := identity.Principal{Subject: "customer-1", TenantID: "tenant", Permissions: map[string]struct{}{"customer:self": {}}, Organizations: map[string]struct{}{"store": {}}}
	handler := journeyHandler(repository, customer, nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, request("POST", "/v1/customer/handovers/handover/reject", `{"organization_id":"store","version":2,"reason_code":"visible-damage","details":"Rayón visible"}`))
	if w.Code != 200 || repository.customer != "customer-1" || !strings.Contains(w.Body.String(), `"state":"open"`) {
		t.Fatalf("reject status=%d subject=%s body=%s", w.Code, repository.customer, w.Body.String())
	}
	repository = &journeyRepo{}
	operator := identity.Principal{Subject: "operator-1", TenantID: "tenant", Permissions: map[string]struct{}{"handover:manage": {}}, Organizations: map[string]struct{}{"store": {}}}
	handler = journeyHandler(repository, operator, nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, request("POST", "/v1/franchise/delivery-exceptions/exception/resolve", `{"organization_id":"store","version":1,"action":"correct-and-represent","notes":"Corregir preparación"}`))
	if w.Code != 200 || repository.customer != "operator-1" || !strings.Contains(w.Body.String(), `"successor_handover"`) {
		t.Fatalf("resolve status=%d subject=%s body=%s", w.Code, repository.customer, w.Body.String())
	}
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, request("GET", "/v1/franchise/delivery-exceptions?organization_id=store", ""))
	if w.Code != 200 || !strings.Contains(w.Body.String(), `"exception"`) {
		t.Fatalf("list status=%d body=%s", w.Code, w.Body.String())
	}
}

func TestReturnReceiptAndDispositionBindOperatorAndClosedFields(t *testing.T) {
	repository := &journeyRepo{}
	operator := identity.Principal{Subject: "operator-1", TenantID: "tenant", Permissions: map[string]struct{}{"handover:manage": {}}, Organizations: map[string]struct{}{"store": {}}}
	handler := journeyHandler(repository, operator, nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, request("GET", "/v1/franchise/returns?organization_id=store", ""))
	if w.Code != 200 || !strings.Contains(w.Body.String(), `"authorization"`) {
		t.Fatalf("return list status=%d body=%s", w.Code, w.Body.String())
	}
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, request("POST", "/v1/franchise/return-authorizations/authorization/receive", `{"organization_id":"store","serial_number":"SERIAL","condition_code":"damaged","notes":"Daño confirmado","evidence_sha256":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}`))
	if w.Code != 200 || repository.customer != "operator-1" || !strings.Contains(w.Body.String(), `"condition_code":"damaged"`) {
		t.Fatalf("return receive status=%d subject=%s body=%s", w.Code, repository.customer, w.Body.String())
	}
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, request("POST", "/v1/franchise/return-receipts/receipt/decide", `{"organization_id":"store","inventory_action":"quarantine","notes":"Separar hasta reconciliar"}`))
	if w.Code != 200 || repository.customer != "operator-1" || !strings.Contains(w.Body.String(), `"customer_remedy":"refund"`) {
		t.Fatalf("return decision status=%d subject=%s body=%s", w.Code, repository.customer, w.Body.String())
	}
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, request("POST", "/v1/franchise/return-receipts/receipt/decide", `{"organization_id":"store","inventory_action":"restock","notes":"No","refund_now":true}`))
	if w.Code != 400 {
		t.Fatalf("return decision accepted browser-owned effect status=%d body=%s", w.Code, w.Body.String())
	}
}

// TestQuoteAcceptanceBrowserPostgres is an opt-in AUTHORED connected fixture.
// Reuses the real journey API, PostgreSQL repository and RS256/JWKS issuer.
// No live login, real customer, provider, payment or business acceptance claim.
func TestQuoteAcceptanceBrowserPostgres(t *testing.T) {
	if os.Getenv("ELITE_QUOTE_E2E") != "1" {
		t.Skip("explicit disposable quote gate not requested")
	}
	for _, project := range []string{"chromium-desktop", "chromium-mobile", "firefox-desktop", "webkit-desktop"} {
		if selected := os.Getenv("ELITE_QUOTE_PROJECT"); selected != "" && selected != project {
			continue
		}
		t.Run(project, func(t *testing.T) { runQuoteAcceptanceBrowser(t, project) })
	}
}

func runQuoteAcceptanceBrowser(t *testing.T, project string) {
	budget := 3 * time.Minute
	if os.Getenv("ELITE_PAYMENT_E2E") == "1" {
		budget = 5 * time.Minute
	}
	ctx, cancel := context.WithTimeout(context.Background(), budget)
	defer cancel()
	config, err := pgxpool.ParseConfig(os.Getenv("TEST_DATABASE_URL"))
	if err != nil || config.ConnConfig.Host != "127.0.0.1" || !strings.HasPrefix(config.ConnConfig.Database, "elite_confirmation_") {
		t.Fatal("requires disposable loopback elite_confirmation_* database")
	}
	if os.Getenv("ELITE_RETURN_MULTITAB_E2E") == "1" && (os.Getenv("ELITE_RETURN_RECOVERY_E2E") != "1" || os.Getenv("ELITE_DELIVERY_READ_E2E") != "1" || os.Getenv("ELITE_ORDER_E2E") != "1") {
		t.Fatal("return multi-tab gate requires return, order and delivery-read fixtures")
	}
	if os.Getenv("ELITE_RETURN_SESSION_E2E") == "1" && os.Getenv("ELITE_RETURN_MULTITAB_E2E") != "1" {
		t.Fatal("return session gate requires multi-tab fixture")
	}
	web := os.Getenv("ELITE_WEB_ROOT")
	if !filepath.IsAbs(web) {
		t.Fatal("absolute materialized web root required")
	}
	artifacts, err := os.MkdirTemp(web, "quote-browser-artifacts-")
	if err != nil {
		t.Fatal(err)
	}
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		t.Fatal(err)
	}
	tenant := (randomid.Generator{}).New()
	fixtures := []string{
		"insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1::uuid,'quote-'||replace(($1::uuid)::text,'-',''),'Synthetic quote','Synthetic quote')",
		"insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'store','store','Synthetic Store','store')",
		"insert into crm.customer_profile(tenant_id,customer_principal_id,display_name,email_normalized)values($1,'customer','Synthetic Customer','quote@example.invalid')",
		"insert into crm.lead(tenant_id,lead_id,organization_id,customer_principal_id,lifecycle_state,source_code,contact_payload)values($1,'lead','store','customer','new','fixture','{}')",
		"insert into catalog.vehicle_model(tenant_id,model_id,model_code,display_name,vehicle_class,lifecycle_state)values($1,'model','model','Synthetic bicycle','bicycle','active')",
		"insert into catalog.vehicle_variant(tenant_id,variant_id,model_id,variant_code,display_name,battery_specification,lifecycle_state)values($1,'variant','model','variant','Synthetic variant','{}','active')",
		"insert into pricing.price_book(tenant_id,price_book_id,market,currency,valid_from,status)values($1,'book','AR','ARS',clock_timestamp()-interval '1 day','active')",
		"insert into pricing.price_book_entry(tenant_id,price_book_id,variant_id,amount_minor_units,tax_mode)values($1,'book','variant',250000,'inclusive')",
		"insert into sales.quotation(tenant_id,quotation_id,organization_id,lead_id,customer_principal_id,variant_id,price_book_id,currency,total_minor_units,valid_until,state,version)select $1,q,'store','lead','customer','variant','book','ARS',250000,clock_timestamp()+interval '1 day','issued',1 from unnest(array['quote-success','quote-lost','quote-race','quote-denied'])q",
		"insert into sales.quotation(tenant_id,quotation_id,organization_id,lead_id,customer_principal_id,variant_id,price_book_id,currency,total_minor_units,created_at,valid_until,state,version)values($1,'quote-expired','store','lead','customer','variant','book','ARS',250000,clock_timestamp()-interval '2 days',clock_timestamp()-interval '1 day','issued',1)",
	}
	for _, sql := range fixtures {
		if _, err = pool.Exec(ctx, sql, tenant); err != nil {
			pool.Close()
			t.Fatal(err)
		}
	}
	pool.Close()
	verifier, token := confirmationTestIssuer(t)
	phases := []string{"initial", "read-failure", "restart"}
	if os.Getenv("ELITE_ORDER_E2E") == "1" {
		phases = append(phases, "operation", "operation-read-failure", "operation-restart")
	}
	if os.Getenv("ELITE_PAYMENT_E2E") == "1" {
		phases = append(phases, "payment", "payment-disabled", "payment-restart")
	}
	if os.Getenv("ELITE_DELIVERY_READ_E2E") == "1" {
		if os.Getenv("ELITE_ORDER_E2E") != "1" {
			t.Fatal("delivery read gate requires connected order/stock gate")
		}
		phases = append(phases, "delivery-read-corrupt", "delivery-read-restored")
	}
	if os.Getenv("ELITE_DELIVERY_ACTION_E2E") == "1" {
		if os.Getenv("ELITE_DELIVERY_READ_E2E") != "1" {
			t.Fatal("delivery action gate requires delivery read fixture")
		}
		phases = append(phases, "delivery-action-loss")
	}
	if os.Getenv("ELITE_OPERATOR_SECTIONS_E2E") == "1" {
		if os.Getenv("ELITE_ORDER_E2E") != "1" {
			t.Fatal("operator sections require connected order gate")
		}
		phases = append(phases, "operation-sections")
	}
	if os.Getenv("ELITE_RESOLUTION_RECOVERY_E2E") == "1" {
		if os.Getenv("ELITE_DELIVERY_ACTION_E2E") != "1" {
			t.Fatal("resolution recovery requires actual rejected delivery fixture")
		}
		phases = append(phases, "operation-resolution")
	}
	if os.Getenv("ELITE_LEAD_RECOVERY_E2E") == "1" {
		phases = append(phases, "operation-lead")
	}
	if os.Getenv("ELITE_AVAILABILITY_RECOVERY_E2E") == "1" {
		phases = append(phases, "operation-availability")
	}
	if os.Getenv("ELITE_QUOTE_CREATE_E2E") == "1" {
		phases = append(phases, "operation-quote")
	}
	if os.Getenv("ELITE_RETURN_RECOVERY_E2E") == "1" {
		phases = append(phases, "operation-returns")
	}
	if os.Getenv("ELITE_CHECKLIST_COMPLETE_E2E") == "1" {
		phases = append(phases, "operation-complete-checklist")
	}
	if os.Getenv("ELITE_CHECKLIST_PUBLISH_E2E") == "1" {
		phases = append(phases, "operation-publish-checklist")
	}
	if os.Getenv("ELITE_SLOT_CREATE_E2E") == "1" {
		phases = append(phases, "operation-create-slot")
	}
	if os.Getenv("ELITE_RESOURCE_CREATE_E2E") == "1" {
		phases = append(phases, "operation-create-resource")
	}
	if os.Getenv("ELITE_AVAILABILITY_CREATE_E2E") == "1" {
		phases = append(phases, "operation-create-availability")
	}
	for _, phase := range phases {
		func() {
			// New pool, repository, API process fixture and Next process each phase.
			p, err := pgxpool.NewWithConfig(ctx, config.Copy())
			if err != nil {
				t.Fatal(err)
			}
			defer p.Close()
			if phase == "operation-quote" {
				for i := 0; i < 3; i++ {
					if _, err := p.Exec(ctx, `insert into crm.lead(tenant_id,lead_id,organization_id,customer_principal_id,lifecycle_state,source_code,contact_payload)values($1,$2,'store','customer','new','fixture','{}')`, tenant, fmt.Sprintf("quote-create-%d", i)); err != nil {
						t.Fatal(err)
					}
				}
			}
			if phase == "operation-availability" {
				repo := postgres.NewFranchiseJourney(p)
				if _, err := p.Exec(ctx, `insert into org.public_location(tenant_id,organization_id,city,region,country,published)values($1,'store','Synthetic','Synthetic','AR',true)`, tenant); err != nil {
					t.Fatal(err)
				}
				for i := 0; i < 5; i++ {
					start := time.Now().UTC().Add(time.Duration(72+i*24) * time.Hour).Truncate(time.Second)
					entryType, reason := "working", ""
					if i%2 == 1 {
						entryType, reason = "unavailable", "synthetic-absence"
					}
					if _, err := repo.CreateAvailability(ctx, tenant, "fixture", franchisejourney.AvailabilityEntry{ID: fmt.Sprintf("cancel-window-%d", i), OrganizationID: "store", EntryType: entryType, ReasonCode: reason, StartsAt: start, EndsAt: start.Add(time.Hour)}, randomid.Generator{}.New()); err != nil {
						t.Fatal(err)
					}
					if i == 4 {
						if _, err := p.Exec(ctx, `insert into crm.appointment(tenant_id,appointment_id,organization_id,lead_id,appointment_kind,starts_at,state,version)values($1,'blocked-window-appointment','store','lead','consultation',$2,'requested',1)`, tenant, start.Add(15*time.Minute)); err != nil {
							t.Fatal(err)
						}
					}
					if i%2 == 0 {
						if _, err := repo.CreateAppointmentSlot(ctx, tenant, franchisejourney.AppointmentSlot{ID: fmt.Sprintf("window-slot-%d", i), OrganizationID: "store", Kind: "consultation", StartsAt: start, EndsAt: start.Add(time.Hour), Capacity: 2}, randomid.Generator{}.New()); err != nil {
							t.Fatal(err)
						}
					}
				}
			}
			if phase == "operation-lead" {
				if _, err := p.Exec(ctx, `insert into crm.lead(tenant_id,lead_id,organization_id,customer_principal_id,lifecycle_state,source_code,contact_payload) values($1,'operator-lead','store','customer','new','fixture','{}')`, tenant); err != nil {
					t.Fatal(err)
				}
			}
			if phase == "operation" {
				if _, err := p.Exec(ctx, `insert into inventory.stock_unit(tenant_id,stock_unit_id,organization_id,variant_id,serial_number,vin,battery_serial_number,state,version,received_at) select $1,s,'store','variant',s,s,s,'available',1,clock_timestamp() from unnest(array['stock-lost','stock-race'])s`, tenant); err != nil {
					t.Fatal(err)
				}
			}
			deliverySerial := ""
			mux := http.NewServeMux()
			if phase == "delivery-read-corrupt" {
				for _, q := range []string{
					`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,'foreign-delivery','foreign-delivery','Synthetic foreign store','store')`,
					`insert into inventory.stock_unit(tenant_id,stock_unit_id,organization_id,variant_id,serial_number,vin,battery_serial_number,state,version,received_at) values($1,'PRIVATE-FOREIGN-STOCK','foreign-delivery','variant','PRIVATE-FOREIGN-STOCK','PRIVATE-FOREIGN-STOCK','PRIVATE-FOREIGN-STOCK','available',1,clock_timestamp())`,
					`insert into sales.delivery_handover(tenant_id,handover_id,organization_id,order_id,customer_principal_id,stock_unit_id,state,version) select $1,'read-fixture-handover','store',o.order_id,o.customer_principal_id,'PRIVATE-FOREIGN-STOCK','prepared',1 from sales.customer_order o join sales.customer_order_line l on l.tenant_id=o.tenant_id and l.order_id=o.order_id where o.tenant_id=$1 and o.organization_id='store' and l.allocated_stock_unit_id is not null order by o.order_id limit 1`,
				} {
					if tag, err := p.Exec(ctx, q, tenant); err != nil || tag.RowsAffected() != 1 {
						t.Fatalf("delivery read fixture failed: %v", err)
					}
				}
			}
			if phase == "delivery-read-restored" {
				if tag, err := p.Exec(ctx, `update sales.delivery_handover h set stock_unit_id=l.allocated_stock_unit_id from sales.customer_order_line l where h.tenant_id=$1 and h.handover_id='read-fixture-handover' and l.tenant_id=h.tenant_id and l.order_id=h.order_id and l.allocated_stock_unit_id is not null`, tenant); err != nil || tag.RowsAffected() != 1 {
					t.Fatalf("restore synthetic relation: %v", err)
				}
			}
			if phase == "delivery-action-loss" {
				if err := p.QueryRow(ctx, `select i.serial_number from sales.delivery_handover h join inventory.stock_unit i on i.tenant_id=h.tenant_id and i.stock_unit_id=h.stock_unit_id where h.tenant_id=$1 and h.handover_id='read-fixture-handover'`, tenant).Scan(&deliverySerial); err != nil {
					t.Fatal(err)
				}
				if _, err := p.Exec(ctx, `insert into sales.delivery_handover(tenant_id,handover_id,organization_id,order_id,customer_principal_id,stock_unit_id,state,version) select tenant_id,'reject-fixture-handover',organization_id,order_id,customer_principal_id,stock_unit_id,'prepared',1 from sales.delivery_handover where tenant_id=$1 and handover_id='read-fixture-handover'`, tenant); err != nil {
					t.Fatal(err)
				}
				repo := postgres.NewFranchiseJourney(p)
				checklist, err := repo.PublishDeliveryChecklist(ctx, tenant, "fixture-operator", franchisejourney.DeliveryChecklist{ID: "delivery-ui-test", OrganizationID: "store", Version: 1, Title: "Synthetic preparation", State: "published", Items: []franchisejourney.ChecklistItem{{ID: "serial", Ordinal: 1, Prompt: "Verify serial", ResponseType: "serial", Required: true}}}, randomid.Generator{}.New())
				if err != nil {
					t.Fatal(err)
				}
				for _, id := range []string{"read-fixture-handover", "reject-fixture-handover"} {
					if _, err := repo.CompleteDeliveryChecklist(ctx, tenant, "store", "fixture-operator", id, 1, checklist.ID, 1, []franchisejourney.ChecklistResponse{{ItemID: "serial", ResponseText: deliverySerial}}, randomid.Generator{}.New()); err != nil {
						t.Fatal(err)
					}
				}
			}
			if phase == "operation-resolution" {
				repo := postgres.NewFranchiseJourney(p)
				for _, id := range []string{"return-fixture-handover", "exchange-fixture-handover"} {
					if _, err := p.Exec(ctx, `insert into sales.delivery_handover(tenant_id,handover_id,organization_id,order_id,customer_principal_id,stock_unit_id,state,version) select tenant_id,$2,organization_id,order_id,customer_principal_id,stock_unit_id,'prepared',1 from sales.delivery_handover where tenant_id=$1 and handover_id='reject-fixture-handover'`, tenant, id); err != nil {
						t.Fatal(err)
					}
					var serial string
					if err := p.QueryRow(ctx, `select i.serial_number from sales.delivery_handover h join inventory.stock_unit i on i.tenant_id=h.tenant_id and i.stock_unit_id=h.stock_unit_id where h.tenant_id=$1 and h.handover_id=$2`, tenant, id).Scan(&serial); err != nil {
						t.Fatal(err)
					}
					if _, err := repo.CompleteDeliveryChecklist(ctx, tenant, "store", "fixture-operator", id, 1, "delivery-ui-test", 1, []franchisejourney.ChecklistResponse{{ItemID: "serial", ResponseText: serial}}, randomid.Generator{}.New()); err != nil {
						t.Fatal(err)
					}
					if _, err := repo.RejectHandover(ctx, tenant, "store", "customer", id, 2, "asset-condition", "Synthetic fixture rejection", strings.Repeat("a", 64), randomid.Generator{}.New(), randomid.Generator{}.New()); err != nil {
						t.Fatal(err)
					}
				}
			}
			FranchiseJourneyModule{Service: franchisejourney.NewService(postgres.NewFranchiseJourney(p), randomid.Generator{}, browserAppointmentClock{})}.Register(mux, verifier)
			paymentProvider := ""
			if strings.HasPrefix(phase, "payment") && phase != "payment-disabled" {
				paymentProvider = "stripe"
			}
			CommerceModule{Service: commerce.NewService(postgres.NewCommerce(p), randomid.Generator{}), PaymentProvider: paymentProvider}.Register(mux, verifier)
			var returnWritten atomic.Bool
			var returnReadFailures atomic.Int32
			var exceptionReadFailures atomic.Int32
			var resolutionWritten atomic.Bool
			var leadWritten atomic.Bool
			var availabilityWritten atomic.Bool
			var availabilityReadFailures atomic.Int32
			var leadReadFailures atomic.Int32
			var resolutionReadFailures atomic.Int32
			api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if phase == "operation-returns" && r.Method == http.MethodPost && r.URL.Path == "/v1/franchise/return-authorizations/ui-return-loss/receive" {
					defer returnWritten.Store(true)
				}
				if phase == "operation-returns" && r.Method == http.MethodGet && r.URL.Path == "/v1/franchise/returns" && returnWritten.Load() && returnReadFailures.Add(1) == 1 {
					writeProblem(w, 503, "INJECTED_RETURN_LIST_FAILURE", "PRIVATE-SYNTHETIC-RETURN-DETAIL")
					return
				}

				if phase == "operation-availability" && r.URL.Path == "/v1/franchise/availability" && r.Method == http.MethodGet && availabilityWritten.Load() && availabilityReadFailures.Add(1) == 1 {
					writeProblem(w, 503, "INJECTED_AVAILABILITY_QUERY_FAILURE", "PRIVATE-SYNTHETIC-AVAILABILITY-DETAIL")
					return
				}
				if phase == "operation-availability" && r.Method == http.MethodPost && strings.Contains(r.URL.Path, "/availability/cancel-window-") {
					defer availabilityWritten.Store(true)
				}
				if phase == "operation-lead" && r.URL.Path == "/v1/franchise/leads" && leadWritten.Load() && leadReadFailures.Add(1) == 1 {
					writeProblem(w, 503, "INJECTED_LEAD_QUERY_FAILURE", "PRIVATE-SYNTHETIC-LEAD-DETAIL")
					return
				}
				if phase == "operation-lead" && r.Method == http.MethodPost && strings.Contains(r.URL.Path, "/leads/operator-lead/") {
					defer leadWritten.Store(true)
				}
				if phase == "operation-resolution" && r.URL.Path == "/v1/franchise/delivery-exceptions" && resolutionWritten.Load() && resolutionReadFailures.Add(1) == 1 {
					writeProblem(w, 503, "INJECTED_RESOLUTION_QUERY_FAILURE", "PRIVATE-SYNTHETIC-RESOLUTION-DETAIL")
					return
				}
				if phase == "operation-resolution" && r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/resolve") {
					defer resolutionWritten.Store(true)
				}
				if phase == "operation-sections" && r.URL.Path == "/v1/franchise/delivery-exceptions" && exceptionReadFailures.Add(1) == 1 {
					writeProblem(w, 503, "INJECTED_SECTION_UNAVAILABLE", "PRIVATE-SYNTHETIC-SECTION-DETAIL")
					return
				}
				if phase == "operation-read-failure" && r.URL.Path == "/v1/commerce/orders" {
					writeProblem(w, 503, "INJECTED_READ_UNAVAILABLE", "synthetic outage")
					return
				}
				if phase == "read-failure" && r.URL.Path == "/v1/customer/journey" {
					w.Header().Set("Content-Type", "application/problem+json")
					w.WriteHeader(http.StatusServiceUnavailable)
					_, _ = w.Write([]byte(`{"code":"INJECTED_READ_UNAVAILABLE"}`))
					return
				}
				mux.ServeHTTP(w, r)
			}))
			defer api.Close()
			listener, err := net.Listen("tcp", "127.0.0.1:0")
			if err != nil {
				t.Fatal(err)
			}
			address := listener.Addr().String()
			listener.Close()
			target, _ := url.Parse("http://" + address)
			edge := httptest.NewTLSServer(httputil.NewSingleHostReverseProxy(target))
			defer edge.Close()
			env := []string{}
			for _, item := range os.Environ() {
				name := strings.ToUpper(strings.SplitN(item, "=", 2)[0])
				if name == "TEST_DATABASE_URL" || name == "DATABASE_URL" || strings.HasPrefix(name, "ELITE_") || name == "ENTERPRISE_API_BASE_URL" || name == "AUTH_SESSION_SECRET" || name == "APP_BASE_URL" {
					continue
				}
				env = append(env, item)
			}

			if phase == "operation-returns" {
				if e := p.QueryRow(ctx, `select i.serial_number from sales.delivery_handover h join inventory.stock_unit i on i.tenant_id=h.tenant_id and i.stock_unit_id=h.stock_unit_id where h.tenant_id=$1 and h.handover_id='read-fixture-handover'`, tenant).Scan(&deliverySerial); e != nil {
					t.Fatal(e)
				}
				repo := postgres.NewFranchiseJourney(p)
				checklist := franchisejourney.DeliveryChecklist{ID: "ui-return-checklist", OrganizationID: "store", Version: 1, Title: "Synthetic returns", Items: []franchisejourney.ChecklistItem{{ID: "serial", Ordinal: 1, Prompt: "Serie", ResponseType: "serial", Required: true}}}
				if _, e := repo.PublishDeliveryChecklist(ctx, tenant, "fixture-author", checklist, randomid.Generator{}.New()); e != nil {
					t.Fatal(e)
				}
				for i, id := range []string{"ui-return-loss", "ui-exchange-invalid", "ui-return-race"} {
					if _, e := p.Exec(ctx, `insert into sales.delivery_handover(tenant_id,handover_id,organization_id,order_id,customer_principal_id,stock_unit_id,state,version) select tenant_id,$2,organization_id,order_id,customer_principal_id,stock_unit_id,'prepared',1 from sales.delivery_handover where tenant_id=$1 and handover_id='read-fixture-handover'`, tenant, id); e != nil {
						t.Fatal(e)
					}
					if _, e := repo.CompleteDeliveryChecklist(ctx, tenant, "store", "fixture", id, 1, checklist.ID, 1, []franchisejourney.ChecklistResponse{{ItemID: "serial", ResponseText: deliverySerial}}, randomid.Generator{}.New()); e != nil {
						t.Fatal(e)
					}
					x, e := repo.RejectHandover(ctx, tenant, "store", "customer", id, 2, "serial-mismatch", "Synthetic", strings.Repeat("a", 64), randomid.Generator{}.New(), randomid.Generator{}.New())
					if e != nil {
						t.Fatal(e)
					}
					action := "return"
					if i == 1 {
						action = "exchange"
					}
					if _, e = repo.ResolveDeliveryException(ctx, tenant, "store", "fixture", x.ID, 1, action, "Synthetic", "", id, randomid.Generator{}.New(), randomid.Generator{}.New()); e != nil {
						t.Fatal(e)
					}
				}
			}
			if phase == "operation-complete-checklist" {
				if e := p.QueryRow(ctx, `select i.serial_number from sales.delivery_handover h join inventory.stock_unit i on i.tenant_id=h.tenant_id and i.stock_unit_id=h.stock_unit_id and i.organization_id=h.organization_id where h.tenant_id=$1 and h.handover_id='read-fixture-handover'`, tenant).Scan(&deliverySerial); e != nil {
					t.Fatal(e)
				}
				repo := postgres.NewFranchiseJourney(p)
				if _, e := repo.PublishDeliveryChecklist(ctx, tenant, "fixture-author", franchisejourney.DeliveryChecklist{ID: "ui-completion-checklist", OrganizationID: "store", Version: 1, Title: "Synthetic completion", Items: []franchisejourney.ChecklistItem{{ID: "serial", Ordinal: 1, Prompt: "Serie", ResponseType: "serial", Required: true}, {ID: "confirmed", Ordinal: 2, Prompt: "Confirmación", ResponseType: "confirmation", Required: true}}}, randomid.Generator{}.New()); e != nil {
					t.Fatal(e)
				}
				for _, id := range []string{"ui-completion-loss", "ui-completion-invalid", "ui-completion-race"} {
					if _, e := p.Exec(ctx, `insert into sales.delivery_handover(tenant_id,handover_id,organization_id,order_id,customer_principal_id,stock_unit_id,state,version) select tenant_id,$2,organization_id,order_id,customer_principal_id,stock_unit_id,'prepared',1 from sales.delivery_handover where tenant_id=$1 and handover_id='read-fixture-handover'`, tenant, id); e != nil {
						t.Fatal(e)
					}
				}
			}
			if phase == "operation-create-slot" {
				for _, days := range []int{30, 31, 200} {
					start := time.Now().UTC().Add(time.Duration(days) * 24 * time.Hour)
					if _, e := postgres.NewFranchiseJourney(p).CreateAvailability(ctx, tenant, "planner", franchisejourney.AvailabilityEntry{ID: fmt.Sprintf("slot-create-working-%d", days), OrganizationID: "store", EntryType: "working", StartsAt: start.Add(-time.Hour), EndsAt: start.Add(2 * time.Hour)}, randomid.Generator{}.New()); e != nil {
						t.Fatal(e)
					}
				}
			}
			env = append(env, "ELITE_RETURN_SESSION_E2E="+os.Getenv("ELITE_RETURN_SESSION_E2E"), "ELITE_RETURN_MULTITAB_E2E="+os.Getenv("ELITE_RETURN_MULTITAB_E2E"), "ELITE_DELIVERY_SERIAL="+deliverySerial, "ELITE_QUOTE_E2E=1", "ELITE_QUOTE_PHASE="+phase, "ELITE_QUOTE_TENANT="+tenant,
				"ELITE_QUOTE_WRITER="+token("writer", tenant, []string{"lead:read", "quote:write"}, []string{"store"}),
				"ELITE_QUOTE_HANDOVER="+token("handover", tenant, []string{"handover:manage"}, []string{"store"}),
				"ELITE_QUOTE_HANDOVER_OTHER="+token("handover_other", tenant, []string{"handover:manage"}, []string{"store"}),
				"ELITE_QUOTE_LEAD="+token("lead-operator", tenant, []string{"lead:read", "lead:assign", "lead:update"}, []string{"store"}),
				"ELITE_QUOTE_SCHEDULER="+token("scheduler", tenant, []string{"availability:read", "availability:manage", "inventory:allocate"}, []string{"store"}),
				"ELITE_QUOTE_PLANNER="+token("planner", tenant, []string{"appointment:manage"}, []string{"store"}),
				"ELITE_QUOTE_RESOURCE="+token("resource", tenant, []string{"resource:manage"}, []string{"store"}),
				"ELITE_QUOTE_AVAILABILITY="+token("availability", tenant, []string{"availability:read"}, []string{"store"}),
				"ELITE_QUOTE_PAYMENT="+token("payment", tenant, []string{"payment:create"}, []string{"store"}),
				"ELITE_QUOTE_OPERATOR="+token("operator", tenant, []string{"inventory:allocate"}, []string{"store"}),
				"ELITE_QUOTE_READER="+token("reader", tenant, []string{"admin:read"}, []string{"store"}),
				"ELITE_QUOTE_CUSTOMER="+token("customer", tenant, []string{"customer:self"}, []string{"store"}),
				"ELITE_QUOTE_STRANGER="+token("stranger", tenant, []string{"customer:self"}, []string{"store"}),
				"ELITE_QUOTE_OBSERVER="+token("customer", tenant, []string{"lead:read"}, []string{"store"}),
				"ELITE_QUOTE_OTHER_ORG="+token("customer", tenant, []string{"customer:self"}, []string{"other"}),
				"ELITE_QUOTE_OTHER_TENANT="+token("customer", (randomid.Generator{}).New(), []string{"customer:self"}, []string{"store"}),
				"ENTERPRISE_API_BASE_URL="+api.URL, "ELITE_WEB_ROOT="+web,
				"AUTH_SESSION_SECRET=synthetic-"+(randomid.Generator{}).New()+(randomid.Generator{}).New(),
				"APP_BASE_URL="+edge.URL, "ELITE_BASE_URL="+edge.URL)
			server := exec.CommandContext(ctx, "node", filepath.Join(web, "node_modules", "next", "dist", "bin", "next"), "start", "--hostname", "127.0.0.1", "--port", strings.Split(address, ":")[1])
			server.Dir, server.Env = web, env
			if err = server.Start(); err != nil {
				t.Fatal(err)
			}
			defer func() { _ = server.Process.Kill(); _ = server.Wait() }()
			ready := false
			client := &http.Client{Timeout: time.Second}
			for deadline := time.Now().Add(20 * time.Second); time.Now().Before(deadline); {
				res, err := client.Get(target.String() + "/icon.svg")
				if err == nil {
					res.Body.Close()
					if res.StatusCode == 200 {
						ready = true
						break
					}
				}
				select {
				case <-ctx.Done():
					t.Fatal("startup cancelled")
				case <-time.After(100 * time.Millisecond):
				}
			}
			if !ready {
				t.Fatal("web did not start")
			}
			gate := filepath.Join(web, "microsoft_playwright_browser_gate")
			pattern := "customer quote acceptance connects"
			if strings.HasPrefix(phase, "operation") {
				pattern = "operator stock continues existing orders"
			}
			if strings.HasPrefix(phase, "payment") {
				pattern = "operator payment request connects"
			}
			cmd := exec.CommandContext(ctx, "node", filepath.Join(gate, "node_modules", "@playwright", "test", "cli.js"), "test", "tests/enterprise-web.spec.mjs", "-g", pattern, "--project="+project, "--output="+filepath.Join(artifacts, project, phase))
			cmd.Dir, cmd.Env = gate, env
			output, err := cmd.CombinedOutput()
			clean := string(output)
			for _, item := range env {
				parts := strings.SplitN(item, "=", 2)
				if len(parts) == 2 && (parts[0] == "AUTH_SESSION_SECRET" || (strings.HasPrefix(parts[0], "ELITE_QUOTE_") && strings.Count(parts[1], ".") == 2)) {
					clean = strings.ReplaceAll(clean, parts[1], "[REDACTED_FIXTURE]")
				}
			}
			t.Log(clean)
			if err != nil {
				t.Fatalf("quote browser phase %s failed: %v", phase, err)
			}
			if phase == "operation-returns" && os.Getenv("ELITE_RETURN_SESSION_E2E") == "1" && !strings.Contains(clean, "RETURN_SESSION_BOUNDARY_PASS transitions=8 denied=12 authorized_reads=2 restored=2") {
				t.Fatal("return session gate did not produce required evidence")
			}
			if phase == "operation-returns" && os.Getenv("ELITE_RETURN_MULTITAB_E2E") == "1" && !strings.Contains(clean, "RETURN_MULTITAB_BROWSER_PASS races=6 posts=12 wins=6 conflicts=6 losses=2 opener_copy=1") {
				t.Fatal("multi-tab browser proof missing; a skipped test is not a PASS")
			}
			checks := []string{
				"select count(*) from sales.quotation where tenant_id=$1 and state='accepted' and order_id is not null",
				"select count(*) from sales.customer_order where tenant_id=$1 and state='placed' and total_minor_units=250000 and customer_principal_id='customer'",
				"select count(*) from sales.customer_order_line where tenant_id=$1 and quantity=1 and unit_price_minor_units=250000",
				"select count(*) from sales.quotation_acceptance where tenant_id=$1 and customer_principal_id='customer' and evidence_sha256_hex ~ '^[a-f0-9]{64}$'",
				"select count(*) from platform.outbox_event where tenant_id=$1 and event_type='quotation.accepted'",
				"select count(*) from platform.outbox_event where tenant_id=$1 and event_type='customer-order.placed'",
			}
			for _, sql := range checks {
				var n int
				if err := p.QueryRow(ctx, sql, tenant).Scan(&n); err != nil {
					t.Fatal(err)
				}
				if n != 3 {
					t.Fatalf("durable invariant: expected3 got%d", n)
				}
			}
			var denied int
			if err := p.QueryRow(ctx, "select count(*) from sales.quotation where tenant_id=$1 and quotation_id in ('quote-denied','quote-expired') and state='issued' and order_id is null", tenant).Scan(&denied); err != nil || denied != 2 {
				t.Fatal("negative requests changed quotes")
			}
			t.Log("QUOTE_ACCEPTANCE_BROWSER_POSTGRES_PASS phase=" + phase + " orders=3 acceptances=3 outbox=6 denied_unchanged=2")

			if phase == "operation-complete-checklist" {
				var rows, answers, events, actors, accepted, rejected int
				e := p.QueryRow(ctx, `select count(*),
 (select count(*) from sales.delivery_checklist_response where tenant_id=$1 and handover_id like 'ui-completion-%'),
 (select count(*) from platform.outbox_event where tenant_id=$1 and aggregate_id like 'ui-completion-%' and event_type='delivery-handover.checklist-completed'),
 (select count(*) from platform.outbox_event where tenant_id=$1 and aggregate_id like 'ui-completion-%' and event_type='delivery-handover.checklist-completed' and payload->>'actor_subject'='handover'),
 count(*)filter(where state='accepted'),count(*)filter(where state='rejected')
 from sales.delivery_handover where tenant_id=$1 and handover_id like 'ui-completion-%' and checklist_completed_at is not null and checklist_completed_by_subject='handover'`, tenant).Scan(&rows, &answers, &events, &actors, &accepted, &rejected)
				if e != nil || rows != 3 || answers != 6 || events != 3 || actors != 3 || accepted != 1 || rejected != 1 {
					t.Fatalf("completion counts=%d/%d/%d/%d/%d/%d error=%v", rows, answers, events, actors, accepted, rejected, e)
				}
				t.Log("CHECKLIST_COMPLETION_BROWSER_PASS handovers=3 answers=6 events=3 actors=3 accepted=1 rejected=1 presented=1")
			}
			if phase == "operation-returns" {
				var receipts, decisions, requests, events int
				e := p.QueryRow(ctx, `select (select count(*) from sales.return_receipt where tenant_id=$1 and authorization_id in ('ui-return-loss','ui-exchange-invalid','ui-return-race') and received_by_subject='handover'),(select count(*) from sales.return_disposition d join sales.return_receipt r on r.tenant_id=d.tenant_id and r.receipt_id=d.receipt_id where d.tenant_id=$1 and r.authorization_id in ('ui-return-loss','ui-exchange-invalid','ui-return-race') and d.decided_by_subject='handover'),(select count(*) from sales.return_effect_request e join sales.return_disposition d on d.tenant_id=e.tenant_id and d.disposition_id=e.disposition_id join sales.return_receipt r on r.tenant_id=d.tenant_id and r.receipt_id=d.receipt_id where e.tenant_id=$1 and r.authorization_id in ('ui-return-loss','ui-exchange-invalid','ui-return-race') and e.state='requested'),(select count(*) from platform.outbox_event where tenant_id=$1 and event_type in ('return.received','return.disposition-requested') and payload->>'actor_subject'='handover')`, tenant).Scan(&receipts, &decisions, &requests, &events)
				if e != nil || receipts != 3 || decisions != 3 || requests != 11 || events != 6 {
					t.Fatal("return browser effects", receipts, decisions, requests, events, e)
				}
				if os.Getenv("ELITE_RETURN_MULTITAB_E2E") == "1" {
					t.Log("RETURN_OPERATIONS_BROWSER_PASS receipts=3 decisions=3 requests=11 events=6 actors=6 multitab=PASS")
				} else {
					t.Log("RETURN_OPERATIONS_BROWSER_PASS receipts=3 decisions=3 requests=11 events=6 actors=6 replay=409 exact_recovery=PASS list_failure=PASS")
				}
			}

			if phase == "operation-publish-checklist" {
				var versions, items, events, actors int
				e := p.QueryRow(ctx, `select count(*),
 (select count(*) from sales.delivery_checklist_item where tenant_id=$1 and checklist_id like 'ui-publication-%'),
 (select count(*) from platform.outbox_event where tenant_id=$1 and aggregate_id like 'ui-publication-%' and event_type='delivery-checklist.published'),
 (select count(*) from platform.outbox_event where tenant_id=$1 and aggregate_id like 'ui-publication-%' and event_type='delivery-checklist.published' and payload->>'actor_subject'='handover')
 from sales.delivery_checklist_template where tenant_id=$1 and checklist_id like 'ui-publication-%' and state='published' and created_by_subject='handover'`, tenant).Scan(&versions, &items, &events, &actors)
				if e != nil || versions != 3 || items != 6 || events != 3 || actors != 3 {
					t.Fatalf("checklist publication counts=%d/%d/%d/%d err=%v", versions, items, events, actors, e)
				}
				t.Log("CHECKLIST_PUBLICATION_BROWSER_PASS versions=3 items=6 events=3 actors=3")
			}
			if phase == "operation-create-slot" {
				var rows, keys, events, actors int
				e := p.QueryRow(ctx, `select count(*),
  (select count(*) from platform.idempotency_record where tenant_id=$1 and scope='franchise-slot'),
  (select count(*) from platform.outbox_event e join platform.idempotency_record i on i.tenant_id=e.tenant_id and i.resource_id=e.aggregate_id and i.scope='franchise-slot' where e.tenant_id=$1 and e.event_type='appointment-slot.created'),
  (select count(*) from platform.outbox_event e join platform.idempotency_record i on i.tenant_id=e.tenant_id and i.resource_id=e.aggregate_id and i.scope='franchise-slot' where e.tenant_id=$1 and e.event_type='appointment-slot.created' and e.payload->>'actor_subject'='planner')
  from crm.appointment_slot a join platform.idempotency_record i on i.tenant_id=a.tenant_id and i.resource_id=a.slot_id and i.scope='franchise-slot' where a.tenant_id=$1`, tenant).Scan(&rows, &keys, &events, &actors)
				if e != nil || rows != 3 || keys != 3 || events != 3 || actors != 3 {
					t.Fatalf("slot creation counts=%d/%d/%d/%d err=%v", rows, keys, events, actors, e)
				}
				t.Log("SLOT_CREATION_BROWSER_PASS slots=3 keys=3 events=3 actors=3")
			}
			if phase == "operation-create-resource" {
				var rows, keys, events, actors, skills int
				e := p.QueryRow(ctx, `select count(*),
  (select count(*) from platform.idempotency_record where tenant_id=$1 and scope='franchise-resource'),
  (select count(*) from platform.outbox_event e join platform.idempotency_record i on i.tenant_id=e.tenant_id and i.resource_id=e.aggregate_id and i.scope='franchise-resource' where e.tenant_id=$1 and e.event_type='service-resource.created'),
  (select count(*) from platform.outbox_event e join platform.idempotency_record i on i.tenant_id=e.tenant_id and i.resource_id=e.aggregate_id and i.scope='franchise-resource' where e.tenant_id=$1 and e.event_type='service-resource.created' and e.payload->>'actor_subject'='resource'),
  (select count(*) from crm.resource_skill rs join platform.idempotency_record i on i.tenant_id=rs.tenant_id and i.resource_id=rs.resource_id and i.scope='franchise-resource' where rs.tenant_id=$1)
  from crm.service_resource r join platform.idempotency_record i on i.tenant_id=r.tenant_id and i.resource_id=r.resource_id and i.scope='franchise-resource' where r.tenant_id=$1`, tenant).Scan(&rows, &keys, &events, &actors, &skills)
				if e != nil || rows != 3 || keys != 3 || events != 3 || actors != 3 || skills != 6 {
					t.Fatalf("resource creation counts=%d/%d/%d/%d/%d err=%v", rows, keys, events, actors, skills, e)
				}
				t.Log("RESOURCE_CREATION_BROWSER_PASS resources=3 keys=3 events=3 actors=3 skills=6")
			}
			if phase == "operation-create-availability" {
				var entries, keys, events, actors, cancelled int
				e := p.QueryRow(ctx, `select count(*),count(*)filter(where a.state='cancelled'),(select count(*) from platform.idempotency_record where tenant_id=$1 and scope='franchise-availability'),(select count(*) from platform.outbox_event e join platform.idempotency_record i on i.tenant_id=e.tenant_id and i.resource_id=e.aggregate_id and i.scope='franchise-availability' where e.tenant_id=$1 and e.event_type='availability-entry.created'),(select count(*) from platform.outbox_event e join platform.idempotency_record i on i.tenant_id=e.tenant_id and i.resource_id=e.aggregate_id and i.scope='franchise-availability' where e.tenant_id=$1 and e.event_type='availability-entry.created' and e.payload->>'actor_subject'='scheduler') from crm.availability_entry a join platform.idempotency_record i on i.tenant_id=a.tenant_id and i.resource_id=a.availability_id and i.scope='franchise-availability' where a.tenant_id=$1`, tenant).Scan(&entries, &cancelled, &keys, &events, &actors)
				if e != nil || entries != 3 || keys != 3 || events != 3 || actors != 3 || cancelled != 1 {
					t.Fatalf("availability creation counts=%d/%d/%d/%d cancelled=%d err=%v", entries, keys, events, actors, cancelled, e)
				}
				t.Log("AVAILABILITY_CREATION_BROWSER_PASS entries=3 keys=3 created_events=3 actors=3 cancelled=1")
			}
			if phase == "operation-quote" {
				var quotes, records, events, actors int
				err := p.QueryRow(ctx, `select (select count(*) from sales.quotation where tenant_id=$1 and lead_id like 'quote-create-%'),(select count(*) from platform.idempotency_record where tenant_id=$1 and scope='franchise-quote' and status='completed'),(select count(*) from platform.outbox_event where tenant_id=$1 and event_type='quotation.issued'),(select count(*) from platform.outbox_event where tenant_id=$1 and event_type='quotation.issued' and payload->>'actor_subject'='writer')`, tenant).Scan(&quotes, &records, &events, &actors)
				if err != nil || quotes != 3 || records != 3 || events != 3 || actors != 3 {
					t.Fatalf("quote creation effects quotes=%d keys=%d events=%d err=%v", quotes, records, events, err)
				}
				t.Log("QUOTE_CREATION_DURABLE_PASS quotes=3 keys=3 events=3 actors=3")
			}
			if phase == "operation-lead" {
				var state, assignee string
				var version, events, actors int
				if err := p.QueryRow(ctx, `select lifecycle_state,assigned_subject,version from crm.lead where tenant_id=$1 and lead_id='operator-lead'`, tenant).Scan(&state, &assignee, &version); err != nil || state != "lost" || assignee != "sales-owner" || version != 5 {
					t.Fatalf("lead final %s %s %d %v", state, assignee, version, err)
				}
				if err := p.QueryRow(ctx, `select count(*),count(*) filter(where payload->>'actor_subject'='lead-operator') from platform.outbox_event where tenant_id=$1 and aggregate_id='operator-lead'`, tenant).Scan(&events, &actors); err != nil || events != 4 || actors != 4 {
					t.Fatalf("lead audit events=%d actors=%d err=%v", events, actors, err)
				}
				t.Log("LEAD_RECOVERY_PASS version=5 events=4 authenticated_actors=4 no_duplicates")
			}
			if phase == "operation-availability" {
				var entries, actors, events, eventActors int
				if err := p.QueryRow(ctx, `select count(*),count(*) filter(where cancelled_by_subject='scheduler' and cancellation_reason_code='schedule-correction') from crm.availability_entry where tenant_id=$1 and availability_id like 'cancel-window-%' and state='cancelled' and version=2`, tenant).Scan(&entries, &actors); err != nil || entries != 4 || actors != 4 {
					t.Fatalf("availability entries=%d actors=%d err=%v", entries, actors, err)
				}
				if err := p.QueryRow(ctx, `select count(*),count(*) filter(where payload->>'actor_subject'='scheduler') from platform.outbox_event where tenant_id=$1 and aggregate_id like 'cancel-window-%' and event_type='availability-entry.cancelled'`, tenant).Scan(&events, &eventActors); err != nil || events != 4 || eventActors != 4 {
					t.Fatalf("availability events=%d actors=%d err=%v", events, eventActors, err)
				}
				t.Log("AVAILABILITY_RECOVERY_PASS cancelled=4 version=2 events=4 authenticated_actors=4")
				var blockedState string
				var blockedVersion, blockedEvents int
				if err := p.QueryRow(ctx, `select state,version,(select count(*) from platform.outbox_event where tenant_id=$1 and aggregate_id='cancel-window-4' and event_type='availability-entry.cancelled') from crm.availability_entry where tenant_id=$1 and availability_id='cancel-window-4'`, tenant).Scan(&blockedState, &blockedVersion, &blockedEvents); err != nil || blockedState != "active" || blockedVersion != 1 || blockedEvents != 0 {
					t.Fatalf("blocked interval=%s/%d events=%d err=%v", blockedState, blockedVersion, blockedEvents, err)
				}
				t.Log("AVAILABILITY_BOOKING_GUARD_PASS active=1 version=1 cancelled_events=0")
				var bookingEffects int
				if err := p.QueryRow(ctx, `select (select count(*) from crm.appointment where tenant_id=$1 and appointment_id<>'blocked-window-appointment')+(select count(*) from platform.idempotency_record where tenant_id=$1 and scope='public-appointment')+(select count(*) from platform.outbox_event where tenant_id=$1 and event_type='appointment.requested')`, tenant).Scan(&bookingEffects); err != nil || bookingEffects != 0 {
					t.Fatalf("booking outside working availability effects=%d err=%v", bookingEffects, err)
				}
				t.Log("AVAILABILITY_PUBLIC_HTTP_PASS withdrawn_slots=2 rejected_bookings=2 durable_booking_effects=0")
			}
			if phase == "operation-resolution" {
				// Existing delivery recovery remains a separate phase.
				for _, q := range []string{
					`select count(*) from sales.delivery_exception where tenant_id=$1 and handover_id='reject-fixture-handover' and state='resolved' and version=2`,
					`select count(*) from sales.delivery_exception_resolution r join sales.delivery_exception e on e.tenant_id=r.tenant_id and e.exception_id=r.exception_id where r.tenant_id=$1 and e.handover_id='reject-fixture-handover' and r.action='correct-and-represent' and r.resolved_by_subject='handover'`,
					`select count(*) from sales.delivery_handover where tenant_id=$1 and supersedes_handover_id='reject-fixture-handover' and state='prepared' and version=1`,
					`select count(*) from platform.outbox_event o join sales.delivery_exception e on e.tenant_id=o.tenant_id and e.exception_id=o.aggregate_id where o.tenant_id=$1 and e.handover_id='reject-fixture-handover' and o.event_type='delivery-exception.resolved'`,
				} {
					var n int
					if err := p.QueryRow(ctx, q, tenant).Scan(&n); err != nil || n != 1 {
						t.Fatalf("resolution durable count=%d err=%v", n, err)
					}
				}
				for _, action := range []string{"return", "exchange"} {
					for _, q := range []string{
						`select count(*) from sales.delivery_exception where tenant_id=$1 and handover_id=$2 and state='resolved' and version=2`,
						`select count(*) from sales.delivery_exception_resolution r join sales.delivery_exception e on e.tenant_id=r.tenant_id and e.exception_id=r.exception_id where r.tenant_id=$1 and e.handover_id=$2 and r.resolved_by_subject='handover'`,
						`select count(*) from sales.return_authorization where tenant_id=$1 and handover_id=$2 and state='authorized'`,
						`select count(*) from platform.outbox_event o join sales.delivery_exception e on e.tenant_id=o.tenant_id and e.exception_id=o.aggregate_id where o.tenant_id=$1 and e.handover_id=$2 and o.event_type='delivery-exception.resolved'`,
					} {
						var n int
						if err := p.QueryRow(ctx, q, tenant, action+"-fixture-handover").Scan(&n); err != nil || n != 1 {
							t.Fatalf("resolution %s count=%d err=%v", action, n, err)
						}
					}
				}
				t.Log("RESOLUTION_RECOVERY_PASS resolutions=3 successor=1 authorizations=2 events=3 actor=handover")
			}
			if phase == "operation-sections" {
				var forbidden int
				if err := p.QueryRow(ctx, `select count(*) from sales.delivery_checklist_template where tenant_id=$1 and checklist_id=$2`, tenant, "denied-probe").Scan(&forbidden); err != nil || forbidden != 0 {
					t.Fatalf("forbidden checklist effect: %d %v", forbidden, err)
				}
				t.Log("OPERATOR_SECTIONS_PASS roles=4 outage_recovered no_forbidden_checklist")
			}
			if phase == "delivery-action-loss" {
				for _, q := range []string{
					`select count(*) from sales.delivery_handover where tenant_id=$1 and handover_id='read-fixture-handover' and state='accepted' and version=3`,
					`select count(*) from sales.delivery_handover where tenant_id=$1 and handover_id='reject-fixture-handover' and state='rejected' and version=3`,
					`select count(*) from sales.delivery_exception where tenant_id=$1 and handover_id='reject-fixture-handover' and state='open'`,
					`select count(*) from platform.outbox_event where tenant_id=$1 and aggregate_id='read-fixture-handover' and event_type='delivery-handover.accepted'`,
					`select count(*) from platform.outbox_event where tenant_id=$1 and aggregate_id='reject-fixture-handover' and event_type='delivery-handover.rejected'`,
				} {
					var n int
					if err := p.QueryRow(ctx, q, tenant).Scan(&n); err != nil || n != 1 {
						t.Fatalf("delivery effect count=%d err=%v", n, err)
					}
				}
				t.Log("DELIVERY_ACTION_LOSS_PASS accepted=1 rejected=1 events=2 recovered_by_GET no_payment_or_shipment")
			}
			if strings.HasPrefix(phase, "delivery-read-") {
				var state string
				var events int
				if err := p.QueryRow(ctx, `select state,(select count(*) from platform.outbox_event where tenant_id=$1 and aggregate_id='read-fixture-handover') from sales.delivery_handover where tenant_id=$1 and handover_id='read-fixture-handover'`, tenant).Scan(&state, &events); err != nil || state != "prepared" || events != 0 {
					t.Fatalf("read mutated fixture state=%s events=%d err=%v", state, events, err)
				}
				t.Log("DELIVERY_READ_HTTP_POSTGRES_PASS phase=" + phase + " no_acceptance_no_outbox")
			}
			if phase == "restart" && os.Getenv("ELITE_ORDER_E2E") != "1" && os.Getenv("ELITE_PAYMENT_E2E") != "1" {
				verifyQuoteCommerceContinuation(t, ctx, p, tenant)
			}
			if strings.HasPrefix(phase, "operation") {
				for _, sql := range []string{
					`select count(*) from inventory.stock_unit where tenant_id=$1 and state='reserved'`,
					`select count(*) from inventory.serial_reservation where tenant_id=$1 and status='reservation'`,
					`select count(*) from sales.customer_order_line where tenant_id=$1 and allocated_stock_unit_id is not null`,
					`select count(*) from platform.outbox_event where tenant_id=$1 and event_type='order.stock-allocated' and payload->>'actor_subject'='operator'`,
				} {
					var n int
					if err := p.QueryRow(ctx, sql, tenant).Scan(&n); err != nil || n != 2 {
						t.Fatalf("operation durable invariant expected2 got%d err=%v", n, err)
					}
				}
				t.Log("ORDER_STOCK_HTTP_BROWSER_PASS phase=" + phase + " reservations=2 allocations=2 outbox=2")
			}
			if strings.HasPrefix(phase, "payment") {
				for _, sql := range []string{`select count(*) from payment.payment_attempt where tenant_id=$1 and state='created' and amount_minor_units=250000 and currency='ARS' and provider_code='stripe'`, `select count(*) from platform.outbox_event where tenant_id=$1 and event_type='payment.requested' and payload->>'actor_subject'='payment'`} {
					var n int
					if err := p.QueryRow(ctx, sql, tenant).Scan(&n); err != nil || n != 2 {
						t.Fatalf("payment invariant expected2 got%d err%v", n, err)
					}
				}
				t.Log("PAYMENT_UI_HTTP_POSTGRES_PASS phase=" + phase + " intents=2 audited_events=2 no_provider_calls")
			}
		}()
	}
}

// AUTHORED continuation of the SAME orders created through the customer browser.
// Exercises existing repository contracts, not operator UI or live payment approval.
func verifyQuoteCommerceContinuation(t *testing.T, ctx context.Context, pool *pgxpool.Pool, tenant string) {
	t.Helper()
	ids := randomid.Generator{}
	repo := postgres.NewCommerce(pool)
	type demand struct{ order, line string }
	demands := make([]demand, 2)
	for i, quote := range []string{"quote-success", "quote-race"} {
		if err := pool.QueryRow(ctx, `select o.order_id,l.line_id from sales.quotation q join sales.customer_order o on o.tenant_id=q.tenant_id and o.order_id=q.order_id join sales.customer_order_line l on l.tenant_id=o.tenant_id and l.order_id=o.order_id where q.tenant_id=$1 and q.quotation_id=$2`, tenant, quote).Scan(&demands[i].order, &demands[i].line); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := pool.Exec(ctx, `insert into inventory.stock_unit(tenant_id,stock_unit_id,organization_id,variant_id,serial_number,state,version,received_at) values($1,'quote-stock','store','variant','QUOTE-SERIAL','available',1,clock_timestamp())`, tenant); err != nil {
		t.Fatal(err)
	}
	if err := repo.AllocateStock(ctx, tenant, "other", demands[0].order, demands[0].line, "quote-stock", 1, 1, ids.New()); !errors.Is(err, commerce.ErrConflict) {
		t.Fatalf("cross-org allocation=%v", err)
	}
	if err := repo.AllocateStock(ctx, ids.New(), "store", demands[0].order, demands[0].line, "quote-stock", 1, 1, ids.New()); !errors.Is(err, commerce.ErrConflict) {
		t.Fatalf("cross-tenant allocation=%v", err)
	}
	if err := repo.AllocateStock(ctx, tenant, "store", demands[0].order, demands[0].line, "quote-stock", 99, 1, ids.New()); !errors.Is(err, commerce.ErrConflict) {
		t.Fatalf("stale allocation=%v", err)
	}
	start := make(chan struct{})
	results := make(chan error, 2)
	for _, d := range demands {
		go func(d demand) {
			<-start
			results <- repo.AllocateStock(ctx, tenant, "store", d.order, d.line, "quote-stock", 1, 1, ids.New())
		}(d)
	}
	close(start)
	success, conflict := 0, 0
	for range demands {
		err := <-results
		if err == nil {
			success++
		} else if errors.Is(err, commerce.ErrConflict) {
			conflict++
		} else {
			t.Fatal(err)
		}
	}
	if success != 1 || conflict != 1 {
		t.Fatalf("stock race success=%d conflict=%d", success, conflict)
	}
	var winner, line string
	if err := pool.QueryRow(ctx, `select order_id,line_id from sales.customer_order_line where tenant_id=$1 and allocated_stock_unit_id='quote-stock'`, tenant).Scan(&winner, &line); err != nil {
		t.Fatal(err)
	}
	if err := repo.AllocateStock(ctx, tenant, "store", winner, line, "quote-stock", 2, 2, ids.New()); !errors.Is(err, commerce.ErrConflict) {
		t.Fatalf("duplicate allocation=%v", err)
	}
	payment := commerce.PaymentAttempt{ID: ids.New(), OrderID: winner, OrganizationID: "other", ProviderCode: "synthetic-no-dispatch", Currency: "ARS", AmountMinorUnits: 250000}
	if err := repo.CreatePaymentAttempt(ctx, tenant, ids.New(), "quote-payment-key", payment); !errors.Is(err, commerce.ErrConflict) {
		t.Fatalf("cross-org payment=%v", err)
	}
	payment.OrganizationID = "store"
	payment.AmountMinorUnits = 1
	if err := repo.CreatePaymentAttempt(ctx, tenant, ids.New(), "quote-payment-key", payment); !errors.Is(err, commerce.ErrConflict) {
		t.Fatalf("wrong amount payment=%v", err)
	}
	payment.AmountMinorUnits = 250000
	payment.Currency = "USD"
	if err := repo.CreatePaymentAttempt(ctx, tenant, ids.New(), "quote-payment-key", payment); !errors.Is(err, commerce.ErrConflict) {
		t.Fatalf("wrong currency payment=%v", err)
	}
	payment.Currency = "ARS"
	if err := repo.CreatePaymentAttempt(ctx, tenant, ids.New(), "quote-payment-key", payment); err != nil {
		t.Fatal(err)
	}
	payment.ID = ids.New()
	if err := repo.CreatePaymentAttempt(ctx, tenant, ids.New(), "quote-payment-key", payment); err == nil {
		t.Fatal("duplicate payment key accepted")
	}
	// A fresh pool proves durable reads, not a PostgreSQL process restart/PITR.
	reopened, err := pgxpool.NewWithConfig(ctx, pool.Config().Copy())
	if err != nil {
		t.Fatal(err)
	}
	defer reopened.Close()
	for _, sql := range []string{
		`select count(*) from inventory.serial_reservation where tenant_id=$1 and stock_unit_id='quote-stock' and status='reservation'`,
		`select count(*) from inventory.stock_unit where tenant_id=$1 and stock_unit_id='quote-stock' and state='reserved' and version=2`,
		`select count(*) from sales.customer_order_line where tenant_id=$1 and allocated_stock_unit_id='quote-stock'`,
		`select count(*) from payment.payment_attempt where tenant_id=$1 and state='created' and currency='ARS' and amount_minor_units=250000`,
		`select count(*) from platform.outbox_event where tenant_id=$1 and event_type='order.stock-allocated'`,
		`select count(*) from platform.outbox_event where tenant_id=$1 and event_type='payment.requested'`,
	} {
		var n int
		if err := reopened.QueryRow(ctx, sql, tenant).Scan(&n); err != nil || n != 1 {
			t.Fatalf("commerce durable count=%d err=%v", n, err)
		}
	}
	t.Log("QUOTE_COMMERCE_CONTINUATION_PASS stock=1 competing_orders=2 payment_created=1 stock_events=1 payment_events=1 provider_calls=0")
}

func (r *journeyRepo) CreateQuoteAs(ctx context.Context, tenant, key string, v franchisejourney.Quote, hash, event, actor string) (franchisejourney.Quote, bool, error) {
	r.quoteActor = actor
	return r.CreateQuote(ctx, tenant, key, v, hash, event)
}

func (r *journeyRepo) CreateAvailabilityOnce(ctx context.Context, tenant, subject, key, hash string, v franchisejourney.AvailabilityEntry, event string) (franchisejourney.AvailabilityEntry, bool, error) {
	value, err := r.CreateAvailability(ctx, tenant, subject, v, event)
	return value, false, err
}

func TestResourceCreationHTTPRecovery(t *testing.T) {
	if os.Getenv("TEST_DATABASE_URL") == "" {
		t.Skip("disposable database required")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cfg, err := pgxpool.ParseConfig(os.Getenv("TEST_DATABASE_URL"))
	if err != nil || cfg.ConnConfig.Host != "127.0.0.1" || !strings.HasPrefix(cfg.ConnConfig.Database, "elite_confirmation_") {
		t.Fatal("requires disposable loopback database")
	}
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := randomid.Generator{}.New()
	for _, q := range []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1::uuid,'resource-'||replace(($1::uuid)::text,'-',''),'Synthetic','Synthetic')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'store','store','Synthetic','store')`,
	} {
		if _, err := pool.Exec(ctx, q, tenant); err != nil {
			t.Fatal(err)
		}
	}
	verifier, token := confirmationTestIssuer(t)
	mux := http.NewServeMux()
	FranchiseJourneyModule{Service: franchisejourney.NewService(postgres.NewFranchiseJourney(pool), randomid.Generator{}, browserAppointmentClock{})}.Register(mux, verifier)
	api := httptest.NewServer(mux)
	defer api.Close()
	bearer := token("resource-manager", tenant, []string{"resource:manage"}, []string{"store"})
	body := `{"organization_id":"store","display_name":"Synthetic bay","kind":"service-bay","skills":["service"]}`
	post := func() (int, map[string]any) {
		t.Helper()
		req, err := http.NewRequestWithContext(ctx, "POST", api.URL+"/v1/franchise/resources", strings.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("content-type", "application/json")
		req.Header.Set("authorization", "Bearer "+bearer)
		req.Header.Set("idempotency-key", "resource-recovery-fixture-key")
		res, err := api.Client().Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		var value map[string]any
		if err = json.NewDecoder(res.Body).Decode(&value); err != nil {
			t.Fatal(err)
		}
		return res.StatusCode, value
	}
	first, a := post()
	second, b := post()
	var count int
	if err = pool.QueryRow(ctx, `select count(*) from crm.service_resource where tenant_id=$1`, tenant).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if first != 201 || second != 200 || a["id"] != b["id"] || count != 1 {
		t.Fatalf("resource replay first=%d second=%d same_id=%v durable_resources=%d", first, second, a["id"] == b["id"], count)
	}
	var actor string
	if err = pool.QueryRow(ctx, `select payload->>'actor_subject' from platform.outbox_event where tenant_id=$1 and event_type='service-resource.created'`, tenant).Scan(&actor); err != nil || actor != "resource-manager" {
		t.Fatalf("actor not retained: %v", err)
	}
}

func TestSlotCreationHTTPRecovery(t *testing.T) {
	if os.Getenv("TEST_DATABASE_URL") == "" {
		t.Skip("disposable database required")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cfg, err := pgxpool.ParseConfig(os.Getenv("TEST_DATABASE_URL"))
	if err != nil || cfg.ConnConfig.Host != "127.0.0.1" || !strings.HasPrefix(cfg.ConnConfig.Database, "elite_confirmation_") {
		t.Fatal("requires disposable loopback database")
	}
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := randomid.Generator{}.New()
	for _, q := range []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1::uuid,'resource-'||replace(($1::uuid)::text,'-',''),'Synthetic','Synthetic')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'store','store','Synthetic','store')`,
	} {
		if _, err := pool.Exec(ctx, q, tenant); err != nil {
			t.Fatal(err)
		}
	}
	verifier, token := confirmationTestIssuer(t)
	mux := http.NewServeMux()
	FranchiseJourneyModule{Service: franchisejourney.NewService(postgres.NewFranchiseJourney(pool), randomid.Generator{}, browserAppointmentClock{})}.Register(mux, verifier)
	api := httptest.NewServer(mux)
	defer api.Close()
	bearer := token("slot-manager", tenant, []string{"appointment:manage"}, []string{"store"})
	start := time.Now().UTC().Add(72 * time.Hour).Truncate(time.Second)
	if _, err := postgres.NewFranchiseJourney(pool).CreateAvailability(ctx, tenant, "slot-manager", franchisejourney.AvailabilityEntry{ID: "working", OrganizationID: "store", EntryType: "working", StartsAt: start.Add(-time.Hour), EndsAt: start.Add(2 * time.Hour)}, randomid.Generator{}.New()); err != nil {
		t.Fatal(err)
	}
	payload, err := json.Marshal(map[string]any{"organization_id": "store", "kind": "service", "starts_at": start, "ends_at": start.Add(time.Hour), "capacity": 2})
	if err != nil {
		t.Fatal(err)
	}
	body := string(payload)
	post := func() (int, map[string]any) {
		t.Helper()
		req, err := http.NewRequestWithContext(ctx, "POST", api.URL+"/v1/franchise/appointment-slots", strings.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("content-type", "application/json")
		req.Header.Set("authorization", "Bearer "+bearer)
		req.Header.Set("idempotency-key", "slot-recovery-fixture-key")
		res, err := api.Client().Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		var value map[string]any
		if err = json.NewDecoder(res.Body).Decode(&value); err != nil {
			t.Fatal(err)
		}
		return res.StatusCode, value
	}
	first, a := post()
	second, b := post()
	var count int
	if err = pool.QueryRow(ctx, `select count(*) from crm.appointment_slot where tenant_id=$1`, tenant).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if first != 201 || second != 200 || a["id"] != b["id"] || count != 1 {
		t.Fatalf("slot replay first=%d second=%d same_id=%v durable_slots=%d", first, second, a["id"] == b["id"], count)
	}
	var actor string
	if err = pool.QueryRow(ctx, `select payload->>'actor_subject' from platform.outbox_event where tenant_id=$1 and event_type='appointment-slot.created'`, tenant).Scan(&actor); err != nil || actor != "slot-manager" {
		t.Fatalf("actor not retained: %v", err)
	}
}

func TestChecklistPublicationHTTPRecovery(t *testing.T) {
	if os.Getenv("TEST_DATABASE_URL") == "" {
		t.Skip("disposable database required")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cfg, err := pgxpool.ParseConfig(os.Getenv("TEST_DATABASE_URL"))
	if err != nil || cfg.ConnConfig.Host != "127.0.0.1" || !strings.HasPrefix(cfg.ConnConfig.Database, "elite_confirmation_") {
		t.Fatal("requires disposable loopback database")
	}
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := randomid.Generator{}.New()
	for _, q := range []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1::uuid,'checklist-'||replace(($1::uuid)::text,'-',''),'Synthetic','Synthetic')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'store','store','Synthetic','store')`,
	} {
		if _, err := pool.Exec(ctx, q, tenant); err != nil {
			t.Fatal(err)
		}
	}
	verifier, token := confirmationTestIssuer(t)
	mux := http.NewServeMux()
	FranchiseJourneyModule{Service: franchisejourney.NewService(postgres.NewFranchiseJourney(pool), randomid.Generator{}, browserAppointmentClock{})}.Register(mux, verifier)
	api := httptest.NewServer(mux)
	defer api.Close()
	bearer := token("checklist-publisher", tenant, []string{"handover:manage"}, []string{"store"})
	body := `{"organization_id":"store","checklist_id":"fixture-checklist","version":1,"title":"Synthetic checklist","items":[{"id":"serial","prompt":"Verify synthetic serial","response_type":"serial","required":true},{"id":"safe","prompt":"Confirm synthetic review","response_type":"confirmation","required":true}]}`
	req, err := http.NewRequestWithContext(ctx, "POST", api.URL+"/v1/franchise/delivery-checklists", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("authorization", "Bearer "+bearer)
	res, err := api.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	var published franchisejourney.DeliveryChecklist
	if res.StatusCode != 201 {
		res.Body.Close()
		t.Fatalf("publication status=%d", res.StatusCode)
	}
	if err = json.NewDecoder(res.Body).Decode(&published); err != nil {
		t.Fatal(err)
	}
	res.Body.Close()
	var actor string
	if err = pool.QueryRow(ctx, `select coalesce(payload->>'actor_subject','') from platform.outbox_event where tenant_id=$1 and event_type='delivery-checklist.published'`, tenant).Scan(&actor); err != nil || actor != "checklist-publisher" {
		t.Errorf("publication actor=%q error=%v", actor, err)
	}
	req, err = http.NewRequestWithContext(ctx, "GET", api.URL+"/v1/franchise/delivery-checklists/result?organization_id=store&checklist_id=fixture-checklist&version=1", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("authorization", "Bearer "+bearer)
	res, err = api.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		t.Fatalf("published version recovery status=%d expected=200", res.StatusCode)
	}
	var recovered franchisejourney.DeliveryChecklist
	if err = json.NewDecoder(res.Body).Decode(&recovered); err != nil {
		t.Fatal(err)
	}
	if recovered.ID != published.ID || recovered.Version != 1 || recovered.State != "published" || len(recovered.Items) != 2 || recovered.Items[0].Ordinal != 1 || recovered.Items[1].Ordinal != 2 {
		t.Fatalf("incoherent recovered checklist: %+v", recovered)
	}
}

func TestChecklistCompletionHTTPRecovery(t *testing.T) {
	raw := os.Getenv("TEST_DATABASE_URL")
	if raw == "" {
		t.Skip("disposable database required")
	}
	cfg, err := pgxpool.ParseConfig(raw)
	if err != nil || cfg.ConnConfig.Host != "127.0.0.1" || !strings.HasPrefix(cfg.ConnConfig.Database, "elite_confirmation_") {
		t.Fatal("requires disposable loopback database")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := randomid.Generator{}.New()
	for _, q := range []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1::uuid,'complete-'||replace(($1::uuid)::text,'-',''),'Synthetic','Synthetic')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'store','store','Synthetic','store')`,
		`insert into catalog.vehicle_model(tenant_id,model_id,model_code,display_name,vehicle_class,lifecycle_state)values($1,'model','model','Synthetic','bicycle','active')`,
		`insert into catalog.vehicle_variant(tenant_id,variant_id,model_id,variant_code,display_name,battery_specification,lifecycle_state)values($1,'variant','model','variant','Synthetic','{}','active')`,
		`insert into crm.customer_profile(tenant_id,customer_principal_id,display_name,email_normalized)values($1,'customer','Synthetic','synthetic@example.test')`,
		`insert into inventory.stock_unit(tenant_id,stock_unit_id,organization_id,variant_id,serial_number,state,version,received_at)values($1,'stock','store','variant','SYNTHETIC-SERIAL','sold',1,clock_timestamp())`,
		`insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version)values($1,'order','store','customer','delivered','ARS',123456,1)`,
		`insert into sales.delivery_handover(tenant_id,handover_id,organization_id,order_id,customer_principal_id,stock_unit_id,state,version)values($1,'completion-fixture','store','order','customer','stock','prepared',1)`,
	} {
		if _, err = pool.Exec(ctx, q, tenant); err != nil {
			t.Fatal(err)
		}
	}
	repo := postgres.NewFranchiseJourney(pool)
	checklist := franchisejourney.DeliveryChecklist{ID: "completion-checklist", OrganizationID: "store", Version: 1, Title: "Synthetic", Items: []franchisejourney.ChecklistItem{{ID: "serial", Ordinal: 1, Prompt: "Serie", ResponseType: "serial", Required: true}, {ID: "confirmed", Ordinal: 2, Prompt: "Confirmación", ResponseType: "confirmation", Required: true}}}
	if _, err = repo.PublishDeliveryChecklist(ctx, tenant, "fixture-author", checklist, randomid.Generator{}.New()); err != nil {
		t.Fatal(err)
	}
	verifier, token := confirmationTestIssuer(t)
	mux := http.NewServeMux()
	FranchiseJourneyModule{Service: franchisejourney.NewService(repo, randomid.Generator{}, browserAppointmentClock{})}.Register(mux, verifier)
	api := httptest.NewServer(mux)
	defer api.Close()
	bearer := token("checklist-operator", tenant, []string{"handover:manage"}, []string{"store"})
	body := `{"organization_id":"store","version":1,"checklist_id":"completion-checklist","checklist_version":1,"responses":[{"item_id":"serial","response_text":"SYNTHETIC-SERIAL"},{"item_id":"confirmed","response_text":"confirmed"}]}`
	for i := 0; i < 2; i++ {
		req, e := http.NewRequestWithContext(ctx, "POST", api.URL+"/v1/franchise/handovers/completion-fixture/complete-checklist", strings.NewReader(body))
		if e != nil {
			t.Fatal(e)
		}
		req.Header.Set("content-type", "application/json")
		req.Header.Set("authorization", "Bearer "+bearer)
		res, e := api.Client().Do(req)
		if e != nil {
			t.Fatal(e)
		}
		res.Body.Close()
		expected := 200
		if i == 1 {
			expected = 409
		}
		if res.StatusCode != expected {
			t.Fatalf("complete/replay status=%d want=%d", res.StatusCode, expected)
		}
	}
	var answers, events int
	if err = pool.QueryRow(ctx, `select (select count(*) from sales.delivery_checklist_response where tenant_id=$1),(select count(*) from platform.outbox_event where tenant_id=$1 and event_type='delivery-handover.checklist-completed' and payload->>'actor_subject'='checklist-operator')`, tenant).Scan(&answers, &events); err != nil || answers != 2 || events != 1 {
		t.Fatal(answers, events, err)
	}
	req, err := http.NewRequestWithContext(ctx, "GET", api.URL+"/v1/franchise/handovers/completion-fixture/checklist-result?organization_id=store", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("authorization", "Bearer "+bearer)
	res, err := api.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		t.Fatalf("completion lookup status=%d; durable answers=%d events=%d", res.StatusCode, answers, events)
	}
	var value struct {
		HandoverID  string                               `json:"handover_id"`
		State       string                               `json:"state"`
		Version     int64                                `json:"version"`
		ChecklistID string                               `json:"checklist_id"`
		Actor       string                               `json:"actor_subject"`
		Responses   []franchisejourney.ChecklistResponse `json:"responses"`
	}
	if err = json.NewDecoder(res.Body).Decode(&value); err != nil || value.HandoverID != "completion-fixture" || value.State != "presented" || value.Version != 2 || value.ChecklistID != checklist.ID || value.Actor != "checklist-operator" || len(value.Responses) != 2 {
		t.Fatal(value, err)
	}
	t.Log("CHECKLIST_COMPLETION_HTTP_PASS complete=200 replay=409 answers=2 events=1 lookup=200 actor=checklist-operator")
}

func TestReturnOperationsHTTPRecovery(t *testing.T) {
	raw := os.Getenv("TEST_DATABASE_URL")
	if raw == "" {
		t.Skip("disposable database required")
	}
	cfg, err := pgxpool.ParseConfig(raw)
	if err != nil || cfg.ConnConfig.Host != "127.0.0.1" || !strings.HasPrefix(cfg.ConnConfig.Database, "elite_confirmation_") {
		t.Fatal("requires disposable loopback database")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := randomid.Generator{}.New()
	for _, q := range []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1::uuid,'complete-'||replace(($1::uuid)::text,'-',''),'Synthetic','Synthetic')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'store','store','Synthetic','store')`,
		`insert into catalog.vehicle_model(tenant_id,model_id,model_code,display_name,vehicle_class,lifecycle_state)values($1,'model','model','Synthetic','bicycle','active')`,
		`insert into catalog.vehicle_variant(tenant_id,variant_id,model_id,variant_code,display_name,battery_specification,lifecycle_state)values($1,'variant','model','variant','Synthetic','{}','active')`,
		`insert into crm.customer_profile(tenant_id,customer_principal_id,display_name,email_normalized)values($1,'customer','Synthetic','synthetic@example.test')`,
		`insert into inventory.stock_unit(tenant_id,stock_unit_id,organization_id,variant_id,serial_number,state,version,received_at)values($1,'stock','store','variant','SYNTHETIC-SERIAL','sold',1,clock_timestamp())`,
		`insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version)values($1,'order','store','customer','delivered','ARS',123456,1)`,
		`insert into sales.delivery_handover(tenant_id,handover_id,organization_id,order_id,customer_principal_id,stock_unit_id,state,version)values($1,'completion-fixture','store','order','customer','stock','prepared',1)`,
	} {
		if _, err = pool.Exec(ctx, q, tenant); err != nil {
			t.Fatal(err)
		}
	}
	repo := postgres.NewFranchiseJourney(pool)
	checklist := franchisejourney.DeliveryChecklist{ID: "completion-checklist", OrganizationID: "store", Version: 1, Title: "Synthetic", Items: []franchisejourney.ChecklistItem{{ID: "serial", Ordinal: 1, Prompt: "Serie", ResponseType: "serial", Required: true}, {ID: "confirmed", Ordinal: 2, Prompt: "Confirmación", ResponseType: "confirmation", Required: true}}}
	if _, err = repo.PublishDeliveryChecklist(ctx, tenant, "fixture-author", checklist, randomid.Generator{}.New()); err != nil {
		t.Fatal(err)
	}

	responses := []franchisejourney.ChecklistResponse{{ItemID: "serial", ResponseText: "SYNTHETIC-SERIAL"}, {ItemID: "confirmed", ResponseText: "confirmed"}}
	if _, err = repo.CompleteDeliveryChecklist(ctx, tenant, "store", "fixture-operator", "completion-fixture", 1, checklist.ID, 1, responses, randomid.Generator{}.New()); err != nil {
		t.Fatal(err)
	}
	exception, err := repo.RejectHandover(ctx, tenant, "store", "customer", "completion-fixture", 2, "serial-mismatch", "Synthetic rejection", strings.Repeat("a", 64), randomid.Generator{}.New(), randomid.Generator{}.New())
	if err != nil {
		t.Fatal(err)
	}
	resolution, err := repo.ResolveDeliveryException(ctx, tenant, "store", "fixture-manager", exception.ID, 1, "return", "Synthetic authorization", "", "return-authorization", randomid.Generator{}.New(), randomid.Generator{}.New())
	if err != nil || resolution.ReturnAuthorizationID != "return-authorization" {
		t.Fatal(resolution, err)
	}
	verifier, token := confirmationTestIssuer(t)
	mux := http.NewServeMux()
	FranchiseJourneyModule{Service: franchisejourney.NewService(repo, randomid.Generator{}, browserAppointmentClock{})}.Register(mux, verifier)
	api := httptest.NewServer(mux)
	defer api.Close()
	bearer := token("return-operator", tenant, []string{"handover:manage"}, []string{"store"})
	post := func(path, body string, target any) int {
		req, e := http.NewRequestWithContext(ctx, "POST", api.URL+path, strings.NewReader(body))
		if e != nil {
			t.Fatal(e)
		}
		req.Header.Set("content-type", "application/json")
		req.Header.Set("authorization", "Bearer "+bearer)
		res, e := api.Client().Do(req)
		if e != nil {
			t.Fatal(e)
		}
		defer res.Body.Close()
		if res.StatusCode == 200 && target != nil {
			if e = json.NewDecoder(res.Body).Decode(target); e != nil {
				t.Fatal(e)
			}
		}
		return res.StatusCode
	}
	body := fmt.Sprintf(`{"organization_id":"store","serial_number":"SYNTHETIC-SERIAL","condition_code":"sealed","notes":"Synthetic receipt","evidence_sha256":"%s"}`, strings.Repeat("b", 64))
	var receipt franchisejourney.ReturnReceipt
	if code := post("/v1/franchise/return-authorizations/return-authorization/receive", body, &receipt); code != 200 || receipt.ID == "" {
		t.Fatal("receive", code, receipt)
	}
	if code := post("/v1/franchise/return-authorizations/return-authorization/receive", body, nil); code != 409 {
		t.Fatal("receive replay", code)
	}
	lookup := func(label string, decided bool) {
		req, e := http.NewRequestWithContext(ctx, "GET", api.URL+"/v1/franchise/returns/result?organization_id=store&authorization_id=return-authorization", nil)
		if e != nil {
			t.Fatal(e)
		}
		req.Header.Set("authorization", "Bearer "+bearer)
		res, e := api.Client().Do(req)
		if e != nil {
			t.Fatal(e)
		}
		defer res.Body.Close()
		if res.StatusCode != 200 {
			t.Errorf("%s lookup=%d after durable commit", label, res.StatusCode)
			return
		}
		var value franchisejourney.ReturnCase
		if e = json.NewDecoder(res.Body).Decode(&value); e != nil {
			t.Fatal(e)
		}
		if value.AuthorizationID != "return-authorization" || value.OrganizationID != "store" || value.Receipt == nil || value.Receipt.ID != receipt.ID || value.Receipt.ReceivedBySubject != "return-operator" {
			t.Fatal("receipt lookup", value)
		}
		if decided && (value.Disposition == nil || value.Disposition.ReceiptID != receipt.ID || value.Disposition.DecidedBySubject != "return-operator" || len(value.Disposition.Effects) != 4) {
			t.Fatal("decision lookup", value)
		}
		if !decided && value.Disposition != nil {
			t.Fatal("uncreated decision exposed")
		}
	}
	lookup("receipt", false)
	body = `{"organization_id":"store","inventory_action":"quarantine","notes":"Synthetic decision"}`
	var disposition franchisejourney.ReturnDisposition
	if code := post("/v1/franchise/return-receipts/"+receipt.ID+"/decide", body, &disposition); code != 200 || disposition.ID == "" || len(disposition.Effects) != 4 {
		t.Fatal("decision", code, disposition)
	}
	if code := post("/v1/franchise/return-receipts/"+receipt.ID+"/decide", body, nil); code != 409 {
		t.Fatal("decision replay", code)
	}
	lookup("decision", true)
	var receipts, decisions, effects, events int
	err = pool.QueryRow(ctx, `select (select count(*) from sales.return_receipt where tenant_id=$1),(select count(*) from sales.return_disposition where tenant_id=$1),(select count(*) from sales.return_effect_request where tenant_id=$1),(select count(*) from platform.outbox_event where tenant_id=$1 and payload->>'actor_subject'='return-operator')`, tenant).Scan(&receipts, &decisions, &effects, &events)
	if err != nil || receipts != 1 || decisions != 1 || effects != 4 || events != 2 {
		t.Fatal(receipts, decisions, effects, events, err)
	}
	if !t.Failed() {
		t.Log("RETURN_OPERATIONS_HTTP_PASS receive=200/replay409 decide=200/replay409 receipts=1 decisions=1 effects=4 events=2 lookup=200")
	}
}
````

### FILE: `db/migrations/0007_appointment_capacity.up.sql`

```yaml
block_id: "GO-FRANCHISE-JOURNEY:file:13"
operation: CREATE
provenance: AUTHORED
source: "local PostgreSQL capacity implementation governed by pinned Microsoft BCApps resource capacity and availability tests"
license: "LicenseRef-Workspace-Owner"
sha256: "ea21676c7f4c2b3a773649090d21b2882adaba5d31b1cc41f2c5c043a57668d0"
variables: []
secrets_allowed: false
```

````sql
begin;

create table crm.appointment_slot (
  tenant_id uuid not null,
  slot_id text not null,
  organization_id text not null,
  appointment_kind text not null check (appointment_kind in ('consultation','test-drive','delivery','service')),
  starts_at timestamptz not null,
  ends_at timestamptz not null,
  capacity integer not null check (capacity between 1 and 100),
  state text not null check (state in ('open','closed')),
  version bigint not null check (version > 0),
  created_at timestamptz not null default clock_timestamp(),
  updated_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, slot_id),
  unique (tenant_id, organization_id, appointment_kind, starts_at),
  foreign key (tenant_id, organization_id) references org.organization (tenant_id, organization_id),
  check (ends_at > starts_at),
  check (ends_at - starts_at <= interval '8 hours'),
  check (updated_at >= created_at)
);

alter table crm.appointment
  add column slot_id text,
  add column ends_at timestamptz,
  add constraint appointment_slot_fk foreign key (tenant_id, slot_id) references crm.appointment_slot (tenant_id, slot_id),
  add constraint appointment_slot_time_check check ((slot_id is null and ends_at is null) or (slot_id is not null and ends_at > starts_at));

create index appointment_slot_public_idx on crm.appointment_slot (tenant_id, organization_id, appointment_kind, state, starts_at);
create index appointment_slot_booking_idx on crm.appointment (tenant_id, slot_id, state) where slot_id is not null;

create or replace function crm.enforce_appointment_slot_non_overlap()
returns trigger
language plpgsql
as $$
begin
  perform pg_advisory_xact_lock(hashtextextended(new.tenant_id::text || '|' || new.organization_id || '|' || new.appointment_kind, 0));
  if new.state = 'open' and exists (
    select 1
    from crm.appointment_slot existing
    where existing.tenant_id = new.tenant_id
      and existing.organization_id = new.organization_id
      and existing.appointment_kind = new.appointment_kind
      and existing.state = 'open'
      and existing.slot_id <> new.slot_id
      and tstzrange(existing.starts_at, existing.ends_at, '[)') && tstzrange(new.starts_at, new.ends_at, '[)')
  ) then
    raise exception using errcode = '23P01', message = 'appointment slots overlap';
  end if;
  return new;
end;
$$;

create trigger appointment_slot_non_overlap
before insert or update of organization_id,appointment_kind,starts_at,ends_at,state
on crm.appointment_slot
for each row execute function crm.enforce_appointment_slot_non_overlap();

commit;
````

### FILE: `db/migrations/0007_appointment_capacity.down.sql`

```yaml
block_id: "GO-FRANCHISE-JOURNEY:file:14"
operation: CREATE
provenance: AUTHORED
source: "local reversible migration governed by pinned Microsoft BCApps resource capacity and availability tests"
license: "LicenseRef-Workspace-Owner"
sha256: "22ff303b1b479b98d208c67e6b231203c8346f495a9a46e40cac1e10576f8465"
variables: []
secrets_allowed: false
```

````sql
begin;
drop index if exists crm.appointment_slot_booking_idx;
alter table crm.appointment drop constraint if exists appointment_slot_time_check, drop constraint if exists appointment_slot_fk, drop column if exists ends_at, drop column if exists slot_id;
drop table if exists crm.appointment_slot;
drop function if exists crm.enforce_appointment_slot_non_overlap();
commit;
````

### FILE: `db/tests/0007_appointment_capacity.test.sql`

```yaml
block_id: "GO-FRANCHISE-JOURNEY:file:15"
operation: CREATE
provenance: AUTHORED
source: "local invariant test governed by pinned Microsoft BCApps capacity, zero-capacity and availability tests"
license: "LicenseRef-Workspace-Owner"
sha256: "43516b4dcbf9b47c938a9f066370d50a9c815bd8901987026b26019b54ef2ec0"
variables: []
secrets_allowed: false
```

````sql
begin;
do $$
begin
  if not exists (select 1 from information_schema.tables where table_schema='crm' and table_name='appointment_slot') then raise exception 'appointment_slot missing'; end if;
  if not exists (select 1 from pg_constraint where conname='appointment_slot_time_check') then raise exception 'appointment slot binding missing'; end if;
  if not exists (select 1 from pg_indexes where schemaname='crm' and indexname='appointment_slot_public_idx') then raise exception 'appointment slot public index missing'; end if;
  if to_regprocedure('crm.enforce_appointment_slot_non_overlap()') is null then raise exception 'appointment slot overlap guard missing'; end if;
end $$;
rollback;
````

### FILE: `db/migrations/0009_appointment_resources.up.sql`

```yaml
block_id: "GO-FRANCHISE-JOURNEY:db-migrations-0009-appointment-resources-up-sql:v4"
operation: CREATE
provenance: AUTHORED
source: "local verified composition governed by pinned Microsoft BCApps Resource, Capacity and ResourceSkill sources/tests"
license: "LicenseRef-Workspace-Owner"
sha256: "45baa236e888947ae6c4be589efe6cae54b1709f9fafa5aea2475e514c057231"
variables: []
secrets_allowed: false
```

````sql
begin;

create table crm.service_resource (
  tenant_id uuid not null,
  resource_id text not null,
  organization_id text not null,
  principal_subject text,
  display_name text not null,
  resource_kind text not null check (resource_kind in ('employee','contractor','service-bay','vehicle','equipment')),
  status text not null check (status in ('active','inactive')),
  version bigint not null check (version > 0),
  created_at timestamptz not null default clock_timestamp(),
  updated_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id,resource_id),
  foreign key (tenant_id,organization_id) references org.organization(tenant_id,organization_id),
  check (length(display_name) between 1 and 160),
  check ((resource_kind in ('employee','contractor')) = (principal_subject is not null))
);

create unique index service_resource_principal_idx on crm.service_resource(tenant_id,organization_id,principal_subject) where principal_subject is not null;

create table crm.resource_skill (
  tenant_id uuid not null,
  resource_id text not null,
  appointment_kind text not null check (appointment_kind in ('consultation','test-drive','delivery','service')),
  primary key (tenant_id,resource_id,appointment_kind),
  foreign key (tenant_id,resource_id) references crm.service_resource(tenant_id,resource_id) on delete cascade
);

create table crm.appointment_resource (
  tenant_id uuid not null,
  appointment_id text not null,
  resource_id text not null,
  assigned_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id,appointment_id),
  foreign key (tenant_id,appointment_id) references crm.appointment(tenant_id,appointment_id) on delete cascade,
  foreign key (tenant_id,resource_id) references crm.service_resource(tenant_id,resource_id)
);

create or replace function crm.enforce_appointment_resource_assignment()
returns trigger language plpgsql as $$
declare candidate crm.appointment%rowtype;
begin
  perform pg_advisory_xact_lock(hashtextextended(new.tenant_id::text || '|appointment-resource|' || new.resource_id,0));
  select * into candidate from crm.appointment a where a.tenant_id=new.tenant_id and a.appointment_id=new.appointment_id;
  if candidate.appointment_id is null or candidate.state <> 'requested' then
    raise exception using errcode='23514', message='appointment is unavailable for resource assignment';
  end if;
  if not exists (
    select 1 from crm.service_resource r join crm.resource_skill s on s.tenant_id=r.tenant_id and s.resource_id=r.resource_id
    where r.tenant_id=new.tenant_id and r.resource_id=new.resource_id and r.organization_id=candidate.organization_id and r.status='active' and s.appointment_kind=candidate.appointment_kind
  ) then
    raise exception using errcode='23514', message='resource is unavailable or lacks appointment skill';
  end if;
  if exists (
    select 1 from crm.appointment_resource ar join crm.appointment a on a.tenant_id=ar.tenant_id and a.appointment_id=ar.appointment_id
    where ar.tenant_id=new.tenant_id and ar.resource_id=new.resource_id and ar.appointment_id<>new.appointment_id and a.state in ('requested','confirmed')
      and tstzrange(a.starts_at,coalesce(a.ends_at,a.starts_at+interval '1 hour'),'[)') && tstzrange(candidate.starts_at,coalesce(candidate.ends_at,candidate.starts_at+interval '1 hour'),'[)')
  ) then
    raise exception using errcode='23P01', message='resource appointment assignments overlap';
  end if;
  return new;
end;
$$;

create trigger appointment_resource_assignment_guard before insert or update of resource_id on crm.appointment_resource for each row execute function crm.enforce_appointment_resource_assignment();
create index service_resource_org_status_idx on crm.service_resource(tenant_id,organization_id,status,resource_id);
create index appointment_resource_schedule_idx on crm.appointment_resource(tenant_id,resource_id,appointment_id);

commit;
````

### FILE: `db/migrations/0009_appointment_resources.down.sql`

```yaml
block_id: "GO-FRANCHISE-JOURNEY:db-migrations-0009-appointment-resources-down-sql:v4"
operation: CREATE
provenance: AUTHORED
source: "local reversible migration governed by pinned Microsoft BCApps Resource, Capacity and ResourceSkill sources/tests"
license: "LicenseRef-Workspace-Owner"
sha256: "46b3d7eac15a1e58369e34c6feb1b76e2aa75545c3ff94b5edd21809e4b3de0b"
variables: []
secrets_allowed: false
```

````sql
begin;
drop index if exists crm.appointment_resource_schedule_idx;
drop index if exists crm.service_resource_org_status_idx;
drop trigger if exists appointment_resource_assignment_guard on crm.appointment_resource;
drop function if exists crm.enforce_appointment_resource_assignment();
drop table if exists crm.appointment_resource;
drop table if exists crm.resource_skill;
drop table if exists crm.service_resource;
commit;
````

### FILE: `db/tests/0009_appointment_resources.test.sql`

```yaml
block_id: "GO-FRANCHISE-JOURNEY:db-tests-0009-appointment-resources-test-sql:v4"
operation: CREATE
provenance: AUTHORED
source: "local invariant test governed by pinned Microsoft BCApps ResourceSkill test authority"
license: "LicenseRef-Workspace-Owner"
sha256: "b09bd47272782104de634a4db29593d88ede6698f56a39fc6d62effe3bb24553"
variables: []
secrets_allowed: false
```

````sql
begin;
do $$
begin
  if to_regclass('crm.service_resource') is null then raise exception 'service_resource missing'; end if;
  if to_regclass('crm.resource_skill') is null then raise exception 'resource_skill missing'; end if;
  if to_regclass('crm.appointment_resource') is null then raise exception 'appointment_resource missing'; end if;
  if not exists(select 1 from pg_indexes where schemaname='crm' and indexname='service_resource_principal_idx' and indexdef like '%WHERE (principal_subject IS NOT NULL)%') then raise exception 'partial principal uniqueness missing'; end if;
  if to_regprocedure('crm.enforce_appointment_resource_assignment()') is null then raise exception 'assignment guard missing'; end if;
end;
$$;
rollback;
````

### FILE: `internal/franchisejourney/availability.go`

```yaml
block_id: "GO-FRANCHISE-JOURNEY:file:19"
operation: CREATE
provenance: AUTHORED
source: "local implementation governed by Microsoft Business Central employee absence and calendar references pinned in metadata"
license: "LicenseRef-Workspace-Owner"
sha256: "cb7de6b97ba5074c4a275241004ff8988e90c116f39b4d604d3b1185903f7e6b"
variables: []
secrets_allowed: false
```

````go
package franchisejourney

import (
	"context"
	"regexp"
	"time"
)

type AvailabilityEntry struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organization_id"`
	ResourceID     string    `json:"resource_id,omitempty"`
	EntryType      string    `json:"entry_type"`
	ReasonCode     string    `json:"reason_code,omitempty"`
	StartsAt       time.Time `json:"starts_at"`
	EndsAt         time.Time `json:"ends_at"`
	State          string    `json:"state"`
	Version        int64     `json:"version"`
}

func (s *Service) CreateAvailability(ctx context.Context, tenant, subject string, value AvailabilityEntry) (AvailabilityEntry, error) {
	value, err := s.prepareAvailability(tenant, subject, value)
	if err != nil {
		return AvailabilityEntry{}, err
	}
	return s.repository.CreateAvailability(ctx, tenant, subject, value, s.ids.New())
}
func (s *Service) prepareAvailability(tenant, subject string, value AvailabilityEntry) (AvailabilityEntry, error) {
	value.ID, value.State, value.Version = s.ids.New(), "active", 1
	validType := value.EntryType == "working" || value.EntryType == "unavailable"
	validReason := value.EntryType == "working" && value.ReasonCode == "" || value.EntryType == "unavailable" && code(value.ReasonCode)
	if tenant == "" || subject == "" || value.OrganizationID == "" || !validType || !validReason || value.StartsAt.IsZero() || !value.EndsAt.After(value.StartsAt) || value.EndsAt.Sub(value.StartsAt) > 366*24*time.Hour {
		return AvailabilityEntry{}, ErrInvalid
	}
	return value, nil
}

func (s *Service) CancelAvailability(ctx context.Context, tenant, organization, entry string, version int64, subject, reason string) (AvailabilityEntry, error) {
	if tenant == "" || organization == "" || entry == "" || version < 1 || subject == "" || !code(reason) {
		return AvailabilityEntry{}, ErrInvalid
	}
	return s.repository.CancelAvailability(ctx, tenant, organization, entry, version, subject, reason, s.ids.New())
}

func (s *Service) Availability(ctx context.Context, tenant, organization, resource string, from, to time.Time) ([]AvailabilityEntry, error) {
	if tenant == "" || organization == "" || from.IsZero() || !to.After(from) || to.Sub(from) > 366*24*time.Hour {
		return nil, ErrInvalid
	}
	return s.repository.Availability(ctx, tenant, organization, resource, from, to)
}

func (s *Service) CancelCustomerAppointment(ctx context.Context, tenant, organization, customer, appointment string, version int64, reason string) (Appointment, error) {
	if tenant == "" || organization == "" || customer == "" || appointment == "" || version < 1 || !code(reason) {
		return Appointment{}, ErrInvalid
	}
	return s.repository.CancelCustomerAppointment(ctx, tenant, organization, customer, appointment, version, reason, s.ids.New(), s.ids.New())
}

var availabilityKeyPattern = regexp.MustCompile(`^[A-Za-z0-9_-]{16,128}$`)

func (s *Service) CreateAvailabilityOnce(ctx context.Context, tenant, subject, key, hash string, value AvailabilityEntry) (AvailabilityEntry, bool, error) {
	if !availabilityKeyPattern.MatchString(key) || !sha256Hex(hash) {
		return AvailabilityEntry{}, false, ErrInvalid
	}
	value, err := s.prepareAvailability(tenant, subject, value)
	if err != nil {
		return AvailabilityEntry{}, false, err
	}
	r, ok := s.repository.(interface {
		CreateAvailabilityOnce(context.Context, string, string, string, string, AvailabilityEntry, string) (AvailabilityEntry, bool, error)
	})
	if !ok {
		return AvailabilityEntry{}, false, ErrConflict
	}
	return r.CreateAvailabilityOnce(ctx, tenant, subject, key, hash, value, s.ids.New())
}
func (s *Service) AvailabilityCreationResult(ctx context.Context, tenant, organization, key string) (AvailabilityEntry, error) {
	if tenant == "" || organization == "" || !availabilityKeyPattern.MatchString(key) {
		return AvailabilityEntry{}, ErrInvalid
	}
	r, ok := s.repository.(interface {
		AvailabilityCreationResult(context.Context, string, string, string) (AvailabilityEntry, error)
	})
	if !ok {
		return AvailabilityEntry{}, ErrConflict
	}
	return r.AvailabilityCreationResult(ctx, tenant, organization, key)
}
````

### FILE: `internal/franchisejourney/availability_test.go`

```yaml
block_id: "GO-FRANCHISE-JOURNEY:file:20"
operation: CREATE
provenance: AUTHORED
source: "local invariant tests governed by Microsoft Business Central employee absence and calendar references pinned in metadata"
license: "LicenseRef-Workspace-Owner"
sha256: "4d6acadf4f6cd184298e40e3e5ab3e3e039e2267ba54c3d37bedbe686a01cbbc"
variables: []
secrets_allowed: false
```

````go
package franchisejourney

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestAvailabilityRequiresExplicitIntervalsAndReasons(t *testing.T) {
	service, repository, _ := newService()
	from := time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC)
	value, err := service.CreateAvailability(context.Background(), "tenant", "scheduler", AvailabilityEntry{OrganizationID: "store", ResourceID: "employee", EntryType: "unavailable", ReasonCode: "annual-leave", StartsAt: from, EndsAt: from.Add(8 * time.Hour)})
	if err != nil || value.State != "active" || repository.availability.EntryType != "unavailable" {
		t.Fatalf("value=%+v err=%v", value, err)
	}
	if _, err = service.CreateAvailability(context.Background(), "tenant", "scheduler", AvailabilityEntry{OrganizationID: "store", EntryType: "unavailable", StartsAt: from, EndsAt: from.Add(time.Hour)}); !errors.Is(err, ErrInvalid) {
		t.Fatalf("reasonless absence accepted: %v", err)
	}
	if _, err = service.TransitionAppointment(context.Background(), "tenant", "store", "appointment", "confirmed", "no-show", 2, "operator", ""); !errors.Is(err, ErrInvalid) {
		t.Fatalf("reasonless no-show accepted: %v", err)
	}
	if _, err = service.TransitionAppointment(context.Background(), "tenant", "store", "appointment", "confirmed", "completed", 2, "operator", "invented-reason"); !errors.Is(err, ErrInvalid) {
		t.Fatalf("irrelevant reason accepted: %v", err)
	}
}
````

### FILE: `internal/platform/postgres/franchisejourney_availability.go`

```yaml
block_id: "GO-FRANCHISE-JOURNEY:file:21"
operation: CREATE
provenance: AUTHORED
source: "local transactional repository governed by Microsoft Business Central calendars and PostgreSQL 18 locking references pinned in metadata"
license: "LicenseRef-Workspace-Owner"
sha256: "8b28584844fb74c6697d04f20d809e06ea121a4748be03f317886df962bd0e50"
variables: []
secrets_allowed: false
```

````go
package postgres

import (
	"context"
	"errors"
	"time"

	"elite.local/enterprise/internal/franchisejourney"
	"github.com/jackc/pgx/v5"
)

func (r *FranchiseJourney) CreateAvailability(ctx context.Context, tenant, subject string, value franchisejourney.AvailabilityEntry, eventID string) (franchisejourney.AvailabilityEntry, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return value, err
	}
	defer tx.Rollback(ctx)
	value, err = r.createAvailabilityTx(ctx, tx, tenant, subject, value, eventID)
	if err != nil {
		return value, err
	}
	return value, tx.Commit(ctx)
}

func (r *FranchiseJourney) createAvailabilityTx(ctx context.Context, tx pgx.Tx, tenant, subject string, value franchisejourney.AvailabilityEntry, eventID string) (franchisejourney.AvailabilityEntry, error) {
	// Match the existing availability trigger's organization lock before taking
	// the statement snapshot. Booking uses this same transaction-scoped fence.
	if value.ResourceID == "" {
		if err := lockOrganizationSchedule(ctx, tx, tenant, value.OrganizationID); err != nil {
			return value, err
		}
	}
	err := tx.QueryRow(ctx, `insert into crm.availability_entry(tenant_id,availability_id,organization_id,resource_id,entry_type,reason_code,starts_at,ends_at,state,version,created_by_subject)
		select $1,$2,o.organization_id,nullif($4,''),$5,nullif($6,''),$7,$8,'active',1,$9 from org.organization o
		where o.tenant_id=$1 and o.organization_id=$3 and o.status='active' and ($4='' or exists(select 1 from crm.service_resource r where r.tenant_id=$1 and r.organization_id=$3 and r.resource_id=$4 and r.status='active'))
		returning availability_id,organization_id,coalesce(resource_id,''),entry_type,coalesce(reason_code,''),starts_at,ends_at,state,version`, tenant, value.ID, value.OrganizationID, value.ResourceID, value.EntryType, value.ReasonCode, value.StartsAt, value.EndsAt, subject).Scan(&value.ID, &value.OrganizationID, &value.ResourceID, &value.EntryType, &value.ReasonCode, &value.StartsAt, &value.EndsAt, &value.State, &value.Version)
	if errors.Is(err, pgx.ErrNoRows) || postgresConflict(err) {
		return value, franchisejourney.ErrConflict
	}
	if err != nil {
		return value, err
	}
	if err = writeJourneyOutbox(ctx, tx, tenant, eventID, "availability-entry", value.ID, 1, "availability-entry.created", map[string]any{"organization_id": value.OrganizationID, "resource_id": value.ResourceID, "entry_type": value.EntryType, "reason_code": value.ReasonCode, "starts_at": value.StartsAt, "ends_at": value.EndsAt, "actor_subject": subject}); err != nil {
		return value, err
	}
	return value, nil
}

func (r *FranchiseJourney) CancelAvailability(ctx context.Context, tenant, organization, entry string, version int64, subject, reason, eventID string) (franchisejourney.AvailabilityEntry, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return franchisejourney.AvailabilityEntry{}, err
	}
	defer tx.Rollback(ctx)
	var value franchisejourney.AvailabilityEntry
	if err = lockOrganizationSchedule(ctx, tx, tenant, organization); err != nil {
		return value, err
	}
	err = tx.QueryRow(ctx, `update crm.availability_entry set state='cancelled',version=version+1,cancelled_by_subject=$5,cancellation_reason_code=$6,cancelled_at=clock_timestamp() where tenant_id=$1 and organization_id=$2 and availability_id=$3 and version=$4 and state='active' returning availability_id,organization_id,coalesce(resource_id,''),entry_type,coalesce(reason_code,''),starts_at,ends_at,state,version`, tenant, organization, entry, version, subject, reason).Scan(&value.ID, &value.OrganizationID, &value.ResourceID, &value.EntryType, &value.ReasonCode, &value.StartsAt, &value.EndsAt, &value.State, &value.Version)
	if errors.Is(err, pgx.ErrNoRows) || postgresConflict(err) {
		return value, franchisejourney.ErrConflict
	}
	if err != nil {
		return value, err
	}
	if err = writeJourneyOutbox(ctx, tx, tenant, eventID, "availability-entry", entry, value.Version, "availability-entry.cancelled", map[string]any{"organization_id": organization, "actor_subject": subject, "reason_code": reason}); err != nil {
		return value, err
	}
	return value, tx.Commit(ctx)
}

func lockOrganizationSchedule(ctx context.Context, tx pgx.Tx, tenant, organization string) error {
	_, err := tx.Exec(ctx, `select pg_advisory_xact_lock(hashtextextended(($1::uuid)::text||'|availability|'||$2::text||'|organization',0))`, tenant, organization)
	return err
}

func (r *FranchiseJourney) Availability(ctx context.Context, tenant, organization, resource string, from, to time.Time) ([]franchisejourney.AvailabilityEntry, error) {
	rows, err := r.pool.Query(ctx, `select availability_id,organization_id,coalesce(resource_id,''),entry_type,coalesce(reason_code,''),starts_at,ends_at,state,version from crm.availability_entry where tenant_id=$1 and organization_id=$2 and ($3='' or resource_id=$3) and starts_at<$5 and ends_at>$4 order by starts_at,availability_id limit 1000`, tenant, organization, resource, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []franchisejourney.AvailabilityEntry{}
	for rows.Next() {
		var value franchisejourney.AvailabilityEntry
		if err = rows.Scan(&value.ID, &value.OrganizationID, &value.ResourceID, &value.EntryType, &value.ReasonCode, &value.StartsAt, &value.EndsAt, &value.State, &value.Version); err != nil {
			return nil, err
		}
		items = append(items, value)
	}
	return items, rows.Err()
}

func (r *FranchiseJourney) CancelCustomerAppointment(ctx context.Context, tenant, organization, customer, appointment string, version int64, reason, transitionID, eventID string) (franchisejourney.Appointment, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return franchisejourney.Appointment{}, err
	}
	defer tx.Rollback(ctx)
	var value franchisejourney.Appointment
	var current string
	err = tx.QueryRow(ctx, `with candidate as (select state from crm.appointment where tenant_id=$1 and organization_id=$2 and customer_principal_id=$3 and appointment_id=$4 and version=$5 and state in ('requested','confirmed') and starts_at>clock_timestamp() for update), changed as (update crm.appointment a set state='cancelled',version=version+1,updated_at=clock_timestamp() from candidate c where a.tenant_id=$1 and a.appointment_id=$4 returning a.appointment_id,a.organization_id,a.lead_id,coalesce(a.model_id,''),a.appointment_kind,a.starts_at,a.state,a.version,coalesce(a.slot_id,''),coalesce(a.ends_at,a.starts_at),c.state) select * from changed`, tenant, organization, customer, appointment, version).Scan(&value.ID, &value.OrganizationID, &value.LeadID, &value.ModelID, &value.Kind, &value.StartsAt, &value.State, &value.Version, &value.SlotID, &value.EndsAt, &current)
	if errors.Is(err, pgx.ErrNoRows) {
		return value, franchisejourney.ErrConflict
	}
	if err != nil {
		return value, err
	}
	if _, err = tx.Exec(ctx, `insert into crm.appointment_transition(tenant_id,transition_id,appointment_id,from_state,to_state,actor_subject,reason_code) values($1,$2,$3,$4,'cancelled',$5,$6)`, tenant, transitionID, appointment, current, customer, reason); err != nil {
		return value, err
	}
	if err = writeJourneyOutbox(ctx, tx, tenant, eventID, "appointment", appointment, value.Version, "appointment.cancelled", map[string]any{"organization_id": organization, "actor_subject": customer, "reason_code": reason, "channel": "customer"}); err != nil {
		return value, err
	}
	return value, tx.Commit(ctx)
}

func (r *FranchiseJourney) CreateAvailabilityOnce(ctx context.Context, tenant, subject, key, hash string, value franchisejourney.AvailabilityEntry, event string) (franchisejourney.AvailabilityEntry, bool, error) {
	var empty franchisejourney.AvailabilityEntry
	if subject == "" {
		return empty, false, franchisejourney.ErrInvalid
	}
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return empty, false, err
	}
	defer tx.Rollback(ctx)
	result, err := tx.Exec(ctx, `insert into platform.idempotency_record(tenant_id,scope,idempotency_key,request_sha256_hex,status,locked_until,expires_at)values($1,'franchise-availability',$2,$3,'processing',clock_timestamp()+interval '30 seconds',clock_timestamp()+interval '24 hours')on conflict do nothing`, tenant, key, hash)
	if err != nil {
		return empty, false, err
	}
	if result.RowsAffected() == 0 {
		replay, storedHash, err := readAvailabilityCreation(ctx, tx, tenant, value.OrganizationID, key)
		if err != nil {
			return empty, false, err
		}
		if hash != storedHash {
			return empty, false, franchisejourney.ErrConflict
		}
		if err = tx.Commit(ctx); err != nil {
			return empty, false, err
		}
		return replay, true, nil
	}
	value, err = r.createAvailabilityTx(ctx, tx, tenant, subject, value, event)
	if err != nil {
		return empty, false, err
	}
	result, err = tx.Exec(ctx, `update platform.idempotency_record set status='completed',response_code=201,response_body=jsonb_build_object('availability_id',$3::text),resource_type='availability-entry',resource_id=$3,locked_until=null where tenant_id=$1 and scope='franchise-availability' and idempotency_key=$2 and status='processing'`, tenant, key, value.ID)
	if err != nil {
		return empty, false, err
	}
	if result.RowsAffected() != 1 {
		return empty, false, franchisejourney.ErrConflict
	}
	if err = tx.Commit(ctx); err != nil {
		return empty, false, err
	}
	return value, false, nil
}

type availabilityRowReader interface {
	QueryRow(context.Context, string, ...any) pgx.Row
}

func readAvailabilityCreation(ctx context.Context, q availabilityRowReader, tenant, organization, key string) (franchisejourney.AvailabilityEntry, string, error) {
	var v franchisejourney.AvailabilityEntry
	var hash string
	err := q.QueryRow(ctx, `select i.request_sha256_hex,a.availability_id,a.organization_id,coalesce(a.resource_id,''),a.entry_type,coalesce(a.reason_code,''),a.starts_at,a.ends_at,a.state,a.version
 from platform.idempotency_record i join crm.availability_entry a on a.tenant_id=i.tenant_id and a.availability_id=i.resource_id
 where i.tenant_id=$1 and i.scope='franchise-availability' and i.idempotency_key=$3 and i.status='completed' and i.resource_type='availability-entry' and i.response_code=201 and i.response_body->>'availability_id'=a.availability_id and a.organization_id=$2`, tenant, organization, key).Scan(&hash, &v.ID, &v.OrganizationID, &v.ResourceID, &v.EntryType, &v.ReasonCode, &v.StartsAt, &v.EndsAt, &v.State, &v.Version)
	if errors.Is(err, pgx.ErrNoRows) {
		return franchisejourney.AvailabilityEntry{}, "", franchisejourney.ErrConflict
	}
	if err != nil {
		return franchisejourney.AvailabilityEntry{}, "", err
	}
	return v, hash, nil
}
func (r *FranchiseJourney) AvailabilityCreationResult(ctx context.Context, tenant, organization, key string) (franchisejourney.AvailabilityEntry, error) {
	v, _, err := readAvailabilityCreation(ctx, r.pool, tenant, organization, key)
	return v, err
}
````

### FILE: `db/migrations/0015_franchise_availability_and_appointment_audit.up.sql`

```yaml
block_id: "GO-FRANCHISE-JOURNEY:file:22"
operation: CREATE
provenance: AUTHORED
source: "local fail-closed schema governed by Microsoft Business Central calendars/absence and PostgreSQL 18 locking references pinned in metadata"
license: "LicenseRef-Workspace-Owner"
sha256: "027b2a51151e980a9ca3dff18ed60b2e73c697534dbacf7177406e7b0a720e1b"
variables: []
secrets_allowed: false
```

````sql
begin;

alter table crm.service_resource add constraint service_resource_org_identity_uq unique (tenant_id,organization_id,resource_id);

create table crm.availability_entry (
  tenant_id uuid not null,
  availability_id text not null,
  organization_id text not null,
  resource_id text,
  entry_type text not null check (entry_type in ('working','unavailable')),
  reason_code text,
  starts_at timestamptz not null,
  ends_at timestamptz not null,
  state text not null check (state in ('active','cancelled')),
  version bigint not null check (version>0),
  created_by_subject text not null check (length(created_by_subject) between 1 and 255),
  created_at timestamptz not null default clock_timestamp(),
  cancelled_by_subject text,
  cancellation_reason_code text,
  cancelled_at timestamptz,
  primary key (tenant_id,availability_id),
  foreign key (tenant_id,organization_id) references org.organization(tenant_id,organization_id),
  foreign key (tenant_id,organization_id,resource_id) references crm.service_resource(tenant_id,organization_id,resource_id),
  check (ends_at>starts_at and ends_at-starts_at<=interval '366 days'),
  check ((entry_type='working' and reason_code is null) or (entry_type='unavailable' and reason_code ~ '^[a-z][a-z0-9]*(-[a-z0-9]+)*$')),
  check ((state='active' and cancelled_by_subject is null and cancellation_reason_code is null and cancelled_at is null) or (state='cancelled' and cancelled_by_subject is not null and cancellation_reason_code is not null and cancelled_at is not null))
);

create index availability_scope_time_idx on crm.availability_entry(tenant_id,organization_id,resource_id,starts_at,ends_at) where state='active';

create or replace function crm.enforce_availability_entry() returns trigger language plpgsql as $$
begin
  perform pg_advisory_xact_lock(hashtextextended(new.tenant_id::text||'|availability|'||new.organization_id||'|'||coalesce(new.resource_id,'organization'),0));
  if exists (
    select 1 from crm.availability_entry e where e.tenant_id=new.tenant_id and e.organization_id=new.organization_id
      and e.resource_id is not distinct from new.resource_id and e.entry_type=new.entry_type and e.state='active'
      and e.availability_id<>new.availability_id and tstzrange(e.starts_at,e.ends_at,'[)') && tstzrange(new.starts_at,new.ends_at,'[)')
  ) then raise exception using errcode='23P01',message='availability entries overlap'; end if;
  if new.entry_type='unavailable' and exists (
    select 1 from crm.appointment a left join crm.appointment_resource ar on ar.tenant_id=a.tenant_id and ar.appointment_id=a.appointment_id
    where a.tenant_id=new.tenant_id and a.organization_id=new.organization_id and a.state in ('requested','confirmed')
      and (new.resource_id is null or ar.resource_id=new.resource_id)
      and tstzrange(a.starts_at,coalesce(a.ends_at,a.starts_at+interval '1 hour'),'[)') && tstzrange(new.starts_at,new.ends_at,'[)')
  ) then raise exception using errcode='23P01',message='unavailability conflicts with active appointments'; end if;
  return new;
end $$;
create trigger availability_entry_guard before insert or update of organization_id,resource_id,entry_type,starts_at,ends_at,state on crm.availability_entry for each row execute function crm.enforce_availability_entry();

create or replace function crm.enforce_availability_cancellation() returns trigger language plpgsql as $$
begin
  if old.state='active' and new.state='cancelled' and old.entry_type='working' and exists (
    select 1 from crm.appointment a left join crm.appointment_resource ar on ar.tenant_id=a.tenant_id and ar.appointment_id=a.appointment_id
    where a.tenant_id=old.tenant_id and a.organization_id=old.organization_id and a.state in ('requested','confirmed')
      and (old.resource_id is null or ar.resource_id=old.resource_id)
      and tstzrange(a.starts_at,coalesce(a.ends_at,a.starts_at+interval '1 hour'),'[)') && tstzrange(old.starts_at,old.ends_at,'[)')
  ) then raise exception using errcode='23P01',message='working availability has active appointments'; end if;
  return new;
end $$;
create trigger availability_cancellation_guard before update of state on crm.availability_entry for each row execute function crm.enforce_availability_cancellation();

create or replace function crm.enforce_appointment_slot_availability() returns trigger language plpgsql as $$
begin
  if new.state='open' and (not exists(select 1 from crm.availability_entry e where e.tenant_id=new.tenant_id and e.organization_id=new.organization_id and e.resource_id is null and e.entry_type='working' and e.state='active' and e.starts_at<=new.starts_at and e.ends_at>=new.ends_at)
    or exists(select 1 from crm.availability_entry e where e.tenant_id=new.tenant_id and e.organization_id=new.organization_id and e.resource_id is null and e.entry_type='unavailable' and e.state='active' and tstzrange(e.starts_at,e.ends_at,'[)') && tstzrange(new.starts_at,new.ends_at,'[)'))) then
    raise exception using errcode='23514',message='appointment slot is outside organization availability';
  end if;
  return new;
end $$;
create trigger appointment_slot_availability_guard before insert or update of organization_id,starts_at,ends_at,state on crm.appointment_slot for each row execute function crm.enforce_appointment_slot_availability();

create or replace function crm.enforce_appointment_resource_assignment() returns trigger language plpgsql as $$
declare candidate crm.appointment%rowtype;
begin
  perform pg_advisory_xact_lock(hashtextextended(new.tenant_id::text || '|appointment-resource|' || new.resource_id,0));
  select * into candidate from crm.appointment a where a.tenant_id=new.tenant_id and a.appointment_id=new.appointment_id;
  if candidate.appointment_id is null or candidate.state <> 'requested' then raise exception using errcode='23514',message='appointment is unavailable for resource assignment'; end if;
  if not exists(select 1 from crm.service_resource r join crm.resource_skill s on s.tenant_id=r.tenant_id and s.resource_id=r.resource_id where r.tenant_id=new.tenant_id and r.resource_id=new.resource_id and r.organization_id=candidate.organization_id and r.status='active' and s.appointment_kind=candidate.appointment_kind) then raise exception using errcode='23514',message='resource is unavailable or lacks appointment skill'; end if;
  if not exists(select 1 from crm.availability_entry e where e.tenant_id=new.tenant_id and e.organization_id=candidate.organization_id and e.resource_id=new.resource_id and e.entry_type='working' and e.state='active' and e.starts_at<=candidate.starts_at and e.ends_at>=coalesce(candidate.ends_at,candidate.starts_at+interval '1 hour'))
    or exists(select 1 from crm.availability_entry e where e.tenant_id=new.tenant_id and e.resource_id=new.resource_id and e.entry_type='unavailable' and e.state='active' and tstzrange(e.starts_at,e.ends_at,'[)') && tstzrange(candidate.starts_at,coalesce(candidate.ends_at,candidate.starts_at+interval '1 hour'),'[)')) then raise exception using errcode='23514',message='resource is outside working availability'; end if;
  if exists(select 1 from crm.appointment_resource ar join crm.appointment a on a.tenant_id=ar.tenant_id and a.appointment_id=ar.appointment_id where ar.tenant_id=new.tenant_id and ar.resource_id=new.resource_id and ar.appointment_id<>new.appointment_id and a.state in ('requested','confirmed') and tstzrange(a.starts_at,coalesce(a.ends_at,a.starts_at+interval '1 hour'),'[)') && tstzrange(candidate.starts_at,coalesce(candidate.ends_at,candidate.starts_at+interval '1 hour'),'[)')) then raise exception using errcode='23P01',message='resource appointment assignments overlap'; end if;
  return new;
end $$;

create table crm.appointment_transition (
  tenant_id uuid not null,
  transition_id uuid not null,
  appointment_id text not null,
  from_state text not null,
  to_state text not null,
  actor_subject text not null check (length(actor_subject) between 1 and 255),
  reason_code text,
  occurred_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id,transition_id),
  foreign key (tenant_id,appointment_id) references crm.appointment(tenant_id,appointment_id),
  check ((to_state in ('cancelled','no-show') and reason_code ~ '^[a-z][a-z0-9]*(-[a-z0-9]+)*$') or (to_state not in ('cancelled','no-show') and reason_code is null))
);
create index appointment_transition_history_idx on crm.appointment_transition(tenant_id,appointment_id,occurred_at,transition_id);
create or replace function crm.prevent_appointment_transition_mutation() returns trigger language plpgsql as $$ begin raise exception 'immutable appointment transition'; end $$;
create trigger appointment_transition_immutable before update or delete on crm.appointment_transition for each row execute function crm.prevent_appointment_transition_mutation();

commit;
````

### FILE: `db/migrations/0015_franchise_availability_and_appointment_audit.down.sql`

```yaml
block_id: "GO-FRANCHISE-JOURNEY:file:23"
operation: CREATE
provenance: AUTHORED
source: "local reversible migration governed by Microsoft Business Central calendars/absence and PostgreSQL 18 references pinned in metadata"
license: "LicenseRef-Workspace-Owner"
sha256: "eecab84a4b5652d3d753cf0ed80f703bebed40be941e88240b7877b316b6eac4"
variables: []
secrets_allowed: false
```

````sql
begin;
drop trigger if exists appointment_transition_immutable on crm.appointment_transition;
drop function if exists crm.prevent_appointment_transition_mutation();
drop table if exists crm.appointment_transition;
drop trigger if exists appointment_slot_availability_guard on crm.appointment_slot;
drop function if exists crm.enforce_appointment_slot_availability();
drop trigger if exists availability_cancellation_guard on crm.availability_entry;
drop function if exists crm.enforce_availability_cancellation();
drop trigger if exists availability_entry_guard on crm.availability_entry;
drop function if exists crm.enforce_availability_entry();
drop table if exists crm.availability_entry;
alter table crm.service_resource drop constraint if exists service_resource_org_identity_uq;
-- Restore the assignment guard owned by migration 0009.
create or replace function crm.enforce_appointment_resource_assignment() returns trigger language plpgsql as $$
declare candidate crm.appointment%rowtype;
begin
  perform pg_advisory_xact_lock(hashtextextended(new.tenant_id::text || '|appointment-resource|' || new.resource_id,0));
  select * into candidate from crm.appointment a where a.tenant_id=new.tenant_id and a.appointment_id=new.appointment_id;
  if candidate.appointment_id is null or candidate.state <> 'requested' then raise exception using errcode='23514',message='appointment is unavailable for resource assignment'; end if;
  if not exists(select 1 from crm.service_resource r join crm.resource_skill s on s.tenant_id=r.tenant_id and s.resource_id=r.resource_id where r.tenant_id=new.tenant_id and r.resource_id=new.resource_id and r.organization_id=candidate.organization_id and r.status='active' and s.appointment_kind=candidate.appointment_kind) then raise exception using errcode='23514',message='resource is unavailable or lacks appointment skill'; end if;
  if exists(select 1 from crm.appointment_resource ar join crm.appointment a on a.tenant_id=ar.tenant_id and a.appointment_id=ar.appointment_id where ar.tenant_id=new.tenant_id and ar.resource_id=new.resource_id and ar.appointment_id<>new.appointment_id and a.state in ('requested','confirmed') and tstzrange(a.starts_at,coalesce(a.ends_at,a.starts_at+interval '1 hour'),'[)') && tstzrange(candidate.starts_at,coalesce(candidate.ends_at,candidate.starts_at+interval '1 hour'),'[)')) then raise exception using errcode='23P01',message='resource appointment assignments overlap'; end if;
  return new;
end $$;
commit;
````

### FILE: `db/tests/0015_franchise_availability_and_appointment_audit.test.sql`

```yaml
block_id: "GO-FRANCHISE-JOURNEY:file:24"
operation: CREATE
provenance: AUTHORED
source: "local invariant test governed by Microsoft Business Central calendars/absence and PostgreSQL 18 references pinned in metadata"
license: "LicenseRef-Workspace-Owner"
sha256: "e376d555a9cbed15cf1f8417f883225e209934c84e890af396173da07380a4e7"
variables: []
secrets_allowed: false
```

````sql
begin;
do $$ begin
  if to_regclass('crm.availability_entry') is null then raise exception 'availability_entry missing'; end if;
  if to_regclass('crm.appointment_transition') is null then raise exception 'appointment_transition missing'; end if;
  if to_regprocedure('crm.enforce_appointment_slot_availability()') is null then raise exception 'slot availability guard missing'; end if;
  if to_regprocedure('crm.enforce_appointment_resource_assignment()') is null then raise exception 'resource availability guard missing'; end if;
  if not exists(select 1 from pg_trigger where tgname='appointment_transition_immutable') then raise exception 'transition immutability missing'; end if;
end $$;
rollback;
````

### FILE: `db/migrations/0016_versioned_delivery_checklist.up.sql`

```yaml
block_id: "GO-FRANCHISE-JOURNEY:file:25"
operation: CREATE
provenance: AUTHORED
source: "local delivery checklist contract governed by Microsoft Field Service inspections and Business Central posting references pinned in metadata"
license: "LicenseRef-Workspace-Owner"
sha256: "67fcc7ac5e5e11bc57f648f1bd5cf87d2e2da174e4be07eae63820fea953098e"
variables: []
secrets_allowed: false
```

````sql
begin;

create table sales.delivery_checklist_template (
  tenant_id uuid not null,
  organization_id text not null,
  checklist_id text not null,
  checklist_version bigint not null check (checklist_version > 0),
  title text not null check (length(title) between 1 and 160),
  state text not null check (state in ('draft','published','retired')),
  created_by_subject text not null check (length(created_by_subject) between 1 and 255),
  created_at timestamptz not null default clock_timestamp(),
  published_at timestamptz,
  primary key (tenant_id,organization_id,checklist_id,checklist_version),
  foreign key (tenant_id,organization_id) references org.organization(tenant_id,organization_id),
  check (checklist_id ~ '^[a-z][a-z0-9]*(-[a-z0-9]+)*$'),
  check ((state='draft' and published_at is null) or (state in ('published','retired') and published_at is not null))
);

create table sales.delivery_checklist_item (
  tenant_id uuid not null,
  organization_id text not null,
  checklist_id text not null,
  checklist_version bigint not null,
  item_id text not null,
  ordinal integer not null check (ordinal between 1 and 1000),
  prompt text not null check (length(prompt) between 1 and 500),
  response_type text not null check (response_type in ('confirmation','text','serial','evidence')),
  required boolean not null,
  primary key (tenant_id,organization_id,checklist_id,checklist_version,item_id),
  unique (tenant_id,organization_id,checklist_id,checklist_version,ordinal),
  foreign key (tenant_id,organization_id,checklist_id,checklist_version) references sales.delivery_checklist_template(tenant_id,organization_id,checklist_id,checklist_version),
  check (item_id ~ '^[a-z][a-z0-9]*(-[a-z0-9]+)*$')
);

alter table sales.delivery_handover
  add column checklist_id text,
  add column checklist_version bigint,
  add column checklist_completed_at timestamptz,
  add column checklist_completed_by_subject text,
  add constraint delivery_handover_checklist_fk foreign key (tenant_id,organization_id,checklist_id,checklist_version) references sales.delivery_checklist_template(tenant_id,organization_id,checklist_id,checklist_version),
  add constraint delivery_handover_checklist_state_ck check (
    (checklist_id is null and checklist_version is null and checklist_completed_at is null and checklist_completed_by_subject is null)
    or
    (checklist_id is not null and checklist_version > 0 and checklist_completed_at is not null and length(checklist_completed_by_subject) between 1 and 255)
  );

create table sales.delivery_checklist_response (
  tenant_id uuid not null,
  organization_id text not null,
  handover_id text not null,
  checklist_id text not null,
  checklist_version bigint not null,
  item_id text not null,
  response_text text not null check (length(response_text) between 1 and 2048),
  evidence_sha256_hex text,
  answered_by_subject text not null check (length(answered_by_subject) between 1 and 255),
  answered_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id,handover_id,item_id),
  foreign key (tenant_id,handover_id) references sales.delivery_handover(tenant_id,handover_id),
  foreign key (tenant_id,organization_id,checklist_id,checklist_version,item_id) references sales.delivery_checklist_item(tenant_id,organization_id,checklist_id,checklist_version,item_id),
  check (evidence_sha256_hex is null or evidence_sha256_hex ~ '^[0-9a-f]{64}$')
);

create or replace function sales.prevent_published_delivery_checklist_mutation() returns trigger language plpgsql as $$
declare parent_state text;
begin
  if tg_table_name='delivery_checklist_template' then
    if old.state in ('published','retired') then raise exception using errcode='23514',message='published delivery checklist is immutable'; end if;
    if tg_op='DELETE' then return old; end if;
    return new;
  end if;
  select state into parent_state from sales.delivery_checklist_template
    where tenant_id=coalesce(new.tenant_id,old.tenant_id) and organization_id=coalesce(new.organization_id,old.organization_id)
      and checklist_id=coalesce(new.checklist_id,old.checklist_id) and checklist_version=coalesce(new.checklist_version,old.checklist_version);
  if parent_state in ('published','retired') then raise exception using errcode='23514',message='published delivery checklist items are immutable'; end if;
  if tg_op='DELETE' then return old; end if;
  return new;
end $$;
create trigger delivery_checklist_template_immutable before update or delete on sales.delivery_checklist_template for each row execute function sales.prevent_published_delivery_checklist_mutation();
create trigger delivery_checklist_item_immutable before insert or update or delete on sales.delivery_checklist_item for each row execute function sales.prevent_published_delivery_checklist_mutation();

create or replace function sales.enforce_delivery_checklist_response() returns trigger language plpgsql as $$
declare expected_type text;
declare completed_at timestamptz;
begin
  select h.checklist_completed_at into completed_at from sales.delivery_handover h where h.tenant_id=new.tenant_id and h.handover_id=new.handover_id for update;
  if completed_at is not null then raise exception using errcode='23514',message='completed delivery checklist responses are immutable'; end if;
  select i.response_type into expected_type from sales.delivery_checklist_item i
    where i.tenant_id=new.tenant_id and i.organization_id=new.organization_id and i.checklist_id=new.checklist_id and i.checklist_version=new.checklist_version and i.item_id=new.item_id;
  if expected_type='confirmation' and new.response_text<>'confirmed' then raise exception using errcode='23514',message='confirmation response must be confirmed'; end if;
  if expected_type='evidence' and new.evidence_sha256_hex is null then raise exception using errcode='23514',message='evidence response requires digest'; end if;
  return new;
end $$;
create trigger delivery_checklist_response_guard before insert on sales.delivery_checklist_response for each row execute function sales.enforce_delivery_checklist_response();

create or replace function sales.prevent_delivery_checklist_response_mutation() returns trigger language plpgsql as $$ begin raise exception using errcode='23514',message='delivery checklist response is immutable'; end $$;
create trigger delivery_checklist_response_immutable before update or delete on sales.delivery_checklist_response for each row execute function sales.prevent_delivery_checklist_response_mutation();

create or replace function sales.enforce_delivery_checklist_completion() returns trigger language plpgsql as $$
begin
  if new.checklist_completed_at is not null and old.checklist_completed_at is null then
    if new.state<>'presented' then raise exception using errcode='23514',message='completed delivery checklist must be presented'; end if;
    if not exists (select 1 from sales.delivery_checklist_template t where t.tenant_id=new.tenant_id and t.organization_id=new.organization_id and t.checklist_id=new.checklist_id and t.checklist_version=new.checklist_version and t.state='published') then
      raise exception using errcode='23514',message='delivery checklist version is not published';
    end if;
    if exists (
      select 1 from sales.delivery_checklist_item i
      left join sales.delivery_checklist_response r on r.tenant_id=i.tenant_id and r.organization_id=i.organization_id and r.checklist_id=i.checklist_id and r.checklist_version=i.checklist_version and r.item_id=i.item_id and r.handover_id=new.handover_id
      where i.tenant_id=new.tenant_id and i.organization_id=new.organization_id and i.checklist_id=new.checklist_id and i.checklist_version=new.checklist_version and i.required and r.item_id is null
    ) then raise exception using errcode='23514',message='required delivery checklist response is missing'; end if;
  elsif old.checklist_completed_at is not null and (new.checklist_id,new.checklist_version,new.checklist_completed_at,new.checklist_completed_by_subject) is distinct from (old.checklist_id,old.checklist_version,old.checklist_completed_at,old.checklist_completed_by_subject) then
    raise exception using errcode='23514',message='completed delivery checklist binding is immutable';
  end if;
  return new;
end $$;
create trigger delivery_checklist_completion_guard before update of state,checklist_id,checklist_version,checklist_completed_at,checklist_completed_by_subject on sales.delivery_handover for each row execute function sales.enforce_delivery_checklist_completion();

create index delivery_checklist_template_state_idx on sales.delivery_checklist_template(tenant_id,organization_id,state,checklist_id,checklist_version desc);
create index delivery_checklist_response_handover_idx on sales.delivery_checklist_response(tenant_id,handover_id,answered_at,item_id);

commit;
````

### FILE: `db/migrations/0016_versioned_delivery_checklist.down.sql`

```yaml
block_id: "GO-FRANCHISE-JOURNEY:file:26"
operation: CREATE
provenance: AUTHORED
source: "local reversible migration governed by Microsoft Field Service inspections and Business Central posting references pinned in metadata"
license: "LicenseRef-Workspace-Owner"
sha256: "c31137a931f72b765b1585863cb08257bf6f77a1c1e6c722e6a7a87aabfe3d15"
variables: []
secrets_allowed: false
```

````sql
begin;
drop trigger if exists delivery_checklist_completion_guard on sales.delivery_handover;
drop function if exists sales.enforce_delivery_checklist_completion();
drop trigger if exists delivery_checklist_response_immutable on sales.delivery_checklist_response;
drop function if exists sales.prevent_delivery_checklist_response_mutation();
drop trigger if exists delivery_checklist_response_guard on sales.delivery_checklist_response;
drop function if exists sales.enforce_delivery_checklist_response();
drop table if exists sales.delivery_checklist_response;
alter table sales.delivery_handover
  drop constraint if exists delivery_handover_checklist_state_ck,
  drop constraint if exists delivery_handover_checklist_fk,
  drop column if exists checklist_completed_by_subject,
  drop column if exists checklist_completed_at,
  drop column if exists checklist_version,
  drop column if exists checklist_id;
drop trigger if exists delivery_checklist_item_immutable on sales.delivery_checklist_item;
drop table if exists sales.delivery_checklist_item;
drop trigger if exists delivery_checklist_template_immutable on sales.delivery_checklist_template;
drop function if exists sales.prevent_published_delivery_checklist_mutation();
drop table if exists sales.delivery_checklist_template;
commit;
````

### FILE: `db/tests/0016_versioned_delivery_checklist.test.sql`

```yaml
block_id: "GO-FRANCHISE-JOURNEY:file:27"
operation: CREATE
provenance: AUTHORED
source: "local invariant test governed by Microsoft Field Service inspections and PostgreSQL 18 references pinned in metadata"
license: "LicenseRef-Workspace-Owner"
sha256: "13181c3d6eb1daef8d913c3cb178ee55d76bfaa1e20e699054c3eb40f6a4e550"
variables: []
secrets_allowed: false
```

````sql
begin;
do $$ begin
  if to_regclass('sales.delivery_checklist_template') is null then raise exception 'delivery checklist template missing'; end if;
  if to_regclass('sales.delivery_checklist_item') is null then raise exception 'delivery checklist item missing'; end if;
  if to_regclass('sales.delivery_checklist_response') is null then raise exception 'delivery checklist response missing'; end if;
  if not exists(select 1 from pg_trigger where tgname='delivery_checklist_template_immutable') then raise exception 'published checklist immutability missing'; end if;
  if not exists(select 1 from pg_trigger where tgname='delivery_checklist_response_immutable') then raise exception 'response immutability missing'; end if;
  if not exists(select 1 from pg_trigger where tgname='delivery_checklist_completion_guard') then raise exception 'completion guard missing'; end if;
  if not exists(select 1 from information_schema.columns where table_schema='sales' and table_name='delivery_handover' and column_name='checklist_version') then raise exception 'handover checklist binding missing'; end if;
end $$;
rollback;
````

### FILE: `db/migrations/0017_delivery_exception_and_return_authorization.up.sql`

```yaml
block_id: "GO-FRANCHISE-JOURNEY:file:28"
operation: CREATE
provenance: AUTHORED
source: "local exception/return boundary governed by https://learn.microsoft.com/en-us/dynamics365/business-central/sales-how-process-sales-returns-cancellations and pinned Microsoft BCApps references"
license: "LicenseRef-Workspace-Owner"
sha256: "b7c0c066915a649eed82ea811ac87fe90b54ca1ad4f7166f4235979ac84a82fc"
variables: []
secrets_allowed: false
```

````sql
begin;

alter table sales.delivery_handover add column supersedes_handover_id text;
alter table sales.delivery_handover add constraint delivery_handover_supersedes_fk foreign key (tenant_id,supersedes_handover_id) references sales.delivery_handover(tenant_id,handover_id);
alter table sales.delivery_handover add constraint delivery_handover_not_self_superseding_ck check (supersedes_handover_id is null or supersedes_handover_id<>handover_id);
create unique index delivery_handover_one_successor_uq on sales.delivery_handover(tenant_id,supersedes_handover_id) where supersedes_handover_id is not null;

create table sales.delivery_exception (
  tenant_id uuid not null,
  exception_id text not null,
  organization_id text not null,
  handover_id text not null,
  customer_principal_id text not null,
  reason_code text not null,
  details text not null check (length(details) between 1 and 1000),
  rejection_evidence_sha256_hex text not null check (rejection_evidence_sha256_hex ~ '^[0-9a-f]{64}$'),
  state text not null check (state in ('open','resolved')),
  version bigint not null check (version>0),
  created_at timestamptz not null default clock_timestamp(),
  resolved_at timestamptz,
  resolved_by_subject text,
  primary key (tenant_id,exception_id),
  unique (tenant_id,handover_id),
  foreign key (tenant_id,handover_id) references sales.delivery_handover(tenant_id,handover_id),
  foreign key (tenant_id,organization_id) references org.organization(tenant_id,organization_id),
  foreign key (tenant_id,customer_principal_id) references crm.customer_profile(tenant_id,customer_principal_id),
  check (reason_code ~ '^[a-z][a-z0-9]*(-[a-z0-9]+)*$'),
  check ((state='open' and resolved_at is null and resolved_by_subject is null) or (state='resolved' and resolved_at is not null and length(resolved_by_subject) between 1 and 255))
);

create table sales.return_authorization (
  tenant_id uuid not null,
  authorization_id text not null,
  organization_id text not null,
  exception_id text not null,
  handover_id text not null,
  order_id text not null,
  stock_unit_id text not null,
  customer_principal_id text not null,
  disposition text not null check (disposition in ('return','exchange')),
  exact_cost_source_order_id text not null,
  state text not null check (state='authorized'),
  authorized_by_subject text not null check (length(authorized_by_subject) between 1 and 255),
  authorized_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id,authorization_id),
  unique (tenant_id,exception_id),
  foreign key (tenant_id,exception_id) references sales.delivery_exception(tenant_id,exception_id),
  foreign key (tenant_id,handover_id) references sales.delivery_handover(tenant_id,handover_id),
  foreign key (tenant_id,order_id) references sales.customer_order(tenant_id,order_id),
  foreign key (tenant_id,stock_unit_id) references inventory.stock_unit(tenant_id,stock_unit_id),
  foreign key (tenant_id,customer_principal_id) references crm.customer_profile(tenant_id,customer_principal_id),
  check (exact_cost_source_order_id=order_id)
);

create table sales.delivery_exception_resolution (
  tenant_id uuid not null,
  resolution_id text not null,
  exception_id text not null,
  action text not null check (action in ('correct-and-represent','return','exchange')),
  notes text not null check (length(notes) between 1 and 1000),
  resolved_by_subject text not null check (length(resolved_by_subject) between 1 and 255),
  resolved_at timestamptz not null default clock_timestamp(),
  successor_handover_id text,
  return_authorization_id text,
  primary key (tenant_id,resolution_id),
  unique (tenant_id,exception_id),
  foreign key (tenant_id,exception_id) references sales.delivery_exception(tenant_id,exception_id),
  foreign key (tenant_id,successor_handover_id) references sales.delivery_handover(tenant_id,handover_id),
  foreign key (tenant_id,return_authorization_id) references sales.return_authorization(tenant_id,authorization_id),
  check ((action='correct-and-represent' and successor_handover_id is not null and return_authorization_id is null) or (action in ('return','exchange') and successor_handover_id is null and return_authorization_id is not null))
);

create or replace function sales.enforce_delivery_exception_lifecycle() returns trigger language plpgsql as $$
begin
  if tg_op='DELETE' then raise exception using errcode='23514',message='delivery exception is immutable'; end if;
  if old.state<>'open' or new.state<>'resolved' or (new.tenant_id,new.exception_id,new.organization_id,new.handover_id,new.customer_principal_id,new.reason_code,new.details,new.rejection_evidence_sha256_hex,new.created_at) is distinct from (old.tenant_id,old.exception_id,old.organization_id,old.handover_id,old.customer_principal_id,old.reason_code,old.details,old.rejection_evidence_sha256_hex,old.created_at) or new.version<>old.version+1 then
    raise exception using errcode='23514',message='invalid delivery exception transition';
  end if;
  return new;
end $$;
create trigger delivery_exception_lifecycle before update or delete on sales.delivery_exception for each row execute function sales.enforce_delivery_exception_lifecycle();

create or replace function sales.prevent_delivery_resolution_mutation() returns trigger language plpgsql as $$ begin raise exception using errcode='23514',message='delivery resolution is immutable'; end $$;
create trigger delivery_exception_resolution_immutable before update or delete on sales.delivery_exception_resolution for each row execute function sales.prevent_delivery_resolution_mutation();
create trigger return_authorization_immutable before update or delete on sales.return_authorization for each row execute function sales.prevent_delivery_resolution_mutation();

create index delivery_exception_org_state_idx on sales.delivery_exception(tenant_id,organization_id,state,created_at,exception_id);
create index delivery_exception_customer_idx on sales.delivery_exception(tenant_id,customer_principal_id,created_at desc,exception_id);

commit;
````

### FILE: `db/migrations/0017_delivery_exception_and_return_authorization.down.sql`

```yaml
block_id: "GO-FRANCHISE-JOURNEY:file:29"
operation: CREATE
provenance: AUTHORED
source: "local reversible migration governed by https://learn.microsoft.com/en-us/dynamics365/business-central/sales-how-process-sales-returns-cancellations"
license: "LicenseRef-Workspace-Owner"
sha256: "bb486c2eb7f8a889fdaf036f162547594bfce9341263723f270ecf0dc28111b6"
variables: []
secrets_allowed: false
```

````sql
begin;
drop index if exists sales.delivery_exception_customer_idx;
drop index if exists sales.delivery_exception_org_state_idx;
drop trigger if exists return_authorization_immutable on sales.return_authorization;
drop trigger if exists delivery_exception_resolution_immutable on sales.delivery_exception_resolution;
drop function if exists sales.prevent_delivery_resolution_mutation();
drop trigger if exists delivery_exception_lifecycle on sales.delivery_exception;
drop function if exists sales.enforce_delivery_exception_lifecycle();
drop table if exists sales.delivery_exception_resolution;
drop table if exists sales.return_authorization;
drop table if exists sales.delivery_exception;
drop index if exists sales.delivery_handover_one_successor_uq;
alter table sales.delivery_handover drop constraint if exists delivery_handover_not_self_superseding_ck;
alter table sales.delivery_handover drop constraint if exists delivery_handover_supersedes_fk;
alter table sales.delivery_handover drop column if exists supersedes_handover_id;
commit;
````

### FILE: `db/tests/0017_delivery_exception_and_return_authorization.test.sql`

```yaml
block_id: "GO-FRANCHISE-JOURNEY:file:30"
operation: CREATE
provenance: AUTHORED
source: "local invariant test governed by Microsoft Business Central return-order separation and PostgreSQL 18 constraints"
license: "LicenseRef-Workspace-Owner"
sha256: "a7643f33eb07d288067467a2c61b1788fe28ed067c03135ae7ecfff22e0b6dc8"
variables: []
secrets_allowed: false
```

````sql
begin;
do $$ begin
  if to_regclass('sales.delivery_exception') is null then raise exception 'delivery exception missing'; end if;
  if to_regclass('sales.delivery_exception_resolution') is null then raise exception 'delivery exception resolution missing'; end if;
  if to_regclass('sales.return_authorization') is null then raise exception 'return authorization missing'; end if;
  if not exists(select 1 from pg_trigger where tgname='delivery_exception_lifecycle') then raise exception 'delivery exception lifecycle guard missing'; end if;
  if not exists(select 1 from pg_trigger where tgname='delivery_exception_resolution_immutable') then raise exception 'resolution immutability missing'; end if;
  if not exists(select 1 from pg_indexes where schemaname='sales' and indexname='delivery_handover_one_successor_uq') then raise exception 'one successor invariant missing'; end if;
end $$;
rollback;
````

### FILE: `db/migrations/0018_return_receipt_disposition_effects.up.sql`

```yaml
block_id: "GO-FRANCHISE-JOURNEY:file:31"
operation: CREATE
provenance: AUTHORED
source: "local return receipt/disposition implementation governed by Microsoft Business Central sales-return and public BCApps return-receipt references"
license: "LicenseRef-Workspace-Owner"
sha256: "db8b62aeab5d3bc23881ce78f1edb6df35dff6c617a6d4223feb8faef6d8ff4a"
variables: []
secrets_allowed: false
```

````sql
begin;

create table sales.return_receipt (
  tenant_id uuid not null,
  receipt_id text not null,
  authorization_id text not null,
  organization_id text not null,
  order_id text not null,
  stock_unit_id text not null,
  customer_principal_id text not null,
  received_serial_number text not null check (length(received_serial_number) between 1 and 128),
  condition_code text not null check (condition_code in ('sealed','opened','damaged','incomplete')),
  notes text not null check (length(notes) between 1 and 1000),
  evidence_sha256_hex text not null check (evidence_sha256_hex ~ '^[0-9a-f]{64}$'),
  received_by_subject text not null check (length(received_by_subject) between 1 and 255),
  received_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id,receipt_id),
  unique (tenant_id,authorization_id),
  foreign key (tenant_id,authorization_id) references sales.return_authorization(tenant_id,authorization_id),
  foreign key (tenant_id,organization_id) references org.organization(tenant_id,organization_id),
  foreign key (tenant_id,order_id) references sales.customer_order(tenant_id,order_id),
  foreign key (tenant_id,stock_unit_id) references inventory.stock_unit(tenant_id,stock_unit_id),
  foreign key (tenant_id,customer_principal_id) references crm.customer_profile(tenant_id,customer_principal_id)
);

create or replace function sales.validate_return_receipt_insert() returns trigger language plpgsql as $$
declare expected record;
begin
  select a.organization_id,a.order_id,a.stock_unit_id,a.customer_principal_id,s.serial_number
    into expected
    from sales.return_authorization a
    join inventory.stock_unit s on s.tenant_id=a.tenant_id and s.stock_unit_id=a.stock_unit_id
   where a.tenant_id=new.tenant_id and a.authorization_id=new.authorization_id and a.state='authorized';
  if not found or (new.organization_id,new.order_id,new.stock_unit_id,new.customer_principal_id,new.received_serial_number)
      is distinct from (expected.organization_id,expected.order_id,expected.stock_unit_id,expected.customer_principal_id,expected.serial_number) then
    raise exception using errcode='23514',message='return receipt does not match authorization and durable serial';
  end if;
  return new;
end $$;
create trigger return_receipt_validate before insert on sales.return_receipt for each row execute function sales.validate_return_receipt_insert();

create table sales.return_disposition (
  tenant_id uuid not null,
  disposition_id text not null,
  receipt_id text not null,
  inventory_action text not null check (inventory_action in ('quarantine','restock','repair','scrap')),
  customer_remedy text not null check (customer_remedy in ('refund','exchange')),
  notes text not null check (length(notes) between 1 and 1000),
  decided_by_subject text not null check (length(decided_by_subject) between 1 and 255),
  decided_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id,disposition_id),
  unique (tenant_id,receipt_id),
  foreign key (tenant_id,receipt_id) references sales.return_receipt(tenant_id,receipt_id)
);

create or replace function sales.validate_return_disposition_insert() returns trigger language plpgsql as $$
declare authorized_disposition text;
begin
  select a.disposition into authorized_disposition
    from sales.return_receipt r
    join sales.return_authorization a on a.tenant_id=r.tenant_id and a.authorization_id=r.authorization_id
   where r.tenant_id=new.tenant_id and r.receipt_id=new.receipt_id;
  if not found or (authorized_disposition='return' and new.customer_remedy<>'refund') or (authorized_disposition='exchange' and new.customer_remedy<>'exchange') then
    raise exception using errcode='23514',message='return disposition contradicts authorization';
  end if;
  return new;
end $$;
create trigger return_disposition_validate before insert on sales.return_disposition for each row execute function sales.validate_return_disposition_insert();

create table sales.return_effect_request (
  tenant_id uuid not null,
  request_id text not null,
  disposition_id text not null,
  effect_kind text not null check (effect_kind in ('inventory','refund','exchange','accounting','fiscal')),
  owner_context text not null check (owner_context in ('inventory','payment','fulfillment','accounting','fiscal')),
  state text not null check (state='requested'),
  idempotency_key text not null check (length(idempotency_key) between 16 and 128),
  requested_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id,request_id),
  unique (tenant_id,disposition_id,effect_kind),
  unique (tenant_id,idempotency_key),
  foreign key (tenant_id,disposition_id) references sales.return_disposition(tenant_id,disposition_id),
  check ((effect_kind='inventory' and owner_context='inventory') or (effect_kind='refund' and owner_context='payment') or (effect_kind='exchange' and owner_context='fulfillment') or (effect_kind='accounting' and owner_context='accounting') or (effect_kind='fiscal' and owner_context='fiscal'))
);

create or replace function sales.prevent_return_processing_mutation() returns trigger language plpgsql as $$
begin
  raise exception using errcode='23514',message='return processing evidence is immutable';
end $$;
create trigger return_receipt_immutable before update or delete on sales.return_receipt for each row execute function sales.prevent_return_processing_mutation();
create trigger return_disposition_immutable before update or delete on sales.return_disposition for each row execute function sales.prevent_return_processing_mutation();
create trigger return_effect_request_immutable before update or delete on sales.return_effect_request for each row execute function sales.prevent_return_processing_mutation();

create index return_receipt_org_time_idx on sales.return_receipt(tenant_id,organization_id,received_at desc,receipt_id);
create index return_effect_request_state_idx on sales.return_effect_request(tenant_id,owner_context,state,requested_at,request_id);

commit;
````

### FILE: `db/migrations/0018_return_receipt_disposition_effects.down.sql`

```yaml
block_id: "GO-FRANCHISE-JOURNEY:file:32"
operation: CREATE
provenance: AUTHORED
source: "local reversible migration governed by PostgreSQL transactional DDL"
license: "LicenseRef-Workspace-Owner"
sha256: "a62f62aafb65fa0c5f6ed2047b972ee260d1c7b5e12c4348eeaf61a479673b19"
variables: []
secrets_allowed: false
```

````sql
begin;

drop index if exists sales.return_effect_request_state_idx;
drop index if exists sales.return_receipt_org_time_idx;
drop trigger if exists return_effect_request_immutable on sales.return_effect_request;
drop trigger if exists return_disposition_immutable on sales.return_disposition;
drop trigger if exists return_receipt_immutable on sales.return_receipt;
drop function if exists sales.prevent_return_processing_mutation();
drop table if exists sales.return_effect_request;
drop trigger if exists return_disposition_validate on sales.return_disposition;
drop function if exists sales.validate_return_disposition_insert();
drop table if exists sales.return_disposition;
drop trigger if exists return_receipt_validate on sales.return_receipt;
drop function if exists sales.validate_return_receipt_insert();
drop table if exists sales.return_receipt;

commit;
````

### FILE: `db/tests/0018_return_receipt_disposition_effects.test.sql`

```yaml
block_id: "GO-FRANCHISE-JOURNEY:file:33"
operation: CREATE
provenance: AUTHORED
source: "local invariant test governed by Microsoft Business Central return receipt separation and PostgreSQL 18 constraints"
license: "LicenseRef-Workspace-Owner"
sha256: "74ac0b6ded07730909db45d8547f372445c46a82f12704bb2c1f878f4e2ac47b"
variables: []
secrets_allowed: false
```

````sql
begin;

do $$
begin
  if to_regclass('sales.return_receipt') is null or to_regclass('sales.return_disposition') is null or to_regclass('sales.return_effect_request') is null then
    raise exception 'return processing tables missing';
  end if;
  if (select count(*) from pg_trigger where tgname in ('return_receipt_validate','return_disposition_validate','return_receipt_immutable','return_disposition_immutable','return_effect_request_immutable') and not tgisinternal) <> 5 then
    raise exception 'return processing trigger contract incomplete';
  end if;
  if not exists (select 1 from pg_indexes where schemaname='sales' and indexname='return_effect_request_state_idx') then
    raise exception 'return effect owner queue index missing';
  end if;
end $$;

rollback;
````

## 6. Configuration surface

No new environment variable or secret is introduced. Runtime uses the composed `DATABASE_URL`, OIDC verifier, organization claims and existing time/ID providers. Market, currency, price book, location publication, permissions and workflow states are persisted governed data; invalid combinations fail before a remote effect.

| Variable | Type | Default | Secret | Validation | Effect |
|---|---|---|---|---|---|
| none | n/a | n/a | no | n/a | consumes existing backend configuration |

## 7. Dependency bill

| Package/image/tool | Pin exacto | Uso | Licencia | Runtime/build | Fuente oficial |
|---|---|---|---|---|---|
| Go | 1.26.7 | service, HTTP and tests | BSD-3-Clause | build/runtime | https://go.dev/dl/ |
| PostgreSQL | 18.6 | constraints, transactions and integration gate | PostgreSQL | runtime/test | https://www.postgresql.org/download/ |
| pgx | 5.10.0 | PostgreSQL driver | MIT | runtime | https://github.com/jackc/pgx |

The Windows PostgreSQL gate used the EDB binary archive linked by PostgreSQL.org: 343,808,005 bytes, locally captured SHA-256 `fbe23da234ee31547bf8a36d29dfd81e82b849df2d2b78d2eecb43d360252f8c`. This records the verification toolchain; it is not redistributed by this pack.

## 8. Apply order

Apply foundation migrations 0001-0004, then 0005, 0006, 0007, 0008, 0009, 0015, 0016, 0017 and 0018. Materialize service, repository and HTTP module, then wire `FranchiseJourneyModule` into the existing variadic enterprise router. Configure working windows, approved checklist/rejection/condition/inventory-action policies and authorized roles before exposure. Map every immutable effect request to the selected inventory, accounting, payment, fulfillment and fiscal owners; never interpret `requested` as completion. Stop on another CRM/quotation/order/handover/calendar/resource/exception/return owner. Rollback removes routing first, drains/reconciles outbox consumers, preserves all evidence, and runs extension downs in reverse order only during an approved reversible release.

## 9. Verification

1. Materialize the full backend profile into an empty directory and run `go test ./...`; expected exit 0.
2. Start PostgreSQL 18.6 in an isolated environment, apply migrations 0001-0018 with `ON_ERROR_STOP=1`, then execute all journey SQL tests; expected exit 0.
3. Set `TEST_DATABASE_URL` to that isolated database and run `go test ./internal/platform/postgres -run TestFranchiseJourneyPersistenceIsolationAndReplay -count=1`; expected exit 0.
4. Execute HTTP and PostgreSQL tests for strict JSON, authentication, permission, organization scope, customer-owned rejection, exact checklist/serial binding, one-time exception resolution, successor constraint, authorization/order/stock/cost linkage, wrong return serial, duplicate receipt/disposition, immutable receipt and exact owner/effect mapping; expected exit 0.
5. Run `go vet ./...` and build every backend command; expected exit 0.
6. Before production, connect the selected authoritative workforce/calendar system, define timezone/recurrence, employment/privacy and cancellation/no-show policy, then prove real IdP, edge anti-abuse, tax/financing, legal quote, evidence-retention, load, browser and recovery acceptance for the selected project.

## 10. Reconstruction evidence

On 2026-08-30, V140 author-tree PostgreSQL 18.6 applied migrations 0001-0018 and the integration passed customer/organization rejection scope, immutable exception/authorization/receipt/disposition/effect evidence, wrong-serial denial, append-only successor, duplicate receipt/disposition denial and exact owner mappings. Go unit/HTTP/repository tests and the web BFF contract were green. Fresh Markdown reconstruction, production build and global gates govern final admission; downstream adapters, reconciliation, exact-cost accounting/fiscal posting and legal policy remain conditions.

## Delta V330 — verificación multi-tab del flujo existente

Cambio AUTHORED sólo de pruebas; no amplía lógica de producto ni admisión.
ELITE_RETURN_MULTITAB_E2E=1 requiere fixtures ORDER/DELIVERY_READ/RETURN_RECOVERY.
Seis carreras entre páginas con cookie compartida y sessionStorage independiente;
200/409, efectos durables únicos y recuperación por digest. Informe
reconstruction_evidence/RETURN_MULTITAB_VERIFICATION_V330.md. Conserva condiciones
de target, providers, seguridad, operación y aceptación aún no demostradas.

## Delta V331 — fronteras de sesión de devolución

Dos archivos de pruebas AUTHORED; producto y dependencias intactos.
ELITE_RETURN_SESSION_E2E=1 requiere ELITE_RETURN_MULTITAB_E2E y sus fixtures.
Comprueba lector, cookie ausente, organización y segundo operador autorizado;
restaura referencia original sin reenviar. No prueba revocación IdP ni borrado
de información ya entregada al cliente. Evidencia RETURN_SESSION_BOUNDARIES_V331.md
en reconstruction_evidence; condiciones productivas permanecen abiertas.

V402 composed delta: Payment/initial-handover composition. New behavior and tests belong to the explicit runtime/portal packs; source and library release claims remain bounded to their evidence. Existing source provenance is preserved.

Canonical V402 integration: selected by the current profile with exact dependencies and caller overlays. Metadata promotion records byte reconstruction, not closure of every admission/release gate. Payload provenance is unchanged.
