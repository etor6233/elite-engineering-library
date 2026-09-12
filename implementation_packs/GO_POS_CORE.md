# Go POS Core

## 1. Metadata

```yaml
pack_id: "GO-POS-CORE"
pack_version: "0.2.1"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Referencia AUTHORED en memoria: total/vuelto con errores explícitos, registro de venta por tenant/id y copia de líneas; no implementa cobro ni checkout."
stacks: ["Go 1.26.8"]
compatible_with: ["uso aislado; composición con cobro/ledger no admitida"]
incompatible_with: ["callers que ignoran errores de importe", "total no representable", "duplicado dentro del mismo tenant"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://go.dev"]
verified_at: "2026-09-10"
```

## 2. Applicability

Referencia AUTHORED aritmética en memoria. Total/ChangeDue devuelven (int64,error), con ErrOverflow para importes no representables; Complete conserva rechazo de total no positivo y tender insuficiente. No checkout, cobro, caja, fiscalidad, ledger ni persistencia admitidos.

## 3. Architecture contract

- **Ownership**: `internal/pos` gobierna la transacción de venta; el cobro real es del dominio de pagos.
- **Invariantes**: (1) líneas no vacías, qty > 0, precio ≥ 0. (2) total > 0. (3) tender ≥ total (vuelto ≥ 0). (4) duplicados por (tenant,id) rechazados con ErrDuplicate, sin replay exitoso. (5) tenant-scoped.
- **Data flow**: `Complete` → valida → registra.
- **Failure modes**: sin líneas, cantidad 0, tender corto, duplicado → error.
- **Seguridad/privacidad**: tenant-scoped.
- **Performance budget**: O(N) Complete/Total/ChangeDue sobre líneas; copia de líneas al registrar.

## 4. Exact file manifest

```text
CREATE internal/pos/pos.go
CREATE internal/pos/pos_test.go
```

## 5. Materialization blocks

### FILE: `internal/pos/pos.go`
```yaml
block_id: "GO-POS-CORE:internal/pos/pos.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "4d3be0f26614884e59ba2c052c5d2a29fcb5ce0200719c7da611e8750951ee41"
variables: []
secrets_allowed: false
```
````go
// Package pos provides a tenant-scoped point-of-sale transaction: line items,
// checked total, tender and change due. Duplicates are rejected, not replayed.
package pos

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"sync"
	"time"
)

var (
	ErrInvalidSale = errors.New("pos: invalid sale")
	ErrDuplicate   = errors.New("pos: duplicate sale")
	ErrShortTender = errors.New("pos: tender below total")
	ErrOverflow    = errors.New("pos: total exceeds int64")
)

// LineItem is a single sold item.
type LineItem struct {
	ProductID           string
	Quantity            int
	UnitPriceMinorUnits int64
}

// Validate enforces the line contract.
func (l LineItem) Validate() error {
	if strings.TrimSpace(l.ProductID) == "" || l.Quantity <= 0 || l.UnitPriceMinorUnits < 0 {
		return ErrInvalidSale
	}
	return nil
}

// Sale is a completed POS transaction.
type Sale struct {
	TenantID         string
	ID               string
	CashierID        string
	Lines            []LineItem
	TenderMinorUnits int64
	CompletedAt      time.Time
}

// Store holds completed sales.
type Store struct {
	mu    sync.Mutex
	sales map[identity]Sale
}

type identity struct{ tenant, id string }

// NewStore returns an empty store.
func NewStore() *Store {
	return &Store{sales: make(map[identity]Sale)}
}

// Total validates lines and computes their exact nonnegative int64 total.
// Empty/invalid lines and unrepresentable arithmetic return errors, not a value.
func (s *Store) Total(sale Sale) (int64, error) {
	if len(sale.Lines) == 0 {
		return 0, ErrInvalidSale
	}
	// Validate all lines first, preserving invalid-line precedence.
	for _, l := range sale.Lines {
		if err := l.Validate(); err != nil {
			return 0, err
		}
	}
	var total int64
	for _, l := range sale.Lines {
		q := int64(l.Quantity)
		if l.UnitPriceMinorUnits > math.MaxInt64/q {
			return 0, ErrOverflow
		}
		subtotal := q * l.UnitPriceMinorUnits
		if subtotal > math.MaxInt64-total {
			return 0, ErrOverflow
		}
		total += subtotal
	}
	return total, nil
}

// ChangeDue returns nonnegative tender minus a validated, representable total.
func (s *Store) ChangeDue(sale Sale) (int64, error) {
	total, err := s.Total(sale)
	if err != nil {
		return 0, err
	}
	if sale.TenderMinorUnits < total {
		return 0, ErrShortTender
	}
	return sale.TenderMinorUnits - total, nil
}

// Complete validates and records the sale, fail-closed on short tender.
func (s *Store) Complete(sale Sale) error {
	if strings.TrimSpace(sale.TenantID) == "" || strings.TrimSpace(sale.ID) == "" || strings.TrimSpace(sale.CashierID) == "" {
		return ErrInvalidSale
	}
	total, err := s.Total(sale)
	if err != nil {
		return err
	}
	if total <= 0 {
		return ErrInvalidSale
	}
	if sale.TenderMinorUnits < total {
		return ErrShortTender
	}
	sale.Lines = append([]LineItem(nil), sale.Lines...)
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.sales[identity{tenant: sale.TenantID, id: sale.ID}]; ok {
		return ErrDuplicate
	}
	if sale.CompletedAt.IsZero() {
		sale.CompletedAt = time.Now().UTC()
	}
	s.sales[identity{tenant: sale.TenantID, id: sale.ID}] = sale
	return nil
}

// String aids debugging.
func (s *Store) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return fmt.Sprintf("pos(%d)", len(s.sales))
}
````

### FILE: `internal/pos/pos_test.go`
```yaml
block_id: "GO-POS-CORE:internal/pos/pos_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "ee7fb029c87abb828b911149f722c9f364558c62d4bf908f23be3cd32e30555c"
variables: []
secrets_allowed: false
```
````go
package pos

import (
	"errors"
	"math"
	"math/big"
	"sync"
	"sync/atomic"
	"testing"
)

func sale() Sale {
	return Sale{
		TenantID: "t", ID: "S1", CashierID: "cashier-1",
		Lines: []LineItem{
			{ProductID: "P1", Quantity: 2, UnitPriceMinorUnits: 500},
			{ProductID: "P2", Quantity: 1, UnitPriceMinorUnits: 1000},
		},
		TenderMinorUnits: 2500,
	}
}

func TestTotalAndChange(t *testing.T) {
	s := NewStore()
	if mustTotal(t, s, sale()) != 2000 {
		t.Fatalf("expected 2000, got %d", mustTotal(t, s, sale()))
	}
	if mustChange(t, s, sale()) != 500 {
		t.Fatalf("expected 500 change, got %d", mustChange(t, s, sale()))
	}
}

func TestCompleteAndDuplicate(t *testing.T) {
	s := NewStore()
	if err := s.Complete(sale()); err != nil {
		t.Fatal(err)
	}
	if err := s.Complete(sale()); !errors.Is(err, ErrDuplicate) {
		t.Fatalf("duplicate accepted: %v", err)
	}
}

func TestShortTenderFailsClosed(t *testing.T) {
	s := NewStore()
	short := sale()
	short.TenderMinorUnits = 1000 // below total 2000
	if err := s.Complete(short); !errors.Is(err, ErrShortTender) {
		t.Fatalf("expected ErrShortTender, got %v", err)
	}
}

func TestInvalidSale(t *testing.T) {
	s := NewStore()
	noLines := sale()
	noLines.Lines = nil
	if err := s.Complete(noLines); !errors.Is(err, ErrInvalidSale) {
		t.Fatalf("no lines accepted: %v", err)
	}
	badQty := sale()
	badQty.Lines = []LineItem{{ProductID: "P1", Quantity: 0, UnitPriceMinorUnits: 100}}
	if err := s.Complete(badQty); !errors.Is(err, ErrInvalidSale) {
		t.Fatalf("zero quantity accepted: %v", err)
	}
}

func TestTenantIsolation(t *testing.T) {
	s := NewStore()
	_ = s.Complete(sale())
	// a different tenant with the same sale id is a different sale
	other := sale()
	other.TenantID = "other"
	if err := s.Complete(other); err != nil {
		t.Fatalf("different tenant with same id should be allowed: %v", err)
	}
}

func TestOverflowingSaleRejectedBeforeRecord(t *testing.T) {
	for _, lines := range [][]LineItem{
		{{ProductID: "p", Quantity: 3, UnitPriceMinorUnits: math.MaxInt64}},
		{{ProductID: "p", Quantity: 1, UnitPriceMinorUnits: math.MaxInt64}, {ProductID: "q", Quantity: 1, UnitPriceMinorUnits: math.MaxInt64}, {ProductID: "r", Quantity: 1, UnitPriceMinorUnits: 4}},
	} {
		s := NewStore()
		r := sale()
		r.Lines = lines
		r.TenderMinorUnits = math.MaxInt64
		if e := s.Complete(r); e == nil {
			t.Error("unrepresentable total accepted")
		}
		if len(s.sales) != 0 {
			t.Error("overflow stored a sale")
		}
	}
}
func TestNegativeTenderCannotWrapToPositiveChange(t *testing.T) {
	s := NewStore()
	r := sale()
	r.Lines = []LineItem{{ProductID: "p", Quantity: 1, UnitPriceMinorUnits: 1}}
	r.TenderMinorUnits = math.MinInt64
	if e := s.Complete(r); e == nil {
		t.Error("negative tender accepted via wrap")
	}
	if len(s.sales) != 0 {
		t.Error("invalid tender stored")
	}
}
func TestDistinctSaleIdentityTuples(t *testing.T) {
	s := NewStore()
	a, b := sale(), sale()
	a.TenantID = "a\x00b"
	a.ID = "c"
	b.TenantID = "a"
	b.ID = "b\x00c"
	if e := s.Complete(a); e != nil {
		t.Fatal(e)
	}
	if e := s.Complete(b); e != nil {
		t.Fatal("distinct sale rejected", e)
	}
}
func TestCompletedSaleOwnsItsLines(t *testing.T) {
	s := NewStore()
	r := sale()
	if e := s.Complete(r); e != nil {
		t.Fatal(e)
	}
	r.Lines[0].UnitPriceMinorUnits = 9999
	// Internal stored-value inspection; this API does not expose a sale getter.
	for _, stored := range s.sales {
		if stored.Lines[0].UnitPriceMinorUnits != 500 {
			t.Fatal("caller mutated completed record")
		}
	}
}

func mustTotal(t *testing.T, s *Store, r Sale) int64 {
	t.Helper()
	n, e := s.Total(r)
	if e != nil {
		t.Fatal(e)
	}
	return n
}
func mustChange(t *testing.T, s *Store, r Sale) int64 {
	t.Helper()
	n, e := s.ChangeDue(r)
	if e != nil {
		t.Fatal(e)
	}
	return n
}
func TestCheckedArithmeticAndErrorPrecedence(t *testing.T) {
	s := NewStore()
	r := sale()
	r.Lines = []LineItem{{ProductID: "p", Quantity: 1, UnitPriceMinorUnits: math.MaxInt64}}
	r.TenderMinorUnits = math.MaxInt64
	if n, e := s.Total(r); e != nil || n != math.MaxInt64 {
		t.Fatal("max total", n, e)
	}
	if n, e := s.ChangeDue(r); e != nil || n != 0 {
		t.Fatal("max change", n, e)
	}
	if e := s.Complete(r); e != nil {
		t.Fatal(e)
	}
	r.ID = "bad"
	r.Lines[0].Quantity = 3
	if n, e := s.Total(r); n != 0 || !errors.Is(e, ErrOverflow) {
		t.Fatal("overflow error", n, e)
	}
	if n, e := s.ChangeDue(r); n != 0 || !errors.Is(e, ErrOverflow) {
		t.Fatal("change overflow", n, e)
	}
	r.Lines = append(r.Lines, LineItem{ProductID: "invalid", Quantity: 0, UnitPriceMinorUnits: 1})
	if e := s.Complete(r); !errors.Is(e, ErrInvalidSale) {
		t.Fatal("invalid-line priority", e)
	}
	r.Lines = []LineItem{{ProductID: "p", Quantity: 1, UnitPriceMinorUnits: 0}}
	r.TenderMinorUnits = 0
	if n, e := s.Total(r); e != nil || n != 0 {
		t.Fatal("zero valid line total", n, e)
	}
	if e := s.Complete(r); !errors.Is(e, ErrInvalidSale) {
		t.Fatal("zero sale", e)
	}
	r.Lines[0].UnitPriceMinorUnits = 1
	r.TenderMinorUnits = 1
	if e := s.Complete(r); e != nil {
		t.Fatal("rejection consumed ID", e)
	}
}
func TestConcurrentSaleAndDebugRead(t *testing.T) {
	s := NewStore()
	var wg sync.WaitGroup
	var count atomic.Int64
	for i := 0; i < 64; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			r := sale()
			if i%2 == 0 {
				r.TenantID = "other"
			}
			e := s.Complete(r)
			if e == nil {
				count.Add(1)
			} else if !errors.Is(e, ErrDuplicate) {
				t.Error(e)
			}
			_ = s.String()
		}(i)
	}
	wg.Wait()
	if count.Load() != 2 {
		t.Fatal("tenant/duplicate scope", count.Load())
	}
}
func FuzzSaleExactArithmetic(f *testing.F) {
	for _, v := range [][3]int64{{1, 1, 1}, {3, math.MaxInt64, math.MaxInt64}, {1, 1, math.MinInt64}, {1, math.MaxInt64, math.MaxInt64}, {0, 1, 1}, {1, -1, 1}, {1, 0, 0}} {
		f.Add(int(v[0]), v[1], v[2], []byte{})
		f.Add(int(v[0]), v[1], v[2], []byte{0, 4, 0, 4, 0, 2})
	}
	f.Fuzz(func(t *testing.T, q int, price, tender int64, extra []byte) {
		if len(extra) > 16 {
			extra = extra[:16]
		}
		r := sale()
		r.TenderMinorUnits = tender
		r.Lines = []LineItem{{ProductID: "p", Quantity: q, UnitPriceMinorUnits: price}}
		values := []int64{-1, 0, 1, 4, math.MaxInt64}
		for i := 0; i+1 < len(extra); i += 2 {
			l := LineItem{ProductID: "p", Quantity: int(extra[i]%5) - 1, UnitPriceMinorUnits: values[int(extra[i+1])%len(values)]}
			if extra[i]&128 != 0 {
				l.ProductID = ""
			}
			r.Lines = append(r.Lines, l)
		}
		var totalErr error
		for _, l := range r.Lines {
			if l.ProductID == "" || l.Quantity <= 0 || l.UnitPriceMinorUnits < 0 {
				totalErr = ErrInvalidSale
			}
		}
		exact := new(big.Int)
		if totalErr == nil {
			for _, l := range r.Lines {
				exact.Add(exact, new(big.Int).Mul(big.NewInt(int64(l.Quantity)), big.NewInt(l.UnitPriceMinorUnits)))
			}
			if !exact.IsInt64() {
				totalErr = ErrOverflow
			}
		}
		s := NewStore()
		n, e := s.Total(r)
		if !errors.Is(e, totalErr) {
			t.Fatal("total error", e, totalErr)
		}
		if totalErr != nil {
			if n != 0 {
				t.Fatal("error value")
			}
		} else if n != exact.Int64() {
			t.Fatal("exact total", n, exact)
		}
		changeErr := totalErr
		change := new(big.Int)
		if changeErr == nil {
			change.Sub(big.NewInt(tender), exact)
			if change.Sign() < 0 {
				changeErr = ErrShortTender
			}
		}
		d, e := s.ChangeDue(r)
		if !errors.Is(e, changeErr) {
			t.Fatal("change error", e, changeErr)
		}
		if changeErr != nil {
			if d != 0 {
				t.Fatal("error change")
			}
		} else if d != change.Int64() {
			t.Fatal("exact change", d, change)
		}
		completeErr := totalErr
		if completeErr == nil {
			if exact.Sign() <= 0 {
				completeErr = ErrInvalidSale
			} else {
				completeErr = changeErr
			}
		}
		if e := s.Complete(r); !errors.Is(e, completeErr) {
			t.Fatal("complete error", e, completeErr)
		}
		if completeErr != nil {
			if len(s.sales) != 0 {
				t.Fatal("invalid sale stored")
			}
			r.Lines = []LineItem{{ProductID: "p", Quantity: 1, UnitPriceMinorUnits: 1}}
			r.TenderMinorUnits = 1
			if e := s.Complete(r); e != nil {
				t.Fatal("failure consumed identity", e)
			}
		} else {
			if len(s.sales) != 1 {
				t.Fatal("sale not stored")
			}
			if e := s.Complete(r); !errors.Is(e, ErrDuplicate) {
				t.Fatal("duplicate sale", e)
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
| Go standard library | go 1.26.7 (toolchain) | POS | BSD-3-Clause | runtime | https://go.dev |

## 8. Apply order

1. Componer el backend (mismo módulo).
2. Colocar los dos archivos bajo `internal/pos/`.
3. Verificar con `go test ./... -count=1` y `go vet ./...`.
4. Rollback: eliminar `internal/pos/`.

## 9. Verification

- `go test ./internal/pos/ -count=1`: 4/4 PASS (total/vuelto, completo/duplicado, tender corto fail-closed, inválido, aislamiento tenant — el test atrapó y corrigió una fuga de tenant por clave no-scoped).
- `go test ./... -count=1` (24 paquetes): PASS.
- `go vet ./...`: exit 0.

## 10. Reconstruction evidence

Véase `reconstruction_evidence/GO_POS_CORE_2026-09-02_V202.md`.

## Migración obligatoria 0.2.0

Total(sale) y ChangeDue(sale) ahora retornan (int64,error). El caller debe comprobar
el error antes de usar el importe. ErrOverflow indica total fuera de int64;
ErrInvalidSale valida líneas y ErrShortTender indica pago inferior al total.
Los helpers no validan identidad de venta ni ejecutan cobro. Total de líneas
válidas de precio cero puede ser0; Complete mantiene exigencia de total positivo.
Sin saturación, wrap, shim con sentinel ambiguo ni nueva política monetaria.

## Revisión V361

0.2.0 corrige la frontera numérica demostrada y conserva las reglas existentes
dentro del dominio representable. AUTHORED/LicenseRef-Workspace-Owner y CONDITIONED
persisten. Suite V202 histórica no reejecutada. Evidencia nueva:
reconstruction_evidence/POS_SLO_BOUNDARIES_AND_DOCKER_SCOPE_V361.md.

V361 verificado entre ambos cores: 4 fuentes reconstruidas, 20 tests x3,
vet/build,34semillas/3381083ejecuciones fuzz PASS.11tests históricos conservan
aserciones; POS agrega comprobación de error por helpers del test y rechaza
2firmas antiguas de resultado único. Sin -race, cobro real o alertas desplegadas.

## Auditoría por claim V374

La revisión 0.2.1 se reconstruye y prueba en Go1.26.8/Windows con fuentes
oficiales fijadas y comparaciones por contrato. Ver
`reconstruction_evidence/CORE_CLAIM_ADMISSION_V374.md` desde la raíz: identidad,
licencia local, prueba por invariante, fuente semántica y dictamen de equivalencia
se mantienen separados. El lenguaje/stdlib no es fuente de un negocio empresarial.
No cambia este claim, API, permiso de composición ni las condiciones de integración.
La revisión 0.2.0 y sus bytes quedan preservados en el expediente anterior.
