# Google Merchant Product Sync Adapter — reconstruction evidence V1

## Scope and authority

Audit date: 2026-08-26. Primary authority is Google's `googleapis/google-cloud-python`, release `google-shopping-merchant-products-v1.8.0`, verified commit `97d7b42cd74b41211f5ec8871cc0dd15debdb1a0`, plus the trusted PyPI distribution `google-shopping-merchant-products 1.8.0`. Google classifies this client `Development Status :: 5 - Production/Stable`, Python >=3.10, Apache-2.0.

Google's current product guide states that product writes require an applicable data source, a successful `productInputs.insert` starts processing but does not prove approval, processed values can differ from input, status is read from `Product.product_status`, and products must be refreshed at least every 30 days. Those constraints govern the pack.

## Exact authority receipts

| Artifact | Exact receipt |
|---|---|
| monorepo source archive | 206,581,382 bytes; SHA-256 `cde3981299462b8346f725d6fd4bf7973972b6b4263cfc7399081e0514ad6b1c` |
| root/package LICENSE | SHA-256 `cfc7749b96f63bd31c3c42b5c471bf756814053e847c10f3eb003417bc523d30` |
| PyPI wheel | 244,077 bytes; SHA-256 `722ef095eca35126255c129964a61818f0d7bc9eeef7f1da12732d81c1dcd505` |
| PyPI sdist | 241,375 bytes; SHA-256 `08f0ab5ccbdf4df65cdf8c0d763ba161acc02f61897f6bf023a20223acbbb7f9` |
| sdist safety | 104 entries; zero absolute/traversal paths; PKG-INFO version 1.8.0; LICENSE hash identical |

## Platform lock and implementation

The verified immediate lane is CPython 3.14 on Windows x86-64. Its `requirements-windows-py314.lock` contains 20 exact wheel URLs/hashes: Google Merchant Products/Shopping Type/API Core/Auth/common protos, grpc/protobuf/proto-plus, requests/trust dependencies and authentication cryptography dependencies. A clean isolated installation passed `pip check`. Other OS/ABI combinations must generate and approve a target-specific lock from the exact top-level wheel input; the Windows lock is fail-closed outside its lane.

`PYTHON-GOOGLE-MERCHANT-PRODUCT-SYNC-ADAPTER 0.1.0` reconstructs eight files. It validates a bounded product schema, builds official v1 `ProductInput`, `ProductAttributes`, `Price` and `InsertProductInputRequest`, disables implicit retry on the external write, persists insert evidence, then calls official `ProductsServiceClient.get_product`. `NotFound` becomes a pending state, not approval. An unexpected post-write status failure preserves durable insert evidence before raising. No local domain write occurs.

## Verification and failure retained

- official SDK version 1.8.0, both client signatures and `Product.product_status` passed runtime probes;
- seven unit/contract/negative/durability tests passed in both authoring and clean reconstructed trees;
- exact Markdown reconstruction produced eight files with zero content hash differences;
- `GOOGLE_MERCHANT_PRODUCT_SYNC_PACK_PLAN.md` composed 20 implementation files plus record;
- `LIB-FAIL-046` records that full monorepo expansion exceeded the command window and was stopped; that partial tree was not used as package evidence. The small official sdist and wheel were independently hashed, inspected and installed.

## Admission decision

`REBUILD_VERIFIED / CONDITIONED`. The code is immediately materializable and invokes Google's official GA client, but no Merchant account or external write was exercised. Promotion requires project/billing/API/ADC, Merchant access, API primary data source, test account, approved offers/product data, product-status/issues reconciliation, refresh cadence, quota/cost, concurrency/idempotency, privacy/security, observability, compensation and production pilot evidence.
