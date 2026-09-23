//go:build !windows

package main

import (
	"context"
	"errors"
)

// Native handles are a Windows-only trusted-launcher protocol.
func nativeStopContext(parent context.Context, raw string) (context.Context, func() error, error) {
	if raw != "" {
		return parent, func() error { return nil }, errors.New("NATIVE_STOP_PROTOCOL_FAILED")
	}
	ctx, cancel := context.WithCancel(parent)
	return ctx, func() error { cancel(); return nil }, nil
}

func nativeStopFailure(context.Context) error { return nil }
