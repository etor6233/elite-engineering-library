package approval

import (
	"strings"
	"testing"
)

func TestDocumentReviewNeverAutoApproves(t *testing.T) {
	r := NewRegistry(Policy{AutoApproveMinorUnits: 10000})
	request := Request{TenantID: "fixture-tenant", ID: "document:fixture", Kind: KindDocumentReview, SubjectID: "fixture-document", Requester: "maker", EvidenceSHA: strings.Repeat("a", 64)}
	state, e := r.Submit(request)
	if e != nil || state != StatePending {
		t.Fatal("document automatic approval", state, e)
	}
	if _, e = r.Approve(request.TenantID, request.ID, "maker", "self"); e != ErrSeparation {
		t.Fatal("self approval", e)
	}
	state, e = r.Approve(request.TenantID, request.ID, "reviewer", "reviewed fields")
	if e != nil || state != StateApproved {
		t.Fatal(state, e)
	}
}
