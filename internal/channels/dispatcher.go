package channels

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
)

// Responder handles the complete inbound envelope so command identity, scope
// and correlation are never lost before a side effect.
type Responder func(ctx context.Context, message Message) (string, error)

// Dispatcher routes inbound messages to a Responder and sends the reply back
// through the same channel. It never widens the tenant or channel scope.
type Dispatcher struct {
	Registry *Registry
	Respond  Responder
}

// Handle processes one inbound message and sends the reply via its channel.
func (d *Dispatcher) Handle(ctx context.Context, m Message) error {
	if err := m.Validate(); err != nil {
		return err
	}
	if m.Direction != DirectionIn {
		return ErrInvalidMessage
	}
	if d.Registry == nil {
		return ErrUnknownChannel
	}
	ch, ok := d.Registry.For(m.ChannelCode)
	if !ok {
		return ErrUnknownChannel
	}
	if d.Respond == nil {
		return ErrUnknownChannel
	}
	reply, err := d.Respond(ctx, m)
	if err != nil {
		return err
	}
	return ch.Send(ctx, Message{
		ChannelCode: m.ChannelCode,
		TenantID:    m.TenantID,
		ExternalID:  m.ExternalID,
		ThreadID:    m.ThreadID,
		DeliveryKey: deliveryKey(m),
		Direction:   DirectionOut,
		Text:        reply,
	})
}

func deliveryKey(m Message) string {
	sum := sha256.Sum256([]byte(m.TenantID + "\x00" + m.ChannelCode + "\x00" + m.ProviderMessageID + "\x00reply"))
	return hex.EncodeToString(sum[:])
}
