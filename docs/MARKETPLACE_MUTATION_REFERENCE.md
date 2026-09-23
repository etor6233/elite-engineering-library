# Mercado Libre: approved catalog mutations

This pack supplies AUTHORED HTTP, SQL and host composition around the existing
catalog, inventory ATP, human approval and outbound-delivery owners. It contains
no Mercado Libre SDK source. Exact official documentation snapshots are listed
in docs/marketplace/official-contract.lock.json. The archived Go SDK is not used.

The closed local claim is an existing, seller-owned User Products item:
PRICE, STOCK, PAUSE and RESUME. MEDIA upload and CREATE initial publication
are connected in MARKETPLACE_INITIAL_REFERENCE.md; existing-item CONTENT
updates are connected in MARKETPLACE_CONTENT_REFERENCE.md. Google Merchant
mapping and remaining T2805 orchestration remain tracked in the library gap map. These endpoints do not
silently approximate those operations.

Profile fixes tenant/organization, seller/site/currency and each variant's SKU,
category and inventory location mode. The reference uses explicit synthetic
values. Listing metadata in the intent records the approved source context;
these four operations never change title, attributes or pictures.

The price is the exact current published minor-unit amount, serialized using
the existing decimal representation owner. No floats, currency conversion,
tax calculation or pricing policy is introduced. The stock quantity comes
from the existing inventory.serial_atp function with a captured horizon.
This is its ATP policy, not a new promise that inbound stock is physically
received. Neither provider write nor reconciliation reserves local inventory.
Live channel allocation and provider orders retain their existing owners.

Supported stock modes are the normal item quantity, Full/Flex selling-address
quantity and one explicitly bound seller warehouse. Full inventory is never
written. The adapter reads item ownership, standard prices and stock locations
before preparing a request. Location writes retain the returned x-version.
Ambiguous/multiple standard prices and unexpected location topology fail closed.
Dynamic pricing blocks PRICE before the PUT. The not-yet-available
/prices/standard write contract is not selected.

Workflow through the same Go application and identity verifier:

1. POST /v1/admin/marketplace/prepare with approval_id (16–64 characters),
   generation (decimal string), variant_id, operation, item_id and expires_at.
   This only reads the provider and the current catalog/ATP source.
2. POST the returned exact intent to /v1/admin/marketplace/requests.
3. A distinct reviewer POSTs request_sha256, approve and reason to
   /v1/admin/marketplace/requests/{id}/decision.
4. POST {} to /v1/admin/marketplace/requests/{id}/send.
5. GET /v1/admin/marketplace/requests/{id} after any missing response.
   If unknown, POST {} to /reconcile, which only reads provider state.

Each permission is explicit: marketplace:request, marketplace:approve,
marketplace:send, marketplace:read and marketplace:reconcile. The original
global administrator permission retains its documented organization semantics.
Unprivileged users are rejected before reading the body. JSON is bounded to
32 KiB with a five-second read deadline; duplicate/extra fields fail. No private
response is cacheable. Preparation expires within fifteen minutes.

The shared PostgreSQL fence permits one provider attempt per approved id.
Another unresolved command for the seller/SKU cannot create a second effect.
Lost acknowledgements, expired sender leases and readback mismatch remain
unknown; no automatic PUT retry exists. A documented rejection or a changed
precondition becomes terminal and requires a new reviewed command.
HTTP200 alone does not mean the requested price/stock/state was applied.
Reconciliation confirms observed state, not a distributed transaction or
exclusive ownership of concurrent edits made outside this application.
Catalog changes after claim may require a subsequent approved synchronization;
the receipt always identifies its original generation and source hash.

Activation uses the existing CATALOG_RELEASE_* settings plus:

    MARKETPLACE_ENABLED=false
    MARKETPLACE_PROFILE_FILE=
    MARKETPLACE_PROFILE_SHA256=
    MARKETPLACE_HMAC_KEY=
    MERCADOLIBRE_ACCESS_TOKEN=

When enabled, the application requires the exact profile/hash, catalog
activation, HMAC material and token. Missing migration/immutable approval
guards prevent startup. The HTTP origin is fixed to api.mercadolibre.com,
TLS >=1.2, no redirects or environment proxy, bounded response/timeouts.
Provider token renewal is supplied by the existing activation environment
owner; credentials never enter browser requests, receipts or logs.

The local reference uses a loopback fixture transport only in tests, synthetic
principals and explicit human decisions. Tests cover exact large amounts,
three stock modes, provider/API acknowledgement loss, twelve concurrent sends,
JSONB canonicalization, cross-product/seller binding, ignored price, version
drift, immutable evidence, startup guard and populated/empty migration rollback.
This is library infrastructure evidence, not live provider or production approval.
