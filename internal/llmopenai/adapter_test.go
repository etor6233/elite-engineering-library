package llmopenai

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"elite.local/enterprise/internal/aifoundation"
)

func testConfig(serverURL string) Config {
	return Config{BaseURL: serverURL, Model: "gpt-4.1-mini", APIKey: "test-key", Timeout: 5 * time.Second}
}

func TestConfigValidate(t *testing.T) {
	if err := testConfig("https://api.openai.com/v1").Validate(); err != nil {
		t.Fatalf("valid config rejected: %v", err)
	}
	bad := testConfig("https://api.openai.com/v1")
	bad.Model = "latest"
	if err := bad.Validate(); err == nil {
		t.Fatal("mutable model accepted")
	}
	noKey := testConfig("https://api.openai.com/v1")
	noKey.APIKey = ""
	if err := noKey.Validate(); err == nil {
		t.Fatal("missing key accepted")
	}
	badURL := testConfig("not-a-url")
	if err := badURL.Validate(); err == nil {
		t.Fatal("bad url accepted")
	}
}

func TestGenerateRoundTrip(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/chat/completions" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("authorization") != "Bearer test-key" {
			t.Errorf("missing bearer")
		}
		var body chatRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body.Model != "gpt-4.1-mini" || body.ResponseFormat == nil || body.ResponseFormat.Type != "json_object" {
			t.Errorf("unexpected request: %+v", body)
		}
		w.Header().Set("content-type", "application/json")
		_, _ = w.Write([]byte(`{"choices":[{"message":{"role":"assistant","content":"{\"answer\":\"hola\"}"},"finish_reason":"stop"}]}`))
	}))
	defer srv.Close()

	c, err := New(testConfig(srv.URL))
	if err != nil {
		t.Fatal(err)
	}
	resp, err := c.Generate(context.Background(), aifoundation.CompletionRequest{
		Model:    aifoundation.ModelRef{Provider: "openai", Model: "gpt-4.1-mini", Version: "v", Digest: strings.Repeat("a", 64)},
		Messages: []aifoundation.Message{{Role: "user", Content: "hola"}},
		Schema:   json.RawMessage(`{"type":"object","properties":{"answer":{"type":"string"}},"required":["answer"],"additionalProperties":false}`),
	})
	if err != nil {
		t.Fatal(err)
	}
	if string(resp.Content) != `{"answer":"hola"}` || resp.FinishReason != "stop" {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestGenerateFailsClosedOnErrorStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":{"message":"bad key"}}`))
	}))
	defer srv.Close()

	c, _ := New(testConfig(srv.URL))
	if _, err := c.Generate(context.Background(), aifoundation.CompletionRequest{
		Model:    aifoundation.ModelRef{Provider: "openai", Model: "gpt-4.1-mini", Version: "v", Digest: strings.Repeat("a", 64)},
		Messages: []aifoundation.Message{{Role: "user", Content: "x"}},
	}); err == nil {
		t.Fatal("expected error on 401")
	}
}

func TestGenerateRejectsEmptyContent(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		_, _ = w.Write([]byte(`{"choices":[{"message":{"role":"assistant","content":""},"finish_reason":"stop"}]}`))
	}))
	defer srv.Close()

	c, _ := New(testConfig(srv.URL))
	if _, err := c.Generate(context.Background(), aifoundation.CompletionRequest{
		Model:    aifoundation.ModelRef{Provider: "openai", Model: "gpt-4.1-mini", Version: "v", Digest: strings.Repeat("a", 64)},
		Messages: []aifoundation.Message{{Role: "user", Content: "x"}},
	}); err == nil {
		t.Fatal("expected error on empty content")
	}
}
