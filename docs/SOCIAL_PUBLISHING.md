# Facebook Page publication infrastructure

This materializable scope publishes exact reviewed text to one configured
Facebook Page. It includes a runnable Go API/worker, immutable human approvals,
durable scheduled jobs, a single-write outbound fence, the official Meta Python
SDK, exact GET observation, public receipts and separately approved revocation.
No marketing policy, moderation, audience inference or eligibility is generated.

## Owners and provenance

The existing approval.Registry owns human separation; approval.request/decision
remain the only approval ledger. The new typed binding supplies zero automatic
approval policy. Its three manual kinds are social_publish, social_revoke and
whatsapp_reply. The generic PostgreSQL owner is shared with WhatsApp.
platform.job owns scheduling, bounded retry budgets and claim generations.
communication.outbound_delivery owns the at-most-one write attempt. Public
provider references/receipts live in platform.outbox_event with immutable social
evidence fields. No social ledger or replacement financial logic is introduced.

All new files are AUTHORED configuration/serialization/authorization/persistence
and process glue. The seven-file Meta SDK adapter is an unchanged dependency.
Provider POST/GET/DELETE behavior is DEPENDENCY_PIN facebook-business26.0.1,
fixed commit788f363d15b1269ab5efb7cd00fb5e3b133cd99b / Graphv26.0. The exact
LicenseRef-Meta-Platform bytes remain owned by that adapter. Nothing here is
attributed to a vendor as a locally adapted marketing algorithm.

## Start

1. Compose the admitted backend, jobs, outbound fence, human approval, seven-file
   Meta Page adapter and this publishing integration. Apply migrations through
   0058 plus0060 (0059 belongs to the independently selected WhatsApp delta).
2. Install the exact18-wheel requirements-windows-py314.lock into an isolated
   CPython3.14 Windows environment using --require-hashes --no-deps. Keep all
   source and license notices supplied by its dependency pack.
3. Review a profile derived from deploy/social/profile.reference.json, replacing
   the clearly synthetic tenant/organization/Page identifiers during deployment.
   Select its SHA256 independently. All worker replicas use the same bytes/hash.
   Queue/timing/attempt bounds are technical configuration, not marketing policy.
4. Supply the empty environment inputs from deploy/social/.env.example when
   activating. No secrets are requested or embedded now. Pin the materialized
   meta_page_write/bridge.py SHA256 and isolated Python executable path.
5. Build `go build -mod=readonly ./cmd/social-publishing` and run the resulting
   executable. It serves on loopback127.0.0.1:8097 by default and runs the worker.
   OIDC verifier owns identity; the caller needs the configured organization and
   social:request/social:approve/social:read/social:reconcile permissions as used.
   Front a non-loopback deployment with the target's admitted TLS/auth gateway.

Before serving, startup rejects missing profile/hash, credentials, runtime,
provider SDK source drift, dependency version drift and a changed Python bridge.
The SDK preflight is offline and verifies fixed source bytes; it does not prove
Page ownership, Meta app access approval or live token permissions.

## Operator contract

POST /v1/social/requests submits the Request JSON documented by
internal/socialbridge/contract.go. Server-owned tenant/Page/profile checks reject
cross-scope values. ApprovalID identifies the immutable request; DeliveryKey is
social_ plus SHA256(tenant NUL page NUL approvalID NUL operation). ContentSHA256
binds the exact UTF-8 message. ScheduledAt and ExpiresAt are explicit UTC instants.
The response provides the canonical RequestSHA256 for exact review.

GET /v1/social/requests/{id} exposes the exact pending payload/hash and current
approval, delivery and receipt in the authorized tenant/organization.
POST /v1/social/requests/{id}/decision requires request_sha256, explicit boolean
approve and reason. Its verified human must differ from the requester. Only the
approved transaction creates the scheduled job. Rejection creates no job.

The worker revalidates durable approval, tenant, profile, page, content,
not-before/expiry and current claim generation in the fence transaction. A replay
cannot obtain another write attempt. Its Python child calls the real SDK POST
once, then GET verifies exact id/from.id/message/is_published before success.

After an uncertain response, known post IDs are persisted and retries perform
GET only. POST /v1/social/requests/{id}/reconcile provides the same GET-only
recovery after the queue budget is exhausted. It preserves the attempt number.
A saved receipt recovers a crash before fence/queue acknowledgement without
another SDK call. Restarting the host preserves approvals, jobs and evidence.

Revocation is another request, operation=revoke, separately reviewed. It names
the original approved publication and its exact public receipt/postID/content.
The worker verifies that binding, then the SDK performs GET and DELETE once.
A different Page, post reference, message or original approval cannot authorize
deletion. A rejected pending request never dispatches. Disabling an action in a
new reviewed profile prevents new admission under that profile; do not assume
changing a file can recall a provider operation already in flight.

## Ambiguous external outcomes

A lost POST response with no provider ID remains UNKNOWN; neither the operator
endpoint nor a restarted worker invents a post ID or resends it. A lost DELETE
acknowledgement also remains UNKNOWN: absence or permission denial cannot prove
deletion. These are explicit preserved external outcomes, not successful sends
or missing implementation. Retain the request and seek provider evidence through
the operator process. No credential alone is claimed to resolve such a case.

## Retention, rollback and boundary

Planned publication text is stored in the access-controlled approval payload;
tokens are never stored there. Provider receipts contain public IDs, hashes and
state. Do not send private customer records as publication text. Target retention
configuration must preserve approval/receipt evidence needed for reconciliation.
Stop admission and drain workers before changing executable/profile versions.
The down migrations refuse removal once their governed evidence exists; use a
reviewed forward change rather than deleting that evidence.

Focused PostgreSQL/real-SDK-fixture tests cover exact approval, replay, scheduled
admission, tenant/Page/content isolation, stale claims, one POST, GET-only
recovery, terminal recovery, receipt crash, expiry under a lock and governed
DELETE. These library results do not certify live app/Page permission, target
load/availability, broader social networks or a fresh composition-wide security
review. The parent composition owns its T2803 SCA/SAST and final admission.
