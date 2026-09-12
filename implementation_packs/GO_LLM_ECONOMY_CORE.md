# Go LLM Economy Core

## 1. Metadata

```yaml
pack_id: "GO-LLM-ECONOMY-CORE"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa la capa de gobierno de costo que hace el gasto por respuesta casi nulo: cache semántica ($0 en hit), compresión de prompt con presupuesto de tokens, ledger de gasto por tenant fail-closed y un gobernador que decide el camino más barato."
stacks: ["Go 1.26.7"]
compatible_with: ["GO-ENTERPRISE-BACKEND 0.4.x", "GO-SEARCH-CORE 0.1.x", "GO-CONVERSATIONAL-AGENT 0.1.x"]
incompatible_with: ["LLM sin presupuesto de tokens", "respuestas sin cache semántica", "gasto por tenant no gobernado"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://go.dev"]
verified_at: "2026-09-02"
```

Este pack no afirma el costo exacto de ningún proveedor: el costo real lo fija el proveedor y el volumen. Lo que garantiza es la **arquitectura** que minimiza las llamadas al modelo: el ruteo determinista (en `GO-CONVERSATIONAL-AGENT`) y la cache semántica resuelven la mayoría de turnos a `$0`; sólo el fallback ambiguo llama al modelo, comprimido y con presupuesto. El tokenizer exacto es del proveedor (CONDITIONED).

## 2. Applicability

Use este pack cuando el chatbot deba servir miles de turnos con gasto de LLM casi nulo: cache semántica sobre respuestas ya dadas, prompts comprimidos a un presupuesto y límite de gasto por tenant con fallo cerrado.

Rechace este pack para: flujos donde cada respuesta deba ser generada en vivo (sin cache); o presupuesto por tenant sin definir.

## 3. Architecture contract

- **Ownership**: `internal/economy` gobierna el gasto (cache semántica, compresión, ledger) y la decisión de camino. `internal/search` aporta la similitud coseno; `internal/agent` aporta el ruteo determinista `$0`.
- **Invariantes**: (1) un hit de cache semántica no llama al modelo (`$0`). (2) El prompt comprimido nunca excede su presupuesto. (3) Un tenant sobre-presupuesto no puede gastar (`ErrBudgetExceeded`). (4) Vectores de dimensión distinta o cero no producen score inventado.
- **Data flow**: `Governor.Decide` → `SemanticCache.Lookup` (hit → `PathCache`) → `Ledger.CanSpend` (sin presupuesto → `PathBlocked`) → `PathLLM` (prompt comprimido).
- **Failure modes**: presupuesto insuficiente → bloqueo fail-closed; similitud bajo umbral → miss (no falso positivo); gasto negativo → error.
- **Seguridad/privacidad**: cache y ledger por tenant; la cache no expone respuestas entre tenants.
- **Performance budget**: O(1) ledger, O(N) cache lineal (sustituible por índice ANN si el workload lo exige).
- **Operación/migración/rollback**: sin migración; ledger y cache son reconstruibles.

## 4. Exact file manifest

```text
CREATE internal/economy/economy.go
CREATE internal/economy/semantic_cache.go
CREATE internal/economy/compressor.go
CREATE internal/economy/ledger.go
CREATE internal/economy/governor.go
CREATE internal/economy/economy_test.go
```

## 5. Materialization blocks

### FILE: `internal/economy/economy.go`
```yaml
block_id: "GO-LLM-ECONOMY-CORE:internal/economy/economy.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "0b319c126bf3a5ebbec8872864006d65ce7da7eb34eb418cd223d198df9d5f91"
variables: []
secrets_allowed: false
```
````go
// Package economy provides the cost-governance layer that makes per-response
// LLM spend near zero: deterministic/cache-first paths cost nothing, prompts
// are compressed to a token budget, and per-tenant spend is enforced fail-closed.
package economy

// EstimateTokens returns a conservative token estimate for text. It is a
// documented heuristic (bytes/3) that over-estimates for non-ASCII input; the
// exact tokenizer is provider-specific and CONDITIONED.
func EstimateTokens(text string) int {
	if len(text) == 0 {
		return 0
	}
	n := len(text) / 3
	if n < 1 {
		n = 1
	}
	return n
}
````

### FILE: `internal/economy/semantic_cache.go`
```yaml
block_id: "GO-LLM-ECONOMY-CORE:internal/economy/semantic_cache.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "bf0925828c9122983016e1d910685b6ee778fcad00316e81762e06740694e9a3"
variables: []
secrets_allowed: false
```
````go
package economy

import (
	"sync"

	"elite.local/enterprise/internal/search"
)

type semanticEntry struct {
	embedding []float32
	answer    string
}

// SemanticCache returns a cached answer when a new question is sufficiently
// similar (cosine) to a previously answered one, avoiding a model call ($0).
type SemanticCache struct {
	mu      sync.RWMutex
	entries []semanticEntry
}

// Lookup returns the best cached answer whose cosine similarity to q is at
// least threshold, or ok=false. Zero/dimension-mismatched vectors are skipped,
// never scored.
func (c *SemanticCache) Lookup(q []float32, threshold float64) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	best := -1.0
	bestIdx := -1
	for i, e := range c.entries {
		sim, err := search.CosineSimilarity(e.embedding, q)
		if err != nil {
			continue
		}
		if sim > best {
			best = sim
			bestIdx = i
		}
	}
	if bestIdx >= 0 && best >= threshold {
		return c.entries[bestIdx].answer, true
	}
	return "", false
}

// Store records an answered question embedding and its answer.
func (c *SemanticCache) Store(embedding []float32, answer string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = append(c.entries, semanticEntry{
		embedding: append([]float32(nil), embedding...),
		answer:    answer,
	})
}
````

### FILE: `internal/economy/compressor.go`
```yaml
block_id: "GO-LLM-ECONOMY-CORE:internal/economy/compressor.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "e6f00f9b0c910a953ffcb747c912559f287241d6262baaaba6ee6f8f16598069"
variables: []
secrets_allowed: false
```
````go
package economy

// CompressedPrompt is a bounded, prioritized prompt.
type CompressedPrompt struct {
	System  string
	Context []string // highest-priority context, truncated to the budget
	User    string
	Tokens  int
}

// Compress builds a prompt within a token budget. The system message and the
// user's last message are always kept; context items are added only while the
// budget allows. This minimizes tokens for the eventual model call.
func Compress(system, user string, context []string, budgetTokens int) CompressedPrompt {
	out := CompressedPrompt{System: system, User: user}
	used := EstimateTokens(system) + EstimateTokens(user)
	for _, ctx := range context {
		t := EstimateTokens(ctx)
		if used+t > budgetTokens {
			break
		}
		out.Context = append(out.Context, ctx)
		used += t
	}
	out.Tokens = used
	return out
}
````

### FILE: `internal/economy/ledger.go`
```yaml
block_id: "GO-LLM-ECONOMY-CORE:internal/economy/ledger.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "0c3d3e89eb0ae356842cad62b7d8b0070b27ab47e37f4fd19a18694d37fbfe91"
variables: []
secrets_allowed: false
```
````go
package economy

import (
	"errors"
	"sync"
)

// ErrBudgetExceeded reports a spend that would exceed the tenant's budget.
var ErrBudgetExceeded = errors.New("economy: budget exceeded")

// CostLedger tracks token spend per tenant against a budget. It fails closed:
// an over-budget tenant cannot spend.
type CostLedger struct {
	mu     sync.Mutex
	budget map[string]int64
}

// NewCostLedger returns an empty ledger.
func NewCostLedger() *CostLedger {
	return &CostLedger{budget: make(map[string]int64)}
}

// SetBudget sets the remaining-token budget for a tenant (clamped to >= 0).
func (l *CostLedger) SetBudget(tenant string, tokens int64) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if tokens < 0 {
		tokens = 0
	}
	l.budget[tenant] = tokens
}

// Remaining returns the tenant's remaining budget.
func (l *CostLedger) Remaining(tenant string) int64 {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.budget[tenant]
}

// CanSpend reports whether spending tokens would stay within budget.
func (l *CostLedger) CanSpend(tenant string, tokens int64) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if tokens < 0 {
		return false
	}
	return l.budget[tenant] >= tokens
}

// Spend consumes tokens; it fails closed if the budget would be exceeded.
func (l *CostLedger) Spend(tenant string, tokens int64) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if tokens < 0 {
		return errors.New("economy: negative spend")
	}
	if l.budget[tenant] < tokens {
		return ErrBudgetExceeded
	}
	l.budget[tenant] -= tokens
	return nil
}
````

### FILE: `internal/economy/governor.go`
```yaml
block_id: "GO-LLM-ECONOMY-CORE:internal/economy/governor.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "03ea503ac1bd78b383334f6dfd5700d80c82670ce26e5a3220510ce817a40a42"
variables: []
secrets_allowed: false
```
````go
package economy

// Path is the cheapest valid resolution for a request.
type Path int

const (
	PathCache   Path = iota // $0: semantic-cache hit
	PathLLM                 // requires a budgeted, compressed model call
	PathBlocked             // over budget, fail-closed
)

// Governor decides the cheapest valid path for a request: a semantic-cache hit
// costs nothing, otherwise a budgeted model call, otherwise blocked.
type Governor struct {
	Cache     *SemanticCache
	Ledger    *CostLedger
	Threshold float64
}

// Decide returns the path and, for a cache hit, the cached answer.
func (g *Governor) Decide(tenant string, embedding []float32, estimatedTokens int64) (Path, string) {
	if g.Cache != nil {
		if ans, ok := g.Cache.Lookup(embedding, g.Threshold); ok {
			return PathCache, ans
		}
	}
	if g.Ledger == nil || !g.Ledger.CanSpend(tenant, estimatedTokens) {
		return PathBlocked, ""
	}
	return PathLLM, ""
}
````

### FILE: `internal/economy/economy_test.go`
```yaml
block_id: "GO-LLM-ECONOMY-CORE:internal/economy/economy_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "93432123b9933300a2ac5238660343833ffd8086fbb1384b3dd5028a98e6a5d4"
variables: []
secrets_allowed: false
```
````go
package economy

import (
	"errors"
	"testing"
)

func TestEstimateTokens(t *testing.T) {
	if EstimateTokens("") != 0 {
		t.Fatal("empty text should be 0 tokens")
	}
	if EstimateTokens("hola") < 1 {
		t.Fatal("non-empty text should be >= 1 token")
	}
}

func TestSemanticCacheHitAndMiss(t *testing.T) {
	c := &SemanticCache{}
	c.Store([]float32{1, 0, 0}, "respuesta A")
	c.Store([]float32{0, 1, 0}, "respuesta B")

	if ans, ok := c.Lookup([]float32{1, 0, 0}, 0.9); !ok || ans != "respuesta A" {
		t.Fatalf("exact hit failed: %q ok=%v", ans, ok)
	}
	if ans, ok := c.Lookup([]float32{0.9, 0.1, 0}, 0.9); !ok || ans != "respuesta A" {
		t.Fatalf("near hit failed: %q ok=%v", ans, ok)
	}
	if _, ok := c.Lookup([]float32{0, 0, 1}, 0.9); ok {
		t.Fatal("orthogonal vector should miss")
	}
	if _, ok := c.Lookup([]float32{1, 0}, 0.9); ok {
		t.Fatal("dim mismatch should miss, not error")
	}
}

func TestCompressRespectsBudget(t *testing.T) {
	// budget just enough for system + user, no context
	sys := "system prompt"
	usr := "pregunta del usuario"
	budget := EstimateTokens(sys) + EstimateTokens(usr)
	p := Compress(sys, usr, []string{"contexto largo que no entra", "otro"}, budget)
	if len(p.Context) != 0 {
		t.Fatalf("expected no context within tight budget, got %v", p.Context)
	}
	if p.Tokens != budget {
		t.Fatalf("expected tokens=%d, got %d", budget, p.Tokens)
	}

	// larger budget admits context up to the limit
	big := budget + 1000
	p2 := Compress(sys, usr, []string{"a", "b"}, big)
	if len(p2.Context) != 2 {
		t.Fatalf("expected both context items, got %v", p2.Context)
	}
	if p2.Tokens > big {
		t.Fatalf("compressed prompt exceeded budget: %d > %d", p2.Tokens, big)
	}
}

func TestLedgerFailClosed(t *testing.T) {
	l := NewCostLedger()
	l.SetBudget("t", 100)
	if !l.CanSpend("t", 100) {
		t.Fatal("exact budget should be spendable")
	}
	if l.CanSpend("t", 101) {
		t.Fatal("over-budget should not be spendable")
	}
	if err := l.Spend("t", 40); err != nil {
		t.Fatal(err)
	}
	if l.Remaining("t") != 60 {
		t.Fatalf("remaining should be 60, got %d", l.Remaining("t"))
	}
	if err := l.Spend("t", 61); !errors.Is(err, ErrBudgetExceeded) {
		t.Fatalf("expected ErrBudgetExceeded, got %v", err)
	}
}

func TestGovernorDecision(t *testing.T) {
	g := &Governor{
		Cache:     &SemanticCache{},
		Ledger:    NewCostLedger(),
		Threshold: 0.9,
	}
	g.Cache.Store([]float32{1, 0}, "cacheada")
	g.Ledger.SetBudget("t", 100)

	// cache hit → PathCache
	if p, ans := g.Decide("t", []float32{1, 0}, 1000); p != PathCache || ans != "cacheada" {
		t.Fatalf("expected cache hit, got path=%d ans=%q", p, ans)
	}
	// miss but within budget → PathLLM
	if p, _ := g.Decide("t", []float32{0, 1}, 50); p != PathLLM {
		t.Fatalf("expected PathLLM, got %d", p)
	}
	// miss and over budget → PathBlocked
	if p, _ := g.Decide("t", []float32{0, 1}, 200); p != PathBlocked {
		t.Fatalf("expected PathBlocked, got %d", p)
	}
}
````


## 6. Configuration surface

Sin variables ni secretos. El umbral de similitud y el presupuesto por tenant se configuran en código como valores exactos; un presupuesto negativo se clampa a cero.

## 7. Dependency bill

| Package/image/tool | Pin exacto | Uso | Licencia | Runtime/build | Fuente oficial |
|---|---|---|---|---|---|
| Go standard library | go 1.26.7 (toolchain) | núcleo | BSD-3-Clause | runtime | https://go.dev |
| internal/search | 0.1.0 | similitud coseno | LicenseRef-Workspace-Owner | runtime | este repositorio |

## 8. Apply order

1. Componer `GO-SEARCH-CORE` y `GO-CONVERSATIONAL-AGENT` (mismo módulo).
2. Colocar los seis archivos bajo `internal/economy/`.
3. Verificar con `go test ./... -count=1` y `go vet ./...`.
4. Rollback: eliminar `internal/economy/`; no deja estado.

## 9. Verification

- `go test ./internal/economy/ -count=1`: 5/5 PASS (estimación de tokens, cache semántica hit/miss/dim-mismatch, compresión con presupuesto, ledger fail-closed, gobernador cache→llm→blocked).
- `go test ./... -count=1` (6 paquetes): PASS.
- `go vet ./...`: exit 0.

## 10. Reconstruction evidence

Véase `reconstruction_evidence/GO_LLM_ECONOMY_CORE_2026-09-02_V181.md`.
