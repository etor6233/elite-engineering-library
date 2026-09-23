package llmopenai

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var appointmentSchema = json.RawMessage(`{"type":"object","properties":{"service":{"type":"string"},"when":{"type":"string"}},"required":["service","when"],"additionalProperties":false}`)

func TestResponsesToolTurnAndContinuationUseOfficialWireContract(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if r.Method != http.MethodPost || r.URL.Path != "/responses" || r.Header.Get("authorization") != "Bearer secret" {
			t.Fatalf("unexpected request %s %s auth=%q", r.Method, r.URL.Path, r.Header.Get("authorization"))
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		w.Header().Set("content-type", "application/json")
		if requests == 1 {
			tools := body["tools"].([]any)
			tool := tools[0].(map[string]any)
			if body["model"] != "gpt-fixed" || body["store"] != true || body["parallel_tool_calls"] != false || body["tool_choice"] != "auto" || tool["type"] != "function" || tool["strict"] != true {
				t.Fatalf("invalid initial contract: %#v", body)
			}
			_, _ = w.Write([]byte(`{"id":"resp_1","status":"completed","output":[{"type":"function_call","call_id":"call_1","name":"book_appointment","arguments":"{\"service\":\"service\",\"when\":\"2026-09-05T15:00:00Z\"}"}],"usage":{"input_tokens":20,"output_tokens":10,"total_tokens":30}}`))
			return
		}
		if body["instructions"] != "Use only offered tools." || body["previous_response_id"] != "resp_1" || body["tool_choice"] != "none" || body["store"] != true {
			t.Fatalf("invalid continuation: %#v", body)
		}
		input := body["input"].([]any)[0].(map[string]any)
		if input["type"] != "function_call_output" || input["call_id"] != "call_1" || input["output"] != "appointment-created" {
			t.Fatalf("invalid function output: %#v", input)
		}
		_, _ = w.Write([]byte(`{"id":"resp_2","status":"completed","output":[{"type":"message","content":[{"type":"output_text","text":"Tu cita quedó registrada."}]}],"usage":{"input_tokens":12,"output_tokens":6,"total_tokens":18}}`))
	}))
	defer server.Close()
	client, err := New(Config{BaseURL: server.URL, Model: "gpt-fixed", APIKey: "secret", Timeout: time.Second})
	if err != nil {
		t.Fatal(err)
	}
	start, err := client.StartToolTurn(context.Background(), ToolTurnRequest{Instructions: "Use only offered tools.", Input: "Necesito una cita", Tools: []FunctionTool{{Name: "book_appointment", Description: "Book an approved service appointment.", Parameters: appointmentSchema}}, MaxOutputTokens: 256, PromptCacheKey: "tenant-hash", StoreApproved: true})
	if err != nil || len(start.FunctionCalls) != 1 || start.FunctionCalls[0].Name != "book_appointment" || start.TotalTokens != 30 {
		t.Fatalf("start=%+v err=%v", start, err)
	}
	final, err := client.CompleteToolTurn(context.Background(), start.ResponseID, "Use only offered tools.", []FunctionOutput{{CallID: start.FunctionCalls[0].CallID, Output: "appointment-created"}}, 256)
	if err != nil || final.OutputText != "Tu cita quedó registrada." || final.TotalTokens != 18 || requests != 2 {
		t.Fatalf("final=%+v err=%v requests=%d", final, err, requests)
	}
}

func TestResponsesToolTurnRejectsUnapprovedStorageAndInvalidCalls(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"id":"resp_1","status":"completed","output":[{"type":"function_call","call_id":"call_1","name":"delete_everything","arguments":"{}"}]}`))
	}))
	defer server.Close()
	client, err := New(Config{BaseURL: server.URL, Model: "gpt-fixed", APIKey: "secret", Timeout: time.Second})
	if err != nil {
		t.Fatal(err)
	}
	req := ToolTurnRequest{Instructions: "bounded", Input: "hello", Tools: []FunctionTool{{Name: "book_appointment", Description: "Book appointment", Parameters: appointmentSchema}}, MaxOutputTokens: 64}
	if _, err := client.StartToolTurn(context.Background(), req); err == nil {
		t.Fatal("unapproved storage accepted")
	}
	req.StoreApproved = true
	if _, err := client.StartToolTurn(context.Background(), req); err == nil {
		t.Fatal("unoffered function accepted")
	}
}

func TestResponsesToolTurnRejectsArgumentsOutsideClosedSchema(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"id":"resp_1","status":"completed","output":[{"type":"function_call","call_id":"call_1","name":"book_appointment","arguments":"{\"service\":\"service\",\"when\":\"tomorrow\",\"admin\":true}"}]}`))
	}))
	defer server.Close()
	client, err := New(Config{BaseURL: server.URL, Model: "gpt-fixed", APIKey: "secret", Timeout: time.Second})
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.StartToolTurn(context.Background(), ToolTurnRequest{Instructions: "bounded", Input: "hello", Tools: []FunctionTool{{Name: "book_appointment", Description: "Book appointment", Parameters: appointmentSchema}}, MaxOutputTokens: 64, StoreApproved: true})
	if err == nil {
		t.Fatal("arguments with unknown key accepted")
	}
}

func TestResponsesContinuationRejectsSecondToolCall(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"id":"resp_2","status":"completed","output":[{"type":"function_call","call_id":"call_2","name":"book_appointment","arguments":"{\"service\":\"x\",\"when\":\"y\"}"}]}`))
	}))
	defer server.Close()
	client, err := New(Config{BaseURL: server.URL, Model: "gpt-fixed", APIKey: "secret", Timeout: time.Second})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := client.CompleteToolTurn(context.Background(), "resp_1", "Use only offered tools.", []FunctionOutput{{CallID: "call_1", Output: "done"}}, 64); err == nil {
		t.Fatal("second tool call accepted")
	}
}

// AUTHORED regression: schema validation is distinct from validating an empty instance.
func TestNXRequiredToolSchemaUsesSchemaBoundary(t *testing.T) {
	schema := json.RawMessage(`{"type":"object","properties":{"item":{"type":"object","properties":{"names":{"type":"array","items":{"type":"string"}}},"required":["names"],"additionalProperties":false}},"required":["item"],"additionalProperties":false}`)
	tools, _, err := validateTools([]FunctionTool{{Name: "review_item", Description: "Review an item", Parameters: schema}})
	if err != nil || len(tools) != 1 {
		t.Fatalf("valid required nested schema rejected: %v", err)
	}
	bad := json.RawMessage(`{"type":"object","properties":{"x":{"type":"string","pattern":".*"}},"required":["x"],"additionalProperties":false}`)
	if _, _, err = validateTools([]FunctionTool{{Name: "bad", Description: "Invalid subset", Parameters: bad}}); err == nil {
		t.Fatal("unsupported schema keyword accepted")
	}
}
