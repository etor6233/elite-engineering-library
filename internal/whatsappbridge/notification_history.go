package whatsappbridge

// AUTHORED read-only collection over existing approvals/status owners. No send
// authorization, inferred latest confirmation, new table or customer identity.
import (
	"context"
	"elite.local/enterprise/internal/platform/identity"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5"
	"net/http"
	"net/url"
	"time"
)

var ErrNotificationHistoryLimit = errors.New("whatsappbridge: notification history exceeds bounded view")

type NotificationHistoryItem struct {
	ConfirmationEventID string             `json:"confirmation_event_id"`
	Status              NotificationStatus `json:"status"`
}
type NotificationHistory struct {
	OrganizationID string                    `json:"organization_id"`
	AppointmentID  string                    `json:"appointment_id"`
	Items          []NotificationHistoryItem `json:"items"`
}

// A consistent read-only snapshot, capped at 20 approvals per appointment. An
// overflow is explicit, never silently presented as complete history.
func (s *PostgresAppointmentApprovals) ReadNotificationHistory(ctx context.Context, p identity.Principal, organization, appointment string) (NotificationHistory, error) {
	value := NotificationHistory{OrganizationID: organization, AppointmentID: appointment, Items: []NotificationHistoryItem{}}
	if s == nil || s.pool == nil || p.Subject == "" || p.TenantID == "" || !p.Allowed("appointment:manage") || !p.AllowedOrganization(organization) || organization == "" || len(organization) > 128 || appointment == "" || len(appointment) > 128 {
		return NotificationHistory{}, ErrNotificationNotFound
	}
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead, AccessMode: pgx.ReadOnly})
	if err != nil {
		return NotificationHistory{}, ErrNotificationRead
	}
	defer tx.Rollback(ctx)
	var exists bool
	err = tx.QueryRow(ctx, `select exists(select 1 from crm.appointment a join crm.lead l on l.tenant_id=a.tenant_id and l.lead_id=a.lead_id and l.organization_id=a.organization_id where a.tenant_id=$1 and a.organization_id=$2 and a.appointment_id=$3)`, p.TenantID, organization, appointment).Scan(&exists)
	if err != nil {
		return NotificationHistory{}, ErrNotificationRead
	}
	if !exists {
		return NotificationHistory{}, ErrNotificationNotFound
	}
	rows, err := tx.Query(ctx, `select confirmation_event_id from communication.whatsapp_appointment_approval where tenant_id=$1 and organization_id=$2 and appointment_id=$3 and channel_code='whatsapp' order by confirmation_event_id limit 21`, p.TenantID, organization, appointment)
	if err != nil {
		return NotificationHistory{}, ErrNotificationRead
	}
	ids := []string{}
	for rows.Next() {
		var id string
		if err = rows.Scan(&id); err != nil {
			break
		}
		ids = append(ids, id)
	}
	rowErr := rows.Err()
	rows.Close()
	if err != nil || rowErr != nil {
		return NotificationHistory{}, ErrNotificationRead
	}
	if len(ids) > 20 {
		return NotificationHistory{}, ErrNotificationHistoryLimit
	}
	for _, id := range ids {
		status, err := readNotificationStatus(ctx, tx, p, organization, appointment, id)
		if err != nil {
			return NotificationHistory{}, err
		}
		value.Items = append(value.Items, NotificationHistoryItem{ConfirmationEventID: id, Status: status})
	}
	if tx.Commit(ctx) != nil {
		return NotificationHistory{}, ErrNotificationRead
	}
	return value, nil
}

func (m *AppointmentNotificationModule) registerNotificationHistory(mux *http.ServeMux, verifier identity.Verifier) {
	mux.HandleFunc("GET /v1/franchise/appointments/{id}/whatsapp-confirmations", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		p, ok := m.notificationPrincipal(ctx, w, r, verifier)
		if !ok {
			return
		}
		q, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil || len(r.URL.RawQuery) > 1024 || len(q) != 1 || len(q["organization_id"]) != 1 || q.Get("organization_id") == "" || len(q.Get("organization_id")) > 128 || len(r.PathValue("id")) > 128 || r.ContentLength != 0 {
			notificationProblem(w, 400, "INVALID_QUERY")
			return
		}
		if !p.AllowedOrganization(q.Get("organization_id")) {
			notificationProblem(w, 403, "ORGANIZATION_FORBIDDEN")
			return
		}
		value, err := m.approvals.ReadNotificationHistory(ctx, p, q.Get("organization_id"), r.PathValue("id"))
		switch {
		case errors.Is(err, ErrNotificationNotFound):
			notificationProblem(w, 404, "NOTIFICATION_NOT_FOUND")
			return
		case errors.Is(err, ErrNotificationIntegrity):
			notificationProblem(w, 409, "NOTIFICATION_EVIDENCE_MISMATCH")
			return
		case errors.Is(err, ErrNotificationHistoryLimit):
			notificationProblem(w, 409, "NOTIFICATION_HISTORY_LIMIT")
			return
		case err != nil:
			notificationProblem(w, 503, "NOTIFICATION_READ_UNAVAILABLE")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-store")
		_ = json.NewEncoder(w).Encode(value)
	})
}
