# V350 — real refund host/processor lifecycle qualification

Library maintenance / T2809; resume142 validated. AUTHORED qualification only,
CONDITIONED / NOT_ADMITTED as product. The canonical refund0.1.4 host loop and
Processor.Step are compiled unchanged. The fixture injects synthetic Store and
Provider implementations; it does not execute the production run() configuration,
PostgreSQL transactions, Stripe/Mercado Pago SDK traffic or financial effects.

## Concrete connection

A Windows-only Go adapter consumes the V348 trusted inherited manual-reset stop
event and cancels the context passed into the actual runRefundLoop/Processor.Step.
It owns a noninheritable SYNCHRONIZE duplicate; cleanup first joins its bounded
wait goroutine, then closes the handle. Parent cancellation causes survive. An
already-signaled event is inspected synchronously before exposing the context,
so the actual host does not make its first claim. Subsequent waits use20ms polls;
scheduling/suspension means this is not a guaranteed20ms shutdown SLO. Only the
trusted launcher/event protocol is assumed: no handle-type authentication, hostile
worker sandbox, service identity, or signal-based proof of domain completion.

The unchanged V349 store reserves the attempt before launch, governs and captures
the fixture process, then publishes its metadata-only process receipt. Fresh
reads reconcile it; repeated same-ID submissions cannot execute again. The
fixture separately writes a run-ID-bound lifecycle acknowledgement and a synthetic
domain result. The process receipt never treats either as business success.

The real processor classifies synthetic provider outcomes as succeeded/retry/
failed and blocks a missing provider. All can coexist with a clean host exit.
Cancellation after one synthetic provider effect can leave no domain result,
because Complete/Finish must not assume a cancelled context committed. Lost
domain acknowledgement is reconciled by reading the synthetic outcome, never by
replaying. A failed report writer preserves host exit1. A provider ignoring the
stop is forcibly terminated after finite grace, leaving no lifecycle/domain ack.

Test records use exclusive checked writes and sync for local observations only.
They are NOT PostgreSQL commit evidence, provider receipts or a new business
result store. A fixture claim/effect without domain result remains unresolved.
No automatic retry/new run ID, restart service, real account, spending or ARCA.

## Findings preserved

FAIL601 / LIB2317: recovery path and patch-context errors corrected from actual
owner/inventory reads; tool output retained in task history.
FAIL602 / LIB2318: initial integrated test invocation rejected --fixture before
ready; original source, log, ten failed cases and receipts retained. Explicit --
terminates Go flag parsing; ten corrected integrated cases passed. The fixture
now dispatches in TestMain so baseline/fuzz suites do not count it as PASS/skip.
FAIL603 / LIB2319: first native adapter returned a live context for a pre-signaled
event. presignaled-red.log preserves the regression; synchronous initial wait
fixes it. Fresh canonical reproduction below is required before closure.

## Authorities and reproducibility

SECURITY_SRE_CLOUD_INFRASTRUCTURE.md, DATABASE_STORAGE_INTERNALS.md,
ENGINEERING_EXECUTION_PLAYBOOK.md and PUBLIC_CODE_ARCHITECTURE_ADMISSION_STANDARD.md
govern lifecycle, acknowledgement, cancellation, provenance and narrow admission.
Official semantics consulted2026-09-09, without acquiring external code:
[WaitForSingleObject](https://learn.microsoft.com/en-us/windows/win32/api/synchapi/nf-synchapi-waitforsingleobject),
[DuplicateHandle](https://learn.microsoft.com/en-us/windows/win32/api/handleapi/nf-handleapi-duplicatehandle),
[Go context](https://pkg.go.dev/context#WithCancelCause).
Closing during a pending wait is undefined; the adapter joins first. Duplicate
ownership is distinct from authentication of an event or hostile-worker security.

Execution stays Go1.26.7 / CPython3.14.4 / PowerShell7.6.5, same trusted Windows x64.
Dependencies/go.mod/go.sum and all15refund product files remain exact; no new
runtime, dependency admission or product/profile promotion. No race-detector or
production durability/security claim. Finite fuzz targets canonical handle text,
not arbitrary live OS handles or business data; SAST/DAST/target gates remain open.

Stage: <LOCALAPPDATA>/Temp/elite-v350-a3669c125ee548bfb3e8d8717ef46234.
Materialize refund0.1.4 anew, add the three qualification Go files below to its
cmd/return-refund-worker directory and build a test binary. Reuse seven unchanged
V345/V347/V348/V349 native qualification files and the strengthened capture test below, add the Python suite below, and
supply an explicit test-binary path/hash fixture. Do not deploy that test binary.

Next: actual PostgreSQL-backed domain reconciliation and production startup/config
mounting, including trusted secret/runtime/identity policy, before admitting any
native service pack. Controlled synthetic lifecycle evidence closes this local
integration gap only; broad T2809 and the other roadmap fronts remain open.

## Qualification source: native_stop_windows.go (AUTHORED / NOT_ADMITTED)

SHA-256: `6bad310eed09d8a884381ab3d9ace446b7731a72df672374102ae8e2b2f00b76`

````go
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
````

## Qualification source: native_stop_windows_test.go (AUTHORED / NOT_ADMITTED)

SHA-256: `bc247f36eea3bfb028d0b85c6ef06da8451b4b055c0e51e1fbe6ad43ec3fcadf`

````go
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
````

## Qualification source: lifecycle_fixture_windows_test.go (AUTHORED / NOT_ADMITTED)

SHA-256: `25bcdd7d343a82f7b724fc34658b0f752938c638b1b04d50c837b45895d70741`

````go
//go:build windows

package main

// Real canonical host loop and Processor.Step; only Store/Provider below are
// synthetic. Test records are NOT PostgreSQL transactions or provider receipts.
import (
	"context"
	"elite.local/return-refund-worker/internal/refundworker"
	"encoding/json"
	"errors"
	"log"
	"os"
	"testing"
	"time"
)

type lifecycleStore struct {
	mode, run string
	claims    int
}

func fixturePublish(name string, v any) error {
	raw, err := json.Marshal(v)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return err
	}
	n, e := f.Write(raw)
	if e == nil && n != len(raw) {
		e = errors.New("fixture short write")
	}
	if e == nil {
		e = f.Sync()
	}
	c := f.Close()
	return errors.Join(e, c)
}
func fixtureReady() {
	if err := fixturePublish("ready", map[string]string{"state": "in-flight"}); err != nil {
		panic(err)
	}
}
func (s *lifecycleStore) Claim(ctx context.Context, worker, token string, lease time.Duration) (*refundworker.Work, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	s.claims++
	if s.claims != 1 {
		_ = fixturePublish("second-claim", s.claims)
		return nil, errors.New("unexpected second claim")
	}
	if err := fixturePublish("domain-claim.json", map[string]any{"run_id": s.run, "claims": s.claims}); err != nil {
		return nil, err
	}
	return &refundworker.Work{TenantID: "synthetic-tenant", RequestID: s.run, Attempt: 1, ClaimToken: token}, nil
}
func (s *lifecycleStore) Prepare(ctx context.Context, w refundworker.Work, worker string) (refundworker.Refund, error) {
	return refundworker.Refund{TenantID: w.TenantID, RequestID: w.RequestID, PaymentAttemptID: "synthetic-payment", Provider: "stripe", ProviderPaymentReference: "synthetic-pi", IdempotencyKey: s.run, Currency: "ARS", AmountMinorUnits: 1000}, ctx.Err()
}
func (s *lifecycleStore) record(ctx context.Context, outcome, code string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if s.mode == "store-failure" || s.mode == "report-failure" {
		fixtureReady()
		<-ctx.Done()
		return errors.New("PRIVATE_FIXTURE_STORE_ERROR")
	}
	err := fixturePublish("domain-result.json", map[string]any{"schema": "synthetic-refund-outcome/v1", "run_id": s.run, "outcome": outcome, "code": code, "claims": s.claims})
	if err != nil {
		return err
	}
	fixtureReady()
	<-ctx.Done()
	if s.mode == "lost-domain-ack" {
		return errors.New("PRIVATE_FIXTURE_ACK_LOSS")
	}
	return nil
}
func (s *lifecycleStore) Complete(ctx context.Context, w refundworker.Work, worker string, r refundworker.Refund, result refundworker.ProviderResult, outcome, code string, retry time.Duration) error {
	return s.record(ctx, outcome, code)
}
func (s *lifecycleStore) Finish(ctx context.Context, w refundworker.Work, worker, outcome, code string, retry time.Duration) error {
	return s.record(ctx, outcome, code)
}

type lifecycleProvider struct{ mode, run string }

func (p lifecycleProvider) Create(ctx context.Context, r refundworker.Refund) (refundworker.ProviderResult, error) {
	if p.mode == "effect-then-cancel" || p.mode == "ignore-stop" {
		if err := fixturePublish("synthetic-effect.json", map[string]string{"run_id": p.run}); err != nil {
			return refundworker.ProviderResult{}, err
		}
		fixtureReady()
		if p.mode == "ignore-stop" {
			time.Sleep(30 * time.Second)
		} else {
			<-ctx.Done()
		}
		return refundworker.ProviderResult{}, ctx.Err()
	}
	status := "succeeded"
	if p.mode == "pending" {
		status = "pending"
	}
	if p.mode == "failed" {
		status = "failed"
	}
	return refundworker.ProviderResult{ProviderPaymentReference: r.ProviderPaymentReference, ProviderRefundReference: "synthetic-re", ProviderStatus: status, Currency: r.Currency, AmountMinorUnits: r.AmountMinorUnits}, nil
}
func (p lifecycleProvider) Retrieve(ctx context.Context, r refundworker.Refund) (refundworker.ProviderResult, error) {
	return p.Create(ctx, r)
}

func TestMain(m *testing.M) {
	if len(os.Args) >= 3 && os.Args[len(os.Args)-3] == "--fixture" {
		runLifecycleFixture()
		os.Exit(0)
	}
	os.Exit(m.Run())
}

func runLifecycleFixture() {
	// Explicit dispatch outside the ordinary test inventory: neither PASS nor skip.
	mode, run := os.Args[len(os.Args)-2], os.Args[len(os.Args)-1]
	ctx, closeStop, err := nativeStopContext(context.Background(), os.Getenv("ELITE_STOP_EVENT_HANDLE"))
	if err != nil {
		exitOnHostFailure(err)
		return
	}
	store := &lifecycleStore{mode: mode, run: run}
	providers := map[string]refundworker.Provider{"stripe": lifecycleProvider{mode, run}}
	if mode == "provider-missing" {
		providers = nil
	}
	if mode == "report-failure" {
		f, e := os.CreateTemp(".", "closed-log-")
		if e != nil {
			panic(e)
		}
		if e = f.Close(); e != nil {
			panic(e)
		}
		log.SetOutput(f)
	}
	processor, err := refundworker.NewProcessor(store, providers, "synthetic-worker", time.Minute, time.Second, 5)
	if err != nil {
		panic(err)
	}
	loopErr := runRefundLoop(ctx, processor.Step, claimToken)
	cause := context.Cause(ctx)
	closeErr := closeStop()
	if cause != errNativeStop && cause != context.Canceled {
		loopErr = errors.Join(loopErr, errNativeProtocol)
	}
	if err := fixturePublish("lifecycle-ack.json", map[string]any{"run_id": run, "loop_returned": true, "stop_observed": cause == errNativeStop, "claims": store.claims, "loop_failed": loopErr != nil}); err != nil {
		loopErr = errors.Join(loopErr, err)
	}
	exitOnHostFailure(errors.Join(loopErr, closeErr))
}
````

## Qualification source: test_worker_lifecycle.py (AUTHORED / NOT_ADMITTED)

SHA-256: `784320aeb60b20aa61c2a2d3240d9598b64a46384ae5f0ffa66230e988080b3e`

````python
"""Real canonical host+processor, synthetic Store/Provider only. NOT_ADMITTED."""
from pathlib import Path
import hashlib,json,os,sys,threading,time,unittest,uuid
from unittest.mock import patch
import result_store as s
ROOT=Path(__file__).resolve().parent/'worker-runs';ROOT.mkdir(exist_ok=True)
FIX=json.loads(Path(sys.argv[2]).read_text());EXE=Path(FIX['path'])
class Tests(unittest.TestCase):
 def setUp(self):
  self.base=ROOT/uuid.uuid4().hex;self.base.mkdir();self.store=self.base/'store';self.store.mkdir();self.run=uuid.uuid4().hex
  self.assertEqual(s.digest(EXE.read_bytes()),FIX['sha256'])
 def profile(self,mode,grace=1):
  p={'schema':'elite-native-launch-qualification/v2','id':'real-loop-synthetic-domain','executable':{'path':str(EXE),'bytes':EXE.stat().st_size,'sha256':FIX['sha256']},'arguments':['-test.run=^TestNativeLifecycleFixture$','--','--fixture',mode,self.run],'cwd':str(self.base),'environment':{'SYSTEMROOT':os.environ['SYSTEMROOT'],'WINDIR':os.environ['WINDIR'],'TEMP':str(self.base),'TMP':str(self.base)},'budgets':{'timeout':8,'output_bytes':4096,'processes':8,'commit_bytes':536870912},'shutdown':{'protocol':'win32-inherited-event/v1','grace_seconds':grace}}
  raw=s.encode(p);return raw,s.digest(raw)
 def execute(self,mode,grace=1):
  raw,digest=self.profile(mode,grace);event=threading.Event();ready=[]
  def stop_ready():
   end=time.monotonic()+7
   while time.monotonic()<end:
    if (self.base/'ready').exists():ready.append(True);break
    time.sleep(.005)
   event.set()
  th=threading.Thread(target=stop_ready);th.start()
  try:r=s.execute(self.store,self.run,raw,digest,cancel=event)
  finally:th.join(8)
  self.assertFalse(th.is_alive());self.assertEqual(ready,[True]);self.assertEqual(r.state,'RECORDED');self.assertEqual(s.inspect(self.store,self.run,digest),r)
  with patch.object(s.g,'launch') as launch:self.assertEqual(s.execute(self.store,self.run,raw,digest).state,'ALREADY_RESERVED');launch.assert_not_called()
  v=json.loads((self.store/self.run/'result.json').read_text());o=v['outcome'];self.assertTrue(o['tree_empty']);self.assertEqual(o['status'],'CANCELLED');self.assertFalse((self.base/'second-claim').exists());return o
 def ack(self):
  a=json.loads((self.base/'lifecycle-ack.json').read_text());self.assertEqual(a['run_id'],self.run);self.assertTrue(a['loop_returned'] and a['stop_observed']);self.assertEqual(a['claims'],1);return a
 def domain(self):
  a=json.loads((self.base/'domain-result.json').read_text());self.assertEqual(a['schema'],'synthetic-refund-outcome/v1');self.assertEqual(a['run_id'],self.run);self.assertEqual(a['claims'],1);return a
 def clean_exit(self,o):self.assertEqual((o['exit_code'],o['shutdown_state']),(0,'EXITED_DURING_GRACE'));self.ack()
 def test_succeeded_domain_is_read_separately_from_process_receipt(self):
  self.clean_exit(self.execute('succeeded'));self.assertEqual(self.domain()['outcome'],'succeeded')
 def test_pending_domain_despite_clean_process_exit(self):
  self.clean_exit(self.execute('pending'));d=self.domain();self.assertEqual((d['outcome'],d['code']),('retry','PROVIDER_PENDING'))
 def test_failed_domain_despite_clean_process_exit(self):
  self.clean_exit(self.execute('failed'));self.assertEqual(self.domain()['outcome'],'failed')
 def test_missing_provider_stays_blocked_without_external_calls(self):
  self.clean_exit(self.execute('provider-missing'));d=self.domain();self.assertEqual((d['outcome'],d['code']),('blocked','PROVIDER_CONFIG_MISSING'))
 def test_effect_before_cancellation_is_unresolved_and_never_replayed(self):
  self.clean_exit(self.execute('effect-then-cancel'));self.assertTrue((self.base/'synthetic-effect.json').exists());self.assertFalse((self.base/'domain-result.json').exists())
 def test_lost_domain_ack_recovered_by_separate_record_read(self):
  self.clean_exit(self.execute('lost-domain-ack'));self.assertEqual(self.domain()['outcome'],'succeeded')
 def test_domain_store_failure_is_not_process_success(self):
  self.clean_exit(self.execute('store-failure'));self.assertFalse((self.base/'domain-result.json').exists())
 def test_report_failure_preserves_host_exit_one_and_unknown_domain(self):
  o=self.execute('report-failure');self.assertEqual((o['exit_code'],o['shutdown_state']),(1,'EXITED_DURING_GRACE'));self.assertTrue(self.ack()['loop_failed']);self.assertFalse((self.base/'domain-result.json').exists())
 def test_noncooperative_provider_is_forced_without_domain_ack(self):
  o=self.execute('ignore-stop',.1);self.assertEqual(o['shutdown_state'],'FORCED');self.assertFalse((self.base/'lifecycle-ack.json').exists());self.assertFalse((self.base/'domain-result.json').exists());self.assertTrue((self.base/'synthetic-effect.json').exists())
 def test_private_domain_errors_are_absent_from_process_store(self):
  self.clean_exit(self.execute('lost-domain-ack'))
  for p in (self.store/self.run).iterdir():self.assertNotIn(b'PRIVATE_FIXTURE',p.read_bytes());self.assertNotIn(b'synthetic-pi',p.read_bytes())

if __name__=='__main__':
 suite=unittest.defaultTestLoader.loadTestsFromTestCase(Tests);result=unittest.TextTestRunner(verbosity=2).run(suite)
 Path(sys.argv[1]).write_text(json.dumps({'pass':result.wasSuccessful(),'tests_run':result.testsRun,'skips':len(result.skipped),'failures':len(result.failures),'errors':len(result.errors),'root':str(ROOT)},indent=2)+'\n')
 sys.exit(0 if result.wasSuccessful() else 1)
````

## Retained fixture correction — FAIL604

The retained open-pipe descendant test passed its assertions but immediate
directory teardown failed32, as in the earlier silent-descendant case. This
successor source extends the exact owned-process membership/termination proof and
finite directory-release check to that sibling.22other tests passed unchanged.
V348 source/history remains intact; seven prior native files are byte-identical,
while this test source is strengthened. Native runtime/product code is unchanged.

## Qualification source: test_capture.py (AUTHORED / NOT_ADMITTED)

SHA-256: `8a78f1c54d788edb18e5d97a1cb2f4ca8c490af5adb0db0aaae2bad28873bf38`

````python
from pathlib import Path
import ctypes as C
from ctypes import wintypes as W
import hashlib,json,os,sys,time,threading,unittest,tempfile,subprocess
from unittest.mock import patch
import native_capture as n
from native_capture import win

class CaptureTests(unittest.TestCase):
    def setUp(self):
        self.temp=tempfile.TemporaryDirectory();self.root=Path(self.temp.name).resolve()
    def tearDown(self):self.temp.cleanup()
    def call(self,code,**kw):return n.capture([sys.executable,'-I','-c',code],self.root,**kw)
    def assert_complete(self,result,code=0):
        self.assertEqual('COMPLETED',result.status);self.assertEqual(code,result.exit_code);self.assertTrue(result.tree_empty)
    def test_separate_binary_streams_at_exact_limit(self):
        result=self.call("import os;[(os.write(1,bytes(range(256))*16),os.write(2,bytes(reversed(range(256)))*16)) for _ in range(256)]",limit=1024*1024)
        self.assert_complete(result)
        self.assertEqual(bytes(range(256))*4096,result.stdout);self.assertEqual(bytes(reversed(range(256)))*4096,result.stderr)
    def test_each_stream_over_limit_stops_tree_before_later_marker(self):
        for fd in (1,2):
            with self.subTest(fd=fd):
                result=self.call(f"import os,time,pathlib;os.write({fd},b'x'*65536);time.sleep(.4);pathlib.Path('later').write_text('bad')",limit=1024)
                self.assertEqual('OUTPUT_LIMIT',result.status);self.assertTrue(result.tree_empty);self.assertLessEqual(len(result.stdout),1024);self.assertLessEqual(len(result.stderr),1024)
                self.assertFalse((self.root/'later').exists())
    def test_stdin_eof_and_unselected_environment_absent(self):
        with patch.dict(os.environ,{'ELITE_CAPTURE_PRIVATE':'DO_NOT_INHERIT'}):
            result=self.call("import sys,os;assert sys.stdin.buffer.read()==b'';assert 'ELITE_CAPTURE_PRIVATE' not in os.environ;print('ok')")
        self.assert_complete(result);self.assertIn(b'ok',result.stdout)
    def test_nonzero_and_still_active_numeric_exit_preserved(self):
        for code in (7,259):
            with self.subTest(code=code):self.assert_complete(self.call(f'raise SystemExit({code})'),code)
    def test_timeout_cleans_child_and_descendant_with_open_pipe(self):
        nested="import os,time,pathlib;pathlib.Path('nested-ready').write_text(str(os.getpid()));time.sleep(20)"
        original=n.OwnedCapture.close;observed=[]
        def close_with_owned_descendant(owner):
            handle=win.open_process(0x100000|0x1000,False,int((self.root/'nested-ready').read_text()))
            self.assertTrue(handle)
            try:
                membership=W.BOOL();self.assertTrue(win.in_job(handle,owner.job,C.byref(membership)));self.assertTrue(membership.value)
                self.assertEqual(win.wait(handle,0),258)
                result=original(owner)
                self.assertEqual(win.wait(handle,5000),0);observed.append('owned-descendant-signaled')
                return result
            finally:win.close(handle)
        with patch.object(n.OwnedCapture,'close',close_with_owned_descendant):
            result=self.call(f"import subprocess,sys,time;subprocess.Popen([sys.executable,'-I','-c',{nested!r}]);time.sleep(20)",timeout=.6)
        self.assertEqual(observed,['owned-descendant-signaled'])
        self.assertTrue((self.root/'nested-ready').exists());self.assertEqual('TIMED_OUT',result.status);self.assertTrue(result.tree_empty)
        # Native tree accounting is distinct from post-termination file release.
        # Observe only this owned fixture, under a finite deadline, after the
        # exact descendant handle has signaled. Never waive expiry or a live child.
        until=time.monotonic()+5;moved=self.root.with_name(self.root.name+'-released')
        while True:
            try:self.root.rename(moved);break
            except OSError as error:
                if error.winerror not in (5,32) or time.monotonic()>=until:raise
                time.sleep(.005)
        moved.rename(self.root)
    def test_silent_descendant_cannot_be_success_when_root_exits(self):
        nested="import time,pathlib,os;pathlib.Path('silent-ready').write_text(str(os.getpid()));time.sleep(20)"
        original=n.OwnedCapture.close
        observed=[]
        def close_with_owned_descendant(owner):
            pid=int((self.root/'silent-ready').read_text())
            handle=win.open_process(0x100000|0x1000,False,pid)
            self.assertTrue(handle)
            try:
                membership=W.BOOL();self.assertTrue(win.in_job(handle,owner.job,C.byref(membership)));self.assertTrue(membership.value)
                self.assertEqual(win.wait(handle,0),258)
                result=original(owner)
                self.assertEqual(win.wait(handle,5000),0)
                observed.append('owned-descendant-signaled')
                return result
            finally:win.close(handle)
        with patch.object(n.OwnedCapture,'close',close_with_owned_descendant):
            result=self.call(f"import subprocess,sys;subprocess.Popen([sys.executable,'-I','-c',{nested!r}],stdout=subprocess.DEVNULL,stderr=subprocess.DEVNULL)",timeout=.6)
        self.assertEqual(observed,['owned-descendant-signaled'])
        self.assertTrue((self.root/'silent-ready').exists());self.assertEqual('TIMED_OUT',result.status);self.assertTrue(result.tree_empty)
        # Fixture teardown is a separate resource-release observation, not an
        # inference from job accounting. Retry only this exact owned rename,
        # after proving descendant termination, under a finite deadline.
        until=time.monotonic()+5;moved=self.root.with_name(self.root.name+'-released')
        while True:
            try:self.root.rename(moved);break
            except OSError as error:
                if error.winerror not in (5,32) or time.monotonic()>=until:raise
                time.sleep(.005)
        moved.rename(self.root)
    def test_finite_descendant_output_after_root_exit_is_collected(self):
        nested="import time,os;time.sleep(.1);os.write(1,b'late-child');os.write(2,b'late-error')"
        result=self.call(f"import subprocess,sys;subprocess.Popen([sys.executable,'-I','-c',{nested!r}])")
        self.assert_complete(result);self.assertEqual(b'late-child',result.stdout);self.assertEqual(b'late-error',result.stderr)
    def test_pre_cancel_does_not_launch(self):
        cancel=threading.Event();cancel.set()
        with patch.object(n,'OwnedCapture') as launch:result=self.call('raise SystemExit(0)',cancel=cancel)
        launch.assert_not_called();self.assertEqual('CANCELLED',result.status);self.assertTrue(result.tree_empty)
    def test_active_cancel_cleans_tree(self):
        cancel=threading.Event()
        def trigger():
            until=time.monotonic()+5
            while not (self.root/'ready').exists() and time.monotonic()<until:time.sleep(.005)
            cancel.set()
        thread=threading.Thread(target=trigger);thread.start()
        try:result=self.call("from pathlib import Path;import time;Path('ready').write_text('yes');time.sleep(20)",cancel=cancel)
        finally:thread.join(timeout=6)
        self.assertFalse(thread.is_alive());self.assertEqual('CANCELLED',result.status);self.assertTrue(result.tree_empty)
    def test_read_failure_is_static_and_cleans(self):
        read=os.read
        def injected(fd,count):raise OSError('PRIVATE_PATH_AND_TOKEN')
        with patch.object(n.os,'read',side_effect=injected):result=self.call('import time;time.sleep(20)')
        self.assertEqual('IO_FAILED',result.status);self.assertTrue(result.tree_empty);self.assertNotIn('PRIVATE_PATH_AND_TOKEN',repr(result))
    def test_unlisted_inheritable_handle_is_not_in_child(self):
        # Give the child the numeric handle of an intentionally inheritable file.
        # GetFinalPathName distinguishes accidental handle-value reuse from leak.
        sentinel=self.root/'not-inherited.txt';sentinel.write_text('PRIVATE_SENTINEL')
        fd=os.open(sentinel,os.O_RDONLY|os.O_BINARY);os.set_inheritable(fd,True);h=msvcrt_handle(fd)
        code=f"""import ctypes as c
from ctypes import wintypes as w
k=c.WinDLL('kernel32',use_last_error=True)
f=k.GetFinalPathNameByHandleW;f.argtypes=[w.HANDLE,w.LPWSTR,w.DWORD,w.DWORD];f.restype=w.DWORD
b=c.create_unicode_buffer(4096);r=f({h},b,len(b),0)
assert r==0 or 'not-inherited.txt' not in b.value
print('excluded')
"""
        try:result=self.call(code)
        finally:os.close(fd)
        self.assert_complete(result);self.assertIn(b'excluded',result.stdout)
    def test_concurrent_capture_has_isolated_streams_and_eof(self):
        results=[];errors=[]
        def child():
            try:results.append(self.call("import os;os.write(1,b'own');os.write(2,b'err')"))
            except BaseException as e:errors.append(type(e).__name__)
        threads=[threading.Thread(target=child) for _ in range(16)]
        for thread in threads:thread.start()
        for thread in threads:thread.join(timeout=8)
        self.assertFalse(errors);self.assertTrue(all(not t.is_alive() for t in threads));self.assertEqual(16,len(results))
        for result in results:self.assert_complete(result);self.assertEqual(b'own',result.stdout);self.assertEqual(b'err',result.stderr)
    def test_bad_limits_fail_before_launch(self):
        for key,values in [('timeout',[True,0,float('inf'),float('nan'),7201]),('limit',[True,0,10485761]),('max_processes',[True,0,65]),('max_memory',[True,0,1073741825])]:
            for value in values:
                with self.subTest(key=key,value=value),patch.object(n,'OwnedCapture') as launch,self.assertRaises(ValueError):self.call('pass',**{key:value})
                launch.assert_not_called()
    def test_launch_failure_does_not_create_child_effect(self):
        result=n.capture([str(self.root/'nonexistent.exe')],self.root)
        self.assertEqual('LAUNCH_FAILED',result.status);self.assertEqual(b'',result.stdout);self.assertEqual(b'',result.stderr)
    def test_cleanup_waits_after_accounting_empty_without_retermination(self):
        owner=n.OwnedCapture([sys.executable,'-I','-c','import time;time.sleep(20)'],self.root,n.clean_environment());owner.start()
        actual_wait=win.wait
        def delayed(handle,milliseconds):
            # A nonblocking wait could see the transitional still-active state;
            # bounded wait must wait for the already-requested kernel cleanup.
            if handle==owner.pi.process and milliseconds==0:return 258
            return actual_wait(handle,milliseconds)
        try:
            with patch.object(win,'wait',side_effect=delayed),patch.object(win,'terminate',wraps=win.terminate) as retry:
                self.assertTrue(owner.close())
                retry.assert_not_called()
            self.assertTrue(owner.empty)
        finally:owner.close()

    def test_real_refund_binary_retains_exit_and_fixed_code(self):
        exe=Path(sys.argv[2]).resolve(strict=True)
        expected=sys.argv[3]
        self.assertEqual(expected,hashlib.sha256(exe.read_bytes()).hexdigest())
        # Explicit hashed local fixture; no DB/provider environment is inherited.
        result=n.capture([str(exe)],self.root)
        self.assert_complete(result,1);self.assertIn(b'REFUND_HOST_FAILED',result.stderr)
    def test_job_commit_budget_refuses_large_allocation(self):
        budget=64*1024*1024
        code="try:\n value=bytearray(128*1024*1024)\nexcept MemoryError:\n print('MEMORY_DENIED')\nelse:\n raise SystemExit(17)"
        result=self.call(code,max_memory=budget)
        self.assert_complete(result);self.assertIn(b'MEMORY_DENIED',result.stdout)
        self.assertIsNotNone(result.observed_peak_job_bytes);self.assertGreater(result.observed_peak_job_bytes,0)

    def test_native_commit_denied_and_granted_by_selected_budget(self):
        code="""import ctypes as c,json
from ctypes import wintypes as w
k=c.WinDLL('kernel32',use_last_error=True)
a=k.VirtualAlloc;a.argtypes=[c.c_void_p,c.c_size_t,w.DWORD,w.DWORD];a.restype=c.c_void_p
f=k.VirtualFree;f.argtypes=[c.c_void_p,c.c_size_t,w.DWORD];f.restype=w.BOOL
large=a(None,128*1024*1024,0x3000,4);error=c.get_last_error()
if large:assert f(large,0,0x8000)
small=a(None,1024*1024,0x3000,4);assert small;assert f(small,0,0x8000)
print(json.dumps({'large_granted':bool(large),'small_granted':bool(small),'large_error':error if not large else None}))
"""
        for mib,allowed in [(64,False),(256,True)]:
            with self.subTest(budget_mib=mib):
                result=self.call(code,max_memory=mib*1024*1024);self.assert_complete(result)
                self.assertEqual({'large_granted':allowed,'small_granted':True,'large_error':None if allowed else 1455},json.loads(result.stdout))
                self.assertGreater(result.observed_peak_job_bytes,0)

    def test_failed_pipe_setup_releases_created_handles(self):
        current=win.api('GetCurrentProcess',[],W.HANDLE)
        count=win.api('GetProcessHandleCount',[W.HANDLE,C.POINTER(W.DWORD)],W.BOOL)
        def observed():
            value=W.DWORD();self.assertTrue(count(current(),C.byref(value)));return value.value
        before=observed();pipe=n.create_pipe;calls=0
        def fail_second(*args):
            nonlocal calls
            calls+=1
            return 0 if calls==2 else pipe(*args)
        with patch.object(n,'create_pipe',side_effect=fail_second),patch.object(win,'create') as launch:
            result=self.call('pass')
        self.assertEqual('LAUNCH_FAILED',result.status);launch.assert_not_called();self.assertEqual(before,observed())

    def test_cleanup_failure_is_not_a_success_receipt(self):
        with patch.object(n,'terminate_job',return_value=0):
            result=self.call('import time;time.sleep(20)',timeout=.1)
        # Closing the owned job still requests termination, but failed explicit
        # cleanup must stay visible. No false proof of an empty tree is emitted.
        self.assertEqual('CLEANUP_FAILED',result.status);self.assertFalse(result.tree_empty)

    def test_active_process_budget_refuses_descendant_creation(self):
        child="from pathlib import Path;Path('quota-escape').write_text('bad')"
        code="import subprocess,sys,json\ntry:\n subprocess.Popen([sys.executable,'-I','-c',"+repr(child)+"])\nexcept OSError as error:\n print(json.dumps({'denied':error.winerror}))\nelse:\n raise SystemExit(17)"
        result=self.call(code,max_processes=1)
        self.assert_complete(result);self.assertEqual({'denied':1816},json.loads(result.stdout));self.assertFalse((self.root/'quota-escape').exists())

    def test_abrupt_owner_death_with_redirected_pipes_kills_tree(self):
        module_dir=Path(__file__).resolve().parent
        script=self.root/'owner.py'
        nested="import os,time,pathlib;pathlib.Path('nested.json').write_text(str(os.getpid()));time.sleep(20)"
        child=f"import os,subprocess,sys,time,pathlib;subprocess.Popen([sys.executable,'-I','-c',{nested!r}],creationflags=0x08000000)\nuntil=time.monotonic()+5\nwhile not pathlib.Path('nested.json').exists() and time.monotonic()<until:time.sleep(.005)\nos.write(1,b'x'*65536);time.sleep(20)"
        script.write_text("import sys,os,time,json\nfrom pathlib import Path\nsys.path.insert(0,sys.argv[1]);import native_capture as n\nroot=Path(sys.argv[2])\nowner=n.OwnedCapture([sys.executable,'-I','-c',"+repr(child)+"],root,n.clean_environment());owner.start()\nuntil=time.monotonic()+5\nwhile not (root/'nested.json').is_file() and time.monotonic()<until:time.sleep(.005)\nn.win.write(root/'owner-ready.json',{'root':owner.pi.pid,'nested':int((root/'nested.json').read_text())})\nn.win.await_file(root/'exit.json');os._exit(0)\n",encoding='utf-8')
        wrapper=subprocess.Popen([sys.executable,'-I',str(script),str(module_dir),str(self.root)],stdin=subprocess.DEVNULL,stdout=subprocess.DEVNULL,stderr=subprocess.DEVNULL,creationflags=subprocess.CREATE_NO_WINDOW,env=n.clean_environment())
        handles=[]
        try:
            info=win.await_file(self.root/'owner-ready.json')
            for key in ('root','nested'):
                handle=win.open_process(0x100000|0x1000,False,info[key]);self.assertTrue(handle);handles.append(handle);self.assertEqual(258,win.wait(handle,0))
            win.write(self.root/'exit.json',{'exit':True});wrapper.wait(timeout=5)
            self.assertEqual(0,wrapper.returncode)
            for handle in handles:self.assertEqual(0,win.wait(handle,5000))
        finally:
            if wrapper.poll() is None:wrapper.kill();wrapper.wait(timeout=5)
            for handle in handles:win.close(handle)

    def test_owned_handles_stable_after_repeated_runs(self):
        current=win.api('GetCurrentProcess',[],W.HANDLE)
        count=win.api('GetProcessHandleCount',[W.HANDLE,C.POINTER(W.DWORD)],W.BOOL)
        def observed():
            v=W.DWORD();self.assertTrue(count(current(),C.byref(v)));return v.value
        self.assert_complete(self.call('pass'));before=observed()
        for _ in range(24):self.assert_complete(self.call('pass'))
        self.assertEqual(before,observed())

def msvcrt_handle(fd):
    import msvcrt
    return msvcrt.get_osfhandle(fd)

if __name__=='__main__':
    outcomes=[]
    class Result(unittest.TextTestResult):
        def addSubTest(self,test,subtest,error):
            outcomes.append({'id':subtest.id(),'status':'PASS' if error is None else 'FAIL','subtest':True});super().addSubTest(test,subtest,error)
        def addSuccess(self,test):outcomes.append({'id':test.id(),'status':'PASS'});super().addSuccess(test)
        def addError(self,test,error):outcomes.append({'id':test.id(),'status':'ERROR'});super().addError(test,error)
        def addFailure(self,test,error):outcomes.append({'id':test.id(),'status':'FAIL'});super().addFailure(test,error)
    suite=unittest.defaultTestLoader.loadTestsFromTestCase(CaptureTests)
    result=unittest.TextTestRunner(verbosity=2,resultclass=Result).run(suite)
    Path(sys.argv[1]).write_text(json.dumps({'tests_run':result.testsRun,'pass':result.wasSuccessful(),'skips':len(result.skipped),'outcomes':outcomes},indent=2)+'\n')
    sys.exit(not result.wasSuccessful())
````

## Canonical integration143 — observed 2026-09-09T14:06:49.641478Z

All15canonical refund files and seven reused native files are unchanged. Five
qualification files reconstruct exactly, including the strengthened capture test.
Ten new native Go tests pass three repetitions; ten integrated worker/process-store
tests pass. The retained90supervisor tests pass with no skips. Go baseline passes
three repetitions, with 1 distinct opt-in test(s) skipped: TestPostgresRefundLineAllocationPartialThenFullAndAmbiguity.
Those skips remain unexecuted, never target evidence. Vet and fixture build pass.
Finite handle-parser fuzz:10seeds, 2922431 executions,10second budget, PASS.
Fuzz gate self-test1positive/2negative and gap self-test3positive/6negative PASS.
No new dependency/runtime/product pack/profile. Full Preflight143 is pending.

The native adapter explicitly mounted the real loop and processor through a
qualification TestMain. The normal startup run(), actual PostgreSQL Store and
provider SDK transports remain unmounted in this native qualification. Synthetic
acknowledgements support test assertions only; they are not financial confirmation.
Broad service/runtime/security/storage/target admission is still CONDITIONED.

| Evidence | SHA-256 |
|---|---|
| worker-lifecycle.json | `2cd6d80bcffced10c717d9eaa5c5e38dfb60a4f49008fb392324472f95a2a525` |
| worker-lifecycle.log | `e16bf2bc73a2ce318cf7260ed4dfbf77f0d51ce05bb747643418c3952345cb6b` |
| worker-lifecycle-green.json | `c66ab8c80396e6c9a09b6780fbf8cc7d5b92f7ebc68305c0133577f70fa87f9d` |
| presignaled-red.log | `2661aacd30ebba2eaebb56ba7e6fa2e2461e633eb55342207c18cea0d5f1eabc` |
| rebuilt-go-tests.log | `272c76f8d287237f325063606f77f8f9add3453483a31469aa4fb399e8f36f41` |
| rebuilt-vet.log | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| rebuilt-worker-fixture.json | `0341a7eaf29525cce76c4f181a3ea66be22e87b5ea23335cd680ee651cd432be` |
| rebuilt-worker-lifecycle.json | `10bab45485f70ebbb99c0d4ac28364d02ef880930ba3587ae7e502c14784d0d1` |
| rebuilt-result.json | `a0d6088680bcaca86b4f5fe5901e8d704ca69d3caef80249bdafd68861e7537e` |
| rebuilt-shutdown.json | `6b0586d112235212d1a8b31b663dd99b84706630eb4de97ea6e80463f6e4c549` |
| rebuilt-launch.json | `9dc3c8efed68851cf6170b0706956b637ab147fd4a1987fce53099d9da149d36` |
| rebuilt-capture.json | `422a65991f16440e67896450afba46df1025150749f10bbc7eb6c42d8e80998f` |
| rebuilt-capture.log | `de5a535a53db0951abc54d738c303d6b101f946ff5053fd8e938bd0ba8a485ef` |
| rebuilt-capture-green.json | `c4b869ed05b0f2082f785da240aab9346736b414b47724192311db00336fa798` |
| canonical-parity.json | `383beda4f7d404608185f7e0ba60e0d24fe9955284de66e07f4ce5df375ca2b4` |
| test-summary.json | `7e96568e1d96c21033e5938d33b6df07c888b1048d82f84acb75dd3171649113` |
| fuzz-profile.json | `ac3b23c5be632dd9f83f3bec91849656b0bf5b723dfabf2dab3ba146360ab7f6` |
| fuzz-receipt.json | `a49e7404b004634be6ff26995843cc7bc0da784b40ba951398e0a6a562d1bcbc` |
| fuzz.log | `91c1ff409325a49863167717d996259554278c04845e768dd6dc21909050073a` |

## Structural inventory correction144 — FAIL605 / LIB2321

Preflight143 rejected stale current inventory787Markdown; actual expected788.
Original failed log SHA-256 `ba2c762611decb7a545389e17fe5be62fed8e6e76c47f64daedabfccecafa9f9`. No successful
Preflight143 receipt exists. Corrected the current canonical inventory header to
162/1461/788/53, preserving historical inventories. The full successor
Preflight144 must pass; code/source/test evidence is unchanged.
## Closure145 — local lifecycle integration qualified

Observed 2026-09-09T14:17:50.458651Z. Preflight144: 154 executed checks PASS; inventory
162packs/1461materializable files/788Markdown/53profiles. Docker absent; opt-in
network/provider/PostgreSQL/target checks remain unexecuted where gated. This
supersedes the pending Preflight144 snapshot above, not historical failure evidence.

Five source blocks reconstruct exactly, seven prior native files are unchanged
and all15refund product files remain exact.20new lifecycle tests and90retained
supervisor tests pass; Go baseline runs three repetitions, vet/build and finite
fuzz2922431executions PASS. The current defect corrections are regression-proven
or recovered in this limited scope; ledger2530unique IDs retains original findings.

The demonstrated chain is native stop → canonical host/processor with synthetic
Store/Provider → separate lifecycle/domain fixture observations → no-replay process
receipt and strict readback. Normal run() startup/config and real PostgreSQL-backed
domain acknowledgement are the next integration work. No product service pack,
financial account, automatic restart/retry or full target admission is implied.

Final successor record/receipt: gap-record145.json and gap-receipt145.json in this
stage. Full service remains CONDITIONED and global roadmap remains incomplete.

| Closure evidence | SHA-256 |
|---|---|
| preflight144.json | `c4ca155b9c3f9589c487eb0ff456d0aac7bbe2b1d7d9cdbe5ae24f1ece0aac27` |
| preflight144.log | `12af763561d2334c76f657bfeb345ad346f1cb51362234b4835a1219395ada96` |
| structural-closure144.json | `8ded7fc00a2158958c5b552c96e1b6be3d7be912f47f3a442abf16173c234f63` |
| gap-record144.json | `bc731366768382cf0d7b521baf2140c661a6bd038e0d578fd9e50219c6da5795` |
| gap-receipt144.json | `9f435c2f8522ab6fbd2947f12a3d5cd74f26d34576d8112d4e5a0a04f5291c4b` |

## Compact cursor repair — FAIL606 / LIB2322

Full Preflight144154checks PASS remains unchanged. Checkpoint creation rejected
the prepared next_action because it exceeded800characters; no EVT0145 was written.
The rejected state is retained as rejected-state145.json. Detailed routing belongs
here, while the cursor stays compact. This is metadata repair, not a code change.

The legacy TestPostgresRefundLineAllocationPartialThenFullAndAmbiguity reads
TEST_DATABASE_URL and only accepts a dedicated loopback elite_refund_test_<unique>
database. It bypasses relational triggers during fixture seeding/cleanup: never
run it against a project/audit database or treat it as originating workflow proof.
The next native integration needs actual guarded PostgreSQL-backed reconciliation
and normal startup/config, preserving the source's narrow compatibility scope.

Final successor decision/receipt: gap-record145-final.json and
gap-receipt145-final.json. Earlier145 record/receipt remains the pre-repair snapshot.
Ledger now2531unique IDs; original605/606 rejection evidence remains preserved.
