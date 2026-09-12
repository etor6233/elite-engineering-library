# Go Cache Core

## 1. Metadata

```yaml
pack_id: "GO-CACHE-CORE"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa un borde de cache gobernado: keys tenant-scoped, store en memoria con TTL y evicción LRU, y carga cache-aside single-flight que evita stampede y no envenena la cache con errores."
stacks: ["Go 1.26.7"]
compatible_with: ["GO-ENTERPRISE-BACKEND 0.4.x", "GO-ML-AI-FOUNDATION 0.1.x", "GO-SEARCH-CORE 0.1.x"]
incompatible_with: ["cache tratada como source of truth", "TTL no positivo", "keys sin tenant obligatorio"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://go.dev", "https://redis.io/docs/latest/"]
verified_at: "2026-09-02"
```

Este pack no afirma durabilidad ni alta disponibilidad. El store en memoria es la implementación de referencia determinista para un proceso; el adapter Redis (protocolo RESP, `redis.io`) es la autoridad para un cache compartido/HA y permanece `CONDITIONED`: se materializa con un cliente fijado y un target real sólo cuando el workload lo exige. Ningún valor cacheado es source of truth.

## 2. Applicability

Use este pack para sesiones, hot data, resultados intermedios y protección contra stampede en lecturas costosas, siempre como capa de aceleración sobre un owner de dominio. Es consumido por el agente conversacional (sesión/contexto) y por consultas read-mostly del portal.

Rechace este pack para: datos que deban ser fuente de verdad transaccional (eso es `TX-DATABASE`); claves sin tenant; o workloads que requieran consistencia fuerte entre réplicas sin un cache compartido admitido.

## 3. Architecture contract

- **Ownership**: el paquete `internal/cache` gobierna el borde `Store`, las keys tenant-scoped, el TTL, la evicción LRU y la carga single-flight. El dominio propietario repuebla la cache; nunca se lee la cache como verdad autoritativa.
- **Invariantes**: (1) `tenant_id` es obligatorio en toda key (`Scope.Key` lo niega vacío). (2) TTL no positivo se rechaza. (3) Un valor mayor al presupuesto se rechaza (`ErrTooLarge`), no se descarta en silencio. (4) Un error de carga no se cachea (sin envenenamiento). (5) Lecturas concurrentes de la misma key en miss ejecutan exactamente una carga.
- **Data flow**: `Loader.Get` → `Store.Get` (hit) → si miss, single-flight → `load(ctx)` → `Store.Set` sólo en éxito → respuesta. `MemoryStore` usa LRU y expiración perezosa.
- **Failure modes**: key inválida → `ErrInvalidKey`; TTL inválido → `ErrInvalidTTL`; valor sobredimensionado → `ErrTooLarge`; carga fallida → error al llamador sin residuo cacheado.
- **Seguridad/privacidad**: keys con namespace por tenant evitan colisiones cross-tenant en un cache compartido; los valores se copian en lectura para impedir mutación del cache.
- **Performance budget**: O(1) amortizado por operación; presupuesto de bytes acotado con evicción LRU.
- **Operación/migración/rollback**: sin migración ni estado durable; reemplazar el store (memoria↔Redis) es una decisión de composición, no de datos.

## 4. Exact file manifest

```text
CREATE internal/cache/cache.go
CREATE internal/cache/memory.go
CREATE internal/cache/load.go
CREATE internal/cache/cache_test.go
CREATE internal/cache/memory_test.go
CREATE internal/cache/load_test.go
```

## 5. Materialization blocks

### FILE: `internal/cache/cache.go`
```yaml
block_id: "GO-CACHE-CORE:internal/cache/cache.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "14fffefe553e98acc6488694109106bf5180e9803cb4f09d720d3712d3b6ecca"
variables: []
secrets_allowed: false
```
````go
// Package cache provides a governed cache boundary: tenant-scoped keys, a
// TTL-aware LRU memory store, and single-flight cache-aside loading. A cache
// never holds source of truth; misses are repopulated from the owning store
// and expiry/invalidation are expected.
package cache

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// Store is the cache boundary. Implementations must be safe for concurrent use
// and must never serve a miss as an authoritative value.
type Store interface {
	Get(ctx context.Context, key string) (value []byte, ok bool, err error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}

var (
	ErrInvalidKey = errors.New("cache: invalid key")
	ErrInvalidTTL = errors.New("cache: invalid ttl")
)

var keyRe = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:/-]{0,255}$`)

// ValidateKey rejects empty, unsafe or overlong keys to keep a shared cache
// namespace collision-free and free of injection surprises.
func ValidateKey(key string) error {
	if !keyRe.MatchString(key) {
		return fmt.Errorf("%w: %q", ErrInvalidKey, key)
	}
	return nil
}

// Scope namespaces a key by tenant (mandatory) and optional kind/id.
type Scope struct {
	TenantID string
	Kind     string
	ID       string
}

// Key builds a tenant-scoped, collision-free key. Empty tenant is denied.
func (s Scope) Key() (string, error) {
	if strings.TrimSpace(s.TenantID) == "" {
		return "", fmt.Errorf("%w: empty tenant", ErrInvalidKey)
	}
	parts := []string{"t", s.TenantID}
	if s.Kind != "" {
		parts = append(parts, s.Kind)
	}
	if s.ID != "" {
		parts = append(parts, s.ID)
	}
	k := strings.Join(parts, ":")
	if err := ValidateKey(k); err != nil {
		return "", err
	}
	return k, nil
}
````

### FILE: `internal/cache/memory.go`
```yaml
block_id: "GO-CACHE-CORE:internal/cache/memory.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "fca7d8abeb45a3ba2729241fdd26abc6f4327c6d36956df10dcf0055cfc1302d"
variables: []
secrets_allowed: false
```
````go
package cache

import (
	"container/list"
	"context"
	"errors"
	"sync"
	"time"
)

// ErrTooLarge reports a value exceeding the store's byte budget.
var ErrTooLarge = errors.New("cache: value exceeds budget")

type entry struct {
	key      string
	value    []byte
	expires  time.Time
	listElem *list.Element
}

// MemoryStore is a bounded, TTL-aware, LRU-evicting cache for a single
// process. It is the deterministic reference implementation; a Redis adapter
// provides the shared/HA equivalent and remains CONDITIONED.
type MemoryStore struct {
	mu        sync.Mutex
	items     map[string]*entry
	lru       *list.List // front = most recently used
	maxBytes  int64
	currBytes int64
	clock     func() time.Time
}

// NewMemoryStore returns a store with the given byte budget (default 1 MiB).
func NewMemoryStore(maxBytes int64) *MemoryStore {
	if maxBytes <= 0 {
		maxBytes = 1 << 20
	}
	return &MemoryStore{
		items:    make(map[string]*entry),
		lru:      list.New(),
		maxBytes: maxBytes,
		clock:    time.Now,
	}
}

func (m *MemoryStore) now() time.Time {
	if m.clock != nil {
		return m.clock()
	}
	return time.Now()
}

// Get returns a copy of the cached value, or ok=false on miss/expiry.
func (m *MemoryStore) Get(_ context.Context, key string) ([]byte, bool, error) {
	if err := ValidateKey(key); err != nil {
		return nil, false, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	e, ok := m.items[key]
	if !ok {
		return nil, false, nil
	}
	if m.now().After(e.expires) {
		m.remove(key)
		return nil, false, nil
	}
	m.lru.MoveToFront(e.listElem)
	out := make([]byte, len(e.value))
	copy(out, e.value)
	return out, true, nil
}

// Set stores a copy of value with a positive TTL and evicts LRU to fit budget.
func (m *MemoryStore) Set(_ context.Context, key string, value []byte, ttl time.Duration) error {
	if err := ValidateKey(key); err != nil {
		return err
	}
	if ttl <= 0 {
		return ErrInvalidTTL
	}
	if int64(len(value)) > m.maxBytes {
		return ErrTooLarge
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.set(key, value, m.now().Add(ttl))
	return nil
}

func (m *MemoryStore) set(key string, value []byte, expires time.Time) {
	if e, ok := m.items[key]; ok {
		m.currBytes -= int64(len(e.value))
		e.value = make([]byte, len(value))
		copy(e.value, value)
		e.expires = expires
		m.currBytes += int64(len(e.value))
		m.lru.MoveToFront(e.listElem)
		return
	}
	e := &entry{key: key, value: append([]byte(nil), value...), expires: expires}
	e.listElem = m.lru.PushFront(e)
	m.items[key] = e
	m.currBytes += int64(len(e.value))

	for m.currBytes > m.maxBytes {
		back := m.lru.Back()
		if back == nil {
			break
		}
		m.remove(back.Value.(*entry).key)
	}
}

// Delete removes a key; it is idempotent.
func (m *MemoryStore) Delete(_ context.Context, key string) error {
	if err := ValidateKey(key); err != nil {
		return err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.items[key]; !ok {
		return nil
	}
	m.remove(key)
	return nil
}

func (m *MemoryStore) remove(key string) {
	if e, ok := m.items[key]; ok {
		m.lru.Remove(e.listElem)
		m.currBytes -= int64(len(e.value))
		delete(m.items, key)
	}
}

// Len reports the number of cached entries.
func (m *MemoryStore) Len() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.items)
}
````

### FILE: `internal/cache/load.go`
```yaml
block_id: "GO-CACHE-CORE:internal/cache/load.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "c50ba4fdbfc7174328b66f16357241c530d7480752722002f0ddb19926347eb8"
variables: []
secrets_allowed: false
```
````go
package cache

import (
	"context"
	"sync"
	"time"
)

// Loader wraps a Store with cache-aside single-flight loading: concurrent
// misses for the same key trigger exactly one load and share the result,
// preventing a thundering herd against the source of truth.
type Loader struct {
	Store  Store
	mu     sync.Mutex
	flight map[string]*call
}

type call struct {
	done chan struct{}
	val  []byte
	err  error
}

// NewLoader returns a Loader over the given Store.
func NewLoader(s Store) *Loader {
	return &Loader{Store: s, flight: make(map[string]*call)}
}

// Get returns the cached value or loads and populates it once per key. A load
// error is returned without poisoning the cache (nothing is stored).
func (l *Loader) Get(ctx context.Context, key string, ttl time.Duration, load func(context.Context) ([]byte, error)) ([]byte, error) {
	if v, ok, err := l.Store.Get(ctx, key); err == nil && ok {
		return v, nil
	}

	l.mu.Lock()
	if c, ok := l.flight[key]; ok {
		l.mu.Unlock()
		<-c.done
		return c.val, c.err
	}
	c := &call{done: make(chan struct{})}
	l.flight[key] = c
	l.mu.Unlock()

	c.val, c.err = load(ctx)
	if c.err == nil {
		_ = l.Store.Set(ctx, key, c.val, ttl)
	}

	close(c.done)
	l.mu.Lock()
	delete(l.flight, key)
	l.mu.Unlock()
	return c.val, c.err
}
````

### FILE: `internal/cache/cache_test.go`
```yaml
block_id: "GO-CACHE-CORE:internal/cache/cache_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "181529d42b172df778f0ed112b98f09631f23a205d06ef9b79ff1a85e88e8952"
variables: []
secrets_allowed: false
```
````go
package cache

import (
	"errors"
	"testing"
)

func TestScopeKey(t *testing.T) {
	k, err := (Scope{TenantID: "tenant-a", Kind: "session", ID: "conv-1"}).Key()
	if err != nil {
		t.Fatal(err)
	}
	if k != "t:tenant-a:session:conv-1" {
		t.Fatalf("unexpected key: %q", k)
	}
	if _, err := (Scope{TenantID: "", Kind: "x"}).Key(); !errors.Is(err, ErrInvalidKey) {
		t.Fatalf("empty tenant accepted: %v", err)
	}
}

func TestValidateKey(t *testing.T) {
	for _, ok := range []string{"a", "t:tenant-a:session:conv-1", "ABC_123.-/"} {
		if err := ValidateKey(ok); err != nil {
			t.Fatalf("valid key %q rejected: %v", ok, err)
		}
	}
	for _, bad := range []string{"", " bad", "bad key!", "k:with:space "} {
		if err := ValidateKey(bad); !errors.Is(err, ErrInvalidKey) {
			t.Fatalf("invalid key %q accepted: %v", bad, err)
		}
	}
}
````

### FILE: `internal/cache/memory_test.go`
```yaml
block_id: "GO-CACHE-CORE:internal/cache/memory_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "4b8fdb2be7de30193ee657324233bfcf8bbb06b2aecddaac0f7236ed9ae41927"
variables: []
secrets_allowed: false
```
````go
package cache

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestMemorySetGet(t *testing.T) {
	s := NewMemoryStore(1 << 20)
	ctx := context.Background()
	if err := s.Set(ctx, "k1", []byte("v1"), time.Minute); err != nil {
		t.Fatal(err)
	}
	v, ok, err := s.Get(ctx, "k1")
	if err != nil || !ok || string(v) != "v1" {
		t.Fatalf("get: %q ok=%v err=%v", v, ok, err)
	}
	if _, ok, _ := s.Get(ctx, "missing"); ok {
		t.Fatal("missing key reported hit")
	}
}

func TestMemoryGetReturnsCopy(t *testing.T) {
	s := NewMemoryStore(1 << 20)
	ctx := context.Background()
	_ = s.Set(ctx, "k", []byte("abc"), time.Minute)
	v, _, _ := s.Get(ctx, "k")
	v[0] = 'z'
	v2, _, _ := s.Get(ctx, "k")
	if string(v2) != "abc" {
		t.Fatalf("cache value mutated: %q", v2)
	}
}

func TestMemoryTTLExpiry(t *testing.T) {
	s := NewMemoryStore(1 << 20)
	var now time.Time
	s.clock = func() time.Time { return now }
	ctx := context.Background()
	_ = s.Set(ctx, "k", []byte("v"), time.Minute)
	now = now.Add(2 * time.Minute)
	if _, ok, _ := s.Get(ctx, "k"); ok {
		t.Fatal("expired key reported hit")
	}
	if s.Len() != 0 {
		t.Fatalf("expired key not evicted, len=%d", s.Len())
	}
}

func TestMemoryLRUEviction(t *testing.T) {
	s := NewMemoryStore(30) // budget: three 10-byte values
	ctx := context.Background()
	_ = s.Set(ctx, "k1", []byte("aaaaaaaaaa"), time.Minute)
	_ = s.Set(ctx, "k2", []byte("bbbbbbbbbb"), time.Minute)
	_ = s.Set(ctx, "k3", []byte("cccccccccc"), time.Minute)
	// touch k1 so it is most recently used
	_, _, _ = s.Get(ctx, "k1")
	_ = s.Set(ctx, "k4", []byte("dddddddddd"), time.Minute) // exceeds budget → evict LRU (k2)
	if _, ok, _ := s.Get(ctx, "k2"); ok {
		t.Fatal("LRU entry k2 should have been evicted")
	}
	for _, k := range []string{"k1", "k3", "k4"} {
		if _, ok, _ := s.Get(ctx, k); !ok {
			t.Fatalf("expected %s to remain", k)
		}
	}
}

func TestMemoryDeleteIdempotent(t *testing.T) {
	s := NewMemoryStore(1 << 20)
	ctx := context.Background()
	_ = s.Set(ctx, "k", []byte("v"), time.Minute)
	if err := s.Delete(ctx, "k"); err != nil {
		t.Fatal(err)
	}
	if err := s.Delete(ctx, "k"); err != nil {
		t.Fatalf("delete not idempotent: %v", err)
	}
}

func TestMemoryValidation(t *testing.T) {
	s := NewMemoryStore(1 << 20)
	ctx := context.Background()
	if err := s.Set(ctx, "bad key!", []byte("v"), time.Minute); !errors.Is(err, ErrInvalidKey) {
		t.Fatalf("invalid key accepted: %v", err)
	}
	if err := s.Set(ctx, "k", []byte("v"), 0); !errors.Is(err, ErrInvalidTTL) {
		t.Fatalf("zero ttl accepted: %v", err)
	}
	if err := s.Set(ctx, "k", make([]byte, 2<<20), time.Minute); !errors.Is(err, ErrTooLarge) {
		t.Fatalf("oversized value accepted: %v", err)
	}
}
````

### FILE: `internal/cache/load_test.go`
```yaml
block_id: "GO-CACHE-CORE:internal/cache/load_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "8391defccf6d63e9c5d1b41479bcb74ed398f65ad83c349b407bbbb6741a01aa"
variables: []
secrets_allowed: false
```
````go
package cache

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestLoaderStampsOnce(t *testing.T) {
	s := NewMemoryStore(1 << 20)
	l := NewLoader(s)
	var calls int32
	load := func(context.Context) ([]byte, error) {
		atomic.AddInt32(&calls, 1)
		time.Sleep(20 * time.Millisecond)
		return []byte("value"), nil
	}
	const n = 20
	var wg sync.WaitGroup
	values := make([]string, n)
	errs := make([]error, n)
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			v, err := l.Get(context.Background(), "k1", time.Minute, load)
			values[i] = string(v)
			errs[i] = err
		}(i)
	}
	wg.Wait()
	if got := atomic.LoadInt32(&calls); got != 1 {
		t.Fatalf("expected exactly 1 load, got %d", got)
	}
	for i := 0; i < n; i++ {
		if errs[i] != nil || values[i] != "value" {
			t.Fatalf("caller %d got %q/%v", i, values[i], errs[i])
		}
	}
}

func TestLoaderPopulatesCache(t *testing.T) {
	s := NewMemoryStore(1 << 20)
	l := NewLoader(s)
	var calls int32
	load := func(context.Context) ([]byte, error) {
		atomic.AddInt32(&calls, 1)
		return []byte("value"), nil
	}
	if _, err := l.Get(context.Background(), "k", time.Minute, load); err != nil {
		t.Fatal(err)
	}
	if _, err := l.Get(context.Background(), "k", time.Minute, load); err != nil {
		t.Fatal(err)
	}
	if got := atomic.LoadInt32(&calls); got != 1 {
		t.Fatalf("expected cache hit on second Get, loads=%d", got)
	}
}

func TestLoaderErrorNotCached(t *testing.T) {
	s := NewMemoryStore(1 << 20)
	l := NewLoader(s)
	calls := 0
	boom := errors.New("boom")
	load := func(context.Context) ([]byte, error) {
		calls++
		return nil, boom
	}
	if _, err := l.Get(context.Background(), "k", time.Minute, load); !errors.Is(err, boom) {
		t.Fatalf("expected boom, got %v", err)
	}
	if _, err := l.Get(context.Background(), "k", time.Minute, load); !errors.Is(err, boom) {
		t.Fatalf("expected boom on retry, got %v", err)
	}
	if calls != 2 {
		t.Fatalf("expected reload after error, calls=%d", calls)
	}
	if s.Len() != 0 {
		t.Fatalf("error poisoned cache, len=%d", s.Len())
	}
}
````


## 6. Configuration surface

Sin variables ni secretos. El presupuesto de bytes se fija en construcción (`NewMemoryStore(maxBytes)`); el TTL es un parámetro de cada `Set`/`Get` y se valida (positivo).

## 7. Dependency bill

| Package/image/tool | Pin exacto | Uso | Licencia | Runtime/build | Fuente oficial |
|---|---|---|---|---|---|
| Go standard library | go 1.26.7 (toolchain) | toda la implementación | BSD-3-Clause | runtime | https://go.dev |

Redis (RESP) permanece como autoridad del adapter compartido/HA no materializado en este pack; se adquiere por expediente separado con cliente fijado si el workload lo exige.

## 8. Apply order

1. Componer el backend (`GO-ENTERPRISE-BACKEND 0.4.x`) o disponer del módulo `elite.local/enterprise` con `go 1.26`.
2. Colocar los seis archivos bajo `internal/cache/`. No hay migración ni conflicto de ruta.
3. Verificar con `go test ./... -count=1` y `go vet ./...`.
4. Rollback: eliminar `internal/cache/`; no deja estado.

## 9. Verification

- `go test ./internal/cache/ -count=1`: 11/11 PASS (keys tenant-scoped con negativos, validación de key, set/get con copia, expiración TTL determinista, evicción LRU con toque MRU, delete idempotente, TTL/oversize inválidos, stampede single-flight 20→1 llamada, populate en hit, error no cacheado).
- `go test ./... -count=1` (aifoundation + search + cache): PASS.
- `go vet ./...`: exit 0.

## 10. Reconstruction evidence

Véase `reconstruction_evidence/GO_CACHE_CORE_2026-09-02_V177.md`.
