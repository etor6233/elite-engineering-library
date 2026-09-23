package llmopenai

import (
	"context"
	"elite.local/enterprise/internal/aifoundation"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func nxReq() aifoundation.CompletionRequest {
	return aifoundation.CompletionRequest{Model: aifoundation.ModelRef{Provider: "openai", Model: "gpt-4.1-mini", Version: "v", Digest: strings.Repeat("a", 64)}, Messages: []aifoundation.Message{{Role: "user", Content: "JSON response"}}, MaxOutputTokens: 20}
}
func TestNXModelMismatchRejectedBeforeCall(t *testing.T) {
	n := 0
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n++
		w.Write([]byte(`{"choices":[{"message":{"role":"assistant","content":"{}"},"finish_reason":"stop"}]}`))
	}))
	defer s.Close()
	c, _ := New(testConfig(s.URL))
	r := nxReq()
	r.Model.Model = "other"
	_, err := c.Generate(context.Background(), r)
	if err == nil || n != 0 {
		t.Fatalf("mismatch reached provider n=%d err=%v", n, err)
	}
}
func TestNXProviderOversizeAndIncompleteRejected(t *testing.T) {
	for _, body := range []string{`{"choices":[{"message":{"role":"assistant","content":"` + strings.Repeat("x", 8388609) + `"},"finish_reason":"stop"}]}`, `{"choices":[{"message":{"role":"assistant","content":"{}"},"finish_reason":"length"}]}`} {
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte(body)) }))
		c, _ := New(testConfig(s.URL))
		if _, err := c.Generate(context.Background(), nxReq()); err == nil {
			t.Error("oversize/incomplete accepted")
		}
		s.Close()
	}
}
