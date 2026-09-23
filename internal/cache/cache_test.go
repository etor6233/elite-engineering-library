package cache

import (
	"errors"
	"testing"
)

func TestScopeKey(t *testing.T) {
	k, err := (Scope{TenantID: "tenant-a", Kind: "session", ID: "conv-1"}).Key()
	if err != nil {
		t.Fatal(err)
	}
	if k != "t:tenant-a:session:conv-1" {
		t.Fatalf("unexpected key: %q", k)
	}
	if _, err := (Scope{TenantID: "", Kind: "x"}).Key(); !errors.Is(err, ErrInvalidKey) {
		t.Fatalf("empty tenant accepted: %v", err)
	}
}

func TestValidateKey(t *testing.T) {
	for _, ok := range []string{"a", "t:tenant-a:session:conv-1", "ABC_123.-/"} {
		if err := ValidateKey(ok); err != nil {
			t.Fatalf("valid key %q rejected: %v", ok, err)
		}
	}
	for _, bad := range []string{"", " bad", "bad key!", "k:with:space "} {
		if err := ValidateKey(bad); !errors.Is(err, ErrInvalidKey) {
			t.Fatalf("invalid key %q accepted: %v", bad, err)
		}
	}
}
