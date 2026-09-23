package royalty

import (
	"context"
	"errors"
	"regexp"
	"time"
)

var ErrConflict = errors.New("royalty conflict")
var ErrInvalid = errors.New("invalid royalty command")
var currencyPattern = regexp.MustCompile(`^[A-Z]{3}$`)

type Policy struct {
	ID              string     `json:"id"`
	AgreementID     string     `json:"agreement_id"`
	OrganizationID  string     `json:"organization_id"`
	Currency        string     `json:"currency"`
	RateBasisPoints int        `json:"rate_basis_points"`
	ValidFrom       time.Time  `json:"valid_from"`
	ValidUntil      *time.Time `json:"valid_until,omitempty"`
}

type PaymentEvent struct {
	ID                     string    `json:"id"`
	OrganizationID         string    `json:"organization_id"`
	PaymentAttemptID       string    `json:"payment_attempt_id"`
	PaymentExpectedVersion int64     `json:"payment_expected_version"`
	State                  string    `json:"state"`
	OccurredAt             time.Time `json:"occurred_at"`
}

type Accrual struct {
	ID                string    `json:"id"`
	PolicyID          string    `json:"policy_id"`
	AgreementID       string    `json:"agreement_id"`
	OrganizationID    string    `json:"organization_id"`
	PaymentAttemptID  string    `json:"payment_attempt_id"`
	SourceEventKey    string    `json:"source_event_key"`
	SourceState       string    `json:"source_state"`
	Currency          string    `json:"currency"`
	BasisMinorUnits   int64     `json:"basis_minor_units"`
	RoyaltyMinorUnits int64     `json:"royalty_minor_units"`
	OccurredAt        time.Time `json:"occurred_at"`
}

type Settlement struct {
	ID                 string    `json:"id"`
	OrganizationID     string    `json:"organization_id"`
	Currency           string    `json:"currency"`
	PeriodStart        time.Time `json:"period_start"`
	PeriodEnd          time.Time `json:"period_end"`
	Status             string    `json:"status"`
	ExpectedMinorUnits int64     `json:"expected_minor_units"`
	Version            int64     `json:"version"`
	ReversalOf         string    `json:"reversal_of,omitempty"`
}

type Reconciliation struct {
	ID                   string `json:"id"`
	SettlementID         string `json:"settlement_id"`
	ExternalReference    string `json:"external_reference"`
	ExpectedMinorUnits   int64  `json:"expected_minor_units"`
	ActualMinorUnits     int64  `json:"actual_minor_units"`
	DifferenceMinorUnits int64  `json:"difference_minor_units"`
	Status               string `json:"status"`
}

type Repository interface {
	CreatePolicy(context.Context, string, string, Policy) error
	AccruePayment(context.Context, string, string, PaymentEvent, string) (Accrual, error)
	OpenSettlement(context.Context, string, string, Settlement) error
	CloseSettlement(context.Context, string, string, string, string, int64) (Settlement, error)
	ReverseSettlement(context.Context, string, string, string, string, string, int64, string) (Settlement, error)
	Reconcile(context.Context, string, string, Reconciliation, string, string) (Reconciliation, error)
}

type IDGenerator interface{ New() string }

type Service struct {
	repository Repository
	ids        IDGenerator
}

func NewService(repository Repository, ids IDGenerator) *Service {
	return &Service{repository: repository, ids: ids}
}

func (s *Service) CreatePolicy(ctx context.Context, tenant string, value Policy) (Policy, error) {
	if tenant == "" || value.AgreementID == "" || value.OrganizationID == "" || !currencyPattern.MatchString(value.Currency) || value.RateBasisPoints < 1 || value.RateBasisPoints > 10000 || value.ValidFrom.IsZero() || (value.ValidUntil != nil && !value.ValidUntil.After(value.ValidFrom)) {
		return value, ErrInvalid
	}
	value.ID = s.ids.New()
	if err := s.repository.CreatePolicy(ctx, tenant, s.ids.New(), value); err != nil {
		return value, err
	}
	return value, nil
}

func (s *Service) AccruePayment(ctx context.Context, tenant, sourceEventKey string, value PaymentEvent) (Accrual, error) {
	if tenant == "" || sourceEventKey == "" || value.OrganizationID == "" || value.PaymentAttemptID == "" || value.PaymentExpectedVersion < 1 || (value.State != "captured" && value.State != "refunded") || value.OccurredAt.IsZero() {
		return Accrual{}, ErrInvalid
	}
	value.ID = s.ids.New()
	return s.repository.AccruePayment(ctx, tenant, sourceEventKey, value, s.ids.New())
}

func (s *Service) OpenSettlement(ctx context.Context, tenant string, value Settlement) (Settlement, error) {
	if tenant == "" || value.OrganizationID == "" || !currencyPattern.MatchString(value.Currency) || value.PeriodStart.IsZero() || !value.PeriodEnd.After(value.PeriodStart) {
		return value, ErrInvalid
	}
	value.ID = s.ids.New()
	value.Status = "draft"
	value.Version = 1
	if err := s.repository.OpenSettlement(ctx, tenant, s.ids.New(), value); err != nil {
		return value, err
	}
	return value, nil
}

func (s *Service) CloseSettlement(ctx context.Context, tenant, organization, id string, expectedVersion int64) (Settlement, error) {
	if tenant == "" || organization == "" || id == "" || expectedVersion < 1 {
		return Settlement{}, ErrInvalid
	}
	return s.repository.CloseSettlement(ctx, tenant, organization, id, s.ids.New(), expectedVersion)
}

func (s *Service) ReverseSettlement(ctx context.Context, tenant, organization, id string, expectedVersion int64, reason string) (Settlement, error) {
	if tenant == "" || organization == "" || id == "" || expectedVersion < 1 || len(reason) < 3 || len(reason) > 500 {
		return Settlement{}, ErrInvalid
	}
	return s.repository.ReverseSettlement(ctx, tenant, organization, id, s.ids.New(), s.ids.New(), expectedVersion, reason)
}

func (s *Service) Reconcile(ctx context.Context, tenant, organization, actor string, value Reconciliation) (Reconciliation, error) {
	if tenant == "" || organization == "" || actor == "" || value.SettlementID == "" || value.ExternalReference == "" || len(value.ExternalReference) > 250 {
		return value, ErrInvalid
	}
	value.ID = s.ids.New()
	return s.repository.Reconcile(ctx, tenant, organization, value, actor, s.ids.New())
}
