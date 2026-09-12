# Go Resilience Core

## 1. Metadata

```yaml
pack_id: "GO-RESILIENCE-CORE"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa la resiliencia: circuit breaker (closed→open→half-open con sonda única), timeout y retry budget acotado con backoff inyectable."
stacks: ["Go 1.26.7"]
compatible_with: ["GO-RELIABLE-ASYNC-WORKERS 0.1.x"]
incompatible_with: ["max failures < 1", "open timeout <= 0", "retry sin tope"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://go.dev", "https://sre.google"]
verified_at: "2026-09-02"
```

## 2. Applicability

Use para proteger llamadas a dependencias frágiles. Rechace para configuración inválida o retry sin tope.

## 3. Architecture contract

- **Ownership**: `internal/resilience` gobierna apertura/cierre y reintentos; la dependencia real es del runtime.
- **Invariantes**: (1) abre tras N fallos consecutivos. (2) en open deniega hasta el timeout, luego una única sonda. (3) sonda exitosa cierra, fallida reabre. (4) retry siempre acotado.
- **Data flow**: `Allow` → `RecordSuccess/RecordFailure`; `Retry(Policy, fn)`.
- **Failure modes**: configuración inválida → error; open → deny.
- **Seguridad/privacidad**: sin PII.
- **Performance budget**: O(1) por decisión.

## 4. Exact file manifest

```text
CREATE internal/resilience/resilience.go
CREATE internal/resilience/resilience_test.go
```

## 5. Materialization blocks

### FILE: `internal/resilience/resilience.go`
```yaml
block_id: "GO-RESILIENCE-CORE:internal/resilience/resilience.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "21f34f08a8d04a23ded5544a3b81dbcd4b26ca273021e647db7302a761c4812e"
variables: []
secrets_allowed: false
```
````go
// Package resilience provides a circuit breaker (closed→open→half-open) with
// a bounded retry budget and injectable backoff. Stdlib-only. AUTHORED.
package resilience

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrInvalidConfig = errors.New("resilience: max failures >= 1 and open timeout > 0")
	ErrInvalidRetry  = errors.New("resilience: max attempts must be >= 1")
)

// State is the circuit-breaker state.
type State int

const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

// CircuitBreaker opens after maxFailures consecutive failures and allows a
// single probe after openTimeout.
type CircuitBreaker struct {
	mu                  sync.Mutex
	maxFailures         int
	openTimeout         time.Duration
	now                 func() time.Time
	state               State
	consecutiveFailures int
	openedAt            time.Time
	probeInFlight       bool
}

// NewCircuitBreaker returns a closed breaker.
func NewCircuitBreaker(maxFailures int, openTimeout time.Duration) (*CircuitBreaker, error) {
	if maxFailures < 1 || openTimeout <= 0 {
		return nil, ErrInvalidConfig
	}
	return &CircuitBreaker{maxFailures: maxFailures, openTimeout: openTimeout, now: time.Now, state: StateClosed}, nil
}

func (c *CircuitBreaker) setClock(fn func() time.Time) { c.now = fn }

// Allow reports whether a request may proceed. In open state it denies until
// the timeout elapses, then admits exactly one probe (half-open).
func (c *CircuitBreaker) Allow() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := c.now()
	switch c.state {
	case StateClosed:
		return true
	case StateOpen:
		if now.Sub(c.openedAt) >= c.openTimeout {
			c.state = StateHalfOpen
			c.probeInFlight = true
			return true
		}
		return false
	case StateHalfOpen:
		if c.probeInFlight {
			return false
		}
		c.probeInFlight = true
		return true
	default:
		return false
	}
}

// RecordSuccess closes the breaker on a successful probe and resets failures.
func (c *CircuitBreaker) RecordSuccess() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.state == StateHalfOpen {
		c.state = StateClosed
		c.probeInFlight = false
	}
	c.consecutiveFailures = 0
}

// RecordFailure counts failures; opens the breaker at the threshold.
func (c *CircuitBreaker) RecordFailure() {
	c.mu.Lock()
	defer c.mu.Unlock()
	switch c.state {
	case StateClosed:
		c.consecutiveFailures++
		if c.consecutiveFailures >= c.maxFailures {
			c.state = StateOpen
			c.openedAt = c.now()
		}
	case StateHalfOpen:
		c.state = StateOpen
		c.openedAt = c.now()
		c.probeInFlight = false
	}
}

// State returns the current state.
func (c *CircuitBreaker) State() State {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.state
}

// RetryPolicy bounds retry attempts with an optional backoff.
type RetryPolicy struct {
	MaxAttempts int
	Backoff     func(attempt int) time.Duration
}

// Retry calls fn up to MaxAttempts times, sleeping backoff between attempts.
// It returns nil on first success, or the last error after exhausting attempts.
func Retry(p RetryPolicy, fn func() error) error {
	if p.MaxAttempts < 1 {
		return ErrInvalidRetry
	}
	var err error
	for attempt := 1; attempt <= p.MaxAttempts; attempt++ {
		if err = fn(); err == nil {
			return nil
		}
		if attempt < p.MaxAttempts && p.Backoff != nil {
			if d := p.Backoff(attempt); d > 0 {
				time.Sleep(d)
			}
		}
	}
	return err
}
````

### FILE: `internal/resilience/resilience_test.go`
```yaml
block_id: "GO-RESILIENCE-CORE:internal/resilience/resilience_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "7dc894ebffecceb734967643f99e95bdf93223bc21402153bca190d92266849d"
variables: []
secrets_allowed: false
```
````go
package resilience

import (
	"errors"
	"testing"
	"time"
)

func newTestBreaker(maxFailures int, timeout time.Duration, start time.Time) *CircuitBreaker {
	cb, err := NewCircuitBreaker(maxFailures, timeout)
	if err != nil {
		panic(err)
	}
	cb.setClock(func() time.Time { return start })
	return cb
}

func TestOpensAfterFailures(t *testing.T) {
	start := time.Unix(0, 0)
	cb := newTestBreaker(3, time.Second, start)
	for i := 0; i < 3; i++ {
		cb.RecordFailure()
	}
	if cb.State() != StateOpen {
		t.Fatalf("expected open, got %v", cb.State())
	}
	if cb.Allow() {
		t.Fatal("open breaker should deny before timeout")
	}
}

func TestHalfOpenProbeAndRecovery(t *testing.T) {
	start := time.Unix(0, 0)
	cb := newTestBreaker(1, time.Second, start)
	cb.RecordFailure()
	if cb.State() != StateOpen {
		t.Fatal("should be open")
	}
	// Advance past timeout → half-open probe allowed.
	cb.setClock(func() time.Time { return start.Add(2 * time.Second) })
	if !cb.Allow() {
		t.Fatal("probe should be allowed after timeout")
	}
	if cb.Allow() {
		t.Fatal("only one probe at a time in half-open")
	}
	cb.RecordSuccess()
	if cb.State() != StateClosed {
		t.Fatalf("expected closed after successful probe, got %v", cb.State())
	}
}

func TestHalfOpenFailureReopens(t *testing.T) {
	start := time.Unix(0, 0)
	cb := newTestBreaker(1, time.Second, start)
	cb.RecordFailure()
	cb.setClock(func() time.Time { return start.Add(2 * time.Second) })
	if !cb.Allow() {
		t.Fatal("probe should be allowed")
	}
	cb.RecordFailure()
	if cb.State() != StateOpen {
		t.Fatalf("expected reopen on failed probe, got %v", cb.State())
	}
}

func TestRetrySucceedsAfterTransientFailure(t *testing.T) {
	attempts := 0
	err := Retry(RetryPolicy{MaxAttempts: 3, Backoff: func(int) time.Duration { return 0 }}, func() error {
		attempts++
		if attempts < 3 {
			return errors.New("transient")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if attempts != 3 {
		t.Fatalf("expected 3 attempts, got %d", attempts)
	}
}

func TestRetryExhausts(t *testing.T) {
	attempts := 0
	err := Retry(RetryPolicy{MaxAttempts: 2, Backoff: func(int) time.Duration { return 0 }}, func() error {
		attempts++
		return errors.New("always fails")
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if attempts != 2 {
		t.Fatalf("expected 2 attempts, got %d", attempts)
	}
}

func TestInvalidConfig(t *testing.T) {
	if _, err := NewCircuitBreaker(0, time.Second); !errors.Is(err, ErrInvalidConfig) {
		t.Fatalf("max failures 0 accepted: %v", err)
	}
	if _, err := NewCircuitBreaker(1, 0); !errors.Is(err, ErrInvalidConfig) {
		t.Fatalf("zero timeout accepted: %v", err)
	}
	if err := Retry(RetryPolicy{MaxAttempts: 0}, func() error { return nil }); !errors.Is(err, ErrInvalidRetry) {
		t.Fatalf("zero attempts accepted: %v", err)
	}
}
````


## 6. Configuration surface

Sin variables ni secretos (maxFailures/openTimeout/RetryPolicy se pasan por constructor).

## 7. Dependency bill

| Package/image/tool | Pin exacto | Uso | Licencia | Runtime/build | Fuente oficial |
|---|---|---|---|---|---|
| Go standard library | go 1.26.7 (toolchain) | resiliencia | BSD-3-Clause | runtime | https://go.dev |

## 8. Apply order

1. Componer el backend (mismo módulo).
2. Colocar los dos archivos bajo `internal/resilience/`.
3. Verificar con `go test ./... -count=1` y `go vet ./...`.
4. Rollback: eliminar `internal/resilience/`.

## 9. Verification

- `go test ./internal/resilience/ -count=1`: 6/6 PASS.
- `go test ./... -count=1` (41 paquetes): PASS.
- `go vet ./...`: exit 0.

## 10. Reconstruction evidence

Véase `reconstruction_evidence/GO_RESILIENCE_CORE_2026-09-02_V221.md`.
