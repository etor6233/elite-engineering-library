// AUTHORED bounded persistence/authorization contract; article rules live in helpcenter.
package helpcms

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

var ErrInvalid = errors.New("invalid help CMS command")
var ErrNotFound = errors.New("help article unavailable")
var ErrConflict = errors.New("help CMS command conflict")
var idPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:-]{0,127}$`)
var categoryPattern = regexp.MustCompile(`^[a-z][a-z0-9_-]{0,63}$`)

const MaxBodyBytes = 16384

func ID(s string) bool { return idPattern.MatchString(s) }
func Text(s string, max int, multiline bool) bool {
	if strings.TrimSpace(s) == "" || len(s) > max || !utf8.ValidString(s) {
		return false
	}
	for _, r := range s {
		if unicode.IsControl(r) && (!multiline || r != '\n' && r != '\r' && r != '\t') {
			return false
		}
	}
	return true
}

type Command struct {
	CommandID      string `json:"command_id"`
	Action         string `json:"action"`
	ArticleID      string `json:"article_id"`
	OrganizationID string `json:"organization_id"`
	Version        string `json:"version,omitempty"`
	Locale         string `json:"locale,omitempty"`
	Category       string `json:"category,omitempty"`
	Title          string `json:"title,omitempty"`
	Body           string `json:"body,omitempty"`
}

func (c Command) Valid() bool {
	raw, e := json.Marshal(c)
	if e != nil || len(raw) > 24576 {
		return false
	}
	if !ID(c.CommandID) || !ID(c.ArticleID) || !ID(c.OrganizationID) {
		return false
	}
	switch c.Action {
	case "create":
		return c.Version == "" && (c.Locale == "es" || c.Locale == "en") && categoryPattern.MatchString(c.Category) && Text(c.Title, 200, false) && Text(c.Body, MaxBodyBytes, true)
	case "update", "publish", "archive":
		v, e := strconv.ParseInt(c.Version, 10, 64)
		if e != nil || v < 1 || v >= math.MaxInt64 || strconv.FormatInt(v, 10) != c.Version || c.Locale != "" || c.Category != "" {
			return false
		}
		if c.Action == "update" {
			return Text(c.Title, 200, false) && Text(c.Body, MaxBodyBytes, true)
		}
		return c.Title == "" && c.Body == ""
	}
	return false
}
func Permission(action string) string {
	switch action {
	case "create", "update":
		return "help:write"
	case "publish", "archive":
		return "help:publish"
	}
	return ""
}
func Authorized(p identity.Principal, permission, org string) bool {
	return permission != "" && p.TenantID != "" && p.Subject != "" && ID(org) && p.Allowed(permission) && p.AllowedOrganization(org)
}
func Editor(p identity.Principal, org string) bool {
	return Authorized(p, "help:write", org) || Authorized(p, "help:publish", org)
}
func Reader(p identity.Principal, org string) bool {
	return Editor(p, org) || Authorized(p, "help:read", org)
}
func Canonical(v any) (json.RawMessage, string, error) {
	raw, e := json.Marshal(v)
	if e != nil {
		return nil, "", ErrInvalid
	}
	body, hash, e := approval.CanonicalPayload(raw)
	if e != nil {
		return nil, "", ErrInvalid
	}
	return body, hash, nil
}

type Article struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organization_id"`
	Locale         string    `json:"locale"`
	Category       string    `json:"category"`
	Title          string    `json:"title"`
	Body           string    `json:"body"`
	State          string    `json:"state"`
	Version        string    `json:"version"`
	UpdatedAt      time.Time `json:"updated_at"`
}
type Receipt struct {
	CommandID     string  `json:"command_id"`
	Action        string  `json:"action"`
	Actor         string  `json:"actor"`
	RequestSHA256 string  `json:"request_sha256"`
	ArticleSHA256 string  `json:"article_sha256"`
	Article       Article `json:"article"`
	Replay        bool    `json:"replay"`
}
type Summary struct {
	ID       string `json:"id"`
	Locale   string `json:"locale"`
	Category string `json:"category"`
	Title    string `json:"title"`
	State    string `json:"state"`
	Version  string `json:"version"`
}
type Page struct {
	Items []Summary `json:"items"`
	Next  string    `json:"next,omitempty"`
}
