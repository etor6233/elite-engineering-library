package approval

import (
	"errors"
	"strings"
	"testing"
)

func req(id, kind string, amount int64) Request {
	return Request{
		TenantID: "tenant-a", ID: id, Kind: Kind(kind), SubjectID: "customer-1",
		AmountMinorUnits: amount, Requester: "seller-1",
		EvidenceSHA: strings.Repeat("a", 64),
	}
}

func policy() Policy {
	return Policy{
		AutoApproveMinorUnits: 1000,
		DualControlMinorUnits: 50000,
		MaxOpenPerSubject:     3,
	}
}

func TestSubmitAutoApprovesLowValueNonMoney(t *testing.T) {
	r := NewRegistry(policy())
	st, err := r.Submit(req("r1", "sale", 500))
	if err != nil || st != StateApproved {
		t.Fatalf("expected auto-approve, got %q/%v", st, err)
	}
	if len(r.Audit()) != 1 || r.Audit()[0].Reviewer != "system" {
		t.Fatalf("expected one system decision, got %v", r.Audit())
	}
}

func TestMoneyNeverAutoApproves(t *testing.T) {
	r := NewRegistry(policy())
	for _, k := range []string{"payment", "refund"} {
		st, err := r.Submit(req("r-"+k, k, 1))
		if err != nil || st != StatePending {
			t.Fatalf("money %s should stay pending, got %q/%v", k, st, err)
		}
	}
}

func TestVelocityFlag(t *testing.T) {
	r := NewRegistry(policy())
	for i := 0; i < 3; i++ {
		if _, err := r.Submit(req("v"+string(rune('0'+i)), "sale", 10000)); err != nil {
			t.Fatalf("unexpected error before bound: %v", err)
		}
	}
	if _, err := r.Submit(req("v3", "sale", 10000)); !errors.Is(err, ErrSubjectFlagged) {
		t.Fatalf("expected velocity flag, got %v", err)
	}
}

func TestDuplicateRejected(t *testing.T) {
	r := NewRegistry(policy())
	_, _ = r.Submit(req("d1", "sale", 500))
	if _, err := r.Submit(req("d1", "sale", 500)); !errors.Is(err, ErrDuplicate) {
		t.Fatalf("expected duplicate, got %v", err)
	}
}

func TestApproveSeparationOfDuties(t *testing.T) {
	r := NewRegistry(policy())
	_, _ = r.Submit(req("a1", "sale", 10000))
	if _, err := r.Approve("tenant-a", "a1", "seller-1", "ok"); !errors.Is(err, ErrSeparation) {
		t.Fatalf("expected separation error, got %v", err)
	}
	st, err := r.Approve("tenant-a", "a1", "manager-1", "ok")
	if err != nil || st != StateApproved {
		t.Fatalf("expected approve, got %q/%v", st, err)
	}
}

func TestDualControlRequiresTwoReviewers(t *testing.T) {
	r := NewRegistry(policy())
	_, _ = r.Submit(req("dc1", "sale", 60000)) // >= DualControlMinorUnits
	st, err := r.Approve("tenant-a", "dc1", "manager-1", "first")
	if err != nil || st != StatePending {
		t.Fatalf("expected still pending after one reviewer, got %q/%v", st, err)
	}
	if _, err := r.Approve("tenant-a", "dc1", "manager-1", "dup"); !errors.Is(err, ErrDuplicateApprover) {
		t.Fatalf("expected duplicate approver error, got %v", err)
	}
	st, err = r.Approve("tenant-a", "dc1", "manager-2", "second")
	if err != nil || st != StateApproved {
		t.Fatalf("expected approved after two reviewers, got %q/%v", st, err)
	}
}

func TestReject(t *testing.T) {
	r := NewRegistry(policy())
	_, _ = r.Submit(req("rj1", "sale", 10000))
	st, err := r.Reject("tenant-a", "rj1", "manager-1", "fraudulent")
	if err != nil || st != StateRejected {
		t.Fatalf("expected reject, got %q/%v", st, err)
	}
	if _, err := r.Approve("tenant-a", "rj1", "manager-2", "late"); !errors.Is(err, ErrNotPending) {
		t.Fatalf("expected not-pending after reject, got %v", err)
	}
}
