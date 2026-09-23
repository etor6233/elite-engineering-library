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
