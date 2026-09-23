package returnexchange

import (
	"context"
	"errors"
	"testing"
	"time"

	"elite.local/enterprise/internal/returneffects"
)

type fakeIDs struct{ value string }

func (f fakeIDs) New() string { return f.value }

type fakeStore struct {
	work       *returneffects.Work
	result     Result
	executeErr error
	completion returneffects.Completion
}

func (f *fakeStore) Claim(context.Context, string, string, string, time.Duration) (*returneffects.Work, error) {
	return f.work, nil
}
func (f *fakeStore) PrepareExchange(context.Context, returneffects.Work, string) (Result, error) {
	return f.result, f.executeErr
}
func (f *fakeStore) Finish(_ context.Context, _ returneffects.Work, _ string, completion returneffects.Completion) error {
	f.completion = completion
	return nil
}

func TestProcessorSuccessAndFailureClassification(t *testing.T) {
	work := &returneffects.Work{RequestID: "exchange-1"}
	store := &fakeStore{work: work, result: Result{RequestID: "exchange-1", ReplacementOrderID: "order-2"}}
	processor, err := NewProcessor(store, fakeIDs{"claim-1"}, "exchange-worker", time.Minute, 5*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	result, err := processor.ProcessOne(context.Background())
	if err != nil || result.ReplacementOrderID != "order-2" {
		t.Fatalf("result=%+v err=%v", result, err)
	}

	store.executeErr = ErrStockUnavailable
	if _, err = processor.ProcessOne(context.Background()); !errors.Is(err, ErrStockUnavailable) {
		t.Fatalf("stock error=%v", err)
	}
	if store.completion.Outcome != "retry" || store.completion.ErrorCode != "REPLACEMENT_STOCK_UNAVAILABLE" {
		t.Fatalf("completion=%+v", store.completion)
	}

	store.executeErr = ErrConflict
	if _, err = processor.ProcessOne(context.Background()); !errors.Is(err, ErrConflict) {
		t.Fatalf("conflict error=%v", err)
	}
	if store.completion.Outcome != "blocked" || store.completion.ErrorCode != "EXCHANGE_CONTRACT_CONFLICT" {
		t.Fatalf("completion=%+v", store.completion)
	}
}

func TestProcessorRejectsInvalidConfigurationAndNoWork(t *testing.T) {
	if _, err := NewProcessor(&fakeStore{}, fakeIDs{"claim"}, "bad worker", time.Minute, time.Second); err == nil {
		t.Fatal("invalid worker accepted")
	}
	processor, err := NewProcessor(&fakeStore{}, fakeIDs{"claim"}, "worker", time.Minute, time.Second)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = processor.ProcessOne(context.Background()); !errors.Is(err, ErrNoWork) {
		t.Fatalf("no work=%v", err)
	}
}
