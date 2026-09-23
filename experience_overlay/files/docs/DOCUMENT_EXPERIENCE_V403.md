# V403 document experience — narrow UI/BFF extension

Owner remains GO-CONNECTED-DOCUMENT-REFERENCE0.1.0, selected in the unchanged
116-pack base. This extension adds AUTHORED UI/transport glue; it changes no Go,
PostgreSQL, AWS SDK, fiscal, security-scanner or approval rule. It is not source
code from xAI, Microsoft or AWS. The bounded reader is ADAPTED from the existing
src/platform/backend/bounded-command.ts pattern, with the original document
owner's2MiB binary budget and a finite read deadline; no new dependency is added.

## Connected operation

The complete120profile applies UI and shared backend transport overlays. The
DocumentWorkspaceV403 component and /experience/documents entry use the fixed
/api/enterprise/documents-v403 BFF. It binds the session tenant/organization,
configured profile and actual documents:write/process/review permissions.
The backend still enforces object ownership, permissions and separate reviewer.
The web process must receive the same DOCUMENTS_ENABLED, DOCUMENTS_MODE,
DOCUMENTS_TENANT_ID, DOCUMENTS_ORGANIZATION_ID and DOCUMENTS_PROFILE_SHA256
configuration as its chosen backend. These are configuration, not credentials.
Missing/incompatible configuration is never renamed a missing account.

Receive exact original bytes under one generated UUID and SHA; process explicitly;
show/correct invoice_number, vendor, total and currency; submit the immutable
proposal; hand its scoped URL to another authorized reviewer; decide on the exact
payload SHA. The decision endpoint owns the atomic commit. There is no client
auto-approval, extra commit request, accounting posting or cross-tenant directory.
The existing backend exposes a document by ID, so this adapter supplies an opaque
review link. It does not claim a new shared inbox, assignment service or document
search endpoint. Those require their own product scope/owner when selected.

Amounts and leading zeroes remain text. Original filenames travel as reversible
percent-encoded ASCII within the existing128-byte name contract, and are decoded
for presentation. Bytes are never converted to change MIME, page count or size.
The UI accepts JPEG/PDF up to2MiB; the existing processor accepts one page.
FIXTURE mode only accepts the owner's exact public JPEG; browser transport fixtures
use synthetic bytes and cannot prove that provider or security admission.
Original/evidence downloads are hash-checked, no-store, nosniff attachments.
No untrusted PDF/HTML renderer or camera decoder is introduced.

## Continuity and errors

Recovery storage contains only scoped opaque document IDs, operation, hashes and
decision boolean. File bytes, extracted fields, names, reasons and tokens are not
written there. The active ID survives reload; an uncertain operation retains its
pending marker. Consult recorded state before continuing; matching receipt/state
clears the marker, never a timeout. A pre-hash guard prevents duplicate submissions.
Refreshing a review preserves edited fields for the same unsubmitted document.
Terminal rejection/quarantine preserves the historical record; a new document can
be opened explicitly. No automatic retries of an uncertain write are implemented.

## Verification and remaining methods

Focused contract/gateway tests use synthetic HTTP responses and the existing
Vitest/TypeScript toolchain; browser tests use the real Next BFF with an explicitly
simulated backend. They cover bounded bytes, scoped permissions, schema/hash
binding, receive/process/review/decision, uncertainty/reload and duplicate prevention.
Exact source hashes and executed counts belong to the finalV403receipt, not this
mutable paragraph. Prior Go/PostgreSQL fixture evidence is reused only for unchanged
owner bytes; these frontend tests do not re-prove database or provider operation.

Real OCR accuracy, scanners, corpus, provider accounts, fiscal/legal policies,
physical phone/camera, assistive-technology and human acceptance remain separate.
Read docs/DOCUMENT_REFERENCE_DECISION.md and the library's
OFFICIAL_DOCUMENT_INTELLIGENCE_PROFILE, SDK_ARTIFACT_LOCK and FIXTURE_CATALOG
before changing class/model/format or admitting automatic storage. Ground truth
is not manufactured from a sample or a simulated extraction response.
