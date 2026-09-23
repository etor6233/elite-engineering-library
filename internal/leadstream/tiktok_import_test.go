package leadstream

import (
	"context"
	"encoding/json"
	"testing"
	"time"
)

func tiktokArtifacts(t *testing.T, mode string, leadData map[string]any, outcomeStatus string) ([]byte, []byte, []byte) {
	t.Helper()
	providerObject := map[string]any{
		"schema": "elite-tiktok-lead-sdk-response/v1", "sdk_version": tiktokSDKVersion, "api_version": tiktokAPIVersion,
		"request_identity": map[string]any{"mode": mode, "lead_source": "INSTANT_FORM", "advertiser_id_sha256": tiktokSHA("123"), "library_id_sha256": "", "page_id_sha256": tiktokSHA("456")},
		"data":             map[string]any{"lead_data": leadData, "meta_data": map[string]any{"lead_source": "INSTANT_FORM", "lead_id": "lead-1", "page_id": "456", "campaign_id": "11", "adgroup_id": "22", "ad_id": "33", "create_time": "2026-09-04 12:00:00"}}, "request_id": "req-1",
	}
	provider, err := json.Marshal(providerObject)
	if err != nil {
		t.Fatal(err)
	}
	provider = append(provider, '\n')
	var candidate any
	codes := []string{}
	if outcomeStatus == "CANDIDATE" {
		fields := make([]map[string]string, 0, len(leadData))
		for key, value := range leadData {
			fields = append(fields, map[string]string{"id": key, "value": value.(string)})
		}
		candidate = map[string]any{"tenant_id": "44444444-4444-4444-8444-444444444444", "organization_id": "store-1", "provider": tiktokProvider, "provider_lead_id": "lead-1", "form_id": "456", "campaign_id": "11", "ad_group_id": "22", "creative_id": "33", "source_kind": "INSTANT_FORM", "submitted_at": "2026-09-04T12:00:00Z", "is_test": mode == "TEST", "fields": fields, "contact_eligibility": "pending_policy"}
	} else {
		codes = []string{"NON_STRING_DYNAMIC_FIELD_REQUIRES_MAPPING"}
	}
	batchObject := map[string]any{"schema": "elite-tiktok-lead-candidate-batch/v1", "provider": tiktokProvider, "sdk_version": tiktokSDKVersion, "sdk_wheel_sha256": tiktokSDKWheelSHA256, "api_version": tiktokAPIVersion, "contract_observed_at": tiktokContractObserved, "tenant_id": "44444444-4444-4444-8444-444444444444", "organization_id": "store-1", "provider_response_sha256": shaHex(provider), "outcomes": []any{map[string]any{"status": outcomeStatus, "normalization_codes": codes, "provider_lead_id_sha256": tiktokSHA("lead-1"), "candidate": candidate}}, "automatic_business_write": false, "automatic_contact_eligibility": false}
	batch, err := json.Marshal(batchObject)
	if err != nil {
		t.Fatal(err)
	}
	batch = append(batch, '\n')
	candidates := 0
	if outcomeStatus == "CANDIDATE" {
		candidates = 1
	}
	receiptObject := map[string]any{"schema": "elite-tiktok-lead-retrieval-receipt/v1", "sdk_version": tiktokSDKVersion, "sdk_wheel_sha256": tiktokSDKWheelSHA256, "api_version": tiktokAPIVersion, "contract_observed_at": tiktokContractObserved, "retrieval_mode": mode, "request_id_sha256": tiktokSHA("req-1"), "provider_response_sha256": shaHex(provider), "candidate_batch_sha256": shaHex(batch), "row_count": 1, "candidate_count": candidates, "rejected_count": 1 - candidates, "automatic_business_write": false, "automatic_contact_eligibility": false}
	receipt, err := json.Marshal(receiptObject)
	if err != nil {
		t.Fatal(err)
	}
	return provider, batch, append(receipt, '\n')
}

func TestImportTikTokEvidenceNormalizesAndReplays(t *testing.T) {
	provider, batch, receipt := tiktokArtifacts(t, "TEST", map[string]any{"Email": "person@example.test"}, "CANDIDATE")
	store := NewMemoryStore()
	first, err := ImportTikTokEvidence(context.Background(), store, provider, batch, receipt, time.Now())
	if err != nil || !first.Complete || first.State != ReceiptNormalized || store.Count() != 1 {
		t.Fatalf("first=%+v err=%v count=%d", first, err, store.Count())
	}
	replay, err := ImportTikTokEvidence(context.Background(), store, provider, batch, receipt, time.Now())
	if err != nil || replay.State != ReceiptDuplicate || store.Count() != 1 {
		t.Fatalf("replay=%+v err=%v count=%d", replay, err, store.Count())
	}
}

func TestImportTikTokEvidenceRejectsTamperingBeforeWrite(t *testing.T) {
	provider, batch, receipt := tiktokArtifacts(t, "LIVE", map[string]any{"Email": "person@example.test"}, "CANDIDATE")
	batch[len(batch)-2] = ' '
	store := NewMemoryStore()
	if _, err := ImportTikTokEvidence(context.Background(), store, provider, batch, receipt, time.Now()); err == nil {
		t.Fatal("tampered batch accepted")
	}
	if store.Count() != 0 {
		t.Fatal("write occurred before complete validation")
	}
}

func TestImportTikTokEvidenceRejectsCandidateAndModeMismatch(t *testing.T) {
	provider, batch, receipt := tiktokArtifacts(t, "LIVE", map[string]any{"Email": "person@example.test"}, "CANDIDATE")
	var object map[string]any
	if err := json.Unmarshal(batch, &object); err != nil {
		t.Fatal(err)
	}
	outcome := object["outcomes"].([]any)[0].(map[string]any)
	outcome["candidate"].(map[string]any)["campaign_id"] = "different"
	batch, _ = json.Marshal(object)
	batch = append(batch, '\n')
	var receiptObject map[string]any
	if err := json.Unmarshal(receipt, &receiptObject); err != nil {
		t.Fatal(err)
	}
	receiptObject["candidate_batch_sha256"] = shaHex(batch)
	receipt, _ = json.Marshal(receiptObject)
	receipt = append(receipt, '\n')
	store := NewMemoryStore()
	if _, err := ImportTikTokEvidence(context.Background(), store, provider, batch, receipt, time.Now()); err == nil {
		t.Fatal("mismatched candidate accepted")
	}
	if store.Count() != 0 {
		t.Fatal("mismatched candidate reached store")
	}

	provider, batch, receipt = tiktokArtifacts(t, "LIVE", map[string]any{"Email": "person@example.test"}, "CANDIDATE")
	if err := json.Unmarshal(receipt, &receiptObject); err != nil {
		t.Fatal(err)
	}
	receiptObject["retrieval_mode"] = "TEST"
	receipt, _ = json.Marshal(receiptObject)
	receipt = append(receipt, '\n')
	if _, err := ImportTikTokEvidence(context.Background(), store, provider, batch, receipt, time.Now()); err == nil {
		t.Fatal("mode mismatch accepted")
	}
}

func TestImportTikTokEvidencePersistsRejectedEvidence(t *testing.T) {
	provider, batch, receipt := tiktokArtifacts(t, "LIVE", map[string]any{"multi": []string{"a", "b"}}, "REJECTED")
	store := NewMemoryStore()
	result, err := ImportTikTokEvidence(context.Background(), store, provider, batch, receipt, time.Now())
	if err != nil || result.State != ReceiptRejected || !result.Complete || store.Count() != 1 {
		t.Fatalf("result=%+v err=%v count=%d", result, err, store.Count())
	}
}

func TestImportTikTokEvidenceRejectsAccountPageAndReceiptIdentity(t *testing.T) {
	provider, batch, receipt := tiktokArtifacts(t, "LIVE", map[string]any{"Email": "person@example.test"}, "CANDIDATE")
	var object map[string]any
	if err := json.Unmarshal(provider, &object); err != nil {
		t.Fatal(err)
	}
	request := object["request_identity"].(map[string]any)
	request["library_id_sha256"] = tiktokSHA("also-set")
	provider, _ = json.Marshal(object)
	provider = append(provider, '\n')
	store := NewMemoryStore()
	if _, err := ImportTikTokEvidence(context.Background(), store, provider, batch, receipt, time.Now()); err == nil {
		t.Fatal("two account selectors accepted")
	}
	if store.Count() != 0 {
		t.Fatal("invalid identity reached store")
	}
}

func TestImportTikTokEvidenceRejectsMalformedSelectorAndOptionalMetaBeforeWrite(t *testing.T) {
	for name, mutate := range map[string]func(map[string]any){
		"malformed account selector": func(object map[string]any) {
			request := object["request_identity"].(map[string]any)
			request["advertiser_id_sha256"] = "not-a-sha256"
		},
		"non-string optional campaign": func(object map[string]any) {
			data := object["data"].(map[string]any)
			meta := data["meta_data"].(map[string]any)
			meta["campaign_id"] = map[string]any{"unexpected": true}
		},
	} {
		t.Run(name, func(t *testing.T) {
			provider, batch, receipt := tiktokArtifacts(t, "LIVE", map[string]any{"Email": "person@example.test"}, "CANDIDATE")
			var object map[string]any
			if err := json.Unmarshal(provider, &object); err != nil {
				t.Fatal(err)
			}
			mutate(object)
			provider, _ = json.Marshal(object)
			provider = append(provider, '\n')
			store := NewMemoryStore()
			if _, err := ImportTikTokEvidence(context.Background(), store, provider, batch, receipt, time.Now()); err == nil {
				t.Fatal("malformed provider evidence accepted")
			}
			if store.Count() != 0 {
				t.Fatal("malformed provider evidence reached store")
			}
		})
	}
}
