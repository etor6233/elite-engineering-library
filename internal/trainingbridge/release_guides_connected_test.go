package trainingbridge

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"testing"
)

func TestTrainingReleaseGuidesConnected(t *testing.T) {
	f := connectedFixture(t)
	ctx := context.Background()
	hash := f.profile.Hash()
	id := uuid.NewString()
	if _, e := f.store.Start(ctx, f.learner, id, "content-onboarding", hash); e != nil {
		t.Fatal(e)
	}
	v, e := f.store.Read(ctx, f.learner, id)
	if e != nil || v.View.ProfileRevision != 2 || len(v.View.Articles) != 3 {
		t.Fatal(v, e)
	}
	expected := []string{"help-cms-view", "catalog-role-view", "training-role-view"}
	for i, article := range v.View.Articles {
		if article.ID != expected[i] || article.Version != "1.0.0" || len(article.Paragraphs) != 3 {
			t.Fatal("wrong same-release content", article)
		}
	}
	answers := map[string]string{"review": "Consulto el comando pendiente, reviso el contenido exacto y preparo una práctica nueva; la evaluación no concede permisos."}
	if _, e = f.store.Submit(ctx, f.learner, id, hash, answers); e == nil {
		t.Fatal("unread curriculum admitted")
	}
	for _, lesson := range expected {
		if _, e = f.store.Acknowledge(ctx, f.learner, id, lesson, hash); e != nil {
			t.Fatal(e)
		}
	}
	a, e := f.store.Submit(ctx, f.learner, id, hash, answers)
	if e != nil {
		t.Fatal(e)
	}
	if _, e = f.store.Assess(ctx, f.learner, a.RequestID, a.PayloadSHA256, true, "Self review"); e == nil {
		t.Fatal("self approval")
	}
	if _, e = f.store.Assess(ctx, f.reviewer, a.RequestID, a.PayloadSHA256, true, "Tres guías de la misma revisión contrastadas con la respuesta"); e != nil {
		t.Fatal(e)
	}
	// Future activation keeps this immutable, approved attempt on revision2.
	d, b := profileFixture(t)
	d.TenantID = f.learner.TenantID
	d.Revision = 3
	raw, _ := json.Marshal(d)
	next, e := ParseProfile(raw, b, Activation{true, d.ID, d.Revision, digest(raw), d.TenantID, d.OrganizationID})
	if e != nil {
		t.Fatal(e)
	}
	restarted, e := NewStore(f.pool, next)
	if e != nil {
		t.Fatal(e)
	}
	v, e = restarted.Read(ctx, f.learner, id)
	if e != nil || v.CurrentProfile || v.View.ProfileRevision != 2 || v.Assessment == nil || v.Assessment.State != "approved" || v.Assessment.Reviewer != "reviewer" {
		t.Fatal("prior curriculum changed", v, e)
	}
	if _, e = restarted.Acknowledge(ctx, f.learner, id, expected[0], hash); e == nil {
		t.Fatal("stale profile mutated")
	}
	var facts, events, requests, decisions, resources int
	e = f.pool.QueryRow(ctx, `select(select count(*)from audit.event where tenant_id=$1 and resource_type='training-attempt'),(select count(*)from platform.outbox_event where tenant_id=$1 and aggregate_type='training-evidence'),(select count(*)from approval.request where tenant_id=$1 and kind='training_assessment'),(select count(*)from approval.decision where tenant_id=$1),(select count(*)from crm.service_resource where tenant_id=$1)`, f.learner.TenantID).Scan(&facts, &events, &requests, &decisions, &resources)
	if e != nil || facts != 6 || events != 6 || requests != 1 || decisions != 1 || resources != 0 {
		t.Fatal("effects", facts, events, requests, decisions, resources, e)
	}
	if f.learner.Allowed("help:publish") || f.learner.Allowed("network:admin") {
		t.Fatal("training granted authority")
	}
	t.Log("HELP_RELEASE_TRAINING_PASS revision2_three_new_guides=true reading_submission_distinct_human_review=true six_facts_six_events=true future_profile_preserves_history=true permission_grants=0")
}
