package cache

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestMemorySetGet(t *testing.T) {
	s := NewMemoryStore(1 << 20)
	ctx := context.Background()
	if err := s.Set(ctx, "k1", []byte("v1"), time.Minute); err != nil {
		t.Fatal(err)
	}
	v, ok, err := s.Get(ctx, "k1")
	if err != nil || !ok || string(v) != "v1" {
		t.Fatalf("get: %q ok=%v err=%v", v, ok, err)
	}
	if _, ok, _ := s.Get(ctx, "missing"); ok {
		t.Fatal("missing key reported hit")
	}
}

func TestMemoryGetReturnsCopy(t *testing.T) {
	s := NewMemoryStore(1 << 20)
	ctx := context.Background()
	_ = s.Set(ctx, "k", []byte("abc"), time.Minute)
	v, _, _ := s.Get(ctx, "k")
	v[0] = 'z'
	v2, _, _ := s.Get(ctx, "k")
	if string(v2) != "abc" {
		t.Fatalf("cache value mutated: %q", v2)
	}
}

func TestMemoryTTLExpiry(t *testing.T) {
	s := NewMemoryStore(1 << 20)
	var now time.Time
	s.clock = func() time.Time { return now }
	ctx := context.Background()
	_ = s.Set(ctx, "k", []byte("v"), time.Minute)
	now = now.Add(2 * time.Minute)
	if _, ok, _ := s.Get(ctx, "k"); ok {
		t.Fatal("expired key reported hit")
	}
	if s.Len() != 0 {
		t.Fatalf("expired key not evicted, len=%d", s.Len())
	}
}

func TestMemoryLRUEviction(t *testing.T) {
	s := NewMemoryStore(30) // budget: three 10-byte values
	ctx := context.Background()
	_ = s.Set(ctx, "k1", []byte("aaaaaaaaaa"), time.Minute)
	_ = s.Set(ctx, "k2", []byte("bbbbbbbbbb"), time.Minute)
	_ = s.Set(ctx, "k3", []byte("cccccccccc"), time.Minute)
	// touch k1 so it is most recently used
	_, _, _ = s.Get(ctx, "k1")
	_ = s.Set(ctx, "k4", []byte("dddddddddd"), time.Minute) // exceeds budget → evict LRU (k2)
	if _, ok, _ := s.Get(ctx, "k2"); ok {
		t.Fatal("LRU entry k2 should have been evicted")
	}
	for _, k := range []string{"k1", "k3", "k4"} {
		if _, ok, _ := s.Get(ctx, k); !ok {
			t.Fatalf("expected %s to remain", k)
		}
	}
}

func TestMemoryDeleteIdempotent(t *testing.T) {
	s := NewMemoryStore(1 << 20)
	ctx := context.Background()
	_ = s.Set(ctx, "k", []byte("v"), time.Minute)
	if err := s.Delete(ctx, "k"); err != nil {
		t.Fatal(err)
	}
	if err := s.Delete(ctx, "k"); err != nil {
		t.Fatalf("delete not idempotent: %v", err)
	}
}

func TestMemoryValidation(t *testing.T) {
	s := NewMemoryStore(1 << 20)
	ctx := context.Background()
	if err := s.Set(ctx, "bad key!", []byte("v"), time.Minute); !errors.Is(err, ErrInvalidKey) {
		t.Fatalf("invalid key accepted: %v", err)
	}
	if err := s.Set(ctx, "k", []byte("v"), 0); !errors.Is(err, ErrInvalidTTL) {
		t.Fatalf("zero ttl accepted: %v", err)
	}
	if err := s.Set(ctx, "k", make([]byte, 2<<20), time.Minute); !errors.Is(err, ErrTooLarge) {
		t.Fatalf("oversized value accepted: %v", err)
	}
}
