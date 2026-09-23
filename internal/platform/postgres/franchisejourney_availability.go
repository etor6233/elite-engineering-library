package postgres

import (
	"context"
	"errors"
	"time"

	"elite.local/enterprise/internal/franchisejourney"
	"github.com/jackc/pgx/v5"
)

func (r *FranchiseJourney) CreateAvailability(ctx context.Context, tenant, subject string, value franchisejourney.AvailabilityEntry, eventID string) (franchisejourney.AvailabilityEntry, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return value, err
	}
	defer tx.Rollback(ctx)
	value, err = r.createAvailabilityTx(ctx, tx, tenant, subject, value, eventID)
	if err != nil {
		return value, err
	}
	return value, tx.Commit(ctx)
}

func (r *FranchiseJourney) createAvailabilityTx(ctx context.Context, tx pgx.Tx, tenant, subject string, value franchisejourney.AvailabilityEntry, eventID string) (franchisejourney.AvailabilityEntry, error) {
	// Match the existing availability trigger's organization lock before taking
	// the statement snapshot. Booking uses this same transaction-scoped fence.
	if value.ResourceID == "" {
		if err := lockOrganizationSchedule(ctx, tx, tenant, value.OrganizationID); err != nil {
			return value, err
		}
	}
	err := tx.QueryRow(ctx, `insert into crm.availability_entry(tenant_id,availability_id,organization_id,resource_id,entry_type,reason_code,starts_at,ends_at,state,version,created_by_subject)
		select $1,$2,o.organization_id,nullif($4,''),$5,nullif($6,''),$7,$8,'active',1,$9 from org.organization o
		where o.tenant_id=$1 and o.organization_id=$3 and o.status='active' and ($4='' or exists(select 1 from crm.service_resource r where r.tenant_id=$1 and r.organization_id=$3 and r.resource_id=$4 and r.status='active'))
		returning availability_id,organization_id,coalesce(resource_id,''),entry_type,coalesce(reason_code,''),starts_at,ends_at,state,version`, tenant, value.ID, value.OrganizationID, value.ResourceID, value.EntryType, value.ReasonCode, value.StartsAt, value.EndsAt, subject).Scan(&value.ID, &value.OrganizationID, &value.ResourceID, &value.EntryType, &value.ReasonCode, &value.StartsAt, &value.EndsAt, &value.State, &value.Version)
	if errors.Is(err, pgx.ErrNoRows) || postgresConflict(err) {
		return value, franchisejourney.ErrConflict
	}
	if err != nil {
		return value, err
	}
	if err = writeJourneyOutbox(ctx, tx, tenant, eventID, "availability-entry", value.ID, 1, "availability-entry.created", map[string]any{"organization_id": value.OrganizationID, "resource_id": value.ResourceID, "entry_type": value.EntryType, "reason_code": value.ReasonCode, "starts_at": value.StartsAt, "ends_at": value.EndsAt, "actor_subject": subject}); err != nil {
		return value, err
	}
	return value, nil
}

func (r *FranchiseJourney) CancelAvailability(ctx context.Context, tenant, organization, entry string, version int64, subject, reason, eventID string) (franchisejourney.AvailabilityEntry, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return franchisejourney.AvailabilityEntry{}, err
	}
	defer tx.Rollback(ctx)
	var value franchisejourney.AvailabilityEntry
	if err = lockOrganizationSchedule(ctx, tx, tenant, organization); err != nil {
		return value, err
	}
	err = tx.QueryRow(ctx, `update crm.availability_entry set state='cancelled',version=version+1,cancelled_by_subject=$5,cancellation_reason_code=$6,cancelled_at=clock_timestamp() where tenant_id=$1 and organization_id=$2 and availability_id=$3 and version=$4 and state='active' returning availability_id,organization_id,coalesce(resource_id,''),entry_type,coalesce(reason_code,''),starts_at,ends_at,state,version`, tenant, organization, entry, version, subject, reason).Scan(&value.ID, &value.OrganizationID, &value.ResourceID, &value.EntryType, &value.ReasonCode, &value.StartsAt, &value.EndsAt, &value.State, &value.Version)
	if errors.Is(err, pgx.ErrNoRows) || postgresConflict(err) {
		return value, franchisejourney.ErrConflict
	}
	if err != nil {
		return value, err
	}
	if err = writeJourneyOutbox(ctx, tx, tenant, eventID, "availability-entry", entry, value.Version, "availability-entry.cancelled", map[string]any{"organization_id": organization, "actor_subject": subject, "reason_code": reason}); err != nil {
		return value, err
	}
	return value, tx.Commit(ctx)
}

func lockOrganizationSchedule(ctx context.Context, tx pgx.Tx, tenant, organization string) error {
	_, err := tx.Exec(ctx, `select pg_advisory_xact_lock(hashtextextended(($1::uuid)::text||'|availability|'||$2::text||'|organization',0))`, tenant, organization)
	return err
}

func (r *FranchiseJourney) Availability(ctx context.Context, tenant, organization, resource string, from, to time.Time) ([]franchisejourney.AvailabilityEntry, error) {
	rows, err := r.pool.Query(ctx, `select availability_id,organization_id,coalesce(resource_id,''),entry_type,coalesce(reason_code,''),starts_at,ends_at,state,version from crm.availability_entry where tenant_id=$1 and organization_id=$2 and ($3='' or resource_id=$3) and starts_at<$5 and ends_at>$4 order by starts_at,availability_id limit 1000`, tenant, organization, resource, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []franchisejourney.AvailabilityEntry{}
	for rows.Next() {
		var value franchisejourney.AvailabilityEntry
		if err = rows.Scan(&value.ID, &value.OrganizationID, &value.ResourceID, &value.EntryType, &value.ReasonCode, &value.StartsAt, &value.EndsAt, &value.State, &value.Version); err != nil {
			return nil, err
		}
		items = append(items, value)
	}
	return items, rows.Err()
}

func (r *FranchiseJourney) CancelCustomerAppointment(ctx context.Context, tenant, organization, customer, appointment string, version int64, reason, transitionID, eventID string) (franchisejourney.Appointment, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return franchisejourney.Appointment{}, err
	}
	defer tx.Rollback(ctx)
	var value franchisejourney.Appointment
	var current string
	err = tx.QueryRow(ctx, `with candidate as (select state from crm.appointment where tenant_id=$1 and organization_id=$2 and customer_principal_id=$3 and appointment_id=$4 and version=$5 and state in ('requested','confirmed') and starts_at>clock_timestamp() for update), changed as (update crm.appointment a set state='cancelled',version=version+1,updated_at=clock_timestamp() from candidate c where a.tenant_id=$1 and a.appointment_id=$4 returning a.appointment_id,a.organization_id,a.lead_id,coalesce(a.model_id,''),a.appointment_kind,a.starts_at,a.state,a.version,coalesce(a.slot_id,''),coalesce(a.ends_at,a.starts_at),c.state) select * from changed`, tenant, organization, customer, appointment, version).Scan(&value.ID, &value.OrganizationID, &value.LeadID, &value.ModelID, &value.Kind, &value.StartsAt, &value.State, &value.Version, &value.SlotID, &value.EndsAt, &current)
	if errors.Is(err, pgx.ErrNoRows) {
		return value, franchisejourney.ErrConflict
	}
	if err != nil {
		return value, err
	}
	if _, err = tx.Exec(ctx, `insert into crm.appointment_transition(tenant_id,transition_id,appointment_id,from_state,to_state,actor_subject,reason_code) values($1,$2,$3,$4,'cancelled',$5,$6)`, tenant, transitionID, appointment, current, customer, reason); err != nil {
		return value, err
	}
	if err = writeJourneyOutbox(ctx, tx, tenant, eventID, "appointment", appointment, value.Version, "appointment.cancelled", map[string]any{"organization_id": organization, "actor_subject": customer, "reason_code": reason, "channel": "customer"}); err != nil {
		return value, err
	}
	return value, tx.Commit(ctx)
}

func (r *FranchiseJourney) CreateAvailabilityOnce(ctx context.Context, tenant, subject, key, hash string, value franchisejourney.AvailabilityEntry, event string) (franchisejourney.AvailabilityEntry, bool, error) {
	var empty franchisejourney.AvailabilityEntry
	if subject == "" {
		return empty, false, franchisejourney.ErrInvalid
	}
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return empty, false, err
	}
	defer tx.Rollback(ctx)
	result, err := tx.Exec(ctx, `insert into platform.idempotency_record(tenant_id,scope,idempotency_key,request_sha256_hex,status,locked_until,expires_at)values($1,'franchise-availability',$2,$3,'processing',clock_timestamp()+interval '30 seconds',clock_timestamp()+interval '24 hours')on conflict do nothing`, tenant, key, hash)
	if err != nil {
		return empty, false, err
	}
	if result.RowsAffected() == 0 {
		replay, storedHash, err := readAvailabilityCreation(ctx, tx, tenant, value.OrganizationID, key)
		if err != nil {
			return empty, false, err
		}
		if hash != storedHash {
			return empty, false, franchisejourney.ErrConflict
		}
		if err = tx.Commit(ctx); err != nil {
			return empty, false, err
		}
		return replay, true, nil
	}
	value, err = r.createAvailabilityTx(ctx, tx, tenant, subject, value, event)
	if err != nil {
		return empty, false, err
	}
	result, err = tx.Exec(ctx, `update platform.idempotency_record set status='completed',response_code=201,response_body=jsonb_build_object('availability_id',$3::text),resource_type='availability-entry',resource_id=$3,locked_until=null where tenant_id=$1 and scope='franchise-availability' and idempotency_key=$2 and status='processing'`, tenant, key, value.ID)
	if err != nil {
		return empty, false, err
	}
	if result.RowsAffected() != 1 {
		return empty, false, franchisejourney.ErrConflict
	}
	if err = tx.Commit(ctx); err != nil {
		return empty, false, err
	}
	return value, false, nil
}

type availabilityRowReader interface {
	QueryRow(context.Context, string, ...any) pgx.Row
}

func readAvailabilityCreation(ctx context.Context, q availabilityRowReader, tenant, organization, key string) (franchisejourney.AvailabilityEntry, string, error) {
	var v franchisejourney.AvailabilityEntry
	var hash string
	err := q.QueryRow(ctx, `select i.request_sha256_hex,a.availability_id,a.organization_id,coalesce(a.resource_id,''),a.entry_type,coalesce(a.reason_code,''),a.starts_at,a.ends_at,a.state,a.version
 from platform.idempotency_record i join crm.availability_entry a on a.tenant_id=i.tenant_id and a.availability_id=i.resource_id
 where i.tenant_id=$1 and i.scope='franchise-availability' and i.idempotency_key=$3 and i.status='completed' and i.resource_type='availability-entry' and i.response_code=201 and i.response_body->>'availability_id'=a.availability_id and a.organization_id=$2`, tenant, organization, key).Scan(&hash, &v.ID, &v.OrganizationID, &v.ResourceID, &v.EntryType, &v.ReasonCode, &v.StartsAt, &v.EndsAt, &v.State, &v.Version)
	if errors.Is(err, pgx.ErrNoRows) {
		return franchisejourney.AvailabilityEntry{}, "", franchisejourney.ErrConflict
	}
	if err != nil {
		return franchisejourney.AvailabilityEntry{}, "", err
	}
	return v, hash, nil
}
func (r *FranchiseJourney) AvailabilityCreationResult(ctx context.Context, tenant, organization, key string) (franchisejourney.AvailabilityEntry, error) {
	v, _, err := readAvailabilityCreation(ctx, r.pool, tenant, organization, key)
	return v, err
}
