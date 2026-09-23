# Official Meta Page write leaf

This AUTHORED adapter invokes the official Meta Business SDK 26.0.1 pinned by
exact wheel graph and commit. It provides one POST of approved text to the selected
Facebook Page, bound GET observation/reconciliation, and separately approved
revocation. It adds no Instagram, TikTok, LinkedIn, ad mutation, media or paid boost.

Install the exact existing CPython 3.14 Windows lock into an isolated environment:

    python -m pip install --require-hashes --no-deps -r meta_page_write/requirements-windows-py314.lock
    python -m pip check
    python -m unittest discover -s meta_page_write/tests -v

The dependency lock is unchanged from the admitted Meta Ads reporting runtime;
its read-only adapter is not reused or relabeled as a write adapter. Meta licenses
permit use with Meta APIs; LICENSE.Meta.txt is mandatory, including for the CAPI
wheel whose license file is absent. This is not an OSI/general-purpose license.

Create an API with create_api(app_id, app_secret, page_access_token), then construct
FacebookPageWriteAdapter(Scope(...), api). Credentials resolve in the caller's
secret manager/environment; this module has no CLI or embedded credentials. Graph
v26.0 is fixed; debug logging is rejected. Calls use the official SDK transport.
Tests mount an in-memory requests transport, so no network publication occurs.

The caller supplies ApprovedIntent only after its existing identity, approval,
revocation, scheduling, quota and durable outbound-fence owners authorize it.
The approval ID is a correlation binding, not proof that this leaf has performed
authorization. Tenant/page/profile/content hashes must match the configured
scope; content is preserved byte-for-byte. The 16384-byte ceiling is a local
transport/memory bound, not a claim about provider character limits.

publish_once performs one POST then a GET. A receipt requires the exact post ID,
Page owner, message and published flag from the GET. A timeout or mismatched result
is UnknownDelivery; do not resend. With a known ID, reconcile performs GET only.
Without an ID, keep the durable fence unresolved pending owner reconciliation;
this adapter never guesses a post from similar text. revoke_once requires a
separate revoke intent, checks the original bytes/Page first, then requires raw
SDK DELETE success=True. It does not infer deletion from an empty object or a
permission/error response. Receipt hashes cover the normalized verified fields,
not an unmodified raw HTTP body. Requests and secrets are not included in receipts.

Scope still open for the complete social core: durable schedule/approval binding,
lease/quotas, dispatch worker wiring to communication.outbound_delivery, receipt
persistence, notification/UI, operator reconciliation and target token lifecycle.
These are local code/integration gaps and are not only missing credentials.
The existing outbound fence may be reused; this leaf does not claim that wiring
has been performed. The six other core gaps are listed in seven-core-decisions.md.

Official Page SDK methods were checked at their immutable commit. Current Meta
Pages API/permissions documentation endpoints returned tool retrieval errors;
no third-party mirror was used as authority. Live permission/app/Page ownership
contract must be verified before enabling a target. No live-readiness claim.

Verification: seven focused tests with real locked SDK objects/requests transport,
including wrong tenant/page/profile/content, unconfirmed observation, GET-only
reconciliation, POST timeout and false/missing/timeout DELETE outcomes. Source
files installed from the wheel match the fixed official source. First revoke
failure and correction are preserved in PROJECT_FAILURE_LESSONS.md. No core,
Commerce, Accounting, PG or previous fuzz suite was rerun.
