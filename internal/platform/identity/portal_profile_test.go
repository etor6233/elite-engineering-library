package identity

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"testing"
)

var portalProfileFixture = []byte(`{
  "schema": "elite.oidc.portal-lifecycle.v1",
  "profile_id": "reference-portal",
  "revision": 1,
  "transport": "TLS",
  "issuer": "https://issuer.invalid",
  "authorization_endpoint": "https://issuer.invalid/authorize",
  "token_endpoint": "https://issuer.invalid/token",
  "jwks_endpoint": "https://issuer.invalid/jwks",
  "revocation_endpoint": "https://issuer.invalid/revoke",
  "end_session_endpoint": "https://issuer.invalid/logout",
  "client_id": "reference-portal",
  "callback_url": "https://portal.invalid/api/auth/callback",
  "post_logout_url": "https://portal.invalid/",
  "bridge_url": "https://backend.invalid",
  "backend_audience": "reference-api",
  "service_client_id": "reference-portal-maintenance",
  "service_subject": "reference-portal-maintenance",
  "service_scopes": [
    "portal-session:manage"
  ],
  "service_audience_parameter": "reference-api",
  "tenant_id": "00000000-0000-4000-8000-000000000001",
  "organization_ids": [
    "00000000-0000-4000-8000-000000000002"
  ],
  "allowed_permissions": [
    "customer:read"
  ],
  "scopes": [
    "openid",
    "offline_access"
  ],
  "maximum_session_seconds": 3600,
  "refresh_before_seconds": 30,
  "retention_seconds": 86400
}`)

func portalProfileFixtureHash(raw []byte) string {
	sum := sha256.Sum256(raw)
	return hex.EncodeToString(sum[:])
}
func TestPortalProfileContract(t *testing.T) {
	p, err := LoadPortalProfile(portalProfileFixture, portalProfileFixtureHash(portalProfileFixture))
	if err != nil {
		t.Fatal(err)
	}
	doc := p.Document()
	doc.Organizations[0] = "foreign"
	if p.Document().Organizations[0] == "foreign" {
		t.Fatal("mutable profile")
	}
	for _, change := range []map[string]any{{"bridge_url": "https://backend.invalid/path"}, {"bridge_url": "https://backend.invalid/"}, {"allowed_permissions": []string{"*"}}, {"transport": "LOOPBACK_FIXTURE"}, {"post_logout_url": "https://foreign.invalid/"}} {
		var value map[string]any
		if json.Unmarshal(portalProfileFixture, &value) != nil {
			t.Fatal("fixture")
		}
		for k, v := range change {
			value[k] = v
		}
		raw, _ := json.Marshal(value)
		if _, err := LoadPortalProfile(raw, portalProfileFixtureHash(raw)); err == nil {
			t.Fatal("invalid profile accepted", change)
		}
	}
	duplicate := append([]byte(`{"profile_id":"foreign",`), portalProfileFixture[1:]...)
	if _, err := LoadPortalProfile(duplicate, portalProfileFixtureHash(duplicate)); err == nil {
		t.Fatal("duplicate")
	}
	if _, err := LoadPortalProfile(portalProfileFixture, "0000000000000000000000000000000000000000000000000000000000000000"); err == nil {
		t.Fatal("wrong hash")
	}
}
func FuzzPortalProfile(f *testing.F) {
	f.Add(portalProfileFixture)
	f.Add([]byte(`{"profile_id":"foreign","profile_id":"fixture"}`))
	f.Add([]byte("null"))
	f.Fuzz(func(t *testing.T, raw []byte) {
		p, err := LoadPortalProfile(raw, portalProfileFixtureHash(raw))
		if err == nil {
			if p.SHA256() != portalProfileFixtureHash(raw) || len(raw) > 16384 {
				t.Fatal("unbound profile")
			}
			d := p.Document()
			if d.MaximumSessionSeconds < 300 || d.MaximumSessionSeconds > 86400 || len(d.Organizations) == 0 {
				t.Fatal("unbounded accepted profile")
			}
		}
	})
}
