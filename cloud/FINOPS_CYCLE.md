# FinOps scheduled cycle — authored orchestration

Observation: 2026-09-14T03:49:40Z. Scope: library infrastructure, development/staging. Status: `PROVEN_LOCAL` for the fake-provider scenarios below; all BigQuery, GCS, Pub/Sub, IAM, Cloud Run Job and Scheduler target execution is `NOT_RUN`. No account, credential, actual billing row or real notification was used.

`finops_cycle.py` is AUTHORED glue over the existing `BillingJournal`, `billing_query`, `normalize_billing_rows` and `execute_step.load_runner`. It adds no SDK/dependency and does not call the restriction or budget-management APIs. Its only remote effects, after an explicit expiring approval, are a bounded read-only billing query, generation-conditioned controller-state writes and optional publication to one named Pub/Sub topic. Persisted alert evaluation, accepted Pub/Sub publication and delivery to a person are separate outcomes.

## Invocation contract

Python API:

```python
init_journal(path)  # Offline, absent destination required.
cycle(config, approval, now, tools, runner, workdir, fixture=False)
```

`now` is an aware UTC `...Z` timestamp. The public CLI derives it from the clock; only fixtures supply their own. `workdir` must be absent, private, and free of symlink/junction ancestors. Every process uses the existing bounded runner with a 180-second deadline, captured output bounds, shell disabled and environment names explicitly allowed in the approved tool references. A fresh directory is required for each container/job attempt. A retry must not delete or bypass its durable intent.

Offline seed and fixture commands (from the cloud directory):

```text
python finops_cycle.py init --output /tmp/elite-finops-seed.sqlite
python -W error::ResourceWarning -m unittest test_finops_cycle -v
```

Future authorized container argv, no shell interpolation:

```json
["python3", "/app/cloud/finops_cycle.py", "run", "--config", "/var/run/elite/finops/config.json", "--approval", "/var/run/elite/finops/approval.json", "--tools", "/var/run/elite/finops/tools.json", "--runner", "/var/run/elite/finops/runner.json", "--workdir", "/tmp/elite-finops-cycle", "--authorize-execution"]
```

The seed must already exist at `journal_uri`, uploaded by a separately approved generation-match=0 step. This cycle never interprets a permission error or unavailable object as an empty state. Job/image/secret-reference/Scheduler composition belongs to `stack_manifest.py` and the cloud integration owner. Schedule interval and retry policy must match this approved config; scheduler activation is a billable/external action and remains `NOT_RUN`. No authorization file is minted by this module.

Exact config fixture (values are illustrative, not an authorized target):

```json
{
  "schema": "elite-finops-cycle/v403.1",
  "production_authorized": false,
  "environment": "development",
  "project": "elite-fixture",
  "location": "US",
  "billing_account": "AAAAAA-BBBBBB-CCCCCC",
  "billing_table": "elite-fixture.billing.export_table",
  "currency": "USD",
  "period_policy": "UTC_CURRENT_MONTH",
  "maximum_bytes_billed": 10000000,
  "interval_seconds": 300,
  "max_state_bytes": 1048576,
  "max_age_seconds": 600,
  "limit_micros": 1000000,
  "journal_uri": "gs://elite-fixture-state/control/billing.sqlite",
  "fence_uri": "gs://elite-fixture-state/control/finops-intents",
  "notification": {
    "mode": "PUBSUB",
    "topic": "projects/elite-fixture/topics/billing-alerts",
    "dedupe_seconds": 3600
  }
}
```

All top-level config fields are required; unknown fields are rejected. `notification: {"mode":"NONE"}` removes the topic/dispatch requirement. Numeric costs are exact integer micros; booleans/floats/string limits are rejected. The current UTC invoice month is queried once per approved interval; closed-month reconciliation, FX and multiple billing accounts require a separate scope. The target must prove the configured table is the complete standard export for the approved account. The module does not infer that ownership from the table name.

Approval fields:

| Field | Required value |
| --- | --- |
| schema | `elite-finops-approval/v403.1` |
| method | `EXECUTED_TARGET`; fixtures exclusively use `SIMULATED_PROVIDER` |
| config_sha256 | `cloud_control.digest(config)` |
| runtime_sha256 | `runtime_identity(tools, runner)` |
| approved_by | Explicit operator identity |
| valid_from / expires_at | UTC validity window including the current run |
| billing_table_account_verified | `true`, supported by target evidence |
| query_cost_authorized / durable_state_authorized | `true` |
| dispatch_authorized | Exactly whether mode is `PUBSUB` |

`tools` has exactly `bq` and `gcloud`, each with the existing binary SHA, method, admission and full runtime admission receipt contract. Prefix arguments are forbidden. Real tool references require `ADMITTED_EXACT_ARTIFACT` and `tool_admission.validate`, including complete runtime inventory and G0–G8 reports. A gcloud admission does not automatically admit a distinct bq launcher; its executable hash and runtime receipt must also pass. `runner` is `{path,sha256}` for the existing bounded runner. Configuration/runtime digests use the existing owner's canonical JSON serializer, not an arbitrary external JSON formatting. Credentials are externally supplied to the workload through approved identity/environment references; inline credential values are absent.

## State and alert behavior

1. A stable config-hash/time-slot ID claims `fence_uri/<id>.json` with GCS generation-match=0 before querying. An existing or ambiguous claim fails closed without an automatic replay.
2. Describe the SQLite object's generation, size and transport checksum; download exactly `gs://...#generation`. Reject oversize or inconsistent downloads. SQLite integrity is checked locally.
3. Execute only the owner's generated account/month `SELECT`, adding an explicit `maximum_bytes_billed`, project/location, empty controlled bigqueryrc and the existing row cap. Reject truncated/capped, duplicate, mixed-currency, boolean/float or overflowing amounts. Signed credits and account-level unallocated cost are preserved.
4. Import a new complete snapshot generation through `BillingJournal`; replace the previous account/month aggregate instead of adding it twice. Compute the existing aggregate alert policy. Store the cycle and an optional outbox claim; create a standalone SQLite backup including committed WAL data.
5. Upload state using the original observed GCS generation as CAS precondition. Verify the resulting generation/size/checksum. A concurrent state writer aborts before any publication. MD5 is used only for GCS transfer equality; source/receipt/command identities use SHA-256.
6. Only after a durable claim, publish one sanitized aggregate payload to the approved Pub/Sub topic. Its stable `id` is available for subscriber deduplication. Persist `PUBLISHED_NOT_DELIVERED` only after a nonempty provider acknowledgement. An ambiguous result keeps `CLAIMED_RECONCILE_BEFORE_RETRY`; a failed second state CAS also leaves the previous durable claim. The adapter has no automatic publish retry. Transport or Pub/Sub delivery can still duplicate; this is not exactly-once delivery.

Deduplication applies to the approved config/month/window. A later configured window may emit a new alert evaluation; that is a new notification policy event, not a replay of an uncertain earlier event. Operators must reconcile uncertain publications before changing policy to force a resend. `NONE` persists the evaluation without dispatch. Local receipt output contains approved aggregates and command hashes, not raw CLI stdout/stderr or token values. Remote state remains the durable authority if a local worker exits.

Fresh observation time proves when the query ran, not how current Google's export is. Every report says `NOT_PROVEN_EXPORT_MAY_BE_DELAYED`, retains residual-cost caveats and sets `hard_cap=false`. No provider, tenant or global spending cap is claimed. `max_age_seconds` concerns observation age only. The cycle does not sum internal estimates into provider billing, infer provider latency, close prior months automatically or stop resources. State history growth fails closed at `max_state_bytes`; retention/archive policy, expiry of durable intents and restore are explicit target responsibilities, not silent deletion in this module.

Target gates remain: complete account/table mapping and currency, export latency/completeness, IAM for read-only billing/query jobs and the exact state/intent namespace, approved query cost/interval, independent bq/gcloud runtime admission, seed upload and restore, CAS against real GCS, workload identity, Job/Scheduler activation, target failure recovery, Pub/Sub subscriber deduplication/delivery and operator acceptance. Actual notices and restrictions require exact separate authorization. None were executed here.

## Local evidence and preserved failures

Command executed in `C:/Users/NL/Desktop/Elite Library Extension V403/cloud`:

```text
C:/Python314/python.exe -W error::ResourceWarning -m unittest test_finops_cycle -v
```

Final observation: Python 3.14 / Windows, 10 tests PASS, 0 failures/errors. Fake runner implements generation-bound reads, create-only claims, conditional uploads and fake Pub/Sub acknowledgements. Tests cover: aggregate including signed credits; subsequent snapshot replacement; durable state before dispatch; repeated-slot query prevention; concurrent CAS rejection before dispatch; lost publish acknowledgement; failed state CAS after publish; deduplication of durable ambiguous claims; explicit NONE dispatch; exact config/approval/method/currency/row-cap/amount/SQL/scope validation; and local seed overwrite refusal. All provider methods are `SIMULATED_PROVIDER`. This does not prove the real CLI parser, IAM, network, physical hardware, external billing export or message delivery.

Preserved progression:

- First 10-test run: 6 errors, `OSError: [Errno 9] Bad file descriptor` when Windows `os.fsync` received a read-only (`rb`) handle. Corrected to `r+b` before the CAS upload; no exception suppression.
- Second run: core assertions passed, but 4 teardown errors exposed a fake-provider SQLite connection left open by its transaction context manager. Explicit `finally: db.close()` fixed the harness. No real state/provider was affected.
- Third run: all 10 PASS with ResourceWarning promoted to error. Parent was notified for the shared failure ledger; earlier observations are not represented as PASS.

Source observations, exact SHA-256 at the timestamp above (the cloud integration owner may subsequently reseal changed owners):

| Local file | SHA-256 |
| --- | --- |
| cloud/finops_cycle.py | `35708a02289e662ed9037900c1ea198fe6db6c941d7b47f0210c383f84f67ec3` |
| cloud/test_finops_cycle.py | `e823886003c93cf729ea559b1ef8a4921fdc542ececabb830ef9d3628a01609f` |
| cloud/cloud_control.py | `8c7de864a3e29774a19705edc15332f8ee0225b36538d303273f1c2b77383f81` |
| cloud/execute_step.py | `63f37cd8add209ebe27409bc4eaef43a54d727fc61731986dfb25dc621a91f05` |
| cloud/tool_admission.py | `0d09f6ccd82fe45464261f9a2eda6d5597a3a8071cd8732ff915d11d0753da05` |
| cloud/owners/secure-ops/production_admission_gate/run_official_tool.py | `41e73b0245c20f5056caf8ae818f0ce61dc502fedf4b28caa245cd3b0f09c5e9` |

Official authority observations (read via web on 2026-09-14; documentation only, no external code acquired; URLs are moving documentation, not artifact admissions):

- [BigQuery bq reference](https://docs.cloud.google.com/bigquery/docs/reference/bq-cli-reference): exact CLI parameters for query, integer maximum billed bytes, rows, project/location and bigqueryrc. A query exceeding the stated byte limit fails without charge; this is per query, not a global account cap.
- [GCS copy command](https://docs.cloud.google.com/sdk/gcloud/reference/storage/cp): generation precondition and content-MD5 parameters.
- [GCS versioned objects](https://docs.cloud.google.com/storage/docs/using-versioned-objects): generation-specific object URL for reads.
- [GCS metadata](https://docs.cloud.google.com/storage/docs/viewing-editing-metadata): describe output includes generation and `md5_hash`.
- [Pub/Sub topic publication CLI](https://docs.cloud.google.com/sdk/gcloud/reference/pubsub/topics/publish): topic and message arguments. Publication acknowledgement is not subscriber or human delivery evidence.

No official document bytes are copied into the implementation. The guide records the URLs and observation method; the listed local file hashes are byte observations, not hashes of remote web pages.
