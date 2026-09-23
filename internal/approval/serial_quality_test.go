package approval

import (
	"strings"
	"testing"
)

func TestSerialQualityAlwaysNeedsDistinctHuman(t *testing.T) {
	r := NewRegistry(Policy{AutoApproveMinorUnits: 1000000})
	request := Request{TenantID: "fixture", ID: "serial-quality", Kind: KindSerialQuality, SubjectID: "unit", Requester: "receiver", EvidenceSHA: strings.Repeat("a", 64)}
	state, err := r.Submit(request)
	if err != nil || state != StatePending {
		t.Fatal("serial quality auto-approved", state, err)
	}
	if _, err = r.Approve("fixture", "serial-quality", "receiver", "self"); err != ErrSeparation {
		t.Fatal("self approval", err)
	}
	state, err = r.Approve("fixture", "serial-quality", "reviewer", "observed evidence")
	if err != nil || state != StateApproved {
		t.Fatal(state, err)
	}
}
