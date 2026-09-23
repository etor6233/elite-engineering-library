package whatsappbridge

import (
	"bufio"
	"bytes"
	"context"
	"elite.local/enterprise/internal/platform/identity"
	"encoding/json"
	"errors"
	"io"
	"net"
	"reflect"
	"sync"
	"testing"
	"time"
)

// AUTHORED transport tests use only local connections and synthetic worker data.
func reportFixture() StatusPollReport {
	return StatusPollReport{Outcome: "COMPLETED", Claimed: true, InsertedObservations: 1, Elapsed: time.Millisecond, NextDelay: time.Second}
}

func FuzzJSONStatusReporterReportInput(f *testing.F) {
	for _, outcome := range []string{"COMPLETED", "IDLE", "AUTH_UNAVAILABLE", "RECONCILE_REQUIRED", "TERMINAL_REVIEW", "RETRY_RECORDED", "PRIVATE\nforged", ""} {
		f.Add(outcome, int64(1), int64(time.Millisecond), int64(time.Second))
	}
	f.Add("COMPLETED", int64(-1), int64(0), int64(time.Second))
	f.Add("COMPLETED", int64(1), int64(-1), int64(time.Second))
	f.Add("COMPLETED", int64(1), int64(0), int64(10*time.Minute+1))
	f.Fuzz(func(t *testing.T, outcome string, count, elapsed, delay int64) {
		c := &reportTestConn{}
		r, _ := NewJSONStatusReporter(c)
		defer r.Close()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		p := StatusPollReport{Outcome: outcome, InsertedObservations: int(count), Elapsed: time.Duration(elapsed), NextDelay: time.Duration(delay)}
		err := r.Report(ctx, p)
		allowed := outcome == "COMPLETED" || outcome == "IDLE" || outcome == "AUTH_UNAVAILABLE" || outcome == "RECONCILE_REQUIRED" || outcome == "TERMINAL_REVIEW" || outcome == "RETRY_RECORDED"
		valid := allowed && p.InsertedObservations >= 0 && elapsed >= 0 && delay >= int64(time.Second) && delay <= int64(10*time.Minute)
		if !valid {
			if err != ErrStatusReport || c.writes != 0 || len(c.bytes) != 0 {
				t.Fatal("invalid report reached transport")
			}
			return
		}
		if err != nil || c.writes != 1 {
			t.Fatal("valid report not emitted")
		}
		assertReportLine(t, c.bytes, outcome)
		var row struct {
			Inserted int   `json:"inserted_observations"`
			Elapsed  int64 `json:"elapsed_ns"`
			Delay    int64 `json:"next_delay_ns"`
		}
		if json.Unmarshal(c.bytes, &row) != nil || row.Inserted != p.InsertedObservations || row.Elapsed != elapsed || row.Delay != delay {
			t.Fatal("numeric value changed")
		}
	})
}

func TestJSONStatusReporterCompletedContextCannotCloseNextReport(t *testing.T) {
	c := &reportTestConn{}
	r, _ := NewJSONStatusReporter(c)
	defer r.Close()
	for i := 0; i < 100; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		if err := r.Report(ctx, reportFixture()); err != nil {
			cancel()
			t.Fatal(err)
		}
		cancel()
	}
	if c.writes != 100 || c.closes != 0 {
		t.Fatal("late cancellation closed next report")
	}
}

func localReportConn(t *testing.T) (net.Conn, net.Conn) {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	a, err := net.DialTimeout("tcp", l.Addr().String(), time.Second)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { a.Close() })
	b, err := l.Accept()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { b.Close() })
	if err := b.SetReadDeadline(time.Now().Add(10 * time.Second)); err != nil {
		t.Fatal(err)
	}
	return a, b
}

func assertReportLine(t *testing.T, line []byte, outcome string) {
	t.Helper()
	var row map[string]json.RawMessage
	if len(line) > 512 || len(line) == 0 || line[len(line)-1] != '\n' || json.Unmarshal(line, &row) != nil {
		t.Fatal("invalid bounded JSON line")
	}
	keys := []string{"outcome", "claimed", "failure_recorded", "inserted_observations", "elapsed_ns", "next_delay_ns"}
	if len(row) != len(keys) {
		t.Fatal("schema drift")
	}
	for _, key := range keys {
		if _, ok := row[key]; !ok {
			t.Fatal("missing fixed field", key)
		}
	}
	var actual string
	if json.Unmarshal(row["outcome"], &actual) != nil || actual != outcome {
		t.Fatal("outcome mismatch")
	}
}

func TestJSONStatusReporterConcurrentTCPAndSchema(t *testing.T) {
	a, b := localReportConn(t)
	r, err := NewJSONStatusReporter(a)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	const count = 32
	errs := make(chan error, count)
	var group sync.WaitGroup
	for i := 0; i < count; i++ {
		group.Add(1)
		go func() { defer group.Done(); errs <- r.Report(ctx, reportFixture()) }()
	}
	reader := bufio.NewReader(b)
	for i := 0; i < count; i++ {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			t.Fatal(err)
		}
		assertReportLine(t, line, "COMPLETED")
	}
	group.Wait()
	for i := 0; i < count; i++ {
		if err := <-errs; err != nil {
			t.Fatal(err)
		}
	}
	if err := r.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := reader.ReadByte(); err != io.EOF {
		t.Fatal("extra bytes or unclosed stream", err)
	}
}

type reportTestConn struct {
	net.Conn
	mode           string
	writes, closes int
	deadline       time.Time
	bytes          []byte
}
type privateReportError struct{}

func (privateReportError) Error() string { panic("private error must not be formatted") }
func (c *reportTestConn) SetWriteDeadline(d time.Time) error {
	c.deadline = d
	if c.mode == "deadline" {
		return privateReportError{}
	}
	return nil
}
func (c *reportTestConn) Close() error { c.closes++; return nil }
func (c *reportTestConn) Write(b []byte) (int, error) {
	c.writes++
	c.bytes = append(c.bytes, b...)
	switch c.mode {
	case "zero":
		return 0, nil
	case "short":
		return len(b) - 1, nil
	case "negative":
		return -1, nil
	case "oversized":
		return len(b) + 1, nil
	case "error":
		return 0, privateReportError{}
	case "full-error":
		return len(b), privateReportError{}
	}
	return len(b), nil
}

func TestJSONStatusReporterUncertainWriteIsTerminal(t *testing.T) {
	for _, mode := range []string{"zero", "short", "negative", "oversized", "error", "full-error", "deadline"} {
		t.Run(mode, func(t *testing.T) {
			c := &reportTestConn{mode: mode}
			r, _ := NewJSONStatusReporter(c)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			if r.Report(ctx, reportFixture()) != ErrStatusReport || r.Report(ctx, reportFixture()) != ErrStatusReport {
				t.Fatal("uncertain write accepted")
			}
			want := 1
			if mode == "deadline" {
				want = 0
			}
			if c.writes != want || c.closes != 1 {
				t.Fatal("retry or unclosed connection", c.writes, c.closes)
			}
		})
	}
}

func TestJSONStatusReporterInputAndDeadlineBoundaries(t *testing.T) {
	if _, err := NewJSONStatusReporter(nil); err != ErrStatusReport {
		t.Fatal("nil connection")
	}
	if (&JSONStatusReporter{}).Report(context.Background(), reportFixture()) != ErrStatusReport {
		t.Fatal("zero reporter")
	}
	for _, mode := range []string{"outcome", "private-outcome", "negative-count", "negative-elapsed", "short-delay", "long-delay", "no-deadline", "cancelled", "nil-context"} {
		t.Run(mode, func(t *testing.T) {
			c := &reportTestConn{}
			r, _ := NewJSONStatusReporter(c)
			defer r.Close()
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			p := reportFixture()
			switch mode {
			case "outcome":
				p.Outcome = ""
			case "private-outcome":
				p.Outcome = "PRIVATE_TOKEN\nforged"
			case "negative-count":
				p.InsertedObservations = -1
			case "negative-elapsed":
				p.Elapsed = -1
			case "short-delay":
				p.NextDelay = time.Second - 1
			case "long-delay":
				p.NextDelay = 10*time.Minute + 1
			case "no-deadline":
				ctx = context.Background()
			case "cancelled":
				cancel()
			case "nil-context":
				ctx = nil
			}
			if r.Report(ctx, p) != ErrStatusReport || c.writes != 0 || c.closes != 0 {
				t.Fatal("invalid input touched transport")
			}
			valid, stop := context.WithTimeout(context.Background(), time.Second)
			defer stop()
			if r.Report(valid, reportFixture()) != nil {
				t.Fatal("invalid input poisoned stream")
			}
		})
	}
	for _, outcome := range []string{"AUTH_UNAVAILABLE", "IDLE", "RECONCILE_REQUIRED", "TERMINAL_REVIEW", "COMPLETED", "RETRY_RECORDED"} {
		c := &reportTestConn{}
		r, _ := NewJSONStatusReporter(c)
		ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
		p := reportFixture()
		p.Outcome = outcome
		p.NextDelay = 10 * time.Minute
		start := time.Now()
		err := r.Report(ctx, p)
		cancel()
		r.Close()
		if err != nil || c.deadline.After(time.Now().Add(3*time.Second)) || c.deadline.Before(start) {
			t.Fatal("deadline not capped")
		}
		assertReportLine(t, c.bytes, outcome)
	}
}

func TestJSONStatusReporterCancellationAndDeadlineInterruptBlockedWrite(t *testing.T) {
	for _, mode := range []string{"cancel", "deadline", "close"} {
		t.Run(mode, func(t *testing.T) {
			a, b := net.Pipe()
			defer b.Close()
			r, _ := NewJSONStatusReporter(a)
			defer r.Close()
			budget := 2 * time.Second
			if mode == "deadline" {
				budget = 50 * time.Millisecond
			}
			ctx, cancel := context.WithTimeout(context.Background(), budget)
			defer cancel()
			done := make(chan error, 1)
			go func() { done <- r.Report(ctx, reportFixture()) }()
			// Waiting for the first byte proves Write is active. The record remains
			// blocked because the receiver intentionally does not consume the rest.
			b.SetReadDeadline(time.Now().Add(time.Second))
			one := make([]byte, 1)
			if _, err := io.ReadFull(b, one); err != nil {
				t.Fatal(err)
			}
			if mode == "cancel" {
				cancel()
			}
			if mode == "close" {
				r.Close()
			}
			select {
			case err := <-done:
				if err != ErrStatusReport {
					t.Fatal("interrupted write accepted")
				}
			case <-time.After(time.Second):
				t.Fatal("write did not unblock")
			}
			fresh, stop := context.WithTimeout(context.Background(), time.Second)
			defer stop()
			if r.Report(fresh, reportFixture()) != ErrStatusReport {
				t.Fatal("truncated stream reused")
			}
		})
	}
}

func TestJSONStatusReporterWaitingCancellationDoesNotPoisonActiveWrite(t *testing.T) {
	a, b := net.Pipe()
	defer b.Close()
	r, _ := NewJSONStatusReporter(a)
	defer r.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	done := make(chan error, 1)
	go func() { done <- r.Report(ctx, reportFixture()) }()
	b.SetReadDeadline(time.Now().Add(time.Second))
	one := make([]byte, 1)
	if _, err := io.ReadFull(b, one); err != nil {
		t.Fatal(err)
	}
	wait, stop := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer stop()
	if r.Report(wait, reportFixture()) != ErrStatusReport {
		t.Fatal("expired waiter accepted")
	}
	reader := bufio.NewReader(b)
	rest, err := reader.ReadBytes('\n')
	if err != nil {
		t.Fatal(err)
	}
	assertReportLine(t, append(one, rest...), "COMPLETED")
	if err := <-done; err != nil {
		t.Fatal("waiter poisoned active report")
	}
}

func TestStatusHostConcreteReporterFailureStopsWithoutRetry(t *testing.T) {
	c := &reportTestConn{mode: "short"}
	r, _ := NewJSONStatusReporter(c)
	defer r.Close()
	calls := 0
	err := runStatusLoop(context.Background(), func(context.Context, identity.Principal) (StatusWorkResult, error) {
		calls++
		return StatusWorkResult{Completed: true}, nil
	}, statusSourceFunc(func(context.Context) (identity.Principal, error) { return identity.Principal{}, nil }), r, time.Second, func(context.Context, time.Duration) error { t.Fatal("wait after report loss"); return nil })
	if err != ErrStatusHost || calls != 1 || c.writes != 1 || c.closes != 1 {
		t.Fatal("host did not stop")
	}
}

func TestStatusHostRealWorkerJSONTransport(t *testing.T) {
	router, _, pool, p, _, _, _ := routerFixture(t, statusBatch("delivered", "1603086314", nil))
	w, err := NewStatusWorker(router, "json-host-test", time.Minute, time.Second)
	if err != nil {
		t.Fatal(err)
	}
	a, b := localReportConn(t)
	reporter, _ := NewJSONStatusReporter(a)
	defer reporter.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() {
		done <- w.Run(ctx, statusSourceFunc(func(context.Context) (identity.Principal, error) { return p, nil }), reporter, time.Second)
	}()
	line, err := bufio.NewReader(b).ReadBytes('\n')
	if err != nil {
		t.Fatal(err)
	}
	cancel()
	assertReportLine(t, line, "COMPLETED")
	if !bytes.Contains(line, []byte(`"inserted_observations":1`)) {
		t.Fatal("observation absent")
	}
	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatal(err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("worker did not stop")
	}
	var n int
	if err := pool.QueryRow(context.Background(), `select count(*) from platform.job where tenant_id=$1 and completed_at is not null and attempts=1`, p.TenantID).Scan(&n); err != nil || n != 1 {
		t.Fatal("job completion", err, n)
	}
}

func TestStatusHostReportLossAfterCommitDoesNotReplayCompletedJob(t *testing.T) {
	router, _, pool, p, _, _, _ := routerFixture(t, statusBatch("delivered", "1603086314", nil))
	w, err := NewStatusWorker(router, "loss-host-test", time.Minute, time.Second)
	if err != nil {
		t.Fatal(err)
	}
	a, b := net.Pipe()
	reporter, _ := NewJSONStatusReporter(a)
	defer reporter.Close()
	defer b.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	done := make(chan error, 1)
	source := statusSourceFunc(func(context.Context) (identity.Principal, error) { return p, nil })
	go func() { done <- w.Run(ctx, source, reporter, time.Second) }()
	// The transaction has committed before the first report byte. Destroy the
	// receiver mid-record to reproduce ambiguous telemetry delivery after commit.
	b.SetReadDeadline(time.Now().Add(4 * time.Second))
	one := make([]byte, 1)
	if _, err := io.ReadFull(b, one); err != nil {
		t.Fatal(err)
	}
	b.Close()
	select {
	case err := <-done:
		if err != ErrStatusHost {
			t.Fatal("report loss not propagated")
		}
	case <-ctx.Done():
		t.Fatal("host did not stop")
	}

	// A supervisor may create a new host and connection, but must not replay the
	// failed report or enqueue the already-committed job again.
	restarted, err := NewStatusWorker(router, "restarted-host-test", time.Minute, time.Second)
	if err != nil {
		t.Fatal(err)
	}
	left, right := localReportConn(t)
	next, _ := NewJSONStatusReporter(left)
	defer next.Close()
	nextCtx, stop := context.WithCancel(context.Background())
	defer stop()
	go func() { done <- restarted.Run(nextCtx, source, next, time.Second) }()
	line, err := bufio.NewReader(right).ReadBytes('\n')
	if err != nil {
		t.Fatal(err)
	}
	stop()
	assertReportLine(t, line, "IDLE")
	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatal(err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("restarted host did not stop")
	}
	var complete, observed, audited int
	err = pool.QueryRow(context.Background(), `select
 (select count(*) from platform.job where tenant_id=$1 and completed_at is not null and attempts=1),
 (select count(*) from communication.whatsapp_status_observation where tenant_id=$1),
 (select count(*) from audit.event where tenant_id=$1 and action='whatsapp.status_job.completed')`, p.TenantID).Scan(&complete, &observed, &audited)
	if err != nil || complete != 1 || observed != 1 || audited != 1 {
		t.Fatal("report loss replayed durable work", err, complete, observed, audited)
	}
}

type statusSourceFunc func(context.Context) (identity.Principal, error)

func (f statusSourceFunc) Resolve(ctx context.Context) (identity.Principal, error) { return f(ctx) }

type statusReportFunc func(context.Context, StatusPollReport) error

func (f statusReportFunc) Report(ctx context.Context, r StatusPollReport) error { return f(ctx, r) }

func TestStatusHostBackoffIdentityAndReport(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	calls, auth := 0, 0
	var reports []StatusPollReport
	source := statusSourceFunc(func(c context.Context) (identity.Principal, error) {
		auth++
		if _, ok := c.Deadline(); !ok {
			t.Fatal("auth without deadline")
		}
		if auth == 1 {
			return identity.Principal{}, errors.New("sensitive auth error")
		}
		return identity.Principal{Subject: "verified"}, nil
	})
	process := func(c context.Context, p identity.Principal) (StatusWorkResult, error) {
		calls++
		if p.Subject != "verified" {
			t.Fatal("wrong source identity")
		}
		switch calls {
		case 1:
			return StatusWorkResult{Claimed: true, FailureRecorded: true}, ErrStatusWork
		case 2:
			return StatusWorkResult{Terminal: true, FailureRecorded: true}, ErrStatusWork
		case 3:
			return StatusWorkResult{Claimed: true, InsertedObservations: 1}, ErrStatusWork
		case 4:
			return StatusWorkResult{Claimed: true, Completed: true}, nil
		default:
			return StatusWorkResult{}, nil
		}
	}
	reporter := statusReportFunc(func(c context.Context, r StatusPollReport) error {
		if _, ok := c.Deadline(); !ok {
			t.Fatal("report without deadline")
		}
		reports = append(reports, r)
		return nil
	})
	waits := 0
	err := runStatusLoop(ctx, process, source, reporter, time.Second, func(c context.Context, d time.Duration) error {
		waits++
		if waits == 6 {
			cancel()
		}
		return nil
	})
	if !errors.Is(err, context.Canceled) || auth != 6 || calls != 5 {
		t.Fatal(err, auth, calls)
	}
	want := []string{"AUTH_UNAVAILABLE", "RETRY_RECORDED", "TERMINAL_REVIEW", "RECONCILE_REQUIRED", "COMPLETED", "IDLE"}
	for i, r := range reports {
		if r.Outcome != want[i] {
			t.Fatal(i, r)
		}
		if i < 4 {
			cap := time.Second * time.Duration(1<<(i+1))
			if r.NextDelay < cap/2 || r.NextDelay > cap {
				t.Fatal("backoff bounds", r)
			}
		} else if r.NextDelay != time.Second {
			t.Fatal("success did not reset", r)
		}
	}
	// Schema allowlist prevents adding arbitrary labels/raw errors accidentally.
	typ := reflect.TypeOf(StatusPollReport{})
	fields := []string{"Outcome", "Claimed", "FailureRecorded", "InsertedObservations", "Elapsed", "NextDelay"}
	if typ.NumField() != len(fields) {
		t.Fatal("report surface changed")
	}
	for i, n := range fields {
		if typ.Field(i).Name != n {
			t.Fatal("unexpected telemetry field")
		}
	}
}

func TestStatusHostReportingFailureStopsWithoutRetry(t *testing.T) {
	calls := 0
	err := runStatusLoop(context.Background(), func(context.Context, identity.Principal) (StatusWorkResult, error) {
		calls++
		return StatusWorkResult{Completed: true}, nil
	}, statusSourceFunc(func(context.Context) (identity.Principal, error) { return identity.Principal{}, nil }), statusReportFunc(func(context.Context, StatusPollReport) error { return errors.New("secret reporter transport") }), time.Second, func(context.Context, time.Duration) error { t.Fatal("wait after report loss"); return nil })
	if err != ErrStatusHost || calls != 1 {
		t.Fatal(err, calls)
	}
}

func TestStatusHostShutdownAndBoundedBackoff(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	n := 0
	err := runStatusLoop(ctx, func(context.Context, identity.Principal) (StatusWorkResult, error) {
		t.Fatal("work after auth failure")
		return StatusWorkResult{}, nil
	}, statusSourceFunc(func(context.Context) (identity.Principal, error) { n++; return identity.Principal{}, ErrStatusHost }), statusReportFunc(func(c context.Context, r StatusPollReport) error {
		if n >= 4 && (r.NextDelay < 5*time.Minute || r.NextDelay > 10*time.Minute) {
			t.Fatal("backoff cap", r)
		}
		return nil
	}), time.Minute, func(context.Context, time.Duration) error {
		if n == 15 {
			cancel()
		}
		return nil
	})
	if !errors.Is(err, context.Canceled) || n != 15 {
		t.Fatal(err, n)
	}
	start := time.Now()
	if !errors.Is(waitStatusPoll(ctx, time.Hour), context.Canceled) || time.Since(start) > time.Second {
		t.Fatal("shutdown did not interrupt wait")
	}
	if (&StatusWorker{}).Run(ctx, nil, nil, time.Second) != ErrStatusHost {
		t.Fatal("zero host accepted")
	}
}

func TestStatusHostRealWorkerCompletesAndStops(t *testing.T) {
	r, _, pool, p, _, _, _ := routerFixture(t, statusBatch("delivered", "1603086314", nil))
	w, err := NewStatusWorker(r, "host-test", time.Minute, time.Second)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var report StatusPollReport
	err = w.Run(ctx, statusSourceFunc(func(context.Context) (identity.Principal, error) { return p, nil }), statusReportFunc(func(c context.Context, r StatusPollReport) error { report = r; cancel(); return nil }), time.Second)
	if !errors.Is(err, context.Canceled) || report.Outcome != "COMPLETED" || !report.Claimed || report.InsertedObservations != 1 {
		t.Fatal(err, report)
	}
	var n int
	if err := pool.QueryRow(context.Background(), `select count(*) from platform.job where tenant_id=$1 and completed_at is not null and attempts=1`, p.TenantID).Scan(&n); err != nil || n != 1 {
		t.Fatal(err, n)
	}
}

func TestStatusHostRejectsLateAuthenticationAndReport(t *testing.T) {
	for _, mode := range []string{"authentication", "report"} {
		t.Run(mode, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			calls := 0
			source := statusSourceFunc(func(c context.Context) (identity.Principal, error) {
				if mode == "authentication" {
					<-c.Done()
				}
				return identity.Principal{Subject: "late-identity"}, nil
			})
			reporter := statusReportFunc(func(c context.Context, r StatusPollReport) error {
				if mode == "report" {
					<-c.Done()
				} else {
					if r.Outcome != "AUTH_UNAVAILABLE" {
						t.Fatal("late authority used", r)
					}
					cancel()
				}
				return nil
			})
			err := runStatusLoop(ctx, func(context.Context, identity.Principal) (StatusWorkResult, error) {
				calls++
				return StatusWorkResult{}, nil
			}, source, reporter, time.Second, func(context.Context, time.Duration) error { t.Fatal("retry after terminal test condition"); return nil })
			if mode == "authentication" {
				if calls != 0 || !errors.Is(err, context.Canceled) {
					t.Fatal(calls, err)
				}
			} else if calls != 1 || err != ErrStatusHost {
				t.Fatal(calls, err)
			}
		})
	}
}
