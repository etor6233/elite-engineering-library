package paymentbridge

// AUTHORED composition glue: signature implementations remain in the pinned
// SDK adapter; persistence reuses ProviderIntegration's durable inbox/job owner.
import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net/http"
	"net/url"
	"strings"
	"time"

	"elite.local/enterprise/internal/providerintegration"
	officialpayments "example.com/elite/official-payment-webhooks/officialpayments"
)

var ErrWebhookConfig = errors.New("payment webhook configuration invalid")

type SecretSource interface {
	PaymentWebhookSecret(context.Context, string, string) (string, error)
}
type WebhookConfig struct {
	TenantID, ConnectionID, ProviderCode string
	StripeAccountID                      string
	LiveMode                             bool
	Tolerance                            time.Duration
	MaxConcurrent                        int
	Secrets                              SecretSource
	Inbox                                providerintegration.Repository
}
type Webhook struct {
	config WebhookConfig
	slots  chan struct{}
}

func NewWebhook(c WebhookConfig) (*Webhook, error) {
	valid := func(s string) bool {
		return s != "" && len(s) <= 200 && strings.TrimSpace(s) == s && !strings.ContainsAny(s, "\r\n")
	}
	if !valid(c.TenantID) || !valid(c.ConnectionID) || (c.ProviderCode != "stripe" && c.ProviderCode != "mercadopago") || c.Tolerance <= 0 || c.Tolerance > 10*time.Minute || c.MaxConcurrent < 1 || c.MaxConcurrent > 32 || c.Secrets == nil || c.Inbox == nil {
		return nil, ErrWebhookConfig
	}
	return &Webhook{config: c, slots: make(chan struct{}, c.MaxConcurrent)}, nil
}

// InboxPayload contains no provider raw body, client secret, signature or payer
// data. BodySHA256 in the receipt binds original bytes. The normalized notice
// only queues GET reconciliation; it never authorizes a payment transition.
type InboxPayload struct {
	Schema       string                               `json:"schema"`
	Notification officialpayments.PaymentNotification `json:"notification"`
}

func (h *Webhook) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fail := func(code int) { w.WriteHeader(code); _, _ = io.WriteString(w, "PAYMENT_WEBHOOK_NOT_ACCEPTED\n") }
	if h == nil {
		fail(503)
		return
	}
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
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
	media, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || media != "application/json" || len(r.Header.Values("Content-Type")) != 1 {
		fail(415)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	_ = http.NewResponseController(w).SetReadDeadline(time.Now().Add(10 * time.Second))
	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, officialpayments.MaxPayloadBytes))
	if err != nil {
		fail(413)
		return
	}
	if len(body) == 0 {
		fail(400)
		return
	}
	one := func(key string) (string, bool) {
		values := r.Header.Values(key)
		return r.Header.Get(key), len(values) == 1 && len(values[0]) > 0 && len(values[0]) <= 16384
	}
	var sig, requestID, dataID string
	if h.config.ProviderCode == "stripe" {
		var ok bool
		sig, ok = one("Stripe-Signature")
		if !ok || r.URL.RawQuery != "" {
			fail(400)
			return
		}
	} else {
		var ok bool
		sig, ok = one("X-Signature")
		if !ok {
			fail(400)
			return
		}
		requestID, ok = one("X-Request-Id")
		if !ok || len(r.URL.RawQuery) > 2048 {
			fail(400)
			return
		}
		query, parseErr := url.ParseQuery(r.URL.RawQuery)
		if parseErr != nil || len(query["data.id"]) != 1 || len(query) > 2 {
			fail(400)
			return
		}
		for key, values := range query {
			if (key != "data.id" && key != "type") || len(values) != 1 {
				fail(400)
				return
			}
		}
		dataID = query.Get("data.id")
	}
	secret, err := h.config.Secrets.PaymentWebhookSecret(ctx, h.config.ProviderCode, h.config.ConnectionID)
	if err != nil || secret == "" || len(secret) > 16384 {
		fail(503)
		return
	}
	var n officialpayments.PaymentNotification
	if h.config.ProviderCode == "stripe" {
		n, err = officialpayments.VerifyStripePaymentNotification(body, sig, secret, h.config.Tolerance, h.config.LiveMode, h.config.StripeAccountID)
	} else {
		n, err = officialpayments.VerifyMercadoPagoPaymentNotification(body, sig, requestID, dataID, secret, h.config.Tolerance)
	}
	if err != nil {
		fail(403)
		return
	}
	normalized, err := json.Marshal(InboxPayload{Schema: "elite-payment-notification/v1", Notification: n})
	if err != nil {
		fail(503)
		return
	}
	_, err = h.config.Inbox.AcceptWebhook(ctx, providerintegration.Receipt{TenantID: h.config.TenantID, ConnectionID: h.config.ConnectionID, ProviderCode: n.ProviderCode, ProviderEventID: n.ProviderEventID, EventType: "payment.reconciliation-requested.v1", BodyHash: n.BodySHA256, Payload: normalized})
	if errors.Is(err, providerintegration.ErrConflict) {
		fail(409)
		return
	}
	if err != nil {
		fail(503)
		return
	}
	w.WriteHeader(200)
	_, _ = io.WriteString(w, "PAYMENT_WEBHOOK_ACCEPTED\n")
}
