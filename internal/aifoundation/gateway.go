package aifoundation

import (
	"bytes"
	"context"
	"errors"
)

// ModelGateway enforces the full serving chain for every generation:
// active-pin check, pre-flight safety, provider call, post-flight safety and
// closed-schema structured output validation.
type ModelGateway struct {
	Provider LLMProvider
	Registry *VersionRegistry
	Safety   SafetyGate
	Schema   []byte // closed object schema applied to every response
}

// Generate runs the guarded generation chain.
func (g *ModelGateway) Generate(ctx context.Context, req CompletionRequest) (CompletionResponse, error) {
	if g.Provider == nil {
		return CompletionResponse{}, errors.New("aifoundation: nil provider")
	}
	if g.Registry == nil {
		return CompletionResponse{}, errors.New("aifoundation: nil registry")
	}
	if g.Safety == nil {
		return CompletionResponse{}, errors.New("aifoundation: nil safety gate")
	}
	if err := req.Model.Validate(); err != nil {
		return CompletionResponse{}, err
	}
	active, ok := g.Registry.Active()
	if !ok || active != req.Model {
		return CompletionResponse{}, errors.New("aifoundation: model is not the active version")
	}
	if g.Safety.PreFlight(SafetyInput{UserText: lastUserText(req.Messages)}) != SafetyAllow {
		return CompletionResponse{}, errors.New("aifoundation: pre-flight safety block")
	}
	// One effective schema governs preflight, wire and returned instance.
	// Distinct non-empty pins are a configuration conflict, never an override.
	schema := req.Schema
	if len(g.Schema) > 0 {
		if len(schema) > 0 && !bytes.Equal(schema, g.Schema) {
			return CompletionResponse{}, errors.New("aifoundation: conflicting schema pins")
		}
		schema = g.Schema
	}
	if len(schema) > 0 {
		if err := ValidateStructuredSchema(schema); err != nil {
			return CompletionResponse{}, err
		}
		schema = append([]byte(nil), schema...)
		req.Schema = append([]byte(nil), schema...)
	}
	resp, err := g.Provider.Generate(ctx, req)
	if err != nil {
		return CompletionResponse{}, err
	}
	if g.Safety.PostFlight(SafetyOutput{Content: string(resp.Content)}) != SafetyAllow {
		return CompletionResponse{}, errors.New("aifoundation: post-flight safety block")
	}
	if len(schema) > 0 {
		if err := ValidateStructuredOutput(schema, resp.Content); err != nil {
			return CompletionResponse{}, err
		}
	}
	return resp, nil
}

func lastUserText(msgs []Message) string {
	for i := len(msgs) - 1; i >= 0; i-- {
		if msgs[i].Role == "user" {
			return msgs[i].Content
		}
	}
	return ""
}
