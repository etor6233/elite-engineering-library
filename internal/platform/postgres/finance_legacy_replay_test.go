package postgres

import (
	"context"
	"elite.local/enterprise/internal/royalty"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"testing"
	"time"
)

type financeIDs struct{}

func (financeIDs) New() string { return uuid.NewString() }
func financeFixture(t *testing.T) (context.Context, *pgxpool.Pool, string) {
	t.Helper()
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Fatal("owned finance test database required")
	}
	ctx := context.Background()
	p, e := pgxpool.New(ctx, dsn)
	if e != nil {
		t.Fatal(e)
	}
	t.Cleanup(p.Close)
	tenant := uuid.NewString()
	for _, q := range []string{`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1::uuid,'finance-'||$1::text,'Finance fixture','Finance fixture')`, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'org','org','Local centro','franchisee'),($1,'other','other','Otro local','franchisee')`, `insert into franchise.agreement(tenant_id,agreement_id,franchise_organization_id,territory_code,terms_version,starts_on,status)values($1,'agreement','org','AR-FINANCE','v1','2026-01-01','active')`} {
		if _, e = p.Exec(ctx, q, tenant); e != nil {
			t.Fatal(e)
		}
	}
	return ctx, p, tenant
}
func TestFinanceLegacyReplayRequirement(t *testing.T) {
	if os.Getenv("FINANCE_LEGACY_PROBE") != "true" {
		t.Skip("historical red retained; successor tests are authoritative")
	}
	ctx, p, tenant := financeFixture(t)
	service := royalty.NewService(NewRoyalty(p), financeIDs{})
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	policy := royalty.Policy{AgreementID: "agreement", OrganizationID: "org", Currency: "USD", RateBasisPoints: 500, ValidFrom: from}
	a, e := service.CreatePolicy(ctx, tenant, policy)
	if e != nil {
		t.Fatal(e)
	}
	b, e := service.CreatePolicy(ctx, tenant, policy)
	if e != nil || a.ID != b.ID {
		t.Errorf("RED policy same request cannot recover original receipt: error=%v original=%s replay=%s", e, a.ID, b.ID)
	}
	settlement := royalty.Settlement{OrganizationID: "org", Currency: "USD", PeriodStart: from, PeriodEnd: from.AddDate(0, 1, 0)}
	x, e := service.OpenSettlement(ctx, tenant, settlement)
	if e != nil {
		t.Fatal(e)
	}
	y, e := service.OpenSettlement(ctx, tenant, settlement)
	if e != nil || x.ID != y.ID {
		t.Errorf("RED settlement same request cannot recover original receipt: error=%v original=%s replay=%s", e, x.ID, y.ID)
	}
	var policies, settlements int
	p.QueryRow(ctx, `select count(*) from royalty.policy where tenant_id=$1`, tenant).Scan(&policies)
	p.QueryRow(ctx, `select count(*) from royalty.settlement_run where tenant_id=$1`, tenant).Scan(&settlements)
	t.Logf("Original safeguards prevented duplicate rows: policies=%d settlements=%d. Failure is replay/recovery, not a proved double posting.", policies, settlements)
}
