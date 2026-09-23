// Package finops provides cloud-cost governance: every resource is tagged and
// inventoried (nothing escapes billing attribution), spend is budgeted per
// tenant fail-closed, and teardown accounts for every resource. It is AUTHORED
// over the COST-FINOPS surface contract and the SRE corpus.
package finops

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	providerRe = regexp.MustCompile(`^[a-z][a-z0-9_-]{0,31}$`)
	kindRe     = regexp.MustCompile(`^[a-z][a-z0-9_-]{0,63}$`)
	regionRe   = regexp.MustCompile(`^[a-z0-9-]{0,63}$`)

	ErrInvalidResource = errors.New("finops: invalid resource")
	ErrMissingTag      = errors.New("finops: missing required tag")
	ErrDuplicate       = errors.New("finops: duplicate resource")
)

// RequiredTags are the tags every resource must carry so cost is attributable
// (tenant, environment, owner). A resource without them cannot be billed
// correctly and is rejected.
var RequiredTags = []string{"tenant", "environment", "owner"}

// Resource is an inventoried cloud resource with a monthly cost estimate.
type Resource struct {
	TenantID               string
	Provider               string // aws | gcp | azure | cloudflare ...
	ID                     string // provider resource id
	Kind                   string // cloud_run | compute | database | storage ...
	Region                 string
	Tags                   map[string]string
	EstMonthlyMinorUnits   int64
}

// Validate enforces tagging and identity so nothing escapes attribution.
func (r Resource) Validate() error {
	if strings.TrimSpace(r.TenantID) == "" || len(r.TenantID) > 64 {
		return fmt.Errorf("%w: tenant", ErrInvalidResource)
	}
	if !providerRe.MatchString(r.Provider) {
		return fmt.Errorf("%w: provider", ErrInvalidResource)
	}
	if strings.TrimSpace(r.ID) == "" || len(r.ID) > 200 {
		return fmt.Errorf("%w: id", ErrInvalidResource)
	}
	if !kindRe.MatchString(r.Kind) {
		return fmt.Errorf("%w: kind", ErrInvalidResource)
	}
	if !regionRe.MatchString(r.Region) {
		return fmt.Errorf("%w: region", ErrInvalidResource)
	}
	for _, tag := range RequiredTags {
		if strings.TrimSpace(r.Tags[tag]) == "" {
			return fmt.Errorf("%w: %s", ErrMissingTag, tag)
		}
	}
	if r.EstMonthlyMinorUnits < 0 {
		return fmt.Errorf("%w: negative cost", ErrInvalidResource)
	}
	return nil
}

// Key returns the provider identity (dedup across tenants is enforced by the
// inventory, not by this key).
func (r Resource) Key() string {
	return r.Provider + "\x00" + r.ID
}
