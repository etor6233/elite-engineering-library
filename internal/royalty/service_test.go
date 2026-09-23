package royalty

import (
	"context"
	"errors"
	"testing"
	"time"
)

type fakeRepo struct {
	policy     Policy
	event      PaymentEvent
	settlement Settlement
}

func (f *fakeRepo) CreatePolicy(_ context.Context, _, _ string, value Policy) error {
	f.policy = value
	return nil
}
func (f *fakeRepo) AccruePayment(_ context.Context, _, _ string, value PaymentEvent, _ string) (Accrual, error) {
	f.event = value
	return Accrual{ID: value.ID, SourceState: value.State}, nil
}
func (f *fakeRepo) OpenSettlement(_ context.Context, _, _ string, value Settlement) error {
	f.settlement = value
	return nil
}
func (f *fakeRepo) CloseSettlement(context.Context, string, string, string, string, int64) (Settlement, error) {
	return Settlement{Status: "closed"}, nil
}
func (f *fakeRepo) ReverseSettlement(context.Context, string, string, string, string, string, int64, string) (Settlement, error) {
	return Settlement{Status: "closed", ReversalOf: "s"}, nil
}
func (f *fakeRepo) Reconcile(_ context.Context, _, _ string, value Reconciliation, _, _ string) (Reconciliation, error) {
	value.Status = "matched"
	return value, nil
}

type seqIDs struct{ n int }

func (s *seqIDs) New() string { s.n++; return "id-" + string(rune('0'+s.n)) }

func TestServiceValidatesAndDelegates(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	repo := &fakeRepo{}
	service := NewService(repo, &seqIDs{})
	policy, err := service.CreatePolicy(context.Background(), "tenant", Policy{AgreementID: "agreement", OrganizationID: "org", Currency: "ARS", RateBasisPoints: 650, ValidFrom: now})
	if err != nil || policy.ID == "" || repo.policy.RateBasisPoints != 650 {
		t.Fatalf("policy=%+v err=%v", policy, err)
	}
	accrual, err := service.AccruePayment(context.Background(), "tenant", "provider:event-1", PaymentEvent{OrganizationID: "org", PaymentAttemptID: "pay", PaymentExpectedVersion: 3, State: "captured", OccurredAt: now})
	if err != nil || accrual.SourceState != "captured" || repo.event.ID == "" {
		t.Fatalf("accrual=%+v err=%v", accrual, err)
	}
	settlement, err := service.OpenSettlement(context.Background(), "tenant", Settlement{OrganizationID: "org", Currency: "ARS", PeriodStart: now, PeriodEnd: now.Add(24 * time.Hour)})
	if err != nil || settlement.Status != "draft" || settlement.Version != 1 {
		t.Fatalf("settlement=%+v err=%v", settlement, err)
	}
}

func TestServiceRejectsUntrustedFinancialShapes(t *testing.T) {
	now := time.Now().UTC()
	service := NewService(&fakeRepo{}, &seqIDs{})
	cases := []error{}
	_, err := service.CreatePolicy(context.Background(), "tenant", Policy{AgreementID: "a", OrganizationID: "o", Currency: "ars", RateBasisPoints: 650, ValidFrom: now})
	cases = append(cases, err)
	_, err = service.AccruePayment(context.Background(), "tenant", "", PaymentEvent{OrganizationID: "o", PaymentAttemptID: "p", PaymentExpectedVersion: 1, State: "captured", OccurredAt: now})
	cases = append(cases, err)
	_, err = service.AccruePayment(context.Background(), "tenant", "e", PaymentEvent{OrganizationID: "o", PaymentAttemptID: "p", PaymentExpectedVersion: 1, State: "pending", OccurredAt: now})
	cases = append(cases, err)
	_, err = service.OpenSettlement(context.Background(), "tenant", Settlement{OrganizationID: "o", Currency: "ARS", PeriodStart: now, PeriodEnd: now})
	cases = append(cases, err)
	_, err = service.ReverseSettlement(context.Background(), "tenant", "o", "s", 1, "")
	cases = append(cases, err)
	for i, got := range cases {
		if !errors.Is(got, ErrInvalid) {
			t.Fatalf("case %d err=%v", i, got)
		}
	}
}
