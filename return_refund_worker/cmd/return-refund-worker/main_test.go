package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"testing"
	"time"

	"elite.local/return-refund-worker/internal/refundworker"
)

type privateHostError struct{ formatted *bool }

func (e privateHostError) Error() string {
	*e.formatted = true
	return "PRIVATE-SYNTHETIC-STEP token=synthetic-secret"
}

func TestRefundHostStepErrorsDoNotFormatPrivateCause(t *testing.T) {
	var output bytes.Buffer
	writer, flags, prefix := log.Writer(), log.Flags(), log.Prefix()
	log.SetOutput(&output)
	log.SetFlags(0)
	log.SetPrefix("")
	defer func() { log.SetOutput(writer); log.SetFlags(flags); log.SetPrefix(prefix) }()
	formatted := false
	logStepError(privateHostError{&formatted})
	if formatted {
		t.Error("host formatted the private error cause")
	}
	if strings.Contains(output.String(), "PRIVATE-SYNTHETIC") || strings.Contains(output.String(), "synthetic-secret") {
		t.Error("private marker escaped into step log")
	}
	if output.String() != "return refund step failed: REFUND_STEP_FAILED\n" {
		t.Error("step failure requires the fixed operator code")
	}
	output.Reset()
	for _, err := range []error{nil, refundworker.ErrNoWork, fmt.Errorf("synthetic wrapper: %w", refundworker.ErrNoWork)} {
		logStepError(err)
	}
	if output.Len() != 0 {
		t.Error("idle results emitted failure logs")
	}
}

func TestRefundHostStartupErrorRedaction(t *testing.T) {
	if os.Getenv("ELITE_REFUND_TEST_CHILD") == "1" {
		main()
		return
	}
	for _, tc := range []struct{ name, dsn string }{
		{"missing_configuration", ""},
		{"invalid_private_dsn", "postgres://synthetic@127.0.0.1:invalid/PRIVATE-SYNTHETIC-DATABASE"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestRefundHostStartupErrorRedaction$")
			for _, key := range []string{"SystemRoot", "WINDIR", "TEMP", "TMP", "PATH"} {
				if value, ok := os.LookupEnv(key); ok {
					cmd.Env = append(cmd.Env, key+"="+value)
				}
			}
			cmd.Env = append(cmd.Env, "ELITE_REFUND_TEST_CHILD=1", "DATABASE_URL="+tc.dsn, "REFUND_WORKER_ID=synthetic-worker")
			output, err := cmd.CombinedOutput()
			var exit *exec.ExitError
			if !errors.As(err, &exit) || exit.ExitCode() != 1 {
				t.Fatal("failed startup must exit1")
			}
			if strings.Contains(string(output), "PRIVATE-SYNTHETIC") {
				t.Error("startup disclosed private configuration")
			}
			if !strings.Contains(string(output), "REFUND_HOST_FAILED") {
				t.Error("startup requires fixed operator code")
			}
		})
	}
}

type privateSinkError struct{}

func (privateSinkError) Error() string { panic("sink error must remain private") }

type failedHostWriter struct {
	mode  string
	calls int
}

func (w *failedHostWriter) Write(p []byte) (int, error) {
	w.calls++
	switch w.mode {
	case "error":
		return 0, privateSinkError{}
	case "full-error":
		return len(p), privateSinkError{}
	case "zero":
		return 0, nil
	case "short":
		return len(p) - 1, nil
	case "negative":
		return -1, nil
	case "oversized":
		return len(p) + 1, nil
	}
	return len(p), nil
}

func TestRefundHostReportDeliveryFailureIsVisible(t *testing.T) {
	for _, mode := range []string{"error", "full-error", "zero", "short", "negative", "oversized"} {
		t.Run(mode, func(t *testing.T) {
			previous := log.Writer()
			defer log.SetOutput(previous)
			w := &failedHostWriter{mode: mode}
			log.SetOutput(w)
			// Reflection permits the same behavioral probe against the old void
			// function and its new error-returning contract without editing baseline.
			result := reflect.ValueOf(logStepError).Call([]reflect.Value{reflect.ValueOf(errors.New("PRIVATE-STEP"))})
			if len(result) != 1 || result[0].IsNil() {
				t.Fatal("failed reporting has no observable error for host")
			}
			err, ok := result[0].Interface().(error)
			if !ok || err.Error() != "return refund report unavailable: REFUND_REPORT_FAILED" {
				t.Fatal("non-static report failure")
			}
			if w.calls != 1 {
				t.Fatal("uncertain report was retried")
			}
		})
	}
}

var _ io.Writer = (*failedHostWriter)(nil)

func TestRefundHostLoopStopsBeforeAnotherClaimOnReportLoss(t *testing.T) {
	for _, mode := range []string{"error", "full-error", "zero", "short", "negative", "oversized"} {
		t.Run(mode, func(t *testing.T) {
			previous := log.Writer()
			defer log.SetOutput(previous)
			w := &failedHostWriter{mode: mode}
			log.SetOutput(w)
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			steps, tokens := 0, 0
			formatted := false
			err := runRefundLoop(ctx, func(context.Context, string) error { steps++; return privateHostError{&formatted} }, func() (string, error) { tokens++; return "synthetic-token", nil })
			if err != errRefundReport || steps != 1 || tokens != 1 || w.calls != 1 || formatted {
				t.Fatal("report failure did not stop the actual host loop privately")
			}
		})
	}
}

func TestRefundHostLoopRetainsHealthyStepAndIdlePolicy(t *testing.T) {
	var output bytes.Buffer
	writer, flags, prefix := log.Writer(), log.Flags(), log.Prefix()
	defer func() { log.SetOutput(writer); log.SetFlags(flags); log.SetPrefix(prefix) }()
	log.SetOutput(&output)
	log.SetFlags(0)
	log.SetPrefix("approved-static: ")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	steps, tokens := 0, 0
	formatted := false
	err := runRefundLoop(ctx, func(context.Context, string) error {
		steps++
		if steps == 1 {
			return privateHostError{&formatted}
		}
		cancel()
		return refundworker.ErrNoWork
	}, func() (string, error) { tokens++; return "synthetic-token", nil })
	if err != nil || steps != 2 || tokens != 2 || formatted {
		t.Fatal("healthy reporting changed step/idle behavior")
	}
	if output.String() != "approved-static: return refund step failed: REFUND_STEP_FAILED\n" {
		t.Fatal("fixed message or configured prefix changed")
	}
}

func TestRefundHostLoopCancellationAndTokenFailure(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if runRefundLoop(ctx, func(context.Context, string) error { t.Fatal("step after cancellation"); return nil }, func() (string, error) { t.Fatal("token after cancellation"); return "", nil }) != nil {
		t.Fatal("cancellation is not normal shutdown")
	}
	private := errors.New("PRIVATE-TOKEN-CAUSE")
	if err := runRefundLoop(context.Background(), func(context.Context, string) error { t.Fatal("step after failed token"); return nil }, func() (string, error) { return "", private }); err != private {
		t.Fatal("token failure not returned")
	}
}

func TestRefundHostRealClosedFileStopsLoop(t *testing.T) {
	p, err := os.CreateTemp(t.TempDir(), "closed-host-log-")
	if err != nil {
		t.Fatal(err)
	}
	name := p.Name()
	if p.Close() != nil {
		t.Fatal("close fixture")
	}
	previous := log.Writer()
	defer log.SetOutput(previous)
	log.SetOutput(p)
	steps := 0
	err = runRefundLoop(context.Background(), func(context.Context, string) error { steps++; return errors.New("PRIVATE-STEP") }, func() (string, error) { return "synthetic-token", nil })
	if err != errRefundReport || steps != 1 {
		t.Fatal("closed file did not stop host")
	}
	data, err := os.ReadFile(name)
	if err != nil || len(data) != 0 {
		t.Fatal("closed sink unexpectedly changed")
	}
}

func TestRefundHostLostReportsExitNonzeroWithoutPrivateCause(t *testing.T) {
	if mode := os.Getenv("ELITE_REFUND_TEST_REPORT_CHILD"); mode != "" {
		if mode == "healthy" {
			exitOnHostFailure(nil)
			_, _ = os.Stdout.WriteString("HOST_NORMAL_RETURN\n")
			return
		}
		// The same loop and exit adapter are called by production run/main. This
		// child substitutes only the already-completed step and failed log sink.
		log.SetOutput(&failedHostWriter{mode: "error"})
		err := runRefundLoop(context.Background(), func(context.Context, string) error { return errors.New("PRIVATE-STEP") }, func() (string, error) { return "synthetic-token", nil })
		exitOnHostFailure(err)
		_, _ = os.Stdout.WriteString("UNEXPECTED_HOST_CONTINUATION\n")
		return
	}
	for _, mode := range []string{"failed", "healthy"} {
		t.Run(mode, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestRefundHostLostReportsExitNonzeroWithoutPrivateCause$")
			for _, key := range []string{"SystemRoot", "WINDIR", "TEMP", "TMP", "PATH"} {
				if v, ok := os.LookupEnv(key); ok {
					cmd.Env = append(cmd.Env, key+"="+v)
				}
			}
			cmd.Env = append(cmd.Env, "ELITE_REFUND_TEST_REPORT_CHILD="+mode)
			output, err := cmd.CombinedOutput()
			if strings.Contains(string(output), "PRIVATE") || strings.Contains(string(output), "UNEXPECTED_HOST_CONTINUATION") {
				t.Fatal("host continued or private text escaped")
			}
			if mode == "failed" {
				var exit *exec.ExitError
				if !errors.As(err, &exit) || exit.ExitCode() != 1 {
					t.Fatal("report loss must exit1 even when stderr cannot report")
				}
			} else if err != nil || !strings.Contains(string(output), "HOST_NORMAL_RETURN") {
				t.Fatal("normal exit changed")
			}
		})
	}
}

type countHostWriter struct {
	count int
	fail  bool
	calls int
}

func (w *countHostWriter) Write([]byte) (int, error) {
	w.calls++
	if w.fail {
		return w.count, privateSinkError{}
	}
	return w.count, nil
}

func FuzzRefundHostCheckedLogWrite(f *testing.F) {
	for _, n := range []int{-1, 0, 1, 6, 7, 8, 1024} {
		f.Add(n, false)
		f.Add(n, true)
	}
	f.Fuzz(func(t *testing.T, n int, failed bool) {
		w := &countHostWriter{count: n, fail: failed}
		actual, err := (checkedRefundLogWriter{writer: w}).Write([]byte("report\n"))
		if w.calls != 1 {
			t.Fatal("log write repeated")
		}
		if !failed && n == 7 {
			if err != nil || actual != 7 {
				t.Fatal("complete write rejected")
			}
			return
		}
		if err != errRefundReport || actual < 0 || actual > 7 {
			t.Fatal("invalid delivery accepted or count escaped")
		}
	})
}

func TestRefundHostIdleDoesNotTouchFailedSink(t *testing.T) {
	w := &failedHostWriter{mode: "error"}
	previous := log.Writer()
	defer log.SetOutput(previous)
	log.SetOutput(w)
	for _, err := range []error{nil, refundworker.ErrNoWork} {
		if logStepError(err) != nil {
			t.Fatal("idle report error")
		}
	}
	if w.calls != 0 {
		t.Fatal("idle wrote to sink")
	}
	if n, err := (checkedRefundLogWriter{}).Write([]byte("report")); n != 0 || err != errRefundReport {
		t.Fatal("nil writer accepted")
	}
}

var _ io.Writer = (*countHostWriter)(nil)
