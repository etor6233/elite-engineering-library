# Go Messaging Channel Adapters

## 1. Metadata

```yaml
pack_id: "GO-MESSAGING-CHANNEL-ADAPTERS"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa los adapters de canales de mensajería: Meta Graph (WhatsApp/Instagram/Facebook Messenger) con verificación HMAC de webhook y envío de texto, y email B2B vía API HTTP (SendGrid-compatible); credenciales inyectadas en runtime."
stacks: ["Go 1.26.7"]
compatible_with: ["GO-CHANNELS-CORE 0.1.x", "GO-CONVERSATIONAL-AGENT 0.1.x"]
incompatible_with: ["API key/token en código", "webhook sin verificación de firma"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://developers.facebook.com/docs/messenger-platform", "https://developers.facebook.com/docs/whatsapp/cloud-api", "https://www.twilio.com/docs/sendgrid/api-reference"]
verified_at: "2026-09-02"
```

Estos adapters implementan el envío y la verificación de webhook para los canales de mensajería. Las credenciales (access token, API key, app secret) se inyectan del entorno en runtime, nunca en código. Los tests usan `httptest` y vectores HMAC conocidos; no requieren red ni cuenta.

## 2. Applicability

Use este pack para conectar el agente a WhatsApp/Instagram/Facebook Messenger (Meta Graph) y a email B2B (relay HTTP), con verificación de firma anti-fraude en los webhooks.

Rechace este pack para: webhooks sin verificación de firma; o credenciales embebidas en código.

## 3. Architecture contract

- **Ownership**: `internal/msgchannels` gobierna envío, verificación de firma y parseo. `internal/channels` gobierna el contrato `Channel` y el ruteo.
- **Invariantes**: (1) firma HMAC-SHA256 del webhook verificada con compare constante. (2) status no-2xx → error (fail-closed). (3) credenciales sólo del entorno. (4) endpoint/URL absolutos http(s).
- **Data flow**: webhook entrante → `VerifyWebhookSignature` → `ParseWebhook` → dispatcher; respuesta → `SendText`/`SendEmail`.
- **Failure modes**: firma inválida → rechazo; status error → error; payload vacío → error.
- **Seguridad/privacidad**: token por Bearer; lectura de error acotada (4 KiB); sin logs de credenciales.
- **Performance budget**: timeout por llamada; sin reintentos automáticos.
- **Operación/migración/rollback**: sin migración; reemplazar credenciales es runtime.

## 4. Exact file manifest

```text
CREATE internal/msgchannels/metagraph.go
CREATE internal/msgchannels/email.go
CREATE internal/msgchannels/metagraph_test.go
CREATE internal/msgchannels/email_test.go
```

## 5. Materialization blocks

### FILE: `internal/msgchannels/metagraph.go`
```yaml
block_id: "GO-MESSAGING-CHANNEL-ADAPTERS:internal/msgchannels/metagraph.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "5dd17bdc4afb336921c86e2e2acf3b29075435bd4bb93f618fb46d579574115d"
variables: []
secrets_allowed: false
```
````go
// Package msgchannels provides channel adapters for the messaging surfaces:
// Meta Graph (WhatsApp/Instagram/Facebook Messenger) with HMAC webhook
// verification and text send, and SMTP email for B2B sales. Credentials are
// injected at runtime; tests use a fake HTTP server and known HMAC vectors.
package msgchannels

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// VerifyWebhookSignature checks the Meta X-Hub-Signature-256 header against the
// raw body using HMAC-SHA256 with the app secret (constant-time compare).
func VerifyWebhookSignature(appSecret string, rawBody []byte, header string) bool {
	if appSecret == "" || len(rawBody) == 0 || header == "" {
		return false
	}
	parts := strings.SplitN(strings.TrimSpace(header), "=", 2)
	if len(parts) != 2 || parts[0] != "sha256" {
		return false
	}
	got, err := hex.DecodeString(parts[1])
	if err != nil {
		return false
	}
	mac := hmac.New(sha256.New, []byte(appSecret))
	mac.Write(rawBody)
	return hmac.Equal(got, mac.Sum(nil))
}

// GraphConfig carries the Meta Graph API settings. AccessToken is a secret from
// the environment, never hardcoded.
type GraphConfig struct {
	GraphURL    string // https://graph.facebook.com
	Version     string // e.g. v22.0
	SenderID    string // phone_number_id (WhatsApp) or page id
	AccessToken string // env
	Timeout     time.Duration
}

// Validate enforces a safe, exact configuration.
func (c GraphConfig) Validate() error {
	u, err := url.Parse(c.GraphURL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return errors.New("msgchannels: graph url must be absolute http(s)")
	}
	if strings.TrimSpace(c.Version) == "" || strings.TrimSpace(c.SenderID) == "" {
		return errors.New("msgchannels: version and sender id required")
	}
	if strings.TrimSpace(c.AccessToken) == "" {
		return errors.New("msgchannels: access token required (from environment)")
	}
	if c.Timeout <= 0 {
		return errors.New("msgchannels: timeout required")
	}
	return nil
}

// SendText sends a text message via the Meta Graph API.
func SendText(ctx context.Context, cfg GraphConfig, to, text string) error {
	if err := cfg.Validate(); err != nil {
		return err
	}
	if strings.TrimSpace(to) == "" || strings.TrimSpace(text) == "" {
		return errors.New("msgchannels: recipient and text required")
	}
	payload := map[string]any{
		"messaging_product": "whatsapp",
		"to":                to,
		"type":              "text",
		"text":              map[string]string{"body": text},
	}
	body, _ := json.Marshal(payload)
	endpoint := strings.TrimRight(cfg.GraphURL, "/") + "/" + cfg.Version + "/" + cfg.SenderID + "/messages"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("authorization", "Bearer "+cfg.AccessToken)
	client := &http.Client{Timeout: cfg.Timeout}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		data, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("msgchannels: graph status %d: %s", resp.StatusCode, strings.TrimSpace(string(data)))
	}
	return nil
}

// InboundMessage is a parsed Meta webhook message.
type InboundMessage struct {
	From    string
	Channel string // "meta" (WhatsApp/IG/FB share the Graph webhook)
	Text    string
}

// ParseWebhook extracts inbound text messages from a Meta webhook payload.
func ParseWebhook(raw []byte) ([]InboundMessage, error) {
	var root struct {
		Entry []struct {
			Messaging []struct {
				Sender struct {
					ID string `json:"id"`
				} `json:"sender"`
				Message struct {
					Text string `json:"text"`
				} `json:"message"`
			} `json:"messaging"`
		} `json:"entry"`
	}
	if err := json.Unmarshal(raw, &root); err != nil {
		return nil, err
	}
	var out []InboundMessage
	for _, e := range root.Entry {
		for _, m := range e.Messaging {
			if m.Sender.ID != "" && m.Message.Text != "" {
				out = append(out, InboundMessage{From: m.Sender.ID, Channel: "meta", Text: m.Message.Text})
			}
		}
	}
	return out, nil
}
````

### FILE: `internal/msgchannels/email.go`
```yaml
block_id: "GO-MESSAGING-CHANNEL-ADAPTERS:internal/msgchannels/email.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "c9cb0931a2ef3bb68c61a13542b05cb01d253e7fe0bf9dde2a0ee9d50f98559a"
variables: []
secrets_allowed: false
```
````go
package msgchannels

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
)

// EmailConfig carries the email-relay HTTP API settings (SendGrid-compatible).
// APIKey is a secret from the environment, never hardcoded.
type EmailConfig struct {
	APIURL  string // e.g. https://api.sendgrid.com/v3/mail/send
	APIKey  string // env
	From    string
	Timeout time.Duration
}

// Validate enforces a safe configuration.
func (c EmailConfig) Validate() error {
	u, err := url.Parse(c.APIURL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return errors.New("msgchannels: api url must be absolute http(s)")
	}
	if strings.TrimSpace(c.APIKey) == "" {
		return errors.New("msgchannels: api key required (from environment)")
	}
	if !strings.Contains(c.From, "@") {
		return errors.New("msgchannels: from must be an email address")
	}
	if c.Timeout <= 0 {
		return errors.New("msgchannels: timeout required")
	}
	return nil
}

type emailPayload struct {
	Personalizations []struct {
		To []struct {
			Email string `json:"email"`
		} `json:"to"`
	} `json:"personalizations"`
	From struct {
		Email string `json:"email"`
	} `json:"from"`
	Subject string `json:"subject"`
	Content []struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"content"`
}

// SendEmail sends a text/plain email via an HTTP relay (SendGrid-compatible).
func SendEmail(ctx context.Context, cfg EmailConfig, to, subject, body string) error {
	if err := cfg.Validate(); err != nil {
		return err
	}
	if strings.TrimSpace(to) == "" || strings.TrimSpace(subject) == "" || strings.TrimSpace(body) == "" {
		return errors.New("msgchannels: to, subject and body required")
	}
	var p emailPayload
	p.Personalizations = []struct {
		To []struct {
			Email string `json:"email"`
		} `json:"to"`
	}{{To: []struct {
		Email string `json:"email"`
	}{{Email: to}}}}
	p.From.Email = cfg.From
	p.Subject = subject
	p.Content = []struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	}{{Type: "text/plain", Value: body}}

	raw, err := json.Marshal(p)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, cfg.APIURL, bytes.NewReader(raw))
	if err != nil {
		return err
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("authorization", "Bearer "+cfg.APIKey)
	client := &http.Client{Timeout: cfg.Timeout}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		data, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("msgchannels: email status %d: %s", resp.StatusCode, strings.TrimSpace(string(data)))
	}
	return nil
}
````

### FILE: `internal/msgchannels/metagraph_test.go`
```yaml
block_id: "GO-MESSAGING-CHANNEL-ADAPTERS:internal/msgchannels/metagraph_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "0a612b41352d56bffa4ec4497333e3fd68fcd95c735083099b59422f2e2dfaf7"
variables: []
secrets_allowed: false
```
````go
package msgchannels

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func sign(secret string, body []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}

func TestVerifyWebhookSignature(t *testing.T) {
	body := []byte(`{"object":"page"}`)
	secret := "app-secret"
	if !VerifyWebhookSignature(secret, body, sign(secret, body)) {
		t.Fatal("valid signature rejected")
	}
	if VerifyWebhookSignature(secret, body, "sha256=deadbeef") {
		t.Fatal("wrong signature accepted")
	}
	if VerifyWebhookSignature(secret, body, "sha1=deadbeef") {
		t.Fatal("non-sha256 accepted")
	}
	if VerifyWebhookSignature("", body, sign(secret, body)) {
		t.Fatal("empty secret accepted")
	}
}

func TestParseWebhook(t *testing.T) {
	raw := []byte(`{"entry":[{"messaging":[{"sender":{"id":"55119999"},"message":{"text":"hola"}}]}]}`)
	msgs, err := ParseWebhook(raw)
	if err != nil {
		t.Fatal(err)
	}
	if len(msgs) != 1 || msgs[0].From != "55119999" || msgs[0].Text != "hola" {
		t.Fatalf("unexpected parse: %+v", msgs)
	}
}

func TestSendTextRoundTrip(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v22.0/phone-1/messages" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("authorization") != "Bearer tok" {
			t.Errorf("missing bearer")
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body["to"] != "55119999" || body["type"] != "text" {
			t.Errorf("unexpected body: %+v", body)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message_id":"m1"}`))
	}))
	defer srv.Close()

	cfg := GraphConfig{GraphURL: srv.URL, Version: "v22.0", SenderID: "phone-1", AccessToken: "tok", Timeout: 5 * time.Second}
	if err := SendText(context.Background(), cfg, "55119999", "hola"); err != nil {
		t.Fatal(err)
	}
}

func TestGraphConfigValidate(t *testing.T) {
	cfg := GraphConfig{GraphURL: "https://graph.facebook.com", Version: "v22.0", SenderID: "p1", AccessToken: "t", Timeout: time.Second}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("valid config rejected: %v", err)
	}
	bad := cfg
	bad.AccessToken = ""
	if err := bad.Validate(); err == nil {
		t.Fatal("missing token accepted")
	}
}
````

### FILE: `internal/msgchannels/email_test.go`
```yaml
block_id: "GO-MESSAGING-CHANNEL-ADAPTERS:internal/msgchannels/email_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "b941460fe1e6df733e0661fb310666e549b8bdfa42bc519d5e19548a25a66c07"
variables: []
secrets_allowed: false
```
````go
package msgchannels

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestEmailConfigValidate(t *testing.T) {
	cfg := EmailConfig{APIURL: "https://api.sendgrid.com/v3/mail/send", APIKey: "k", From: "ventas@empresa.com", Timeout: time.Second}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("valid config rejected: %v", err)
	}
	bad := cfg
	bad.From = "no-arroba"
	if err := bad.Validate(); err == nil {
		t.Fatal("bad from accepted")
	}
	noKey := cfg
	noKey.APIKey = ""
	if err := noKey.Validate(); err == nil {
		t.Fatal("missing key accepted")
	}
}

func TestSendEmailRoundTrip(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("authorization") != "Bearer k" {
			t.Errorf("missing bearer")
		}
		var p emailPayload
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			t.Fatal(err)
		}
		if p.From.Email != "ventas@empresa.com" || p.Subject != "Oferta" {
			t.Errorf("unexpected payload: %+v", p)
		}
		if len(p.Personalizations) != 1 || p.Personalizations[0].To[0].Email != "cliente@x.com" {
			t.Errorf("unexpected recipient: %+v", p.Personalizations)
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	cfg := EmailConfig{APIURL: srv.URL, APIKey: "k", From: "ventas@empresa.com", Timeout: 5 * time.Second}
	if err := SendEmail(context.Background(), cfg, "cliente@x.com", "Oferta", "Hola"); err != nil {
		t.Fatal(err)
	}
}
````


## 6. Configuration surface

| Variable | Tipo | Secreto | Efecto |
|---|---|---|---|
| `GraphConfig.AccessToken` | string | sí | Bearer del Graph API |
| `AppSecret` (verify) | string | sí | firma del webhook |
| `EmailConfig.APIKey` | string | sí | Bearer del relay email |
| `Timeout` | duration | no | deadline de cada llamada |

## 7. Dependency bill

| Package/image/tool | Pin exacto | Uso | Licencia | Runtime/build | Fuente oficial |
|---|---|---|---|---|---|
| Go standard library | go 1.26.7 (toolchain) | net/http, crypto/hmac, encoding/json | BSD-3-Clause | runtime | https://go.dev |

## 8. Apply order

1. Componer `GO-CHANNELS-CORE` (mismo módulo).
2. Colocar los cuatro archivos bajo `internal/msgchannels/`.
3. Inyectar `AccessToken`/`APIKey`/`AppSecret` del entorno del target.
4. Verificar con `go test ./... -count=1` y `go vet ./...`.
5. Rollback: eliminar `internal/msgchannels/`; no deja estado.

## 9. Verification

- `go test ./internal/msgchannels/ -count=1`: 6/6 PASS (firma HMAC válida/inválida/no-sha256, parseo de webhook, envío de texto con servidor falso verificando path/Bearer/payload, config de Graph con negativos, envío de email con relay falso, config de email con negativos).
- `go test ./... -count=1` (11 paquetes): PASS.
- `go vet ./...`: exit 0.

## 10. Reconstruction evidence

Véase `reconstruction_evidence/GO_MESSAGING_CHANNEL_ADAPTERS_2026-09-02_V187.md`.
