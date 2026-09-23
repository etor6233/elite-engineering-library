package qr

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

// ErrUnknown conceals whether the target is absent or unauthorized.
var ErrUnknown = errors.New("qr: unknown target")
var ErrInvalidScope = errors.New("qr: invalid trusted tenant scope")
var ErrTenantMismatch = errors.New("qr: tenant scope mismatch")

// Resolver must be bound by the server to the authenticated actor and intended
// operation. It returns true only for an existing object that actor may access
// in tenantID. An existence-only lookup does not satisfy this contract.
// Tenant and actor may not be derived from the scanned reference. Actual domain
// commands must recheck authorization, versions and idempotency when executed.
type Resolver interface {
	Resolve(ctx context.Context, tenantID string, kind Kind, id string) (bool, error)
}

// Verify checks an untrusted reference against a tenant obtained from the
// authenticated server session, then the actor/operation-bound resolver. It
// grants no business effect and does not prove the reference issuer's identity.
// expectedTenant must never come from the QR, request body or user-selected ID.
func Verify(ctx context.Context, encoded string, expectedTenant string, r Resolver) (Payload, error) {
	if strings.TrimSpace(expectedTenant) == "" || strings.TrimSpace(expectedTenant) != expectedTenant || len(expectedTenant) > 64 {
		return Payload{}, ErrInvalidScope
	}
	if err := ctx.Err(); err != nil {
		return Payload{}, err
	}
	p, err := Decode(encoded)
	if err != nil {
		return Payload{}, err
	}
	if p.TenantID != expectedTenant {
		return Payload{}, ErrTenantMismatch
	}
	if r == nil {
		return Payload{}, errors.New("qr: nil resolver")
	}
	ok, err := r.Resolve(ctx, expectedTenant, p.Kind, p.ID)
	if err != nil {
		return Payload{}, err
	}
	if !ok {
		return Payload{}, fmt.Errorf("%w: %s/%s", ErrUnknown, p.Kind, p.ID)
	}
	return p, nil
}
