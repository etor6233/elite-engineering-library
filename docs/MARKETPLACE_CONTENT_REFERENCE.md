# Existing item content from the approved catalog

CONTENT extends the same Mercado Libre owner, Go host, role API, approval kind
and outbound fence. It sends one PUT/items/{id} containing exactly family_name,
attributes and pictures. It does not include price, quantity, status or a legacy
title/variations field. Family name comes from the current immutable catalog
model; attributes are the explicitly bound provider profile; the picture ID
comes from an accepted MEDIA upload for the same generation and variant.

Use the workflow in MARKETPLACE_INITIAL_REFERENCE.md to upload the source PNG.
Then POST /v1/admin/marketplace/prepare with operation CONTENT, a new approval_id,
accepted media_approval_id, current generation, variant_id, existing item_id and
expires_at within15minutes. Submit the returned exact intent, obtain approval
from a distinct authorized reviewer, and send through the existing request path.
GET status after a missing API response. Reconciliation is read-only.

Official User Products documentation states that family_name propagates to all
selling conditions of a User Product, and cannot be changed after an associated
selling condition has sales. This reference therefore requires a bounded
authenticated search returning exactly the selected item and an explicit
sold_quantity zero. Missing sales metadata, more than one associated item,
unexpected owner/UP or changed preconditions fail closed before PUT. No inferred
permission to rename other listings is granted.

Readback confirms the same owner/item/UP and exact requested name, attributes
and picture ID. Provider image bytes remain transformed, not SHA-equivalent to
the normalized PNG. Concurrent changes outside this application cannot be
locked by a local transaction; the receipt identifies the observed scope.

Migration0081 extends the immutable effect's operation constraint. Populated
rollback refuses to erase CONTENT evidence; empty down/up is supported.
Host activation requires that constraint plus the existing immutable guards.

TestMarketplaceContentContracts checks eight contract scenarios. The connected
source/API/PG proof uses3approvals/3fences/3attempts,1multipart,1contentPUT,
lost-response recovery and12replays without another effect. Unsafe scope,
client title tampering and tenant/org/permission boundaries are rejected.
Price and stock remain byte-semantically unchanged by the three-field write.

All new implementation/proof is AUTHORED glue around the existing catalog,
media, human approval, ATP, decimal and fence owners. Eight fixed official
contracts remain in docs/marketplace/official-contract.lock.json. No dependency,
SDK or upstream algorithm changed. Other provider feeds and communications
automation remain separate T2805 work.
