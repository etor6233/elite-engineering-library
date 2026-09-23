package outbounddelivery

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"elite.local/enterprise/internal/channels"
)

type memoryStore struct {
	hash, state string
	receipt     Receipt
}

func (s *memoryStore) Claim(_ context.Context, _ channels.Message, hash string) (Claim, error) {
	if s.hash != "" && s.hash != hash {
		return Claim{}, ErrConflict
	}
	if s.state == "accepted" {
		return Claim{Replay: true}, nil
	}
	if s.state == "sending" {
		return Claim{}, ErrInProgress
	}
	if s.state == "unknown" {
		return Claim{}, ErrUnknown
	}
	if s.state == "failed_terminal" {
		return Claim{}, ErrTerminal
	}
	s.hash, s.state = hash, "sending"
	return Claim{}, nil
}
func (s *memoryStore) Complete(_ context.Context, _ channels.Message, hash string, receipt Receipt) error {
	if s.hash != hash || s.state != "sending" {
		return ErrConflict
	}
	s.state, s.receipt = "accepted", receipt
	return nil
}
func (s *memoryStore) MarkUnknown(_ context.Context, _ channels.Message, hash, _ string) error {
	if s.hash != hash || s.state != "sending" {
		return ErrConflict
	}
	s.state = "unknown"
	return nil
}
func (s *memoryStore) MarkFailed(_ context.Context, _ channels.Message, hash, evidence, code string) error {
	if s.hash != hash || s.state != "sending" || !hex64RE.MatchString(evidence) || code == "" {
		return ErrConflict
	}
	s.state = "failed_terminal"
	return nil
}

type testReceiver struct{ code string }

func (r testReceiver) Code() string                                      { return r.code }
func (testReceiver) Receive(context.Context) ([]channels.Message, error) { return nil, nil }
func (testReceiver) Send(context.Context, channels.Message) error {
	return errors.New("direct send must not be used")
}

type testSender struct {
	calls int
	err   error
}

func (s *testSender) SendWithReceipt(context.Context, channels.Message) (Receipt, error) {
	s.calls++
	if s.err != nil {
		return Receipt{}, s.err
	}
	return Receipt{ProviderMessageID: "provider-1", EvidenceSHA256: strings.Repeat("a", 64), AcceptedAt: time.Date(2026, 9, 4, 12, 0, 0, 0, time.UTC)}, nil
}

func outboundMessage() channels.Message {
	return channels.Message{ChannelCode: "whatsapp", TenantID: "tenant", ExternalID: "contact", ThreadID: "thread", DeliveryKey: strings.Repeat("d", 64), Direction: channels.DirectionOut, Text: "respuesta"}
}

func TestChannelSendsOnceAndReplaysWithoutProvider(t *testing.T) {
	store, sender := &memoryStore{}, &testSender{}
	c := &Channel{CodeValue: "whatsapp", Receiver: testReceiver{code: "whatsapp"}, Sender: sender, Store: store}
	message := outboundMessage()
	if err := c.Send(context.Background(), message); err != nil {
		t.Fatal(err)
	}
	if err := c.Send(context.Background(), message); err != nil {
		t.Fatal(err)
	}
	if sender.calls != 1 || store.state != "accepted" {
		t.Fatalf("calls=%d state=%s", sender.calls, store.state)
	}
	message.Text = "divergent"
	if err := c.Send(context.Background(), message); !errors.Is(err, ErrConflict) {
		t.Fatalf("divergent replay=%v", err)
	}
}

func TestAmbiguousProviderFailureNeverRetriesAutomatically(t *testing.T) {
	store, sender := &memoryStore{}, &testSender{err: errors.New("timeout after write")}
	c := &Channel{CodeValue: "whatsapp", Receiver: testReceiver{code: "whatsapp"}, Sender: sender, Store: store}
	message := outboundMessage()
	if err := c.Send(context.Background(), message); !errors.Is(err, ErrUnknown) {
		t.Fatalf("first=%v", err)
	}
	if err := c.Send(context.Background(), message); !errors.Is(err, ErrUnknown) {
		t.Fatalf("retry=%v", err)
	}
	if sender.calls != 1 || store.state != "unknown" {
		t.Fatalf("calls=%d state=%s", sender.calls, store.state)
	}
}

func TestProvenTerminalProviderFailureClosesWithoutUnknown(t *testing.T) {
	store := &memoryStore{}
	sender := &testSender{err: NewTerminalFailure("PROVIDER_REJECTED", strings.Repeat("e", 64), errors.New("invalid request"))}
	c := &Channel{CodeValue: "whatsapp", Receiver: testReceiver{code: "whatsapp"}, Sender: sender, Store: store}
	message := outboundMessage()
	if err := c.Send(context.Background(), message); !errors.Is(err, ErrTerminal) {
		t.Fatalf("first=%v", err)
	}
	if err := c.Send(context.Background(), message); !errors.Is(err, ErrTerminal) {
		t.Fatalf("replay=%v", err)
	}
	if sender.calls != 1 || store.state != "failed_terminal" {
		t.Fatalf("calls=%d state=%s", sender.calls, store.state)
	}
}

func TestMessageHashBindsRecipientAndPayload(t *testing.T) {
	a, err := MessageSHA256(outboundMessage())
	if err != nil {
		t.Fatal(err)
	}
	b := outboundMessage()
	b.ExternalID = "other"
	bh, _ := MessageSHA256(b)
	if a == bh || len(a) != 64 {
		t.Fatalf("a=%s b=%s", a, bh)
	}
}
