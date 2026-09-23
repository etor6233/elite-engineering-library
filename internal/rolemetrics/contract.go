// AUTHORED read-model contract. Source domains own records, money and lifecycle.
package rolemetrics

import (
	"elite.local/enterprise/internal/platform/identity"
	"errors"
	"regexp"
	"time"
)

var ErrInvalid = errors.New("invalid metric request")
var ErrUnavailable = errors.New("metric unavailable")
var ErrCardinality = errors.New("metric cardinality exceeds reference limit")
var idPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:-]{0,127}$`)

func ID(v string) bool { return idPattern.MatchString(v) }

type Row struct {
	State      string `json:"state"`
	Currency   string `json:"currency,omitempty"`
	Count      string `json:"count"`
	TotalMinor string `json:"total_minor_units,omitempty"`
}
type Snapshot struct {
	Kind           string    `json:"kind"`
	OrganizationID string    `json:"organization_id"`
	Scope          string    `json:"scope"`
	Source         string    `json:"source"`
	Basis          string    `json:"basis"`
	ObservedAt     time.Time `json:"observed_at"`
	Rows           []Row     `json:"rows"`
}

// Kind is an explicit projection; no permission inheritance between roles.
func Permission(kind string) string {
	switch kind {
	case "orders", "stock", "cases", "shipments":
		return "admin:read"
	case "own-orders", "own-cases", "own-appointments":
		return "customer:self"
	case "leads":
		return "lead:read"
	case "appointments":
		return "appointment:manage"
	case "factory-destination":
		return "factory:read"
	case "factory-owned":
		return "supply:factory-read"
	case "supply":
		return "supply:read"
	}
	return ""
}
func Authorized(p identity.Principal, kind, org string) bool {
	return p.TenantID != "" && p.Subject != "" && ID(org) && Permission(kind) != "" && p.Allowed(Permission(kind)) && p.AllowedOrganization(org)
}
