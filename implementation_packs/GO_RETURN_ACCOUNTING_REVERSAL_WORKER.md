# Go Return Accounting Reversal Worker

## 1. Metadata
```yaml
pack_id: "GO-RETURN-ACCOUNTING-REVERSAL-WORKER"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Consumes the accounting-owned return effect, reverses the unique posted SALE journal exactly line by line with the existing ledger, reconciles a crash after posting and records immutable execution evidence; it never invents accounts, tax or cost values."
stacks: ["Go 1.26.7", "PostgreSQL 18.6", "pgx 5.10.0"]
compatible_with: ["GO-ENTERPRISE-ACCOUNTING-LEDGER-API 0.1.x", "GO-RETURN-EFFECT-EXECUTION-WORKER 0.1.x", "GO-FRANCHISE-CUSTOMER-JOURNEY-API 0.8.x"]
incompatible_with: ["alternate accounting owner or mutable ledger without reconciliation contract"]
license_expression: "LicenseRef-Workspace-Owner AND MIT dependency"
upstream_sources: ["https://learn.microsoft.com/en-us/dynamics365/supply-chain/sales-marketing/sales-returns", "https://learn.microsoft.com/en-us/dynamics365/business-central/sales-how-process-sales-returns-cancellations", "https://learn.microsoft.com/en-us/dynamics365/business-central/finance-how-reverse-journal-posting", "https://github.com/microsoft/BCApps/tree/31a860b527f0dc72c7a44a255d7e7d403cfa4789", "https://www.postgresql.org/docs/18/explicit-locking.html", "https://docs.aws.amazon.com/wellarchitected/latest/reliability-pillar/welcome.html"]
verified_at: "2026-08-31"
```

All eight files are local `AUTHORED` implementation. Microsoft governs return credit/reversal, original-sale linkage and equal/opposite journal reversal; the pinned MIT BCApps evidence already governs the ledger implementation. PostgreSQL governs transactions and locking; AWS governs durable reconciliation. No local line is attributed to those organizations.

## 2. Applicability
Use only when the project has one posted `SALE` journal for the original order and has approved that journal's account, revenue, cost and tax composition. Inventory and the selected remedy (`refund` or `exchange`) must already have succeeded. Reject missing/ambiguous/incomplete source journals; never derive accounts from a document or order at runtime.

## 3. Architecture contract
The worker uses the existing accounting ledger as the single source of truth. It claims only `owner_context=accounting`, locates the exact original order journal and invokes its established equal/opposite reversal. Stable journal/register/event identities make a crash after ledger commit reconcilable: the next lease verifies the same reversal and completes the effect without reposting. The immutable link, attempt, execution result and return-accounting outbox event close atomically. Fiscal credit notes, statutory books and provider settlement remain separate owners.

## 4. Exact file manifest
```text
CREATE db/migrations/0022_return_accounting_reversal.up.sql
CREATE db/migrations/0022_return_accounting_reversal.down.sql
CREATE db/tests/0022_return_accounting_reversal.test.sql
CREATE internal/returnaccounting/processor.go
CREATE internal/returnaccounting/processor_test.go
CREATE internal/platform/postgres/returnaccounting.go
CREATE internal/platform/postgres/returnaccounting_integration_test.go
CREATE cmd/return-accounting-worker/main.go
```

## 5. Materialization blocks

### FILE: `db/migrations/0022_return_accounting_reversal.up.sql`
```yaml
block_id: "GO-RETURN-ACCOUNTING:file:01"
operation: CREATE
provenance: AUTHORED
source: "local implementation governed by official sources in metadata"
license: "LicenseRef-Workspace-Owner"
sha256: "6caab03910b27e5d7f58f4255078afef8c05e98d2a58eedae85f88b198b4aeee"
variables: []
secrets_allowed: false
```
````sql
begin;

create table accounting.return_effect_posting (
  tenant_id uuid not null,
  request_id text not null,
  disposition_id text not null,
  source_order_id text not null,
  remedy text not null check (remedy in ('refund','exchange')),
  original_journal_id text not null,
  reversal_journal_id text not null,
  currency text not null check (currency ~ '^[A-Z]{3}$'),
  reversed_minor_units bigint not null check (reversed_minor_units > 0),
  status text not null check (status='posted'),
  result_sha256_hex text not null check (result_sha256_hex ~ '^[0-9a-f]{64}$'),
  posted_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id,request_id),
  unique (tenant_id,disposition_id),
  unique (tenant_id,reversal_journal_id),
  foreign key (tenant_id,request_id) references sales.return_effect_request(tenant_id,request_id),
  foreign key (tenant_id,disposition_id) references sales.return_disposition(tenant_id,disposition_id),
  foreign key (tenant_id,source_order_id) references sales.customer_order(tenant_id,order_id),
  foreign key (tenant_id,original_journal_id) references accounting.journal(tenant_id,journal_id),
  foreign key (tenant_id,reversal_journal_id) references accounting.journal(tenant_id,journal_id),
  check (original_journal_id<>reversal_journal_id)
);

create or replace function accounting.prevent_return_effect_posting_mutation() returns trigger language plpgsql as $$
begin raise exception using errcode='23514',message='return accounting posting is immutable'; end $$;
create trigger return_effect_posting_immutable before update or delete on accounting.return_effect_posting for each row execute function accounting.prevent_return_effect_posting_mutation();
create index return_effect_posting_source_idx on accounting.return_effect_posting(tenant_id,source_order_id,posted_at desc,request_id);

commit;
````

### FILE: `db/migrations/0022_return_accounting_reversal.down.sql`
```yaml
block_id: "GO-RETURN-ACCOUNTING:file:02"
operation: CREATE
provenance: AUTHORED
source: "local implementation governed by official sources in metadata"
license: "LicenseRef-Workspace-Owner"
sha256: "ab31a81859825e1e47969eb761874583b8d97b5cdae5e40e3aebb5ae2e6c318c"
variables: []
secrets_allowed: false
```
````sql
begin;
drop index if exists accounting.return_effect_posting_source_idx;
drop trigger if exists return_effect_posting_immutable on accounting.return_effect_posting;
drop function if exists accounting.prevent_return_effect_posting_mutation();
drop table if exists accounting.return_effect_posting;
commit;
````

### FILE: `db/tests/0022_return_accounting_reversal.test.sql`
```yaml
block_id: "GO-RETURN-ACCOUNTING:file:03"
operation: CREATE
provenance: AUTHORED
source: "local implementation governed by official sources in metadata"
license: "LicenseRef-Workspace-Owner"
sha256: "905445ccf491ac8bef7d5d48fc6992067fd4ad60f47b9aeec655d4e8b35eb268"
variables: []
secrets_allowed: false
```
````sql
begin;
do $$ begin
  if to_regclass('accounting.return_effect_posting') is null then raise exception 'return accounting table missing'; end if;
  if (select count(*) from pg_trigger where tgname='return_effect_posting_immutable' and not tgisinternal)<>1 then raise exception 'return accounting immutability missing'; end if;
  if not exists(select 1 from pg_indexes where schemaname='accounting' and indexname='return_effect_posting_source_idx') then raise exception 'return accounting index missing'; end if;
end $$;
rollback;
````

### FILE: `internal/returnaccounting/processor.go`
```yaml
block_id: "GO-RETURN-ACCOUNTING:file:04"
operation: CREATE
provenance: AUTHORED
source: "local implementation governed by official sources in metadata"
license: "LicenseRef-Workspace-Owner"
sha256: "1b69ede07b06f06f3be7bf1840910726a86c0110fc4fe710d5c27fd9edb575c1"
variables: []
secrets_allowed: false
```
````go
package returnaccounting

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	"elite.local/enterprise/internal/returneffects"
)

var (
	ErrNoWork   = errors.New("no return accounting work")
	ErrConflict = errors.New("return accounting contract conflict")
)
var workerPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:-]{0,63}$`)

type Result struct {
	RequestID, OriginalJournalID, ReversalJournalID, ResultSHA256 string
	ReversedMinorUnits                                            int64
}
type Store interface {
	Claim(context.Context, string, string, string, time.Duration) (*returneffects.Work, error)
	PostOrReconcile(context.Context, returneffects.Work, string) (Result, error)
	Finish(context.Context, returneffects.Work, string, returneffects.Completion) error
}
type IDGenerator interface{ New() string }
type Processor struct {
	store        Store
	ids          IDGenerator
	workerID     string
	lease, retry time.Duration
}

func NewProcessor(store Store, ids IDGenerator, workerID string, lease, retry time.Duration) (*Processor, error) {
	if store == nil || ids == nil || !workerPattern.MatchString(workerID) || lease < time.Second || lease > 15*time.Minute || retry < time.Second || retry > time.Hour {
		return nil, fmt.Errorf("invalid return accounting processor configuration")
	}
	return &Processor{store: store, ids: ids, workerID: workerID, lease: lease, retry: retry}, nil
}
func (p *Processor) ProcessOne(ctx context.Context) (Result, error) {
	work, err := p.store.Claim(ctx, "accounting", p.workerID, p.ids.New(), p.lease)
	if err != nil {
		return Result{}, err
	}
	if work == nil {
		return Result{}, ErrNoWork
	}
	result, postErr := p.store.PostOrReconcile(ctx, *work, p.workerID)
	if postErr == nil {
		return result, nil
	}
	completion := returneffects.Completion{Outcome: "retry", ErrorCode: "ACCOUNTING_TRANSIENT_FAILURE", RetryAfter: p.retry}
	if errors.Is(postErr, ErrConflict) {
		completion = returneffects.Completion{Outcome: "blocked", ErrorCode: "ACCOUNTING_SOURCE_CONFLICT"}
	}
	if finishErr := p.store.Finish(ctx, *work, p.workerID, completion); finishErr != nil {
		return Result{}, errors.Join(postErr, finishErr)
	}
	return Result{}, postErr
}
````

### FILE: `internal/returnaccounting/processor_test.go`
```yaml
block_id: "GO-RETURN-ACCOUNTING:file:05"
operation: CREATE
provenance: AUTHORED
source: "local implementation governed by official sources in metadata"
license: "LicenseRef-Workspace-Owner"
sha256: "f9eabd178d60fd0196a96a198c8adde53e6c34ddea1df1ac75ff93b6f6680487"
variables: []
secrets_allowed: false
```
````go
package returnaccounting

import (
	"context"
	"elite.local/enterprise/internal/returneffects"
	"errors"
	"testing"
	"time"
)

type ids struct{ v string }

func (i ids) New() string { return i.v }

type store struct {
	work       *returneffects.Work
	err        error
	completion returneffects.Completion
}

func (s *store) Claim(context.Context, string, string, string, time.Duration) (*returneffects.Work, error) {
	return s.work, nil
}
func (s *store) PostOrReconcile(context.Context, returneffects.Work, string) (Result, error) {
	return Result{RequestID: "a", ReversalJournalID: "r"}, s.err
}
func (s *store) Finish(_ context.Context, _ returneffects.Work, _ string, c returneffects.Completion) error {
	s.completion = c
	return nil
}
func TestProcessor(t *testing.T) {
	s := &store{work: &returneffects.Work{RequestID: "a"}}
	p, e := NewProcessor(s, ids{"c"}, "worker", time.Minute, time.Second)
	if e != nil {
		t.Fatal(e)
	}
	r, e := p.ProcessOne(context.Background())
	if e != nil || r.ReversalJournalID != "r" {
		t.Fatalf("%+v %v", r, e)
	}
	s.err = ErrConflict
	if _, e = p.ProcessOne(context.Background()); !errors.Is(e, ErrConflict) || s.completion.Outcome != "blocked" {
		t.Fatalf("%v %+v", e, s.completion)
	}
}
func TestNoWorkAndConfig(t *testing.T) {
	if _, e := NewProcessor(&store{}, ids{"c"}, "bad worker", time.Minute, time.Second); e == nil {
		t.Fatal("invalid accepted")
	}
	p, e := NewProcessor(&store{}, ids{"c"}, "worker", time.Minute, time.Second)
	if e != nil {
		t.Fatal(e)
	}
	if _, e = p.ProcessOne(context.Background()); !errors.Is(e, ErrNoWork) {
		t.Fatal(e)
	}
}
````

### FILE: `internal/platform/postgres/returnaccounting.go`
```yaml
block_id: "GO-RETURN-ACCOUNTING:file:06"
operation: CREATE
provenance: AUTHORED
source: "local implementation governed by official sources in metadata"
license: "LicenseRef-Workspace-Owner"
sha256: "0fadb731322c9bef409165bc212739b3fb4c44341fe21820aa0a5644ee42123c"
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

	"elite.local/enterprise/internal/accounting"
	"elite.local/enterprise/internal/returnaccounting"
	"elite.local/enterprise/internal/returneffects"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReturnAccounting struct {
	pool    *pgxpool.Pool
	effects *ReturnEffects
	ledger  *Accounting
}

func NewReturnAccounting(pool *pgxpool.Pool) *ReturnAccounting {
	return &ReturnAccounting{pool: pool, effects: NewReturnEffects(pool), ledger: NewAccounting(pool)}
}
func (r *ReturnAccounting) Claim(ctx context.Context, owner, worker, token string, lease time.Duration) (*returneffects.Work, error) {
	return r.effects.Claim(ctx, owner, worker, token, lease)
}
func (r *ReturnAccounting) Finish(ctx context.Context, work returneffects.Work, worker string, c returneffects.Completion) error {
	return r.effects.Finish(ctx, work, worker, c)
}

func stableReturnAccountingUUID(scope, value string) string {
	s := sha256.Sum256([]byte(scope + "\x00" + value))
	b := s[:16]
	b[6] = (b[6] & 0x0f) | 0x50
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

type accountingSource struct {
	disposition, order, organization, remedy, journal, status, currency string
	version, total                                                      int64
}

func (r *ReturnAccounting) PostOrReconcile(ctx context.Context, work returneffects.Work, worker string) (returnaccounting.Result, error) {
	if work.EffectKind != "accounting" || work.OwnerContext != "accounting" {
		return returnaccounting.Result{}, returnaccounting.ErrConflict
	}
	var s accountingSource
	err := r.pool.QueryRow(ctx, `
select e.disposition_id,rr.order_id,rr.organization_id,d.customer_remedy,j.journal_id,j.status,j.version,j.currency,j.total_debit_minor_units
from sales.return_effect_execution x
join sales.return_effect_request e on e.tenant_id=x.tenant_id and e.request_id=x.request_id
join sales.return_disposition d on d.tenant_id=e.tenant_id and d.disposition_id=e.disposition_id
join sales.return_receipt rr on rr.tenant_id=d.tenant_id and rr.receipt_id=d.receipt_id
join accounting.journal j on j.tenant_id=rr.tenant_id and j.organization_id=rr.organization_id and j.source_type='SALE' and j.source_id=rr.order_id and j.reversal_of is null
join sales.return_effect_request inv on inv.tenant_id=e.tenant_id and inv.disposition_id=e.disposition_id and inv.effect_kind='inventory'
join sales.return_effect_execution invx on invx.tenant_id=inv.tenant_id and invx.request_id=inv.request_id and invx.status='succeeded'
join sales.return_effect_request remedy on remedy.tenant_id=e.tenant_id and remedy.disposition_id=e.disposition_id and remedy.effect_kind=d.customer_remedy
join sales.return_effect_execution remedyx on remedyx.tenant_id=remedy.tenant_id and remedyx.request_id=remedy.request_id and remedyx.status='succeeded'
where x.tenant_id=$1 and x.request_id=$2 and x.status='claimed' and x.claimed_by=$3 and x.claim_token=$4 and x.claimed_until>=clock_timestamp()
and e.effect_kind='accounting' and e.owner_context='accounting' and j.status in('posted','reversed')`, work.TenantID, work.RequestID, worker, work.ClaimToken).Scan(&s.disposition, &s.order, &s.organization, &s.remedy, &s.journal, &s.status, &s.version, &s.currency, &s.total)
	if errors.Is(err, pgx.ErrNoRows) {
		return returnaccounting.Result{}, returnaccounting.ErrConflict
	}
	if err != nil {
		return returnaccounting.Result{}, err
	}
	reversalID := stableReturnAccountingUUID("journal", work.TenantID+"/"+work.RequestID)
	registerID := stableReturnAccountingUUID("register", work.TenantID+"/"+work.RequestID)
	eventID := stableReturnAccountingUUID("event", work.TenantID+"/"+work.RequestID)
	if s.status == "posted" {
		_, err = r.ledger.ReverseJournal(ctx, work.TenantID, s.organization, s.journal, reversalID, registerID, eventID, s.version, worker, "approved customer return "+work.RequestID)
		if err != nil && !errors.Is(err, accounting.ErrConflict) {
			return returnaccounting.Result{}, err
		}
	}
	var reversalStatus string
	var reversedTotal int64
	var exactLines int
	err = r.pool.QueryRow(ctx, `select r.status,r.total_debit_minor_units,count(*)::int from accounting.journal r join accounting.entry re on re.tenant_id=r.tenant_id and re.journal_id=r.journal_id join accounting.entry oe on oe.tenant_id=re.tenant_id and oe.entry_id=re.original_entry_id where r.tenant_id=$1 and r.journal_id=$2 and r.reversal_of=$3 and r.source_type='REVERSAL' and r.source_id=$3 and re.debit_minor_units=oe.credit_minor_units and re.credit_minor_units=oe.debit_minor_units group by r.status,r.total_debit_minor_units`, work.TenantID, reversalID, s.journal).Scan(&reversalStatus, &reversedTotal, &exactLines)
	if errors.Is(err, pgx.ErrNoRows) || reversalStatus != "posted" || reversedTotal != s.total || exactLines < 2 {
		return returnaccounting.Result{}, returnaccounting.ErrConflict
	}
	if err != nil {
		return returnaccounting.Result{}, err
	}

	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return returnaccounting.Result{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var attempt int
	var started time.Time
	err = tx.QueryRow(ctx, `select attempt_count,claimed_at from sales.return_effect_execution where tenant_id=$1 and request_id=$2 and status='claimed' and claimed_by=$3 and claim_token=$4 and claimed_until>=clock_timestamp() for update`, work.TenantID, work.RequestID, worker, work.ClaimToken).Scan(&attempt, &started)
	if errors.Is(err, pgx.ErrNoRows) || attempt != work.Attempt {
		return returnaccounting.Result{}, returnaccounting.ErrConflict
	}
	if err != nil {
		return returnaccounting.Result{}, err
	}
	command, err := json.Marshal(struct {
		RequestID, OrderID, OriginalJournalID, ReversalJournalID, Remedy, Currency string
		ReversedMinorUnits                                                         int64
	}{work.RequestID, s.order, s.journal, reversalID, s.remedy, s.currency, reversedTotal})
	if err != nil {
		return returnaccounting.Result{}, err
	}
	digest := sha256.Sum256(command)
	hash := hex.EncodeToString(digest[:])
	if _, err = tx.Exec(ctx, `insert into accounting.return_effect_posting(tenant_id,request_id,disposition_id,source_order_id,remedy,original_journal_id,reversal_journal_id,currency,reversed_minor_units,status,result_sha256_hex)values($1,$2,$3,$4,$5,$6,$7,$8,$9,'posted',$10)`, work.TenantID, work.RequestID, s.disposition, s.order, s.remedy, s.journal, reversalID, s.currency, reversedTotal, hash); err != nil {
		return returnaccounting.Result{}, err
	}
	if _, err = tx.Exec(ctx, `insert into sales.return_effect_attempt(tenant_id,attempt_id,request_id,attempt_no,worker_id,outcome,provider_reference,result_sha256_hex,started_at)values($1,$2,$3,$4,$5,'succeeded',$6,$7,$8)`, work.TenantID, work.ClaimToken, work.RequestID, attempt, worker, reversalID, hash, started); err != nil {
		return returnaccounting.Result{}, err
	}
	updated, err := tx.Exec(ctx, `update sales.return_effect_execution set status='succeeded',claimed_at=null,claimed_by=null,claim_token=null,claimed_until=null,last_error_code=null,provider_reference=$5,result_sha256_hex=$6,updated_at=clock_timestamp() where tenant_id=$1 and request_id=$2 and status='claimed' and claimed_by=$3 and claim_token=$4`, work.TenantID, work.RequestID, worker, work.ClaimToken, reversalID, hash)
	if err != nil {
		return returnaccounting.Result{}, err
	}
	if updated.RowsAffected() != 1 {
		return returnaccounting.Result{}, returnaccounting.ErrConflict
	}
	if _, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)values($1,gen_random_uuid(),'return-accounting',$2,1,'return.accounting-posted',1,clock_timestamp(),$3)`, work.TenantID, work.RequestID, command); err != nil {
		return returnaccounting.Result{}, err
	}
	if err = tx.Commit(ctx); err != nil {
		return returnaccounting.Result{}, err
	}
	return returnaccounting.Result{RequestID: work.RequestID, OriginalJournalID: s.journal, ReversalJournalID: reversalID, ReversedMinorUnits: reversedTotal, ResultSHA256: hash}, nil
}
````

### FILE: `internal/platform/postgres/returnaccounting_integration_test.go`
```yaml
block_id: "GO-RETURN-ACCOUNTING:file:07"
operation: CREATE
provenance: AUTHORED
source: "local implementation governed by official sources in metadata"
license: "LicenseRef-Workspace-Owner"
sha256: "8ef1458b6930e9bcdb23451406d49bd0d7329d3581363cae91d8a238563da142"
variables: []
secrets_allowed: false
```
````go
package postgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"testing"
	"time"
)

func TestReturnAccountingExactReversalAndCrashReconciliation(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	ctx := context.Background()
	pool, e := pgxpool.New(ctx, url)
	if e != nil {
		t.Fatal(e)
	}
	defer pool.Close()
	tenant := "018f4d4a-7b36-7a21-8d10-2f4c54c2e144"
	cleanup := func() {
		tx, err := pool.Begin(ctx)
		if err != nil {
			return
		}
		defer tx.Rollback(ctx)
		_, _ = tx.Exec(ctx, `set local session_replication_role=replica`)
		for _, table := range []string{"accounting.return_effect_posting", "sales.return_effect_resume", "sales.return_effect_attempt", "sales.return_effect_execution", "sales.return_effect_request", "sales.return_disposition", "sales.return_receipt", "sales.return_authorization", "sales.delivery_exception", "sales.delivery_handover", "sales.customer_order_line", "sales.customer_order", "platform.outbox_event", "accounting.entry", "accounting.register", "accounting.journal_line", "accounting.journal", "accounting.period", "accounting.account", "inventory.stock_unit", "crm.customer_profile", "catalog.vehicle_variant", "catalog.vehicle_model", "org.organization", "platform.tenant"} {
			_, _ = tx.Exec(ctx, `delete from `+table+` where tenant_id=$1`, tenant)
		}
		_ = tx.Commit(ctx)
	}
	cleanup()
	defer cleanup()
	tx, e := pool.Begin(ctx)
	if e != nil {
		t.Fatal(e)
	}
	defer tx.Rollback(ctx)
	_, e = tx.Exec(ctx, `set local session_replication_role=replica`)
	if e != nil {
		t.Fatal(e)
	}
	queries := []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1,'return-accounting','Return Accounting','Return Accounting')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'store','return-accounting-store','Store','store')`,
		`insert into catalog.vehicle_model(tenant_id,model_id,model_code,display_name,vehicle_class,lifecycle_state)values($1,'model','ra-model','Model','motorcycle','active')`,
		`insert into catalog.vehicle_variant(tenant_id,variant_id,model_id,variant_code,display_name,battery_specification,lifecycle_state)values($1,'variant','model','ra-variant','Variant','{}','active')`,
		`insert into crm.customer_profile(tenant_id,customer_principal_id,display_name)values($1,'customer','Customer')`,
		`insert into accounting.account(tenant_id,account_code,display_name,account_type)values($1,'CASH','Cash','asset'),($1,'REVENUE','Revenue','revenue')`,
		`insert into accounting.period(tenant_id,period_id,starts_on,ends_on,status,version)values($1,'2026-08','2026-08-01','2026-09-01','open',1)`,
		`insert into inventory.stock_unit(tenant_id,stock_unit_id,organization_id,variant_id,serial_number,vin,battery_serial_number,state,version,received_at)values($1,'stock-1','store','variant','SERIAL-1','VIN-1','BAT-1','quarantine',2,clock_timestamp()),($1,'stock-2','store','variant','SERIAL-2','VIN-2','BAT-2','quarantine',2,clock_timestamp())`,
		`insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version)values($1,'order-1','store','customer','delivered','ARS',10000,4),($1,'order-2','store','customer','delivered','ARS',20000,4)`,
		`insert into sales.customer_order_line(tenant_id,order_id,line_id,variant_id,quantity,unit_price_minor_units,allocated_stock_unit_id)values($1,'order-1','line-1','variant',1,10000,'stock-1'),($1,'order-2','line-2','variant',1,20000,'stock-2')`,
		`insert into sales.delivery_handover(tenant_id,handover_id,organization_id,order_id,customer_principal_id,stock_unit_id,state,acceptance_evidence_sha256_hex,customer_accepted_at,version)values($1,'hand-1','store','order-1','customer','stock-1','accepted',repeat('a',64),clock_timestamp(),2),($1,'hand-2','store','order-2','customer','stock-2','accepted',repeat('b',64),clock_timestamp(),2)`,
		`insert into sales.delivery_exception(tenant_id,exception_id,organization_id,handover_id,customer_principal_id,reason_code,details,rejection_evidence_sha256_hex,state,version,resolved_at,resolved_by_subject)values($1,'ex-1','store','hand-1','customer','defect','Defect',repeat('c',64),'resolved',2,clock_timestamp(),'operator'),($1,'ex-2','store','hand-2','customer','defect','Defect',repeat('d',64),'resolved',2,clock_timestamp(),'operator')`,
		`insert into sales.return_authorization(tenant_id,authorization_id,organization_id,exception_id,handover_id,order_id,stock_unit_id,customer_principal_id,disposition,exact_cost_source_order_id,state,authorized_by_subject)values($1,'auth-1','store','ex-1','hand-1','order-1','stock-1','customer','return','order-1','authorized','operator'),($1,'auth-2','store','ex-2','hand-2','order-2','stock-2','customer','return','order-2','authorized','operator')`,
		`insert into sales.return_receipt(tenant_id,receipt_id,authorization_id,organization_id,order_id,stock_unit_id,customer_principal_id,received_serial_number,condition_code,notes,evidence_sha256_hex,received_by_subject)values($1,'receipt-1','auth-1','store','order-1','stock-1','customer','SERIAL-1','damaged','Received',repeat('e',64),'operator'),($1,'receipt-2','auth-2','store','order-2','stock-2','customer','SERIAL-2','damaged','Received',repeat('f',64),'operator')`,
		`insert into sales.return_disposition(tenant_id,disposition_id,receipt_id,inventory_action,customer_remedy,notes,decided_by_subject)values($1,'disp-1','receipt-1','quarantine','refund','Return','operator'),($1,'disp-2','receipt-2','quarantine','refund','Return','operator')`,
	}
	for _, q := range queries {
		if _, e = tx.Exec(ctx, q, tenant); e != nil {
			t.Fatal(e)
		}
	}
	for _, v := range []struct {
		id     string
		amount int64
	}{{"1", 10000}, {"2", 20000}} {
		if _, e = tx.Exec(ctx, `insert into accounting.journal(tenant_id,journal_id,organization_id,period_id,source_type,source_id,currency,posting_date,status,total_debit_minor_units,total_credit_minor_units,version,posted_by,posted_at)values($1,$2,'store','2026-08','SALE',$3,'ARS','2026-08-15','posted',$4,$4,2,'controller',clock_timestamp())`, tenant, "sale-"+v.id, "order-"+v.id, v.amount); e != nil {
			t.Fatal(e)
		}
		if _, e = tx.Exec(ctx, `insert into accounting.journal_line(tenant_id,journal_id,line_no,account_code,description,debit_minor_units,credit_minor_units)values($1,$2,1,'CASH','Sale cash',$3,0),($1,$2,2,'REVENUE','Sale revenue',0,$3)`, tenant, "sale-"+v.id, v.amount); e != nil {
			t.Fatal(e)
		}
		if _, e = tx.Exec(ctx, `insert into accounting.register(tenant_id,register_id,journal_id,organization_id,period_id,posted_by,total_debit_minor_units,total_credit_minor_units)values($1,$2,$3,'store','2026-08','controller',$4,$4)`, tenant, "sale-register-"+v.id, "sale-"+v.id, v.amount); e != nil {
			t.Fatal(e)
		}
		if _, e = tx.Exec(ctx, `insert into accounting.entry(tenant_id,entry_id,register_id,journal_id,line_no,organization_id,period_id,account_code,currency,posting_date,description,debit_minor_units,credit_minor_units)values($1,$2,$3,$4,1,'store','2026-08','CASH','ARS','2026-08-15','Sale cash',$5,0),($1,$6,$3,$4,2,'store','2026-08','REVENUE','ARS','2026-08-15','Sale revenue',0,$5)`, tenant, "sale-entry-"+v.id+"-1", "sale-register-"+v.id, "sale-"+v.id, v.amount, "sale-entry-"+v.id+"-2"); e != nil {
			t.Fatal(e)
		}
	}
	if e = tx.Commit(ctx); e != nil {
		t.Fatal(e)
	}
	for _, id := range []string{"1", "2"} {
		if _, e = pool.Exec(ctx, `insert into sales.return_effect_request(tenant_id,request_id,disposition_id,effect_kind,owner_context,state,idempotency_key)values($1,$2,$3,'inventory','inventory','requested',$4),($1,$5,$3,'refund','payment','requested',$6),($1,$7,$3,'accounting','accounting','requested',$8)`, tenant, "inventory-"+id, "disp-"+id, "inventory-key-000000"+id, "refund-"+id, "refund-key-000000000"+id, "accounting-"+id, "accounting-key-000000"+id); e != nil {
			t.Fatal(e)
		}
		if _, e = pool.Exec(ctx, `update sales.return_effect_execution set status='succeeded',result_sha256_hex=repeat('a',64) where tenant_id=$1 and request_id in($2,$3)`, tenant, "inventory-"+id, "refund-"+id); e != nil {
			t.Fatal(e)
		}
	}
	store := NewReturnAccounting(pool)
	work, e := store.Claim(ctx, "accounting", "worker-a", "claim-a", time.Minute)
	if e != nil || work == nil || work.RequestID != "accounting-1" {
		t.Fatalf("%+v %v", work, e)
	}
	result, e := store.PostOrReconcile(ctx, *work, "worker-a")
	if e != nil || result.ReversedMinorUnits != 10000 {
		t.Fatalf("%+v %v", result, e)
	}
	assertAccountingReturn(t, ctx, pool, tenant, "accounting-1", "sale-1", result.ReversalJournalID, 10000)
	reversalID := stableReturnAccountingUUID("journal", tenant+"/accounting-2")
	_, e = NewAccounting(pool).ReverseJournal(ctx, tenant, "store", "sale-2", reversalID, stableReturnAccountingUUID("register", tenant+"/accounting-2"), stableReturnAccountingUUID("event", tenant+"/accounting-2"), 2, "crashed-worker", "approved customer return accounting-2")
	if e != nil {
		t.Fatal(e)
	}
	work, e = store.Claim(ctx, "accounting", "worker-b", "claim-b", time.Minute)
	if e != nil || work == nil || work.RequestID != "accounting-2" {
		t.Fatalf("%+v %v", work, e)
	}
	result, e = store.PostOrReconcile(ctx, *work, "worker-b")
	if e != nil || result.ReversalJournalID != reversalID {
		t.Fatalf("%+v %v", result, e)
	}
	assertAccountingReturn(t, ctx, pool, tenant, "accounting-2", "sale-2", reversalID, 20000)
	if _, e = pool.Exec(ctx, `update accounting.return_effect_posting set status='posted' where tenant_id=$1 and request_id='accounting-1'`, tenant); e == nil {
		t.Fatal("posting mutation accepted")
	}
}
func assertAccountingReturn(t *testing.T, ctx context.Context, pool *pgxpool.Pool, tenant, request, original, reversal string, amount int64) {
	t.Helper()
	var os, rs, xs string
	var posting, attempts, swapped int
	if e := pool.QueryRow(ctx, `select status from accounting.journal where tenant_id=$1 and journal_id=$2`, tenant, original).Scan(&os); e != nil {
		t.Fatal(e)
	}
	if e := pool.QueryRow(ctx, `select status from accounting.journal where tenant_id=$1 and journal_id=$2`, tenant, reversal).Scan(&rs); e != nil {
		t.Fatal(e)
	}
	if e := pool.QueryRow(ctx, `select status from sales.return_effect_execution where tenant_id=$1 and request_id=$2`, tenant, request).Scan(&xs); e != nil {
		t.Fatal(e)
	}
	if e := pool.QueryRow(ctx, `select count(*) from accounting.return_effect_posting where tenant_id=$1 and request_id=$2 and reversed_minor_units=$3`, tenant, request, amount).Scan(&posting); e != nil {
		t.Fatal(e)
	}
	if e := pool.QueryRow(ctx, `select count(*) from sales.return_effect_attempt where tenant_id=$1 and request_id=$2 and outcome='succeeded'`, tenant, request).Scan(&attempts); e != nil {
		t.Fatal(e)
	}
	if e := pool.QueryRow(ctx, `select count(*) from accounting.entry r join accounting.entry o on o.tenant_id=r.tenant_id and o.entry_id=r.original_entry_id where r.tenant_id=$1 and r.journal_id=$2 and r.debit_minor_units=o.credit_minor_units and r.credit_minor_units=o.debit_minor_units`, tenant, reversal).Scan(&swapped); e != nil {
		t.Fatal(e)
	}
	if os != "reversed" || rs != "posted" || xs != "succeeded" || posting != 1 || attempts != 1 || swapped != 2 {
		t.Fatalf("original=%s reversal=%s execution=%s posting=%d attempts=%d swapped=%d", os, rs, xs, posting, attempts, swapped)
	}
}
````

### FILE: `cmd/return-accounting-worker/main.go`
```yaml
block_id: "GO-RETURN-ACCOUNTING:file:08"
operation: CREATE
provenance: AUTHORED
source: "local implementation governed by official sources in metadata"
license: "LicenseRef-Workspace-Owner"
sha256: "c26592a9caad543cd9ad6e3c7506e2f748460309c3ae2bd0c85bb218ae7a5e3f"
variables: []
secrets_allowed: false
```
````go
package main

import (
	"context"
	"elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/platform/randomid"
	"elite.local/enterprise/internal/returnaccounting"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	url, id := os.Getenv("DATABASE_URL"), os.Getenv("RETURN_ACCOUNTING_WORKER_ID")
	if url == "" || id == "" {
		slog.Error("DATABASE_URL and RETURN_ACCOUNTING_WORKER_ID are required")
		os.Exit(2)
	}
	cfg, e := pgxpool.ParseConfig(url)
	if e != nil {
		os.Exit(2)
	}
	cfg.MaxConns = 4
	p, e := pgxpool.NewWithConfig(ctx, cfg)
	if e != nil {
		os.Exit(1)
	}
	defer p.Close()
	if e = p.Ping(ctx); e != nil {
		os.Exit(1)
	}
	processor, e := returnaccounting.NewProcessor(postgres.NewReturnAccounting(p), randomid.Generator{}, id, 2*time.Minute, 30*time.Second)
	if e != nil {
		os.Exit(2)
	}
	if e = run(ctx, processor); e != nil && !errors.Is(e, context.Canceled) {
		os.Exit(1)
	}
}
func run(ctx context.Context, p *returnaccounting.Processor) error {
	timer := time.NewTimer(0)
	if !timer.Stop() {
		<-timer.C
	}
	defer timer.Stop()
	for {
		_, e := p.ProcessOne(ctx)
		if e == nil {
			continue
		}
		if errors.Is(e, context.Canceled) {
			return e
		}
		wait := 2 * time.Second
		if errors.Is(e, returnaccounting.ErrNoWork) {
			wait = 500 * time.Millisecond
		} else {
			slog.Warn("return accounting deferred")
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
Runtime requires `DATABASE_URL` by secret reference and `RETURN_ACCOUNTING_WORKER_ID`. Project readiness must approve chart of accounts, source-journal completeness, posting period, cost/tax treatment and operator responsibility before enabling this worker.

## 7. Dependency bill
Uses the root Go 1.26.7 module and `github.com/jackc/pgx/v5 v5.10.0` under MIT. It adds no provider SDK and reuses the existing accounting reversal implementation governed by pinned Microsoft BCApps evidence.

## 8. Apply order
Apply after migrations 0001–0021, the accounting ledger and return-effect worker. Run migration 0022, prove source-journal policy, then enable the binary only after inventory and remedy owners are active.

## 9. Verification
Materialize eight files; verify SHA; run full Go test/vet and every root build. On PostgreSQL 18.6 apply 0001–0022, SQL test and integration. Prove exact swapped entries, one reversal, terminal evidence, immutability, crash reconciliation and down/up.

## 10. Reconstruction evidence
Governed by `reconstruction_evidence/RETURN_ACCOUNTING_REVERSAL_EXECUTION_INVENTORY_2026-08-31_V144.md`. No accountant approval, live statutory ledger, ARCA credit note or audited close was executed.
