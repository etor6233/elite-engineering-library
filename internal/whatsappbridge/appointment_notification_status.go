package whatsappbridge

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"elite.local/enterprise/internal/platform/identity"
	"github.com/jackc/pgx/v5"
)

var (
	ErrNotificationNotFound  = errors.New("whatsappbridge: scoped notification not found")
	ErrNotificationRead      = errors.New("whatsappbridge: notification read unavailable")
	ErrNotificationIntegrity = errors.New("whatsappbridge: notification evidence mismatch")
	notificationEventID      = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
)

// NotificationStatus is a bounded local snapshot, not a provider delivery claim
// or permission to send. Hashes, recipient, actor and template are not exposed.
type NotificationStatus struct {
	DeliveryKey            string     `json:"delivery_key"`
	FenceState             string     `json:"fence_state"`
	DeliveryStatus         string     `json:"delivery_status"`
	ProviderTimestamp      *int64     `json:"provider_timestamp,omitempty"`
	ProviderEventCount     int64      `json:"provider_event_count"`
	ApprovalExpiresAt      time.Time  `json:"approval_expires_at"`
	ApprovalExpired        bool       `json:"approval_expired"`
	AcceptedAt             *time.Time `json:"accepted_at,omitempty"`
	UpdatedAt              *time.Time `json:"updated_at,omitempty"`
	ObservedAt             time.Time  `json:"observed_at"`
	ReconciliationRequired bool       `json:"reconciliation_required"`
}

// ReadNotificationStatus deliberately does not call ResolveWhatsAppApproval:
// an expired/revoked sending grant must not hide the historical outcome from
// a currently authorized operator. Current appointment/lead/org scope is checked.
func (s *PostgresAppointmentApprovals) ReadNotificationStatus(ctx context.Context, p identity.Principal, organization, appointment, event string) (NotificationStatus, error) {
	if s == nil || s.pool == nil {
		return NotificationStatus{}, ErrNotificationNotFound
	}
	return readNotificationStatus(ctx, s.pool, p, organization, appointment, event)
}

type notificationQueryer interface {
	QueryRow(context.Context, string, ...any) pgx.Row
}

func readNotificationStatus(ctx context.Context, db notificationQueryer, p identity.Principal, organization, appointment, event string) (NotificationStatus, error) {
	value := NotificationStatus{}
	if db == nil || p.Subject == "" || p.TenantID == "" || !p.AllowedOrganization(organization) || !((appointment != "" && p.Allowed("appointment:manage") && notificationEventID.MatchString(event)) || (appointment == "" && p.Allowed("whatsapp:approve") && strings.HasPrefix(event, "wa-reply:") && validDigest(strings.TrimPrefix(event, "wa-reply:"))) || (appointment == "" && p.Allowed("notification:read") && strings.HasPrefix(event, "wa-schedule:") && validDigest(strings.TrimPrefix(event, "wa-schedule:")))) {
		return value, ErrNotificationNotFound
	}
	var requestHash, recipientHash, expectedHash, expectedRecipient, providerHash, evidenceHash string
	var lease *time.Time
	err := db.QueryRow(ctx, `select g.delivery_key,g.expires_at,g.expires_at<=statement_timestamp(),
 coalesce(d.state,'not_started'),d.accepted_at,d.updated_at,statement_timestamp(),d.locked_until,
 coalesce(d.request_sha256_hex,''),coalesce(d.recipient_hmac,''),g.message_sha256,g.external_id_hmac,
 coalesce(d.provider_message_hmac,''),coalesce(d.evidence_sha256_hex,''), totals.event_count,totals.last_time,
 case when totals.event_count=0 then 'not_observed_by_this_reader' else latest.status end
 from communication.whatsapp_delivery_approval g
 join crm.lead l on l.tenant_id=g.tenant_id and l.lead_id=g.lead_id and l.organization_id=g.organization_id
 left join communication.outbound_delivery d on d.tenant_id=g.tenant_id and d.channel_code=g.channel_code and d.delivery_key=g.delivery_key
 left join lateral (select count(*) event_count,max(provider_timestamp) last_time from communication.whatsapp_status_observation s where s.tenant_id=g.tenant_id and s.delivery_key=g.delivery_key) totals on true
 left join lateral (select case when count(distinct provider_status)=1 then 'observed_'||min(provider_status) else 'ambiguous_latest_timestamp' end status from communication.whatsapp_status_observation s where s.tenant_id=g.tenant_id and s.delivery_key=g.delivery_key and s.provider_timestamp=totals.last_time) latest on true
 where g.tenant_id=$1 and g.organization_id=$2 and g.appointment_id=$3 and g.approval_event_id=$4 and g.channel_code='whatsapp'`, p.TenantID, organization, appointment, event).Scan(&value.DeliveryKey, &value.ApprovalExpiresAt, &value.ApprovalExpired, &value.FenceState, &value.AcceptedAt, &value.UpdatedAt, &value.ObservedAt, &lease, &requestHash, &recipientHash, &expectedHash, &expectedRecipient, &providerHash, &evidenceHash, &value.ProviderEventCount, &value.ProviderTimestamp, &value.DeliveryStatus)
	if errors.Is(err, pgx.ErrNoRows) {
		return NotificationStatus{}, ErrNotificationNotFound
	}
	if err != nil {
		return NotificationStatus{}, ErrNotificationRead
	}
	if value.FenceState != "not_started" && (requestHash != expectedHash || recipientHash != expectedRecipient) {
		return NotificationStatus{}, ErrNotificationIntegrity
	}
	switch value.FenceState {
	case "not_started":
	case "accepted":
		if value.AcceptedAt == nil || value.AcceptedAt.IsZero() || !validDigest(providerHash) || !validDigest(evidenceHash) {
			return NotificationStatus{}, ErrNotificationIntegrity
		}
	case "sending":
		// Read-only: report stale lease without changing state or creating events.
		value.ReconciliationRequired = lease == nil || !lease.After(value.ObservedAt)
	case "unknown":
		value.ReconciliationRequired = true
	case "failed_terminal":
		if !validDigest(evidenceHash) {
			return NotificationStatus{}, ErrNotificationIntegrity
		}
	default:
		return NotificationStatus{}, ErrNotificationIntegrity
	}
	if value.ProviderEventCount > 0 && value.FenceState != "accepted" {
		return NotificationStatus{}, ErrNotificationIntegrity
	}
	return value, nil
}

func (m *AppointmentNotificationModule) registerNotificationStatus(mux *http.ServeMux, verifier identity.Verifier) {
	mux.HandleFunc("GET /v1/franchise/appointments/{id}/whatsapp-confirmation", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		p, ok := m.notificationPrincipal(ctx, w, r, verifier)
		if !ok {
			return
		}
		query, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil || len(r.URL.RawQuery) > 2048 || len(query) != 2 || len(query["organization_id"]) != 1 || len(query["confirmation_event_id"]) != 1 || query.Get("organization_id") == "" || !notificationEventID.MatchString(query.Get("confirmation_event_id")) || r.ContentLength != 0 {
			notificationProblem(w, 400, "INVALID_QUERY")
			return
		}
		organization := query.Get("organization_id")
		if !p.AllowedOrganization(organization) {
			notificationProblem(w, 403, "ORGANIZATION_FORBIDDEN")
			return
		}
		value, err := m.approvals.ReadNotificationStatus(ctx, p, organization, r.PathValue("id"), query.Get("confirmation_event_id"))
		if errors.Is(err, ErrNotificationNotFound) {
			notificationProblem(w, 404, "NOTIFICATION_NOT_FOUND")
			return
		}
		if errors.Is(err, ErrNotificationIntegrity) {
			notificationProblem(w, 409, "NOTIFICATION_EVIDENCE_MISMATCH")
			return
		}
		if err != nil {
			notificationProblem(w, 503, "NOTIFICATION_READ_UNAVAILABLE")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-store")
		_ = json.NewEncoder(w).Encode(value)
	})
}
