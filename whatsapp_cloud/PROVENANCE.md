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
