package trainingbridge

import (
	"context"
	"elite.local/enterprise/internal/approval"
	"elite.local/enterprise/internal/platform/postgres"
	"encoding/json"
	"github.com/google/uuid"
	"testing"
)

func TestTrainingGenericBypassAndAuditImmutability(t *testing.T) {
	f := connectedFixture(t)
	ctx := context.Background()
	id := uuid.NewString()
	hash := f.profile.Hash()
	if _, e := f.store.Start(ctx, f.learner, id, "resource-onboarding", hash); e != nil {
		t.Fatal(e)
	}
	if _, e := f.store.Acknowledge(ctx, f.learner, id, "resource-create-view", hash); e != nil {
		t.Fatal(e)
	}
	a, e := f.store.Submit(ctx, f.learner, id, hash, map[string]string{"recovery": "Consulto la referencia antes de repetir."})
	if e != nil {
		t.Fatal(e)
	}
	if _, e = f.store.reviews.Decide(ctx, f.reviewer, f.learner.TenantID, a.RequestID, "store-1", a.PayloadSHA256, true, "Generic bypass", "training:review", nil); e == nil {
		t.Fatal("generic decision bypassed training fact")
	}
	current, e := f.store.Assessment(ctx, f.learner, a.RequestID)
	if e != nil || current.State != approval.StatePending {
		t.Fatal("generic decision leaked", current.State, e)
	}
	raw, _ := json.Marshal(a.Payload)
	_, e = f.store.reviews.Submit(ctx, f.learner, postgres.HumanApprovalSpec{Request: approval.Request{TenantID: f.learner.TenantID, ID: "duplicate-training-identity", Kind: approval.KindTrainingAssessment, SubjectID: f.learner.Subject, Requester: f.learner.Subject, EvidenceSHA: a.PayloadSHA256}, OrganizationID: "store-1", Payload: raw}, "training:learn", nil)
	if e == nil {
		t.Fatal("one attempt acquired a second assessment identity")
	}
	for _, q := range []string{`update audit.event set evidence='{}' where tenant_id=$1 and resource_type='training-attempt'`, `delete from audit.event where tenant_id=$1 and resource_type='training-attempt'`} {
		if _, e = f.pool.Exec(ctx, q, f.learner.TenantID); e == nil {
			t.Fatal("training facts altered")
		}
	}
	if _, e = f.store.Assess(ctx, f.reviewer, a.RequestID, a.PayloadSHA256, true, "Distinct human review recorded"); e != nil {
		t.Fatal(e)
	}
	var requests, decisions, facts int
	e = f.pool.QueryRow(ctx, `select(select count(*) from approval.request where tenant_id=$1),(select count(*) from approval.decision where tenant_id=$1),(select count(*) from audit.event where tenant_id=$1 and resource_type='training-attempt')`, f.learner.TenantID).Scan(&requests, &decisions, &facts)
	if e != nil || requests != 1 || decisions != 1 || facts != 4 {
		t.Fatal("rejected mutation left effects", requests, decisions, facts, e)
	}
	t.Log("TRAINING_PHYSICAL_GUARDS_PASS generic decision without fact rejected; duplicate assessment refused; audit update/delete denied; normal human review preserved")
}
