package paymentbridge

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"elite.local/enterprise/internal/providerintegration"
	"github.com/stripe/stripe-go/v86"
	stripewebhook "github.com/stripe/stripe-go/v86/webhook"
)

type fixtureSecret struct{}

func (fixtureSecret) PaymentWebhookSecret(context.Context, string, string) (string, error) {
	return "whsec_fixture", nil
}

type fixtureInbox struct {
	mu      sync.Mutex
	records map[string]providerintegration.Receipt
	err     error
}

func (f *fixtureInbox) AcceptWebhook(_ context.Context, r providerintegration.Receipt) (bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.err != nil {
		return false, f.err
	}
	old, ok := f.records[r.ProviderEventID]
	if ok {
		if old.BodyHash != r.BodyHash {
			return false, providerintegration.ErrConflict
		}
		return true, nil
	}
	f.records[r.ProviderEventID] = r
	return false, nil
}
func signedRequest(t *testing.T, body []byte) *http.Request {
	t.Helper()
	signed := stripewebhook.GenerateTestSignedPayload(&stripewebhook.UnsignedPayload{Payload: body, Secret: "whsec_fixture", Timestamp: time.Now()})
	r := httptest.NewRequest("POST", "https://fixture.test/webhooks/payment", bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Stripe-Signature", signed.Header)
	return r
}
func paymentEvent() []byte {
	return []byte(fmt.Sprintf(`{"id":"evt_fixture","object":"event","api_version":%q,"type":"payment_intent.succeeded","livemode":false,"data":{"object":{"id":"pi_fixture","object":"payment_intent","client_secret":"must-never-persist","payer":{"email":"private@example.test"}}}}`, stripe.APIVersion))
}
func TestWebhookOfficialSignatureToScopedInboxAndReplay(t *testing.T) {
	inbox := &fixtureInbox{records: map[string]providerintegration.Receipt{}}
	h, err := NewWebhook(WebhookConfig{TenantID: "tenant-fixture", ConnectionID: "connection-fixture", ProviderCode: "stripe", Tolerance: 5 * time.Minute, MaxConcurrent: 2, Secrets: fixtureSecret{}, Inbox: inbox})
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 2; i++ {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, signedRequest(t, paymentEvent()))
		if w.Code != 200 {
			t.Fatalf("code%d: %s", w.Code, w.Body.String())
		}
	}
	if len(inbox.records) != 1 {
		t.Fatal("duplicate inbox")
	}
	row := inbox.records["evt_fixture"]
	if row.TenantID != "tenant-fixture" || row.ConnectionID != "connection-fixture" || row.ProviderCode != "stripe" {
		t.Fatal("scope mismatch")
	}
	var normalized InboxPayload
	if json.Unmarshal(row.Payload, &normalized) != nil || normalized.Notification.ProviderReference != "pi_fixture" {
		t.Fatal("notification mismatch")
	}
	if strings.Contains(string(row.Payload), "secret") || strings.Contains(string(row.Payload), "private") || strings.Contains(string(row.Payload), "succeeded\",\"amount") {
		t.Fatal("raw data leaked")
	}
}
func TestWebhookRejectsBeforeInboxAndReportsCommitFailure(t *testing.T) {
	inbox := &fixtureInbox{records: map[string]providerintegration.Receipt{}}
	h, _ := NewWebhook(WebhookConfig{TenantID: "tenant", ConnectionID: "connection", ProviderCode: "stripe", Tolerance: time.Minute, MaxConcurrent: 1, Secrets: fixtureSecret{}, Inbox: inbox})
	cases := []struct {
		name string
		edit func(*http.Request)
		want int
	}{{"missing-signature", func(r *http.Request) { r.Header.Del("Stripe-Signature") }, 400}, {"duplicate-header", func(r *http.Request) { r.Header.Add("Stripe-Signature", r.Header.Get("Stripe-Signature")) }, 400}, {"wrong-signature", func(r *http.Request) { r.Header.Set("Stripe-Signature", "t=1,v1=bad") }, 403}, {"query-scope", func(r *http.Request) { r.URL.RawQuery = "tenant=other" }, 400}, {"method", func(r *http.Request) { r.Method = "GET" }, 405}}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := signedRequest(t, paymentEvent())
			tc.edit(r)
			w := httptest.NewRecorder()
			h.ServeHTTP(w, r)
			if w.Code != tc.want {
				t.Fatalf("%d", w.Code)
			}
		})
	}
	if len(inbox.records) != 0 {
		t.Fatal("invalid receipt persisted")
	}
	inbox.err = errors.New("database unavailable with private info")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, signedRequest(t, paymentEvent()))
	if w.Code != 503 || strings.Contains(w.Body.String(), "private") {
		t.Fatal("failed commit accepted or leaked")
	}
}
