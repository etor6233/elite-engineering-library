package workers

import (
	"context"
	db "elite.local/enterprise/internal/platform/postgres"
	"fmt"
	"time"
)

type OutboxStore interface {
	Claim(context.Context, string, time.Duration, int) ([]db.OutboxEvent, error)
	RemainingLease(context.Context, string, string, string, int) (time.Duration, error)
	MarkPublished(context.Context, string, string, string, int) error
	Release(context.Context, string, string, string, string, int, time.Duration) error
}
type Publisher interface {
	Publish(context.Context, db.OutboxEvent, string) error
}
type OutboxProcessor struct {
	Store             OutboxStore
	Publisher         Publisher
	WorkerID          string
	Lease, RetryDelay time.Duration
	BatchSize         int
}

func (p OutboxProcessor) ProcessOnce(ctx context.Context) (int, error) {
	if p.Store == nil || p.Publisher == nil || p.WorkerID == "" || p.Lease < time.Microsecond || p.RetryDelay < 0 || p.BatchSize < 1 || p.BatchSize > 1000 {
		return 0, fmt.Errorf("invalid outbox processor")
	}
	events, err := p.Store.Claim(ctx, p.WorkerID, p.Lease, p.BatchSize)
	if err != nil {
		return 0, err
	}
	published := 0
	for _, event := range events {
		started := time.Now()
		remaining, err := p.Store.RemainingLease(ctx, event.TenantID, event.EventID, p.WorkerID, event.Attempts)
		remaining -= time.Since(started)
		if err != nil {
			return published, err
		}
		if remaining <= 0 || ctx.Err() != nil {
			return published, db.ErrOutboxClaimLost
		}
		publishCtx, cancel := context.WithTimeout(ctx, remaining)
		publishErr := p.Publisher.Publish(publishCtx, event, event.EventID)
		cancel()
		if publishErr != nil {
			if releaseErr := p.Store.Release(ctx, event.TenantID, event.EventID, p.WorkerID, "PUBLISH_FAILED", event.Attempts, p.RetryDelay); releaseErr != nil {
				return published, releaseErr
			}
			continue
		}
		if err := p.Store.MarkPublished(ctx, event.TenantID, event.EventID, p.WorkerID, event.Attempts); err != nil {
			return published, err
		}
		published++
	}
	return published, nil
}

type JobStore interface {
	Claim(context.Context, string, string, time.Duration, int) ([]db.Job, error)
	RemainingLease(context.Context, string, string, string, int) (time.Duration, error)
	Complete(context.Context, string, string, string, int) error
	Fail(context.Context, string, string, string, string, int, time.Duration) (bool, error)
}
type JobHandler interface {
	Handle(context.Context, db.Job) error
}
type JobProcessor struct {
	Store             JobStore
	Handler           JobHandler
	Queue, WorkerID   string
	Lease, RetryDelay time.Duration
	BatchSize         int
}

func (p JobProcessor) ProcessOnce(ctx context.Context) (int, error) {
	if p.Store == nil || p.Handler == nil || p.Queue == "" || p.WorkerID == "" {
		return 0, fmt.Errorf("invalid job processor")
	}
	jobs, err := p.Store.Claim(ctx, p.Queue, p.WorkerID, p.Lease, p.BatchSize)
	if err != nil {
		return 0, err
	}
	completed := 0
	for _, job := range jobs {
		// Subtract the complete DB round trip conservatively. A delayed response
		// must not extend a lease, and later batch entries may have expired.
		started := time.Now()
		remaining, err := p.Store.RemainingLease(ctx, job.TenantID, job.JobID, p.WorkerID, job.Attempts)
		remaining -= time.Since(started)
		if err != nil {
			return completed, err
		}
		if remaining <= 0 || ctx.Err() != nil {
			return completed, db.ErrJobClaimLost
		}
		jobCtx, cancel := context.WithTimeout(ctx, remaining)
		handleErr := p.Handler.Handle(jobCtx, job)
		cancel()
		if handleErr != nil {
			if _, failErr := p.Store.Fail(ctx, job.TenantID, job.JobID, p.WorkerID, "HANDLER_FAILED", job.Attempts, p.RetryDelay); failErr != nil {
				return completed, failErr
			}
			continue
		}
		if err := p.Store.Complete(ctx, job.TenantID, job.JobID, p.WorkerID, job.Attempts); err != nil {
			return completed, err
		}
		completed++
	}
	return completed, nil
}
