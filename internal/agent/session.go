// Package agent provides a deterministic, framework-agnostic conversational
// agent core: a bounded intent router, per-intent tool authorization, guarded
// knowledge retrieval, fail-closed human handoff and append-only turn
// evidence. It consumes aifoundation (safety) and search (knowledge); the LLM
// provider and channel adapters are CONDITIONED.
package agent

import "time"

// SessionState is the lifecycle state of a conversation.
type SessionState string

const (
	StateOpen      SessionState = "open"
	StateHandedOff SessionState = "handed_off"
	StateClosed    SessionState = "closed"
)

// Role is a conversation participant.
type Role string

const (
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleTool      Role = "tool"
)

// Message is a single turn of a conversation.
type Message struct {
	Role Role
	Text string
}

// Turn records one user message, the resolved intent, any tool calls and the
// assistant reply. Turns are append-only evidence.
type Turn struct {
	UserText      string
	Intent        Intent
	AssistantText string
	ToolCalls     []ToolCall
	HandedOff     bool
	At            time.Time
}

// Session is an in-memory conversation. It is append-only by convention: a
// durable store (cache/outbox) is the project's persistence integration.
type Session struct {
	ID       string
	TenantID string
	State    SessionState
	Messages []Message
	Turns    []Turn
}

// NewSession returns an open session scoped to one tenant.
func NewSession(id, tenantID string) *Session {
	return &Session{ID: id, TenantID: tenantID, State: StateOpen}
}
