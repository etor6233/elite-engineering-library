package main

import (
	"context"
	"elite.local/return-refund-worker/internal/operationaltelemetry"
	"elite.local/return-refund-worker/internal/refundworker"
	"errors"
	"testing"
	"time"
)

func TestObservedLoopStopsOnUnknownDeliveryWithoutNextClaim(t *testing.T) {
	for _, tc := range []struct {
		name    string
		stepErr error
		want    string
	}{{"idle", refundworker.ErrNoWork, "idle"}, {"handled", nil, "handled"}, {"error", errors.New("private-customer-diagnostic"), "error"}} {
		t.Run(tc.name, func(t *testing.T) {
			calls, tokens, reports := 0, 0, 0
			err := runObservedRefundLoop(context.Background(), func(context.Context, string) error { calls++; return tc.stepErr }, func() (string, error) { tokens++; return "synthetic-token", nil }, func(ctx context.Context, outcome string, elapsed time.Duration) error {
				reports++
				if outcome != tc.want || elapsed < 0 {
					t.Fatal("wrong closed outcome")
				}
				return operationaltelemetry.ErrUnavailable
			})
			if !errors.Is(err, operationaltelemetry.ErrUnavailable) || calls != 1 || tokens != 1 || reports != 1 {
				t.Fatalf("unexpected counts: %v %d %d %d", err, calls, tokens, reports)
			}
		})
	}
}
func TestObservedLoopPrecancellationMakesNoClaimOrReport(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	fail := func() { t.Fatal("pre-cancelled work") }
	if err := runObservedRefundLoop(ctx, func(context.Context, string) error { fail(); return nil }, func() (string, error) { fail(); return "", nil }, func(context.Context, string, time.Duration) error { fail(); return nil }); err != nil {
		t.Fatal(err)
	}
}
