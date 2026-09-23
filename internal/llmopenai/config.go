// Package llmopenai implements aifoundation.LLMProvider over an
// OpenAI-compatible Chat Completions endpoint. The API key is injected from
// the environment at runtime, never hardcoded; tests use a fake HTTP server.
package llmopenai

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"
)

// Config carries the runtime endpoint settings. The API key is a secret
// injected from the environment, never stored in source.
type Config struct {
	BaseURL string // e.g. https://api.openai.com/v1
	Model   string // exact model name, never "latest"/"auto"
	APIKey  string // secret, from env
	Timeout time.Duration
}

// Validate enforces an exact, safe configuration.
func (c Config) Validate() error {
	u, err := url.Parse(c.BaseURL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return errors.New("llmopenai: base url must be absolute http(s)")
	}
	m := strings.ToLower(strings.TrimSpace(c.Model))
	if m == "" || m == "latest" || m == "auto" {
		return fmt.Errorf("llmopenai: model must be an exact name, got %q", c.Model)
	}
	if strings.TrimSpace(c.APIKey) == "" {
		return errors.New("llmopenai: api key required (inject from environment)")
	}
	if c.Timeout <= 0 {
		return errors.New("llmopenai: timeout required")
	}
	return nil
}
