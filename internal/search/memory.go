package search

import (
	"context"
	"errors"
	"math"
	"sort"
	"strings"
	"sync"
)

var (
	ErrNotFound = errors.New("search: not found")
	ErrDenied   = errors.New("search: tenant required")
)

// MemoryStore is a deterministic in-memory Store. Full-text relevance is a
// documented token-overlap approximation for tests; the authoritative
// full-text and ranking are PostgreSQL tsvector/ts_rank (verified in SQL).
type MemoryStore struct {
	mu   sync.RWMutex
	docs map[string]Document // key: tenantID + "\x00" + id
}

// NewMemoryStore returns an empty MemoryStore.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{docs: make(map[string]Document)}
}

func key(tenantID, id string) string { return tenantID + "\x00" + id }

// Upsert validates then inserts or replaces the tenant-scoped document.
func (m *MemoryStore) Upsert(_ context.Context, doc Document) error {
	if err := doc.Validate(); err != nil {
		return err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.docs[key(doc.TenantID, doc.ID)] = cloneDocument(doc)
	return nil
}

// Delete removes a document; it fails closed on an empty tenant.
func (m *MemoryStore) Delete(_ context.Context, tenantID, id string) error {
	if strings.TrimSpace(tenantID) == "" {
		return ErrDenied
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	k := key(tenantID, id)
	if _, ok := m.docs[k]; !ok {
		return ErrNotFound
	}
	delete(m.docs, k)
	return nil
}

// Search returns documents whose title/body overlap the query tokens, scoped
// and ranked. It denies an empty tenant and rejects an empty query.
func (m *MemoryStore) Search(_ context.Context, q Query) ([]Result, error) {
	if strings.TrimSpace(q.TenantID) == "" {
		return nil, ErrDenied
	}
	if err := ValidateQueryInput(q.Text); err != nil {
		return nil, err
	}
	if math.IsNaN(q.MinRank) || math.IsInf(q.MinRank, 0) {
		return nil, errors.New("search: invalid rank threshold")
	}
	qTokens := tokenize(q.Text)
	m.mu.RLock()
	defer m.mu.RUnlock()
	var out []Result
	for _, d := range m.docs {
		if !matchesScope(d, q.TenantID, q.OrganizationID, q.Kind, q.Facets) {
			continue
		}
		score := overlapScore(qTokens, tokenize(d.Title), tokenize(d.Body))
		if score <= 0 || score < q.MinRank {
			continue
		}
		out = append(out, Result{Document: cloneDocument(d), Score: score})
	}
	sortResults(out)
	out = limitResults(out, q.Limit)
	return out, nil
}

// VectorSearch returns the top-K documents by cosine similarity.
func (m *MemoryStore) VectorSearch(_ context.Context, q VectorQuery) ([]Result, error) {
	if strings.TrimSpace(q.TenantID) == "" {
		return nil, ErrDenied
	}
	if len(q.Embedding) > 4096 || math.IsNaN(q.MinSimilarity) || math.IsInf(q.MinSimilarity, 0) || q.MinSimilarity < -1 || q.MinSimilarity > 1 {
		return nil, errors.New("search: invalid vector query")
	}
	if _, err := CosineSimilarity(q.Embedding, q.Embedding); err != nil {
		return nil, err
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	var out []Result
	for _, d := range m.docs {
		if !matchesScope(d, q.TenantID, q.OrganizationID, q.Kind, nil) {
			continue
		}
		if d.Embedding == nil {
			continue
		}
		sim, err := CosineSimilarity(d.Embedding, q.Embedding)
		if err != nil {
			continue // dim mismatch or zero vector: skip, never invent a score
		}
		if sim < q.MinSimilarity {
			continue
		}
		out = append(out, Result{Document: cloneDocument(d), Score: sim})
	}
	sortResults(out)
	topK := q.TopK
	if topK <= 0 {
		topK = 10
	}
	topK = ClampTopK(topK)
	if topK < len(out) {
		out = out[:topK]
	}
	return out, nil
}

func matchesScope(d Document, tenantID, orgID, kind string, facets map[string]string) bool {
	if d.TenantID != tenantID {
		return false
	}
	if orgID != "" && d.OrganizationID != orgID {
		return false
	}
	if kind != "" && d.Kind != kind {
		return false
	}
	for fk, fv := range facets {
		if d.Facets[fk] != fv {
			return false
		}
	}
	return true
}

func overlapScore(query, title, body []string) float64 {
	if len(query) == 0 {
		return 0
	}
	var score float64
	for _, qt := range query {
		if contains(title, qt) {
			score += 2
		}
		if contains(body, qt) {
			score++
		}
	}
	return score
}

func contains(list []string, s string) bool {
	for _, x := range list {
		if x == s {
			return true
		}
	}
	return false
}

func sortResults(out []Result) {
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Score != out[j].Score {
			return out[i].Score > out[j].Score
		}
		return out[i].Document.ID < out[j].Document.ID
	})
}

func limitResults(out []Result, limit int) []Result {
	if limit <= 0 {
		limit = ClampTopK(10)
	} else {
		limit = ClampTopK(limit)
	}
	if limit < len(out) {
		return out[:limit]
	}
	return out
}

// Detach mutable members at both sides of the Store boundary.
func cloneDocument(d Document) Document {
	if d.Facets != nil {
		f := make(map[string]string, len(d.Facets))
		for k, v := range d.Facets {
			f[k] = v
		}
		d.Facets = f
	}
	if d.Embedding != nil {
		d.Embedding = append([]float32{}, d.Embedding...)
	}
	return d
}
