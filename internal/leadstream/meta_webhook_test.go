package leadstream

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

const validMetaLeadWebhook = `{"object":"page","entry":[{"id":"138124712925158","time":1630087927,"changes":[{"value":{"form_id":"172036835026334","leadgen_id":"233370042045345","created_time":1630087926,"page_id":"138124712925158"},"field":"leadgen"}]}]}`

func metaSignature(payload, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(payload))
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}

func metaHandler(store MetaWebhookStore) MetaLeadWebhookHandler {
	return MetaLeadWebhookHandler{Store: store, TenantID: "tenant-1", OrganizationID: "store-1", AppSecret: "app-secret", VerifyToken: "verify-token", Clock: fixedTime}
}

func TestDecodeMetaLeadWebhookMatchesOfficialLeadgenEnvelope(t *testing.T) {
	batch, err := DecodeMetaLeadWebhook([]byte(validMetaLeadWebhook), metaSignature(validMetaLeadWebhook, "app-secret"), "app-secret", "tenant-1", "store-1", fixedTime())
	if err != nil {
		t.Fatal(err)
	}
	if batch.PayloadSHA256 != payloadHash([]byte(validMetaLeadWebhook)) || len(batch.Signals) != 1 {
		t.Fatalf("batch=%+v", batch)
	}
	signal := batch.Signals[0]
	if signal.LeadgenID != "233370042045345" || signal.FormID != "172036835026334" || signal.PageID != "138124712925158" || !signal.OccurredAt.Equal(time.Unix(1630087926, 0).UTC()) {
		t.Fatalf("signal=%+v", signal)
	}
}

func TestMetaSignatureUsesExactRawBodyAndSHA256(t *testing.T) {
	valid := metaSignature(validMetaLeadWebhook, "app-secret")
	if _, err := DecodeMetaLeadWebhook([]byte(validMetaLeadWebhook), valid, "app-secret", "tenant-1", "store-1", fixedTime()); err != nil {
		t.Fatal(err)
	}
	for _, signature := range []string{"", "sha1=bad", "sha256=bad", metaSignature(validMetaLeadWebhook+" ", "app-secret")} {
		if _, err := DecodeMetaLeadWebhook([]byte(validMetaLeadWebhook), signature, "app-secret", "tenant-1", "store-1", fixedTime()); !errors.Is(err, ErrInvalidMetaSignature) {
			t.Fatalf("signature %q: %v", signature, err)
		}
	}
}

func TestMetaWebhookSubscriptionVerification(t *testing.T) {
	handler := metaHandler(NewMemoryMetaWebhookStore())
	request := httptest.NewRequest(http.MethodGet, "/webhooks/meta-leads?hub.mode=subscribe&hub.challenge=challenge-123&hub.verify_token=verify-token", nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusOK || response.Body.String() != "challenge-123" {
		t.Fatalf("status=%d body=%q", response.Code, response.Body.String())
	}
	request = httptest.NewRequest(http.MethodGet, "/webhooks/meta-leads?hub.mode=subscribe&hub.challenge=challenge-123&hub.verify_token=wrong", nil)
	response = httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusForbidden || strings.Contains(response.Body.String(), "wrong") {
		t.Fatalf("status=%d body=%q", response.Code, response.Body.String())
	}
}

func TestMetaWebhookAcknowledgesOnlyAfterDurableRecord(t *testing.T) {
	store := NewMemoryMetaWebhookStore()
	handler := metaHandler(store)
	call := func(payload, signature string) *httptest.ResponseRecorder {
		request := httptest.NewRequest(http.MethodPost, "/webhooks/meta-leads", strings.NewReader(payload))
		request.Header.Set("content-type", "application/json")
		request.Header.Set("X-Hub-Signature-256", signature)
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)
		return response
	}
	if response := call(validMetaLeadWebhook, metaSignature(validMetaLeadWebhook, "app-secret")); response.Code != http.StatusOK || response.Body.String() != "EVENT_RECEIVED" || store.SignalCount() != 1 {
		t.Fatalf("status=%d body=%q count=%d", response.Code, response.Body.String(), store.SignalCount())
	}
	if response := call(validMetaLeadWebhook, metaSignature(validMetaLeadWebhook, "app-secret")); response.Code != http.StatusOK || store.SignalCount() != 1 {
		t.Fatalf("replay status=%d count=%d", response.Code, store.SignalCount())
	}
	store.Fail = errors.New("database unavailable")
	second := strings.Replace(validMetaLeadWebhook, "233370042045345", "233370042045346", 1)
	if response := call(second, metaSignature(second, "app-secret")); response.Code != http.StatusInternalServerError || store.SignalCount() != 1 {
		t.Fatalf("failure status=%d count=%d", response.Code, store.SignalCount())
	}
}

func TestMetaWebhookRejectsInvalidBoundaryWithoutStorage(t *testing.T) {
	for _, test := range []struct {
		name, payload, signature, contentType string
		maxBytes                              int64
		want                                  int
	}{
		{name: "signature", payload: validMetaLeadWebhook, signature: "sha256=bad", contentType: "application/json", want: http.StatusForbidden},
		{name: "content type", payload: validMetaLeadWebhook, signature: metaSignature(validMetaLeadWebhook, "app-secret"), contentType: "text/plain", want: http.StatusUnsupportedMediaType},
		{name: "malformed", payload: `{`, signature: metaSignature(`{`, "app-secret"), contentType: "application/json", want: http.StatusBadRequest},
		{name: "oversize", payload: validMetaLeadWebhook, signature: metaSignature(validMetaLeadWebhook, "app-secret"), contentType: "application/json", maxBytes: 8, want: http.StatusRequestEntityTooLarge},
		{name: "wrong object", payload: `{"object":"user","entry":[]}`, signature: metaSignature(`{"object":"user","entry":[]}`, "app-secret"), contentType: "application/json", want: http.StatusBadRequest},
	} {
		t.Run(test.name, func(t *testing.T) {
			store := NewMemoryMetaWebhookStore()
			handler := metaHandler(store)
			handler.MaxBytes = test.maxBytes
			request := httptest.NewRequest(http.MethodPost, "/webhooks/meta-leads", strings.NewReader(test.payload))
			request.Header.Set("content-type", test.contentType)
			request.Header.Set("X-Hub-Signature-256", test.signature)
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, request)
			if response.Code != test.want || store.SignalCount() != 0 {
				t.Fatalf("status=%d want=%d stored=%d", response.Code, test.want, store.SignalCount())
			}
		})
	}
}

func TestMetaWebhookBatchedEntriesAndDivergence(t *testing.T) {
	second := `{"object":"page","entry":[{"id":"138124712925158","time":1630087927,"changes":[{"value":{"form_id":"172036835026334","leadgen_id":"233370042045345","created_time":1630087926,"page_id":"138124712925158"},"field":"leadgen"}]},{"id":"138124712925158","time":1630087928,"changes":[{"value":{"form_id":"172036835026334","leadgen_id":"233370042045346","created_time":1630087927,"page_id":"138124712925158"},"field":"leadgen"}]}]}`
	batch, err := DecodeMetaLeadWebhook([]byte(second), metaSignature(second, "app-secret"), "app-secret", "tenant-1", "store-1", fixedTime())
	if err != nil || len(batch.Signals) != 2 {
		t.Fatalf("signals=%d err=%v payload=%s", len(batch.Signals), err, second)
	}
	store := NewMemoryMetaWebhookStore()
	if receipt, err := store.RecordMetaWebhook(t.Context(), batch); err != nil || receipt.NewSignals != 2 {
		t.Fatalf("receipt=%+v err=%v", receipt, err)
	}
	changed := strings.Replace(validMetaLeadWebhook, "172036835026334", "different-form", 1)
	changedBatch, err := DecodeMetaLeadWebhook([]byte(changed), metaSignature(changed, "app-secret"), "app-secret", "tenant-1", "store-1", fixedTime())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.RecordMetaWebhook(t.Context(), changedBatch); !errors.Is(err, ErrDivergentDuplicate) {
		t.Fatalf("expected divergence, got %v", err)
	}
}

func TestMemoryMetaWebhookConcurrentReplay(t *testing.T) {
	batch, err := DecodeMetaLeadWebhook([]byte(validMetaLeadWebhook), metaSignature(validMetaLeadWebhook, "app-secret"), "app-secret", "tenant-1", "store-1", fixedTime())
	if err != nil {
		t.Fatal(err)
	}
	store := NewMemoryMetaWebhookStore()
	const workers = 24
	var wg sync.WaitGroup
	errs := make(chan error, workers)
	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, recordErr := store.RecordMetaWebhook(t.Context(), batch)
			errs <- recordErr
		}()
	}
	wg.Wait()
	close(errs)
	for recordErr := range errs {
		if recordErr != nil {
			t.Fatal(recordErr)
		}
	}
	if store.SignalCount() != 1 {
		t.Fatalf("signals=%d", store.SignalCount())
	}
}
