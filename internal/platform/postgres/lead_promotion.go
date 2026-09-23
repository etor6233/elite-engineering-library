package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"elite.local/enterprise/internal/leadpromotion"
	"elite.local/enterprise/internal/leadstream"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LeadPromotion struct{ pool *pgxpool.Pool }

func NewLeadPromotion(pool *pgxpool.Pool) *LeadPromotion { return &LeadPromotion{pool: pool} }

func (r *LeadPromotion) Promote(ctx context.Context, c leadpromotion.Command) (leadpromotion.Receipt, error) {
	if r == nil || r.pool == nil {
		return leadpromotion.Receipt{}, errors.New("postgres lead promotion: nil pool")
	}
	requestHash, err := c.RequestSHA256()
	if err != nil {
		return leadpromotion.Receipt{}, err
	}
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return leadpromotion.Receipt{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	result, err := tx.Exec(ctx, `insert into platform.idempotency_record(tenant_id,scope,idempotency_key,request_sha256_hex,status,locked_until,expires_at) values($1,'lead-promotion',$2,$3,'processing',clock_timestamp()+interval '30 seconds',clock_timestamp()+interval '24 hours') on conflict do nothing`, c.TenantID, c.IdempotencyKey, requestHash)
	if err != nil {
		return leadpromotion.Receipt{}, err
	}
	if result.RowsAffected() == 0 {
		var priorHash, status, resourceID string
		if err = tx.QueryRow(ctx, `select request_sha256_hex,status,coalesce(resource_id,'') from platform.idempotency_record where tenant_id=$1 and scope='lead-promotion' and idempotency_key=$2`, c.TenantID, c.IdempotencyKey).Scan(&priorHash, &status, &resourceID); err != nil {
			return leadpromotion.Receipt{}, err
		}
		if priorHash != requestHash || status != "completed" || resourceID == "" {
			return leadpromotion.Receipt{}, leadpromotion.ErrConflict
		}
		if err = tx.Commit(ctx); err != nil {
			return leadpromotion.Receipt{}, err
		}
		return leadpromotion.Receipt{LeadID: resourceID, Replayed: true}, nil
	}
	var organization string
	var isTest bool
	var fieldsJSON []byte
	var sourceHash string
	err = tx.QueryRow(ctx, `select c.organization_id,c.is_test,c.fields,r.source_payload_sha256 from integration.lead_candidate c join integration.lead_ingress_raw r using(tenant_id,provider,provider_event_id) where c.tenant_id=$1 and c.provider=$2 and c.provider_event_id=$3 and c.contact_eligibility='pending_policy' for update of c`, c.TenantID, c.Provider, c.ProviderEventID).Scan(&organization, &isTest, &fieldsJSON, &sourceHash)
	if errors.Is(err, pgx.ErrNoRows) {
		return leadpromotion.Receipt{}, leadpromotion.ErrNotFound
	}
	if err != nil {
		return leadpromotion.Receipt{}, err
	}
	if isTest {
		return leadpromotion.Receipt{}, leadpromotion.ErrTestLead
	}
	var fields []leadpromotion.Field
	if err = json.Unmarshal(fieldsJSON, &fields); err != nil {
		return leadpromotion.Receipt{}, fmt.Errorf("decode candidate fields: %w", err)
	}
	contact, err := leadpromotion.ContactPayload(fields, c.FieldMapping)
	if err != nil {
		return leadpromotion.Receipt{}, err
	}
	mappingJSON, err := json.Marshal(c.FieldMapping)
	if err != nil {
		return leadpromotion.Receipt{}, err
	}
	_, err = tx.Exec(ctx, `insert into crm.lead(tenant_id,lead_id,organization_id,lifecycle_state,source_code,contact_payload,consent_required) values($1,$2,$3,'new',$4,$5,true)`, c.TenantID, c.LeadID, organization, c.Provider, contact)
	if err != nil {
		return leadpromotion.Receipt{}, err
	}
	_, err = tx.Exec(ctx, `insert into crm.consent_evidence(tenant_id,consent_id,lead_id,purpose_code,policy_version,decision,occurred_at,evidence_sha256_hex) values($1,$2,$3,$4,$5,'granted',$6,$7)`, c.TenantID, c.ConsentID, c.LeadID, c.PurposeCode, c.PolicyVersion, c.DecisionAt.UTC(), c.EvidenceSHA256)
	if err != nil {
		return leadpromotion.Receipt{}, err
	}
	_, err = tx.Exec(ctx, `insert into integration.lead_promotion(tenant_id,provider,provider_event_id,lead_id,consent_id,mapping_version,field_mapping,source_payload_sha256,decision_at) values($1,$2,$3,$4,$5,$6,$7,$8,$9)`, c.TenantID, c.Provider, c.ProviderEventID, c.LeadID, c.ConsentID, c.MappingVersion, mappingJSON, sourceHash, c.DecisionAt.UTC())
	if err != nil {
		return leadpromotion.Receipt{}, err
	}
	eventID := leadstream.StableUUID(c.TenantID, c.Provider, c.ProviderEventID, "promoted")
	payload, _ := json.Marshal(map[string]string{"provider": c.Provider, "provider_event_id": c.ProviderEventID, "lead_id": c.LeadID, "organization_id": organization, "purpose_code": c.PurposeCode})
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload) values($1,$2,'lead',$3,1,'lead.promoted',1,clock_timestamp(),$4)`, c.TenantID, eventID, c.LeadID, payload)
	if err != nil {
		return leadpromotion.Receipt{}, err
	}
	_, err = tx.Exec(ctx, `update platform.idempotency_record set status='completed',response_code=201,response_body=jsonb_build_object('lead_id',$3::text),resource_type='lead',resource_id=$3,locked_until=null where tenant_id=$1 and scope='lead-promotion' and idempotency_key=$2 and status='processing'`, c.TenantID, c.IdempotencyKey, c.LeadID)
	if err != nil {
		return leadpromotion.Receipt{}, err
	}
	if err = tx.Commit(ctx); err != nil {
		return leadpromotion.Receipt{}, err
	}
	return leadpromotion.Receipt{LeadID: c.LeadID}, nil
}
