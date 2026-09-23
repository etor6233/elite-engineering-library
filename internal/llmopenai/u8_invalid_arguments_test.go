package llmopenai

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestU8ResponsesInvalidArgumentsClassification(t *testing.T) {
	for _, tc := range []struct {
		name, args string
		status     int
		invalid    bool
	}{{"duplicate", `{"service":"one","service":"two","when":"tomorrow"}`, 200, true}, {"missing", `{"service":"one"}`, 200, true}, {"malformed", `{`, 200, true}, {"valid", `{"service":"one","when":"tomorrow"}`, 200, false}, {"transport", `{}`, 503, false}} {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls++
				w.WriteHeader(tc.status)
				_ = json.NewEncoder(w).Encode(map[string]any{"id": "fixture", "status": "completed", "output": []any{map[string]any{"type": "function_call", "call_id": "call-1", "name": "book_appointment", "arguments": tc.args}}})
			}))
			defer server.Close()
			client, e := New(Config{BaseURL: server.URL, Model: "fixture", APIKey: "synthetic", Timeout: time.Second})
			if e != nil {
				t.Fatal(e)
			}
			got, e := client.StartToolTurn(context.Background(), ToolTurnRequest{Instructions: "Fixture only", Input: "Fixture request", Tools: []FunctionTool{{Name: "book_appointment", Description: "Fixture appointment proposal", Parameters: appointmentSchema}}, MaxOutputTokens: 32, PromptCacheKey: "fixture", StoreApproved: true})
			if errors.Is(e, ErrInvalidToolArguments) != tc.invalid || calls != 1 {
				t.Fatalf("classification %v calls=%d", e, calls)
			}
			if tc.invalid && (len(got.FunctionCalls) > 0 || e.Error() != ErrInvalidToolArguments.Error()) {
				t.Fatal("invalid arguments leaked into effects/diagnostics")
			}
			if tc.name == "valid" && (e != nil || len(got.FunctionCalls) != 1) {
				t.Fatal("valid rejected")
			}
			if tc.name == "transport" && e == nil {
				t.Fatal("transport success")
			}
		})
	}
}
