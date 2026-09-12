# Go Warranty & Claims Core

## 1. Metadata

```yaml
pack_id: "GO-WARRANTY-CLAIMS-CORE"
pack_version: "0.1.2"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa el ciclo de garantía/reclamos tenant-scoped: open→in_review→approved/rejected→resolved/closed, fail-closed y con nota obligatoria al rechazar/resolver."
stacks: ["Go 1.26.8"]
compatible_with: ["GO-RETURN-EXCHANGE-FULFILLMENT-WORKER 0.1.x"]
incompatible_with: ["rechazo o resolución sin nota", "transición inválida"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://go.dev", "https://franpos.com/franchise-back-office-royalty-management-guide/"]
verified_at: "2026-09-10"
```

## 2. Applicability

Referencia AUTHORED en memoria para estados locales de reclamos. Reject/Resolve requieren nota antes de consultar identidad/estado. No determina elegibilidad, términos legales, devolución, reembolso ni reemplazo; el caller autoriza actor/tenant y conecta efectos admitidos por separado.

## 3. Architecture contract

- **Ownership**: `internal/warranty` gobierna el ciclo; la ejecución del efecto (reembolso/reemplazo) es del dominio de devoluciones.
- **Invariantes**: (1) estados gobernados. (2) rechazo y resolución exigen nota. (3) id de reclamo único por tenant.
- **Data flow**: `Open` → `Review` → (`Approve` → `Resolve` | `Reject`) → `Close`.
- **Failure modes**: nota ausente, transición inválida, duplicado → error.
- **Seguridad/privacidad**: tenant-scoped.
- **Performance budget**: O(1) por transición.

## 4. Exact file manifest

```text
CREATE internal/warranty/warranty.go
CREATE internal/warranty/warranty_test.go
```

## 5. Materialization blocks

### FILE: `internal/warranty/warranty.go`
```yaml
block_id: "GO-WARRANTY-CLAIMS-CORE:internal/warranty/warranty.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "d3f7799a357dc018014ea03bf48b8e9f2cc314f754fa325b9f74f23be7c34bbf"
variables: []
secrets_allowed: false
```
````go
// Package warranty provides a tenant-scoped warranty/claim lifecycle:
// open → in_review → approved/rejected → resolved/closed, fail-closed.
// Rejection and resolution require a note. AUTHORED.
package warranty

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"
)

// State is the claim lifecycle.
type State string

const (
	StateOpen     State = "open"
	StateInReview State = "in_review"
	StateApproved State = "approved"
	StateRejected State = "rejected"
	StateResolved State = "resolved"
	StateClosed   State = "closed"
)

var (
	idRe = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:-]{0,127}$`)

	ErrInvalidClaim  = errors.New("warranty: invalid claim")
	ErrDuplicate     = errors.New("warranty: duplicate claim")
	ErrNotFound      = errors.New("warranty: not found")
	ErrBadTransition = errors.New("warranty: bad state transition")
	ErrNoteRequired  = errors.New("warranty: note required")
)

// Claim is a single warranty/reclamation request.
type Claim struct {
	TenantID  string
	ID        string
	SubjectID string
	State     State
	Note      string
	At        time.Time
}

// Store holds claims.
type Store struct {
	mu     sync.Mutex
	claims map[identity]Claim
}

// NewStore returns an empty store.
func NewStore() *Store {
	return &Store{claims: make(map[identity]Claim)}
}

type identity struct{ tenant, id string }

func key(tenant, id string) identity { return identity{tenant: tenant, id: id} }

func validate(c Claim) error {
	if strings.TrimSpace(c.TenantID) == "" || !idRe.MatchString(c.ID) || !idRe.MatchString(c.SubjectID) {
		return ErrInvalidClaim
	}
	return nil
}

// Open registers a new claim in state open.
func (s *Store) Open(c Claim) error {
	if err := validate(c); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.claims[key(c.TenantID, c.ID)]; ok {
		return ErrDuplicate
	}
	c.State = StateOpen
	if c.At.IsZero() {
		c.At = time.Now().UTC()
	}
	s.claims[key(c.TenantID, c.ID)] = c
	return nil
}

// Review marks open → in_review.
func (s *Store) Review(tenant, id string) error {
	return s.transition(tenant, id, StateOpen, StateInReview, "")
}

// Approve marks in_review → approved.
func (s *Store) Approve(tenant, id string) error {
	return s.transition(tenant, id, StateInReview, StateApproved, "")
}

// Reject marks in_review → rejected; a reason note is required.
func (s *Store) Reject(tenant, id, note string) error {
	if strings.TrimSpace(note) == "" {
		return ErrNoteRequired
	}
	return s.transition(tenant, id, StateInReview, StateRejected, note)
}

// Resolve marks approved → resolved; a resolution note is required.
func (s *Store) Resolve(tenant, id, note string) error {
	if strings.TrimSpace(note) == "" {
		return ErrNoteRequired
	}
	return s.transition(tenant, id, StateApproved, StateResolved, note)
}

// Close marks rejected/resolved → closed.
func (s *Store) Close(tenant, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	c, ok := s.claims[key(tenant, id)]
	if !ok {
		return ErrNotFound
	}
	if c.State != StateRejected && c.State != StateResolved {
		return ErrBadTransition
	}
	c.State = StateClosed
	s.claims[key(tenant, id)] = c
	return nil
}

func (s *Store) transition(tenant, id string, from, to State, note string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	c, ok := s.claims[key(tenant, id)]
	if !ok {
		return ErrNotFound
	}
	if c.State != from {
		return ErrBadTransition
	}
	c.State = to
	if note != "" {
		c.Note = note
	}
	s.claims[key(tenant, id)] = c
	return nil
}

// State returns the current state.
func (s *Store) State(tenant, id string) (State, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	c, ok := s.claims[key(tenant, id)]
	return c.State, ok
}

// String aids debugging.
func (s *Store) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return fmt.Sprintf("warranty(%d)", len(s.claims))
}
````

### FILE: `internal/warranty/warranty_test.go`
```yaml
block_id: "GO-WARRANTY-CLAIMS-CORE:internal/warranty/warranty_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "a77d369a631557626ecfabf22f83bfa3ded9ebbce094faec5ce90dcdc2fabe20"
variables: []
secrets_allowed: false
```
````go
package warranty

import (
	"errors"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
)

func TestHappyPath(t *testing.T) {
	s := NewStore()
	c := Claim{TenantID: "t", ID: "c1", SubjectID: "sub1"}
	if err := s.Open(c); err != nil {
		t.Fatal(err)
	}
	if err := s.Review("t", "c1"); err != nil {
		t.Fatal(err)
	}
	if err := s.Approve("t", "c1"); err != nil {
		t.Fatal(err)
	}
	if err := s.Resolve("t", "c1", "refund issued"); err != nil {
		t.Fatal(err)
	}
	if err := s.Close("t", "c1"); err != nil {
		t.Fatal(err)
	}
	if st, _ := s.State("t", "c1"); st != StateClosed {
		t.Fatalf("expected closed, got %s", st)
	}
}

func TestRejectRequiresNote(t *testing.T) {
	s := NewStore()
	if err := s.Open(Claim{TenantID: "t", ID: "c1", SubjectID: "sub1"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Review("t", "c1"); err != nil {
		t.Fatal(err)
	}
	if err := s.Reject("t", "c1", ""); !errors.Is(err, ErrNoteRequired) {
		t.Fatalf("expected note required, got %v", err)
	}
	if err := s.Reject("t", "c1", "out of warranty"); err != nil {
		t.Fatal(err)
	}
	if st, _ := s.State("t", "c1"); st != StateRejected {
		t.Fatalf("expected rejected, got %s", st)
	}
}

func TestResolveRequiresNoteAndBadTransition(t *testing.T) {
	s := NewStore()
	if err := s.Open(Claim{TenantID: "t", ID: "c1", SubjectID: "sub1"}); err != nil {
		t.Fatal(err)
	}
	// Approve before review → bad transition.
	if err := s.Approve("t", "c1"); !errors.Is(err, ErrBadTransition) {
		t.Fatalf("expected bad transition, got %v", err)
	}
	if err := s.Review("t", "c1"); err != nil {
		t.Fatal(err)
	}
	if err := s.Approve("t", "c1"); err != nil {
		t.Fatal(err)
	}
	if err := s.Resolve("t", "c1", ""); !errors.Is(err, ErrNoteRequired) {
		t.Fatalf("expected note required, got %v", err)
	}
}

func TestDuplicateAndTenantIsolation(t *testing.T) {
	s := NewStore()
	c := Claim{TenantID: "t", ID: "c1", SubjectID: "sub1"}
	if err := s.Open(c); err != nil {
		t.Fatal(err)
	}
	if err := s.Open(c); !errors.Is(err, ErrDuplicate) {
		t.Fatalf("expected duplicate, got %v", err)
	}
	if _, ok := s.State("t2", "c1"); ok {
		t.Fatal("claim leaked across tenants")
	}
}

func TestAliasedClaimCannotReadForeignState(t *testing.T) {
	s := NewStore()
	if e := s.Open(Claim{TenantID: "a\x00b", ID: "c", SubjectID: "item"}); e != nil {
		t.Fatal(e)
	}
	if st, ok := s.State("a", "b\x00c"); ok || st != "" {
		t.Fatalf("foreign state: %q %v", st, ok)
	}
}
func TestAliasedClaimCannotTransition(t *testing.T) {
	for _, action := range []string{"review", "approve", "reject", "resolve", "close-resolved", "close-rejected"} {
		t.Run(action, func(t *testing.T) {
			s := NewStore()
			if e := s.Open(Claim{TenantID: "a\x00b", ID: "c", SubjectID: "item"}); e != nil {
				t.Fatal(e)
			}
			want := StateOpen
			if action != "review" {
				if e := s.Review("a\x00b", "c"); e != nil {
					t.Fatal(e)
				}
				want = StateInReview
			}
			if action == "resolve" || action == "close-resolved" {
				if e := s.Approve("a\x00b", "c"); e != nil {
					t.Fatal(e)
				}
				want = StateApproved
			}
			if action == "close-resolved" {
				if e := s.Resolve("a\x00b", "c", "resolved"); e != nil {
					t.Fatal(e)
				}
				want = StateResolved
			}
			if action == "close-rejected" {
				if e := s.Reject("a\x00b", "c", "rejected"); e != nil {
					t.Fatal(e)
				}
				want = StateRejected
			}
			var e error
			switch action {
			case "review":
				e = s.Review("a", "b\x00c")
			case "approve":
				e = s.Approve("a", "b\x00c")
			case "reject":
				e = s.Reject("a", "b\x00c", "reason")
			case "resolve":
				e = s.Resolve("a", "b\x00c", "resolution")
			default:
				e = s.Close("a", "b\x00c")
			}
			if !errors.Is(e, ErrNotFound) {
				t.Errorf("foreign transition got %v", e)
			}
			if st, ok := s.State("a\x00b", "c"); !ok || st != want {
				t.Errorf("original state changed: %q want %q", st, want)
			}
		})
	}
}

func TestRequiredNotePrecedesLookupAndDoesNotTransition(t *testing.T) {
	s := NewStore()
	for _, note := range []string{"", " \t\n"} {
		if e := s.Reject("missing", "claim", note); !errors.Is(e, ErrNoteRequired) {
			t.Fatal("reject error precedence changed", e)
		}
		if e := s.Resolve("missing", "claim", note); !errors.Is(e, ErrNoteRequired) {
			t.Fatal("resolve error precedence changed", e)
		}
	}
	if e := s.Open(Claim{TenantID: "t", ID: "c", SubjectID: "item", State: StateClosed}); e != nil {
		t.Fatal(e)
	}
	if st, _ := s.State("t", "c"); st != StateOpen {
		t.Fatal("caller bypassed open")
	}
	if e := s.Review("t", "c"); e != nil {
		t.Fatal(e)
	}
	if e := s.Reject("t", "c", " "); !errors.Is(e, ErrNoteRequired) {
		t.Fatal(e)
	}
	if st, _ := s.State("t", "c"); st != StateInReview {
		t.Fatal("blank reject mutated state")
	}
	if e := s.Approve("t", "c"); e != nil {
		t.Fatal(e)
	}
	if e := s.Resolve("t", "c", " "); !errors.Is(e, ErrNoteRequired) {
		t.Fatal(e)
	}
	if st, _ := s.State("t", "c"); st != StateApproved {
		t.Fatal("blank resolution mutated state")
	}
	if e := s.Resolve("t", "c", "resolved"); e != nil {
		t.Fatal(e)
	}
	if e := s.Close("t", "c"); e != nil {
		t.Fatal(e)
	}
	if e := s.Open(Claim{TenantID: "t", ID: "c", SubjectID: "item"}); !errors.Is(e, ErrDuplicate) {
		t.Fatal("closed identity unexpectedly reused", e)
	}
}
func TestConcurrentClaimLifecycle(t *testing.T) {
	s := NewStore()
	var success atomic.Int64
	phase := func(action func() error, repeated error) {
		var wg sync.WaitGroup
		success.Store(0)
		for i := 0; i < 64; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				e := action()
				if e == nil {
					success.Add(1)
				} else if !errors.Is(e, repeated) {
					t.Error(e)
				}
				s.State("t", "c")
				_ = s.String()
			}()
		}
		wg.Wait()
		if success.Load() != 1 {
			t.Fatal("expected one successful mutation", success.Load())
		}
	}
	phase(func() error { return s.Open(Claim{TenantID: "t", ID: "c", SubjectID: "item"}) }, ErrDuplicate)
	phase(func() error { return s.Review("t", "c") }, ErrBadTransition)
	phase(func() error { return s.Approve("t", "c") }, ErrBadTransition)
	phase(func() error { return s.Resolve("t", "c", "resolved") }, ErrBadTransition)
	phase(func() error { return s.Close("t", "c") }, ErrBadTransition)
	if st, ok := s.State("t", "c"); !ok || st != StateClosed {
		t.Fatal("wrong final state")
	}
}
func FuzzClaimLifecycleModel(f *testing.F) {
	for _, x := range []string{"a", "a\x00b", "", " ", "客户", "b\x00c"} {
		f.Add(x, []byte{0, 1, 0, 0, 1, 0, 1, 0, 2, 0, 1, 0, 3, 0, 1, 0, 4, 0, 1, 0, 5, 0, 1, 0})
		f.Add(x, []byte{0, 0, 0, 0, 1, 0, 0, 0, 2, 0, 0, 0, 4, 0, 0, 1, 4, 0, 0, 0, 5, 0, 0, 0})
		f.Add(x, []byte{0, 0, 0, 0, 1, 0, 0, 0, 3, 0, 0, 0, 5, 0, 0, 0})
	}
	f.Fuzz(func(t *testing.T, x string, ops []byte) {
		if len(x) > 64 {
			x = x[:64]
		}
		if len(ops) > 160 {
			ops = ops[:160]
		}
		tenants := []string{"a", "a\x00b", x, " "}
		ids := []string{"c", "b\x00c", x, " "}
		subjects := []string{"item", x, " "}
		notes := []string{"reason", "", " \t", x}
		valid := func(id string) bool {
			if len(id) < 1 || len(id) > 128 {
				return false
			}
			for i, c := range []byte(id) {
				alnum := (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
				if !alnum && (i == 0 || (c != '.' && c != '_' && c != ':' && c != '-')) {
					return false
				}
			}
			return true
		}
		var model []Claim
		s := NewStore()
		lookup := func(tenant, id string) int {
			for i, c := range model {
				if c.TenantID == tenant && c.ID == id {
					return i
				}
			}
			return -1
		}
		for i := 0; i+3 < len(ops); i += 4 {
			tenant := tenants[int(ops[i+1])%len(tenants)]
			id := ids[int(ops[i+2])%len(ids)]
			index := lookup(tenant, id)
			action := ops[i] % 7
			note := notes[int(ops[i+3])%len(notes)]
			var got, want error
			if action == 0 {
				c := Claim{TenantID: tenant, ID: id, SubjectID: subjects[int(ops[i+3])%len(subjects)], State: StateClosed}
				if strings.TrimSpace(tenant) == "" || !valid(id) || !valid(c.SubjectID) {
					want = ErrInvalidClaim
				} else if index >= 0 {
					want = ErrDuplicate
				}
				got = s.Open(c)
				if want == nil {
					c.State = StateOpen
					model = append(model, c)
				}
			} else if action < 6 {
				// Notes are validated before identity or state, preserving the public API.
				if (action == 3 || action == 4) && strings.TrimSpace(note) == "" {
					want = ErrNoteRequired
				} else if index < 0 {
					want = ErrNotFound
				} else {
					st := model[index].State
					next := State("")
					switch action {
					case 1:
						if st == StateOpen {
							next = StateInReview
						}
					case 2:
						if st == StateInReview {
							next = StateApproved
						}
					case 3:
						if st == StateInReview {
							next = StateRejected
						}
					case 4:
						if st == StateApproved {
							next = StateResolved
						}
					case 5:
						if st == StateRejected || st == StateResolved {
							next = StateClosed
						}
					}
					if next == "" {
						want = ErrBadTransition
					} else {
						model[index].State = next
					}
				}
				switch action {
				case 1:
					got = s.Review(tenant, id)
				case 2:
					got = s.Approve(tenant, id)
				case 3:
					got = s.Reject(tenant, id, note)
				case 4:
					got = s.Resolve(tenant, id, note)
				case 5:
					got = s.Close(tenant, id)
				}
			} else {
				_ = s.String()
			}
			if !errors.Is(got, want) {
				t.Fatalf("step %d got %v want %v", i/4, got, want)
			}
			for _, tenant := range tenants {
				for _, id := range ids {
					idx := lookup(tenant, id)
					st, ok := s.State(tenant, id)
					if ok != (idx >= 0) {
						t.Fatal("foreign/missing identity")
					}
					if idx < 0 && st != "" {
						t.Fatal("absent nonempty state")
					}
					if idx >= 0 && st != model[idx].State {
						t.Fatal("lifecycle mismatch")
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
| Go standard library | go 1.26.7 (toolchain) | garantías/reclamos | BSD-3-Clause | runtime | https://go.dev |

## 8. Apply order

1. Componer el backend (mismo módulo).
2. Colocar los dos archivos bajo `internal/warranty/`.
3. Verificar con `go test ./... -count=1` y `go vet ./...`.
4. Rollback: eliminar `internal/warranty/`.

## 9. Verification

- `go test ./internal/warranty/ -count=1`: 4/4 PASS.
- `go test ./... -count=1` (37 paquetes): PASS.
- `go vet ./...`: exit 0.

## 10. Reconstruction evidence

Véase `reconstruction_evidence/GO_WARRANTY_CLAIMS_CORE_2026-09-02_V214.md`.

## Revisión V359 — identidad sin alias

0.1.1 usa tuplas tenant/identidad y preserva firmas, validación, prioridad de
errores y estados. AUTHORED/LicenseRef-Workspace-Owner, CONDITIONED y fuera
del perfil persisten. Compatibilidad histórica no acredita montaje actual.
La suite V214 es histórica, no composición reejecutada en V359.
Evidencia: reconstruction_evidence/ONBOARDING_CLAIM_ISOLATION_V359.md.

V359 verificado entre ambos cores: 4/4 fuentes reconstruidas, 9 tests históricos
intactos, 4 grupos red, 17 tests x3, vet/build y 30 semillas/4760595 ejecuciones
fuzz PASS. Sin -race, activación, elegibilidad legal ni efectos de proveedor.

## Auditoría por claim V374

La revisión 0.1.2 se reconstruye y prueba en Go1.26.8/Windows con fuentes
oficiales fijadas y comparaciones por contrato. Ver
`reconstruction_evidence/CORE_CLAIM_ADMISSION_V374.md` desde la raíz: identidad,
licencia local, prueba por invariante, fuente semántica y dictamen de equivalencia
se mantienen separados. El lenguaje/stdlib no es fuente de un negocio empresarial.
No cambia este claim, API, permiso de composición ni las condiciones de integración.
La revisión 0.1.1 y sus bytes quedan preservados en el expediente anterior.
