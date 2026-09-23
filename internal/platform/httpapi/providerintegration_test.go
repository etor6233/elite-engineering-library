package httpapi

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"elite.local/enterprise/internal/providerintegration"
)

type webhookRepository struct{ calls int }

func (r *webhookRepository) AcceptWebhook(_ context.Context, _ providerintegration.Receipt) (bool, error) {
	r.calls++
	return r.calls > 1, nil
}

func TestProviderWebhookHTTPAdmissionAndReplay(t *testing.T) {
	now := time.Unix(1_800_000_000, 0)
	secret := []byte("01234567890123456789012345678901")
	registry, err := providerintegration.LoadHMACConnections([]byte(`[{"tenant_id":"tenant","connection_id":"primary","provider_code":"sandbox","secret_env":"SANDBOX_SECRET"}]`), func(string) (string, bool) { return string(secret), true })
	if err != nil {
		t.Fatal(err)
	}
	repository := &webhookRepository{}
	service := providerintegration.NewService(repository, registry, providerintegration.HMACSHA256{Now: func() time.Time { return now }, Tolerance: 5 * time.Minute})
	mux := http.NewServeMux()
	ProviderIntegrationModule{Service: service}.Register(mux, nil)
	body := `{"event_id":"evt-1","type":"payment.updated","data":{}}`
	for index, expected := range []int{202, 200} {
		request := httptest.NewRequest("POST", "/v1/integrations/sandbox/connections/primary/webhooks", strings.NewReader(body))
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("X-Elite-Webhook-Timestamp", "1800000000")
		request.Header.Set("X-Elite-Webhook-Signature", providerSignature(secret, "1800000000", []byte(body)))
		response := httptest.NewRecorder()
		mux.ServeHTTP(response, request)
		if response.Code != expected {
			t.Fatalf("index=%d status=%d body=%s", index, response.Code, response.Body.String())
		}
	}
	request := httptest.NewRequest("POST", "/v1/integrations/sandbox/connections/primary/webhooks", strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Elite-Webhook-Timestamp", "1800000000")
	request.Header.Set("X-Elite-Webhook-Signature", "v1=00")
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != 401 {
		t.Fatalf("tampered status=%d", response.Code)
	}
}

func providerSignature(secret []byte, timestamp string, body []byte) string {
	mac := hmac.New(sha256.New, secret)
	_, _ = mac.Write([]byte(timestamp + "."))
	_, _ = mac.Write(body)
	return "v1=" + hex.EncodeToString(mac.Sum(nil))
}
