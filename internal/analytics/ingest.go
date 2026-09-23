package analytics

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

var sourceRe = regexp.MustCompile(`^[a-z][a-z0-9._-]{0,63}$`)

// IngestRecord is a validated raw intake: it carries source ownership, an
// idempotency replay key and a bounded late-data window.
type IngestRecord struct {
	SourceID      string
	ExternalID    string
	EventTime     time.Time
	ReceivedAt    time.Time
	PayloadSHA256 string
	MaxLateness   time.Duration // 0 means unbounded (no late-data rejection)
}

// Validate enforces the ingest contract: source identity, external id, payload
// digest, event/received ordering and the late-data bound.
func (r IngestRecord) Validate() error {
	if !sourceRe.MatchString(r.SourceID) {
		return fmt.Errorf("%w: source", ErrInvalidRecord)
	}
	if strings.TrimSpace(r.ExternalID) == "" || len(r.ExternalID) > 200 {
		return fmt.Errorf("%w: external id", ErrInvalidRecord)
	}
	if !hex64Re.MatchString(strings.ToLower(r.PayloadSHA256)) {
		return fmt.Errorf("%w: payload sha", ErrInvalidRecord)
	}
	if r.EventTime.IsZero() || r.ReceivedAt.IsZero() {
		return fmt.Errorf("%w: timestamps", ErrInvalidRecord)
	}
	if r.ReceivedAt.Before(r.EventTime) {
		return fmt.Errorf("%w: received before event", ErrInvalidRecord)
	}
	if r.MaxLateness > 0 && r.ReceivedAt.Sub(r.EventTime) > r.MaxLateness {
		return fmt.Errorf("%w: exceeds lateness bound", ErrInvalidRecord)
	}
	return nil
}

// ReplayKey returns the idempotency key for deduplicated replay.
func (r IngestRecord) ReplayKey() string {
	return r.SourceID + "\x00" + r.ExternalID + "\x00" + strings.ToLower(r.PayloadSHA256)
}
