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
