// Package cache provides a governed cache boundary: tenant-scoped keys, a
// TTL-aware LRU memory store, and single-flight cache-aside loading. A cache
// never holds source of truth; misses are repopulated from the owning store
// and expiry/invalidation are expected.
package cache

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// Store is the cache boundary. Implementations must be safe for concurrent use
// and must never serve a miss as an authoritative value.
type Store interface {
	Get(ctx context.Context, key string) (value []byte, ok bool, err error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}

var (
	ErrInvalidKey = errors.New("cache: invalid key")
	ErrInvalidTTL = errors.New("cache: invalid ttl")
)

var keyRe = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:/-]{0,255}$`)

// ValidateKey rejects empty, unsafe or overlong keys to keep a shared cache
// namespace collision-free and free of injection surprises.
func ValidateKey(key string) error {
	if !keyRe.MatchString(key) {
		return fmt.Errorf("%w: %q", ErrInvalidKey, key)
	}
	return nil
}

// Scope namespaces a key by tenant (mandatory) and optional kind/id.
type Scope struct {
	TenantID string
	Kind     string
	ID       string
}

// Key builds a tenant-scoped, collision-free key. Empty tenant is denied.
func (s Scope) Key() (string, error) {
	if strings.TrimSpace(s.TenantID) == "" {
		return "", fmt.Errorf("%w: empty tenant", ErrInvalidKey)
	}
	parts := []string{"t", s.TenantID}
	if s.Kind != "" {
		parts = append(parts, s.Kind)
	}
	if s.ID != "" {
		parts = append(parts, s.ID)
	}
	k := strings.Join(parts, ":")
	if err := ValidateKey(k); err != nil {
		return "", err
	}
	return k, nil
}
