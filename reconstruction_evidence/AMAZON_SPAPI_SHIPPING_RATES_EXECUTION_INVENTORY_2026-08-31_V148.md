# Amazon SP-API Shipping Rates Execution Inventory — V148

Date: 2026-08-31  
Scope: strict one-package Shipping V2 rate quotation added to the V147 tracking lane.

## Result

`PYTHON-AMAZON-SPAPI-SHIPPING-TRACKING-ADAPTER@0.2.0` now materializes ten files and its profile composes 2 packs / 35 files. Rates constructs Amazon's official `Address`, `Dimensions`, `Weight`, `Currency`, `Item`, `Package`, `ChannelDetails` and `GetRatesRequest` models and invokes the exact `ShippingApi.get_rates(body)` method. The receipt commits atomically, hashes the canonical request and fixes both `shipment_purchased=false` and `automatic_delivery_completion=false`.

## Authority and provenance

- official SDK `amzn-sp-api` 1.11.1 / commit `8e792ae345a8d334ccdbdd03181f05f040e6a4fc` / Apache-2.0: https://github.com/amzn/selling-partner-api-sdk
- official models commit `8e429486005c4ebdce5099e48cc48515a65359bb` / Apache-2.0: https://github.com/amzn/selling-partner-api-models
- Amazon's official samples state that samples are educational and require project hardening: https://github.com/amzn/selling-partner-api-samples
- open Shipping V2 issue #5270 reports `getRates` rejecting more than one package despite no model `maxItems`: https://github.com/amzn/selling-partner-api-models/issues/5270

The SDK/models are Amazon-authored. The ten materialized files are local `AUTHORED` integration and are not represented as Amazon source. The one-package limit is a conservative fail-closed response to current upstream evidence, not an invented claim that the API universally supports only one package.

## Executed gates

- `pip check` PASS on the exact seven-wheel environment;
- 11 tests PASS: six tracking plus five rates;
- official model construction proves required insured value, items and external channel; malformed keys, units, values, dates/countries, provider failure and existing output reject closed;
- pack materialization 10/10;
- all ten reconstructed files match the tested author tree byte-for-byte;
- clean profile composition 2 packs / 35 files.
- global `VERIFY_LIBRARY_PASS` and `VERIFY_EXECUTABLE_LIBRARY_PASS`: 89 packs / 958 blocks / 494 Markdown / 37 profiles / 121 locked sources.

## Boundary retained

No live Amazon call, rate, account, sandbox, carrier/geography, headers, quota or price was demonstrated. Python's official SDK has basic API support but no built-in rate limiter. V148 does not purchase/cancel a shipment, fetch documents, accept a handover, support multiple packages or prove delivery. Purchase must bind the exact provider `requestToken`, selected `rateId` and returned document specification before it can be admitted. Production still requires real access, privacy, retry/throttle evidence, reconciliation, load, alarms, rollback and business acceptance.
