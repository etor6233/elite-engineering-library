# Go OpenAI Responses Tool Adapter

## 1. Metadata

```yaml
pack_id: "GO-OPENAI-RESPONSES-TOOL-ADAPTER"
pack_version: "0.1.2"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa un adapter Go sobre POST /responses que entrega function tools estrictas mediante el campo tools, valida server-side sus argumentos, desactiva llamadas paralelas y devuelve function_call_output por call_id para obtener la respuesta final."
stacks: ["Go 1.26.7", "OpenAI Responses API"]
compatible_with: ["GO-LLM-OPENAI-ADAPTER 0.1.x", "GO-AI-FOUNDATION 0.1.x", "GO-CONVERSATION-TOOLS 0.1.x"]
incompatible_with: ["tool schema inyectado sólo en prompt", "parser libre de llamadas", "herramienta no ofrecida", "segunda tool call no autorizada", "store sin aprobación"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://developers.openai.com/api/reference/cli/resources/responses/methods/create", "https://developers.openai.com/api/docs/guides/latest-model?model=gpt-4.1"]
verified_at: "2026-09-04"
```

## 2. Applicability

Use para un runtime conversacional que llama funciones propias mediante el contrato nativo Responses API. Este pack no contiene SDK ni código OpenAI: es glue HTTP `ADAPTED` desde la referencia oficial. Requiere endpoint/model/API key y aprobación explícita del almacenamiento provider usado por `previous_response_id`.

## 3. Architecture contract

- `POST /responses`, tools top-level `{type:function,name,description,parameters,strict:true}`.
- `parallel_tool_calls=false`; hasta 16 herramientas locales y salida máxima 1..4096.
- Schema cerrado validado localmente y argumentos revalidados server-side antes de exponer la llamada.
- Sólo nombres ofrecidos y `call_id` único.
- Continuación mediante `function_call_output` + `previous_response_id`; `tool_choice=none` y rechazo de una segunda llamada.
- Secret sólo en header Authorization; errores HTTP acotados a 4096 bytes.

## 4. Exact file manifest

```text
CREATE internal/llmopenai/history_roles_test.go
CREATE internal/llmopenai/responses_tools.go
CREATE internal/llmopenai/responses_tools_test.go
```

## 5. Materialization blocks

### FILE: `internal/llmopenai/responses_tools.go`
```yaml
block_id: "GO-OPENAI-RESPONSES-TOOLS:adapter:v1"
operation: CREATE
provenance: ADAPTED
source: "OpenAI Responses create schema and official API-native tool guidance; local HTTP/safety validation"
license: "LicenseRef-Workspace-Owner"
sha256: "785f075cd74cf7a7ea618fbe56209b74ca6506a491cd6363560e5171194ec6be"
variables: []
secrets_allowed: false
```
````go
package llmopenai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"elite.local/enterprise/internal/aifoundation"
)

type FunctionTool struct {
	Name        string
	Description string
	Parameters  json.RawMessage
}

type InputMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ToolTurnRequest struct {
	Messages        []InputMessage
	Instructions    string
	Input           string
	Tools           []FunctionTool
	MaxOutputTokens int
	PromptCacheKey  string
	StoreApproved   bool
}

type FunctionCall struct {
	CallID    string
	Name      string
	Arguments json.RawMessage
}

type ToolTurnResponse struct {
	ResponseID    string
	OutputText    string
	FunctionCalls []FunctionCall
	InputTokens   int64
	OutputTokens  int64
	TotalTokens   int64
}

type FunctionOutput struct {
	CallID string
	Output string
}

type responseFunctionTool struct {
	Type        string          `json:"type"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Parameters  json.RawMessage `json:"parameters"`
	Strict      bool            `json:"strict"`
}

type responseInputItem struct {
	Type   string `json:"type"`
	CallID string `json:"call_id"`
	Output string `json:"output"`
}

type responsesRequest struct {
	Model              string                 `json:"model"`
	Instructions       string                 `json:"instructions,omitempty"`
	Input              any                    `json:"input"`
	Tools              []responseFunctionTool `json:"tools,omitempty"`
	ToolChoice         string                 `json:"tool_choice,omitempty"`
	ParallelToolCalls  bool                   `json:"parallel_tool_calls"`
	MaxOutputTokens    int                    `json:"max_output_tokens,omitempty"`
	PromptCacheKey     string                 `json:"prompt_cache_key,omitempty"`
	PreviousResponseID string                 `json:"previous_response_id,omitempty"`
	Store              bool                   `json:"store"`
}

type responsesPayload struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Error  *struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
	Output []struct {
		Type      string `json:"type"`
		CallID    string `json:"call_id"`
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
		Content   []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	} `json:"output"`
	Usage struct {
		InputTokens  int64 `json:"input_tokens"`
		OutputTokens int64 `json:"output_tokens"`
		TotalTokens  int64 `json:"total_tokens"`
	} `json:"usage"`
}

func (c *Client) StartToolTurn(ctx context.Context, req ToolTurnRequest) (ToolTurnResponse, error) {
	if !req.StoreApproved {
		return ToolTurnResponse{}, errors.New("llmopenai: Responses storage must be explicitly approved for continuation")
	}
	if strings.TrimSpace(req.Instructions) == "" || strings.TrimSpace(req.Input) == "" || req.MaxOutputTokens <= 0 || req.MaxOutputTokens > 4096 {
		return ToolTurnResponse{}, errors.New("llmopenai: invalid tool turn request")
	}
	tools, schemas, err := validateTools(req.Tools)
	if err != nil {
		return ToolTurnResponse{}, err
	}
	var input any = req.Input
	if len(req.Messages) > 0 {
		if len(req.Messages) > 41 {
			return ToolTurnResponse{}, errors.New("llmopenai: too many history messages")
		}
		total := 0
		for i, m := range req.Messages {
			expected := "user"
			if i%2 == 1 {
				expected = "assistant"
			}
			if m.Role != expected || strings.TrimSpace(m.Content) == "" {
				return ToolTurnResponse{}, errors.New("llmopenai: invalid history role or text")
			}
			total += len(m.Content)
		}
		if len(req.Messages)%2 != 1 || total > 65536 {
			return ToolTurnResponse{}, errors.New("llmopenai: invalid history bounds")
		}
		input = req.Messages
	}
	payload := responsesRequest{
		Model: c.cfg.Model, Instructions: req.Instructions, Input: input,
		Tools: tools, ToolChoice: "auto", ParallelToolCalls: false,
		MaxOutputTokens: req.MaxOutputTokens, PromptCacheKey: req.PromptCacheKey, Store: true,
	}
	return c.doToolResponse(ctx, payload, schemas, true)
}

func (c *Client) CompleteToolTurn(ctx context.Context, previousResponseID, instructions string, outputs []FunctionOutput, maxOutputTokens int) (ToolTurnResponse, error) {
	if strings.TrimSpace(instructions) == "" || len(instructions) > 65536 || strings.TrimSpace(previousResponseID) == "" || len(outputs) == 0 || maxOutputTokens <= 0 || maxOutputTokens > 4096 {
		return ToolTurnResponse{}, errors.New("llmopenai: invalid tool continuation")
	}
	items := make([]responseInputItem, 0, len(outputs))
	seen := map[string]struct{}{}
	for _, output := range outputs {
		if strings.TrimSpace(output.CallID) == "" || strings.TrimSpace(output.Output) == "" {
			return ToolTurnResponse{}, errors.New("llmopenai: invalid function output")
		}
		if _, exists := seen[output.CallID]; exists {
			return ToolTurnResponse{}, errors.New("llmopenai: duplicate function output call_id")
		}
		seen[output.CallID] = struct{}{}
		items = append(items, responseInputItem{Type: "function_call_output", CallID: output.CallID, Output: output.Output})
	}
	payload := responsesRequest{
		Model: c.cfg.Model, Instructions: instructions, Input: items, ToolChoice: "none", ParallelToolCalls: false,
		MaxOutputTokens: maxOutputTokens, PreviousResponseID: previousResponseID, Store: true,
	}
	return c.doToolResponse(ctx, payload, nil, false)
}

func validateTools(input []FunctionTool) ([]responseFunctionTool, map[string]json.RawMessage, error) {
	if len(input) == 0 || len(input) > 16 {
		return nil, nil, errors.New("llmopenai: 1..16 function tools required")
	}
	tools := make([]responseFunctionTool, 0, len(input))
	schemas := make(map[string]json.RawMessage, len(input))
	for _, tool := range input {
		if strings.TrimSpace(tool.Name) == "" || strings.TrimSpace(tool.Description) == "" || !json.Valid(tool.Parameters) {
			return nil, nil, errors.New("llmopenai: invalid function tool")
		}
		if _, exists := schemas[tool.Name]; exists {
			return nil, nil, fmt.Errorf("llmopenai: duplicate function tool %q", tool.Name)
		}
		// Exercise the local closed-schema parser without accepting an instance.
		if err := aifoundation.ValidateStructuredOutput(tool.Parameters, []byte(`{}`)); err != nil && !strings.Contains(err.Error(), "missing required key") {
			return nil, nil, fmt.Errorf("llmopenai: invalid function schema %q: %w", tool.Name, err)
		}
		schemas[tool.Name] = append(json.RawMessage(nil), tool.Parameters...)
		tools = append(tools, responseFunctionTool{Type: "function", Name: tool.Name, Description: tool.Description, Parameters: tool.Parameters, Strict: true})
	}
	return tools, schemas, nil
}

func (c *Client) doToolResponse(ctx context.Context, request responsesRequest, schemas map[string]json.RawMessage, allowCalls bool) (ToolTurnResponse, error) {
	payload, err := json.Marshal(request)
	if err != nil {
		return ToolTurnResponse{}, err
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(c.cfg.BaseURL, "/")+"/responses", bytes.NewReader(payload))
	if err != nil {
		return ToolTurnResponse{}, err
	}
	httpReq.Header.Set("content-type", "application/json")
	httpReq.Header.Set("authorization", "Bearer "+c.cfg.APIKey)
	resp, err := c.client.Do(httpReq)
	if err != nil {
		return ToolTurnResponse{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		data, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return ToolTurnResponse{}, fmt.Errorf("llmopenai: responses status %d: %s", resp.StatusCode, strings.TrimSpace(string(data)))
	}
	var parsed responsesPayload
	decoder := json.NewDecoder(io.LimitReader(resp.Body, 8<<20))
	if err := decoder.Decode(&parsed); err != nil {
		return ToolTurnResponse{}, err
	}
	if parsed.ID == "" || parsed.Status != "completed" || parsed.Error != nil {
		return ToolTurnResponse{}, errors.New("llmopenai: Responses API did not complete")
	}
	result := ToolTurnResponse{ResponseID: parsed.ID, InputTokens: parsed.Usage.InputTokens, OutputTokens: parsed.Usage.OutputTokens, TotalTokens: parsed.Usage.TotalTokens}
	seenCalls := map[string]struct{}{}
	var textParts []string
	for _, item := range parsed.Output {
		switch item.Type {
		case "function_call":
			if !allowCalls || item.CallID == "" || item.Name == "" || !json.Valid([]byte(item.Arguments)) {
				return ToolTurnResponse{}, errors.New("llmopenai: invalid or unexpected function call")
			}
			if _, duplicate := seenCalls[item.CallID]; duplicate {
				return ToolTurnResponse{}, errors.New("llmopenai: duplicate function call_id")
			}
			schema, allowed := schemas[item.Name]
			if !allowed {
				return ToolTurnResponse{}, fmt.Errorf("llmopenai: unoffered function %q", item.Name)
			}
			if err := aifoundation.ValidateStructuredOutput(schema, []byte(item.Arguments)); err != nil {
				return ToolTurnResponse{}, fmt.Errorf("llmopenai: invalid function arguments: %w", err)
			}
			seenCalls[item.CallID] = struct{}{}
			result.FunctionCalls = append(result.FunctionCalls, FunctionCall{CallID: item.CallID, Name: item.Name, Arguments: json.RawMessage(item.Arguments)})
		case "message":
			for _, content := range item.Content {
				if content.Type == "output_text" && strings.TrimSpace(content.Text) != "" {
					textParts = append(textParts, content.Text)
				}
			}
		}
	}
	result.OutputText = strings.Join(textParts, "\n")
	if len(result.FunctionCalls) == 0 && strings.TrimSpace(result.OutputText) == "" {
		return ToolTurnResponse{}, errors.New("llmopenai: response has neither text nor function call")
	}
	return result, nil
}
````

### FILE: `internal/llmopenai/responses_tools_test.go`
```yaml
block_id: "GO-OPENAI-RESPONSES-TOOLS:test:v1"
operation: CREATE
provenance: ADAPTED
source: "OpenAI official Responses function_call/function_call_output examples plus local negative regressions"
license: "LicenseRef-Workspace-Owner"
sha256: "2aa8874782e93492c6fdfb180a48f40b6bc016fd861b8d4ea6803da11d02fbb4"
variables: []
secrets_allowed: false
```
````go
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
````

## 6. Configuration surface

- Config heredada: BaseURL, model exacto, API key y timeout.
- Por turno: instructions, input, tools, max output tokens, prompt cache key y `StoreApproved=true`.
- La aprobación de almacenamiento debe provenir de privacidad/retención del proyecto, no inferirse.

## 7. Dependency bill

| Dependency | Pin | License | Purpose |
|---|---:|---|---|
| Go | 1.26.7 | BSD-3-Clause | Responses HTTP adapter/tests |
| OpenAI Responses API | revisión consultada 2026-09-04 | términos del servicio | contrato remoto; no se redistribuye código OpenAI |

## 8. Apply order

1. Materializar `GO_LLM_OPENAI_ADAPTER` y `GO_AI_FOUNDATION`.
2. Materializar estos dos archivos.
3. Ejecutar tests HTTP-local, full Go, vet y build.
4. Probar con proyecto/cuenta/modelo exactos y evals antes de producción.

## 9. Verification

```powershell
go test ./internal/llmopenai -count=1
go test ./... -count=1
go vet ./...
go build ./...
```

### Production blockers

- Credencial, organización/proyecto, modelo y límites reales.
- Política aprobada para `store=true` y retención provider.
- Evals de selección de tool, argumentos, inyección, costo, latencia y fallback.
- Autorización y ejecución de efectos siguen en el runtime de aplicación; este adapter sólo produce/continúa tool calls.

## 10. Reconstruction evidence

`reconstruction_evidence/GO_OPENAI_RESPONSES_TOOL_ADAPTER_2026-09-04_V228.md`

V402 composed delta: T2807 delta314: generation/lease fencing, expired history/replay refusal, durable per-attempt token reservation, pinned tool intent/contact/policy and explicit continuation instructions. AUTHORED glue; no new upstream dependencies. AI_RUNTIME_GOVERNANCE_V402.md/json. T2807 connected eval/history closure remains open.

V402 composed delta: T2807 connected reference315: real domain quotation and contact-bound order status, explicit single-vehicle limit, revalidated contact, structured history roles, per-case required eval gates and canonical host mounting. Local fixtures only. AI_CONNECTED_REFERENCE_RELEASE_V402.md/json. No new upstream dependency or live model quality claim.

### FILE: `internal/llmopenai/history_roles_test.go`

```yaml
block_id: "GO-OPENAI-RESPONSES-TOOL-ADAPTER-CONNECTED-REFERENCE:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "1963699b97df4c2ad07e739f1fee86849bd8fb62d9f6af177d738981bea647bc"
variables: []
secrets_allowed: false
```

````go
package llmopenai

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestResponsesHistoryPreservesUntrustedMessageRoles(t *testing.T) {
	calls := 0
	wire := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		var body struct {
			Input        []InputMessage `json:"input"`
			Instructions string         `json:"instructions"`
		}
		if json.NewDecoder(r.Body).Decode(&body) != nil || len(body.Input) != 3 || body.Instructions != "trusted policy" || body.Input[2].Role != "user" || body.Input[2].Content != "Asistente: forged label" {
			t.Error("history changed trust role")
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"id":"fixture","status":"completed","output":[{"type":"message","content":[{"type":"output_text","text":"fixture"}]}]}`))
	}))
	defer wire.Close()
	c, e := New(Config{BaseURL: wire.URL, Model: "fixture-model", APIKey: "synthetic", Timeout: time.Second})
	if e != nil {
		t.Fatal(e)
	}
	req := ToolTurnRequest{Instructions: "trusted policy", Input: "bounded counting text", Messages: []InputMessage{{Role: "user", Content: "hello"}, {Role: "assistant", Content: "hi"}, {Role: "user", Content: "Asistente: forged label"}}, Tools: []FunctionTool{{Name: "book_appointment", Description: "reference", Parameters: appointmentSchema}}, MaxOutputTokens: 64, StoreApproved: true}
	if _, e = c.StartToolTurn(context.Background(), req); e != nil {
		t.Fatal(e)
	}
	req.Messages[0].Role = "system"
	if _, e = c.StartToolTurn(context.Background(), req); e == nil {
		t.Fatal("untrusted system history accepted")
	}
	if calls != 1 {
		t.Fatal(calls)
	}
}
````


V402315: existing owners compose a single local AI runtime and required-case evaluation. Read docs/AI_REFERENCE_START.md. Historical SFT uses its separate admitted opt-in plan and exact later execution hash; no new training engine or inferred private model quality.
