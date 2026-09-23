package returnfiscal

import (
	"context"
	"errors"
	"testing"
	"time"

	"elite.local/enterprise/internal/returneffects"
)

type ids struct{}

func (ids) New() string { return "claim" }

type fakeStore struct {
	work       *returneffects.Work
	result     Result
	err        error
	completion returneffects.Completion
}

func (f *fakeStore) Claim(context.Context, string, string, string, time.Duration) (*returneffects.Work, error) {
	return f.work, nil
}
func (f *fakeStore) RequestOrObserve(context.Context, returneffects.Work, string) (Result, error) {
	return f.result, f.err
}
func (f *fakeStore) Finish(_ context.Context, _ returneffects.Work, _ string, c returneffects.Completion) error {
	f.completion = c
	return nil
}

func TestProcessorWaitsForAuthorizationAndClosesOnlyAuthorized(t *testing.T) {
	work := &returneffects.Work{RequestID: "fiscal", EffectKind: "fiscal", OwnerContext: "fiscal"}
	pending := &fakeStore{work: work, result: Result{CreditInvoiceID: "credit", Status: "queued"}, err: ErrPending}
	p, _ := NewProcessor(pending, ids{}, "worker", time.Minute, time.Second)
	if _, err := p.ProcessOne(context.Background()); !errors.Is(err, ErrPending) || pending.completion.Outcome != "retry" || pending.completion.ErrorCode != "FISCAL_AUTHORIZATION_PENDING" {
		t.Fatalf("pending=%+v err=%v", pending.completion, err)
	}
	success := &fakeStore{work: work, result: Result{CreditInvoiceID: "credit", Status: "authorized", ResultSHA256: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}}
	p, _ = NewProcessor(success, ids{}, "worker", time.Minute, time.Second)
	if result, err := p.ProcessOne(context.Background()); err != nil || result.Status != "authorized" || success.completion.Outcome != "succeeded" {
		t.Fatalf("result=%+v completion=%+v err=%v", result, success.completion, err)
	}
}
