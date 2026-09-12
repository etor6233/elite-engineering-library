# Go Referrals Core

## 1. Metadata

```yaml
pack_id: "GO-REFERRALS-CORE"
pack_version: "0.1.2"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa el programa de referidos tenant-scoped: códigos únicos, sin auto-referido, un referido por persona y máquina de estados pending→converted→rewarded fail-closed."
stacks: ["Go 1.26.8"]
compatible_with: ["GO-LOYALTY-CORE 0.1.x", "GO-ENTERPRISE-BACKEND 0.4.x"]
incompatible_with: ["auto-referido", "referido duplicado", "recompensa sin conversión"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://go.dev"]
verified_at: "2026-09-10"
```

## 2. Applicability

Referencia AUTHORED en memoria: un registro por referee dentro de cada tenant y estados locales. Convert/Reward no prueban conversión comercial ni emiten recompensa, crédito o cupón. El caller debe autenticar y autorizar tenant/identidad; no hay persistencia, antifraude ni entrega durable. Rechaza auto-referidos y transiciones fuera de orden.

## 3. Architecture contract

- **Ownership**: `internal/referrals` gobierna el registro y la máquina de estados; la recompensa (loyalty/cupón) es composición del dominio.
- **Invariantes**: (1) referrer ≠ referee. (2) código y referido únicos dentro del tenant; tuplas distintas no comparten identidad. (3) estados: pending→converted→rewarded. (4) tenant-scoped.
- **Data flow**: `Create` → `Convert` → `Reward`.
- **Failure modes**: auto-referido, duplicado, transición inválida → error.
- **Seguridad/privacidad**: tenant-scoped.
- **Performance budget**: O(1).

## 4. Exact file manifest

```text
CREATE internal/referrals/referrals.go
CREATE internal/referrals/referrals_test.go
```

## 5. Materialization blocks

### FILE: `internal/referrals/referrals.go`
```yaml
block_id: "GO-REFERRALS-CORE:internal/referrals/referrals.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "23f3fb2be5f341caf748b8bcf364ffcff4e571e234549c84e8b312a0e9d83231"
variables: []
secrets_allowed: false
```
````go
// Package referrals provides tenant-scoped referral tracking: unique codes,
// no self-referral, one referral per referee within a tenant, and a governed state machine
// pending → converted → rewarded (fail-closed).
package referrals

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"
)

// State is the referral lifecycle.
type State string

const (
	StatePending   State = "pending"
	StateConverted State = "converted"
	StateRewarded  State = "rewarded"
)

var (
	codeRe = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]{0,31}$`)

	ErrInvalidReferral  = errors.New("referrals: invalid referral")
	ErrSelfReferral     = errors.New("referrals: self-referral forbidden")
	ErrDuplicateCode    = errors.New("referrals: duplicate code")
	ErrDuplicateReferee = errors.New("referrals: referee already referred")
	ErrNotFound         = errors.New("referrals: not found")
	ErrBadTransition    = errors.New("referrals: bad state transition")
)

// Referral links a referrer code to a referee.
type Referral struct {
	TenantID   string
	Code       string
	ReferrerID string
	RefereeID  string
	State      State
	At         time.Time
}

// Store is the referral registry.
type Store struct {
	mu        sync.Mutex
	byCode    map[identity]Referral
	byReferee map[identity]string // (tenant, referee) -> code
}

// NewStore returns an empty registry.
func NewStore() *Store {
	return &Store{byCode: make(map[identity]Referral), byReferee: make(map[identity]string)}
}

type identity struct{ tenant, id string }

func ck(tenant, id string) identity { return identity{tenant: tenant, id: id} }

// Create registers a referral, fail-closed on self-referral or duplicates.
func (s *Store) Create(r Referral) error {
	if err := r.validate(); err != nil {
		return err
	}
	if r.ReferrerID == r.RefereeID {
		return ErrSelfReferral
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.byCode[ck(r.TenantID, r.Code)]; ok {
		return ErrDuplicateCode
	}
	if _, ok := s.byReferee[ck(r.TenantID, r.RefereeID)]; ok {
		return ErrDuplicateReferee
	}
	r.State = StatePending
	if r.At.IsZero() {
		r.At = time.Now().UTC()
	}
	s.byCode[ck(r.TenantID, r.Code)] = r
	s.byReferee[ck(r.TenantID, r.RefereeID)] = r.Code
	return nil
}

// Convert marks pending → converted.
func (s *Store) Convert(tenant, refereeID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	code, ok := s.byReferee[ck(tenant, refereeID)]
	if !ok {
		return ErrNotFound
	}
	r := s.byCode[ck(tenant, code)]
	if r.State != StatePending {
		return ErrBadTransition
	}
	r.State = StateConverted
	s.byCode[ck(tenant, code)] = r
	return nil
}

// Reward marks converted → rewarded (fail-closed: must be converted first).
func (s *Store) Reward(tenant, refereeID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	code, ok := s.byReferee[ck(tenant, refereeID)]
	if !ok {
		return ErrNotFound
	}
	r := s.byCode[ck(tenant, code)]
	if r.State != StateConverted {
		return ErrBadTransition
	}
	r.State = StateRewarded
	s.byCode[ck(tenant, code)] = r
	return nil
}

// State returns the current state.
func (s *Store) State(tenant, refereeID string) (State, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	code, ok := s.byReferee[ck(tenant, refereeID)]
	if !ok {
		return "", false
	}
	return s.byCode[ck(tenant, code)].State, true
}

func (r Referral) validate() error {
	if strings.TrimSpace(r.TenantID) == "" {
		return fmt.Errorf("%w: tenant", ErrInvalidReferral)
	}
	if !codeRe.MatchString(r.Code) {
		return fmt.Errorf("%w: code", ErrInvalidReferral)
	}
	if strings.TrimSpace(r.ReferrerID) == "" || strings.TrimSpace(r.RefereeID) == "" {
		return fmt.Errorf("%w: ids", ErrInvalidReferral)
	}
	return nil
}
````

### FILE: `internal/referrals/referrals_test.go`
```yaml
block_id: "GO-REFERRALS-CORE:internal/referrals/referrals_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "9314db0d1bf1730673bf51105d23959367309c671aa29d3abf9f7a493e399219"
variables: []
secrets_allowed: false
```
````go
package referrals

import (
	"errors"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
)

func ref(code, referrer, referee string) Referral {
	return Referral{TenantID: "t", Code: code, ReferrerID: referrer, RefereeID: referee}
}

func TestCreateAndStateMachine(t *testing.T) {
	s := NewStore()
	if err := s.Create(ref("R1", "a", "b")); err != nil {
		t.Fatal(err)
	}
	st, ok := s.State("t", "b")
	if !ok || st != StatePending {
		t.Fatalf("expected pending, got %q/%v", st, ok)
	}
	if err := s.Convert("t", "b"); err != nil {
		t.Fatal(err)
	}
	if err := s.Reward("t", "b"); err != nil {
		t.Fatal(err)
	}
	if st, _ := s.State("t", "b"); st != StateRewarded {
		t.Fatalf("expected rewarded, got %q", st)
	}
}

func TestRewardRequiresConverted(t *testing.T) {
	s := NewStore()
	_ = s.Create(ref("R1", "a", "b"))
	if err := s.Reward("t", "b"); !errors.Is(err, ErrBadTransition) {
		t.Fatalf("expected ErrBadTransition, got %v", err)
	}
}

func TestSelfReferralRejected(t *testing.T) {
	s := NewStore()
	if err := s.Create(ref("R1", "a", "a")); !errors.Is(err, ErrSelfReferral) {
		t.Fatalf("expected self-referral rejection, got %v", err)
	}
}

func TestDuplicateCodeAndReferee(t *testing.T) {
	s := NewStore()
	_ = s.Create(ref("R1", "a", "b"))
	if err := s.Create(ref("R1", "c", "d")); !errors.Is(err, ErrDuplicateCode) {
		t.Fatalf("duplicate code accepted: %v", err)
	}
	if err := s.Create(ref("R2", "a", "b")); !errors.Is(err, ErrDuplicateReferee) {
		t.Fatalf("duplicate referee accepted: %v", err)
	}
}

func TestTenantIsolation(t *testing.T) {
	s := NewStore()
	_ = s.Create(ref("R1", "a", "b"))
	if _, ok := s.State("other", "b"); ok {
		t.Fatal("cross-tenant leak")
	}
	if err := s.Convert("other", "b"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("cross-tenant should be not found, got %v", err)
	}
}

func identityReferral(tenant, code, referee string) Referral {
	return Referral{TenantID: tenant, Code: code, ReferrerID: "source", RefereeID: referee}
}
func TestDistinctReferralTuplesDoNotCollide(t *testing.T) {
	s := NewStore()
	if e := s.Create(identityReferral("a\x00b", "FIRST", "c")); e != nil {
		t.Fatal(e)
	}
	if e := s.Create(identityReferral("a", "SECOND", "b\x00c")); e != nil {
		t.Fatalf("distinct identity rejected: %v", e)
	}
}
func TestAbsentAliasedRefereeIsNotFound(t *testing.T) {
	s := NewStore()
	if e := s.Create(identityReferral("a\x00b", "SHARED", "c")); e != nil {
		t.Fatal(e)
	}
	if st, ok := s.State("a", "b\x00c"); ok || st != "" {
		t.Errorf("absent referee reported present: %q %v", st, ok)
	}
	if e := s.Convert("a", "b\x00c"); !errors.Is(e, ErrNotFound) {
		t.Errorf("convert got %v", e)
	}
	if e := s.Reward("a", "b\x00c"); !errors.Is(e, ErrNotFound) {
		t.Errorf("reward got %v", e)
	}
}
func TestAliasedLookupCannotTransitionAnotherReferee(t *testing.T) {
	s := NewStore()
	for _, r := range []Referral{identityReferral("a\x00b", "SHARED", "c"), identityReferral("a", "SHARED", "legitimate")} {
		if e := s.Create(r); e != nil {
			t.Fatal(e)
		}
	}
	if st, ok := s.State("a", "b\x00c"); ok || st != "" {
		t.Errorf("read unrelated referral: %q %v", st, ok)
	}
	if e := s.Convert("a", "b\x00c"); !errors.Is(e, ErrNotFound) {
		t.Errorf("converted unrelated referral: %v", e)
	}
	if e := s.Reward("a", "b\x00c"); !errors.Is(e, ErrNotFound) {
		t.Errorf("rewarded unrelated referral: %v", e)
	}
	for _, p := range [][2]string{{"a\x00b", "c"}, {"a", "legitimate"}} {
		if st, ok := s.State(p[0], p[1]); !ok || st != StatePending {
			t.Errorf("other identity changed: %q %v", st, ok)
		}
	}
}

func TestRejectedCreatesDoNotConsumeIdentity(t *testing.T) {
	for _, r := range []Referral{
		{TenantID: " ", Code: "OK", ReferrerID: "source", RefereeID: "dest"},
		{TenantID: "t", Code: "bad code", ReferrerID: "source", RefereeID: "dest"},
		{TenantID: "t", Code: "OK", ReferrerID: " ", RefereeID: "dest"},
		{TenantID: "t", Code: "OK", ReferrerID: "source", RefereeID: " "},
		{TenantID: "t", Code: "OK", ReferrerID: "dest", RefereeID: "dest"},
	} {
		s := NewStore()
		e := s.Create(r)
		want := ErrInvalidReferral
		if r.ReferrerID == r.RefereeID {
			want = ErrSelfReferral
		}
		if !errors.Is(e, want) {
			t.Fatalf("wrong reject %v", e)
		}
		if _, ok := s.State(r.TenantID, r.RefereeID); ok {
			t.Fatal("rejected create stored")
		}
		if e := s.Create(identityReferral("t", "OK", "dest")); e != nil {
			t.Fatal("rejection consumed identity", e)
		}
	}
	s := NewStore()
	r := identityReferral("t", "OK", "dest")
	r.State = StateRewarded
	if e := s.Create(r); e != nil {
		t.Fatal(e)
	}
	if st, _ := s.State("t", "dest"); st != StatePending {
		t.Fatal("caller bypassed pending")
	}
	r.ReferrerID = "dest"
	if e := s.Create(r); !errors.Is(e, ErrSelfReferral) {
		t.Fatal("validation priority changed", e)
	}
	r.ReferrerID = "source"
	if e := s.Create(r); !errors.Is(e, ErrDuplicateCode) {
		t.Fatal("duplicate-code priority changed", e)
	}
}
func TestConcurrentReferralLifecycleAtMostOneTransition(t *testing.T) {
	s := NewStore()
	var success atomic.Int64
	race := func(action func() error, want error) {
		success.Store(0)
		var wg sync.WaitGroup
		for i := 0; i < 64; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				e := action()
				if e == nil {
					success.Add(1)
				} else if !errors.Is(e, want) {
					t.Error(e)
				}
			}()
		}
		wg.Wait()
		if success.Load() != 1 {
			t.Fatalf("expected one success, got %d", success.Load())
		}
	}
	race(func() error { return s.Create(identityReferral("t", "CODE", "dest")) }, ErrDuplicateCode)
	race(func() error { return s.Convert("t", "dest") }, ErrBadTransition)
	race(func() error { return s.Reward("t", "dest") }, ErrBadTransition)
	if st, ok := s.State("t", "dest"); !ok || st != StateRewarded {
		t.Fatal("final state", st, ok)
	}
	if e := s.Create(identityReferral("other", "CODE", "dest")); e != nil {
		t.Fatal("foreign identity consumed", e)
	}
}

// The oracle uses a linear list and direct field comparison, not store keys or
// its validation/transition helpers. Each step checks every candidate identity.
func FuzzReferralSequenceModel(f *testing.F) {
	for _, p := range [][2]string{{"a", "b"}, {"a\x00b", "c"}, {"", " "}, {"á", "客户"}, {"\x00", "\x00"}, {"a", "a"}} {
		f.Add(p[0], p[1], []byte{0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 2, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 1, 1, 0, 2})
		f.Add(p[0], p[1], []byte{0, 2, 1, 0, 0, 0, 1, 2, 1, 0, 1, 1, 2, 0, 0, 2, 1, 2, 0, 0})
	}
	f.Fuzz(func(t *testing.T, x, y string, ops []byte) {
		if len(x) > 64 {
			x = x[:64]
		}
		if len(y) > 64 {
			y = y[:64]
		}
		if len(ops) > 160 {
			ops = ops[:160]
		}
		tenants := []string{"t", x, x + "\x00" + y, "other", " ", ""}
		people := []string{"dest", y, y + "\x00dest", x, "source", " "}
		codes := []string{"CODE", "SECOND", "bad code", "", "A-._9"}
		s := NewStore()
		var model []Referral
		validCode := func(code string) bool {
			if len(code) < 1 || len(code) > 32 {
				return false
			}
			for i, c := range []byte(code) {
				alnum := (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
				if !alnum && (i == 0 || (c != '.' && c != '_' && c != '-')) {
					return false
				}
			}
			return true
		}
		lookup := func(tenant, person string) int {
			for i, r := range model {
				if r.TenantID == tenant && r.RefereeID == person {
					return i
				}
			}
			return -1
		}
		for i := 0; i+4 < len(ops); i += 5 {
			tenant := tenants[int(ops[i+1])%len(tenants)]
			person := people[int(ops[i+2])%len(people)]
			var want, got error
			switch ops[i] % 4 {
			case 0:
				r := Referral{TenantID: tenant, Code: codes[int(ops[i+3])%len(codes)], ReferrerID: people[int(ops[i+4])%len(people)], RefereeID: person, State: StateRewarded}
				if strings.TrimSpace(r.TenantID) == "" || !validCode(r.Code) || strings.TrimSpace(r.ReferrerID) == "" || strings.TrimSpace(r.RefereeID) == "" {
					want = ErrInvalidReferral
				} else if r.ReferrerID == r.RefereeID {
					want = ErrSelfReferral
				} else {
					for _, old := range model {
						if old.TenantID == r.TenantID && old.Code == r.Code {
							want = ErrDuplicateCode
							break
						}
					}
					if want == nil && lookup(tenant, person) >= 0 {
						want = ErrDuplicateReferee
					}
				}
				got = s.Create(r)
				if want == nil {
					r.State = StatePending
					model = append(model, r)
				}
			case 1, 2:
				index := lookup(tenant, person)
				previous, next := StatePending, StateConverted
				if ops[i]%4 == 2 {
					previous, next = StateConverted, StateRewarded
				}
				if index < 0 {
					want = ErrNotFound
				} else if model[index].State != previous {
					want = ErrBadTransition
				} else {
					model[index].State = next
				}
				if ops[i]%4 == 1 {
					got = s.Convert(tenant, person)
				} else {
					got = s.Reward(tenant, person)
				}
			case 3:
				// State reads are checked exhaustively below and must have no effect.
				s.State(tenant, person)
			}
			if !errors.Is(got, want) {
				t.Fatalf("step %d error got=%v want=%v", i/5, got, want)
			}
			for _, tenant := range tenants {
				for _, person := range people {
					index := lookup(tenant, person)
					st, ok := s.State(tenant, person)
					if ok != (index >= 0) {
						t.Fatalf("existence mismatch %q %q", tenant, person)
					}
					if index >= 0 && st != model[index].State {
						t.Fatal("state mismatch")
					}
					if index < 0 && st != "" {
						t.Fatal("absent state not empty")
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
| Go standard library | go 1.26.7 (toolchain) | referidos | BSD-3-Clause | runtime | https://go.dev |

## 8. Apply order

1. Componer el backend (mismo módulo).
2. Colocar los dos archivos bajo `internal/referrals/`.
3. Verificar con `go test ./... -count=1` y `go vet ./...`.
4. Rollback: eliminar `internal/referrals/`.

## 9. Verification

- `go test ./internal/referrals/ -count=1`: 5/5 PASS (máquina de estados, reward requiere converted, auto-referido rechazado, código/referido duplicados, aislamiento tenant).
- `go test ./... -count=1` (18 paquetes): PASS.
- `go vet ./...`: exit 0.

## 10. Reconstruction evidence

Véase `reconstruction_evidence/GO_REFERRALS_CORE_2026-09-02_V196.md`.

## Revisión V357 — aislamiento de identidad

0.1.1 reemplaza concatenación ambigua por claves estructuradas tenant/id.
API y orden de validación/errores permanecen. Create reinicia State a pending;
los duplicados y transiciones repetidas retornan error, sin replay exitoso.
AUTHORED/LicenseRef-Workspace-Owner y CONDITIONED persisten; compatibilidad
histórica no acredita integración actual con loyalty ni issuance de recompensa.
La suite V196 es evidencia histórica, no una composición reejecutada en V357.
Evidencia nueva: reconstruction_evidence/REFERRAL_IDENTITY_ISOLATION_V357.md.

V357 verificado:2/2fuentes reconstruidas exactas;5tests históricos preservados,
3regresiones antes rojas y10tests x3 PASS;vet/build,12semillas y1709469
ejecuciones fuzz (10s,4workers) PASS. No -race, auth ni recompensa real probados.

## Auditoría por claim V374

La revisión 0.1.2 se reconstruye y prueba en Go1.26.8/Windows con fuentes
oficiales fijadas y comparaciones por contrato. Ver
`reconstruction_evidence/CORE_CLAIM_ADMISSION_V374.md` desde la raíz: identidad,
licencia local, prueba por invariante, fuente semántica y dictamen de equivalencia
se mantienen separados. El lenguaje/stdlib no es fuente de un negocio empresarial.
No cambia este claim, API, permiso de composición ni las condiciones de integración.
La revisión 0.1.1 y sus bytes quedan preservados en el expediente anterior.
