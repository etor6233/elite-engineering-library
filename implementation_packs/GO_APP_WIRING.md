# Go App Wiring

## 1. Metadata

```yaml
pack_id: "GO-APP-WIRING"
pack_version: "0.2.1"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Ensambla un único runtime conectado de canal, Responses tool calling, safety, presupuesto, aprobación, dominio con scope dinámico, persistencia durable y respuesta; falla cerrado si falta un input obligatorio."
stacks: ["Go 1.26.7", "OpenAI Responses API"]
compatible_with: ["GO-CONNECTED-CONVERSATION-RUNTIME 0.1.x", "GO-CHANNELS-CORE 0.4.x", "GO-LLM-OPENAI-ADAPTER 0.2.x", "GO-AGENT-DOMAIN-BINDING 0.4.x", "GO-LLM-ECONOMY-CORE 0.1.x", "GO-HUMAN-APPROVAL-CORE 0.1.x", "GO-FINOPS-CORE 0.1.x"]
incompatible_with: ["componentes construidos pero fuera del recorrido", "canales sin adapter registrado", "scope inferido por el modelo", "presupuesto LLM no configurado"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://go.dev", "https://developers.openai.com/api/docs/guides/function-calling", "https://sre.google/sre-book/reliable-product-launches/"]
verified_at: "2026-09-04"
```

Los bloques son `AUTHORED`; las fuentes gobiernan el contrato, no la autoría. El pack queda `CONDITIONED` por credenciales, adapters live, configuración comercial, política de privacidad, target y aceptación, no por wiring local faltante.

## 2. Applicability

Use esta fábrica como composition root del agente omnicanal. Deben inyectarse store durable, resolver autorizado de contacto, al menos un canal, token rotatorio del dominio, presupuesto y configuración explícita del runtime. No mantenga un segundo agente paralelo.

## 3. Architecture contract

- **Flujo**: `App.Dispatcher` → `Conversation.Handle` → LLM/economía/aprobación/gateway/store → canal.
- **Fail closed**: `MissingRequired` enumera credenciales, stores, resolver, registry y presupuesto ausentes; `Runtime.Validate` gobierna límites y consentimiento de almacenamiento Responses.
- **Scope**: organización y lead vienen del resolver por contacto; los defaults globales son opcionales y sólo compatibles con llamadas heredadas.
- **Costo**: el ledger recibe presupuesto tenant-scoped antes de procesar. Su persistencia distribuida sigue siendo condición del proyecto.
- **FinOps**: inventario/budget de infraestructura permanece en el composition root, pero no se presenta falsamente como costo durable por conversación.

## 4. Exact file manifest

```text
CREATE internal/app/config.go
CREATE internal/app/app.go
CREATE internal/app/app_test.go
```

## 5. Materialization blocks

### FILE: `internal/app/config.go`
```yaml
block_id: "GO-APP-WIRING:internal/app/config.go:v2"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "ea32755899fdaa84f062521bbdf09f7c9d9263bb16a3f4d4ca978c1d22ae4b8a"
variables: []
secrets_allowed: false
```
````go
// Package app wires the whole franchise runtime into one assembly: the
// conversational agent, the LLM adapter, the domain binding, the cost governor,
// the anti-fraud approval queue and the FinOps inventory/budget. It reports
// exactly which credentials are missing, fail-closed, so the agent never blocks
// on a guess — it tells the operator the precise gap.
package app

import (
	"sort"
	"strings"

	"elite.local/enterprise/internal/approval"
	"elite.local/enterprise/internal/channels"
	"elite.local/enterprise/internal/conversationruntime"
	"elite.local/enterprise/internal/domainbind"
)

// Config carries every runtime input. Secrets are injected from the
// environment; the assembly only validates that they are present.
type Config struct {
	LLMAPIKey           string
	LLMModel            string
	LLMBaseURL          string
	DomainBaseURL       string
	DomainTokenProvider domainbind.AccessTokenProvider
	TenantID            string
	TenantCode          string
	OrganizationID      string
	LeadID              string
	ServiceKinds        map[string]string
	ProductVariants     map[string]string
	PriceBookID         string
	ConversationStore   conversationruntime.Store
	ContactResolver     conversationruntime.ContactResolver
	ChannelRegistry     *channels.Registry
	Conversation        conversationruntime.Config
	LLMTokenBudget      int64
	ApprovalPolicy      approval.Policy
}

// MissingRequired returns the names of the required inputs still absent. An
// empty list means the assembly can proceed.
func (c Config) MissingRequired() []string {
	var missing []string
	trim := strings.TrimSpace
	if trim(c.LLMAPIKey) == "" {
		missing = append(missing, "LLM_API_KEY")
	}
	if trim(c.LLMModel) == "" {
		missing = append(missing, "LLM_MODEL")
	}
	if trim(c.LLMBaseURL) == "" {
		missing = append(missing, "LLM_BASE_URL")
	}
	if trim(c.DomainBaseURL) == "" {
		missing = append(missing, "DOMAIN_BASE_URL")
	}
	if c.DomainTokenProvider == nil {
		missing = append(missing, "DOMAIN_TOKEN_PROVIDER")
	}
	if trim(c.TenantID) == "" {
		missing = append(missing, "TENANT_ID")
	}
	if trim(c.TenantCode) == "" {
		missing = append(missing, "TENANT_CODE")
	}
	if len(c.ServiceKinds) == 0 {
		missing = append(missing, "SERVICE_KINDS")
	}
	if len(c.ProductVariants) == 0 {
		missing = append(missing, "PRODUCT_VARIANTS")
	}
	if trim(c.PriceBookID) == "" {
		missing = append(missing, "PRICE_BOOK_ID")
	}
	if c.ConversationStore == nil {
		missing = append(missing, "CONVERSATION_STORE")
	}
	if c.ContactResolver == nil {
		missing = append(missing, "CONTACT_RESOLVER")
	}
	if c.ChannelRegistry == nil || c.ChannelRegistry.Count() == 0 {
		missing = append(missing, "CHANNEL_REGISTRY")
	}
	if c.LLMTokenBudget <= 0 {
		missing = append(missing, "LLM_TOKEN_BUDGET")
	}
	sort.Strings(missing)
	return missing
}
````

### FILE: `internal/app/app.go`
```yaml
block_id: "GO-APP-WIRING:internal/app/app.go:v2"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "b79f89913f11c4afc73ff61b4ca3a3a57c0cbed3a0bc4d827bbe62214c9aab28"
variables: []
secrets_allowed: false
```
````go
package app

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"elite.local/enterprise/internal/aifoundation"
	"elite.local/enterprise/internal/approval"
	"elite.local/enterprise/internal/channels"
	"elite.local/enterprise/internal/conversationruntime"
	"elite.local/enterprise/internal/domainbind"
	"elite.local/enterprise/internal/economy"
	"elite.local/enterprise/internal/finops"
	"elite.local/enterprise/internal/llmopenai"
)

// App is the fully-wired franchise runtime. Every component is connected:
// the agent delegates to the domain, governs cost, and hands off to human
// approval; FinOps tracks every deployed resource.
type App struct {
	Conversation *conversationruntime.Runtime
	Dispatcher   *channels.Dispatcher
	LLM          *llmopenai.Client
	Gateway      *domainbind.Gateway
	Economy      *economy.Governor
	Approval     *approval.Registry
	Inventory    *finops.Inventory
	Budget       *finops.Budget
}

// New assembles the runtime from a validated config. It fails closed with a
// list of missing inputs rather than guessing.
func New(cfg Config) (*App, error) {
	if missing := cfg.MissingRequired(); len(missing) > 0 {
		return nil, fmt.Errorf("%w: missing required inputs: %s", ErrIncomplete, strings.Join(missing, ", "))
	}

	llm, err := llmopenai.New(llmopenai.Config{
		BaseURL: cfg.LLMBaseURL,
		Model:   cfg.LLMModel,
		APIKey:  cfg.LLMAPIKey,
		Timeout: 30 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	gw, err := domainbind.NewGateway(domainbind.Config{
		BaseURL:         cfg.DomainBaseURL,
		TenantID:        cfg.TenantID,
		TenantCode:      cfg.TenantCode,
		OrganizationID:  cfg.OrganizationID,
		LeadID:          cfg.LeadID,
		ServiceKinds:    cfg.ServiceKinds,
		ProductVariants: cfg.ProductVariants,
		PriceBookID:     cfg.PriceBookID,
		ValidMinutes:    60,
		Timeout:         10 * time.Second,
		TokenProvider:   cfg.DomainTokenProvider,
	})
	if err != nil {
		return nil, err
	}

	safety := aifoundation.NewSafetyGate()
	ledger := economy.NewCostLedger()
	ledger.SetBudget(cfg.TenantID, cfg.LLMTokenBudget)
	eco := &economy.Governor{
		Cache:     &economy.SemanticCache{},
		Ledger:    ledger,
		Threshold: 0.9,
	}
	appr := approval.NewRegistry(cfg.ApprovalPolicy)
	cfg.Conversation.TokenBudget = cfg.LLMTokenBudget
	cfg.Conversation.ModelRevision = cfg.LLMModel
	runtime := &conversationruntime.Runtime{Config: cfg.Conversation, Store: cfg.ConversationStore, Model: llm, Resolver: cfg.ContactResolver, Domain: gw, Approval: appr, Safety: safety}
	if err := runtime.Validate(); err != nil {
		return nil, err
	}
	dispatcher := &channels.Dispatcher{Registry: cfg.ChannelRegistry, Respond: runtime.Handle}

	return &App{
		Conversation: runtime,
		Dispatcher:   dispatcher,
		LLM:          llm,
		Gateway:      gw,
		Economy:      eco,
		Approval:     appr,
		Inventory:    finops.NewInventory(),
		Budget:       finops.NewBudget(),
	}, nil
}

// ErrIncomplete is the sentinel for a config missing required inputs.
var ErrIncomplete = errors.New("app: incomplete config")
````

### FILE: `internal/app/app_test.go`
```yaml
block_id: "GO-APP-WIRING:internal/app/app_test.go:v2"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "e583d959f14c6ac7c0d862e25fe5184d1f34fdce5fe8f3740c7c43364a91d11f"
variables: []
secrets_allowed: false
```
````go
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
````

## 6. Configuration surface

Requiere API/model/base URL exactos, token provider, tenant/código, mappings, price book, store durable, contact resolver, registry con al menos un canal, presupuesto y política de aprobación. `Conversation.Config` fija límites, handoff y consentimiento de storage Responses.

## 7. Dependency bill

| Dependency | Pin | License | Purpose |
|---|---:|---|---|
| Go | 1.26.7 | BSD-3-Clause | composition root |
| OpenAI Responses API | revisión consultada 2026-09-04 | términos del servicio | model/tool runtime, no código redistribuido |

## 8. Apply order

Materializar todos los packs compatibles, implementar store/resolver/adapters, construir `App`, exponer sólo `Dispatcher` y cerrar configuración antes de recibir tráfico.

## 9. Verification

Materializar 3/3, aplicar `gofmt`, ejecutar tests del app con servidores HTTP controlados y con PostgreSQL real, y luego `go test ./...`, `go vet ./...`, `go build ./...`. La prueba local no autoriza cuentas live ni producción.

## 10. Reconstruction evidence

V236: tres archivos reconstruidos byte-exactos/gofmt; E2E de cita en memoria y cotización con PostgreSQL real probaron que replay no repite inferencia ni dominio. Véase `reconstruction_evidence/GO_CONNECTED_CONVERSATION_RUNTIME_2026-09-04_V236.md`.

V402 composed delta: T2807 delta314: generation/lease fencing, expired history/replay refusal, durable per-attempt token reservation, pinned tool intent/contact/policy and explicit continuation instructions. AUTHORED glue; no new upstream dependencies. AI_RUNTIME_GOVERNANCE_V402.md/json. T2807 connected eval/history closure remains open.
