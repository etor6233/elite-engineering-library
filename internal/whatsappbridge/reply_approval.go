package whatsappbridge

// AUTHORED typed projection of conversation/contact/consent and the shared
// durable approval owner. It defines no consent or business approval policy.
import (
	"context"
	"encoding/json"
	"time"

	"elite.local/enterprise/internal/approval"
	"elite.local/enterprise/internal/channels"
	"elite.local/enterprise/internal/contactidentity"
	"elite.local/enterprise/internal/outbounddelivery"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/postgres"
	"github.com/jackc/pgx/v5"
)

type ReplyContext struct {
	Schema                string           `json:"schema"`
	OrganizationID        string           `json:"organization_id"`
	ConnectionID          string           `json:"connection_id"`
	SourceEventID         string           `json:"source_event_id"`
	ProviderMessageID     string           `json:"provider_message_id"`
	ExternalIDHMAC        string           `json:"external_id_hmac"`
	BindingVersion        int64            `json:"binding_version"`
	LeadID                string           `json:"lead_id"`
	SubjectID             string           `json:"subject_id"`
	PolicyVersion         string           `json:"policy_version"`
	ConsentID             string           `json:"consent_id"`
	ConsentPurpose        string           `json:"consent_purpose"`
	ConsentEvidenceSHA256 string           `json:"consent_evidence_sha256"`
	ProfileSHA256         string           `json:"profile_sha256"`
	MessageSHA256         string           `json:"message_sha256"`
	ExpiresAt             time.Time        `json:"expires_at"`
	Message               channels.Message `json:"message"`
}

type ReplyProposal struct {
	Status        *NotificationStatus `json:"status,omitempty"`
	RequestID     string              `json:"request_id"`
	State         string              `json:"state"`
	PayloadSHA256 string              `json:"payload_sha256"`
	Context       ReplyContext        `json:"context"`
}

type PostgresReplyApprovals struct {
	base                                      *PostgresAppointmentApprovals
	tenant, organization, connection, profile string
}

func NewPostgresReplyApprovals(base *PostgresAppointmentApprovals, tenant, organization, connection string, profile json.RawMessage) (*PostgresReplyApprovals, error) {
	if base == nil || base.pool == nil || !notificationEventID.MatchString(tenant) || organization == "" || !webhookConnection.MatchString(connection) || !json.Valid(profile) || len(profile) > 32768 {
		return nil, ErrApproval
	}
	return &PostgresReplyApprovals{base: base, tenant: tenant, organization: organization, connection: connection, profile: digest(profile)}, nil
}
func ReplyDeliveryKey(tenant, connection, sourceMessage string) string {
	return "wa-reply:" + digest([]byte(tenant+"\x00"+connection+"\x00"+sourceMessage))
}

func (s *PostgresReplyApprovals) validPrincipal(p identity.Principal, permission string) bool {
	return s != nil && p.TenantID == s.tenant && p.Subject != "" && len(p.Subject) <= 128 && p.Allowed(permission) && p.AllowedOrganization(s.organization)
}

// buildProposal is called only from the verified raw-inbox router. It projects
// the exact terminal conversation turn and current explicit contact/consent.
func (s *PostgresReplyApprovals) buildProposal(ctx context.Context, p identity.Principal, eventID string, in channels.Message) (ReplyContext, error) {
	var c ReplyContext
	if !s.validPrincipal(p, "whatsapp:process") || in.TenantID != s.tenant || in.ChannelCode != "whatsapp" || in.Direction != channels.DirectionIn || in.ThreadID != conversationThread(s.connection) || in.Validate() != nil {
		return c, ErrApproval
	}
	recipient, err := contactidentity.ExternalDigest(s.base.hmacKey, s.tenant, "whatsapp", in.ExternalID)
	if err != nil {
		return c, ErrApproval
	}
	var text, responseHash string
	err = s.base.pool.QueryRow(ctx, `select t.assistant_text,t.response_sha256_hex,b.version,b.lead_id,b.subject_id,b.policy_version,consent.consent_id,consent.evidence_sha256_hex
 from communication.conversation_turn t
 join communication.contact_channel_binding b on b.tenant_id=t.tenant_id and b.channel_code=t.channel_code and b.external_id_hmac=$7
 join crm.lead l on l.tenant_id=b.tenant_id and l.lead_id=b.lead_id and l.organization_id=$8
 join integration.provider_connection pc on pc.tenant_id=t.tenant_id and pc.connection_id=$9 and pc.provider_code='meta-whatsapp' and pc.organization_id=l.organization_id and pc.state='active'
 join integration.webhook_event inbox on inbox.tenant_id=pc.tenant_id and inbox.connection_id=pc.connection_id and inbox.provider_code=pc.provider_code and inbox.provider_event_id=$10 and inbox.event_type='whatsapp.raw_webhook.received.v1' and inbox.state in ('received','processed')
 join crm.consent_evidence consent on consent.tenant_id=b.tenant_id and consent.lead_id=b.lead_id and consent.purpose_code=$11 and consent.policy_version=b.policy_version and consent.decision='granted' and consent.occurred_at<=statement_timestamp()
 where t.tenant_id=$1 and t.channel_code='whatsapp' and t.provider_message_id=$2 and t.external_id=$3 and t.thread_id=$4 and t.occurred_at=$5 and t.user_text=$6
 and t.state in ('completed','handed_off') and t.expires_at>statement_timestamp() and t.occurred_at<=statement_timestamp() and t.occurred_at+interval '24 hours'>statement_timestamp()+interval '40 seconds'
 and b.state='active' and b.pii_allowed and b.effective_at<=statement_timestamp() and b.policy_version=$12
 and not exists(select 1 from crm.consent_evidence newer where newer.tenant_id=consent.tenant_id and newer.lead_id=consent.lead_id and newer.purpose_code=consent.purpose_code and newer.consent_id<>consent.consent_id and newer.occurred_at>=consent.occurred_at and newer.occurred_at<=statement_timestamp())`, s.tenant, in.ProviderMessageID, in.ExternalID, in.ThreadID, in.OccurredAt, in.Text, recipient, s.organization, s.connection, eventID, s.base.purpose, s.base.policy).Scan(&text, &responseHash, &c.BindingVersion, &c.LeadID, &c.SubjectID, &c.PolicyVersion, &c.ConsentID, &c.ConsentEvidenceSHA256)
	if err != nil || digest([]byte(text)) != responseHash {
		return ReplyContext{}, ErrApproval
	}
	request := ReplyRequest{Kind: "text_reply", Recipient: in.ExternalID, Text: text, SourceMessageID: in.ProviderMessageID, LastInboundAt: in.OccurredAt.Unix(), WindowExpiresAt: in.OccurredAt.Unix() + 86400}
	raw, _ := json.Marshal(request)
	c.Message = channels.Message{ChannelCode: "whatsapp", TenantID: s.tenant, ExternalID: in.ExternalID, ThreadID: in.ThreadID, Direction: channels.DirectionOut, DeliveryKey: ReplyDeliveryKey(s.tenant, s.connection, in.ProviderMessageID), Text: string(raw)}
	if _, err = replyRequestFor(c.Message); err != nil {
		return ReplyContext{}, ErrApproval
	}
	c.MessageSHA256, err = outbounddelivery.MessageSHA256(c.Message)
	if err != nil {
		return ReplyContext{}, ErrApproval
	}
	c.Schema = "elite-whatsapp-reply-approval/v1"
	c.OrganizationID = s.organization
	c.ConnectionID = s.connection
	c.SourceEventID = eventID
	c.ProviderMessageID = in.ProviderMessageID
	c.ExternalIDHMAC = recipient
	c.ConsentPurpose = s.base.purpose
	c.ProfileSHA256 = s.profile
	c.ExpiresAt = time.Unix(request.WindowExpiresAt, 0).UTC()
	return c, nil
}

// current binds approval to the current tenant/org/connection, source turn,
// contact version, explicit latest consent and unexpired customer-service window.
// It is a pre-dispatch snapshot, not a promise to recall an in-flight request.
func (s *PostgresReplyApprovals) current(ctx context.Context, db statusQuery, c ReplyContext) error {
	if s == nil || c.Schema != "elite-whatsapp-reply-approval/v1" || c.OrganizationID != s.organization || c.ConnectionID != s.connection || c.ProfileSHA256 != s.profile || c.PolicyVersion != s.base.policy || c.ConsentPurpose != s.base.purpose || c.Message.TenantID != s.tenant || c.Message.DeliveryKey != ReplyDeliveryKey(s.tenant, s.connection, c.ProviderMessageID) || c.Message.ThreadID != conversationThread(s.connection) {
		return ErrApproval
	}
	hash, err := outbounddelivery.MessageSHA256(c.Message)
	if err != nil || hash != c.MessageSHA256 {
		return ErrApproval
	}
	recipient, err := contactidentity.ExternalDigest(s.base.hmacKey, s.tenant, "whatsapp", c.Message.ExternalID)
	if err != nil || recipient != c.ExternalIDHMAC {
		return ErrApproval
	}
	raw, err := replyRequestFor(c.Message)
	if err != nil {
		return ErrApproval
	}
	var request ReplyRequest
	_ = json.Unmarshal(raw, &request)
	if request.SourceMessageID != c.ProviderMessageID || !c.ExpiresAt.Equal(time.Unix(request.WindowExpiresAt, 0)) {
		return ErrApproval
	}
	var found bool
	err = db.QueryRow(ctx, `select true from communication.contact_channel_binding b
 join crm.lead l on l.tenant_id=b.tenant_id and l.lead_id=b.lead_id and l.organization_id=$3
 join integration.provider_connection pc on pc.tenant_id=b.tenant_id and pc.connection_id=$4 and pc.provider_code='meta-whatsapp' and pc.organization_id=l.organization_id and pc.state='active'
 join communication.conversation_turn t on t.tenant_id=b.tenant_id and t.channel_code=b.channel_code and t.provider_message_id=$5 and t.external_id=$6 and t.thread_id=$7 and t.occurred_at=to_timestamp($8) and t.assistant_text=$9 and t.state in ('completed','handed_off') and t.expires_at>statement_timestamp()
 join crm.consent_evidence consent on consent.tenant_id=b.tenant_id and consent.lead_id=b.lead_id and consent.consent_id=$10 and consent.purpose_code=$11 and consent.policy_version=b.policy_version and consent.evidence_sha256_hex=$12 and consent.decision='granted' and consent.occurred_at<=statement_timestamp()
 where b.tenant_id=$1 and b.channel_code='whatsapp' and b.external_id_hmac=$2 and b.version=$13 and b.lead_id=$14 and b.subject_id=$15 and b.policy_version=$16 and b.state='active' and b.pii_allowed and b.effective_at<=statement_timestamp()
 and t.occurred_at<=statement_timestamp() and t.occurred_at+interval '24 hours'>statement_timestamp()+interval '40 seconds'
 and not exists(select 1 from crm.consent_evidence newer where newer.tenant_id=consent.tenant_id and newer.lead_id=consent.lead_id and newer.purpose_code=consent.purpose_code and newer.consent_id<>consent.consent_id and newer.occurred_at>=consent.occurred_at and newer.occurred_at<=statement_timestamp())`, s.tenant, c.ExternalIDHMAC, s.organization, s.connection, c.ProviderMessageID, c.Message.ExternalID, c.Message.ThreadID, request.LastInboundAt, request.Text, c.ConsentID, c.ConsentPurpose, c.ConsentEvidenceSHA256, c.BindingVersion, c.LeadID, c.SubjectID, c.PolicyVersion).Scan(&found)
	if err != nil || !found {
		return ErrApproval
	}
	return nil
}

// Read never derives tenant/org from the request body. It returns the exact
// stored proposal for the human to review, together with its approval hash.
func (s *PostgresReplyApprovals) Read(ctx context.Context, p identity.Principal, key string) (ReplyProposal, error) {
	var value ReplyProposal
	var raw []byte
	if !s.validPrincipal(p, "whatsapp:approve") {
		return value, ErrApproval
	}
	err := s.base.pool.QueryRow(ctx, `select request_id,state,evidence_sha,payload from approval.request where tenant_id=$1 and request_id=$2 and organization_id=$3 and kind='whatsapp_reply'`, s.tenant, key, s.organization).Scan(&value.RequestID, &value.State, &value.PayloadSHA256, &raw)
	_, actualHash, hashErr := approval.CanonicalPayload(raw)
	if err != nil || hashErr != nil || actualHash != value.PayloadSHA256 || json.Unmarshal(raw, &value.Context) != nil {
		return ReplyProposal{}, ErrApproval
	}
	if value.State == "approved" {
		status, err := readNotificationStatus(ctx, s.base.pool, p, s.organization, "", key)
		if err != nil {
			return ReplyProposal{}, err
		}
		value.Status = &status
	}
	return value, nil
}

func (s *PostgresReplyApprovals) ResolveWhatsAppApproval(ctx context.Context, tenant, key string) (Approval, error) {
	var value Approval
	var c ReplyContext
	var raw []byte
	if s == nil || tenant != s.tenant {
		return value, ErrApproval
	}
	err := s.base.pool.QueryRow(ctx, `select r.payload,r.evidence_sha,d.reviewer,r.decided_at from approval.request r join approval.decision d on d.tenant_id=r.tenant_id and d.request_id=r.request_id and d.approved and d.reviewer<>r.requester where r.tenant_id=$1 and r.request_id=$2 and r.organization_id=$3 and r.kind='whatsapp_reply' and r.state='approved' and (select count(*) from approval.decision d2 where d2.tenant_id=r.tenant_id and d2.request_id=r.request_id)=1`, tenant, key, s.organization).Scan(&raw, &value.EvidenceSHA256, &value.ApprovedBy, &value.NotBefore)
	_, actualHash, hashErr := approval.CanonicalPayload(raw)
	if err != nil || hashErr != nil || actualHash != value.EvidenceSHA256 || json.Unmarshal(raw, &c) != nil || s.current(ctx, s.base.pool, c) != nil {
		return Approval{}, ErrApproval
	}
	value.MessageSHA256 = c.MessageSHA256
	value.ProfileSHA256 = c.ProfileSHA256
	value.ExpiresAt = c.ExpiresAt
	return value, nil
}

// Propose stores a pending request only; the worker principal cannot approve it.
func (s *PostgresReplyApprovals) Propose(ctx context.Context, p identity.Principal, event string, in channels.Message) (ReplyProposal, error) {
	c, err := s.buildProposal(ctx, p, event, in)
	if err != nil {
		// The runtime already persisted a governed handoff (no model/tool effect for
		// unresolved identity). It is deliberately not made into an outbound request.
		if s.validPrincipal(p, "whatsapp:process") && in.TenantID == s.tenant && in.ThreadID == conversationThread(s.connection) {
			var found bool
			if s.base.pool.QueryRow(ctx, `select true from communication.conversation_turn where tenant_id=$1 and channel_code='whatsapp' and provider_message_id=$2 and external_id=$3 and thread_id=$4 and occurred_at=$5 and user_text=$6 and state='handed_off'`, s.tenant, in.ProviderMessageID, in.ExternalID, in.ThreadID, in.OccurredAt, in.Text).Scan(&found) == nil && found {
				return ReplyProposal{State: "HANDOFF_REVIEW_REQUIRED"}, nil
			}
		}
		return ReplyProposal{}, err
	}
	raw, _ := json.Marshal(c)
	canonical, hash, err := approval.CanonicalPayload(raw)
	if err != nil {
		return ReplyProposal{}, ErrApproval
	}
	spec := postgres.HumanApprovalSpec{Request: approval.Request{TenantID: s.tenant, ID: c.Message.DeliveryKey, Kind: approval.KindWhatsAppReply, SubjectID: c.ExternalIDHMAC, Requester: p.Subject, EvidenceSHA: hash}, OrganizationID: s.organization, Payload: canonical}
	_, err = postgres.NewHumanApprovals(s.base.pool).Submit(ctx, p, spec, "whatsapp:process", func(ctx context.Context, tx pgx.Tx) error { return s.current(ctx, tx, c) })
	if err != nil {
		return ReplyProposal{}, err
	}
	var state string
	if s.base.pool.QueryRow(ctx, `select state from approval.request where tenant_id=$1 and request_id=$2`, s.tenant, c.Message.DeliveryKey).Scan(&state) != nil {
		return ReplyProposal{}, ErrApproval
	}
	return ReplyProposal{RequestID: c.Message.DeliveryKey, State: state, PayloadSHA256: hash, Context: c}, nil
}

func (s *PostgresReplyApprovals) Decide(ctx context.Context, p identity.Principal, key, hash string, approved bool, reason string) (ReplyProposal, error) {
	value, err := s.Read(ctx, p, key)
	if err != nil || value.PayloadSHA256 != hash {
		return ReplyProposal{}, ErrApproval
	}
	guard := func(ctx context.Context, tx pgx.Tx) error {
		if !approved {
			return nil
		}
		return s.current(ctx, tx, value.Context)
	}
	_, err = postgres.NewHumanApprovals(s.base.pool).Decide(ctx, p, s.tenant, key, s.organization, hash, approved, reason, "whatsapp:approve", guard)
	if err != nil {
		return ReplyProposal{}, err
	}
	return s.Read(ctx, p, key)
}
