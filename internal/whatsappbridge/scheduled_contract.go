package whatsappbridge

// AUTHORED source/time/request binding over the existing CRM and template owner.
import (
	"context"
	"elite.local/enterprise/internal/approval"
	cr "elite.local/enterprise/internal/catalogrelease"
	"elite.local/enterprise/internal/channels"
	"elite.local/enterprise/internal/contactidentity"
	"elite.local/enterprise/internal/outbounddelivery"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/postgres"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5"
	"strings"
	"time"
)

// Optional source resolver shares the original approval/job/fence machinery.
type scheduledSourceExtension interface {
	Source(context.Context, pgx.Tx, ScheduleRequest, bool) (ScheduleContext, error)
	Enqueue(context.Context, pgx.Tx, ScheduleContext, string) error
	Admit(context.Context, pgx.Tx, ScheduleContext) error
}
type ScheduleRequest struct {
	CampaignID     string    `json:"campaign_id,omitempty"`
	LeadID         string    `json:"lead_id,omitempty"`
	Step           int       `json:"step,omitempty"`
	RequestID      string    `json:"request_id"`
	AppointmentID  string    `json:"appointment_id"`
	Recipient      string    `json:"recipient"`
	TemplateName   string    `json:"template_name"`
	LanguageCode   string    `json:"language_code"`
	BodyParameters []string  `json:"body_parameters"`
	NotBefore      time.Time `json:"not_before"`
	ExpiresAt      time.Time `json:"expires_at"`
}
type ScheduleContext struct {
	SourceSHA256          string           `json:"source_sha256,omitempty"`
	MemberSHA256          string           `json:"member_sha256,omitempty"`
	Schema                string           `json:"schema"`
	Request               ScheduleRequest  `json:"request"`
	DeliveryKey           string           `json:"delivery_key"`
	OrganizationID        string           `json:"organization_id"`
	ConnectionID          string           `json:"connection_id"`
	AppointmentID         string           `json:"appointment_id"`
	AppointmentVersion    int64            `json:"appointment_version,string"`
	AppointmentStartsAt   time.Time        `json:"appointment_starts_at"`
	LeadID                string           `json:"lead_id"`
	SubjectID             string           `json:"subject_id"`
	ExternalIDHMAC        string           `json:"external_id_hmac"`
	BindingVersion        int64            `json:"binding_version,string"`
	PolicyVersion         string           `json:"policy_version"`
	ConsentID             string           `json:"consent_id"`
	ConsentPurpose        string           `json:"consent_purpose"`
	ConsentEvidenceSHA256 string           `json:"consent_evidence_sha256"`
	ProfileSHA256         string           `json:"profile_sha256"`
	MessageSHA256         string           `json:"message_sha256"`
	NotBefore             time.Time        `json:"not_before"`
	ExpiresAt             time.Time        `json:"expires_at"`
	Message               channels.Message `json:"message"`
}
type ScheduleApprovals struct {
	extension               scheduledSourceExtension
	base                    *PostgresAppointmentApprovals
	tenant, org, connection string
	profile                 json.RawMessage
	human                   *postgres.HumanApprovals
}

func NewScheduleApprovals(base *PostgresAppointmentApprovals, tenant, org, connection string, profile json.RawMessage) (*ScheduleApprovals, error) {
	if base == nil || base.pool == nil || !notificationEventID.MatchString(tenant) || !cr.ValidID(org) || !webhookConnection.MatchString(connection) || len(profile) > 32768 || !json.Valid(profile) {
		return nil, ErrApproval
	}
	return &ScheduleApprovals{base: base, tenant: tenant, org: org, connection: connection, profile: append(json.RawMessage(nil), profile...), human: postgres.NewHumanApprovals(base.pool)}, nil
}
func (s *ScheduleApprovals) authorized(p identity.Principal, permission string) bool {
	return s != nil && p.TenantID == s.tenant && p.Subject != "" && len(p.Subject) <= 128 && p.Allowed(permission) && p.AllowedOrganization(s.org)
}
func scheduleKey(tenant, id string) string { return "wa-schedule:" + digest([]byte(tenant+"\x00"+id)) }
func validScheduleKey(key string) bool {
	return strings.HasPrefix(key, "wa-schedule:") && validDigest(strings.TrimPrefix(key, "wa-schedule:"))
}
func (s *ScheduleApprovals) source(ctx context.Context, tx pgx.Tx, r ScheduleRequest, preparing bool) (ScheduleContext, error) {
	if r.CampaignID != "" {
		if s.extension == nil || r.AppointmentID != "" {
			return ScheduleContext{}, ErrApproval
		}
		return s.extension.Source(ctx, tx, r, preparing)
	}
	if r.LeadID != "" || r.Step != 0 {
		return ScheduleContext{}, ErrApproval
	}
	var c ScheduleContext
	var now time.Time
	if !cr.ValidID(r.RequestID) || len(r.RequestID) < 16 || !cr.ValidID(r.AppointmentID) || r.BodyParameters == nil || len(r.BodyParameters) > 20 || r.NotBefore.IsZero() || !r.ExpiresAt.After(r.NotBefore) || r.ExpiresAt.Sub(r.NotBefore) > 24*time.Hour {
		return c, ErrApproval
	}
	for _, v := range r.BodyParameters {
		if !cr.ValidText(v, 1024) {
			return c, ErrApproval
		}
	}
	var profile struct {
		Templates []struct {
			Name     string `json:"name"`
			Language string `json:"language_code"`
			Count    int    `json:"body_parameter_count"`
		} `json:"approved_templates"`
	}
	if json.Unmarshal(s.profile, &profile) != nil {
		return c, ErrApproval
	}
	matches := 0
	for _, t := range profile.Templates {
		if t.Name == r.TemplateName && t.Language == r.LanguageCode && t.Count == len(r.BodyParameters) {
			matches++
		}
	}
	if matches != 1 {
		return c, ErrApproval
	}
	h, e := contactidentity.ExternalDigest(s.base.hmacKey, s.tenant, "whatsapp", r.Recipient)
	if e != nil {
		return c, ErrApproval
	}
	e = tx.QueryRow(ctx, `select clock_timestamp(),a.version,a.starts_at,l.lead_id,b.subject_id,b.version,b.policy_version,consent.consent_id,consent.evidence_sha256_hex
 from crm.appointment a join crm.lead l on l.tenant_id=a.tenant_id and l.lead_id=a.lead_id and l.organization_id=a.organization_id
 join org.organization o on o.tenant_id=a.tenant_id and o.organization_id=a.organization_id and o.status='active'
 join platform.tenant tenant on tenant.tenant_id=a.tenant_id and tenant.status='active'
 join integration.provider_connection pc on pc.tenant_id=a.tenant_id and pc.connection_id=$4 and pc.organization_id=a.organization_id and pc.provider_code='meta-whatsapp' and pc.state='active'
 join communication.contact_channel_binding b on b.tenant_id=a.tenant_id and b.lead_id=a.lead_id and b.channel_code='whatsapp' and b.external_id_hmac=$5 and b.state='active' and b.pii_allowed and b.effective_at<=statement_timestamp()
 join crm.consent_evidence consent on consent.tenant_id=l.tenant_id and consent.lead_id=l.lead_id and consent.purpose_code=$6 and consent.policy_version=$7 and consent.decision='granted' and consent.occurred_at<=statement_timestamp()
 where a.tenant_id=$1 and a.organization_id=$2 and a.appointment_id=$3 and a.state='confirmed' and a.starts_at>statement_timestamp() and b.policy_version=$7
 and not exists(select 1 from crm.consent_evidence newer where newer.tenant_id=consent.tenant_id and newer.lead_id=consent.lead_id and newer.purpose_code=consent.purpose_code and newer.consent_id<>consent.consent_id and newer.occurred_at>=consent.occurred_at and newer.occurred_at<=statement_timestamp())
 for share of a,l,b,consent,o,tenant,pc`, s.tenant, s.org, r.AppointmentID, s.connection, h, s.base.purpose, s.base.policy).Scan(&now, &c.AppointmentVersion, &c.AppointmentStartsAt, &c.LeadID, &c.SubjectID, &c.BindingVersion, &c.PolicyVersion, &c.ConsentID, &c.ConsentEvidenceSHA256)
	if errors.Is(e, pgx.ErrNoRows) {
		return c, ErrApproval
	}
	if e != nil {
		return c, e
	}
	if e = tx.QueryRow(ctx, `select clock_timestamp()`).Scan(&now); e != nil {
		return c, e
	}
	if !r.ExpiresAt.After(now) || r.ExpiresAt.After(c.AppointmentStartsAt) || preparing && !r.NotBefore.After(now) || r.NotBefore.After(now.Add(30*24*time.Hour)) {
		return c, ErrApproval
	}
	body, _ := json.Marshal(struct {
		Recipient      string   `json:"recipient"`
		TemplateName   string   `json:"template_name"`
		LanguageCode   string   `json:"language_code"`
		BodyParameters []string `json:"body_parameters"`
	}{r.Recipient, r.TemplateName, r.LanguageCode, r.BodyParameters})
	c.Schema = "elite-whatsapp-schedule/v1"
	c.Request = r
	c.DeliveryKey = scheduleKey(s.tenant, r.RequestID)
	c.OrganizationID = s.org
	c.ConnectionID = s.connection
	c.AppointmentID = r.AppointmentID
	c.ExternalIDHMAC = h
	c.ConsentPurpose = s.base.purpose
	c.ProfileSHA256 = digest(s.profile)
	c.NotBefore = r.NotBefore
	c.ExpiresAt = r.ExpiresAt
	c.Message = channels.Message{TenantID: s.tenant, ChannelCode: "whatsapp", Direction: channels.DirectionOut, DeliveryKey: c.DeliveryKey, ExternalID: r.Recipient, Text: string(body)}
	c.MessageSHA256, e = outbounddelivery.MessageSHA256(c.Message)
	if e != nil {
		return ScheduleContext{}, ErrApproval
	}
	return c, nil
}
func (s *ScheduleApprovals) Prepare(ctx context.Context, p identity.Principal, r ScheduleRequest) (ScheduleContext, error) {
	if !s.authorized(p, "notification:request") {
		return ScheduleContext{}, ErrApproval
	}
	tx, e := s.base.pool.Begin(ctx)
	if e != nil {
		return ScheduleContext{}, e
	}
	defer tx.Rollback(ctx)
	out, e := s.source(ctx, tx, r, true)
	if e != nil {
		return out, e
	}
	return out, tx.Commit(ctx)
}
func (s *ScheduleApprovals) guard(ctx context.Context, tx pgx.Tx, c ScheduleContext) error {
	rebuilt, e := s.source(ctx, tx, c.Request, false)
	if e != nil {
		return e
	}
	raw, _ := json.Marshal(c)
	again, _ := json.Marshal(rebuilt)
	_, h, e := approval.CanonicalPayload(raw)
	_, other, e2 := approval.CanonicalPayload(again)
	if e != nil || e2 != nil || h != other {
		return ErrApproval
	}
	var cancelled bool
	if e = tx.QueryRow(ctx, `select exists(select 1 from communication.whatsapp_schedule_cancellation where tenant_id=$1 and delivery_key=$2)`, s.tenant, c.DeliveryKey).Scan(&cancelled); e != nil {
		return e
	}
	if cancelled {
		return ErrApproval
	}
	return nil
}
func (s *ScheduleApprovals) Submit(ctx context.Context, p identity.Principal, c ScheduleContext) (string, bool, error) {
	if !s.authorized(p, "notification:request") {
		return "", false, ErrApproval
	}
	raw, e := json.Marshal(c)
	if e != nil {
		return "", false, e
	}
	_, h, e := approval.CanonicalPayload(raw)
	if e != nil {
		return "", false, e
	}
	subject := c.AppointmentID
	if c.Request.CampaignID != "" {
		subject = c.LeadID
	}
	replay, e := s.human.Submit(ctx, p, postgres.HumanApprovalSpec{Request: approval.Request{TenantID: s.tenant, ID: c.DeliveryKey, Kind: approval.KindWhatsAppSchedule, SubjectID: subject, Requester: p.Subject, EvidenceSHA: h}, OrganizationID: s.org, Payload: raw}, "notification:request", func(ctx context.Context, tx pgx.Tx) error { return s.guard(ctx, tx, c) })
	return h, replay, e
}
func (s *ScheduleApprovals) request(ctx context.Context, key string) (ScheduleContext, string, error) {
	var c ScheduleContext
	var raw []byte
	var h string
	if !validScheduleKey(key) {
		return c, "", ErrApproval
	}
	e := s.base.pool.QueryRow(ctx, `select payload,evidence_sha from approval.request where tenant_id=$1 and request_id=$2 and organization_id=$3 and kind='whatsapp_schedule'`, s.tenant, key, s.org).Scan(&raw, &h)
	if e != nil {
		return c, "", ErrApproval
	}
	_, actual, e := approval.CanonicalPayload(raw)
	if e != nil || h != actual || json.Unmarshal(raw, &c) != nil || c.OrganizationID != s.org || c.ConnectionID != s.connection || c.DeliveryKey != key || c.ProfileSHA256 != digest(s.profile) {
		return c, "", ErrApproval
	}
	return c, h, nil
}
