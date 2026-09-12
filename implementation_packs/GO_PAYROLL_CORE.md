# Go Payroll Core

## 1. Metadata

```yaml
pack_id: "GO-PAYROLL-CORE"
pack_version: "0.1.3"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa el cálculo de nómina tenant-scoped config-driven: deducciones fijas o porcentuales inyectadas en runtime, net fail-closed (nunca negativo) e idempotente. Sin tablas impositivas hardcodeadas."
stacks: ["Go 1.26.8"]
compatible_with: ["módulo aislado Go 1.26.7; integración con owners empresariales no admitida"]
incompatible_with: ["net negativo", "deducción inválida", "run duplicado"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://go.dev"]
verified_at: "2026-09-10"
```

## 2. Applicability

Referencia AUTHORED aritmética en memoria; las deducciones dadas por el caller no prueban legalidad laboral/fiscal. No implementa liquidación, recibos, aprobación, pago ni ledger durable. Conserva porcentajes sobre gross y truncado individual; errores de deducción inválida prevalecen sobre net matemáticamente negativo.

## 3. Architecture contract

- **Ownership**: `internal/payroll` calcula; las reglas impositivas/legales son configuración del operador (no se inventan tablas).
- **Invariantes**: (1) net ≥ 0. (2) deducciones fijas ≥ 0 y porcentuales 0..10000 bps. (3) duplicados por tenant/run id rechazados con ErrDuplicate, sin replay exitoso. (4) net válido nunca supera gross por overflow; reglas inválidas conservan prioridad de error.
- **Data flow**: `Execute(Run)` → net.
- **Failure modes**: net negativo, deducción inválida, duplicado → error.
- **Seguridad/privacidad**: tenant-scoped.
- **Performance budget**: O(D) por run.

## 4. Exact file manifest

```text
CREATE internal/payroll/payroll.go
CREATE internal/payroll/payroll_test.go
```

## 5. Materialization blocks

### FILE: `internal/payroll/payroll.go`
```yaml
block_id: "GO-PAYROLL-CORE:internal/payroll/payroll.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "5a9b0f24dd5f18eb87f3285a715e0595a57bf8e42e866a644d6635f2469b6b9d"
variables: []
secrets_allowed: false
```
````go
// Package payroll provides a tenant-scoped payroll run: config-driven
// deductions (fixed amount or percentage of gross) applied to gross, with
// idempotency and fail-closed net (never negative). No tax tables are
// hardcoded — deduction rules are injected at runtime by jurisdiction.
// AUTHORED.
package payroll

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"
)

// Kind is the deduction calculation.
type Kind string

const (
	KindFixed   Kind = "fixed"
	KindPercent Kind = "percent"
)

var (
	idRe = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:-]{0,127}$`)

	ErrInvalidRun  = errors.New("payroll: invalid run")
	ErrNegativeNet = errors.New("payroll: net would be negative")
	ErrDuplicate   = errors.New("payroll: duplicate run")
	ErrNotFound    = errors.New("payroll: not found")
)

// Deduction is a single config-driven deduction.
type Deduction struct {
	Code             string
	Kind             Kind
	AmountMinorUnits int64 // fixed amount in minor units
	BasisPoints      int64 // percent in basis points (0..10000)
}

// Run is a single payroll computation.
type Run struct {
	TenantID        string
	ID              string
	EmployeeID      string
	GrossMinorUnits int64
	Deductions      []Deduction
	NetMinorUnits   int64
}

// Store holds executed runs.
type Store struct {
	mu   sync.Mutex
	runs map[string]Run
}

// NewStore returns an empty store.
func NewStore() *Store {
	return &Store{runs: make(map[string]Run)}
}

func key(tenant, id string) string { return tenant + "\x00" + id }

// Execute computes net = gross - sum(deductions), idempotent by run id,
// fail-closed on negative net or invalid deductions.
func (s *Store) Execute(run Run) (int64, error) {
	if strings.TrimSpace(run.TenantID) == "" || !idRe.MatchString(run.ID) || !idRe.MatchString(run.EmployeeID) {
		return 0, ErrInvalidRun
	}
	if run.GrossMinorUnits < 0 {
		return 0, ErrInvalidRun
	}
	net := run.GrossMinorUnits
	exceedsGross := false
	for _, d := range run.Deductions {
		var amount int64
		switch d.Kind {
		case KindFixed:
			if d.AmountMinorUnits < 0 {
				return 0, ErrInvalidRun
			}
			amount = d.AmountMinorUnits
		case KindPercent:
			if d.BasisPoints < 0 || d.BasisPoints > 10000 {
				return 0, ErrInvalidRun
			}
			// Each percentage still applies to gross and truncates independently.
			amount = (run.GrossMinorUnits/10000)*d.BasisPoints + (run.GrossMinorUnits%10000)*d.BasisPoints/10000
		default:
			return 0, ErrInvalidRun
		}
		// Keep validating later deductions so invalid-run error precedence is
		// unchanged. Nonnegative deductions cannot recover an exceeded gross.
		if exceedsGross || amount > net {
			exceedsGross = true
			continue
		}
		net -= amount
	}
	if exceedsGross {
		return 0, ErrNegativeNet
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	k := key(run.TenantID, run.ID)
	if _, ok := s.runs[k]; ok {
		return 0, ErrDuplicate
	}
	run.NetMinorUnits = net
	s.runs[k] = run
	return net, nil
}

// Net returns the executed net for a run.
func (s *Store) Net(tenant, id string) (int64, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	r, ok := s.runs[key(tenant, id)]
	return r.NetMinorUnits, ok
}

// String aids debugging.
func (s *Store) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return fmt.Sprintf("payroll(%d)", len(s.runs))
}
````

### FILE: `internal/payroll/payroll_test.go`
```yaml
block_id: "GO-PAYROLL-CORE:internal/payroll/payroll_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "0d44af51784430e279fb9ab796079387b9f728f5b2c555fed829baa7049f8457"
variables: []
secrets_allowed: false
```
````go
package payroll

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"sync"
	"sync/atomic"
	"testing"
)

func TestFixedDeductions(t *testing.T) {
	s := NewStore()
	run := Run{
		TenantID: "t", ID: "r1", EmployeeID: "e1", GrossMinorUnits: 100000,
		Deductions: []Deduction{
			{Code: "health", Kind: KindFixed, AmountMinorUnits: 10000},
			{Code: "pension", Kind: KindFixed, AmountMinorUnits: 5000},
		},
	}
	net, err := s.Execute(run)
	if err != nil {
		t.Fatal(err)
	}
	if net != 85000 {
		t.Fatalf("expected net 85000, got %d", net)
	}
}

func TestPercentDeduction(t *testing.T) {
	s := NewStore()
	run := Run{
		TenantID: "t", ID: "r1", EmployeeID: "e1", GrossMinorUnits: 100000,
		Deductions: []Deduction{{Code: "tax", Kind: KindPercent, BasisPoints: 1000}}, // 10%
	}
	net, err := s.Execute(run)
	if err != nil {
		t.Fatal(err)
	}
	if net != 90000 {
		t.Fatalf("expected net 90000, got %d", net)
	}
}

func TestNegativeNetRejected(t *testing.T) {
	s := NewStore()
	run := Run{
		TenantID: "t", ID: "r1", EmployeeID: "e1", GrossMinorUnits: 1000,
		Deductions: []Deduction{{Code: "x", Kind: KindFixed, AmountMinorUnits: 5000}},
	}
	if _, err := s.Execute(run); !errors.Is(err, ErrNegativeNet) {
		t.Fatalf("expected negative-net, got %v", err)
	}
}

func TestInvalidDeduction(t *testing.T) {
	s := NewStore()
	run := Run{
		TenantID: "t", ID: "r1", EmployeeID: "e1", GrossMinorUnits: 1000,
		Deductions: []Deduction{{Code: "x", Kind: KindPercent, BasisPoints: 20000}},
	}
	if _, err := s.Execute(run); !errors.Is(err, ErrInvalidRun) {
		t.Fatalf("expected invalid-run, got %v", err)
	}
}

func TestDuplicateAndTenantIsolation(t *testing.T) {
	s := NewStore()
	run := Run{TenantID: "t", ID: "r1", EmployeeID: "e1", GrossMinorUnits: 1000}
	if _, err := s.Execute(run); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Execute(run); !errors.Is(err, ErrDuplicate) {
		t.Fatalf("expected duplicate, got %v", err)
	}
	if _, ok := s.Net("t2", "r1"); ok {
		t.Fatal("run leaked across tenants")
	}
}

func TestFullRangePercentageNet(t *testing.T) {
	for _, bps := range []int64{1, 2, 1000, 5000, 9999, 10000} {
		t.Run(fmt.Sprint(bps), func(t *testing.T) {
			s := NewStore()
			r := Run{TenantID: "t", ID: "r", EmployeeID: "e", GrossMinorUnits: math.MaxInt64, Deductions: []Deduction{{Kind: KindPercent, BasisPoints: bps}}}
			deduction := new(big.Int).Mul(big.NewInt(math.MaxInt64), big.NewInt(bps))
			deduction.Quo(deduction, big.NewInt(10000))
			want := new(big.Int).Sub(big.NewInt(math.MaxInt64), deduction)
			got, e := s.Execute(r)
			if e != nil || got != want.Int64() {
				t.Fatalf("net=%d want=%d error=%v", got, want.Int64(), e)
			}
		})
	}
}
func TestNegativeNetCannotWrapToStoredSuccess(t *testing.T) {
	s := NewStore()
	r := Run{TenantID: "t", ID: "r", EmployeeID: "e", GrossMinorUnits: 0, Deductions: []Deduction{{Kind: KindFixed, AmountMinorUnits: math.MaxInt64}, {Kind: KindFixed, AmountMinorUnits: math.MaxInt64}}}
	if got, e := s.Execute(r); !errors.Is(e, ErrNegativeNet) || got != 0 {
		t.Fatalf("negative mathematical net accepted: %d %v", got, e)
	}
	if _, ok := s.Net("t", "r"); ok {
		t.Fatal("rejected run stored")
	}
	r.Deductions = nil
	if got, e := s.Execute(r); e != nil || got != 0 {
		t.Fatalf("failure consumed identity: %d %v", got, e)
	}
}

func TestInvalidDeductionPrecedenceAfterExcess(t *testing.T) {
	s := NewStore()
	r := Run{TenantID: "t", ID: "r", EmployeeID: "e", GrossMinorUnits: 0, Deductions: []Deduction{{Kind: KindFixed, AmountMinorUnits: math.MaxInt64}, {Kind: Kind("invalid")}}}
	if n, e := s.Execute(r); !errors.Is(e, ErrInvalidRun) || n != 0 {
		t.Fatalf("invalid-rule precedence changed: %d %v", n, e)
	}
	if _, ok := s.Net("t", "r"); ok {
		t.Fatal("failed run stored")
	}
}
func TestConcurrentMaximumPayrollRunUniquePerTenant(t *testing.T) {
	s := NewStore()
	var wg sync.WaitGroup
	var success atomic.Int64
	for i := 0; i < 64; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			tenant := "a"
			if i%2 == 0 {
				tenant = "b"
			}
			r := Run{TenantID: tenant, ID: "same", EmployeeID: "e", GrossMinorUnits: math.MaxInt64, Deductions: []Deduction{{Kind: KindPercent, BasisPoints: 10000}}}
			n, e := s.Execute(r)
			if e == nil {
				success.Add(1)
				if n != 0 {
					t.Error("100 percent not zero net")
				}
			} else if !errors.Is(e, ErrDuplicate) {
				t.Error(e)
			}
		}(i)
	}
	wg.Wait()
	if success.Load() != 2 {
		t.Fatal("duplicate/tenant scope")
	}
	for _, tenant := range []string{"a", "b"} {
		if n, ok := s.Net(tenant, "same"); !ok || n != 0 {
			t.Fatal("stored net")
		}
	}
}
func FuzzPayrollExactDeductions(f *testing.F) {
	for _, g := range []int64{0, 1, 9999, 10000, math.MaxInt64, -1} {
		for _, ops := range [][]byte{{0, 5, 0, 5}, {1, 3}, {1, 2, 0, 1}, {0, 5, 2, 0}} {
			f.Add(g, ops)
		}
	}
	f.Fuzz(func(t *testing.T, gross int64, ops []byte) {
		if len(ops) > 32 {
			ops = ops[:32]
		}
		values := []int64{0, 1, 9999, 10000, 10001, math.MaxInt64, -1, gross}
		r := Run{TenantID: "t", ID: "r", EmployeeID: "e", GrossMinorUnits: gross}
		for i := 0; i+1 < len(ops); i += 2 {
			d := Deduction{Code: fmt.Sprint(i)}
			v := values[int(ops[i+1])%len(values)]
			switch ops[i] % 3 {
			case 0:
				d.Kind = KindFixed
				d.AmountMinorUnits = v
			case 1:
				d.Kind = KindPercent
				d.BasisPoints = v
			case 2:
				d.Kind = "invalid"
			}
			r.Deductions = append(r.Deductions, d)
		}
		wantErr := error(nil)
		if gross < 0 {
			wantErr = ErrInvalidRun
		}
		// Validate all inputs first, then use exact arithmetic. This independently
		// preserves invalid-rule precedence over a negative mathematical net.
		for _, d := range r.Deductions {
			if d.Kind == KindFixed {
				if d.AmountMinorUnits < 0 {
					wantErr = ErrInvalidRun
				}
			} else if d.Kind == KindPercent {
				if d.BasisPoints < 0 || d.BasisPoints > 10000 {
					wantErr = ErrInvalidRun
				}
			} else {
				wantErr = ErrInvalidRun
			}
		}
		net := big.NewInt(gross)
		if wantErr == nil {
			for _, d := range r.Deductions {
				a := big.NewInt(d.AmountMinorUnits)
				if d.Kind == KindPercent {
					a.Mul(big.NewInt(gross), big.NewInt(d.BasisPoints))
					a.Quo(a, big.NewInt(10000))
				}
				net.Sub(net, a)
			}
			if net.Sign() < 0 {
				wantErr = ErrNegativeNet
			}
		}
		s := NewStore()
		n, e := s.Execute(r)
		if !errors.Is(e, wantErr) {
			t.Fatalf("classification got=%v want=%v", e, wantErr)
		}
		if wantErr != nil {
			if n != 0 {
				t.Fatal("error returned nonzero")
			}
			if _, ok := s.Net("t", "r"); ok {
				t.Fatal("failed result stored")
			}
			r.GrossMinorUnits = 0
			r.Deductions = nil
			if n, e := s.Execute(r); e != nil || n != 0 {
				t.Fatal("failure consumed identity")
			}
		} else {
			if n != net.Int64() || n < 0 || n > gross {
				t.Fatalf("net got=%d want=%s", n, net)
			}
			if stored, ok := s.Net("t", "r"); !ok || stored != n {
				t.Fatal("stored net mismatch")
			}
			if _, e := s.Execute(r); !errors.Is(e, ErrDuplicate) {
				t.Fatal("duplicate not rejected")
			}
		}
	})
}
````


## 6. Configuration surface

Sin variables ni secretos (las reglas de deducción se inyectan por run).

## 7. Dependency bill

| Package/image/tool | Pin exacto | Uso | Licencia | Runtime/build | Fuente oficial |
|---|---|---|---|---|---|
| Go standard library | go 1.26.7 (toolchain) | nómina | BSD-3-Clause | runtime | https://go.dev |

## 8. Apply order

1. Componer el backend (mismo módulo).
2. Colocar los dos archivos bajo `internal/payroll/`.
3. Verificar con `go test ./... -count=1` y `go vet ./...`.
4. Rollback: eliminar `internal/payroll/`.

## 9. Verification

- `go test ./internal/payroll/ -count=1`: 5/5 PASS.
- `go test ./... -count=1` (37 paquetes): PASS.
- `go vet ./...`: exit 0.

## 10. Reconstruction evidence

Véase `reconstruction_evidence/GO_PAYROLL_CORE_2026-09-02_V216.md`.

## Revisión V356 — límites de int64

0.1.1 corrige overflow sin cambiar firmas ni reglas. Cálculo percentual con
cociente/resto conserva floor(amount*bps/10000); payroll limita cada resta
y sigue validando reglas posteriores para conservar prioridad de error.
AUTHORED/LicenseRef-Workspace-Owner, SUPPORTED_REFERENCE/CONDITIONED persisten.
Compatibilidad histórica declarada no demuestra integración con owners actuales.
Ver reconstruction_evidence/BASIS_POINT_ARITHMETIC_V356.md para red/green,
fuentes canónicas, modelo exacto y límites. Las suites V195/V216 son históricas;
no se presenta una composición empresarial antigua como reejecutada.

V356 verificado:4/4fuentes entre ambos cores,18tests ordinarios x3,vet/build
y gate2targets fuzz10s/4workers PASS;37semillas y4614438ejecuciones.
Los11tests históricos conservan cuerpos/aserciones. Sin cambio de API, reglas,
truncado ni nueva dependencia; no equivale a admisión financiera/productiva.

## Revisión de compatibilidad V365

Metadata 0.1.2; las dos fuentes materializables conservan sus SHA y APIs.
Declaración anterior: `["GO-ENTERPRISE-BACKEND-CORE 0.1.x"]`. Se conserva como historial, no como
compatibilidad actual demostrada. El módulo importa sólo stdlib; no existe un
adapter a los owners citados en estos dos archivos. No se sustituye el rango por
un latest ni se admite un CANDIDATE. El análisis de versión/identidad y las suites
del conjunto están en reconstruction_evidence/CORE_COMPATIBILITY_AUDIT_V365.md.
La equivalencia funcional con el perfil integral sigue pendiente por los owners
de V293; esta corrección no agrega módulos ni reduce requisitos.

## Auditoría por claim V374

La revisión 0.1.3 se reconstruye y prueba en Go1.26.8/Windows con fuentes
oficiales fijadas y comparaciones por contrato. Ver
`reconstruction_evidence/CORE_CLAIM_ADMISSION_V374.md` desde la raíz: identidad,
licencia local, prueba por invariante, fuente semántica y dictamen de equivalencia
se mantienen separados. El lenguaje/stdlib no es fuente de un negocio empresarial.
No cambia este claim, API, permiso de composición ni las condiciones de integración.
La revisión 0.1.2 y sus bytes quedan preservados en el expediente anterior.
