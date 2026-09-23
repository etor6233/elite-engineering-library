package whatsappbridge

// AUTHORED campaign source implements the existing schedule extension.
import (
	"context"
	"elite.local/enterprise/internal/channels"
	"elite.local/enterprise/internal/outbounddelivery"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"time"
)

type campaignQuery interface {
	QueryRow(context.Context, string, ...any) pgx.Row
}

func (c *Campaigns) load(ctx context.Context, db campaignQuery, id string) (CampaignSnapshot, string, error) {
	var out CampaignSnapshot
	var raw []byte
	var h string
	s := c.schedule.approvals
	e := db.QueryRow(ctx, `select payload,campaign_sha256 from communication.whatsapp_campaign where tenant_id=$1 and organization_id=$2 and campaign_id=$3 for share`, s.tenant, s.org, id).Scan(&raw, &h)
	if errors.Is(e, pgx.ErrNoRows) {
		return out, "", ErrApproval
	}
	if e != nil {
		return out, "", e
	}
	if json.Unmarshal(raw, &out) != nil {
		return out, "", ErrApproval
	}
	_, actual, e := campaignHash(out)
	if e != nil || actual != h || out.Schema != "elite-whatsapp-campaign/v1" || out.TenantID != s.tenant || out.OrganizationID != s.org || out.ConnectionID != s.connection || out.ProfileSHA256 != digest(s.profile) || out.Policy != c.policy || out.Request.ID != id || len(out.Members) != len(out.Request.Members) || c.validate(out.Request) != nil {
		return out, "", ErrApproval
	}
	return out, h, nil
}
func campaignScheduleRequest(plan CampaignSnapshot, member CampaignMember, step int) ScheduleRequest {
	r := plan.Request.Steps[step-1]
	return ScheduleRequest{RequestID: "campaign-" + digest([]byte(plan.Request.ID+"\x00"+member.LeadID+"\x00"+fmt.Sprint(step))), CampaignID: plan.Request.ID, LeadID: member.LeadID, Step: step, Recipient: member.Recipient, TemplateName: r.TemplateName, LanguageCode: r.LanguageCode, BodyParameters: append([]string{}, r.BodyParameters...), NotBefore: r.NotBefore, ExpiresAt: r.ExpiresAt}
}
func (c *Campaigns) stoppedOrConverted(ctx context.Context, db campaignQuery, id, lead string) (bool, error) {
	var stopped bool
	s := c.schedule.approvals
	e := db.QueryRow(ctx, `select exists(select 1 from communication.whatsapp_campaign_stop where tenant_id=$1 and campaign_id=$2)
 or exists(select 1 from sales.quotation q join sales.quotation_acceptance a on a.tenant_id=q.tenant_id and a.quotation_id=q.quotation_id
 join sales.customer_order o on o.tenant_id=a.tenant_id and o.order_id=a.order_id
 where q.tenant_id=$1 and q.organization_id=$3 and q.lead_id=$4 and q.order_id=a.order_id and a.accepted_at<=clock_timestamp())`, s.tenant, id, s.org, lead).Scan(&stopped)
	return stopped, e
}
func (c *Campaigns) Source(ctx context.Context, tx pgx.Tx, r ScheduleRequest, preparing bool) (ScheduleContext, error) {
	var out ScheduleContext
	plan, h, e := c.load(ctx, tx, r.CampaignID)
	if e != nil {
		return out, e
	}
	if r.Step < 1 || r.Step > len(plan.Request.Steps) || r.AppointmentID != "" {
		return out, ErrApproval
	}
	var selected *CampaignMember
	for i := range plan.Members {
		if plan.Members[i].LeadID == r.LeadID {
			selected = &plan.Members[i]
			break
		}
	}
	if selected == nil {
		return out, ErrApproval
	}
	_, requestHash, e := campaignHash(r)
	_, expected, e2 := campaignHash(campaignScheduleRequest(plan, *selected, r.Step))
	if e != nil || e2 != nil || requestHash != expected {
		return out, ErrApproval
	}
	fresh, e := c.member(ctx, tx, CampaignRecipient{selected.LeadID, selected.Recipient})
	if e != nil {
		return out, e
	}
	_, memberHash, e := campaignHash(fresh)
	_, before, e2 := campaignHash(*selected)
	if e != nil || e2 != nil || memberHash != before || !contains(plan.Request.Sources, fresh.Source) || !contains(plan.Request.States, fresh.Lifecycle) {
		return out, ErrApproval
	}
	if stopped, e := c.stoppedOrConverted(ctx, tx, r.CampaignID, r.LeadID); e != nil {
		return out, e
	} else if stopped {
		return out, ErrApproval
	}
	var now time.Time
	if e = tx.QueryRow(ctx, "select clock_timestamp()").Scan(&now); e != nil {
		return out, e
	}
	if !r.ExpiresAt.After(now) || preparing && !r.NotBefore.After(now) {
		return out, ErrApproval
	}
	s := c.schedule.approvals
	body, _ := json.Marshal(struct {
		Recipient      string   `json:"recipient"`
		TemplateName   string   `json:"template_name"`
		LanguageCode   string   `json:"language_code"`
		BodyParameters []string `json:"body_parameters"`
	}{r.Recipient, r.TemplateName, r.LanguageCode, r.BodyParameters})
	out = ScheduleContext{Schema: "elite-whatsapp-schedule/v1", Request: r, DeliveryKey: scheduleKey(s.tenant, r.RequestID), OrganizationID: s.org, ConnectionID: s.connection, LeadID: r.LeadID, SubjectID: fresh.SubjectID, ExternalIDHMAC: fresh.ExternalIDHMAC, BindingVersion: fresh.BindingVersion, PolicyVersion: fresh.BindingPolicy, ConsentID: fresh.ConsentID, ConsentPurpose: c.policy.ConsentPurpose, ConsentEvidenceSHA256: fresh.ConsentSHA256, ProfileSHA256: digest(s.profile), NotBefore: r.NotBefore, ExpiresAt: r.ExpiresAt, SourceSHA256: h, MemberSHA256: memberHash}
	out.Message = channels.Message{TenantID: s.tenant, ChannelCode: "whatsapp", Direction: channels.DirectionOut, DeliveryKey: out.DeliveryKey, ExternalID: r.Recipient, Text: string(body)}
	out.MessageSHA256, e = outbounddelivery.MessageSHA256(out.Message)
	return out, e
}
func (c *Campaigns) Enqueue(ctx context.Context, tx pgx.Tx, b ScheduleContext, h string) error {
	s := c.schedule.approvals
	_, e := tx.Exec(ctx, `insert into communication.whatsapp_schedule(tenant_id,delivery_key,job_id,organization_id,appointment_id,appointment_version,not_before,expires_at,request_sha256,campaign_id,lead_id,campaign_step)
 select m.tenant_id,$2,gen_random_uuid(),p.organization_id,null,0,$3,$4,$5,m.campaign_id,m.lead_id,$6
 from communication.whatsapp_campaign_member m join communication.whatsapp_campaign p using(tenant_id,campaign_id)
 where m.tenant_id=$1 and m.campaign_id=$7 and m.lead_id=$8 and m.member_sha256=$9 and p.campaign_sha256=$10 and p.organization_id=$11
 on conflict(tenant_id,delivery_key)do nothing`, s.tenant, b.DeliveryKey, b.NotBefore, b.ExpiresAt, h, b.Request.Step, b.Request.CampaignID, b.LeadID, b.MemberSHA256, b.SourceSHA256, s.org)
	if e != nil {
		return e
	}
	var same bool
	e = tx.QueryRow(ctx, `select exists(select 1 from communication.whatsapp_schedule where tenant_id=$1 and delivery_key=$2 and request_sha256=$3 and campaign_id=$4 and lead_id=$5 and campaign_step=$6)`, s.tenant, b.DeliveryKey, h, b.Request.CampaignID, b.LeadID, b.Request.Step).Scan(&same)
	if e != nil {
		return e
	}
	if !same {
		return ErrApproval
	}
	return nil
}
func (c *Campaigns) Admit(ctx context.Context, tx pgx.Tx, b ScheduleContext) error {
	if b.Request.Step <= 1 {
		return nil
	}
	plan, _, e := c.load(ctx, tx, b.Request.CampaignID)
	if e != nil {
		return e
	}
	for _, m := range plan.Members {
		if m.LeadID == b.LeadID {
			previous := campaignScheduleRequest(plan, m, b.Request.Step-1)
			var accepted bool
			e = tx.QueryRow(ctx, `select exists(select 1 from communication.outbound_delivery where tenant_id=$1 and delivery_key=$2 and channel_code='whatsapp' and state='accepted')`, b.Message.TenantID, scheduleKey(b.Message.TenantID, previous.RequestID)).Scan(&accepted)
			if e != nil {
				return e
			}
			if !accepted {
				return ErrApproval
			}
			return nil
		}
	}
	return ErrApproval
}
