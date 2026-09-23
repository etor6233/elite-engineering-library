package whatsappbridge

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"elite.local/enterprise/internal/channels"
	"elite.local/enterprise/internal/contactidentity"
	"elite.local/enterprise/internal/outbounddelivery"
	"elite.local/enterprise/internal/platform/identity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrApproval = errors.New("whatsappbridge: appointment approval absent, divergent or no longer authorized")

// AppointmentApprovalCommand references proof from the authorized workflow.
// These hashes do not themselves prove consent or approve a Meta template.
// The caller must supply a verified Principal, never JSON-authored identity.
type AppointmentApprovalCommand struct {
	Message               channels.Message
	OrganizationID        string
	AppointmentID         string
	ConfirmationEventID   string
	AppointmentVersion    int64
	BindingVersion        int64
	PolicyVersion         string
	ConsentID             string
	ConsentPurpose        string
	ConsentEvidenceSHA256 string
	ProfileSHA256         string
	EvidenceSHA256        string
	ExpiresAt             time.Time
}

type PostgresAppointmentApprovals struct {
	pool    *pgxpool.Pool
	hmacKey []byte
	purpose string
	policy  string
}

func NewPostgresAppointmentApprovals(pool *pgxpool.Pool, key []byte, purpose, policy string) (*PostgresAppointmentApprovals, error) {
	if pool == nil || len(key) < 32 || strings.TrimSpace(purpose) == "" || len(purpose) > 128 || strings.TrimSpace(policy) == "" || len(policy) > 128 {
		return nil, ErrApproval
	}
	return &PostgresAppointmentApprovals{pool: pool, hmacKey: append([]byte(nil), key...), purpose: purpose, policy: policy}, nil
}

// One business event has one notification identity, regardless of recipient or
// retries. A different key cannot bypass the durable fence for the same event.
func AppointmentConfirmationDeliveryKey(tenant, event string) string {
	return "wa-appointment-confirmed:" + digest([]byte(tenant+"\x00"+event))
}

func (s *PostgresAppointmentApprovals) Approve(ctx context.Context, p identity.Principal, c AppointmentApprovalCommand) (Approval, error) {
	empty := Approval{}
	now := time.Now().UTC()
	if s == nil || c.ConsentPurpose != s.purpose || c.PolicyVersion != s.policy {
		return empty, ErrApproval
	}
	if c.ConsentID == "" || strings.TrimSpace(c.ConsentPurpose) == "" || len(c.ConsentPurpose) > 128 {
		return empty, ErrApproval
	}
	if s == nil || p.Subject == "" || len(p.Subject) > 255 || p.TenantID != c.Message.TenantID || !p.Allowed("appointment:manage") || !p.AllowedOrganization(c.OrganizationID) || c.Message.ChannelCode != "whatsapp" || c.AppointmentID == "" || c.ConfirmationEventID == "" || c.Message.DeliveryKey != AppointmentConfirmationDeliveryKey(p.TenantID, c.ConfirmationEventID) || c.AppointmentVersion < 1 || c.BindingVersion < 1 || strings.TrimSpace(c.PolicyVersion) == "" || len(c.PolicyVersion) > 128 || !validDigest(c.ConsentEvidenceSHA256) || !validDigest(c.ProfileSHA256) || !validDigest(c.EvidenceSHA256) || !c.ExpiresAt.After(now) || c.ExpiresAt.After(now.Add(24*time.Hour)) {
		return empty, ErrApproval
	}
	hash, err := outbounddelivery.MessageSHA256(c.Message)
	if err != nil {
		return empty, ErrApproval
	}
	if _, err = requestFor(c.Message); err != nil {
		return empty, ErrApproval
	}
	recipient, err := contactidentity.ExternalDigest(s.hmacKey, p.TenantID, "whatsapp", c.Message.ExternalID)
	if err != nil {
		return empty, ErrApproval
	}
	// The fingerprint is never persisted in clear: it binds the complete approved
	// request and actor without storing recipient/template body as database PII.
	fingerprintBytes, err := json.Marshal(struct {
		Actor   string
		Command AppointmentApprovalCommand
	}{p.Subject, c})
	if err != nil {
		return empty, ErrApproval
	}
	fingerprint := digest(fingerprintBytes)
	// A single INSERT..SELECT observes only committed appointment+outbox+binding.
	// It does not own or mutate those aggregates. No provider effect happens here.
	_, err = s.pool.Exec(ctx, `insert into communication.whatsapp_appointment_approval
 (tenant_id,delivery_key,appointment_id,appointment_version,organization_id,confirmation_event_id,external_id_hmac,binding_version,lead_id,subject_id,policy_version,consent_evidence_sha256,message_sha256,profile_sha256,approval_request_sha256,evidence_sha256,approved_by,expires_at,consent_id,consent_purpose)
 select a.tenant_id,$2,a.appointment_id,a.version,a.organization_id,t.transition_id,b.external_id_hmac,b.version,b.lead_id,b.subject_id,b.policy_version,$10,$11,$12,$13,$14,$15,$16,consent.consent_id,consent.purpose_code
 from crm.appointment a
 join crm.lead l on l.tenant_id=a.tenant_id and l.lead_id=a.lead_id and l.organization_id=a.organization_id
 join crm.appointment_transition t on t.tenant_id=a.tenant_id and t.appointment_id=a.appointment_id and t.transition_id=$5 and t.from_state='requested' and t.to_state='confirmed'
 join platform.outbox_event e on e.tenant_id=t.tenant_id and e.event_id=t.transition_id and e.aggregate_type='appointment' and e.aggregate_id=a.appointment_id and e.aggregate_version=a.version and e.event_type='appointment.confirmed' and e.schema_version=1
 join communication.contact_channel_binding b on b.tenant_id=a.tenant_id and b.lead_id=a.lead_id and b.channel_code='whatsapp' and b.external_id_hmac=$7
 join crm.consent_evidence consent on consent.tenant_id=a.tenant_id and consent.lead_id=a.lead_id and consent.consent_id=$17 and consent.purpose_code=$18 and consent.policy_version=b.policy_version and consent.decision='granted' and consent.evidence_sha256_hex=$10 and consent.occurred_at<=statement_timestamp()
 where a.tenant_id=$1 and a.organization_id=$3 and a.appointment_id=$4 and a.version=$6 and a.state='confirmed' and a.starts_at>statement_timestamp()
 and b.state='active' and b.pii_allowed and b.version=$8 and b.policy_version=$9 and b.effective_at<=statement_timestamp()
 and not exists(select 1 from crm.consent_evidence newer where newer.tenant_id=consent.tenant_id and newer.lead_id=consent.lead_id and newer.purpose_code=consent.purpose_code and newer.consent_id<>consent.consent_id and newer.occurred_at>=consent.occurred_at and newer.occurred_at<=statement_timestamp())
 on conflict do nothing`, p.TenantID, c.Message.DeliveryKey, c.OrganizationID, c.AppointmentID, c.ConfirmationEventID, c.AppointmentVersion, recipient, c.BindingVersion, c.PolicyVersion, c.ConsentEvidenceSHA256, hash, c.ProfileSHA256, fingerprint, c.EvidenceSHA256, p.Subject, c.ExpiresAt.UTC(), c.ConsentID, c.ConsentPurpose)
	if err != nil {
		return empty, ErrApproval
	}
	var stored string
	if s.pool.QueryRow(ctx, `select approval_request_sha256 from communication.whatsapp_appointment_approval where tenant_id=$1 and delivery_key=$2`, p.TenantID, c.Message.DeliveryKey).Scan(&stored) != nil || stored != fingerprint {
		return empty, ErrApproval
	}
	return s.ResolveWhatsAppApproval(ctx, p.TenantID, c.Message.DeliveryKey)
}

// Resolve rechecks the current binding and appointment. It does not renew an
// approval on expiry, version change, cancellation, reassignment or revocation.
// It is a pre-send snapshot, not a promise to recall an already in-flight send.
func (s *PostgresAppointmentApprovals) ResolveWhatsAppApproval(ctx context.Context, tenant, key string) (Approval, error) {
	value := Approval{}
	if s == nil || strings.TrimSpace(tenant) == "" || strings.TrimSpace(key) == "" {
		return value, ErrApproval
	}
	err := s.pool.QueryRow(ctx, `select g.message_sha256,g.profile_sha256,g.evidence_sha256,g.approved_by,g.approved_at,g.expires_at
 from communication.whatsapp_appointment_approval g
 join crm.appointment a on a.tenant_id=g.tenant_id and a.appointment_id=g.appointment_id and a.organization_id=g.organization_id and a.lead_id=g.lead_id and a.version=g.appointment_version and a.state='confirmed' and a.starts_at>statement_timestamp()
 join crm.lead l on l.tenant_id=a.tenant_id and l.lead_id=a.lead_id and l.organization_id=a.organization_id
 join communication.contact_channel_binding b on b.tenant_id=g.tenant_id and b.channel_code=g.channel_code and b.external_id_hmac=g.external_id_hmac and b.version=g.binding_version and b.lead_id=g.lead_id and b.subject_id=g.subject_id and b.policy_version=g.policy_version and b.state='active' and b.pii_allowed and b.effective_at<=statement_timestamp()
 join crm.consent_evidence consent on consent.tenant_id=g.tenant_id and consent.lead_id=g.lead_id and consent.consent_id=g.consent_id and consent.purpose_code=g.consent_purpose and consent.policy_version=g.policy_version and consent.evidence_sha256_hex=g.consent_evidence_sha256 and consent.decision='granted' and consent.occurred_at<=statement_timestamp()
 where g.tenant_id=$1 and g.delivery_key=$2 and g.policy_version=$3 and g.consent_purpose=$4 and g.approved_at<=statement_timestamp() and g.expires_at>statement_timestamp()
 and not exists(select 1 from crm.consent_evidence newer where newer.tenant_id=consent.tenant_id and newer.lead_id=consent.lead_id and newer.purpose_code=consent.purpose_code and newer.consent_id<>consent.consent_id and newer.occurred_at>=consent.occurred_at and newer.occurred_at<=statement_timestamp())`, tenant, key, s.policy, s.purpose).Scan(&value.MessageSHA256, &value.ProfileSHA256, &value.EvidenceSHA256, &value.ApprovedBy, &value.NotBefore, &value.ExpiresAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Approval{}, ErrApproval
		}
		return Approval{}, ErrApproval
	}
	return value, nil
}
