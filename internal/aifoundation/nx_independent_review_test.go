package aifoundation

import (
	"context"
	"encoding/json"
	"testing"
)

const reviewNestedSchema = `{"type":"object","properties":{"profile":{"type":"object","properties":{"name":{"type":"string"},"active":{"type":"boolean"}},"required":["name"],"additionalProperties":false},"rows":{"type":"array","items":{"type":"object","properties":{"qty":{"type":"integer"},"tags":{"type":"array","items":{"type":"string"}}},"required":["qty","tags"],"additionalProperties":false}},"matrix":{"type":"array","items":{"type":"array","items":{"type":"number"}}},"optional":{"type":"null"}},"required":["profile","rows","matrix"],"additionalProperties":false}`
const reviewValidNested = `{"profile":{"name":"Ana","active":true},"rows":[{"qty":1e2,"tags":["a","b"]},{"qty":1.0,"tags":[]}],"matrix":[[1.25,-2],[]],"optional":null}`

type reviewProvider struct {
	calls   int
	content string
}

func (p *reviewProvider) Generate(_ context.Context, _ CompletionRequest) (CompletionResponse, error) {
	p.calls++
	return CompletionResponse{Content: json.RawMessage(p.content), FinishReason: "stop"}, nil
}
func TestReviewPositiveNestedSchema(t *testing.T) {
	if err := ValidateStructuredSchema([]byte(reviewNestedSchema)); err != nil {
		t.Fatal(err)
	}
	for _, instance := range []string{reviewValidNested, `{"profile":{"name":"Ana"},"rows":[],"matrix":[]}`} {
		if err := ValidateStructuredOutput([]byte(reviewNestedSchema), []byte(instance)); err != nil {
			t.Fatalf("valid nested instance: %v", err)
		}
	}
	for name, instance := range map[string]string{"nestedRequired": `{"profile":{},"rows":[],"matrix":[]}`, "nestedExtra": `{"profile":{"name":"Ana","extra":1},"rows":[],"matrix":[]}`, "wrongArrayItem": `{"profile":{"name":"Ana"},"rows":[{"qty":1.2,"tags":[]}],"matrix":[]}`, "nestedArrayScalar": `{"profile":{"name":"Ana"},"rows":[],"matrix":[[true]]}`, "nullArray": `{"profile":{"name":"Ana"},"rows":null,"matrix":[]}`} {
		t.Run(name, func(t *testing.T) {
			if ValidateStructuredOutput([]byte(reviewNestedSchema), []byte(instance)) == nil {
				t.Fatal("invalid nested output accepted")
			}
		})
	}
}
func TestReviewPositiveGatewayConfiguredSchema(t *testing.T) {
	g, m := newGateway(t, CompletionResponse{})
	p := &reviewProvider{content: reviewValidNested}
	g.Provider = p
	g.Schema = []byte(reviewNestedSchema)
	if _, err := g.Generate(context.Background(), CompletionRequest{Model: m, Messages: []Message{{Role: "user", Content: "hello"}}}); err != nil || p.calls != 1 {
		t.Fatalf("nested gateway err=%v calls=%d", err, p.calls)
	}
	for _, schema := range []string{`{"type":"object","properties":{"x":{"type":"string","minLength":2}},"additionalProperties":false}`, `{"type":"object","properties":{"x":{"type":"array","items":{"type":"object","properties":{},"required":["missing"],"additionalProperties":false}}},"additionalProperties":false}`} {
		p.calls = 0
		g.Schema = []byte(schema)
		if _, err := g.Generate(context.Background(), CompletionRequest{Model: m}); err == nil || p.calls != 0 {
			t.Fatalf("invalid configured schema err=%v calls=%d", err, p.calls)
		}
	}
}
func TestReviewDefectGatewayRequestSchema(t *testing.T) {
	for name, tc := range map[string]struct {
		schema, output string
		wantCalls      int
	}{"invalidSchemaPreCall": {`{"type":"object","properties":{"x":{"type":"string","minLength":2}},"additionalProperties":false}`, `{"x":"ok"}`, 0}, "invalidOutputPostCall": {`{"type":"object","properties":{"x":{"type":"string"}},"required":["x"],"additionalProperties":false}`, `{"x":123}`, 1}} {
		t.Run(name, func(t *testing.T) {
			g, m := newGateway(t, CompletionResponse{})
			p := &reviewProvider{content: tc.output}
			g.Provider = p
			g.Schema = nil
			resp, err := g.Generate(context.Background(), CompletionRequest{Model: m, Schema: json.RawMessage(tc.schema)})
			if err == nil || p.calls != tc.wantCalls {
				t.Fatalf("request schema not enforced: err=%v calls=%d wantCalls=%d content=%s", err, p.calls, tc.wantCalls, resp.Content)
			}
		})
	}
}

func TestNXGatewaySchemaPinsCannotConflict(t *testing.T) {
	g, m := newGateway(t, CompletionResponse{})
	p := &reviewProvider{content: reviewValidNested}
	g.Provider = p
	g.Schema = []byte(reviewNestedSchema)
	_, e := g.Generate(context.Background(), CompletionRequest{Model: m, Schema: json.RawMessage(`{"type":"object","properties":{},"additionalProperties":false}`)})
	if e == nil || p.calls != 0 {
		t.Fatal("conflicting schema replaced authority")
	}
}
