package whatsappbridge

// AUTHORED explicit, resumable composition of original per-message approvals.
import (
	"context"
	"elite.local/enterprise/internal/approval"
	cr "elite.local/enterprise/internal/catalogrelease"
	"elite.local/enterprise/internal/platform/identity"
	"encoding/json"
	"errors"
)

type CampaignItem struct {
	LeadID        string `json:"lead_id"`
	Step          int    `json:"step"`
	DeliveryKey   string `json:"delivery_key"`
	RequestSHA256 string `json:"request_sha256,omitempty"`
	State         string `json:"state"`
	Code          string `json:"code,omitempty"`
}
type CampaignBatch struct {
	CampaignID     string         `json:"campaign_id"`
	CampaignSHA256 string         `json:"campaign_sha256"`
	Complete       bool           `json:"complete"`
	Items          []CampaignItem `json:"items"`
}

func (c *Campaigns) Create(ctx context.Context, p identity.Principal, in CampaignSnapshot) (CampaignBatch, error) {
	var out CampaignBatch
	if !c.authorized(p, "marketing:request") || !c.authorized(p, "notification:request") || in.Requester != p.Subject {
		return out, ErrApproval
	}
	raw, h, e := campaignHash(in)
	if e != nil || len(raw) > 32768 {
		return out, ErrApproval
	}
	s := c.schedule.approvals
	tx, e := s.base.pool.Begin(ctx)
	if e != nil {
		return out, e
	}
	defer tx.Rollback(ctx)
	var exists bool
	e = tx.QueryRow(ctx, "select exists(select 1 from communication.whatsapp_campaign where tenant_id=$1 and campaign_id=$2)", s.tenant, in.Request.ID).Scan(&exists)
	if e != nil {
		return out, e
	}
	if !exists {
		fresh, e := c.prepare(ctx, tx, p, in.Request)
		if e != nil {
			return out, e
		}
		_, actual, e := campaignHash(fresh)
		if e != nil || actual != h {
			return out, ErrApproval
		}
		_, e = tx.Exec(ctx, `insert into communication.whatsapp_campaign(tenant_id,campaign_id,organization_id,creator_subject,campaign_sha256,payload)values($1,$2,$3,$4,$5,$6)on conflict(tenant_id,campaign_id)do nothing`, s.tenant, in.Request.ID, s.org, p.Subject, h, json.RawMessage(raw))
		if e != nil {
			return out, e
		}
		stored, actual, e := c.load(ctx, tx, in.Request.ID)
		if e != nil {
			return out, e
		}
		if actual != h || stored.Requester != p.Subject {
			return out, ErrApproval
		}
		for _, member := range fresh.Members {
			_, mh, e := campaignHash(member)
			if e != nil {
				return out, e
			}
			_, e = tx.Exec(ctx, `insert into communication.whatsapp_campaign_member(tenant_id,campaign_id,lead_id,member_sha256)values($1,$2,$3,$4)on conflict do nothing`, s.tenant, in.Request.ID, member.LeadID, mh)
			if e != nil {
				return out, e
			}
		}
	} else {
		stored, actual, e := c.load(ctx, tx, in.Request.ID)
		if e != nil {
			return out, e
		}
		if actual != h || stored.Requester != p.Subject {
			return out, ErrApproval
		}
	}
	if e = tx.Commit(ctx); e != nil {
		return out, e
	}
	return c.Resume(ctx, p, in.Request.ID, h)
}
func (c *Campaigns) Resume(ctx context.Context, p identity.Principal, id, hash string) (CampaignBatch, error) {
	var out CampaignBatch
	if !c.authorized(p, "marketing:request") || !c.authorized(p, "notification:request") || !validDigest(hash) {
		return out, ErrApproval
	}
	s := c.schedule.approvals
	plan, h, e := c.load(ctx, s.base.pool, id)
	if e != nil {
		return out, e
	}
	if hash != h || plan.Requester != p.Subject {
		return out, ErrApproval
	}
	out = CampaignBatch{CampaignID: id, CampaignSHA256: h, Complete: true}
	for _, m := range plan.Members {
		for step := 1; step <= len(plan.Request.Steps); step++ {
			request := campaignScheduleRequest(plan, m, step)
			item := CampaignItem{LeadID: m.LeadID, Step: step, DeliveryKey: scheduleKey(s.tenant, request.RequestID)}
			var exists bool
			if e = s.base.pool.QueryRow(ctx, "select exists(select 1 from approval.request where tenant_id=$1 and request_id=$2)", s.tenant, item.DeliveryKey).Scan(&exists); e != nil {
				return out, e
			}
			if exists {
				bound, bh, err := s.request(ctx, item.DeliveryKey)
				if err != nil || bound.SourceSHA256 != h || bound.Request.CampaignID != id || bound.LeadID != m.LeadID || bound.Request.Step != step {
					item.Code = "EXISTING_REQUEST_UNVERIFIED"
					out.Complete = false
				} else {
					item.RequestSHA256 = bh
					item.State = "PREPARED"
				}
			} else {
				bound, err := s.Prepare(ctx, p, request)
				if err != nil {
					item.Code = "SOURCE_NOT_CURRENT"
					out.Complete = false
				} else {
					bh, _, err := s.Submit(ctx, p, bound)
					if err != nil {
						item.Code = "PREPARATION_UNCONFIRMED"
						out.Complete = false
					} else {
						item.RequestSHA256 = bh
						item.State = "PREPARED"
					}
				}
			}
			out.Items = append(out.Items, item)
		}
	}
	return out, nil
}
func (c *Campaigns) Review(ctx context.Context, p identity.Principal, id, hash string, yes bool, reason string) (CampaignBatch, error) {
	var out CampaignBatch
	if !c.authorized(p, "marketing:approve") || !c.authorized(p, "notification:approve") || !validDigest(hash) || !cr.ValidText(reason, 2048) {
		return out, ErrApproval
	}
	s := c.schedule.approvals
	plan, h, e := c.load(ctx, s.base.pool, id)
	if e != nil {
		return out, e
	}
	if h != hash || plan.Requester == p.Subject {
		return out, ErrApproval
	}
	out = CampaignBatch{CampaignID: id, CampaignSHA256: h, Complete: true}
	for _, m := range plan.Members {
		for step := 1; step <= len(plan.Request.Steps); step++ {
			request := campaignScheduleRequest(plan, m, step)
			item := CampaignItem{LeadID: m.LeadID, Step: step, DeliveryKey: scheduleKey(s.tenant, request.RequestID)}
			bound, bh, err := s.request(ctx, item.DeliveryKey)
			if err != nil || bound.SourceSHA256 != h || bound.Request.CampaignID != id || bound.LeadID != m.LeadID || bound.Request.Step != step {
				item.Code = "REQUEST_UNAVAILABLE"
				out.Complete = false
			} else {
				state, err := s.Decide(ctx, p, item.DeliveryKey, bh, yes, reason)
				if errors.Is(err, approval.ErrNotPending) {
					state, err = c.decisionReplay(ctx, p, item.DeliveryKey, bh, yes, reason)
				}
				item.RequestSHA256 = bh
				item.State = string(state)
				if err != nil {
					item.Code = "DECISION_UNCONFIRMED"
					out.Complete = false
				}
			}
			out.Items = append(out.Items, item)
		}
	}
	return out, nil
}
func (c *Campaigns) Stop(ctx context.Context, p identity.Principal, id, hash, reason string) error {
	if !c.authorized(p, "marketing:cancel") || !validDigest(hash) || !cr.ValidText(reason, 2048) {
		return ErrApproval
	}
	s := c.schedule.approvals
	tx, e := s.base.pool.Begin(ctx)
	if e != nil {
		return e
	}
	defer tx.Rollback(ctx)
	var actual string
	if e = tx.QueryRow(ctx, `select campaign_sha256 from communication.whatsapp_campaign where tenant_id=$1 and campaign_id=$2 and organization_id=$3 for update`, s.tenant, id, s.org).Scan(&actual); e != nil {
		return e
	}
	if actual != hash {
		return ErrApproval
	}
	_, e = tx.Exec(ctx, `insert into communication.whatsapp_campaign_stop(tenant_id,campaign_id,campaign_sha256,actor_subject,reason)values($1,$2,$3,$4,$5)on conflict do nothing`, s.tenant, id, hash, p.Subject, reason)
	if e != nil {
		return e
	}
	var actor, storedReason string
	e = tx.QueryRow(ctx, "select actor_subject,reason from communication.whatsapp_campaign_stop where tenant_id=$1 and campaign_id=$2", s.tenant, id).Scan(&actor, &storedReason)
	if e != nil {
		return e
	}
	if actor != p.Subject || storedReason != reason {
		return ErrApproval
	}
	return tx.Commit(ctx)
}

// A replay reads the original immutable decision; it creates no approval/job.
func (c *Campaigns) decisionReplay(ctx context.Context, p identity.Principal, key, hash string, yes bool, reason string) (approval.State, error) {
	s := c.schedule.approvals
	want := approval.StateRejected
	if yes {
		want = approval.StateApproved
	}
	var exact bool
	e := s.base.pool.QueryRow(ctx, `select exists(
      select 1 from approval.request r join approval.decision d on d.tenant_id=r.tenant_id and d.request_id=r.request_id
      where r.tenant_id=$1 and r.organization_id=$2 and r.request_id=$3 and r.evidence_sha=$4 and r.kind='whatsapp_schedule' and r.state=$5
      and d.reviewer=$6 and d.reviewer<>r.requester and d.approved=$7 and d.reason=$8
      and(select count(*)from approval.decision d2 where d2.tenant_id=r.tenant_id and d2.request_id=r.request_id)=1)`, s.tenant, s.org, key, hash, want, p.Subject, yes, reason).Scan(&exact)
	if e != nil {
		return "", e
	}
	if !exact {
		return "", ErrApproval
	}
	return want, nil
}
