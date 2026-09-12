# Go Applicant Onboarding Core

## 1. Metadata

```yaml
pack_id: "GO-APPLICANT-ONBOARDING-CORE"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa el alta segura anti-fake (franquiciado o cliente online): submitted→contact_verified→approved→activated|rejected, con deduplicación de contacto, velocidad anti-falsos, doble control y activación sólo desde aprobado+verificado."
stacks: ["Go 1.26.7"]
compatible_with: ["GO-OTP-VERIFICATION-CORE 0.1.x", "GO-HUMAN-APPROVAL-CORE 0.1.x", "GO-ONBOARDING-CORE 0.1.x"]
incompatible_with: ["email/teléfono duplicado", "velocidad excedida", "auto-aprobación", "activación sin aprobación"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://go.dev", "https://sre.google"]
verified_at: "2026-09-02"
```

## 2. Applicability

Use como frontera de alta para franquiciado y para cliente que compra online. Rechace para contacto duplicado, velocidad excedida, auto-aprobación o activación sin aprobación.

## 3. Architecture contract

- **Ownership**: `internal/applicant` gobierna el ciclo; el KYC documental/biométrico real es un adapter CONDITIONED (proveedor externo), no se inventa.
- **Invariantes**: (1) contacto único por tenant. (2) velocidad acotada por tenant. (3) aprobación por revisor ≠ solicitante. (4) activación sólo desde approved (que exige contact_verified).
- **Data flow**: `Submit` → `MarkContactVerified` → `Approve` → `Activate` (o `Reject`).
- **Failure modes**: duplicado, velocidad, auto-aprobación, transición inválida, sin motivo → error.
- **Seguridad/privacidad**: tenant-scoped.
- **Performance budget**: O(1) por transición.

## 4. Exact file manifest

```text
CREATE internal/applicant/applicant.go
CREATE internal/applicant/applicant_test.go
```

## 5. Materialization blocks

### FILE: `internal/applicant/applicant.go`
```yaml
block_id: "GO-APPLICANT-ONBOARDING-CORE:internal/applicant/applicant.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "dd7a75914ea9dd00fd2aeecd9208cb185cf4294d11951f5f5a6bec45fec15a6f"
variables: []
secrets_allowed: false
```
````go
// Package applicant provides the secure, anti-fake registration lifecycle for
// a franchisee applicant or an online customer: submitted → contact_verified →
// approved → activated | rejected. It enforces duplicate-contact detection,
// submission velocity, dual control (reviewer != submitter) and activation
// only from an approved+verified applicant. Stdlib-only. AUTHORED.
package applicant

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"
)

// State is the registration lifecycle.
type State string

const (
	StateSubmitted       State = "submitted"
	StateContactVerified State = "contact_verified"
	StateApproved        State = "approved"
	StateActivated       State = "activated"
	StateRejected        State = "rejected"
)

var (
	idRe    = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:-]{0,127}$`)
	emailRe = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)
	phoneRe = regexp.MustCompile(`^\+?[0-9][0-9 .-]{4,31}$`)

	ErrInvalidApplicant = errors.New("applicant: invalid applicant")
	ErrDuplicateContact = errors.New("applicant: duplicate email or phone")
	ErrVelocityExceeded = errors.New("applicant: submission velocity exceeded")
	ErrNotFound         = errors.New("applicant: not found")
	ErrBadTransition    = errors.New("applicant: bad state transition")
	ErrSelfApproval     = errors.New("applicant: reviewer must differ from submitter")
	ErrReasonRequired   = errors.New("applicant: reject reason required")
)

// Applicant is a registration request.
type Applicant struct {
	TenantID     string
	ID           string
	Kind         string // "franchisee" or "customer"
	Email        string
	Phone        string // optional
	SubmittedBy  string // actor who submitted (self or an operator)
	State        State
	RejectReason string
}

// Store holds applicants and anti-fake counters.
type Store struct {
	mu             sync.Mutex
	maxSubmissions int
	window         time.Duration
	items          map[string]Applicant
	byEmail        map[string]string
	byPhone        map[string]string
	submissions    map[string][]time.Time // tenant -> submission times
	now            func() time.Time
}

// NewStore returns a store with the given per-tenant velocity limit.
func NewStore(maxSubmissions int, window time.Duration) (*Store, error) {
	if maxSubmissions < 1 || window <= 0 {
		return nil, ErrInvalidApplicant
	}
	return &Store{
		maxSubmissions: maxSubmissions,
		window:         window,
		items:          make(map[string]Applicant),
		byEmail:        make(map[string]string),
		byPhone:        make(map[string]string),
		submissions:    make(map[string][]time.Time),
		now:            time.Now,
	}, nil
}

func (s *Store) setClock(fn func() time.Time) { s.now = fn }

func validKind(k string) bool { return k == "franchisee" || k == "customer" }

func key(tenant, id string) string { return tenant + "\x00" + id }

// Submit registers a new applicant, fail-closed on duplicate contact, bad
// input or velocity exceeded.
func (s *Store) Submit(a Applicant) error {
	if strings.TrimSpace(a.TenantID) == "" || !idRe.MatchString(a.ID) ||
		!validKind(a.Kind) || !emailRe.MatchString(a.Email) || !idRe.MatchString(a.SubmittedBy) {
		return ErrInvalidApplicant
	}
	if a.Phone != "" && !phoneRe.MatchString(a.Phone) {
		return ErrInvalidApplicant
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	ekey := a.TenantID + "\x00email\x00" + a.Email
	if _, ok := s.byEmail[ekey]; ok {
		return ErrDuplicateContact
	}
	if a.Phone != "" {
		pkey := a.TenantID + "\x00phone\x00" + a.Phone
		if _, ok := s.byPhone[pkey]; ok {
			return ErrDuplicateContact
		}
	}

	now := s.now()
	cutoff := now.Add(-s.window)
	var kept []time.Time
	for _, ts := range s.submissions[a.TenantID] {
		if ts.After(cutoff) {
			kept = append(kept, ts)
		}
	}
	if len(kept) >= s.maxSubmissions {
		s.submissions[a.TenantID] = kept
		return ErrVelocityExceeded
	}
	s.submissions[a.TenantID] = append(kept, now)

	a.State = StateSubmitted
	s.items[key(a.TenantID, a.ID)] = a
	s.byEmail[ekey] = a.ID
	if a.Phone != "" {
		s.byPhone[a.TenantID+"\x00phone\x00"+a.Phone] = a.ID
	}
	return nil
}

// MarkContactVerified moves submitted → contact_verified (after OTP).
func (s *Store) MarkContactVerified(tenant, id string) error {
	return s.transition(tenant, id, StateSubmitted, StateContactVerified)
}

// Approve moves contact_verified → approved; the reviewer must differ from the
// submitter (dual control) and contact must be verified first.
func (s *Store) Approve(tenant, id, reviewer string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	a, ok := s.items[key(tenant, id)]
	if !ok {
		return ErrNotFound
	}
	if a.State != StateContactVerified {
		return ErrBadTransition
	}
	if reviewer == a.SubmittedBy {
		return ErrSelfApproval
	}
	a.State = StateApproved
	s.items[key(tenant, id)] = a
	return nil
}

// Activate moves approved → activated. Access/money is never granted before
// this point (and only from a verified+approved applicant).
func (s *Store) Activate(tenant, id string) error {
	return s.transition(tenant, id, StateApproved, StateActivated)
}

// Reject rejects any non-terminal applicant; reason and dual control required.
func (s *Store) Reject(tenant, id, reviewer, reason string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	a, ok := s.items[key(tenant, id)]
	if !ok {
		return ErrNotFound
	}
	if a.State == StateActivated || a.State == StateRejected {
		return ErrBadTransition
	}
	if strings.TrimSpace(reason) == "" {
		return ErrReasonRequired
	}
	if reviewer == a.SubmittedBy {
		return ErrSelfApproval
	}
	a.State = StateRejected
	a.RejectReason = reason
	s.items[key(tenant, id)] = a
	return nil
}

func (s *Store) transition(tenant, id string, from, to State) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	a, ok := s.items[key(tenant, id)]
	if !ok {
		return ErrNotFound
	}
	if a.State != from {
		return ErrBadTransition
	}
	a.State = to
	s.items[key(tenant, id)] = a
	return nil
}

// State returns the current state.
func (s *Store) State(tenant, id string) (State, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	a, ok := s.items[key(tenant, id)]
	return a.State, ok
}

// String aids debugging.
func (s *Store) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return fmt.Sprintf("applicant(%d)", len(s.items))
}
````

### FILE: `internal/applicant/applicant_test.go`
```yaml
block_id: "GO-APPLICANT-ONBOARDING-CORE:internal/applicant/applicant_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "5f06bda12ef32e80418b730afabb9c2d3461cb8da3f84df334120ea100cfe32c"
variables: []
secrets_allowed: false
```
````go
package applicant

import (
	"errors"
	"testing"
	"time"
)

func validApplicant(tenant, id, email string) Applicant {
	return Applicant{TenantID: tenant, ID: id, Kind: "franchisee", Email: email, SubmittedBy: "self"}
}

func TestSubmitAndVerify(t *testing.T) {
	s, _ := NewStore(10, time.Hour)
	a := validApplicant("t", "a1", "x@y.z")
	if err := s.Submit(a); err != nil {
		t.Fatal(err)
	}
	if st, _ := s.State("t", "a1"); st != StateSubmitted {
		t.Fatalf("expected submitted, got %s", st)
	}
	if err := s.MarkContactVerified("t", "a1"); err != nil {
		t.Fatal(err)
	}
	if st, _ := s.State("t", "a1"); st != StateContactVerified {
		t.Fatalf("expected contact_verified, got %s", st)
	}
}

func TestDuplicateContact(t *testing.T) {
	s, _ := NewStore(10, time.Hour)
	if err := s.Submit(validApplicant("t", "a1", "x@y.z")); err != nil {
		t.Fatal(err)
	}
	if err := s.Submit(validApplicant("t", "a2", "x@y.z")); !errors.Is(err, ErrDuplicateContact) {
		t.Fatalf("duplicate email accepted: %v", err)
	}
	a2 := validApplicant("t", "a2", "other@y.z")
	a2.Phone = "+5491100000000"
	if err := s.Submit(a2); err != nil {
		t.Fatal(err)
	}
	a3 := validApplicant("t", "a3", "third@y.z")
	a3.Phone = "+5491100000000"
	if err := s.Submit(a3); !errors.Is(err, ErrDuplicateContact) {
		t.Fatalf("duplicate phone accepted: %v", err)
	}
}

func TestApproveDualControl(t *testing.T) {
	s, _ := NewStore(10, time.Hour)
	if err := s.Submit(validApplicant("t", "a1", "x@y.z")); err != nil {
		t.Fatal(err)
	}
	if err := s.MarkContactVerified("t", "a1"); err != nil {
		t.Fatal(err)
	}
	// Same reviewer as submitter → rejected (dual control).
	if err := s.Approve("t", "a1", "self"); !errors.Is(err, ErrSelfApproval) {
		t.Fatalf("self-approval accepted: %v", err)
	}
	if err := s.Approve("t", "a1", "operator"); err != nil {
		t.Fatal(err)
	}
	if st, _ := s.State("t", "a1"); st != StateApproved {
		t.Fatalf("expected approved, got %s", st)
	}
}

func TestActivateOnlyFromApproved(t *testing.T) {
	s, _ := NewStore(10, time.Hour)
	if err := s.Submit(validApplicant("t", "a1", "x@y.z")); err != nil {
		t.Fatal(err)
	}
	if err := s.Activate("t", "a1"); !errors.Is(err, ErrBadTransition) {
		t.Fatalf("activate before approve accepted: %v", err)
	}
	if err := s.MarkContactVerified("t", "a1"); err != nil {
		t.Fatal(err)
	}
	if err := s.Approve("t", "a1", "op"); err != nil {
		t.Fatal(err)
	}
	if err := s.Activate("t", "a1"); err != nil {
		t.Fatal(err)
	}
	if st, _ := s.State("t", "a1"); st != StateActivated {
		t.Fatalf("expected activated, got %s", st)
	}
}

func TestRejectRequiresReasonAndDualControl(t *testing.T) {
	s, _ := NewStore(10, time.Hour)
	if err := s.Submit(validApplicant("t", "a1", "x@y.z")); err != nil {
		t.Fatal(err)
	}
	if err := s.Reject("t", "a1", "op", ""); !errors.Is(err, ErrReasonRequired) {
		t.Fatalf("empty reason accepted: %v", err)
	}
	if err := s.Reject("t", "a1", "self", "fake"); !errors.Is(err, ErrSelfApproval) {
		t.Fatalf("self-reject accepted: %v", err)
	}
	if err := s.Reject("t", "a1", "op", "unverifiable"); err != nil {
		t.Fatal(err)
	}
	if st, _ := s.State("t", "a1"); st != StateRejected {
		t.Fatalf("expected rejected, got %s", st)
	}
}

func TestVelocityExceeded(t *testing.T) {
	s, _ := NewStore(2, time.Hour)
	if err := s.Submit(validApplicant("t", "a1", "x1@y.z")); err != nil {
		t.Fatal(err)
	}
	if err := s.Submit(validApplicant("t", "a2", "x2@y.z")); err != nil {
		t.Fatal(err)
	}
	if err := s.Submit(validApplicant("t", "a3", "x3@y.z")); !errors.Is(err, ErrVelocityExceeded) {
		t.Fatalf("velocity exceeded accepted: %v", err)
	}
}

func TestInvalidInput(t *testing.T) {
	s, _ := NewStore(10, time.Hour)
	a := validApplicant("t", "a1", "not-an-email")
	if err := s.Submit(a); !errors.Is(err, ErrInvalidApplicant) {
		t.Fatalf("bad email accepted: %v", err)
	}
	b := validApplicant("t", "a1", "x@y.z")
	b.Kind = "admin"
	if err := s.Submit(b); !errors.Is(err, ErrInvalidApplicant) {
		t.Fatalf("bad kind accepted: %v", err)
	}
}

func TestTenantIsolation(t *testing.T) {
	s, _ := NewStore(10, time.Hour)
	if err := s.Submit(validApplicant("t1", "a1", "x@y.z")); err != nil {
		t.Fatal(err)
	}
	// Same email in another tenant is fine.
	if err := s.Submit(validApplicant("t2", "a1", "x@y.z")); err != nil {
		t.Fatalf("tenant isolation failed: %v", err)
	}
}
````


## 6. Configuration surface

Sin variables ni secretos (límite de velocidad y ventana se pasan por `NewStore`).

## 7. Dependency bill

| Package/image/tool | Pin exacto | Uso | Licencia | Runtime/build | Fuente oficial |
|---|---|---|---|---|---|
| Go standard library | go 1.26.7 (toolchain) | alta anti-fake | BSD-3-Clause | runtime | https://go.dev |

## 8. Apply order

1. Componer el backend (mismo módulo).
2. Colocar los dos archivos bajo `internal/applicant/`.
3. Verificar con `go test ./... -count=1` y `go vet ./...`.
4. Rollback: eliminar `internal/applicant/`.

## 9. Verification

- `go test ./internal/applicant/ -count=1`: 8/8 PASS.
- `go test ./... -count=1` (43 paquetes): PASS.
- `go vet ./...`: exit 0.

## 10. Reconstruction evidence

Véase `reconstruction_evidence/GO_APPLICANT_ONBOARDING_CORE_2026-09-02_V223.md`.
