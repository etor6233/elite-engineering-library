package franchisejourney

import (
	"errors"
	"testing"
)

func TestCommercialReleaseRequiresExplicitVersionedEffect(t *testing.T) {
	for _, tc := range []struct {
		name               string
		revision           int
		effect             string
		admitted, writable bool
	}{
		{"readonly", 1, "READ_ONLY_ELIGIBILITY", true, false},
		{"commercial", 2, CommercialReleaseEffect, true, true},
		{"readonly-cannot-upgrade", 1, CommercialReleaseEffect, false, false},
		{"revision-alone-cannot-upgrade", 2, "READ_ONLY_ELIGIBILITY", false, false},
		{"physical-shipment-denied", 2, "POST_SHIPMENT", false, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			d, a := profileFixture(false)
			d.AlgorithmRevision = tc.revision
			d.Options.ReleaseEffect = tc.effect
			raw, a := lockProfile(t, d, a)
			p, err := LoadHandoverProfile(raw, a)
			if tc.admitted {
				if err != nil || p.AllowsCommercialRelease(d.TenantID, d.OrganizationID) != tc.writable {
					t.Fatal("effect admission", err)
				}
			} else if !errors.Is(err, ErrReleaseConditioned) {
				t.Fatal("unsupported effect admitted", err)
			}
			if p.AllowsCommercialRelease(d.TenantID, "foreign") {
				t.Fatal("foreign organization admitted")
			}
		})
	}
}
