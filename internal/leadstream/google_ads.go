package leadstream

import (
	"bytes"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"
)

type googleLeadColumn struct {
	ColumnName  string `json:"column_name"`
	StringValue string `json:"string_value"`
	ColumnID    string `json:"column_id"`
}

type googleWebhookLead struct {
	LeadID         string             `json:"lead_id"`
	UserData       []googleLeadColumn `json:"user_column_data"`
	APIVersion     string             `json:"api_version"`
	FormID         json.Number        `json:"form_id"`
	CampaignID     json.Number        `json:"campaign_id"`
	GoogleKey      string             `json:"google_key"`
	IsTest         bool               `json:"is_test"`
	GCLID          string             `json:"gcl_id"`
	AdGroupID      json.Number        `json:"adgroup_id"`
	CreativeID     json.Number        `json:"creative_id"`
	AssetGroupID   json.Number        `json:"asset_group_id"`
	LeadStage      string             `json:"lead_stage"`
	LeadSubmitTime string             `json:"lead_submit_time"`
	LeadSource     string             `json:"lead_source"`
}

type GoogleDecoded struct {
	Raw               RawEvent
	Candidate         *LeadCandidate
	NormalizationCode string
}

// DecodeGoogleWebhook implements the public Google Lead Form webhook contract.
// Unknown JSON fields are deliberately ignored for forward compatibility.
// Once the provider key is authenticated, a structurally invalid lead still
// returns a RawEvent so the durable store can retain rejected evidence.
func DecodeGoogleWebhook(payload []byte, expectedKey, tenantID, organizationID string, receivedAt time.Time) (GoogleDecoded, error) {
	if len(payload) == 0 || int64(len(payload)) > MaxPayloadBytes {
		return GoogleDecoded{}, fmt.Errorf("%w: payload size", ErrInvalidPayload)
	}
	dec := json.NewDecoder(bytes.NewReader(payload))
	dec.UseNumber()
	var in googleWebhookLead
	if err := dec.Decode(&in); err != nil {
		return GoogleDecoded{}, fmt.Errorf("%w: json", ErrInvalidPayload)
	}
	if err := ensureEOF(dec); err != nil {
		return GoogleDecoded{}, err
	}
	if !secretEqual(in.GoogleKey, expectedKey) {
		return GoogleDecoded{}, ErrInvalidCredential
	}
	storedPayload, err := redactGoogleCredential(payload)
	if err != nil {
		return GoogleDecoded{}, err
	}

	eventID := strings.TrimSpace(in.LeadID)
	if eventID == "" {
		eventID = "sha256:" + payloadHash(payload)
	}
	schema := strings.TrimSpace(in.APIVersion)
	if schema == "" {
		schema = "unspecified"
	}
	occurred := receivedAt.UTC()
	code := ""
	if in.LeadSubmitTime != "" {
		parsed, err := time.Parse(time.RFC3339, in.LeadSubmitTime)
		if err != nil {
			code = "INVALID_SUBMIT_TIME"
		} else {
			occurred = parsed.UTC()
		}
	}
	raw, err := NewRedactedRawEvent(tenantID, organizationID, "google_ads", eventID, "google.ads.lead.v"+schema, "https://googleads.googleapis.com/lead-form", schema, occurred, receivedAt, payload, storedPayload, "google_ads:remove-google_key:v1")
	if err != nil {
		return GoogleDecoded{}, err
	}
	if strings.TrimSpace(in.LeadID) == "" {
		return GoogleDecoded{Raw: raw, NormalizationCode: "MISSING_LEAD_ID"}, nil
	}
	fields := make([]Field, 0, len(in.UserData))
	for _, f := range in.UserData {
		fields = append(fields, Field{ID: strings.TrimSpace(f.ColumnID), Value: strings.TrimSpace(f.StringValue)})
	}
	candidate := LeadCandidate{
		TenantID: tenantID, OrganizationID: organizationID, Provider: "google_ads",
		ProviderLeadID: in.LeadID, FormID: in.FormID.String(), CampaignID: in.CampaignID.String(),
		AdGroupID: in.AdGroupID.String(), CreativeID: in.CreativeID.String(), AssetGroupID: in.AssetGroupID.String(),
		ClickID: in.GCLID, LeadStage: in.LeadStage, SourceKind: in.LeadSource,
		SubmittedAt: occurred, IsTest: in.IsTest, Fields: fields, ContactEligibility: "pending_policy",
	}
	if code == "" {
		if err := candidate.Validate(); err != nil {
			code = "INVALID_NORMALIZED_LEAD"
		}
	}
	if code != "" {
		return GoogleDecoded{Raw: raw, NormalizationCode: code}, nil
	}
	return GoogleDecoded{Raw: raw, Candidate: &candidate}, nil
}

func redactGoogleCredential(payload []byte) ([]byte, error) {
	dec := json.NewDecoder(bytes.NewReader(payload))
	dec.UseNumber()
	var object map[string]json.RawMessage
	if err := dec.Decode(&object); err != nil {
		return nil, fmt.Errorf("%w: redaction json", ErrInvalidPayload)
	}
	if err := ensureEOF(dec); err != nil {
		return nil, err
	}
	delete(object, "google_key")
	redacted, err := json.Marshal(object)
	if err != nil {
		return nil, fmt.Errorf("%w: redaction", ErrInvalidPayload)
	}
	return redacted, nil
}

func ensureEOF(dec *json.Decoder) error {
	var extra any
	if err := dec.Decode(&extra); err == io.EOF {
		return nil
	}
	return fmt.Errorf("%w: trailing json", ErrInvalidPayload)
}

func secretEqual(got, want string) bool {
	if got == "" || want == "" || len(got) != len(want) {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(got), []byte(want)) == 1
}

func payloadHash(payload []byte) string {
	sum := sha256.Sum256(payload)
	return hex.EncodeToString(sum[:])
}
