# Connected Merchant contract decisions

Contract decisions for the connected local claim; final receipts and remaining T2805 scope are in MERCHANT_CONNECTED_REFERENCE.md.

Reuse PYTHON_GOOGLE_MERCHANT_PRODUCT_SYNC_ADAPTER1.8.0 SDK pin, its exact
build_insert_request and official ProductInput/ProductAttributes/Price models.
Reuse the admitted catalog source, serial ATP, human approval and outbound fence.
The Go mapping, bounded Python IPC, PostgreSQL binding and fixture transport are
AUTHORED glue. No copied tutorial or new pricing/inventory algorithm.

Selected reference: one offer in a dedicated API primary data source. Offer,
language/feed label, account/data source and product metadata are explicitly
profile/source-bound. Source model specification must contain reviewed
merchant_description; never manufacture a description, brand or GTIN.
The product link and immutable PNG URL use the current catalog origin/routes.
Price minor units convert exactly to integer micros; overflow or unrepresentable
currency precision is rejected. Availability follows original current serial ATP
with a recorded horizon, without claiming physical receipt or reserving stock.

Use official ProductInput.version_number for primary-source insertion and compare
Product.version_number when observing the processed product. Bind it to source
generation for this dedicated reference namespace; lower versions are rejected
by the provider. Refresh at the same generation is supported by the documented
contract. Existing foreign version namespaces require target admission; no
automatic source takeover or shared-account migration is implied.

One manual approval binds the exact provider request, source generation, profile,
offer and ATP. One provider attempt is admitted by the original durable fence.
Persist the exact SDK response before subsequent processing observation. Missing
acknowledgements remain UNKNOWN; recovery is GET-only by the deterministic
account/language/feed-label/offer resource name. Confirm identity/data source,
generation and selected fields. NotFound is pending, never proof of no effect.

An accepted ProductInput is not Google approval. Processed attributes can differ
because of rules/supplemental sources. Preserve product_status/issues separately,
report mismatch, and never mark ads/product approval from an insert200.
An exact GET match after an ambiguous insert proves observed desired state, not
that the provider's30day freshness deadline was refreshed by that write.

The explicit refresh-due/manual approval queue and operator reconciliation path are connected. Queue reads never re-send expired manual requests. Automated communications remain the separate T2805 job/notification scope.

Source acquisition: merchant-connected-authority/acquisition.json reacquires the
same20official wheels with exact hashes and fixes three current official docs;
runtime-receipt.json restores the pinned Windows/CPython3.14 environment offline,
passes pip check and loads actual SDK1.8.0 insert/get/product_status interfaces.
No accounts, billing action or live product write occurred.
