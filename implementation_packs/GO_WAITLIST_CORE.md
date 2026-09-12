# Go Waitlist Core

## 1. Metadata

```yaml
pack_id: "GO-WAITLIST-CORE"
pack_version: "0.1.3"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa la lista de espera tenant-scoped: join con posición por orden de llegada, y máquina de estados pending→notified→seated/cancelled fail-closed."
stacks: ["Go 1.26.8"]
compatible_with: ["uso aislado; integración de recordatorios no admitida"]
incompatible_with: ["sujeto duplicado", "transición inválida"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://go.dev", "https://franpos.com/franchise-back-office-royalty-management-guide/"]
verified_at: "2026-09-10"
```

## 2. Applicability

Referencia AUTHORED en memoria: posición entre pending por orden de Join exitoso dentro del tenant. Notify sólo marca estado local; Seat admite pending o notified. No reserva capacidad ni demuestra envío, oferta, expiración o persistencia. Identidades terminales no se reutilizan.

## 3. Architecture contract

- **Ownership**: `internal/waitlist` gestiona posición y estado; no contiene adapter de notificación real admitido (límite V355).
- **Invariantes**: (1) un entry por sujeto. (2) la posición cuenta sólo entries `pending` en orden de llegada. (3) transiciones gobernadas.
- **Data flow**: `Join` → `Position` / `Notify` / `Seat` / `Cancel`.
- **Failure modes**: duplicado, transición inválida, entrada inválida → error.
- **Seguridad/privacidad**: tenant-scoped.
- **Performance budget**: O(N) por `Position`.

## 4. Exact file manifest

```text
CREATE internal/waitlist/waitlist.go
CREATE internal/waitlist/waitlist_test.go
```

## 5. Materialization blocks

### FILE: `internal/waitlist/waitlist.go`
```yaml
block_id: "GO-WAITLIST-CORE:internal/waitlist/waitlist.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "a83273d7fa047571ebc5953e1a97a97c30170973bafddda6f97024c0fb907794"
variables: []
secrets_allowed: false
```
````go
// Package waitlist provides a tenant-scoped waitlist: join to reserve a
// position, then notify/seat/cancel through a governed state machine.
// Position is computed only over pending entries in join order. AUTHORED.
package waitlist

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"
)

// State is the waitlist lifecycle.
type State string

const (
	StatePending   State = "pending"
	StateNotified  State = "notified"
	StateSeated    State = "seated"
	StateCancelled State = "cancelled"
)

var (
	idRe = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:-]{0,127}$`)

	ErrInvalidEntry  = errors.New("waitlist: invalid entry")
	ErrDuplicate     = errors.New("waitlist: duplicate subject")
	ErrNotFound      = errors.New("waitlist: not found")
	ErrBadTransition = errors.New("waitlist: bad state transition")
)

// Entry is a single waitlist position.
type Entry struct {
	TenantID  string
	SubjectID string
	State     State
	At        time.Time
}

// Store holds entries in join order.
type Store struct {
	mu    sync.Mutex
	items map[identity]Entry
	order []identity // keys in join order
}

// NewStore returns an empty waitlist store.
func NewStore() *Store {
	return &Store{items: make(map[identity]Entry)}
}

type identity struct{ tenant, subject string }

func key(tenant, subject string) identity { return identity{tenant: tenant, subject: subject} }

// Join adds a subject to the waitlist (pending). One entry per subject.
func (s *Store) Join(tenant, subject string) error {
	if strings.TrimSpace(tenant) == "" || !idRe.MatchString(subject) {
		return ErrInvalidEntry
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	k := key(tenant, subject)
	if _, ok := s.items[k]; ok {
		return ErrDuplicate
	}
	s.items[k] = Entry{TenantID: tenant, SubjectID: subject, State: StatePending, At: time.Now().UTC()}
	s.order = append(s.order, k)
	return nil
}

// Notify marks pending → notified.
func (s *Store) Notify(tenant, subject string) error {
	return s.transition(tenant, subject, StatePending, StateNotified)
}

// Seat marks pending/notified → seated.
func (s *Store) Seat(tenant, subject string) error {
	return s.transitionAny(tenant, subject, StateSeated, StatePending, StateNotified)
}

// Cancel marks pending/notified → cancelled.
func (s *Store) Cancel(tenant, subject string) error {
	return s.transitionAny(tenant, subject, StateCancelled, StatePending, StateNotified)
}

func (s *Store) transition(tenant, subject string, from, to State) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	k := key(tenant, subject)
	e, ok := s.items[k]
	if !ok {
		return ErrNotFound
	}
	if e.State != from {
		return ErrBadTransition
	}
	e.State = to
	s.items[k] = e
	return nil
}

func (s *Store) transitionAny(tenant, subject string, to State, froms ...State) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	k := key(tenant, subject)
	e, ok := s.items[k]
	if !ok {
		return ErrNotFound
	}
	for _, f := range froms {
		if e.State == f {
			e.State = to
			s.items[k] = e
			return nil
		}
	}
	return ErrBadTransition
}

// Position returns the 1-based rank among pending entries for the tenant.
func (s *Store) Position(tenant, subject string) (int, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	k := key(tenant, subject)
	if _, ok := s.items[k]; !ok {
		return 0, false
	}
	pos := 0
	for _, ok := range s.order {
		e := s.items[ok]
		if e.TenantID == tenant && e.State == StatePending {
			pos++
			if ok == k {
				return pos, true
			}
		}
	}
	return 0, false
}

// State returns the current state.
func (s *Store) State(tenant, subject string) (State, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.items[key(tenant, subject)]
	return e.State, ok
}

// String aids debugging.
func (s *Store) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return fmt.Sprintf("waitlist(%d)", len(s.items))
}
````

### FILE: `internal/waitlist/waitlist_test.go`
```yaml
block_id: "GO-WAITLIST-CORE:internal/waitlist/waitlist_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "89a167dd428fcec6e91beb640e1e5b96a8e8dd929333deee6be452ed61b54963"
variables: []
secrets_allowed: false
```
````go
package waitlist

import (
	"errors"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
)

func TestJoinAndPosition(t *testing.T) {
	s := NewStore()
	if err := s.Join("t", "a"); err != nil {
		t.Fatal(err)
	}
	if err := s.Join("t", "b"); err != nil {
		t.Fatal(err)
	}
	if p, _ := s.Position("t", "a"); p != 1 {
		t.Fatalf("expected position 1, got %d", p)
	}
	if p, _ := s.Position("t", "b"); p != 2 {
		t.Fatalf("expected position 2, got %d", p)
	}
}

func TestDuplicateSubject(t *testing.T) {
	s := NewStore()
	if err := s.Join("t", "a"); err != nil {
		t.Fatal(err)
	}
	if err := s.Join("t", "a"); !errors.Is(err, ErrDuplicate) {
		t.Fatalf("expected duplicate, got %v", err)
	}
}

func TestNotifySeatExcludesFromPosition(t *testing.T) {
	s := NewStore()
	for _, x := range []string{"a", "b", "c"} {
		if err := s.Join("t", x); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Notify("t", "a"); err != nil {
		t.Fatal(err)
	}
	if err := s.Seat("t", "a"); err != nil {
		t.Fatal(err)
	}
	// a is now seated; b becomes position 1.
	if p, _ := s.Position("t", "b"); p != 1 {
		t.Fatalf("expected b at position 1, got %d", p)
	}
	if p, _ := s.Position("t", "c"); p != 2 {
		t.Fatalf("expected c at position 2, got %d", p)
	}
}

func TestBadTransition(t *testing.T) {
	s := NewStore()
	if err := s.Join("t", "a"); err != nil {
		t.Fatal(err)
	}
	if err := s.Cancel("t", "a"); err != nil {
		t.Fatal(err)
	}
	if err := s.Seat("t", "a"); !errors.Is(err, ErrBadTransition) {
		t.Fatalf("expected bad transition (seat a cancelled), got %v", err)
	}
}

func TestInvalidEntry(t *testing.T) {
	s := NewStore()
	if err := s.Join("", "a"); !errors.Is(err, ErrInvalidEntry) {
		t.Fatalf("empty tenant accepted: %v", err)
	}
	if err := s.Join("t", ""); !errors.Is(err, ErrInvalidEntry) {
		t.Fatalf("empty subject accepted: %v", err)
	}
}

func TestTenantIsolation(t *testing.T) {
	s := NewStore()
	if err := s.Join("t1", "a"); err != nil {
		t.Fatal(err)
	}
	if _, ok := s.State("t2", "a"); ok {
		t.Fatal("subject leaked across tenants")
	}
}

func TestAliasedWaitlistLookupDoesNotReadForeignState(t *testing.T) {
	s := NewStore()
	if e := s.Join("a\x00b", "c"); e != nil {
		t.Fatal(e)
	}
	// This subject cannot Join, but read/transition APIs used to accept its alias.
	if st, ok := s.State("a", "b\x00c"); ok || st != "" {
		t.Fatalf("foreign state %q %v", st, ok)
	}
}
func TestAliasedWaitlistCannotNotifySeatOrCancel(t *testing.T) {
	for _, action := range []string{"notify", "seat", "cancel"} {
		s := NewStore()
		if e := s.Join("a\x00b", "c"); e != nil {
			t.Fatal(e)
		}
		var e error
		switch action {
		case "notify":
			e = s.Notify("a", "b\x00c")
		case "seat":
			e = s.Seat("a", "b\x00c")
		case "cancel":
			e = s.Cancel("a", "b\x00c")
		}
		if !errors.Is(e, ErrNotFound) {
			t.Errorf("%s foreign effect: %v", action, e)
		}
		if st, ok := s.State("a\x00b", "c"); !ok || st != StatePending {
			t.Errorf("%s changed foreign state: %q %v", action, st, ok)
		}
		if pos, ok := s.Position("a\x00b", "c"); !ok || pos != 1 {
			t.Errorf("%s removed foreign pending position", action)
		}
	}
}

func TestPendingFIFOAndTerminalIdentityPreserved(t *testing.T) {
	s := NewStore()
	for _, pair := range [][2]string{{"t", "a"}, {"other", "a"}, {"t", "b"}, {"t", "c"}} {
		if e := s.Join(pair[0], pair[1]); e != nil {
			t.Fatal(e)
		}
	}
	if e := s.Notify("t", "a"); e != nil {
		t.Fatal(e)
	}
	if n, ok := s.Position("t", "b"); !ok || n != 1 {
		t.Fatal("notified counted as pending")
	}
	if e := s.Seat("t", "b"); e != nil {
		t.Fatal("pending-to-seated contract changed", e)
	}
	if n, ok := s.Position("t", "c"); !ok || n != 1 {
		t.Fatal("seated counted as pending")
	}
	if e := s.Cancel("t", "a"); e != nil {
		t.Fatal(e)
	}
	if e := s.Join("t", "a"); !errors.Is(e, ErrDuplicate) {
		t.Fatal("terminal identity unexpectedly reusable", e)
	}
	if n, ok := s.Position("other", "a"); !ok || n != 1 {
		t.Fatal("foreign queue changed")
	}
}
func TestConcurrentWaitlistJoinAndTerminalTransition(t *testing.T) {
	s := NewStore()
	var wg sync.WaitGroup
	var joins, changes atomic.Int64
	for i := 0; i < 64; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			e := s.Join("t", "a")
			if e == nil {
				joins.Add(1)
			} else if !errors.Is(e, ErrDuplicate) {
				t.Error(e)
			}
			_ = s.String()
		}()
	}
	wg.Wait()
	if joins.Load() != 1 {
		t.Fatal("multiple joins")
	}
	for i := 0; i < 64; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			var e error
			if i%2 == 0 {
				e = s.Seat("t", "a")
			} else {
				e = s.Cancel("t", "a")
			}
			if e == nil {
				changes.Add(1)
			} else if !errors.Is(e, ErrBadTransition) {
				t.Error(e)
			}
			s.State("t", "a")
			s.Position("t", "a")
		}(i)
	}
	wg.Wait()
	if changes.Load() != 1 {
		t.Fatal("multiple terminal transitions")
	}
}
func FuzzWaitlistFIFOModel(f *testing.F) {
	for _, x := range []string{"a", "a\x00b", "", " ", "客户", "b\x00c"} {
		f.Add(x, []byte{0, 1, 0, 1, 0, 1, 2, 0, 1, 3, 0, 1})
		f.Add(x, []byte{0, 0, 0, 0, 0, 2, 0, 2, 0, 1, 0, 0, 2, 0, 2, 3, 0, 0, 0, 0, 0})
	}
	f.Fuzz(func(t *testing.T, x string, ops []byte) {
		if len(x) > 64 {
			x = x[:64]
		}
		if len(ops) > 120 {
			ops = ops[:120]
		}
		tenants := []string{"a", "a\x00b", "other", x, " "}
		subjects := []string{"c", "b\x00c", "next", x, " "}
		s := NewStore()
		var model []Entry
		valid := func(id string) bool {
			if len(id) < 1 || len(id) > 128 {
				return false
			}
			for i, c := range []byte(id) {
				alpha := (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
				if !alpha && (i == 0 || (c != '.' && c != '_' && c != ':' && c != '-')) {
					return false
				}
			}
			return true
		}
		lookup := func(tenant, subject string) int {
			for i, r := range model {
				if r.TenantID == tenant && r.SubjectID == subject {
					return i
				}
			}
			return -1
		}
		for i := 0; i+2 < len(ops); i += 3 {
			tenant := tenants[int(ops[i+1])%len(tenants)]
			subject := subjects[int(ops[i+2])%len(subjects)]
			index := lookup(tenant, subject)
			var got, want error
			action := ops[i] % 5
			if action == 0 {
				if strings.TrimSpace(tenant) == "" || !valid(subject) {
					want = ErrInvalidEntry
				} else if index >= 0 {
					want = ErrDuplicate
				}
				got = s.Join(tenant, subject)
				if want == nil {
					model = append(model, Entry{TenantID: tenant, SubjectID: subject, State: StatePending})
				}
			} else if action < 4 {
				if index < 0 {
					want = ErrNotFound
				} else {
					st := model[index].State
					allowed := st == StatePending || (action != 1 && st == StateNotified)
					if !allowed {
						want = ErrBadTransition
					} else {
						switch action {
						case 1:
							model[index].State = StateNotified
						case 2:
							model[index].State = StateSeated
						case 3:
							model[index].State = StateCancelled
						}
					}
				}
				switch action {
				case 1:
					got = s.Notify(tenant, subject)
				case 2:
					got = s.Seat(tenant, subject)
				case 3:
					got = s.Cancel(tenant, subject)
				}
			} else {
				_ = s.String()
			}
			if !errors.Is(got, want) {
				t.Fatalf("step %d got=%v want=%v", i/3, got, want)
			}
			for _, tenant := range tenants {
				for _, subject := range subjects {
					index := lookup(tenant, subject)
					st, ok := s.State(tenant, subject)
					if ok != (index >= 0) {
						t.Fatal("state existence mismatch")
					}
					if index < 0 && st != "" {
						t.Fatal("absent state")
					}
					if index >= 0 && st != model[index].State {
						t.Fatal("state mismatch")
					}
					rank, wantRank := 0, 0
					for _, e := range model {
						if e.TenantID == tenant && e.State == StatePending {
							rank++
							if e.SubjectID == subject {
								wantRank = rank
							}
						}
					}
					pos, ok := s.Position(tenant, subject)
					if pos != wantRank || ok != (wantRank > 0) {
						t.Fatalf("FIFO got %d %v want %d", pos, ok, wantRank)
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
| Go standard library | go 1.26.7 (toolchain) | lista de espera | BSD-3-Clause | runtime | https://go.dev |

## 8. Apply order

1. Componer el backend (mismo módulo).
2. Colocar los dos archivos bajo `internal/waitlist/`.
3. Verificar con `go test ./... -count=1` y `go vet ./...`.
4. Rollback: eliminar `internal/waitlist/`.

## 9. Verification

- `go test ./internal/waitlist/ -count=1`: 6/6 PASS.
- `go test ./... -count=1` (37 paquetes): PASS.
- `go vet ./...`: exit 0.

## 10. Reconstruction evidence

Véase `reconstruction_evidence/GO_WAITLIST_CORE_2026-09-02_V213.md`.

V355: metadata0.1.1 retira compatibilidad con GO-REMINDERS-CORE0.1.x por FAIL624.
Sus2fuentes siguen idénticas; no hay imports/calls a reminders. La API0.2.0 exige
tenant explícito y no demuestra envío. Ningún adapter waitlist→reminders ni
promoción durable fue implementado o aprobado. Evidencia: reconstruction_evidence/REMINDER_TENANT_API_V355.md.

## Revisión V358 — identidad sin alias

0.1.2 usa tuplas de identidad, preserva firmas, validación y transiciones.
AUTHORED/LicenseRef-Workspace-Owner, CONDITIONED y fuera del perfil persisten.
Suite V213 es histórica, no composición reejecutada. No integración ni
permiso externo inferido. Evidencia: reconstruction_evidence/MODERATION_WAITLIST_ISOLATION_V358.md.

V358 verificado entre ambos cores:4/4fuentes reconstruidas,11tests históricos
intactos,4regresiones red,19tests x3,vet/build y24semillas/4216603ejecuciones
fuzz PASS. Concurrencia local probada; sin -race, auth ni efectos de proveedor.

## Auditoría por claim V374

La revisión 0.1.3 se reconstruye y prueba en Go1.26.8/Windows con fuentes
oficiales fijadas y comparaciones por contrato. Ver
`reconstruction_evidence/CORE_CLAIM_ADMISSION_V374.md` desde la raíz: identidad,
licencia local, prueba por invariante, fuente semántica y dictamen de equivalencia
se mantienen separados. El lenguaje/stdlib no es fuente de un negocio empresarial.
No cambia este claim, API, permiso de composición ni las condiciones de integración.
La revisión 0.1.2 y sus bytes quedan preservados en el expediente anterior.
