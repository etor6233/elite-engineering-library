package app

import "fmt"

// NewProduction adds the delivery invariant that every configured channel is
// fenced by a durable attempt/receipt store. New remains available for local
// tests and compositions that intentionally do not send externally.
func NewProduction(cfg Config) (*App, error) {
	if cfg.ChannelRegistry == nil || !cfg.ChannelRegistry.AllDurable() {
		return nil, fmt.Errorf("%w: all channels require durable delivery", ErrIncomplete)
	}
	return New(cfg)
}
