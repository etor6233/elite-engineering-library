# Amazon SP-API Shipping Tracking Execution Inventory — V147

Date: 2026-08-31  
Scope: read-only Amazon Shipping V2 tracking lane.

## Result

`PYTHON-AMAZON-SPAPI-SHIPPING-TRACKING-ADAPTER@0.1.0` materializes seven files. Its profile composes 2 packs / 32 files: the exact upstream acquisition lock plus the adapter. The wrapper invokes Amazon's official `ShippingApi.get_tracking(tracking_id, carrier_id)`, atomically preserves the provider response, hashes the tracking ID in its receipt and fixes `automatic_delivery_completion=false`.

## Authority and provenance

- official SDK: `amzn/selling-partner-api-sdk`, release `python-1.11.1`, commit `8e792ae345a8d334ccdbdd03181f05f040e6a4fc`, Apache-2.0;
- official model repository: `amzn/selling-partner-api-models`, commit `8e429486005c4ebdce5099e48cc48515a65359bb`, Apache-2.0;
- official wheel: `amzn-sp-api==1.11.1`, SHA-256 `a88e4059c9fbfd454bfef67b82288a926f73cbdc608a8831f8383565736c70b0`;
- official repository: https://github.com/amzn/selling-partner-api-sdk
- official models: https://github.com/amzn/selling-partner-api-models

The SDK/model code is Amazon-authored. The seven materialized wrapper/profile/lock/test/runbook files are local `AUTHORED` integration and are not represented as Amazon source.

## Executed gates

- official installed SDK reports Shipping V2 and exact `get_tracking(self, tracking_id, carrier_id, **kwargs)` signature;
- `pip check` PASS for the existing exact seven-wheel environment;
- six tests PASS: official signature, atomic positive output, invalid input/output rejection, provider-failure cleanup, profile closure and SDK-version rejection;
- pack materialization 7/7;
- all seven reconstructed files match the tested author tree byte-for-byte;
- clean profile composition 2 packs / 32 files;
- global `VERIFY_LIBRARY_PASS` and `VERIFY_EXECUTABLE_LIBRARY_PASS` after V147 promotion: 89 packs / 955 blocks / 493 Markdown / 37 composed profiles / 121 locked sources.

## Boundary retained

No Amazon call, account, sandbox, region, carrier support, provider response, header/rate behavior or live reconciliation was claimed. This is not a universal carrier, quote/purchase/cancel adapter, webhook consumer, proof of physical delivery or production evidence. A project must prove real credentials, Shipping role/scope, sandbox, geography/carrier, privacy/retention, limits/retries, duplicates/out-of-order events, reconciliation, load, alarms, rollback and business acceptance. No project is production-ready until its target CDN/WAF, IdP, providers, PostgreSQL recovery, load, offensive security, deployment and acceptance are demonstrated.
