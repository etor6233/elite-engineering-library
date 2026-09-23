package agent

import (
	"context"
	"errors"
	"strings"
	"time"

	"elite.local/enterprise/internal/aifoundation"
)

var (
	ErrNilRouter = errors.New("agent: nil router")
	ErrNilSafety = errors.New("agent: nil safety gate")
)

// Agent is the conversational core. It is deterministic and testable without a
// live model; an LLM provider and channel adapters are CONDITIONED.
type Agent struct {
	Router    Router
	Tools     *ToolRegistry
	Knowledge Knowledge
	Safety    aifoundation.SafetyGate
	Greeting  string
	HandoffText  string
}

// NewAgent returns an agent with sensible defaults and nil-checked guards.
func NewAgent(r Router, tools *ToolRegistry, safety aifoundation.SafetyGate) *Agent {
	return &Agent{
		Router:     r,
		Tools:      tools,
		Safety:     safety,
		Greeting:   "¡Hola! Soy tu asistente. ¿En qué puedo ayudarte hoy?",
		HandoffText: "Voy a derivar tu consulta a un asesor humano, que te contactará a la brevedad.",
	}
}

// Process handles one user message: safety pre-flight, routing, tool/knowledge
// execution, safety post-flight and evidence finalization. It fails closed into
// human handoff on any unsafe or unknown path.
func (a *Agent) Process(ctx context.Context, s *Session, userText string) (Turn, error) {
	if a.Router == nil {
		return Turn{}, ErrNilRouter
	}
	if a.Safety == nil {
		return Turn{}, ErrNilSafety
	}
	turn := Turn{UserText: userText, At: time.Now().UTC()}
	s.Messages = append(s.Messages, Message{Role: RoleUser, Text: userText})

	if a.Safety.PreFlight(aifoundation.SafetyInput{UserText: userText}) != aifoundation.SafetyAllow {
		a.handoff(s, &turn, IntentFallback, "pre-flight safety block")
		return a.finalize(s, turn), nil
	}

	intent := a.Router.Route(userText)
	turn.Intent = intent

	switch intent {
	case IntentGreeting:
		turn.AssistantText = a.Greeting
	case IntentHandoff:
		a.handoff(s, &turn, intent, "user requested human")
	case IntentFAQ:
		a.answerFAQ(ctx, s, userText, &turn)
	case IntentAppointment, IntentQuote, IntentOrderStatus, IntentReturnRequest:
		a.runTool(ctx, s, intent, userText, &turn)
	default:
		a.handoff(s, &turn, IntentFallback, "unroutable input")
	}

	if a.Safety.PostFlight(aifoundation.SafetyOutput{Content: turn.AssistantText}) != aifoundation.SafetyAllow {
		a.handoff(s, &turn, turn.Intent, "post-flight safety block")
	}

	return a.finalize(s, turn), nil
}

func (a *Agent) answerFAQ(ctx context.Context, s *Session, text string, turn *Turn) {
	if a.Knowledge == nil {
		a.handoff(s, turn, IntentFAQ, "knowledge unavailable")
		return
	}
	docs, err := a.Knowledge.Retrieve(ctx, s.TenantID, text, 3)
	if err != nil || len(docs) == 0 {
		a.handoff(s, turn, IntentFAQ, "no knowledge or retrieval error")
		return
	}
	turn.AssistantText = strings.Join(docs, "\n")
}

func (a *Agent) runTool(ctx context.Context, s *Session, intent Intent, text string, turn *Turn) {
	if a.Tools == nil {
		a.handoff(s, turn, intent, "tool registry unavailable")
		return
	}
	t, ok := a.Tools.For(intent)
	if !ok {
		a.handoff(s, turn, intent, "no tool authorized for intent")
		return
	}
	call := ToolCall{Intent: intent, Name: t.Name(), Args: text}
	result, err := t.Run(ctx, s.TenantID, text)
	switch {
	case errors.Is(err, ErrNeedsInfo):
		call.Result = result
		turn.ToolCalls = append(turn.ToolCalls, call)
		turn.AssistantText = result // prompt for more info; session stays open
	case err != nil:
		call.Err = err.Error()
		turn.ToolCalls = append(turn.ToolCalls, call)
		a.handoff(s, turn, intent, "tool failure: "+err.Error())
	default:
		call.Result = result
		turn.ToolCalls = append(turn.ToolCalls, call)
		turn.AssistantText = result
	}
}

func (a *Agent) handoff(s *Session, turn *Turn, intent Intent, reason string) {
	turn.HandedOff = true
	turn.Intent = intent
	turn.AssistantText = a.HandoffText
	s.State = StateHandedOff
}

func (a *Agent) finalize(s *Session, turn Turn) Turn {
	s.Messages = append(s.Messages, Message{Role: RoleAssistant, Text: turn.AssistantText})
	s.Turns = append(s.Turns, turn)
	return turn
}
