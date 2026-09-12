# Go Channels Core

## 1. Metadata

```yaml
pack_id: "GO-CHANNELS-CORE"
pack_version: "0.4.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa envelopes y dispatch tenant-scoped que conservan identidad ingress y generan una DeliveryKey outbound estable; el registry demuestra al composition root que existe al menos un adapter."
stacks: ["Go 1.26.7"]
compatible_with: ["GO-CONNECTED-CONVERSATION-RUNTIME 0.1.x", "PYTHON_META_WHATSAPP_CLOUD_ADAPTER"]
incompatible_with: ["respuesta por canal distinto", "tenant inferido del contenido", "inbound sin identidad provider", "outbound sin receipt/deduplicación"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://go.dev", "https://sre.google/sre-book/reliable-product-launches/"]
verified_at: "2026-09-04"
```

Todos los bloques son `AUTHORED`. El core está probado; cada adapter live conserva su propio gate de autenticidad, webhook, rate limits, receipt y reconciliación.

## 2. Applicability

Use para abstraer WhatsApp, email u otra mensajería sin perder la identidad que gobierna idempotencia. El tenant y el contacto deben estar pre-resueltos fuera del texto.

## 3. Architecture contract

- Inbound exige `ProviderMessageID` y `OccurredAt`; outbound exige `DeliveryKey`.
- Dispatcher entrega el envelope completo al responder y devuelve por el mismo tenant/canal/contacto/thread.
- La key SHA-256 es determinista por mensaje/purpose; el adapter debe usarla o persistir un receipt propio.
- Registry rechaza duplicados y expone `Count` para que el composition root falle si no hay adapters.

## 4. Exact file manifest

```text
CREATE internal/channels/channel.go
CREATE internal/channels/registry.go
CREATE internal/channels/dispatcher.go
CREATE internal/channels/channels_test.go
```

## 5. Materialization blocks

### FILE: `internal/channels/channel.go`
```yaml
block_id: "GO-CHANNELS-CORE:internal/channels/channel.go:v4"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "2bd88fc11511fac919d8313b60784ee17b9d86c9fbe81877b3d49d0661624a0a"
variables: []
secrets_allowed: false
```
````go
// Package channels provides the channel abstraction for the conversational
// agent: a Channel contract (receive/send), a registry keyed by channel code,
// and a dispatcher that routes inbound messages to a responder and sends the
// reply back. The real adapters (Meta Graph, Twilio, SMTP) are CONDITIONED.
package channels

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// Direction is inbound or outbound.
type Direction string

const (
	DirectionIn  Direction = "in"
	DirectionOut Direction = "out"
)

var (
	codeRe = regexp.MustCompile(`^[a-z][a-z0-9_]{0,31}$`)

	ErrInvalidMessage = errors.New("channels: invalid message")
	ErrUnknownChannel = errors.New("channels: unknown channel")
	ErrDuplicateCode  = errors.New("channels: duplicate channel code")
)

// Message is a channel-agnostic inbound/outbound message.
type Message struct {
	ChannelCode       string
	TenantID          string
	ExternalID        string // contact/recipient id in the channel
	ThreadID          string // conversation thread (optional)
	ProviderMessageID string // provider-owned inbound identity; never LLM-authored
	DeliveryKey       string // caller-owned outbound deduplication identity
	OccurredAt        time.Time
	Direction         Direction
	Text              string
}

// Validate enforces the message contract.
func (m Message) Validate() error {
	if !codeRe.MatchString(m.ChannelCode) {
		return fmt.Errorf("%w: channel code", ErrInvalidMessage)
	}
	if strings.TrimSpace(m.TenantID) == "" || len(m.TenantID) > 64 {
		return fmt.Errorf("%w: tenant", ErrInvalidMessage)
	}
	if strings.TrimSpace(m.ExternalID) == "" || len(m.ExternalID) > 256 {
		return fmt.Errorf("%w: external id", ErrInvalidMessage)
	}
	switch m.Direction {
	case DirectionIn:
		if strings.TrimSpace(m.ProviderMessageID) == "" || len(m.ProviderMessageID) > 256 || m.OccurredAt.IsZero() {
			return fmt.Errorf("%w: inbound identity", ErrInvalidMessage)
		}
	case DirectionOut:
		if len(strings.TrimSpace(m.DeliveryKey)) < 16 || len(m.DeliveryKey) > 128 {
			return fmt.Errorf("%w: outbound delivery key", ErrInvalidMessage)
		}
	default:
		return fmt.Errorf("%w: direction", ErrInvalidMessage)
	}
	if strings.TrimSpace(m.Text) == "" {
		return fmt.Errorf("%w: empty text", ErrInvalidMessage)
	}
	return nil
}

// Channel is the receive/send boundary for one messaging surface.
type Channel interface {
	Code() string
	Receive(ctx context.Context) ([]Message, error)
	Send(ctx context.Context, m Message) error
}
````

### FILE: `internal/channels/registry.go`
```yaml
block_id: "GO-CHANNELS-CORE:internal/channels/registry.go:v4"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "f937df0b061875a80cbf6206c312b576b2ef76397c03dfa59fe88e85b587e4d9"
variables: []
secrets_allowed: false
```
````go
package channels

import (
	"fmt"
	"sync"
)

// Registry authorizes at most one channel per code.
type Registry struct {
	mu       sync.RWMutex
	channels map[string]Channel
}

// NewRegistry returns an empty registry.
func NewRegistry() *Registry {
	return &Registry{channels: make(map[string]Channel)}
}

// Register binds a channel to its code, rejecting duplicates and invalid codes.
func (r *Registry) Register(c Channel) error {
	if c == nil {
		return fmt.Errorf("channels: nil channel")
	}
	if !codeRe.MatchString(c.Code()) {
		return fmt.Errorf("%w: invalid code %q", ErrInvalidMessage, c.Code())
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.channels[c.Code()]; ok {
		return fmt.Errorf("%w: %s", ErrDuplicateCode, c.Code())
	}
	r.channels[c.Code()] = c
	return nil
}

// For returns the channel bound to a code, if any.
func (r *Registry) For(code string) (Channel, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.channels[code]
	return c, ok
}

// Count returns the number of explicitly registered channel adapters.
func (r *Registry) Count() int {
	if r == nil {
		return 0
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.channels)
}
````

### FILE: `internal/channels/dispatcher.go`
```yaml
block_id: "GO-CHANNELS-CORE:internal/channels/dispatcher.go:v4"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "eebb41f403b53f3dfba4df644e4d391e9856d23b5e731a79320717b86571e857"
variables: []
secrets_allowed: false
```
````go
package channels

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
)

// Responder handles the complete inbound envelope so command identity, scope
// and correlation are never lost before a side effect.
type Responder func(ctx context.Context, message Message) (string, error)

// Dispatcher routes inbound messages to a Responder and sends the reply back
// through the same channel. It never widens the tenant or channel scope.
type Dispatcher struct {
	Registry *Registry
	Respond  Responder
}

// Handle processes one inbound message and sends the reply via its channel.
func (d *Dispatcher) Handle(ctx context.Context, m Message) error {
	if err := m.Validate(); err != nil {
		return err
	}
	if m.Direction != DirectionIn {
		return ErrInvalidMessage
	}
	if d.Registry == nil {
		return ErrUnknownChannel
	}
	ch, ok := d.Registry.For(m.ChannelCode)
	if !ok {
		return ErrUnknownChannel
	}
	if d.Respond == nil {
		return ErrUnknownChannel
	}
	reply, err := d.Respond(ctx, m)
	if err != nil {
		return err
	}
	return ch.Send(ctx, Message{
		ChannelCode: m.ChannelCode,
		TenantID:    m.TenantID,
		ExternalID:  m.ExternalID,
		ThreadID:    m.ThreadID,
		DeliveryKey: deliveryKey(m),
		Direction:   DirectionOut,
		Text:        reply,
	})
}

func deliveryKey(m Message) string {
	sum := sha256.Sum256([]byte(m.TenantID + "\x00" + m.ChannelCode + "\x00" + m.ProviderMessageID + "\x00reply"))
	return hex.EncodeToString(sum[:])
}
````

### FILE: `internal/channels/channels_test.go`
```yaml
block_id: "GO-CHANNELS-CORE:internal/channels/channels_test.go:v4"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "5d8f5ff50abc0dca1ac89b66d2d76142d9268e28e1c157d489ab3bbbeab8d85c"
variables: []
secrets_allowed: false
```
````go
package channels

import (
	"context"
	"errors"
	"testing"
	"time"
)

type fakeChannel struct {
	code string
	in   []Message
	out  []Message
}

func (f *fakeChannel) Code() string { return f.code }
func (f *fakeChannel) Receive(context.Context) ([]Message, error) {
	m := f.in
	f.in = nil
	return m, nil
}
func (f *fakeChannel) Send(_ context.Context, m Message) error {
	f.out = append(f.out, m)
	return nil
}

func TestMessageValidate(t *testing.T) {
	ok := Message{ChannelCode: "whatsapp", TenantID: "t", ExternalID: "e", ProviderMessageID: "wamid.1", OccurredAt: time.Now(), Direction: DirectionIn, Text: "hola"}
	if err := ok.Validate(); err != nil {
		t.Fatalf("valid message rejected: %v", err)
	}
	bad := ok
	bad.Direction = "sideways"
	if err := bad.Validate(); !errors.Is(err, ErrInvalidMessage) {
		t.Fatalf("bad direction accepted: %v", err)
	}
	empty := ok
	empty.Text = "  "
	if err := empty.Validate(); !errors.Is(err, ErrInvalidMessage) {
		t.Fatalf("empty text accepted: %v", err)
	}
	missingIdentity := ok
	missingIdentity.ProviderMessageID = ""
	if err := missingIdentity.Validate(); !errors.Is(err, ErrInvalidMessage) {
		t.Fatalf("missing provider identity accepted: %v", err)
	}
}

func TestRegistryOnePerCode(t *testing.T) {
	r := NewRegistry()
	if err := r.Register(&fakeChannel{code: "whatsapp"}); err != nil {
		t.Fatal(err)
	}
	if err := r.Register(&fakeChannel{code: "whatsapp"}); !errors.Is(err, ErrDuplicateCode) {
		t.Fatalf("duplicate code accepted: %v", err)
	}
	if err := r.Register(&fakeChannel{code: "BAD!"}); err == nil {
		t.Fatal("invalid code accepted")
	}
	if _, ok := r.For("whatsapp"); !ok {
		t.Fatal("expected whatsapp registered")
	}
	if _, ok := r.For("telegram"); ok {
		t.Fatal("unexpected telegram")
	}
}

func TestDispatcherRoutesAndReplies(t *testing.T) {
	wa := &fakeChannel{code: "whatsapp"}
	r := NewRegistry()
	_ = r.Register(wa)

	d := &Dispatcher{
		Registry: r,
		Respond: func(_ context.Context, message Message) (string, error) {
			if message.ProviderMessageID == "" || message.ThreadID != "thread-1" {
				t.Fatalf("responder lost envelope: %+v", message)
			}
			return "respuesta a: " + message.Text, nil
		},
	}
	in := Message{ChannelCode: "whatsapp", TenantID: "t", ExternalID: "e1", ThreadID: "thread-1", ProviderMessageID: "wamid.1", OccurredAt: time.Now(), Direction: DirectionIn, Text: "hola"}
	if err := d.Handle(context.Background(), in); err != nil {
		t.Fatal(err)
	}
	if len(wa.out) != 1 {
		t.Fatalf("expected 1 outbound, got %d", len(wa.out))
	}
	out := wa.out[0]
	if out.Text != "respuesta a: hola" || out.ExternalID != "e1" || out.Direction != DirectionOut || len(out.DeliveryKey) != 64 {
		t.Fatalf("wrong outbound: %+v", out)
	}
	firstKey := out.DeliveryKey
	if err := d.Handle(context.Background(), in); err != nil || len(wa.out) != 2 || wa.out[1].DeliveryKey != firstKey {
		t.Fatalf("replay changed delivery key: err=%v out=%+v", err, wa.out)
	}
	in.ProviderMessageID = "wamid.2"
	if err := d.Handle(context.Background(), in); err != nil || wa.out[2].DeliveryKey == firstKey {
		t.Fatalf("distinct inbound reused delivery key: err=%v out=%+v", err, wa.out)
	}
}

func TestDispatcherUnknownChannel(t *testing.T) {
	r := NewRegistry()
	d := &Dispatcher{Registry: r, Respond: func(context.Context, Message) (string, error) { return "x", nil }}
	in := Message{ChannelCode: "telegram", TenantID: "t", ExternalID: "e", ProviderMessageID: "msg-1", OccurredAt: time.Now(), Direction: DirectionIn, Text: "hi"}
	if err := d.Handle(context.Background(), in); !errors.Is(err, ErrUnknownChannel) {
		t.Fatalf("expected unknown channel, got %v", err)
	}
}
````

## 6. Configuration surface

El proyecto registra explícitamente cada adapter por código. Tenant/contacto/message ID/timestamp llegan desde la frontera autenticada del provider; no hay defaults ni secretos en este core.

## 7. Dependency bill

| Dependency | Pin | License | Purpose |
|---|---:|---|---|
| Go | 1.26.7 | BSD-3-Clause | envelopes, registry y dispatcher |

## 8. Apply order

Materializar antes del runtime y después registrar adapters live. Conectar receipts/reconciliación outbound antes de aceptar tráfico real.

## 9. Verification

Materializar 4/4, comparar hashes, ejecutar `gofmt` y tests. En el proyecto probar además webhook/provider real, duplicados, reordenamiento, reintentos, receipts outbound y reconciliación.

## 10. Reconstruction evidence

V236: cuatro archivos reconstruidos byte-exactos/gofmt; identidad ingress, delivery key estable/distinta y registry no vacío pasaron junto con el E2E app. Véase `reconstruction_evidence/GO_CONNECTED_CONVERSATION_RUNTIME_2026-09-04_V236.md`.
