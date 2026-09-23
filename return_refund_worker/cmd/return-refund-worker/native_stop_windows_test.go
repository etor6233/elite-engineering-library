//go:build windows

package main

import (
	"context"
	"errors"
	"strconv"
	"sync"
	"testing"
	"time"
	"unsafe"
)

func stopTestEvent(t *testing.T) uintptr {
	t.Helper()
	h, _, _ := stopKernel.NewProc("CreateEventW").Call(0, 1, 0, 0)
	if h == 0 {
		t.Fatal("event creation")
	}
	t.Cleanup(func() {
		if ok, _, _ := stopClose.Call(h); ok == 0 {
			t.Error("event close")
		}
	})
	return h
}

func FuzzNativeStopHandle(f *testing.F) {
	for _, v := range []string{"", "0", "4", "0004", "-1", "+4", "18446744073709551615", "9223372036854775807", "private", "4\x00"} {
		f.Add(v)
	}
	f.Fuzz(func(t *testing.T, raw string) {
		h, e := parseStopHandle(raw)
		if e != nil {
			if e != errNativeProtocol || h != 0 {
				t.Fatal("non-static rejection")
			}
			return
		}
		if h == 0 || h > ^uintptr(0)>>1 || strconv.FormatUint(uint64(h), 10) != raw {
			t.Fatal("accepted noncanonical or pseudo handle")
		}
	})
}

func TestNativeStopPreSignaledPreventsFirstRealHostClaim(t *testing.T) {
	h := stopTestEvent(t)
	stopKernel.NewProc("SetEvent").Call(h)
	ctx, close, e := nativeStopContext(context.Background(), strconv.FormatUint(uint64(h), 10))
	if e != nil {
		t.Fatal(e)
	}
	defer close()
	if ctx.Err() == nil {
		t.Fatal("pre-signaled event returned a live context")
	}
	claims := 0
	if err := runRefundLoop(ctx, func(context.Context, string) error { claims++; return nil }, claimToken); err != nil {
		t.Fatal(err)
	}
	if claims != 0 {
		t.Fatal("claim after pre-signaled stop")
	}
}
func TestNativeStopParse(t *testing.T) {
	for _, raw := range []string{"0", "-1", "+1", "01", " 1", "1 ", "18446744073709551615", "9223372036854775808", "1.0", "private-secret"} {
		if _, e := parseStopHandle(raw); e != errNativeProtocol {
			t.Fatalf("invalid accepted: %q", raw)
		}
	}
	if h, e := parseStopHandle("4"); e != nil || h != 4 {
		t.Fatal(h, e)
	}
}
func TestNativeStopInvalidClosedHandle(t *testing.T) {
	h, _, _ := stopKernel.NewProc("CreateEventW").Call(0, 1, 0, 0)
	if h == 0 {
		t.Fatal("create")
	}
	stopClose.Call(h)
	ctx, close, e := nativeStopContext(context.Background(), strconv.FormatUint(uint64(h), 10))
	if e != errNativeProtocol || context.Cause(ctx) != errNativeProtocol {
		t.Fatal("expected fixed failure")
	}
	if close() != nil {
		t.Fatal("cleanup")
	}
}
func TestNativeStopSignal(t *testing.T) {
	h := stopTestEvent(t)
	ctx, close, e := nativeStopContext(context.Background(), strconv.FormatUint(uint64(h), 10))
	if e != nil {
		t.Fatal(e)
	}
	defer close()
	stopKernel.NewProc("SetEvent").Call(h)
	select {
	case <-ctx.Done():
	case <-time.After(time.Second):
		t.Fatal("event ignored")
	}
	if context.Cause(ctx) != errNativeStop {
		t.Fatal(context.Cause(ctx))
	}
}
func TestNativeStopAlreadySignaled(t *testing.T) {
	h := stopTestEvent(t)
	stopKernel.NewProc("SetEvent").Call(h)
	ctx, close, e := nativeStopContext(context.Background(), strconv.FormatUint(uint64(h), 10))
	if e != nil {
		t.Fatal(e)
	}
	defer close()
	select {
	case <-ctx.Done():
	case <-time.After(time.Second):
		t.Fatal("event ignored")
	}
}
func TestNativeStopParentCause(t *testing.T) {
	parent, cancel := context.WithCancelCause(context.Background())
	cause := errors.New("parent cause")
	h := stopTestEvent(t)
	ctx, close, e := nativeStopContext(parent, strconv.FormatUint(uint64(h), 10))
	if e != nil {
		t.Fatal(e)
	}
	cancel(cause)
	if close() != nil || context.Cause(ctx) != cause {
		t.Fatal("cause or cleanup")
	}
}
func TestNativeStopEmpty(t *testing.T) {
	ctx, close, e := nativeStopContext(context.Background(), "")
	if e != nil || ctx.Err() != nil {
		t.Fatal(e)
	}
	if close() != nil || ctx.Err() != context.Canceled {
		t.Fatal("empty cleanup")
	}
}
func TestNativeStopConcurrentClose(t *testing.T) {
	h := stopTestEvent(t)
	_, close, e := nativeStopContext(context.Background(), strconv.FormatUint(uint64(h), 10))
	if e != nil {
		t.Fatal(e)
	}
	var wg sync.WaitGroup
	for range 24 {
		wg.Go(func() {
			if close() != nil {
				t.Error("cleanup")
			}
		})
	}
	wg.Wait()
}
func TestNativeStopDuplicateSurvivesSourceClosure(t *testing.T) {
	h, _, _ := stopKernel.NewProc("CreateEventW").Call(0, 1, 1, 0)
	if h == 0 {
		t.Fatal("create")
	}
	ctx, close, e := nativeStopContext(context.Background(), strconv.FormatUint(uint64(h), 10))
	if e != nil {
		stopClose.Call(h)
		t.Fatal(e)
	}
	stopClose.Call(h)
	defer close()
	select {
	case <-ctx.Done():
	case <-time.After(time.Second):
		t.Fatal("duplicate not independent")
	}
	if context.Cause(ctx) != errNativeStop {
		t.Fatal(context.Cause(ctx))
	}
}
func TestNativeStopRepeatedNoHandleGrowth(t *testing.T) {
	count := func() uint32 {
		var n uint32
		ok, _, _ := stopKernel.NewProc("GetProcessHandleCount").Call(^uintptr(0), uintptr(unsafe.Pointer(&n)))
		if ok == 0 {
			t.Fatal("count")
		}
		return n
	}
	h := stopTestEvent(t)
	raw := strconv.FormatUint(uint64(h), 10)
	_, close, _ := nativeStopContext(context.Background(), raw)
	close()
	before := count()
	for range 64 {
		_, close, e := nativeStopContext(context.Background(), raw)
		if e != nil || close() != nil {
			t.Fatal("cycle")
		}
	}
	if count() != before {
		t.Fatal("native handle growth")
	}
}

func TestNativeProtocolFailureRemainsHostFailure(t *testing.T) {
	ctx, cancel := context.WithCancelCause(context.Background())
	cancel(errNativeProtocol)
	if !errors.Is(nativeStopFailure(ctx), errNativeProtocol) {
		t.Fatal("native protocol failure suppressed")
	}
	ctx, cancel = context.WithCancelCause(context.Background())
	cancel(errNativeStop)
	if nativeStopFailure(ctx) != nil {
		t.Fatal("cooperative stop misclassified")
	}
}
