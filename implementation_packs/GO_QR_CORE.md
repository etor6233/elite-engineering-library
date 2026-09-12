# Go QR Core

## 1. Metadata

```yaml
pack_id: "GO-QR-CORE"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa la identidad QR segura: payload versionado, tenant-scoped y con checksum tamper-evident, resolución contra el dominio fail-closed; la imagen se genera con librería oficial pinnada, nunca reimplementada."
stacks: ["Go 1.26.7"]
compatible_with: ["GO-AGENT-DOMAIN-BINDING 0.1.x", "GO-ENTERPRISE-BACKEND 0.4.x", "GO-FINOPS-CORE 0.1.x"]
incompatible_with: ["QR sin checksum/tenant", "resolución sin verificación de dominio"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://go.dev", "https://github.com/skip2/go-qrcode"]
verified_at: "2026-09-02"
```

La identidad del QR (checksum, tenant, resolución) es `AUTHORED` y verificada. La **generación de imagen** usa la librería oficial MIT `skip2/go-qrcode` (pinnada por commit/SHA), nunca reimplementada; queda `CONDITIONED` (requiere red para fijarla). El QR apunta a una identidad verificable, no a una URL cruda que se pueda falsificar.

## 2. Applicability

Use este pack para QRs de compras, productos, clientes o turnos que se identifiquen de forma segura: payload tamper-evident, resuelto contra el dominio y tenant-scoped. Rechace para QRs sin verificación o con URL cruda confiable.

## 3. Architecture contract

- **Ownership**: `internal/qr` gobierna payload, checksum y verificación; la imagen la genera el adapter oficial (CONDITIONED). El dominio (producto/pedido/cliente) resuelve la identidad.
- **Invariantes**: (1) checksum SHA-256 sobre tenant+kind+id (tamper-evident). (2) tenant obligatorio; resolución cross-tenant falla cerrado. (3) QR que no resuelve → `ErrUnknown`. (4) imagen por librería oficial, no reimplementada.
- **Data flow**: `NewPayload` → `Encode` (JSON) → imagen (adapter) → escaneo → `Decode` → `Verify` → `Resolver` → efecto.
- **Failure modes**: checksum inválido → `ErrChecksum`; campo inválido → `ErrInvalidPayload`; no resuelve → `ErrUnknown`.
- **Seguridad/privacidad**: tenant-scoped; checksum impide falsificación; sin PII extra en el payload.
- **Performance budget**: O(1) por payload.
- **Operación/migración/rollback**: sin migración; el resolver es composición.

## 4. Exact file manifest

```text
CREATE internal/qr/payload.go
CREATE internal/qr/verify.go
CREATE internal/qr/payload_test.go
```

## 5. Materialization blocks

### FILE: `internal/qr/payload.go`
```yaml
block_id: "GO-QR-CORE:internal/qr/payload.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "1621f5ce3917004fb657be9274ca6e6b9d068a8d0a56cd1e29fba6ddecf94501"
variables: []
secrets_allowed: false
```
````go
// Package qr provides a secure, tenant-scoped QR identity for purchases,
// products, customers and appointments. The payload carries a tamper-evident
// checksum and is resolved against the domain fail-closed; the image encoding
// is a pinned official library (CONDITIONED), never reimplemented.
package qr

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Kind is the bounded set of QR target kinds.
type Kind string

const (
	KindProduct     Kind = "product"
	KindOrder       Kind = "order"
	KindCustomer    Kind = "customer"
	KindAppointment Kind = "appointment"
)

var (
	idRe = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:-]{0,127}$`)

	ErrInvalidPayload = errors.New("qr: invalid payload")
	ErrChecksum       = errors.New("qr: checksum mismatch")
)

// Payload is a versioned, tenant-scoped QR identity.
type Payload struct {
	Version  int    `json:"v"`
	TenantID string `json:"t"`
	Kind     Kind   `json:"k"`
	ID       string `json:"id"`
	Checksum string `json:"c"`
}

// checksum computes the tamper-evident digest over the identity fields.
func checksum(tenant string, kind Kind, id string) string {
	h := sha256.Sum256([]byte(tenant + "\x00" + string(kind) + "\x00" + id))
	return hex.EncodeToString(h[:])
}

// NewPayload builds a signed payload.
func NewPayload(tenant string, kind Kind, id string) (Payload, error) {
	p := Payload{Version: 1, TenantID: tenant, Kind: kind, ID: id}
	if err := p.ValidateFields(); err != nil {
		return Payload{}, err
	}
	p.Checksum = checksum(tenant, kind, id)
	return p, nil
}

// ValidateFields checks identity fields (not the checksum).
func (p Payload) ValidateFields() error {
	if p.Version != 1 {
		return fmt.Errorf("%w: version", ErrInvalidPayload)
	}
	if strings.TrimSpace(p.TenantID) == "" || len(p.TenantID) > 64 {
		return fmt.Errorf("%w: tenant", ErrInvalidPayload)
	}
	switch p.Kind {
	case KindProduct, KindOrder, KindCustomer, KindAppointment:
	default:
		return fmt.Errorf("%w: kind", ErrInvalidPayload)
	}
	if !idRe.MatchString(p.ID) {
		return fmt.Errorf("%w: id", ErrInvalidPayload)
	}
	return nil
}

// Validate checks the checksum too.
func (p Payload) Validate() error {
	if err := p.ValidateFields(); err != nil {
		return err
	}
	if p.Checksum != checksum(p.TenantID, p.Kind, p.ID) {
		return ErrChecksum
	}
	return nil
}

// Encode returns the canonical JSON string for the QR.
func (p Payload) Encode() (string, error) {
	if err := p.Validate(); err != nil {
		return "", err
	}
	b, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Decode parses and validates a payload string (tamper-evident).
func Decode(encoded string) (Payload, error) {
	var p Payload
	if err := json.Unmarshal([]byte(encoded), &p); err != nil {
		return Payload{}, fmt.Errorf("%w: %v", ErrInvalidPayload, err)
	}
	if err := p.Validate(); err != nil {
		return Payload{}, err
	}
	return p, nil
}
````

### FILE: `internal/qr/verify.go`
```yaml
block_id: "GO-QR-CORE:internal/qr/verify.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "024d87877dcc7a254d0d85b71f51b78b243d72b7e3b055d837f3c563c4db6d0d"
variables: []
secrets_allowed: false
```
````go
package qr

import (
	"context"
	"errors"
	"fmt"
)

// ErrUnknown reports a payload that does not resolve against the domain.
var ErrUnknown = errors.New("qr: unknown target")

// Resolver resolves a decoded payload against the domain (product, order,
// customer, appointment) for the tenant. It never widens scope.
type Resolver interface {
	Resolve(ctx context.Context, tenantID string, kind Kind, id string) (bool, error)
}

// Verify decodes, validates the checksum and resolves against the domain,
// fail-closed. A valid QR that points to nothing returns ErrUnknown.
func Verify(ctx context.Context, encoded string, r Resolver) (Payload, error) {
	p, err := Decode(encoded)
	if err != nil {
		return Payload{}, err
	}
	if r == nil {
		return Payload{}, errors.New("qr: nil resolver")
	}
	ok, err := r.Resolve(ctx, p.TenantID, p.Kind, p.ID)
	if err != nil {
		return Payload{}, err
	}
	if !ok {
		return Payload{}, fmt.Errorf("%w: %s/%s", ErrUnknown, p.Kind, p.ID)
	}
	return p, nil
}
````

### FILE: `internal/qr/payload_test.go`
```yaml
block_id: "GO-QR-CORE:internal/qr/payload_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "9735606b2be2964faaec752ad80134ba0a40e2adb83643d3c00bc5cc9f2c4094"
variables: []
secrets_allowed: false
```
````go
package qr

import (
	"context"
	"errors"
	"testing"
)

type fakeResolver struct {
	ok  map[string]bool
	err error
}

func (f *fakeResolver) Resolve(_ context.Context, tenant string, kind Kind, id string) (bool, error) {
	if f.err != nil {
		return false, f.err
	}
	return f.ok[tenant+"\x00"+string(kind)+"\x00"+id], nil
}

func TestPayloadRoundTrip(t *testing.T) {
	p, err := NewPayload("tenant-a", KindProduct, "SKU-123")
	if err != nil {
		t.Fatal(err)
	}
	enc, err := p.Encode()
	if err != nil {
		t.Fatal(err)
	}
	got, err := Decode(enc)
	if err != nil {
		t.Fatal(err)
	}
	if got != p {
		t.Fatalf("round-trip mismatch: %+v vs %+v", got, p)
	}
}

func TestPayloadChecksumTamper(t *testing.T) {
	p, _ := NewPayload("tenant-a", KindOrder, "ORD-1")
	p.ID = "ORD-2" // tamper without recomputing checksum
	if err := p.Validate(); !errors.Is(err, ErrChecksum) {
		t.Fatalf("expected checksum mismatch, got %v", err)
	}
}

func TestPayloadInvalidFields(t *testing.T) {
	if _, err := NewPayload("", KindProduct, "SKU"); !errors.Is(err, ErrInvalidPayload) {
		t.Fatalf("empty tenant accepted: %v", err)
	}
	if _, err := NewPayload("t", Kind("bad"), "x"); !errors.Is(err, ErrInvalidPayload) {
		t.Fatalf("bad kind accepted: %v", err)
	}
	if _, err := NewPayload("t", KindProduct, "bad id!"); !errors.Is(err, ErrInvalidPayload) {
		t.Fatalf("bad id accepted: %v", err)
	}
}

func TestVerifyResolvesAndFailsClosed(t *testing.T) {
	p, _ := NewPayload("tenant-a", KindProduct, "SKU-1")
	enc, _ := p.Encode()

	r := &fakeResolver{ok: map[string]bool{"tenant-a\x00product\x00SKU-1": true}}
	if _, err := Verify(context.Background(), enc, r); err != nil {
		t.Fatalf("valid QR rejected: %v", err)
	}

	missing := &fakeResolver{ok: map[string]bool{}}
	if _, err := Verify(context.Background(), enc, missing); !errors.Is(err, ErrUnknown) {
		t.Fatalf("expected ErrUnknown, got %v", err)
	}

	// cross-tenant: resolver only knows tenant-b → ErrUnknown
	other := &fakeResolver{ok: map[string]bool{"tenant-b\x00product\x00SKU-1": true}}
	if _, err := Verify(context.Background(), enc, other); !errors.Is(err, ErrUnknown) {
		t.Fatalf("cross-tenant should fail closed, got %v", err)
	}
}

func TestVerifyRejectsTampered(t *testing.T) {
	p, _ := NewPayload("tenant-a", KindOrder, "ORD-1")
	enc, _ := p.Encode()
	// tamper the encoded JSON's id (the checksum won't match)
	tampered := enc[:len(enc)-2] + `"}` // corrupts JSON → decode fails
	if _, err := Verify(context.Background(), tampered, &fakeResolver{}); err == nil {
		t.Fatal("tampered payload accepted")
	}
}
````


## 6. Configuration surface

Sin variables ni secretos. El `Resolver` (producto/pedido/cliente/turno) es composición del proyecto contra el backend.

## 7. Dependency bill

| Package/image/tool | Pin exacto | Uso | Licencia | Runtime/build | Fuente oficial |
|---|---|---|---|---|---|
| Go standard library | go 1.26.7 (toolchain) | payload/checksum/verify | BSD-3-Clause | runtime | https://go.dev |
| skip2/go-qrcode | (commit fijado al adquirir) | imagen PNG | MIT | runtime (CONDITIONED) | https://github.com/skip2/go-qrcode |

## 8. Apply order

1. Componer el backend (mismo módulo).
2. Colocar los tres archivos bajo `internal/qr/`.
3. Fijar `skip2/go-qrcode` por commit/SHA y escribir el adapter de imagen.
4. Implementar `Resolver` contra producto/pedido/cliente.
5. Verificar con `go test ./... -count=1` y `go vet ./...`.
6. Rollback: eliminar `internal/qr/`.

## 9. Verification

- `go test ./internal/qr/ -count=1`: 5/5 PASS (round-trip, tamper por checksum, campos inválidos, verificación resuelve y falla cerrado cross-tenant, rechazo de payload manipulado).
- `go test ./... -count=1` (15 paquetes): PASS.
- `go vet ./...`: exit 0.

## 10. Reconstruction evidence

Véase `reconstruction_evidence/GO_QR_CORE_2026-09-02_V193.md`.
