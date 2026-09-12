# Go Meta Lead Evidence Import

## 1. Metadata

```yaml
pack_id: "GO-META-LEAD-EVIDENCE-IMPORT"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa el enlace verificable desde los tres artefactos exactos V226 de Meta hasta el owner PostgreSQL V224, validando hashes, identidad, scope, cardinalidad, raw/candidato y replay antes de permitir la promoción V225."
stacks: ["Go 1.26.7", "PostgreSQL 18.6", "Meta V226 artifact contract"]
compatible_with: ["GO-OMNICHANNEL-LEAD-INGRESS 0.2.x", "PYTHON-META-LEAD-RECONCILIATION-ADAPTER 0.1.x", "GO-LEAD-CANDIDATE-PROMOTION 0.1.x"]
incompatible_with: ["webhook Meta no demostrado", "candidate batch no V226", "write CRM directo", "contacto automático", "artefacto sin receipt/hash"]
license_expression: "LicenseRef-Workspace-Owner AND LicenseRef-Meta-Platform"
upstream_sources: ["https://github.com/facebook/facebook-python-business-sdk/blob/788f363d15b1269ab5efb7cd00fb5e3b133cd99b/facebook_business/adobjects/lead.py", "https://github.com/facebook/facebook-python-business-sdk/blob/788f363d15b1269ab5efb7cd00fb5e3b133cd99b/facebook_business/adobjects/leadgenform.py"]
verified_at: "2026-09-04"
```

## 2. Applicability

Use después de V226 y V224. El importer verifica `provider-response.json`, `lead-candidates.json` y `RETRIEVAL_RECEIPT.json` completos antes de la primera escritura. No llama Meta, no recibe webhooks y no promociona al CRM: entrega evidencia raw/candidato al owner durable existente.

El código es `ADAPTED` a partir del contrato oficial Meta fijado, V226 y los owners V224/V225. Ninguna línea se atribuye falsamente a Meta. Las fixtures son salidas locales generadas por el productor V226 exacto.

## 3. Architecture contract

- Verifica SDK/API/source commit, schemas, hashes response→batch→receipt, cardinalidades y flags de no-write/no-contact.
- Comprueba que lead/form/fecha/campaña/adset/ad/platform/fields normalizados coincidan con el raw exportado.
- Exige scope tenant/organization a nivel batch y su igualdad con cada candidato.
- Persiste rechazo con su scope aunque `candidate=null`.
- Procesa sólo después de validar el lote completo; una falla de Store puede dejar un prefijo durable explícito y el replay converge por idempotencia `(tenant, provider, provider_event_id)`.
- No oculta truncamiento de la fuente.

## 4. Exact file manifest

```text
CREATE internal/leadstream/meta_import.go
CREATE internal/leadstream/meta_import_test.go
CREATE internal/leadstream/testdata/meta-v226-valid/provider-response.json
CREATE internal/leadstream/testdata/meta-v226-valid/lead-candidates.json
CREATE internal/leadstream/testdata/meta-v226-valid/RETRIEVAL_RECEIPT.json
CREATE internal/platform/postgres/meta_import_integration_test.go
CREATE cmd/meta-lead-import/main.go
```

## 5. Materialization blocks

### FILE: `internal/leadstream/meta_import.go`
```yaml
block_id: "GO-META-LEAD-IMPORT:core:v1"
operation: CREATE
provenance: ADAPTED
source: "exact Meta Lead/LeadgenForm contracts plus V224/V226 local artifact and durability contracts"
license: "LicenseRef-Meta-Platform AND LicenseRef-Workspace-Owner"
sha256: "1a51c555ceb30c42e339dfcabd216e76ae1781b5bcd218b1e9077b680dda2f2f"
variables: []
secrets_allowed: false
```
````go
package leadstream

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

const (
	metaSDKVersion   = "26.0.1"
	metaAPIVersion   = "v26.0"
	metaSourceCommit = "788f363d15b1269ab5efb7cd00fb5e3b133cd99b"
)

var metaNormalizationCode = regexp.MustCompile(`^[A-Z][A-Z0-9_]{0,63}$`)

type metaProviderResponse struct {
	GraphAPIVersion       string            `json:"graph_api_version"`
	RetrievalMode         string            `json:"retrieval_mode"`
	Rows                  []json.RawMessage `json:"rows"`
	TruncatedByLocalBound bool              `json:"truncated_by_local_bound"`
}

type metaProviderRow struct {
	ID          string          `json:"id"`
	CreatedTime string          `json:"created_time"`
	FormID      string          `json:"form_id"`
	CampaignID  string          `json:"campaign_id"`
	AdsetID     string          `json:"adset_id"`
	AdID        string          `json:"ad_id"`
	Platform    string          `json:"platform"`
	FieldData   []metaFieldData `json:"field_data"`
}

type metaFieldData struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

type metaCandidateBatch struct {
	Schema                      string        `json:"schema"`
	Provider                    string        `json:"provider"`
	SourceCommit                string        `json:"source_commit"`
	TenantID                    string        `json:"tenant_id"`
	OrganizationID              string        `json:"organization_id"`
	ProviderResponseSHA256      string        `json:"provider_response_sha256"`
	Outcomes                    []metaOutcome `json:"outcomes"`
	AutomaticBusinessWrite      bool          `json:"automatic_business_write"`
	AutomaticContactEligibility bool          `json:"automatic_contact_eligibility"`
}

type metaOutcome struct {
	Status               string             `json:"status"`
	NormalizationCodes   []string           `json:"normalization_codes"`
	Candidate            *metaWireCandidate `json:"candidate"`
	ProviderLeadIDSHA256 string             `json:"provider_lead_id_sha256"`
}

type metaWireCandidate struct {
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

type metaRetrievalReceipt struct {
	Schema                      string `json:"schema"`
	SDKVersion                  string `json:"sdk_version"`
	GraphAPIVersion             string `json:"graph_api_version"`
	SourceCommit                string `json:"source_commit"`
	FormIDSHA256                string `json:"form_id_sha256"`
	ProviderResponseSHA256      string `json:"provider_response_sha256"`
	CandidateBatchSHA256        string `json:"candidate_batch_sha256"`
	RowCount                    int    `json:"row_count"`
	CandidateCount              int    `json:"candidate_count"`
	RejectedCount               int    `json:"rejected_count"`
	TruncatedByLocalBound       bool   `json:"truncated_by_local_bound"`
	AutomaticBusinessWrite      bool   `json:"automatic_business_write"`
	AutomaticContactEligibility bool   `json:"automatic_contact_eligibility"`
}

type MetaImportReceipt struct {
	Rows            int  `json:"rows"`
	Normalized      int  `json:"normalized"`
	Rejected        int  `json:"rejected"`
	Duplicates      int  `json:"duplicates"`
	Processed       int  `json:"processed"`
	TruncatedSource bool `json:"truncated_source"`
	Complete        bool `json:"complete"`
}

type preparedMetaRecord struct {
	raw       RawEvent
	candidate *LeadCandidate
	code      string
}

// ImportMetaEvidence verifies all three V226 artifacts before making the first
// durable write. Store errors may leave a strict prefix committed; the receipt
// reports progress and replay safely converges through provider idempotency.
func ImportMetaEvidence(ctx context.Context, store Store, providerResponse, candidateBatch, retrievalReceipt []byte, receivedAt time.Time) (MetaImportReceipt, error) {
	var result MetaImportReceipt
	if store == nil || receivedAt.IsZero() {
		return result, errors.New("leadstream: store and received_at are required")
	}
	records, truncated, err := prepareMetaRecords(providerResponse, candidateBatch, retrievalReceipt, receivedAt.UTC())
	if err != nil {
		return result, err
	}
	result.Rows, result.TruncatedSource = len(records), truncated
	for index, record := range records {
		receipt, recordErr := store.Record(ctx, record.raw, record.candidate, record.code)
		if recordErr != nil {
			result.Processed = index
			return result, fmt.Errorf("leadstream: Meta import stopped at row %d after %d durable records: %w", index, index, recordErr)
		}
		result.Processed = index + 1
		switch receipt.State {
		case ReceiptNormalized:
			result.Normalized++
		case ReceiptRejected:
			result.Rejected++
		case ReceiptDuplicate:
			result.Duplicates++
		default:
			return result, fmt.Errorf("leadstream: unexpected store receipt state %q", receipt.State)
		}
	}
	result.Complete = true
	return result, nil
}

func prepareMetaRecords(providerBytes, candidateBytes, receiptBytes []byte, receivedAt time.Time) ([]preparedMetaRecord, bool, error) {
	if len(providerBytes) == 0 || len(candidateBytes) == 0 || len(receiptBytes) == 0 {
		return nil, false, errors.New("leadstream: all Meta evidence artifacts are required")
	}
	var provider metaProviderResponse
	var batch metaCandidateBatch
	var receipt metaRetrievalReceipt
	if err := json.Unmarshal(providerBytes, &provider); err != nil {
		return nil, false, fmt.Errorf("leadstream: invalid Meta provider response: %w", err)
	}
	if err := json.Unmarshal(candidateBytes, &batch); err != nil {
		return nil, false, fmt.Errorf("leadstream: invalid Meta candidate batch: %w", err)
	}
	if err := json.Unmarshal(receiptBytes, &receipt); err != nil {
		return nil, false, fmt.Errorf("leadstream: invalid Meta retrieval receipt: %w", err)
	}
	providerHash, candidateHash := shaHex(providerBytes), shaHex(candidateBytes)
	if provider.GraphAPIVersion != metaAPIVersion || (provider.RetrievalMode != "TEST" && provider.RetrievalMode != "LIVE") {
		return nil, false, errors.New("leadstream: unsupported Meta provider evidence identity")
	}
	if batch.Schema != "elite-meta-lead-candidate-batch/v1" || batch.Provider != "meta_lead_ads" || batch.SourceCommit != metaSourceCommit || batch.AutomaticBusinessWrite || batch.AutomaticContactEligibility {
		return nil, false, errors.New("leadstream: unsupported or unsafe Meta candidate batch")
	}
	if receipt.Schema != "elite-meta-lead-reconciliation-receipt/v1" || receipt.SDKVersion != metaSDKVersion || receipt.GraphAPIVersion != metaAPIVersion || receipt.SourceCommit != metaSourceCommit || receipt.AutomaticBusinessWrite || receipt.AutomaticContactEligibility {
		return nil, false, errors.New("leadstream: unsupported or unsafe Meta retrieval receipt")
	}
	if batch.ProviderResponseSHA256 != providerHash || receipt.ProviderResponseSHA256 != providerHash || receipt.CandidateBatchSHA256 != candidateHash {
		return nil, false, errors.New("leadstream: Meta evidence hash mismatch")
	}
	if provider.TruncatedByLocalBound != receipt.TruncatedByLocalBound || len(provider.Rows) != len(batch.Outcomes) || len(provider.Rows) != receipt.RowCount {
		return nil, false, errors.New("leadstream: Meta evidence cardinality mismatch")
	}
	records := make([]preparedMetaRecord, 0, len(provider.Rows))
	candidates, rejected := 0, 0
	for index := range provider.Rows {
		record, status, err := prepareMetaRecord(provider.Rows[index], batch.Outcomes[index], batch.TenantID, batch.OrganizationID, provider.RetrievalMode, receipt.FormIDSHA256, receivedAt)
		if err != nil {
			return nil, false, fmt.Errorf("leadstream: invalid Meta row %d: %w", index, err)
		}
		if status == "CANDIDATE" {
			candidates++
		} else {
			rejected++
		}
		records = append(records, record)
	}
	if candidates != receipt.CandidateCount || rejected != receipt.RejectedCount || candidates+rejected != receipt.RowCount {
		return nil, false, errors.New("leadstream: Meta receipt outcome counts mismatch")
	}
	return records, provider.TruncatedByLocalBound, nil
}

func prepareMetaRecord(rawJSON json.RawMessage, outcome metaOutcome, tenantID, organizationID, mode, formHash string, receivedAt time.Time) (preparedMetaRecord, string, error) {
	var row metaProviderRow
	if err := json.Unmarshal(rawJSON, &row); err != nil {
		return preparedMetaRecord{}, "", err
	}
	if strings.TrimSpace(row.ID) == "" || len(row.ID) > 256 || shaHex([]byte(row.ID)) != outcome.ProviderLeadIDSHA256 || shaHex([]byte(row.FormID)) != formHash {
		return preparedMetaRecord{}, "", errors.New("provider lead/form identity mismatch")
	}
	occurredAt, err := parseMetaTime(row.CreatedTime)
	if err != nil {
		return preparedMetaRecord{}, "", err
	}
	raw, err := NewRawEvent(tenantID, organizationID, "meta_lead_ads", row.ID, "meta.lead.reconciled.v26.0", "meta-business-sdk:leadgenform", metaAPIVersion, occurredAt, receivedAt, rawJSON)
	if err != nil {
		return preparedMetaRecord{}, "", err
	}
	if outcome.Status == "REJECTED" {
		if outcome.Candidate != nil || len(outcome.NormalizationCodes) == 0 {
			return preparedMetaRecord{}, "", errors.New("rejected outcome contract mismatch")
		}
		for _, code := range outcome.NormalizationCodes {
			if !metaNormalizationCode.MatchString(code) {
				return preparedMetaRecord{}, "", errors.New("invalid normalization code")
			}
		}
		code := strings.Join(outcome.NormalizationCodes, "+")
		if len(code) > 128 {
			return preparedMetaRecord{}, "", errors.New("normalization code set is too long")
		}
		return preparedMetaRecord{raw: raw, code: code}, "REJECTED", nil
	}
	if outcome.Status != "CANDIDATE" || outcome.Candidate == nil || len(outcome.NormalizationCodes) != 0 {
		return preparedMetaRecord{}, "", errors.New("candidate outcome contract mismatch")
	}
	wire := outcome.Candidate
	if wire.TenantID != tenantID || wire.OrganizationID != organizationID || wire.Provider != "meta_lead_ads" || wire.ProviderLeadID != row.ID || wire.FormID != row.FormID || wire.CampaignID != row.CampaignID || wire.AdGroupID != row.AdsetID || wire.CreativeID != row.AdID || wire.SourceKind != row.Platform || wire.ContactEligibility != "pending_policy" || wire.IsTest != (mode == "TEST") {
		return preparedMetaRecord{}, "", errors.New("candidate does not match provider row")
	}
	submittedAt, err := parseMetaTime(wire.SubmittedAt)
	if err != nil || !submittedAt.Equal(occurredAt) {
		return preparedMetaRecord{}, "", errors.New("candidate submitted_at mismatch")
	}
	if err := matchMetaFields(row.FieldData, wire.Fields); err != nil {
		return preparedMetaRecord{}, "", err
	}
	candidate := &LeadCandidate{
		TenantID: wire.TenantID, OrganizationID: wire.OrganizationID, Provider: wire.Provider,
		ProviderLeadID: wire.ProviderLeadID, FormID: wire.FormID, CampaignID: wire.CampaignID,
		AdGroupID: wire.AdGroupID, CreativeID: wire.CreativeID, SourceKind: wire.SourceKind,
		SubmittedAt: submittedAt, IsTest: wire.IsTest, Fields: wire.Fields,
		ContactEligibility: wire.ContactEligibility,
	}
	if err := candidate.Validate(); err != nil {
		return preparedMetaRecord{}, "", err
	}
	return preparedMetaRecord{raw: raw, candidate: candidate}, "CANDIDATE", nil
}

func matchMetaFields(raw []metaFieldData, normalized []Field) error {
	if len(raw) != len(normalized) {
		return errors.New("candidate field cardinality mismatch")
	}
	values := make(map[string]string, len(raw))
	for _, field := range raw {
		if strings.TrimSpace(field.Name) == "" || len(field.Values) != 1 || strings.TrimSpace(field.Values[0]) == "" {
			return errors.New("provider field is not losslessly normalizable")
		}
		if _, exists := values[field.Name]; exists {
			return errors.New("duplicate provider field")
		}
		values[field.Name] = field.Values[0]
	}
	for _, field := range normalized {
		if values[field.ID] != field.Value {
			return errors.New("normalized field does not match provider evidence")
		}
	}
	return nil
}

func parseMetaTime(value string) (time.Time, error) {
	for _, layout := range []string{time.RFC3339Nano, "2006-01-02T15:04:05-0700"} {
		if parsed, err := time.Parse(layout, value); err == nil {
			return parsed.UTC(), nil
		}
	}
	return time.Time{}, errors.New("unsupported Meta created_time")
}

func shaHex(value []byte) string {
	hash := sha256.Sum256(value)
	return hex.EncodeToString(hash[:])
}
````

### FILE: `internal/leadstream/meta_import_test.go`
```yaml
block_id: "GO-META-LEAD-IMPORT:test:v1"
operation: CREATE
provenance: ADAPTED
source: "V224 replay/divergence regressions plus exact V226 artifact contract"
license: "LicenseRef-Meta-Platform AND LicenseRef-Workspace-Owner"
sha256: "29e3cefea572e1c16b2ec0a5d4f09dc258481f0ed58198e6bcf6e98a8e08c936"
variables: []
secrets_allowed: false
```
````go
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
````

### FILE: `internal/leadstream/testdata/meta-v226-valid/provider-response.json`
```yaml
block_id: "GO-META-LEAD-IMPORT:fixture-provider:v1"
operation: CREATE
provenance: AUTHORED
source: "deterministic test output produced by V226 with a fake official-SDK row"
license: "LicenseRef-Workspace-Owner"
sha256: "fec8986d87cb71e9dfdffb40494ac8aa579c7a56a613605ec6978fe920a2d798"
variables: []
secrets_allowed: false
```
````json
{
  "graph_api_version": "v26.0",
  "retrieval_mode": "LIVE",
  "rows": [
    {
      "ad_id": "ad-1",
      "adset_id": "adset-1",
      "campaign_id": "campaign-1",
      "created_time": "2026-09-04T12:00:00+0000",
      "field_data": [
        {
          "name": "email",
          "values": [
            "person@example.test"
          ]
        }
      ],
      "form_id": "123456789",
      "id": "meta-cross-language-1",
      "platform": "facebook"
    }
  ],
  "truncated_by_local_bound": false
}
````

### FILE: `internal/leadstream/testdata/meta-v226-valid/lead-candidates.json`
```yaml
block_id: "GO-META-LEAD-IMPORT:fixture-candidate:v1"
operation: CREATE
provenance: AUTHORED
source: "deterministic test output produced by V226; example-only identities"
license: "LicenseRef-Workspace-Owner"
sha256: "57a5fc4e5a4fce062b244f5c0937fb8e4cb1b43caed8dd00e4efbfe22ebab0c9"
variables: []
secrets_allowed: false
```
````json
{
  "automatic_business_write": false,
  "automatic_contact_eligibility": false,
  "organization_id": "store-1",
  "outcomes": [
    {
      "candidate": {
        "ad_group_id": "adset-1",
        "campaign_id": "campaign-1",
        "contact_eligibility": "pending_policy",
        "creative_id": "ad-1",
        "fields": [
          {
            "id": "email",
            "value": "person@example.test"
          }
        ],
        "form_id": "123456789",
        "is_test": false,
        "organization_id": "store-1",
        "provider": "meta_lead_ads",
        "provider_lead_id": "meta-cross-language-1",
        "source_kind": "facebook",
        "submitted_at": "2026-09-04T12:00:00+0000",
        "tenant_id": "44444444-4444-4444-8444-444444444444"
      },
      "normalization_codes": [],
      "provider_lead_id_sha256": "a5f7a209773b0eba1bc806daa8f72d48573ebd6ed803200c5401dcfcd347e5f6",
      "status": "CANDIDATE"
    }
  ],
  "provider": "meta_lead_ads",
  "provider_response_sha256": "fec8986d87cb71e9dfdffb40494ac8aa579c7a56a613605ec6978fe920a2d798",
  "schema": "elite-meta-lead-candidate-batch/v1",
  "source_commit": "788f363d15b1269ab5efb7cd00fb5e3b133cd99b",
  "tenant_id": "44444444-4444-4444-8444-444444444444"
}
````

### FILE: `internal/leadstream/testdata/meta-v226-valid/RETRIEVAL_RECEIPT.json`
```yaml
block_id: "GO-META-LEAD-IMPORT:fixture-receipt:v1"
operation: CREATE
provenance: AUTHORED
source: "hash-linked deterministic-format test receipt produced by V226"
license: "LicenseRef-Workspace-Owner"
sha256: "dd37223aaf77a733b474496ed2f446ba92d7c91c9343360ebe44602e876c491e"
variables: []
secrets_allowed: false
```
````json
{
  "automatic_business_write": false,
  "automatic_contact_eligibility": false,
  "candidate_batch_sha256": "57a5fc4e5a4fce062b244f5c0937fb8e4cb1b43caed8dd00e4efbfe22ebab0c9",
  "candidate_count": 1,
  "created_at": "2026-09-04T18:28:33.686967+00:00",
  "form_id_sha256": "15e2b0d3c33891ebb0f1ef609ec419420c20e320ce94c65fbc8c3312448eb225",
  "graph_api_version": "v26.0",
  "provider": "Meta Marketing API",
  "provider_response_sha256": "fec8986d87cb71e9dfdffb40494ac8aa579c7a56a613605ec6978fe920a2d798",
  "query_sha256": "bc967b21ffaf80ed9f0e4aa41e2bf1ef3c321867e6ceac95b1dc48348f665c6a",
  "rejected_count": 0,
  "row_count": 1,
  "schema": "elite-meta-lead-reconciliation-receipt/v1",
  "sdk_version": "26.0.1",
  "source_commit": "788f363d15b1269ab5efb7cd00fb5e3b133cd99b",
  "truncated_by_local_bound": false
}
````

### FILE: `internal/platform/postgres/meta_import_integration_test.go`
```yaml
block_id: "GO-META-LEAD-IMPORT:postgres-test:v1"
operation: CREATE
provenance: ADAPTED
source: "V224 PostgreSQL owner tests extended to V226 Meta artifacts"
license: "LicenseRef-Workspace-Owner"
sha256: "062b40c83027aa6ead8687174f12294c14db604c9daf8209c35988b3acbc12c8"
variables: []
secrets_allowed: false
```
````go
package postgres

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"os"
	"testing"
	"time"

	"elite.local/enterprise/internal/leadstream"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestMetaEvidenceImportDurabilityRejectionAndReplay(t *testing.T) {
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
	const organization = "meta-store"
	if _, err := pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,$2,'Meta Import Test','Meta Import Test')`, tenant, "meta-"+tenant[:8]); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,$2,$2,'Meta Store','store')`, tenant, organization); err != nil {
		t.Fatal(err)
	}
	defer func() {
		_, _ = pool.Exec(context.Background(), `update platform.outbox_event set published_at=greatest(clock_timestamp(),occurred_at),claimed_by=null,claimed_until=null where tenant_id=$1 and published_at is null`, tenant)
	}()

	provider, batch, receipt := postgresMetaArtifacts(t, tenant, organization)
	repository := NewLeadIngress(pool)
	first, err := leadstream.ImportMetaEvidence(ctx, repository, provider, batch, receipt, time.Now().UTC())
	if err != nil || !first.Complete || first.Normalized != 1 || first.Rejected != 1 {
		t.Fatalf("first=%+v err=%v", first, err)
	}
	replay, err := leadstream.ImportMetaEvidence(ctx, repository, provider, batch, receipt, time.Now().UTC())
	if err != nil || replay.Duplicates != 2 || !replay.Complete {
		t.Fatalf("replay=%+v err=%v", replay, err)
	}

	var rawCount, candidateCount, rejectedCount, outboxCount int
	if err := pool.QueryRow(ctx, `select count(*) from integration.lead_ingress_raw where tenant_id=$1 and provider='meta_lead_ads' and provider_event_id in ('meta-live-1','meta-rejected-1')`, tenant).Scan(&rawCount); err != nil {
		t.Fatal(err)
	}
	if err := pool.QueryRow(ctx, `select count(*) from integration.lead_candidate where tenant_id=$1 and provider='meta_lead_ads' and provider_event_id='meta-live-1' and contact_eligibility='pending_policy'`, tenant).Scan(&candidateCount); err != nil {
		t.Fatal(err)
	}
	if err := pool.QueryRow(ctx, `select count(*) from integration.lead_ingress_raw where tenant_id=$1 and provider='meta_lead_ads' and provider_event_id='meta-rejected-1' and state='rejected' and normalization_error_code='UNAPPROVED_FORM_FIELD'`, tenant).Scan(&rejectedCount); err != nil {
		t.Fatal(err)
	}
	if err := pool.QueryRow(ctx, `select count(*) from platform.outbox_event where tenant_id=$1 and aggregate_id in ('meta_lead_ads:meta-live-1','meta_lead_ads:meta-rejected-1')`, tenant).Scan(&outboxCount); err != nil {
		t.Fatal(err)
	}
	if rawCount != 2 || candidateCount != 1 || rejectedCount != 1 || outboxCount != 2 {
		t.Fatalf("raw=%d candidate=%d rejected=%d outbox=%d", rawCount, candidateCount, rejectedCount, outboxCount)
	}
}

func postgresMetaArtifacts(t *testing.T, tenant, organization string) ([]byte, []byte, []byte) {
	t.Helper()
	rows := []map[string]any{
		{"id": "meta-live-1", "created_time": "2026-09-04T12:00:00+0000", "form_id": "123456789", "campaign_id": "campaign-1", "adset_id": "adset-1", "ad_id": "ad-1", "platform": "facebook", "field_data": []map[string]any{{"name": "email", "values": []string{"person@example.test"}}}},
		{"id": "meta-rejected-1", "created_time": "2026-09-04T12:01:00+0000", "form_id": "123456789", "campaign_id": "campaign-1", "adset_id": "adset-1", "ad_id": "ad-1", "platform": "facebook", "field_data": []map[string]any{{"name": "unapproved", "values": []string{"retained-only-in-raw"}}}},
	}
	provider := append(mustJSONIndent(t, map[string]any{"graph_api_version": "v26.0", "retrieval_mode": "LIVE", "rows": rows, "truncated_by_local_bound": false}), '\n')
	providerHash := testSHA(provider)
	outcomes := []map[string]any{
		{"status": "CANDIDATE", "normalization_codes": []string{}, "provider_lead_id_sha256": testSHA([]byte("meta-live-1")), "candidate": map[string]any{"tenant_id": tenant, "organization_id": organization, "provider": "meta_lead_ads", "provider_lead_id": "meta-live-1", "form_id": "123456789", "campaign_id": "campaign-1", "ad_group_id": "adset-1", "creative_id": "ad-1", "source_kind": "facebook", "submitted_at": "2026-09-04T12:00:00+0000", "is_test": false, "fields": []map[string]string{{"id": "email", "value": "person@example.test"}}, "contact_eligibility": "pending_policy"}},
		{"status": "REJECTED", "normalization_codes": []string{"UNAPPROVED_FORM_FIELD"}, "provider_lead_id_sha256": testSHA([]byte("meta-rejected-1")), "candidate": nil},
	}
	batch := append(mustJSONIndent(t, map[string]any{"schema": "elite-meta-lead-candidate-batch/v1", "provider": "meta_lead_ads", "source_commit": "788f363d15b1269ab5efb7cd00fb5e3b133cd99b", "tenant_id": tenant, "organization_id": organization, "provider_response_sha256": providerHash, "outcomes": outcomes, "automatic_business_write": false, "automatic_contact_eligibility": false}), '\n')
	receipt := append(mustJSONIndent(t, map[string]any{"schema": "elite-meta-lead-reconciliation-receipt/v1", "sdk_version": "26.0.1", "graph_api_version": "v26.0", "source_commit": "788f363d15b1269ab5efb7cd00fb5e3b133cd99b", "form_id_sha256": testSHA([]byte("123456789")), "provider_response_sha256": providerHash, "candidate_batch_sha256": testSHA(batch), "row_count": 2, "candidate_count": 1, "rejected_count": 1, "truncated_by_local_bound": false, "automatic_business_write": false, "automatic_contact_eligibility": false}), '\n')
	return provider, batch, receipt
}

func mustJSONIndent(t *testing.T, value any) []byte {
	t.Helper()
	encoded, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	return encoded
}

func testSHA(value []byte) string {
	const digits = "0123456789abcdef"
	return sha256Hex(value, digits)
}

func sha256Hex(value []byte, digits string) string {
	h := sha256.Sum256(value)
	encoded := make([]byte, len(h)*2)
	for i, b := range h {
		encoded[i*2], encoded[i*2+1] = digits[b>>4], digits[b&15]
	}
	return string(encoded)
}
````

### FILE: `cmd/meta-lead-import/main.go`
```yaml
block_id: "GO-META-LEAD-IMPORT:command:v1"
operation: CREATE
provenance: ADAPTED
source: "local enterprise command/pool patterns plus V226 artifact contract"
license: "LicenseRef-Workspace-Owner"
sha256: "446cf69e93fb0817d7523152f69d5e689f477dc4b0fe67bcde7c4b5799e2ad15"
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
	providerPath := flag.String("provider-response", "", "V226 provider-response.json")
	batchPath := flag.String("candidate-batch", "", "V226 lead-candidates.json")
	receiptPath := flag.String("retrieval-receipt", "", "V226 RETRIEVAL_RECEIPT.json")
	flag.Parse()
	if *providerPath == "" || *batchPath == "" || *receiptPath == "" || os.Getenv("DATABASE_URL") == "" {
		slog.Error("DATABASE_URL and all three V226 artifact paths are required")
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
	result, err := leadstream.ImportMetaEvidence(ctx, postgres.NewLeadIngress(pool), providerBytes, batchBytes, receiptBytes, time.Now().UTC())
	encoded, _ := json.Marshal(result)
	if err != nil {
		fmt.Fprintln(os.Stderr, string(encoded))
		fail("Meta lead import incomplete; replay is safe", err)
	}
	fmt.Println(string(encoded))
}

func readBounded(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	bytes, err := io.ReadAll(io.LimitReader(file, maxArtifactBytes+1))
	if err != nil {
		return nil, err
	}
	if int64(len(bytes)) > maxArtifactBytes {
		return nil, fmt.Errorf("artifact exceeds %d bytes", maxArtifactBytes)
	}
	if len(bytes) == 0 {
		return nil, fmt.Errorf("artifact is empty")
	}
	return bytes, nil
}

func fail(message string, err error) {
	slog.Error(message, "error", err)
	os.Exit(1)
}
````

## 6. Configuration surface

- `DATABASE_URL` sólo por variable de entorno.
- Paths explícitos a los tres artefactos V226; máximo 32 MiB por archivo.
- Tenant, organización, form, modo y campos provienen del expediente V226 firmado por hashes y se revalidan.

## 7. Dependency bill

| Dependency | Pin | License | Purpose |
|---|---:|---|---|
| Go | 1.26.7 | BSD-3-Clause | artifact validation/import command |
| PostgreSQL | 18.6 | PostgreSQL | durable ingress owner |
| pgx | 5.10.0 | MIT | PostgreSQL driver |
| Meta Business SDK artifact | 26.0.1 / locked wheel SHA-256 | Meta Platform license | upstream evidence produced by V226, not embedded here |

## 8. Apply order

1. Materializar V224, V225 y V226.
2. Materializar este pack sin colisiones.
3. Aplicar migraciones hasta 0045.
4. Ejecutar unit/integration, full Go, vet y build.
5. Generar `/test_leads` real con V226 e importar mediante el comando.
6. Reconciliar raw/candidato/rechazo/outbox y repetir para demostrar replay.
7. Sólo la política/mapping explícitos de V225 pueden crear el lead CRM.

## 9. Verification

```powershell
go test ./internal/leadstream -count=1
$env:TEST_DATABASE_URL = '<ephemeral PostgreSQL 18.6 URL with migrations through 0045>'
go test ./internal/platform/postgres -run TestMetaEvidenceImportDurabilityRejectionAndReplay -count=1
go test ./... -count=1
go vet ./...
go build ./...
```

### Production blockers

- Cuenta, permisos, Page/form y `/test_leads`/`/leads` reales.
- Import completo cuando V226 declare truncamiento.
- Scheduler/worker y archivo cifrado/retención según target.
- Política de consentimiento, mapping y promoción V225 con corpus real.
- Webhook live, conversación, outcome y postback de conversión siguen separados.

## 10. Reconstruction evidence

`reconstruction_evidence/GO_META_LEAD_EVIDENCE_IMPORT_2026-09-04_V227.md`
