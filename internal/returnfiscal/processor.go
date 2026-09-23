package returnfiscal

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	"elite.local/enterprise/internal/returneffects"
)

var (
	ErrNoWork   = errors.New("no return fiscal work")
	ErrPending  = errors.New("credit note authorization pending")
	ErrRejected = errors.New("credit note rejected")
	ErrConflict = errors.New("return fiscal source conflict")
)

var workerPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:-]{0,63}$`)

type Result struct {
	RequestID, OriginalInvoiceID, CreditInvoiceID, Status, ResultSHA256 string
}

type Store interface {
	Claim(context.Context, string, string, string, time.Duration) (*returneffects.Work, error)
	RequestOrObserve(context.Context, returneffects.Work, string) (Result, error)
	Finish(context.Context, returneffects.Work, string, returneffects.Completion) error
}

type IDGenerator interface{ New() string }

type Processor struct {
	store        Store
	ids          IDGenerator
	workerID     string
	lease, retry time.Duration
}

func NewProcessor(store Store, ids IDGenerator, workerID string, lease, retry time.Duration) (*Processor, error) {
	if store == nil || ids == nil || !workerPattern.MatchString(workerID) || lease < time.Second || lease > 15*time.Minute || retry < time.Second || retry > time.Hour {
		return nil, fmt.Errorf("invalid return fiscal processor configuration")
	}
	return &Processor{store: store, ids: ids, workerID: workerID, lease: lease, retry: retry}, nil
}

func (p *Processor) ProcessOne(ctx context.Context) (Result, error) {
	work, err := p.store.Claim(ctx, "fiscal", p.workerID, p.ids.New(), p.lease)
	if err != nil {
		return Result{}, err
	}
	if work == nil {
		return Result{}, ErrNoWork
	}
	result, observeErr := p.store.RequestOrObserve(ctx, *work, p.workerID)
	completion := returneffects.Completion{Outcome: "succeeded", ProviderReference: result.CreditInvoiceID, ResultSHA256: result.ResultSHA256}
	if observeErr != nil {
		completion = returneffects.Completion{Outcome: "retry", ErrorCode: "FISCAL_TRANSIENT_FAILURE", RetryAfter: p.retry}
		if errors.Is(observeErr, ErrPending) {
			completion.ErrorCode = "FISCAL_AUTHORIZATION_PENDING"
		}
		if errors.Is(observeErr, ErrRejected) {
			completion = returneffects.Completion{Outcome: "blocked", ErrorCode: "FISCAL_AUTHORIZATION_REJECTED", ProviderReference: result.CreditInvoiceID}
		}
		if errors.Is(observeErr, ErrConflict) {
			completion = returneffects.Completion{Outcome: "blocked", ErrorCode: "FISCAL_SOURCE_CONFLICT"}
		}
	}
	if finishErr := p.store.Finish(ctx, *work, p.workerID, completion); finishErr != nil {
		return Result{}, errors.Join(observeErr, finishErr)
	}
	if observeErr != nil {
		return Result{}, observeErr
	}
	return result, nil
}
