# Go OTP Verification Core

## 1. Metadata

```yaml
pack_id: "GO-OTP-VERIFICATION-CORE"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa la verificación de email/teléfono por código de un solo uso: el código se guarda como hash SHA-256 (nunca en claro), expira tras TTL y admite intentos acotados; fail-closed."
stacks: ["Go 1.26.7"]
compatible_with: ["GO-APPLICANT-ONBOARDING-CORE 0.1.x"]
incompatible_with: ["canal desconocido", "TTL <= 0", "código reutilizado", "intentos agotados"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://go.dev"]
verified_at: "2026-09-02"
```

## 2. Applicability

Use para probar que el solicitante controla su email/teléfono antes de avanzar en el alta. Rechace para canal desconocido, código reutilizado, expirado o intentos agotados.

## 3. Architecture contract

- **Ownership**: `internal/otp` emite/verifica el código; el envío real (email/SMS) es del runtime.
- **Invariantes**: (1) el código sólo existe como hash. (2) un solo uso. (3) expira tras TTL. (4) intentos acotados (lock).
- **Data flow**: `Issue` → `Verify`.
- **Failure modes**: entrada inválida, no encontrado, expirado, lock → error (fail-closed).
- **Seguridad/privacidad**: sin almacenar el código en claro.
- **Performance budget**: O(1) por operación.

## 4. Exact file manifest

```text
CREATE internal/otp/otp.go
CREATE internal/otp/otp_test.go
```

## 5. Materialization blocks

### FILE: `internal/otp/otp.go`
```yaml
block_id: "GO-OTP-VERIFICATION-CORE:internal/otp/otp.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "c72b8920acc63eb109cfcbb6d59b4d81102abd6f5f0f893b75cb18596d1ea0ae"
variables: []
secrets_allowed: false
```
````go
// Package otp provides one-time-code verification for a contact channel
// (email/phone): the code is stored as a SHA-256 hash (never plaintext),
// expires after a TTL and allows a bounded number of attempts. Fail-closed.
// Stdlib-only. AUTHORED.
package otp

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"
)

var (
	ErrInvalidInput = errors.New("otp: invalid input")
	ErrNotFound     = errors.New("otp: not found")
	ErrExpired      = errors.New("otp: expired")
	ErrLocked       = errors.New("otp: too many attempts")
)

const (
	codeDigits  = 6
	maxAttempts = 5
)

type record struct {
	hash      string
	expiresAt time.Time
	attempts  int
}

// Store holds issued codes by tenant+channel+identifier.
type Store struct {
	mu    sync.Mutex
	items map[string]*record
	now   func() time.Time
}

// NewStore returns an empty store with a real clock.
func NewStore() *Store {
	return &Store{items: make(map[string]*record), now: time.Now}
}

func (s *Store) setClock(fn func() time.Time) { s.now = fn }

func key(tenant, channel, identifier string) string {
	return tenant + "\x00" + channel + "\x00" + identifier
}

func generateCode() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

// Issue creates a code for the channel, stores its hash, and returns the code
// (the caller sends it via email/SMS). One code per (tenant, channel, id).
func (s *Store) Issue(tenant, channel, identifier string, ttl time.Duration) (string, error) {
	if strings.TrimSpace(tenant) == "" || (channel != "email" && channel != "phone") ||
		strings.TrimSpace(identifier) == "" || ttl <= 0 {
		return "", ErrInvalidInput
	}
	code, err := generateCode()
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256([]byte(code))
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[key(tenant, channel, identifier)] = &record{
		hash:      hex.EncodeToString(sum[:]),
		expiresAt: s.now().Add(ttl),
		attempts:  0,
	}
	return code, nil
}

// Verify checks a code; on success the code is consumed (one-time). A code
// that expires or exceeds attempts is removed (fail-closed).
func (s *Store) Verify(tenant, channel, identifier, code string) (bool, error) {
	if strings.TrimSpace(tenant) == "" || strings.TrimSpace(identifier) == "" || code == "" {
		return false, ErrInvalidInput
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	k := key(tenant, channel, identifier)
	r, ok := s.items[k]
	if !ok {
		return false, ErrNotFound
	}
	if s.now().After(r.expiresAt) {
		delete(s.items, k)
		return false, ErrExpired
	}
	if r.attempts >= maxAttempts {
		delete(s.items, k)
		return false, ErrLocked
	}
	sum := sha256.Sum256([]byte(code))
	if hex.EncodeToString(sum[:]) == r.hash {
		delete(s.items, k)
		return true, nil
	}
	r.attempts++
	return false, nil
}

// String aids debugging.
func (s *Store) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return fmt.Sprintf("otp(%d)", len(s.items))
}
````

### FILE: `internal/otp/otp_test.go`
```yaml
block_id: "GO-OTP-VERIFICATION-CORE:internal/otp/otp_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "d8a26827750d130a65e8b52bc8790896fd266efec5e32b2f0016c493d518d49c"
variables: []
secrets_allowed: false
```
````go
package otp

import (
	"errors"
	"testing"
	"time"
)

func TestIssueAndVerifyOneTime(t *testing.T) {
	s := NewStore()
	code, err := s.Issue("t", "email", "a@b.c", time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	if len(code) != 6 {
		t.Fatalf("expected 6-digit code, got %q", code)
	}
	ok, err := s.Verify("t", "email", "a@b.c", code)
	if err != nil || !ok {
		t.Fatalf("correct code rejected: %v %v", ok, err)
	}
	// One-time: second verify fails.
	if _, err := s.Verify("t", "email", "a@b.c", code); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected not-found on second use, got %v", err)
	}
}

func TestWrongCodeAndLock(t *testing.T) {
	s := NewStore()
	code, _ := s.Issue("t", "email", "a@b.c", time.Minute)
	for i := 0; i < maxAttempts; i++ {
		if _, err := s.Verify("t", "email", "a@b.c", "000000"); err != nil {
			t.Fatalf("attempt %d unexpected err: %v", i, err)
		}
	}
	// Now even the correct code is locked.
	if _, err := s.Verify("t", "email", "a@b.c", code); !errors.Is(err, ErrLocked) {
		t.Fatalf("expected locked, got %v", err)
	}
}

func TestExpired(t *testing.T) {
	s := NewStore()
	start := time.Unix(0, 0)
	s.setClock(func() time.Time { return start })
	code, _ := s.Issue("t", "email", "a@b.c", time.Minute)
	s.setClock(func() time.Time { return start.Add(2 * time.Minute) })
	if _, err := s.Verify("t", "email", "a@b.c", code); !errors.Is(err, ErrExpired) {
		t.Fatalf("expected expired, got %v", err)
	}
}

func TestInvalidInput(t *testing.T) {
	s := NewStore()
	if _, err := s.Issue("", "email", "a@b.c", time.Minute); !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("empty tenant accepted: %v", err)
	}
	if _, err := s.Issue("t", "sms", "a@b.c", time.Minute); !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("bad channel accepted: %v", err)
	}
	if _, err := s.Issue("t", "email", "a@b.c", 0); !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("zero ttl accepted: %v", err)
	}
}

func TestTenantIsolation(t *testing.T) {
	s := NewStore()
	code, _ := s.Issue("t1", "email", "a@b.c", time.Minute)
	if _, err := s.Verify("t2", "email", "a@b.c", code); !errors.Is(err, ErrNotFound) {
		t.Fatalf("code leaked across tenants: %v", err)
	}
}
````


## 6. Configuration surface

Sin variables ni secretos (TTL se pasa por `Issue`).

## 7. Dependency bill

| Package/image/tool | Pin exacto | Uso | Licencia | Runtime/build | Fuente oficial |
|---|---|---|---|---|---|
| Go standard library | go 1.26.7 (toolchain) | verificación OTP | BSD-3-Clause | runtime | https://go.dev |

## 8. Apply order

1. Componer el backend (mismo módulo).
2. Colocar los dos archivos bajo `internal/otp/`.
3. Verificar con `go test ./... -count=1` y `go vet ./...`.
4. Rollback: eliminar `internal/otp/`.

## 9. Verification

- `go test ./internal/otp/ -count=1`: 5/5 PASS.
- `go test ./... -count=1` (43 paquetes): PASS.
- `go vet ./...`: exit 0.

## 10. Reconstruction evidence

Véase `reconstruction_evidence/GO_OTP_VERIFICATION_CORE_2026-09-02_V222.md`.
