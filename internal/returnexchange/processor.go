package returnexchange

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	"elite.local/enterprise/internal/returneffects"
)

var (
	ErrNoWork           = errors.New("no return exchange work")
	ErrConflict         = errors.New("return exchange contract conflict")
	ErrStockUnavailable = errors.New("replacement stock unavailable")
)

var workerIDPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:-]{0,63}$`)

type Result struct {
	RequestID             string
	ReplacementOrderID    string
	ReplacementStockID    string
	ReplacementHandoverID string
	ResultSHA256          string
}

type Store interface {
	Claim(context.Context, string, string, string, time.Duration) (*returneffects.Work, error)
	PrepareExchange(context.Context, returneffects.Work, string) (Result, error)
	Finish(context.Context, returneffects.Work, string, returneffects.Completion) error
}

type IDGenerator interface{ New() string }

type Processor struct {
	store      Store
	ids        IDGenerator
	workerID   string
	lease      time.Duration
	retryDelay time.Duration
}

func NewProcessor(store Store, ids IDGenerator, workerID string, lease, retryDelay time.Duration) (*Processor, error) {
	if store == nil || ids == nil || !workerIDPattern.MatchString(workerID) || lease < time.Second || lease > 15*time.Minute || retryDelay < time.Second || retryDelay > time.Hour {
		return nil, fmt.Errorf("invalid return exchange processor configuration")
	}
	return &Processor{store: store, ids: ids, workerID: workerID, lease: lease, retryDelay: retryDelay}, nil
}

func (p *Processor) ProcessOne(ctx context.Context) (Result, error) {
	work, err := p.store.Claim(ctx, "fulfillment", p.workerID, p.ids.New(), p.lease)
	if err != nil {
		return Result{}, err
	}
	if work == nil {
		return Result{}, ErrNoWork
	}
	result, executeErr := p.store.PrepareExchange(ctx, *work, p.workerID)
	if executeErr == nil {
		return result, nil
	}
	completion := returneffects.Completion{Outcome: "retry", ErrorCode: "EXCHANGE_TRANSIENT_FAILURE", RetryAfter: p.retryDelay}
	if errors.Is(executeErr, ErrStockUnavailable) {
		completion.ErrorCode = "REPLACEMENT_STOCK_UNAVAILABLE"
	}
	if errors.Is(executeErr, ErrConflict) {
		completion = returneffects.Completion{Outcome: "blocked", ErrorCode: "EXCHANGE_CONTRACT_CONFLICT"}
	}
	if finishErr := p.store.Finish(ctx, *work, p.workerID, completion); finishErr != nil {
		return Result{}, errors.Join(executeErr, finishErr)
	}
	return Result{}, executeErr
}
