# Go Help Center Core

## 1. Metadata

```yaml
pack_id: "GO-HELP-CENTER-CORE"
pack_version: "0.1.2"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Referencia AUTHORED de artículos en memoria por tupla tenant/id, concurrencia optimista acotada a int64 y estados draft→published→archived; búsqueda local substring, sin adapter externo ni autorización del actor."
stacks: ["Go 1.26.8"]
compatible_with: ["uso aislado; integración con backend/buscador aún no admitida"]
incompatible_with: ["artículo sin título/cuerpo", "versión obsoleta", "transición inválida"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://go.dev"]
verified_at: "2026-09-10"
```

## 2. Applicability

Referencia local de artículos y estados: no ofrece UI, autenticación, autorización, persistencia ni integración con buscador. Las transiciones sólo cambian estado en memoria.

## 3. Architecture contract

- **Ownership**: `internal/helpcenter` gobierna artículos y publicación; la búsqueda actual recorre artículos publicados con substring case-insensitive, sin conexión a `GO-SEARCH-CORE`.
- **Invariantes**: (1) título/cuerpo no vacíos, categoría segura. (2) versionado con concurrencia optimista. (3) estados draft→published→archived. (4) búsqueda sólo publicados. (5) tenant-scoped.
- **Data flow**: `Create` → `Update`/`Publish`/`Archive` (version-gated) → `Published`/`Search`.
- **Failure modes**: inválido, duplicado, versión obsoleta, transición inválida → error.
- **Seguridad/privacidad**: tenant-scoped.
- **Performance budget**: lookup/mutación por map; validación depende del texto. Published recorre N artículos y Search también el texto publicado; sin índice ni orden garantizado.

## 4. Exact file manifest

```text
CREATE internal/helpcenter/helpcenter.go
CREATE internal/helpcenter/helpcenter_test.go
```

## 5. Materialization blocks

### FILE: `internal/helpcenter/helpcenter.go`
```yaml
block_id: "GO-HELP-CENTER-CORE:internal/helpcenter/helpcenter.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "5c8e958ab4ef1262fbfed2583b99f23d50cce20f27715031cb347e2ac47cd844"
variables: []
secrets_allowed: false
```
````go
// Package helpcenter provides a tenant-scoped knowledge base: articles with
// categories, versioned updates (optimistic concurrency) and a state machine
// draft → published → archived. Search is a local substring scan.
package helpcenter

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strings"
	"sync"
	"time"
)

// State is the publication lifecycle.
type State string

const (
	StateDraft     State = "draft"
	StatePublished State = "published"
	StateArchived  State = "archived"
)

var (
	categoryRe = regexp.MustCompile(`^[a-z][a-z0-9_-]{0,63}$`)

	ErrInvalidArticle   = errors.New("helpcenter: invalid article")
	ErrDuplicate        = errors.New("helpcenter: duplicate article")
	ErrNotFound         = errors.New("helpcenter: not found")
	ErrVersion          = errors.New("helpcenter: version conflict")
	ErrBadTransition    = errors.New("helpcenter: bad state transition")
	ErrVersionExhausted = errors.New("helpcenter: version exhausted")
)

// Article is a single help document.
type Article struct {
	TenantID  string
	ID        string
	Category  string
	Title     string
	Body      string
	State     State
	Version   int64
	UpdatedAt time.Time
}

// Store holds articles.
type Store struct {
	mu       sync.Mutex
	articles map[articleKey]Article
}

// NewStore returns an empty store.
func NewStore() *Store {
	return &Store{articles: make(map[articleKey]Article)}
}

type articleKey struct{ tenant, id string }

func key(tenant, id string) articleKey { return articleKey{tenant, id} }

func validate(a Article) error {
	if strings.TrimSpace(a.TenantID) == "" || strings.TrimSpace(a.ID) == "" {
		return ErrInvalidArticle
	}
	if !categoryRe.MatchString(a.Category) {
		return ErrInvalidArticle
	}
	if strings.TrimSpace(a.Title) == "" || strings.TrimSpace(a.Body) == "" {
		return ErrInvalidArticle
	}
	return nil
}

// Create adds a draft article (version 1).
func (s *Store) Create(a Article) error {
	if err := validate(a); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.articles[key(a.TenantID, a.ID)]; ok {
		return ErrDuplicate
	}
	a.State = StateDraft
	a.Version = 1
	a.UpdatedAt = time.Now().UTC()
	s.articles[key(a.TenantID, a.ID)] = a
	return nil
}

// Update edits title/body with optimistic concurrency (expected version).
func (s *Store) Update(tenant, id string, title, body string, expectedVersion int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	a, ok := s.articles[key(tenant, id)]
	if !ok {
		return ErrNotFound
	}
	if a.Version != expectedVersion {
		return ErrVersion
	}
	if strings.TrimSpace(title) == "" || strings.TrimSpace(body) == "" {
		return ErrInvalidArticle
	}
	if a.Version == math.MaxInt64 {
		return ErrVersionExhausted
	}
	a.Title = title
	a.Body = body
	a.Version++
	a.UpdatedAt = time.Now().UTC()
	s.articles[key(tenant, id)] = a
	return nil
}

// Publish marks draft → published (version-gated).
func (s *Store) Publish(tenant, id string, expectedVersion int64) error {
	return s.transition(tenant, id, expectedVersion, StateDraft, StatePublished)
}

// Archive marks published → archived (version-gated).
func (s *Store) Archive(tenant, id string, expectedVersion int64) error {
	return s.transition(tenant, id, expectedVersion, StatePublished, StateArchived)
}

func (s *Store) transition(tenant, id string, expectedVersion int64, from, to State) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	a, ok := s.articles[key(tenant, id)]
	if !ok {
		return ErrNotFound
	}
	if a.Version != expectedVersion {
		return ErrVersion
	}
	if a.State != from {
		return ErrBadTransition
	}
	if a.Version == math.MaxInt64 {
		return ErrVersionExhausted
	}
	a.State = to
	a.Version++
	a.UpdatedAt = time.Now().UTC()
	s.articles[key(tenant, id)] = a
	return nil
}

// Published returns tenant-scoped published articles, optionally by category.
func (s *Store) Published(tenant, category string) []Article {
	s.mu.Lock()
	defer s.mu.Unlock()
	var out []Article
	for _, a := range s.articles {
		if a.TenantID == tenant && a.State == StatePublished {
			if category == "" || a.Category == category {
				out = append(out, a)
			}
		}
	}
	return out
}

// Search returns published articles whose title/body contain the query
// (case-insensitive). No external search adapter is connected.
func (s *Store) Search(tenant, query string) []Article {
	q := strings.ToLower(strings.TrimSpace(query))
	if q == "" {
		return nil
	}
	var out []Article
	for _, a := range s.Published(tenant, "") {
		if strings.Contains(strings.ToLower(a.Title), q) || strings.Contains(strings.ToLower(a.Body), q) {
			out = append(out, a)
		}
	}
	return out
}

// String aids debugging.
func (s *Store) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return fmt.Sprintf("helpcenter(%d)", len(s.articles))
}
````

### FILE: `internal/helpcenter/helpcenter_test.go`
```yaml
block_id: "GO-HELP-CENTER-CORE:internal/helpcenter/helpcenter_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "86b8d27917677867cfdc5f9c7536e0bfb4ddb3196d3756825c046861246ecbc6"
variables: []
secrets_allowed: false
```
````go
package helpcenter

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
)

func art() Article {
	return Article{TenantID: "t", ID: "a1", Category: "faq", Title: "Cómo reservar", Body: "Elegí un turno disponible."}
}

func TestCreatePublishArchive(t *testing.T) {
	s := NewStore()
	if err := s.Create(art()); err != nil {
		t.Fatal(err)
	}
	if len(s.Published("t", "")) != 0 {
		t.Fatal("draft should not be published")
	}
	if err := s.Publish("t", "a1", 1); err != nil {
		t.Fatal(err)
	}
	if len(s.Published("t", "")) != 1 {
		t.Fatal("published should be visible")
	}
	if err := s.Archive("t", "a1", 2); err != nil {
		t.Fatal(err)
	}
	if len(s.Published("t", "")) != 0 {
		t.Fatal("archived should not be published")
	}
}

func TestOptimisticConcurrency(t *testing.T) {
	s := NewStore()
	_ = s.Create(art())
	if err := s.Update("t", "a1", "Nuevo", "Contenido", 99); !errors.Is(err, ErrVersion) {
		t.Fatalf("stale version accepted: %v", err)
	}
	if err := s.Update("t", "a1", "Nuevo", "Contenido", 1); err != nil {
		t.Fatal(err)
	}
	// now version is 2
	if err := s.Publish("t", "a1", 2); err != nil {
		t.Fatal(err)
	}
}

func TestBadTransition(t *testing.T) {
	s := NewStore()
	_ = s.Create(art())
	_ = s.Publish("t", "a1", 1)
	if err := s.Publish("t", "a1", 2); !errors.Is(err, ErrBadTransition) {
		t.Fatalf("double publish accepted: %v", err)
	}
}

func TestSearchPublishedOnly(t *testing.T) {
	s := NewStore()
	_ = s.Create(art())
	_ = s.Create(Article{TenantID: "t", ID: "a2", Category: "faq", Title: "Devoluciones", Body: "Cómo devolver"})
	// only a1 is published
	_ = s.Publish("t", "a1", 1)
	results := s.Search("t", "reservar")
	if len(results) != 1 || results[0].ID != "a1" {
		t.Fatalf("search should return only published a1, got %+v", results)
	}
	if len(s.Search("t", "devolver")) != 0 {
		t.Fatal("draft a2 should not appear in search")
	}
}

func TestInvalidAndTenantIsolation(t *testing.T) {
	s := NewStore()
	bad := art()
	bad.Category = "BAD CAT"
	if err := s.Create(bad); !errors.Is(err, ErrInvalidArticle) {
		t.Fatalf("bad category accepted: %v", err)
	}
	_ = s.Create(art())
	if len(s.Published("other", "")) != 0 {
		t.Fatal("cross-tenant leak")
	}
}

func TestDistinctArticleTuples(t *testing.T) {
	s := NewStore()
	a, b := art(), art()
	a.TenantID = "a\x00b"
	a.ID = "c"
	b.TenantID = "a"
	b.ID = "b\x00c"
	if e := s.Create(a); e != nil {
		t.Fatal(e)
	}
	if e := s.Create(b); e != nil {
		t.Fatal("distinct article rejected", e)
	}
}
func TestForeignArticleUpdateRejected(t *testing.T) {
	s := NewStore()
	a := art()
	a.TenantID = "a\x00b"
	a.ID = "c"
	if e := s.Create(a); e != nil {
		t.Fatal(e)
	}
	if e := s.Update("a", "b\x00c", "foreign", "foreign", 1); !errors.Is(e, ErrNotFound) {
		t.Fatal("foreign edit admitted", e)
	}
	if e := s.Publish(a.TenantID, a.ID, 1); e != nil {
		t.Fatal("foreign edit consumed version", e)
	}
	if rows := s.Published(a.TenantID, ""); len(rows) != 1 || rows[0].Title != a.Title {
		t.Fatal("owner article changed", rows)
	}
}
func TestForeignArticleTransitionsRejected(t *testing.T) {
	for _, archive := range []bool{false, true} {
		t.Run(map[bool]string{false: "publish", true: "archive"}[archive], func(t *testing.T) {
			s := NewStore()
			a := art()
			a.TenantID = "a\x00b"
			a.ID = "c"
			if e := s.Create(a); e != nil {
				t.Fatal(e)
			}
			var e error
			if archive {
				if e = s.Publish(a.TenantID, a.ID, 1); e != nil {
					t.Fatal(e)
				}
				e = s.Archive("a", "b\x00c", 2)
			} else {
				e = s.Publish("a", "b\x00c", 1)
			}
			if !errors.Is(e, ErrNotFound) {
				t.Fatal("foreign transition admitted", e)
			}
			for _, row := range s.articles {
				want := StateDraft
				if archive {
					want = StatePublished
				}
				if row.State != want {
					t.Fatal("owner state changed")
				}
			}
		})
	}
}
func TestVersionExhaustionRejected(t *testing.T) {
	for _, op := range []string{"update", "publish", "archive"} {
		t.Run(op, func(t *testing.T) {
			s := NewStore()
			if e := s.Create(art()); e != nil {
				t.Fatal(e)
			}
			// Synthetic internal boundary: no public API initializes a MaxInt64 version.
			var before Article
			for k, a := range s.articles {
				a.Version = math.MaxInt64
				if op == "archive" {
					a.State = StatePublished
				}
				s.articles[k] = a
				before = a
			}
			var e error
			switch op {
			case "update":
				e = s.Update("t", "a1", "new", "body", math.MaxInt64)
			case "publish":
				e = s.Publish("t", "a1", math.MaxInt64)
			case "archive":
				e = s.Archive("t", "a1", math.MaxInt64)
			}
			if !errors.Is(e, ErrVersionExhausted) {
				t.Error("exhausted version accepted", e)
			}
			for _, a := range s.articles {
				if a != before {
					t.Error("failed operation changed article")
				}
			}
		})
	}
}

func TestVersionBoundaryAndErrorPrecedence(t *testing.T) {
	s := NewStore()
	if e := s.Create(art()); e != nil {
		t.Fatal(e)
	}
	for k, a := range s.articles {
		a.Version = math.MaxInt64 - 1
		s.articles[k] = a
	}
	if e := s.Update("t", "a1", "last", "body", math.MaxInt64-1); e != nil {
		t.Fatal(e)
	}
	if e := s.Update("t", "a1", "x", "body", 1); !errors.Is(e, ErrVersion) {
		t.Fatal("stale precedence", e)
	}
	if e := s.Update("t", "a1", " ", "body", math.MaxInt64); !errors.Is(e, ErrInvalidArticle) {
		t.Fatal("invalid precedence", e)
	}
	if e := s.Archive("t", "a1", math.MaxInt64); !errors.Is(e, ErrBadTransition) {
		t.Fatal("state precedence", e)
	}
	if e := s.Publish("missing", "a1", math.MaxInt64); !errors.Is(e, ErrNotFound) {
		t.Fatal("lookup precedence", e)
	}
	if e := s.Update("t", "a1", "x", "body", math.MaxInt64); !errors.Is(e, ErrVersionExhausted) {
		t.Fatal("overflow accepted", e)
	}
}
func TestConcurrentVersionHasSingleWinner(t *testing.T) {
	s := NewStore()
	if e := s.Create(art()); e != nil {
		t.Fatal(e)
	}
	var wg sync.WaitGroup
	var winners atomic.Int32
	for i := 0; i < 64; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			e := s.Update("t", "a1", fmt.Sprint(i), "body", 1)
			if e == nil {
				winners.Add(1)
			} else if !errors.Is(e, ErrVersion) {
				t.Error(e)
			}
			_ = s.String()
			s.Published("t", "")
			s.Search("t", "body")
		}(i)
	}
	wg.Wait()
	if winners.Load() != 1 {
		t.Fatal("optimistic concurrency lost", winners.Load())
	}
	if e := s.Publish("t", "a1", 2); e != nil {
		t.Fatal(e)
	}
}
func TestPublishedSearchCopyAndArchivedUpdate(t *testing.T) {
	s := NewStore()
	if e := s.Create(art()); e != nil {
		t.Fatal(e)
	}
	if e := s.Publish("t", "a1", 1); e != nil {
		t.Fatal(e)
	}
	rows := s.Search("t", " RESERVAR ")
	if len(rows) != 1 {
		t.Fatal("substring search", rows)
	}
	rows[0].Title = "mutated"
	if s.Published("t", "")[0].Title != art().Title {
		t.Fatal("returned article aliases stored value")
	}
	if len(s.Published("t", "other")) != 0 || len(s.Search("t", " ")) != 0 || len(s.Search("other", "reservar")) != 0 {
		t.Fatal("filter failure")
	}
	if e := s.Archive("t", "a1", 2); e != nil {
		t.Fatal(e)
	}
	if e := s.Update("t", "a1", "archived edit", "body", 3); e != nil {
		t.Fatal("existing archived edit semantics changed", e)
	}
	if len(s.Published("t", "")) != 0 {
		t.Fatal("archived edit republished article")
	}
}
func FuzzArticleHistoryModel(f *testing.F) {
	for _, prefix := range []string{"a", "x\x00y", "ñ", " "} {
		for _, ops := range [][]byte{{0, 0, 0, 1, 2, 0, 3, 0}, {0, 0, 0, 1, 1, 0, 2, 0, 3, 0}, {0, 0, 4, 0, 1, 0, 2, 0}, {0, 0, 2, 0, 4, 0, 3, 0}} {
			f.Add(prefix, ops)
		}
	}
	f.Fuzz(func(t *testing.T, prefix string, ops []byte) {
		if len(prefix) > 64 {
			return
		}
		if len(ops) > 128 {
			ops = ops[:128]
		}
		s := NewStore()
		model := map[string]map[string]Article{}
		tenants := []string{"t" + prefix + "\x00b", "t" + prefix, "other", "t" + prefix}
		ids := []string{"c", "b\x00c", "c", "c"}
		for i := 0; i+1 < len(ops); i += 2 {
			op := ops[i] % 5
			slot := int(ops[i+1] % 4)
			tenant, id := tenants[slot], ids[slot]
			if model[tenant] == nil {
				model[tenant] = map[string]Article{}
			}
			a, exists := model[tenant][id]
			expected := a.Version
			switch (ops[i+1] / 4) % 4 {
			case 1:
				expected--
			case 2:
				expected = math.MaxInt64
			case 3:
				expected = 0
			}
			title, body := "Title", "body"
			if ops[i]&128 != 0 {
				title = " "
			}
			var want, got error
			switch op {
			case 0:
				fresh := Article{TenantID: tenant, ID: id, Category: "faq", Title: title, Body: body}
				if strings.TrimSpace(title) == "" {
					want = ErrInvalidArticle
				} else if exists {
					want = ErrDuplicate
				} else {
					a = fresh
					a.Version = 1
					a.State = StateDraft
					model[tenant][id] = a
				}
				got = s.Create(fresh)
			case 1, 2, 3:
				if !exists {
					want = ErrNotFound
				} else if a.Version != expected {
					want = ErrVersion
				} else if op == 1 && strings.TrimSpace(title) == "" {
					want = ErrInvalidArticle
				} else if op == 2 && a.State != StateDraft || op == 3 && a.State != StatePublished {
					want = ErrBadTransition
				} else if a.Version == math.MaxInt64 {
					want = ErrVersionExhausted
				} else {
					a.Version++
					if op == 1 {
						a.Title = title
						a.Body = body
					} else if op == 2 {
						a.State = StatePublished
					} else {
						a.State = StateArchived
					}
					model[tenant][id] = a
				}
				if op == 1 {
					got = s.Update(tenant, id, title, body, expected)
				} else if op == 2 {
					got = s.Publish(tenant, id, expected)
				} else {
					got = s.Archive(tenant, id, expected)
				}
			case 4:
				// Synthetic private-state boundary only; no restored/persisted source claimed.
				if exists {
					a.Version = math.MaxInt64
					model[tenant][id] = a
					for k, v := range s.articles {
						if v.TenantID == tenant && v.ID == id {
							v.Version = math.MaxInt64
							s.articles[k] = v
						}
					}
				}
			}
			if !errors.Is(got, want) {
				t.Fatal("model error", op, slot, got, want)
			}
			total := 0
			for mt, articles := range model {
				total += len(articles)
				expectedPublished := map[string]Article{}
				for mi, m := range articles {
					found := false
					for _, actual := range s.articles {
						if actual.TenantID == mt && actual.ID == mi {
							found = true
							if actual.Title != m.Title || actual.Body != m.Body || actual.Category != m.Category || actual.Version != m.Version || actual.State != m.State {
								t.Fatal("stored model mismatch", actual, m)
							}
						}
					}
					if !found {
						t.Fatal("model article missing")
					}
					if m.State == StatePublished {
						expectedPublished[mi] = m
					}
				}
				check := func(rows []Article) {
					if len(rows) != len(expectedPublished) {
						t.Fatal("published count", len(rows), len(expectedPublished))
					}
					seen := map[string]bool{}
					for _, row := range rows {
						m, ok := expectedPublished[row.ID]
						if !ok || seen[row.ID] || row.TenantID != mt || row.Version != m.Version || row.Title != m.Title {
							t.Fatal("published scope/model mismatch")
						}
						seen[row.ID] = true
					}
				}
				check(s.Published(mt, "faq"))
				check(s.Search(mt, " BODY "))
				if len(s.Published(mt, "missing")) != 0 {
					t.Fatal("category ignored")
				}
			}
			if len(s.articles) != total {
				t.Fatal("article cardinality mismatch")
			}
			_ = s.String()
		}
	})
}
````


## 6. Configuration surface

Sin variables ni secretos.

## 7. Dependency bill

| Package/image/tool | Pin exacto | Uso | Licencia | Runtime/build | Fuente oficial |
|---|---|---|---|---|---|
| Go standard library | go 1.26.7 (toolchain) | artículos | BSD-3-Clause | runtime | https://go.dev |

## 8. Apply order

1. Componer el backend (mismo módulo).
2. Colocar los dos archivos bajo `internal/helpcenter/`.
3. Verificar con `go test ./... -count=1` y `go vet ./...`.
4. Rollback: eliminar `internal/helpcenter/`.

## 9. Verification

- `go test ./internal/helpcenter/ -count=1`: 5/5 PASS (estados, concurrencia optimista, transición inválida, búsqueda sólo publicados, inválido y aislamiento tenant).
- `go test ./... -count=1` (25 paquetes): PASS.
- `go vet ./...`: exit 0.

## 10. Reconstruction evidence

Véase `reconstruction_evidence/GO_HELP_CENTER_CORE_2026-09-02_V203.md`.

## Revisión V362

0.1.1 usa tupla tenant/id y sincroniza String. Conserva firmas, duplicados,
precedencia not-found→versión→contenido/estado y Update en cualquier estado
(incluido archived); editar no republica. Publicar/archivar sólo muta memoria.
ErrVersionExhausted rechaza incremento desde MaxInt64 sin cambiar ningún campo;
MaxInt64-1 aún puede incrementar. No reset/wrap ni política de recuperación nueva.
La prueba del extremo int64 inyecta el estado interno sintético: Create siempre
inicia en1; no se simulan años de uso, restore ni entrada de snapshot externo.
Search usa substring local; metadata anterior de delegación no era integración
ejecutable. V203 queda como historial; evidencia nueva:
reconstruction_evidence/DASHBOARD_HELP_BOUNDARIES_V362.md.

V362: entre ambos cores4fuentes reconstruidas,19tests x3,8tests históricos
intactos,vet/build,34semillas y4266789ejecuciones fuzz PASS.
Sin -race, auth externa, fuente KPI conectada, búsqueda externa o publicación real.

## Auditoría por claim V374

La revisión 0.1.2 se reconstruye y prueba en Go1.26.8/Windows con fuentes
oficiales fijadas y comparaciones por contrato. Ver
`reconstruction_evidence/CORE_CLAIM_ADMISSION_V374.md` desde la raíz: identidad,
licencia local, prueba por invariante, fuente semántica y dictamen de equivalencia
se mantienen separados. El lenguaje/stdlib no es fuente de un negocio empresarial.
No cambia este claim, API, permiso de composición ni las condiciones de integración.
La revisión 0.1.1 y sus bytes quedan preservados en el expediente anterior.
