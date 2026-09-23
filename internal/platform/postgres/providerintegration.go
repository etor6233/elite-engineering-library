package postgres

import (
	"context"
	"errors"
	"fmt"

	"elite.local/enterprise/internal/providerintegration"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProviderIntegration struct{ pool *pgxpool.Pool }

func NewProviderIntegration(pool *pgxpool.Pool) *ProviderIntegration {
	return &ProviderIntegration{pool: pool}
}

func (r *ProviderIntegration) AcceptWebhook(ctx context.Context, value providerintegration.Receipt) (bool, error) {
	return r.acceptWebhookWithEffect(ctx, value, "provider-events", nil)
}

// The optional local effect commits with a newly inserted inbox entry/job.
// Replays never repeat the effect; failure rolls back all three components.
func (r *ProviderIntegration) acceptWebhookWithEffect(ctx context.Context, value providerintegration.Receipt, queue string, effect func(context.Context, pgx.Tx) error) (bool, error) {
	if r == nil || r.pool == nil {
		return false, providerintegration.ErrConnection
	}
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return false, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	// A startup registry is not durable authority. Lock the current connection
	// through the event/job commit; a concurrent disable must serialize with it.
	var connection string
	err = tx.QueryRow(ctx, `select connection_id from integration.provider_connection
where tenant_id=$1 and connection_id=$2 and provider_code=$3 and state='active'
for share`, value.TenantID, value.ConnectionID, value.ProviderCode).Scan(&connection)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, providerintegration.ErrConnection
	}
	if err != nil {
		return false, fmt.Errorf("resolve provider connection: %w", err)
	}
	var inserted bool
	err = tx.QueryRow(ctx, `with accepted as (
insert into integration.webhook_event(tenant_id,connection_id,provider_code,provider_event_id,event_type,body_sha256_hex,payload)
values($1,$2,$3,$4,$5,$6,$7) on conflict do nothing returning true)
select coalesce((select true from accepted),false)`, value.TenantID, value.ConnectionID, value.ProviderCode, value.ProviderEventID, value.EventType, value.BodyHash, value.Payload).Scan(&inserted)
	if err != nil {
		return false, fmt.Errorf("accept provider webhook: %w", err)
	}
	if !inserted {
		var equivalent bool
		if err = tx.QueryRow(ctx, `select provider_code=$4 and event_type=$5 and body_sha256_hex=$6 and payload=$7::jsonb
from integration.webhook_event where tenant_id=$1 and connection_id=$2 and provider_event_id=$3`, value.TenantID, value.ConnectionID, value.ProviderEventID, value.ProviderCode, value.EventType, value.BodyHash, value.Payload).Scan(&equivalent); err != nil {
			return false, err
		}
		if !equivalent {
			return false, providerintegration.ErrConflict
		}
		if err = tx.Commit(ctx); err != nil {
			return false, err
		}
		return true, nil
	}
	if effect != nil {
		if err = effect(ctx, tx); err != nil {
			return false, err
		}
	}
	_, err = tx.Exec(ctx, `insert into platform.job(tenant_id,job_id,queue,job_type,schema_version,payload,max_attempts)
values($1,gen_random_uuid(),$6,'provider.webhook.received',1,jsonb_build_object('connection_id',$2::text,'provider_code',$3::text,'provider_event_id',$4::text,'event_type',$5::text),10)`, value.TenantID, value.ConnectionID, value.ProviderCode, value.ProviderEventID, value.EventType, queue)
	if err != nil {
		return false, fmt.Errorf("enqueue provider webhook: %w", err)
	}
	if err = tx.Commit(ctx); err != nil {
		return false, err
	}
	return false, nil
}

func (r *ProviderIntegration) UpsertMapping(ctx context.Context, value providerintegration.Mapping) error {
	_, err := r.pool.Exec(ctx, `insert into integration.external_mapping(tenant_id,provider_code,resource_type,internal_id,external_id,version_token)
values($1,$2,$3,$4,$5,nullif($6,''))
on conflict(tenant_id,provider_code,resource_type,internal_id) do update
set external_id=excluded.external_id,version_token=excluded.version_token,updated_at=clock_timestamp()`, value.TenantID, value.ProviderCode, value.ResourceType, value.InternalID, value.ExternalID, value.VersionToken)
	return err
}

func (r *ProviderIntegration) RecordReconciliation(ctx context.Context, value providerintegration.Reconciliation) error {
	_, err := r.pool.Exec(ctx, `insert into integration.reconciliation_item(tenant_id,reconciliation_id,provider_code,resource_type,internal_id,external_id,state,evidence,observed_at)
values($1,$2,$3,$4,nullif($5,''),nullif($6,''),$7,$8,$9)`, value.TenantID, value.ID, value.ProviderCode, value.ResourceType, value.InternalID, value.ExternalID, value.State, value.Evidence, value.ObservedAt)
	return err
}

func (r *ProviderIntegration) ResolveReconciliation(ctx context.Context, tenantID, providerCode, reconciliationID string) error {
	result, err := r.pool.Exec(ctx, `update integration.reconciliation_item set state='resolved',resolved_at=clock_timestamp() where tenant_id=$1 and provider_code=$2 and reconciliation_id=$3 and state<>'resolved'`, tenantID, providerCode, reconciliationID)
	if err != nil {
		return err
	}
	if result.RowsAffected() != 1 {
		return providerintegration.ErrConnection
	}
	return nil
}
