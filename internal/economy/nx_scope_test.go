package economy

import (
	"testing"
	"time"
)

func TestNXCacheScopeTTLAndRevisions(t *testing.T) {
	now := time.Unix(1000, 0)
	c := &SemanticCache{now: func() time.Time { return now }}
	scope := CacheScope{"t", "u", "p1", "k1"}
	if err := c.StoreScoped(scope, []float32{1, 0}, "answer", time.Minute); err != nil {
		t.Fatal(err)
	}
	if _, ok := c.LookupScoped(scope, []float32{1, 0}, .9); !ok {
		t.Fatal("scoped hit missing")
	}
	for _, s := range []CacheScope{{"other", "u", "p1", "k1"}, {"t", "other", "p1", "k1"}, {"t", "u", "p2", "k1"}, {"t", "u", "p1", "k2"}} {
		if _, ok := c.LookupScoped(s, []float32{1, 0}, .9); ok {
			t.Fatal("scope leak")
		}
	}
	now = now.Add(time.Minute)
	if _, ok := c.LookupScoped(scope, []float32{1, 0}, .9); ok {
		t.Fatal("expired cache returned")
	}
	if c.Store([]float32{1}, "x") == nil {
		t.Fatal("unscoped write accepted")
	}
}
