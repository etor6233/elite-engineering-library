package search

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

var (
	ErrInvalidQuery = errors.New("search: invalid query")

	maxQueryRunes = 500
	maxTopK       = 100
)

// ValidateQueryInput rejects empty, whitespace-only or overlong query text.
// The accepted string is intended for PostgreSQL websearch_to_tsquery, which
// parses a user-friendly subset and never treats input as SQL.
func ValidateQueryInput(text string) error {
	t := strings.TrimSpace(text)
	if t == "" {
		return fmt.Errorf("%w: empty", ErrInvalidQuery)
	}
	if len([]rune(t)) > maxQueryRunes {
		return fmt.Errorf("%w: too long", ErrInvalidQuery)
	}
	return nil
}

// ClampTopK bounds a caller-provided limit into the safe range.
func ClampTopK(k int) int {
	if k < 1 {
		return 1
	}
	if k > maxTopK {
		return maxTopK
	}
	return k
}

// tokenize splits text into lowercase word tokens for the in-memory index.
func tokenize(text string) []string {
	return strings.FieldsFunc(strings.ToLower(text), func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
}
