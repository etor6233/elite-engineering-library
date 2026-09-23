package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"elite.local/enterprise/internal/leadstream"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LeadIngress struct{ pool *pgxpool.Pool }

func NewLeadIngress(pool *pgxpool.Pool) *LeadIngress { return &LeadIngress{pool: pool} }

func (r *LeadIngress) Record(ctx context.Context, raw leadstream.RawEvent, candidate *leadstream.LeadCandidate, normalizationCode string) (leadstream.Receipt, error) {
	if r == nil || r.pool == nil {
		return leadstream.Receipt{}, errors.New("postgres lead ingress: nil pool")
	}
	if err := raw.Validate(); err != nil {
		return leadstream.Receipt{}, err
	}
	if candidate == nil && normalizationCode == "" {
		return leadstream.Receipt{}, errors.New("postgres lead ingress: candidate or rejection required")
	}
	if candidate != nil {
		if normalizationCode != "" {
			return leadstream.Receipt{}, errors.New("postgres lead ingress: candidate and rejection conflict")
		}
		if err := candidate.Validate(); err != nil {
			return leadstream.Receipt{}, err
		}
	}
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return leadstream.Receipt{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	state := leadstream.ReceiptNormalized
	if candidate == nil {
		state = leadstream.ReceiptRejected
	}
	result, err := tx.Exec(ctx, `insert into integration.lead_ingress_raw
(tenant_id,organization_id,provider,provider_event_id,event_type,source,schema_version,occurred_at,received_at,payload_redacted,source_payload_sha256,stored_payload_sha256,redaction_profile,state,normalization_error_code)
values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,nullif($15,'')) on conflict do nothing`,
		raw.TenantID, raw.OrganizationID, raw.Provider, raw.ProviderEventID, raw.EventType, raw.Source,
		raw.SchemaVersion, raw.OccurredAt, raw.ReceivedAt, raw.Payload, raw.SourceSHA256, raw.StoredSHA256, raw.RedactionProfile, string(state), normalizationCode)
	if err != nil {
		return leadstream.Receipt{}, err
	}
	if result.RowsAffected() == 0 {
		var hash string
		if err := tx.QueryRow(ctx, `select source_payload_sha256 from integration.lead_ingress_raw where tenant_id=$1 and provider=$2 and provider_event_id=$3`, raw.TenantID, raw.Provider, raw.ProviderEventID).Scan(&hash); err != nil {
			return leadstream.Receipt{}, err
		}
		if hash != raw.SourceSHA256 {
			return leadstream.Receipt{}, leadstream.ErrDivergentDuplicate
		}
		if err := tx.Commit(ctx); err != nil {
			return leadstream.Receipt{}, err
		}
		return leadstream.Receipt{State: leadstream.ReceiptDuplicate, SourceSHA256: hash}, nil
	}
	if candidate != nil {
		fields, err := json.Marshal(candidate.Fields)
		if err != nil {
			return leadstream.Receipt{}, err
		}
		_, err = tx.Exec(ctx, `insert into integration.lead_candidate
(tenant_id,provider,provider_event_id,organization_id,provider_lead_id,form_id,campaign_id,ad_group_id,creative_id,asset_group_id,click_id,lead_stage,source_kind,submitted_at,is_test,fields,contact_eligibility)
values($1,$2,$3,$4,$5,nullif($6,''),nullif($7,''),nullif($8,''),nullif($9,''),nullif($10,''),nullif($11,''),nullif($12,''),nullif($13,''),$14,$15,$16,'pending_policy')`,
			candidate.TenantID, candidate.Provider, raw.ProviderEventID, candidate.OrganizationID, candidate.ProviderLeadID,
			candidate.FormID, candidate.CampaignID, candidate.AdGroupID, candidate.CreativeID, candidate.AssetGroupID,
			candidate.ClickID, candidate.LeadStage, candidate.SourceKind, candidate.SubmittedAt, candidate.IsTest, fields)
		if err != nil {
			return leadstream.Receipt{}, err
		}
	}
	eventPayload, _ := json.Marshal(map[string]any{
		"provider": raw.Provider, "provider_event_id": raw.ProviderEventID,
		"organization_id": raw.OrganizationID, "state": state, "is_test": candidate != nil && candidate.IsTest,
	})
	outboxID := leadstream.StableUUID(raw.TenantID, raw.Provider, raw.ProviderEventID, "received")
	_, err = tx.Exec(ctx, `insert into platform.outbox_event
(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)
values($1,$2,'provider-lead',$3,1,$4,1,clock_timestamp(),$5)`, raw.TenantID, outboxID, raw.Provider+":"+raw.ProviderEventID,
		"provider-lead."+string(state), eventPayload)
	if err != nil {
		return leadstream.Receipt{}, fmt.Errorf("write lead ingress outbox: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return leadstream.Receipt{}, err
	}
	return leadstream.Receipt{State: state, SourceSHA256: raw.SourceSHA256}, nil
}
