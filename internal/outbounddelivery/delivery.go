// Package outbounddelivery fences provider sends so a lost acknowledgement
// never causes an automatic duplicate. Ambiguous outcomes require explicit
// reconciliation before any further effect.
package outbounddelivery

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"regexp"
	"strings"
	"time"

	"elite.local/enterprise/internal/channels"
)

var (
	ErrInvalid     = errors.New("outbounddelivery: invalid configuration or message")
	ErrConflict    = errors.New("outbounddelivery: payload conflict")
	ErrInProgress  = errors.New("outbounddelivery: delivery already in progress")
	ErrUnknown     = errors.New("outbounddelivery: provider effect unknown; reconcile before retry")
	ErrTerminal    = errors.New("outbounddelivery: delivery terminal")
	hex64RE        = regexp.MustCompile(`^[0-9a-f]{64}$`)
	terminalCodeRE = regexp.MustCompile(`^[A-Z][A-Z0-9_]{2,63}$`)
)

type Claim struct{ Replay bool }

type Receipt struct {
	ProviderMessageID string
	EvidenceSHA256    string
	AcceptedAt        time.Time
}

type TerminalFailure struct {
	Code           string
	EvidenceSHA256 string
	Cause          error
}

func (e *TerminalFailure) Error() string {
	if e == nil {
		return ErrTerminal.Error()
	}
	return ErrTerminal.Error() + ": " + e.Code
}

func (e *TerminalFailure) Unwrap() error { return e.Cause }

func NewTerminalFailure(code, evidenceSHA256 string, cause error) error {
	if !terminalCodeRE.MatchString(code) || !hex64RE.MatchString(evidenceSHA256) {
		return ErrInvalid
	}
	return &TerminalFailure{Code: code, EvidenceSHA256: evidenceSHA256, Cause: cause}
}

func (r Receipt) Validate() error {
	if strings.TrimSpace(r.ProviderMessageID) == "" || len(r.ProviderMessageID) > 256 || !hex64RE.MatchString(r.EvidenceSHA256) || r.AcceptedAt.IsZero() {
		return ErrInvalid
	}
	return nil
}

type Sender interface {
	SendWithReceipt(context.Context, channels.Message) (Receipt, error)
}

type Store interface {
	Claim(context.Context, channels.Message, string) (Claim, error)
	Complete(context.Context, channels.Message, string, Receipt) error
	MarkUnknown(context.Context, channels.Message, string, string) error
	MarkFailed(context.Context, channels.Message, string, string, string) error
}

type Channel struct {
	CodeValue string
	Receiver  channels.Channel
	Sender    Sender
	Store     Store
}

func (c *Channel) Code() string { return c.CodeValue }

func (c *Channel) DurableDelivery() bool {
	return c != nil && c.CodeValue != "" && c.Receiver != nil && c.Sender != nil && c.Store != nil
}

func (c *Channel) Receive(ctx context.Context) ([]channels.Message, error) {
	if !c.DurableDelivery() || c.Receiver.Code() != c.CodeValue {
		return nil, ErrInvalid
	}
	return c.Receiver.Receive(ctx)
}

func (c *Channel) Send(ctx context.Context, message channels.Message) error {
	if !c.DurableDelivery() || c.Receiver.Code() != c.CodeValue || message.ChannelCode != c.CodeValue || message.Direction != channels.DirectionOut {
		return ErrInvalid
	}
	if err := message.Validate(); err != nil {
		return err
	}
	hash, err := MessageSHA256(message)
	if err != nil {
		return err
	}
	claim, err := c.Store.Claim(ctx, message, hash)
	if err != nil || claim.Replay {
		return err
	}
	receipt, sendErr := c.Sender.SendWithReceipt(ctx, message)
	if sendErr != nil {
		var terminal *TerminalFailure
		if errors.As(sendErr, &terminal) {
			markErr := c.Store.MarkFailed(ctx, message, hash, terminal.EvidenceSHA256, terminal.Code)
			return errors.Join(ErrTerminal, sendErr, markErr)
		}
		markErr := c.Store.MarkUnknown(ctx, message, hash, "PROVIDER_CALL_UNCERTAIN")
		return errors.Join(ErrUnknown, sendErr, markErr)
	}
	if err = receipt.Validate(); err != nil {
		markErr := c.Store.MarkUnknown(ctx, message, hash, "PROVIDER_RECEIPT_INVALID")
		return errors.Join(ErrUnknown, err, markErr)
	}
	if err = c.Store.Complete(ctx, message, hash, receipt); err != nil {
		return errors.Join(ErrUnknown, err)
	}
	return nil
}

func MessageSHA256(message channels.Message) (string, error) {
	if err := message.Validate(); err != nil || message.Direction != channels.DirectionOut {
		return "", ErrInvalid
	}
	canonical := struct {
		ChannelCode, TenantID, ExternalID, ThreadID, DeliveryKey, Text string
	}{message.ChannelCode, message.TenantID, message.ExternalID, message.ThreadID, message.DeliveryKey, message.Text}
	raw, err := json.Marshal(canonical)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(raw)
	return hex.EncodeToString(sum[:]), nil
}
