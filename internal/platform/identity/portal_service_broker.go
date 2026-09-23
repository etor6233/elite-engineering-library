package identity

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

// Derived strictly from the hash-bound portal service identity; no body tenant.
func NewPortalServiceTokenBroker(ctx context.Context, p PortalProfile, reader ServiceSecretReader) (*ServiceTokenBroker, error) {
	if p.SHA256() == "" {
		return nil, ErrServiceTokenConfiguration
	}
	d := p.Document()
	document := serviceTokenDocument{Schema: "elite.oidc.service-token.v1", ProfileID: d.ProfileID, Revision: 1, Transport: d.Transport, Issuer: d.Issuer, TokenEndpoint: d.TokenEndpoint, JWKSEndpoint: d.JWKSEndpoint, ClientID: d.ServiceClientID, ClientAuthentication: "client_secret_basic", Audience: d.BackendAudience, Subject: d.ServiceSubject, TenantID: d.TenantID, Organizations: d.Organizations, Permissions: []string{"portal-session:manage"}, Scopes: d.ServiceScopes, AudienceParameter: d.ServiceAudienceParameter, MaximumLifetimeSeconds: 3600, RefreshBeforeSeconds: 10}
	raw, err := json.Marshal(document)
	if err != nil {
		return nil, ErrServiceTokenConfiguration
	}
	sum := sha256.Sum256(raw)
	profile, err := LoadServiceTokenProfile(raw, hex.EncodeToString(sum[:]))
	if err != nil {
		return nil, err
	}
	return NewServiceTokenBroker(ctx, profile, reader)
}
