package trainingbridge

import (
	"context"
	"elite.local/enterprise/internal/platform/identity"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/url"
	"os"
	"strings"
	"sync"
	"testing"
)

type trainingFixture struct {
	store                      *Store
	pool                       *pgxpool.Pool
	learner, reviewer, foreign identity.Principal
	profile                    *Profile
}

func connectedFixture(t *testing.T) *trainingFixture {
	t.Helper()
	raw := os.Getenv("ELITE_TRAINING_DATABASE_URL")
	if raw == "" {
		t.Skip("explicit owned training database required")
	}
	u, e := url.Parse(raw)
	if e != nil || u.Scheme != "postgres" || u.Hostname() != "127.0.0.1" || !strings.HasPrefix(u.Path, "/elite_training_connected_") {
		t.Fatal("dedicated loopback training database required")
	}
	ctx := context.Background()
	pool, e := pgxpool.New(ctx, raw)
	if e != nil {
		t.Fatal(e)
	}
	t.Cleanup(pool.Close)
	d, b := profileFixture(t)
	d.TenantID = uuid.NewString()
	bytes, _ := json.Marshal(d)
	p, e := ParseProfile(bytes, b, Activation{true, d.ID, d.Revision, digest(bytes), d.TenantID, d.OrganizationID})
	if e != nil {
		t.Fatal(e)
	}
	for _, q := range []string{`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1::uuid,'training-'||$1::text,'Fixture','Fixture')`, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'store-1','store-1','Fixture','store'),($1,'other','other','Other','store')`} {
		if _, e = pool.Exec(ctx, q, d.TenantID); e != nil {
			t.Fatal(e)
		}
	}
	s, e := NewStore(pool, p)
	if e != nil {
		t.Fatal(e)
	}
	learner := identity.Principal{TenantID: d.TenantID, Subject: "learner", Organizations: map[string]struct{}{"store-1": {}}, Permissions: map[string]struct{}{"training:learn": {}, "training:review": {}}}
	reviewer := identity.Principal{TenantID: d.TenantID, Subject: "reviewer", Organizations: map[string]struct{}{"store-1": {}}, Permissions: map[string]struct{}{"training:review": {}}}
	foreign := identity.Principal{TenantID: d.TenantID, Subject: "foreign", Organizations: map[string]struct{}{"other": {}}, Permissions: map[string]struct{}{"training:learn": {}, "training:review": {}}}
	return &trainingFixture{s, pool, learner, reviewer, foreign, p}
}
func TestTrainingDurableBindingsAndAtomicity(t *testing.T) {
	f := connectedFixture(t)
	ctx := context.Background()
	id := uuid.NewString()
	hash := f.profile.Hash()
	start := func() error { _, e := f.store.Start(ctx, f.learner, id, "resource-onboarding", hash); return e }
	var wg sync.WaitGroup
	errs := make(chan error, 4)
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); errs <- start() }()
	}
	wg.Wait()
	close(errs)
	for e := range errs {
		if e != nil {
			t.Fatal(e)
		}
	}
	if _, e := f.store.Start(ctx, f.learner, id, "admin-onboarding", hash); e == nil {
		t.Fatal("attempt identity changed course")
	}
	if _, e := f.store.Read(ctx, f.foreign, id); e == nil {
		t.Fatal("cross organization attempt visible")
	}
	other := f.learner
	other.Subject = "other-learner"
	other.Permissions = map[string]struct{}{"training:learn": {}}
	if _, e := f.store.Read(ctx, other, id); e == nil {
		t.Fatal("another learner read attempt")
	}
	answers := map[string]string{"recovery": "Consulto la referencia guardada y no creo otro recurso para forzar un reintento."}
	if _, e := f.store.Submit(ctx, f.learner, id, hash, answers); e == nil {
		t.Fatal("response submitted before declared reading")
	}
	if _, e := f.pool.Exec(ctx, `create function audit.reject_training_fixture() returns trigger language plpgsql as $$ begin if new.event_type='training.lesson-read' then raise exception 'fixture rollback';end if;return new;end $$;create trigger reject_training_fixture before insert on platform.outbox_event for each row execute function audit.reject_training_fixture()`); e != nil {
		t.Fatal(e)
	}
	if _, e := f.store.Acknowledge(ctx, f.learner, id, "resource-create-view", hash); e == nil {
		t.Fatal("outbox failure accepted reading")
	}
	v, e := f.store.Read(ctx, f.learner, id)
	if e != nil || len(v.ReadLessons) != 0 {
		t.Fatal("reading fact survived rollback", e)
	}
	if _, e = f.pool.Exec(ctx, `drop trigger reject_training_fixture on platform.outbox_event;drop function audit.reject_training_fixture()`); e != nil {
		t.Fatal(e)
	}
	for i := 0; i < 2; i++ {
		if _, e = f.store.Acknowledge(ctx, f.learner, id, "resource-create-view", hash); e != nil {
			t.Fatal(e)
		}
	}
	a, e := f.store.Submit(ctx, f.learner, id, hash, answers)
	if e != nil {
		t.Fatal(e)
	}
	if _, e = f.store.Submit(ctx, f.learner, id, hash, map[string]string{"recovery": "Changed immutable response"}); e == nil {
		t.Fatal("submission overwritten")
	}
	if _, e = f.store.Assess(ctx, f.learner, a.RequestID, a.PayloadSHA256, true, "Self review"); e == nil {
		t.Fatal("self review accepted")
	}
	if _, e = f.store.Assess(ctx, f.foreign, a.RequestID, a.PayloadSHA256, true, "Other organization"); e == nil {
		t.Fatal("foreign review accepted")
	}
	if _, e = f.store.Assess(ctx, f.reviewer, a.RequestID, strings.Repeat("f", 64), true, "Wrong payload"); e == nil {
		t.Fatal("different payload assessed")
	}
	if _, e = f.store.Assess(ctx, f.reviewer, a.RequestID, a.PayloadSHA256, true, "Respuesta revisada contra el contenido de esta versión"); e != nil {
		t.Fatal(e)
	}
	restarted, e := NewStore(f.pool, f.profile)
	if e != nil {
		t.Fatal(e)
	}
	v, e = restarted.Read(ctx, f.learner, id)
	if e != nil || v.Assessment == nil || v.Assessment.State != "approved" || v.Assessment.Reviewer != "reviewer" {
		t.Fatal("assessment not durable", e)
	}
	if _, e = f.store.Assess(ctx, f.reviewer, a.RequestID, a.PayloadSHA256, false, "Changed decision"); e == nil {
		t.Fatal("decision changed")
	}
	var facts, events, requests, decisions, resources int
	e = f.pool.QueryRow(ctx, `select (select count(*) from audit.event where tenant_id=$1 and resource_type='training-attempt'),(select count(*) from platform.outbox_event where tenant_id=$1 and aggregate_type='training-evidence'),(select count(*) from approval.request where tenant_id=$1 and kind='training_assessment'),(select count(*) from approval.decision where tenant_id=$1),(select count(*) from crm.service_resource where tenant_id=$1)`, f.learner.TenantID).Scan(&facts, &events, &requests, &decisions, &resources)
	if e != nil || facts != 4 || events != 4 || requests != 1 || decisions != 1 || resources != 0 {
		t.Fatal("unexpected durable effects", facts, events, requests, decisions, resources, e)
	}
	if !f.learner.Allowed("training:learn") || f.learner.Allowed("resource:manage") {
		t.Fatal("training changed permissions")
	}
	if _, e = f.pool.Exec(ctx, `update approval.request set payload=jsonb_set(payload,'{answers,recovery}','"overwrite"') where tenant_id=$1 and request_id=$2`, f.learner.TenantID, a.RequestID); e == nil {
		t.Fatal("training payload update accepted")
	}
	t.Log("TRAINING_BINDINGS_PASS concurrent_start_one_fact=true reading_outbox_rollback=true own_subject_org_hash_bound=true human_review_durable=true duplicate_submission_immutable=true permission_grants=0 resource_creation=0")
}
