// Package socialbridge is AUTHORED integration glue. It binds an explicit
// deployment profile and reviewed bytes to existing approval/jobs/fence owners.
// Meta owns provider behavior; no marketing policy or audience logic is inferred.
package socialbridge

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"os"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"elite.local/enterprise/internal/channels"
	"elite.local/enterprise/internal/platform/identity"
)

var ErrBinding = errors.New("social: invalid or unauthorized binding")
var ErrUnknown = errors.New("social: remote outcome unresolved; no repeated write permitted")
var hashRE = regexp.MustCompile(`^[0-9a-f]{64}$`)
var idRE = regexp.MustCompile(`^[A-Za-z0-9_-]{1,128}$`)
var pageRE = regexp.MustCompile(`^[1-9][0-9]{0,31}$`)

type Config struct {
	Schema         string `json:"schema"`
	TenantID       string `json:"tenant_id"`
	OrganizationID string `json:"organization_id"`
	PageID         string `json:"page_id"`
	Queue          string `json:"queue"`
	Publish        bool   `json:"publish"`
	Revoke         bool   `json:"revoke"`
	Review         string `json:"review"`
	LeaseSeconds   int    `json:"lease_seconds"`
	RetrySeconds   int    `json:"retry_seconds"`
	PollSeconds    int    `json:"poll_seconds"`
	MaxAttempts    int    `json:"max_attempts"`
}
type Profile struct {
	config Config
	hash   string
}

func Hash(raw []byte) string { v := sha256.Sum256(raw); return hex.EncodeToString(v[:]) }
func Decode(raw []byte, into any) error {
	if len(raw) == 0 || len(raw) > 32768 || !utf8.Valid(raw) {
		return ErrBinding
	}
	// Duplicate/case-alias rejection is protocol glue, not a business policy.
	d := json.NewDecoder(bytes.NewReader(raw))
	if unique(d, 0) != nil {
		return ErrBinding
	}
	if _, e := d.Token(); e != io.EOF {
		return ErrBinding
	}
	d = json.NewDecoder(bytes.NewReader(raw))
	d.DisallowUnknownFields()
	if d.Decode(into) != nil {
		return ErrBinding
	}
	return nil
}
func unique(d *json.Decoder, depth int) error {
	if depth > 16 {
		return ErrBinding
	}
	t, e := d.Token()
	if e != nil {
		return e
	}
	v, ok := t.(json.Delim)
	if !ok {
		return nil
	}
	switch v {
	case '{':
		seen := map[string]bool{}
		for d.More() {
			k, e := d.Token()
			if e != nil {
				return e
			}
			s, ok := k.(string)
			if !ok {
				return ErrBinding
			}
			s = strings.ToLower(s)
			if seen[s] {
				return ErrBinding
			}
			seen[s] = true
			if unique(d, depth+1) != nil {
				return ErrBinding
			}
		}
	case '[':
		for d.More() {
			if unique(d, depth+1) != nil {
				return ErrBinding
			}
		}
	default:
		return ErrBinding
	}
	_, e = d.Token()
	return e
}
func Load(raw []byte, expected string) (*Profile, error) {
	var c Config
	if !hashRE.MatchString(expected) || Hash(raw) != expected || Decode(raw, &c) != nil {
		return nil, ErrBinding
	}
	if c.Schema != "elite.meta-page-publishing.v1" || !idRE.MatchString(c.TenantID) || !idRE.MatchString(c.OrganizationID) || !pageRE.MatchString(c.PageID) || !idRE.MatchString(c.Queue) || c.Review != "one_distinct_human" || c.LeaseSeconds < 60 || c.LeaseSeconds > 600 || c.RetrySeconds < 1 || c.RetrySeconds > 3600 || c.PollSeconds < 1 || c.PollSeconds > 60 || c.MaxAttempts < 2 || c.MaxAttempts > 100 {
		return nil, ErrBinding
	}
	return &Profile{c, expected}, nil
}
func LoadFile(path, hash string) (*Profile, error) {
	f, e := os.Open(path)
	if e != nil {
		return nil, ErrBinding
	}
	defer f.Close()
	b, e := io.ReadAll(io.LimitReader(f, 32769))
	if e != nil {
		return nil, ErrBinding
	}
	return Load(b, hash)
}
func (p *Profile) Config() Config { return p.config }
func (p *Profile) SHA256() string { return p.hash }
func (p *Profile) Authorize(a identity.Principal, permission string) bool {
	return p != nil && a.Subject != "" && len(a.Subject) <= 128 && a.TenantID == p.config.TenantID && a.AllowedOrganization(p.config.OrganizationID) && a.Allowed(permission)
}

type Intent struct {
	TenantID      string `json:"tenant_id"`
	PageID        string `json:"page_id"`
	ProfileSHA256 string `json:"profile_sha256"`
	ApprovalID    string `json:"approval_id"`
	DeliveryKey   string `json:"delivery_key"`
	Operation     string `json:"operation"`
	Message       string `json:"message"`
	ContentSHA256 string `json:"content_sha256"`
}
type Request struct {
	Intent             Intent    `json:"intent"`
	ScheduledAt        time.Time `json:"scheduled_at"`
	ExpiresAt          time.Time `json:"expires_at"`
	OriginalApprovalID string    `json:"original_approval_id"`
	ProviderReference  string    `json:"provider_reference"`
}

func (p *Profile) Validate(r Request) error {
	i := r.Intent
	c := p.config
	if i.TenantID != c.TenantID || i.PageID != c.PageID || i.ProfileSHA256 != p.hash || !idRE.MatchString(i.ApprovalID) || i.DeliveryKey != DeliveryKey(i) || !utf8.ValidString(i.Message) || len(i.Message) > 16384 || strings.TrimSpace(i.Message) == "" || Hash([]byte(i.Message)) != i.ContentSHA256 || r.ScheduledAt.IsZero() || !r.ExpiresAt.After(r.ScheduledAt) {
		return ErrBinding
	}
	if i.Operation == "publish" {
		if !c.Publish || r.OriginalApprovalID != "" || r.ProviderReference != "" {
			return ErrBinding
		}
	} else if i.Operation == "revoke" {
		if !c.Revoke || !idRE.MatchString(r.OriginalApprovalID) || !ValidReference(c.PageID, r.ProviderReference) {
			return ErrBinding
		}
	} else {
		return ErrBinding
	}
	return nil
}
func DeliveryKey(i Intent) string {
	return "social_" + Hash([]byte(i.TenantID+"\x00"+i.PageID+"\x00"+i.ApprovalID+"\x00"+i.Operation))
}
func ValidReference(page, ref string) bool {
	return regexp.MustCompile(`^` + regexp.QuoteMeta(page) + `_[1-9][0-9]{0,63}$`).MatchString(ref)
}
func (r Request) Hash() string { b, _ := json.Marshal(r); return Hash(b) }
func (r Request) Message() channels.Message {
	return channels.Message{ChannelCode: "facebook_pages", TenantID: r.Intent.TenantID, ExternalID: r.Intent.PageID, ThreadID: r.Hash(), DeliveryKey: r.Intent.DeliveryKey, Direction: channels.DirectionOut, Text: r.Intent.Message}
}

type Receipt struct {
	TenantID          string `json:"tenant_id"`
	PageID            string `json:"page_id"`
	ProfileSHA256     string `json:"profile_sha256"`
	ApprovalID        string `json:"approval_id"`
	DeliveryKey       string `json:"delivery_key"`
	ProviderReference string `json:"provider_reference"`
	State             string `json:"state"`
	ContentSHA256     string `json:"content_sha256"`
	EvidenceSHA256    string `json:"evidence_sha256"`
}

func (r Receipt) Validate(i Intent) error {
	want := "published"
	if i.Operation == "revoke" {
		want = "revoked"
	}
	if r.TenantID != i.TenantID || r.PageID != i.PageID || r.ProfileSHA256 != i.ProfileSHA256 || r.ApprovalID != i.ApprovalID || r.DeliveryKey != i.DeliveryKey || r.ContentSHA256 != i.ContentSHA256 || r.State != want || !ValidReference(i.PageID, r.ProviderReference) || !hashRE.MatchString(r.EvidenceSHA256) {
		return ErrBinding
	}
	return nil
}

type Result struct {
	Receipt           *Receipt `json:"receipt,omitempty"`
	Unknown           bool     `json:"unknown,omitempty"`
	ProviderReference string   `json:"provider_reference,omitempty"`
	Code              string   `json:"code,omitempty"`
}
type Adapter interface {
	Execute(context.Context, Request, bool) (Result, error)
}
