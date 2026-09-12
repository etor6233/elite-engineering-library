# Go Data Analytics Core

## 1. Metadata

```yaml
pack_id: "GO-DATA-ANALYTICS-CORE"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa el contrato de ingesta (landing inmutable con replay idempotente, late-data bound y ownership de fuente) y de métricas semánticas (definiciones versionadas, freshness y reconciliación), tenant-scoped."
stacks: ["Go 1.26.7", "PostgreSQL 18.6"]
compatible_with: ["GO-ENTERPRISE-BACKEND 0.4.x", "PG-TX-FOUNDATION 0.1.x"]
incompatible_with: ["métricas sin tenant", "definición de métrica re-versionada", "replay sin idempotency key"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://go.dev", "https://www.postgresql.org/docs/18/"]
verified_at: "2026-09-02"
```

Este pack cubre las superficies `DATA-INGEST` y `ANALYTICS-BI` a nivel de contrato ejecutable. No afirma un warehouse/lakehouse ni un pipeline de CDC productivo: la ingesta define el landing y sus invariantes, y las métricas definen versionado, freshness y reconciliación. El motor de agregación real sobre el source of truth y el lakehouse son decisión del proyecto.

## 2. Applicability

Use este pack cuando la franquicia necesite ingestar registros de fuentes (POS/ERP/marketplaces) con replay idempotente y late-data bound, y materializar métricas semánticas (ventas, devoluciones, citas) con freshness y reconciliación, todo aislado por tenant.

Rechace este pack para: un warehouse analítico masivo sin admitir la plataforma; métricas sin un tenant obligatorio; o fuentes cuyo contrato/schema no esté versionado.

## 3. Architecture contract

- **Ownership**: `internal/analytics` gobierna la validación de métrica, freshness, reconciliación y registro de ingesta. La migración `0041` gobierna `data.*` (landing) y `analytics.*` (métricas). El source of truth de cada métrica sigue en su owner de dominio.
- **Invariantes**: (1) `tenant_id` obligatorio. (2) Definición de métrica versionada e inmutable (code+version únicos). (3) Replay idempotente por (source, external_id, payload sha). (4) `received_at >= event_time`; late-data bound por `MaxLateness`. (5) Una métrica no es fresca si excede su `FreshnessTTL`; la reconciliación no pasa con tolerancia negativa.
- **Data flow**: ingesta `IngestRecord.Validate` → landing (replay dedup) → validación; métrica `MetricDefinition` → agregación tenant-scoped → `metric_snapshot` con `reconciled` y `computed_at`.
- **Failure modes**: registro inválido/late-data → `ErrInvalidRecord`; métrica inválida → `ErrInvalidMetric`; snapshot stale → no se considera válido; reconciliación fallida → no se marca `reconciled`.
- **Seguridad/privacidad**: aislamiento por tenant en todo; payload aterriza como `jsonb` con sha de contenido; sin PII adicional.
- **Performance budget**: core Go sin I/O; agregaciones delegan en índices de PostgreSQL.
- **Operación/migración/rollback**: migración `0041` up/down atómica; el landing y las métricas se reconstruyen desde el source of truth.

## 4. Exact file manifest

```text
CREATE internal/analytics/metric.go
CREATE internal/analytics/snapshot.go
CREATE internal/analytics/ingest.go
CREATE internal/analytics/analytics_test.go
CREATE db/migrations/0041_data_analytics.up.sql
CREATE db/migrations/0041_data_analytics.down.sql
CREATE db/tests/0041_data_analytics.test.sql
```

## 5. Materialization blocks

### FILE: `internal/analytics/metric.go`
```yaml
block_id: "GO-DATA-ANALYTICS-CORE:internal/analytics/metric.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "9fd7359576bc1f14e1a91aacb90ea9fba8ea3e74aae1e9e3a22066a40a900637"
variables: []
secrets_allowed: false
```
````go
// Package analytics provides a governed semantic-metrics and data-ingest
// contract: versioned metric definitions, freshness and reconciliation gates,
// and a validated ingest record with replay idempotency and a late-data bound.
package analytics

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// Aggregation is the bounded set of metric aggregations.
type Aggregation string

const (
	AggCount Aggregation = "count"
	AggSum   Aggregation = "sum"
	AggMin   Aggregation = "min"
	AggMax   Aggregation = "max"
	AggAvg   Aggregation = "avg"
)

var (
	codeRe  = regexp.MustCompile(`^[a-z][a-z0-9._-]{0,63}$`)
	dimRe   = regexp.MustCompile(`^[a-z][a-z0-9._-]{0,63}$`)
	hex64Re = regexp.MustCompile(`^[0-9a-f]{64}$`)
)

var (
	ErrInvalidMetric = errors.New("analytics: invalid metric")
	ErrInvalidRecord = errors.New("analytics: invalid record")
)

// MetricDefinition is a versioned, tenant-scoped semantic metric. A version is
// immutable: the same code+version cannot be redefined.
type MetricDefinition struct {
	TenantID     string
	Code         string
	Version      int64
	Aggregation  Aggregation
	Source       string
	Dimensions   []string
	FreshnessTTL time.Duration
}

// Validate enforces the metric contract.
func (m MetricDefinition) Validate() error {
	if strings.TrimSpace(m.TenantID) == "" {
		return fmt.Errorf("%w: tenant", ErrInvalidMetric)
	}
	if !codeRe.MatchString(m.Code) {
		return fmt.Errorf("%w: code", ErrInvalidMetric)
	}
	if m.Version <= 0 {
		return fmt.Errorf("%w: version", ErrInvalidMetric)
	}
	switch m.Aggregation {
	case AggCount, AggSum, AggMin, AggMax, AggAvg:
	default:
		return fmt.Errorf("%w: aggregation", ErrInvalidMetric)
	}
	if strings.TrimSpace(m.Source) == "" {
		return fmt.Errorf("%w: source", ErrInvalidMetric)
	}
	if len(m.Dimensions) > 16 {
		return fmt.Errorf("%w: dimensions", ErrInvalidMetric)
	}
	seen := map[string]bool{}
	for _, d := range m.Dimensions {
		if !dimRe.MatchString(d) {
			return fmt.Errorf("%w: dimension %q", ErrInvalidMetric, d)
		}
		if seen[d] {
			return fmt.Errorf("%w: duplicate dimension %q", ErrInvalidMetric, d)
		}
		seen[d] = true
	}
	if m.FreshnessTTL <= 0 {
		return fmt.Errorf("%w: freshness ttl", ErrInvalidMetric)
	}
	return nil
}
````

### FILE: `internal/analytics/snapshot.go`
```yaml
block_id: "GO-DATA-ANALYTICS-CORE:internal/analytics/snapshot.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "59e1d536351cccb279961a5ff6ab099dc9a64f36575d6f123f30a0e73f6673a8"
variables: []
secrets_allowed: false
```
````go
package analytics

import "time"

// Fresh reports whether a value computed at computedAt is still within its
// freshness TTL as of now. A zero/negative TTL is never fresh (fail-closed).
func Fresh(computedAt, now time.Time, ttl time.Duration) bool {
	if ttl <= 0 {
		return false
	}
	return now.Sub(computedAt) <= ttl
}

// Reconciles reports whether the sum of parts equals the total within a
// non-negative absolute tolerance. Reconciliation never passes on a negative
// tolerance.
func Reconciles(parts []float64, total, tolerance float64) bool {
	if tolerance < 0 {
		return false
	}
	var sum float64
	for _, p := range parts {
		sum += p
	}
	d := sum - total
	if d < 0 {
		d = -d
	}
	return d <= tolerance
}
````

### FILE: `internal/analytics/ingest.go`
```yaml
block_id: "GO-DATA-ANALYTICS-CORE:internal/analytics/ingest.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "5623eec2192192bada127e1492f5e7bac6d75855fd8cd039f88f5f202a4ec8a7"
variables: []
secrets_allowed: false
```
````go
package analytics

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

var sourceRe = regexp.MustCompile(`^[a-z][a-z0-9._-]{0,63}$`)

// IngestRecord is a validated raw intake: it carries source ownership, an
// idempotency replay key and a bounded late-data window.
type IngestRecord struct {
	SourceID      string
	ExternalID    string
	EventTime     time.Time
	ReceivedAt    time.Time
	PayloadSHA256 string
	MaxLateness   time.Duration // 0 means unbounded (no late-data rejection)
}

// Validate enforces the ingest contract: source identity, external id, payload
// digest, event/received ordering and the late-data bound.
func (r IngestRecord) Validate() error {
	if !sourceRe.MatchString(r.SourceID) {
		return fmt.Errorf("%w: source", ErrInvalidRecord)
	}
	if strings.TrimSpace(r.ExternalID) == "" || len(r.ExternalID) > 200 {
		return fmt.Errorf("%w: external id", ErrInvalidRecord)
	}
	if !hex64Re.MatchString(strings.ToLower(r.PayloadSHA256)) {
		return fmt.Errorf("%w: payload sha", ErrInvalidRecord)
	}
	if r.EventTime.IsZero() || r.ReceivedAt.IsZero() {
		return fmt.Errorf("%w: timestamps", ErrInvalidRecord)
	}
	if r.ReceivedAt.Before(r.EventTime) {
		return fmt.Errorf("%w: received before event", ErrInvalidRecord)
	}
	if r.MaxLateness > 0 && r.ReceivedAt.Sub(r.EventTime) > r.MaxLateness {
		return fmt.Errorf("%w: exceeds lateness bound", ErrInvalidRecord)
	}
	return nil
}

// ReplayKey returns the idempotency key for deduplicated replay.
func (r IngestRecord) ReplayKey() string {
	return r.SourceID + "\x00" + r.ExternalID + "\x00" + strings.ToLower(r.PayloadSHA256)
}
````

### FILE: `internal/analytics/analytics_test.go`
```yaml
block_id: "GO-DATA-ANALYTICS-CORE:internal/analytics/analytics_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "1682eea431dc85c0309305cb00df809647c12344fc1afa882da8936288809ba3"
variables: []
secrets_allowed: false
```
````go
package analytics

import (
	"errors"
	"testing"
	"time"
)

func validMetric() MetricDefinition {
	return MetricDefinition{
		TenantID: "tenant-a", Code: "sales.revenue", Version: 1,
		Aggregation: AggSum, Source: "sales.customer_order",
		Dimensions: []string{"organization", "month"}, FreshnessTTL: time.Hour,
	}
}

func TestMetricValidate(t *testing.T) {
	if err := validMetric().Validate(); err != nil {
		t.Fatalf("valid metric rejected: %v", err)
	}
	bad := validMetric()
	bad.Aggregation = "median"
	if err := bad.Validate(); !errors.Is(err, ErrInvalidMetric) {
		t.Fatalf("invalid aggregation accepted: %v", err)
	}
	dup := validMetric()
	dup.Dimensions = []string{"a", "a"}
	if err := dup.Validate(); !errors.Is(err, ErrInvalidMetric) {
		t.Fatalf("duplicate dimension accepted: %v", err)
	}
	zero := validMetric()
	zero.FreshnessTTL = 0
	if err := zero.Validate(); !errors.Is(err, ErrInvalidMetric) {
		t.Fatalf("zero freshness accepted: %v", err)
	}
}

func TestFresh(t *testing.T) {
	now := time.Now()
	if !Fresh(now.Add(-30*time.Minute), now, time.Hour) {
		t.Fatal("fresh value reported stale")
	}
	if Fresh(now.Add(-2*time.Hour), now, time.Hour) {
		t.Fatal("stale value reported fresh")
	}
	if Fresh(now, now, 0) {
		t.Fatal("zero ttl reported fresh (should fail closed)")
	}
}

func TestReconciles(t *testing.T) {
	if !Reconciles([]float64{1, 2, 3}, 6, 0) {
		t.Fatal("exact reconcile failed")
	}
	if !Reconciles([]float64{1, 2, 3}, 6.05, 0.1) {
		t.Fatal("tolerated reconcile failed")
	}
	if Reconciles([]float64{1, 2}, 6, 0.1) {
		t.Fatal("mismatched reconcile passed")
	}
	if Reconciles([]float64{1}, 1, -0.1) {
		t.Fatal("negative tolerance passed")
	}
}

func validRecord() IngestRecord {
	return IngestRecord{
		SourceID: "pos", ExternalID: "rec-1",
		EventTime: time.Now().Add(-time.Minute), ReceivedAt: time.Now(),
		PayloadSHA256: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		MaxLateness: time.Hour,
	}
}

func TestIngestValidate(t *testing.T) {
	if err := validRecord().Validate(); err != nil {
		t.Fatalf("valid record rejected: %v", err)
	}
	late := validRecord()
	late.EventTime = time.Now().Add(-3 * time.Hour) // exceeds MaxLateness (1h)
	if err := late.Validate(); !errors.Is(err, ErrInvalidRecord) {
		t.Fatalf("late record accepted: %v", err)
	}
	future := validRecord()
	future.ReceivedAt = future.EventTime.Add(-time.Minute) // received before event
	if err := future.Validate(); !errors.Is(err, ErrInvalidRecord) {
		t.Fatalf("received-before-event accepted: %v", err)
	}
	badSHA := validRecord()
	badSHA.PayloadSHA256 = "zz"
	if err := badSHA.Validate(); !errors.Is(err, ErrInvalidRecord) {
		t.Fatalf("bad sha accepted: %v", err)
	}
}

func TestIngestReplayKey(t *testing.T) {
	a := validRecord()
	b := validRecord()
	if a.ReplayKey() != b.ReplayKey() {
		t.Fatal("identical records produced different replay keys")
	}
	b.ExternalID = "rec-2"
	if a.ReplayKey() == b.ReplayKey() {
		t.Fatal("different records produced same replay key")
	}
}
````

### FILE: `db/migrations/0041_data_analytics.up.sql`
```yaml
block_id: "GO-DATA-ANALYTICS-CORE:db/migrations/0041_data_analytics.up.sql:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "bf98c3f6cceda9883c15c0e72615dc124d0f13e3a8401052a84a3d395e587f9e"
variables: []
secrets_allowed: false
```
````sql
begin;

create schema if not exists data;
create schema if not exists analytics;

-- ingest source ownership and contract
create table data.ingest_source (
  tenant_id uuid not null,
  source_id text not null,
  owner_organization_id text,
  contract_schema text not null,
  freshness_ttl interval not null,
  created_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, source_id),
  foreign key (tenant_id) references platform.tenant (tenant_id),
  check (source_id ~ '^[a-z][a-z0-9._-]{0,63}$'),
  check (length(contract_schema) between 1 and 200),
  check (freshness_ttl > interval '0')
);

-- immutable landing record with replay idempotency and event/received ordering
create table data.landing_record (
  tenant_id uuid not null,
  source_id text not null,
  external_id text not null,
  event_time timestamptz not null,
  received_at timestamptz not null default clock_timestamp(),
  payload jsonb not null,
  source_sha256 text not null,
  lineage text not null default '',
  state text not null default 'landed' check (state in ('landed', 'validated', 'rejected')),
  primary key (tenant_id, source_id, external_id, source_sha256),
  foreign key (tenant_id, source_id) references data.ingest_source (tenant_id, source_id),
  check (received_at >= event_time),
  check (jsonb_typeof(payload) = 'object'),
  check (source_sha256 ~ '^[0-9a-f]{64}$'),
  check (length(external_id) between 1 and 200)
);

create index landing_record_source_time_idx
  on data.landing_record (tenant_id, source_id, received_at);

-- versioned, immutable semantic metric definition
create table analytics.metric_definition (
  tenant_id uuid not null,
  metric_code text not null,
  version bigint not null check (version > 0),
  aggregation text not null check (aggregation in ('count', 'sum', 'min', 'max', 'avg')),
  source_ref text not null,
  dimensions jsonb not null default '[]'::jsonb,
  freshness_ttl interval not null,
  created_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, metric_code, version),
  foreign key (tenant_id) references platform.tenant (tenant_id),
  check (metric_code ~ '^[a-z][a-z0-9._-]{0,63}$'),
  check (length(source_ref) between 1 and 200),
  check (jsonb_typeof(dimensions) = 'array'),
  check (freshness_ttl > interval '0')
);

-- materialized metric snapshot
create table analytics.metric_snapshot (
  tenant_id uuid not null,
  metric_code text not null,
  version bigint not null,
  period_start timestamptz not null,
  period_end timestamptz not null,
  value numeric not null,
  dimensions jsonb not null default '{}'::jsonb,
  computed_at timestamptz not null default clock_timestamp(),
  reconciled boolean not null default false,
  primary key (tenant_id, metric_code, version, period_start, period_end),
  foreign key (tenant_id, metric_code, version)
    references analytics.metric_definition (tenant_id, metric_code, version),
  check (period_end > period_start),
  check (jsonb_typeof(dimensions) = 'object')
);

commit;
````

### FILE: `db/migrations/0041_data_analytics.down.sql`
```yaml
block_id: "GO-DATA-ANALYTICS-CORE:db/migrations/0041_data_analytics.down.sql:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "427c5b8c10fec82410d31f1b335188dbb47fa9e20c3f284ef4ec648b4ca379aa"
variables: []
secrets_allowed: false
```
````sql
begin;

drop table if exists analytics.metric_snapshot;
drop table if exists analytics.metric_definition;
drop table if exists data.landing_record;
drop table if exists data.ingest_source;
drop schema if exists analytics;
drop schema if exists data;

commit;
````

### FILE: `db/tests/0041_data_analytics.test.sql`
```yaml
block_id: "GO-DATA-ANALYTICS-CORE:db/tests/0041_data_analytics.test.sql:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "2e52a4823898877095946c4467e646c2975ba66499e5595ed4da2f2e8dc36191"
variables: []
secrets_allowed: false
```
````sql
-- 0041_data_analytics.test.sql — verifica replay idempotente, rechazo de
-- received-before-event, definición versionada de métrica y snapshot.
-- Precondición: migración 0001 (platform.tenant) y 0041 aplicadas.

begin;

insert into platform.tenant (tenant_id, tenant_code, legal_name, display_name) values
  ('11111111-1111-1111-1111-111111111111', 'tenant-a', 'A', 'Tenant A');

insert into data.ingest_source (tenant_id, source_id, contract_schema, freshness_ttl) values
  ('11111111-1111-1111-1111-111111111111', 'pos', 'pos-events/v1', interval '1 hour');

insert into data.landing_record
  (tenant_id, source_id, external_id, event_time, received_at, payload, source_sha256)
values
  ('11111111-1111-1111-1111-111111111111', 'pos', 'rec-1', now() - interval '1 minute', now(),
   '{"total":10}', repeat('a',64));

insert into data.landing_record
  (tenant_id, source_id, external_id, event_time, received_at, payload, source_sha256)
values
  ('11111111-1111-1111-1111-111111111111', 'pos', 'rec-1', now() - interval '1 minute', now(),
   '{"total":10}', repeat('a',64))
on conflict do nothing;

-- replay: exactly one row despite the duplicate insert
do $$
declare n int;
begin
  select count(*) into n from data.landing_record
   where tenant_id = '11111111-1111-1111-1111-111111111111' and external_id = 'rec-1';
  if n <> 1 then raise exception 'replay dedup expected 1, got %', n; end if;
end $$;

-- ordering: received_at < event_time must violate the check constraint
do $$
begin
  begin
    insert into data.landing_record
      (tenant_id, source_id, external_id, event_time, received_at, payload, source_sha256)
    values
      ('11111111-1111-1111-1111-111111111111', 'pos', 'rec-2', now(), now() - interval '1 minute',
       '{"x":1}', repeat('b',64));
    raise exception 'received-before-event was accepted';
  exception when check_violation then
    null; -- expected
  end;
end $$;

-- versioned metric definition + snapshot
insert into analytics.metric_definition
  (tenant_id, metric_code, version, aggregation, source_ref, dimensions, freshness_ttl)
values
  ('11111111-1111-1111-1111-111111111111', 'sales.revenue', 1, 'sum', 'sales.customer_order',
   '["organization","month"]', interval '1 hour');

insert into analytics.metric_snapshot
  (tenant_id, metric_code, version, period_start, period_end, value, dimensions, reconciled)
values
  ('11111111-1111-1111-1111-111111111111', 'sales.revenue', 1,
   '2026-09-01 00:00:00+00', '2026-09-02 00:00:00+00', 1234.56,
   '{"organization":"org-1"}', true);

do $$
declare n int;
begin
  select count(*) into n from analytics.metric_snapshot
   where tenant_id = '11111111-1111-1111-1111-111111111111';
  if n <> 1 then raise exception 'snapshot expected 1, got %', n; end if;
end $$;

rollback;
````


## 6. Configuration surface

Sin variables ni secretos. `MaxLateness` y `FreshnessTTL` son parámetros de cada registro/métrica y se validan (positivos o acotados).

## 7. Dependency bill

| Package/image/tool | Pin exacto | Uso | Licencia | Runtime/build | Fuente oficial |
|---|---|---|---|---|---|
| Go standard library | go 1.26.7 (toolchain) | core Go | BSD-3-Clause | runtime | https://go.dev |
| PostgreSQL 18.6 | 18.6 | landing/métricas/índices | PostgreSQL License | runtime | https://www.postgresql.org |

## 8. Apply order

1. Componer `PG-TX-FOUNDATION` (migración `0001`, `platform.tenant`) y el backend.
2. Aplicar migración `0041` sobre PostgreSQL 18.6.
3. Colocar los cuatro archivos Go bajo `internal/analytics/` y los tres SQL bajo `db/`.
4. Verificar: `psql ... -v ON_ERROR_STOP=1 -f db/tests/0041_data_analytics.test.sql` y `go test ./... -count=1`.
5. Rollback: migración `0041` down y eliminar `internal/analytics/`.

## 9. Verification

- `go test ./internal/analytics/ -count=1`: 5/5 PASS (métrica con negativos de agregación/dimensión/TTL, freshness fail-closed, reconciliación con tolerancia, ingesta con late-data/recibido-antes-evento/sha y replay key).
- `go test ./... -count=1` (5 paquetes): PASS.
- `go vet ./...`: exit 0.
- PostgreSQL 18.6 real (initdb → up → test → down): `ON_ERROR_STOP=1` exit 0; aserciones DO: replay dedup `INSERT 0 0`, rechazo received-before-event por `check_violation`, métrica versionada + snapshot; down deja `data_ns=true analytics_ns=true`.

## 10. Reconstruction evidence

Véase `reconstruction_evidence/GO_DATA_ANALYTICS_CORE_2026-09-02_V179.md`.
