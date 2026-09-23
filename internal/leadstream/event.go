// Package leadstream defines the provider-neutral, fail-closed boundary for
// lead events. Provider payloads remain byte-exact while normalized candidates
// are versioned separately; receiving a payload never grants permission to
// contact a person or to write a business lead automatically.
package leadstream

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

const MaxPayloadBytes int64 = 1 << 20

var (
	ErrInvalidEvent       = errors.New("leadstream: invalid event")
	ErrInvalidCredential  = errors.New("leadstream: invalid provider credential")
	ErrInvalidPayload     = errors.New("leadstream: invalid provider payload")
	ErrDivergentDuplicate = errors.New("leadstream: divergent duplicate")
	providerRe            = regexp.MustCompile(`^[a-z][a-z0-9_]{0,31}$`)
)

type RawEvent struct {
	TenantID         string
	OrganizationID   string
	Provider         string
	ProviderEventID  string
	EventType        string
	Source           string
	SchemaVersion    string
	OccurredAt       time.Time
	ReceivedAt       time.Time
	Payload          []byte // durable representation; provider credentials must be removed first
	SourceSHA256     string // hash of the exact provider body before redaction
	StoredSHA256     string // hash of Payload
	RedactionProfile string
}

func NewRawEvent(tenantID, organizationID, provider, eventID, eventType, source, schemaVersion string, occurredAt, receivedAt time.Time, payload []byte) (RawEvent, error) {
	return NewRedactedRawEvent(tenantID, organizationID, provider, eventID, eventType, source, schemaVersion, occurredAt, receivedAt, payload, payload, "none:v1")
}

// NewRedactedRawEvent binds a safe durable representation to the hash of the
// exact source body. It never retains sourcePayload and therefore prevents a
// provider credential removed by the caller from being persisted accidentally.
func NewRedactedRawEvent(tenantID, organizationID, provider, eventID, eventType, source, schemaVersion string, occurredAt, receivedAt time.Time, sourcePayload, storedPayload []byte, redactionProfile string) (RawEvent, error) {
	sourceHash := sha256.Sum256(sourcePayload)
	storedHash := sha256.Sum256(storedPayload)
	e := RawEvent{
		TenantID: tenantID, OrganizationID: organizationID, Provider: provider,
		ProviderEventID: eventID, EventType: eventType, Source: source,
		SchemaVersion: schemaVersion, OccurredAt: occurredAt.UTC(), ReceivedAt: receivedAt.UTC(),
		Payload: append([]byte(nil), storedPayload...), SourceSHA256: hex.EncodeToString(sourceHash[:]), StoredSHA256: hex.EncodeToString(storedHash[:]),
		RedactionProfile: redactionProfile,
	}
	return e, e.Validate()
}

func (e RawEvent) Validate() error {
	if strings.TrimSpace(e.TenantID) == "" || strings.TrimSpace(e.OrganizationID) == "" {
		return fmt.Errorf("%w: tenant and organization required", ErrInvalidEvent)
	}
	if !providerRe.MatchString(e.Provider) || strings.TrimSpace(e.ProviderEventID) == "" || len(e.ProviderEventID) > 256 {
		return fmt.Errorf("%w: provider identity", ErrInvalidEvent)
	}
	if strings.TrimSpace(e.EventType) == "" || strings.TrimSpace(e.Source) == "" || strings.TrimSpace(e.SchemaVersion) == "" {
		return fmt.Errorf("%w: event context", ErrInvalidEvent)
	}
	if e.ReceivedAt.IsZero() || e.OccurredAt.IsZero() || len(e.Payload) == 0 || int64(len(e.Payload)) > MaxPayloadBytes {
		return fmt.Errorf("%w: time or payload", ErrInvalidEvent)
	}
	if strings.TrimSpace(e.RedactionProfile) == "" || len(e.RedactionProfile) > 128 {
		return fmt.Errorf("%w: redaction profile", ErrInvalidEvent)
	}
	h := sha256.Sum256(e.Payload)
	if e.StoredSHA256 != hex.EncodeToString(h[:]) || !validSHA256(e.SourceSHA256) {
		return fmt.Errorf("%w: payload hash", ErrInvalidEvent)
	}
	return nil
}

func validSHA256(value string) bool {
	if len(value) != sha256.Size*2 {
		return false
	}
	_, err := hex.DecodeString(value)
	return err == nil
}

type Field struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

type LeadCandidate struct {
	TenantID           string    `json:"tenant_id"`
	OrganizationID     string    `json:"organization_id"`
	Provider           string    `json:"provider"`
	ProviderLeadID     string    `json:"provider_lead_id"`
	FormID             string    `json:"form_id,omitempty"`
	CampaignID         string    `json:"campaign_id,omitempty"`
	AdGroupID          string    `json:"ad_group_id,omitempty"`
	CreativeID         string    `json:"creative_id,omitempty"`
	AssetGroupID       string    `json:"asset_group_id,omitempty"`
	ClickID            string    `json:"click_id,omitempty"`
	LeadStage          string    `json:"lead_stage,omitempty"`
	SourceKind         string    `json:"source_kind,omitempty"`
	SubmittedAt        time.Time `json:"submitted_at"`
	IsTest             bool      `json:"is_test"`
	Fields             []Field   `json:"fields"`
	ContactEligibility string    `json:"contact_eligibility"`
}

func (c LeadCandidate) Validate() error {
	if c.TenantID == "" || c.OrganizationID == "" || !providerRe.MatchString(c.Provider) || c.ProviderLeadID == "" || c.SubmittedAt.IsZero() {
		return fmt.Errorf("%w: candidate identity", ErrInvalidPayload)
	}
	if c.ContactEligibility != "pending_policy" {
		return fmt.Errorf("%w: ingress cannot grant contact eligibility", ErrInvalidPayload)
	}
	seen := map[string]struct{}{}
	for _, f := range c.Fields {
		if strings.TrimSpace(f.ID) == "" || strings.TrimSpace(f.Value) == "" {
			return fmt.Errorf("%w: empty field", ErrInvalidPayload)
		}
		if _, ok := seen[f.ID]; ok {
			return fmt.Errorf("%w: duplicate field %s", ErrInvalidPayload, f.ID)
		}
		seen[f.ID] = struct{}{}
	}
	return nil
}

func StableUUID(parts ...string) string {
	sum := sha256.Sum256([]byte(strings.Join(parts, "\x00")))
	b := sum[:16]
	b[6] = (b[6] & 0x0f) | 0x50
	b[8] = (b[8] & 0x3f) | 0x80
	hexv := hex.EncodeToString(b)
	return hexv[0:8] + "-" + hexv[8:12] + "-" + hexv[12:16] + "-" + hexv[16:20] + "-" + hexv[20:32]
}
