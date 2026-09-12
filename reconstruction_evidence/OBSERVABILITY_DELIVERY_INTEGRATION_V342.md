# V342 — explicit delivery failure and local telemetry integration

Date:2026-09-08. Corrective library maintenance; existing candidate owner
GO-OBSERVABILITY-CORE0.3.0 →0.3.1. Two AUTHORED files, unchanged file count.

## Defect and repair

Four tests against exact0.3.0 source reproduced false successful delivery when a
configured writer returned zero/short/negative/oversized byte count with nil
error (FAIL581). Explicit error/full-write controls behaved correctly. These
writers violate io.Writer's contract; this is local defensive validation of the
candidate's delivery claim, not a vulnerability or fix attributed to Go.

New checkedJSONWriter validates the count around the existing writer before
Go JSONHandler can report success. It makes exactly one Write call; a short count
without error becomes io.ErrShortWrite, and an impossible count becomes0 plus
io.ErrShortWrite. TryEmit preserves the existing static ErrDeliveryUnknown
boundary without formatting the underlying error or panic. There is no blind
retry. Full-length writes still do not prove flush/fsync, remote acknowledgement,
retention or a complete business transaction. The writer may still block under
the existing synchronous contract; lifecycle/deadlines belong to its owner.

The repair was returned to both canonical blocks and independently reconstructed
with2/2 identical hashes. No new dependency, source archive or production wiring.
Admission remains CANDIDATE; no composition guard was weakened.

Sources checked: [io.Writer contract](https://pkg.go.dev/io#Writer),
[JSONHandler.Handle](https://pkg.go.dev/log/slog#JSONHandler.Handle).
Execution uses observed Go1.26.7 and its locally hashed io/slog sources, not the
current version selected by those documentation pages. Security/SRE manual §15
requires private bounded telemetry, observable failures and actionable diagnosis.

## Executed integration

A canonical isolated reference test starts an actual loopback HTTP server and
an actual restricted-permission local file.32concurrent synthetic requests carry
a private marker in query/header; the test handler does not forward those values
to telemetry. Approved numeric correlation slots are confined to1..32, without
creating metric labels. Each request writes start/completion through the real
policy logger and Go JSONHandler into the file, and updates Counter/Histogram.

After explicit file Sync/Close, a fresh read verifies64complete JSON records,
32paired correlations in order, exact allowed keys, no private marker, count32
and histogram count/bucket sum32. This measures the synthetic reference path;
it does not prove the actual application's routing or auth, universal privacy,
OTLP/Prometheus export, production SLO or physical-storage durability after power loss.

A second real-file probe closes its sink before emission: TryEmit returns
ErrDeliveryUnknown and the closed file remains empty. The caller selects a new
file/logger and writes a new heartbeat successfully. It does not retry the
uncertain event or claim automatic repair of a truncated stream.

## Verification

- Combined candidate and independently rebuilt source:28unique top-level tests,
  each repeated3times;17distinct fuzz seed cases, zero skips.
- All6writer boundary modes plus complete-write control; four baseline failures
  become PASS without removing assertions. Underlying writer called once/event.
- test/vet/build PASS for candidate and rebuilt source with offline dependencies.
- GO-NATIVE-FUZZ-GATE0.1.0 materialized4/4 with hash verification;3existing targets,
  10seconds each, GOMAXPROCS2, PASS. No failed corpus was discarded.
- Go1.26.7 Windows/amd64 executable and stdlib source hashes recorded. No new SCA,
  race-detector or load/SLO claim. CGO_ENABLED0; prior compiler/race condition remains.

TEST-05 gains this integrated reference evidence but remains blocked for the
selected runtime's real host/reporter/supervisor/export/retention/alert path.
TEST-11 gains the corrected delivery-boundary evidence. Neither test's acceptance
criteria were weakened to reduce the pending count; no extra contract test row
was created. Integration test glue remains a test, not secretly mounted middleware.

## Continuity and scope

Discovery mistakes578–580/582 were read-only and recovered by exact inventory
or native flags. Their local lessons remain;581is the actual code defect.
The previous pnpm source/licence work is preserved but is not the focus of this
slice. HistoricalD, deferred ARCA, no real accounts/providers/spend and all target
conditions remain. The repaired candidate is absent from all selected composition
profiles, so unrelated154-step Preflight is retained with its original date;
the affected isolated full suite and structural gate are executed instead.

Expected structural inventory:162packs/1461materializable files/780Markdown/
53profiles. Source file counts and composition selections remain unchanged.

Stage: %TEMP%/elite-v342-9b34063622f44f5782a8e63857306dc1.

| Evidence | SHA256 |
|---|---|
| baseline-pack.md | `e6fc1d218d0a2790bcd5d5c06f307d37bde7a0612d6ecae5cad2bf2e5a61970f` |
| baseline-manifest.json | `84e50d8cb6d6850c011a0ed7cabb21f078563ee1913f477f33493716141bd01a` |
| delivery-red.log | `9fbee3ad7ccf0b48e1ff3da14e790d678eabd0afdd7eb6b2f807c6c2f5b24873` |
| delivery-green.log | `e0c393bedf0f3e4150fd52f61b7e0327709aab957491b36f1f7962a6a158d475` |
| reference-first.log | `7432c6efc0a4ec0d7d57c20536bba6c7fe4e25199c74e7aaf26bd48d86896ff5` |
| canonical-parity.json | `d04751d3f9638642ececc29c75c9b1d6f435576b73c0909527abda072a6f8881` |
| combined-gates.json | `b7b8a01506059a9c805a209ef166b290704bc77933a1228a76a4211593d30e7e` |
| rebuilt-gates.json | `a5977c06fd220194153d6db410d6cbbff372ec137fcca15a9726a4582d3b109c` |
| test-summary.json | `038e376542f5fbb3c655fc2086fb14df3770d16e5be98514f6341ad64e8e999d` |
| toolchain-receipt.json | `294a8a9ef0cc90ac16273652a8fc5535934c198969ac1fde3bcc6050f29dafa7` |
| rebuilt-tests.log | `b34657b239934a7560f55f4008da4abe2bb982292d63ea453430e4aa271a3d9f` |
| fuzz-profile.json | `43da59b89c64a4973bd0bee222ba0956b3a500f587fefa6ac988bb4b63e52824` |
| fuzz-receipt.json | `10a17d7e84d0dfe71067f6c77b8f012768f5f3eead037cd0bdc6f336803afecd` |
| fuzz.log | `637b2629675dda1468730ec3eb0835ce554e63aaa5c2a47f791954cd495212bf` |
| discovery-errors.txt | `86e175f4c0d275507413a8051afc63f699a8aecee3f117956189d6f7eb6ddca8` |

Integration helper count correction (FAIL583): the ledger contained2507unique IDs, not the manually expected2508. Recovery derives uniqueness/delta from the preserved baseline and records this incident separately. No code or test result changed; state write had not occurred.

## Closure126

VERIFY_LIBRARY125 PASS:162packs/1461materializable files/780Markdown and53unchanged composition counts. Exactly one implementation-pack Markdown changed (GO_OBSERVABILITY_CORE); its two source/test blocks equal the tested clean rebuild. The other162directory Markdown files are unchanged. Ledger uniqueness matches the recorded baseline/delta. No code changed after its green tests/fuzz or the structural gate. Final checkpoint/resume and contract-plan checks follow.

This closes FAIL581 for the local candidate and adds a real HTTP/file reference test. It does not close selected runtime mounting, supervisor, export, retention, alert channel, race/SLO or production admission; TEST05 remains explicitly open.

| Closure evidence | SHA256 |
|---|---|
| verify-library125.log | `c6b8c07c96e9a87d70aa5d20d2f1e4acaf8713222b5e5eb4db5c10085bcb529e` |
| structural-closure.json | `8d4ac534c4f7d36cc495c332449c66717e7454061c41a6771e0b0369d9ba0ea5` |
| final-uniqueness.json | `1329e638220595416b0499a1532f7b3fff279d0d4ee6655581f6081fb3189a1b` |
| ledger-delta.json | `a7f35795f220747b0f99e8b13b8c4aaea9864276f2749e1e5469c072a975622a` |
