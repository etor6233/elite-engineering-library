package qr

import (
	"context"
	"errors"
	"strings"
	"testing"
)

// Test resolver is bound to an actor and operation, not merely a database row.
type actorResolver struct {
	actor, operation string
	calls            int
	failure          error
}

func (r *actorResolver) Resolve(_ context.Context, tenant string, kind Kind, id string) (bool, error) {
	r.calls++
	if r.failure != nil {
		return false, r.failure
	}
	exists := tenant == "tenant-a" && kind == KindOrder && (id == "ORD-public" || id == "ORD-restricted")
	authorized := r.actor == "operator-a" && r.operation == "view" && id == "ORD-public"
	return exists && authorized, nil
}
func reference(t *testing.T, id string) string {
	t.Helper()
	p, err := NewPayload("tenant-a", KindOrder, id)
	if err != nil {
		t.Fatal(err)
	}
	s, err := p.Encode()
	if err != nil {
		t.Fatal(err)
	}
	return s
}
func TestObjectAuthorizationIsDistinctFromExistence(t *testing.T) {
	for _, tc := range []struct {
		actor, operation, id string
		allow                bool
	}{
		{"operator-a", "view", "ORD-public", true},
		{"operator-a", "view", "ORD-restricted", false},
		{"operator-b", "view", "ORD-public", false},
		{"operator-a", "release", "ORD-public", false},
		{"operator-a", "view", "ORD-missing", false},
	} {
		t.Run(tc.actor+"/"+tc.operation+"/"+tc.id, func(t *testing.T) {
			r := &actorResolver{actor: tc.actor, operation: tc.operation}
			got, err := Verify(context.Background(), reference(t, tc.id), "tenant-a", r)
			if tc.allow {
				if err != nil || got.ID != tc.id {
					t.Fatalf("authorized reference rejected: %v", err)
				}
			} else if !errors.Is(err, ErrUnknown) || got != (Payload{}) {
				t.Fatalf("unauthorized object accepted/leaked: %+v %v", got, err)
			}
		})
	}
}
func TestInvalidTrustedScopeDoesNotCallResolver(t *testing.T) {
	for _, scope := range []string{"", " ", " tenant-a", "tenant-a ", strings.Repeat("a", 65)} {
		r := &actorResolver{}
		got, err := Verify(context.Background(), reference(t, "ORD-public"), scope, r)
		if !errors.Is(err, ErrInvalidScope) || r.calls != 0 || got != (Payload{}) {
			t.Fatalf("scope %q accepted: %v", scope, err)
		}
	}
}
func TestResolverFailureAndNilFailClosed(t *testing.T) {
	sentinel := errors.New("authorization store unavailable")
	r := &actorResolver{failure: sentinel}
	if got, err := Verify(context.Background(), reference(t, "ORD-public"), "tenant-a", r); !errors.Is(err, sentinel) || got != (Payload{}) {
		t.Fatalf("failure hidden: %v", err)
	}
	if _, err := Verify(context.Background(), reference(t, "ORD-public"), "tenant-a", nil); err == nil {
		t.Fatal("nil resolver accepted")
	}
}
func TestCancelledContextDoesNotResolve(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	r := &actorResolver{}
	if _, err := Verify(ctx, reference(t, "ORD-public"), "tenant-a", r); !errors.Is(err, context.Canceled) || r.calls != 0 {
		t.Fatalf("cancel ignored: %v", err)
	}
}
