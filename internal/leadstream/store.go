package leadstream

import (
	"context"
	"errors"
	"sync"
)

type ReceiptState string

const (
	ReceiptNormalized ReceiptState = "normalized"
	ReceiptRejected   ReceiptState = "rejected"
	ReceiptDuplicate  ReceiptState = "duplicate"
)

type Receipt struct {
	State        ReceiptState
	SourceSHA256 string
}

type Store interface {
	Record(ctx context.Context, raw RawEvent, candidate *LeadCandidate, normalizationCode string) (Receipt, error)
}

type memoryRecord struct {
	raw       RawEvent
	candidate *LeadCandidate
	code      string
}

// MemoryStore is only a deterministic test/reference adapter. Production must
// inject a durable Store such as postgres.LeadIngress.
type MemoryStore struct {
	mu      sync.Mutex
	records map[string]memoryRecord
	Fail    error
}

func NewMemoryStore() *MemoryStore { return &MemoryStore{records: map[string]memoryRecord{}} }

func eventKey(raw RawEvent) string {
	return raw.TenantID + "\x00" + raw.Provider + "\x00" + raw.ProviderEventID
}

func (s *MemoryStore) Record(_ context.Context, raw RawEvent, candidate *LeadCandidate, normalizationCode string) (Receipt, error) {
	if s == nil {
		return Receipt{}, errors.New("leadstream: nil store")
	}
	if err := raw.Validate(); err != nil {
		return Receipt{}, err
	}
	if candidate == nil && normalizationCode == "" {
		return Receipt{}, errors.New("leadstream: candidate or rejection required")
	}
	if candidate != nil {
		if normalizationCode != "" {
			return Receipt{}, errors.New("leadstream: candidate and rejection are mutually exclusive")
		}
		if err := candidate.Validate(); err != nil {
			return Receipt{}, err
		}
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.Fail != nil {
		return Receipt{}, s.Fail
	}
	key := eventKey(raw)
	if prior, ok := s.records[key]; ok {
		if prior.raw.SourceSHA256 != raw.SourceSHA256 {
			return Receipt{}, ErrDivergentDuplicate
		}
		return Receipt{State: ReceiptDuplicate, SourceSHA256: raw.SourceSHA256}, nil
	}
	var cloned *LeadCandidate
	if candidate != nil {
		copy := *candidate
		copy.Fields = append([]Field(nil), candidate.Fields...)
		cloned = &copy
	}
	s.records[key] = memoryRecord{raw: raw, candidate: cloned, code: normalizationCode}
	state := ReceiptNormalized
	if cloned == nil {
		state = ReceiptRejected
	}
	return Receipt{State: state, SourceSHA256: raw.SourceSHA256}, nil
}

func (s *MemoryStore) Count() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.records)
}
