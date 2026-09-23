package postgres

import (
	"context"
	"encoding/json"
	"errors"

	"elite.local/enterprise/internal/leadstream"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MetaLeadWebhookStore struct{ pool *pgxpool.Pool }

func NewMetaLeadWebhookStore(pool *pgxpool.Pool) *MetaLeadWebhookStore {
	return &MetaLeadWebhookStore{pool: pool}
}

func (s *MetaLeadWebhookStore) RecordMetaWebhook(ctx context.Context, batch leadstream.MetaWebhookBatch) (leadstream.MetaWebhookReceipt, error) {
	if s == nil || s.pool == nil {
		return leadstream.MetaWebhookReceipt{}, errors.New("postgres Meta webhook: nil pool")
	}
	if err := batch.Validate(); err != nil {
		return leadstream.MetaWebhookReceipt{}, err
	}
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return leadstream.MetaWebhookReceipt{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	batchResult, err := tx.Exec(ctx, `insert into integration.meta_lead_webhook_batch
(tenant_id,organization_id,payload_sha256,payload,received_at,signal_count)
values($1,$2,$3,$4,$5,$6) on conflict do nothing`, batch.TenantID, batch.OrganizationID, batch.PayloadSHA256, batch.Payload, batch.ReceivedAt, len(batch.Signals))
	if err != nil {
		return leadstream.MetaWebhookReceipt{}, err
	}
	if batchResult.RowsAffected() == 0 {
		var organizationID string
		if err = tx.QueryRow(ctx, `select organization_id from integration.meta_lead_webhook_batch where tenant_id=$1 and payload_sha256=$2`, batch.TenantID, batch.PayloadSHA256).Scan(&organizationID); err != nil {
			return leadstream.MetaWebhookReceipt{}, err
		}
		if organizationID != batch.OrganizationID {
			return leadstream.MetaWebhookReceipt{}, leadstream.ErrDivergentDuplicate
		}
	}
	receipt := leadstream.MetaWebhookReceipt{}
	for _, signal := range batch.Signals {
		result, execErr := tx.Exec(ctx, `insert into integration.meta_lead_webhook_signal
(tenant_id,leadgen_id,organization_id,batch_payload_sha256,value_sha256,form_id,page_id,ad_id,ad_group_id,occurred_at,state)
values($1,$2,$3,$4,$5,$6,$7,nullif($8,''),nullif($9,''),$10,'pending_retrieval') on conflict do nothing`,
			batch.TenantID, signal.LeadgenID, batch.OrganizationID, batch.PayloadSHA256, signal.ValueSHA256,
			signal.FormID, signal.PageID, signal.AdID, signal.AdGroupID, signal.OccurredAt)
		if execErr != nil {
			return leadstream.MetaWebhookReceipt{}, execErr
		}
		if result.RowsAffected() == 0 {
			var hash string
			if err = tx.QueryRow(ctx, `select value_sha256 from integration.meta_lead_webhook_signal where tenant_id=$1 and leadgen_id=$2`, batch.TenantID, signal.LeadgenID).Scan(&hash); err != nil {
				return leadstream.MetaWebhookReceipt{}, err
			}
			if hash != signal.ValueSHA256 {
				return leadstream.MetaWebhookReceipt{}, leadstream.ErrDivergentDuplicate
			}
			receipt.DuplicateSignals++
			continue
		}
		eventPayload, _ := json.Marshal(map[string]any{
			"provider": "meta_lead_ads", "leadgen_id": signal.LeadgenID, "form_id": signal.FormID,
			"page_id": signal.PageID, "organization_id": batch.OrganizationID, "state": "pending_retrieval",
			"source_payload_sha256": batch.PayloadSHA256,
		})
		eventID := leadstream.StableUUID(batch.TenantID, "meta_lead_ads", signal.LeadgenID, "signal")
		if _, err = tx.Exec(ctx, `insert into platform.outbox_event
(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)
values($1,$2,'provider-lead-signal',$3,1,'meta-lead.signal-received',1,$4,$5)`, batch.TenantID, eventID, "meta_lead_ads:"+signal.LeadgenID, signal.OccurredAt, eventPayload); err != nil {
			return leadstream.MetaWebhookReceipt{}, err
		}
		receipt.NewSignals++
	}
	if err = tx.Commit(ctx); err != nil {
		return leadstream.MetaWebhookReceipt{}, err
	}
	return receipt, nil
}
