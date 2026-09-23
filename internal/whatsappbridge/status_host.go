package whatsappbridge

// AUTHORED opt-in host for the existing worker; no daemon, identity, provider
// subscription or monitoring destination is activated by materialization.
import (
	"context"
	"encoding/json"
	"errors"
	"math/rand/v2"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"elite.local/enterprise/internal/platform/identity"
)

var ErrStatusHost = errors.New("whatsappbridge: status host configuration or reporting unavailable")

var ErrStatusReport = errors.New("whatsappbridge: status report unavailable or delivery uncertain")

// JSONStatusReporter owns one already-established connection exclusively: no
// other writer or deadline owner may use it. The caller establishes authenticated
// TLS and selects its approved collector; this constructor never dials. Success
// means a complete local Write, not collector acknowledgement or durable retention.
// Do not copy a reporter after construction. Use a fresh connection/reporter after
// any uncertain write; never retry an uncertain report automatically.
type JSONStatusReporter struct {
	conn      net.Conn
	gate      chan struct{}
	closed    atomic.Bool
	closeOnce sync.Once
	closeErr  error
}

var _ StatusReporter = (*JSONStatusReporter)(nil)

func NewJSONStatusReporter(conn net.Conn) (*JSONStatusReporter, error) {
	if conn == nil {
		return nil, ErrStatusReport
	}
	r := &JSONStatusReporter{conn: conn, gate: make(chan struct{}, 1)}
	r.gate <- struct{}{}
	return r, nil
}

// Close interrupts pending network I/O under the net.Conn contract. The
// connection is terminal even if Close fails. Underlying error text is private.
func (r *JSONStatusReporter) Close() error {
	if r == nil || r.conn == nil {
		return ErrStatusReport
	}
	r.closeOnce.Do(func() { r.closed.Store(true); r.closeErr = r.conn.Close() })
	if r.closeErr != nil {
		return ErrStatusReport
	}
	return nil
}

// Report emits one bounded JSON line with six fixed fields. A deadline is
// mandatory and is capped at three seconds, including serialization contention.
// Deadline/cancellation/write failures close the stream to prevent appending to
// a possibly truncated record. Invalid input and cancellation while waiting for
// another writer do not touch the connection. Supplied net.Conn implementations
// must honor concurrent Close and SetWriteDeadline, as required by Go's contract.
func (r *JSONStatusReporter) Report(ctx context.Context, p StatusPollReport) error {
	if r == nil || r.conn == nil || r.gate == nil || ctx == nil || r.closed.Load() {
		return ErrStatusReport
	}
	switch p.Outcome {
	case "AUTH_UNAVAILABLE", "IDLE", "RECONCILE_REQUIRED", "TERMINAL_REVIEW", "COMPLETED", "RETRY_RECORDED":
	default:
		return ErrStatusReport
	}
	if p.InsertedObservations < 0 || p.Elapsed < 0 || p.NextDelay < time.Second || p.NextDelay > 10*time.Minute {
		return ErrStatusReport
	}
	deadline, ok := ctx.Deadline()
	if !ok || ctx.Err() != nil {
		return ErrStatusReport
	}
	if cap := time.Now().Add(3 * time.Second); deadline.After(cap) {
		deadline = cap
	}
	ctx, cancel := context.WithDeadline(ctx, deadline)
	defer cancel()
	select {
	case <-ctx.Done():
		return ErrStatusReport
	case <-r.gate:
	}
	defer func() { r.gate <- struct{}{} }()
	if ctx.Err() != nil || r.closed.Load() {
		return ErrStatusReport
	}
	line, err := json.Marshal(struct {
		Outcome              string `json:"outcome"`
		Claimed              bool   `json:"claimed"`
		FailureRecorded      bool   `json:"failure_recorded"`
		InsertedObservations int    `json:"inserted_observations"`
		ElapsedNS            int64  `json:"elapsed_ns"`
		NextDelayNS          int64  `json:"next_delay_ns"`
	}{p.Outcome, p.Claimed, p.FailureRecorded, p.InsertedObservations, int64(p.Elapsed), int64(p.NextDelay)})
	if err != nil {
		return ErrStatusReport
	}
	line = append(line, '\n')
	if r.conn.SetWriteDeadline(deadline) != nil {
		_ = r.Close()
		return ErrStatusReport
	}
	done := make(chan struct{})
	stop := context.AfterFunc(ctx, func() { _ = r.Close(); close(done) })
	n, err := r.conn.Write(line)
	// stop does not wait for an already-running callback. Join it before releasing
	// the gate so it cannot close a connection during the next report.
	if !stop() {
		<-done
	}
	if err != nil || n != len(line) || ctx.Err() != nil || r.closed.Load() {
		_ = r.Close()
		return ErrStatusReport
	}
	return nil
}

// Resolve must authenticate/revalidate a real service identity and current
// grants on every poll. It must honor context; never derive grants from payload.
type StatusPrincipalSource interface {
	Resolve(context.Context) (identity.Principal, error)
}

// StatusPollReport deliberately has no contact, payload, token, provider ID,
// tenant label or arbitrary error text. A host adds only approved static labels.
type StatusPollReport struct {
	Outcome                  string
	Claimed, FailureRecorded bool
	InsertedObservations     int
	Elapsed, NextDelay       time.Duration
}

// Report must honor its deadline and return an error when delivery fails. The
// host stops then, so loss of its reporting channel is never silently ignored.
type StatusReporter interface {
	Report(context.Context, StatusPollReport) error
}

// Run blocks until cancellation or reporting failure. Poll interval is explicit
// and bounded; each loop takes at most one job and revalidates identity. Failed
// polls back off with jitter, capped at ten minutes. No send retry/redrive.
func (w *StatusWorker) Run(ctx context.Context, source StatusPrincipalSource, reporter StatusReporter, interval time.Duration) error {
	if w == nil || w.router == nil || w.jobs == nil || source == nil || reporter == nil || interval < time.Second || interval > time.Minute {
		return ErrStatusHost
	}
	return runStatusLoop(ctx, w.ProcessOnce, source, reporter, interval, waitStatusPoll)
}

func waitStatusPoll(ctx context.Context, delay time.Duration) error {
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func runStatusLoop(ctx context.Context, process func(context.Context, identity.Principal) (StatusWorkResult, error), source StatusPrincipalSource, reporter StatusReporter, interval time.Duration, wait func(context.Context, time.Duration) error) error {
	failures := 0
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		started := time.Now()
		authCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		principal, err := source.Resolve(authCtx)
		if err == nil {
			err = authCtx.Err()
		}
		cancel()
		result := StatusWorkResult{}
		outcome := "AUTH_UNAVAILABLE"
		if err == nil && ctx.Err() == nil {
			result, err = process(ctx, principal)
			outcome = "IDLE"
			if err != nil {
				outcome = "RECONCILE_REQUIRED"
			}
			if result.Terminal {
				outcome = "TERMINAL_REVIEW"
			} else if result.Completed && err == nil {
				outcome = "COMPLETED"
			} else if err != nil && result.FailureRecorded {
				outcome = "RETRY_RECORDED"
			}
		}
		// Report the completed attempt even if shutdown arrived during processing.
		// WithoutCancel is solely for bounded telemetry, never another job attempt.
		if err != nil {
			if failures < 10 {
				failures++
			}
		} else {
			failures = 0
		}
		delay := interval
		if failures > 0 {
			ceiling := interval * time.Duration(1<<failures)
			if ceiling > 10*time.Minute {
				ceiling = 10 * time.Minute
			}
			delay = ceiling/2 + time.Duration(rand.Int64N(int64(ceiling/2)+1))
		}
		reportCtx, finish := context.WithTimeout(context.WithoutCancel(ctx), 3*time.Second)
		reportErr := reporter.Report(reportCtx, StatusPollReport{Outcome: outcome, Claimed: result.Claimed, FailureRecorded: result.FailureRecorded, InsertedObservations: result.InsertedObservations, Elapsed: time.Since(started), NextDelay: delay})
		if reportErr == nil {
			reportErr = reportCtx.Err()
		}
		finish()
		if reportErr != nil {
			return ErrStatusHost
		}
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := wait(ctx, delay); err != nil {
			return err
		}
	}
}
