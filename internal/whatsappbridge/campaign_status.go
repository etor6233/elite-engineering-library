package whatsappbridge

// AUTHORED read-only observation of original message and quotation/order evidence.
import (
	"context"
	"elite.local/enterprise/internal/platform/identity"
	"errors"
	"github.com/jackc/pgx/v5"
	"time"
)

type CampaignConversion struct {
	LeadID         string    `json:"lead_id"`
	QuotationID    string    `json:"quotation_id"`
	OrderID        string    `json:"order_id"`
	OrderState     string    `json:"order_state"`
	AcceptedAt     time.Time `json:"accepted_at"`
	EvidenceSHA256 string    `json:"evidence_sha256"`
}
type CampaignItemStatus struct {
	LeadID      string          `json:"lead_id"`
	Step        int             `json:"step"`
	DeliveryKey string          `json:"delivery_key"`
	Preparation string          `json:"preparation"`
	Schedule    *ScheduleStatus `json:"schedule,omitempty"`
}
type CampaignStatus struct {
	CampaignSHA256    string               `json:"campaign_sha256"`
	Snapshot          CampaignSnapshot     `json:"snapshot"`
	Stopped           bool                 `json:"stopped"`
	Items             []CampaignItemStatus `json:"items"`
	Conversions       []CampaignConversion `json:"conversions"`
	ConversionMeaning string               `json:"conversion_meaning"`
}

func (c *Campaigns) Status(ctx context.Context, p identity.Principal, id string) (CampaignStatus, error) {
	var out CampaignStatus
	if !c.authorized(p, "marketing:read") || !c.authorized(p, "notification:read") {
		return out, ErrApproval
	}
	s := c.schedule.approvals
	plan, h, e := c.load(ctx, s.base.pool, id)
	if e != nil {
		return out, e
	}
	out = CampaignStatus{CampaignSHA256: h, Snapshot: plan, ConversionMeaning: "First original quote/order acceptance after an accepted campaign message; observed sequence, not causal attribution or payment proof", Conversions: []CampaignConversion{}}
	if e = s.base.pool.QueryRow(ctx, "select exists(select 1 from communication.whatsapp_campaign_stop where tenant_id=$1 and campaign_id=$2)", s.tenant, id).Scan(&out.Stopped); e != nil {
		return out, e
	}
	for _, m := range plan.Members {
		for step := 1; step <= len(plan.Request.Steps); step++ {
			request := campaignScheduleRequest(plan, m, step)
			key := scheduleKey(s.tenant, request.RequestID)
			item := CampaignItemStatus{LeadID: m.LeadID, Step: step, DeliveryKey: key, Preparation: "NOT_PREPARED"}
			var exists bool
			e = s.base.pool.QueryRow(ctx, "select exists(select 1 from approval.request where tenant_id=$1 and request_id=$2)", s.tenant, key).Scan(&exists)
			if e != nil {
				return out, e
			}
			if exists {
				status, e := s.Status(ctx, p, key)
				if e != nil {
					return out, e
				}
				if status.Context.SourceSHA256 != h {
					return out, ErrApproval
				}
				item.Preparation = "PREPARED"
				item.Schedule = &status
			}
			out.Items = append(out.Items, item)
		}
		conversion := CampaignConversion{LeadID: m.LeadID}
		e = s.base.pool.QueryRow(ctx, `select a.quotation_id,a.order_id,o.state,a.accepted_at,a.evidence_sha256_hex
 from sales.quotation q join sales.quotation_acceptance a on a.tenant_id=q.tenant_id and a.quotation_id=q.quotation_id and a.order_id=q.order_id
 join sales.customer_order o on o.tenant_id=a.tenant_id and o.order_id=a.order_id and o.organization_id=q.organization_id and o.customer_principal_id=a.customer_principal_id
 where q.tenant_id=$1 and q.organization_id=$2 and q.lead_id=$3 and a.accepted_at<=clock_timestamp()
 and exists(select 1 from communication.whatsapp_schedule s join communication.outbound_delivery d on d.tenant_id=s.tenant_id and d.delivery_key=s.delivery_key and d.channel_code='whatsapp'
 where s.tenant_id=q.tenant_id and s.campaign_id=$4 and s.lead_id=q.lead_id and d.state='accepted' and d.accepted_at<=a.accepted_at)
 order by a.accepted_at,a.quotation_id limit 1`, s.tenant, s.org, m.LeadID, id).Scan(&conversion.QuotationID, &conversion.OrderID, &conversion.OrderState, &conversion.AcceptedAt, &conversion.EvidenceSHA256)
		if e != nil && !errors.Is(e, pgx.ErrNoRows) {
			return out, e
		}
		if e == nil {
			out.Conversions = append(out.Conversions, conversion)
		}
	}
	return out, nil
}
func (c *Campaigns) VerifyInfrastructure(ctx context.Context) error {
	var valid bool
	e := c.schedule.approvals.base.pool.QueryRow(ctx, `select
 (select count(*)from pg_trigger where not tgisinternal and tgenabled in('O','A') and
 (tgrelid=to_regclass('communication.whatsapp_campaign')and tgname='whatsapp_campaign_immutable'
 or tgrelid=to_regclass('communication.whatsapp_campaign_member')and tgname='whatsapp_campaign_member_immutable'
 or tgrelid=to_regclass('communication.whatsapp_campaign_stop')and tgname='whatsapp_campaign_stop_immutable'))=3
 and(select count(*)from pg_constraint where conrelid=to_regclass('communication.whatsapp_schedule') and convalidated
 and conname in('whatsapp_schedule_source_check','whatsapp_schedule_campaign_member_fk','whatsapp_schedule_campaign_step_unique'))=3`).Scan(&valid)
	if e != nil || !valid {
		return ErrApproval
	}
	return nil
}
