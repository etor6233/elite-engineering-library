package identity

// AUTHORED configuration and claim binding glue. OAuth/OIDC runs in pinned SDKs.
import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/netip"
	"net/url"
	"os"
	"slices"
	"strings"
)

var ErrServiceTokenConfiguration = errors.New("service token configuration rejected")
var ErrServiceTokenUnavailable = errors.New("service token unavailable")

type serviceTokenDocument struct {
	Schema                 string   `json:"schema"`
	ProfileID              string   `json:"profile_id"`
	Revision               int      `json:"revision"`
	Transport              string   `json:"transport"`
	Issuer                 string   `json:"issuer"`
	TokenEndpoint          string   `json:"token_endpoint"`
	JWKSEndpoint           string   `json:"jwks_endpoint"`
	ClientID               string   `json:"client_id"`
	ClientAuthentication   string   `json:"client_authentication"`
	Audience               string   `json:"audience"`
	Subject                string   `json:"subject"`
	TenantID               string   `json:"tenant_id"`
	Organizations          []string `json:"organization_ids"`
	Permissions            []string `json:"permissions"`
	Scopes                 []string `json:"scopes"`
	Resource               string   `json:"resource,omitempty"`
	AudienceParameter      string   `json:"audience_parameter,omitempty"`
	MaximumLifetimeSeconds int      `json:"maximum_lifetime_seconds"`
	RefreshBeforeSeconds   int      `json:"refresh_before_seconds"`
}

// ServiceTokenProfile is immutable after hash verification; zero value is invalid.
type ServiceTokenProfile struct {
	document serviceTokenDocument
	hash     string
}
type ServiceTokenBinding struct {
	ProfileID, DocumentSHA256, Transport, Issuer, Audience, ClientID, Subject, TenantID string
	Organizations, Permissions, Scopes                                                  []string
}

func (p ServiceTokenProfile) Binding() ServiceTokenBinding {
	d := p.document
	return ServiceTokenBinding{d.ProfileID, p.hash, d.Transport, d.Issuer, d.Audience, d.ClientID, d.Subject, d.TenantID, slices.Clone(d.Organizations), slices.Clone(d.Permissions), slices.Clone(d.Scopes)}
}
func LoadServiceTokenProfileFile(path, digest string) (ServiceTokenProfile, error) {
	f, err := os.Open(path)
	if err != nil {
		return ServiceTokenProfile{}, ErrServiceTokenConfiguration
	}
	defer f.Close()
	raw, err := io.ReadAll(io.LimitReader(f, 16385))
	if err != nil {
		return ServiceTokenProfile{}, ErrServiceTokenConfiguration
	}
	return LoadServiceTokenProfile(raw, digest)
}
func LoadServiceTokenProfile(raw []byte, digest string) (ServiceTokenProfile, error) {
	reject := func() (ServiceTokenProfile, error) { return ServiceTokenProfile{}, ErrServiceTokenConfiguration }
	if len(raw) == 0 || len(raw) > 16384 || len(digest) != 64 {
		return reject()
	}
	sum := sha256.Sum256(raw)
	if hex.EncodeToString(sum[:]) != digest {
		return reject()
	}
	// Reject duplicate member names before normal decoding; this schema has no objects nested in members.
	first := json.NewDecoder(bytes.NewReader(raw))
	token, err := first.Token()
	if err != nil || token != json.Delim('{') {
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
		var value json.RawMessage
		if first.Decode(&value) != nil {
			return reject()
		}
	}
	if token, err = first.Token(); err != nil || token != json.Delim('}') {
		return reject()
	}
	if _, err = first.Token(); err != io.EOF {
		return reject()
	}
	var d serviceTokenDocument
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	if dec.Decode(&d) != nil {
		return reject()
	}
	if d.Schema != "elite.oidc.service-token.v1" || d.Revision != 1 || !serviceAtom(d.ProfileID) || !serviceAtom(d.ClientID) || !serviceAtom(d.Audience) || !serviceAtom(d.Subject) || !serviceAtom(d.TenantID) {
		return reject()
	}
	if d.Transport != "TLS" && d.Transport != "LOOPBACK_FIXTURE" {
		return reject()
	}
	if !serviceURL(d.Issuer, d.Transport) || !serviceURL(d.TokenEndpoint, d.Transport) || !serviceURL(d.JWKSEndpoint, d.Transport) {
		return reject()
	}
	if d.ClientAuthentication != "client_secret_basic" && d.ClientAuthentication != "client_secret_post" {
		return reject()
	}
	if !serviceSet(d.Organizations) || !serviceSet(d.Permissions) || !serviceSet(d.Scopes) {
		return reject()
	}
	if d.Resource != "" && (!serviceURL(d.Resource, "TLS") || d.AudienceParameter != "") {
		return reject()
	}
	if d.AudienceParameter != "" && d.AudienceParameter != d.Audience {
		return reject()
	}
	if d.MaximumLifetimeSeconds < 60 || d.MaximumLifetimeSeconds > 86400 || d.RefreshBeforeSeconds < 5 || d.RefreshBeforeSeconds > 300 || d.RefreshBeforeSeconds*2 >= d.MaximumLifetimeSeconds {
		return reject()
	}
	return ServiceTokenProfile{document: d, hash: digest}, nil
}
func serviceAtom(s string) bool {
	return len(s) > 0 && len(s) <= 200 && strings.TrimSpace(s) == s && !strings.ContainsAny(s, "\x00\r\n\t ") && s != "*"
}
func serviceSet(values []string) bool {
	if len(values) == 0 || len(values) > 100 {
		return false
	}
	seen := map[string]bool{}
	for _, v := range values {
		if !serviceAtom(v) || seen[v] {
			return false
		}
		seen[v] = true
	}
	return true
}
func serviceEqual(a, b []string) bool {
	if !serviceSet(a) || len(a) != len(b) {
		return false
	}
	a = slices.Clone(a)
	b = slices.Clone(b)
	slices.Sort(a)
	slices.Sort(b)
	return slices.Equal(a, b)
}
func serviceURL(raw, mode string) bool {
	u, err := url.Parse(raw)
	if err != nil || len(raw) > 2048 || u.Host == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" || u.Opaque != "" || u.String() != raw {
		return false
	}
	if mode == "LOOPBACK_FIXTURE" {
		ip, err := netip.ParseAddr(u.Hostname())
		return err == nil && ip.IsLoopback() && (u.Scheme == "http" || u.Scheme == "https")
	}
	return u.Scheme == "https"
}
