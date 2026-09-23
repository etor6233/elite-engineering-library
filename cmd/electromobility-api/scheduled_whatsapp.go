package main

// AUTHORED optional same-host composition; all credential and identity paths
// remain the original WhatsApp host's verified runtime sources.
import (
	"context"
	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/whatsappbridge"
	"github.com/jackc/pgx/v5/pgxpool"
)

func init() { whatsappScheduleFactory = selectedWhatsAppSchedule }
func selectedWhatsAppSchedule(ctx context.Context, h *whatsappHost, c whatsappHostConfig, pool *pgxpool.Pool, base *whatsappbridge.PostgresAppointmentApprovals, sender *whatsappbridge.Sender, fence *postgres.OutboundDeliveryStore) error {
	if h == nil || pool == nil {
		return errWhatsAppHost
	}
	p, e := h.identity.Resolve(ctx)
	if e != nil || !p.Allowed("notification:dispatch") {
		return errWhatsAppHost
	}
	approvals, e := whatsappbridge.NewScheduleApprovals(base, c.TenantID, c.OrganizationID, c.ConnectionID, sender.Profile)
	if e != nil {
		return e
	}
	if e = approvals.VerifyInfrastructure(ctx); e != nil {
		return e
	}
	module, e := whatsappbridge.NewScheduledNotifications(approvals, sender, fence, whatsappbridge.VerifiedInboxChannel{}, c.WorkerID+"-scheduled")
	if e != nil {
		return e
	}
	if c.CampaignPolicyFile != "" || c.CampaignPolicySHA256 != "" {
		if whatsappCampaignFactory == nil {
			return errWhatsAppHost
		}
		campaign, e := whatsappCampaignFactory(ctx, module, c)
		if e != nil {
			return e
		}
		h.modules = append(h.modules, campaign)
	}
	h.modules = append(h.modules, module)
	h.extraRun = func(ctx context.Context) error { return module.Run(ctx, h.identity, h.reporter, h.interval) }
	return nil
}

var whatsappCampaignFactory func(context.Context, *whatsappbridge.ScheduledNotifications, whatsappHostConfig) (httpapi.EnterpriseModule, error)
