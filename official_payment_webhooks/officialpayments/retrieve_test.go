package officialpayments

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stripe/stripe-go/v86"
	stripewebhook "github.com/stripe/stripe-go/v86/webhook"
)

func TestStripeRetrieveUsesOfficialGETAndExcludesSecrets(t *testing.T) {
	var calls atomic.Int32
	response := `{"id":"pi_123","object":"payment_intent","amount":1099,"amount_received":1099,"amount_capturable":0,"currency":"usd","status":"succeeded","livemode":false,"metadata":{"order_id":"order-123"},"client_secret":"never-export-secret","latest_charge":{"id":"ch_123","amount":1099,"amount_captured":1099,"amount_refunded":100,"currency":"usd","captured":true,"livemode":false}}`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		if r.Method != "GET" || r.URL.Path != "/v1/payment_intents/pi_123" || r.URL.Query().Get("expand[0]") != "latest_charge" || r.Header.Get("Authorization") != "Bearer sk_test_local" {
			t.Errorf("unexpected SDK GET: %s %s", r.Method, r.URL)
		}
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, response)
	}))
	defer server.Close()
	backend := stripe.GetBackendWithConfig(stripe.APIBackend, &stripe.BackendConfig{URL: stripe.String(server.URL), HTTPClient: server.Client(), MaxNetworkRetries: stripe.Int64(0)})
	client, _ := newStripeIntentClient("sk_test_local", backend)
	snapshot, err := client.RetrieveIntent(context.Background(), "pi_123")
	if err != nil {
		t.Fatal(err)
	}
	if snapshot.ProviderCode != "stripe" || snapshot.OrderID != "order-123" || snapshot.AmountMinor != 1099 || snapshot.ReceivedMinor != 1099 || snapshot.RefundedMinor != 100 || snapshot.Currency != "USD" {
		t.Fatalf("snapshot=%+v", snapshot)
	}
	data, _ := json.Marshal(snapshot)
	if strings.Contains(string(data), "secret") || len(snapshot.SHA256()) != 64 || calls.Load() != 1 {
		t.Fatal("observation leak or extra request")
	}
	for _, id := range []string{"", "pi_123/other", "../pi_123", "pi_123?secret=x"} {
		if _, err = client.RetrieveIntent(context.Background(), id); err == nil {
			t.Fatal("invalid id accepted")
		}
	}
	if calls.Load() != 1 {
		t.Fatal("invalid identifier sent to provider")
	}
	response = strings.Replace(response, `"amount_captured":1099`, `"amount_captured":1098`, 1)
	if _, err = client.RetrieveIntent(context.Background(), "pi_123"); err == nil {
		t.Fatal("inconsistent charge accepted")
	}
}

func TestMercadoPagoRetrieveOfficialGETAndExactMinorUnits(t *testing.T) {
	calls := 0
	transport := requesterFunc(func(r *http.Request) (*http.Response, error) {
		calls++
		if r.Method != "GET" || r.URL.String() != "https://api.mercadopago.com/v1/payments/987" || r.Header.Get("Authorization") != "Bearer TEST-local" {
			t.Errorf("unexpected SDK GET: %s %s", r.Method, r.URL)
		}
		body := `{"id":987,"transaction_amount":1250.50,"transaction_amount_refunded":0.01,"currency_id":"ARS","external_reference":"order-456","status":"approved","captured":true,"collector_id":123,"live_mode":false,"payer":{"email":"private@example.test"}}`
		return &http.Response{StatusCode: 200, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body)), Request: r}, nil
	})
	client, _ := newMercadoPagoPaymentClient("TEST-local", transport)
	snapshot, err := client.RetrievePayment(context.Background(), "987", 2)
	if err != nil {
		t.Fatal(err)
	}
	if snapshot.ProviderCode != "mercadopago" || snapshot.AmountMinor != 125050 || snapshot.ReceivedMinor != 125050 || snapshot.RefundedMinor != 1 || snapshot.OrderID != "order-456" || snapshot.CollectorID != "123" {
		t.Fatalf("snapshot=%+v", snapshot)
	}
	data, _ := json.Marshal(snapshot)
	if strings.Contains(string(data), "private") {
		t.Fatal("PII leaked")
	}
	for _, id := range []string{"", "0", "-1", "1/2", "001", "9223372036854775808"} {
		if _, err = client.RetrievePayment(context.Background(), id, 2); err == nil {
			t.Fatal("invalid numeric id")
		}
	}
	if calls != 1 {
		t.Fatal("invalid identifier sent to provider")
	}
}

func TestDecimalConversionRejectsRoundingAndUnsupportedRange(t *testing.T) {
	for _, value := range []float64{-1, math.Inf(1), math.NaN(), 1.001, 1e20} {
		if _, err := decimalMinor(value, 2); err == nil {
			t.Fatalf("accepted %v", value)
		}
	}
	for _, row := range []struct {
		value float64
		exp   int
		want  int64
	}{{1.23, 2, 123}, {0, 2, 0}, {1.001, 3, 1001}, {123, 0, 123}} {
		got, err := decimalMinor(row.value, row.exp)
		if err != nil || got != row.want {
			t.Fatalf("got %d %v", got, err)
		}
	}
}

func TestStripePaymentNotificationBindsModeAccountAndResource(t *testing.T) {
	body := []byte(fmt.Sprintf(`{"id":"evt_123","object":"event","api_version":%q,"type":"payment_intent.succeeded","livemode":false,"data":{"object":{"id":"pi_123","object":"payment_intent"}}}`, stripe.APIVersion))
	signed := stripewebhook.GenerateTestSignedPayload(&stripewebhook.UnsignedPayload{Payload: body, Secret: "whsec_test", Timestamp: time.Now()})
	notification, err := VerifyStripePaymentNotification(body, signed.Header, "whsec_test", 5*time.Minute, false, "")
	if err != nil || notification.ProviderReference != "pi_123" || notification.ProviderEventID != "evt_123" || notification.ProviderCode != "stripe" {
		t.Fatalf("notification=%+v err=%v", notification, err)
	}
	if _, err = VerifyStripePaymentNotification(body, signed.Header, "whsec_test", 5*time.Minute, true, ""); err == nil {
		t.Fatal("live mismatch accepted")
	}
	if _, err = VerifyStripePaymentNotification(body, signed.Header, "whsec_test", 5*time.Minute, false, "acct_wrong"); err == nil {
		t.Fatal("account mismatch accepted")
	}
	if _, err = VerifyStripePaymentNotification(body, signed.Header, "wrong", 5*time.Minute, false, ""); err == nil {
		t.Fatal("bad signature accepted")
	}
}

func TestMercadoPagoNotificationUsesSignedResourceNotBodyFinancialClaims(t *testing.T) {
	body := []byte(`{"type":"payment","data":{"id":"987"},"status":"approved","transaction_amount":999999,"user_id":666}`)
	header := mercadoPagoHeader("987", "request-1", "secret", time.Now())
	notification, err := VerifyMercadoPagoPaymentNotification(body, header, "request-1", "987", "secret", 5*time.Minute)
	if err != nil || notification.ProviderCode != "mercadopago" || notification.ProviderReference != "987" || notification.Type != "payment.reconciliation-requested" {
		t.Fatalf("%+v %v", notification, err)
	}
	data, _ := json.Marshal(notification)
	if strings.Contains(string(data), "approved") || strings.Contains(string(data), "999999") {
		t.Fatal("unsigned financial claims propagated")
	}
	if _, err = VerifyMercadoPagoPaymentNotification(body, header, "other-request", "987", "secret", 5*time.Minute); err == nil {
		t.Fatal("request substitution accepted")
	}
	if _, err = VerifyMercadoPagoPaymentNotification(body, header, "request-1", "988", "secret", 5*time.Minute); err == nil {
		t.Fatal("resource substitution accepted")
	}
}
