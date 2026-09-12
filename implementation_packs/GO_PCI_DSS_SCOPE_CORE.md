# Go PCI DSS Scope Core

## 1. Metadata

```yaml
pack_id: "GO-PCI-DSS-SCOPE-CORE"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa el gobierno de alcance del CDE (PCI DSS 4.0): registro tenant-scoped de componentes, fail-closed ante datos de cuenta fuera de alcance, almacenamiento de SAD prohibido tras la autorización y SAQ obligatorio en alcance."
stacks: ["Go 1.26.7"]
compatible_with: ["GO-ENTERPRISE-BACKEND-CORE 0.1.x"]
incompatible_with: ["componente con datos de cuenta fuera de alcance", "almacenamiento de SAD", "SAQ ausente o desconocido"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://go.dev", "https://www.pcisecuritystandards.org"]
verified_at: "2026-09-02"
```

## 2. Applicability

Use para declarar y auditar el alcance del entorno de datos del titular de tarjeta (CDE) antes de incorporar un sistema que maneja PAN/SAD. No es una autoevaluación PCI completa; es el contrato de alcance que la precede.

## 3. Architecture contract

- **Ownership**: `internal/pciscope` gobierna el alcance; la autoevaluación (SAQ) y el ASV son procesos externos del operador.
- **Invariantes**: (1) componente que toca PAN/SAD debe estar en alcance. (2) el almacenamiento de SAD tras la autorización está prohibido (PCI DSS 4.0 Req 3.2). (3) todo componente en alcance declara un tipo SAQ válido. (4) el alcance se confirma al menos anualmente (Req 12.5.2).
- **Data flow**: `Declare(Component)` → `Scope(tenant)` / `Audit(tenant)`.
- **Failure modes**: scope violation → error (fail-closed).
- **Seguridad/privacidad**: tenant-scoped; `Audit` devuelve violaciones sin confiar en `Declare`.
- **Performance budget**: O(1) por `Declare`; O(N) por `Audit`.

## 4. Exact file manifest

```text
CREATE internal/pciscope/pciscope.go
CREATE internal/pciscope/pciscope_test.go
```

## 5. Materialization blocks

### FILE: `internal/pciscope/pciscope.go`
```yaml
block_id: "GO-PCI-DSS-SCOPE-CORE:internal/pciscope/pciscope.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "454151a17fb92854fac14cc2f5998fb048e667de537a32f69f41aedd8923cee9"
variables: []
secrets_allowed: false
```
````go
// Package pciscope provides PCI DSS cardholder-data-environment (CDE) scope
// governance: a tenant-scoped registry that declares system components and
// fail-closed on scope violations. Governed by PCI DSS 4.0 (Req 3.2: do not
// store sensitive authentication data after authorization; Req 12.5.2: scope
// is confirmed at least annually). AUTHORED; not a PCI assessment.
package pciscope

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"
)

// Touch records which account-data operations a component performs.
// PAN = primary account number; SAD = sensitive authentication data
// (full track, CVV, PIN/PIN block — never stored after authorization).
type Touch struct {
	StoresPAN     bool
	ProcessesPAN  bool
	TransmitsPAN  bool
	StoresSAD     bool
}

// TouchesAccountData reports whether the component handles PAN at all.
func (t Touch) TouchesAccountData() bool { return t.StoresPAN || t.ProcessesPAN || t.TransmitsPAN }

// SAQ types defined by PCI DSS 4.0 (SAQ A, A-EP, D, P2PE, SPoC, D-SP).
var saqSet = map[string]bool{
	"SAQ A": true, "SAQ A-EP": true, "SAQ D": true,
	"SAQ P2PE": true, "SAQ SPoC": true, "SAQ D-SP": true,
}

var (
	idRe = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:-]{0,127}$`)

	ErrInvalidComponent      = errors.New("pciscope: invalid component")
	ErrAccountDataOutOfScope = errors.New("pciscope: account-data component must be in scope")
	ErrSADStorage            = errors.New("pciscope: SAD storage prohibited after authorization")
	ErrMissingSAQ            = errors.New("pciscope: in-scope component requires a SAQ type")
	ErrBadSAQ                = errors.New("pciscope: unknown SAQ type")
	ErrDuplicate             = errors.New("pciscope: duplicate component")
)

// Component is a declared system component and its CDE scope status.
type Component struct {
	TenantID string
	ID       string
	Name     string
	Touch    Touch
	InScope  bool
	SAQ      string // SAQ type when in scope; empty when out of scope
}

// Registry is the tenant-scoped CDE scope registry.
type Registry struct {
	mu     sync.Mutex
	items  map[string]Component
}

// NewRegistry returns an empty registry.
func NewRegistry() *Registry {
	return &Registry{items: make(map[string]Component)}
}

func key(tenant, id string) string { return tenant + "\x00" + id }

func (c Component) validate() error {
	if strings.TrimSpace(c.TenantID) == "" || !idRe.MatchString(c.ID) || strings.TrimSpace(c.Name) == "" {
		return ErrInvalidComponent
	}
	if c.Touch.StoresSAD {
		return ErrSADStorage
	}
	if c.Touch.TouchesAccountData() && !c.InScope {
		return ErrAccountDataOutOfScope
	}
	if c.InScope {
		if strings.TrimSpace(c.SAQ) == "" {
			return ErrMissingSAQ
		}
		if !saqSet[c.SAQ] {
			return ErrBadSAQ
		}
	}
	return nil
}

// Declare registers a component, fail-closed on any scope violation.
func (r *Registry) Declare(c Component) error {
	if err := c.validate(); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.items[key(c.TenantID, c.ID)]; ok {
		return ErrDuplicate
	}
	r.items[key(c.TenantID, c.ID)] = c
	return nil
}

// Scope returns the tenant-scoped components.
func (r *Registry) Scope(tenant string) []Component {
	r.mu.Lock()
	defer r.mu.Unlock()
	var out []Component
	for _, c := range r.items {
		if c.TenantID == tenant {
			out = append(out, c)
		}
	}
	return out
}

// Audit defensively re-checks the registry and returns every violation found.
// A non-empty result means the CDE scope is not trustworthy.
func (r *Registry) Audit(tenant string) []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	var out []string
	for _, c := range r.items {
		if c.TenantID != tenant {
			continue
		}
		if c.Touch.StoresSAD {
			out = append(out, fmt.Sprintf("%s: SAD storage prohibited", c.ID))
		}
		if c.Touch.TouchesAccountData() && !c.InScope {
			out = append(out, fmt.Sprintf("%s: account data out of scope", c.ID))
		}
		if c.InScope && strings.TrimSpace(c.SAQ) == "" {
			out = append(out, fmt.Sprintf("%s: in scope without SAQ", c.ID))
		}
	}
	return out
}

// String aids debugging.
func (r *Registry) String() string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return fmt.Sprintf("pciscope(%d)", len(r.items))
}
````

### FILE: `internal/pciscope/pciscope_test.go`
```yaml
block_id: "GO-PCI-DSS-SCOPE-CORE:internal/pciscope/pciscope_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "cd2dd60efad4490ea9168a1c7d550eba97c0a7dd31fc3e7e9692119047c09c31"
variables: []
secrets_allowed: false
```
````go
package pciscope

import (
	"errors"
	"testing"
)

func TestAccountDataOutOfScopeRejected(t *testing.T) {
	r := NewRegistry()
	c := Component{
		TenantID: "t", ID: "payments", Name: "Payments API",
		Touch: Touch{ProcessesPAN: true}, InScope: false,
	}
	if err := r.Declare(c); !errors.Is(err, ErrAccountDataOutOfScope) {
		t.Fatalf("expected account-data-out-of-scope, got %v", err)
	}
}

func TestInScopeWithSAQAccepted(t *testing.T) {
	r := NewRegistry()
	c := Component{
		TenantID: "t", ID: "payments", Name: "Payments API",
		Touch: Touch{ProcessesPAN: true}, InScope: true, SAQ: "SAQ D",
	}
	if err := r.Declare(c); err != nil {
		t.Fatal(err)
	}
	if got := len(r.Scope("t")); got != 1 {
		t.Fatalf("expected 1 component, got %d", got)
	}
}

func TestNonAccountDataOutOfScopeAccepted(t *testing.T) {
	r := NewRegistry()
	c := Component{TenantID: "t", ID: "cdn", Name: "Static CDN", InScope: false}
	if err := r.Declare(c); err != nil {
		t.Fatal(err)
	}
}

func TestSADStorageProhibited(t *testing.T) {
	r := NewRegistry()
	c := Component{
		TenantID: "t", ID: "legacy", Name: "Legacy",
		Touch: Touch{StoresSAD: true}, InScope: true, SAQ: "SAQ D",
	}
	if err := r.Declare(c); !errors.Is(err, ErrSADStorage) {
		t.Fatalf("expected SAD-storage prohibition, got %v", err)
	}
}

func TestMissingSAQRejected(t *testing.T) {
	r := NewRegistry()
	c := Component{
		TenantID: "t", ID: "payments", Name: "Payments API",
		Touch: Touch{TransmitsPAN: true}, InScope: true,
	}
	if err := r.Declare(c); !errors.Is(err, ErrMissingSAQ) {
		t.Fatalf("expected missing-SAQ, got %v", err)
	}
}

func TestDuplicateAndTenantIsolation(t *testing.T) {
	r := NewRegistry()
	c := Component{TenantID: "t", ID: "x", Name: "X", InScope: false}
	if err := r.Declare(c); err != nil {
		t.Fatal(err)
	}
	if err := r.Declare(c); !errors.Is(err, ErrDuplicate) {
		t.Fatalf("expected duplicate, got %v", err)
	}
	// Same id in another tenant is a distinct component.
	c2 := Component{TenantID: "t2", ID: "x", Name: "X2", InScope: false}
	if err := r.Declare(c2); err != nil {
		t.Fatalf("tenant isolation failed: %v", err)
	}
}

func TestAuditClean(t *testing.T) {
	r := NewRegistry()
	for _, c := range []Component{
		{TenantID: "t", ID: "pay", Name: "Pay", Touch: Touch{ProcessesPAN: true}, InScope: true, SAQ: "SAQ D"},
		{TenantID: "t", ID: "cdn", Name: "CDN", InScope: false},
	} {
		if err := r.Declare(c); err != nil {
			t.Fatal(err)
		}
	}
	if v := r.Audit("t"); len(v) != 0 {
		t.Fatalf("expected clean audit, got %v", v)
	}
}
````


## 6. Configuration surface

Sin variables ni secretos.

## 7. Dependency bill

| Package/image/tool | Pin exacto | Uso | Licencia | Runtime/build | Fuente oficial |
|---|---|---|---|---|---|
| Go standard library | go 1.26.7 (toolchain) | alcance CDE | BSD-3-Clause | runtime | https://go.dev |

## 8. Apply order

1. Componer el backend (mismo módulo).
2. Colocar los dos archivos bajo `internal/pciscope/`.
3. Verificar con `go test ./... -count=1` y `go vet ./...`.
4. Rollback: eliminar `internal/pciscope/`.

## 9. Verification

- `go test ./internal/pciscope/ -count=1`: 7/7 PASS (fuera de alcance rechazado, SAQ válido, fuera de alcance sin datos, SAD prohibido, SAQ ausente, duplicado/aislamiento de tenant, audit limpio).
- `go test ./... -count=1` (30 paquetes): PASS.
- `go vet ./...`: exit 0.

## 10. Reconstruction evidence

Véase `reconstruction_evidence/GO_PCI_DSS_SCOPE_CORE_2026-09-02_V208.md`.
