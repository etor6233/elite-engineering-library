// AUTHORED bounded transport and durable-result contract around existing fulfillment.
package networkrole

import (
	"elite.local/enterprise/internal/approval"
	"elite.local/enterprise/internal/platform/identity"
	"encoding/json"
	"errors"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

var ErrInvalid = errors.New("invalid network command")
var ErrNotFound = errors.New("network result unavailable")
var ErrConflict = errors.New("network command conflict")

type Option struct {
 Value string `json:"value"`
 Label string `json:"label"`
}

// Existing organization_organization_code_check, surfaced before transport.
var organizationCode = regexp.MustCompile(`^[a-z][a-z0-9]*(-[a-z0-9]+)*$`)
var idPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:-]{0,127}$`)

func ID(s string) bool { return idPattern.MatchString(s) }
func Text(s string, max int) bool {
	if len(s) < 2 || len(s) > max || !utf8.ValidString(s) || strings.TrimSpace(s) != s {
		return false
	}
	for _, r := range s {
		if unicode.IsControl(r) {
			return false
		}
	}
	return true
}

type Command struct {
	CommandID           string `json:"command_id"`
	Action              string `json:"action"`
	ScopeOrganizationID string `json:"scope_organization_id"`
	EntityID            string `json:"entity_id"`
	Code                string `json:"code,omitempty"`
	DisplayName         string `json:"display_name,omitempty"`
	Type                string `json:"type,omitempty"`
	TerritoryCode       string `json:"territory_code,omitempty"`
	TermsVersion        string `json:"terms_version,omitempty"`
	StartsOn            string `json:"starts_on,omitempty"`
	EndsOn              string `json:"ends_on,omitempty"`
	Current             string `json:"current,omitempty"`
	Target              string `json:"target,omitempty"`
	Version             string `json:"version,omitempty"`
}

func Permission(action string) string {
	switch action {
	case "create-organization", "transition-organization":
		return "network:admin"
	case "create-agreement", "transition-agreement":
		return "franchise:write"
	}
	return ""
}
func (c Command) Valid() bool {
	if !ID(c.CommandID) || !ID(c.EntityID) || (c.ScopeOrganizationID != "" && !ID(c.ScopeOrganizationID)) {
		return false
	}
	switch c.Action {
	case "create-organization":
		types := map[string]bool{"enterprise": true, "franchisor": true, "franchisee": true, "factory": true, "warehouse": true, "store": true, "service_center": true}
		root := c.Type == "enterprise" || c.Type == "franchisor"
		if !types[c.Type] || (root && c.ScopeOrganizationID != "") || (!root && c.ScopeOrganizationID == "") {
			return false
		}
		return len(c.Code) <= 128 && organizationCode.MatchString(c.Code) && Text(c.DisplayName, 100) && ID(c.Type) && c.TerritoryCode == "" && c.TermsVersion == "" && c.StartsOn == "" && c.EndsOn == "" && c.Current == "" && c.Target == "" && c.Version == ""
	case "create-agreement":
		start, e := time.Parse("2006-01-02", c.StartsOn)
		if e != nil || start.Year() < 1 {
			return false
		}
		if c.EndsOn != "" {
			end, e := time.Parse("2006-01-02", c.EndsOn)
			if e != nil || !end.After(start) {
				return false
			}
		}
		return ID(c.ScopeOrganizationID) && ID(c.TerritoryCode) && ID(c.TermsVersion) && c.Code == "" && c.DisplayName == "" && c.Type == "" && c.Current == "" && c.Target == "" && c.Version == ""
	case "transition-organization", "transition-agreement":
		v, e := strconv.ParseInt(c.Version, 10, 64)
		return e == nil && v > 0 && v < math.MaxInt64 && strconv.FormatInt(v, 10) == c.Version && ID(c.ScopeOrganizationID) && ID(c.Current) && ID(c.Target) && c.Code == "" && c.DisplayName == "" && c.Type == "" && c.TerritoryCode == "" && c.TermsVersion == "" && c.StartsOn == "" && c.EndsOn == "" && (c.Action != "transition-organization" || c.ScopeOrganizationID == c.EntityID)
	}
	return false
}
func Authorized(p identity.Principal, action, scope string) bool {
	perm := Permission(action)
	if p.TenantID == "" || p.Subject == "" || perm == "" || !p.Allowed(perm) {
		return false
	}
	if scope == "" {
		return action == "create-organization" && p.Allowed("network:bootstrap")
	}
	return p.AllowedOrganization(scope)
}
func Canonical(v any) (json.RawMessage, string, error) {
	raw, e := json.Marshal(v)
	if e != nil {
		return nil, "", e
	}
	return approval.CanonicalPayload(raw)
}

type Entity struct {
	Kind                 string `json:"kind"`
	ID                   string `json:"id"`
	OrganizationID       string `json:"organization_id"`
	Version              string `json:"version"`
	State                string `json:"state"`
	Code                 string `json:"code,omitempty"`
	DisplayName          string `json:"display_name,omitempty"`
	Type                 string `json:"type,omitempty"`
	ParentOrganizationID string `json:"parent_organization_id,omitempty"`
	TerritoryCode        string `json:"territory_code,omitempty"`
	TermsVersion         string `json:"terms_version,omitempty"`
	StartsOn             string `json:"starts_on,omitempty"`
	EndsOn               string `json:"ends_on,omitempty"`
}
type Receipt struct {
	CommandID           string    `json:"command_id"`
	Action              string    `json:"action"`
	ScopeOrganizationID string    `json:"scope_organization_id"`
	Actor               string    `json:"actor"`
	RequestSHA256       string    `json:"request_sha256"`
	Entity              Entity    `json:"entity"`
	RecordedAt          time.Time `json:"recorded_at"`
	Replay              bool      `json:"replay"`
}
