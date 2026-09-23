package finops

import (
	"errors"
	"sync"
)

// ErrOverBudget reports a spend that would exceed the tenant's budget.
var ErrOverBudget = errors.New("finops: over budget")

// Budget tracks spend per tenant against a cap, fail-closed.
type Budget struct {
	mu    sync.Mutex
	caps  map[string]int64
	spent map[string]int64
}

// NewBudget returns an empty budget tracker.
func NewBudget() *Budget {
	return &Budget{caps: make(map[string]int64), spent: make(map[string]int64)}
}

// SetCap sets the tenant's spending cap (clamped to >= 0).
func (b *Budget) SetCap(tenant string, cap int64) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if cap < 0 {
		cap = 0
	}
	b.caps[tenant] = cap
}

// Record adds spend; it fails closed if the cap would be exceeded.
func (b *Budget) Record(tenant string, amount int64) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if amount < 0 {
		return errors.New("finops: negative spend")
	}
	if b.spent[tenant]+amount > b.caps[tenant] {
		return ErrOverBudget
	}
	b.spent[tenant] += amount
	return nil
}

// Remaining returns the tenant's remaining budget.
func (b *Budget) Remaining(tenant string) int64 {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.caps[tenant] - b.spent[tenant]
}

// OverBudget reports whether the tenant has no remaining budget.
func (b *Budget) OverBudget(tenant string) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.spent[tenant] >= b.caps[tenant]
}
