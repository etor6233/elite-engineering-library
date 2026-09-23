package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Job struct {
	TenantID, JobID, Queue, JobType      string
	SchemaVersion, Attempts, MaxAttempts int
	Payload                              json.RawMessage
}
type Jobs struct{ pool *pgxpool.Pool }

var ErrJobClaimLost = errors.New("job claim lost")

// ExhaustJobInTx preserves a crashed final attempt as terminal evidence. It
// cannot expire a live claim, reset attempts, or silently start another effect.
func ExhaustJobInTx(ctx context.Context, tx pgx.Tx, job Job) error {
	if tx == nil || job.Attempts < 1 || !json.Valid(job.Payload) {
		return ErrJobClaimLost
	}
	tag, err := tx.Exec(ctx, lockedJob+`update platform.job j set claimed_by=null,claimed_until=null,terminal_error_code='LEASE_EXHAUSTED'
 from owned o where j.tenant_id=o.tenant_id and j.job_id=o.job_id and o.attempts=$3
 and o.attempts>=o.max_attempts and o.claimed_until<clock_timestamp()
 and o.queue=$4 and o.job_type=$5 and o.schema_version=$6 and o.payload=$7::jsonb`, job.TenantID, job.JobID, job.Attempts, job.Queue, job.JobType, job.SchemaVersion, job.Payload)
	if err != nil {
		return err
	}
	if tag.RowsAffected() != 1 {
		return ErrJobClaimLost
	}
	return nil
}

func NewJobs(pool *pgxpool.Pool) *Jobs { return &Jobs{pool: pool} }

// JobScope is server-side selection, not authentication or tenant authorization.
type JobScope struct {
	TenantID, JobType string
	SchemaVersion     int
	PayloadMatch      json.RawMessage
}

func (j *Jobs) ClaimScoped(ctx context.Context, queue, worker string, lease time.Duration, limit int, scope JobScope) ([]Job, error) {
	var fields map[string]json.RawMessage
	if scope.TenantID == "" || scope.JobType == "" || scope.SchemaVersion < 1 || len(scope.PayloadMatch) > 4096 || json.Unmarshal(scope.PayloadMatch, &fields) != nil || fields == nil {
		return nil, fmt.Errorf("invalid job scope")
	}
	return j.claim(ctx, queue, worker, lease, limit, &scope)
}

func (j *Jobs) Claim(ctx context.Context, queue, worker string, lease time.Duration, limit int) ([]Job, error) {
	return j.claim(ctx, queue, worker, lease, limit, nil)
}

func (j *Jobs) claim(ctx context.Context, queue, worker string, lease time.Duration, limit int, scope *JobScope) ([]Job, error) {
	if j == nil || j.pool == nil || queue == "" || worker == "" || lease < time.Microsecond || limit < 1 || limit > 1000 {
		return nil, fmt.Errorf("invalid job claim parameters")
	}
	var tenant any
	jobType, schema, match := "", 0, json.RawMessage(`{}`)
	if scope != nil {
		tenant = scope.TenantID
		jobType = scope.JobType
		schema = scope.SchemaVersion
		match = scope.PayloadMatch
	}
	rows, err := j.pool.Query(ctx, `with candidates as (
select tenant_id,job_id from platform.job where queue=$1 and completed_at is null and terminal_error_code is null and available_at<=clock_timestamp() and attempts<max_attempts and (claimed_until is null or claimed_until<clock_timestamp())
and ($5::uuid is null or tenant_id=$5) and ($6::text='' or job_type=$6) and ($7::integer=0 or schema_version=$7) and payload @> $8::jsonb
order by priority desc,available_at,job_id for update skip locked limit $2)
update platform.job j set claimed_by=$3,claimed_until=clock_timestamp()+$4::interval,attempts=attempts+1 from candidates c where j.tenant_id=c.tenant_id and j.job_id=c.job_id returning j.tenant_id,j.job_id,j.queue,j.job_type,j.schema_version,j.payload,j.attempts,j.max_attempts`, queue, limit, worker, lease.String(), tenant, jobType, schema, match)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	jobs := make([]Job, 0, limit)
	for rows.Next() {
		var job Job
		if err := rows.Scan(&job.TenantID, &job.JobID, &job.Queue, &job.JobType, &job.SchemaVersion, &job.Payload, &job.Attempts, &job.MaxAttempts); err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}
	return jobs, rows.Err()
}

// Attempts is a monotonically increasing claim generation, not an operator-
// resettable counter. A worker name alone is not a fencing token.
// Lock first, then check the DB clock: waiting for a row lock may expire a lease.
const lockedJob = `with owned as materialized (
 select tenant_id,job_id,claimed_by,claimed_until,attempts,max_attempts,queue,job_type,schema_version,payload
 from platform.job where tenant_id=$1 and job_id=$2
 and completed_at is null and terminal_error_code is null for update
) `

func (j *Jobs) Complete(ctx context.Context, tenant, jobID, worker string, attempt int) error {
	if j == nil || j.pool == nil || worker == "" || attempt < 1 {
		return ErrJobClaimLost
	}
	return completeJob(ctx, j.pool, tenant, jobID, worker, attempt, nil)
}

type jobExecutor interface {
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	QueryRow(context.Context, string, ...any) pgx.Row
}

// CompleteJobInTx does not commit. The caller atomically records its durable
// receipt/ack in the same transaction. Full claimed metadata must still match.
func CompleteJobInTx(ctx context.Context, tx pgx.Tx, job Job, worker string) error {
	if tx == nil || worker == "" || job.Attempts < 1 || job.Queue == "" || job.JobType == "" || job.SchemaVersion < 1 || len(job.Payload) > 65536 || !json.Valid(job.Payload) {
		return ErrJobClaimLost
	}
	return completeJob(ctx, tx, job.TenantID, job.JobID, worker, job.Attempts, &job)
}

func completeJob(ctx context.Context, db jobExecutor, tenant, jobID, worker string, attempt int, expected *Job) error {
	var payload any
	queue, kind, schema := "", "", 0
	if expected != nil {
		payload = expected.Payload
		queue = expected.Queue
		kind = expected.JobType
		schema = expected.SchemaVersion
	}
	r, err := db.Exec(ctx, lockedJob+`update platform.job j set completed_at=clock_timestamp(),claimed_by=null,claimed_until=null
 from owned o where j.tenant_id=o.tenant_id and j.job_id=o.job_id
 and o.claimed_by=$3 and o.attempts=$4 and o.claimed_until>clock_timestamp()
 and ($5::jsonb is null or (o.payload=$5 and o.queue=$6 and o.job_type=$7 and o.schema_version=$8))`, tenant, jobID, worker, attempt, payload, queue, kind, schema)
	if err != nil {
		return err
	}
	if r.RowsAffected() != 1 {
		return ErrJobClaimLost
	}
	return nil
}

func (j *Jobs) Fail(ctx context.Context, tenant, jobID, worker, code string, attempt int, retryAfter time.Duration) (bool, error) {
	if j == nil || j.pool == nil || worker == "" || attempt < 1 || code == "" || retryAfter < 0 {
		return false, fmt.Errorf("invalid job failure parameters")
	}
	return failJob(ctx, j.pool, tenant, jobID, worker, code, attempt, retryAfter, nil)
}

// FailJobInTx shares the transition implementation with Jobs.Fail. The caller
// can append a safe failure audit and update its inbox in the same transaction.
func FailJobInTx(ctx context.Context, tx pgx.Tx, job Job, worker, code string, retryAfter time.Duration) (bool, error) {
	if tx == nil || worker == "" || job.Attempts < 1 || job.Queue == "" || job.JobType == "" || job.SchemaVersion < 1 || code == "" || retryAfter < 0 || len(job.Payload) > 65536 || !json.Valid(job.Payload) {
		return false, ErrJobClaimLost
	}
	return failJob(ctx, tx, job.TenantID, job.JobID, worker, code, job.Attempts, retryAfter, &job)
}

func failJob(ctx context.Context, db jobExecutor, tenant, jobID, worker, code string, attempt int, retryAfter time.Duration, expected *Job) (bool, error) {
	var payload any
	queue, kind, schema := "", "", 0
	if expected != nil {
		payload = expected.Payload
		queue = expected.Queue
		kind = expected.JobType
		schema = expected.SchemaVersion
	}
	var terminal bool
	err := db.QueryRow(ctx, lockedJob+`update platform.job j set claimed_by=null,claimed_until=null,terminal_error_code=case when o.attempts>=o.max_attempts then $4 else null end,available_at=case when o.attempts>=o.max_attempts then j.available_at else clock_timestamp()+$5::interval end
 from owned o where j.tenant_id=o.tenant_id and j.job_id=o.job_id
 and o.claimed_by=$3 and o.attempts=$6 and o.claimed_until>clock_timestamp()
 and ($7::jsonb is null or (o.payload=$7 and o.queue=$8 and o.job_type=$9 and o.schema_version=$10))
 returning j.terminal_error_code is not null`, tenant, jobID, worker, code, retryAfter.String(), attempt, payload, queue, kind, schema).Scan(&terminal)
	if err != nil {
		return false, fmt.Errorf("job claim lost or update failed: %w", err)
	}
	return terminal, nil
}

// RemainingLease is a preflight check, not permission to perform arbitrary
// remote effects. Handlers still require transactional/idempotent effects and
// reconciliation; cancellation cannot undo an already-issued external request.
func (j *Jobs) RemainingLease(ctx context.Context, tenant, jobID, worker string, attempt int) (time.Duration, error) {
	if j == nil || j.pool == nil || worker == "" || attempt < 1 {
		return 0, ErrJobClaimLost
	}
	var micros int64
	err := j.pool.QueryRow(ctx, `select floor(extract(epoch from (claimed_until-clock_timestamp()))*1000000)::bigint
 from platform.job where tenant_id=$1 and job_id=$2 and claimed_by=$3 and attempts=$4
 and completed_at is null and terminal_error_code is null and claimed_until>clock_timestamp()`, tenant, jobID, worker, attempt).Scan(&micros)
	if err != nil {
		return 0, fmt.Errorf("job lease unavailable: %w", err)
	}
	if micros <= 0 || micros > int64((time.Duration(1<<63-1))/time.Microsecond) {
		return 0, ErrJobClaimLost
	}
	return time.Duration(micros) * time.Microsecond, nil
}
