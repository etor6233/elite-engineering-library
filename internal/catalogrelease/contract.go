// AUTHORED publication binding contracts for existing catalog/price/approval
// owners. No legal, tax, pricing or ranking rules are inferred.
package catalogrelease

import (
	"elite.local/enterprise/internal/approval"
	"elite.local/enterprise/internal/electromobility"
	"encoding/json"
	"errors"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

var ErrInvalid = errors.New("invalid catalog publication request")
var ErrNotFound = errors.New("catalog publication outside scope")
var ErrConflict = errors.New("catalog publication conflict")
var identifier = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:-]{0,63}$`)
var digestPattern = regexp.MustCompile(`^[0-9a-f]{64}$`)

const PolicyCode = "catalog-reference-v1"

func ValidID(s string) bool  { return identifier.MatchString(s) }
func ValidSHA(s string) bool { return digestPattern.MatchString(s) }
func ValidText(s string, n int) bool {
	if s == "" || len(s) > n || !utf8.ValidString(s) || strings.TrimSpace(s) != s {
		return false
	}
	for _, r := range s {
		if unicode.IsControl(r) {
			return false
		}
	}
	return true
}
func Canonical(v any) (json.RawMessage, string, error) {
	raw, err := json.Marshal(v)
	if err != nil {
		return nil, "", err
	}
	return approval.CanonicalPayload(raw)
}

type Profile struct {
	Schema         string `json:"schema"`
	TenantID       string `json:"tenant_id"`
	OrganizationID string `json:"organization_id"`
	Market         string `json:"market"`
	Currency       string `json:"currency"`
	Origin         string `json:"origin"`
	PolicyCode     string `json:"policy_code"`
}

func (p Profile) Valid() bool {
	if p.Schema != "elite-catalog-publication/v1" || !regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`).MatchString(p.TenantID) || !ValidID(p.OrganizationID) || !regexp.MustCompile(`^[A-Z]{2}$`).MatchString(p.Market) || !regexp.MustCompile(`^[A-Z]{3}$`).MatchString(p.Currency) || p.PolicyCode != PolicyCode {
		return false
	}
	u, e := url.Parse(p.Origin)
	return e == nil && u.Scheme == "https" && u.Host != "" && u.User == nil && u.RawQuery == "" && u.Fragment == "" && u.Path == "" && u.Opaque == "" && len(p.Origin) <= 512
}

type DraftRequest struct {
	CommandID   string           `json:"command_id"`
	PriceBookID string           `json:"price_book_id"`
	Models      []ModelReference `json:"models"`
}
type ModelReference struct {
	ModelID string `json:"model_id"`
	MediaID string `json:"media_id"`
}

func (r DraftRequest) Valid() bool {
	if !ValidID(r.CommandID) || !ValidID(r.PriceBookID) || len(r.Models) == 0 || len(r.Models) > 32 {
		return false
	}
	seen := map[string]bool{}
	for _, m := range r.Models {
		if !ValidID(m.ModelID) || !ValidID(m.MediaID) || seen[m.ModelID] {
			return false
		}
		seen[m.ModelID] = true
	}
	return true
}

type ReviewRequest struct {
	CommandID      string `json:"command_id"`
	DraftID        string `json:"draft_id"`
	Stage          string `json:"stage"`
	SnapshotSHA256 string `json:"snapshot_sha256"`
	Approved       bool   `json:"approved"`
	Reason         string `json:"reason"`
	EvidenceSHA256 string `json:"evidence_sha256"`
}

func (r ReviewRequest) Valid() bool {
	return ValidID(r.CommandID) && ValidID(r.DraftID) && ValidSHA(r.SnapshotSHA256) && ValidSHA(r.EvidenceSHA256) && ValidText(r.Reason, 2000) && map[string]bool{"legal": true, "technical": true, "media": true, "publication": true}[r.Stage]
}

type PublishRequest struct {
	CommandID          string `json:"command_id"`
	DraftID            string `json:"draft_id"`
	SnapshotSHA256     string `json:"snapshot_sha256"`
	ExpectedGeneration int64  `json:"expected_generation,string"`
	Reason             string `json:"reason"`
}

func (r PublishRequest) Valid() bool {
	return ValidID(r.CommandID) && ValidID(r.DraftID) && ValidSHA(r.SnapshotSHA256) && r.ExpectedGeneration >= 0 && ValidText(r.Reason, 2000)
}

type Receipt struct {
	SourcePriceBookID    string   `json:"source_price_book_id,omitempty"`
	SourceVariantIDs     []string `json:"source_variant_ids,omitempty"`
	EffectivePriceBookID string   `json:"effective_price_book_id,omitempty"`
	CommandID            string   `json:"command_id"`
	Actor                string   `json:"actor"`
	Kind                 string   `json:"kind"`
	ResourceID           string   `json:"resource_id"`
	RequestSHA256        string   `json:"request_sha256"`
	SnapshotSHA256       string   `json:"snapshot_sha256,omitempty"`
	Generation           int64    `json:"generation,string"`
	Replay               bool     `json:"replay"`
}
type Media struct {
	ID             string `json:"id"`
	OriginalSHA256 string `json:"original_sha256"`
	SHA256         string `json:"sha256"`
	Width          int    `json:"width"`
	Height         int    `json:"height"`
}
type Variant struct {
	ID                   string          `json:"id"`
	ModelID              string          `json:"model_id"`
	Code                 string          `json:"code"`
	DisplayName          string          `json:"display_name"`
	BatterySpecification json.RawMessage `json:"battery_specification"`
	HomologationState    string          `json:"homologation_state"`
	AmountMinorUnits     int64           `json:"amount_minor_units,string"`
	TaxMode              string          `json:"tax_mode"`
}
type Model struct {
	electromobility.Model
	Media        Media  `json:"media"`
	CanonicalURL string `json:"canonical_url"`
}
type Snapshot struct {
	TenantCode  string     `json:"tenant_code"`
	Schema      string     `json:"schema"`
	Profile     Profile    `json:"profile"`
	PriceBookID string     `json:"price_book_id"`
	ValidFrom   time.Time  `json:"valid_from"`
	ValidUntil  *time.Time `json:"valid_until"`
	Models      []Model    `json:"models"`
	Variants    []Variant  `json:"variants"`
}
type Draft struct {
	ID       string            `json:"id"`
	Maker    string            `json:"maker"`
	SHA256   string            `json:"sha256"`
	Snapshot json.RawMessage   `json:"snapshot"`
	Reviews  map[string]string `json:"reviews"`
}
type Publication struct {
	EffectivePriceBookID string          `json:"effective_price_book_id,omitempty"`
	Generation           int64           `json:"generation,string"`
	DraftID              string          `json:"draft_id"`
	SHA256               string          `json:"sha256"`
	Snapshot             json.RawMessage `json:"snapshot"`
}

type SearchHit struct {
	ModelID      string `json:"model_id"`
	Title        string `json:"title"`
	URL          string `json:"url"`
	Generation   int64  `json:"generation,string"`
	SourceSHA256 string `json:"source_sha256"`
}

// PublicDocument omits private stream ownership and review configuration.
type PublicDocument struct {
	TenantCode           string    `json:"tenant_code"`
	EffectivePriceBookID string    `json:"effective_price_book_id,omitempty"`
	Generation           int64     `json:"generation,string"`
	SourceSHA256         string    `json:"source_sha256"`
	Market               string    `json:"market"`
	Currency             string    `json:"currency"`
	Models               []Model   `json:"models"`
	Variants             []Variant `json:"variants"`
}

func ProjectPublic(p Publication) (PublicDocument, error) {
	var s Snapshot
	if err := json.Unmarshal(p.Snapshot, &s); err != nil {
		return PublicDocument{}, err
	}
	return PublicDocument{TenantCode: s.TenantCode, EffectivePriceBookID: p.EffectivePriceBookID, Generation: p.Generation, SourceSHA256: p.SHA256, Market: s.Profile.Market, Currency: s.Profile.Currency, Models: s.Models, Variants: s.Variants}, nil
}
