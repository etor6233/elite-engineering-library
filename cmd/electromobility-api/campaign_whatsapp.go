package main

// AUTHORED optional fixed-policy configuration; no new token or provider path.
import (
	"bytes"
	"context"
	"elite.local/enterprise/internal/approval"
	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/whatsappbridge"
	"encoding/json"
	"io"
)

func init() { whatsappCampaignFactory = selectedWhatsAppCampaign }
func selectedWhatsAppCampaign(ctx context.Context, schedule *whatsappbridge.ScheduledNotifications, c whatsappHostConfig) (httpapi.EnterpriseModule, error) {
	raw, e := readWhatsAppExact(c.CampaignPolicyFile, c.CampaignPolicySHA256, 4096)
	if e != nil {
		return nil, e
	}
	if _, _, e = approval.CanonicalPayload(raw); e != nil {
		return nil, errWhatsAppHost
	}
	var policy whatsappbridge.CampaignPolicy
	d := json.NewDecoder(bytes.NewReader(raw))
	d.DisallowUnknownFields()
	if d.Decode(&policy) != nil || d.Decode(new(any)) != io.EOF {
		return nil, errWhatsAppHost
	}
	module, e := whatsappbridge.NewCampaigns(schedule, policy)
	if e != nil {
		return nil, e
	}
	if e = module.VerifyInfrastructure(ctx); e != nil {
		return nil, e
	}
	return module, nil
}
