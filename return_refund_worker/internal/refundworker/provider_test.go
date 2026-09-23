package refundworker

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stripe/stripe-go/v86"
)

func TestStripeOfficialRefundCreateAndRetrieve(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer sk_test_local" {
			t.Fatal("Stripe authorization missing")
		}
		if r.Method == http.MethodPost {
			if r.URL.Path != "/v1/refunds" || r.Header.Get("Idempotency-Key") != "return-effect-key-0001" {
				t.Fatalf("unexpected create request %s %s", r.Method, r.URL.Path)
			}
			body, _ := io.ReadAll(r.Body)
			form, err := url.ParseQuery(string(body))
			if err != nil {
				t.Fatal(err)
			}
			if form.Get("amount") != "125050" || form.Get("payment_intent") != "pi_123" || form.Get("reason") != "requested_by_customer" || form.Get("metadata[return_request_id]") != "effect-refund" {
				t.Fatalf("unexpected Stripe form: %s", body)
			}
		} else if r.Method != http.MethodGet || r.URL.Path != "/v1/refunds/re_123" {
			t.Fatalf("unexpected retrieve request %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"id":"re_123","object":"refund","amount":125050,"currency":"ars","payment_intent":"pi_123","status":"pending"}`)
	}))
	defer server.Close()
	backend := stripe.GetBackendWithConfig(stripe.APIBackend, &stripe.BackendConfig{URL: stripe.String(server.URL), HTTPClient: server.Client(), MaxNetworkRetries: stripe.Int64(0)})
	provider, err := newStripeProvider("sk_test_local", backend)
	if err != nil {
		t.Fatal(err)
	}
	request := Refund{RequestID: "effect-refund", Provider: "stripe", ProviderPaymentReference: "pi_123", IdempotencyKey: "return-effect-key-0001", Currency: "ARS", AmountMinorUnits: 125050}
	created, err := provider.Create(context.Background(), request)
	if err != nil || created.ProviderRefundReference != "re_123" || created.ProviderStatus != "pending" {
		t.Fatalf("created=%+v err=%v", created, err)
	}
	request.ProviderRefundReference = created.ProviderRefundReference
	retrieved, err := provider.Retrieve(context.Background(), request)
	if err != nil || retrieved != created {
		t.Fatalf("retrieved=%+v err=%v", retrieved, err)
	}
}

type requesterFunc func(*http.Request) (*http.Response, error)

func (f requesterFunc) Do(request *http.Request) (*http.Response, error) { return f(request) }

func TestMercadoPagoOfficialRefundCreateAndRetrieve(t *testing.T) {
	paymentReads := 0
	transport := requesterFunc(func(request *http.Request) (*http.Response, error) {
		if request.Header.Get("Authorization") != "Bearer TEST-local" {
			t.Fatal("Mercado Pago authorization missing")
		}
		path := request.URL.Path
		var response string
		switch {
		case request.Method == http.MethodGet && path == "/v1/payments/7186040733":
			paymentReads++
			if paymentReads == 1 {
				response = `{"id":7186040733,"status":"approved","captured":true,"currency_id":"ARS","transaction_amount":2000.00}`
			} else {
				response = `{"id":7186040733,"status":"refunded","captured":false,"currency_id":"ARS","transaction_amount":2000.00,"transaction_amount_refunded":1250.50}`
			}
		case request.Method == http.MethodPost && path == "/v1/payments/7186040733/refunds":
			if request.Header.Get("X-Idempotency-Key") != "return-effect-key-0002" {
				t.Fatal("Mercado Pago idempotency key missing")
			}
			var body map[string]any
			if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
				t.Fatal(err)
			}
			if body["amount"] != 1250.5 {
				t.Fatalf("unexpected amount: %#v", body)
			}
			response = `{"id":1622029222,"payment_id":7186040733,"status":"approved","amount":1250.50}`
		case request.Method == http.MethodGet && path == "/v1/payments/7186040733/refunds/1622029222":
			response = `{"id":1622029222,"payment_id":7186040733,"status":"approved","amount":1250.50}`
		default:
			t.Fatalf("unexpected Mercado Pago request %s %s", request.Method, path)
		}
		return &http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(response)), Request: request}, nil
	})
	provider, err := newMercadoPagoProvider("TEST-local", transport)
	if err != nil {
		t.Fatal(err)
	}
	request := Refund{RequestID: "effect-refund", Provider: "mercado_pago", ProviderPaymentReference: "7186040733", IdempotencyKey: "return-effect-key-0002", Currency: "ARS", AmountMinorUnits: 125050}
	created, err := provider.Create(context.Background(), request)
	if err != nil || created.ProviderRefundReference != "1622029222" || created.ProviderStatus != "approved" {
		t.Fatalf("created=%+v err=%v", created, err)
	}
	request.ProviderRefundReference = created.ProviderRefundReference
	retrieved, err := provider.Retrieve(context.Background(), request)
	if err != nil || retrieved != created {
		t.Fatalf("retrieved=%+v err=%v", retrieved, err)
	}
}

func TestProvidersFailClosed(t *testing.T) {
	stripeProvider, _ := NewStripeProvider("sk_test")
	if _, err := stripeProvider.Create(context.Background(), Refund{Provider: "stripe", ProviderPaymentReference: "invented", IdempotencyKey: "return-effect-key-0003", Currency: "ARS", AmountMinorUnits: 1}); err == nil {
		t.Fatal("invalid Stripe payment reference accepted")
	}
	mpProvider, _ := NewMercadoPagoProvider("TEST")
	if _, err := mpProvider.Create(context.Background(), Refund{Provider: "mercado_pago", ProviderPaymentReference: "1", IdempotencyKey: "return-effect-key-0004", Currency: "USD", AmountMinorUnits: 1}); err == nil {
		t.Fatal("non-ARS Mercado Pago refund accepted without an exact exponent contract")
	}
}
