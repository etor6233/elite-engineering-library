package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"elite.local/enterprise/internal/electromobility"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var electromobilityPublicModelsSQL = `select m.model_id,m.model_code,m.display_name,m.vehicle_class,m.specification from catalog.vehicle_model m join platform.tenant t on t.tenant_id=m.tenant_id where t.tenant_code=$1 and t.status='active' and m.lifecycle_state='active' and m.publicly_visible=true order by m.display_name,m.model_id`

type Electromobility struct{ pool *pgxpool.Pool }

func NewElectromobility(pool *pgxpool.Pool) *Electromobility { return &Electromobility{pool: pool} }

func (r *Electromobility) ListPublicModels(ctx context.Context, tenantCode string) ([]electromobility.Model, error) {
	rows, err := r.pool.Query(ctx, electromobilityPublicModelsSQL, tenantCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	models := []electromobility.Model{}
	for rows.Next() {
		var model electromobility.Model
		if err := rows.Scan(&model.ID, &model.Code, &model.DisplayName, &model.VehicleClass, &model.Specification); err != nil {
			return nil, err
		}
		models = append(models, model)
	}
	return models, rows.Err()
}
func (r *Electromobility) CreateModel(ctx context.Context, tenantID, eventID string, model electromobility.Model) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err = r.createModelTx(ctx, tx, tenantID, eventID, model); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// AUTHORED extraction: original model insert and outbox SQL/arguments retained.
func (r *Electromobility) createModelTx(ctx context.Context, tx pgx.Tx, tenantID, eventID string, model electromobility.Model) error {
	_, err := tx.Exec(ctx, `insert into catalog.vehicle_model(tenant_id,model_id,model_code,display_name,vehicle_class,lifecycle_state,publicly_visible,specification)values($1,$2,$3,$4,$5,'draft',false,$6)`, tenantID, model.ID, model.Code, model.DisplayName, model.VehicleClass, model.Specification)
	if err != nil {
		return err
	}
	payload, _ := json.Marshal(model)
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)values($1,$2,'vehicle-model',$3,1,'vehicle-model.created',1,clock_timestamp(),$4)`, tenantID, eventID, model.ID, payload)
	if err != nil {
		return err
	}
	return nil
}
func (r *Electromobility) ResolvePublicOrganization(ctx context.Context, tenantCode, organizationCode string) (string, string, error) {
	var tenantID, organizationID string
	err := r.pool.QueryRow(ctx, `select t.tenant_id,o.organization_id from platform.tenant t join org.organization o on o.tenant_id=t.tenant_id where t.tenant_code=$1 and t.status='active' and o.organization_code=$2 and o.status='active'`, tenantCode, organizationCode).Scan(&tenantID, &organizationID)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", "", electromobility.ErrNotFound
	}
	return tenantID, organizationID, err
}
func (r *Electromobility) CreateLead(ctx context.Context, lead electromobility.Lead, eventID, idempotencyKey string) (string, bool, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return "", false, err
	}
	defer tx.Rollback(ctx)
	result, err := tx.Exec(ctx, `insert into platform.idempotency_record(tenant_id,scope,idempotency_key,request_sha256_hex,status,locked_until,expires_at)values($1,'public-lead',$2,$3,'processing',clock_timestamp()+interval '30 seconds',clock_timestamp()+interval '24 hours') on conflict do nothing`, lead.TenantID, idempotencyKey, lead.ConsentEvidenceHash)
	if err != nil {
		return "", false, err
	}
	if result.RowsAffected() == 0 {
		var requestHash, status, resourceID string
		err = tx.QueryRow(ctx, `select request_sha256_hex,status,coalesce(resource_id,'') from platform.idempotency_record where tenant_id=$1 and scope='public-lead' and idempotency_key=$2`, lead.TenantID, idempotencyKey).Scan(&requestHash, &status, &resourceID)
		if err != nil {
			return "", false, err
		}
		if requestHash != lead.ConsentEvidenceHash || status != "completed" || resourceID == "" {
			return "", false, electromobility.ErrConflict
		}
		if err = tx.Commit(ctx); err != nil {
			return "", false, err
		}
		return resourceID, true, nil
	}
	var model any = nil
	if lead.ModelID != "" {
		model = lead.ModelID
	}
	_, err = tx.Exec(ctx, `insert into crm.lead(tenant_id,lead_id,organization_id,model_id,lifecycle_state,source_code,contact_payload,consent_required)values($1,$2,$3,$4,'new',$5,$6,true)`, lead.TenantID, lead.ID, lead.OrganizationID, model, lead.SourceCode, lead.Contact)
	if err != nil {
		return "", false, err
	}
	_, err = tx.Exec(ctx, `insert into crm.consent_evidence(tenant_id,consent_id,lead_id,purpose_code,policy_version,decision,occurred_at,evidence_sha256_hex)values($1,$2,$3,'sales-contact','v1','granted',clock_timestamp(),$4)`, lead.TenantID, lead.ConsentID, lead.ID, lead.ConsentEvidenceHash)
	if err != nil {
		return "", false, err
	}
	payload, _ := json.Marshal(map[string]string{"lead_id": lead.ID, "organization_id": lead.OrganizationID, "source_code": lead.SourceCode})
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)values($1,$2,'lead',$3,1,'lead.captured',1,clock_timestamp(),$4)`, lead.TenantID, eventID, lead.ID, payload)
	if err != nil {
		return "", false, fmt.Errorf("write lead outbox: %w", err)
	}
	_, err = tx.Exec(ctx, `update platform.idempotency_record set status='completed',response_code=202,response_body=jsonb_build_object('lead_id',$3::text,'status','accepted'),resource_type='lead',resource_id=$3,locked_until=null where tenant_id=$1 and scope='public-lead' and idempotency_key=$2 and status='processing'`, lead.TenantID, idempotencyKey, lead.ID)
	if err != nil {
		return "", false, err
	}
	if err = tx.Commit(ctx); err != nil {
		return "", false, err
	}
	return lead.ID, false, nil
}
