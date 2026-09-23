package officialpayments

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/stripe/stripe-go/v86"
	stripewebhook "github.com/stripe/stripe-go/v86/webhook"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func checkoutFixtureRequest() CheckoutRequest {
	return CheckoutRequest{OrderID: "order-1", PaymentAttemptID: "attempt-1", AmountMinor: 125050, Currency: "ARS", MinorUnitExponent: 2, DisplayName: "Order order-1", SuccessURL: "https://shop.example.test/payment-return", CancelURL: "https://shop.example.test/payment-cancel", NotificationURL: "https://api.example.test/payment-webhook", IdempotencyKey: "attempt-1-checkout", ExpiresAt: time.Now().Add(time.Hour)}
}
func TestStripeHostedCheckoutSDKCreatesAndRetrievesBoundSession(t *testing.T) {
	calls := []string{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, r.Method)
		if r.Method == "POST" {
			if err := r.ParseForm(); err != nil {
				t.Error(err)
			}
			for key, want := range map[string]string{"mode": "payment", "client_reference_id": "attempt-1", "metadata[order_id]": "order-1", "metadata[payment_attempt_id]": "attempt-1", "payment_intent_data[metadata][payment_attempt_id]": "attempt-1", "line_items[0][price_data][unit_amount]": "125050", "line_items[0][price_data][currency]": "ars", "line_items[0][quantity]": "1"} {
				if r.Form.Get(key) != want {
					t.Errorf("%s=%s", key, r.Form.Get(key))
				}
			}
			if r.Header.Get("Idempotency-Key") != "attempt-1-checkout" {
				t.Error("missing idempotency")
			}
		}
		if r.Method == "GET" && r.URL.Path != "/v1/checkout/sessions/cs_test_abc" {
			t.Errorf("unexpected GET %s", r.URL)
		}
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"id":"cs_test_abc","object":"checkout.session","mode":"payment","amount_total":125050,"currency":"ars","client_reference_id":"attempt-1","metadata":{"order_id":"order-1","payment_attempt_id":"attempt-1"},"status":"open","payment_status":"unpaid","livemode":false,"url":"https://checkout.stripe.com/c/pay/cs_test_abc#provider-fragment","payment_intent":"pi_abc","expires_at":2000000000}`)
	}))
	defer server.Close()
	backend := stripe.GetBackendWithConfig(stripe.APIBackend, &stripe.BackendConfig{URL: stripe.String(server.URL), HTTPClient: server.Client(), MaxNetworkRetries: stripe.Int64(0)})
	c, _ := newStripeIntentClient("sk_test_fixture", backend)
	created, err := c.CreateCheckout(context.Background(), checkoutFixtureRequest())
	if err != nil {
		t.Fatal(err)
	}
	read, err := c.RetrieveCheckout(context.Background(), created.SessionID)
	if err != nil {
		t.Fatal(err)
	}
	if read != created || read.PaymentReference != "pi_abc" || read.PaymentAttemptID != "attempt-1" || len(calls) != 2 {
		t.Fatalf("%+v calls%v", read, calls)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
func TestMercadoPagoHostedCheckoutSDKCreateAndGetNoPayerData(t *testing.T) {
	calls := 0
	client := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		calls++
		if r.URL.Host != "api.mercadopago.com" || !strings.HasPrefix(r.URL.Path, "/checkout/preferences") {
			t.Errorf("wrong endpoint %s", r.URL)
		}
		if r.Method == "POST" {
			var body map[string]any
			if json.NewDecoder(r.Body).Decode(&body) != nil {
				t.Error("bad body")
			}
			metadata := body["metadata"].(map[string]any)
			if body["external_reference"] != "order-1" || metadata["payment_attempt_id"] != "attempt-1" || body["payer"] != nil {
				t.Errorf("wrong contract %v", body)
			}
			if r.Header.Get("X-Idempotency-Key") != "attempt-1-checkout" {
				t.Error("missing key")
			}
		}
		body := `{"id":"123-abc","external_reference":"order-1","metadata":{"order_id":"order-1","payment_attempt_id":"attempt-1"},"collector_id":123,"items":[{"currency_id":"ARS","quantity":1,"unit_price":1250.50}],"init_point":"https://www.mercadopago.com.ar/checkout/v1/redirect?pref_id=123-abc","sandbox_init_point":"https://sandbox.mercadopago.com.ar/checkout/v1/redirect?pref_id=123-abc","expiration_date_to":"2030-01-01T00:00:00Z"}`
		return &http.Response{StatusCode: 200, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body)), Request: r}, nil
	})}
	c, _ := NewMercadoPagoCheckoutClientWithHTTPClient("TEST-fixture", true, client)
	created, err := c.CreateCheckout(context.Background(), checkoutFixtureRequest())
	if err != nil {
		t.Fatal(err)
	}
	read, err := c.RetrieveCheckout(context.Background(), created.SessionID, 2)
	if err != nil {
		t.Fatal(err)
	}
	if read != created || read.Currency != "ARS" || read.CollectorID != "123" || read.AmountMinor != 125050 || read.LiveMode || calls != 2 {
		t.Fatalf("%+v calls%d", read, calls)
	}
}
func TestHostedCheckoutRejectsUntrustedRedirectAndIncompleteRequest(t *testing.T) {
	for _, raw := range []string{"http://checkout.stripe.com/x", "https://checkout.stripe.com.evil.test/x", "https://user@checkout.stripe.com/x", "https://checkout.stripe.com:444/x", "javascript:alert(1)"} {
		if hostedURL(raw, "checkout.stripe.com") {
			t.Fatalf("unsafe URL%s", raw)
		}
	}
	req := checkoutFixtureRequest()
	req.PaymentAttemptID = ""
	if validCheckout(req) {
		t.Fatal("missing attempt accepted")
	}
	req = checkoutFixtureRequest()
	req.ExpiresAt = time.Time{}
	if validCheckout(req) {
		t.Fatal("unbounded checkout accepted")
	}
}
func TestStripeRefundDisputeAndCheckoutCallbacksRouteOnlyToOfficialGET(t *testing.T) {
	for _, row := range []struct{ kind, object, id, intent, resource, want string }{{"checkout.session.completed", "checkout.session", "cs_test_abc", "", "checkout_session", "cs_test_abc"}, {"charge.refunded", "charge", "ch_abc", "pi_abc", "payment_intent", "pi_abc"}, {"charge.dispute.created", "dispute", "dp_abc", "pi_abc", "payment_intent", "pi_abc"}} {
		body := []byte(fmt.Sprintf(`{"id":"evt_abc","object":"event","api_version":%q,"type":%q,"livemode":false,"data":{"object":{"id":%q,"object":%q,"payment_intent":%q}}}`, stripe.APIVersion, row.kind, row.id, row.object, row.intent))
		signed := stripewebhook.GenerateTestSignedPayload(&stripewebhook.UnsignedPayload{Payload: body, Secret: "secret", Timestamp: time.Now()})
		n, err := VerifyStripePaymentNotification(body, signed.Header, "secret", time.Minute, false, "")
		if err != nil || n.ResourceType != row.resource || n.ProviderReference != row.want {
			t.Fatalf("%s %+v %v", row.kind, n, err)
		}
	}
}
func TestStripeDisputedSnapshotCannotMasqueradeAsUncontestedReceipt(t *testing.T) {
	c, _ := NewStripeIntentClientWithHTTPClient("sk_test_fixture", &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		body := `{"id":"pi_abc","amount":100,"amount_received":100,"currency":"usd","status":"succeeded","metadata":{"order_id":"order-1","payment_attempt_id":"attempt-1"},"latest_charge":{"id":"ch_abc","amount":100,"amount_captured":100,"amount_refunded":0,"currency":"usd","captured":true,"disputed":true}}`
		return &http.Response{StatusCode: 200, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body)), Request: r}, nil
	})})
	snapshot, err := c.RetrieveIntent(context.Background(), "pi_abc")
	if err != nil || !snapshot.Disputed || snapshot.PaymentAttemptID != "attempt-1" {
		t.Fatalf("%+v %v", snapshot, err)
	}
}

func TestStripeAccountProbeUsesCredentialBoundOfficialEndpoint(t *testing.T) {
	var seen bool
	c, _ := NewStripeIntentClientWithHTTPClient("sk_test_fixture", &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		seen = true
		if r.Method != "GET" || r.URL.String() != "https://api.stripe.com/v1/account" || r.Header.Get("Authorization") != "Bearer sk_test_fixture" {
			t.Errorf("wrong account probe %s", r.URL)
		}
		return &http.Response{StatusCode: 200, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(`{"id":"acct_fixture","object":"account","email":"not-exported@example.test"}`)), Request: r}, nil
	})})
	id, err := c.ProbeAccount(context.Background())
	if err != nil || id != "acct_fixture" || !seen {
		t.Fatalf("probe %s %v", id, err)
	}
}
