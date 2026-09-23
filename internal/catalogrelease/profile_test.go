package catalogrelease

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestCatalogProfilePinnedBoundary(t *testing.T) {
	p := Profile{Schema: "elite-catalog-publication/v1", TenantID: "ca6c9cb7-e548-4f19-a2e6-27d2aa2d4893", OrganizationID: "store", Market: "AR", Currency: "ARS", Origin: "https://catalog.example.invalid", PolicyCode: PolicyCode}
	raw, _ := json.Marshal(p)
	digest := sha256.Sum256(raw)
	hash := hex.EncodeToString(digest[:])
	path := filepath.Join(t.TempDir(), "profile.json")
	if e := os.WriteFile(path, raw, 0600); e != nil {
		t.Fatal(e)
	}
	if got, e := LoadProfile(path, hash); e != nil || got != p {
		t.Fatal("valid profile", got, e)
	}
	if _, e := LoadProfile("relative.json", hash); e == nil {
		t.Fatal("relative path")
	}
	if _, e := LoadProfile(filepath.Dir(path), hash); e == nil {
		t.Fatal("directory path")
	}
	changed := append(append([]byte(nil), raw...), '\n')
	if e := os.WriteFile(path, changed, 0600); e != nil {
		t.Fatal(e)
	}
	if _, e := LoadProfile(path, hash); e == nil {
		t.Fatal("changed bytes under old hash")
	}
	bad := append(append([]byte(nil), raw[:len(raw)-1]...), []byte(`,"secret":"not-an-input"}`)...)
	digest = sha256.Sum256(bad)
	if e := os.WriteFile(path, bad, 0600); e != nil {
		t.Fatal(e)
	}
	if _, e := LoadProfile(path, hex.EncodeToString(digest[:])); e == nil {
		t.Fatal("unknown field admitted")
	}
	p.Origin = "https://user:password@example.invalid"
	if p.Valid() {
		t.Fatal("credential URL")
	}
	p.Origin = "https://example.invalid/path"
	if p.Valid() {
		t.Fatal("origin path")
	}
}
