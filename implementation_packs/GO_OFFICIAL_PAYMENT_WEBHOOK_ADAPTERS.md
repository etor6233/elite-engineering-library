# Go Official Payment Webhook Adapters

## 1. Metadata

```yaml
pack_id: "GO-OFFICIAL-PAYMENT-WEBHOOK-ADAPTERS"
pack_version: "0.3.2"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "SDKs oficiales fijados para creaciÃ³n idempotente, GET de reconciliaciÃ³n y notificaciones autenticadas Stripe/Mercado Pago; normalizaciÃ³n local acotada sin decidir cobro ni liberaciÃ³n del dominio."
stacks: ["Go 1.26.8", "stripe-go/v86 86.3.0", "mercadopago/sdk-go 1.14.0"]
compatible_with: ["GO-PROVIDER-INTEGRATION-CORE 0.1.x", "GO-ENTERPRISE-BACKEND 0.4.x", "OFFICIAL-UPSTREAM-ACQUISITION-CORE 0.4.x"]
incompatible_with: ["body transformado antes de firma", "secretos en cÃ³digo", "QR Mercado Pago no firmado", "webhooks sin inbox/reconciliation"]
license_expression: "LicenseRef-Workspace-Owner AND MIT"
upstream_sources: ["https://github.com/stripe/stripe-go/tree/a2df585a800a97fe8ec4ebf551b4449bdb3d90a1", "https://github.com/mercadopago/sdk-go/tree/f910ee53fbb6819e435eaf3d0f800cb1fe74ae09"]
verified_at: "2026-09-11"
```

## 2. Applicability

Use cuando el proyecto haya seleccionado Stripe o Mercado Pago y necesite autenticar webhooks Go antes del inbox durable. Requiere recibir el body crudo sin transformaciÃ³n, secretos por referencia, headers/query originales, lÃ­mites HTTP y una polÃ­tica de tolerancia positiva. Rechazar para QR Mercado Pago â€”el SDK oficial declara que esas notificaciones no estÃ¡n firmadasâ€”, para capturar/reembolsar pagos sin workflow seleccionado o cuando el proveedor/versiÃ³n no coincida.

El glue de normalizaciÃ³n/outbound es `AUTHORED`; la criptografÃ­a, parsing de signature, API-version enforcement y serializaciÃ³n HTTP se delegan a los SDKs oficiales. El outbound crea la request oficial PaymentIntent/Payment con idempotency y metadata de orden; el GET oficial devuelve observaciÃ³n mÃ­nima para reconciliaciÃ³n. Las nuevas normalizaciones no mutan el ledger. El claim no incluye cuenta, cobro confirmado, PCI, OAuth, settlement, chargeback, refund, paÃ­s/producto, sandbox ni conciliaciÃ³n.

## 3. Architecture contract

La frontera acepta bytes inmutables, metadata del edge y un secreto resuelto en runtime; primero limita tamaÃ±o y entradas, despuÃ©s invoca el verifier oficial y reciÃ©n entonces parsea/normaliza. Stripe conserva el API-version check y tolerancia del SDK. Mercado Pago conserva su manifest/HMAC/tolerancia oficial y aÃ±ade un cross-check entre `data.id` firmado en query y `data.id` del body cuando existe.

Nunca persiste ni registra secretos. Devuelve una copia del payload autenticado para que el caller escriba atÃ³micamente su inbox; no confirma procesamiento al proveedor por sÃ­ mismo. Fallos producen error cerrado y ningÃºn efecto externo. Performance budget: body mÃ¡ximo 1 MiB y verificaciÃ³n local sub-10 ms p95 en hardware del target a demostrar. Rollback: retirar wiring del adapter y conservar raw inbox/reconciliation; actualizar sÃ³lo por dependency update con SDK/source/hash/evidencia nuevos.

## 4. Exact file manifest

```text
CREATE official_payment_webhooks/go.mod
CREATE official_payment_webhooks/go.sum
CREATE official_payment_webhooks/officialpayments/client.go
CREATE official_payment_webhooks/officialpayments/client_test.go
CREATE official_payment_webhooks/officialpayments/verifier.go
CREATE official_payment_webhooks/officialpayments/verifier_test.go
CREATE official_payment_webhooks/officialpayments/retrieve.go
CREATE official_payment_webhooks/officialpayments/notification.go
CREATE official_payment_webhooks/officialpayments/retrieve_test.go
CREATE official_payment_webhooks/officialpayments/checkout.go
CREATE official_payment_webhooks/officialpayments/checkout_test.go
CREATE official_payment_webhooks/README.md
CREATE official_payment_webhooks/officialpayments/preference_recovery.go
```

## 5. Materialization blocks

### FILE: `official_payment_webhooks/go.mod`
```yaml
block_id: "GO-OFFICIAL-PAYMENT-WEBHOOKS:gomod:v1"
operation: CREATE
provenance: AUTHORED
source: "dependency pins for official provider SDKs"
license: "LicenseRef-Workspace-Owner"
sha256: "fc6106ab8e8e4ab2dde71a6210e5d4b31cc674e62088947e29eb8e025387b141"
variables: []
secrets_allowed: false
```
````go
module example.com/elite/official-payment-webhooks

go 1.26.0

require (
	github.com/mercadopago/sdk-go v1.14.0
	github.com/stripe/stripe-go/v86 v86.3.0
)

require github.com/google/uuid v1.6.0 // indirect
````

### FILE: `official_payment_webhooks/go.sum`
```yaml
block_id: "GO-OFFICIAL-PAYMENT-WEBHOOKS:gosum:v1"
operation: CREATE
provenance: AUTHORED
source: "Go checksum database resolution for exact module pins"
license: "LicenseRef-Workspace-Owner"
sha256: "6f564169198e386498d553a0c9a1e7bce16b619d4e88d640429ad1f3ff74dc06"
variables: []
secrets_allowed: false
```
````text
github.com/davecgh/go-spew v1.1.0 h1:ZDRjVQ15GmhC3fiQ8ni8+OwkZQO4DARzQgrnXU1Liz8=
github.com/davecgh/go-spew v1.1.0/go.mod h1:J7Y8YcW2NihsgmVo/mv3lAwl/skON4iLHjSsI+c5H38=
github.com/google/uuid v1.6.0 h1:NIvaJDMOsjHA8n1jAhLSgzrAzy1Hgr+hNrb57e+94F0=
github.com/google/uuid v1.6.0/go.mod h1:TIyPZe4MgqvfeYDBFedMoGGpEw/LqOeaOT+nhxU+yHo=
github.com/mercadopago/sdk-go v1.14.0 h1:3PYp9GPa+iysx2lcaKbpBEkXgEw4IIpbw3T6Jhl0IzI=
github.com/mercadopago/sdk-go v1.14.0/go.mod h1:hvQlOYb3MuYPGfjox7jeGBFbeL0nS+iPwUnxAv6yGWI=
github.com/pmezard/go-difflib v1.0.0 h1:4DBwDE0NGyQoBHbLQYPwSUPoCMWR5BEzIk/f1lZbAQM=
github.com/pmezard/go-difflib v1.0.0/go.mod h1:iKH77koFhYxTK1pcRnkKkqfTogsbg7gZNVY4sRDYZ/4=
github.com/stretchr/testify v1.7.0 h1:nwc3DEeHmmLAfoZucVR881uASk0Mfjw8xYJ99tb5CcY=
github.com/stretchr/testify v1.7.0/go.mod h1:6Fq8oRcR53rry900zMqJjRRixrwX3KX962/h/Wwjteg=
github.com/stripe/stripe-go/v86 v86.3.0 h1:BKtYc3NtRa4EGzKAmp4jvl5q7kk2rwMZ+llF18N5vHI=
github.com/stripe/stripe-go/v86 v86.3.0/go.mod h1:Co7QRXCKGNOPTugAdvjgRo+KcMtd9hxy+pZMN0yThsQ=
gopkg.in/yaml.v3 v3.0.1 h1:fxVm/GzAzEWqLHuvctI91KS9hhNmmWOoWu0XTYJS7CA=
gopkg.in/yaml.v3 v3.0.1/go.mod h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZviUmN+EgEM=
````

### FILE: `official_payment_webhooks/officialpayments/verifier.go`
```yaml
block_id: "GO-OFFICIAL-PAYMENT-WEBHOOKS:verifier:v1"
operation: CREATE
provenance: AUTHORED
source: "local adapter invoking only public APIs of exact official Stripe and Mercado Pago SDKs"
license: "LicenseRef-Workspace-Owner"
sha256: "64ba006387b43979648306f4b0acb0e4de3db7ea152f33f69ea0becd6ed74eac"
variables: []
secrets_allowed: false
```
````go
package officialpayments

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	mpwebhook "github.com/mercadopago/sdk-go/pkg/webhook"
	stripewebhook "github.com/stripe/stripe-go/v86/webhook"
)

const MaxPayloadBytes = 1 << 20

var (
	ErrInvalidInput      = errors.New("invalid payment webhook input")
	ErrPayloadTooLarge   = errors.New("payment webhook payload too large")
	ErrPayloadIDMismatch = errors.New("payment webhook payload id mismatch")
)

type Event struct {
	Provider  string
	ID        string
	Type      string
	RequestID string
	Payload   []byte
}

func VerifyStripe(payload []byte, signatureHeader, signingSecret string, tolerance time.Duration) (Event, error) {
	if len(payload) == 0 || strings.TrimSpace(signatureHeader) == "" || strings.TrimSpace(signingSecret) == "" || tolerance <= 0 {
		return Event{}, ErrInvalidInput
	}
	if len(payload) > MaxPayloadBytes {
		return Event{}, ErrPayloadTooLarge
	}
	event, err := stripewebhook.ConstructEventWithTolerance(payload, signatureHeader, signingSecret, tolerance)
	if err != nil {
		return Event{}, fmt.Errorf("stripe official webhook verification: %w", err)
	}
	return Event{Provider: "stripe", ID: event.ID, Type: string(event.Type), Payload: append([]byte(nil), payload...)}, nil
}

func VerifyMercadoPago(payload []byte, signatureHeader, requestID, dataID, signingSecret string, tolerance time.Duration) (Event, error) {
	if len(payload) == 0 || strings.TrimSpace(signatureHeader) == "" || strings.TrimSpace(signingSecret) == "" || strings.TrimSpace(dataID) == "" || tolerance <= 0 {
		return Event{}, ErrInvalidInput
	}
	if len(payload) > MaxPayloadBytes {
		return Event{}, ErrPayloadTooLarge
	}
	if err := mpwebhook.ValidateSignature(signatureHeader, requestID, dataID, signingSecret, mpwebhook.WithTolerance(tolerance)); err != nil {
		return Event{}, fmt.Errorf("mercado pago official webhook verification: %w", err)
	}
	var body struct {
		Type string `json:"type"`
		Data struct {
			ID json.RawMessage `json:"id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(payload, &body); err != nil {
		return Event{}, fmt.Errorf("mercado pago payload JSON: %w", err)
	}
	bodyID := strings.Trim(string(body.Data.ID), `"`)
	if bodyID != "" && !strings.EqualFold(bodyID, dataID) {
		return Event{}, ErrPayloadIDMismatch
	}
	return Event{Provider: "mercado_pago", ID: dataID, Type: body.Type, RequestID: requestID, Payload: append([]byte(nil), payload...)}, nil
}
````

### FILE: `official_payment_webhooks/officialpayments/client.go`
```yaml
block_id: "GO-OFFICIAL-PAYMENT-WEBHOOKS:client:v1"
operation: CREATE
provenance: AUTHORED
source: "local outbound boundary invoking exact official Stripe PaymentIntent and Mercado Pago Payment clients"
license: "LicenseRef-Workspace-Owner"
sha256: "9afae1052d458ab227a748b1a9f79c59e2dba1d91607bbd9cb17ab4132324c41"
variables: []
secrets_allowed: false
```
````go
package officialpayments

import (
	"context"
	"errors"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/mercadopago/sdk-go/pkg/config"
	"github.com/mercadopago/sdk-go/pkg/payment"
	"github.com/mercadopago/sdk-go/pkg/requester"
	"github.com/mercadopago/sdk-go/pkg/requestoptions"
	"github.com/stripe/stripe-go/v86"
)

var (
	ErrInvalidPaymentRequest = errors.New("invalid payment request")
	idempotencyPattern       = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:-]{0,127}$`)
	currencyPattern          = regexp.MustCompile(`^[a-z]{3}$`)
)

type PaymentResult struct {
	Provider     string
	ID           string
	Status       string
	ClientSecret string
}

type StripeIntentRequest struct {
	AmountMinor    int64
	Currency       string
	OrderID        string
	IdempotencyKey string
}

type StripeIntentClient struct{ client *stripe.Client }

func NewStripeIntentClient(secretKey string) (*StripeIntentClient, error) {
	return NewStripeIntentClientWithHTTPClient(secretKey, &http.Client{Timeout: 10 * time.Second})
}

// HTTP client injection is a trusted composition boundary. Production uses the
// SDK's fixed official URL; offline tests replace transport without provider IO.
func NewStripeIntentClientWithHTTPClient(secretKey string, client *http.Client) (*StripeIntentClient, error) {
	if client == nil {
		return nil, ErrInvalidPaymentRequest
	}
	bounded := *client
	bounded.Timeout = 10 * time.Second
	bounded.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	// The SDK's default error logger can include raw response samples. Configure
	// this backend only; the application emits its own bounded error codes.
	backend := stripe.GetBackendWithConfig(stripe.APIBackend, &stripe.BackendConfig{HTTPClient: &bounded, MaxNetworkRetries: stripe.Int64(0), LeveledLogger: &stripe.LeveledLogger{Level: stripe.LevelNull}})
	return newStripeIntentClient(secretKey, backend)
}

func newStripeIntentClient(secretKey string, backend stripe.Backend) (*StripeIntentClient, error) {
	if strings.TrimSpace(secretKey) == "" || backend == nil {
		return nil, ErrInvalidPaymentRequest
	}
	backends := &stripe.Backends{API: backend, Connect: backend, Uploads: backend, MeterEvents: backend}
	return &StripeIntentClient{client: stripe.NewClient(secretKey, stripe.WithBackends(backends))}, nil
}

func (c *StripeIntentClient) CreateIntent(ctx context.Context, request StripeIntentRequest) (PaymentResult, error) {
	if c == nil || c.client == nil || ctx == nil || request.AmountMinor <= 0 || !currencyPattern.MatchString(request.Currency) || strings.TrimSpace(request.OrderID) == "" || !idempotencyPattern.MatchString(request.IdempotencyKey) {
		return PaymentResult{}, ErrInvalidPaymentRequest
	}
	params := &stripe.PaymentIntentCreateParams{Amount: stripe.Int64(request.AmountMinor), Currency: stripe.String(request.Currency), Metadata: map[string]string{"order_id": request.OrderID}}
	params.SetIdempotencyKey(request.IdempotencyKey)
	intent, err := c.client.V1PaymentIntents.Create(ctx, params)
	if err != nil {
		return PaymentResult{}, err
	}
	return PaymentResult{Provider: "stripe", ID: intent.ID, Status: string(intent.Status), ClientSecret: intent.ClientSecret}, nil
}

type MercadoPagoPaymentRequest struct {
	AmountMinor       int64
	MinorUnitExponent int
	Currency          string
	OrderID           string
	PaymentMethodID   string
	PaymentToken      string
	PayerEmail        string
	NotificationURL   string
	IdempotencyKey    string
}

type MercadoPagoPaymentClient struct{ client payment.Client }

func NewMercadoPagoPaymentClient(accessToken string) (*MercadoPagoPaymentClient, error) {
	return NewMercadoPagoPaymentClientWithHTTPClient(accessToken, &http.Client{Timeout: 10 * time.Second})
}

func NewMercadoPagoPaymentClientWithHTTPClient(accessToken string, client *http.Client) (*MercadoPagoPaymentClient, error) {
	if client == nil {
		return nil, ErrInvalidPaymentRequest
	}
	bounded := *client
	bounded.Timeout = 10 * time.Second
	bounded.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	return newMercadoPagoPaymentClient(accessToken, &bounded)
}

func newMercadoPagoPaymentClient(accessToken string, transport requester.Requester) (*MercadoPagoPaymentClient, error) {
	if strings.TrimSpace(accessToken) == "" {
		return nil, ErrInvalidPaymentRequest
	}
	options := []config.Option{config.WithTimeout(10 * time.Second), config.WithMaxRetries(0)}
	if transport != nil {
		options = []config.Option{config.WithHTTPClient(transport), config.WithMaxRetries(0)}
	}
	cfg, err := config.New(accessToken, options...)
	if err != nil {
		return nil, err
	}
	return &MercadoPagoPaymentClient{client: payment.NewClient(cfg)}, nil
}

func (c *MercadoPagoPaymentClient) CreatePayment(ctx context.Context, request MercadoPagoPaymentRequest) (PaymentResult, error) {
	if c == nil || c.client == nil || ctx == nil || request.AmountMinor <= 0 || request.MinorUnitExponent < 0 || request.MinorUnitExponent > 3 || !currencyPattern.MatchString(strings.ToLower(request.Currency)) || strings.TrimSpace(request.OrderID) == "" || strings.TrimSpace(request.PaymentMethodID) == "" || strings.TrimSpace(request.PaymentToken) == "" || strings.TrimSpace(request.PayerEmail) == "" || !idempotencyPattern.MatchString(request.IdempotencyKey) {
		return PaymentResult{}, ErrInvalidPaymentRequest
	}
	amount := float64(request.AmountMinor) / math.Pow10(request.MinorUnitExponent)
	ctx = requestoptions.WithIdempotencyKey(ctx, request.IdempotencyKey)
	resource, err := c.client.Create(ctx, payment.Request{
		TransactionAmount: amount,
		PaymentMethodID:   request.PaymentMethodID,
		Payer:             &payment.PayerRequest{Email: request.PayerEmail},
		Token:             request.PaymentToken,
		Installments:      1,
		ExternalReference: request.OrderID,
		NotificationURL:   request.NotificationURL,
		Metadata:          map[string]any{"order_id": request.OrderID, "currency": strings.ToUpper(request.Currency)},
	})
	if err != nil {
		return PaymentResult{}, err
	}
	return PaymentResult{Provider: "mercado_pago", ID: fmtInt(resource.ID), Status: resource.Status}, nil
}

func fmtInt(value int) string {
	return strconv.Itoa(value)
}
````

### FILE: `official_payment_webhooks/officialpayments/client_test.go`
```yaml
block_id: "GO-OFFICIAL-PAYMENT-WEBHOOKS:client-test:v1"
operation: CREATE
provenance: AUTHORED
source: "local contract tests through official SDK HTTP clients and controlled local transports"
license: "LicenseRef-Workspace-Owner"
sha256: "6f0f8596f193a91e8063ef88ce2f5ebd54d5bbce10645145961fa38b1872c2d5"
variables: []
secrets_allowed: false
```
````go
package officialpayments

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/stripe/stripe-go/v86"
)

func TestStripeBackendDoesNotLogMalformedResponseBody(t *testing.T) {
	const marker = "ELITE_PAYMENT_LOG_FIXTURE"
	if os.Getenv(marker) == "child" {
		transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
			if r.URL.String() != "https://api.stripe.com/v1/account" {
				t.Fatal("unexpected provider request")
			}
			return &http.Response{StatusCode: 200, Header: http.Header{"Content-Type": []string{"application/json"}}, Body: io.NopCloser(strings.NewReader(`{"PRIVATE_RESPONSE_SENTINEL":`)), Request: r}, nil
		})
		client, err := NewStripeIntentClientWithHTTPClient("sk_test_fixture", &http.Client{Transport: transport})
		if err != nil {
			t.Fatal("fixture construction failed")
		}
		if _, err = client.ProbeAccount(context.Background()); err == nil {
			t.Fatal("malformed response accepted")
		}
		return
	}
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	command := exec.CommandContext(ctx, executable, "-test.run=^TestStripeBackendDoesNotLogMalformedResponseBody$")
	command.Env = append(os.Environ(), marker+"=child")
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatal("fixture subprocess failed")
	}
	if strings.Contains(string(output), "PRIVATE_RESPONSE_SENTINEL") || strings.Contains(string(output), "body sample") {
		t.Fatal("provider raw response reached process logs")
	}
}

func TestStripeCreateIntentOfficialClientContract(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/payment_intents" {
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer sk_test_local" {
			t.Fatalf("authorization header missing")
		}
		if r.Header.Get("Idempotency-Key") != "order-123-attempt-1" {
			t.Fatalf("idempotency header missing")
		}
		body, _ := io.ReadAll(r.Body)
		form, err := url.ParseQuery(string(body))
		if err != nil {
			t.Fatal(err)
		}
		if form.Get("amount") != "1099" || form.Get("currency") != "usd" || form.Get("metadata[order_id]") != "order-123" {
			t.Fatalf("unexpected form: %s", body)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"id":"pi_123","object":"payment_intent","status":"requires_payment_method","client_secret":"pi_secret"}`)
	}))
	defer server.Close()
	backend := stripe.GetBackendWithConfig(stripe.APIBackend, &stripe.BackendConfig{URL: stripe.String(server.URL), HTTPClient: server.Client(), MaxNetworkRetries: stripe.Int64(0)})
	client, err := newStripeIntentClient("sk_test_local", backend)
	if err != nil {
		t.Fatal(err)
	}
	result, err := client.CreateIntent(context.Background(), StripeIntentRequest{AmountMinor: 1099, Currency: "usd", OrderID: "order-123", IdempotencyKey: "order-123-attempt-1"})
	if err != nil {
		t.Fatal(err)
	}
	if result.ID != "pi_123" || result.Status != "requires_payment_method" || result.ClientSecret != "pi_secret" {
		t.Fatalf("unexpected result: %#v", result)
	}
}

type requesterFunc func(*http.Request) (*http.Response, error)

func (f requesterFunc) Do(r *http.Request) (*http.Response, error) { return f(r) }

func TestMercadoPagoCreatePaymentOfficialClientContract(t *testing.T) {
	transport := requesterFunc(func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodPost || r.URL.String() != "https://api.mercadopago.com/v1/payments" {
			t.Fatalf("unexpected request %s %s", r.Method, r.URL)
		}
		if r.Header.Get("Authorization") != "Bearer TEST-local" {
			t.Fatalf("authorization header missing")
		}
		if r.Header.Get("X-Idempotency-Key") != "order-456-attempt-1" {
			t.Fatalf("idempotency header missing: %#v", r.Header)
		}
		body, _ := io.ReadAll(r.Body)
		var payload map[string]any
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatal(err)
		}
		if payload["transaction_amount"] != 1250.5 || payload["external_reference"] != "order-456" || payload["payment_method_id"] != "visa" {
			t.Fatalf("unexpected body: %s", body)
		}
		return &http.Response{StatusCode: 201, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(`{"id":987,"status":"pending"}`)), Request: r}, nil
	})
	client, err := newMercadoPagoPaymentClient("TEST-local", transport)
	if err != nil {
		t.Fatal(err)
	}
	result, err := client.CreatePayment(context.Background(), MercadoPagoPaymentRequest{AmountMinor: 125050, MinorUnitExponent: 2, Currency: "ars", OrderID: "order-456", PaymentMethodID: "visa", PaymentToken: "card-token", PayerEmail: "payer@example.test", NotificationURL: "https://example.test/webhooks/mp", IdempotencyKey: "order-456-attempt-1"})
	if err != nil {
		t.Fatal(err)
	}
	if result.ID != "987" || result.Status != "pending" {
		t.Fatalf("unexpected result: %#v", result)
	}
}

func TestOutboundPaymentInputFailsClosed(t *testing.T) {
	stripeClient, _ := NewStripeIntentClient("sk_test")
	if _, err := stripeClient.CreateIntent(context.Background(), StripeIntentRequest{AmountMinor: 1, Currency: "USD", OrderID: "order", IdempotencyKey: "id"}); err == nil {
		t.Fatal("uppercase currency accepted")
	}
	mpClient, _ := NewMercadoPagoPaymentClient("TEST")
	if _, err := mpClient.CreatePayment(context.Background(), MercadoPagoPaymentRequest{AmountMinor: 1, MinorUnitExponent: 4}); err == nil {
		t.Fatal("invalid exponent accepted")
	}
}
````

### FILE: `official_payment_webhooks/officialpayments/verifier_test.go`
```yaml
block_id: "GO-OFFICIAL-PAYMENT-WEBHOOKS:test:v1"
operation: CREATE
provenance: AUTHORED
source: "local regression using official SDK verifier and Stripe official test-signature helper"
license: "LicenseRef-Workspace-Owner"
sha256: "e924989d2efbe1a3c2be3ecbbbc72a5bf4367e5e5b63e4a73902a6afea1b1939"
variables: []
secrets_allowed: false
```
````go
package officialpayments

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stripe/stripe-go/v86"
	stripewebhook "github.com/stripe/stripe-go/v86/webhook"
)

func TestVerifyStripeUsesOfficialSignatureAndAPIVersion(t *testing.T) {
	payload := []byte(fmt.Sprintf(`{"id":"evt_123","object":"event","api_version":%q,"type":"payment_intent.succeeded","data":{"object":{}}}`, stripe.APIVersion))
	signed := stripewebhook.GenerateTestSignedPayload(&stripewebhook.UnsignedPayload{Payload: payload, Secret: "whsec_test", Timestamp: time.Now()})
	event, err := VerifyStripe(payload, signed.Header, "whsec_test", 5*time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	if event.Provider != "stripe" || event.ID != "evt_123" || event.Type != "payment_intent.succeeded" {
		t.Fatalf("unexpected event: %#v", event)
	}

	tampered := append([]byte(nil), payload...)
	tampered[len(tampered)-2] ^= 1
	if _, err := VerifyStripe(tampered, signed.Header, "whsec_test", 5*time.Minute); err == nil {
		t.Fatal("tampered Stripe payload accepted")
	}
}

func mercadoPagoHeader(dataID, requestID, secret string, timestamp time.Time) string {
	ts := fmt.Sprintf("%d", timestamp.UnixMilli())
	manifest := "id:" + dataID + ";request-id:" + requestID + ";ts:" + ts + ";"
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(manifest))
	return "ts=" + ts + ",v1=" + hex.EncodeToString(mac.Sum(nil))
}

func TestVerifyMercadoPagoUsesOfficialSignatureAndCrossChecksID(t *testing.T) {
	const dataID = "123456"
	const requestID = "request-123"
	const secret = "mp_secret"
	payload := []byte(`{"type":"payment","data":{"id":"123456"}}`)
	header := mercadoPagoHeader(dataID, requestID, secret, time.Now())
	event, err := VerifyMercadoPago(payload, header, requestID, dataID, secret, 5*time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	if event.Provider != "mercado_pago" || event.ID != dataID || event.Type != "payment" || event.RequestID != requestID {
		t.Fatalf("unexpected event: %#v", event)
	}

	mismatch := []byte(`{"type":"payment","data":{"id":"different"}}`)
	if _, err := VerifyMercadoPago(mismatch, header, requestID, dataID, secret, 5*time.Minute); !errors.Is(err, ErrPayloadIDMismatch) {
		t.Fatalf("expected ID mismatch, got %v", err)
	}
	if _, err := VerifyMercadoPago(payload, strings.Replace(header, "v1=", "v1=00", 1), requestID, dataID, secret, 5*time.Minute); err == nil {
		t.Fatal("tampered Mercado Pago signature accepted")
	}
}

func TestFailClosedInputAndSize(t *testing.T) {
	if _, err := VerifyStripe([]byte(`{}`), "", "secret", time.Minute); !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("expected invalid input, got %v", err)
	}
	tooLarge := make([]byte, MaxPayloadBytes+1)
	if _, err := VerifyStripe(tooLarge, "header", "secret", time.Minute); !errors.Is(err, ErrPayloadTooLarge) {
		t.Fatalf("expected too large, got %v", err)
	}
	if _, err := VerifyMercadoPago([]byte(`{}`), "header", "request", "", "secret", time.Minute); !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("expected invalid input, got %v", err)
	}
}
````

### FILE: `official_payment_webhooks/README.md`
```yaml
block_id: "GO-OFFICIAL-PAYMENT-WEBHOOKS:readme:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "63d63590cf39c65efa80488bf8aaf2ffb71ad2ba46fa47d1ebaf8e12a846bbaf"
variables: []
secrets_allowed: false
```
````markdown
# Official payment SDK adapters

This standalone module invokes the exact official Stripe Go86.3.0 and MercadoPago Go1.14.0 dependency pins. Local files are AUTHORED boundary glue; their ownership is not attributed to either provider.

CreateIntent/CreatePayment submit idempotent requests. RetrieveIntent/RetrievePayment use official GET methods and return a bounded PaymentSnapshot for reconciliation. Stripe expands latest_charge and checks amount/capture/refund coherence, including Disputed. MercadoPago converts the SDK decimal representation to integer minor units without rounding; the selected contract supplies its currency exponent. The snapshot omits client secrets, payer/card details and raw response bodies. Its SHA-256 identifies the normalized local observation, not a provider signature.

StripeIntentClient.CreateCheckout/RetrieveCheckout provide the official hosted Checkout Session flow. MercadoPagoCheckoutClient.CreateCheckout/RetrieveCheckout provide official Checkout Pro Preference APIs for the admitted Argentina lane. CheckoutRequest binds OrderID, PaymentAttemptID, amount/currency, configured HTTPS return/callback URLs, idempotency key and finite expiry. Redirect URLs must use checkout.stripe.com or the selected official MercadoPago Argentina host. Stripe session and underlying PaymentIntent metadata include both internal IDs. MercadoPago keeps ExternalReference=OrderID plus metadata IDs. Its CheckoutResult.LiveMode describes the selected redirect lane; only subsequent PaymentSnapshot.LiveMode observes a payment's actual provider mode. A browser return or completed session never substitutes the payment GET observation. No payer/card input is accepted by these hosted-checkout methods.

VerifyStripePaymentNotification validates SDK signature/API version, account/mode and payment-intent identity. VerifyMercadoPagoPaymentNotification treats only signed data.id/request-id/timestamp as authenticated routing inputs; status, amount and user_id in the body cannot authorize money or delivery. Its normalized provider is mercadopago, matching Commerce; the older generic VerifyMercadoPago API retains its historical mercado_pago value for compatibility.

Production constructors use fixed SDK API URLs, 10-second HTTP budgets and zero automatic network retries. WithHTTPClient constructors are trusted composition injection points for offline transports; they do not select endpoints from user request data. The surrounding owner must bind tenant/connection/provider reference/order/currency/amount/mode/account to a durable intent and handle ambiguous effects before retrying. Never log SDK errors or PaymentResult.ClientSecret without redaction.

These adapters do not create a durable inbox, advance the enterprise payment ledger or authorize commercial release. The composition must supply those owners and gates. A Stripe PaymentIntent alone is not confirmation or collection. Connect-account routing, partial payments, refunds and disputes require their explicit selected workflow; no live effect runs in tests.

Offline verification: GOPROXY=off GOTOOLCHAIN=local go test ./...; go vet ./...; go build ./... with pinned toolchain and previously verified module cache. Fixtures exercise provider HTTP serialization, signed callbacks, exact reconciliation data, refunds, invalid IDs/mode/account, and unsigned MercadoPago body claims. Credentials are unnecessary for this reference verification.
````

### FILE: `official_payment_webhooks/officialpayments/retrieve.go`
```yaml
block_id: "GO-OFFICIAL-PAYMENT-WEBHOOKS:retrieve-go:v1"
operation: CREATE
provenance: AUTHORED
source: "bounded mapping glue over exact official Stripe/MercadoPago SDK public APIs; see PAYMENT_SDK_RECONCILIATION_V402.md"
license: "LicenseRef-Workspace-Owner"
sha256: "5148b8b7d9195cb300bbdba3b4b6f16f7b161958770b2a395d214a12cae22b26"
variables: []
secrets_allowed: false
```
````go
package officialpayments

// AUTHORED boundary mapping over the exact official SDK pins. This adapter
// observes provider resources; it never decides order release or edits a ledger.
import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"math"
	"math/big"
	"regexp"
	"strconv"
	"strings"

	"github.com/stripe/stripe-go/v86"
)

var ErrInvalidObservation = errors.New("payment observation does not match contract")
var stripeIntentID = regexp.MustCompile(`^pi_[A-Za-z0-9]{1,200}$`)
var decimalPaymentID = regexp.MustCompile(`^[1-9][0-9]{0,18}$`)
var stripeAccountID = regexp.MustCompile(`^acct_[A-Za-z0-9]{1,200}$`)

// ProbeAccount observes the account associated with this exact credential.
// It uses /v1/account, never an ID chosen from request input or /accounts/{id}.
func (c *StripeIntentClient) ProbeAccount(ctx context.Context) (string, error) {
	if c == nil || c.client == nil || ctx == nil {
		return "", ErrInvalidPaymentRequest
	}
	account, err := c.client.V1Accounts.Retrieve(ctx, &stripe.AccountRetrieveParams{})
	if err != nil {
		return "", err
	}
	if account == nil || !stripeAccountID.MatchString(account.ID) {
		return "", ErrInvalidObservation
	}
	return account.ID, nil
}

// PaymentSnapshot contains only reconciliation fields, never client secrets,
// payer details or card data. Status retains the provider's original vocabulary.
type PaymentSnapshot struct {
	ProviderCode      string `json:"provider_code"`
	ProviderReference string `json:"provider_reference"`
	OrderID           string `json:"order_id"`
	PaymentAttemptID  string `json:"payment_attempt_id,omitempty"`
	Currency          string `json:"currency"`
	AmountMinor       int64  `json:"amount_minor"`
	ReceivedMinor     int64  `json:"received_minor"`
	RefundedMinor     int64  `json:"refunded_minor"`
	CapturableMinor   int64  `json:"capturable_minor"`
	Status            string `json:"status"`
	LiveMode          bool   `json:"live_mode"`
	Disputed          bool   `json:"disputed"`
	CollectorID       string `json:"collector_id,omitempty"`
}

// SHA256 binds the normalized observation, not a provider-signed receipt.
func (s PaymentSnapshot) SHA256() string {
	data, _ := json.Marshal(s)
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func (c *StripeIntentClient) RetrieveIntent(ctx context.Context, id string) (PaymentSnapshot, error) {
	if c == nil || c.client == nil || ctx == nil || !stripeIntentID.MatchString(id) {
		return PaymentSnapshot{}, ErrInvalidPaymentRequest
	}
	params := &stripe.PaymentIntentRetrieveParams{}
	params.AddExpand("latest_charge")
	p, err := c.client.V1PaymentIntents.Retrieve(ctx, id, params)
	if err != nil {
		return PaymentSnapshot{}, err
	}
	if p == nil || p.ID != id || p.Amount <= 0 || p.AmountReceived < 0 || p.AmountReceived > p.Amount || p.AmountCapturable < 0 || p.AmountCapturable > p.Amount || !currencyPattern.MatchString(string(p.Currency)) || strings.TrimSpace(p.Metadata["order_id"]) == "" || len(p.Metadata["order_id"]) > 200 || p.Status == "" {
		return PaymentSnapshot{}, ErrInvalidObservation
	}
	refunded := int64(0)
	disputed := false
	if p.AmountReceived > 0 {
		charge := p.LatestCharge
		if charge == nil || charge.ID == "" || charge.Amount != p.Amount || charge.AmountCaptured != p.AmountReceived || charge.AmountRefunded < 0 || charge.AmountRefunded > charge.AmountCaptured || charge.Currency != p.Currency || !charge.Captured || charge.Livemode != p.Livemode {
			return PaymentSnapshot{}, ErrInvalidObservation
		}
		refunded = charge.AmountRefunded
		disputed = charge.Disputed
	}
	return PaymentSnapshot{ProviderCode: "stripe", ProviderReference: p.ID, OrderID: p.Metadata["order_id"], PaymentAttemptID: p.Metadata["payment_attempt_id"], Currency: strings.ToUpper(string(p.Currency)), AmountMinor: p.Amount, ReceivedMinor: p.AmountReceived, RefundedMinor: refunded, CapturableMinor: p.AmountCapturable, Status: string(p.Status), LiveMode: p.Livemode, Disputed: disputed}, nil
}

// decimalMinor converts the SDK's decimal JSON float representation without
// rounding a fractional minor unit. The bounded exact integer is checked again
// by the calling owner against the durable order currency/amount.
func decimalMinor(value float64, exponent int) (int64, error) {
	if math.IsNaN(value) || math.IsInf(value, 0) || value < 0 || exponent < 0 || exponent > 3 {
		return 0, ErrInvalidObservation
	}
	amount, ok := new(big.Rat).SetString(strconv.FormatFloat(value, 'f', -1, 64))
	if !ok {
		return 0, ErrInvalidObservation
	}
	amount.Mul(amount, new(big.Rat).SetInt64(int64(math.Pow10(exponent))))
	if !amount.IsInt() || !amount.Num().IsInt64() {
		return 0, ErrInvalidObservation
	}
	minor := amount.Num().Int64()
	if minor > 1<<53-1 {
		return 0, ErrInvalidObservation
	}
	return minor, nil
}

// RetrievePayment requires the currency exponent from the selected provider
// contract; no currency table or merchant financial policy is invented here.
func (c *MercadoPagoPaymentClient) RetrievePayment(ctx context.Context, id string, exponent int) (PaymentSnapshot, error) {
	if c == nil || c.client == nil || ctx == nil || !decimalPaymentID.MatchString(id) || exponent < 0 || exponent > 3 {
		return PaymentSnapshot{}, ErrInvalidPaymentRequest
	}
	numeric, err := strconv.Atoi(id)
	if err != nil || numeric <= 0 {
		return PaymentSnapshot{}, ErrInvalidPaymentRequest
	}
	p, err := c.client.Get(ctx, numeric)
	if err != nil {
		return PaymentSnapshot{}, err
	}
	if p == nil || p.ID != numeric || !currencyPattern.MatchString(strings.ToLower(p.CurrencyID)) || strings.TrimSpace(p.ExternalReference) == "" || len(p.ExternalReference) > 200 || p.Status == "" || p.CollectorID <= 0 {
		return PaymentSnapshot{}, ErrInvalidObservation
	}
	amount, err := decimalMinor(p.TransactionAmount, exponent)
	if err != nil || amount <= 0 {
		return PaymentSnapshot{}, ErrInvalidObservation
	}
	refunded, err := decimalMinor(p.TransactionAmountRefunded, exponent)
	if err != nil || refunded > amount {
		return PaymentSnapshot{}, ErrInvalidObservation
	}
	received := int64(0)
	if p.Captured {
		received = amount
	}
	attempt, _ := p.Metadata["payment_attempt_id"].(string)
	return PaymentSnapshot{ProviderCode: "mercadopago", ProviderReference: id, OrderID: p.ExternalReference, PaymentAttemptID: attempt, Currency: strings.ToUpper(p.CurrencyID), AmountMinor: amount, ReceivedMinor: received, RefundedMinor: refunded, Status: p.Status, LiveMode: p.LiveMode, Disputed: p.Status == "charged_back", CollectorID: strconv.FormatInt(p.CollectorID, 10)}, nil
}
````

### FILE: `official_payment_webhooks/officialpayments/notification.go`
```yaml
block_id: "GO-OFFICIAL-PAYMENT-WEBHOOKS:notification-go:v1"
operation: CREATE
provenance: AUTHORED
source: "bounded mapping glue over exact official Stripe/MercadoPago SDK public APIs; see PAYMENT_SDK_RECONCILIATION_V402.md"
license: "LicenseRef-Workspace-Owner"
sha256: "11ad6b9441b06c474ff2e5e7862c6fd118347b24d5993db8311b2ec631752da9"
variables: []
secrets_allowed: false
```
````go
package officialpayments

// AUTHORED normalization only. SDK verifiers remain the sole provider signature
// implementation. Notifications schedule GET reconciliation; they never attest
// payment success merely because a callback body declares it.
import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strings"
	"time"
)

type PaymentNotification struct {
	ProviderCode      string `json:"provider_code"`
	ProviderEventID   string `json:"provider_event_id"`
	ProviderReference string `json:"provider_reference"`
	ResourceType      string `json:"resource_type"`
	Type              string `json:"type"`
	BodySHA256        string `json:"body_sha256"`
}

func bodyDigest(payload []byte) string {
	sum := sha256.Sum256(payload)
	return hex.EncodeToString(sum[:])
}

func VerifyStripePaymentNotification(payload []byte, signature, secret string, tolerance time.Duration, expectedLive bool, expectedAccount string) (PaymentNotification, error) {
	event, err := VerifyStripe(payload, signature, secret, tolerance)
	if err != nil {
		return PaymentNotification{}, err
	}
	var body struct {
		Account string `json:"account"`
		Live    *bool  `json:"livemode"`
		Data    struct {
			Object struct {
				ID            string          `json:"id"`
				Object        string          `json:"object"`
				PaymentIntent json.RawMessage `json:"payment_intent"`
			} `json:"object"`
		} `json:"data"`
	}
	if json.Unmarshal(payload, &body) != nil || event.ID == "" || len(event.ID) > 200 || body.Account != expectedAccount || body.Live == nil || *body.Live != expectedLive {
		return PaymentNotification{}, ErrInvalidObservation
	}
	resource := ""
	if strings.HasPrefix(event.Type, "payment_intent.") && body.Data.Object.Object == "payment_intent" && stripeIntentID.MatchString(body.Data.Object.ID) {
		resource = "payment_intent"
	}
	if strings.HasPrefix(event.Type, "checkout.session.") && body.Data.Object.Object == "checkout.session" && checkoutSessionID.MatchString(body.Data.Object.ID) {
		resource = "checkout_session"
	}
	if (strings.HasPrefix(event.Type, "charge.") || strings.HasPrefix(event.Type, "refund.")) && (body.Data.Object.Object == "charge" || body.Data.Object.Object == "dispute" || body.Data.Object.Object == "refund") {
		var intent string
		if json.Unmarshal(body.Data.Object.PaymentIntent, &intent) != nil {
			var object struct {
				ID string `json:"id"`
			}
			if json.Unmarshal(body.Data.Object.PaymentIntent, &object) == nil {
				intent = object.ID
			}
		}
		if stripeIntentID.MatchString(intent) {
			body.Data.Object.ID = intent
			resource = "payment_intent"
		}
	}
	if resource == "" {
		return PaymentNotification{}, ErrInvalidObservation
	}
	return PaymentNotification{ProviderCode: "stripe", ProviderEventID: event.ID, ProviderReference: body.Data.Object.ID, ResourceType: resource, Type: event.Type, BodySHA256: bodyDigest(payload)}, nil
}

func VerifyMercadoPagoPaymentNotification(payload []byte, signature, requestID, dataID, secret string, tolerance time.Duration) (PaymentNotification, error) {
	// Mercado Pago signs data.id/request-id/timestamp, not the full JSON body.
	// Only that authenticated resource ID selects GET; unsigned body status,
	// merchant and amounts cannot authorize domain effects.
	if strings.TrimSpace(requestID) == "" || len(requestID) > 128 || strings.ContainsAny(requestID, "\r\n") || !decimalPaymentID.MatchString(dataID) {
		return PaymentNotification{}, ErrInvalidInput
	}
	_, err := VerifyMercadoPago(payload, signature, requestID, dataID, secret, tolerance)
	if err != nil {
		return PaymentNotification{}, err
	}
	return PaymentNotification{ProviderCode: "mercadopago", ProviderEventID: "mp:" + requestID + ":" + dataID, ProviderReference: dataID, ResourceType: "payment", Type: "payment.reconciliation-requested", BodySHA256: bodyDigest(payload)}, nil
}
````

### FILE: `official_payment_webhooks/officialpayments/retrieve_test.go`
```yaml
block_id: "GO-OFFICIAL-PAYMENT-WEBHOOKS:retrieve-test-go:v1"
operation: CREATE
provenance: AUTHORED
source: "bounded mapping glue over exact official Stripe/MercadoPago SDK public APIs; see PAYMENT_SDK_RECONCILIATION_V402.md"
license: "LicenseRef-Workspace-Owner"
sha256: "0a6b926e852f8a5480f953dbf8c79df0622c70ab2b20f6405ed5f341c3da6913"
variables: []
secrets_allowed: false
```
````go
package officialpayments

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stripe/stripe-go/v86"
	stripewebhook "github.com/stripe/stripe-go/v86/webhook"
)

func TestStripeRetrieveUsesOfficialGETAndExcludesSecrets(t *testing.T) {
	var calls atomic.Int32
	response := `{"id":"pi_123","object":"payment_intent","amount":1099,"amount_received":1099,"amount_capturable":0,"currency":"usd","status":"succeeded","livemode":false,"metadata":{"order_id":"order-123"},"client_secret":"never-export-secret","latest_charge":{"id":"ch_123","amount":1099,"amount_captured":1099,"amount_refunded":100,"currency":"usd","captured":true,"livemode":false}}`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		if r.Method != "GET" || r.URL.Path != "/v1/payment_intents/pi_123" || r.URL.Query().Get("expand[0]") != "latest_charge" || r.Header.Get("Authorization") != "Bearer sk_test_local" {
			t.Errorf("unexpected SDK GET: %s %s", r.Method, r.URL)
		}
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, response)
	}))
	defer server.Close()
	backend := stripe.GetBackendWithConfig(stripe.APIBackend, &stripe.BackendConfig{URL: stripe.String(server.URL), HTTPClient: server.Client(), MaxNetworkRetries: stripe.Int64(0)})
	client, _ := newStripeIntentClient("sk_test_local", backend)
	snapshot, err := client.RetrieveIntent(context.Background(), "pi_123")
	if err != nil {
		t.Fatal(err)
	}
	if snapshot.ProviderCode != "stripe" || snapshot.OrderID != "order-123" || snapshot.AmountMinor != 1099 || snapshot.ReceivedMinor != 1099 || snapshot.RefundedMinor != 100 || snapshot.Currency != "USD" {
		t.Fatalf("snapshot=%+v", snapshot)
	}
	data, _ := json.Marshal(snapshot)
	if strings.Contains(string(data), "secret") || len(snapshot.SHA256()) != 64 || calls.Load() != 1 {
		t.Fatal("observation leak or extra request")
	}
	for _, id := range []string{"", "pi_123/other", "../pi_123", "pi_123?secret=x"} {
		if _, err = client.RetrieveIntent(context.Background(), id); err == nil {
			t.Fatal("invalid id accepted")
		}
	}
	if calls.Load() != 1 {
		t.Fatal("invalid identifier sent to provider")
	}
	response = strings.Replace(response, `"amount_captured":1099`, `"amount_captured":1098`, 1)
	if _, err = client.RetrieveIntent(context.Background(), "pi_123"); err == nil {
		t.Fatal("inconsistent charge accepted")
	}
}

func TestMercadoPagoRetrieveOfficialGETAndExactMinorUnits(t *testing.T) {
	calls := 0
	transport := requesterFunc(func(r *http.Request) (*http.Response, error) {
		calls++
		if r.Method != "GET" || r.URL.String() != "https://api.mercadopago.com/v1/payments/987" || r.Header.Get("Authorization") != "Bearer TEST-local" {
			t.Errorf("unexpected SDK GET: %s %s", r.Method, r.URL)
		}
		body := `{"id":987,"transaction_amount":1250.50,"transaction_amount_refunded":0.01,"currency_id":"ARS","external_reference":"order-456","status":"approved","captured":true,"collector_id":123,"live_mode":false,"payer":{"email":"private@example.test"}}`
		return &http.Response{StatusCode: 200, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body)), Request: r}, nil
	})
	client, _ := newMercadoPagoPaymentClient("TEST-local", transport)
	snapshot, err := client.RetrievePayment(context.Background(), "987", 2)
	if err != nil {
		t.Fatal(err)
	}
	if snapshot.ProviderCode != "mercadopago" || snapshot.AmountMinor != 125050 || snapshot.ReceivedMinor != 125050 || snapshot.RefundedMinor != 1 || snapshot.OrderID != "order-456" || snapshot.CollectorID != "123" {
		t.Fatalf("snapshot=%+v", snapshot)
	}
	data, _ := json.Marshal(snapshot)
	if strings.Contains(string(data), "private") {
		t.Fatal("PII leaked")
	}
	for _, id := range []string{"", "0", "-1", "1/2", "001", "9223372036854775808"} {
		if _, err = client.RetrievePayment(context.Background(), id, 2); err == nil {
			t.Fatal("invalid numeric id")
		}
	}
	if calls != 1 {
		t.Fatal("invalid identifier sent to provider")
	}
}

func TestDecimalConversionRejectsRoundingAndUnsupportedRange(t *testing.T) {
	for _, value := range []float64{-1, math.Inf(1), math.NaN(), 1.001, 1e20} {
		if _, err := decimalMinor(value, 2); err == nil {
			t.Fatalf("accepted %v", value)
		}
	}
	for _, row := range []struct {
		value float64
		exp   int
		want  int64
	}{{1.23, 2, 123}, {0, 2, 0}, {1.001, 3, 1001}, {123, 0, 123}} {
		got, err := decimalMinor(row.value, row.exp)
		if err != nil || got != row.want {
			t.Fatalf("got %d %v", got, err)
		}
	}
}

func TestStripePaymentNotificationBindsModeAccountAndResource(t *testing.T) {
	body := []byte(fmt.Sprintf(`{"id":"evt_123","object":"event","api_version":%q,"type":"payment_intent.succeeded","livemode":false,"data":{"object":{"id":"pi_123","object":"payment_intent"}}}`, stripe.APIVersion))
	signed := stripewebhook.GenerateTestSignedPayload(&stripewebhook.UnsignedPayload{Payload: body, Secret: "whsec_test", Timestamp: time.Now()})
	notification, err := VerifyStripePaymentNotification(body, signed.Header, "whsec_test", 5*time.Minute, false, "")
	if err != nil || notification.ProviderReference != "pi_123" || notification.ProviderEventID != "evt_123" || notification.ProviderCode != "stripe" {
		t.Fatalf("notification=%+v err=%v", notification, err)
	}
	if _, err = VerifyStripePaymentNotification(body, signed.Header, "whsec_test", 5*time.Minute, true, ""); err == nil {
		t.Fatal("live mismatch accepted")
	}
	if _, err = VerifyStripePaymentNotification(body, signed.Header, "whsec_test", 5*time.Minute, false, "acct_wrong"); err == nil {
		t.Fatal("account mismatch accepted")
	}
	if _, err = VerifyStripePaymentNotification(body, signed.Header, "wrong", 5*time.Minute, false, ""); err == nil {
		t.Fatal("bad signature accepted")
	}
}

func TestMercadoPagoNotificationUsesSignedResourceNotBodyFinancialClaims(t *testing.T) {
	body := []byte(`{"type":"payment","data":{"id":"987"},"status":"approved","transaction_amount":999999,"user_id":666}`)
	header := mercadoPagoHeader("987", "request-1", "secret", time.Now())
	notification, err := VerifyMercadoPagoPaymentNotification(body, header, "request-1", "987", "secret", 5*time.Minute)
	if err != nil || notification.ProviderCode != "mercadopago" || notification.ProviderReference != "987" || notification.Type != "payment.reconciliation-requested" {
		t.Fatalf("%+v %v", notification, err)
	}
	data, _ := json.Marshal(notification)
	if strings.Contains(string(data), "approved") || strings.Contains(string(data), "999999") {
		t.Fatal("unsigned financial claims propagated")
	}
	if _, err = VerifyMercadoPagoPaymentNotification(body, header, "other-request", "987", "secret", 5*time.Minute); err == nil {
		t.Fatal("request substitution accepted")
	}
	if _, err = VerifyMercadoPagoPaymentNotification(body, header, "request-1", "988", "secret", 5*time.Minute); err == nil {
		t.Fatal("resource substitution accepted")
	}
}
````

### FILE: `official_payment_webhooks/officialpayments/checkout.go`
```yaml
block_id: "GO-OFFICIAL-PAYMENT-WEBHOOKS:checkout-go:v1"
operation: CREATE
provenance: AUTHORED
source: "bounded mapping glue over exact official Stripe/MercadoPago SDK public APIs; see PAYMENT_SDK_RECONCILIATION_V402.md"
license: "LicenseRef-Workspace-Owner"
sha256: "1533b728ca09a06016aa8e0c7b2828cb9ea630b9b3f1bab6834f2d70999b5338"
variables: []
secrets_allowed: false
```
````go
package officialpayments

// AUTHORED request/response mapping for official hosted payment UIs. Payment
// credentials/card entry remain at the provider; a return URL is never receipt.
import (
	"context"
	"math"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/mercadopago/sdk-go/pkg/config"
	"github.com/mercadopago/sdk-go/pkg/preference"
	"github.com/mercadopago/sdk-go/pkg/requestoptions"
	"github.com/stripe/stripe-go/v86"
)

var checkoutSessionID = regexp.MustCompile(`^cs_(test|live)_[A-Za-z0-9]{1,200}$`)
var preferenceID = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9-]{0,200}$`)

type CheckoutRequest struct {
	OrderID, PaymentAttemptID, DisplayName string
	AmountMinor                            int64
	Currency                               string
	MinorUnitExponent                      int
	SuccessURL, CancelURL, NotificationURL string
	IdempotencyKey                         string
	ExpiresAt                              time.Time
}
type CheckoutResult struct {
	ProviderCode     string    `json:"provider_code"`
	SessionID        string    `json:"session_id"`
	URL              string    `json:"url,omitempty"`
	PaymentReference string    `json:"payment_reference,omitempty"`
	OrderID          string    `json:"order_id"`
	PaymentAttemptID string    `json:"payment_attempt_id"`
	Currency         string    `json:"currency"`
	AmountMinor      int64     `json:"amount_minor"`
	Status           string    `json:"status"`
	PaymentStatus    string    `json:"payment_status,omitempty"`
	LiveMode         bool      `json:"live_mode"`
	CollectorID      string    `json:"collector_id,omitempty"`
	ExpiresAt        time.Time `json:"expires_at"`
}

func secureURL(raw string) bool {
	u, err := url.Parse(raw)
	return err == nil && len(raw) <= 8192 && u.Scheme == "https" && u.Hostname() != "" && u.User == nil && !strings.ContainsAny(raw, "\r\n") && (u.Port() == "" || u.Port() == "443")
}
func hostedURL(raw, host string) bool {
	if !secureURL(raw) {
		return false
	}
	u, _ := url.Parse(raw)
	return strings.EqualFold(u.Hostname(), host)
}
func validCheckout(r CheckoutRequest) bool {
	return r.AmountMinor > 0 && r.AmountMinor <= 1<<53-1 && currencyPattern.MatchString(strings.ToLower(r.Currency)) && r.OrderID != "" && len(r.OrderID) <= 200 && r.PaymentAttemptID != "" && len(r.PaymentAttemptID) <= 200 && strings.TrimSpace(r.DisplayName) != "" && len(r.DisplayName) <= 200 && idempotencyPattern.MatchString(r.IdempotencyKey) && secureURL(r.SuccessURL) && secureURL(r.CancelURL) && !r.ExpiresAt.IsZero() && r.ExpiresAt.After(time.Now().Add(30*time.Minute)) && r.ExpiresAt.Before(time.Now().Add(24*time.Hour))
}
func (c *StripeIntentClient) CreateCheckout(ctx context.Context, r CheckoutRequest) (CheckoutResult, error) {
	if c == nil || c.client == nil || ctx == nil || !validCheckout(r) {
		return CheckoutResult{}, ErrInvalidPaymentRequest
	}
	metadata := map[string]string{"order_id": r.OrderID, "payment_attempt_id": r.PaymentAttemptID}
	params := &stripe.CheckoutSessionCreateParams{Mode: stripe.String("payment"), SuccessURL: stripe.String(r.SuccessURL), CancelURL: stripe.String(r.CancelURL), ClientReferenceID: stripe.String(r.PaymentAttemptID), ExpiresAt: stripe.Int64(r.ExpiresAt.Unix()), Metadata: metadata, PaymentIntentData: &stripe.CheckoutSessionCreatePaymentIntentDataParams{Metadata: metadata}, LineItems: []*stripe.CheckoutSessionCreateLineItemParams{{Quantity: stripe.Int64(1), PriceData: &stripe.CheckoutSessionCreateLineItemPriceDataParams{Currency: stripe.String(strings.ToLower(r.Currency)), UnitAmount: stripe.Int64(r.AmountMinor), ProductData: &stripe.CheckoutSessionCreateLineItemPriceDataProductDataParams{Name: stripe.String(r.DisplayName)}}}}}
	params.SetIdempotencyKey(r.IdempotencyKey)
	session, err := c.client.V1CheckoutSessions.Create(ctx, params)
	if err != nil {
		return CheckoutResult{}, err
	}
	result, err := stripeCheckoutResult(session)
	if err != nil || result.OrderID != r.OrderID || result.PaymentAttemptID != r.PaymentAttemptID || result.Currency != strings.ToUpper(r.Currency) || result.AmountMinor != r.AmountMinor || result.URL == "" {
		return CheckoutResult{}, ErrInvalidObservation
	}
	return result, nil
}
func (c *StripeIntentClient) RetrieveCheckout(ctx context.Context, id string) (CheckoutResult, error) {
	if c == nil || c.client == nil || ctx == nil || !checkoutSessionID.MatchString(id) {
		return CheckoutResult{}, ErrInvalidPaymentRequest
	}
	session, err := c.client.V1CheckoutSessions.Retrieve(ctx, id, &stripe.CheckoutSessionRetrieveParams{})
	if err != nil {
		return CheckoutResult{}, err
	}
	result, err := stripeCheckoutResult(session)
	if err == nil && result.SessionID != id {
		return CheckoutResult{}, ErrInvalidObservation
	}
	return result, err
}
func stripeCheckoutResult(s *stripe.CheckoutSession) (CheckoutResult, error) {
	if s == nil || !checkoutSessionID.MatchString(s.ID) || s.Mode != "payment" || s.AmountTotal <= 0 || !currencyPattern.MatchString(string(s.Currency)) || s.Metadata["order_id"] == "" || s.Metadata["payment_attempt_id"] == "" || s.ClientReferenceID != s.Metadata["payment_attempt_id"] || s.Status == "" || s.PaymentStatus == "" || (s.URL != "" && !hostedURL(s.URL, "checkout.stripe.com")) {
		return CheckoutResult{}, ErrInvalidObservation
	}
	ref := ""
	if s.PaymentIntent != nil {
		ref = s.PaymentIntent.ID
		if !stripeIntentID.MatchString(ref) {
			return CheckoutResult{}, ErrInvalidObservation
		}
	}
	return CheckoutResult{ProviderCode: "stripe", SessionID: s.ID, URL: s.URL, PaymentReference: ref, OrderID: s.Metadata["order_id"], PaymentAttemptID: s.Metadata["payment_attempt_id"], Currency: strings.ToUpper(string(s.Currency)), AmountMinor: s.AmountTotal, Status: string(s.Status), PaymentStatus: string(s.PaymentStatus), LiveMode: s.Livemode, ExpiresAt: time.Unix(s.ExpiresAt, 0).UTC()}, nil
}

// Argentina is the admitted Checkout Pro country lane for this reference.
// Other countries/custom checkout domains require their explicit provider lock.
type MercadoPagoCheckoutClient struct {
	client  preference.Client
	sandbox bool
}

func NewMercadoPagoCheckoutClient(token string, sandbox bool) (*MercadoPagoCheckoutClient, error) {
	return NewMercadoPagoCheckoutClientWithHTTPClient(token, sandbox, &http.Client{Timeout: 10 * time.Second})
}
func NewMercadoPagoCheckoutClientWithHTTPClient(token string, sandbox bool, httpClient *http.Client) (*MercadoPagoCheckoutClient, error) {
	if strings.TrimSpace(token) == "" || httpClient == nil {
		return nil, ErrInvalidPaymentRequest
	}
	bounded := *httpClient
	bounded.Timeout = 10 * time.Second
	bounded.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	cfg, err := config.New(token, config.WithHTTPClient(&bounded), config.WithMaxRetries(0))
	if err != nil {
		return nil, err
	}
	return &MercadoPagoCheckoutClient{client: preference.NewClient(cfg), sandbox: sandbox}, nil
}
func (c *MercadoPagoCheckoutClient) CreateCheckout(ctx context.Context, r CheckoutRequest) (CheckoutResult, error) {
	if c == nil || c.client == nil || ctx == nil || !validCheckout(r) || r.MinorUnitExponent < 0 || r.MinorUnitExponent > 3 || !secureURL(r.NotificationURL) {
		return CheckoutResult{}, ErrInvalidPaymentRequest
	}
	amount := float64(r.AmountMinor) / math.Pow10(r.MinorUnitExponent)
	roundTrip, err := decimalMinor(amount, r.MinorUnitExponent)
	if err != nil || roundTrip != r.AmountMinor {
		return CheckoutResult{}, ErrInvalidPaymentRequest
	}
	ctx = requestoptions.WithIdempotencyKey(ctx, r.IdempotencyKey)
	expires := r.ExpiresAt
	p, err := c.client.Create(ctx, preference.Request{Items: []preference.ItemRequest{{ID: r.OrderID, Title: r.DisplayName, Quantity: 1, CurrencyID: strings.ToUpper(r.Currency), UnitPrice: amount}}, ExternalReference: r.OrderID, Metadata: map[string]any{"order_id": r.OrderID, "payment_attempt_id": r.PaymentAttemptID}, NotificationURL: r.NotificationURL, BackURLs: &preference.BackURLsRequest{Success: r.SuccessURL, Pending: r.SuccessURL, Failure: r.CancelURL}, Expires: true, ExpirationDateTo: &expires})
	if err != nil {
		return CheckoutResult{}, err
	}
	result, err := c.preferenceResult(p, r.MinorUnitExponent)
	if err != nil || result.OrderID != r.OrderID || result.PaymentAttemptID != r.PaymentAttemptID || result.Currency != strings.ToUpper(r.Currency) || result.AmountMinor != r.AmountMinor {
		return CheckoutResult{}, ErrInvalidObservation
	}
	return result, nil
}
func (c *MercadoPagoCheckoutClient) RetrieveCheckout(ctx context.Context, id string, exponent int) (CheckoutResult, error) {
	if c == nil || c.client == nil || ctx == nil || !preferenceID.MatchString(id) {
		return CheckoutResult{}, ErrInvalidPaymentRequest
	}
	p, err := c.client.Get(ctx, id)
	if err != nil {
		return CheckoutResult{}, err
	}
	result, err := c.preferenceResult(p, exponent)
	if err == nil && result.SessionID != id {
		return CheckoutResult{}, ErrInvalidObservation
	}
	return result, err
}
func (c *MercadoPagoCheckoutClient) preferenceResult(p *preference.Response, exponent int) (CheckoutResult, error) {
	if p == nil || !preferenceID.MatchString(p.ID) || p.ExternalReference == "" || p.CollectorID <= 0 || len(p.Items) != 1 || p.Items[0].Quantity != 1 {
		return CheckoutResult{}, ErrInvalidObservation
	}
	order, ok := p.Metadata["order_id"].(string)
	attempt, ok2 := p.Metadata["payment_attempt_id"].(string)
	if !ok || !ok2 || order != p.ExternalReference || attempt == "" || !currencyPattern.MatchString(strings.ToLower(p.Items[0].CurrencyID)) {
		return CheckoutResult{}, ErrInvalidObservation
	}
	amount, err := decimalMinor(p.Items[0].UnitPrice, exponent)
	if err != nil || amount <= 0 {
		return CheckoutResult{}, ErrInvalidObservation
	}
	location, host := p.InitPoint, "www.mercadopago.com.ar"
	if c.sandbox {
		location, host = p.SandboxInitPoint, "sandbox.mercadopago.com.ar"
	}
	if !hostedURL(location, host) {
		return CheckoutResult{}, ErrInvalidObservation
	}
	return CheckoutResult{ProviderCode: "mercadopago", SessionID: p.ID, URL: location, OrderID: order, PaymentAttemptID: attempt, Currency: strings.ToUpper(p.Items[0].CurrencyID), AmountMinor: amount, Status: "preference-created", LiveMode: !c.sandbox, CollectorID: fmtInt(int(p.CollectorID)), ExpiresAt: p.ExpirationDateTo}, nil
}
````

### FILE: `official_payment_webhooks/officialpayments/checkout_test.go`
```yaml
block_id: "GO-OFFICIAL-PAYMENT-WEBHOOKS:checkout-test-go:v1"
operation: CREATE
provenance: AUTHORED
source: "bounded mapping glue over exact official Stripe/MercadoPago SDK public APIs; see PAYMENT_SDK_RECONCILIATION_V402.md"
license: "LicenseRef-Workspace-Owner"
sha256: "754bef6d11f6d66574cf5555b95c111ed35cc70d69cbe44242797dfb7f21efa3"
variables: []
secrets_allowed: false
```
````go
package officialpayments

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/stripe/stripe-go/v86"
	stripewebhook "github.com/stripe/stripe-go/v86/webhook"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func checkoutFixtureRequest() CheckoutRequest {
	return CheckoutRequest{OrderID: "order-1", PaymentAttemptID: "attempt-1", AmountMinor: 125050, Currency: "ARS", MinorUnitExponent: 2, DisplayName: "Order order-1", SuccessURL: "https://shop.example.test/payment-return", CancelURL: "https://shop.example.test/payment-cancel", NotificationURL: "https://api.example.test/payment-webhook", IdempotencyKey: "attempt-1-checkout", ExpiresAt: time.Now().Add(time.Hour)}
}
func TestStripeHostedCheckoutSDKCreatesAndRetrievesBoundSession(t *testing.T) {
	calls := []string{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, r.Method)
		if r.Method == "POST" {
			if err := r.ParseForm(); err != nil {
				t.Error(err)
			}
			for key, want := range map[string]string{"mode": "payment", "client_reference_id": "attempt-1", "metadata[order_id]": "order-1", "metadata[payment_attempt_id]": "attempt-1", "payment_intent_data[metadata][payment_attempt_id]": "attempt-1", "line_items[0][price_data][unit_amount]": "125050", "line_items[0][price_data][currency]": "ars", "line_items[0][quantity]": "1"} {
				if r.Form.Get(key) != want {
					t.Errorf("%s=%s", key, r.Form.Get(key))
				}
			}
			if r.Header.Get("Idempotency-Key") != "attempt-1-checkout" {
				t.Error("missing idempotency")
			}
		}
		if r.Method == "GET" && r.URL.Path != "/v1/checkout/sessions/cs_test_abc" {
			t.Errorf("unexpected GET %s", r.URL)
		}
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"id":"cs_test_abc","object":"checkout.session","mode":"payment","amount_total":125050,"currency":"ars","client_reference_id":"attempt-1","metadata":{"order_id":"order-1","payment_attempt_id":"attempt-1"},"status":"open","payment_status":"unpaid","livemode":false,"url":"https://checkout.stripe.com/c/pay/cs_test_abc#provider-fragment","payment_intent":"pi_abc","expires_at":2000000000}`)
	}))
	defer server.Close()
	backend := stripe.GetBackendWithConfig(stripe.APIBackend, &stripe.BackendConfig{URL: stripe.String(server.URL), HTTPClient: server.Client(), MaxNetworkRetries: stripe.Int64(0)})
	c, _ := newStripeIntentClient("sk_test_fixture", backend)
	created, err := c.CreateCheckout(context.Background(), checkoutFixtureRequest())
	if err != nil {
		t.Fatal(err)
	}
	read, err := c.RetrieveCheckout(context.Background(), created.SessionID)
	if err != nil {
		t.Fatal(err)
	}
	if read != created || read.PaymentReference != "pi_abc" || read.PaymentAttemptID != "attempt-1" || len(calls) != 2 {
		t.Fatalf("%+v calls%v", read, calls)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
func TestMercadoPagoHostedCheckoutSDKCreateAndGetNoPayerData(t *testing.T) {
	calls := 0
	client := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		calls++
		if r.URL.Host != "api.mercadopago.com" || !strings.HasPrefix(r.URL.Path, "/checkout/preferences") {
			t.Errorf("wrong endpoint %s", r.URL)
		}
		if r.Method == "POST" {
			var body map[string]any
			if json.NewDecoder(r.Body).Decode(&body) != nil {
				t.Error("bad body")
			}
			metadata := body["metadata"].(map[string]any)
			if body["external_reference"] != "order-1" || metadata["payment_attempt_id"] != "attempt-1" || body["payer"] != nil {
				t.Errorf("wrong contract %v", body)
			}
			if r.Header.Get("X-Idempotency-Key") != "attempt-1-checkout" {
				t.Error("missing key")
			}
		}
		body := `{"id":"123-abc","external_reference":"order-1","metadata":{"order_id":"order-1","payment_attempt_id":"attempt-1"},"collector_id":123,"items":[{"currency_id":"ARS","quantity":1,"unit_price":1250.50}],"init_point":"https://www.mercadopago.com.ar/checkout/v1/redirect?pref_id=123-abc","sandbox_init_point":"https://sandbox.mercadopago.com.ar/checkout/v1/redirect?pref_id=123-abc","expiration_date_to":"2030-01-01T00:00:00Z"}`
		return &http.Response{StatusCode: 200, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body)), Request: r}, nil
	})}
	c, _ := NewMercadoPagoCheckoutClientWithHTTPClient("TEST-fixture", true, client)
	created, err := c.CreateCheckout(context.Background(), checkoutFixtureRequest())
	if err != nil {
		t.Fatal(err)
	}
	read, err := c.RetrieveCheckout(context.Background(), created.SessionID, 2)
	if err != nil {
		t.Fatal(err)
	}
	if read != created || read.Currency != "ARS" || read.CollectorID != "123" || read.AmountMinor != 125050 || read.LiveMode || calls != 2 {
		t.Fatalf("%+v calls%d", read, calls)
	}
}
func TestHostedCheckoutRejectsUntrustedRedirectAndIncompleteRequest(t *testing.T) {
	for _, raw := range []string{"http://checkout.stripe.com/x", "https://checkout.stripe.com.evil.test/x", "https://user@checkout.stripe.com/x", "https://checkout.stripe.com:444/x", "javascript:alert(1)"} {
		if hostedURL(raw, "checkout.stripe.com") {
			t.Fatalf("unsafe URL%s", raw)
		}
	}
	req := checkoutFixtureRequest()
	req.PaymentAttemptID = ""
	if validCheckout(req) {
		t.Fatal("missing attempt accepted")
	}
	req = checkoutFixtureRequest()
	req.ExpiresAt = time.Time{}
	if validCheckout(req) {
		t.Fatal("unbounded checkout accepted")
	}
}
func TestStripeRefundDisputeAndCheckoutCallbacksRouteOnlyToOfficialGET(t *testing.T) {
	for _, row := range []struct{ kind, object, id, intent, resource, want string }{{"checkout.session.completed", "checkout.session", "cs_test_abc", "", "checkout_session", "cs_test_abc"}, {"charge.refunded", "charge", "ch_abc", "pi_abc", "payment_intent", "pi_abc"}, {"charge.dispute.created", "dispute", "dp_abc", "pi_abc", "payment_intent", "pi_abc"}} {
		body := []byte(fmt.Sprintf(`{"id":"evt_abc","object":"event","api_version":%q,"type":%q,"livemode":false,"data":{"object":{"id":%q,"object":%q,"payment_intent":%q}}}`, stripe.APIVersion, row.kind, row.id, row.object, row.intent))
		signed := stripewebhook.GenerateTestSignedPayload(&stripewebhook.UnsignedPayload{Payload: body, Secret: "secret", Timestamp: time.Now()})
		n, err := VerifyStripePaymentNotification(body, signed.Header, "secret", time.Minute, false, "")
		if err != nil || n.ResourceType != row.resource || n.ProviderReference != row.want {
			t.Fatalf("%s %+v %v", row.kind, n, err)
		}
	}
}
func TestStripeDisputedSnapshotCannotMasqueradeAsUncontestedReceipt(t *testing.T) {
	c, _ := NewStripeIntentClientWithHTTPClient("sk_test_fixture", &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		body := `{"id":"pi_abc","amount":100,"amount_received":100,"currency":"usd","status":"succeeded","metadata":{"order_id":"order-1","payment_attempt_id":"attempt-1"},"latest_charge":{"id":"ch_abc","amount":100,"amount_captured":100,"amount_refunded":0,"currency":"usd","captured":true,"disputed":true}}`
		return &http.Response{StatusCode: 200, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body)), Request: r}, nil
	})})
	snapshot, err := c.RetrieveIntent(context.Background(), "pi_abc")
	if err != nil || !snapshot.Disputed || snapshot.PaymentAttemptID != "attempt-1" {
		t.Fatalf("%+v %v", snapshot, err)
	}
}

func TestStripeAccountProbeUsesCredentialBoundOfficialEndpoint(t *testing.T) {
	var seen bool
	c, _ := NewStripeIntentClientWithHTTPClient("sk_test_fixture", &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		seen = true
		if r.Method != "GET" || r.URL.String() != "https://api.stripe.com/v1/account" || r.Header.Get("Authorization") != "Bearer sk_test_fixture" {
			t.Errorf("wrong account probe %s", r.URL)
		}
		return &http.Response{StatusCode: 200, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(`{"id":"acct_fixture","object":"account","email":"not-exported@example.test"}`)), Request: r}, nil
	})})
	id, err := c.ProbeAccount(context.Background())
	if err != nil || id != "acct_fixture" || !seen {
		t.Fatalf("probe %s %v", id, err)
	}
}
````

### FILE: `official_payment_webhooks/officialpayments/preference_recovery.go`
```yaml
block_id: "GO-OFFICIAL-PAYMENT-WEBHOOKS:preference-recovery:v1"
operation: CREATE
provenance: AUTHORED
source: "bounded identity recovery through Search/Get public APIs in Mercado Pago SDK 1.14.0 at f910ee53fbb6819e435eaf3d0f800cb1fe74ae09; local validation glue, no vendor authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "fd3dad0004905048125f0017a3e0efd0733aa782628c624303f0f1a12941d6a8"
variables: []
secrets_allowed: false
```
````go
package officialpayments

// AUTHORED bounded recovery mapping over preference.Search/Get in the pinned
// Mercado Pago SDK. Search summaries locate identity; complete GET proves it.
import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/mercadopago/sdk-go/pkg/preference"
)

type PreferenceRecoveryRequest struct {
	OrderID, PaymentAttemptID, Currency, CollectorID string
	AmountMinor                                      int64
	MinorUnitExponent                                int
	LiveMode                                         bool
	ExpiresAt                                        time.Time
}

// Exactly one total result is required; ambiguous/truncated or older-than-search
// history cannot authorize a guessed preference. No mutation API is called.
func (c *MercadoPagoCheckoutClient) RecoverCheckout(ctx context.Context, r PreferenceRecoveryRequest) (CheckoutResult, error) {
	if c == nil || c.client == nil || ctx == nil || r.OrderID == "" || len(r.OrderID) > 200 || r.PaymentAttemptID == "" || len(r.PaymentAttemptID) > 200 || r.AmountMinor <= 0 || !currencyPattern.MatchString(strings.ToLower(r.Currency)) || r.MinorUnitExponent < 0 || r.MinorUnitExponent > 3 || r.ExpiresAt.IsZero() || r.LiveMode == c.sandbox {
		return CheckoutResult{}, ErrInvalidPaymentRequest
	}
	collector, err := strconv.ParseInt(r.CollectorID, 10, 64)
	if err != nil || collector <= 0 || strconv.FormatInt(collector, 10) != r.CollectorID {
		return CheckoutResult{}, ErrInvalidPaymentRequest
	}
	page, err := c.client.Search(ctx, preference.SearchRequest{Limit: 2, Offset: 0, Filters: map[string]string{"external_reference": r.OrderID, "site_id": "MLA"}})
	if err != nil {
		return CheckoutResult{}, err
	}
	if page == nil || page.Total != 1 || len(page.Elements) != 1 || page.NextOffset < 0 || page.NextOffset > 1 {
		return CheckoutResult{}, ErrInvalidObservation
	}
	item := page.Elements[0]
	if !preferenceID.MatchString(item.ID) || item.ExternalReference != r.OrderID || item.CollectorID != collector || item.SiteID != "MLA" || item.LiveMode != r.LiveMode || !item.Expires || !item.ExpirationDateTo.Equal(r.ExpiresAt) {
		return CheckoutResult{}, ErrInvalidObservation
	}
	complete, err := c.client.Get(ctx, item.ID)
	if err != nil {
		return CheckoutResult{}, err
	}
	result, err := c.preferenceResult(complete, r.MinorUnitExponent)
	if err != nil {
		return CheckoutResult{}, err
	}
	if result.SessionID != item.ID || result.OrderID != r.OrderID || result.PaymentAttemptID != r.PaymentAttemptID || result.Currency != r.Currency || result.AmountMinor != r.AmountMinor || result.CollectorID != r.CollectorID || result.LiveMode != r.LiveMode || !result.ExpiresAt.Equal(r.ExpiresAt) || !complete.Expires || complete.SiteID != "MLA" {
		return CheckoutResult{}, ErrInvalidObservation
	}
	return result, nil
}
````

## 6. Configuration surface

| ParÃ¡metro | Tipo/default | ValidaciÃ³n | Secreto | Mutabilidad/efecto |
|---|---|---|---|---|
| payload | bytes / ninguno | 1..1 MiB, body crudo | puede contener PII; no log | por request; se copia tras verificar |
| signature header | string | no vacÃ­o; parser SDK | no | por request |
| request/data ID MP | string | data ID obligatorio; body/query iguales | no | por request |
| signing secret | string | no vacÃ­o, secret manager | sÃ­ | rotaciÃ³n exige ventana/endpoint coordinados |
| tolerance | duration | `>0` | no | config versionada; afecta replay window |

Combinaciones invÃ¡lidas fallan antes de devolver evento. El edge debe preservar body/headers, limitar request total y no aceptar tenant/provider desde el body.

## 7. Dependency bill

| Package/image/tool | Pin exacto | Uso | Licencia | Runtime/build | Fuente oficial |
|---|---|---|---|---|---|
| Go | 1.26.7 windows/amd64 en evidencia | build/test | BSD-3-Clause | build | go.dev |
| `github.com/stripe/stripe-go/v86` | 86.3.0 | signature/API version/event | MIT | runtime | Stripe |
| `github.com/mercadopago/sdk-go` | 1.14.0 | signature/tolerance/reasons | MIT | runtime | Mercado Pago |
| transitive test modules | exactos en `go.sum` | tests del graph | segÃºn mÃ³dulo | build | Go checksum DB/upstreams |

Las dos fuentes oficiales exactas tambiÃ©n estÃ¡n fijadas por commit/archive/SHA/licencia en `OFFICIAL-UPSTREAM-ACQUISITION-CORE`.

## 8. Apply order

1. Materializar y ejecutar `go mod download`/`go test` con toolchain fijado.
2. Seleccionar proveedor/cuenta/producto y registrar tÃ©rminos/secret references.
3. Conectar despuÃ©s de body-preserving edge y antes de inbox durable.
4. Mapear el evento autenticado a dedup/job/reconciliation del core comÃºn.
5. Ejecutar fixtures locales, webhook sandbox oficial, replay/conflict y reconciliaciÃ³n.

En workspace existente, integrar como mÃ³dulo/import o adaptar paths explÃ­citamente; no copiar sobre `go.mod`. Rollback desconecta routing y conserva inbox/raw evidence. No borrar eventos ni marcar pagos reconciliados durante rollback.

## 9. Verification

```powershell
pwsh -NoProfile -File .\materialize_markdown_pack.ps1 -PackFile .\implementation_packs\GO_OFFICIAL_PAYMENT_WEBHOOK_ADAPTERS.md -Destination <empty>
cd <empty>\official_payment_webhooks
go mod download
go test ./...
```

Resultado esperado: build y seis suites PASS. Los tests ejercen firma Stripe oficial, API version, body alterado, firma Mercado Pago vÃ¡lida/alterada, ID query/body divergente, inputs incompletos, payload >1 MiB y requests outbound de ambos SDKs contra transports locales, incluyendo authorization, idempotency, amount/currency y order metadata. `GOPROXY=off go test ./...` debe pasar despuÃ©s de poblar el cache, demostrando reconstrucciÃ³n offline del graph ya adquirido. Sandbox y reconciliaciÃ³n son gates de proyecto, no falseados por fixtures.

## 10. Reconstruction evidence

- workspace limpio: `%LOCALAPPDATA%\Temp\elite-official-payment-adapters-verify-20260826`;
- Go oficial: 1.26.7 windows/amd64;
- SDKs: Stripe 86.3.0 y Mercado Pago 1.14.0, sources adquiridos por lock oficial;
- archivos: siete, SHA-256/manifest verificados por materializador;
- tests: `go test ./...` PASS; signature/tamper/size/ID e outbound request/idempotency incluidos;
- fallo aprendido: `stripe.EventType` requiriÃ³ conversiÃ³n explÃ­cita (`LIB-FAIL-035`);
- fecha/revisor: 2026-08-26 / Codex; expediente `GO_OFFICIAL_PAYMENT_WEBHOOK_ADAPTERS_2026-08-26_V1.md`.

## V402 current delta â€” 0.3.0

Historical V1 receipts remain above. Current12files use the same SDK pins. Added official GET/retrieve, credential-bound account probe, hosted checkout, minimal snapshot, refund/dispute observation, callback normalization and eleven targeted tests; production HTTP constructors have finite10s timeout and zero automatic retries. This is AUTHORED integration glue around DEPENDENCY_PIN, not vendor source or enterprise settlement logic. Source/cache binding, local gates and current SCA are recorded in `reconstruction_evidence/PAYMENT_SDK_RECONCILIATION_V402.md`. The main franchise wiring/inbox/consumer/handover gates remain separate and cannot be inferred from this pack PASS.

## V402 logger correction â€” 0.3.1

The default Stripe backend logger emitted a raw response body sample in the real malformed-response E2E fixture. This patch configures the official per-backend LevelNull logger without changing global logging. The host emits bounded operation/error codes. A subprocess regression verifies malformed provider JSON fails while stdout/stderr contain no response sample. No dependency pin, monetary algorithm, SDK operation, materialized file count or license changed.

V402 preference identity recovery: Search is bounded to two summaries and requires
exactly one total result for the external order reference/site. A complete Get
then verifies attempt metadata, collector, mode, currency, amount and expiry.
The provider's documented 90-day search window remains explicit; absent,
ambiguous, truncated or inconsistent identity fails closed, without another POST.
The connected callback processor retains separate financial-observation and
checkout-identity receipts, including recovery after SaveCheckout/fence crashes.
Official reference: https://www.mercadopago.com.ar/developers/es/reference/online-payments/checkout-pro-preferences/search-preferences/get

The new file is AUTHORED orchestration/validation glue around the unchanged fixed
SDK APIs. All 12 previous payloads are byte-identical. Existing dependency/SCA and
license evidence applies to the unchanged module graph, not to untested new code.
New behavior evidence: V402 stage mp-recovery-agent/result.json and source-lock.json;
three top-level/nine PostgreSQL cases, one creation POST and one money outbox per
successful fixture, with missing/ambiguous/binding negatives and vet PASS.
Exact pack reconstruction is required before REBUILD_VERIFIED promotion.

Exact reconstruction V402: materialize_markdown_pack.ps1 emitted all13 declared files into sdk-recovery-roundtrip, and all13 SHA-256 values match the candidate manifest. No previous payload was changed. The PostgreSQL recovery receipt remains separate behavior evidence; no live production or complete-library claim.

Canonical V402 integration: selected by the current profile with exact dependencies and caller overlays. Metadata promotion records byte reconstruction, not closure of every admission/release gate. Payload provenance is unchanged.
