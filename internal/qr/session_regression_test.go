package qr

import (
	"context"
	"testing"
)

// Both tenants genuinely exist in this synthetic domain. This test does not
// rely on the old fake resolver accidentally lacking the other tenant.
type multiTenantResolver struct{ calls int }

func (r *multiTenantResolver) Resolve(_ context.Context, tenant string, kind Kind, id string) (bool, error) {
	r.calls++
	return (tenant == "tenant-a" || tenant == "tenant-b") && kind == KindProduct && id == "SKU-1", nil
}
func TestExistingCrossTenantTargetCannotPassTrustedSession(t *testing.T) {
	p, _ := NewPayload("tenant-b", KindProduct, "SKU-1")
	encoded, _ := p.Encode()
	r := &multiTenantResolver{}
	if _, err := Verify(context.Background(), encoded, "tenant-a", r); err == nil {
		t.Fatal("cross-tenant QR accepted while trusted session is tenant-a")
	}
	if r.calls != 0 {
		t.Fatalf("cross-tenant resolver was called %d times", r.calls)
	}
}
func TestRecomputedChecksumDoesNotGrantTenantAccess(t *testing.T) {
	p, _ := NewPayload("tenant-a", KindProduct, "SKU-1")
	// An untrusted sender can recompute every unkeyed checksum field. Decode must
	// not be mislabeled authentication; session scope is checked independently.
	forged, _ := NewPayload("tenant-b", p.Kind, p.ID)
	encoded, _ := forged.Encode()
	if _, err := Decode(encoded); err != nil {
		t.Fatalf("fixture needs a structurally valid recomputed checksum: %v", err)
	}
	r := &multiTenantResolver{}
	if _, err := Verify(context.Background(), encoded, "tenant-a", r); err == nil {
		t.Fatal("recomputed checksum granted another tenant's target")
	}
	if r.calls != 0 {
		t.Fatal("forged tenant reached the resolver")
	}
}
