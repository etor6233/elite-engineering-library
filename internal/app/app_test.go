package app

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"

	"elite.local/enterprise/internal/approval"
	"elite.local/enterprise/internal/channels"
	"elite.local/enterprise/internal/conversationruntime"
	"elite.local/enterprise/internal/domainbind"
	platformpostgres "elite.local/enterprise/internal/platform/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

type appTokenProvider struct{}

func (appTokenProvider) AccessToken(context.Context) (string, error) { return "token", nil }

type appStore struct {
	mu       sync.Mutex
	hash     string
	response string
}

func (s *appStore) Claim(_ context.Context, _ channels.Message, hash string) (conversationruntime.Claim, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.response != "" {
		return conversationruntime.Claim{Replay: true, State: conversationruntime.StateCompleted, ResponseText: s.response}, nil
	}
	s.hash = hash
	return conversationruntime.Claim{Generation: 1}, nil
}
func (*appStore) History(context.Context, channels.Message, int) ([]conversationruntime.HistoryMessage, error) {
	return nil, nil
}
func (s *appStore) Complete(_ context.Context, _ channels.Message, hash string, generation int, completion conversationruntime.Completion) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if hash != s.hash {
		return conversationruntime.ErrReplayConflict
	}
	s.response = completion.ResponseText
	return nil
}
func (*appStore) FailRetryable(context.Context, channels.Message, string, int, string) error {
	return nil
}

type appResolver struct{}

func (appResolver) Resolve(context.Context, channels.Message) (conversationruntime.Contact, error) {
	return conversationruntime.Contact{Scope: domainbind.Scope{OrganizationID: "north", LeadID: "lead-1"}, SubjectID: "subject-1", PIIAllowed: true}, nil
}

type appChannel struct {
	mu   sync.Mutex
	sent []channels.Message
}

func (*appChannel) Code() string                                        { return "whatsapp" }
func (*appChannel) Receive(context.Context) ([]channels.Message, error) { return nil, nil }
func (c *appChannel) Send(_ context.Context, message channels.Message) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.sent = append(c.sent, message)
	return nil
}

func fullConfig() Config {
	registry := channels.NewRegistry()
	_ = registry.Register(&appChannel{})
	return Config{
		LLMAPIKey:           "sk-test",
		LLMModel:            "gpt-4.1-mini",
		LLMBaseURL:          "https://api.openai.com/v1",
		DomainBaseURL:       "http://localhost:8080",
		DomainTokenProvider: appTokenProvider{},
		TenantID:            "018f4d4a-7b36-7a21-8d10-2f4c54c20001",
		TenantCode:          "acme",
		ServiceKinds:        map[string]string{"corte": "service"},
		ProductVariants:     map[string]string{"scooter": "variant-1"},
		PriceBookID:         "retail",
		ConversationStore:   &appStore{},
		ContactResolver:     appResolver{},
		ChannelRegistry:     registry,
		Conversation: conversationruntime.Config{
			Instructions: "Usa sólo las herramientas ofrecidas.", PromptCacheKey: "franchise-v1", MaxOutputTokens: 128,
			MaxHistory: 4, MaxInputBytes: 4096, MaxToolBytes: 4096, StoreApproved: true, HandoffText: "Escalado a un especialista.",
		},
		LLMTokenBudget: 100000,
		ApprovalPolicy: approval.Policy{AutoApproveMinorUnits: 1, MaxOpenPerSubject: 3},
	}
}

func TestMissingRequiredListsGaps(t *testing.T) {
	missing := (Config{}).MissingRequired()
	if len(missing) == 0 || !contains(missing, "LLM_API_KEY") || !contains(missing, "CONVERSATION_STORE") || !contains(missing, "CHANNEL_REGISTRY") {
		t.Fatalf("missing=%v", missing)
	}
	if missing := fullConfig().MissingRequired(); len(missing) != 0 {
		t.Fatalf("full config missing=%v", missing)
	}
}

func TestNewWiresSingleConnectedRuntime(t *testing.T) {
	a, err := New(fullConfig())
	if err != nil {
		t.Fatal(err)
	}
	if a.Conversation == nil || a.Dispatcher == nil || a.LLM == nil || a.Gateway == nil || a.Economy == nil || a.Approval == nil || a.Inventory == nil || a.Budget == nil {
		t.Fatalf("wiring left a nil component: %+v", a)
	}
}

func TestDispatcherConnectsChannelLLMApprovalDomainPersistenceAndReply(t *testing.T) {
	llmCalls := 0
	llm := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		llmCalls++
		w.Header().Set("content-type", "application/json")
		if llmCalls == 1 {
			_, _ = w.Write([]byte(`{"id":"resp-1","status":"completed","output":[{"type":"function_call","call_id":"call-1","name":"book_appointment","arguments":"{\"service\":\"corte\",\"when\":\"2026-09-05T15:00:00Z\"}"}],"usage":{"input_tokens":12,"output_tokens":4,"total_tokens":16}}`))
			return
		}
		_, _ = w.Write([]byte(`{"id":"resp-2","status":"completed","output":[{"type":"message","content":[{"type":"output_text","text":"Tu cita quedó confirmada."}]}],"usage":{"input_tokens":8,"output_tokens":5,"total_tokens":13}}`))
	}))
	defer llm.Close()
	domainCalls := 0
	domain := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		domainCalls++
		if r.URL.Path != "/v1/public/acme/north/appointments" || len(r.Header.Get("Idempotency-Key")) != 64 {
			t.Errorf("domain request path=%q key=%q", r.URL.Path, r.Header.Get("Idempotency-Key"))
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body["lead_id"] != "lead-1" {
			t.Errorf("domain body=%v err=%v", body, err)
		}
		w.Header().Set("content-type", "application/json")
		_, _ = w.Write([]byte(`{"appointment_id":"apt-1"}`))
	}))
	defer domain.Close()
	channel := &appChannel{}
	registry := channels.NewRegistry()
	if err := registry.Register(channel); err != nil {
		t.Fatal(err)
	}
	store := &appStore{}
	cfg := fullConfig()
	cfg.LLMBaseURL, cfg.DomainBaseURL, cfg.ChannelRegistry, cfg.ConversationStore = llm.URL, domain.URL, registry, store
	a, err := New(cfg)
	if err != nil {
		t.Fatal(err)
	}
	message := channels.Message{ChannelCode: "whatsapp", TenantID: cfg.TenantID, ExternalID: "wa-1", ThreadID: "thread-1", ProviderMessageID: "wamid.e2e.1", OccurredAt: time.Date(2026, 9, 4, 12, 0, 0, 0, time.UTC), Direction: channels.DirectionIn, Text: "Quiero un turno de corte mañana"}
	if err := a.Dispatcher.Handle(context.Background(), message); err != nil {
		t.Fatal(err)
	}
	if err := a.Dispatcher.Handle(context.Background(), message); err != nil {
		t.Fatal(err)
	}
	if llmCalls != 2 || domainCalls != 1 || len(channel.sent) != 2 {
		t.Fatalf("llm=%d domain=%d sent=%d", llmCalls, domainCalls, len(channel.sent))
	}
	if channel.sent[0].Text != "Tu cita quedó confirmada." || channel.sent[0].DeliveryKey != channel.sent[1].DeliveryKey {
		t.Fatalf("outbound=%#v", channel.sent)
	}
}

func TestDispatcherWithPostgresStoreIsExactlyOnceAcrossReplay(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	tenant := "018f4d4a-7b36-7a21-8d10-2f4c54c20002"
	_, _ = pool.Exec(ctx, `delete from communication.conversation_turn where tenant_id=$1`, tenant)
	_, _ = pool.Exec(ctx, `delete from platform.tenant where tenant_id=$1`, tenant)
	defer func() {
		_, _ = pool.Exec(ctx, `delete from communication.conversation_turn where tenant_id=$1`, tenant)
		_, _ = pool.Exec(ctx, `delete from platform.tenant where tenant_id=$1`, tenant)
	}()
	if _, err := pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,'app-e2e','App E2E','App')`, tenant); err != nil {
		t.Fatal(err)
	}
	store, err := platformpostgres.NewConversationStore(pool, 5*time.Second, 24*time.Hour, 3)
	if err != nil {
		t.Fatal(err)
	}
	llmCalls := 0
	llm := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		llmCalls++
		w.Header().Set("content-type", "application/json")
		if llmCalls == 1 {
			_, _ = w.Write([]byte(`{"id":"resp-pg-1","status":"completed","output":[{"type":"function_call","call_id":"call-pg-1","name":"create_quote","arguments":"{\"product\":\"scooter\",\"quantity\":1}"}],"usage":{"input_tokens":12,"output_tokens":4,"total_tokens":16}}`))
			return
		}
		_, _ = w.Write([]byte(`{"id":"resp-pg-2","status":"completed","output":[{"type":"message","content":[{"type":"output_text","text":"Cotización creada."}]}],"usage":{"input_tokens":8,"output_tokens":5,"total_tokens":13}}`))
	}))
	defer llm.Close()
	domainCalls := 0
	domain := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		domainCalls++
		if r.URL.Path != "/v1/franchise/quotes" || r.Header.Get("Authorization") != "Bearer token" || len(r.Header.Get("Idempotency-Key")) != 64 {
			t.Errorf("domain request path=%q auth=%q key=%q", r.URL.Path, r.Header.Get("Authorization"), r.Header.Get("Idempotency-Key"))
		}
		w.Header().Set("content-type", "application/json")
		_, _ = w.Write([]byte(`{"quotation_id":"quote-pg-1"}`))
	}))
	defer domain.Close()
	channel := &appChannel{}
	registry := channels.NewRegistry()
	if err := registry.Register(channel); err != nil {
		t.Fatal(err)
	}
	cfg := fullConfig()
	cfg.TenantID, cfg.LLMBaseURL, cfg.DomainBaseURL = tenant, llm.URL, domain.URL
	cfg.ChannelRegistry, cfg.ConversationStore = registry, store
	a, err := New(cfg)
	if err != nil {
		t.Fatal(err)
	}
	message := channels.Message{ChannelCode: "whatsapp", TenantID: tenant, ExternalID: "wa-pg", ThreadID: "thread-pg", ProviderMessageID: "wamid.app.pg.1", OccurredAt: time.Date(2026, 9, 4, 12, 0, 0, 0, time.UTC), Direction: channels.DirectionIn, Text: "Cotizame un scooter"}
	if err := a.Dispatcher.Handle(ctx, message); err != nil {
		t.Fatal(err)
	}
	if err := a.Dispatcher.Handle(ctx, message); err != nil {
		t.Fatal(err)
	}
	var state, assistant string
	var attempts int
	if err := pool.QueryRow(ctx, `select state,assistant_text,attempt_count from communication.conversation_turn where tenant_id=$1 and channel_code=$2 and provider_message_id=$3`, tenant, message.ChannelCode, message.ProviderMessageID).Scan(&state, &assistant, &attempts); err != nil {
		t.Fatal(err)
	}
	if state != "completed" || assistant != "Cotización creada." || attempts != 1 || llmCalls != 2 || domainCalls != 1 || len(channel.sent) != 2 {
		t.Fatalf("state=%s assistant=%q attempts=%d llm=%d domain=%d sent=%d", state, assistant, attempts, llmCalls, domainCalls, len(channel.sent))
	}
}

func TestNewFailsClosedOnMissing(t *testing.T) {
	if _, err := New(Config{}); err == nil {
		t.Fatal("expected error on missing config")
	}
}

func contains(list []string, s string) bool {
	for _, x := range list {
		if x == s {
			return true
		}
	}
	return false
}

func (*appStore) Reserve(_ context.Context, _ channels.Message, _ string, _ int, cap, tokens int64) error {
	if tokens > cap {
		return conversationruntime.ErrBudgetExceeded
	}
	return nil
}
func (*appStore) PinTool(context.Context, channels.Message, string, int, string) error { return nil }
