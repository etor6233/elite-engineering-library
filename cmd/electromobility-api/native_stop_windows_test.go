//go:build windows

package main

import (
	"context"
	"strconv"
	"testing"
	"time"
)

func TestAPINativeStopLifecycle(t *testing.T) {
	h, _, _ := stopKernel.NewProc("CreateEventW").Call(0, 1, 0, 0)
	if h == 0 {
		t.Fatal("event creation")
	}
	defer stopClose.Call(h)
	ctx, cleanup, err := nativeStopContext(context.Background(), strconv.FormatUint(uint64(h), 10))
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()
	if ctx.Err() != nil {
		t.Fatal("premature stop")
	}
	stopKernel.NewProc("SetEvent").Call(h)
	select {
	case <-ctx.Done():
	case <-time.After(time.Second):
		t.Fatal("stop did not cancel host")
	}
	if err := cleanup(); err != nil {
		t.Fatal(err)
	}
	ctx, cleanup, err = nativeStopContext(context.Background(), strconv.FormatUint(uint64(h), 10))
	if err != nil || ctx.Err() == nil {
		t.Fatal("pre-signaled stop must precede host startup")
	}
	if err := cleanup(); err != nil {
		t.Fatal(err)
	}
}

func TestAPINativeStopInvalidHandle(t *testing.T) {
	for _, raw := range []string{"0", "-1", "+4", "04", "999999999999999999999999"} {
		if _, _, err := nativeStopContext(context.Background(), raw); err == nil {
			t.Fatal("invalid handle accepted")
		}
	}
}
