# Go Rate Limit Core

## 1. Metadata

```yaml
pack_id: "GO-RATE-LIMIT-CORE"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa el rate limiting / anti-abuse tenant-scoped con token bucket: deny fail-closed cuando se agotan los tokens, refill a tasa fija y burst acotado."
stacks: ["Go 1.26.7"]
compatible_with: ["GO-ENTERPRISE-BACKEND-CORE 0.1.x", "GO_ELECTROMOBILITY_PUBLIC_CRM_API 0.2.x"]
incompatible_with: ["tasa o burst no positivos", "tenant/clave vacía"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://go.dev", "https://sre.google"]
verified_at: "2026-09-02"
```

## 2. Applicability

Use ante cada endpoint público para limitar por clave (anti-abuse). Rechace para tasa/burst no positivos o tenant/clave vacía.

## 3. Architecture contract

- **Ownership**: `internal/ratelimit` decide allow/deny; el edge (CDN/WAF) replica la política distribuida.
- **Invariantes**: (1) burst acotado. (2) refill a tasa fija. (3) deny fail-closed (entrada inválida o tokens agotados). (4) buckets por tenant+clave.
- **Data flow**: `NewLimiter(Policy)` → `Allow/AllowN`.
- **Failure modes**: política inválida, tokens agotados → deny.
- **Seguridad/privacidad**: tenant-scoped; sin almacenar PII (sólo contadores).
- **Performance budget**: O(1) por allow.

## 4. Exact file manifest

```text
CREATE internal/ratelimit/ratelimit.go
CREATE internal/ratelimit/ratelimit_test.go
```

## 5. Materialization blocks

### FILE: `internal/ratelimit/ratelimit.go`
```yaml
block_id: "GO-RATE-LIMIT-CORE:internal/ratelimit/ratelimit.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "9d94c564e4f0d291de05289e6afbc3168bda31dedfa9f47c1cf16c72e971477f"
variables: []
secrets_allowed: false
```
````go
// Package ratelimit provides a tenant-scoped token-bucket rate limiter:
// per-key buckets, deny (fail-closed) when tokens are exhausted, refilled at
// a fixed rate with a bounded burst. Stdlib-only. AUTHORED.
package ratelimit

import (
	"errors"
	"strings"
	"sync"
	"time"
)

var (
	ErrInvalidPolicy = errors.New("ratelimit: rate and burst must be positive")
)

// Policy is a single refill policy applied to every key.
type Policy struct {
	Rate  float64       // tokens per second
	Burst float64       // maximum bucket capacity
}

// Limiter is the tenant-scoped limiter.
type Limiter struct {
	mu     sync.Mutex
	policy Policy
	now    func() time.Time
	items  map[string]*bucket
}

type bucket struct {
	tokens float64
	last   time.Time
}

// NewLimiter returns a limiter with the given policy and a real clock.
func NewLimiter(p Policy) (*Limiter, error) {
	if p.Rate <= 0 || p.Burst <= 0 {
		return nil, ErrInvalidPolicy
	}
	return &Limiter{policy: p, now: time.Now, items: make(map[string]*bucket)}, nil
}

// setClock is used by tests to control time.
func (l *Limiter) setClock(fn func() time.Time) { l.now = fn }

func (l *Limiter) key(tenant, id string) string { return tenant + "\x00" + id }

func (l *Limiter) refill(b *bucket, now time.Time) {
	elapsed := now.Sub(b.last).Seconds()
	if elapsed < 0 {
		elapsed = 0
	}
	b.tokens += elapsed * l.policy.Rate
	if b.tokens > l.policy.Burst {
		b.tokens = l.policy.Burst
	}
	b.last = now
}

// AllowN consumes n tokens for the key if available; deny (false) otherwise
// and on any invalid input (fail-closed).
func (l *Limiter) AllowN(tenant, id string, n float64) bool {
	if strings.TrimSpace(tenant) == "" || strings.TrimSpace(id) == "" || n <= 0 {
		return false
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	now := l.now()
	k := l.key(tenant, id)
	b, ok := l.items[k]
	if !ok {
		b = &bucket{tokens: l.policy.Burst, last: now}
		l.items[k] = b
	}
	l.refill(b, now)
	if b.tokens < n {
		return false
	}
	b.tokens -= n
	return true
}

// Allow consumes one token.
func (l *Limiter) Allow(tenant, id string) bool { return l.AllowN(tenant, id, 1) }

// String aids debugging.
func (l *Limiter) String() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return "ratelimit(" + itoa(len(l.items)) + ")"
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var b [20]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	return string(b[i:])
}
````

### FILE: `internal/ratelimit/ratelimit_test.go`
```yaml
block_id: "GO-RATE-LIMIT-CORE:internal/ratelimit/ratelimit_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "eaeffab126de36dc0ffebd94b49ddf654103b14b83d9553711275709f04a13ce"
variables: []
secrets_allowed: false
```
````go
package ratelimit

import (
	"testing"
	"time"
)

func newTestLimiter(rate, burst float64, start time.Time) *Limiter {
	l, err := NewLimiter(Policy{Rate: rate, Burst: burst})
	if err != nil {
		panic(err)
	}
	l.setClock(func() time.Time { return start })
	return l
}

func TestBurstLimited(t *testing.T) {
	start := time.Unix(0, 0)
	l := newTestLimiter(1, 3, start)
	// 3 allowed, 4th denied.
	for i := 0; i < 3; i++ {
		if !l.Allow("t", "k") {
			t.Fatalf("attempt %d should be allowed", i)
		}
	}
	if l.Allow("t", "k") {
		t.Fatal("burst exhausted but allowed")
	}
}

func TestRefillOverTime(t *testing.T) {
	start := time.Unix(0, 0)
	l := newTestLimiter(1, 1, start)
	if !l.Allow("t", "k") {
		t.Fatal("first allow failed")
	}
	if l.Allow("t", "k") {
		t.Fatal("should be empty")
	}
	// Advance 1 second → 1 token refilled.
	l.setClock(func() time.Time { return start.Add(time.Second) })
	if !l.Allow("t", "k") {
		t.Fatal("refill did not happen")
	}
}

func TestDenyInvalidInput(t *testing.T) {
	l, _ := NewLimiter(Policy{Rate: 1, Burst: 1})
	if l.Allow("", "k") {
		t.Fatal("empty tenant allowed")
	}
	if l.Allow("t", "") {
		t.Fatal("empty key allowed")
	}
	if l.AllowN("t", "k", 0) {
		t.Fatal("zero n allowed")
	}
}

func TestTenantIsolation(t *testing.T) {
	start := time.Unix(0, 0)
	l := newTestLimiter(1, 1, start)
	if !l.Allow("t1", "k") {
		t.Fatal("t1 should be allowed")
	}
	if !l.Allow("t2", "k") {
		t.Fatal("t2 should have its own bucket")
	}
}

func TestInvalidPolicy(t *testing.T) {
	if _, err := NewLimiter(Policy{Rate: 0, Burst: 1}); err != ErrInvalidPolicy {
		t.Fatalf("zero rate accepted: %v", err)
	}
	if _, err := NewLimiter(Policy{Rate: 1, Burst: 0}); err != ErrInvalidPolicy {
		t.Fatalf("zero burst accepted: %v", err)
	}
}

func TestAllowNOverBurst(t *testing.T) {
	start := time.Unix(0, 0)
	l := newTestLimiter(1, 2, start)
	if l.AllowN("t", "k", 3) {
		t.Fatal("n > burst should be denied")
	}
	if !l.AllowN("t", "k", 2) {
		t.Fatal("n == burst should be allowed")
	}
}
````


## 6. Configuration surface

Sin variables ni secretos (tasa y burst se pasan por `Policy`).

## 7. Dependency bill

| Package/image/tool | Pin exacto | Uso | Licencia | Runtime/build | Fuente oficial |
|---|---|---|---|---|---|
| Go standard library | go 1.26.7 (toolchain) | rate limiting | BSD-3-Clause | runtime | https://go.dev |

## 8. Apply order

1. Componer el backend (mismo módulo).
2. Colocar los dos archivos bajo `internal/ratelimit/`.
3. Verificar con `go test ./... -count=1` y `go vet ./...`.
4. Rollback: eliminar `internal/ratelimit/`.

## 9. Verification

- `go test ./internal/ratelimit/ -count=1`: 6/6 PASS.
- `go test ./... -count=1` (41 paquetes): PASS.
- `go vet ./...`: exit 0.

## 10. Reconstruction evidence

Véase `reconstruction_evidence/GO_RATE_LIMIT_CORE_2026-09-02_V218.md`.
