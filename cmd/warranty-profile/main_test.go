package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	wc "elite.local/enterprise/internal/warrantyclaim"
)

func TestWarrantyProfileMaterializesIdenticallyAndRefusesOverwrite(t *testing.T) {
	source := filepath.Join("..", "..", "deploy", "warranty", "profile.reference.json")
	raw, err := os.ReadFile(source)
	if err != nil {
		t.Fatal(err)
	}
	sha := wc.SHA(raw)
	root := t.TempDir()
	args := func(dest, hash string) []string {
		return []string{"--input", source, "--sha256", hash, "--output", dest, "--tenant-id", "bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbbbb", "--organization-id", "branch", "--factory-organization-id", "maker"}
	}
	first := filepath.Join(root, "first")
	second := filepath.Join(root, "second")
	a, err := materialize(args(first, sha))
	if err != nil {
		t.Fatal(err)
	}
	b, err := materialize(args(second, sha))
	if err != nil {
		t.Fatal(err)
	}
	if a.ProfileSHA256 != b.ProfileSHA256 || a.SourceSHA256 != sha || a.TenantID != "bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbbbb" || a.OrganizationID != "branch" || a.FactoryOrganizationID != "maker" || a.ProductionAuthorized {
		t.Fatal("activation binding", a, b)
	}
	for _, name := range []string{"profile.json", "activation.json"} {
		one, err := os.ReadFile(filepath.Join(first, name))
		if err != nil {
			t.Fatal(err)
		}
		two, err := os.ReadFile(filepath.Join(second, name))
		if err != nil || !bytes.Equal(one, two) {
			t.Fatal("nonportable rebuild", name, err)
		}
	}
	if _, err = materialize(args(first, sha)); err == nil {
		t.Fatal("existing destination overwritten")
	}
	bad := filepath.Join(root, "bad-hash")
	if _, err = materialize(args(bad, strings.Repeat("f", 64))); err == nil {
		t.Fatal("unlocked source admitted")
	}
	if _, err = os.Lstat(bad); !os.IsNotExist(err) {
		t.Fatal("invalid source created output")
	}
	invalid := filepath.Join(root, "invalid-override")
	wrong := args(invalid, sha)
	wrong[7] = "not-a-tenant"
	if _, err = materialize(wrong); err == nil {
		t.Fatal("invalid override admitted")
	}
	if _, err = os.Lstat(invalid); !os.IsNotExist(err) {
		t.Fatal("invalid scope created output")
	}
}
