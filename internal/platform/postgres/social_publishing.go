package postgres

// AUTHORED binding glue. Request/decision, scheduling, single-write fence and
// public provider observations retain their existing shared persistence owners.
import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"elite.local/enterprise/internal/approval"
	"elite.local/enterprise/internal/outbounddelivery"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/socialbridge"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SocialPublishing struct {
	pool      *pgxpool.Pool
	profile   *socialbridge.Profile
	approvals *HumanApprovals
	jobs      *Jobs
	fence     *OutboundDeliveryStore
	adapter   socialbridge.Adapter
}
type socialJob struct {
	ApprovalID    string `json:"approval_id"`
	ProfileSHA256 string `json:"profile_sha256"`
	PageID        string `json:"page_id"`
	RequestSHA256 string `json:"request_sha256"`
}

func NewSocialPublishing(pool *pgxpool.Pool, p *socialbridge.Profile, key []byte, a socialbridge.Adapter) (*SocialPublishing, error) {
	if p == nil || a == nil {
		return nil, socialbridge.ErrBinding
	}
	f, e := NewOutboundDeliveryStore(pool, key, time.Duration(p.Config().LeaseSeconds)*time.Second)
	if e != nil {
		return nil, e
	}
	return &SocialPublishing{pool, p, NewHumanApprovals(pool), NewJobs(pool), f, a}, nil
}
func socialKind(r socialbridge.Request) approval.Kind {
	if r.Intent.Operation == "revoke" {
		return approval.KindSocialRevoke
	}
	return approval.KindSocialPublish
}
func socialJobID(tenant, id string) string {
	return uuid.NewSHA1(uuid.NameSpaceOID, []byte("elite.social.job\x00"+tenant+"\x00"+id)).String()
}
func (s *SocialPublishing) request(ctx context.Context, id string) (socialbridge.Request, string, error) {
	var raw []byte
	var hash string
	c := s.profile.Config()
	e := s.pool.QueryRow(ctx, `select payload,evidence_sha from approval.request where tenant_id=$1 and request_id=$2 and organization_id=$3 and kind in ('social_publish','social_revoke')`, c.TenantID, id, c.OrganizationID).Scan(&raw, &hash)
	var r socialbridge.Request
	if e != nil {
		return r, "", e
	}
	if socialbridge.Decode(raw, &r) != nil || s.profile.Validate(r) != nil || r.Intent.ApprovalID != id {
		return r, "", socialbridge.ErrBinding
	}
	return r, hash, nil
}
func (s *SocialPublishing) guard(ctx context.Context, tx pgx.Tx, r socialbridge.Request, needDue bool) error {
	if s.profile.Validate(r) != nil {
		return socialbridge.ErrBinding
	}
	c := s.profile.Config()
	var valid bool
	e := tx.QueryRow(ctx, `with locked as materialized(select status from platform.tenant where tenant_id=$1 for share)
 select status='active' and $2::timestamptz>clock_timestamp() and (not $3::boolean or $4::timestamptz<=clock_timestamp()) from locked`, c.TenantID, r.ExpiresAt, needDue, r.ScheduledAt).Scan(&valid)
	if e != nil || !valid {
		return socialbridge.ErrBinding
	}
	if r.Intent.Operation == "revoke" {
		// The exact publish receipt belongs to this page, original approved content
		// and tenant. A user-supplied post id alone never authorizes a DELETE.
		var raw []byte
		e = tx.QueryRow(ctx, `select payload from platform.outbox_event where tenant_id=$1 and aggregate_type='social_approval' and aggregate_id=$2 and event_type='social.receipt' and aggregate_version=1 for share`, c.TenantID, r.OriginalApprovalID).Scan(&raw)
		var receipt socialbridge.Receipt
		if e != nil || socialbridge.Decode(raw, &receipt) != nil || receipt.State != "published" || receipt.PageID != c.PageID || receipt.TenantID != c.TenantID || receipt.ProfileSHA256 != s.profile.SHA256() || receipt.ContentSHA256 != r.Intent.ContentSHA256 || receipt.ProviderReference != r.ProviderReference {
			return socialbridge.ErrBinding
		}
	}
	return nil
}
func (s *SocialPublishing) Submit(ctx context.Context, p identity.Principal, r socialbridge.Request) (string, bool, error) {
	if !s.profile.Authorize(p, "social:request") || s.profile.Validate(r) != nil {
		return "", false, socialbridge.ErrBinding
	}
	raw, _ := json.Marshal(r)
	payload, hash, e := approval.CanonicalPayload(raw)
	if e != nil {
		return "", false, e
	}
	spec := HumanApprovalSpec{Request: approval.Request{TenantID: p.TenantID, ID: r.Intent.ApprovalID, Kind: socialKind(r), SubjectID: r.Intent.PageID, Requester: p.Subject, EvidenceSHA: hash}, OrganizationID: s.profile.Config().OrganizationID, Payload: payload}
	replay, e := s.approvals.Submit(ctx, p, spec, "social:request", func(ctx context.Context, tx pgx.Tx) error { return s.guard(ctx, tx, r, false) })
	return hash, replay, e
}
func (s *SocialPublishing) Decide(ctx context.Context, p identity.Principal, id, hash string, approve bool, reason string) (approval.State, error) {
	if !s.profile.Authorize(p, "social:approve") {
		return "", socialbridge.ErrBinding
	}
	r, stored, e := s.request(ctx, id)
	if e != nil || stored != hash {
		return "", socialbridge.ErrBinding
	}
	c := s.profile.Config()
	return s.approvals.Decide(ctx, p, c.TenantID, id, c.OrganizationID, hash, approve, reason, "social:approve", func(ctx context.Context, tx pgx.Tx) error {
		if !approve {
			return nil
		}
		if e := s.guard(ctx, tx, r, false); e != nil {
			return e
		}
		payload, _ := json.Marshal(socialJob{id, s.profile.SHA256(), c.PageID, hash})
		_, e := tx.Exec(ctx, `insert into platform.job(tenant_id,job_id,queue,job_type,schema_version,payload,max_attempts,available_at)values($1,$2,$3,'social.publish',1,$4,$5,$6)`, c.TenantID, socialJobID(c.TenantID, id), c.Queue, payload, c.MaxAttempts, r.ScheduledAt)
		return e
	})
}

type SocialStatus struct {
	ApprovalID    string          `json:"approval_id"`
	State         string          `json:"approval_state"`
	RequestSHA256 string          `json:"request_sha256"`
	Payload       json.RawMessage `json:"payload"`
	DeliveryState string          `json:"delivery_state"`
	Receipt       json.RawMessage `json:"receipt,omitempty"`
	TerminalCode  string          `json:"terminal_code,omitempty"`
}

func (s *SocialPublishing) Status(ctx context.Context, p identity.Principal, id string) (SocialStatus, error) {
	var v SocialStatus
	if !s.profile.Authorize(p, "social:read") {
		return v, socialbridge.ErrBinding
	}
	c := s.profile.Config()
	e := s.pool.QueryRow(ctx, `select a.request_id,a.state,a.evidence_sha,a.payload,coalesce(d.state,''),e.payload,coalesce(j.terminal_error_code,'') from approval.request a
 left join communication.outbound_delivery d on d.tenant_id=a.tenant_id and d.channel_code='facebook_pages' and d.delivery_key=a.payload->'intent'->>'delivery_key'
 left join platform.outbox_event e on e.tenant_id=a.tenant_id and e.aggregate_type='social_approval' and e.aggregate_id=a.request_id and e.event_type='social.receipt' and e.aggregate_version=1
 left join platform.job j on j.tenant_id=a.tenant_id and j.job_id=$4
 where a.tenant_id=$1 and a.request_id=$2 and a.organization_id=$3 and a.kind in ('social_publish','social_revoke')`, c.TenantID, id, c.OrganizationID, socialJobID(c.TenantID, id)).Scan(&v.ApprovalID, &v.State, &v.RequestSHA256, &v.Payload, &v.DeliveryState, &v.Receipt, &v.TerminalCode)
	if e != nil {
		return SocialStatus{}, e
	}
	if _, e = s.verifiedWorkspaceRequest(v.ApprovalID, v.Payload, v.RequestSHA256); e != nil {
		return SocialStatus{}, e
	}
	var request socialbridge.Request
	if socialbridge.Decode(v.Payload, &request) != nil {
		return SocialStatus{}, socialbridge.ErrBinding
	}
	if len(v.Receipt) != 0 {
		var receipt socialbridge.Receipt
		if socialbridge.Decode(v.Receipt, &receipt) != nil || receipt.Validate(request.Intent) != nil {
			return SocialStatus{}, socialbridge.ErrBinding
		}
	}
	if v.DeliveryState == "accepted" && len(v.Receipt) == 0 {
		return SocialStatus{}, socialbridge.ErrBinding
	}
	return v, nil
}
func (s *SocialPublishing) Claim(ctx context.Context, worker string) ([]Job, error) {
	c := s.profile.Config()
	match, _ := json.Marshal(map[string]string{"profile_sha256": s.profile.SHA256(), "page_id": c.PageID})
	return s.jobs.ClaimScoped(ctx, c.Queue, worker, time.Duration(c.LeaseSeconds)*time.Second, 1, JobScope{TenantID: c.TenantID, JobType: "social.publish", SchemaVersion: 1, PayloadMatch: match})
}
func (s *SocialPublishing) event(ctx context.Context, r socialbridge.Request, kind string, payload any) error {
	raw, _ := json.Marshal(payload)
	id := uuid.NewSHA1(uuid.NameSpaceOID, []byte("elite.social.event\x00"+r.Intent.TenantID+"\x00"+r.Intent.ApprovalID+"\x00"+kind)).String()
	_, e := s.pool.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload)values($1,$2,'social_approval',$3,1,$4,1,clock_timestamp(),$5) on conflict do nothing`, r.Intent.TenantID, id, r.Intent.ApprovalID, kind, raw)
	if e != nil {
		return e
	}
	var same bool
	e = s.pool.QueryRow(ctx, `select payload=$3::jsonb from platform.outbox_event where tenant_id=$1 and event_id=$2`, r.Intent.TenantID, id, raw).Scan(&same)
	if e != nil {
		return e
	}
	if !same {
		return socialbridge.ErrBinding
	}
	return nil
}
func (s *SocialPublishing) observation(ctx context.Context, r socialbridge.Request) (string, *socialbridge.Receipt, error) {
	var raw []byte
	e := s.pool.QueryRow(ctx, `select payload from platform.outbox_event where tenant_id=$1 and aggregate_type='social_approval' and aggregate_id=$2 and event_type='social.receipt' and aggregate_version=1`, r.Intent.TenantID, r.Intent.ApprovalID).Scan(&raw)
	if e == nil {
		var v socialbridge.Receipt
		if socialbridge.Decode(raw, &v) != nil || v.Validate(r.Intent) != nil {
			return "", nil, socialbridge.ErrBinding
		}
		return v.ProviderReference, &v, nil
	}
	if !errors.Is(e, pgx.ErrNoRows) {
		return "", nil, e
	}
	var ref string
	e = s.pool.QueryRow(ctx, `select payload->>'provider_reference' from platform.outbox_event where tenant_id=$1 and aggregate_type='social_approval' and aggregate_id=$2 and event_type='social.observation' and aggregate_version=1`, r.Intent.TenantID, r.Intent.ApprovalID).Scan(&ref)
	if errors.Is(e, pgx.ErrNoRows) {
		return "", nil, nil
	}
	return ref, nil, e
}
func (s *SocialPublishing) Process(ctx context.Context, job Job, worker string) error {
	var payload socialJob
	c := s.profile.Config()
	if job.TenantID != c.TenantID || job.JobType != "social.publish" || job.Queue != c.Queue || job.SchemaVersion != 1 || socialbridge.Decode(job.Payload, &payload) != nil || payload.ProfileSHA256 != s.profile.SHA256() || payload.PageID != c.PageID || job.JobID != socialJobID(c.TenantID, payload.ApprovalID) {
		return socialbridge.ErrBinding
	}
	r, hash, e := s.request(ctx, payload.ApprovalID)
	if e != nil || hash != payload.RequestSHA256 {
		return socialbridge.ErrBinding
	}
	message := r.Message()
	messageHash, e := outbounddelivery.MessageSHA256(message)
	if e != nil {
		return e
	}
	claim, e := s.fence.claimWithAdmission(ctx, message, messageHash, func(ctx context.Context, tx pgx.Tx) error {
		if _, e := LookupHumanApproval(ctx, tx, c.TenantID, payload.ApprovalID, c.OrganizationID, socialKind(r), hash); e != nil {
			return e
		}
		if e := s.guard(ctx, tx, r, true); e != nil {
			return e
		}
		// An obsolete job generation cannot acquire the remote write fence.
		var live bool
		e := tx.QueryRow(ctx, `with locked as materialized(select claimed_by,attempts,claimed_until,completed_at,terminal_error_code,payload from platform.job where tenant_id=$1 and job_id=$2 for update)
 select claimed_by=$3 and attempts=$4 and claimed_until>clock_timestamp() and completed_at is null and terminal_error_code is null and payload=$5::jsonb from locked`, job.TenantID, job.JobID, worker, job.Attempts, job.Payload).Scan(&live)
		if e != nil || !live {
			return ErrJobClaimLost
		}
		return nil
	})
	if e == nil && claim.Replay {
		return s.jobs.Complete(ctx, job.TenantID, job.JobID, worker, job.Attempts)
	}
	unknown := errors.Is(e, outbounddelivery.ErrUnknown)
	if e != nil && !unknown {
		return e
	}
	ref, receipt, e := s.observation(ctx, r)
	if e != nil {
		return e
	}
	if receipt == nil {
		// Every ambiguous delivery is read-only from here onward. No reference means
		// there is no safe GET identity; retain the case for provider evidence.
		if unknown && ref == "" {
			_, _ = s.jobs.Fail(ctx, job.TenantID, job.JobID, worker, "PROVIDER_REFERENCE_REQUIRED", job.Attempts, time.Duration(c.RetrySeconds)*time.Second)
			return socialbridge.ErrUnknown
		}
		call := r
		if unknown {
			call.ProviderReference = ref
		}
		result, callErr := s.adapter.Execute(ctx, call, unknown)
		if callErr != nil || result.Receipt == nil || result.Receipt.Validate(r.Intent) != nil {
			if result.ProviderReference != "" && socialbridge.ValidReference(c.PageID, result.ProviderReference) {
				if e = s.event(ctx, r, "social.observation", map[string]string{"provider_reference": result.ProviderReference, "request_sha256": hash}); e != nil {
					return e
				}
			}
			if !unknown {
				if e = s.fence.MarkUnknown(ctx, message, messageHash, "PROVIDER_OUTCOME_UNCONFIRMED"); e != nil {
					return e
				}
			}
			_, _ = s.jobs.Fail(ctx, job.TenantID, job.JobID, worker, "PROVIDER_OUTCOME_UNCONFIRMED", job.Attempts, time.Duration(c.RetrySeconds)*time.Second)
			return socialbridge.ErrUnknown
		}
		receipt = result.Receipt
		if e = s.event(ctx, r, "social.receipt", receipt); e != nil {
			return e
		}
	}
	durable := outbounddelivery.Receipt{ProviderMessageID: receipt.ProviderReference, EvidenceSHA256: receipt.EvidenceSHA256, AcceptedAt: time.Now().UTC()}
	if unknown {
		e = s.fence.ReconcileAccepted(ctx, message, messageHash, durable)
	} else {
		e = s.fence.Complete(ctx, message, messageHash, durable)
	}
	if e != nil {
		return e
	}
	return s.jobs.Complete(ctx, job.TenantID, job.JobID, worker, job.Attempts)
}

// Reconcile is an authorized GET-only recovery path even after a job's retry
// budget is exhausted. It never resets a generation or reissues POST/DELETE.
func (s *SocialPublishing) Reconcile(ctx context.Context, p identity.Principal, id string) error {
	if !s.profile.Authorize(p, "social:reconcile") {
		return socialbridge.ErrBinding
	}
	r, _, e := s.request(ctx, id)
	if e != nil {
		return e
	}
	message := r.Message()
	hash, _ := outbounddelivery.MessageSHA256(message)
	var state string
	e = s.pool.QueryRow(ctx, `select state from communication.outbound_delivery where tenant_id=$1 and channel_code='facebook_pages' and delivery_key=$2`, r.Intent.TenantID, r.Intent.DeliveryKey).Scan(&state)
	if e != nil {
		return e
	}
	if state != "accepted" {
		_, e = s.fence.Claim(ctx, message, hash)
		if !errors.Is(e, outbounddelivery.ErrUnknown) {
			return socialbridge.ErrUnknown
		}
		ref, receipt, e := s.observation(ctx, r)
		if e != nil {
			return e
		}
		if receipt == nil {
			if ref == "" {
				return socialbridge.ErrUnknown
			}
			call := r
			call.ProviderReference = ref
			result, e := s.adapter.Execute(ctx, call, true)
			if e != nil || result.Receipt == nil || result.Receipt.Validate(r.Intent) != nil {
				return socialbridge.ErrUnknown
			}
			receipt = result.Receipt
			if e = s.event(ctx, r, "social.receipt", receipt); e != nil {
				return e
			}
		}
		if e = s.fence.ReconcileAccepted(ctx, message, hash, outbounddelivery.Receipt{ProviderMessageID: receipt.ProviderReference, EvidenceSHA256: receipt.EvidenceSHA256, AcceptedAt: time.Now().UTC()}); e != nil {
			return e
		}
	}
	// Reconcile completion is only allowed without a live worker. The accepted
	// fence and exact original payload remain authority; attempts are unchanged.
	_, e = s.pool.Exec(ctx, `update platform.job j set completed_at=clock_timestamp(),terminal_error_code=null,claimed_by=null,claimed_until=null
 where j.tenant_id=$1 and j.job_id=$2 and j.job_type='social.publish' and j.payload->>'approval_id'=$3 and j.payload->>'profile_sha256'=$4
 and j.completed_at is null and (j.claimed_until is null or j.claimed_until<=clock_timestamp())
 and exists(select 1 from communication.outbound_delivery d where d.tenant_id=j.tenant_id and d.channel_code='facebook_pages' and d.delivery_key=$5 and d.state='accepted' and d.request_sha256_hex=$6)`, r.Intent.TenantID, socialJobID(r.Intent.TenantID, id), id, s.profile.SHA256(), r.Intent.DeliveryKey, hash)
	return e
}
