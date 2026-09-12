# Go SLO Core

## 1. Metadata

```yaml
pack_id: "GO-SLO-CORE"
pack_version: "0.1.3"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Calculadora AUTHORED en memoria de disponibilidad/presupuesto/burn rate para target finito en (0,1), más comparación de dos budgets suministrados por el caller; sin ventanas temporales ni alertas desplegadas."
stacks: ["Go 1.26.8"]
compatible_with: ["módulo aislado Go 1.26.7; integración con owners empresariales no admitida"]
incompatible_with: ["target fuera de (0,1)"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://go.dev", "https://sre.google"]
verified_at: "2026-09-10"
```

## 2. Applicability

Referencia AUTHORED en memoria: rechaza NaN y target fuera de (0,1). Budget sólo acumula eventos; no implementa ventana temporal, expiración, reset, collector o paginación. MultiWindowAlert compara dos budgets suministrados por el caller; no prueba despliegue de alertas ni validación de thresholds. Counters int64 requieren una vida/presupuesto de eventos acotada en el target.

## 3. Architecture contract

- **Ownership**: `internal/slo` calcula el presupuesto; la telemetría que alimenta `Record` es del runtime.
- **Invariantes**: (1) target en (0,1). (2) disponibilidad 1 cuando vacío. (3) burn rate 0 sin errores. (4) alerta sólo cuando ambas ventanas queman.
- **Data flow**: `NewBudget` → `Record` → `Availability/ErrorBudgetRemaining/BurnRate/Exhausted`.
- **Failure modes**: target inválido → error.
- **Seguridad/privacidad**: sin PII.
- **Performance budget**: O(1) por record.

## 4. Exact file manifest

```text
CREATE internal/slo/slo.go
CREATE internal/slo/slo_test.go
```

## 5. Materialization blocks

### FILE: `internal/slo/slo.go`
```yaml
block_id: "GO-SLO-CORE:internal/slo/slo.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "900b34d5ae9ef847176fda3b6f33b131df02cfe88002222e8cb54492ec520962"
variables: []
secrets_allowed: false
```
````go
// Package slo provides AUTHORED in-memory availability and error-budget
// arithmetic. Burn-rate formulas are compared with Google SRE; this package
// has no time-window lifecycle, SLI ingestion or deployed alerting.
package slo

import (
	"errors"
	"math"
	"sync"
)

var (
	ErrInvalidTarget = errors.New("slo: target must be in (0,1)")
)

// Budget tracks total events and errors toward an availability target.
type Budget struct {
	mu     sync.Mutex
	target float64 // availability target in (0,1), e.g. 0.999
	total  int64
	errors int64
}

// NewBudget returns a budget with the given availability target.
func NewBudget(target float64) (*Budget, error) {
	if math.IsNaN(target) || target <= 0 || target >= 1 {
		return nil, ErrInvalidTarget
	}
	return &Budget{target: target}, nil
}

// Record registers one event (ok=true) or one error (ok=false).
func (b *Budget) Record(ok bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.total++
	if !ok {
		b.errors++
	}
}

// Availability returns the observed availability (1 when empty).
func (b *Budget) Availability() float64 {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.total == 0 {
		return 1
	}
	return float64(b.total-b.errors) / float64(b.total)
}

// Errors returns the error count.
func (b *Budget) Errors() int64 {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.errors
}

// Total returns the event count.
func (b *Budget) Total() int64 {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.total
}

// allowedErrors returns the error budget in events.
func (b *Budget) allowedErrors() float64 {
	return float64(b.total) * (1 - b.target)
}

// ErrorBudgetRemaining returns the remaining error budget as a fraction in
// [0,1]; 1 when nothing consumed, 0 when exhausted.
func (b *Budget) ErrorBudgetRemaining() float64 {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.total == 0 {
		return 1
	}
	allowed := b.allowedErrors()
	if allowed <= 0 {
		return 0
	}
	remaining := 1 - float64(b.errors)/allowed
	if remaining < 0 {
		return 0
	}
	if remaining > 1 {
		return 1
	}
	return remaining
}

// BurnRate returns how many times faster than allowed the budget is burning
// (errors/total divided by 1-target). 0 when empty or no errors.
func (b *Budget) BurnRate() float64 {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.total == 0 || b.errors == 0 {
		return 0
	}
	allowed := 1 - b.target
	if allowed <= 0 {
		return 0
	}
	return (float64(b.errors) / float64(b.total)) / allowed
}

// Exhausted reports whether the error budget is fully consumed.
func (b *Budget) Exhausted() bool { return b.ErrorBudgetRemaining() <= 0 }

// MultiWindowAlert compares two caller-supplied cumulative budgets against
// their thresholds. The caller owns window population, alignment and lifecycle;
// this Boolean helper does not implement a Google SRE alerting deployment.
func MultiWindowAlert(short, long *Budget, shortThreshold, longThreshold float64) bool {
	if short == nil || long == nil {
		return false
	}
	return short.BurnRate() > shortThreshold && long.BurnRate() > longThreshold
}
````

### FILE: `internal/slo/slo_test.go`
```yaml
block_id: "GO-SLO-CORE:internal/slo/slo_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "8a704820fa654b0844ef4746838e5e769a4c95c314e6d03961254e9602aac08a"
variables: []
secrets_allowed: false
```
````go
package slo

import (
	"errors"
	"math"
	"math/big"
	"sync"
	"testing"
)

func TestInvalidTarget(t *testing.T) {
	for _, v := range []float64{0, 1, 1.5, -0.1} {
		if _, err := NewBudget(v); !errors.Is(err, ErrInvalidTarget) {
			t.Fatalf("target %v accepted: %v", v, err)
		}
	}
}

func TestAvailability(t *testing.T) {
	b, _ := NewBudget(0.99)
	if b.Availability() != 1 {
		t.Fatalf("empty availability should be 1, got %v", b.Availability())
	}
	for i := 0; i < 99; i++ {
		b.Record(true)
	}
	b.Record(false)
	if b.Availability() != 0.99 {
		t.Fatalf("expected 0.99, got %v", b.Availability())
	}
}

func TestErrorBudgetRemaining(t *testing.T) {
	b, _ := NewBudget(0.99)
	for i := 0; i < 100; i++ {
		b.Record(true)
	}
	if b.ErrorBudgetRemaining() != 1 {
		t.Fatalf("expected 1 remaining, got %v", b.ErrorBudgetRemaining())
	}
	b.Record(false) // 101 total, 1 error; allowed = 1.01 → remaining ~0.0099
	if b.ErrorBudgetRemaining() <= 0 || b.ErrorBudgetRemaining() >= 0.1 {
		t.Fatalf("unexpected remaining: %v", b.ErrorBudgetRemaining())
	}
}

func TestBurnRate(t *testing.T) {
	b, _ := NewBudget(0.99)
	// 1000 events, 10 errors = 1% error rate = exactly the SLO budget → burn ~1.
	for i := 0; i < 990; i++ {
		b.Record(true)
	}
	for i := 0; i < 10; i++ {
		b.Record(false)
	}
	br := b.BurnRate()
	if br < 0.9 || br > 1.1 {
		t.Fatalf("expected burn ~1.0, got %v", br)
	}
}

func TestExhausted(t *testing.T) {
	b, _ := NewBudget(0.5) // allowed error rate 50%
	for i := 0; i < 4; i++ {
		b.Record(true)
	}
	for i := 0; i < 6; i++ {
		b.Record(false)
	}
	// 10 total, 6 errors; allowed = 5 → budget exhausted.
	if !b.Exhausted() {
		t.Fatal("budget should be exhausted at 6/10 errors on 0.5 target")
	}
}

func TestMultiWindowAlert(t *testing.T) {
	short, _ := NewBudget(0.99)
	long, _ := NewBudget(0.99)
	// Short spike burns fast; long window still fine → no alert.
	for i := 0; i < 100; i++ {
		short.Record(false)
		long.Record(true)
	}
	if MultiWindowAlert(short, long, 5, 5) {
		t.Fatal("alert should not fire when only short window burns")
	}
	// Both windows burn fast → alert.
	for i := 0; i < 100; i++ {
		long.Record(false)
	}
	if !MultiWindowAlert(short, long, 5, 5) {
		t.Fatal("alert should fire when both windows burn")
	}
}

func TestNaNTargetRejected(t *testing.T) {
	b, e := NewBudget(math.NaN())
	if b != nil || !errors.Is(e, ErrInvalidTarget) {
		t.Fatalf("invalid NaN target admitted: budget=%v err=%v", b, e)
	}
}

func TestFiniteTargetBoundariesAndEmptyBudget(t *testing.T) {
	for _, target := range []float64{math.Inf(1), math.Inf(-1), 0, 1, -1} {
		if b, e := NewBudget(target); b != nil || !errors.Is(e, ErrInvalidTarget) {
			t.Fatal("invalid target", target, e)
		}
	}
	for _, target := range []float64{math.SmallestNonzeroFloat64, math.Nextafter(1, 0), 0.5} {
		b, e := NewBudget(target)
		if e != nil {
			t.Fatal(e)
		}
		if b.Total() != 0 || b.Errors() != 0 || b.Availability() != 1 || b.ErrorBudgetRemaining() != 1 || b.BurnRate() != 0 || b.Exhausted() {
			t.Fatal("empty budget")
		}
	}
}
func TestConcurrentBudgetRecords(t *testing.T) {
	b, e := NewBudget(0.99)
	if e != nil {
		t.Fatal(e)
	}
	var wg sync.WaitGroup
	for i := 0; i < 64; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				b.Record(i%2 == 0)
				b.Availability()
				b.ErrorBudgetRemaining()
				b.BurnRate()
			}
		}(i)
	}
	wg.Wait()
	if b.Total() != 6400 || b.Errors() != 3200 || b.Availability() != 0.5 {
		t.Fatal("lost records")
	}
}
func FuzzFiniteBudgetModel(f *testing.F) {
	for _, target := range []float64{math.NaN(), math.Inf(1), math.Inf(-1), -1, 0, 1, 0.5, 0.99, math.SmallestNonzeroFloat64, math.Nextafter(1, 0)} {
		f.Add(target, []byte{})
		f.Add(target, []byte{0, 1, 0, 1, 1, 1, 0})
	}
	f.Fuzz(func(t *testing.T, target float64, ops []byte) {
		b, e := NewBudget(target)
		if !(target > 0 && target < 1) {
			if b != nil || !errors.Is(e, ErrInvalidTarget) {
				t.Fatal("invalid target accepted", target, e)
			}
			return
		}
		if e != nil || b == nil {
			t.Fatal("valid target rejected", e)
		}
		if len(ops) > 256 {
			ops = ops[:256]
		}
		bad := int64(0)
		close := func(got, want float64) {
			if math.IsNaN(got) || math.IsInf(got, 0) || math.Abs(got-want) > 1e-12*math.Max(1, math.Abs(want)) {
				t.Fatal("rational oracle", got, want)
			}
		}
		for i, op := range ops {
			ok := op%2 == 0
			if !ok {
				bad++
			}
			b.Record(ok)
			count := int64(i + 1)
			if b.Total() != count || b.Errors() != bad {
				t.Fatal("counter mismatch")
			}
			avail, _ := new(big.Rat).SetFrac(big.NewInt(count-bad), big.NewInt(count)).Float64()
			close(b.Availability(), avail)
			allowed := new(big.Rat).Sub(big.NewRat(1, 1), new(big.Rat).SetFloat64(target))
			allowed.Mul(allowed, big.NewRat(count, 1))
			burn := new(big.Rat).Quo(big.NewRat(bad, 1), allowed)
			burnFloat, _ := burn.Float64()
			close(b.BurnRate(), burnFloat)
			left := new(big.Rat).Sub(big.NewRat(1, 1), burn)
			if left.Sign() < 0 {
				left.SetInt64(0)
			}
			remaining, _ := left.Float64()
			close(b.ErrorBudgetRemaining(), remaining)
			if b.Exhausted() != (b.ErrorBudgetRemaining() <= 0) {
				t.Fatal("exhausted mismatch")
			}
		}
	})
}

// V374: SRE error-rate/error-budget formula; this helper has no rolling window.
func TestSourceBurnRateGoldenV374(t *testing.T) {
	b, e := NewBudget(0.99)
	if e != nil {
		t.Fatal(e)
	}
	for i := 0; i < 1000; i++ {
		b.Record(i >= 60)
	}
	if math.Abs(b.BurnRate()-6) > 1e-12 || b.ErrorBudgetRemaining() != 0 {
		t.Fatal("SRE formula mismatch")
	}
	// New observations accumulate; there is no window expiry or SLI ingestion here.
	for i := 0; i < 1000; i++ {
		b.Record(true)
	}
	if b.Total() != 2000 || b.Errors() != 60 || math.Abs(b.BurnRate()-3) > 1e-12 {
		t.Fatal("cumulative contract changed")
	}
}
````


## 6. Configuration surface

Sin variables ni secretos (target se pasa por `NewBudget`).

## 7. Dependency bill

| Package/image/tool | Pin exacto | Uso | Licencia | Runtime/build | Fuente oficial |
|---|---|---|---|---|---|
| Go standard library | go 1.26.7 (toolchain) | SLO/error budget | BSD-3-Clause | runtime | https://go.dev |

## 8. Apply order

1. Componer el backend (mismo módulo).
2. Colocar los dos archivos bajo `internal/slo/`.
3. Verificar con `go test ./... -count=1` y `go vet ./...`.
4. Rollback: eliminar `internal/slo/`.

## 9. Verification

- `go test ./internal/slo/ -count=1`: 6/6 PASS.
- `go test ./... -count=1` (41 paquetes): PASS.
- `go vet ./...`: exit 0.

## 10. Reconstruction evidence

Véase `reconstruction_evidence/GO_SLO_CORE_2026-09-02_V220.md`.

## Revisión V361

0.1.1 corrige la frontera numérica demostrada y conserva las reglas existentes
dentro del dominio representable. AUTHORED/LicenseRef-Workspace-Owner y CONDITIONED
persisten. Suite V220 histórica no reejecutada. Evidencia nueva:
reconstruction_evidence/POS_SLO_BOUNDARIES_AND_DOCKER_SCOPE_V361.md.

V361 verificado entre ambos cores: 4 fuentes reconstruidas, 20 tests x3,
vet/build,34semillas/3381083ejecuciones fuzz PASS.11tests históricos conservan
aserciones; POS agrega comprobación de error por helpers del test y rechaza
2firmas antiguas de resultado único. Sin -race, cobro real o alertas desplegadas.

## Revisión de compatibilidad V365

Metadata 0.1.2; las dos fuentes materializables conservan sus SHA y APIs.
Declaración anterior: `["GO-OBSERVABILITY-CORE 0.1.x"]`. Se conserva como historial, no como
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
