package postgres

import (
	"context"
	"elite.local/enterprise/internal/fulfillment"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Fulfillment struct{ pool *pgxpool.Pool }

func NewFulfillment(pool *pgxpool.Pool) *Fulfillment { return &Fulfillment{pool: pool} }
func (r *Fulfillment) CreateShipment(ctx context.Context, tenant, eventID string, v fulfillment.Shipment) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, `insert into logistics.shipment(tenant_id,shipment_id,provider_code,provider_reference,origin_organization_id,destination_organization_id,state)values($1,$2,$3,nullif($4,''),$5,$6,'planned')`, tenant, v.ID, v.ProviderCode, v.ProviderReference, v.OriginOrganizationID, v.DestinationOrganizationID)
	if err != nil {
		return err
	}
	if err = outbox(ctx, tx, tenant, eventID, "shipment", v.ID, 1, "shipment.planned", `jsonb_build_object('provider_code',$7::text)`, v.ProviderCode); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
func (r *Fulfillment) TransitionShipment(ctx context.Context, tenant, originOrganization, destinationOrganization, id, current, target, reference, eventID string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	result, err := tx.Exec(ctx, `update logistics.shipment set state=$6,provider_reference=coalesce(nullif($7,''),provider_reference),last_event_at=clock_timestamp() where tenant_id=$1 and origin_organization_id=$2 and destination_organization_id=$3 and shipment_id=$4 and state=$5`, tenant, originOrganization, destinationOrganization, id, current, target, reference)
	if err != nil {
		return err
	}
	if result.RowsAffected() != 1 {
		return fulfillment.ErrConflict
	}
	if err = outbox(ctx, tx, tenant, eventID, "shipment", id, 1, "shipment."+target, `'{}'::jsonb`); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
func (r *Fulfillment) OpenServiceCase(ctx context.Context, tenant, eventID string, v fulfillment.ServiceCase) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err = openServiceCaseInTx(ctx, tx, tenant, eventID, v); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// AUTHORED transaction composition only: original case writer SQL preserved.
func openServiceCaseInTx(ctx context.Context, tx pgx.Tx, tenant, eventID string, v fulfillment.ServiceCase) error {
	var err error
	result, err := tx.Exec(ctx, `insert into service_ops.service_case(tenant_id,service_case_id,stock_unit_id,organization_id,state,severity,description,version)select $1,$2,s.stock_unit_id,$4,'opened',$5,$6,1 from inventory.stock_unit s where s.tenant_id=$1 and s.stock_unit_id=$3 and s.organization_id=$4`, tenant, v.ID, v.StockUnitID, v.OrganizationID, v.Severity, v.Description)
	if err != nil {
		return err
	}
	if result.RowsAffected() != 1 {
		return fulfillment.ErrConflict
	}
	if err = outbox(ctx, tx, tenant, eventID, "service-case", v.ID, 1, "service-case.opened", `jsonb_build_object('severity',$7::text)`, v.Severity); err != nil {
		return err
	}
	return nil
}
func (r *Fulfillment) TransitionServiceCase(ctx context.Context, tenant, organization, id, current, target string, version int64, eventID string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	result, err := tx.Exec(ctx, `update service_ops.service_case set state=$6,version=version+1,closed_at=case when $6='closed' then clock_timestamp() else null end where tenant_id=$1 and organization_id=$2 and service_case_id=$3 and state=$4 and version=$5`, tenant, organization, id, current, version, target)
	if err != nil {
		return err
	}
	if result.RowsAffected() != 1 {
		return fulfillment.ErrConflict
	}
	if err = outbox(ctx, tx, tenant, eventID, "service-case", id, version+1, "service-case."+target, `'{}'::jsonb`); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
func (r *Fulfillment) CreateRecall(ctx context.Context, tenant, eventID string, v fulfillment.Recall) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, `insert into service_ops.recall(tenant_id,recall_id,recall_code,title,severity,status)values($1,$2,$3,$4,$5,'draft')`, tenant, v.ID, v.Code, v.Title, v.Severity)
	if err != nil {
		return err
	}
	if err = outbox(ctx, tx, tenant, eventID, "recall", v.ID, 1, "recall.created", `jsonb_build_object('recall_code',$7::text)`, v.Code); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
func (r *Fulfillment) ActivateRecall(ctx context.Context, tenant, id, eventID string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	result, err := tx.Exec(ctx, `update service_ops.recall set status='active',published_at=clock_timestamp() where tenant_id=$1 and recall_id=$2 and status='draft' and published_at is null`, tenant, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() != 1 {
		return fulfillment.ErrConflict
	}
	if err = outbox(ctx, tx, tenant, eventID, "recall", id, 2, "recall.activated", `'{}'::jsonb`); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
func (r *Fulfillment) AddRecallUnit(ctx context.Context, tenant, organization, recallID, stockID, eventID string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	result, err := tx.Exec(ctx, `insert into service_ops.recall_unit(tenant_id,recall_id,stock_unit_id,state)select r.tenant_id,r.recall_id,s.stock_unit_id,'identified' from service_ops.recall r join inventory.stock_unit s on s.tenant_id=r.tenant_id where r.tenant_id=$1 and s.organization_id=$2 and r.recall_id=$3 and r.status='active' and s.stock_unit_id=$4 on conflict do nothing`, tenant, organization, recallID, stockID)
	if err != nil {
		return err
	}
	if result.RowsAffected() != 1 {
		return fulfillment.ErrConflict
	}
	if err = outbox(ctx, tx, tenant, eventID, "recall", recallID, 3, "recall.unit-identified", `jsonb_build_object('stock_unit_id',$7::text)`, stockID); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
func (r *Fulfillment) CreateOrganization(ctx context.Context, tenant, eventID string, v fulfillment.Organization) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err = createOrganizationInTx(ctx, tx, tenant, eventID, v); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// AUTHORED transaction extraction; original owner SQL and guards preserved.
func createOrganizationInTx(ctx context.Context, tx pgx.Tx, tenant, eventID string, v fulfillment.Organization) error {
	result, err := tx.Exec(ctx, `insert into org.organization(tenant_id,organization_id,parent_organization_id,organization_code,display_name,organization_type,status,version)
		select t.tenant_id,$2,nullif($3,''),$4,$5,$6,'provisioning',1
		from platform.tenant t left join org.organization p on p.tenant_id=t.tenant_id and p.organization_id=nullif($3,'')
		where t.tenant_id=$1 and t.status='active' and (
		  ($3='' and $6 in ('enterprise','franchisor')) or
		  ($3<>'' and p.status='active' and (($6='franchisee' and p.organization_type in ('enterprise','franchisor')) or ($6 in ('store','service_center','warehouse') and p.organization_type in ('enterprise','franchisor','franchisee')) or ($6='factory' and p.organization_type in ('enterprise','franchisor'))))
		)`, tenant, v.ID, v.ParentOrganizationID, v.Code, v.DisplayName, v.Type)
	if err != nil {
		return fulfillmentConstraint(err)
	}
	if result.RowsAffected() != 1 {
		return fulfillment.ErrConflict
	}
	if err = outbox(ctx, tx, tenant, eventID, "organization", v.ID, 1, "organization.provisioning", `jsonb_build_object('parent_organization_id',nullif($7::text,''),'organization_type',$8::text)`, v.ParentOrganizationID, v.Type); err != nil {
		return err
	}
	return nil
}
func (r *Fulfillment) TransitionOrganization(ctx context.Context, tenant, id, current, target string, version int64, eventID string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err = transitionOrganizationInTx(ctx, tx, tenant, id, current, target, version, eventID); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// AUTHORED transaction extraction; original owner SQL and guards preserved.
func transitionOrganizationInTx(ctx context.Context, tx pgx.Tx, tenant, id, current, target string, version int64, eventID string) error {
	result, err := tx.Exec(ctx, `update org.organization o set status=$4,version=version+1,updated_at=clock_timestamp()
		where tenant_id=$1 and organization_id=$2 and status=$3 and version=$5
		and not ($4='closed' and (exists(select 1 from org.organization child where child.tenant_id=o.tenant_id and child.parent_organization_id=o.organization_id and child.status<>'closed') or exists(select 1 from franchise.agreement a where a.tenant_id=o.tenant_id and a.franchise_organization_id=o.organization_id and a.status in ('active','suspended'))))`, tenant, id, current, target, version)
	if err != nil {
		return err
	}
	if result.RowsAffected() != 1 {
		return fulfillment.ErrConflict
	}
	if err = outbox(ctx, tx, tenant, eventID, "organization", id, version+1, "organization."+target, `'{}'::jsonb`); err != nil {
		return err
	}
	return nil
}
func (r *Fulfillment) CreateAgreement(ctx context.Context, tenant, eventID string, v fulfillment.Agreement) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err = createAgreementInTx(ctx, tx, tenant, eventID, v); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// AUTHORED transaction extraction; original owner SQL and guards preserved.
func createAgreementInTx(ctx context.Context, tx pgx.Tx, tenant, eventID string, v fulfillment.Agreement) error {
	result, err := tx.Exec(ctx, `insert into franchise.agreement(tenant_id,agreement_id,franchise_organization_id,territory_code,terms_version,starts_on,ends_on,status,version)
		select $1,$2,o.organization_id,$4,$5,$6,$7,'draft',1 from org.organization o where o.tenant_id=$1 and o.organization_id=$3 and o.organization_type='franchisee' and o.status='active'`, tenant, v.ID, v.OrganizationID, v.TerritoryCode, v.TermsVersion, v.StartsOn, v.EndsOn)
	if err != nil {
		return fulfillmentConstraint(err)
	}
	if result.RowsAffected() != 1 {
		return fulfillment.ErrConflict
	}
	if err = outbox(ctx, tx, tenant, eventID, "franchise-agreement", v.ID, 1, "franchise-agreement.created", `jsonb_build_object('organization_id',$7::text)`, v.OrganizationID); err != nil {
		return err
	}
	return nil
}
func (r *Fulfillment) TransitionAgreement(ctx context.Context, tenant, organization, id, current, target string, version int64, eventID string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err = transitionAgreementInTx(ctx, tx, tenant, organization, id, current, target, version, eventID); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// AUTHORED transaction extraction; original owner SQL and guards preserved.
func transitionAgreementInTx(ctx context.Context, tx pgx.Tx, tenant, organization, id, current, target string, version int64, eventID string) error {
	result, err := tx.Exec(ctx, `update franchise.agreement a set status=$5,version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and franchise_organization_id=$2 and agreement_id=$3 and status=$4 and version=$6 and ($5<>'active' or exists(select 1 from org.organization o where o.tenant_id=a.tenant_id and o.organization_id=a.franchise_organization_id and o.organization_type='franchisee' and o.status='active'))`, tenant, organization, id, current, target, version)
	if err != nil {
		return fulfillmentConstraint(err)
	}
	if result.RowsAffected() != 1 {
		return fulfillment.ErrConflict
	}
	if err = outbox(ctx, tx, tenant, eventID, "franchise-agreement", id, version+1, "franchise-agreement."+target, `'{}'::jsonb`); err != nil {
		return err
	}
	return nil
}
func (r *Fulfillment) QueueMessage(ctx context.Context, tenant, eventID string, v fulfillment.Message) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, `insert into communication.message(tenant_id,message_id,recipient_principal_id,channel,template_code,template_version,state)values($1,$2,$3,$4,$5,$6,'queued')`, tenant, v.ID, v.RecipientPrincipalID, v.Channel, v.TemplateCode, v.TemplateVersion)
	if err != nil {
		return err
	}
	if err = outbox(ctx, tx, tenant, eventID, "message", v.ID, 1, "message.queued", `jsonb_build_object('channel',$7::text)`, v.Channel); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
func (r *Fulfillment) TransitionMessage(ctx context.Context, tenant, id, current, target, reference, errorCode, eventID string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	result, err := tx.Exec(ctx, `update communication.message set state=$4,provider_reference=coalesce(nullif($5,''),provider_reference),last_error_code=nullif($6,''),updated_at=clock_timestamp() where tenant_id=$1 and message_id=$2 and state=$3`, tenant, id, current, target, reference, errorCode)
	if err != nil {
		return err
	}
	if result.RowsAffected() != 1 {
		return fulfillment.ErrConflict
	}
	if err = outbox(ctx, tx, tenant, eventID, "message", id, 1, "message."+target, `'{}'::jsonb`); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
func outbox(ctx context.Context, tx interface {
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
}, tenant, eventID, aggregateType, aggregateID string, version int64, eventType, payloadExpression string, payloadValues ...any) error {
	query := `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)values($1,$2,$3,$4,$5,$6,1,clock_timestamp(),` + payloadExpression + `)`
	args := []any{tenant, eventID, aggregateType, aggregateID, version, eventType}
	args = append(args, payloadValues...)
	_, err := tx.Exec(ctx, query, args...)
	return err
}

func fulfillmentConstraint(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && (pgErr.Code == "23503" || pgErr.Code == "23505" || pgErr.Code == "23514" || pgErr.Code == "23P01") {
		return fulfillment.ErrConflict
	}
	return err
}
