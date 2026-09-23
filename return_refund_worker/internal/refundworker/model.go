package refundworker

import (
	"context"
	"errors"
	"time"
)

var (
	ErrNoWork           = errors.New("no refund work")
	ErrMappingConflict  = errors.New("refund source mapping conflict")
	ErrStaleClaim       = errors.New("stale refund claim")
	ErrInvalidProvider  = errors.New("invalid refund provider request")
	ErrResponseMismatch = errors.New("refund provider response mismatch")
)

type Work struct {
	TenantID       string
	RequestID      string
	IdempotencyKey string
	Attempt        int
	ClaimToken     string
}

type Refund struct {
	TenantID                 string
	RequestID                string
	PaymentAttemptID         string
	OrderID                  string
	LineID                   string
	StockUnitID              string
	Provider                 string
	ProviderPaymentReference string
	IdempotencyKey           string
	Currency                 string
	AmountMinorUnits         int64
	ProviderRefundReference  string
	ProviderStatus           string
}

type ProviderResult struct {
	ProviderPaymentReference string
	ProviderRefundReference  string
	ProviderStatus           string
	Currency                 string
	AmountMinorUnits         int64
}

type Provider interface {
	Create(context.Context, Refund) (ProviderResult, error)
	Retrieve(context.Context, Refund) (ProviderResult, error)
}

type Store interface {
	Claim(context.Context, string, string, time.Duration) (*Work, error)
	Prepare(context.Context, Work, string) (Refund, error)
	Complete(context.Context, Work, string, Refund, ProviderResult, string, string, time.Duration) error
	Finish(context.Context, Work, string, string, string, time.Duration) error
}
