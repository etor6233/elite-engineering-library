# Go Return Effect Execution Worker

## 1. Metadata

```yaml
pack_id: "GO-RETURN-EFFECT-EXECUTION-WORKER"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Añade proyección durable, claim/lease/fencing, historial inmutable, reanudación auditada y un worker que aplica atómicamente la disposición física de inventario solicitada por el retorno; no simula refunds, exchange, accounting ni fiscal."
stacks: ["Go 1.26.7", "PostgreSQL 18.6", "pgx 5.10.0"]
compatible_with: ["GO-ENTERPRISE-BACKEND 0.4.x", "GO-RELIABLE-ASYNC-WORKERS 0.1.x", "GO-SUPPLY-FACTORY-INVENTORY-API 0.12.x", "GO-FRANCHISE-CUSTOMER-JOURNEY-API 0.8.x", "GO-ELECTROMOBILITY-APPLICATION 1.9.x"]
incompatible_with: ["alternate return-effect owner or mutable inventory projection without reconciliation contract"]
license_expression: "LicenseRef-Workspace-Owner AND MIT dependency"
upstream_sources: ["https://docs.aws.amazon.com/wellarchitected/latest/reliability-pillar/welcome.html", "https://www.postgresql.org/docs/18/explicit-locking.html", "https://www.postgresql.org/docs/18/sql-select.html", "https://learn.microsoft.com/en-us/dynamics365/business-central/sales-how-process-sales-returns-cancellations", "https://docs.stripe.com/api/idempotent_requests", "https://docs.stripe.com/refunds", "https://www.arca.gob.ar/fe/ayuda/webservice.asp"]
verified_at: "2026-08-30"
```

All eight blocks are local `AUTHORED` implementation. AWS governs persisted idempotency and retry expectations; PostgreSQL governs locking and `SKIP LOCKED`; Microsoft governs linkage to the original sale and separated return effects; Stripe governs provider idempotency and asynchronous refund reconciliation; ARCA governs fiscal homologation. No local line is attributed to those companies.

## 2. Applicability

Use after the V140 return receipt/disposition pack when a project needs durable execution instead of leaving owner requests as inert records. It adds one execution projection and audit ledger without duplicating inventory, payment, fulfillment, accounting or fiscal ownership. Reject it when an external ERP owns inventory unless a synchronized adapter and reconciliation contract replace the local inventory consumer.

## 3. Architecture contract

Every `sales.return_effect_request` receives one mutable execution projection; request, receipt, disposition and attempt/resume evidence remain append-only. Claims use `FOR UPDATE SKIP LOCKED`, a bounded lease, attempt number and fencing token. Expired claims may be reclaimed; an old worker cannot complete a new claim. Retry, blocked, failed and succeeded are explicit. A blocked request returns to the queue only through a durable operator resume record.

The included worker selects only `owner_context=inventory`. It locks the claimed execution and exact stock unit in a serializable transaction, requires the physical stock to remain `sold`, maps `quarantine→quarantine`, `restock→available`, `repair→service`, and `scrap→retired`, increments the existing stock version, writes an outbox event, immutable success attempt and execution result atomically. The result digest is calculated server-side over the exact request/stock/from/to/version command. Any state, owner, token, lease, action or version mismatch fails closed.

Refund, exchange, accounting and fiscal requests remain unclaimed by this binary. Stripe requires a real provider result plus webhook/retrieval reconciliation; Microsoft requires exact source linkage and accounting separation; ARCA requires credentials, point of sale, approved document mappings and homologation. Their later consumers must use the same stable request idempotency key, persist provider references/status and never translate “requested” into “succeeded” without target evidence.

## 4. Exact file manifest

```text
CREATE db/migrations/0019_return_effect_execution.up.sql
CREATE db/migrations/0019_return_effect_execution.down.sql
CREATE db/tests/0019_return_effect_execution.test.sql
CREATE internal/returneffects/processor.go
CREATE internal/returneffects/processor_test.go
CREATE internal/platform/postgres/returneffects.go
CREATE internal/platform/postgres/returneffects_integration_test.go
CREATE cmd/return-effect-worker/main.go
```

## 5. Materialization blocks

### FILE: `db/migrations/0019_return_effect_execution.up.sql`

```yaml
block_id: "GO-RETURN-EFFECT:file:01"
operation: CREATE
provenance: AUTHORED
source: "local implementation governed by official sources in metadata"
license: "LicenseRef-Workspace-Owner"
sha256: "483a9b57e8d5f0e5bfd95f455f8820ed46f0a11e2e8cdd502bbce5d259432dfd"
variables: []
secrets_allowed: false
```

````sql
begin;

create table sales.return_effect_execution (
  tenant_id uuid not null,
  request_id text not null,
  status text not null check (status in ('requested','claimed','retry','blocked','failed','succeeded')),
  attempt_count integer not null default 0 check (attempt_count >= 0),
  available_at timestamptz not null default clock_timestamp(),
  claimed_at timestamptz,
  claimed_by text,
  claim_token text,
  claimed_until timestamptz,
  last_error_code text,
  provider_reference text,
  result_sha256_hex text,
  updated_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id,request_id),
  foreign key (tenant_id,request_id) references sales.return_effect_request(tenant_id,request_id),
  unique (tenant_id,claim_token),
  check ((status='claimed')=(claimed_at is not null and claimed_by is not null and claim_token is not null and claimed_until is not null)),
  check (status<>'claimed' or claimed_until>claimed_at),
  check (result_sha256_hex is null or result_sha256_hex ~ '^[0-9a-f]{64}$'),
  check (last_error_code is null or last_error_code ~ '^[A-Z][A-Z0-9_]{1,63}$'),
  check (status<>'succeeded' or (result_sha256_hex is not null and last_error_code is null)),
  check (status not in ('retry','blocked','failed') or last_error_code is not null)
);

create table sales.return_effect_attempt (
  tenant_id uuid not null,
  attempt_id text not null,
  request_id text not null,
  attempt_no integer not null check (attempt_no>0),
  worker_id text not null check (worker_id ~ '^[A-Za-z0-9][A-Za-z0-9._:-]{0,63}$'),
  outcome text not null check (outcome in ('retry','blocked','failed','succeeded')),
  error_code text,
  provider_reference text,
  result_sha256_hex text,
  started_at timestamptz not null,
  completed_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id,attempt_id),
  unique (tenant_id,request_id,attempt_no),
  foreign key (tenant_id,request_id) references sales.return_effect_request(tenant_id,request_id),
  check (completed_at>=started_at),
  check (error_code is null or error_code ~ '^[A-Z][A-Z0-9_]{1,63}$'),
  check (result_sha256_hex is null or result_sha256_hex ~ '^[0-9a-f]{64}$'),
  check ((outcome='succeeded' and error_code is null and result_sha256_hex is not null) or (outcome<>'succeeded' and error_code is not null and result_sha256_hex is null))
);

create table sales.return_effect_resume (
  tenant_id uuid not null,
  resume_id text not null,
  request_id text not null,
  reason_code text not null check (reason_code ~ '^[A-Z][A-Z0-9_]{1,63}$'),
  requested_by_subject text not null check (length(requested_by_subject) between 1 and 255),
  requested_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id,resume_id),
  foreign key (tenant_id,request_id) references sales.return_effect_request(tenant_id,request_id)
);

create or replace function sales.enqueue_return_effect_execution() returns trigger language plpgsql as $$
begin
  insert into sales.return_effect_execution(tenant_id,request_id,status)
  values(new.tenant_id,new.request_id,'requested');
  return new;
end $$;

create trigger return_effect_execution_enqueue
after insert on sales.return_effect_request
for each row execute function sales.enqueue_return_effect_execution();

insert into sales.return_effect_execution(tenant_id,request_id,status)
select tenant_id,request_id,'requested' from sales.return_effect_request
on conflict do nothing;

create or replace function sales.prevent_return_effect_audit_mutation() returns trigger language plpgsql as $$
begin
  raise exception using errcode='23514',message='return effect audit evidence is immutable';
end $$;

create trigger return_effect_attempt_immutable before update or delete on sales.return_effect_attempt for each row execute function sales.prevent_return_effect_audit_mutation();
create trigger return_effect_resume_immutable before update or delete on sales.return_effect_resume for each row execute function sales.prevent_return_effect_audit_mutation();

create index return_effect_execution_claim_idx on sales.return_effect_execution(status,available_at,claimed_until,tenant_id,request_id);
create index return_effect_attempt_request_idx on sales.return_effect_attempt(tenant_id,request_id,attempt_no desc);

commit;
````

### FILE: `db/migrations/0019_return_effect_execution.down.sql`

```yaml
block_id: "GO-RETURN-EFFECT:file:02"
operation: CREATE
provenance: AUTHORED
source: "local implementation governed by official sources in metadata"
license: "LicenseRef-Workspace-Owner"
sha256: "3ece816ebe15ccaf167eca4c6ab321d9f235903500ace999033d41d958e9b7a3"
variables: []
secrets_allowed: false
```

````sql
begin;

drop index if exists sales.return_effect_attempt_request_idx;
drop index if exists sales.return_effect_execution_claim_idx;
drop trigger if exists return_effect_resume_immutable on sales.return_effect_resume;
drop trigger if exists return_effect_attempt_immutable on sales.return_effect_attempt;
drop function if exists sales.prevent_return_effect_audit_mutation();
drop trigger if exists return_effect_execution_enqueue on sales.return_effect_request;
drop function if exists sales.enqueue_return_effect_execution();
drop table if exists sales.return_effect_resume;
drop table if exists sales.return_effect_attempt;
drop table if exists sales.return_effect_execution;

commit;
````

### FILE: `db/tests/0019_return_effect_execution.test.sql`

```yaml
block_id: "GO-RETURN-EFFECT:file:03"
operation: CREATE
provenance: AUTHORED
source: "local implementation governed by official sources in metadata"
license: "LicenseRef-Workspace-Owner"
sha256: "a59b744c1986f923ce19002c6a3b195ef79b09cd669e54f4ada1c6f1e25da762"
variables: []
secrets_allowed: false
```

````sql
begin;

do $$
begin
  if to_regclass('sales.return_effect_execution') is null or to_regclass('sales.return_effect_attempt') is null or to_regclass('sales.return_effect_resume') is null then
    raise exception 'return effect execution tables missing';
  end if;
  if (select count(*) from pg_trigger where tgname in ('return_effect_execution_enqueue','return_effect_attempt_immutable','return_effect_resume_immutable') and not tgisinternal) <> 3 then
    raise exception 'return effect execution trigger contract incomplete';
  end if;
  if not exists (select 1 from pg_indexes where schemaname='sales' and indexname='return_effect_execution_claim_idx') then
    raise exception 'return effect claim index missing';
  end if;
  if exists (select 1 from sales.return_effect_request r left join sales.return_effect_execution e on e.tenant_id=r.tenant_id and e.request_id=r.request_id where e.request_id is null) then
    raise exception 'return effect request missing execution projection';
  end if;
end $$;

rollback;
````

### FILE: `internal/returneffects/processor.go`

```yaml
block_id: "GO-RETURN-EFFECT:file:04"
operation: CREATE
provenance: AUTHORED
source: "local implementation governed by official sources in metadata"
license: "LicenseRef-Workspace-Owner"
sha256: "be09b8fbd847dcc149359538e57b38b1bceb44b05e7743c75ab3e0561f088214"
variables: []
secrets_allowed: false
```

````go
package returneffects

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"
)

var (
	ErrNoWork         = errors.New("no return effect work")
	ErrInventoryState = errors.New("return inventory state conflict")
)

var workerIDPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:-]{0,63}$`)

type Work struct {
	TenantID        string
	RequestID       string
	DispositionID   string
	EffectKind      string
	OwnerContext    string
	IdempotencyKey  string
	Attempt         int
	ClaimToken      string
	InventoryAction string
	OrganizationID  string
	StockUnitID     string
}

type Result struct {
	RequestID      string
	StockUnitID    string
	InventoryState string
	StockVersion   int64
	ResultSHA256   string
}

type Completion struct {
	Outcome           string
	ErrorCode         string
	ProviderReference string
	ResultSHA256      string
	RetryAfter        time.Duration
}

type Store interface {
	Claim(context.Context, string, string, string, time.Duration) (*Work, error)
	ApplyInventory(context.Context, Work, string) (Result, error)
	Finish(context.Context, Work, string, Completion) error
	ResumeBlocked(context.Context, string, string, string, string, string) error
}

type IDGenerator interface{ New() string }

type Processor struct {
	store      Store
	ids        IDGenerator
	workerID   string
	lease      time.Duration
	retryDelay time.Duration
}

func NewInventoryProcessor(store Store, ids IDGenerator, workerID string, lease, retryDelay time.Duration) (*Processor, error) {
	if store == nil || ids == nil || !workerIDPattern.MatchString(workerID) || lease < time.Second || lease > 15*time.Minute || retryDelay < time.Second || retryDelay > time.Hour {
		return nil, fmt.Errorf("invalid return effect processor configuration")
	}
	return &Processor{store: store, ids: ids, workerID: workerID, lease: lease, retryDelay: retryDelay}, nil
}

func (p *Processor) ProcessOne(ctx context.Context) (Result, error) {
	claimToken := p.ids.New()
	work, err := p.store.Claim(ctx, "inventory", p.workerID, claimToken, p.lease)
	if err != nil {
		return Result{}, err
	}
	if work == nil {
		return Result{}, ErrNoWork
	}
	result, applyErr := p.store.ApplyInventory(ctx, *work, p.workerID)
	if applyErr == nil {
		return result, nil
	}
	completion := Completion{Outcome: "retry", ErrorCode: "INVENTORY_TRANSIENT_FAILURE", RetryAfter: p.retryDelay}
	if errors.Is(applyErr, ErrInventoryState) {
		completion = Completion{Outcome: "blocked", ErrorCode: "INVENTORY_STATE_CONFLICT"}
	}
	if finishErr := p.store.Finish(ctx, *work, p.workerID, completion); finishErr != nil {
		return Result{}, errors.Join(applyErr, finishErr)
	}
	return Result{}, applyErr
}
````

### FILE: `internal/returneffects/processor_test.go`

```yaml
block_id: "GO-RETURN-EFFECT:file:05"
operation: CREATE
provenance: AUTHORED
source: "local implementation governed by official sources in metadata"
license: "LicenseRef-Workspace-Owner"
sha256: "701622e3ffd79125c420326f56c0df0401fc0adf94f22e10155596a261f09338"
variables: []
secrets_allowed: false
```

````go
package returneffects

import (
	"context"
	"errors"
	"testing"
	"time"
)

type fixedIDs struct{ value string }

func (f fixedIDs) New() string { return f.value }

type fakeStore struct {
	work       *Work
	claimErr   error
	applyErr   error
	result     Result
	completion Completion
	finishes   int
}

func (f *fakeStore) Claim(context.Context, string, string, string, time.Duration) (*Work, error) {
	return f.work, f.claimErr
}
func (f *fakeStore) ApplyInventory(context.Context, Work, string) (Result, error) {
	return f.result, f.applyErr
}
func (f *fakeStore) Finish(_ context.Context, _ Work, _ string, c Completion) error {
	f.finishes++
	f.completion = c
	return nil
}
func (f *fakeStore) ResumeBlocked(context.Context, string, string, string, string, string) error {
	return nil
}

func TestInventoryProcessorSuccessNoWorkAndFailClosed(t *testing.T) {
	if _, err := NewInventoryProcessor(nil, fixedIDs{"id"}, "worker", time.Minute, time.Second); err == nil {
		t.Fatal("nil store accepted")
	}
	noWork := &fakeStore{}
	processor, err := NewInventoryProcessor(noWork, fixedIDs{"claim"}, "worker-1", time.Minute, time.Second)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = processor.ProcessOne(context.Background()); !errors.Is(err, ErrNoWork) {
		t.Fatalf("no work err=%v", err)
	}

	work := &Work{RequestID: "request", ClaimToken: "claim"}
	success := &fakeStore{work: work, result: Result{RequestID: "request", InventoryState: "quarantine", StockVersion: 2}}
	processor, _ = NewInventoryProcessor(success, fixedIDs{"claim"}, "worker-1", time.Minute, time.Second)
	result, err := processor.ProcessOne(context.Background())
	if err != nil || result.InventoryState != "quarantine" || success.finishes != 0 {
		t.Fatalf("result=%+v finishes=%d err=%v", result, success.finishes, err)
	}

	blocked := &fakeStore{work: work, applyErr: ErrInventoryState}
	processor, _ = NewInventoryProcessor(blocked, fixedIDs{"claim"}, "worker-1", time.Minute, time.Second)
	if _, err = processor.ProcessOne(context.Background()); !errors.Is(err, ErrInventoryState) {
		t.Fatalf("blocked err=%v", err)
	}
	if blocked.finishes != 1 || blocked.completion.Outcome != "blocked" || blocked.completion.ErrorCode != "INVENTORY_STATE_CONFLICT" {
		t.Fatalf("completion=%+v", blocked.completion)
	}

	transient := &fakeStore{work: work, applyErr: errors.New("database unavailable")}
	processor, _ = NewInventoryProcessor(transient, fixedIDs{"claim"}, "worker-1", time.Minute, 3*time.Second)
	if _, err = processor.ProcessOne(context.Background()); err == nil {
		t.Fatal("transient failure hidden")
	}
	if transient.completion.Outcome != "retry" || transient.completion.RetryAfter != 3*time.Second {
		t.Fatalf("completion=%+v", transient.completion)
	}
}
````

### FILE: `internal/platform/postgres/returneffects.go`

```yaml
block_id: "GO-RETURN-EFFECT:file:06"
operation: CREATE
provenance: AUTHORED
source: "local implementation governed by official sources in metadata"
license: "LicenseRef-Workspace-Owner"
sha256: "3ccb86aae49c251307cae50aeedb0834a1cf698402b16f6ee8dcd28507005ca0"
variables: []
secrets_allowed: false
```

````go
package postgres

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"elite.local/enterprise/internal/returneffects"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReturnEffects struct{ pool *pgxpool.Pool }

func NewReturnEffects(pool *pgxpool.Pool) *ReturnEffects { return &ReturnEffects{pool: pool} }

func (r *ReturnEffects) Claim(ctx context.Context, owner, worker, claimToken string, lease time.Duration) (*returneffects.Work, error) {
	if owner == "" || worker == "" || claimToken == "" || lease <= 0 {
		return nil, fmt.Errorf("invalid return effect claim")
	}
	var work returneffects.Work
	err := r.pool.QueryRow(ctx, `
with candidate as (
  select x.tenant_id,x.request_id
    from sales.return_effect_execution x
    join sales.return_effect_request e on e.tenant_id=x.tenant_id and e.request_id=x.request_id
   where e.owner_context=$1
     and ((x.status in ('requested','retry') and x.available_at<=clock_timestamp()) or (x.status='claimed' and x.claimed_until<clock_timestamp()))
   order by x.available_at,e.requested_at,x.request_id
   for update of x skip locked
   limit 1
), claimed as (
  update sales.return_effect_execution x
     set status='claimed',attempt_count=x.attempt_count+1,claimed_at=clock_timestamp(),claimed_by=$2,claim_token=$3,claimed_until=clock_timestamp()+$4::interval,last_error_code=null,provider_reference=null,result_sha256_hex=null,updated_at=clock_timestamp()
    from candidate c
   where x.tenant_id=c.tenant_id and x.request_id=c.request_id
  returning x.tenant_id,x.request_id,x.attempt_count,x.claim_token
)
select c.tenant_id,c.request_id,e.disposition_id,e.effect_kind,e.owner_context,e.idempotency_key,c.attempt_count,c.claim_token,d.inventory_action,r.organization_id,r.stock_unit_id
  from claimed c
  join sales.return_effect_request e on e.tenant_id=c.tenant_id and e.request_id=c.request_id
  join sales.return_disposition d on d.tenant_id=e.tenant_id and d.disposition_id=e.disposition_id
  join sales.return_receipt r on r.tenant_id=d.tenant_id and r.receipt_id=d.receipt_id`, owner, worker, claimToken, lease.String()).Scan(
		&work.TenantID, &work.RequestID, &work.DispositionID, &work.EffectKind, &work.OwnerContext, &work.IdempotencyKey, &work.Attempt, &work.ClaimToken, &work.InventoryAction, &work.OrganizationID, &work.StockUnitID,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &work, nil
}

func inventoryTarget(action string) (string, bool) {
	targets := map[string]string{"quarantine": "quarantine", "restock": "available", "repair": "service", "scrap": "retired"}
	target, ok := targets[action]
	return target, ok
}

func (r *ReturnEffects) ApplyInventory(ctx context.Context, work returneffects.Work, worker string) (returneffects.Result, error) {
	target, ok := inventoryTarget(work.InventoryAction)
	if !ok || work.OwnerContext != "inventory" || work.EffectKind != "inventory" {
		return returneffects.Result{}, returneffects.ErrInventoryState
	}
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return returneffects.Result{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var attempt int
	var startedAt time.Time
	var stockState string
	var stockVersion int64
	err = tx.QueryRow(ctx, `
select x.attempt_count,x.claimed_at,s.state,s.version
  from sales.return_effect_execution x
  join sales.return_effect_request e on e.tenant_id=x.tenant_id and e.request_id=x.request_id
  join sales.return_disposition d on d.tenant_id=e.tenant_id and d.disposition_id=e.disposition_id
  join sales.return_receipt r on r.tenant_id=d.tenant_id and r.receipt_id=d.receipt_id
  join inventory.stock_unit s on s.tenant_id=r.tenant_id and s.stock_unit_id=r.stock_unit_id and s.organization_id=r.organization_id
 where x.tenant_id=$1 and x.request_id=$2 and x.status='claimed' and x.claimed_by=$3 and x.claim_token=$4 and x.claimed_until>=clock_timestamp()
   and e.effect_kind='inventory' and e.owner_context='inventory' and d.inventory_action=$5 and r.organization_id=$6 and r.stock_unit_id=$7
 for update of x,s`, work.TenantID, work.RequestID, worker, work.ClaimToken, work.InventoryAction, work.OrganizationID, work.StockUnitID).Scan(&attempt, &startedAt, &stockState, &stockVersion)
	if errors.Is(err, pgx.ErrNoRows) {
		return returneffects.Result{}, returneffects.ErrInventoryState
	}
	if err != nil {
		return returneffects.Result{}, err
	}
	if attempt != work.Attempt || stockState != "sold" {
		return returneffects.Result{}, returneffects.ErrInventoryState
	}

	newVersion := stockVersion + 1
	command, err := json.Marshal(map[string]any{"request_id": work.RequestID, "stock_unit_id": work.StockUnitID, "from_state": stockState, "to_state": target, "version": newVersion})
	if err != nil {
		return returneffects.Result{}, err
	}
	digest := sha256.Sum256(command)
	resultHash := hex.EncodeToString(digest[:])

	updated, err := tx.Exec(ctx, `update inventory.stock_unit set state=$5,version=$6,received_at=coalesce(received_at,clock_timestamp()),updated_at=clock_timestamp() where tenant_id=$1 and organization_id=$2 and stock_unit_id=$3 and state=$4 and version=$7`, work.TenantID, work.OrganizationID, work.StockUnitID, stockState, target, newVersion, stockVersion)
	if err != nil {
		return returneffects.Result{}, err
	}
	if updated.RowsAffected() != 1 {
		return returneffects.Result{}, returneffects.ErrInventoryState
	}

	if _, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload) values($1,gen_random_uuid(),'stock-unit',$2,$3,$4,1,clock_timestamp(),$5)`, work.TenantID, work.StockUnitID, newVersion, "stock-unit.return-"+target, command); err != nil {
		return returneffects.Result{}, err
	}
	if _, err = tx.Exec(ctx, `insert into sales.return_effect_attempt(tenant_id,attempt_id,request_id,attempt_no,worker_id,outcome,provider_reference,result_sha256_hex,started_at) values($1,$2,$3,$4,$5,'succeeded',$6,$7,$8)`, work.TenantID, work.ClaimToken, work.RequestID, attempt, worker, work.StockUnitID, resultHash, startedAt); err != nil {
		return returneffects.Result{}, err
	}
	completed, err := tx.Exec(ctx, `update sales.return_effect_execution set status='succeeded',claimed_at=null,claimed_by=null,claim_token=null,claimed_until=null,last_error_code=null,provider_reference=$5,result_sha256_hex=$6,updated_at=clock_timestamp() where tenant_id=$1 and request_id=$2 and status='claimed' and claimed_by=$3 and claim_token=$4`, work.TenantID, work.RequestID, worker, work.ClaimToken, work.StockUnitID, resultHash)
	if err != nil {
		return returneffects.Result{}, err
	}
	if completed.RowsAffected() != 1 {
		return returneffects.Result{}, returneffects.ErrInventoryState
	}
	if err = tx.Commit(ctx); err != nil {
		return returneffects.Result{}, err
	}
	return returneffects.Result{RequestID: work.RequestID, StockUnitID: work.StockUnitID, InventoryState: target, StockVersion: newVersion, ResultSHA256: resultHash}, nil
}

func (r *ReturnEffects) Finish(ctx context.Context, work returneffects.Work, worker string, completion returneffects.Completion) error {
	if completion.Outcome != "retry" && completion.Outcome != "blocked" && completion.Outcome != "failed" && completion.Outcome != "succeeded" {
		return fmt.Errorf("invalid completion outcome")
	}
	if completion.Outcome == "succeeded" {
		if completion.ErrorCode != "" || len(completion.ResultSHA256) != 64 {
			return fmt.Errorf("invalid successful completion")
		}
	} else if completion.ErrorCode == "" || completion.ResultSHA256 != "" {
		return fmt.Errorf("invalid unsuccessful completion")
	}
	if completion.Outcome == "retry" && completion.RetryAfter <= 0 {
		return fmt.Errorf("retry delay required")
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var attempt int
	var startedAt time.Time
	err = tx.QueryRow(ctx, `select attempt_count,claimed_at from sales.return_effect_execution where tenant_id=$1 and request_id=$2 and status='claimed' and claimed_by=$3 and claim_token=$4 and claimed_until>=clock_timestamp() for update`, work.TenantID, work.RequestID, worker, work.ClaimToken).Scan(&attempt, &startedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return returneffects.ErrInventoryState
	}
	if err != nil {
		return err
	}
	if attempt != work.Attempt {
		return returneffects.ErrInventoryState
	}
	if _, err = tx.Exec(ctx, `insert into sales.return_effect_attempt(tenant_id,attempt_id,request_id,attempt_no,worker_id,outcome,error_code,provider_reference,result_sha256_hex,started_at) values($1,$2,$3,$4,$5,$6,nullif($7,''),nullif($8,''),nullif($9,''),$10)`, work.TenantID, work.ClaimToken, work.RequestID, attempt, worker, completion.Outcome, completion.ErrorCode, completion.ProviderReference, completion.ResultSHA256, startedAt); err != nil {
		return err
	}
	available := time.Now().UTC()
	if completion.Outcome == "retry" {
		available = available.Add(completion.RetryAfter)
	}
	updated, err := tx.Exec(ctx, `update sales.return_effect_execution set status=$5,available_at=$6,claimed_at=null,claimed_by=null,claim_token=null,claimed_until=null,last_error_code=nullif($7,''),provider_reference=nullif($8,''),result_sha256_hex=nullif($9,''),updated_at=clock_timestamp() where tenant_id=$1 and request_id=$2 and status='claimed' and claimed_by=$3 and claim_token=$4`, work.TenantID, work.RequestID, worker, work.ClaimToken, completion.Outcome, available, completion.ErrorCode, completion.ProviderReference, completion.ResultSHA256)
	if err != nil {
		return err
	}
	if updated.RowsAffected() != 1 {
		return returneffects.ErrInventoryState
	}
	return tx.Commit(ctx)
}

func (r *ReturnEffects) ResumeBlocked(ctx context.Context, tenant, requestID, resumeID, actor, reason string) error {
	if tenant == "" || requestID == "" || resumeID == "" || actor == "" || reason == "" {
		return fmt.Errorf("invalid return effect resume")
	}
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	if _, err = tx.Exec(ctx, `insert into sales.return_effect_resume(tenant_id,resume_id,request_id,reason_code,requested_by_subject) values($1,$2,$3,$4,$5)`, tenant, resumeID, requestID, reason, actor); err != nil {
		return err
	}
	updated, err := tx.Exec(ctx, `update sales.return_effect_execution set status='requested',available_at=clock_timestamp(),last_error_code=null,provider_reference=null,result_sha256_hex=null,updated_at=clock_timestamp() where tenant_id=$1 and request_id=$2 and status='blocked'`, tenant, requestID)
	if err != nil {
		return err
	}
	if updated.RowsAffected() != 1 {
		return returneffects.ErrInventoryState
	}
	return tx.Commit(ctx)
}
````

### FILE: `internal/platform/postgres/returneffects_integration_test.go`

```yaml
block_id: "GO-RETURN-EFFECT:file:07"
operation: CREATE
provenance: AUTHORED
source: "local implementation governed by official sources in metadata"
license: "LicenseRef-Workspace-Owner"
sha256: "c93ca2701f59e559a46da3830220bcf7a6b370f7295c8b285db3128cfd3a0ae3"
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
	"time"

	"elite.local/enterprise/internal/returneffects"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestReturnEffectsInventoryClaimFencingAndResume(t *testing.T) {
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
	tenant := "018f4d4a-7b36-7a21-8d10-2f4c54c2e101"
	cleanup := func() {
		tx, cleanupErr := pool.Begin(ctx)
		if cleanupErr != nil {
			return
		}
		defer tx.Rollback(ctx)
		_, _ = tx.Exec(ctx, `set local session_replication_role=replica`)
		for _, table := range []string{"sales.return_effect_resume", "sales.return_effect_attempt", "sales.return_effect_execution", "sales.return_effect_request", "sales.return_disposition", "sales.return_receipt", "sales.return_authorization", "platform.outbox_event", "inventory.stock_unit", "org.organization", "platform.tenant"} {
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
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,'return-effects','Return Effects','Return Effects')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,'store','return-effects-store','Store','store')`,
		`insert into inventory.stock_unit(tenant_id,stock_unit_id,organization_id,variant_id,serial_number,vin,battery_serial_number,state,version,received_at) values($1,'stock-1','store','variant','SERIAL-1','VIN-1','BATTERY-1','sold',1,clock_timestamp()),($1,'stock-2','store','variant','SERIAL-2','VIN-2','BATTERY-2','sold',1,clock_timestamp())`,
		`insert into sales.return_authorization(tenant_id,authorization_id,organization_id,exception_id,handover_id,order_id,stock_unit_id,customer_principal_id,disposition,exact_cost_source_order_id,state,authorized_by_subject) values($1,'auth-1','store','exception-1','handover-1','order-1','stock-1','customer','return','order-1','authorized','operator'),($1,'auth-2','store','exception-2','handover-2','order-2','stock-2','customer','return','order-2','authorized','operator')`,
		`insert into sales.return_receipt(tenant_id,receipt_id,authorization_id,organization_id,order_id,stock_unit_id,customer_principal_id,received_serial_number,condition_code,notes,evidence_sha256_hex,received_by_subject) values($1,'receipt-1','auth-1','store','order-1','stock-1','customer','SERIAL-1','damaged','received','aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa','operator'),($1,'receipt-2','auth-2','store','order-2','stock-2','customer','SERIAL-2','opened','received','bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb','operator')`,
		`insert into sales.return_disposition(tenant_id,disposition_id,receipt_id,inventory_action,customer_remedy,notes,decided_by_subject) values($1,'disposition-1','receipt-1','quarantine','refund','inspect','operator'),($1,'disposition-2','receipt-2','restock','refund','inspect','operator')`,
	}
	for _, query := range fixtures {
		if _, err = tx.Exec(ctx, query, tenant); err != nil {
			t.Fatal(err)
		}
	}
	if err = tx.Commit(ctx); err != nil {
		t.Fatal(err)
	}

	if _, err = pool.Exec(ctx, `insert into sales.return_effect_request(tenant_id,request_id,disposition_id,effect_kind,owner_context,state,idempotency_key) values($1,'effect-1','disposition-1','inventory','inventory','requested','return-effect-key-0001'),($1,'effect-2','disposition-2','inventory','inventory','requested','return-effect-key-0002')`, tenant); err != nil {
		t.Fatal(err)
	}
	store := NewReturnEffects(pool)

	type claimResult struct {
		work *returneffects.Work
		err  error
	}
	start := make(chan struct{})
	claims := make(chan claimResult, 2)
	for _, input := range []struct{ worker, token string }{{"worker-a", "claim-a"}, {"worker-b", "claim-b"}} {
		go func(worker, token string) {
			<-start
			work, claimErr := store.Claim(ctx, "inventory", worker, token, time.Minute)
			claims <- claimResult{work, claimErr}
		}(input.worker, input.token)
	}
	close(start)
	claimed := map[string]*returneffects.Work{}
	for range 2 {
		result := <-claims
		if result.err != nil || result.work == nil {
			t.Fatalf("claim=%+v err=%v", result.work, result.err)
		}
		if _, exists := claimed[result.work.RequestID]; exists {
			t.Fatal("same request claimed concurrently")
		}
		claimed[result.work.RequestID] = result.work
	}
	if len(claimed) != 2 {
		t.Fatalf("claimed=%d", len(claimed))
	}

	first := claimed["effect-1"]
	if first == nil {
		t.Fatal("effect-1 not claimed")
	}
	firstWorker := "worker-a"
	if first.ClaimToken == "claim-b" {
		firstWorker = "worker-b"
	}
	if _, err = store.ApplyInventory(ctx, *first, "wrong-worker"); !errors.Is(err, returneffects.ErrInventoryState) {
		t.Fatalf("wrong worker applied effect: %v", err)
	}
	result, err := store.ApplyInventory(ctx, *first, firstWorker)
	if err != nil || result.InventoryState != "quarantine" || result.StockVersion != 2 || len(result.ResultSHA256) != 64 {
		t.Fatalf("result=%+v err=%v", result, err)
	}
	var state, executionStatus string
	var version int64
	var attempts, events int
	if err = pool.QueryRow(ctx, `select state,version from inventory.stock_unit where tenant_id=$1 and stock_unit_id='stock-1'`, tenant).Scan(&state, &version); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select status from sales.return_effect_execution where tenant_id=$1 and request_id='effect-1'`, tenant).Scan(&executionStatus); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select count(*) from sales.return_effect_attempt where tenant_id=$1 and request_id='effect-1'`, tenant).Scan(&attempts); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select count(*) from platform.outbox_event where tenant_id=$1 and aggregate_type='stock-unit' and aggregate_id='stock-1' and event_type='stock-unit.return-quarantine' and payload->>'request_id'='effect-1'`, tenant).Scan(&events); err != nil {
		t.Fatal(err)
	}
	if state != "quarantine" || version != 2 || executionStatus != "succeeded" || attempts != 1 || events != 1 {
		t.Fatalf("state=%s version=%d status=%s attempts=%d events=%d", state, version, executionStatus, attempts, events)
	}
	if work, claimErr := store.Claim(ctx, "inventory", "worker-c", "claim-c", time.Minute); claimErr != nil || work != nil {
		t.Fatalf("completed work reclaimed: %+v %v", work, claimErr)
	}

	second := claimed["effect-2"]
	if second == nil {
		t.Fatal("effect-2 not claimed")
	}
	secondWorker := "worker-a"
	if second.ClaimToken == "claim-b" {
		secondWorker = "worker-b"
	}
	if err = store.Finish(ctx, *second, secondWorker, returneffects.Completion{Outcome: "blocked", ErrorCode: "INVENTORY_STATE_CONFLICT"}); err != nil {
		t.Fatal(err)
	}
	if err = store.ResumeBlocked(ctx, tenant, "effect-2", "resume-1", "operator", "INVENTORY_RECONCILED"); err != nil {
		t.Fatal(err)
	}
	reclaimed, err := store.Claim(ctx, "inventory", "worker-c", "claim-c", time.Minute)
	if err != nil || reclaimed == nil || reclaimed.RequestID != "effect-2" || reclaimed.Attempt != 2 {
		t.Fatalf("reclaimed=%+v err=%v", reclaimed, err)
	}
	if err = store.Finish(ctx, *second, secondWorker, returneffects.Completion{Outcome: "failed", ErrorCode: "STALE_CLAIM"}); !errors.Is(err, returneffects.ErrInventoryState) {
		t.Fatalf("stale claim completed: %v", err)
	}
	result, err = store.ApplyInventory(ctx, *reclaimed, "worker-c")
	if err != nil || result.InventoryState != "available" || result.StockVersion != 2 {
		t.Fatalf("restock result=%+v err=%v", result, err)
	}
	if _, err = pool.Exec(ctx, `update sales.return_effect_attempt set error_code='MUTATED' where tenant_id=$1 and request_id='effect-2'`, tenant); err == nil {
		t.Fatal("attempt audit mutation accepted")
	}
	if _, err = pool.Exec(ctx, `delete from sales.return_effect_resume where tenant_id=$1 and request_id='effect-2'`, tenant); err == nil {
		t.Fatal("resume audit deletion accepted")
	}
}
````

### FILE: `cmd/return-effect-worker/main.go`

```yaml
block_id: "GO-RETURN-EFFECT:file:08"
operation: CREATE
provenance: AUTHORED
source: "local implementation governed by official sources in metadata"
license: "LicenseRef-Workspace-Owner"
sha256: "cf9ab4ed3239e21f0a192725c34eb6f0b4a5e8a25fa61f904bdaa3ff31dd4a18"
variables: []
secrets_allowed: false
```

````go
package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/platform/randomid"
	"elite.local/enterprise/internal/returneffects"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	databaseURL := os.Getenv("DATABASE_URL")
	workerID := os.Getenv("RETURN_EFFECT_WORKER_ID")
	if databaseURL == "" || workerID == "" {
		slog.Error("DATABASE_URL and RETURN_EFFECT_WORKER_ID are required")
		os.Exit(2)
	}
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		slog.Error("database configuration is invalid")
		os.Exit(2)
	}
	config.MaxConns = 4
	config.MinConns = 0
	config.MaxConnLifetime = 30 * time.Minute
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		slog.Error("database pool creation failed")
		os.Exit(1)
	}
	defer pool.Close()
	if err = pool.Ping(ctx); err != nil {
		slog.Error("database unavailable")
		os.Exit(1)
	}
	processor, err := returneffects.NewInventoryProcessor(postgres.NewReturnEffects(pool), randomid.Generator{}, workerID, 2*time.Minute, 5*time.Second)
	if err != nil {
		slog.Error("return effect worker configuration is invalid")
		os.Exit(2)
	}
	if err = run(ctx, processor); err != nil && !errors.Is(err, context.Canceled) {
		slog.Error("return effect worker stopped unexpectedly")
		os.Exit(1)
	}
}

func run(ctx context.Context, processor *returneffects.Processor) error {
	timer := time.NewTimer(0)
	if !timer.Stop() {
		<-timer.C
	}
	defer timer.Stop()
	for {
		_, err := processor.ProcessOne(ctx)
		if err == nil {
			continue
		}
		if errors.Is(err, context.Canceled) {
			return err
		}
		wait := 2 * time.Second
		if errors.Is(err, returneffects.ErrNoWork) {
			wait = 500 * time.Millisecond
		} else {
			slog.Warn("return inventory effect deferred")
		}
		timer.Reset(wait)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
		}
	}
}
````

## 6. Configuration surface

`DATABASE_URL` and `RETURN_EFFECT_WORKER_ID` are mandatory runtime inputs. Database credentials remain external secrets. The worker uses fixed bounded defaults: four database connections, two-minute lease, five-second retry and 500 ms idle poll. A selected project must size these against measured load and deployment termination budgets. Only one owner (`inventory`) is admitted by this binary.

## 7. Dependency bill

| Dependency | Fixed baseline | Purpose | License/authority |
|---|---:|---|---|
| Go | 1.26.7 | worker/runtime/tests | Go project |
| PostgreSQL | 18.6 | projection, locks, audit, transaction, outbox | PostgreSQL License |
| pgx | 5.10.0 | PostgreSQL driver/pool | MIT |
| Existing journey/inventory/outbox schemas | compatible pack versions above | authoritative return and stock ownership | workspace owner |

No Stripe, ARCA, cloud or queue SDK is added. Those adapters remain separate conditioned packs with their own official artifacts, credentials and target gates.

## 8. Apply order

Compose after migrations 0001–0018, the reliable worker primitives and the existing application wiring. Apply 0019, run its SQL contract, then build all Go binaries. Deploy the inventory worker only after a target database backup and rollback rehearsal. Stop claims, drain or let leases expire, then run the down migration only when no project depends on execution history; production rollback should normally preserve audit rows and deploy the prior binary.

## 9. Verification

Required gates: materialize byte-for-byte; `gofmt`; unit tests; full `go test ./...`; `go vet ./...`; build every actual `cmd` main; apply all 19 migrations to clean PostgreSQL 18.6; run the 0019 SQL test; run concurrent real-database integration; prove wrong-worker rejection, claim fencing, atomic stock/outbox/result, duplicate suppression, blocked resume, immutable attempts/resumes and down/up. A project remains conditioned until target load, security, recovery, observability, deployment and business acceptance pass.

## 10. Reconstruction evidence

V141 evidence records clean Markdown composition, hashes, Go gates, PostgreSQL 18.6 migrations/rollback and the focused concurrency integration. The evidence admits only the reusable library component. It does not claim a live refund provider, replacement shipment, accounting reversal, ARCA credit note or production readiness.
