package order

import (
	"errors"
	"testing"
)

func TestTransitionRequiresPermissionAndIncrementsVersion(t *testing.T) {
	o, err := New("order-1", "018f4d4a-7b36-7a21-8d10-2f4c54c28a01", "org-1", "customer-1", "USD", 100)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = o.Transition(Placed, map[string]struct{}{}); !errors.Is(err, ErrConflict) {
		t.Fatalf("expected conflict, got %v", err)
	}
	next, err := o.Transition(Placed, map[string]struct{}{"order:create": {}})
	if err != nil {
		t.Fatal(err)
	}
	if next.State != Placed || next.Version != 2 {
		t.Fatalf("unexpected order: %+v", next)
	}
}

func TestTransitionRejectsInvalidStateJump(t *testing.T) {
	o, _ := New("order-1", "018f4d4a-7b36-7a21-8d10-2f4c54c28a01", "org-1", "customer-1", "USD", 100)
	if _, err := o.Transition(Delivered, map[string]struct{}{"*": {}}); !errors.Is(err, ErrConflict) {
		t.Fatalf("expected conflict, got %v", err)
	}
}
