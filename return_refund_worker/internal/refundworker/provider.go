package refundworker

import (
	"context"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/mercadopago/sdk-go/pkg/config"
	"github.com/mercadopago/sdk-go/pkg/payment"
	mprefund "github.com/mercadopago/sdk-go/pkg/refund"
	"github.com/mercadopago/sdk-go/pkg/requester"
	"github.com/mercadopago/sdk-go/pkg/requestoptions"
	"github.com/stripe/stripe-go/v86"
)

var (
	idempotencyPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:-]{15,127}$`)
	currencyPattern    = regexp.MustCompile(`^[A-Z]{3}$`)
)

type stripeRefundAPI interface {
	Create(context.Context, *stripe.RefundCreateParams) (*stripe.Refund, error)
	Retrieve(context.Context, string, *stripe.RefundRetrieveParams) (*stripe.Refund, error)
}

type StripeProvider struct{ refunds stripeRefundAPI }

func NewStripeProvider(secretKey string) (*StripeProvider, error) {
	if strings.TrimSpace(secretKey) == "" {
		return nil, ErrInvalidProvider
	}
	return &StripeProvider{refunds: stripe.NewClient(secretKey).V1Refunds}, nil
}

func newStripeProvider(secretKey string, backend stripe.Backend) (*StripeProvider, error) {
	if strings.TrimSpace(secretKey) == "" || backend == nil {
		return nil, ErrInvalidProvider
	}
	backends := &stripe.Backends{API: backend, Connect: backend, Uploads: backend, MeterEvents: backend}
	return &StripeProvider{refunds: stripe.NewClient(secretKey, stripe.WithBackends(backends)).V1Refunds}, nil
}

func validateRefundInput(value Refund, provider string) error {
	if value.Provider != provider || value.AmountMinorUnits <= 0 || !currencyPattern.MatchString(value.Currency) || strings.TrimSpace(value.ProviderPaymentReference) == "" || !idempotencyPattern.MatchString(value.IdempotencyKey) {
		return ErrInvalidProvider
	}
	return nil
}

func stripeResult(value Refund, resource *stripe.Refund) (ProviderResult, error) {
	if resource == nil || resource.ID == "" || resource.Amount != value.AmountMinorUnits || strings.ToUpper(string(resource.Currency)) != value.Currency {
		return ProviderResult{}, ErrResponseMismatch
	}
	// The requested reference determines the identity domain. A Stripe refund
	// can contain both IDs; its PaymentIntent is not a substitute for its Charge.
	if value.ProviderRefundReference != "" && resource.ID != value.ProviderRefundReference {
		return ProviderResult{}, ErrResponseMismatch
	}
	paymentReference := ""
	switch {
	case strings.HasPrefix(value.ProviderPaymentReference, "pi_"):
		if resource.PaymentIntent != nil {
			paymentReference = resource.PaymentIntent.ID
		}
	case strings.HasPrefix(value.ProviderPaymentReference, "ch_"):
		if resource.Charge != nil {
			paymentReference = resource.Charge.ID
		}
	}
	if paymentReference != value.ProviderPaymentReference {
		return ProviderResult{}, ErrResponseMismatch
	}
	return ProviderResult{ProviderPaymentReference: paymentReference, ProviderRefundReference: resource.ID, ProviderStatus: string(resource.Status), Currency: value.Currency, AmountMinorUnits: resource.Amount}, nil
}

func (p *StripeProvider) Create(ctx context.Context, value Refund) (ProviderResult, error) {
	if p == nil || p.refunds == nil || ctx == nil || validateRefundInput(value, "stripe") != nil || value.ProviderRefundReference != "" {
		return ProviderResult{}, ErrInvalidProvider
	}
	params := &stripe.RefundCreateParams{Amount: stripe.Int64(value.AmountMinorUnits), Reason: stripe.String(string(stripe.RefundReasonRequestedByCustomer)), Metadata: map[string]string{"return_request_id": value.RequestID}}
	switch {
	case strings.HasPrefix(value.ProviderPaymentReference, "pi_"):
		params.PaymentIntent = stripe.String(value.ProviderPaymentReference)
	case strings.HasPrefix(value.ProviderPaymentReference, "ch_"):
		params.Charge = stripe.String(value.ProviderPaymentReference)
	default:
		return ProviderResult{}, ErrInvalidProvider
	}
	params.SetIdempotencyKey(value.IdempotencyKey)
	resource, err := p.refunds.Create(ctx, params)
	if err != nil {
		return ProviderResult{}, err
	}
	return stripeResult(value, resource)
}

func (p *StripeProvider) Retrieve(ctx context.Context, value Refund) (ProviderResult, error) {
	if p == nil || p.refunds == nil || ctx == nil || validateRefundInput(value, "stripe") != nil || strings.TrimSpace(value.ProviderRefundReference) == "" {
		return ProviderResult{}, ErrInvalidProvider
	}
	resource, err := p.refunds.Retrieve(ctx, value.ProviderRefundReference, nil)
	if err != nil {
		return ProviderResult{}, err
	}
	return stripeResult(value, resource)
}

type mercadoPagoRefundAPI interface {
	CreatePartialRefund(context.Context, int, float64) (*mprefund.Response, error)
	Get(context.Context, int, int) (*mprefund.Response, error)
}

type mercadoPagoPaymentAPI interface {
	Get(context.Context, int) (*payment.Response, error)
}

type MercadoPagoProvider struct {
	refunds  mercadoPagoRefundAPI
	payments mercadoPagoPaymentAPI
}

func NewMercadoPagoProvider(accessToken string) (*MercadoPagoProvider, error) {
	return newMercadoPagoProvider(accessToken, nil)
}

func newMercadoPagoProvider(accessToken string, transport requester.Requester) (*MercadoPagoProvider, error) {
	if strings.TrimSpace(accessToken) == "" {
		return nil, ErrInvalidProvider
	}
	options := []config.Option{config.WithTimeout(10 * time.Second), config.WithMaxRetries(2)}
	if transport != nil {
		options = []config.Option{config.WithHTTPClient(transport)}
	}
	cfg, err := config.New(accessToken, options...)
	if err != nil {
		return nil, err
	}
	return &MercadoPagoProvider{refunds: mprefund.NewClient(cfg), payments: payment.NewClient(cfg)}, nil
}

func parsePositiveInt(value string) (int, error) {
	n, err := strconv.Atoi(value)
	if err != nil || n <= 0 {
		return 0, ErrInvalidProvider
	}
	return n, nil
}

func arsMinorToMajor(value int64) (float64, error) {
	if value <= 0 || value > 9_000_000_000_000_000 {
		return 0, ErrInvalidProvider
	}
	major := float64(value) / 100
	if int64(math.Round(major*100)) != value {
		return 0, ErrInvalidProvider
	}
	return major, nil
}

func arsMajorToMinor(value float64) (int64, error) {
	if value <= 0 || math.IsNaN(value) || math.IsInf(value, 0) {
		return 0, ErrResponseMismatch
	}
	minor := int64(math.Round(value * 100))
	if math.Abs(value-float64(minor)/100) > 0.0000001 {
		return 0, ErrResponseMismatch
	}
	return minor, nil
}

func (p *MercadoPagoProvider) verifyPayment(ctx context.Context, value Refund, paymentID int, requireRefundable bool) error {
	resource, err := p.payments.Get(ctx, paymentID)
	if err != nil {
		return err
	}
	amount, convErr := arsMajorToMinor(resource.TransactionAmount)
	validState := (resource.Status == "approved" && resource.Captured) || (!requireRefundable && resource.Status == "refunded")
	if convErr != nil || resource.ID != paymentID || strings.ToUpper(resource.CurrencyID) != value.Currency || !validState || amount < value.AmountMinorUnits {
		return ErrResponseMismatch
	}
	return nil
}

func mercadoPagoResult(value Refund, paymentID int, resource *mprefund.Response) (ProviderResult, error) {
	if resource == nil || resource.ID <= 0 || resource.PaymentID != paymentID {
		return ProviderResult{}, ErrResponseMismatch
	}
	amount, err := arsMajorToMinor(resource.Amount)
	if err != nil || amount != value.AmountMinorUnits {
		return ProviderResult{}, ErrResponseMismatch
	}
	return ProviderResult{ProviderPaymentReference: strconv.Itoa(paymentID), ProviderRefundReference: strconv.Itoa(resource.ID), ProviderStatus: resource.Status, Currency: value.Currency, AmountMinorUnits: amount}, nil
}

func (p *MercadoPagoProvider) Create(ctx context.Context, value Refund) (ProviderResult, error) {
	if p == nil || p.refunds == nil || p.payments == nil || ctx == nil || validateRefundInput(value, "mercado_pago") != nil || value.Currency != "ARS" || value.ProviderRefundReference != "" {
		return ProviderResult{}, ErrInvalidProvider
	}
	paymentID, err := parsePositiveInt(value.ProviderPaymentReference)
	if err != nil {
		return ProviderResult{}, err
	}
	if err = p.verifyPayment(ctx, value, paymentID, true); err != nil {
		return ProviderResult{}, err
	}
	amount, err := arsMinorToMajor(value.AmountMinorUnits)
	if err != nil {
		return ProviderResult{}, err
	}
	resource, err := p.refunds.CreatePartialRefund(requestoptions.WithIdempotencyKey(ctx, value.IdempotencyKey), paymentID, amount)
	if err != nil {
		return ProviderResult{}, err
	}
	return mercadoPagoResult(value, paymentID, resource)
}

func (p *MercadoPagoProvider) Retrieve(ctx context.Context, value Refund) (ProviderResult, error) {
	if p == nil || p.refunds == nil || p.payments == nil || ctx == nil || validateRefundInput(value, "mercado_pago") != nil || value.Currency != "ARS" || value.ProviderRefundReference == "" {
		return ProviderResult{}, ErrInvalidProvider
	}
	paymentID, err := parsePositiveInt(value.ProviderPaymentReference)
	if err != nil {
		return ProviderResult{}, err
	}
	refundID, err := parsePositiveInt(value.ProviderRefundReference)
	if err != nil {
		return ProviderResult{}, err
	}
	if err = p.verifyPayment(ctx, value, paymentID, false); err != nil {
		return ProviderResult{}, err
	}
	resource, err := p.refunds.Get(ctx, paymentID, refundID)
	if err != nil {
		return ProviderResult{}, err
	}
	return mercadoPagoResult(value, paymentID, resource)
}
