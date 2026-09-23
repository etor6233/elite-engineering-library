package economy

import (
	"math"
	"testing"
	"time"
)

func TestReviewPositiveCacheScopeCloneAndExpiry(t *testing.T) {
	now := time.Unix(1000, 0)
	c := &SemanticCache{now: func() time.Time { return now }}
	scope := CacheScope{"t", "u", "p1", "k1"}
	embedding := []float32{1, 0}
	if e := c.StoreScoped(scope, embedding, "answer", time.Minute); e != nil {
		t.Fatal(e)
	}
	embedding[0] = 0
	for name, s := range map[string]CacheScope{"same": scope, "tenant": {"other", "u", "p1", "k1"}, "reader": {"t", "other", "p1", "k1"}, "policy": {"t", "u", "p2", "k1"}, "knowledge": {"t", "u", "p1", "k2"}} {
		t.Run(name, func(t *testing.T) {
			answer, ok := c.LookupScoped(s, []float32{1, 0}, .99)
			if ok != (name == "same") || (ok && answer != "answer") {
				t.Fatalf("scope or clone incorrect: %q %v", answer, ok)
			}
		})
	}
	for _, q := range [][]float32{{float32(math.NaN()), 0}, {float32(math.Inf(1)), 0}} {
		if _, ok := c.LookupScoped(scope, q, .99); ok {
			t.Fatal("non-finite lookup returned answer")
		}
	}
	now = now.Add(time.Minute)
	if _, ok := c.LookupScoped(scope, []float32{1, 0}, .99); ok {
		t.Fatal("expired returned")
	}
}
func TestReviewDefectCacheStoresNonfinite(t *testing.T) {
	for name, value := range map[string]float32{"NaN": float32(math.NaN()), "Inf": float32(math.Inf(1))} {
		t.Run(name, func(t *testing.T) {
			c := &SemanticCache{}
			if err := c.StoreScoped(CacheScope{"t", "u", "p1", "k1"}, []float32{value, 1}, "answer", time.Minute); err == nil {
				t.Fatal("invalid non-finite cache embedding stored; lookup still misses")
			}
		})
	}
}
