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
