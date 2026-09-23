// Package analytics provides a governed semantic-metrics and data-ingest
// contract: versioned metric definitions, freshness and reconciliation gates,
// and a validated ingest record with replay idempotency and a late-data bound.
package analytics

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// Aggregation is the bounded set of metric aggregations.
type Aggregation string

const (
	AggCount Aggregation = "count"
	AggSum   Aggregation = "sum"
	AggMin   Aggregation = "min"
	AggMax   Aggregation = "max"
	AggAvg   Aggregation = "avg"
)

var (
	codeRe  = regexp.MustCompile(`^[a-z][a-z0-9._-]{0,63}$`)
	dimRe   = regexp.MustCompile(`^[a-z][a-z0-9._-]{0,63}$`)
	hex64Re = regexp.MustCompile(`^[0-9a-f]{64}$`)
)

var (
	ErrInvalidMetric = errors.New("analytics: invalid metric")
	ErrInvalidRecord = errors.New("analytics: invalid record")
)

// MetricDefinition is a versioned, tenant-scoped semantic metric. A version is
// immutable: the same code+version cannot be redefined.
type MetricDefinition struct {
	TenantID     string
	Code         string
	Version      int64
	Aggregation  Aggregation
	Source       string
	Dimensions   []string
	FreshnessTTL time.Duration
}

// Validate enforces the metric contract.
func (m MetricDefinition) Validate() error {
	if strings.TrimSpace(m.TenantID) == "" {
		return fmt.Errorf("%w: tenant", ErrInvalidMetric)
	}
	if !codeRe.MatchString(m.Code) {
		return fmt.Errorf("%w: code", ErrInvalidMetric)
	}
	if m.Version <= 0 {
		return fmt.Errorf("%w: version", ErrInvalidMetric)
	}
	switch m.Aggregation {
	case AggCount, AggSum, AggMin, AggMax, AggAvg:
	default:
		return fmt.Errorf("%w: aggregation", ErrInvalidMetric)
	}
	if strings.TrimSpace(m.Source) == "" {
		return fmt.Errorf("%w: source", ErrInvalidMetric)
	}
	if len(m.Dimensions) > 16 {
		return fmt.Errorf("%w: dimensions", ErrInvalidMetric)
	}
	seen := map[string]bool{}
	for _, d := range m.Dimensions {
		if !dimRe.MatchString(d) {
			return fmt.Errorf("%w: dimension %q", ErrInvalidMetric, d)
		}
		if seen[d] {
			return fmt.Errorf("%w: duplicate dimension %q", ErrInvalidMetric, d)
		}
		seen[d] = true
	}
	if m.FreshnessTTL <= 0 {
		return fmt.Errorf("%w: freshness ttl", ErrInvalidMetric)
	}
	return nil
}
