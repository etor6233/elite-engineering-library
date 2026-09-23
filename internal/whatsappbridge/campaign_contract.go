package whatsappbridge

// AUTHORED bounded audience/source binding around existing CRM and consent rows.
import (
	"context"
	"elite.local/enterprise/internal/approval"
	cr "elite.local/enterprise/internal/catalogrelease"
	"elite.local/enterprise/internal/contactidentity"
	"elite.local/enterprise/internal/platform/identity"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5"
	"sort"
	"time"
)

type CampaignPolicy struct {
	Schema         string `json:"schema"`
	ConsentPurpose string `json:"consent_purpose"`
	PolicyVersion  string `json:"policy_version"`
	MaxMembers     int    `json:"max_members"`
	MaxSteps       int    `json:"max_steps"`
}
type CampaignRecipient struct {
	LeadID    string `json:"lead_id"`
	Recipient string `json:"recipient"`
}
type CampaignStep struct {
	TemplateName   string    `json:"template_name"`
	LanguageCode   string    `json:"language_code"`
	BodyParameters []string  `json:"body_parameters"`
	NotBefore      time.Time `json:"not_before"`
	ExpiresAt      time.Time `json:"expires_at"`
}
type CampaignRequest struct {
	ID      string              `json:"campaign_id"`
	Sources []string            `json:"sources"`
	States  []string            `json:"states"`
	Members []CampaignRecipient `json:"members"`
	Steps   []CampaignStep      `json:"steps"`
}
type CampaignMember struct {
	LeadID         string    `json:"lead_id"`
	Recipient      string    `json:"recipient"`
	Source         string    `json:"source"`
	Lifecycle      string    `json:"lifecycle"`
	LeadUpdatedAt  time.Time `json:"lead_updated_at"`
	SubjectID      string    `json:"subject_id"`
	ExternalIDHMAC string    `json:"external_id_hmac"`
	BindingVersion int64     `json:"binding_version,string"`
	BindingPolicy  string    `json:"binding_policy"`
	ConsentID      string    `json:"consent_id"`
	ConsentSHA256  string    `json:"consent_sha256"`
}
type CampaignSnapshot struct {
	Schema         string           `json:"schema"`
	Request        CampaignRequest  `json:"request"`
	Requester      string           `json:"requester"`
	TenantID       string           `json:"tenant_id"`
	OrganizationID string           `json:"organization_id"`
	ConnectionID   string           `json:"connection_id"`
	ProfileSHA256  string           `json:"profile_sha256"`
	Policy         CampaignPolicy   `json:"policy"`
	Members        []CampaignMember `json:"members"`
}
type Campaigns struct {
	schedule *ScheduledNotifications
	policy   CampaignPolicy
}

func NewCampaigns(schedule *ScheduledNotifications, policy CampaignPolicy) (*Campaigns, error) {
	if schedule == nil || schedule.approvals == nil || schedule.approvals.extension != nil ||
		policy.Schema != "elite-whatsapp-campaign-policy/v1" || !cr.ValidID(policy.ConsentPurpose) || !cr.ValidID(policy.PolicyVersion) ||
		policy.ConsentPurpose == schedule.approvals.base.purpose || policy.MaxMembers < 1 || policy.MaxMembers > 20 || policy.MaxSteps < 1 || policy.MaxSteps > 3 {
		return nil, ErrApproval
	}
	c := &Campaigns{schedule: schedule, policy: policy}
	schedule.approvals.extension = c
	return c, nil
}
func (c *Campaigns) authorized(p identity.Principal, permission string) bool {
	return c != nil && c.schedule.approvals.authorized(p, permission)
}
func campaignHash(v any) ([]byte, string, error) {
	raw, e := json.Marshal(v)
	if e != nil {
		return nil, "", e
	}
	return approval.CanonicalPayload(raw)
}
func (c *Campaigns) validate(r CampaignRequest) error {
	if !cr.ValidID(r.ID) || len(r.ID) < 16 || len(r.Members) < 1 || len(r.Members) > c.policy.MaxMembers || len(r.Steps) < 1 || len(r.Steps) > c.policy.MaxSteps || len(r.Sources) < 1 || len(r.Sources) > 8 || len(r.States) < 1 || len(r.States) > 3 {
		return ErrApproval
	}
	unique := func(v []string) bool {
		seen := map[string]bool{}
		for _, x := range v {
			if !cr.ValidID(x) || seen[x] {
				return false
			}
			seen[x] = true
		}
		return true
	}
	if !unique(r.Sources) || !unique(r.States) {
		return ErrApproval
	}
	for _, state := range r.States {
		if state != "new" && state != "qualified" && state != "contacted" {
			return ErrApproval
		}
	}
	seen := map[string]bool{}
	for _, m := range r.Members {
		if !cr.ValidID(m.LeadID) || seen[m.LeadID] {
			return ErrApproval
		}
		seen[m.LeadID] = true
	}
	var previous time.Time
	var profile struct {
		Templates []struct {
			Name     string `json:"name"`
			Language string `json:"language_code"`
			Count    int    `json:"body_parameter_count"`
		} `json:"approved_templates"`
	}
	if json.Unmarshal(c.schedule.approvals.profile, &profile) != nil {
		return ErrApproval
	}
	for _, s := range r.Steps {
		if s.NotBefore.IsZero() || !s.ExpiresAt.After(s.NotBefore) || s.ExpiresAt.Sub(s.NotBefore) > 24*time.Hour || !previous.IsZero() && !s.NotBefore.After(previous) || s.BodyParameters == nil || len(s.BodyParameters) > 20 {
			return ErrApproval
		}
		matches := 0
		for _, t := range profile.Templates {
			if t.Name == s.TemplateName && t.Language == s.LanguageCode && t.Count == len(s.BodyParameters) {
				matches++
			}
		}
		if matches != 1 {
			return ErrApproval
		}
		for _, v := range s.BodyParameters {
			if !cr.ValidText(v, 1024) {
				return ErrApproval
			}
		}
		previous = s.ExpiresAt
	}
	return nil
}
func (c *Campaigns) member(ctx context.Context, tx pgx.Tx, recipient CampaignRecipient) (CampaignMember, error) {
	s := c.schedule.approvals
	m := CampaignMember{LeadID: recipient.LeadID, Recipient: recipient.Recipient}
	h, e := contactidentity.ExternalDigest(s.base.hmacKey, s.tenant, "whatsapp", recipient.Recipient)
	if e != nil {
		return m, ErrApproval
	}
	e = tx.QueryRow(ctx, `select l.source_code,l.lifecycle_state,l.updated_at,b.subject_id,b.version,b.policy_version,consent.consent_id,consent.evidence_sha256_hex
 from crm.lead l join org.organization o on o.tenant_id=l.tenant_id and o.organization_id=l.organization_id
 join platform.tenant tenant on tenant.tenant_id=l.tenant_id
 join communication.contact_channel_binding b on b.tenant_id=l.tenant_id and b.lead_id=l.lead_id
 join crm.consent_evidence consent on consent.tenant_id=l.tenant_id and consent.lead_id=l.lead_id
 join integration.provider_connection pc on pc.tenant_id=l.tenant_id and pc.organization_id=l.organization_id
 where l.tenant_id=$1 and l.organization_id=$2 and l.lead_id=$3 and o.status='active' and tenant.status='active'
 and b.channel_code='whatsapp' and b.external_id_hmac=$4 and b.state='active' and b.pii_allowed
 and b.effective_at<=statement_timestamp()
 and b.policy_version=$5 and consent.purpose_code=$6 and consent.policy_version=$5 and consent.decision='granted' and consent.occurred_at<=statement_timestamp()
 and not exists(select 1 from crm.consent_evidence newer where newer.tenant_id=consent.tenant_id and newer.lead_id=consent.lead_id and newer.purpose_code=consent.purpose_code and newer.consent_id<>consent.consent_id and newer.occurred_at>=consent.occurred_at and newer.occurred_at<=statement_timestamp())
 and pc.connection_id=$7 and pc.provider_code='meta-whatsapp' and pc.state='active'
 for share of l,o,tenant,b,consent,pc`, s.tenant, s.org, m.LeadID, h, c.policy.PolicyVersion, c.policy.ConsentPurpose, s.connection).Scan(&m.Source, &m.Lifecycle, &m.LeadUpdatedAt, &m.SubjectID, &m.BindingVersion, &m.BindingPolicy, &m.ConsentID, &m.ConsentSHA256)
	if errors.Is(e, pgx.ErrNoRows) {
		return m, ErrApproval
	}
	m.ExternalIDHMAC = h
	return m, e
}
func contains(v []string, s string) bool {
	for _, x := range v {
		if x == s {
			return true
		}
	}
	return false
}
func (c *Campaigns) prepare(ctx context.Context, tx pgx.Tx, p identity.Principal, r CampaignRequest) (CampaignSnapshot, error) {
	var out CampaignSnapshot
	if c.validate(r) != nil {
		return out, ErrApproval
	}
	// Copy before sorting: a caller's request is never mutated.
	r.Members = append([]CampaignRecipient(nil), r.Members...)
	sort.Slice(r.Members, func(i, j int) bool { return r.Members[i].LeadID < r.Members[j].LeadID })
	r.Sources = append([]string(nil), r.Sources...)
	sort.Strings(r.Sources)
	r.States = append([]string(nil), r.States...)
	sort.Strings(r.States)
	s := c.schedule.approvals
	out = CampaignSnapshot{Schema: "elite-whatsapp-campaign/v1", Request: r, Requester: p.Subject, TenantID: s.tenant, OrganizationID: s.org, ConnectionID: s.connection, ProfileSHA256: digest(s.profile), Policy: c.policy}
	for _, selected := range r.Members {
		m, e := c.member(ctx, tx, selected)
		if e != nil {
			return out, e
		}
		if !contains(r.Sources, m.Source) || !contains(r.States, m.Lifecycle) {
			return out, ErrApproval
		}
		out.Members = append(out.Members, m)
	}
	var now time.Time
	if e := tx.QueryRow(ctx, "select clock_timestamp()").Scan(&now); e != nil {
		return out, e
	}
	for _, step := range r.Steps {
		if !step.NotBefore.After(now) || step.NotBefore.After(now.Add(30*24*time.Hour)) {
			return out, ErrApproval
		}
	}
	if raw, _, e := campaignHash(out); e != nil || len(raw) > 32768 {
		return out, ErrApproval
	}
	return out, nil
}
func (c *Campaigns) Prepare(ctx context.Context, p identity.Principal, r CampaignRequest) (CampaignSnapshot, error) {
	if !c.authorized(p, "marketing:request") {
		return CampaignSnapshot{}, ErrApproval
	}
	tx, e := c.schedule.approvals.base.pool.Begin(ctx)
	if e != nil {
		return CampaignSnapshot{}, e
	}
	defer tx.Rollback(ctx)
	out, e := c.prepare(ctx, tx, p, r)
	if e != nil {
		return out, e
	}
	return out, tx.Commit(ctx)
}
