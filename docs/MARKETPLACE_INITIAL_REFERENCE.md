# Approved catalog to initial Mercado Libre publication

Scope: local/fixtures, User Products buy_it_now, one selected physical variant,
one immutable normalized PNG, and initial item stock mode. Existing location
mutation modes retain their separate contract. No classified vehicle, Full-stock
write, automatic price rule or whole-provider certification is implied.

## Run the connected owner

Apply migrations through0080 and enable the existing catalog and marketplace
modules in the same Go host. Use the existing MARKETPLACE_* profile/hash/HMAC
settings and future MERCADOLIBRE_ACCESS_TOKEN. The reference profile selects
item stock mode for initial publication. Provider category, listing type,
condition/attributes and seller account must match the target product; provider
validation rejects incompatible settings without choosing commercial policy.

The existing role API is /v1/admin/marketplace. Every request is authenticated,
tenant/organization scoped, private/no-store and bounded. No secrets enter JSON.

1. POST /prepare with operation MEDIA, a new approval_id, current generation,
   variant_id and expires_at (within15minutes). No item_id or media_approval_id.
2. POST the returned exact intent to /requests; a distinct authorized human
   approves its request_sha256 at /requests/{id}/decision.
3. POST {} to /requests/{id}/send. The original durable fence admits one
   multipart POST. Only the current catalog owner's normalized PNG bytes,
   verified against their SHA256, can be uploaded. GET request status after
   any lost API acknowledgement; do not generate another identifier to retry.
4. Prepare operation CREATE with a new approval_id and the accepted
   media_approval_id, same generation/variant. No item_id. The server binds the
   uploaded provider picture ID, current source price/attributes/family name and
   original serial ATP. Submit, obtain distinct human approval and send.
5. POST {} to /requests/{id}/reconcile when the original send is unknown.
   Creation recovery uses a persisted provider ID or bounded seller-SKU search,
   then original item/prices/stock GETs. It compares seller/site/category/SKU,
   User Product, family, listing type, mode, attributes, exact decimal amount,
   quantity and the accepted picture ID. It never issues another item POST.

available_to_promise keeps its original owner and horizon semantics; it can
include inbound supply. Publication neither reserves stock nor invents physical
receipt. An initial quantity zero can produce a paused listing. A confirmed
publication is not a claim that Mercado Libre approved or activated the listing.

## Durable ambiguity and recovery

Each media upload and creation has a separate immutable human request, original
outbound fence and immutable catalog effect. Parsed provider upload/create
acknowledgements persist before subsequent confirmation. A lost local API reply
is recovered by GET. A lost item POST reply is recovered by exact readback.
Existing accepted CREATE evidence prevents another local creation of that SKU,
even if provider search later returns empty. A replacement selling condition or
deleted listing requires explicit target policy and a new admission; it is not
an automatic duplicate.

An upload whose provider response is completely lost and whose picture ID was
never durably observed remains UNKNOWN. There is no documented source-hash
lookup contract and no invented image equivalence. That original operation
cannot resend or be marked accepted from a guessed ID. A separately reviewed
MEDIA request may upload a replacement PNG; it does not clear the unknown
orphan. Unassociated provider media has no product publication effect.

Mercado Libre transforms images. Evidence binds the exact catalog PNG hash to
the multipart request, authenticated provider response/picture ID and listing
readback. It does not claim that provider JPEG bytes equal the original PNG.

## Provenance and verification

All new HTTP/JSON/PG/host/test code is AUTHORED composition, with no corporate
authorship claim. Existing catalog normalization, amounts, serial ATP, manual
approval and fence algorithms retain their original owners and notices.
Official contracts are fixed in docs/marketplace/official-contract.lock.json.
No archived Mercado Libre SDK or unpinned source code is included.

Run TestMarketplaceInitialContracts for seven transport/recovery scenarios;
TestMarketplaceInitialConnected against an owned database with all migrations
for source→API→PG→provider effects,12callers→1upload, human separation, tampering,
unknown/no-retry, accepted upload binding, duplicate CREATE denial and immutable
evidence. Host guards require both effect and observation triggers. Populated
0080down fails and rolls back; empty down/up passes. The exact command-boundary
fuzz campaign includes MEDIA/CREATE plus the prior mutation invariants.

This closes initial publication/media in local fixtures. Existing-item CONTENT
updates are connected in MARKETPLACE_CONTENT_REFERENCE.md. Remaining feeds/comms
and the global TEST02/03/07/T2805 closure remain separately tracked. Provider credentials are supplied later by the user.
