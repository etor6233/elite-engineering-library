package marketplacebridge_test

import (
	"elite.local/enterprise/internal/approval"
	"strings"
	"testing"
)

func TestMarketplaceApprovalAlwaysManual(t *testing.T) {
	r := approval.NewRegistry(approval.Policy{AutoApproveMinorUnits: 1000})
	state, e := r.Submit(approval.Request{TenantID: "fixture", ID: "marketplace-approval-0001", Kind: approval.KindMarketplaceMutation, SubjectID: "fixture-sku", Requester: "maker", EvidenceSHA: strings.Repeat("a", 64)})
	if e != nil || state != approval.StatePending {
		t.Fatal("provider mutation automatically approved", state, e)
	}
	if _, e = r.Approve("fixture", "marketplace-approval-0001", "maker", "self"); e == nil {
		t.Fatal("separation of duties")
	}
	if state, e = r.Approve("fixture", "marketplace-approval-0001", "reviewer", "exact reviewed payload"); e != nil || state != approval.StateApproved {
		t.Fatal(state, e)
	}
}
