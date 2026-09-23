package economy

import (
	"errors"
	"sync"
)

// ErrBudgetExceeded reports a spend that would exceed the tenant's budget.
var ErrBudgetExceeded = errors.New("economy: budget exceeded")

// CostLedger tracks token spend per tenant against a budget. It fails closed:
// an over-budget tenant cannot spend.
type CostLedger struct {
	mu     sync.Mutex
	budget map[string]int64
}

// NewCostLedger returns an empty ledger.
func NewCostLedger() *CostLedger {
	return &CostLedger{budget: make(map[string]int64)}
}

// SetBudget sets the remaining-token budget for a tenant (clamped to >= 0).
func (l *CostLedger) SetBudget(tenant string, tokens int64) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if tokens < 0 {
		tokens = 0
	}
	l.budget[tenant] = tokens
}

// Remaining returns the tenant's remaining budget.
func (l *CostLedger) Remaining(tenant string) int64 {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.budget[tenant]
}

// CanSpend reports whether spending tokens would stay within budget.
func (l *CostLedger) CanSpend(tenant string, tokens int64) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if tokens < 0 {
		return false
	}
	return l.budget[tenant] >= tokens
}

// Spend consumes tokens; it fails closed if the budget would be exceeded.
func (l *CostLedger) Spend(tenant string, tokens int64) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if tokens < 0 {
		return errors.New("economy: negative spend")
	}
	if l.budget[tenant] < tokens {
		return ErrBudgetExceeded
	}
	l.budget[tenant] -= tokens
	return nil
}
