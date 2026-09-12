# Go Data Residency Core

## 1. Metadata

```yaml
pack_id: "GO-DATA-RESIDENCY-CORE"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa la aplicación de política de residencia de datos: tenant-scoped, deny-by-default y fail-closed; cada categoría declara sus regiones permitidas y todo lo demás se rechaza."
stacks: ["Go 1.26.7"]
compatible_with: ["GO-ENTERPRISE-BACKEND-CORE 0.1.x"]
incompatible_with: ["categoría sin política", "región fuera del conjunto permitido", "política vacía"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://go.dev", "https://gdpr-info.eu"]
verified_at: "2026-09-02"
```

## 2. Applicability

Use para decidir dónde puede residir una categoría de datos (PII, pagos, etc.) antes de escribir o replicar. Rechace para categoría sin política o región fuera del conjunto permitido.

## 3. Architecture contract

- **Ownership**: `internal/residency` declara y aplica la política; la colocación física real es del proyecto.
- **Invariantes**: (1) deny-by-default: categoría sin política → rechazo. (2) región fuera del conjunto permitido → rechazo. (3) la política exige al menos una región válida. (4) tenant-scoped.
- **Data flow**: `SetPolicy(Policy)` → `Check(category, region)`.
- **Failure modes**: política inválida, sin política, región no permitida → error (fail-closed).
- **Seguridad/privacidad**: tenant-scoped.
- **Performance budget**: O(1) por `Check`.

## 4. Exact file manifest

```text
CREATE internal/residency/residency.go
CREATE internal/residency/residency_test.go
```

## 5. Materialization blocks

### FILE: `internal/residency/residency.go`
```yaml
block_id: "GO-DATA-RESIDENCY-CORE:internal/residency/residency.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "1e59c31a34c8ea8af06d0b27fa86488c0ebf30981b140b33dfcc489a25b572ae"
variables: []
secrets_allowed: false
```
````go
// Package residency provides data-residency policy enforcement: tenant-scoped,
// deny-by-default, fail-closed. Governed by GDPR Chapter V (transfers) and
// regional data-sovereignty requirements. AUTHORED.
package residency

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"
)

var (
	categoryRe = regexp.MustCompile(`^[a-z][a-z0-9_-]{0,63}$`)
	regionRe   = regexp.MustCompile(`^[a-z0-9-]{2,32}$`)

	ErrInvalidPolicy = errors.New("residency: invalid policy")
	ErrNoPolicy      = errors.New("residency: no policy for category (deny by default)")
	ErrDenied        = errors.New("residency: region not allowed for category")
)

// Policy declares where a data category may reside.
type Policy struct {
	TenantID       string
	Category       string
	AllowedRegions []string
}

// Registry is the tenant-scoped residency policy registry.
type Registry struct {
	mu       sync.Mutex
	policies map[string]Policy
}

// NewRegistry returns an empty registry.
func NewRegistry() *Registry {
	return &Registry{policies: make(map[string]Policy)}
}

func key(tenant, category string) string { return tenant + "\x00" + category }

// SetPolicy declares (or replaces) the allowed regions for a category.
// At least one region is required; duplicates are removed.
func (r *Registry) SetPolicy(p Policy) error {
	if strings.TrimSpace(p.TenantID) == "" || !categoryRe.MatchString(p.Category) {
		return ErrInvalidPolicy
	}
	seen := make(map[string]bool)
	regions := make([]string, 0, len(p.AllowedRegions))
	for _, rg := range p.AllowedRegions {
		rg = strings.ToLower(strings.TrimSpace(rg))
		if !regionRe.MatchString(rg) {
			return ErrInvalidPolicy
		}
		if !seen[rg] {
			seen[rg] = true
			regions = append(regions, rg)
		}
	}
	if len(regions) == 0 {
		return ErrInvalidPolicy
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	p.AllowedRegions = regions
	r.policies[key(p.TenantID, p.Category)] = p
	return nil
}

// Check fails closed unless the category has a policy that allows the region.
func (r *Registry) Check(tenant, category, region string) error {
	region = strings.ToLower(strings.TrimSpace(region))
	if strings.TrimSpace(tenant) == "" || !categoryRe.MatchString(category) || !regionRe.MatchString(region) {
		return ErrInvalidPolicy
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	p, ok := r.policies[key(tenant, category)]
	if !ok {
		return ErrNoPolicy
	}
	for _, rg := range p.AllowedRegions {
		if rg == region {
			return nil
		}
	}
	return ErrDenied
}

// Regions returns the allowed regions for a category (nil if no policy).
func (r *Registry) Regions(tenant, category string) []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	p, ok := r.policies[key(tenant, category)]
	if !ok {
		return nil
	}
	out := make([]string, len(p.AllowedRegions))
	copy(out, p.AllowedRegions)
	return out
}

// Categories returns the declared categories for a tenant.
func (r *Registry) Categories(tenant string) []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	var out []string
	for _, p := range r.policies {
		if p.TenantID == tenant {
			out = append(out, p.Category)
		}
	}
	return out
}

// String aids debugging.
func (r *Registry) String() string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return fmt.Sprintf("residency(%d)", len(r.policies))
}
````

### FILE: `internal/residency/residency_test.go`
```yaml
block_id: "GO-DATA-RESIDENCY-CORE:internal/residency/residency_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "4f962bba4b9a4560262442547b4e47034a1292936899c14ede9590c72e8a0e89"
variables: []
secrets_allowed: false
```
````go
package residency

import (
	"errors"
	"testing"
)

func TestAllowedRegion(t *testing.T) {
	r := NewRegistry()
	if err := r.SetPolicy(Policy{TenantID: "t", Category: "pii", AllowedRegions: []string{"eu", "sa-east-1"}}); err != nil {
		t.Fatal(err)
	}
	if err := r.Check("t", "pii", "eu"); err != nil {
		t.Fatalf("allowed region rejected: %v", err)
	}
}

func TestDisallowedRegion(t *testing.T) {
	r := NewRegistry()
	if err := r.SetPolicy(Policy{TenantID: "t", Category: "pii", AllowedRegions: []string{"eu"}}); err != nil {
		t.Fatal(err)
	}
	if err := r.Check("t", "pii", "us-east-1"); !errors.Is(err, ErrDenied) {
		t.Fatalf("expected deny, got %v", err)
	}
}

func TestDenyByDefault(t *testing.T) {
	r := NewRegistry()
	if err := r.Check("t", "pii", "eu"); !errors.Is(err, ErrNoPolicy) {
		t.Fatalf("expected no-policy (deny by default), got %v", err)
	}
}

func TestInvalidPolicy(t *testing.T) {
	r := NewRegistry()
	if err := r.SetPolicy(Policy{TenantID: "", Category: "pii", AllowedRegions: []string{"eu"}}); !errors.Is(err, ErrInvalidPolicy) {
		t.Fatalf("empty tenant accepted: %v", err)
	}
	if err := r.SetPolicy(Policy{TenantID: "t", Category: "pii", AllowedRegions: []string{}}); !errors.Is(err, ErrInvalidPolicy) {
		t.Fatalf("empty regions accepted: %v", err)
	}
	if err := r.SetPolicy(Policy{TenantID: "t", Category: "pii", AllowedRegions: []string{"BAD REGION!"}}); !errors.Is(err, ErrInvalidPolicy) {
		t.Fatalf("bad region accepted: %v", err)
	}
}

func TestPolicyReplaceAndDedup(t *testing.T) {
	r := NewRegistry()
	if err := r.SetPolicy(Policy{TenantID: "t", Category: "pii", AllowedRegions: []string{"eu", "eu", "us-east-1"}}); err != nil {
		t.Fatal(err)
	}
	if got := len(r.Regions("t", "pii")); got != 2 {
		t.Fatalf("expected 2 deduped regions, got %d", got)
	}
	// Replace policy: us now denied.
	if err := r.SetPolicy(Policy{TenantID: "t", Category: "pii", AllowedRegions: []string{"eu"}}); err != nil {
		t.Fatal(err)
	}
	if err := r.Check("t", "pii", "us-east-1"); !errors.Is(err, ErrDenied) {
		t.Fatalf("expected deny after replace, got %v", err)
	}
}

func TestTenantIsolation(t *testing.T) {
	r := NewRegistry()
	if err := r.SetPolicy(Policy{TenantID: "t1", Category: "pii", AllowedRegions: []string{"eu"}}); err != nil {
		t.Fatal(err)
	}
	if err := r.Check("t2", "pii", "eu"); !errors.Is(err, ErrNoPolicy) {
		t.Fatalf("tenant leak: %v", err)
	}
}
````


## 6. Configuration surface

Sin variables ni secretos.

## 7. Dependency bill

| Package/image/tool | Pin exacto | Uso | Licencia | Runtime/build | Fuente oficial |
|---|---|---|---|---|---|
| Go standard library | go 1.26.7 (toolchain) | residencia de datos | BSD-3-Clause | runtime | https://go.dev |

## 8. Apply order

1. Componer el backend (mismo módulo).
2. Colocar los dos archivos bajo `internal/residency/`.
3. Verificar con `go test ./... -count=1` y `go vet ./...`.
4. Rollback: eliminar `internal/residency/`.

## 9. Verification

- `go test ./internal/residency/ -count=1`: 6/6 PASS (región permitida, denegada, deny-by-default, política inválida, reemplazo/dedup, aislamiento de tenant).
- `go test ./... -count=1` (32 paquetes): PASS.
- `go vet ./...`: exit 0.

## 10. Reconstruction evidence

Véase `reconstruction_evidence/GO_DATA_RESIDENCY_CORE_2026-09-02_V211.md`.
