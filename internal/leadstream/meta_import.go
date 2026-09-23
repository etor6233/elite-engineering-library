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
