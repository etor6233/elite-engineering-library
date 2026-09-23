package whatsappbridge

// AUTHORED selection of the existing owner's crashed-final-attempt transition.
// This never resets a lease/attempt or creates a replacement provider effect.
import (
	"context"
	"elite.local/enterprise/internal/platform/postgres"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5"
)

func (m *ScheduledNotifications) quarantineExpired(ctx context.Context, match json.RawMessage, jobType string) (ScheduledWorkResult, error) {
	var out ScheduledWorkResult
	tx, e := m.approvals.base.pool.Begin(ctx)
	if e != nil {
		return out, e
	}
	defer tx.Rollback(ctx)
	var job postgres.Job
	e = tx.QueryRow(ctx, `select tenant_id,job_id,queue,job_type,schema_version,payload,attempts,max_attempts from platform.job
 where tenant_id=$1 and queue='whatsapp-scheduled' and job_type=$3 and schema_version=1 and payload @> $2::jsonb
 and completed_at is null and terminal_error_code is null and attempts>=max_attempts and claimed_until<clock_timestamp()
 order by available_at,job_id for update skip locked limit 1`, m.approvals.tenant, match, jobType).Scan(&job.TenantID, &job.JobID, &job.Queue, &job.JobType, &job.SchemaVersion, &job.Payload, &job.Attempts, &job.MaxAttempts)
	if errors.Is(e, pgx.ErrNoRows) {
		return out, nil
	}
	if e != nil {
		return out, e
	}
	if e = postgres.ExhaustJobInTx(ctx, tx, job); e != nil {
		return out, e
	}
	var payload struct {
		DeliveryKey string `json:"delivery_key"`
	}
	if json.Unmarshal(job.Payload, &payload) != nil || !validScheduleKey(payload.DeliveryKey) {
		return out, ErrApproval
	}
	out = ScheduledWorkResult{Claimed: true, DeliveryKey: payload.DeliveryKey, Outcome: "LEASE_EXHAUSTED"}
	return out, tx.Commit(ctx)
}
