package whatsappbridge

// AUTHORED composition of the existing job claim generation and send fence.
import (
	"context"
	"elite.local/enterprise/internal/approval"
	"elite.local/enterprise/internal/channels"
	"elite.local/enterprise/internal/outbounddelivery"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/postgres"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5"
	"sync/atomic"
	"time"
)

func (s *ScheduleApprovals) admit(ctx context.Context, tx pgx.Tx, c ScheduleContext, h string) error {
	var stored string
	e := tx.QueryRow(ctx, `select request_sha256 from communication.whatsapp_schedule where tenant_id=$1 and delivery_key=$2 and organization_id=$3 for update`, s.tenant, c.DeliveryKey, s.org).Scan(&stored)
	if errors.Is(e, pgx.ErrNoRows) {
		return ErrApproval
	}
	if e != nil {
		return e
	}
	if stored != h {
		return ErrApproval
	}
	if e = s.guard(ctx, tx, c); e != nil {
		return e
	}
	if _, e = postgres.LookupHumanApproval(ctx, tx, s.tenant, c.DeliveryKey, s.org, approval.KindWhatsAppSchedule, h); e != nil {
		return e
	}
	var now time.Time
	if e = tx.QueryRow(ctx, `select clock_timestamp()`).Scan(&now); e != nil {
		return e
	}
	if c.NotBefore.After(now) || !c.ExpiresAt.After(now) {
		return ErrApproval
	}
	if c.Request.CampaignID != "" {
		if s.extension == nil {
			return ErrApproval
		}
		return s.extension.Admit(ctx, tx, c)
	}
	var ambiguous bool
	e = tx.QueryRow(ctx, `select exists(select 1 from communication.whatsapp_schedule s
 join communication.outbound_delivery d on d.tenant_id=s.tenant_id and d.delivery_key=s.delivery_key and d.channel_code='whatsapp'
 where s.tenant_id=$1 and s.appointment_id=$2 and s.delivery_key<>$3 and d.state in('unknown','sending'))`, s.tenant, c.AppointmentID, c.DeliveryKey).Scan(&ambiguous)
	if e != nil {
		return e
	}
	if ambiguous {
		return outbounddelivery.ErrUnknown
	}
	return nil
}

type scheduledBound struct {
	*postgres.OutboundDeliveryStore
	approvals *ScheduleApprovals
	c         ScheduleContext
	hash      string
}

func (b *scheduledBound) Claim(ctx context.Context, m channels.Message, h string) (outbounddelivery.Claim, error) {
	if m.DeliveryKey != b.c.DeliveryKey || h != b.c.MessageSHA256 {
		return outbounddelivery.Claim{}, ErrApproval
	}
	return b.OutboundDeliveryStore.ClaimWithDomainAdmission(ctx, m, h, func(ctx context.Context, tx pgx.Tx) error { return b.approvals.admit(ctx, tx, b.c, b.hash) })
}

type ScheduledNotifications struct {
	turn      atomic.Uint32
	approvals *ScheduleApprovals
	sender    *Sender
	store     *postgres.OutboundDeliveryStore
	receiver  channels.Channel
	jobs      *postgres.Jobs
	worker    string
}

func NewScheduledNotifications(s *ScheduleApprovals, sender *Sender, store *postgres.OutboundDeliveryStore, receiver channels.Channel, worker string) (*ScheduledNotifications, error) {
	if s == nil || sender == nil || sender.TenantID != s.tenant || digest(sender.Profile) != digest(s.profile) || sender.Tokens == nil || store == nil || receiver == nil || receiver.Code() != "whatsapp" || !webhookConnection.MatchString(worker) {
		return nil, ErrApproval
	}
	frozen := *sender
	frozen.Profile = append(json.RawMessage(nil), sender.Profile...)
	frozen.Approvals = s
	return &ScheduledNotifications{approvals: s, sender: &frozen, store: store, receiver: receiver, jobs: postgres.NewJobs(s.base.pool), worker: worker}, nil
}

type ScheduledWorkResult struct {
	Claimed     bool   `json:"claimed"`
	DeliveryKey string `json:"delivery_key"`
	Outcome     string `json:"outcome"`
}

func (m *ScheduledNotifications) ProcessOnce(ctx context.Context, p identity.Principal) (ScheduledWorkResult, error) {
	var out ScheduledWorkResult
	if m == nil || !m.approvals.authorized(p, "notification:dispatch") {
		return out, ErrApproval
	}
	s := m.approvals
	ctx, cancel := context.WithTimeout(ctx, 45*time.Second)
	defer cancel()
	match, _ := json.Marshal(map[string]string{"connection_id": s.connection, "profile_sha256": digest(s.profile)})
	kinds := []string{"whatsapp.scheduled.v1"}
	if s.extension != nil {
		kinds = append(kinds, "whatsapp.campaign.scheduled.v1")
		if m.turn.Add(1)%2 == 0 {
			kinds[0], kinds[1] = kinds[1], kinds[0]
		}
	}
	var jobs []postgres.Job
	var e error
	for _, kind := range kinds {
		if expired, err := m.quarantineExpired(ctx, match, kind); err != nil || expired.Claimed {
			return expired, err
		}
		jobs, e = m.jobs.ClaimScoped(ctx, "whatsapp-scheduled", m.worker, 2*time.Minute, 1, postgres.JobScope{TenantID: s.tenant, JobType: kind, SchemaVersion: 1, PayloadMatch: match})
		if e != nil {
			return out, e
		}
		if len(jobs) > 0 {
			break
		}
	}
	if len(jobs) == 0 {
		return out, nil
	}

	job := jobs[0]
	out.Claimed = true
	var payload struct {
		DeliveryKey   string `json:"delivery_key"`
		ConnectionID  string `json:"connection_id"`
		ProfileSHA256 string `json:"profile_sha256"`
	}
	if json.Unmarshal(job.Payload, &payload) != nil || payload.ConnectionID != s.connection || payload.ProfileSHA256 != digest(s.profile) {
		_, e = m.jobs.Fail(ctx, job.TenantID, job.JobID, m.worker, "SCHEDULE_JOB_BINDING", job.Attempts, time.Second)
		return out, e
	}
	out.DeliveryKey = payload.DeliveryKey
	c, h, e := s.request(ctx, payload.DeliveryKey)
	if e != nil {
		_, failed := m.jobs.Fail(ctx, job.TenantID, job.JobID, m.worker, "SCHEDULE_REQUEST_UNAVAILABLE", job.Attempts, time.Second)
		if failed != nil {
			return out, failed
		}
		return out, e
	}
	var same bool
	e = s.base.pool.QueryRow(ctx, `select exists(select 1 from communication.whatsapp_schedule where tenant_id=$1 and delivery_key=$2 and job_id=$3 and request_sha256=$4)`, s.tenant, c.DeliveryKey, job.JobID, h).Scan(&same)
	if e != nil || !same {
		_, failed := m.jobs.Fail(ctx, job.TenantID, job.JobID, m.worker, "SCHEDULE_JOB_DIVERGENT", job.Attempts, time.Second)
		if failed != nil {
			return out, failed
		}
		return out, ErrApproval
	}
	bound := &scheduledBound{OutboundDeliveryStore: m.store, approvals: s, c: c, hash: h}
	channel := &outbounddelivery.Channel{CodeValue: "whatsapp", Receiver: m.receiver, Sender: m.sender, Store: bound}
	sendErr := channel.Send(ctx, c.Message)
	var state string
	e = s.base.pool.QueryRow(ctx, `select coalesce((select state from communication.outbound_delivery where tenant_id=$1 and delivery_key=$2 and channel_code='whatsapp'),'')`, s.tenant, c.DeliveryKey).Scan(&state)
	if e != nil {
		return out, e
	}
	switch state {
	case "accepted":
		out.Outcome = "ACCEPTED"
	case "failed_terminal":
		out.Outcome = "FAILED_TERMINAL"
	case "unknown", "sending":
		out.Outcome = "RECONCILIATION_REQUIRED"
	default:
		if errors.Is(sendErr, ErrApproval) || errors.Is(sendErr, approval.ErrInvalidRequest) || errors.Is(sendErr, approval.ErrSeparation) {
			out.Outcome = "SUPPRESSED"
		} else {
			_, e = m.jobs.Fail(ctx, job.TenantID, job.JobID, m.worker, "SCHEDULE_EFFECT_UNCONFIRMED", job.Attempts, time.Second)
			if e != nil {
				return out, e
			}
			return out, sendErr
		}
	}
	tx, e := s.base.pool.Begin(ctx)
	if e != nil {
		return out, e
	}
	defer tx.Rollback(ctx)
	_, e = tx.Exec(ctx, `insert into communication.whatsapp_schedule_result(tenant_id,delivery_key,job_id,outcome)values($1,$2,$3,$4)on conflict do nothing`, s.tenant, c.DeliveryKey, job.JobID, out.Outcome)
	if e != nil {
		return out, e
	}
	if e = postgres.CompleteJobInTx(ctx, tx, job, m.worker); e != nil {
		return out, e
	}
	return out, tx.Commit(ctx)
}
