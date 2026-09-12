# Python Meta WhatsApp Cloud Adapter

## 1. Metadata

```yaml
pack_id: "PYTHON-META-WHATSAPP-CLOUD-ADAPTER"
pack_version: "0.14.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Referencias Meta exactas y adaptación scoped; composición Go AUTHORED de envío aprobado, inbox/job, worker atómico, host opt-in con reportador JSON acotado y lectura autorizada de historial. Vista de agenda, recuperación GET y ayuda versionada probadas localmente con el pack de portales; sin identidad de servicio, alertas, redrive ni suscripción live desplegados."
stacks: ["CPython 3.12+ stdlib", "Go 1.26.7 para el bridge opcional", "Meta WhatsApp Cloud API", "fbsamples signed commit de70ee90"]
compatible_with: ["OFFICIAL-UPSTREAM-ACQUISITION-CORE 0.4.x", "GO-PROVIDER-INTEGRATION-CORE 0.1.x", "GO-PG-OUTBOUND-DELIVERY-FENCE 0.2.x", "GO-CHANNELS-CORE 0.4.x", "GO-RELIABLE-ASYNC-WORKERS 0.3.x"]
incompatible_with: ["SDK Node archivado", "mutación de dominio automática", "secreto en archivo/CLI", "template o destinatario no aprobado", "uso fuera de Meta APIs"]
license_expression: "LicenseRef-Workspace-Owner AND LicenseRef-Meta-Platform-API-Only"
upstream_sources: ["https://github.com/fbsamples/whatsapp-api-examples/tree/de70ee908a67026e642aaee3703d20464e2a9466", "https://developers.facebook.com/docs/whatsapp/cloud-api"]
verified_at: "2026-09-08"
```

## 2. Applicability

Use sólo cuando el proyecto elija Meta WhatsApp Cloud API y acepte licencia/Platform Policy, pruebe Business/app/phone, template/recipient consent, webhook subscription, versión Graph vigente, retención, quota/costo y reconciliación. No usar el SDK Node oficial archivado como baseline nuevo. Tres archivos de código Meta se incluyen `VERBATIM` únicamente como autoridad/referencia; LICENSE conserva texto exacto con un newline final normalizado y queda `ADAPTED`. La ruta ejecutable es `whatsapp_cloud.py`, declarada `ADAPTED` y gobernada por la misma licencia restringida.

## 3. Architecture contract

La plantilla nace bloqueada y nunca infiere una versión móvil de Graph API. Template/language/arity se aprueban exactamente; recipient/phone IDs son numéricos; secrets sólo se resuelven externamente. El outbound conserva patrón oficial `POST /{version}/{phone_number_id}/messages` con Bearer JSON, exige message ID y escribe respuesta/receipt atómicos con identificadores hasheados.

Ingress valida el raw body completo con HMAC-SHA256 y `compare_digest`, exige perfil y WABA/teléfono exactos, verifica challenge y acepta sólo envelope `whatsapp_business_account`/field `messages`. Evidence v2 conserva scope, destinatario, identidad/estado/tiempo y key determinista pseudonimizados; rechaza estructuras ambiguas y lotes mixtos antes de publicar. El writer offline no persiste raw PII inbound; el receptor HTTP 0.9.0 retiene bytes originales sólo con expediente explícito de retención/acceso, sin confundir base64 con cifrado; la respuesta outbound cruda exige controles de acceso/retención. El provider inbox común debe aportar binding de tenant, idempotencia/durable processing y correlación; este pack no escribe dominio ni transforma aceptación HTTP en entrega. Errores eliminan staging, pero no autorizan reintentar un send incierto.

## 4. Exact file manifest

```text
CREATE internal/whatsappbridge/status_host.go
CREATE internal/whatsappbridge/status_host_test.go
CREATE internal/whatsappbridge/notification_history.go
CREATE internal/whatsappbridge/notification_history_test.go
CREATE internal/whatsappbridge/notification_browser_test.go
CREATE whatsapp_cloud/official-source.lock.json
CREATE whatsapp_cloud/provider-profile.template.json
CREATE whatsapp_cloud/template-message.template.json
CREATE whatsapp_cloud/whatsapp_cloud.py
CREATE whatsapp_cloud/test_whatsapp_cloud.py
CREATE whatsapp_cloud/status_reconciliation.py
CREATE whatsapp_cloud/test_status_reconciliation.py
CREATE whatsapp_cloud/PROVENANCE.md
CREATE whatsapp_cloud/README.md
CREATE whatsapp_cloud/upstream/LICENSE
CREATE whatsapp_cloud/upstream/signature_validation_app.py
CREATE whatsapp_cloud/upstream/message_helper.py
CREATE whatsapp_cloud/upstream/incomingWebhook.js
CREATE internal/whatsappbridge/sender.go
CREATE internal/whatsappbridge/sender_test.go
CREATE internal/whatsappbridge/appointment_approval.go
CREATE internal/whatsappbridge/appointment_approval_test.go
CREATE db/migrations/0050_whatsapp_appointment_approval.up.sql
CREATE db/migrations/0050_whatsapp_appointment_approval.down.sql
CREATE db/tests/0050_whatsapp_appointment_approval.test.sql
CREATE internal/whatsappbridge/appointment_notification.go
CREATE internal/whatsappbridge/appointment_notification_test.go
CREATE internal/whatsappbridge/appointment_notification_status.go
CREATE internal/whatsappbridge/appointment_notification_status_test.go
CREATE internal/whatsappbridge/status_observer.go
CREATE internal/whatsappbridge/status_observer_test.go
CREATE internal/whatsappbridge/webhook_receiver.go
CREATE internal/whatsappbridge/webhook_receiver_test.go
CREATE internal/whatsappbridge/status_router.go
CREATE internal/whatsappbridge/status_router_test.go
CREATE internal/whatsappbridge/status_worker.go
CREATE internal/whatsappbridge/status_worker_test.go
CREATE db/migrations/0052_whatsapp_status_routing.up.sql
CREATE db/migrations/0052_whatsapp_status_routing.down.sql
CREATE db/migrations/0051_whatsapp_status_observations.up.sql
CREATE db/migrations/0051_whatsapp_status_observations.down.sql
CREATE db/migrations/0059_whatsapp_conversation_reply.down.sql
CREATE db/migrations/0059_whatsapp_conversation_reply.up.sql
CREATE docs/whatsapp-conversation-operations.md
CREATE internal/whatsappbridge/conversation_browser_issuer_test.go
CREATE internal/whatsappbridge/conversation_browser_test.go
CREATE internal/whatsappbridge/conversation_connected_test.go
CREATE internal/whatsappbridge/conversation_projection_fuzz_test.go
CREATE internal/whatsappbridge/conversation_router.go
CREATE internal/whatsappbridge/reply_approval.go
CREATE internal/whatsappbridge/reply_module.go
CREATE internal/whatsappbridge/reply_recovery.go
CREATE internal/whatsappbridge/reply_request.go
CREATE microsoft_playwright_browser_gate/tests/whatsapp-connected.spec.mjs
CREATE whatsapp_cloud/test_conversation_contract.py
CREATE whatsapp_cloud/test_profile_validation_bridge.py
```

## 5. Materialization blocks

### FILE: `internal/whatsappbridge/status_host.go`

```yaml
block_id: "PYTHON-META-WHATSAPP-CLOUD-ADAPTER:v343:1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "1537eb802b42843c73bca4276cf41d4f8fd67dca5b9a75acfcb313786e39d85a"
variables: []
secrets_allowed: false
```

````go
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
````

### FILE: `internal/whatsappbridge/status_host_test.go`

```yaml
block_id: "PYTHON-META-WHATSAPP-CLOUD-ADAPTER:v343:2"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "aa1b0082cd2789d63180eff1374078a9f84ed117760e58f780910cd919e1b04d"
variables: []
secrets_allowed: false
```

````go
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
````

### FILE: `internal/whatsappbridge/notification_history.go`

```yaml
block_id: "PYTHON-META-WHATSAPP-CLOUD-ADAPTER:v278:3"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "620e8d645eb9f711658afc1be3d61b4fce5ea626040dffb659c92d43de66b4ca"
variables: []
secrets_allowed: false
```

````go
package whatsappbridge

// AUTHORED read-only collection over existing approvals/status owners. No send
// authorization, inferred latest confirmation, new table or customer identity.
import (
	"context"
	"elite.local/enterprise/internal/platform/identity"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5"
	"net/http"
	"net/url"
	"time"
)

var ErrNotificationHistoryLimit = errors.New("whatsappbridge: notification history exceeds bounded view")

type NotificationHistoryItem struct {
	ConfirmationEventID string             `json:"confirmation_event_id"`
	Status              NotificationStatus `json:"status"`
}
type NotificationHistory struct {
	OrganizationID string                    `json:"organization_id"`
	AppointmentID  string                    `json:"appointment_id"`
	Items          []NotificationHistoryItem `json:"items"`
}

// A consistent read-only snapshot, capped at 20 approvals per appointment. An
// overflow is explicit, never silently presented as complete history.
func (s *PostgresAppointmentApprovals) ReadNotificationHistory(ctx context.Context, p identity.Principal, organization, appointment string) (NotificationHistory, error) {
	value := NotificationHistory{OrganizationID: organization, AppointmentID: appointment, Items: []NotificationHistoryItem{}}
	if s == nil || s.pool == nil || p.Subject == "" || p.TenantID == "" || !p.Allowed("appointment:manage") || !p.AllowedOrganization(organization) || organization == "" || len(organization) > 128 || appointment == "" || len(appointment) > 128 {
		return NotificationHistory{}, ErrNotificationNotFound
	}
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead, AccessMode: pgx.ReadOnly})
	if err != nil {
		return NotificationHistory{}, ErrNotificationRead
	}
	defer tx.Rollback(ctx)
	var exists bool
	err = tx.QueryRow(ctx, `select exists(select 1 from crm.appointment a join crm.lead l on l.tenant_id=a.tenant_id and l.lead_id=a.lead_id and l.organization_id=a.organization_id where a.tenant_id=$1 and a.organization_id=$2 and a.appointment_id=$3)`, p.TenantID, organization, appointment).Scan(&exists)
	if err != nil {
		return NotificationHistory{}, ErrNotificationRead
	}
	if !exists {
		return NotificationHistory{}, ErrNotificationNotFound
	}
	rows, err := tx.Query(ctx, `select confirmation_event_id from communication.whatsapp_appointment_approval where tenant_id=$1 and organization_id=$2 and appointment_id=$3 and channel_code='whatsapp' order by confirmation_event_id limit 21`, p.TenantID, organization, appointment)
	if err != nil {
		return NotificationHistory{}, ErrNotificationRead
	}
	ids := []string{}
	for rows.Next() {
		var id string
		if err = rows.Scan(&id); err != nil {
			break
		}
		ids = append(ids, id)
	}
	rowErr := rows.Err()
	rows.Close()
	if err != nil || rowErr != nil {
		return NotificationHistory{}, ErrNotificationRead
	}
	if len(ids) > 20 {
		return NotificationHistory{}, ErrNotificationHistoryLimit
	}
	for _, id := range ids {
		status, err := readNotificationStatus(ctx, tx, p, organization, appointment, id)
		if err != nil {
			return NotificationHistory{}, err
		}
		value.Items = append(value.Items, NotificationHistoryItem{ConfirmationEventID: id, Status: status})
	}
	if tx.Commit(ctx) != nil {
		return NotificationHistory{}, ErrNotificationRead
	}
	return value, nil
}

func (m *AppointmentNotificationModule) registerNotificationHistory(mux *http.ServeMux, verifier identity.Verifier) {
	mux.HandleFunc("GET /v1/franchise/appointments/{id}/whatsapp-confirmations", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		p, ok := m.notificationPrincipal(ctx, w, r, verifier)
		if !ok {
			return
		}
		q, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil || len(r.URL.RawQuery) > 1024 || len(q) != 1 || len(q["organization_id"]) != 1 || q.Get("organization_id") == "" || len(q.Get("organization_id")) > 128 || len(r.PathValue("id")) > 128 || r.ContentLength != 0 {
			notificationProblem(w, 400, "INVALID_QUERY")
			return
		}
		if !p.AllowedOrganization(q.Get("organization_id")) {
			notificationProblem(w, 403, "ORGANIZATION_FORBIDDEN")
			return
		}
		value, err := m.approvals.ReadNotificationHistory(ctx, p, q.Get("organization_id"), r.PathValue("id"))
		switch {
		case errors.Is(err, ErrNotificationNotFound):
			notificationProblem(w, 404, "NOTIFICATION_NOT_FOUND")
			return
		case errors.Is(err, ErrNotificationIntegrity):
			notificationProblem(w, 409, "NOTIFICATION_EVIDENCE_MISMATCH")
			return
		case errors.Is(err, ErrNotificationHistoryLimit):
			notificationProblem(w, 409, "NOTIFICATION_HISTORY_LIMIT")
			return
		case err != nil:
			notificationProblem(w, 503, "NOTIFICATION_READ_UNAVAILABLE")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-store")
		_ = json.NewEncoder(w).Encode(value)
	})
}
````

### FILE: `internal/whatsappbridge/notification_history_test.go`

```yaml
block_id: "PYTHON-META-WHATSAPP-CLOUD-ADAPTER:v278:4"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "fdba06ff3e9224bd0504496aaed6c9fc3413a706fb182dc577bd90a8d6674cf8"
variables: []
secrets_allowed: false
```

````go
package whatsappbridge

import (
	"context"
	"elite.local/enterprise/internal/platform/identity"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNotificationHistoryReadOnlyAndScoped(t *testing.T) {
	r, o, pool, p, c, _, _ := routerFixture(t, statusBatch("delivered", "1603086314", nil))
	ctx := context.Background()
	value, err := o.Approvals.ReadNotificationHistory(ctx, p, c.OrganizationID, c.AppointmentID)
	if err != nil || len(value.Items) != 1 || value.Items[0].ConfirmationEventID != c.ConfirmationEventID || value.Items[0].Status.DeliveryStatus != "not_observed_by_this_reader" {
		t.Fatal(err, value)
	}
	w, err := NewStatusWorker(r, "history-worker", time.Minute, time.Second)
	if err != nil {
		t.Fatal(err)
	}
	if result, err := w.ProcessOnce(ctx, p); err != nil || !result.Completed {
		t.Fatal(err, result)
	}
	value, err = o.Approvals.ReadNotificationHistory(ctx, p, c.OrganizationID, c.AppointmentID)
	if err != nil || len(value.Items) != 1 || value.Items[0].Status.DeliveryStatus != "observed_delivered" {
		t.Fatal(err, value)
	}
	for _, mode := range []string{"permission", "organization", "tenant", "appointment"} {
		t.Run(mode, func(t *testing.T) {
			bad := p
			appointment := c.AppointmentID
			switch mode {
			case "permission":
				bad.Permissions = map[string]struct{}{}
			case "organization":
				bad.Organizations = map[string]struct{}{}
			case "tenant":
				bad.TenantID = "00000000-0000-4000-8000-000000000000"
			case "appointment":
				appointment = "other"
			}
			if _, err := o.Approvals.ReadNotificationHistory(ctx, bad, c.OrganizationID, appointment); !errors.Is(err, ErrNotificationNotFound) {
				t.Fatal("foreign history", err)
			}
		})
	}
	other, command, empty := approvalFixtureData(t, pool)
	value, err = empty.ReadNotificationHistory(ctx, other, command.OrganizationID, command.AppointmentID)
	if err != nil || len(value.Items) != 0 {
		t.Fatal("empty history", err, value)
	}
	if _, err := pool.Exec(ctx, `update communication.outbound_delivery set recipient_hmac=repeat('f',64) where tenant_id=$1`, p.TenantID); err != nil {
		t.Fatal(err)
	}
	if _, err := o.Approvals.ReadNotificationHistory(ctx, p, c.OrganizationID, c.AppointmentID); !errors.Is(err, ErrNotificationIntegrity) {
		t.Fatal("tampered history", err)
	}
}

func TestNotificationHistoryHTTPBoundary(t *testing.T) {
	_, o, _, p, c, _, _ := routerFixture(t, statusBatch("delivered", "1603086314", nil))
	v, sign := notificationIssuer(t)
	module := &AppointmentNotificationModule{approvals: o.Approvals, sender: &Sender{TenantID: p.TenantID}}
	mux := http.NewServeMux()
	module.Register(mux, v)
	api := httptest.NewServer(mux)
	defer api.Close()
	for _, test := range []struct {
		name, query, token string
		status             int
	}{{"valid", "organization_id=store-1", sign(p), 200}, {"duplicate", "organization_id=store-1&organization_id=store-1", sign(p), 400}, {"extra", "organization_id=store-1&tenant=other", sign(p), 400}, {"foreign", "organization_id=other", sign(p), 403}, {"no-token", "organization_id=store-1", "", 401}, {"no-permission", "organization_id=store-1", sign(identity.Principal{TenantID: p.TenantID, Subject: p.Subject}), 403}} {
		t.Run(test.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", api.URL+"/v1/franchise/appointments/"+c.AppointmentID+"/whatsapp-confirmations?"+test.query, nil)
			if test.token != "" {
				req.Header.Set("Authorization", "Bearer "+test.token)
			}
			response, err := api.Client().Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer response.Body.Close()
			if response.StatusCode != test.status || response.Header.Get("Cache-Control") != "no-store" {
				t.Fatal(response.StatusCode)
			}
			var body map[string]any
			if json.NewDecoder(response.Body).Decode(&body) != nil {
				t.Fatal("response JSON")
			}
			raw, _ := json.Marshal(body)
			for _, secret := range []string{"5491112345678", "wamid.synthetic", "test-app-secret"} {
				if strings.Contains(string(raw), secret) {
					t.Fatal("sensitive history")
				}
			}
		})
	}
}

func TestNotificationHistoryOverflowIsExplicit(t *testing.T) {
	_, o, pool, p, c, _, _ := routerFixture(t, statusBatch("delivered", "1603086314", nil))
	ctx := context.Background()
	// Synthetic historical approvals test the read limit, not the send admission.
	for i := 0; i < 20; i++ {
		var id string
		if err := pool.QueryRow(ctx, `insert into crm.appointment_transition(tenant_id,transition_id,appointment_id,from_state,to_state,actor_subject) values($1,gen_random_uuid(),$2,'requested','confirmed','fixture-history') returning transition_id`, p.TenantID, c.AppointmentID).Scan(&id); err != nil {
			t.Fatal(err)
		}
		if _, err := pool.Exec(ctx, `insert into communication.whatsapp_appointment_approval(tenant_id,delivery_key,channel_code,appointment_id,appointment_version,organization_id,confirmation_event_id,external_id_hmac,binding_version,lead_id,subject_id,policy_version,consent_id,consent_purpose,consent_evidence_sha256,message_sha256,profile_sha256,approval_request_sha256,evidence_sha256,approved_by,approved_at,expires_at) select tenant_id,'fixture-history-'||$2::text,channel_code,appointment_id,appointment_version,organization_id,$2::uuid,external_id_hmac,binding_version,lead_id,subject_id,policy_version,consent_id,consent_purpose,consent_evidence_sha256,message_sha256,profile_sha256,approval_request_sha256,evidence_sha256,approved_by,approved_at,expires_at from communication.whatsapp_appointment_approval where tenant_id=$1 and delivery_key=$3`, p.TenantID, id, c.Message.DeliveryKey); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := o.Approvals.ReadNotificationHistory(ctx, p, c.OrganizationID, c.AppointmentID); !errors.Is(err, ErrNotificationHistoryLimit) {
		t.Fatal("silently truncated history", err)
	}
}
````

### FILE: `internal/whatsappbridge/notification_browser_test.go`

```yaml
block_id: "PYTHON-META-WHATSAPP-CLOUD-ADAPTER:v278:5"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "165507a3c0703de4d08beab732100d1f466fa9daa301ee23650b86548ba1dba0"
variables: []
secrets_allowed: false
```

````go
package whatsappbridge

// AUTHORED opt-in browser fixture. Real Go/PG/Next/TLS; synthetic issuer,
// approval, receipt and provider webhook. Never a live Meta or IdP assertion.
import (
	"context"
	"elite.local/enterprise/internal/franchisejourney"
	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/platform/postgres"
	"elite.local/enterprise/internal/platform/randomid"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

type notificationClock struct{}

func (notificationClock) Now() time.Time { return time.Now() }

func TestNotificationStatusBrowserPostgres(t *testing.T) {
	if os.Getenv("ELITE_NOTIFICATION_BROWSER_E2E") != "1" {
		t.Skip("explicit isolated notification browser gate not selected")
	}
	web := os.Getenv("ELITE_WEB_ROOT")
	if !filepath.IsAbs(web) {
		t.Fatal("absolute materialized web root required")
	}
	for _, relative := range []string{"node_modules/next/dist/bin/next", ".next/BUILD_ID", "microsoft_playwright_browser_gate/node_modules/@playwright/test/cli.js"} {
		if _, err := os.Stat(filepath.Join(web, relative)); err != nil {
			t.Fatal("browser prerequisite missing; install gate with --ignore-workspace and build exact web first: " + relative)
		}
	}
	for _, project := range []string{"chromium-desktop", "chromium-mobile", "firefox-desktop", "webkit-desktop"} {
		t.Run(project, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
			defer cancel()
			r, o, pool, p, c, _, _ := routerFixture(t, statusBatch("delivered", "1603086314", nil))
			w, err := NewStatusWorker(r, "browser-worker", time.Minute, time.Second)
			if err != nil {
				t.Fatal(err)
			}
			if result, err := w.ProcessOnce(ctx, p); err != nil || !result.Completed {
				t.Fatal("worker prerequisite", err, result)
			}
			verifier, sign := notificationIssuer(t)
			module := &AppointmentNotificationModule{approvals: o.Approvals, sender: &Sender{TenantID: p.TenantID}}
			mux := http.NewServeMux()
			module.Register(mux, verifier)
			httpapi.FranchiseJourneyModule{Service: franchisejourney.NewService(postgres.NewFranchiseJourney(pool), randomid.Generator{}, notificationClock{})}.Register(mux, verifier)
			api := httptest.NewServer(mux)
			defer api.Close()
			target, _ := url.Parse("http://127.0.0.1:4173")
			edge := httptest.NewTLSServer(httputil.NewSingleHostReverseProxy(target))
			defer edge.Close()
			var start time.Time
			if err := pool.QueryRow(ctx, `select starts_at from crm.appointment where tenant_id=$1 and appointment_id=$2`, p.TenantID, c.AppointmentID).Scan(&start); err != nil {
				t.Fatal(err)
			}
			env := []string{}
			for _, entry := range os.Environ() {
				name := strings.ToUpper(strings.SplitN(entry, "=", 2)[0])
				if name == "BUSINESS_CONFIG_FILE" || name == "TEST_DATABASE_URL" || name == "DATABASE_URL" || name == "ELITE_WHATSAPP_TEST_DATABASE_URL" || strings.HasPrefix(name, "ELITE_NOTIFICATION_") || strings.HasPrefix(name, "ELITE_CONFIRMATION_") || name == "ENTERPRISE_API_BASE_URL" || name == "AUTH_SESSION_SECRET" || name == "APP_BASE_URL" || name == "ELITE_BASE_URL" || name == "ELITE_RUNTIME_ONLY" || name == "ELITE_OPERATOR_AGENDA_E2E" || name == "ELITE_WEB_ROOT" {
					continue
				}
				env = append(env, entry)
			}
			token := sign(p)
			secret := "synthetic-" + (randomid.Generator{}).New() + (randomid.Generator{}).New()
			// Use the existing business configuration owner, never a new runtime flag.
			rawConfig, err := os.ReadFile(filepath.Join(web, "config", "business.example.json"))
			if err != nil {
				t.Fatal(err)
			}
			var config map[string]any
			if err := json.Unmarshal(rawConfig, &config); err != nil {
				t.Fatal(err)
			}
			features, ok := config["features"].(map[string]any)
			if !ok {
				t.Fatal("business fixture features missing")
			}
			features["whatsapp_status_history"] = true
			configBytes, err := json.Marshal(config)
			if err != nil {
				t.Fatal(err)
			}
			configName := "notification-fixture-" + (randomid.Generator{}).New() + ".json"
			configPath := filepath.Join(web, "config", configName)
			configFile, err := os.OpenFile(configPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() {
				if err := os.Remove(configPath); err != nil {
					t.Error(err)
				}
			})
			_, writeErr := configFile.Write(configBytes)
			closeErr := configFile.Close()
			if writeErr != nil || closeErr != nil {
				t.Fatal("business fixture could not be persisted")
			}
			env = append(env, "BUSINESS_CONFIG_FILE="+configName)
			env = append(env, "ENTERPRISE_API_BASE_URL="+api.URL, "APP_BASE_URL="+edge.URL, "ELITE_BASE_URL="+edge.URL, "ELITE_WEB_ROOT="+web, "ELITE_NOTIFICATION_BROWSER_E2E=1", "AUTH_SESSION_SECRET="+secret, "ELITE_NOTIFICATION_TOKEN="+token, "ELITE_NOTIFICATION_TENANT="+p.TenantID, "ELITE_NOTIFICATION_ORGANIZATION="+c.OrganizationID, "ELITE_NOTIFICATION_APPOINTMENT="+c.AppointmentID, "ELITE_NOTIFICATION_EVENT="+c.ConfirmationEventID, "ELITE_NOTIFICATION_DAY="+start.UTC().Format("2006-01-02"))
			server := exec.CommandContext(ctx, "node", filepath.Join(web, "node_modules", "next", "dist", "bin", "next"), "start", "--hostname", "127.0.0.1", "--port", "4173")
			server.Dir = web
			server.Env = env
			if err := server.Start(); err != nil {
				t.Fatal(err)
			}
			defer func() { _ = server.Process.Kill(); _ = server.Wait() }()
			ready := false
			client := &http.Client{Timeout: time.Second}
			for deadline := time.Now().Add(20 * time.Second); time.Now().Before(deadline); {
				response, err := client.Get(target.String() + "/icon.svg")
				if err == nil {
					response.Body.Close()
					if response.StatusCode == 200 {
						ready = true
						break
					}
				}
				select {
				case <-ctx.Done():
					t.Fatal("web startup cancelled")
				case <-time.After(100 * time.Millisecond):
				}
			}
			if !ready {
				t.Fatal("web fixture did not start")
			}
			cmd := exec.CommandContext(ctx, "node", filepath.Join(web, "microsoft_playwright_browser_gate", "node_modules", "@playwright", "test", "cli.js"), "test", "tests/notification-status.spec.mjs", "--project", project, "--workers=1", "--retries=0", "--max-failures=1", "--output", filepath.Join(web, "microsoft_playwright_browser_gate", "test-results", "notifications-"+project))
			cmd.Dir = filepath.Join(web, "microsoft_playwright_browser_gate")
			cmd.Env = env
			output, err := cmd.CombinedOutput()
			clean := strings.ReplaceAll(strings.ReplaceAll(string(output), token, "[REDACTED_TEST_TOKEN]"), secret, "[REDACTED_TEST_SECRET]")
			t.Log(clean)
			if err != nil {
				t.Fatal("notification browser gate", err)
			}
			var jobs, attempts, observations int
			if err := pool.QueryRow(ctx, `select (select count(*) from platform.job where tenant_id=$1 and completed_at is not null),(select sum(attempt_count)::int from communication.outbound_delivery where tenant_id=$1),(select count(*) from communication.whatsapp_status_observation where tenant_id=$1)`, p.TenantID).Scan(&jobs, &attempts, &observations); err != nil || jobs != 1 || attempts != 1 || observations != 1 {
				t.Fatal("read changed durable effects", err, jobs, attempts, observations)
			}
			t.Log("NOTIFICATION_BROWSER_POSTGRES_PASS status=observed_delivered jobs=1 send_attempts=1 observations=1 no_resend=true")
		})
	}
}
````


### FILE: `internal/whatsappbridge/status_worker.go`

```yaml
block_id: "META-WHATSAPP:status-worker:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "a94051c187923cb19a10807369c28f78dea7c3e676d2da5bd2659ec85e7b5ffa"
variables: []
secrets_allowed: false
```

````go
package whatsappbridge

// AUTHORED opt-in worker. Reuses the existing job queue, router, observer and
// append-only audit. No daemon is registered, no send or customer reply occurs.
import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/postgres"
	"github.com/jackc/pgx/v5"
)

var ErrStatusWork = errors.New("whatsappbridge: status job unresolved; inspect scoped job and audit")

type StatusWorkResult struct {
	Claimed, Completed, FailureRecorded, Terminal bool
	InsertedObservations                          int
}

type StatusWorker struct {
	router       *StatusRouter
	jobs         *postgres.Jobs
	worker       string
	lease, retry time.Duration
	busy         chan struct{}
}

func NewStatusWorker(router *StatusRouter, worker string, lease, retry time.Duration) (*StatusWorker, error) {
	if router == nil || router.observer.Approvals == nil || router.observer.Approvals.pool == nil || router.busy == nil || !webhookConnection.MatchString(worker) || lease < time.Second || lease > 2*time.Minute || retry < 0 || retry > time.Hour {
		return nil, ErrStatusWork
	}
	return &StatusWorker{router: router, jobs: postgres.NewJobs(router.observer.Approvals.pool), worker: worker, lease: lease, retry: retry, busy: make(chan struct{}, 1)}, nil
}

// ProcessOnce claims at most one job, in the configured tenant/connection/type.
// p must come from actual authentication; no identity is derived from payload.
// A caller must handle errors even when observations were partially committed.
func (w *StatusWorker) ProcessOnce(ctx context.Context, p identity.Principal) (StatusWorkResult, error) {
	var result StatusWorkResult
	if w == nil || p.TenantID != w.router.observer.TenantID || p.Subject == "" || len(p.Subject) > 255 || !p.Allowed("appointment:manage") || ctx.Err() != nil {
		return result, ErrStatusWork
	}
	select {
	case w.busy <- struct{}{}:
		defer func() { <-w.busy }()
	default:
		return result, ErrStatusWork
	}
	ctx, cancelCall := context.WithTimeout(ctx, w.lease+5*time.Second)
	defer cancelCall()
	var connectionOrg *string
	err := w.router.observer.Approvals.pool.QueryRow(ctx, `select organization_id from integration.provider_connection where tenant_id=$1 and connection_id=$2 and provider_code='meta-whatsapp' and state='active'`, p.TenantID, w.router.connection).Scan(&connectionOrg)
	if err != nil || (connectionOrg != nil && !p.AllowedOrganization(*connectionOrg)) {
		return result, ErrStatusWork
	}
	match, _ := json.Marshal(map[string]string{"connection_id": w.router.connection, "provider_code": "meta-whatsapp", "event_type": "whatsapp.raw_webhook.received.v1"})
	if expired, err := w.quarantineExpired(ctx, p, match); err != nil {
		return result, ErrStatusWork
	} else if expired {
		result.Terminal = true
		result.FailureRecorded = true
		return result, ErrStatusWork
	}
	jobs, err := w.jobs.ClaimScoped(ctx, "provider-events", w.worker, w.lease, 1, postgres.JobScope{TenantID: p.TenantID, JobType: "provider.webhook.received", SchemaVersion: 1, PayloadMatch: match})
	if err != nil {
		return result, ErrStatusWork
	}
	if len(jobs) == 0 {
		return result, nil
	}
	result.Claimed = true
	job := jobs[0]
	event, err := w.eventID(job)
	if err != nil {
		// No trustworthy inbox identity: retain job/audit, do not guess an event.
		return w.recordFailure(ctx, p, job, "", "JOB_CONTRACT_INVALID", result)
	}
	started := time.Now()
	remaining, err := w.jobs.RemainingLease(ctx, job.TenantID, job.JobID, w.worker, job.Attempts)
	remaining -= time.Since(started)
	if err != nil || remaining <= 0 {
		return result, ErrStatusWork
	}
	workCtx, cancel := context.WithTimeout(ctx, remaining)
	result.InsertedObservations, err = w.router.observeRetained(workCtx, p, event, func(ctx context.Context, tx pgx.Tx) error {
		if err := postgres.CompleteJobInTx(ctx, tx, job, w.worker); err != nil {
			return err
		}
		tag, err := tx.Exec(ctx, `update integration.webhook_event set state='processed',processed_at=coalesce(processed_at,clock_timestamp()),last_error_code=null
 where tenant_id=$1 and connection_id=$2 and provider_event_id=$3 and provider_code='meta-whatsapp'
 and event_type='whatsapp.raw_webhook.received.v1' and state in ('received','processed')`, p.TenantID, w.router.connection, event)
		if err != nil {
			return err
		}
		if tag.RowsAffected() != 1 {
			return ErrStatusWork
		}
		return statusWorkAudit(ctx, tx, p, job, event, "completed", "STATUS_OBSERVED", false)
	})
	cancel()
	if err != nil {
		return w.recordFailure(ctx, p, job, event, "STATUS_UNRESOLVED", result)
	}
	result.Completed = true
	return result, nil
}

func (w *StatusWorker) quarantineExpired(ctx context.Context, p identity.Principal, match json.RawMessage) (bool, error) {
	tx, err := w.router.observer.Approvals.pool.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer tx.Rollback(ctx)
	var job postgres.Job
	err = tx.QueryRow(ctx, `select tenant_id,job_id,queue,job_type,schema_version,payload,attempts,max_attempts from platform.job
 where tenant_id=$1 and queue='provider-events' and job_type='provider.webhook.received' and schema_version=1 and payload @> $2::jsonb
 and completed_at is null and terminal_error_code is null and attempts>=max_attempts and claimed_until<clock_timestamp()
 order by claimed_until,job_id for update skip locked limit 1`, p.TenantID, match).Scan(&job.TenantID, &job.JobID, &job.Queue, &job.JobType, &job.SchemaVersion, &job.Payload, &job.Attempts, &job.MaxAttempts)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if err = postgres.ExhaustJobInTx(ctx, tx, job); err != nil {
		return false, err
	}
	event, _ := w.eventID(job)
	if event != "" {
		if _, err = tx.Exec(ctx, `update integration.webhook_event set state='failed',last_error_code='LEASE_EXHAUSTED' where tenant_id=$1 and connection_id=$2 and provider_event_id=$3 and provider_code='meta-whatsapp' and event_type='whatsapp.raw_webhook.received.v1' and state='received'`, p.TenantID, w.router.connection, event); err != nil {
			return false, err
		}
	}
	if err = statusWorkAudit(ctx, tx, p, job, event, "failed", "LEASE_EXHAUSTED", true); err != nil {
		return false, err
	}
	if err = tx.Commit(ctx); err != nil {
		return false, err
	}
	return true, nil
}

func (w *StatusWorker) eventID(job postgres.Job) (string, error) {
	var fields map[string]string
	if job.TenantID != w.router.observer.TenantID || job.Queue != "provider-events" || job.JobType != "provider.webhook.received" || job.SchemaVersion != 1 || len(job.Payload) > 4096 || json.Unmarshal(job.Payload, &fields) != nil || len(fields) != 4 || fields["connection_id"] != w.router.connection || fields["provider_code"] != "meta-whatsapp" || fields["event_type"] != "whatsapp.raw_webhook.received.v1" {
		return "", ErrStatusWork
	}
	event := fields["provider_event_id"]
	if !strings.HasPrefix(event, "wa:") || !validDigest(strings.TrimPrefix(event, "wa:")) {
		return "", ErrStatusWork
	}
	return event, nil
}

func (w *StatusWorker) recordFailure(ctx context.Context, p identity.Principal, job postgres.Job, event, code string, result StatusWorkResult) (StatusWorkResult, error) {
	// Bounded metadata cleanup after cancellation is not another processing try.
	cleanup, cancel := context.WithTimeout(context.WithoutCancel(ctx), 3*time.Second)
	defer cancel()
	tx, err := w.router.observer.Approvals.pool.Begin(cleanup)
	if err != nil {
		return result, ErrStatusWork
	}
	defer tx.Rollback(cleanup)
	terminal, err := postgres.FailJobInTx(cleanup, tx, job, w.worker, code, w.retry)
	if err != nil {
		return result, ErrStatusWork
	}
	if event != "" {
		// Missing/failed/processed inbox cannot be rewritten as a retryable one.
		// Job failure and audit remain useful even if its inbox is absent.
		_, err = tx.Exec(cleanup, `update integration.webhook_event set state=case when $4 then 'failed' else 'received' end,last_error_code=$5
 where tenant_id=$1 and connection_id=$2 and provider_event_id=$3 and provider_code='meta-whatsapp'
 and event_type='whatsapp.raw_webhook.received.v1' and state='received'`, p.TenantID, w.router.connection, event, terminal, code)
		if err != nil {
			return result, ErrStatusWork
		}
	}
	if statusWorkAudit(cleanup, tx, p, job, event, "failed", code, terminal) != nil || tx.Commit(cleanup) != nil {
		return result, ErrStatusWork
	}
	result.FailureRecorded = true
	result.Terminal = terminal
	return result, ErrStatusWork
}

func statusWorkAudit(ctx context.Context, tx pgx.Tx, p identity.Principal, job postgres.Job, event, action, code string, terminal bool) error {
	_, err := tx.Exec(ctx, `insert into audit.event(tenant_id,event_id,actor_subject,action,resource_type,resource_id,decision,reason_code,evidence)
 values($1,gen_random_uuid(),$2,$3,'provider_job',$4,'system',$5,jsonb_build_object('attempt',$6::int,'terminal',$7::bool,'provider_event_sha256',$8::text))`, p.TenantID, p.Subject, "whatsapp.status_job."+action, job.JobID, code, job.Attempts, terminal, digest([]byte(event)))
	return err
}
````

### FILE: `internal/whatsappbridge/status_worker_test.go`

```yaml
block_id: "META-WHATSAPP:status-worker-tests:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "917692f3144e005b1e3c5968cbc1f9ef5c0657963ef3d2b5edb58cb13d3c2128"
variables: []
secrets_allowed: false
```

````go
package whatsappbridge

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"elite.local/enterprise/internal/platform/postgres"
)

func TestStatusWorkerRejectsUninitializedRouter(t *testing.T) {
	for _, router := range []*StatusRouter{nil, {}} {
		if _, err := NewStatusWorker(router, "worker-main", time.Minute, 0); err == nil {
			t.Fatal("uninitialized router accepted")
		}
	}
}

func TestStatusWorkerScopeCompletionAndAudit(t *testing.T) {
	r, _, pool, p, _, event, _ := routerFixture(t, statusBatch("delivered", "1603086314", nil))
	ctx := context.Background()
	var original postgres.Job
	if err := pool.QueryRow(ctx, `select tenant_id,job_id,queue,job_type,schema_version,payload,attempts,max_attempts from platform.job where tenant_id=$1`, p.TenantID).Scan(&original.TenantID, &original.JobID, &original.Queue, &original.JobType, &original.SchemaVersion, &original.Payload, &original.Attempts, &original.MaxAttempts); err != nil {
		t.Fatal(err)
	}
	var otherTenant string
	if err := pool.QueryRow(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values(gen_random_uuid(),'scope-'||gen_random_uuid()::text,'Synthetic scope','Scope') returning tenant_id`).Scan(&otherTenant); err != nil {
		t.Fatal(err)
	}
	var excluded []string
	for _, mode := range []string{"tenant", "queue", "job-type", "schema", "connection", "provider", "event-type"} {
		var payload map[string]string
		if json.Unmarshal(original.Payload, &payload) != nil {
			t.Fatal("fixture payload")
		}
		tenant, queue, kind, version := p.TenantID, original.Queue, original.JobType, 1
		switch mode {
		case "tenant":
			tenant = otherTenant
		case "queue":
			queue = "other"
		case "job-type":
			kind = "other"
		case "schema":
			version = 2
		case "connection":
			payload["connection_id"] = "other"
		case "provider":
			payload["provider_code"] = "other"
		case "event-type":
			payload["event_type"] = "other"
		}
		raw, _ := json.Marshal(payload)
		var id string
		if err := pool.QueryRow(ctx, `insert into platform.job(tenant_id,job_id,queue,job_type,schema_version,payload,max_attempts) values($1,gen_random_uuid(),$2,$3,$4,$5,2) returning job_id`, tenant, queue, kind, version, raw).Scan(&id); err != nil {
			t.Fatal(err)
		}
		excluded = append(excluded, id)
	}
	w, err := NewStatusWorker(r, "worker-main", time.Minute, 0)
	if err != nil {
		t.Fatal(err)
	}
	result, err := w.ProcessOnce(ctx, p)
	if err != nil || !result.Claimed || !result.Completed || result.InsertedObservations != 1 || result.FailureRecorded {
		t.Fatalf("completion: %+v %v", result, err)
	}
	// A new instance after the commit must not execute the job again.
	w, err = NewStatusWorker(r, "worker-restarted", time.Minute, 0)
	if err != nil {
		t.Fatal(err)
	}
	result, err = w.ProcessOnce(ctx, p)
	if err != nil || result.Claimed || result.Completed {
		t.Fatalf("replay: %+v %v", result, err)
	}
	var completed, untouched bool
	if err := pool.QueryRow(ctx, `select j.completed_at is not null and j.attempts=1 and j.terminal_error_code is null and e.state='processed' and e.processed_at is not null from platform.job j join integration.webhook_event e on e.tenant_id=j.tenant_id and e.provider_event_id=$2 where j.tenant_id=$1 and j.job_id=$3`, p.TenantID, event, original.JobID).Scan(&completed); err != nil || !completed {
		t.Fatal("job/inbox completion", completed, err)
	}
	if err := pool.QueryRow(ctx, `select count(*)=7 and bool_and(attempts=0 and completed_at is null and terminal_error_code is null) from platform.job where job_id=any($1::uuid[])`, excluded).Scan(&untouched); err != nil || !untouched {
		t.Fatal("claimed another scope", untouched, err)
	}
	var audits int
	var safe bool
	if err := pool.QueryRow(ctx, `select count(*),bool_and(reason_code='STATUS_OBSERVED' and evidence ? 'attempt' and evidence::text not like '%wamid%' and evidence::text not like '%549111%') from audit.event where tenant_id=$1 and action='whatsapp.status_job.completed'`, p.TenantID).Scan(&audits, &safe); err != nil || audits != 1 || !safe {
		t.Fatal("audit", audits, safe, err)
	}
}

func TestStatusWorkerTransientAndTerminalFailures(t *testing.T) {
	for _, terminal := range []bool{false, true} {
		t.Run(fmt.Sprint(terminal), func(t *testing.T) {
			batch := statusBatch("delivered", "1603086314", nil)
			if terminal {
				batch = statusBatch("delivered", "1603086314", map[string]any{"id": "wamid.unresolved"})
			}
			r, _, pool, p, _, event, path := routerFixture(t, batch)
			ctx := context.Background()
			if _, err := pool.Exec(ctx, `update platform.job set max_attempts=2 where tenant_id=$1`, p.TenantID); err != nil {
				t.Fatal(err)
			}
			raw, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if !terminal {
				if err := os.Remove(path); err != nil {
					t.Fatal(err)
				}
			}
			w, err := NewStatusWorker(r, "worker-retry", time.Minute, 0)
			if err != nil {
				t.Fatal(err)
			}
			first, err := w.ProcessOnce(ctx, p)
			if err == nil || !first.Claimed || first.Completed || !first.FailureRecorded || first.Terminal {
				t.Fatalf("retry: %+v %v", first, err)
			}
			if !terminal {
				if err := os.WriteFile(path, raw, 0600); err != nil {
					t.Fatal(err)
				}
			}
			second, err := w.ProcessOnce(ctx, p)
			if terminal {
				if err == nil || !second.Terminal || !second.FailureRecorded || second.Completed {
					t.Fatalf("terminal: %+v %v", second, err)
				}
			} else if err != nil || !second.Completed || second.InsertedObservations != 1 {
				t.Fatalf("recovery: %+v %v", second, err)
			}
			third, err := w.ProcessOnce(ctx, p)
			if err != nil || third.Claimed {
				t.Fatalf("unbounded retry: %+v %v", third, err)
			}
			var state string
			var audits int
			if err := pool.QueryRow(ctx, `select state from integration.webhook_event where tenant_id=$1 and provider_event_id=$2`, p.TenantID, event).Scan(&state); err != nil {
				t.Fatal(err)
			}
			if (terminal && state != "failed") || (!terminal && state != "processed") {
				t.Fatal("inbox state", state)
			}
			if err := pool.QueryRow(ctx, `select count(*) from audit.event where tenant_id=$1 and action like 'whatsapp.status_job.%'`, p.TenantID).Scan(&audits); err != nil || audits != 2 {
				t.Fatal("missing attempt audit", audits, err)
			}
		})
	}
}

func TestStatusWorkerAtomicCompletionRollback(t *testing.T) {
	r, _, pool, p, _, event, _ := routerFixture(t, statusBatch("delivered", "1603086314", nil))
	ctx := context.Background()
	// Fault injection is limited to this synthetic tenant in a disposable DB.
	if !notificationEventID.MatchString(p.TenantID) {
		t.Fatal("unsafe fixture tenant")
	}
	sql := fmt.Sprintf(`create function integration.reject_v277_completion() returns trigger language plpgsql as $body$ begin if new.tenant_id='%s'::uuid and new.state='processed' then raise exception 'synthetic ACK failure'; end if; return new; end $body$; create trigger reject_v277_completion before update on integration.webhook_event for each row execute function integration.reject_v277_completion();`, p.TenantID)
	if _, err := pool.Exec(ctx, sql); err != nil {
		t.Fatal(err)
	}
	cleanup := func() {
		if _, err := pool.Exec(ctx, `drop trigger if exists reject_v277_completion on integration.webhook_event; drop function if exists integration.reject_v277_completion();`); err != nil {
			t.Fatal(err)
		}
	}
	t.Cleanup(cleanup)
	w, err := NewStatusWorker(r, "worker-atomic", time.Minute, 0)
	if err != nil {
		t.Fatal(err)
	}
	first, err := w.ProcessOnce(ctx, p)
	if err == nil || first.Completed || first.InsertedObservations != 1 || !first.FailureRecorded {
		t.Fatalf("failed ACK: %+v %v", first, err)
	}
	var pending bool
	if err := pool.QueryRow(ctx, `select j.completed_at is null and e.state='received' and e.processed_at is null from platform.job j join integration.webhook_event e on e.tenant_id=j.tenant_id and e.provider_event_id=$2 where j.tenant_id=$1`, p.TenantID, event).Scan(&pending); err != nil || !pending {
		t.Fatal("partial job/inbox commit", pending, err)
	}
	cleanup()
	second, err := w.ProcessOnce(ctx, p)
	if err != nil || !second.Completed || second.InsertedObservations != 0 {
		t.Fatalf("ACK recovery: %+v %v", second, err)
	}
	var count int
	if err := pool.QueryRow(ctx, `select count(*) from communication.whatsapp_status_observation where tenant_id=$1`, p.TenantID).Scan(&count); err != nil || count != 1 {
		t.Fatal("duplicate observation", count, err)
	}
}

func TestStatusWorkerReclaimsAndQuarantinesCrashedAttempt(t *testing.T) {
	for _, exhausted := range []bool{false, true} {
		t.Run(fmt.Sprint(exhausted), func(t *testing.T) {
			r, _, pool, p, _, event, _ := routerFixture(t, statusBatch("delivered", "1603086314", nil))
			ctx := context.Background()
			attempt := 1
			if exhausted {
				attempt = 10
			}
			if _, err := pool.Exec(ctx, `update platform.job set attempts=$2,claimed_by='crashed',claimed_until=clock_timestamp()-interval '1 second' where tenant_id=$1`, p.TenantID, attempt); err != nil {
				t.Fatal(err)
			}
			w, err := NewStatusWorker(r, "worker-recovery", time.Minute, 0)
			if err != nil {
				t.Fatal(err)
			}
			result, err := w.ProcessOnce(ctx, p)
			if exhausted {
				if err == nil || result.Completed || !result.Terminal || !result.FailureRecorded {
					t.Fatalf("orphan: %+v %v", result, err)
				}
			} else if err != nil || !result.Completed {
				t.Fatalf("reclaimed: %+v %v", result, err)
			}
			var state string
			var count int
			if err := pool.QueryRow(ctx, `select state from integration.webhook_event where tenant_id=$1 and provider_event_id=$2`, p.TenantID, event).Scan(&state); err != nil {
				t.Fatal(err)
			}
			if exhausted && state != "failed" {
				t.Fatal(state)
			}
			if err := pool.QueryRow(ctx, `select count(*) from communication.whatsapp_status_observation where tenant_id=$1`, p.TenantID).Scan(&count); err != nil {
				t.Fatal(err)
			}
			if exhausted && count != 0 {
				t.Fatal("quarantine executed effect")
			}
			again, err := w.ProcessOnce(ctx, p)
			if err != nil || again.Claimed || again.Terminal {
				t.Fatalf("repeat: %+v %v", again, err)
			}
		})
	}
}

func TestStatusWorkerConcurrentClaims(t *testing.T) {
	r, _, pool, p, _, _, _ := routerFixture(t, statusBatch("delivered", "1603086314", nil))
	other, err := NewStatusRouter(&r.observer, r.connection, r.retention, r.key, r.maxRoutes)
	if err != nil {
		t.Fatal(err)
	}
	w1, err := NewStatusWorker(r, "worker-a", time.Minute, 0)
	if err != nil {
		t.Fatal(err)
	}
	w2, err := NewStatusWorker(other, "worker-b", time.Minute, 0)
	if err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	start := make(chan struct{})
	counts := make(chan bool, 2)
	errs := make(chan error, 2)
	for _, w := range []*StatusWorker{w1, w2} {
		wg.Add(1)
		go func(w *StatusWorker) {
			defer wg.Done()
			<-start
			res, err := w.ProcessOnce(context.Background(), p)
			counts <- res.Completed
			errs <- err
		}(w)
	}
	close(start)
	wg.Wait()
	close(counts)
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatal(err)
		}
	}
	completed := 0
	for done := range counts {
		if done {
			completed++
		}
	}
	if completed != 1 {
		t.Fatal("completion count", completed)
	}
	var attempt int
	if err := pool.QueryRow(context.Background(), `select attempts from platform.job where tenant_id=$1`, p.TenantID).Scan(&attempt); err != nil || attempt != 1 {
		t.Fatal("claim race", attempt, err)
	}
}

func TestStatusWorkerClaimLostAfterObservation(t *testing.T) {
	r, _, pool, p, _, event, _ := routerFixture(t, statusBatch("delivered", "1603086314", nil))
	ctx := context.Background()
	w, err := NewStatusWorker(r, "worker-old", time.Minute, 0)
	if err != nil {
		t.Fatal(err)
	}
	calls := 0
	r.observer.Secrets = statusSecretHook(func(context.Context) (string, error) {
		calls++
		if calls == 2 {
			if _, err := pool.Exec(ctx, `update platform.job set claimed_until=clock_timestamp()-interval '1 second' where tenant_id=$1`, p.TenantID); err != nil {
				t.Fatal(err)
			}
			match, _ := json.Marshal(map[string]string{"connection_id": r.connection, "provider_code": "meta-whatsapp", "event_type": "whatsapp.raw_webhook.received.v1"})
			jobs, err := w.jobs.ClaimScoped(ctx, "provider-events", "worker-new", time.Minute, 1, postgres.JobScope{TenantID: p.TenantID, JobType: "provider.webhook.received", SchemaVersion: 1, PayloadMatch: match})
			if err != nil || len(jobs) != 1 || jobs[0].Attempts != 2 {
				t.Fatal("takeover", jobs, err)
			}
		}
		return "test-app-secret", nil
	})
	result, err := w.ProcessOnce(ctx, p)
	if err == nil || result.Completed || result.FailureRecorded || result.InsertedObservations != 1 {
		t.Fatalf("stale result: %+v %v", result, err)
	}
	var retained bool
	if err := pool.QueryRow(ctx, `select j.completed_at is null and j.claimed_by='worker-new' and j.attempts=2 and e.state='received' from platform.job j join integration.webhook_event e on e.tenant_id=j.tenant_id and e.provider_event_id=$2 where j.tenant_id=$1`, p.TenantID, event).Scan(&retained); err != nil || !retained {
		t.Fatal("old worker damaged new claim", retained, err)
	}
	// After another simulated crash, a fresh claim completes without new effects.
	r.observer.Secrets = appSecretFixture{}
	if _, err := pool.Exec(ctx, `update platform.job set claimed_until=clock_timestamp()-interval '1 second' where tenant_id=$1`, p.TenantID); err != nil {
		t.Fatal(err)
	}
	result, err = w.ProcessOnce(ctx, p)
	if err != nil || !result.Completed || result.InsertedObservations != 0 {
		t.Fatalf("recovery after takeover: %+v %v", result, err)
	}
}

func TestStatusWorkerRejectsBeforeClaim(t *testing.T) {
	for _, mode := range []string{"permission", "tenant", "disabled", "cancelled", "busy"} {
		t.Run(mode, func(t *testing.T) {
			r, _, pool, p, _, _, _ := routerFixture(t, statusBatch("delivered", "1603086314", nil))
			tenant := p.TenantID
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			w, err := NewStatusWorker(r, "worker-guard", time.Minute, 0)
			if err != nil {
				t.Fatal(err)
			}
			switch mode {
			case "permission":
				p.Permissions = map[string]struct{}{}
			case "tenant":
				p.TenantID = "00000000-0000-0000-0000-000000000001"
			case "disabled":
				if _, err := pool.Exec(ctx, `update integration.provider_connection set state='disabled' where tenant_id=$1`, tenant); err != nil {
					t.Fatal(err)
				}
			case "cancelled":
				cancel()
			case "busy":
				w.busy <- struct{}{}
				defer func() { <-w.busy }()
			}
			result, err := w.ProcessOnce(ctx, p)
			if err == nil || result.Claimed || result.Completed {
				t.Fatalf("unauthorized claim: %+v %v", result, err)
			}
			var attempts int
			if err := pool.QueryRow(context.Background(), `select attempts from platform.job where tenant_id=$1`, tenant).Scan(&attempts); err != nil || attempts != 0 {
				t.Fatal("claim changed", attempts, err)
			}
		})
	}
}

func TestStatusWorkerClaimPayloadDriftCannotComplete(t *testing.T) {
	r, _, pool, p, _, event, _ := routerFixture(t, statusBatch("delivered", "1603086314", nil))
	ctx := context.Background()
	w, err := NewStatusWorker(r, "worker-drift", time.Minute, 0)
	if err != nil {
		t.Fatal(err)
	}
	calls := 0
	r.observer.Secrets = statusSecretHook(func(context.Context) (string, error) {
		calls++
		if calls == 2 {
			if _, err := pool.Exec(ctx, `update platform.job set payload=jsonb_set(payload,'{unexpected}','true') where tenant_id=$1`, p.TenantID); err != nil {
				t.Fatal(err)
			}
		}
		return "test-app-secret", nil
	})
	result, err := w.ProcessOnce(ctx, p)
	if err == nil || result.Completed || result.FailureRecorded || result.InsertedObservations != 1 {
		t.Fatalf("drift accepted: %+v %v", result, err)
	}
	var pending bool
	if err := pool.QueryRow(ctx, `select j.completed_at is null and j.payload ? 'unexpected' and e.state='received' from platform.job j join integration.webhook_event e on e.tenant_id=j.tenant_id and e.provider_event_id=$2 where j.tenant_id=$1`, p.TenantID, event).Scan(&pending); err != nil || !pending {
		t.Fatal("changed claim overwritten", pending, err)
	}
}
````

### FILE: `internal/whatsappbridge/status_router.go`

```yaml
block_id: "META-WHATSAPP:router:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "c8fdabb5b659d152908e47739a0094b137ea3a341bc52e9acab89d74b3113823"
variables: []
secrets_allowed: false
```

````go
package whatsappbridge

// AUTHORED composition of the existing Meta-adapted verifier and SQL owners.
// Optional verified customer-message projection uses the existing runtime.
// No automatic reply sending, new queue, or principal fabrication.
import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"elite.local/enterprise/internal/contactidentity"
	"elite.local/enterprise/internal/platform/identity"
	"github.com/jackc/pgx/v5"
)

var ErrStatusRouting = errors.New("whatsappbridge: retained status routing unavailable or unverified")

type StatusRouter struct {
	observer              StatusObserver
	connection, retention string
	key                   []byte
	maxRoutes             int
	conversation          *ConversationRoute
	busy                  chan struct{}
}

// key must be the outbound owner's exact HMAC key (resolved server-side).
// A real verified principal is required per call; this is not a public handler.
func NewStatusRouter(o *StatusObserver, connection, retention string, key []byte, maxRoutes int) (*StatusRouter, error) {
	if o == nil || o.Approvals == nil || o.Approvals.pool == nil || o.Approvals.pool.Config().MaxConns < 2 || o.Secrets == nil || !notificationEventID.MatchString(o.TenantID) || !webhookConnection.MatchString(connection) || !validDigest(retention) || len(key) < 32 || maxRoutes < 1 || maxRoutes > 8 || len(o.Profile) > 32768 || !json.Valid(o.Profile) || !filepath.IsAbs(o.Process.EvidenceDirectory) {
		return nil, ErrStatusRouting
	}
	copyObserver := *o
	copyObserver.Profile = append(json.RawMessage(nil), o.Profile...)
	return &StatusRouter{observer: copyObserver, connection: connection, retention: retention, key: append([]byte(nil), key...), maxRoutes: maxRoutes, busy: make(chan struct{}, 1)}, nil
}

type retainedStatus struct {
	Schema        string `json:"schema"`
	Body          []byte `json:"body"`
	Signature     string `json:"signature"`
	ProfileHash   string `json:"profile_sha256"`
	RetentionHash string `json:"retention_approval_sha256"`
	Count         int    `json:"verified_event_count"`
}

type statusIdentity struct{ Message, Recipient string }
type resolvedStatus struct {
	organization, appointment, event, delivery, anchor string
	identity                                           statusIdentity
	receipt                                            []byte
}

// ObserveRetained verifies all routes before invoking the existing owners.
// It never marks inbox/job complete. A caller must retain/retry/reconcile on
// error; earlier observer commits may exist and are idempotently replayable.
func (r *StatusRouter) ObserveRetained(ctx context.Context, p identity.Principal, eventID string) (int, error) {
	return r.observeRetained(ctx, p, eventID, nil)
}

// finish is internal-only: it may atomically ACK the exact claimed job/inbox
// after every admitted route was observed. Public routing still performs no ACK.
func (r *StatusRouter) observeRetained(ctx context.Context, p identity.Principal, eventID string, finish func(context.Context, pgx.Tx) error) (int, error) {
	if r == nil || p.TenantID != r.observer.TenantID || p.Subject == "" || len(p.Subject) > 255 || !p.Allowed("appointment:manage") || !strings.HasPrefix(eventID, "wa:") || !validDigest(strings.TrimPrefix(eventID, "wa:")) {
		return 0, ErrStatusRouting
	}
	select {
	case r.busy <- struct{}{}:
		defer func() { <-r.busy }()
	default:
		return 0, ErrStatusRouting
	}
	timeout := 30 * time.Second
	if r.conversation != nil {
		timeout = 110 * time.Second
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	o := &r.observer
	tx, err := o.Approvals.pool.Begin(ctx)
	if err != nil {
		return 0, ErrStatusRouting
	}
	defer tx.Rollback(ctx)
	// Hold connection/inbox scope stable through verification and observation.
	// Observer needs another pool connection; this router admits one call at a time.
	var payload []byte
	var bodyHash string
	var connectionOrg *string
	err = tx.QueryRow(ctx, `select e.payload,e.body_sha256_hex,c.organization_id
 from integration.webhook_event e join integration.provider_connection c
 on c.tenant_id=e.tenant_id and c.connection_id=e.connection_id and c.provider_code=e.provider_code
 where e.tenant_id=$1 and e.connection_id=$2 and e.provider_event_id=$3
 and c.state='active' and e.provider_code='meta-whatsapp'
 and e.event_type='whatsapp.raw_webhook.received.v1' and e.state in ('received','processed')
 and octet_length(e.payload::text)<=2097152 for share of c,e`, p.TenantID, r.connection, eventID).Scan(&payload, &bodyHash, &connectionOrg)
	if err != nil {
		return 0, ErrStatusRouting
	}
	var retained retainedStatus
	d := json.NewDecoder(bytes.NewReader(payload))
	d.DisallowUnknownFields()
	if d.Decode(&retained) != nil || d.Decode(new(any)) != io.EOF || retained.Schema != "elite-whatsapp-retained-webhook/v1" || len(retained.Body) == 0 || len(retained.Body) > 1<<20 || retained.ProfileHash != digest(o.Profile) || retained.RetentionHash != r.retention || retained.Count < 1 || retained.Count > 1000 || digest(retained.Body) != bodyHash || eventID != "wa:"+bodyHash {
		return 0, ErrStatusRouting
	}
	secret, err := o.Secrets.WhatsAppAppSecret(ctx)
	if err != nil {
		return 0, ErrStatusRouting
	}
	verifier := &WebhookReceiver{config: WebhookReceiverConfig{Profile: o.Profile, Process: o.Process}}
	verified, err := verifier.verify(ctx, "receive", retained.Body, retained.Signature, secret, map[string]string{})
	if err != nil || verified.BodyHash != bodyHash || verified.Count != retained.Count || verified.Challenge != "" {
		return 0, ErrStatusRouting
	}
	// This is identity projection AFTER the shared verifier has checked raw HMAC,
	// unique JSON keys, complete envelope, WABA/phone, events, and status semantics.
	// It is not a replacement signature/scope/normalization implementation.
	messages, messageCount, err := verifiedConversationMessages(retained.Body, p.TenantID, r.connection, r.maxRoutes)
	if err != nil || (messageCount > 0 && (r.conversation == nil || connectionOrg == nil || *connectionOrg != r.conversation.OrganizationID || !p.AllowedOrganization(*connectionOrg))) {
		return 0, ErrStatusRouting
	}
	identities, err := verifiedStatusIdentities(retained.Body, verified.Count, r.maxRoutes)
	if err != nil {
		return 0, ErrStatusRouting
	}
	root, err := os.OpenRoot(o.Process.EvidenceDirectory)
	if err != nil {
		return 0, ErrStatusRouting
	}
	defer root.Close()
	routes := make([]resolvedStatus, 0, len(identities))
	for _, id := range identities {
		messageHMAC, e1 := contactidentity.ExternalDigest(r.key, p.TenantID, "whatsapp", id.Message)
		recipientHMAC, e2 := contactidentity.ExternalDigest(r.key, p.TenantID, "whatsapp", id.Recipient)
		if e1 != nil || e2 != nil {
			return 0, ErrStatusRouting
		}
		rows, err := tx.Query(ctx, `select g.organization_id,g.appointment_id,g.approval_event_id,g.delivery_key,d.evidence_sha256_hex
 from communication.outbound_delivery d join communication.whatsapp_delivery_approval g
 on g.tenant_id=d.tenant_id and g.channel_code=d.channel_code and g.delivery_key=d.delivery_key
 join crm.lead l on l.tenant_id=g.tenant_id and l.lead_id=g.lead_id and l.organization_id=g.organization_id
 where d.tenant_id=$1 and d.channel_code='whatsapp' and d.state='accepted' and d.accepted_at is not null
 and d.provider_message_hmac=$2 and d.recipient_hmac=$3 and g.external_id_hmac=d.recipient_hmac
 and g.profile_sha256=$4 and d.request_sha256_hex=g.message_sha256
 and ($5::text is null or g.organization_id=$5) limit 2`, p.TenantID, messageHMAC, recipientHMAC, digest(o.Profile), connectionOrg)
		if err != nil {
			return 0, ErrStatusRouting
		}
		matches := 0
		var route resolvedStatus
		for rows.Next() {
			matches++
			if rows.Scan(&route.organization, &route.appointment, &route.event, &route.delivery, &route.anchor) != nil {
				rows.Close()
				return 0, ErrStatusRouting
			}
		}
		queryErr := rows.Err()
		rows.Close()
		if queryErr != nil || matches != 1 || !p.AllowedOrganization(route.organization) || !validDigest(route.anchor) || !validApprovalRoute(p.TenantID, route) {
			return 0, ErrStatusRouting
		}
		route.identity = id
		route.receipt, err = routedReceipt(root, p.TenantID, route)
		if err != nil {
			return 0, ErrStatusRouting
		}
		routes = append(routes, route)
	}
	// All routes and message shapes are validated before effects. Runtime replay
	// handles a crash after proposal commit and before the shared job ACK.
	for _, message := range messages {
		if _, err := r.conversation.Runtime.Handle(ctx, message); err != nil {
			return 0, ErrStatusRouting
		}
		if _, err := r.conversation.Proposals.Propose(ctx, p, eventID, message); err != nil {
			return 0, ErrStatusRouting
		}
	}
	inserted := 0
	for _, route := range routes {
		n, err := o.Observe(ctx, p, route.organization, route.appointment, route.event, route.receipt, []SignedStatusWebhook{{Body: retained.Body, Signature: retained.Signature}})
		inserted += n
		if err != nil {
			return inserted, ErrStatusRouting
		}
	}
	if finish != nil {
		if err := finish(ctx, tx); err != nil {
			return inserted, ErrStatusRouting
		}
	}
	if tx.Commit(ctx) != nil {
		return inserted, ErrStatusRouting
	}
	return inserted, nil
}

func verifiedStatusIdentities(raw []byte, count, limit int) ([]statusIdentity, error) {
	// Exact map keys, not struct decoding: encoding/json matches struct fields
	// case-insensitively, unlike the already-verified Python provider parser.
	var envelope map[string]json.RawMessage
	if json.Unmarshal(raw, &envelope) != nil {
		return nil, ErrStatusRouting
	}
	var entries []map[string]json.RawMessage
	if json.Unmarshal(envelope["entry"], &entries) != nil {
		return nil, ErrStatusRouting
	}
	seen := map[statusIdentity]bool{}
	total := 0
	for _, entry := range entries {
		var changes []map[string]json.RawMessage
		if json.Unmarshal(entry["changes"], &changes) != nil {
			return nil, ErrStatusRouting
		}
		for _, change := range changes {
			var value map[string]json.RawMessage
			if json.Unmarshal(change["value"], &value) != nil {
				return nil, ErrStatusRouting
			}
			var messages []json.RawMessage
			if v, ok := value["messages"]; ok {
				if json.Unmarshal(v, &messages) != nil {
					return nil, ErrStatusRouting
				}
			}
			total += len(messages)
			var statuses []map[string]json.RawMessage
			if v, ok := value["statuses"]; ok && json.Unmarshal(v, &statuses) != nil {
				return nil, ErrStatusRouting
			}
			for _, status := range statuses {
				var message, recipient string
				if json.Unmarshal(status["id"], &message) != nil || json.Unmarshal(status["recipient_id"], &recipient) != nil {
					return nil, ErrStatusRouting
				}
				total++
				if message == "" || recipient == "" || len(message) > 256 || len(recipient) > 256 {
					return nil, ErrStatusRouting
				}
				seen[statusIdentity{message, recipient}] = true
				if len(seen) > limit {
					return nil, ErrStatusRouting
				}
			}
		}
	}
	if total != count || total == 0 {
		return nil, ErrStatusRouting
	}
	ids := make([]statusIdentity, 0, len(seen))
	for id := range seen {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool {
		if ids[i].Message == ids[j].Message {
			return ids[i].Recipient < ids[j].Recipient
		}
		return ids[i].Message < ids[j].Message
	})
	return ids, nil
}

func routedReceipt(root *os.Root, tenant string, route resolvedStatus) ([]byte, error) {
	// Same deterministic location as Sender, but root-confined and bounded.
	file, err := root.Open(filepath.Join(digest([]byte(tenant+"\x00"+route.delivery)), "SEND_RECEIPT.json"))
	if err != nil {
		return nil, ErrStatusRouting
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil || !info.Mode().IsRegular() || info.Size() < 1 || info.Size() > 65536 {
		return nil, ErrStatusRouting
	}
	raw, err := io.ReadAll(io.LimitReader(file, 65537))
	if err != nil || len(raw) > 65536 || digest(raw) != route.anchor {
		return nil, ErrStatusRouting
	}
	var receipt struct {
		Message   string `json:"message_id_sha256"`
		Recipient string `json:"recipient_sha256"`
	}
	if json.Unmarshal(raw, &receipt) != nil || receipt.Message != digest([]byte(route.identity.Message)) || receipt.Recipient != digest([]byte(route.identity.Recipient)) {
		return nil, ErrStatusRouting
	}
	return raw, nil
}

func validApprovalRoute(tenant string, r resolvedStatus) bool {
	if r.appointment == "" {
		return r.delivery == r.event && strings.HasPrefix(r.event, "wa-reply:") && validDigest(strings.TrimPrefix(r.event, "wa-reply:"))
	}
	return r.delivery == AppointmentConfirmationDeliveryKey(tenant, r.event)
}
````

### FILE: `internal/whatsappbridge/status_router_test.go`

```yaml
block_id: "META-WHATSAPP:router-test:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "cb230a423fdd290355a5753cbb6f94ac407515d18bfeaa7d57ad66a87f2f4dd5"
variables: []
secrets_allowed: false
```

````go
package whatsappbridge

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"elite.local/enterprise/internal/contactidentity"
	"elite.local/enterprise/internal/platform/identity"
	"github.com/jackc/pgx/v5/pgxpool"
)

func alterStatusBatch(t *testing.T, b SignedStatusWebhook, edit func(map[string]any)) SignedStatusWebhook {
	t.Helper()
	var body map[string]any
	if json.Unmarshal(b.Body, &body) != nil {
		t.Fatal("fixture JSON")
	}
	edit(body)
	raw, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	mac := hmac.New(sha256.New, []byte("test-app-secret"))
	mac.Write(raw)
	return SignedStatusWebhook{Body: raw, Signature: "sha256=" + hex.EncodeToString(mac.Sum(nil))}
}
func statusValue(body map[string]any) map[string]any {
	return body["entry"].([]any)[0].(map[string]any)["changes"].([]any)[0].(map[string]any)["value"].(map[string]any)
}
func routerFixture(t *testing.T, batch SignedStatusWebhook) (*StatusRouter, *StatusObserver, *pgxpool.Pool, identity.Principal, AppointmentApprovalCommand, string, string) {
	t.Helper()
	pool := approvalPool(t)
	o, p, c, receipt := statusObserverFixture(t, pool)
	if _, err := pool.Exec(context.Background(), `insert into integration.provider_connection(tenant_id,connection_id,provider_code,secret_ref) values($1,'wa-primary','meta-whatsapp','META_APP_SECRET')`, p.TenantID); err != nil {
		t.Fatal(err)
	}
	config := receiverConfig(o)
	h, err := NewWebhookReceiver(config)
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(h)
	defer server.Close()
	response, _ := webhookRequest(t, server, batch)
	if response.StatusCode != 200 {
		t.Fatalf("receiver %d", response.StatusCode)
	}
	directory := filepath.Join(o.Process.EvidenceDirectory, digest([]byte(p.TenantID+"\x00"+c.Message.DeliveryKey)))
	if err := os.MkdirAll(directory, 0700); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(directory, "SEND_RECEIPT.json")
	if err := os.WriteFile(path, receipt, 0600); err != nil {
		t.Fatal(err)
	}
	r, err := NewStatusRouter(o, config.ConnectionID, config.RetentionApprovalSHA256, []byte("0123456789abcdef0123456789abcdef"), 4)
	if err != nil {
		t.Fatal(err)
	}
	return r, o, pool, p, c, "wa:" + digest(batch.Body), path
}

func TestStatusRouterRetainedToObservedAndReplay(t *testing.T) {
	batch := alterStatusBatch(t, statusBatch("delivered", "1603086314", nil), func(body map[string]any) {
		v := statusValue(body)
		v["statuses"] = append(v["statuses"].([]any), map[string]any{"id": "wamid.synthetic", "recipient_id": "5491112345678", "status": "sent", "timestamp": "1603086313"})
		// Case variants must not override the exact fields verified by Python.
		v["Statuses"] = []any{map[string]any{"id": "wamid.foreign", "recipient_id": "000", "status": "read", "timestamp": "1603086315"}}
	})
	r, o, pool, p, c, event, _ := routerFixture(t, batch)
	for _, expected := range []int{2, 0, 0} {
		n, err := r.ObserveRetained(context.Background(), p, event)
		if err != nil || n != expected {
			t.Fatalf("route inserted=%d expected=%d err=%v", n, expected, err)
		}
	}
	value, err := o.Approvals.ReadNotificationStatus(context.Background(), p, c.OrganizationID, c.AppointmentID, c.ConfirmationEventID)
	if err != nil || value.DeliveryStatus != "observed_delivered" || value.ProviderEventCount != 2 {
		t.Fatalf("read: %+v %v", value, err)
	}
	var untouched bool
	if err := pool.QueryRow(context.Background(), `select
 (select state='received' and processed_at is null from integration.webhook_event where tenant_id=$1 and provider_event_id=$2)
 and (select count(*)=1 and bool_and(completed_at is null) from platform.job where tenant_id=$1 and queue='provider-events')
 and (select attempt_count=1 and state='accepted' from communication.outbound_delivery where tenant_id=$1 and delivery_key=$3)`, p.TenantID, event, c.Message.DeliveryKey).Scan(&untouched); err != nil || !untouched {
		t.Fatalf("router changed receipt/job/send: %v %v", untouched, err)
	}
}

func TestStatusRouterRejectsBeforeObservation(t *testing.T) {
	for _, mode := range []string{"unknown", "recipient", "mixed", "mixed-unknown", "permission", "organization", "tenant", "disabled", "connection-org", "retention", "profile", "signature", "payload", "count", "missing-receipt", "changed-receipt", "wrong-key", "busy", "scope-changed", "ambiguous"} {
		t.Run(mode, func(t *testing.T) {
			batch := statusBatch("delivered", "1603086314", nil)
			if mode == "unknown" {
				batch = statusBatch("delivered", "1603086314", map[string]any{"id": "wamid.unknown"})
			}
			if mode == "recipient" {
				batch = statusBatch("delivered", "1603086314", map[string]any{"recipient_id": "9999999"})
			}
			if mode == "mixed" || mode == "mixed-unknown" {
				batch = alterStatusBatch(t, batch, func(body map[string]any) {
					v := statusValue(body)
					if mode == "mixed" {
						v["messages"] = []any{map[string]any{"id": "wamid.inbound", "from": "5491112345678", "type": "text", "timestamp": "1603086315", "text": map[string]any{"body": "Synthetic"}}}
					} else {
						v["statuses"] = append(v["statuses"].([]any), map[string]any{"id": "wamid.unknown", "recipient_id": "5491112345678", "status": "sent", "timestamp": "1603086313"})
					}
				})
			}
			r, o, pool, p, c, event, path := routerFixture(t, batch)
			originalTenant := p.TenantID
			execSQL := func(sql string, args ...any) {
				t.Helper()
				if _, err := pool.Exec(context.Background(), sql, args...); err != nil {
					t.Fatal(err)
				}
			}
			switch mode {
			case "permission":
				p.Permissions = map[string]struct{}{}
			case "organization":
				p.Organizations = map[string]struct{}{"other": {}}
			case "tenant":
				p.TenantID = "00000000-0000-0000-0000-000000000001"
			case "disabled":
				execSQL(`update integration.provider_connection set state='disabled' where tenant_id=$1`, p.TenantID)
			case "connection-org", "scope-changed":
				execSQL(`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,'store-2','store-2','Synthetic other','store')`, p.TenantID)
				if mode == "connection-org" {
					execSQL(`update integration.provider_connection set organization_id='store-2' where tenant_id=$1`, p.TenantID)
				} else {
					execSQL(`update crm.lead set organization_id='store-2' where tenant_id=$1`, p.TenantID)
				}
			case "retention":
				r.retention = strings.Repeat("b", 64)
			case "profile":
				r.observer.Profile = append(r.observer.Profile, ' ')
			case "signature":
				r.observer.Secrets = statusSecretHook(func(context.Context) (string, error) { return "wrong-secret", nil })
			case "payload":
				execSQL(`update integration.webhook_event set payload=jsonb_set(payload,'{body}','"e30="') where tenant_id=$1`, p.TenantID)
			case "count":
				execSQL(`update integration.webhook_event set payload=jsonb_set(payload,'{verified_event_count}','2') where tenant_id=$1`, p.TenantID)
			case "missing-receipt":
				if err := os.Remove(path); err != nil {
					t.Fatal(err)
				}
			case "changed-receipt":
				if err := os.WriteFile(path, []byte("{}"), 0600); err != nil {
					t.Fatal(err)
				}
			case "wrong-key":
				r.key = []byte(strings.Repeat("x", 32))
			case "busy":
				r.busy <- struct{}{}
				defer func() { <-r.busy }()
			case "ambiguous":
				otherEvent := "00000000-0000-0000-0000-000000000002"
				otherKey := AppointmentConfirmationDeliveryKey(p.TenantID, otherEvent)
				execSQL(`insert into crm.appointment_transition(tenant_id,transition_id,appointment_id,from_state,to_state,actor_subject) values($1,$2,$3,'requested','confirmed','synthetic-duplicate')`, p.TenantID, otherEvent, c.AppointmentID)
				execSQL(`insert into communication.whatsapp_appointment_approval select tenant_id,$3,channel_code,appointment_id,appointment_version,organization_id,$2,external_id_hmac,binding_version,lead_id,subject_id,policy_version,consent_id,consent_purpose,consent_evidence_sha256,message_sha256,profile_sha256,approval_request_sha256,evidence_sha256,approved_by,approved_at,expires_at from communication.whatsapp_appointment_approval where tenant_id=$1`, p.TenantID, otherEvent, otherKey)
				execSQL(`insert into communication.outbound_delivery select tenant_id,channel_code,$2,request_sha256_hex,recipient_hmac,state,attempt_count,locked_until,provider_message_hmac,evidence_sha256_hex,failure_code,accepted_at,created_at,updated_at from communication.outbound_delivery where tenant_id=$1`, p.TenantID, otherKey)
			}
			n, err := r.ObserveRetained(context.Background(), p, event)
			if err == nil || n != 0 {
				t.Fatalf("unsafe routing accepted: %d %v", n, err)
			}
			var count int
			if err := o.Approvals.pool.QueryRow(context.Background(), `select count(*) from communication.whatsapp_status_observation where tenant_id=$1`, originalTenant).Scan(&count); err != nil || count != 0 {
				t.Fatalf("partial write: %d %v", count, err)
			}
		})
	}
}

func TestStatusRouterProjectionBudgetAndExactKeys(t *testing.T) {
	batch := statusBatch("delivered", "1603086314", nil)
	ids, err := verifiedStatusIdentities(batch.Body, 1, 1)
	if err != nil || len(ids) != 1 || ids[0].Message != "wamid.synthetic" {
		t.Fatal(ids, err)
	}
	if _, err := verifiedStatusIdentities(batch.Body, 2, 1); err == nil {
		t.Fatal("count mismatch accepted")
	}
	if _, err := verifiedStatusIdentities(batch.Body, 1, 0); err == nil {
		t.Fatal("budget ignored")
	}
	upper := alterStatusBatch(t, batch, func(body map[string]any) {
		v := statusValue(body)
		v["Statuses"] = v["statuses"]
		delete(v, "statuses")
	})
	if _, err := verifiedStatusIdentities(upper.Body, 1, 1); err == nil {
		t.Fatal("case-only provider key accepted")
	}
	if _, err := NewStatusRouter(nil, "wa-primary", strings.Repeat("a", 64), []byte(strings.Repeat("x", 32)), 1); err == nil {
		t.Fatal("nil observer accepted")
	}
}

func TestStatusRouterMultipleMessagesRecoverPartialObservation(t *testing.T) {
	batch := alterStatusBatch(t, statusBatch("delivered", "1603086314", nil), func(body map[string]any) {
		v := statusValue(body)
		v["statuses"] = append(v["statuses"].([]any), map[string]any{"id": "wamid.zz-second", "recipient_id": "5491112345678", "status": "read", "timestamp": "1603086315"})
	})
	r, o, pool, p, c, event, path := routerFixture(t, batch)
	ctx := context.Background()
	otherEvent := "00000000-0000-0000-0000-000000000003"
	otherKey := AppointmentConfirmationDeliveryKey(p.TenantID, otherEvent)
	// A second explicitly synthetic, durably anchored send receipt. This is not
	// evidence of a second Meta API call or of live provider acceptance.
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var receipt map[string]any
	if json.Unmarshal(raw, &receipt) != nil {
		t.Fatal("synthetic receipt JSON")
	}
	receipt["message_id_sha256"] = digest([]byte("wamid.zz-second"))
	secondRaw, err := json.Marshal(receipt)
	if err != nil {
		t.Fatal(err)
	}
	messageHMAC, err := contactidentity.ExternalDigest(r.key, p.TenantID, "whatsapp", "wamid.zz-second")
	if err != nil {
		t.Fatal(err)
	}
	execSQL := func(sql string, args ...any) {
		t.Helper()
		if _, err := pool.Exec(ctx, sql, args...); err != nil {
			t.Fatal(err)
		}
	}
	execSQL(`insert into crm.appointment_transition(tenant_id,transition_id,appointment_id,from_state,to_state,actor_subject) values($1,$2,$3,'requested','confirmed','synthetic-second')`, p.TenantID, otherEvent, c.AppointmentID)
	execSQL(`insert into communication.whatsapp_appointment_approval select tenant_id,$3,channel_code,appointment_id,appointment_version,organization_id,$2,external_id_hmac,binding_version,lead_id,subject_id,policy_version,consent_id,consent_purpose,consent_evidence_sha256,message_sha256,profile_sha256,approval_request_sha256,evidence_sha256,approved_by,approved_at,expires_at from communication.whatsapp_appointment_approval where tenant_id=$1`, p.TenantID, otherEvent, otherKey)
	execSQL(`insert into communication.outbound_delivery select tenant_id,channel_code,$2,request_sha256_hex,recipient_hmac,state,attempt_count,locked_until,$3,$4,failure_code,accepted_at,created_at,updated_at from communication.outbound_delivery where tenant_id=$1`, p.TenantID, otherKey, messageHMAC, digest(secondRaw))
	directory := filepath.Join(o.Process.EvidenceDirectory, digest([]byte(p.TenantID+"\x00"+otherKey)))
	if err := os.MkdirAll(directory, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(directory, "SEND_RECEIPT.json"), secondRaw, 0600); err != nil {
		t.Fatal(err)
	}
	// Signature validation, first observer, then fail the second observer.
	calls := 0
	r.observer.Secrets = statusSecretHook(func(context.Context) (string, error) {
		calls++
		if calls == 3 {
			return "", errors.New("synthetic secret source unavailable")
		}
		return "test-app-secret", nil
	})
	n, err := r.ObserveRetained(ctx, p, event)
	if err == nil || n != 1 || calls != 3 {
		t.Fatalf("partial result not preserved: %d %d %v", n, calls, err)
	}
	var pending bool
	if err := pool.QueryRow(ctx, `select state='received' and processed_at is null from integration.webhook_event where tenant_id=$1 and provider_event_id=$2`, p.TenantID, event).Scan(&pending); err != nil || !pending {
		t.Fatalf("premature inbox completion: %v %v", pending, err)
	}
	r.observer.Secrets = appSecretFixture{}
	for _, expected := range []int{1, 0} {
		n, err := r.ObserveRetained(ctx, p, event)
		if err != nil || n != expected {
			t.Fatalf("recovery: %d %v", n, err)
		}
	}
	var total int
	if err := pool.QueryRow(ctx, `select count(*) from communication.whatsapp_status_observation where tenant_id=$1`, p.TenantID).Scan(&total); err != nil || total != 2 {
		t.Fatalf("duplicated/lost observations: %d %v", total, err)
	}
}
````

### FILE: `db/migrations/0052_whatsapp_status_routing.up.sql`

```yaml
block_id: "META-WHATSAPP:router-index-up:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "21a16e09525c2622c99f6df06c1f14cc58be4046d250e7afa41e6e86e9e28106"
variables: []
secrets_allowed: false
```

````sql
begin;
-- AUTHORED lookup support; not UNIQUE: ambiguity must be rejected, not erased.
create index whatsapp_status_message_route_idx
on communication.outbound_delivery(tenant_id,provider_message_hmac,recipient_hmac)
where channel_code='whatsapp' and state='accepted';
commit;
````

### FILE: `db/migrations/0052_whatsapp_status_routing.down.sql`

```yaml
block_id: "META-WHATSAPP:router-index-down:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "1354387e5e45e34b199ba14a4392203d7e40e97a5e50f2beaa476950bcc9caae"
variables: []
secrets_allowed: false
```

````sql
begin;
drop index communication.whatsapp_status_message_route_idx;
commit;
````



### FILE: `internal/whatsappbridge/webhook_receiver.go`
```yaml
block_id: "WHATSAPP:webhook-receiver:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "349ef32e52bcacdd216029cff466b976781ecf9ae83277a7b7fb016245a08724"
variables: []
secrets_allowed: false
```
````go
package whatsappbridge

// AUTHORED HTTP/SQL composition. Meta authentication remains in the pinned
// Python adapter; this file does not invent another provider signature scheme.
import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"time"

	"elite.local/enterprise/internal/providerintegration"
)

var ErrWebhookIngress = errors.New("whatsappbridge: ingress unavailable or unverified")
var webhookConnection = regexp.MustCompile(`^[a-z0-9][a-z0-9._-]{0,127}$`)

type VerifyTokenSource interface {
	WhatsAppVerifyToken(context.Context) (string, error)
}

// RetentionApprovalSHA256 links the project's approved raw-payload policy.
// Base64 is NOT encryption. Protect the database/backups and restrict/purge raw
// PII under that policy before exposing this handler. No default activation.
type WebhookReceiverConfig struct {
	TenantID, ConnectionID, RetentionApprovalSHA256 string
	Profile                                         json.RawMessage
	Process                                         Process
	Secrets                                         AppSecretSource
	Verification                                    VerifyTokenSource
	Store                                           providerintegration.Repository
	MaxConcurrent                                   int
}

type WebhookReceiver struct {
	config WebhookReceiverConfig
	slots  chan struct{}
}

func NewWebhookReceiver(c WebhookReceiverConfig) (*WebhookReceiver, error) {
	if !notificationEventID.MatchString(c.TenantID) || !webhookConnection.MatchString(c.ConnectionID) || !validDigest(c.RetentionApprovalSHA256) || len(c.Profile) > 32768 || !json.Valid(c.Profile) || c.Store == nil || c.Secrets == nil || c.Verification == nil || c.MaxConcurrent < 1 || c.MaxConcurrent > 16 {
		return nil, ErrWebhookIngress
	}
	c.Profile = append(json.RawMessage(nil), c.Profile...)
	return &WebhookReceiver{config: c, slots: make(chan struct{}, c.MaxConcurrent)}, nil
}

type ingressResult struct {
	Schema    string `json:"schema"`
	Binding   string `json:"binding_sha256"`
	BodyHash  string `json:"body_sha256"`
	Count     int    `json:"event_count"`
	Challenge string `json:"challenge"`
}

func (h *WebhookReceiver) verify(ctx context.Context, mode string, body []byte, signature, secret string, query map[string]string) (ingressResult, error) {
	var out ingressResult
	p := h.config.Process
	script := filepath.Join(p.AdapterDirectory, "whatsapp_cloud.py")
	if !filepath.IsAbs(p.AdapterDirectory) || exactFile(p.PythonExecutable, p.PythonSHA256) != nil || exactFile(script, p.AdapterSHA256) != nil || secret == "" || len(secret) > 16384 {
		return out, ErrWebhookIngress
	}
	frame := struct {
		Schema    string            `json:"schema"`
		Mode      string            `json:"mode"`
		Profile   json.RawMessage   `json:"profile"`
		Body      []byte            `json:"body"`
		Signature string            `json:"signature"`
		Secret    string            `json:"secret"`
		Query     map[string]string `json:"query"`
	}{"elite-whatsapp-ingress-bridge/v1", mode, h.config.Profile, body, signature, secret, query}
	input, err := json.Marshal(frame)
	if err != nil || len(input) > 2<<20 {
		return out, ErrWebhookIngress
	}
	cmd := exec.CommandContext(ctx, p.PythonExecutable, "-I", "-B", script, "--ingress-bridge")
	cmd.Env = []string{}
	if root := os.Getenv("SystemRoot"); root != "" {
		cmd.Env = append(cmd.Env, "SystemRoot="+root)
	}
	cmd.Stdin = bytes.NewReader(input)
	cmd.Stderr = io.Discard
	cmd.WaitDelay = time.Second
	var output boundedOutput
	cmd.Stdout = &output
	if cmd.Run() != nil || output.overflow {
		return out, ErrWebhookIngress
	}
	decoder := json.NewDecoder(bytes.NewReader(output.Bytes()))
	decoder.DisallowUnknownFields()
	if decoder.Decode(&out) != nil || decoder.Decode(new(any)) != io.EOF || out.Schema != "elite-whatsapp-ingress-result/v1" || out.Binding != digest(input) {
		return ingressResult{}, ErrWebhookIngress
	}
	return out, nil
}

// Mount on one server-configured callback URL. No URL/header/body selects the
// tenant or connection. Configure server read/header/idle timeouts and TLS/edge
// protections; redact subscription query tokens in access logs.
func (h *WebhookReceiver) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fail := func(code int) { w.WriteHeader(code); _, _ = io.WriteString(w, "WEBHOOK_NOT_ACCEPTED\n") }
	if h == nil {
		fail(503)
		return
	}
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		w.Header().Set("Allow", "GET, POST")
		fail(405)
		return
	}
	select {
	case h.slots <- struct{}{}:
		defer func() { <-h.slots }()
	default:
		fail(503)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	_ = http.NewResponseController(w).SetReadDeadline(time.Now().Add(5 * time.Second))
	if r.Method == http.MethodGet {
		if len(r.URL.RawQuery) > 2048 {
			fail(400)
			return
		}
		// ParseQuery reports malformed escaping instead of silently dropping it.
		values, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil || len(values) != 3 {
			fail(400)
			return
		}
		query := map[string]string{}
		for _, name := range []string{"hub.mode", "hub.verify_token", "hub.challenge"} {
			if len(values[name]) != 1 {
				fail(400)
				return
			}
			query[name] = values[name][0]
		}
		secret, err := h.config.Verification.WhatsAppVerifyToken(ctx)
		if err != nil {
			fail(503)
			return
		}
		result, err := h.verify(ctx, "subscribe", []byte{}, "", secret, query)
		if err != nil || result.Challenge != query["hub.challenge"] || result.Challenge == "" || result.Count != 0 || result.BodyHash != "" {
			fail(403)
			return
		}
		w.WriteHeader(200)
		_, _ = io.WriteString(w, result.Challenge)
		return
	}
	if r.URL.RawQuery != "" || len(r.Header.Values("X-Hub-Signature-256")) != 1 || len(r.Header.Get("X-Hub-Signature-256")) != 71 {
		fail(400)
		return
	}
	media, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || media != "application/json" || len(r.Header.Values("Content-Type")) != 1 || r.Header.Get("Content-Encoding") != "" {
		fail(415)
		return
	}
	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 1<<20))
	if err != nil || len(body) == 0 {
		fail(413)
		return
	}
	secret, err := h.config.Secrets.WhatsAppAppSecret(ctx)
	if err != nil {
		fail(503)
		return
	}
	signature := r.Header.Get("X-Hub-Signature-256")
	result, err := h.verify(ctx, "receive", body, signature, secret, map[string]string{})
	if err != nil || result.BodyHash != digest(body) || result.Count < 1 || result.Count > 1000 || result.Challenge != "" {
		fail(403)
		return
	}
	// Deterministic payload preserves exact original bytes for re-verification.
	// Admission does not mark the job processed or mutate CRM/delivery status.
	payload, err := json.Marshal(struct {
		Schema        string `json:"schema"`
		Body          []byte `json:"body"`
		Signature     string `json:"signature"`
		ProfileHash   string `json:"profile_sha256"`
		RetentionHash string `json:"retention_approval_sha256"`
		Count         int    `json:"verified_event_count"`
	}{"elite-whatsapp-retained-webhook/v1", body, signature, digest(h.config.Profile), h.config.RetentionApprovalSHA256, result.Count})
	if err != nil {
		fail(503)
		return
	}
	_, err = h.config.Store.AcceptWebhook(ctx, providerintegration.Receipt{TenantID: h.config.TenantID, ConnectionID: h.config.ConnectionID, ProviderCode: "meta-whatsapp", ProviderEventID: "wa:" + result.BodyHash, EventType: "whatsapp.raw_webhook.received.v1", BodyHash: result.BodyHash, Payload: payload})
	if err != nil {
		fail(503)
		return
	}
	w.WriteHeader(200)
	_, _ = io.WriteString(w, "EVENT_RECEIVED\n")
}
````

### FILE: `internal/whatsappbridge/webhook_receiver_test.go`
```yaml
block_id: "WHATSAPP:webhook-receiver-test:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "c15b1a8352d9b75d8e35670d502b3e625ce945b4864759389e281bbfbb0e6c40"
variables: []
secrets_allowed: false
```
````go
package whatsappbridge

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"elite.local/enterprise/internal/platform/postgres"
)

type verifyTokenFixture struct{}

func (verifyTokenFixture) WhatsAppVerifyToken(context.Context) (string, error) {
	return "test-verify-token", nil
}

func receiverConfig(o *StatusObserver) WebhookReceiverConfig {
	return WebhookReceiverConfig{TenantID: o.TenantID, ConnectionID: "wa-primary", RetentionApprovalSHA256: strings.Repeat("a", 64), Profile: o.Profile, Process: o.Process, Secrets: appSecretFixture{}, Verification: verifyTokenFixture{}, Store: postgres.NewProviderIntegration(o.Approvals.pool), MaxConcurrent: 4}
}

func webhookRequest(t *testing.T, server *httptest.Server, batch SignedStatusWebhook) (*http.Response, []byte) {
	t.Helper()
	req, _ := http.NewRequest(http.MethodPost, server.URL, bytes.NewReader(batch.Body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Hub-Signature-256", batch.Signature)
	res, err := server.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	return res, body
}

func TestWebhookReceiverPostgresReceiptAndObservation(t *testing.T) {
	pool := approvalPool(t)
	o, p, c, receipt := statusObserverFixture(t, pool)
	if _, err := pool.Exec(context.Background(), `insert into integration.provider_connection(tenant_id,connection_id,provider_code,secret_ref) values($1,'wa-primary','meta-whatsapp','META_APP_SECRET')`, p.TenantID); err != nil {
		t.Fatal(err)
	}
	config := receiverConfig(o)
	h, err := NewWebhookReceiver(config)
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(h)
	defer server.Close()
	res, err := server.Client().Get(server.URL + "?hub.mode=subscribe&hub.verify_token=test-verify-token&hub.challenge=12345")
	if err != nil {
		t.Fatal(err)
	}
	raw, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil || res.StatusCode != 200 || string(raw) != "12345" {
		t.Fatal("challenge", res.StatusCode, string(raw), err)
	}
	batch := statusBatch("delivered", "1603086314", nil)
	// Drop the transport only after the receiver has produced a successful ACK.
	// Its 200 was written to a recorder, never sent to the provider-side client.
	lost := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured := httptest.NewRecorder()
		h.ServeHTTP(captured, r)
		if captured.Code != 200 {
			http.Error(w, "fixture did not commit", 500)
			return
		}
		connection, _, err := w.(http.Hijacker).Hijack()
		if err != nil {
			t.Error("fixture hijack", err)
			return
		}
		_ = connection.Close()
	}))
	defer lost.Close()
	lostRequest, _ := http.NewRequest(http.MethodPost, lost.URL, bytes.NewReader(batch.Body))
	lostRequest.Header.Set("Content-Type", "application/json")
	lostRequest.Header.Set("X-Hub-Signature-256", batch.Signature)
	if lostResponse, err := lost.Client().Do(lostRequest); err == nil {
		lostResponse.Body.Close()
		t.Fatal("fixture unexpectedly delivered acknowledgement")
	}
	var committed int
	if err := pool.QueryRow(context.Background(), `select count(*) from integration.webhook_event where tenant_id=$1`, p.TenantID).Scan(&committed); err != nil || committed != 1 {
		t.Fatal("lost ACK without commit", committed, err)
	}
	var wg sync.WaitGroup
	results := make(chan int, 4)
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			req, _ := http.NewRequest(http.MethodPost, server.URL, bytes.NewReader(batch.Body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-Hub-Signature-256", batch.Signature)
			response, err := server.Client().Do(req)
			if err != nil {
				results <- 0
				return
			}
			defer response.Body.Close()
			_, _ = io.Copy(io.Discard, response.Body)
			results <- response.StatusCode
		}()
	}
	wg.Wait()
	close(results)
	for status := range results {
		if status != 200 {
			t.Fatalf("concurrent receipt status %d", status)
		}
	}
	var events, jobs int
	if err := pool.QueryRow(context.Background(), `select (select count(*) from integration.webhook_event where tenant_id=$1 and state='received' and processed_at is null),(select count(*) from platform.job where tenant_id=$1 and job_type='provider.webhook.received' and completed_at is null)`, p.TenantID).Scan(&events, &jobs); err != nil || events != 1 || jobs != 1 {
		t.Fatal("not durable or duplicated", events, jobs, err)
	}
	var stored []byte
	if err := pool.QueryRow(context.Background(), `select payload from integration.webhook_event where tenant_id=$1 and connection_id='wa-primary'`, p.TenantID).Scan(&stored); err != nil {
		t.Fatal(err)
	}
	var retained struct {
		Schema        string `json:"schema"`
		Body          []byte `json:"body"`
		Signature     string `json:"signature"`
		ProfileHash   string `json:"profile_sha256"`
		RetentionHash string `json:"retention_approval_sha256"`
		Count         int    `json:"verified_event_count"`
	}
	if json.Unmarshal(stored, &retained) != nil || retained.Schema != "elite-whatsapp-retained-webhook/v1" || !bytes.Equal(retained.Body, batch.Body) || retained.Signature != batch.Signature || retained.ProfileHash != digest(o.Profile) || retained.RetentionHash != config.RetentionApprovalSHA256 || retained.Count != 1 {
		t.Fatal("retained original differs")
	}
	if strings.Contains(string(stored), "test-app-secret") || strings.Contains(string(stored), "test-verify-token") {
		t.Fatal("secret persisted")
	}
	// Exercise the existing observer using the actual reloaded signed original.
	// Target routing is supplied by this synthetic fixture, not an automatic worker.
	if n, err := o.Observe(context.Background(), p, c.OrganizationID, c.AppointmentID, c.ConfirmationEventID, receipt, []SignedStatusWebhook{{Body: retained.Body, Signature: retained.Signature}}); n != 1 || err != nil {
		t.Fatal("observation from retained original", n, err)
	}
	value, err := o.Approvals.ReadNotificationStatus(context.Background(), p, c.OrganizationID, c.AppointmentID, c.ConfirmationEventID)
	if err != nil || value.DeliveryStatus != "observed_delivered" {
		t.Fatal(value, err)
	}
	// A lost HTTP acknowledgement is recoverable by replaying the exact original.
	res, raw = webhookRequest(t, server, batch)
	if res.StatusCode != 200 || string(raw) != "EVENT_RECEIVED\n" || res.Header.Get("Cache-Control") != "no-store" {
		t.Fatal(res.StatusCode, string(raw))
	}
	if _, err := pool.Exec(context.Background(), `update integration.provider_connection set state='disabled' where tenant_id=$1`, p.TenantID); err != nil {
		t.Fatal(err)
	}
	res, _ = webhookRequest(t, server, batch)
	if res.StatusCode != 503 {
		t.Fatal("disabled acknowledged", res.StatusCode)
	}
	if err := pool.QueryRow(context.Background(), `select (select count(*) from integration.webhook_event where tenant_id=$1),(select count(*) from platform.job where tenant_id=$1 and job_type='provider.webhook.received')`, p.TenantID).Scan(&events, &jobs); err != nil || events != 1 || jobs != 1 {
		t.Fatal(events, jobs, err)
	}
}

func TestWebhookReceiverPostgresRejectsBeforeReceipt(t *testing.T) {
	pool := approvalPool(t)
	o, p, _, _ := statusObserverFixture(t, pool)
	if _, err := pool.Exec(context.Background(), `insert into integration.provider_connection(tenant_id,connection_id,provider_code,secret_ref) values($1,'wa-primary','meta-whatsapp','META_APP_SECRET')`, p.TenantID); err != nil {
		t.Fatal(err)
	}
	config := receiverConfig(o)
	bad := config
	bad.RetentionApprovalSHA256 = ""
	if _, err := NewWebhookReceiver(bad); err == nil {
		t.Fatal("retention approval missing")
	}
	bad = config
	bad.MaxConcurrent = 0
	if _, err := NewWebhookReceiver(bad); err == nil {
		t.Fatal("missing concurrency budget")
	}
	h, err := NewWebhookReceiver(config)
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(h)
	defer server.Close()
	for _, mode := range []string{"tampered", "foreign_waba", "duplicate_signature", "media", "encoding", "oversize", "query", "invalid_json", "runtime_hash", "busy", "get_token", "get_duplicate", "get_malformed", "method"} {
		t.Run(mode, func(t *testing.T) {
			batch := statusBatch("read", "1603086315", nil)
			req, _ := http.NewRequest(http.MethodPost, server.URL, bytes.NewReader(batch.Body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-Hub-Signature-256", batch.Signature)
			want := 403
			sign := func(raw []byte) {
				mac := hmac.New(sha256.New, []byte("test-app-secret"))
				mac.Write(raw)
				req.Body = io.NopCloser(bytes.NewReader(raw))
				req.ContentLength = int64(len(raw))
				req.Header.Set("X-Hub-Signature-256", "sha256="+hex.EncodeToString(mac.Sum(nil)))
			}
			switch mode {
			case "tampered":
				req.Header.Set("X-Hub-Signature-256", "sha256="+strings.Repeat("0", 64))
			case "foreign_waba":
				sign(bytes.ReplaceAll(batch.Body, []byte("987654321"), []byte("987654322")))
			case "duplicate_signature":
				req.Header.Add("X-Hub-Signature-256", batch.Signature)
				want = 400
			case "media":
				req.Header.Set("Content-Type", "text/plain")
				want = 415
			case "encoding":
				req.Header.Set("Content-Encoding", "gzip")
				want = 415
			case "oversize":
				req.Body = io.NopCloser(strings.NewReader(strings.Repeat("x", (1<<20)+1)))
				req.ContentLength = (1 << 20) + 1
				want = 413
			case "query":
				req.URL.RawQuery = "tenant_id=other"
				want = 400
			case "invalid_json":
				sign([]byte(`{"object":`))
			case "runtime_hash":
				previous := h.config.Process.AdapterSHA256
				h.config.Process.AdapterSHA256 = strings.Repeat("0", 64)
				defer func() { h.config.Process.AdapterSHA256 = previous }()
			case "busy":
				for i := 0; i < cap(h.slots); i++ {
					h.slots <- struct{}{}
				}
				defer func() {
					for i := 0; i < cap(h.slots); i++ {
						<-h.slots
					}
				}()
				want = 503
			case "get_token":
				req.Method = http.MethodGet
				req.URL.RawQuery = "hub.mode=subscribe&hub.verify_token=wrong&hub.challenge=12345"
			case "get_duplicate":
				req.Method = http.MethodGet
				req.URL.RawQuery = "hub.mode=subscribe&hub.verify_token=test-verify-token&hub.challenge=1&hub.challenge=2"
				want = 400
			case "get_malformed":
				req.Method = http.MethodGet
				req.URL.RawQuery = "hub.mode=subscribe&hub.verify_token=%XX&hub.challenge=1"
				want = 400
			case "method":
				req.Method = http.MethodPut
				want = 405
			}
			res, err := server.Client().Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			raw, err := io.ReadAll(res.Body)
			if err != nil || res.StatusCode != want || string(raw) != "WEBHOOK_NOT_ACCEPTED\n" {
				t.Fatal(mode, res.StatusCode, string(raw), err)
			}
		})
	}
	var count int
	if err := pool.QueryRow(context.Background(), `select count(*) from integration.webhook_event where tenant_id=$1`, p.TenantID).Scan(&count); err != nil || count != 0 {
		t.Fatal("rejected callback persisted", count, err)
	}
}
````

### FILE: `internal/whatsappbridge/status_observer.go`
```yaml
block_id: "PY-META-WHATSAPP:durable-status-observer:v1"
operation: CREATE
provenance: AUTHORED
source: "Local verified integration of the admitted Meta signature boundary and PostgreSQL observation transactions; not upstream company code"
license: "LicenseRef-Workspace-Owner"
sha256: "36a91717fa18a63534497e148e0d505a9cde82bf220a8c41aad7a0a73104e9a6"
variables: []
secrets_allowed: false
```
````go
package whatsappbridge

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"elite.local/enterprise/internal/platform/identity"
	"github.com/jackc/pgx/v5"
)

var ErrStatusObservation = errors.New("whatsappbridge: status observation unverified, unavailable or divergent")

type AppSecretSource interface {
	WhatsAppAppSecret(context.Context) (string, error)
}
type SignedStatusWebhook struct {
	Body      []byte `json:"body"`
	Signature string `json:"signature"`
}

// StatusObserver is a service integration, not a public HTTP endpoint. The
// caller supplies an actually verified Principal and durable inbox raw bytes.
// The database, never the caller, supplies the trusted outbound receipt anchor.
type StatusObserver struct {
	TenantID         string
	Profile          json.RawMessage
	Process          Process
	ReconcilerSHA256 string
	Secrets          AppSecretSource
	Approvals        *PostgresAppointmentApprovals
}

type statusEvent struct {
	Kind          string `json:"kind"`
	MessageHash   string `json:"id_sha256"`
	RecipientHash string `json:"recipient_sha256"`
	AccountHash   string `json:"business_account_id_sha256"`
	PhoneHash     string `json:"phone_number_id_sha256"`
	Timestamp     string `json:"timestamp"`
	Status        string `json:"status"`
	EventHash     string `json:"event_sha256"`
	EventKey      string `json:"event_key_sha256"`
}
type statusResult struct {
	Schema  string          `json:"schema"`
	Binding string          `json:"binding_sha256"`
	Receipt json.RawMessage `json:"receipt"`
	Events  []statusEvent   `json:"events"`
}
type statusReceipt struct {
	Schema           string `json:"schema"`
	Anchor           string `json:"send_receipt_sha256"`
	ObservationsHash string `json:"observations_sha256"`
	Matching         int    `json:"matching_unique_events"`
	BusinessWrite    bool   `json:"automatic_business_write"`
	Resend           bool   `json:"resend_authorized"`
}

type statusQuery interface {
	QueryRow(context.Context, string, ...any) pgx.Row
}

func (o *StatusObserver) anchor(ctx context.Context, db statusQuery, p identity.Principal, organization, appointment, event string, locked bool) (string, string, error) {
	suffix := ""
	if locked {
		suffix = " for update of d for share of l"
	}
	var anchor, key string
	err := db.QueryRow(ctx, `select d.evidence_sha256_hex,g.delivery_key
 from communication.whatsapp_delivery_approval g
 join crm.lead l on l.tenant_id=g.tenant_id and l.lead_id=g.lead_id and l.organization_id=g.organization_id
 join communication.outbound_delivery d on d.tenant_id=g.tenant_id and d.delivery_key=g.delivery_key and d.channel_code=g.channel_code
 where g.tenant_id=$1 and g.organization_id=$2 and g.appointment_id=$3 and g.approval_event_id=$4
 and g.profile_sha256=$5 and d.state='accepted' and d.request_sha256_hex=g.message_sha256
 and d.recipient_hmac=g.external_id_hmac and d.accepted_at is not null`+suffix,
		p.TenantID, organization, appointment, event, digest(o.Profile)).Scan(&anchor, &key)
	if err != nil || !validDigest(anchor) {
		return "", "", ErrStatusObservation
	}
	return anchor, key, nil
}

// Observe validates all signatures in an isolated, hash-locked Python process,
// then stores observations atomically. It never changes the send fence, creates
// a resend, or acknowledges a durable inbox before the caller sees success.
func (o *StatusObserver) Observe(ctx context.Context, p identity.Principal, organization, appointment, event string, receipt []byte, batches []SignedStatusWebhook) (int, error) {
	if o == nil || o.Approvals == nil || o.Approvals.pool == nil || o.Secrets == nil || p.TenantID == "" || p.TenantID != o.TenantID || p.Subject == "" || len(p.Subject) > 255 || !p.Allowed("appointment:manage") || !p.AllowedOrganization(organization) || !((appointment != "" && notificationEventID.MatchString(event)) || (appointment == "" && strings.HasPrefix(event, "wa-reply:") && validDigest(strings.TrimPrefix(event, "wa-reply:")))) || len(o.Profile) > 32768 || !json.Valid(o.Profile) || len(receipt) == 0 || len(receipt) > 65536 || len(batches) == 0 || len(batches) > 8 {
		return 0, ErrStatusObservation
	}
	for _, b := range batches {
		if len(b.Body) == 0 || len(b.Body) > 1<<20 || len(b.Signature) != 71 {
			return 0, ErrStatusObservation
		}
	}
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	anchor, key, err := o.anchor(ctx, o.Approvals.pool, p, organization, appointment, event, false)
	if err != nil || digest(receipt) != anchor {
		return 0, ErrStatusObservation
	}
	proc := o.Process
	script := filepath.Join(proc.AdapterDirectory, "whatsapp_cloud.py")
	if !filepath.IsAbs(proc.AdapterDirectory) || !filepath.IsAbs(proc.EvidenceDirectory) || exactFile(proc.PythonExecutable, proc.PythonSHA256) != nil || exactFile(script, proc.AdapterSHA256) != nil || exactFile(filepath.Join(proc.AdapterDirectory, "status_reconciliation.py"), o.ReconcilerSHA256) != nil {
		return 0, ErrStatusObservation
	}
	secret, err := o.Secrets.WhatsAppAppSecret(ctx)
	if err != nil || secret == "" || len(secret) > 16384 {
		return 0, ErrStatusObservation
	}
	input, err := json.Marshal(struct {
		Schema    string                `json:"schema"`
		Profile   json.RawMessage       `json:"profile"`
		Receipt   []byte                `json:"send_receipt"`
		Anchor    string                `json:"expected_send_receipt_sha256"`
		Webhooks  []SignedStatusWebhook `json:"webhooks"`
		Secret    string                `json:"app_secret"`
		Directory string                `json:"evidence_directory"`
	}{"elite-whatsapp-status-bridge/v1", o.Profile, receipt, anchor, batches, secret, proc.EvidenceDirectory})
	if err != nil || len(input) > 12<<20 {
		return 0, ErrStatusObservation
	}
	cmd := exec.CommandContext(ctx, proc.PythonExecutable, "-I", "-B", script, "--status-bridge")
	cmd.Env = []string{}
	if root := os.Getenv("SystemRoot"); root != "" {
		cmd.Env = append(cmd.Env, "SystemRoot="+root)
	}
	cmd.Stdin = bytes.NewReader(input)
	cmd.Stderr = io.Discard
	cmd.WaitDelay = 2 * time.Second
	stdout := statusOutput{}
	cmd.Stdout = &stdout
	if cmd.Run() != nil {
		return 0, ErrStatusObservation
	}
	var result statusResult
	decoder := json.NewDecoder(bytes.NewReader(stdout.Bytes()))
	decoder.DisallowUnknownFields()
	if decoder.Decode(&result) != nil || decoder.Decode(new(any)) != io.EOF || result.Schema != "elite-whatsapp-status-result/v1" || result.Binding != digest(input) || len(result.Events) > 8000 {
		return 0, ErrStatusObservation
	}
	var evidence statusReceipt
	if len(result.Receipt) > 65536 || json.Unmarshal(result.Receipt, &evidence) != nil || evidence.Schema != "elite-whatsapp-status-reconciliation/v1" || evidence.Anchor != anchor || !validDigest(evidence.ObservationsHash) || evidence.Matching != len(result.Events) || evidence.BusinessWrite || evidence.Resend {
		return 0, ErrStatusObservation
	}
	for _, e := range result.Events {
		stamp, err := strconv.ParseInt(e.Timestamp, 10, 64)
		if err != nil || stamp < 1 || stamp > 999999999999 || strconv.FormatInt(stamp, 10) != e.Timestamp || !validDigest(e.EventKey) || !validDigest(e.EventHash) || e.Kind != "status" {
			return 0, ErrStatusObservation
		}
		switch e.Status {
		case "sent", "delivered", "read", "failed", "deleted":
		default:
			return 0, ErrStatusObservation
		}
	}
	tx, err := o.Approvals.pool.Begin(ctx)
	if err != nil {
		return 0, ErrStatusObservation
	}
	defer tx.Rollback(ctx)
	current, currentKey, err := o.anchor(ctx, tx, p, organization, appointment, event, true)
	if err != nil || current != anchor || currentKey != key {
		return 0, ErrStatusObservation
	}
	if len(result.Events) == 0 {
		return 0, nil
	}
	_, err = tx.Exec(ctx, `insert into communication.whatsapp_status_batch(tenant_id,delivery_key,observations_sha256,send_receipt_sha256,verification_receipt,observed_by) values($1,$2,$3,$4,$5,$6) on conflict do nothing`, p.TenantID, key, evidence.ObservationsHash, anchor, result.Receipt, p.Subject)
	if err != nil {
		return 0, ErrStatusObservation
	}
	inserted := 0
	for _, e := range result.Events {
		timestamp, _ := strconv.ParseInt(e.Timestamp, 10, 64)
		tag, err := tx.Exec(ctx, `insert into communication.whatsapp_status_observation(tenant_id,delivery_key,event_key_sha256,event_sha256,provider_timestamp,provider_status,observations_sha256) values($1,$2,$3,$4,$5,$6,$7) on conflict do nothing`, p.TenantID, key, e.EventKey, e.EventHash, timestamp, e.Status, evidence.ObservationsHash)
		if err != nil {
			return 0, ErrStatusObservation
		}
		if tag.RowsAffected() == 1 {
			inserted++
			continue
		}
		var hash, state string
		var storedTime int64
		if tx.QueryRow(ctx, `select event_sha256,provider_status,provider_timestamp from communication.whatsapp_status_observation where tenant_id=$1 and delivery_key=$2 and event_key_sha256=$3`, p.TenantID, key, e.EventKey).Scan(&hash, &state, &storedTime) != nil || hash != e.EventHash || state != e.Status || storedTime != timestamp {
			return 0, ErrStatusObservation
		}
	}
	if tx.Commit(ctx) != nil {
		return 0, ErrStatusObservation
	}
	return inserted, nil
}

type statusOutput struct{ bytes.Buffer }

func (b *statusOutput) Write(p []byte) (int, error) {
	if b.Len()+len(p) > 8<<20 {
		return 0, ErrStatusObservation
	}
	return b.Buffer.Write(p)
}
````

### FILE: `internal/whatsappbridge/status_observer_test.go`
```yaml
block_id: "PY-META-WHATSAPP:durable-status-observer-test:v1"
operation: CREATE
provenance: AUTHORED
source: "Local verified integration of the admitted Meta signature boundary and PostgreSQL observation transactions; not upstream company code"
license: "LicenseRef-Workspace-Owner"
sha256: "6ed7e3c4c777baf6ca4c031a0d94b9adfe9d712634898ffbe55dbcd971e0db35"
variables: []
secrets_allowed: false
```
````go
package whatsappbridge

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"elite.local/enterprise/internal/outbounddelivery"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

type appSecretFixture struct{}

type statusSecretHook func(context.Context) (string, error)

func (f statusSecretHook) WhatsAppAppSecret(ctx context.Context) (string, error) { return f(ctx) }

func (appSecretFixture) WhatsAppAppSecret(context.Context) (string, error) {
	return "test-app-secret", nil
}

func statusObserverFixture(t *testing.T, pool *pgxpool.Pool) (*StatusObserver, identity.Principal, AppointmentApprovalCommand, []byte) {
	t.Helper()
	python := os.Getenv("ELITE_WHATSAPP_PYTHON")
	if python == "" {
		t.Fatal("ELITE_WHATSAPP_PYTHON required for actual isolated status verifier")
	}
	directory, err := filepath.Abs("../../whatsapp_cloud")
	if err != nil {
		t.Fatal(err)
	}
	// Reuse the Python send-bridge fixture; its only provider transport is local.
	code := `import sys,json,base64
sys.path.insert(0,sys.argv[1])
from test_status_reconciliation import StatusReconciliationTests,profile
c=StatusReconciliationTests();c.setUp()
try: print(json.dumps({"profile":profile(),"receipt":base64.b64encode(c.receipt).decode(),"message":__import__('test_whatsapp_cloud').message()}))
finally: c.doCleanups()
`
	raw, err := exec.Command(python, "-I", "-B", "-c", code, directory).Output()
	if err != nil {
		t.Fatal("synthetic send fixture", err)
	}
	var data struct {
		Profile json.RawMessage `json:"profile"`
		Receipt []byte          `json:"receipt"`
		Message json.RawMessage `json:"message"`
	}
	if json.Unmarshal(raw, &data) != nil {
		t.Fatal("fixture JSON")
	}
	p, c, approvals := approvalFixtureData(t, pool)
	c.Message.Text = string(data.Message)
	c.ProfileSHA256 = digest(data.Profile)
	if _, err = approvals.Approve(context.Background(), p, c); err != nil {
		t.Fatal(err)
	}
	store, err := postgres.NewOutboundDeliveryStore(pool, []byte("0123456789abcdef0123456789abcdef"), time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	hash, _ := outbounddelivery.MessageSHA256(c.Message)
	if _, err = store.Claim(context.Background(), c.Message, hash); err != nil {
		t.Fatal(err)
	}
	if err = store.Complete(context.Background(), c.Message, hash, outbounddelivery.Receipt{ProviderMessageID: "wamid.synthetic", EvidenceSHA256: digest(data.Receipt), AcceptedAt: time.Now()}); err != nil {
		t.Fatal(err)
	}
	fileHash := func(path string) string {
		b, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		return digest(b)
	}
	o := &StatusObserver{TenantID: p.TenantID, Profile: data.Profile, Approvals: approvals, Secrets: appSecretFixture{}, ReconcilerSHA256: fileHash(filepath.Join(directory, "status_reconciliation.py")), Process: Process{PythonExecutable: python, PythonSHA256: fileHash(python), AdapterDirectory: directory, AdapterSHA256: fileHash(filepath.Join(directory, "whatsapp_cloud.py")), EvidenceDirectory: t.TempDir()}}
	return o, p, c, data.Receipt
}

func statusBatch(status, timestamp string, extra map[string]any) SignedStatusWebhook {
	event := map[string]any{"id": "wamid.synthetic", "recipient_id": "5491112345678", "status": status, "timestamp": timestamp}
	for k, v := range extra {
		event[k] = v
	}
	raw, _ := json.Marshal(map[string]any{"object": "whatsapp_business_account", "entry": []any{map[string]any{"id": "987654321", "changes": []any{map[string]any{"field": "messages", "value": map[string]any{"messaging_product": "whatsapp", "metadata": map[string]any{"phone_number_id": "123456789"}, "statuses": []any{event}}}}}}})
	mac := hmac.New(sha256.New, []byte("test-app-secret"))
	mac.Write(raw)
	return SignedStatusWebhook{Body: raw, Signature: "sha256=" + hex.EncodeToString(mac.Sum(nil))}
}

func TestStatusObserverPostgresConcurrentAndRead(t *testing.T) {
	pool := approvalPool(t)
	o, p, c, receipt := statusObserverFixture(t, pool)
	batches := []SignedStatusWebhook{statusBatch("read", "1603086315", nil), statusBatch("sent", "1603086313", nil), statusBatch("delivered", "1603086314", nil)}
	var wg sync.WaitGroup
	errs := make(chan error, 4)
	counts := make(chan int, 4)
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			n, e := o.Observe(context.Background(), p, c.OrganizationID, c.AppointmentID, c.ConfirmationEventID, receipt, batches)
			counts <- n
			errs <- e
		}()
	}
	wg.Wait()
	close(errs)
	close(counts)
	for e := range errs {
		if e != nil {
			t.Fatal(e)
		}
	}
	total := 0
	for n := range counts {
		total += n
	}
	if total != 3 {
		t.Fatalf("inserted=%d", total)
	}
	value, err := o.Approvals.ReadNotificationStatus(context.Background(), p, c.OrganizationID, c.AppointmentID, c.ConfirmationEventID)
	if err != nil || value.FenceState != "accepted" || value.DeliveryStatus != "observed_read" || value.ProviderEventCount != 3 || value.ProviderTimestamp == nil || *value.ProviderTimestamp != 1603086315 {
		t.Fatal(value, err)
	}
	verifier, sign := notificationIssuer(t)
	config, err := pgxpool.ParseConfig(os.Getenv("ELITE_WHATSAPP_TEST_DATABASE_URL"))
	if err != nil {
		t.Fatal(err)
	}
	config.ConnConfig.RuntimeParams["default_transaction_read_only"] = "on"
	readPool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}
	defer readPool.Close()
	readApprovals := *o.Approvals
	readApprovals.pool = readPool
	module := &AppointmentNotificationModule{approvals: &readApprovals, sender: &Sender{TenantID: p.TenantID}}
	mux := http.NewServeMux()
	module.registerNotificationStatus(mux, verifier)
	server := httptest.NewServer(mux)
	defer server.Close()
	req, _ := http.NewRequest(http.MethodGet, server.URL+"/v1/franchise/appointments/"+c.AppointmentID+"/whatsapp-confirmation?organization_id="+c.OrganizationID+"&confirmation_event_id="+c.ConfirmationEventID, nil)
	req.Header.Set("Authorization", "Bearer "+sign(p))
	response, err := server.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	var httpValue NotificationStatus
	if response.StatusCode != 200 || response.Header.Get("Cache-Control") != "no-store" || json.NewDecoder(response.Body).Decode(&httpValue) != nil || httpValue.DeliveryStatus != "observed_read" || httpValue.ProviderEventCount != 3 {
		t.Fatal("HTTP read-only observation", httpValue, response.StatusCode)
	}
	if n, err := o.Observe(context.Background(), p, c.OrganizationID, c.AppointmentID, c.ConfirmationEventID, receipt, []SignedStatusWebhook{statusBatch("failed", "1603086315", nil)}); err != nil || n != 1 {
		t.Fatal(n, err)
	}
	value, err = o.Approvals.ReadNotificationStatus(context.Background(), p, c.OrganizationID, c.AppointmentID, c.ConfirmationEventID)
	if err != nil || value.DeliveryStatus != "ambiguous_latest_timestamp" || value.ProviderEventCount != 4 {
		t.Fatal(value, err)
	}
	var events, attempts int
	if err = pool.QueryRow(context.Background(), `select attempt_count,(select count(*) from communication.outbound_delivery_event e where e.tenant_id=d.tenant_id) from communication.outbound_delivery d where tenant_id=$1`, p.TenantID).Scan(&attempts, &events); err != nil || attempts != 1 || events != 2 {
		t.Fatal("send fence changed", attempts, events, err)
	}
	// Evidence is append-only and does not leak into the public response.
	if _, err = pool.Exec(context.Background(), `update communication.whatsapp_status_observation set provider_status='sent' where tenant_id=$1`, p.TenantID); err == nil {
		t.Fatal("immutable evidence updated")
	}
	encoded, _ := json.Marshal(value)
	for _, private := range []string{"5491112345678", "wamid.synthetic", "test-app-secret", digest(receipt)} {
		if strings.Contains(string(encoded), private) {
			t.Fatal("private data in response")
		}
	}
}

func TestStatusObserverPostgresRejections(t *testing.T) {
	pool := approvalPool(t)
	for _, mode := range []string{"signature", "receipt", "tenant", "permission", "organization", "recipient", "adapter", "reconciler", "unknown", "reassigned", "profile"} {
		t.Run(mode, func(t *testing.T) {
			o, p, c, receipt := statusObserverFixture(t, pool)
			batches := []SignedStatusWebhook{statusBatch("delivered", "1603086314", nil)}
			switch mode {
			case "signature":
				batches[0].Signature = "sha256=" + strings.Repeat("0", 64)
			case "receipt":
				receipt = append(receipt, ' ')
			case "tenant":
				p.TenantID = "00000000-0000-0000-0000-000000000000"
			case "permission":
				p.Permissions = map[string]struct{}{}
			case "organization":
				c.OrganizationID = "other"
			case "recipient":
				batches[0] = statusBatch("delivered", "1603086314", map[string]any{"recipient_id": "5491112345679"})
			case "adapter":
				o.Process.AdapterSHA256 = strings.Repeat("0", 64)
			case "reconciler":
				o.ReconcilerSHA256 = strings.Repeat("0", 64)
			case "profile":
				o.Profile = json.RawMessage(`{"wrong":true}`)
			case "unknown":
				if _, err := pool.Exec(context.Background(), `update communication.outbound_delivery set state='unknown' where tenant_id=$1`, o.TenantID); err != nil {
					t.Fatal(err)
				}
			case "reassigned":
				if _, err := pool.Exec(context.Background(), `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,'other','other','Other','store')`, o.TenantID); err != nil {
					t.Fatal(err)
				}
				if _, err := pool.Exec(context.Background(), `update crm.lead set organization_id='other' where tenant_id=$1`, o.TenantID); err != nil {
					t.Fatal(err)
				}
			}
			if n, err := o.Observe(context.Background(), p, c.OrganizationID, c.AppointmentID, c.ConfirmationEventID, receipt, batches); !errors.Is(err, ErrStatusObservation) || n != 0 {
				t.Fatal(n, err)
			}
			var count int
			if err := pool.QueryRow(context.Background(), `select count(*) from communication.whatsapp_status_batch where tenant_id=$1`, o.TenantID).Scan(&count); err != nil || count != 0 {
				t.Fatal("partial evidence", count, err)
			}
		})
	}
}

func TestStatusObserverPostgresConflictRollsBackBatch(t *testing.T) {
	pool := approvalPool(t)
	o, p, c, receipt := statusObserverFixture(t, pool)
	original := statusBatch("failed", "1603086314", map[string]any{"errors": []any{map[string]any{"code": 1}}})
	if _, err := o.Observe(context.Background(), p, c.OrganizationID, c.AppointmentID, c.ConfirmationEventID, receipt, []SignedStatusWebhook{original}); err != nil {
		t.Fatal(err)
	}
	divergent := statusBatch("failed", "1603086314", map[string]any{"errors": []any{map[string]any{"code": 2}}})
	if _, err := o.Observe(context.Background(), p, c.OrganizationID, c.AppointmentID, c.ConfirmationEventID, receipt, []SignedStatusWebhook{statusBatch("sent", "1603086313", nil), divergent}); err == nil {
		t.Fatal("divergence accepted")
	}
	var events, batches int
	if err := pool.QueryRow(context.Background(), `select (select count(*) from communication.whatsapp_status_observation where tenant_id=$1),(select count(*) from communication.whatsapp_status_batch where tenant_id=$1)`, p.TenantID).Scan(&events, &batches); err != nil || events != 1 || batches != 1 {
		t.Fatal("not atomic", events, batches, err)
	}
}

// The secret is requested after the first database anchor lookup and before the
// isolated verifier runs. This hook changes durable state in that exact window;
// no sleeps, replacement verifier or weakened production guard are involved.
func TestStatusObserverPostgresRechecksAfterVerification(t *testing.T) {
	pool := approvalPool(t)
	for _, mode := range []string{"lead_scope", "fence_state", "receipt_anchor"} {
		t.Run(mode, func(t *testing.T) {
			o, p, c, receipt := statusObserverFixture(t, pool)
			called := false
			var mutationErr error
			o.Secrets = statusSecretHook(func(ctx context.Context) (secret string, err error) {
				defer func() { mutationErr = err }()
				called = true
				switch mode {
				case "lead_scope":
					if _, err := pool.Exec(ctx, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,'other','other','Other','store')`, p.TenantID); err != nil {
						return "", err
					}
					if _, err := pool.Exec(ctx, `update crm.lead set organization_id='other' where tenant_id=$1`, p.TenantID); err != nil {
						return "", err
					}
				case "fence_state":
					if _, err := pool.Exec(ctx, `update communication.outbound_delivery set state='unknown' where tenant_id=$1`, p.TenantID); err != nil {
						return "", err
					}
				case "receipt_anchor":
					if _, err := pool.Exec(ctx, `update communication.outbound_delivery set evidence_sha256_hex=repeat('f',64) where tenant_id=$1`, p.TenantID); err != nil {
						return "", err
					}
				}
				return "test-app-secret", nil
			})
			n, err := o.Observe(context.Background(), p, c.OrganizationID, c.AppointmentID, c.ConfirmationEventID, receipt, []SignedStatusWebhook{statusBatch("read", "1603086315", nil)})
			if !called || mutationErr != nil || n != 0 || !errors.Is(err, ErrStatusObservation) {
				t.Fatal("stale authority accepted or mutation not exercised", called, mutationErr, n, err)
			}
			assertStatusRows(t, pool, p.TenantID, 0, 0)
		})
	}
}

func assertStatusRows(t *testing.T, pool *pgxpool.Pool, tenant string, wantEvents, wantBatches int) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var events, batches int
	err := pool.QueryRow(ctx, `select (select count(*) from communication.whatsapp_status_observation where tenant_id=$1),(select count(*) from communication.whatsapp_status_batch where tenant_id=$1)`, tenant).Scan(&events, &batches)
	if err != nil || events != wantEvents || batches != wantBatches {
		t.Fatal("unexpected durable observations", events, batches, err)
	}
}

func TestStatusObserverPostgresCanceledWriteAndRecovery(t *testing.T) {
	pool := approvalPool(t)
	o, p, c, receipt := statusObserverFixture(t, pool)
	batches := []SignedStatusWebhook{statusBatch("read", "1603086315", nil)}
	lock, err := pool.Begin(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer lock.Rollback(context.Background())
	// Block the event INSERT after its parent batch has been inserted. The lock
	// is confined to the dedicated loopback test database guarded by approvalPool.
	if _, err := lock.Exec(context.Background(), `lock table communication.whatsapp_status_observation in exclusive mode`); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	type outcome struct {
		n   int
		err error
	}
	done := make(chan outcome, 1)
	go func() {
		n, err := o.Observe(ctx, p, c.OrganizationID, c.AppointmentID, c.ConfirmationEventID, receipt, batches)
		done <- outcome{n, err}
	}()
	probe, stop := context.WithTimeout(context.Background(), 5*time.Second)
	defer stop()
	ticker := time.NewTicker(20 * time.Millisecond)
	defer ticker.Stop()
	for {
		var blocked bool
		err := pool.QueryRow(probe, `select exists(select 1 from pg_stat_activity where datname=current_database() and wait_event_type='Lock' and query like 'insert into communication.whatsapp_status_observation%')`).Scan(&blocked)
		if err != nil {
			t.Fatal("could not prove blocked INSERT", err)
		}
		if blocked {
			break
		}
		select {
		case early := <-done:
			t.Fatal("observer ended before blocked INSERT", early)
		case <-probe.Done():
			t.Fatal("blocked INSERT not observed")
		case <-ticker.C:
		}
	}
	cancel()
	select {
	case result := <-done:
		if result.n != 0 || !errors.Is(result.err, ErrStatusObservation) {
			t.Fatal("canceled write reported success", result)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("cancellation did not release observer")
	}
	if err := lock.Rollback(context.Background()); err != nil {
		t.Fatal(err)
	}
	assertStatusRows(t, pool, p.TenantID, 0, 0)
	if n, err := o.Observe(context.Background(), p, c.OrganizationID, c.AppointmentID, c.ConfirmationEventID, receipt, batches); n != 1 || err != nil {
		t.Fatal("recovery failed", n, err)
	}
	assertStatusRows(t, pool, p.TenantID, 1, 1)
}

func TestStatusObserverPostgresReplayAcrossConnectionRestart(t *testing.T) {
	pool := approvalPool(t)
	o, p, c, receipt := statusObserverFixture(t, pool)
	batches := []SignedStatusWebhook{statusBatch("delivered", "1603086314", nil)}
	if n, err := o.Observe(context.Background(), p, c.OrganizationID, c.AppointmentID, c.ConfirmationEventID, receipt, batches); n != 1 || err != nil {
		t.Fatal(n, err)
	}
	pool.Close()
	// Fresh database connections and service values; no in-memory dedup cache.
	reopened := approvalPool(t)
	approvals := *o.Approvals
	approvals.pool = reopened
	observer := *o
	observer.Approvals = &approvals
	if n, err := observer.Observe(context.Background(), p, c.OrganizationID, c.AppointmentID, c.ConfirmationEventID, receipt, batches); n != 0 || err != nil {
		t.Fatal("replay after connection restart", n, err)
	}
	assertStatusRows(t, reopened, p.TenantID, 1, 1)
	var attempts int
	if err := reopened.QueryRow(context.Background(), `select attempt_count from communication.outbound_delivery where tenant_id=$1`, p.TenantID).Scan(&attempts); err != nil || attempts != 1 {
		t.Fatal("replay changed send attempts", attempts, err)
	}
}
````

### FILE: `db/migrations/0051_whatsapp_status_observations.up.sql`
```yaml
block_id: "PY-META-WHATSAPP:status-observations-up:v1"
operation: CREATE
provenance: AUTHORED
source: "Local verified integration of the admitted Meta signature boundary and PostgreSQL observation transactions; not upstream company code"
license: "LicenseRef-Workspace-Owner"
sha256: "00514d0f9c3285e38a10839d291edc966b96f235c1e97375823ec625319ac346"
variables: []
secrets_allowed: false
```
````sql
begin;
create table communication.whatsapp_status_batch (
 tenant_id uuid not null,
 delivery_key text not null,
 observations_sha256 text not null check(observations_sha256 ~ '^[0-9a-f]{64}$'),
 send_receipt_sha256 text not null check(send_receipt_sha256 ~ '^[0-9a-f]{64}$'),
 verification_receipt jsonb not null check(jsonb_typeof(verification_receipt)='object' and octet_length(verification_receipt::text)<=65536),
 observed_by text not null check(length(observed_by) between 1 and 255),
 recorded_at timestamptz not null default statement_timestamp(),
 primary key(tenant_id,delivery_key,observations_sha256),
 foreign key(tenant_id,delivery_key) references communication.whatsapp_appointment_approval(tenant_id,delivery_key)
);
create table communication.whatsapp_status_observation (
 tenant_id uuid not null,
 delivery_key text not null,
 event_key_sha256 text not null check(event_key_sha256 ~ '^[0-9a-f]{64}$'),
 event_sha256 text not null check(event_sha256 ~ '^[0-9a-f]{64}$'),
 provider_timestamp bigint not null check(provider_timestamp between 1 and 999999999999),
 provider_status text not null check(provider_status in ('sent','delivered','read','failed','deleted')),
 observations_sha256 text not null,
 primary key(tenant_id,delivery_key,event_key_sha256),
 foreign key(tenant_id,delivery_key,observations_sha256) references communication.whatsapp_status_batch(tenant_id,delivery_key,observations_sha256)
);
create index whatsapp_status_latest_idx on communication.whatsapp_status_observation(tenant_id,delivery_key,provider_timestamp desc);
create trigger whatsapp_status_batch_immutable before update or delete on communication.whatsapp_status_batch
 for each row execute function communication.reject_whatsapp_approval_mutation();
create trigger whatsapp_status_observation_immutable before update or delete on communication.whatsapp_status_observation
 for each row execute function communication.reject_whatsapp_approval_mutation();
commit;
````

### FILE: `db/migrations/0051_whatsapp_status_observations.down.sql`
```yaml
block_id: "PY-META-WHATSAPP:status-observations-down:v1"
operation: CREATE
provenance: AUTHORED
source: "Local verified integration of the admitted Meta signature boundary and PostgreSQL observation transactions; not upstream company code"
license: "LicenseRef-Workspace-Owner"
sha256: "fce14de4dae9007ccc2dfa0139251161fce502ebaf8a9b39eb41bc41253c2ef5"
variables: []
secrets_allowed: false
```
````sql
begin;
drop table communication.whatsapp_status_observation;
drop table communication.whatsapp_status_batch;
commit;
````

### FILE: `whatsapp_cloud/status_reconciliation.py`
```yaml
block_id: "PY-META-WHATSAPP:status-reconciliation:v1"
operation: CREATE
provenance: AUTHORED
source: "Local anchored correlation from the official Meta timestamp/identity contract, reusing admitted signature validation; not upstream code"
license: "LicenseRef-Workspace-Owner"
sha256: "ec31c2442ef176e89c4cd03de179345d1805d4f528b0710a9a4ed060ba7484ec"
variables: []
secrets_allowed: false
```
````python
"""AUTHORED correlation of an anchored send receipt with signed Meta statuses.

Reuses the existing Meta-adapted signature/scope boundary. Never sends messages
or changes the PostgreSQL fence. Hash anchors must come from trusted storage,
not from the webhook caller or a newly computed hash of an untrusted file.
"""
from __future__ import annotations

from datetime import datetime, timezone
import hmac
import json
from pathlib import Path
import re
import shutil
from typing import Any
import uuid
import base64
import tempfile

from whatsapp_cloud import SOURCE_COMMIT, _sha256, _unique_object, normalize_verified_webhook, validate_profile


def status_bridge(raw: bytes) -> dict[str, Any]:
    """Internal isolated-process contract. Secrets and raw bodies stay in memory."""
    if not 0 < len(raw) <= 12 * 1024 * 1024:
        raise ValueError("status frame exceeds local budget")
    frame = json.loads(raw.decode("utf-8"), object_pairs_hook=_unique_object)
    fields = {"schema", "profile", "send_receipt", "expected_send_receipt_sha256", "webhooks", "app_secret", "evidence_directory"}
    if not isinstance(frame, dict) or set(frame) != fields or frame["schema"] != "elite-whatsapp-status-bridge/v1":
        raise ValueError("status frame contract mismatch")
    if not isinstance(frame["app_secret"], str) or not 0 < len(frame["app_secret"]) <= 16384:
        raise ValueError("app secret unavailable")
    if not isinstance(frame["evidence_directory"], str) or not Path(frame["evidence_directory"]).is_absolute():
        raise ValueError("absolute evidence directory required")
    batches = frame["webhooks"]
    if not isinstance(batches, list) or not 1 <= len(batches) <= 8:
        raise ValueError("bounded webhook batch required")
    signed = []
    for item in batches:
        if not isinstance(item, dict) or set(item) != {"body", "signature"} or not isinstance(item["body"], str):
            raise ValueError("status batch contract mismatch")
        signed.append((base64.b64decode(item["body"], validate=True), item["signature"]))
    from whatsapp_cloud import verify_packaged_official_source
    root = Path(__file__).resolve().parent
    verify_packaged_official_source(root, root / "official-source.lock.json")
    with tempfile.TemporaryDirectory(prefix="whatsapp-status-", dir=frame["evidence_directory"]) as temp:
        output = Path(temp) / "verified"
        receipt = reconcile_status_webhooks(profile=frame["profile"],
            send_receipt_bytes=base64.b64decode(frame["send_receipt"], validate=True),
            expected_send_receipt_sha256=frame["expected_send_receipt_sha256"],
            signed_webhooks=signed, app_secret=frame["app_secret"], output_directory=output)
        observations = json.loads((output / "status-observations.json").read_bytes())
        return {"schema": "elite-whatsapp-status-result/v1", "binding_sha256": _sha256(raw),
                "receipt": receipt, "events": observations["events"]}


def reconcile_status_webhooks(*, profile: dict[str, Any], send_receipt_bytes: bytes,
                             expected_send_receipt_sha256: str,
                             signed_webhooks: list[tuple[bytes, str]], app_secret: str,
                             output_directory: Path) -> dict[str, Any]:
    """Publish a bounded observation snapshot, not a global delivery truth.

    Supply SEND_RECEIPT.json bytes and the EvidenceSHA256 already stored by the
    Go outbound fence. Each webhook is verified on its original bytes, including
    unrelated events in a shared batch. The snapshot covers only these inputs.
    """
    validate_profile(profile)
    if not isinstance(send_receipt_bytes, bytes) or not 0 < len(send_receipt_bytes) <= 65536:
        raise ValueError("send receipt exceeds local budget")
    if not isinstance(expected_send_receipt_sha256, str) or not re.fullmatch(r"[0-9a-f]{64}", expected_send_receipt_sha256):
        raise ValueError("trusted send receipt anchor is required")
    if not hmac.compare_digest(_sha256(send_receipt_bytes), expected_send_receipt_sha256):
        raise PermissionError("send receipt anchor mismatch")
    try:
        send = json.loads(send_receipt_bytes.decode("utf-8"), object_pairs_hook=_unique_object)
    except (UnicodeDecodeError, json.JSONDecodeError, RecursionError) as error:
        raise ValueError("send receipt must be UTF-8 JSON") from error
    if not isinstance(send, dict) or send.get("schema") != "elite-whatsapp-cloud-send-receipt/v1" or send.get("source_commit") != SOURCE_COMMIT or send.get("automatic_business_write") is not False:
        raise ValueError("send receipt contract mismatch")
    for key in ("phone_number_id_sha256", "recipient_sha256", "request_sha256", "response_sha256", "message_id_sha256"):
        if not isinstance(send.get(key), str) or not re.fullmatch(r"[0-9a-f]{64}", send[key]):
            raise ValueError("send receipt digest is invalid")
    if send.get("graph_api_version") != profile["graph_api_version"] or send["phone_number_id_sha256"] != _sha256(profile["phone_number_id"].encode("ascii")):
        raise PermissionError("send receipt profile mismatch")
    if not isinstance(signed_webhooks, list) or not 1 <= len(signed_webhooks) <= 8:
        raise ValueError("one to eight signed webhook bodies are required")

    matched: dict[str, dict[str, Any]] = {}
    raw_hashes: set[str] = set()
    excluded = 0
    duplicate = 0
    for item in signed_webhooks:
        if not isinstance(item, tuple) or len(item) != 2 or not isinstance(item[1], str):
            raise ValueError("webhook input must be raw bytes and signature")
        raw, signature = item
        scope, events = normalize_verified_webhook(raw, signature, app_secret, profile=profile)
        raw_hashes.add(_sha256(raw))
        for event in events:
            if event["kind"] != "status" or event["id_sha256"] != send["message_id_sha256"]:
                excluded += 1
                continue
            if event["recipient_sha256"] != send["recipient_sha256"]:
                raise PermissionError("matched provider message has another recipient")
            key = event["event_key_sha256"]
            if key in matched:
                if matched[key]["event_sha256"] != event["event_sha256"]:
                    raise ValueError("same status identity has divergent provider evidence")
                duplicate += 1
                continue
            matched[key] = event

    # Meta documents out-of-order delivery. Sort by provider time, never arrival.
    # Equal-time differing states remain explicit; do not invent precedence.
    timeline = sorted(matched.values(), key=lambda event: (int(event["timestamp"]), event["event_key_sha256"]))
    last_time = timeline[-1]["timestamp"] if timeline else None
    last_states = sorted({event["status"] for event in timeline if event["timestamp"] == last_time})
    decision = "NOT_OBSERVED_IN_INPUT" if not timeline else "AMBIGUOUS_LATEST_TIMESTAMP" if len(last_states) != 1 else "MATCHED_OBSERVATIONS"
    observations = {"schema": "elite-whatsapp-status-observations/v1", "events": timeline}
    normalized = (json.dumps(observations, sort_keys=True, indent=2) + "\n").encode("utf-8")
    receipt = {
        "schema": "elite-whatsapp-status-reconciliation/v1", "source_commit": SOURCE_COMMIT,
        "created_at": datetime.now(timezone.utc).isoformat(), **scope,
        "send_receipt_sha256": expected_send_receipt_sha256,
        "message_id_sha256": send["message_id_sha256"], "recipient_sha256": send["recipient_sha256"],
        "profile_sha256": _sha256(json.dumps(profile, sort_keys=True, separators=(",", ":")).encode("utf-8")),
        "webhook_body_sha256": sorted(raw_hashes), "observations_sha256": _sha256(normalized),
        "decision": decision, "last_observed_timestamp": last_time,
        "last_observed_status": last_states[0] if len(last_states) == 1 else None,
        "observed_statuses": sorted({event["status"] for event in timeline}),
        "matching_unique_events": len(timeline), "duplicate_events": duplicate,
        "excluded_events": excluded, "raw_payload_persisted": False,
        "automatic_business_write": False, "resend_authorized": False,
    }
    output = output_directory.resolve()
    if output.exists():
        raise FileExistsError("output_directory must not exist")
    output.parent.mkdir(parents=True, exist_ok=True)
    stage = output.parent / f".whatsapp-status-stage-{uuid.uuid4().hex}"
    stage.mkdir()
    try:
        (stage / "status-observations.json").write_bytes(normalized)
        (stage / "STATUS_RECONCILIATION.json").write_text(json.dumps(receipt, sort_keys=True, indent=2) + "\n", encoding="utf-8")
        stage.replace(output)
        return receipt
    except BaseException:
        shutil.rmtree(stage, ignore_errors=True)
        raise
````

### FILE: `whatsapp_cloud/test_status_reconciliation.py`
```yaml
block_id: "PY-META-WHATSAPP:test-status-reconciliation:v1"
operation: CREATE
provenance: AUTHORED
source: "Local tests connecting the existing send bridge receipt and synthetic signed status callbacks; no live Meta evidence"
license: "LicenseRef-Workspace-Owner"
sha256: "1513c03f4e07e63118af547bdcee365262ec237e86a5624ee906b2aea6c18d15"
variables: []
secrets_allowed: false
```
````python
from __future__ import annotations

import hashlib
import hmac
import itertools
import json
from pathlib import Path
import tempfile
import unittest
import base64
import os
import subprocess
import sys
from unittest.mock import patch

from status_reconciliation import reconcile_status_webhooks
from test_whatsapp_cloud import message, profile
from whatsapp_cloud import send_bridge


def event(status="delivered", timestamp="1603086314", **extra):
    return {"id": "wamid.synthetic", "recipient_id": "5491112345678", "timestamp": timestamp, "status": status, **extra}


def webhook(events, *, account="987654321", phone="123456789"):
    raw = json.dumps({"object": "whatsapp_business_account", "entry": [{"id": account, "changes": [{"field": "messages", "value": {"messaging_product": "whatsapp", "metadata": {"phone_number_id": phone}, "statuses": events}}]}]}, separators=(",", ":")).encode()
    return raw, "sha256=" + hmac.new(b"test-app-secret", raw, hashlib.sha256).hexdigest()


class StatusReconciliationTests(unittest.TestCase):
    def setUp(self):
        self.temp = tempfile.TemporaryDirectory()
        self.addCleanup(self.temp.cleanup)
        self.root = Path(self.temp.name)
        self.calls = 0
        def transport(*args):
            self.calls += 1
            return 200, {}, b'{"messages":[{"id":"wamid.synthetic"}]}'
        frame = {"schema": "elite-whatsapp-send-bridge/v1", "binding_sha256": "a" * 64, "profile": profile(), "request": message(), "access_token": "test-access-token", "output_directory": str(self.root / "send")}
        self.send = send_bridge(json.dumps(frame).encode(), transport)
        self.receipt = (self.root / "send" / "SEND_RECEIPT.json").read_bytes()

    def reconcile(self, batches=None, **overrides):
        args = dict(profile=profile(), send_receipt_bytes=self.receipt,
                    expected_send_receipt_sha256=self.send["evidence_sha256"],
                    signed_webhooks=batches if batches is not None else [webhook([event()])],
                    app_secret="test-app-secret", output_directory=self.root / "result")
        args.update(overrides)
        return reconcile_status_webhooks(**args)

    def test_send_bridge_to_signed_statuses_uses_durable_receipt_anchor_without_resend(self):
        with patch("whatsapp_cloud.urlopen", side_effect=AssertionError("no provider calls allowed")):
            result = self.reconcile()
        self.assertEqual(result["last_observed_status"], "delivered")
        self.assertEqual(result["send_receipt_sha256"], self.send["evidence_sha256"])
        self.assertEqual(self.calls, 1)
        self.assertFalse(result["resend_authorized"])
        self.assertFalse(result["automatic_business_write"])
        observed = (self.root / "result" / "status-observations.json").read_bytes()
        self.assertEqual(hashlib.sha256(observed).hexdigest(), result["observations_sha256"])
        for secret in ("test-app-secret", "test-access-token", "5491112345678", "wamid.synthetic"):
            self.assertNotIn(secret, observed.decode() + json.dumps(result))

    def test_all_arrival_permutations_have_identical_observation_hash(self):
        events = [event("sent", "1603086313"), event("delivered", "1603086314"), event("read", "1603086315")]
        hashes = set()
        for index, order in enumerate(itertools.permutations(events)):
            result = self.reconcile([webhook([item]) for item in order], output_directory=self.root / str(index))
            self.assertEqual(result["last_observed_status"], "read")
            hashes.add(result["observations_sha256"])
        self.assertEqual(len(hashes), 1)

    def test_repeated_batch_and_event_are_deduplicated(self):
        batch = webhook([event(), event()])
        result = self.reconcile([batch, batch])
        self.assertEqual(result["matching_unique_events"], 1)
        self.assertEqual(result["duplicate_events"], 3)

    def test_same_timestamp_does_not_invent_precedence(self):
        result = self.reconcile([webhook([event("sent"), event("delivered")])])
        self.assertEqual(result["decision"], "AMBIGUOUS_LATEST_TIMESTAMP")
        self.assertIsNone(result["last_observed_status"])
        self.assertEqual(result["matching_unique_events"], 2)

    def test_other_message_is_not_evidence_of_our_delivery(self):
        result = self.reconcile([webhook([event(id="wamid.other")])])
        self.assertEqual(result["decision"], "NOT_OBSERVED_IN_INPUT")
        self.assertIsNone(result["last_observed_status"])
        self.assertEqual(result["excluded_events"], 1)

    def test_mixed_batch_does_not_correlate_other_message(self):
        result = self.reconcile([webhook([event(), event("read", id="wamid.other")])])
        self.assertEqual(result["last_observed_status"], "delivered")
        self.assertEqual(result["excluded_events"], 1)

    def test_foreign_account_phone_or_recipient_rejects_atomically(self):
        for batch in [webhook([event()], account="999999999"), webhook([event()], phone="999999999"), webhook([event(recipient_id="5491112345679")])]:
            with self.subTest(batch_digest=hashlib.sha256(batch[0]).hexdigest()):
                with self.assertRaises(PermissionError): self.reconcile([batch])
                self.assertFalse((self.root / "result").exists())

    def test_signature_tamper_even_in_unrelated_message_rejects_entire_batch(self):
        raw, sig = webhook([event(id="wamid.other")])
        with self.assertRaises(PermissionError): self.reconcile([(raw + b" ", sig)])
        self.assertFalse((self.root / "result").exists())

    def test_send_anchor_tamper_rejects_before_output(self):
        with self.assertRaises(PermissionError): self.reconcile(send_receipt_bytes=self.receipt + b" ")
        self.assertFalse((self.root / "result").exists())

    def test_invalid_anchor_type_and_shape(self):
        for value in (None, "", "z" * 64):
            with self.subTest(value=value), self.assertRaises(ValueError):
                self.reconcile(expected_send_receipt_sha256=value)

    def test_same_event_identity_with_different_payload_is_not_silently_deduplicated(self):
        with self.assertRaisesRegex(ValueError, "divergent"):
            self.reconcile([webhook([event("failed", errors=[{"code": 1}])]), webhook([event("failed", errors=[{"code": 2}])])])
        self.assertFalse((self.root / "result").exists())

    def test_failed_and_deleted_are_retained_without_authorizing_retry(self):
        result = self.reconcile([webhook([event("failed", "1603086313"), event("deleted", "1603086314")])])
        self.assertEqual(result["observed_statuses"], ["deleted", "failed"])
        self.assertFalse(result["resend_authorized"])

    def test_profile_mismatch(self):
        changed = profile(); changed["graph_api_version"] = "v98.0"
        with self.assertRaises(PermissionError): self.reconcile(profile=changed)

    def test_empty_excessive_and_malformed_input(self):
        for value in ([], [webhook([event()])] * 9, [(b"{}",)], [(b"{}", None)]):
            with self.subTest(size=len(value)), self.assertRaises(ValueError): self.reconcile(value)
        self.assertFalse((self.root / "result").exists())

    def test_existing_output_is_never_overwritten(self):
        self.reconcile()
        previous = (self.root / "result" / "STATUS_RECONCILIATION.json").read_bytes()
        with self.assertRaises(FileExistsError): self.reconcile()
        self.assertEqual((self.root / "result" / "STATUS_RECONCILIATION.json").read_bytes(), previous)

    def test_write_failure_does_not_publish_partial_result(self):
        with patch.object(Path, "write_text", side_effect=OSError("fixture disk unavailable")):
            with self.assertRaises(OSError): self.reconcile()
        self.assertFalse((self.root / "result").exists())
        self.assertEqual(list(self.root.glob(".whatsapp-status-stage-*")), [])

    def test_late_bad_batch_leaves_no_partial_observations(self):
        raw, sig = webhook([event("read")])
        with self.assertRaises(PermissionError):
            self.reconcile([webhook([event()]), (raw + b" ", sig)])
        self.assertFalse((self.root / "result").exists())

    def test_bounded_raw_input_before_signature_parsing(self):
        with self.assertRaises(ValueError):
            self.reconcile([(b" " * (1024 * 1024 + 1), "sha256=" + "0" * 64)])
        with self.assertRaises(ValueError):
            self.reconcile(send_receipt_bytes=b" " * 65537)

    def test_receipt_contract_is_checked_even_when_hash_matches(self):
        for changes in ({"schema": "foreign"}, {"recipient_sha256": ""}, {"automatic_business_write": True}):
            altered = json.loads(self.receipt); altered.update(changes)
            raw = json.dumps(altered).encode()
            with self.subTest(changes=changes), self.assertRaises(ValueError):
                self.reconcile(send_receipt_bytes=raw, expected_send_receipt_sha256=hashlib.sha256(raw).hexdigest())

    def test_status_bridge_isolated_process_and_safe_failure(self):
        raw, signature = webhook([event()])
        frame = {"schema": "elite-whatsapp-status-bridge/v1", "profile": profile(),
                 "send_receipt": base64.b64encode(self.receipt).decode(),
                 "expected_send_receipt_sha256": self.send["evidence_sha256"],
                 "webhooks": [{"body": base64.b64encode(raw).decode(), "signature": signature}],
                 "app_secret": "test-app-secret", "evidence_directory": str(self.root)}
        environment = {key: os.environ[key] for key in ("SystemRoot",) if key in os.environ}
        args = [sys.executable, "-I", "-B", str(Path(__file__).resolve().parent / "whatsapp_cloud.py"), "--status-bridge"]
        encoded = json.dumps(frame).encode()
        result = subprocess.run(args, input=encoded, capture_output=True, env=environment, timeout=10)
        self.assertEqual(result.returncode, 0, result.stderr)
        reply = json.loads(result.stdout)
        self.assertEqual(reply["binding_sha256"], hashlib.sha256(encoded).hexdigest())
        self.assertEqual(reply["receipt"]["last_observed_status"], "delivered")
        self.assertEqual(list(self.root.glob("whatsapp-status-*")), [])
        frame["webhooks"][0]["signature"] = "sha256=" + "0" * 64
        result = subprocess.run(args, input=json.dumps(frame).encode(), capture_output=True, env=environment, timeout=10)
        self.assertEqual(result.returncode, 2)
        self.assertEqual(result.stdout, b"")
        self.assertEqual(result.stderr, b"WHATSAPP_STATUS_UNVERIFIED\n")


if __name__ == "__main__":
    unittest.main()
````

### FILE: `whatsapp_cloud/official-source.lock.json`
```yaml
block_id: "PY-META-WHATSAPP:source-lock:v1"
operation: CREATE
provenance: AUTHORED
source: "local lock derived from official GitHub metadata and raw bytes"
license: "LicenseRef-Workspace-Owner"
sha256: "044f5b73195a44afca1fb1be1ceb1a7ec9d5c15092dfe03bb491cb17b52453da"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-adapted-source-lock/v1",
  "source_id": "meta-whatsapp-api-examples",
  "repository": "fbsamples/whatsapp-api-examples",
  "commit": "de70ee908a67026e642aaee3703d20464e2a9466",
  "commit_signature_verified": true,
  "archive_bytes": 348229,
  "archive_sha256": "38d183a0d041d116dcf51356dd840c96690015f34c14c74b20681dd626ca9a7c",
  "license_expression": "LicenseRef-Meta-Platform-API-Only",
  "license_sha256": "ef5c10eeba318e71ebf81a7300fc7c97a115e230dcb4647dc597946976069300",
  "files": [
    {"packaged_path": "upstream/LICENSE", "upstream_path": "LICENSE", "upstream_sha256": "ef5c10eeba318e71ebf81a7300fc7c97a115e230dcb4647dc597946976069300", "sha256": "48d97b3c936203a750a3288c2b924327769d62327383a1060244b3d3c05ad01f", "provenance": "ADAPTED_FINAL_NEWLINE_ONLY"},
    {"packaged_path": "upstream/signature_validation_app.py", "upstream_path": "signature-validation-with-webhooks-payloads/app.py", "sha256": "6c052c136a0be059faee9e9ad93ff056b05b7086b819fce067ebd27c8ef2412c", "provenance": "VERBATIM_REFERENCE_ONLY"},
    {"packaged_path": "upstream/message_helper.py", "upstream_path": "send-messages-flight-app-python/message_helper.py", "sha256": "fef1aa035aa2fad4769a08edbe4a347a87dac475394a821bc71b301d1ef45702", "provenance": "VERBATIM_REFERENCE_ONLY"},
    {"packaged_path": "upstream/incomingWebhook.js", "upstream_path": "template-for-ecommerce-js/routes/incomingWebhook.js", "sha256": "a92a62bff637b8d5b0a3fbd3ce82dbd08c67904edff8df9adeb1be0e2d92330f", "provenance": "VERBATIM_REFERENCE_ONLY"}
  ],
  "verified_at": "2026-08-26"
}
````

### FILE: `whatsapp_cloud/provider-profile.template.json`
```yaml
block_id: "PY-META-WHATSAPP:profile:v1"
operation: CREATE
provenance: AUTHORED
source: "local fail-closed provider/terms/access profile"
license: "LicenseRef-Workspace-Owner"
sha256: "3029f068e1d67dcdc774b43c28f89705a7440247b7df02b8869c3813601bf841"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-whatsapp-cloud-profile/v1",
  "provider": "Meta WhatsApp Cloud API",
  "source_id": "meta-whatsapp-api-examples",
  "source_commit": "de70ee908a67026e642aaee3703d20464e2a9466",
  "decision": "BLOCKED_ACCESS_TERMS_AND_RECONCILIATION_REQUIRED",
  "official_source_license_accepted": false,
  "platform_terms_accepted": false,
  "business_account_proven": false,
  "app_registration_proven": false,
  "phone_number_id_proven": false,
  "message_template_approval_proven": false,
  "webhook_subscription_proven": false,
  "test_recipient_consent_proven": false,
  "quota_and_cost_approved": false,
  "reconciliation_approved": false,
  "graph_api_version": "",
  "business_account_id": "",
  "phone_number_id": "",
  "access_token_environment_variable": "WHATSAPP_ACCESS_TOKEN",
  "app_secret_environment_variable": "META_APP_SECRET",
  "verify_token_environment_variable": "WHATSAPP_VERIFY_TOKEN",
  "approved_templates": [],
  "data_retention": "",
  "automatic_business_write": false
}
````

### FILE: `whatsapp_cloud/template-message.template.json`
```yaml
block_id: "PY-META-WHATSAPP:message-template:v1"
operation: CREATE
provenance: AUTHORED
source: "local bounded message request template"
license: "LicenseRef-Workspace-Owner"
sha256: "7044c3885068429fa2915a68092d206701a181237cb37a198b2c1221d7667996"
variables: []
secrets_allowed: false
```
````json
{
  "recipient": "",
  "template_name": "",
  "language_code": "",
  "body_parameters": []
}
````

### FILE: `whatsapp_cloud/whatsapp_cloud.py`
```yaml
block_id: "PY-META-WHATSAPP:adapter:v1"
operation: CREATE
provenance: ADAPTED
source: "Meta examples signature-validation-with-webhooks-payloads/app.py, send-messages-flight-app-python/message_helper.py and template-for-ecommerce-js/routes/incomingWebhook.js at de70ee90; changes enumerated in PROVENANCE.md"
license: "LicenseRef-Meta-Platform-API-Only"
sha256: "04e3b9f80c2d471de3f105968e17283243bef350c9f64c0c9fc0e7ff844306c8"
variables: []
secrets_allowed: false
```
````python
"""
Copyright (c) Meta Platforms, Inc. and affiliates.
All rights reserved.

Adapted from the exact official examples listed in PROVENANCE.md under the
root upstream/LICENSE. Safety, atomic evidence and validation changes are
Copyright (c) the Elite Engineering Library owner.
"""

from __future__ import annotations

from datetime import datetime, timezone
import hashlib
import hmac
import json
import os
from pathlib import Path
import re
import shutil
import sys
from typing import Any, Callable
from urllib.error import HTTPError
from urllib.request import Request, urlopen
import uuid


SOURCE_ID = "meta-whatsapp-api-examples"
SOURCE_COMMIT = "de70ee908a67026e642aaee3703d20464e2a9466"
Transport = Callable[[str, str, dict[str, str], bytes, float], tuple[int, dict[str, str], bytes]]


def _sha256(value: bytes) -> str:
    return hashlib.sha256(value).hexdigest()


def _object(path: Path) -> dict[str, Any]:
    value = json.loads(path.read_text(encoding="utf-8"))
    if not isinstance(value, dict):
        raise ValueError(f"{path.name} must contain one JSON object")
    return value


def verify_packaged_official_source(root: Path, lock_path: Path) -> None:
    lock = _object(lock_path)
    if lock.get("source_id") != SOURCE_ID or lock.get("commit") != SOURCE_COMMIT or lock.get("commit_signature_verified") is not True:
        raise RuntimeError("official Meta source identity/signature mismatch")
    resolved = root.resolve()
    for entry in lock.get("files", []):
        relative = str(entry.get("packaged_path", ""))
        if not relative or Path(relative).is_absolute() or ".." in Path(relative).parts:
            raise RuntimeError("unsafe official source path")
        target = (resolved / relative).resolve()
        if resolved not in target.parents or not target.is_file() or _sha256(target.read_bytes()) != entry.get("sha256"):
            raise RuntimeError(f"packaged official source mismatch: {relative}")


def validate_profile(profile: dict[str, Any]) -> None:
    required = (
        "official_source_license_accepted", "platform_terms_accepted", "business_account_proven",
        "app_registration_proven", "phone_number_id_proven", "message_template_approval_proven",
        "webhook_subscription_proven", "test_recipient_consent_proven", "quota_and_cost_approved",
        "reconciliation_approved",
    )
    if profile.get("provider") != "Meta WhatsApp Cloud API" or profile.get("source_id") != SOURCE_ID or profile.get("source_commit") != SOURCE_COMMIT:
        raise ValueError("WhatsApp profile authority mismatch")
    if profile.get("decision") != "PROVEN" or any(profile.get(key) is not True for key in required):
        raise PermissionError("WhatsApp access, terms or reconciliation remain blocked")
    if profile.get("automatic_business_write") is not False or not str(profile.get("data_retention", "")).strip():
        raise ValueError("retention is required and automatic_business_write must remain false")
    if not re.fullmatch(r"v[0-9]{1,3}\.[0-9]{1,2}", str(profile.get("graph_api_version", ""))):
        raise ValueError("graph_api_version must be explicitly approved, for example vNN.0")
    for key in ("business_account_id", "phone_number_id"):
        if not isinstance(profile.get(key), str) or not re.fullmatch(r"[0-9]{5,30}", profile[key]):
            raise ValueError(f"{key} must be a configured numeric string")
    for key in ("access_token_environment_variable", "app_secret_environment_variable", "verify_token_environment_variable"):
        if not re.fullmatch(r"[A-Z][A-Z0-9_]{2,80}", str(profile.get(key, ""))):
            raise ValueError(f"{key} must name an environment variable")
    approved = profile.get("approved_templates")
    if not isinstance(approved, list) or not approved:
        raise PermissionError("at least one exact approved template is required")
    seen: set[tuple[str, str]] = set()
    for item in approved:
        if not isinstance(item, dict):
            raise ValueError("approved template entries must be objects")
        name, language, count = item.get("name"), item.get("language_code"), item.get("body_parameter_count")
        if not isinstance(name, str) or not re.fullmatch(r"[a-z0-9_]{1,512}", name):
            raise ValueError("approved template name is invalid")
        if not isinstance(language, str) or not re.fullmatch(r"[a-z]{2}(?:_[A-Z]{2})?", language):
            raise ValueError("approved template language is invalid")
        if not isinstance(count, int) or isinstance(count, bool) or not 0 <= count <= 20 or (name, language) in seen:
            raise ValueError("approved template parameter count/identity is invalid")
        seen.add((name, language))


def build_template_payload(profile: dict[str, Any], request: dict[str, Any]) -> tuple[str, dict[str, Any]]:
    validate_profile(profile)
    recipient = str(request.get("recipient", ""))
    if not re.fullmatch(r"[1-9][0-9]{7,14}", recipient):
        raise ValueError("recipient must be an E.164 number without plus sign")
    name, language, parameters = request.get("template_name"), request.get("language_code"), request.get("body_parameters")
    if not isinstance(parameters, list) or len(parameters) > 20 or any(not isinstance(value, str) or not value or len(value) > 1024 for value in parameters):
        raise ValueError("body_parameters must be bounded non-empty strings")
    approved = {(item["name"], item["language_code"]): item["body_parameter_count"] for item in profile["approved_templates"]}
    if approved.get((name, language)) != len(parameters):
        raise PermissionError("template, language or parameter count is not approved")
    components = []
    if parameters:
        components.append({"type": "body", "parameters": [{"type": "text", "text": value} for value in parameters]})
    payload = {"messaging_product": "whatsapp", "recipient_type": "individual", "to": recipient, "type": "template", "template": {"name": name, "language": {"code": language}, "components": components}}
    return recipient, payload


def build_reply_payload(profile: dict[str, Any], request: dict[str, Any], *, observed_at: float | None = None) -> tuple[str, dict[str, Any]]:
    """ADAPTED get_text_message_input at the already-pinned Meta commit.

    Scope/window guards are local composition of a human-approved proposal;
    a text message is never interpreted as an approved Meta template.
    """
    validate_profile(profile)
    if set(request) != {"kind", "recipient", "text", "source_message_id", "last_inbound_at", "window_expires_at"} or request.get("kind") != "text_reply":
        raise ValueError("text reply contract mismatch")
    recipient, text = request.get("recipient"), request.get("text")
    if not isinstance(recipient, str) or not re.fullmatch(r"[1-9][0-9]{7,14}", recipient):
        raise ValueError("reply recipient must be exact E.164 digits")
    # Conservative local budget, not a statement of the provider maximum.
    if not isinstance(text, str) or not 0 < len(text) <= 1024 or len(text.encode("utf-8", "strict")) > 4096 or any(ord(c)<32 and c not in "\n\t" for c in text):
        raise ValueError("reply text exceeds local contract")
    _event_identity(request.get("source_message_id"), "source message")
    start, end = request.get("last_inbound_at"), request.get("window_expires_at")
    if type(start) is not int or type(end) is not int or start<1 or end!=start+86400:
        raise ValueError("reply service window is invalid")
    now=datetime.now(timezone.utc).timestamp() if observed_at is None else observed_at
    if now<start or now>=end:
        raise PermissionError("reply service window is closed")
    return recipient, {"messaging_product":"whatsapp", "recipient_type":"individual", "to":recipient, "type":"text", "text":{"body":text}}


def stdlib_transport(method: str, url: str, headers: dict[str, str], body: bytes, timeout: float) -> tuple[int, dict[str, str], bytes]:
    request = Request(url, data=body, headers=headers, method=method)
    try:
        with urlopen(request, timeout=timeout) as response:
            return response.status, dict(response.headers.items()), response.read()
    except HTTPError as error:
        return error.code, dict(error.headers.items()), error.read()


def send_template_to_evidence(profile: dict[str, Any], request: dict[str, Any], access_token: str, output_directory: Path, transport: Transport = stdlib_transport) -> dict[str, Any]:
    if not access_token:
        raise PermissionError("WhatsApp access token is unavailable")
    recipient, payload = (build_reply_payload(profile, request) if request.get("kind") == "text_reply" else build_template_payload(profile, request))
    output = output_directory.resolve()
    if output.exists():
        raise FileExistsError("output_directory must not exist")
    output.parent.mkdir(parents=True, exist_ok=True)
    stage = output.parent / f".whatsapp-send-stage-{uuid.uuid4().hex}"
    stage.mkdir()
    try:
        request_bytes = (json.dumps(payload, ensure_ascii=False, sort_keys=True, separators=(",", ":")) + "\n").encode("utf-8")
        version, phone_id = profile["graph_api_version"], profile["phone_number_id"]
        url = f"https://graph.facebook.com/{version}/{phone_id}/messages"
        status, _, response_bytes = transport("POST", url, {"Authorization": f"Bearer {access_token}", "Content-Type": "application/json"}, request_bytes, 30.0)
        try:
            response = json.loads(response_bytes.decode("utf-8"))
        except (UnicodeDecodeError, json.JSONDecodeError) as error:
            raise RuntimeError("WhatsApp response is not UTF-8 JSON") from error
        if status < 200 or status >= 300 or not isinstance(response, dict):
            raise RuntimeError(f"WhatsApp provider rejected message with HTTP {status}")
        messages = response.get("messages")
        message_id = messages[0].get("id") if isinstance(messages, list) and messages and isinstance(messages[0], dict) else None
        if not isinstance(message_id, str) or not message_id:
            raise RuntimeError("WhatsApp success response is missing message id")
        (stage / "provider-response.json").write_bytes(response_bytes)
        receipt = {"schema": "elite-whatsapp-cloud-send-receipt/v1", "created_at": datetime.now(timezone.utc).isoformat(), "source_commit": SOURCE_COMMIT, "graph_api_version": version, "phone_number_id_sha256": _sha256(phone_id.encode("ascii")), "recipient_sha256": _sha256(recipient.encode("ascii")), "request_sha256": _sha256(request_bytes), "response_sha256": _sha256(response_bytes), "message_id_sha256": _sha256(message_id.encode("utf-8")), "automatic_business_write": False}
        (stage / "SEND_RECEIPT.json").write_text(json.dumps(receipt, sort_keys=True, indent=2) + "\n", encoding="utf-8")
        stage.replace(output)
        return receipt
    except BaseException:
        shutil.rmtree(stage, ignore_errors=True)
        raise


def verify_webhook_signature(raw_body: bytes, signature_header: str, app_secret: str) -> None:
    if not app_secret:
        raise PermissionError("Meta app secret is unavailable")
    if not re.fullmatch(r"sha256=[0-9a-f]{64}", signature_header or ""):
        raise PermissionError("WhatsApp webhook signature header is missing or malformed")
    expected = "sha256=" + hmac.new(app_secret.encode("utf-8"), raw_body, hashlib.sha256).hexdigest()
    if not hmac.compare_digest(signature_header, expected):
        raise PermissionError("WhatsApp webhook signature mismatch")


def verify_subscription(query: dict[str, str], verify_token: str) -> str:
    challenge = query.get("hub.challenge", "")
    if not verify_token or query.get("hub.mode") != "subscribe" or not challenge or not hmac.compare_digest(query.get("hub.verify_token", ""), verify_token):
        raise PermissionError("WhatsApp webhook subscription verification failed")
    return challenge


def _unique_object(pairs: list[tuple[str, Any]]) -> dict[str, Any]:
    result: dict[str, Any] = {}
    for key, value in pairs:
        if key in result:
            raise ValueError("WhatsApp webhook contains duplicate JSON keys")
        result[key] = value
    return result


def _event_array(value: Any) -> list[Any]:
    if not isinstance(value, list) or len(value) > 1000:
        raise ValueError("WhatsApp webhook arrays must be bounded lists")
    return value


def _event_identity(value: Any, label: str) -> str:
    if not isinstance(value, str) or not value or len(value) > 512 or any(ord(c) < 33 or ord(c) > 126 for c in value):
        raise ValueError(f"WhatsApp {label} must be a bounded non-empty identity")
    return value


def normalize_verified_webhook(raw_body: bytes, signature_header: str, app_secret: str, *, profile: dict[str, Any]) -> tuple[dict[str, str], list[dict[str, Any]]]:
    """The shared signature/scope/parser boundary; no disk or business effects."""
    validate_profile(profile)
    # Local resource budget, not a claimed Meta API limit. Reject before parsing.
    if not isinstance(raw_body, bytes) or not 0 < len(raw_body) <= 1024 * 1024:
        raise ValueError("WhatsApp webhook exceeds local body budget")
    verify_webhook_signature(raw_body, signature_header, app_secret)
    try:
        body = json.loads(raw_body.decode("utf-8"), object_pairs_hook=_unique_object)
    except (UnicodeDecodeError, json.JSONDecodeError, RecursionError) as error:
        raise ValueError("WhatsApp webhook must be UTF-8 JSON") from error
    if not isinstance(body, dict) or body.get("object") != "whatsapp_business_account" or not isinstance(body.get("entry"), list):
        raise ValueError("WhatsApp webhook envelope is invalid")
    events: list[dict[str, Any]] = []
    scope = {"business_account_id_sha256": _sha256(profile["business_account_id"].encode("ascii")), "phone_number_id_sha256": _sha256(profile["phone_number_id"].encode("ascii"))}
    for entry in _event_array(body["entry"]):
        if not isinstance(entry, dict):
            raise ValueError("WhatsApp entry is invalid")
        if entry.get("id") != profile["business_account_id"]:
            raise PermissionError("WhatsApp business account scope mismatch")
        for change in _event_array(entry.get("changes")):
            if not isinstance(change, dict) or change.get("field") != "messages" or not isinstance(change.get("value"), dict):
                raise ValueError("WhatsApp change is not an admitted messages envelope")
            value = change["value"]
            metadata = value.get("metadata")
            if value.get("messaging_product") != "whatsapp" or not isinstance(metadata, dict):
                raise ValueError("WhatsApp product/metadata is missing")
            if metadata.get("phone_number_id") != profile["phone_number_id"]:
                raise PermissionError("WhatsApp phone number scope mismatch")
            for field, kind, contact_key in (("messages", "message", "from"), ("statuses", "status", "recipient_id")):
                for event in _event_array(value.get(field, [])):
                    if not isinstance(event, dict):
                        raise ValueError("WhatsApp event must be an object")
                    identity = _event_identity(event.get("id"), "message id")
                    contact = _event_identity(event.get(contact_key), "contact id")
                    timestamp = event.get("timestamp")
                    if not isinstance(timestamp, str) or not re.fullmatch(r"[1-9][0-9]{0,11}", timestamp):
                        raise ValueError("WhatsApp timestamp must be canonical positive Unix seconds")
                    normalized_event = {**scope, "kind": kind, "id_sha256": _sha256(identity.encode("ascii")), "timestamp": timestamp, "event_sha256": _sha256(json.dumps(event, sort_keys=True, separators=(",", ":")).encode("utf-8"))}
                    normalized_event["sender_sha256" if kind == "message" else "recipient_sha256"] = _sha256(contact.encode("ascii"))
                    if kind == "status":
                        status = event.get("status")
                        if not isinstance(status, str) or status not in {"sent", "delivered", "read", "failed", "deleted"}:
                            raise ValueError("WhatsApp status is not admitted")
                        normalized_event["status"] = status
                    else:
                        message_type = event.get("type")
                        if not isinstance(message_type, str) or not re.fullmatch(r"[a-z][a-z0-9_]{0,63}", message_type):
                            raise ValueError("WhatsApp message type is invalid")
                        if message_type == "text":
                            text = event.get("text")
                            if not isinstance(text, dict) or not isinstance(text.get("body"), str) or not 0 < len(text["body"].encode("utf-8", "strict")) <= 16384:
                                raise ValueError("text message body is invalid")
                        normalized_event["type"] = message_type
                    # Stable provider-event identity independent of batch order and JSON layout.
                    identity_fields = {key: val for key, val in normalized_event.items() if key != "event_sha256"}
                    normalized_event["event_key_sha256"] = _sha256(json.dumps(identity_fields, sort_keys=True, separators=(",", ":")).encode("ascii"))
                    events.append(normalized_event)
                    if len(events) > 1000:
                        raise ValueError("WhatsApp webhook exceeds local event budget")
    if not events:
        raise ValueError("WhatsApp webhook contains no admitted message/status events")
    return scope, events


def webhook_to_evidence(raw_body: bytes, signature_header: str, app_secret: str, output_directory: Path, *, profile: dict[str, Any]) -> dict[str, Any]:
    scope, events = normalize_verified_webhook(raw_body, signature_header, app_secret, profile=profile)
    output = output_directory.resolve()
    if output.exists():
        raise FileExistsError("output_directory must not exist")
    output.parent.mkdir(parents=True, exist_ok=True)
    stage = output.parent / f".whatsapp-webhook-stage-{uuid.uuid4().hex}"
    stage.mkdir()
    try:
        normalized = (json.dumps({"schema": "elite-whatsapp-cloud-events/v2", "events": events}, sort_keys=True, indent=2) + "\n").encode("utf-8")
        (stage / "normalized-events.json").write_bytes(normalized)
        receipt = {"schema": "elite-whatsapp-cloud-webhook-receipt/v2", **scope, "profile_sha256": _sha256(json.dumps(profile, sort_keys=True, separators=(",", ":")).encode("utf-8")), "created_at": datetime.now(timezone.utc).isoformat(), "source_commit": SOURCE_COMMIT, "raw_body_sha256": _sha256(raw_body), "normalized_sha256": _sha256(normalized), "event_count": len(events), "raw_payload_persisted": False, "automatic_business_write": False}
        (stage / "WEBHOOK_RECEIPT.json").write_text(json.dumps(receipt, sort_keys=True, indent=2) + "\n", encoding="utf-8")
        stage.replace(output)
        return receipt
    except BaseException:
        shutil.rmtree(stage, ignore_errors=True)
        raise


def send_from_files(profile_path: Path, request_path: Path, output: Path) -> dict[str, Any]:
    profile, message = _object(profile_path), _object(request_path)
    secret_name = str(profile.get("access_token_environment_variable", ""))
    return send_template_to_evidence(profile, message, os.environ.get(secret_name, ""), output)


def send_bridge(raw: bytes, transport: Transport = stdlib_transport) -> dict[str, Any]:
    """One approved send, invoked only behind the existing durable Go fence.

    The frame contains secrets in memory; neither the frame nor exceptions may
    be printed/persisted. This bridge adds no retry or delivery-state owner.
    """
    if not 0 < len(raw) <= 131072:
        raise ValueError("bridge frame exceeds local budget")
    frame = json.loads(raw.decode("utf-8"), object_pairs_hook=_unique_object)
    fields = {"schema", "binding_sha256", "profile", "request", "access_token", "output_directory"}
    if not isinstance(frame, dict) or set(frame) != fields or frame["schema"] != "elite-whatsapp-send-bridge/v1":
        raise ValueError("bridge frame contract mismatch")
    binding = frame["binding_sha256"]
    if not isinstance(binding, str) or not re.fullmatch(r"[0-9a-f]{64}", binding):
        raise ValueError("bridge binding is invalid")
    if not isinstance(frame["profile"], dict) or not isinstance(frame["request"], dict):
        raise ValueError("bridge profile/request must be objects")
    token, destination = frame["access_token"], frame["output_directory"]
    if not isinstance(token, str) or not 0 < len(token) <= 16384 or not isinstance(destination, str) or not Path(destination).is_absolute():
        raise ValueError("bridge secret or output unavailable")
    root = Path(__file__).resolve().parent
    verify_packaged_official_source(root, root / "official-source.lock.json")
    receipt = send_template_to_evidence(frame["profile"], frame["request"], token, Path(destination), transport)
    evidence = Path(destination)
    provider_bytes = (evidence / "provider-response.json").read_bytes()
    if _sha256(provider_bytes) != receipt["response_sha256"]:
        raise ValueError("provider evidence changed")
    response = json.loads(provider_bytes.decode("utf-8"), object_pairs_hook=_unique_object)
    messages = response.get("messages")
    if not isinstance(messages, list) or len(messages) != 1 or not isinstance(messages[0], dict):
        raise ValueError("provider message identity is ambiguous")
    identity = _event_identity(messages[0].get("id"), "provider message id")
    if len(identity) > 256 or _sha256(identity.encode("ascii")) != receipt["message_id_sha256"]:
        raise ValueError("provider message identity mismatch")
    return {"schema": "elite-whatsapp-send-result/v1", "binding_sha256": binding,
            "provider_message_id": identity, "evidence_sha256": _sha256((evidence / "SEND_RECEIPT.json").read_bytes()),
            "accepted_at": receipt["created_at"]}


def ingress_bridge(raw: bytes) -> dict[str, Any]:
    """AUTHORED stdio boundary reusing the existing Meta verifier; no I/O effects."""
    import base64
    if not 0 < len(raw) <= 2 * 1024 * 1024:
        raise ValueError("ingress frame exceeds local budget")
    frame = json.loads(raw.decode("utf-8"), object_pairs_hook=_unique_object)
    if not isinstance(frame, dict) or set(frame) != {"schema", "mode", "profile", "body", "signature", "secret", "query"} or frame["schema"] != "elite-whatsapp-ingress-bridge/v1":
        raise ValueError("ingress frame contract mismatch")
    if not isinstance(frame["secret"], str) or not 0 < len(frame["secret"]) <= 16384:
        raise PermissionError("ingress secret unavailable")
    root = Path(__file__).resolve().parent
    verify_packaged_official_source(root, root / "official-source.lock.json")
    reply = {"schema": "elite-whatsapp-ingress-result/v1", "binding_sha256": _sha256(raw), "body_sha256": "", "event_count": 0, "challenge": ""}
    if frame["mode"] == "subscribe":
        query = frame["query"]
        if not isinstance(query, dict) or set(query) != {"hub.mode", "hub.verify_token", "hub.challenge"} or any(not isinstance(v, str) for v in query.values()):
            raise ValueError("subscription query mismatch")
        if not re.fullmatch(r"[0-9]{1,256}", query["hub.challenge"]):
            raise ValueError("challenge exceeds local contract")
        if frame["body"] or frame["signature"]:
            raise ValueError("subscription cannot contain body or signature")
        reply["challenge"] = verify_subscription(query, frame["secret"])
    elif frame["mode"] == "receive":
        if frame["query"] or not isinstance(frame["body"], str):
            raise ValueError("ingress body contract mismatch")
        body = base64.b64decode(frame["body"], validate=True)
        _, events = normalize_verified_webhook(body, frame["signature"], frame["secret"], profile=frame["profile"])
        reply.update(body_sha256=_sha256(body), event_count=len(events))
    else:
        raise ValueError("ingress mode rejected")
    return reply


def recover_send_bridge(raw: bytes) -> dict[str, Any]:
    """Pure evidence validation; never calls transport or changes delivery state."""
    import base64
    if not 0<len(raw)<=262144:
        raise ValueError("recovery frame exceeds local budget")
    frame=json.loads(raw.decode("utf-8"),object_pairs_hook=_unique_object)
    if not isinstance(frame,dict) or set(frame)!={"schema","binding_sha256","profile","request","receipt","response"} or frame["schema"]!="elite-whatsapp-recover-send/v1":
        raise ValueError("recovery frame contract mismatch")
    if not re.fullmatch(r"[0-9a-f]{64}",str(frame["binding_sha256"])):
        raise ValueError("recovery binding invalid")
    root=Path(__file__).resolve().parent
    verify_packaged_official_source(root,root/"official-source.lock.json")
    receipt_raw=base64.b64decode(frame["receipt"],validate=True)
    response_raw=base64.b64decode(frame["response"],validate=True)
    if not 0<len(receipt_raw)<=65536 or not 0<len(response_raw)<=65536:
        raise ValueError("recovery evidence exceeds local budget")
    receipt=json.loads(receipt_raw.decode("utf-8"),object_pairs_hook=_unique_object)
    response=json.loads(response_raw.decode("utf-8"),object_pairs_hook=_unique_object)
    at=datetime.fromisoformat(receipt["created_at"])
    if at.tzinfo is None or at.timestamp()>datetime.now(timezone.utc).timestamp()+60:
        raise ValueError("recovery acceptance time invalid")
    request,profile=frame["request"],frame["profile"]
    if request.get("kind")!="text_reply":
        raise ValueError("recovery supports exact conversational reply only")
    recipient,payload=build_reply_payload(profile,request,observed_at=at.timestamp())
    expected=(json.dumps(payload,ensure_ascii=False,sort_keys=True,separators=(",",":"))+"\n").encode("utf-8")
    if receipt.get("schema")!="elite-whatsapp-cloud-send-receipt/v1" or receipt.get("source_commit")!=SOURCE_COMMIT or receipt.get("automatic_business_write") is not False or receipt.get("graph_api_version")!=profile["graph_api_version"] or receipt.get("phone_number_id_sha256")!=_sha256(profile["phone_number_id"].encode("ascii")) or receipt.get("recipient_sha256")!=_sha256(recipient.encode("ascii")) or receipt.get("request_sha256")!=_sha256(expected) or receipt.get("response_sha256")!=_sha256(response_raw):
        raise ValueError("recovery evidence binding mismatch")
    messages=response.get("messages")
    if not isinstance(messages,list) or len(messages)!=1 or not isinstance(messages[0],dict):
        raise ValueError("recovery provider identity ambiguous")
    identity=_event_identity(messages[0].get("id"),"provider message id")
    if len(identity)>256 or receipt.get("message_id_sha256")!=_sha256(identity.encode("ascii")):
        raise ValueError("recovery provider identity mismatch")
    return {"schema":"elite-whatsapp-send-result/v1","binding_sha256":frame["binding_sha256"],"provider_message_id":identity,"evidence_sha256":_sha256(receipt_raw),"accepted_at":receipt["created_at"]}


def validate_profile_bridge(raw: bytes) -> dict[str, str]:
    """Bounded startup validation of the exact configuration bytes; no I/O."""
    import base64
    if not 0 < len(raw) <= 65536:
        raise ValueError("profile validation frame exceeds local budget")
    frame=json.loads(raw.decode("utf-8"),object_pairs_hook=_unique_object)
    if not isinstance(frame,dict) or set(frame)!={"schema","profile"} or frame["schema"]!="elite-whatsapp-validate-profile/v1" or not isinstance(frame["profile"],str):
        raise ValueError("profile validation contract mismatch")
    profile_raw=base64.b64decode(frame["profile"],validate=True)
    if not 0<len(profile_raw)<=32768:
        raise ValueError("profile exceeds local budget")
    profile=json.loads(profile_raw.decode("utf-8"),object_pairs_hook=_unique_object)
    if not isinstance(profile,dict):
        raise ValueError("profile must be an object")
    validate_profile(profile)
    return {"schema":"elite-whatsapp-profile-validation/v1","profile_sha256":_sha256(profile_raw)}


if __name__ == "__main__":
    if sys.argv[1:] not in (["--send-bridge"], ["--status-bridge"], ["--ingress-bridge"], ["--recover-send-bridge"], ["--validate-profile-bridge"]):
        sys.exit(2)
    try:
        if sys.argv[1:] == ["--send-bridge"]:
            reply = send_bridge(sys.stdin.buffer.read(131073))
        elif sys.argv[1:] == ["--validate-profile-bridge"]:
            reply = validate_profile_bridge(sys.stdin.buffer.read(65537))
        elif sys.argv[1:] == ["--recover-send-bridge"]:
            reply = recover_send_bridge(sys.stdin.buffer.read(262145))
        elif sys.argv[1:] == ["--ingress-bridge"]:
            reply = ingress_bridge(sys.stdin.buffer.read(2 * 1024 * 1024 + 1))
        else:
            # Fixed local module; Go verifies both source hashes before -I/-B.
            import importlib.util
            sys.modules["whatsapp_cloud"] = sys.modules[__name__]
            spec = importlib.util.spec_from_file_location("status_reconciliation", Path(__file__).resolve().parent / "status_reconciliation.py")
            module = importlib.util.module_from_spec(spec)
            spec.loader.exec_module(module)
            reply = module.status_bridge(sys.stdin.buffer.read(12 * 1024 * 1024 + 1))
        sys.stdout.buffer.write((json.dumps(reply, sort_keys=True, separators=(",", ":")) + "\n").encode("utf-8"))
    except Exception:
        # No provider body, token, recipient or exception text crosses stderr.
        sys.stderr.buffer.write(b"WHATSAPP_INGRESS_UNVERIFIED\n" if sys.argv[1:] == ["--ingress-bridge"] else b"WHATSAPP_STATUS_UNVERIFIED\n" if sys.argv[1:] == ["--status-bridge"] else b"WHATSAPP_SEND_UNCERTAIN_OR_REJECTED\n")
        sys.exit(2)
````

### FILE: `whatsapp_cloud/test_whatsapp_cloud.py`
```yaml
block_id: "PY-META-WHATSAPP:tests:v1"
operation: CREATE
provenance: AUTHORED
source: "local source-integrity, failure, signature, redaction and atomicity tests"
license: "LicenseRef-Workspace-Owner"
sha256: "c45a97b3584b6f308316d9760610fca4a5a4354c3d9c2146af5180402f42f421"
variables: []
secrets_allowed: false
```
````python
from __future__ import annotations

import hashlib
import hmac
import json
from pathlib import Path
import tempfile
import unittest
import subprocess
import sys

from whatsapp_cloud import build_template_payload, send_template_to_evidence, verify_packaged_official_source, verify_subscription, verify_webhook_signature, webhook_to_evidence, send_bridge


def profile():
    return {"provider": "Meta WhatsApp Cloud API", "source_id": "meta-whatsapp-api-examples", "source_commit": "de70ee908a67026e642aaee3703d20464e2a9466", "decision": "PROVEN", "official_source_license_accepted": True, "platform_terms_accepted": True, "business_account_proven": True, "app_registration_proven": True, "phone_number_id_proven": True, "message_template_approval_proven": True, "webhook_subscription_proven": True, "test_recipient_consent_proven": True, "quota_and_cost_approved": True, "reconciliation_approved": True, "business_account_id": "987654321", "graph_api_version": "v99.0", "phone_number_id": "123456789", "access_token_environment_variable": "WHATSAPP_ACCESS_TOKEN", "app_secret_environment_variable": "META_APP_SECRET", "verify_token_environment_variable": "WHATSAPP_VERIFY_TOKEN", "approved_templates": [{"name": "order_update", "language_code": "es_AR", "body_parameter_count": 2}], "data_retention": "30 days", "automatic_business_write": False}


def message():
    return {"recipient": "5491112345678", "template_name": "order_update", "language_code": "es_AR", "body_parameters": ["A-1", "despachado"]}


class WhatsAppCloudTests(unittest.TestCase):
    def test_signed_foreign_scope_is_rejected(self):
        payload = {"object": "whatsapp_business_account", "entry": [{"id": "999999999", "changes": [{"field": "messages", "value": {"messaging_product": "whatsapp", "metadata": {"phone_number_id": "999999999"}, "statuses": [{"id": "wamid.1", "status": "delivered", "timestamp": "1603086313", "recipient_id": "5491112345678"}]}}]}]}
        raw = json.dumps(payload).encode()
        signature = "sha256=" + hmac.new(b"app-secret", raw, hashlib.sha256).hexdigest()
        with tempfile.TemporaryDirectory() as temp:
            with self.assertRaises(PermissionError):
                webhook_to_evidence(raw, signature, "app-secret", Path(temp) / "evidence", profile=profile())

    def test_packaged_official_source_hashes(self):
        root = Path(__file__).parent
        verify_packaged_official_source(root, root / "official-source.lock.json")

    def test_profile_and_template_fail_closed(self):
        blocked = profile(); blocked["decision"] = "BLOCKED"
        with self.assertRaises(PermissionError): build_template_payload(blocked, message())
        bad = message(); bad["body_parameters"] = ["only-one"]
        with self.assertRaises(PermissionError): build_template_payload(profile(), bad)

    def test_send_preserves_response_and_redacts_identifiers(self):
        calls = []
        def transport(method, url, headers, body, timeout):
            calls.append((method, url, headers, body, timeout))
            return 200, {"content-type": "application/json"}, b'{"messages":[{"id":"wamid.secret-id"}]}'
        with tempfile.TemporaryDirectory() as temp:
            output = Path(temp) / "evidence"
            receipt = send_template_to_evidence(profile(), message(), "secret-token", output, transport)
            serialized = json.dumps(receipt)
            self.assertNotIn("5491112345678", serialized); self.assertNotIn("123456789", serialized); self.assertNotIn("secret", serialized)
            self.assertEqual(calls[0][0], "POST"); self.assertEqual(calls[0][1], "https://graph.facebook.com/v99.0/123456789/messages")
            self.assertEqual(calls[0][2]["Authorization"], "Bearer secret-token")
            self.assertTrue((output / "provider-response.json").is_file())

    def test_send_provider_failure_is_atomic(self):
        def transport(*_): return 403, {}, b'{"error":{"message":"denied"}}'
        with tempfile.TemporaryDirectory() as temp:
            output = Path(temp) / "evidence"
            with self.assertRaisesRegex(RuntimeError, "HTTP 403"):
                send_template_to_evidence(profile(), message(), "token", output, transport)
            self.assertFalse(output.exists()); self.assertEqual(list(Path(temp).iterdir()), [])

    def test_subscription_is_exact(self):
        self.assertEqual(verify_subscription({"hub.mode": "subscribe", "hub.verify_token": "verify", "hub.challenge": "42"}, "verify"), "42")
        with self.assertRaises(PermissionError): verify_subscription({"hub.mode": "subscribe", "hub.verify_token": "wrong", "hub.challenge": "42"}, "verify")

    def test_signature_uses_raw_body_and_rejects_tamper(self):
        body, secret = b'{"object":"whatsapp_business_account"}', "app-secret"
        signature = "sha256=" + hmac.new(secret.encode(), body, hashlib.sha256).hexdigest()
        verify_webhook_signature(body, signature, secret)
        with self.assertRaises(PermissionError): verify_webhook_signature(body + b" ", signature, secret)

    def test_webhook_normalizes_and_redacts(self):
        payload = {"object": "whatsapp_business_account", "entry": [{"id": "987654321", "changes": [{"field": "messages", "value": {"messaging_product": "whatsapp", "metadata": {"phone_number_id": "123456789"}, "messages": [{"id": "wamid.1", "from": "5491112345678", "type": "text", "timestamp": "1", "text": {"body": "private"}}], "statuses": [{"id": "wamid.2", "status": "delivered", "timestamp": "2", "recipient_id": "5491112345678"}]}}]}]}
        raw, secret = json.dumps(payload, separators=(",", ":")).encode(), "app-secret"
        signature = "sha256=" + hmac.new(secret.encode(), raw, hashlib.sha256).hexdigest()
        with tempfile.TemporaryDirectory() as temp:
            output = Path(temp) / "evidence"
            receipt = webhook_to_evidence(raw, signature, secret, output, profile=profile())
            normalized = (output / "normalized-events.json").read_text(encoding="utf-8")
            self.assertEqual(receipt["event_count"], 2); self.assertFalse(receipt["raw_payload_persisted"])
            self.assertNotIn("5491112345678", normalized); self.assertNotIn("private", normalized)

    def test_webhook_invalid_envelope_is_atomic(self):
        raw, secret = b'{"object":"wrong","entry":[]}', "app-secret"
        signature = "sha256=" + hmac.new(secret.encode(), raw, hashlib.sha256).hexdigest()
        with tempfile.TemporaryDirectory() as temp:
            output = Path(temp) / "evidence"
            with self.assertRaises(ValueError): webhook_to_evidence(raw, signature, secret, output, profile=profile())
            self.assertFalse(output.exists())


class ScopedWebhookTests(unittest.TestCase):
    def payload(self):
        return {"object": "whatsapp_business_account", "entry": [{"id": "987654321", "changes": [{"field": "messages", "value": {"messaging_product": "whatsapp", "metadata": {"phone_number_id": "123456789"}, "statuses": [{"id": "wamid.example", "status": "delivered", "timestamp": "1603086313", "recipient_id": "5491112345678"}]}}]}]}

    def run_payload(self, payload, output, configured=None):
        raw = payload if isinstance(payload, bytes) else json.dumps(payload).encode()
        signature = "sha256=" + hmac.new(b"app-secret", raw, hashlib.sha256).hexdigest()
        return webhook_to_evidence(raw, signature, "app-secret", output, profile=configured or profile())

    def test_independent_account_and_phone_scope_rejection(self):
        for target in ("account", "phone", "missing_account", "missing_phone", "mixed_batch"):
            with self.subTest(target=target), tempfile.TemporaryDirectory() as temp:
                payload = self.payload()
                entry = payload["entry"][0]
                metadata = entry["changes"][0]["value"]["metadata"]
                if target == "account": entry["id"] = "999999999"
                if target == "phone": metadata["phone_number_id"] = "999999999"
                if target == "missing_account": del entry["id"]
                if target == "missing_phone": metadata.clear()
                if target == "mixed_batch":
                    foreign = self.payload()["entry"][0]
                    foreign["id"] = "999999999"
                    payload["entry"].append(foreign)
                with self.assertRaises(PermissionError):
                    self.run_payload(payload, Path(temp) / "evidence")
                self.assertEqual(list(Path(temp).iterdir()), [])

    def test_malformed_status_is_atomic(self):
        variants = [("id", ""), ("id", "a" * 513), ("id", "private\ntext"), ("recipient_id", None), ("recipient_id", ""), ("timestamp", 1603086313), ("timestamp", "-1"), ("timestamp", "01"), ("timestamp", ""), ("status", "accepted"), ("status", "unknown"), ("status", {})]
        for key, value in variants:
            with self.subTest(key=key, value=value), tempfile.TemporaryDirectory() as temp:
                payload = self.payload()
                payload["entry"][0]["changes"][0]["value"]["statuses"][0][key] = value
                with self.assertRaises(ValueError): self.run_payload(payload, Path(temp) / "evidence")
                self.assertEqual(list(Path(temp).iterdir()), [])

    def test_malformed_containers_rejected(self):
        for target in ("entry", "changes", "metadata", "product", "statuses", "messages", "field"):
            with self.subTest(target=target), tempfile.TemporaryDirectory() as temp:
                payload = self.payload()
                entry = payload["entry"][0]
                change = entry["changes"][0]
                value = change["value"]
                if target == "entry": payload["entry"] = {}
                elif target == "changes": entry["changes"] = None
                elif target == "field": change["field"] = "other"
                elif target == "product": value["messaging_product"] = "other"
                else: value[target] = None
                with self.assertRaises(ValueError): self.run_payload(payload, Path(temp) / "evidence")

    def test_scope_profile_cannot_be_omitted_or_unproven(self):
        with tempfile.TemporaryDirectory() as temp:
            for key, value in (("business_account_id", ""), ("phone_number_id", 123456789), ("decision", "BLOCKED")):
                configured = profile(); configured[key] = value
                with self.subTest(key=key), self.assertRaises((ValueError, PermissionError)):
                    self.run_payload(self.payload(), Path(temp) / "evidence", configured)
            with self.assertRaises(TypeError):
                webhook_to_evidence(b"{}", "signature", "app-secret", Path(temp) / "evidence")
            self.assertEqual(list(Path(temp).iterdir()), [])

    def test_duplicate_json_keys_rejected(self):
        raw = json.dumps(self.payload()).replace('"status": "delivered"', '"status": "sent", "status": "delivered"').encode()
        with tempfile.TemporaryDirectory() as temp, self.assertRaisesRegex(ValueError, "duplicate JSON"):
            self.run_payload(raw, Path(temp) / "evidence")

    def test_local_body_and_event_budgets(self):
        with tempfile.TemporaryDirectory() as temp:
            with self.assertRaisesRegex(ValueError, "body budget"):
                self.run_payload(b" " * (1024 * 1024 + 1), Path(temp) / "body")
            payload = self.payload()
            value = payload["entry"][0]["changes"][0]["value"]
            value["statuses"] *= 1001
            with self.assertRaisesRegex(ValueError, "bounded lists"):
                self.run_payload(payload, Path(temp) / "events")
            value["statuses"] = value["statuses"][:600]
            payload["entry"][0]["changes"] *= 2
            with self.assertRaisesRegex(ValueError, "event budget"):
                self.run_payload(payload, Path(temp) / "aggregate")
            self.assertEqual(list(Path(temp).iterdir()), [])

    def test_status_identity_preserved_without_ordering_or_delivery_inference(self):
        with tempfile.TemporaryDirectory() as temp:
            payload = self.payload()
            statuses = payload["entry"][0]["changes"][0]["value"]["statuses"]
            statuses[:] = [{**statuses[0], "status": state} for state in ("read", "sent", "failed", "delivered", "deleted")]
            output = Path(temp) / "first"
            receipt = self.run_payload(payload, output)
            events = json.loads((output / "normalized-events.json").read_bytes())["events"]
            self.assertEqual([e["status"] for e in events], ["read", "sent", "failed", "delivered", "deleted"])
            self.assertEqual(len({e["event_key_sha256"] for e in events}), 5)
            self.assertEqual(receipt["schema"], "elite-whatsapp-cloud-webhook-receipt/v2")
            self.assertFalse(receipt["automatic_business_write"])
            self.assertEqual(receipt["normalized_sha256"], hashlib.sha256((output / "normalized-events.json").read_bytes()).hexdigest())
            # Format/order changes must not invent new provider-event identities.
            statuses.reverse()
            second = Path(temp) / "second"
            self.run_payload(json.dumps(payload, indent=2).encode(), second)
            later = json.loads((second / "normalized-events.json").read_bytes())["events"]
            self.assertEqual({e["event_key_sha256"] for e in events}, {e["event_key_sha256"] for e in later})
            for value in ("987654321", "123456789", "5491112345678", "wamid.example"):
                self.assertNotIn(value, json.dumps(events) + json.dumps(receipt))
            self.assertEqual(events[0]["recipient_sha256"], hashlib.sha256(b"5491112345678").hexdigest())

    def test_existing_output_is_not_overwritten(self):
        with tempfile.TemporaryDirectory() as temp:
            output = Path(temp) / "evidence"
            self.run_payload(self.payload(), output)
            before = (output / "normalized-events.json").read_bytes()
            with self.assertRaises(FileExistsError): self.run_payload(self.payload(), output)
            self.assertEqual(before, (output / "normalized-events.json").read_bytes())


class SendBridgeTests(unittest.TestCase):
    def frame(self, output):
        return {"schema": "elite-whatsapp-send-bridge/v1", "binding_sha256": "b" * 64, "profile": profile(), "request": message(), "access_token": "synthetic-secret", "output_directory": str(output)}

    def test_bridge_uses_existing_adapter_and_binds_receipt(self):
        calls = []
        def transport(method, url, headers, body, timeout):
            calls.append((method, url, json.loads(body)))
            return 200, {}, b'{"messages":[{"id":"wamid.bridge"}]}'
        with tempfile.TemporaryDirectory() as temp:
            output = Path(temp) / "evidence"
            result = send_bridge(json.dumps(self.frame(output)).encode(), transport)
            self.assertEqual(len(calls), 1)
            self.assertEqual(calls[0][0], "POST")
            self.assertEqual(calls[0][2]["to"], message()["recipient"])
            self.assertEqual(result["binding_sha256"], "b" * 64)
            self.assertEqual(result["provider_message_id"], "wamid.bridge")
            self.assertEqual(result["evidence_sha256"], hashlib.sha256((output / "SEND_RECEIPT.json").read_bytes()).hexdigest())
            self.assertNotIn("synthetic-secret", json.dumps(result))
            with self.assertRaises(FileExistsError): send_bridge(json.dumps(self.frame(output)).encode(), transport)
            self.assertEqual(len(calls), 1)

    def test_bridge_rejects_frame_without_provider(self):
        def forbidden(*args): self.fail("provider must not be invoked")
        with tempfile.TemporaryDirectory() as temp:
            for key, value in (("schema", "wrong"), ("binding_sha256", "bad"), ("profile", []), ("request", []), ("access_token", ""), ("output_directory", "relative")):
                frame = self.frame(Path(temp) / "evidence"); frame[key] = value
                with self.subTest(key=key), self.assertRaises(ValueError): send_bridge(json.dumps(frame).encode(), forbidden)
            frame = self.frame(Path(temp) / "evidence"); frame["extra"] = True
            with self.assertRaises(ValueError): send_bridge(json.dumps(frame).encode(), forbidden)
            with self.assertRaises(ValueError): send_bridge(b" " * 131073, forbidden)

    def test_uncertain_provider_never_retried_inside_bridge(self):
        calls = []
        def lost(*args): calls.append(1); raise TimeoutError("synthetic private transport")
        with tempfile.TemporaryDirectory() as temp, self.assertRaises(TimeoutError):
            send_bridge(json.dumps(self.frame(Path(temp) / "evidence")).encode(), lost)
        self.assertEqual(calls, [1])

    def test_ambiguous_provider_identity_is_not_accepted(self):
        with tempfile.TemporaryDirectory() as temp:
            for index, body in enumerate((b'{"messages":[{"id":"one"},{"id":"two"}]}', b'{"messages":[{"id":"bad\\nidentity"}]}')):
                with self.subTest(index=index), self.assertRaises(ValueError):
                    send_bridge(json.dumps(self.frame(Path(temp) / str(index))).encode(), lambda *args: (200, {}, body))

    def test_real_cli_fails_closed_without_logging_secret(self):
        frame = self.frame(Path(tempfile.gettempdir()) / "must-not-be-created-whatsapp")
        frame["profile"]["decision"] = "BLOCKED"
        result = subprocess.run([sys.executable, "-I", "-B", str(Path(__file__).with_name("whatsapp_cloud.py")), "--send-bridge"], input=json.dumps(frame).encode(), capture_output=True, timeout=10)
        self.assertEqual(result.returncode, 2)
        self.assertEqual(result.stdout, b"")
        self.assertEqual(result.stderr, b"WHATSAPP_SEND_UNCERTAIN_OR_REJECTED\n")


if __name__ == "__main__": unittest.main()
````

### FILE: `whatsapp_cloud/PROVENANCE.md`
```yaml
block_id: "PYTHON-META-WHATSAPP-CLOUD-ADAPTER:v343:4"
operation: CREATE
provenance: AUTHORED
source: "local exact adaptation record"
license: "LicenseRef-Workspace-Owner"
sha256: "c02fbbceb503df5329fa814ed8d24cad897c393ea2f0f7e77fe58ef1e28b4320"
variables: []
secrets_allowed: false
```
````markdown
# WhatsApp Cloud adaptation provenance

## V402 — governed conversational replies (candidate)

The active V402 claim extends the admitted owner to verified inbound text,
the existing conversation runtime, exact human approval and the existing send
fence/status ledger. The versioned sections below are historical; their
customer-message/UI/host-pending limitations are superseded only by the local
V402 evidence and explicit opt-in composition. No live account is certified.

- build_reply_payload is ADAPTED from get_text_message_input in the already
  pinned official Meta message_helper.py at de70ee908a67026e642aaee3703d20464e2a9466.
  The original stays VERBATIM; fields recipient/type=text/text.body are preserved.
- Window/shape guards, startup-profile stdio validation and local receipt
  recovery are AUTHORED integration glue, not Meta code. They reuse the exact
  verifier, stdlib transport, receipt format and fixed source-lock verification.
- Go projection, approval.request typed context, shared observer/read view,
  HTTP review/send/recovery and frontend are AUTHORED glue between admitted
  owners. No new runtime, approval ledger, outbound ledger or external module.
- Current official authority: https://whatsappbusiness.com/policy/ section2
  governs the 24-hour user-response window and escalation; Meta's official
  https://www.postman.com/meta/whatsapp-business-platform/folder/o48mro7/messages
  describes unique message IDs and status via webhooks. Snapshot hashes live
  in the V402 admission evidence. Developer-doc endpoints returned429, not proof.
- The existing contact resolver gains effective_at filtering before exposing
  a contact to the runtime. This enforces its admitted activation field; no new
  identity/consent policy is inferred.

Tests execute Go1.26.8, Python3.14 and synthetic HTTP/PostgreSQL only. Fixture
v99.0 is deliberately not a real Meta version assertion. Native fuzz covers
projection shape/scope, not HMAC, SAST, DAST or deployment security. Secrets,
terms, consent, retention and live provider operation remain target activation
responsibilities; the infrastructure does not fabricate their acceptance.


## 0.13.0 — concrete status transport, 2026-09-08

JSONStatusReporter and its regression/integration/fuzz tests extend the existing
status_host.go/status_host_test.go owners as AUTHORED Go standard-library code.
No Meta source, upstream snapshot, external module or licence changes. Existing
VERBATIM/ADAPTED Meta files remain byte-identical. GO-OBSERVABILITY-CORE is not
imported or promoted by this change.

Method sources: https://pkg.go.dev/context#AfterFunc (stopping a callback does
not wait for a running callback) and https://pkg.go.dev/net#Conn (concurrent Close
and deadlines). Execution pins Go 1.26.7; documentation's displayed newer release
does not change the toolchain. The code and tests are not attributed to Go/Meta.
Local TCP/PostgreSQL evidence is V343, with synthetic data and no external sends.

## 0.12.0 local host and read-only projection

status_host.go, notification_history.go, their tests and the browser harness are
AUTHORED composition of the admitted worker, authorization, PostgreSQL transaction,
existing notification projection and Microsoft Playwright runtime. No additional
Meta source, upstream revision, dependency or license is introduced.
Method sources consulted 2026-09-06: https://go.dev/blog/pipelines
(cooperative cancellation), https://sre.google/sre-book/monitoring-distributed-systems/
(actionable bounded reports), https://playwright.dev/docs/best-practices
(user-visible tests with isolated controlled data). These sources govern method,
not an assertion that the integration was written or certified by those companies.

The current section above governs the V278 local claim. The versioned sections
below retain their historical scope; earlier host/UI-pending statements are
superseded only for this tested local integration, never for target deployment.

## V277 scoped status job completion — 2026-09-06

status_worker.go, its tests, router transaction callback and shared jobs.go
extensions are AUTHORED. No new code copied from Meta/AWS/PostgreSQL, upstream
revision, dependency or protocol. Existing Meta VERBATIM references and adapted
signature/status verifier retain their exact bytes, license and source lock.

Official method authorities consulted 2026-09-06:
- https://aws.amazon.com/es/builders-library/leader-election-in-distributed-systems/
  for lease checks and pause/failure considerations, not this implementation.
- https://www.postgresql.org/docs/18/tutorial-transactions.html
  for all-or-nothing local transaction semantics, not remote exactly-once.

Scoped claiming and transactional completion/failure extend the existing queue
owner. Router/observer remain the only identity resolution/observation owners.
Tests demonstrate local contracts and injected recovery, not endorsement by
those organizations, live account readiness or REUSABLE_PACK promotion.

## V276 retained-status routing — 2026-09-06

status_router.go, its tests and index migration 0052 are AUTHORED, not code
published by Meta or Google. They reuse the existing Meta-adapted verifier,
contactidentity HMAC function, outbound approval/fence and StatusObserver. No
new dependency, provider protocol, parser authority, queue or upstream pin.

Official authorities checked on 2026-09-06: Meta/Postman Message Status Update
Notifications (id, recipient_id, timestamp and delivery states), and
https://go.dev/blog/osroot (Go's traversal-resistant file APIs, including platform
limitations). The router uses Go 1.26.7 os.Root already available in the pinned
runtime, not a locally invented filesystem sandbox. Windows fixtures are not
cross-platform security certification. Embedded Meta source/license hashes stay
unchanged; runtime integration remains local responsibility and CONDITIONED.

## V274 HTTP-to-durable-inbox composition — 2026-09-06

webhook_receiver.go, its tests and ingress_bridge are AUTHORED. The new stdio
entry point reuses the existing Meta-adapted normalize_verified_webhook and
verify_subscription; it is not new official Meta code. The pinned Meta examples,
license and source lock are unchanged. The receiver reuses provider inbox/job
0.1.2 instead of adding another queue or claiming that the reference Elite HMAC
is Meta's signature protocol. Raw-byte retention is opt-in and explicit; base64
does not provide encryption. Live access/retention/worker/UI gates remain open.

The official Meta/Postman status payload page was available on 2026-09-06.
The current developers.facebook.com overview returned HTTP 429 and its .md
variant was unavailable. The official Node SDK webhook documentation explicitly
states that it is archived; it is historical GET/challenge/POST/200 protocol
context, not an admitted runtime or proof of current retry policies. No archived
SDK was added. Fresh live callback/subscription contract probes remain required.

## V272 interruption and stale-scope regression — 2026-09-06

Three new Go tests (one with three subcases) are AUTHORED. They exercise the
unchanged StatusObserver, isolated Python verifier and PostgreSQL transaction.
Google SRE Testing for Reliability (https://sre.google/sre-book/testing-reliability/)
supports fault-injection testing; PostgreSQL 18 transaction and row-lock docs
define the atomicity/locking semantics under test. Neither source publishes
these tests or certifies this integration. No new Meta code, dependency, schema,
public ingress, production authorization or reusable-pack promotion is claimed.
Connection-pool replacement is not a PostgreSQL-server crash or a host restart.

## V271 durable observation integration — 2026-09-06

StatusObserver, SQL migration 0051, tests and the reader extension are AUTHORED.
The isolated process still reuses the existing Meta-adapted signature/parser and
V270 correlation code. New database tables store only evidence/observations;
the existing outbound delivery owner remains the only send fence.

Meta's official Message Status Update Notifications page below was rechecked:
provider timestamps, not callback arrival, govern observation order. PostgreSQL
18 row locking (https://www.postgresql.org/docs/18/explicit-locking.html) governs
the transaction's FOR UPDATE/FOR SHARE lifetime. Python's official importlib
source-file recipe (https://docs.python.org/3/library/importlib.html#importing-a-source-file-directly)
supports loading a fixed local module under isolated mode; it is not a license
to load caller-supplied paths. These sources do not publish this authored service
or certify production. Upstream Meta bytes, license, notices and dependencies
are unchanged; scripts have new hashes and must be revalidated in target pins.

## V270 anchored status observations — 2026-09-06

status_reconciliation.py and test_status_reconciliation.py are AUTHORED. They
reuse normalize_verified_webhook, extracted without changing the existing
Meta-adapted HMAC/scope/parser contract or evidence v2. This is not a new Meta
SDK, copied company implementation or provider certification. Official sources,
license/notices and upstream bytes remain unchanged.

[Meta Message Status Update Notifications](https://www.postman.com/meta/whatsapp-business-platform/request/rgtfq23/message-status-update-notifications)
was consulted on 2026-09-06: callbacks identify message/recipient and carry
status/timestamp; their arrival order can differ from status time. The existing
admitted fbsamples commit de70ee908a67026e642aaee3703d20464e2a9466 supplies
signature-validation reference code, not this local reconciliation module.

Trusted send anchor, atomic snapshots, local budgets, deduplication-conflict
handling and explicit tie ambiguity are local engineering decisions and are
tested as such. The module preserves observations without inventing a provider
state machine, a retry right or a global latest status. No PostgreSQL business
state or delivery fence is changed. Rebuild evidence lives in the library at
reconstruction_evidence/WHATSAPP_ANCHORED_STATUS_RECONCILIATION_V270.md.

## V268 read-only outcome — 2026-09-06

The status reader/tests and shared HTTP authentication helper are AUTHORED.
They reuse the same approval and outbound tables without a migration, new
dependency, provider call or invented delivery transition. Meta sources and
the executable Python adapter remain unchanged.

[Meta status notifications](https://www.postman.com/meta/whatsapp-business-platform/request/rgtfq23/message-status-update-notifications)
separate sent/delivered/read with timestamps; a local acceptance is not proof
of delivery. [PostgreSQL 18 snapshots](https://www.postgresql.org/docs/18/transaction-iso.html)
explain the single SELECT's committed snapshot. Both consulted 2026-09-06.
Neither authority publishes this local reader or certifies the complete journey.

## V267 authenticated dispatch — 2026-09-06

`appointment_notification.go` and its tests are AUTHORED API composition, not
Meta/AWS source. They reuse the V266 resolver, existing identity.Verifier,
existing Channel/PostgreSQL fence and V265 Sender/Python adapter unchanged.
They introduce no SQL owner, queue, framework, dependency or upstream pin.

Official [Go encoding/json](https://pkg.go.dev/encoding/json) documents
case-insensitive field matching and decoder behavior; the local flat contract
rejects aliases, duplicate keys, nulls and trailing values explicitly.
The AWS outbox method and WhatsApp policy references below were consulted again
on 2026-09-06. Their documentation does not publish or certify this API module.

Tests use real HTTP, local signed RS256/JWKS and the existing OIDC verifier,
real PostgreSQL and a clearly synthetic pinned Python child. SQL appointment
and consent fixtures do not prove a real captured consent, login or Meta send.

## V266 durable appointment approval — 2026-09-06

`appointment_approval.go`, its tests and migration/test 0050 are AUTHORED,
not source published by Meta or AWS. They connect existing appointment/outbox,
contact identity and consent owners to the V265 delivery bridge. No upstream
bytes, license, runtime or dependency admission changed.

[AWS transactional outbox](https://docs.aws.amazon.com/prescriptive-guidance/latest/cloud-design-patterns/transactional-outbox.html),
consulted 2026-09-06, supports acting on committed events and retaining
idempotency; it does not publish this resolver or certify its implementation.
[WhatsApp Business Messaging Policy](https://whatsappbusiness.com/policy/)
governs actual permission and opt-out. Database evidence is not legal approval
or proof that a business captured valid consent. The target must establish it.

The immutable grant binds tenant, organization, confirmed appointment/event,
contact version, fixed policy/purpose, current exact consent, message/profile
hashes, actor and expiry. The existing PostgreSQL fence owns send replay.
Tests use synthetic SQL domain fixtures and a synthetic pinned provider child,
not real Meta or the authenticated appointment confirmation HTTP journey.
Revalidation is a pre-send snapshot, not atomic revocation of in-flight sends.

## V265 bridge — 2026-09-06

`internal/whatsappbridge/sender.go` and its tests are AUTHORED integration.
`whatsapp_cloud.py` adds a bounded stdio entry point around its existing ADAPTED
template sender; no second HTTP implementation is added. The three Meta source
files, license text and upstream revision remain unchanged. No new dependency
or runtime version is admitted merely by adding the bridge.

The caller requires hash-bound workflow approval and validates the child receipt.
The existing Go/PostgreSQL fence owns replay and ambiguous outcomes. The CLI
does not create consent, an approval database or an independent retry loop.
Official method: [AWS idempotent API guidance](https://aws.amazon.com/builders-library/making-retries-safe-with-idempotent-APIs/),
consulted 2026-09-06, governs deliberate request identity and ambiguous effects;
it is not a claim that Meta implements AWS's idempotency semantics.
[Meta status notifications](https://www.postman.com/meta/whatsapp-business-platform/request/rgtfq23/message-status-update-notifications)
govern the distinction between accepted/sent and subsequent delivery evidence.

Process tests use a synthetic pinned child, not a hidden live-send switch in
the production adapter. Python tests inject transport only into the public
testable function. The production CLI never selects a fake provider.

Authority: `fbsamples/whatsapp-api-examples@de70ee908a67026e642aaee3703d20464e2a9466`, signed GitHub commit, official Meta repository. License is restricted to use with Facebook web services/APIs and Platform Policy; preserve `upstream/LICENSE` in every copy.

Reference files and SHA-256 are recorded in `official-source.lock.json`. The three code files are byte-verbatim. Markdown materialization requires a terminating newline, so the root LICENSE copy is explicitly `ADAPTED_FINAL_NEWLINE_ONLY`: upstream SHA-256 `ef5c10ee…9300`, packaged normalized SHA-256 `48d97b3c…d01f`; its legal text is unchanged. They are evidence, not the executable production path. The official signature Python example contains an unconditional invalid-signature return after its comparison; the active e-commerce JavaScript example validates `X-Hub-Signature-256` but directly mutates shared sample objects and logs/handles messages without durable idempotency.

`whatsapp_cloud.py` is explicitly `ADAPTED`, not represented as unmodified Meta code. It retains the official Graph `/{VERSION}/{PHONE_NUMBER_ID}/messages` Bearer/JSON template pattern, GET subscription challenge and raw-body HMAC-SHA256 pattern. Changes: fail-closed typed profile; exact template/language/arity allowlist; secrets only through environment references; constant-time full-header comparison; bounded timeouts; injected HTTP transport; atomic response/receipt; hashed identifiers; normalized webhook evidence; no raw inbound PII persistence; no automatic domain write. The outbound provider-response.json is raw provider data and requires access/retention controls; hashes are pseudonyms, not anonymization. Project provider edge/outbox/inbox remains responsible for durable business idempotency and reconciliation.

## V264 scoped webhook correction — 2026-09-06

The existing Meta revision/license and four packaged reference files are unchanged.
The new validation, local budgets and tests are local changes to the `ADAPTED`
integration, not code copied from Meta. Before correction, a signed event with
another account/phone was accepted by the regression fixture.

Official contracts consulted on 2026-09-06:

- [Meta webhook payload reference](https://www.postman.com/meta/whatsapp-business-platform/folder/vzaxn16/webhook-payload-reference): entry.id identifies WABA; value.metadata.phone_number_id identifies the phone; value.messaging_product and field identify the envelope.
- [Meta message status notifications](https://www.postman.com/meta/whatsapp-business-platform/request/rgtfq23/message-status-update-notifications): statuses include message id, recipient_id, status and timestamp. Arrival order can differ from event timing.

The adaptation now requires the approved profile as a keyword argument, validates
its business_account_id and phone_number_id, and rejects a mixed/foreign batch
without publishing partial evidence. Events v2 preserve hashed account, phone,
contact and provider message identity, timestamp, exact admitted status and a
deterministic event key independent of batch order/JSON whitespace. The receipt
binds the profile and normalized bytes by SHA-256. Unknown shapes/states fail
closed; they require contract review, not silent normalization to delivered.

Local limits are 1 MiB raw body, 1,000 array entries/events, 512 ASCII identity
characters and canonical positive Unix seconds up to 12 digits. These are
implementation budgets, NOT Meta's limits or a timestamp-freshness guarantee.
No sorting, business-state transition, durable deduplication, tenant assignment,
contact binding, notification dispatch or automatic retry is introduced here.
The common inbox/outbound fence remains the owner of those effects.
````

### FILE: `whatsapp_cloud/README.md`
```yaml
block_id: "PYTHON-META-WHATSAPP-CLOUD-ADAPTER:v343:3"
operation: CREATE
provenance: AUTHORED
source: "local operator contract"
license: "LicenseRef-Workspace-Owner"
sha256: "6ae1fc4c0800e02773899e3ff80550b70452593305e479372779f4eaf3bdd2bf"
variables: []
secrets_allowed: false
```
````markdown
# Meta WhatsApp Cloud API adapted integration

## V402 — receive, propose, review, send and reconcile

Explicitly enable ConversationRoute on the existing StatusRouter. One durable
provider-events worker now routes signed mixed text/status envelopes. Original
signed bytes are preserved and reverified; an unbound contact reaches the
existing durable handoff without model/tool/provider calls. Runtime.Handle
stores a proposal and never sends it. The old status-only mode remains default.

Human review in /franchise/whatsapp binds the exact text, recipient, tenant,
organization, connection, contact version, consent, profile and user-message
window through the shared approval.request/decision owner. Approval and sending
are separate actions. A rejected decision cannot be changed by retry. A later
binding/consent change or closed window blocks a new send.

Provider acceptance, delivery and read status remain separate. Unknown sends
never auto-retry. Recovery verifies an existing bounded local SEND_RECEIPT and
provider-response against the exact approved request, then uses the same
outbound owner's reconciliation; it never performs HTTP or invents a messageID.
A crashed expired sending lease first follows the existing unknown transition.
No local receipt means explicit review, not an inferred acceptance or secondPOST.

See docs/whatsapp-conversation-operations.md for exact host/API activation.
The profile template stays blocked and empty. Startup uses
--validate-profile-bridge with the exact file bytes as base64 and validates them
through the existing validate_profile. This proves configuration structure only.
The host must still obtain actual identity/tokens and use an authenticated
reporting channel; it must not synthesize an approved account or service identity.

The sections below preserve earlier releases. Their no-customer-routing and
UI/host-pending statements describe those versions, not an automatic V402
production entitlement. Unsupported media and unmappable status events remain
retained for review; no business/data retention permissions are inferred.


## 0.13.0 — bounded JSON transport for the existing status host

NewJSONStatusReporter(conn) supplies the existing StatusReporter interface with
an AUTHORED standard-library transport. Pass an already-established net.Conn;
the reporter owns all writes, write deadlines and Close exclusively. The caller
must select/authenticate the collector (including TLS where required), retain the
reporter for the host lifetime and close it at shutdown. No dial, global logger,
service identity, daemon, collector, alert rule or live subscription is installed.

Report requires a context deadline, capped at three seconds including waiting
for a concurrent report. The host already supplies that deadline. The six JSON
fields are outcome, claimed, failure_recorded, inserted_observations, elapsed_ns
and next_delay_ns. Outcome accepts only AUTH_UNAVAILABLE, IDLE,
RECONCILE_REQUIRED, TERMINAL_REVIEW, COMPLETED or RETRY_RECORDED. Counts and elapsed
time are nonnegative; next delay is one second through ten minutes. No raw error,
payload, tenant/contact/provider ID, token, free-form label or message is accepted.

One complete Write emits one JSON line. Short/invalid counts, errors or active
write cancellation return static ErrStatusReport and close the stream; Run then
returns ErrStatusHost without retrying. A fresh connection/reporter is required
for future reports. An uncertain record must not be replayed automatically.
Rejected input or an expired waiter does not poison another active write.
Cancellation relies on the supplied net.Conn honoring Close and write deadlines.
The cancellation callback is joined before releasing the writer gate. Success
proves local Write acceptance only, not remote acknowledgement or durable retention.

Use reporter, err := NewJSONStatusReporter(approvedConnection), defer
reporter.Close(), then worker.Run(ctx, realPrincipalSource, reporter, interval),
checking both constructor and Run errors in the owning service. A supervisor
must observe Run failure independently of the failed report connection. This
recipe deliberately provides no production identity or transport credentials.

Local verification covers concurrent loopback TCP, exact schema, partial and
invalid writes, early cancellation, deadline expiry, terminal streams and the
real PostgreSQL status worker feeding TCP. The AUTHORED test fixture is not a
deployed service. Run the existing Go suite with the dedicated PostgreSQL/Python
variables described below, and the finite FuzzJSONStatusReporterReportInput
target through GO-NATIVE-FUZZ-GATE. TLS/auth, collector durability, retention,
supervision, alerts, race/SLO and live Meta acceptance remain target gates.

## 0.12.0 — opt-in host and operator status history

The existing StatusWorker now exposes Run(ctx, principalSource, reporter, interval).
Configure one host per worker; the host authenticates/revalidates its service
principal on each poll and reports bounded outcomes without payload/contact/token
fields. Both callbacks must honor context: authentication has a five-second budget,
reporting three seconds. Failed polls back off with jitter up to ten minutes.
Reporting failure stops the host; cancellation prevents another claim. This is
cooperative cancellation, not a force-kill guarantee for a broken callback.
No daemon, service credential, subscription or monitoring destination is installed.

GET /v1/franchise/appointments/{id}/whatsapp-confirmations returns the exact
authorized appointment history in one read-only repeatable-read snapshot. More
than twenty approvals is an explicit conflict, not a silently truncated history.
The existing appointment operator portal can read this through its canonical BFF;
it cannot send, retry or change a business state from this view.
Enable features.whatsapp_status_history in the existing server-side business
configuration only after mounting/probing this Go module. Default absent/false
hides the view and closes its BFF without a backend call. It grants no permission
or consent and creates no new configuration registry.
Use docs/whatsapp-status-operations.md (view version whatsapp-status-view/1.0.0).

Verify with ELITE_WHATSAPP_TEST_DATABASE_URL pointing only to the loopback
disposable database and ELITE_WHATSAPP_PYTHON pointing to Python:
go test ./internal/whatsappbridge -v -count=1 -timeout=5m
For the four-browser connected fixture, additionally set
ELITE_NOTIFICATION_BROWSER_E2E=1 and ELITE_WEB_ROOT to the built composition.
Install the isolated browser gate with pnpm install --ignore-workspace
--frozen-lockfile. The opt-in gate fails if its exact CLI/build is absent.

This local integration is AUTHORED. It does not certify live Meta delivery,
production service identity, alert delivery, operator training completion,
deployment, edge security, load or the complete franchise journey.

The current section above governs the V278 local claim. The versioned sections
below retain their historical scope; earlier host/UI-pending statements are
superseded only for this tested local integration, never for target deployment.

## Scoped status job worker (0.11.0)

This section supersedes earlier worker-pending statements for retained STATUS
batches only. Compose GO-RELIABLE-ASYNC-WORKERS 0.3.0 and this pack together;
apply the existing 52 migrations. No new queue, schema or dependency is added.
Create the observer/router below, then NewStatusWorker(router, uniqueWorkerID,
lease, retryDelay). Call ProcessOnce(ctx, genuinelyVerifiedPrincipal) from an
explicitly configured host. Nothing starts automatically. Lease 1s..2min,
retry 0..1h and one call per instance are local budgets, not provider limits.
Use a nonzero retry delay/backoff in the host; stop polling on shutdown and
inspect errors. Do not put this event type behind a generic auto-ACK worker.

Claims filter tenant, connection, provider, queue, job type, schema and event
type before taking work. Existing router/observer owners revalidate retained
bytes, current authority and exact receipt identity. Completion commits job,
inbox processed and safe audit together, with the exact claim generation, live
lease and payload/metadata still matching. Observation commits are separate:
on a later failure they remain durable and replay without duplicate observations.
An insertion count is NOT completion; check Completed and the returned error.

Recoverable failures requeue within max_attempts. Terminal attempts set job and
inbox failure with audit; a crashed expired final attempt is quarantined as
LEASE_EXHAUSTED without starting another effect. Existing job/inbox tables are
the terminal store, not a new DLQ. Audit contains actor, scoped job ID, attempt,
safe reason, terminal flag and event hash; never raw messages/contacts/secrets.
FailureRecorded=false means no durable failure record was proven (e.g. lost
claim, changed payload or unavailable DB). The host MUST surface a safe error
and reconcile the current job; never acknowledge it, hide it, reset attempts or
overwrite a newer claim. Configure authorized alerting and retention in target.

Unknown or mixed customer-message batches are not silently discarded or sent
to an invented handler. They stay unresolved/terminal for authorized triage.
No replies, sends, business acceptance or automated terminal redrive occur.
Diagnose using existing routing guidance; preserve evidence and generations.
Rollback pauses this host and preserves jobs/inbox/audit/observations/fences;
never revert to pre-0.2.0 unsafe generic workers or clear terminal history.

Verified locally: scoped exclusion, concurrent consumers, retry/exhaustion,
crash recovery, failure between completion writes, stale claim and payload
drift. These are injected faults with real PostgreSQL, not a killed production
host. Frontend, operational redrive policy, live Meta/account/consent/retention,
capacity, security, deployment and complete business acceptance remain gates.
Code is AUTHORED; source/license pins and Meta-adapted verification are unchanged.

## Retained-status routing (0.10.0)

After migrations through 0052, use `NewStatusRouter(observer, connectionID,
retentionApprovalSHA256, outboundHMACKey, maxRoutes)` and
`ObserveRetained(ctx, verifiedPrincipal, providerEventID)`. The event ID comes
from the existing durable inbox/job, not a user-supplied tenant or appointment.
The caller must supply a genuinely verified human/service identity with
appointment:manage and the current organization grants; never manufacture one
from webhook fields. Keep pool capacity >=2 and budget capacity across router
instances/other consumers. One concurrent call per router and 1..8 distinct
message/recipient routes per batch are local budgets, not Meta limits.

The router reads only an active matching connection and retained receipt,
checks raw/profile/retention hashes and reuses the existing pinned Python
signature/scope verifier. Only then does it project exact, case-sensitive JSON
identity keys, HMAC them using the outbound owner's key, and find exactly one
accepted approved delivery. Connection organization scope and principal grants
must match. All routes and anchored receipts are resolved before observations.
Receipts use the Sender's deterministic path, Go os.Root confinement, a 64 KiB
budget and exact DB/hash/message/recipient binding. Protect the configured
evidence root, runtime files and database from untrusted writers; filesystem
confinement is not content authenticity, authorization or a kernel sandbox.

StatusObserver remains the only observation writer. Replays add no duplicate
events. A later failure may leave earlier observations committed: preserve the
inbox/job and replay, do not resend. The router never completes a job, marks an
inbox processed, replies to a customer or changes business state. Mixed inbound
messages, unknown/ambiguous send identities, mismatched receipts and unsupported
batches remain unresolved for a suitable handler or operator. The caller must
check the error, not interpret an insertion count (including zero) as completion.

For ErrStatusRouting diagnose with safe hashes only: check active connection and
organization/permissions; retained body/profile/policy binding; pinned runtime
and secret availability; accepted delivery/HMAC-key identity; receipt existence
and digest. Never log raw payload, provider IDs, contacts or secrets. On failure
keep evidence and apply the library failure/recovery protocol. SQL index 0052
supports lookup but throughput/SLO still require measurement in the target.

This is a service integration, not the completed queue worker. Admission still
requires queue claiming/completion/retry/DLQ/retention, frontend/live Meta journey
and operation. Do not register it automatically. Rollback stops its consumer and
preserves raw receipts, observations and send fences; dropping only index 0052
does not delete business data. Rebuild and verify before resuming.

## Opt-in HTTP reception into the existing provider inbox (0.9.0)

In the complete Go/PostgreSQL profile, compose NewWebhookReceiver with a fixed
WebhookReceiverConfig: TenantID, ConnectionID, raw approved Profile, hash-pinned
Process, server-side AppSecretSource and VerifyTokenSource, Store from
postgres.NewProviderIntegration(pool), MaxConcurrent (1..16), and
RetentionApprovalSHA256 referencing the project's approved raw-data policy.
Missing retention approval is rejected. The hash is a link, not self-proving
consent: the agent must validate the referenced project decision before exposure.
Use GO-PROVIDER-INTEGRATION-CORE 0.1.2 or a separately verified compatible fix;
earlier versions do not enforce current connection state/replay scope.

Mount that http.Handler on one explicitly configured callback URL. No callback
body, URL argument or header selects a tenant. Provision an active
integration.provider_connection with the same tenant/connection and local
provider_code meta-whatsapp; supply secrets externally. Do not use the generic
X-Elite-Webhook-Signature protocol for Meta. GET verifies the subscription token
and returns the bounded numeric challenge; POST invokes the existing Python
Meta verifier under -I/-B and verifies exact body/WABA/phone/profile scope.

The receiver accepts <=1 MiB raw JSON, bounded parallel work, a five-second
processing deadline and bounded stdio. Configure server TLS, header/read/idle
timeouts, edge rate limits and safe access logs as well; the handler alone does
not establish a production ingress SLO. Never log hub.verify_token or raw bodies.

Only after the shared event+job SQL commit does POST return 200 EVENT_RECEIVED.
Busy/unavailable/disabled storage returns 503, invalid signature/scope 403,
invalid query/duplicate signature 400, media/encoding 415, body budget 413.
It does not retry sends, promote a lead, complete a job or claim delivery.
Provider event identity is wa:<SHA256(exact body)> within tenant/connection.
Exact redelivery does not create another event/job. Profile, signature or
retention-policy changes on the same body fail closed as receipt conflicts;
rotation/reconfiguration must reconcile retained work, not overwrite history.

Unlike the offline pseudonymous evidence writer, this opt-in HTTP lane retains
the exact original bytes as base64 plus signature/profile/policy hashes and
verified count in integration.webhook_event.payload, schema
elite-whatsapp-retained-webhook/v1. BASE64 IS NOT ENCRYPTION. This is sensitive
raw data, including possible message text/contacts. Restrict database/backups,
prove encryption at rest, implement the approved retention/deletion policy and
audit worker access before enabling the endpoint. No key, consent or retention
duration is invented by materialization. App secrets/tokens are not persisted.

The existing provider-events queue receives provider.webhook.received with only
connection/provider/event identifiers. Do not let a generic worker acknowledge
this event_type until a WhatsApp-specific handler has loaded the exact retained
bytes, revalidated current authority/profile/signature, resolved the approved
tenant/message/appointment and committed the corresponding effect. The current
tests explicitly supply that routing; the automated routing/worker, retry/DLQ
operations, operator UI and live subscription are still project/library gaps.

The fixture proves HTTP lost ACK after commit, four concurrent replays, one
event/job, original-byte reload -> StatusObserver -> scoped status read, and 14
negative HTTP cases. This is not a live Meta call or automated worker certification.
Run the complete internal/whatsappbridge suite with its existing dedicated
loopback PostgreSQL and Python environment requirements. See V274 evidence.

## Durable observations and scoped GET (0.8.1)

After migration 0051, compose a tenant-bound whatsappbridge.StatusObserver
with the existing PostgresAppointmentApprovals, approved raw Profile, Process,
ReconcilerSHA256 and AppSecretSource. Freeze configuration before concurrent
use. Process pins Python and whatsapp_cloud.py; ReconcilerSHA256 additionally
pins status_reconciliation.py. Both scripts and their package must be protected
from writes in deployment. Secret lookup stays server-side, never in CLI args.

Call Observe(ctx, verifiedPrincipal, organization, appointment, confirmationEvent,
sendReceiptBytes, signedWebhooks) from the authorized service/durable ingress
owner. A Principal assembled from request JSON is not authentication. The method
loads the trusted SEND_RECEIPT anchor from the existing accepted outbound fence;
the caller cannot supply that anchor. Unknown sends without an accepted receipt
remain blocked. Scope, approval profile and fence integrity must match.

The same hash-locked Python code runs with -I/-B and a minimal environment,
reverifies signatures and returns bound observations over size-limited stdio.
No provider request or send occurs. A 20-second deadline includes verification
and persistence. The SQL transaction rechecks current scope, locks the fence
and current appointment/lead, then inserts append-only verification batches and
unique provider events. A conflicting event rolls the whole transaction back.
Duplicates return zero inserted events; the send fence and attempt count never
change. Caller acknowledgement of an ingress message must happen only after
success; this module does not invent or activate a queue/worker/retry schedule.

The existing authenticated GET now reads the observations in the same SQL
snapshot. It adds provider_event_count and optional provider_timestamp (Unix
seconds); delivery_status is observed_sent/delivered/read/failed/deleted, or
ambiguous_latest_timestamp for equal latest timestamps with different states.
No observations retains not_observed_by_this_reader. These are recorded provider
observations, not a globally current status, complete callback history, delivery
SLA, business-state change or permission to resend. No raw IDs/hashes/phones are
returned. Migrations are required before using the updated reader.

No public ingestion endpoint, subscription, scheduler, UI or live account is
activated here. The durable owner must supply raw authenticated-channel input,
tenant/object routing, retention, least-privilege DB roles, recovery and monitoring.
Evidence receipts contain pseudonymous hashes: they still require privacy controls.
The 0051 down migration removes these observation tables; use it only with an
approved evidence backup/export or disposable test data, never blindly in production.

Verification: ELITE_WHATSAPP_PYTHON and the existing dedicated loopback
ELITE_WHATSAPP_TEST_DATABASE_URL guard, 51 migrations, then
go test ./internal/whatsappbridge -v -count=1. Tests use synthetic Meta payloads
and secrets, the real isolated Python verifier, PostgreSQL and signed local
RS256/JWKS HTTP GET, including read-only SQL sessions. See V271 evidence.

The 0.8.1 regression suite also changes lead scope, fence state and receipt anchor
after the initial lookup; the transactional recheck must reject all three. It
observes a real blocked event INSERT through pg_stat_activity, cancels the call,
asserts that neither its batch nor event survived, and then proves safe recovery.
A fresh connection pool/service can replay a committed batch without another
event or send attempt. These are local synthetic tests, not a database-server
crash, live IdP revocation, provider redelivery, backup/restore or production drill.
No additional inbox, sender ledger, migration or runtime dependency is introduced.

## Anchored signed-status reconciliation (0.7.0)

The existing send bridge returns evidence_sha256 for SEND_RECEIPT.json. The Go
outbound fence stores it as EvidenceSHA256. Use that trusted, tenant/object-scoped
anchor to correlate original signed webhook bytes with the exact sent message:

```python
from status_reconciliation import reconcile_status_webhooks

result = reconcile_status_webhooks(
    profile=approved_profile,
    send_receipt_bytes=retained_send_receipt_bytes,
    expected_send_receipt_sha256=trusted_fence_evidence_sha256,
    signed_webhooks=original_body_and_signature_pairs,
    app_secret=app_secret_from_secret_store,
    output_directory=absent_evidence_directory,
)
```

These inputs must come from the application's authorized evidence/ingress
owners, not arbitrary browser input. Do not compute a new expected hash from
an untrusted receipt: the anchor must already exist in the durable fence.
The function validates the send anchor/profile and reuses the one existing
signature/scope/parser boundary. It does not accept normalized JSON as proof
of a signature. Retain original webhook bytes under the selected privacy
policy until verified processing completes; they are only supplied in memory
here and are not persisted by this function.

It emits STATUS_RECONCILIATION.json and status-observations.json in one new
snapshot directory. Unrelated messages are excluded, a matched ID with another
recipient rejects, exact repetitions deduplicate, and identical event keys
with divergent provider evidence reject. All supplied batches must validate
before publication. Local budgets: eight batches, one MiB/1,000 events each,
and 64 KiB for the send receipt. These are local limits, not Meta quotas.

The observation timeline uses provider timestamps, not arrival order. Equal
latest timestamps with different states yield AMBIGUOUS_LATEST_TIMESTAMP;
there is no guessed precedence. NOT_OBSERVED_IN_INPUT means only that these
inputs contain no matching status. MATCHED_OBSERVATIONS and last_observed_status
describe the supplied snapshot, not a globally current status. All distinct
states are retained, including failed/deleted; none authorizes a resend.

No network call, SQL update, durable inbox acknowledgement or payment/business
mutation occurs. SHA digests are pseudonymous identifiers, not anonymous data;
access, retention, tenant scope and external anchoring remain mandatory.
In 0.8.0 StatusObserver can persist this evidence and the GET reader consumes it;
mounting the ingress runner and operator UI in the project remains pending.
This function does not close an unknown send if no trusted SEND_RECEIPT exists.

Run python -m unittest discover -v from whatsapp_cloud. Tests connect the real
send-bridge artifact path to signatures generated with a synthetic app secret,
including all six sent/delivered/read arrival permutations. No live Meta account
or delivery is claimed. Source/method and limitations: PROVENANCE.md, V270.
The changed Python adapter bytes require revalidating its hash in the project's
Sender configuration through its normal dependency/change gates; do not bypass
the hash check or silently overwrite a deployed pin.

## Read-only notification outcome (0.6.0)

The same module now registers:
`GET /v1/franchise/appointments/{id}/whatsapp-confirmation?organization_id=...&confirmation_event_id=...`.
Use the same Bearer verifier and appointment:manage permission. Tenant is still
server-bound; organization and the current appointment/lead association must
match the historical grant. No body, unknown or repeated query parameters.

The reader makes one SELECT over existing owners and returns no phone, template,
actor, consent ID, provider ID or evidence hashes. Its fields are delivery_key,
fence_state, delivery_status, approval_expires_at, approval_expired, accepted_at
and updated_at when present, observed_at and reconciliation_required.

- not_started: an approval exists, no matching local fence row is observed;
- sending: stored attempt state; an expired lease sets reconciliation_required
  without mutating the row or inferring that a second send is safe;
- accepted: stored provider acceptance, NOT confirmed delivery;
- unknown: reconcile through the existing owner before any further effect;
- failed_terminal: recorded terminal outcome, not an instruction to retry.

In the historical 0.6.0 reader delivery_status was always
`not_observed_by_this_reader`; 0.8.0 adds the durable observations above. Absence
of local observations does not claim that no such evidence exists elsewhere.
A grant that expired, a cancelled appointment or revoked
sending consent does not hide its historical outcome from an operator who still
has current scope. Reading is not authorization to send. Reassigned appointment
or lead scope is rejected; retention/privacy policy remains a target obligation.

GET is bounded to five seconds, returns no-store and 404 for absent/scoped-out
records, 409 for grant/fence hash mismatch, 503 for unavailable storage. It never
calls Claim, ResolveWhatsAppApproval, a token source or the provider. Do not use
POST to discover history, nor change event/expiry to force a resend. The tests
exercise reads with transaction_read_only=on and assert unchanged event counts.

## Authenticated notification operation (0.5.0)

The composed Go profile provides `NewAppointmentNotificationModule`.
Construct it with the existing PostgresAppointmentApprovals, tenant-bound
Sender, PostgreSQL OutboundDeliveryStore and verified WhatsApp receiver. It
copies the sender/profile and forces the same approval resolver into the
sender. Register it on the existing mux with the existing identity.Verifier:

```go
module, err := whatsappbridge.NewAppointmentNotificationModule(approvals, sender, deliveryStore, receiver)
if err != nil { return err }
module.Register(mux, verifier)
```

No listener, cron, worker or account is activated by materialization. Register
once for the configured tenant on its existing host/router; multi-tenant routing
must select server-approved configurations, never a profile from the request.
The application must supply the durable PostgreSQL store, not a test double.

`POST /v1/franchise/appointments/{id}/whatsapp-confirmation` requires a Bearer
token accepted by the existing verifier, `appointment:manage`, the configured
tenant and the allowed organization. Cookie-only requests are not accepted.
The request uses exactly the twelve JSON fields in
`AppointmentNotificationRequest`: organization_id, confirmation_event_id,
appointment_version, binding_version, consent_id, consent_evidence_sha256,
evidence_sha256, expires_at, recipient, template_name, language_code and
body_parameters. Every field is required. Unknown/duplicate/case-aliased keys,
nulls, trailing JSON and bodies over 64 KiB fail closed. The caller approves
the exact content and preserves the request; no business wording is generated.

Tenant/actor/channel/purpose/policy/profile and event-derived DeliveryKey come
from verified identity/configuration, not the body. The endpoint approves first,
then calls ONLY the existing durable Channel; Sender revalidates immediately
before attempting the provider. All endpoint errors use non-sensitive Problem
Details and no-store. Its 45-second budget is local, not a provider SLA; configure
the existing server's finite read/write/proxy timeouts for this operation.

200 `accepted` means the provider accepted (or the fence already holds that
result), NOT delivered. A lost HTTP response can be retried with the identical
request/event while its grant remains current. A divergent/expired grant returns
409; do not change expiry/event/recipient to force replay. In-progress, terminal
and uncertain delivery return separate 409 codes; reconciliation remains owned
by the existing delivery system. No automatic HTTP retry is installed.

The endpoint does not enqueue work for eventual delivery. A crash between grant
and send needs the same explicit request to resume, not an assertion that a job
was queued. Local outcome lookup is available above; UI, opt-out capture, provider inbox/status,
real IdP/login/Meta, template purpose, secret store, rate/cost/egress, telemetry,
support and deployment remain target gates. Never expose this route publicly
before those applicable conditions are demonstrated.

This pack includes four exact official Meta reference files plus an explicitly adapted Python integration. It does not use the archived official Node SDK and never calls the adaptation “official unmodified code.” The license permits use only with Facebook web services/APIs and Platform Policy.

## Durable appointment approval (0.4.0)

In the composed franchise Go profile, apply migrations 0001 through 0050.
Construct `NewPostgresAppointmentApprovals(pool, contactHMACKey, purpose, policy)`
with the existing contact-identity HMAC key and server-approved consent purpose
and policy. Pass a Principal from the existing identity verifier to `Approve`,
never identity supplied by a request body. Use the exact key returned by
`AppointmentConfirmationDeliveryKey(tenant, confirmationEventID)`.

The single approval INSERT requires a committed confirmation transition/outbox,
current confirmed appointment version and organization, active verified contact
binding, PII permission, and the exact granted consent evidence. A later or
equal-time ambiguous consent decision rejects approval. No recipient or consent
is inferred from template text. The immutable grant stores hashes, not message
text or raw phone; hashes remain sensitive pseudonymous evidence.

Supply this store as the Sender's ApprovalResolver. Resolution rechecks the
appointment, contact version, policy, current unambiguous consent and expiry
before each attempted send. The 24-hour maximum approval window is a LOCAL
safety budget, not a Meta rule. This is a pre-send snapshot: it cannot recall
an HTTP request already in flight. Normal outbox cleanup after approval is
allowed; the immutable confirmation transition remains referenced.

This adds no second consent, appointment, contact, queue or delivery owner.
The approval table is a provider-specific authorization record. V267 adds the
explicit authenticated API above. Operator UI, outbox worker wiring, approved
template content, actual consent capture and
real Meta account/status/reconciliation remain target integration work.

## Go delivery-fence bridge (introduced in 0.3.0)

The composed Go profile also provides `internal/whatsappbridge.Sender`, which
implements the existing `outbounddelivery.Sender` contract. This bridge is
AUTHORED, not Meta Go source. Python continues to own the template POST and
provider evidence; PostgreSQL continues to own the send fence. No new database,
message queue, framework or paid dependency is added by the bridge itself.

Construct one sender per authorized tenant/profile. Supply an ApprovalResolver
backed by the existing authorized business workflow and a TokenSource backed
by the secret store. Approval binds the exact `outbounddelivery.MessageSHA256`,
profile bytes SHA-256, evidence/actor and validity window. The message Text must
be an approved template-request JSON; its recipient must equal ExternalID.
Never resolve approval or recipient from LLM text or fuzzy contact matching.
The durable resolver above can supply approval; TokenSource remains supplied
by the target secret store. Neither interface manufactures user consent.

Wrap this sender in `outbounddelivery.Channel` with the existing PostgreSQL
OutboundDeliveryStore and a separately verified receiver. Register only the
durable channel. No direct calls to Sender/CLI are allowed in the app: the CLI
does not independently deduplicate. Do not invent a new DeliveryKey or delete
an evidence directory to retry an uncertain operation.

`Process` requires absolute runtime/adapter/evidence paths and approved SHA-256
for the Python executable and materialized whatsapp_cloud.py. Use the exact
canonical pack hash, not a digest blindly accepted from arbitrary code. Deploy
the complete verified Python distribution and library in a protected read-only
tree: executable hash alone is not proof of every Python DLL/stdlib dependency
and path hashing is not protection against an attacker with write access.
The parent launches `-I -B` without shell, inherited secrets, PYTHONPATH or proxy
variables. The single stdin frame holds the token in memory, not CLI arguments
or files; stderr is discarded by Go and bounded stdout must match the request
binding. Timeout is finite. Any uncertain result remains unknown in the fence.

The evidence root contains sensitive provider response/receipt data and needs
restricted ACLs, retention, capacity and backup policy. Its stable subdirectory
is derived from tenant+DeliveryKey. The returned hash binds the stored receipt;
the receipt binds the provider response. An accepted receipt is NOT delivered.

Verify in the composed Go profile (not the Python-only profile):

```powershell
$env:ELITE_WHATSAPP_PYTHON = '[approved absolute Python executable]'
go test ./internal/whatsappbridge ./internal/outbounddelivery
```

For the durable integration test, provision the composed migrations into a
dedicated disposable database named `elite_whatsapp_*` on 127.0.0.1:55959, and set
`ELITE_WHATSAPP_TEST_DATABASE_URL`. The test retains synthetic immutable events;
dispose of that dedicated database afterward, never disable audit triggers.
Tests use a clearly synthetic provider child; Python adapter tests separately
exercise the real send owner with an injected transport. Neither is a live Meta
probe. Missing local test runtime/database is an explicit skip, not a PASS.

Remaining gates: deployed API/operator UI/worker wiring, real account and
template, provider protocol, real consent/opt-out, cost, TLS/egress, webhook
ingestion, status reconciliation and operation. SQL-fixture approval tests
are not a browser-to-Meta end-to-end journey.

The profile starts blocked. Before any call, prove/accept the official source license and current platform terms, Business/app/phone setup, approved template and recipient consent, webhook subscription, retention, quota/cost and reconciliation. Choose the current Graph API version from official Meta authority at project time; the pack does not guess a moving version. Tokens and app/verify secrets belong only in the named environment variables.

Run offline tests:

```powershell
Push-Location .\whatsapp_cloud
python -m unittest -v test_whatsapp_cloud.py
Pop-Location
```

Tests prove exact embedded upstream hashes, fail-closed template policy, official URL/header shape, atomic provider failure, subscription challenge, raw-body HMAC and scoped/pseudonymized normalized webhook evidence. They do not prove a Meta account, template approval, delivery, pricing, user consent, opt-out, throughput or reconciliation. Production ingress must connect normalized events to the common durable provider inbox and return within Meta's required webhook behavior after the target contract is verified.

## Upgrade from 0.1.0 to 0.2.0

Configure the exact approved `business_account_id` as well as `phone_number_id`.
Call `webhook_to_evidence(raw, signature, app_secret, new_output, profile=approved_profile)`.
The profile comes from server-side configuration, never from the webhook or LLM.
The keyword argument is mandatory: older calls fail rather than discard scope.
Consume `elite-whatsapp-cloud-events/v2` and webhook receipt `/v2`; do not coerce
old unscoped evidence into v2. No database schema or upstream pin changes.

The existing provider inbox must bind the verified WABA/phone to its tenant,
deduplicate by scoped event key, check the outbound message/recipient association
and retain out-of-order statuses. A valid app signature alone does not authorize
a tenant. No status is synthesized from the send HTTP response or arrival order.
`sent` is not `delivered`, and `delivered` is not a sale or appointment acceptance.

Malformed/foreign/oversized batches publish no evidence. Route them to authorized
incident handling with safe hashes and bounded retry/DLQ policy; do not acknowledge
them as processed or blindly retry a provider send. Ingress endpoint ACK behavior
and durable inbox integration remain project gates. Do not delete or reuse an
outbound evidence path after an uncertain send: this adapter has no durable send
fence of its own and must be wrapped by the existing shared delivery owner.

SHA-256 contact identifiers are pseudonymous, not anonymous. Restrict/retain all
evidence, especially raw outbound `provider-response.json`; the offline writer
does not retain a raw webhook body. The 0.9.0 opt-in HTTP lane above does retain
raw bytes and requires its separate proven access/retention controls. Neither
lane activates a provider call, account or paid service by materialization.

Rollback: preserve v2 evidence and pause its consumer if deployment fails. Do not
reactivate the older unscoped path. Rebuild the corrected version and prove the
target's scope binding, inbox, receipts, recovery and account before resuming.
````

### FILE: `whatsapp_cloud/upstream/LICENSE`
```yaml
block_id: "PY-META-WHATSAPP:upstream-license:v1"
operation: CREATE
provenance: ADAPTED
source: "https://github.com/fbsamples/whatsapp-api-examples/blob/de70ee908a67026e642aaee3703d20464e2a9466/LICENSE; legal text unchanged, final newline added by Markdown materialization and both hashes retained"
license: "LicenseRef-Meta-Platform-API-Only"
sha256: "48d97b3c936203a750a3288c2b924327769d62327383a1060244b3d3c05ad01f"
variables: []
secrets_allowed: false
```
````text
Copyright (c) Meta Platforms, Inc. and affiliates.
All rights reserved.

You are hereby granted a non-exclusive, worldwide, royalty-free license to use,
copy, modify, and distribute this software in source code or binary form for use
in connection with the web services and APIs provided by Facebook.

As with any software that integrates with the Facebook platform, your use of
this software is subject to the Facebook Platform Policy
[http://developers.facebook.com/policy/]. This copyright notice shall be
included in all copies or substantial portions of the software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
````

### FILE: `whatsapp_cloud/upstream/signature_validation_app.py`
```yaml
block_id: "PY-META-WHATSAPP:upstream-signature-example:v1"
operation: CREATE
provenance: VERBATIM
source: "https://github.com/fbsamples/whatsapp-api-examples/blob/de70ee908a67026e642aaee3703d20464e2a9466/signature-validation-with-webhooks-payloads/app.py; reference-only due retained control-flow defect"
license: "LicenseRef-Meta-Platform-API-Only"
sha256: "6c052c136a0be059faee9e9ad93ff056b05b7086b819fce067ebd27c8ef2412c"
variables: []
secrets_allowed: false
```
````python
"""
Copyright (c) Meta Platforms, Inc. and affiliates.
All rights reserved.

This source code is licensed under the license found in the
LICENSE file in the root directory of this source tree.
"""

from flask import Flask, request
import json
import hmac
import hashlib
import requests

app = Flask(__name__)
TOKEN = "your-secret-token"
PAGE_ACCESS_TOKEN = "secret_page_access_token"

@app.route('/webhook', methods=['GET', 'POST'])
def webhook():
    # Webhook verification
    if request.method == 'GET':
        if request.args.get("hub.mode") == "subscribe" and request.args.get("hub.challenge"):
            if not request.args.get("hub.verify_token") == TOKEN:
                return "Verification token mismatch", 403
            print("WEBHOOK_VERIFIED")
            return request.args["hub.challenge"], 200
    elif request.method == 'POST':
        # Validate payload
        signature = request.headers["X-Hub-Signature-256"].split('=')[1]
        payload = request.get_data()
        expected_signature = hmac.new(TOKEN.encode('utf-8'), payload, hashlib.sha256).hexdigest()

        if signature != expected_signature:
            print("Signature hash does not match")
        return 'INVALID SIGNATURE HASH', 403

        body = json.loads(payload.decode('utf-8'))

        if 'object' in body and body['object'] == 'page':
            entries = body['entry']
            # Iterate through each entry as multiple entries can sometimes be batched
            for entry in entries:
                # Fetch the 'changes' element field
                change_event = entry['changes'][0]
                # Verify it is a change in the 'feed' field
                if change_event['field'] != 'feed':
                    continue
                # Fetch the 'post_id' in the 'value' element
                post_id = change_event['value']['post_id']
                comment_on_post(post_id)

            return 'WEBHOOK EVENT HANDLED', 200
        return 'INVALID WEBHOOK EVENT', 403

def comment_on_post(post_id):
    payload = {
        'message': 'Lovely post!'
    }
    headers = {
        'content-type': 'application/json'
    }
    url = "https://graph.facebook.com/{}/comments?access_token={}".format(post_id, PAGE_ACCESS_TOKEN)

    r = requests.post(url, json=payload, headers=headers)
    print(r.text)

if __name__ == '__main__':
    app.run(debug=True)
````

### FILE: `whatsapp_cloud/upstream/message_helper.py`
```yaml
block_id: "PY-META-WHATSAPP:upstream-message-helper:v1"
operation: CREATE
provenance: VERBATIM
source: "https://github.com/fbsamples/whatsapp-api-examples/blob/de70ee908a67026e642aaee3703d20464e2a9466/send-messages-flight-app-python/message_helper.py; reference-only"
license: "LicenseRef-Meta-Platform-API-Only"
sha256: "fef1aa035aa2fad4769a08edbe4a347a87dac475394a821bc71b301d1ef45702"
variables: []
secrets_allowed: false
```
````python
"""
Copyright (c) Meta Platforms, Inc. and affiliates.
All rights reserved.

This source code is licensed under the license found in the
LICENSE file in the root directory of this source tree.
"""

import aiohttp
import json
from flask import current_app

async def send_message(data):
  headers = {
    "Content-type": "application/json",
    "Authorization": f"Bearer {current_app.config['ACCESS_TOKEN']}",
    }
  
  async with aiohttp.ClientSession() as session:
    url = 'https://graph.facebook.com' + f"/{current_app.config['VERSION']}/{current_app.config['PHONE_NUMBER_ID']}/messages"
    try:
      async with session.post(url, data=data, headers=headers) as response:
        if response.status == 200:
          print("Status:", response.status)
          print("Content-type:", response.headers['content-type'])

          html = await response.text()
          print("Body:", html)
        else:
          print(response.status)        
          print(response)        
    except aiohttp.ClientConnectorError as e:
      print('Connection Error', str(e))

def get_text_message_input(recipient, text):
  return json.dumps({
    "messaging_product": "whatsapp",
    "preview_url": False,
    "recipient_type": "individual",
    "to": recipient,
    "type": "text",
    "text": {
        "body": text
    }
  })

def get_templated_message_input(recipient, flight):
  return json.dumps({
    "messaging_product": "whatsapp",
    "to": recipient,
    "type": "template",
    "template": {
      "name": "sample_flight_confirmation",
      "language": {
        "code": "en_US"
      },
      "components": [
        {
          "type": "header",
          "parameters": [
            {
              "type": "document",
              "document": {
                "filename": "FlightConfirmation.pdf",
                "link": flight['document']
              }
            }
          ]
        },
        {
          "type": "body",
          "parameters": [
            {
              "type": "text",
              "text": flight['origin']
            },
            {
              "type": "text",
              "text": flight['destination']
            },
            {
              "type": "text",
              "text": flight['time']
            }
          ]
        }
      ]
    }
  })
````

### FILE: `whatsapp_cloud/upstream/incomingWebhook.js`
```yaml
block_id: "PY-META-WHATSAPP:upstream-webhook-example:v1"
operation: CREATE
provenance: VERBATIM
source: "https://github.com/fbsamples/whatsapp-api-examples/blob/de70ee908a67026e642aaee3703d20464e2a9466/template-for-ecommerce-js/routes/incomingWebhook.js; reference-only"
license: "LicenseRef-Meta-Platform-API-Only"
sha256: "a92a62bff637b8d5b0a3fbd3ce82dbd08c67904edff8df9adeb1be0e2d92330f"
variables: []
secrets_allowed: false
```
````javascript
/* Copyright (c) Meta Platforms, Inc. and affiliates.
* All rights reserved.
*
* This source code is licensed under the license found in the
* LICENSE file in the root directory of this source tree.
*/

const express = require('express');
const router = express.Router();
const { messageStatuses } = require("../public/javascripts/messageStatuses");
const { interactiveList, interactiveReplyButton } = require("../public/javascripts/interactiveMessages");
const { products } = require("../public/javascripts/products");
const { createProductsList, updateWhatsAppMessage, sendWhatsAppMessage } = require("../messageHelper")
const XHubSignature = require('x-hub-signature');

const verificationToken = process.env.WEBHOOK_VERIFICATION_TOKEN;
const appSecret = process.env.APP_SECRET;
const xhub = new XHubSignature('SHA256', appSecret);


async function processMessage(message) {
  const customerPhoneNumber = message.from;
  const messageType = message.type;

  if (messageType === "text") {
    const textMessage = message.text.body;
    console.log( textMessage );

  try {
    let replyButtonMessage = interactiveReplyButton;
    replyButtonMessage.to = process.env.RECIPIENT_PHONE_NUMBER;
    const replyButtonSent = await sendWhatsAppMessage(replyButtonMessage);
    console.log(replyButtonSent);
  } catch (error) {
    console.log(error);
  }

  } else if (messageType === "interactive") {
    const interactiveType = message.interactive.type;

    if (interactiveType === "button_reply") {
      const buttonId = message.interactive.button_reply.id;
      const buttonTitle = message.interactive.button_reply.title;

      if (buttonId == 1) {
        try {
          let productsList = interactiveList;
          productsList.to = process.env.RECIPIENT_PHONE_NUMBER;
          productsList.interactive.action.sections[0].rows = products.map(createProductsList);

          // List messages have a 10 item limit total
          productsList.interactive.action.sections[0].rows.length = 10;
          const sendProductLists = await sendWhatsAppMessage(productsList);
          console.log(sendProductLists);
        } catch (error) {
          console.log(error);
        }
      }
    }
    else if (interactiveType === "list_reply") {
      const itemId = message.interactive.list_reply.id;
      const itemTitle = message.interactive.list_reply.title;
      const itemDescrption = message.interactive.list_reply.description;
    }
  }
}

router.post('/', async function (req, res, next) {
  // Calculate x-hub signature value to check with value in request header
  const calcXHubSignature = xhub.sign(req.rawBody).toLowerCase();

  if(req.headers['x-hub-signature-256'] != calcXHubSignature)
  {
    console.log(
      "Warning - request header X-Hub-Signature not present or invalid"
    );
    res.sendStatus(401);
    return;
  }

  console.log("request header X-Hub-Signature validated");

  const body = req.body.entry[0].changes[0];

  // Verify this is from the messages webhook, not other updates
  if(body.field !== 'messages'){
    // not from the messages webhook so dont process
    return res.sendStatus(400)
  }

  if(body.value.hasOwnProperty("messages")) {

    // Mark an incoming message as read
    try {
      let sendReadStatus = messageStatuses.read;
      sendReadStatus.message_id = body.value.messages[0].id;
      const readSent = await sendWhatsAppMessage(sendReadStatus);
      console.log(readSent);
    } catch (error) {
      console.log(error);
    }

    body.value.messages.forEach(processMessage);
  }

  res.sendStatus( 200 );

});

router.get('/', function (req, res, next) {
  if (
    req.query['hub.mode'] == 'subscribe' &&
    req.query['hub.verify_token'] == verificationToken
  ) {
    res.send(req.query['hub.challenge']);
  } else {
    res.sendStatus(400);
  }
});

module.exports = router;
````

### FILE: `internal/whatsappbridge/sender.go`
```yaml
block_id: "PY-META-WHATSAPP:internal/whatsappbridge/sender.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed process integration with the existing outbounddelivery fence; no Meta Go source copied"
license: "LicenseRef-Workspace-Owner"
sha256: "107b0b5a8de7e032e9dcf708a4efc87c72dffee928b7fa18a87de208de95ded3"
variables: []
secrets_allowed: false
```
````go
// Package whatsappbridge connects the admitted Python adapter to the existing
// outbound delivery fence. This is AUTHORED integration, not Meta source code.
package whatsappbridge

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"elite.local/enterprise/internal/channels"
	"elite.local/enterprise/internal/outbounddelivery"
)

var ErrBridge = errors.New("whatsappbridge: unverified or uncertain provider result")

// Approval must be resolved from an authorized durable workflow, not message
// text. The bridge never generates consent, templates or business approvals.
type Approval struct {
	MessageSHA256  string
	ProfileSHA256  string
	EvidenceSHA256 string
	ApprovedBy     string
	NotBefore      time.Time
	ExpiresAt      time.Time
}

type ApprovalResolver interface {
	ResolveWhatsAppApproval(context.Context, string, string) (Approval, error)
}

type TokenSource interface {
	WhatsAppToken(context.Context) (string, error)
}

type Process struct {
	PythonExecutable  string
	PythonSHA256      string
	AdapterDirectory  string
	AdapterSHA256     string
	EvidenceDirectory string
}

// Sender is tenant-bound at composition time. Always wrap it in
// outbounddelivery.Channel with the shared PostgreSQL Store; never invoke a
// provider retry from an error. Accepted receipts do not assert delivery.
type Sender struct {
	TenantID  string
	Profile   json.RawMessage
	Approvals ApprovalResolver
	Tokens    TokenSource
	Process   Process
}

type frame struct {
	Schema          string          `json:"schema"`
	BindingSHA256   string          `json:"binding_sha256"`
	Profile         json.RawMessage `json:"profile"`
	Request         json.RawMessage `json:"request"`
	AccessToken     string          `json:"access_token"`
	OutputDirectory string          `json:"output_directory"`
}

type result struct {
	Schema            string    `json:"schema"`
	BindingSHA256     string    `json:"binding_sha256"`
	ProviderMessageID string    `json:"provider_message_id"`
	EvidenceSHA256    string    `json:"evidence_sha256"`
	AcceptedAt        time.Time `json:"accepted_at"`
}

func digest(value []byte) string { sum := sha256.Sum256(value); return hex.EncodeToString(sum[:]) }
func validDigest(value string) bool {
	b, e := hex.DecodeString(value)
	return e == nil && len(b) == 32 && value == strings.ToLower(value)
}

func exactFile(path, hash string) error {
	if !filepath.IsAbs(path) || !validDigest(hash) {
		return ErrBridge
	}
	info, err := os.Lstat(path)
	if err != nil || !info.Mode().IsRegular() {
		return ErrBridge
	}
	file, err := os.Open(path)
	if err != nil {
		return ErrBridge
	}
	defer file.Close()
	h := sha256.New()
	if _, err = io.Copy(h, file); err != nil || hex.EncodeToString(h.Sum(nil)) != hash {
		return ErrBridge
	}
	return nil
}

// Decode only the exact template shape; the Python owner validates policy and
// arity. The recipient must be the same one fenced in channels.Message.
func requestFor(message channels.Message) (json.RawMessage, error) {
	if len(message.Text) > 32768 {
		return nil, ErrBridge
	}
	var shape map[string]json.RawMessage
	if json.Unmarshal([]byte(message.Text), &shape) == nil {
		if _, ok := shape["kind"]; ok {
			return replyRequestFor(message)
		}
	}
	var body struct {
		Recipient      string   `json:"recipient"`
		TemplateName   string   `json:"template_name"`
		LanguageCode   string   `json:"language_code"`
		BodyParameters []string `json:"body_parameters"`
	}
	decoder := json.NewDecoder(strings.NewReader(message.Text))
	decoder.DisallowUnknownFields()
	if decoder.Decode(&body) != nil || decoder.Decode(new(any)) != io.EOF || body.Recipient != message.ExternalID || body.TemplateName == "" || body.LanguageCode == "" || body.BodyParameters == nil {
		return nil, ErrBridge
	}
	raw, err := json.Marshal(body)
	if err != nil {
		return nil, ErrBridge
	}
	return raw, nil
}

func (s *Sender) SendWithReceipt(ctx context.Context, message channels.Message) (outbounddelivery.Receipt, error) {
	empty := outbounddelivery.Receipt{}
	if s == nil || s.TenantID == "" || message.TenantID != s.TenantID || message.ChannelCode != "whatsapp" || message.Direction != channels.DirectionOut || s.Approvals == nil || s.Tokens == nil || len(s.Profile) > 32768 || !json.Valid(s.Profile) {
		return empty, ErrBridge
	}
	binding, err := outbounddelivery.MessageSHA256(message)
	if err != nil {
		return empty, ErrBridge
	}
	request, err := requestFor(message)
	if err != nil {
		return empty, err
	}
	approval, err := s.Approvals.ResolveWhatsAppApproval(ctx, message.TenantID, message.DeliveryKey)
	now := time.Now().UTC()
	if err != nil || approval.MessageSHA256 != binding || approval.ProfileSHA256 != digest(s.Profile) || !validDigest(approval.EvidenceSHA256) || strings.TrimSpace(approval.ApprovedBy) == "" || approval.NotBefore.IsZero() || now.Before(approval.NotBefore) || !now.Before(approval.ExpiresAt) {
		return empty, ErrBridge
	}
	p := s.Process
	script := filepath.Join(p.AdapterDirectory, "whatsapp_cloud.py")
	if !filepath.IsAbs(p.AdapterDirectory) || !filepath.IsAbs(p.EvidenceDirectory) || exactFile(p.PythonExecutable, p.PythonSHA256) != nil || exactFile(script, p.AdapterSHA256) != nil {
		return empty, ErrBridge
	}
	token, err := s.Tokens.WhatsAppToken(ctx)
	if err != nil || token == "" || len(token) > 16384 {
		return empty, ErrBridge
	}
	// Delivery identity determines one evidence destination. It does not replace
	// the database fence and must never be changed merely to force another send.
	output := filepath.Join(p.EvidenceDirectory, digest([]byte(message.TenantID+"\x00"+message.DeliveryKey)))
	input, err := json.Marshal(frame{"elite-whatsapp-send-bridge/v1", binding, s.Profile, request, token, output})
	if err != nil {
		return empty, ErrBridge
	}
	limitedCtx, cancel := context.WithTimeout(ctx, 40*time.Second)
	defer cancel()
	command := exec.CommandContext(limitedCtx, p.PythonExecutable, "-I", "-B", script, "--send-bridge")
	// No shell and no inherited API tokens, PYTHONPATH, proxies or startup hooks.
	command.Env = []string{}
	if root := os.Getenv("SystemRoot"); root != "" {
		command.Env = append(command.Env, "SystemRoot="+root)
	}
	command.Stdin = bytes.NewReader(input)
	var stdout boundedOutput
	command.Stdout = &stdout
	command.Stderr = io.Discard
	command.WaitDelay = 2 * time.Second
	if command.Run() != nil || stdout.overflow {
		return empty, ErrBridge
	}
	var response result
	decoder := json.NewDecoder(bytes.NewReader(stdout.Bytes()))
	decoder.DisallowUnknownFields()
	if decoder.Decode(&response) != nil || decoder.Decode(new(any)) != io.EOF || response.Schema != "elite-whatsapp-send-result/v1" || response.BindingSHA256 != binding {
		return empty, ErrBridge
	}
	receipt := outbounddelivery.Receipt{ProviderMessageID: response.ProviderMessageID, EvidenceSHA256: response.EvidenceSHA256, AcceptedAt: response.AcceptedAt}
	if receipt.Validate() != nil || receipt.AcceptedAt.Before(now.Add(-time.Minute)) || receipt.AcceptedAt.After(time.Now().Add(time.Minute)) {
		return empty, ErrBridge
	}
	return receipt, nil
}

type boundedOutput struct {
	bytes.Buffer
	overflow bool
}

func (b *boundedOutput) Write(p []byte) (int, error) {
	if b.Len()+len(p) > 4096 {
		b.overflow = true
		return 0, ErrBridge
	}
	return b.Buffer.Write(p)
}
````

### FILE: `internal/whatsappbridge/sender_test.go`
```yaml
block_id: "PY-META-WHATSAPP:internal/whatsappbridge/sender_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed process integration with the existing outbounddelivery fence; no Meta Go source copied"
license: "LicenseRef-Workspace-Owner"
sha256: "4f7d4b04e4a0e8811b14f3234873f4b0308a7c164f5609e8a315389ffb3e045d"
variables: []
secrets_allowed: false
```
````go
package whatsappbridge

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"elite.local/enterprise/internal/channels"
	"elite.local/enterprise/internal/leadstream"
	"elite.local/enterprise/internal/outbounddelivery"
	"elite.local/enterprise/internal/platform/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

type approvalFixture struct{ value Approval }

func (a approvalFixture) ResolveWhatsAppApproval(context.Context, string, string) (Approval, error) {
	return a.value, nil
}

type tokenFixture struct{}

func (tokenFixture) WhatsAppToken(context.Context) (string, error) { return "synthetic-token", nil }

type receiverFixture struct{}

func (receiverFixture) Code() string                                        { return "whatsapp" }
func (receiverFixture) Receive(context.Context) ([]channels.Message, error) { return nil, nil }
func (receiverFixture) Send(context.Context, channels.Message) error {
	return errors.New("receiver cannot send")
}

type storeFixture struct {
	hash, state string
	receipt     outbounddelivery.Receipt
}

func (s *storeFixture) Claim(_ context.Context, _ channels.Message, h string) (outbounddelivery.Claim, error) {
	if s.hash != "" && s.hash != h {
		return outbounddelivery.Claim{}, outbounddelivery.ErrConflict
	}
	switch s.state {
	case "accepted":
		return outbounddelivery.Claim{Replay: true}, nil
	case "unknown":
		return outbounddelivery.Claim{}, outbounddelivery.ErrUnknown
	case "sending":
		return outbounddelivery.Claim{}, outbounddelivery.ErrInProgress
	}
	s.hash = h
	s.state = "sending"
	return outbounddelivery.Claim{}, nil
}
func (s *storeFixture) Complete(_ context.Context, _ channels.Message, _ string, r outbounddelivery.Receipt) error {
	s.state = "accepted"
	s.receipt = r
	return nil
}
func (s *storeFixture) MarkUnknown(context.Context, channels.Message, string, string) error {
	s.state = "unknown"
	return nil
}
func (s *storeFixture) MarkFailed(context.Context, channels.Message, string, string, string) error {
	s.state = "failed_terminal"
	return nil
}

func messageFixture() channels.Message {
	return channels.Message{ChannelCode: "whatsapp", TenantID: "tenant-fixture", ExternalID: "5491112345678", Direction: channels.DirectionOut, DeliveryKey: "appointment-confirmed-fixture-key", Text: `{"recipient":"5491112345678","template_name":"appointment_confirmed","language_code":"es_AR","body_parameters":["2026-09-07T10:00Z"]}`}
}

func senderFixture(t *testing.T, mode string) (*Sender, channels.Message, string) {
	t.Helper()
	python := os.Getenv("ELITE_WHATSAPP_PYTHON")
	if python == "" {
		t.Skip("process contract requires explicit local ELITE_WHATSAPP_PYTHON; never a live provider")
	}
	pythonBytes, err := os.ReadFile(python)
	if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	marker := filepath.Join(root, "invocations.txt")
	markerJSON, _ := json.Marshal(marker)
	// A clearly synthetic child tests the real process boundary, not Meta.
	script := `import sys,json,datetime,pathlib
f=json.load(sys.stdin)
p=pathlib.Path(` + string(markerJSON) + `)
with p.open('ab') as stream: stream.write(b'call\n')
mode=` + `"` + mode + `"` + `
if mode=='lost': sys.stderr.write('synthetic-token private recipient'); sys.exit(2)
r={'schema':'elite-whatsapp-send-result/v1','binding_sha256':f['binding_sha256'],'provider_message_id':'wamid.fixture','evidence_sha256':'a'*64,'accepted_at':datetime.datetime.now(datetime.timezone.utc).isoformat()}
if mode=='wrong': r['binding_sha256']='f'*64
if mode=='missing': r['provider_message_id']=''
if mode=='overflow': print('x'*8192); sys.exit(0)
print(json.dumps(r))
`
	if err = os.WriteFile(filepath.Join(root, "whatsapp_cloud.py"), []byte(script), 0600); err != nil {
		t.Fatal(err)
	}
	message := messageFixture()
	h, err := outbounddelivery.MessageSHA256(message)
	if err != nil {
		t.Fatal(err)
	}
	profile := json.RawMessage(`{"synthetic_fixture":true}`)
	approval := Approval{MessageSHA256: h, ProfileSHA256: digest(profile), EvidenceSHA256: strings.Repeat("b", 64), ApprovedBy: "fixture-owner", NotBefore: time.Now().Add(-time.Minute), ExpiresAt: time.Now().Add(time.Minute)}
	s := &Sender{TenantID: message.TenantID, Profile: profile, Approvals: approvalFixture{approval}, Tokens: tokenFixture{}, Process: Process{PythonExecutable: python, PythonSHA256: digest(pythonBytes), AdapterDirectory: root, AdapterSHA256: digest([]byte(script)), EvidenceDirectory: filepath.Join(root, "evidence")}}
	return s, message, marker
}

func TestProcessBridgeFencedReplayAndAmbiguity(t *testing.T) {
	for _, mode := range []string{"accepted", "lost", "wrong", "missing", "overflow"} {
		t.Run(mode, func(t *testing.T) {
			sender, message, marker := senderFixture(t, mode)
			store := &storeFixture{}
			channel := &outbounddelivery.Channel{CodeValue: "whatsapp", Receiver: receiverFixture{}, Sender: sender, Store: store}
			for i := 0; i < 2; i++ {
				err := channel.Send(context.Background(), message)
				if mode == "accepted" && err != nil {
					t.Fatal(err)
				}
				if mode != "accepted" && !errors.Is(err, outbounddelivery.ErrUnknown) {
					t.Fatalf("expected unknown, got %v", err)
				}
				if err != nil && strings.Contains(err.Error(), "synthetic-token") {
					t.Fatal("secret in error")
				}
			}
			calls, err := os.ReadFile(marker)
			if err != nil || string(calls) != "call\n" {
				t.Fatalf("expected one process invocation: %q %v", calls, err)
			}
			if mode == "accepted" && store.state != "accepted" {
				t.Fatal(store.state)
			}
			if mode != "accepted" && store.state != "unknown" {
				t.Fatal(store.state)
			}
		})
	}
}

func TestBridgeRejectsUnapprovedOrDriftedRequestBeforeChild(t *testing.T) {
	for _, mode := range []string{"tenant", "recipient", "approval", "expired", "profile", "python", "adapter", "no_approvals", "channel"} {
		t.Run(mode, func(t *testing.T) {
			sender, message, marker := senderFixture(t, "accepted")
			switch mode {
			case "tenant":
				message.TenantID = "other"
			case "recipient":
				message.ExternalID = "5491199999999"
			case "approval":
				a := sender.Approvals.(approvalFixture)
				a.value.MessageSHA256 = strings.Repeat("c", 64)
				sender.Approvals = a
			case "expired":
				a := sender.Approvals.(approvalFixture)
				a.value.ExpiresAt = time.Now().Add(-time.Second)
				sender.Approvals = a
			case "profile":
				sender.Profile = json.RawMessage(`{"changed":true}`)
			case "python":
				sender.Process.PythonSHA256 = strings.Repeat("c", 64)
			case "adapter":
				sender.Process.AdapterSHA256 = strings.Repeat("c", 64)
			case "no_approvals":
				sender.Approvals = nil
			case "channel":
				message.ChannelCode = "email"
			}
			if _, err := sender.SendWithReceipt(context.Background(), message); !errors.Is(err, ErrBridge) {
				t.Fatal("unapproved call accepted", err)
			}
			if _, err := os.Stat(marker); !os.IsNotExist(err) {
				t.Fatal("child invoked")
			}
		})
	}
}

func TestRequestContractIsExact(t *testing.T) {
	for _, text := range []string{`{}`, messageFixture().Text + ` {}`, strings.Replace(messageFixture().Text, `"recipient":`, `"extra":1,"recipient":`, 1), strings.Repeat("x", 32769)} {
		message := messageFixture()
		message.Text = text
		if _, err := requestFor(message); err == nil {
			t.Fatal("malformed request accepted")
		}
	}
}

func TestProcessBridgePostgresDurableReplay(t *testing.T) {
	raw := os.Getenv("ELITE_WHATSAPP_TEST_DATABASE_URL")
	if raw == "" {
		t.Skip("requires dedicated local elite_whatsapp_ database")
	}
	endpoint, err := url.Parse(raw)
	if err != nil || endpoint.Hostname() != "127.0.0.1" || !strings.HasPrefix(endpoint.Path, "/elite_whatsapp_") {
		t.Fatal("unsafe test database target")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()
	pool, err := pgxpool.New(ctx, raw)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	for _, mode := range []string{"accepted", "lost", "wrong", "missing", "overflow"} {
		t.Run(mode, func(t *testing.T) {
			sender, message, marker := senderFixture(t, mode)
			tenant := leadstream.StableUUID(t.Name(), time.Now().UTC().Format(time.RFC3339Nano))
			_, err := pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,$2,'Synthetic bridge','Synthetic bridge')`, tenant, "wa-"+tenant[:8])
			if err != nil {
				t.Fatal(err)
			}
			sender.TenantID = tenant
			message.TenantID = tenant
			binding, err := outbounddelivery.MessageSHA256(message)
			if err != nil {
				t.Fatal(err)
			}
			approved := sender.Approvals.(approvalFixture)
			approved.value.MessageSHA256 = binding
			sender.Approvals = approved
			for attempt := 0; attempt < 2; attempt++ {
				// A new Store instance demonstrates that replay/unknown is in PostgreSQL,
				// not in the Go fixture's memory.
				store, err := postgres.NewOutboundDeliveryStore(pool, []byte("0123456789abcdef0123456789abcdef"), time.Second)
				if err != nil {
					t.Fatal(err)
				}
				channel := &outbounddelivery.Channel{CodeValue: "whatsapp", Receiver: receiverFixture{}, Sender: sender, Store: store}
				err = channel.Send(ctx, message)
				if mode == "accepted" && err != nil {
					t.Fatal(err)
				}
				if mode != "accepted" && !errors.Is(err, outbounddelivery.ErrUnknown) {
					t.Fatalf("expected durable unknown: %v", err)
				}
				changed := message
				changed.Text += " "
				if err = channel.Send(ctx, changed); !errors.Is(err, outbounddelivery.ErrConflict) {
					t.Fatalf("divergent request: %v", err)
				}
			}
			calls, err := os.ReadFile(marker)
			if err != nil || string(calls) != "call\n" {
				t.Fatalf("provider invocation count: %q %v", calls, err)
			}
			var state string
			var attempts, events int
			err = pool.QueryRow(ctx, `select state,attempt_count,(select count(*) from communication.outbound_delivery_event e where e.tenant_id=d.tenant_id and e.channel_code=d.channel_code and e.delivery_key=d.delivery_key) from communication.outbound_delivery d where tenant_id=$1 and channel_code='whatsapp' and delivery_key=$2`, tenant, message.DeliveryKey).Scan(&state, &attempts, &events)
			expected := "unknown"
			if mode == "accepted" {
				expected = "accepted"
			}
			if err != nil || state != expected || attempts != 1 || events != 2 {
				t.Fatalf("durable evidence state=%s attempts=%d events=%d err=%v", state, attempts, events, err)
			}
		})
	}
}
````

### FILE: `internal/whatsappbridge/appointment_approval.go`
```yaml
block_id: "PY-META-WHATSAPP:appointment-approval-0:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration of existing appointment, consent, contact and outbound owners; AWS outbox method"
license: "LicenseRef-Workspace-Owner"
sha256: "50f8f778be5cd5ced634037fe29c642bb3faf3f106d176e64239305e30644adc"
variables: []
secrets_allowed: false
```
````go
package whatsappbridge

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"elite.local/enterprise/internal/channels"
	"elite.local/enterprise/internal/contactidentity"
	"elite.local/enterprise/internal/outbounddelivery"
	"elite.local/enterprise/internal/platform/identity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrApproval = errors.New("whatsappbridge: appointment approval absent, divergent or no longer authorized")

// AppointmentApprovalCommand references proof from the authorized workflow.
// These hashes do not themselves prove consent or approve a Meta template.
// The caller must supply a verified Principal, never JSON-authored identity.
type AppointmentApprovalCommand struct {
	Message               channels.Message
	OrganizationID        string
	AppointmentID         string
	ConfirmationEventID   string
	AppointmentVersion    int64
	BindingVersion        int64
	PolicyVersion         string
	ConsentID             string
	ConsentPurpose        string
	ConsentEvidenceSHA256 string
	ProfileSHA256         string
	EvidenceSHA256        string
	ExpiresAt             time.Time
}

type PostgresAppointmentApprovals struct {
	pool    *pgxpool.Pool
	hmacKey []byte
	purpose string
	policy  string
}

func NewPostgresAppointmentApprovals(pool *pgxpool.Pool, key []byte, purpose, policy string) (*PostgresAppointmentApprovals, error) {
	if pool == nil || len(key) < 32 || strings.TrimSpace(purpose) == "" || len(purpose) > 128 || strings.TrimSpace(policy) == "" || len(policy) > 128 {
		return nil, ErrApproval
	}
	return &PostgresAppointmentApprovals{pool: pool, hmacKey: append([]byte(nil), key...), purpose: purpose, policy: policy}, nil
}

// One business event has one notification identity, regardless of recipient or
// retries. A different key cannot bypass the durable fence for the same event.
func AppointmentConfirmationDeliveryKey(tenant, event string) string {
	return "wa-appointment-confirmed:" + digest([]byte(tenant+"\x00"+event))
}

func (s *PostgresAppointmentApprovals) Approve(ctx context.Context, p identity.Principal, c AppointmentApprovalCommand) (Approval, error) {
	empty := Approval{}
	now := time.Now().UTC()
	if s == nil || c.ConsentPurpose != s.purpose || c.PolicyVersion != s.policy {
		return empty, ErrApproval
	}
	if c.ConsentID == "" || strings.TrimSpace(c.ConsentPurpose) == "" || len(c.ConsentPurpose) > 128 {
		return empty, ErrApproval
	}
	if s == nil || p.Subject == "" || len(p.Subject) > 255 || p.TenantID != c.Message.TenantID || !p.Allowed("appointment:manage") || !p.AllowedOrganization(c.OrganizationID) || c.Message.ChannelCode != "whatsapp" || c.AppointmentID == "" || c.ConfirmationEventID == "" || c.Message.DeliveryKey != AppointmentConfirmationDeliveryKey(p.TenantID, c.ConfirmationEventID) || c.AppointmentVersion < 1 || c.BindingVersion < 1 || strings.TrimSpace(c.PolicyVersion) == "" || len(c.PolicyVersion) > 128 || !validDigest(c.ConsentEvidenceSHA256) || !validDigest(c.ProfileSHA256) || !validDigest(c.EvidenceSHA256) || !c.ExpiresAt.After(now) || c.ExpiresAt.After(now.Add(24*time.Hour)) {
		return empty, ErrApproval
	}
	hash, err := outbounddelivery.MessageSHA256(c.Message)
	if err != nil {
		return empty, ErrApproval
	}
	if _, err = requestFor(c.Message); err != nil {
		return empty, ErrApproval
	}
	recipient, err := contactidentity.ExternalDigest(s.hmacKey, p.TenantID, "whatsapp", c.Message.ExternalID)
	if err != nil {
		return empty, ErrApproval
	}
	// The fingerprint is never persisted in clear: it binds the complete approved
	// request and actor without storing recipient/template body as database PII.
	fingerprintBytes, err := json.Marshal(struct {
		Actor   string
		Command AppointmentApprovalCommand
	}{p.Subject, c})
	if err != nil {
		return empty, ErrApproval
	}
	fingerprint := digest(fingerprintBytes)
	// A single INSERT..SELECT observes only committed appointment+outbox+binding.
	// It does not own or mutate those aggregates. No provider effect happens here.
	_, err = s.pool.Exec(ctx, `insert into communication.whatsapp_appointment_approval
 (tenant_id,delivery_key,appointment_id,appointment_version,organization_id,confirmation_event_id,external_id_hmac,binding_version,lead_id,subject_id,policy_version,consent_evidence_sha256,message_sha256,profile_sha256,approval_request_sha256,evidence_sha256,approved_by,expires_at,consent_id,consent_purpose)
 select a.tenant_id,$2,a.appointment_id,a.version,a.organization_id,t.transition_id,b.external_id_hmac,b.version,b.lead_id,b.subject_id,b.policy_version,$10,$11,$12,$13,$14,$15,$16,consent.consent_id,consent.purpose_code
 from crm.appointment a
 join crm.lead l on l.tenant_id=a.tenant_id and l.lead_id=a.lead_id and l.organization_id=a.organization_id
 join crm.appointment_transition t on t.tenant_id=a.tenant_id and t.appointment_id=a.appointment_id and t.transition_id=$5 and t.from_state='requested' and t.to_state='confirmed'
 join platform.outbox_event e on e.tenant_id=t.tenant_id and e.event_id=t.transition_id and e.aggregate_type='appointment' and e.aggregate_id=a.appointment_id and e.aggregate_version=a.version and e.event_type='appointment.confirmed' and e.schema_version=1
 join communication.contact_channel_binding b on b.tenant_id=a.tenant_id and b.lead_id=a.lead_id and b.channel_code='whatsapp' and b.external_id_hmac=$7
 join crm.consent_evidence consent on consent.tenant_id=a.tenant_id and consent.lead_id=a.lead_id and consent.consent_id=$17 and consent.purpose_code=$18 and consent.policy_version=b.policy_version and consent.decision='granted' and consent.evidence_sha256_hex=$10 and consent.occurred_at<=statement_timestamp()
 where a.tenant_id=$1 and a.organization_id=$3 and a.appointment_id=$4 and a.version=$6 and a.state='confirmed' and a.starts_at>statement_timestamp()
 and b.state='active' and b.pii_allowed and b.version=$8 and b.policy_version=$9 and b.effective_at<=statement_timestamp()
 and not exists(select 1 from crm.consent_evidence newer where newer.tenant_id=consent.tenant_id and newer.lead_id=consent.lead_id and newer.purpose_code=consent.purpose_code and newer.consent_id<>consent.consent_id and newer.occurred_at>=consent.occurred_at and newer.occurred_at<=statement_timestamp())
 on conflict do nothing`, p.TenantID, c.Message.DeliveryKey, c.OrganizationID, c.AppointmentID, c.ConfirmationEventID, c.AppointmentVersion, recipient, c.BindingVersion, c.PolicyVersion, c.ConsentEvidenceSHA256, hash, c.ProfileSHA256, fingerprint, c.EvidenceSHA256, p.Subject, c.ExpiresAt.UTC(), c.ConsentID, c.ConsentPurpose)
	if err != nil {
		return empty, ErrApproval
	}
	var stored string
	if s.pool.QueryRow(ctx, `select approval_request_sha256 from communication.whatsapp_appointment_approval where tenant_id=$1 and delivery_key=$2`, p.TenantID, c.Message.DeliveryKey).Scan(&stored) != nil || stored != fingerprint {
		return empty, ErrApproval
	}
	return s.ResolveWhatsAppApproval(ctx, p.TenantID, c.Message.DeliveryKey)
}

// Resolve rechecks the current binding and appointment. It does not renew an
// approval on expiry, version change, cancellation, reassignment or revocation.
// It is a pre-send snapshot, not a promise to recall an already in-flight send.
func (s *PostgresAppointmentApprovals) ResolveWhatsAppApproval(ctx context.Context, tenant, key string) (Approval, error) {
	value := Approval{}
	if s == nil || strings.TrimSpace(tenant) == "" || strings.TrimSpace(key) == "" {
		return value, ErrApproval
	}
	err := s.pool.QueryRow(ctx, `select g.message_sha256,g.profile_sha256,g.evidence_sha256,g.approved_by,g.approved_at,g.expires_at
 from communication.whatsapp_appointment_approval g
 join crm.appointment a on a.tenant_id=g.tenant_id and a.appointment_id=g.appointment_id and a.organization_id=g.organization_id and a.lead_id=g.lead_id and a.version=g.appointment_version and a.state='confirmed' and a.starts_at>statement_timestamp()
 join crm.lead l on l.tenant_id=a.tenant_id and l.lead_id=a.lead_id and l.organization_id=a.organization_id
 join communication.contact_channel_binding b on b.tenant_id=g.tenant_id and b.channel_code=g.channel_code and b.external_id_hmac=g.external_id_hmac and b.version=g.binding_version and b.lead_id=g.lead_id and b.subject_id=g.subject_id and b.policy_version=g.policy_version and b.state='active' and b.pii_allowed and b.effective_at<=statement_timestamp()
 join crm.consent_evidence consent on consent.tenant_id=g.tenant_id and consent.lead_id=g.lead_id and consent.consent_id=g.consent_id and consent.purpose_code=g.consent_purpose and consent.policy_version=g.policy_version and consent.evidence_sha256_hex=g.consent_evidence_sha256 and consent.decision='granted' and consent.occurred_at<=statement_timestamp()
 where g.tenant_id=$1 and g.delivery_key=$2 and g.policy_version=$3 and g.consent_purpose=$4 and g.approved_at<=statement_timestamp() and g.expires_at>statement_timestamp()
 and not exists(select 1 from crm.consent_evidence newer where newer.tenant_id=consent.tenant_id and newer.lead_id=consent.lead_id and newer.purpose_code=consent.purpose_code and newer.consent_id<>consent.consent_id and newer.occurred_at>=consent.occurred_at and newer.occurred_at<=statement_timestamp())`, tenant, key, s.policy, s.purpose).Scan(&value.MessageSHA256, &value.ProfileSHA256, &value.EvidenceSHA256, &value.ApprovedBy, &value.NotBefore, &value.ExpiresAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Approval{}, ErrApproval
		}
		return Approval{}, ErrApproval
	}
	return value, nil
}
````

### FILE: `internal/whatsappbridge/appointment_approval_test.go`
```yaml
block_id: "PY-META-WHATSAPP:appointment-approval-1:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration of existing appointment, consent, contact and outbound owners; AWS outbox method"
license: "LicenseRef-Workspace-Owner"
sha256: "057c9fc0031b6566fe27bfe496f229282d69b88e44b2614c73990d2e720c34e2"
variables: []
secrets_allowed: false
```
````go
package whatsappbridge

import (
	"context"
	"errors"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"elite.local/enterprise/internal/contactidentity"
	"elite.local/enterprise/internal/leadstream"
	"elite.local/enterprise/internal/outbounddelivery"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

func approvalPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	raw := os.Getenv("ELITE_WHATSAPP_TEST_DATABASE_URL")
	if raw == "" {
		t.Skip("requires dedicated local elite_whatsapp_ database and 50 migrations")
	}
	// Same local/dedicated guard as the existing process integration test.
	if !strings.HasPrefix(raw, "postgres://postgres@127.0.0.1:55959/elite_whatsapp_") {
		t.Fatal("unsafe test database target")
	}
	pool, err := pgxpool.New(context.Background(), raw)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)
	return pool
}

func approvalFixtureData(t *testing.T, pool *pgxpool.Pool) (identity.Principal, AppointmentApprovalCommand, *PostgresAppointmentApprovals) {
	t.Helper()
	ctx := context.Background()
	tenant := leadstream.StableUUID(t.Name(), time.Now().Format(time.RFC3339Nano))
	event := leadstream.StableUUID(tenant, "confirmation")
	statements := []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1::uuid,'wa-'||$1::text,'Synthetic approval','Synthetic approval')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,'store-1','store-1','Store 1','store')`,
		`insert into crm.lead(tenant_id,lead_id,organization_id,lifecycle_state,source_code,contact_payload,consent_required) values($1,'lead-1','store-1','new','fixture','{}',true)`,
		`insert into crm.appointment(tenant_id,appointment_id,organization_id,lead_id,appointment_kind,starts_at,state,version) values($1,'appointment-1','store-1','lead-1','consultation',clock_timestamp()+interval '1 day','confirmed',3)`,
		`insert into crm.consent_evidence(tenant_id,consent_id,lead_id,purpose_code,policy_version,decision,occurred_at,evidence_sha256_hex) values($1,'consent-1','lead-1','fixture-appointment-notification','policy-1','granted',clock_timestamp()-interval '1 minute',repeat('c',64))`,
	}
	for _, q := range statements {
		if _, err := pool.Exec(ctx, q, tenant); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := pool.Exec(ctx, `insert into crm.appointment_transition(tenant_id,transition_id,appointment_id,from_state,to_state,actor_subject) values($1,$2,'appointment-1','requested','confirmed','fixture-operator')`, tenant, event); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload) values($1,$2,'appointment','appointment-1',3,'appointment.confirmed',1,clock_timestamp(),'{}')`, tenant, event); err != nil {
		t.Fatal(err)
	}
	key := []byte("0123456789abcdef0123456789abcdef")
	binding, err := postgres.NewContactIdentityStore(pool, key)
	if err != nil {
		t.Fatal(err)
	}
	m := messageFixture()
	m.TenantID = tenant
	m.DeliveryKey = AppointmentConfirmationDeliveryKey(tenant, event)
	_, err = binding.Apply(ctx, contactidentity.Command{TenantID: tenant, ChannelCode: "whatsapp", ExternalID: m.ExternalID, LeadID: "lead-1", SubjectID: "lead:lead-1", PIIAllowed: true, State: contactidentity.StateActive, PolicyVersion: "policy-1", EvidenceSHA256: strings.Repeat("e", 64), EffectiveAt: time.Now().Add(-time.Minute), RequestID: "binding-1"})
	if err != nil {
		t.Fatal(err)
	}
	resolver, err := NewPostgresAppointmentApprovals(pool, key, "fixture-appointment-notification", "policy-1")
	if err != nil {
		t.Fatal(err)
	}
	p := identity.Principal{Subject: "fixture-operator", TenantID: tenant, Permissions: map[string]struct{}{"appointment:manage": {}}, Organizations: map[string]struct{}{"store-1": {}}}
	c := AppointmentApprovalCommand{Message: m, OrganizationID: "store-1", AppointmentID: "appointment-1", ConfirmationEventID: event, AppointmentVersion: 3, BindingVersion: 1, PolicyVersion: "policy-1", ConsentID: "consent-1", ConsentPurpose: "fixture-appointment-notification", ConsentEvidenceSHA256: strings.Repeat("c", 64), ProfileSHA256: digest([]byte(`{"synthetic_fixture":true}`)), EvidenceSHA256: strings.Repeat("d", 64), ExpiresAt: time.Now().Add(time.Hour)}
	return p, c, resolver
}

func TestAppointmentApprovalDurableConcurrentAndFenced(t *testing.T) {
	pool := approvalPool(t)
	p, c, resolver := approvalFixtureData(t, pool)
	ctx := context.Background()
	var wg sync.WaitGroup
	errs := make(chan error, 8)
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); _, err := resolver.Approve(ctx, p, c); errs <- err }()
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatal(err)
		}
	}
	var count int
	if err := pool.QueryRow(ctx, `select count(*) from communication.whatsapp_appointment_approval where tenant_id=$1`, p.TenantID).Scan(&count); err != nil || count != 1 {
		t.Fatal(count, err)
	}
	if _, err := resolver.ResolveWhatsAppApproval(ctx, leadstream.StableUUID("other-tenant"), c.Message.DeliveryKey); !errors.Is(err, ErrApproval) {
		t.Fatal("foreign tenant approval", err)
	}
	changed := c
	changed.Message.Text += " "
	if _, err := resolver.Approve(ctx, p, changed); !errors.Is(err, ErrApproval) {
		t.Fatal("divergent replay accepted", err)
	}
	for _, mutation := range []string{`update communication.whatsapp_appointment_approval set approved_by='changed' where tenant_id=$1`, `delete from communication.whatsapp_appointment_approval where tenant_id=$1`} {
		if _, err := pool.Exec(ctx, mutation, p.TenantID); err == nil {
			t.Fatal("immutable evidence mutated")
		}
	}
	sender, _, marker := senderFixture(t, "accepted")
	sender.TenantID = p.TenantID
	sender.Approvals = resolver
	for i := 0; i < 2; i++ {
		store, err := postgres.NewOutboundDeliveryStore(pool, []byte("0123456789abcdef0123456789abcdef"), time.Second)
		if err != nil {
			t.Fatal(err)
		}
		channel := &outbounddelivery.Channel{CodeValue: "whatsapp", Receiver: receiverFixture{}, Sender: sender, Store: store}
		if err = channel.Send(ctx, c.Message); err != nil {
			t.Fatal(err)
		}
	}
	calls, err := os.ReadFile(marker)
	if err != nil || string(calls) != "call\n" {
		t.Fatal("duplicate child send", err)
	}
	var state string
	if err = pool.QueryRow(ctx, `select state from communication.outbound_delivery where tenant_id=$1 and delivery_key=$2`, p.TenantID, c.Message.DeliveryKey).Scan(&state); err != nil || state != "accepted" {
		t.Fatal(state, err)
	}
}

func TestAppointmentApprovalRequiresExactProof(t *testing.T) {
	pool := approvalPool(t)
	for _, mode := range []string{"role", "organization", "tenant", "recipient", "event", "appointment_version", "binding_version", "policy", "consent_id", "consent_hash", "consent_purpose", "profile_hash", "expired", "key", "no_outbox", "requested", "withdrawn", "tied_consent", "future_binding"} {
		t.Run(mode, func(t *testing.T) {
			p, c, resolver := approvalFixtureData(t, pool)
			ctx := context.Background()
			switch mode {
			case "role":
				p.Permissions = map[string]struct{}{}
			case "organization":
				p.Organizations = map[string]struct{}{}
			case "tenant":
				p.TenantID = leadstream.StableUUID("different")
			case "recipient":
				c.Message.ExternalID = "5491199999999"
			case "event":
				c.ConfirmationEventID = leadstream.StableUUID("other-event")
				c.Message.DeliveryKey = AppointmentConfirmationDeliveryKey(p.TenantID, c.ConfirmationEventID)
			case "appointment_version":
				c.AppointmentVersion = 2
			case "binding_version":
				c.BindingVersion = 2
			case "policy":
				c.PolicyVersion = "other-policy"
			case "consent_id":
				c.ConsentID = "other-consent"
			case "consent_hash":
				c.ConsentEvidenceSHA256 = strings.Repeat("a", 64)
			case "consent_purpose":
				c.ConsentPurpose = "other-purpose"
			case "profile_hash":
				c.ProfileSHA256 = ""
			case "expired":
				c.ExpiresAt = time.Now().Add(-time.Second)
			case "key":
				c.Message.DeliveryKey = strings.Repeat("k", 64)
			case "no_outbox":
				if _, err := pool.Exec(ctx, `delete from platform.outbox_event where tenant_id=$1`, p.TenantID); err != nil {
					t.Fatal(err)
				}
			case "requested":
				if _, err := pool.Exec(ctx, `update crm.appointment set state='requested' where tenant_id=$1`, p.TenantID); err != nil {
					t.Fatal(err)
				}
			case "withdrawn", "tied_consent":
				timing := "clock_timestamp()"
				if mode == "tied_consent" {
					timing = "occurred_at"
				}
				q := `insert into crm.consent_evidence(tenant_id,consent_id,lead_id,purpose_code,policy_version,decision,occurred_at,evidence_sha256_hex) select tenant_id,'newer',lead_id,purpose_code,policy_version,'withdrawn',` + timing + `,evidence_sha256_hex from crm.consent_evidence where tenant_id=$1`
				if _, err := pool.Exec(ctx, q, p.TenantID); err != nil {
					t.Fatal(err)
				}
			case "future_binding":
				if _, err := pool.Exec(ctx, `update communication.contact_channel_binding set effective_at=clock_timestamp()+interval '1 day' where tenant_id=$1`, p.TenantID); err != nil {
					t.Fatal(err)
				}
			}
			if _, err := resolver.Approve(ctx, p, c); !errors.Is(err, ErrApproval) {
				t.Fatalf("missing proof accepted: %v", err)
			}
			var count int
			if err := pool.QueryRow(ctx, `select count(*) from communication.whatsapp_appointment_approval where tenant_id=$1`, c.Message.TenantID).Scan(&count); err != nil || count != 0 {
				t.Fatal("unexpected durable approval", count, err)
			}
		})
	}
}

func TestAppointmentApprovalRechecksRevocationAndChanges(t *testing.T) {
	pool := approvalPool(t)
	for _, mode := range []string{"revoked", "pii", "version", "cancelled", "consent_withdrawn", "expiry"} {
		t.Run(mode, func(t *testing.T) {
			p, c, resolver := approvalFixtureData(t, pool)
			ctx := context.Background()
			if mode == "expiry" {
				c.ExpiresAt = time.Now().Add(300 * time.Millisecond)
			}
			if _, err := resolver.Approve(ctx, p, c); err != nil {
				t.Fatal(err)
			}
			queries := map[string]string{
				"revoked":           `update communication.contact_channel_binding set state='revoked',pii_allowed=false,version=version+1 where tenant_id=$1`,
				"pii":               `update communication.contact_channel_binding set pii_allowed=false,version=version+1 where tenant_id=$1`,
				"version":           `update crm.appointment set version=version+1 where tenant_id=$1`,
				"cancelled":         `update crm.appointment set state='cancelled',version=version+1 where tenant_id=$1`,
				"consent_withdrawn": `insert into crm.consent_evidence(tenant_id,consent_id,lead_id,purpose_code,policy_version,decision,occurred_at,evidence_sha256_hex) select tenant_id,'withdrawal',lead_id,purpose_code,policy_version,'withdrawn',clock_timestamp(),evidence_sha256_hex from crm.consent_evidence where tenant_id=$1`,
			}
			if mode == "expiry" {
				time.Sleep(350 * time.Millisecond)
			} else {
				if _, err := pool.Exec(ctx, queries[mode], p.TenantID); err != nil {
					t.Fatal(err)
				}
			}
			if _, err := resolver.ResolveWhatsAppApproval(ctx, p.TenantID, c.Message.DeliveryKey); !errors.Is(err, ErrApproval) {
				t.Fatal("stale approval returned", err)
			}
			sender, _, marker := senderFixture(t, "accepted")
			sender.TenantID = p.TenantID
			sender.Approvals = resolver
			if _, err := sender.SendWithReceipt(ctx, c.Message); !errors.Is(err, ErrBridge) {
				t.Fatal("stale approval sent", err)
			}
			if _, err := os.Stat(marker); !os.IsNotExist(err) {
				t.Fatal("child invoked after revocation")
			}
		})
	}
}
````

### FILE: `db/migrations/0050_whatsapp_appointment_approval.up.sql`
```yaml
block_id: "PY-META-WHATSAPP:appointment-approval-2:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration of existing appointment, consent, contact and outbound owners; AWS outbox method"
license: "LicenseRef-Workspace-Owner"
sha256: "e880a343999d9d79de15885a56ea6aab21776ee94ac53f489ce711f74418e8a8"
variables: []
secrets_allowed: false
```
````sql
begin;

create table communication.whatsapp_appointment_approval (
  tenant_id uuid not null,
  delivery_key text not null check(length(delivery_key) between 16 and 128),
  channel_code text not null default 'whatsapp' check(channel_code='whatsapp'),
  appointment_id text not null,
  appointment_version bigint not null check(appointment_version>0),
  organization_id text not null,
  confirmation_event_id uuid not null,
  external_id_hmac text not null check(external_id_hmac ~ '^[0-9a-f]{64}$'),
  binding_version bigint not null check(binding_version>0),
  lead_id text not null,
  subject_id text not null,
  policy_version text not null check(length(policy_version) between 1 and 128),
  consent_id text not null,
  consent_purpose text not null check(length(consent_purpose) between 1 and 128),
  consent_evidence_sha256 text not null check(consent_evidence_sha256 ~ '^[0-9a-f]{64}$'),
  message_sha256 text not null check(message_sha256 ~ '^[0-9a-f]{64}$'),
  profile_sha256 text not null check(profile_sha256 ~ '^[0-9a-f]{64}$'),
  approval_request_sha256 text not null check(approval_request_sha256 ~ '^[0-9a-f]{64}$'),
  evidence_sha256 text not null check(evidence_sha256 ~ '^[0-9a-f]{64}$'),
  approved_by text not null check(length(approved_by) between 1 and 255),
  approved_at timestamptz not null default statement_timestamp(),
  expires_at timestamptz not null,
  primary key(tenant_id,delivery_key),
  unique(tenant_id,confirmation_event_id,channel_code),
  foreign key(tenant_id,appointment_id) references crm.appointment(tenant_id,appointment_id),
  foreign key(tenant_id,confirmation_event_id) references crm.appointment_transition(tenant_id,transition_id),
  foreign key(tenant_id,consent_id) references crm.consent_evidence(tenant_id,consent_id),
  foreign key(tenant_id,channel_code,external_id_hmac) references communication.contact_channel_binding(tenant_id,channel_code,external_id_hmac),
  check(expires_at>approved_at)
);

create function communication.reject_whatsapp_approval_mutation() returns trigger
language plpgsql as $$ begin raise exception using errcode='55000',message='WhatsApp approval evidence is immutable'; end $$;
create trigger whatsapp_approval_immutable before update or delete
on communication.whatsapp_appointment_approval for each row
execute function communication.reject_whatsapp_approval_mutation();

create index whatsapp_notification_consent_lookup_idx
on crm.consent_evidence(tenant_id,lead_id,purpose_code,occurred_at desc);

commit;
````

### FILE: `db/migrations/0050_whatsapp_appointment_approval.down.sql`
```yaml
block_id: "PY-META-WHATSAPP:appointment-approval-3:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration of existing appointment, consent, contact and outbound owners; AWS outbox method"
license: "LicenseRef-Workspace-Owner"
sha256: "1aff4e7b751d1db0349074c0a9e36bcd55bd0e45b0fe5aa9da8d2fe9c91e2084"
variables: []
secrets_allowed: false
```
````sql
begin;
-- Destructive schema rollback. Export/retain evidence and stop the consumer
-- before using this on an approved target; not a retry mechanism.
drop trigger whatsapp_approval_immutable on communication.whatsapp_appointment_approval;
drop function communication.reject_whatsapp_approval_mutation();
drop table communication.whatsapp_appointment_approval;
drop index crm.whatsapp_notification_consent_lookup_idx;
commit;
````

### FILE: `db/tests/0050_whatsapp_appointment_approval.test.sql`
```yaml
block_id: "PY-META-WHATSAPP:appointment-approval-4:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration of existing appointment, consent, contact and outbound owners; AWS outbox method"
license: "LicenseRef-Workspace-Owner"
sha256: "db9e4dce1fb6c19b8d5590ed72f77cce5a2a012fdca53ebb62a8ce13f3a7fc41"
variables: []
secrets_allowed: false
```
````sql
begin;
do $$ begin
  if to_regclass('communication.whatsapp_appointment_approval') is null then
    raise exception 'approval table is required';
  end if;
  if not exists(select 1 from pg_trigger where tgrelid='communication.whatsapp_appointment_approval'::regclass and tgname='whatsapp_approval_immutable' and tgenabled='O') then
    raise exception 'approval immutability trigger is required';
  end if;
end $$;
rollback;
````


### FILE: `internal/whatsappbridge/appointment_notification.go`
```yaml
block_id: "PY-META-WHATSAPP:notification-http-0:v1"
operation: CREATE
provenance: AUTHORED
source: "local authenticated composition of admitted identity, approval and delivery owners"
license: "LicenseRef-Workspace-Owner"
sha256: "939cafb14e13a689c415ad40e8e3e6c25336d732591fd5b90d702d876af9f8bb"
variables: []
secrets_allowed: false
```
````go
package whatsappbridge

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net/http"
	"strings"
	"time"

	"elite.local/enterprise/internal/channels"
	"elite.local/enterprise/internal/outbounddelivery"
	"elite.local/enterprise/internal/platform/identity"
)

// AppointmentNotificationRequest is an explicit operator-approved request.
// Tenant, actor, channel, purpose, policy, profile and delivery identity are
// server-derived. This endpoint neither captures consent nor confirms a turn.
type AppointmentNotificationRequest struct {
	OrganizationID        string    `json:"organization_id"`
	ConfirmationEventID   string    `json:"confirmation_event_id"`
	AppointmentVersion    int64     `json:"appointment_version"`
	BindingVersion        int64     `json:"binding_version"`
	ConsentID             string    `json:"consent_id"`
	ConsentEvidenceSHA256 string    `json:"consent_evidence_sha256"`
	EvidenceSHA256        string    `json:"evidence_sha256"`
	ExpiresAt             time.Time `json:"expires_at"`
	Recipient             string    `json:"recipient"`
	TemplateName          string    `json:"template_name"`
	LanguageCode          string    `json:"language_code"`
	BodyParameters        []string  `json:"body_parameters"`
}

// AppointmentNotificationModule follows the existing module Register contract.
// Mount on the existing protected server; it creates no listener or retry loop.
type AppointmentNotificationModule struct {
	approvals *PostgresAppointmentApprovals
	sender    *Sender
	channel   *outbounddelivery.Channel
}

func NewAppointmentNotificationModule(approvals *PostgresAppointmentApprovals, sender *Sender, store outbounddelivery.Store, receiver channels.Channel) (*AppointmentNotificationModule, error) {
	if approvals == nil || approvals.pool == nil || sender == nil || sender.TenantID == "" || sender.Tokens == nil || !json.Valid(sender.Profile) || len(sender.Profile) > 32768 || store == nil || receiver == nil || receiver.Code() != "whatsapp" {
		return nil, ErrBridge
	}
	// Freeze the caller's struct/profile and force the same durable resolver into
	// the sender. The caller cannot accidentally wire a permissive second one.
	copySender := *sender
	copySender.Profile = append(json.RawMessage(nil), sender.Profile...)
	copySender.Approvals = approvals
	return &AppointmentNotificationModule{approvals: approvals, sender: &copySender, channel: &outbounddelivery.Channel{CodeValue: "whatsapp", Receiver: receiver, Sender: &copySender, Store: store}}, nil
}

func (m *AppointmentNotificationModule) Register(mux *http.ServeMux, verifier identity.Verifier) {
	m.registerNotificationStatus(mux, verifier)
	m.registerNotificationHistory(mux, verifier)
	mux.HandleFunc("POST /v1/franchise/appointments/{id}/whatsapp-confirmation", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 45*time.Second)
		defer cancel()
		p, ok := m.notificationPrincipal(ctx, w, r, verifier)
		if !ok {
			return
		}
		media, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
		if err != nil || media != "application/json" || len(params) > 1 || (len(params) == 1 && !strings.EqualFold(params["charset"], "utf-8")) {
			notificationProblem(w, 415, "UNSUPPORTED_MEDIA_TYPE")
			return
		}
		input, err := readNotificationRequest(http.MaxBytesReader(w, r.Body, 64<<10))
		if err != nil {
			notificationProblem(w, 400, "INVALID_BODY")
			return
		}
		if !p.AllowedOrganization(input.OrganizationID) {
			notificationProblem(w, 403, "ORGANIZATION_FORBIDDEN")
			return
		}
		request, _ := json.Marshal(struct {
			Recipient      string   `json:"recipient"`
			TemplateName   string   `json:"template_name"`
			LanguageCode   string   `json:"language_code"`
			BodyParameters []string `json:"body_parameters"`
		}{input.Recipient, input.TemplateName, input.LanguageCode, input.BodyParameters})
		key := AppointmentConfirmationDeliveryKey(p.TenantID, input.ConfirmationEventID)
		message := channels.Message{TenantID: p.TenantID, ChannelCode: "whatsapp", Direction: channels.DirectionOut, DeliveryKey: key, ExternalID: input.Recipient, Text: string(request)}
		command := AppointmentApprovalCommand{Message: message, OrganizationID: input.OrganizationID, AppointmentID: r.PathValue("id"), ConfirmationEventID: input.ConfirmationEventID, AppointmentVersion: input.AppointmentVersion, BindingVersion: input.BindingVersion, PolicyVersion: m.approvals.policy, ConsentID: input.ConsentID, ConsentPurpose: m.approvals.purpose, ConsentEvidenceSHA256: input.ConsentEvidenceSHA256, ProfileSHA256: digest(m.sender.Profile), EvidenceSHA256: input.EvidenceSHA256, ExpiresAt: input.ExpiresAt.UTC()}
		if _, err := m.approvals.Approve(ctx, p, command); err != nil {
			notificationProblem(w, 409, "APPROVAL_NOT_CURRENT_OR_DIVERGENT")
			return
		}
		// Preflight occurs before Claim; the Sender still revalidates after Claim.
		// Only the existing durable Channel may invoke the provider. A crash after
		// approval is recovered with the SAME body/event, never a new identity.
		if err := m.channel.Send(ctx, message); err != nil {
			code := "DELIVERY_RECONCILIATION_REQUIRED"
			if errors.Is(err, outbounddelivery.ErrInProgress) {
				code = "DELIVERY_IN_PROGRESS"
			} else if errors.Is(err, outbounddelivery.ErrTerminal) {
				code = "DELIVERY_TERMINAL"
			}
			notificationProblem(w, 409, code)
			return
		}
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "accepted", "delivery_key": key})
	})
}

func (m *AppointmentNotificationModule) notificationPrincipal(ctx context.Context, w http.ResponseWriter, r *http.Request, verifier identity.Verifier) (identity.Principal, bool) {
	// Shared Bearer-only boundary for reads and writes; no cookie/body identity.
	values := r.Header.Values("Authorization")
	if verifier == nil || len(values) != 1 || !strings.HasPrefix(values[0], "Bearer ") || len(values[0]) > 16391 || len(strings.Fields(values[0])) != 2 {
		notificationProblem(w, 401, "UNAUTHENTICATED")
		return identity.Principal{}, false
	}
	p, err := verifier.Verify(ctx, strings.TrimPrefix(values[0], "Bearer "))
	if err != nil || p.Subject == "" || p.TenantID == "" {
		notificationProblem(w, 401, "UNAUTHENTICATED")
		return identity.Principal{}, false
	}
	if m == nil || m.sender == nil || p.TenantID != m.sender.TenantID || !p.Allowed("appointment:manage") {
		notificationProblem(w, 403, "FORBIDDEN")
		return identity.Principal{}, false
	}
	return p, true
}

func notificationProblem(w http.ResponseWriter, status int, code string) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{"type": "urn:elite:problem:" + code, "status": status, "code": code})
}

// encoding/json accepts duplicate keys by default. Reject duplicates, aliases,
// nulls and trailing data before decoding this flat, versioned command shape.
func readNotificationRequest(reader io.Reader) (AppointmentNotificationRequest, error) {
	var input AppointmentNotificationRequest
	allowed := map[string]bool{"organization_id": true, "confirmation_event_id": true, "appointment_version": true, "binding_version": true, "consent_id": true, "consent_evidence_sha256": true, "evidence_sha256": true, "expires_at": true, "recipient": true, "template_name": true, "language_code": true, "body_parameters": true}
	decoder := json.NewDecoder(reader)
	first, err := decoder.Token()
	if err != nil || first != json.Delim('{') {
		return input, ErrBridge
	}
	fields := map[string]json.RawMessage{}
	for decoder.More() {
		name, err := decoder.Token()
		key, ok := name.(string)
		if err != nil || !ok || !allowed[key] || fields[key] != nil {
			return input, ErrBridge
		}
		var raw json.RawMessage
		if decoder.Decode(&raw) != nil || bytes.Equal(bytes.TrimSpace(raw), []byte("null")) {
			return input, ErrBridge
		}
		if key == "body_parameters" {
			var values []json.RawMessage
			if json.Unmarshal(raw, &values) != nil {
				return input, ErrBridge
			}
			for _, value := range values {
				var text string
				if bytes.Equal(bytes.TrimSpace(value), []byte("null")) || json.Unmarshal(value, &text) != nil {
					return input, ErrBridge
				}
			}
		}
		fields[key] = raw
	}
	last, err := decoder.Token()
	if err != nil || last != json.Delim('}') || decoder.Decode(new(any)) != io.EOF || len(fields) != len(allowed) {
		return input, ErrBridge
	}
	raw, err := json.Marshal(fields)
	if err != nil || json.Unmarshal(raw, &input) != nil {
		return input, ErrBridge
	}
	return input, nil
}
````

### FILE: `internal/whatsappbridge/appointment_notification_test.go`
```yaml
block_id: "PY-META-WHATSAPP:notification-http-1:v1"
operation: CREATE
provenance: AUTHORED
source: "local authenticated composition of admitted identity, approval and delivery owners"
license: "LicenseRef-Workspace-Owner"
sha256: "32ae8990b26158e727bc990fca3109183c5d4d98c32dda33dcd99a5172283a51"
variables: []
secrets_allowed: false
```
````go
package whatsappbridge

import (
	"bytes"
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/postgres"
)

// Same local RS256/discovery/JWKS method as the existing confirmation fixture.
// Uses the real identity verifier; the issuer and its identities are synthetic.
func notificationIssuer(t *testing.T) (identity.Verifier, func(identity.Principal) string) {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	enc := base64.RawURLEncoding.EncodeToString
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/.well-known/openid-configuration" {
			_ = json.NewEncoder(w).Encode(map[string]any{"issuer": server.URL, "jwks_uri": server.URL + "/keys", "id_token_signing_alg_values_supported": []string{"RS256"}})
		} else if r.URL.Path == "/keys" {
			_ = json.NewEncoder(w).Encode(map[string]any{"keys": []any{map[string]any{"kty": "RSA", "use": "sig", "alg": "RS256", "kid": "notification-fixture", "n": enc(key.N.Bytes()), "e": enc(big.NewInt(int64(key.E)).Bytes())}}})
		} else {
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(server.Close)
	verifier, err := identity.NewOIDCVerifier(context.Background(), server.URL, "notification-fixture")
	if err != nil {
		t.Fatal(err)
	}
	sign := func(p identity.Principal) string {
		permissions, organizations := []string{}, []string{}
		for s := range p.Permissions {
			permissions = append(permissions, s)
		}
		for s := range p.Organizations {
			organizations = append(organizations, s)
		}
		header, _ := json.Marshal(map[string]string{"alg": "RS256", "kid": "notification-fixture"})
		claims, _ := json.Marshal(map[string]any{"iss": server.URL, "aud": "notification-fixture", "sub": p.Subject, "tenant_id": p.TenantID, "permissions": permissions, "organization_ids": organizations, "exp": time.Now().Add(time.Minute).Unix()})
		unsigned := enc(header) + "." + enc(claims)
		hash := sha256.Sum256([]byte(unsigned))
		signature, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, hash[:])
		if err != nil {
			t.Fatal(err)
		}
		return unsigned + "." + enc(signature)
	}
	return verifier, sign
}

func notificationRequest(c AppointmentApprovalCommand) AppointmentNotificationRequest {
	return AppointmentNotificationRequest{OrganizationID: c.OrganizationID, ConfirmationEventID: c.ConfirmationEventID, AppointmentVersion: c.AppointmentVersion, BindingVersion: c.BindingVersion, ConsentID: c.ConsentID, ConsentEvidenceSHA256: c.ConsentEvidenceSHA256, EvidenceSHA256: c.EvidenceSHA256, ExpiresAt: c.ExpiresAt, Recipient: c.Message.ExternalID, TemplateName: "appointment_confirmed", LanguageCode: "es_AR", BodyParameters: []string{"2026-09-07T10:00Z"}}
}

func TestNotificationRequestRejectsAmbiguity(t *testing.T) {
	valid, _ := json.Marshal(notificationRequest(AppointmentApprovalCommand{ExpiresAt: time.Now()}))
	if _, err := readNotificationRequest(bytes.NewReader(valid)); err != nil {
		t.Fatal(err)
	}
	for name, body := range map[string]string{
		"null_parameter": strings.Replace(string(valid), `"body_parameters":["2026-09-07T10:00Z"]`, `"body_parameters":[null]`, 1),
		"duplicate":      strings.Replace(string(valid), `"recipient":`, `"recipient":"other","recipient":`, 1),
		"case_alias":     strings.Replace(string(valid), `"recipient":`, `"Recipient":`, 1),
		"unknown":        strings.Replace(string(valid), `"recipient":`, `"tenant_id":`, 1),
		"null":           strings.Replace(string(valid), `"recipient":""`, `"recipient":null`, 1),
		"trailing":       string(valid) + "{}",
		"array":          "[]", "empty": "{}", "number": "1", "malformed": "{",
		"bad_type": strings.Replace(string(valid), `"binding_version":0`, `"binding_version":"0"`, 1),
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := readNotificationRequest(strings.NewReader(body)); err == nil {
				t.Fatal("ambiguous request accepted")
			}
		})
	}
}

func TestNotificationHTTPPostgresReplayAndRecovery(t *testing.T) {
	pool := approvalPool(t)
	verifier, sign := notificationIssuer(t)
	for _, mode := range []string{"accepted", "provider_lost", "http_lost"} {
		t.Run(mode, func(t *testing.T) {
			p, c, resolver := approvalFixtureData(t, pool)
			providerMode := "accepted"
			if mode == "provider_lost" {
				providerMode = "lost"
			}
			sender, _, marker := senderFixture(t, providerMode)
			sender.TenantID = p.TenantID
			store, err := postgres.NewOutboundDeliveryStore(pool, []byte("0123456789abcdef0123456789abcdef"), time.Minute)
			if err != nil {
				t.Fatal(err)
			}
			module, err := NewAppointmentNotificationModule(resolver, sender, store, receiverFixture{})
			if err != nil {
				t.Fatal(err)
			}
			mux := http.NewServeMux()
			module.Register(mux, verifier)
			var dropped atomic.Bool
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if mode == "http_lost" && dropped.CompareAndSwap(false, true) {
					record := httptest.NewRecorder()
					mux.ServeHTTP(record, r)
					if record.Code != 200 {
						t.Errorf("expected commit before lost HTTP response: %d %s", record.Code, record.Body.String())
					}
					panic(http.ErrAbortHandler)
				}
				mux.ServeHTTP(w, r)
			}))
			defer server.Close()
			body, _ := json.Marshal(notificationRequest(c))
			token := sign(p)
			call := func() (int, string, error) {
				req, err := http.NewRequest(http.MethodPost, server.URL+"/v1/franchise/appointments/appointment-1/whatsapp-confirmation", bytes.NewReader(body))
				if err != nil {
					return 0, "", err
				}
				req.Header.Set("Authorization", "Bearer "+token)
				req.Header.Set("Content-Type", "application/json; charset=utf-8")
				response, err := server.Client().Do(req)
				if err != nil {
					return 0, "", err
				}
				defer response.Body.Close()
				raw, err := io.ReadAll(response.Body)
				return response.StatusCode, string(raw), err
			}
			if mode == "http_lost" {
				if _, _, err := call(); err == nil {
					t.Fatal("HTTP response loss was not observed")
				}
			}
			if mode == "accepted" {
				var wg sync.WaitGroup
				for i := 0; i < 8; i++ {
					wg.Add(1)
					go func() {
						defer wg.Done()
						status, _, err := call()
						if err != nil || (status != 200 && status != 409) {
							t.Errorf("concurrent call: %d %v", status, err)
						}
					}()
				}
				wg.Wait()
			}
			for i := 0; i < 2; i++ {
				status, response, err := call()
				want := 200
				if mode == "provider_lost" {
					want = 409
				}
				if err != nil || status != want {
					t.Fatalf("%d %s %v", status, response, err)
				}
				if want == 200 && (!strings.Contains(response, `"status":"accepted"`) || strings.Contains(response, "delivered")) {
					t.Fatal(response)
				}
				if want == 409 && !strings.Contains(response, "DELIVERY_RECONCILIATION_REQUIRED") {
					t.Fatal(response)
				}
				if strings.Contains(response, c.Message.ExternalID) || strings.Contains(response, "synthetic-token") {
					t.Fatal("response leaked sensitive data")
				}
			}
			calls, err := os.ReadFile(marker)
			if err != nil || string(calls) != "call\n" {
				t.Fatal("duplicate provider call", string(calls), err)
			}
			var grants, attempts int
			var state string
			if err = pool.QueryRow(context.Background(), `select count(*) from communication.whatsapp_appointment_approval where tenant_id=$1`, p.TenantID).Scan(&grants); err != nil || grants != 1 {
				t.Fatal(grants, err)
			}
			if err = pool.QueryRow(context.Background(), `select state,attempt_count from communication.outbound_delivery where tenant_id=$1`, p.TenantID).Scan(&state, &attempts); err != nil || attempts != 1 {
				t.Fatal(state, attempts, err)
			}
			wantState := "accepted"
			if mode == "provider_lost" {
				wantState = "unknown"
			}
			if state != wantState {
				t.Fatal(state)
			}
			// The identical event cannot silently change template parameters on replay.
			body = bytes.Replace(body, []byte("2026-09-07T10:00Z"), []byte("changed"), 1)
			if status, _, err := call(); err != nil || status != 409 {
				t.Fatal("divergent replay", status, err)
			}
		})
	}
}

func TestNotificationHTTPRejectsBeforeFence(t *testing.T) {
	pool := approvalPool(t)
	verifier, sign := notificationIssuer(t)
	for _, mode := range []string{"no_token", "invalid_signature", "role", "organization", "tenant", "recipient", "consent", "revoked", "path", "oversized", "media", "injected_identity", "duplicate_auth"} {
		t.Run(mode, func(t *testing.T) {
			p, c, resolver := approvalFixtureData(t, pool)
			sender, _, marker := senderFixture(t, "accepted")
			sender.TenantID = p.TenantID
			store, err := postgres.NewOutboundDeliveryStore(pool, []byte("0123456789abcdef0123456789abcdef"), time.Minute)
			if err != nil {
				t.Fatal(err)
			}
			module, err := NewAppointmentNotificationModule(resolver, sender, store, receiverFixture{})
			if err != nil {
				t.Fatal(err)
			}
			mux := http.NewServeMux()
			module.Register(mux, verifier)
			input := notificationRequest(c)
			status := 409
			pathID := "appointment-1"
			switch mode {
			case "no_token", "invalid_signature", "duplicate_auth":
				status = 401
			case "role":
				p.Permissions = map[string]struct{}{}
				status = 403
			case "organization":
				p.Organizations = map[string]struct{}{}
				status = 403
			case "tenant":
				p.TenantID = "other-tenant"
				status = 403
			case "recipient":
				input.Recipient = "5491199999999"
			case "consent":
				input.ConsentEvidenceSHA256 = strings.Repeat("f", 64)
			case "path":
				pathID = "other-appointment"
			case "revoked":
				if _, err := pool.Exec(context.Background(), `update communication.contact_channel_binding set state='revoked',pii_allowed=false,version=version+1 where tenant_id=$1`, p.TenantID); err != nil {
					t.Fatal(err)
				}
			case "oversized", "injected_identity":
				status = 400
			case "media":
				status = 415
			}
			body, _ := json.Marshal(input)
			if mode == "oversized" {
				body = bytes.Replace(body, []byte(input.Recipient), bytes.Repeat([]byte("1"), 65536), 1)
			}
			if mode == "injected_identity" {
				body = append([]byte(`{"tenant_id":"injected",`), body[1:]...)
			}
			req := httptest.NewRequest(http.MethodPost, "/v1/franchise/appointments/"+pathID+"/whatsapp-confirmation", bytes.NewReader(body))
			token := sign(p)
			if mode == "invalid_signature" {
				token += "tamper"
			}
			if mode != "no_token" {
				req.Header.Set("Authorization", "Bearer "+token)
			}
			if mode == "duplicate_auth" {
				req.Header.Add("Authorization", "Bearer "+token)
			}
			req.Header.Set("Content-Type", "application/json")
			if mode == "media" {
				req.Header.Set("Content-Type", "application/jsonp")
			}
			record := httptest.NewRecorder()
			mux.ServeHTTP(record, req)
			if record.Code != status || record.Header().Get("Cache-Control") != "no-store" || record.Header().Get("Content-Type") != "application/problem+json" {
				t.Fatal(record.Code, record.Body.String())
			}
			if _, err := os.Stat(marker); !os.IsNotExist(err) {
				t.Fatal("child invoked for rejected request")
			}
			for _, table := range []string{"communication.whatsapp_appointment_approval", "communication.outbound_delivery"} {
				var count int
				if err := pool.QueryRow(context.Background(), `select count(*) from `+table+` where tenant_id=$1`, c.Message.TenantID).Scan(&count); err != nil || count != 0 {
					t.Fatal(table, count, err)
				}
			}
		})
	}
}
````


### FILE: `internal/whatsappbridge/appointment_notification_status.go`
```yaml
block_id: "PY-META-WHATSAPP:notification-status-0:v1"
operation: CREATE
provenance: AUTHORED
source: "local read-only composition of admitted approval and delivery owners"
license: "LicenseRef-Workspace-Owner"
sha256: "b7beff2f10f586a096fe1e40c57de272fdf12814e5674f7fe31530d2e8e3a5fe"
variables: []
secrets_allowed: false
```
````go
package whatsappbridge

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"elite.local/enterprise/internal/platform/identity"
	"github.com/jackc/pgx/v5"
)

var (
	ErrNotificationNotFound  = errors.New("whatsappbridge: scoped notification not found")
	ErrNotificationRead      = errors.New("whatsappbridge: notification read unavailable")
	ErrNotificationIntegrity = errors.New("whatsappbridge: notification evidence mismatch")
	notificationEventID      = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
)

// NotificationStatus is a bounded local snapshot, not a provider delivery claim
// or permission to send. Hashes, recipient, actor and template are not exposed.
type NotificationStatus struct {
	DeliveryKey            string     `json:"delivery_key"`
	FenceState             string     `json:"fence_state"`
	DeliveryStatus         string     `json:"delivery_status"`
	ProviderTimestamp      *int64     `json:"provider_timestamp,omitempty"`
	ProviderEventCount     int64      `json:"provider_event_count"`
	ApprovalExpiresAt      time.Time  `json:"approval_expires_at"`
	ApprovalExpired        bool       `json:"approval_expired"`
	AcceptedAt             *time.Time `json:"accepted_at,omitempty"`
	UpdatedAt              *time.Time `json:"updated_at,omitempty"`
	ObservedAt             time.Time  `json:"observed_at"`
	ReconciliationRequired bool       `json:"reconciliation_required"`
}

// ReadNotificationStatus deliberately does not call ResolveWhatsAppApproval:
// an expired/revoked sending grant must not hide the historical outcome from
// a currently authorized operator. Current appointment/lead/org scope is checked.
func (s *PostgresAppointmentApprovals) ReadNotificationStatus(ctx context.Context, p identity.Principal, organization, appointment, event string) (NotificationStatus, error) {
	if s == nil || s.pool == nil {
		return NotificationStatus{}, ErrNotificationNotFound
	}
	return readNotificationStatus(ctx, s.pool, p, organization, appointment, event)
}

type notificationQueryer interface {
	QueryRow(context.Context, string, ...any) pgx.Row
}

func readNotificationStatus(ctx context.Context, db notificationQueryer, p identity.Principal, organization, appointment, event string) (NotificationStatus, error) {
	value := NotificationStatus{}
	if db == nil || p.Subject == "" || p.TenantID == "" || !p.AllowedOrganization(organization) || !((appointment != "" && p.Allowed("appointment:manage") && notificationEventID.MatchString(event)) || (appointment == "" && p.Allowed("whatsapp:approve") && strings.HasPrefix(event, "wa-reply:") && validDigest(strings.TrimPrefix(event, "wa-reply:")))) {
		return value, ErrNotificationNotFound
	}
	var requestHash, recipientHash, expectedHash, expectedRecipient, providerHash, evidenceHash string
	var lease *time.Time
	err := db.QueryRow(ctx, `select g.delivery_key,g.expires_at,g.expires_at<=statement_timestamp(),
 coalesce(d.state,'not_started'),d.accepted_at,d.updated_at,statement_timestamp(),d.locked_until,
 coalesce(d.request_sha256_hex,''),coalesce(d.recipient_hmac,''),g.message_sha256,g.external_id_hmac,
 coalesce(d.provider_message_hmac,''),coalesce(d.evidence_sha256_hex,''), totals.event_count,totals.last_time,
 case when totals.event_count=0 then 'not_observed_by_this_reader' else latest.status end
 from communication.whatsapp_delivery_approval g
 join crm.lead l on l.tenant_id=g.tenant_id and l.lead_id=g.lead_id and l.organization_id=g.organization_id
 left join communication.outbound_delivery d on d.tenant_id=g.tenant_id and d.channel_code=g.channel_code and d.delivery_key=g.delivery_key
 left join lateral (select count(*) event_count,max(provider_timestamp) last_time from communication.whatsapp_status_observation s where s.tenant_id=g.tenant_id and s.delivery_key=g.delivery_key) totals on true
 left join lateral (select case when count(distinct provider_status)=1 then 'observed_'||min(provider_status) else 'ambiguous_latest_timestamp' end status from communication.whatsapp_status_observation s where s.tenant_id=g.tenant_id and s.delivery_key=g.delivery_key and s.provider_timestamp=totals.last_time) latest on true
 where g.tenant_id=$1 and g.organization_id=$2 and g.appointment_id=$3 and g.approval_event_id=$4 and g.channel_code='whatsapp'`, p.TenantID, organization, appointment, event).Scan(&value.DeliveryKey, &value.ApprovalExpiresAt, &value.ApprovalExpired, &value.FenceState, &value.AcceptedAt, &value.UpdatedAt, &value.ObservedAt, &lease, &requestHash, &recipientHash, &expectedHash, &expectedRecipient, &providerHash, &evidenceHash, &value.ProviderEventCount, &value.ProviderTimestamp, &value.DeliveryStatus)
	if errors.Is(err, pgx.ErrNoRows) {
		return NotificationStatus{}, ErrNotificationNotFound
	}
	if err != nil {
		return NotificationStatus{}, ErrNotificationRead
	}
	if value.FenceState != "not_started" && (requestHash != expectedHash || recipientHash != expectedRecipient) {
		return NotificationStatus{}, ErrNotificationIntegrity
	}
	switch value.FenceState {
	case "not_started":
	case "accepted":
		if value.AcceptedAt == nil || value.AcceptedAt.IsZero() || !validDigest(providerHash) || !validDigest(evidenceHash) {
			return NotificationStatus{}, ErrNotificationIntegrity
		}
	case "sending":
		// Read-only: report stale lease without changing state or creating events.
		value.ReconciliationRequired = lease == nil || !lease.After(value.ObservedAt)
	case "unknown":
		value.ReconciliationRequired = true
	case "failed_terminal":
		if !validDigest(evidenceHash) {
			return NotificationStatus{}, ErrNotificationIntegrity
		}
	default:
		return NotificationStatus{}, ErrNotificationIntegrity
	}
	if value.ProviderEventCount > 0 && value.FenceState != "accepted" {
		return NotificationStatus{}, ErrNotificationIntegrity
	}
	return value, nil
}

func (m *AppointmentNotificationModule) registerNotificationStatus(mux *http.ServeMux, verifier identity.Verifier) {
	mux.HandleFunc("GET /v1/franchise/appointments/{id}/whatsapp-confirmation", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		p, ok := m.notificationPrincipal(ctx, w, r, verifier)
		if !ok {
			return
		}
		query, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil || len(r.URL.RawQuery) > 2048 || len(query) != 2 || len(query["organization_id"]) != 1 || len(query["confirmation_event_id"]) != 1 || query.Get("organization_id") == "" || !notificationEventID.MatchString(query.Get("confirmation_event_id")) || r.ContentLength != 0 {
			notificationProblem(w, 400, "INVALID_QUERY")
			return
		}
		organization := query.Get("organization_id")
		if !p.AllowedOrganization(organization) {
			notificationProblem(w, 403, "ORGANIZATION_FORBIDDEN")
			return
		}
		value, err := m.approvals.ReadNotificationStatus(ctx, p, organization, r.PathValue("id"), query.Get("confirmation_event_id"))
		if errors.Is(err, ErrNotificationNotFound) {
			notificationProblem(w, 404, "NOTIFICATION_NOT_FOUND")
			return
		}
		if errors.Is(err, ErrNotificationIntegrity) {
			notificationProblem(w, 409, "NOTIFICATION_EVIDENCE_MISMATCH")
			return
		}
		if err != nil {
			notificationProblem(w, 503, "NOTIFICATION_READ_UNAVAILABLE")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-store")
		_ = json.NewEncoder(w).Encode(value)
	})
}
````

### FILE: `internal/whatsappbridge/appointment_notification_status_test.go`
```yaml
block_id: "PY-META-WHATSAPP:notification-status-1:v1"
operation: CREATE
provenance: AUTHORED
source: "local read-only composition of admitted approval and delivery owners"
license: "LicenseRef-Workspace-Owner"
sha256: "8158977f24dddc9cf2cd9725ae4f918347c3d7485cfd46b3edf62a45fb793d97"
variables: []
secrets_allowed: false
```
````go
package whatsappbridge

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"elite.local/enterprise/internal/outbounddelivery"
	"elite.local/enterprise/internal/platform/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestNotificationStatusReadOnlyStates(t *testing.T) {
	pool := approvalPool(t)
	verifier, sign := notificationIssuer(t)
	config, err := pgxpool.ParseConfig(os.Getenv("ELITE_WHATSAPP_TEST_DATABASE_URL"))
	if err != nil {
		t.Fatal(err)
	}
	config.ConnConfig.RuntimeParams["default_transaction_read_only"] = "on"
	readPool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}
	defer readPool.Close()
	var readOnly string
	if err := readPool.QueryRow(context.Background(), "show transaction_read_only").Scan(&readOnly); err != nil || readOnly != "on" {
		t.Fatal("read-only test session not demonstrated", readOnly, err)
	}
	for _, mode := range []string{"not_started", "sending", "stale_sending", "accepted", "unknown", "failed_terminal", "expired_accepted", "cancelled", "revoked", "request_mismatch", "recipient_mismatch"} {
		t.Run(mode, func(t *testing.T) {
			ctx := context.Background()
			p, c, writer := approvalFixtureData(t, pool)
			if mode == "expired_accepted" {
				c.ExpiresAt = time.Now().Add(500 * time.Millisecond)
			}
			if _, err := writer.Approve(ctx, p, c); err != nil {
				t.Fatal(err)
			}
			store, err := postgres.NewOutboundDeliveryStore(pool, []byte("0123456789abcdef0123456789abcdef"), time.Minute)
			if err != nil {
				t.Fatal(err)
			}
			hash, err := outbounddelivery.MessageSHA256(c.Message)
			if err != nil {
				t.Fatal(err)
			}
			want := mode
			if mode != "not_started" {
				if _, err := store.Claim(ctx, c.Message, hash); err != nil {
					t.Fatal(err)
				}
			}
			switch mode {
			case "unknown":
				if err := store.MarkUnknown(ctx, c.Message, hash, "PROVIDER_CALL_UNCERTAIN"); err != nil {
					t.Fatal(err)
				}
			case "failed_terminal":
				if err := store.MarkFailed(ctx, c.Message, hash, strings.Repeat("e", 64), "PROVEN_REJECTED"); err != nil {
					t.Fatal(err)
				}
			case "stale_sending":
				want = "sending"
				if _, err := pool.Exec(ctx, `update communication.outbound_delivery set locked_until=clock_timestamp()-interval '1 second' where tenant_id=$1`, p.TenantID); err != nil {
					t.Fatal(err)
				}
			case "accepted", "expired_accepted", "cancelled", "revoked", "request_mismatch", "recipient_mismatch":
				want = "accepted"
				if err := store.Complete(ctx, c.Message, hash, outbounddelivery.Receipt{ProviderMessageID: "wamid.status-fixture", EvidenceSHA256: strings.Repeat("a", 64), AcceptedAt: time.Now()}); err != nil {
					t.Fatal(err)
				}
				queries := map[string]string{
					"cancelled":          `update crm.appointment set state='cancelled',version=version+1 where tenant_id=$1`,
					"revoked":            `update communication.contact_channel_binding set state='revoked',pii_allowed=false,version=version+1 where tenant_id=$1`,
					"request_mismatch":   `update communication.outbound_delivery set request_sha256_hex=repeat('f',64) where tenant_id=$1`,
					"recipient_mismatch": `update communication.outbound_delivery set recipient_hmac=repeat('f',64) where tenant_id=$1`,
				}
				if q := queries[mode]; q != "" {
					if _, err := pool.Exec(ctx, q, p.TenantID); err != nil {
						t.Fatal(err)
					}
				}
				if mode == "expired_accepted" {
					time.Sleep(550 * time.Millisecond)
				}
			}
			reader, err := NewPostgresAppointmentApprovals(readPool, []byte("0123456789abcdef0123456789abcdef"), "fixture-appointment-notification", "policy-1")
			if err != nil {
				t.Fatal(err)
			}
			sender, _, marker := senderFixture(t, "accepted")
			sender.TenantID = p.TenantID
			module, err := NewAppointmentNotificationModule(reader, sender, store, receiverFixture{})
			if err != nil {
				t.Fatal(err)
			}
			mux := http.NewServeMux()
			module.Register(mux, verifier)
			var before int
			if err := pool.QueryRow(ctx, `select count(*) from communication.outbound_delivery_event where tenant_id=$1`, p.TenantID).Scan(&before); err != nil {
				t.Fatal(err)
			}
			for i := 0; i < 2; i++ {
				req := httptest.NewRequest(http.MethodGet, "/v1/franchise/appointments/appointment-1/whatsapp-confirmation?organization_id=store-1&confirmation_event_id="+c.ConfirmationEventID, nil)
				req.Header.Set("Authorization", "Bearer "+sign(p))
				response := httptest.NewRecorder()
				mux.ServeHTTP(response, req)
				if mode == "request_mismatch" || mode == "recipient_mismatch" {
					if response.Code != 409 || !strings.Contains(response.Body.String(), "NOTIFICATION_EVIDENCE_MISMATCH") {
						t.Fatal(response.Code, response.Body.String())
					}
					continue
				}
				if response.Code != 200 || response.Header().Get("Cache-Control") != "no-store" {
					t.Fatal(response.Code, response.Body.String())
				}
				var value NotificationStatus
				if err := json.Unmarshal(response.Body.Bytes(), &value); err != nil {
					t.Fatal(err)
				}
				if value.FenceState != want || value.DeliveryStatus != "not_observed_by_this_reader" || value.DeliveryKey != c.Message.DeliveryKey || value.ObservedAt.IsZero() {
					t.Fatal(value)
				}
				if value.ReconciliationRequired != (mode == "unknown" || mode == "stale_sending") {
					t.Fatal(value)
				}
				if value.ApprovalExpired != (mode == "expired_accepted") {
					t.Fatal(value)
				}
				if (value.AcceptedAt != nil) != (want == "accepted") {
					t.Fatal(value)
				}
				for _, secret := range []string{c.Message.ExternalID, p.Subject, c.ConsentID, c.EvidenceSHA256, "wamid.status-fixture", "appointment_confirmed", "synthetic-token"} {
					if strings.Contains(response.Body.String(), secret) {
						t.Fatal("sensitive field exposed")
					}
				}
			}
			var after int
			if err := pool.QueryRow(ctx, `select count(*) from communication.outbound_delivery_event where tenant_id=$1`, p.TenantID).Scan(&after); err != nil || after != before {
				t.Fatal("GET mutated events", before, after, err)
			}
			if _, err := os.Stat(marker); !os.IsNotExist(err) {
				t.Fatal("GET invoked provider")
			}
			if mode == "stale_sending" {
				var state string
				if err := pool.QueryRow(ctx, `select state from communication.outbound_delivery where tenant_id=$1`, p.TenantID).Scan(&state); err != nil || state != "sending" {
					t.Fatal("GET mutated stale fence", state, err)
				}
			}
		})
	}
}

func TestNotificationStatusScopeAndUnavailable(t *testing.T) {
	pool := approvalPool(t)
	verifier, sign := notificationIssuer(t)
	for _, mode := range []string{"no_token", "role", "organization", "tenant", "appointment", "event", "missing_grant", "reassigned", "duplicate_query", "unknown_query", "bad_uuid", "body", "unavailable"} {
		t.Run(mode, func(t *testing.T) {
			ctx := context.Background()
			p, c, resolver := approvalFixtureData(t, pool)
			if mode != "missing_grant" {
				if _, err := resolver.Approve(ctx, p, c); err != nil {
					t.Fatal(err)
				}
			}
			sender, _, marker := senderFixture(t, "accepted")
			sender.TenantID = p.TenantID
			store, err := postgres.NewOutboundDeliveryStore(pool, []byte("0123456789abcdef0123456789abcdef"), time.Minute)
			if err != nil {
				t.Fatal(err)
			}
			if mode == "unavailable" {
				closed := approvalPool(t)
				closed.Close()
				resolver, err = NewPostgresAppointmentApprovals(closed, []byte("0123456789abcdef0123456789abcdef"), "fixture-appointment-notification", "policy-1")
				if err != nil {
					t.Fatal(err)
				}
			}
			module, err := NewAppointmentNotificationModule(resolver, sender, store, receiverFixture{})
			if err != nil {
				t.Fatal(err)
			}
			mux := http.NewServeMux()
			module.Register(mux, verifier)
			query := url.Values{"organization_id": {"store-1"}, "confirmation_event_id": {c.ConfirmationEventID}}
			pathID := "appointment-1"
			want := 404
			var body *strings.Reader = strings.NewReader("")
			switch mode {
			case "no_token":
				want = 401
			case "role":
				p.Permissions = map[string]struct{}{}
				want = 403
			case "organization":
				p.Organizations = map[string]struct{}{}
				want = 403
			case "tenant":
				p.TenantID = "other-tenant"
				want = 403
			case "appointment":
				pathID = "other"
			case "event":
				query.Set("confirmation_event_id", "00000000-0000-0000-0000-000000000000")
			case "reassigned":
				if _, err := pool.Exec(ctx, `insert into crm.lead(tenant_id,lead_id,organization_id,lifecycle_state,source_code,contact_payload,consent_required) values($1,'other-lead','store-1','new','fixture','{}',true)`, p.TenantID); err != nil {
					t.Fatal(err)
				}
				if _, err := pool.Exec(ctx, `update crm.appointment set lead_id='other-lead' where tenant_id=$1`, p.TenantID); err != nil {
					t.Fatal(err)
				}
			case "duplicate_query":
				query.Add("organization_id", "store-1")
				want = 400
			case "unknown_query":
				query.Set("tenant_id", p.TenantID)
				want = 400
			case "bad_uuid":
				query.Set("confirmation_event_id", "not-an-event")
				want = 400
			case "body":
				body = strings.NewReader("{}")
				want = 400
			case "unavailable":
				want = 503
			}
			req := httptest.NewRequest(http.MethodGet, "/v1/franchise/appointments/"+pathID+"/whatsapp-confirmation?"+query.Encode(), body)
			if mode != "no_token" {
				req.Header.Set("Authorization", "Bearer "+sign(p))
			}
			response := httptest.NewRecorder()
			mux.ServeHTTP(response, req)
			if response.Code != want || response.Header().Get("Cache-Control") != "no-store" || response.Header().Get("Content-Type") != "application/problem+json" {
				t.Fatal(response.Code, response.Body.String())
			}
			if strings.Contains(response.Body.String(), "closed") || strings.Contains(response.Body.String(), c.Message.ExternalID) {
				t.Fatal("internal error leaked")
			}
			if _, err := os.Stat(marker); !os.IsNotExist(err) {
				t.Fatal("GET invoked child")
			}
		})
	}
}

func TestNotificationPostThenStatus(t *testing.T) {
	pool := approvalPool(t)
	p, c, resolver := approvalFixtureData(t, pool)
	verifier, sign := notificationIssuer(t)
	sender, _, marker := senderFixture(t, "accepted")
	sender.TenantID = p.TenantID
	store, err := postgres.NewOutboundDeliveryStore(pool, []byte("0123456789abcdef0123456789abcdef"), time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	module, err := NewAppointmentNotificationModule(resolver, sender, store, receiverFixture{})
	if err != nil {
		t.Fatal(err)
	}
	mux := http.NewServeMux()
	module.Register(mux, verifier)
	server := httptest.NewServer(mux)
	defer server.Close()
	body, _ := json.Marshal(notificationRequest(c))
	token := sign(p)
	for _, method := range []string{http.MethodPost, http.MethodGet, http.MethodGet} {
		path := "/v1/franchise/appointments/appointment-1/whatsapp-confirmation"
		data := []byte(nil)
		if method == http.MethodPost {
			data = body
		} else {
			path += "?organization_id=store-1&confirmation_event_id=" + c.ConfirmationEventID
		}
		req, err := http.NewRequest(method, server.URL+path, bytes.NewReader(data))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		response, err := server.Client().Do(req)
		if err != nil {
			t.Fatal(err)
		}
		if response.StatusCode != 200 {
			response.Body.Close()
			t.Fatal(response.StatusCode)
		}
		if method == http.MethodGet {
			var value NotificationStatus
			err = json.NewDecoder(response.Body).Decode(&value)
			if err != nil || value.FenceState != "accepted" || value.DeliveryStatus != "not_observed_by_this_reader" {
				response.Body.Close()
				t.Fatal(value, err)
			}
		}
		response.Body.Close()
	}
	calls, err := os.ReadFile(marker)
	if err != nil || string(calls) != "call\n" {
		t.Fatal("read caused duplicate", string(calls), err)
	}
}
````


### FILE: `db/migrations/0059_whatsapp_conversation_reply.down.sql`

```yaml
block_id: "COMMUNICATIONS-DELTA-PYTHON_META_WHATSAPP_CLOUD_ADAPTER:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "c7c6b2274846501afb0630401de6c0fef1b5fe7f9d5ee6396625b4a6e7e9d8f4"
variables: []
secrets_allowed: false
```

````sql
begin;
-- Removing this owner while conversational status evidence exists would sever
-- its FK. Fail closed; export/retain evidence before a deliberate downgrade.
do $$ begin if exists(select 1 from communication.whatsapp_status_batch b where not exists(select 1 from communication.whatsapp_appointment_approval a where a.tenant_id=b.tenant_id and a.delivery_key=b.delivery_key)) then raise exception 'conversational status evidence prevents downgrade'; end if; end $$;
drop view communication.whatsapp_delivery_approval;
alter table communication.whatsapp_status_batch drop constraint whatsapp_status_batch_outbound_fkey;
alter table communication.whatsapp_status_batch drop column channel_code;
alter table communication.whatsapp_status_batch add foreign key(tenant_id,delivery_key) references communication.whatsapp_appointment_approval(tenant_id,delivery_key);
commit;
````

### FILE: `db/migrations/0059_whatsapp_conversation_reply.up.sql`

```yaml
block_id: "COMMUNICATIONS-DELTA-PYTHON_META_WHATSAPP_CLOUD_ADAPTER:file2:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "9fab24f18f1539433d23ea6a022683f925152748f2ddb32eafdf424c77969601"
variables: []
secrets_allowed: false
```

````sql
begin;
-- AUTHORED extension of the existing status owner: every status batch anchors
-- an already accepted outbound row, whether its approval was appointment or
-- an exact human-reviewed conversational response.
alter table communication.whatsapp_status_batch drop constraint whatsapp_status_batch_tenant_id_delivery_key_fkey;
alter table communication.whatsapp_status_batch add column channel_code text not null default 'whatsapp' check(channel_code='whatsapp');
alter table communication.whatsapp_status_batch add constraint whatsapp_status_batch_outbound_fkey foreign key(tenant_id,channel_code,delivery_key) references communication.outbound_delivery(tenant_id,channel_code,delivery_key);

create view communication.whatsapp_delivery_approval as
 select g.tenant_id,g.channel_code,g.delivery_key,g.organization_id,g.appointment_id,g.confirmation_event_id::text as approval_event_id,g.lead_id,g.external_id_hmac,g.profile_sha256,g.message_sha256,g.expires_at
 from communication.whatsapp_appointment_approval g
 union all
 select r.tenant_id,'whatsapp',r.request_id,r.organization_id,''::text,r.request_id,
 r.payload->>'lead_id',r.payload->>'external_id_hmac',r.payload->>'profile_sha256',r.payload->>'message_sha256',(r.payload->>'expires_at')::timestamptz
 from approval.request r
 where r.kind='whatsapp_reply' and r.state='approved' and r.payload->>'schema'='elite-whatsapp-reply-approval/v1'
 and r.payload->>'organization_id'=r.organization_id;
commit;
````

### FILE: `docs/whatsapp-conversation-operations.md`

```yaml
block_id: "COMMUNICATIONS-DELTA-PYTHON_META_WHATSAPP_CLOUD_ADAPTER:file3:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "4358ce0dd86e282793af762b47552ddac673ca6484ea1bf952080fbf23caf1e6"
variables: []
secrets_allowed: false
```

````markdown
# V402 WhatsApp connected host contract

Version: V402. This supplements the existing whatsapp_cloud/provider-profile.template.json and docs/whatsapp-status-operations.md. It is infrastructure configuration, not an account proof. Values remain absent until supplied/configured by the user; no secrets are embedded in profile artifacts.

## Exact composition

1. Load a bounded, exact-SHA provider profile from a private configured path. Python validate_profile remains normative. Keep all template/readiness/source fields, exact source_commit de70ee908a67026e642aaee3703d20464e2a9466, Graph version/WABA/phone and explicit approved_templates. Template default values are blocked; fixture profile() in whatsapp_cloud/test_whatsapp_cloud.py is synthetic only. No automatic PROVEN defaults.
2. Resolve the same server-side HMAC key for postgres.NewContactIdentityStore, postgres.NewOutboundDeliveryStore(pool,key,time.Minute), NewPostgresAppointmentApprovals and NewStatusRouter. Scope is fixed tenant+organization+connection; require durable integration.provider_connection with provider_code meta-whatsapp and that organization, active.
3. base = whatsappbridge.NewPostgresAppointmentApprovals(pool,key,purpose,policy). Purpose/policy are explicit configured consent contract values, not invented business defaults. replies = NewPostgresReplyApprovals(base,tenant,organization,connection,profileJSON).
4. Process holds PythonExecutable/PythonSHA256, AdapterDirectory/AdapterSHA256 (whatsapp_cloud.py), EvidenceDirectory. ReconcilerSHA256 hashes status_reconciliation.py independently. Hashes must match the materialized revision; directories private. Existing Process launches -I -B and no inherited secrets/proxies/PYTHONPATH.
5. Sender{TenantID,Profile,Tokens,Process}. The token source is runtime-only; other sources implement WhatsAppAppSecret(ctx) and WhatsAppVerifyToken(ctx). NewReplyModule(replies,sender,outboundStore) forces the same approval resolver and registered durable channel. Register(mux,real OIDC verifier). Exports VerifiedInboxChannel{} for the receiver side of an outbounddelivery.Channel or NewAppointmentNotificationModule; it refuses direct receive/send because inbox+job owns ingress.
6. NewWebhookReceiver(WebhookReceiverConfig{TenantID,ConnectionID,RetentionApprovalSHA256,Profile,Process,Secrets,Verification,Store:postgres.NewProviderIntegration(pool),MaxConcurrent}) mounts raw signed /messages ingress. The exact existing callback path is host-selected. Ingress returns receipt only, not a conversation response.
7. Construct the existing app.New(Config) once. OrganizationID and LeadID legacy defaults BOTH EMPTY: contact resolver provides scope. ContactResolver=whatsappbridge.ScopedContactResolver{TenantID,OrganizationID,Resolver:postgresContactStore}. Existing ConversationStore, LLM API/model/base, domain token source/base, variants/services/pricebook, finite token budget, approved retention and handoff text remain explicit configuration. Register a durable WhatsApp channel, but do not execute Dispatcher on webhook replies. The admitted Runtime.Handle is the only callable conversation entry here.
8. StatusObserver{TenantID,Profile,Process,ReconcilerSHA256,Secrets,Approvals:base}; router=NewStatusRouter(observer,connection,retentionSHA,key,8); router.EnableConversation(ConversationRoute{OrganizationID,Runtime:app.Conversation,Proposals:replies}). NewStatusWorker(router,workerID,2*time.Minute,retry). Run(ctx,StatusPrincipalSource,StatusReporter,interval). One existing provider-events queue processes mixed messages/statuses. No parallel status-only consumer.
9. Service identity is revalidated for each poll using actual OIDC/service auth; requires appointment:manage (existing status owner), whatsapp:process and configured organization. A provider webhook never creates a principal. StatusReporter must use the existing bounded exclusive authenticated connection. Failure to authenticate or report stops/holds processing; no invented subject/permissions.

## Human HTTP contract

Bearer-only same verifier/tenant/org, no cookie/body identity. All responses no-store. Human review permission whatsapp:approve; sends/recovery additionally whatsapp:send. Service requester differs from human reviewer through the common approval owner.

- GET /v1/franchise/whatsapp/replies returns at most 50 exact pending/completed request projections, newest first.
- GET /v1/franchise/whatsapp/replies/{request_id} returns request_id,state,payload_sha256,context (exact message text/recipient/service-window/scope).
- POST /.../{request_id}/decision body {payload_sha256,approved,reason}. Approve/reject immutable; no message text/contact fields accepted. Approval itself sends nothing.
- POST /.../{request_id}/send body {payload_sha256}. Reads current bindings/consent/window and uses the existing fence. Accepted does not mean delivered. Retry retains exact identity and never creates a second POST after unknown.
- POST /.../{request_id}/recover body {payload_sha256}. Reads bounded existing local SEND_RECEIPT/provider-response through os.Root and the hash-locked Python evidence validator. Never performs HTTP. With a proven receipt, calls existing ReconcileAccepted; without it returns reconciliation required. A bare webhook ID or user assertion cannot manufacture an acceptance.

Unsupported inbound media and unknown statuses remain retained for explicit review; text with unknown contact produces the existing durable handoff and never reaches the model/tool/provider. An unresolved send remains unknown; do not provide a resend button. A new response requires a new genuine user message and a new reviewed proposal; changing text/key to force retry is prohibited.

## Empty activation inputs

Profile file/hash; fixed tenant/org/connection; purpose/policy/retention approval hash; private evidence directory; Python/source hashes from manifest; WhatsApp access-token/app-secret/verify-token sources; contact-HMAC source; OIDC verifier/service-token source; domain endpoint/token; LLM endpoint/model/token/budget; service/product/pricebook configuration; authenticated status reporter connection. These are configuration/account activation, not proof of deployment, provider reachability or live delivery.

Go integration tests use ELITE_WHATSAPP_CONNECTED_DATABASE_URL for a disposable loopback database named elite_whatsapp_connected_* and explicit ELITE_WHATSAPP_PYTHON. Test fixtures never contact Meta or the model provider.
````

### FILE: `internal/whatsappbridge/conversation_browser_issuer_test.go`

```yaml
block_id: "COMMUNICATIONS-DELTA-PYTHON_META_WHATSAPP_CLOUD_ADAPTER:file4:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "8c19b120891ee2bc4b96650ab683656e5b42e9ddfbe8edd2e0f1e6e2e68e8fdd"
variables: []
secrets_allowed: false
```

````go
package whatsappbridge

// AUTHORED synthetic RS256/JWKS fixture reused from admitted local browser tests.
import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"elite.local/enterprise/internal/platform/identity"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func whatsappBrowserIssuer(t *testing.T) (identity.Verifier, func(string, string, []string, []string) string) {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	encode := base64.RawURLEncoding.EncodeToString
	var issuer *httptest.Server
	issuer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/.well-known/openid-configuration":
			_ = json.NewEncoder(w).Encode(map[string]any{"issuer": issuer.URL, "jwks_uri": issuer.URL + "/keys", "id_token_signing_alg_values_supported": []string{"RS256"}})
		case "/keys":
			_ = json.NewEncoder(w).Encode(map[string]any{"keys": []any{map[string]any{"kty": "RSA", "use": "sig", "alg": "RS256", "kid": "local-confirmation", "n": encode(key.N.Bytes()), "e": encode(big.NewInt(int64(key.E)).Bytes())}}})
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(issuer.Close)
	verifier, err := identity.NewOIDCVerifier(context.Background(), issuer.URL, "confirmation-api")
	if err != nil {
		t.Fatal(err)
	}
	token := func(subject, tenant string, permissions, organizations []string) string {
		header, err := json.Marshal(map[string]string{"alg": "RS256", "kid": "local-confirmation", "typ": "JWT"})
		if err != nil {
			t.Fatal(err)
		}
		claims, err := json.Marshal(map[string]any{"iss": issuer.URL, "aud": "confirmation-api", "sub": subject, "iat": time.Now().Unix(), "exp": time.Now().Add(5 * time.Minute).Unix(), "tenant_id": tenant, "permissions": permissions, "organization_ids": organizations})
		if err != nil {
			t.Fatal(err)
		}
		unsigned := encode(header) + "." + encode(claims)
		digest := sha256.Sum256([]byte(unsigned))
		signature, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, digest[:])
		if err != nil {
			t.Fatal(err)
		}
		return unsigned + "." + encode(signature)
	}
	return verifier, token
}
````

### FILE: `internal/whatsappbridge/conversation_browser_test.go`

```yaml
block_id: "COMMUNICATIONS-DELTA-PYTHON_META_WHATSAPP_CLOUD_ADAPTER:file5:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "3073def90bb846be7bff37f91e7a0d1d84710267f6ca7128362346d03af432ac"
variables: []
secrets_allowed: false
```

````go
package whatsappbridge

// AUTHORED browser composition fixture. Existing conversation fixture supplies
// inbox/runtime/proposal/transport; real OIDC verifies the browser's backend token.
import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestWhatsAppHumanBrowserPostgres(t *testing.T) {
	if os.Getenv("ELITE_WHATSAPP_BROWSER") != "1" {
		t.Fatal("explicit WhatsApp browser fixture required")
	}
	f := newConnectedReplyFixture(t)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	for _, id := range []string{"wamid.browser-normal", "wamid.lost", "wamid.browser-reject"} {
		f.ingest(t, f.inbound(t, id, time.Now().Add(-time.Second), "5491112345678"))
		f.process(t)
	}
	if f.metaCalls.Load() != 0 || f.llmCalls.Load() != 6 || f.domainCalls.Load() != 3 {
		t.Fatal("proposal effects", f.metaCalls.Load(), f.llmCalls.Load(), f.domainCalls.Load())
	}
	verifier, token := whatsappBrowserIssuer(t)
	identities := map[string]any{}
	for _, name := range []string{"human", "review-only", "reader", "foreign-org"} {
		permissions := []string{"whatsapp:approve", "whatsapp:send"}
		orgs := []string{"store-1"}
		subject := "fixture-human"
		if name == "review-only" {
			permissions = []string{"whatsapp:approve"}
			subject = "review-only"
		}
		if name == "reader" {
			permissions = []string{"lead:read"}
			subject = "reader"
		}
		if name == "foreign-org" {
			orgs = []string{"store-2"}
			subject = "foreign-org"
		}
		identities[name] = map[string]any{"subject": subject, "tenantId": f.human.TenantID, "permissions": permissions, "organizations": orgs, "accessToken": token(subject, f.human.TenantID, permissions, orgs)}
	}
	mux := http.NewServeMux()
	f.module.Register(mux, verifier)
	controlToken := fmt.Sprintf("fixture-%d", time.Now().UnixNano())
	mux.HandleFunc("POST /__fixture/delivered", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Fixture-Token") != controlToken {
			w.WriteHeader(403)
			return
		}
		var body struct {
			Index int `json:"index"`
		}
		if json.NewDecoder(r.Body).Decode(&body) != nil || body.Index < 1 || body.Index > 2 {
			w.WriteHeader(400)
			return
		}
		f.ingest(t, statusBatch("delivered", fmt.Sprint(time.Now().Unix()), map[string]any{"id": fmt.Sprintf("wamid.reply.%d", body.Index)}))
		f.process(t)
		w.WriteHeader(200)
	})
	var mu sync.Mutex
	counts := map[string]int{}
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" && !strings.HasPrefix(r.URL.Path, "/__fixture/") {
			mu.Lock()
			counts[r.URL.Path]++
			mu.Unlock()
		}
		mux.ServeHTTP(w, r)
	}))
	defer api.Close()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	address := listener.Addr().String()
	listener.Close()
	target, _ := url.Parse("http://" + address)
	edge := httptest.NewTLSServer(httputil.NewSingleHostReverseProxy(target))
	defer edge.Close()
	web := os.Getenv("ELITE_WEB_ROOT")
	node := os.Getenv("ELITE_NODE_BIN")
	if !filepath.IsAbs(web) || !filepath.IsAbs(node) {
		t.Fatal("fixed web/runtime paths required")
	}
	env := []string{}
	for _, item := range os.Environ() {
		name := strings.ToUpper(strings.SplitN(item, "=", 2)[0])
		if strings.HasPrefix(name, "ELITE_") || strings.Contains(name, "DATABASE_URL") || name == "APP_BASE_URL" || name == "ENTERPRISE_API_BASE_URL" || name == "AUTH_SESSION_SECRET" {
			continue
		}
		env = append(env, item)
	}
	encoded, _ := json.Marshal(identities)
	env = append(env, "ELITE_WHATSAPP_BROWSER=1", "ELITE_WORKSPACE_E2E=enabled", "ELITE_WHATSAPP_IDENTITIES="+string(encoded), "ELITE_WEB_ROOT="+web, "ELITE_BASE_URL="+edge.URL, "APP_BASE_URL="+edge.URL, "ENTERPRISE_API_BASE_URL="+api.URL, "AUTH_SESSION_SECRET=synthetic-browser-only-"+controlToken, "ELITE_WHATSAPP_CONTROL="+api.URL+"/__fixture/delivered", "ELITE_WHATSAPP_CONTROL_TOKEN="+controlToken)
	artifacts, err := os.MkdirTemp(web, "whatsapp-browser-artifacts-")
	if err != nil {
		t.Fatal(err)
	}
	log, err := os.Create(filepath.Join(artifacts, "next.log"))
	if err != nil {
		t.Fatal(err)
	}
	defer log.Close()
	server := exec.CommandContext(ctx, node, filepath.Join(web, "node_modules", "next", "dist", "bin", "next"), "start", "--hostname", "127.0.0.1", "--port", strings.Split(address, ":")[1])
	server.Dir, server.Env = web, env
	server.Stdout, server.Stderr = log, log
	if err = server.Start(); err != nil {
		t.Fatal(err)
	}
	defer func() { server.Process.Kill(); server.Wait() }()
	client := &http.Client{Timeout: time.Second}
	ready := false
	for deadline := time.Now().Add(20 * time.Second); time.Now().Before(deadline); {
		resp, e := client.Get(target.String() + "/icon.svg")
		if e == nil {
			resp.Body.Close()
			if resp.StatusCode == 200 {
				ready = true
				break
			}
		}
		select {
		case <-ctx.Done():
			t.Fatal("startup cancelled")
		case <-time.After(100 * time.Millisecond):
		}
	}
	if !ready {
		t.Fatal("Next did not start", artifacts)
	}
	gate := filepath.Join(web, "microsoft_playwright_browser_gate")
	command := exec.CommandContext(ctx, node, filepath.Join(gate, "node_modules", "@playwright", "test", "cli.js"), "test", "tests/whatsapp-connected.spec.mjs", "--project=chromium-desktop", "--timeout=90000", "--output="+filepath.Join(artifacts, "browsers"))
	command.Dir, command.Env = gate, env
	output, err := command.CombinedOutput()
	if e := os.WriteFile(filepath.Join(artifacts, "browser.log"), output, 0600); e != nil {
		t.Fatal(e)
	}
	if err != nil || !strings.Contains(string(output), "1 passed") {
		t.Fatalf("browser %v\n%s\nartifacts=%s", err, output, artifacts)
	}
	var accepted, attempts, observations, approved, rejected int
	err = f.pool.QueryRow(ctx, `select (select count(*) from communication.outbound_delivery where tenant_id=$1 and state='accepted'),(select coalesce(sum(attempt_count),0) from communication.outbound_delivery where tenant_id=$1),(select count(*) from communication.whatsapp_status_observation where tenant_id=$1),(select count(*) from approval.request where tenant_id=$1 and state='approved'),(select count(*) from approval.request where tenant_id=$1 and state='rejected')`, f.human.TenantID).Scan(&accepted, &attempts, &observations, &approved, &rejected)
	if err != nil || accepted != 2 || attempts != 2 || observations != 2 || approved != 2 || rejected != 1 || f.metaCalls.Load() != 2 {
		t.Fatal("durable outcomes", accepted, attempts, observations, approved, rejected, f.metaCalls.Load(), err)
	}
	mu.Lock()
	countsJSON, _ := json.Marshal(counts)
	mu.Unlock()
	if err = os.WriteFile(filepath.Join(artifacts, "api-post-counts.json"), countsJSON, 0600); err != nil {
		t.Fatal(err)
	}
	t.Logf("WHATSAPP_BROWSER_POSTGRES_PASS actual_BFF_SQL=true JWE_RS256_JWKS=true runtime_proposals=3 approvals=2 rejected=1 provider_POSTs=2 accepted=2 delivered=2 lost_receipt_recovered_without_POST=true reader_navigation_hidden=true send_permission_enforced=true artifacts=%s", artifacts)
}
````

### FILE: `internal/whatsappbridge/conversation_connected_test.go`

```yaml
block_id: "COMMUNICATIONS-DELTA-PYTHON_META_WHATSAPP_CLOUD_ADAPTER:file6:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "e6ff9deac9e5e1736fe6e4abfb8ad91f0247aed3337151830cec2bd695d723fd"
variables: []
secrets_allowed: false
```

````go
package whatsappbridge

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"elite.local/enterprise/internal/app"
	"elite.local/enterprise/internal/channels"
	"elite.local/enterprise/internal/contactidentity"
	"elite.local/enterprise/internal/conversationruntime"
	"elite.local/enterprise/internal/leadstream"
	"elite.local/enterprise/internal/outbounddelivery"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

type replyDomainToken struct{}

func (replyDomainToken) AccessToken(context.Context) (string, error) {
	return "fixture-domain-token", nil
}

type replyVerifier map[string]identity.Principal

func (v replyVerifier) Verify(_ context.Context, token string) (identity.Principal, error) {
	p, ok := v[token]
	if !ok {
		return p, ErrApproval
	}
	return p, nil
}

type connectedReplyFixture struct {
	pool                             *pgxpool.Pool
	observer                         *StatusObserver
	worker                           *StatusWorker
	router                           *StatusRouter
	module                           *ReplyModule
	receiver                         *httptest.Server
	api                              *httptest.Server
	service, human                   identity.Principal
	metaCalls, llmCalls, domainCalls *atomic.Int32
}

func newConnectedReplyFixture(t *testing.T) *connectedReplyFixture {
	t.Helper()
	ctx := context.Background()
	db := os.Getenv("ELITE_WHATSAPP_CONNECTED_DATABASE_URL")
	u, err := url.Parse(db)
	if err != nil || u.Scheme != "postgres" || u.Hostname() != "127.0.0.1" || !strings.HasPrefix(u.Path, "/elite_whatsapp_connected_") {
		t.Fatal("dedicated loopback WhatsApp database required")
	}
	pool, err := pgxpool.New(ctx, db)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)
	tenant := leadstream.StableUUID(t.Name(), time.Now().Format(time.RFC3339Nano))
	key := []byte("0123456789abcdef0123456789abcdef")
	for _, sql := range []string{
		`insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name)values($1::uuid,'wa-'||$1::text,'Fixture','Fixture')`,
		`insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values($1,'store-1','store-1','Fixture','store')`,
		`insert into crm.lead(tenant_id,lead_id,organization_id,lifecycle_state,source_code,contact_payload,consent_required)values($1,'lead-1','store-1','new','fixture','{}',true)`,
		`insert into crm.consent_evidence(tenant_id,consent_id,lead_id,purpose_code,policy_version,decision,occurred_at,evidence_sha256_hex)values($1,'consent-1','lead-1','fixture-conversation','policy-1','granted',clock_timestamp()-interval '1 minute',repeat('c',64))`,
		`insert into integration.provider_connection(tenant_id,connection_id,provider_code,secret_ref,organization_id)values($1,'wa-primary','meta-whatsapp','META_APP_SECRET','store-1')`,
	} {
		if _, err = pool.Exec(ctx, sql, tenant); err != nil {
			t.Fatal(err)
		}
	}
	contacts, err := postgres.NewContactIdentityStore(pool, key)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = contacts.Apply(ctx, contactidentity.Command{TenantID: tenant, ChannelCode: "whatsapp", ExternalID: "5491112345678", LeadID: "lead-1", SubjectID: "lead:lead-1", PIIAllowed: true, State: contactidentity.StateActive, PolicyVersion: "policy-1", EvidenceSHA256: strings.Repeat("e", 64), EffectiveAt: time.Now().Add(-time.Minute), RequestID: "binding-1"}); err != nil {
		t.Fatal(err)
	}
	base, err := NewPostgresAppointmentApprovals(pool, key, "fixture-conversation", "policy-1")
	if err != nil {
		t.Fatal(err)
	}
	python := os.Getenv("ELITE_WHATSAPP_PYTHON")
	if python == "" {
		t.Fatal("fixed Python required")
	}
	actual, err := filepath.Abs("../../whatsapp_cloud")
	if err != nil {
		t.Fatal(err)
	}
	raw, err := exec.Command(python, "-I", "-B", "-c", `import sys,json;sys.path.insert(0,sys.argv[1]);from test_whatsapp_cloud import profile;print(json.dumps(profile()))`, actual).Output()
	if err != nil {
		t.Fatal("profile fixture", err)
	}
	var profile json.RawMessage = raw
	profile = append(json.RawMessage(nil), []byte(strings.TrimSpace(string(profile)))...)
	var meta, llmCount, domainCount atomic.Int32
	metaServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" || r.URL.Path != "/messages" || r.Header.Get("Authorization") != "Bearer synthetic-token" {
			t.Error("unexpected fixture Meta transport")
		}
		body, _ := io.ReadAll(r.Body)
		var m map[string]any
		_ = json.Unmarshal(body, &m)
		if m["type"] != "text" || m["to"] != "5491112345678" || m["messaging_product"] != "whatsapp" {
			t.Errorf("provider payload: %s", body)
		}
		n := meta.Add(1)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"messaging_product":"whatsapp","messages":[{"id":"wamid.reply.%d"}]}`, n)
	}))
	t.Cleanup(metaServer.Close)
	// The production owner runs unchanged inside a hash-locked fixture wrapper.
	// Its only transport maps the exact official URL to this loopback HTTP server.
	// No live network/credential is used; lost stdout happens AFTER real receipt write.
	wrap := t.TempDir()
	quote := func(s string) string { b, _ := json.Marshal(s); return string(b) }
	actualCode, _ := os.ReadFile(filepath.Join(actual, "whatsapp_cloud.py"))
	statusCode, _ := os.ReadFile(filepath.Join(actual, "status_reconciliation.py"))
	script := `import sys,json,hashlib,pathlib,importlib.util,urllib.request
root=pathlib.Path(` + quote(actual) + `)
assert hashlib.sha256((root/'whatsapp_cloud.py').read_bytes()).hexdigest()==` + quote(digest(actualCode)) + `
assert hashlib.sha256((root/'status_reconciliation.py').read_bytes()).hexdigest()==` + quote(digest(statusCode)) + `
spec=importlib.util.spec_from_file_location('whatsapp_cloud',root/'whatsapp_cloud.py');m=importlib.util.module_from_spec(spec);sys.modules['whatsapp_cloud']=m;spec.loader.exec_module(m)
raw=sys.stdin.buffer.read(13*1024*1024)
def transport(method,url,headers,body,timeout):
 assert method=='POST' and url=='https://graph.facebook.com/v99.0/123456789/messages'
 req=urllib.request.Request(` + quote(metaServer.URL+"/messages") + `,data=body,headers=headers,method=method)
 with urllib.request.urlopen(req,timeout=timeout) as r:return r.status,dict(r.headers),r.read(65537)
mode=sys.argv[1]
if mode=='--send-bridge':
 reply=m.send_bridge(raw,transport)
 if json.loads(raw)['request'].get('source_message_id')=='wamid.lost':sys.exit(2)
elif mode=='--ingress-bridge':reply=m.ingress_bridge(raw)
elif mode=='--recover-send-bridge':reply=m.recover_send_bridge(raw)
else:
 spec=importlib.util.spec_from_file_location('status_reconciliation',root/'status_reconciliation.py');q=importlib.util.module_from_spec(spec);spec.loader.exec_module(q);reply=q.status_bridge(raw)
print(json.dumps(reply))
`
	if err = os.WriteFile(filepath.Join(wrap, "whatsapp_cloud.py"), []byte(script), 0600); err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(filepath.Join(wrap, "status_reconciliation.py"), statusCode, 0600); err != nil {
		t.Fatal(err)
	}
	pythonCode, _ := os.ReadFile(python)
	process := Process{PythonExecutable: python, PythonSHA256: digest(pythonCode), AdapterDirectory: wrap, AdapterSHA256: digest([]byte(script)), EvidenceDirectory: t.TempDir()}
	observer := &StatusObserver{TenantID: tenant, Profile: profile, Process: process, ReconcilerSHA256: digest(statusCode), Secrets: appSecretFixture{}, Approvals: base}
	replies, err := NewPostgresReplyApprovals(base, tenant, "store-1", "wa-primary", profile)
	if err != nil {
		t.Fatal(err)
	}
	store, err := postgres.NewOutboundDeliveryStore(pool, key, time.Second)
	if err != nil {
		t.Fatal(err)
	}
	sender := &Sender{TenantID: tenant, Profile: profile, Process: process, Tokens: tokenFixture{}}
	module, err := NewReplyModule(replies, sender, store)
	if err != nil {
		t.Fatal(err)
	}
	service := identity.Principal{Subject: "fixture-service", TenantID: tenant, Organizations: map[string]struct{}{"store-1": {}}, Permissions: map[string]struct{}{"appointment:manage": {}, "whatsapp:process": {}}}
	human := identity.Principal{Subject: "fixture-human", TenantID: tenant, Organizations: map[string]struct{}{"store-1": {}}, Permissions: map[string]struct{}{"whatsapp:approve": {}, "whatsapp:send": {}}}
	llm := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := llmCount.Add(1)
		w.Header().Set("Content-Type", "application/json")
		if n%2 == 1 {
			fmt.Fprintf(w, `{"id":"resp%d","status":"completed","output":[{"type":"function_call","call_id":"call%d","name":"create_quote","arguments":"{\"product\":\"scooter\",\"quantity\":1}"}],"usage":{"input_tokens":12,"output_tokens":4,"total_tokens":16}}`, n, n)
			return
		}
		fmt.Fprintf(w, `{"id":"resp%d","status":"completed","output":[{"type":"message","content":[{"type":"output_text","text":"Cotización preparada para revisar."}]}],"usage":{"input_tokens":8,"output_tokens":5,"total_tokens":13}}`, n)
	}))
	t.Cleanup(llm.Close)
	domain := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		domainCount.Add(1)
		if r.URL.Path != "/v1/franchise/quotes" || r.Header.Get("Authorization") != "Bearer fixture-domain-token" || len(r.Header.Get("Idempotency-Key")) != 64 {
			t.Error("unbound domain tool")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"quotation_id":"fixture-quote"}`))
	}))
	t.Cleanup(domain.Close)
	turns, err := postgres.NewConversationStore(pool, time.Minute, 48*time.Hour, 3)
	if err != nil {
		t.Fatal(err)
	}
	registry := channels.NewRegistry()
	_ = registry.Register(VerifiedInboxChannel{})
	runtime, err := app.New(app.Config{LLMAPIKey: "synthetic-llm-key", LLMModel: "fixture-model", LLMBaseURL: llm.URL, DomainBaseURL: domain.URL, DomainTokenProvider: replyDomainToken{}, TenantID: tenant, TenantCode: "fixture", ServiceKinds: map[string]string{"consulta": "consultation"}, ProductVariants: map[string]string{"scooter": "scooter-1"}, PriceBookID: "fixture-book", ConversationStore: turns, ContactResolver: ScopedContactResolver{TenantID: tenant, OrganizationID: "store-1", Resolver: contacts}, ChannelRegistry: registry, LLMTokenBudget: 100000, Conversation: conversationruntime.Config{Instructions: "Synthetic fixture; propose and use only admitted tools.", PromptCacheKey: "fixture-wa", MaxOutputTokens: 128, MaxHistory: 4, MaxInputBytes: 16384, MaxToolBytes: 16384, StoreApproved: true, HandoffText: "Requiere revisión humana."}})
	if err != nil {
		t.Fatal(err)
	}
	router, err := NewStatusRouter(observer, "wa-primary", strings.Repeat("a", 64), key, 8)
	if err != nil {
		t.Fatal(err)
	}
	if err = router.EnableConversation(ConversationRoute{OrganizationID: "store-1", Runtime: runtime.Conversation, Proposals: replies}); err != nil {
		t.Fatal(err)
	}
	worker, err := NewStatusWorker(router, "wa-connected", 2*time.Minute, 0)
	if err != nil {
		t.Fatal(err)
	}
	receiver, err := NewWebhookReceiver(receiverConfig(observer))
	if err != nil {
		t.Fatal(err)
	}
	webhook := httptest.NewServer(receiver)
	t.Cleanup(webhook.Close)
	mux := http.NewServeMux()
	module.Register(mux, replyVerifier{"human": human, "service": service})
	api := httptest.NewServer(mux)
	t.Cleanup(api.Close)
	return &connectedReplyFixture{pool: pool, observer: observer, worker: worker, router: router, module: module, receiver: webhook, api: api, service: service, human: human, metaCalls: &meta, llmCalls: &llmCount, domainCalls: &domainCount}
}
func (f *connectedReplyFixture) inbound(t *testing.T, id string, at time.Time, recipient string) SignedStatusWebhook {
	t.Helper()
	return alterStatusBatch(t, statusBatch("sent", fmt.Sprint(at.Unix()), nil), func(body map[string]any) {
		v := statusValue(body)
		delete(v, "statuses")
		v["messages"] = []any{map[string]any{"id": id, "from": recipient, "type": "text", "timestamp": fmt.Sprint(at.Unix()), "text": map[string]any{"body": "Cotizame un scooter"}}}
	})
}
func (f *connectedReplyFixture) ingest(t *testing.T, b SignedStatusWebhook) {
	t.Helper()
	res, _ := webhookRequest(t, f.receiver, b)
	if res.StatusCode != 200 {
		t.Fatalf("ingress: %d", res.StatusCode)
	}
}
func (f *connectedReplyFixture) process(t *testing.T) {
	t.Helper()
	r, err := f.worker.ProcessOnce(context.Background(), f.service)
	if err != nil || !r.Completed {
		t.Fatalf("worker %+v %v", r, err)
	}
}
func (f *connectedReplyFixture) command(t *testing.T, key, op, token string, body any) int {
	t.Helper()
	raw, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", f.api.URL+"/v1/franchise/whatsapp/replies/"+key+"/"+op, strings.NewReader(string(raw)))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	res, err := f.api.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	b, _ := io.ReadAll(res.Body)
	t.Logf("%s %d %s", op, res.StatusCode, string(b))
	return res.StatusCode
}
func (f *connectedReplyFixture) proposal(t *testing.T, id string) ReplyProposal {
	t.Helper()
	v, err := f.module.approvals.Read(context.Background(), f.human, ReplyDeliveryKey(f.service.TenantID, "wa-primary", id))
	if err != nil {
		t.Fatal(err)
	}
	return v
}
func TestWhatsAppConnectedHumanReplyAndStatus(t *testing.T) {
	f := newConnectedReplyFixture(t)
	b := f.inbound(t, "wamid.inbound", time.Now().Add(-time.Second), "5491112345678")
	f.ingest(t, b)
	f.ingest(t, b)
	f.process(t)
	proposal := f.proposal(t, "wamid.inbound")
	if proposal.State != "pending" || f.metaCalls.Load() != 0 || f.llmCalls.Load() != 2 || f.domainCalls.Load() != 1 {
		t.Fatal("runtime did not stop at proposal", proposal.State, f.metaCalls.Load(), f.llmCalls.Load(), f.domainCalls.Load())
	}
	if f.command(t, proposal.RequestID, "send", "human", map[string]any{"payload_sha256": proposal.PayloadSHA256}) != 409 {
		t.Fatal("pending sent")
	}
	if f.command(t, proposal.RequestID, "decision", "human", map[string]any{"payload_sha256": strings.Repeat("f", 64), "approved": true, "reason": "review"}) != 409 {
		t.Fatal("changed hash approved")
	}
	if f.command(t, proposal.RequestID, "decision", "human", map[string]any{"payload_sha256": proposal.PayloadSHA256, "approved": true, "reason": "Reviewed exact response"}) != 200 {
		t.Fatal("approval failed")
	}
	for i := 0; i < 2; i++ {
		if f.command(t, proposal.RequestID, "send", "human", map[string]any{"payload_sha256": proposal.PayloadSHA256}) != 200 {
			t.Fatal("send/replay failed")
		}
	}
	if f.metaCalls.Load() != 1 {
		t.Fatal("provider send duplicated")
	}
	status := statusBatch("delivered", fmt.Sprint(time.Now().Unix()), map[string]any{"id": "wamid.reply.1"})
	f.ingest(t, status)
	f.process(t)
	f.ingest(t, status)
	result, err := f.worker.ProcessOnce(context.Background(), f.service)
	if err != nil || result.Claimed {
		t.Fatal("status replay reprocessed", result, err)
	}
	var observations, attempts int
	var state string
	if err = f.pool.QueryRow(context.Background(), `select state,attempt_count from communication.outbound_delivery where tenant_id=$1 and delivery_key=$2`, f.service.TenantID, proposal.RequestID).Scan(&state, &attempts); err != nil {
		t.Fatal(err)
	}
	if err = f.pool.QueryRow(context.Background(), `select count(*) from communication.whatsapp_status_observation where tenant_id=$1`, f.service.TenantID).Scan(&observations); err != nil {
		t.Fatal(err)
	}
	read := f.proposal(t, "wamid.inbound")
	if read.Status == nil || read.Status.FenceState != "accepted" || read.Status.DeliveryStatus != "observed_delivered" {
		t.Fatal("human status projection missing", read.Status)
	}
	if state != "accepted" || attempts != 1 || observations != 1 || f.domainCalls.Load() != 1 {
		t.Fatal("unexpected connected effects", state, attempts, observations)
	}
}
func TestWhatsAppConnectedLostResponseRecoveryNoSecondPost(t *testing.T) {
	f := newConnectedReplyFixture(t)
	f.ingest(t, f.inbound(t, "wamid.lost", time.Now().Add(-time.Second), "5491112345678"))
	f.process(t)
	v := f.proposal(t, "wamid.lost")
	if f.command(t, v.RequestID, "decision", "human", map[string]any{"payload_sha256": v.PayloadSHA256, "approved": true, "reason": "Reviewed"}) != 200 {
		t.Fatal("approval")
	}
	for i := 0; i < 2; i++ {
		if f.command(t, v.RequestID, "send", "human", map[string]any{"payload_sha256": v.PayloadSHA256}) != 409 {
			t.Fatal("uncertain accepted")
		}
	}
	if f.metaCalls.Load() != 1 {
		t.Fatal("expected exactly one fixture POST", f.metaCalls.Load())
	}
	if f.command(t, v.RequestID, "recover", "human", map[string]any{"payload_sha256": v.PayloadSHA256}) != 200 {
		t.Fatal("evidence recovery")
	}
	f.ingest(t, statusBatch("read", fmt.Sprint(time.Now().Unix()), map[string]any{"id": "wamid.reply.1"}))
	f.process(t)
	if f.metaCalls.Load() != 1 {
		t.Fatal("reconciliation sent again")
	}
}
func TestWhatsAppConnectedRejectWindowAndUnboundContact(t *testing.T) {
	f := newConnectedReplyFixture(t)
	f.ingest(t, f.inbound(t, "wamid.reject", time.Now().Add(-time.Second), "5491112345678"))
	f.process(t)
	v := f.proposal(t, "wamid.reject")
	if f.command(t, v.RequestID, "decision", "human", map[string]any{"payload_sha256": v.PayloadSHA256, "approved": false, "reason": "Rejected"}) != 200 {
		t.Fatal("reject")
	}
	if f.command(t, v.RequestID, "decision", "human", map[string]any{"payload_sha256": v.PayloadSHA256, "approved": true, "reason": "Retry approval"}) != 409 {
		t.Fatal("reject mutated")
	}
	if f.command(t, v.RequestID, "send", "human", map[string]any{"payload_sha256": v.PayloadSHA256}) != 409 {
		t.Fatal("reject sent")
	}
	f.ingest(t, f.inbound(t, "wamid.unknown", time.Now().Add(-time.Second), "5491199999999"))
	f.process(t)
	if f.llmCalls.Load() != 2 || f.domainCalls.Load() != 1 || f.metaCalls.Load() != 0 {
		t.Fatal("unbound contact reached model/tool/provider")
	}
	var handed bool
	if err := f.pool.QueryRow(context.Background(), `select state='handed_off' and failure_code='CONTACT_SCOPE_UNAVAILABLE' from communication.conversation_turn where tenant_id=$1 and provider_message_id='wamid.unknown'`, f.service.TenantID).Scan(&handed); err != nil || !handed {
		t.Fatal("missing durable handoff", err)
	}
	// Local clock boundary for expired text is also enforced independently before
	// the provider process, even with an otherwise exact typed message.
	expired := v.Context.Message
	var request ReplyRequest
	_ = json.Unmarshal([]byte(expired.Text), &request)
	request.LastInboundAt = time.Now().Add(-25 * time.Hour).Unix()
	request.WindowExpiresAt = request.LastInboundAt + 86400
	raw, _ := json.Marshal(request)
	expired.Text = string(raw)
	if _, err := replyRequestFor(expired); err == nil {
		t.Fatal("closed window admitted")
	}
}

// Failure injection at the existing durable Store boundary: provider returns a
// receipt, then process loss prevents Complete. No SQL fabricates a send state.
type replyCrashStore struct {
	*postgres.OutboundDeliveryStore
}

func (replyCrashStore) Complete(context.Context, channels.Message, string, outbounddelivery.Receipt) error {
	return outbounddelivery.ErrUnknown
}
func TestWhatsAppConnectedLeaseCrashAndMissingEvidence(t *testing.T) {
	f := newConnectedReplyFixture(t)
	f.ingest(t, f.inbound(t, "wamid.crash", time.Now().Add(-time.Second), "5491112345678"))
	f.process(t)
	v := f.proposal(t, "wamid.crash")
	if f.command(t, v.RequestID, "decision", "human", map[string]any{"payload_sha256": v.PayloadSHA256, "approved": true, "reason": "Reviewed"}) != 200 {
		t.Fatal("approval")
	}
	f.module.channel.Store = replyCrashStore{f.module.store}
	if f.command(t, v.RequestID, "send", "human", map[string]any{"payload_sha256": v.PayloadSHA256}) != 409 {
		t.Fatal("crash not surfaced")
	}
	dir := filepath.Join(f.module.sender.Process.EvidenceDirectory, digest([]byte(f.service.TenantID+"\x00"+v.RequestID)))
	receipt := filepath.Join(dir, "SEND_RECEIPT.json")
	saved, err := os.ReadFile(receipt)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(receipt, []byte(`{}`), 0600); err != nil {
		t.Fatal(err)
	}
	if f.command(t, v.RequestID, "recover", "human", map[string]any{"payload_sha256": v.PayloadSHA256}) != 409 {
		t.Fatal("unproved receipt accepted")
	}
	if err = os.WriteFile(receipt, saved, 0600); err != nil {
		t.Fatal(err)
	}
	time.Sleep(1100 * time.Millisecond)
	for i := 0; i < 2; i++ {
		if f.command(t, v.RequestID, "recover", "human", map[string]any{"payload_sha256": v.PayloadSHA256}) != 200 {
			t.Fatal("expired lease recovery/replay failed")
		}
	}
	if f.metaCalls.Load() != 1 {
		t.Fatal("crash recovery duplicated POST")
	}
}

func TestWhatsAppFutureContactCannotReachRuntime(t *testing.T) {
	f := newConnectedReplyFixture(t)
	binding, err := postgres.NewContactIdentityStore(f.pool, []byte("0123456789abcdef0123456789abcdef"))
	if err != nil {
		t.Fatal(err)
	}
	_, err = binding.Apply(context.Background(), contactidentity.Command{TenantID: f.service.TenantID, ChannelCode: "whatsapp", ExternalID: "5491112345678", LeadID: "lead-1", SubjectID: "lead:lead-1", PIIAllowed: true, State: contactidentity.StateActive, PolicyVersion: "policy-1", EvidenceSHA256: strings.Repeat("e", 64), EffectiveAt: time.Now().Add(time.Hour), ExpectedVersion: 1, RequestID: "binding-future"})
	if err != nil {
		t.Fatal(err)
	}
	f.ingest(t, f.inbound(t, "wamid.future", time.Now().Add(-time.Second), "5491112345678"))
	result, err := f.worker.ProcessOnce(context.Background(), f.service)
	if f.llmCalls.Load() != 0 || f.domainCalls.Load() != 0 || f.metaCalls.Load() != 0 {
		t.Fatalf("future contact escaped before effective_at: model=%d domain=%d provider=%d", f.llmCalls.Load(), f.domainCalls.Load(), f.metaCalls.Load())
	}
	if err != nil || !result.Completed {
		t.Fatal("future contact should produce durable handoff", result, err)
	}
}

func TestWhatsAppReplyApprovalCannotSurviveScopeOrBindingDrift(t *testing.T) {
	f := newConnectedReplyFixture(t)
	ctx := context.Background()
	f.ingest(t, f.inbound(t, "wamid.binding", time.Now().Add(-time.Second), "5491112345678"))
	f.process(t)
	v := f.proposal(t, "wamid.binding")
	self := f.service
	self.Permissions = map[string]struct{}{"whatsapp:approve": {}}
	if _, err := f.module.approvals.Decide(ctx, self, v.RequestID, v.PayloadSHA256, true, "self review"); err == nil {
		t.Fatal("service requester approved own reply")
	}
	foreign := f.human
	foreign.Organizations = map[string]struct{}{"other": {}}
	if _, err := f.module.approvals.Read(ctx, foreign, v.RequestID); err == nil {
		t.Fatal("foreign organization read exact reply")
	}
	if f.command(t, v.RequestID, "decision", "human", map[string]any{"payload_sha256": v.PayloadSHA256, "approved": true, "reason": "Reviewed"}) != 200 {
		t.Fatal("approval")
	}
	store, err := postgres.NewContactIdentityStore(f.pool, []byte("0123456789abcdef0123456789abcdef"))
	if err != nil {
		t.Fatal(err)
	}
	_, err = store.Apply(ctx, contactidentity.Command{TenantID: f.service.TenantID, ChannelCode: "whatsapp", ExternalID: "5491112345678", LeadID: "lead-1", SubjectID: "lead:lead-1", PIIAllowed: false, State: contactidentity.StateRevoked, PolicyVersion: "policy-1", EvidenceSHA256: strings.Repeat("e", 64), EffectiveAt: time.Now().Add(-time.Second), ExpectedVersion: 1, RequestID: "revoke-contact"})
	if err != nil {
		t.Fatal(err)
	}
	if f.command(t, v.RequestID, "send", "human", map[string]any{"payload_sha256": v.PayloadSHA256}) != 409 {
		t.Fatal("revoked binding sent")
	}
	_, err = store.Apply(ctx, contactidentity.Command{TenantID: f.service.TenantID, ChannelCode: "whatsapp", ExternalID: "5491112345678", LeadID: "lead-1", SubjectID: "lead:lead-1", PIIAllowed: true, State: contactidentity.StateActive, PolicyVersion: "policy-1", EvidenceSHA256: strings.Repeat("e", 64), EffectiveAt: time.Now().Add(-time.Second), ExpectedVersion: 2, RequestID: "rebind-contact"})
	if err != nil {
		t.Fatal(err)
	}
	if f.command(t, v.RequestID, "send", "human", map[string]any{"payload_sha256": v.PayloadSHA256}) != 409 {
		t.Fatal("new binding silently renewed old approval")
	}
	if f.metaCalls.Load() != 0 {
		t.Fatal("drift reached provider")
	}
}
````

### FILE: `internal/whatsappbridge/conversation_projection_fuzz_test.go`

```yaml
block_id: "COMMUNICATIONS-DELTA-PYTHON_META_WHATSAPP_CLOUD_ADAPTER:file7:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "2036266d0df0c779dd30afeaf1bfd05a0acbe9c9029921ad4ebcb7dd6717f565"
variables: []
secrets_allowed: false
```

````go
package whatsappbridge

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"
	"unicode/utf8"
)

// Native fuzz covers the bounded typed projection AFTER the separate signed
// Python trust boundary. This does not claim signature verification or DAST.
func FuzzWhatsAppConversationProjection(f *testing.F) {
	seed := fmt.Sprintf(`{"object":"whatsapp_business_account","entry":[{"id":"987654321","changes":[{"field":"messages","value":{"messaging_product":"whatsapp","metadata":{"phone_number_id":"123456789"},"messages":[{"id":"wamid.fixture","from":"5491112345678","timestamp":"%d","type":"text","text":{"body":"Cotizame un scooter"}}]}}]}]}`, time.Now().Add(-time.Minute).Unix())
	for _, v := range []string{seed, `{}`, `null`, strings.Replace(seed, `"type":"text"`, `"type":"image"`, 1), strings.Replace(seed, `"body":"Cotizame un scooter"`, `"body":null`, 1)} {
		f.Add([]byte(v))
	}
	f.Fuzz(func(t *testing.T, raw []byte) {
		if len(raw) > 16384 {
			return
		}
		messages, total, err := verifiedConversationMessages(raw, "fixture-tenant", "fixture-connection", 8)
		if err != nil {
			return
		}
		if len(messages) > 8 || total < len(messages) {
			t.Fatal("projection escaped route budget")
		}
		seen := map[string]bool{}
		for _, m := range messages {
			if m.TenantID != "fixture-tenant" || m.ChannelCode != "whatsapp" || m.ThreadID != conversationThread("fixture-connection") || m.Direction != "in" || !utf8.ValidString(m.Text) || m.Validate() != nil || seen[m.ProviderMessageID] {
				t.Fatal("projection escaped source scope or identity")
			}
			seen[m.ProviderMessageID] = true
			encoded, err := json.Marshal(m)
			if err != nil {
				t.Fatal(err)
			}
			var again map[string]any
			if json.Unmarshal(encoded, &again) != nil || again["Text"] != m.Text || again["ProviderMessageID"] != m.ProviderMessageID {
				t.Fatal("projection JSON changed visible text/identity")
			}
		}
	})
}
````

### FILE: `internal/whatsappbridge/conversation_router.go`

```yaml
block_id: "COMMUNICATIONS-DELTA-PYTHON_META_WHATSAPP_CLOUD_ADAPTER:file8:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "252b0ed658c23c5dc2ec5da01adabaebe65586ba5fc343f279ea8d3ef6801c2d"
variables: []
secrets_allowed: false
```

````go
package whatsappbridge

// AUTHORED typed projection after the existing Meta signature/scope verifier.
// The existing conversation runtime stores a proposal; this path never sends it.
import (
	"context"
	"encoding/json"
	"strconv"
	"time"
	"unicode/utf8"

	"elite.local/enterprise/internal/channels"
	"elite.local/enterprise/internal/conversationruntime"
	"elite.local/enterprise/internal/platform/identity"
)

type ConversationHandler interface {
	Handle(context.Context, channels.Message) (string, error)
}
type ReplyProposer interface {
	Propose(context.Context, identity.Principal, string, channels.Message) (ReplyProposal, error)
}

type ConversationRoute struct {
	OrganizationID string
	Runtime        ConversationHandler
	Proposals      ReplyProposer
}

// EnableConversation is explicit composition; its dependencies cannot be
// supplied by a webhook. The organization must match the durable connection.
func (r *StatusRouter) EnableConversation(route ConversationRoute) error {
	if r == nil || route.OrganizationID == "" || route.Runtime == nil || route.Proposals == nil {
		return ErrStatusRouting
	}
	r.conversation = &route
	return nil
}

// ScopedContactResolver prevents a tenant-wide contact binding from crossing
// the connection's selected organization before tools or model execution.
type ScopedContactResolver struct {
	TenantID, OrganizationID string
	Resolver                 conversationruntime.ContactResolver
}

func (s ScopedContactResolver) Resolve(ctx context.Context, m channels.Message) (conversationruntime.Contact, error) {
	if s.Resolver == nil || s.TenantID == "" || s.OrganizationID == "" || m.TenantID != s.TenantID || m.ChannelCode != "whatsapp" {
		return conversationruntime.Contact{}, ErrStatusRouting
	}
	c, err := s.Resolver.Resolve(ctx, m)
	if err != nil || c.Scope.OrganizationID != s.OrganizationID {
		return conversationruntime.Contact{}, ErrStatusRouting
	}
	return c, nil
}

func conversationThread(connection string) string {
	return "wa-connection:" + digest([]byte(connection))
}

// Raw has already passed Python's exact-key, HMAC, scope and event validation.
// This projection still validates supported message shape before any runtime
// effect. Unsupported media remains in the inbox for explicit handoff/review.
func verifiedConversationMessages(raw []byte, tenant, connection string, limit int) ([]channels.Message, int, error) {
	var envelope map[string]json.RawMessage
	if json.Unmarshal(raw, &envelope) != nil {
		return nil, 0, ErrStatusRouting
	}
	var entries []map[string]json.RawMessage
	if json.Unmarshal(envelope["entry"], &entries) != nil {
		return nil, 0, ErrStatusRouting
	}
	var out []channels.Message
	seen := map[string]string{}
	total := 0
	for _, entry := range entries {
		var changes []map[string]json.RawMessage
		if json.Unmarshal(entry["changes"], &changes) != nil {
			return nil, 0, ErrStatusRouting
		}
		for _, change := range changes {
			var value map[string]json.RawMessage
			if json.Unmarshal(change["value"], &value) != nil {
				return nil, 0, ErrStatusRouting
			}
			var messages []map[string]json.RawMessage
			if v, ok := value["messages"]; ok && json.Unmarshal(v, &messages) != nil {
				return nil, 0, ErrStatusRouting
			}
			for _, m := range messages {
				total++
				var id, from, stamp, kind string
				if json.Unmarshal(m["id"], &id) != nil || json.Unmarshal(m["from"], &from) != nil || json.Unmarshal(m["timestamp"], &stamp) != nil || json.Unmarshal(m["type"], &kind) != nil || kind != "text" {
					return nil, 0, ErrStatusRouting
				}
				var body map[string]json.RawMessage
				var text string
				if json.Unmarshal(m["text"], &body) != nil || json.Unmarshal(body["body"], &text) != nil || !utf8.ValidString(text) || len(text) == 0 || len(text) > 16384 {
					return nil, 0, ErrStatusRouting
				}
				unix, err := strconv.ParseInt(stamp, 10, 64)
				if err != nil || unix < 1 || strconv.FormatInt(unix, 10) != stamp || time.Unix(unix, 0).After(time.Now().Add(time.Minute)) {
					return nil, 0, ErrStatusRouting
				}
				message := channels.Message{ChannelCode: "whatsapp", TenantID: tenant, ExternalID: from, ThreadID: conversationThread(connection), ProviderMessageID: id, OccurredAt: time.Unix(unix, 0).UTC(), Direction: channels.DirectionIn, Text: text}
				if message.Validate() != nil {
					return nil, 0, ErrStatusRouting
				}
				encoded, _ := json.Marshal(message)
				hash := digest(encoded)
				if previous, ok := seen[id]; ok {
					if previous != hash {
						return nil, 0, ErrStatusRouting
					}
					continue
				}
				seen[id] = hash
				out = append(out, message)
				if len(out) > limit {
					return nil, 0, ErrStatusRouting
				}
			}
		}
	}
	return out, total, nil
}
````

### FILE: `internal/whatsappbridge/reply_approval.go`

```yaml
block_id: "COMMUNICATIONS-DELTA-PYTHON_META_WHATSAPP_CLOUD_ADAPTER:file9:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "09d9fa5f22e52ff85dd7003461c35a02775c7e5c96b89a7226a280910e607f47"
variables: []
secrets_allowed: false
```

````go
package whatsappbridge

// AUTHORED typed projection of conversation/contact/consent and the shared
// durable approval owner. It defines no consent or business approval policy.
import (
	"context"
	"encoding/json"
	"time"

	"elite.local/enterprise/internal/approval"
	"elite.local/enterprise/internal/channels"
	"elite.local/enterprise/internal/contactidentity"
	"elite.local/enterprise/internal/outbounddelivery"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/postgres"
	"github.com/jackc/pgx/v5"
)

type ReplyContext struct {
	Schema                string           `json:"schema"`
	OrganizationID        string           `json:"organization_id"`
	ConnectionID          string           `json:"connection_id"`
	SourceEventID         string           `json:"source_event_id"`
	ProviderMessageID     string           `json:"provider_message_id"`
	ExternalIDHMAC        string           `json:"external_id_hmac"`
	BindingVersion        int64            `json:"binding_version"`
	LeadID                string           `json:"lead_id"`
	SubjectID             string           `json:"subject_id"`
	PolicyVersion         string           `json:"policy_version"`
	ConsentID             string           `json:"consent_id"`
	ConsentPurpose        string           `json:"consent_purpose"`
	ConsentEvidenceSHA256 string           `json:"consent_evidence_sha256"`
	ProfileSHA256         string           `json:"profile_sha256"`
	MessageSHA256         string           `json:"message_sha256"`
	ExpiresAt             time.Time        `json:"expires_at"`
	Message               channels.Message `json:"message"`
}

type ReplyProposal struct {
	Status        *NotificationStatus `json:"status,omitempty"`
	RequestID     string              `json:"request_id"`
	State         string              `json:"state"`
	PayloadSHA256 string              `json:"payload_sha256"`
	Context       ReplyContext        `json:"context"`
}

type PostgresReplyApprovals struct {
	base                                      *PostgresAppointmentApprovals
	tenant, organization, connection, profile string
}

func NewPostgresReplyApprovals(base *PostgresAppointmentApprovals, tenant, organization, connection string, profile json.RawMessage) (*PostgresReplyApprovals, error) {
	if base == nil || base.pool == nil || !notificationEventID.MatchString(tenant) || organization == "" || !webhookConnection.MatchString(connection) || !json.Valid(profile) || len(profile) > 32768 {
		return nil, ErrApproval
	}
	return &PostgresReplyApprovals{base: base, tenant: tenant, organization: organization, connection: connection, profile: digest(profile)}, nil
}
func ReplyDeliveryKey(tenant, connection, sourceMessage string) string {
	return "wa-reply:" + digest([]byte(tenant+"\x00"+connection+"\x00"+sourceMessage))
}

func (s *PostgresReplyApprovals) validPrincipal(p identity.Principal, permission string) bool {
	return s != nil && p.TenantID == s.tenant && p.Subject != "" && len(p.Subject) <= 128 && p.Allowed(permission) && p.AllowedOrganization(s.organization)
}

// buildProposal is called only from the verified raw-inbox router. It projects
// the exact terminal conversation turn and current explicit contact/consent.
func (s *PostgresReplyApprovals) buildProposal(ctx context.Context, p identity.Principal, eventID string, in channels.Message) (ReplyContext, error) {
	var c ReplyContext
	if !s.validPrincipal(p, "whatsapp:process") || in.TenantID != s.tenant || in.ChannelCode != "whatsapp" || in.Direction != channels.DirectionIn || in.ThreadID != conversationThread(s.connection) || in.Validate() != nil {
		return c, ErrApproval
	}
	recipient, err := contactidentity.ExternalDigest(s.base.hmacKey, s.tenant, "whatsapp", in.ExternalID)
	if err != nil {
		return c, ErrApproval
	}
	var text, responseHash string
	err = s.base.pool.QueryRow(ctx, `select t.assistant_text,t.response_sha256_hex,b.version,b.lead_id,b.subject_id,b.policy_version,consent.consent_id,consent.evidence_sha256_hex
 from communication.conversation_turn t
 join communication.contact_channel_binding b on b.tenant_id=t.tenant_id and b.channel_code=t.channel_code and b.external_id_hmac=$7
 join crm.lead l on l.tenant_id=b.tenant_id and l.lead_id=b.lead_id and l.organization_id=$8
 join integration.provider_connection pc on pc.tenant_id=t.tenant_id and pc.connection_id=$9 and pc.provider_code='meta-whatsapp' and pc.organization_id=l.organization_id and pc.state='active'
 join integration.webhook_event inbox on inbox.tenant_id=pc.tenant_id and inbox.connection_id=pc.connection_id and inbox.provider_code=pc.provider_code and inbox.provider_event_id=$10 and inbox.event_type='whatsapp.raw_webhook.received.v1' and inbox.state in ('received','processed')
 join crm.consent_evidence consent on consent.tenant_id=b.tenant_id and consent.lead_id=b.lead_id and consent.purpose_code=$11 and consent.policy_version=b.policy_version and consent.decision='granted' and consent.occurred_at<=statement_timestamp()
 where t.tenant_id=$1 and t.channel_code='whatsapp' and t.provider_message_id=$2 and t.external_id=$3 and t.thread_id=$4 and t.occurred_at=$5 and t.user_text=$6
 and t.state in ('completed','handed_off') and t.expires_at>statement_timestamp() and t.occurred_at<=statement_timestamp() and t.occurred_at+interval '24 hours'>statement_timestamp()+interval '40 seconds'
 and b.state='active' and b.pii_allowed and b.effective_at<=statement_timestamp() and b.policy_version=$12
 and not exists(select 1 from crm.consent_evidence newer where newer.tenant_id=consent.tenant_id and newer.lead_id=consent.lead_id and newer.purpose_code=consent.purpose_code and newer.consent_id<>consent.consent_id and newer.occurred_at>=consent.occurred_at and newer.occurred_at<=statement_timestamp())`, s.tenant, in.ProviderMessageID, in.ExternalID, in.ThreadID, in.OccurredAt, in.Text, recipient, s.organization, s.connection, eventID, s.base.purpose, s.base.policy).Scan(&text, &responseHash, &c.BindingVersion, &c.LeadID, &c.SubjectID, &c.PolicyVersion, &c.ConsentID, &c.ConsentEvidenceSHA256)
	if err != nil || digest([]byte(text)) != responseHash {
		return ReplyContext{}, ErrApproval
	}
	request := ReplyRequest{Kind: "text_reply", Recipient: in.ExternalID, Text: text, SourceMessageID: in.ProviderMessageID, LastInboundAt: in.OccurredAt.Unix(), WindowExpiresAt: in.OccurredAt.Unix() + 86400}
	raw, _ := json.Marshal(request)
	c.Message = channels.Message{ChannelCode: "whatsapp", TenantID: s.tenant, ExternalID: in.ExternalID, ThreadID: in.ThreadID, Direction: channels.DirectionOut, DeliveryKey: ReplyDeliveryKey(s.tenant, s.connection, in.ProviderMessageID), Text: string(raw)}
	if _, err = replyRequestFor(c.Message); err != nil {
		return ReplyContext{}, ErrApproval
	}
	c.MessageSHA256, err = outbounddelivery.MessageSHA256(c.Message)
	if err != nil {
		return ReplyContext{}, ErrApproval
	}
	c.Schema = "elite-whatsapp-reply-approval/v1"
	c.OrganizationID = s.organization
	c.ConnectionID = s.connection
	c.SourceEventID = eventID
	c.ProviderMessageID = in.ProviderMessageID
	c.ExternalIDHMAC = recipient
	c.ConsentPurpose = s.base.purpose
	c.ProfileSHA256 = s.profile
	c.ExpiresAt = time.Unix(request.WindowExpiresAt, 0).UTC()
	return c, nil
}

// current binds approval to the current tenant/org/connection, source turn,
// contact version, explicit latest consent and unexpired customer-service window.
// It is a pre-dispatch snapshot, not a promise to recall an in-flight request.
func (s *PostgresReplyApprovals) current(ctx context.Context, db statusQuery, c ReplyContext) error {
	if s == nil || c.Schema != "elite-whatsapp-reply-approval/v1" || c.OrganizationID != s.organization || c.ConnectionID != s.connection || c.ProfileSHA256 != s.profile || c.PolicyVersion != s.base.policy || c.ConsentPurpose != s.base.purpose || c.Message.TenantID != s.tenant || c.Message.DeliveryKey != ReplyDeliveryKey(s.tenant, s.connection, c.ProviderMessageID) || c.Message.ThreadID != conversationThread(s.connection) {
		return ErrApproval
	}
	hash, err := outbounddelivery.MessageSHA256(c.Message)
	if err != nil || hash != c.MessageSHA256 {
		return ErrApproval
	}
	recipient, err := contactidentity.ExternalDigest(s.base.hmacKey, s.tenant, "whatsapp", c.Message.ExternalID)
	if err != nil || recipient != c.ExternalIDHMAC {
		return ErrApproval
	}
	raw, err := replyRequestFor(c.Message)
	if err != nil {
		return ErrApproval
	}
	var request ReplyRequest
	_ = json.Unmarshal(raw, &request)
	if request.SourceMessageID != c.ProviderMessageID || !c.ExpiresAt.Equal(time.Unix(request.WindowExpiresAt, 0)) {
		return ErrApproval
	}
	var found bool
	err = db.QueryRow(ctx, `select true from communication.contact_channel_binding b
 join crm.lead l on l.tenant_id=b.tenant_id and l.lead_id=b.lead_id and l.organization_id=$3
 join integration.provider_connection pc on pc.tenant_id=b.tenant_id and pc.connection_id=$4 and pc.provider_code='meta-whatsapp' and pc.organization_id=l.organization_id and pc.state='active'
 join communication.conversation_turn t on t.tenant_id=b.tenant_id and t.channel_code=b.channel_code and t.provider_message_id=$5 and t.external_id=$6 and t.thread_id=$7 and t.occurred_at=to_timestamp($8) and t.assistant_text=$9 and t.state in ('completed','handed_off') and t.expires_at>statement_timestamp()
 join crm.consent_evidence consent on consent.tenant_id=b.tenant_id and consent.lead_id=b.lead_id and consent.consent_id=$10 and consent.purpose_code=$11 and consent.policy_version=b.policy_version and consent.evidence_sha256_hex=$12 and consent.decision='granted' and consent.occurred_at<=statement_timestamp()
 where b.tenant_id=$1 and b.channel_code='whatsapp' and b.external_id_hmac=$2 and b.version=$13 and b.lead_id=$14 and b.subject_id=$15 and b.policy_version=$16 and b.state='active' and b.pii_allowed and b.effective_at<=statement_timestamp()
 and t.occurred_at<=statement_timestamp() and t.occurred_at+interval '24 hours'>statement_timestamp()+interval '40 seconds'
 and not exists(select 1 from crm.consent_evidence newer where newer.tenant_id=consent.tenant_id and newer.lead_id=consent.lead_id and newer.purpose_code=consent.purpose_code and newer.consent_id<>consent.consent_id and newer.occurred_at>=consent.occurred_at and newer.occurred_at<=statement_timestamp())`, s.tenant, c.ExternalIDHMAC, s.organization, s.connection, c.ProviderMessageID, c.Message.ExternalID, c.Message.ThreadID, request.LastInboundAt, request.Text, c.ConsentID, c.ConsentPurpose, c.ConsentEvidenceSHA256, c.BindingVersion, c.LeadID, c.SubjectID, c.PolicyVersion).Scan(&found)
	if err != nil || !found {
		return ErrApproval
	}
	return nil
}

// Read never derives tenant/org from the request body. It returns the exact
// stored proposal for the human to review, together with its approval hash.
func (s *PostgresReplyApprovals) Read(ctx context.Context, p identity.Principal, key string) (ReplyProposal, error) {
	var value ReplyProposal
	var raw []byte
	if !s.validPrincipal(p, "whatsapp:approve") {
		return value, ErrApproval
	}
	err := s.base.pool.QueryRow(ctx, `select request_id,state,evidence_sha,payload from approval.request where tenant_id=$1 and request_id=$2 and organization_id=$3 and kind='whatsapp_reply'`, s.tenant, key, s.organization).Scan(&value.RequestID, &value.State, &value.PayloadSHA256, &raw)
	_, actualHash, hashErr := approval.CanonicalPayload(raw)
	if err != nil || hashErr != nil || actualHash != value.PayloadSHA256 || json.Unmarshal(raw, &value.Context) != nil {
		return ReplyProposal{}, ErrApproval
	}
	if value.State == "approved" {
		status, err := readNotificationStatus(ctx, s.base.pool, p, s.organization, "", key)
		if err != nil {
			return ReplyProposal{}, err
		}
		value.Status = &status
	}
	return value, nil
}

func (s *PostgresReplyApprovals) ResolveWhatsAppApproval(ctx context.Context, tenant, key string) (Approval, error) {
	var value Approval
	var c ReplyContext
	var raw []byte
	if s == nil || tenant != s.tenant {
		return value, ErrApproval
	}
	err := s.base.pool.QueryRow(ctx, `select r.payload,r.evidence_sha,d.reviewer,r.decided_at from approval.request r join approval.decision d on d.tenant_id=r.tenant_id and d.request_id=r.request_id and d.approved and d.reviewer<>r.requester where r.tenant_id=$1 and r.request_id=$2 and r.organization_id=$3 and r.kind='whatsapp_reply' and r.state='approved' and (select count(*) from approval.decision d2 where d2.tenant_id=r.tenant_id and d2.request_id=r.request_id)=1`, tenant, key, s.organization).Scan(&raw, &value.EvidenceSHA256, &value.ApprovedBy, &value.NotBefore)
	_, actualHash, hashErr := approval.CanonicalPayload(raw)
	if err != nil || hashErr != nil || actualHash != value.EvidenceSHA256 || json.Unmarshal(raw, &c) != nil || s.current(ctx, s.base.pool, c) != nil {
		return Approval{}, ErrApproval
	}
	value.MessageSHA256 = c.MessageSHA256
	value.ProfileSHA256 = c.ProfileSHA256
	value.ExpiresAt = c.ExpiresAt
	return value, nil
}

// Propose stores a pending request only; the worker principal cannot approve it.
func (s *PostgresReplyApprovals) Propose(ctx context.Context, p identity.Principal, event string, in channels.Message) (ReplyProposal, error) {
	c, err := s.buildProposal(ctx, p, event, in)
	if err != nil {
		// The runtime already persisted a governed handoff (no model/tool effect for
		// unresolved identity). It is deliberately not made into an outbound request.
		if s.validPrincipal(p, "whatsapp:process") && in.TenantID == s.tenant && in.ThreadID == conversationThread(s.connection) {
			var found bool
			if s.base.pool.QueryRow(ctx, `select true from communication.conversation_turn where tenant_id=$1 and channel_code='whatsapp' and provider_message_id=$2 and external_id=$3 and thread_id=$4 and occurred_at=$5 and user_text=$6 and state='handed_off'`, s.tenant, in.ProviderMessageID, in.ExternalID, in.ThreadID, in.OccurredAt, in.Text).Scan(&found) == nil && found {
				return ReplyProposal{State: "HANDOFF_REVIEW_REQUIRED"}, nil
			}
		}
		return ReplyProposal{}, err
	}
	raw, _ := json.Marshal(c)
	canonical, hash, err := approval.CanonicalPayload(raw)
	if err != nil {
		return ReplyProposal{}, ErrApproval
	}
	spec := postgres.HumanApprovalSpec{Request: approval.Request{TenantID: s.tenant, ID: c.Message.DeliveryKey, Kind: approval.KindWhatsAppReply, SubjectID: c.ExternalIDHMAC, Requester: p.Subject, EvidenceSHA: hash}, OrganizationID: s.organization, Payload: canonical}
	_, err = postgres.NewHumanApprovals(s.base.pool).Submit(ctx, p, spec, "whatsapp:process", func(ctx context.Context, tx pgx.Tx) error { return s.current(ctx, tx, c) })
	if err != nil {
		return ReplyProposal{}, err
	}
	var state string
	if s.base.pool.QueryRow(ctx, `select state from approval.request where tenant_id=$1 and request_id=$2`, s.tenant, c.Message.DeliveryKey).Scan(&state) != nil {
		return ReplyProposal{}, ErrApproval
	}
	return ReplyProposal{RequestID: c.Message.DeliveryKey, State: state, PayloadSHA256: hash, Context: c}, nil
}

func (s *PostgresReplyApprovals) Decide(ctx context.Context, p identity.Principal, key, hash string, approved bool, reason string) (ReplyProposal, error) {
	value, err := s.Read(ctx, p, key)
	if err != nil || value.PayloadSHA256 != hash {
		return ReplyProposal{}, ErrApproval
	}
	guard := func(ctx context.Context, tx pgx.Tx) error {
		if !approved {
			return nil
		}
		return s.current(ctx, tx, value.Context)
	}
	_, err = postgres.NewHumanApprovals(s.base.pool).Decide(ctx, p, s.tenant, key, s.organization, hash, approved, reason, "whatsapp:approve", guard)
	if err != nil {
		return ReplyProposal{}, err
	}
	return s.Read(ctx, p, key)
}
````

### FILE: `internal/whatsappbridge/reply_module.go`

```yaml
block_id: "COMMUNICATIONS-DELTA-PYTHON_META_WHATSAPP_CLOUD_ADAPTER:file10:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "6251e73e897e2e6eada1111ce5eee0cce0bc574e6ccd57b3b85a2d30818c4318"
variables: []
secrets_allowed: false
```

````go
package whatsappbridge

// AUTHORED protected review/send boundary. The HTTP caller can approve/reject
// an exact stored proposal but cannot supply message text or contact identity.
import (
	"bytes"
	"context"
	"elite.local/enterprise/internal/approval"
	"elite.local/enterprise/internal/channels"
	"elite.local/enterprise/internal/outbounddelivery"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/postgres"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"
)

type ReplyModule struct {
	approvals *PostgresReplyApprovals
	sender    *Sender
	store     *postgres.OutboundDeliveryStore
	channel   *outbounddelivery.Channel
}
type VerifiedInboxChannel struct{}

func (VerifiedInboxChannel) Code() string { return "whatsapp" }
func (VerifiedInboxChannel) Receive(context.Context) ([]channels.Message, error) {
	return nil, ErrBridge
}
func (VerifiedInboxChannel) Send(context.Context, channels.Message) error { return ErrBridge }
func NewReplyModule(approvals *PostgresReplyApprovals, sender *Sender, store *postgres.OutboundDeliveryStore) (*ReplyModule, error) {
	if approvals == nil || sender == nil || store == nil || sender.TenantID != approvals.tenant || digest(sender.Profile) != approvals.profile {
		return nil, ErrBridge
	}
	cp := *sender
	cp.Profile = append(json.RawMessage(nil), sender.Profile...)
	cp.Approvals = approvals
	return &ReplyModule{approvals: approvals, sender: &cp, store: store, channel: &outbounddelivery.Channel{CodeValue: "whatsapp", Receiver: VerifiedInboxChannel{}, Sender: &cp, Store: store}}, nil
}
func (m *ReplyModule) principal(w http.ResponseWriter, r *http.Request, v identity.Verifier, permission string) (identity.Principal, bool) {
	values := r.Header.Values("Authorization")
	if v == nil || len(values) != 1 || !strings.HasPrefix(values[0], "Bearer ") || len(values[0]) > 16391 || len(strings.Fields(values[0])) != 2 {
		notificationProblem(w, 401, "UNAUTHENTICATED")
		return identity.Principal{}, false
	}
	p, err := v.Verify(r.Context(), strings.TrimPrefix(values[0], "Bearer "))
	if err != nil {
		notificationProblem(w, 401, "UNAUTHENTICATED")
		return identity.Principal{}, false
	}
	if m == nil || !m.approvals.validPrincipal(p, permission) {
		notificationProblem(w, 403, "FORBIDDEN")
		return identity.Principal{}, false
	}
	return p, true
}
func replyJSON(w http.ResponseWriter, value any) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(value)
}
func readReplyCommand(w http.ResponseWriter, r *http.Request, decision bool) (string, bool, string, error) {
	if r.Header.Get("Content-Type") != "application/json" || r.URL.RawQuery != "" {
		return "", false, "", ErrApproval
	}
	raw, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 4096))
	if err != nil {
		return "", false, "", ErrApproval
	}
	canonical, _, err := approval.CanonicalPayload(raw)
	if err != nil {
		return "", false, "", ErrApproval
	}
	var fields map[string]json.RawMessage
	if json.Unmarshal(canonical, &fields) != nil {
		return "", false, "", ErrApproval
	}
	expected := 1
	if decision {
		expected = 3
	}
	if len(fields) != expected {
		return "", false, "", ErrApproval
	}
	for _, raw := range fields {
		if bytes.Equal(bytes.TrimSpace(raw), []byte("null")) {
			return "", false, "", ErrApproval
		}
	}
	var hash, reason string
	var approved bool
	if json.Unmarshal(fields["payload_sha256"], &hash) != nil || !validDigest(hash) {
		return "", false, "", ErrApproval
	}
	if decision {
		if json.Unmarshal(fields["approved"], &approved) != nil || json.Unmarshal(fields["reason"], &reason) != nil || len(reason) > 2048 {
			return "", false, "", ErrApproval
		}
	}
	return hash, approved, reason, nil
}
func (m *ReplyModule) Register(mux *http.ServeMux, v identity.Verifier) {
	mux.HandleFunc("GET /v1/franchise/whatsapp/replies", func(w http.ResponseWriter, r *http.Request) {
		p, ok := m.principal(w, r, v, "whatsapp:approve")
		if !ok {
			return
		}
		if r.URL.RawQuery != "" {
			notificationProblem(w, 400, "INVALID_QUERY")
			return
		}
		rows, err := m.approvals.base.pool.Query(r.Context(), `select request_id,state,evidence_sha,payload from approval.request where tenant_id=$1 and organization_id=$2 and kind='whatsapp_reply' order by created_at desc,request_id limit 50`, p.TenantID, m.approvals.organization)
		if err != nil {
			notificationProblem(w, 503, "UNAVAILABLE")
			return
		}
		defer rows.Close()
		values := []ReplyProposal{}
		for rows.Next() {
			var value ReplyProposal
			var raw []byte
			if rows.Scan(&value.RequestID, &value.State, &value.PayloadSHA256, &raw) != nil || json.Unmarshal(raw, &value.Context) != nil {
				notificationProblem(w, 503, "UNAVAILABLE")
				return
			}
			values = append(values, value)
		}
		if rows.Err() != nil {
			notificationProblem(w, 503, "UNAVAILABLE")
			return
		}
		rows.Close()
		for i := range values {
			value, err := m.approvals.Read(r.Context(), p, values[i].RequestID)
			if err != nil {
				notificationProblem(w, 503, "UNAVAILABLE")
				return
			}
			values[i] = value
		}
		replyJSON(w, values)
	})
	mux.HandleFunc("GET /v1/franchise/whatsapp/replies/{id}", func(w http.ResponseWriter, r *http.Request) {
		p, ok := m.principal(w, r, v, "whatsapp:approve")
		if !ok {
			return
		}
		value, err := m.approvals.Read(r.Context(), p, r.PathValue("id"))
		if err != nil {
			notificationProblem(w, 404, "NOT_FOUND")
			return
		}
		replyJSON(w, value)
	})
	mux.HandleFunc("POST /v1/franchise/whatsapp/replies/{id}/decision", func(w http.ResponseWriter, r *http.Request) {
		p, ok := m.principal(w, r, v, "whatsapp:approve")
		if !ok {
			return
		}
		hash, approved, reason, err := readReplyCommand(w, r, true)
		if err != nil {
			notificationProblem(w, 400, "INVALID_BODY")
			return
		}
		value, err := m.approvals.Decide(r.Context(), p, r.PathValue("id"), hash, approved, reason)
		if err != nil {
			notificationProblem(w, 409, "APPROVAL_NOT_CURRENT_OR_DIVERGENT")
			return
		}
		replyJSON(w, value)
	})
	mux.HandleFunc("POST /v1/franchise/whatsapp/replies/{id}/recover", func(w http.ResponseWriter, r *http.Request) {
		p, ok := m.principal(w, r, v, "whatsapp:send")
		if !ok {
			return
		}
		hash, _, _, err := readReplyCommand(w, r, false)
		if err != nil {
			notificationProblem(w, 400, "INVALID_BODY")
			return
		}
		_, err = m.Recover(r.Context(), p, r.PathValue("id"), hash)
		if err != nil {
			notificationProblem(w, 409, "DELIVERY_RECONCILIATION_REQUIRED")
			return
		}
		replyJSON(w, map[string]string{"status": "accepted", "delivery_key": r.PathValue("id")})
	})

	mux.HandleFunc("POST /v1/franchise/whatsapp/replies/{id}/send", func(w http.ResponseWriter, r *http.Request) {
		p, ok := m.principal(w, r, v, "whatsapp:send")
		if !ok {
			return
		}
		// Read requires review permission too: a sender must be able to inspect the
		// exact proposal being released to the provider.
		hash, _, _, err := readReplyCommand(w, r, false)
		if err != nil {
			notificationProblem(w, 400, "INVALID_BODY")
			return
		}
		value, err := m.approvals.Read(r.Context(), p, r.PathValue("id"))
		if err != nil || value.State != "approved" || value.PayloadSHA256 != hash {
			notificationProblem(w, 409, "APPROVAL_NOT_CURRENT_OR_DIVERGENT")
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 45*time.Second)
		defer cancel()
		if _, err = m.approvals.ResolveWhatsAppApproval(ctx, p.TenantID, value.RequestID); err != nil {
			notificationProblem(w, 409, "APPROVAL_NOT_CURRENT_OR_DIVERGENT")
			return
		}
		if err = m.channel.Send(ctx, value.Context.Message); err != nil {
			notificationProblem(w, 409, "DELIVERY_RECONCILIATION_REQUIRED")
			return
		}
		replyJSON(w, map[string]string{"status": "accepted", "delivery_key": value.RequestID})
	})
}
````

### FILE: `internal/whatsappbridge/reply_recovery.go`

```yaml
block_id: "COMMUNICATIONS-DELTA-PYTHON_META_WHATSAPP_CLOUD_ADAPTER:file11:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "e93e50776e7ea8f9704aa125e63f9046f157b175137fdf1722802a6cd87b55b9"
variables: []
secrets_allowed: false
```

````go
package whatsappbridge

// AUTHORED recovery of an existing local provider acceptance receipt. It never
// performs POST, invents a provider messageID or releases an unproved unknown.
import (
	"bytes"
	"context"
	"elite.local/enterprise/internal/outbounddelivery"
	"elite.local/enterprise/internal/platform/identity"
	"encoding/json"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func recoverEvidence(root *os.Root, path string) ([]byte, error) {
	f, err := root.Open(path)
	if err != nil {
		return nil, ErrBridge
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil || !info.Mode().IsRegular() || info.Size() < 1 || info.Size() > 65536 {
		return nil, ErrBridge
	}
	raw, err := io.ReadAll(io.LimitReader(f, 65537))
	if err != nil || len(raw) > 65536 {
		return nil, ErrBridge
	}
	return raw, nil
}
func (m *ReplyModule) Recover(ctx context.Context, p identity.Principal, key, hash string) (outbounddelivery.Receipt, error) {
	if m == nil || !m.approvals.validPrincipal(p, "whatsapp:send") {
		return outbounddelivery.Receipt{}, ErrApproval
	}
	value, err := m.approvals.Read(ctx, p, key)
	if err != nil || value.State != "approved" || value.PayloadSHA256 != hash {
		return outbounddelivery.Receipt{}, ErrApproval
	}
	if value.Status == nil || (value.Status.FenceState != "unknown" && value.Status.FenceState != "sending" && value.Status.FenceState != "accepted") {
		return outbounddelivery.Receipt{}, ErrBridge
	}
	c := value.Context
	var approvalAt time.Time
	var reviewer string
	err = m.approvals.base.pool.QueryRow(ctx, `select r.decided_at,d.reviewer from approval.request r join approval.decision d on d.tenant_id=r.tenant_id and d.request_id=r.request_id and d.approved and d.reviewer<>r.requester where r.tenant_id=$1 and r.request_id=$2 and r.organization_id=$3 and r.kind='whatsapp_reply' and r.state='approved' and (select count(*) from approval.decision d2 where d2.tenant_id=r.tenant_id and d2.request_id=r.request_id)=1`, p.TenantID, key, m.approvals.organization).Scan(&approvalAt, &reviewer)
	if err != nil || reviewer == "" || c.ConnectionID != m.approvals.connection || c.ProfileSHA256 != m.approvals.profile {
		return outbounddelivery.Receipt{}, ErrApproval
	}
	request := json.RawMessage(c.Message.Text)
	actual, err := outbounddelivery.MessageSHA256(c.Message)
	if err != nil || actual != c.MessageSHA256 {
		return outbounddelivery.Receipt{}, ErrApproval
	}
	proc := m.sender.Process
	script := filepath.Join(proc.AdapterDirectory, "whatsapp_cloud.py")
	if !filepath.IsAbs(proc.EvidenceDirectory) || exactFile(proc.PythonExecutable, proc.PythonSHA256) != nil || exactFile(script, proc.AdapterSHA256) != nil {
		return outbounddelivery.Receipt{}, ErrBridge
	}
	root, err := os.OpenRoot(proc.EvidenceDirectory)
	if err != nil {
		return outbounddelivery.Receipt{}, ErrBridge
	}
	defer root.Close()
	dir := digest([]byte(p.TenantID + "\x00" + key))
	receiptRaw, err := recoverEvidence(root, filepath.Join(dir, "SEND_RECEIPT.json"))
	if err != nil {
		return outbounddelivery.Receipt{}, ErrBridge
	}
	responseRaw, err := recoverEvidence(root, filepath.Join(dir, "provider-response.json"))
	if err != nil {
		return outbounddelivery.Receipt{}, ErrBridge
	}
	raw, err := json.Marshal(struct {
		Schema   string          `json:"schema"`
		Binding  string          `json:"binding_sha256"`
		Profile  json.RawMessage `json:"profile"`
		Request  json.RawMessage `json:"request"`
		Receipt  []byte          `json:"receipt"`
		Response []byte          `json:"response"`
	}{"elite-whatsapp-recover-send/v1", actual, m.sender.Profile, request, receiptRaw, responseRaw})
	if err != nil {
		return outbounddelivery.Receipt{}, ErrBridge
	}
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, proc.PythonExecutable, "-I", "-B", script, "--recover-send-bridge")
	cmd.Stdin = bytes.NewReader(raw)
	cmd.Stderr = io.Discard
	cmd.Env = []string{}
	cmd.WaitDelay = 2 * time.Second
	if root := os.Getenv("SystemRoot"); root != "" {
		cmd.Env = append(cmd.Env, "SystemRoot="+root)
	}
	var stdout boundedOutput
	cmd.Stdout = &stdout
	if cmd.Run() != nil || stdout.overflow {
		return outbounddelivery.Receipt{}, ErrBridge
	}
	var result result
	d := json.NewDecoder(bytes.NewReader(stdout.Bytes()))
	d.DisallowUnknownFields()
	if d.Decode(&result) != nil || d.Decode(new(any)) != io.EOF || result.Schema != "elite-whatsapp-send-result/v1" || result.BindingSHA256 != actual || result.EvidenceSHA256 != digest(receiptRaw) {
		return outbounddelivery.Receipt{}, ErrBridge
	}
	receipt := outbounddelivery.Receipt{ProviderMessageID: result.ProviderMessageID, EvidenceSHA256: result.EvidenceSHA256, AcceptedAt: result.AcceptedAt}
	if receipt.Validate() != nil || receipt.AcceptedAt.Before(approvalAt) || !receipt.AcceptedAt.Before(c.ExpiresAt) {
		return outbounddelivery.Receipt{}, ErrBridge
	}
	// Use the same owner's expired-lease transition; it never calls the provider.
	// The existing row was read above, and this path cannot create a new delivery.
	if value.Status.FenceState == "sending" {
		claim, claimErr := m.store.Claim(ctx, c.Message, actual)
		if claim.Replay {
			value.Status.FenceState = "accepted"
		} else if !errors.Is(claimErr, outbounddelivery.ErrUnknown) {
			return outbounddelivery.Receipt{}, ErrBridge
		}
	}
	if value.Status.FenceState == "accepted" {
		var same bool
		err = m.approvals.base.pool.QueryRow(ctx, `select evidence_sha256_hex=$3 and request_sha256_hex=$4 from communication.outbound_delivery where tenant_id=$1 and channel_code='whatsapp' and delivery_key=$2 and state='accepted'`, p.TenantID, key, receipt.EvidenceSHA256, actual).Scan(&same)
		if err != nil || !same {
			return outbounddelivery.Receipt{}, ErrBridge
		}
		return receipt, nil
	}
	if err = m.store.ReconcileAccepted(ctx, c.Message, actual, receipt); err != nil {
		return outbounddelivery.Receipt{}, err
	}
	return receipt, nil
}
````

### FILE: `internal/whatsappbridge/reply_request.go`

```yaml
block_id: "COMMUNICATIONS-DELTA-PYTHON_META_WHATSAPP_CLOUD_ADAPTER:file12:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "fc054e3a26040fe128d944d4e16062984b252af394700d3b186391e2f263e274"
variables: []
secrets_allowed: false
```

````go
package whatsappbridge

// AUTHORED projection guard for one exact human-approved Meta text request.
import (
	"bytes"
	"elite.local/enterprise/internal/channels"
	"encoding/json"
	"io"
	"time"
	"unicode/utf8"
)

type ReplyRequest struct {
	Kind            string `json:"kind"`
	Recipient       string `json:"recipient"`
	Text            string `json:"text"`
	SourceMessageID string `json:"source_message_id"`
	LastInboundAt   int64  `json:"last_inbound_at"`
	WindowExpiresAt int64  `json:"window_expires_at"`
}

func replyRequestFor(message channels.Message) (json.RawMessage, error) {
	var body ReplyRequest
	d := json.NewDecoder(bytes.NewBufferString(message.Text))
	d.DisallowUnknownFields()
	if d.Decode(&body) != nil || d.Decode(new(any)) != io.EOF || body.Kind != "text_reply" || body.Recipient != message.ExternalID || !utf8.ValidString(body.Text) || len(body.Text) == 0 || len(body.Text) > 4096 || utf8.RuneCountInString(body.Text) > 1024 || body.SourceMessageID == "" || len(body.SourceMessageID) > 256 || body.LastInboundAt < 1 || body.WindowExpiresAt != body.LastInboundAt+86400 {
		return nil, ErrBridge
	}
	for _, r := range body.Text {
		if r < 32 && r != '\n' && r != '\t' {
			return nil, ErrBridge
		}
	}
	now := time.Now().Unix()
	if now < body.LastInboundAt || now >= body.WindowExpiresAt {
		return nil, ErrBridge
	}
	raw, err := json.Marshal(body)
	if err != nil {
		return nil, ErrBridge
	}
	return raw, nil
}
````

### FILE: `microsoft_playwright_browser_gate/tests/whatsapp-connected.spec.mjs`

```yaml
block_id: "COMMUNICATIONS-DELTA-PYTHON_META_WHATSAPP_CLOUD_ADAPTER:file13:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "4df0f7aabdccd2efd0a0abd79867ebf85de12457cf7eb845a529cd33ec5964af"
variables: []
secrets_allowed: false
```

````text
import {test,expect} from '@playwright/test';
import {createHash} from 'node:crypto';
import {createRequire} from 'node:module';
import {resolve} from 'node:path';
import {pathToFileURL} from 'node:url';

test('human reviews WhatsApp proposals and reconciles an uncertain send through real BFF Go and PostgreSQL',async({page,context},info)=>{
 if(process.env.ELITE_WHATSAPP_BROWSER!=='1')throw new Error('explicit local WhatsApp fixture required');
 const base=process.env.ELITE_BASE_URL;expect(new URL(base).protocol).toBe('https:');expect(new URL(base).hostname).toBe('127.0.0.1');
 const identities=JSON.parse(process.env.ELITE_WHATSAPP_IDENTITIES);
 const require=createRequire(resolve(process.env.ELITE_WEB_ROOT,'package.json'));
 const {EncryptJWT}=await import(pathToFileURL(require.resolve('jose')).href);
 async function identity(name){
  const jwt=await new EncryptJWT(identities[name]).setProtectedHeader({alg:'dir',enc:'A256GCM',typ:'JWT'}).setIssuedAt().setExpirationTime('300s').encrypt(createHash('sha256').update(process.env.AUTH_SESSION_SECRET).digest());
  await context.clearCookies();await context.addCookies([{name:'__Host-elite_session',value:jwt,url:base+'/',secure:true,httpOnly:true,sameSite:'Lax'}]);
 }
 const errors=[];page.on('pageerror',error=>errors.push(error.message));let sends=0,recoveries=0;
 page.on('request',r=>{if(r.method()==='POST'&&r.url().includes('/api/enterprise/franchise/whatsapp/replies')){const b=r.postDataJSON();if(b.action==='send')sends++;if(b.action==='recover')recoveries++}});
 await identity('reader');await page.goto('/franchise');await expect(page.getByRole('link',{name:'Revisar respuestas de WhatsApp'})).toHaveCount(0);
 await page.goto('/franchise/whatsapp');await expect(page.getByText('Tu sesión no permite revisar respuestas.')).toBeVisible();
 await identity('foreign-org');await page.goto('/franchise/whatsapp');await expect(page.getByRole('alert').filter({hasText:'No pudimos consultar'})).toBeVisible();await expect(page.locator('article')).toHaveCount(0);
 await identity('human');await page.goto('/franchise');await page.getByRole('link',{name:'Revisar respuestas de WhatsApp'}).click();
 await expect(page.getByRole('heading',{name:'Respuestas de WhatsApp',exact:true})).toBeVisible();await expect(page.locator('article')).toHaveCount(3);
 const rejected=page.locator('article').nth(0),lost=page.locator('article').nth(1),normal=page.locator('article').nth(2);
 await expect(normal).toContainText('Cotización preparada para revisar.');await expect(normal).toContainText('5491112345678');
 await rejected.getByRole('button',{name:'Rechazar',exact:true}).click();await expect(rejected).toContainText('Revisión: Rechazada');await expect(rejected.getByRole('button')).toHaveCount(0);
 await normal.getByRole('button',{name:'Aprobar este texto'}).click();await expect(normal).toContainText('Revisión: Aprobada');expect(sends).toBe(0);
 await identity('review-only');await page.reload();await expect(normal.getByRole('button',{name:'Enviar respuesta aprobada'})).toHaveCount(0);
 await identity('human');await page.reload();await normal.getByRole('button',{name:'Enviar respuesta aprobada'}).click();await expect(normal).toContainText('Envío: Aceptado por WhatsApp');expect(sends).toBe(1);
 expect((await context.request.post(process.env.ELITE_WHATSAPP_CONTROL,{headers:{'X-Fixture-Token':process.env.ELITE_WHATSAPP_CONTROL_TOKEN},data:{index:1}})).status()).toBe(200);
 await page.getByRole('link',{name:'Consultar estado',exact:true}).click();await expect(normal).toContainText('Entregado según WhatsApp');
 await lost.getByRole('button',{name:'Aprobar este texto'}).click();await expect(lost).toContainText('Revisión: Aprobada');
 await lost.getByRole('button',{name:'Enviar respuesta aprobada'}).click();await expect(page.getByRole('status')).toContainText('La acción no está confirmada');expect(sends).toBe(2);
 await expect(lost.getByRole('button',{name:'Enviar respuesta aprobada'})).toBeDisabled();
 await page.getByRole('link',{name:'Consultar estado',exact:true}).click();await expect(lost).toContainText('El resultado del envío necesita reconciliación.');await expect(lost.getByRole('button',{name:'Enviar respuesta aprobada'})).toHaveCount(0);
 await page.screenshot({path:info.outputPath('whatsapp-reconciliation-desktop.png'),fullPage:true});
 await lost.getByRole('button',{name:'Verificar comprobante guardado'}).click();await expect(lost).toContainText('Envío: Aceptado por WhatsApp');expect(recoveries).toBe(1);expect(sends).toBe(2);
 expect((await context.request.post(process.env.ELITE_WHATSAPP_CONTROL,{headers:{'X-Fixture-Token':process.env.ELITE_WHATSAPP_CONTROL_TOKEN},data:{index:2}})).status()).toBe(200);
 await page.getByRole('link',{name:'Consultar estado',exact:true}).click();await expect(lost).toContainText('Entregado según WhatsApp');await expect(normal).toContainText('Entregado según WhatsApp');
 await page.setViewportSize({width:390,height:844});expect(await page.evaluate(()=>document.documentElement.scrollWidth<=window.innerWidth)).toBe(true);await page.screenshot({path:info.outputPath('whatsapp-recovered-mobile.png'),fullPage:true});expect(errors).toEqual([]);
});
````

### FILE: `whatsapp_cloud/test_conversation_contract.py`

```yaml
block_id: "COMMUNICATIONS-DELTA-PYTHON_META_WHATSAPP_CLOUD_ADAPTER:file14:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "9589539bf7c1c9c3c7044915787c4bc9f44564997708ef2e9d6db735ed261afa"
variables: []
secrets_allowed: false
```

````python
import base64,hashlib,hmac,json,tempfile,time,unittest
from pathlib import Path
from test_whatsapp_cloud import profile
from whatsapp_cloud import build_reply_payload,normalize_verified_webhook,send_bridge,recover_send_bridge

class ConversationContractTests(unittest.TestCase):
    def request(self):
        start=int(time.time())-10
        return {"kind":"text_reply","recipient":"5491112345678","text":"Texto exacto <revisado>","source_message_id":"wamid.fixture","last_inbound_at":start,"window_expires_at":start+86400}
    def test_exact_text_payload_window_and_utf8(self):
        request=self.request()
        recipient,payload=build_reply_payload(profile(),request)
        self.assertEqual(recipient,request["recipient"])
        self.assertEqual(payload,{"messaging_product":"whatsapp","recipient_type":"individual","to":recipient,"type":"text","text":{"body":request["text"]}})
        for change in [{"text":"\ud800"},{"window_expires_at":request["window_expires_at"]+1},{"last_inbound_at":1,"window_expires_at":86401},{"template_name":"invented"}]:
            with self.subTest(change=list(change)),self.assertRaises((ValueError,PermissionError,UnicodeError)):
                build_reply_payload(profile(),{**request,**change})
    def test_recovery_requires_bound_original_evidence_without_transport(self):
        request=self.request()
        with tempfile.TemporaryDirectory() as tmp:
            folder=Path(tmp)/"send"
            frame={"schema":"elite-whatsapp-send-bridge/v1","binding_sha256":"a"*64,"profile":profile(),"request":request,"access_token":"synthetic","output_directory":str(folder)}
            sent=send_bridge(json.dumps(frame).encode(),lambda *args:(200,{},b'{"messages":[{"id":"wamid.original"}]}'))
            original={"schema":"elite-whatsapp-recover-send/v1","binding_sha256":"a"*64,"profile":profile(),"request":request,"receipt":base64.b64encode((folder/"SEND_RECEIPT.json").read_bytes()).decode(),"response":base64.b64encode((folder/"provider-response.json").read_bytes()).decode()}
            self.assertEqual(recover_send_bridge(json.dumps(original).encode()),sent)
            changed={**original,"request":{**request,"text":"otro texto"}}
            with self.assertRaises(ValueError):recover_send_bridge(json.dumps(changed).encode())
            changed={**original,"response":base64.b64encode(b'{"messages":[{"id":"wamid.forged"}]}').decode()}
            with self.assertRaises(ValueError):recover_send_bridge(json.dumps(changed).encode())

if __name__=="__main__":unittest.main()
````

### FILE: `whatsapp_cloud/test_profile_validation_bridge.py`

```yaml
block_id: "COMMUNICATIONS-DELTA-PYTHON_META_WHATSAPP_CLOUD_ADAPTER:file15:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "4d0882c99183078e0703ff0a4bb7da5bdda5046724891ba2eb69c123c99c1a71"
variables: []
secrets_allowed: false
```

````python
import base64, hashlib, json, subprocess, sys, unittest
from pathlib import Path
from test_whatsapp_cloud import profile

class ProfileValidationBridgeTests(unittest.TestCase):
    def test_blocked_default_and_exact_fixture(self):
        root=Path(__file__).resolve().parent
        for good, raw in [(False,(root/"provider-profile.template.json").read_bytes()),(True,(json.dumps(profile(),ensure_ascii=False,indent=2)+"\n").encode("utf-8"))]:
            with self.subTest(accepted=good):
                frame=json.dumps({"schema":"elite-whatsapp-validate-profile/v1","profile":base64.b64encode(raw).decode()}).encode()
                result=subprocess.run([sys.executable,"-I","-B",str(root/"whatsapp_cloud.py"),"--validate-profile-bridge"],input=frame,capture_output=True)
                if good:
                    self.assertEqual(result.returncode,0)
                    self.assertEqual(json.loads(result.stdout),{"schema":"elite-whatsapp-profile-validation/v1","profile_sha256":hashlib.sha256(raw).hexdigest()})
                else:
                    self.assertNotEqual(result.returncode,0)
                    self.assertEqual(result.stdout,b"")

if __name__=="__main__":unittest.main()
````

## 6. Configuration surface

El profile exige licencia/terms, Business/app/phone/template/webhook/consent, cuota/costo y reconciliación `PROVEN`; versión Graph y templates exactos; retención y tres referencias de secrets. La request separa recipient/template/language/parameters. No hay valores de secret ni aprobación embebida.

## 7. Dependency bill

| Dependencia | Pin | Licencia | Uso |
|---|---|---|---|
| `fbsamples/whatsapp-api-examples` | commit firmado `de70ee90…9466`, archive SHA `38d183a0…a7c` | Meta Platform API-only | patrones y cuatro referencias exactas |
| Python stdlib | 3.12+ | PSF | HTTP, HMAC, JSON, atomic evidence |

El pack no depende del `WhatsApp-Nodejs-SDK` archivado. La licencia Meta y Platform Policy aplican también al archivo adaptado; todo redistribuidor conserva LICENSE, source lock y PROVENANCE.

## 8. Apply order

1. Materializar y verificar 22 hashes, incluidos cuatro archivos de referencia Meta; Go y SQL requieren el perfil integral compuesto con los owners de turnos, consentimiento, contacto, channels y outbounddelivery. El perfil Python separado no es un módulo Go standalone: no aplicar 0050 sin las migraciones anteriores.
2. Cerrar licencia/terms/acceso/costo/consent/template/webhook/reconciliación en readiness.
3. Elegir y registrar la versión Graph vigente desde autoridad oficial.
4. Aplicar 0001–0050 en base PostgreSQL local desechable; ejecutar los 22 tests Python offline y, en el perfil Go compuesto, los contratos approval/bridge/fence con Python exacto y ELITE_WHATSAPP_TEST_DATABASE_URL configurada. Construir el resolver con propósito/política server-side y Principal verificado, según README.
5. Para la operación explícita, construir NewAppointmentNotificationModule con los owners probados y registrarlo en el mux existente con identity.Verifier. No crea listener ni job. Automatización outbox e ingress al provider inbox durable permanecen como conexiones separadas que deben probarse en el proyecto.
6. Probar test number/template, firma real, duplicados/reintentos/status/opt-out y conciliación antes de habilitar tráfico.

## 9. Verification

```powershell
Push-Location .\whatsapp_cloud
python -m unittest -v test_whatsapp_cloud.py
Pop-Location
```

Resultado local V264: hashes oficiales exactos y 17 tests PASS. El PASS histórico V1 de ocho tests permanece en su evidencia, pero no demostraba aislamiento WABA/teléfono. No se ejecuta el ejemplo oficial defectuoso. Cuenta, términos actuales, versión Graph, template/test phone, delivery, costo, opt-out, retries y reconciliación siguen `PROJECT_CONDITIONED`.

V265: 22 tests Python, cuatro tests Go superiores/19 subcasos, incluyendo cinco
con PostgreSQL real y child sintético; full Go test/vet/build y reconstrucción
6/6 PASS. El bridge ejecutable existe; el resolver durable de aprobaciones,
turno/contacto y correlación de statuses todavía deben conectarse.

V266 cierra la implementación local del resolver durable: 19 negativos, seis
invalidaciones y ocho aprobaciones concurrentes; fence/PG real con child sintético,
full Go/vet/build, 22 Python, migraciones/down-up y reconstrucción 7/7 PASS.
El registro V265 anterior es histórico. API/worker autenticados, captura real
de consentimiento, cuenta y status/reconciliación siguen pendientes.

V267 añade API explícita protegida mediante el verifier existente. Once negativos
del parser, trece rechazos HTTP antes del fence y tres recorridos HTTP/RS256/PG
(concurrencia, pérdida de respuesta del proveedor y pérdida HTTP post-commit)
pasan; una sola invocación por evento. Full Go/vet/build, 22 Python y 4/4 hashes
reconstruidos PASS. No cierra UI, worker autónomo ni proveedor real.

## 10. Reconstruction evidence

- `reconstruction_evidence/WHATSAPP_READ_ONLY_OUTCOME_V268.md`: GET autenticado y snapshot local, 11 casos de estado/integridad, 13 de scope/errores y POST→GET real; sesión PostgreSQL read-only, full Go/vet/build, 22 Python y 5/5 hashes reconstruidos PASS. No habilita envío por consultar ni deduce entrega de accepted.
- `reconstruction_evidence/WHATSAPP_AUTHENTICATED_DISPATCH_V267.md`
- `reconstruction_evidence/WHATSAPP_APPOINTMENT_APPROVAL_V266.md`
- `reconstruction_evidence/WHATSAPP_GO_DURABLE_BRIDGE_V265.md`
- `reconstruction_evidence/WHATSAPP_SCOPED_WEBHOOK_V264.md`
- `reconstruction_evidence/META_WHATSAPP_CLOUD_ADAPTER_2026-08-26_V1.md`
- El SDK Node archivado permanece rechazado; se usa autoridad sample activa/firmada con licencia restringida.
- El defecto upstream de control flow queda retenido y no se relabela como PASS.

V402 composed delta: Signed raw inbox messages are projected into the existing conversation runtime with scoped current contact identity. Runtime output becomes an immutable proposal, human review gates the existing outbound fence, and status/recovery uses the same durable owners. Existing Python source retains its declared adaptation; new Go persistence/transport/tests are AUTHORED glue. Requires selected conversation, approval, contact, app and fence owners for Go composition.
