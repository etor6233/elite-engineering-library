# PostgreSQL Transactional Foundation

## 1. Metadata

```yaml
pack_id: "PG-TX-FOUNDATION"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa un baseline PostgreSQL multi-tenant para configuración versionada, idempotencia, outbox, inbox, jobs y auditoría append-only."
stacks:
  - "PostgreSQL 18.6"
compatible_with:
  - "PBC-CORE >=0.1.0 <1.0.0"
incompatible_with:
  - "bases que no soporten generated identity, jsonb o índices parciales"
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources:
  - "https://www.postgresql.org/docs/18/"
  - "https://www.postgresql.org/support/versioning/"
verified_at: "2026-08-25"
```

Los bloques son `AUTHORED`; usan sintaxis y conceptos documentados públicamente, sin copiar código de PostgreSQL. La versión 18.6 se fija porque es el minor estable actual verificado en la fecha indicada. Una actualización exige repetir migración, seguridad y restore gates.

## 2. Applicability

Adoptar para un sistema transaccional nuevo que necesite:

- múltiples tenants con ownership explícito;
- configuración empresarial versionada;
- reintentos seguros mediante idempotencia;
- publicación confiable mediante transactional outbox;
- deduplicación de consumers;
- jobs con lease/retry/dead-letter;
- auditoría inmutable desde el rol de aplicación.

No adoptar sin cambios para analytics masivo, event sourcing puro, hard multi-tenancy regulatorio que exija bases separadas, ni escrituras multi-región activas.

Límites:

- RLS no se activa en este pack. Requiere diseño de roles, pool/session reset y pruebas de bypass específicas del runtime;
- append-only no detiene a superusuarios ni propietarios que modifiquen DDL;
- outbox ofrece entrega al menos una vez, nunca exactamente una vez;
- no contiene tablas particulares de catálogo, ventas, fábrica o garantías.

## 3. Architecture contract

Ownership:

```text
platform.tenant
  ├── platform.business_profile_version
  ├── platform.idempotency_record
  ├── platform.outbox_event
  ├── platform.consumer_inbox
  ├── platform.job
  └── audit.event
```

Invariantes:

- todas las filas operativas incluyen `tenant_id` y FK al tenant;
- IDs son generados por la aplicación; no hay extensión oculta;
- misma idempotency key con request hash distinto es conflicto que resuelve el adapter;
- evento de dominio único por tenant, aggregate, version y type;
- inbox y write local se confirman en la misma transacción del consumer;
- audit no admite update/delete para el rol de aplicación;
- los side effects externos ocurren después del commit y toleran duplicados.

Failure/degraded modes:

- publisher caído: outbox acumula; transacciones de negocio continúan hasta el límite operativo;
- consumer caído: redelivery y deduplicación por inbox;
- worker caído: lease vence y otro worker reclama;
- DB no disponible: fail closed para writes; no simular éxito;
- migration falla: `ON_ERROR_STOP` y transacción abortan sin estado parcial.

## 4. Exact file manifest

```text
CREATE db/migrations/0001_platform_foundation.up.sql
CREATE db/migrations/0001_platform_foundation.down.sql
CREATE db/tests/0001_platform_foundation.test.sql
```

## 5. Materialization blocks

### FILE: `db/migrations/0001_platform_foundation.up.sql`

```yaml
block_id: "PG-TX-FOUNDATION:migration-up:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "9818f77246aded2ec4a2f12580dc5b3ca833b29e4d06a3949cc313ac4db64045"
variables: []
secrets_allowed: false
```

````sql
begin;

create schema platform;
create schema audit;

create table platform.tenant (
  tenant_id uuid primary key,
  tenant_code text not null unique,
  legal_name text not null,
  display_name text not null,
  status text not null default 'active'
    check (status in ('provisioning', 'active', 'suspended', 'closed')),
  created_at timestamptz not null default clock_timestamp(),
  closed_at timestamptz,
  check ((status = 'closed') = (closed_at is not null)),
  check (tenant_code ~ '^[a-z][a-z0-9]*(-[a-z0-9]+)*$')
);

create table platform.business_profile_version (
  tenant_id uuid not null,
  profile_version bigint not null check (profile_version > 0),
  schema_version text not null,
  profile jsonb not null,
  sha256_hex text not null,
  status text not null
    check (status in ('draft', 'admitted', 'active', 'retired', 'rejected')),
  created_by_subject text not null,
  created_at timestamptz not null default clock_timestamp(),
  activated_at timestamptz,
  primary key (tenant_id, profile_version),
  foreign key (tenant_id) references platform.tenant (tenant_id),
  check (jsonb_typeof(profile) = 'object'),
  check (sha256_hex ~ '^[0-9a-f]{64}$'),
  check ((status = 'active') = (activated_at is not null))
);

create unique index business_profile_one_active_per_tenant_uq
  on platform.business_profile_version (tenant_id)
  where status = 'active';

create table platform.idempotency_record (
  tenant_id uuid not null,
  scope text not null,
  idempotency_key text not null,
  request_sha256_hex text not null,
  status text not null
    check (status in ('processing', 'completed', 'failed_terminal')),
  response_code integer,
  response_body jsonb,
  resource_type text,
  resource_id text,
  locked_until timestamptz,
  created_at timestamptz not null default clock_timestamp(),
  expires_at timestamptz not null,
  primary key (tenant_id, scope, idempotency_key),
  foreign key (tenant_id) references platform.tenant (tenant_id),
  check (length(scope) between 1 and 128),
  check (length(idempotency_key) between 8 and 200),
  check (request_sha256_hex ~ '^[0-9a-f]{64}$'),
  check (expires_at > created_at),
  check (
    (status = 'processing' and response_code is null)
    or
    (status in ('completed', 'failed_terminal') and response_code between 100 and 599)
  )
);

create index idempotency_expiry_idx
  on platform.idempotency_record (expires_at);

create table platform.outbox_event (
  tenant_id uuid not null,
  event_id uuid not null,
  aggregate_type text not null,
  aggregate_id text not null,
  aggregate_version bigint not null check (aggregate_version > 0),
  event_type text not null,
  schema_version integer not null check (schema_version > 0),
  occurred_at timestamptz not null,
  trace_id text,
  payload jsonb not null,
  headers jsonb not null default '{}'::jsonb,
  available_at timestamptz not null default clock_timestamp(),
  claimed_by text,
  claimed_until timestamptz,
  attempts integer not null default 0 check (attempts >= 0),
  published_at timestamptz,
  last_error_code text,
  primary key (tenant_id, event_id),
  foreign key (tenant_id) references platform.tenant (tenant_id),
  unique (tenant_id, aggregate_type, aggregate_id, aggregate_version, event_type),
  check (jsonb_typeof(payload) = 'object'),
  check (jsonb_typeof(headers) = 'object'),
  check ((claimed_by is null) = (claimed_until is null)),
  check (published_at is null or published_at >= occurred_at)
);

create index outbox_ready_idx
  on platform.outbox_event (available_at, occurred_at, event_id)
  where published_at is null;

create table platform.consumer_inbox (
  tenant_id uuid not null,
  consumer_name text not null,
  event_id uuid not null,
  received_at timestamptz not null default clock_timestamp(),
  completed_at timestamptz,
  result_sha256_hex text,
  primary key (tenant_id, consumer_name, event_id),
  foreign key (tenant_id) references platform.tenant (tenant_id),
  check (length(consumer_name) between 1 and 128),
  check (completed_at is null or completed_at >= received_at),
  check (result_sha256_hex is null or result_sha256_hex ~ '^[0-9a-f]{64}$')
);

create index inbox_incomplete_idx
  on platform.consumer_inbox (received_at)
  where completed_at is null;

create table platform.job (
  tenant_id uuid not null,
  job_id uuid not null,
  queue text not null,
  job_type text not null,
  schema_version integer not null check (schema_version > 0),
  payload jsonb not null,
  priority smallint not null default 0,
  available_at timestamptz not null default clock_timestamp(),
  attempts integer not null default 0 check (attempts >= 0),
  max_attempts integer not null check (max_attempts between 1 and 100),
  claimed_by text,
  claimed_until timestamptz,
  completed_at timestamptz,
  terminal_error_code text,
  created_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, job_id),
  foreign key (tenant_id) references platform.tenant (tenant_id),
  check (jsonb_typeof(payload) = 'object'),
  check ((claimed_by is null) = (claimed_until is null)),
  check (not (completed_at is not null and terminal_error_code is not null)),
  check (completed_at is null or completed_at >= created_at)
);

create index job_ready_idx
  on platform.job (queue, priority desc, available_at, job_id)
  where completed_at is null and terminal_error_code is null;

create table audit.event (
  audit_sequence bigint generated always as identity primary key,
  tenant_id uuid not null,
  event_id uuid not null,
  occurred_at timestamptz not null default clock_timestamp(),
  actor_subject text not null,
  action text not null,
  resource_type text not null,
  resource_id text not null,
  trace_id text,
  source_ip inet,
  decision text not null check (decision in ('allowed', 'denied', 'system')),
  reason_code text,
  evidence jsonb not null default '{}'::jsonb,
  foreign key (tenant_id) references platform.tenant (tenant_id),
  unique (tenant_id, event_id),
  check (jsonb_typeof(evidence) = 'object')
);

create index audit_tenant_time_idx
  on audit.event (tenant_id, occurred_at desc, audit_sequence desc);

create function audit.reject_mutation()
returns trigger
language plpgsql
as $function$
begin
  raise exception using
    errcode = '55000',
    message = 'audit.event is append-only';
end;
$function$;

create trigger audit_event_reject_update_or_delete
before update or delete on audit.event
for each row execute function audit.reject_mutation();

commit;
````

### FILE: `db/migrations/0001_platform_foundation.down.sql`

```yaml
block_id: "PG-TX-FOUNDATION:migration-down:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "ba5a09127aaf4967a39617a86bf7795610d46bda0cdf832905983a4c568db8c7"
variables: []
secrets_allowed: false
```

````sql
begin;

drop schema audit cascade;
drop schema platform cascade;

commit;
````

### FILE: `db/tests/0001_platform_foundation.test.sql`

```yaml
block_id: "PG-TX-FOUNDATION:migration-test:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "dac0420b59fa9d4ad41368da882abc6b9bd5ab3e75e728a1c870d460fda68613"
variables: []
secrets_allowed: false
```

````sql
begin;

insert into platform.tenant (
  tenant_id,
  tenant_code,
  legal_name,
  display_name
) values (
  '018f4d4a-7b36-7a21-8d10-2f4c54c28a01',
  'test-tenant',
  'Test Tenant S.A.',
  'Test Tenant'
);

insert into platform.business_profile_version (
  tenant_id,
  profile_version,
  schema_version,
  profile,
  sha256_hex,
  status,
  created_by_subject,
  activated_at
) values (
  '018f4d4a-7b36-7a21-8d10-2f4c54c28a01',
  1,
  '1.0.0',
  '{"schemaVersion":"1.0.0"}'::jsonb,
  repeat('a', 64),
  'active',
  'test:operator',
  clock_timestamp()
);

do $test$
begin
  begin
    insert into platform.business_profile_version (
      tenant_id, profile_version, schema_version, profile,
      sha256_hex, status, created_by_subject, activated_at
    ) values (
      '018f4d4a-7b36-7a21-8d10-2f4c54c28a01', 2, '1.0.0',
      '{"schemaVersion":"1.0.0"}'::jsonb, repeat('b', 64),
      'active', 'test:operator', clock_timestamp()
    );
    raise exception 'expected one-active-profile unique violation';
  exception
    when unique_violation then null;
  end;
end;
$test$;

insert into platform.idempotency_record (
  tenant_id,
  scope,
  idempotency_key,
  request_sha256_hex,
  status,
  locked_until,
  expires_at
) values (
  '018f4d4a-7b36-7a21-8d10-2f4c54c28a01',
  'order:create',
  'idem-test-0001',
  repeat('c', 64),
  'processing',
  clock_timestamp() + interval '30 seconds',
  clock_timestamp() + interval '24 hours'
);

insert into platform.outbox_event (
  tenant_id,
  event_id,
  aggregate_type,
  aggregate_id,
  aggregate_version,
  event_type,
  schema_version,
  occurred_at,
  payload
) values (
  '018f4d4a-7b36-7a21-8d10-2f4c54c28a01',
  '018f4d4a-7b36-7a21-8d10-2f4c54c28a10',
  'order',
  'order-test-1',
  1,
  'order.created',
  1,
  clock_timestamp(),
  '{"orderId":"order-test-1"}'::jsonb
);

insert into audit.event (
  tenant_id,
  event_id,
  actor_subject,
  action,
  resource_type,
  resource_id,
  decision
) values (
  '018f4d4a-7b36-7a21-8d10-2f4c54c28a01',
  '018f4d4a-7b36-7a21-8d10-2f4c54c28a20',
  'test:operator',
  'business-profile.activate',
  'business-profile',
  '1',
  'allowed'
);

do $test$
begin
  begin
    update audit.event
       set action = 'tampered'
     where event_id = '018f4d4a-7b36-7a21-8d10-2f4c54c28a20';
    raise exception 'expected append-only mutation rejection';
  exception
    when object_not_in_prerequisite_state then null;
  end;
end;
$test$;

do $test$
declare
  actual_count bigint;
begin
  select count(*) into actual_count
    from platform.outbox_event
   where tenant_id = '018f4d4a-7b36-7a21-8d10-2f4c54c28a01';

  if actual_count <> 1 then
    raise exception 'expected exactly one outbox event, got %', actual_count;
  end if;
end;
$test$;

rollback;
````

## 6. Configuration surface

| Variable | Tipo | Default | Secreto | Validación | Efecto |
|---|---|---|---|---|---|
| `DATABASE_URL` | URI | ninguno | sí | TLS/host/db/role permitidos | conexión del adapter |
| `DB_STATEMENT_TIMEOUT_MS` | integer | proyecto define | no | `1..120000` | limita queries |
| `DB_LOCK_TIMEOUT_MS` | integer | proyecto define | no | menor que statement timeout | limita espera de locks |
| `OUTBOX_BATCH_SIZE` | integer | `100` sugerido, no incorporado | no | `1..1000` | presión y throughput |
| `OUTBOX_LEASE_MS` | integer | proyecto define | no | mayor que p99 publish con margen | recovery |
| `JOB_MAX_ATTEMPTS` | integer | por job | no | `1..100` | retry/dead-letter |
| `IDEMPOTENCY_RETENTION` | duration | por endpoint | no | mayor que retry window | replay seguro |

Los defaults operativos no se codifican en SQL porque deben derivarse de SLO, workload y proveedor.

## 7. Dependency bill

| Package/image/tool | Pin exacto | Uso | Licencia | Runtime/build | Fuente oficial |
|---|---|---|---|---|---|
| PostgreSQL server | `18.6` | source of truth | PostgreSQL License | runtime | <https://www.postgresql.org/support/versioning/> |
| `psql` | `18.6` | migrations/tests | PostgreSQL License | build/ops | <https://www.postgresql.org/docs/18/app-psql.html> |

Si se usa imagen, registrar digest OCI además del tag. No usar `postgres:latest`.

## 8. Apply order

Workspace/base vacía:

```powershell
psql $env:DATABASE_URL --set ON_ERROR_STOP=1 --file db/migrations/0001_platform_foundation.up.sql
psql $env:DATABASE_URL --set ON_ERROR_STOP=1 --file db/tests/0001_platform_foundation.test.sql
```

Éxito: ambos comandos devuelven código 0; el test termina en `ROLLBACK` y no deja fixtures.

Base existente:

- comprobar colisiones de schemas/tablas/functions;
- registrar migration en la herramienta elegida, con checksum;
- ejecutar primero en restore reciente y staging;
- medir locks y duración;
- no ejecutar el down migration sobre datos sin autorización destructiva y backup verificado.

## 9. Verification

Gates obligatorios antes de producción:

- migration up/test/down/up en base efímera;
- restricciones negativas para estados, hashes y JSON;
- concurrencia de idempotency same-key/same-hash y different-hash;
- outbox crash antes de publish, después de publish y antes de mark;
- inbox con duplicados concurrentes;
- jobs con lease vencido, retry cap y poison payload;
- grants mínimos con rol migration, app read/write, worker y audit reader separados;
- tenant isolation en cada query y, si se añade RLS, pruebas de pool reset/bypass;
- workload con depth/oldest-age de outbox y jobs;
- backup + PITR restore con invariantes muestreadas;
- `EXPLAIN (ANALYZE, BUFFERS)` sobre queries críticas con datos representativos.

## 10. Reconstruction evidence

Estado: `REBUILD_VERIFIED / CONDITIONED`, evidencia `PG-TX-20260824-V1`.

- 3 archivos materializados con hashes verificados;
- PostgreSQL y `psql` 18.6 reales en cluster efímero con checksums de página;
- `up → test → down → up → test`: PASS;
- restricciones de perfil activo, idempotencia, outbox y auditoría append-only: PASS;
- composición posterior con migration de organización/órdenes y test del repository Go: PASS.

Condiciones restantes: concurrencia multi-session para idempotencia/outbox/inbox/jobs, roles/grants mínimos, RLS si se adopta, dataset/workload representativo, upgrade locks, backup/PITR/restore y failover. La prueba local no convierte el pack en producción ni sustituye un restore drill.
