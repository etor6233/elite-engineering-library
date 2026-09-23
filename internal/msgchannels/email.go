package msgchannels

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// EmailConfig carries the email-relay HTTP API settings (SendGrid-compatible).
// APIKey is a secret from the environment, never hardcoded.
type EmailConfig struct {
	APIURL  string // e.g. https://api.sendgrid.com/v3/mail/send
	APIKey  string // env
	From    string
	Timeout time.Duration
}

// Validate enforces a safe configuration.
func (c EmailConfig) Validate() error {
	u, err := url.Parse(c.APIURL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return errors.New("msgchannels: api url must be absolute http(s)")
	}
	if strings.TrimSpace(c.APIKey) == "" {
		return errors.New("msgchannels: api key required (from environment)")
	}
	if !strings.Contains(c.From, "@") {
		return errors.New("msgchannels: from must be an email address")
	}
	if c.Timeout <= 0 {
		return errors.New("msgchannels: timeout required")
	}
	return nil
}

type emailPayload struct {
	Personalizations []struct {
		To []struct {
			Email string `json:"email"`
		} `json:"to"`
	} `json:"personalizations"`
	From struct {
		Email string `json:"email"`
	} `json:"from"`
	Subject string `json:"subject"`
	Content []struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"content"`
}

// SendEmail sends a text/plain email via an HTTP relay (SendGrid-compatible).
func SendEmail(ctx context.Context, cfg EmailConfig, to, subject, body string) error {
	if err := cfg.Validate(); err != nil {
		return err
	}
	if strings.TrimSpace(to) == "" || strings.TrimSpace(subject) == "" || strings.TrimSpace(body) == "" {
		return errors.New("msgchannels: to, subject and body required")
	}
	var p emailPayload
	p.Personalizations = []struct {
		To []struct {
			Email string `json:"email"`
		} `json:"to"`
	}{{To: []struct {
		Email string `json:"email"`
	}{{Email: to}}}}
	p.From.Email = cfg.From
	p.Subject = subject
	p.Content = []struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	}{{Type: "text/plain", Value: body}}

	raw, err := json.Marshal(p)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, cfg.APIURL, bytes.NewReader(raw))
	if err != nil {
		return err
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("authorization", "Bearer "+cfg.APIKey)
	client := &http.Client{Timeout: cfg.Timeout}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		data, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("msgchannels: email status %d: %s", resp.StatusCode, strings.TrimSpace(string(data)))
	}
	return nil
}
