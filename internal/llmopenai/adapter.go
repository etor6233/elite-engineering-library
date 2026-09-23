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

// Client is an OpenAI-compatible Chat Completions adapter.
type Client struct {
	cfg    Config
	client *http.Client
}

// New returns a validated client.
func New(cfg Config) (*Client, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return &Client{cfg: cfg, client: &http.Client{Timeout: cfg.Timeout}}, nil
}

type chatRequest struct {
	Model          string              `json:"model"`
	Messages       []chatMessage       `json:"messages"`
	MaxTokens      int                 `json:"max_tokens,omitempty"`
	Temperature    *float64            `json:"temperature,omitempty"`
	Seed           *int64              `json:"seed,omitempty"`
	ResponseFormat *chatResponseFormat `json:"response_format,omitempty"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResponseFormat struct {
	Type string `json:"type"`
}

type chatResponse struct {
	Choices []struct {
		Message      chatMessage `json:"message"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
}

// Generate maps aifoundation.CompletionRequest to a Chat Completions call and
// returns the raw content. It fails closed on non-2xx, empty content or an
// unexpected payload.
func (c *Client) Generate(ctx context.Context, req aifoundation.CompletionRequest) (aifoundation.CompletionResponse, error) {
	if err := req.Model.Validate(); err != nil {
		return aifoundation.CompletionResponse{}, err
	}
	if req.Model.Provider != "openai" || req.Model.Model != c.cfg.Model {
		return aifoundation.CompletionResponse{}, errors.New("llmopenai: request/config model mismatch")
	}
	if len(req.Schema) > 0 {
		if err := aifoundation.ValidateStructuredSchema(req.Schema); err != nil {
			return aifoundation.CompletionResponse{}, err
		}
	}
	messages := make([]chatMessage, 0, len(req.Messages))
	for _, m := range req.Messages {
		messages = append(messages, chatMessage{Role: m.Role, Content: m.Content})
	}
	body := chatRequest{
		Model:       c.cfg.Model,
		Messages:    messages,
		MaxTokens:   req.MaxOutputTokens,
		Temperature: req.Temperature,
		Seed:        req.Seed,
	}
	if len(req.Schema) > 0 {
		body.ResponseFormat = &chatResponseFormat{Type: "json_object"}
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return aifoundation.CompletionResponse{}, err
	}
	endpoint := strings.TrimRight(c.cfg.BaseURL, "/") + "/chat/completions"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(payload))
	if err != nil {
		return aifoundation.CompletionResponse{}, err
	}
	httpReq.Header.Set("content-type", "application/json")
	httpReq.Header.Set("authorization", "Bearer "+c.cfg.APIKey)

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return aifoundation.CompletionResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return aifoundation.CompletionResponse{}, fmt.Errorf("llmopenai: status %d", resp.StatusCode)
	}

	var parsed chatResponse
	data, err := io.ReadAll(io.LimitReader(resp.Body, (8<<20)+1))
	if err != nil {
		return aifoundation.CompletionResponse{}, err
	}
	if len(data) > 8<<20 {
		return aifoundation.CompletionResponse{}, errors.New("llmopenai: response limit")
	}
	if err := json.Unmarshal(data, &parsed); err != nil {
		return aifoundation.CompletionResponse{}, err
	}
	if len(parsed.Choices) == 0 {
		return aifoundation.CompletionResponse{}, errors.New("llmopenai: empty choices")
	}
	if parsed.Choices[0].FinishReason != "stop" {
		return aifoundation.CompletionResponse{}, errors.New("llmopenai: incomplete response")
	}
	content := parsed.Choices[0].Message.Content
	if strings.TrimSpace(content) == "" {
		return aifoundation.CompletionResponse{}, errors.New("llmopenai: empty content")
	}
	if len(req.Schema) > 0 {
		if err := aifoundation.ValidateStructuredOutput(req.Schema, []byte(content)); err != nil {
			return aifoundation.CompletionResponse{}, err
		}
	}
	return aifoundation.CompletionResponse{
		Content:      json.RawMessage(content),
		FinishReason: parsed.Choices[0].FinishReason,
	}, nil
}
