package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"

	"elite.local/enterprise/internal/platform/identity"
)

type hostWAIdentityVerifier struct {
	calls int
	p     identity.Principal
	token string
}

func (v *hostWAIdentityVerifier) Verify(_ context.Context, raw string) (identity.Principal, error) {
	v.calls++
	if raw != v.token {
		return identity.Principal{}, identity.ErrUnauthenticated
	}
	return v.p, nil
}

func TestWhatsAppHostActivationNoImplicitInputs(t *testing.T) {
	for _, enabled := range []string{"", "false"} {
		calls := 0
		h, e := selectedWhatsAppHost(context.Background(), nil, nil, func(k string) string {
			calls++
			if k != "WHATSAPP_ENABLED" {
				t.Fatal("disabled host read inputs")
			}
			return enabled
		})
		if e != nil || h != nil || calls != 1 {
			t.Fatalf("%v %v calls%d", h, e, calls)
		}
	}
	for _, enabled := range []string{"yes", "TRUE", "true"} {
		if _, e := selectedWhatsAppHost(context.Background(), nil, nil, func(string) string { return enabled }); e == nil {
			t.Fatal("activation without actual dependencies")
		}
	}
}

func TestWhatsAppHostExactFilesAndSecretRotation(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "profile.json")
	raw := []byte("{\"mode\":\"fixture\"}\n")
	if e := os.WriteFile(path, raw, 0600); e != nil {
		t.Fatal(e)
	}
	h := sha256.Sum256(raw)
	if b, e := readWhatsAppExact(path, hex.EncodeToString(h[:]), 64); e != nil || string(b) != string(raw) {
		t.Fatal(e)
	}
	for _, tc := range []struct {
		path, hash string
		limit      int64
	}{{path, hex.EncodeToString(h[:]), 1}, {"relative", hex.EncodeToString(h[:]), 64}, {dir, hex.EncodeToString(h[:]), 64}, {path, "00", 64}} {
		if _, e := readWhatsAppExact(tc.path, tc.hash, tc.limit); e == nil {
			t.Fatal("invalid exact file")
		}
	}
	tokenFile := filepath.Join(dir, "rotating-token")
	if e := os.WriteFile(tokenFile, []byte("first\n"), 0600); e != nil {
		t.Fatal(e)
	}
	s := whatsappHostSecrets{func(string) string { return tokenFile }}
	if v, e := s.WhatsAppToken(context.Background()); e != nil || v != "first" {
		t.Fatal(v, e)
	}
	if e := os.WriteFile(tokenFile, []byte("second\n"), 0600); e != nil {
		t.Fatal(e)
	}
	if v, e := s.WhatsAppToken(context.Background()); e != nil || v != "second" {
		t.Fatal(v, e)
	}
	if e := os.WriteFile(tokenFile, []byte("two\nlines"), 0600); e != nil {
		t.Fatal(e)
	}
	if _, e := s.WhatsAppToken(context.Background()); e == nil {
		t.Fatal("multiline secret")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, e := s.WhatsAppToken(ctx); e == nil {
		t.Fatal("cancelled secret read")
	}
}

func TestWhatsAppHostIdentityRevalidatedAndScoped(t *testing.T) {
	path := filepath.Join(t.TempDir(), "token")
	if e := os.WriteFile(path, []byte("current"), 0600); e != nil {
		t.Fatal(e)
	}
	p := identity.Principal{Subject: "worker", TenantID: "tenant", Organizations: map[string]struct{}{"org": {}}, Permissions: map[string]struct{}{"appointment:manage": {}, "whatsapp:process": {}}}
	v := &hostWAIdentityVerifier{p: p, token: "current"}
	s := whatsappHostIdentity{whatsappFileTokenSource{whatsappHostSecrets{func(string) string { return path }}}, v, "tenant", "org"}
	if _, e := s.Resolve(context.Background()); e != nil {
		t.Fatal(e)
	}
	if raw, e := s.AccessToken(context.Background()); e != nil || raw != "current" {
		t.Fatal(e)
	}
	if v.calls != 2 {
		t.Fatal("identity cached")
	}
	delete(v.p.Permissions, "whatsapp:process")
	if _, e := s.Resolve(context.Background()); e == nil {
		t.Fatal("revoked grant")
	}
	v.p.Permissions["whatsapp:process"] = struct{}{}
	v.p.TenantID = "other"
	if _, e := s.Resolve(context.Background()); e == nil {
		t.Fatal("wrong tenant")
	}
	v.p.TenantID = "tenant"
	delete(v.p.Organizations, "org")
	if _, e := s.Resolve(context.Background()); e == nil {
		t.Fatal("wrong org")
	}
	v.p.Organizations["org"] = struct{}{}
	if e := os.WriteFile(path, []byte("expired"), 0600); e != nil {
		t.Fatal(e)
	}
	if _, e := s.Resolve(context.Background()); e == nil {
		t.Fatal("token verification ignored")
	}
}
