# Offline full-runtime stack composition

`stack_manifest.py` is AUTHORED orchestration over the selected library owners.
It emits configurations, Google Cloud Run manifests and bounded argv arrays.
`build_stack(config, images, source_root=..., composition=...)` returns a plan;
`validate_stack(plan)` checks structural integrity. Neither function executes a
provider. A structural PASS is not authorization, image admission or target
runtime evidence. All generated commands retain `execution: NOT_RUN`.

From the `cloud` directory, using the already available Python runtime:

```powershell
C:/Python314/python.exe -B -m unittest -v test_stack_manifest
C:/Python314/python.exe -B stack_manifest.py emit --config stack.fixture.json --output evidence/stack/generated-review
C:/Python314/python.exe -B stack_manifest.py validate --plan evidence/stack/generated-review/stack-plan.json
```

The output directory must be absent. Existing observations are not overwritten.
For a sealed consumer, supply `--source-root <consumer-composition>` and
`--composition <sealed-consumer-composition.json>`. The config must bind that
exact composition and inventory. Base116 is retained; a 120-pack consumer is
accepted only with the exact four V403 extension identities and extension
policy. Consumer provenance verification still belongs to `seal_consumer`.

The fixture deliberately names a nonexistent project, `.invalid` principals
and synthetic immutable image digests. Its method is `SIMULATED_ARTIFACTS`.
Six logical image families represent five OCI images: API and workers can
share one digest containing all `/app/bin/<cmd>` binaries; web runs
`node /app/server.js`; migrations run the PostgreSQL client; ARCA runs the
existing certificate launcher; controller runs `python3 /app/cloud/finops_cycle.py`.
Actual source-to-image provenance, manifests,
executable files, licenses and scans must be verified before target execution.
Do not turn fixture declarations into execution receipts.

| Configuration | Generated behavior and source responsibility |
| --- | --- |
| api | Real `cmd/electromobility-api` HTTP process and existing payment/WhatsApp/portal loops; minimum one instance and continuously allocated CPU. |
| web | Next standalone service; minimum zero, bounded maximum; explicit public catalog tenant/org, OIDC/session refs and metadata service identity. |
| go-workers | Return effect, exchange and accounting loops are manually scaled WorkerPools, each one instance with its stable worker ID. No invented HTTP endpoint. |
| arca | When enabled, API, Go fiscal/parameter workers and .NET launcher share `/run/arca`; return-fiscal remains a WorkerPool. All related processes are absent when deferred. |
| migrations | One finite Job, no retries, exact selected pending SQL suffix and hashes. Cursor is an index into 84 selected files; missing historical ordinals 63/64 are intentional. Each SQL owns its transaction. |
| database | Private Cloud SQL PostgreSQL 18, explicit tier/database/user refs, PITR, retained backups and deletion protection. Existing VPC/private service access and password provisioning are target prerequisites. |
| storage | Separate quarantine, retained and controller buckets, public access prevention, uniform IAM, explicit versioning and soft deletion. Controller CAS state has no bucket retention policy; document retention remains explicit. |
| identity | Separate runtime, migration, scheduler/controller accounts; SQL client only for database users; Secret Manager access bound per secret; direct IAP browser principals and independent web-to-API invoker grant. |
| scheduler | Explicit tenant/org and bounded retention command; Cloud Scheduler invokes the finite Run Job through the Google API with OAuth, zero automatic retries. |
| documents | REVIEW_ONLY and automatic persistence false, exact profile reference/hash required; no document provider is silently chosen. |
| cost-controller | Finite `finops_cycle` Job over existing BillingJournal/bq/bounded runner; separate generation-CAS SQLite and create-only intents, exact versioned mounts and optional Pub/Sub. Scheduler has no automatic retries. This is an aggregate alert evaluation, not a hard spending cap. |
| ci | Separate library/project repositories, pinned library commit, target/repository-scoped WIF principal and deploy identity. Existing portable build/execution owners remain authoritative. |

All 14 selected Go entrypoints are inventoried and hashed. The older reference
`cmd/api`, file import CLIs and warranty operator CLI cannot be deployed as
servers or unattended jobs. Social publishing and Meta lead HTTP ingress are
explicitly deferred in the fixture because their provider bindings remain
absent. Embedded API loops remain governed by their original profiles and
durable claim/reconciliation logic; materializing 116 packs does not admit
every external business operation.

The ARCA helper imports a PFX into CurrentUser/My, using fixed separate secret
mount directories `/var/run/secrets/arca/certificate/cert.pfx` and
`/var/run/secrets/arca/password/password`. Its source hash is bound when enabled.
The existing launcher copies mounted files into UID65532 private HOME with
0700/0600 permissions, imports, then removes private staging. The manifests do
not fabricate a TCP probe for the UDS-only .NET worker. Linux permissions,
socket readiness, sidecar restart, stale socket recovery and same-UID access
still require target acceptance. Only homologation is configured.

Google's [v1 YAML contract](https://docs.cloud.google.com/run/docs/reference/yaml/v1)
defines the three resource shapes. Persistent loops use
[manual worker pool scaling](https://docs.cloud.google.com/run/docs/configuring/workerpools/manual-scaling).
The scheduled job uses the documented
[OAuth Google API invocation](https://docs.cloud.google.com/run/docs/execute/jobs-on-schedule).
Database and storage argv follow the official
[Cloud SQL create reference](https://docs.cloud.google.com/sdk/gcloud/reference/sql/instances/create),
[bucket create reference](https://docs.cloud.google.com/sdk/gcloud/reference/storage/buckets/create)
and [bucket update reference](https://docs.cloud.google.com/sdk/gcloud/reference/storage/buckets/update).
Secret configuration follows the
[Secret Manager mount contract](https://docs.cloud.google.com/run/docs/configuring/services/secrets).
The fixture selects `identity.browser_access=IAP_DIRECT`, following the current
[direct Cloud Run IAP contract](https://docs.cloud.google.com/run/docs/securing/identity-aware-proxy-cloud-run)
and [IAP access policy command](https://docs.cloud.google.com/sdk/gcloud/reference/iap/web/add-iam-policy-binding).
The web manifest enables IAP from its initial configuration, grants the IAP
service agent Run invoker, and assigns allowed users/groups the IAP accessor
role. Application OIDC callback and secure cookies retain the same HTTPS
origin. The API does not enable IAP and keeps the separate metadata/user-token
transport. No browser login or target configuration was executed.

IAP's OAuth setup is an explicit prerequisite: Google-managed clients apply
inside an organization; external/no-organization cases require the documented
custom-client/Console setup. The generator does not fabricate programmatic
OAuth-client creation or credentials. The `iap_oauth_setup_allowed_users_and_app_oidc_callback`
and `iap_service_agent_exists` receipts remain required. The alternative
`IAM_API_ONLY` emits developer proxy argv but labels the browser infrastructure
`INFRASTRUCTURE_INCOMPLETE_PROTECTED_API_ONLY`; a local proxy alone does not
prove HTTPS callback or cookie compatibility.

The official HTML observations, acquisition dates, last-updated dates and
SHA-256 values are preserved in `evidence/stack/method-sources*.json` and
`official-docs/`. Google documentation is attributed under CC BY 4.0 and sample
code under Apache-2.0; this adapter copies no sample. No dependency was added.

Before execution, the existing bounded runner integration must regenerate the
plan from locked input and compare every manifest byte. Dependencies and each
named `required_target_receipts` must be checked against the same plan, image
set and target; a standalone structural validation cannot replace that gate.
Target acceptance still includes exact Linux images, private DB access,
migrations/restore, secret values, OIDC callbacks and dual authorization
headers, forced stop recovery, document corpus, billing freshness and CAS,
and semantic receipts for each actual provider operation. This delivery
records local fixture evidence only; no account, cloud resource, secret or
franchise was created.

## FinOps Job and Scheduler extension — 2026-09-14

The earlier cost-controller entry was configuration only. This extension adds
the executable finite Job and its Scheduler binding to the existing
`finops_cycle.py` implementation described in `FINOPS_CYCLE.md`. Configuration
now contains the complete `cycle` object, `job_name`, matching UTC `schedule`,
finite `timeout_seconds`, exact `cycle_owner_sha256`, four `mounts` references
and a hash-bound `seed_receipt_ref`. The four mount keys are `config`,
`approval`, `tools`, `runner`; each has only `secret_ref: name:numeric-version`
and `sha256`. The config SHA must equal the canonical generated cycle digest.
The other referenced bytes require target receipt verification; fixture hashes
do not certify mounted files. No secret values or approvals are generated.

Each mount has its own `/var/run/elite/finops/<key>/<key>.json` path. The Job has
one task, parallelism one, no ports, no Cloud SQL/VPC annotation and zero task
retries. Its task timeout cannot exceed the approved interval. The fixture
uses 900 seconds and `*/15 * * * *`; supported intervals are whole-minute
divisors of an hour, validated against the exact cron expression. The Scheduler
uses the Google Run API OAuth invocation, zero retry attempts and zero retry
duration. Creation/activation is an external action behind the explicit
schedule approval gate, not something executed by `emit`.

The controller's service account must differ from every application worker,
web, scheduler, migrator and deployer identity. It receives BigQuery Job User
on the query project and Data Viewer conditioned on the exact approved
export table's project/name/type/service. No wildcard table selection is
generated. GCS access uses a custom role containing only
`storage.objects.get/create/delete`, conditioned on
`control/billing.sqlite`; overwrite requires create+delete. The intent prefix
`control/finops-intents/` receives only Object Creator. It cannot rewrite the
deployer's `control/state.json` or delete another intent through these grants.
Four individual Secret Manager bindings serve the four mounted files. Pub/Sub
Publisher is emitted only for an explicitly configured `PUBSUB` topic; default
`NONE` creates neither a publish binding nor a topic. Target IAM must also prove
the identity has no broader inherited bindings.

Seed initialization is a mandatory BARRIER before Job configuration and
scheduling. `configs/cost-controller.json` emits the exact offline init argv
and `gcloud storage cp ... --if-generation-match=0` argv under
`seed_initialization`, with `execution=NOT_RUN`. The operator runs these only
after authorization, retains seed bytes/hash/generation and demonstrates
restore; the semantic seed receipt must bind that evidence to this plan.
There is no automatic empty-state fallback, upload in the generator, overwrite
of a prior seed or fabricated provider receipt.

New step IDs and gate responsibilities:

| Step | Required target gate |
| --- | --- |
| configure-finops-controller | artifact_same_digest, finops_runtime_and_mount_hashes |
| finops-query-user, finops-billing-reader | finops_billing_scope_and_query_budget |
| finops-journal-role | new_resource_absent, finops_namespace_iam |
| finops-journal-access, finops-intent-create | finops_namespace_iam |
| finops-seed (BARRIER) | finops_seed_generation_zero_and_restore |
| finops-invoker | Reconciled configured Job dependency |
| finops-scheduler | finops_schedule_approval_and_reconciliation |
| finops-pubsub-publisher (PUBSUB only) | finops_pubsub_topic_and_dispatch_approval |

The default fixture emits 50 steps and 8 manifests. Local command
`python -m unittest test_stack_manifest -q` ran 38 tests: 37 PASS and one
existing Windows symlink test SKIPPED because OS privilege 1314 is unavailable.
That skip is explicitly unproven here; no target or real provider ran. New tests
cover finite task shape, seed gate ordering, exact mounts, scheduler interval
and retry bounds, isolated IAM, PUBSUB/NONE distinction, scope/hash/byte-limit
drift, omitted controller image and forbidden identity reuse. The generated
plan digest is `8f44cc4ed30cfb00075d72912b10568af60790cbb3045e7031c9d19f8e38f2c1`.

The original fixture SHA
`7a23db82320b124afbbb80efee29763f77908005fb09031f1537f5c000a63823`
belongs to the earlier four-OCI qualification and remains recorded in its
historical receipts. This extension's fixture SHA is
`b9917f7b8a8ba8d67eeadcb0e5b75b98ed4169b8ba40b976fabb2d6968a9b7f1`;
generator SHA is
`8ca237c17f40ffa29b0ca656c379b1368ed0d0927ce9cc5347940cd01400b81e`;
test SHA is
`14a9d0c19bd04ea3c0d66c8c969db0b2e00daa799fd7b6f9995675bd54427f7c`.
Subsequent integrated source changes require a new receipt, not rewriting
these observations. Image runtime admission and semantic provider execution
are owned by the cloud integration gate and remain `NOT_RUN` here.

Additional current official authorities read on 2026-09-14 through the web
tool, documentation only; no external code or dependency incorporated:
[BigQuery IAM Conditions](https://docs.cloud.google.com/bigquery/docs/conditions)
permits inherited project conditions with explicit table name/type/service;
[GCS permissions](https://docs.cloud.google.com/storage/docs/access-control/iam-permissions)
requires create and delete for replacing an object;
[custom role creation](https://docs.cloud.google.com/sdk/gcloud/reference/iam/roles/create)
defines the generated finite permission set;
[bucket IAM binding](https://docs.cloud.google.com/sdk/gcloud/reference/storage/buckets/add-iam-policy-binding)
defines the condition parameter;
[Cloud Run volumes](https://docs.cloud.google.com/run/docs/reference/rest/v1/Volume)
defines numeric versions and read-only mode 0444. These new observations are
live documentation reads with URLs/date, not new downloaded document-byte
locks or provider acceptance evidence.
