package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OutboxEvent struct {
	TenantID         string
	EventID          string
	AggregateType    string
	AggregateID      string
	AggregateVersion int64
	EventType        string
	SchemaVersion    int
	OccurredAt       time.Time
	Payload          json.RawMessage
	Headers          json.RawMessage
	// Attempts is the monotonically increasing claim generation. Never reset it.
	Attempts int
}

var ErrOutboxClaimLost = errors.New("outbox event claim lost")

type Outbox struct{ pool *pgxpool.Pool }

func NewOutbox(pool *pgxpool.Pool) *Outbox { return &Outbox{pool: pool} }

func (o *Outbox) Claim(ctx context.Context, workerID string, lease time.Duration, batchSize int) ([]OutboxEvent, error) {
	if o == nil || o.pool == nil || workerID == "" || lease < time.Microsecond || batchSize < 1 || batchSize > 1000 {
		return nil, fmt.Errorf("invalid outbox claim parameters")
	}
	rows, err := o.pool.Query(ctx, `
with candidates as (
  select tenant_id,event_id
    from platform.outbox_event
   where published_at is null
     and available_at <= clock_timestamp()
     and (claimed_until is null or claimed_until < clock_timestamp())
   order by available_at,occurred_at,event_id
   for update skip locked
   limit $1
)
update platform.outbox_event as event
   set claimed_by=$2,
       claimed_until=clock_timestamp()+$3::interval,
       attempts=attempts+1
  from candidates
 where event.tenant_id=candidates.tenant_id
   and event.event_id=candidates.event_id
returning event.tenant_id,event.event_id,event.aggregate_type,event.aggregate_id,
          event.aggregate_version,event.event_type,event.schema_version,
          event.occurred_at,event.payload,event.headers,event.attempts`, batchSize, workerID, lease.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	events := make([]OutboxEvent, 0, batchSize)
	for rows.Next() {
		var event OutboxEvent
		if err := rows.Scan(&event.TenantID, &event.EventID, &event.AggregateType, &event.AggregateID, &event.AggregateVersion, &event.EventType, &event.SchemaVersion, &event.OccurredAt, &event.Payload, &event.Headers, &event.Attempts); err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, rows.Err()
}

// Lock before checking the DB clock: a wait on this row may outlive the lease.
// Worker identity alone cannot distinguish a restarted or concurrent attempt.
const lockedOutbox = `with owned as materialized (
 select tenant_id,event_id,claimed_by,claimed_until,attempts
 from platform.outbox_event where tenant_id=$1 and event_id=$2
 and published_at is null for update
) `

func (o *Outbox) MarkPublished(ctx context.Context, tenantID, eventID, workerID string, attempt int) error {
	if o == nil || o.pool == nil || workerID == "" || attempt < 1 {
		return ErrOutboxClaimLost
	}
	result, err := o.pool.Exec(ctx, lockedOutbox+`update platform.outbox_event e set published_at=clock_timestamp(),claimed_by=null,claimed_until=null,last_error_code=null
 from owned o where e.tenant_id=o.tenant_id and e.event_id=o.event_id
 and o.claimed_by=$3 and o.attempts=$4 and o.claimed_until>clock_timestamp()`, tenantID, eventID, workerID, attempt)
	if err != nil {
		return err
	}
	if result.RowsAffected() != 1 {
		return ErrOutboxClaimLost
	}
	return nil
}

func (o *Outbox) Release(ctx context.Context, tenantID, eventID, workerID, errorCode string, attempt int, retryAfter time.Duration) error {
	if o == nil || o.pool == nil || workerID == "" || attempt < 1 {
		return ErrOutboxClaimLost
	}
	if retryAfter < 0 || errorCode == "" {
		return fmt.Errorf("invalid outbox release parameters")
	}
	result, err := o.pool.Exec(ctx, lockedOutbox+`update platform.outbox_event e set claimed_by=null,claimed_until=null,last_error_code=$4,available_at=clock_timestamp()+$5::interval
 from owned o where e.tenant_id=o.tenant_id and e.event_id=o.event_id
 and o.claimed_by=$3 and o.attempts=$6 and o.claimed_until>clock_timestamp()`, tenantID, eventID, workerID, errorCode, retryAfter.String(), attempt)
	if err != nil {
		return err
	}
	if result.RowsAffected() != 1 {
		return ErrOutboxClaimLost
	}
	return nil
}

// RemainingLease is a conservative preflight, not a remote-effect fence. The
// publisher still needs stable idempotency and reconciliation after ambiguity.
func (o *Outbox) RemainingLease(ctx context.Context, tenantID, eventID, workerID string, attempt int) (time.Duration, error) {
	if o == nil || o.pool == nil || workerID == "" || attempt < 1 {
		return 0, ErrOutboxClaimLost
	}
	var micros int64
	err := o.pool.QueryRow(ctx, `select floor(extract(epoch from (claimed_until-clock_timestamp()))*1000000)::bigint
 from platform.outbox_event where tenant_id=$1 and event_id=$2 and claimed_by=$3 and attempts=$4
 and published_at is null and claimed_until>clock_timestamp()`, tenantID, eventID, workerID, attempt).Scan(&micros)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, ErrOutboxClaimLost
	}
	if err != nil {
		return 0, err
	}
	if micros <= 0 || micros > int64((time.Duration(1<<63-1))/time.Microsecond) {
		return 0, ErrOutboxClaimLost
	}
	return time.Duration(micros) * time.Microsecond, nil
}
