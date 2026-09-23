package paymentbridge

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"elite.local/enterprise/internal/bcamounts"
	"elite.local/enterprise/internal/providerintegration"
	"github.com/stripe/stripe-go/v86"
	stripewebhook "github.com/stripe/stripe-go/v86/webhook"
)

type localProviderTransport struct {
	destination  *url.URL
	providerHost string
}

func (t localProviderTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	if r.URL.Scheme != "https" || r.URL.Host != t.providerHost {
		return nil, fmt.Errorf("fixture forbids network host")
	}
	cloned := r.Clone(r.Context())
	u := *r.URL
	u.Scheme = t.destination.Scheme
	u.Host = t.destination.Host
	cloned.URL = &u
	cloned.Host = ""
	return http.DefaultTransport.RoundTrip(cloned)
}
func fixtureHTTPClient(t *testing.T, host string, handler http.HandlerFunc) *http.Client {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	parsed, _ := url.Parse(server.URL)
	return &http.Client{Transport: localProviderTransport{parsed, host}}
}
func TestSDKDriverHostedSessionCallbackAndRefundDisputeUseOfficialHTTP(t *testing.T) {
	var mu sync.Mutex
	account := "acct_fixture"
	refund := int64(0)
	disputed := false
	posts := 0
	expires := time.Now().Add(time.Hour).UTC().Truncate(time.Second)
	client := fixtureHTTPClient(t, "api.stripe.com", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		if r.Header.Get("Authorization") != "Bearer sk_test_fixture" {
			t.Error("missing fixture credential")
		}
		switch r.URL.Path {
		case "/v1/account":
			if r.Method != "GET" {
				t.Error("account write")
			}
			fmt.Fprintf(w, `{"id":%q}`, account)
		case "/v1/checkout/sessions", "/v1/checkout/sessions/cs_test_fixture":
			if r.Method == "POST" {
				posts++
				r.ParseForm()
				if r.Form.Get("metadata[payment_attempt_id]") != "attempt-1" || r.Form.Get("payment_intent_data[metadata][order_id]") != "order-1" || r.Form.Get("expires_at") != strconv.FormatInt(expires.Unix(), 10) || r.Header.Get("Idempotency-Key") != "attempt-1" {
					t.Errorf("wrong stable checkout request")
				}
			}
			fmt.Fprintf(w, `{"id":"cs_test_fixture","object":"checkout.session","mode":"payment","amount_total":125050,"currency":"ars","client_reference_id":"attempt-1","metadata":{"order_id":"order-1","payment_attempt_id":"attempt-1"},"status":"open","payment_status":"unpaid","livemode":false,"url":"https://checkout.stripe.com/c/pay/cs_test_fixture","payment_intent":"pi_fixture","expires_at":%d}`, expires.Unix())
		case "/v1/payment_intents/pi_fixture":
			if r.Method != "GET" || r.URL.Query().Get("expand[0]") != "latest_charge" {
				t.Error("payment observation is not expanded GET")
			}
			fmt.Fprintf(w, `{"id":"pi_fixture","amount":125050,"amount_received":125050,"currency":"ars","status":"succeeded","metadata":{"order_id":"order-1","payment_attempt_id":"attempt-1"},"latest_charge":{"id":"ch_fixture","amount":125050,"amount_captured":125050,"amount_refunded":%d,"currency":"ars","captured":true,"disputed":%t}}`, refund, disputed)
		default:
			t.Errorf("unexpected path%s", r.URL)
			w.WriteHeader(500)
		}
	})
	scope := Scope{TenantID: "tenant", OrganizationID: "store", ConnectionID: "connection", ProviderCode: "stripe", AccountRef: "acct_fixture", Currency: "ARS", MinorUnitExponent: 2}
	d, err := NewSDKDriver(SDKDriverConfig{Scope: scope, SuccessURL: "https://shop.example.test/return", CancelURL: "https://shop.example.test/cancel", DisplayName: "Order"}, "sk_test_fixture", client)
	if err != nil {
		t.Fatal(err)
	}
	req := Request{TenantID: scope.TenantID, OrganizationID: scope.OrganizationID, ProviderCode: scope.ProviderCode, Currency: scope.Currency, PaymentAttemptID: "attempt-1", OrderID: "order-1", CustomerSubject: "customer", AmountMinor: 125050, CheckoutExpiresAt: expires}
	created, err := d.CreateCheckout(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	inbox := &fixtureInbox{records: map[string]providerintegration.Receipt{}}
	h, _ := NewWebhook(WebhookConfig{TenantID: scope.TenantID, ConnectionID: scope.ConnectionID, ProviderCode: "stripe", Tolerance: time.Minute, MaxConcurrent: 1, Secrets: fixtureSecret{}, Inbox: inbox})
	body := []byte(fmt.Sprintf(`{"id":"evt_session","object":"event","api_version":%q,"type":"checkout.session.completed","livemode":false,"data":{"object":{"id":"cs_test_fixture","object":"checkout.session"}}}`, stripe.APIVersion))
	signed := stripewebhook.GenerateTestSignedPayload(&stripewebhook.UnsignedPayload{Payload: body, Secret: "whsec_fixture", Timestamp: time.Now()})
	callback := httptest.NewRequest("POST", "https://api.example.test/payment-webhook", strings.NewReader(string(body)))
	callback.Header.Set("Content-Type", "application/json")
	callback.Header.Set("Stripe-Signature", signed.Header)
	response := httptest.NewRecorder()
	h.ServeHTTP(response, callback)
	if response.Code != 200 {
		t.Fatal(response.Code)
	}
	var payload InboxPayload
	if json.Unmarshal(inbox.records["evt_session"].Payload, &payload) != nil || payload.Notification.ResourceType != "checkout_session" {
		t.Fatal("wrong durable notice")
	}
	observedCheckout, err := d.RetrieveCheckout(context.Background(), payload.Notification.ProviderReference)
	if err != nil || observedCheckout != created {
		t.Fatalf("checkout %+v %v", observedCheckout, err)
	}
	paid, err := d.RetrievePayment(context.Background(), observedCheckout.PaymentReference)
	if err != nil || paid.ReceivedMinor != req.AmountMinor || paid.RefundedMinor != 0 || paid.Disputed {
		t.Fatalf("paid %+v %v", paid, err)
	}
	mu.Lock()
	refund = 1
	mu.Unlock()
	refunded, err := d.RetrievePayment(context.Background(), observedCheckout.PaymentReference)
	if err != nil || refunded.RefundedMinor != 1 {
		t.Fatalf("refund %+v %v", refunded, err)
	}
	mu.Lock()
	refund = 0
	disputed = true
	mu.Unlock()
	dispute, err := d.RetrievePayment(context.Background(), observedCheckout.PaymentReference)
	if err != nil || !dispute.Disputed {
		t.Fatalf("dispute %+v %v", dispute, err)
	}
	mu.Lock()
	account = "acct_other"
	mu.Unlock()
	if _, err = d.CreateCheckout(context.Background(), req); err == nil {
		t.Fatal("wrong credential account allowed")
	}
	mu.Lock()
	postCount := posts
	mu.Unlock()
	if postCount != 1 {
		t.Fatalf("unexpected provider writes%d", postCount)
	}
}
func TestSDKDriverMercadoPagoProbesAccountBeforeHostedCheckout(t *testing.T) {
	var mu sync.Mutex
	expires := time.Now().Add(time.Hour).UTC().Truncate(time.Second)
	posts := 0
	account := 123
	client := fixtureHTTPClient(t, "api.mercadopago.com", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/users/me" {
			fmt.Fprintf(w, `{"id":%d,"country_id":"AR","site_id":"MLA"}`, account)
			return
		}
		if r.URL.Path != "/checkout/preferences" || r.Method != "POST" {
			t.Error("unexpected MP effect")
			w.WriteHeader(500)
			return
		}
		posts++
		io.Copy(io.Discard, r.Body)
		fmt.Fprintf(w, `{"id":"123-preference","external_reference":"order-1","metadata":{"order_id":"order-1","payment_attempt_id":"attempt-1"},"collector_id":123,"items":[{"currency_id":"ARS","quantity":1,"unit_price":1250.50}],"sandbox_init_point":"https://sandbox.mercadopago.com.ar/checkout/v1/redirect?pref_id=123-preference","expiration_date_to":%q}`, expires.Format(time.RFC3339))
	})
	scope := Scope{TenantID: "tenant", OrganizationID: "store", ConnectionID: "connection", ProviderCode: "mercadopago", AccountRef: "123", Currency: "ARS", MinorUnitExponent: 2}
	d, err := NewSDKDriver(SDKDriverConfig{Scope: scope, SuccessURL: "https://shop.example.test/return", CancelURL: "https://shop.example.test/cancel", NotificationURL: "https://api.example.test/payment", DisplayName: "Order"}, "TEST-fixture", client)
	if err != nil {
		t.Fatal(err)
	}
	req := Request{TenantID: "tenant", OrganizationID: "store", ProviderCode: "mercadopago", Currency: "ARS", OrderID: "order-1", PaymentAttemptID: "attempt-1", CustomerSubject: "customer", AmountMinor: 125050, CheckoutExpiresAt: expires}
	result, err := d.CreateCheckout(context.Background(), req)
	if err != nil || result.AccountRef != "123" {
		t.Fatalf("%+v %v", result, err)
	}
	mu.Lock()
	account = 456
	mu.Unlock()
	if _, err = d.CreateCheckout(context.Background(), req); err == nil {
		t.Fatal("foreign account allowed")
	}
	mu.Lock()
	postCount := posts
	mu.Unlock()
	if postCount != 1 {
		t.Fatal("unexpected repeat effect")
	}
}

var _ Driver = (*SDKDriver)(nil)

func TestExactBCAmountCheckoutWebhookFixture(t *testing.T) {
	const unitPriceMinor = int64(41_789)
	const quantity = int64(3)
	expectedMinor, err := bcamounts.LineAmount(quantity, unitPriceMinor, 0)
	if err != nil || expectedMinor != 125_367 {
		t.Fatalf("bcamounts oracle: got=%d err=%v want=125367", expectedMinor, err)
	}
	var mu sync.Mutex
	account := "acct_fixture"
	expires := time.Now().Add(time.Hour).UTC().Truncate(time.Second)
	client := fixtureHTTPClient(t, "api.stripe.com", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		if r.Header.Get("Authorization") != "Bearer sk_test_fixture" {
			t.Error("missing fixture credential")
		}
		switch r.URL.Path {
		case "/v1/account":
			fmt.Fprintf(w, `{"id":%q}`, account)
		case "/v1/checkout/sessions", "/v1/checkout/sessions/cs_test_fixture":
			if r.Method == "POST" {
				r.ParseForm()
				if r.Form.Get("line_items[0][price_data][unit_amount]") != strconv.FormatInt(expectedMinor, 10) {
					t.Errorf("checkout unit_amount=%s want=%d", r.Form.Get("line_items[0][price_data][unit_amount]"), expectedMinor)
				}
			}
			fmt.Fprintf(w, `{"id":"cs_test_fixture","object":"checkout.session","mode":"payment","amount_total":%d,"currency":"ars","client_reference_id":"attempt-1","metadata":{"order_id":"order-1","payment_attempt_id":"attempt-1"},"status":"open","payment_status":"unpaid","livemode":false,"url":"https://checkout.stripe.com/c/pay/cs_test_fixture","payment_intent":"pi_fixture","expires_at":%d}`, expectedMinor, expires.Unix())
		case "/v1/payment_intents/pi_fixture":
			if r.Method != "GET" || r.URL.Query().Get("expand[0]") != "latest_charge" {
				t.Error("payment observation is not expanded GET")
			}
			fmt.Fprintf(w, `{"id":"pi_fixture","amount":%d,"amount_received":%d,"currency":"ars","status":"succeeded","metadata":{"order_id":"order-1","payment_attempt_id":"attempt-1"},"latest_charge":{"id":"ch_fixture","amount":%d,"amount_captured":%d,"amount_refunded":0,"currency":"ars","captured":true,"disputed":false}}`, expectedMinor, expectedMinor, expectedMinor, expectedMinor)
		default:
			t.Errorf("unexpected path %s", r.URL.Path)
			w.WriteHeader(500)
		}
	})
	scope := Scope{TenantID: "tenant", OrganizationID: "store", ConnectionID: "connection", ProviderCode: "stripe", AccountRef: account, Currency: "ARS", MinorUnitExponent: 2}
	driver, err := NewSDKDriver(SDKDriverConfig{Scope: scope, SuccessURL: "https://shop.example.test/return", CancelURL: "https://shop.example.test/cancel", DisplayName: "Order"}, "sk_test_fixture", client)
	if err != nil {
		t.Fatal(err)
	}
	req := Request{TenantID: scope.TenantID, OrganizationID: scope.OrganizationID, ProviderCode: scope.ProviderCode, Currency: scope.Currency, PaymentAttemptID: "attempt-1", OrderID: "order-1", CustomerSubject: "customer", AmountMinor: expectedMinor, CheckoutExpiresAt: expires}
	created, err := driver.CreateCheckout(context.Background(), req)
	if err != nil || created.AmountMinor != expectedMinor {
		t.Fatalf("checkout %+v err=%v", created, err)
	}
	inbox := &fixtureInbox{records: map[string]providerintegration.Receipt{}}
	webhook, err := NewWebhook(WebhookConfig{TenantID: scope.TenantID, ConnectionID: scope.ConnectionID, ProviderCode: "stripe", Tolerance: time.Minute, MaxConcurrent: 1, Secrets: fixtureSecret{}, Inbox: inbox})
	if err != nil {
		t.Fatal(err)
	}
	body := []byte(fmt.Sprintf(`{"id":"evt_exact","object":"event","api_version":%q,"type":"checkout.session.completed","livemode":false,"data":{"object":{"id":"cs_test_fixture","object":"checkout.session"}}}`, stripe.APIVersion))
	signed := stripewebhook.GenerateTestSignedPayload(&stripewebhook.UnsignedPayload{Payload: body, Secret: "whsec_fixture", Timestamp: time.Now()})
	callback := httptest.NewRequest("POST", "https://api.example.test/payment-webhook", strings.NewReader(string(body)))
	callback.Header.Set("Content-Type", "application/json")
	callback.Header.Set("Stripe-Signature", signed.Header)
	response := httptest.NewRecorder()
	webhook.ServeHTTP(response, callback)
	if response.Code != 200 {
		t.Fatalf("webhook status=%d body=%s", response.Code, response.Body.String())
	}
	var payload InboxPayload
	if json.Unmarshal(inbox.records["evt_exact"].Payload, &payload) != nil || payload.Notification.ResourceType != "checkout_session" {
		t.Fatal("durable webhook notice mismatch")
	}
	observedCheckout, err := driver.RetrieveCheckout(context.Background(), payload.Notification.ProviderReference)
	if err != nil || observedCheckout.AmountMinor != expectedMinor {
		t.Fatalf("checkout GET %+v err=%v", observedCheckout, err)
	}
	paid, err := driver.RetrievePayment(context.Background(), observedCheckout.PaymentReference)
	if err != nil || paid.ReceivedMinor != expectedMinor {
		t.Fatalf("payment GET %+v err=%v", paid, err)
	}
	t.Logf("PAYMENT_EXACT_BC_AMOUNT_WEBHOOK_FIXTURE_PASS amount_minor=%d quantity=%d unit_price_minor=%d bcamounts_oracle=true webhook_inbox=true production_claim=false", expectedMinor, quantity, unitPriceMinor)
}

func TestSDKDriverRejectsMismatchedStripeCredentialModeBeforeHTTP(t *testing.T) {
	config := SDKDriverConfig{Scope: Scope{TenantID: "tenant", OrganizationID: "store", ConnectionID: "connection", ProviderCode: "stripe", AccountRef: "acct_fixture", Currency: "ARS", MinorUnitExponent: 2}, SuccessURL: "https://shop.example.test/return", CancelURL: "https://shop.example.test/cancel", DisplayName: "Order"}
	client := fixtureHTTPClient(t, "api.stripe.com", func(w http.ResponseWriter, r *http.Request) {
		t.Error("constructor must not perform HTTP")
		w.WriteHeader(500)
	})
	for _, token := range []string{"sk_live_fixture", "rk_live_fixture", "pk_test_fixture", "sk_test_", "sk_test_fixture\n"} {
		if _, err := NewSDKDriver(config, token, client); err == nil {
			t.Error("invalid credential mode accepted")
		}
	}
	for _, token := range []string{"sk_test_fixture", "rk_test_fixture"} {
		if _, err := NewSDKDriver(config, token, client); err != nil {
			t.Fatal(err)
		}
	}
	config.Scope.LiveMode = true
	if _, err := NewSDKDriver(config, "sk_test_fixture", client); err == nil {
		t.Error("test credential accepted for live profile")
	}
	if _, err := NewSDKDriver(config, "rk_live_fixture", client); err != nil {
		t.Fatal(err)
	}
}
