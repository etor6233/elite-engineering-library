// Package economy supplies local token-allowance and scoped-cache helpers.
// Durable generation reservation remains the connected PostgreSQL owner.
package economy

// EstimateTokens is a bytes/3 heuristic, not a provider tokenizer or upper bound.
// Its budget is only an estimate; the provider runtime governs actual tokens.
func EstimateTokens(text string) int {
	if len(text) == 0 {
		return 0
	}
	n := (len(text) + 2) / 3
	if n < 1 {
		n = 1
	}
	return n
}
