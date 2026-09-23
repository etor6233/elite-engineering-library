// Package contactidentity defines the evidence-bound association between a
// provider channel identity and an existing CRM lead. It never infers identity
// or consent from message text.
package contactidentity

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"regexp"
	"strings"
	"time"
)

var (
	ErrInvalid  = errors.New("contactidentity: invalid command")
	ErrConflict = errors.New("contactidentity: version or replay conflict")
	ErrNotFound = errors.New("contactidentity: active binding not found")
	channelRE   = regexp.MustCompile(`^[a-z][a-z0-9_]{0,31}$`)
	hexRE       = regexp.MustCompile(`^[0-9a-f]{64}$`)
)

type State string

const (
	StateActive  State = "active"
	StateRevoked State = "revoked"
)

type Command struct {
	TenantID        string
	ChannelCode     string
	ExternalID      string
	LeadID          string
	SubjectID       string
	PIIAllowed      bool
	State           State
	PolicyVersion   string
	EvidenceSHA256  string
	EffectiveAt     time.Time
	ExpectedVersion int64
	RequestID       string
}

type Receipt struct {
	Version  int64
	State    State
	Replayed bool
}

func (c Command) Validate() error {
	if strings.TrimSpace(c.TenantID) == "" || !channelRE.MatchString(c.ChannelCode) || strings.TrimSpace(c.ExternalID) == "" || len(c.ExternalID) > 256 || strings.TrimSpace(c.LeadID) == "" || strings.TrimSpace(c.SubjectID) == "" || len(c.SubjectID) > 256 || strings.TrimSpace(c.PolicyVersion) == "" || len(c.PolicyVersion) > 128 || !hexRE.MatchString(c.EvidenceSHA256) || c.EffectiveAt.IsZero() || c.ExpectedVersion < 0 || strings.TrimSpace(c.RequestID) == "" || len(c.RequestID) > 128 {
		return ErrInvalid
	}
	if c.State != StateActive && c.State != StateRevoked {
		return ErrInvalid
	}
	if c.State == StateRevoked && c.PIIAllowed {
		return ErrInvalid
	}
	return nil
}

func ExternalDigest(secret []byte, tenantID, channelCode, externalID string) (string, error) {
	if len(secret) < 32 || strings.TrimSpace(tenantID) == "" || !channelRE.MatchString(channelCode) || strings.TrimSpace(externalID) == "" || len(externalID) > 256 {
		return "", ErrInvalid
	}
	mac := hmac.New(sha256.New, secret)
	_, _ = mac.Write([]byte(tenantID))
	_, _ = mac.Write([]byte{0})
	_, _ = mac.Write([]byte(channelCode))
	_, _ = mac.Write([]byte{0})
	_, _ = mac.Write([]byte(externalID))
	return hex.EncodeToString(mac.Sum(nil)), nil
}

func (c Command) RequestSHA256(secret []byte) (string, string, error) {
	if err := c.Validate(); err != nil {
		return "", "", err
	}
	digest, err := ExternalDigest(secret, c.TenantID, c.ChannelCode, c.ExternalID)
	if err != nil {
		return "", "", err
	}
	canonical := struct {
		TenantID, ChannelCode, ExternalDigest, LeadID, SubjectID, PolicyVersion, EvidenceSHA256, EffectiveAt, RequestID string
		PIIAllowed                                                                                                      bool
		State                                                                                                           State
		ExpectedVersion                                                                                                 int64
	}{c.TenantID, c.ChannelCode, digest, c.LeadID, c.SubjectID, c.PolicyVersion, c.EvidenceSHA256, c.EffectiveAt.UTC().Format(time.RFC3339Nano), c.RequestID, c.PIIAllowed, c.State, c.ExpectedVersion}
	b, err := json.Marshal(canonical)
	if err != nil {
		return "", "", err
	}
	sum := sha256.Sum256(b)
	return digest, hex.EncodeToString(sum[:]), nil
}
