# Go GDPR Consent & Erasure Core

## 1. Metadata

```yaml
pack_id: "GO-GDPR-CONSENT-ERASURE-CORE"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa el gobierno GDPR de consentimiento (Art. 6/7) y derecho de supresión (Art. 17): consentimiento por finalidad y versión de política revocable, y máquina de estados de borrado fail-closed con retención legal."
stacks: ["Go 1.26.7"]
compatible_with: ["GO-ENTERPRISE-BACKEND-CORE 0.1.x"]
incompatible_with: ["consentimiento sin finalidad o versión de política", "rechazo de borrado sin base legal", "borrado con retención legal activa"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://go.dev", "https://gdpr-info.eu"]
verified_at: "2026-09-02"
```

## 2. Applicability

Use para registrar consentimiento y ejecutar el derecho de supresión en un tenant. Rechace para consentimiento sin finalidad/versión, rechazo sin base legal o borrado con retención legal activa.

## 3. Architecture contract

- **Ownership**: `internal/gdpr` gobierna consentimiento y borrado; la ejecución real del borrado en los almacenes de dominio es del proyecto.
- **Invariantes**: (1) consentimiento específico e informado (finalidad + versión de política). (2) retirada tan fácil como otorgarlo (Art. 7(3)). (3) borrado sólo si no hay retención legal activa (Art. 17(3)). (4) rechazo exige base legal.
- **Data flow**: `Record/Withdraw` (consent) y `RequestErasure→Begin→Complete|Reject` (erasure).
- **Failure modes**: entrada inválida, transición inválida, retención activa, base legal ausente → error.
- **Seguridad/privacidad**: tenant-scoped.
- **Performance budget**: O(1) por operación; O(H) por hold.

## 4. Exact file manifest

```text
CREATE internal/gdpr/gdpr.go
CREATE internal/gdpr/gdpr_test.go
```

## 5. Materialization blocks

### FILE: `internal/gdpr/gdpr.go`
```yaml
block_id: "GO-GDPR-CONSENT-ERASURE-CORE:internal/gdpr/gdpr.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "295da7e159b6de92dfaf25498802fdfcc5a0366c1ee47a8ecd9e431a373c6da8"
variables: []
secrets_allowed: false
```
````go
// Package gdpr provides GDPR consent and right-to-erasure governance:
// tenant-scoped, fail-closed. Governed by GDPR Art. 6/7 (consent must be
// specific, informed, freely given and withdrawable) and Art. 17 (erasure
// with retention/legal-hold exceptions). AUTHORED.
package gdpr

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"
)

var (
	idRe = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:-]{0,127}$`)

	ErrInvalidInput     = errors.New("gdpr: invalid input")
	ErrNotFound         = errors.New("gdpr: not found")
	ErrDuplicate        = errors.New("gdpr: duplicate request")
	ErrBadTransition    = errors.New("gdpr: bad state transition")
	ErrRetentionHold    = errors.New("gdpr: erasure blocked by retention hold")
	ErrNoLegalBasis     = errors.New("gdpr: rejection requires a legal basis")
)

// ---- Consent (Art. 6/7) ----

// Consent is a purpose-specific consent grant, bound to a policy version.
type Consent struct {
	TenantID      string
	SubjectID     string
	Purpose       string
	PolicyVersion string
	Granted       bool
	At            time.Time
}

// ConsentStore holds per-subject/purpose consent records.
type ConsentStore struct {
	mu    sync.Mutex
	items map[string]Consent
}

// NewConsentStore returns an empty consent store.
func NewConsentStore() *ConsentStore {
	return &ConsentStore{items: make(map[string]Consent)}
}

func consentKey(tenant, subject, purpose string) string {
	return tenant + "\x00" + subject + "\x00" + purpose
}

func validateConsent(tenant, subject, purpose, policyVersion string) error {
	if strings.TrimSpace(tenant) == "" || !idRe.MatchString(subject) ||
		strings.TrimSpace(purpose) == "" || strings.TrimSpace(policyVersion) == "" {
		return ErrInvalidInput
	}
	return nil
}

// Record grants consent for a purpose, bound to a policy version (informed).
// Re-recording an existing consent re-grants it with the new version.
func (s *ConsentStore) Record(tenant, subject, purpose, policyVersion string) error {
	if err := validateConsent(tenant, subject, purpose, policyVersion); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	k := consentKey(tenant, subject, purpose)
	at := time.Now().UTC()
	if c, ok := s.items[k]; ok {
		at = c.At // preserve the original grant time
	}
	s.items[k] = Consent{
		TenantID: tenant, SubjectID: subject, Purpose: purpose,
		PolicyVersion: policyVersion, Granted: true, At: at,
	}
	return nil
}

// Withdraw revokes consent; withdrawal is as easy as giving it (Art. 7(3)).
func (s *ConsentStore) Withdraw(tenant, subject, purpose string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	k := consentKey(tenant, subject, purpose)
	c, ok := s.items[k]
	if !ok {
		return ErrNotFound
	}
	c.Granted = false
	c.At = time.Now().UTC()
	s.items[k] = c
	return nil
}

// IsGranted reports whether consent is currently granted.
func (s *ConsentStore) IsGranted(tenant, subject, purpose string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.items[consentKey(tenant, subject, purpose)].Granted
}

// ---- Erasure (Art. 17) ----

// State is the erasure-request lifecycle.
type State string

const (
	StateRequested  State = "requested"
	StateInProgress State = "in_progress"
	StateCompleted  State = "completed"
	StateRejected   State = "rejected"
)

// ErasureRequest is a single right-to-erasure request.
type ErasureRequest struct {
	TenantID  string
	SubjectID string
	ID        string
	State     State
	LegalBasis string // set on rejection (e.g. tax retention obligation)
}

// ErasureStore holds requests and retention holds.
type ErasureStore struct {
	mu       sync.Mutex
	requests map[string]ErasureRequest
	holds    map[string][]string // tenant+subject -> reasons
}

// NewErasureStore returns an empty erasure store.
func NewErasureStore() *ErasureStore {
	return &ErasureStore{requests: make(map[string]ErasureRequest), holds: make(map[string][]string)}
}

func holdKey(tenant, subject string) string { return tenant + "\x00" + subject }

// RequestErasure registers a request (Art. 17). Idempotent per request id.
func (s *ErasureStore) RequestErasure(tenant, subject, id string) error {
	if strings.TrimSpace(tenant) == "" || !idRe.MatchString(subject) || !idRe.MatchString(id) {
		return ErrInvalidInput
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.requests[id]; ok {
		return ErrDuplicate
	}
	s.requests[id] = ErasureRequest{TenantID: tenant, SubjectID: subject, ID: id, State: StateRequested}
	return nil
}

// AddRetentionHold blocks erasure completion for a subject (e.g. legal/tax
// retention obligation, Art. 17(3)).
func (s *ErasureStore) AddRetentionHold(tenant, subject, reason string) error {
	if strings.TrimSpace(tenant) == "" || !idRe.MatchString(subject) || strings.TrimSpace(reason) == "" {
		return ErrInvalidInput
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	k := holdKey(tenant, subject)
	s.holds[k] = append(s.holds[k], reason)
	return nil
}

// RemoveRetentionHold releases one hold reason.
func (s *ErasureStore) RemoveRetentionHold(tenant, subject, reason string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	k := holdKey(tenant, subject)
	list := s.holds[k]
	for i, r := range list {
		if r == reason {
			s.holds[k] = append(list[:i], list[i+1:]...)
			return nil
		}
	}
	return ErrNotFound
}

// Begin marks requested → in_progress.
func (s *ErasureStore) Begin(id string) error {
	return s.transition(id, StateRequested, StateInProgress, "")
}

// Complete marks in_progress → completed, fail-closed if a hold remains.
func (s *ErasureStore) Complete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	r, ok := s.requests[id]
	if !ok {
		return ErrNotFound
	}
	if r.State != StateInProgress {
		return ErrBadTransition
	}
	if len(s.holds[holdKey(r.TenantID, r.SubjectID)]) > 0 {
		return ErrRetentionHold
	}
	r.State = StateCompleted
	s.requests[id] = r
	return nil
}

// Reject marks requested/in_progress → rejected with a legal basis.
func (s *ErasureStore) Reject(id, legalBasis string) error {
	if strings.TrimSpace(legalBasis) == "" {
		return ErrNoLegalBasis
	}
	return s.transitionAny(id, legalBasis)
}

func (s *ErasureStore) transition(id string, from, to State, basis string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	r, ok := s.requests[id]
	if !ok {
		return ErrNotFound
	}
	if r.State != from {
		return ErrBadTransition
	}
	r.State = to
	r.LegalBasis = basis
	s.requests[id] = r
	return nil
}

func (s *ErasureStore) transitionAny(id, basis string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	r, ok := s.requests[id]
	if !ok {
		return ErrNotFound
	}
	if r.State == StateCompleted || r.State == StateRejected {
		return ErrBadTransition
	}
	r.State = StateRejected
	r.LegalBasis = basis
	s.requests[id] = r
	return nil
}

// State returns the request state.
func (s *ErasureStore) State(id string) (State, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	r, ok := s.requests[id]
	return r.State, ok
}

// String aids debugging.
func (s *ErasureStore) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return fmt.Sprintf("gdpr.erasure(%d)", len(s.requests))
}
````

### FILE: `internal/gdpr/gdpr_test.go`
```yaml
block_id: "GO-GDPR-CONSENT-ERASURE-CORE:internal/gdpr/gdpr_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "722b7ddd9debd48757793ea10129c1e88cee7b5ca851611baf67d269bcdbbec5"
variables: []
secrets_allowed: false
```
````go
package gdpr

import (
	"errors"
	"testing"
)

func TestConsentRecordWithdrawRegrant(t *testing.T) {
	s := NewConsentStore()
	if err := s.Record("t", "sub1", "marketing", "v2"); err != nil {
		t.Fatal(err)
	}
	if !s.IsGranted("t", "sub1", "marketing") {
		t.Fatal("consent should be granted")
	}
	if err := s.Withdraw("t", "sub1", "marketing"); err != nil {
		t.Fatal(err)
	}
	if s.IsGranted("t", "sub1", "marketing") {
		t.Fatal("consent should be withdrawn")
	}
	// Re-record re-grants with a new policy version.
	if err := s.Record("t", "sub1", "marketing", "v3"); err != nil {
		t.Fatal(err)
	}
	if !s.IsGranted("t", "sub1", "marketing") {
		t.Fatal("consent should be re-granted")
	}
}

func TestConsentInvalid(t *testing.T) {
	s := NewConsentStore()
	if err := s.Record("t", "sub1", "", "v1"); !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("empty purpose accepted: %v", err)
	}
	if err := s.Record("t", "sub1", "marketing", ""); !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("empty policy version accepted: %v", err)
	}
	if err := s.Record("", "sub1", "marketing", "v1"); !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("empty tenant accepted: %v", err)
	}
}

func TestErasureHappyPath(t *testing.T) {
	s := NewErasureStore()
	if err := s.RequestErasure("t", "sub1", "r1"); err != nil {
		t.Fatal(err)
	}
	if err := s.Begin("r1"); err != nil {
		t.Fatal(err)
	}
	if err := s.Complete("r1"); err != nil {
		t.Fatal(err)
	}
	if st, _ := s.State("r1"); st != StateCompleted {
		t.Fatalf("expected completed, got %s", st)
	}
}

func TestErasureBlockedByRetentionHold(t *testing.T) {
	s := NewErasureStore()
	if err := s.RequestErasure("t", "sub1", "r1"); err != nil {
		t.Fatal(err)
	}
	if err := s.Begin("r1"); err != nil {
		t.Fatal(err)
	}
	if err := s.AddRetentionHold("t", "sub1", "tax-retention-10y"); err != nil {
		t.Fatal(err)
	}
	if err := s.Complete("r1"); !errors.Is(err, ErrRetentionHold) {
		t.Fatalf("expected retention hold, got %v", err)
	}
	if err := s.RemoveRetentionHold("t", "sub1", "tax-retention-10y"); err != nil {
		t.Fatal(err)
	}
	if err := s.Complete("r1"); err != nil {
		t.Fatalf("completion after hold release failed: %v", err)
	}
}

func TestRejectRequiresLegalBasis(t *testing.T) {
	s := NewErasureStore()
	if err := s.RequestErasure("t", "sub1", "r1"); err != nil {
		t.Fatal(err)
	}
	if err := s.Reject("r1", ""); !errors.Is(err, ErrNoLegalBasis) {
		t.Fatalf("expected no-legal-basis, got %v", err)
	}
	if err := s.Reject("r1", "legal-obligation-Art17(3)"); err != nil {
		t.Fatal(err)
	}
	if st, _ := s.State("r1"); st != StateRejected {
		t.Fatalf("expected rejected, got %s", st)
	}
}

func TestBadTransitionAndDuplicate(t *testing.T) {
	s := NewErasureStore()
	if err := s.RequestErasure("t", "sub1", "r1"); err != nil {
		t.Fatal(err)
	}
	if err := s.Complete("r1"); !errors.Is(err, ErrBadTransition) {
		t.Fatalf("complete-before-begin accepted: %v", err)
	}
	if err := s.RequestErasure("t", "sub1", "r1"); !errors.Is(err, ErrDuplicate) {
		t.Fatalf("duplicate request accepted: %v", err)
	}
}

func TestTenantIsolation(t *testing.T) {
	cs := NewConsentStore()
	if err := cs.Record("t1", "sub1", "marketing", "v1"); err != nil {
		t.Fatal(err)
	}
	if cs.IsGranted("t2", "sub1", "marketing") {
		t.Fatal("consent leaked across tenants")
	}
}
````


## 6. Configuration surface

Sin variables ni secretos.

## 7. Dependency bill

| Package/image/tool | Pin exacto | Uso | Licencia | Runtime/build | Fuente oficial |
|---|---|---|---|---|---|
| Go standard library | go 1.26.7 (toolchain) | consentimiento/borrado | BSD-3-Clause | runtime | https://go.dev |

## 8. Apply order

1. Componer el backend (mismo módulo).
2. Colocar los dos archivos bajo `internal/gdpr/`.
3. Verificar con `go test ./... -count=1` y `go vet ./...`.
4. Rollback: eliminar `internal/gdpr/`.

## 9. Verification

- `go test ./internal/gdpr/ -count=1`: 7/7 PASS (consentimiento revocable/re-otorgable, inválido, borrado happy path, retención bloquea, base legal exigida, transición inválida/duplicado, aislamiento de tenant).
- `go test ./... -count=1` (31 paquetes): PASS.
- `go vet ./...`: exit 0.

## 10. Reconstruction evidence

Véase `reconstruction_evidence/GO_GDPR_CONSENT_ERASURE_CORE_2026-09-02_V209.md`.
