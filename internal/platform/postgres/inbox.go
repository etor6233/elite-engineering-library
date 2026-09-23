package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Inbox struct{ pool *pgxpool.Pool }

func NewInbox(pool *pgxpool.Pool) *Inbox { return &Inbox{pool: pool} }

// Process executes local database effects and the inbox receipt atomically.
// Returning processed=false means a completed receipt already existed.
func (i *Inbox) Process(ctx context.Context, tenantID, consumerName, eventID string, handler func(context.Context, pgx.Tx) (string, error)) (processed bool, err error) {
	if tenantID == "" || consumerName == "" || len(consumerName) > 128 || eventID == "" || handler == nil {
		return false, fmt.Errorf("invalid inbox parameters")
	}
	tx, err := i.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return false, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var inserted bool
	err = tx.QueryRow(ctx, `
with receipt as (
  insert into platform.consumer_inbox(tenant_id,consumer_name,event_id)
  values($1,$2,$3)
  on conflict do nothing
  returning true
)
select coalesce((select true from receipt),false)`, tenantID, consumerName, eventID).Scan(&inserted)
	if err != nil {
		return false, err
	}
	if !inserted {
		var completed bool
		if err := tx.QueryRow(ctx, `select completed_at is not null from platform.consumer_inbox where tenant_id=$1 and consumer_name=$2 and event_id=$3`, tenantID, consumerName, eventID).Scan(&completed); err != nil {
			return false, err
		}
		if !completed {
			return false, fmt.Errorf("inbox receipt exists but is incomplete")
		}
		if err := tx.Commit(ctx); err != nil {
			return false, err
		}
		return false, nil
	}
	resultHash, err := handler(ctx, tx)
	if err != nil {
		return false, err
	}
	if resultHash == "" {
		resultHash = fmt.Sprintf("%064x", 0)
	}
	result, err := tx.Exec(ctx, `update platform.consumer_inbox set completed_at=clock_timestamp(),result_sha256_hex=$4 where tenant_id=$1 and consumer_name=$2 and event_id=$3 and completed_at is null`, tenantID, consumerName, eventID, resultHash)
	if err != nil {
		return false, err
	}
	if result.RowsAffected() != 1 {
		return false, fmt.Errorf("inbox completion lost")
	}
	if err := tx.Commit(ctx); err != nil {
		return false, err
	}
	return true, nil
}
