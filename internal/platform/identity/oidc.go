package identity

import (
	"context"
	"errors"
	"fmt"

	"github.com/coreos/go-oidc/v3/oidc"
)

var ErrUnauthenticated = errors.New("unauthenticated")

type Principal struct {
	Subject       string
	TenantID      string
	Permissions   map[string]struct{}
	Organizations map[string]struct{}
}

func (p Principal) Allowed(permission string) bool {
	if _, all := p.Permissions["*"]; all {
		return true
	}
	_, allowed := p.Permissions[permission]
	return allowed
}

func (p Principal) AllowedOrganization(organizationID string) bool {
	if organizationID == "" {
		return false
	}
	if _, all := p.Permissions["*"]; all {
		return true
	}
	_, allowed := p.Organizations[organizationID]
	return allowed
}

type Verifier interface {
	Verify(context.Context, string) (Principal, error)
}

type OIDCVerifier struct{ verifier *oidc.IDTokenVerifier }

func NewOIDCVerifier(ctx context.Context, issuer, audience string) (*OIDCVerifier, error) {
	provider, err := oidc.NewProvider(ctx, issuer)
	if err != nil {
		return nil, err
	}
	return &OIDCVerifier{verifier: provider.Verifier(&oidc.Config{
		ClientID:             audience,
		SupportedSigningAlgs: []string{oidc.RS256},
	})}, nil
}

func (v *OIDCVerifier) Verify(ctx context.Context, raw string) (Principal, error) {
	token, err := v.verifier.Verify(ctx, raw)
	if err != nil {
		return Principal{}, fmt.Errorf("%w: token rejected", ErrUnauthenticated)
	}
	var claims struct {
		TenantID        string   `json:"tenant_id"`
		Permissions     []string `json:"permissions"`
		OrganizationIDs []string `json:"organization_ids"`
	}
	if err := token.Claims(&claims); err != nil || token.Subject == "" || claims.TenantID == "" {
		return Principal{}, fmt.Errorf("%w: required claims missing", ErrUnauthenticated)
	}
	permissions := make(map[string]struct{}, len(claims.Permissions))
	for _, permission := range claims.Permissions {
		if permission != "" {
			permissions[permission] = struct{}{}
		}
	}
	organizations := make(map[string]struct{}, len(claims.OrganizationIDs))
	for _, organizationID := range claims.OrganizationIDs {
		if organizationID != "" {
			organizations[organizationID] = struct{}{}
		}
	}
	return Principal{Subject: token.Subject, TenantID: claims.TenantID, Permissions: permissions, Organizations: organizations}, nil
}
