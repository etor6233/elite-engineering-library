package main

// AUTHORED host composition. Secrets are future environment inputs and never
// appear in configuration errors, URLs, logs, or customer responses.
import (
	"context"
	"encoding/base64"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"elite.local/enterprise/internal/paymentbridge"
	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/platform/workers"
	"github.com/jackc/pgx/v5/pgxpool"
)

var errPaymentConfiguration = errors.New("payment checkout configuration is incomplete or invalid")

type paymentConfiguration struct {
	Driver                         paymentbridge.SDKDriverConfig
	Token, WebhookSecret, WorkerID string
	HMACKey                        []byte
}

func loadPaymentConfiguration(lookup func(string) string) (*paymentConfiguration, error) {
	if lookup == nil {
		return nil, errPaymentConfiguration
	}
	switch lookup("PAYMENT_CHECKOUT_ENABLED") {
	case "", "false":
		return nil, nil
	case "true":
	default:
		return nil, errPaymentConfiguration
	}
	scope := paymentbridge.Scope{TenantID: lookup("PAYMENT_TENANT_ID"), OrganizationID: lookup("PAYMENT_ORGANIZATION_ID"), ConnectionID: lookup("PAYMENT_CONNECTION_ID"), ProviderCode: lookup("PAYMENT_REQUEST_PROVIDER"), AccountRef: lookup("PAYMENT_ACCOUNT_REF"), Currency: lookup("PAYMENT_CURRENCY")}
	exponent, err := strconv.Atoi(lookup("PAYMENT_MINOR_UNIT_EXPONENT"))
	if err != nil {
		return nil, errPaymentConfiguration
	}
	scope.MinorUnitExponent = exponent
	switch lookup("PAYMENT_LIVE_MODE") {
	case "true":
		scope.LiveMode = true
	case "false":
	default:
		return nil, errPaymentConfiguration
	}
	c := &paymentConfiguration{Driver: paymentbridge.SDKDriverConfig{Scope: scope, SuccessURL: lookup("PAYMENT_SUCCESS_URL"), CancelURL: lookup("PAYMENT_CANCEL_URL"), NotificationURL: lookup("PAYMENT_NOTIFICATION_URL"), DisplayName: lookup("PAYMENT_DISPLAY_NAME")}, Token: lookup("PAYMENT_PROVIDER_SECRET"), WebhookSecret: lookup("PAYMENT_WEBHOOK_SECRET"), WorkerID: lookup("PAYMENT_WORKER_ID")}
	c.HMACKey, err = base64.StdEncoding.Strict().DecodeString(lookup("PAYMENT_OUTBOUND_HMAC_KEY_BASE64"))
	if err != nil || len(c.HMACKey) < 32 || len(c.HMACKey) > 64 || len(c.WebhookSecret) < 16 || len(c.WebhookSecret) > 16384 || strings.TrimSpace(c.WebhookSecret) != c.WebhookSecret || strings.ContainsAny(c.WebhookSecret, "\r\n") || c.WorkerID == "" || len(c.WorkerID) > 128 || strings.TrimSpace(c.WorkerID) != c.WorkerID || strings.ContainsAny(c.WorkerID, "\r\n") {
		return nil, errPaymentConfiguration
	}
	// Construction checks scope, exact key mode and configured HTTPS redirects;
	// identity is subsequently probed through the official SDK before claims.
	if _, err = paymentbridge.NewSDKDriver(c.Driver, c.Token, nil); err != nil {
		return nil, errPaymentConfiguration
	}
	return c, nil
}

type scopedPaymentSecret struct{ provider, connection, secret string }

func (s scopedPaymentSecret) PaymentWebhookSecret(_ context.Context, provider, connection string) (string, error) {
	if provider != s.provider || connection != s.connection {
		return "", errPaymentConfiguration
	}
	return s.secret, nil
}

type paymentRuntime struct {
	module    httpapi.PaymentCheckoutModule
	worker    paymentbridge.Worker
	callbacks workers.JobProcessor
	driver    *paymentbridge.SDKDriver
}

func preparePaymentRuntime(ctx context.Context, pool *pgxpool.Pool, lookup func(string) string, client *http.Client) (*paymentRuntime, error) {
	c, err := loadPaymentConfiguration(lookup)
	if err != nil || c == nil {
		return nil, err
	}
	if ctx == nil || pool == nil {
		return nil, errPaymentConfiguration
	}
	driver, err := paymentbridge.NewSDKDriver(c.Driver, c.Token, client)
	if err != nil {
		return nil, errPaymentConfiguration
	}
	probe, cancel := context.WithTimeout(ctx, 15*time.Second)
	err = driver.ValidateCredential(probe)
	cancel()
	if err != nil {
		return nil, errPaymentConfiguration
	}
	scope := c.Driver.Scope
	store := postgres.NewPaymentCheckoutStore(pool)
	baseFence, err := postgres.NewOutboundDeliveryStore(pool, c.HMACKey, 2*time.Minute)
	if err != nil {
		return nil, errPaymentConfiguration
	}
	worker := paymentbridge.Worker{Scope: scope, Store: store, Fence: &postgres.PaymentDispatchFence{OutboundDeliveryStore: baseFence, Scope: scope}, Driver: driver}
	callback, err := postgres.NewPaymentCallbackProcessor(pool, worker, c.WorkerID)
	if err != nil {
		return nil, errPaymentConfiguration
	}
	reader, err := postgres.NewCustomerCheckoutReader(pool, scope)
	if err != nil {
		return nil, errPaymentConfiguration
	}
	webhook, err := paymentbridge.NewWebhook(paymentbridge.WebhookConfig{TenantID: scope.TenantID, ConnectionID: scope.ConnectionID, ProviderCode: scope.ProviderCode, LiveMode: scope.LiveMode, Tolerance: 5 * time.Minute, MaxConcurrent: 8, Secrets: scopedPaymentSecret{scope.ProviderCode, scope.ConnectionID, c.WebhookSecret}, Inbox: store})
	if err != nil {
		return nil, errPaymentConfiguration
	}
	return &paymentRuntime{module: httpapi.PaymentCheckoutModule{Reader: reader, Webhook: webhook}, worker: worker, driver: driver, callbacks: workers.JobProcessor{Store: &postgres.PaymentCallbackJobs{Jobs: postgres.NewJobs(pool), Scope: scope}, Handler: postgres.PaymentCallbackRecoveryProcessor{Base: callback}, Queue: "payment-provider-events", WorkerID: c.WorkerID, Lease: 2 * time.Minute, RetryDelay: 30 * time.Second, BatchSize: 1}}, nil
}

// One bounded batch per iteration keeps callbacks responsive. Every provider
// mutation still passes the persisted request binding and outbound fence.
func (r *paymentRuntime) run(ctx context.Context) {
	timer := time.NewTimer(0)
	defer timer.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
		}
		iteration, cancel := context.WithTimeout(ctx, 50*time.Second)
		_, sendErr := r.worker.ProcessOnce(iteration, 1)
		_, callbackErr := r.callbacks.ProcessOnce(iteration)
		cancel()
		if ctx.Err() != nil {
			return
		}
		delay := 500 * time.Millisecond
		if sendErr != nil || callbackErr != nil {
			slog.Warn("payment processing deferred; durable state retained")
			delay = 2 * time.Second
		}
		timer.Reset(delay)
	}
}
