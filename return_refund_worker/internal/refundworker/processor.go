package refundworker

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

type Processor struct {
	store       Store
	providers   map[string]Provider
	workerID    string
	lease       time.Duration
	retryAfter  time.Duration
	maxAttempts int
}

func NewProcessor(store Store, providers map[string]Provider, workerID string, lease, retryAfter time.Duration, maxAttempts int) (*Processor, error) {
	if store == nil || strings.TrimSpace(workerID) == "" || lease <= 0 || retryAfter <= 0 || maxAttempts < 1 {
		return nil, fmt.Errorf("invalid refund processor configuration")
	}
	copyProviders := make(map[string]Provider, len(providers))
	for name, provider := range providers {
		if provider != nil {
			copyProviders[name] = provider
		}
	}
	return &Processor{store: store, providers: copyProviders, workerID: workerID, lease: lease, retryAfter: retryAfter, maxAttempts: maxAttempts}, nil
}

func classifyStatus(provider, status string) (outcome, code string) {
	status = strings.ToLower(strings.TrimSpace(status))
	switch provider {
	case "stripe":
		switch status {
		case "succeeded":
			return "succeeded", ""
		case "pending":
			return "retry", "PROVIDER_PENDING"
		case "requires_action":
			return "retry", "PROVIDER_ACTION_REQUIRED"
		case "failed":
			return "failed", "PROVIDER_REFUND_FAILED"
		case "canceled":
			return "failed", "PROVIDER_REFUND_CANCELED"
		}
	case "mercado_pago":
		switch status {
		case "approved":
			return "succeeded", ""
		case "pending", "in_process":
			return "retry", "PROVIDER_PENDING"
		case "rejected", "failed":
			return "failed", "PROVIDER_REFUND_FAILED"
		case "cancelled", "canceled":
			return "failed", "PROVIDER_REFUND_CANCELED"
		}
	}
	return "blocked", "UNKNOWN_PROVIDER_STATUS"
}

func validateResult(refund Refund, result ProviderResult) error {
	if result.ProviderPaymentReference != refund.ProviderPaymentReference || result.ProviderRefundReference == "" || result.ProviderStatus == "" || result.Currency != refund.Currency || result.AmountMinorUnits != refund.AmountMinorUnits {
		return ErrResponseMismatch
	}
	if refund.ProviderRefundReference != "" && result.ProviderRefundReference != refund.ProviderRefundReference {
		return ErrResponseMismatch
	}
	return nil
}

func (p *Processor) Step(ctx context.Context, claimToken string) error {
	if p == nil || ctx == nil || strings.TrimSpace(claimToken) == "" {
		return fmt.Errorf("invalid refund step")
	}
	work, err := p.store.Claim(ctx, p.workerID, claimToken, p.lease)
	if err != nil {
		return err
	}
	if work == nil {
		return ErrNoWork
	}
	refund, err := p.store.Prepare(ctx, *work, p.workerID)
	if err != nil {
		outcome, code, retry := "blocked", "REFUND_PREPARATION_FAILED", time.Duration(0)
		if errors.Is(err, ErrMappingConflict) {
			code = "REFUND_MAPPING_CONFLICT"
		} else if retryablePreparationError(err) {
			code = "REFUND_PREPARATION_EXHAUSTED"
			if work.Attempt < p.maxAttempts {
				outcome, code, retry = "retry", "REFUND_PREPARATION_TRANSIENT", p.retryAfter
			}
		}
		if finishErr := p.store.Finish(ctx, *work, p.workerID, outcome, code, retry); finishErr != nil {
			return errors.Join(err, finishErr)
		}
		return err
	}
	provider := p.providers[refund.Provider]
	if provider == nil {
		return p.store.Finish(ctx, *work, p.workerID, "blocked", "PROVIDER_CONFIG_MISSING", 0)
	}
	var result ProviderResult
	if refund.ProviderRefundReference == "" {
		result, err = provider.Create(ctx, refund)
	} else {
		result, err = provider.Retrieve(ctx, refund)
	}
	if err != nil {
		outcome := "retry"
		if work.Attempt >= p.maxAttempts || errors.Is(err, ErrInvalidProvider) || errors.Is(err, ErrResponseMismatch) {
			outcome = "blocked"
		}
		if finishErr := p.store.Finish(ctx, *work, p.workerID, outcome, "PROVIDER_CALL_FAILED", p.retryAfter); finishErr != nil {
			return errors.Join(err, finishErr)
		}
		return err
	}
	if err = validateResult(refund, result); err != nil {
		if finishErr := p.store.Complete(ctx, *work, p.workerID, refund, result, "blocked", "PROVIDER_RESPONSE_MISMATCH", 0); finishErr != nil {
			return errors.Join(err, finishErr)
		}
		return err
	}
	outcome, code := classifyStatus(refund.Provider, result.ProviderStatus)
	retry := time.Duration(0)
	if outcome == "retry" {
		retry = p.retryAfter
		if work.Attempt >= p.maxAttempts {
			outcome = "blocked"
			code = "PROVIDER_RECONCILIATION_EXHAUSTED"
			retry = 0
		}
	}
	return p.store.Complete(ctx, *work, p.workerID, refund, result, outcome, code, retry)
}

// Preparation has not called the provider. Retry the complete local transaction
// only for recognized transient failures, preserving its request/idempotency key.
// Unknown errors, source ambiguity and integrity violations require inspection.
func retryablePreparationError(err error) bool {
	var sqlErr *pgconn.PgError
	if errors.As(err, &sqlErr) {
		switch sqlErr.Code {
		case "40001", "40P01", "55P03", "57P01", "57P02", "57P03":
			return true
		}
		return strings.HasPrefix(sqlErr.Code, "08")
	}
	return pgconn.SafeToRetry(err) || errors.Is(err, context.DeadlineExceeded)
}
