# Architecture Pack — PostgreSQL Transactional Modular Monolith

> **Estado:** `CANDIDATE_PACK`.  
> **Claim:** baseline transaccional con ownership modular, idempotencia, outbox, inbox, jobs y migrations.  
> **No claim:** SQL siguiente es un contrato de forma; tipos, políticas, volumen y retention se instancian por proyecto.

## 1. Reglas no negociables

- constraints protegen invariantes además de la aplicación;
- una transacción modifica un aggregate/module ownership boundary;
- side effects externos ocurren después del commit mediante outbox;
- cada command externo reintentable tiene idempotency key;
- cada consumer persiste dedup antes/con el side effect local;
- jobs tienen lease, attempt limit, next attempt y dead-letter operable;
- migrations son forward/backward compatibles durante el rollout;
- backup/PITR se valida restaurando.

## 2. Convenciones

```text
UUID/ULID: IDs opacos; generación y monotonicidad documentadas
timestamptz: instantes; business timezone en columna separada
money: bigint minor_units + char(3) currency
version: bigint optimistic concurrency
state: constrained text/enum con migration strategy
jsonb: payload/evidence versionado, no sustituto de schema
```

Schemas iniciales sugeridos: `org`, `catalog`, `crm`, `supply`, `inventory`, `sales`, `service`, `integration`, `audit`, `platform`.

## 3. Idempotency record

```sql
create table platform.idempotency_record (
  scope text not null,
  idempotency_key text not null,
  request_hash bytea not null,
  status text not null check (status in ('processing','completed','failed_terminal')),
  response_code integer,
  response_body jsonb,
  resource_id uuid,
  locked_until timestamptz,
  created_at timestamptz not null default now(),
  expires_at timestamptz not null,
  primary key (scope, idempotency_key)
);
```

Semántica:

1. mismo key + hash distinto → conflict;
2. completed → replay exacto del resultado estable;
3. processing con lease vivo → in-progress/retry-after;
4. lease vencido → reclaim seguro;
5. retention mayor que la ventana de retry del cliente/proveedor.

## 4. Transactional outbox

```sql
create table platform.outbox_event (
  event_id uuid primary key,
  aggregate_type text not null,
  aggregate_id uuid not null,
  aggregate_version bigint not null,
  event_type text not null,
  schema_version integer not null,
  occurred_at timestamptz not null,
  trace_id text,
  tenant_id uuid,
  payload jsonb not null,
  headers jsonb not null default '{}'::jsonb,
  available_at timestamptz not null default now(),
  claimed_by text,
  claimed_until timestamptz,
  attempts integer not null default 0,
  published_at timestamptz,
  last_error_code text
);

create index outbox_ready_idx
  on platform.outbox_event (available_at, occurred_at)
  where published_at is null;

create unique index outbox_aggregate_version_uq
  on platform.outbox_event (aggregate_type, aggregate_id, aggregate_version, event_type);
```

Claim batch con `for update skip locked`; límite de batch y lease explícitos. Publicar fuera de la transacción de claim; marcar éxito en transacción corta. Duplicados son posibles: consumers idempotentes.

## 5. Consumer inbox

```sql
create table platform.consumer_inbox (
  consumer_name text not null,
  event_id uuid not null,
  received_at timestamptz not null default now(),
  completed_at timestamptz,
  result_hash bytea,
  primary key (consumer_name, event_id)
);
```

La inserción inbox y el write de dominio local ocurren en la misma transacción. Un `completed_at` vacío necesita recovery policy; no asumir que insert significa side effect completo.

## 6. Bounded job queue

```sql
create table platform.job (
  job_id uuid primary key,
  queue text not null,
  job_type text not null,
  schema_version integer not null,
  payload jsonb not null,
  priority smallint not null default 0,
  available_at timestamptz not null default now(),
  attempts integer not null default 0,
  max_attempts integer not null,
  claimed_by text,
  claimed_until timestamptz,
  completed_at timestamptz,
  terminal_error_code text,
  created_at timestamptz not null default now()
);

create index job_ready_idx
  on platform.job (queue, priority desc, available_at)
  where completed_at is null and terminal_error_code is null;
```

Backoff con jitter y cap. La queue expone depth, oldest age, attempts, terminal rate y execution duration. Poison jobs no bloquean la cola.

## 7. Ejemplo de stock serializado

```sql
create table inventory.stock_item (
  stock_item_id uuid primary key,
  serial_number text not null unique,
  product_variant_id uuid not null,
  location_id uuid not null,
  state text not null check (state in ('in_transit','quality_hold','available','reserved','sold','service','retired')),
  reservation_id uuid,
  version bigint not null default 0,
  updated_at timestamptz not null default now(),
  check ((state = 'reserved') = (reservation_id is not null))
);
```

Reserva usa conditional update sobre `state='available'` y verifica exactamente una fila. No se hace read-then-write sin lock/condition. Para cantidades fungibles, usar ledger/movements + balance invariant apropiado, no adaptar este esquema a ciegas.

## 8. Optimistic concurrency

```sql
update sales.customer_order
set state = :next_state,
    version = version + 1,
    updated_at = now()
where order_id = :order_id
  and version = :expected_version
  and state = :expected_state;
```

Cero filas significa conflict/stale command, no 404 genérico. El state machine se valida antes y en constraints/triggers sólo cuando ayuden sin ocultar lógica crítica.

## 9. Multi-tenancy

Baseline:

- `organization_id`/`branch_id` forman parte de los aggregates que corresponden;
- authorization construye query scope server-side;
- índices comienzan por tenant cuando las queries lo requieren;
- foreign keys incluyen ownership cuando evita referencias cross-tenant;
- RLS es defense-in-depth sólo si connection/session context es seguro y testeado.

No activar RLS sin pool reset tests, migration/admin paths, background worker identity y pruebas de bypass. RLS no reemplaza checks de acción/recurso.

## 10. Migrations

Secuencia expand/contract:

```text
add nullable/new structure → deploy dual-read/write if needed
→ backfill bounded/resumable with metrics
→ validate constraints concurrently/online where supported
→ switch reads
→ stop old writes
→ remove old structure in later release
```

Cada migration declara lock risk, expected duration, rollback/roll-forward, disk growth y compatible app versions. Nunca ejecutar DDL riesgoso implícitamente al levantar todas las replicas.

## 11. Pool y timeouts

- pool total menor que capacidad DB después de reservar admin/maintenance;
- `statement_timeout`, `lock_timeout`, idle transaction timeout;
- transaction duration observable;
- retries sólo para errores clasificados y operaciones idempotentes;
- connection acquisition forma parte del latency budget;
- slow query logs con bind values redacted.

## 12. Backup, restore y retention

- base backups + WAL/PITR; encryption y access control;
- restore automatizado a entorno aislado;
- medir RPO/RTO reales y consistency de object storage/search projections;
- retention/partitioning para audit/outbox/inbox/jobs;
- legal hold impide eliminación aplicable;
- cryptographic erasure sólo si diseño de keys lo soporta.

## 13. Tests mínimos

- idempotency: same key/same hash, same key/different hash, lease expiry;
- outbox: crash before publish, after publish before mark, duplicate delivery;
- inbox: concurrent duplicate, crash mid-handler;
- jobs: lease loss, retry cap, poison, clock skew assumptions;
- stock/order: competing reservations, cancellation race, payment race;
- tenant: cross-tenant select/update/delete attempts;
- migration: old app/new schema y new app/old-compatible schema;
- restore: sampled business invariants after PITR.

## 14. Gate para alternativas

Agregar Valkey, OpenSearch, Kafka/NATS, Temporal o sharding sólo con:

```text
measured bottleneck or required capability
correctness model
failure/degraded mode
operational owner and SLO
backup/rebuild/reconciliation plan
cost and exit strategy
benchmark/test evidence
```
