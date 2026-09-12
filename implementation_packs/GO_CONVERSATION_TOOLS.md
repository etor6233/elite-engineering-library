# Go Conversation Tools

## 1. Metadata

```yaml
pack_id: "GO-CONVERSATION-TOOLS"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa el contrato de negocio agnóstico del chatbot (citas, cotización/venta, estado de pedido, devolución) y las tools que lo bindean al agente por intención, sin que el bot toque dinero/inventario directamente."
stacks: ["Go 1.26.7"]
compatible_with: ["GO-CONVERSATIONAL-AGENT 0.1.x", "GO-FRANCHISE-CUSTOMER-JOURNEY-API 0.8.x", "GO-COMMERCE-PRICING-PAYMENT-API 0.3.x"]
incompatible_with: ["efectos de dinero/inventario dentro del bot", "tool sin intención autorizada"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://go.dev"]
verified_at: "2026-09-02"
```

Este pack hace al chatbot **adaptable a cualquier negocio**: define el contrato `Domain` (cuatro operaciones) que cualquier franquicia implementa, y las tools deterministas que bindean cada intención a ese contrato. El parser `clave: valor` es el camino determinista; convertir texto libre a estructura es trabajo del LLM (CONDITIONED). El dinero, el inventario y la agenda nunca los decide el bot.

## 2. Applicability

Use este pack para conectar el agente conversacional a citas, cotización/venta, estado de pedido y devolución de un negocio concreto, sin acoplar el bot a un dominio.

Rechace este pack para: un bot que ejecute efectos de dinero/inventario por sí mismo; o un dominio sin implementar el contrato `Domain`.

## 3. Architecture contract

- **Ownership**: `internal/agenttools` gobierna el contrato `Domain` y las tools; `internal/agent` gobierna el ruteo y la autorización por intención.
- **Invariantes**: (1) una intención → una tool → un método del `Domain`. (2) campos requeridos ausentes → `ErrNeedsInfo` (nunca se inventa el dato). (3) cantidad inválida → `ErrNeedsInfo`. (4) el bot delega, no ejecuta efectos.
- **Data flow**: `Agent.Process` → intent → `agenttools.Tool.Run` → `parseKV` → `Domain.<operación>` → confirmación.
- **Failure modes**: campo ausente → `ErrNeedsInfo`; dominio falla → error propaga (el agente deriva a humano).
- **Seguridad/privacidad**: tenant-scoped; sin PII adicional; evidencia del lado del dominio.
- **Performance budget**: O(1) por tool; parser determinista.
- **Operación/migración/rollback**: sin migración; reemplazar `Domain` es composición.

## 4. Exact file manifest

```text
CREATE internal/agenttools/domain.go
CREATE internal/agenttools/parse.go
CREATE internal/agenttools/tools.go
CREATE internal/agenttools/toolkit.go
CREATE internal/agenttools/agenttools_test.go
```

## 5. Materialization blocks

### FILE: `internal/agenttools/domain.go`
```yaml
block_id: "GO-CONVERSATION-TOOLS:internal/agenttools/domain.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "b4e134e2843f3202f22a5c66531f0dc69db729445400490c9640086f4f264f0c"
variables: []
secrets_allowed: false
```
````go
// Package agenttools binds the conversational agent's intents to a
// business-agnostic Domain contract (appointments, quotes, order status,
// returns). The chatbot never touches money/inventory directly: it delegates
// to the domain, which any franchise implements.
package agenttools

import "context"

// Domain is the business-agnostic capability contract the chatbot delegates
// to. Any franchise implements these four operations.
type Domain interface {
	BookAppointment(ctx context.Context, tenantID string, in AppointmentInput) (string, error)
	CreateQuote(ctx context.Context, tenantID string, in QuoteInput) (string, error)
	OrderStatus(ctx context.Context, tenantID, orderID string) (string, error)
	RequestReturn(ctx context.Context, tenantID string, in ReturnInput) (string, error)
}

// AppointmentInput carries a booking request.
type AppointmentInput struct {
	Service string
	When    string
}

// QuoteInput carries a quote request.
type QuoteInput struct {
	Product  string
	Quantity int
}

// ReturnInput carries a return request.
type ReturnInput struct {
	OrderID string
	Reason  string
}
````

### FILE: `internal/agenttools/parse.go`
```yaml
block_id: "GO-CONVERSATION-TOOLS:internal/agenttools/parse.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "12fb9fc6aefd5cbaf884683bf543c310330cc00ba5bec50f408af8da8c8fa3a2"
variables: []
secrets_allowed: false
```
````go
package agenttools

import (
	"regexp"
	"strings"
)

var kvRe = regexp.MustCompile(`(?i)([a-záéíóúñ]+):\s*([^,\n;]+)`)

// parseKV extracts deterministic "key: value" pairs (case-insensitive keys).
// This is the deterministic path; converting free text to structured input is
// the LLM's job (CONDITIONED).
func parseKV(text string) map[string]string {
	out := map[string]string{}
	for _, m := range kvRe.FindAllStringSubmatch(text, -1) {
		k := strings.ToLower(m[1])
		v := strings.TrimSpace(m[2])
		if v != "" {
			out[k] = v
		}
	}
	return out
}
````

### FILE: `internal/agenttools/tools.go`
```yaml
block_id: "GO-CONVERSATION-TOOLS:internal/agenttools/tools.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "4eca0fc21c16602dae6f778b0ca78d280ed310d158ce2b7b6b458cbfefacea7a"
variables: []
secrets_allowed: false
```
````go
package agenttools

import (
	"context"
	"strconv"

	"elite.local/enterprise/internal/agent"
)

type appointmentTool struct{ d Domain }

func (appointmentTool) Name() string         { return "book-appointment" }
func (appointmentTool) Intent() agent.Intent { return agent.IntentAppointment }
func (t appointmentTool) Run(ctx context.Context, tenantID, text string) (string, error) {
	kv := parseKV(text)
	svc, hasSvc := kv["servicio"]
	when, hasWhen := kv["cuando"]
	if !hasSvc || !hasWhen {
		return "Necesito servicio y cuándo (ej. servicio: corte, cuando: mañana).", agent.ErrNeedsInfo
	}
	return t.d.BookAppointment(ctx, tenantID, AppointmentInput{Service: svc, When: when})
}

type quoteTool struct{ d Domain }

func (quoteTool) Name() string         { return "create-quote" }
func (quoteTool) Intent() agent.Intent { return agent.IntentQuote }
func (t quoteTool) Run(ctx context.Context, tenantID, text string) (string, error) {
	kv := parseKV(text)
	product, has := kv["producto"]
	if !has {
		return "Necesito el producto (ej. producto: scooter).", agent.ErrNeedsInfo
	}
	qty := 1
	if q, ok := kv["cantidad"]; ok {
		n, err := strconv.Atoi(q)
		if err != nil || n <= 0 {
			return "Cantidad inválida.", agent.ErrNeedsInfo
		}
		qty = n
	}
	return t.d.CreateQuote(ctx, tenantID, QuoteInput{Product: product, Quantity: qty})
}

type orderStatusTool struct{ d Domain }

func (orderStatusTool) Name() string         { return "order-status" }
func (orderStatusTool) Intent() agent.Intent { return agent.IntentOrderStatus }
func (t orderStatusTool) Run(ctx context.Context, tenantID, text string) (string, error) {
	kv := parseKV(text)
	order, has := kv["pedido"]
	if !has {
		return "Necesito el número de pedido (ej. pedido: ORD-123).", agent.ErrNeedsInfo
	}
	return t.d.OrderStatus(ctx, tenantID, order)
}

type returnTool struct{ d Domain }

func (returnTool) Name() string         { return "request-return" }
func (returnTool) Intent() agent.Intent { return agent.IntentReturnRequest }
func (t returnTool) Run(ctx context.Context, tenantID, text string) (string, error) {
	kv := parseKV(text)
	order, hasOrder := kv["pedido"]
	reason, hasReason := kv["motivo"]
	if !hasOrder || !hasReason {
		return "Necesito pedido y motivo (ej. pedido: ORD-1, motivo: no sirve).", agent.ErrNeedsInfo
	}
	return t.d.RequestReturn(ctx, tenantID, ReturnInput{OrderID: order, Reason: reason})
}
````

### FILE: `internal/agenttools/toolkit.go`
```yaml
block_id: "GO-CONVERSATION-TOOLS:internal/agenttools/toolkit.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "aa1a9d9c53e8c914a9ca6089bdebe1844ffc93dd3dcb13eb815f38de88cb4b43"
variables: []
secrets_allowed: false
```
````go
package agenttools

import "elite.local/enterprise/internal/agent"

// Toolkit wires the four business tools into a ToolRegistry (one per intent).
func Toolkit(d Domain) *agent.ToolRegistry {
	r := agent.NewToolRegistry()
	_ = r.Register(appointmentTool{d: d})
	_ = r.Register(quoteTool{d: d})
	_ = r.Register(orderStatusTool{d: d})
	_ = r.Register(returnTool{d: d})
	return r
}
````

### FILE: `internal/agenttools/agenttools_test.go`
```yaml
block_id: "GO-CONVERSATION-TOOLS:internal/agenttools/agenttools_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "3dac1f760ff224e33afe826b823426892a29cf849a2ce90d58b69b9eebad439e"
variables: []
secrets_allowed: false
```
````go
package agenttools

import (
	"context"
	"errors"
	"testing"

	"elite.local/enterprise/internal/agent"
)

type fakeDomain struct {
	appt  AppointmentInput
	quote QuoteInput
	order string
	ret   ReturnInput
}

func (f *fakeDomain) BookAppointment(_ context.Context, _ string, in AppointmentInput) (string, error) {
	f.appt = in
	return "turno confirmado", nil
}
func (f *fakeDomain) CreateQuote(_ context.Context, _ string, in QuoteInput) (string, error) {
	f.quote = in
	return "cotización lista", nil
}
func (f *fakeDomain) OrderStatus(_ context.Context, _ string, orderID string) (string, error) {
	f.order = orderID
	return "en camino", nil
}
func (f *fakeDomain) RequestReturn(_ context.Context, _ string, in ReturnInput) (string, error) {
	f.ret = in
	return "devolución solicitada", nil
}

func TestToolkitWiresFourIntents(t *testing.T) {
	r := Toolkit(&fakeDomain{})
	for _, intent := range []agent.Intent{
		agent.IntentAppointment, agent.IntentQuote, agent.IntentOrderStatus, agent.IntentReturnRequest,
	} {
		if _, ok := r.For(intent); !ok {
			t.Fatalf("missing tool for intent %q", intent)
		}
	}
}

func TestAppointmentToolExtractsAndCalls(t *testing.T) {
	d := &fakeDomain{}
	tool, _ := Toolkit(d).For(agent.IntentAppointment)
	out, err := tool.Run(context.Background(), "t", "servicio: corte, cuando: mañana")
	if err != nil || out != "turno confirmado" {
		t.Fatalf("unexpected: %q/%v", out, err)
	}
	if d.appt.Service != "corte" || d.appt.When != "mañana" {
		t.Fatalf("wrong args: %+v", d.appt)
	}
}

func TestAppointmentToolNeedsInfo(t *testing.T) {
	tool, _ := Toolkit(&fakeDomain{}).For(agent.IntentAppointment)
	if _, err := tool.Run(context.Background(), "t", "quiero un turno"); !errors.Is(err, agent.ErrNeedsInfo) {
		t.Fatalf("expected ErrNeedsInfo, got %v", err)
	}
}

func TestQuoteToolQuantity(t *testing.T) {
	d := &fakeDomain{}
	tool, _ := Toolkit(d).For(agent.IntentQuote)
	if _, err := tool.Run(context.Background(), "t", "producto: scooter, cantidad: 2"); err != nil {
		t.Fatal(err)
	}
	if d.quote.Product != "scooter" || d.quote.Quantity != 2 {
		t.Fatalf("wrong quote args: %+v", d.quote)
	}
}

func TestQuoteToolInvalidQuantity(t *testing.T) {
	tool, _ := Toolkit(&fakeDomain{}).For(agent.IntentQuote)
	if _, err := tool.Run(context.Background(), "t", "producto: x, cantidad: cero"); !errors.Is(err, agent.ErrNeedsInfo) {
		t.Fatalf("expected ErrNeedsInfo for bad quantity, got %v", err)
	}
}

func TestOrderStatusAndReturn(t *testing.T) {
	d := &fakeDomain{}
	tool, _ := Toolkit(d).For(agent.IntentOrderStatus)
	if _, err := tool.Run(context.Background(), "t", "pedido: ORD-9"); err != nil {
		t.Fatal(err)
	}
	if d.order != "ORD-9" {
		t.Fatalf("wrong order id: %q", d.order)
	}

	ret, _ := Toolkit(d).For(agent.IntentReturnRequest)
	if _, err := ret.Run(context.Background(), "t", "pedido: ORD-9, motivo: no sirve"); err != nil {
		t.Fatal(err)
	}
	if d.ret.OrderID != "ORD-9" || d.ret.Reason != "no sirve" {
		t.Fatalf("wrong return args: %+v", d.ret)
	}
}
````


## 6. Configuration surface

Sin variables ni secretos. El binding del `Domain` real (endpoints Go existentes) es composición del proyecto.

## 7. Dependency bill

| Package/image/tool | Pin exacto | Uso | Licencia | Runtime/build | Fuente oficial |
|---|---|---|---|---|---|
| Go standard library | go 1.26.7 (toolchain) | núcleo | BSD-3-Clause | runtime | https://go.dev |
| internal/agent | 0.1.0 | intents/tools/registry | LicenseRef-Workspace-Owner | runtime | este repositorio |

## 8. Apply order

1. Componer `GO-CONVERSATIONAL-AGENT` (mismo módulo).
2. Colocar los cinco archivos bajo `internal/agenttools/`.
3. Implementar `Domain` con los endpoints Go reales y registrar `Toolkit(domain)` en el agente.
4. Verificar con `go test ./... -count=1` y `go vet ./...`.
5. Rollback: eliminar `internal/agenttools/`; no deja estado.

## 9. Verification

- `go test ./internal/agenttools/ -count=1`: 6/6 PASS (toolkit con 4 intents, extracción y llamada de cita, `ErrNeedsInfo` por campo ausente, cantidad parseada e inválida, estado de pedido y devolución).
- `go test ./... -count=1` (8 paquetes): PASS.
- `go vet ./...`: exit 0.

## 10. Reconstruction evidence

Véase `reconstruction_evidence/GO_CONVERSATION_TOOLS_2026-09-02_V183.md`.
