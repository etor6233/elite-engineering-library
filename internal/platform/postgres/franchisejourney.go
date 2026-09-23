package postgres

import (
	"context"
	"elite.local/enterprise/internal/businesspolicy"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"elite.local/enterprise/internal/bcsales"
	"elite.local/enterprise/internal/franchisejourney"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FranchiseJourney struct {
	pool   *pgxpool.Pool
	policy *businesspolicy.Profile
}

func NewFranchiseJourney(pool *pgxpool.Pool) *FranchiseJourney {
	return &FranchiseJourney{pool: pool, policy: businesspolicy.Reference()}
}

func NewFranchiseJourneyWithProfile(pool *pgxpool.Pool, policy *businesspolicy.Profile) (*FranchiseJourney, error) {
	if pool == nil || !policy.Valid() {
		return nil, businesspolicy.ErrProfile
	}
	return &FranchiseJourney{pool: pool, policy: policy}, nil
}
func (r *FranchiseJourney) BusinessPolicySHA256() string {
	if r == nil {
		return ""
	}
	return r.policy.SHA256()
}

// Read model only; assignment/transition owners remain authoritative for writes.
func (r *FranchiseJourney) AppointmentAgenda(ctx context.Context, tenant, organization string, from, to time.Time) (franchisejourney.AppointmentAgenda, error) {
	result := franchisejourney.AppointmentAgenda{Appointments: []franchisejourney.Appointment{}, Resources: []franchisejourney.ServiceResource{}}
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead, AccessMode: pgx.ReadOnly})
	if err != nil {
		return result, err
	}
	defer tx.Rollback(ctx)
	rows, err := tx.Query(ctx, `select a.appointment_id,a.organization_id,a.lead_id,coalesce(a.model_id,''),a.appointment_kind,a.starts_at,a.state,a.version,coalesce(a.slot_id,''),coalesce(a.ends_at,a.starts_at),coalesce(ar.resource_id,'')
 from crm.appointment a left join crm.appointment_resource ar on ar.tenant_id=a.tenant_id and ar.appointment_id=a.appointment_id
 where a.tenant_id=$1 and a.organization_id=$2 and a.starts_at >= $3 and a.starts_at < $4
 order by a.starts_at,a.appointment_id limit 201`, tenant, organization, from, to)
	if err != nil {
		return result, err
	}
	for rows.Next() {
		var a franchisejourney.Appointment
		if err := rows.Scan(&a.ID, &a.OrganizationID, &a.LeadID, &a.ModelID, &a.Kind, &a.StartsAt, &a.State, &a.Version, &a.SlotID, &a.EndsAt, &a.ResourceID); err != nil {
			rows.Close()
			return result, err
		}
		result.Appointments = append(result.Appointments, a)
	}
	err = rows.Err()
	rows.Close()
	if err != nil {
		return result, err
	}
	rows, err = tx.Query(ctx, `select r.resource_id,r.organization_id,r.display_name,r.resource_kind,r.status,r.version,
 coalesce(array(select s.appointment_kind from crm.resource_skill s where s.tenant_id=r.tenant_id and s.resource_id=r.resource_id order by s.appointment_kind),'{}')
 from crm.service_resource r where r.tenant_id=$1 and r.organization_id=$2 and r.status='active'
 order by r.display_name,r.resource_id limit 201`, tenant, organization)
	if err != nil {
		return result, err
	}
	for rows.Next() {
		var item franchisejourney.ServiceResource
		if err := rows.Scan(&item.ID, &item.OrganizationID, &item.DisplayName, &item.Kind, &item.Status, &item.Version, &item.Skills); err != nil {
			rows.Close()
			return result, err
		}
		result.Resources = append(result.Resources, item)
	}
	err = rows.Err()
	rows.Close()
	if err != nil {
		return result, err
	}
	if len(result.Appointments) > 200 {
		result.Truncated = true
		result.Appointments = result.Appointments[:200]
	}
	if len(result.Resources) > 200 {
		result.Truncated = true
		result.Resources = result.Resources[:200]
	}
	return result, tx.Commit(ctx)
}

func (r *FranchiseJourney) PublicLocations(ctx context.Context, tenantCode string) ([]franchisejourney.Location, error) {
	rows, err := r.pool.Query(ctx, `select o.organization_id,o.organization_code,o.display_name,l.city,l.region,l.country,coalesce(l.contact_phone,''),coalesce(l.contact_email,'') from platform.tenant t join org.organization o on o.tenant_id=t.tenant_id join org.public_location l on l.tenant_id=o.tenant_id and l.organization_id=o.organization_id where t.tenant_code=$1 and t.status='active' and o.status='active' and l.published=true order by l.sort_order,o.organization_id`, tenantCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []franchisejourney.Location{}
	for rows.Next() {
		var item franchisejourney.Location
		if err := rows.Scan(&item.OrganizationID, &item.Code, &item.Name, &item.City, &item.Region, &item.Country, &item.ContactPhone, &item.ContactEmail); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *FranchiseJourney) PublicAppointmentSlots(ctx context.Context, tenantCode, organizationCode, kind string, from, to time.Time) ([]franchisejourney.AppointmentSlot, error) {
	rows, err := r.pool.Query(ctx, `select s.slot_id,s.organization_id,s.appointment_kind,s.starts_at,s.ends_at,s.capacity,count(a.appointment_id),s.state,s.version
		from platform.tenant t
		join org.organization o on o.tenant_id=t.tenant_id
		join org.public_location l on l.tenant_id=o.tenant_id and l.organization_id=o.organization_id and l.published=true
		join crm.appointment_slot s on s.tenant_id=o.tenant_id and s.organization_id=o.organization_id
		left join crm.appointment a on a.tenant_id=s.tenant_id and a.slot_id=s.slot_id and a.state in ('requested','confirmed')
		where t.tenant_code=$1 and t.status='active' and o.organization_code=$2 and o.status='active'
		  and s.appointment_kind=$3 and s.state='open' and s.starts_at >= $4 and s.starts_at < $5
		  and s.starts_at >= clock_timestamp()+make_interval(secs=>$6::double precision)
		  and exists(select 1 from crm.availability_entry e where e.tenant_id=s.tenant_id and e.organization_id=s.organization_id and e.resource_id is null and e.entry_type='working' and e.state='active' and e.starts_at<=s.starts_at and e.ends_at>=s.ends_at)
		  and not exists(select 1 from crm.availability_entry e where e.tenant_id=s.tenant_id and e.organization_id=s.organization_id and e.resource_id is null and e.entry_type='unavailable' and e.state='active' and tstzrange(e.starts_at,e.ends_at,'[)') && tstzrange(s.starts_at,s.ends_at,'[)'))
		group by s.tenant_id,s.slot_id
		having count(a.appointment_id) < s.capacity
		order by s.starts_at,s.slot_id limit 500`, tenantCode, organizationCode, kind, from, to, r.policy.LeadTime().Seconds())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []franchisejourney.AppointmentSlot{}
	for rows.Next() {
		var item franchisejourney.AppointmentSlot
		if err := rows.Scan(&item.ID, &item.OrganizationID, &item.Kind, &item.StartsAt, &item.EndsAt, &item.Capacity, &item.Booked, &item.State, &item.Version); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *FranchiseJourney) CreateAppointmentSlot(ctx context.Context, tenant string, value franchisejourney.AppointmentSlot, eventID string) (franchisejourney.AppointmentSlot, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return franchisejourney.AppointmentSlot{}, err
	}
	defer tx.Rollback(ctx)
	value, err = r.createAppointmentSlotTx(ctx, tx, tenant, "", value, eventID)
	if err != nil {
		return franchisejourney.AppointmentSlot{}, err
	}
	return value, tx.Commit(ctx)
}

func (r *FranchiseJourney) createAppointmentSlotTx(ctx context.Context, tx pgx.Tx, tenant, subject string, value franchisejourney.AppointmentSlot, eventID string) (franchisejourney.AppointmentSlot, error) {
	// Direct repository callers receive the same profile bounds as the service.
	// The persisted slot retains its own capacity; later profiles are prospective.
	var now time.Time
	if err := tx.QueryRow(ctx, `select clock_timestamp()`).Scan(&now); err != nil {
		return franchisejourney.AppointmentSlot{}, err
	}
	if !r.policy.AllowsSlot(value.StartsAt, value.EndsAt, now, value.Capacity) {
		return franchisejourney.AppointmentSlot{}, franchisejourney.ErrConflict
	}
	err := tx.QueryRow(ctx, `insert into crm.appointment_slot(tenant_id,slot_id,organization_id,appointment_kind,starts_at,ends_at,capacity,state,version)
		select $1,$2,o.organization_id,$4,$5,$6,$7,'open',1 from org.organization o
		where o.tenant_id=$1 and o.organization_id=$3 and o.status='active'
		returning slot_id,organization_id,appointment_kind,starts_at,ends_at,capacity,0,state,version`, tenant, value.ID, value.OrganizationID, value.Kind, value.StartsAt, value.EndsAt, value.Capacity).Scan(&value.ID, &value.OrganizationID, &value.Kind, &value.StartsAt, &value.EndsAt, &value.Capacity, &value.Booked, &value.State, &value.Version)
	if errors.Is(err, pgx.ErrNoRows) || postgresConflict(err) {
		return franchisejourney.AppointmentSlot{}, franchisejourney.ErrConflict
	}
	if err != nil {
		return franchisejourney.AppointmentSlot{}, err
	}
	payload := map[string]any{"organization_id": value.OrganizationID, "appointment_kind": value.Kind, "starts_at": value.StartsAt, "capacity": value.Capacity, "policy_sha256": r.policy.SHA256()}
	if subject != "" {
		payload["actor_subject"] = subject
	}
	if err = writeJourneyOutbox(ctx, tx, tenant, eventID, "appointment-slot", value.ID, 1, "appointment-slot.created", payload); err != nil {
		return franchisejourney.AppointmentSlot{}, err
	}
	return value, nil
}

func (r *FranchiseJourney) CreateServiceResource(ctx context.Context, tenant string, value franchisejourney.ServiceResource, eventID string) (franchisejourney.ServiceResource, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return franchisejourney.ServiceResource{}, err
	}
	defer tx.Rollback(ctx)
	value, err = r.createServiceResourceTx(ctx, tx, tenant, "", value, eventID)
	if err != nil {
		return franchisejourney.ServiceResource{}, err
	}
	return value, tx.Commit(ctx)
}

func (r *FranchiseJourney) createServiceResourceTx(ctx context.Context, tx pgx.Tx, tenant, subject string, value franchisejourney.ServiceResource, eventID string) (franchisejourney.ServiceResource, error) {
	result, err := tx.Exec(ctx, `insert into crm.service_resource(tenant_id,resource_id,organization_id,principal_subject,display_name,resource_kind,status,version)
		select $1,$2,o.organization_id,nullif($4,''),$5,$6,'active',1 from org.organization o where o.tenant_id=$1 and o.organization_id=$3 and o.status='active'`, tenant, value.ID, value.OrganizationID, value.PrincipalSubject, value.DisplayName, value.Kind)
	if err != nil {
		if postgresConflict(err) {
			return franchisejourney.ServiceResource{}, franchisejourney.ErrConflict
		}
		return franchisejourney.ServiceResource{}, err
	}
	if result.RowsAffected() != 1 {
		return franchisejourney.ServiceResource{}, franchisejourney.ErrConflict
	}
	if _, err = tx.Exec(ctx, `insert into crm.resource_skill(tenant_id,resource_id,appointment_kind) select $1,$2,unnest($3::text[])`, tenant, value.ID, value.Skills); err != nil {
		if postgresConflict(err) {
			return franchisejourney.ServiceResource{}, franchisejourney.ErrConflict
		}
		return franchisejourney.ServiceResource{}, err
	}
	payload := map[string]any{"organization_id": value.OrganizationID, "resource_kind": value.Kind, "skills": value.Skills}
	if subject != "" {
		payload["actor_subject"] = subject
	}
	if err = writeJourneyOutbox(ctx, tx, tenant, eventID, "service-resource", value.ID, 1, "service-resource.created", payload); err != nil {
		return franchisejourney.ServiceResource{}, err
	}
	return value, nil
}

func (r *FranchiseJourney) AssignAppointmentResource(ctx context.Context, tenant, organization, appointment, resource string, version int64, eventID string) (franchisejourney.Appointment, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return franchisejourney.Appointment{}, err
	}
	defer tx.Rollback(ctx)
	if _, err = tx.Exec(ctx, `insert into crm.appointment_resource(tenant_id,appointment_id,resource_id) values($1,$2,$3)`, tenant, appointment, resource); err != nil {
		if postgresConflict(err) {
			return franchisejourney.Appointment{}, franchisejourney.ErrConflict
		}
		return franchisejourney.Appointment{}, err
	}
	var value franchisejourney.Appointment
	err = tx.QueryRow(ctx, `update crm.appointment set version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and organization_id=$2 and appointment_id=$3 and state='requested' and version=$4 returning appointment_id,organization_id,lead_id,coalesce(model_id,''),appointment_kind,starts_at,state,version,coalesce(slot_id,''),coalesce(ends_at,starts_at)`, tenant, organization, appointment, version).Scan(&value.ID, &value.OrganizationID, &value.LeadID, &value.ModelID, &value.Kind, &value.StartsAt, &value.State, &value.Version, &value.SlotID, &value.EndsAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return franchisejourney.Appointment{}, franchisejourney.ErrConflict
	}
	if err != nil {
		return franchisejourney.Appointment{}, err
	}
	value.ResourceID = resource
	if err = writeJourneyOutbox(ctx, tx, tenant, eventID, "appointment", appointment, value.Version, "appointment.resource-assigned", map[string]any{"organization_id": organization, "resource_id": resource}); err != nil {
		return franchisejourney.Appointment{}, err
	}
	return value, tx.Commit(ctx)
}

func (r *FranchiseJourney) TransitionAppointment(ctx context.Context, tenant, organization, appointment, current, target string, version int64, subject, reason, eventID string) (franchisejourney.Appointment, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return franchisejourney.Appointment{}, err
	}
	defer tx.Rollback(ctx)
	var value franchisejourney.Appointment
	err = tx.QueryRow(ctx, `update crm.appointment a set state=$5,version=version+1,updated_at=clock_timestamp()
		where tenant_id=$1 and organization_id=$2 and appointment_id=$3 and state=$4 and version=$6
		and ($5<>'confirmed' or exists(select 1 from crm.appointment_resource ar where ar.tenant_id=a.tenant_id and ar.appointment_id=a.appointment_id))
		and ($5 not in ('completed','no-show') or a.starts_at<=clock_timestamp())
		returning appointment_id,organization_id,lead_id,coalesce(model_id,''),appointment_kind,starts_at,state,version,coalesce(slot_id,''),coalesce(ends_at,starts_at),coalesce((select ar.resource_id from crm.appointment_resource ar where ar.tenant_id=a.tenant_id and ar.appointment_id=a.appointment_id),'')`, tenant, organization, appointment, current, target, version).Scan(&value.ID, &value.OrganizationID, &value.LeadID, &value.ModelID, &value.Kind, &value.StartsAt, &value.State, &value.Version, &value.SlotID, &value.EndsAt, &value.ResourceID)
	if errors.Is(err, pgx.ErrNoRows) {
		return franchisejourney.Appointment{}, franchisejourney.ErrConflict
	}
	if err != nil {
		return franchisejourney.Appointment{}, err
	}
	if _, err = tx.Exec(ctx, `insert into crm.appointment_transition(tenant_id,transition_id,appointment_id,from_state,to_state,actor_subject,reason_code) values($1,$2,$3,$4,$5,$6,nullif($7,''))`, tenant, eventID, appointment, current, target, subject, reason); err != nil {
		return franchisejourney.Appointment{}, err
	}
	if err = writeJourneyOutbox(ctx, tx, tenant, eventID, "appointment", appointment, value.Version, "appointment."+target, map[string]any{"organization_id": organization, "resource_id": value.ResourceID, "actor_subject": subject, "reason_code": reason}); err != nil {
		return franchisejourney.Appointment{}, err
	}
	return value, tx.Commit(ctx)
}

func (r *FranchiseJourney) RequestAppointment(ctx context.Context, tenantCode, organizationCode, idempotencyKey string, value franchisejourney.Appointment, requestHash, eventID string) (franchisejourney.Appointment, bool, error) {
	boundHash, bindErr := r.policy.BindRequestHash(requestHash)
	if bindErr != nil {
		return franchisejourney.Appointment{}, false, franchisejourney.ErrInvalid
	}
	requestHash = boundHash
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return franchisejourney.Appointment{}, false, err
	}
	defer tx.Rollback(ctx)
	var tenant, organization string
	err = tx.QueryRow(ctx, `select t.tenant_id,o.organization_id from platform.tenant t join org.organization o on o.tenant_id=t.tenant_id join org.public_location l on l.tenant_id=o.tenant_id and l.organization_id=o.organization_id where t.tenant_code=$1 and t.status='active' and o.organization_code=$2 and o.status='active' and l.published=true`, tenantCode, organizationCode).Scan(&tenant, &organization)
	if errors.Is(err, pgx.ErrNoRows) {
		return franchisejourney.Appointment{}, false, franchisejourney.ErrNotFound
	}
	if err != nil {
		return franchisejourney.Appointment{}, false, err
	}
	result, err := tx.Exec(ctx, `insert into platform.idempotency_record(tenant_id,scope,idempotency_key,request_sha256_hex,status,locked_until,expires_at) values($1,'public-appointment',$2,$3,'processing',clock_timestamp()+interval '30 seconds',clock_timestamp()+interval '24 hours') on conflict do nothing`, tenant, idempotencyKey, requestHash)
	if err != nil {
		return franchisejourney.Appointment{}, false, err
	}
	if result.RowsAffected() == 0 {
		var storedHash, status, resourceID string
		if err = tx.QueryRow(ctx, `select request_sha256_hex,status,coalesce(resource_id,'') from platform.idempotency_record where tenant_id=$1 and scope='public-appointment' and idempotency_key=$2`, tenant, idempotencyKey).Scan(&storedHash, &status, &resourceID); err != nil {
			return franchisejourney.Appointment{}, false, err
		}
		if storedHash != requestHash || status != "completed" || resourceID == "" {
			return franchisejourney.Appointment{}, false, franchisejourney.ErrConflict
		}
		var replay franchisejourney.Appointment
		err = tx.QueryRow(ctx, `select appointment_id,organization_id,lead_id,coalesce(model_id,''),appointment_kind,starts_at,state,version,coalesce(slot_id,''),coalesce(ends_at,starts_at) from crm.appointment where tenant_id=$1 and appointment_id=$2`, tenant, resourceID).Scan(&replay.ID, &replay.OrganizationID, &replay.LeadID, &replay.ModelID, &replay.Kind, &replay.StartsAt, &replay.State, &replay.Version, &replay.SlotID, &replay.EndsAt)
		if err != nil {
			return franchisejourney.Appointment{}, false, err
		}
		if err = tx.Commit(ctx); err != nil {
			return franchisejourney.Appointment{}, false, err
		}
		return replay, true, nil
	}
	value.OrganizationID = organization
	if err = lockOrganizationSchedule(ctx, tx, tenant, organization); err != nil {
		return franchisejourney.Appointment{}, false, err
	}
	var capacity int
	err = tx.QueryRow(ctx, `select s.slot_id,s.ends_at,s.capacity from crm.appointment_slot s where s.tenant_id=$1 and s.organization_id=$2 and s.appointment_kind=$3 and s.starts_at=$4 and s.state='open' and s.starts_at>=clock_timestamp()+make_interval(secs=>$5::double precision)
		and exists(select 1 from crm.availability_entry e where e.tenant_id=s.tenant_id and e.organization_id=s.organization_id and e.resource_id is null and e.entry_type='working' and e.state='active' and e.starts_at<=s.starts_at and e.ends_at>=s.ends_at)
		and not exists(select 1 from crm.availability_entry e where e.tenant_id=s.tenant_id and e.organization_id=s.organization_id and e.resource_id is null and e.entry_type='unavailable' and e.state='active' and tstzrange(e.starts_at,e.ends_at,'[)') && tstzrange(s.starts_at,s.ends_at,'[)')) for update of s`, tenant, organization, value.Kind, value.StartsAt, r.policy.LeadTime().Seconds()).Scan(&value.SlotID, &value.EndsAt, &capacity)
	if errors.Is(err, pgx.ErrNoRows) {
		return franchisejourney.Appointment{}, false, franchisejourney.ErrConflict
	}
	if err != nil {
		return franchisejourney.Appointment{}, false, err
	}
	var booked int
	if err = tx.QueryRow(ctx, `select count(*) from crm.appointment where tenant_id=$1 and slot_id=$2 and state in ('requested','confirmed')`, tenant, value.SlotID).Scan(&booked); err != nil {
		return franchisejourney.Appointment{}, false, err
	}
	if booked >= capacity {
		return franchisejourney.Appointment{}, false, franchisejourney.ErrConflict
	}
	result, err = tx.Exec(ctx, `insert into crm.appointment(tenant_id,appointment_id,organization_id,lead_id,customer_principal_id,model_id,appointment_kind,starts_at,ends_at,slot_id,state,version) select $1,$2,$3,l.lead_id,l.customer_principal_id,nullif($5,''),$6,$7,$8,$9,'requested',1 from crm.lead l where l.tenant_id=$1 and l.organization_id=$3 and l.lead_id=$4`, tenant, value.ID, organization, value.LeadID, value.ModelID, value.Kind, value.StartsAt, value.EndsAt, value.SlotID)
	if err != nil {
		return franchisejourney.Appointment{}, false, err
	}
	if result.RowsAffected() != 1 {
		return franchisejourney.Appointment{}, false, franchisejourney.ErrConflict
	}
	if err = writeJourneyOutbox(ctx, tx, tenant, eventID, "appointment", value.ID, 1, "appointment.requested", map[string]any{"organization_id": organization, "lead_id": value.LeadID, "policy_sha256": r.policy.SHA256()}); err != nil {
		return franchisejourney.Appointment{}, false, err
	}
	_, err = tx.Exec(ctx, `update platform.idempotency_record set status='completed',response_code=202,response_body=jsonb_build_object('appointment_id',$3::text),resource_type='appointment',resource_id=$3,locked_until=null where tenant_id=$1 and scope='public-appointment' and idempotency_key=$2 and status='processing'`, tenant, idempotencyKey, value.ID)
	if err != nil {
		return franchisejourney.Appointment{}, false, err
	}
	if err = tx.Commit(ctx); err != nil {
		return franchisejourney.Appointment{}, false, err
	}
	return value, false, nil
}

func (r *FranchiseJourney) Leads(ctx context.Context, tenant, organization string, limit int, after string) (franchisejourney.Page[franchisejourney.Lead], error) {
	rows, err := r.pool.Query(ctx, `select lead_id,organization_id,coalesce(model_id,''),lifecycle_state,source_code,coalesce(assigned_subject,''),created_at,version from crm.lead where tenant_id=$1 and organization_id=$2 and ($3='' or lead_id>$3) order by lead_id limit $4`, tenant, organization, after, limit+1)
	if err != nil {
		return franchisejourney.Page[franchisejourney.Lead]{}, err
	}
	defer rows.Close()
	items := []franchisejourney.Lead{}
	for rows.Next() {
		var item franchisejourney.Lead
		if err := rows.Scan(&item.ID, &item.OrganizationID, &item.ModelID, &item.State, &item.SourceCode, &item.AssignedSubject, &item.CreatedAt, &item.Version); err != nil {
			return franchisejourney.Page[franchisejourney.Lead]{}, err
		}
		items = append(items, item)
	}
	page := franchisejourney.Page[franchisejourney.Lead]{Items: items}
	if len(items) > limit {
		page.NextCursor = items[limit-1].ID
		page.Items = items[:limit]
	}
	return page, rows.Err()
}

func (r *FranchiseJourney) AssignLead(ctx context.Context, tenant, organization, lead, subject string, version int64, eventID string) (franchisejourney.Lead, error) {
	return r.assignLead(ctx, tenant, organization, lead, subject, version, eventID, "")
}

func (r *FranchiseJourney) AssignLeadAs(ctx context.Context, tenant, organization, lead, subject string, version int64, eventID, actor string) (franchisejourney.Lead, error) {
	if actor == "" {
		return franchisejourney.Lead{}, franchisejourney.ErrInvalid
	}
	return r.assignLead(ctx, tenant, organization, lead, subject, version, eventID, actor)
}

func (r *FranchiseJourney) assignLead(ctx context.Context, tenant, organization, lead, subject string, version int64, eventID, actor string) (franchisejourney.Lead, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return franchisejourney.Lead{}, err
	}
	defer tx.Rollback(ctx)
	var value franchisejourney.Lead
	err = tx.QueryRow(ctx, `update crm.lead set assigned_subject=$5,version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and organization_id=$2 and lead_id=$3 and version=$4 and lifecycle_state not in ('converted','lost') returning lead_id,organization_id,coalesce(model_id,''),lifecycle_state,source_code,assigned_subject,created_at,version`, tenant, organization, lead, version, subject).Scan(&value.ID, &value.OrganizationID, &value.ModelID, &value.State, &value.SourceCode, &value.AssignedSubject, &value.CreatedAt, &value.Version)
	if errors.Is(err, pgx.ErrNoRows) {
		return franchisejourney.Lead{}, franchisejourney.ErrConflict
	}
	if err != nil {
		return franchisejourney.Lead{}, err
	}
	payload := map[string]any{"organization_id": organization, "assigned_subject": subject}
	if actor != "" {
		payload["actor_subject"] = actor
	}
	if err = writeJourneyOutbox(ctx, tx, tenant, eventID, "lead", lead, value.Version, "lead.assigned", payload); err != nil {
		return franchisejourney.Lead{}, err
	}
	return value, tx.Commit(ctx)
}

func (r *FranchiseJourney) TransitionLead(ctx context.Context, tenant, organization, lead, current, target string, version int64, eventID string) (franchisejourney.Lead, error) {
	return r.transitionLead(ctx, tenant, organization, lead, current, target, version, eventID, "")
}

func (r *FranchiseJourney) TransitionLeadAs(ctx context.Context, tenant, organization, lead, current, target string, version int64, eventID, actor string) (franchisejourney.Lead, error) {
	if actor == "" {
		return franchisejourney.Lead{}, franchisejourney.ErrInvalid
	}
	return r.transitionLead(ctx, tenant, organization, lead, current, target, version, eventID, actor)
}

func (r *FranchiseJourney) transitionLead(ctx context.Context, tenant, organization, lead, current, target string, version int64, eventID, actor string) (franchisejourney.Lead, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return franchisejourney.Lead{}, err
	}
	defer tx.Rollback(ctx)
	var value franchisejourney.Lead
	err = tx.QueryRow(ctx, `update crm.lead set lifecycle_state=$6,version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and organization_id=$2 and lead_id=$3 and lifecycle_state=$4 and version=$5 returning lead_id,organization_id,coalesce(model_id,''),lifecycle_state,source_code,coalesce(assigned_subject,''),created_at,version`, tenant, organization, lead, current, version, target).Scan(&value.ID, &value.OrganizationID, &value.ModelID, &value.State, &value.SourceCode, &value.AssignedSubject, &value.CreatedAt, &value.Version)
	if errors.Is(err, pgx.ErrNoRows) {
		return franchisejourney.Lead{}, franchisejourney.ErrConflict
	}
	if err != nil {
		return franchisejourney.Lead{}, err
	}
	payload := map[string]any{"organization_id": organization}
	if actor != "" {
		payload["actor_subject"] = actor
	}
	if err = writeJourneyOutbox(ctx, tx, tenant, eventID, "lead", lead, value.Version, "lead."+target, payload); err != nil {
		return franchisejourney.Lead{}, err
	}
	return value, tx.Commit(ctx)
}

func (r *FranchiseJourney) CreateQuote(ctx context.Context, tenant, idempotencyKey string, value franchisejourney.Quote, requestHash, eventID string) (franchisejourney.Quote, bool, error) {
	return r.createQuote(ctx, tenant, idempotencyKey, value, requestHash, eventID, "")
}
func (r *FranchiseJourney) CreateQuoteAs(ctx context.Context, tenant, key string, value franchisejourney.Quote, hash, event, actor string) (franchisejourney.Quote, bool, error) {
	if actor == "" {
		return franchisejourney.Quote{}, false, franchisejourney.ErrInvalid
	}
	return r.createQuote(ctx, tenant, key, value, hash, event, actor)
}
func (r *FranchiseJourney) createQuote(ctx context.Context, tenant, idempotencyKey string, value franchisejourney.Quote, requestHash, eventID, actor string) (franchisejourney.Quote, bool, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return franchisejourney.Quote{}, false, err
	}
	defer tx.Rollback(ctx)
	result, err := tx.Exec(ctx, `insert into platform.idempotency_record(tenant_id,scope,idempotency_key,request_sha256_hex,status,locked_until,expires_at) values($1,'franchise-quote',$2,$3,'processing',clock_timestamp()+interval '30 seconds',clock_timestamp()+interval '24 hours') on conflict do nothing`, tenant, idempotencyKey, requestHash)
	if err != nil {
		return franchisejourney.Quote{}, false, err
	}
	if result.RowsAffected() == 0 {
		var storedHash, status, resourceID string
		if err = tx.QueryRow(ctx, `select request_sha256_hex,status,coalesce(resource_id,'') from platform.idempotency_record where tenant_id=$1 and scope='franchise-quote' and idempotency_key=$2`, tenant, idempotencyKey).Scan(&storedHash, &status, &resourceID); err != nil {
			return franchisejourney.Quote{}, false, err
		}
		if storedHash != requestHash || status != "completed" || resourceID == "" {
			return franchisejourney.Quote{}, false, franchisejourney.ErrConflict
		}
		var replay franchisejourney.Quote
		err = tx.QueryRow(ctx, `select quotation_id,organization_id,lead_id,coalesce(customer_principal_id,''),variant_id,price_book_id,currency,total_minor_units,valid_until,state,version,coalesce(order_id,'') from sales.quotation where tenant_id=$1 and quotation_id=$2`, tenant, resourceID).Scan(&replay.ID, &replay.OrganizationID, &replay.LeadID, &replay.CustomerSubject, &replay.VariantID, &replay.PriceBookID, &replay.Currency, &replay.TotalMinorUnits, &replay.ValidUntil, &replay.State, &replay.Version, &replay.OrderID)
		if err != nil {
			return franchisejourney.Quote{}, false, err
		}
		if err = tx.Commit(ctx); err != nil {
			return franchisejourney.Quote{}, false, err
		}
		return replay, true, nil
	}
	err = tx.QueryRow(ctx, `select b.currency,e.amount_minor_units from pricing.price_book b join pricing.price_book_entry e on e.tenant_id=b.tenant_id and e.price_book_id=b.price_book_id `+bcPriceTimeContextSQL+` where b.tenant_id=$1 and b.price_book_id=$2 and e.variant_id=$3 and `+bcPriceEligibilitySQL+``, tenant, value.PriceBookID, value.VariantID).Scan(&value.Currency, &value.TotalMinorUnits)
	if errors.Is(err, pgx.ErrNoRows) {
		return franchisejourney.Quote{}, false, franchisejourney.ErrConflict
	}
	if err != nil {
		return franchisejourney.Quote{}, false, err
	}
	result, err = tx.Exec(ctx, `insert into sales.quotation(tenant_id,quotation_id,organization_id,lead_id,customer_principal_id,variant_id,price_book_id,currency,total_minor_units,valid_until,state,version) select $1,$2,$3,$4,customer_principal_id,$5,$6,$7,$8,$9,'issued',1 from crm.lead where tenant_id=$1 and organization_id=$3 and lead_id=$4`, tenant, value.ID, value.OrganizationID, value.LeadID, value.VariantID, value.PriceBookID, value.Currency, value.TotalMinorUnits, value.ValidUntil)
	if err != nil {
		return franchisejourney.Quote{}, false, err
	}
	if result.RowsAffected() != 1 {
		return franchisejourney.Quote{}, false, franchisejourney.ErrConflict
	}
	if err = writeJourneyOutbox(ctx, tx, tenant, eventID, "quotation", value.ID, 1, "quotation.issued", map[string]any{"organization_id": value.OrganizationID, "lead_id": value.LeadID, "actor_subject": actor}); err != nil {
		return franchisejourney.Quote{}, false, err
	}
	_, err = tx.Exec(ctx, `update platform.idempotency_record set status='completed',response_code=201,response_body=jsonb_build_object('quotation_id',$3::text),resource_type='quotation',resource_id=$3,locked_until=null where tenant_id=$1 and scope='franchise-quote' and idempotency_key=$2 and status='processing'`, tenant, idempotencyKey, value.ID)
	if err != nil {
		return franchisejourney.Quote{}, false, err
	}
	if err = tx.Commit(ctx); err != nil {
		return franchisejourney.Quote{}, false, err
	}
	return value, false, nil
}

func (r *FranchiseJourney) AcceptQuote(ctx context.Context, tenant, organization, customer, quote string, version int64, evidence, orderID, lineID, quoteEventID, orderEventID string) (franchisejourney.Quote, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return franchisejourney.Quote{}, err
	}
	defer tx.Rollback(ctx)
	var value franchisejourney.Quote
	err = tx.QueryRow(ctx, `select quotation_id,organization_id,lead_id,customer_principal_id,variant_id,price_book_id,currency,total_minor_units,valid_until,state,version from sales.quotation where tenant_id=$1 and organization_id=$2 and customer_principal_id=$3 and quotation_id=$4 and version=$5 and state='issued' and valid_until>clock_timestamp() for update`, tenant, organization, customer, quote, version).Scan(&value.ID, &value.OrganizationID, &value.LeadID, &value.CustomerSubject, &value.VariantID, &value.PriceBookID, &value.Currency, &value.TotalMinorUnits, &value.ValidUntil, &value.State, &value.Version)
	if errors.Is(err, pgx.ErrNoRows) {
		return franchisejourney.Quote{}, franchisejourney.ErrConflict
	}
	if err != nil {
		return franchisejourney.Quote{}, err
	}
	// AUTHORED mapping from the already authorized one-line quote contract.
	// BC supplies document/line copying; consent, states and retention stay here.
	orderHeader, orderLines, err := bcsales.TransferQuoteToOrder(
		bcsales.Header{DocumentType: bcsales.Quote, ID: value.ID, CustomerID: value.CustomerSubject, Currency: value.Currency},
		[]bcsales.Line{{DocumentType: bcsales.Quote, DocumentID: value.ID, VariantID: value.VariantID, Quantity: 1, UnitPriceMinor: value.TotalMinorUnits}}, orderID)
	if err != nil || len(orderLines) != 1 {
		return franchisejourney.Quote{}, franchisejourney.ErrConflict
	}
	orderLine := orderLines[0]
	if _, err = tx.Exec(ctx, `insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version) values($1,$2,$3,$4,'placed',$5,$6,1)`, tenant, orderHeader.ID, organization, orderHeader.CustomerID, orderHeader.Currency, value.TotalMinorUnits); err != nil {
		return franchisejourney.Quote{}, err
	}
	if _, err = tx.Exec(ctx, `insert into sales.customer_order_line(tenant_id,order_id,line_id,variant_id,quantity,unit_price_minor_units) values($1,$2,$3,$4,$5,$6)`, tenant, orderLine.DocumentID, lineID, orderLine.VariantID, orderLine.Quantity, orderLine.UnitPriceMinor); err != nil {
		return franchisejourney.Quote{}, err
	}
	if _, err = tx.Exec(ctx, `update sales.quotation set state='accepted',order_id=$6,version=version+1,updated_at=clock_timestamp() where tenant_id=$1 and organization_id=$2 and customer_principal_id=$3 and quotation_id=$4 and version=$5`, tenant, organization, customer, quote, version, orderID); err != nil {
		return franchisejourney.Quote{}, err
	}
	if _, err = tx.Exec(ctx, `insert into sales.quotation_acceptance(tenant_id,quotation_id,order_id,customer_principal_id,evidence_sha256_hex) values($1,$2,$3,$4,$5)`, tenant, quote, orderID, customer, evidence); err != nil {
		return franchisejourney.Quote{}, err
	}
	value.State = "accepted"
	value.Version++
	value.OrderID = orderID
	if err = writeJourneyOutbox(ctx, tx, tenant, quoteEventID, "quotation", quote, value.Version, "quotation.accepted", map[string]any{"organization_id": organization, "customer_subject": customer, "order_id": orderID}); err != nil {
		return franchisejourney.Quote{}, err
	}
	if err = writeJourneyOutbox(ctx, tx, tenant, orderEventID, "customer-order", orderID, 1, "customer-order.placed", map[string]any{"organization_id": organization, "customer_subject": customer, "quotation_id": quote, "currency": value.Currency, "total_minor_units": value.TotalMinorUnits}); err != nil {
		return franchisejourney.Quote{}, err
	}
	return value, tx.Commit(ctx)
}

func (r *FranchiseJourney) PublishDeliveryChecklist(ctx context.Context, tenant, subject string, value franchisejourney.DeliveryChecklist, eventID string) (franchisejourney.DeliveryChecklist, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return franchisejourney.DeliveryChecklist{}, err
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, `insert into sales.delivery_checklist_template(tenant_id,organization_id,checklist_id,checklist_version,title,state,created_by_subject) values($1,$2,$3,$4,$5,'draft',$6)`, tenant, value.OrganizationID, value.ID, value.Version, value.Title, subject)
	if postgresConflict(err) {
		return franchisejourney.DeliveryChecklist{}, franchisejourney.ErrConflict
	}
	if err != nil {
		return franchisejourney.DeliveryChecklist{}, err
	}
	for _, item := range value.Items {
		_, err = tx.Exec(ctx, `insert into sales.delivery_checklist_item(tenant_id,organization_id,checklist_id,checklist_version,item_id,ordinal,prompt,response_type,required) values($1,$2,$3,$4,$5,$6,$7,$8,$9)`, tenant, value.OrganizationID, value.ID, value.Version, item.ID, item.Ordinal, item.Prompt, item.ResponseType, item.Required)
		if postgresConflict(err) {
			return franchisejourney.DeliveryChecklist{}, franchisejourney.ErrConflict
		}
		if err != nil {
			return franchisejourney.DeliveryChecklist{}, err
		}
	}
	if err = tx.QueryRow(ctx, `update sales.delivery_checklist_template set state='published',published_at=clock_timestamp() where tenant_id=$1 and organization_id=$2 and checklist_id=$3 and checklist_version=$4 and state='draft' returning state`, tenant, value.OrganizationID, value.ID, value.Version).Scan(&value.State); err != nil {
		return franchisejourney.DeliveryChecklist{}, err
	}
	if err = writeJourneyOutbox(ctx, tx, tenant, eventID, "delivery-checklist", value.ID, value.Version, "delivery-checklist.published", map[string]any{"organization_id": value.OrganizationID, "checklist_version": value.Version, "item_count": len(value.Items), "actor_subject": subject}); err != nil {
		return franchisejourney.DeliveryChecklist{}, err
	}
	return value, tx.Commit(ctx)
}

func (r *FranchiseJourney) CompleteDeliveryChecklist(ctx context.Context, tenant, organization, subject, handover string, version int64, checklistID string, checklistVersion int64, responses []franchisejourney.ChecklistResponse, eventID string) (franchisejourney.Handover, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return franchisejourney.Handover{}, err
	}
	defer tx.Rollback(ctx)
	if err = lockDeliveryScope(ctx, tx, tenant, organization, handover); err != nil {
		return franchisejourney.Handover{}, err
	}
	for _, response := range responses {
		var evidence any
		if response.EvidenceSHA256 != "" {
			evidence = response.EvidenceSHA256
		}
		_, err = tx.Exec(ctx, `insert into sales.delivery_checklist_response(tenant_id,organization_id,handover_id,checklist_id,checklist_version,item_id,response_text,evidence_sha256_hex,answered_by_subject) values($1,$2,$3,$4,$5,$6,$7,$8,$9)`, tenant, organization, handover, checklistID, checklistVersion, response.ItemID, response.ResponseText, evidence, subject)
		if postgresConflict(err) {
			return franchisejourney.Handover{}, franchisejourney.ErrConflict
		}
		if err != nil {
			return franchisejourney.Handover{}, err
		}
	}
	var value franchisejourney.Handover
	value.ChecklistItems = []franchisejourney.ChecklistItem{}
	err = tx.QueryRow(ctx, `update sales.delivery_handover h set state='presented',checklist_id=$5,checklist_version=$6,checklist_completed_at=clock_timestamp(),checklist_completed_by_subject=$7,updated_at=clock_timestamp(),version=h.version+1 from sales.delivery_checklist_template t where h.tenant_id=$1 and h.organization_id=$2 and h.handover_id=$3 and h.version=$4 and h.state='prepared' and t.tenant_id=h.tenant_id and t.organization_id=h.organization_id and t.checklist_id=$5 and t.checklist_version=$6 and t.state='published' returning h.handover_id,h.organization_id,h.order_id,h.customer_principal_id,h.stock_unit_id,h.state,h.version,h.checklist_id,h.checklist_version,t.title,h.checklist_completed_at`, tenant, organization, handover, version, checklistID, checklistVersion, subject).Scan(&value.ID, &value.OrganizationID, &value.OrderID, &value.CustomerSubject, &value.StockUnitID, &value.State, &value.Version, &value.ChecklistID, &value.ChecklistVersion, &value.ChecklistTitle, &value.ChecklistCompletedAt)
	if errors.Is(err, pgx.ErrNoRows) || postgresConflict(err) {
		return franchisejourney.Handover{}, franchisejourney.ErrConflict
	}
	if err != nil {
		return franchisejourney.Handover{}, err
	}
	if err = writeJourneyOutbox(ctx, tx, tenant, eventID, "delivery-handover", handover, value.Version, "delivery-handover.checklist-completed", map[string]any{"organization_id": organization, "checklist_id": checklistID, "checklist_version": checklistVersion, "response_count": len(responses), "actor_subject": subject}); err != nil {
		return franchisejourney.Handover{}, err
	}
	return value, tx.Commit(ctx)
}

func (r *FranchiseJourney) RejectHandover(ctx context.Context, tenant, organization, customer, handover string, version int64, reasonCode, details, evidence, exceptionID, eventID string) (franchisejourney.DeliveryException, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return franchisejourney.DeliveryException{}, err
	}
	defer tx.Rollback(ctx)
	if err = lockDeliveryScope(ctx, tx, tenant, organization, handover); err != nil {
		return franchisejourney.DeliveryException{}, err
	}
	var nextVersion int64
	err = tx.QueryRow(ctx, `update sales.delivery_handover set state='rejected',updated_at=clock_timestamp(),version=version+1 where tenant_id=$1 and organization_id=$2 and customer_principal_id=$3 and handover_id=$4 and version=$5 and state='presented' and checklist_completed_at is not null returning version`, tenant, organization, customer, handover, version).Scan(&nextVersion)
	if errors.Is(err, pgx.ErrNoRows) {
		return franchisejourney.DeliveryException{}, franchisejourney.ErrConflict
	}
	if err != nil {
		return franchisejourney.DeliveryException{}, err
	}
	value := franchisejourney.DeliveryException{ID: exceptionID, OrganizationID: organization, HandoverID: handover, CustomerSubject: customer, ReasonCode: reasonCode, Details: details, State: "open", Version: 1}
	err = tx.QueryRow(ctx, `insert into sales.delivery_exception(tenant_id,exception_id,organization_id,handover_id,customer_principal_id,reason_code,details,rejection_evidence_sha256_hex,state,version) values($1,$2,$3,$4,$5,$6,$7,$8,'open',1) returning created_at`, tenant, exceptionID, organization, handover, customer, reasonCode, details, evidence).Scan(&value.CreatedAt)
	if postgresConflict(err) {
		return franchisejourney.DeliveryException{}, franchisejourney.ErrConflict
	}
	if err != nil {
		return franchisejourney.DeliveryException{}, err
	}
	if err = writeJourneyOutbox(ctx, tx, tenant, eventID, "delivery-handover", handover, nextVersion, "delivery-handover.rejected", map[string]any{"organization_id": organization, "customer_subject": customer, "exception_id": exceptionID, "reason_code": reasonCode, "rejection_evidence_sha256": evidence}); err != nil {
		return franchisejourney.DeliveryException{}, err
	}
	return value, tx.Commit(ctx)
}

func (r *FranchiseJourney) DeliveryExceptions(ctx context.Context, tenant, organization string, limit int) ([]franchisejourney.DeliveryException, error) {
	rows, err := r.pool.Query(ctx, `select e.exception_id,e.organization_id,e.handover_id,e.customer_principal_id,e.reason_code,e.details,e.state,e.version,e.created_at,e.resolved_at,coalesce(x.action,''),coalesce(x.successor_handover_id,''),coalesce(x.return_authorization_id,''),exists(select 1 from sales.delivery_handover h join sales.customer_order o on o.tenant_id=h.tenant_id and o.order_id=h.order_id and o.organization_id=h.organization_id and o.customer_principal_id=h.customer_principal_id join inventory.stock_unit su on su.tenant_id=h.tenant_id and su.stock_unit_id=h.stock_unit_id and su.organization_id=h.organization_id where h.tenant_id=e.tenant_id and h.handover_id=e.handover_id and h.organization_id=e.organization_id and h.customer_principal_id=e.customer_principal_id) from sales.delivery_exception e left join sales.delivery_exception_resolution x on x.tenant_id=e.tenant_id and x.exception_id=e.exception_id where e.tenant_id=$1 and e.organization_id=$2 order by e.created_at desc,e.exception_id limit $3`, tenant, organization, limit)
	if err != nil {
		return nil, err
	}
	return scanDeliveryExceptions(rows)
}

func scanDeliveryExceptions(rows pgx.Rows) ([]franchisejourney.DeliveryException, error) {
	defer rows.Close()
	items := []franchisejourney.DeliveryException{}
	for rows.Next() {
		var item franchisejourney.DeliveryException
		var coherent bool
		if err := rows.Scan(&item.ID, &item.OrganizationID, &item.HandoverID, &item.CustomerSubject, &item.ReasonCode, &item.Details, &item.State, &item.Version, &item.CreatedAt, &item.ResolvedAt, &item.ResolutionAction, &item.SuccessorHandoverID, &item.ReturnAuthorizationID, &coherent); err != nil {
			return nil, err
		}
		if !coherent {
			return nil, franchisejourney.ErrConflict
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *FranchiseJourney) ResolveDeliveryException(ctx context.Context, tenant, organization, subject, exceptionID string, version int64, action, notes, successorID, authorizationID, resolutionID, eventID string) (franchisejourney.DeliveryResolution, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return franchisejourney.DeliveryResolution{}, err
	}
	defer tx.Rollback(ctx)
	var handover franchisejourney.Handover
	var customer, orderID, stockUnitID string
	if err = tx.QueryRow(ctx, `select h.handover_id,h.customer_principal_id,h.order_id,h.stock_unit_id from sales.delivery_exception e join sales.delivery_handover h on h.tenant_id=e.tenant_id and h.handover_id=e.handover_id where e.tenant_id=$1 and e.organization_id=$2 and e.exception_id=$3 and e.state='open' and e.version=$4 and h.state='rejected' for update of e,h`, tenant, organization, exceptionID, version).Scan(&handover.ID, &customer, &orderID, &stockUnitID); errors.Is(err, pgx.ErrNoRows) {
		return franchisejourney.DeliveryResolution{}, franchisejourney.ErrConflict
	}
	if err != nil {
		return franchisejourney.DeliveryResolution{}, err
	}
	result := franchisejourney.DeliveryResolution{}
	if err = lockDeliveryScope(ctx, tx, tenant, organization, handover.ID); err != nil {
		return franchisejourney.DeliveryResolution{}, err
	}
	originalHandoverID := handover.ID
	var successor any
	var authorization any
	if action == "correct-and-represent" {
		successor = successorID
		handover = franchisejourney.Handover{ID: successorID, OrganizationID: organization, OrderID: orderID, CustomerSubject: customer, StockUnitID: stockUnitID, State: "prepared", Version: 1, ChecklistItems: []franchisejourney.ChecklistItem{}, SupersedesHandoverID: originalHandoverID}
		_, err = tx.Exec(ctx, `insert into sales.delivery_handover(tenant_id,handover_id,organization_id,order_id,customer_principal_id,stock_unit_id,state,version,supersedes_handover_id) values($1,$2,$3,$4,$5,$6,'prepared',1,$7)`, tenant, successorID, organization, orderID, customer, stockUnitID, originalHandoverID)
		result.SuccessorHandover = &handover
	} else {
		authorization = authorizationID
		_, err = tx.Exec(ctx, `insert into sales.return_authorization(tenant_id,authorization_id,organization_id,exception_id,handover_id,order_id,stock_unit_id,customer_principal_id,disposition,exact_cost_source_order_id,state,authorized_by_subject) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$6,'authorized',$10)`, tenant, authorizationID, organization, exceptionID, originalHandoverID, orderID, stockUnitID, customer, action, subject)
		result.ReturnAuthorizationID = authorizationID
		result.Disposition = action
	}
	if postgresConflict(err) {
		return franchisejourney.DeliveryResolution{}, franchisejourney.ErrConflict
	}
	if err != nil {
		return franchisejourney.DeliveryResolution{}, err
	}
	_, err = tx.Exec(ctx, `insert into sales.delivery_exception_resolution(tenant_id,resolution_id,exception_id,action,notes,resolved_by_subject,successor_handover_id,return_authorization_id) values($1,$2,$3,$4,$5,$6,$7,$8)`, tenant, resolutionID, exceptionID, action, notes, subject, successor, authorization)
	if postgresConflict(err) {
		return franchisejourney.DeliveryResolution{}, franchisejourney.ErrConflict
	}
	if err != nil {
		return franchisejourney.DeliveryResolution{}, err
	}
	resolved := franchisejourney.DeliveryException{ID: exceptionID, OrganizationID: organization, HandoverID: originalHandoverID, CustomerSubject: customer, State: "resolved", Version: version + 1, ResolutionAction: action}
	if result.SuccessorHandover != nil {
		resolved.SuccessorHandoverID = successorID
	} else {
		resolved.ReturnAuthorizationID = authorizationID
	}
	err = tx.QueryRow(ctx, `update sales.delivery_exception set state='resolved',resolved_at=clock_timestamp(),resolved_by_subject=$5,version=version+1 where tenant_id=$1 and organization_id=$2 and exception_id=$3 and version=$4 and state='open' returning reason_code,details,created_at,resolved_at`, tenant, organization, exceptionID, version, subject).Scan(&resolved.ReasonCode, &resolved.Details, &resolved.CreatedAt, &resolved.ResolvedAt)
	if errors.Is(err, pgx.ErrNoRows) || postgresConflict(err) {
		return franchisejourney.DeliveryResolution{}, franchisejourney.ErrConflict
	}
	if err != nil {
		return franchisejourney.DeliveryResolution{}, err
	}
	result.Exception = resolved
	if err = writeJourneyOutbox(ctx, tx, tenant, eventID, "delivery-exception", exceptionID, resolved.Version, "delivery-exception.resolved", map[string]any{"organization_id": organization, "action": action, "successor_handover_id": resolved.SuccessorHandoverID, "return_authorization_id": resolved.ReturnAuthorizationID, "actor_subject": subject}); err != nil {
		return franchisejourney.DeliveryResolution{}, err
	}
	return result, tx.Commit(ctx)
}

func (r *FranchiseJourney) ReturnCases(ctx context.Context, tenant, organization string, limit int) ([]franchisejourney.ReturnCase, error) {
	if limit < 1 || limit > 100 {
		return nil, franchisejourney.ErrInvalid
	}
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead, AccessMode: pgx.ReadOnly})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	rows, err := tx.Query(ctx, `select a.authorization_id,a.organization_id,a.order_id,a.stock_unit_id,a.customer_principal_id,a.disposition,a.authorized_at,r.receipt_id,r.received_serial_number,r.condition_code,r.notes,r.evidence_sha256_hex,r.received_by_subject,r.received_at,d.disposition_id,d.inventory_action,d.customer_remedy,d.notes,d.decided_by_subject,d.decided_at, `+returnScopePredicate+` and (r.receipt_id is null or (r.organization_id,r.order_id,r.stock_unit_id,r.customer_principal_id)=(a.organization_id,a.order_id,a.stock_unit_id,a.customer_principal_id)) from sales.return_authorization a left join sales.return_receipt r on r.tenant_id=a.tenant_id and r.authorization_id=a.authorization_id left join sales.return_disposition d on d.tenant_id=r.tenant_id and d.receipt_id=r.receipt_id where a.tenant_id=$1 and a.organization_id=$2 order by a.authorized_at desc,a.authorization_id limit $3`, tenant, organization, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []franchisejourney.ReturnCase{}
	dispositions := map[string]*franchisejourney.ReturnDisposition{}
	for rows.Next() {
		var item franchisejourney.ReturnCase
		var receiptID, serial, condition, receiptNotes, evidence, receiver *string
		var receivedAt *time.Time
		var dispositionID, inventoryAction, remedy, dispositionNotes, decider *string
		var decidedAt *time.Time
		var scopeValid bool
		if err = rows.Scan(&item.AuthorizationID, &item.OrganizationID, &item.OrderID, &item.StockUnitID, &item.CustomerSubject, &item.AuthorizedAction, &item.AuthorizedAt, &receiptID, &serial, &condition, &receiptNotes, &evidence, &receiver, &receivedAt, &dispositionID, &inventoryAction, &remedy, &dispositionNotes, &decider, &decidedAt, &scopeValid); err != nil {
			return nil, err
		}
		if !scopeValid {
			return nil, franchisejourney.ErrConflict
		}
		if receiptID != nil {
			item.Receipt = &franchisejourney.ReturnReceipt{ID: *receiptID, AuthorizationID: item.AuthorizationID, OrganizationID: item.OrganizationID, OrderID: item.OrderID, StockUnitID: item.StockUnitID, CustomerSubject: item.CustomerSubject, ReceivedSerialNumber: *serial, ConditionCode: *condition, Notes: *receiptNotes, EvidenceSHA256: *evidence, ReceivedBySubject: *receiver, ReceivedAt: *receivedAt}
		}
		if dispositionID != nil {
			item.Disposition = &franchisejourney.ReturnDisposition{ID: *dispositionID, ReceiptID: *receiptID, InventoryAction: *inventoryAction, CustomerRemedy: *remedy, Notes: *dispositionNotes, DecidedBySubject: *decider, DecidedAt: *decidedAt, Effects: []franchisejourney.ReturnEffectRequest{}}
			dispositions[*dispositionID] = item.Disposition
		}
		items = append(items, item)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	rows.Close()
	if len(dispositions) == 0 {
		return items, tx.Commit(ctx)
	}
	dispositionIDs := make([]string, 0, len(dispositions))
	for id := range dispositions {
		dispositionIDs = append(dispositionIDs, id)
	}
	effectRows, err := tx.Query(ctx, `select e.disposition_id,e.request_id,e.effect_kind,e.owner_context,e.state,e.idempotency_key,e.requested_at from sales.return_effect_request e join sales.return_disposition d on d.tenant_id=e.tenant_id and d.disposition_id=e.disposition_id join sales.return_receipt r on r.tenant_id=d.tenant_id and r.receipt_id=d.receipt_id where e.tenant_id=$1 and r.organization_id=$2 and e.disposition_id=any($3::text[]) order by e.requested_at,e.request_id`, tenant, organization, dispositionIDs)
	if err != nil {
		return nil, err
	}
	defer effectRows.Close()
	for effectRows.Next() {
		var dispositionID string
		var effect franchisejourney.ReturnEffectRequest
		if err = effectRows.Scan(&dispositionID, &effect.ID, &effect.EffectKind, &effect.OwnerContext, &effect.State, &effect.IdempotencyKey, &effect.RequestedAt); err != nil {
			return nil, err
		}
		if disposition := dispositions[dispositionID]; disposition != nil {
			disposition.Effects = append(disposition.Effects, effect)
		}
	}
	effectRows.Close()
	if err = effectRows.Err(); err != nil {
		return nil, err
	}
	return items, tx.Commit(ctx)
}

func (r *FranchiseJourney) ReceiveReturn(ctx context.Context, tenant, organization, subject, authorizationID, serialNumber, conditionCode, notes, evidence, receiptID, eventID string) (franchisejourney.ReturnReceipt, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return franchisejourney.ReturnReceipt{}, err
	}
	defer tx.Rollback(ctx)
	if err = lockReturnAuthorizationScope(ctx, tx, tenant, organization, authorizationID); err != nil {
		return franchisejourney.ReturnReceipt{}, err
	}
	value := franchisejourney.ReturnReceipt{ID: receiptID, AuthorizationID: authorizationID, OrganizationID: organization, ReceivedSerialNumber: serialNumber, ConditionCode: conditionCode, Notes: notes, EvidenceSHA256: evidence, ReceivedBySubject: subject}
	err = tx.QueryRow(ctx, `insert into sales.return_receipt(tenant_id,receipt_id,authorization_id,organization_id,order_id,stock_unit_id,customer_principal_id,received_serial_number,condition_code,notes,evidence_sha256_hex,received_by_subject) select a.tenant_id,$4,a.authorization_id,a.organization_id,a.order_id,a.stock_unit_id,a.customer_principal_id,$5,$6,$7,$8,$9 from sales.return_authorization a where a.tenant_id=$1 and a.organization_id=$2 and a.authorization_id=$3 and a.state='authorized' returning order_id,stock_unit_id,customer_principal_id,received_at`, tenant, organization, authorizationID, receiptID, serialNumber, conditionCode, notes, evidence, subject).Scan(&value.OrderID, &value.StockUnitID, &value.CustomerSubject, &value.ReceivedAt)
	if errors.Is(err, pgx.ErrNoRows) || postgresConflict(err) {
		return franchisejourney.ReturnReceipt{}, franchisejourney.ErrConflict
	}
	if err != nil {
		return franchisejourney.ReturnReceipt{}, err
	}
	if err = writeJourneyOutbox(ctx, tx, tenant, eventID, "return-authorization", authorizationID, 1, "return.received", map[string]any{"organization_id": organization, "receipt_id": receiptID, "order_id": value.OrderID, "stock_unit_id": value.StockUnitID, "condition_code": conditionCode, "evidence_sha256": evidence, "actor_subject": subject}); err != nil {
		return franchisejourney.ReturnReceipt{}, err
	}
	return value, tx.Commit(ctx)
}

func (r *FranchiseJourney) DecideReturn(ctx context.Context, tenant, organization, subject, receiptID, inventoryAction, notes, dispositionID, inventoryRequestID, remedyRequestID, accountingRequestID, fiscalRequestID, eventID string) (franchisejourney.ReturnDisposition, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return franchisejourney.ReturnDisposition{}, err
	}
	defer tx.Rollback(ctx)
	var authorizationID string
	err = tx.QueryRow(ctx, `select authorization_id from sales.return_receipt where tenant_id=$1 and organization_id=$2 and receipt_id=$3`, tenant, organization, receiptID).Scan(&authorizationID)
	if errors.Is(err, pgx.ErrNoRows) {
		return franchisejourney.ReturnDisposition{}, franchisejourney.ErrConflict
	}
	if err != nil {
		return franchisejourney.ReturnDisposition{}, err
	}
	if err = lockReturnAuthorizationScope(ctx, tx, tenant, organization, authorizationID); err != nil {
		return franchisejourney.ReturnDisposition{}, err
	}

	var remedy string
	err = tx.QueryRow(ctx, `select case a.disposition when 'return' then 'refund' else 'exchange' end from sales.return_receipt r join sales.return_authorization a on a.tenant_id=r.tenant_id and a.authorization_id=r.authorization_id where r.tenant_id=$1 and r.organization_id=$2 and r.receipt_id=$3 and (r.organization_id,r.order_id,r.stock_unit_id,r.customer_principal_id)=(a.organization_id,a.order_id,a.stock_unit_id,a.customer_principal_id) for update of r`, tenant, organization, receiptID).Scan(&remedy)
	if errors.Is(err, pgx.ErrNoRows) {
		return franchisejourney.ReturnDisposition{}, franchisejourney.ErrConflict
	}
	if err != nil {
		return franchisejourney.ReturnDisposition{}, err
	}
	value := franchisejourney.ReturnDisposition{ID: dispositionID, ReceiptID: receiptID, InventoryAction: inventoryAction, CustomerRemedy: remedy, Notes: notes, DecidedBySubject: subject, Effects: []franchisejourney.ReturnEffectRequest{}}
	err = tx.QueryRow(ctx, `insert into sales.return_disposition(tenant_id,disposition_id,receipt_id,inventory_action,customer_remedy,notes,decided_by_subject) values($1,$2,$3,$4,$5,$6,$7) returning decided_at`, tenant, dispositionID, receiptID, inventoryAction, remedy, notes, subject).Scan(&value.DecidedAt)
	if postgresConflict(err) {
		return franchisejourney.ReturnDisposition{}, franchisejourney.ErrConflict
	}
	if err != nil {
		return franchisejourney.ReturnDisposition{}, err
	}
	effects := []struct{ id, kind, owner string }{{inventoryRequestID, "inventory", "inventory"}, {remedyRequestID, remedy, map[string]string{"refund": "payment", "exchange": "fulfillment"}[remedy]}, {accountingRequestID, "accounting", "accounting"}}
	if remedy == "refund" {
		effects = append(effects, struct{ id, kind, owner string }{fiscalRequestID, "fiscal", "fiscal"})
	}
	for _, requested := range effects {
		var effect franchisejourney.ReturnEffectRequest
		err = tx.QueryRow(ctx, `insert into sales.return_effect_request(tenant_id,request_id,disposition_id,effect_kind,owner_context,state,idempotency_key) values($1,$2,$3,$4,$5,'requested',$2) returning request_id,effect_kind,owner_context,state,idempotency_key,requested_at`, tenant, requested.id, dispositionID, requested.kind, requested.owner).Scan(&effect.ID, &effect.EffectKind, &effect.OwnerContext, &effect.State, &effect.IdempotencyKey, &effect.RequestedAt)
		if postgresConflict(err) {
			return franchisejourney.ReturnDisposition{}, franchisejourney.ErrConflict
		}
		if err != nil {
			return franchisejourney.ReturnDisposition{}, err
		}
		value.Effects = append(value.Effects, effect)
	}
	if err = writeJourneyOutbox(ctx, tx, tenant, eventID, "return-disposition", dispositionID, 1, "return.disposition-requested", map[string]any{"organization_id": organization, "receipt_id": receiptID, "inventory_action": inventoryAction, "customer_remedy": remedy, "effect_count": len(value.Effects), "actor_subject": subject}); err != nil {
		return franchisejourney.ReturnDisposition{}, err
	}
	return value, tx.Commit(ctx)
}

func (r *FranchiseJourney) CustomerJourney(ctx context.Context, tenant, organization, customer string) (franchisejourney.CustomerJourney, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead, AccessMode: pgx.ReadOnly})
	if err != nil {
		return franchisejourney.CustomerJourney{}, err
	}
	defer tx.Rollback(ctx)
	result, err := readCustomerJourney(ctx, tx, tenant, organization, customer)
	if err != nil {
		return franchisejourney.CustomerJourney{}, err
	}
	if err = tx.Commit(ctx); err != nil {
		return franchisejourney.CustomerJourney{}, err
	}
	return result, nil
}

func readCustomerJourney(ctx context.Context, tx pgx.Tx, tenant, organization, customer string) (franchisejourney.CustomerJourney, error) {
	result := franchisejourney.CustomerJourney{Appointments: []franchisejourney.Appointment{}, Quotes: []franchisejourney.Quote{}, Handovers: []franchisejourney.Handover{}, Exceptions: []franchisejourney.DeliveryException{}}
	rows, err := tx.Query(ctx, `select appointment_id,organization_id,lead_id,coalesce(model_id,''),appointment_kind,starts_at,state,version,coalesce(slot_id,''),coalesce(ends_at,starts_at) from crm.appointment where tenant_id=$1 and organization_id=$2 and customer_principal_id=$3 order by starts_at desc limit 100`, tenant, organization, customer)
	if err != nil {
		return result, err
	}
	for rows.Next() {
		var v franchisejourney.Appointment
		if err = rows.Scan(&v.ID, &v.OrganizationID, &v.LeadID, &v.ModelID, &v.Kind, &v.StartsAt, &v.State, &v.Version, &v.SlotID, &v.EndsAt); err != nil {
			rows.Close()
			return result, err
		}
		result.Appointments = append(result.Appointments, v)
	}
	rows.Close()
	if err = rows.Err(); err != nil {
		return result, err
	}
	rows, err = tx.Query(ctx, `select quotation_id,organization_id,lead_id,coalesce(customer_principal_id,''),variant_id,price_book_id,currency,total_minor_units,valid_until,state,version,coalesce(order_id,'') from sales.quotation where tenant_id=$1 and organization_id=$2 and customer_principal_id=$3 order by created_at desc limit 100`, tenant, organization, customer)
	if err != nil {
		return result, err
	}
	for rows.Next() {
		var v franchisejourney.Quote
		if err = rows.Scan(&v.ID, &v.OrganizationID, &v.LeadID, &v.CustomerSubject, &v.VariantID, &v.PriceBookID, &v.Currency, &v.TotalMinorUnits, &v.ValidUntil, &v.State, &v.Version, &v.OrderID); err != nil {
			rows.Close()
			return result, err
		}
		result.Quotes = append(result.Quotes, v)
	}
	rows.Close()
	if err = rows.Err(); err != nil {
		return result, err
	}
	rows, err = tx.Query(ctx, `select h.handover_id,h.organization_id,h.order_id,h.customer_principal_id,h.stock_unit_id,h.state,h.version,h.customer_accepted_at,coalesce(h.acceptance_evidence_sha256_hex,''),coalesce(h.checklist_id,''),coalesce(h.checklist_version,0),coalesce(t.title,''),h.checklist_completed_at,coalesce(h.supersedes_handover_id,''),exists(select 1 from sales.customer_order o join inventory.stock_unit su on su.tenant_id=o.tenant_id and su.stock_unit_id=h.stock_unit_id and su.organization_id=h.organization_id where o.tenant_id=h.tenant_id and o.order_id=h.order_id and o.organization_id=h.organization_id and o.customer_principal_id=h.customer_principal_id) from sales.delivery_handover h left join sales.delivery_checklist_template t on t.tenant_id=h.tenant_id and t.organization_id=h.organization_id and t.checklist_id=h.checklist_id and t.checklist_version=h.checklist_version where h.tenant_id=$1 and h.organization_id=$2 and h.customer_principal_id=$3 order by h.created_at desc limit 100`, tenant, organization, customer)
	if err != nil {
		return result, err
	}
	for rows.Next() {
		var v franchisejourney.Handover
		var coherent bool
		v.ChecklistItems = []franchisejourney.ChecklistItem{}
		if err = rows.Scan(&v.ID, &v.OrganizationID, &v.OrderID, &v.CustomerSubject, &v.StockUnitID, &v.State, &v.Version, &v.CustomerAcceptedAt, &v.AcceptanceEvidence, &v.ChecklistID, &v.ChecklistVersion, &v.ChecklistTitle, &v.ChecklistCompletedAt, &v.SupersedesHandoverID, &coherent); err != nil {
			rows.Close()
			return result, err
		}
		if !coherent {
			rows.Close()
			return result, franchisejourney.ErrConflict
		}
		result.Handovers = append(result.Handovers, v)
	}
	rows.Close()
	if err = rows.Err(); err != nil {
		return result, err
	}
	indexes := map[string]int{}
	handoverIDs := []string{}
	for index := range result.Handovers {
		indexes[result.Handovers[index].ID] = index
		handoverIDs = append(handoverIDs, result.Handovers[index].ID)
	}
	rows, err = tx.Query(ctx, `select h.handover_id,i.item_id,i.ordinal,i.prompt,i.response_type,i.required from sales.delivery_handover h join sales.delivery_checklist_item i on i.tenant_id=h.tenant_id and i.organization_id=h.organization_id and i.checklist_id=h.checklist_id and i.checklist_version=h.checklist_version where h.tenant_id=$1 and h.organization_id=$2 and h.customer_principal_id=$3 and h.handover_id=any($4::text[]) order by h.created_at desc,h.handover_id,i.ordinal`, tenant, organization, customer, handoverIDs)
	if err != nil {
		return result, err
	}
	for rows.Next() {
		var handoverID string
		var item franchisejourney.ChecklistItem
		if err = rows.Scan(&handoverID, &item.ID, &item.Ordinal, &item.Prompt, &item.ResponseType, &item.Required); err != nil {
			rows.Close()
			return result, err
		}
		if index, ok := indexes[handoverID]; ok {
			result.Handovers[index].ChecklistItems = append(result.Handovers[index].ChecklistItems, item)
		}
	}
	rows.Close()
	if err = rows.Err(); err != nil {
		return result, err
	}
	rows, err = tx.Query(ctx, `select e.exception_id,e.organization_id,e.handover_id,e.customer_principal_id,e.reason_code,e.details,e.state,e.version,e.created_at,e.resolved_at,coalesce(x.action,''),coalesce(x.successor_handover_id,''),coalesce(x.return_authorization_id,''),exists(select 1 from sales.delivery_handover h join sales.customer_order o on o.tenant_id=h.tenant_id and o.order_id=h.order_id and o.organization_id=h.organization_id and o.customer_principal_id=h.customer_principal_id join inventory.stock_unit su on su.tenant_id=h.tenant_id and su.stock_unit_id=h.stock_unit_id and su.organization_id=h.organization_id where h.tenant_id=e.tenant_id and h.handover_id=e.handover_id and h.organization_id=e.organization_id and h.customer_principal_id=e.customer_principal_id) from sales.delivery_exception e left join sales.delivery_exception_resolution x on x.tenant_id=e.tenant_id and x.exception_id=e.exception_id where e.tenant_id=$1 and e.organization_id=$2 and e.customer_principal_id=$3 order by e.created_at desc,e.exception_id limit 100`, tenant, organization, customer)
	if err != nil {
		return result, err
	}
	result.Exceptions, err = scanDeliveryExceptions(rows)
	return result, err
}

func postgresConflict(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && (pgErr.Code == "23503" || pgErr.Code == "23505" || pgErr.Code == "23514" || pgErr.Code == "23P01")
}

func (r *FranchiseJourney) AcceptHandover(ctx context.Context, tenant, organization, customer, handover string, version int64, serialNumber, checklistID string, checklistVersion int64, evidence, eventID string) (franchisejourney.Handover, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return franchisejourney.Handover{}, err
	}
	defer tx.Rollback(ctx)
	if err = lockDeliveryScope(ctx, tx, tenant, organization, handover); err != nil {
		return franchisejourney.Handover{}, err
	}
	var value franchisejourney.Handover
	value.ChecklistItems = []franchisejourney.ChecklistItem{}
	err = tx.QueryRow(ctx, `update sales.delivery_handover h set state='accepted',acceptance_evidence_sha256_hex=$6,customer_accepted_at=clock_timestamp(),updated_at=clock_timestamp(),version=h.version+1 from inventory.stock_unit i,sales.delivery_checklist_template t where h.tenant_id=$1 and h.organization_id=$2 and h.customer_principal_id=$3 and h.handover_id=$4 and h.version=$5 and h.state='presented' and h.checklist_completed_at is not null and h.checklist_id=$8 and h.checklist_version=$9 and i.tenant_id=h.tenant_id and i.stock_unit_id=h.stock_unit_id and i.serial_number=$7 and t.tenant_id=h.tenant_id and t.organization_id=h.organization_id and t.checklist_id=h.checklist_id and t.checklist_version=h.checklist_version returning h.handover_id,h.organization_id,h.order_id,h.customer_principal_id,h.stock_unit_id,h.state,h.version,h.customer_accepted_at,h.acceptance_evidence_sha256_hex,h.checklist_id,h.checklist_version,t.title,h.checklist_completed_at`, tenant, organization, customer, handover, version, evidence, serialNumber, checklistID, checklistVersion).Scan(&value.ID, &value.OrganizationID, &value.OrderID, &value.CustomerSubject, &value.StockUnitID, &value.State, &value.Version, &value.CustomerAcceptedAt, &value.AcceptanceEvidence, &value.ChecklistID, &value.ChecklistVersion, &value.ChecklistTitle, &value.ChecklistCompletedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return franchisejourney.Handover{}, franchisejourney.ErrConflict
	}
	if err != nil {
		return franchisejourney.Handover{}, err
	}
	if err = writeJourneyOutbox(ctx, tx, tenant, eventID, "delivery-handover", handover, value.Version, "delivery-handover.accepted", map[string]any{"organization_id": organization, "customer_subject": customer, "acceptance_evidence_sha256": evidence, "checklist_id": checklistID, "checklist_version": checklistVersion}); err != nil {
		return franchisejourney.Handover{}, err
	}
	return value, tx.Commit(ctx)
}

// Tenant-scoped foreign keys alone do not authorize a relationship. Keep the
// linked rows stable until the command and its audit commit together. This is
// an integrity check, NOT proof of payment, financial release or delivery.
func lockDeliveryScope(ctx context.Context, tx pgx.Tx, tenant, organization, handover string) error {
	var id string
	err := tx.QueryRow(ctx, `select h.handover_id
		from sales.delivery_handover h
		join sales.customer_order o on o.tenant_id=h.tenant_id and o.order_id=h.order_id
			and o.organization_id=h.organization_id and o.customer_principal_id=h.customer_principal_id
		join inventory.stock_unit i on i.tenant_id=h.tenant_id and i.stock_unit_id=h.stock_unit_id
			and i.organization_id=h.organization_id
		where h.tenant_id=$1 and h.organization_id=$2 and h.handover_id=$3
		for update of h for share of o,i`, tenant, organization, handover).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return franchisejourney.ErrConflict
	}
	return err
}

func writeJourneyOutbox(ctx context.Context, tx pgx.Tx, tenant, eventID, aggregateType, aggregateID string, version int64, eventType string, payload any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload) values($1,$2,$3,$4,$5,$6,1,clock_timestamp(),$7)`, tenant, eventID, aggregateType, aggregateID, version, eventType, body)
	if err != nil {
		return fmt.Errorf("write journey outbox: %w", err)
	}
	return nil
}

func (r *FranchiseJourney) QuoteResult(ctx context.Context, tenant, organization, lead, key string) (franchisejourney.Quote, error) {
	var v franchisejourney.Quote
	err := r.pool.QueryRow(ctx, `select q.quotation_id,q.organization_id,q.lead_id,q.variant_id,q.price_book_id,q.currency,q.total_minor_units,q.valid_until,q.state,q.version,coalesce(q.order_id,'')
 from platform.idempotency_record i join sales.quotation q on q.tenant_id=i.tenant_id and q.quotation_id=i.resource_id
 join crm.lead l on l.tenant_id=q.tenant_id and l.organization_id=q.organization_id and l.lead_id=q.lead_id
 where i.tenant_id=$1 and i.scope='franchise-quote' and i.idempotency_key=$4 and i.status='completed' and i.resource_type='quotation' and i.response_code=201
 and i.response_body->>'quotation_id'=q.quotation_id and q.organization_id=$2 and q.lead_id=$3
 and q.customer_principal_id is not distinct from l.customer_principal_id`, tenant, organization, lead, key).Scan(&v.ID, &v.OrganizationID, &v.LeadID, &v.VariantID, &v.PriceBookID, &v.Currency, &v.TotalMinorUnits, &v.ValidUntil, &v.State, &v.Version, &v.OrderID)
	if errors.Is(err, pgx.ErrNoRows) {
		return franchisejourney.Quote{}, franchisejourney.ErrConflict
	}
	if err != nil {
		return franchisejourney.Quote{}, err
	}
	return v, nil
}

func (r *FranchiseJourney) CreateServiceResourceOnce(ctx context.Context, tenant, subject, key, hash string, value franchisejourney.ServiceResource, event string) (franchisejourney.ServiceResource, bool, error) {
	var empty franchisejourney.ServiceResource
	if subject == "" {
		return empty, false, franchisejourney.ErrInvalid
	}
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return empty, false, err
	}
	defer tx.Rollback(ctx)
	result, err := tx.Exec(ctx, `insert into platform.idempotency_record(tenant_id,scope,idempotency_key,request_sha256_hex,status,locked_until,expires_at)values($1,'franchise-resource',$2,$3,'processing',clock_timestamp()+interval '30 seconds',clock_timestamp()+interval '24 hours')on conflict do nothing`, tenant, key, hash)
	if err != nil {
		return empty, false, err
	}
	if result.RowsAffected() == 0 {
		replay, storedHash, err := readServiceResourceCreation(ctx, tx, tenant, value.OrganizationID, key)
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
	value, err = r.createServiceResourceTx(ctx, tx, tenant, subject, value, event)
	if err != nil {
		return empty, false, err
	}
	result, err = tx.Exec(ctx, `update platform.idempotency_record set status='completed',response_code=201,response_body=jsonb_build_object('resource_id',$3::text),resource_type='service-resource',resource_id=$3,locked_until=null where tenant_id=$1 and scope='franchise-resource' and idempotency_key=$2 and status='processing'`, tenant, key, value.ID)
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

type serviceResourceRowReader interface {
	QueryRow(context.Context, string, ...any) pgx.Row
}

func readServiceResourceCreation(ctx context.Context, q serviceResourceRowReader, tenant, organization, key string) (franchisejourney.ServiceResource, string, error) {
	var v franchisejourney.ServiceResource
	var hash string
	err := q.QueryRow(ctx, `select i.request_sha256_hex,a.resource_id,a.organization_id,coalesce(a.principal_subject,''),a.display_name,a.resource_kind,a.status,a.version,
array(select rs.appointment_kind from crm.resource_skill rs where rs.tenant_id=a.tenant_id and rs.resource_id=a.resource_id order by rs.appointment_kind)
from platform.idempotency_record i join crm.service_resource a on a.tenant_id=i.tenant_id and a.resource_id=i.resource_id
where i.tenant_id=$1 and i.scope='franchise-resource' and i.idempotency_key=$3 and i.status='completed' and i.resource_type='service-resource' and i.response_code=201 and i.response_body->>'resource_id'=a.resource_id and a.organization_id=$2`, tenant, organization, key).Scan(&hash, &v.ID, &v.OrganizationID, &v.PrincipalSubject, &v.DisplayName, &v.Kind, &v.Status, &v.Version, &v.Skills)
	if errors.Is(err, pgx.ErrNoRows) {
		return franchisejourney.ServiceResource{}, "", franchisejourney.ErrConflict
	}
	if err != nil {
		return franchisejourney.ServiceResource{}, "", err
	}
	return v, hash, nil
}
func (r *FranchiseJourney) ServiceResourceCreationResult(ctx context.Context, tenant, organization, key string) (franchisejourney.ServiceResource, error) {
	v, _, err := readServiceResourceCreation(ctx, r.pool, tenant, organization, key)
	return v, err
}

func (r *FranchiseJourney) CreateAppointmentSlotOnce(ctx context.Context, tenant, subject, key, hash string, value franchisejourney.AppointmentSlot, event string) (franchisejourney.AppointmentSlot, bool, error) {
	boundHash, bindErr := r.policy.BindRequestHash(hash)
	if bindErr != nil {
		return franchisejourney.AppointmentSlot{}, false, franchisejourney.ErrInvalid
	}
	hash = boundHash
	var empty franchisejourney.AppointmentSlot
	if subject == "" {
		return empty, false, franchisejourney.ErrInvalid
	}
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return empty, false, err
	}
	defer tx.Rollback(ctx)
	result, err := tx.Exec(ctx, `insert into platform.idempotency_record(tenant_id,scope,idempotency_key,request_sha256_hex,status,locked_until,expires_at)values($1,'franchise-slot',$2,$3,'processing',clock_timestamp()+interval '30 seconds',clock_timestamp()+interval '24 hours')on conflict do nothing`, tenant, key, hash)
	if err != nil {
		return empty, false, err
	}
	if result.RowsAffected() == 0 {
		replay, storedHash, err := readAppointmentSlotCreation(ctx, tx, tenant, value.OrganizationID, key)
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
	value, err = r.createAppointmentSlotTx(ctx, tx, tenant, subject, value, event)
	if err != nil {
		return empty, false, err
	}
	result, err = tx.Exec(ctx, `update platform.idempotency_record set status='completed',response_code=201,response_body=jsonb_build_object('slot_id',$3::text),resource_type='appointment-slot',resource_id=$3,locked_until=null where tenant_id=$1 and scope='franchise-slot' and idempotency_key=$2 and status='processing'`, tenant, key, value.ID)
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

type appointmentSlotRowReader interface {
	QueryRow(context.Context, string, ...any) pgx.Row
}

func readAppointmentSlotCreation(ctx context.Context, q appointmentSlotRowReader, tenant, organization, key string) (franchisejourney.AppointmentSlot, string, error) {
	var v franchisejourney.AppointmentSlot
	var hash string
	err := q.QueryRow(ctx, `select i.request_sha256_hex,a.slot_id,a.organization_id,a.appointment_kind,a.starts_at,a.ends_at,a.capacity,
(select count(*) from crm.appointment p where p.tenant_id=a.tenant_id and p.slot_id=a.slot_id and p.state in ('requested','confirmed')),a.state,a.version
from platform.idempotency_record i join crm.appointment_slot a on a.tenant_id=i.tenant_id and a.slot_id=i.resource_id
where i.tenant_id=$1 and i.scope='franchise-slot' and i.idempotency_key=$3 and i.status='completed' and i.resource_type='appointment-slot' and i.response_code=201 and i.response_body->>'slot_id'=a.slot_id and a.organization_id=$2`, tenant, organization, key).Scan(&hash, &v.ID, &v.OrganizationID, &v.Kind, &v.StartsAt, &v.EndsAt, &v.Capacity, &v.Booked, &v.State, &v.Version)
	if errors.Is(err, pgx.ErrNoRows) {
		return franchisejourney.AppointmentSlot{}, "", franchisejourney.ErrConflict
	}
	if err != nil {
		return franchisejourney.AppointmentSlot{}, "", err
	}
	return v, hash, nil
}
func (r *FranchiseJourney) AppointmentSlotCreationResult(ctx context.Context, tenant, organization, key string) (franchisejourney.AppointmentSlot, error) {
	v, _, err := readAppointmentSlotCreation(ctx, r.pool, tenant, organization, key)
	return v, err
}

func (r *FranchiseJourney) PublishedDeliveryChecklist(ctx context.Context, tenant, organization, id string, version int64) (franchisejourney.DeliveryChecklist, error) {
	var value franchisejourney.DeliveryChecklist
	err := r.pool.QueryRow(ctx, `select t.checklist_id,t.organization_id,t.checklist_version,t.title,t.state,
 coalesce((select jsonb_agg(jsonb_build_object('id',i.item_id,'ordinal',i.ordinal,'prompt',i.prompt,'response_type',i.response_type,'required',i.required) order by i.ordinal)
 from sales.delivery_checklist_item i where i.tenant_id=t.tenant_id and i.organization_id=t.organization_id and i.checklist_id=t.checklist_id and i.checklist_version=t.checklist_version),'[]'::jsonb)
 from sales.delivery_checklist_template t where t.tenant_id=$1 and t.organization_id=$2 and t.checklist_id=$3 and t.checklist_version=$4 and t.state='published'`, tenant, organization, id, version).Scan(&value.ID, &value.OrganizationID, &value.Version, &value.Title, &value.State, &value.Items)
	if errors.Is(err, pgx.ErrNoRows) {
		return franchisejourney.DeliveryChecklist{}, franchisejourney.ErrConflict
	}
	if err != nil {
		return franchisejourney.DeliveryChecklist{}, err
	}
	return value, nil
}

func (r *FranchiseJourney) DeliveryChecklistCompletion(ctx context.Context, tenant, organization, handover string) (franchisejourney.ChecklistCompletion, error) {
	var v franchisejourney.ChecklistCompletion
	err := r.pool.QueryRow(ctx, `select h.handover_id,h.organization_id,h.state,h.version,h.checklist_id,h.checklist_version,h.checklist_completed_at,h.checklist_completed_by_subject,
 coalesce((select jsonb_agg(jsonb_build_object('item_id',r.item_id,'response_text',r.response_text,'evidence_sha256',coalesce(r.evidence_sha256_hex,'')) order by r.item_id collate "C") from sales.delivery_checklist_response r where r.tenant_id=h.tenant_id and r.handover_id=h.handover_id),'[]'::jsonb)
 from sales.delivery_handover h
 join sales.customer_order o on o.tenant_id=h.tenant_id and o.order_id=h.order_id and o.organization_id=h.organization_id and o.customer_principal_id=h.customer_principal_id
 join inventory.stock_unit i on i.tenant_id=h.tenant_id and i.stock_unit_id=h.stock_unit_id and i.organization_id=h.organization_id
 join sales.delivery_checklist_template t on t.tenant_id=h.tenant_id and t.organization_id=h.organization_id and t.checklist_id=h.checklist_id and t.checklist_version=h.checklist_version and t.state='published'
 where h.tenant_id=$1 and h.organization_id=$2 and h.handover_id=$3 and h.checklist_completed_at is not null
 and not exists(select 1 from sales.delivery_checklist_response r where r.tenant_id=h.tenant_id and r.handover_id=h.handover_id and (r.organization_id,r.checklist_id,r.checklist_version,r.answered_by_subject) is distinct from (h.organization_id,h.checklist_id,h.checklist_version,h.checklist_completed_by_subject))`, tenant, organization, handover).Scan(&v.HandoverID, &v.OrganizationID, &v.State, &v.Version, &v.ChecklistID, &v.ChecklistVersion, &v.CompletedAt, &v.ActorSubject, &v.Responses)
	if errors.Is(err, pgx.ErrNoRows) {
		return franchisejourney.ChecklistCompletion{}, franchisejourney.ErrConflict
	}
	if err != nil {
		return franchisejourney.ChecklistCompletion{}, err
	}
	return v, nil
}

// Shared read invariant; the authorization alias is always scoped by tenant and organization.
const returnScopePredicate = `exists(select 1 from sales.delivery_handover h
 join sales.customer_order o on o.tenant_id=h.tenant_id and o.order_id=h.order_id and o.organization_id=h.organization_id and o.customer_principal_id=h.customer_principal_id
 join inventory.stock_unit i on i.tenant_id=h.tenant_id and i.stock_unit_id=h.stock_unit_id and i.organization_id=h.organization_id
 join sales.delivery_exception x on x.tenant_id=h.tenant_id and x.handover_id=h.handover_id and x.organization_id=h.organization_id and x.customer_principal_id=h.customer_principal_id
 where h.tenant_id=a.tenant_id and h.handover_id=a.handover_id and x.exception_id=a.exception_id
 and (h.organization_id,h.order_id,h.stock_unit_id,h.customer_principal_id)=(a.organization_id,a.order_id,a.stock_unit_id,a.customer_principal_id))`

func lockReturnAuthorizationScope(ctx context.Context, tx pgx.Tx, tenant, organization, authorization string) error {
	var id string
	err := tx.QueryRow(ctx, `select a.authorization_id from sales.return_authorization a
 join sales.delivery_handover h on h.tenant_id=a.tenant_id and h.handover_id=a.handover_id and (h.organization_id,h.order_id,h.stock_unit_id,h.customer_principal_id)=(a.organization_id,a.order_id,a.stock_unit_id,a.customer_principal_id)
 join sales.customer_order o on o.tenant_id=h.tenant_id and o.order_id=h.order_id and o.organization_id=h.organization_id and o.customer_principal_id=h.customer_principal_id
 join inventory.stock_unit i on i.tenant_id=h.tenant_id and i.stock_unit_id=h.stock_unit_id and i.organization_id=h.organization_id
 join sales.delivery_exception x on x.tenant_id=a.tenant_id and x.exception_id=a.exception_id and x.handover_id=h.handover_id and x.organization_id=h.organization_id and x.customer_principal_id=h.customer_principal_id
 where a.tenant_id=$1 and a.organization_id=$2 and a.authorization_id=$3 for update of h for share of a,x,o,i`, tenant, organization, authorization).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return franchisejourney.ErrConflict
	}
	return err
}

func (r *FranchiseJourney) ReturnCaseResult(ctx context.Context, tenant, organization, authorization string) (franchisejourney.ReturnCase, error) {
	var v franchisejourney.ReturnCase
	var scopeValid bool
	err := r.pool.QueryRow(ctx, `select a.authorization_id,a.organization_id,a.order_id,a.stock_unit_id,a.customer_principal_id,a.disposition,a.authorized_at,
 case when r.receipt_id is null then null else jsonb_build_object('id',r.receipt_id,'authorization_id',r.authorization_id,'organization_id',r.organization_id,'order_id',r.order_id,'stock_unit_id',r.stock_unit_id,'customer_subject',r.customer_principal_id,'received_serial_number',r.received_serial_number,'condition_code',r.condition_code,'notes',r.notes,'evidence_sha256',r.evidence_sha256_hex,'received_by_subject',r.received_by_subject,'received_at',r.received_at) end,
 case when d.disposition_id is null then null else jsonb_build_object('id',d.disposition_id,'receipt_id',d.receipt_id,'inventory_action',d.inventory_action,'customer_remedy',d.customer_remedy,'notes',d.notes,'decided_by_subject',d.decided_by_subject,'decided_at',d.decided_at,'effect_requests',coalesce((select jsonb_agg(jsonb_build_object('id',e.request_id,'effect_kind',e.effect_kind,'owner_context',e.owner_context,'state',e.state,'idempotency_key',e.idempotency_key,'requested_at',e.requested_at) order by e.effect_kind) from sales.return_effect_request e where e.tenant_id=d.tenant_id and e.disposition_id=d.disposition_id),'[]'::jsonb)) end,
 `+returnScopePredicate+` and (r.receipt_id is null or (r.organization_id,r.order_id,r.stock_unit_id,r.customer_principal_id)=(a.organization_id,a.order_id,a.stock_unit_id,a.customer_principal_id))
 from sales.return_authorization a left join sales.return_receipt r on r.tenant_id=a.tenant_id and r.authorization_id=a.authorization_id left join sales.return_disposition d on d.tenant_id=r.tenant_id and d.receipt_id=r.receipt_id
 where a.tenant_id=$1 and a.organization_id=$2 and a.authorization_id=$3`, tenant, organization, authorization).Scan(&v.AuthorizationID, &v.OrganizationID, &v.OrderID, &v.StockUnitID, &v.CustomerSubject, &v.AuthorizedAction, &v.AuthorizedAt, &v.Receipt, &v.Disposition, &scopeValid)
	if errors.Is(err, pgx.ErrNoRows) {
		return franchisejourney.ReturnCase{}, franchisejourney.ErrConflict
	}
	if err != nil {
		return franchisejourney.ReturnCase{}, err
	}
	if !scopeValid {
		return franchisejourney.ReturnCase{}, franchisejourney.ErrConflict
	}
	return v, nil
}
