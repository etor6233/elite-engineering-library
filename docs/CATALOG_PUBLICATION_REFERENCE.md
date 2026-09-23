# Connected catalog publication — J3

LIBRARY_INFRASTRUCTURE / local reference. AUTHORED orchestration reuses the
existing Catalog, Commerce, BC eligibility, shared Approval, Search and
OutboundDelivery owners. No new price, tax, legal, ranking or business policy.

Apply migrations through0074 in order. Select the complete franchise composition.
CATALOG_RELEASE_ENABLED=true requires CATALOG_RELEASE_PROFILE_FILE (absolute
regular file) and CATALOG_RELEASE_PROFILE_SHA256 (SHA256 of exact UTF-8 bytes).
deploy/catalog/publication.reference.json is an explicit synthetic profile.
Choose an existing active tenant/organization, market/currency and HTTPS origin;
never infer fiscal settings. Host requires the same selected Commerce repository
and12enabled database guards. Disabled mode reads no profile or provider keys.

Draft creation captures existing validated models/variants, a bounded existing
price book and normalized PNG assets.1–32models/1–128variants; snapshot at most
256KiB. PNG input at most1MiB, dimensions at most2048each/1048576totalpixels,
normalized output at most4MiB. Go1.26.8 image/png supplies decoding/reencoding.
This removes uninterpreted trailers; it is not antivirus or libxml2 assurance.

Four shared human review requests bind the exact draft SHA: legal, technical,
media and publication. Each reviewer differs from the maker; the final stage
requires three preceding approvals. These are recorded decisions, not automatic
legal/technical certification. Publish requires all four approvals and current
database-clock price validity. Generic approval bypass is rejected.

Each publication/rollback copies the exact approved price header and entries
into a fresh effective book via the original Commerce transaction writer, then
activates it once. This preserves created1/activated2 outbox semantics, old quote
references and immutable history. One transaction publishes effective pricing,
search projection, receipt and outbox. Rollback appends a generation and refuses
expired prices. Source model/price edits never mutate a captured snapshot.

Private routes require a single organization_id query and verifier-supplied
tenant/actor. Bodies cannot choose identity. POST /v1/admin/catalog/media/{command}
uses image/png and catalog:draft; POST /drafts and GET /drafts/{id} use catalog:draft
and catalog:read respectively. POST /drafts/{id}/review/{stage} uses
catalog:review:{stage}; POST /drafts/{id}/publish uses catalog:publish.
GET /commands/{command} recovers actor/request-hash-bound receipts after a lost
response. Do not issue another command ID merely because the transport failed.
Versions and minor amounts are quoted int64 decimal strings. JSON32KiB boundary,
duplicate/unknown fields and wrong organization/permissions are rejected.

GET /v1/public/catalog and /feed return identical typed current publication bytes.
They omit private profile/reviewer configuration, include the public tenant code,
and emit ETag generation/sourceSHA with max-age=0,must-revalidate. Current-only
normalized media and scoped PG full-text search use the same approved source.
The existing public model and price owners read this publication when selected;
unconfigured tenants and narrower Go profiles preserve their existing behavior.

Optional reference feed: CATALOG_FEED_ENABLED=true additionally requires exact
CATALOG_FEED_PROFILE_FILE/SHA256, CATALOG_FEED_HMAC_KEY (base6432bytes) and
CATALOG_FEED_TOKEN. The fixture profile uses literal loopback HTTP; other profiles
require HTTPS. No provider call occurs during activation. Real values are supplied
later from target secret management, never in the profile/receipt. Enabling feed
without these inputs or its immutable database guard fails closed.

POST /v1/admin/catalog/feed/{generation} and POST the same path/reconcile require
catalog:feed; GET uses catalog:read. The existing OutboundDelivery state machine
and fence own the send. An immutable intent binds actor/generation/source/payload/
receiver profile hashes in the same claim transaction. Receiver protocol:
POST and GET /publications/{idempotency-key}; Authorization Bearer,
Idempotency-Key, X-Catalog-Generation and X-Content-SHA256 headers. ACK JSON binds
status accepted/rejected, quoted generation, payload_sha256, remote_id and
accepted_at. Accepted200/201 must match the request; terminal rejected422 remains
rejected. An unknown send is reconciled by GET, never automatically re-POSTed.
The receiver preserves per-key receipts and rejects older generations.
This is an explicit reference receiver contract. Mapping to the actual admitted
Google/Meta/ML SDKs belongs to T2805, not a claim of credential-only readiness yet.

Actual PG18.6/71migrations, Go1.26.8 and Next16.3.4/Node24.20.0/Chromium:
3snapshots/12reviews/5publications;1new/7replays;13table rollback; source isolation;
lost publication response recovered by GET;3stale ETags; four reference feed sends,
one terminal rejection and one lost accepted response reconciled once. Actual
storefront follows generations1/2/3/5 with canonical/sitemap/robots/PNG/404.
73/74down refuses populated history intact; empty74→73down then73→74up passes.
Fuzz finite2s/7seeds/121059executions. Narrow compile is not a runtime claim.
Canonical CATALOG_CONNECTED_RELEASE_V402.md/json binds exact source and receipts.
Role editor UI T2804, provider mappings T2805, identity/SCA T2803, operations and
signed release remain their own gates. No production or corporate glue attribution.
