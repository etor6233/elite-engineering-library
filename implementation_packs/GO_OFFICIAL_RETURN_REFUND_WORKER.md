# Go Official Return Refund Worker

V372: montaje opt-in OTLP/HTTP con mTLS y configuración hash-bound en main real;
logs/métricas/trazas de cinco outcomes cerrados, entrega desconocida detiene
claims sin replay. Adaptador Win32 cooperativo, fallo de watcher propagado y
probe CLI finito. Referencia Windows con PostgreSQL, Collector y Prometheus
verificada por el pack WINDOWS-REFERENCE-TELEMETRY-RUNTIME. No garantiza pagos,
identidad/retención/alertas productivas ni sustituye gates de cada proyecto.
Evidencia reconstruction_evidence/INTEGRATED_TELEMETRY_CONTROL_V372.md.

V344: pérdida de escritura del log detiene el host antes de otro claim; error
estático y exit1 aun con stderr inutilizable. Se conservan privacidad, formato
y política de pasos/idle cuando el reporte funciona. No cambia reembolsos,
SDKs, resultados durables ni activa providers. Evidencia
reconstruction_evidence/REFUND_HOST_REPORT_DELIVERY_V344.md.

V318: host conserva exit1 y códigos operativos fijos sin serializar errores de
configuración/DB/provider. Dos regresiones raíz incluyen proceso real y causa
con marcador privado. Sin cambio de retry, resultados durables o reglas de refund.
Evidencia reconstruction_evidence/REFUND_HOST_LOG_PRIVACY_V318.md.

V313: actualización de seguridad por GO-2026-5970/CVE-2026-56852.
x/text0.29.0 ->0.39.0, x/sync0.17.0 ->0.21.0 y piso explícito
x/mod0.40.0 por GO-2026-6179/6180; conservar override al mantener locks;
artefactos oficiales Go con sumdb, commit y BSD-3-Clause verificados.
Sin backport local ni cambio de negocio. El fixture de reembolso exige DB
descartable exclusiva; su seed no prueba constraints del journey de origen.
Evidencia reconstruction_evidence/COMPOSITION_DEPENDENCY_SECURITY_V313.md.
No promoción integral, runtime monitor ni aprobación de producción.

V313: actualización de seguridad por GO-2026-5970/CVE-2026-56852.
x/text0.29.0 ->0.39.0 y mínimo transitivo x/sync0.17.0 ->0.21.0;
artefactos oficiales Go con sumdb, commit y BSD-3-Clause verificados.
Sin backport local ni cambio de negocio. El fixture de reembolso exige DB
descartable exclusiva; su seed no prueba constraints del journey de origen.
Evidencia reconstruction_evidence/COMPOSITION_DEPENDENCY_SECURITY_V313.md.
No promoción integral, runtime monitor ni aprobación de producción.

V402 / 0.1.6: one explicit compatibility mapping at refund preparation accepts
Commerce `mercadopago` and historical `mercado_pago`, preserving the original
payment row and all IDs/idempotency keys. Existing refund rows, SDKs and financial
rules are unchanged. PostgreSQL fixtures with the real pinned SDK and an
in-process transport passed both aliases, partial/full refunds and retrieval.
Source remains AUTHORED; only this delta is compatibility glue. Evidence:
reconstruction_evidence/REFUND_PROVIDER_ALIAS_COMPATIBILITY_V402.md.

## 1. Metadata

```yaml
pack_id: "GO-OFFICIAL-RETURN-REFUND-WORKER"
pack_version: "0.1.6"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa un worker Go autónomo que resuelve la línea y el pago exactos, crea/reconcilia refunds parciales mediante los SDK oficiales Stripe y Mercado Pago, conserva observaciones inmutables y sólo completa el efecto ante estado final verificable."
stacks: ["Go 1.26.8", "PostgreSQL 18.6", "pgx 5.10.0", "stripe-go/v86 86.3.0", "mercadopago/sdk-go 1.14.0"]
compatible_with: ["GO-RETURN-EFFECT-EXECUTION-WORKER 0.1.x", "GO-COMMERCE-PRICING-PAYMENT-API 0.3.x", "GO-COMMERCE-PRICING-PAYMENT-API 0.6.x", "GO-OFFICIAL-PAYMENT-WEBHOOK-ADAPTERS 0.2.x", "GO-ELECTROMOBILITY-APPLICATION 1.9.x"]
incompatible_with: ["unmapped multi-payment order", "client-supplied refund amount", "provider HTTP calls without stable idempotency", "Mercado Pago currency other than ARS without admitted exponent contract"]
license_expression: "LicenseRef-Workspace-Owner AND MIT dependencies"
upstream_sources: ["https://github.com/stripe/stripe-go/tree/a2df585a800a97fe8ec4ebf551b4449bdb3d90a1", "https://docs.stripe.com/api/refunds/create", "https://docs.stripe.com/api/idempotent_requests", "https://github.com/mercadopago/sdk-go/tree/f910ee53fbb6819e435eaf3d0f800cb1fe74ae09", "https://www.mercadopago.com.ar/developers/es/reference/online-payments/checkout-api-payments/create-refund/post", "https://www.mercadopago.com.ar/developers/es/reference/online-payments/checkout-api-payments/get-refund/get", "https://www.mercadopago.com.ar/developers/es/docs/sales-processing/cancellations-and-refunds", "https://www.postgresql.org/docs/18/explicit-locking.html", "https://www.postgresql.org/docs/18/sql-select.html"]
verified_at: "2026-09-11"
```

All twenty-five materialized files are local `AUTHORED` implementation. Provider serialization, authentication, retry primitives and response types are invoked through the exact official MIT SDK releases above; no local line is represented as copied Stripe or Mercado Pago source.

## 2. Applicability

Use after return receipt/disposition and durable return-effect execution when a project has selected Stripe or Mercado Pago and needs a real payment owner. The common path requires one sold stock unit allocated to exactly one quantity-one order line and exactly one eligible captured/disputed payment with sufficient remaining balance. It rejects missing or ambiguous mappings instead of selecting a payment or amount heuristically.

The Mercado Pago adapter is deliberately limited to ARS with two decimal minor units. Other currencies/markets require a separately admitted exponent and sandbox corpus. Production remains conditioned on provider account enablement, target secret manager, live webhook/reconciliation, refund policy approval, accounting/fiscal effects and target load/security/recovery evidence.

## 3. Architecture contract

`payment.return_refund` is the durable projection and `payment.return_refund_observation` is append-only provider evidence. A claim uses `FOR UPDATE SKIP LOCKED`, lease, attempt number and fencing token. Preparation locks the exact return, line and eligible payments; it derives amount/currency only from server-owned order data and subtracts earlier confirmed refunds. Multiple eligible payments fail closed.

The remote call occurs outside the database transaction. Creation uses the return request's stable idempotency key through the official SDK. If the process fails after the provider accepted the request, retry repeats that key; once an ID is persisted, later attempts retrieve rather than create. Stripe requires exact payment/charge reference, refund ID, amount, currency and known status. Mercado Pago additionally retrieves the original payment through its official SDK, verifies payment ID/currency/original amount, sends `X-Idempotency-Key`, and reconciles the refund by payment/refund IDs.

Pending and action-required responses remain retryable with immutable observations. Unknown statuses and response divergence block. Success writes observation, refund projection, attempt, execution result and outbox atomically. Partial refunds do not mutate the payment state; only an exact sum equal to the captured amount transitions the payment to `refunded`. Provider secrets never enter Markdown, persistence, events or logs.

## 4. Exact file manifest

```text
CREATE db/migrations/0020_return_refund_execution.up.sql
CREATE db/migrations/0020_return_refund_execution.down.sql
CREATE db/tests/0020_return_refund_execution.test.sql
CREATE return_refund_worker/go.mod
CREATE return_refund_worker/go.sum
CREATE return_refund_worker/internal/refundworker/model.go
CREATE return_refund_worker/internal/refundworker/provider.go
CREATE return_refund_worker/internal/refundworker/processor.go
CREATE return_refund_worker/internal/refundworker/postgres.go
CREATE return_refund_worker/internal/refundworker/provider_test.go
CREATE return_refund_worker/internal/refundworker/processor_test.go
CREATE return_refund_worker/internal/refundworker/postgres_integration_test.go
CREATE return_refund_worker/cmd/return-refund-worker/main.go
CREATE return_refund_worker/cmd/return-refund-worker/main_test.go
CREATE return_refund_worker/README.md
CREATE return_refund_worker/internal/operationaltelemetry/reporter.go
CREATE return_refund_worker/internal/operationaltelemetry/reporter_test.go
CREATE return_refund_worker/cmd/return-refund-worker/native_stop_other.go
CREATE return_refund_worker/cmd/return-refund-worker/native_stop_windows.go
CREATE return_refund_worker/cmd/return-refund-worker/native_stop_windows_test.go
CREATE return_refund_worker/cmd/return-refund-worker/operational_telemetry_test.go
CREATE return_refund_worker/cmd/telemetry-reference/certificates.go
CREATE return_refund_worker/cmd/telemetry-reference/main.go
CREATE return_refund_worker/docs/COMMERCE_PROVIDER_ALIAS.md
CREATE return_refund_worker/internal/refundworker/provider_alias_integration_test.go
```

## 5. Materialization blocks

### FILE: `db/migrations/0020_return_refund_execution.up.sql`
```yaml
block_id: "GO-OFFICIAL-RETURN-REFUND:file:01"
operation: CREATE
provenance: AUTHORED
source: "local implementation governed by official sources in metadata"
license: "LicenseRef-Workspace-Owner"
sha256: "08fcbcd6a3515423f4c624315d2dbf1e101fff250fac4d699cceca65f8eb4b79"
variables: []
secrets_allowed: false
```
````sql
begin;

create table payment.return_refund (
  tenant_id uuid not null,
  request_id text not null,
  payment_attempt_id text not null,
  order_id text not null,
  line_id text not null,
  stock_unit_id text not null,
  provider_code text not null check (provider_code in ('stripe','mercado_pago')),
  provider_payment_reference text not null check (length(provider_payment_reference) between 1 and 255),
  idempotency_key text not null check (length(idempotency_key) between 16 and 128),
  currency text not null check (currency ~ '^[A-Z]{3}$'),
  amount_minor_units bigint not null check (amount_minor_units>0),
  provider_refund_reference text,
  provider_status text,
  state text not null check (state in ('prepared','pending','requires_action','blocked','succeeded','failed','canceled')),
  response_sha256_hex text,
  version bigint not null default 1 check (version>0),
  created_at timestamptz not null default clock_timestamp(),
  updated_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id,request_id),
  foreign key (tenant_id,request_id) references sales.return_effect_request(tenant_id,request_id),
  foreign key (tenant_id,payment_attempt_id) references payment.payment_attempt(tenant_id,payment_attempt_id),
  foreign key (tenant_id,order_id,line_id) references sales.customer_order_line(tenant_id,order_id,line_id),
  foreign key (tenant_id,stock_unit_id) references inventory.stock_unit(tenant_id,stock_unit_id),
  unique (tenant_id,idempotency_key),
  check (provider_refund_reference is null or length(provider_refund_reference) between 1 and 255),
  check (response_sha256_hex is null or response_sha256_hex ~ '^[0-9a-f]{64}$'),
  check ((state='prepared' and provider_refund_reference is null and provider_status is null and response_sha256_hex is null) or
         (state<>'prepared' and provider_refund_reference is not null and provider_status is not null and response_sha256_hex is not null))
);

create unique index return_refund_provider_reference_uq
  on payment.return_refund(tenant_id,provider_code,provider_refund_reference)
  where provider_refund_reference is not null;

create table payment.return_refund_observation (
  tenant_id uuid not null,
  observation_id text not null,
  request_id text not null,
  provider_refund_reference text not null,
  provider_status text not null,
  amount_minor_units bigint not null check (amount_minor_units>0),
  currency text not null check (currency ~ '^[A-Z]{3}$'),
  response_sha256_hex text not null check (response_sha256_hex ~ '^[0-9a-f]{64}$'),
  observed_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id,observation_id),
  foreign key (tenant_id,request_id) references payment.return_refund(tenant_id,request_id)
);

create or replace function payment.enforce_return_refund_lifecycle() returns trigger language plpgsql as $$
begin
  if tg_op='DELETE' then
    raise exception using errcode='23514',message='return refund projection cannot be deleted';
  end if;
  if (new.tenant_id,new.request_id,new.payment_attempt_id,new.order_id,new.line_id,new.stock_unit_id,new.provider_code,new.provider_payment_reference,new.idempotency_key,new.currency,new.amount_minor_units,new.created_at)
       is distinct from
     (old.tenant_id,old.request_id,old.payment_attempt_id,old.order_id,old.line_id,old.stock_unit_id,old.provider_code,old.provider_payment_reference,old.idempotency_key,old.currency,old.amount_minor_units,old.created_at)
     or new.version<>old.version+1 or new.updated_at<old.updated_at
     or old.state in ('succeeded','failed','canceled')
     or (old.state='prepared' and new.state not in ('pending','requires_action','blocked','succeeded','failed','canceled'))
     or (old.state='pending' and new.state not in ('pending','requires_action','blocked','succeeded','failed','canceled'))
     or (old.state='requires_action' and new.state not in ('requires_action','pending','blocked','succeeded','failed','canceled'))
     or (old.state='blocked' and new.state not in ('blocked','pending','requires_action','succeeded','failed','canceled')) then
    raise exception using errcode='23514',message='invalid return refund transition';
  end if;
  return new;
end $$;

create trigger return_refund_lifecycle before update or delete on payment.return_refund
for each row execute function payment.enforce_return_refund_lifecycle();

create or replace function payment.prevent_return_refund_observation_mutation() returns trigger language plpgsql as $$
begin
  raise exception using errcode='23514',message='return refund observation is immutable';
end $$;

create trigger return_refund_observation_immutable before update or delete on payment.return_refund_observation
for each row execute function payment.prevent_return_refund_observation_mutation();

create index return_refund_payment_idx on payment.return_refund(tenant_id,payment_attempt_id,state,request_id);
create index return_refund_observation_request_idx on payment.return_refund_observation(tenant_id,request_id,observed_at,observation_id);

commit;
````

### FILE: `db/migrations/0020_return_refund_execution.down.sql`
```yaml
block_id: "GO-OFFICIAL-RETURN-REFUND:file:02"
operation: CREATE
provenance: AUTHORED
source: "local implementation governed by official sources in metadata"
license: "LicenseRef-Workspace-Owner"
sha256: "87f20db13360f6b3236cb5d34a4a526a47d953f81af93974f7447043e2d291d8"
variables: []
secrets_allowed: false
```
````sql
begin;

drop index if exists payment.return_refund_observation_request_idx;
drop index if exists payment.return_refund_payment_idx;
drop trigger if exists return_refund_observation_immutable on payment.return_refund_observation;
drop function if exists payment.prevent_return_refund_observation_mutation();
drop trigger if exists return_refund_lifecycle on payment.return_refund;
drop function if exists payment.enforce_return_refund_lifecycle();
drop table if exists payment.return_refund_observation;
drop index if exists payment.return_refund_provider_reference_uq;
drop table if exists payment.return_refund;

commit;
````

### FILE: `db/tests/0020_return_refund_execution.test.sql`
```yaml
block_id: "GO-OFFICIAL-RETURN-REFUND:file:03"
operation: CREATE
provenance: AUTHORED
source: "local implementation governed by official sources in metadata"
license: "LicenseRef-Workspace-Owner"
sha256: "89616f7b55b225a59cd76efd3154ca7d04ac556f89907b141cdb9381827cb285"
variables: []
secrets_allowed: false
```
````sql
begin;

do $$
begin
  if to_regclass('payment.return_refund') is null or to_regclass('payment.return_refund_observation') is null then
    raise exception 'return refund tables missing';
  end if;
  if (select count(*) from pg_trigger where tgname in ('return_refund_lifecycle','return_refund_observation_immutable') and not tgisinternal)<>2 then
    raise exception 'return refund lifecycle triggers missing';
  end if;
  if not exists (select 1 from pg_indexes where schemaname='payment' and indexname='return_refund_provider_reference_uq') then
    raise exception 'return refund provider identity index missing';
  end if;
end $$;

rollback;
````

### FILE: `return_refund_worker/go.mod`
```yaml
block_id: "GO-OFFICIAL-RETURN-REFUND:file:04"
operation: CREATE
provenance: AUTHORED
source: "dependency pins for official provider SDKs and pgx"
license: "LicenseRef-Workspace-Owner"
sha256: "32f96cb7d5255316b23ab3f597d5fa6ef508a64fd444c0ece5cd9e3ca9803aba"
variables: []
secrets_allowed: false
```
````go
module elite.local/return-refund-worker

go 1.26.0

require (
	github.com/jackc/pgx/v5 v5.10.0
	github.com/mercadopago/sdk-go v1.14.0
	github.com/stripe/stripe-go/v86 v86.3.0
)

require (
	github.com/google/uuid v1.6.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	// Security floor for selected upstream tools: GO-2026-6179/6180.
	// Preserve this constraint when tidying; validate the complete module graph.
	golang.org/x/mod v0.40.0 // indirect
	golang.org/x/sync v0.21.0 // indirect
	golang.org/x/text v0.39.0 // indirect
)
````

### FILE: `return_refund_worker/go.sum`
```yaml
block_id: "GO-OFFICIAL-RETURN-REFUND:file:05"
operation: CREATE
provenance: AUTHORED
source: "Go checksum database resolution for exact module pins"
license: "LicenseRef-Workspace-Owner"
sha256: "ea80bd93ebcdd4c1cd14c9b759c419f1c903ccd007b57014e16827fbfb14a13a"
variables: []
secrets_allowed: false
```
````text
github.com/davecgh/go-spew v1.1.0/go.mod h1:J7Y8YcW2NihsgmVo/mv3lAwl/skON4iLHjSsI+c5H38=
github.com/davecgh/go-spew v1.1.1 h1:vj9j/u1bqnvCEfJOwUhtlOARqs3+rkHYY13jYWTU97c=
github.com/davecgh/go-spew v1.1.1/go.mod h1:J7Y8YcW2NihsgmVo/mv3lAwl/skON4iLHjSsI+c5H38=
github.com/google/uuid v1.6.0 h1:NIvaJDMOsjHA8n1jAhLSgzrAzy1Hgr+hNrb57e+94F0=
github.com/google/uuid v1.6.0/go.mod h1:TIyPZe4MgqvfeYDBFedMoGGpEw/LqOeaOT+nhxU+yHo=
github.com/jackc/pgpassfile v1.0.0 h1:/6Hmqy13Ss2zCq62VdNG8tM1wchn8zjSGOBJ6icpsIM=
github.com/jackc/pgpassfile v1.0.0/go.mod h1:CEx0iS5ambNFdcRtxPj5JhEz+xB6uRky5eyVu/W2HEg=
github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 h1:iCEnooe7UlwOQYpKFhBabPMi4aNAfoODPEFNiAnClxo=
github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761/go.mod h1:5TJZWKEWniPve33vlWYSoGYefn3gLQRzjfDlhSJ9ZKM=
github.com/jackc/pgx/v5 v5.10.0 h1:VhSvgU2jSli8o3AqIEOTJr7rZwAEUVo4E4XhR94Zfr0=
github.com/jackc/pgx/v5 v5.10.0/go.mod h1:mal1tBGAFfLHvZzaYh77YS/eC6IX9OWbRV1QIIM0Jn4=
github.com/jackc/puddle/v2 v2.2.2 h1:PR8nw+E/1w0GLuRFSmiioY6UooMp6KJv0/61nB7icHo=
github.com/jackc/puddle/v2 v2.2.2/go.mod h1:vriiEXHvEE654aYKXXjOvZM39qJ0q+azkZFrfEOc3H4=
github.com/mercadopago/sdk-go v1.14.0 h1:3PYp9GPa+iysx2lcaKbpBEkXgEw4IIpbw3T6Jhl0IzI=
github.com/mercadopago/sdk-go v1.14.0/go.mod h1:hvQlOYb3MuYPGfjox7jeGBFbeL0nS+iPwUnxAv6yGWI=
github.com/pmezard/go-difflib v1.0.0 h1:4DBwDE0NGyQoBHbLQYPwSUPoCMWR5BEzIk/f1lZbAQM=
github.com/pmezard/go-difflib v1.0.0/go.mod h1:iKH77koFhYxTK1pcRnkKkqfTogsbg7gZNVY4sRDYZ/4=
github.com/stretchr/objx v0.1.0/go.mod h1:HFkY916IF+rwdDfMAkV7OtwuqBVzrE8GR6GFx+wExME=
github.com/stretchr/testify v1.3.0/go.mod h1:M5WIy9Dh21IEIfnGCwXGc5bZfKNJtfHm1UVUgZn+9EI=
github.com/stretchr/testify v1.7.0/go.mod h1:6Fq8oRcR53rry900zMqJjRRixrwX3KX962/h/Wwjteg=
github.com/stretchr/testify v1.11.1 h1:7s2iGBzp5EwR7/aIZr8ao5+dra3wiQyKjjFuvgVKu7U=
github.com/stretchr/testify v1.11.1/go.mod h1:wZwfW3scLgRK+23gO65QZefKpKQRnfz6sD981Nm4B6U=
github.com/stripe/stripe-go/v86 v86.3.0 h1:BKtYc3NtRa4EGzKAmp4jvl5q7kk2rwMZ+llF18N5vHI=
github.com/stripe/stripe-go/v86 v86.3.0/go.mod h1:Co7QRXCKGNOPTugAdvjgRo+KcMtd9hxy+pZMN0yThsQ=
golang.org/x/mod v0.40.0 h1:hUv+3cXcdRHz08UmSiOob7sadHig73uo5bkXxQ/tvUs=
golang.org/x/mod v0.40.0/go.mod h1:0/weTWkPWGBikyTWAX3dkjVztMmBA5hM0DH6BElSupE=
golang.org/x/sync v0.21.0 h1:HLII4xRRTtCRkxYp4HNFF0Js/Og6q2i++KXbg0gHCwM=
golang.org/x/sync v0.21.0/go.mod h1:9xrNwdLfx4jkKbNva9FpL6vEN7evnE43NNNJQ2LF3+0=
golang.org/x/text v0.39.0 h1:UbZz4pLOvn600D6Oh6GGEI6VAmndrEBLv8/6BEXzyus=
golang.org/x/text v0.39.0/go.mod h1:3UwRclnC2g0TU9x8PZiyfOajCd1zaUNHF9cvqcQZ+ZM=
gopkg.in/check.v1 v0.0.0-20161208181325-20d25e280405/go.mod h1:Co6ibVJAznAaIkqp8huTwlJQCZ016jof/cbN4VW5Yz0=
gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c/go.mod h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZviUmN+EgEM=
gopkg.in/yaml.v3 v3.0.1 h1:fxVm/GzAzEWqLHuvctI91KS9hhNmmWOoWu0XTYJS7CA=
gopkg.in/yaml.v3 v3.0.1/go.mod h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZviUmN+EgEM=
````

### FILE: `return_refund_worker/internal/refundworker/model.go`
```yaml
block_id: "GO-OFFICIAL-RETURN-REFUND:file:06"
operation: CREATE
provenance: AUTHORED
source: "local domain boundary governed by official sources in metadata"
license: "LicenseRef-Workspace-Owner"
sha256: "2ff957347c81bcc5d18508261da50acac3d27fc5fb5e0c2cef39582b7c5c4a5b"
variables: []
secrets_allowed: false
```
````go
package refundworker

import (
	"context"
	"errors"
	"time"
)

var (
	ErrNoWork           = errors.New("no refund work")
	ErrMappingConflict  = errors.New("refund source mapping conflict")
	ErrStaleClaim       = errors.New("stale refund claim")
	ErrInvalidProvider  = errors.New("invalid refund provider request")
	ErrResponseMismatch = errors.New("refund provider response mismatch")
)

type Work struct {
	TenantID       string
	RequestID      string
	IdempotencyKey string
	Attempt        int
	ClaimToken     string
}

type Refund struct {
	TenantID                 string
	RequestID                string
	PaymentAttemptID         string
	OrderID                  string
	LineID                   string
	StockUnitID              string
	Provider                 string
	ProviderPaymentReference string
	IdempotencyKey           string
	Currency                 string
	AmountMinorUnits         int64
	ProviderRefundReference  string
	ProviderStatus           string
}

type ProviderResult struct {
	ProviderPaymentReference string
	ProviderRefundReference  string
	ProviderStatus           string
	Currency                 string
	AmountMinorUnits         int64
}

type Provider interface {
	Create(context.Context, Refund) (ProviderResult, error)
	Retrieve(context.Context, Refund) (ProviderResult, error)
}

type Store interface {
	Claim(context.Context, string, string, time.Duration) (*Work, error)
	Prepare(context.Context, Work, string) (Refund, error)
	Complete(context.Context, Work, string, Refund, ProviderResult, string, string, time.Duration) error
	Finish(context.Context, Work, string, string, string, time.Duration) error
}
````

### FILE: `return_refund_worker/internal/refundworker/provider.go`
```yaml
block_id: "GO-OFFICIAL-RETURN-REFUND:file:07"
operation: CREATE
provenance: AUTHORED
source: "local adapter invoking exact official Stripe and Mercado Pago SDK APIs"
license: "LicenseRef-Workspace-Owner"
sha256: "3e9332e64d571b620f025dca88013f83da772fdad253078ffbffaa17aee31b74"
variables: []
secrets_allowed: false
```
````go
package refundworker

import (
	"context"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/mercadopago/sdk-go/pkg/config"
	"github.com/mercadopago/sdk-go/pkg/payment"
	mprefund "github.com/mercadopago/sdk-go/pkg/refund"
	"github.com/mercadopago/sdk-go/pkg/requester"
	"github.com/mercadopago/sdk-go/pkg/requestoptions"
	"github.com/stripe/stripe-go/v86"
)

var (
	idempotencyPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:-]{15,127}$`)
	currencyPattern    = regexp.MustCompile(`^[A-Z]{3}$`)
)

type stripeRefundAPI interface {
	Create(context.Context, *stripe.RefundCreateParams) (*stripe.Refund, error)
	Retrieve(context.Context, string, *stripe.RefundRetrieveParams) (*stripe.Refund, error)
}

type StripeProvider struct{ refunds stripeRefundAPI }

func NewStripeProvider(secretKey string) (*StripeProvider, error) {
	if strings.TrimSpace(secretKey) == "" {
		return nil, ErrInvalidProvider
	}
	return &StripeProvider{refunds: stripe.NewClient(secretKey).V1Refunds}, nil
}

func newStripeProvider(secretKey string, backend stripe.Backend) (*StripeProvider, error) {
	if strings.TrimSpace(secretKey) == "" || backend == nil {
		return nil, ErrInvalidProvider
	}
	backends := &stripe.Backends{API: backend, Connect: backend, Uploads: backend, MeterEvents: backend}
	return &StripeProvider{refunds: stripe.NewClient(secretKey, stripe.WithBackends(backends)).V1Refunds}, nil
}

func validateRefundInput(value Refund, provider string) error {
	if value.Provider != provider || value.AmountMinorUnits <= 0 || !currencyPattern.MatchString(value.Currency) || strings.TrimSpace(value.ProviderPaymentReference) == "" || !idempotencyPattern.MatchString(value.IdempotencyKey) {
		return ErrInvalidProvider
	}
	return nil
}

func stripeResult(value Refund, resource *stripe.Refund) (ProviderResult, error) {
	if resource == nil || resource.ID == "" || resource.Amount != value.AmountMinorUnits || strings.ToUpper(string(resource.Currency)) != value.Currency {
		return ProviderResult{}, ErrResponseMismatch
	}
	paymentReference := ""
	if resource.PaymentIntent != nil {
		paymentReference = resource.PaymentIntent.ID
	}
	if paymentReference == "" && resource.Charge != nil {
		paymentReference = resource.Charge.ID
	}
	if paymentReference != value.ProviderPaymentReference {
		return ProviderResult{}, ErrResponseMismatch
	}
	return ProviderResult{ProviderPaymentReference: paymentReference, ProviderRefundReference: resource.ID, ProviderStatus: string(resource.Status), Currency: value.Currency, AmountMinorUnits: resource.Amount}, nil
}

func (p *StripeProvider) Create(ctx context.Context, value Refund) (ProviderResult, error) {
	if p == nil || p.refunds == nil || ctx == nil || validateRefundInput(value, "stripe") != nil || value.ProviderRefundReference != "" {
		return ProviderResult{}, ErrInvalidProvider
	}
	params := &stripe.RefundCreateParams{Amount: stripe.Int64(value.AmountMinorUnits), Reason: stripe.String(string(stripe.RefundReasonRequestedByCustomer)), Metadata: map[string]string{"return_request_id": value.RequestID}}
	switch {
	case strings.HasPrefix(value.ProviderPaymentReference, "pi_"):
		params.PaymentIntent = stripe.String(value.ProviderPaymentReference)
	case strings.HasPrefix(value.ProviderPaymentReference, "ch_"):
		params.Charge = stripe.String(value.ProviderPaymentReference)
	default:
		return ProviderResult{}, ErrInvalidProvider
	}
	params.SetIdempotencyKey(value.IdempotencyKey)
	resource, err := p.refunds.Create(ctx, params)
	if err != nil {
		return ProviderResult{}, err
	}
	return stripeResult(value, resource)
}

func (p *StripeProvider) Retrieve(ctx context.Context, value Refund) (ProviderResult, error) {
	if p == nil || p.refunds == nil || ctx == nil || validateRefundInput(value, "stripe") != nil || strings.TrimSpace(value.ProviderRefundReference) == "" {
		return ProviderResult{}, ErrInvalidProvider
	}
	resource, err := p.refunds.Retrieve(ctx, value.ProviderRefundReference, nil)
	if err != nil {
		return ProviderResult{}, err
	}
	return stripeResult(value, resource)
}

type mercadoPagoRefundAPI interface {
	CreatePartialRefund(context.Context, int, float64) (*mprefund.Response, error)
	Get(context.Context, int, int) (*mprefund.Response, error)
}

type mercadoPagoPaymentAPI interface {
	Get(context.Context, int) (*payment.Response, error)
}

type MercadoPagoProvider struct {
	refunds  mercadoPagoRefundAPI
	payments mercadoPagoPaymentAPI
}

func NewMercadoPagoProvider(accessToken string) (*MercadoPagoProvider, error) {
	return newMercadoPagoProvider(accessToken, nil)
}

func newMercadoPagoProvider(accessToken string, transport requester.Requester) (*MercadoPagoProvider, error) {
	if strings.TrimSpace(accessToken) == "" {
		return nil, ErrInvalidProvider
	}
	options := []config.Option{config.WithTimeout(10 * time.Second), config.WithMaxRetries(2)}
	if transport != nil {
		options = []config.Option{config.WithHTTPClient(transport)}
	}
	cfg, err := config.New(accessToken, options...)
	if err != nil {
		return nil, err
	}
	return &MercadoPagoProvider{refunds: mprefund.NewClient(cfg), payments: payment.NewClient(cfg)}, nil
}

func parsePositiveInt(value string) (int, error) {
	n, err := strconv.Atoi(value)
	if err != nil || n <= 0 {
		return 0, ErrInvalidProvider
	}
	return n, nil
}

func arsMinorToMajor(value int64) (float64, error) {
	if value <= 0 || value > 9_000_000_000_000_000 {
		return 0, ErrInvalidProvider
	}
	major := float64(value) / 100
	if int64(math.Round(major*100)) != value {
		return 0, ErrInvalidProvider
	}
	return major, nil
}

func arsMajorToMinor(value float64) (int64, error) {
	if value <= 0 || math.IsNaN(value) || math.IsInf(value, 0) {
		return 0, ErrResponseMismatch
	}
	minor := int64(math.Round(value * 100))
	if math.Abs(value-float64(minor)/100) > 0.0000001 {
		return 0, ErrResponseMismatch
	}
	return minor, nil
}

func (p *MercadoPagoProvider) verifyPayment(ctx context.Context, value Refund, paymentID int, requireRefundable bool) error {
	resource, err := p.payments.Get(ctx, paymentID)
	if err != nil {
		return err
	}
	amount, convErr := arsMajorToMinor(resource.TransactionAmount)
	validState := (resource.Status == "approved" && resource.Captured) || (!requireRefundable && resource.Status == "refunded")
	if convErr != nil || resource.ID != paymentID || strings.ToUpper(resource.CurrencyID) != value.Currency || !validState || amount < value.AmountMinorUnits {
		return ErrResponseMismatch
	}
	return nil
}

func mercadoPagoResult(value Refund, paymentID int, resource *mprefund.Response) (ProviderResult, error) {
	if resource == nil || resource.ID <= 0 || resource.PaymentID != paymentID {
		return ProviderResult{}, ErrResponseMismatch
	}
	amount, err := arsMajorToMinor(resource.Amount)
	if err != nil || amount != value.AmountMinorUnits {
		return ProviderResult{}, ErrResponseMismatch
	}
	return ProviderResult{ProviderPaymentReference: strconv.Itoa(paymentID), ProviderRefundReference: strconv.Itoa(resource.ID), ProviderStatus: resource.Status, Currency: value.Currency, AmountMinorUnits: amount}, nil
}

func (p *MercadoPagoProvider) Create(ctx context.Context, value Refund) (ProviderResult, error) {
	if p == nil || p.refunds == nil || p.payments == nil || ctx == nil || validateRefundInput(value, "mercado_pago") != nil || value.Currency != "ARS" || value.ProviderRefundReference != "" {
		return ProviderResult{}, ErrInvalidProvider
	}
	paymentID, err := parsePositiveInt(value.ProviderPaymentReference)
	if err != nil {
		return ProviderResult{}, err
	}
	if err = p.verifyPayment(ctx, value, paymentID, true); err != nil {
		return ProviderResult{}, err
	}
	amount, err := arsMinorToMajor(value.AmountMinorUnits)
	if err != nil {
		return ProviderResult{}, err
	}
	resource, err := p.refunds.CreatePartialRefund(requestoptions.WithIdempotencyKey(ctx, value.IdempotencyKey), paymentID, amount)
	if err != nil {
		return ProviderResult{}, err
	}
	return mercadoPagoResult(value, paymentID, resource)
}

func (p *MercadoPagoProvider) Retrieve(ctx context.Context, value Refund) (ProviderResult, error) {
	if p == nil || p.refunds == nil || p.payments == nil || ctx == nil || validateRefundInput(value, "mercado_pago") != nil || value.Currency != "ARS" || value.ProviderRefundReference == "" {
		return ProviderResult{}, ErrInvalidProvider
	}
	paymentID, err := parsePositiveInt(value.ProviderPaymentReference)
	if err != nil {
		return ProviderResult{}, err
	}
	refundID, err := parsePositiveInt(value.ProviderRefundReference)
	if err != nil {
		return ProviderResult{}, err
	}
	if err = p.verifyPayment(ctx, value, paymentID, false); err != nil {
		return ProviderResult{}, err
	}
	resource, err := p.refunds.Get(ctx, paymentID, refundID)
	if err != nil {
		return ProviderResult{}, err
	}
	return mercadoPagoResult(value, paymentID, resource)
}
````

### FILE: `return_refund_worker/internal/refundworker/processor.go`
```yaml
block_id: "GO-OFFICIAL-RETURN-REFUND:file:08"
operation: CREATE
provenance: AUTHORED
source: "local orchestration governed by official sources in metadata"
license: "LicenseRef-Workspace-Owner"
sha256: "2ba4258f7477d1ea454f52dcde46b96e7768195e954af0d54b9b90d460ca882b"
variables: []
secrets_allowed: false
```
````go
package refundworker

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

type Processor struct {
	store       Store
	providers   map[string]Provider
	workerID    string
	lease       time.Duration
	retryAfter  time.Duration
	maxAttempts int
}

func NewProcessor(store Store, providers map[string]Provider, workerID string, lease, retryAfter time.Duration, maxAttempts int) (*Processor, error) {
	if store == nil || strings.TrimSpace(workerID) == "" || lease <= 0 || retryAfter <= 0 || maxAttempts < 1 {
		return nil, fmt.Errorf("invalid refund processor configuration")
	}
	copyProviders := make(map[string]Provider, len(providers))
	for name, provider := range providers {
		if provider != nil {
			copyProviders[name] = provider
		}
	}
	return &Processor{store: store, providers: copyProviders, workerID: workerID, lease: lease, retryAfter: retryAfter, maxAttempts: maxAttempts}, nil
}

func classifyStatus(provider, status string) (outcome, code string) {
	status = strings.ToLower(strings.TrimSpace(status))
	switch provider {
	case "stripe":
		switch status {
		case "succeeded":
			return "succeeded", ""
		case "pending":
			return "retry", "PROVIDER_PENDING"
		case "requires_action":
			return "retry", "PROVIDER_ACTION_REQUIRED"
		case "failed":
			return "failed", "PROVIDER_REFUND_FAILED"
		case "canceled":
			return "failed", "PROVIDER_REFUND_CANCELED"
		}
	case "mercado_pago":
		switch status {
		case "approved":
			return "succeeded", ""
		case "pending", "in_process":
			return "retry", "PROVIDER_PENDING"
		case "rejected", "failed":
			return "failed", "PROVIDER_REFUND_FAILED"
		case "cancelled", "canceled":
			return "failed", "PROVIDER_REFUND_CANCELED"
		}
	}
	return "blocked", "UNKNOWN_PROVIDER_STATUS"
}

func validateResult(refund Refund, result ProviderResult) error {
	if result.ProviderPaymentReference != refund.ProviderPaymentReference || result.ProviderRefundReference == "" || result.ProviderStatus == "" || result.Currency != refund.Currency || result.AmountMinorUnits != refund.AmountMinorUnits {
		return ErrResponseMismatch
	}
	return nil
}

func (p *Processor) Step(ctx context.Context, claimToken string) error {
	if p == nil || ctx == nil || strings.TrimSpace(claimToken) == "" {
		return fmt.Errorf("invalid refund step")
	}
	work, err := p.store.Claim(ctx, p.workerID, claimToken, p.lease)
	if err != nil {
		return err
	}
	if work == nil {
		return ErrNoWork
	}
	refund, err := p.store.Prepare(ctx, *work, p.workerID)
	if err != nil {
		code := "REFUND_MAPPING_CONFLICT"
		if !errors.Is(err, ErrMappingConflict) {
			code = "REFUND_PREPARATION_FAILED"
		}
		if finishErr := p.store.Finish(ctx, *work, p.workerID, "blocked", code, 0); finishErr != nil {
			return errors.Join(err, finishErr)
		}
		return err
	}
	provider := p.providers[refund.Provider]
	if provider == nil {
		return p.store.Finish(ctx, *work, p.workerID, "blocked", "PROVIDER_CONFIG_MISSING", 0)
	}
	var result ProviderResult
	if refund.ProviderRefundReference == "" {
		result, err = provider.Create(ctx, refund)
	} else {
		result, err = provider.Retrieve(ctx, refund)
	}
	if err != nil {
		outcome := "retry"
		if work.Attempt >= p.maxAttempts || errors.Is(err, ErrInvalidProvider) || errors.Is(err, ErrResponseMismatch) {
			outcome = "blocked"
		}
		if finishErr := p.store.Finish(ctx, *work, p.workerID, outcome, "PROVIDER_CALL_FAILED", p.retryAfter); finishErr != nil {
			return errors.Join(err, finishErr)
		}
		return err
	}
	if err = validateResult(refund, result); err != nil {
		if finishErr := p.store.Complete(ctx, *work, p.workerID, refund, result, "blocked", "PROVIDER_RESPONSE_MISMATCH", 0); finishErr != nil {
			return errors.Join(err, finishErr)
		}
		return err
	}
	outcome, code := classifyStatus(refund.Provider, result.ProviderStatus)
	retry := time.Duration(0)
	if outcome == "retry" {
		retry = p.retryAfter
		if work.Attempt >= p.maxAttempts {
			outcome = "blocked"
			code = "PROVIDER_RECONCILIATION_EXHAUSTED"
			retry = 0
		}
	}
	return p.store.Complete(ctx, *work, p.workerID, refund, result, outcome, code, retry)
}
````

### FILE: `return_refund_worker/internal/refundworker/postgres.go`
```yaml
block_id: "GO-OFFICIAL-RETURN-REFUND:file:09"
operation: CREATE
provenance: AUTHORED
source: "local PostgreSQL implementation governed by official sources in metadata"
license: "LicenseRef-Workspace-Owner"
sha256: "06a0915c86c29d30516ad05a4ab283a1f4f38624224e7a7736e5737230648ae3"
variables: []
secrets_allowed: false
```
````go
package refundworker

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStore struct{ pool *pgxpool.Pool }

func NewPostgresStore(pool *pgxpool.Pool) *PostgresStore { return &PostgresStore{pool: pool} }

func (s *PostgresStore) Claim(ctx context.Context, workerID, claimToken string, lease time.Duration) (*Work, error) {
	if s == nil || s.pool == nil || workerID == "" || claimToken == "" || lease <= 0 {
		return nil, fmt.Errorf("invalid refund claim")
	}
	var work Work
	err := s.pool.QueryRow(ctx, `
with candidate as (
  select x.tenant_id,x.request_id
    from sales.return_effect_execution x
    join sales.return_effect_request e on e.tenant_id=x.tenant_id and e.request_id=x.request_id
   where e.effect_kind='refund' and e.owner_context='payment'
     and ((x.status in ('requested','retry') and x.available_at<=clock_timestamp()) or (x.status='claimed' and x.claimed_until<clock_timestamp()))
   order by x.available_at,e.requested_at,x.request_id
   for update of x skip locked limit 1
), claimed as (
  update sales.return_effect_execution x
     set status='claimed',attempt_count=x.attempt_count+1,claimed_at=clock_timestamp(),claimed_by=$1,claim_token=$2,claimed_until=clock_timestamp()+$3::interval,last_error_code=null,result_sha256_hex=null,updated_at=clock_timestamp()
    from candidate c where x.tenant_id=c.tenant_id and x.request_id=c.request_id
  returning x.tenant_id,x.request_id,x.attempt_count,x.claim_token
)
select c.tenant_id,c.request_id,e.idempotency_key,c.attempt_count,c.claim_token
  from claimed c join sales.return_effect_request e on e.tenant_id=c.tenant_id and e.request_id=c.request_id`, workerID, claimToken, lease.String()).Scan(&work.TenantID, &work.RequestID, &work.IdempotencyKey, &work.Attempt, &work.ClaimToken)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &work, nil
}

func scanRefund(row pgx.Row) (Refund, error) {
	var value Refund
	err := row.Scan(&value.TenantID, &value.RequestID, &value.PaymentAttemptID, &value.OrderID, &value.LineID, &value.StockUnitID, &value.Provider, &value.ProviderPaymentReference, &value.IdempotencyKey, &value.Currency, &value.AmountMinorUnits, &value.ProviderRefundReference, &value.ProviderStatus)
	return value, err
}

func (s *PostgresStore) Prepare(ctx context.Context, work Work, workerID string) (Refund, error) {
	if s == nil || s.pool == nil {
		return Refund{}, fmt.Errorf("invalid refund store")
	}
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return Refund{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var attempt int
	err = tx.QueryRow(ctx, `select attempt_count from sales.return_effect_execution where tenant_id=$1 and request_id=$2 and status='claimed' and claimed_by=$3 and claim_token=$4 and claimed_until>=clock_timestamp() for update`, work.TenantID, work.RequestID, workerID, work.ClaimToken).Scan(&attempt)
	if errors.Is(err, pgx.ErrNoRows) || attempt != work.Attempt {
		return Refund{}, ErrStaleClaim
	}
	if err != nil {
		return Refund{}, err
	}

	existing, existingErr := scanRefund(tx.QueryRow(ctx, `select tenant_id,request_id,payment_attempt_id,order_id,line_id,stock_unit_id,provider_code,provider_payment_reference,idempotency_key,currency,amount_minor_units,coalesce(provider_refund_reference,''),coalesce(provider_status,'') from payment.return_refund where tenant_id=$1 and request_id=$2 for update`, work.TenantID, work.RequestID))
	if existingErr == nil {
		if err = tx.Commit(ctx); err != nil {
			return Refund{}, err
		}
		return existing, nil
	}
	if !errors.Is(existingErr, pgx.ErrNoRows) {
		return Refund{}, existingErr
	}

	type lineSource struct {
		orderID, lineID, stockID, currency string
		amount                             int64
	}
	rows, err := tx.Query(ctx, `
select r.order_id,l.line_id,r.stock_unit_id,o.currency,l.unit_price_minor_units
  from sales.return_effect_request e
  join sales.return_disposition d on d.tenant_id=e.tenant_id and d.disposition_id=e.disposition_id
  join sales.return_receipt r on r.tenant_id=d.tenant_id and r.receipt_id=d.receipt_id
  join sales.customer_order o on o.tenant_id=r.tenant_id and o.order_id=r.order_id and o.organization_id=r.organization_id
  join sales.customer_order_line l on l.tenant_id=r.tenant_id and l.order_id=r.order_id and l.allocated_stock_unit_id=r.stock_unit_id and l.quantity=1
 where e.tenant_id=$1 and e.request_id=$2 and e.effect_kind='refund' and e.owner_context='payment' and e.state='requested' and d.customer_remedy='refund'
 for update of l,o`, work.TenantID, work.RequestID)
	if err != nil {
		return Refund{}, err
	}
	var lines []lineSource
	for rows.Next() {
		var line lineSource
		if err = rows.Scan(&line.orderID, &line.lineID, &line.stockID, &line.currency, &line.amount); err != nil {
			rows.Close()
			return Refund{}, err
		}
		lines = append(lines, line)
	}
	rows.Close()
	if err = rows.Err(); err != nil {
		return Refund{}, err
	}
	if len(lines) != 1 || lines[0].amount <= 0 {
		return Refund{}, ErrMappingConflict
	}

	type paymentSource struct {
		id, provider, reference, currency string
		amount, refunded                  int64
	}
	rows, err = tx.Query(ctx, `
-- AUTHORED provider-name compatibility at the Commerce -> refund-owner boundary.
-- Preserve payment/refund IDs and idempotency keys; existing refund rows retain
-- their historical mercado_pago provider code and are returned above unchanged.
select p.payment_attempt_id,case p.provider_code when 'mercadopago' then 'mercado_pago' else p.provider_code end,p.provider_reference,p.currency,p.amount_minor_units,
       coalesce((select sum(rr.amount_minor_units) from payment.return_refund rr where rr.tenant_id=p.tenant_id and rr.payment_attempt_id=p.payment_attempt_id and rr.state='succeeded'),0)
  from payment.payment_attempt p
 where p.tenant_id=$1 and p.order_id=$2 and p.state in ('captured','disputed') and p.provider_reference is not null and p.provider_code in ('stripe','mercado_pago','mercadopago')
 for update of p`, work.TenantID, lines[0].orderID)
	if err != nil {
		return Refund{}, err
	}
	var payments []paymentSource
	for rows.Next() {
		var value paymentSource
		if err = rows.Scan(&value.id, &value.provider, &value.reference, &value.currency, &value.amount, &value.refunded); err != nil {
			rows.Close()
			return Refund{}, err
		}
		if value.currency == lines[0].currency && value.amount-value.refunded >= lines[0].amount {
			payments = append(payments, value)
		}
	}
	rows.Close()
	if err = rows.Err(); err != nil {
		return Refund{}, err
	}
	if len(payments) != 1 {
		return Refund{}, ErrMappingConflict
	}

	value := Refund{TenantID: work.TenantID, RequestID: work.RequestID, PaymentAttemptID: payments[0].id, OrderID: lines[0].orderID, LineID: lines[0].lineID, StockUnitID: lines[0].stockID, Provider: payments[0].provider, ProviderPaymentReference: payments[0].reference, IdempotencyKey: work.IdempotencyKey, Currency: lines[0].currency, AmountMinorUnits: lines[0].amount}
	_, err = tx.Exec(ctx, `insert into payment.return_refund(tenant_id,request_id,payment_attempt_id,order_id,line_id,stock_unit_id,provider_code,provider_payment_reference,idempotency_key,currency,amount_minor_units,state) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,'prepared')`, value.TenantID, value.RequestID, value.PaymentAttemptID, value.OrderID, value.LineID, value.StockUnitID, value.Provider, value.ProviderPaymentReference, value.IdempotencyKey, value.Currency, value.AmountMinorUnits)
	if err != nil {
		return Refund{}, err
	}
	if err = tx.Commit(ctx); err != nil {
		return Refund{}, err
	}
	return value, nil
}

func canonicalResult(result ProviderResult) (string, error) {
	encoded, err := json.Marshal(result)
	if err != nil {
		return "", err
	}
	digest := sha256.Sum256(encoded)
	return hex.EncodeToString(digest[:]), nil
}

func refundState(outcome, provider, status string) string {
	if outcome == "succeeded" {
		return "succeeded"
	}
	if outcome == "blocked" {
		return "blocked"
	}
	if outcome == "failed" {
		if (provider == "stripe" && status == "canceled") || (provider == "mercado_pago" && (status == "cancelled" || status == "canceled")) {
			return "canceled"
		}
		return "failed"
	}
	if status == "requires_action" {
		return "requires_action"
	}
	return "pending"
}

func validOutcome(outcome, code string, retry time.Duration) bool {
	if outcome == "succeeded" {
		return code == "" && retry == 0
	}
	if outcome == "retry" {
		return code != "" && retry > 0
	}
	return (outcome == "blocked" || outcome == "failed") && code != "" && retry == 0
}

func (s *PostgresStore) Complete(ctx context.Context, work Work, workerID string, refund Refund, result ProviderResult, outcome, code string, retryAfter time.Duration) error {
	if s == nil || s.pool == nil || !validOutcome(outcome, code, retryAfter) {
		return fmt.Errorf("invalid refund completion")
	}
	resultHash, err := canonicalResult(result)
	if err != nil {
		return err
	}
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var attempt int
	var startedAt time.Time
	err = tx.QueryRow(ctx, `select attempt_count,claimed_at from sales.return_effect_execution where tenant_id=$1 and request_id=$2 and status='claimed' and claimed_by=$3 and claim_token=$4 and claimed_until>=clock_timestamp() for update`, work.TenantID, work.RequestID, workerID, work.ClaimToken).Scan(&attempt, &startedAt)
	if errors.Is(err, pgx.ErrNoRows) || attempt != work.Attempt {
		return ErrStaleClaim
	}
	if err != nil {
		return err
	}
	var stored Refund
	stored, err = scanRefund(tx.QueryRow(ctx, `select tenant_id,request_id,payment_attempt_id,order_id,line_id,stock_unit_id,provider_code,provider_payment_reference,idempotency_key,currency,amount_minor_units,coalesce(provider_refund_reference,''),coalesce(provider_status,'') from payment.return_refund where tenant_id=$1 and request_id=$2 for update`, work.TenantID, work.RequestID))
	if err != nil || stored.PaymentAttemptID != refund.PaymentAttemptID || stored.AmountMinorUnits != refund.AmountMinorUnits || stored.Currency != refund.Currency {
		if err != nil {
			return err
		}
		return ErrMappingConflict
	}
	state := refundState(outcome, refund.Provider, result.ProviderStatus)
	updated, err := tx.Exec(ctx, `update payment.return_refund set provider_refund_reference=$3,provider_status=$4,state=$5,response_sha256_hex=$6,version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and request_id=$2`, work.TenantID, work.RequestID, result.ProviderRefundReference, result.ProviderStatus, state, resultHash)
	if err != nil || updated.RowsAffected() != 1 {
		if err != nil {
			return err
		}
		return ErrMappingConflict
	}
	_, err = tx.Exec(ctx, `insert into payment.return_refund_observation(tenant_id,observation_id,request_id,provider_refund_reference,provider_status,amount_minor_units,currency,response_sha256_hex) values($1,$2,$3,$4,$5,$6,$7,$8)`, work.TenantID, work.ClaimToken, work.RequestID, result.ProviderRefundReference, result.ProviderStatus, result.AmountMinorUnits, result.Currency, resultHash)
	if err != nil {
		return err
	}
	if outcome == "succeeded" {
		var paymentAmount, paymentVersion int64
		var paymentState string
		err = tx.QueryRow(ctx, `select amount_minor_units,state,version from payment.payment_attempt where tenant_id=$1 and payment_attempt_id=$2 for update`, work.TenantID, refund.PaymentAttemptID).Scan(&paymentAmount, &paymentState, &paymentVersion)
		if err != nil {
			return err
		}
		var refunded int64
		if err = tx.QueryRow(ctx, `select coalesce(sum(amount_minor_units),0) from payment.return_refund where tenant_id=$1 and payment_attempt_id=$2 and state='succeeded'`, work.TenantID, refund.PaymentAttemptID).Scan(&refunded); err != nil {
			return err
		}
		if refunded > paymentAmount {
			return ErrMappingConflict
		}
		if refunded == paymentAmount && (paymentState == "captured" || paymentState == "disputed") {
			if _, err = tx.Exec(ctx, `update payment.payment_attempt set state='refunded',version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and payment_attempt_id=$2 and version=$3`, work.TenantID, refund.PaymentAttemptID, paymentVersion); err != nil {
				return err
			}
			if _, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload) values($1,gen_random_uuid(),'payment',$2,$3,'payment.refunded',1,clock_timestamp(),jsonb_build_object('return_request_id',$4::text,'refund_reference',$5::text,'amount_minor_units',$6::bigint,'currency',$7::text))`, work.TenantID, refund.PaymentAttemptID, paymentVersion+1, work.RequestID, result.ProviderRefundReference, refund.AmountMinorUnits, refund.Currency); err != nil {
				return err
			}
		}
		if _, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload) values($1,gen_random_uuid(),'return-refund',$2,1,'return-refund.succeeded',1,clock_timestamp(),jsonb_build_object('payment_attempt_id',$3::text,'line_id',$4::text,'stock_unit_id',$5::text,'provider',$6::text,'refund_reference',$7::text,'amount_minor_units',$8::bigint,'currency',$9::text))`, work.TenantID, work.RequestID, refund.PaymentAttemptID, refund.LineID, refund.StockUnitID, refund.Provider, result.ProviderRefundReference, refund.AmountMinorUnits, refund.Currency); err != nil {
			return err
		}
	}
	if _, err = tx.Exec(ctx, `insert into sales.return_effect_attempt(tenant_id,attempt_id,request_id,attempt_no,worker_id,outcome,error_code,provider_reference,result_sha256_hex,started_at) values($1,$2,$3,$4,$5,$6,nullif($7,''),$8,case when $6='succeeded' then $9 else null end,$10)`, work.TenantID, work.ClaimToken, work.RequestID, attempt, workerID, outcome, code, result.ProviderRefundReference, resultHash, startedAt); err != nil {
		return err
	}
	available := time.Now().UTC()
	if outcome == "retry" {
		available = available.Add(retryAfter)
	}
	updated, err = tx.Exec(ctx, `update sales.return_effect_execution set status=$5,available_at=$6,claimed_at=null,claimed_by=null,claim_token=null,claimed_until=null,last_error_code=nullif($7,''),provider_reference=$8,result_sha256_hex=case when $5='succeeded' then $9 else null end,updated_at=clock_timestamp() where tenant_id=$1 and request_id=$2 and status='claimed' and claimed_by=$3 and claim_token=$4`, work.TenantID, work.RequestID, workerID, work.ClaimToken, outcome, available, code, result.ProviderRefundReference, resultHash)
	if err != nil || updated.RowsAffected() != 1 {
		if err != nil {
			return err
		}
		return ErrStaleClaim
	}
	return tx.Commit(ctx)
}

func (s *PostgresStore) Finish(ctx context.Context, work Work, workerID, outcome, code string, retryAfter time.Duration) error {
	if s == nil || s.pool == nil || !validOutcome(outcome, code, retryAfter) || outcome == "succeeded" {
		return fmt.Errorf("invalid refund finish")
	}
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var attempt int
	var startedAt time.Time
	err = tx.QueryRow(ctx, `select attempt_count,claimed_at from sales.return_effect_execution where tenant_id=$1 and request_id=$2 and status='claimed' and claimed_by=$3 and claim_token=$4 and claimed_until>=clock_timestamp() for update`, work.TenantID, work.RequestID, workerID, work.ClaimToken).Scan(&attempt, &startedAt)
	if errors.Is(err, pgx.ErrNoRows) || attempt != work.Attempt {
		return ErrStaleClaim
	}
	if err != nil {
		return err
	}
	if _, err = tx.Exec(ctx, `insert into sales.return_effect_attempt(tenant_id,attempt_id,request_id,attempt_no,worker_id,outcome,error_code,started_at) values($1,$2,$3,$4,$5,$6,$7,$8)`, work.TenantID, work.ClaimToken, work.RequestID, attempt, workerID, outcome, code, startedAt); err != nil {
		return err
	}
	available := time.Now().UTC()
	if outcome == "retry" {
		available = available.Add(retryAfter)
	}
	updated, err := tx.Exec(ctx, `update sales.return_effect_execution set status=$5,available_at=$6,claimed_at=null,claimed_by=null,claim_token=null,claimed_until=null,last_error_code=$7,result_sha256_hex=null,updated_at=clock_timestamp() where tenant_id=$1 and request_id=$2 and status='claimed' and claimed_by=$3 and claim_token=$4`, work.TenantID, work.RequestID, workerID, work.ClaimToken, outcome, available, code)
	if err != nil || updated.RowsAffected() != 1 {
		if err != nil {
			return err
		}
		return ErrStaleClaim
	}
	return tx.Commit(ctx)
}
````

### FILE: `return_refund_worker/internal/refundworker/provider_test.go`
```yaml
block_id: "GO-OFFICIAL-RETURN-REFUND:file:10"
operation: CREATE
provenance: AUTHORED
source: "local contract tests through exact official SDK HTTP clients"
license: "LicenseRef-Workspace-Owner"
sha256: "d4dce5054edd333a8f016e2b7679b4da6f588b4d00ba11f2f5808aa1d93b7b1b"
variables: []
secrets_allowed: false
```
````go
package refundworker

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stripe/stripe-go/v86"
)

func TestStripeOfficialRefundCreateAndRetrieve(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer sk_test_local" {
			t.Fatal("Stripe authorization missing")
		}
		if r.Method == http.MethodPost {
			if r.URL.Path != "/v1/refunds" || r.Header.Get("Idempotency-Key") != "return-effect-key-0001" {
				t.Fatalf("unexpected create request %s %s", r.Method, r.URL.Path)
			}
			body, _ := io.ReadAll(r.Body)
			form, err := url.ParseQuery(string(body))
			if err != nil {
				t.Fatal(err)
			}
			if form.Get("amount") != "125050" || form.Get("payment_intent") != "pi_123" || form.Get("reason") != "requested_by_customer" || form.Get("metadata[return_request_id]") != "effect-refund" {
				t.Fatalf("unexpected Stripe form: %s", body)
			}
		} else if r.Method != http.MethodGet || r.URL.Path != "/v1/refunds/re_123" {
			t.Fatalf("unexpected retrieve request %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"id":"re_123","object":"refund","amount":125050,"currency":"ars","payment_intent":"pi_123","status":"pending"}`)
	}))
	defer server.Close()
	backend := stripe.GetBackendWithConfig(stripe.APIBackend, &stripe.BackendConfig{URL: stripe.String(server.URL), HTTPClient: server.Client(), MaxNetworkRetries: stripe.Int64(0)})
	provider, err := newStripeProvider("sk_test_local", backend)
	if err != nil {
		t.Fatal(err)
	}
	request := Refund{RequestID: "effect-refund", Provider: "stripe", ProviderPaymentReference: "pi_123", IdempotencyKey: "return-effect-key-0001", Currency: "ARS", AmountMinorUnits: 125050}
	created, err := provider.Create(context.Background(), request)
	if err != nil || created.ProviderRefundReference != "re_123" || created.ProviderStatus != "pending" {
		t.Fatalf("created=%+v err=%v", created, err)
	}
	request.ProviderRefundReference = created.ProviderRefundReference
	retrieved, err := provider.Retrieve(context.Background(), request)
	if err != nil || retrieved != created {
		t.Fatalf("retrieved=%+v err=%v", retrieved, err)
	}
}

type requesterFunc func(*http.Request) (*http.Response, error)

func (f requesterFunc) Do(request *http.Request) (*http.Response, error) { return f(request) }

func TestMercadoPagoOfficialRefundCreateAndRetrieve(t *testing.T) {
	paymentReads := 0
	transport := requesterFunc(func(request *http.Request) (*http.Response, error) {
		if request.Header.Get("Authorization") != "Bearer TEST-local" {
			t.Fatal("Mercado Pago authorization missing")
		}
		path := request.URL.Path
		var response string
		switch {
		case request.Method == http.MethodGet && path == "/v1/payments/7186040733":
			paymentReads++
			if paymentReads == 1 {
				response = `{"id":7186040733,"status":"approved","captured":true,"currency_id":"ARS","transaction_amount":2000.00}`
			} else {
				response = `{"id":7186040733,"status":"refunded","captured":false,"currency_id":"ARS","transaction_amount":2000.00,"transaction_amount_refunded":1250.50}`
			}
		case request.Method == http.MethodPost && path == "/v1/payments/7186040733/refunds":
			if request.Header.Get("X-Idempotency-Key") != "return-effect-key-0002" {
				t.Fatal("Mercado Pago idempotency key missing")
			}
			var body map[string]any
			if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
				t.Fatal(err)
			}
			if body["amount"] != 1250.5 {
				t.Fatalf("unexpected amount: %#v", body)
			}
			response = `{"id":1622029222,"payment_id":7186040733,"status":"approved","amount":1250.50}`
		case request.Method == http.MethodGet && path == "/v1/payments/7186040733/refunds/1622029222":
			response = `{"id":1622029222,"payment_id":7186040733,"status":"approved","amount":1250.50}`
		default:
			t.Fatalf("unexpected Mercado Pago request %s %s", request.Method, path)
		}
		return &http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(response)), Request: request}, nil
	})
	provider, err := newMercadoPagoProvider("TEST-local", transport)
	if err != nil {
		t.Fatal(err)
	}
	request := Refund{RequestID: "effect-refund", Provider: "mercado_pago", ProviderPaymentReference: "7186040733", IdempotencyKey: "return-effect-key-0002", Currency: "ARS", AmountMinorUnits: 125050}
	created, err := provider.Create(context.Background(), request)
	if err != nil || created.ProviderRefundReference != "1622029222" || created.ProviderStatus != "approved" {
		t.Fatalf("created=%+v err=%v", created, err)
	}
	request.ProviderRefundReference = created.ProviderRefundReference
	retrieved, err := provider.Retrieve(context.Background(), request)
	if err != nil || retrieved != created {
		t.Fatalf("retrieved=%+v err=%v", retrieved, err)
	}
}

func TestProvidersFailClosed(t *testing.T) {
	stripeProvider, _ := NewStripeProvider("sk_test")
	if _, err := stripeProvider.Create(context.Background(), Refund{Provider: "stripe", ProviderPaymentReference: "invented", IdempotencyKey: "return-effect-key-0003", Currency: "ARS", AmountMinorUnits: 1}); err == nil {
		t.Fatal("invalid Stripe payment reference accepted")
	}
	mpProvider, _ := NewMercadoPagoProvider("TEST")
	if _, err := mpProvider.Create(context.Background(), Refund{Provider: "mercado_pago", ProviderPaymentReference: "1", IdempotencyKey: "return-effect-key-0004", Currency: "USD", AmountMinorUnits: 1}); err == nil {
		t.Fatal("non-ARS Mercado Pago refund accepted without an exact exponent contract")
	}
}
````

### FILE: `return_refund_worker/internal/refundworker/processor_test.go`
```yaml
block_id: "GO-OFFICIAL-RETURN-REFUND:file:11"
operation: CREATE
provenance: AUTHORED
source: "local orchestration regressions"
license: "LicenseRef-Workspace-Owner"
sha256: "f4208fe5cbfb8c821c2c9e14c8b5aa34947052bf6e91f418dccc87b0524d94a5"
variables: []
secrets_allowed: false
```
````go
package refundworker

import (
	"context"
	"errors"
	"testing"
	"time"
)

type fakeStore struct {
	work       *Work
	refund     Refund
	prepareErr error
	completed  bool
	finished   bool
	outcome    string
	code       string
	result     ProviderResult
}

func (f *fakeStore) Claim(context.Context, string, string, time.Duration) (*Work, error) {
	return f.work, nil
}
func (f *fakeStore) Prepare(context.Context, Work, string) (Refund, error) {
	return f.refund, f.prepareErr
}
func (f *fakeStore) Complete(_ context.Context, _ Work, _ string, _ Refund, result ProviderResult, outcome, code string, _ time.Duration) error {
	f.completed, f.outcome, f.code, f.result = true, outcome, code, result
	return nil
}
func (f *fakeStore) Finish(_ context.Context, _ Work, _ string, outcome, code string, _ time.Duration) error {
	f.finished, f.outcome, f.code = true, outcome, code
	return nil
}

type fakeProvider struct {
	created, retrieved bool
	result             ProviderResult
	err                error
}

func (f *fakeProvider) Create(context.Context, Refund) (ProviderResult, error) {
	f.created = true
	return f.result, f.err
}
func (f *fakeProvider) Retrieve(context.Context, Refund) (ProviderResult, error) {
	f.retrieved = true
	return f.result, f.err
}

func baseRefund() Refund {
	return Refund{TenantID: "tenant", RequestID: "request", PaymentAttemptID: "payment", Provider: "stripe", ProviderPaymentReference: "pi_1", IdempotencyKey: "return-effect-key-0001", Currency: "ARS", AmountMinorUnits: 1000}
}

func TestProcessorCreatesAndCompletesSucceededRefund(t *testing.T) {
	work := &Work{TenantID: "tenant", RequestID: "request", Attempt: 1, ClaimToken: "claim"}
	store := &fakeStore{work: work, refund: baseRefund()}
	provider := &fakeProvider{result: ProviderResult{ProviderPaymentReference: "pi_1", ProviderRefundReference: "re_1", ProviderStatus: "succeeded", Currency: "ARS", AmountMinorUnits: 1000}}
	processor, err := NewProcessor(store, map[string]Provider{"stripe": provider}, "worker", time.Minute, time.Second, 5)
	if err != nil {
		t.Fatal(err)
	}
	if err = processor.Step(context.Background(), "claim"); err != nil {
		t.Fatal(err)
	}
	if !provider.created || provider.retrieved || !store.completed || store.outcome != "succeeded" || store.code != "" {
		t.Fatalf("provider=%+v store=%+v", provider, store)
	}
}

func TestProcessorRetrievesPendingRefundAndSchedulesRetry(t *testing.T) {
	work := &Work{TenantID: "tenant", RequestID: "request", Attempt: 2, ClaimToken: "claim"}
	refund := baseRefund()
	refund.ProviderRefundReference = "re_1"
	store := &fakeStore{work: work, refund: refund}
	provider := &fakeProvider{result: ProviderResult{ProviderPaymentReference: "pi_1", ProviderRefundReference: "re_1", ProviderStatus: "pending", Currency: "ARS", AmountMinorUnits: 1000}}
	processor, _ := NewProcessor(store, map[string]Provider{"stripe": provider}, "worker", time.Minute, time.Second, 5)
	if err := processor.Step(context.Background(), "claim"); err != nil {
		t.Fatal(err)
	}
	if provider.created || !provider.retrieved || store.outcome != "retry" || store.code != "PROVIDER_PENDING" {
		t.Fatalf("provider=%+v store=%+v", provider, store)
	}
}

func TestProcessorBlocksMappingProviderAndExhaustionFailures(t *testing.T) {
	work := &Work{TenantID: "tenant", RequestID: "request", Attempt: 5, ClaimToken: "claim"}
	store := &fakeStore{work: work, refund: baseRefund(), prepareErr: ErrMappingConflict}
	processor, _ := NewProcessor(store, nil, "worker", time.Minute, time.Second, 5)
	if err := processor.Step(context.Background(), "claim"); !errors.Is(err, ErrMappingConflict) || !store.finished || store.outcome != "blocked" {
		t.Fatalf("mapping err=%v store=%+v", err, store)
	}
	store = &fakeStore{work: work, refund: baseRefund()}
	processor, _ = NewProcessor(store, nil, "worker", time.Minute, time.Second, 5)
	if err := processor.Step(context.Background(), "claim"); err != nil || store.outcome != "blocked" || store.code != "PROVIDER_CONFIG_MISSING" {
		t.Fatalf("missing provider err=%v store=%+v", err, store)
	}
	store = &fakeStore{work: work, refund: baseRefund()}
	provider := &fakeProvider{err: errors.New("transport")}
	processor, _ = NewProcessor(store, map[string]Provider{"stripe": provider}, "worker", time.Minute, time.Second, 5)
	if err := processor.Step(context.Background(), "claim"); err == nil || store.outcome != "blocked" {
		t.Fatalf("exhausted err=%v store=%+v", err, store)
	}
}

func TestProcessorBlocksDivergentProviderResult(t *testing.T) {
	work := &Work{TenantID: "tenant", RequestID: "request", Attempt: 1, ClaimToken: "claim"}
	store := &fakeStore{work: work, refund: baseRefund()}
	provider := &fakeProvider{result: ProviderResult{ProviderPaymentReference: "pi_1", ProviderRefundReference: "re_1", ProviderStatus: "succeeded", Currency: "ARS", AmountMinorUnits: 999}}
	processor, _ := NewProcessor(store, map[string]Provider{"stripe": provider}, "worker", time.Minute, time.Second, 5)
	if err := processor.Step(context.Background(), "claim"); !errors.Is(err, ErrResponseMismatch) || !store.completed || store.outcome != "blocked" {
		t.Fatalf("mismatch err=%v store=%+v", err, store)
	}
}
````

### FILE: `return_refund_worker/internal/refundworker/postgres_integration_test.go`
```yaml
block_id: "GO-OFFICIAL-RETURN-REFUND:file:12"
operation: CREATE
provenance: AUTHORED
source: "local PostgreSQL 18.6 integration regression"
license: "LicenseRef-Workspace-Owner"
sha256: "f20113caa74aadd0a3fb28b94ecaf45c6cc49656dd2e7376bafd054b694e5d13"
variables: []
secrets_allowed: false
```
````go
package refundworker

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type succeedingProvider struct{}

func (succeedingProvider) Create(_ context.Context, value Refund) (ProviderResult, error) {
	return ProviderResult{ProviderPaymentReference: value.ProviderPaymentReference, ProviderRefundReference: "provider-" + value.RequestID, ProviderStatus: "succeeded", Currency: value.Currency, AmountMinorUnits: value.AmountMinorUnits}, nil
}
func (succeedingProvider) Retrieve(_ context.Context, value Refund) (ProviderResult, error) {
	return ProviderResult{ProviderPaymentReference: value.ProviderPaymentReference, ProviderRefundReference: value.ProviderRefundReference, ProviderStatus: "succeeded", Currency: value.Currency, AmountMinorUnits: value.AmountMinorUnits}, nil
}

// This legacy compatibility fixture bypasses relational triggers when seeding
// and removing synthetic data. It does not validate the originating business
// workflow. Require a separate disposable database, never a project/audit DB.
func refundTestConfig(raw string) (*pgxpool.Config, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, errors.New("explicit disposable database is required")
	}
	cfg, err := pgxpool.ParseConfig(raw)
	if err != nil || cfg.ConnConfig.Host != "127.0.0.1" || !strings.HasPrefix(cfg.ConnConfig.Database, "elite_refund_test_") || len(strings.TrimPrefix(cfg.ConnConfig.Database, "elite_refund_test_")) < 16 {
		return nil, errors.New("requires a dedicated disposable loopback elite_refund_test_<unique> database")
	}
	for _, fallback := range cfg.ConnConfig.Fallbacks {
		if fallback.Host != "127.0.0.1" {
			return nil, errors.New("non-loopback fallback forbidden")
		}
	}
	return cfg, nil
}

func TestRefundDatabaseGuard(t *testing.T) {
	for _, value := range []string{"", "invalid", "postgres://u@remote/elite_refund_test_0123456789abcdef", "postgres://u@localhost/elite_refund_test_0123456789abcdef", "postgres://u@127.0.0.1/production", "postgres://u@127.0.0.1/elite_confirmation_audit", "postgres://u@127.0.0.1/elite_refund_test_short", "host=127.0.0.1,remote dbname=elite_refund_test_0123456789abcdef"} {
		if cfg, err := refundTestConfig(value); err == nil || cfg != nil {
			t.Fatal("unsafe target accepted")
		}
	}
	if cfg, err := refundTestConfig("postgres://u@127.0.0.1/elite_refund_test_0123456789abcdef"); err != nil || cfg == nil {
		t.Fatal("disposable target rejected", err)
	}
}

func TestPostgresRefundLineAllocationPartialThenFullAndAmbiguity(t *testing.T) {
	testRefundLineAllocationPartialThenFullAndAmbiguity(t, "stripe", succeedingProvider{})
}

func testRefundLineAllocationPartialThenFullAndAmbiguity(t *testing.T, sourceProvider string, provider Provider) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	cfg, err := refundTestConfig(url)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := "018f4d4a-7b36-7a21-8d10-2f4c54c2f201"
	cleanup := func() {
		tx, cleanupErr := pool.Begin(ctx)
		if cleanupErr != nil {
			return
		}
		defer tx.Rollback(ctx)
		_, _ = tx.Exec(ctx, `set local session_replication_role=replica`)
		for _, table := range []string{"payment.return_refund_observation", "payment.return_refund", "sales.return_effect_attempt", "sales.return_effect_execution", "sales.return_effect_request", "sales.return_disposition", "sales.return_receipt", "sales.return_authorization", "payment.payment_attempt", "sales.customer_order_line", "sales.customer_order", "platform.outbox_event", "inventory.stock_unit", "org.organization", "platform.tenant"} {
			_, _ = tx.Exec(ctx, `delete from `+table+` where tenant_id=$1`, tenant)
		}
		_ = tx.Commit(ctx)
	}
	cleanup()
	defer cleanup()

	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback(ctx)
	if _, err = tx.Exec(ctx, `set local session_replication_role=replica`); err != nil {
		t.Fatal(err)
	}
	fixtures := []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,'return-refund','Return Refund','Return Refund')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,'store','return-refund-store','Store','store')`,
		`insert into inventory.stock_unit(tenant_id,stock_unit_id,organization_id,variant_id,serial_number,vin,battery_serial_number,state,version,received_at) values($1,'stock-1','store','variant','SERIAL-1','VIN-1','BATTERY-1','sold',1,clock_timestamp()),($1,'stock-2','store','variant','SERIAL-2','VIN-2','BATTERY-2','sold',1,clock_timestamp()),($1,'stock-3','store','variant','SERIAL-3','VIN-3','BATTERY-3','sold',1,clock_timestamp())`,
		`insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version) values($1,'order-1','store','customer','delivered','ARS',2000,3),($1,'order-2','store','customer','delivered','ARS',1000,3)`,
		`insert into sales.customer_order_line(tenant_id,order_id,line_id,variant_id,quantity,unit_price_minor_units,allocated_stock_unit_id) values($1,'order-1','line-1','variant',1,1000,'stock-1'),($1,'order-1','line-2','variant',1,1000,'stock-2'),($1,'order-2','line-3','variant',1,1000,'stock-3')`,
		`insert into payment.payment_attempt(tenant_id,payment_attempt_id,order_id,provider_code,provider_reference,idempotency_key,state,currency,amount_minor_units,version) values($1,'payment-1','order-1','stripe','pi_payment_1','payment-key-000001','captured','ARS',2000,3),($1,'payment-2a','order-2','stripe','pi_payment_2a','payment-key-000002','captured','ARS',1000,3),($1,'payment-2b','order-2','stripe','pi_payment_2b','payment-key-000003','captured','ARS',1000,3)`,
		`insert into sales.return_authorization(tenant_id,authorization_id,organization_id,exception_id,handover_id,order_id,stock_unit_id,customer_principal_id,disposition,exact_cost_source_order_id,state,authorized_by_subject) values($1,'auth-1','store','exception-1','handover-1','order-1','stock-1','customer','return','order-1','authorized','operator'),($1,'auth-2','store','exception-2','handover-2','order-1','stock-2','customer','return','order-1','authorized','operator'),($1,'auth-3','store','exception-3','handover-3','order-2','stock-3','customer','return','order-2','authorized','operator')`,
		`insert into sales.return_receipt(tenant_id,receipt_id,authorization_id,organization_id,order_id,stock_unit_id,customer_principal_id,received_serial_number,condition_code,notes,evidence_sha256_hex,received_by_subject) values($1,'receipt-1','auth-1','store','order-1','stock-1','customer','SERIAL-1','opened','received','aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa','operator'),($1,'receipt-2','auth-2','store','order-1','stock-2','customer','SERIAL-2','opened','received','bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb','operator'),($1,'receipt-3','auth-3','store','order-2','stock-3','customer','SERIAL-3','opened','received','cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc','operator')`,
		`insert into sales.return_disposition(tenant_id,disposition_id,receipt_id,inventory_action,customer_remedy,notes,decided_by_subject) values($1,'disposition-1','receipt-1','restock','refund','approved','operator'),($1,'disposition-2','receipt-2','restock','refund','approved','operator'),($1,'disposition-3','receipt-3','restock','refund','approved','operator')`,
	}
	for _, query := range fixtures {
		if _, err = tx.Exec(ctx, query, tenant); err != nil {
			t.Fatal(err)
		}
	}
	if sourceProvider != "stripe" {
		if _, err = tx.Exec(ctx, `update payment.payment_attempt set provider_code=$2,provider_reference=case payment_attempt_id when 'payment-1' then '7186040733' when 'payment-2a' then '7186040734' else '7186040735' end where tenant_id=$1`, tenant, sourceProvider); err != nil {
			t.Fatal(err)
		}
	}
	if err = tx.Commit(ctx); err != nil {
		t.Fatal(err)
	}
	if _, err = pool.Exec(ctx, `insert into sales.return_effect_request(tenant_id,request_id,disposition_id,effect_kind,owner_context,state,idempotency_key) values($1,'refund-1','disposition-1','refund','payment','requested','return-refund-key-0001'),($1,'refund-2','disposition-2','refund','payment','requested','return-refund-key-0002'),($1,'refund-3','disposition-3','refund','payment','requested','return-refund-key-0003')`, tenant); err != nil {
		t.Fatal(err)
	}

	store := NewPostgresStore(pool)
	refundProvider := sourceProvider
	if sourceProvider == "mercadopago" {
		refundProvider = "mercado_pago"
	}
	processor, err := NewProcessor(store, map[string]Provider{refundProvider: provider}, "refund-worker", time.Minute, time.Millisecond, 5)
	if err != nil {
		t.Fatal(err)
	}
	if err = processor.Step(ctx, "claim-refund-1"); err != nil {
		t.Fatal(err)
	}
	var paymentState, executionState string
	var paymentVersion int64
	var observations, events int
	if err = pool.QueryRow(ctx, `select state,version from payment.payment_attempt where tenant_id=$1 and payment_attempt_id='payment-1'`, tenant).Scan(&paymentState, &paymentVersion); err != nil {
		t.Fatal(err)
	}
	if paymentState != "captured" || paymentVersion != 3 {
		t.Fatalf("partial refund incorrectly closed payment: state=%s version=%d", paymentState, paymentVersion)
	}
	// The boundary alias must never rewrite historical payment identity or keys.
	var storedProvider, paymentKey string
	if err = pool.QueryRow(ctx, `select provider_code,idempotency_key from payment.payment_attempt where tenant_id=$1 and payment_attempt_id='payment-1'`, tenant).Scan(&storedProvider, &paymentKey); err != nil || storedProvider != sourceProvider || paymentKey != "payment-key-000001" {
		t.Fatalf("source identity changed: provider=%s key=%s err=%v", storedProvider, paymentKey, err)
	}
	prepared, err := scanRefund(pool.QueryRow(ctx, `select tenant_id,request_id,payment_attempt_id,order_id,line_id,stock_unit_id,provider_code,provider_payment_reference,idempotency_key,currency,amount_minor_units,coalesce(provider_refund_reference,''),coalesce(provider_status,'') from payment.return_refund where tenant_id=$1 and request_id='refund-1'`, tenant))
	if err != nil || prepared.Provider != refundProvider || prepared.PaymentAttemptID != "payment-1" || prepared.RequestID != "refund-1" || prepared.IdempotencyKey != "return-refund-key-0001" || prepared.AmountMinorUnits != 1000 {
		t.Fatalf("refund identity or contract changed: %+v err=%v", prepared, err)
	}
	retrieved, err := provider.Retrieve(ctx, prepared)
	if err != nil || validateResult(prepared, retrieved) != nil || retrieved.ProviderRefundReference != prepared.ProviderRefundReference {
		t.Fatalf("persisted refund could not reconcile using its unchanged identity: %+v err=%v", retrieved, err)
	}
	if err = processor.Step(ctx, "claim-refund-2"); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select state,version from payment.payment_attempt where tenant_id=$1 and payment_attempt_id='payment-1'`, tenant).Scan(&paymentState, &paymentVersion); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select status from sales.return_effect_execution where tenant_id=$1 and request_id='refund-2'`, tenant).Scan(&executionState); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select count(*) from payment.return_refund_observation where tenant_id=$1 and request_id in ('refund-1','refund-2')`, tenant).Scan(&observations); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select count(*) from platform.outbox_event where tenant_id=$1 and event_type in ('return-refund.succeeded','payment.refunded')`, tenant).Scan(&events); err != nil {
		t.Fatal(err)
	}
	if paymentState != "refunded" || paymentVersion != 4 || executionState != "succeeded" || observations != 2 || events != 3 {
		t.Fatalf("state=%s version=%d execution=%s observations=%d events=%d", paymentState, paymentVersion, executionState, observations, events)
	}
	if err = processor.Step(ctx, "claim-refund-3"); !errors.Is(err, ErrMappingConflict) {
		t.Fatalf("ambiguous payment mapping accepted: %v", err)
	}
	if err = pool.QueryRow(ctx, `select status from sales.return_effect_execution where tenant_id=$1 and request_id='refund-3'`, tenant).Scan(&executionState); err != nil {
		t.Fatal(err)
	}
	if executionState != "blocked" {
		t.Fatalf("ambiguous mapping status=%s", executionState)
	}
	if _, err = pool.Exec(ctx, `update payment.return_refund_observation set provider_status='invented' where tenant_id=$1 and request_id='refund-1'`, tenant); err == nil {
		t.Fatal("immutable refund observation was mutated")
	}
	if work, claimErr := store.Claim(ctx, "refund-worker", "claim-extra", time.Minute); claimErr != nil || work != nil {
		t.Fatalf("terminal refund work reclaimed: %+v %v", work, claimErr)
	}
}
````

### FILE: `return_refund_worker/cmd/return-refund-worker/main.go`
```yaml
block_id: "GO-OFFICIAL-RETURN-REFUND:file:13"
operation: CREATE
provenance: AUTHORED
source: "local runtime wiring"
license: "LicenseRef-Workspace-Owner"
sha256: "5ac3b6c37645bacc7dae27f6848fcde3905fdaf347e12a613586716de7012f43"
variables: []
secrets_allowed: false
```
````go
package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"elite.local/return-refund-worker/internal/operationaltelemetry"
	"elite.local/return-refund-worker/internal/refundworker"
	"github.com/jackc/pgx/v5/pgxpool"
)

func claimToken() (string, error) {
	value := make([]byte, 16)
	if _, err := rand.Read(value); err != nil {
		return "", err
	}
	return hex.EncodeToString(value), nil
}

func run() (result error) {
	databaseURL := os.Getenv("DATABASE_URL")
	workerID := os.Getenv("REFUND_WORKER_ID")
	if databaseURL == "" || workerID == "" {
		return errors.New("DATABASE_URL and REFUND_WORKER_ID are required")
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	ctx, cleanup, err := nativeStopContext(ctx, os.Getenv("ELITE_STOP_EVENT_HANDLE"))
	if err != nil {
		return err
	}
	defer func() {
		if cleanupErr := cleanup(); cleanupErr != nil && result == nil {
			result = cleanupErr
		}
	}()
	var reporter *operationaltelemetry.Reporter
	config, digest := os.Getenv("REFUND_TELEMETRY_CONFIG"), os.Getenv("REFUND_TELEMETRY_CONFIG_SHA256")
	if config != "" || digest != "" {
		reporter, err = operationaltelemetry.FromFile(config, digest)
		if err != nil {
			return err
		}
		defer reporter.Close()
	}
	observe := func(ctx context.Context, outcome string, elapsed time.Duration) error {
		if reporter == nil {
			return nil
		}
		return reporter.Observe(ctx, outcome, elapsed)
	}
	if ctx.Err() != nil {
		return nil
	}
	if err = observe(ctx, "started", 0); err != nil {
		return err
	}
	defer func() {
		outcome := "stopped"
		if result != nil {
			outcome = "error"
		}
		// A separate finite reporting budget records cooperative stop after parent
		// cancellation. It never retries an observation whose delivery is unknown.
		if !errors.Is(result, operationaltelemetry.ErrUnavailable) {
			if reportErr := observe(context.Background(), outcome, 0); reportErr != nil {
				result = reportErr
			}
		}
	}()
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()
	if err = pool.Ping(ctx); err != nil {
		return err
	}
	providers := map[string]refundworker.Provider{}
	if secret := os.Getenv("STRIPE_SECRET_KEY"); secret != "" {
		providers["stripe"], err = refundworker.NewStripeProvider(secret)
		if err != nil {
			return err
		}
	}
	if token := os.Getenv("MERCADO_PAGO_ACCESS_TOKEN"); token != "" {
		providers["mercado_pago"], err = refundworker.NewMercadoPagoProvider(token)
		if err != nil {
			return err
		}
	}
	processor, err := refundworker.NewProcessor(refundworker.NewPostgresStore(pool), providers, workerID, 2*time.Minute, 30*time.Second, 20)
	if err != nil {
		return err
	}
	if err = runObservedRefundLoop(ctx, processor.Step, claimToken, observe); err != nil {
		return err
	}
	return nativeStopFailure(ctx)
}

// The real host stops claiming work when its error report cannot be written.
// The existing durable processor remains responsible for reconciliation. This
// does not retry the uncertain report or change a refund result after commit.
func runRefundLoop(ctx context.Context, step func(context.Context, string) error, nextToken func() (string, error)) error {
	return runObservedRefundLoop(ctx, step, nextToken, nil)
}

// handled describes Step returning nil; the durable domain record, not this
// operational outcome, determines whether money moved or a refund succeeded.
func runObservedRefundLoop(ctx context.Context, step func(context.Context, string) error, nextToken func() (string, error), observe func(context.Context, string, time.Duration) error) error {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		if ctx.Err() != nil {
			return nil
		}
		token, tokenErr := nextToken()
		if tokenErr != nil {
			return tokenErr
		}
		began := time.Now()
		stepErr := step(ctx, token)
		if reportErr := logStepError(stepErr); reportErr != nil {
			return reportErr
		}
		if observe != nil {
			outcome := "handled"
			if errors.Is(stepErr, refundworker.ErrNoWork) {
				outcome = "idle"
			} else if stepErr != nil {
				outcome = "error"
			}
			if reportErr := observe(ctx, outcome, time.Since(began)); reportErr != nil {
				return reportErr
			}
		}
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
		}
	}
}

func main() {
	exitOnHostFailure(run())
}

func exitOnHostFailure(err error) {
	if err != nil {
		log.Fatal("return refund host failed: REFUND_HOST_FAILED")
	}
}

// Log only host-owned operational outcomes. Durable results own diagnostics.
var errRefundReport = errors.New("return refund report unavailable: REFUND_REPORT_FAILED")

type checkedRefundLogWriter struct{ writer io.Writer }

func (w checkedRefundLogWriter) Write(p []byte) (int, error) {
	if w.writer == nil {
		return 0, errRefundReport
	}
	n, err := w.writer.Write(p)
	if n < 0 || n > len(p) {
		return 0, errRefundReport
	}
	if err != nil || n != len(p) {
		return n, errRefundReport
	}
	return n, nil
}

// Preserve the host logger's configured prefix/flags and private fixed message,
// but surface delivery loss. The host loop emits serially. A blocking writer
// still requires target-owned process supervision; no I/O deadline is claimed.
// Full local Write acceptance does not prove log retention or alert delivery.
func logStepError(err error) error {
	if err != nil && !errors.Is(err, refundworker.ErrNoWork) {
		logger := log.New(checkedRefundLogWriter{log.Writer()}, log.Prefix(), log.Flags())
		if logger.Output(2, "return refund step failed: REFUND_STEP_FAILED") != nil {
			return errRefundReport
		}
	}
	return nil
}
````

### FILE: `return_refund_worker/README.md`
```yaml
block_id: "GO-OFFICIAL-RETURN-REFUND:file:14"
operation: CREATE
provenance: AUTHORED
source: "local operator documentation"
license: "LicenseRef-Workspace-Owner"
sha256: "8236d2d1174055f254389735d1df35c5a542b7e481024c1b47cffd02e9e6ad30"
variables: []
secrets_allowed: false
```
````markdown
# Official return refund worker

## 0.1.4 — observable host report failure

The real run loop now checks log delivery before claiming another job. A failed,
short or invalid step-log write returns static REFUND_REPORT_FAILED; the owning
main exits1 via its existing private REFUND_HOST_FAILED boundary. Each step log
makes one Write, with no retry of the uncertain record and no formatting of the
underlying step/sink error. The terminal host message is a separate best-effort
write; exit status remains available even when stderr is unusable.

Successful reports preserve the configured logger prefix/flags and fixed
REFUND_STEP_FAILED message. Nil/ErrNoWork still emit nothing; ordinary step
errors continue under the existing polling policy when their report succeeds.
An already-cancelled loop does not generate a token or start another step.

This is AUTHORED runtime wiring, not a provider/financial-policy change. Durable
refund reconciliation, idempotency, amounts, currencies and results are unchanged.
The owner must supervise the exit status and diagnose durable results before a
restart; this patch does not install a supervisor, endpoint, alert or retry daemon.
A blocking writer still needs target-owned process supervision. Complete local
Write acceptance does not prove retention, remote acknowledgement or alerting.

Tests cover six failed-write modes against exact0.1.3 before the fix, the actual
host loop, a real closed file, privacy, healthy continuation/idle/cancellation,
and a subprocess exiting1 with an unusable report sink. Finite fuzz covers write
counts/error combinations. No provider credentials, payment calls or live effects
are used. Source contracts: https://pkg.go.dev/log#Logger.Output and
https://pkg.go.dev/io#Writer; executed toolchain stays Go1.26.7.

This standalone Go worker claims only `refund/payment` return effects from the enterprise PostgreSQL schema. It resolves the exact allocated order line and one unambiguous captured payment, creates a partial refund through the exact official Stripe or Mercado Pago SDK, persists immutable provider observations, retrieves non-final refunds, and completes the effect only after an exact final provider result.

Required runtime values are `DATABASE_URL` and `REFUND_WORKER_ID`; add `STRIPE_SECRET_KEY` and/or `MERCADO_PAGO_ACCESS_TOKEN` only through the deployment secret manager. A missing provider credential blocks that provider's effect without exposing a secret. The Mercado Pago adapter is deliberately ARS/two-decimal until another market has an admitted currency-exponent contract.

Production admission still requires provider account/country enablement, sandbox and live contractual tests, webhook delivery into the durable inbox, reconciliation alerts, business approval of line-level refund allocation, and load/security/recovery evidence for the selected target.

Host stderr reports only REFUND_HOST_FAILED (startup exit1) or REFUND_STEP_FAILED (step failure; inspect the durable result). Arbitrary database/provider errors are not formatted in host logs. Idle results do not emit failure messages. This does not prove target log delivery, retention or alerts.


## Optional closed-outcome operational telemetry

Set both REFUND_TELEMETRY_CONFIG and REFUND_TELEMETRY_CONFIG_SHA256 to an exact
local JSON config accepted by operationaltelemetry.Config. Pin CA, client
certificate and private-key bytes individually. TLS1.3 mutual authentication,
no proxy or redirects, response and request byte caps, and one shared finite
observation deadline are required. With neither variable, the previous local
logger behavior remains. A required-observability project must configure and
verify this path explicitly before its own admission; absence is not readiness.

Outcomes are started, idle, handled, error and stopped. handled only means Step
returned nil; the database and provider receipts own the refund result. No DSN,
customer, claim token, provider response or free-form error enters OTLP. Export
uses cumulative metrics and a correlated log/span. An unknown delivery stops
new claims without retrying the observation or changing a committed refund.
The HTTP acknowledgement is not an fsync or business-commit acknowledgement.

The inherited ELITE_STOP_EVENT_HANDLE protocol is optional and Windows-only.
Only a trusted launcher may supply the manual-reset event. Native protocol
failure remains a host failure; the supervisor receipt owns process status.
This does not isolate hostile same-user workers or install a production service.

cmd/telemetry-reference supplies a finite synthetic load probe and ephemeral
certificate fixture for reference_telemetry/run_reference_runtime.py. It does
not call the database, provider or refund processor. Certificate keys are for
that short-lived reference only; the CA private key is never saved. The full
reference runner requires a hash-bound runtime lock, exact migration identities
and a fresh owned PostgreSQL cluster. See the Windows reference telemetry pack.
````

### FILE: `return_refund_worker/cmd/return-refund-worker/main_test.go`
```yaml
block_id: "GO-OFFICIAL-RETURN-REFUND:file:15"
operation: CREATE
provenance: AUTHORED
source: "local host privacy regression and subprocess fixture"
license: "LicenseRef-Workspace-Owner"
sha256: "98c991d3ba6cce7ed03a6a7d0dae411cb599822a8064ab17056e0616e0fb8448"
variables: []
secrets_allowed: false
```
````go
package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"testing"
	"time"

	"elite.local/return-refund-worker/internal/refundworker"
)

type privateHostError struct{ formatted *bool }

func (e privateHostError) Error() string {
	*e.formatted = true
	return "PRIVATE-SYNTHETIC-STEP token=synthetic-secret"
}

func TestRefundHostStepErrorsDoNotFormatPrivateCause(t *testing.T) {
	var output bytes.Buffer
	writer, flags, prefix := log.Writer(), log.Flags(), log.Prefix()
	log.SetOutput(&output)
	log.SetFlags(0)
	log.SetPrefix("")
	defer func() { log.SetOutput(writer); log.SetFlags(flags); log.SetPrefix(prefix) }()
	formatted := false
	logStepError(privateHostError{&formatted})
	if formatted {
		t.Error("host formatted the private error cause")
	}
	if strings.Contains(output.String(), "PRIVATE-SYNTHETIC") || strings.Contains(output.String(), "synthetic-secret") {
		t.Error("private marker escaped into step log")
	}
	if output.String() != "return refund step failed: REFUND_STEP_FAILED\n" {
		t.Error("step failure requires the fixed operator code")
	}
	output.Reset()
	for _, err := range []error{nil, refundworker.ErrNoWork, fmt.Errorf("synthetic wrapper: %w", refundworker.ErrNoWork)} {
		logStepError(err)
	}
	if output.Len() != 0 {
		t.Error("idle results emitted failure logs")
	}
}

func TestRefundHostStartupErrorRedaction(t *testing.T) {
	if os.Getenv("ELITE_REFUND_TEST_CHILD") == "1" {
		main()
		return
	}
	for _, tc := range []struct{ name, dsn string }{
		{"missing_configuration", ""},
		{"invalid_private_dsn", "postgres://synthetic@127.0.0.1:invalid/PRIVATE-SYNTHETIC-DATABASE"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestRefundHostStartupErrorRedaction$")
			for _, key := range []string{"SystemRoot", "WINDIR", "TEMP", "TMP", "PATH"} {
				if value, ok := os.LookupEnv(key); ok {
					cmd.Env = append(cmd.Env, key+"="+value)
				}
			}
			cmd.Env = append(cmd.Env, "ELITE_REFUND_TEST_CHILD=1", "DATABASE_URL="+tc.dsn, "REFUND_WORKER_ID=synthetic-worker")
			output, err := cmd.CombinedOutput()
			var exit *exec.ExitError
			if !errors.As(err, &exit) || exit.ExitCode() != 1 {
				t.Fatal("failed startup must exit1")
			}
			if strings.Contains(string(output), "PRIVATE-SYNTHETIC") {
				t.Error("startup disclosed private configuration")
			}
			if !strings.Contains(string(output), "REFUND_HOST_FAILED") {
				t.Error("startup requires fixed operator code")
			}
		})
	}
}

type privateSinkError struct{}

func (privateSinkError) Error() string { panic("sink error must remain private") }

type failedHostWriter struct {
	mode  string
	calls int
}

func (w *failedHostWriter) Write(p []byte) (int, error) {
	w.calls++
	switch w.mode {
	case "error":
		return 0, privateSinkError{}
	case "full-error":
		return len(p), privateSinkError{}
	case "zero":
		return 0, nil
	case "short":
		return len(p) - 1, nil
	case "negative":
		return -1, nil
	case "oversized":
		return len(p) + 1, nil
	}
	return len(p), nil
}

func TestRefundHostReportDeliveryFailureIsVisible(t *testing.T) {
	for _, mode := range []string{"error", "full-error", "zero", "short", "negative", "oversized"} {
		t.Run(mode, func(t *testing.T) {
			previous := log.Writer()
			defer log.SetOutput(previous)
			w := &failedHostWriter{mode: mode}
			log.SetOutput(w)
			// Reflection permits the same behavioral probe against the old void
			// function and its new error-returning contract without editing baseline.
			result := reflect.ValueOf(logStepError).Call([]reflect.Value{reflect.ValueOf(errors.New("PRIVATE-STEP"))})
			if len(result) != 1 || result[0].IsNil() {
				t.Fatal("failed reporting has no observable error for host")
			}
			err, ok := result[0].Interface().(error)
			if !ok || err.Error() != "return refund report unavailable: REFUND_REPORT_FAILED" {
				t.Fatal("non-static report failure")
			}
			if w.calls != 1 {
				t.Fatal("uncertain report was retried")
			}
		})
	}
}

var _ io.Writer = (*failedHostWriter)(nil)

func TestRefundHostLoopStopsBeforeAnotherClaimOnReportLoss(t *testing.T) {
	for _, mode := range []string{"error", "full-error", "zero", "short", "negative", "oversized"} {
		t.Run(mode, func(t *testing.T) {
			previous := log.Writer()
			defer log.SetOutput(previous)
			w := &failedHostWriter{mode: mode}
			log.SetOutput(w)
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			steps, tokens := 0, 0
			formatted := false
			err := runRefundLoop(ctx, func(context.Context, string) error { steps++; return privateHostError{&formatted} }, func() (string, error) { tokens++; return "synthetic-token", nil })
			if err != errRefundReport || steps != 1 || tokens != 1 || w.calls != 1 || formatted {
				t.Fatal("report failure did not stop the actual host loop privately")
			}
		})
	}
}

func TestRefundHostLoopRetainsHealthyStepAndIdlePolicy(t *testing.T) {
	var output bytes.Buffer
	writer, flags, prefix := log.Writer(), log.Flags(), log.Prefix()
	defer func() { log.SetOutput(writer); log.SetFlags(flags); log.SetPrefix(prefix) }()
	log.SetOutput(&output)
	log.SetFlags(0)
	log.SetPrefix("approved-static: ")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	steps, tokens := 0, 0
	formatted := false
	err := runRefundLoop(ctx, func(context.Context, string) error {
		steps++
		if steps == 1 {
			return privateHostError{&formatted}
		}
		cancel()
		return refundworker.ErrNoWork
	}, func() (string, error) { tokens++; return "synthetic-token", nil })
	if err != nil || steps != 2 || tokens != 2 || formatted {
		t.Fatal("healthy reporting changed step/idle behavior")
	}
	if output.String() != "approved-static: return refund step failed: REFUND_STEP_FAILED\n" {
		t.Fatal("fixed message or configured prefix changed")
	}
}

func TestRefundHostLoopCancellationAndTokenFailure(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if runRefundLoop(ctx, func(context.Context, string) error { t.Fatal("step after cancellation"); return nil }, func() (string, error) { t.Fatal("token after cancellation"); return "", nil }) != nil {
		t.Fatal("cancellation is not normal shutdown")
	}
	private := errors.New("PRIVATE-TOKEN-CAUSE")
	if err := runRefundLoop(context.Background(), func(context.Context, string) error { t.Fatal("step after failed token"); return nil }, func() (string, error) { return "", private }); err != private {
		t.Fatal("token failure not returned")
	}
}

func TestRefundHostRealClosedFileStopsLoop(t *testing.T) {
	p, err := os.CreateTemp(t.TempDir(), "closed-host-log-")
	if err != nil {
		t.Fatal(err)
	}
	name := p.Name()
	if p.Close() != nil {
		t.Fatal("close fixture")
	}
	previous := log.Writer()
	defer log.SetOutput(previous)
	log.SetOutput(p)
	steps := 0
	err = runRefundLoop(context.Background(), func(context.Context, string) error { steps++; return errors.New("PRIVATE-STEP") }, func() (string, error) { return "synthetic-token", nil })
	if err != errRefundReport || steps != 1 {
		t.Fatal("closed file did not stop host")
	}
	data, err := os.ReadFile(name)
	if err != nil || len(data) != 0 {
		t.Fatal("closed sink unexpectedly changed")
	}
}

func TestRefundHostLostReportsExitNonzeroWithoutPrivateCause(t *testing.T) {
	if mode := os.Getenv("ELITE_REFUND_TEST_REPORT_CHILD"); mode != "" {
		if mode == "healthy" {
			exitOnHostFailure(nil)
			_, _ = os.Stdout.WriteString("HOST_NORMAL_RETURN\n")
			return
		}
		// The same loop and exit adapter are called by production run/main. This
		// child substitutes only the already-completed step and failed log sink.
		log.SetOutput(&failedHostWriter{mode: "error"})
		err := runRefundLoop(context.Background(), func(context.Context, string) error { return errors.New("PRIVATE-STEP") }, func() (string, error) { return "synthetic-token", nil })
		exitOnHostFailure(err)
		_, _ = os.Stdout.WriteString("UNEXPECTED_HOST_CONTINUATION\n")
		return
	}
	for _, mode := range []string{"failed", "healthy"} {
		t.Run(mode, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestRefundHostLostReportsExitNonzeroWithoutPrivateCause$")
			for _, key := range []string{"SystemRoot", "WINDIR", "TEMP", "TMP", "PATH"} {
				if v, ok := os.LookupEnv(key); ok {
					cmd.Env = append(cmd.Env, key+"="+v)
				}
			}
			cmd.Env = append(cmd.Env, "ELITE_REFUND_TEST_REPORT_CHILD="+mode)
			output, err := cmd.CombinedOutput()
			if strings.Contains(string(output), "PRIVATE") || strings.Contains(string(output), "UNEXPECTED_HOST_CONTINUATION") {
				t.Fatal("host continued or private text escaped")
			}
			if mode == "failed" {
				var exit *exec.ExitError
				if !errors.As(err, &exit) || exit.ExitCode() != 1 {
					t.Fatal("report loss must exit1 even when stderr cannot report")
				}
			} else if err != nil || !strings.Contains(string(output), "HOST_NORMAL_RETURN") {
				t.Fatal("normal exit changed")
			}
		})
	}
}

type countHostWriter struct {
	count int
	fail  bool
	calls int
}

func (w *countHostWriter) Write([]byte) (int, error) {
	w.calls++
	if w.fail {
		return w.count, privateSinkError{}
	}
	return w.count, nil
}

func FuzzRefundHostCheckedLogWrite(f *testing.F) {
	for _, n := range []int{-1, 0, 1, 6, 7, 8, 1024} {
		f.Add(n, false)
		f.Add(n, true)
	}
	f.Fuzz(func(t *testing.T, n int, failed bool) {
		w := &countHostWriter{count: n, fail: failed}
		actual, err := (checkedRefundLogWriter{writer: w}).Write([]byte("report\n"))
		if w.calls != 1 {
			t.Fatal("log write repeated")
		}
		if !failed && n == 7 {
			if err != nil || actual != 7 {
				t.Fatal("complete write rejected")
			}
			return
		}
		if err != errRefundReport || actual < 0 || actual > 7 {
			t.Fatal("invalid delivery accepted or count escaped")
		}
	})
}

func TestRefundHostIdleDoesNotTouchFailedSink(t *testing.T) {
	w := &failedHostWriter{mode: "error"}
	previous := log.Writer()
	defer log.SetOutput(previous)
	log.SetOutput(w)
	for _, err := range []error{nil, refundworker.ErrNoWork} {
		if logStepError(err) != nil {
			t.Fatal("idle report error")
		}
	}
	if w.calls != 0 {
		t.Fatal("idle wrote to sink")
	}
	if n, err := (checkedRefundLogWriter{}).Write([]byte("report")); n != 0 || err != errRefundReport {
		t.Fatal("nil writer accepted")
	}
}

var _ io.Writer = (*countHostWriter)(nil)
````

## 6. Configuration surface

`DATABASE_URL` and `REFUND_WORKER_ID` are mandatory runtime inputs. `STRIPE_SECRET_KEY` and `MERCADO_PAGO_ACCESS_TOKEN` are optional provider selections and must be resolved by the target secret manager. Fixed safe defaults are a two-minute lease, thirty-second provider reconciliation delay, twenty attempts and one-second idle poll; a selected project must size them against measured API latency, provider rate limits and termination budgets. Mercado Pago is admitted only for ARS/two decimals in this pack.

## 7. Dependency bill

| Dependency | Fixed baseline | Purpose | License/authority |
|---|---:|---|---|
| Go | 1.26.7 | runtime/tests | Go project |
| PostgreSQL | 18.6 | locks, projection, immutable evidence, transaction/outbox | PostgreSQL License |
| pgx | 5.10.0 | PostgreSQL driver/pool | MIT |
| Stripe Go SDK | 86.3.0, commit above | refund create/retrieve, types, idempotency header | MIT |
| Mercado Pago Go SDK | 1.14.0, commit above | payment/refund get/create and idempotency header | MIT |
| Existing return/order/payment schemas | compatible packs above | authoritative mapping and effect execution | workspace owner |

## 8. Apply order

Compose after backend migrations 0001–0019 and the official payment adapter/source locks. Apply 0020 and its SQL test, then build/test the standalone module. Configure exactly the providers selected by the project. Deploy only after provider sandbox create/retrieve/webhook reconciliation, target secrets, alerting and approved refund allocation rules pass. Stop claims and preserve audit rows for rollback; use the down migration only in a non-production rehearsal or after an explicit data-retention decision.

## 9. Verification

Required gates: exact materialization and hashes; frozen module resolution; `gofmt`; full unit/contract tests; `go vet`; binary build; clean PostgreSQL 18.6 migrations 0001–0020; 0020 SQL test; down/absence/up; real database integration proving line-derived amounts, partial-then-full payment transition, outbox, immutable observations, ambiguity block and no terminal reclaim. Project production admission additionally requires selected provider sandbox/live contract, webhook/reconciliation, target load/security/recovery/observability and business/fiscal/accounting acceptance.

## 10. Reconstruction evidence

V142 evidence records exact SDK pins, official API contracts, byte reconstruction, Go gates and PostgreSQL 18.6 integration. It admits the reusable refund owner component only. It does not claim live provider credentials/account approval, refund policy correctness for discounts/tax/shipping not represented by the current line model, accounting reversal, ARCA credit note or production readiness.

### FILE: `return_refund_worker/internal/operationaltelemetry/reporter.go`
```yaml
block_id: "REFUND-TELEMETRY:reporter:v1"
operation: CREATE
provenance: AUTHORED
source: "local implementation; official OTLP HTTP JSON and Microsoft native lifecycle contracts; V372 exact integration evidence"
license: "LicenseRef-Workspace-Owner"
sha256: "5126360538dcb26a20ef6b6b76d9bca6b892d0f3df96fb9f9f53859c7261ccdb"
variables: []
secrets_allowed: false
```
````go
// Package operationaltelemetry is an AUTHORED, closed-outcome OTLP/HTTP JSON
// reporter for the refund host. It is not a general OpenTelemetry SDK, universal
// PII filter, business ledger, durable acknowledgement or automatic retry queue.
package operationaltelemetry

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"math"
	"mime"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

var ErrConfig = errors.New("operational telemetry configuration rejected: TELEMETRY_CONFIG_REJECTED")
var ErrUnavailable = errors.New("operational telemetry delivery unknown: TELEMETRY_DELIVERY_UNKNOWN")
var ErrObservation = errors.New("operational telemetry observation rejected: TELEMETRY_OBSERVATION_REJECTED")

type Config struct {
	Endpoint          string `json:"endpoint"`
	CAFile            string `json:"ca_file"`
	CASHA256          string `json:"ca_sha256"`
	CertificateFile   string `json:"certificate_file"`
	CertificateSHA256 string `json:"certificate_sha256"`
	PrivateKeyFile    string `json:"private_key_file"`
	PrivateKeySHA256  string `json:"private_key_sha256"`
	TimeoutMS         int    `json:"timeout_ms"`
}

var outcomes = [...]string{"started", "idle", "handled", "error", "stopped"}
var bounds = [...]float64{0.001, 0.01, 0.1, 1, 10, 60}

type series struct {
	count   int64
	sum     float64
	buckets [7]int64
}
type Reporter struct {
	client   *http.Client
	base     url.URL
	timeout  time.Duration
	gate     chan struct{}
	closed   atomic.Bool
	instance string
	started  time.Time
	sequence int64
	values   [5]series
}

// strictObject preserves duplicate-key detection instead of relying on the
// last-value-wins behavior of encoding/json. All input is already byte-bounded.
func strictObject(data []byte, allowed map[string]bool) (map[string]json.RawMessage, error) {
	d := json.NewDecoder(bytes.NewReader(data))
	token, err := d.Token()
	if err != nil || token != json.Delim('{') {
		return nil, ErrConfig
	}
	values := map[string]json.RawMessage{}
	for d.More() {
		token, err = d.Token()
		if err != nil {
			return nil, ErrConfig
		}
		key, ok := token.(string)
		if !ok || !allowed[key] {
			return nil, ErrConfig
		}
		if _, exists := values[key]; exists {
			return nil, ErrConfig
		}
		var raw json.RawMessage
		if d.Decode(&raw) != nil {
			return nil, ErrConfig
		}
		values[key] = raw
	}
	token, err = d.Token()
	if err != nil || token != json.Delim('}') {
		return nil, ErrConfig
	}
	if _, err = d.Token(); err != io.EOF {
		return nil, ErrConfig
	}
	return values, nil
}
func fixedFile(path, digest string, limit int64) ([]byte, error) {
	expected, err := hex.DecodeString(digest)
	if err != nil || len(expected) != 32 || len(digest) != 64 {
		return nil, ErrConfig
	}
	abs, err := filepath.Abs(path)
	if err != nil || path == "" {
		return nil, ErrConfig
	}
	for p := abs; ; p = filepath.Dir(p) {
		st, e := os.Lstat(p)
		if e != nil || st.Mode()&os.ModeSymlink != 0 {
			return nil, ErrConfig
		}
		parent := filepath.Dir(p)
		if parent == p {
			break
		}
	}
	f, err := os.Open(abs)
	if err != nil {
		return nil, ErrConfig
	}
	defer f.Close()
	st, err := f.Stat()
	if err != nil || !st.Mode().IsRegular() || st.Size() > limit {
		return nil, ErrConfig
	}
	data, err := io.ReadAll(io.LimitReader(f, limit+1))
	if err != nil || int64(len(data)) > limit {
		return nil, ErrConfig
	}
	got := sha256.Sum256(data)
	if !bytes.Equal(got[:], expected) {
		return nil, ErrConfig
	}
	return data, nil
}

// FromFile binds the config and every credential file to caller-approved hashes.
// Private key bytes are parsed in memory and never added to diagnostics/telemetry.
func FromFile(path, digest string) (*Reporter, error) {
	data, err := fixedFile(path, digest, 8192)
	if err != nil {
		return nil, ErrConfig
	}
	allowed := map[string]bool{"endpoint": true, "ca_file": true, "ca_sha256": true, "certificate_file": true, "certificate_sha256": true, "private_key_file": true, "private_key_sha256": true, "timeout_ms": true}
	values, err := strictObject(data, allowed)
	if err != nil || len(values) != len(allowed) {
		return nil, ErrConfig
	}
	var cfg Config
	if json.Unmarshal(data, &cfg) != nil {
		return nil, ErrConfig
	}
	ca, e1 := fixedFile(cfg.CAFile, cfg.CASHA256, 1<<20)
	cert, e2 := fixedFile(cfg.CertificateFile, cfg.CertificateSHA256, 1<<20)
	key, e3 := fixedFile(cfg.PrivateKeyFile, cfg.PrivateKeySHA256, 1<<20)
	if e1 != nil || e2 != nil || e3 != nil {
		return nil, ErrConfig
	}
	if cfg.TimeoutMS < 10 || cfg.TimeoutMS > 3000 {
		return nil, ErrConfig
	}
	return New(cfg.Endpoint, time.Duration(cfg.TimeoutMS)*time.Millisecond, ca, cert, key)
}

func New(endpoint string, timeout time.Duration, caPEM, certPEM, keyPEM []byte) (*Reporter, error) {
	u, err := url.Parse(endpoint)
	if err != nil || len(endpoint) > 2048 || u.Scheme != "https" || u.Hostname() == "" || u.User != nil || u.RawQuery != "" || u.ForceQuery || u.Fragment != "" || u.Opaque != "" || u.RawPath != "" || (u.Path != "" && u.Path != "/") || timeout < 10*time.Millisecond || timeout > 3*time.Second {
		return nil, ErrConfig
	}
	ca := x509.NewCertPool()
	if !ca.AppendCertsFromPEM(caPEM) {
		return nil, ErrConfig
	}
	certificate, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return nil, ErrConfig
	}
	instance := make([]byte, 16)
	if _, err = rand.Read(instance); err != nil {
		return nil, ErrConfig
	}
	transport := &http.Transport{Proxy: nil, DialContext: (&net.Dialer{Timeout: timeout}).DialContext, TLSClientConfig: &tls.Config{MinVersion: tls.VersionTLS13, RootCAs: ca, Certificates: []tls.Certificate{certificate}}, TLSHandshakeTimeout: timeout, ResponseHeaderTimeout: timeout, MaxResponseHeaderBytes: 8192, DisableKeepAlives: true, ForceAttemptHTTP2: false}
	return &Reporter{client: &http.Client{Transport: transport, Timeout: timeout, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}, base: *u, timeout: timeout, gate: make(chan struct{}, 1), instance: hex.EncodeToString(instance), started: time.Now()}, nil
}
func (r *Reporter) Close() {
	if r != nil {
		r.closed.Store(true)
		r.client.CloseIdleConnections()
	}
}
func textAttr(key, value string) any {
	return map[string]any{"key": key, "value": map[string]any{"stringValue": value}}
}
func numberAttr(key string, value int64) any {
	return map[string]any{"key": key, "value": map[string]any{"intValue": strconv.FormatInt(value, 10)}}
}
func stamp(t time.Time) string { return strconv.FormatInt(t.UnixNano(), 10) }
func (r *Reporter) resource() any {
	return map[string]any{"attributes": []any{textAttr("service.name", "elite-return-refund-worker"), textAttr("deployment.environment.name", "reference"), textAttr("service.instance.id", r.instance)}}
}
func scope() any { return map[string]any{"name": "elite.refund.host", "version": "1.0.0"} }

// Observe serializes calls within one shared deadline (including waiting for
// another observer). Input is a closed outcome and monotonic elapsed duration;
// it accepts no message, error, request, tenant, business ID or arbitrary fields.
// A failed or partial three-signal export returns unknown and is never retried.
func (r *Reporter) Observe(parent context.Context, outcome string, elapsed time.Duration) error {
	idx := -1
	for i, value := range outcomes {
		if value == outcome {
			idx = i
		}
	}
	if r == nil || parent == nil || idx < 0 || elapsed < 0 {
		return ErrObservation
	}
	ctx, cancel := context.WithTimeout(parent, r.timeout)
	defer cancel()
	select {
	case r.gate <- struct{}{}:
		defer func() { <-r.gate }()
	case <-ctx.Done():
		return ErrUnavailable
	}
	if r.closed.Load() || ctx.Err() != nil {
		return ErrUnavailable
	}
	if r.sequence == math.MaxInt64 || r.values[idx].count == math.MaxInt64 {
		return ErrObservation
	}
	seconds := elapsed.Seconds()
	sum := r.values[idx].sum + seconds
	if math.IsInf(sum, 0) || math.IsNaN(sum) {
		return ErrObservation
	}
	ids := make([]byte, 24)
	if _, err := rand.Read(ids); err != nil {
		return ErrUnavailable
	}
	trace, span := hex.EncodeToString(ids[:16]), hex.EncodeToString(ids[16:])
	end := time.Now()
	start := end.Add(-elapsed)
	r.sequence++
	r.values[idx].count++
	r.values[idx].sum = sum
	bucket := len(bounds)
	for i, b := range bounds {
		if seconds <= b {
			bucket = i
			break
		}
	}
	r.values[idx].buckets[bucket]++
	attrs := []any{textAttr("outcome", outcome), numberAttr("sequence", r.sequence)}
	severity := 9
	status := map[string]any{"code": 1}
	if outcome == "error" {
		severity = 17
		status = map[string]any{"code": 2}
	}
	log := map[string]any{"resourceLogs": []any{map[string]any{"resource": r.resource(), "scopeLogs": []any{map[string]any{"scope": scope(), "logRecords": []any{map[string]any{"timeUnixNano": stamp(end), "observedTimeUnixNano": stamp(end), "severityNumber": severity, "body": map[string]any{"stringValue": "refund_host_observation"}, "attributes": attrs, "traceId": trace, "spanId": span}}}}}}}
	spans := map[string]any{"resourceSpans": []any{map[string]any{"resource": r.resource(), "scopeSpans": []any{map[string]any{"scope": scope(), "spans": []any{map[string]any{"traceId": trace, "spanId": span, "name": "refund.host.observation", "kind": 1, "startTimeUnixNano": stamp(start), "endTimeUnixNano": stamp(end), "attributes": attrs, "status": status}}}}}}}
	counts := []any{}
	hist := []any{}
	for i, value := range r.values {
		a := []any{textAttr("outcome", outcomes[i])}
		counts = append(counts, map[string]any{"attributes": a, "startTimeUnixNano": stamp(r.started), "timeUnixNano": stamp(end), "asInt": strconv.FormatInt(value.count, 10)})
		buckets := []string{}
		for _, n := range value.buckets {
			buckets = append(buckets, strconv.FormatInt(n, 10))
		}
		hist = append(hist, map[string]any{"attributes": a, "startTimeUnixNano": stamp(r.started), "timeUnixNano": stamp(end), "count": strconv.FormatInt(value.count, 10), "sum": value.sum, "explicitBounds": bounds[:], "bucketCounts": buckets})
	}
	metrics := map[string]any{"resourceMetrics": []any{map[string]any{"resource": r.resource(), "scopeMetrics": []any{map[string]any{"scope": scope(), "metrics": []any{map[string]any{"name": "elite_refund_observations", "sum": map[string]any{"dataPoints": counts, "aggregationTemporality": 2, "isMonotonic": true}}, map[string]any{"name": "elite_refund_observation_duration", "unit": "s", "histogram": map[string]any{"dataPoints": hist, "aggregationTemporality": 2}}}}}}}}
	for _, signal := range []struct {
		path string
		data any
	}{{"logs", log}, {"metrics", metrics}, {"traces", spans}} {
		if r.send(ctx, signal.path, signal.data) != nil {
			return ErrUnavailable
		}
	}
	return nil
}
func (r *Reporter) send(ctx context.Context, signal string, value any) error {
	data, err := json.Marshal(value)
	if err != nil || len(data) > 16384 {
		return ErrUnavailable
	}
	u := r.base
	u.Path = "/v1/" + signal
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewReader(data))
	if err != nil {
		return ErrUnavailable
	}
	req.Header.Set("Content-Type", "application/json")
	response, err := r.client.Do(req)
	if err != nil {
		return ErrUnavailable
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return ErrUnavailable
	}
	typ, _, err := mime.ParseMediaType(response.Header.Get("Content-Type"))
	if err != nil || typ != "application/json" {
		return ErrUnavailable
	}
	b, err := io.ReadAll(io.LimitReader(response.Body, 1025))
	if err != nil || len(b) > 1024 {
		return ErrUnavailable
	}
	fields, err := strictObject(b, map[string]bool{"partialSuccess": true})
	if err != nil {
		return ErrUnavailable
	}
	if p, ok := fields["partialSuccess"]; ok {
		name := map[string]string{"logs": "rejectedLogRecords", "metrics": "rejectedDataPoints", "traces": "rejectedSpans"}[signal]
		partial, e := strictObject(p, map[string]bool{name: true, "errorMessage": true})
		if e != nil {
			return ErrUnavailable
		}
		for key, raw := range partial {
			if key == "errorMessage" {
				var message string
				if json.Unmarshal(raw, &message) != nil || message != "" {
					return ErrUnavailable
				}
			} else if strings.TrimSpace(string(raw)) != "0" && strings.TrimSpace(string(raw)) != "\"0\"" {
				return ErrUnavailable
			}
		}
	}
	return nil
}
````

### FILE: `return_refund_worker/internal/operationaltelemetry/reporter_test.go`
```yaml
block_id: "REFUND-TELEMETRY:reporter-test:v1"
operation: CREATE
provenance: AUTHORED
source: "local implementation; official OTLP HTTP JSON and Microsoft native lifecycle contracts; V372 exact integration evidence"
license: "LicenseRef-Workspace-Owner"
sha256: "f7123f4a30b263534c09db1aad02ebe6e6fac051805d31b2984b0dedde51369c"
variables: []
secrets_allowed: false
```
````go
package operationaltelemetry

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"io"
	"math"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type certs struct {
	ca, client, key []byte
	server          tls.Certificate
	pool            *x509.CertPool
}

func certificates(t *testing.T) certs {
	t.Helper()
	now := time.Now()
	pub, priv, e := ed25519.GenerateKey(rand.Reader)
	if e != nil {
		t.Fatal(e)
	}
	root := &x509.Certificate{SerialNumber: big.NewInt(1), Subject: pkix.Name{CommonName: "reference-test-ca"}, NotBefore: now.Add(-time.Minute), NotAfter: now.Add(time.Hour), IsCA: true, BasicConstraintsValid: true, KeyUsage: x509.KeyUsageCertSign}
	raw, e := x509.CreateCertificate(rand.Reader, root, root, pub, priv)
	if e != nil {
		t.Fatal(e)
	}
	ca := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: raw})
	makeLeaf := func(serial int64, server bool) ([]byte, []byte) {
		p, k, e := ed25519.GenerateKey(rand.Reader)
		if e != nil {
			t.Fatal(e)
		}
		usage := x509.ExtKeyUsageClientAuth
		if server {
			usage = x509.ExtKeyUsageServerAuth
		}
		c := &x509.Certificate{SerialNumber: big.NewInt(serial), Subject: pkix.Name{CommonName: "reference-test-peer"}, NotBefore: now.Add(-time.Minute), NotAfter: now.Add(time.Hour), KeyUsage: x509.KeyUsageDigitalSignature, ExtKeyUsage: []x509.ExtKeyUsage{usage}, IPAddresses: []net.IP{net.ParseIP("127.0.0.1")}}
		b, e := x509.CreateCertificate(rand.Reader, c, root, p, priv)
		if e != nil {
			t.Fatal(e)
		}
		pk, e := x509.MarshalPKCS8PrivateKey(k)
		if e != nil {
			t.Fatal(e)
		}
		return pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: b}), pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: pk})
	}
	client, key := makeLeaf(2, false)
	sc, sk := makeLeaf(3, true)
	pair, e := tls.X509KeyPair(sc, sk)
	if e != nil {
		t.Fatal(e)
	}
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(ca)
	return certs{ca, client, key, pair, pool}
}
func serverFor(t *testing.T, c certs, h http.HandlerFunc) *httptest.Server {
	t.Helper()
	s := httptest.NewUnstartedServer(h)
	s.TLS = &tls.Config{MinVersion: tls.VersionTLS13, Certificates: []tls.Certificate{c.server}, ClientAuth: tls.RequireAndVerifyClientCert, ClientCAs: c.pool}
	s.StartTLS()
	t.Cleanup(s.Close)
	return s
}
func newReporter(t *testing.T, c certs, s *httptest.Server) *Reporter {
	t.Helper()
	r, e := New(s.URL, 3*time.Second, c.ca, c.client, c.key)
	if e != nil {
		t.Fatal(e)
	}
	t.Cleanup(r.Close)
	return r
}
func accepted(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = io.WriteString(w, "{}")
}
func decodeMap(t *testing.T, b []byte) map[string]any {
	t.Helper()
	var j map[string]any
	if json.Unmarshal(b, &j) != nil {
		t.Fatalf("invalid JSON: %q", b)
	}
	return j
}
func child(v any, key string) any {
	if m, ok := v.(map[string]any); ok {
		return m[key]
	}
	return nil
}
func first(v any) any {
	if a, ok := v.([]any); ok && len(a) > 0 {
		return a[0]
	}
	return nil
}
func TestThreeSignalsPrivacyCorrelationAndMonotonicity(t *testing.T) {
	c := certificates(t)
	var mu sync.Mutex
	var paths []string
	var payloads [][]byte
	s := serverFor(t, c, func(w http.ResponseWriter, q *http.Request) {
		if q.Method != "POST" || q.Header.Get("Content-Type") != "application/json" || len(q.TLS.VerifiedChains) != 1 {
			t.Error("transport contract")
		}
		b, e := io.ReadAll(io.LimitReader(q.Body, 16385))
		if e != nil || len(b) > 16384 {
			t.Error("body budget")
		}
		mu.Lock()
		paths = append(paths, q.URL.Path)
		payloads = append(payloads, b)
		mu.Unlock()
		accepted(w)
	})
	r := newReporter(t, c, s)
	for _, o := range []string{"started", "idle", "handled", "error", "stopped"} {
		if e := r.Observe(context.Background(), o, 10*time.Millisecond); e != nil {
			t.Fatal(e)
		}
	}
	if len(paths) != 15 {
		t.Fatal(paths)
	}
	for n := 0; n < 5; n++ {
		if strings.Join(paths[n*3:n*3+3], ",") != "/v1/logs,/v1/metrics,/v1/traces" {
			t.Fatal(paths)
		}
		l := first(child(first(child(first(child(decodeMap(t, payloads[n*3]), "resourceLogs")), "scopeLogs")), "logRecords"))
		tr := first(child(first(child(first(child(decodeMap(t, payloads[n*3+2]), "resourceSpans")), "scopeSpans")), "spans"))
		if child(l, "traceId") != child(tr, "traceId") || child(l, "spanId") != child(tr, "spanId") || len(child(l, "traceId").(string)) != 32 {
			t.Fatal("correlation mismatch")
		}
		if child(child(l, "body"), "stringValue") != "refund_host_observation" {
			t.Fatal("unapproved message")
		}
	}
	if r.sequence != 5 {
		t.Fatal(r.sequence)
	}
	for _, v := range r.values {
		if v.count != 1 || v.buckets[1] != 1 {
			t.Fatal(v)
		}
	}
	before := len(paths)
	for _, o := range []string{"private@example.invalid", "password=private", "STARTED", ""} {
		if e := r.Observe(context.Background(), o, 0); !errors.Is(e, ErrObservation) {
			t.Fatal(e)
		}
	}
	if e := r.Observe(context.Background(), "idle", -1); !errors.Is(e, ErrObservation) {
		t.Fatal(e)
	}
	if len(paths) != before {
		t.Fatal("rejected input touched sink")
	}
	for _, b := range payloads {
		if strings.Contains(string(b), "private") || strings.Contains(string(b), "password") {
			t.Fatal("private marker")
		}
	}
}
func TestRemoteRejectionNeverRetriesOrLeaksBody(t *testing.T) {
	for _, scenario := range []string{"redirect", "failure", "rate", "partial", "warning", "oversized", "invalid", "duplicate", "wrong-type"} {
		t.Run(scenario, func(t *testing.T) {
			c := certificates(t)
			var calls atomic.Int64
			s := serverFor(t, c, func(w http.ResponseWriter, q *http.Request) {
				calls.Add(1)
				w.Header().Set("Content-Type", "application/json")
				switch scenario {
				case "redirect":
					w.Header().Set("Location", "https://invalid.example.invalid/private")
					w.WriteHeader(302)
				case "failure":
					w.WriteHeader(500)
				case "rate":
					w.WriteHeader(429)
				case "partial":
					_, _ = io.WriteString(w, `{"partialSuccess":{"rejectedLogRecords":"1"}}`)
				case "warning":
					_, _ = io.WriteString(w, `{"partialSuccess":{"errorMessage":"private@example.invalid"}}`)
				case "oversized":
					_, _ = io.WriteString(w, strings.Repeat("x", 1025))
				case "invalid":
					_, _ = io.WriteString(w, "private@example.invalid")
				case "duplicate":
					_, _ = io.WriteString(w, `{"partialSuccess":{},"partialSuccess":{}}`)
				case "wrong-type":
					w.Header().Set("Content-Type", "text/plain")
					_, _ = io.WriteString(w, "{}")
				}
			})
			r := newReporter(t, c, s)
			e := r.Observe(context.Background(), "error", 0)
			if !errors.Is(e, ErrUnavailable) || strings.Contains(e.Error(), "private") {
				t.Fatal(e)
			}
			if calls.Load() != 1 {
				t.Fatal("automatic retry/extra signal", calls.Load())
			}
		})
	}
}
func TestMidExportFailureIsUnknownWithoutReplay(t *testing.T) {
	c := certificates(t)
	var n atomic.Int64
	s := serverFor(t, c, func(w http.ResponseWriter, q *http.Request) {
		if n.Add(1) == 2 {
			w.WriteHeader(503)
			return
		}
		accepted(w)
	})
	r := newReporter(t, c, s)
	if !errors.Is(r.Observe(context.Background(), "handled", 0), ErrUnavailable) || n.Load() != 2 || r.values[2].count != 1 {
		t.Fatal("partial export not retained as uncertain")
	}
}
func TestCancellationBoundsSinkAndQueuedObserver(t *testing.T) {
	c := certificates(t)
	entered := make(chan struct{})
	release := make(chan struct{})
	var once sync.Once
	s := serverFor(t, c, func(w http.ResponseWriter, q *http.Request) {
		once.Do(func() { close(entered) })
		select {
		case <-release:
		case <-q.Context().Done():
		}
	})
	defer close(release)
	r := newReporter(t, c, s)
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- r.Observe(ctx, "idle", 0) }()
	select {
	case <-entered:
	case <-time.After(time.Second):
		cancel()
		t.Fatal("handler did not receive request")
	}
	qctx, qcancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer qcancel()
	start := time.Now()
	if !errors.Is(r.Observe(qctx, "idle", 0), ErrUnavailable) || time.Since(start) > time.Second {
		t.Fatal("waiter deadline")
	}
	cancel()
	select {
	case e := <-done:
		if !errors.Is(e, ErrUnavailable) {
			t.Fatal(e)
		}
	case <-time.After(time.Second):
		t.Fatal("active cancellation")
	}
}
func TestParallelObserversAndOverflow(t *testing.T) {
	c := certificates(t)
	var n atomic.Int64
	s := serverFor(t, c, func(w http.ResponseWriter, q *http.Request) { n.Add(1); accepted(w) })
	r := newReporter(t, c, s)
	var wg sync.WaitGroup
	for i := 0; i < 24; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if e := r.Observe(context.Background(), "handled", time.Millisecond); e != nil {
				t.Error(e)
			}
		}()
	}
	wg.Wait()
	if r.sequence != 24 || r.values[2].count != 24 || n.Load() != 72 {
		t.Fatal(r.sequence, n.Load())
	}
	r.sequence = math.MaxInt64
	if !errors.Is(r.Observe(context.Background(), "idle", 0), ErrObservation) || n.Load() != 72 {
		t.Fatal("sequence overflow")
	}
	r.sequence = 24
	r.values[1].count = math.MaxInt64
	if !errors.Is(r.Observe(context.Background(), "idle", 0), ErrObservation) || n.Load() != 72 {
		t.Fatal("counter overflow")
	}
	r.Close()
	if !errors.Is(r.Observe(context.Background(), "idle", 0), ErrUnavailable) {
		t.Fatal("closed observer")
	}
}
func TestTLSIdentityAndEndpointValidation(t *testing.T) {
	c := certificates(t)
	other := certificates(t)
	s := serverFor(t, c, func(w http.ResponseWriter, q *http.Request) { accepted(w) })
	for _, b := range []struct{ ca, cert, key []byte }{{other.ca, c.client, c.key}, {c.ca, other.client, other.key}} {
		r, e := New(s.URL, time.Second, b.ca, b.cert, b.key)
		if e != nil {
			t.Fatal(e)
		}
		if !errors.Is(r.Observe(context.Background(), "started", 0), ErrUnavailable) {
			t.Fatal("foreign identity accepted")
		}
		r.Close()
	}
	for _, u := range []string{"http://localhost", "https://user:password@localhost", "https://localhost/path", "https://localhost?x=y", "https://localhost#fragment", "https://"} {
		if _, e := New(u, time.Second, c.ca, c.client, c.key); !errors.Is(e, ErrConfig) {
			t.Fatal(u, e)
		}
	}
	for _, d := range []time.Duration{0, 9 * time.Millisecond, 4 * time.Second, time.Duration(math.MaxInt64)} {
		if _, e := New(s.URL, d, c.ca, c.client, c.key); !errors.Is(e, ErrConfig) {
			t.Fatal(d, e)
		}
	}
}
func TestExactConfigCredentialHashesAndNoDiagnostics(t *testing.T) {
	c := certificates(t)
	dir := t.TempDir()
	put := func(name string, b []byte) (string, string) {
		p := filepath.Join(dir, name)
		if e := os.WriteFile(p, b, 0600); e != nil {
			t.Fatal(e)
		}
		sum := sha256.Sum256(b)
		return p, hex.EncodeToString(sum[:])
	}
	caf, cah := put("ca.pem", c.ca)
	cf, ch := put("client.pem", c.client)
	kf, kh := put("client.key", c.key)
	cfg := Config{Endpoint: "https://127.0.0.1:1", CAFile: caf, CASHA256: cah, CertificateFile: cf, CertificateSHA256: ch, PrivateKeyFile: kf, PrivateKeySHA256: kh, TimeoutMS: 1000}
	b, _ := json.Marshal(cfg)
	p, h := put("config.json", b)
	r, e := FromFile(p, h)
	if e != nil {
		t.Fatal(e)
	}
	r.Close()
	for i, raw := range [][]byte{append(b[:len(b)-1], []byte(`,"endpoint":"https://other.invalid"}`)...), []byte(`{"endpoint":null}`), []byte(`{"Endpoint":"https://other.invalid"}`), append(b, []byte(` {}`)...)} {
		p, h = put("bad"+string(rune('a'+i)), raw)
		if _, e = FromFile(p, h); !errors.Is(e, ErrConfig) {
			t.Fatal(e)
		}
	}
	p, h = put("config2.json", b)
	if e = os.WriteFile(kf, []byte("private@example.invalid"), 0600); e != nil {
		t.Fatal(e)
	}
	_, e = FromFile(p, h)
	if !errors.Is(e, ErrConfig) || strings.Contains(e.Error(), "private") {
		t.Fatal(e)
	}
	if _, e = FromFile(dir, h); !errors.Is(e, ErrConfig) {
		t.Fatal("directory accepted")
	}
}

// Protocol seeds include exact duplicate keys, trailing values and wrong types.
func FuzzStrictOTLPResponseObject(f *testing.F) {
	for _, seed := range []string{`{}`, `{"partialSuccess":{}}`, `{"partialSuccess":{"rejectedLogRecords":"1"}}`, `{"partialSuccess":{},"partialSuccess":{}}`, `[]`, `{"unexpected":1}`, `{} {}`, `{"partialSuccess":null}`} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, raw string) {
		if len(raw) > 2048 {
			return
		}
		values, err := strictObject([]byte(raw), map[string]bool{"partialSuccess": true})
		if err != nil {
			if err.Error() != ErrConfig.Error() {
				t.Fatal("unbounded diagnostic")
			}
			return
		}
		if !json.Valid([]byte(raw)) {
			t.Fatal("invalid JSON accepted")
		}
		if len(values) > 1 {
			t.Fatal("unknown or duplicate field accepted")
		}
		for k := range values {
			if k != "partialSuccess" {
				t.Fatal("unknown key")
			}
		}
		if strings.HasPrefix(strings.TrimSpace(raw), "[") {
			t.Fatal("array accepted")
		}
	})
}
func TestConfigTimeoutCannotOverflowIntoAcceptedDuration(t *testing.T) {
	c := certificates(t)
	dir := t.TempDir()
	put := func(n string, b []byte) (string, string) {
		p := filepath.Join(dir, n)
		if err := os.WriteFile(p, b, 0600); err != nil {
			t.Fatal(err)
		}
		h := sha256.Sum256(b)
		return p, hex.EncodeToString(h[:])
	}
	caf, cah := put("ca", c.ca)
	cf, ch := put("cert", c.client)
	kf, kh := put("key", c.key)
	for _, n := range []int{0, 9, 3001, math.MaxInt64, 9223372036855} {
		cfg := Config{Endpoint: "https://127.0.0.1:1", CAFile: caf, CASHA256: cah, CertificateFile: cf, CertificateSHA256: ch, PrivateKeyFile: kf, PrivateKeySHA256: kh, TimeoutMS: n}
		b, _ := json.Marshal(cfg)
		p, h := put("cfg", b)
		if _, err := FromFile(p, h); !errors.Is(err, ErrConfig) {
			t.Fatal("timeout accepted", n)
		}
	}
}
````

### FILE: `return_refund_worker/cmd/return-refund-worker/native_stop_other.go`
```yaml
block_id: "REFUND-TELEMETRY:native-stop-other:v1"
operation: CREATE
provenance: AUTHORED
source: "local implementation; official OTLP HTTP JSON and Microsoft native lifecycle contracts; V372 exact integration evidence"
license: "LicenseRef-Workspace-Owner"
sha256: "9bade57a4db708cdc464970eec0e7d129cba9cffdf48b69db89351683bbf9efa"
variables: []
secrets_allowed: false
```
````go
//go:build !windows

package main

import (
	"context"
	"errors"
)

// Native handles are a Windows-only trusted-launcher protocol.
func nativeStopContext(parent context.Context, raw string) (context.Context, func() error, error) {
	if raw != "" {
		return parent, func() error { return nil }, errors.New("NATIVE_STOP_PROTOCOL_FAILED")
	}
	ctx, cancel := context.WithCancel(parent)
	return ctx, func() error { cancel(); return nil }, nil
}

func nativeStopFailure(context.Context) error { return nil }
````

### FILE: `return_refund_worker/cmd/return-refund-worker/native_stop_windows.go`
```yaml
block_id: "REFUND-TELEMETRY:native-stop-windows:v1"
operation: CREATE
provenance: AUTHORED
source: "local implementation; official OTLP HTTP JSON and Microsoft native lifecycle contracts; V372 exact integration evidence"
license: "LicenseRef-Workspace-Owner"
sha256: "97d705c2674f7d364d80d22a1029ab27978ebb5dfa15f72fad94c5a0b9668a91"
variables: []
secrets_allowed: false
```
````go
//go:build windows

package main

// AUTHORED qualification only. The launcher supplies a trusted manual-reset
// event; this is neither object-type authentication nor hostile-worker isolation.
import (
	"context"
	"errors"
	"strconv"
	"sync"
	"syscall"
	"unsafe"
)

var errNativeStop = errors.New("NATIVE_STOP_REQUESTED")
var errNativeProtocol = errors.New("NATIVE_STOP_PROTOCOL_FAILED")
var stopKernel = syscall.NewLazyDLL("kernel32.dll")
var stopDuplicate = stopKernel.NewProc("DuplicateHandle")
var stopWait = stopKernel.NewProc("WaitForSingleObject")
var stopClose = stopKernel.NewProc("CloseHandle")

func parseStopHandle(raw string) (uintptr, error) {
	n, err := strconv.ParseUint(raw, 10, strconv.IntSize)
	if err != nil || n == 0 || n > uint64(^uintptr(0)>>1) || strconv.FormatUint(n, 10) != raw {
		return 0, errNativeProtocol
	}
	return uintptr(n), nil
}

// Duplicate before waiting. Cleanup joins the waiter BEFORE closing its handle:
// CloseHandle during a pending Windows wait has undefined behavior.
// Empty raw retains the ordinary parent-context lifecycle. Cancellation stops
// claims through the real host loop; it does not guarantee a domain commit.
func nativeStopContext(parent context.Context, raw string) (context.Context, func() error, error) {
	ctx, cancel := context.WithCancelCause(parent)
	if raw == "" {
		return ctx, func() error { cancel(context.Canceled); return nil }, nil
	}
	source, err := parseStopHandle(raw)
	if err != nil {
		cancel(err)
		return ctx, func() error { return nil }, err
	}
	var owned uintptr
	ok, _, _ := stopDuplicate.Call(^uintptr(0), source, ^uintptr(0), uintptr(unsafe.Pointer(&owned)), 0x100000, 0, 0)
	if ok == 0 {
		cancel(errNativeProtocol)
		return ctx, func() error { return nil }, errNativeProtocol
	}
	done, joined := make(chan struct{}), make(chan struct{})
	var once sync.Once
	var closeErr error
	cleanup := func() error {
		once.Do(func() {
			close(done)
			<-joined
			ok, _, _ := stopClose.Call(owned)
			if ok == 0 {
				closeErr = errNativeProtocol
			}
			cancel(context.Canceled)
		})
		return closeErr
	}
	// A pre-signaled event must not expose a live context to the first claim.
	// The trusted launcher supplies a manual-reset event, so this read is not
	// destructive. Other waitable object types are outside the protocol.
	initial, _, _ := stopWait.Call(owned, 0)
	if initial != 258 {
		close(joined)
		if initial == 0 {
			cancel(errNativeStop)
			return ctx, cleanup, nil
		}
		cancel(errNativeProtocol)
		_ = cleanup()
		return ctx, cleanup, errNativeProtocol
	}
	go func() {
		defer close(joined)
		for {
			select {
			case <-done:
				return
			case <-ctx.Done():
				return
			default:
			}
			status, _, _ := stopWait.Call(owned, 20)
			switch status {
			case 0:
				cancel(errNativeStop)
				return
			case 258:
			default:
				cancel(errNativeProtocol)
				return
			}
		}
	}()
	return ctx, cleanup, nil
}

func nativeStopFailure(ctx context.Context) error {
	if errors.Is(context.Cause(ctx), errNativeProtocol) {
		return errNativeProtocol
	}
	return nil
}
````

### FILE: `return_refund_worker/cmd/return-refund-worker/native_stop_windows_test.go`
```yaml
block_id: "REFUND-TELEMETRY:native-stop-windows-test:v1"
operation: CREATE
provenance: AUTHORED
source: "local implementation; official OTLP HTTP JSON and Microsoft native lifecycle contracts; V372 exact integration evidence"
license: "LicenseRef-Workspace-Owner"
sha256: "a03917b6468cd93d6ea27f5bb5fa519eff761454a7f43b67a36b83e9439778fd"
variables: []
secrets_allowed: false
```
````go
//go:build windows

package main

import (
	"context"
	"errors"
	"strconv"
	"sync"
	"testing"
	"time"
	"unsafe"
)

func stopTestEvent(t *testing.T) uintptr {
	t.Helper()
	h, _, _ := stopKernel.NewProc("CreateEventW").Call(0, 1, 0, 0)
	if h == 0 {
		t.Fatal("event creation")
	}
	t.Cleanup(func() {
		if ok, _, _ := stopClose.Call(h); ok == 0 {
			t.Error("event close")
		}
	})
	return h
}

func FuzzNativeStopHandle(f *testing.F) {
	for _, v := range []string{"", "0", "4", "0004", "-1", "+4", "18446744073709551615", "9223372036854775807", "private", "4\x00"} {
		f.Add(v)
	}
	f.Fuzz(func(t *testing.T, raw string) {
		h, e := parseStopHandle(raw)
		if e != nil {
			if e != errNativeProtocol || h != 0 {
				t.Fatal("non-static rejection")
			}
			return
		}
		if h == 0 || h > ^uintptr(0)>>1 || strconv.FormatUint(uint64(h), 10) != raw {
			t.Fatal("accepted noncanonical or pseudo handle")
		}
	})
}

func TestNativeStopPreSignaledPreventsFirstRealHostClaim(t *testing.T) {
	h := stopTestEvent(t)
	stopKernel.NewProc("SetEvent").Call(h)
	ctx, close, e := nativeStopContext(context.Background(), strconv.FormatUint(uint64(h), 10))
	if e != nil {
		t.Fatal(e)
	}
	defer close()
	if ctx.Err() == nil {
		t.Fatal("pre-signaled event returned a live context")
	}
	claims := 0
	if err := runRefundLoop(ctx, func(context.Context, string) error { claims++; return nil }, claimToken); err != nil {
		t.Fatal(err)
	}
	if claims != 0 {
		t.Fatal("claim after pre-signaled stop")
	}
}
func TestNativeStopParse(t *testing.T) {
	for _, raw := range []string{"0", "-1", "+1", "01", " 1", "1 ", "18446744073709551615", "9223372036854775808", "1.0", "private-secret"} {
		if _, e := parseStopHandle(raw); e != errNativeProtocol {
			t.Fatalf("invalid accepted: %q", raw)
		}
	}
	if h, e := parseStopHandle("4"); e != nil || h != 4 {
		t.Fatal(h, e)
	}
}
func TestNativeStopInvalidClosedHandle(t *testing.T) {
	h, _, _ := stopKernel.NewProc("CreateEventW").Call(0, 1, 0, 0)
	if h == 0 {
		t.Fatal("create")
	}
	stopClose.Call(h)
	ctx, close, e := nativeStopContext(context.Background(), strconv.FormatUint(uint64(h), 10))
	if e != errNativeProtocol || context.Cause(ctx) != errNativeProtocol {
		t.Fatal("expected fixed failure")
	}
	if close() != nil {
		t.Fatal("cleanup")
	}
}
func TestNativeStopSignal(t *testing.T) {
	h := stopTestEvent(t)
	ctx, close, e := nativeStopContext(context.Background(), strconv.FormatUint(uint64(h), 10))
	if e != nil {
		t.Fatal(e)
	}
	defer close()
	stopKernel.NewProc("SetEvent").Call(h)
	select {
	case <-ctx.Done():
	case <-time.After(time.Second):
		t.Fatal("event ignored")
	}
	if context.Cause(ctx) != errNativeStop {
		t.Fatal(context.Cause(ctx))
	}
}
func TestNativeStopAlreadySignaled(t *testing.T) {
	h := stopTestEvent(t)
	stopKernel.NewProc("SetEvent").Call(h)
	ctx, close, e := nativeStopContext(context.Background(), strconv.FormatUint(uint64(h), 10))
	if e != nil {
		t.Fatal(e)
	}
	defer close()
	select {
	case <-ctx.Done():
	case <-time.After(time.Second):
		t.Fatal("event ignored")
	}
}
func TestNativeStopParentCause(t *testing.T) {
	parent, cancel := context.WithCancelCause(context.Background())
	cause := errors.New("parent cause")
	h := stopTestEvent(t)
	ctx, close, e := nativeStopContext(parent, strconv.FormatUint(uint64(h), 10))
	if e != nil {
		t.Fatal(e)
	}
	cancel(cause)
	if close() != nil || context.Cause(ctx) != cause {
		t.Fatal("cause or cleanup")
	}
}
func TestNativeStopEmpty(t *testing.T) {
	ctx, close, e := nativeStopContext(context.Background(), "")
	if e != nil || ctx.Err() != nil {
		t.Fatal(e)
	}
	if close() != nil || ctx.Err() != context.Canceled {
		t.Fatal("empty cleanup")
	}
}
func TestNativeStopConcurrentClose(t *testing.T) {
	h := stopTestEvent(t)
	_, close, e := nativeStopContext(context.Background(), strconv.FormatUint(uint64(h), 10))
	if e != nil {
		t.Fatal(e)
	}
	var wg sync.WaitGroup
	for range 24 {
		wg.Go(func() {
			if close() != nil {
				t.Error("cleanup")
			}
		})
	}
	wg.Wait()
}
func TestNativeStopDuplicateSurvivesSourceClosure(t *testing.T) {
	h, _, _ := stopKernel.NewProc("CreateEventW").Call(0, 1, 1, 0)
	if h == 0 {
		t.Fatal("create")
	}
	ctx, close, e := nativeStopContext(context.Background(), strconv.FormatUint(uint64(h), 10))
	if e != nil {
		stopClose.Call(h)
		t.Fatal(e)
	}
	stopClose.Call(h)
	defer close()
	select {
	case <-ctx.Done():
	case <-time.After(time.Second):
		t.Fatal("duplicate not independent")
	}
	if context.Cause(ctx) != errNativeStop {
		t.Fatal(context.Cause(ctx))
	}
}
func TestNativeStopRepeatedNoHandleGrowth(t *testing.T) {
	count := func() uint32 {
		var n uint32
		ok, _, _ := stopKernel.NewProc("GetProcessHandleCount").Call(^uintptr(0), uintptr(unsafe.Pointer(&n)))
		if ok == 0 {
			t.Fatal("count")
		}
		return n
	}
	h := stopTestEvent(t)
	raw := strconv.FormatUint(uint64(h), 10)
	_, close, _ := nativeStopContext(context.Background(), raw)
	close()
	before := count()
	for range 64 {
		_, close, e := nativeStopContext(context.Background(), raw)
		if e != nil || close() != nil {
			t.Fatal("cycle")
		}
	}
	if count() != before {
		t.Fatal("native handle growth")
	}
}

func TestNativeProtocolFailureRemainsHostFailure(t *testing.T) {
	ctx, cancel := context.WithCancelCause(context.Background())
	cancel(errNativeProtocol)
	if !errors.Is(nativeStopFailure(ctx), errNativeProtocol) {
		t.Fatal("native protocol failure suppressed")
	}
	ctx, cancel = context.WithCancelCause(context.Background())
	cancel(errNativeStop)
	if nativeStopFailure(ctx) != nil {
		t.Fatal("cooperative stop misclassified")
	}
}
````

### FILE: `return_refund_worker/cmd/return-refund-worker/operational_telemetry_test.go`
```yaml
block_id: "REFUND-TELEMETRY:operational-telemetry-test:v1"
operation: CREATE
provenance: AUTHORED
source: "local implementation; official OTLP HTTP JSON and Microsoft native lifecycle contracts; V372 exact integration evidence"
license: "LicenseRef-Workspace-Owner"
sha256: "03881b8240d3160fcf67901300228f86c5e4685f6e423e562b527b96a27fa876"
variables: []
secrets_allowed: false
```
````go
package main

import (
	"context"
	"elite.local/return-refund-worker/internal/operationaltelemetry"
	"elite.local/return-refund-worker/internal/refundworker"
	"errors"
	"testing"
	"time"
)

func TestObservedLoopStopsOnUnknownDeliveryWithoutNextClaim(t *testing.T) {
	for _, tc := range []struct {
		name    string
		stepErr error
		want    string
	}{{"idle", refundworker.ErrNoWork, "idle"}, {"handled", nil, "handled"}, {"error", errors.New("private-customer-diagnostic"), "error"}} {
		t.Run(tc.name, func(t *testing.T) {
			calls, tokens, reports := 0, 0, 0
			err := runObservedRefundLoop(context.Background(), func(context.Context, string) error { calls++; return tc.stepErr }, func() (string, error) { tokens++; return "synthetic-token", nil }, func(ctx context.Context, outcome string, elapsed time.Duration) error {
				reports++
				if outcome != tc.want || elapsed < 0 {
					t.Fatal("wrong closed outcome")
				}
				return operationaltelemetry.ErrUnavailable
			})
			if !errors.Is(err, operationaltelemetry.ErrUnavailable) || calls != 1 || tokens != 1 || reports != 1 {
				t.Fatalf("unexpected counts: %v %d %d %d", err, calls, tokens, reports)
			}
		})
	}
}
func TestObservedLoopPrecancellationMakesNoClaimOrReport(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	fail := func() { t.Fatal("pre-cancelled work") }
	if err := runObservedRefundLoop(ctx, func(context.Context, string) error { fail(); return nil }, func() (string, error) { fail(); return "", nil }, func(context.Context, string, time.Duration) error { fail(); return nil }); err != nil {
		t.Fatal(err)
	}
}
````

### FILE: `return_refund_worker/cmd/telemetry-reference/certificates.go`
```yaml
block_id: "REFUND-TELEMETRY:certificates:v1"
operation: CREATE
provenance: AUTHORED
source: "local implementation; official OTLP HTTP JSON and Microsoft native lifecycle contracts; V372 exact integration evidence"
license: "LicenseRef-Workspace-Owner"
sha256: "ef6b4a375f6763f352feedd9a1a15b28cf16371c5f9c14c6609b150dadd74a12"
variables: []
secrets_allowed: false
```
````go
// AUTHORED ephemeral reference fixture; not a certificate authority service.
package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"
)

func referenceCertificates(directory string) error {
	info, e := os.Lstat(directory)
	if e != nil || !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return errors.New("REFERENCE_CERTIFICATE_REJECTED")
	}
	write := func(name string, data []byte) error {
		f, e := os.OpenFile(filepath.Join(directory, name), os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
		if e != nil {
			return e
		}
		_, e = f.Write(data)
		if e == nil {
			e = f.Sync()
		}
		ce := f.Close()
		if e != nil {
			return e
		}
		return ce
	}
	serial := func() (*big.Int, error) { return rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128)) }
	pub, key, e := ed25519.GenerateKey(rand.Reader)
	if e != nil {
		return e
	}
	id := sha256.Sum256(pub)
	sn, e := serial()
	if e != nil {
		return e
	}
	now := time.Now()
	ca := &x509.Certificate{SerialNumber: sn, Subject: pkix.Name{CommonName: "Elite ephemeral reference CA"}, NotBefore: now.Add(-time.Minute), NotAfter: now.Add(12 * time.Hour), IsCA: true, BasicConstraintsValid: true, MaxPathLen: 0, MaxPathLenZero: true, KeyUsage: x509.KeyUsageCertSign | x509.KeyUsageCRLSign, SubjectKeyId: id[:20]}
	der, e := x509.CreateCertificate(rand.Reader, ca, ca, pub, key)
	if e != nil {
		return e
	}
	if e = write("ca.pem", pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})); e != nil {
		return e
	}
	for _, role := range []string{"server", "client", "rogue"} {
		lp, lk, e := ed25519.GenerateKey(rand.Reader)
		if e != nil {
			return e
		}
		sn, e = serial()
		if e != nil {
			return e
		}
		lid := sha256.Sum256(lp)
		leaf := &x509.Certificate{SerialNumber: sn, Subject: pkix.Name{CommonName: "elite-reference-" + role}, NotBefore: now.Add(-time.Minute), NotAfter: now.Add(2 * time.Hour), BasicConstraintsValid: true, KeyUsage: x509.KeyUsageDigitalSignature, SubjectKeyId: lid[:20], AuthorityKeyId: ca.SubjectKeyId, ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}}
		if role == "server" {
			leaf.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
			leaf.IPAddresses = []net.IP{net.ParseIP("127.0.0.1")}
		}
		signer := key
		if role == "rogue" {
			signer = lk
		}
		// The rogue fixture is a valid self-issued leaf whose issuer name resembles
		// the real CA, but whose signing key is unrelated. It must never authenticate.
		parent := ca
		if role == "rogue" {
			copyCA := *ca
			copyCA.PublicKey = lp
			parent = &copyCA
		}
		der, e = x509.CreateCertificate(rand.Reader, leaf, parent, lp, signer)
		if e != nil {
			return e
		}
		if e = write(role+".pem", pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})); e != nil {
			return e
		}
		kb, e := x509.MarshalPKCS8PrivateKey(lk)
		if e != nil {
			return e
		}
		if e = write(role+".key", pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: kb})); e != nil {
			return e
		}
	}
	return nil
}
````

### FILE: `return_refund_worker/cmd/telemetry-reference/main.go`
```yaml
block_id: "REFUND-TELEMETRY:main:v1"
operation: CREATE
provenance: AUTHORED
source: "local implementation; official OTLP HTTP JSON and Microsoft native lifecycle contracts; V372 exact integration evidence"
license: "LicenseRef-Workspace-Owner"
sha256: "18b17eb8d65be1306904c02a528b8ed639e244932f9206ff00e4c07ae741ee55"
variables: []
secrets_allowed: false
```
````go
// AUTHORED finite reference load probe. No business data, provider or refund call.
package main

import (
	"context"
	"elite.local/return-refund-worker/internal/operationaltelemetry"
	"encoding/json"
	"os"
	"strconv"
	"sync"
	"time"
)

func run() error {
	if len(os.Args) == 3 && os.Args[1] == "--certificates" {
		return referenceCertificates(os.Args[2])
	}
	if len(os.Args) != 5 {
		return operationaltelemetry.ErrConfig
	}
	n, e := strconv.Atoi(os.Args[3])
	if e != nil || n < 1 || n > 2000 {
		return operationaltelemetry.ErrConfig
	}
	parallel, e := strconv.Atoi(os.Args[4])
	if e != nil || parallel < 1 || parallel > 16 || n%parallel != 0 {
		return operationaltelemetry.ErrConfig
	}
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	began := time.Now()
	var wg sync.WaitGroup
	errs := make(chan error, parallel)
	for range parallel {
		wg.Go(func() {
			r, e := operationaltelemetry.FromFile(os.Args[1], os.Args[2])
			if e != nil {
				errs <- e
				return
			}
			defer r.Close()
			for range n / parallel {
				if e = r.Observe(ctx, "idle", time.Millisecond); e != nil {
					errs <- e
					cancel()
					return
				}
			}
		})
	}
	wg.Wait()
	close(errs)
	for e := range errs {
		return e
	}
	return json.NewEncoder(os.Stdout).Encode(struct {
		Observations int     `json:"observations"`
		Instances    int     `json:"instances"`
		Seconds      float64 `json:"seconds"`
	}{n, parallel, time.Since(began).Seconds()})
}
func main() {
	if run() != nil {
		os.Stderr.WriteString("REFERENCE_TELEMETRY_PROBE_FAILED\n")
		os.Exit(1)
	}
}
````


### FILE: `return_refund_worker/docs/COMMERCE_PROVIDER_ALIAS.md`
```yaml
block_id: "GO-OFFICIAL-RETURN-REFUND:provider-alias-contract:v1"
operation: CREATE
provenance: AUTHORED
source: "local compatibility verification/documentation; pinned official SDK invoked, not represented as vendor-authored code"
license: "LicenseRef-Workspace-Owner"
sha256: "21d6e4003b6c911acdf690288114d1737e5af419296f15fa769792969dbe3e94"
variables: []
secrets_allowed: false
```
````markdown
# Commerce to refund provider compatibility

`mercadopago` is the Commerce payment code. `mercado_pago` remains the existing
refund owner's code, provider registration name and durable refund projection.
`PostgresStore.Prepare` maps the former to the latter once, while selecting a new
refund's source payment. Both exact spellings are accepted; no case folding or
other alias is introduced. Stripe is unchanged.

The source payment row is never renamed. Existing refund rows are returned before
this mapping and remain byte-for-byte unchanged in their identity fields. Payment
attempt IDs, source provider references, return request IDs and idempotency keys
are carried through without modification. The existing SQL constraint and SDK
adapter continue to receive `mercado_pago`; no historical migration is needed.

This is AUTHORED compatibility glue. The official Mercado Pago Go SDK remains
dependency-pinned at 1.14.0. No financial rules, rounding, refund policy, SDK source
or attribution changes in this delta.

The integration fixture seeds the same explicitly disposable database as the
legacy refund test, including its documented trigger-bypass setup/cleanup. It
then runs actual store queries/transactions and pinned SDK requests through an
in-process transport for both source spellings. It proves partial and full refund
states, persisted identity, GET reconciliation, exact idempotency headers,
ambiguity rejection and immutable observations. It does not prove originating
return authorization journeys or live provider acceptance.
````


### FILE: `return_refund_worker/internal/refundworker/provider_alias_integration_test.go`
```yaml
block_id: "GO-OFFICIAL-RETURN-REFUND:provider-alias-test:v1"
operation: CREATE
provenance: AUTHORED
source: "local compatibility verification/documentation; pinned official SDK invoked, not represented as vendor-authored code"
license: "LicenseRef-Workspace-Owner"
sha256: "4a3a9d2b9ef20ac08006926b3ed48a5d7009171f14dcc700f8e46cf906626fc9"
variables: []
secrets_allowed: false
```
````go
package refundworker

// AUTHORED compatibility fixtures: real pinned Mercado Pago SDK calls use an
// in-process requester. No network, account, secret or provider capture occurs.
import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestPostgresMercadoPagoProviderAliases(t *testing.T) {
	for _, sourceProvider := range []string{"mercado_pago", "mercadopago"} {
		t.Run(sourceProvider, func(t *testing.T) {
			creates := make(map[string]int)
			paymentReads, refundReads := 0, 0
			transport := requesterFunc(func(request *http.Request) (*http.Response, error) {
				if request.Header.Get("Authorization") != "Bearer TEST-local" || request.URL.Host != "api.mercadopago.com" {
					return nil, fmt.Errorf("unexpected SDK authority")
				}
				var response string
				switch {
				case request.Method == http.MethodGet && request.URL.Path == "/v1/payments/7186040733":
					paymentReads++
					response = `{"id":7186040733,"status":"approved","captured":true,"currency_id":"ARS","transaction_amount":20.00}`
				case request.Method == http.MethodPost && request.URL.Path == "/v1/payments/7186040733/refunds":
					key := request.Header.Get("X-Idempotency-Key")
					if key != "return-refund-key-0001" && key != "return-refund-key-0002" {
						return nil, fmt.Errorf("historical refund idempotency key changed")
					}
					var body map[string]any
					if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
						return nil, err
					}
					if body["amount"] != float64(10) {
						return nil, fmt.Errorf("refund amount changed: %v", body["amount"])
					}
					creates[key]++
					id := 1622029221
					if key == "return-refund-key-0002" {
						id++
					}
					response = fmt.Sprintf(`{"id":%d,"payment_id":7186040733,"status":"approved","amount":10.00}`, id)
				case request.Method == http.MethodGet && request.URL.Path == "/v1/payments/7186040733/refunds/1622029221":
					refundReads++
					response = `{"id":1622029221,"payment_id":7186040733,"status":"approved","amount":10.00}`
				default:
					return nil, fmt.Errorf("unexpected SDK request %s %s", request.Method, request.URL.Path)
				}
				return &http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(response)), Request: request}, nil
			})
			provider, err := newMercadoPagoProvider("TEST-local", transport)
			if err != nil {
				t.Fatal(err)
			}
			testRefundLineAllocationPartialThenFullAndAmbiguity(t, sourceProvider, provider)
			if len(creates) != 2 || creates["return-refund-key-0001"] != 1 || creates["return-refund-key-0002"] != 1 || paymentReads != 3 || refundReads != 1 {
				t.Fatalf("SDK effects/reconciliation changed: creates=%v payment_reads=%d refund_reads=%d", creates, paymentReads, refundReads)
			}
		})
	}
}
````
