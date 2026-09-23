package leadstream

import (
	"bytes"
	"context"
	"errors"
	"testing"
)

const validMercadoLibreQuestion = `{"id":11751825075,"seller_id":123456789,"buyer_id":56801932,"item_id":"MLA739200576","status":"UNANSWERED","text":"Necesito información","date_created":"2026-09-05T12:00:00Z","future_field":{"preserved":true}}`

func TestMercadoLibreQuestionNormalizesIntoCommonPendingPolicyCandidate(t *testing.T) {
	d, err := DecodeMercadoLibreQuestion([]byte(validMercadoLibreQuestion), "123456789", "tenant-1", "store-1", fixedTime())
	if err != nil {
		t.Fatal(err)
	}
	if d.Candidate == nil || d.Candidate.Provider != "mercadolibre" || d.Candidate.ProviderLeadID != "11751825075" || d.Candidate.ContactEligibility != "pending_policy" {
		t.Fatalf("unexpected candidate: %+v", d.Candidate)
	}
	if d.Raw.SourceSHA256 != payloadHash([]byte(validMercadoLibreQuestion)) || !bytes.Contains(d.Raw.Payload, []byte("future_field")) {
		t.Fatal("provider evidence identity or forward-compatible raw body was lost")
	}
	store := NewMemoryStore()
	receipt, err := store.Record(context.Background(), d.Raw, d.Candidate, "")
	if err != nil || receipt.State != ReceiptNormalized || store.Count() != 1 {
		t.Fatalf("receipt=%+v err=%v count=%d", receipt, err, store.Count())
	}
}

func TestMercadoLibreVehicleContactDataNeverGrantsContactEligibility(t *testing.T) {
	payload := `{"id":11949565740,"seller_id":123456789,"text":"Quiero coordinar","status":"UNANSWERED","item_id":"MLA595976788","date_created":"2026-09-05T12:07:18.109-04:00","from":{"id":21547449,"first_name":"Juan","last_name":"Lead","phone":{"number":"95712582","area_code":"9"},"email":"juan@example.com"}}`
	d, err := DecodeMercadoLibreQuestion([]byte(payload), "123456789", "tenant-1", "store-1", fixedTime())
	if err != nil || d.Candidate == nil {
		t.Fatalf("decode=%+v err=%v", d, err)
	}
	if d.Candidate.ContactEligibility != "pending_policy" || len(d.Candidate.Fields) != 8 {
		t.Fatalf("PII bypassed policy or was not normalized: %+v", d.Candidate)
	}
}

func TestMercadoLibreQuestionFailClosedBoundaries(t *testing.T) {
	if _, err := DecodeMercadoLibreQuestion([]byte(validMercadoLibreQuestion), "999999999", "tenant-1", "store-1", fixedTime()); !errors.Is(err, ErrInvalidCredential) {
		t.Fatalf("seller mismatch must fail before storage: %v", err)
	}
	banned := `{"id":11751825075,"seller_id":123456789,"buyer_id":56801932,"item_id":"MLA739200576","status":"BANNED","text":"","date_created":"2026-09-05T12:00:00Z"}`
	d, err := DecodeMercadoLibreQuestion([]byte(banned), "123456789", "tenant-1", "store-1", fixedTime())
	if err != nil || d.Candidate != nil || d.NormalizationCode != "QUESTION_BANNED" {
		t.Fatalf("banned question must be retained as rejection: %+v %v", d, err)
	}
	ambiguous := `{"id":11751825075,"seller_id":123456789,"buyer_id":56801932,"item_id":"MLA739200576","status":"UNANSWERED","text":"x","date_created":"2026-09-05T12:00:00Z","from":{"id":56801933}}`
	d, err = DecodeMercadoLibreQuestion([]byte(ambiguous), "123456789", "tenant-1", "store-1", fixedTime())
	if err != nil || d.Candidate != nil || d.NormalizationCode != "AMBIGUOUS_BUYER_ID" {
		t.Fatalf("ambiguous buyer must not normalize: %+v %v", d, err)
	}
}
