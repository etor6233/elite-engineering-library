# Go Return Exchange Fulfillment Worker

## 1. Metadata

```yaml
pack_id: "GO-RETURN-EXCHANGE-FULFILLMENT-WORKER"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Consumes the durable fulfillment-owned exchange effect and atomically creates a linked replacement order, reserves an exact stock unit, prepares a delivery handover and records immutable execution evidence; it does not invent price, tax, accounting, carrier or remote-delivery policy."
stacks: ["Go 1.26.7", "PostgreSQL 18.6", "pgx 5.10.0"]
compatible_with: ["GO-ENTERPRISE-BACKEND 0.4.x", "GO-RETURN-EFFECT-EXECUTION-WORKER 0.1.x", "GO-FRANCHISE-CUSTOMER-JOURNEY-API 0.8.x", "GO-ELECTROMOBILITY-APPLICATION 1.5.x"]
incompatible_with: ["alternate replacement-order or stock-reservation owner without reconciliation contract"]
license_expression: "LicenseRef-Workspace-Owner AND MIT dependency"
upstream_sources: ["https://learn.microsoft.com/en-us/dynamics365/supply-chain/sales-marketing/sales-returns", "https://learn.microsoft.com/en-us/dynamics365/supply-chain/sales-marketing/create-item-replacement-order", "https://learn.microsoft.com/en-us/dynamics365/commerce/orderexchanges", "https://learn.microsoft.com/en-us/dynamics365/business-central/sales-how-process-sales-returns-cancellations", "https://www.postgresql.org/docs/18/explicit-locking.html", "https://www.postgresql.org/docs/18/sql-select.html", "https://docs.aws.amazon.com/wellarchitected/latest/reliability-pillar/welcome.html"]
verified_at: "2026-08-31"
```

All eight files are local `AUTHORED` implementation. Microsoft governs independent replacement orders, original-sale linkage, disposition timing, exact-cost separation and net exchange settlement; PostgreSQL governs serializable transactions, row locking and `SKIP LOCKED`; AWS governs durable retry and idempotent recovery. No local line is attributed to those companies.

## 2. Applicability

Use after return authorization, physical receipt, disposition and durable effect execution when the business admits an even, like-for-like exchange at the same organization. The admitted automatic lane is one delivered order with exactly one quantity-one line, exact original stock linkage, identical variant replacement, no implicit discounts/taxes/shipping allocation, completed physical inventory effect and a separate accounting request. Different variant, split order, price difference, cross-organization fulfillment or remote carrier delivery remains blocked until a project supplies approved policy and its executable owner.

## 3. Architecture contract

Microsoft's return model creates a replacement sales order associated with the RMA instead of mutating the original sale. This pack follows that boundary: one immutable `sales.return_exchange` links original order/line/stock to a new zero-balance replacement order/line, reserved stock, prepared handover and accounting request. `prepared` proves only local fulfillment preparation; it does not mean delivered, financially posted, fiscally accepted or customer accepted.

The worker claims only `effect_kind=exchange` and `owner_context=fulfillment`. A serializable transaction fences the lease, proves the original single-line sale, proves the inventory effect succeeded, locks one exact available same-variant stock unit with `FOR UPDATE SKIP LOCKED`, creates the replacement order and handover, reserves stock, writes three outbox events, immutable attempt and terminal execution together. No stock is retryable; a contract ambiguity is blocked. A stale worker cannot complete another lease.

## 4. Exact file manifest

```text
CREATE db/migrations/0021_return_exchange_fulfillment.up.sql
CREATE db/migrations/0021_return_exchange_fulfillment.down.sql
CREATE db/tests/0021_return_exchange_fulfillment.test.sql
CREATE internal/returnexchange/processor.go
CREATE internal/returnexchange/processor_test.go
CREATE internal/platform/postgres/returnexchange.go
CREATE internal/platform/postgres/returnexchange_integration_test.go
CREATE cmd/return-exchange-worker/main.go
```

## 5. Materialization blocks

### FILE: `db/migrations/0021_return_exchange_fulfillment.up.sql`
```yaml
block_id: "GO-RETURN-EXCHANGE:file:01"
operation: CREATE
provenance: AUTHORED
source: "local implementation governed by official sources in metadata"
license: "LicenseRef-Workspace-Owner"
sha256: "4c668ae536f898a5833ef5f01743b5acb78447d078494460a87562173ff189b9"
variables: []
secrets_allowed: false
```
````sql
begin;

create table sales.return_exchange (
  tenant_id uuid not null,
  request_id text not null,
  disposition_id text not null,
  original_order_id text not null,
  original_line_id text not null,
  original_stock_unit_id text not null,
  replacement_order_id text not null,
  replacement_line_id text not null,
  replacement_stock_unit_id text not null,
  replacement_handover_id text not null,
  accounting_request_id text not null,
  currency text not null check (currency ~ '^[A-Z]{3}$'),
  original_unit_price_minor_units bigint not null check (original_unit_price_minor_units > 0),
  settlement_mode text not null check (settlement_mode='even-exchange-zero-balance'),
  status text not null check (status='prepared'),
  result_sha256_hex text not null check (result_sha256_hex ~ '^[0-9a-f]{64}$'),
  prepared_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id,request_id),
  unique (tenant_id,disposition_id),
  unique (tenant_id,replacement_order_id),
  unique (tenant_id,replacement_stock_unit_id),
  unique (tenant_id,replacement_handover_id),
  foreign key (tenant_id,request_id) references sales.return_effect_request(tenant_id,request_id),
  foreign key (tenant_id,disposition_id) references sales.return_disposition(tenant_id,disposition_id),
  foreign key (tenant_id,original_order_id,original_line_id) references sales.customer_order_line(tenant_id,order_id,line_id),
  foreign key (tenant_id,original_stock_unit_id) references inventory.stock_unit(tenant_id,stock_unit_id),
  foreign key (tenant_id,replacement_order_id,replacement_line_id) references sales.customer_order_line(tenant_id,order_id,line_id),
  foreign key (tenant_id,replacement_stock_unit_id) references inventory.stock_unit(tenant_id,stock_unit_id),
  foreign key (tenant_id,replacement_handover_id) references sales.delivery_handover(tenant_id,handover_id),
  foreign key (tenant_id,accounting_request_id) references sales.return_effect_request(tenant_id,request_id),
  check (original_order_id<>replacement_order_id),
  check (original_stock_unit_id<>replacement_stock_unit_id)
);

create or replace function sales.prevent_return_exchange_mutation() returns trigger language plpgsql as $$
begin
  raise exception using errcode='23514',message='prepared return exchange is immutable';
end $$;

create trigger return_exchange_immutable before update or delete on sales.return_exchange for each row execute function sales.prevent_return_exchange_mutation();
create index return_exchange_original_order_idx on sales.return_exchange(tenant_id,original_order_id,prepared_at desc,request_id);

commit;
````

### FILE: `db/migrations/0021_return_exchange_fulfillment.down.sql`
```yaml
block_id: "GO-RETURN-EXCHANGE:file:02"
operation: CREATE
provenance: AUTHORED
source: "local implementation governed by official sources in metadata"
license: "LicenseRef-Workspace-Owner"
sha256: "8a082205951f03b68266260e5ade81df27fe8302afd41ce4e8afd95aaeec018c"
variables: []
secrets_allowed: false
```
````sql
begin;

drop index if exists sales.return_exchange_original_order_idx;
drop trigger if exists return_exchange_immutable on sales.return_exchange;
drop function if exists sales.prevent_return_exchange_mutation();
drop table if exists sales.return_exchange;

commit;
````

### FILE: `db/tests/0021_return_exchange_fulfillment.test.sql`
```yaml
block_id: "GO-RETURN-EXCHANGE:file:03"
operation: CREATE
provenance: AUTHORED
source: "local implementation governed by official sources in metadata"
license: "LicenseRef-Workspace-Owner"
sha256: "76ac9e9d79c0a894f00d8a776f0b63344b030917aa0ce03a70ad76a889958d50"
variables: []
secrets_allowed: false
```
````sql
begin;

do $$
begin
  if to_regclass('sales.return_exchange') is null then
    raise exception 'return exchange table missing';
  end if;
  if (select count(*) from pg_trigger where tgname='return_exchange_immutable' and not tgisinternal) <> 1 then
    raise exception 'return exchange immutability trigger missing';
  end if;
  if not exists (select 1 from pg_indexes where schemaname='sales' and indexname='return_exchange_original_order_idx') then
    raise exception 'return exchange lookup index missing';
  end if;
  if exists (select 1 from sales.return_exchange where original_order_id=replacement_order_id or original_stock_unit_id=replacement_stock_unit_id) then
    raise exception 'return exchange reused original order or stock';
  end if;
end $$;

rollback;
````

### FILE: `internal/returnexchange/processor.go`
```yaml
block_id: "GO-RETURN-EXCHANGE:file:04"
operation: CREATE
provenance: AUTHORED
source: "local implementation governed by official sources in metadata"
license: "LicenseRef-Workspace-Owner"
sha256: "c932aa8f159473ed54e94cf985ace6e251b82ac94393e6b7fda93ffcbe60f4a6"
variables: []
secrets_allowed: false
```
````go
package returnexchange

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	"elite.local/enterprise/internal/returneffects"
)

var (
	ErrNoWork           = errors.New("no return exchange work")
	ErrConflict         = errors.New("return exchange contract conflict")
	ErrStockUnavailable = errors.New("replacement stock unavailable")
)

var workerIDPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:-]{0,63}$`)

type Result struct {
	RequestID             string
	ReplacementOrderID    string
	ReplacementStockID    string
	ReplacementHandoverID string
	ResultSHA256          string
}

type Store interface {
	Claim(context.Context, string, string, string, time.Duration) (*returneffects.Work, error)
	PrepareExchange(context.Context, returneffects.Work, string) (Result, error)
	Finish(context.Context, returneffects.Work, string, returneffects.Completion) error
}

type IDGenerator interface{ New() string }

type Processor struct {
	store      Store
	ids        IDGenerator
	workerID   string
	lease      time.Duration
	retryDelay time.Duration
}

func NewProcessor(store Store, ids IDGenerator, workerID string, lease, retryDelay time.Duration) (*Processor, error) {
	if store == nil || ids == nil || !workerIDPattern.MatchString(workerID) || lease < time.Second || lease > 15*time.Minute || retryDelay < time.Second || retryDelay > time.Hour {
		return nil, fmt.Errorf("invalid return exchange processor configuration")
	}
	return &Processor{store: store, ids: ids, workerID: workerID, lease: lease, retryDelay: retryDelay}, nil
}

func (p *Processor) ProcessOne(ctx context.Context) (Result, error) {
	work, err := p.store.Claim(ctx, "fulfillment", p.workerID, p.ids.New(), p.lease)
	if err != nil {
		return Result{}, err
	}
	if work == nil {
		return Result{}, ErrNoWork
	}
	result, executeErr := p.store.PrepareExchange(ctx, *work, p.workerID)
	if executeErr == nil {
		return result, nil
	}
	completion := returneffects.Completion{Outcome: "retry", ErrorCode: "EXCHANGE_TRANSIENT_FAILURE", RetryAfter: p.retryDelay}
	if errors.Is(executeErr, ErrStockUnavailable) {
		completion.ErrorCode = "REPLACEMENT_STOCK_UNAVAILABLE"
	}
	if errors.Is(executeErr, ErrConflict) {
		completion = returneffects.Completion{Outcome: "blocked", ErrorCode: "EXCHANGE_CONTRACT_CONFLICT"}
	}
	if finishErr := p.store.Finish(ctx, *work, p.workerID, completion); finishErr != nil {
		return Result{}, errors.Join(executeErr, finishErr)
	}
	return Result{}, executeErr
}
````

### FILE: `internal/returnexchange/processor_test.go`
```yaml
block_id: "GO-RETURN-EXCHANGE:file:05"
operation: CREATE
provenance: AUTHORED
source: "local implementation governed by official sources in metadata"
license: "LicenseRef-Workspace-Owner"
sha256: "fb7352ec5bd50fc48cc910b473cfd858d73be3d098d999f8b3faa2819d14e2ca"
variables: []
secrets_allowed: false
```
````go
package returnexchange

import (
	"context"
	"errors"
	"testing"
	"time"

	"elite.local/enterprise/internal/returneffects"
)

type fakeIDs struct{ value string }

func (f fakeIDs) New() string { return f.value }

type fakeStore struct {
	work       *returneffects.Work
	result     Result
	executeErr error
	completion returneffects.Completion
}

func (f *fakeStore) Claim(context.Context, string, string, string, time.Duration) (*returneffects.Work, error) {
	return f.work, nil
}
func (f *fakeStore) PrepareExchange(context.Context, returneffects.Work, string) (Result, error) {
	return f.result, f.executeErr
}
func (f *fakeStore) Finish(_ context.Context, _ returneffects.Work, _ string, completion returneffects.Completion) error {
	f.completion = completion
	return nil
}

func TestProcessorSuccessAndFailureClassification(t *testing.T) {
	work := &returneffects.Work{RequestID: "exchange-1"}
	store := &fakeStore{work: work, result: Result{RequestID: "exchange-1", ReplacementOrderID: "order-2"}}
	processor, err := NewProcessor(store, fakeIDs{"claim-1"}, "exchange-worker", time.Minute, 5*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	result, err := processor.ProcessOne(context.Background())
	if err != nil || result.ReplacementOrderID != "order-2" {
		t.Fatalf("result=%+v err=%v", result, err)
	}

	store.executeErr = ErrStockUnavailable
	if _, err = processor.ProcessOne(context.Background()); !errors.Is(err, ErrStockUnavailable) {
		t.Fatalf("stock error=%v", err)
	}
	if store.completion.Outcome != "retry" || store.completion.ErrorCode != "REPLACEMENT_STOCK_UNAVAILABLE" {
		t.Fatalf("completion=%+v", store.completion)
	}

	store.executeErr = ErrConflict
	if _, err = processor.ProcessOne(context.Background()); !errors.Is(err, ErrConflict) {
		t.Fatalf("conflict error=%v", err)
	}
	if store.completion.Outcome != "blocked" || store.completion.ErrorCode != "EXCHANGE_CONTRACT_CONFLICT" {
		t.Fatalf("completion=%+v", store.completion)
	}
}

func TestProcessorRejectsInvalidConfigurationAndNoWork(t *testing.T) {
	if _, err := NewProcessor(&fakeStore{}, fakeIDs{"claim"}, "bad worker", time.Minute, time.Second); err == nil {
		t.Fatal("invalid worker accepted")
	}
	processor, err := NewProcessor(&fakeStore{}, fakeIDs{"claim"}, "worker", time.Minute, time.Second)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = processor.ProcessOne(context.Background()); !errors.Is(err, ErrNoWork) {
		t.Fatalf("no work=%v", err)
	}
}
````

### FILE: `internal/platform/postgres/returnexchange.go`
```yaml
block_id: "GO-RETURN-EXCHANGE:file:06"
operation: CREATE
provenance: AUTHORED
source: "local implementation governed by official sources in metadata"
license: "LicenseRef-Workspace-Owner"
sha256: "f6843aba7f8e1d63ce13ee5be8c8fc98cba99bf335397758df39e3c6fe4edc6f"
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
	"time"

	"elite.local/enterprise/internal/returneffects"
	"elite.local/enterprise/internal/returnexchange"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReturnExchange struct {
	pool    *pgxpool.Pool
	effects *ReturnEffects
}

func NewReturnExchange(pool *pgxpool.Pool) *ReturnExchange {
	return &ReturnExchange{pool: pool, effects: NewReturnEffects(pool)}
}

func (r *ReturnExchange) Claim(ctx context.Context, owner, worker, token string, lease time.Duration) (*returneffects.Work, error) {
	return r.effects.Claim(ctx, owner, worker, token, lease)
}

func (r *ReturnExchange) Finish(ctx context.Context, work returneffects.Work, worker string, completion returneffects.Completion) error {
	return r.effects.Finish(ctx, work, worker, completion)
}

type exchangeSource struct {
	attempt           int
	startedAt         time.Time
	dispositionID     string
	organizationID    string
	originalOrderID   string
	originalLineID    string
	originalStockID   string
	customerID        string
	variantID         string
	currency          string
	unitPrice         int64
	accountingRequest string
}

func (r *ReturnExchange) PrepareExchange(ctx context.Context, work returneffects.Work, worker string) (returnexchange.Result, error) {
	if work.EffectKind != "exchange" || work.OwnerContext != "fulfillment" {
		return returnexchange.Result{}, returnexchange.ErrConflict
	}
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return returnexchange.Result{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var source exchangeSource
	err = tx.QueryRow(ctx, `
select x.attempt_count,x.claimed_at,e.disposition_id,r.organization_id,r.order_id,l.line_id,r.stock_unit_id,r.customer_principal_id,l.variant_id,o.currency,l.unit_price_minor_units,accounting.request_id
  from sales.return_effect_execution x
  join sales.return_effect_request e on e.tenant_id=x.tenant_id and e.request_id=x.request_id
  join sales.return_disposition d on d.tenant_id=e.tenant_id and d.disposition_id=e.disposition_id
  join sales.return_receipt r on r.tenant_id=d.tenant_id and r.receipt_id=d.receipt_id
  join sales.return_authorization a on a.tenant_id=r.tenant_id and a.authorization_id=r.authorization_id
  join sales.customer_order o on o.tenant_id=r.tenant_id and o.order_id=r.order_id
  join sales.customer_order_line l on l.tenant_id=o.tenant_id and l.order_id=o.order_id and l.allocated_stock_unit_id=r.stock_unit_id
  join sales.return_effect_request inventory on inventory.tenant_id=e.tenant_id and inventory.disposition_id=e.disposition_id and inventory.effect_kind='inventory' and inventory.owner_context='inventory'
  join sales.return_effect_execution inventory_x on inventory_x.tenant_id=inventory.tenant_id and inventory_x.request_id=inventory.request_id and inventory_x.status='succeeded'
  join sales.return_effect_request accounting on accounting.tenant_id=e.tenant_id and accounting.disposition_id=e.disposition_id and accounting.effect_kind='accounting' and accounting.owner_context='accounting'
 where x.tenant_id=$1 and x.request_id=$2 and x.status='claimed' and x.claimed_by=$3 and x.claim_token=$4 and x.claimed_until>=clock_timestamp()
   and e.effect_kind='exchange' and e.owner_context='fulfillment' and d.customer_remedy='exchange' and a.disposition='exchange'
   and o.state='delivered' and l.quantity=1 and l.unit_price_minor_units>0 and o.total_minor_units=l.unit_price_minor_units
   and (select count(*) from sales.customer_order_line all_lines where all_lines.tenant_id=o.tenant_id and all_lines.order_id=o.order_id)=1
 for update of x,o,l`, work.TenantID, work.RequestID, worker, work.ClaimToken).Scan(
		&source.attempt, &source.startedAt, &source.dispositionID, &source.organizationID, &source.originalOrderID, &source.originalLineID, &source.originalStockID, &source.customerID, &source.variantID, &source.currency, &source.unitPrice, &source.accountingRequest,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return returnexchange.Result{}, returnexchange.ErrConflict
	}
	if err != nil {
		return returnexchange.Result{}, err
	}
	if source.attempt != work.Attempt {
		return returnexchange.Result{}, returnexchange.ErrConflict
	}

	var replacementStock string
	var stockVersion int64
	err = tx.QueryRow(ctx, `select stock_unit_id,version from inventory.stock_unit where tenant_id=$1 and organization_id=$2 and variant_id=$3 and state='available' and stock_unit_id<>$4 order by received_at nulls last,stock_unit_id for update skip locked limit 1`, work.TenantID, source.organizationID, source.variantID, source.originalStockID).Scan(&replacementStock, &stockVersion)
	if errors.Is(err, pgx.ErrNoRows) {
		return returnexchange.Result{}, returnexchange.ErrStockUnavailable
	}
	if err != nil {
		return returnexchange.Result{}, err
	}

	var replacementOrder, replacementLine, replacementHandover string
	if err = tx.QueryRow(ctx, `select gen_random_uuid()::text,gen_random_uuid()::text,gen_random_uuid()::text`).Scan(&replacementOrder, &replacementLine, &replacementHandover); err != nil {
		return returnexchange.Result{}, err
	}
	if _, err = tx.Exec(ctx, `insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version) values($1,$2,$3,$4,'allocated',$5,0,1)`, work.TenantID, replacementOrder, source.organizationID, source.customerID, source.currency); err != nil {
		return returnexchange.Result{}, err
	}
	if _, err = tx.Exec(ctx, `insert into sales.customer_order_line(tenant_id,order_id,line_id,variant_id,quantity,unit_price_minor_units,allocated_stock_unit_id) values($1,$2,$3,$4,1,0,$5)`, work.TenantID, replacementOrder, replacementLine, source.variantID, replacementStock); err != nil {
		return returnexchange.Result{}, err
	}
	updated, err := tx.Exec(ctx, `update inventory.stock_unit set state='reserved',version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and organization_id=$2 and stock_unit_id=$3 and variant_id=$4 and state='available' and version=$5`, work.TenantID, source.organizationID, replacementStock, source.variantID, stockVersion)
	if err != nil {
		return returnexchange.Result{}, err
	}
	if updated.RowsAffected() != 1 {
		return returnexchange.Result{}, returnexchange.ErrStockUnavailable
	}
	if _, err = tx.Exec(ctx, `insert into sales.delivery_handover(tenant_id,handover_id,organization_id,order_id,customer_principal_id,stock_unit_id,state,version) values($1,$2,$3,$4,$5,$6,'prepared',1)`, work.TenantID, replacementHandover, source.organizationID, replacementOrder, source.customerID, replacementStock); err != nil {
		return returnexchange.Result{}, err
	}

	command, err := json.Marshal(struct {
		RequestID             string `json:"request_id"`
		OriginalOrderID       string `json:"original_order_id"`
		ReplacementOrderID    string `json:"replacement_order_id"`
		ReplacementStockID    string `json:"replacement_stock_unit_id"`
		ReplacementHandoverID string `json:"replacement_handover_id"`
	}{work.RequestID, source.originalOrderID, replacementOrder, replacementStock, replacementHandover})
	if err != nil {
		return returnexchange.Result{}, err
	}
	digest := sha256.Sum256(command)
	resultHash := hex.EncodeToString(digest[:])

	if _, err = tx.Exec(ctx, `insert into sales.return_exchange(tenant_id,request_id,disposition_id,original_order_id,original_line_id,original_stock_unit_id,replacement_order_id,replacement_line_id,replacement_stock_unit_id,replacement_handover_id,accounting_request_id,currency,original_unit_price_minor_units,settlement_mode,status,result_sha256_hex) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,'even-exchange-zero-balance','prepared',$14)`, work.TenantID, work.RequestID, source.dispositionID, source.originalOrderID, source.originalLineID, source.originalStockID, replacementOrder, replacementLine, replacementStock, replacementHandover, source.accountingRequest, source.currency, source.unitPrice, resultHash); err != nil {
		return returnexchange.Result{}, err
	}
	if _, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload) values($1,gen_random_uuid(),'return-exchange',$2,1,'return.exchange-prepared',1,clock_timestamp(),$3),($1,gen_random_uuid(),'order',$4,1,'order.exchange-replacement-created',1,clock_timestamp(),$3),($1,gen_random_uuid(),'stock-unit',$5,$6,'stock-unit.exchange-reserved',1,clock_timestamp(),$3)`, work.TenantID, work.RequestID, command, replacementOrder, replacementStock, stockVersion+1); err != nil {
		return returnexchange.Result{}, err
	}
	if _, err = tx.Exec(ctx, `insert into sales.return_effect_attempt(tenant_id,attempt_id,request_id,attempt_no,worker_id,outcome,provider_reference,result_sha256_hex,started_at) values($1,$2,$3,$4,$5,'succeeded',$6,$7,$8)`, work.TenantID, work.ClaimToken, work.RequestID, source.attempt, worker, replacementOrder, resultHash, source.startedAt); err != nil {
		return returnexchange.Result{}, err
	}
	completed, err := tx.Exec(ctx, `update sales.return_effect_execution set status='succeeded',claimed_at=null,claimed_by=null,claim_token=null,claimed_until=null,last_error_code=null,provider_reference=$5,result_sha256_hex=$6,updated_at=clock_timestamp() where tenant_id=$1 and request_id=$2 and status='claimed' and claimed_by=$3 and claim_token=$4`, work.TenantID, work.RequestID, worker, work.ClaimToken, replacementOrder, resultHash)
	if err != nil {
		return returnexchange.Result{}, err
	}
	if completed.RowsAffected() != 1 {
		return returnexchange.Result{}, returnexchange.ErrConflict
	}
	if err = tx.Commit(ctx); err != nil {
		return returnexchange.Result{}, err
	}
	return returnexchange.Result{RequestID: work.RequestID, ReplacementOrderID: replacementOrder, ReplacementStockID: replacementStock, ReplacementHandoverID: replacementHandover, ResultSHA256: resultHash}, nil
}
````

### FILE: `internal/platform/postgres/returnexchange_integration_test.go`
```yaml
block_id: "GO-RETURN-EXCHANGE:file:07"
operation: CREATE
provenance: AUTHORED
source: "local implementation governed by official sources in metadata"
license: "LicenseRef-Workspace-Owner"
sha256: "7bbd3a6bedb8541a5d51a20cb9cf4f9dfb7868bfbc90306fc28871f319acbd8c"
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
	"elite.local/enterprise/internal/returnexchange"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestReturnExchangeAtomicReservationAndNoStockRetry(t *testing.T) {
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
	tenant := "018f4d4a-7b36-7a21-8d10-2f4c54c2e143"
	cleanup := func() {
		tx, cleanupErr := pool.Begin(ctx)
		if cleanupErr != nil {
			return
		}
		defer tx.Rollback(ctx)
		_, _ = tx.Exec(ctx, `set local session_replication_role=replica`)
		for _, table := range []string{"sales.return_exchange", "sales.return_effect_resume", "sales.return_effect_attempt", "sales.return_effect_execution", "sales.return_effect_request", "sales.return_disposition", "sales.return_receipt", "sales.return_authorization", "sales.delivery_exception", "sales.delivery_handover", "sales.customer_order_line", "sales.customer_order", "platform.outbox_event", "inventory.stock_unit", "crm.customer_profile", "catalog.vehicle_variant", "catalog.vehicle_model", "org.organization", "platform.tenant"} {
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
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,'return-exchange','Return Exchange','Return Exchange')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,'store','return-exchange-store','Store','store')`,
		`insert into catalog.vehicle_model(tenant_id,model_id,model_code,display_name,vehicle_class,lifecycle_state) values($1,'model','return-exchange-model','Model','motorcycle','active')`,
		`insert into catalog.vehicle_variant(tenant_id,variant_id,model_id,variant_code,display_name,battery_specification,lifecycle_state) values($1,'variant','model','return-exchange-variant','Variant','{}','active')`,
		`insert into crm.customer_profile(tenant_id,customer_principal_id,display_name) values($1,'customer','Customer')`,
		`insert into inventory.stock_unit(tenant_id,stock_unit_id,organization_id,variant_id,serial_number,vin,battery_serial_number,state,version,received_at) values($1,'original-1','store','variant','ORIGINAL-1','VIN-O1','BAT-O1','quarantine',2,clock_timestamp()),($1,'original-2','store','variant','ORIGINAL-2','VIN-O2','BAT-O2','quarantine',2,clock_timestamp()),($1,'replacement-1','store','variant','REPLACEMENT-1','VIN-R1','BAT-R1','available',7,clock_timestamp())`,
		`insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version) values($1,'original-order-1','store','customer','delivered','ARS',500000,4),($1,'original-order-2','store','customer','delivered','ARS',500000,4)`,
		`insert into sales.customer_order_line(tenant_id,order_id,line_id,variant_id,quantity,unit_price_minor_units,allocated_stock_unit_id) values($1,'original-order-1','original-line-1','variant',1,500000,'original-1'),($1,'original-order-2','original-line-2','variant',1,500000,'original-2')`,
		`insert into sales.delivery_handover(tenant_id,handover_id,organization_id,order_id,customer_principal_id,stock_unit_id,state,acceptance_evidence_sha256_hex,customer_accepted_at,version) values($1,'original-handover-1','store','original-order-1','customer','original-1','accepted','aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa',clock_timestamp(),2),($1,'original-handover-2','store','original-order-2','customer','original-2','accepted','bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb',clock_timestamp(),2)`,
		`insert into sales.delivery_exception(tenant_id,exception_id,organization_id,handover_id,customer_principal_id,reason_code,details,rejection_evidence_sha256_hex,state,version,resolved_at,resolved_by_subject) values($1,'exception-1','store','original-handover-1','customer','defect','Defect','cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc','resolved',2,clock_timestamp(),'operator'),($1,'exception-2','store','original-handover-2','customer','defect','Defect','dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd','resolved',2,clock_timestamp(),'operator')`,
		`insert into sales.return_authorization(tenant_id,authorization_id,organization_id,exception_id,handover_id,order_id,stock_unit_id,customer_principal_id,disposition,exact_cost_source_order_id,state,authorized_by_subject) values($1,'auth-1','store','exception-1','original-handover-1','original-order-1','original-1','customer','exchange','original-order-1','authorized','operator'),($1,'auth-2','store','exception-2','original-handover-2','original-order-2','original-2','customer','exchange','original-order-2','authorized','operator')`,
		`insert into sales.return_receipt(tenant_id,receipt_id,authorization_id,organization_id,order_id,stock_unit_id,customer_principal_id,received_serial_number,condition_code,notes,evidence_sha256_hex,received_by_subject) values($1,'receipt-1','auth-1','store','original-order-1','original-1','customer','ORIGINAL-1','damaged','Received','eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee','operator'),($1,'receipt-2','auth-2','store','original-order-2','original-2','customer','ORIGINAL-2','damaged','Received','ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff','operator')`,
		`insert into sales.return_disposition(tenant_id,disposition_id,receipt_id,inventory_action,customer_remedy,notes,decided_by_subject) values($1,'disposition-1','receipt-1','quarantine','exchange','Replace','operator'),($1,'disposition-2','receipt-2','quarantine','exchange','Replace','operator')`,
	}
	for _, query := range fixtures {
		if _, err = tx.Exec(ctx, query, tenant); err != nil {
			t.Fatal(err)
		}
	}
	if err = tx.Commit(ctx); err != nil {
		t.Fatal(err)
	}

	for _, set := range []struct{ suffix string }{{"1"}, {"2"}} {
		if _, err = pool.Exec(ctx, `insert into sales.return_effect_request(tenant_id,request_id,disposition_id,effect_kind,owner_context,state,idempotency_key) values($1,$2,$3,'inventory','inventory','requested',$4),($1,$5,$3,'exchange','fulfillment','requested',$6),($1,$7,$3,'accounting','accounting','requested',$8)`, tenant, "inventory-"+set.suffix, "disposition-"+set.suffix, "inventory-request-key-000"+set.suffix, "exchange-"+set.suffix, "exchange-request-key-000"+set.suffix, "accounting-"+set.suffix, "accounting-request-key-000"+set.suffix); err != nil {
			t.Fatal(err)
		}
		if _, err = pool.Exec(ctx, `update sales.return_effect_execution set status='succeeded',result_sha256_hex=repeat('a',64) where tenant_id=$1 and request_id=$2`, tenant, "inventory-"+set.suffix); err != nil {
			t.Fatal(err)
		}
	}

	store := NewReturnExchange(pool)
	work, err := store.Claim(ctx, "fulfillment", "worker-a", "claim-a", time.Minute)
	if err != nil || work == nil || work.RequestID != "exchange-1" {
		t.Fatalf("work=%+v err=%v", work, err)
	}
	result, err := store.PrepareExchange(ctx, *work, "worker-a")
	if err != nil || result.ReplacementOrderID == "" || result.ReplacementStockID != "replacement-1" || len(result.ResultSHA256) != 64 {
		t.Fatalf("result=%+v err=%v", result, err)
	}
	var orderState, stockState, executionState, exchangeStatus string
	var total, stockVersion int64
	var handovers, events, attempts int
	if err = pool.QueryRow(ctx, `select state,total_minor_units from sales.customer_order where tenant_id=$1 and order_id=$2`, tenant, result.ReplacementOrderID).Scan(&orderState, &total); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select state,version from inventory.stock_unit where tenant_id=$1 and stock_unit_id='replacement-1'`, tenant).Scan(&stockState, &stockVersion); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select status from sales.return_effect_execution where tenant_id=$1 and request_id='exchange-1'`, tenant).Scan(&executionState); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select status from sales.return_exchange where tenant_id=$1 and request_id='exchange-1'`, tenant).Scan(&exchangeStatus); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select count(*) from sales.delivery_handover where tenant_id=$1 and handover_id=$2 and state='prepared'`, tenant, result.ReplacementHandoverID).Scan(&handovers); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select count(*) from platform.outbox_event where tenant_id=$1 and payload->>'request_id'='exchange-1'`, tenant).Scan(&events); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select count(*) from sales.return_effect_attempt where tenant_id=$1 and request_id='exchange-1' and outcome='succeeded'`, tenant).Scan(&attempts); err != nil {
		t.Fatal(err)
	}
	if orderState != "allocated" || total != 0 || stockState != "reserved" || stockVersion != 8 || executionState != "succeeded" || exchangeStatus != "prepared" || handovers != 1 || events != 3 || attempts != 1 {
		t.Fatalf("order=%s/%d stock=%s/%d execution=%s exchange=%s handovers=%d events=%d attempts=%d", orderState, total, stockState, stockVersion, executionState, exchangeStatus, handovers, events, attempts)
	}

	work, err = store.Claim(ctx, "fulfillment", "worker-b", "claim-b", time.Minute)
	if err != nil || work == nil || work.RequestID != "exchange-2" {
		t.Fatalf("second work=%+v err=%v", work, err)
	}
	if _, err = store.PrepareExchange(ctx, *work, "worker-b"); !errors.Is(err, returnexchange.ErrStockUnavailable) {
		t.Fatalf("no-stock err=%v", err)
	}
	if err = store.Finish(ctx, *work, "worker-b", returneffects.Completion{Outcome: "retry", ErrorCode: "REPLACEMENT_STOCK_UNAVAILABLE", RetryAfter: time.Minute}); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select status from sales.return_effect_execution where tenant_id=$1 and request_id='exchange-2'`, tenant).Scan(&executionState); err != nil {
		t.Fatal(err)
	}
	if executionState != "retry" {
		t.Fatalf("second execution=%s", executionState)
	}
	if _, err = pool.Exec(ctx, `update sales.return_exchange set status='prepared' where tenant_id=$1 and request_id='exchange-1'`, tenant); err == nil {
		t.Fatal("exchange mutation accepted")
	}
	if work, err = store.Claim(ctx, "fulfillment", "worker-c", "claim-c", time.Minute); err != nil || work != nil {
		t.Fatalf("terminal or delayed work reclaimed: %+v %v", work, err)
	}
}
````

### FILE: `cmd/return-exchange-worker/main.go`
```yaml
block_id: "GO-RETURN-EXCHANGE:file:08"
operation: CREATE
provenance: AUTHORED
source: "local implementation governed by official sources in metadata"
license: "LicenseRef-Workspace-Owner"
sha256: "f003a99e7e697d29f5333c952398246d21bec4462298e732d9afdfa41068612f"
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
	"elite.local/enterprise/internal/returnexchange"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	databaseURL, workerID := os.Getenv("DATABASE_URL"), os.Getenv("RETURN_EXCHANGE_WORKER_ID")
	if databaseURL == "" || workerID == "" {
		slog.Error("DATABASE_URL and RETURN_EXCHANGE_WORKER_ID are required")
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
	processor, err := returnexchange.NewProcessor(postgres.NewReturnExchange(pool), randomid.Generator{}, workerID, 2*time.Minute, 30*time.Second)
	if err != nil {
		slog.Error("return exchange worker configuration is invalid")
		os.Exit(2)
	}
	if err = run(ctx, processor); err != nil && !errors.Is(err, context.Canceled) {
		slog.Error("return exchange worker stopped unexpectedly")
		os.Exit(1)
	}
}

func run(ctx context.Context, processor *returnexchange.Processor) error {
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
		if errors.Is(err, returnexchange.ErrNoWork) {
			wait = 500 * time.Millisecond
		} else {
			slog.Warn("return exchange deferred")
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

Required runtime variables are `DATABASE_URL` by secret reference and a non-secret `RETURN_EXCHANGE_WORKER_ID`. The built-in lease is two minutes and stock retry delay is thirty seconds. Project admission must additionally record the allowed exchange variants/locations, financial settlement policy, accounting owner and delivery mode; this pack admits only exact same-variant local preparation.

## 7. Dependency bill

The root module already pins Go 1.26.7 and `github.com/jackc/pgx/v5 v5.10.0` under MIT. Runtime code adds no network SDK and invokes no carrier, payment or fiscal provider. Microsoft documentation is an architectural authority, not a copied runtime dependency.

## 8. Apply order

Apply after database migrations 0001 through 0020 and after `GO-RETURN-EFFECT-EXECUTION-WORKER`. Run migration 0021, build `cmd/return-exchange-worker`, provide database secret and worker identity, then enable only after the project's exchange and accounting decisions are approved.

## 9. Verification

Materialize all eight files and verify hashes. Run `go test ./...`, `go vet ./...` and build every root `cmd`, including `return-exchange-worker`. On PostgreSQL 18.6, apply 0001–0021, run the 0021 SQL test and the focused integration test. Prove one stock unit cannot satisfy two exchanges, unavailable stock becomes delayed retry, the prepared record is immutable, terminal work cannot be reclaimed, and down/up reapplication is clean.

## 10. Reconstruction evidence

Governed by `reconstruction_evidence/RETURN_EXCHANGE_FULFILLMENT_EXECUTION_INVENTORY_2026-08-31_V143.md`. No live customer, carrier, accounting ledger or fiscal provider was used; production readiness remains project-specific and blocked until those target gates pass.
