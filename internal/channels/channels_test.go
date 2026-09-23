package channels

import (
	"context"
	"errors"
	"testing"
	"time"
)

type fakeChannel struct {
	code string
	in   []Message
	out  []Message
}

func (f *fakeChannel) Code() string { return f.code }
func (f *fakeChannel) Receive(context.Context) ([]Message, error) {
	m := f.in
	f.in = nil
	return m, nil
}
func (f *fakeChannel) Send(_ context.Context, m Message) error {
	f.out = append(f.out, m)
	return nil
}

func TestMessageValidate(t *testing.T) {
	ok := Message{ChannelCode: "whatsapp", TenantID: "t", ExternalID: "e", ProviderMessageID: "wamid.1", OccurredAt: time.Now(), Direction: DirectionIn, Text: "hola"}
	if err := ok.Validate(); err != nil {
		t.Fatalf("valid message rejected: %v", err)
	}
	bad := ok
	bad.Direction = "sideways"
	if err := bad.Validate(); !errors.Is(err, ErrInvalidMessage) {
		t.Fatalf("bad direction accepted: %v", err)
	}
	empty := ok
	empty.Text = "  "
	if err := empty.Validate(); !errors.Is(err, ErrInvalidMessage) {
		t.Fatalf("empty text accepted: %v", err)
	}
	missingIdentity := ok
	missingIdentity.ProviderMessageID = ""
	if err := missingIdentity.Validate(); !errors.Is(err, ErrInvalidMessage) {
		t.Fatalf("missing provider identity accepted: %v", err)
	}
}

func TestRegistryOnePerCode(t *testing.T) {
	r := NewRegistry()
	if err := r.Register(&fakeChannel{code: "whatsapp"}); err != nil {
		t.Fatal(err)
	}
	if err := r.Register(&fakeChannel{code: "whatsapp"}); !errors.Is(err, ErrDuplicateCode) {
		t.Fatalf("duplicate code accepted: %v", err)
	}
	if err := r.Register(&fakeChannel{code: "BAD!"}); err == nil {
		t.Fatal("invalid code accepted")
	}
	if _, ok := r.For("whatsapp"); !ok {
		t.Fatal("expected whatsapp registered")
	}
	if _, ok := r.For("telegram"); ok {
		t.Fatal("unexpected telegram")
	}
}

func TestDispatcherRoutesAndReplies(t *testing.T) {
	wa := &fakeChannel{code: "whatsapp"}
	r := NewRegistry()
	_ = r.Register(wa)

	d := &Dispatcher{
		Registry: r,
		Respond: func(_ context.Context, message Message) (string, error) {
			if message.ProviderMessageID == "" || message.ThreadID != "thread-1" {
				t.Fatalf("responder lost envelope: %+v", message)
			}
			return "respuesta a: " + message.Text, nil
		},
	}
	in := Message{ChannelCode: "whatsapp", TenantID: "t", ExternalID: "e1", ThreadID: "thread-1", ProviderMessageID: "wamid.1", OccurredAt: time.Now(), Direction: DirectionIn, Text: "hola"}
	if err := d.Handle(context.Background(), in); err != nil {
		t.Fatal(err)
	}
	if len(wa.out) != 1 {
		t.Fatalf("expected 1 outbound, got %d", len(wa.out))
	}
	out := wa.out[0]
	if out.Text != "respuesta a: hola" || out.ExternalID != "e1" || out.Direction != DirectionOut || len(out.DeliveryKey) != 64 {
		t.Fatalf("wrong outbound: %+v", out)
	}
	firstKey := out.DeliveryKey
	if err := d.Handle(context.Background(), in); err != nil || len(wa.out) != 2 || wa.out[1].DeliveryKey != firstKey {
		t.Fatalf("replay changed delivery key: err=%v out=%+v", err, wa.out)
	}
	in.ProviderMessageID = "wamid.2"
	if err := d.Handle(context.Background(), in); err != nil || wa.out[2].DeliveryKey == firstKey {
		t.Fatalf("distinct inbound reused delivery key: err=%v out=%+v", err, wa.out)
	}
}

func TestDispatcherUnknownChannel(t *testing.T) {
	r := NewRegistry()
	d := &Dispatcher{Registry: r, Respond: func(context.Context, Message) (string, error) { return "x", nil }}
	in := Message{ChannelCode: "telegram", TenantID: "t", ExternalID: "e", ProviderMessageID: "msg-1", OccurredAt: time.Now(), Direction: DirectionIn, Text: "hi"}
	if err := d.Handle(context.Background(), in); !errors.Is(err, ErrUnknownChannel) {
		t.Fatalf("expected unknown channel, got %v", err)
	}
}
