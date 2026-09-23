package economy

import "strings"

// Path is the cheapest valid resolution for a request.
type Path int

const (
	PathCache   Path = iota // generation avoided; infrastructure/embedding costs are not measured
	PathLLM                 // requires a budgeted, compressed model call
	PathBlocked             // over budget, fail-closed
)

// Governor consumes a local token allowance atomically. Decide never consults
// the unscoped legacy cache. Use DecideScoped for tenant/subject/policy-scoped hits.
// This in-memory allowance is not provider billing, durable or distributed.
type Governor struct {
	Cache     *SemanticCache
	Ledger    *CostLedger
	Threshold float64
}

// Decide returns the path and, for a cache hit, the cached answer.
func (g *Governor) Decide(tenant string, embedding []float32, estimatedTokens int64) (Path, string) {
	if strings.TrimSpace(tenant) == "" || estimatedTokens <= 0 {
		return PathBlocked, ""
	}
	if g.Ledger == nil || g.Ledger.Spend(tenant, estimatedTokens) != nil {
		return PathBlocked, ""
	}
	return PathLLM, ""
}

// DecideScoped uses reader/revision identity before consuming a local allowance.
// Durable production reservation remains the connected PostgreSQL owner.
func (g *Governor) DecideScoped(scope CacheScope, embedding []float32, estimatedTokens int64) (Path, string) {
	if !scope.valid() || estimatedTokens <= 0 {
		return PathBlocked, ""
	}
	if g.Cache != nil {
		if answer, ok := g.Cache.LookupScoped(scope, embedding, g.Threshold); ok {
			return PathCache, answer
		}
	}
	return g.Decide(scope.Tenant, embedding, estimatedTokens)
}
