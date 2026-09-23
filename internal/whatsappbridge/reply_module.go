package whatsappbridge

// AUTHORED protected review/send boundary. The HTTP caller can approve/reject
// an exact stored proposal but cannot supply message text or contact identity.
import (
	"bytes"
	"context"
	"elite.local/enterprise/internal/approval"
	"elite.local/enterprise/internal/channels"
	"elite.local/enterprise/internal/outbounddelivery"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/postgres"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"
)

type ReplyModule struct {
	approvals *PostgresReplyApprovals
	sender    *Sender
	store     *postgres.OutboundDeliveryStore
	channel   *outbounddelivery.Channel
}
type VerifiedInboxChannel struct{}

func (VerifiedInboxChannel) Code() string { return "whatsapp" }
func (VerifiedInboxChannel) Receive(context.Context) ([]channels.Message, error) {
	return nil, ErrBridge
}
func (VerifiedInboxChannel) Send(context.Context, channels.Message) error { return ErrBridge }
func NewReplyModule(approvals *PostgresReplyApprovals, sender *Sender, store *postgres.OutboundDeliveryStore) (*ReplyModule, error) {
	if approvals == nil || sender == nil || store == nil || sender.TenantID != approvals.tenant || digest(sender.Profile) != approvals.profile {
		return nil, ErrBridge
	}
	cp := *sender
	cp.Profile = append(json.RawMessage(nil), sender.Profile...)
	cp.Approvals = approvals
	return &ReplyModule{approvals: approvals, sender: &cp, store: store, channel: &outbounddelivery.Channel{CodeValue: "whatsapp", Receiver: VerifiedInboxChannel{}, Sender: &cp, Store: store}}, nil
}
func (m *ReplyModule) principal(w http.ResponseWriter, r *http.Request, v identity.Verifier, permission string) (identity.Principal, bool) {
	values := r.Header.Values("Authorization")
	if v == nil || len(values) != 1 || !strings.HasPrefix(values[0], "Bearer ") || len(values[0]) > 16391 || len(strings.Fields(values[0])) != 2 {
		notificationProblem(w, 401, "UNAUTHENTICATED")
		return identity.Principal{}, false
	}
	p, err := v.Verify(r.Context(), strings.TrimPrefix(values[0], "Bearer "))
	if err != nil {
		notificationProblem(w, 401, "UNAUTHENTICATED")
		return identity.Principal{}, false
	}
	if m == nil || !m.approvals.validPrincipal(p, permission) {
		notificationProblem(w, 403, "FORBIDDEN")
		return identity.Principal{}, false
	}
	return p, true
}
func replyJSON(w http.ResponseWriter, value any) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(value)
}
func readReplyCommand(w http.ResponseWriter, r *http.Request, decision bool) (string, bool, string, error) {
	if r.Header.Get("Content-Type") != "application/json" || r.URL.RawQuery != "" {
		return "", false, "", ErrApproval
	}
	raw, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 4096))
	if err != nil {
		return "", false, "", ErrApproval
	}
	canonical, _, err := approval.CanonicalPayload(raw)
	if err != nil {
		return "", false, "", ErrApproval
	}
	var fields map[string]json.RawMessage
	if json.Unmarshal(canonical, &fields) != nil {
		return "", false, "", ErrApproval
	}
	expected := 1
	if decision {
		expected = 3
	}
	if len(fields) != expected {
		return "", false, "", ErrApproval
	}
	for _, raw := range fields {
		if bytes.Equal(bytes.TrimSpace(raw), []byte("null")) {
			return "", false, "", ErrApproval
		}
	}
	var hash, reason string
	var approved bool
	if json.Unmarshal(fields["payload_sha256"], &hash) != nil || !validDigest(hash) {
		return "", false, "", ErrApproval
	}
	if decision {
		if json.Unmarshal(fields["approved"], &approved) != nil || json.Unmarshal(fields["reason"], &reason) != nil || len(reason) > 2048 {
			return "", false, "", ErrApproval
		}
	}
	return hash, approved, reason, nil
}
func (m *ReplyModule) Register(mux *http.ServeMux, v identity.Verifier) {
	m.registerAIActivity(mux, v)
	mux.HandleFunc("GET /v1/franchise/whatsapp/replies", func(w http.ResponseWriter, r *http.Request) {
		p, ok := m.principal(w, r, v, "whatsapp:approve")
		if !ok {
			return
		}
		if r.URL.RawQuery != "" {
			notificationProblem(w, 400, "INVALID_QUERY")
			return
		}
		rows, err := m.approvals.base.pool.Query(r.Context(), `select request_id,state,evidence_sha,payload from approval.request where tenant_id=$1 and organization_id=$2 and kind='whatsapp_reply' order by created_at desc,request_id limit 50`, p.TenantID, m.approvals.organization)
		if err != nil {
			notificationProblem(w, 503, "UNAVAILABLE")
			return
		}
		defer rows.Close()
		values := []ReplyProposal{}
		for rows.Next() {
			var value ReplyProposal
			var raw []byte
			if rows.Scan(&value.RequestID, &value.State, &value.PayloadSHA256, &raw) != nil || json.Unmarshal(raw, &value.Context) != nil {
				notificationProblem(w, 503, "UNAVAILABLE")
				return
			}
			values = append(values, value)
		}
		if rows.Err() != nil {
			notificationProblem(w, 503, "UNAVAILABLE")
			return
		}
		rows.Close()
		for i := range values {
			value, err := m.approvals.Read(r.Context(), p, values[i].RequestID)
			if err != nil {
				notificationProblem(w, 503, "UNAVAILABLE")
				return
			}
			values[i] = value
		}
		replyJSON(w, values)
	})
	mux.HandleFunc("GET /v1/franchise/whatsapp/replies/{id}", func(w http.ResponseWriter, r *http.Request) {
		p, ok := m.principal(w, r, v, "whatsapp:approve")
		if !ok {
			return
		}
		value, err := m.approvals.Read(r.Context(), p, r.PathValue("id"))
		if err != nil {
			notificationProblem(w, 404, "NOT_FOUND")
			return
		}
		replyJSON(w, value)
	})
	mux.HandleFunc("POST /v1/franchise/whatsapp/replies/{id}/decision", func(w http.ResponseWriter, r *http.Request) {
		p, ok := m.principal(w, r, v, "whatsapp:approve")
		if !ok {
			return
		}
		hash, approved, reason, err := readReplyCommand(w, r, true)
		if err != nil {
			notificationProblem(w, 400, "INVALID_BODY")
			return
		}
		value, err := m.approvals.Decide(r.Context(), p, r.PathValue("id"), hash, approved, reason)
		if err != nil {
			notificationProblem(w, 409, "APPROVAL_NOT_CURRENT_OR_DIVERGENT")
			return
		}
		replyJSON(w, value)
	})
	mux.HandleFunc("POST /v1/franchise/whatsapp/replies/{id}/recover", func(w http.ResponseWriter, r *http.Request) {
		p, ok := m.principal(w, r, v, "whatsapp:send")
		if !ok {
			return
		}
		hash, _, _, err := readReplyCommand(w, r, false)
		if err != nil {
			notificationProblem(w, 400, "INVALID_BODY")
			return
		}
		_, err = m.Recover(r.Context(), p, r.PathValue("id"), hash)
		if err != nil {
			notificationProblem(w, 409, "DELIVERY_RECONCILIATION_REQUIRED")
			return
		}
		replyJSON(w, map[string]string{"status": "accepted", "delivery_key": r.PathValue("id")})
	})

	mux.HandleFunc("POST /v1/franchise/whatsapp/replies/{id}/send", func(w http.ResponseWriter, r *http.Request) {
		p, ok := m.principal(w, r, v, "whatsapp:send")
		if !ok {
			return
		}
		// Read requires review permission too: a sender must be able to inspect the
		// exact proposal being released to the provider.
		hash, _, _, err := readReplyCommand(w, r, false)
		if err != nil {
			notificationProblem(w, 400, "INVALID_BODY")
			return
		}
		value, err := m.approvals.Read(r.Context(), p, r.PathValue("id"))
		if err != nil || value.State != "approved" || value.PayloadSHA256 != hash {
			notificationProblem(w, 409, "APPROVAL_NOT_CURRENT_OR_DIVERGENT")
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 45*time.Second)
		defer cancel()
		if _, err = m.approvals.ResolveWhatsAppApproval(ctx, p.TenantID, value.RequestID); err != nil {
			notificationProblem(w, 409, "APPROVAL_NOT_CURRENT_OR_DIVERGENT")
			return
		}
		if err = m.channel.Send(ctx, value.Context.Message); err != nil {
			notificationProblem(w, 409, "DELIVERY_RECONCILIATION_REQUIRED")
			return
		}
		replyJSON(w, map[string]string{"status": "accepted", "delivery_key": value.RequestID})
	})
}
