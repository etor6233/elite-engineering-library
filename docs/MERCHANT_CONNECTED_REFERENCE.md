# Connected Google Merchant reference

Library fixture claim: current reviewed catalog, original serial ATP and exact
integer minor-unit price -> manual distinct approval -> official SDK1.8.0
ProductInput insertion/upsert -> immutable PostgreSQL receipt -> processing GET.
The Go/Python process, SQL and HTTP binding is AUTHORED glue. Google owns the
unchanged generated SDK; original source/SDK wheel/20 dependencies are fixed in
the existing sdk-artifact.lock.json and requirements-windows-py314.lock.
docs/merchant/official-contract.lock.json fixes current official API docs.

Select a dedicated API primary data source and a stable offer namespace in
config/merchant.reference.json. It contains synthetic examples, no credential.
Bind its variant to the reviewed catalog and set model.specification's
merchant_description explicitly. Source generation is the official positive
int64 version_number. Updates and refreshes use insert/upsert at that version;
source rollback is a new monotonically increasing local publication generation.
No takeover of a foreign version namespace is inferred. Supplemental/rule
transformations remain observable divergence, not approval.

## Runtime and host

Use the fixed Windows x64 / CPython3.14 lane. Acquire the exact twenty wheels
from their URLs/hashes in connected-runtime.lock.json (the original hash-locked
requirements file is the same graph). Run install_connected_runtime.py with
--wheels ABS_CACHE --destination ABS_ABSENT_RUNTIME --python-sha256 BASE_SHA
--lock-sha256 LOCK_SHA using the admitted base Python. It validates all wheels
before creating a venv, installs offline without dependency resolution, runs
pip check, verifies1041 installed payload files and emits merchant-runtime.json
plus its hash in setup-receipt.json. The script never reads an account/secret.
Source-specific fixtures run with test_connected_worker.py --python VENV_PYTHON
--receipt ABS_ABSENT_RECEIPT. This uses the actual fixed REST SDK on loopback.

At future activation supply MERCHANT_ENABLED=true,
MERCHANT_PROFILE_FILE / MERCHANT_PROFILE_SHA256,
MERCHANT_RUNTIME_FILE / MERCHANT_RUNTIME_SHA256,
MERCHANT_HMAC_KEY (32 bytes, base64) and GOOGLE_APPLICATION_CREDENTIALS through
the runtime secret environment. Catalog activation/profile/origin must match;
migration0082 and original immutable approval/fence guards must be active.
The emitted runtime defaults to CREDENTIALS; only fixture harnesses use
LOCAL_FIXTURES plus an explicit127.0.0.1 origin and no credential file.
Enabled invalid configurations fail startup; disabled modules make no calls.
The old standalone sync_product.py remains an admitted lower-level utility;
the connected host uses its exact request builder/models, with PostgreSQL
receipts and the original outbound fence replacing its standalone file runner.

## Operator workflow

Authenticated tenant/organization role endpoints under /v1/admin/merchant:
- GET /queue: bounded per-profile view (maximum32 offers). Read only.
- POST /prepare: approval_id, generation as string, variant_id, expires_at.
- POST /requests: submit the exact returned intent.
- POST /requests/{id}/decision: request_sha256, explicit approve, reason;
  reviewer must differ from requester.
- POST /requests/{id}/send with {}: one provider attempt for that approval.
- GET /requests/{id}: durable approval, fence, input receipt and processing.
- POST /requests/{id}/reconcile with {}: GET only, never resends the insertion.

Permissions are merchant:read/request/approve/send/reconcile, with verifier
tenant and organization scope. All request bodies/time/process/response sizes
are bounded; duplicate/case-folded/unknown JSON fields fail closed.

CURRENT_INPUT means an acknowledged input with a future refresh deadline,
not Google product approval. A GET match after a lost insertion proves desired
state but leaves freshness unconfirmed. Input acknowledgements alone anchor
refresh_due_at. At REFRESH_APPROVAL_REQUIRED or changed source, prepare a new
request and review it; expired/manual approvals are never auto-renewed.
UNKNOWN/SENDING blocks new attempts for that offer across profile changes.
Missing resources, incorrect data source/version/fields and rejected writes
stay explicit. Product status/issues are preserved as exact SDK JSON bytes.
Polling /queue does not create jobs or send products; its operator read path is
the selected manual feed-renewal workflow. Automated campaigns/reminders are
the separate T2805 notification owner. This reference selects upsert/refresh
of its bound offers; bulk feeds, multi-source takeover and automated removal
are not claimed as equivalent provider operations.

## Evidence and limits

MERCHANT_CONNECTED_RELEASE_V402.md/json binds the current pack, SDK/runtime
restore, connected API/PG proof, host/migration and finite fuzz results.
Input processing can be pending or altered by Google; no insert200 or local
PASS demonstrates product approval, live account acceptance or advertising.
Future account credentials/access are the connected provider activation
condition. Local fixtures require none. T2805 remains open until its other
selected mappings and communications workflow are connected.
