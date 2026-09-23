package aifoundation

import (
	"errors"
	"strings"
	"sync"
	"testing"
)

func ref(provider, model, version string) ModelRef {
	return ModelRef{Provider: provider, Model: model, Version: version, Digest: strings.Repeat("a", 64)}
}

func TestModelRefValidate(t *testing.T) {
	valid := ref("openai", "gpt-4.1-mini", "2025-04-14")
	if err := valid.Validate(); err != nil {
		t.Fatalf("valid ref rejected: %v", err)
	}
	if err := (ModelRef{Model: "x", Version: "v1", Digest: strings.Repeat("a", 64)}).Validate(); !errors.Is(err, ErrInvalidModelRef) {
		t.Fatalf("missing provider accepted: %v", err)
	}
	mutable := ref("openai", "gpt-4.1-mini", "latest")
	if err := mutable.Validate(); !errors.Is(err, ErrInvalidModelRef) {
		t.Fatalf("mutable version accepted: %v", err)
	}
	badDigest := ref("openai", "gpt-4.1-mini", "v1")
	badDigest.Digest = "zz"
	if err := badDigest.Validate(); !errors.Is(err, ErrInvalidModelRef) {
		t.Fatalf("bad digest accepted: %v", err)
	}
}

func TestRegistrySingleActive(t *testing.T) {
	var r VersionRegistry
	a := ref("openai", "gpt-4.1-mini", "2025-04-14")
	b := ref("openai", "gpt-4.1-mini", "2025-06-01")
	b.Digest = strings.Repeat("b", 64)
	if err := r.Register(a); err != nil {
		t.Fatal(err)
	}
	if err := r.Register(b); err != nil {
		t.Fatal(err)
	}
	if err := r.Promote(a, true); err != nil {
		t.Fatal(err)
	}
	active, ok := r.Active()
	if !ok || active != a {
		t.Fatalf("expected %+v active, got %+v (ok=%v)", a, active, ok)
	}
	if err := r.Promote(b, true); err != nil {
		t.Fatal(err)
	}
	active, ok = r.Active()
	if !ok || active != b {
		t.Fatalf("expected %+v active after promote, got %+v", b, active)
	}
}

func TestRegistryPromoteRequiresGate(t *testing.T) {
	var r VersionRegistry
	a := ref("openai", "gpt-4.1-mini", "2025-04-14")
	if err := r.Register(a); err != nil {
		t.Fatal(err)
	}
	if err := r.Promote(a, false); err == nil {
		t.Fatal("promote without gate passed was accepted")
	}
	if _, ok := r.Active(); ok {
		t.Fatal("no version should be active after rejected promote")
	}
}

func TestRegistryDuplicateRejected(t *testing.T) {
	var r VersionRegistry
	a := ref("openai", "gpt-4.1-mini", "2025-04-14")
	if err := r.Register(a); err != nil {
		t.Fatal(err)
	}
	if err := r.Register(a); err == nil {
		t.Fatal("duplicate digest accepted")
	}
}

func TestRegistryRollback(t *testing.T) {
	var r VersionRegistry
	a := ref("openai", "gpt-4.1-mini", "2025-04-14")
	if err := r.Register(a); err != nil {
		t.Fatal(err)
	}
	if err := r.Promote(a, true); err != nil {
		t.Fatal(err)
	}
	if err := r.Rollback(a, "incident"); err != nil {
		t.Fatal(err)
	}
	if _, ok := r.Active(); ok {
		t.Fatal("active version still present after rollback")
	}
}

func TestRegistryConcurrentPromoteOneActive(t *testing.T) {
	var r VersionRegistry
	const n = 8
	refs := make([]ModelRef, n)
	for i := 0; i < n; i++ {
		refs[i] = ref("openai", "gpt-4.1-mini", "v")
		refs[i].Version = "v" + string(rune('0'+i))
		refs[i].Digest = strings.Repeat(string("0123456789abcdef"[i]), 64)
		if err := r.Register(refs[i]); err != nil {
			t.Fatal(err)
		}
	}
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(m ModelRef) {
			defer wg.Done()
			_ = r.Promote(m, true)
		}(refs[i])
	}
	wg.Wait()

	active, ok := r.Active()
	if !ok {
		t.Fatal("expected exactly one active version")
	}
	count := 0
	for _, ev := range r.Audit() {
		if ev.To == StateActive {
			count++
		}
	}
	if count != n {
		t.Fatalf("expected %d promote events, got %d", n, count)
	}
	if active.Digest == "" {
		t.Fatal("active digest empty")
	}
}
