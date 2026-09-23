package whatsappbridge

// AUTHORED opt-in worker. Reuses the existing job queue, router, observer and
// append-only audit. No daemon is registered, no send or customer reply occurs.
import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/postgres"
	"github.com/jackc/pgx/v5"
)

var ErrStatusWork = errors.New("whatsappbridge: status job unresolved; inspect scoped job and audit")

type StatusWorkResult struct {
	Claimed, Completed, FailureRecorded, Terminal bool
	InsertedObservations                          int
}

type StatusWorker struct {
	router       *StatusRouter
	jobs         *postgres.Jobs
	worker       string
	lease, retry time.Duration
	busy         chan struct{}
}

func NewStatusWorker(router *StatusRouter, worker string, lease, retry time.Duration) (*StatusWorker, error) {
	if router == nil || router.observer.Approvals == nil || router.observer.Approvals.pool == nil || router.busy == nil || !webhookConnection.MatchString(worker) || lease < time.Second || lease > 2*time.Minute || retry < 0 || retry > time.Hour {
		return nil, ErrStatusWork
	}
	return &StatusWorker{router: router, jobs: postgres.NewJobs(router.observer.Approvals.pool), worker: worker, lease: lease, retry: retry, busy: make(chan struct{}, 1)}, nil
}

// ProcessOnce claims at most one job, in the configured tenant/connection/type.
// p must come from actual authentication; no identity is derived from payload.
// A caller must handle errors even when observations were partially committed.
func (w *StatusWorker) ProcessOnce(ctx context.Context, p identity.Principal) (StatusWorkResult, error) {
	var result StatusWorkResult
	if w == nil || p.TenantID != w.router.observer.TenantID || p.Subject == "" || len(p.Subject) > 255 || !p.Allowed("appointment:manage") || ctx.Err() != nil {
		return result, ErrStatusWork
	}
	select {
	case w.busy <- struct{}{}:
		defer func() { <-w.busy }()
	default:
		return result, ErrStatusWork
	}
	ctx, cancelCall := context.WithTimeout(ctx, w.lease+5*time.Second)
	defer cancelCall()
	var connectionOrg *string
	err := w.router.observer.Approvals.pool.QueryRow(ctx, `select organization_id from integration.provider_connection where tenant_id=$1 and connection_id=$2 and provider_code='meta-whatsapp' and state='active'`, p.TenantID, w.router.connection).Scan(&connectionOrg)
	if err != nil || (connectionOrg != nil && !p.AllowedOrganization(*connectionOrg)) {
		return result, ErrStatusWork
	}
	match, _ := json.Marshal(map[string]string{"connection_id": w.router.connection, "provider_code": "meta-whatsapp", "event_type": "whatsapp.raw_webhook.received.v1"})
	if expired, err := w.quarantineExpired(ctx, p, match); err != nil {
		return result, ErrStatusWork
	} else if expired {
		result.Terminal = true
		result.FailureRecorded = true
		return result, ErrStatusWork
	}
	jobs, err := w.jobs.ClaimScoped(ctx, "provider-events", w.worker, w.lease, 1, postgres.JobScope{TenantID: p.TenantID, JobType: "provider.webhook.received", SchemaVersion: 1, PayloadMatch: match})
	if err != nil {
		return result, ErrStatusWork
	}
	if len(jobs) == 0 {
		return result, nil
	}
	result.Claimed = true
	job := jobs[0]
	event, err := w.eventID(job)
	if err != nil {
		// No trustworthy inbox identity: retain job/audit, do not guess an event.
		return w.recordFailure(ctx, p, job, "", "JOB_CONTRACT_INVALID", result)
	}
	started := time.Now()
	remaining, err := w.jobs.RemainingLease(ctx, job.TenantID, job.JobID, w.worker, job.Attempts)
	remaining -= time.Since(started)
	if err != nil || remaining <= 0 {
		return result, ErrStatusWork
	}
	workCtx, cancel := context.WithTimeout(ctx, remaining)
	result.InsertedObservations, err = w.router.observeRetained(workCtx, p, event, func(ctx context.Context, tx pgx.Tx) error {
		if err := postgres.CompleteJobInTx(ctx, tx, job, w.worker); err != nil {
			return err
		}
		tag, err := tx.Exec(ctx, `update integration.webhook_event set state='processed',processed_at=coalesce(processed_at,clock_timestamp()),last_error_code=null
 where tenant_id=$1 and connection_id=$2 and provider_event_id=$3 and provider_code='meta-whatsapp'
 and event_type='whatsapp.raw_webhook.received.v1' and state in ('received','processed')`, p.TenantID, w.router.connection, event)
		if err != nil {
			return err
		}
		if tag.RowsAffected() != 1 {
			return ErrStatusWork
		}
		return statusWorkAudit(ctx, tx, p, job, event, "completed", "STATUS_OBSERVED", false)
	})
	cancel()
	if err != nil {
		return w.recordFailure(ctx, p, job, event, "STATUS_UNRESOLVED", result)
	}
	result.Completed = true
	return result, nil
}

func (w *StatusWorker) quarantineExpired(ctx context.Context, p identity.Principal, match json.RawMessage) (bool, error) {
	tx, err := w.router.observer.Approvals.pool.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer tx.Rollback(ctx)
	var job postgres.Job
	err = tx.QueryRow(ctx, `select tenant_id,job_id,queue,job_type,schema_version,payload,attempts,max_attempts from platform.job
 where tenant_id=$1 and queue='provider-events' and job_type='provider.webhook.received' and schema_version=1 and payload @> $2::jsonb
 and completed_at is null and terminal_error_code is null and attempts>=max_attempts and claimed_until<clock_timestamp()
 order by claimed_until,job_id for update skip locked limit 1`, p.TenantID, match).Scan(&job.TenantID, &job.JobID, &job.Queue, &job.JobType, &job.SchemaVersion, &job.Payload, &job.Attempts, &job.MaxAttempts)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if err = postgres.ExhaustJobInTx(ctx, tx, job); err != nil {
		return false, err
	}
	event, _ := w.eventID(job)
	if event != "" {
		if _, err = tx.Exec(ctx, `update integration.webhook_event set state='failed',last_error_code='LEASE_EXHAUSTED' where tenant_id=$1 and connection_id=$2 and provider_event_id=$3 and provider_code='meta-whatsapp' and event_type='whatsapp.raw_webhook.received.v1' and state='received'`, p.TenantID, w.router.connection, event); err != nil {
			return false, err
		}
	}
	if err = statusWorkAudit(ctx, tx, p, job, event, "failed", "LEASE_EXHAUSTED", true); err != nil {
		return false, err
	}
	if err = tx.Commit(ctx); err != nil {
		return false, err
	}
	return true, nil
}

func (w *StatusWorker) eventID(job postgres.Job) (string, error) {
	var fields map[string]string
	if job.TenantID != w.router.observer.TenantID || job.Queue != "provider-events" || job.JobType != "provider.webhook.received" || job.SchemaVersion != 1 || len(job.Payload) > 4096 || json.Unmarshal(job.Payload, &fields) != nil || len(fields) != 4 || fields["connection_id"] != w.router.connection || fields["provider_code"] != "meta-whatsapp" || fields["event_type"] != "whatsapp.raw_webhook.received.v1" {
		return "", ErrStatusWork
	}
	event := fields["provider_event_id"]
	if !strings.HasPrefix(event, "wa:") || !validDigest(strings.TrimPrefix(event, "wa:")) {
		return "", ErrStatusWork
	}
	return event, nil
}

func (w *StatusWorker) recordFailure(ctx context.Context, p identity.Principal, job postgres.Job, event, code string, result StatusWorkResult) (StatusWorkResult, error) {
	// Bounded metadata cleanup after cancellation is not another processing try.
	cleanup, cancel := context.WithTimeout(context.WithoutCancel(ctx), 3*time.Second)
	defer cancel()
	tx, err := w.router.observer.Approvals.pool.Begin(cleanup)
	if err != nil {
		return result, ErrStatusWork
	}
	defer tx.Rollback(cleanup)
	terminal, err := postgres.FailJobInTx(cleanup, tx, job, w.worker, code, w.retry)
	if err != nil {
		return result, ErrStatusWork
	}
	if event != "" {
		// Missing/failed/processed inbox cannot be rewritten as a retryable one.
		// Job failure and audit remain useful even if its inbox is absent.
		_, err = tx.Exec(cleanup, `update integration.webhook_event set state=case when $4 then 'failed' else 'received' end,last_error_code=$5
 where tenant_id=$1 and connection_id=$2 and provider_event_id=$3 and provider_code='meta-whatsapp'
 and event_type='whatsapp.raw_webhook.received.v1' and state='received'`, p.TenantID, w.router.connection, event, terminal, code)
		if err != nil {
			return result, ErrStatusWork
		}
	}
	if statusWorkAudit(cleanup, tx, p, job, event, "failed", code, terminal) != nil || tx.Commit(cleanup) != nil {
		return result, ErrStatusWork
	}
	result.FailureRecorded = true
	result.Terminal = terminal
	return result, ErrStatusWork
}

func statusWorkAudit(ctx context.Context, tx pgx.Tx, p identity.Principal, job postgres.Job, event, action, code string, terminal bool) error {
	_, err := tx.Exec(ctx, `insert into audit.event(tenant_id,event_id,actor_subject,action,resource_type,resource_id,decision,reason_code,evidence)
 values($1,gen_random_uuid(),$2,$3,'provider_job',$4,'system',$5,jsonb_build_object('attempt',$6::int,'terminal',$7::bool,'provider_event_sha256',$8::text))`, p.TenantID, p.Subject, "whatsapp.status_job."+action, job.JobID, code, job.Attempts, terminal, digest([]byte(event)))
	return err
}
