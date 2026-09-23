package app

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"elite.local/enterprise/internal/channels"
	"elite.local/enterprise/internal/contactidentity"
	"elite.local/enterprise/internal/leadstream"
	"elite.local/enterprise/internal/outbounddelivery"
	platformpostgres "elite.local/enterprise/internal/platform/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

type appReceiptSender struct{ calls int }

func (s *appReceiptSender) SendWithReceipt(context.Context, channels.Message) (outbounddelivery.Receipt, error) {
	s.calls++
	return outbounddelivery.Receipt{ProviderMessageID: "wamid.outbound.1", EvidenceSHA256: strings.Repeat("f", 64), AcceptedAt: time.Date(2026, 9, 4, 12, 1, 0, 0, time.UTC)}, nil
}

func TestProductionJourneyFencesOutboundReplay(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := leadstream.StableUUID(t.Name(), time.Now().UTC().Format(time.RFC3339Nano))
	if _, err = pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,$2,'Durable E2E','Durable E2E')`, tenant, "durable-"+tenant[:8]); err != nil {
		t.Fatal(err)
	}
	if _, err = pool.Exec(ctx, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,'org','org','Store','store')`, tenant); err != nil {
		t.Fatal(err)
	}
	if _, err = pool.Exec(ctx, `insert into crm.lead(tenant_id,lead_id,organization_id,lifecycle_state,source_code,contact_payload,consent_required) values($1,'lead','org','new','meta','{}',true)`, tenant); err != nil {
		t.Fatal(err)
	}
	key := []byte("0123456789abcdef0123456789abcdef")
	identityStore, err := platformpostgres.NewContactIdentityStore(pool, key)
	if err != nil {
		t.Fatal(err)
	}
	when := time.Date(2026, 9, 4, 12, 0, 0, 0, time.UTC)
	_, err = identityStore.Apply(ctx, contactidentity.Command{TenantID: tenant, ChannelCode: "whatsapp", ExternalID: "wa-durable", LeadID: "lead", SubjectID: "lead:lead", PIIAllowed: true, State: contactidentity.StateActive, PolicyVersion: "sales-contact-v1", EvidenceSHA256: strings.Repeat("a", 64), EffectiveAt: when, RequestID: "durable-binding"})
	if err != nil {
		t.Fatal(err)
	}
	turnStore, err := platformpostgres.NewConversationStore(pool, 5*time.Second, 24*time.Hour, 3)
	if err != nil {
		t.Fatal(err)
	}
	deliveryStore, err := platformpostgres.NewOutboundDeliveryStore(pool, key, 5*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	llmCalls := 0
	llm := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		llmCalls++
		w.Header().Set("content-type", "application/json")
		if llmCalls == 1 {
			_, _ = w.Write([]byte(`{"id":"r1","status":"completed","output":[{"type":"function_call","call_id":"c1","name":"create_quote","arguments":"{\"product\":\"scooter\",\"quantity\":1}"}],"usage":{"input_tokens":1,"output_tokens":1,"total_tokens":2}}`))
			return
		}
		_, _ = w.Write([]byte(`{"id":"r2","status":"completed","output":[{"type":"message","content":[{"type":"output_text","text":"Cotización creada."}]}],"usage":{"input_tokens":1,"output_tokens":1,"total_tokens":2}}`))
	}))
	defer llm.Close()
	domainCalls := 0
	domain := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		domainCalls++
		w.Header().Set("content-type", "application/json")
		_, _ = w.Write([]byte(`{"quotation_id":"q1"}`))
	}))
	defer domain.Close()
	receiver := &appChannel{}
	sender := &appReceiptSender{}
	durable := &outbounddelivery.Channel{CodeValue: "whatsapp", Receiver: receiver, Sender: sender, Store: deliveryStore}
	registry := channels.NewRegistry()
	if err = registry.Register(durable); err != nil {
		t.Fatal(err)
	}
	cfg := fullConfig()
	cfg.TenantID, cfg.LLMBaseURL, cfg.DomainBaseURL = tenant, llm.URL, domain.URL
	cfg.ChannelRegistry, cfg.ConversationStore, cfg.ContactResolver = registry, turnStore, identityStore
	a, err := NewProduction(cfg)
	if err != nil {
		t.Fatal(err)
	}
	message := channels.Message{ChannelCode: "whatsapp", TenantID: tenant, ExternalID: "wa-durable", ThreadID: "thread", ProviderMessageID: "wamid.inbound.1", OccurredAt: when, Direction: channels.DirectionIn, Text: "Cotizame un scooter"}
	if err = a.Dispatcher.Handle(ctx, message); err != nil {
		t.Fatal(err)
	}
	if err = a.Dispatcher.Handle(ctx, message); err != nil {
		t.Fatal(err)
	}
	var state string
	var attempts, events int
	if err = pool.QueryRow(ctx, `select d.state,d.attempt_count,count(e.*) from communication.outbound_delivery d join communication.outbound_delivery_event e using(tenant_id,channel_code,delivery_key) where d.tenant_id=$1 group by d.state,d.attempt_count`, tenant).Scan(&state, &attempts, &events); err != nil {
		t.Fatal(err)
	}
	if state != "accepted" || attempts != 1 || events != 2 || sender.calls != 1 || llmCalls != 2 || domainCalls != 1 {
		t.Fatalf("state=%s attempts=%d events=%d sends=%d llm=%d domain=%d", state, attempts, events, sender.calls, llmCalls, domainCalls)
	}
}

func TestNewProductionRejectsUnfencedChannel(t *testing.T) {
	cfg := fullConfig()
	if _, err := NewProduction(cfg); err == nil {
		t.Fatal("unfenced channel accepted")
	}
}
