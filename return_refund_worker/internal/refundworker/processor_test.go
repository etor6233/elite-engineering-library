package refundworker

import (
	"context"
	"errors"
	"testing"
	"time"
)

type fakeStore struct {
	work       *Work
	refund     Refund
	prepareErr error
	completed  bool
	finished   bool
	outcome    string
	code       string
	result     ProviderResult
}

func (f *fakeStore) Claim(context.Context, string, string, time.Duration) (*Work, error) {
	return f.work, nil
}
func (f *fakeStore) Prepare(context.Context, Work, string) (Refund, error) {
	return f.refund, f.prepareErr
}
func (f *fakeStore) Complete(_ context.Context, _ Work, _ string, _ Refund, result ProviderResult, outcome, code string, _ time.Duration) error {
	f.completed, f.outcome, f.code, f.result = true, outcome, code, result
	return nil
}
func (f *fakeStore) Finish(_ context.Context, _ Work, _ string, outcome, code string, _ time.Duration) error {
	f.finished, f.outcome, f.code = true, outcome, code
	return nil
}

type fakeProvider struct {
	created, retrieved bool
	result             ProviderResult
	err                error
}

func (f *fakeProvider) Create(context.Context, Refund) (ProviderResult, error) {
	f.created = true
	return f.result, f.err
}
func (f *fakeProvider) Retrieve(context.Context, Refund) (ProviderResult, error) {
	f.retrieved = true
	return f.result, f.err
}

func baseRefund() Refund {
	return Refund{TenantID: "tenant", RequestID: "request", PaymentAttemptID: "payment", Provider: "stripe", ProviderPaymentReference: "pi_1", IdempotencyKey: "return-effect-key-0001", Currency: "ARS", AmountMinorUnits: 1000}
}

func TestProcessorCreatesAndCompletesSucceededRefund(t *testing.T) {
	work := &Work{TenantID: "tenant", RequestID: "request", Attempt: 1, ClaimToken: "claim"}
	store := &fakeStore{work: work, refund: baseRefund()}
	provider := &fakeProvider{result: ProviderResult{ProviderPaymentReference: "pi_1", ProviderRefundReference: "re_1", ProviderStatus: "succeeded", Currency: "ARS", AmountMinorUnits: 1000}}
	processor, err := NewProcessor(store, map[string]Provider{"stripe": provider}, "worker", time.Minute, time.Second, 5)
	if err != nil {
		t.Fatal(err)
	}
	if err = processor.Step(context.Background(), "claim"); err != nil {
		t.Fatal(err)
	}
	if !provider.created || provider.retrieved || !store.completed || store.outcome != "succeeded" || store.code != "" {
		t.Fatalf("provider=%+v store=%+v", provider, store)
	}
}

func TestProcessorRetrievesPendingRefundAndSchedulesRetry(t *testing.T) {
	work := &Work{TenantID: "tenant", RequestID: "request", Attempt: 2, ClaimToken: "claim"}
	refund := baseRefund()
	refund.ProviderRefundReference = "re_1"
	store := &fakeStore{work: work, refund: refund}
	provider := &fakeProvider{result: ProviderResult{ProviderPaymentReference: "pi_1", ProviderRefundReference: "re_1", ProviderStatus: "pending", Currency: "ARS", AmountMinorUnits: 1000}}
	processor, _ := NewProcessor(store, map[string]Provider{"stripe": provider}, "worker", time.Minute, time.Second, 5)
	if err := processor.Step(context.Background(), "claim"); err != nil {
		t.Fatal(err)
	}
	if provider.created || !provider.retrieved || store.outcome != "retry" || store.code != "PROVIDER_PENDING" {
		t.Fatalf("provider=%+v store=%+v", provider, store)
	}
}

func TestProcessorBlocksMappingProviderAndExhaustionFailures(t *testing.T) {
	work := &Work{TenantID: "tenant", RequestID: "request", Attempt: 5, ClaimToken: "claim"}
	store := &fakeStore{work: work, refund: baseRefund(), prepareErr: ErrMappingConflict}
	processor, _ := NewProcessor(store, nil, "worker", time.Minute, time.Second, 5)
	if err := processor.Step(context.Background(), "claim"); !errors.Is(err, ErrMappingConflict) || !store.finished || store.outcome != "blocked" {
		t.Fatalf("mapping err=%v store=%+v", err, store)
	}
	store = &fakeStore{work: work, refund: baseRefund()}
	processor, _ = NewProcessor(store, nil, "worker", time.Minute, time.Second, 5)
	if err := processor.Step(context.Background(), "claim"); err != nil || store.outcome != "blocked" || store.code != "PROVIDER_CONFIG_MISSING" {
		t.Fatalf("missing provider err=%v store=%+v", err, store)
	}
	store = &fakeStore{work: work, refund: baseRefund()}
	provider := &fakeProvider{err: errors.New("transport")}
	processor, _ = NewProcessor(store, map[string]Provider{"stripe": provider}, "worker", time.Minute, time.Second, 5)
	if err := processor.Step(context.Background(), "claim"); err == nil || store.outcome != "blocked" {
		t.Fatalf("exhausted err=%v store=%+v", err, store)
	}
}

func TestProcessorBlocksDivergentProviderResult(t *testing.T) {
	work := &Work{TenantID: "tenant", RequestID: "request", Attempt: 1, ClaimToken: "claim"}
	store := &fakeStore{work: work, refund: baseRefund()}
	provider := &fakeProvider{result: ProviderResult{ProviderPaymentReference: "pi_1", ProviderRefundReference: "re_1", ProviderStatus: "succeeded", Currency: "ARS", AmountMinorUnits: 999}}
	processor, _ := NewProcessor(store, map[string]Provider{"stripe": provider}, "worker", time.Minute, time.Second, 5)
	if err := processor.Step(context.Background(), "claim"); !errors.Is(err, ErrResponseMismatch) || !store.completed || store.outcome != "blocked" {
		t.Fatalf("mismatch err=%v store=%+v", err, store)
	}
}
