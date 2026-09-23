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
