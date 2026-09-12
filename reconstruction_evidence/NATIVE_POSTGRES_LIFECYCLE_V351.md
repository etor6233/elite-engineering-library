# V351 — PostgreSQL-backed native worker lifecycle and closure accounting

Library maintenance / T2809 resumed from validated checkpoint145. AUTHORED local
qualification, CONDITIONED / NOT_ADMITTED as a product service. The refund pack
and selected67pack/746file profile remain unchanged. Candidate main.go overlays
only the existing signal context with the V350 native-stop context and propagates
watcher/cleanup errors. runRefundLoop and the complete PostgresStore/Processor
implementations are unchanged. No new dependencies or payment-provider traffic.

## How far is the complete library from finished?

The maintenance contract records41passed,5blocked and2planned tests out of48:
85.4167% registered checks passed and14.5833% not passed. These are counts, NOT
percentages of total engineering effort, elapsed time or full-library completion.
The roadmap explicitly has1complete owner task and10open tasks with partial work.
No measured effort weights or complete remaining-work inventory supports an exact
global percentage or a reliable delivery date. Local tests added to TEST05 evidence
do not close that broad integration gate or inflate the contract denominator.

| Remaining contract gate | Concrete remaining closure |
|---|---|
| TEST02 | admission/equivalence of21cores by exact claim, source, license and evidence |
| TEST03 | integrated security/privacy/supply-chain acceptance of the actual composition |
| TEST05 | complete runtime identity, reporting, retention, alerts and operational evidence |
| TEST06 | partial-install recovery from the exact distributed candidate |
| TEST07 | final independent reproducible archive, SCA/notices/SBOM/signature acceptance |
| TEST08 | NEW/EXISTING usability/help/recovery on that exact final artifact |
| TEST09 | authorized historical corpus plus identified channel/format and import proof |

BENCH01 and readinessD–H also remain open. A user question asks which historical
channel and format are required; no missing answer/corpus is inferred. RoundD
was rendered in an isolated readiness copy and matches the existing root advisory
exactly (1119f31c5085584724576dbb6f06c2afecf0b3093cc5b01712d701a50a1281e0).
Raw originals and derivative lineage remain required; no importer was invented.

## Concrete PostgreSQL connection

The native supervisor launches a hash-pinned Go qualification binary. The startup
lane invokes run() itself: existing environment validation, database Ping,
provider absence, PostgresStore construction and the real polling loop. Only the
Windows stop wiring is overlaid. Without provider credentials the durable request
must be blocked with PROVIDER_CONFIG_MISSING, with no synthetic provider effect.
The test wrapper loads a bounded, hash-bound, exact loopback-only DSN fixture into
the child's environment. It is NOT a production secret manager/configuration API,
authenticated service identity or a deployed service. The standalone main binary
is built; startup-lane evidence is explicitly through the qualification TestMain.

Other lanes combine the actual loop, Processor and PostgresStore with a synthetic
Provider. Its independent qualification.effect insert models one observed effect,
never a real payment. Success, pending and failure are read back from actual domain
tables. Lost post-commit acknowledgement is reconciled by database read. A failure
in the immutable-observation insert proves rollback of that Complete transaction.
Cancellation after an effect and forced shutdown retain claimed/prepared state:
neither a clean process exit nor an absent domain result authorizes replay.

The existing process store publishes metadata separately from PostgreSQL state.
Its publication can fail after the domain commit; same-ID retries are refused and
database reads recover the committed result without a second synthetic effect.
Config digest mismatch and pre-cancellation prevent the first claim. Each case
gets a distinct database cloned from a fresh, private, loopback-only fixture.

Seeded captured order and authorized return are synthetic starting conditions,
not proof of an originating sales/approval/financial workflow. Every canonical
constraint and trigger stays enabled; session_replication_role stays origin.
There is no cleanup that deletes immutable business rows or disables triggers:
the entire owned cluster is retained outside the library after it is stopped.
The legacy trigger-bypassing compatibility test remains skipped in ordinary Go
baseline and is not used as new integration evidence.

Controlled PostgreSQL stop/start compares complete ordered snapshots of seven
tables across11case databases. This is local process restart with fsync,
synchronous_commit and full_page_writes on, not PITR/backup restore, power-cycle,
replication, exactly-once external payment processing or production DR.

## Observed corrections and boundaries

FAIL607/LIB2323: repeated guessed recovery paths corrected from actual inventory;
no result inferred from failed reads. Discovery-dependent reads must be sequential.
FAIL608/LIB2324: installed PostgreSQL18.6 lacks share/timezone/UTC and rejects an
explicit UTC setting. Preserve the original startup failure; keep initdb's observed
America/Buenos_Aires setting and record it. This does not admit complete timezone
coverage or alter/fetch the upstream runtime. Production locale/runtime remains a gate.
FAIL609/LIB2325: initial10s shutdown budget expired during a14.885s checkpoint,
flushing12841files. Cleanup subsequently completed. Fixture now has an explicit40s
stop budget (45scommand bound); this is not a production shutdown SLO.
FAIL610/LIB2326: third run compared post-restart snapshots equal, then incorrectly
used the16KiB process-metadata encoder for the larger evidence snapshot. Separate
canonical JSON hashing fixes the evidence layer; result-store limits are unchanged.
FAIL611/LIB2327: renderer rejected an absolute output reference before writing;
isolated canonical-relative output re-render matches the historical roundD prompt.
All failed logs/sources/receipts are retained. Fresh canonical gates are pending.

## Authority and reproduction

Authorities: DATABASE_STORAGE_INTERNALS.md (WAL/commit/recovery distinctions),
SECURITY_SRE_CLOUD_INFRASTRUCTURE.md, PUBLIC_CODE_ARCHITECTURE_ADMISSION_STANDARD.md
and ENGINEERING_EXECUTION_PLAYBOOK.md. Official PostgreSQL18
[WAL settings](https://www.postgresql.org/docs/18/runtime-config-wal.html) and
[initdb](https://www.postgresql.org/docs/18/app-initdb.html) consulted2026-09-09;
configuration observations do not substitute for hardware/target durability proof.
No upstream code is copied or relabeled. Runtime stays Go1.26.7, PostgreSQL18.6,
CPython3.14.4 and PowerShell7.6.5, same owned Windows x64 host.

Stage: C:/Users/NL/AppData/Local/Temp/elite-v351-bea959e7a04a42d79829a9c8ffe8135b.
Recompose FRANCHISE_COMPLETE_PACK_PLAN to a new consumer, copy its refund module
to a separate qualification worker, overlay main.go and the Go test below, and
reuse V350 native_stop_windows.go/native_stop_windows_test.go byte-identically.
Build the test binary and preserve its path/hash in build-receipt.json. Reuse the
eight current native Python files from V350's rebuilt directory. Place the SQL
and Python harness below beside the consumer/candidate directories, then run the
harness with an absent output receipt path. Owned cluster directories are unique;
do not point this apparatus at a project, audit, production or existing cluster.

## Qualification source: main.go (AUTHORED / NOT_ADMITTED)

SHA-256: `56f475587cd13eeb17df7e355f0ae42f2014aca04293b8ebc263dad81125ef77`

````go
package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"elite.local/return-refund-worker/internal/refundworker"
	"github.com/jackc/pgx/v5/pgxpool"
)

func claimToken() (string, error) {
	value := make([]byte, 16)
	if _, err := rand.Read(value); err != nil {
		return "", err
	}
	return hex.EncodeToString(value), nil
}

func run() (result error) {
	databaseURL := os.Getenv("DATABASE_URL")
	workerID := os.Getenv("REFUND_WORKER_ID")
	if databaseURL == "" || workerID == "" {
		return errors.New("DATABASE_URL and REFUND_WORKER_ID are required")
	}
	signalCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	ctx, closeNative, err := nativeStopContext(signalCtx, os.Getenv("ELITE_STOP_EVENT_HANDLE"))
	if err != nil {
		return err
	}
	defer func() {
		cause := context.Cause(ctx)
		if cause != nil && cause != errNativeStop && cause != context.Canceled {
			result = errors.Join(result, cause)
		}
		result = errors.Join(result, closeNative())
	}()
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()
	if err = pool.Ping(ctx); err != nil {
		return err
	}
	providers := map[string]refundworker.Provider{}
	if secret := os.Getenv("STRIPE_SECRET_KEY"); secret != "" {
		providers["stripe"], err = refundworker.NewStripeProvider(secret)
		if err != nil {
			return err
		}
	}
	if token := os.Getenv("MERCADO_PAGO_ACCESS_TOKEN"); token != "" {
		providers["mercado_pago"], err = refundworker.NewMercadoPagoProvider(token)
		if err != nil {
			return err
		}
	}
	processor, err := refundworker.NewProcessor(refundworker.NewPostgresStore(pool), providers, workerID, 2*time.Minute, 30*time.Second, 20)
	if err != nil {
		return err
	}
	return runRefundLoop(ctx, processor.Step, claimToken)
}

// The real host stops claiming work when its error report cannot be written.
// The existing durable processor remains responsible for reconciliation. This
// does not retry the uncertain report or change a refund result after commit.
func runRefundLoop(ctx context.Context, step func(context.Context, string) error, nextToken func() (string, error)) error {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		if ctx.Err() != nil {
			return nil
		}
		token, tokenErr := nextToken()
		if tokenErr != nil {
			return tokenErr
		}
		if reportErr := logStepError(step(ctx, token)); reportErr != nil {
			return reportErr
		}
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
		}
	}
}

func main() {
	exitOnHostFailure(run())
}

func exitOnHostFailure(err error) {
	if err != nil {
		log.Fatal("return refund host failed: REFUND_HOST_FAILED")
	}
}

// Log only host-owned operational outcomes. Durable results own diagnostics.
var errRefundReport = errors.New("return refund report unavailable: REFUND_REPORT_FAILED")

type checkedRefundLogWriter struct{ writer io.Writer }

func (w checkedRefundLogWriter) Write(p []byte) (int, error) {
	if w.writer == nil {
		return 0, errRefundReport
	}
	n, err := w.writer.Write(p)
	if n < 0 || n > len(p) {
		return 0, errRefundReport
	}
	if err != nil || n != len(p) {
		return n, errRefundReport
	}
	return n, nil
}

// Preserve the host logger's configured prefix/flags and private fixed message,
// but surface delivery loss. The host loop emits serially. A blocking writer
// still requires target-owned process supervision; no I/O deadline is claimed.
// Full local Write acceptance does not prove log retention or alert delivery.
func logStepError(err error) error {
	if err != nil && !errors.Is(err, refundworker.ErrNoWork) {
		logger := log.New(checkedRefundLogWriter{log.Writer()}, log.Prefix(), log.Flags())
		if logger.Output(2, "return refund step failed: REFUND_STEP_FAILED") != nil {
			return errRefundReport
		}
	}
	return nil
}
````

## Qualification source: native_postgres_fixture_windows_test.go (AUTHORED / NOT_ADMITTED)

SHA-256: `6bbcba3be7247f0fd55cafaf653c996d9587866c641af4df897311ed273f0d66`

````go
//go:build windows

package main

// AUTHORED qualification, NOT_ADMITTED. Real PostgreSQL Store and host; synthetic
// provider effects stay in a dedicated local test schema. No payment transport.
import (
	"context"
	"crypto/sha256"
	"elite.local/return-refund-worker/internal/refundworker"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"io"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"
)

var errFixtureConfig = errors.New("FIXTURE_CONFIG_REJECTED")
var fixtureDatabase = regexp.MustCompile(`^elite_refund_test_[0-9a-f]{32}$`)
var fixtureDSN = regexp.MustCompile(`^postgres://postgres@127\.0\.0\.1:[0-9]{1,5}/elite_refund_test_[0-9a-f]{32}\?sslmode=disable$`)

func fixtureConfig(raw []byte) (*pgxpool.Config, error) {
	if len(raw) == 0 || len(raw) > 2048 || !fixtureDSN.Match(raw) {
		return nil, errFixtureConfig
	}
	c, e := pgxpool.ParseConfig(string(raw))
	if e != nil {
		return nil, errFixtureConfig
	}
	if c.ConnConfig.Host != "127.0.0.1" || c.ConnConfig.Port == 0 || !fixtureDatabase.MatchString(c.ConnConfig.Database) || c.ConnConfig.User != "postgres" || c.ConnConfig.Password != "" || c.ConnConfig.TLSConfig != nil || len(c.ConnConfig.Fallbacks) != 0 {
		return nil, errFixtureConfig
	}
	if len(c.ConnConfig.RuntimeParams) != 0 {
		return nil, errFixtureConfig
	}
	c.ConnConfig.ConnectTimeout = 3 * time.Second
	return c, nil
}
func readFixtureConfig(path, digest string) (string, *pgxpool.Config, error) {
	info, e := os.Lstat(path)
	if e != nil || !info.Mode().IsRegular() || info.Size() > 2048 {
		return "", nil, errFixtureConfig
	}
	f, e := os.Open(path)
	if e != nil {
		return "", nil, errFixtureConfig
	}
	defer f.Close()
	raw, e := io.ReadAll(io.LimitReader(f, 2049))
	if e != nil {
		return "", nil, errFixtureConfig
	}
	hash := sha256.Sum256(raw)
	if hex.EncodeToString(hash[:]) != digest {
		return "", nil, errFixtureConfig
	}
	cfg, e := fixtureConfig(raw)
	if e != nil {
		return "", nil, e
	}
	return string(raw), cfg, nil
}
func pgFixtureWrite(name string, value any) error {
	raw, e := json.Marshal(value)
	if e != nil {
		return e
	}
	f, e := os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if e != nil {
		return e
	}
	n, e := f.Write(raw)
	if e == nil && n != len(raw) {
		e = io.ErrShortWrite
	}
	if e == nil {
		e = f.Sync()
	}
	return errors.Join(e, f.Close())
}

type pgFixtureProvider struct {
	pool      *pgxpool.Pool
	mode, run string
}

func (p pgFixtureProvider) Create(ctx context.Context, r refundworker.Refund) (refundworker.ProviderResult, error) {
	effectCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if _, e := p.pool.Exec(effectCtx, `insert into qualification.effect(run_id,idempotency_key) values($1,$2)`, p.run, r.IdempotencyKey); e != nil {
		return refundworker.ProviderResult{}, e
	}
	if p.mode == "effect-cancel" || p.mode == "ignore-stop" {
		if e := pgFixtureWrite("ready", map[string]string{"stage": "effect-committed"}); e != nil {
			return refundworker.ProviderResult{}, e
		}
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
	return refundworker.ProviderResult{ProviderPaymentReference: r.ProviderPaymentReference, ProviderRefundReference: "fixture-" + p.run, ProviderStatus: status, Currency: r.Currency, AmountMinorUnits: r.AmountMinorUnits}, nil
}
func (p pgFixtureProvider) Retrieve(ctx context.Context, r refundworker.Refund) (refundworker.ProviderResult, error) {
	return p.Create(ctx, r)
}

type pgFixtureStore struct {
	*refundworker.PostgresStore
	mode string
}

func (s pgFixtureStore) Complete(ctx context.Context, w refundworker.Work, id string, r refundworker.Refund, p refundworker.ProviderResult, outcome, code string, retry time.Duration) error {
	e := s.PostgresStore.Complete(ctx, w, id, r, p, outcome, code, retry)
	if markErr := pgFixtureWrite("ready", map[string]any{"stage": "complete-returned", "error": e != nil}); markErr != nil {
		return errors.Join(e, markErr)
	}
	<-ctx.Done()
	if s.mode == "lost-ack" {
		return errors.New("PRIVATE_FIXTURE_LOST_ACK")
	}
	return e
}
func runPGFixture(mode, run, path, digest string) error {
	raw, cfg, e := readFixtureConfig(path, digest)
	if e != nil {
		return e
	}
	if e = pgFixtureWrite("config-accepted", map[string]string{"run_id": run}); e != nil {
		return e
	}
	if mode == "startup" {
		if e = os.Setenv("DATABASE_URL", raw); e != nil {
			return e
		}
		if e = os.Setenv("REFUND_WORKER_ID", "fixture-"+run); e != nil {
			return e
		}
		return runHostFixtureStartup()
	}
	ctx, closeStop, e := nativeStopContext(context.Background(), os.Getenv("ELITE_STOP_EVENT_HANDLE"))
	if e != nil {
		return e
	}
	defer closeStop()
	pool, e := pgxpool.NewWithConfig(ctx, cfg)
	if e != nil {
		return e
	}
	defer pool.Close()
	if e = pool.Ping(ctx); e != nil {
		return e
	}
	store := pgFixtureStore{refundworker.NewPostgresStore(pool), mode}
	processor, e := refundworker.NewProcessor(store, map[string]refundworker.Provider{"stripe": pgFixtureProvider{pool, mode, run}}, "fixture-"+run, 2*time.Minute, time.Hour, 20)
	if e != nil {
		return e
	}
	e = runRefundLoop(ctx, processor.Step, claimToken)
	cause := context.Cause(ctx)
	if cause != errNativeStop && cause != context.Canceled {
		e = errors.Join(e, errNativeProtocol)
	}
	return errors.Join(e, closeStop())
}

// The startup lane executes run() itself: actual configuration, Ping, provider
// absence, PostgresStore construction and host loop, with only native stop wiring
// overlaid in main.go. Credentials remain absent and the domain must be blocked.
func runHostFixtureStartup() error { return run() }
func TestMain(m *testing.M) {
	if len(os.Args) >= 5 && os.Args[len(os.Args)-5] == "--pg-fixture" {
		args := os.Args[len(os.Args)-4:]
		e := runPGFixture(args[0], args[1], args[2], args[3])
		if markErr := pgFixtureWrite("lifecycle-ack.json", map[string]any{"run_id": args[1], "returned": true, "failed": e != nil}); markErr != nil {
			e = errors.Join(e, markErr)
		}
		exitOnHostFailure(e)
		os.Exit(0)
	}
	os.Exit(m.Run())
}

func TestPGFixtureConfigGuard(t *testing.T) {
	good := "postgres://postgres@127.0.0.1:5432/elite_refund_test_0123456789abcdef0123456789abcdef?sslmode=disable"
	if _, e := fixtureConfig([]byte(good)); e != nil {
		t.Fatal(e)
	}
	for _, bad := range []string{"", good + "\n", strings.Replace(good, "127.0.0.1", "localhost", 1), strings.Replace(good, "127.0.0.1", "example.test", 1), strings.Replace(good, "elite_refund_test_0123456789abcdef0123456789abcdef", "production", 1), strings.Replace(good, "postgres@", "postgres:secret@", 1), good + "&application_name=private", strings.Replace(good, "sslmode=disable", "sslmode=prefer", 1)} {
		if _, e := fixtureConfig([]byte(bad)); e != errFixtureConfig {
			t.Fatal("unsafe config accepted")
		}
	}
}
func TestPGFixtureConfigDigestMismatch(t *testing.T) {
	p := t.TempDir() + "/config"
	if e := os.WriteFile(p, []byte("private"), 0600); e != nil {
		t.Fatal(e)
	}
	if _, _, e := readFixtureConfig(p, strings.Repeat("0", 64)); e != errFixtureConfig {
		t.Fatal("digest accepted")
	}
}

func FuzzPGFixtureConfig(f *testing.F) {
	for _, seed := range []string{"", "postgres://postgres@127.0.0.1:5432/elite_refund_test_0123456789abcdef0123456789abcdef?sslmode=disable", "postgres://postgres@localhost/production", "postgres://user:secret@example.test/db", "\x00", "service=private", "passfile=private"} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, raw string) {
		c, e := fixtureConfig([]byte(raw))
		if e != nil {
			if e != errFixtureConfig || c != nil {
				t.Fatal("nonstatic rejection")
			}
			return
		}
		if !fixtureDSN.MatchString(raw) || c.ConnConfig.Host != "127.0.0.1" || c.ConnConfig.Port == 0 || !fixtureDatabase.MatchString(c.ConnConfig.Database) || c.ConnConfig.Password != "" || c.ConnConfig.TLSConfig != nil || len(c.ConnConfig.Fallbacks) != 0 || len(c.ConnConfig.RuntimeParams) != 0 {
			t.Fatal("unsafe fixture accepted")
		}
	})
}
````

## Qualification source: seed.sql (AUTHORED / NOT_ADMITTED)

SHA-256: `6fa405b703985ccbde52f1d3f37f130826d10d0056afbcccd667b0a557792b9a`

````sql
-- AUTHORED local fixture. Pre-existing captured order / authorized return are
-- synthetic starting conditions, NOT an originating sales/approval workflow.
-- All actual constraints/triggers remain enabled. Fresh database only.
begin;
insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values('018f4d4a-7b36-7a21-8d10-2f4c54c2f351','fixture351','Synthetic','Synthetic');
insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type)values('018f4d4a-7b36-7a21-8d10-2f4c54c2f351','store','store','Synthetic','store');
insert into catalog.vehicle_model(tenant_id,model_id,model_code,display_name,vehicle_class,lifecycle_state)values('018f4d4a-7b36-7a21-8d10-2f4c54c2f351','model','model','Synthetic','bicycle','active');
insert into catalog.vehicle_variant(tenant_id,variant_id,model_id,variant_code,display_name,battery_specification,lifecycle_state)values('018f4d4a-7b36-7a21-8d10-2f4c54c2f351','variant','model','variant','Synthetic','{}','active');
insert into crm.customer_profile(tenant_id,customer_principal_id,display_name,email_normalized)values('018f4d4a-7b36-7a21-8d10-2f4c54c2f351','customer','Synthetic','synthetic@example.test');
insert into inventory.stock_unit(tenant_id,stock_unit_id,organization_id,variant_id,serial_number,state,version,received_at)values('018f4d4a-7b36-7a21-8d10-2f4c54c2f351','stock','store','variant','SYNTHETIC-SERIAL','sold',1,clock_timestamp());
insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version)values('018f4d4a-7b36-7a21-8d10-2f4c54c2f351','order','store','customer','delivered','ARS',1000,1);
insert into sales.customer_order_line(tenant_id,order_id,line_id,variant_id,quantity,unit_price_minor_units,allocated_stock_unit_id)values('018f4d4a-7b36-7a21-8d10-2f4c54c2f351','order','line','variant',1,1000,'stock');
insert into payment.payment_attempt(tenant_id,payment_attempt_id,order_id,provider_code,provider_reference,idempotency_key,state,currency,amount_minor_units,version)values('018f4d4a-7b36-7a21-8d10-2f4c54c2f351','payment','order','stripe','synthetic-pi','fixture-payment-00000001','captured','ARS',1000,3);
insert into sales.delivery_handover(tenant_id,handover_id,organization_id,order_id,customer_principal_id,stock_unit_id,state,version)values('018f4d4a-7b36-7a21-8d10-2f4c54c2f351','handover','store','order','customer','stock','prepared',1);
insert into sales.delivery_exception(tenant_id,exception_id,organization_id,handover_id,customer_principal_id,reason_code,details,rejection_evidence_sha256_hex,state,version)values('018f4d4a-7b36-7a21-8d10-2f4c54c2f351','exception','store','handover','customer','synthetic-return','Synthetic',repeat('a',64),'open',1);
insert into sales.return_authorization(tenant_id,authorization_id,organization_id,exception_id,handover_id,order_id,stock_unit_id,customer_principal_id,disposition,exact_cost_source_order_id,state,authorized_by_subject)values('018f4d4a-7b36-7a21-8d10-2f4c54c2f351','authorization','store','exception','handover','order','stock','customer','return','order','authorized','synthetic-operator');
insert into sales.return_receipt(tenant_id,receipt_id,authorization_id,organization_id,order_id,stock_unit_id,customer_principal_id,received_serial_number,condition_code,notes,evidence_sha256_hex,received_by_subject)values('018f4d4a-7b36-7a21-8d10-2f4c54c2f351','receipt','authorization','store','order','stock','customer','SYNTHETIC-SERIAL','opened','Synthetic',repeat('a',64),'synthetic-operator');
insert into sales.return_disposition(tenant_id,disposition_id,receipt_id,inventory_action,customer_remedy,notes,decided_by_subject)values('018f4d4a-7b36-7a21-8d10-2f4c54c2f351','disposition','receipt','restock','refund','Synthetic','synthetic-operator');
insert into sales.return_effect_request(tenant_id,request_id,disposition_id,effect_kind,owner_context,state,idempotency_key)values('018f4d4a-7b36-7a21-8d10-2f4c54c2f351','refund','disposition','refund','payment','requested','synthetic-refund-0000000001');
create schema qualification;
create table qualification.effect(run_id text primary key,idempotency_key text not null unique);
commit;
````

## Qualification source: run_postgres_lifecycle.py (AUTHORED / NOT_ADMITTED)

SHA-256: `7c02a0d10f17d49480bf810edf1edc210e60326cbd9cad69b0893f250356beba`

````python
"""AUTHORED qualification: isolated PostgreSQL + real refund host/store.
Synthetic captured-order/authorized-return fixtures and provider effects only.
No existing cluster/database, credentials, payment transport or release deploy.
"""
from pathlib import Path
from contextlib import nullcontext
from unittest.mock import patch
import hashlib,json,os,re,socket,subprocess,sys,threading,time,traceback,uuid
S=Path(__file__).resolve().parent
sys.path.insert(0,str(S/'candidate'))
import result_store as rs
PG=Path('C:/Users/NL/AppData/Local/Temp/elite-franchise-build-f3200751e6f644c9ae8c08a62d6b14e5/postgresql-18.6/pgsql/bin')
BUILD=json.loads((S/'build-receipt.json').read_text());sha=lambda p:hashlib.sha256(p.read_bytes()).hexdigest()
WORK=S/('pg-run-'+uuid.uuid4().hex);WORK.mkdir();DATA=WORK/'pgdata'
ENV={k:v for k,v in os.environ.items() if not k.startswith('PG')};ENV['PGCONNECT_TIMEOUT']='3'
CHILD=subprocess.CREATE_NO_WINDOW
with socket.socket() as port_probe:port_probe.bind(('127.0.0.1',0));PORT=port_probe.getsockname()[1]
DB=re.compile(r'^elite_refund_test_[0-9a-f]{32}$')
class Fixture:
 def __init__(self):self.proc=None;self.logs=[];self.dbnames=[];self.cases=[];self.fail=None
 def command(self,args,log=None,timeout=60):
  p=subprocess.run([str(x) for x in args],env=ENV,stdin=subprocess.DEVNULL,stdout=subprocess.PIPE,stderr=subprocess.PIPE,creationflags=CHILD,timeout=timeout)
  if log:(WORK/log).write_bytes(p.stdout+p.stderr)
  if p.returncode:raise RuntimeError(f'owned command failed {Path(str(args[0])).name}: '+p.stderr.decode('utf-8',errors='replace')[:2000])
  return p.stdout.decode('utf-8').strip()
 def sql(self,db,sql):
  assert db=='postgres' or DB.fullmatch(db)
  return self.command([PG/'psql.exe','-X','-w','-qAt','-h','127.0.0.1','-p',PORT,'-U','postgres','-d',db,'-v','ON_ERROR_STOP=1','-c',sql],timeout=12)
 def file(self,db,path):
  return self.command([PG/'psql.exe','-X','-w','-qAt','-h','127.0.0.1','-p',PORT,'-U','postgres','-d',db,'-v','ON_ERROR_STOP=1','-f',path],timeout=30)
 def start(self):
  assert DATA.resolve().is_relative_to(WORK.resolve())
  out=(WORK/f'postgres-{len(self.logs)}.log').open('xb');self.logs.append(out)
  self.proc=subprocess.Popen([str(PG/'postgres.exe'),'-D',str(DATA)],env=ENV,stdin=subprocess.DEVNULL,stdout=out,stderr=subprocess.STDOUT,creationflags=CHILD)
  end=time.monotonic()+15
  while time.monotonic()<end:
   if self.proc.poll() is not None:raise RuntimeError('owned postgres exited during startup')
   p=subprocess.run([str(PG/'pg_isready.exe'),'-q','-h','127.0.0.1','-p',str(PORT),'-U','postgres','-d','postgres'],env=ENV,stdin=subprocess.DEVNULL,stdout=subprocess.DEVNULL,stderr=subprocess.DEVNULL,creationflags=CHILD,timeout=3)
   if p.returncode==0:return
   time.sleep(.05)
  raise TimeoutError('owned postgres readiness')
 def stop(self):
  if self.proc is not None and self.proc.poll() is None:
   assert DATA.resolve().is_relative_to(WORK.resolve()) and (DATA/'PG_VERSION').read_text().strip()=='18'
   # Observed shutdown checkpoint needed14.885s to flush12841 files. The old10s
   # test budget was insufficient. Keep a finite40s bound and preserve expiry;
   # this is a fixture budget, not an admitted production shutdown SLO.
   self.command([PG/'pg_ctl.exe','-D',DATA,'-m','fast','-w','-t','40','stop'],log=f'stop-{len(self.logs)}.log',timeout=45)
   assert self.proc.wait(timeout=10)==0
  for f in self.logs:
   if not f.closed:f.close()
 def initialize(self):
  version=self.command([PG/'postgres.exe','--version']);assert version=='postgres (PostgreSQL) 18.6'
  self.command([PG/'initdb.exe','-D',DATA,'-U','postgres','-A','trust','--no-locale','-E','UTF8'],log='initdb.log')
  # Preserve initdb's validated zone. This installed runtime lacks the UTC alias;
  # observe the actual timezone instead of inferring complete timezone coverage.
  with (DATA/'postgresql.conf').open('a') as f:f.write(f"\nlisten_addresses='127.0.0.1'\nport={PORT}\nfsync=on\nsynchronous_commit=on\nfull_page_writes=on\n")
  self.start();self.template='elite_refund_test_'+uuid.uuid4().hex;self.sql('postgres',f'create database {self.template}')
  migrations=sorted((S/'consumer/db/migrations').glob('*.up.sql'));assert migrations
  with (WORK/'migrations.log').open('w',encoding='utf-8') as log:
   for path in migrations:log.write(path.name+'\n'+self.file(self.template,path)+'\n');log.flush()
  self.file(self.template,S/'seed.sql')
  settings=json.loads(self.sql(self.template,"select json_build_object('fsync',current_setting('fsync'),'synchronous_commit',current_setting('synchronous_commit'),'full_page_writes',current_setting('full_page_writes'),'replication_role',current_setting('session_replication_role'),'checksums',current_setting('data_checksums'))"))
  assert settings=={'fsync':'on','synchronous_commit':'on','full_page_writes':'on','replication_role':'origin','checksums':'on'}
  assert self.sql(self.template,"select count(*) from pg_trigger where not tgisinternal and tgenabled<>'O'")=='0'
  settings['timezone']=self.sql(self.template,'show timezone')
  self.runtime={'version':version,'executables':{n:sha(PG/n) for n in ['postgres.exe','initdb.exe','pg_ctl.exe','psql.exe']},'settings':settings,'migrations':{p.name:sha(p) for p in migrations},'seed_sha256':sha(S/'seed.sql'),'utc_alias_present':(PG.parent/'share/timezone/UTC').exists()}
 def snapshot(self,db):
  names=['payment.payment_attempt','payment.return_refund','payment.return_refund_observation','sales.return_effect_execution','sales.return_effect_attempt','platform.outbox_event','qualification.effect']
  return {n:json.loads(self.sql(db,f"select coalesce(jsonb_agg(to_jsonb(t) order by to_jsonb(t)::text),'[]'::jsonb) from {n} t")) for n in names}
 def case(self,mode,publication_loss=False,tamper=False,precancel=False):
  root=WORK/(mode+'-'+uuid.uuid4().hex);root.mkdir();store=root/'store';store.mkdir();run=uuid.uuid4().hex
  db='elite_refund_test_'+uuid.uuid4().hex;self.sql('postgres',f'create database {db} template {self.template}');self.dbnames.append(db)
  if mode=='commit-failure':
   self.sql(db,"create function qualification.reject_observation() returns trigger language plpgsql as $$begin raise exception 'PRIVATE_FIXTURE_COMMIT_FAILURE'; end$$; create trigger fixture_reject before insert on payment.return_refund_observation for each row execute function qualification.reject_observation()")
  cfg=root/'fixture.dsn';cfg.write_bytes(f'postgres://postgres@127.0.0.1:{PORT}/{db}?sslmode=disable'.encode());cfg_sha=sha(cfg)
  if tamper:cfg.write_bytes(cfg.read_bytes()+b'\n')
  exe=Path(BUILD['fixture']['path']);assert sha(exe)==BUILD['fixture']['sha256']
  p={'schema':'elite-native-launch-qualification/v2','id':'postgres-native-fixture','executable':{'path':str(exe),'bytes':exe.stat().st_size,'sha256':sha(exe)},'arguments':['-test.run=^$','--','--pg-fixture',mode,run,str(cfg),cfg_sha],'cwd':str(root),'environment':{'SYSTEMROOT':os.environ['SYSTEMROOT'],'WINDIR':os.environ['WINDIR'],'TEMP':str(root),'TMP':str(root)},'budgets':{'timeout':15,'output_bytes':4096,'processes':8,'commit_bytes':536870912},'shutdown':{'protocol':'win32-inherited-event/v1','grace_seconds':.15 if mode=='ignore-stop' else 2}}
  raw=rs.encode(p);digest=rs.digest(raw);cancel=threading.Event();ready=[];errors=[]
  def watcher():
   try:
    end=time.monotonic()+12
    while time.monotonic()<end:
     signal=(self.sql(db,"select status from sales.return_effect_execution where request_id='refund'")=='blocked') if mode=='startup' else (root/'ready').exists()
     if signal:ready.append(True);break
     time.sleep(.01)
   except Exception as e:errors.append(type(e).__name__)
   finally:cancel.set()
  if precancel:cancel.set()
  th=None
  if not tamper and not precancel:th=threading.Thread(target=watcher);th.start()
  original=rs.publish_file
  def lost(directory,name,data):
   if name=='result.json':raise OSError('fixture publication unavailable')
   return original(directory,name,data)
  with patch.object(rs,'publish_file',side_effect=lost) if publication_loss else nullcontext():
   receipt=rs.execute(store,run,raw,digest,cancel=cancel)
  if th:th.join(14);assert not th.is_alive();assert not errors,errors;assert ready==[True],root
  reopened=rs.inspect(store,run,digest)
  if publication_loss:assert receipt.state=='RESULT_UNCERTAIN' and reopened.state=='UNRESOLVED';out=None
  else:
   assert receipt.state=='RECORDED' and receipt==reopened,(root,receipt,reopened)
   out=json.loads((store/run/'result.json').read_text())['outcome']
  before=self.snapshot(db)
  with patch.object(rs.g,'launch') as launch:assert rs.execute(store,run,raw,digest).state=='ALREADY_RESERVED';launch.assert_not_called()
  after=self.snapshot(db);assert after==before,'replay mutated domain'
  execution=after['sales.return_effect_execution'];assert len(execution)==1
  status=execution[0]['status'];attempts=execution[0]['attempt_count'];observations=len(after['payment.return_refund_observation']);effects=len(after['qualification.effect']);outbox=len(after['platform.outbox_event'])
  if tamper:
   assert status=='requested' and attempts==0 and effects==0 and out['exit_code']==1;assert not (root/'config-accepted').exists()
  elif precancel:
   assert status=='requested' and attempts==0 and effects==0;assert not (root/'config-accepted').exists()
  else:
   assert attempts==1
   if out:assert out['status']=='CANCELLED' and out['tree_empty'];assert out['shutdown_state']==('FORCED' if mode=='ignore-stop' else 'EXITED_DURING_GRACE')
   if mode!='ignore-stop':assert json.loads((root/'lifecycle-ack.json').read_text())=={'run_id':run,'returned':True,'failed':False}
   if mode=='startup':assert status=='blocked' and execution[0]['last_error_code']=='PROVIDER_CONFIG_MISSING' and effects==0 and observations==0 and outbox==0
   elif mode in ('succeeded','lost-ack'):
    assert status=='succeeded' and observations==1 and effects==1 and outbox==2
    assert after['payment.payment_attempt'][0]['state']=='refunded' and after['payment.return_refund'][0]['state']=='succeeded'
   elif mode in ('pending','failed'):
    assert status==('retry' if mode=='pending' else 'failed') and observations==1 and effects==1 and outbox==0
    assert after['payment.payment_attempt'][0]['state']=='captured'
   elif mode in ('effect-cancel','ignore-stop','commit-failure'):
    assert status=='claimed' and observations==0 and effects==1 and outbox==0
    assert after['payment.return_refund'][0]['state']=='prepared' and after['payment.payment_attempt'][0]['state']=='captured'
   else:raise AssertionError('unknown fixture mode')
  receipt_info={'mode':mode,'publication_loss':publication_loss,'tampered_config':tamper,'precancel':precancel,'database':db,'root':str(root),'process_result_state':receipt.state,'reopened':reopened.state,'outcome':out,'domain_status':status,'attempts':attempts,'observations':observations,'synthetic_effects':effects,'outbox':outbox,'snapshot':after}
  (root/'asserted-evidence.json').write_text(json.dumps(receipt_info,indent=2)+'\n');self.cases.append(receipt_info);print('PG_NATIVE_CASE_PASS',mode,('publication-loss' if publication_loss else 'tamper' if tamper else 'precancel' if precancel else ''),status,flush=True)
 def execute(self):
  try:
   self.initialize()
   for mode in ['startup','succeeded','pending','failed','lost-ack','effect-cancel','commit-failure','ignore-stop']:self.case(mode)
   self.case('succeeded',publication_loss=True);self.case('startup',tamper=True);self.case('startup',precancel=True)
   before={db:self.snapshot(db) for db in self.dbnames};self.stop();self.start();after={db:self.snapshot(db) for db in self.dbnames};assert before==after
   # Full test snapshots are evidence, not16KiB operational result records.
   snapshot_digest=lambda value:hashlib.sha256(json.dumps(value,sort_keys=True,separators=(',',':')).encode('utf-8')).hexdigest()
   self.restart={'databases':len(self.dbnames),'before_sha256':snapshot_digest(before),'after_sha256':snapshot_digest(after),'identical':True,'scope':'controlled local PostgreSQL restart, not PITR/power-cycle/DR'}
   print('PG_NATIVE_RESTART_PASS databases',len(self.dbnames),flush=True)
  except BaseException as e:self.fail={'type':type(e).__name__,'detail':str(e)};traceback.print_exc()
  finally:
   try:self.stop()
   except Exception as e:self.fail={'type':type(e).__name__,'detail':'owned postgres stop failed: '+str(e)};traceback.print_exc()
   result={'pass':self.fail is None,'cases_passed':len(self.cases),'cases':self.cases,'error':self.fail,'runtime':getattr(self,'runtime',None),'restart':getattr(self,'restart',None),'work':str(WORK),'owned_postgres_stopped':self.proc is None or self.proc.poll() is not None}
   Path(sys.argv[1]).write_text(json.dumps(result,indent=2)+'\n')
  return self.fail is None
if __name__=='__main__':raise SystemExit(0 if Fixture().execute() else 1)
````

## Canonical integration146 — 2026-09-09T14:57:03.930034Z

Fresh canonical reconstruction PASS:746product files equal, four qualification
source blocks, eight reused native Python files and two reused Go native files exact.
Go baseline passes three repetitions, including10native-stop and2configuration
tests; vet, test-binary and host builds PASS. Distinct opt-in skipped tests:
TestPostgresRefundLineAllocationPartialThenFullAndAmbiguity. These are not executed evidence.

All11native PostgreSQL lifecycle cases PASS and all11database snapshots are
byte-equivalent under canonical JSON after a controlled PostgreSQL restart.
Owned server stopped successfully. This demonstrates local transaction/restart
behavior with fsync/synchronous_commit/full_page_writes enabled and constraints
active. It does not prove hardware power-loss, PITR, external provider effect or
originating business approval: provider effects and seed conditions are synthetic.
The normal run() startup is invoked through a qualification TestMain with local
configuration; it is not an installed production service or production secret gate.

Retained supervisor90tests PASS, zero skips. Fuzz configuration parser:7seeds,
3687878executions, finite10s PASS. Fuzz kit1positive/2negative and capability gap
kit3positive/6negative PASS. FAIL608–612 recovered with original evidence retained;
ledger2537unique IDs. Full Preflight146 is pending.

Registered contract progress remains41/48passed,5blocked,2planned:7/48=14.5833%
not passed, not percentage of remaining effort. Roadmap1complete/10open. RoundD
historical corpus channel/format question remains unanswered. No invented approval.

| Canonical evidence | SHA-256 |
|---|---|
| postgres-first.json | `d534564bb755c57f1b75eb32b8c29b0b8f872cb0a01281a07156fb9f379cd14e` |
| postgres-second.json | `b0216fcf43ceff0c2f0092fc3b04c7bc4100d74e8cfe8ae82cb7f3faf84942d8` |
| postgres-third.json | `88272960780178be8e9a9569384ea5bc4e9f832481106515975f54e1313fae01` |
| postgres-rebuilt.json | `364369b5da63ea5e4bb072c82c23f38947503ed319cee0f22afc0fe545948ef3` |
| postgres-rebuilt.log | `ceb4973f2c9ded68cd16874e4c3d2ee3a04edf73a2f37957ce645c5811953fff` |
| canonical-parity.json | `2c3b2f6aead307f45f3c625a51cd8c287c531b017ba95493f0d88b2ec58499ef` |
| canonical/build-receipt.json | `521ad9ca275acc3c6db9ff1fa9e2ab3f701317a1b7961da26052191ee48f2a19` |
| rebuilt-go-tests.log | `a3f73166fbcf38d8093ae66da15127b5e565bd4c903ff34037c2b6b20a7faf3a` |
| rebuilt-go-vet.log | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| rebuilt-result.json | `a0d6088680bcaca86b4f5fe5901e8d704ca69d3caef80249bdafd68861e7537e` |
| rebuilt-shutdown.json | `6b0586d112235212d1a8b31b663dd99b84706630eb4de97ea6e80463f6e4c549` |
| rebuilt-launch.json | `404a8164e7aff98a87c13aec6d47652d4c0a67473a0ed4273b4083eebf96ce61` |
| rebuilt-capture.json | `c4b869ed05b0f2082f785da240aab9346736b414b47724192311db00336fa798` |
| test-summary.json | `4616be41bce98d37181d4b22de9b5a272ec50a9792b218a7a0edb312b2abb480` |
| fuzz-profile.json | `a7d6a2b2ef74573ec8cbf843fc0c9c8a1aab85c16ed2c2f7ebb10ceb017e4669` |
| fuzz-receipt.json | `a262731a2fd5685679c9e5d3091409fdfe0beaa6289835789697a14306686103` |
| fuzz.log | `bb3e630bd3a64b7e960893a4e95624caa8d808d1919f2af2408dcd4144df7128` |
| fuzz-self.log | `716c481531ec334c8b120b6a4128753f395d20c68c6f8aecac1f5dfa084707ac` |
| gap-self.log | `dbe711c12c2333308f23c27838212995c31b16ae7bdeac14051225577a79dcec` |
## Closure147 — local PostgreSQL lifecycle qualification complete

Observed 2026-09-09T15:01:11.058136Z. Full Preflight146: 154 executed checks PASS;
162packs/1461materializable files/789Markdown/53profiles. Docker unavailable.
This supersedes the pending Preflight146 snapshot above. Conditional network,
provider, opt-in PostgreSQL and target gates are not promoted by library Preflight.

Exact reconstruction preserves746product files; four qualification sources and
all reused native sources match their declared digests.11PostgreSQL cases and
controlled restart of11databases PASS; before/after digest:
`cfd873afadd3ede65b15c3a9f2b639388609a216e7c0550255e5c27716a6055d`. Owned PostgreSQL was stopped successfully.
90retained native tests, Go baseline x3/vet/build and configuration fuzz3687878PASS.
One legacy opt-in PostgreSQL test stays skipped; the separate eleven-case fixture
uses real transactions with constraints/triggers enabled and synthetic provider.

The concrete advance is startup run() through qualification TestMain → native
lifecycle → real Processor/PostgresStore → separate process receipt and durable
domain observations → no-replay readback → controlled PostgreSQL restart.
Installed service identity, production configuration/secrets, collector/retention/
alerts, integrated security and target admission remain open. The existing product
packs stay unchanged and no real payment/provider/account or final release ran.

Contract still41passed/5blocked/2planned of48; no exact effort or total-completion
percentage follows. Next T2809 boundary is existing service identity/configuration
and governed operational collection, then whole-gate security/retention/alerts.
Historical channel/format answer remains pending; preserve raw history and derived
lineage. Final successor gap-record147.json/gap-receipt147.json remains conditioned.

| Closure evidence | SHA-256 |
|---|---|
| preflight146.json | `a34843f272402553625990c5b93c244adef9436dcb296cd52f5dd5d5cf42c12e` |
| preflight146.log | `c2aeda252ae56776e818587d71058ac4b60396090192983884f7e8b6cb091906` |
| structural-closure146.json | `721e299c7b1460111938bb3c5cf4b9013fc8cbf23a21680d36ff79f4b58eb4ca` |
| gap-record146.json | `e197bb67dd9c75fa263c7c741641ed4f4ff9c0a01b5414a136855b7a4ece47a4` |
| gap-receipt146.json | `3536e907c16bf1a61ee13db51dd83626cb5124f3726e760c080fa8d022ff7ad2` |
| resume146.log | `86d3390412a2a3886678dea24a0b976505496b1d96fe639c7539236cfd2a1042` |
| plan146.log | `2fd65a7ad27c8b968c64f8079abf8610894d6e7732b0b0f514dba962cf93b565` |
