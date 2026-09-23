package whatsappbridge

// AUTHORED read projection. Existing runtime/approval/delivery owners retain all
// writes. The historical approval organization, not a mutable connection or a
// tenant-wide turn list, authorizes this view. No message text/contact is emitted.
import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"elite.local/enterprise/internal/approval"
	"elite.local/enterprise/internal/platform/identity"
	"github.com/jackc/pgx/v5"
)

type AIActivityItem struct {
	RequestID      string    `json:"request_id"`
	State          string    `json:"state"`
	ReviewState    string    `json:"review_state"`
	Tool           string    `json:"tool"`
	Attempts       int       `json:"attempts"`
	ReservedTokens string    `json:"reserved_tokens"`
	InputTokens    string    `json:"input_tokens"`
	OutputTokens   string    `json:"output_tokens"`
	IntentPinned   bool      `json:"intent_pinned"`
	OccurredAt     time.Time `json:"occurred_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
type AIActivity struct {
	OrganizationID string           `json:"organization_id"`
	ObservedAt     time.Time        `json:"observed_at"`
	Scope          string           `json:"scope"`
	HasMore        bool             `json:"has_more"`
	Items          []AIActivityItem `json:"items"`
}

func validAIActivity(v AIActivityItem) bool {
	if len(v.RequestID) != 73 || v.RequestID[:9] != "wa-reply:" || !validDigest(v.RequestID[9:]) || (v.State != "completed" && v.State != "handed_off") || (v.ReviewState != "pending" && v.ReviewState != "approved" && v.ReviewState != "rejected") || v.Attempts < 1 || v.Attempts > 100 || v.OccurredAt.IsZero() || v.UpdatedAt.Before(v.OccurredAt) {
		return false
	}
	switch v.Tool {
	case "none", "book_appointment", "create_quote", "order_status":
	default:
		return false
	}
	for _, n := range []string{v.ReservedTokens, v.InputTokens, v.OutputTokens} {
		v, e := strconv.ParseInt(n, 10, 64)
		if e != nil || v < 0 || strconv.FormatInt(v, 10) != n {
			return false
		}
	}
	return true
}

// AIActivity reads at most 50 latest retained, exactly scoped reviewable turns.
// Unbound handoffs (no approval request) and in-flight turns are intentionally
// excluded: conversation_turn has no historical organization authority.
// Reservations are conservative internal token units, never provider invoices.
func (s *PostgresReplyApprovals) AIActivity(ctx context.Context, p identity.Principal) (AIActivity, error) {
	out := AIActivity{Scope: "reviewable_proposals", Items: []AIActivityItem{}}
	if !s.validPrincipal(p, "whatsapp:approve") || s.base == nil || s.base.pool == nil {
		return AIActivity{}, ErrApproval
	}
	tx, e := s.base.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead, AccessMode: pgx.ReadOnly})
	if e != nil {
		return AIActivity{}, ErrNotificationRead
	}
	defer tx.Rollback(ctx)
	out.OrganizationID = s.organization
	if e = tx.QueryRow(ctx, `select statement_timestamp()`).Scan(&out.ObservedAt); e != nil {
		return AIActivity{}, ErrNotificationRead
	}
	rows, e := tx.Query(ctx, `select r.request_id,r.state,r.evidence_sha,r.payload from approval.request r
 where r.tenant_id=$1 and r.organization_id=$2 and r.kind='whatsapp_reply'
 and r.payload->>'connection_id'=$3 and r.payload->>'profile_sha256'=$4
 order by r.created_at desc,r.request_id limit 51`, p.TenantID, s.organization, s.connection, s.profile)
	if e != nil {
		return AIActivity{}, ErrNotificationRead
	}
	type proposal struct {
		id, state, hash string
		raw             []byte
	}
	pending := []proposal{}
	for rows.Next() {
		var x proposal
		if e = rows.Scan(&x.id, &x.state, &x.hash, &x.raw); e != nil {
			rows.Close()
			return AIActivity{}, ErrNotificationRead
		}
		pending = append(pending, x)
	}
	e = rows.Err()
	rows.Close()
	if e != nil {
		return AIActivity{}, ErrNotificationRead
	}
	if len(pending) > 50 {
		out.HasMore = true
		pending = pending[:50]
	}
	for _, x := range pending {
		var c ReplyContext
		_, h, he := approval.CanonicalPayload(x.raw)
		if he != nil || h != x.hash || json.Unmarshal(x.raw, &c) != nil || c.Schema != "elite-whatsapp-reply-approval/v1" || c.OrganizationID != s.organization || c.ConnectionID != s.connection || c.ProfileSHA256 != s.profile || c.Message.TenantID != p.TenantID || c.Message.ChannelCode != "whatsapp" || c.Message.ThreadID != conversationThread(s.connection) || c.Message.DeliveryKey != x.id || x.id != ReplyDeliveryKey(p.TenantID, s.connection, c.ProviderMessageID) {
			return AIActivity{}, ErrNotificationIntegrity
		}
		raw, e := replyRequestFor(c.Message)
		if e != nil {
			return AIActivity{}, ErrNotificationIntegrity
		}
		var reply ReplyRequest
		if json.Unmarshal(raw, &reply) != nil || reply.SourceMessageID != c.ProviderMessageID {
			return AIActivity{}, ErrNotificationIntegrity
		}
		v := AIActivityItem{RequestID: x.id, ReviewState: x.state}
		var responseHash string
		var expired bool
		e = tx.QueryRow(ctx, `select t.state,coalesce(t.tool_name,'none'),t.attempt_count,
 coalesce((select sum(tokens)::text from communication.conversation_reservation b where b.tenant_id=t.tenant_id and b.channel_code=t.channel_code and b.provider_message_id=t.provider_message_id),'0'),
 t.input_tokens::text,t.output_tokens::text,
 exists(select 1 from communication.conversation_intent i where i.tenant_id=t.tenant_id and i.channel_code=t.channel_code and i.provider_message_id=t.provider_message_id),
 t.occurred_at,t.updated_at,t.response_sha256_hex,t.expires_at<=statement_timestamp()
 from communication.conversation_turn t where t.tenant_id=$1 and t.channel_code='whatsapp' and t.provider_message_id=$2 and t.external_id=$3 and t.thread_id=$4
 and t.occurred_at=to_timestamp($5) and t.assistant_text=$6 and t.state in ('completed','handed_off')`, p.TenantID, c.ProviderMessageID, c.Message.ExternalID, c.Message.ThreadID, reply.LastInboundAt, reply.Text).Scan(&v.State, &v.Tool, &v.Attempts, &v.ReservedTokens, &v.InputTokens, &v.OutputTokens, &v.IntentPinned, &v.OccurredAt, &v.UpdatedAt, &responseHash, &expired)
		// Expired turn is intentionally outside retention; do not make a new copy.
		if e == pgx.ErrNoRows {
			return AIActivity{}, ErrNotificationIntegrity
		}
		if e != nil {
			return AIActivity{}, ErrNotificationRead
		}
		if responseHash != digest([]byte(reply.Text)) || !validAIActivity(v) {
			return AIActivity{}, ErrNotificationIntegrity
		}
		if !expired {
			out.Items = append(out.Items, v)
		}
	}
	if e = tx.Commit(ctx); e != nil {
		return AIActivity{}, ErrNotificationRead
	}
	return out, nil
}
