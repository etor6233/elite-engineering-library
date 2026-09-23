package whatsappbridge

// AUTHORED negative probe of the already-created immutable partial-review receipt.
import (
	"context"
	"elite.local/enterprise/internal/platform/identity"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/url"
	"os"
	"strings"
	"testing"
)

func TestCampaignDecisionReplayBinding(t *testing.T) {
	dsn, tenant := os.Getenv("ELITE_WHATSAPP_CONNECTED_DATABASE_URL"), os.Getenv("ELITE_CAMPAIGN_PARTIAL_TENANT")
	if dsn == "" || tenant == "" {
		t.Skip("exact owned partial-progress fixture required")
	}
	u, e := url.Parse(dsn)
	if e != nil || u.Hostname() != "127.0.0.1" || !strings.HasPrefix(u.Path, "/elite_whatsapp_connected_") {
		t.Fatal("fixture scope")
	}
	ctx := context.Background()
	pool, e := pgxpool.New(ctx, dsn)
	if e != nil {
		t.Fatal(e)
	}
	defer pool.Close()
	var key, hash, reviewer, reason string
	e = pool.QueryRow(ctx, `select r.request_id,r.evidence_sha,d.reviewer,d.reason from approval.request r join approval.decision d using(tenant_id,request_id)
 where r.tenant_id=$1 and r.kind='whatsapp_schedule' and r.payload->'request'->>'campaign_id'='campaign-partial-0001' and r.payload->'request'->>'step'='1' and r.state='approved'`, tenant).Scan(&key, &hash, &reviewer, &reason)
	if e != nil {
		t.Fatal(e)
	}
	c := &Campaigns{schedule: &ScheduledNotifications{approvals: &ScheduleApprovals{tenant: tenant, org: "store-1", base: &PostgresAppointmentApprovals{pool: pool}}}}
	p := identity.Principal{Subject: reviewer, TenantID: tenant}
	if state, e := c.decisionReplay(ctx, p, key, hash, true, reason); e != nil || state != "approved" {
		t.Fatal("exact replay", state, e)
	}
	for _, v := range []struct {
		actor, hash, reason string
		yes                 bool
	}{
		{"another-reviewer", hash, reason, true}, {reviewer, strings.Repeat("f", 64), reason, true}, {reviewer, hash, reason + " changed", true}, {reviewer, hash, reason, false},
	} {
		q := p
		q.Subject = v.actor
		if _, e := c.decisionReplay(ctx, q, key, v.hash, v.yes, v.reason); e == nil {
			t.Fatal("unbound decision replay accepted")
		}
	}
	t.Log("CAMPAIGN_DECISION_REPLAY_BOUNDARY_PASS exact=1 actor_hash_reason_decision_negatives=4 no_writes")
}
