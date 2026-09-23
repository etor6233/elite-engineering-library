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
	platformpostgres "elite.local/enterprise/internal/platform/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestCapturedLeadIdentityReachesConversationAndQuoteExactlyOnce(t *testing.T) {
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
	if _, err = pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,$2,'Contact E2E','Contact E2E')`, tenant, "contact-e2e-"+tenant[:8]); err != nil {
		t.Fatal(err)
	}
	if _, err = pool.Exec(ctx, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,'org','org','Contact Store','store')`, tenant); err != nil {
		t.Fatal(err)
	}
	if _, err = pool.Exec(ctx, `insert into crm.lead(tenant_id,lead_id,organization_id,lifecycle_state,source_code,contact_payload,consent_required) values($1,'lead','org','new','meta','{}',true)`, tenant); err != nil {
		t.Fatal(err)
	}
	identityStore, err := platformpostgres.NewContactIdentityStore(pool, []byte("0123456789abcdef0123456789abcdef"))
	if err != nil {
		t.Fatal(err)
	}
	when := time.Date(2026, 9, 4, 12, 0, 0, 0, time.UTC)
	_, err = identityStore.Apply(ctx, contactidentity.Command{TenantID: tenant, ChannelCode: "whatsapp", ExternalID: "wa-captured-lead", LeadID: "lead", SubjectID: "lead:lead", PIIAllowed: true, State: contactidentity.StateActive, PolicyVersion: "sales-contact-v1", EvidenceSHA256: strings.Repeat("d", 64), EffectiveAt: when.Add(-time.Minute), RequestID: "contact-e2e-binding-1"})
	if err != nil {
		t.Fatal(err)
	}
	turnStore, err := platformpostgres.NewConversationStore(pool, 5*time.Second, 24*time.Hour, 3)
	if err != nil {
		t.Fatal(err)
	}
	llmCalls := 0
	llm := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		llmCalls++
		w.Header().Set("content-type", "application/json")
		if llmCalls == 1 {
			_, _ = w.Write([]byte(`{"id":"resp-contact-1","status":"completed","output":[{"type":"function_call","call_id":"call-contact-1","name":"create_quote","arguments":"{\"product\":\"scooter\",\"quantity\":1}"}],"usage":{"input_tokens":12,"output_tokens":4,"total_tokens":16}}`))
			return
		}
		_, _ = w.Write([]byte(`{"id":"resp-contact-2","status":"completed","output":[{"type":"message","content":[{"type":"output_text","text":"Cotización creada."}]}],"usage":{"input_tokens":8,"output_tokens":5,"total_tokens":13}}`))
	}))
	defer llm.Close()
	domainCalls := 0
	domain := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		domainCalls++
		if r.URL.Path != "/v1/franchise/quotes" || r.Header.Get("Authorization") != "Bearer token" || len(r.Header.Get("Idempotency-Key")) != 64 {
			t.Errorf("domain request path=%q auth=%q key=%q", r.URL.Path, r.Header.Get("Authorization"), r.Header.Get("Idempotency-Key"))
		}
		w.Header().Set("content-type", "application/json")
		_, _ = w.Write([]byte(`{"quotation_id":"quote-contact-1"}`))
	}))
	defer domain.Close()
	channel := &appChannel{}
	registry := channels.NewRegistry()
	if err = registry.Register(channel); err != nil {
		t.Fatal(err)
	}
	cfg := fullConfig()
	cfg.TenantID, cfg.LLMBaseURL, cfg.DomainBaseURL = tenant, llm.URL, domain.URL
	cfg.ChannelRegistry, cfg.ConversationStore, cfg.ContactResolver = registry, turnStore, identityStore
	a, err := New(cfg)
	if err != nil {
		t.Fatal(err)
	}
	message := channels.Message{ChannelCode: "whatsapp", TenantID: tenant, ExternalID: "wa-captured-lead", ThreadID: "thread-contact", ProviderMessageID: "wamid.contact.1", OccurredAt: when, Direction: channels.DirectionIn, Text: "Cotizame un scooter"}
	if err = a.Dispatcher.Handle(ctx, message); err != nil {
		t.Fatal(err)
	}
	if err = a.Dispatcher.Handle(ctx, message); err != nil {
		t.Fatal(err)
	}
	var state, toolName string
	var attempts int
	if err = pool.QueryRow(ctx, `select state,tool_name,attempt_count from communication.conversation_turn where tenant_id=$1 and channel_code=$2 and provider_message_id=$3`, tenant, message.ChannelCode, message.ProviderMessageID).Scan(&state, &toolName, &attempts); err != nil {
		t.Fatal(err)
	}
	if state != "completed" || toolName != "create_quote" || attempts != 1 || llmCalls != 2 || domainCalls != 1 || len(channel.sent) != 2 {
		t.Fatalf("state=%s tool=%s attempts=%d llm=%d domain=%d sent=%d", state, toolName, attempts, llmCalls, domainCalls, len(channel.sent))
	}
}
