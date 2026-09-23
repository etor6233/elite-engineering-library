// Package msgchannels provides a WhatsApp text-send leaf, typed inbound text
// parsers for WhatsApp/Page/Instagram, HMAC verification and SMTP. It is not a
// durable multichannel runtime; consent, auth and reconciliation belong to the
// connected owners. No Instagram/Messenger outbound claim is made here.
package msgchannels

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// VerifyWebhookSignature checks the Meta X-Hub-Signature-256 header against the
// raw body using HMAC-SHA256 with the app secret (constant-time compare).
func VerifyWebhookSignature(appSecret string, rawBody []byte, header string) bool {
	if appSecret == "" || len(rawBody) == 0 || header == "" {
		return false
	}
	parts := strings.SplitN(strings.TrimSpace(header), "=", 2)
	if len(parts) != 2 || parts[0] != "sha256" {
		return false
	}
	got, err := hex.DecodeString(parts[1])
	if err != nil {
		return false
	}
	mac := hmac.New(sha256.New, []byte(appSecret))
	mac.Write(rawBody)
	return hmac.Equal(got, mac.Sum(nil))
}

// GraphConfig carries the Meta Graph API settings. AccessToken is a secret from
// the environment, never hardcoded.
type GraphConfig struct {
	GraphURL    string // https://graph.facebook.com
	Version     string // e.g. v22.0
	SenderID    string // WhatsApp phone_number_id only
	AccessToken string // env
	Timeout     time.Duration
}

// Validate enforces a safe, exact configuration.
func (c GraphConfig) Validate() error {
	u, err := url.Parse(c.GraphURL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return errors.New("msgchannels: graph url must be absolute http(s)")
	}
	if strings.TrimSpace(c.Version) == "" || strings.TrimSpace(c.SenderID) == "" {
		return errors.New("msgchannels: version and sender id required")
	}
	if strings.TrimSpace(c.AccessToken) == "" {
		return errors.New("msgchannels: access token required (from environment)")
	}
	if c.Timeout <= 0 {
		return errors.New("msgchannels: timeout required")
	}
	return nil
}

// SendText sends WhatsApp text only; callers own consent and delivery receipts.
func SendText(ctx context.Context, cfg GraphConfig, to, text string) error {
	if err := cfg.Validate(); err != nil {
		return err
	}
	if strings.TrimSpace(to) == "" || strings.TrimSpace(text) == "" {
		return errors.New("msgchannels: recipient and text required")
	}
	payload := map[string]any{
		"messaging_product": "whatsapp",
		"to":                to,
		"type":              "text",
		"text":              map[string]string{"body": text},
	}
	body, _ := json.Marshal(payload)
	endpoint := strings.TrimRight(cfg.GraphURL, "/") + "/" + cfg.Version + "/" + cfg.SenderID + "/messages"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("authorization", "Bearer "+cfg.AccessToken)
	client := &http.Client{Timeout: cfg.Timeout}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("msgchannels: graph status %d", resp.StatusCode)
	}
	return nil
}

// InboundMessage is text decoded from a typed webhook. Parsing is not trust.
type InboundMessage struct{ From, Channel, Text, MessageID, RecipientID string }

func ParseWebhook(raw []byte) ([]InboundMessage, error) {
	if len(raw) == 0 || len(raw) > 1<<20 {
		return nil, errors.New("msgchannels: payload limit")
	}
	var root struct {
		Object string `json:"object"`
		Entry  []struct {
			Messaging []struct {
				Sender struct {
					ID string `json:"id"`
				} `json:"sender"`
				Recipient struct {
					ID string `json:"id"`
				} `json:"recipient"`
				Message struct {
					Text string `json:"text"`
					MID  string `json:"mid"`
					Echo bool   `json:"is_echo"`
				} `json:"message"`
			} `json:"messaging"`
			Changes []struct {
				Field string `json:"field"`
				Value struct {
					Product  string `json:"messaging_product"`
					Metadata struct {
						Phone string `json:"phone_number_id"`
					} `json:"metadata"`
					Messages []struct {
						From string `json:"from"`
						ID   string `json:"id"`
						Type string `json:"type"`
						Text struct {
							Body string `json:"body"`
						} `json:"text"`
					} `json:"messages"`
				} `json:"value"`
			} `json:"changes"`
		} `json:"entry"`
	}
	if err := json.Unmarshal(raw, &root); err != nil {
		return nil, err
	}
	if root.Object != "whatsapp_business_account" && root.Object != "page" && root.Object != "instagram" {
		return nil, errors.New("msgchannels: unsupported envelope")
	}
	out := []InboundMessage{}
	for _, e := range root.Entry {
		if root.Object == "whatsapp_business_account" {
			for _, c := range e.Changes {
				if c.Field != "messages" {
					continue
				}
				if c.Value.Product != "whatsapp" || c.Value.Metadata.Phone == "" {
					return nil, errors.New("msgchannels: invalid WhatsApp envelope")
				}
				for _, m := range c.Value.Messages {
					if m.Type != "text" {
						continue
					}
					if m.From == "" || m.ID == "" || m.Text.Body == "" {
						return nil, errors.New("msgchannels: invalid WhatsApp text")
					}
					out = append(out, InboundMessage{m.From, "whatsapp", m.Text.Body, m.ID, c.Value.Metadata.Phone})
				}
			}
		} else {
			for _, m := range e.Messaging {
				if m.Message.Text == "" || m.Message.Echo {
					continue
				}
				if m.Sender.ID == "" || m.Recipient.ID == "" || m.Message.MID == "" {
					return nil, errors.New("msgchannels: invalid page/Instagram text")
				}
				out = append(out, InboundMessage{m.Sender.ID, root.Object, m.Message.Text, m.Message.MID, m.Recipient.ID})
			}
		}
	}
	return out, nil
}

// ParseVerifiedWebhook verifies bytes before parsing and binds recipient+channel
// to server configuration. The caller then resolves tenant/contact and dedupes
// MessageID through the existing durable inbound owner before any action.
func ParseVerifiedWebhook(secret string, raw []byte, signature, channel, recipient string) ([]InboundMessage, error) {
	if recipient == "" || channel == "" || !VerifyWebhookSignature(secret, raw, signature) {
		return nil, errors.New("msgchannels: webhook denied")
	}
	rows, err := ParseWebhook(raw)
	if err != nil {
		return nil, err
	}
	for _, m := range rows {
		if m.Channel != channel || m.RecipientID != recipient {
			return nil, errors.New("msgchannels: recipient/channel mismatch")
		}
	}
	return rows, nil
}
