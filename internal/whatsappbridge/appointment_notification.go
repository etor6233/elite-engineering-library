package whatsappbridge

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net/http"
	"strings"
	"time"

	"elite.local/enterprise/internal/channels"
	"elite.local/enterprise/internal/outbounddelivery"
	"elite.local/enterprise/internal/platform/identity"
)

// AppointmentNotificationRequest is an explicit operator-approved request.
// Tenant, actor, channel, purpose, policy, profile and delivery identity are
// server-derived. This endpoint neither captures consent nor confirms a turn.
type AppointmentNotificationRequest struct {
	OrganizationID        string    `json:"organization_id"`
	ConfirmationEventID   string    `json:"confirmation_event_id"`
	AppointmentVersion    int64     `json:"appointment_version"`
	BindingVersion        int64     `json:"binding_version"`
	ConsentID             string    `json:"consent_id"`
	ConsentEvidenceSHA256 string    `json:"consent_evidence_sha256"`
	EvidenceSHA256        string    `json:"evidence_sha256"`
	ExpiresAt             time.Time `json:"expires_at"`
	Recipient             string    `json:"recipient"`
	TemplateName          string    `json:"template_name"`
	LanguageCode          string    `json:"language_code"`
	BodyParameters        []string  `json:"body_parameters"`
}

// AppointmentNotificationModule follows the existing module Register contract.
// Mount on the existing protected server; it creates no listener or retry loop.
type AppointmentNotificationModule struct {
	approvals *PostgresAppointmentApprovals
	sender    *Sender
	channel   *outbounddelivery.Channel
}

func NewAppointmentNotificationModule(approvals *PostgresAppointmentApprovals, sender *Sender, store outbounddelivery.Store, receiver channels.Channel) (*AppointmentNotificationModule, error) {
	if approvals == nil || approvals.pool == nil || sender == nil || sender.TenantID == "" || sender.Tokens == nil || !json.Valid(sender.Profile) || len(sender.Profile) > 32768 || store == nil || receiver == nil || receiver.Code() != "whatsapp" {
		return nil, ErrBridge
	}
	// Freeze the caller's struct/profile and force the same durable resolver into
	// the sender. The caller cannot accidentally wire a permissive second one.
	copySender := *sender
	copySender.Profile = append(json.RawMessage(nil), sender.Profile...)
	copySender.Approvals = approvals
	return &AppointmentNotificationModule{approvals: approvals, sender: &copySender, channel: &outbounddelivery.Channel{CodeValue: "whatsapp", Receiver: receiver, Sender: &copySender, Store: store}}, nil
}

func (m *AppointmentNotificationModule) Register(mux *http.ServeMux, verifier identity.Verifier) {
	m.registerNotificationStatus(mux, verifier)
	m.registerNotificationHistory(mux, verifier)
	mux.HandleFunc("POST /v1/franchise/appointments/{id}/whatsapp-confirmation", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 45*time.Second)
		defer cancel()
		p, ok := m.notificationPrincipal(ctx, w, r, verifier)
		if !ok {
			return
		}
		media, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
		if err != nil || media != "application/json" || len(params) > 1 || (len(params) == 1 && !strings.EqualFold(params["charset"], "utf-8")) {
			notificationProblem(w, 415, "UNSUPPORTED_MEDIA_TYPE")
			return
		}
		input, err := readNotificationRequest(http.MaxBytesReader(w, r.Body, 64<<10))
		if err != nil {
			notificationProblem(w, 400, "INVALID_BODY")
			return
		}
		if !p.AllowedOrganization(input.OrganizationID) {
			notificationProblem(w, 403, "ORGANIZATION_FORBIDDEN")
			return
		}
		request, _ := json.Marshal(struct {
			Recipient      string   `json:"recipient"`
			TemplateName   string   `json:"template_name"`
			LanguageCode   string   `json:"language_code"`
			BodyParameters []string `json:"body_parameters"`
		}{input.Recipient, input.TemplateName, input.LanguageCode, input.BodyParameters})
		key := AppointmentConfirmationDeliveryKey(p.TenantID, input.ConfirmationEventID)
		message := channels.Message{TenantID: p.TenantID, ChannelCode: "whatsapp", Direction: channels.DirectionOut, DeliveryKey: key, ExternalID: input.Recipient, Text: string(request)}
		command := AppointmentApprovalCommand{Message: message, OrganizationID: input.OrganizationID, AppointmentID: r.PathValue("id"), ConfirmationEventID: input.ConfirmationEventID, AppointmentVersion: input.AppointmentVersion, BindingVersion: input.BindingVersion, PolicyVersion: m.approvals.policy, ConsentID: input.ConsentID, ConsentPurpose: m.approvals.purpose, ConsentEvidenceSHA256: input.ConsentEvidenceSHA256, ProfileSHA256: digest(m.sender.Profile), EvidenceSHA256: input.EvidenceSHA256, ExpiresAt: input.ExpiresAt.UTC()}
		if _, err := m.approvals.Approve(ctx, p, command); err != nil {
			notificationProblem(w, 409, "APPROVAL_NOT_CURRENT_OR_DIVERGENT")
			return
		}
		// Preflight occurs before Claim; the Sender still revalidates after Claim.
		// Only the existing durable Channel may invoke the provider. A crash after
		// approval is recovered with the SAME body/event, never a new identity.
		if err := m.channel.Send(ctx, message); err != nil {
			code := "DELIVERY_RECONCILIATION_REQUIRED"
			if errors.Is(err, outbounddelivery.ErrInProgress) {
				code = "DELIVERY_IN_PROGRESS"
			} else if errors.Is(err, outbounddelivery.ErrTerminal) {
				code = "DELIVERY_TERMINAL"
			}
			notificationProblem(w, 409, code)
			return
		}
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "accepted", "delivery_key": key})
	})
}

func (m *AppointmentNotificationModule) notificationPrincipal(ctx context.Context, w http.ResponseWriter, r *http.Request, verifier identity.Verifier) (identity.Principal, bool) {
	// Shared Bearer-only boundary for reads and writes; no cookie/body identity.
	values := r.Header.Values("Authorization")
	if verifier == nil || len(values) != 1 || !strings.HasPrefix(values[0], "Bearer ") || len(values[0]) > 16391 || len(strings.Fields(values[0])) != 2 {
		notificationProblem(w, 401, "UNAUTHENTICATED")
		return identity.Principal{}, false
	}
	p, err := verifier.Verify(ctx, strings.TrimPrefix(values[0], "Bearer "))
	if err != nil || p.Subject == "" || p.TenantID == "" {
		notificationProblem(w, 401, "UNAUTHENTICATED")
		return identity.Principal{}, false
	}
	if m == nil || m.sender == nil || p.TenantID != m.sender.TenantID || !p.Allowed("appointment:manage") {
		notificationProblem(w, 403, "FORBIDDEN")
		return identity.Principal{}, false
	}
	return p, true
}

func notificationProblem(w http.ResponseWriter, status int, code string) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{"type": "urn:elite:problem:" + code, "status": status, "code": code})
}

// encoding/json accepts duplicate keys by default. Reject duplicates, aliases,
// nulls and trailing data before decoding this flat, versioned command shape.
func readNotificationRequest(reader io.Reader) (AppointmentNotificationRequest, error) {
	var input AppointmentNotificationRequest
	allowed := map[string]bool{"organization_id": true, "confirmation_event_id": true, "appointment_version": true, "binding_version": true, "consent_id": true, "consent_evidence_sha256": true, "evidence_sha256": true, "expires_at": true, "recipient": true, "template_name": true, "language_code": true, "body_parameters": true}
	decoder := json.NewDecoder(reader)
	first, err := decoder.Token()
	if err != nil || first != json.Delim('{') {
		return input, ErrBridge
	}
	fields := map[string]json.RawMessage{}
	for decoder.More() {
		name, err := decoder.Token()
		key, ok := name.(string)
		if err != nil || !ok || !allowed[key] || fields[key] != nil {
			return input, ErrBridge
		}
		var raw json.RawMessage
		if decoder.Decode(&raw) != nil || bytes.Equal(bytes.TrimSpace(raw), []byte("null")) {
			return input, ErrBridge
		}
		if key == "body_parameters" {
			var values []json.RawMessage
			if json.Unmarshal(raw, &values) != nil {
				return input, ErrBridge
			}
			for _, value := range values {
				var text string
				if bytes.Equal(bytes.TrimSpace(value), []byte("null")) || json.Unmarshal(value, &text) != nil {
					return input, ErrBridge
				}
			}
		}
		fields[key] = raw
	}
	last, err := decoder.Token()
	if err != nil || last != json.Delim('}') || decoder.Decode(new(any)) != io.EOF || len(fields) != len(allowed) {
		return input, ErrBridge
	}
	raw, err := json.Marshal(fields)
	if err != nil || json.Unmarshal(raw, &input) != nil {
		return input, ErrBridge
	}
	return input, nil
}
