# Go Reliable Async Workers

V316: outbox usa attempts monotónico como generación; MarkPublished/Release
exigen generación y lease vigente después del row lock. RemainingLease revalida
al publisher y limita su contexto por evento, incluyendo entradas tardías de lote.
Core0.4.4 y workers0.3.1 se reconstruyen juntos: cambian firmas internas de store;
no mezclar con consumers anteriores ni reiniciar attempts. Sin migración nueva.
Tests PostgreSQL requieren OUTBOX_TEST_DATABASE_URL y
OUTBOX_PROCESSOR_TEST_DATABASE_URL en DB loopback distintas elite_outbox_test_<unique>.
TEST_DATABASE_URL genérica no habilita esas pruebas; no aceptar SKIP como PASS.
Evidencia OUTBOX_GENERATION_FENCING_V316.md: seis rojos DB,tres publisher,
33 tests raíz por corte conectado y4/4 rebuild. Entrega externa sigue al menos
una vez: publisher debe honrar contexto/idempotency key y reconciliar ambigüedad.
No exactly-once remoto, supervisor/alerta, dead-letter policy o producción probados.

## 1. Metadata

```yaml
pack_id: "GO-RELIABLE-ASYNC-WORKERS"
pack_version: "0.3.1"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Inbox transaccional y processors; los jobs exigen generación y lease vigente al confirmar/reprogramar, y revalidan presupuesto antes del handler. No equivale a fencing de efectos remotos."
stacks: ["Go 1.26.7", "PostgreSQL 18.6", "pgx 5.10.0"]
compatible_with: ["GO-ENTERPRISE-BACKEND 0.4.4", "PG-TX-FOUNDATION 0.1.x"]
incompatible_with: []
license_expression: "LicenseRef-Workspace-Owner AND MIT dependency"
upstream_sources: ["https://go.dev", "https://www.postgresql.org/docs/18/", "https://github.com/jackc/pgx"]
verified_at: "2026-09-06"
```

Bloques `AUTHORED`. El outbox ofrece entrega al menos una vez; el publisher externo debe aceptar una idempotency key. Inbox sólo protege una transacción local. Un handler de inbox no puede ejecutar pagos, emails u otros side effects remotos dentro de esa transacción.

## 2. Applicability

### Delta 0.3.0 — scoped claim y transacciones compartidas

`ClaimScoped` filtra tenant/type/schema y campos JSON antes de tomar trabajo;
no autentica al caller. `CompleteJobInTx` y `FailJobInTx` reutilizan el SQL de
las APIs existentes y además exigen payload/queue/type/schema exactos. No hacen
commit: el caller añade su inbox y audit en la misma transacción. Un claim
perdido no autoriza sobrescribir al nuevo. `ExhaustJobInTx` conserva como
LEASE_EXHAUSTED sólo un último intento ya vencido; nunca reinicia generaciones.
Estos helpers son AUTHORED, no código publicado por AWS. Atomicidad local según
[PostgreSQL 18](https://www.postgresql.org/docs/18/tutorial-transactions.html),
consultado 2026-09-06; no garantiza efectos remotos exactly-once.

APIs 0.2.0 permanecen compatibles. El worker específico está en
PYTHON-META-WHATSAPP-CLOUD-ADAPTER 0.11.0: no se duplica aquí. No hay migración
ni dependencia nueva; no promover un publisher/outbox por estas pruebas.
Evidencia: `reconstruction_evidence/WHATSAPP_SCOPED_JOB_COMPLETION_V277.md`.

### Delta 0.2.0 — job claim fencing

Mantenimiento `AUTHORED`, no código copiado de AWS. La guía oficial [Amazon Builders' Library](https://aws.amazon.com/es/builders-library/leader-election-in-distributed-systems/) exige revisar concesiones antes de efectos y considerar pausas/red; [PostgreSQL 18 SELECT](https://www.postgresql.org/docs/18/sql-select.html) documenta row locking y SKIP LOCKED para consumidores de cola. Consultadas 2026-09-06. La implementación y sus pruebas son responsabilidad local, no una certificación de esas empresas.

`JobStore.Complete` y `Fail` requieren ahora el `Attempts` exacto devuelto por `Claim`; `RemainingLease` es obligatorio en el processor. Actualizar conjuntamente store, interface, mocks y callers; drenar/parar binarios 0.1.x antes de activar 0.2.0. No mezclar workers antiguos: conservan el SQL inseguro. No hay nueva migración, dependencia ni cola. No restablecer attempts ni reutilizar IDs de jobs: la generación deja de ser fiable si un operador reescribe su historia.

El check final toma el row lock antes de comparar generación y reloj PostgreSQL. Cada handler recibe un contexto acotado por tiempo restante menos el round trip de validación; cooperar con cancelación sigue siendo obligatorio. Un side effect que ya salió no se deshace: exige contrato idempotente, fence propio y reconciliación. El lease no constituye autorización empresarial.

Límites explícitos: no implementa el worker/routing WhatsApp; no cambia el publisher/outbox ni demuestra fencing de sus claims; un job cuyo último intento murió debe conservarse para reconciliación operativa, no reiniciarse a ciegas. No hay renovación automática de lease. Evidencia: `reconstruction_evidence/JOB_CLAIM_GENERATION_AND_EXPIRY_V275.md`.

Use when local transactions must publish events or execute retriable work without dual writes. Reject it for exactly-once claims, long work without lease renewal design, or remote effects whose provider lacks idempotency/reconciliation. It complements, rather than replaces, an admitted transport.

## 3. Architecture contract

Outbox, inbox and job state live in PostgreSQL. Claims use bounded batches and leases; success/retry/dead states are explicit. Delivery is at least once, so handlers/publishers carry stable idempotency keys. Remote side effects never occur inside an inbox database transaction. Shutdown stops new claims, drains bounded work and releases/lets leases expire safely.

## 4. Exact file manifest

```text
CREATE internal/platform/postgres/inbox.go
CREATE internal/platform/postgres/inbox_integration_test.go
CREATE internal/platform/postgres/jobs.go
CREATE internal/platform/postgres/jobs_integration_test.go
CREATE internal/platform/workers/processors.go
CREATE internal/platform/workers/processors_test.go
```

## 5. Materialization blocks

### FILE: `internal/platform/postgres/inbox.go`

```yaml
block_id: "GO-ASYNC:inbox:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "75974f50b0e7ffb8f7244d6cf8f96647aaf9932cb70488af986b1e30514db2db"
variables: []
secrets_allowed: false
```

````go
package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Inbox struct{ pool *pgxpool.Pool }

func NewInbox(pool *pgxpool.Pool) *Inbox { return &Inbox{pool: pool} }

// Process executes local database effects and the inbox receipt atomically.
// Returning processed=false means a completed receipt already existed.
func (i *Inbox) Process(ctx context.Context, tenantID, consumerName, eventID string, handler func(context.Context, pgx.Tx) (string, error)) (processed bool, err error) {
	if tenantID == "" || consumerName == "" || len(consumerName) > 128 || eventID == "" || handler == nil {
		return false, fmt.Errorf("invalid inbox parameters")
	}
	tx, err := i.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return false, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var inserted bool
	err = tx.QueryRow(ctx, `
with receipt as (
  insert into platform.consumer_inbox(tenant_id,consumer_name,event_id)
  values($1,$2,$3)
  on conflict do nothing
  returning true
)
select coalesce((select true from receipt),false)`, tenantID, consumerName, eventID).Scan(&inserted)
	if err != nil {
		return false, err
	}
	if !inserted {
		var completed bool
		if err := tx.QueryRow(ctx, `select completed_at is not null from platform.consumer_inbox where tenant_id=$1 and consumer_name=$2 and event_id=$3`, tenantID, consumerName, eventID).Scan(&completed); err != nil {
			return false, err
		}
		if !completed {
			return false, fmt.Errorf("inbox receipt exists but is incomplete")
		}
		if err := tx.Commit(ctx); err != nil {
			return false, err
		}
		return false, nil
	}
	resultHash, err := handler(ctx, tx)
	if err != nil {
		return false, err
	}
	if resultHash == "" {
		resultHash = fmt.Sprintf("%064x", 0)
	}
	result, err := tx.Exec(ctx, `update platform.consumer_inbox set completed_at=clock_timestamp(),result_sha256_hex=$4 where tenant_id=$1 and consumer_name=$2 and event_id=$3 and completed_at is null`, tenantID, consumerName, eventID, resultHash)
	if err != nil {
		return false, err
	}
	if result.RowsAffected() != 1 {
		return false, fmt.Errorf("inbox completion lost")
	}
	if err := tx.Commit(ctx); err != nil {
		return false, err
	}
	return true, nil
}
````

### FILE: `internal/platform/postgres/inbox_integration_test.go`

```yaml
block_id: "GO-ASYNC:inbox-test:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "becfc0b889fb81ed1fb20db2a90933ef7b60b97eb69dc30a6d3df4cc2d0929bd"
variables: []
secrets_allowed: false
```

````go
package postgres

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestInboxAtomicDedupAndRollback(t *testing.T) {
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
	tenant := "018f4d4a-7b36-7a21-8d10-2f4c54c28c01"
	event := "018f4d4a-7b36-7a21-8d10-2f4c54c28c02"
	_, _ = pool.Exec(ctx, `delete from platform.consumer_inbox where tenant_id=$1`, tenant)
	_, _ = pool.Exec(ctx, `delete from platform.tenant where tenant_id=$1`, tenant)
	defer func() {
		_, _ = pool.Exec(ctx, `delete from platform.consumer_inbox where tenant_id=$1`, tenant)
		_, _ = pool.Exec(ctx, `delete from platform.tenant where tenant_id=$1`, tenant)
	}()
	if _, err := pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,'inbox-test','Inbox Test','Inbox')`, tenant); err != nil {
		t.Fatal(err)
	}
	inbox := NewInbox(pool)
	calls := 0
	failing := func(context.Context, pgx.Tx) (string, error) { calls++; return "", errors.New("handler failed") }
	if _, err := inbox.Process(ctx, tenant, "inventory", event, failing); err == nil {
		t.Fatal("expected handler failure")
	}
	var count int
	if err := pool.QueryRow(ctx, `select count(*) from platform.consumer_inbox where tenant_id=$1`, tenant).Scan(&count); err != nil || count != 0 {
		t.Fatalf("failed handler left receipt: count=%d err=%v", count, err)
	}
	handler := func(ctx context.Context, tx pgx.Tx) (string, error) {
		calls++
		_, err := tx.Exec(ctx, `update platform.tenant set display_name='Inbox Applied' where tenant_id=$1`, tenant)
		return "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", err
	}
	processed, err := inbox.Process(ctx, tenant, "inventory", event, handler)
	if err != nil || !processed {
		t.Fatalf("first process: %v %v", processed, err)
	}
	processed, err = inbox.Process(ctx, tenant, "inventory", event, handler)
	if err != nil || processed {
		t.Fatalf("duplicate process: %v %v", processed, err)
	}
	if calls != 2 {
		t.Fatalf("handler calls=%d want 2 (one rollback, one commit)", calls)
	}
}
````

### FILE: `internal/platform/postgres/jobs.go`

```yaml
block_id: "GO-ASYNC:jobs:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "417bff0d640ea13146a873e1e556e61cc428d507ab714ef10ef4d073ec8f8345"
variables: []
secrets_allowed: false
```

````go
package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Job struct {
	TenantID, JobID, Queue, JobType      string
	SchemaVersion, Attempts, MaxAttempts int
	Payload                              json.RawMessage
}
type Jobs struct{ pool *pgxpool.Pool }

var ErrJobClaimLost = errors.New("job claim lost")

// ExhaustJobInTx preserves a crashed final attempt as terminal evidence. It
// cannot expire a live claim, reset attempts, or silently start another effect.
func ExhaustJobInTx(ctx context.Context, tx pgx.Tx, job Job) error {
	if tx == nil || job.Attempts < 1 || !json.Valid(job.Payload) {
		return ErrJobClaimLost
	}
	tag, err := tx.Exec(ctx, lockedJob+`update platform.job j set claimed_by=null,claimed_until=null,terminal_error_code='LEASE_EXHAUSTED'
 from owned o where j.tenant_id=o.tenant_id and j.job_id=o.job_id and o.attempts=$3
 and o.attempts>=o.max_attempts and o.claimed_until<clock_timestamp()
 and o.queue=$4 and o.job_type=$5 and o.schema_version=$6 and o.payload=$7::jsonb`, job.TenantID, job.JobID, job.Attempts, job.Queue, job.JobType, job.SchemaVersion, job.Payload)
	if err != nil {
		return err
	}
	if tag.RowsAffected() != 1 {
		return ErrJobClaimLost
	}
	return nil
}

func NewJobs(pool *pgxpool.Pool) *Jobs { return &Jobs{pool: pool} }

// JobScope is server-side selection, not authentication or tenant authorization.
type JobScope struct {
	TenantID, JobType string
	SchemaVersion     int
	PayloadMatch      json.RawMessage
}

func (j *Jobs) ClaimScoped(ctx context.Context, queue, worker string, lease time.Duration, limit int, scope JobScope) ([]Job, error) {
	var fields map[string]json.RawMessage
	if scope.TenantID == "" || scope.JobType == "" || scope.SchemaVersion < 1 || len(scope.PayloadMatch) > 4096 || json.Unmarshal(scope.PayloadMatch, &fields) != nil || fields == nil {
		return nil, fmt.Errorf("invalid job scope")
	}
	return j.claim(ctx, queue, worker, lease, limit, &scope)
}

func (j *Jobs) Claim(ctx context.Context, queue, worker string, lease time.Duration, limit int) ([]Job, error) {
	return j.claim(ctx, queue, worker, lease, limit, nil)
}

func (j *Jobs) claim(ctx context.Context, queue, worker string, lease time.Duration, limit int, scope *JobScope) ([]Job, error) {
	if j == nil || j.pool == nil || queue == "" || worker == "" || lease < time.Microsecond || limit < 1 || limit > 1000 {
		return nil, fmt.Errorf("invalid job claim parameters")
	}
	var tenant any
	jobType, schema, match := "", 0, json.RawMessage(`{}`)
	if scope != nil {
		tenant = scope.TenantID
		jobType = scope.JobType
		schema = scope.SchemaVersion
		match = scope.PayloadMatch
	}
	rows, err := j.pool.Query(ctx, `with candidates as (
select tenant_id,job_id from platform.job where queue=$1 and completed_at is null and terminal_error_code is null and available_at<=clock_timestamp() and attempts<max_attempts and (claimed_until is null or claimed_until<clock_timestamp())
and ($5::uuid is null or tenant_id=$5) and ($6::text='' or job_type=$6) and ($7::integer=0 or schema_version=$7) and payload @> $8::jsonb
order by priority desc,available_at,job_id for update skip locked limit $2)
update platform.job j set claimed_by=$3,claimed_until=clock_timestamp()+$4::interval,attempts=attempts+1 from candidates c where j.tenant_id=c.tenant_id and j.job_id=c.job_id returning j.tenant_id,j.job_id,j.queue,j.job_type,j.schema_version,j.payload,j.attempts,j.max_attempts`, queue, limit, worker, lease.String(), tenant, jobType, schema, match)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	jobs := make([]Job, 0, limit)
	for rows.Next() {
		var job Job
		if err := rows.Scan(&job.TenantID, &job.JobID, &job.Queue, &job.JobType, &job.SchemaVersion, &job.Payload, &job.Attempts, &job.MaxAttempts); err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}
	return jobs, rows.Err()
}

// Attempts is a monotonically increasing claim generation, not an operator-
// resettable counter. A worker name alone is not a fencing token.
// Lock first, then check the DB clock: waiting for a row lock may expire a lease.
const lockedJob = `with owned as materialized (
 select tenant_id,job_id,claimed_by,claimed_until,attempts,max_attempts,queue,job_type,schema_version,payload
 from platform.job where tenant_id=$1 and job_id=$2
 and completed_at is null and terminal_error_code is null for update
) `

func (j *Jobs) Complete(ctx context.Context, tenant, jobID, worker string, attempt int) error {
	if j == nil || j.pool == nil || worker == "" || attempt < 1 {
		return ErrJobClaimLost
	}
	return completeJob(ctx, j.pool, tenant, jobID, worker, attempt, nil)
}

type jobExecutor interface {
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	QueryRow(context.Context, string, ...any) pgx.Row
}

// CompleteJobInTx does not commit. The caller atomically records its durable
// receipt/ack in the same transaction. Full claimed metadata must still match.
func CompleteJobInTx(ctx context.Context, tx pgx.Tx, job Job, worker string) error {
	if tx == nil || worker == "" || job.Attempts < 1 || job.Queue == "" || job.JobType == "" || job.SchemaVersion < 1 || len(job.Payload) > 65536 || !json.Valid(job.Payload) {
		return ErrJobClaimLost
	}
	return completeJob(ctx, tx, job.TenantID, job.JobID, worker, job.Attempts, &job)
}

func completeJob(ctx context.Context, db jobExecutor, tenant, jobID, worker string, attempt int, expected *Job) error {
	var payload any
	queue, kind, schema := "", "", 0
	if expected != nil {
		payload = expected.Payload
		queue = expected.Queue
		kind = expected.JobType
		schema = expected.SchemaVersion
	}
	r, err := db.Exec(ctx, lockedJob+`update platform.job j set completed_at=clock_timestamp(),claimed_by=null,claimed_until=null
 from owned o where j.tenant_id=o.tenant_id and j.job_id=o.job_id
 and o.claimed_by=$3 and o.attempts=$4 and o.claimed_until>clock_timestamp()
 and ($5::jsonb is null or (o.payload=$5 and o.queue=$6 and o.job_type=$7 and o.schema_version=$8))`, tenant, jobID, worker, attempt, payload, queue, kind, schema)
	if err != nil {
		return err
	}
	if r.RowsAffected() != 1 {
		return ErrJobClaimLost
	}
	return nil
}

func (j *Jobs) Fail(ctx context.Context, tenant, jobID, worker, code string, attempt int, retryAfter time.Duration) (bool, error) {
	if j == nil || j.pool == nil || worker == "" || attempt < 1 || code == "" || retryAfter < 0 {
		return false, fmt.Errorf("invalid job failure parameters")
	}
	return failJob(ctx, j.pool, tenant, jobID, worker, code, attempt, retryAfter, nil)
}

// FailJobInTx shares the transition implementation with Jobs.Fail. The caller
// can append a safe failure audit and update its inbox in the same transaction.
func FailJobInTx(ctx context.Context, tx pgx.Tx, job Job, worker, code string, retryAfter time.Duration) (bool, error) {
	if tx == nil || worker == "" || job.Attempts < 1 || job.Queue == "" || job.JobType == "" || job.SchemaVersion < 1 || code == "" || retryAfter < 0 || len(job.Payload) > 65536 || !json.Valid(job.Payload) {
		return false, ErrJobClaimLost
	}
	return failJob(ctx, tx, job.TenantID, job.JobID, worker, code, job.Attempts, retryAfter, &job)
}

func failJob(ctx context.Context, db jobExecutor, tenant, jobID, worker, code string, attempt int, retryAfter time.Duration, expected *Job) (bool, error) {
	var payload any
	queue, kind, schema := "", "", 0
	if expected != nil {
		payload = expected.Payload
		queue = expected.Queue
		kind = expected.JobType
		schema = expected.SchemaVersion
	}
	var terminal bool
	err := db.QueryRow(ctx, lockedJob+`update platform.job j set claimed_by=null,claimed_until=null,terminal_error_code=case when o.attempts>=o.max_attempts then $4 else null end,available_at=case when o.attempts>=o.max_attempts then j.available_at else clock_timestamp()+$5::interval end
 from owned o where j.tenant_id=o.tenant_id and j.job_id=o.job_id
 and o.claimed_by=$3 and o.attempts=$6 and o.claimed_until>clock_timestamp()
 and ($7::jsonb is null or (o.payload=$7 and o.queue=$8 and o.job_type=$9 and o.schema_version=$10))
 returning j.terminal_error_code is not null`, tenant, jobID, worker, code, retryAfter.String(), attempt, payload, queue, kind, schema).Scan(&terminal)
	if err != nil {
		return false, fmt.Errorf("job claim lost or update failed: %w", err)
	}
	return terminal, nil
}

// RemainingLease is a preflight check, not permission to perform arbitrary
// remote effects. Handlers still require transactional/idempotent effects and
// reconciliation; cancellation cannot undo an already-issued external request.
func (j *Jobs) RemainingLease(ctx context.Context, tenant, jobID, worker string, attempt int) (time.Duration, error) {
	if j == nil || j.pool == nil || worker == "" || attempt < 1 {
		return 0, ErrJobClaimLost
	}
	var micros int64
	err := j.pool.QueryRow(ctx, `select floor(extract(epoch from (claimed_until-clock_timestamp()))*1000000)::bigint
 from platform.job where tenant_id=$1 and job_id=$2 and claimed_by=$3 and attempts=$4
 and completed_at is null and terminal_error_code is null and claimed_until>clock_timestamp()`, tenant, jobID, worker, attempt).Scan(&micros)
	if err != nil {
		return 0, fmt.Errorf("job lease unavailable: %w", err)
	}
	if micros <= 0 || micros > int64((time.Duration(1<<63-1))/time.Microsecond) {
		return 0, ErrJobClaimLost
	}
	return time.Duration(micros) * time.Microsecond, nil
}
````

### FILE: `internal/platform/postgres/jobs_integration_test.go`

```yaml
block_id: "GO-ASYNC:jobs-test:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "97ffdd220c22757d9ac86994eebea4870001f5ab5af9196892eaad50c9a207aa"
variables: []
secrets_allowed: false
```

````go
package postgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"strings"
	"testing"
	"time"
)

// This fault-injection suite only uses a dedicated local test database.
func TestJobsExpiredAndReusedWorkerClaims(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		t.Fatal(err)
	}
	if config.ConnConfig.Host != "127.0.0.1" || !strings.HasPrefix(config.ConnConfig.Database, "elite_jobs_") {
		t.Fatal("dedicated local elite_jobs_ database required")
	}
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	for _, takeover := range []bool{false, true} {
		for _, operation := range []string{"complete", "fail"} {
			name := operation
			if takeover {
				name += "-reused-worker"
			} else {
				name += "-expired"
			}
			t.Run(name, func(t *testing.T) {
				var tenant, job string
				if err := pool.QueryRow(ctx, `select gen_random_uuid()::text,gen_random_uuid()::text`).Scan(&tenant, &job); err != nil {
					t.Fatal(err)
				}
				if _, err := pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,$2,'Synthetic lease test','Lease')`, tenant, "lease-"+job); err != nil {
					t.Fatal(err)
				}
				if _, err := pool.Exec(ctx, `insert into platform.job(tenant_id,job_id,queue,job_type,schema_version,payload,max_attempts) values($1,$2,$3,'lease.test',1,'{}',3)`, tenant, job, job); err != nil {
					t.Fatal(err)
				}
				store := NewJobs(pool)
				first, err := store.Claim(ctx, job, "same-worker", time.Minute, 1)
				if err != nil || len(first) != 1 {
					t.Fatalf("first claim: %v %v", first, err)
				}
				// Deterministic DB-clock expiry, not a timing-sensitive sleep.
				if _, err := pool.Exec(ctx, `update platform.job set claimed_until=clock_timestamp()-interval '1 second' where tenant_id=$1 and job_id=$2`, tenant, job); err != nil {
					t.Fatal(err)
				}
				if takeover {
					second, err := store.Claim(ctx, job, "same-worker", time.Minute, 1)
					if err != nil || len(second) != 1 || second[0].Attempts != 2 {
						t.Fatalf("takeover: %v %v", second, err)
					}
				}
				if _, err := store.RemainingLease(ctx, tenant, job, "same-worker", first[0].Attempts); err == nil {
					t.Fatal("stale preflight accepted")
				}
				if operation == "complete" {
					err = store.Complete(ctx, tenant, job, "same-worker", first[0].Attempts)
				} else {
					_, err = store.Fail(ctx, tenant, job, "same-worker", "STALE", first[0].Attempts, 0)
				}
				if err == nil {
					t.Fatal("stale execution changed job state")
				}
				var untouched bool
				if err := pool.QueryRow(ctx, `select completed_at is null and terminal_error_code is null and claimed_by='same-worker' from platform.job where tenant_id=$1 and job_id=$2`, tenant, job).Scan(&untouched); err != nil || !untouched {
					t.Fatalf("lost current claim: %v %v", untouched, err)
				}
				if takeover {
					if remaining, err := store.RemainingLease(ctx, tenant, job, "same-worker", 2); err != nil || remaining <= 0 || remaining > time.Minute {
						t.Fatalf("current lease %v %v", remaining, err)
					}
					if operation == "complete" {
						err = store.Complete(ctx, tenant, job, "same-worker", 2)
					} else {
						_, err = store.Fail(ctx, tenant, job, "same-worker", "CURRENT", 2, 0)
					}
					if err != nil {
						t.Fatalf("current generation rejected: %v", err)
					}
				}
			})
		}
	}
}

func TestJobsOwnershipRetryAndTerminal(t *testing.T) {
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
	tenant := "018f4d4a-7b36-7a21-8d10-2f4c54c28d01"
	jobID := "018f4d4a-7b36-7a21-8d10-2f4c54c28d02"
	_, _ = pool.Exec(ctx, `delete from platform.job where tenant_id=$1`, tenant)
	_, _ = pool.Exec(ctx, `delete from platform.tenant where tenant_id=$1`, tenant)
	defer func() {
		_, _ = pool.Exec(ctx, `delete from platform.job where tenant_id=$1`, tenant)
		_, _ = pool.Exec(ctx, `delete from platform.tenant where tenant_id=$1`, tenant)
	}()
	_, err = pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1,'jobs-test','Jobs Test','Jobs')`, tenant)
	if err != nil {
		t.Fatal(err)
	}
	_, err = pool.Exec(ctx, `insert into platform.job(tenant_id,job_id,queue,job_type,schema_version,payload,max_attempts)values($1,$2,'email','welcome',1,'{}',2)`, tenant, jobID)
	if err != nil {
		t.Fatal(err)
	}
	store := NewJobs(pool)
	first, err := store.Claim(ctx, "email", "worker-a", time.Minute, 1)
	if err != nil || len(first) != 1 || first[0].Attempts != 1 {
		t.Fatalf("first claim %+v %v", first, err)
	}
	if err := store.Complete(ctx, tenant, jobID, "worker-b", first[0].Attempts); err == nil {
		t.Fatal("wrong worker completed job")
	}
	terminal, err := store.Fail(ctx, tenant, jobID, "worker-a", "TEMP", first[0].Attempts, 0)
	if err != nil || terminal {
		t.Fatalf("first failure terminal=%v err=%v", terminal, err)
	}
	second, err := store.Claim(ctx, "email", "worker-b", time.Minute, 1)
	if err != nil || len(second) != 1 || second[0].Attempts != 2 {
		t.Fatalf("second claim %+v %v", second, err)
	}
	terminal, err = store.Fail(ctx, tenant, jobID, "worker-b", "EXHAUSTED", second[0].Attempts, 0)
	if err != nil || !terminal {
		t.Fatalf("terminal failure=%v err=%v", terminal, err)
	}
	third, err := store.Claim(ctx, "email", "worker-c", time.Minute, 1)
	if err != nil || len(third) != 0 {
		t.Fatalf("terminal job reclaimed %+v %v", third, err)
	}
}

func TestJobsExpiryWhileWaitingForRowLock(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		t.Fatal(err)
	}
	if config.ConnConfig.Host != "127.0.0.1" || !strings.HasPrefix(config.ConnConfig.Database, "elite_jobs_") {
		t.Fatal("dedicated local elite_jobs_ database required")
	}
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	for _, operation := range []string{"complete", "fail"} {
		t.Run(operation, func(t *testing.T) {
			var tenant, job string
			if err := pool.QueryRow(ctx, `select gen_random_uuid()::text,gen_random_uuid()::text`).Scan(&tenant, &job); err != nil {
				t.Fatal(err)
			}
			if _, err := pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,$2,'Synthetic lock wait','Lease')`, tenant, "wait-"+job); err != nil {
				t.Fatal(err)
			}
			if _, err := pool.Exec(ctx, `insert into platform.job(tenant_id,job_id,queue,job_type,schema_version,payload,max_attempts) values($1,$2,$3,'lease.test',1,'{}',3)`, tenant, job, job); err != nil {
				t.Fatal(err)
			}
			store := NewJobs(pool)
			claimed, err := store.Claim(ctx, job, "wait-worker", time.Second, 1)
			if err != nil || len(claimed) != 1 {
				t.Fatalf("claim %v %v", claimed, err)
			}
			tx, err := pool.Begin(ctx)
			if err != nil {
				t.Fatal(err)
			}
			defer tx.Rollback(ctx)
			if _, err := tx.Exec(ctx, `select job_id from platform.job where tenant_id=$1 and job_id=$2 for update`, tenant, job); err != nil {
				t.Fatal(err)
			}
			done := make(chan error, 1)
			go func() {
				if operation == "complete" {
					done <- store.Complete(ctx, tenant, job, "wait-worker", 1)
				} else {
					_, err := store.Fail(ctx, tenant, job, "wait-worker", "DELAYED", 1, 0)
					done <- err
				}
			}()
			// Confirm actual server lock contention, not merely scheduling a goroutine.
			for {
				var waiting bool
				if err := pool.QueryRow(ctx, `select exists(select 1 from pg_stat_activity where datname=current_database() and pid<>pg_backend_pid() and state='active' and wait_event_type='Lock' and query like 'with owned as materialized (%')`).Scan(&waiting); err != nil {
					t.Fatal(err)
				}
				if waiting {
					break
				}
				select {
				case err := <-done:
					t.Fatalf("transition did not wait: %v", err)
				case <-ctx.Done():
					t.Fatal(ctx.Err())
				case <-time.After(5 * time.Millisecond):
				}
			}
			// Hold the row lock through expiry without changing the row version.
			if _, err := tx.Exec(ctx, `select pg_sleep(greatest(0,extract(epoch from (claimed_until-clock_timestamp())))+0.02) from platform.job where tenant_id=$1 and job_id=$2`, tenant, job); err != nil {
				t.Fatal(err)
			}
			if err := tx.Commit(ctx); err != nil {
				t.Fatal(err)
			}
			if err := <-done; err == nil {
				t.Fatal("lease expired during lock wait but transition succeeded")
			}
			var untouched bool
			if err := pool.QueryRow(ctx, `select completed_at is null and terminal_error_code is null and claimed_by='wait-worker' from platform.job where tenant_id=$1 and job_id=$2`, tenant, job).Scan(&untouched); err != nil || !untouched {
				t.Fatalf("changed expired job: %v %v", untouched, err)
			}
		})
	}
}
````

### FILE: `internal/platform/workers/processors.go`

```yaml
block_id: "GO-ASYNC:processors:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "7bff6315aae3f6c0a2b04253860c9056649ee6527208c5191b9f84a560c62314"
variables: []
secrets_allowed: false
```

````go
package workers

import (
	"context"
	db "elite.local/enterprise/internal/platform/postgres"
	"fmt"
	"time"
)

type OutboxStore interface {
	Claim(context.Context, string, time.Duration, int) ([]db.OutboxEvent, error)
	RemainingLease(context.Context, string, string, string, int) (time.Duration, error)
	MarkPublished(context.Context, string, string, string, int) error
	Release(context.Context, string, string, string, string, int, time.Duration) error
}
type Publisher interface {
	Publish(context.Context, db.OutboxEvent, string) error
}
type OutboxProcessor struct {
	Store             OutboxStore
	Publisher         Publisher
	WorkerID          string
	Lease, RetryDelay time.Duration
	BatchSize         int
}

func (p OutboxProcessor) ProcessOnce(ctx context.Context) (int, error) {
	if p.Store == nil || p.Publisher == nil || p.WorkerID == "" || p.Lease < time.Microsecond || p.RetryDelay < 0 || p.BatchSize < 1 || p.BatchSize > 1000 {
		return 0, fmt.Errorf("invalid outbox processor")
	}
	events, err := p.Store.Claim(ctx, p.WorkerID, p.Lease, p.BatchSize)
	if err != nil {
		return 0, err
	}
	published := 0
	for _, event := range events {
		started := time.Now()
		remaining, err := p.Store.RemainingLease(ctx, event.TenantID, event.EventID, p.WorkerID, event.Attempts)
		remaining -= time.Since(started)
		if err != nil {
			return published, err
		}
		if remaining <= 0 || ctx.Err() != nil {
			return published, db.ErrOutboxClaimLost
		}
		publishCtx, cancel := context.WithTimeout(ctx, remaining)
		publishErr := p.Publisher.Publish(publishCtx, event, event.EventID)
		cancel()
		if publishErr != nil {
			if releaseErr := p.Store.Release(ctx, event.TenantID, event.EventID, p.WorkerID, "PUBLISH_FAILED", event.Attempts, p.RetryDelay); releaseErr != nil {
				return published, releaseErr
			}
			continue
		}
		if err := p.Store.MarkPublished(ctx, event.TenantID, event.EventID, p.WorkerID, event.Attempts); err != nil {
			return published, err
		}
		published++
	}
	return published, nil
}

type JobStore interface {
	Claim(context.Context, string, string, time.Duration, int) ([]db.Job, error)
	RemainingLease(context.Context, string, string, string, int) (time.Duration, error)
	Complete(context.Context, string, string, string, int) error
	Fail(context.Context, string, string, string, string, int, time.Duration) (bool, error)
}
type JobHandler interface {
	Handle(context.Context, db.Job) error
}
type JobProcessor struct {
	Store             JobStore
	Handler           JobHandler
	Queue, WorkerID   string
	Lease, RetryDelay time.Duration
	BatchSize         int
}

func (p JobProcessor) ProcessOnce(ctx context.Context) (int, error) {
	if p.Store == nil || p.Handler == nil || p.Queue == "" || p.WorkerID == "" {
		return 0, fmt.Errorf("invalid job processor")
	}
	jobs, err := p.Store.Claim(ctx, p.Queue, p.WorkerID, p.Lease, p.BatchSize)
	if err != nil {
		return 0, err
	}
	completed := 0
	for _, job := range jobs {
		// Subtract the complete DB round trip conservatively. A delayed response
		// must not extend a lease, and later batch entries may have expired.
		started := time.Now()
		remaining, err := p.Store.RemainingLease(ctx, job.TenantID, job.JobID, p.WorkerID, job.Attempts)
		remaining -= time.Since(started)
		if err != nil {
			return completed, err
		}
		if remaining <= 0 || ctx.Err() != nil {
			return completed, db.ErrJobClaimLost
		}
		jobCtx, cancel := context.WithTimeout(ctx, remaining)
		handleErr := p.Handler.Handle(jobCtx, job)
		cancel()
		if handleErr != nil {
			if _, failErr := p.Store.Fail(ctx, job.TenantID, job.JobID, p.WorkerID, "HANDLER_FAILED", job.Attempts, p.RetryDelay); failErr != nil {
				return completed, failErr
			}
			continue
		}
		if err := p.Store.Complete(ctx, job.TenantID, job.JobID, p.WorkerID, job.Attempts); err != nil {
			return completed, err
		}
		completed++
	}
	return completed, nil
}
````

### FILE: `internal/platform/workers/processors_test.go`

```yaml
block_id: "GO-ASYNC:processors-test:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "c07d2a932b516fda912226bc44ff3dd55604e6dd6f641ce5ba0ebb34ed10e00e"
variables: []
secrets_allowed: false
```

````go
package workers

import (
	"context"
	db "elite.local/enterprise/internal/platform/postgres"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type fakeOutbox struct {
	events           []db.OutboxEvent
	marked, released int
	remaining        time.Duration
	remainingByEvent map[string]time.Duration
}

func (f *fakeOutbox) Claim(context.Context, string, time.Duration, int) ([]db.OutboxEvent, error) {
	return f.events, nil
}
func (f *fakeOutbox) RemainingLease(_ context.Context, _, event, _ string, attempt int) (time.Duration, error) {
	if attempt != 1 {
		return 0, errors.New("wrong outbox preflight generation")
	}
	if f.remainingByEvent != nil {
		return f.remainingByEvent[event], nil
	}
	return f.remaining, nil
}
func (f *fakeOutbox) MarkPublished(_ context.Context, _, _, _ string, attempt int) error {
	if attempt != 1 {
		return errors.New("wrong outbox completion generation")
	}
	f.marked++
	return nil
}
func (f *fakeOutbox) Release(_ context.Context, _, _, _, _ string, attempt int, _ time.Duration) error {
	if attempt != 1 {
		return errors.New("wrong outbox release generation")
	}
	f.released++
	return nil
}

type fakePublisher struct{ fail string }

func (f fakePublisher) Publish(_ context.Context, e db.OutboxEvent, key string) error {
	if key != e.EventID {
		panic("missing idempotency key")
	}
	if e.EventID == f.fail {
		return errors.New("down")
	}
	return nil
}

type publisherFunc func(context.Context, db.OutboxEvent, string) error

func (f publisherFunc) Publish(ctx context.Context, e db.OutboxEvent, key string) error {
	return f(ctx, e, key)
}

func TestOutboxProcessorRejectsExpiredBudget(t *testing.T) {
	out := &fakeOutbox{events: []db.OutboxEvent{{TenantID: "t", EventID: "old", Attempts: 1}}, remaining: 0}
	called := false
	publisher := publisherFunc(func(context.Context, db.OutboxEvent, string) error { called = true; return nil })
	n, err := (OutboxProcessor{Store: out, Publisher: publisher, WorkerID: "w", Lease: time.Second, RetryDelay: time.Second, BatchSize: 1}).ProcessOnce(context.Background())
	if called || n != 0 || err == nil || out.marked != 0 || out.released != 0 {
		t.Fatalf("expired claim reached publisher: called=%v n=%d err=%v", called, n, err)
	}
}

func TestOutboxProcessorBoundsPublisherContext(t *testing.T) {
	out := &fakeOutbox{events: []db.OutboxEvent{{TenantID: "t", EventID: "fresh", Attempts: 1}}, remaining: 50 * time.Millisecond}
	publisher := publisherFunc(func(ctx context.Context, _ db.OutboxEvent, _ string) error {
		deadline, ok := ctx.Deadline()
		if !ok || time.Until(deadline) > 50*time.Millisecond {
			t.Error("publisher context exceeds remaining lease")
		}
		return nil
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if n, err := (OutboxProcessor{Store: out, Publisher: publisher, WorkerID: "w", Lease: time.Second, RetryDelay: time.Second, BatchSize: 1}).ProcessOnce(ctx); err != nil || n != 1 {
		t.Fatalf("fresh publish n=%d err=%v", n, err)
	}
}

func TestOutboxProcessorRechecksLaterBatchEntry(t *testing.T) {
	out := &fakeOutbox{events: []db.OutboxEvent{{TenantID: "t", EventID: "fresh", Attempts: 1}, {TenantID: "t", EventID: "late", Attempts: 1}}, remainingByEvent: map[string]time.Duration{"fresh": time.Second, "late": 0}}
	calls := 0
	publisher := publisherFunc(func(context.Context, db.OutboxEvent, string) error { calls++; return nil })
	n, err := (OutboxProcessor{Store: out, Publisher: publisher, WorkerID: "w", Lease: time.Second, RetryDelay: time.Second, BatchSize: 2}).ProcessOnce(context.Background())
	if n != 1 || err == nil || calls != 1 || out.marked != 1 {
		t.Fatalf("late batch entry published: n=%d err=%v calls=%d", n, err, calls)
	}
}

func TestOutboxProcessorPostgresRetryKeepsEventKey(t *testing.T) {
	// Separate database from the postgres package: Claim is intentionally global,
	// and Go can run package suites concurrently. Never use TEST_DATABASE_URL.
	raw := os.Getenv("OUTBOX_PROCESSOR_TEST_DATABASE_URL")
	if raw == "" {
		t.Skip("OUTBOX_PROCESSOR_TEST_DATABASE_URL is not set; separate isolated database required")
	}
	cfg, err := pgxpool.ParseConfig(raw)
	if err != nil || cfg.ConnConfig.Host != "127.0.0.1" || !strings.HasPrefix(cfg.ConnConfig.Database, "elite_outbox_test_") || len(strings.TrimPrefix(cfg.ConnConfig.Database, "elite_outbox_test_")) < 16 {
		t.Fatal("processor requires dedicated loopback elite_outbox_test_<unique> database")
	}
	for _, fallback := range cfg.ConnConfig.Fallbacks {
		if fallback.Host != "127.0.0.1" {
			t.Fatal("non-loopback fallback forbidden")
		}
	}
	if other := os.Getenv("OUTBOX_TEST_DATABASE_URL"); other != "" {
		otherCfg, e := pgxpool.ParseConfig(other)
		if e == nil && otherCfg.ConnConfig.Host == cfg.ConnConfig.Host && otherCfg.ConnConfig.Port == cfg.ConnConfig.Port && otherCfg.ConnConfig.Database == cfg.ConnConfig.Database {
			t.Fatal("processor and store tests require distinct databases")
		}
	}
	cfg.MaxConns = 4
	cfg.ConnConfig.ConnectTimeout = 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	const tenant = "018f4d4a-7b36-7a21-8d10-2f4c54c28506"
	const first = "018f4d4a-7b36-7a21-8d10-2f4c54c28507"
	const retry = "018f4d4a-7b36-7a21-8d10-2f4c54c28508"
	cleanup := func() {
		clean, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if _, e := pool.Exec(clean, `delete from platform.outbox_event where tenant_id=$1`, tenant); e != nil {
			t.Error(e)
		}
		if _, e := pool.Exec(clean, `delete from platform.tenant where tenant_id=$1`, tenant); e != nil {
			t.Error(e)
		}
	}
	cleanup()
	defer cleanup()
	if _, err = pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,'processor-v316','Synthetic','Synthetic')`, tenant); err != nil {
		t.Fatal(err)
	}
	if _, err = pool.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload) values($1,$2,'test','a',1,'test.created',1,clock_timestamp(),'{}'),($1,$3,'test','b',1,'test.created',1,clock_timestamp(),'{}')`, tenant, first, retry); err != nil {
		t.Fatal(err)
	}
	store := db.NewOutbox(pool)
	seen := map[string][]int{}
	publisher := publisherFunc(func(publishCtx context.Context, event db.OutboxEvent, key string) error {
		deadline, ok := publishCtx.Deadline()
		if !ok || time.Until(deadline) > time.Minute || publishCtx.Err() != nil {
			t.Error("publisher lacks current bounded context")
		}
		if key != event.EventID || event.TenantID != tenant {
			t.Error("publisher event binding changed")
		}
		seen[key] = append(seen[key], event.Attempts)
		if event.EventID == retry && event.Attempts == 1 {
			return errors.New("synthetic transport failure")
		}
		return nil
	})
	processor := OutboxProcessor{Store: store, Publisher: publisher, WorkerID: "processor-v316", Lease: time.Minute, RetryDelay: 0, BatchSize: 2}
	if n, err := processor.ProcessOnce(ctx); err != nil || n != 1 {
		t.Fatalf("first batch n=%d err=%v", n, err)
	}
	var published bool
	var code string
	var attempts int
	if err = pool.QueryRow(ctx, `select published_at is not null,coalesce(last_error_code,''),attempts from platform.outbox_event where tenant_id=$1 and event_id=$2`, tenant, retry).Scan(&published, &code, &attempts); err != nil || published || code != "PUBLISH_FAILED" || attempts != 1 {
		t.Fatal("retry state not durable", published, code, attempts, err)
	}
	if n, err := processor.ProcessOnce(ctx); err != nil || n != 1 {
		t.Fatalf("retry batch n=%d err=%v", n, err)
	}
	if len(seen[first]) != 1 || len(seen[retry]) != 2 || seen[retry][0] != 1 || seen[retry][1] != 2 {
		t.Fatal("retry key/generation changed", seen)
	}
	var count int
	if err = pool.QueryRow(ctx, `select count(*) from platform.outbox_event where tenant_id=$1 and published_at is not null and claimed_by is null and claimed_until is null and last_error_code is null`, tenant).Scan(&count); err != nil || count != 2 {
		t.Fatal("ack state incomplete", count, err)
	}
	t.Log("OUTBOX_PROCESSOR_POSTGRES_PASS published=2 stable_retry_key=1 generations=1,2 provider=synthetic")
}

type fakeJobs struct {
	jobs              []db.Job
	completed, failed int
	remaining         time.Duration
	leaseErr          error
	generation        int
}

func (f *fakeJobs) Claim(context.Context, string, string, time.Duration, int) ([]db.Job, error) {
	return f.jobs, nil
}
func (f *fakeJobs) RemainingLease(_ context.Context, _, _, _ string, attempt int) (time.Duration, error) {
	if attempt != f.generation {
		return 0, errors.New("wrong preflight generation")
	}
	return f.remaining, f.leaseErr
}
func (f *fakeJobs) Complete(_ context.Context, _, _, _ string, attempt int) error {
	if attempt != f.generation {
		return errors.New("wrong completion generation")
	}
	f.completed++
	return nil
}
func (f *fakeJobs) Fail(_ context.Context, _, _, _, _ string, attempt int, _ time.Duration) (bool, error) {
	if attempt != f.generation {
		return false, errors.New("wrong failure generation")
	}
	f.failed++
	return false, nil
}

type fakeHandler struct{ fail string }

func (f fakeHandler) Handle(_ context.Context, j db.Job) error {
	if j.JobID == f.fail {
		return errors.New("bad")
	}
	return nil
}
func TestProcessorsSeparateSuccessAndRetry(t *testing.T) {
	ctx := context.Background()
	out := &fakeOutbox{events: []db.OutboxEvent{{TenantID: "t", EventID: "ok", Attempts: 1}, {TenantID: "t", EventID: "retry", Attempts: 1}}, remaining: time.Minute}
	n, err := (OutboxProcessor{Store: out, Publisher: fakePublisher{fail: "retry"}, WorkerID: "w", Lease: time.Minute, RetryDelay: time.Second, BatchSize: 2}).ProcessOnce(ctx)
	if err != nil || n != 1 || out.marked != 1 || out.released != 1 {
		t.Fatalf("outbox n=%d marked=%d released=%d err=%v", n, out.marked, out.released, err)
	}
	jobs := &fakeJobs{remaining: time.Minute, generation: 7, jobs: []db.Job{{TenantID: "t", JobID: "ok", Attempts: 7}, {TenantID: "t", JobID: "retry", Attempts: 7}}}
	n, err = (JobProcessor{Store: jobs, Handler: fakeHandler{fail: "retry"}, Queue: "q", WorkerID: "w", Lease: time.Minute, RetryDelay: time.Second, BatchSize: 2}).ProcessOnce(ctx)
	if err != nil || n != 1 || jobs.completed != 1 || jobs.failed != 1 {
		t.Fatalf("jobs n=%d completed=%d failed=%d err=%v", n, jobs.completed, jobs.failed, err)
	}
}

type handlerFunc func(context.Context, db.Job) error

func (f handlerFunc) Handle(ctx context.Context, job db.Job) error { return f(ctx, job) }

func TestJobProcessorLeaseBudgetAndLoss(t *testing.T) {
	for _, mode := range []string{"lost", "expired", "later-batch-expired", "bounded", "cancelled"} {
		t.Run(mode, func(t *testing.T) {
			store := &fakeJobs{generation: 3, remaining: time.Minute, jobs: []db.Job{{TenantID: "t", JobID: "one", Attempts: 3}}}
			if mode == "lost" {
				store.leaseErr = db.ErrJobClaimLost
			}
			if mode == "expired" {
				store.remaining = 0
			}
			if mode == "later-batch-expired" {
				store.jobs = append(store.jobs, db.Job{TenantID: "t", JobID: "two", Attempts: 3})
			}
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			if mode == "cancelled" {
				cancel()
			}
			calls := 0
			handler := handlerFunc(func(ctx context.Context, job db.Job) error {
				calls++
				deadline, ok := ctx.Deadline()
				if !ok || time.Until(deadline) <= 0 || time.Until(deadline) > time.Minute {
					t.Fatal("missing bounded lease deadline")
				}
				if mode == "later-batch-expired" {
					store.remaining = 0
				}
				return nil
			})
			n, err := (JobProcessor{Store: store, Handler: handler, Queue: "q", WorkerID: "w", Lease: time.Minute, BatchSize: 2}).ProcessOnce(ctx)
			if mode == "bounded" {
				if err != nil || n != 1 || calls != 1 || store.completed != 1 {
					t.Fatalf("bounded: %d %d %v", n, calls, err)
				}
			} else if mode == "later-batch-expired" {
				if err == nil || n != 1 || calls != 1 || store.completed != 1 {
					t.Fatalf("later batch: %d %d %v", n, calls, err)
				}
			} else if err == nil || n != 0 || calls != 0 || store.completed != 0 || store.failed != 0 {
				t.Fatalf("stale handler/effect: %d %d %v", n, calls, err)
			}
		})
	}
}
````

## 6. Configuration surface

Worker ID, queue/topic, batch size, lease, retry delay, attempt budget and poll/backpressure policy are typed constructor inputs. Defaults must be bounded; nonpositive durations/counts fail validation. Transport credentials remain external secrets owned by the selected adapter.

## 7. Dependency bill

| Dependency | Pin | Use | License | Scope | Source |
|---|---|---|---|---|---|
| Go | `1.26.7` | processors/tests | BSD-3-Clause | build/runtime | `go.dev` |
| PostgreSQL | `18.6` verified baseline | inbox/outbox/jobs/leases | PostgreSQL | runtime/test | `postgresql.org` |
| `pgx` | `5.10.0` | PostgreSQL adapter | MIT | build/runtime | `github.com/jackc/pgx` |

## 8. Apply order

Compose after PostgreSQL foundation, then materialize stores/processors and wire a selected transport/handler. Run duplicate, lease expiry, partial batch, retry and crash recovery tests. Existing queues need stable idempotency mapping. Rollback stops claim loops, drains, redeploys the prior compatible worker and preserves pending rows.

## 9. Verification

La composición limpia con Go core y PostgreSQL foundation pasó `gofmt` sin diff, `go test ./...`, `go vet ./...`, build y tests de integración en PostgreSQL 18.6; evidencia en `reconstruction_evidence/GO_RELIABLE_ASYNC_WORKERS_2026-08-24_V1.md`. Antes de producción faltan transport adapter real, idempotencia comprobada en el proveedor, shutdown/drain, métricas, backpressure, jitter/budgets, dead-letter tooling y carga/crash tests.

## 10. Reconstruction evidence

Clean reconstruction and PostgreSQL worker gates are recorded in `reconstruction_evidence/GO_RELIABLE_ASYNC_WORKERS_2026-08-24_V1.md`; final evidence validates version 0.1.1 in composition.
