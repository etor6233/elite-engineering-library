# Go Loyalty Core

## 1. Metadata

```yaml
pack_id: "GO-LOYALTY-CORE"
pack_version: "0.1.2"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Referencia AUTHORED en memoria de puntos int64: earn/burn con saldo no negativo ni overflow, cuentas por tupla tenant/cliente y movimientos únicos por tenant. Historial devuelto como copia; no ledger durable, autorización ni política de recompensas."
stacks: ["Go 1.26.8"]
compatible_with: ["GO-ENTERPRISE-BACKEND 0.4.x", "GO-FINOPS-CORE 0.1.x"]
incompatible_with: ["saldo negativo", "doble gasto del mismo entry", "puntos cross-tenant"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://go.dev"]
verified_at: "2026-09-10"
```

Referencia `AUTHORED` en memoria; no es código de un proveedor financiero. El caller debe aportar identidad autorizada y reglas de puntos aprobadas; el pack no las infiere.

## 2. Applicability

Referencia local de movimientos de puntos. La integración compra→earn, permisos, reglas, vigencia, reversión y persistencia siguen pendientes del proyecto; estos tests no las admiten.

## 3. Architecture contract

- **Ownership**: `internal/loyalty` gobierna el ledger; el pedido (earn en compra) es del dominio.
- **Invariantes**: (1) points > 0. (2) burn ≤ balance (fail-closed). (3) entry ID único dentro del tenant entre earn/burn y sus clientes; duplicado retorna ErrDuplicate sin efecto. (4) cuentas separadas por tupla, historial retornado como copia. (5) crédito que excede int64 retorna ErrOverflow sin consumir ID ni mutar historia/saldo.
- **Data flow**: `Earn`/`Burn` → balance + history.
- **Failure modes**: puntos inválidos, saldo insuficiente, duplicado → error.
- **Seguridad/privacidad**: tenant-scoped; sin PII extra.
- **Performance budget**: mapas O(1) esperado; History recorre O(total de movimientos en memoria). No benchmark productivo.
- **Operación/migración/rollback**: sólo memoria del proceso; no hay journal durable ni API de replay/restauración. No afirmar recovery desde History por estos tests.

## 4. Exact file manifest

```text
CREATE internal/loyalty/loyalty.go
CREATE internal/loyalty/loyalty_test.go
```

## 5. Materialization blocks

### FILE: `internal/loyalty/loyalty.go`
```yaml
block_id: "GO-LOYALTY-CORE:internal/loyalty/loyalty.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "14a339db71422e0c247a74dcd73998a62c1c7e1da4549acb0e3f0ada317cabf1"
variables: []
secrets_allowed: false
```
````go
// Package loyalty provides a tenant-scoped points ledger: earn/burn with
// idempotent entries, immutable history and fail-closed balance (never
// negative, never over-burn). AUTHORED over the standard loyalty-ledger
// pattern, connected to identity/tenant by the caller.
package loyalty

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strings"
	"sync"
	"time"
)

// EntryKind is earn or burn.
type EntryKind string

const (
	KindEarn EntryKind = "earn"
	KindBurn EntryKind = "burn"
)

var (
	idRe = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:-]{0,127}$`)

	ErrInvalidPoints = errors.New("loyalty: invalid points")
	ErrInsufficient  = errors.New("loyalty: insufficient balance")
	ErrDuplicate     = errors.New("loyalty: duplicate entry")
	ErrOverflow      = errors.New("loyalty: balance exceeds int64 capacity")
)

// Entry is an immutable ledger movement.
type Entry struct {
	TenantID   string
	CustomerID string
	ID         string
	Kind       EntryKind
	Points     int64 // positive amount
	Reason     string
	At         time.Time
}

// Ledger is the points ledger, safe for concurrent use.
type Ledger struct {
	mu       sync.Mutex
	balances map[accountID]int64
	entries  map[entryID]bool
	history  []Entry
}

// NewLedger returns an empty ledger.
func NewLedger() *Ledger {
	return &Ledger{
		balances: make(map[accountID]int64),
		entries:  make(map[entryID]bool),
	}
}

type accountID struct{ tenant, customer string }
type entryID struct{ tenant, id string }

func accountKey(tenant, customer string) accountID { return accountID{tenant, customer} }

func validate(tenant, customer, id string, points int64) error {
	if strings.TrimSpace(tenant) == "" || strings.TrimSpace(customer) == "" {
		return ErrInvalidPoints
	}
	if !idRe.MatchString(id) {
		return fmt.Errorf("loyalty: invalid entry id %q", id)
	}
	if points <= 0 {
		return ErrInvalidPoints
	}
	return nil
}

// Earn credits positive points within int64 capacity. Duplicate IDs are
// rejected within a tenant; rejected attempts consume neither balance nor ID.
func (l *Ledger) Earn(tenant, customer, id string, points int64, reason string) error {
	if err := validate(tenant, customer, id, points); err != nil {
		return err
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.entries[entryID{tenant, id}] {
		return ErrDuplicate
	}
	if points > math.MaxInt64-l.balances[accountKey(tenant, customer)] {
		return ErrOverflow
	}
	l.entries[entryID{tenant, id}] = true
	l.balances[accountKey(tenant, customer)] += points
	l.history = append(l.history, Entry{
		TenantID: tenant, CustomerID: customer, ID: id, Kind: KindEarn,
		Points: points, Reason: reason, At: time.Now().UTC(),
	})
	return nil
}

// Burn debits points, fail-closed if insufficient. Successful entry IDs are
// unique within a tenant across both Earn and Burn; duplicates return ErrDuplicate.
func (l *Ledger) Burn(tenant, customer, id string, points int64, reason string) error {
	if err := validate(tenant, customer, id, points); err != nil {
		return err
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.entries[entryID{tenant, id}] {
		return ErrDuplicate
	}
	if l.balances[accountKey(tenant, customer)] < points {
		return ErrInsufficient
	}
	l.entries[entryID{tenant, id}] = true
	l.balances[accountKey(tenant, customer)] -= points
	l.history = append(l.history, Entry{
		TenantID: tenant, CustomerID: customer, ID: id, Kind: KindBurn,
		Points: points, Reason: reason, At: time.Now().UTC(),
	})
	return nil
}

// Balance returns the tenant-scoped points balance.
func (l *Ledger) Balance(tenant, customer string) int64 {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.balances[accountKey(tenant, customer)]
}

// History returns the immutable movement history for a customer.
func (l *Ledger) History(tenant, customer string) []Entry {
	l.mu.Lock()
	defer l.mu.Unlock()
	var out []Entry
	for _, e := range l.history {
		if e.TenantID == tenant && e.CustomerID == customer {
			out = append(out, e)
		}
	}
	return out
}
````

### FILE: `internal/loyalty/loyalty_test.go`
```yaml
block_id: "GO-LOYALTY-CORE:internal/loyalty/loyalty_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "105a99ffc6c9a12fb8ec0ae7c1627aec807ea68ad1eb2d125686a8081274d3a1"
variables: []
secrets_allowed: false
```
````go
package loyalty

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"sync"
	"sync/atomic"
	"testing"
)

func TestEarnBurnBalance(t *testing.T) {
	l := NewLedger()
	if err := l.Earn("t", "c1", "e1", 100, "compra"); err != nil {
		t.Fatal(err)
	}
	if err := l.Earn("t", "c1", "e2", 50, "compra"); err != nil {
		t.Fatal(err)
	}
	if l.Balance("t", "c1") != 150 {
		t.Fatalf("expected 150, got %d", l.Balance("t", "c1"))
	}
	if err := l.Burn("t", "c1", "b1", 60, "canje"); err != nil {
		t.Fatal(err)
	}
	if l.Balance("t", "c1") != 90 {
		t.Fatalf("expected 90, got %d", l.Balance("t", "c1"))
	}
}

func TestBurnInsufficientFailsClosed(t *testing.T) {
	l := NewLedger()
	_ = l.Earn("t", "c1", "e1", 50, "x")
	if err := l.Burn("t", "c1", "b1", 51, "y"); !errors.Is(err, ErrInsufficient) {
		t.Fatalf("expected ErrInsufficient, got %v", err)
	}
	if l.Balance("t", "c1") != 50 {
		t.Fatalf("balance changed on failed burn: %d", l.Balance("t", "c1"))
	}
}

func TestDuplicateEntryRejected(t *testing.T) {
	l := NewLedger()
	_ = l.Earn("t", "c1", "e1", 10, "x")
	if err := l.Earn("t", "c1", "e1", 10, "x"); !errors.Is(err, ErrDuplicate) {
		t.Fatalf("expected ErrDuplicate, got %v", err)
	}
}

func TestInvalidPointsRejected(t *testing.T) {
	l := NewLedger()
	if err := l.Earn("t", "c1", "e1", 0, "x"); !errors.Is(err, ErrInvalidPoints) {
		t.Fatalf("zero points accepted: %v", err)
	}
	if err := l.Earn("", "c1", "e1", 10, "x"); !errors.Is(err, ErrInvalidPoints) {
		t.Fatalf("empty tenant accepted: %v", err)
	}
}

func TestTenantIsolation(t *testing.T) {
	l := NewLedger()
	_ = l.Earn("tenant-a", "c1", "e1", 100, "x")
	if l.Balance("tenant-b", "c1") != 0 {
		t.Fatalf("cross-tenant leak: %d", l.Balance("tenant-b", "c1"))
	}
	if len(l.History("tenant-b", "c1")) != 0 {
		t.Fatal("cross-tenant history leak")
	}
}

func TestHistoryImmutableOrder(t *testing.T) {
	l := NewLedger()
	_ = l.Earn("t", "c1", "e1", 100, "a")
	_ = l.Burn("t", "c1", "b1", 30, "b")
	h := l.History("t", "c1")
	if len(h) != 2 || h[0].Kind != KindEarn || h[1].Kind != KindBurn {
		t.Fatalf("unexpected history: %+v", h)
	}
}

func TestCreditOverflowAtomicAndRetry(t *testing.T) {
	l := NewLedger()
	if e := l.Earn("t", "c", "max", math.MaxInt64, "initial"); e != nil {
		t.Fatal(e)
	}
	if e := l.Earn("t", "c", "next", 1, "overflow"); e == nil {
		t.Fatal("overflow accepted")
	}
	if l.Balance("t", "c") != math.MaxInt64 || len(l.History("t", "c")) != 1 {
		t.Fatal("rejection mutated balance/history")
	}
	if e := l.Burn("t", "c", "space", 1, "release"); e != nil {
		t.Fatal(e)
	}
	if e := l.Earn("t", "c", "next", 1, "retry"); e != nil {
		t.Fatalf("failed attempt consumed ID: %v", e)
	}
	if e := l.Earn("t", "c", "next", 1, "replay"); !errors.Is(e, ErrDuplicate) {
		t.Fatalf("duplicate accepted: %v", e)
	}
}
func TestAccountTupleIsolation(t *testing.T) {
	l := NewLedger()
	if e := l.Earn("a\x00b", "c", "first", 20, "initial"); e != nil {
		t.Fatal(e)
	}
	if l.Balance("a", "b\x00c") != 0 {
		t.Fatal("distinct account tuples alias")
	}
	if e := l.Burn("a", "b\x00c", "steal", 1, "other"); !errors.Is(e, ErrInsufficient) {
		t.Fatalf("other account debit: %v", e)
	}
	if l.Balance("a\x00b", "c") != 20 || len(l.History("a", "b\x00c")) != 0 {
		t.Fatal("cross-account mutation")
	}
}
func TestEntryIdentityScopedByTenant(t *testing.T) {
	l := NewLedger()
	for _, tenant := range []string{"a", "b"} {
		if e := l.Earn(tenant, "c", "shared-earn", 10, "initial"); e != nil {
			t.Fatalf("tenant %s cannot earn: %v", tenant, e)
		}
		if e := l.Burn(tenant, "c", "shared-burn", 1, "burn"); e != nil {
			t.Fatal(e)
		}
	}
	if e := l.Earn("a", "other", "shared-earn", 10, "reuse"); !errors.Is(e, ErrDuplicate) {
		t.Fatalf("same-tenant duplicate: %v", e)
	}
}

func TestConcurrentDuplicateCreditAndHistoryCopy(t *testing.T) {
	l := NewLedger()
	var wg sync.WaitGroup
	var successes atomic.Int64
	for i := 0; i < 64; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			e := l.Earn("t", "c", "same", 10, "reason")
			if e == nil {
				successes.Add(1)
			} else if !errors.Is(e, ErrDuplicate) {
				t.Error(e)
			}
		}()
	}
	wg.Wait()
	if successes.Load() != 1 || l.Balance("t", "c") != 10 {
		t.Fatal("duplicate credit under concurrency")
	}
	h := l.History("t", "c")
	if len(h) != 1 {
		t.Fatal("duplicate history")
	}
	h[0].Points = 999
	h[0].CustomerID = "changed"
	if fresh := l.History("t", "c"); len(fresh) != 1 || fresh[0].Points != 10 || fresh[0].CustomerID != "c" {
		t.Fatal("history exposed mutable state")
	}
	successes.Store(0)
	for i := 0; i < 64; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			e := l.Burn("t", "c", fmt.Sprintf("burn-%d", i), 1, "reason")
			if e == nil {
				successes.Add(1)
			} else if !errors.Is(e, ErrInsufficient) {
				t.Error(e)
			}
		}(i)
	}
	wg.Wait()
	if successes.Load() != 10 || l.Balance("t", "c") != 0 || len(l.History("t", "c")) != 11 {
		t.Fatal("concurrent over-burn/history mismatch")
	}
}

// This oracle uses arbitrary-precision balances and nested identity maps, not
// the production key representation or capacity arithmetic.
func FuzzLedgerStateMachine(f *testing.F) {
	for _, seed := range []uint64{0, 1, 2, 1<<63 - 1, 1 << 63, ^uint64(0)} {
		f.Add(seed, []byte{0, 4, 1, 0, 0, 4, 0, 4, 1, 1, 3, 4, 2, 3, 2, 1, 255, 2})
	}
	f.Fuzz(func(t *testing.T, seed uint64, ops []byte) {
		if len(ops) > 96 {
			ops = ops[:96]
		}
		tenants := []string{"a", "b", "a\x00b", "a"}
		customers := []string{"c", "b\x00c", "other", "c"}
		amounts := []int64{0, -1, 1, 2, math.MaxInt64, math.MaxInt64 - 1, int64(seed)}
		balances := map[string]map[string]*big.Int{}
		used := map[string]map[string]bool{}
		history := []Entry{}
		l := NewLedger()
		balance := func(tenant, customer string) *big.Int {
			if balances[tenant] == nil {
				balances[tenant] = map[string]*big.Int{}
			}
			if balances[tenant][customer] == nil {
				balances[tenant][customer] = new(big.Int)
			}
			return balances[tenant][customer]
		}
		for i := 0; i+2 < len(ops); i += 3 {
			b := ops[i]
			tenant, customer := tenants[int(b)%4], customers[int(b>>2)%4]
			id := fmt.Sprintf("id-%d", b>>4)
			points := amounts[int(ops[i+1])%len(amounts)]
			earn := ops[i+2]&1 == 0
			if used[tenant] == nil {
				used[tenant] = map[string]bool{}
			}
			expected := error(nil)
			next := new(big.Int).Set(balance(tenant, customer))
			kind := KindBurn
			if points <= 0 {
				expected = ErrInvalidPoints
			} else if used[tenant][id] {
				expected = ErrDuplicate
			} else if earn {
				kind = KindEarn
				next.Add(next, big.NewInt(points))
				if !next.IsInt64() {
					expected = ErrOverflow
				}
			} else {
				next.Sub(next, big.NewInt(points))
				if next.Sign() < 0 {
					expected = ErrInsufficient
				}
			}
			var got error
			if earn {
				got = l.Earn(tenant, customer, id, points, "fuzz")
			} else {
				got = l.Burn(tenant, customer, id, points, "fuzz")
			}
			if !errors.Is(got, expected) {
				t.Fatalf("error classification got=%v want=%v", got, expected)
			}
			if expected == nil {
				balance(tenant, customer).Set(next)
				used[tenant][id] = true
				history = append(history, Entry{TenantID: tenant, CustomerID: customer, ID: id, Kind: kind, Points: points, Reason: "fuzz"})
			}
			for _, tn := range tenants {
				for _, cu := range customers {
					if got := l.Balance(tn, cu); got != balance(tn, cu).Int64() {
						t.Fatalf("balance mismatch tenant=%q customer=%q", tn, cu)
					}
					actual := l.History(tn, cu)
					j := 0
					for _, e := range history {
						if e.TenantID == tn && e.CustomerID == cu {
							if j >= len(actual) {
								t.Fatal("history missing")
							}
							g := actual[j]
							if g.ID != e.ID || g.Kind != e.Kind || g.Points != e.Points || g.Reason != e.Reason || g.TenantID != tn || g.CustomerID != cu || g.At.IsZero() {
								t.Fatal("history differs")
							}
							j++
						}
					}
					if j != len(actual) {
						t.Fatal("history extra after rejected operation")
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
| Go standard library | go 1.26.7 (toolchain) | ledger | BSD-3-Clause | runtime | https://go.dev |

## 8. Apply order

1. Componer el backend (mismo módulo).
2. Colocar los dos archivos bajo `internal/loyalty/`.
3. Verificar con `go test ./... -count=1` y `go vet ./...`.
4. Rollback: eliminar `internal/loyalty/`.

## 9. Verification

V354 conserva los tests históricos: reconstrucción canónica4/4fuentes entre
ambos packs;19tests ordinarios (giftcards9/loyalty10) x3, vet/build PASS.
Dos gates fuzz de ambos targets:10s por target,6semillas por target,4workers,
mismas invariantes y fuentes PASS. La primera ejecución24workers terminó
con deadline en giftcards; se conserva como fallo de calificación, sin causa
raíz demostrada ni falsa reparación del runtime. Ver recibos V354.
No se declara detector de races, persistencia, autorización ni integración financiera.
Las composiciones completas V194/V198 son históricas; este delta no las reejecuta.

## 10. Reconstruction evidence

V354: reconstruction_evidence/CREDIT_CORE_ISOLATION_V354.md.
Antes: reconstruction_evidence/GO_LOYALTY_CORE_2026-09-02_V194.md.

AUTHORED/LicenseRef-Workspace-Owner, SUPPORTED_REFERENCE y CONDITIONED persisten.
No seleccionado por el perfil integral. Las compatibilidades históricas declaradas
no demuestran integración con las versiones actuales de los owners empresariales.

## Auditoría por claim V374

La revisión 0.1.2 se reconstruye y prueba en Go1.26.8/Windows con fuentes
oficiales fijadas y comparaciones por contrato. Ver
`reconstruction_evidence/CORE_CLAIM_ADMISSION_V374.md` desde la raíz: identidad,
licencia local, prueba por invariante, fuente semántica y dictamen de equivalencia
se mantienen separados. El lenguaje/stdlib no es fuente de un negocio empresarial.
No cambia este claim, API, permiso de composición ni las condiciones de integración.
La revisión 0.1.1 y sus bytes quedan preservados en el expediente anterior.
