package commerce

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"
)

var ErrConflict = errors.New("commerce conflict")
var ErrPaymentUnavailable = errors.New("payment intent unavailable")
var currencyPattern = regexp.MustCompile(`^[A-Z]{3}$`)
var marketPattern = regexp.MustCompile(`^[A-Z]{2}$`)

type PriceEntry struct {
	VariantID        string `json:"variant_id"`
	AmountMinorUnits int64  `json:"amount_minor_units"`
	TaxMode          string `json:"tax_mode"`
}
type PriceBook struct {
	ID         string       `json:"id"`
	Market     string       `json:"market"`
	Currency   string       `json:"currency"`
	ValidFrom  time.Time    `json:"valid_from"`
	ValidUntil *time.Time   `json:"valid_until,omitempty"`
	Status     string       `json:"status"`
	Entries    []PriceEntry `json:"entries"`
}
type OrderLine struct {
	OrderID             string `json:"order_id"`
	LineID              string `json:"line_id"`
	OrganizationID      string `json:"organization_id"`
	PriceBookID         string `json:"price_book_id"`
	VariantID           string `json:"variant_id"`
	Quantity            int    `json:"quantity"`
	UnitPriceMinorUnits int64  `json:"unit_price_minor_units"`
	OrderVersion        int64  `json:"order_version"`
}
type PaymentAttempt struct {
	ID               string `json:"id"`
	OrderID          string `json:"order_id"`
	OrganizationID   string `json:"organization_id"`
	ProviderCode     string `json:"provider_code"`
	State            string `json:"state"`
	Currency         string `json:"currency"`
	AmountMinorUnits int64  `json:"amount_minor_units"`
	Version          int64  `json:"version"`
}
type Repository interface {
	CreatePriceBook(context.Context, string, string, PriceBook) error
	ActivatePriceBook(context.Context, string, string, string) error
	PublicPrice(context.Context, string, string, string) (PriceEntry, error)
	AddOrderLine(context.Context, string, string, OrderLine, int64) (OrderLine, error)
	PlaceOrder(context.Context, string, string, string, int64, string) error
	AllocateStock(context.Context, string, string, string, string, string, int64, int64, string) error
	CreatePaymentAttempt(context.Context, string, string, string, PaymentAttempt) error
	TransitionPayment(context.Context, string, string, string, string, string, int64, string, string) error
}
type IDGenerator interface{ New() string }

// OperationReader is the read side of the existing commerce owner. It never
// creates a second inventory, payment ledger or delivery state machine.
type OperationReader interface {
	Operations(context.Context, string, string) (OperationSnapshot, error)
}
type AuditedAllocator interface {
	AllocateStockAs(context.Context, string, string, string, string, string, int64, int64, string, string) error
}
type PaymentIntentRecorder interface {
	RecordPaymentIntent(context.Context, string, string, string, PaymentAttempt, string) (PaymentAttempt, error)
}
type OrderPaymentRecorder interface {
	RecordOrderPayment(context.Context, string, string, string, PaymentAttempt, string) (PaymentAttempt, error)
}
type OperationSnapshot struct {
	PaymentProvider string           `json:"payment_provider"`
	Orders          []OperationOrder `json:"orders"`
	Stock           []OperationStock `json:"stock"`
	Truncated       bool             `json:"truncated"`
}
type OperationOrder struct {
	ID        string            `json:"id"`
	State     string            `json:"state"`
	Version   int64             `json:"version"`
	Lines     []OperationLine   `json:"lines"`
	Payments  []OperationStatus `json:"payments"`
	Handovers []OperationStatus `json:"handovers"`
}
type OperationLine struct {
	ID        string `json:"id"`
	VariantID string `json:"variant_id"`
	Name      string `json:"name"`
	Quantity  int    `json:"quantity"`
	StockID   string `json:"stock_id"`
}
type OperationStock struct {
	ID        string `json:"id"`
	VariantID string `json:"variant_id"`
	Serial    string `json:"serial"`
	Version   int64  `json:"version"`
}
type OperationStatus struct {
	ID    string `json:"id"`
	State string `json:"state"`
}
type Service struct {
	repository Repository
	ids        IDGenerator
}

func NewService(r Repository, ids IDGenerator) *Service { return &Service{repository: r, ids: ids} }
func (s *Service) Operations(ctx context.Context, tenant, organization string) (OperationSnapshot, error) {
	reader, ok := s.repository.(OperationReader)
	if !ok || tenant == "" || organization == "" {
		return OperationSnapshot{}, errors.New("operation reader unavailable or invalid scope")
	}
	return reader.Operations(ctx, tenant, organization)
}
func (s *Service) CreatePriceBook(ctx context.Context, tenant string, input PriceBook) (PriceBook, error) {
	if tenant == "" || !marketPattern.MatchString(input.Market) || !currencyPattern.MatchString(input.Currency) || input.ValidFrom.IsZero() || len(input.Entries) == 0 {
		return PriceBook{}, fmt.Errorf("invalid price book")
	}
	seen := map[string]bool{}
	for _, entry := range input.Entries {
		if entry.VariantID == "" || entry.AmountMinorUnits < 0 || !map[string]bool{"exclusive": true, "inclusive": true, "not-applicable": true}[entry.TaxMode] || seen[entry.VariantID] {
			return PriceBook{}, fmt.Errorf("invalid price entry")
		}
		seen[entry.VariantID] = true
	}
	input.ID = s.ids.New()
	input.Status = "draft"
	if err := s.repository.CreatePriceBook(ctx, tenant, s.ids.New(), input); err != nil {
		return PriceBook{}, err
	}
	return input, nil
}
func (s *Service) ActivatePriceBook(ctx context.Context, tenant, id string) error {
	if tenant == "" || id == "" {
		return fmt.Errorf("invalid price book")
	}
	return s.repository.ActivatePriceBook(ctx, tenant, id, s.ids.New())
}
func (s *Service) PublicPrice(ctx context.Context, tenantCode, market, variant string) (PriceEntry, error) {
	if tenantCode == "" || !marketPattern.MatchString(market) || variant == "" {
		return PriceEntry{}, fmt.Errorf("invalid price query")
	}
	return s.repository.PublicPrice(ctx, tenantCode, market, variant)
}
func (s *Service) AddOrderLine(ctx context.Context, tenant string, input OrderLine, expectedVersion int64) (OrderLine, error) {
	if tenant == "" || input.OrderID == "" || input.OrganizationID == "" || input.PriceBookID == "" || input.VariantID == "" || input.Quantity < 1 || expectedVersion < 1 {
		return OrderLine{}, fmt.Errorf("invalid order line")
	}
	input.LineID = s.ids.New()
	return s.repository.AddOrderLine(ctx, tenant, s.ids.New(), input, expectedVersion)
}
func (s *Service) PlaceOrder(ctx context.Context, tenant, organization, orderID string, version int64) error {
	if tenant == "" || organization == "" || orderID == "" || version < 1 {
		return fmt.Errorf("invalid order placement")
	}
	return s.repository.PlaceOrder(ctx, tenant, organization, orderID, version, s.ids.New())
}
func (s *Service) AllocateStock(ctx context.Context, tenant, organization, orderID, lineID, stockID string, orderVersion, stockVersion int64) error {
	if tenant == "" || organization == "" || orderID == "" || lineID == "" || stockID == "" || orderVersion < 1 || stockVersion < 1 {
		return fmt.Errorf("invalid allocation")
	}
	return s.repository.AllocateStock(ctx, tenant, organization, orderID, lineID, stockID, orderVersion, stockVersion, s.ids.New())
}

func (s *Service) AllocateStockAs(ctx context.Context, tenant, organization, orderID, lineID, stockID string, orderVersion, stockVersion int64, actor string) error {
	repo, ok := s.repository.(AuditedAllocator)
	if !ok || actor == "" || tenant == "" || organization == "" || orderID == "" || lineID == "" || stockID == "" || orderVersion < 1 || stockVersion < 1 {
		return fmt.Errorf("invalid audited allocation")
	}
	return repo.AllocateStockAs(ctx, tenant, organization, orderID, lineID, stockID, orderVersion, stockVersion, s.ids.New(), actor)
}

var paymentTransitions = map[string]map[string]bool{"created": {"pending": true, "authorized": true, "failed": true}, "pending": {"authorized": true, "failed": true}, "authorized": {"captured": true, "failed": true}, "captured": {"refunded": true, "disputed": true}, "disputed": {"refunded": true}}

// RequestOrderPayment is the initial whole-order intent only. It does not
// authorize a provider call, another attempt, partial payment or credit policy.
func (s *Service) RequestOrderPayment(ctx context.Context, tenant, organization, orderID, provider, key, actor string) (PaymentAttempt, error) {
	repo, ok := s.repository.(OrderPaymentRecorder)
	if !ok {
		return PaymentAttempt{}, ErrPaymentUnavailable
	}
	if tenant == "" || organization == "" || orderID == "" || actor == "" || (provider != "stripe" && provider != "mercadopago") || len(key) < 16 || len(key) > 128 {
		return PaymentAttempt{}, fmt.Errorf("invalid order payment request")
	}
	value, err := repo.RecordOrderPayment(ctx, tenant, s.ids.New(), key, PaymentAttempt{ID: s.ids.New(), OrderID: orderID, OrganizationID: organization, ProviderCode: provider}, actor)
	if err != nil && !errors.Is(err, ErrConflict) {
		return PaymentAttempt{}, ErrPaymentUnavailable
	}
	return value, err
}

func (s *Service) CreatePaymentAttempt(ctx context.Context, tenant, idempotency string, input PaymentAttempt) (PaymentAttempt, error) {
	return s.createPaymentAttempt(ctx, tenant, idempotency, input, "")
}
func (s *Service) CreatePaymentAttemptAs(ctx context.Context, tenant, idempotency string, input PaymentAttempt, actor string) (PaymentAttempt, error) {
	if actor == "" {
		return PaymentAttempt{}, fmt.Errorf("authenticated payment actor required")
	}
	return s.createPaymentAttempt(ctx, tenant, idempotency, input, actor)
}
func (s *Service) createPaymentAttempt(ctx context.Context, tenant, idempotency string, input PaymentAttempt, actor string) (PaymentAttempt, error) {
	if tenant == "" || len(idempotency) < 16 || len(idempotency) > 128 || input.OrderID == "" || input.OrganizationID == "" || input.ProviderCode == "" || !currencyPattern.MatchString(input.Currency) || input.AmountMinorUnits <= 0 {
		return PaymentAttempt{}, fmt.Errorf("invalid payment attempt")
	}
	recorder, ok := s.repository.(PaymentIntentRecorder)
	if !ok {
		return PaymentAttempt{}, fmt.Errorf("durable payment recorder unavailable")
	}
	input.ID = s.ids.New()
	input.State = "created"
	input.Version = 1
	return recorder.RecordPaymentIntent(ctx, tenant, s.ids.New(), idempotency, input, actor)
}
func (s *Service) TransitionPayment(ctx context.Context, tenant, organization, id, current, target string, version int64, providerReference string) error {
	if organization == "" || !paymentTransitions[current][target] || version < 1 {
		return fmt.Errorf("%w: invalid payment transition", ErrConflict)
	}
	return s.repository.TransitionPayment(ctx, tenant, organization, id, current, target, version, providerReference, s.ids.New())
}
