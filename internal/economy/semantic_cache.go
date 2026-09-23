package economy

import (
	"elite.local/enterprise/internal/search"
	"errors"
	"math"
	"strings"
	"sync"
	"time"
)

// CacheScope binds content to server-resolved reader and revisions. Similarity
// never authorizes access; callers must check current object authorization.
type CacheScope struct{ Tenant, Subject, PolicyVersion, KnowledgeVersion string }

func (s CacheScope) valid() bool {
	for _, v := range []string{s.Tenant, s.Subject, s.PolicyVersion, s.KnowledgeVersion} {
		if strings.TrimSpace(v) == "" || len(v) > 256 {
			return false
		}
	}
	return true
}

type semanticEntry struct {
	scope     CacheScope
	embedding []float32
	answer    string
	expires   time.Time
}
type SemanticCache struct {
	mu      sync.RWMutex
	entries []semanticEntry
	now     func() time.Time
}

func (c *SemanticCache) clock() time.Time {
	if c.now != nil {
		return c.now()
	}
	return time.Now()
}

// Deprecated: refuse unscoped writes. Use StoreScoped.
func (c *SemanticCache) Store(_ []float32, _ string) error {
	return errors.New("economy: explicit cache scope required")
}

// Deprecated: unscoped lookups always miss. Use LookupScoped.
func (c *SemanticCache) Lookup(_ []float32, _ float64) (string, bool) { return "", false }
func (c *SemanticCache) StoreScoped(scope CacheScope, embedding []float32, answer string, ttl time.Duration) error {
	if !scope.valid() || ttl <= 0 || ttl > 24*time.Hour || len(embedding) == 0 || len(embedding) > 4096 || len(answer) == 0 || len(answer) > 1<<20 {
		return errors.New("economy: invalid cache entry")
	}
	if _, err := search.CosineSimilarity(embedding, embedding); err != nil {
		return err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	now := c.clock()
	kept := c.entries[:0]
	for _, e := range c.entries {
		if now.Before(e.expires) {
			kept = append(kept, e)
		}
	}
	c.entries = kept
	if len(c.entries) >= 1024 {
		copy(c.entries, c.entries[1:])
		c.entries = c.entries[:1023]
	}
	c.entries = append(c.entries, semanticEntry{scope, append([]float32(nil), embedding...), answer, now.Add(ttl)})
	return nil
}
func (c *SemanticCache) LookupScoped(scope CacheScope, q []float32, threshold float64) (string, bool) {
	if !scope.valid() || math.IsNaN(threshold) || math.IsInf(threshold, 0) || threshold <= 0 || threshold > 1 {
		return "", false
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	now := c.clock()
	best := -1.0
	answer := ""
	found := false
	for _, e := range c.entries {
		if e.scope != scope || !now.Before(e.expires) {
			continue
		}
		v, err := search.CosineSimilarity(e.embedding, q)
		if err == nil && v >= threshold && v > best {
			best = v
			answer = e.answer
			found = true
		}
	}
	return answer, found
}
