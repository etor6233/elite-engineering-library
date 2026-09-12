# Go Gift Cards Core

## 1. Metadata

```yaml
pack_id: "GO-GIFT-CARDS-CORE"
pack_version: "0.1.3"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Referencia AUTHORED en memoria para emitir saldo positivo y canjear sin sobreconsumo, con identificador no vacío de canje único por tenant. Rechaza duplicados; no hace top-up, replay exitoso, persistencia, autorización ni checkout productivo."
stacks: ["Go 1.26.8"]
compatible_with: ["módulo aislado Go 1.26.7; integración con owners empresariales no admitida"]
incompatible_with: ["saldo negativo", "doble canje", "crédito cross-tenant"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://go.dev"]
verified_at: "2026-09-10"
```

## 2. Applicability

Referencia local para emisión/canje; entradas positivas, canje con identidad no vacía, saldo suficiente. La API no implementa moneda, checkout, recarga, reversión, persistencia ni autorización. Deben resolverse en los owners reales antes de admitir el proyecto.

## 3. Architecture contract

- **Ownership**: `internal/giftcards` gobierna el crédito; el pago con gift card es composición del dominio.
- **Invariantes**: (1) amount > 0. (2) redeem ≤ balance (fail-closed). (3) redeem ID no vacío, único dentro del tenant entre todas sus tarjetas; repetición retorna ErrSpent sin nuevo débito. (4) tenant-scoped.
- **Data flow**: `Issue` → `Redeem` → `Balance`.
- **Failure modes**: inactivo, insuficiente, duplicado, gastado → error.
- **Seguridad/privacidad**: tenant-scoped.
- **Performance budget**: O(1).

## 4. Exact file manifest

```text
CREATE internal/giftcards/giftcards.go
CREATE internal/giftcards/giftcards_test.go
```

## 5. Materialization blocks

### FILE: `internal/giftcards/giftcards.go`
```yaml
block_id: "GO-GIFT-CARDS-CORE:internal/giftcards/giftcards.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "a088ad1335038bc2f5e6658197b1aa8d3e8a9456cba5d9527523f6874e51a878"
variables: []
secrets_allowed: false
```
````go
// Package giftcards provides tenant-scoped store credit: issue/redeem with
// idempotency and fail-closed balance (never negative, never over-redeem).
package giftcards

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"
)

var (
	codeRe = regexp.MustCompile(`^[A-Z0-9][A-Z0-9._-]{0,31}$`)

	ErrInvalidCard  = errors.New("giftcards: invalid card")
	ErrInactive     = errors.New("giftcards: inactive")
	ErrInsufficient = errors.New("giftcards: insufficient balance")
	ErrDuplicate    = errors.New("giftcards: duplicate code")
	ErrNotFound     = errors.New("giftcards: not found")
	ErrSpent        = errors.New("giftcards: idempotency id already used")
)

// Card is a store-credit card.
type Card struct {
	TenantID          string
	Code              string
	BalanceMinorUnits int64
	Active            bool
}

// Store holds cards and redeem idempotency.
type Store struct {
	mu    sync.Mutex
	cards map[string]Card
	spent map[redemptionID]bool
}

// NewStore returns an empty store.
func NewStore() *Store {
	return &Store{cards: make(map[string]Card), spent: make(map[redemptionID]bool)}
}

type redemptionID struct{ tenant, id string }

func key(tenant, code string) string { return tenant + "\x00" + code }

// Issue creates a card. Duplicate codes are rejected; it does not top up.
func (s *Store) Issue(tenant, code string, amountMinorUnits int64) error {
	if !codeRe.MatchString(code) || strings.TrimSpace(tenant) == "" || amountMinorUnits <= 0 {
		return ErrInvalidCard
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.cards[key(tenant, code)]; ok {
		return ErrDuplicate
	}
	s.cards[key(tenant, code)] = Card{TenantID: tenant, Code: code, BalanceMinorUnits: amountMinorUnits, Active: true}
	return nil
}

// Redeem debits the card. A successful redeem ID is unique within its tenant
// across all cards; repeats return ErrSpent without another debit.
func (s *Store) Redeem(tenant, code, redeemID string, amountMinorUnits int64) error {
	if amountMinorUnits <= 0 || strings.TrimSpace(tenant) == "" || !codeRe.MatchString(code) || strings.TrimSpace(redeemID) == "" {
		return ErrInvalidCard
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.spent[redemptionID{tenant, redeemID}] {
		return ErrSpent
	}
	c, ok := s.cards[key(tenant, code)]
	if !ok {
		return ErrNotFound
	}
	if !c.Active {
		return ErrInactive
	}
	if c.BalanceMinorUnits < amountMinorUnits {
		return ErrInsufficient
	}
	s.spent[redemptionID{tenant, redeemID}] = true
	c.BalanceMinorUnits -= amountMinorUnits
	s.cards[key(tenant, code)] = c
	return nil
}

// Balance returns the remaining credit.
func (s *Store) Balance(tenant, code string) (int64, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	c, ok := s.cards[key(tenant, code)]
	if !ok {
		return 0, false
	}
	return c.BalanceMinorUnits, true
}

// String aids debugging.
func (s *Store) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return fmt.Sprintf("giftcards(%d)", len(s.cards))
}
````

### FILE: `internal/giftcards/giftcards_test.go`
```yaml
block_id: "GO-GIFT-CARDS-CORE:internal/giftcards/giftcards_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "3d3d247d2f3531707ccc4088e97aa2c89d11c995c6a6b9344ceece566670fe7a"
variables: []
secrets_allowed: false
```
````go
package giftcards

import (
	"errors"
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"testing"
)

func TestIssueRedeemBalance(t *testing.T) {
	s := NewStore()
	if err := s.Issue("t", "GC1", 1000); err != nil {
		t.Fatal(err)
	}
	if b, _ := s.Balance("t", "GC1"); b != 1000 {
		t.Fatalf("expected 1000, got %d", b)
	}
	if err := s.Redeem("t", "GC1", "r1", 300); err != nil {
		t.Fatal(err)
	}
	if b, _ := s.Balance("t", "GC1"); b != 700 {
		t.Fatalf("expected 700, got %d", b)
	}
}

func TestRedeemInsufficientFailsClosed(t *testing.T) {
	s := NewStore()
	_ = s.Issue("t", "GC1", 100)
	if err := s.Redeem("t", "GC1", "r1", 101); !errors.Is(err, ErrInsufficient) {
		t.Fatalf("expected ErrInsufficient, got %v", err)
	}
	if b, _ := s.Balance("t", "GC1"); b != 100 {
		t.Fatalf("balance changed on failed redeem: %d", b)
	}
}

func TestIdempotencyByRedeemID(t *testing.T) {
	s := NewStore()
	_ = s.Issue("t", "GC1", 1000)
	_ = s.Redeem("t", "GC1", "r1", 100)
	if err := s.Redeem("t", "GC1", "r1", 100); !errors.Is(err, ErrSpent) {
		t.Fatalf("expected ErrSpent, got %v", err)
	}
	if b, _ := s.Balance("t", "GC1"); b != 900 {
		t.Fatalf("double spend: %d", b)
	}
}

func TestDuplicateAndInvalid(t *testing.T) {
	s := NewStore()
	_ = s.Issue("t", "GC1", 100)
	if err := s.Issue("t", "GC1", 100); !errors.Is(err, ErrDuplicate) {
		t.Fatalf("duplicate accepted: %v", err)
	}
	if err := s.Issue("t", "bad code!", 100); !errors.Is(err, ErrInvalidCard) {
		t.Fatalf("bad code accepted: %v", err)
	}
	if err := s.Issue("t", "GC2", 0); !errors.Is(err, ErrInvalidCard) {
		t.Fatalf("zero amount accepted: %v", err)
	}
}

func TestTenantIsolation(t *testing.T) {
	s := NewStore()
	_ = s.Issue("t", "GC1", 100)
	if _, ok := s.Balance("other", "GC1"); ok {
		t.Fatal("cross-tenant leak")
	}
	if err := s.Redeem("other", "GC1", "r1", 50); !errors.Is(err, ErrNotFound) {
		t.Fatalf("cross-tenant should be not found, got %v", err)
	}
}

func TestRedeemIdentityScopedByTenant(t *testing.T) {
	s := NewStore()
	for _, tenant := range []string{"a", "b"} {
		if e := s.Issue(tenant, "CARD", 100); e != nil {
			t.Fatal(e)
		}
		if e := s.Redeem(tenant, "CARD", "shared", 10); e != nil {
			t.Fatalf("tenant %s cannot redeem: %v", tenant, e)
		}
	}
	if e := s.Issue("a", "OTHER", 100); e != nil {
		t.Fatal(e)
	}
	if e := s.Redeem("a", "OTHER", "shared", 10); !errors.Is(e, ErrSpent) {
		t.Fatalf("same-tenant duplicate: %v", e)
	}
	if b, _ := s.Balance("a", "OTHER"); b != 100 {
		t.Fatal("duplicate debited other card")
	}
}
func TestBlankRedeemIdentityRejected(t *testing.T) {
	for _, id := range []string{"", " ", "\t\n"} {
		t.Run(id, func(t *testing.T) {
			s := NewStore()
			if e := s.Issue("t", "CARD", 100); e != nil {
				t.Fatal(e)
			}
			if e := s.Redeem("t", "CARD", id, 1); !errors.Is(e, ErrInvalidCard) {
				t.Fatalf("blank identity accepted: %v", e)
			}
			if b, _ := s.Balance("t", "CARD"); b != 100 {
				t.Fatal("invalid identity mutated balance")
			}
		})
	}
}

func TestFailedRedemptionCanRetryAndConcurrentDebit(t *testing.T) {
	s := NewStore()
	if e := s.Issue("t", "CARD", 10); e != nil {
		t.Fatal(e)
	}
	if e := s.Redeem("t", "CARD", "retry", 11); !errors.Is(e, ErrInsufficient) {
		t.Fatal(e)
	}
	if e := s.Redeem("t", "CARD", "retry", 1); e != nil {
		t.Fatalf("failed attempt consumed identity: %v", e)
	}
	var wg sync.WaitGroup
	var successes atomic.Int64
	for i := 0; i < 64; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if e := s.Issue("t", fmt.Sprintf("CARD-%d", i), 1); e != nil {
				t.Error(e)
			}
			_ = s.String()
			e := s.Redeem("t", "CARD", fmt.Sprintf("debit-%d", i), 1)
			if e == nil {
				successes.Add(1)
			} else if !errors.Is(e, ErrInsufficient) {
				t.Error(e)
			}
		}(i)
	}
	wg.Wait()
	if b, _ := s.Balance("t", "CARD"); b != 0 || successes.Load() != 9 {
		t.Fatal("concurrent overspend")
	}
}
func TestRedemptionTupleDoesNotAlias(t *testing.T) {
	s := NewStore()
	for _, tn := range []string{"a\x00b", "a"} {
		if e := s.Issue(tn, "CARD", 10); e != nil {
			t.Fatal(e)
		}
	}
	if e := s.Redeem("a\x00b", "CARD", "c", 1); e != nil {
		t.Fatal(e)
	}
	if e := s.Redeem("a", "CARD", "b\x00c", 1); e != nil {
		t.Fatalf("delimiter tuple alias: %v", e)
	}
}
func FuzzRedemptionStateMachine(f *testing.F) {
	for _, seed := range []uint64{0, 1, 2, 1<<63 - 1, 1 << 63, ^uint64(0)} {
		f.Add(seed, []byte{0, 2, 1, 2, 0, 2, 5, 3, 2, 4, 3, 1, 4, 5})
	}
	f.Fuzz(func(t *testing.T, seed uint64, ops []byte) {
		if len(ops) > 96 {
			ops = ops[:96]
		}
		credit := int64(seed & math.MaxInt64)
		if credit == 0 {
			credit = 1
		}
		tenants := []string{"a", "b", "a\x00b"}
		codes := []string{"CARD", "OTHER"}
		ids := []string{"shared", "next", "", " ", "c", "b\x00c"}
		amounts := []int64{0, -1, 1, 2, math.MaxInt64, int64(seed)}
		model := map[string]map[string]int64{}
		used := map[string]map[string]bool{}
		s := NewStore()
		for _, tn := range tenants {
			model[tn] = map[string]int64{}
			used[tn] = map[string]bool{}
			for _, code := range codes {
				if e := s.Issue(tn, code, credit); e != nil {
					t.Fatal(e)
				}
				model[tn][code] = credit
			}
		}
		for i := 0; i+1 < len(ops); i += 2 {
			b := ops[i]
			tn, code := tenants[int(b)%len(tenants)], codes[int(b>>2)%2]
			id := ids[int(b>>3)%len(ids)]
			amount := amounts[int(ops[i+1])%len(amounts)]
			expected := error(nil)
			if amount <= 0 || id == "" || id == " " {
				expected = ErrInvalidCard
			} else if used[tn][id] {
				expected = ErrSpent
			} else if amount > model[tn][code] {
				expected = ErrInsufficient
			}
			if e := s.Redeem(tn, code, id, amount); !errors.Is(e, expected) {
				t.Fatalf("classification got=%v want=%v", e, expected)
			}
			if expected == nil {
				model[tn][code] -= amount
				used[tn][id] = true
			}
			for _, tenant := range tenants {
				for _, card := range codes {
					if b, ok := s.Balance(tenant, card); !ok || b != model[tenant][card] {
						t.Fatal("balance mutation mismatch")
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
| Go standard library | go 1.26.7 (toolchain) | gift cards | BSD-3-Clause | runtime | https://go.dev |

## 8. Apply order

1. Componer el backend (mismo módulo).
2. Colocar los dos archivos bajo `internal/giftcards/`.
3. Verificar con `go test ./... -count=1` y `go vet ./...`.
4. Rollback: eliminar `internal/giftcards/`.

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
Antes: reconstruction_evidence/GO_GIFT_CARDS_CORE_2026-09-02_V198.md.

AUTHORED/LicenseRef-Workspace-Owner, SUPPORTED_REFERENCE y CONDITIONED persisten.
No seleccionado por el perfil integral. Las compatibilidades históricas declaradas
no demuestran integración con las versiones actuales de los owners empresariales.

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
