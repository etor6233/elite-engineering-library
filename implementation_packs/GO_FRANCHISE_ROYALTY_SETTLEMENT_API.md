# Go Franchise Royalty Settlement API

## 1. Metadata

```yaml
pack_id: "GO-FRANCHISE-ROYALTY-SETTLEMENT-API"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Añade política versionada, devengos append-only desde pagos confirmados, liquidación concurrente, reversión compensatoria y conciliación explícita por franquicia."
stacks: ["Go 1.26.7", "PostgreSQL 18.6", "pgx 5.10.0"]
compatible_with: ["ELECTROMOBILITY-FRANCHISE-MODULES 0.1.x", "GO-COMMERCE-PRICING-PAYMENT-API 0.3.x", "GO-ELECTROMOBILITY-APPLICATION 1.x"]
incompatible_with: []
license_expression: "LicenseRef-Workspace-Owner AND MIT upstream evidence"
upstream_sources: ["https://github.com/microsoft/BCApps/tree/31a860b527f0dc72c7a44a255d7e7d403cfa4789", "https://www.postgresql.org/docs/18/", "https://go.dev"]
verified_at: "2026-08-30"
```

Los nueve bloques ejecutables son `AUTHORED`; ninguna línea Go/SQL se atribuye a Microsoft. El diseño está gobernado por fuentes públicas exactas de `microsoft/BCApps@31a860b527f0dc72c7a44a255d7e7d403cfa4789`: `SalespersonCommission.Report.al` (`4cff73d35c960415e6c77149b3a6e2533b0925f90fa1ba312a6de52f127b36fd`), `GLEntry.Table.al` (`08655960b539eb5afbea0c91f6189af07b4f83c68f7de08d5d590da6591fd9dd`), `ReversalEntry.Table.al` (`13e17b83c3a2d685060ac151a9f5ef0b156238036f034515ccb260d0b2eb835d`), `BankAccReconciliation.Table.al` (`2fac6264945000836bd27f0d48e74be6db1dd1bab1cf6675bc3ae32d06ab4b62`), `ShpfyPayout.Table.al` (`7503e221f987cecdc3401aec1d3b41e73e00e4f04dcc2f877a555b351ec75f31`) y `ShpfyPaymentTransaction.Table.al` (`e8f8cb2bf3d6b1b276bc2cc96f7cfbcfb6538e7bc2084499dc3b0f25a7c42db6`). El snapshot tiene ZIP SHA-256 `e3151b040df39cada83d41fcf5d9e6cdff1d8fddf226934c1ad21a2bdebc7210` y licencia MIT raíz SHA-256 `c2cfccb812fe482101a8f04597dfc5a9991a6b2748266c47ac91b6a5aae15383`.

## 2. Applicability

Use after the franchise agreement and commerce/payment owners. It records a contractual royalty subledger; it does not replace the payment provider ledger, general ledger, tax engine, statutory invoicing or bank adapter. A target project must supply approved contractual rates, currencies, periods, provider event mapping and reconciliation authority before production.

## 3. Architecture contract

The browser never supplies authoritative monetary bases or computed royalties. A protected command references an existing organization-scoped payment in the exact `captured` or `refunded` state and version; PostgreSQL resolves one active agreement and one effective policy, then computes integer minor units transactionally. Accrual, settlement lines and reconciliation evidence are immutable. Closing uses serializable transactions, advisory locking and optimistic versions; one concurrent close wins. Reversal posts equal-and-opposite lines and preserves the original. Reconciliation persists expected, actual and difference instead of silently marking mismatches as success. Every mutation writes its outbox event in the same transaction.

## 4. Exact file manifest

```text
CREATE internal/royalty/service.go
CREATE internal/royalty/service_test.go
CREATE internal/platform/postgres/royalty.go
CREATE internal/platform/postgres/royalty_integration_test.go
CREATE internal/platform/httpapi/royalty.go
CREATE internal/platform/httpapi/royalty_test.go
CREATE db/migrations/0010_franchise_royalty_settlement.up.sql
CREATE db/migrations/0010_franchise_royalty_settlement.down.sql
CREATE db/tests/0010_franchise_royalty_settlement.test.sql
```

## 5. Materialization blocks

### FILE: `internal/royalty/service.go`

```yaml
block_id: "GO-FRANCHISE-ROYALTY:service:v1"
operation: CREATE
provenance: AUTHORED
source: "local implementation governed by pinned Microsoft BCApps evidence"
license: "LicenseRef-Workspace-Owner"
sha256: "c3ae6d1aa7f780fa821235da443bebeef5cf6407f786f7240d768224a3680c10"
variables: []
secrets_allowed: false
```

````go
package royalty

import (
	"context"
	"errors"
	"regexp"
	"time"
)

var ErrConflict = errors.New("royalty conflict")
var ErrInvalid = errors.New("invalid royalty command")
var currencyPattern = regexp.MustCompile(`^[A-Z]{3}$`)

type Policy struct {
	ID              string     `json:"id"`
	AgreementID     string     `json:"agreement_id"`
	OrganizationID  string     `json:"organization_id"`
	Currency        string     `json:"currency"`
	RateBasisPoints int        `json:"rate_basis_points"`
	ValidFrom       time.Time  `json:"valid_from"`
	ValidUntil      *time.Time `json:"valid_until,omitempty"`
}

type PaymentEvent struct {
	ID                     string    `json:"id"`
	OrganizationID         string    `json:"organization_id"`
	PaymentAttemptID       string    `json:"payment_attempt_id"`
	PaymentExpectedVersion int64     `json:"payment_expected_version"`
	State                  string    `json:"state"`
	OccurredAt             time.Time `json:"occurred_at"`
}

type Accrual struct {
	ID                string    `json:"id"`
	PolicyID          string    `json:"policy_id"`
	AgreementID       string    `json:"agreement_id"`
	OrganizationID    string    `json:"organization_id"`
	PaymentAttemptID  string    `json:"payment_attempt_id"`
	SourceEventKey    string    `json:"source_event_key"`
	SourceState       string    `json:"source_state"`
	Currency          string    `json:"currency"`
	BasisMinorUnits   int64     `json:"basis_minor_units"`
	RoyaltyMinorUnits int64     `json:"royalty_minor_units"`
	OccurredAt        time.Time `json:"occurred_at"`
}

type Settlement struct {
	ID                 string    `json:"id"`
	OrganizationID     string    `json:"organization_id"`
	Currency           string    `json:"currency"`
	PeriodStart        time.Time `json:"period_start"`
	PeriodEnd          time.Time `json:"period_end"`
	Status             string    `json:"status"`
	ExpectedMinorUnits int64     `json:"expected_minor_units"`
	Version            int64     `json:"version"`
	ReversalOf         string    `json:"reversal_of,omitempty"`
}

type Reconciliation struct {
	ID                   string `json:"id"`
	SettlementID         string `json:"settlement_id"`
	ExternalReference    string `json:"external_reference"`
	ExpectedMinorUnits   int64  `json:"expected_minor_units"`
	ActualMinorUnits     int64  `json:"actual_minor_units"`
	DifferenceMinorUnits int64  `json:"difference_minor_units"`
	Status               string `json:"status"`
}

type Repository interface {
	CreatePolicy(context.Context, string, string, Policy) error
	AccruePayment(context.Context, string, string, PaymentEvent, string) (Accrual, error)
	OpenSettlement(context.Context, string, string, Settlement) error
	CloseSettlement(context.Context, string, string, string, string, int64) (Settlement, error)
	ReverseSettlement(context.Context, string, string, string, string, string, int64, string) (Settlement, error)
	Reconcile(context.Context, string, string, Reconciliation, string, string) (Reconciliation, error)
}

type IDGenerator interface{ New() string }

type Service struct {
	repository Repository
	ids        IDGenerator
}

func NewService(repository Repository, ids IDGenerator) *Service {
	return &Service{repository: repository, ids: ids}
}

func (s *Service) CreatePolicy(ctx context.Context, tenant string, value Policy) (Policy, error) {
	if tenant == "" || value.AgreementID == "" || value.OrganizationID == "" || !currencyPattern.MatchString(value.Currency) || value.RateBasisPoints < 1 || value.RateBasisPoints > 10000 || value.ValidFrom.IsZero() || (value.ValidUntil != nil && !value.ValidUntil.After(value.ValidFrom)) {
		return value, ErrInvalid
	}
	value.ID = s.ids.New()
	if err := s.repository.CreatePolicy(ctx, tenant, s.ids.New(), value); err != nil {
		return value, err
	}
	return value, nil
}

func (s *Service) AccruePayment(ctx context.Context, tenant, sourceEventKey string, value PaymentEvent) (Accrual, error) {
	if tenant == "" || sourceEventKey == "" || value.OrganizationID == "" || value.PaymentAttemptID == "" || value.PaymentExpectedVersion < 1 || (value.State != "captured" && value.State != "refunded") || value.OccurredAt.IsZero() {
		return Accrual{}, ErrInvalid
	}
	value.ID = s.ids.New()
	return s.repository.AccruePayment(ctx, tenant, sourceEventKey, value, s.ids.New())
}

func (s *Service) OpenSettlement(ctx context.Context, tenant string, value Settlement) (Settlement, error) {
	if tenant == "" || value.OrganizationID == "" || !currencyPattern.MatchString(value.Currency) || value.PeriodStart.IsZero() || !value.PeriodEnd.After(value.PeriodStart) {
		return value, ErrInvalid
	}
	value.ID = s.ids.New()
	value.Status = "draft"
	value.Version = 1
	if err := s.repository.OpenSettlement(ctx, tenant, s.ids.New(), value); err != nil {
		return value, err
	}
	return value, nil
}

func (s *Service) CloseSettlement(ctx context.Context, tenant, organization, id string, expectedVersion int64) (Settlement, error) {
	if tenant == "" || organization == "" || id == "" || expectedVersion < 1 {
		return Settlement{}, ErrInvalid
	}
	return s.repository.CloseSettlement(ctx, tenant, organization, id, s.ids.New(), expectedVersion)
}

func (s *Service) ReverseSettlement(ctx context.Context, tenant, organization, id string, expectedVersion int64, reason string) (Settlement, error) {
	if tenant == "" || organization == "" || id == "" || expectedVersion < 1 || len(reason) < 3 || len(reason) > 500 {
		return Settlement{}, ErrInvalid
	}
	return s.repository.ReverseSettlement(ctx, tenant, organization, id, s.ids.New(), s.ids.New(), expectedVersion, reason)
}

func (s *Service) Reconcile(ctx context.Context, tenant, organization, actor string, value Reconciliation) (Reconciliation, error) {
	if tenant == "" || organization == "" || actor == "" || value.SettlementID == "" || value.ExternalReference == "" || len(value.ExternalReference) > 250 {
		return value, ErrInvalid
	}
	value.ID = s.ids.New()
	return s.repository.Reconcile(ctx, tenant, organization, value, actor, s.ids.New())
}
````

### FILE: `internal/royalty/service_test.go`

```yaml
block_id: "GO-FRANCHISE-ROYALTY:service-test:v1"
operation: CREATE
provenance: AUTHORED
source: "local executable regression"
license: "LicenseRef-Workspace-Owner"
sha256: "d05d15b00a7e1cb5c01394d85cdc2020194e4e7bf4eaca2e67ec666d5e86ba95"
variables: []
secrets_allowed: false
```

````go
package royalty

import (
	"context"
	"errors"
	"testing"
	"time"
)

type fakeRepo struct {
	policy     Policy
	event      PaymentEvent
	settlement Settlement
}

func (f *fakeRepo) CreatePolicy(_ context.Context, _, _ string, value Policy) error {
	f.policy = value
	return nil
}
func (f *fakeRepo) AccruePayment(_ context.Context, _, _ string, value PaymentEvent, _ string) (Accrual, error) {
	f.event = value
	return Accrual{ID: value.ID, SourceState: value.State}, nil
}
func (f *fakeRepo) OpenSettlement(_ context.Context, _, _ string, value Settlement) error {
	f.settlement = value
	return nil
}
func (f *fakeRepo) CloseSettlement(context.Context, string, string, string, string, int64) (Settlement, error) {
	return Settlement{Status: "closed"}, nil
}
func (f *fakeRepo) ReverseSettlement(context.Context, string, string, string, string, string, int64, string) (Settlement, error) {
	return Settlement{Status: "closed", ReversalOf: "s"}, nil
}
func (f *fakeRepo) Reconcile(_ context.Context, _, _ string, value Reconciliation, _, _ string) (Reconciliation, error) {
	value.Status = "matched"
	return value, nil
}

type seqIDs struct{ n int }

func (s *seqIDs) New() string { s.n++; return "id-" + string(rune('0'+s.n)) }

func TestServiceValidatesAndDelegates(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	repo := &fakeRepo{}
	service := NewService(repo, &seqIDs{})
	policy, err := service.CreatePolicy(context.Background(), "tenant", Policy{AgreementID: "agreement", OrganizationID: "org", Currency: "ARS", RateBasisPoints: 650, ValidFrom: now})
	if err != nil || policy.ID == "" || repo.policy.RateBasisPoints != 650 {
		t.Fatalf("policy=%+v err=%v", policy, err)
	}
	accrual, err := service.AccruePayment(context.Background(), "tenant", "provider:event-1", PaymentEvent{OrganizationID: "org", PaymentAttemptID: "pay", PaymentExpectedVersion: 3, State: "captured", OccurredAt: now})
	if err != nil || accrual.SourceState != "captured" || repo.event.ID == "" {
		t.Fatalf("accrual=%+v err=%v", accrual, err)
	}
	settlement, err := service.OpenSettlement(context.Background(), "tenant", Settlement{OrganizationID: "org", Currency: "ARS", PeriodStart: now, PeriodEnd: now.Add(24 * time.Hour)})
	if err != nil || settlement.Status != "draft" || settlement.Version != 1 {
		t.Fatalf("settlement=%+v err=%v", settlement, err)
	}
}

func TestServiceRejectsUntrustedFinancialShapes(t *testing.T) {
	now := time.Now().UTC()
	service := NewService(&fakeRepo{}, &seqIDs{})
	cases := []error{}
	_, err := service.CreatePolicy(context.Background(), "tenant", Policy{AgreementID: "a", OrganizationID: "o", Currency: "ars", RateBasisPoints: 650, ValidFrom: now})
	cases = append(cases, err)
	_, err = service.AccruePayment(context.Background(), "tenant", "", PaymentEvent{OrganizationID: "o", PaymentAttemptID: "p", PaymentExpectedVersion: 1, State: "captured", OccurredAt: now})
	cases = append(cases, err)
	_, err = service.AccruePayment(context.Background(), "tenant", "e", PaymentEvent{OrganizationID: "o", PaymentAttemptID: "p", PaymentExpectedVersion: 1, State: "pending", OccurredAt: now})
	cases = append(cases, err)
	_, err = service.OpenSettlement(context.Background(), "tenant", Settlement{OrganizationID: "o", Currency: "ARS", PeriodStart: now, PeriodEnd: now})
	cases = append(cases, err)
	_, err = service.ReverseSettlement(context.Background(), "tenant", "o", "s", 1, "")
	cases = append(cases, err)
	for i, got := range cases {
		if !errors.Is(got, ErrInvalid) {
			t.Fatalf("case %d err=%v", i, got)
		}
	}
}
````

### FILE: `internal/platform/postgres/royalty.go`

```yaml
block_id: "GO-FRANCHISE-ROYALTY:postgres:v1"
operation: CREATE
provenance: AUTHORED
source: "local implementation governed by pinned Microsoft BCApps and PostgreSQL evidence"
license: "LicenseRef-Workspace-Owner"
sha256: "8a65ba89abe094c5e2dac3ba95d82dc92238d8ed4bcbbf3fbdf11b108d769853"
variables: []
secrets_allowed: false
```

````go
package postgres

import (
	"context"
	"elite.local/enterprise/internal/royalty"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Royalty struct{ pool *pgxpool.Pool }

func NewRoyalty(pool *pgxpool.Pool) *Royalty { return &Royalty{pool: pool} }

func royaltyConflict(err error) error {
	if err == nil {
		return nil
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && (pgErr.Code == "23503" || pgErr.Code == "23505" || pgErr.Code == "23514" || pgErr.Code == "40001") {
		return royalty.ErrConflict
	}
	return err
}

func (r *Royalty) CreatePolicy(ctx context.Context, tenant, eventID string, value royalty.Policy) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if _, err = tx.Exec(ctx, `select pg_advisory_xact_lock(hashtextextended($1,0))`, tenant+":"+value.AgreementID+":"+value.Currency); err != nil {
		return err
	}
	var agreement string
	err = tx.QueryRow(ctx, `select agreement_id from franchise.agreement where tenant_id=$1 and agreement_id=$2 and franchise_organization_id=$3 and status='active' and starts_on <= ($4::timestamptz at time zone 'UTC')::date and (ends_on is null or ends_on > ($4::timestamptz at time zone 'UTC')::date) for update`, tenant, value.AgreementID, value.OrganizationID, value.ValidFrom).Scan(&agreement)
	if errors.Is(err, pgx.ErrNoRows) {
		return royalty.ErrConflict
	}
	if err != nil {
		return err
	}
	var overlap bool
	err = tx.QueryRow(ctx, `select exists(select 1 from royalty.policy where tenant_id=$1 and agreement_id=$2 and currency=$3 and valid_from < coalesce($5::timestamptz,'infinity'::timestamptz) and coalesce(valid_until,'infinity'::timestamptz) > $4)`, tenant, value.AgreementID, value.Currency, value.ValidFrom, value.ValidUntil).Scan(&overlap)
	if err != nil {
		return err
	}
	if overlap {
		return royalty.ErrConflict
	}
	_, err = tx.Exec(ctx, `insert into royalty.policy(tenant_id,policy_id,agreement_id,franchise_organization_id,currency,rate_basis_points,valid_from,valid_until) values($1,$2,$3,$4,$5,$6,$7,$8)`, tenant, value.ID, value.AgreementID, value.OrganizationID, value.Currency, value.RateBasisPoints, value.ValidFrom, value.ValidUntil)
	if err != nil {
		return royaltyConflict(err)
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload) values($1,$2,'royalty-policy',$3,1,'royalty-policy.created',1,clock_timestamp(),jsonb_build_object('agreement_id',$4::text,'organization_id',$5::text,'currency',$6::text,'rate_basis_points',$7::integer))`, tenant, eventID, value.ID, value.AgreementID, value.OrganizationID, value.Currency, value.RateBasisPoints)
	if err != nil {
		return royaltyConflict(err)
	}
	return royaltyConflict(tx.Commit(ctx))
}

func (r *Royalty) AccruePayment(ctx context.Context, tenant, sourceEventKey string, value royalty.PaymentEvent, eventID string) (royalty.Accrual, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return royalty.Accrual{}, err
	}
	defer tx.Rollback(ctx)
	var amount int64
	var currency string
	err = tx.QueryRow(ctx, `select p.amount_minor_units,p.currency from payment.payment_attempt p join sales.customer_order o on o.tenant_id=p.tenant_id and o.order_id=p.order_id where p.tenant_id=$1 and p.payment_attempt_id=$2 and p.state=$3 and p.version=$4 and o.organization_id=$5 for update of p`, tenant, value.PaymentAttemptID, value.State, value.PaymentExpectedVersion, value.OrganizationID).Scan(&amount, &currency)
	if errors.Is(err, pgx.ErrNoRows) {
		return royalty.Accrual{}, royalty.ErrConflict
	}
	if err != nil {
		return royalty.Accrual{}, err
	}
	rows, err := tx.Query(ctx, `select p.policy_id,p.agreement_id,p.rate_basis_points from royalty.policy p join franchise.agreement a on a.tenant_id=p.tenant_id and a.agreement_id=p.agreement_id where p.tenant_id=$1 and p.franchise_organization_id=$2 and p.currency=$3 and p.valid_from <= $4 and (p.valid_until is null or p.valid_until > $4) and a.status='active' and a.starts_on <= ($4::timestamptz at time zone 'UTC')::date and (a.ends_on is null or a.ends_on > ($4::timestamptz at time zone 'UTC')::date) order by p.valid_from desc limit 2`, tenant, value.OrganizationID, currency, value.OccurredAt)
	if err != nil {
		return royalty.Accrual{}, err
	}
	defer rows.Close()
	type resolved struct {
		policy, agreement string
		rate              int
	}
	found := []resolved{}
	for rows.Next() {
		var v resolved
		if err = rows.Scan(&v.policy, &v.agreement, &v.rate); err != nil {
			return royalty.Accrual{}, err
		}
		found = append(found, v)
	}
	if err = rows.Err(); err != nil {
		return royalty.Accrual{}, err
	}
	if len(found) != 1 {
		return royalty.Accrual{}, royalty.ErrConflict
	}
	var royaltyAmount int64
	err = tx.QueryRow(ctx, `select round(($1::numeric*$2::numeric)/10000)::bigint`, amount, found[0].rate).Scan(&royaltyAmount)
	if err != nil {
		return royalty.Accrual{}, err
	}
	if royaltyAmount == 0 {
		return royalty.Accrual{}, royalty.ErrConflict
	}
	if value.State == "refunded" {
		royaltyAmount = -royaltyAmount
	}
	result := royalty.Accrual{ID: value.ID, PolicyID: found[0].policy, AgreementID: found[0].agreement, OrganizationID: value.OrganizationID, PaymentAttemptID: value.PaymentAttemptID, SourceEventKey: sourceEventKey, SourceState: value.State, Currency: currency, BasisMinorUnits: amount, RoyaltyMinorUnits: royaltyAmount, OccurredAt: value.OccurredAt}
	_, err = tx.Exec(ctx, `insert into royalty.accrual(tenant_id,accrual_id,policy_id,agreement_id,franchise_organization_id,payment_attempt_id,source_event_key,source_state,currency,basis_minor_units,royalty_minor_units,occurred_at) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`, tenant, result.ID, result.PolicyID, result.AgreementID, result.OrganizationID, result.PaymentAttemptID, result.SourceEventKey, result.SourceState, result.Currency, result.BasisMinorUnits, result.RoyaltyMinorUnits, result.OccurredAt)
	if err != nil {
		return royalty.Accrual{}, royaltyConflict(err)
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload) values($1,$2,'royalty-accrual',$3,1,'royalty-accrual.posted',1,clock_timestamp(),jsonb_build_object('payment_attempt_id',$4::text,'organization_id',$5::text,'currency',$6::text,'royalty_minor_units',$7::bigint,'source_state',$8::text))`, tenant, eventID, result.ID, result.PaymentAttemptID, result.OrganizationID, result.Currency, result.RoyaltyMinorUnits, result.SourceState)
	if err != nil {
		return royalty.Accrual{}, royaltyConflict(err)
	}
	if err = tx.Commit(ctx); err != nil {
		return royalty.Accrual{}, royaltyConflict(err)
	}
	return result, nil
}

func (r *Royalty) OpenSettlement(ctx context.Context, tenant, eventID string, value royalty.Settlement) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, `insert into royalty.settlement_run(tenant_id,settlement_id,franchise_organization_id,currency,period_start,period_end,status,expected_minor_units,version) select $1,$2,$3,$4,$5,$6,'draft',0,1 where exists(select 1 from franchise.agreement where tenant_id=$1 and franchise_organization_id=$3 and status='active')`, tenant, value.ID, value.OrganizationID, value.Currency, value.PeriodStart, value.PeriodEnd)
	if err != nil {
		return royaltyConflict(err)
	}
	var exists bool
	if err = tx.QueryRow(ctx, `select exists(select 1 from royalty.settlement_run where tenant_id=$1 and settlement_id=$2)`, tenant, value.ID).Scan(&exists); err != nil {
		return err
	}
	if !exists {
		return royalty.ErrConflict
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload) values($1,$2,'royalty-settlement',$3,1,'royalty-settlement.opened',1,clock_timestamp(),jsonb_build_object('organization_id',$4::text,'currency',$5::text))`, tenant, eventID, value.ID, value.OrganizationID, value.Currency)
	if err != nil {
		return royaltyConflict(err)
	}
	return royaltyConflict(tx.Commit(ctx))
}

func (r *Royalty) CloseSettlement(ctx context.Context, tenant, organization, id, eventID string, expectedVersion int64) (royalty.Settlement, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return royalty.Settlement{}, err
	}
	defer tx.Rollback(ctx)
	var value royalty.Settlement
	err = tx.QueryRow(ctx, `select settlement_id,franchise_organization_id,currency,period_start,period_end,status,expected_minor_units,version,coalesce(reversal_of,'') from royalty.settlement_run where tenant_id=$1 and settlement_id=$2 and franchise_organization_id=$3 and status='draft' and version=$4 for update`, tenant, id, organization, expectedVersion).Scan(&value.ID, &value.OrganizationID, &value.Currency, &value.PeriodStart, &value.PeriodEnd, &value.Status, &value.ExpectedMinorUnits, &value.Version, &value.ReversalOf)
	if errors.Is(err, pgx.ErrNoRows) {
		return royalty.Settlement{}, royalty.ErrConflict
	}
	if err != nil {
		return royalty.Settlement{}, err
	}
	if _, err = tx.Exec(ctx, `select pg_advisory_xact_lock(hashtextextended($1,0))`, tenant+":"+organization+":"+value.Currency); err != nil {
		return royalty.Settlement{}, err
	}
	result, err := tx.Exec(ctx, `insert into royalty.settlement_line(tenant_id,settlement_id,line_id,accrual_id,amount_minor_units) select a.tenant_id,$2,a.accrual_id,a.accrual_id,a.royalty_minor_units from royalty.accrual a where a.tenant_id=$1 and a.franchise_organization_id=$3 and a.currency=$4 and a.occurred_at >= $5 and a.occurred_at < $6 and not exists(select 1 from royalty.settlement_line l where l.tenant_id=a.tenant_id and l.accrual_id=a.accrual_id and l.original_line_id is null) order by a.occurred_at,a.accrual_id`, tenant, id, organization, value.Currency, value.PeriodStart, value.PeriodEnd)
	if err != nil {
		return royalty.Settlement{}, royaltyConflict(err)
	}
	if result.RowsAffected() == 0 {
		return royalty.Settlement{}, royalty.ErrConflict
	}
	err = tx.QueryRow(ctx, `update royalty.settlement_run s set status='closed',expected_minor_units=x.total,version=version+1,closed_at=clock_timestamp() from (select sum(amount_minor_units)::bigint total from royalty.settlement_line where tenant_id=$1 and settlement_id=$2) x where s.tenant_id=$1 and s.settlement_id=$2 and s.status='draft' and s.version=$3 returning s.expected_minor_units,s.version,s.status`, tenant, id, expectedVersion).Scan(&value.ExpectedMinorUnits, &value.Version, &value.Status)
	if errors.Is(err, pgx.ErrNoRows) {
		return royalty.Settlement{}, royalty.ErrConflict
	}
	if err != nil {
		return royalty.Settlement{}, err
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload) values($1,$2,'royalty-settlement',$3,$4,'royalty-settlement.closed',1,clock_timestamp(),jsonb_build_object('organization_id',$5::text,'currency',$6::text,'expected_minor_units',$7::bigint))`, tenant, eventID, id, value.Version, organization, value.Currency, value.ExpectedMinorUnits)
	if err != nil {
		return royalty.Settlement{}, royaltyConflict(err)
	}
	if err = tx.Commit(ctx); err != nil {
		return royalty.Settlement{}, royaltyConflict(err)
	}
	return value, nil
}

func (r *Royalty) ReverseSettlement(ctx context.Context, tenant, organization, id, reversalID, eventID string, expectedVersion int64, reason string) (royalty.Settlement, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return royalty.Settlement{}, err
	}
	defer tx.Rollback(ctx)
	var value royalty.Settlement
	err = tx.QueryRow(ctx, `select settlement_id,franchise_organization_id,currency,period_start,period_end,status,expected_minor_units,version from royalty.settlement_run where tenant_id=$1 and settlement_id=$2 and franchise_organization_id=$3 and status='closed' and version=$4 and reversal_of is null for update`, tenant, id, organization, expectedVersion).Scan(&value.ID, &value.OrganizationID, &value.Currency, &value.PeriodStart, &value.PeriodEnd, &value.Status, &value.ExpectedMinorUnits, &value.Version)
	if errors.Is(err, pgx.ErrNoRows) {
		return royalty.Settlement{}, royalty.ErrConflict
	}
	if err != nil {
		return royalty.Settlement{}, err
	}
	_, err = tx.Exec(ctx, `insert into royalty.settlement_run(tenant_id,settlement_id,franchise_organization_id,currency,period_start,period_end,status,expected_minor_units,version,reversal_of,reversal_reason,closed_at) values($1,$2,$3,$4,$5,$6,'closed',$7,1,$8,$9,clock_timestamp())`, tenant, reversalID, organization, value.Currency, value.PeriodStart, value.PeriodEnd, -value.ExpectedMinorUnits, id, reason)
	if err != nil {
		return royalty.Settlement{}, royaltyConflict(err)
	}
	result, err := tx.Exec(ctx, `insert into royalty.settlement_line(tenant_id,settlement_id,line_id,accrual_id,amount_minor_units,original_line_id) select tenant_id,$3,$3||':'||line_id,accrual_id,-amount_minor_units,$2||':'||line_id from royalty.settlement_line where tenant_id=$1 and settlement_id=$2`, tenant, id, reversalID)
	if err != nil {
		return royalty.Settlement{}, royaltyConflict(err)
	}
	if result.RowsAffected() == 0 {
		return royalty.Settlement{}, royalty.ErrConflict
	}
	result, err = tx.Exec(ctx, `update royalty.settlement_run set status='reversed',version=version+1 where tenant_id=$1 and settlement_id=$2 and status='closed' and version=$3`, tenant, id, expectedVersion)
	if err != nil {
		return royalty.Settlement{}, err
	}
	if result.RowsAffected() != 1 {
		return royalty.Settlement{}, royalty.ErrConflict
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload) values($1,$2,'royalty-settlement',$3,1,'royalty-settlement.reversal-posted',1,clock_timestamp(),jsonb_build_object('reversal_of',$4::text,'organization_id',$5::text,'expected_minor_units',$6::bigint,'reason',$7::text))`, tenant, eventID, reversalID, id, organization, -value.ExpectedMinorUnits, reason)
	if err != nil {
		return royalty.Settlement{}, royaltyConflict(err)
	}
	if err = tx.Commit(ctx); err != nil {
		return royalty.Settlement{}, royaltyConflict(err)
	}
	return royalty.Settlement{ID: reversalID, OrganizationID: organization, Currency: value.Currency, PeriodStart: value.PeriodStart, PeriodEnd: value.PeriodEnd, Status: "closed", ExpectedMinorUnits: -value.ExpectedMinorUnits, Version: 1, ReversalOf: id}, nil
}

func (r *Royalty) Reconcile(ctx context.Context, tenant, organization string, value royalty.Reconciliation, actor, eventID string) (royalty.Reconciliation, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return value, err
	}
	defer tx.Rollback(ctx)
	err = tx.QueryRow(ctx, `select expected_minor_units from royalty.settlement_run where tenant_id=$1 and settlement_id=$2 and franchise_organization_id=$3 and status='closed' for update`, tenant, value.SettlementID, organization).Scan(&value.ExpectedMinorUnits)
	if errors.Is(err, pgx.ErrNoRows) {
		return value, royalty.ErrConflict
	}
	if err != nil {
		return value, err
	}
	value.DifferenceMinorUnits = value.ActualMinorUnits - value.ExpectedMinorUnits
	value.Status = "mismatch"
	if value.DifferenceMinorUnits == 0 {
		value.Status = "matched"
	}
	_, err = tx.Exec(ctx, `insert into royalty.reconciliation(tenant_id,reconciliation_id,settlement_id,external_reference,expected_minor_units,actual_minor_units,difference_minor_units,status,recorded_by) values($1,$2,$3,$4,$5,$6,$7,$8,$9)`, tenant, value.ID, value.SettlementID, value.ExternalReference, value.ExpectedMinorUnits, value.ActualMinorUnits, value.DifferenceMinorUnits, value.Status, actor)
	if err != nil {
		return value, royaltyConflict(err)
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload) values($1,$2,'royalty-reconciliation',$3,1,'royalty-reconciliation.recorded',1,clock_timestamp(),jsonb_build_object('settlement_id',$4::text,'status',$5::text,'difference_minor_units',$6::bigint))`, tenant, eventID, value.ID, value.SettlementID, value.Status, value.DifferenceMinorUnits)
	if err != nil {
		return value, royaltyConflict(err)
	}
	if err = tx.Commit(ctx); err != nil {
		return value, royaltyConflict(err)
	}
	return value, nil
}
````

### FILE: `internal/platform/postgres/royalty_integration_test.go`

```yaml
block_id: "GO-FRANCHISE-ROYALTY:postgres-test:v1"
operation: CREATE
provenance: AUTHORED
source: "local PostgreSQL integration and concurrency regression"
license: "LicenseRef-Workspace-Owner"
sha256: "d16d3821f44226c5ddc367cf48a9eb30e6c9057e7d0c9e3a5001fd744607aef1"
variables: []
secrets_allowed: false
```

````go
package postgres

import (
	"context"
	"crypto/rand"
	"elite.local/enterprise/internal/royalty"
	"errors"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestRoyaltyLedgerSettlementConcurrencyReversalAndReconciliation(t *testing.T) {
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
		tx, cleanupErr := pool.Begin(ctx)
		if cleanupErr != nil {
			t.Errorf("cleanup begin: %v", cleanupErr)
			return
		}
		defer tx.Rollback(ctx)
		if _, cleanupErr = tx.Exec(ctx, `set local session_replication_role = replica`); cleanupErr != nil {
			t.Errorf("cleanup role: %v", cleanupErr)
			return
		}
		for _, query := range []string{`delete from platform.outbox_event where tenant_id=$1`, `delete from royalty.reconciliation where tenant_id=$1`, `delete from royalty.settlement_line where tenant_id=$1`, `delete from royalty.settlement_run where tenant_id=$1`, `delete from royalty.accrual where tenant_id=$1`, `delete from royalty.policy where tenant_id=$1`, `delete from payment.payment_attempt where tenant_id=$1`, `delete from sales.customer_order where tenant_id=$1`, `delete from franchise.agreement where tenant_id=$1`, `delete from org.organization where tenant_id=$1`, `delete from platform.tenant where tenant_id=$1`} {
			if _, cleanupErr = tx.Exec(ctx, query, tenant); cleanupErr != nil {
				t.Errorf("cleanup query: %v", cleanupErr)
				return
			}
		}
		if cleanupErr = tx.Commit(ctx); cleanupErr != nil {
			t.Errorf("cleanup commit: %v", cleanupErr)
		}
	}()
	fixtures := []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1::uuid,'royalty-'||substring($1::text,1,8),'Royalty','Royalty')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'franchise','franchise','Franchise','franchisee'),($1,'other','other','Other','franchisee')`,
		`insert into franchise.agreement(tenant_id,agreement_id,franchise_organization_id,territory_code,terms_version,starts_on,status)values($1,'agreement','franchise','AR-BUE-ROYALTY','v1','2026-01-01','active')`,
		`insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version)values($1,'order','franchise','customer','confirmed','ARS',10001,1)`,
		`insert into payment.payment_attempt(tenant_id,payment_attempt_id,order_id,provider_code,provider_reference,idempotency_key,state,currency,amount_minor_units,version)values($1,'payment','order','provider','provider-payment','idem-payment','captured','ARS',10001,3)`,
	}
	for _, q := range fixtures {
		if _, err = pool.Exec(ctx, q, tenant); err != nil {
			t.Fatal(err)
		}
	}
	repo := NewRoyalty(pool)
	from := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	until := from.AddDate(0, 1, 0)
	if err = repo.CreatePolicy(ctx, tenant, "018f4d4a-7b36-7a21-8d10-2f4c54c29c01", royalty.Policy{ID: "policy", AgreementID: "agreement", OrganizationID: "franchise", Currency: "ARS", RateBasisPoints: 650, ValidFrom: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)}); err != nil {
		t.Fatal(err)
	}
	if err = repo.CreatePolicy(ctx, tenant, "unused", royalty.Policy{ID: "overlap", AgreementID: "agreement", OrganizationID: "franchise", Currency: "ARS", RateBasisPoints: 700, ValidFrom: time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)}); !errors.Is(err, royalty.ErrConflict) {
		t.Fatalf("overlapping policy accepted: %v", err)
	}
	accrual, err := repo.AccruePayment(ctx, tenant, "provider:capture-1", royalty.PaymentEvent{ID: "accrual-capture", OrganizationID: "franchise", PaymentAttemptID: "payment", PaymentExpectedVersion: 3, State: "captured", OccurredAt: from.Add(time.Hour)}, "018f4d4a-7b36-7a21-8d10-2f4c54c29c02")
	if err != nil || accrual.RoyaltyMinorUnits != 650 || accrual.BasisMinorUnits != 10001 {
		t.Fatalf("accrual=%+v err=%v", accrual, err)
	}
	if _, err = repo.AccruePayment(ctx, tenant, "provider:capture-1", royalty.PaymentEvent{ID: "accrual-duplicate", OrganizationID: "franchise", PaymentAttemptID: "payment", PaymentExpectedVersion: 3, State: "captured", OccurredAt: from.Add(time.Hour)}, "unused"); !errors.Is(err, royalty.ErrConflict) {
		t.Fatalf("duplicate accrual accepted: %v", err)
	}
	if _, err = repo.AccruePayment(ctx, tenant, "provider:cross-scope", royalty.PaymentEvent{ID: "accrual-cross", OrganizationID: "other", PaymentAttemptID: "payment", PaymentExpectedVersion: 3, State: "captured", OccurredAt: from.Add(time.Hour)}, "unused"); !errors.Is(err, royalty.ErrConflict) {
		t.Fatalf("cross-scope accrual accepted: %v", err)
	}
	settlement := royalty.Settlement{ID: "settlement", OrganizationID: "franchise", Currency: "ARS", PeriodStart: from, PeriodEnd: until, Status: "draft", Version: 1}
	if err = repo.OpenSettlement(ctx, tenant, "018f4d4a-7b36-7a21-8d10-2f4c54c29c03", settlement); err != nil {
		t.Fatal(err)
	}
	type closeResult struct {
		value royalty.Settlement
		err   error
	}
	results := make(chan closeResult, 2)
	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-start
			v, e := repo.CloseSettlement(ctx, tenant, "franchise", "settlement", []string{"018f4d4a-7b36-7a21-8d10-2f4c54c29c04", "018f4d4a-7b36-7a21-8d10-2f4c54c29c05"}[i], 1)
			results <- closeResult{v, e}
		}(i)
	}
	close(start)
	wg.Wait()
	close(results)
	successes, conflicts := 0, 0
	var closed royalty.Settlement
	for result := range results {
		if result.err == nil {
			successes++
			closed = result.value
		} else if errors.Is(result.err, royalty.ErrConflict) {
			conflicts++
		} else {
			t.Fatal(result.err)
		}
	}
	if successes != 1 || conflicts != 1 || closed.ExpectedMinorUnits != 650 || closed.Version != 2 {
		t.Fatalf("success=%d conflicts=%d closed=%+v", successes, conflicts, closed)
	}
	reversal, err := repo.ReverseSettlement(ctx, tenant, "franchise", "settlement", "settlement-reversal", "018f4d4a-7b36-7a21-8d10-2f4c54c29c06", 2, "contract correction")
	if err != nil || reversal.ExpectedMinorUnits != -650 || reversal.ReversalOf != "settlement" {
		t.Fatalf("reversal=%+v err=%v", reversal, err)
	}
	if _, err = repo.ReverseSettlement(ctx, tenant, "franchise", "settlement", "second-reversal", "unused", 3, "duplicate reversal"); !errors.Is(err, royalty.ErrConflict) {
		t.Fatalf("second reversal accepted: %v", err)
	}
	reconciliation, err := repo.Reconcile(ctx, tenant, "franchise", royalty.Reconciliation{ID: "reconciliation", SettlementID: "settlement-reversal", ExternalReference: "bank-line-1", ActualMinorUnits: -649}, "controller", "018f4d4a-7b36-7a21-8d10-2f4c54c29c07")
	if err != nil || reconciliation.Status != "mismatch" || reconciliation.DifferenceMinorUnits != 1 {
		t.Fatalf("reconciliation=%+v err=%v", reconciliation, err)
	}
	var lineTotal int64
	var lines int
	if err = pool.QueryRow(ctx, `select coalesce(sum(amount_minor_units),0),count(*) from royalty.settlement_line where tenant_id=$1 and settlement_id in ('settlement','settlement-reversal')`, tenant).Scan(&lineTotal, &lines); err != nil {
		t.Fatal(err)
	}
	if lineTotal != 0 || lines != 2 {
		t.Fatalf("compensating history total=%d lines=%d", lineTotal, lines)
	}
}
````

### FILE: `internal/platform/httpapi/royalty.go`

```yaml
block_id: "GO-FRANCHISE-ROYALTY:http:v1"
operation: CREATE
provenance: AUTHORED
source: "local protected HTTP boundary"
license: "LicenseRef-Workspace-Owner"
sha256: "7eb3a70511bb1f21a1957144b39b6649fc39074e58ff4216d9339275f60d7de4"
variables: []
secrets_allowed: false
```

````go
package httpapi

import (
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/royalty"
	"errors"
	"net/http"
)

type RoyaltyModule struct{ Service *royalty.Service }

func (m RoyaltyModule) Register(mux *http.ServeMux, verifier identity.Verifier) {
	api := royaltyAPI{service: m.Service, verifier: verifier}
	mux.HandleFunc("POST /v1/royalties/policies", api.createPolicy)
	mux.HandleFunc("POST /v1/royalties/accruals/from-payment", api.accruePayment)
	mux.HandleFunc("POST /v1/royalties/settlements", api.openSettlement)
	mux.HandleFunc("POST /v1/royalties/settlements/{id}/close", api.closeSettlement)
	mux.HandleFunc("POST /v1/royalties/settlements/{id}/reverse", api.reverseSettlement)
	mux.HandleFunc("POST /v1/royalties/settlements/{id}/reconciliations", api.reconcile)
}

type royaltyAPI struct {
	service  *royalty.Service
	verifier identity.Verifier
}

func (a royaltyAPI) auth(w http.ResponseWriter, r *http.Request, permission string) (identity.Principal, bool) {
	p, err := authenticate(r.Context(), r.Header.Get("Authorization"), a.verifier)
	if err != nil {
		writeProblem(w, 401, "UNAUTHENTICATED", "a valid bearer token is required")
		return p, false
	}
	if !p.Allowed(permission) {
		writeProblem(w, 403, "FORBIDDEN", permission+" permission is required")
		return p, false
	}
	if r.Header.Get("Content-Type") != "application/json" {
		writeProblem(w, 415, "UNSUPPORTED_MEDIA_TYPE", "Content-Type must be application/json")
		return p, false
	}
	return p, true
}
func royaltyOrganization(w http.ResponseWriter, p identity.Principal, organization string) bool {
	if !p.AllowedOrganization(organization) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return false
	}
	return true
}

func (a royaltyAPI) createPolicy(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "royalty:policy")
	if !ok {
		return
	}
	var input royalty.Policy
	if !decodeStrict(w, r, &input) {
		return
	}
	if !royaltyOrganization(w, p, input.OrganizationID) {
		return
	}
	value, err := a.service.CreatePolicy(r.Context(), p.TenantID, input)
	if err != nil {
		writeRoyaltyResult(w, err)
		return
	}
	writeJSON(w, 201, value)
}
func (a royaltyAPI) accruePayment(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "royalty:post")
	if !ok {
		return
	}
	var input royalty.PaymentEvent
	if !decodeStrict(w, r, &input) {
		return
	}
	if !royaltyOrganization(w, p, input.OrganizationID) {
		return
	}
	value, err := a.service.AccruePayment(r.Context(), p.TenantID, r.Header.Get("Idempotency-Key"), input)
	if err != nil {
		writeRoyaltyResult(w, err)
		return
	}
	writeJSON(w, 201, value)
}
func (a royaltyAPI) openSettlement(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "royalty:settle")
	if !ok {
		return
	}
	var input royalty.Settlement
	if !decodeStrict(w, r, &input) {
		return
	}
	if !royaltyOrganization(w, p, input.OrganizationID) {
		return
	}
	value, err := a.service.OpenSettlement(r.Context(), p.TenantID, input)
	if err != nil {
		writeRoyaltyResult(w, err)
		return
	}
	writeJSON(w, 201, value)
}
func (a royaltyAPI) closeSettlement(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "royalty:settle")
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
	if !royaltyOrganization(w, p, input.OrganizationID) {
		return
	}
	value, err := a.service.CloseSettlement(r.Context(), p.TenantID, input.OrganizationID, r.PathValue("id"), input.ExpectedVersion)
	if err != nil {
		writeRoyaltyResult(w, err)
		return
	}
	writeJSON(w, 200, value)
}
func (a royaltyAPI) reverseSettlement(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "royalty:reverse")
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
	if !royaltyOrganization(w, p, input.OrganizationID) {
		return
	}
	value, err := a.service.ReverseSettlement(r.Context(), p.TenantID, input.OrganizationID, r.PathValue("id"), input.ExpectedVersion, input.Reason)
	if err != nil {
		writeRoyaltyResult(w, err)
		return
	}
	writeJSON(w, 201, value)
}
func (a royaltyAPI) reconcile(w http.ResponseWriter, r *http.Request) {
	p, ok := a.auth(w, r, "royalty:reconcile")
	if !ok {
		return
	}
	var input struct {
		OrganizationID    string `json:"organization_id"`
		ExternalReference string `json:"external_reference"`
		ActualMinorUnits  int64  `json:"actual_minor_units"`
	}
	if !decodeStrict(w, r, &input) {
		return
	}
	if !royaltyOrganization(w, p, input.OrganizationID) {
		return
	}
	value, err := a.service.Reconcile(r.Context(), p.TenantID, input.OrganizationID, p.Subject, royalty.Reconciliation{SettlementID: r.PathValue("id"), ExternalReference: input.ExternalReference, ActualMinorUnits: input.ActualMinorUnits})
	if err != nil {
		writeRoyaltyResult(w, err)
		return
	}
	writeJSON(w, 201, value)
}
func writeRoyaltyResult(w http.ResponseWriter, err error) {
	if errors.Is(err, royalty.ErrInvalid) {
		writeProblem(w, 400, "INVALID_ROYALTY_COMMAND", "royalty command does not match contract")
		return
	}
	if errors.Is(err, royalty.ErrConflict) {
		writeProblem(w, 409, "ROYALTY_CONFLICT", "royalty state, scope, idempotency or version conflict")
		return
	}
	writeProblem(w, 500, "INTERNAL", "request failed")
}
````

### FILE: `internal/platform/httpapi/royalty_test.go`

```yaml
block_id: "GO-FRANCHISE-ROYALTY:http-test:v1"
operation: CREATE
provenance: AUTHORED
source: "local authorization and scope regression"
license: "LicenseRef-Workspace-Owner"
sha256: "cbbf1ed637a69b50e6f566d7a1700ce7a5ea38f5017caed93297ee7cfb57764c"
variables: []
secrets_allowed: false
```

````go
package httpapi

import (
	"context"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/royalty"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type royaltyRepo struct {
	policies int
	accruals int
}

func (r *royaltyRepo) CreatePolicy(context.Context, string, string, royalty.Policy) error {
	r.policies++
	return nil
}
func (r *royaltyRepo) AccruePayment(_ context.Context, _, _ string, value royalty.PaymentEvent, _ string) (royalty.Accrual, error) {
	r.accruals++
	return royalty.Accrual{ID: value.ID, OrganizationID: value.OrganizationID, SourceState: value.State}, nil
}
func (*royaltyRepo) OpenSettlement(context.Context, string, string, royalty.Settlement) error {
	return nil
}
func (*royaltyRepo) CloseSettlement(context.Context, string, string, string, string, int64) (royalty.Settlement, error) {
	return royalty.Settlement{Status: "closed"}, nil
}
func (*royaltyRepo) ReverseSettlement(context.Context, string, string, string, string, string, int64, string) (royalty.Settlement, error) {
	return royalty.Settlement{Status: "closed", ReversalOf: "s"}, nil
}
func (*royaltyRepo) Reconcile(_ context.Context, _, _ string, value royalty.Reconciliation, _, _ string) (royalty.Reconciliation, error) {
	value.Status = "matched"
	return value, nil
}

type royaltyIDs struct{ n int }

func (i *royaltyIDs) New() string { i.n++; return "id" + string(rune('0'+i.n)) }

type royaltyVerifier struct{ principal identity.Principal }

func (v royaltyVerifier) Verify(context.Context, string) (identity.Principal, error) {
	return v.principal, nil
}

func TestRoyaltyHTTPAuthorizationScopeAndStrictCommands(t *testing.T) {
	repo := &royaltyRepo{}
	service := royalty.NewService(repo, &royaltyIDs{})
	mux := http.NewServeMux()
	principal := identity.Principal{Subject: "controller", TenantID: "018f4d4a-7b36-7a21-8d10-2f4c54c29d00", Permissions: map[string]struct{}{"royalty:policy": {}, "royalty:post": {}}, Organizations: map[string]struct{}{"franchise": {}}}
	RoyaltyModule{Service: service}.Register(mux, royaltyVerifier{principal: principal})
	now := time.Now().UTC().Truncate(time.Second)
	request := httptest.NewRequest("POST", "/v1/royalties/policies", strings.NewReader(`{"agreement_id":"agreement","organization_id":"franchise","currency":"ARS","rate_basis_points":650,"valid_from":"`+now.Format(time.RFC3339)+`"}`))
	request.Header.Set("Authorization", "Bearer valid")
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 201 || repo.policies != 1 {
		t.Fatalf("policy status=%d count=%d body=%s", response.Code, repo.policies, response.Body.String())
	}
	request = httptest.NewRequest("POST", "/v1/royalties/accruals/from-payment", strings.NewReader(`{"organization_id":"other","payment_attempt_id":"payment","payment_expected_version":3,"state":"captured","occurred_at":"`+now.Format(time.RFC3339)+`"}`))
	request.Header.Set("Authorization", "Bearer valid")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Idempotency-Key", "provider:event")
	response = httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 403 || repo.accruals != 0 {
		t.Fatalf("cross-scope status=%d accruals=%d", response.Code, repo.accruals)
	}
	request = httptest.NewRequest("POST", "/v1/royalties/settlements/s/close", strings.NewReader(`{"organization_id":"franchise","expected_version":1}`))
	request.Header.Set("Authorization", "Bearer valid")
	request.Header.Set("Content-Type", "application/json")
	response = httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 403 {
		t.Fatalf("permission status=%d", response.Code)
	}
}
````

### FILE: `db/migrations/0010_franchise_royalty_settlement.up.sql`

```yaml
block_id: "GO-FRANCHISE-ROYALTY:migration-up:v1"
operation: CREATE
provenance: AUTHORED
source: "local PostgreSQL schema governed by pinned Microsoft BCApps evidence"
license: "LicenseRef-Workspace-Owner"
sha256: "8980ea12f21d75ed8d398559ef4ab2cb8c58e716453ea253156ba5fda39c482b"
variables: []
secrets_allowed: false
```

````sql
create schema if not exists royalty;

create table royalty.policy (
  tenant_id uuid not null,
  policy_id text not null,
  agreement_id text not null,
  franchise_organization_id text not null,
  currency text not null check (currency ~ '^[A-Z]{3}$'),
  rate_basis_points integer not null check (rate_basis_points between 1 and 10000),
  basis text not null default 'captured-payment' check (basis = 'captured-payment'),
  valid_from timestamptz not null,
  valid_until timestamptz,
  created_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, policy_id),
  foreign key (tenant_id, agreement_id) references franchise.agreement (tenant_id, agreement_id),
  foreign key (tenant_id, franchise_organization_id) references org.organization (tenant_id, organization_id),
  check (valid_until is null or valid_until > valid_from)
);

create index royalty_policy_resolution_idx
  on royalty.policy (tenant_id, franchise_organization_id, currency, valid_from desc);

create table royalty.accrual (
  tenant_id uuid not null,
  accrual_id text not null,
  policy_id text not null,
  agreement_id text not null,
  franchise_organization_id text not null,
  payment_attempt_id text not null,
  source_event_key text not null,
  source_state text not null check (source_state in ('captured','refunded')),
  currency text not null check (currency ~ '^[A-Z]{3}$'),
  basis_minor_units bigint not null check (basis_minor_units > 0),
  royalty_minor_units bigint not null check (royalty_minor_units <> 0),
  occurred_at timestamptz not null,
  created_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, accrual_id),
  foreign key (tenant_id, policy_id) references royalty.policy (tenant_id, policy_id),
  foreign key (tenant_id, agreement_id) references franchise.agreement (tenant_id, agreement_id),
  foreign key (tenant_id, franchise_organization_id) references org.organization (tenant_id, organization_id),
  foreign key (tenant_id, payment_attempt_id) references payment.payment_attempt (tenant_id, payment_attempt_id),
  unique (tenant_id, source_event_key),
  unique (tenant_id, payment_attempt_id, source_state),
  check ((source_state = 'captured' and royalty_minor_units > 0) or (source_state = 'refunded' and royalty_minor_units < 0))
);

create index royalty_accrual_unsettled_idx
  on royalty.accrual (tenant_id, franchise_organization_id, currency, occurred_at, accrual_id);

create table royalty.settlement_run (
  tenant_id uuid not null,
  settlement_id text not null,
  franchise_organization_id text not null,
  currency text not null check (currency ~ '^[A-Z]{3}$'),
  period_start timestamptz not null,
  period_end timestamptz not null,
  status text not null check (status in ('draft','closed','reversed')),
  expected_minor_units bigint not null default 0,
  version bigint not null check (version > 0),
  reversal_of text,
  reversal_reason text,
  created_at timestamptz not null default clock_timestamp(),
  closed_at timestamptz,
  primary key (tenant_id, settlement_id),
  foreign key (tenant_id, franchise_organization_id) references org.organization (tenant_id, organization_id),
  foreign key (tenant_id, reversal_of) references royalty.settlement_run (tenant_id, settlement_id),
  check (period_end > period_start),
  check ((status = 'draft' and closed_at is null) or (status in ('closed','reversed') and closed_at is not null)),
  check ((reversal_of is null and reversal_reason is null) or (reversal_of is not null and length(reversal_reason) between 3 and 500))
);

create unique index royalty_settlement_one_reversal_idx
  on royalty.settlement_run (tenant_id, reversal_of)
  where reversal_of is not null;

create unique index royalty_settlement_one_draft_period_idx
  on royalty.settlement_run (tenant_id, franchise_organization_id, currency, period_start, period_end)
  where reversal_of is null and status = 'draft';

create table royalty.settlement_line (
  tenant_id uuid not null,
  settlement_id text not null,
  line_id text not null,
  accrual_id text not null,
  amount_minor_units bigint not null check (amount_minor_units <> 0),
  original_line_id text,
  primary key (tenant_id, settlement_id, line_id),
  foreign key (tenant_id, settlement_id) references royalty.settlement_run (tenant_id, settlement_id),
  foreign key (tenant_id, accrual_id) references royalty.accrual (tenant_id, accrual_id)
);

create unique index royalty_settlement_line_once_idx
  on royalty.settlement_line (tenant_id, accrual_id)
  where original_line_id is null;

create unique index royalty_settlement_reversal_once_idx
  on royalty.settlement_line (tenant_id, original_line_id)
  where original_line_id is not null;

create table royalty.reconciliation (
  tenant_id uuid not null,
  reconciliation_id text not null,
  settlement_id text not null,
  external_reference text not null,
  expected_minor_units bigint not null,
  actual_minor_units bigint not null,
  difference_minor_units bigint not null,
  status text not null check (status in ('matched','mismatch')),
  recorded_by text not null,
  created_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, reconciliation_id),
  foreign key (tenant_id, settlement_id) references royalty.settlement_run (tenant_id, settlement_id),
  unique (tenant_id, settlement_id),
  unique (tenant_id, external_reference),
  check (difference_minor_units = actual_minor_units - expected_minor_units),
  check ((status = 'matched') = (difference_minor_units = 0))
);

create or replace function royalty.prevent_financial_history_mutation()
returns trigger language plpgsql as $$
begin
  raise exception 'immutable royalty history';
end;
$$;

create trigger royalty_policy_immutable before update or delete on royalty.policy
for each row execute function royalty.prevent_financial_history_mutation();
create trigger royalty_accrual_immutable before update or delete on royalty.accrual
for each row execute function royalty.prevent_financial_history_mutation();
create trigger royalty_line_immutable before update or delete on royalty.settlement_line
for each row execute function royalty.prevent_financial_history_mutation();
create trigger royalty_reconciliation_immutable before update or delete on royalty.reconciliation
for each row execute function royalty.prevent_financial_history_mutation();
````

### FILE: `db/migrations/0010_franchise_royalty_settlement.down.sql`

```yaml
block_id: "GO-FRANCHISE-ROYALTY:migration-down:v1"
operation: CREATE
provenance: AUTHORED
source: "local development rollback"
license: "LicenseRef-Workspace-Owner"
sha256: "a3883753722804aae217322fdf58005781af83e3698b2ac8c2354b249d1b971c"
variables: []
secrets_allowed: false
```

````sql
drop trigger if exists royalty_reconciliation_immutable on royalty.reconciliation;
drop trigger if exists royalty_line_immutable on royalty.settlement_line;
drop trigger if exists royalty_accrual_immutable on royalty.accrual;
drop trigger if exists royalty_policy_immutable on royalty.policy;
drop function if exists royalty.prevent_financial_history_mutation();
drop table if exists royalty.reconciliation;
drop table if exists royalty.settlement_line;
drop table if exists royalty.settlement_run;
drop table if exists royalty.accrual;
drop table if exists royalty.policy;
drop schema if exists royalty;
````

### FILE: `db/tests/0010_franchise_royalty_settlement.test.sql`

```yaml
block_id: "GO-FRANCHISE-ROYALTY:migration-test:v1"
operation: CREATE
provenance: AUTHORED
source: "local SQL invariant regression"
license: "LicenseRef-Workspace-Owner"
sha256: "c6b70bb53fab087181000ea9ab0d24e6176f5fda93a01be11d340d279334509d"
variables: []
secrets_allowed: false
```

````sql
begin;

insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c29000','royalty-test','Royalty Test','Royalty Test');
insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c29000','franchise-1','franchise-1','Franchise One','franchisee');
insert into franchise.agreement(tenant_id,agreement_id,franchise_organization_id,territory_code,terms_version,starts_on,status)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c29000','agreement-1','franchise-1','AR-BUE-001','terms-v1','2026-01-01','active');
insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c29000','order-1','franchise-1','customer-1','confirmed','ARS',10001,1);
insert into payment.payment_attempt(tenant_id,payment_attempt_id,order_id,provider_code,provider_reference,idempotency_key,state,currency,amount_minor_units,version)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c29000','payment-1','order-1','provider','provider-1','idem-1','captured','ARS',10001,3);
insert into royalty.policy(tenant_id,policy_id,agreement_id,franchise_organization_id,currency,rate_basis_points,valid_from)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c29000','policy-1','agreement-1','franchise-1','ARS',650,'2026-01-01');
insert into royalty.accrual(tenant_id,accrual_id,policy_id,agreement_id,franchise_organization_id,payment_attempt_id,source_event_key,source_state,currency,basis_minor_units,royalty_minor_units,occurred_at)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c29000','accrual-1','policy-1','agreement-1','franchise-1','payment-1','provider:capture-1','captured','ARS',10001,650,'2026-08-01');
insert into royalty.settlement_run(tenant_id,settlement_id,franchise_organization_id,currency,period_start,period_end,status,expected_minor_units,version,closed_at)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c29000','settlement-1','franchise-1','ARS','2026-08-01','2026-09-01','closed',650,2,clock_timestamp());
insert into royalty.settlement_line(tenant_id,settlement_id,line_id,accrual_id,amount_minor_units)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c29000','settlement-1','line-1','accrual-1',650);
insert into royalty.reconciliation(tenant_id,reconciliation_id,settlement_id,external_reference,expected_minor_units,actual_minor_units,difference_minor_units,status,recorded_by)
values ('018f4d4a-7b36-7a21-8d10-2f4c54c29000','reconciliation-1','settlement-1','bank-statement-1',650,649,-1,'mismatch','controller-1');

do $$
declare total bigint; difference bigint;
begin
  select sum(amount_minor_units) into total from royalty.settlement_line where tenant_id='018f4d4a-7b36-7a21-8d10-2f4c54c29000' and settlement_id='settlement-1';
  if total <> 650 then raise exception 'settlement total mismatch'; end if;
  select difference_minor_units into difference from royalty.reconciliation where tenant_id='018f4d4a-7b36-7a21-8d10-2f4c54c29000' and settlement_id='settlement-1';
  if difference <> -1 then raise exception 'reconciliation difference mismatch'; end if;
  begin
    update royalty.accrual set royalty_minor_units=1 where tenant_id='018f4d4a-7b36-7a21-8d10-2f4c54c29000' and accrual_id='accrual-1';
    raise exception 'immutable accrual accepted mutation';
  exception when raise_exception then
    if sqlerrm <> 'immutable royalty history' then raise; end if;
  end;
end $$;

rollback;
````

## 6. Configuration surface

No new process variable or secret. The project supplies policy rows only through the protected API after contract approval. Provider adapters must map their authenticated immutable event identity into `Idempotency-Key`; never accept an event key produced by a browser.

## 7. Dependency bill

| Dependency | Pin | Use | License | Source |
|---|---|---|---|---|
| Go | `1.26.7` | domain and HTTP | BSD-3-Clause | `go.dev` |
| PostgreSQL | `18.6` | durable subledger, locking and constraints | PostgreSQL | `postgresql.org` |
| pgx | composed pin `5.10.0` | transactional repository | MIT | `github.com/jackc/pgx` |
| Microsoft BCApps | commit `31a860b...4789` | exact architecture evidence only | MIT | `github.com/microsoft/BCApps` |

## 8. Apply order

Materialize after migration 0003, commerce and franchise administration; apply migration 0010 before serving routes. Wire `RoyaltyModule` into the application composition root. Existing systems must map settled provider events, contract policies and finance permissions explicitly. Rollback disables new posting first; dropping immutable financial history is development-only and never an automatic production rollback.

## 9. Verification

Require clean Markdown reconstruction, exact SHA comparison, formatting, all Go tests/vet/build, empty PostgreSQL 18.6 migrations 0001–0010, SQL invariant test, repository integration, one-winner concurrent close, duplicate/cross-scope/stale rejection, equal-and-opposite reversal, explicit reconciliation difference, migration down/up and global library gates. These gates prove the reusable local subledger, not contractual correctness, taxes, statutory accounting, provider payout truth, production recovery/security/load or business acceptance.

## 10. Reconstruction evidence

The exact rebuild, source hashes, failures and gates are recorded in `reconstruction_evidence/FRANCHISE_ROYALTY_SETTLEMENT_2026-08-30_V120.md`.
