# Go FX Core

## 1. Metadata

```yaml
pack_id: "GO-FX-CORE"
pack_version: "0.1.3"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa un helper AUTHORED en memoria para multiplicar montos binary64 por tasas positivas finitas, aislado por tenant/par, con rechazo sin tasa y error ante resultado no finito o cero. No implementa política monetaria, proveedor de tasas, redondeo decimal ni ledger durable."
stacks: ["Go 1.26.8"]
compatible_with: ["módulo aislado Go 1.26.7; integración con owners empresariales no admitida"]
incompatible_with: ["par sin tasa", "tasa no positiva", "monto no positivo"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://go.dev"]
verified_at: "2026-09-10"
```

## 2. Applicability

Use como referencia aritmética en memoria con tasas inyectadas y códigos de tres letras ASCII mayúsculas. Rechace par sin tasa, input no finito/no positivo o resultado no representable como binary64 positivo finito. No valida pertenencia a ISO4217, vigencia de tasas ni política monetaria; no admite por sí solo checkout, contabilidad o FX productivo.

## 3. Architecture contract

- **Ownership**: `internal/fx` convierte; las tasas de cambio reales son runtime (proveedor de FX).
- **Invariantes**: (1) tasa finita > 0. (2) monto finito > 0. (3) par sin tasa → error. (4) códigos de tres letras ASCII mayúsculas; no catálogo ISO4217. (5) resultado finito > 0 o ErrInvalidResult. (6) una tasa inválida no reemplaza la anterior.
- **Data flow**: `SetRate` → `Convert`.
- **Failure modes**: par sin tasa, NaN/infinito/no positivo → error; overflow a infinito y underflow a cero → ErrInvalidResult. Subnormales positivos representables se conservan; no agrega redondeo monetario.
- **Seguridad/privacidad**: tenant-scoped.
- **Performance budget**: O(1) por conversión.

## 4. Exact file manifest

```text
CREATE internal/fx/fx.go
CREATE internal/fx/fx_test.go
```

## 5. Materialization blocks

### FILE: `internal/fx/fx.go`
```yaml
block_id: "GO-FX-CORE:internal/fx/fx.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "d7e7aa5764c35a72dd99856f4cb93e3eeef67ac3b15740d7e63e7e43b0318bf9"
variables: []
secrets_allowed: false
```
````go
// Package fx provides tenant-scoped currency conversion with rates injected
// at runtime; deny-by-default (no rate → error). No exchange rates are
// hardcoded. AUTHORED.
package fx

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strings"
	"sync"
)

var (
	ccyRe = regexp.MustCompile(`^[A-Z]{3}$`)

	ErrInvalidCurrency = errors.New("fx: invalid currency")
	ErrInvalidRate     = errors.New("fx: rate must be positive")
	ErrInvalidAmount   = errors.New("fx: amount must be positive")
	ErrNoRate          = errors.New("fx: no rate for pair (deny by default)")
	ErrInvalidResult   = errors.New("fx: converted amount must be finite and positive")
)

// Store holds conversion rates keyed by tenant+from+to.
type Store struct {
	mu    sync.Mutex
	rates map[string]float64
}

// NewStore returns an empty store.
func NewStore() *Store {
	return &Store{rates: make(map[string]float64)}
}

func key(tenant, from, to string) string { return tenant + "\x00" + from + "\x00" + to }

// SetRate declares a direct rate (units of `to` per 1 `from`). Positive only.
func (s *Store) SetRate(tenant, from, to string, rate float64) error {
	if strings.TrimSpace(tenant) == "" || !ccyRe.MatchString(from) || !ccyRe.MatchString(to) {
		return ErrInvalidCurrency
	}
	if !finitePositive(rate) {
		return ErrInvalidRate
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rates[key(tenant, from, to)] = rate
	return nil
}

// Convert converts amount from `from` to `to`. Deny-by-default: no direct
// rate → ErrNoRate.
func (s *Store) Convert(tenant, from, to string, amount float64) (float64, error) {
	if strings.TrimSpace(tenant) == "" || !ccyRe.MatchString(from) || !ccyRe.MatchString(to) {
		return 0, ErrInvalidCurrency
	}
	if !finitePositive(amount) {
		return 0, ErrInvalidAmount
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	rate, ok := s.rates[key(tenant, from, to)]
	if !ok {
		return 0, ErrNoRate
	}
	result := amount * rate
	if !finitePositive(result) {
		return 0, ErrInvalidResult
	}
	return result, nil
}

// Rate returns the direct rate, if declared.
func (s *Store) Rate(tenant, from, to string) (float64, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	rate, ok := s.rates[key(tenant, from, to)]
	return rate, ok
}

// String aids debugging.
func (s *Store) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return fmt.Sprintf("fx(%d)", len(s.rates))
}

// finitePositive rejects NaN because its comparison with zero is false.
func finitePositive(value float64) bool { return value > 0 && !math.IsInf(value, 0) }
````

### FILE: `internal/fx/fx_test.go`
```yaml
block_id: "GO-FX-CORE:internal/fx/fx_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "8c9a9bb0b45c284c6c171d9f1306d971a23ac0f120fbc8fd5206306fc1d8b557"
variables: []
secrets_allowed: false
```
````go
package fx

import (
	"errors"
	"math"
	"math/big"
	"testing"
)

func TestConvert(t *testing.T) {
	s := NewStore()
	if err := s.SetRate("t", "USD", "ARS", 1000.0); err != nil {
		t.Fatal(err)
	}
	got, err := s.Convert("t", "USD", "ARS", 10.0)
	if err != nil {
		t.Fatal(err)
	}
	if got != 10000.0 {
		t.Fatalf("expected 10000, got %v", got)
	}
}

func TestDenyByDefault(t *testing.T) {
	s := NewStore()
	if _, err := s.Convert("t", "USD", "ARS", 10.0); !errors.Is(err, ErrNoRate) {
		t.Fatalf("expected no-rate (deny), got %v", err)
	}
}

func TestInvalidRate(t *testing.T) {
	s := NewStore()
	if err := s.SetRate("t", "USD", "ARS", 0); !errors.Is(err, ErrInvalidRate) {
		t.Fatalf("zero rate accepted: %v", err)
	}
	if err := s.SetRate("t", "USD", "ARS", -1); !errors.Is(err, ErrInvalidRate) {
		t.Fatalf("negative rate accepted: %v", err)
	}
}

func TestInvalidInput(t *testing.T) {
	s := NewStore()
	if err := s.SetRate("t", "us", "ARS", 1); !errors.Is(err, ErrInvalidCurrency) {
		t.Fatalf("bad from currency accepted: %v", err)
	}
	if _, err := s.Convert("t", "USD", "ARS", 0); !errors.Is(err, ErrInvalidAmount) {
		t.Fatalf("zero amount accepted: %v", err)
	}
}

func TestTenantIsolation(t *testing.T) {
	s := NewStore()
	if err := s.SetRate("t1", "USD", "ARS", 1000.0); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Convert("t2", "USD", "ARS", 10.0); !errors.Is(err, ErrNoRate) {
		t.Fatalf("rate leaked across tenants: %v", err)
	}
}

func TestNonFiniteRateDoesNotReplacePrior(t *testing.T) {
	for name, v := range map[string]float64{"nan": math.NaN(), "positive-infinity": math.Inf(1), "negative-infinity": math.Inf(-1)} {
		t.Run(name, func(t *testing.T) {
			s := NewStore()
			if e := s.SetRate("t", "USD", "ARS", 2); e != nil {
				t.Fatal(e)
			}
			if e := s.SetRate("t", "USD", "ARS", v); !errors.Is(e, ErrInvalidRate) {
				t.Errorf("invalid rate accepted: %v", e)
			}
			if got, ok := s.Rate("t", "USD", "ARS"); !ok || got != 2 {
				t.Errorf("invalid update replaced valid rate: %v", got)
			}
		})
	}
}
func TestNonFiniteAmountRejected(t *testing.T) {
	for name, v := range map[string]float64{"nan": math.NaN(), "positive-infinity": math.Inf(1), "negative-infinity": math.Inf(-1)} {
		t.Run(name, func(t *testing.T) {
			s := NewStore()
			_ = s.SetRate("t", "USD", "ARS", 2)
			got, e := s.Convert("t", "USD", "ARS", v)
			if !errors.Is(e, ErrInvalidAmount) || got != 0 {
				t.Errorf("invalid amount returned result=%v error=%v", got, e)
			}
		})
	}
}
func TestNonRepresentableProductRejected(t *testing.T) {
	for name, p := range map[string][2]float64{"overflow": {math.MaxFloat64, 2}, "underflow-to-zero": {math.SmallestNonzeroFloat64, 0.25}} {
		t.Run(name, func(t *testing.T) {
			s := NewStore()
			if e := s.SetRate("t", "USD", "ARS", p[0]); e != nil {
				t.Fatal(e)
			}
			got, e := s.Convert("t", "USD", "ARS", p[1])
			if !errors.Is(e, ErrInvalidResult) || got != 0 {
				t.Errorf("nonrepresentable product returned result=%v error=%v", got, e)
			}
		})
	}
}

func TestRepresentablePositiveResultsPreserved(t *testing.T) {
	for _, p := range [][2]float64{{1, math.SmallestNonzeroFloat64}, {0.5, math.SmallestNonzeroFloat64 * 2}, {1, math.MaxFloat64}, {0.1, 3}} {
		s := NewStore()
		if e := s.SetRate("t", "USD", "ARS", p[0]); e != nil {
			t.Fatal(e)
		}
		got, e := s.Convert("t", "USD", "ARS", p[1])
		if e != nil || got != p[0]*p[1] {
			t.Fatalf("representable conversion changed: %v %v", got, e)
		}
	}
}

// Numeric oracle decodes IEEE754 inputs and uses exact rational multiplication;
// it does not call the production validation helper or its multiplication result.
func FuzzFiniteConversion(f *testing.F) {
	seeds := [][2]float64{{1, 1}, {math.NaN(), 1}, {math.Inf(1), 1}, {math.Inf(-1), 1}, {1, math.NaN()}, {1, math.Inf(1)}, {1, math.Inf(-1)}, {math.MaxFloat64, 2}, {math.SmallestNonzeroFloat64, 0.25}, {1, math.SmallestNonzeroFloat64}, {0, 1}, {math.Copysign(0, -1), 1}, {1, -1}, {0.1, 3}}
	for _, p := range seeds {
		f.Add(math.Float64bits(p[0]), math.Float64bits(p[1]))
	}
	f.Fuzz(func(t *testing.T, rb, ab uint64) {
		valid := func(bits uint64) bool { return bits>>63 == 0 && bits != 0 && (bits>>52)&0x7ff != 0x7ff }
		rate, amount := math.Float64frombits(rb), math.Float64frombits(ab)
		s := NewStore()
		e := s.SetRate("t", "USD", "ARS", rate)
		if !valid(rb) {
			if !errors.Is(e, ErrInvalidRate) {
				t.Fatal("rate guard")
			}
			if _, ok := s.Rate("t", "USD", "ARS"); ok {
				t.Fatal("invalid rate stored")
			}
			return
		}
		if e != nil {
			t.Fatal(e)
		}
		got, e := s.Convert("t", "USD", "ARS", amount)
		if !valid(ab) {
			if !errors.Is(e, ErrInvalidAmount) || got != 0 {
				t.Fatal("amount guard")
			}
			return
		}
		exact := new(big.Rat).Mul(new(big.Rat).SetFloat64(rate), new(big.Rat).SetFloat64(amount))
		want, _ := exact.Float64()
		if !valid(math.Float64bits(want)) {
			if !errors.Is(e, ErrInvalidResult) || got != 0 {
				t.Fatal("result guard")
			}
		} else if e != nil || math.Float64bits(got) != math.Float64bits(want) {
			t.Fatalf("rational oracle mismatch: got=%g want=%g error=%v", got, want, e)
		}
		if r, ok := s.Rate("t", "USD", "ARS"); !ok || math.Float64bits(r) != rb {
			t.Fatal("conversion mutated rate")
		}
	})
}
````


## 6. Configuration surface

Sin variables ni secretos (las tasas se inyectan por par).

## 7. Dependency bill

| Package/image/tool | Pin exacto | Uso | Licencia | Runtime/build | Fuente oficial |
|---|---|---|---|---|---|
| Go standard library | go 1.26.7 (toolchain) | multi-moneda | BSD-3-Clause | runtime | https://go.dev |

## 8. Apply order

1. Componer el backend (mismo módulo).
2. Colocar los dos archivos bajo `internal/fx/`.
3. Verificar con `go test ./... -count=1` y `go vet ./...`.
4. Rollback: eliminar `internal/fx/`.

## 9. Verification

V353 conserva5tests históricos y añade4tests de regresión. FuzzFiniteConversion
usa14semillas, interpreta inputs IEEE754 por bits y contrasta el resultado con
multiplicación racional exacta independiente del helper de producción.

- Reconstrucción canónica:2/2fuentes idénticas al candidato corregido.
- go test ./... -count=3:9tests ordinarios y14semillas PASS; go vet/build PASS.
- GO_NATIVE_FUZZ_GATE10s:14semillas,13.047.867ejecuciones PASS en Go1.26.7 Windows/amd64.
- Este resultado prueba el contrato aritmético local; no cierra TEST02 ni seguridad integral.
- V217/37paquetes permanece evidencia histórica; no afirma reejecución de esa
  composición ni validación de una política financiera por este delta.

## 10. Reconstruction evidence

V353: reconstruction_evidence/ACCEPTANCE_REAUDIT_FX_BOUNDARY_V353.md.
Antes: reconstruction_evidence/GO_FX_CORE_2026-09-02_V217.md.

Ambos archivos siguen AUTHORED/LicenseRef-Workspace-Owner. Go/math gobierna el
lenguaje/runtime, no representa arquitectura financiera copiada de una empresa.
Mantener CONDITIONED hasta policy/fuente/versiones de tasas, precisión/redondeo,
autorización/persistencia e integración reales del proyecto.

## Revisión de compatibilidad V365

Metadata 0.1.2; las dos fuentes materializables conservan sus SHA y APIs.
Declaración anterior: `["GO-COMMERCE-PRICING-PAYMENT-API 0.1.x"]`. Se conserva como historial, no como
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
