package leadstream

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

const validGooglePayload = `{"lead_id":"lead-123","user_column_data":[{"column_id":"FULL_NAME","string_value":"Ada Lovelace"},{"column_id":"EMAIL","string_value":"ada@example.com"}],"api_version":"3","form_id":9223372036854775806,"campaign_id":42,"google_key":"secret-key","is_test":false,"gcl_id":"click-1","adgroup_id":7,"creative_id":8,"asset_group_id":9,"lead_stage":"SUBMITTED","lead_submit_time":"2026-09-04T12:30:00Z","lead_source":"LEAD_FORM","future_field":{"kept":"in raw, ignored by parser"}}`

func fixedTime() time.Time { return time.Date(2026, 9, 4, 12, 31, 0, 0, time.UTC) }

func TestGoogleDecodePreservesSourceIdentityRedactsCredentialAndIgnoresUnknownFields(t *testing.T) {
	d, err := DecodeGoogleWebhook([]byte(validGooglePayload), "secret-key", "tenant-1", "store-1", fixedTime())
	if err != nil {
		t.Fatal(err)
	}
	if d.Candidate == nil || d.Candidate.FormID != "9223372036854775806" || d.Candidate.ContactEligibility != "pending_policy" {
		t.Fatalf("unexpected candidate: %+v", d.Candidate)
	}
	if d.Raw.ProviderEventID != "lead-123" || d.Raw.SourceSHA256 != payloadHash([]byte(validGooglePayload)) {
		t.Fatal("source hash or provider identity changed")
	}
	if d.Raw.StoredSHA256 != payloadHash(d.Raw.Payload) || d.Raw.RedactionProfile != "google_ads:remove-google_key:v1" {
		t.Fatal("stored payload evidence is not self-consistent")
	}
	if bytes.Contains(d.Raw.Payload, []byte("secret-key")) || bytes.Contains(d.Raw.Payload, []byte("google_key")) {
		t.Fatal("provider credential survived durable redaction")
	}
	if !bytes.Contains(d.Raw.Payload, []byte("future_field")) {
		t.Fatal("unknown provider field was not retained in redacted evidence")
	}
}

func TestGoogleDecodeRejectsWrongCredentialBeforeStorage(t *testing.T) {
	if _, err := DecodeGoogleWebhook([]byte(validGooglePayload), "wrong", "tenant-1", "store-1", fixedTime()); !errors.Is(err, ErrInvalidCredential) {
		t.Fatalf("expected credential error, got %v", err)
	}
}

func TestMemoryStoreExactReplayAndDivergence(t *testing.T) {
	d, _ := DecodeGoogleWebhook([]byte(validGooglePayload), "secret-key", "tenant-1", "store-1", fixedTime())
	s := NewMemoryStore()
	first, err := s.Record(context.Background(), d.Raw, d.Candidate, "")
	if err != nil || first.State != ReceiptNormalized {
		t.Fatalf("first record: %+v %v", first, err)
	}
	replay, err := s.Record(context.Background(), d.Raw, d.Candidate, "")
	if err != nil || replay.State != ReceiptDuplicate || s.Count() != 1 {
		t.Fatalf("replay: %+v %v count=%d", replay, err, s.Count())
	}
	changedPayload := strings.Replace(validGooglePayload, "Ada Lovelace", "Grace Hopper", 1)
	other, decodeErr := DecodeGoogleWebhook([]byte(changedPayload), "secret-key", "tenant-1", "store-1", fixedTime())
	if decodeErr != nil {
		t.Fatal(decodeErr)
	}
	if _, err := s.Record(context.Background(), other.Raw, other.Candidate, ""); !errors.Is(err, ErrDivergentDuplicate) {
		t.Fatalf("expected divergent duplicate, got %v", err)
	}
}

func TestGoogleHandlerAcknowledgesOnlyAfterStore(t *testing.T) {
	s := NewMemoryStore()
	h := GoogleHandler{Store: s, TenantID: "tenant-1", OrganizationID: "store-1", GoogleKey: "secret-key", Clock: fixedTime}
	r := httptest.NewRequest(http.MethodPost, "/webhooks/google-ads", strings.NewReader(validGooglePayload))
	r.Header.Set("content-type", "application/json; charset=utf-8")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != http.StatusOK || strings.TrimSpace(w.Body.String()) != "{}" || s.Count() != 1 {
		t.Fatalf("status=%d body=%q count=%d", w.Code, w.Body.String(), s.Count())
	}
	s.Fail = errors.New("durability unavailable")
	r = httptest.NewRequest(http.MethodPost, "/webhooks/google-ads", strings.NewReader(strings.Replace(validGooglePayload, "lead-123", "lead-124", 1)))
	r.Header.Set("content-type", "application/json")
	w = httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("store failure must be retryable 5xx, got %d", w.Code)
	}
}

func TestGoogleHandlerExactReplayAndDivergentIdentity(t *testing.T) {
	s := NewMemoryStore()
	h := GoogleHandler{Store: s, TenantID: "tenant-1", OrganizationID: "store-1", GoogleKey: "secret-key", Clock: fixedTime}
	call := func(payload string) *httptest.ResponseRecorder {
		r := httptest.NewRequest(http.MethodPost, "/webhooks/google-ads", strings.NewReader(payload))
		r.Header.Set("content-type", "application/json")
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		return w
	}
	if w := call(validGooglePayload); w.Code != http.StatusOK {
		t.Fatalf("first delivery status=%d body=%q", w.Code, w.Body.String())
	}
	if w := call(validGooglePayload); w.Code != http.StatusOK || s.Count() != 1 {
		t.Fatalf("exact replay status=%d count=%d", w.Code, s.Count())
	}
	changed := strings.Replace(validGooglePayload, "Ada Lovelace", "Grace Hopper", 1)
	if w := call(changed); w.Code != http.StatusConflict || s.Count() != 1 {
		t.Fatalf("divergent identity status=%d count=%d", w.Code, s.Count())
	}
}

func TestGoogleHandlerRequestBoundary(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		contentType string
		payload     string
		maxBytes    int64
		want        int
	}{
		{name: "method", method: http.MethodGet, contentType: "application/json", payload: validGooglePayload, want: http.StatusBadRequest},
		{name: "content type", method: http.MethodPost, contentType: "text/plain", payload: validGooglePayload, want: http.StatusUnsupportedMediaType},
		{name: "too large", method: http.MethodPost, contentType: "application/json", payload: validGooglePayload, maxBytes: 8, want: http.StatusRequestEntityTooLarge},
		{name: "malformed", method: http.MethodPost, contentType: "application/json", payload: `{`, want: http.StatusBadRequest},
		{name: "trailing json", method: http.MethodPost, contentType: "application/json", payload: validGooglePayload + `{}`, want: http.StatusBadRequest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewMemoryStore()
			h := GoogleHandler{Store: s, TenantID: "tenant-1", OrganizationID: "store-1", GoogleKey: "secret-key", Clock: fixedTime, MaxBytes: tt.maxBytes}
			r := httptest.NewRequest(tt.method, "/webhooks/google-ads", strings.NewReader(tt.payload))
			r.Header.Set("content-type", tt.contentType)
			w := httptest.NewRecorder()
			h.ServeHTTP(w, r)
			if w.Code != tt.want || s.Count() != 0 {
				t.Fatalf("status=%d want=%d stored=%d", w.Code, tt.want, s.Count())
			}
		})
	}
}

func TestGoogleHandlerRejectsCredentialAndRetainsAuthenticatedNormalizationFailure(t *testing.T) {
	s := NewMemoryStore()
	h := GoogleHandler{Store: s, TenantID: "tenant-1", OrganizationID: "store-1", GoogleKey: "secret-key", Clock: fixedTime}
	badKey := strings.Replace(validGooglePayload, "secret-key", "wrong-key", 1)
	r := httptest.NewRequest(http.MethodPost, "/webhooks/google-ads", strings.NewReader(badKey))
	r.Header.Set("content-type", "application/json")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != http.StatusBadRequest || s.Count() != 0 {
		t.Fatalf("untrusted payload stored: status=%d count=%d", w.Code, s.Count())
	}
	badTime := strings.Replace(validGooglePayload, "2026-09-04T12:30:00Z", "not-a-time", 1)
	r = httptest.NewRequest(http.MethodPost, "/webhooks/google-ads", strings.NewReader(badTime))
	r.Header.Set("content-type", "application/json")
	w = httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != http.StatusBadRequest || s.Count() != 1 {
		t.Fatalf("authenticated rejection not retained: status=%d count=%d", w.Code, s.Count())
	}
}

func TestCandidateRejectsAmbiguousDuplicateField(t *testing.T) {
	p := strings.Replace(validGooglePayload, `{"column_id":"EMAIL","string_value":"ada@example.com"}`, `{"column_id":"FULL_NAME","string_value":"Other"}`, 1)
	d, err := DecodeGoogleWebhook([]byte(p), "secret-key", "tenant-1", "store-1", fixedTime())
	if err != nil || d.Candidate != nil || d.NormalizationCode == "" {
		t.Fatalf("ambiguous lead must become retained rejection: %+v %v", d, err)
	}
}

func TestMissingProviderLeadIDIsRetainedAsRejectedEvidence(t *testing.T) {
	p := strings.Replace(validGooglePayload, `"lead_id":"lead-123"`, `"lead_id":""`, 1)
	d, err := DecodeGoogleWebhook([]byte(p), "secret-key", "tenant-1", "store-1", fixedTime())
	if err != nil || d.Candidate != nil || d.NormalizationCode != "MISSING_LEAD_ID" || !strings.HasPrefix(d.Raw.ProviderEventID, "sha256:") {
		t.Fatalf("unexpected decode: %+v %v", d, err)
	}
	s := NewMemoryStore()
	h := GoogleHandler{Store: s, TenantID: "tenant-1", OrganizationID: "store-1", GoogleKey: "secret-key", Clock: fixedTime}
	r := httptest.NewRequest(http.MethodPost, "/webhooks/google-ads", bytes.NewReader([]byte(p)))
	r.Header.Set("content-type", "application/json")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != http.StatusBadRequest || s.Count() != 1 {
		t.Fatalf("status=%d stored=%d", w.Code, s.Count())
	}
}

func TestMemoryStoreConcurrentExactReplayHasOneDurableIdentity(t *testing.T) {
	d, err := DecodeGoogleWebhook([]byte(validGooglePayload), "secret-key", "tenant-1", "store-1", fixedTime())
	if err != nil {
		t.Fatal(err)
	}
	s := NewMemoryStore()
	const workers = 32
	states := make(chan ReceiptState, workers)
	errs := make(chan error, workers)
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			receipt, recordErr := s.Record(context.Background(), d.Raw, d.Candidate, "")
			states <- receipt.State
			errs <- recordErr
		}()
	}
	wg.Wait()
	close(states)
	close(errs)
	for recordErr := range errs {
		if recordErr != nil {
			t.Fatal(recordErr)
		}
	}
	normalized, duplicates := 0, 0
	for state := range states {
		switch state {
		case ReceiptNormalized:
			normalized++
		case ReceiptDuplicate:
			duplicates++
		default:
			t.Fatalf("unexpected state %q", state)
		}
	}
	if normalized != 1 || duplicates != workers-1 || s.Count() != 1 {
		t.Fatalf("normalized=%d duplicates=%d count=%d", normalized, duplicates, s.Count())
	}
}
