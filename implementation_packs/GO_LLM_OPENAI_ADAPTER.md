# Go OpenAI-Compatible LLM Adapter

## 1. Metadata

```yaml
pack_id: "GO-LLM-OPENAI-ADAPTER"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa el adapter que implementa aifoundation.LLMProvider sobre un endpoint Chat Completions OpenAI-compatible; la API key se inyecta del entorno en runtime, nunca en código, y los tests usan un servidor HTTP falso."
stacks: ["Go 1.26.7"]
compatible_with: ["GO-ML-AI-FOUNDATION 0.1.x", "GO-CONVERSATIONAL-AGENT 0.1.x", "GO-LLM-ECONOMY-CORE 0.1.x"]
incompatible_with: ["modelo mutable (latest/auto)", "API key en código", "endpoint sin URL absoluta http(s)"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://platform.openai.com/docs/api-reference/chat"]
verified_at: "2026-09-02"
```

Este adapter implementa el contrato `aifoundation.LLMProvider`; no afirma que un proveedor concreto funcione ni su costo. La credencial se inyecta del entorno (`Config.APIKey`), nunca se hardcodea. El modelo es un nombre exacto (nunca `latest`/`auto`). Los tests usan `httptest` y no requieren red ni cuenta.

## 2. Applicability

Use este pack para conectar el agente a un LLM OpenAI-compatible (OpenAI o cualquier endpoint compatible) con salida estructurada (`response_format: json_object`).

Rechace este pack para: un proveedor con protocolo distinto a Chat Completions; o una API key embebida en código.

## 3. Architecture contract

- **Ownership**: `internal/llmopenai` gobierna la configuración y el mapeo request/response. `internal/aifoundation` gobierna el contrato y el pin de modelo.
- **Invariantes**: (1) URL absoluta http(s) y modelo exacto (no mutable). (2) API key sólo del entorno. (3) status no-2xx, contenido vacío o payload inesperado → error (fail-closed). (4) `response_format: json_object` sólo si hay schema.
- **Data flow**: `Generate` → mapear `CompletionRequest` → POST `/chat/completions` → parsear → `CompletionResponse`.
- **Failure modes**: 4xx/5xx, contenido vacío, JSON inválido → error.
- **Seguridad/privacidad**: key por header Bearer; sin logs de key; límite de lectura del body de error (4 KiB).
- **Performance budget**: timeout configurado; sin reintentos automáticos (los gobierna la capa de economía).
- **Operación/migración/rollback**: sin migración; reemplazar `Config` es composición.

## 4. Exact file manifest

```text
CREATE internal/llmopenai/config.go
CREATE internal/llmopenai/adapter.go
CREATE internal/llmopenai/adapter_test.go
```

## 5. Materialization blocks

### FILE: `internal/llmopenai/config.go`
```yaml
block_id: "GO-LLM-OPENAI-ADAPTER:internal/llmopenai/config.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "e361c4e7cb32ecd28da64843dec667b254acc700babf120eb9e8d46067a0d28e"
variables: []
secrets_allowed: false
```
````go
// Package llmopenai implements aifoundation.LLMProvider over an
// OpenAI-compatible Chat Completions endpoint. The API key is injected from
// the environment at runtime, never hardcoded; tests use a fake HTTP server.
package llmopenai

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"
)

// Config carries the runtime endpoint settings. The API key is a secret
// injected from the environment, never stored in source.
type Config struct {
	BaseURL string // e.g. https://api.openai.com/v1
	Model   string // exact model name, never "latest"/"auto"
	APIKey  string // secret, from env
	Timeout time.Duration
}

// Validate enforces an exact, safe configuration.
func (c Config) Validate() error {
	u, err := url.Parse(c.BaseURL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return errors.New("llmopenai: base url must be absolute http(s)")
	}
	m := strings.ToLower(strings.TrimSpace(c.Model))
	if m == "" || m == "latest" || m == "auto" {
		return fmt.Errorf("llmopenai: model must be an exact name, got %q", c.Model)
	}
	if strings.TrimSpace(c.APIKey) == "" {
		return errors.New("llmopenai: api key required (inject from environment)")
	}
	if c.Timeout <= 0 {
		return errors.New("llmopenai: timeout required")
	}
	return nil
}
````

### FILE: `internal/llmopenai/adapter.go`
```yaml
block_id: "GO-LLM-OPENAI-ADAPTER:internal/llmopenai/adapter.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "32f42b87994db1a225decc753989e543a5d38a588b69be1cc4fecca36c77be72"
variables: []
secrets_allowed: false
```
````go
package llmopenai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"elite.local/enterprise/internal/aifoundation"
)

// Client is an OpenAI-compatible Chat Completions adapter.
type Client struct {
	cfg    Config
	client *http.Client
}

// New returns a validated client.
func New(cfg Config) (*Client, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return &Client{cfg: cfg, client: &http.Client{Timeout: cfg.Timeout}}, nil
}

type chatRequest struct {
	Model          string              `json:"model"`
	Messages       []chatMessage       `json:"messages"`
	MaxTokens      int                 `json:"max_tokens,omitempty"`
	Temperature    *float64            `json:"temperature,omitempty"`
	Seed           *int64              `json:"seed,omitempty"`
	ResponseFormat *chatResponseFormat `json:"response_format,omitempty"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResponseFormat struct {
	Type string `json:"type"`
}

type chatResponse struct {
	Choices []struct {
		Message      chatMessage `json:"message"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
}

// Generate maps aifoundation.CompletionRequest to a Chat Completions call and
// returns the raw content. It fails closed on non-2xx, empty content or an
// unexpected payload.
func (c *Client) Generate(ctx context.Context, req aifoundation.CompletionRequest) (aifoundation.CompletionResponse, error) {
	messages := make([]chatMessage, 0, len(req.Messages))
	for _, m := range req.Messages {
		messages = append(messages, chatMessage{Role: m.Role, Content: m.Content})
	}
	body := chatRequest{
		Model:       c.cfg.Model,
		Messages:    messages,
		MaxTokens:   req.MaxOutputTokens,
		Temperature: req.Temperature,
		Seed:        req.Seed,
	}
	if len(req.Schema) > 0 {
		body.ResponseFormat = &chatResponseFormat{Type: "json_object"}
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return aifoundation.CompletionResponse{}, err
	}
	endpoint := strings.TrimRight(c.cfg.BaseURL, "/") + "/chat/completions"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(payload))
	if err != nil {
		return aifoundation.CompletionResponse{}, err
	}
	httpReq.Header.Set("content-type", "application/json")
	httpReq.Header.Set("authorization", "Bearer "+c.cfg.APIKey)

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return aifoundation.CompletionResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		data, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return aifoundation.CompletionResponse{}, fmt.Errorf("llmopenai: status %d: %s", resp.StatusCode, strings.TrimSpace(string(data)))
	}

	var parsed chatResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return aifoundation.CompletionResponse{}, err
	}
	if len(parsed.Choices) == 0 {
		return aifoundation.CompletionResponse{}, errors.New("llmopenai: empty choices")
	}
	content := parsed.Choices[0].Message.Content
	if strings.TrimSpace(content) == "" {
		return aifoundation.CompletionResponse{}, errors.New("llmopenai: empty content")
	}
	return aifoundation.CompletionResponse{
		Content:      json.RawMessage(content),
		FinishReason: parsed.Choices[0].FinishReason,
	}, nil
}
````

### FILE: `internal/llmopenai/adapter_test.go`
```yaml
block_id: "GO-LLM-OPENAI-ADAPTER:internal/llmopenai/adapter_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "b01956c6264ce79144b99128bbddf1799f0e6bfa5c5f283f3102a1ac397f65ab"
variables: []
secrets_allowed: false
```
````go
package llmopenai

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"elite.local/enterprise/internal/aifoundation"
)

func testConfig(serverURL string) Config {
	return Config{BaseURL: serverURL, Model: "gpt-4.1-mini", APIKey: "test-key", Timeout: 5 * time.Second}
}

func TestConfigValidate(t *testing.T) {
	if err := testConfig("https://api.openai.com/v1").Validate(); err != nil {
		t.Fatalf("valid config rejected: %v", err)
	}
	bad := testConfig("https://api.openai.com/v1")
	bad.Model = "latest"
	if err := bad.Validate(); err == nil {
		t.Fatal("mutable model accepted")
	}
	noKey := testConfig("https://api.openai.com/v1")
	noKey.APIKey = ""
	if err := noKey.Validate(); err == nil {
		t.Fatal("missing key accepted")
	}
	badURL := testConfig("not-a-url")
	if err := badURL.Validate(); err == nil {
		t.Fatal("bad url accepted")
	}
}

func TestGenerateRoundTrip(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/chat/completions" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("authorization") != "Bearer test-key" {
			t.Errorf("missing bearer")
		}
		var body chatRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body.Model != "gpt-4.1-mini" || body.ResponseFormat == nil || body.ResponseFormat.Type != "json_object" {
			t.Errorf("unexpected request: %+v", body)
		}
		w.Header().Set("content-type", "application/json")
		_, _ = w.Write([]byte(`{"choices":[{"message":{"role":"assistant","content":"{\"answer\":\"hola\"}"},"finish_reason":"stop"}]}`))
	}))
	defer srv.Close()

	c, err := New(testConfig(srv.URL))
	if err != nil {
		t.Fatal(err)
	}
	resp, err := c.Generate(context.Background(), aifoundation.CompletionRequest{
		Model:    aifoundation.ModelRef{Provider: "openai", Model: "gpt-4.1-mini", Version: "v", Digest: strings.Repeat("a", 64)},
		Messages: []aifoundation.Message{{Role: "user", Content: "hola"}},
		Schema:   json.RawMessage(`{"type":"object"}`),
	})
	if err != nil {
		t.Fatal(err)
	}
	if string(resp.Content) != `{"answer":"hola"}` || resp.FinishReason != "stop" {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestGenerateFailsClosedOnErrorStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":{"message":"bad key"}}`))
	}))
	defer srv.Close()

	c, _ := New(testConfig(srv.URL))
	if _, err := c.Generate(context.Background(), aifoundation.CompletionRequest{
		Model:    aifoundation.ModelRef{Provider: "openai", Model: "m", Version: "v", Digest: strings.Repeat("a", 64)},
		Messages: []aifoundation.Message{{Role: "user", Content: "x"}},
	}); err == nil {
		t.Fatal("expected error on 401")
	}
}

func TestGenerateRejectsEmptyContent(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		_, _ = w.Write([]byte(`{"choices":[{"message":{"role":"assistant","content":""},"finish_reason":"stop"}]}`))
	}))
	defer srv.Close()

	c, _ := New(testConfig(srv.URL))
	if _, err := c.Generate(context.Background(), aifoundation.CompletionRequest{
		Model:    aifoundation.ModelRef{Provider: "openai", Model: "m", Version: "v", Digest: strings.Repeat("a", 64)},
		Messages: []aifoundation.Message{{Role: "user", Content: "x"}},
	}); err == nil {
		t.Fatal("expected error on empty content")
	}
}
````


## 6. Configuration surface

| Variable | Tipo | Default seguro | Secreto | Efecto |
|---|---|---|---|---|
| `BaseURL` | URL http(s) | ninguno | no | endpoint del provider |
| `Model` | string exacto | ninguno | no | modelo pinnado |
| `APIKey` | string | ninguno | **sí** | Bearer token (del entorno) |
| `Timeout` | duration | ninguno | no | deadline de la llamada |

## 7. Dependency bill

| Package/image/tool | Pin exacto | Uso | Licencia | Runtime/build | Fuente oficial |
|---|---|---|---|---|---|
| Go standard library | go 1.26.7 (toolchain) | net/http, encoding/json | BSD-3-Clause | runtime | https://go.dev |
| internal/aifoundation | 0.1.0 | contrato LLMProvider | LicenseRef-Workspace-Owner | runtime | este repositorio |

## 8. Apply order

1. Componer `GO-ML-AI-FOUNDATION` (mismo módulo).
2. Colocar los tres archivos bajo `internal/llmopenai/`.
3. Inyectar `APIKey` del entorno del target (nunca en código).
4. Verificar con `go test ./... -count=1` y `go vet ./...`.
5. Rollback: eliminar `internal/llmopenai/`; no deja estado.

## 9. Verification

- `go test ./internal/llmopenai/ -count=1`: 4/4 PASS (config válida/negativos, round-trip con servidor falso verificando path/Bearer/modelo/response_format, fail-closed en 401, rechazo de contenido vacío).
- `go test ./... -count=1` (10 paquetes): PASS.
- `go vet ./...`: exit 0.

## 10. Reconstruction evidence

Véase `reconstruction_evidence/GO_LLM_OPENAI_ADAPTER_2026-09-02_V186.md`.
