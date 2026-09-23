package agent

import (
	"context"
	"testing"
)

type stubTool struct {
	name   string
	intent Intent
}

func (s stubTool) Name() string          { return s.name }
func (s stubTool) Intent() Intent        { return s.intent }
func (s stubTool) Run(context.Context, string, string) (string, error) { return "ok", nil }

func TestToolRegistryAuthorizesOnePerIntent(t *testing.T) {
	r := NewToolRegistry()
	if err := r.Register(stubTool{name: "appointments", intent: IntentAppointment}); err != nil {
		t.Fatal(err)
	}
	if err := r.Register(stubTool{name: "appointments2", intent: IntentAppointment}); err == nil {
		t.Fatal("duplicate intent accepted")
	}
	if _, ok := r.For(IntentAppointment); !ok {
		t.Fatal("expected tool for appointment intent")
	}
	if _, ok := r.For(IntentQuote); ok {
		t.Fatal("unexpected tool for quote intent")
	}
}

func TestToolRegistryRejectsNilAndEmptyName(t *testing.T) {
	r := NewToolRegistry()
	if err := r.Register(nil); err == nil {
		t.Fatal("nil tool accepted")
	}
	if err := r.Register(stubTool{name: "", intent: IntentQuote}); err == nil {
		t.Fatal("empty name accepted")
	}
}
