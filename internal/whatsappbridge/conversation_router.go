package whatsappbridge

// AUTHORED typed projection after the existing Meta signature/scope verifier.
// The existing conversation runtime stores a proposal; this path never sends it.
import (
	"context"
	"encoding/json"
	"strconv"
	"time"
	"unicode/utf8"

	"elite.local/enterprise/internal/channels"
	"elite.local/enterprise/internal/conversationruntime"
	"elite.local/enterprise/internal/platform/identity"
)

type ConversationHandler interface {
	Handle(context.Context, channels.Message) (string, error)
}
type ReplyProposer interface {
	Propose(context.Context, identity.Principal, string, channels.Message) (ReplyProposal, error)
}

type ConversationRoute struct {
	OrganizationID string
	Runtime        ConversationHandler
	Proposals      ReplyProposer
}

// EnableConversation is explicit composition; its dependencies cannot be
// supplied by a webhook. The organization must match the durable connection.
func (r *StatusRouter) EnableConversation(route ConversationRoute) error {
	if r == nil || route.OrganizationID == "" || route.Runtime == nil || route.Proposals == nil {
		return ErrStatusRouting
	}
	r.conversation = &route
	return nil
}

// ScopedContactResolver prevents a tenant-wide contact binding from crossing
// the connection's selected organization before tools or model execution.
type ScopedContactResolver struct {
	TenantID, OrganizationID string
	Resolver                 conversationruntime.ContactResolver
}

func (s ScopedContactResolver) Resolve(ctx context.Context, m channels.Message) (conversationruntime.Contact, error) {
	if s.Resolver == nil || s.TenantID == "" || s.OrganizationID == "" || m.TenantID != s.TenantID || m.ChannelCode != "whatsapp" {
		return conversationruntime.Contact{}, ErrStatusRouting
	}
	c, err := s.Resolver.Resolve(ctx, m)
	if err != nil || c.Scope.OrganizationID != s.OrganizationID {
		return conversationruntime.Contact{}, ErrStatusRouting
	}
	return c, nil
}

func conversationThread(connection string) string {
	return "wa-connection:" + digest([]byte(connection))
}

// Raw has already passed Python's exact-key, HMAC, scope and event validation.
// This projection still validates supported message shape before any runtime
// effect. Unsupported media remains in the inbox for explicit handoff/review.
func verifiedConversationMessages(raw []byte, tenant, connection string, limit int) ([]channels.Message, int, error) {
	var envelope map[string]json.RawMessage
	if json.Unmarshal(raw, &envelope) != nil {
		return nil, 0, ErrStatusRouting
	}
	var entries []map[string]json.RawMessage
	if json.Unmarshal(envelope["entry"], &entries) != nil {
		return nil, 0, ErrStatusRouting
	}
	var out []channels.Message
	seen := map[string]string{}
	total := 0
	for _, entry := range entries {
		var changes []map[string]json.RawMessage
		if json.Unmarshal(entry["changes"], &changes) != nil {
			return nil, 0, ErrStatusRouting
		}
		for _, change := range changes {
			var value map[string]json.RawMessage
			if json.Unmarshal(change["value"], &value) != nil {
				return nil, 0, ErrStatusRouting
			}
			var messages []map[string]json.RawMessage
			if v, ok := value["messages"]; ok && json.Unmarshal(v, &messages) != nil {
				return nil, 0, ErrStatusRouting
			}
			for _, m := range messages {
				total++
				var id, from, stamp, kind string
				if json.Unmarshal(m["id"], &id) != nil || json.Unmarshal(m["from"], &from) != nil || json.Unmarshal(m["timestamp"], &stamp) != nil || json.Unmarshal(m["type"], &kind) != nil || kind != "text" {
					return nil, 0, ErrStatusRouting
				}
				var body map[string]json.RawMessage
				var text string
				if json.Unmarshal(m["text"], &body) != nil || json.Unmarshal(body["body"], &text) != nil || !utf8.ValidString(text) || len(text) == 0 || len(text) > 16384 {
					return nil, 0, ErrStatusRouting
				}
				unix, err := strconv.ParseInt(stamp, 10, 64)
				if err != nil || unix < 1 || strconv.FormatInt(unix, 10) != stamp || time.Unix(unix, 0).After(time.Now().Add(time.Minute)) {
					return nil, 0, ErrStatusRouting
				}
				message := channels.Message{ChannelCode: "whatsapp", TenantID: tenant, ExternalID: from, ThreadID: conversationThread(connection), ProviderMessageID: id, OccurredAt: time.Unix(unix, 0).UTC(), Direction: channels.DirectionIn, Text: text}
				if message.Validate() != nil {
					return nil, 0, ErrStatusRouting
				}
				encoded, _ := json.Marshal(message)
				hash := digest(encoded)
				if previous, ok := seen[id]; ok {
					if previous != hash {
						return nil, 0, ErrStatusRouting
					}
					continue
				}
				seen[id] = hash
				out = append(out, message)
				if len(out) > limit {
					return nil, 0, ErrStatusRouting
				}
			}
		}
	}
	return out, total, nil
}
