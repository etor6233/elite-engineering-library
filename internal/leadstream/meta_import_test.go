package leadstream

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func metaArtifacts(t *testing.T, rows []map[string]any, outcomes []map[string]any, mode string, truncated bool) ([]byte, []byte, []byte) {
	t.Helper()
	provider, err := json.MarshalIndent(map[string]any{"graph_api_version": metaAPIVersion, "retrieval_mode": mode, "rows": rows, "truncated_by_local_bound": truncated}, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	provider = append(provider, '\n')
	providerHash := shaHex(provider)
	batch, err := json.MarshalIndent(map[string]any{
		"schema": "elite-meta-lead-candidate-batch/v1", "provider": "meta_lead_ads", "source_commit": metaSourceCommit,
		"tenant_id": "44444444-4444-4444-8444-444444444444", "organization_id": "store-1",
		"provider_response_sha256": providerHash, "outcomes": outcomes,
		"automatic_business_write": false, "automatic_contact_eligibility": false,
	}, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	batch = append(batch, '\n')
	formID := rows[0]["form_id"].(string)
	candidates := 0
	for _, outcome := range outcomes {
		if outcome["status"] == "CANDIDATE" {
			candidates++
		}
	}
	receipt, err := json.MarshalIndent(map[string]any{
		"schema": "elite-meta-lead-reconciliation-receipt/v1", "sdk_version": metaSDKVersion,
		"graph_api_version": metaAPIVersion, "source_commit": metaSourceCommit,
		"form_id_sha256": shaHex([]byte(formID)), "provider_response_sha256": providerHash,
		"candidate_batch_sha256": shaHex(batch), "row_count": len(rows), "candidate_count": candidates,
		"rejected_count": len(rows) - candidates, "truncated_by_local_bound": truncated,
		"automatic_business_write": false, "automatic_contact_eligibility": false,
	}, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	return provider, batch, append(receipt, '\n')
}

func metaCandidateOutcome(id string, isTest bool) map[string]any {
	h := sha256.Sum256([]byte(id))
	return map[string]any{
		"status": "CANDIDATE", "normalization_codes": []string{}, "provider_lead_id_sha256": hex.EncodeToString(h[:]),
		"candidate": map[string]any{
			"tenant_id": "44444444-4444-4444-8444-444444444444", "organization_id": "store-1",
			"provider": "meta_lead_ads", "provider_lead_id": id, "form_id": "123456789",
			"campaign_id": "campaign-1", "ad_group_id": "adset-1", "creative_id": "ad-1",
			"source_kind": "facebook", "submitted_at": "2026-09-04T12:00:00+0000", "is_test": isTest,
			"fields":              []map[string]string{{"id": "email", "value": "person@example.test"}},
			"contact_eligibility": "pending_policy",
		},
	}
}

func metaRow(id string) map[string]any {
	return map[string]any{
		"id": id, "created_time": "2026-09-04T12:00:00+0000", "form_id": "123456789",
		"campaign_id": "campaign-1", "adset_id": "adset-1", "ad_id": "ad-1", "platform": "facebook",
		"field_data": []map[string]any{{"name": "email", "values": []string{"person@example.test"}}},
	}
}

func TestImportMetaEvidenceNormalizesAndReplays(t *testing.T) {
	provider, batch, receipt := metaArtifacts(t, []map[string]any{metaRow("lead-1")}, []map[string]any{metaCandidateOutcome("lead-1", false)}, "LIVE", false)
	store := NewMemoryStore()
	first, err := ImportMetaEvidence(context.Background(), store, provider, batch, receipt, time.Now())
	if err != nil || !first.Complete || first.Normalized != 1 || store.Count() != 1 {
		t.Fatalf("first=%+v err=%v count=%d", first, err, store.Count())
	}
	replay, err := ImportMetaEvidence(context.Background(), store, provider, batch, receipt, time.Now())
	if err != nil || replay.Duplicates != 1 || store.Count() != 1 {
		t.Fatalf("replay=%+v err=%v count=%d", replay, err, store.Count())
	}
}

func TestImportMetaEvidenceRejectsTamperingBeforeWrite(t *testing.T) {
	provider, batch, receipt := metaArtifacts(t, []map[string]any{metaRow("lead-1")}, []map[string]any{metaCandidateOutcome("lead-1", false)}, "LIVE", false)
	batch[len(batch)-2] = ' '
	store := NewMemoryStore()
	if _, err := ImportMetaEvidence(context.Background(), store, provider, batch, receipt, time.Now()); err == nil {
		t.Fatal("tampered batch accepted")
	}
	if store.Count() != 0 {
		t.Fatal("write occurred before complete validation")
	}
}

func TestImportMetaEvidenceRejectsCandidateMismatchAndModeMismatch(t *testing.T) {
	row := metaRow("lead-1")
	outcome := metaCandidateOutcome("lead-1", false)
	outcome["candidate"].(map[string]any)["campaign_id"] = "different"
	provider, batch, receipt := metaArtifacts(t, []map[string]any{row}, []map[string]any{outcome}, "LIVE", false)
	store := NewMemoryStore()
	if _, err := ImportMetaEvidence(context.Background(), store, provider, batch, receipt, time.Now()); err == nil {
		t.Fatal("mismatched candidate accepted")
	}
	if store.Count() != 0 {
		t.Fatal("mismatched candidate reached store")
	}

	outcome = metaCandidateOutcome("lead-1", true)
	provider, batch, receipt = metaArtifacts(t, []map[string]any{row}, []map[string]any{outcome}, "LIVE", false)
	if _, err := ImportMetaEvidence(context.Background(), store, provider, batch, receipt, time.Now()); err == nil {
		t.Fatal("test candidate accepted as live")
	}
}

func TestImportMetaEvidenceRejectsMultiValueField(t *testing.T) {
	row := metaRow("lead-1")
	row["field_data"] = []map[string]any{{"name": "email", "values": []string{"a@example.test", "b@example.test"}}}
	provider, batch, receipt := metaArtifacts(t, []map[string]any{row}, []map[string]any{metaCandidateOutcome("lead-1", false)}, "LIVE", false)
	store := NewMemoryStore()
	if _, err := ImportMetaEvidence(context.Background(), store, provider, batch, receipt, time.Now()); err == nil {
		t.Fatal("multi-value field accepted")
	}
	if store.Count() != 0 {
		t.Fatal("ambiguous field reached store")
	}
}

func TestImportMetaEvidencePersistsRejectedEvidenceWithBatchScope(t *testing.T) {
	row := metaRow("lead-rejected")
	h := sha256.Sum256([]byte("lead-rejected"))
	outcome := map[string]any{"status": "REJECTED", "normalization_codes": []string{"UNAPPROVED_FORM_FIELD"}, "candidate": nil, "provider_lead_id_sha256": hex.EncodeToString(h[:])}
	provider, batch, receipt := metaArtifacts(t, []map[string]any{row}, []map[string]any{outcome}, "LIVE", false)
	store := NewMemoryStore()
	result, err := ImportMetaEvidence(context.Background(), store, provider, batch, receipt, time.Now())
	if err != nil || result.Rejected != 1 || !result.Complete || store.Count() != 1 {
		t.Fatalf("result=%+v err=%v count=%d", result, err, store.Count())
	}
}

func TestImportMetaEvidencePreservesTruncationAndStoreFailure(t *testing.T) {
	provider, batch, receipt := metaArtifacts(t, []map[string]any{metaRow("lead-1")}, []map[string]any{metaCandidateOutcome("lead-1", false)}, "LIVE", true)
	store := NewMemoryStore()
	store.Fail = errors.New("database unavailable")
	result, err := ImportMetaEvidence(context.Background(), store, provider, batch, receipt, time.Now())
	if err == nil || result.Processed != 0 || result.Complete || !result.TruncatedSource {
		t.Fatalf("result=%+v err=%v", result, err)
	}
}

func TestParseMetaTimeAcceptsOfficialOffsetAndRFC3339(t *testing.T) {
	for _, value := range []string{"2026-09-04T12:00:00+0000", "2026-09-04T12:00:00Z"} {
		if _, err := parseMetaTime(value); err != nil {
			t.Fatalf("%s: %v", value, err)
		}
	}
	if _, err := parseMetaTime("2026/09/04"); err == nil {
		t.Fatal("invalid time accepted")
	}
}

func TestImportMetaEvidenceConsumesExactV226PythonArtifacts(t *testing.T) {
	directory := filepath.Join("testdata", "meta-v226-valid")
	provider, err := os.ReadFile(filepath.Join(directory, "provider-response.json"))
	if err != nil {
		t.Fatal(err)
	}
	batch, err := os.ReadFile(filepath.Join(directory, "lead-candidates.json"))
	if err != nil {
		t.Fatal(err)
	}
	receipt, err := os.ReadFile(filepath.Join(directory, "RETRIEVAL_RECEIPT.json"))
	if err != nil {
		t.Fatal(err)
	}
	store := NewMemoryStore()
	result, err := ImportMetaEvidence(context.Background(), store, provider, batch, receipt, time.Now())
	if err != nil || !result.Complete || result.Normalized != 1 || store.Count() != 1 {
		t.Fatalf("result=%+v err=%v count=%d", result, err, store.Count())
	}
}
