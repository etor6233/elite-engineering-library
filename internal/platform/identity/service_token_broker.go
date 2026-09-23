package identity

// AUTHORED lifecycle/transport glue; protocol and JWT verification are DEPENDENCY_PIN.
import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

type ServiceSecretReader func(context.Context) (string, error)

// ServiceSecretFile reads on each grant, allowing externally managed secret rotation.
// Files and their directory/ACL are provisioned by the host secret manager, never the pack.
func ServiceSecretFile(path string) ServiceSecretReader {
	return func(ctx context.Context) (string, error) {
		if ctx.Err() != nil {
			return "", ErrServiceTokenUnavailable
		}
		f, err := os.Open(path)
		if err != nil {
			return "", ErrServiceTokenUnavailable
		}
		defer f.Close()
		info, err := f.Stat()
		if err != nil || !info.Mode().IsRegular() || info.Size() > 4096 {
			return "", ErrServiceTokenUnavailable
		}
		raw, err := io.ReadAll(io.LimitReader(f, 4097))
		if err != nil || len(raw) > 4096 {
			return "", ErrServiceTokenUnavailable
		}
		value := strings.TrimSuffix(strings.TrimSuffix(string(raw), "\n"), "\r")
		if len(value) < 1 || strings.ContainsAny(value, "\x00\r\n") {
			return "", ErrServiceTokenUnavailable
		}
		return value, nil
	}
}

type ServiceTokenBroker struct {
	profile  ServiceTokenProfile
	secret   ServiceSecretReader
	client   *http.Client
	verifier *oidc.IDTokenVerifier
	gate     chan struct{}
	token    *oauth2.Token
	now      func() time.Time
}

func NewServiceTokenBroker(ctx context.Context, p ServiceTokenProfile, secret ServiceSecretReader) (*ServiceTokenBroker, error) {
	return newServiceTokenBroker(ctx, p, secret, http.DefaultTransport)
}

// The transport seam is private and used for a loopback TLS fixture with its own CA.
func newServiceTokenBroker(ctx context.Context, p ServiceTokenProfile, secret ServiceSecretReader, transport http.RoundTripper) (*ServiceTokenBroker, error) {
	if p.hash == "" || secret == nil {
		return nil, ErrServiceTokenConfiguration
	}
	d := p.document
	client := &http.Client{Timeout: 15 * time.Second, Transport: serviceTokenTransport{base: transport, allowed: map[string]bool{strings.TrimSuffix(d.Issuer, "/") + "/.well-known/openid-configuration": true, d.TokenEndpoint: true, d.JWKSEndpoint: true}}, CheckRedirect: func(*http.Request, []*http.Request) error { return ErrServiceTokenUnavailable }}
	provider, err := oidc.NewProvider(oidc.ClientContext(ctx, client), d.Issuer)
	if err != nil {
		return nil, ErrServiceTokenConfiguration
	}
	var metadata struct {
		JWKSEndpoint string   `json:"jwks_uri"`
		AuthMethods  []string `json:"token_endpoint_auth_methods_supported"`
		Grants       []string `json:"grant_types_supported"`
	}
	if provider.Claims(&metadata) != nil || provider.Endpoint().TokenURL != d.TokenEndpoint || metadata.JWKSEndpoint != d.JWKSEndpoint {
		return nil, ErrServiceTokenConfiguration
	}
	contains := func(list []string, value string) bool {
		for _, item := range list {
			if item == value {
				return true
			}
		}
		return false
	}
	if len(metadata.AuthMethods) > 0 && !contains(metadata.AuthMethods, d.ClientAuthentication) || len(metadata.Grants) > 0 && !contains(metadata.Grants, "client_credentials") {
		return nil, ErrServiceTokenConfiguration
	}
	b := &ServiceTokenBroker{profile: p, secret: secret, client: client, gate: make(chan struct{}, 1), now: time.Now}
	b.verifier = provider.Verifier(&oidc.Config{ClientID: d.Audience, SupportedSigningAlgs: []string{oidc.RS256}, Now: func() time.Time { return b.now() }})
	return b, nil
}
func (b *ServiceTokenBroker) Binding() ServiceTokenBinding { return b.profile.Binding() }

// AccessToken returns only a verified short-lived bearer. Failed renewals do not
// return the previous token and do not cache malformed/foreign provider responses.
func (b *ServiceTokenBroker) AccessToken(ctx context.Context) (string, error) {
	if b == nil || ctx == nil {
		return "", ErrServiceTokenUnavailable
	}
	select {
	case b.gate <- struct{}{}:
		defer func() { <-b.gate }()
	case <-ctx.Done():
		return "", ErrServiceTokenUnavailable
	}
	if ctx.Err() != nil {
		return "", ErrServiceTokenUnavailable
	}
	d := b.profile.document
	threshold := b.now().Add(time.Duration(d.RefreshBeforeSeconds) * time.Second)
	if b.token != nil && b.token.Expiry.After(threshold) {
		return b.token.AccessToken, nil
	}
	b.token = nil
	secret, err := b.secret(ctx)
	if err != nil || len(secret) < 1 || len(secret) > 4096 || strings.ContainsAny(secret, "\x00\r\n") {
		return "", ErrServiceTokenUnavailable
	}
	style := oauth2.AuthStyleInHeader
	if d.ClientAuthentication == "client_secret_post" {
		style = oauth2.AuthStyleInParams
	}
	params := url.Values{}
	if d.Resource != "" {
		params.Set("resource", d.Resource)
	}
	if d.AudienceParameter != "" {
		params.Set("audience", d.AudienceParameter)
	}
	config := clientcredentials.Config{ClientID: d.ClientID, ClientSecret: secret, TokenURL: d.TokenEndpoint, Scopes: append([]string(nil), d.Scopes...), AuthStyle: style, EndpointParams: params}
	token, err := config.Token(context.WithValue(ctx, oauth2.HTTPClient, b.client))
	if err != nil || token == nil || !strings.EqualFold(token.TokenType, "Bearer") || len(token.AccessToken) == 0 || len(token.AccessToken) > 16384 || token.RefreshToken != "" || token.Expiry.IsZero() {
		return "", ErrServiceTokenUnavailable
	}
	jwt, err := b.verifier.Verify(ctx, token.AccessToken)
	if err != nil {
		return "", ErrServiceTokenUnavailable
	}
	var claims struct {
		TenantID        string   `json:"tenant_id"`
		Organizations   []string `json:"organization_ids"`
		Permissions     []string `json:"permissions"`
		Scope           string   `json:"scope"`
		NotBefore       int64    `json:"nbf"`
		AuthorizedParty string   `json:"azp"`
		ClientID        string   `json:"client_id"`
	}
	if jwt.Claims(&claims) != nil || jwt.Issuer != d.Issuer || jwt.Subject != d.Subject || claims.TenantID != d.TenantID || !serviceEqual(claims.Organizations, d.Organizations) || !serviceEqual(claims.Permissions, d.Permissions) || !serviceEqual(strings.Fields(claims.Scope), d.Scopes) || len(jwt.Audience) != 1 || jwt.Audience[0] != d.Audience {
		return "", ErrServiceTokenUnavailable
	}
	now := b.now()
	if claims.NotBefore > now.Unix() || jwt.IssuedAt.IsZero() || jwt.IssuedAt.After(now.Add(30*time.Second)) || jwt.Expiry.Sub(jwt.IssuedAt) > time.Duration(d.MaximumLifetimeSeconds)*time.Second || !jwt.Expiry.After(threshold) {
		return "", ErrServiceTokenUnavailable
	}
	if claims.AuthorizedParty != "" && claims.AuthorizedParty != d.ClientID || claims.ClientID != "" && claims.ClientID != d.ClientID {
		return "", ErrServiceTokenUnavailable
	}
	if raw := token.Extra("scope"); raw != nil {
		scope, ok := raw.(string)
		if !ok || !serviceEqual(strings.Fields(scope), d.Scopes) {
			return "", ErrServiceTokenUnavailable
		}
	}
	if jwt.Expiry.Before(token.Expiry) {
		token.Expiry = jwt.Expiry
	}
	if !token.Expiry.After(threshold) || token.Expiry.Sub(now) > time.Duration(d.MaximumLifetimeSeconds+30)*time.Second {
		return "", ErrServiceTokenUnavailable
	}
	b.token = token
	return token.AccessToken, nil
}

type serviceTokenTransport struct {
	base    http.RoundTripper
	allowed map[string]bool
}

func (t serviceTokenTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	if !t.allowed[request.URL.String()] || (request.Method != http.MethodGet && request.Method != http.MethodPost) {
		return nil, ErrServiceTokenUnavailable
	}
	response, err := t.base.RoundTrip(request)
	if err != nil {
		return nil, ErrServiceTokenUnavailable
	}
	defer response.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(response.Body, 262145))
	if err != nil || len(raw) > 262144 {
		return nil, ErrServiceTokenUnavailable
	}
	response.Body = io.NopCloser(bytes.NewReader(raw))
	return response, nil
}
