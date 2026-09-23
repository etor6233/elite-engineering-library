package whatsappbridge

// AUTHORED transaction wiring. Approval decisions and queue ownership remain
// in their original owners; this module binds their exact source and timing.
import (
	"context"
	"elite.local/enterprise/internal/approval"
	"elite.local/enterprise/internal/platform/identity"
	"encoding/json"
	"github.com/jackc/pgx/v5"
	"time"
)

func (s *ScheduleApprovals) Decide(ctx context.Context, p identity.Principal, key, hash string, yes bool, reason string) (approval.State, error) {
	if !s.authorized(p, "notification:approve") {
		return "", ErrApproval
	}
	c, h, e := s.request(ctx, key)
	if e != nil || h != hash {
		return "", ErrApproval
	}
	guard := func(ctx context.Context, tx pgx.Tx) error {
		if !yes {
			return nil
		}
		if e := s.guard(ctx, tx, c); e != nil {
			return e
		}
		var now time.Time
		if e := tx.QueryRow(ctx, `select clock_timestamp()`).Scan(&now); e != nil || !c.NotBefore.After(now) {
			return ErrApproval
		}
		var e error
		if c.Request.CampaignID != "" {
			if s.extension == nil {
				return ErrApproval
			}
			e = s.extension.Enqueue(ctx, tx, c, h)
			if e != nil {
				return e
			}
		} else {
			_, e = tx.Exec(ctx, `insert into communication.whatsapp_schedule(tenant_id,delivery_key,job_id,organization_id,appointment_id,appointment_version,not_before,expires_at,request_sha256)
 values($1,$2,gen_random_uuid(),$3,$4,$5,$6,$7,$8) on conflict(tenant_id,delivery_key) do nothing`, s.tenant, key, s.org, c.AppointmentID, c.AppointmentVersion, c.NotBefore, c.ExpiresAt, h)
			if e != nil {
				return e
			}
		}
		jobType := "whatsapp.scheduled.v1"
		if c.Request.CampaignID != "" {
			jobType = "whatsapp.campaign.scheduled.v1"
		}
		_, e = tx.Exec(ctx, `insert into platform.job(tenant_id,job_id,queue,job_type,schema_version,payload,available_at,max_attempts)
 select tenant_id,job_id,'whatsapp-scheduled',$6::text,1,
 jsonb_build_object('delivery_key',delivery_key,'connection_id',$3::text,'profile_sha256',$4::text),
 not_before,3 from communication.whatsapp_schedule
 where tenant_id=$1 and delivery_key=$2 and request_sha256=$5 on conflict(tenant_id,job_id) do nothing`, s.tenant, key, s.connection, digest(s.profile), h, jobType)
		return e
	}
	return s.human.Decide(ctx, p, s.tenant, key, s.org, hash, yes, reason, "notification:approve", guard)
}
func (s *ScheduleApprovals) Cancel(ctx context.Context, p identity.Principal, key, hash, reason string) error {
	if !s.authorized(p, "notification:cancel") || reason == "" || len(reason) > 2048 {
		return ErrApproval
	}
	_, h, e := s.request(ctx, key)
	if e != nil || h != hash {
		return ErrApproval
	}
	tx, e := s.base.pool.Begin(ctx)
	if e != nil {
		return e
	}
	defer tx.Rollback(ctx)
	var found string
	e = tx.QueryRow(ctx, `select request_sha256 from communication.whatsapp_schedule where tenant_id=$1 and delivery_key=$2 for update`, s.tenant, key).Scan(&found)
	if e != nil || found != h {
		return ErrApproval
	}
	var effect bool
	e = tx.QueryRow(ctx, `select exists(select 1 from communication.outbound_delivery where tenant_id=$1 and delivery_key=$2 and channel_code='whatsapp' and state in('sending','unknown','accepted'))`, s.tenant, key).Scan(&effect)
	if e != nil || effect {
		return ErrApproval
	}
	_, e = tx.Exec(ctx, `insert into communication.whatsapp_schedule_cancellation(tenant_id,delivery_key,actor_subject,reason,request_sha256)values($1,$2,$3,$4,$5)on conflict do nothing`, s.tenant, key, p.Subject, reason, h)
	if e != nil {
		return e
	}
	var same bool
	e = tx.QueryRow(ctx, `select actor_subject=$3 and reason=$4 and request_sha256=$5 from communication.whatsapp_schedule_cancellation where tenant_id=$1 and delivery_key=$2`, s.tenant, key, p.Subject, reason, h).Scan(&same)
	if e != nil || !same {
		return ErrApproval
	}
	return tx.Commit(ctx)
}

type ScheduleStatus struct {
	Notification  *NotificationStatus `json:"notification,omitempty"`
	Context       ScheduleContext     `json:"context"`
	RequestSHA256 string              `json:"request_sha256"`
	ApprovalState string              `json:"approval_state"`
	DeliveryState string              `json:"delivery_state"`
	Cancelled     bool                `json:"cancelled"`
	Outcome       string              `json:"outcome"`
	Attempts      int                 `json:"attempts"`
	JobCompleted  bool                `json:"job_completed"`
	TerminalError string              `json:"terminal_error"`
}

func (s *ScheduleApprovals) Status(ctx context.Context, p identity.Principal, key string) (ScheduleStatus, error) {
	var out ScheduleStatus
	if !s.authorized(p, "notification:read") {
		return out, ErrApproval
	}
	c, h, e := s.request(ctx, key)
	if e != nil {
		return out, e
	}
	out.Context = c
	out.RequestSHA256 = h
	e = s.base.pool.QueryRow(ctx, `select r.state,coalesce(d.state,''),cancel.delivery_key is not null,coalesce(result.outcome,''),
 coalesce(j.attempts,0),j.completed_at is not null,coalesce(j.terminal_error_code,'')
 from approval.request r left join communication.whatsapp_schedule s on s.tenant_id=r.tenant_id and s.delivery_key=r.request_id
 left join communication.whatsapp_schedule_cancellation cancel on cancel.tenant_id=s.tenant_id and cancel.delivery_key=s.delivery_key
 left join communication.whatsapp_schedule_result result on result.tenant_id=s.tenant_id and result.delivery_key=s.delivery_key
 left join platform.job j on j.tenant_id=s.tenant_id and j.job_id=s.job_id
 left join communication.outbound_delivery d on d.tenant_id=r.tenant_id and d.delivery_key=r.request_id and d.channel_code='whatsapp'
 where r.tenant_id=$1 and r.request_id=$2 and r.organization_id=$3 and r.kind='whatsapp_schedule'`, s.tenant, key, s.org).Scan(&out.ApprovalState, &out.DeliveryState, &out.Cancelled, &out.Outcome, &out.Attempts, &out.JobCompleted, &out.TerminalError)
	if e != nil {
		return out, e
	}
	if out.ApprovalState == "approved" {
		status, err := readNotificationStatus(ctx, s.base.pool, p, s.org, "", key)
		if err != nil {
			return out, err
		}
		out.Notification = &status
	}
	return out, nil
}
func (s *ScheduleApprovals) ResolveWhatsAppApproval(ctx context.Context, tenant, key string) (Approval, error) {
	var out Approval
	if tenant != s.tenant {
		return out, ErrApproval
	}
	c, h, e := s.request(ctx, key)
	if e != nil {
		return out, e
	}
	tx, e := s.base.pool.Begin(ctx)
	if e != nil {
		return out, e
	}
	defer tx.Rollback(ctx)
	if e = s.admit(ctx, tx, c, h); e != nil {
		return out, e
	}
	e = tx.QueryRow(ctx, `select d.reviewer from approval.request r join approval.decision d on d.tenant_id=r.tenant_id and d.request_id=r.request_id and d.approved and d.reviewer<>r.requester where r.tenant_id=$1 and r.request_id=$2 and r.state='approved'`, tenant, key).Scan(&out.ApprovedBy)
	if e != nil {
		return out, e
	}
	out.MessageSHA256 = c.MessageSHA256
	out.ProfileSHA256 = c.ProfileSHA256
	out.EvidenceSHA256 = h
	out.NotBefore = c.NotBefore
	out.ExpiresAt = c.ExpiresAt
	return out, tx.Commit(ctx)
}

// Keep the canonical approval payload check available to the bounded HTTP
// wrapper; unknown/duplicate fields are rejected by its strict typed decoder.
func scheduleCanonical(v any) ([]byte, string, error) {
	raw, e := json.Marshal(v)
	if e != nil {
		return nil, "", e
	}
	return approval.CanonicalPayload(raw)
}
