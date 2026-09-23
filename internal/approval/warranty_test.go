package approval

import (
	"strings"
	"testing"
)

func TestWarrantyRepairAlwaysNeedsDistinctHuman(t *testing.T) {
	r := NewRegistry(Policy{AutoApproveMinorUnits: 1000000})
	request := Request{TenantID: "fixture", ID: "warranty", Kind: KindWarrantyRepair, SubjectID: "case", Requester: "technician", EvidenceSHA: strings.Repeat("a", 64)}
	state, err := r.Submit(request)
	if err != nil || state != StatePending {
		t.Fatal("warranty auto-approved", state, err)
	}
	if _, err = r.Approve("fixture", "warranty", "technician", "self"); err != ErrSeparation {
		t.Fatal("self approval", err)
	}
	state, err = r.Approve("fixture", "warranty", "reviewer", "observed evidence")
	if err != nil || state != StateApproved {
		t.Fatal(state, err)
	}
}
