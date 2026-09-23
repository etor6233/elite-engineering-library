package returnaccounting

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	"elite.local/enterprise/internal/returneffects"
)

var (
	ErrNoWork   = errors.New("no return accounting work")
	ErrConflict = errors.New("return accounting contract conflict")
)
var workerPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:-]{0,63}$`)

type Result struct {
	RequestID, OriginalJournalID, ReversalJournalID, ResultSHA256 string
	ReversedMinorUnits                                            int64
}
type Store interface {
	Claim(context.Context, string, string, string, time.Duration) (*returneffects.Work, error)
	PostOrReconcile(context.Context, returneffects.Work, string) (Result, error)
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
		return nil, fmt.Errorf("invalid return accounting processor configuration")
	}
	return &Processor{store: store, ids: ids, workerID: workerID, lease: lease, retry: retry}, nil
}
func (p *Processor) ProcessOne(ctx context.Context) (Result, error) {
	work, err := p.store.Claim(ctx, "accounting", p.workerID, p.ids.New(), p.lease)
	if err != nil {
		return Result{}, err
	}
	if work == nil {
		return Result{}, ErrNoWork
	}
	result, postErr := p.store.PostOrReconcile(ctx, *work, p.workerID)
	if postErr == nil {
		return result, nil
	}
	completion := returneffects.Completion{Outcome: "retry", ErrorCode: "ACCOUNTING_TRANSIENT_FAILURE", RetryAfter: p.retry}
	if errors.Is(postErr, ErrConflict) {
		completion = returneffects.Completion{Outcome: "blocked", ErrorCode: "ACCOUNTING_SOURCE_CONFLICT"}
	}
	if finishErr := p.store.Finish(ctx, *work, p.workerID, completion); finishErr != nil {
		return Result{}, errors.Join(postErr, finishErr)
	}
	return Result{}, postErr
}
