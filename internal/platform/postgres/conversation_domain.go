package postgres

// AUTHORED local composition glue. Existing domain owners implement quotation
// and appointment mutations; customer_order remains the authority for status.
import (
	"context"
	"elite.local/enterprise/internal/agenttools"
	"elite.local/enterprise/internal/conversationruntime"
	"elite.local/enterprise/internal/domainbind"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"strings"
)

type ConversationDomain struct {
	pool                 *pgxpool.Pool
	delegate             conversationruntime.Domain
	tenant, organization string
}

func NewConversationDomain(pool *pgxpool.Pool, delegate conversationruntime.Domain, tenant, organization string) (*ConversationDomain, error) {
	if pool == nil || delegate == nil || strings.TrimSpace(tenant) == "" || strings.TrimSpace(organization) == "" {
		return nil, conversationruntime.ErrInvalidRuntime
	}
	return &ConversationDomain{pool: pool, delegate: delegate, tenant: tenant, organization: organization}, nil
}
func (d *ConversationDomain) allowed(tenant string, s domainbind.Scope) bool {
	return d != nil && tenant == d.tenant && s.OrganizationID == d.organization && strings.TrimSpace(s.LeadID) != ""
}
func (d *ConversationDomain) BookAppointmentFor(ctx context.Context, tenant string, s domainbind.Scope, key domainbind.CommandIdentity, in agenttools.AppointmentInput) (string, error) {
	if !d.allowed(tenant, s) {
		return "", conversationruntime.ErrInvalidRuntime
	}
	return d.delegate.BookAppointmentFor(ctx, tenant, s, key, in)
}
func (d *ConversationDomain) CreateQuoteFor(ctx context.Context, tenant string, s domainbind.Scope, key domainbind.CommandIdentity, in agenttools.QuoteInput) (string, error) {
	if !d.allowed(tenant, s) {
		return "", conversationruntime.ErrInvalidRuntime
	}
	return d.delegate.CreateQuoteFor(ctx, tenant, s, key, in)
}
func (d *ConversationDomain) OrderStatusFor(ctx context.Context, tenant string, s domainbind.Scope, id string) (string, error) {
	if !d.allowed(tenant, s) || strings.TrimSpace(id) == "" || len(id) > 128 {
		return "", conversationruntime.ErrInvalidRuntime
	}
	var state string
	err := d.pool.QueryRow(ctx, `select o.state from sales.customer_order o join crm.lead l on l.tenant_id=o.tenant_id and l.organization_id=o.organization_id and l.customer_principal_id=o.customer_principal_id where o.tenant_id=$1 and o.organization_id=$2 and l.lead_id=$3 and o.order_id=$4`, tenant, s.OrganizationID, s.LeadID, id).Scan(&state)
	if err != nil {
		return "", errors.New("conversation domain: order unavailable for resolved contact")
	}
	return fmt.Sprintf("Pedido %s: %s.", id, state), nil
}
