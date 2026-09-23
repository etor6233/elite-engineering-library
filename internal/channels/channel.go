// Package channels provides the channel abstraction for the conversational
// agent: a Channel contract (receive/send), a registry keyed by channel code,
// and a dispatcher that routes inbound messages to a responder and sends the
// reply back. The real adapters (Meta Graph, Twilio, SMTP) are CONDITIONED.
package channels

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// Direction is inbound or outbound.
type Direction string

const (
	DirectionIn  Direction = "in"
	DirectionOut Direction = "out"
)

var (
	codeRe = regexp.MustCompile(`^[a-z][a-z0-9_]{0,31}$`)

	ErrInvalidMessage = errors.New("channels: invalid message")
	ErrUnknownChannel = errors.New("channels: unknown channel")
	ErrDuplicateCode  = errors.New("channels: duplicate channel code")
)

// Message is a channel-agnostic inbound/outbound message.
type Message struct {
	ChannelCode       string
	TenantID          string
	ExternalID        string // contact/recipient id in the channel
	ThreadID          string // conversation thread (optional)
	ProviderMessageID string // provider-owned inbound identity; never LLM-authored
	DeliveryKey       string // caller-owned outbound deduplication identity
	OccurredAt        time.Time
	Direction         Direction
	Text              string
}

// Validate enforces the message contract.
func (m Message) Validate() error {
	if !codeRe.MatchString(m.ChannelCode) {
		return fmt.Errorf("%w: channel code", ErrInvalidMessage)
	}
	if strings.TrimSpace(m.TenantID) == "" || len(m.TenantID) > 64 {
		return fmt.Errorf("%w: tenant", ErrInvalidMessage)
	}
	if strings.TrimSpace(m.ExternalID) == "" || len(m.ExternalID) > 256 {
		return fmt.Errorf("%w: external id", ErrInvalidMessage)
	}
	switch m.Direction {
	case DirectionIn:
		if strings.TrimSpace(m.ProviderMessageID) == "" || len(m.ProviderMessageID) > 256 || m.OccurredAt.IsZero() {
			return fmt.Errorf("%w: inbound identity", ErrInvalidMessage)
		}
	case DirectionOut:
		if len(strings.TrimSpace(m.DeliveryKey)) < 16 || len(m.DeliveryKey) > 128 {
			return fmt.Errorf("%w: outbound delivery key", ErrInvalidMessage)
		}
	default:
		return fmt.Errorf("%w: direction", ErrInvalidMessage)
	}
	if strings.TrimSpace(m.Text) == "" {
		return fmt.Errorf("%w: empty text", ErrInvalidMessage)
	}
	return nil
}

// Channel is the receive/send boundary for one messaging surface.
type Channel interface {
	Code() string
	Receive(ctx context.Context) ([]Message, error)
	Send(ctx context.Context, m Message) error
}
