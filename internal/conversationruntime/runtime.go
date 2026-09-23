// Package conversationruntime connects a validated channel envelope to
// budgeted Responses tool calling, authorized contact scope, domain commands,
// approval, durable turn evidence and a replayable reply.
package conversationruntime

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"elite.local/enterprise/internal/agenttools"
	"elite.local/enterprise/internal/aifoundation"
	"elite.local/enterprise/internal/approval"
	"elite.local/enterprise/internal/channels"
	"elite.local/enterprise/internal/domainbind"

	"elite.local/enterprise/internal/llmopenai"
)

var (
	ErrBudgetExceeded      = errors.New("conversationruntime: durable budget exceeded")
	ErrBudgetConfiguration = errors.New("conversationruntime: changed durable budget cap")
	ErrToolIntentConflict  = errors.New("conversationruntime: changed tool intent")
	ErrInvalidRuntime      = errors.New("conversationruntime: invalid runtime")
	ErrReplayConflict      = errors.New("conversationruntime: replay payload conflict")
	ErrInProgress          = errors.New("conversationruntime: message is already processing")
	ErrTerminal            = errors.New("conversationruntime: message reached terminal failure")
)

type TurnState string

const (
	StateCompleted TurnState = "completed"
	StateHandedOff TurnState = "handed_off"
)

type Claim struct {
	Generation   int
	Replay       bool
	State        TurnState
	ResponseText string
}

type HistoryMessage struct {
	UserText      string
	AssistantText string
}

type Completion struct {
	State          TurnState
	ResponseText   string
	LLMResponseID  string
	ToolName       string
	ToolArguments  json.RawMessage
	ToolResult     string
	InputTokens    int64
	OutputTokens   int64
	ReservedTokens int64
	FailureCode    string
}

type Store interface {
	Reserve(context.Context, channels.Message, string, int, int64, int64) error
	PinTool(context.Context, channels.Message, string, int, string) error
	Claim(context.Context, channels.Message, string) (Claim, error)
	History(context.Context, channels.Message, int) ([]HistoryMessage, error)
	Complete(context.Context, channels.Message, string, int, Completion) error
	FailRetryable(context.Context, channels.Message, string, int, string) error
}

type Model interface {
	StartToolTurn(context.Context, llmopenai.ToolTurnRequest) (llmopenai.ToolTurnResponse, error)
	CompleteToolTurn(context.Context, string, string, []llmopenai.FunctionOutput, int) (llmopenai.ToolTurnResponse, error)
}

type Contact struct {
	Scope      domainbind.Scope
	SubjectID  string
	PIIAllowed bool
}

type ContactResolver interface {
	Resolve(context.Context, channels.Message) (Contact, error)
}

type Domain interface {
	BookAppointmentFor(context.Context, string, domainbind.Scope, domainbind.CommandIdentity, agenttools.AppointmentInput) (string, error)
	CreateQuoteFor(context.Context, string, domainbind.Scope, domainbind.CommandIdentity, agenttools.QuoteInput) (string, error)
	OrderStatusFor(context.Context, string, domainbind.Scope, string) (string, error)
}

type ApprovalGate interface {
	Submit(approval.Request) (approval.State, error)
	State(string, string) (approval.State, bool)
}

type Config struct {
	TokenBudget     int64
	ModelRevision   string
	Instructions    string
	PromptCacheKey  string
	MaxOutputTokens int
	MaxHistory      int
	MaxInputBytes   int
	MaxToolBytes    int
	StoreApproved   bool
	HandoffText     string
}

type Runtime struct {
	Config   Config
	Store    Store
	Model    Model
	Resolver ContactResolver
	Domain   Domain
	Approval ApprovalGate
	Safety   aifoundation.SafetyGate
}

func (r *Runtime) Validate() error {
	if r.Store == nil || r.Model == nil || r.Resolver == nil || r.Domain == nil || r.Approval == nil || r.Safety == nil {
		return ErrInvalidRuntime
	}
	if r.Config.TokenBudget < 1 || r.Config.TokenBudget > 1000000000000 || strings.TrimSpace(r.Config.ModelRevision) == "" || len(r.Config.Instructions) > 65536 || strings.TrimSpace(r.Config.Instructions) == "" || strings.TrimSpace(r.Config.PromptCacheKey) == "" || r.Config.MaxOutputTokens < 1 || r.Config.MaxOutputTokens > 4096 || r.Config.MaxHistory < 0 || r.Config.MaxHistory > 20 || r.Config.MaxInputBytes < 1 || r.Config.MaxInputBytes > 65536 || r.Config.MaxToolBytes < 1 || r.Config.MaxToolBytes > 65536 || !r.Config.StoreApproved || strings.TrimSpace(r.Config.HandoffText) == "" {
		return ErrInvalidRuntime
	}
	return nil
}

func (r *Runtime) Handle(ctx context.Context, message channels.Message) (string, error) {
	if err := r.Validate(); err != nil {
		return "", err
	}
	if err := message.Validate(); err != nil || message.Direction != channels.DirectionIn {
		return "", channels.ErrInvalidMessage
	}
	if len(message.Text) > r.Config.MaxInputBytes {
		return "", channels.ErrInvalidMessage
	}
	requestHash, err := hashMessage(message)
	if err != nil {
		return "", err
	}
	claim, err := r.Store.Claim(ctx, message, requestHash)
	if err != nil {
		return "", err
	}
	if claim.Replay {
		return claim.ResponseText, nil
	}
	if claim.Generation < 1 {
		return "", ErrInvalidRuntime
	}
	contact, err := r.Resolver.Resolve(ctx, message)
	if err != nil || strings.TrimSpace(contact.SubjectID) == "" || contact.Scope.OrganizationID == "" || contact.Scope.LeadID == "" {
		return r.handoff(ctx, message, requestHash, claim.Generation, "CONTACT_SCOPE_UNAVAILABLE", Completion{})
	}
	if r.Safety.PreFlight(aifoundation.SafetyInput{UserText: message.Text, ContainsPII: true, Allowlisted: contact.PIIAllowed}) != aifoundation.SafetyAllow {
		return r.handoff(ctx, message, requestHash, claim.Generation, "SAFETY_PREFLIGHT", Completion{})
	}
	history, err := r.Store.History(ctx, message, r.Config.MaxHistory)
	if err != nil {
		return "", r.failRetryable(ctx, message, requestHash, claim.Generation, "HISTORY_FAILED", err)
	}
	input := buildInput(history, message.Text)
	if len(input) > r.Config.MaxInputBytes {
		return r.handoff(ctx, message, requestHash, claim.Generation, "CONTEXT_TOO_LARGE", Completion{})
	}
	// Byte-based conservative capacity reservation for two bounded requests.
	// This is not a tokenizer or a monetary cost guarantee. Actual provider usage
	// remains evidence; each uncertain retry consumes another reservation.
	schemaBytes := 0
	for _, tool := range runtimeTools() {
		schemaBytes += len(tool.Name) + len(tool.Description) + len(tool.Parameters)
	}
	reserved := int64(2*(len(r.Config.Instructions)+len(input)+schemaBytes+2048) + 3*r.Config.MaxOutputTokens + r.Config.MaxToolBytes)
	if err := r.Store.Reserve(ctx, message, requestHash, claim.Generation, r.Config.TokenBudget, reserved); err != nil {
		if errors.Is(err, ErrBudgetExceeded) {
			return r.handoff(ctx, message, requestHash, claim.Generation, "BUDGET_BLOCKED", Completion{})
		}
		return "", r.failRetryable(ctx, message, requestHash, claim.Generation, "BUDGET_RESERVATION_FAILED", err)
	}
	first, err := r.Model.StartToolTurn(ctx, llmopenai.ToolTurnRequest{
		Instructions: r.Config.Instructions, Input: input, Messages: buildMessages(history, message.Text), Tools: runtimeTools(),
		MaxOutputTokens: r.Config.MaxOutputTokens, PromptCacheKey: r.Config.PromptCacheKey,
		StoreApproved: r.Config.StoreApproved,
	})
	if err != nil {
		// Invalid arguments from a completed response are not a transport retry.
		// Preserve the durable, non-executing human handoff contract.
		if errors.Is(err, llmopenai.ErrInvalidToolArguments) {
			return r.handoff(ctx, message, requestHash, claim.Generation, "INVALID_TOOL_ARGUMENTS", Completion{})
		}
		return "", r.failRetryable(ctx, message, requestHash, claim.Generation, "MODEL_START_FAILED", err)
	}
	base := Completion{LLMResponseID: first.ResponseID, InputTokens: first.InputTokens, OutputTokens: first.OutputTokens, ReservedTokens: reserved}
	if len(first.FunctionCalls) == 0 {
		return r.completeText(ctx, message, requestHash, claim.Generation, first.OutputText, base)
	}
	if len(first.FunctionCalls) != 1 {
		return r.handoff(ctx, message, requestHash, claim.Generation, "MULTIPLE_TOOL_CALLS", base)
	}
	call := first.FunctionCalls[0]
	if len(call.Arguments) > r.Config.MaxToolBytes || call.CallID == "" {
		return r.handoff(ctx, message, requestHash, claim.Generation, "INVALID_TOOL_ARGUMENTS", base)
	}
	var schema json.RawMessage
	for _, tool := range runtimeTools() {
		if tool.Name == call.Name {
			schema = tool.Parameters
		}
	}
	if schema == nil {
		return r.handoff(ctx, message, requestHash, claim.Generation, "UNAUTHORIZED_TOOL", base)
	}
	if aifoundation.ValidateStructuredOutput(schema, call.Arguments) != nil {
		return r.handoff(ctx, message, requestHash, claim.Generation, "INVALID_TOOL_ARGUMENTS", base)
	}
	canonical, _, err := approval.CanonicalPayload(call.Arguments)
	if err != nil {
		return r.handoff(ctx, message, requestHash, claim.Generation, "INVALID_TOOL_ARGUMENTS", base)
	}
	call.Arguments = canonical
	intentRaw, _ := json.Marshal(struct {
		Contact   Contact
		Tool      string
		Arguments json.RawMessage
		Policy    Config
	}{contact, call.Name, canonical, r.Config})
	intentSum := sha256.Sum256(intentRaw)
	intent := hex.EncodeToString(intentSum[:])
	if err := r.Store.PinTool(ctx, message, requestHash, claim.Generation, intent); err != nil {
		if errors.Is(err, ErrToolIntentConflict) {
			return r.handoff(ctx, message, requestHash, claim.Generation, "TOOL_INTENT_CHANGED", base)
		}
		return "", r.failRetryable(ctx, message, requestHash, claim.Generation, "TOOL_FENCE_FAILED", err)
	}
	base.ToolName, base.ToolArguments = call.Name, append(json.RawMessage(nil), call.Arguments...)
	currentContact, contactErr := r.Resolver.Resolve(ctx, message)
	if contactErr != nil || currentContact != contact {
		return r.handoff(ctx, message, requestHash, claim.Generation, "CONTACT_CHANGED", base)
	}
	toolResult, pending, terminal, err := r.execute(ctx, message, contact, intent, call)
	if err != nil {
		return "", r.failRetryable(ctx, message, requestHash, claim.Generation, "DOMAIN_CALL_FAILED", err)
	}
	if terminal != "" {
		return r.handoff(ctx, message, requestHash, claim.Generation, terminal, base)
	}
	if pending {
		base.ToolResult = toolResult
		return r.handoff(ctx, message, requestHash, claim.Generation, "HUMAN_APPROVAL_REQUIRED", base)
	}
	if len(toolResult) > r.Config.MaxToolBytes {
		return r.handoff(ctx, message, requestHash, claim.Generation, "TOOL_RESULT_TOO_LARGE", base)
	}
	base.ToolResult = toolResult
	second, err := r.Model.CompleteToolTurn(ctx, first.ResponseID, r.Config.Instructions, []llmopenai.FunctionOutput{{CallID: call.CallID, Output: toolResult}}, r.Config.MaxOutputTokens)
	if err != nil {
		return "", r.failRetryable(ctx, message, requestHash, claim.Generation, "MODEL_COMPLETE_FAILED", err)
	}
	base.LLMResponseID = second.ResponseID
	base.InputTokens += second.InputTokens
	base.OutputTokens += second.OutputTokens
	if len(second.FunctionCalls) != 0 {
		return r.handoff(ctx, message, requestHash, claim.Generation, "UNEXPECTED_CONTINUATION_TOOL", base)
	}
	return r.completeText(ctx, message, requestHash, claim.Generation, second.OutputText, base)
}

func (r *Runtime) execute(ctx context.Context, message channels.Message, contact Contact, requestHash string, call llmopenai.FunctionCall) (result string, pending bool, terminal string, retryable error) {
	command := domainbind.CommandIdentity{IdempotencyKey: commandKey(message, call.Name), OccurredAt: message.OccurredAt}
	switch call.Name {
	case "book_appointment":
		var args struct {
			Service string `json:"service"`
			When    string `json:"when"`
		}
		if decodeStrict(call.Arguments, &args) != nil {
			return "", false, "INVALID_TOOL_ARGUMENTS", nil
		}
		state, err := r.Approval.Submit(approval.Request{TenantID: message.TenantID, ID: command.IdempotencyKey, Kind: approval.KindReservation, SubjectID: contact.SubjectID, AmountMinorUnits: 0, Requester: "agent-runtime", EvidenceSHA: requestHash})
		if errors.Is(err, approval.ErrDuplicate) {
			var ok bool
			state, ok = r.Approval.State(message.TenantID, command.IdempotencyKey)
			if !ok {
				return "", false, "APPROVAL_STATE_LOST", nil
			}
		} else if err != nil {
			return "", false, "APPROVAL_REJECTED", nil
		}
		if state != approval.StateApproved {
			return "La solicitud quedó pendiente de aprobación humana.", true, "", nil
		}
		value, err := r.Domain.BookAppointmentFor(ctx, message.TenantID, contact.Scope, command, agenttools.AppointmentInput{Service: args.Service, When: args.When})
		if err != nil {
			return "", false, "", err
		}
		return value, false, "", nil
	case "create_quote":
		var args struct {
			Product  string `json:"product"`
			Quantity int    `json:"quantity"`
		}
		if decodeStrict(call.Arguments, &args) != nil || args.Quantity < 1 {
			return "", false, "INVALID_TOOL_ARGUMENTS", nil
		}
		if args.Quantity != 1 {
			return "", false, "SINGLE_VEHICLE_QUOTE_REQUIRED", nil
		}
		value, err := r.Domain.CreateQuoteFor(ctx, message.TenantID, contact.Scope, command, agenttools.QuoteInput{Product: args.Product, Quantity: args.Quantity})
		if err != nil {
			return "", false, "", err
		}
		return value, false, "", nil
	case "order_status":
		var args struct {
			OrderID string `json:"order_id"`
		}
		if decodeStrict(call.Arguments, &args) != nil {
			return "", false, "INVALID_TOOL_ARGUMENTS", nil
		}
		value, err := r.Domain.OrderStatusFor(ctx, message.TenantID, contact.Scope, args.OrderID)
		if err != nil {
			return "", false, "", err
		}
		return value, false, "", nil
	default:
		return "", false, "UNAUTHORIZED_TOOL", nil
	}
}

func (r *Runtime) completeText(ctx context.Context, message channels.Message, requestHash string, generation int, text string, completion Completion) (string, error) {
	if len(text) > r.Config.MaxInputBytes || strings.TrimSpace(text) == "" || r.Safety.PostFlight(aifoundation.SafetyOutput{Content: text}) != aifoundation.SafetyAllow {
		return r.handoff(ctx, message, requestHash, generation, "SAFETY_POSTFLIGHT", completion)
	}
	completion.State, completion.ResponseText = StateCompleted, text
	if err := r.Store.Complete(ctx, message, requestHash, generation, completion); err != nil {
		return "", err
	}
	return text, nil
}

func (r *Runtime) handoff(ctx context.Context, message channels.Message, requestHash string, generation int, code string, completion Completion) (string, error) {
	completion.State, completion.ResponseText = StateHandedOff, r.Config.HandoffText
	completion.FailureCode = code
	if completion.ToolResult == "" {
		completion.ToolResult = code
	}
	if err := r.Store.Complete(ctx, message, requestHash, generation, completion); err != nil {
		return "", err
	}
	return completion.ResponseText, nil
}

func (r *Runtime) failRetryable(ctx context.Context, message channels.Message, requestHash string, generation int, code string, cause error) error {
	if err := r.Store.FailRetryable(ctx, message, requestHash, generation, code); err != nil {
		return errors.Join(cause, err)
	}
	return cause
}

func runtimeTools() []llmopenai.FunctionTool {
	return []llmopenai.FunctionTool{
		{Name: "book_appointment", Description: "Solicitar una cita para un servicio y momento definidos por el cliente.", Parameters: json.RawMessage(`{"type":"object","properties":{"service":{"type":"string"},"when":{"type":"string"}},"required":["service","when"],"additionalProperties":false}`)},
		{Name: "create_quote", Description: "Crear una cotización no vinculante para un producto y cantidad.", Parameters: json.RawMessage(`{"type":"object","properties":{"product":{"type":"string"},"quantity":{"type":"integer"}},"required":["product","quantity"],"additionalProperties":false}`)},
		{Name: "order_status", Description: "Consultar el estado de un pedido del contacto autenticado.", Parameters: json.RawMessage(`{"type":"object","properties":{"order_id":{"type":"string"}},"required":["order_id"],"additionalProperties":false}`)},
	}
}

func decodeStrict(raw []byte, destination any) error {
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destination); err != nil {
		return err
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return errors.New("conversationruntime: trailing tool arguments")
	}
	return nil
}

func buildInput(history []HistoryMessage, current string) string {
	var b strings.Builder
	for _, item := range history {
		fmt.Fprintf(&b, "Usuario: %s\nAsistente: %s\n", item.UserText, item.AssistantText)
	}
	b.WriteString("Usuario: ")
	b.WriteString(current)
	return b.String()
}

func hashMessage(message channels.Message) (string, error) {
	payload, err := json.Marshal(struct {
		Channel   string `json:"channel"`
		Tenant    string `json:"tenant"`
		External  string `json:"external"`
		Thread    string `json:"thread"`
		MessageID string `json:"message_id"`
		Occurred  string `json:"occurred_at"`
		Text      string `json:"text"`
	}{message.ChannelCode, message.TenantID, message.ExternalID, message.ThreadID, message.ProviderMessageID, message.OccurredAt.UTC().Format("2006-01-02T15:04:05.999999999Z07:00"), message.Text})
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(payload)
	return hex.EncodeToString(sum[:]), nil
}

func commandKey(message channels.Message, tool string) string {
	sum := sha256.Sum256([]byte(message.TenantID + "\x00" + message.ChannelCode + "\x00" + message.ProviderMessageID + "\x00" + tool))
	return hex.EncodeToString(sum[:])
}

func buildMessages(history []HistoryMessage, current string) []llmopenai.InputMessage {
	result := make([]llmopenai.InputMessage, 0, len(history)*2+1)
	for _, h := range history {
		result = append(result, llmopenai.InputMessage{Role: "user", Content: h.UserText}, llmopenai.InputMessage{Role: "assistant", Content: h.AssistantText})
	}
	return append(result, llmopenai.InputMessage{Role: "user", Content: current})
}
