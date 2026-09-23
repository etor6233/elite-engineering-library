# Google Merchant Product Sync evidence adapter

This module invokes Google's official GA `google-shopping-merchant-products` Python client `1.8.0`, `ProductInputsServiceClient.insert_product_input`, and `ProductsServiceClient.get_product`. The SDK and generated v1 models are Google code under Apache-2.0. `sync_product.py`, profiles, validation, and receipts are local integration glue and are not represented as Google-authored code.

The provided profile is blocked. Before use, prove the Google Cloud project and billing decision, Merchant API enablement, ADC identity, Merchant account access, an API primary data source, a test-account contract, allowed offer IDs, quota/cost, retention, reconciliation ownership, and explicit authorization for the external product write. Credentials remain in Application Default Credentials or an authorized runtime secret store; never place tokens or keys in Markdown, JSON, CLI arguments, logs, or receipts.

For the verified CPython 3.14 / Windows x86-64 lane:

```text
python -m venv .venv
.venv\Scripts\python -m pip install --require-hashes -r requirements-windows-py314.lock
.venv\Scripts\python -m pip check
.venv\Scripts\python -m unittest -v test_sync_product.py
.venv\Scripts\python sync_product.py --profile C:\secure\merchant-profile.json --product C:\approved\product.json --output C:\evidence\merchant-SKU12345
```

For another OS/Python ABI, treat `requirements-direct.in` as the exact top-level input and produce a new fully hash-pinned candidate lock under the dependency update contract before execution. Do not reuse the Windows lock on another target.

The adapter stores the exact insert request/response and then reads the processed Product status. A successful insert is never labelled approved: Google states that insertion starts processing and approval must be read from `product_status`. If the product is not yet visible, the receipt remains pending; if the post-write status check fails, durable insert evidence is retained before the error is raised. No local business record is written automatically.
