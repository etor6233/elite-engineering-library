package franchisejourney

import (
	"context"
	"strings"
	"testing"
	"time"
)

type outcomeFake struct {
	fakeRepository
	value ReturnOutcome
	calls int
}

func (r *outcomeFake) ReturnOutcome(_ context.Context, _, _, _ string) (ReturnOutcome, error) {
	r.calls++
	return r.value, nil
}
func outcomeFixture() ReturnOutcome {
	now := time.Unix(1000, 0)
	v := ReturnOutcome{AuthorizationID: "auth", OrganizationID: "store", OrderID: "order", DispositionID: "disp", Remedy: "refund", ObservedAt: now, Stages: []ReturnOutcomeStage{}}
	for _, kind := range []string{"inventory", "refund", "accounting", "fiscal"} {
		v.Stages = append(v.Stages, ReturnOutcomeStage{RequestID: kind, Kind: kind, Status: "requested", UpdatedAt: now})
	}
	return v
}
func TestCommercialCareOutcomeRejectsIncoherentClaims(t *testing.T) {
	for name, change := range map[string]func(*ReturnOutcome){"wrongOrg": func(v *ReturnOutcome) { v.OrganizationID = "other" }, "wrongAuthorization": func(v *ReturnOutcome) { v.AuthorizationID = "other" }, "missingStage": func(v *ReturnOutcome) { v.Stages = v.Stages[:3] }, "duplicate": func(v *ReturnOutcome) { v.Stages[1] = v.Stages[0] }, "successWithoutEvidence": func(v *ReturnOutcome) { v.Stages[0].Status = "succeeded" }, "refundSuccessWithoutOwner": func(v *ReturnOutcome) {
		v.Stages[1].Status = "succeeded"
		v.Stages[1].ResultSHA256 = strings.Repeat("a", 64)
	}, "unknownState": func(v *ReturnOutcome) { v.Stages[0].Status = "complete" }, "retryWithoutReason": func(v *ReturnOutcome) { v.Stages[0].Status = "retry" }, "queuedFiscalClaimedSuccess": func(v *ReturnOutcome) {
		v.Stages[3].Status = "succeeded"
		v.Stages[3].ResultSHA256 = strings.Repeat("a", 64)
		v.Stages[3].Owner = &ReturnOwnerOutcome{State: "queued", Reference: "invoice", AmountMinorUnits: "100", Currency: "ARS"}
	}} {
		t.Run(name, func(t *testing.T) {
			v := outcomeFixture()
			change(&v)
			repo := &outcomeFake{value: v}
			s := NewService(repo, nil, nil)
			if _, e := s.ReturnOutcome(context.Background(), "tenant", "store", "auth"); e == nil {
				t.Fatal("incoherent outcome accepted")
			}
		})
	}
}
func TestCommercialCareOutcomePreservesIndependentStatuses(t *testing.T) {
	v := outcomeFixture()
	v.Stages[0].Status = "succeeded"
	v.Stages[0].ResultSHA256 = strings.Repeat("a", 64)
	v.Stages[1].Status = "retry"
	v.Stages[1].ErrorCode = "PROVIDER_PENDING"
	v.Stages[1].Owner = &ReturnOwnerOutcome{State: "pending", Reference: "refund-fixture", AmountMinorUnits: "9007199254740993", Currency: "ARS"}
	if e := v.Validate("store", "auth"); e != nil {
		t.Fatal(e)
	}
	repo := &outcomeFake{value: v}
	got, e := NewService(repo, nil, nil).ReturnOutcome(context.Background(), "tenant", "store", "auth")
	if e != nil || got.Stages[1].Owner.AmountMinorUnits != "9007199254740993" || got.Stages[1].Status != "retry" {
		t.Fatalf("projection loss %v %v", got, e)
	}
}
