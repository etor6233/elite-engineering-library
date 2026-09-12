# Amazon SP-API Catalog Adapter — reconstruction evidence V1

## Scope and authority

Audit date: 2026-08-26. Primary authority is Amazon's public `amzn/selling-partner-api-sdk`, release/tag `python-1.11.1`, commit `8e792ae345a8d334ccdbdd03181f05f040e6a4fc`. GitHub reported the commit signature verified and the repository active under Apache-2.0. Amazon's official PyPI distribution is `amzn-sp-api 1.11.1`.

This evidence covers source/artifact identity, local reconstruction, dependency installation and contract tests. It does not cover developer registration, application approval, roles, LWA credentials, account data, a sandbox call, rate limits or production.

## Exact authority receipts

| Artifact | Exact receipt |
|---|---|
| source archive | 26,447,333 bytes; SHA-256 `5bf712836bda76619c796fdb2b3a0a13acb281ba8bcf1b39c8f4b92f8f4799a1` |
| repository LICENSE | SHA-256 `09e8a9bcec8067104652c168685ab0931e7868f9c8284b66f5ae6edae5f1130b` |
| repository NOTICE | SHA-256 `d4290ed64c2edd0fce1d84e3f9dfb2881240fe534def76b8cd29ed6af683e287` |
| PyPI wheel | `amzn_sp_api-1.11.1-py3-none-any.whl`; 3,557,108 bytes; SHA-256 `a88e4059c9fbfd454bfef67b82288a926f73cbdc608a8831f8383565736c70b0` |
| PyPI sdist | `amzn_sp_api-1.11.1.tar.gz`; 1,056,178 bytes; SHA-256 `f76a8f76a357e69e0a9ba2d2ba81143ef59e1aa7bb12a9f04743c6aeaac1a3e6` |

The release source contains a stale distribution `pyproject.toml` version `1.11.0`, while the trusted PyPI wheel metadata reports `Name: amzn-sp-api`, `Version: 1.11.1`, `License-Expression: Apache-2.0`, `Requires-Python: >=3.9`. This contradiction is retained as `UP-FAIL-020`; the pack verifies installed metadata 1.11.1 and never derives runtime identity from the stale file.

## Materialized implementation

`PYTHON-AMAZON-SPAPI-CATALOG-ADAPTER 0.1.0` reconstructs seven files. The complete graph is hash-fixed to Amazon SDK 1.11.1 plus requests 2.34.2, urllib3 2.7.0, certifi 2026.7.22, charset-normalizer 3.5.1, idna 3.19 and six 1.17.0. The adapter invokes official `SPAPIConfig`, `SPAPIClient`, `CatalogApi` and Catalog Items `2022-04-01`; local validation/evidence code remains `AUTHORED`.

The profile blocks execution until developer/application registration, roles/data access, sandbox contract, region/marketplace, ASIN/datasets, quota/cost and reconciliation are proven. Secrets are referenced only by three canonical environment variable names. One `get_catalog_item` result is preserved with an atomic receipt; ASIN and marketplace are hashed in the receipt and no local business write occurs.

## Verification and failures retained

- complete seven-wheel installation finished and `pip check` reported no broken requirements;
- official distribution version, `CatalogApi.get_catalog_item(asin, marketplace_ids, ...)`, `SPAPIConfig(refresh_token=...)` and Catalog API model namespace passed;
- five unit/contract/negative tests passed in 0.013 seconds;
- the canonical Markdown reconstructed seven files with zero content hash differences;
- `AMAZON_SPAPI_CATALOG_PACK_PLAN.md` composed 19 implementation files plus its materialization record;
- `LIB-FAIL-044` records the rejected wrong materializer parameter;
- `LIB-FAIL-045` records two stuck install processes after files were installed; only their verified audit command lines were terminated, and metadata/tests/checks were rerun in separate processes.

## Admission decision

`REBUILD_VERIFIED / CONDITIONED`. This is real runnable integration code over Amazon's official SDK, not a production-proven Amazon integration. Promotion requires a real sandbox/account, role mapping, LWA lifecycle, contract fixtures, rate/retry policy (the Python SDK advertises neither RDT nor rate limiter), reconciliation, privacy/security, load and rollback evidence.
