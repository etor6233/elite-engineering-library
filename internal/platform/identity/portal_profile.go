package identity

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/url"
	"os"
	"slices"
	"strings"
)

// PortalProfile is hash-bound application configuration, not issuer credentials.
type PortalProfile struct {
	document PortalProfileDocument
	digest   string
}
type PortalProfileDocument struct {
	Schema                   string   `json:"schema"`
	ProfileID                string   `json:"profile_id"`
	Revision                 int      `json:"revision"`
	Transport                string   `json:"transport"`
	Issuer                   string   `json:"issuer"`
	AuthorizationEndpoint    string   `json:"authorization_endpoint"`
	TokenEndpoint            string   `json:"token_endpoint"`
	JWKSEndpoint             string   `json:"jwks_endpoint"`
	RevocationEndpoint       string   `json:"revocation_endpoint"`
	EndSessionEndpoint       string   `json:"end_session_endpoint"`
	ClientID                 string   `json:"client_id"`
	CallbackURL              string   `json:"callback_url"`
	PostLogoutURL            string   `json:"post_logout_url"`
	BridgeURL                string   `json:"bridge_url"`
	BackendAudience          string   `json:"backend_audience"`
	ServiceClientID          string   `json:"service_client_id"`
	ServiceSubject           string   `json:"service_subject"`
	ServiceScopes            []string `json:"service_scopes"`
	ServiceAudienceParameter string   `json:"service_audience_parameter"`
	TenantID                 string   `json:"tenant_id"`
	Organizations            []string `json:"organization_ids"`
	AllowedPermissions       []string `json:"allowed_permissions"`
	Scopes                   []string `json:"scopes"`
	MaximumSessionSeconds    int      `json:"maximum_session_seconds"`
	RefreshBeforeSeconds     int      `json:"refresh_before_seconds"`
	RetentionSeconds         int      `json:"retention_seconds"`
}

func (p PortalProfile) SHA256() string { return p.digest }
func (p PortalProfile) Document() PortalProfileDocument {
	d := p.document
	d.Organizations = slices.Clone(d.Organizations)
	d.AllowedPermissions = slices.Clone(d.AllowedPermissions)
	d.Scopes = slices.Clone(d.Scopes)
	d.ServiceScopes = slices.Clone(d.ServiceScopes)
	return d
}
func LoadPortalProfileFile(path, digest string) (PortalProfile, error) {
	f, err := os.Open(path)
	if err != nil {
		return PortalProfile{}, ErrServiceTokenConfiguration
	}
	defer f.Close()
	raw, err := io.ReadAll(io.LimitReader(f, 16385))
	if err != nil {
		return PortalProfile{}, ErrServiceTokenConfiguration
	}
	return LoadPortalProfile(raw, digest)
}
func LoadPortalProfile(raw []byte, digest string) (PortalProfile, error) {
	reject := func() (PortalProfile, error) { return PortalProfile{}, ErrServiceTokenConfiguration }
	if len(raw) == 0 || len(raw) > 16384 || len(digest) != 64 {
		return reject()
	}
	sum := sha256.Sum256(raw)
	if hex.EncodeToString(sum[:]) != digest {
		return reject()
	}
	first := json.NewDecoder(bytes.NewReader(raw))
	tok, err := first.Token()
	if err != nil || tok != json.Delim('{') {
		return reject()
	}
	seen := map[string]bool{}
	for first.More() {
		key, err := first.Token()
		name, ok := key.(string)
		if err != nil || !ok || seen[name] {
			return reject()
		}
		seen[name] = true
		var item json.RawMessage
		if first.Decode(&item) != nil {
			return reject()
		}
	}
	if _, err = first.Token(); err != nil {
		return reject()
	}
	if _, err = first.Token(); err != io.EOF {
		return reject()
	}
	var d PortalProfileDocument
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	if dec.Decode(&d) != nil {
		return reject()
	}
	if d.Schema != "elite.oidc.portal-lifecycle.v1" || d.Revision != 1 || !serviceAtom(d.ProfileID) || !serviceAtom(d.ClientID) || !serviceAtom(d.ServiceClientID) || !serviceAtom(d.ServiceSubject) || !serviceAtom(d.TenantID) || !serviceAtom(d.BackendAudience) {
		return reject()
	}
	if d.Transport != "TLS" && d.Transport != "LOOPBACK_FIXTURE" {
		return reject()
	}
	for _, u := range []string{d.Issuer, d.AuthorizationEndpoint, d.TokenEndpoint, d.JWKSEndpoint, d.RevocationEndpoint, d.EndSessionEndpoint, d.CallbackURL, d.PostLogoutURL, d.BridgeURL} {
		if !serviceURL(u, d.Transport) {
			return reject()
		}
	}
	if !strings.HasSuffix(d.CallbackURL, "/api/auth/callback") || strings.TrimSuffix(d.CallbackURL, "/api/auth/callback")+"/" != d.PostLogoutURL || strings.HasSuffix(d.BridgeURL, "/") {
		return reject()
	}
	bridge, err := url.Parse(d.BridgeURL)
	if err != nil || bridge.Path != "" && bridge.Path != "/" {
		return reject()
	}
	if !serviceSet(d.Organizations) || !serviceSet(d.AllowedPermissions) || !serviceSet(d.ServiceScopes) || !serviceSet(d.Scopes) || !slices.Contains(d.Scopes, "openid") || !slices.Contains(d.Scopes, "offline_access") {
		return reject()
	}
	if d.ServiceAudienceParameter != "" && d.ServiceAudienceParameter != d.BackendAudience {
		return reject()
	}
	if d.MaximumSessionSeconds < 300 || d.MaximumSessionSeconds > 86400 || d.RefreshBeforeSeconds < 5 || d.RefreshBeforeSeconds > 300 || d.RefreshBeforeSeconds*2 >= d.MaximumSessionSeconds || d.RetentionSeconds < 3600 || d.RetentionSeconds > 604800 {
		return reject()
	}
	return PortalProfile{d, digest}, nil
}

// ServiceAllowed grants only this private storage capability under a fixed profile.
func (p PortalProfile) ServiceAllowed(principal Principal) bool {
	if p.digest == "" || principal.Subject != p.document.ServiceSubject || principal.TenantID != p.document.TenantID || len(principal.Permissions) != 1 {
		return false
	}
	if _, ok := principal.Permissions["portal-session:manage"]; !ok {
		return false
	}
	organizations := make([]string, 0, len(principal.Organizations))
	for id := range principal.Organizations {
		organizations = append(organizations, id)
	}
	return serviceEqual(organizations, p.document.Organizations)
}
