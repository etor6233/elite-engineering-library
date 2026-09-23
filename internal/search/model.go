// Package search provides a governed, tenant-scoped search index contract:
// document validation, cosine similarity, full-text query input validation and
// a deterministic in-memory Store for tests. The authoritative full-text and
// ranking live in PostgreSQL (tsvector/ts_rank); this package supplies the
// stdlib-only core that any Store implementation shares.
package search

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strings"
)

var (
	kindRe     = regexp.MustCompile(`^[a-z][a-z0-9._-]{0,63}$`)
	hex64Re    = regexp.MustCompile(`^[0-9a-f]{64}$`)
	facetKeyRe = regexp.MustCompile(`^[a-z][a-z0-9._-]{0,63}$`)
)

// Document is an immutable searchable index entry. It is NOT the source of
// truth: the owning domain (catalog, knowledge base) owns the canonical row
// and republishes into this index.
type Document struct {
	TenantID       string
	OrganizationID string // optional; empty means tenant-scoped only
	ID             string
	Kind           string
	ExternalID     string
	Title          string
	Body           string
	Facets         map[string]string
	Embedding      []float32
	ContentSHA256  string
}

var (
	// ErrInvalidDocument reports an index entry that violates the contract.
	ErrInvalidDocument = errors.New("search: invalid document")
)

// Validate enforces the index contract: non-empty tenant/id/external/title,
// a safe kind, well-formed facet keys, finite embeddings and a 64-hex content
// digest. It fails closed on any malformed input.
func (d Document) Validate() error {
	if strings.TrimSpace(d.TenantID) == "" || len(d.TenantID) > 64 {
		return fmt.Errorf("%w: bad tenant id", ErrInvalidDocument)
	}
	if !kindRe.MatchString(d.Kind) {
		return fmt.Errorf("%w: bad kind", ErrInvalidDocument)
	}
	if strings.TrimSpace(d.ID) == "" || len(d.ID) > 128 {
		return fmt.Errorf("%w: bad id", ErrInvalidDocument)
	}
	if strings.TrimSpace(d.ExternalID) == "" || len(d.ExternalID) > 200 {
		return fmt.Errorf("%w: bad external id", ErrInvalidDocument)
	}
	if strings.TrimSpace(d.Title) == "" {
		return fmt.Errorf("%w: empty title", ErrInvalidDocument)
	}
	if len(d.OrganizationID) > 128 {
		return fmt.Errorf("%w: bad organization id", ErrInvalidDocument)
	}
	for k := range d.Facets {
		if !facetKeyRe.MatchString(k) {
			return fmt.Errorf("%w: bad facet key %q", ErrInvalidDocument, k)
		}
	}
	if d.Embedding != nil {
		if len(d.Embedding) == 0 {
			return fmt.Errorf("%w: embedding slice present but empty", ErrInvalidDocument)
		}
		for _, v := range d.Embedding {
			if math.IsNaN(float64(v)) || math.IsInf(float64(v), 0) {
				return fmt.Errorf("%w: non-finite embedding", ErrInvalidDocument)
			}
		}
	}
	if !hex64Re.MatchString(strings.ToLower(d.ContentSHA256)) {
		return fmt.Errorf("%w: content sha256 must be 64 hex chars", ErrInvalidDocument)
	}
	return nil
}
