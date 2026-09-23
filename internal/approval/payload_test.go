package approval

import (
	"strings"
	"testing"
)

func TestBoundPayloadAndManualKinds(t *testing.T) {
	a, h, e := CanonicalPayload([]byte(`{"b":2,"a":"reviewed"}`))
	if e != nil || string(a) != `{"a":"reviewed","b":2}` || len(h) != 64 {
		t.Fatal(string(a), h, e)
	}
	for _, raw := range []string{`{"a":1,"a":2}`, `{"a":1,"A":2}`, `{} {}`, `null`, `[]`, `{"a":{"x":1,"X":2}}`} {
		if _, _, e = CanonicalPayload([]byte(raw)); e == nil {
			t.Fatal("invalid payload", raw)
		}
	}
	for _, kind := range []Kind{KindWhatsAppReply, KindSocialPublish, KindSocialRevoke} {
		r := NewRegistry(Policy{AutoApproveMinorUnits: 100})
		state, e := r.Submit(Request{TenantID: "tenant", ID: "request", Kind: kind, SubjectID: "page", Requester: "requester", EvidenceSHA: strings.Repeat("a", 64)})
		if e != nil || state != StatePending {
			t.Fatal("manual kind auto-approved", kind, state, e)
		}
	}
}
