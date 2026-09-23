package returneffects

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"
)

var (
	ErrNoWork         = errors.New("no return effect work")
	ErrInventoryState = errors.New("return inventory state conflict")
)

var workerIDPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:-]{0,63}$`)

type Work struct {
	TenantID        string
	RequestID       string
	DispositionID   string
	EffectKind      string
	OwnerContext    string
	IdempotencyKey  string
	Attempt         int
	ClaimToken      string
	InventoryAction string
	OrganizationID  string
	StockUnitID     string
}

type Result struct {
	RequestID      string
	StockUnitID    string
	InventoryState string
	StockVersion   int64
	ResultSHA256   string
}

type Completion struct {
	Outcome           string
	ErrorCode         string
	ProviderReference string
	ResultSHA256      string
	RetryAfter        time.Duration
}

type Store interface {
	Claim(context.Context, string, string, string, time.Duration) (*Work, error)
	ApplyInventory(context.Context, Work, string) (Result, error)
	Finish(context.Context, Work, string, Completion) error
	ResumeBlocked(context.Context, string, string, string, string, string) error
}

type IDGenerator interface{ New() string }

type Processor struct {
	store      Store
	ids        IDGenerator
	workerID   string
	lease      time.Duration
	retryDelay time.Duration
}

func NewInventoryProcessor(store Store, ids IDGenerator, workerID string, lease, retryDelay time.Duration) (*Processor, error) {
	if store == nil || ids == nil || !workerIDPattern.MatchString(workerID) || lease < time.Second || lease > 15*time.Minute || retryDelay < time.Second || retryDelay > time.Hour {
		return nil, fmt.Errorf("invalid return effect processor configuration")
	}
	return &Processor{store: store, ids: ids, workerID: workerID, lease: lease, retryDelay: retryDelay}, nil
}

func (p *Processor) ProcessOne(ctx context.Context) (Result, error) {
	claimToken := p.ids.New()
	work, err := p.store.Claim(ctx, "inventory", p.workerID, claimToken, p.lease)
	if err != nil {
		return Result{}, err
	}
	if work == nil {
		return Result{}, ErrNoWork
	}
	result, applyErr := p.store.ApplyInventory(ctx, *work, p.workerID)
	if applyErr == nil {
		return result, nil
	}
	completion := Completion{Outcome: "retry", ErrorCode: "INVENTORY_TRANSIENT_FAILURE", RetryAfter: p.retryDelay}
	if errors.Is(applyErr, ErrInventoryState) {
		completion = Completion{Outcome: "blocked", ErrorCode: "INVENTORY_STATE_CONFLICT"}
	}
	if finishErr := p.store.Finish(ctx, *work, p.workerID, completion); finishErr != nil {
		return Result{}, errors.Join(applyErr, finishErr)
	}
	return Result{}, applyErr
}
