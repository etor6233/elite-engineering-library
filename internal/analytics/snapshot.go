package analytics

import "time"

// Fresh reports whether a value computed at computedAt is still within its
// freshness TTL as of now. A zero/negative TTL is never fresh (fail-closed).
func Fresh(computedAt, now time.Time, ttl time.Duration) bool {
	if ttl <= 0 || computedAt.IsZero() || now.Before(computedAt) {
		return false
	}
	return now.Sub(computedAt) <= ttl
}

// Reconciles reports whether the sum of parts equals the total within a
// non-negative absolute tolerance. Reconciliation never passes on a negative
// tolerance.
func Reconciles(parts []float64, total, tolerance float64) bool {
	if tolerance < 0 {
		return false
	}
	var sum float64
	for _, p := range parts {
		sum += p
	}
	d := sum - total
	if d < 0 {
		d = -d
	}
	return d <= tolerance
}
