# Go SMS Channel Adapter

## 1. Metadata

```yaml
pack_id: "GO-SMS-CHANNEL-ADAPTER"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa el adapter SMS (Twilio-compatible) con basic auth y body form-encoded; credenciales inyectadas en runtime y tests con servidor falso."
stacks: ["Go 1.26.7"]
compatible_with: ["GO-MESSAGING-CHANNEL-ADAPTERS 0.1.x", "GO-CHANNELS-CORE 0.1.x"]
incompatible_with: ["credenciales en código", "endpoint sin URL absoluta http(s)"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://www.twilio.com/docs/messaging/api"]
verified_at: "2026-09-02"
```

Completa los canales de mensajería (WhatsApp/IG/FB + email + SMS). Las credenciales (AccountSid/AuthToken) se inyectan del entorno; los tests usan `httptest` y no requieren red ni cuenta.

## 2. Applicability

Use este pack para enviar SMS/Messages desde el agente. Rechace para credenciales embebidas.

## 3. Architecture contract

- **Ownership**: `internal/msgchannels` gobierna el envío SMS; `internal/channels` gobierna el ruteo.
- **Invariantes**: (1) basic auth con credenciales del entorno. (2) status no-2xx → error (fail-closed). (3) URL absoluta http(s).
- **Data flow**: `SendSMS` → POST form-encoded con basic auth → respuesta.
- **Failure modes**: credencial vacía, status error → error.
- **Seguridad/privacidad**: basic auth; sin logs de credenciales.
- **Performance budget**: timeout por llamada.

## 4. Exact file manifest

```text
CREATE internal/msgchannels/sms.go
CREATE internal/msgchannels/sms_test.go
```

## 5. Materialization blocks

### FILE: `internal/msgchannels/sms.go`
```yaml
block_id: "GO-SMS-CHANNEL-ADAPTER:internal/msgchannels/sms.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "f24a4c06c6c7c9c288506f9c1809e9d092ff689fd8623379bb5be25a1c0a1617"
variables: []
secrets_allowed: false
```
````go
package msgchannels

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// SMSConfig carries the Twilio-compatible SMS API settings. AuthToken is a
// secret from the environment, never hardcoded.
type SMSConfig struct {
	APIURL     string // https://api.twilio.com/2010-04-01/Accounts/{sid}/Messages.json
	AccountSid string // env (basic auth username)
	AuthToken  string // env (basic auth password)
	From       string // sender phone number / messaging service sid
	Timeout    time.Duration
}

// Validate enforces a safe configuration.
func (c SMSConfig) Validate() error {
	u, err := url.Parse(c.APIURL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return errors.New("msgchannels: sms url must be absolute http(s)")
	}
	if strings.TrimSpace(c.AccountSid) == "" || strings.TrimSpace(c.AuthToken) == "" {
		return errors.New("msgchannels: account sid and auth token required (from environment)")
	}
	if strings.TrimSpace(c.From) == "" {
		return errors.New("msgchannels: from required")
	}
	if c.Timeout <= 0 {
		return errors.New("msgchannels: timeout required")
	}
	return nil
}

// SendSMS sends a text SMS via a Twilio-compatible API (basic auth + form body).
func SendSMS(ctx context.Context, cfg SMSConfig, to, body string) error {
	if err := cfg.Validate(); err != nil {
		return err
	}
	if strings.TrimSpace(to) == "" || strings.TrimSpace(body) == "" {
		return errors.New("msgchannels: to and body required")
	}
	form := url.Values{}
	form.Set("To", to)
	form.Set("From", cfg.From)
	form.Set("Body", body)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, cfg.APIURL, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(cfg.AccountSid, cfg.AuthToken)

	client := &http.Client{Timeout: cfg.Timeout}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		data, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("msgchannels: sms status %d: %s", resp.StatusCode, strings.TrimSpace(string(data)))
	}
	return nil
}
````

### FILE: `internal/msgchannels/sms_test.go`
```yaml
block_id: "GO-SMS-CHANNEL-ADAPTER:internal/msgchannels/sms_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "4b8e9093b8eae56b705faf7decf8bab3332435ee0aadc06c56edc5d88b4c23ba"
variables: []
secrets_allowed: false
```
````go
package msgchannels

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSMSConfigValidate(t *testing.T) {
	cfg := SMSConfig{APIURL: "https://api.twilio.com/2010-04-01/Accounts/sid/Messages.json", AccountSid: "sid", AuthToken: "tok", From: "+123", Timeout: time.Second}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("valid config rejected: %v", err)
	}
	bad := cfg
	bad.AuthToken = ""
	if err := bad.Validate(); err == nil {
		t.Fatal("missing auth token accepted")
	}
}

func TestSendSMSRoundTrip(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		user, pass, ok := r.BasicAuth()
		if !ok || user != "sid" || pass != "tok" {
			t.Errorf("bad basic auth: %q %q %v", user, pass, ok)
		}
		_ = r.ParseForm()
		if r.Form.Get("To") != "+55119999" || r.Form.Get("Body") != "hola" {
			t.Errorf("unexpected form: %+v", r.Form)
		}
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"sid":"SM1"}`))
	}))
	defer srv.Close()

	cfg := SMSConfig{APIURL: srv.URL, AccountSid: "sid", AuthToken: "tok", From: "+123", Timeout: 5 * time.Second}
	if err := SendSMS(context.Background(), cfg, "+55119999", "hola"); err != nil {
		t.Fatal(err)
	}
}
````


## 6. Configuration surface

| Variable | Tipo | Secreto | Efecto |
|---|---|---|---|
| `APIURL` | URL http(s) | no | endpoint |
| `AccountSid` | string | sí | basic auth username |
| `AuthToken` | string | sí | basic auth password |
| `From` | string | no | remitente |
| `Timeout` | duration | no | deadline |

## 7. Dependency bill

| Package/image/tool | Pin exacto | Uso | Licencia | Runtime/build | Fuente oficial |
|---|---|---|---|---|---|
| Go standard library | go 1.26.7 (toolchain) | net/http | BSD-3-Clause | runtime | https://go.dev |

## 8. Apply order

1. Componer `GO-MESSAGING-CHANNEL-ADAPTERS` (mismo paquete).
2. Colocar los dos archivos bajo `internal/msgchannels/`.
3. Inyectar `AccountSid`/`AuthToken` del entorno.
4. Verificar con `go test ./... -count=1` y `go vet ./...`.
5. Rollback: eliminar los dos archivos.

## 9. Verification

- `go test ./internal/msgchannels/ -count=1`: 8/8 PASS (suma SMS config + envío con basic auth/form).
- `go test ./... -count=1` (12 paquetes): PASS.
- `go vet ./...`: exit 0.

## 10. Reconstruction evidence

Véase `reconstruction_evidence/GO_SMS_CHANNEL_ADAPTER_2026-09-02_V189.md`.
