# Go Promotions Core

## 1. Metadata

```yaml
pack_id: "GO-PROMOTIONS-CORE"
pack_version: "0.1.3"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa los cupones/descuentos tenant-scoped: percent o fixed, con expiración, monto mínimo, límite de uso y redención fail-closed."
stacks: ["Go 1.26.8"]
compatible_with: ["módulo aislado Go 1.26.7; integración con owners empresariales no admitida"]
incompatible_with: ["cupón sin expiración", "uso sin límite", "descuento cross-tenant"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://go.dev"]
verified_at: "2026-09-10"
```

## 2. Applicability

Referencia AUTHORED en memoria para el cálculo/contador de cupones. No admite por sí sola checkout, política monetaria, autorización, persistencia ni vínculo idempotente por pedido. Conserva truncado entero existente, límite de usos y expiración inclusiva; cada Redeem exitoso consume un uso, incluso con descuento cero.

## 3. Architecture contract

- **Ownership**: `internal/promotions` gobierna cupones y uso; el precio final lo resuelve el dominio.
- **Invariantes**: (1) percent ≤ 10000 basis points, fixed > 0. (2) expiración obligatoria. (3) uso ≤ MaxUses. (4) pedido ≥ mínimo. (5) tenant-scoped.
- **Data flow**: `Create` → `Redeem` (validar + consumir + descuento).
- **Failure modes**: expirado, inactivo, bajo mínimo, agotado, duplicado → error.
- **Seguridad/privacidad**: tenant-scoped.
- **Performance budget**: O(1).
- **Operación/migración/rollback**: sin migración.

## 4. Exact file manifest

```text
CREATE internal/promotions/promotions.go
CREATE internal/promotions/promotions_test.go
```

## 5. Materialization blocks

### FILE: `internal/promotions/promotions.go`
```yaml
block_id: "GO-PROMOTIONS-CORE:internal/promotions/promotions.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "81da2ae2b9743532f1bcbd2a691980ad71b797ec6dc8ed0a9d24cb1bdcda8467"
variables: []
secrets_allowed: false
```
````go
// Package promotions provides tenant-scoped discount coupons: percent or fixed
// discount with expiry, minimum order, usage limits and fail-closed redemption.
package promotions

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"
)

// Kind is percent or fixed discount.
type Kind string

const (
	KindPercent Kind = "percent"
	KindFixed   Kind = "fixed"
)

var (
	codeRe = regexp.MustCompile(`^[A-Z0-9][A-Z0-9._-]{0,31}$`)

	ErrInvalidCoupon = errors.New("promotions: invalid coupon")
	ErrExpired       = errors.New("promotions: expired")
	ErrInactive      = errors.New("promotions: inactive")
	ErrBelowMinimum  = errors.New("promotions: below minimum order")
	ErrExhausted     = errors.New("promotions: usage limit reached")
	ErrNotFound      = errors.New("promotions: coupon not found")
)

// Coupon is a tenant-scoped discount code.
type Coupon struct {
	TenantID           string
	Code               string
	Kind               Kind
	Value              int64 // percent basis points (0..10000) or fixed minor units
	MinOrderMinorUnits int64
	MaxUses            int
	ExpiresAt          time.Time
	Active             bool
}

// Validate enforces the coupon contract.
func (c Coupon) Validate() error {
	if strings.TrimSpace(c.TenantID) == "" || !codeRe.MatchString(c.Code) {
		return ErrInvalidCoupon
	}
	switch c.Kind {
	case KindPercent:
		if c.Value <= 0 || c.Value > 10000 {
			return ErrInvalidCoupon
		}
	case KindFixed:
		if c.Value <= 0 {
			return ErrInvalidCoupon
		}
	default:
		return ErrInvalidCoupon
	}
	if c.MinOrderMinorUnits < 0 || c.MaxUses < 1 {
		return ErrInvalidCoupon
	}
	if c.ExpiresAt.IsZero() {
		return ErrInvalidCoupon
	}
	return nil
}

// Store holds coupons and their usage.
type Store struct {
	mu      sync.Mutex
	coupons map[string]Coupon
	used    map[string]int
	clock   func() time.Time
}

// NewStore returns an empty coupon store.
func NewStore() *Store {
	return &Store{
		coupons: make(map[string]Coupon),
		used:    make(map[string]int),
		clock:   time.Now,
	}
}

func (s *Store) now() time.Time {
	if s.clock != nil {
		return s.clock()
	}
	return time.Now()
}

func key(tenant, code string) string { return tenant + "\x00" + code }

// Create adds a validated coupon.
func (s *Store) Create(c Coupon) error {
	if err := c.Validate(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.coupons[key(c.TenantID, c.Code)]; ok {
		return fmt.Errorf("promotions: duplicate code %q", c.Code)
	}
	s.coupons[key(c.TenantID, c.Code)] = c
	s.used[key(c.TenantID, c.Code)] = 0
	return nil
}

// Redeem validates and consumes one use, returning the discount amount.
func (s *Store) Redeem(tenant, code string, orderMinorUnits int64) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	c, ok := s.coupons[key(tenant, code)]
	if !ok {
		return 0, ErrNotFound
	}
	if !c.Active {
		return 0, ErrInactive
	}
	if s.now().After(c.ExpiresAt) {
		return 0, ErrExpired
	}
	if orderMinorUnits < c.MinOrderMinorUnits {
		return 0, ErrBelowMinimum
	}
	if s.used[key(tenant, code)] >= c.MaxUses {
		return 0, ErrExhausted
	}
	s.used[key(tenant, code)]++

	switch c.Kind {
	case KindPercent:
		// Preserve floor(order*bps/10000) without overflowing the product.
		return (orderMinorUnits/10000)*c.Value + (orderMinorUnits%10000)*c.Value/10000, nil
	case KindFixed:
		if c.Value < orderMinorUnits {
			return c.Value, nil
		}
		return orderMinorUnits, nil
	default:
		return 0, ErrInvalidCoupon
	}
}
````

### FILE: `internal/promotions/promotions_test.go`
```yaml
block_id: "GO-PROMOTIONS-CORE:internal/promotions/promotions_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "837205419f3c2e4b11ceaf1810e1fac600d1600239305760de21394bc8c5e795"
variables: []
secrets_allowed: false
```
````go
package promotions

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"sync"
	"sync/atomic"
	"testing"

	"time"
)

func coupon(code string) Coupon {
	return Coupon{
		TenantID: "t", Code: code, Kind: KindPercent, Value: 1000,
		MinOrderMinorUnits: 1000, MaxUses: 2, ExpiresAt: time.Now().Add(24 * time.Hour), Active: true,
	}
}

func TestRedeemPercent(t *testing.T) {
	s := NewStore()
	_ = s.Create(coupon("SAVE10"))
	d, err := s.Redeem("t", "SAVE10", 5000)
	if err != nil || d != 500 {
		t.Fatalf("expected 500 discount, got %d/%v", d, err)
	}
}

func TestRedeemFixedCappedAtOrder(t *testing.T) {
	s := NewStore()
	c := coupon("FIXED")
	c.Kind = KindFixed
	c.Value = 1000
	c.MinOrderMinorUnits = 0
	_ = s.Create(c)
	d, _ := s.Redeem("t", "FIXED", 500)
	if d != 500 {
		t.Fatalf("fixed should cap at order, got %d", d)
	}
}

func TestUsageLimitExhausted(t *testing.T) {
	s := NewStore()
	_ = s.Create(coupon("TWICE"))
	_, _ = s.Redeem("t", "TWICE", 5000)
	_, _ = s.Redeem("t", "TWICE", 5000)
	if _, err := s.Redeem("t", "TWICE", 5000); !errors.Is(err, ErrExhausted) {
		t.Fatalf("expected ErrExhausted, got %v", err)
	}
}

func TestExpiredAndInactiveAndMinimum(t *testing.T) {
	s := NewStore()
	_ = s.Create(coupon("EXPIRED"))
	expired := coupon("EXPIRED")
	expired.ExpiresAt = time.Now().Add(-time.Hour)
	s.coupons[key("t", "EXPIRED")] = expired
	if _, err := s.Redeem("t", "EXPIRED", 5000); !errors.Is(err, ErrExpired) {
		t.Fatalf("expected ErrExpired, got %v", err)
	}

	_ = s.Create(coupon("OFF"))
	inactive := coupon("OFF")
	inactive.Active = false
	s.coupons[key("t", "OFF")] = inactive
	if _, err := s.Redeem("t", "OFF", 5000); !errors.Is(err, ErrInactive) {
		t.Fatalf("expected ErrInactive, got %v", err)
	}

	_ = s.Create(coupon("MIN"))
	if _, err := s.Redeem("t", "MIN", 500); !errors.Is(err, ErrBelowMinimum) {
		t.Fatalf("expected ErrBelowMinimum, got %v", err)
	}
}

func TestValidateRejectsBadValues(t *testing.T) {
	s := NewStore()
	bad := coupon("BAD")
	bad.Value = 0
	if err := s.Create(bad); !errors.Is(err, ErrInvalidCoupon) {
		t.Fatalf("zero value accepted: %v", err)
	}
	dup := coupon("DUP")
	_ = s.Create(dup)
	if err := s.Create(dup); err == nil {
		t.Fatal("duplicate code accepted")
	}
}

func TestTenantIsolation(t *testing.T) {
	s := NewStore()
	_ = s.Create(coupon("SAVE10"))
	if _, err := s.Redeem("other", "SAVE10", 5000); !errors.Is(err, ErrNotFound) {
		t.Fatalf("cross-tenant should be not found, got %v", err)
	}
}

func TestFullRangePercentageDiscount(t *testing.T) {
	for _, bps := range []int64{1, 2, 1000, 5000, 9999, 10000} {
		t.Run(fmt.Sprint(bps), func(t *testing.T) {
			s := NewStore()
			now := time.Unix(1000000, 0)
			s.clock = func() time.Time { return now }
			c := Coupon{TenantID: "t", Code: "C", Kind: KindPercent, Value: bps, MaxUses: 1, ExpiresAt: now.Add(time.Hour), Active: true}
			if e := s.Create(c); e != nil {
				t.Fatal(e)
			}
			exact := new(big.Int).Mul(big.NewInt(math.MaxInt64), big.NewInt(bps))
			exact.Quo(exact, big.NewInt(10000))
			got, e := s.Redeem("t", "C", math.MaxInt64)
			if e != nil || got != exact.Int64() {
				t.Fatalf("discount=%d want=%d error=%v", got, exact.Int64(), e)
			}
		})
	}
}

func TestConcurrentCouponLimitAtMaximumAmount(t *testing.T) {
	s := NewStore()
	now := time.Unix(1000000, 0)
	s.clock = func() time.Time { return now }
	c := Coupon{TenantID: "t", Code: "C", Kind: KindPercent, Value: 10000, MaxUses: 7, ExpiresAt: now.Add(time.Hour), Active: true}
	if e := s.Create(c); e != nil {
		t.Fatal(e)
	}
	var wg sync.WaitGroup
	var success atomic.Int64
	for i := 0; i < 64; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			d, e := s.Redeem("t", "C", math.MaxInt64)
			if e == nil {
				if d != math.MaxInt64 {
					t.Errorf("wrong discount %d", d)
				}
				success.Add(1)
			} else if !errors.Is(e, ErrExhausted) {
				t.Error(e)
			}
		}()
	}
	wg.Wait()
	if success.Load() != 7 {
		t.Fatal("usage limit not atomic")
	}
}
func TestCouponBoundaryAndRejectedUsePreserved(t *testing.T) {
	s := NewStore()
	now := time.Unix(1000000, 0)
	s.clock = func() time.Time { return now }
	c := Coupon{TenantID: "t", Code: "C", Kind: KindPercent, Value: 1000, MinOrderMinorUnits: 1, MaxUses: 1, ExpiresAt: now, Active: true}
	if e := s.Create(c); e != nil {
		t.Fatal(e)
	}
	if d, e := s.Redeem("other", "C", math.MaxInt64); !errors.Is(e, ErrNotFound) || d != 0 {
		t.Fatal("foreign redemption")
	}
	if d, e := s.Redeem("t", "C", 0); !errors.Is(e, ErrBelowMinimum) || d != 0 {
		t.Fatal("invalid order")
	}
	// Preserve the existing inclusive expiry boundary and truncation rule.
	if d, e := s.Redeem("t", "C", 19); e != nil || d != 1 {
		t.Fatalf("boundary/rounding changed: %d %v", d, e)
	}
}
func FuzzPercentageCoupon(f *testing.F) {
	for _, p := range [][2]int64{{0, 1}, {1, 1}, {9999, 9999}, {10000, 10000}, {10001, 5000}, {math.MaxInt64, 1}, {math.MaxInt64, 2}, {math.MaxInt64, 1000}, {math.MaxInt64, 9999}, {math.MaxInt64, 10000}, {-1, 1000}, {1, 0}, {1, 10001}} {
		f.Add(p[0], p[1])
	}
	f.Fuzz(func(t *testing.T, order, bps int64) {
		s := NewStore()
		now := time.Unix(1000000, 0)
		s.clock = func() time.Time { return now }
		c := Coupon{TenantID: "t", Code: "C", Kind: KindPercent, Value: bps, MaxUses: 2, ExpiresAt: now.Add(time.Hour), Active: true}
		e := s.Create(c)
		if bps <= 0 || bps > 10000 {
			if !errors.Is(e, ErrInvalidCoupon) {
				t.Fatal("invalid bps accepted")
			}
			return
		}
		if e != nil {
			t.Fatal(e)
		}
		d, e := s.Redeem("t", "C", order)
		remaining := 1
		if order < 0 {
			if !errors.Is(e, ErrBelowMinimum) || d != 0 {
				t.Fatal("negative order accepted")
			}
			remaining = 2
		} else {
			want := new(big.Int).Mul(big.NewInt(order), big.NewInt(bps))
			want.Quo(want, big.NewInt(10000))
			if e != nil || d != want.Int64() || d < 0 || d > order {
				t.Fatalf("oracle got=%d want=%s err=%v", d, want, e)
			}
		}
		for i := 0; i < remaining; i++ {
			if d, e := s.Redeem("t", "C", 10000); e != nil || d != bps {
				t.Fatal("usage mutated incorrectly")
			}
		}
		if d, e := s.Redeem("t", "C", 10000); !errors.Is(e, ErrExhausted) || d != 0 {
			t.Fatal("usage not exhausted")
		}
	})
}
````


## 6. Configuration surface

Sin variables ni secretos.

## 7. Dependency bill

| Package/image/tool | Pin exacto | Uso | Licencia | Runtime/build | Fuente oficial |
|---|---|---|---|---|---|
| Go standard library | go 1.26.7 (toolchain) | cupones | BSD-3-Clause | runtime | https://go.dev |

## 8. Apply order

1. Componer el backend (mismo módulo).
2. Colocar los dos archivos bajo `internal/promotions/`.
3. Verificar con `go test ./... -count=1` y `go vet ./...`.
4. Rollback: eliminar `internal/promotions/`.

## 9. Verification

- `go test ./internal/promotions/ -count=1`: 6/6 PASS (percent, fixed con tope, límite agotado, expirado/inactivo/mínimo, valores inválidos, aislamiento tenant).
- `go test ./... -count=1` (17 paquetes): PASS.
- `go vet ./...`: exit 0.

## 10. Reconstruction evidence

Véase `reconstruction_evidence/GO_PROMOTIONS_CORE_2026-09-02_V195.md`.

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
Declaración anterior: `["GO-ENTERPRISE-BACKEND 0.4.x", "GO-COMMERCE-PRICING-PAYMENT-API 0.3.x"]`. Se conserva como historial, no como
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
