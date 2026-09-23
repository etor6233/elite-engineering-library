//go:build windows

package main

// AUTHORED qualification only. The launcher supplies a trusted manual-reset
// event; this is neither object-type authentication nor hostile-worker isolation.
import (
	"context"
	"errors"
	"strconv"
	"sync"
	"syscall"
	"unsafe"
)

var errNativeStop = errors.New("NATIVE_STOP_REQUESTED")
var errNativeProtocol = errors.New("NATIVE_STOP_PROTOCOL_FAILED")
var stopKernel = syscall.NewLazyDLL("kernel32.dll")
var stopDuplicate = stopKernel.NewProc("DuplicateHandle")
var stopWait = stopKernel.NewProc("WaitForSingleObject")
var stopClose = stopKernel.NewProc("CloseHandle")

func parseStopHandle(raw string) (uintptr, error) {
	n, err := strconv.ParseUint(raw, 10, strconv.IntSize)
	if err != nil || n == 0 || n > uint64(^uintptr(0)>>1) || strconv.FormatUint(n, 10) != raw {
		return 0, errNativeProtocol
	}
	return uintptr(n), nil
}

// Duplicate before waiting. Cleanup joins the waiter BEFORE closing its handle:
// CloseHandle during a pending Windows wait has undefined behavior.
// Empty raw retains the ordinary parent-context lifecycle. Cancellation stops
// claims through the real host loop; it does not guarantee a domain commit.
func nativeStopContext(parent context.Context, raw string) (context.Context, func() error, error) {
	ctx, cancel := context.WithCancelCause(parent)
	if raw == "" {
		return ctx, func() error { cancel(context.Canceled); return nil }, nil
	}
	source, err := parseStopHandle(raw)
	if err != nil {
		cancel(err)
		return ctx, func() error { return nil }, err
	}
	var owned uintptr
	ok, _, _ := stopDuplicate.Call(^uintptr(0), source, ^uintptr(0), uintptr(unsafe.Pointer(&owned)), 0x100000, 0, 0)
	if ok == 0 {
		cancel(errNativeProtocol)
		return ctx, func() error { return nil }, errNativeProtocol
	}
	done, joined := make(chan struct{}), make(chan struct{})
	var once sync.Once
	var closeErr error
	cleanup := func() error {
		once.Do(func() {
			close(done)
			<-joined
			ok, _, _ := stopClose.Call(owned)
			if ok == 0 {
				closeErr = errNativeProtocol
			}
			cancel(context.Canceled)
		})
		return closeErr
	}
	// A pre-signaled event must not expose a live context to the first claim.
	// The trusted launcher supplies a manual-reset event, so this read is not
	// destructive. Other waitable object types are outside the protocol.
	initial, _, _ := stopWait.Call(owned, 0)
	if initial != 258 {
		close(joined)
		if initial == 0 {
			cancel(errNativeStop)
			return ctx, cleanup, nil
		}
		cancel(errNativeProtocol)
		_ = cleanup()
		return ctx, cleanup, errNativeProtocol
	}
	go func() {
		defer close(joined)
		for {
			select {
			case <-done:
				return
			case <-ctx.Done():
				return
			default:
			}
			status, _, _ := stopWait.Call(owned, 20)
			switch status {
			case 0:
				cancel(errNativeStop)
				return
			case 258:
			default:
				cancel(errNativeProtocol)
				return
			}
		}
	}()
	return ctx, cleanup, nil
}

func nativeStopFailure(ctx context.Context) error {
	if errors.Is(context.Cause(ctx), errNativeProtocol) {
		return errNativeProtocol
	}
	return nil
}
