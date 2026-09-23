package refundworker

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	stripe "github.com/stripe/stripe-go/v86"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestBusinessPreparationErrorPolicy(t *testing.T) {
	for _, tc := range []struct {
		name          string
		err           error
		attempt       int
		outcome, code string
	}{
		{"serialization", &pgconn.PgError{Code: "40001"}, 1, "retry", "REFUND_PREPARATION_TRANSIENT"},
		{"deadlock", &pgconn.PgError{Code: "40P01"}, 1, "retry", "REFUND_PREPARATION_TRANSIENT"},
		{"connection", &pgconn.PgError{Code: "08006"}, 1, "retry", "REFUND_PREPARATION_TRANSIENT"},
		{"wrapped", fmt.Errorf("prepare: %w", &pgconn.PgError{Code: "40001"}), 1, "retry", "REFUND_PREPARATION_TRANSIENT"},
		{"mapping", ErrMappingConflict, 1, "blocked", "REFUND_MAPPING_CONFLICT"},
		{"syntax", &pgconn.PgError{Code: "42601"}, 1, "blocked", "REFUND_PREPARATION_FAILED"},
		{"unique", &pgconn.PgError{Code: "23505"}, 1, "blocked", "REFUND_PREPARATION_FAILED"},
		{"unknown", errors.New("configuration invalid"), 1, "blocked", "REFUND_PREPARATION_FAILED"},
		{"exhausted", &pgconn.PgError{Code: "40001"}, 5, "blocked", "REFUND_PREPARATION_EXHAUSTED"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			store := &fakeStore{work: &Work{TenantID: "tenant", RequestID: "request", Attempt: tc.attempt, ClaimToken: "claim"}, refund: baseRefund(), prepareErr: tc.err}
			provider := &fakeProvider{}
			p, _ := NewProcessor(store, map[string]Provider{"stripe": provider}, "worker", time.Minute, time.Second, 5)
			err := p.Step(context.Background(), "claim")
			if !errors.Is(err, tc.err) || !store.finished || store.outcome != tc.outcome || store.code != tc.code || provider.created || provider.retrieved {
				t.Fatalf("err=%v outcome=%s code=%s provider=%+v; want %s/%s", err, store.outcome, store.code, provider, tc.outcome, tc.code)
			}
		})
	}
}

func TestBusinessPreparationRecoveryDispatchesOnce(t *testing.T) {
	store := &fakeStore{work: &Work{TenantID: "tenant", RequestID: "request", Attempt: 1, ClaimToken: "claim-1"}, refund: baseRefund(), prepareErr: &pgconn.PgError{Code: "40001"}}
	provider := &fakeProvider{result: ProviderResult{ProviderPaymentReference: "pi_1", ProviderRefundReference: "re_1", ProviderStatus: "succeeded", Currency: "ARS", AmountMinorUnits: 1000}}
	p, _ := NewProcessor(store, map[string]Provider{"stripe": provider}, "worker", time.Minute, time.Second, 5)
	if err := p.Step(context.Background(), "claim-1"); err == nil || store.outcome != "retry" || provider.created {
		t.Fatalf("first err=%v outcome=%s created=%v", err, store.outcome, provider.created)
	}
	store.prepareErr = nil
	store.work.Attempt = 2
	store.work.ClaimToken = "claim-2"
	if err := p.Step(context.Background(), "claim-2"); err != nil || store.outcome != "succeeded" || !provider.created {
		t.Fatalf("recovery err=%v outcome=%s created=%v", err, store.outcome, provider.created)
	}
}

func TestBusinessStripeRefundIdentifier(t *testing.T) {
	for _, ref := range []string{"ch_original", "pi_related"} {
		t.Run(ref, func(t *testing.T) {
			calls := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls++
				if r.Method == http.MethodPost {
					body, _ := io.ReadAll(r.Body)
					form, err := url.ParseQuery(string(body))
					if err != nil {
						t.Fatal(err)
					}
					field := "charge"
					if ref == "pi_related" {
						field = "payment_intent"
					}
					if form.Get(field) != ref || r.Header.Get("Idempotency-Key") != "return-effect-key-0001" {
						t.Errorf("wrong identity/idempotency")
					}
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = io.WriteString(w, `{"id":"re_fixture","object":"refund","amount":1000,"currency":"ars","charge":"ch_original","payment_intent":"pi_related","status":"succeeded"}`)
			}))
			defer server.Close()
			backend := stripe.GetBackendWithConfig(stripe.APIBackend, &stripe.BackendConfig{URL: stripe.String(server.URL), HTTPClient: server.Client(), MaxNetworkRetries: stripe.Int64(0)})
			p, err := newStripeProvider("sk_test_local", backend)
			if err != nil {
				t.Fatal(err)
			}
			value := baseRefund()
			value.ProviderPaymentReference = ref
			got, err := p.Create(context.Background(), value)
			if err != nil || got.ProviderPaymentReference != ref {
				t.Fatalf("create got=%+v err=%v", got, err)
			}
			value.ProviderRefundReference = got.ProviderRefundReference
			read, err := p.Retrieve(context.Background(), value)
			if err != nil || read != got || calls != 2 {
				t.Fatalf("retrieve=%+v err=%v calls=%d", read, err, calls)
			}
		})
	}
}

func TestBusinessStripeRefundIdentityFailClosed(t *testing.T) {
	for _, tc := range []struct {
		name, ref, charge, pi, currency string
		amount                          int64
	}{
		{"charge-missing", "ch_original", "", "pi_related", "ars", 1000},
		{"charge-wrong", "ch_original", "ch_other", "pi_related", "ars", 1000},
		{"intent-missing", "pi_related", "ch_original", "", "ars", 1000},
		{"intent-wrong", "pi_related", "ch_original", "pi_other", "ars", 1000},
		{"currency", "ch_original", "ch_original", "pi_related", "usd", 1000},
		{"amount", "ch_original", "ch_original", "pi_related", "ars", 999},
	} {
		t.Run(tc.name, func(t *testing.T) {
			v := baseRefund()
			v.ProviderPaymentReference = tc.ref
			r := &stripe.Refund{ID: "re_fixture", Amount: tc.amount, Currency: stripe.Currency(tc.currency)}
			if tc.charge != "" {
				r.Charge = &stripe.Charge{ID: tc.charge}
			}
			if tc.pi != "" {
				r.PaymentIntent = &stripe.PaymentIntent{ID: tc.pi}
			}
			if _, err := stripeResult(v, r); !errors.Is(err, ErrResponseMismatch) {
				t.Fatalf("got %v", err)
			}
		})
	}
}
