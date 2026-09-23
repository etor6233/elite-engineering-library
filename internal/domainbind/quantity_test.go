package domainbind

import (
	"context"
	"elite.local/enterprise/internal/agenttools"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestQuoteQuantityCannotBeSilentlyDiscarded(t *testing.T) {
	calls := 0
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls++; w.WriteHeader(201) }))
	defer s.Close()
	g, e := NewGateway(testConfig(s.URL))
	if e != nil {
		t.Fatal(e)
	}
	for _, n := range []int{0, -1, 2, 100} {
		if _, e = g.CreateQuoteCommand(context.Background(), "tenant-id", CommandIdentity{IdempotencyKey: "synthetic-quantity-key", OccurredAt: time.Now()}, agenttools.QuoteInput{Product: "scooter", Quantity: n}); e == nil {
			t.Fatal("quantity accepted", n)
		}
	}
	if calls != 0 {
		t.Fatal("invalid quantities caused effects", calls)
	}
}
