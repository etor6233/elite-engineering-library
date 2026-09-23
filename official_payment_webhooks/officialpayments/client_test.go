package officialpayments

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/stripe/stripe-go/v86"
)

func TestStripeBackendDoesNotLogMalformedResponseBody(t *testing.T) {
	const marker = "ELITE_PAYMENT_LOG_FIXTURE"
	if os.Getenv(marker) == "child" {
		transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
			if r.URL.String() != "https://api.stripe.com/v1/account" {
				t.Fatal("unexpected provider request")
			}
			return &http.Response{StatusCode: 200, Header: http.Header{"Content-Type": []string{"application/json"}}, Body: io.NopCloser(strings.NewReader(`{"PRIVATE_RESPONSE_SENTINEL":`)), Request: r}, nil
		})
		client, err := NewStripeIntentClientWithHTTPClient("sk_test_fixture", &http.Client{Transport: transport})
		if err != nil {
			t.Fatal("fixture construction failed")
		}
		if _, err = client.ProbeAccount(context.Background()); err == nil {
			t.Fatal("malformed response accepted")
		}
		return
	}
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	command := exec.CommandContext(ctx, executable, "-test.run=^TestStripeBackendDoesNotLogMalformedResponseBody$")
	command.Env = append(os.Environ(), marker+"=child")
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatal("fixture subprocess failed")
	}
	if strings.Contains(string(output), "PRIVATE_RESPONSE_SENTINEL") || strings.Contains(string(output), "body sample") {
		t.Fatal("provider raw response reached process logs")
	}
}

func TestStripeCreateIntentOfficialClientContract(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/payment_intents" {
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer sk_test_local" {
			t.Fatalf("authorization header missing")
		}
		if r.Header.Get("Idempotency-Key") != "order-123-attempt-1" {
			t.Fatalf("idempotency header missing")
		}
		body, _ := io.ReadAll(r.Body)
		form, err := url.ParseQuery(string(body))
		if err != nil {
			t.Fatal(err)
		}
		if form.Get("amount") != "1099" || form.Get("currency") != "usd" || form.Get("metadata[order_id]") != "order-123" {
			t.Fatalf("unexpected form: %s", body)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"id":"pi_123","object":"payment_intent","status":"requires_payment_method","client_secret":"pi_secret"}`)
	}))
	defer server.Close()
	backend := stripe.GetBackendWithConfig(stripe.APIBackend, &stripe.BackendConfig{URL: stripe.String(server.URL), HTTPClient: server.Client(), MaxNetworkRetries: stripe.Int64(0)})
	client, err := newStripeIntentClient("sk_test_local", backend)
	if err != nil {
		t.Fatal(err)
	}
	result, err := client.CreateIntent(context.Background(), StripeIntentRequest{AmountMinor: 1099, Currency: "usd", OrderID: "order-123", IdempotencyKey: "order-123-attempt-1"})
	if err != nil {
		t.Fatal(err)
	}
	if result.ID != "pi_123" || result.Status != "requires_payment_method" || result.ClientSecret != "pi_secret" {
		t.Fatalf("unexpected result: %#v", result)
	}
}

type requesterFunc func(*http.Request) (*http.Response, error)

func (f requesterFunc) Do(r *http.Request) (*http.Response, error) { return f(r) }

func TestMercadoPagoCreatePaymentOfficialClientContract(t *testing.T) {
	transport := requesterFunc(func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodPost || r.URL.String() != "https://api.mercadopago.com/v1/payments" {
			t.Fatalf("unexpected request %s %s", r.Method, r.URL)
		}
		if r.Header.Get("Authorization") != "Bearer TEST-local" {
			t.Fatalf("authorization header missing")
		}
		if r.Header.Get("X-Idempotency-Key") != "order-456-attempt-1" {
			t.Fatalf("idempotency header missing: %#v", r.Header)
		}
		body, _ := io.ReadAll(r.Body)
		var payload map[string]any
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatal(err)
		}
		if payload["transaction_amount"] != 1250.5 || payload["external_reference"] != "order-456" || payload["payment_method_id"] != "visa" {
			t.Fatalf("unexpected body: %s", body)
		}
		return &http.Response{StatusCode: 201, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(`{"id":987,"status":"pending"}`)), Request: r}, nil
	})
	client, err := newMercadoPagoPaymentClient("TEST-local", transport)
	if err != nil {
		t.Fatal(err)
	}
	result, err := client.CreatePayment(context.Background(), MercadoPagoPaymentRequest{AmountMinor: 125050, MinorUnitExponent: 2, Currency: "ars", OrderID: "order-456", PaymentMethodID: "visa", PaymentToken: "card-token", PayerEmail: "payer@example.test", NotificationURL: "https://example.test/webhooks/mp", IdempotencyKey: "order-456-attempt-1"})
	if err != nil {
		t.Fatal(err)
	}
	if result.ID != "987" || result.Status != "pending" {
		t.Fatalf("unexpected result: %#v", result)
	}
}

func TestOutboundPaymentInputFailsClosed(t *testing.T) {
	stripeClient, _ := NewStripeIntentClient("sk_test")
	if _, err := stripeClient.CreateIntent(context.Background(), StripeIntentRequest{AmountMinor: 1, Currency: "USD", OrderID: "order", IdempotencyKey: "id"}); err == nil {
		t.Fatal("uppercase currency accepted")
	}
	mpClient, _ := NewMercadoPagoPaymentClient("TEST")
	if _, err := mpClient.CreatePayment(context.Background(), MercadoPagoPaymentRequest{AmountMinor: 1, MinorUnitExponent: 4}); err == nil {
		t.Fatal("invalid exponent accepted")
	}
}
