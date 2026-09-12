# Go TikTok Lead Durable Import

## 1. Metadata

```yaml
pack_id: "GO-TIKTOK-LEAD-DURABLE-IMPORT"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Verifica los tres artefactos hash-linked del adapter TikTok Lead 0.2.0 y reutiliza el único Store provider-neutral/PostgreSQL/outbox existente para persistir candidato o rechazo antes de permitir promoción de contacto."
stacks: ["Go 1.26.7", "PostgreSQL 18.6", "TikTok Lead artifacts v1"]
compatible_with: ["GO-OMNICHANNEL-LEAD-INGRESS", "PYTHON-TIKTOK-LEAD-ADAPTER 0.2.0"]
incompatible_with: ["persistencia directa del webhook", "CRM paralelo", "artefactos sin hash", "contact eligibility automática"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://business-api.tiktok.com/portal/docs/get-an-instant-form-lead-or-a-direct-message-lead/v1.3", "https://github.com/tiktok/tiktok-business-api-sdk"]
verified_at: "2026-09-04"
```

## 2. Applicability

Use sólo junto con `PYTHON-TIKTOK-LEAD-ADAPTER` 0.2.0 y el owner `GO-OMNICHANNEL-LEAD-INGRESS`. El comando ingiere exactamente provider response, candidate batch y receipt; valida identidades, hashes, modo, selectores, metadata, fields y cardinalidad antes de la primera escritura. No crea tablas, CRM ni outbox alternativos.

## 3. Architecture contract

El artefacto de provider autenticado gobierna el lead; un webhook TikTok no autenticado nunca entra directamente al store. El importador conserva la respuesta SDK como raw evidence, usa `(tenant, provider=tiktok_lead_generation, lead_id)` para idempotencia y deja todo candidato en `pending_policy`. Un duplicado byte-idéntico converge; identidad reutilizada con bytes distintos falla cerrada. Fields no string ya llegan como `REJECTED` desde el adapter y se conservan sin contacto. La prueba PostgreSQL usa las tablas/outbox de migración 0044.

## 4. Exact file manifest

```text
CREATE internal/leadstream/tiktok_import.go
CREATE internal/leadstream/tiktok_import_test.go
CREATE cmd/tiktok-lead-import/main.go
CREATE internal/platform/postgres/tiktok_import_integration_test.go
CREATE docs/tiktok-lead-import.md
```

## 5. Materialization blocks

### FILE: `internal/leadstream/tiktok_import.go`
```yaml
block_id: "GO-TIKTOK-LEAD-IMPORT:domain:v1"
operation: CREATE
provenance: AUTHORED
source: "local importer governed by TikTok official Lead v1.3 and exact adapter artifact contracts"
license: "LicenseRef-Workspace-Owner"
sha256: "374ac77ec797d5147782f08fbb137199af52f74a05a819de1a3b75c4ed8cbee5"
variables: []
secrets_allowed: false
```
````go
package leadstream

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"
)

const (
	tiktokSDKVersion       = "1.1.3"
	tiktokSDKWheelSHA256   = "663b4a1f2585c4b386144d968646f695befd733c464418bbaeeac41befff33b7"
	tiktokAPIVersion       = "v1.3"
	tiktokContractObserved = "2026-09-04"
	tiktokProvider         = "tiktok_lead_generation"
)

var tiktokNormalizationCode = regexp.MustCompile(`^[A-Z][A-Z0-9_]{0,63}$`)

type tiktokRequestIdentity struct {
	Mode               string `json:"mode"`
	LeadSource         string `json:"lead_source"`
	AdvertiserIDSHA256 string `json:"advertiser_id_sha256"`
	LibraryIDSHA256    string `json:"library_id_sha256"`
	PageIDSHA256       string `json:"page_id_sha256"`
}

type tiktokProviderResponse struct {
	Schema          string                `json:"schema"`
	SDKVersion      string                `json:"sdk_version"`
	APIVersion      string                `json:"api_version"`
	RequestIdentity tiktokRequestIdentity `json:"request_identity"`
	Data            struct {
		LeadData map[string]json.RawMessage `json:"lead_data"`
		MetaData map[string]json.RawMessage `json:"meta_data"`
	} `json:"data"`
	RequestID string `json:"request_id"`
}

type tiktokCandidateBatch struct {
	Schema                      string          `json:"schema"`
	Provider                    string          `json:"provider"`
	SDKVersion                  string          `json:"sdk_version"`
	SDKWheelSHA256              string          `json:"sdk_wheel_sha256"`
	APIVersion                  string          `json:"api_version"`
	ContractObservedAt          string          `json:"contract_observed_at"`
	TenantID                    string          `json:"tenant_id"`
	OrganizationID              string          `json:"organization_id"`
	ProviderResponseSHA256      string          `json:"provider_response_sha256"`
	Outcomes                    []tiktokOutcome `json:"outcomes"`
	AutomaticBusinessWrite      bool            `json:"automatic_business_write"`
	AutomaticContactEligibility bool            `json:"automatic_contact_eligibility"`
}

type tiktokOutcome struct {
	Status               string           `json:"status"`
	NormalizationCodes   []string         `json:"normalization_codes"`
	ProviderLeadIDSHA256 string           `json:"provider_lead_id_sha256"`
	Candidate            *tiktokCandidate `json:"candidate"`
}

type tiktokCandidate struct {
	TenantID           string  `json:"tenant_id"`
	OrganizationID     string  `json:"organization_id"`
	Provider           string  `json:"provider"`
	ProviderLeadID     string  `json:"provider_lead_id"`
	FormID             string  `json:"form_id"`
	CampaignID         string  `json:"campaign_id"`
	AdGroupID          string  `json:"ad_group_id"`
	CreativeID         string  `json:"creative_id"`
	SourceKind         string  `json:"source_kind"`
	SubmittedAt        string  `json:"submitted_at"`
	IsTest             bool    `json:"is_test"`
	Fields             []Field `json:"fields"`
	ContactEligibility string  `json:"contact_eligibility"`
}

type tiktokRetrievalReceipt struct {
	Schema                      string `json:"schema"`
	SDKVersion                  string `json:"sdk_version"`
	SDKWheelSHA256              string `json:"sdk_wheel_sha256"`
	APIVersion                  string `json:"api_version"`
	ContractObservedAt          string `json:"contract_observed_at"`
	RetrievalMode               string `json:"retrieval_mode"`
	RequestIDSHA256             string `json:"request_id_sha256"`
	ProviderResponseSHA256      string `json:"provider_response_sha256"`
	CandidateBatchSHA256        string `json:"candidate_batch_sha256"`
	RowCount                    int    `json:"row_count"`
	CandidateCount              int    `json:"candidate_count"`
	RejectedCount               int    `json:"rejected_count"`
	AutomaticBusinessWrite      bool   `json:"automatic_business_write"`
	AutomaticContactEligibility bool   `json:"automatic_contact_eligibility"`
}

type TikTokImportReceipt struct {
	State    ReceiptState `json:"state"`
	Complete bool         `json:"complete"`
}

func ImportTikTokEvidence(ctx context.Context, store Store, providerResponse, candidateBatch, retrievalReceipt []byte, receivedAt time.Time) (TikTokImportReceipt, error) {
	var result TikTokImportReceipt
	if store == nil || receivedAt.IsZero() {
		return result, errors.New("leadstream: store and received_at are required")
	}
	raw, candidate, code, err := prepareTikTokRecord(providerResponse, candidateBatch, retrievalReceipt, receivedAt.UTC())
	if err != nil {
		return result, err
	}
	receipt, err := store.Record(ctx, raw, candidate, code)
	if err != nil {
		return result, err
	}
	result.State, result.Complete = receipt.State, true
	return result, nil
}

func prepareTikTokRecord(providerBytes, batchBytes, receiptBytes []byte, receivedAt time.Time) (RawEvent, *LeadCandidate, string, error) {
	if len(providerBytes) == 0 || len(batchBytes) == 0 || len(receiptBytes) == 0 {
		return RawEvent{}, nil, "", errors.New("leadstream: all TikTok evidence artifacts are required")
	}
	var provider tiktokProviderResponse
	var batch tiktokCandidateBatch
	var receipt tiktokRetrievalReceipt
	if err := decodeTikTokStrict(providerBytes, &provider); err != nil {
		return RawEvent{}, nil, "", fmt.Errorf("leadstream: invalid TikTok provider artifact: %w", err)
	}
	if err := decodeTikTokStrict(batchBytes, &batch); err != nil {
		return RawEvent{}, nil, "", fmt.Errorf("leadstream: invalid TikTok candidate artifact: %w", err)
	}
	if err := decodeTikTokStrict(receiptBytes, &receipt); err != nil {
		return RawEvent{}, nil, "", fmt.Errorf("leadstream: invalid TikTok receipt: %w", err)
	}
	if provider.Schema != "elite-tiktok-lead-sdk-response/v1" || provider.SDKVersion != tiktokSDKVersion || provider.APIVersion != tiktokAPIVersion || strings.TrimSpace(provider.RequestID) == "" {
		return RawEvent{}, nil, "", errors.New("leadstream: unsupported TikTok provider identity")
	}
	if batch.Schema != "elite-tiktok-lead-candidate-batch/v1" || batch.Provider != tiktokProvider || batch.SDKVersion != tiktokSDKVersion || batch.SDKWheelSHA256 != tiktokSDKWheelSHA256 || batch.APIVersion != tiktokAPIVersion || batch.ContractObservedAt != tiktokContractObserved || batch.AutomaticBusinessWrite || batch.AutomaticContactEligibility || len(batch.Outcomes) != 1 {
		return RawEvent{}, nil, "", errors.New("leadstream: unsupported or unsafe TikTok candidate batch")
	}
	if receipt.Schema != "elite-tiktok-lead-retrieval-receipt/v1" || receipt.SDKVersion != tiktokSDKVersion || receipt.SDKWheelSHA256 != tiktokSDKWheelSHA256 || receipt.APIVersion != tiktokAPIVersion || receipt.ContractObservedAt != tiktokContractObserved || receipt.AutomaticBusinessWrite || receipt.AutomaticContactEligibility {
		return RawEvent{}, nil, "", errors.New("leadstream: unsupported or unsafe TikTok retrieval receipt")
	}
	if shaHex(providerBytes) != batch.ProviderResponseSHA256 || receipt.ProviderResponseSHA256 != batch.ProviderResponseSHA256 || shaHex(batchBytes) != receipt.CandidateBatchSHA256 || shaHex([]byte(provider.RequestID)) != receipt.RequestIDSHA256 {
		return RawEvent{}, nil, "", errors.New("leadstream: TikTok evidence hash mismatch")
	}
	if receipt.RowCount != 1 || receipt.CandidateCount+receipt.RejectedCount != 1 || receipt.RetrievalMode != provider.RequestIdentity.Mode || (receipt.RetrievalMode != "TEST" && receipt.RetrievalMode != "LIVE") {
		return RawEvent{}, nil, "", errors.New("leadstream: TikTok receipt cardinality or mode mismatch")
	}
	advertiserSelectorSet := provider.RequestIdentity.AdvertiserIDSHA256 != ""
	librarySelectorSet := provider.RequestIdentity.LibraryIDSHA256 != ""
	if (advertiserSelectorSet && !validSHA256(provider.RequestIdentity.AdvertiserIDSHA256)) ||
		(librarySelectorSet && !validSHA256(provider.RequestIdentity.LibraryIDSHA256)) ||
		advertiserSelectorSet == librarySelectorSet {
		return RawEvent{}, nil, "", errors.New("leadstream: TikTok account selector identity mismatch")
	}
	leadID, err := tiktokMetaString(provider.Data.MetaData, "lead_id", true)
	if err != nil || shaHex([]byte(leadID)) != batch.Outcomes[0].ProviderLeadIDSHA256 {
		return RawEvent{}, nil, "", errors.New("leadstream: TikTok lead identity mismatch")
	}
	leadSource, err := tiktokMetaString(provider.Data.MetaData, "lead_source", true)
	if err != nil || leadSource != provider.RequestIdentity.LeadSource || (leadSource != "INSTANT_FORM" && leadSource != "DIRECT_MESSAGE") {
		return RawEvent{}, nil, "", errors.New("leadstream: TikTok lead source mismatch")
	}
	pageID, err := tiktokMetaString(provider.Data.MetaData, "page_id", false)
	if err != nil {
		return RawEvent{}, nil, "", err
	}
	if (leadSource == "INSTANT_FORM" && (!validSHA256(provider.RequestIdentity.PageIDSHA256) || shaHex([]byte(pageID)) != provider.RequestIdentity.PageIDSHA256)) || (leadSource == "DIRECT_MESSAGE" && (pageID != "" || provider.RequestIdentity.PageIDSHA256 != "")) {
		return RawEvent{}, nil, "", errors.New("leadstream: TikTok page identity mismatch")
	}
	created, err := tiktokMetaString(provider.Data.MetaData, "create_time", true)
	if err != nil {
		return RawEvent{}, nil, "", err
	}
	occurredAt, err := time.ParseInLocation("2006-01-02 15:04:05", created, time.UTC)
	if err != nil {
		return RawEvent{}, nil, "", errors.New("leadstream: invalid TikTok create_time")
	}
	raw, err := NewRawEvent(batch.TenantID, batch.OrganizationID, tiktokProvider, leadID, "tiktok.lead.reconciled.v1.3", "tiktok-business-api-sdk-official:lead/get", tiktokAPIVersion, occurredAt, receivedAt, providerBytes)
	if err != nil {
		return RawEvent{}, nil, "", err
	}
	outcome := batch.Outcomes[0]
	if outcome.Status == "REJECTED" {
		if outcome.Candidate != nil || len(outcome.NormalizationCodes) == 0 || receipt.RejectedCount != 1 || receipt.CandidateCount != 0 {
			return RawEvent{}, nil, "", errors.New("leadstream: TikTok rejected outcome mismatch")
		}
		codes := append([]string(nil), outcome.NormalizationCodes...)
		sort.Strings(codes)
		for _, code := range codes {
			if !tiktokNormalizationCode.MatchString(code) {
				return RawEvent{}, nil, "", errors.New("leadstream: invalid TikTok normalization code")
			}
		}
		joined := strings.Join(codes, "+")
		if len(joined) > 128 {
			return RawEvent{}, nil, "", errors.New("leadstream: TikTok normalization code set too long")
		}
		return raw, nil, joined, nil
	}
	if outcome.Status != "CANDIDATE" || outcome.Candidate == nil || len(outcome.NormalizationCodes) != 0 || receipt.CandidateCount != 1 || receipt.RejectedCount != 0 {
		return RawEvent{}, nil, "", errors.New("leadstream: TikTok candidate outcome mismatch")
	}
	wire := outcome.Candidate
	campaignID, err := tiktokMetaString(provider.Data.MetaData, "campaign_id", false)
	if err != nil {
		return RawEvent{}, nil, "", err
	}
	adGroupID, err := tiktokMetaString(provider.Data.MetaData, "adgroup_id", false)
	if err != nil {
		return RawEvent{}, nil, "", err
	}
	creativeID, err := tiktokMetaString(provider.Data.MetaData, "ad_id", false)
	if err != nil {
		return RawEvent{}, nil, "", err
	}
	if wire.TenantID != batch.TenantID || wire.OrganizationID != batch.OrganizationID || wire.Provider != tiktokProvider || wire.ProviderLeadID != leadID || wire.FormID != pageID || wire.CampaignID != campaignID || wire.AdGroupID != adGroupID || wire.CreativeID != creativeID || wire.SourceKind != leadSource || wire.IsTest != (provider.RequestIdentity.Mode == "TEST") || wire.ContactEligibility != "pending_policy" {
		return RawEvent{}, nil, "", errors.New("leadstream: TikTok candidate does not match provider evidence")
	}
	submitted, err := time.Parse(time.RFC3339, wire.SubmittedAt)
	if err != nil || !submitted.Equal(occurredAt) {
		return RawEvent{}, nil, "", errors.New("leadstream: TikTok submitted_at mismatch")
	}
	if err := matchTikTokFields(provider.Data.LeadData, wire.Fields); err != nil {
		return RawEvent{}, nil, "", err
	}
	candidate := &LeadCandidate{TenantID: wire.TenantID, OrganizationID: wire.OrganizationID, Provider: wire.Provider, ProviderLeadID: wire.ProviderLeadID, FormID: wire.FormID, CampaignID: wire.CampaignID, AdGroupID: wire.AdGroupID, CreativeID: wire.CreativeID, SourceKind: wire.SourceKind, SubmittedAt: submitted, IsTest: wire.IsTest, Fields: wire.Fields, ContactEligibility: wire.ContactEligibility}
	if err := candidate.Validate(); err != nil {
		return RawEvent{}, nil, "", err
	}
	return raw, candidate, "", nil
}

func decodeTikTokStrict(value []byte, target any) error {
	decoder := json.NewDecoder(bytes.NewReader(value))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return err
	}
	return ensureEOF(decoder)
}

func tiktokMetaString(meta map[string]json.RawMessage, name string, required bool) (string, error) {
	raw, exists := meta[name]
	if !exists || bytes.Equal(raw, []byte("null")) {
		if required {
			return "", fmt.Errorf("leadstream: TikTok meta field required: %s", name)
		}
		return "", nil
	}
	var value string
	if err := json.Unmarshal(raw, &value); err != nil || strings.TrimSpace(value) == "" {
		return "", fmt.Errorf("leadstream: invalid TikTok meta field: %s", name)
	}
	return value, nil
}

func matchTikTokFields(raw map[string]json.RawMessage, normalized []Field) error {
	if len(raw) != len(normalized) {
		return errors.New("leadstream: TikTok field cardinality mismatch")
	}
	values := make(map[string]string, len(raw))
	for name, encoded := range raw {
		var value string
		if err := json.Unmarshal(encoded, &value); err != nil || strings.TrimSpace(name) == "" || strings.TrimSpace(value) == "" {
			return errors.New("leadstream: TikTok field is not losslessly string-normalizable")
		}
		values[name] = value
	}
	for _, field := range normalized {
		if values[field.ID] != field.Value {
			return errors.New("leadstream: TikTok normalized field mismatch")
		}
	}
	return nil
}

func tiktokSHA(value string) string {
	hash := sha256.Sum256([]byte(value))
	return hex.EncodeToString(hash[:])
}
````

### FILE: `internal/leadstream/tiktok_import_test.go`
```yaml
block_id: "GO-TIKTOK-LEAD-IMPORT:domain-tests:v1"
operation: CREATE
provenance: AUTHORED
source: "local integrity, replay, mismatch and rejection tests"
license: "LicenseRef-Workspace-Owner"
sha256: "16fa30a6fb8adc7518f583378dea7b3e67be95ae60ec3627cc7f2e8579271003"
variables: []
secrets_allowed: false
```
````go
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
````

### FILE: `cmd/tiktok-lead-import/main.go`
```yaml
block_id: "GO-TIKTOK-LEAD-IMPORT:command:v1"
operation: CREATE
provenance: AUTHORED
source: "local bounded command over the canonical PostgreSQL lead owner"
license: "LicenseRef-Workspace-Owner"
sha256: "03166a3a0364b7d0ab82bfeb5f5683c11205c29032db0b9e7a84ae40b46d6f88"
variables: []
secrets_allowed: false
```
````go
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"

	"elite.local/enterprise/internal/leadstream"
	"elite.local/enterprise/internal/platform/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

const maxArtifactBytes int64 = 32 << 20

func main() {
	providerPath := flag.String("provider-response", "", "TikTok provider-response.json")
	batchPath := flag.String("candidate-batch", "", "TikTok candidate-batch.json")
	receiptPath := flag.String("retrieval-receipt", "", "TikTok retrieval-receipt.json")
	flag.Parse()
	if *providerPath == "" || *batchPath == "" || *receiptPath == "" || os.Getenv("DATABASE_URL") == "" {
		slog.Error("DATABASE_URL and all three TikTok artifact paths are required")
		os.Exit(2)
	}
	providerBytes, err := readBounded(*providerPath)
	if err != nil {
		fail("provider response unavailable", err)
	}
	batchBytes, err := readBounded(*batchPath)
	if err != nil {
		fail("candidate batch unavailable", err)
	}
	receiptBytes, err := readBounded(*receiptPath)
	if err != nil {
		fail("retrieval receipt unavailable", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	config, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		fail("database configuration invalid", err)
	}
	config.MaxConns = 4
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		fail("database pool unavailable", err)
	}
	defer pool.Close()
	if err := pool.Ping(ctx); err != nil {
		fail("database unavailable", err)
	}
	result, err := leadstream.ImportTikTokEvidence(ctx, postgres.NewLeadIngress(pool), providerBytes, batchBytes, receiptBytes, time.Now().UTC())
	encoded, _ := json.Marshal(result)
	if err != nil {
		fmt.Fprintln(os.Stderr, string(encoded))
		fail("TikTok lead import incomplete; replay is safe", err)
	}
	fmt.Println(string(encoded))
}

func readBounded(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	value, err := io.ReadAll(io.LimitReader(file, maxArtifactBytes+1))
	if err != nil {
		return nil, err
	}
	if int64(len(value)) > maxArtifactBytes {
		return nil, fmt.Errorf("artifact exceeds %d bytes", maxArtifactBytes)
	}
	if len(value) == 0 {
		return nil, fmt.Errorf("artifact is empty")
	}
	return value, nil
}

func fail(message string, err error) {
	slog.Error(message, "error", err)
	os.Exit(1)
}
````

### FILE: `internal/platform/postgres/tiktok_import_integration_test.go`
```yaml
block_id: "GO-TIKTOK-LEAD-IMPORT:postgres-test:v1"
operation: CREATE
provenance: AUTHORED
source: "local real PostgreSQL shared-owner/outbox integration regression"
license: "LicenseRef-Workspace-Owner"
sha256: "91ec0a5ca4c45b2523b4918fe85a266c141d314e21bb3aa5f8b10f0e8f5915d1"
variables: []
secrets_allowed: false
```
````go
package postgres

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"testing"
	"time"

	"elite.local/enterprise/internal/leadstream"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestTikTokImportUsesSharedDurableOwnerAndReplays(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := leadstream.StableUUID(t.Name(), time.Now().UTC().Format(time.RFC3339Nano))
	organization := "store-1"
	if _, err = pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,$2,'TikTok Import Test','TikTok Import Test')`, tenant, "tiktok-"+tenant[:8]); err != nil {
		t.Fatal(err)
	}
	if _, err = pool.Exec(ctx, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,$2,$2,'TikTok Store','store')`, tenant, organization); err != nil {
		t.Fatal(err)
	}
	defer func() {
		_, _ = pool.Exec(context.Background(), `update platform.outbox_event set published_at=greatest(clock_timestamp(),occurred_at),claimed_by=null,claimed_until=null where tenant_id=$1 and published_at is null`, tenant)
	}()
	provider, batch, receipt := tiktokIntegrationArtifacts(t, tenant, organization)
	store := NewLeadIngress(pool)
	first, err := leadstream.ImportTikTokEvidence(ctx, store, provider, batch, receipt, time.Now().UTC())
	if err != nil || first.State != leadstream.ReceiptNormalized || !first.Complete {
		t.Fatalf("first=%+v err=%v", first, err)
	}
	replay, err := leadstream.ImportTikTokEvidence(ctx, store, provider, batch, receipt, time.Now().UTC())
	if err != nil || replay.State != leadstream.ReceiptDuplicate {
		t.Fatalf("replay=%+v err=%v", replay, err)
	}
	var rawCount, candidateCount, outboxCount int
	if err = pool.QueryRow(ctx, `select count(*) from integration.lead_ingress_raw where tenant_id=$1 and provider='tiktok_lead_generation' and provider_event_id='lead-1'`, tenant).Scan(&rawCount); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select count(*) from integration.lead_candidate where tenant_id=$1 and provider='tiktok_lead_generation' and provider_event_id='lead-1' and contact_eligibility='pending_policy'`, tenant).Scan(&candidateCount); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `select count(*) from platform.outbox_event where tenant_id=$1 and aggregate_id='tiktok_lead_generation:lead-1'`, tenant).Scan(&outboxCount); err != nil {
		t.Fatal(err)
	}
	if rawCount != 1 || candidateCount != 1 || outboxCount != 1 {
		t.Fatalf("raw=%d candidate=%d outbox=%d", rawCount, candidateCount, outboxCount)
	}
}

func tiktokIntegrationArtifacts(t *testing.T, tenant, organization string) ([]byte, []byte, []byte) {
	t.Helper()
	hash := func(value []byte) string {
		sum := sha256.Sum256(value)
		return hex.EncodeToString(sum[:])
	}
	providerObject := map[string]any{"schema": "elite-tiktok-lead-sdk-response/v1", "sdk_version": "1.1.3", "api_version": "v1.3", "request_identity": map[string]any{"mode": "TEST", "lead_source": "INSTANT_FORM", "advertiser_id_sha256": hash([]byte("123")), "library_id_sha256": "", "page_id_sha256": hash([]byte("456"))}, "data": map[string]any{"lead_data": map[string]any{"Email": "person@example.test"}, "meta_data": map[string]any{"lead_source": "INSTANT_FORM", "lead_id": "lead-1", "page_id": "456", "campaign_id": "11", "adgroup_id": "22", "ad_id": "33", "create_time": "2026-09-04 12:00:00"}}, "request_id": "req-1"}
	provider, _ := json.Marshal(providerObject)
	provider = append(provider, '\n')
	candidate := map[string]any{"tenant_id": tenant, "organization_id": organization, "provider": "tiktok_lead_generation", "provider_lead_id": "lead-1", "form_id": "456", "campaign_id": "11", "ad_group_id": "22", "creative_id": "33", "source_kind": "INSTANT_FORM", "submitted_at": "2026-09-04T12:00:00Z", "is_test": true, "fields": []map[string]string{{"id": "Email", "value": "person@example.test"}}, "contact_eligibility": "pending_policy"}
	batchObject := map[string]any{"schema": "elite-tiktok-lead-candidate-batch/v1", "provider": "tiktok_lead_generation", "sdk_version": "1.1.3", "sdk_wheel_sha256": "663b4a1f2585c4b386144d968646f695befd733c464418bbaeeac41befff33b7", "api_version": "v1.3", "contract_observed_at": "2026-09-04", "tenant_id": tenant, "organization_id": organization, "provider_response_sha256": hash(provider), "outcomes": []any{map[string]any{"status": "CANDIDATE", "normalization_codes": []string{}, "provider_lead_id_sha256": hash([]byte("lead-1")), "candidate": candidate}}, "automatic_business_write": false, "automatic_contact_eligibility": false}
	batch, _ := json.Marshal(batchObject)
	batch = append(batch, '\n')
	receiptObject := map[string]any{"schema": "elite-tiktok-lead-retrieval-receipt/v1", "sdk_version": "1.1.3", "sdk_wheel_sha256": "663b4a1f2585c4b386144d968646f695befd733c464418bbaeeac41befff33b7", "api_version": "v1.3", "contract_observed_at": "2026-09-04", "retrieval_mode": "TEST", "request_id_sha256": hash([]byte("req-1")), "provider_response_sha256": hash(provider), "candidate_batch_sha256": hash(batch), "row_count": 1, "candidate_count": 1, "rejected_count": 0, "automatic_business_write": false, "automatic_contact_eligibility": false}
	receipt, _ := json.Marshal(receiptObject)
	return provider, batch, append(receipt, '\n')
}
````

### FILE: `docs/tiktok-lead-import.md`
```yaml
block_id: "GO-TIKTOK-LEAD-IMPORT:guide:v1"
operation: CREATE
provenance: AUTHORED
source: "local operator and limitations contract"
license: "LicenseRef-Workspace-Owner"
sha256: "5264ad5866527933049a3f394de9c192892350fc1eddf7b24f19fae034b17c64"
variables: []
secrets_allowed: false
```
````markdown
# TikTok Lead durable import

This integration consumes only the three hash-linked artifacts emitted by `PYTHON-TIKTOK-LEAD-ADAPTER` 0.2.0. It verifies SDK/wheel/API/contract identity, hashes, retrieval mode, account/page identity hashes, provider metadata, candidate fields and outcome counts before the first write.

The import reuses the provider-neutral `leadstream.Store` and PostgreSQL migration 0044. It creates no parallel CRM, raw-event table, candidate table or outbox. A valid candidate is stored with `contact_eligibility=pending_policy`; a rejected dynamic field remains evidence with a normalization code. Replay is idempotent and a reused provider lead ID with different source bytes fails closed.

This pack does not authenticate TikTok webhook calls. The webhook remains an untrusted hint until the project demonstrates a provider-supported or independently approved ingress control. The authenticated retrieval response—not the webhook—is the input to durable candidate import. Live account, terms/consent, callback, delivery, reconciliation, retention and cost still require project evidence.
````

## 6. Configuration surface

Requiere `DATABASE_URL` sólo en runtime y tres paths de artefactos sin secretos. Tenant, organización, modo y provider quedan dentro de artefactos hash-linked producidos por el profile TikTok `PROVEN`. No modifica policies, mappings ni consentimiento.

## 7. Dependency bill

No añade módulos Go. Reutiliza el `go.mod`, pgx, `leadstream.Store`, `postgres.LeadIngress`, migración 0044 y outbox canónicos del backend compuesto.

## 8. Apply order

1. Componer después de `GO-OMNICHANNEL-LEAD-INGRESS` y junto con `PYTHON-TIKTOK-LEAD-ADAPTER` 0.2.0.
2. Ejecutar los cinco tests focales de `internal/leadstream` y full Go test/vet/build.
3. Aplicar migraciones existentes hasta 0049 y ejecutar el test PostgreSQL TikTok sin skip.
4. Ejecutar retrieval TEST real, importar tres artefactos y verificar raw/candidate/outbox/replay.
5. Mantener promoción/contacto en el owner explícito posterior.

## 9. Verification

```powershell
go test ./internal/leadstream
go test ./internal/platform/postgres -run TestTikTokImportUsesSharedDurableOwnerAndReplays -count=1
go test ./...
go vet ./...
go build ./...
```

Resultado local demostrado antes de promoción final: cinco tests domain y compilación del command. La admisión completa exige round-trip Markdown, PostgreSQL 18.6 real sin skip, full Go/vet/build y composición 62+1 desde vacío. Cuenta TikTok, delivery/retrieval live y producción siguen condicionados al proyecto.

## 10. Reconstruction evidence

- Input exacto: los tres artifacts del pack TikTok Lead 0.2.0, wheel SHA `663b4a…f33b7` y contrato v1.3 observado.
- Output único: `leadstream.Store`; PostgreSQL owner 0044 y outbox existente.
- Procedencia: cinco bloques `AUTHORED`; no se presentan como código TikTok ni copian un SDK Lead inexistente.
