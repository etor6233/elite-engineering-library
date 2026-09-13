# Go Agent Domain Binding

## 1. Metadata

```yaml
pack_id: "GO-AGENT-DOMAIN-BINDING"
pack_version: "0.4.1"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa un gateway HTTP que transforma tools autorizadas en comandos de cita, cotización, consulta y devolución con bearer actual, identidad idempotente estable y scope organización/lead resuelto por contacto."
stacks: ["Go 1.26.7", "HTTP/JSON"]
compatible_with: ["GO-CONNECTED-CONVERSATION-RUNTIME 0.1.x", "GO-FRANCHISE-CUSTOMER-JOURNEY-API 0.9.x", "GO-ENTERPRISE-QUERY-API 0.1.x"]
incompatible_with: ["credencial hardcodeada", "scope generado por el LLM", "efecto sin identidad ingress", "mapping comercial no aprobado"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://go.dev", "https://datatracker.ietf.org/doc/html/rfc6750", "https://sre.google/sre-book/reliable-product-launches/"]
verified_at: "2026-09-04"
```

Los bloques son `AUTHORED`. El código fue verificado contra APIs locales controladas; credenciales, IdP y backend live siguen condicionados al proyecto.

## 2. Applicability

Use para conectar el runtime a los owners existentes, sin duplicar CRM, agenda, pricing, orders o returns. Los mappings servicio/producto/precio se inyectan y los valores desconocidos fallan cerrados.

## 3. Architecture contract

- `Scope` dinámico proviene del resolver autorizado y exige organización+lead.
- `CommandIdentity` proviene del mensaje upstream; citas y cotizaciones no aceptan reintentos sin key/timestamp estable.
- El token se obtiene por request y se rechaza ausente, vacío o con CR/LF.
- Los defaults organización/lead son opcionales y sólo sirven para métodos heredados; ambos presentes o ambos ausentes. El runtime conectado usa exclusivamente métodos `*For`.
- Dos contactos con scopes diferentes no comparten paths ni bodies.

## 4. Exact file manifest

```text
CREATE internal/domainbind/quantity_test.go
CREATE internal/domainbind/config.go
CREATE internal/domainbind/gateway.go
CREATE internal/domainbind/gateway_test.go
```

## 5. Materialization blocks

### FILE: `internal/domainbind/config.go`
```yaml
block_id: "GO-AGENT-DOMAIN-BINDING:internal/domainbind/config.go:v4"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "eb65f4ba6bb392fe7bb985e8794fac7ee40a1c9980cbaa3451b04361aea1d578"
variables: []
secrets_allowed: false
```
````go
// Package domainbind binds the conversational agent's business intents to the
// franchise backend HTTP endpoints. Business-specific mappings (service→kind,
// product→variant, price book) are injected as configuration, never invented;
// unknown values fail closed.
package domainbind

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"
)

// Config carries the backend endpoint and the business mappings.
type Config struct {
	BaseURL         string            // backend root, e.g. http://localhost:8080
	TenantID        string            // exact tenant carried by the authenticated session
	TenantCode      string            // public tenant code
	OrganizationID  string            // optional legacy single-contact default
	LeadID          string            // optional legacy single-contact default
	ServiceKinds    map[string]string // "corte" -> "service"
	ProductVariants map[string]string // "scooter" -> "variant-1"
	PriceBookID     string
	ValidMinutes    int
	Timeout         time.Duration
	TokenProvider   AccessTokenProvider // current short-lived token source for protected routes
}

// AccessTokenProvider returns the current token at request time. Its
// implementation owns acquisition, refresh and secret storage.
type AccessTokenProvider interface {
	AccessToken(context.Context) (string, error)
}

// Validate enforces a safe, exact configuration.
func (c Config) Validate() error {
	u, err := url.Parse(c.BaseURL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return errors.New("domainbind: base url must be absolute http(s)")
	}
	if strings.TrimSpace(c.TenantID) == "" || strings.TrimSpace(c.TenantCode) == "" {
		return errors.New("domainbind: tenant and tenant code required")
	}
	if (strings.TrimSpace(c.OrganizationID) == "") != (strings.TrimSpace(c.LeadID) == "") {
		return errors.New("domainbind: legacy organization and lead defaults must be both present or both absent")
	}
	if len(c.ServiceKinds) == 0 || len(c.ProductVariants) == 0 || strings.TrimSpace(c.PriceBookID) == "" {
		return errors.New("domainbind: service kinds, product variants and price book required")
	}
	if c.ValidMinutes <= 0 || c.Timeout <= 0 {
		return errors.New("domainbind: valid minutes and timeout required")
	}
	return nil
}

// ResolveService returns the appointment kind for a service name, fail-closed.
func (c Config) ResolveService(service string) (string, error) {
	kind, ok := c.ServiceKinds[strings.ToLower(strings.TrimSpace(service))]
	if !ok || kind == "" {
		return "", fmt.Errorf("domainbind: unknown service %q", service)
	}
	return kind, nil
}

// ResolveProduct returns the variant id for a product name, fail-closed.
func (c Config) ResolveProduct(product string) (string, error) {
	variant, ok := c.ProductVariants[strings.ToLower(strings.TrimSpace(product))]
	if !ok || variant == "" {
		return "", fmt.Errorf("domainbind: unknown product %q", product)
	}
	return variant, nil
}
````

### FILE: `internal/domainbind/gateway.go`
```yaml
block_id: "GO-AGENT-DOMAIN-BINDING:internal/domainbind/gateway.go:v4"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "b4302bc524b5338ac4299f5011f6d0301ba73f1dcb64a1b0a1f2798713053261"
variables: []
secrets_allowed: false
```
````go
package domainbind

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"elite.local/enterprise/internal/agenttools"
)

var (
	ErrIdempotencyRequired = errors.New("domainbind: stable command identity required")
	ErrUnauthenticated     = errors.New("domainbind: protected operation requires an access token")
)

// CommandIdentity is supplied by the ingress owner, never by the LLM. The
// external message identity provides the key and timestamp, making retries
// byte-stable even when the customer used a relative date.
type CommandIdentity struct {
	IdempotencyKey string
	OccurredAt     time.Time
}

// Scope is resolved from the authenticated channel/contact mapping. It is not
// inferred from user text or generated by the model.
type Scope struct {
	OrganizationID string
	LeadID         string
}

func (s Scope) validate() error {
	if strings.TrimSpace(s.OrganizationID) == "" || len(s.OrganizationID) > 128 || strings.TrimSpace(s.LeadID) == "" || len(s.LeadID) > 128 {
		return errors.New("domainbind: organization and lead scope required")
	}
	return nil
}

func (c CommandIdentity) validate() error {
	if len(strings.TrimSpace(c.IdempotencyKey)) < 16 || len(c.IdempotencyKey) > 128 || c.OccurredAt.IsZero() {
		return ErrIdempotencyRequired
	}
	return nil
}

// Gateway implements agenttools.Domain over the franchise backend HTTP API.
type Gateway struct {
	cfg    Config
	client *http.Client
}

// NewGateway returns a validated gateway.
func NewGateway(cfg Config) (*Gateway, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return &Gateway{cfg: cfg, client: &http.Client{Timeout: cfg.Timeout}}, nil
}

func (g *Gateway) post(ctx context.Context, path string, body any, commandKey string, protected bool) error {
	raw, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(g.cfg.BaseURL, "/")+path, bytes.NewReader(raw))
	if err != nil {
		return err
	}
	req.Header.Set("content-type", "application/json")
	if commandKey != "" {
		req.Header.Set("Idempotency-Key", commandKey)
	}
	if protected {
		if err := g.authorize(ctx, req); err != nil {
			return err
		}
	}
	return g.do(req)
}

func (g *Gateway) get(ctx context.Context, path string, protected bool) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, strings.TrimRight(g.cfg.BaseURL, "/")+path, nil)
	if err != nil {
		return nil, err
	}
	if protected {
		if err := g.authorize(ctx, req); err != nil {
			return nil, err
		}
	}
	resp, err := g.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		data, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return nil, fmt.Errorf("domainbind: status %d: %s", resp.StatusCode, strings.TrimSpace(string(data)))
	}
	return io.ReadAll(resp.Body)
}

func (g *Gateway) authorize(ctx context.Context, req *http.Request) error {
	if g.cfg.TokenProvider == nil {
		return ErrUnauthenticated
	}
	token, err := g.cfg.TokenProvider.AccessToken(ctx)
	if err != nil {
		return fmt.Errorf("domainbind: access token: %w", err)
	}
	token = strings.TrimSpace(token)
	if token == "" || strings.ContainsAny(token, "\r\n") {
		return ErrUnauthenticated
	}
	req.Header.Set("Authorization", "Bearer "+token)
	return nil
}

func (g *Gateway) validateTenant(tenantID string) error {
	if strings.TrimSpace(tenantID) == "" || tenantID != g.cfg.TenantID {
		return errors.New("domainbind: tenant scope mismatch")
	}
	return nil
}

func (g *Gateway) do(req *http.Request) error {
	resp, err := g.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		data, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("domainbind: status %d: %s", resp.StatusCode, strings.TrimSpace(string(data)))
	}
	return nil
}

// BookAppointment books a public appointment for the resolved kind and time.
func (g *Gateway) BookAppointment(ctx context.Context, tenantID string, in agenttools.AppointmentInput) (string, error) {
	return "", ErrIdempotencyRequired
}

// BookAppointmentCommand books an appointment with ingress-owned replay identity.
func (g *Gateway) BookAppointmentCommand(ctx context.Context, tenantID string, command CommandIdentity, in agenttools.AppointmentInput) (string, error) {
	return g.BookAppointmentFor(ctx, tenantID, Scope{OrganizationID: g.cfg.OrganizationID, LeadID: g.cfg.LeadID}, command, in)
}

// BookAppointmentFor books an appointment for the resolved contact scope.
func (g *Gateway) BookAppointmentFor(ctx context.Context, tenantID string, scope Scope, command CommandIdentity, in agenttools.AppointmentInput) (string, error) {
	if err := g.validateTenant(tenantID); err != nil {
		return "", err
	}
	if err := scope.validate(); err != nil {
		return "", err
	}
	if err := command.validate(); err != nil {
		return "", err
	}
	kind, err := g.cfg.ResolveService(in.Service)
	if err != nil {
		return "", err
	}
	startsAt, err := parseWhen(in.When, command.OccurredAt)
	if err != nil {
		return "", err
	}
	path := fmt.Sprintf("/v1/public/%s/%s/appointments", urlEscape(g.cfg.TenantCode), urlEscape(scope.OrganizationID))
	body := map[string]any{"lead_id": scope.LeadID, "kind": kind, "starts_at": startsAt.UTC().Format(time.RFC3339)}
	if err := g.post(ctx, path, body, command.IdempotencyKey, false); err != nil {
		return "", err
	}
	return fmt.Sprintf("Turno solicitado para %s (%s).", in.Service, startsAt.Format("2006-01-02 15:04")), nil
}

// CreateQuote creates a franchise quote for the resolved variant.
func (g *Gateway) CreateQuote(ctx context.Context, tenantID string, in agenttools.QuoteInput) (string, error) {
	return "", ErrIdempotencyRequired
}

// CreateQuoteCommand creates a quote with ingress-owned replay identity.
func (g *Gateway) CreateQuoteCommand(ctx context.Context, tenantID string, command CommandIdentity, in agenttools.QuoteInput) (string, error) {
	return g.CreateQuoteFor(ctx, tenantID, Scope{OrganizationID: g.cfg.OrganizationID, LeadID: g.cfg.LeadID}, command, in)
}

// CreateQuoteFor creates a quote for the resolved contact scope.
func (g *Gateway) CreateQuoteFor(ctx context.Context, tenantID string, scope Scope, command CommandIdentity, in agenttools.QuoteInput) (string, error) {
	// The selected quotation owner represents one vehicle, not a quantity line.
	if in.Quantity != 1 {
		return "", errors.New("domainbind: reference quote requires exactly one vehicle")
	}
	if err := g.validateTenant(tenantID); err != nil {
		return "", err
	}
	if err := scope.validate(); err != nil {
		return "", err
	}
	if err := command.validate(); err != nil {
		return "", err
	}
	variant, err := g.cfg.ResolveProduct(in.Product)
	if err != nil {
		return "", err
	}
	validUntil := command.OccurredAt.UTC().Add(time.Duration(g.cfg.ValidMinutes) * time.Minute)
	body := map[string]any{
		"organization_id": scope.OrganizationID,
		"lead_id":         scope.LeadID,
		"variant_id":      variant,
		"price_book_id":   g.cfg.PriceBookID,
		"valid_until":     validUntil.Format(time.RFC3339),
	}
	if err := g.post(ctx, "/v1/franchise/quotes", body, command.IdempotencyKey, true); err != nil {
		return "", err
	}
	return fmt.Sprintf("Cotización creada para %s.", in.Product), nil
}

// OrderStatus returns the state of the customer's orders matching orderID.
func (g *Gateway) OrderStatus(ctx context.Context, tenantID, orderID string) (string, error) {
	return g.OrderStatusFor(ctx, tenantID, Scope{OrganizationID: g.cfg.OrganizationID, LeadID: g.cfg.LeadID}, orderID)
}

// OrderStatusFor returns order status within the resolved organization scope.
func (g *Gateway) OrderStatusFor(ctx context.Context, tenantID string, scope Scope, orderID string) (string, error) {
	if err := g.validateTenant(tenantID); err != nil {
		return "", err
	}
	if err := scope.validate(); err != nil {
		return "", err
	}
	path := fmt.Sprintf("/v1/customer/orders?organization_id=%s", urlEscape(scope.OrganizationID))
	data, err := g.get(ctx, path, true)
	if err != nil {
		return "", err
	}
	var page struct {
		Items []struct {
			OrderID string `json:"order_id"`
			State   string `json:"state"`
		} `json:"items"`
	}
	if err := json.Unmarshal(data, &page); err != nil {
		return "", err
	}
	for _, item := range page.Items {
		if item.OrderID == orderID {
			return fmt.Sprintf("Pedido %s: %s.", orderID, item.State), nil
		}
	}
	return "", fmt.Errorf("domainbind: order %q not found", orderID)
}

// RequestReturn rejects the handover for the order, starting the governed
// return flow. It never executes money/inventory directly.
func (g *Gateway) RequestReturn(ctx context.Context, tenantID string, in agenttools.ReturnInput) (string, error) {
	return g.RequestReturnFor(ctx, tenantID, Scope{OrganizationID: g.cfg.OrganizationID, LeadID: g.cfg.LeadID}, in)
}

// RequestReturnFor starts the governed return flow for a resolved scope.
func (g *Gateway) RequestReturnFor(ctx context.Context, tenantID string, scope Scope, in agenttools.ReturnInput) (string, error) {
	if err := g.validateTenant(tenantID); err != nil {
		return "", err
	}
	if err := scope.validate(); err != nil {
		return "", err
	}
	path := fmt.Sprintf("/v1/customer/journey?organization_id=%s", urlEscape(scope.OrganizationID))
	data, err := g.get(ctx, path, true)
	if err != nil {
		return "", err
	}
	var journey struct {
		Handovers []struct {
			ID      string `json:"id"`
			OrderID string `json:"order_id"`
			Version int64  `json:"version"`
		} `json:"handovers"`
	}
	if err := json.Unmarshal(data, &journey); err != nil {
		return "", err
	}
	for _, h := range journey.Handovers {
		if h.OrderID == in.OrderID {
			body := map[string]any{
				"organization_id": scope.OrganizationID,
				"version":         h.Version,
				"reason_code":     "customer-return",
				"details":         in.Reason,
			}
			if err := g.post(ctx, fmt.Sprintf("/v1/customer/handovers/%s/reject", urlEscape(h.ID)), body, "", true); err != nil {
				return "", err
			}
			return "Devolución solicitada; quedará gobernada por aprobación y los workers de devolución.", nil
		}
	}
	return "", fmt.Errorf("domainbind: no entregable para el pedido %q", in.OrderID)
}

func urlEscape(s string) string {
	return url.PathEscape(s)
}

func parseWhen(when string, now time.Time) (time.Time, error) {
	s := strings.TrimSpace(when)
	switch strings.ToLower(s) {
	case "manana", "mañana":
		return now.Add(24 * time.Hour), nil
	case "hoy":
		return now, nil
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Time{}, errors.New("domainbind: when must be an ISO timestamp or hoy/manana")
	}
	if t.Before(now) {
		return time.Time{}, errors.New("domainbind: when must be in the future")
	}
	return t, nil
}
````

### FILE: `internal/domainbind/gateway_test.go`
```yaml
block_id: "GO-AGENT-DOMAIN-BINDING:internal/domainbind/gateway_test.go:v4"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "3626783c467fb0a5ed84a08bf23995ec3665435c61b7e1ad682986f50fa01f2d"
variables: []
secrets_allowed: false
```
````go
package domainbind

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"elite.local/enterprise/internal/agenttools"
)

type tokenProvider struct{ token string }

func (p *tokenProvider) AccessToken(context.Context) (string, error) { return p.token, nil }

func testConfig(serverURL string) Config {
	return Config{
		BaseURL:         serverURL,
		TenantID:        "tenant-id",
		TenantCode:      "acme",
		OrganizationID:  "store",
		LeadID:          "lead-1",
		ServiceKinds:    map[string]string{"corte": "service", "prueba": "test-drive"},
		ProductVariants: map[string]string{"scooter": "variant-1"},
		PriceBookID:     "retail",
		ValidMinutes:    60,
		Timeout:         5 * time.Second,
		TokenProvider:   &tokenProvider{token: "token-1"},
	}
}

func TestConfigValidateAndResolve(t *testing.T) {
	cfg := testConfig("https://backend.example")
	if err := cfg.Validate(); err != nil {
		t.Fatalf("valid config rejected: %v", err)
	}
	if _, err := cfg.ResolveService("desconocido"); err == nil {
		t.Fatal("unknown service accepted")
	}
	if _, err := cfg.ResolveProduct("desconocido"); err == nil {
		t.Fatal("unknown product accepted")
	}
	noLead := testConfig("https://backend.example")
	noLead.LeadID = ""
	if err := noLead.Validate(); err == nil {
		t.Fatal("missing lead accepted")
	}
	dynamic := testConfig("https://backend.example")
	dynamic.OrganizationID, dynamic.LeadID = "", ""
	if err := dynamic.Validate(); err != nil {
		t.Fatalf("dynamic-scope config rejected: %v", err)
	}
}

func TestBookAppointmentRoundTrip(t *testing.T) {
	occurredAt := time.Date(2026, 9, 4, 12, 0, 0, 0, time.UTC)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/public/acme/store/appointments" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)
		if body["kind"] != "service" || body["lead_id"] != "lead-1" {
			t.Errorf("unexpected body: %+v", body)
		}
		if r.Header.Get("Idempotency-Key") != "meta-message-0001" {
			t.Errorf("missing idempotency key: %q", r.Header.Get("Idempotency-Key"))
		}
		if body["starts_at"] != occurredAt.Add(24*time.Hour).Format(time.RFC3339) {
			t.Errorf("relative time is not retry-stable: %+v", body)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":"a1","state":"requested"}`))
	}))
	defer srv.Close()

	g, err := NewGateway(testConfig(srv.URL))
	if err != nil {
		t.Fatal(err)
	}
	out, err := g.BookAppointmentCommand(context.Background(), "tenant-id", CommandIdentity{IdempotencyKey: "meta-message-0001", OccurredAt: occurredAt}, agenttools.AppointmentInput{Service: "corte", When: "manana"})
	if err != nil {
		t.Fatal(err)
	}
	if out == "" {
		t.Fatal("empty confirmation")
	}
}

func TestBookAppointmentUnknownService(t *testing.T) {
	g, _ := NewGateway(testConfig("https://backend.example"))
	if _, err := g.BookAppointmentCommand(context.Background(), "tenant-id", CommandIdentity{IdempotencyKey: "meta-message-0002", OccurredAt: time.Now()}, agenttools.AppointmentInput{Service: "spa", When: "manana"}); err == nil {
		t.Fatal("unknown service accepted")
	}
	if _, err := g.BookAppointment(context.Background(), "tenant-id", agenttools.AppointmentInput{Service: "corte", When: "manana"}); !errors.Is(err, ErrIdempotencyRequired) {
		t.Fatal("legacy appointment path did not fail closed")
	}
}

func TestCreateQuoteRoundTrip(t *testing.T) {
	provider := &tokenProvider{token: "token-1"}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/franchise/quotes" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)
		if body["variant_id"] != "variant-1" || body["price_book_id"] != "retail" {
			t.Errorf("unexpected body: %+v", body)
		}
		if r.Header.Get("Idempotency-Key") != "meta-message-quote-1" || r.Header.Get("Authorization") != "Bearer "+provider.token {
			t.Errorf("missing command identity or current bearer")
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":"q1","state":"draft"}`))
	}))
	defer srv.Close()

	cfg := testConfig(srv.URL)
	cfg.TokenProvider = provider
	g, _ := NewGateway(cfg)
	command := CommandIdentity{IdempotencyKey: "meta-message-quote-1", OccurredAt: time.Date(2026, 9, 4, 12, 0, 0, 0, time.UTC)}
	if _, err := g.CreateQuoteCommand(context.Background(), "tenant-id", command, agenttools.QuoteInput{Product: "scooter", Quantity: 1}); err != nil {
		t.Fatal(err)
	}
	provider.token = "token-2"
	if _, err := g.CreateQuoteCommand(context.Background(), "tenant-id", command, agenttools.QuoteInput{Product: "scooter", Quantity: 1}); err != nil {
		t.Fatal(err)
	}
}

func TestOrderStatusRoundTrip(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/customer/orders" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"items":[{"order_id":"ORD-1","state":"delivered"}]}`))
	}))
	defer srv.Close()

	g, _ := NewGateway(testConfig(srv.URL))
	out, err := g.OrderStatus(context.Background(), "tenant-id", "ORD-1")
	if err != nil {
		t.Fatal(err)
	}
	if out != "Pedido ORD-1: delivered." {
		t.Fatalf("unexpected: %q", out)
	}
}

func TestRequestReturnRejectsHandover(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/customer/journey":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"handovers":[{"id":"h1","order_id":"ORD-1","version":2}]}`))
		case "/v1/customer/handovers/h1/reject":
			if r.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", r.Method)
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"state":"rejected"}`))
		default:
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer srv.Close()

	g, _ := NewGateway(testConfig(srv.URL))
	if _, err := g.RequestReturn(context.Background(), "tenant-id", agenttools.ReturnInput{OrderID: "ORD-1", Reason: "no sirve"}); err != nil {
		t.Fatal(err)
	}
}

func TestDynamicContactScopeIsNotShared(t *testing.T) {
	seen := []string{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)
		seen = append(seen, r.URL.Path+":"+body["lead_id"].(string))
		w.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()
	g, _ := NewGateway(testConfig(srv.URL))
	command := CommandIdentity{IdempotencyKey: "provider-message-0001", OccurredAt: time.Date(2026, 9, 4, 12, 0, 0, 0, time.UTC)}
	if _, err := g.BookAppointmentFor(context.Background(), "tenant-id", Scope{OrganizationID: "north", LeadID: "lead-a"}, command, agenttools.AppointmentInput{Service: "corte", When: "manana"}); err != nil {
		t.Fatal(err)
	}
	command.IdempotencyKey = "provider-message-0002"
	if _, err := g.BookAppointmentFor(context.Background(), "tenant-id", Scope{OrganizationID: "south", LeadID: "lead-b"}, command, agenttools.AppointmentInput{Service: "corte", When: "manana"}); err != nil {
		t.Fatal(err)
	}
	if len(seen) != 2 || seen[0] != "/v1/public/acme/north/appointments:lead-a" || seen[1] != "/v1/public/acme/south/appointments:lead-b" {
		t.Fatalf("contact scopes crossed: %v", seen)
	}
	if _, err := g.BookAppointmentFor(context.Background(), "tenant-id", Scope{OrganizationID: "north"}, command, agenttools.AppointmentInput{Service: "corte", When: "manana"}); err == nil {
		t.Fatal("incomplete contact scope accepted")
	}
}

func TestProtectedAndTenantScopeFailClosed(t *testing.T) {
	cfg := testConfig("https://backend.example")
	cfg.TokenProvider = nil
	g, _ := NewGateway(cfg)
	command := CommandIdentity{IdempotencyKey: "meta-message-quote-2", OccurredAt: time.Now()}
	if _, err := g.CreateQuoteCommand(context.Background(), "tenant-id", command, agenttools.QuoteInput{Product: "scooter", Quantity: 1}); !errors.Is(err, ErrUnauthenticated) {
		t.Fatalf("missing token provider not rejected: %v", err)
	}
	if _, err := g.OrderStatus(context.Background(), "other-tenant", "ORD-1"); err == nil {
		t.Fatal("cross-tenant call accepted")
	}
}
````

## 6. Configuration surface

`DOMAIN_BASE_URL`, tenant code, token provider, service/product mappings, price book, timeouts and quote validity are explicit project inputs. Legacy organization/lead defaults are optional as a pair and are not used by the connected runtime.

## 7. Dependency bill

| Dependency | Pin | License | Purpose |
|---|---:|---|---|
| Go | 1.26.7 | BSD-3-Clause | HTTP/JSON gateway and tests |

## 8. Apply order

Materialize conversation tools first, then this gateway, the connected runtime and app wiring. Configure mappings and rotating token provider before exposure.

## 9. Verification

Materializar 3/3, comparar hashes, ejecutar `gofmt`, tests de HTTP/token/replay/scopes y la suite completa. Revalidar rutas/contratos contra el backend materializado y probar IdP real antes de promoción productiva.

## 10. Reconstruction evidence

V236: tres archivos reconstruidos byte-exactos y gofmt-idempotentes; dos contactos conservaron scopes diferentes; bearer actual y claves/timestamps estables pasaron en tests focales y en la suite de 48 paquetes. Véase `reconstruction_evidence/GO_CONNECTED_CONVERSATION_RUNTIME_2026-09-04_V236.md`.

V402 composed delta: T2807 connected reference315: real domain quotation and contact-bound order status, explicit single-vehicle limit, revalidated contact, structured history roles, per-case required eval gates and canonical host mounting. Local fixtures only. AI_CONNECTED_REFERENCE_RELEASE_V402.md/json. No new upstream dependency or live model quality claim.

### FILE: `internal/domainbind/quantity_test.go`

```yaml
block_id: "GO-AGENT-DOMAIN-BINDING-CONNECTED-REFERENCE:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "c3d4e439bff82189076ed7e89d82351a6ab52dd3492b88c23ac0ff40b54cd87a"
variables: []
secrets_allowed: false
```

````go
package domainbind

import (
	"context"
	"elite.local/enterprise/internal/agenttools"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestQuoteQuantityCannotBeSilentlyDiscarded(t *testing.T) {
	calls := 0
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls++; w.WriteHeader(201) }))
	defer s.Close()
	g, e := NewGateway(testConfig(s.URL))
	if e != nil {
		t.Fatal(e)
	}
	for _, n := range []int{0, -1, 2, 100} {
		if _, e = g.CreateQuoteCommand(context.Background(), "tenant-id", CommandIdentity{IdempotencyKey: "synthetic-quantity-key", OccurredAt: time.Now()}, agenttools.QuoteInput{Product: "scooter", Quantity: n}); e == nil {
			t.Fatal("quantity accepted", n)
		}
	}
	if calls != 0 {
		t.Fatal("invalid quantities caused effects", calls)
	}
}
````


V402315: existing owners compose a single local AI runtime and required-case evaluation. Read docs/AI_REFERENCE_START.md. Historical SFT uses its separate admitted opt-in plan and exact later execution hash; no new training engine or inferred private model quality.
