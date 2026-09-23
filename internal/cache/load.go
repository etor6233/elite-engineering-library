package cache

import (
	"context"
	"sync"
	"time"
)

// Loader wraps a Store with cache-aside single-flight loading: concurrent
// misses for the same key trigger exactly one load and share the result,
// preventing a thundering herd against the source of truth.
type Loader struct {
	Store  Store
	mu     sync.Mutex
	flight map[string]*call
}

type call struct {
	done chan struct{}
	val  []byte
	err  error
}

// NewLoader returns a Loader over the given Store.
func NewLoader(s Store) *Loader {
	return &Loader{Store: s, flight: make(map[string]*call)}
}

// Get returns the cached value or loads and populates it once per key. A load
// error is returned without poisoning the cache (nothing is stored).
func (l *Loader) Get(ctx context.Context, key string, ttl time.Duration, load func(context.Context) ([]byte, error)) ([]byte, error) {
	if v, ok, err := l.Store.Get(ctx, key); err == nil && ok {
		return v, nil
	}

	l.mu.Lock()
	if c, ok := l.flight[key]; ok {
		l.mu.Unlock()
		<-c.done
		return c.val, c.err
	}
	c := &call{done: make(chan struct{})}
	l.flight[key] = c
	l.mu.Unlock()

	c.val, c.err = load(ctx)
	if c.err == nil {
		_ = l.Store.Set(ctx, key, c.val, ttl)
	}

	close(c.done)
	l.mu.Lock()
	delete(l.flight, key)
	l.mu.Unlock()
	return c.val, c.err
}
