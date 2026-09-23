# Windows reference telemetry runtime

This is an opt-in, finite reference acceptance runner for the library. It mounts
the actual refund main and CLI, official-source minimal Collector and Prometheus,
and a fresh isolated PostgreSQL cluster. It is not a production service installer,
secret manager, payment environment, notification recipient or training pipeline.
Use only on a trusted Windows x64 host under the current local account.

Materialize GO-OFFICIAL-RETURN-REFUND-WORKER0.1.5 as part of the compatible full
consumer composition. Build its two cmd binaries with the admitted Go1.26.8.
Acquire and qualify the exact Collector0.160.0 and Prometheus3.14.0 sources through
the maintenance-telemetry-source-inspection profile. The V372 evidence preserves
the exact OCB selection, generated source, module locks, build commands, security
patches, licenses, tests and artifact hashes; official release binaries are not
interchangeable with this minimal locally built distribution. Verify the local
PostgreSQL18.6 runtime independently. Do not fill the runtime lock from reputation.

Copy runtime-lock.example.json outside the source tree. Record each verified
binary path, byte size and SHA-256, every selected *.up.sql filename/hash from the
consumer, and the exact seed.sql hash. Bind that entire JSON by an independently
trusted SHA-256 at invocation. The example is intentionally rejected unchanged.
No runtime downloads, installations or external accounts occur in this runner.

Run:

```text
python -B reference_telemetry/run_reference_runtime.py --runtime-lock LOCK.json --runtime-lock-sha256 TRUSTED_SHA256 --consumer MATERIALIZED_CONSUMER --work-root OWNED_EVIDENCE_DIRECTORY
```

The evidence root must already exist and be short enough for the native path
protocol. A new UUID directory receives a current-user/SYSTEM ACL, ephemeral
certificates, exact configs, a fresh cluster and bounded process receipts. The
CA private key is never saved. Client/server keys expire after two hours and
must never be published. The cluster identity is checked before any mutation.
All owned processes stop in finally; files and failed evidence are retained.

Five host attempts verify bad-config rejection before claims, actual missing
provider blocking, DB-error signals, graceful stop, collector loss and recovery.
A handled step does not imply a successful refund. No provider key is inherited.
OTLP, metrics and the Prometheus API require mutual TLS; rogue/missing clients
fail. The separate fixed health response is bound to local HTTP.
Actual OTLP output proves log/span correlation and cumulative metric invariants.
Prometheus executes the error and collector-health rules and their recovery.
A2000-observation/8-instance CLI probe drives1MiB rotation,2backup retention and
expiry; retained files remain byte-identical after collector process restart.
This is local file/process recovery, not power-loss durability, PITR or DR.
Prometheus64MB/1h retention is a configured TSDB policy, not a hard total-disk cap.

Regression suites:

```text
python -B reference_telemetry/test_reference_profile.py policy.json
python -B reference_telemetry/test_capture.py capture.json REFUND_EXE REFUND_SHA256
python -B reference_telemetry/test_governed_launch.py launch.json REFUND_FIXTURE_JSON
python -B reference_telemetry/test_shutdown.py shutdown.json
python -B reference_telemetry/test_result_store.py result.json REFUND_FIXTURE_JSON
```

REFUND_FIXTURE_JSON contains path and sha256 for that exact executable. The
supervisor is a trusted-launcher primitive with finite capture, job limits,
cooperative stop and durable local metadata. It never retries an uncertain run.
It does not sandbox administrators or hostile same-user processes. Runtime,
service identity, secret custody, retention obligations, alert routing, workload,
security and acceptance must be requalified in every real consumer target.
