# Go Search Core

## 1. Metadata

```yaml
pack_id: "GO-SEARCH-CORE"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa un índice de búsqueda gobernado por tenant: full-text PostgreSQL (tsvector/GIN + websearch_to_tsquery + ts_rank) con filtros de kind/facets, almacenamiento de embeddings real[] y búsqueda vectorial por coseno, con core Go stdlib-only y negación por defecto de tenant."
stacks: ["Go 1.26.7", "PostgreSQL 18.6"]
compatible_with: ["GO-ENTERPRISE-BACKEND 0.4.x", "PG-TX-FOUNDATION 0.1.x", "GO-ML-AI-FOUNDATION 0.1.x"]
incompatible_with: ["full-text sin configuración fija", "búsqueda sin tenant_id obligatorio", "configuración de idioma/ranking no declarada"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://go.dev", "https://www.postgresql.org/docs/18/textsearch.html"]
verified_at: "2026-09-02"
```

Este pack no afirma recall universal de un modelo de embeddings ni escalado a corpus masivo. El full-text usa la configuración `simple` (sin stemming) de forma fija y declarada. La búsqueda vectorial es una exploración lineal por coseno honesta para corpus acotados; para escala se fijan Faiss/Meta `7059eaf7da7eddda62e71367e684d4bdedd7f94f` y DiskANN/Microsoft `860cf47bc11b6c2b938818773831a9374553d976` (ambos MIT) como backends de retrieval autorizados, que se materializan sólo si el workload lo justifica.

## 2. Applicability

Use este pack cuando una franquicia necesite búsqueda de catálogo/productos (full-text + facets) o recuperación semántica sobre una base de conocimiento (embeddings), siempre aislada por tenant y opcionalmente por organización. Es consumido por el portal web (búsqueda de catálogo) y por el agente conversacional (recuperación de conocimiento).

Rechace este pack para: corpus que exija ANN a gran escala sin admitir Faiss/DiskANN; idiomas que requieran stemming/ranking propio sin configurarlo; o búsqueda sin un tenant obligatorio. El índice `search.document` es una proyección derivada: el source of truth del catálogo/conocimiento sigue en su owner de dominio.

## 3. Architecture contract

- **Ownership**: el paquete `internal/search` gobierna validación de documento, similitud coseno, validación de query y el borde `Store`. La migración `0040` gobierna el índice persistente `search.document`. El dominio propietario (catálogo/conocimiento) republica en el índice.
- **Invariantes**: (1) `tenant_id` es obligatorio en toda búsqueda (deny-by-default). (2) El índice no es source of truth; se reconstruye desde el owner. (3) El embedding es `real[]` de dimensión consistente por `kind`; dimensión distinta o vector cero no produce score inventado. (4) El `tsv` se genera con configuración `simple` fija (title peso A, body peso B).
- **Data flow**: `Store.Upsert(validado)` → fila inmutable indexada; `Store.Search` → `websearch_to_tsquery` sobre `tsv` con filtros tenant/org/kind/facets + `ts_rank`; `Store.VectorSearch` → coseno Go sobre embeddings del tenant.
- **Failure modes**: tenant vacío → `ErrDenied`; query vacía/sobrelarga → `ErrInvalidQuery`; documento inválido → `ErrInvalidDocument`; embedding de dimensión distinta se omite (nunca se inventa score).
- **Seguridad/privacidad**: aislamiento por `tenant_id` y `organization_id` en toda query; el store nunca amplía el scope solicitado.
- **Performance budget**: core Go sin I/O; el full-text delega en índices GIN de PostgreSQL; la búsqueda vectorial es lineal y se sustituye por Faiss/DiskANN sólo ante evidencia de workload.
- **Operación/migración/rollback**: migración `0040` up/down atómica; el índice se puede reconstruir borrando y republicando desde el owner de dominio.

## 4. Exact file manifest

```text
CREATE internal/search/model.go
CREATE internal/search/cosine.go
CREATE internal/search/fulltext.go
CREATE internal/search/store.go
CREATE internal/search/memory.go
CREATE internal/search/model_test.go
CREATE internal/search/cosine_test.go
CREATE internal/search/memory_test.go
CREATE db/migrations/0040_search.up.sql
CREATE db/migrations/0040_search.down.sql
CREATE db/tests/0040_search.test.sql
```

## 5. Materialization blocks

### FILE: `internal/search/model.go`
```yaml
block_id: "GO-SEARCH-CORE:internal/search/model.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "955eb3492c284cbcc69590c7a98a3398b8d5d34ded6e7c831976cea328cda901"
variables: []
secrets_allowed: false
```
````go
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
````

### FILE: `internal/search/cosine.go`
```yaml
block_id: "GO-SEARCH-CORE:internal/search/cosine.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "ed136fa4504cf069a83f60cc07bf35ad094ca405c68dfb542f65a6ae8222312c"
variables: []
secrets_allowed: false
```
````go
package search

import (
	"errors"
	"fmt"
	"math"
)

var (
	ErrEmptyEmbedding = errors.New("search: empty embedding")
	ErrDimMismatch    = errors.New("search: embedding dimension mismatch")
	ErrZeroVector     = errors.New("search: zero vector")
)

// CosineSimilarity returns the cosine of two float32 embeddings. Equal,
// non-empty dimensions are required; a zero vector yields an error (the
// similarity is undefined) rather than an invented score.
func CosineSimilarity(a, b []float32) (float64, error) {
	if len(a) == 0 || len(b) == 0 {
		return 0, ErrEmptyEmbedding
	}
	if len(a) != len(b) {
		return 0, fmt.Errorf("%w: %d vs %d", ErrDimMismatch, len(a), len(b))
	}
	var dot, na, nb float64
	for i := range a {
		av := float64(a[i])
		bv := float64(b[i])
		dot += av * bv
		na += av * av
		nb += bv * bv
	}
	if na == 0 || nb == 0 {
		return 0, ErrZeroVector
	}
	return dot / (math.Sqrt(na) * math.Sqrt(nb)), nil
}
````

### FILE: `internal/search/fulltext.go`
```yaml
block_id: "GO-SEARCH-CORE:internal/search/fulltext.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "0018e2cc8c396a49e0390cd99f4a12a73775e7a457f9702465fd26a4c5ea8a38"
variables: []
secrets_allowed: false
```
````go
package search

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

var (
	ErrInvalidQuery = errors.New("search: invalid query")

	maxQueryRunes = 500
	maxTopK       = 100
)

// ValidateQueryInput rejects empty, whitespace-only or overlong query text.
// The accepted string is intended for PostgreSQL websearch_to_tsquery, which
// parses a user-friendly subset and never treats input as SQL.
func ValidateQueryInput(text string) error {
	t := strings.TrimSpace(text)
	if t == "" {
		return fmt.Errorf("%w: empty", ErrInvalidQuery)
	}
	if len([]rune(t)) > maxQueryRunes {
		return fmt.Errorf("%w: too long", ErrInvalidQuery)
	}
	return nil
}

// ClampTopK bounds a caller-provided limit into the safe range.
func ClampTopK(k int) int {
	if k < 1 {
		return 1
	}
	if k > maxTopK {
		return maxTopK
	}
	return k
}

// tokenize splits text into lowercase word tokens for the in-memory index.
func tokenize(text string) []string {
	return strings.FieldsFunc(strings.ToLower(text), func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
}
````

### FILE: `internal/search/store.go`
```yaml
block_id: "GO-SEARCH-CORE:internal/search/store.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "6c080f74820ccfa0679786706d2d13d04f4a9c9f368f4e2ac0abee7fcd8a5a44"
variables: []
secrets_allowed: false
```
````go
package search

import "context"

// Store is the search index boundary. A PostgreSQL implementation backs it
// with tsvector full-text and real[] embeddings; MemoryStore is a
// deterministic implementation for tests and small tenants.
type Store interface {
	Upsert(ctx context.Context, doc Document) error
	Delete(ctx context.Context, tenantID, id string) error
	Search(ctx context.Context, q Query) ([]Result, error)
	VectorSearch(ctx context.Context, q VectorQuery) ([]Result, error)
}

// Query is a full-text search scoped by tenant (mandatory) and optionally by
// organization, kind and facets. OrganizationID is enforced deny-by-default by
// the caller: the store never widens scope.
type Query struct {
	TenantID       string
	OrganizationID string
	Kind           string
	Text           string
	Facets         map[string]string
	Limit          int
	MinRank        float64
}

// VectorQuery is a cosine similarity search over embeddings.
type VectorQuery struct {
	TenantID       string
	OrganizationID string
	Kind           string
	Embedding      []float32
	TopK           int
	MinSimilarity  float64
}

// Result pairs a document with its relevance score (rank or similarity).
type Result struct {
	Document Document
	Score    float64
}
````

### FILE: `internal/search/memory.go`
```yaml
block_id: "GO-SEARCH-CORE:internal/search/memory.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "765217bf4874b67cdf23bbb8bc6bc01ac2573c1a1f80b736a010951050c91609"
variables: []
secrets_allowed: false
```
````go
package search

import (
	"context"
	"errors"
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
	m.docs[key(doc.TenantID, doc.ID)] = doc
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
		out = append(out, Result{Document: d, Score: score})
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
	if len(q.Embedding) == 0 {
		return nil, ErrEmptyEmbedding
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
		out = append(out, Result{Document: d, Score: sim})
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
````

### FILE: `internal/search/model_test.go`
```yaml
block_id: "GO-SEARCH-CORE:internal/search/model_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "ba2875badd4e8608976b41071704af391c129b9bffa78a96e9d4db33cb55f79b"
variables: []
secrets_allowed: false
```
````go
package search

import (
	"errors"
	"math"
	"strings"
	"testing"
)

func doc(tenant, org, id, kind, title, body string) Document {
	return Document{
		TenantID:       tenant,
		OrganizationID: org,
		ID:             id,
		Kind:           kind,
		ExternalID:     "ext-" + id,
		Title:          title,
		Body:           body,
		Facets:         map[string]string{},
		ContentSHA256:  strings.Repeat("a", 64),
	}
}

func TestDocumentValidate(t *testing.T) {
	valid := doc("t-1", "org-1", "d1", "product", "Scooter Urbano", "movilidad electrica")
	if err := valid.Validate(); err != nil {
		t.Fatalf("valid document rejected: %v", err)
	}

	bad := doc("t-1", "", "d1", "Product", "x", "y") // kind uppercase
	if err := bad.Validate(); !errors.Is(err, ErrInvalidDocument) {
		t.Fatalf("bad kind accepted: %v", err)
	}

	empty := doc("t-1", "", "d1", "product", "  ", "y")
	if err := empty.Validate(); !errors.Is(err, ErrInvalidDocument) {
		t.Fatalf("empty title accepted: %v", err)
	}

	badSHA := doc("t-1", "", "d1", "product", "x", "y")
	badSHA.ContentSHA256 = "zz"
	if err := badSHA.Validate(); !errors.Is(err, ErrInvalidDocument) {
		t.Fatalf("bad content sha accepted: %v", err)
	}

	nan := doc("t-1", "", "d1", "product", "x", "y")
	nan.Embedding = []float32{1.0, float32(math.NaN())} // NaN
	if err := nan.Validate(); !errors.Is(err, ErrInvalidDocument) {
		t.Fatalf("non-finite embedding accepted: %v", err)
	}
}

func TestDocumentValidateEmbedding(t *testing.T) {
	d := doc("t-1", "", "d1", "product", "x", "y")
	d.Embedding = []float32{}
	if err := d.Validate(); !errors.Is(err, ErrInvalidDocument) {
		t.Fatalf("empty embedding slice accepted: %v", err)
	}
	d.Embedding = []float32{1.0, 0.0}
	if err := d.Validate(); err != nil {
		t.Fatalf("finite embedding rejected: %v", err)
	}
}
````

### FILE: `internal/search/cosine_test.go`
```yaml
block_id: "GO-SEARCH-CORE:internal/search/cosine_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "5c86909c0615ea7a0b545324b7c04a8ec9c5465883534998ca553cdf56189be7"
variables: []
secrets_allowed: false
```
````go
package search

import (
	"errors"
	"math"
	"testing"
)

func TestCosineSimilarity(t *testing.T) {
	if s, err := CosineSimilarity([]float32{1, 0}, []float32{1, 0}); err != nil || math.Abs(s-1.0) > 1e-6 {
		t.Fatalf("identical vectors: got %v, %v", s, err)
	}
	if s, err := CosineSimilarity([]float32{1, 0}, []float32{0, 1}); err != nil || math.Abs(s-0.0) > 1e-6 {
		t.Fatalf("orthogonal vectors: got %v, %v", s, err)
	}
	if s, err := CosineSimilarity([]float32{1, 0}, []float32{-1, 0}); err != nil || math.Abs(s+1.0) > 1e-6 {
		t.Fatalf("opposite vectors: got %v, %v", s, err)
	}
}

func TestCosineDimMismatch(t *testing.T) {
	if _, err := CosineSimilarity([]float32{1, 0, 0}, []float32{1, 0}); !errors.Is(err, ErrDimMismatch) {
		t.Fatalf("dim mismatch accepted: %v", err)
	}
}

func TestCosineZeroVector(t *testing.T) {
	if _, err := CosineSimilarity([]float32{0, 0}, []float32{1, 0}); !errors.Is(err, ErrZeroVector) {
		t.Fatalf("zero vector accepted: %v", err)
	}
}

func TestCosineEmpty(t *testing.T) {
	if _, err := CosineSimilarity(nil, []float32{1}); !errors.Is(err, ErrEmptyEmbedding) {
		t.Fatalf("empty embedding accepted: %v", err)
	}
}
````

### FILE: `internal/search/memory_test.go`
```yaml
block_id: "GO-SEARCH-CORE:internal/search/memory_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "a3a5fad0ab1530699fb2fdbd970ff70ee8c4dfcaab73239fc40142a538c7f0cc"
variables: []
secrets_allowed: false
```
````go
package search

import (
	"context"
	"errors"
	"testing"
)

func TestSearchRankingAndScope(t *testing.T) {
	s := NewMemoryStore()
	a := doc("t-1", "org-1", "d1", "product", "Scooter Urbano Pro", "scooter electrico urbano 350W")
	a.Facets = map[string]string{"category": "movilidad"}
	b := doc("t-1", "org-1", "d2", "product", "Scooter Plegable", "scooter plegable ligero")
	c := doc("t-1", "org-2", "d3", "product", "Scooter Urbano", "scooter urbano") // other org
	d := doc("t-2", "org-1", "d4", "product", "Scooter Urbano", "scooter urbano")  // other tenant
	for _, x := range []Document{a, b, c, d} {
		if err := s.Upsert(context.Background(), x); err != nil {
			t.Fatal(err)
		}
	}

	res, err := s.Search(context.Background(), Query{TenantID: "t-1", OrganizationID: "org-1", Text: "scooter urbano"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 2 {
		t.Fatalf("expected 2 results scoped to org-1, got %d", len(res))
	}
	if res[0].Document.ID != "d1" {
		t.Fatalf("expected d1 ranked first (title hit), got %s", res[0].Document.ID)
	}
}

func TestSearchFacetAndKindFilter(t *testing.T) {
	s := NewMemoryStore()
	a := doc("t-1", "", "d1", "product", "Casco Integral", "casco seguridad")
	a.Facets = map[string]string{"category": "accesorio"}
	b := doc("t-1", "", "d2", "product", "Casco Abierto", "casco")
	b.Facets = map[string]string{"category": "accesorio"}
	c := doc("t-1", "", "d3", "faq", "Casco", "como elegir casco")
	for _, x := range []Document{a, b, c} {
		if err := s.Upsert(context.Background(), x); err != nil {
			t.Fatal(err)
		}
	}

	res, err := s.Search(context.Background(), Query{TenantID: "t-1", Kind: "product", Text: "casco"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 2 {
		t.Fatalf("kind filter failed: got %d", len(res))
	}

	res, err = s.Search(context.Background(), Query{TenantID: "t-1", Text: "casco", Facets: map[string]string{"category": "accesorio"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 2 {
		t.Fatalf("facet filter failed: got %d", len(res))
	}
}

func TestSearchDeniesEmptyTenantAndQuery(t *testing.T) {
	s := NewMemoryStore()
	if _, err := s.Search(context.Background(), Query{TenantID: "", Text: "x"}); !errors.Is(err, ErrDenied) {
		t.Fatalf("empty tenant accepted: %v", err)
	}
	if _, err := s.Search(context.Background(), Query{TenantID: "t-1", Text: "   "}); !errors.Is(err, ErrInvalidQuery) {
		t.Fatalf("empty query accepted: %v", err)
	}
}

func TestVectorSearchTopKAndThreshold(t *testing.T) {
	s := NewMemoryStore()
	a := doc("t-1", "", "d1", "product", "A", "x")
	a.Embedding = []float32{1, 0, 0}
	b := doc("t-1", "", "d2", "product", "B", "x")
	b.Embedding = []float32{0, 1, 0}
	c := doc("t-1", "", "d3", "product", "C", "x")
	c.Embedding = []float32{1, 0, 0}
	for _, x := range []Document{a, b, c} {
		if err := s.Upsert(context.Background(), x); err != nil {
			t.Fatal(err)
		}
	}

	res, err := s.VectorSearch(context.Background(), VectorQuery{
		TenantID: "t-1", Embedding: []float32{1, 0, 0}, TopK: 2, MinSimilarity: 0.99,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 2 {
		t.Fatalf("expected top-2, got %d", len(res))
	}
	if res[0].Score < 0.999 || res[1].Score < 0.999 {
		t.Fatalf("expected near-1 similarities, got %v %v", res[0].Score, res[1].Score)
	}
}

func TestDelete(t *testing.T) {
	s := NewMemoryStore()
	a := doc("t-1", "", "d1", "product", "A", "x")
	if err := s.Upsert(context.Background(), a); err != nil {
		t.Fatal(err)
	}
	if err := s.Delete(context.Background(), "t-1", "d1"); err != nil {
		t.Fatal(err)
	}
	if err := s.Delete(context.Background(), "t-1", "d1"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected not found after delete, got %v", err)
	}
	if err := s.Delete(context.Background(), "", "d1"); !errors.Is(err, ErrDenied) {
		t.Fatalf("empty tenant delete accepted: %v", err)
	}
}
````

### FILE: `db/migrations/0040_search.up.sql`
```yaml
block_id: "GO-SEARCH-CORE:db/migrations/0040_search.up.sql:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "1c407427f5d9cbeae0d800e50db3edd208e91030695fabf9c2945c3dbf43c8d1"
variables: []
secrets_allowed: false
```
````sql
begin;

create schema if not exists search;

create table search.document (
  tenant_id uuid not null,
  organization_id text,
  document_id text not null,
  kind text not null,
  external_id text not null,
  title text not null,
  body text not null default '',
  facets jsonb not null default '{}'::jsonb,
  embedding real[],
  content_sha256 text not null,
  tsv tsvector generated always as (
    setweight(to_tsvector('simple', title), 'A') ||
    setweight(to_tsvector('simple', body), 'B')
  ) stored,
  indexed_at timestamptz not null default clock_timestamp(),
  updated_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, document_id),
  foreign key (tenant_id) references platform.tenant (tenant_id),
  check (length(document_id) between 1 and 128),
  check (length(external_id) between 1 and 200),
  check (title <> ''),
  check (kind ~ '^[a-z][a-z0-9._-]{0,63}$'),
  check (organization_id is null or length(organization_id) between 1 and 128),
  check (jsonb_typeof(facets) = 'object'),
  check (embedding is null or cardinality(embedding) > 0),
  check (content_sha256 ~ '^[0-9a-f]{64}$'),
  check (updated_at >= indexed_at)
);

create index search_document_tsv_idx on search.document using gin (tsv);
create index search_document_scope_idx on search.document (tenant_id, kind, external_id);
create index search_document_facets_idx on search.document using gin (facets);

commit;
````

### FILE: `db/migrations/0040_search.down.sql`
```yaml
block_id: "GO-SEARCH-CORE:db/migrations/0040_search.down.sql:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "a6ed0dbcd81c88d0dcc541d3830cc94ee231778d54a9f503cf6cac194889cd52"
variables: []
secrets_allowed: false
```
````sql
begin;

drop table if exists search.document;
drop schema if exists search;

commit;
````

### FILE: `db/tests/0040_search.test.sql`
```yaml
block_id: "GO-SEARCH-CORE:db/tests/0040_search.test.sql:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "513527780cf4956f73f4dd236e78128cc5e31ca4b586908e481d12adfb767995"
variables: []
secrets_allowed: false
```
````sql
-- 0040_search.test.sql — verifica búsqueda full-text (tsvector/GIN + websearch_to_tsquery),
-- ranking por peso, aislamiento ACL por tenant/organization, filtro de facets y
-- almacenamiento round-trip de embeddings real[].
-- Precondición: migración 0001 (platform.tenant) y 0040 (search.document) aplicadas.

begin;

insert into platform.tenant (tenant_id, tenant_code, legal_name, display_name) values
  ('11111111-1111-1111-1111-111111111111', 'tenant-a', 'A', 'Tenant A'),
  ('22222222-2222-2222-2222-222222222222', 'tenant-b', 'B', 'Tenant B');

insert into search.document
  (tenant_id, organization_id, document_id, kind, external_id, title, body, facets, embedding, content_sha256)
values
  ('11111111-1111-1111-1111-111111111111', 'org-1', 'd1', 'product', 'SKU-1',
   'Scooter Urbano Pro', 'vehiculo electrico ligero', '{"category":"movilidad"}',
   ARRAY[1.0,0.0,0.0]::real[], repeat('a',64)),
  ('11111111-1111-1111-1111-111111111111', 'org-1', 'd2', 'product', 'SKU-2',
   'Vehiculo Plegable', 'scooter plegable portatil', '{"category":"movilidad"}',
   ARRAY[0.0,1.0,0.0]::real[], repeat('b',64)),
  ('11111111-1111-1111-1111-111111111111', 'org-2', 'd3', 'product', 'SKU-3',
   'Scooter Urbano', 'scooter urbano', '{"category":"movilidad"}',
   null, repeat('c',64)),
  ('22222222-2222-2222-2222-222222222222', 'org-1', 'd4', 'product', 'SKU-4',
   'Scooter Urbano', 'scooter urbano', '{"category":"movilidad"}',
   null, repeat('d',64));

-- full-text 'scooter' en (tenant-a, org-1): 2 resultados, d1 (title, peso A) primero
do $$
declare n int; top text;
begin
  select count(*) into n from search.document
   where tenant_id = '11111111-1111-1111-1111-111111111111'
     and organization_id = 'org-1'
     and tsv @@ websearch_to_tsquery('simple', 'scooter');
  if n <> 2 then raise exception 'fulltext expected 2, got %', n; end if;

  select document_id into top from search.document
   where tenant_id = '11111111-1111-1111-1111-111111111111'
     and organization_id = 'org-1'
     and tsv @@ websearch_to_tsquery('simple', 'scooter')
   order by ts_rank(tsv, websearch_to_tsquery('simple', 'scooter')) desc, document_id
   limit 1;
  if top <> 'd1' then raise exception 'expected d1 top-ranked (title weight A), got %', top; end if;
end $$;

-- término específico 'urbano' en (tenant-a, org-1): sólo d1
do $$
declare n int;
begin
  select count(*) into n from search.document
   where tenant_id = '11111111-1111-1111-1111-111111111111'
     and organization_id = 'org-1'
     and tsv @@ websearch_to_tsquery('simple', 'urbano');
  if n <> 1 then raise exception 'urbano expected 1, got %', n; end if;
end $$;

-- aislamiento de tenant: tenant-b ve 1, tenant-a ve 0 de tenant-b
do $$
declare a int; b int;
begin
  select count(*) into b from search.document
   where tenant_id = '22222222-2222-2222-2222-222222222222'
     and tsv @@ websearch_to_tsquery('simple', 'scooter');
  if b <> 1 then raise exception 'tenant-b expected 1, got %', b; end if;
  select count(*) into a from search.document
   where tenant_id = '11111111-1111-1111-1111-111111111111'
     and tsv @@ websearch_to_tsquery('simple', 'scooter')
     and organization_id = 'org-2';
  if a <> 1 then raise exception 'org-2 expected 1, got %', a; end if;
end $$;

-- filtro de facets
do $$
declare n int;
begin
  select count(*) into n from search.document
   where tenant_id = '11111111-1111-1111-1111-111111111111'
     and facets @> '{"category":"movilidad"}';
  if n <> 3 then raise exception 'facet expected 3, got %', n; end if;
end $$;

-- round-trip de embedding real[]
do $$
declare dim int; v real;
begin
  select array_length(embedding,1), embedding[1] into dim, v from search.document
   where tenant_id = '11111111-1111-1111-1111-111111111111' and document_id = 'd1';
  if dim <> 3 or v <> 1.0 then raise exception 'embedding round-trip failed: dim=%, v=%', dim, v; end if;
end $$;

rollback;
````


## 6. Configuration surface

Sin variables ni secretos. La configuración de idioma (`simple`) y el ranking (`ts_rank`) son fijos y declarados; un cambio de idioma o modelo de embeddings es una nueva versión del pack, no una variable de runtime.

## 7. Dependency bill

| Package/image/tool | Pin exacto | Uso | Licencia | Runtime/build | Fuente oficial |
|---|---|---|---|---|---|
| Go standard library | go 1.26.7 (toolchain) | core Go | BSD-3-Clause | runtime | https://go.dev |
| PostgreSQL 18.6 | 18.6 | full-text/índices/embeddings | PostgreSQL License | runtime | https://www.postgresql.org |

Faiss/Meta y DiskANN/Microsoft permanecen como backends de retrieval autorizados (MIT) no materializados en este pack; se adquieren por expediente separado si el workload lo exige.

## 8. Apply order

1. Componer el backend (`GO-ENTERPRISE-BACKEND 0.4.x`) o disponer del módulo `elite.local/enterprise` con `go 1.26`.
2. Aplicar migración `0040` sobre PostgreSQL 18.6 con la cadena `0001` ya aplicada (`platform.tenant`).
3. Colocar los ocho archivos Go bajo `internal/search/` y los tres SQL bajo `db/`.
4. Verificar: `psql ... -v ON_ERROR_STOP=1 -f db/tests/0040_search.test.sql` y `go test ./... -count=1`.
5. Rollback: migración `0040` down y eliminar `internal/search/`; el índice se reconstruye desde el owner.

## 9. Verification

- `go test ./internal/search/ -count=1`: 11/11 PASS (validación de documento con negativos, coseno con ortogonal/idéntico/opuesto/dim/cero, ranking por peso title>body, aislamiento tenant/org, filtro kind/facets, top-K y umbral vectorial, delete con not-found/denied).
- `go test ./... -count=1` (aifoundation + search): PASS.
- `go vet ./...`: exit 0.
- PostgreSQL 18.6 real (initdb → start → up → test → down): `ON_ERROR_STOP=1` exit 0; 5 aserciones DO pasan (full-text 2 resultados con d1 top por peso A, término específico, aislamiento tenant/org, filtro facets, round-trip embedding real[] dim 3); down deja `search_ns_null=true`.

## 10. Reconstruction evidence

Véase `reconstruction_evidence/GO_SEARCH_CORE_2026-09-02_V176.md`.
