package aifoundation

import (
	"context"
	"math"
	"testing"
)

func TestEvalRequiredCasesCannotHideBehindAverage(t *testing.T) {
	s := EvalSuite{ID: "reference-v1", MinPassRate: 0.5, Cases: []EvalCase{{ID: "ordinary", Assert: func(CompletionResponse) bool { return true }}, {ID: "tenant-isolation", Required: true, Assert: func(CompletionResponse) bool { return false }}}}
	r, e := s.Run(fakeProvider{})
	if e != nil || r.GatePassed || len(r.Cases) != 2 || r.Cases[1].Passed {
		t.Fatal(r, e)
	}
}
func TestEvalRejectsUnboundedOrAmbiguousDecision(t *testing.T) {
	yes := func(CompletionResponse) bool { return true }
	for _, s := range []EvalSuite{
		{ID: "nan", MinPassRate: math.NaN(), Cases: []EvalCase{{ID: "one", Assert: yes}}},
		{ID: "nil", MinPassRate: 1, Cases: []EvalCase{{ID: "one"}}},
		{ID: "duplicate", MinPassRate: 1, Cases: []EvalCase{{ID: "same", Assert: yes}, {ID: "same", Assert: yes}}},
	} {
		if _, e := s.Run(fakeProvider{}); e == nil {
			t.Fatal("invalid eval accepted", s.ID)
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	s := EvalSuite{ID: "canceled", MinPassRate: 1, Cases: []EvalCase{{ID: "one", Assert: yes}}}
	if _, e := s.RunContext(ctx, fakeProvider{}); e != context.Canceled {
		t.Fatal(e)
	}
}
