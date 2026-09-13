# Go Enterprise Accounting Ledger API

## 1. Metadata

```yaml
pack_id: "GO-ENTERPRISE-ACCOUNTING-LEDGER-API"
pack_version: "0.1.2"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Añade chart of accounts, períodos, journals balanceados, posting concurrente, ledger inmutable, trial balance, reversión compensatoria y cierre."
stacks: ["Go 1.26.7", "PostgreSQL 18.6", "pgx 5.10.0"]
compatible_with: ["GO-ENTERPRISE-BACKEND 0.4.x", "GO-FRANCHISE-ROYALTY-SETTLEMENT-API 0.1.x", "GO-ELECTROMOBILITY-APPLICATION 1.x", "GO-BC-EXACT-AMOUNT-ADAPTER 0.1.0"]
incompatible_with: []
license_expression: "LicenseRef-Workspace-Owner AND MIT upstream evidence"
upstream_sources: ["https://github.com/microsoft/BCApps/tree/31a860b527f0dc72c7a44a255d7e7d403cfa4789", "https://www.postgresql.org/docs/18/", "https://go.dev"]
verified_at: "2026-08-30"
```

Los nueve bloques son `AUTHORED`; no se atribuye Go/SQL local a Microsoft. La autoridad exacta es `microsoft/BCApps@31a860b527f0dc72c7a44a255d7e7d403cfa4789` (ZIP `e3151b040df39cada83d41fcf5d9e6cdff1d8fddf226934c1ad21a2bdebc7210`, MIT raíz `c2cfccb812fe482101a8f04597dfc5a9991a6b2748266c47ac91b6a5aae15383`): `GenJournalLine.Table.al` `8fb790053a04517d02f425f45e0d02dacecaf8b5ef2bd2bc23fa6b398342769b`, `GLEntry.Table.al` `08655960b539eb5afbea0c91f6189af07b4f83c68f7de08d5d590da6591fd9dd`, `GLRegister.Table.al` `2f7a870f2868a39c53d89b97ed831741579afef506458a5d701bd8bf6c6c807d`, `AccountingPeriod.Table.al` `d0b2628475e1abb64208139256ee3ec9f7d45e81564ab94c250564b2779e1ed6`, `VATEntry.Table.al` `c7aff50430b505d30ebf2ad6b3ee7f0a213510a2bf32828f5ffcf54a6d8d1fe4`, `GeneralPostingSetup.Table.al` `033fac045efae2df61be9d1d7ab22006acbc9f53c87a828bf74d0e96388ac580`, `ReversalEntry.Table.al` `13e17b83c3a2d685060ac151a9f5ef0b156238036f034515ccb260d0b2eb835d` y `GenJnlPostLine.Codeunit.al` `fdecc9a5e52831552e9addfe4c7e66d1b53fa3425c95b2d49d606478a24d54f3`.

V402 / 0.1.1: calls GO-BC-EXACT-AMOUNT-ADAPTER0.1.0 for the precisely mapped amount/balance/reversal functions. Existing caller blocks remain AUTHORED; other business semantics are not reclassified. Source derivation, FAIL808 and exact local tests: reconstruction_evidence/BC_EXACT_AMOUNT_ADAPTATION_V402.md. Composition now requires that adapter.

## 2. Applicability

Use as the operational general-ledger boundary after tenant/organization foundation. Upstream business owners map approved events into balanced journals; this pack never infers accounts from arbitrary documents. Reject direct browser exposure without accounting permissions. It is not an Argentine localization, tax engine, statutory invoice system, bank feed or audited financial statement product.

## 3. Architecture contract

One account and period owner, one source identity per journal, positive debit-or-credit lines and exact minor-unit balance. PostgreSQL verifies open period, posting date, active accounts, organization and optimistic version. Posting creates register and immutable entries atomically with outbox; serializable concurrency admits one winner. Reversal swaps debit/credit into a new posted journal and links every original entry. A period closes only with no draft journals; closed periods reject new journals. Trial balance is tenant/organization/period scoped.

## 4. Exact file manifest

```text
CREATE internal/accounting/service.go
CREATE internal/accounting/service_test.go
CREATE internal/platform/postgres/accounting.go
CREATE internal/platform/postgres/accounting_integration_test.go
CREATE internal/platform/httpapi/accounting.go
CREATE internal/platform/httpapi/accounting_test.go
CREATE db/migrations/0011_accounting_ledger.up.sql
CREATE db/migrations/0011_accounting_ledger.down.sql
CREATE db/tests/0011_accounting_ledger.test.sql
```

## 5. Materialization blocks

### FILE: `internal/accounting/service.go`
```yaml
block_id: "GO-ACCOUNTING:service:v1"
operation: CREATE
provenance: AUTHORED
source: "local implementation governed by pinned Microsoft BCApps evidence"
license: "LicenseRef-Workspace-Owner"
sha256: "97a1d5bcc5691333b82f76eb1dac27b732c9f0514dbc6dc626d15c6a2eefacee"
variables: []
secrets_allowed: false
```
````go
package accounting

import (
	"context"
	"elite.local/enterprise/internal/bcamounts"
	"errors"
	"regexp"
	"time"
)

var ErrInvalid = errors.New("invalid accounting command")
var ErrConflict = errors.New("accounting conflict")
var codePattern = regexp.MustCompile(`^[A-Z0-9][A-Z0-9._-]{0,31}$`)
var currencyPattern = regexp.MustCompile(`^[A-Z]{3}$`)

type Account struct {
	Code string `json:"code"`
	Name string `json:"name"`
	Type string `json:"type"`
}
type Period struct {
	ID       string    `json:"id"`
	StartsOn time.Time `json:"starts_on"`
	EndsOn   time.Time `json:"ends_on"`
	Status   string    `json:"status"`
	Version  int64     `json:"version"`
}
type Line struct {
	LineNo           int    `json:"line_no"`
	AccountCode      string `json:"account_code"`
	Description      string `json:"description"`
	DebitMinorUnits  int64  `json:"debit_minor_units"`
	CreditMinorUnits int64  `json:"credit_minor_units"`
}
type Journal struct {
	ID                    string    `json:"id"`
	OrganizationID        string    `json:"organization_id"`
	PeriodID              string    `json:"period_id"`
	SourceType            string    `json:"source_type"`
	SourceID              string    `json:"source_id"`
	Currency              string    `json:"currency"`
	PostingDate           time.Time `json:"posting_date"`
	Status                string    `json:"status"`
	TotalDebitMinorUnits  int64     `json:"total_debit_minor_units"`
	TotalCreditMinorUnits int64     `json:"total_credit_minor_units"`
	Version               int64     `json:"version"`
	ReversalOf            string    `json:"reversal_of,omitempty"`
	Lines                 []Line    `json:"lines,omitempty"`
}
type Balance struct {
	AccountCode      string `json:"account_code"`
	DebitMinorUnits  int64  `json:"debit_minor_units"`
	CreditMinorUnits int64  `json:"credit_minor_units"`
	NetMinorUnits    int64  `json:"net_minor_units"`
}
type Repository interface {
	CreateAccount(context.Context, string, string, Account) error
	OpenPeriod(context.Context, string, string, Period) error
	CreateJournal(context.Context, string, string, Journal) error
	PostJournal(context.Context, string, string, string, string, int64, string, string) (Journal, error)
	ReverseJournal(context.Context, string, string, string, string, string, string, int64, string, string) (Journal, error)
	ClosePeriod(context.Context, string, string, int64, string) (Period, error)
	TrialBalance(context.Context, string, string, string) ([]Balance, error)
}
type IDGenerator interface{ New() string }
type Service struct {
	repository Repository
	ids        IDGenerator
}

func NewService(repository Repository, ids IDGenerator) *Service {
	return &Service{repository: repository, ids: ids}
}

func (s *Service) CreateAccount(ctx context.Context, tenant string, value Account) (Account, error) {
	validType := value.Type == "asset" || value.Type == "liability" || value.Type == "equity" || value.Type == "revenue" || value.Type == "expense"
	if tenant == "" || !codePattern.MatchString(value.Code) || len(value.Name) < 2 || len(value.Name) > 120 || !validType {
		return value, ErrInvalid
	}
	if err := s.repository.CreateAccount(ctx, tenant, s.ids.New(), value); err != nil {
		return value, err
	}
	return value, nil
}
func (s *Service) OpenPeriod(ctx context.Context, tenant string, value Period) (Period, error) {
	if tenant == "" || value.StartsOn.IsZero() || !value.EndsOn.After(value.StartsOn) {
		return value, ErrInvalid
	}
	value.ID = s.ids.New()
	value.Status = "open"
	value.Version = 1
	if err := s.repository.OpenPeriod(ctx, tenant, s.ids.New(), value); err != nil {
		return value, err
	}
	return value, nil
}
func (s *Service) CreateJournal(ctx context.Context, tenant string, value Journal) (Journal, error) {
	// Reserved for the typed conversion binding; generic input cannot forge it.
	if value.SourceType == "FX_CONVERSION" {
		return value, ErrInvalid
	}
	if tenant == "" || value.OrganizationID == "" || value.PeriodID == "" || !codePattern.MatchString(value.SourceType) || value.SourceID == "" || !currencyPattern.MatchString(value.Currency) || value.PostingDate.IsZero() || len(value.Lines) < 2 {
		return value, ErrInvalid
	}
	seen := map[int]struct{}{}
	entries := make([]bcamounts.Entry, 0, len(value.Lines))
	for _, line := range value.Lines {
		if line.LineNo < 1 || !codePattern.MatchString(line.AccountCode) || len(line.Description) > 250 || ((line.DebitMinorUnits > 0) == (line.CreditMinorUnits > 0)) {
			return value, ErrInvalid
		}
		if _, ok := seen[line.LineNo]; ok {
			return value, ErrInvalid
		}
		seen[line.LineNo] = struct{}{}
		entries = append(entries, bcamounts.Entry{Debit: line.DebitMinorUnits, Credit: line.CreditMinorUnits})
	}
	debit, credit, amountErr := bcamounts.JournalTotals(entries)
	if amountErr != nil {
		return value, ErrInvalid
	}
	value.ID = s.ids.New()
	value.Status = "draft"
	value.Version = 1
	value.TotalDebitMinorUnits = debit
	value.TotalCreditMinorUnits = credit
	if err := s.repository.CreateJournal(ctx, tenant, s.ids.New(), value); err != nil {
		return value, err
	}
	return value, nil
}
func (s *Service) PostJournal(ctx context.Context, tenant, organization, id, actor string, version int64) (Journal, error) {
	if tenant == "" || organization == "" || id == "" || actor == "" || version < 1 {
		return Journal{}, ErrInvalid
	}
	return s.repository.PostJournal(ctx, tenant, organization, id, actor, version, s.ids.New(), s.ids.New())
}
func (s *Service) ReverseJournal(ctx context.Context, tenant, organization, id, actor string, version int64, reason string) (Journal, error) {
	if tenant == "" || organization == "" || id == "" || actor == "" || version < 1 || len(reason) < 3 || len(reason) > 500 {
		return Journal{}, ErrInvalid
	}
	return s.repository.ReverseJournal(ctx, tenant, organization, id, s.ids.New(), s.ids.New(), s.ids.New(), version, actor, reason)
}
func (s *Service) ClosePeriod(ctx context.Context, tenant, id string, version int64) (Period, error) {
	if tenant == "" || id == "" || version < 1 {
		return Period{}, ErrInvalid
	}
	return s.repository.ClosePeriod(ctx, tenant, id, version, s.ids.New())
}
func (s *Service) TrialBalance(ctx context.Context, tenant, organization, period string) ([]Balance, error) {
	if tenant == "" || organization == "" || period == "" {
		return nil, ErrInvalid
	}
	return s.repository.TrialBalance(ctx, tenant, organization, period)
}
````

### FILE: `internal/accounting/service_test.go`
```yaml
block_id: "GO-ACCOUNTING:service-test:v1"
operation: CREATE
provenance: AUTHORED
source: "local balance regression"
license: "LicenseRef-Workspace-Owner"
sha256: "52ee0db68deda10768e206482c46e6e93d98a6ced5e3d4ec66bddaba37aae6f2"
variables: []
secrets_allowed: false
```
````go
package accounting

import (
	"context"
	"errors"
	"testing"
	"time"
)

type fakeRepo struct{ journal Journal }

func (*fakeRepo) CreateAccount(context.Context, string, string, Account) error { return nil }
func (*fakeRepo) OpenPeriod(context.Context, string, string, Period) error     { return nil }
func (f *fakeRepo) CreateJournal(_ context.Context, _, _ string, v Journal) error {
	f.journal = v
	return nil
}
func (*fakeRepo) PostJournal(context.Context, string, string, string, string, int64, string, string) (Journal, error) {
	return Journal{Status: "posted"}, nil
}
func (*fakeRepo) ReverseJournal(context.Context, string, string, string, string, string, string, int64, string, string) (Journal, error) {
	return Journal{Status: "posted", ReversalOf: "j"}, nil
}
func (*fakeRepo) ClosePeriod(context.Context, string, string, int64, string) (Period, error) {
	return Period{Status: "closed"}, nil
}
func (*fakeRepo) TrialBalance(context.Context, string, string, string) ([]Balance, error) {
	return []Balance{{AccountCode: "CASH"}}, nil
}

type ids struct{ n int }

func (i *ids) New() string { i.n++; return "id" + string(rune('0'+i.n)) }
func TestBalancedJournal(t *testing.T) {
	r := &fakeRepo{}
	s := NewService(r, &ids{})
	j, err := s.CreateJournal(context.Background(), "tenant", Journal{OrganizationID: "org", PeriodID: "p", SourceType: "SALE", SourceID: "order", Currency: "ARS", PostingDate: time.Now(), Lines: []Line{{LineNo: 1, AccountCode: "CASH", DebitMinorUnits: 100}, {LineNo: 2, AccountCode: "REVENUE", CreditMinorUnits: 100}}})
	if err != nil || j.TotalDebitMinorUnits != 100 || r.journal.ID == "" {
		t.Fatalf("journal=%+v err=%v", j, err)
	}
}
func TestRejectsUnbalancedOrDualSided(t *testing.T) {
	s := NewService(&fakeRepo{}, &ids{})
	for _, lines := range [][]Line{{{LineNo: 1, AccountCode: "CASH", DebitMinorUnits: 100}, {LineNo: 2, AccountCode: "REVENUE", CreditMinorUnits: 99}}, {{LineNo: 1, AccountCode: "CASH", DebitMinorUnits: 100, CreditMinorUnits: 1}, {LineNo: 2, AccountCode: "REVENUE", CreditMinorUnits: 100}}} {
		_, err := s.CreateJournal(context.Background(), "tenant", Journal{OrganizationID: "org", PeriodID: "p", SourceType: "SALE", SourceID: "order", Currency: "ARS", PostingDate: time.Now(), Lines: lines})
		if !errors.Is(err, ErrInvalid) {
			t.Fatalf("accepted lines=%+v err=%v", lines, err)
		}
	}
}
````

### FILE: `internal/platform/postgres/accounting.go`
```yaml
block_id: "GO-ACCOUNTING:postgres:v1"
operation: CREATE
provenance: AUTHORED
source: "local PostgreSQL implementation governed by pinned Microsoft BCApps evidence"
license: "LicenseRef-Workspace-Owner"
sha256: "d31c8cdeffd16e7a308027b19d50470b6ef64ee2cdafba4b8b0524ce690d204d"
variables: []
secrets_allowed: false
```
````go
package postgres

import (
	"context"
	"elite.local/enterprise/internal/accounting"
	"elite.local/enterprise/internal/bcamounts"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Accounting struct{ pool *pgxpool.Pool }

func NewAccounting(pool *pgxpool.Pool) *Accounting { return &Accounting{pool: pool} }
func accountingConflict(err error) error {
	if err == nil {
		return nil
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && (pgErr.Code == "23503" || pgErr.Code == "23505" || pgErr.Code == "23514" || pgErr.Code == "40001") {
		return accounting.ErrConflict
	}
	return err
}

func (r *Accounting) CreateAccount(ctx context.Context, tenant, eventID string, value accounting.Account) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, `insert into accounting.account(tenant_id,account_code,display_name,account_type)values($1,$2,$3,$4)`, tenant, value.Code, value.Name, value.Type)
	if err != nil {
		return accountingConflict(err)
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)values($1,$2,'account',$3,1,'account.created',1,clock_timestamp(),jsonb_build_object('account_type',$4::text))`, tenant, eventID, value.Code, value.Type)
	if err != nil {
		return accountingConflict(err)
	}
	return accountingConflict(tx.Commit(ctx))
}

func (r *Accounting) OpenPeriod(ctx context.Context, tenant, eventID string, value accounting.Period) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if _, err = tx.Exec(ctx, `select pg_advisory_xact_lock(hashtextextended($1,0))`, tenant+":accounting-period"); err != nil {
		return err
	}
	var overlap bool
	err = tx.QueryRow(ctx, `select exists(select 1 from accounting.period where tenant_id=$1 and starts_on < ($3::timestamptz at time zone 'UTC')::date and ends_on > ($2::timestamptz at time zone 'UTC')::date)`, tenant, value.StartsOn, value.EndsOn).Scan(&overlap)
	if err != nil {
		return err
	}
	if overlap {
		return accounting.ErrConflict
	}
	_, err = tx.Exec(ctx, `insert into accounting.period(tenant_id,period_id,starts_on,ends_on,status,version)values($1,$2,($3::timestamptz at time zone 'UTC')::date,($4::timestamptz at time zone 'UTC')::date,'open',1)`, tenant, value.ID, value.StartsOn, value.EndsOn)
	if err != nil {
		return accountingConflict(err)
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)values($1,$2,'accounting-period',$3,1,'accounting-period.opened',1,clock_timestamp(),'{}')`, tenant, eventID, value.ID)
	if err != nil {
		return accountingConflict(err)
	}
	return accountingConflict(tx.Commit(ctx))
}

func (r *Accounting) CreateJournal(ctx context.Context, tenant, eventID string, value accounting.Journal) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err = createJournalInTx(ctx, tx, tenant, eventID, value); err != nil {
		return err
	}
	return accountingConflict(tx.Commit(ctx))
}

// AUTHORED transaction composition seam. This is the same existing draft writer.
func createJournalInTx(ctx context.Context, tx pgx.Tx, tenant, eventID string, value accounting.Journal) error {
	var period string
	err := tx.QueryRow(ctx, `select period_id from accounting.period where tenant_id=$1 and period_id=$2 and status='open' and starts_on <= ($3::timestamptz at time zone 'UTC')::date and ends_on > ($3::timestamptz at time zone 'UTC')::date for share`, tenant, value.PeriodID, value.PostingDate).Scan(&period)
	if errors.Is(err, pgx.ErrNoRows) {
		return accounting.ErrConflict
	}
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `insert into accounting.journal(tenant_id,journal_id,organization_id,period_id,source_type,source_id,currency,posting_date,status,total_debit_minor_units,total_credit_minor_units,version)values($1,$2,$3,$4,$5,$6,$7,($8::timestamptz at time zone 'UTC')::date,'draft',$9,$10,1)`, tenant, value.ID, value.OrganizationID, value.PeriodID, value.SourceType, value.SourceID, value.Currency, value.PostingDate, value.TotalDebitMinorUnits, value.TotalCreditMinorUnits)
	if err != nil {
		return accountingConflict(err)
	}
	for _, line := range value.Lines {
		result, lineErr := tx.Exec(ctx, `insert into accounting.journal_line(tenant_id,journal_id,line_no,account_code,description,debit_minor_units,credit_minor_units)select $1,$2,$3,$4,$5,$6,$7 where exists(select 1 from accounting.account where tenant_id=$1 and account_code=$4 and active)`, tenant, value.ID, line.LineNo, line.AccountCode, line.Description, line.DebitMinorUnits, line.CreditMinorUnits)
		if lineErr != nil {
			return accountingConflict(lineErr)
		}
		if result.RowsAffected() != 1 {
			return accounting.ErrConflict
		}
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)values($1,$2,'journal',$3,1,'journal.created',1,clock_timestamp(),jsonb_build_object('organization_id',$4::text,'source_type',$5::text,'source_id',$6::text,'currency',$7::text))`, tenant, eventID, value.ID, value.OrganizationID, value.SourceType, value.SourceID, value.Currency)
	if err != nil {
		return accountingConflict(err)
	}
	return nil
}

func (r *Accounting) PostJournal(ctx context.Context, tenant, organization, id, actor string, expectedVersion int64, registerID, eventID string) (accounting.Journal, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return accounting.Journal{}, err
	}
	defer tx.Rollback(ctx)
	var value accounting.Journal
	err = tx.QueryRow(ctx, `select j.journal_id,j.organization_id,j.period_id,j.source_type,j.source_id,j.currency,j.posting_date,j.status,j.total_debit_minor_units,j.total_credit_minor_units,j.version from accounting.journal j join accounting.period p on p.tenant_id=j.tenant_id and p.period_id=j.period_id where j.tenant_id=$1 and j.journal_id=$2 and j.organization_id=$3 and j.status='draft' and j.version=$4 and p.status='open' for update of j,p`, tenant, id, organization, expectedVersion).Scan(&value.ID, &value.OrganizationID, &value.PeriodID, &value.SourceType, &value.SourceID, &value.Currency, &value.PostingDate, &value.Status, &value.TotalDebitMinorUnits, &value.TotalCreditMinorUnits, &value.Version)
	if errors.Is(err, pgx.ErrNoRows) {
		return value, accounting.ErrConflict
	}
	if err != nil {
		return value, err
	}
	var debit, credit int64
	err = tx.QueryRow(ctx, `select coalesce(sum(debit_minor_units),0),coalesce(sum(credit_minor_units),0) from accounting.journal_line where tenant_id=$1 and journal_id=$2`, tenant, id).Scan(&debit, &credit)
	if err != nil {
		return value, err
	}
	if bcamounts.CheckBalance(debit, credit) != nil || debit != value.TotalDebitMinorUnits {
		return value, accounting.ErrConflict
	}
	_, err = tx.Exec(ctx, `update accounting.journal set status='posted',version=version+1,posted_by=$4,posted_at=clock_timestamp() where tenant_id=$1 and journal_id=$2 and version=$3`, tenant, id, expectedVersion, actor)
	if err != nil {
		return value, err
	}
	_, err = tx.Exec(ctx, `insert into accounting.register(tenant_id,register_id,journal_id,organization_id,period_id,posted_by,total_debit_minor_units,total_credit_minor_units)values($1,$2,$3,$4,$5,$6,$7,$7)`, tenant, registerID, id, organization, value.PeriodID, actor, debit)
	if err != nil {
		return value, accountingConflict(err)
	}
	result, err := tx.Exec(ctx, `insert into accounting.entry(tenant_id,entry_id,register_id,journal_id,line_no,organization_id,period_id,account_code,currency,posting_date,description,debit_minor_units,credit_minor_units)select l.tenant_id,$3||':'||l.line_no::text,$3,l.journal_id,l.line_no,$4,$5,l.account_code,$6,$7,l.description,l.debit_minor_units,l.credit_minor_units from accounting.journal_line l where l.tenant_id=$1 and l.journal_id=$2 order by l.line_no`, tenant, id, registerID, organization, value.PeriodID, value.Currency, value.PostingDate)
	if err != nil {
		return value, accountingConflict(err)
	}
	if result.RowsAffected() < 2 {
		return value, accounting.ErrConflict
	}
	value.Status = "posted"
	value.Version = expectedVersion + 1
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)values($1,$2,'journal',$3,$4,'journal.posted',1,clock_timestamp(),jsonb_build_object('organization_id',$5::text,'register_id',$6::text,'currency',$7::text,'total_minor_units',$8::bigint))`, tenant, eventID, id, value.Version, organization, registerID, value.Currency, debit)
	if err != nil {
		return value, accountingConflict(err)
	}
	if err = tx.Commit(ctx); err != nil {
		return value, accountingConflict(err)
	}
	return value, nil
}

func (r *Accounting) ReverseJournal(ctx context.Context, tenant, organization, id, reversalID, registerID, eventID string, expectedVersion int64, actor, reason string) (accounting.Journal, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return accounting.Journal{}, err
	}
	defer tx.Rollback(ctx)
	var original accounting.Journal
	err = tx.QueryRow(ctx, `select j.period_id,j.currency,j.posting_date,j.total_debit_minor_units,j.total_credit_minor_units from accounting.journal j join accounting.period p on p.tenant_id=j.tenant_id and p.period_id=j.period_id where j.tenant_id=$1 and j.journal_id=$2 and j.organization_id=$3 and j.status='posted' and j.version=$4 and j.reversal_of is null and p.status='open' for update of j,p`, tenant, id, organization, expectedVersion).Scan(&original.PeriodID, &original.Currency, &original.PostingDate, &original.TotalDebitMinorUnits, &original.TotalCreditMinorUnits)
	if errors.Is(err, pgx.ErrNoRows) {
		return original, accounting.ErrConflict
	}
	if err != nil {
		return original, err
	}
	_, err = tx.Exec(ctx, `insert into accounting.journal(tenant_id,journal_id,organization_id,period_id,source_type,source_id,currency,posting_date,status,total_debit_minor_units,total_credit_minor_units,version,reversal_of,reversal_reason,posted_by,posted_at)values($1,$2,$3,$4,'REVERSAL',$5,$6,$7,'posted',$8,$8,1,$5,$9,$10,clock_timestamp())`, tenant, reversalID, organization, original.PeriodID, id, original.Currency, original.PostingDate, original.TotalDebitMinorUnits, reason, actor)
	if err != nil {
		return original, accountingConflict(err)
	}
	_, err = tx.Exec(ctx, reverseJournalLinesBCSQL, tenant, id, reversalID)
	if err != nil {
		return original, accountingConflict(err)
	}
	_, err = tx.Exec(ctx, `insert into accounting.register(tenant_id,register_id,journal_id,organization_id,period_id,posted_by,total_debit_minor_units,total_credit_minor_units)values($1,$2,$3,$4,$5,$6,$7,$7)`, tenant, registerID, reversalID, organization, original.PeriodID, actor, original.TotalDebitMinorUnits)
	if err != nil {
		return original, accountingConflict(err)
	}
	result, err := tx.Exec(ctx, `insert into accounting.entry(tenant_id,entry_id,register_id,journal_id,line_no,organization_id,period_id,account_code,currency,posting_date,description,debit_minor_units,credit_minor_units,original_entry_id)select l.tenant_id,$4||':'||l.line_no::text,$4,l.journal_id,l.line_no,$5,$6,l.account_code,$7,$8,l.description,l.debit_minor_units,l.credit_minor_units,e.entry_id from accounting.journal_line l join accounting.entry e on e.tenant_id=l.tenant_id and e.journal_id=$2 and e.line_no=l.line_no where l.tenant_id=$1 and l.journal_id=$3`, tenant, id, reversalID, registerID, organization, original.PeriodID, original.Currency, original.PostingDate)
	if err != nil {
		return original, accountingConflict(err)
	}
	if result.RowsAffected() < 2 {
		return original, accounting.ErrConflict
	}
	result, err = tx.Exec(ctx, `update accounting.journal set status='reversed',version=version+1 where tenant_id=$1 and journal_id=$2 and status='posted' and version=$3`, tenant, id, expectedVersion)
	if err != nil {
		return original, err
	}
	if result.RowsAffected() != 1 {
		return original, accounting.ErrConflict
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)values($1,$2,'journal',$3,1,'journal.reversal-posted',1,clock_timestamp(),jsonb_build_object('reversal_of',$4::text,'organization_id',$5::text,'register_id',$6::text,'reason',$7::text))`, tenant, eventID, reversalID, id, organization, registerID, reason)
	if err != nil {
		return original, accountingConflict(err)
	}
	if err = tx.Commit(ctx); err != nil {
		return original, accountingConflict(err)
	}
	return accounting.Journal{ID: reversalID, OrganizationID: organization, PeriodID: original.PeriodID, SourceType: "REVERSAL", SourceID: id, Currency: original.Currency, PostingDate: original.PostingDate, Status: "posted", TotalDebitMinorUnits: original.TotalDebitMinorUnits, TotalCreditMinorUnits: original.TotalCreditMinorUnits, Version: 1, ReversalOf: id}, nil
}

func (r *Accounting) ClosePeriod(ctx context.Context, tenant, id string, expectedVersion int64, eventID string) (accounting.Period, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return accounting.Period{}, err
	}
	defer tx.Rollback(ctx)
	if _, err = tx.Exec(ctx, `select pg_advisory_xact_lock(hashtextextended($1,0))`, tenant+":accounting-period"); err != nil {
		return accounting.Period{}, err
	}
	var value accounting.Period
	err = tx.QueryRow(ctx, `select period_id,starts_on,ends_on,status,version from accounting.period where tenant_id=$1 and period_id=$2 and status='open' and version=$3 and not exists(select 1 from accounting.journal where tenant_id=$1 and period_id=$2 and status='draft') for update`, tenant, id, expectedVersion).Scan(&value.ID, &value.StartsOn, &value.EndsOn, &value.Status, &value.Version)
	if errors.Is(err, pgx.ErrNoRows) {
		return value, accounting.ErrConflict
	}
	if err != nil {
		return value, err
	}
	_, err = tx.Exec(ctx, `update accounting.period set status='closed',version=version+1,closed_at=clock_timestamp() where tenant_id=$1 and period_id=$2 and version=$3`, tenant, id, expectedVersion)
	if err != nil {
		return value, err
	}
	value.Status = "closed"
	value.Version++
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)values($1,$2,'accounting-period',$3,$4,'accounting-period.closed',1,clock_timestamp(),'{}')`, tenant, eventID, id, value.Version)
	if err != nil {
		return value, accountingConflict(err)
	}
	if err = tx.Commit(ctx); err != nil {
		return value, accountingConflict(err)
	}
	return value, nil
}

func (r *Accounting) TrialBalance(ctx context.Context, tenant, organization, period string) ([]accounting.Balance, error) {
	rows, err := r.pool.Query(ctx, `select account_code,sum(debit_minor_units)::bigint,sum(credit_minor_units)::bigint,(sum(debit_minor_units)-sum(credit_minor_units))::bigint from accounting.entry where tenant_id=$1 and organization_id=$2 and period_id=$3 group by account_code order by account_code`, tenant, organization, period)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	values := []accounting.Balance{}
	for rows.Next() {
		var v accounting.Balance
		if err = rows.Scan(&v.AccountCode, &v.DebitMinorUnits, &v.CreditMinorUnits, &v.NetMinorUnits); err != nil {
			return nil, err
		}
		values = append(values, v)
	}
	return values, rows.Err()
}
````

### FILE: `internal/platform/postgres/accounting_integration_test.go`
```yaml
block_id: "GO-ACCOUNTING:postgres-test:v1"
operation: CREATE
provenance: AUTHORED
source: "local posting concurrency and reversal regression"
license: "LicenseRef-Workspace-Owner"
sha256: "f2a45ec05adcfbeb78d087f292441c80ad5aec54f365df2a193e6bbc31e9b682"
variables: []
secrets_allowed: false
```
````go
package postgres

import (
	"context"
	"crypto/rand"
	"elite.local/enterprise/internal/accounting"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"sync"
	"testing"
	"time"
)

func TestAccountingPostingConcurrencyReversalAndClose(t *testing.T) {
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
	var raw [16]byte
	if _, err = rand.Read(raw[:]); err != nil {
		t.Fatal(err)
	}
	raw[6] = (raw[6] & 0x0f) | 0x40
	raw[8] = (raw[8] & 0x3f) | 0x80
	tenant := fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", raw[0:4], raw[4:6], raw[6:8], raw[8:10], raw[10:16])
	defer func() {
		tx, e := pool.Begin(ctx)
		if e != nil {
			t.Errorf("cleanup begin: %v", e)
			return
		}
		defer tx.Rollback(ctx)
		if _, e = tx.Exec(ctx, `set local session_replication_role=replica`); e != nil {
			t.Errorf("cleanup role: %v", e)
			return
		}
		for _, q := range []string{`delete from platform.outbox_event where tenant_id=$1`, `delete from accounting.entry where tenant_id=$1`, `delete from accounting.register where tenant_id=$1`, `delete from accounting.journal_line where tenant_id=$1`, `delete from accounting.journal where tenant_id=$1`, `delete from accounting.period where tenant_id=$1`, `delete from accounting.account where tenant_id=$1`, `delete from org.organization where tenant_id=$1`, `delete from platform.tenant where tenant_id=$1`} {
			if _, e = tx.Exec(ctx, q, tenant); e != nil {
				t.Errorf("cleanup: %v", e)
				return
			}
		}
		if e = tx.Commit(ctx); e != nil {
			t.Errorf("cleanup commit: %v", e)
		}
	}()
	for _, q := range []string{`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1::uuid,'acct-'||substring($1::text,1,8),'Accounting','Accounting')`, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'franchise','franchise','Franchise','franchisee'),($1,'other','other','Other','franchisee')`} {
		if _, err = pool.Exec(ctx, q, tenant); err != nil {
			t.Fatal(err)
		}
	}
	repo := NewAccounting(pool)
	if err = repo.CreateAccount(ctx, tenant, "018f4d4a-7b36-7a21-8d10-2f4c54c29e01", accounting.Account{Code: "CASH", Name: "Cash", Type: "asset"}); err != nil {
		t.Fatal(err)
	}
	if err = repo.CreateAccount(ctx, tenant, "018f4d4a-7b36-7a21-8d10-2f4c54c29e02", accounting.Account{Code: "REVENUE", Name: "Revenue", Type: "revenue"}); err != nil {
		t.Fatal(err)
	}
	from := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	until := from.AddDate(0, 1, 0)
	if err = repo.OpenPeriod(ctx, tenant, "018f4d4a-7b36-7a21-8d10-2f4c54c29e03", accounting.Period{ID: "2026-08", StartsOn: from, EndsOn: until, Status: "open", Version: 1}); err != nil {
		t.Fatal(err)
	}
	if err = repo.OpenPeriod(ctx, tenant, "unused", accounting.Period{ID: "overlap", StartsOn: from.AddDate(0, 0, 15), EndsOn: until.AddDate(0, 0, 15), Status: "open", Version: 1}); !errors.Is(err, accounting.ErrConflict) {
		t.Fatalf("overlap accepted: %v", err)
	}
	journal := accounting.Journal{ID: "journal", OrganizationID: "franchise", PeriodID: "2026-08", SourceType: "SALE", SourceID: "order-1", Currency: "ARS", PostingDate: from.Add(12 * time.Hour), Status: "draft", TotalDebitMinorUnits: 10001, TotalCreditMinorUnits: 10001, Version: 1, Lines: []accounting.Line{{LineNo: 1, AccountCode: "CASH", Description: "cash", DebitMinorUnits: 10001}, {LineNo: 2, AccountCode: "REVENUE", Description: "sale", CreditMinorUnits: 10001}}}
	if err = repo.CreateJournal(ctx, tenant, "018f4d4a-7b36-7a21-8d10-2f4c54c29e04", journal); err != nil {
		t.Fatal(err)
	}
	type result struct {
		journal accounting.Journal
		err     error
	}
	results := make(chan result, 2)
	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-start
			j, e := repo.PostJournal(ctx, tenant, "franchise", "journal", "controller", 1, []string{"register-a", "register-b"}[i], []string{"018f4d4a-7b36-7a21-8d10-2f4c54c29e05", "018f4d4a-7b36-7a21-8d10-2f4c54c29e06"}[i])
			results <- result{j, e}
		}(i)
	}
	close(start)
	wg.Wait()
	close(results)
	success, conflict := 0, 0
	for got := range results {
		if got.err == nil {
			success++
		} else if errors.Is(got.err, accounting.ErrConflict) {
			conflict++
		} else {
			t.Fatal(got.err)
		}
	}
	if success != 1 || conflict != 1 {
		t.Fatalf("posting success=%d conflict=%d", success, conflict)
	}
	balances, err := repo.TrialBalance(ctx, tenant, "franchise", "2026-08")
	if err != nil || len(balances) != 2 {
		t.Fatalf("balances=%+v err=%v", balances, err)
	}
	reversal, err := repo.ReverseJournal(ctx, tenant, "franchise", "journal", "journal-reversal", "register-reversal", "018f4d4a-7b36-7a21-8d10-2f4c54c29e07", 2, "controller", "duplicate sale correction")
	if err != nil || reversal.ReversalOf != "journal" {
		t.Fatalf("reversal=%+v err=%v", reversal, err)
	}
	if _, err = repo.ReverseJournal(ctx, tenant, "franchise", "journal", "second", "register-second", "unused", 3, "controller", "again"); !errors.Is(err, accounting.ErrConflict) {
		t.Fatalf("second reversal accepted: %v", err)
	}
	balances, err = repo.TrialBalance(ctx, tenant, "franchise", "2026-08")
	if err != nil {
		t.Fatal(err)
	}
	for _, balance := range balances {
		if balance.NetMinorUnits != 0 {
			t.Fatalf("non-zero after reversal: %+v", balances)
		}
	}
	period, err := repo.ClosePeriod(ctx, tenant, "2026-08", 1, "018f4d4a-7b36-7a21-8d10-2f4c54c29e08")
	if err != nil || period.Status != "closed" || period.Version != 2 {
		t.Fatalf("period=%+v err=%v", period, err)
	}
	journal.ID = "after-close"
	journal.SourceID = "order-2"
	if err = repo.CreateJournal(ctx, tenant, "unused", journal); !errors.Is(err, accounting.ErrConflict) {
		t.Fatalf("journal in closed period accepted: %v", err)
	}
}
````

### FILE: `internal/platform/httpapi/accounting.go`
```yaml
block_id: "GO-ACCOUNTING:http:v1"
operation: CREATE
provenance: AUTHORED
source: "local protected accounting boundary"
license: "LicenseRef-Workspace-Owner"
sha256: "6b11017d819271d3a4c7f604caa1c5a7d05f71ba9cdebaf93110b57e8036556e"
variables: []
secrets_allowed: false
```
````go
package httpapi

import (
	"elite.local/enterprise/internal/accounting"
	"elite.local/enterprise/internal/platform/identity"
	"errors"
	"net/http"
)

type AccountingModule struct{ Service *accounting.Service }

func (m AccountingModule) Register(mux *http.ServeMux, verifier identity.Verifier) {
	api := accountingAPI{service: m.Service, verifier: verifier}
	mux.HandleFunc("POST /v1/accounting/accounts", api.createAccount)
	mux.HandleFunc("POST /v1/accounting/periods", api.openPeriod)
	mux.HandleFunc("POST /v1/accounting/journals", api.createJournal)
	mux.HandleFunc("POST /v1/accounting/journals/{id}/post", api.postJournal)
	mux.HandleFunc("POST /v1/accounting/journals/{id}/reverse", api.reverseJournal)
	mux.HandleFunc("POST /v1/accounting/periods/{id}/close", api.closePeriod)
	mux.HandleFunc("GET /v1/accounting/trial-balance", api.trialBalance)
}

type accountingAPI struct {
	service  *accounting.Service
	verifier identity.Verifier
}

func (a accountingAPI) auth(w http.ResponseWriter, r *http.Request, permission string, jsonBody bool) (identity.Principal, bool) {
	p, err := authenticate(r.Context(), r.Header.Get("Authorization"), a.verifier)
	if err != nil {
		writeProblem(w, 401, "UNAUTHENTICATED", "a valid bearer token is required")
		return p, false
	}
	if !p.Allowed(permission) {
		writeProblem(w, 403, "FORBIDDEN", permission+" permission is required")
		return p, false
	}
	if jsonBody && r.Header.Get("Content-Type") != "application/json" {
		writeProblem(w, 415, "UNSUPPORTED_MEDIA_TYPE", "Content-Type must be application/json")
		return p, false
	}
	return p, true
}
func accountingOrg(w http.ResponseWriter, p identity.Principal, id string) bool {
	if !p.AllowedOrganization(id) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return false
	}
	return true
}
func (a accountingAPI) createAccount(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "accounting:manage", true)
	if !ok {
		return
	}
	var input accounting.Account
	if !decodeStrict(w, r, &input) {
		return
	}
	value, err := a.service.CreateAccount(r.Context(), p.TenantID, input)
	if err != nil {
		writeAccountingResult(w, err)
		return
	}
	writeJSON(w, 201, value)
}
func (a accountingAPI) openPeriod(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "accounting:manage", true)
	if !ok {
		return
	}
	var input accounting.Period
	if !decodeStrict(w, r, &input) {
		return
	}
	value, err := a.service.OpenPeriod(r.Context(), p.TenantID, input)
	if err != nil {
		writeAccountingResult(w, err)
		return
	}
	writeJSON(w, 201, value)
}
func (a accountingAPI) createJournal(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "accounting:write", true)
	if !ok {
		return
	}
	var input accounting.Journal
	if !decodeStrict(w, r, &input) {
		return
	}
	if !accountingOrg(w, p, input.OrganizationID) {
		return
	}
	value, err := a.service.CreateJournal(r.Context(), p.TenantID, input)
	if err != nil {
		writeAccountingResult(w, err)
		return
	}
	writeJSON(w, 201, value)
}
func (a accountingAPI) postJournal(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "accounting:post", true)
	if !ok {
		return
	}
	var input struct {
		OrganizationID  string `json:"organization_id"`
		ExpectedVersion int64  `json:"expected_version"`
	}
	if !decodeStrict(w, r, &input) {
		return
	}
	if !accountingOrg(w, p, input.OrganizationID) {
		return
	}
	value, err := a.service.PostJournal(r.Context(), p.TenantID, input.OrganizationID, r.PathValue("id"), p.Subject, input.ExpectedVersion)
	if err != nil {
		writeAccountingResult(w, err)
		return
	}
	writeJSON(w, 200, value)
}
func (a accountingAPI) reverseJournal(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "accounting:reverse", true)
	if !ok {
		return
	}
	var input struct {
		OrganizationID  string `json:"organization_id"`
		ExpectedVersion int64  `json:"expected_version"`
		Reason          string `json:"reason"`
	}
	if !decodeStrict(w, r, &input) {
		return
	}
	if !accountingOrg(w, p, input.OrganizationID) {
		return
	}
	value, err := a.service.ReverseJournal(r.Context(), p.TenantID, input.OrganizationID, r.PathValue("id"), p.Subject, input.ExpectedVersion, input.Reason)
	if err != nil {
		writeAccountingResult(w, err)
		return
	}
	writeJSON(w, 201, value)
}
func (a accountingAPI) closePeriod(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "accounting:close", true)
	if !ok {
		return
	}
	var input struct {
		ExpectedVersion int64 `json:"expected_version"`
	}
	if !decodeStrict(w, r, &input) {
		return
	}
	value, err := a.service.ClosePeriod(r.Context(), p.TenantID, r.PathValue("id"), input.ExpectedVersion)
	if err != nil {
		writeAccountingResult(w, err)
		return
	}
	writeJSON(w, 200, value)
}
func (a accountingAPI) trialBalance(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "accounting:read", false)
	if !ok {
		return
	}
	organization := r.URL.Query().Get("organization_id")
	if !accountingOrg(w, p, organization) {
		return
	}
	value, err := a.service.TrialBalance(r.Context(), p.TenantID, organization, r.URL.Query().Get("period_id"))
	if err != nil {
		writeAccountingResult(w, err)
		return
	}
	writeJSON(w, 200, value)
}
func writeAccountingResult(w http.ResponseWriter, err error) {
	if errors.Is(err, accounting.ErrInvalid) {
		writeProblem(w, 400, "INVALID_ACCOUNTING_COMMAND", "accounting command does not match contract")
		return
	}
	if errors.Is(err, accounting.ErrConflict) {
		writeProblem(w, 409, "ACCOUNTING_CONFLICT", "accounting state, scope, balance or version conflict")
		return
	}
	writeProblem(w, 500, "INTERNAL", "request failed")
}
````

### FILE: `internal/platform/httpapi/accounting_test.go`
```yaml
block_id: "GO-ACCOUNTING:http-test:v1"
operation: CREATE
provenance: AUTHORED
source: "local scope and schema regression"
license: "LicenseRef-Workspace-Owner"
sha256: "d655b5b405d1eef4c4c2808d8dac06bf5a1ca0c343c9052bae6e41f9445f95d0"
variables: []
secrets_allowed: false
```
````go
package httpapi

import (
	"context"
	"elite.local/enterprise/internal/accounting"
	"elite.local/enterprise/internal/platform/identity"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type accountingRepo struct{ journals int }

func (*accountingRepo) CreateAccount(context.Context, string, string, accounting.Account) error {
	return nil
}
func (*accountingRepo) OpenPeriod(context.Context, string, string, accounting.Period) error {
	return nil
}
func (r *accountingRepo) CreateJournal(_ context.Context, _, _ string, _ accounting.Journal) error {
	r.journals++
	return nil
}
func (*accountingRepo) PostJournal(context.Context, string, string, string, string, int64, string, string) (accounting.Journal, error) {
	return accounting.Journal{Status: "posted"}, nil
}
func (*accountingRepo) ReverseJournal(context.Context, string, string, string, string, string, string, int64, string, string) (accounting.Journal, error) {
	return accounting.Journal{Status: "posted"}, nil
}
func (*accountingRepo) ClosePeriod(context.Context, string, string, int64, string) (accounting.Period, error) {
	return accounting.Period{Status: "closed"}, nil
}
func (*accountingRepo) TrialBalance(context.Context, string, string, string) ([]accounting.Balance, error) {
	return []accounting.Balance{}, nil
}

type accountingIDs struct{ n int }

func (i *accountingIDs) New() string { i.n++; return "id" + string(rune('0'+i.n)) }

type accountingVerifier struct{ p identity.Principal }

func (v accountingVerifier) Verify(context.Context, string) (identity.Principal, error) {
	return v.p, nil
}
func TestAccountingHTTPRejectsCrossScopeAndUnbalanced(t *testing.T) {
	repo := &accountingRepo{}
	service := accounting.NewService(repo, &accountingIDs{})
	mux := http.NewServeMux()
	p := identity.Principal{Subject: "controller", TenantID: "018f4d4a-7b36-7a21-8d10-2f4c54c29f00", Permissions: map[string]struct{}{"accounting:write": {}}, Organizations: map[string]struct{}{"franchise": {}}}
	AccountingModule{Service: service}.Register(mux, accountingVerifier{p: p})
	now := time.Now().UTC()
	body := `{"organization_id":"other","period_id":"p","source_type":"SALE","source_id":"o","currency":"ARS","posting_date":"` + now.Format(time.RFC3339) + `","lines":[{"line_no":1,"account_code":"CASH","debit_minor_units":100},{"line_no":2,"account_code":"REVENUE","credit_minor_units":100}]}`
	request := httptest.NewRequest("POST", "/v1/accounting/journals", strings.NewReader(body))
	request.Header.Set("Authorization", "Bearer valid")
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 403 || repo.journals != 0 {
		t.Fatalf("status=%d journals=%d", response.Code, repo.journals)
	}
	request = httptest.NewRequest("POST", "/v1/accounting/journals", strings.NewReader(strings.Replace(body, `"other"`, `"franchise"`, 1)))
	request.Header.Set("Authorization", "Bearer valid")
	request.Header.Set("Content-Type", "application/json")
	response = httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 201 || repo.journals != 1 {
		t.Fatalf("status=%d journals=%d body=%s", response.Code, repo.journals, response.Body.String())
	}
}
````

### FILE: `db/migrations/0011_accounting_ledger.up.sql`
```yaml
block_id: "GO-ACCOUNTING:migration-up:v1"
operation: CREATE
provenance: AUTHORED
source: "local PostgreSQL schema governed by pinned Microsoft BCApps evidence"
license: "LicenseRef-Workspace-Owner"
sha256: "7dd6d1aeaa4610a1ce66e13309ac652c7126c6e63d947535e4026fc9a2754f92"
variables: []
secrets_allowed: false
```
````sql
create schema if not exists accounting;

create table accounting.account (
  tenant_id uuid not null,
  account_code text not null check (account_code ~ '^[A-Z0-9][A-Z0-9._-]{0,31}$'),
  display_name text not null check (length(display_name) between 2 and 120),
  account_type text not null check (account_type in ('asset','liability','equity','revenue','expense')),
  active boolean not null default true,
  created_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, account_code),
  foreign key (tenant_id) references platform.tenant (tenant_id)
);

create table accounting.period (
  tenant_id uuid not null,
  period_id text not null,
  starts_on date not null,
  ends_on date not null,
  status text not null check (status in ('open','closed')),
  version bigint not null check (version > 0),
  closed_at timestamptz,
  created_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, period_id),
  foreign key (tenant_id) references platform.tenant (tenant_id),
  unique (tenant_id, starts_on, ends_on),
  check (ends_on > starts_on),
  check ((status='open' and closed_at is null) or (status='closed' and closed_at is not null))
);

create table accounting.journal (
  tenant_id uuid not null,
  journal_id text not null,
  organization_id text not null,
  period_id text not null,
  source_type text not null check (source_type ~ '^[A-Z0-9][A-Z0-9._-]{0,31}$'),
  source_id text not null,
  currency text not null check (currency ~ '^[A-Z]{3}$'),
  posting_date date not null,
  status text not null check (status in ('draft','posted','reversed')),
  total_debit_minor_units bigint not null check (total_debit_minor_units > 0),
  total_credit_minor_units bigint not null check (total_credit_minor_units > 0),
  version bigint not null check (version > 0),
  reversal_of text,
  reversal_reason text,
  posted_by text,
  posted_at timestamptz,
  created_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, journal_id),
  foreign key (tenant_id, organization_id) references org.organization (tenant_id, organization_id),
  foreign key (tenant_id, period_id) references accounting.period (tenant_id, period_id),
  foreign key (tenant_id, reversal_of) references accounting.journal (tenant_id, journal_id),
  unique (tenant_id, source_type, source_id),
  check (total_debit_minor_units = total_credit_minor_units),
  check ((status='draft' and posted_at is null and posted_by is null) or (status in ('posted','reversed') and posted_at is not null and posted_by is not null)),
  check ((reversal_of is null and reversal_reason is null) or (reversal_of is not null and length(reversal_reason) between 3 and 500))
);

create unique index accounting_one_reversal_idx on accounting.journal(tenant_id,reversal_of) where reversal_of is not null;

create table accounting.journal_line (
  tenant_id uuid not null,
  journal_id text not null,
  line_no integer not null check (line_no > 0),
  account_code text not null,
  description text not null default '' check (length(description) <= 250),
  debit_minor_units bigint not null default 0 check (debit_minor_units >= 0),
  credit_minor_units bigint not null default 0 check (credit_minor_units >= 0),
  primary key (tenant_id,journal_id,line_no),
  foreign key (tenant_id,journal_id) references accounting.journal(tenant_id,journal_id),
  foreign key (tenant_id,account_code) references accounting.account(tenant_id,account_code),
  check ((debit_minor_units > 0) <> (credit_minor_units > 0))
);

create table accounting.register (
  tenant_id uuid not null,
  register_id text not null,
  journal_id text not null,
  organization_id text not null,
  period_id text not null,
  posted_by text not null,
  total_debit_minor_units bigint not null check (total_debit_minor_units > 0),
  total_credit_minor_units bigint not null check (total_credit_minor_units > 0),
  posted_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id,register_id),
  foreign key (tenant_id,journal_id) references accounting.journal(tenant_id,journal_id),
  unique (tenant_id,journal_id),
  check (total_debit_minor_units=total_credit_minor_units)
);

create table accounting.entry (
  tenant_id uuid not null,
  entry_id text not null,
  register_id text not null,
  journal_id text not null,
  line_no integer not null,
  organization_id text not null,
  period_id text not null,
  account_code text not null,
  currency text not null check (currency ~ '^[A-Z]{3}$'),
  posting_date date not null,
  description text not null,
  debit_minor_units bigint not null check (debit_minor_units >= 0),
  credit_minor_units bigint not null check (credit_minor_units >= 0),
  original_entry_id text,
  created_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id,entry_id),
  foreign key (tenant_id,register_id) references accounting.register(tenant_id,register_id),
  foreign key (tenant_id,journal_id,line_no) references accounting.journal_line(tenant_id,journal_id,line_no),
  foreign key (tenant_id,account_code) references accounting.account(tenant_id,account_code),
  unique (tenant_id,journal_id,line_no),
  check ((debit_minor_units > 0) <> (credit_minor_units > 0))
);

create unique index accounting_entry_one_reversal_idx on accounting.entry(tenant_id,original_entry_id) where original_entry_id is not null;

create index accounting_trial_balance_idx on accounting.entry(tenant_id,organization_id,period_id,account_code);

create or replace function accounting.prevent_posted_history_mutation() returns trigger language plpgsql as $$ begin raise exception 'immutable accounting history'; end; $$;
create trigger accounting_line_immutable before update or delete on accounting.journal_line for each row execute function accounting.prevent_posted_history_mutation();
create trigger accounting_register_immutable before update or delete on accounting.register for each row execute function accounting.prevent_posted_history_mutation();
create trigger accounting_entry_immutable before update or delete on accounting.entry for each row execute function accounting.prevent_posted_history_mutation();
````

### FILE: `db/migrations/0011_accounting_ledger.down.sql`
```yaml
block_id: "GO-ACCOUNTING:migration-down:v1"
operation: CREATE
provenance: AUTHORED
source: "local development rollback"
license: "LicenseRef-Workspace-Owner"
sha256: "3745aa3c0af7c437e83c7e4d0660bc5dfcf161df8358ad694059e818fbea2aad"
variables: []
secrets_allowed: false
```
````sql
drop trigger if exists accounting_entry_immutable on accounting.entry;
drop trigger if exists accounting_register_immutable on accounting.register;
drop trigger if exists accounting_line_immutable on accounting.journal_line;
drop function if exists accounting.prevent_posted_history_mutation();
drop table if exists accounting.entry;
drop table if exists accounting.register;
drop table if exists accounting.journal_line;
drop table if exists accounting.journal;
drop table if exists accounting.period;
drop table if exists accounting.account;
drop schema if exists accounting;
````

### FILE: `db/tests/0011_accounting_ledger.test.sql`
```yaml
block_id: "GO-ACCOUNTING:migration-test:v1"
operation: CREATE
provenance: AUTHORED
source: "local SQL invariant regression"
license: "LicenseRef-Workspace-Owner"
sha256: "1c597bf733e744e95f28a197c416ffffb1ae0a207143d13671ebce02e052e1c8"
variables: []
secrets_allowed: false
```
````sql
begin;
insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values('018f4d4a-7b36-7a21-8d10-2f4c54c29e00','accounting-test','Accounting','Accounting');
insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values('018f4d4a-7b36-7a21-8d10-2f4c54c29e00','franchise','franchise','Franchise','franchisee');
insert into accounting.account values('018f4d4a-7b36-7a21-8d10-2f4c54c29e00','CASH','Cash','asset',true,clock_timestamp()),('018f4d4a-7b36-7a21-8d10-2f4c54c29e00','REVENUE','Revenue','revenue',true,clock_timestamp());
insert into accounting.period(tenant_id,period_id,starts_on,ends_on,status,version)values('018f4d4a-7b36-7a21-8d10-2f4c54c29e00','2026-08','2026-08-01','2026-09-01','open',1);
insert into accounting.journal(tenant_id,journal_id,organization_id,period_id,source_type,source_id,currency,posting_date,status,total_debit_minor_units,total_credit_minor_units,version)values('018f4d4a-7b36-7a21-8d10-2f4c54c29e00','journal','franchise','2026-08','SALE','order-1','ARS','2026-08-15','draft',100,100,1);
insert into accounting.journal_line values('018f4d4a-7b36-7a21-8d10-2f4c54c29e00','journal',1,'CASH','cash',100,0),('018f4d4a-7b36-7a21-8d10-2f4c54c29e00','journal',2,'REVENUE','sale',0,100);
do $$ begin
  begin update accounting.journal_line set debit_minor_units=1 where tenant_id='018f4d4a-7b36-7a21-8d10-2f4c54c29e00'; raise exception 'line mutation accepted'; exception when raise_exception then if sqlerrm <> 'immutable accounting history' then raise; end if; end;
  begin insert into accounting.journal_line values('018f4d4a-7b36-7a21-8d10-2f4c54c29e00','journal',3,'CASH','bad',1,1); raise exception 'dual-sided line accepted'; exception when check_violation then null; end;
end $$;
rollback;
````

## 6. Configuration surface

No process secret or environment variable. Chart of accounts, periods and source mappings are project business configuration and require an accountable owner. Browser-originated account mappings are prohibited.

## 7. Dependency bill

| Dependency | Pin | Use | License | Source |
|---|---|---|---|---|
| Go | `1.26.7` | domain/HTTP | BSD-3-Clause | `go.dev` |
| PostgreSQL | `18.6` | ledger/locking/constraints | PostgreSQL | `postgresql.org` |
| pgx | `5.10.0` composed | transactions | MIT | `github.com/jackc/pgx` |
| Microsoft BCApps | commit `31a860b...4789` | exact architecture evidence only | MIT | `github.com/microsoft/BCApps` |

## 8. Apply order

Materialize after foundation and business-event owners; apply migration 0011 after 0010 and wire `AccountingModule`. Configure chart/period/source mappings before posting. Production rollback redeploys compatible code and preserves posted history; schema drop is development-only.

## 9. Verification

Require clean reconstruction/hash, Go format/test twice/vet/build, fresh PostgreSQL 18.6 migrations 0001–0011, SQL invariants, exact balance, source uniqueness, cross-scope rejection, one-winner posting, trial balance, equal/opposite reversal, closed-period rejection, down/up and global verifier. This does not prove fiscal localization, statutory books, tax, banking, consolidation, FX, audit opinion or target acceptance.

## 10. Reconstruction evidence

Recorded in `reconstruction_evidence/ENTERPRISE_ACCOUNTING_LEDGER_2026-08-30_V121.md`.

V402 composed delta: Connect the immutable FX conversion to the existing draft writer; preserve posting/reversal, explicit actor/account/period selection and no corporate authorship. Evidence FX_JOURNAL_CONNECTION_V402.md.
