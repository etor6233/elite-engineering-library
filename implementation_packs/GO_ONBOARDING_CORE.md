# Go Onboarding Core

## 1. Metadata

```yaml
pack_id: "GO-ONBOARDING-CORE"
pack_version: "0.1.2"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa el onboarding del franquiciado: pasos ordenados, checklist con progreso y completado idempotente, tenant-scoped."
stacks: ["Go 1.26.8"]
compatible_with: ["GO-ENTERPRISE-BACKEND 0.4.x", "GO-FRANCHISE-ROYALTY-SETTLEMENT-API 0.1.x"]
incompatible_with: ["pasos vacíos", "paso desconocido", "onboarding duplicado"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://go.dev"]
verified_at: "2026-09-10"
```

## 2. Applicability

Referencia AUTHORED en memoria para checklist; no activa una franquicia ni prueba capacitación/evaluación. Conserva orden de alta de pasos, pero permite completar cualquier paso conocido sin prerequisites. El caller autoriza tenant y franquiciado; no hay persistencia ni aprobación empresarial.

## 3. Architecture contract

- **Ownership**: `internal/onboarding` gobierna el checklist; la activación del franquiciado es del dominio.
- **Invariantes**: (1) pasos no vacíos y deduplicados. (2) paso debe existir. (3) completado idempotente. (4) tenant-scoped.
- **Data flow**: `Start` → `Complete` → `Progress`/`Done`.
- **Failure modes**: pasos inválidos, duplicado, paso desconocido, no encontrado → error.
- **Seguridad/privacidad**: tenant-scoped.
- **Performance budget**: O(N) complete/progress sobre el número de pasos; Start O(N) esperado.

## 4. Exact file manifest

```text
CREATE internal/onboarding/onboarding.go
CREATE internal/onboarding/onboarding_test.go
```

## 5. Materialization blocks

### FILE: `internal/onboarding/onboarding.go`
```yaml
block_id: "GO-ONBOARDING-CORE:internal/onboarding/onboarding.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "57bc85e39c18cbea30ed7f3477160fdc0faf2f317ce807223f08bd9be2215e53"
variables: []
secrets_allowed: false
```
````go
// Package onboarding provides tenant-scoped franchisee onboarding: an ordered
// step checklist with progress tracking and idempotent step completion.
package onboarding

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

var (
	ErrInvalidSteps = errors.New("onboarding: invalid steps")
	ErrDuplicate    = errors.New("onboarding: duplicate onboarding")
	ErrNotFound     = errors.New("onboarding: not found")
	ErrUnknownStep  = errors.New("onboarding: unknown step")
)

// Onboarding is a franchisee's step checklist.
type Onboarding struct {
	TenantID     string
	FranchiseeID string
	Steps        []string
	Completed    map[string]bool
}

// Store holds onboardings.
type Store struct {
	mu    sync.Mutex
	items map[identity]*Onboarding
}

// NewStore returns an empty store.
func NewStore() *Store {
	return &Store{items: make(map[identity]*Onboarding)}
}

type identity struct{ tenant, franchisee string }

func key(tenant, franchisee string) identity { return identity{tenant: tenant, franchisee: franchisee} }

// Start begins an onboarding with deduplicated, non-empty steps.
func (s *Store) Start(tenant, franchisee string, steps []string) error {
	if strings.TrimSpace(tenant) == "" || strings.TrimSpace(franchisee) == "" || len(steps) == 0 {
		return ErrInvalidSteps
	}
	seen := make(map[string]bool)
	dedup := make([]string, 0, len(steps))
	for _, st := range steps {
		if strings.TrimSpace(st) == "" {
			return ErrInvalidSteps
		}
		if !seen[st] {
			seen[st] = true
			dedup = append(dedup, st)
		}
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.items[key(tenant, franchisee)]; ok {
		return ErrDuplicate
	}
	s.items[key(tenant, franchisee)] = &Onboarding{
		TenantID: tenant, FranchiseeID: franchisee, Steps: dedup, Completed: make(map[string]bool),
	}
	return nil
}

// Complete marks a step done (idempotent, must be a known step).
func (s *Store) Complete(tenant, franchisee, step string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	o, ok := s.items[key(tenant, franchisee)]
	if !ok {
		return ErrNotFound
	}
	found := false
	for _, st := range o.Steps {
		if st == step {
			found = true
			break
		}
	}
	if !found {
		return ErrUnknownStep
	}
	o.Completed[step] = true
	return nil
}

// Progress returns (done, total, ok).
func (s *Store) Progress(tenant, franchisee string) (int, int, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	o, ok := s.items[key(tenant, franchisee)]
	if !ok {
		return 0, 0, false
	}
	done := 0
	for _, st := range o.Steps {
		if o.Completed[st] {
			done++
		}
	}
	return done, len(o.Steps), true
}

// Done reports whether all steps are complete.
func (s *Store) Done(tenant, franchisee string) bool {
	done, total, ok := s.Progress(tenant, franchisee)
	return ok && done == total
}

// String aids debugging.
func (s *Store) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return fmt.Sprintf("onboarding(%d)", len(s.items))
}
````

### FILE: `internal/onboarding/onboarding_test.go`
```yaml
block_id: "GO-ONBOARDING-CORE:internal/onboarding/onboarding_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "94506a1638046a14f30fce623a017ef986e18962815e6c42f40fd3f5a007ee72"
variables: []
secrets_allowed: false
```
````go
package onboarding

import (
	"errors"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
)

func TestProgressAndDone(t *testing.T) {
	s := NewStore()
	if err := s.Start("t", "f1", []string{"alta", "pago", "capacitación"}); err != nil {
		t.Fatal(err)
	}
	done, total, ok := s.Progress("t", "f1")
	if !ok || done != 0 || total != 3 {
		t.Fatalf("expected 0/3, got %d/%d/%v", done, total, ok)
	}
	_ = s.Complete("t", "f1", "alta")
	_ = s.Complete("t", "f1", "pago")
	if s.Done("t", "f1") {
		t.Fatal("should not be done at 2/3")
	}
	_ = s.Complete("t", "f1", "capacitación")
	if !s.Done("t", "f1") {
		t.Fatal("should be done at 3/3")
	}
}

func TestCompleteIdempotent(t *testing.T) {
	s := NewStore()
	_ = s.Start("t", "f1", []string{"a"})
	_ = s.Complete("t", "f1", "a")
	_ = s.Complete("t", "f1", "a") // idempotent
	done, total, _ := s.Progress("t", "f1")
	if done != 1 || total != 1 {
		t.Fatalf("expected 1/1, got %d/%d", done, total)
	}
}

func TestUnknownStepRejected(t *testing.T) {
	s := NewStore()
	_ = s.Start("t", "f1", []string{"a"})
	if err := s.Complete("t", "f1", "zzz"); !errors.Is(err, ErrUnknownStep) {
		t.Fatalf("unknown step accepted: %v", err)
	}
}

func TestStartDedupAndInvalid(t *testing.T) {
	s := NewStore()
	_ = s.Start("t", "f1", []string{"a", "b", "a"})
	_, total, _ := s.Progress("t", "f1")
	if total != 2 {
		t.Fatalf("expected 2 dedup steps, got %d", total)
	}
	if err := s.Start("t", "f1", []string{"c"}); !errors.Is(err, ErrDuplicate) {
		t.Fatalf("duplicate onboarding accepted: %v", err)
	}
	if err := s.Start("t", "f2", nil); !errors.Is(err, ErrInvalidSteps) {
		t.Fatalf("empty steps accepted: %v", err)
	}
}

func TestTenantIsolation(t *testing.T) {
	s := NewStore()
	_ = s.Start("t", "f1", []string{"a"})
	if _, _, ok := s.Progress("other", "f1"); ok {
		t.Fatal("cross-tenant leak")
	}
	if err := s.Complete("other", "f1", "a"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("cross-tenant should be not found, got %v", err)
	}
}

func TestDistinctOnboardingTuplesDoNotCollide(t *testing.T) {
	s := NewStore()
	if e := s.Start("a\x00b", "c", []string{"step"}); e != nil {
		t.Fatal(e)
	}
	if e := s.Start("a", "b\x00c", []string{"other"}); e != nil {
		t.Fatalf("distinct onboarding rejected: %v", e)
	}
}
func TestAliasedOnboardingCannotReadOrComplete(t *testing.T) {
	s := NewStore()
	if e := s.Start("a\x00b", "c", []string{"step"}); e != nil {
		t.Fatal(e)
	}
	if n, total, ok := s.Progress("a", "b\x00c"); ok || n != 0 || total != 0 {
		t.Errorf("foreign progress: %d %d %v", n, total, ok)
	}
	if e := s.Complete("a", "b\x00c", "step"); !errors.Is(e, ErrNotFound) {
		t.Errorf("foreign completion: %v", e)
	}
	if s.Done("a", "b\x00c") {
		t.Error("foreign Done")
	}
	if n, total, ok := s.Progress("a\x00b", "c"); !ok || n != 0 || total != 1 {
		t.Errorf("original progress changed: %d %d %v", n, total, ok)
	}
}

func TestChecklistCopiesInputAndAllowsAnyKnownStep(t *testing.T) {
	s := NewStore()
	steps := []string{"first", "last", "first", " first"}
	if e := s.Start("t", "f", steps); e != nil {
		t.Fatal(e)
	}
	steps[0] = "mutated"
	if e := s.Complete("t", "f", "last"); e != nil {
		t.Fatal("existing completion order changed", e)
	}
	if e := s.Complete("t", "f", "mutated"); !errors.Is(e, ErrUnknownStep) {
		t.Fatal("caller slice aliases store", e)
	}
	if e := s.Complete("t", "f", "first"); e != nil {
		t.Fatal(e)
	}
	if e := s.Complete("t", "f", "first"); e != nil {
		t.Fatal("completion lost idempotence", e)
	}
	if n, total, ok := s.Progress("t", "f"); !ok || n != 2 || total != 3 {
		t.Fatal("dedup/exact step identity changed", n, total, ok)
	}
	if e := s.Complete("t", "f", " first"); e != nil {
		t.Fatal(e)
	}
	if !s.Done("t", "f") {
		t.Fatal("not done")
	}
}
func TestConcurrentChecklistStartAndIdempotentCompletion(t *testing.T) {
	s := NewStore()
	var wg sync.WaitGroup
	var starts atomic.Int64
	for i := 0; i < 64; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			e := s.Start("t", "f", []string{"a", "b"})
			if e == nil {
				starts.Add(1)
			} else if !errors.Is(e, ErrDuplicate) {
				t.Error(e)
			}
			_ = s.String()
		}()
	}
	wg.Wait()
	if starts.Load() != 1 {
		t.Fatal("multiple starts")
	}
	for i := 0; i < 64; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			step := "a"
			if i%2 == 0 {
				step = "b"
			}
			if e := s.Complete("t", "f", step); e != nil {
				t.Error(e)
			}
			s.Progress("t", "f")
			_ = s.String()
		}(i)
	}
	wg.Wait()
	if n, total, ok := s.Progress("t", "f"); !ok || n != 2 || total != 2 {
		t.Fatal("idempotent concurrent progress", n, total, ok)
	}
}
func FuzzChecklistProgressModel(f *testing.F) {
	for _, x := range []string{"a", "a\x00b", "", " ", "客户", "b\x00c"} {
		f.Add(x, []byte{0, 1, 0, 0, 0, 0, 1, 1, 1, 0, 1, 0, 1, 1, 0, 0})
		f.Add(x, []byte{0, 0, 0, 0, 1, 0, 0, 1, 1, 0, 0, 0, 1, 0, 0, 0, 2, 0, 0, 0})
	}
	f.Fuzz(func(t *testing.T, x string, ops []byte) {
		if len(x) > 64 {
			x = x[:64]
		}
		if len(ops) > 128 {
			ops = ops[:128]
		}
		tenants := []string{"a", "a\x00b", x, " "}
		people := []string{"c", "b\x00c", x, " "}
		stepNames := []string{"first", "last", x, " first", " "}
		plans := [][]string{{"first", "last", "first"}, {x, "first"}, nil, {" "}, {" first", "first"}}
		type entry struct {
			tenant, person string
			steps          []string
			done           []bool
		}
		var model []entry
		s := NewStore()
		lookup := func(tenant, person string) int {
			for i, e := range model {
				if e.tenant == tenant && e.person == person {
					return i
				}
			}
			return -1
		}
		for i := 0; i+3 < len(ops); i += 4 {
			tenant := tenants[int(ops[i+1])%len(tenants)]
			person := people[int(ops[i+2])%len(people)]
			index := lookup(tenant, person)
			var got, want error
			switch ops[i] % 3 {
			case 0:
				steps := plans[int(ops[i+3])%len(plans)]
				unique := []string{}
				if strings.TrimSpace(tenant) == "" || strings.TrimSpace(person) == "" || len(steps) == 0 {
					want = ErrInvalidSteps
				}
				for _, st := range steps {
					if strings.TrimSpace(st) == "" {
						want = ErrInvalidSteps
					}
					found := false
					for _, old := range unique {
						if st == old {
							found = true
						}
					}
					if !found {
						unique = append(unique, st)
					}
				}
				if want == nil && index >= 0 {
					want = ErrDuplicate
				}
				got = s.Start(tenant, person, steps)
				if want == nil {
					model = append(model, entry{tenant: tenant, person: person, steps: unique, done: make([]bool, len(unique))})
				}
			case 1:
				step := stepNames[int(ops[i+3])%len(stepNames)]
				if index < 0 {
					want = ErrNotFound
				} else {
					found := -1
					for j, st := range model[index].steps {
						if st == step {
							found = j
						}
					}
					if found < 0 {
						want = ErrUnknownStep
					} else {
						model[index].done[found] = true
					}
				}
				got = s.Complete(tenant, person, step)
			case 2:
				_ = s.String()
			}
			if !errors.Is(got, want) {
				t.Fatalf("step %d got %v want %v", i/4, got, want)
			}
			for _, tenant := range tenants {
				for _, person := range people {
					idx := lookup(tenant, person)
					n, total := 0, 0
					if idx >= 0 {
						total = len(model[idx].steps)
						for _, done := range model[idx].done {
							if done {
								n++
							}
						}
					}
					actual, all, ok := s.Progress(tenant, person)
					if actual != n || all != total || ok != (idx >= 0) {
						t.Fatal("progress differs from list model")
					}
					if s.Done(tenant, person) != (idx >= 0 && n == total) {
						t.Fatal("Done differs from list model")
					}
				}
			}
		}
	})
}
````


## 6. Configuration surface

Sin variables ni secretos. Los pasos son config del proyecto.

## 7. Dependency bill

| Package/image/tool | Pin exacto | Uso | Licencia | Runtime/build | Fuente oficial |
|---|---|---|---|---|---|
| Go standard library | go 1.26.7 (toolchain) | onboarding | BSD-3-Clause | runtime | https://go.dev |

## 8. Apply order

1. Componer el backend (mismo módulo).
2. Colocar los dos archivos bajo `internal/onboarding/`.
3. Verificar con `go test ./... -count=1` y `go vet ./...`.
4. Rollback: eliminar `internal/onboarding/`.

## 9. Verification

- `go test ./internal/onboarding/ -count=1`: 5/5 PASS (progreso/done, idempotente, paso desconocido, dedup/inválido, aislamiento tenant).
- `go test ./... -count=1` (28 paquetes): PASS.
- `go vet ./...`: exit 0.

## 10. Reconstruction evidence

Véase `reconstruction_evidence/GO_ONBOARDING_CORE_2026-09-02_V206.md`.

## Revisión V359 — identidad sin alias

0.1.1 usa tuplas tenant/identidad y preserva firmas, validación, prioridad de
errores y estados. AUTHORED/LicenseRef-Workspace-Owner, CONDITIONED y fuera
del perfil persisten. Compatibilidad histórica no acredita montaje actual.
La suite V206 es histórica, no composición reejecutada en V359.
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
