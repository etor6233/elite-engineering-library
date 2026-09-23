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
