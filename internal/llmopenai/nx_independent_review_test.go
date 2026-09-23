package llmopenai

import (
	"context"
	"elite.local/enterprise/internal/aifoundation"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
)

const reviewSchema = `{"type":"object","properties":{"rows":{"type":"array","items":{"type":"object","properties":{"name":{"type":"string"},"count":{"type":"integer"}},"required":["name","count"],"additionalProperties":false}}},"required":["rows"],"additionalProperties":false}`

func reviewRequest(schema string) aifoundation.CompletionRequest {
	return aifoundation.CompletionRequest{Model: aifoundation.ModelRef{Provider: "openai", Model: "gpt-4.1-mini", Version: "v", Digest: strings.Repeat("a", 64)}, Messages: []aifoundation.Message{{Role: "user", Content: "Return JSON."}}, Schema: json.RawMessage(schema)}
}
func TestReviewPositiveAdapterRejectsInvalidBeforeCall(t *testing.T) {
	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1); w.WriteHeader(500) }))
	defer srv.Close()
	c, e := New(testConfig(srv.URL))
	if e != nil {
		t.Fatal(e)
	}
	for _, schema := range []string{`{"type":"object","properties":{"n":{"type":"string","minLength":2}},"additionalProperties":false}`, `{"type":"object","properties":{"rows":{"type":"array","items":{"type":"object","properties":{},"required":["missing"],"additionalProperties":false}}},"additionalProperties":false}`} {
		if _, e := c.Generate(context.Background(), reviewRequest(schema)); e == nil {
			t.Fatal("invalid schema accepted")
		}
	}
	if calls.Load() != 0 {
		t.Fatalf("provider called %d times", calls.Load())
	}
}
func TestReviewPositiveAdapterNestedAndMalformedOutput(t *testing.T) {
	for name, content := range map[string]string{"valid": `{"rows":[{"name":"a","count":1e2}]}`, "missingRequired": `{"rows":[{"name":"a"}]}`, "extraProperty": `{"rows":[{"name":"a","count":1,"extra":true}]}`, "wrongScalar": `{"rows":[{"name":"a","count":1.5}]}`, "duplicateKey": `{"rows":[],"rows":[]}`} {
		t.Run(name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var body chatRequest
				if e := json.NewDecoder(r.Body).Decode(&body); e != nil {
					t.Error(e)
				}
				if body.ResponseFormat == nil || body.ResponseFormat.Type != "json_object" {
					t.Error("expected local-validation JSON mode")
				}
				json.NewEncoder(w).Encode(map[string]any{"choices": []any{map[string]any{"message": map[string]any{"role": "assistant", "content": content}, "finish_reason": "stop"}}})
			}))
			defer srv.Close()
			c, e := New(testConfig(srv.URL))
			if e != nil {
				t.Fatal(e)
			}
			_, e = c.Generate(context.Background(), reviewRequest(reviewSchema))
			if (e == nil) != (name == "valid") {
				t.Fatalf("output expectation failed: %v", e)
			}
		})
	}
}
