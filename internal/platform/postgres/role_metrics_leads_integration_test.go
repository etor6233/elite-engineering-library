package postgres_test

import (
	"context"
	db "elite.local/enterprise/internal/platform/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"testing"
)

func TestRoleMetricsConvertedLeadBaseline(t *testing.T) {
	raw := os.Getenv("PAYMENT_CONNECTED_DB_URL")
	if raw == "" {
		t.Skip("owned database required")
	}
	ctx := context.Background()
	pool, e := pgxpool.New(ctx, raw)
	if e != nil {
		t.Fatal(e)
	}
	defer pool.Close()
	tenant := uuid.NewString()
	for _, q := range []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1,'metric-lead','Synthetic','Synthetic')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'root','root','Synthetic','franchisor')`,
		`insert into crm.lead(tenant_id,lead_id,organization_id,lifecycle_state,source_code,contact_payload)values($1,'new','root','new','fixture','{}'),($1,'converted','root','converted','fixture','{}'),($1,'lost','root','lost','fixture','{}')`,
	} {
		if _, e = pool.Exec(ctx, q, tenant); e != nil {
			t.Fatal(e)
		}
	}
	v, e := db.NewEnterpriseQuery(pool).Overview(ctx, tenant, "root")
	if e != nil || v.OpenLeads != 1 {
		t.Fatalf("converted lead must not remain open: %+v %v", v, e)
	}
	t.Log("ROLE_METRIC_CONVERTED_LEAD_PASS converted_and_lost_excluded=true")
}
