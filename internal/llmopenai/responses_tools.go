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

// ErrInvalidToolArguments classifies a completed provider response whose offered
// tool arguments fail the pinned schema. It contains no argument values.
// Callers may use errors.Is through wrappers; transport failures stay distinct.
var ErrInvalidToolArguments = errors.New("llmopenai: invalid function arguments")

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
		// Validate the schema itself; an empty instance cannot validate required fields.
		if err := aifoundation.ValidateStructuredSchema(tool.Parameters); err != nil {
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
			if !allowCalls || item.CallID == "" || item.Name == "" {
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
				return ToolTurnResponse{}, ErrInvalidToolArguments
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
