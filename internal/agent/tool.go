package agent

import (
	"context"
	"errors"
	"fmt"
)

// ErrNeedsInfo signals a tool that needs more information from the user.
var ErrNeedsInfo = errors.New("agent: tool needs more information")

// ToolCall records an authorized tool invocation and its outcome.
type ToolCall struct {
	Intent Intent
	Name   string
	Args   string
	Result string
	Err    string
}

// Tool is a single, authorized capability bound to exactly one intent. It
// receives the tenant and the raw user text and returns a confirmation or an
// error. It never performs a side effect outside its declared intent.
type Tool interface {
	Name() string
	Intent() Intent
	Run(ctx context.Context, tenantID string, text string) (string, error)
}

// ToolRegistry authorizes at most one tool per intent.
type ToolRegistry struct {
	tools map[Intent]Tool
}

// NewToolRegistry returns an empty registry.
func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{tools: make(map[Intent]Tool)}
}

// Register binds a tool to its intent, rejecting duplicates.
func (r *ToolRegistry) Register(t Tool) error {
	if t == nil {
		return errors.New("agent: nil tool")
	}
	if t.Name() == "" {
		return errors.New("agent: tool with empty name")
	}
	if _, ok := r.tools[t.Intent()]; ok {
		return fmt.Errorf("agent: duplicate tool for intent %q", t.Intent())
	}
	r.tools[t.Intent()] = t
	return nil
}

// For returns the tool authorized for an intent, if any.
func (r *ToolRegistry) For(intent Intent) (Tool, bool) {
	t, ok := r.tools[intent]
	return t, ok
}
