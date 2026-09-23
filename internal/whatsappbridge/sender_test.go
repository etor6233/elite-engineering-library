package whatsappbridge

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"elite.local/enterprise/internal/channels"
	"elite.local/enterprise/internal/leadstream"
	"elite.local/enterprise/internal/outbounddelivery"
	"elite.local/enterprise/internal/platform/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

type approvalFixture struct{ value Approval }

func (a approvalFixture) ResolveWhatsAppApproval(context.Context, string, string) (Approval, error) {
	return a.value, nil
}

type tokenFixture struct{}

func (tokenFixture) WhatsAppToken(context.Context) (string, error) { return "synthetic-token", nil }

type receiverFixture struct{}

func (receiverFixture) Code() string                                        { return "whatsapp" }
func (receiverFixture) Receive(context.Context) ([]channels.Message, error) { return nil, nil }
func (receiverFixture) Send(context.Context, channels.Message) error {
	return errors.New("receiver cannot send")
}

type storeFixture struct {
	hash, state string
	receipt     outbounddelivery.Receipt
}

func (s *storeFixture) Claim(_ context.Context, _ channels.Message, h string) (outbounddelivery.Claim, error) {
	if s.hash != "" && s.hash != h {
		return outbounddelivery.Claim{}, outbounddelivery.ErrConflict
	}
	switch s.state {
	case "accepted":
		return outbounddelivery.Claim{Replay: true}, nil
	case "unknown":
		return outbounddelivery.Claim{}, outbounddelivery.ErrUnknown
	case "sending":
		return outbounddelivery.Claim{}, outbounddelivery.ErrInProgress
	}
	s.hash = h
	s.state = "sending"
	return outbounddelivery.Claim{}, nil
}
func (s *storeFixture) Complete(_ context.Context, _ channels.Message, _ string, r outbounddelivery.Receipt) error {
	s.state = "accepted"
	s.receipt = r
	return nil
}
func (s *storeFixture) MarkUnknown(context.Context, channels.Message, string, string) error {
	s.state = "unknown"
	return nil
}
func (s *storeFixture) MarkFailed(context.Context, channels.Message, string, string, string) error {
	s.state = "failed_terminal"
	return nil
}

func messageFixture() channels.Message {
	return channels.Message{ChannelCode: "whatsapp", TenantID: "tenant-fixture", ExternalID: "5491112345678", Direction: channels.DirectionOut, DeliveryKey: "appointment-confirmed-fixture-key", Text: `{"recipient":"5491112345678","template_name":"appointment_confirmed","language_code":"es_AR","body_parameters":["2026-09-07T10:00Z"]}`}
}

func senderFixture(t *testing.T, mode string) (*Sender, channels.Message, string) {
	t.Helper()
	python := os.Getenv("ELITE_WHATSAPP_PYTHON")
	if python == "" {
		t.Skip("process contract requires explicit local ELITE_WHATSAPP_PYTHON; never a live provider")
	}
	pythonBytes, err := os.ReadFile(python)
	if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	marker := filepath.Join(root, "invocations.txt")
	markerJSON, _ := json.Marshal(marker)
	// A clearly synthetic child tests the real process boundary, not Meta.
	script := `import sys,json,datetime,pathlib
f=json.load(sys.stdin)
p=pathlib.Path(` + string(markerJSON) + `)
with p.open('ab') as stream: stream.write(b'call\n')
mode=` + `"` + mode + `"` + `
if mode=='lost': sys.stderr.write('synthetic-token private recipient'); sys.exit(2)
r={'schema':'elite-whatsapp-send-result/v1','binding_sha256':f['binding_sha256'],'provider_message_id':'wamid.fixture','evidence_sha256':'a'*64,'accepted_at':datetime.datetime.now(datetime.timezone.utc).isoformat()}
if mode=='wrong': r['binding_sha256']='f'*64
if mode=='missing': r['provider_message_id']=''
if mode=='overflow': print('x'*8192); sys.exit(0)
print(json.dumps(r))
`
	if err = os.WriteFile(filepath.Join(root, "whatsapp_cloud.py"), []byte(script), 0600); err != nil {
		t.Fatal(err)
	}
	message := messageFixture()
	h, err := outbounddelivery.MessageSHA256(message)
	if err != nil {
		t.Fatal(err)
	}
	profile := json.RawMessage(`{"synthetic_fixture":true}`)
	approval := Approval{MessageSHA256: h, ProfileSHA256: digest(profile), EvidenceSHA256: strings.Repeat("b", 64), ApprovedBy: "fixture-owner", NotBefore: time.Now().Add(-time.Minute), ExpiresAt: time.Now().Add(time.Minute)}
	s := &Sender{TenantID: message.TenantID, Profile: profile, Approvals: approvalFixture{approval}, Tokens: tokenFixture{}, Process: Process{PythonExecutable: python, PythonSHA256: digest(pythonBytes), AdapterDirectory: root, AdapterSHA256: digest([]byte(script)), EvidenceDirectory: filepath.Join(root, "evidence")}}
	return s, message, marker
}

func TestProcessBridgeFencedReplayAndAmbiguity(t *testing.T) {
	for _, mode := range []string{"accepted", "lost", "wrong", "missing", "overflow"} {
		t.Run(mode, func(t *testing.T) {
			sender, message, marker := senderFixture(t, mode)
			store := &storeFixture{}
			channel := &outbounddelivery.Channel{CodeValue: "whatsapp", Receiver: receiverFixture{}, Sender: sender, Store: store}
			for i := 0; i < 2; i++ {
				err := channel.Send(context.Background(), message)
				if mode == "accepted" && err != nil {
					t.Fatal(err)
				}
				if mode != "accepted" && !errors.Is(err, outbounddelivery.ErrUnknown) {
					t.Fatalf("expected unknown, got %v", err)
				}
				if err != nil && strings.Contains(err.Error(), "synthetic-token") {
					t.Fatal("secret in error")
				}
			}
			calls, err := os.ReadFile(marker)
			if err != nil || string(calls) != "call\n" {
				t.Fatalf("expected one process invocation: %q %v", calls, err)
			}
			if mode == "accepted" && store.state != "accepted" {
				t.Fatal(store.state)
			}
			if mode != "accepted" && store.state != "unknown" {
				t.Fatal(store.state)
			}
		})
	}
}

func TestBridgeRejectsUnapprovedOrDriftedRequestBeforeChild(t *testing.T) {
	for _, mode := range []string{"tenant", "recipient", "approval", "expired", "profile", "python", "adapter", "no_approvals", "channel"} {
		t.Run(mode, func(t *testing.T) {
			sender, message, marker := senderFixture(t, "accepted")
			switch mode {
			case "tenant":
				message.TenantID = "other"
			case "recipient":
				message.ExternalID = "5491199999999"
			case "approval":
				a := sender.Approvals.(approvalFixture)
				a.value.MessageSHA256 = strings.Repeat("c", 64)
				sender.Approvals = a
			case "expired":
				a := sender.Approvals.(approvalFixture)
				a.value.ExpiresAt = time.Now().Add(-time.Second)
				sender.Approvals = a
			case "profile":
				sender.Profile = json.RawMessage(`{"changed":true}`)
			case "python":
				sender.Process.PythonSHA256 = strings.Repeat("c", 64)
			case "adapter":
				sender.Process.AdapterSHA256 = strings.Repeat("c", 64)
			case "no_approvals":
				sender.Approvals = nil
			case "channel":
				message.ChannelCode = "email"
			}
			if _, err := sender.SendWithReceipt(context.Background(), message); !errors.Is(err, ErrBridge) {
				t.Fatal("unapproved call accepted", err)
			}
			if _, err := os.Stat(marker); !os.IsNotExist(err) {
				t.Fatal("child invoked")
			}
		})
	}
}

func TestRequestContractIsExact(t *testing.T) {
	for _, text := range []string{`{}`, messageFixture().Text + ` {}`, strings.Replace(messageFixture().Text, `"recipient":`, `"extra":1,"recipient":`, 1), strings.Repeat("x", 32769)} {
		message := messageFixture()
		message.Text = text
		if _, err := requestFor(message); err == nil {
			t.Fatal("malformed request accepted")
		}
	}
}

func TestProcessBridgePostgresDurableReplay(t *testing.T) {
	raw := os.Getenv("ELITE_WHATSAPP_TEST_DATABASE_URL")
	if raw == "" {
		t.Skip("requires dedicated local elite_whatsapp_ database")
	}
	endpoint, err := url.Parse(raw)
	if err != nil || endpoint.Hostname() != "127.0.0.1" || !strings.HasPrefix(endpoint.Path, "/elite_whatsapp_") {
		t.Fatal("unsafe test database target")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()
	pool, err := pgxpool.New(ctx, raw)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	for _, mode := range []string{"accepted", "lost", "wrong", "missing", "overflow"} {
		t.Run(mode, func(t *testing.T) {
			sender, message, marker := senderFixture(t, mode)
			tenant := leadstream.StableUUID(t.Name(), time.Now().UTC().Format(time.RFC3339Nano))
			_, err := pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,$2,'Synthetic bridge','Synthetic bridge')`, tenant, "wa-"+tenant[:8])
			if err != nil {
				t.Fatal(err)
			}
			sender.TenantID = tenant
			message.TenantID = tenant
			binding, err := outbounddelivery.MessageSHA256(message)
			if err != nil {
				t.Fatal(err)
			}
			approved := sender.Approvals.(approvalFixture)
			approved.value.MessageSHA256 = binding
			sender.Approvals = approved
			for attempt := 0; attempt < 2; attempt++ {
				// A new Store instance demonstrates that replay/unknown is in PostgreSQL,
				// not in the Go fixture's memory.
				store, err := postgres.NewOutboundDeliveryStore(pool, []byte("0123456789abcdef0123456789abcdef"), time.Second)
				if err != nil {
					t.Fatal(err)
				}
				channel := &outbounddelivery.Channel{CodeValue: "whatsapp", Receiver: receiverFixture{}, Sender: sender, Store: store}
				err = channel.Send(ctx, message)
				if mode == "accepted" && err != nil {
					t.Fatal(err)
				}
				if mode != "accepted" && !errors.Is(err, outbounddelivery.ErrUnknown) {
					t.Fatalf("expected durable unknown: %v", err)
				}
				changed := message
				changed.Text += " "
				if err = channel.Send(ctx, changed); !errors.Is(err, outbounddelivery.ErrConflict) {
					t.Fatalf("divergent request: %v", err)
				}
			}
			calls, err := os.ReadFile(marker)
			if err != nil || string(calls) != "call\n" {
				t.Fatalf("provider invocation count: %q %v", calls, err)
			}
			var state string
			var attempts, events int
			err = pool.QueryRow(ctx, `select state,attempt_count,(select count(*) from communication.outbound_delivery_event e where e.tenant_id=d.tenant_id and e.channel_code=d.channel_code and e.delivery_key=d.delivery_key) from communication.outbound_delivery d where tenant_id=$1 and channel_code='whatsapp' and delivery_key=$2`, tenant, message.DeliveryKey).Scan(&state, &attempts, &events)
			expected := "unknown"
			if mode == "accepted" {
				expected = "accepted"
			}
			if err != nil || state != expected || attempts != 1 || events != 2 {
				t.Fatalf("durable evidence state=%s attempts=%d events=%d err=%v", state, attempts, events, err)
			}
		})
	}
}
