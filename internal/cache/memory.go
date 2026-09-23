package cache

import (
	"container/list"
	"context"
	"errors"
	"sync"
	"time"
)

// ErrTooLarge reports a value exceeding the store's byte budget.
var ErrTooLarge = errors.New("cache: value exceeds budget")

type entry struct {
	key      string
	value    []byte
	expires  time.Time
	listElem *list.Element
}

// MemoryStore is a bounded, TTL-aware, LRU-evicting cache for a single
// process. It is the deterministic reference implementation; a Redis adapter
// provides the shared/HA equivalent and remains CONDITIONED.
type MemoryStore struct {
	mu        sync.Mutex
	items     map[string]*entry
	lru       *list.List // front = most recently used
	maxBytes  int64
	currBytes int64
	clock     func() time.Time
}

// NewMemoryStore returns a store with the given byte budget (default 1 MiB).
func NewMemoryStore(maxBytes int64) *MemoryStore {
	if maxBytes <= 0 {
		maxBytes = 1 << 20
	}
	return &MemoryStore{
		items:    make(map[string]*entry),
		lru:      list.New(),
		maxBytes: maxBytes,
		clock:    time.Now,
	}
}

func (m *MemoryStore) now() time.Time {
	if m.clock != nil {
		return m.clock()
	}
	return time.Now()
}

// Get returns a copy of the cached value, or ok=false on miss/expiry.
func (m *MemoryStore) Get(_ context.Context, key string) ([]byte, bool, error) {
	if err := ValidateKey(key); err != nil {
		return nil, false, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	e, ok := m.items[key]
	if !ok {
		return nil, false, nil
	}
	if m.now().After(e.expires) {
		m.remove(key)
		return nil, false, nil
	}
	m.lru.MoveToFront(e.listElem)
	out := make([]byte, len(e.value))
	copy(out, e.value)
	return out, true, nil
}

// Set stores a copy of value with a positive TTL and evicts LRU to fit budget.
func (m *MemoryStore) Set(_ context.Context, key string, value []byte, ttl time.Duration) error {
	if err := ValidateKey(key); err != nil {
		return err
	}
	if ttl <= 0 {
		return ErrInvalidTTL
	}
	if int64(len(value)) > m.maxBytes {
		return ErrTooLarge
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.set(key, value, m.now().Add(ttl))
	return nil
}

func (m *MemoryStore) set(key string, value []byte, expires time.Time) {
	if e, ok := m.items[key]; ok {
		m.currBytes -= int64(len(e.value))
		e.value = make([]byte, len(value))
		copy(e.value, value)
		e.expires = expires
		m.currBytes += int64(len(e.value))
		m.lru.MoveToFront(e.listElem)
		return
	}
	e := &entry{key: key, value: append([]byte(nil), value...), expires: expires}
	e.listElem = m.lru.PushFront(e)
	m.items[key] = e
	m.currBytes += int64(len(e.value))

	for m.currBytes > m.maxBytes {
		back := m.lru.Back()
		if back == nil {
			break
		}
		m.remove(back.Value.(*entry).key)
	}
}

// Delete removes a key; it is idempotent.
func (m *MemoryStore) Delete(_ context.Context, key string) error {
	if err := ValidateKey(key); err != nil {
		return err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.items[key]; !ok {
		return nil
	}
	m.remove(key)
	return nil
}

func (m *MemoryStore) remove(key string) {
	if e, ok := m.items[key]; ok {
		m.lru.Remove(e.listElem)
		m.currBytes -= int64(len(e.value))
		delete(m.items, key)
	}
}

// Len reports the number of cached entries.
func (m *MemoryStore) Len() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.items)
}
