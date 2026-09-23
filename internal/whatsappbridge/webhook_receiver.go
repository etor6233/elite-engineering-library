package whatsappbridge

// AUTHORED HTTP/SQL composition. Meta authentication remains in the pinned
// Python adapter; this file does not invent another provider signature scheme.
import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"time"

	"elite.local/enterprise/internal/providerintegration"
)

var ErrWebhookIngress = errors.New("whatsappbridge: ingress unavailable or unverified")
var webhookConnection = regexp.MustCompile(`^[a-z0-9][a-z0-9._-]{0,127}$`)

type VerifyTokenSource interface {
	WhatsAppVerifyToken(context.Context) (string, error)
}

// RetentionApprovalSHA256 links the project's approved raw-payload policy.
// Base64 is NOT encryption. Protect the database/backups and restrict/purge raw
// PII under that policy before exposing this handler. No default activation.
type WebhookReceiverConfig struct {
	TenantID, ConnectionID, RetentionApprovalSHA256 string
	Profile                                         json.RawMessage
	Process                                         Process
	Secrets                                         AppSecretSource
	Verification                                    VerifyTokenSource
	Store                                           providerintegration.Repository
	MaxConcurrent                                   int
}

type WebhookReceiver struct {
	config WebhookReceiverConfig
	slots  chan struct{}
}

func NewWebhookReceiver(c WebhookReceiverConfig) (*WebhookReceiver, error) {
	if !notificationEventID.MatchString(c.TenantID) || !webhookConnection.MatchString(c.ConnectionID) || !validDigest(c.RetentionApprovalSHA256) || len(c.Profile) > 32768 || !json.Valid(c.Profile) || c.Store == nil || c.Secrets == nil || c.Verification == nil || c.MaxConcurrent < 1 || c.MaxConcurrent > 16 {
		return nil, ErrWebhookIngress
	}
	c.Profile = append(json.RawMessage(nil), c.Profile...)
	return &WebhookReceiver{config: c, slots: make(chan struct{}, c.MaxConcurrent)}, nil
}

type ingressResult struct {
	Schema    string `json:"schema"`
	Binding   string `json:"binding_sha256"`
	BodyHash  string `json:"body_sha256"`
	Count     int    `json:"event_count"`
	Challenge string `json:"challenge"`
}

func (h *WebhookReceiver) verify(ctx context.Context, mode string, body []byte, signature, secret string, query map[string]string) (ingressResult, error) {
	var out ingressResult
	p := h.config.Process
	script := filepath.Join(p.AdapterDirectory, "whatsapp_cloud.py")
	if !filepath.IsAbs(p.AdapterDirectory) || exactFile(p.PythonExecutable, p.PythonSHA256) != nil || exactFile(script, p.AdapterSHA256) != nil || secret == "" || len(secret) > 16384 {
		return out, ErrWebhookIngress
	}
	frame := struct {
		Schema    string            `json:"schema"`
		Mode      string            `json:"mode"`
		Profile   json.RawMessage   `json:"profile"`
		Body      []byte            `json:"body"`
		Signature string            `json:"signature"`
		Secret    string            `json:"secret"`
		Query     map[string]string `json:"query"`
	}{"elite-whatsapp-ingress-bridge/v1", mode, h.config.Profile, body, signature, secret, query}
	input, err := json.Marshal(frame)
	if err != nil || len(input) > 2<<20 {
		return out, ErrWebhookIngress
	}
	cmd := exec.CommandContext(ctx, p.PythonExecutable, "-I", "-B", script, "--ingress-bridge")
	cmd.Env = []string{}
	if root := os.Getenv("SystemRoot"); root != "" {
		cmd.Env = append(cmd.Env, "SystemRoot="+root)
	}
	cmd.Stdin = bytes.NewReader(input)
	cmd.Stderr = io.Discard
	cmd.WaitDelay = time.Second
	var output boundedOutput
	cmd.Stdout = &output
	if cmd.Run() != nil || output.overflow {
		return out, ErrWebhookIngress
	}
	decoder := json.NewDecoder(bytes.NewReader(output.Bytes()))
	decoder.DisallowUnknownFields()
	if decoder.Decode(&out) != nil || decoder.Decode(new(any)) != io.EOF || out.Schema != "elite-whatsapp-ingress-result/v1" || out.Binding != digest(input) {
		return ingressResult{}, ErrWebhookIngress
	}
	return out, nil
}

// Mount on one server-configured callback URL. No URL/header/body selects the
// tenant or connection. Configure server read/header/idle timeouts and TLS/edge
// protections; redact subscription query tokens in access logs.
func (h *WebhookReceiver) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fail := func(code int) { w.WriteHeader(code); _, _ = io.WriteString(w, "WEBHOOK_NOT_ACCEPTED\n") }
	if h == nil {
		fail(503)
		return
	}
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		w.Header().Set("Allow", "GET, POST")
		fail(405)
		return
	}
	select {
	case h.slots <- struct{}{}:
		defer func() { <-h.slots }()
	default:
		fail(503)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	_ = http.NewResponseController(w).SetReadDeadline(time.Now().Add(5 * time.Second))
	if r.Method == http.MethodGet {
		if len(r.URL.RawQuery) > 2048 {
			fail(400)
			return
		}
		// ParseQuery reports malformed escaping instead of silently dropping it.
		values, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil || len(values) != 3 {
			fail(400)
			return
		}
		query := map[string]string{}
		for _, name := range []string{"hub.mode", "hub.verify_token", "hub.challenge"} {
			if len(values[name]) != 1 {
				fail(400)
				return
			}
			query[name] = values[name][0]
		}
		secret, err := h.config.Verification.WhatsAppVerifyToken(ctx)
		if err != nil {
			fail(503)
			return
		}
		result, err := h.verify(ctx, "subscribe", []byte{}, "", secret, query)
		if err != nil || result.Challenge != query["hub.challenge"] || result.Challenge == "" || result.Count != 0 || result.BodyHash != "" {
			fail(403)
			return
		}
		w.WriteHeader(200)
		_, _ = io.WriteString(w, result.Challenge)
		return
	}
	if r.URL.RawQuery != "" || len(r.Header.Values("X-Hub-Signature-256")) != 1 || len(r.Header.Get("X-Hub-Signature-256")) != 71 {
		fail(400)
		return
	}
	media, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || media != "application/json" || len(r.Header.Values("Content-Type")) != 1 || r.Header.Get("Content-Encoding") != "" {
		fail(415)
		return
	}
	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 1<<20))
	if err != nil || len(body) == 0 {
		fail(413)
		return
	}
	secret, err := h.config.Secrets.WhatsAppAppSecret(ctx)
	if err != nil {
		fail(503)
		return
	}
	signature := r.Header.Get("X-Hub-Signature-256")
	result, err := h.verify(ctx, "receive", body, signature, secret, map[string]string{})
	if err != nil || result.BodyHash != digest(body) || result.Count < 1 || result.Count > 1000 || result.Challenge != "" {
		fail(403)
		return
	}
	// Deterministic payload preserves exact original bytes for re-verification.
	// Admission does not mark the job processed or mutate CRM/delivery status.
	payload, err := json.Marshal(struct {
		Schema        string `json:"schema"`
		Body          []byte `json:"body"`
		Signature     string `json:"signature"`
		ProfileHash   string `json:"profile_sha256"`
		RetentionHash string `json:"retention_approval_sha256"`
		Count         int    `json:"verified_event_count"`
	}{"elite-whatsapp-retained-webhook/v1", body, signature, digest(h.config.Profile), h.config.RetentionApprovalSHA256, result.Count})
	if err != nil {
		fail(503)
		return
	}
	_, err = h.config.Store.AcceptWebhook(ctx, providerintegration.Receipt{TenantID: h.config.TenantID, ConnectionID: h.config.ConnectionID, ProviderCode: "meta-whatsapp", ProviderEventID: "wa:" + result.BodyHash, EventType: "whatsapp.raw_webhook.received.v1", BodyHash: result.BodyHash, Payload: payload})
	if err != nil {
		fail(503)
		return
	}
	w.WriteHeader(200)
	_, _ = io.WriteString(w, "EVENT_RECEIVED\n")
}
