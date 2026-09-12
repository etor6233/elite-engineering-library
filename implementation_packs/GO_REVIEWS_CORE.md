# Go Reviews Core

## 1. Metadata

```yaml
pack_id: "GO-REVIEWS-CORE"
pack_version: "0.1.2"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa reseñas/ratings tenant-scoped con moderación: rating 1-5, una por autor por sujeto, estados pending→approved/rejected y promedio sólo sobre aprobadas."
stacks: ["Go 1.26.8"]
compatible_with: ["GO-ENTERPRISE-BACKEND 0.4.x"]
incompatible_with: ["rating fuera de 1-5", "reseña duplicada", "promedio con no-aprobadas"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://go.dev"]
verified_at: "2026-09-10"
```

## 2. Applicability

Referencia AUTHORED en memoria para estado de moderación y promedio de aprobadas. El caller autoriza tenant, subject, author y moderador; no prueba compra verificada, protección antiabuso, publicación web ni retención durable.

## 3. Architecture contract

- **Ownership**: `internal/reviews` gobierna reseñas y moderación.
- **Invariantes**: (1) rating 1-5. (2) una reseña por autor por sujeto. (3) promedio sólo aprobadas. (4) tenant-scoped.
- **Data flow**: `Submit` → `Approve`/`Reject` → `Approved`/`AverageRating`.
- **Failure modes**: rating inválido, duplicado, transición inválida → error.
- **Seguridad/privacidad**: tenant-scoped.
- **Performance budget**: O(1) submit, O(N) average.

## 4. Exact file manifest

```text
CREATE internal/reviews/reviews.go
CREATE internal/reviews/reviews_test.go
```

## 5. Materialization blocks

### FILE: `internal/reviews/reviews.go`
```yaml
block_id: "GO-REVIEWS-CORE:internal/reviews/reviews.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "6dc3aa976fb612ee7157a009c7de4638ece3558b9e83f0eeccb4bf4df2bd908c"
variables: []
secrets_allowed: false
```
````go
// Package reviews provides tenant-scoped customer reviews with moderation: a
// rating 1-5, one review per author per subject, and a state machine
// pending → approved/rejected; averages are computed over approved only.
package reviews

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
)

// State is the moderation lifecycle.
type State string

const (
	StatePending  State = "pending"
	StateApproved State = "approved"
	StateRejected State = "rejected"
)

var (
	ErrInvalidRating = errors.New("reviews: rating must be 1..5")
	ErrDuplicate     = errors.New("reviews: author already reviewed this subject")
	ErrNotFound      = errors.New("reviews: not found")
	ErrBadTransition = errors.New("reviews: bad state transition")
)

// Review is a single moderated review.
type Review struct {
	TenantID  string
	SubjectID string
	AuthorID  string
	Rating    int
	Text      string
	State     State
	At        time.Time
}

// Store holds reviews.
type Store struct {
	mu      sync.Mutex
	reviews map[identity]Review
}

// NewStore returns an empty store.
func NewStore() *Store {
	return &Store{reviews: make(map[identity]Review)}
}

type identity struct{ tenant, subject, author string }

func key(tenant, subject, author string) identity {
	return identity{tenant: tenant, subject: subject, author: author}
}

// Submit adds a review pending moderation; one per author per subject.
func (s *Store) Submit(r Review) error {
	if r.Rating < 1 || r.Rating > 5 {
		return ErrInvalidRating
	}
	if strings.TrimSpace(r.TenantID) == "" || strings.TrimSpace(r.SubjectID) == "" || strings.TrimSpace(r.AuthorID) == "" {
		return ErrInvalidRating
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.reviews[key(r.TenantID, r.SubjectID, r.AuthorID)]; ok {
		return ErrDuplicate
	}
	r.State = StatePending
	if r.At.IsZero() {
		r.At = time.Now().UTC()
	}
	s.reviews[key(r.TenantID, r.SubjectID, r.AuthorID)] = r
	return nil
}

// Approve marks pending → approved.
func (s *Store) Approve(tenant, subject, author string) error {
	return s.transition(tenant, subject, author, StatePending, StateApproved)
}

// Reject marks pending → rejected.
func (s *Store) Reject(tenant, subject, author string) error {
	return s.transition(tenant, subject, author, StatePending, StateRejected)
}

func (s *Store) transition(tenant, subject, author string, from, to State) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	k := key(tenant, subject, author)
	r, ok := s.reviews[k]
	if !ok {
		return ErrNotFound
	}
	if r.State != from {
		return ErrBadTransition
	}
	r.State = to
	s.reviews[k] = r
	return nil
}

// Approved returns the tenant-scoped approved reviews for a subject.
func (s *Store) Approved(tenant, subject string) []Review {
	s.mu.Lock()
	defer s.mu.Unlock()
	var out []Review
	for _, r := range s.reviews {
		if r.TenantID == tenant && r.SubjectID == subject && r.State == StateApproved {
			out = append(out, r)
		}
	}
	return out
}

// AverageRating returns the average over approved reviews, 0 if none.
func (s *Store) AverageRating(tenant, subject string) float64 {
	approved := s.Approved(tenant, subject)
	if len(approved) == 0 {
		return 0
	}
	var sum int
	for _, r := range approved {
		sum += r.Rating
	}
	return float64(sum) / float64(len(approved))
}

// String helps debug formatting.
func (s *Store) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return fmt.Sprintf("reviews(%d)", len(s.reviews))
}
````

### FILE: `internal/reviews/reviews_test.go`
```yaml
block_id: "GO-REVIEWS-CORE:internal/reviews/reviews_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "71b05f0a282a482d103219e3bbbf101cf78b033eef4cbb5453bed9c38d73e395"
variables: []
secrets_allowed: false
```
````go
package reviews

import (
	"errors"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
)

func rev(rating int) Review {
	return Review{TenantID: "t", SubjectID: "P1", AuthorID: "a1", Rating: rating, Text: "ok"}
}

func TestSubmitModerationAndAverage(t *testing.T) {
	s := NewStore()
	_ = s.Submit(rev(5))
	_ = s.Submit(Review{TenantID: "t", SubjectID: "P1", AuthorID: "a2", Rating: 3, Text: "ok"})
	// not yet approved → average 0
	if s.AverageRating("t", "P1") != 0 {
		t.Fatalf("pending reviews should not count, got %v", s.AverageRating("t", "P1"))
	}
	_ = s.Approve("t", "P1", "a1")
	_ = s.Reject("t", "P1", "a2")
	if s.AverageRating("t", "P1") != 5 {
		t.Fatalf("expected 5 (only approved), got %v", s.AverageRating("t", "P1"))
	}
	if len(s.Approved("t", "P1")) != 1 {
		t.Fatalf("expected 1 approved, got %d", len(s.Approved("t", "P1")))
	}
}

func TestInvalidRating(t *testing.T) {
	s := NewStore()
	if err := s.Submit(rev(0)); !errors.Is(err, ErrInvalidRating) {
		t.Fatalf("rating 0 accepted: %v", err)
	}
	if err := s.Submit(rev(6)); !errors.Is(err, ErrInvalidRating) {
		t.Fatalf("rating 6 accepted: %v", err)
	}
}

func TestDuplicateAuthorRejected(t *testing.T) {
	s := NewStore()
	_ = s.Submit(rev(4))
	if err := s.Submit(rev(5)); !errors.Is(err, ErrDuplicate) {
		t.Fatalf("duplicate author accepted: %v", err)
	}
}

func TestBadTransition(t *testing.T) {
	s := NewStore()
	_ = s.Submit(rev(5))
	_ = s.Approve("t", "P1", "a1")
	if err := s.Approve("t", "P1", "a1"); !errors.Is(err, ErrBadTransition) {
		t.Fatalf("double approve accepted: %v", err)
	}
}

func TestTenantIsolation(t *testing.T) {
	s := NewStore()
	_ = s.Submit(rev(5))
	if len(s.Approved("other", "P1")) != 0 {
		t.Fatal("cross-tenant leak")
	}
}

func TestDistinctModerationTuplesDoNotCollide(t *testing.T) {
	for _, p := range [][2]Review{
		{{TenantID: "a\x00b", SubjectID: "c", AuthorID: "d", Rating: 5}, {TenantID: "a", SubjectID: "b\x00c", AuthorID: "d", Rating: 3}},
		{{TenantID: "a", SubjectID: "b\x00c", AuthorID: "d", Rating: 5}, {TenantID: "a", SubjectID: "b", AuthorID: "c\x00d", Rating: 3}},
	} {
		s := NewStore()
		if e := s.Submit(p[0]); e != nil {
			t.Fatal(e)
		}
		if e := s.Submit(p[1]); e != nil {
			t.Errorf("distinct review rejected: %v", e)
		}
	}
}
func TestAliasedModerationCannotPublishOrReject(t *testing.T) {
	for _, approve := range []bool{true, false} {
		s := NewStore()
		if e := s.Submit(Review{TenantID: "a\x00b", SubjectID: "c", AuthorID: "d", Rating: 5}); e != nil {
			t.Fatal(e)
		}
		var e error
		if approve {
			e = s.Approve("a", "b\x00c", "d")
		} else {
			e = s.Reject("a", "b\x00c", "d")
		}
		if !errors.Is(e, ErrNotFound) {
			t.Errorf("foreign moderation got %v", e)
		}
		if len(s.Approved("a\x00b", "c")) != 0 || s.AverageRating("a\x00b", "c") != 0 {
			t.Error("foreign moderation published review")
		}
		if e := s.Approve("a\x00b", "c", "d"); e != nil {
			t.Errorf("original pending state changed: %v", e)
		}
	}
}

func TestApprovedSnapshotCannotMutateStore(t *testing.T) {
	s := NewStore()
	r := rev(5)
	r.State = StateApproved
	if e := s.Submit(r); e != nil {
		t.Fatal(e)
	}
	if s.AverageRating("t", "P1") != 0 {
		t.Fatal("caller bypassed moderation")
	}
	if e := s.Approve("t", "P1", "a1"); e != nil {
		t.Fatal(e)
	}
	out := s.Approved("t", "P1")
	out[0].Rating = 1
	out[0].Text = "changed"
	out[0].TenantID = "other"
	again := s.Approved("t", "P1")
	if len(again) != 1 || again[0].Rating != 5 || again[0].Text != "ok" || s.AverageRating("t", "P1") != 5 {
		t.Fatal("returned copy aliases store")
	}
}
func TestConcurrentModerationAndDebugRead(t *testing.T) {
	s := NewStore()
	var wg sync.WaitGroup
	var submits, transitions atomic.Int64
	for i := 0; i < 64; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			e := s.Submit(rev(5))
			if e == nil {
				submits.Add(1)
			} else if !errors.Is(e, ErrDuplicate) {
				t.Error(e)
			}
			_ = s.String()
		}()
	}
	wg.Wait()
	if submits.Load() != 1 {
		t.Fatal("duplicate submission")
	}
	for i := 0; i < 64; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			var e error
			if i%2 == 0 {
				e = s.Approve("t", "P1", "a1")
			} else {
				e = s.Reject("t", "P1", "a1")
			}
			if e == nil {
				transitions.Add(1)
			} else if !errors.Is(e, ErrBadTransition) {
				t.Error(e)
			}
			_ = s.String()
			s.Approved("t", "P1")
		}(i)
	}
	wg.Wait()
	if transitions.Load() != 1 {
		t.Fatal("multiple terminal transitions")
	}
}
func FuzzReviewModerationModel(f *testing.F) {
	for _, x := range []string{"a", "a\x00b", "", " ", "客户", "b\x00c"} {
		f.Add(x, []byte{0, 1, 0, 0, 4, 0, 0, 1, 0, 2, 1, 0, 1, 0, 0, 2, 1, 0, 0, 0})
		f.Add(x, []byte{0, 0, 0, 0, 5, 0, 0, 0, 0, 4, 1, 0, 0, 0, 0, 2, 0, 0, 0, 0})
	}
	f.Fuzz(func(t *testing.T, x string, ops []byte) {
		if len(x) > 64 {
			x = x[:64]
		}
		if len(ops) > 120 {
			ops = ops[:120]
		}
		tenants := []string{"a", x, "a\x00b", " "}
		subjects := []string{"c", "b\x00c", x, " "}
		authors := []string{"d", "c\x00d", x, " "}
		s := NewStore()
		var model []Review
		lookup := func(tenant, subject, author string) int {
			for i, r := range model {
				if r.TenantID == tenant && r.SubjectID == subject && r.AuthorID == author {
					return i
				}
			}
			return -1
		}
		for i := 0; i+4 < len(ops); i += 5 {
			tenant := tenants[int(ops[i+1])%len(tenants)]
			subject := subjects[int(ops[i+2])%len(subjects)]
			author := authors[int(ops[i+3])%len(authors)]
			index := lookup(tenant, subject, author)
			var got, want error
			switch ops[i] % 4 {
			case 0:
				r := Review{TenantID: tenant, SubjectID: subject, AuthorID: author, Rating: int(ops[i+4])%8 - 1, State: StateApproved, Text: "synthetic"}
				if r.Rating < 1 || r.Rating > 5 || strings.TrimSpace(tenant) == "" || strings.TrimSpace(subject) == "" || strings.TrimSpace(author) == "" {
					want = ErrInvalidRating
				} else if index >= 0 {
					want = ErrDuplicate
				}
				got = s.Submit(r)
				if want == nil {
					r.State = StatePending
					model = append(model, r)
				}
			case 1, 2:
				if index < 0 {
					want = ErrNotFound
				} else if model[index].State != StatePending {
					want = ErrBadTransition
				} else {
					if ops[i]%4 == 1 {
						model[index].State = StateApproved
					} else {
						model[index].State = StateRejected
					}
				}
				if ops[i]%4 == 1 {
					got = s.Approve(tenant, subject, author)
				} else {
					got = s.Reject(tenant, subject, author)
				}
			case 3:
				_ = s.String()
			}
			if !errors.Is(got, want) {
				t.Fatalf("step %d got=%v want=%v", i/5, got, want)
			}
			for _, tenant := range tenants {
				for _, subject := range subjects {
					var expected []Review
					sum := 0
					for _, r := range model {
						if r.TenantID == tenant && r.SubjectID == subject && r.State == StateApproved {
							expected = append(expected, r)
							sum += r.Rating
						}
					}
					out := s.Approved(tenant, subject)
					if len(out) != len(expected) {
						t.Fatal("approved count mismatch")
					}
					seen := make([]bool, len(expected))
					for _, r := range out {
						found := -1
						for j, e := range expected {
							if r.TenantID == e.TenantID && r.SubjectID == e.SubjectID && r.AuthorID == e.AuthorID && r.Rating == e.Rating && r.Text == e.Text && r.State == e.State {
								found = j
								break
							}
						}
						if found < 0 || seen[found] {
							t.Fatal("foreign/duplicate publication")
						}
						seen[found] = true
					}
					avg := float64(0)
					if len(expected) > 0 {
						avg = float64(sum) / float64(len(expected))
					}
					if s.AverageRating(tenant, subject) != avg {
						t.Fatal("average not approved model")
					}
				}
			}
		}
	})
}
````


## 6. Configuration surface

Sin variables ni secretos.

## 7. Dependency bill

| Package/image/tool | Pin exacto | Uso | Licencia | Runtime/build | Fuente oficial |
|---|---|---|---|---|---|
| Go standard library | go 1.26.7 (toolchain) | reseñas | BSD-3-Clause | runtime | https://go.dev |

## 8. Apply order

1. Componer el backend (mismo módulo).
2. Colocar los dos archivos bajo `internal/reviews/`.
3. Verificar con `go test ./... -count=1` y `go vet ./...`.
4. Rollback: eliminar `internal/reviews/`.

## 9. Verification

- `go test ./internal/reviews/ -count=1`: 5/5 PASS (moderación y promedio, rating inválido, duplicado, transición inválida, aislamiento tenant).
- `go test ./... -count=1` (19 paquetes): PASS.
- `go vet ./...`: exit 0.

## 10. Reconstruction evidence

Véase `reconstruction_evidence/GO_REVIEWS_CORE_2026-09-02_V197.md`.

## Revisión V358 — identidad sin alias

0.1.1 usa tuplas de identidad, preserva firmas, validación y transiciones.
AUTHORED/LicenseRef-Workspace-Owner, CONDITIONED y fuera del perfil persisten.
Suite V197 es histórica, no composición reejecutada. No integración ni
permiso externo inferido. Evidencia: reconstruction_evidence/MODERATION_WAITLIST_ISOLATION_V358.md.

V358 verificado entre ambos cores:4/4fuentes reconstruidas,11tests históricos
intactos,4regresiones red,19tests x3,vet/build y24semillas/4216603ejecuciones
fuzz PASS. Concurrencia local probada; sin -race, auth ni efectos de proveedor.

## Auditoría por claim V374

La revisión 0.1.2 se reconstruye y prueba en Go1.26.8/Windows con fuentes
oficiales fijadas y comparaciones por contrato. Ver
`reconstruction_evidence/CORE_CLAIM_ADMISSION_V374.md` desde la raíz: identidad,
licencia local, prueba por invariante, fuente semántica y dictamen de equivalencia
se mantienen separados. El lenguaje/stdlib no es fuente de un negocio empresarial.
No cambia este claim, API, permiso de composición ni las condiciones de integración.
La revisión 0.1.1 y sus bytes quedan preservados en el expediente anterior.
