package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"elite.local/enterprise/internal/fiscal"
	"github.com/jackc/pgx/v5"
)

func (r *Fiscal) GetParameterSnapshot(ctx context.Context, tenant, organization, snapshotID string) (fiscal.ParameterSnapshot, error) {
	var value fiscal.ParameterSnapshot
	var codes []byte
	err := r.pool.QueryRow(ctx, `select tenant_id::text,snapshot_id,organization_id,taxpayer_cuit,parameter_kind,voucher_class,response_hash,provider_codes,fetched_at from fiscal.parameter_snapshot where tenant_id=$1 and organization_id=$2 and snapshot_id=$3`, tenant, organization, snapshotID).Scan(&value.TenantID, &value.ID, &value.OrganizationID, &value.TaxpayerCUIT, &value.Kind, &value.VoucherClass, &value.ResponseHash, &codes, &value.FetchedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return value, fiscal.ErrNotFound
	}
	if err != nil {
		return value, err
	}
	if err = json.Unmarshal(codes, &value.ProviderCodes); err != nil {
		return value, err
	}
	rows, err := r.pool.Query(ctx, `select parameter_code,description,valid_from::text,valid_until::text,voucher_class,emission_type,blocked,deregistered_on::text from fiscal.parameter_item where tenant_id=$1 and snapshot_id=$2 order by parameter_code`, tenant, snapshotID)
	if err != nil {
		return value, err
	}
	defer rows.Close()
	for rows.Next() {
		var item fiscal.ParameterItem
		if err = rows.Scan(&item.Code, &item.Description, &item.ValidFrom, &item.ValidUntil, &item.VoucherClass, &item.EmissionType, &item.Blocked, &item.DeregisteredOn); err != nil {
			return value, err
		}
		value.Items = append(value.Items, item)
	}
	return value, rows.Err()
}

func (r *Fiscal) StoreParameterSnapshot(ctx context.Context, value fiscal.ParameterSnapshot, eventID string) (fiscal.ParameterSnapshot, bool, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return value, false, err
	}
	defer tx.Rollback(ctx)
	class := any(value.VoucherClass)
	if _, err = tx.Exec(ctx, `select pg_advisory_xact_lock(hashtextextended($1::text||'|'||$2::text||'|'||$3::text||'|'||coalesce($4::text,''),0))`, value.TenantID, value.OrganizationID, value.Kind, class); err != nil {
		return value, false, err
	}
	var existing string
	err = tx.QueryRow(ctx, `select snapshot_id from fiscal.parameter_snapshot where tenant_id=$1 and organization_id=$2 and taxpayer_cuit=$3 and parameter_kind=$4 and voucher_class is not distinct from $5 and response_hash=$6`, value.TenantID, value.OrganizationID, value.TaxpayerCUIT, value.Kind, class, value.ResponseHash).Scan(&existing)
	if err == nil {
		value.ID = existing
		return value, true, tx.Commit(ctx)
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return value, false, err
	}
	codes, err := json.Marshal(append([]string{}, value.ProviderCodes...))
	if err != nil {
		return value, false, fiscal.ErrInvalid
	}
	result, err := tx.Exec(ctx, `insert into fiscal.parameter_snapshot(tenant_id,snapshot_id,organization_id,taxpayer_cuit,parameter_kind,voucher_class,response_hash,provider_codes,fetched_at) select $1,$2,$3,$4,$5,$6,$7,$8,$9 where exists(select 1 from org.organization where tenant_id=$1 and organization_id=$3 and status='active')`, value.TenantID, value.ID, value.OrganizationID, value.TaxpayerCUIT, value.Kind, class, value.ResponseHash, codes, value.FetchedAt)
	if err != nil {
		return value, false, fiscalConflict(err)
	}
	if result.RowsAffected() != 1 {
		return value, false, fiscal.ErrConflict
	}
	for _, item := range value.Items {
		_, err = tx.Exec(ctx, `insert into fiscal.parameter_item(tenant_id,snapshot_id,parameter_code,description,valid_from,valid_until,voucher_class,emission_type,blocked,deregistered_on) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`, value.TenantID, value.ID, item.Code, item.Description, item.ValidFrom, item.ValidUntil, item.VoucherClass, item.EmissionType, item.Blocked, item.DeregisteredOn)
		if err != nil {
			return value, false, fiscalConflict(err)
		}
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload) values($1,$2,'fiscal-parameter-snapshot',$3,1,'fiscal-parameter-snapshot.recorded',1,clock_timestamp(),jsonb_build_object('organization_id',$4::text,'taxpayer_cuit',$5::text,'parameter_kind',$6::text,'voucher_class',$7::text,'response_hash',$8::text,'item_count',$9::integer,'approved',false))`, value.TenantID, eventID, value.ID, value.OrganizationID, value.TaxpayerCUIT, value.Kind, class, value.ResponseHash, len(value.Items))
	if err != nil {
		return value, false, fiscalConflict(err)
	}
	return value, false, fiscalConflict(tx.Commit(ctx))
}

func (r *Fiscal) DecideParameterSnapshot(ctx context.Context, tenant, organization, snapshotID, subject string, approved bool, reason, eventID string) error {
	decision := "rejected"
	if approved {
		decision = "approved"
	}
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	result, err := tx.Exec(ctx, `insert into fiscal.parameter_decision(tenant_id,snapshot_id,decision,decided_by_subject,decision_reason) select s.tenant_id,s.snapshot_id,$4,$5,$6 from fiscal.parameter_snapshot s where s.tenant_id=$1 and s.organization_id=$2 and s.snapshot_id=$3`, tenant, organization, snapshotID, decision, subject, reason)
	if err != nil {
		return fiscalConflict(err)
	}
	if result.RowsAffected() != 1 {
		return fiscal.ErrNotFound
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload) values($1,$2,'fiscal-parameter-snapshot',$3,2,'fiscal-parameter-snapshot.decided',1,clock_timestamp(),jsonb_build_object('organization_id',$4::text,'decision',$5::text,'decided_by_subject',$6::text,'decision_reason',$7::text))`, tenant, eventID, snapshotID, organization, decision, subject, reason)
	if err != nil {
		return fiscalConflict(err)
	}
	return fiscalConflict(tx.Commit(ctx))
}

func (r *Fiscal) ConfigureParameterSchedule(ctx context.Context, value fiscal.ParameterSchedule, subject, eventID string) (fiscal.ParameterSchedule, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return value, err
	}
	defer tx.Rollback(ctx)
	class := any(value.VoucherClass)
	err = tx.QueryRow(ctx, `insert into fiscal.parameter_refresh_schedule(tenant_id,schedule_id,organization_id,taxpayer_cuit,parameter_kind,voucher_class,interval_seconds,next_run_at,active,version,configured_by_subject) select $1,$2,$3,$4,$5,$6,$7,$8,true,1,$9 where exists(select 1 from org.organization where tenant_id=$1 and organization_id=$3 and status='active') on conflict (tenant_id,organization_id,taxpayer_cuit,parameter_kind,voucher_class) do update set interval_seconds=excluded.interval_seconds,next_run_at=excluded.next_run_at,active=true,version=fiscal.parameter_refresh_schedule.version+1,configured_by_subject=excluded.configured_by_subject,configured_at=clock_timestamp(),lease_owner=null,lease_until=null returning schedule_id,version`, value.TenantID, value.ID, value.OrganizationID, value.TaxpayerCUIT, value.Kind, class, value.IntervalSeconds, value.NextRunAt, subject).Scan(&value.ID, &value.Version)
	if errors.Is(err, pgx.ErrNoRows) {
		return value, fiscal.ErrConflict
	}
	if err != nil {
		return value, fiscalConflict(err)
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload) values($1,$2,'fiscal-parameter-schedule',$3,$4,'fiscal-parameter-schedule.configured',1,clock_timestamp(),jsonb_build_object('organization_id',$5::text,'taxpayer_cuit',$6::text,'parameter_kind',$7::text,'voucher_class',$8::text,'interval_seconds',$9::bigint,'configured_by_subject',$10::text))`, value.TenantID, eventID, value.ID, value.Version, value.OrganizationID, value.TaxpayerCUIT, value.Kind, class, value.IntervalSeconds, subject)
	if err != nil {
		return value, fiscalConflict(err)
	}
	return value, fiscalConflict(tx.Commit(ctx))
}

func (r *Fiscal) ClaimParameterSchedule(ctx context.Context, worker string, lease time.Duration) (fiscal.ParameterSchedule, error) {
	var value fiscal.ParameterSchedule
	err := r.pool.QueryRow(ctx, `with candidate as (select tenant_id,schedule_id from fiscal.parameter_refresh_schedule where active and next_run_at<=clock_timestamp() and (lease_until is null or lease_until<clock_timestamp()) order by next_run_at,tenant_id,schedule_id for update skip locked limit 1) update fiscal.parameter_refresh_schedule s set lease_owner=$1,lease_until=clock_timestamp()+$2::interval,attempts=attempts+1 from candidate c where s.tenant_id=c.tenant_id and s.schedule_id=c.schedule_id returning s.tenant_id::text,s.schedule_id,s.organization_id,s.taxpayer_cuit,s.parameter_kind,s.voucher_class,s.interval_seconds,s.next_run_at,s.active,s.version,s.lease_owner`, worker, lease.String()).Scan(&value.TenantID, &value.ID, &value.OrganizationID, &value.TaxpayerCUIT, &value.Kind, &value.VoucherClass, &value.IntervalSeconds, &value.NextRunAt, &value.Active, &value.Version, &value.LeaseOwner)
	if errors.Is(err, pgx.ErrNoRows) {
		return value, fiscal.ErrNoParameterWork
	}
	return value, err
}

func (r *Fiscal) CompleteParameterSchedule(ctx context.Context, value fiscal.ParameterSchedule, worker, snapshotID, eventID string) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	result, err := tx.Exec(ctx, `update fiscal.parameter_refresh_schedule set next_run_at=clock_timestamp()+make_interval(secs=>interval_seconds::double precision),lease_owner=null,lease_until=null,last_snapshot_id=$4,last_succeeded_at=clock_timestamp(),last_error_code=null,version=version+1 where tenant_id=$1 and schedule_id=$2 and lease_owner=$3 and lease_until>=clock_timestamp()`, value.TenantID, value.ID, worker, snapshotID)
	if err != nil {
		return fiscalConflict(err)
	}
	if result.RowsAffected() != 1 {
		return fiscal.ErrConflict
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload) values($1,$2,'fiscal-parameter-schedule',$3,$4+1,'fiscal-parameter-schedule.refreshed',1,clock_timestamp(),jsonb_build_object('organization_id',$5::text,'snapshot_id',$6::text))`, value.TenantID, eventID, value.ID, value.Version, value.OrganizationID, snapshotID)
	if err != nil {
		return fiscalConflict(err)
	}
	return fiscalConflict(tx.Commit(ctx))
}

func (r *Fiscal) DeferParameterSchedule(ctx context.Context, value fiscal.ParameterSchedule, worker, code string, retry time.Duration) error {
	result, err := r.pool.Exec(ctx, `update fiscal.parameter_refresh_schedule set next_run_at=clock_timestamp()+$4::interval,lease_owner=null,lease_until=null,last_error_code=$5,version=version+1 where tenant_id=$1 and schedule_id=$2 and lease_owner=$3`, value.TenantID, value.ID, worker, retry.String(), code)
	if err != nil {
		return fiscalConflict(err)
	}
	if result.RowsAffected() != 1 {
		return fiscal.ErrConflict
	}
	return nil
}
