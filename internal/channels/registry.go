package channels

import (
	"fmt"
	"sync"
)

// Registry authorizes at most one channel per code.
type Registry struct {
	mu       sync.RWMutex
	channels map[string]Channel
}

// NewRegistry returns an empty registry.
func NewRegistry() *Registry {
	return &Registry{channels: make(map[string]Channel)}
}

// Register binds a channel to its code, rejecting duplicates and invalid codes.
func (r *Registry) Register(c Channel) error {
	if c == nil {
		return fmt.Errorf("channels: nil channel")
	}
	if !codeRe.MatchString(c.Code()) {
		return fmt.Errorf("%w: invalid code %q", ErrInvalidMessage, c.Code())
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.channels[c.Code()]; ok {
		return fmt.Errorf("%w: %s", ErrDuplicateCode, c.Code())
	}
	r.channels[c.Code()] = c
	return nil
}

// For returns the channel bound to a code, if any.
func (r *Registry) For(code string) (Channel, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.channels[code]
	return c, ok
}

// Count returns the number of explicitly registered channel adapters.
func (r *Registry) Count() int {
	if r == nil {
		return 0
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.channels)
}
