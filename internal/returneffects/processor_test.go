package returneffects

import (
	"context"
	"errors"
	"testing"
	"time"
)

type fixedIDs struct{ value string }

func (f fixedIDs) New() string { return f.value }

type fakeStore struct {
	work       *Work
	claimErr   error
	applyErr   error
	result     Result
	completion Completion
	finishes   int
}

func (f *fakeStore) Claim(context.Context, string, string, string, time.Duration) (*Work, error) {
	return f.work, f.claimErr
}
func (f *fakeStore) ApplyInventory(context.Context, Work, string) (Result, error) {
	return f.result, f.applyErr
}
func (f *fakeStore) Finish(_ context.Context, _ Work, _ string, c Completion) error {
	f.finishes++
	f.completion = c
	return nil
}
func (f *fakeStore) ResumeBlocked(context.Context, string, string, string, string, string) error {
	return nil
}

func TestInventoryProcessorSuccessNoWorkAndFailClosed(t *testing.T) {
	if _, err := NewInventoryProcessor(nil, fixedIDs{"id"}, "worker", time.Minute, time.Second); err == nil {
		t.Fatal("nil store accepted")
	}
	noWork := &fakeStore{}
	processor, err := NewInventoryProcessor(noWork, fixedIDs{"claim"}, "worker-1", time.Minute, time.Second)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = processor.ProcessOne(context.Background()); !errors.Is(err, ErrNoWork) {
		t.Fatalf("no work err=%v", err)
	}

	work := &Work{RequestID: "request", ClaimToken: "claim"}
	success := &fakeStore{work: work, result: Result{RequestID: "request", InventoryState: "quarantine", StockVersion: 2}}
	processor, _ = NewInventoryProcessor(success, fixedIDs{"claim"}, "worker-1", time.Minute, time.Second)
	result, err := processor.ProcessOne(context.Background())
	if err != nil || result.InventoryState != "quarantine" || success.finishes != 0 {
		t.Fatalf("result=%+v finishes=%d err=%v", result, success.finishes, err)
	}

	blocked := &fakeStore{work: work, applyErr: ErrInventoryState}
	processor, _ = NewInventoryProcessor(blocked, fixedIDs{"claim"}, "worker-1", time.Minute, time.Second)
	if _, err = processor.ProcessOne(context.Background()); !errors.Is(err, ErrInventoryState) {
		t.Fatalf("blocked err=%v", err)
	}
	if blocked.finishes != 1 || blocked.completion.Outcome != "blocked" || blocked.completion.ErrorCode != "INVENTORY_STATE_CONFLICT" {
		t.Fatalf("completion=%+v", blocked.completion)
	}

	transient := &fakeStore{work: work, applyErr: errors.New("database unavailable")}
	processor, _ = NewInventoryProcessor(transient, fixedIDs{"claim"}, "worker-1", time.Minute, 3*time.Second)
	if _, err = processor.ProcessOne(context.Background()); err == nil {
		t.Fatal("transient failure hidden")
	}
	if transient.completion.Outcome != "retry" || transient.completion.RetryAfter != 3*time.Second {
		t.Fatalf("completion=%+v", transient.completion)
	}
}
