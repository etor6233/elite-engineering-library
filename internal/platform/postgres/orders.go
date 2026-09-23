package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"elite.local/enterprise/internal/order"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Orders struct{ pool *pgxpool.Pool }

func NewOrders(pool *pgxpool.Pool) *Orders { return &Orders{pool: pool} }

func (r *Orders) Create(ctx context.Context, entity order.Order, key, requestHash string) (order.Order, bool, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return order.Order{}, false, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	claim, err := tx.Exec(ctx, `insert into platform.idempotency_record(tenant_id,scope,idempotency_key,request_sha256_hex,status,locked_until,expires_at) values($1,'order:create',$2,$3,'processing',clock_timestamp()+interval '30 seconds',clock_timestamp()+interval '24 hours') on conflict (tenant_id,scope,idempotency_key) do nothing`, entity.TenantID, key, requestHash)
	if err != nil {
		return order.Order{}, false, err
	}
	if claim.RowsAffected() == 0 {
		var existingHash, status string
		var response []byte
		err = tx.QueryRow(ctx, `select request_sha256_hex,status,response_body from platform.idempotency_record where tenant_id=$1 and scope='order:create' and idempotency_key=$2 for update`, entity.TenantID, key).Scan(&existingHash, &status, &response)
		if err != nil {
			return order.Order{}, false, err
		}
		if existingHash != requestHash {
			return order.Order{}, false, fmt.Errorf("%w: idempotency key reused", order.ErrConflict)
		}
		if status != "completed" {
			return order.Order{}, false, fmt.Errorf("%w: request in progress", order.ErrConflict)
		}
		var replay order.Order
		if err := json.Unmarshal(response, &replay); err != nil {
			return order.Order{}, false, err
		}
		return replay, true, tx.Commit(ctx)
	}
	_, err = tx.Exec(ctx, `insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version) values($1,$2,$3,$4,$5,$6,$7,$8)`, entity.TenantID, entity.ID, entity.OrganizationID, entity.CustomerPrincipal, entity.State, entity.Currency, entity.TotalMinorUnits, entity.Version)
	if err != nil {
		return order.Order{}, false, err
	}
	payload, _ := json.Marshal(entity)
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload) values($1,gen_random_uuid(),'order',$2,$3,'order.created',1,clock_timestamp(),$4)`, entity.TenantID, entity.ID, entity.Version, payload)
	if err != nil {
		return order.Order{}, false, err
	}
	_, err = tx.Exec(ctx, `update platform.idempotency_record set status='completed',response_code=201,response_body=$3,resource_type='order',resource_id=$4,locked_until=null where tenant_id=$1 and scope='order:create' and idempotency_key=$2`, entity.TenantID, key, payload, entity.ID)
	if err != nil {
		return order.Order{}, false, err
	}
	return entity, false, tx.Commit(ctx)
}

func (r *Orders) Get(ctx context.Context, tenantID, organizationID, id string) (order.Order, error) {
	var entity order.Order
	err := r.pool.QueryRow(ctx, `select order_id,tenant_id,organization_id,customer_principal_id,state,currency,total_minor_units,version from sales.customer_order where tenant_id=$1 and organization_id=$2 and order_id=$3`, tenantID, organizationID, id).Scan(&entity.ID, &entity.TenantID, &entity.OrganizationID, &entity.CustomerPrincipal, &entity.State, &entity.Currency, &entity.TotalMinorUnits, &entity.Version)
	if errors.Is(err, pgx.ErrNoRows) {
		return order.Order{}, order.ErrNotFound
	}
	return entity, err
}

func (r *Orders) SaveTransition(ctx context.Context, entity order.Order, expectedVersion int64) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	result, err := tx.Exec(ctx, `update sales.customer_order set state=$4,version=$5,updated_at=clock_timestamp() where tenant_id=$1 and organization_id=$2 and order_id=$3 and version=$6`, entity.TenantID, entity.OrganizationID, entity.ID, entity.State, entity.Version, expectedVersion)
	if err != nil {
		return err
	}
	if result.RowsAffected() != 1 {
		return fmt.Errorf("%w: concurrent update", order.ErrConflict)
	}
	payload, _ := json.Marshal(entity)
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload) values($1,gen_random_uuid(),'order',$2,$3,$4,1,clock_timestamp(),$5)`, entity.TenantID, entity.ID, entity.Version, "order."+string(entity.State), payload)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}
