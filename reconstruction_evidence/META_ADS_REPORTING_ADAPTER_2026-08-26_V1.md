# Meta Ads Reporting Adapter — reconstruction evidence V1

## Scope and authority

Audit date: 2026-08-26. Primary authority is Meta's active `facebook/facebook-python-business-sdk`, release/tag `26.0.1`, commit `788f363d15b1269ab5efb7cd00fb5e3b133cd99b`, and PyPI distribution `facebook-business 26.0.1` owned by Meta Platforms. The official `AdAccount.get_insights(fields, params, ...)` implementation creates a GET request to `/insights` for Graph API v26.0.

The SDK and its `facebook/capi-param-builder` dependency use the Meta platform license: use, copy, modification and distribution are allowed only in connection with Facebook web services/APIs and remain subject to Platform Policy. This is not an OSI/general-purpose license. Local adapter code is explicitly `AUTHORED`; it is not attributed to Meta.

## Exact authority receipts

| Artifact | Exact receipt |
|---|---|
| Business SDK source archive | 2,070,568 bytes; SHA-256 `bd3bc1f14072662eecb05fb6f43888eb652c9a1f039894f99a7e7609e4452db7` |
| Business SDK wheel | 1,554,832 bytes; SHA-256 `41555ce87ac105617ed678a489a630108ac162fbdccfbbf65fa20ae1de7cd0ca` |
| Business SDK sdist | 732,948 bytes; SHA-256 `2397aaa1da3f070ae18b6477697682ecfac9cd308521821f08e21b91cc1e4d56` |
| CAPI source | tag `v1.3.0-python`, commit `66a7eb84e29a65aa999089366c1f7c3c0e49bdb7`; 674,665 bytes; SHA-256 `be8d21efab5a1ea16b1f6e84d4219d9e343436b8789ede407d0b94d0a54dabf5` |
| Meta LICENSE | 1,035 bytes; SHA-256 `48d97b3c936203a750a3288c2b924327769d62327383a1060244b3d3c05ad01f` in both sources |
| signatures | GitHub API returned `verified=false`, reason `unsigned`, for both exact tagged commits |

## Platform lock and implementation

The verified immediate lane is CPython 3.14 Windows x86-64. Its lock contains 18 exact wheel URLs/hashes: Meta Business/CAPI plus requests, aiohttp and transitives. Clean `--require-hashes` installation and `pip check` passed. The CAPI 1.3.0 wheel omits its LICENSE; `UP-FAIL-021` requires preserving the exact official source LICENSE. Other OS/ABI targets require an independently admitted hash-complete lock.

`PYTHON-META-ADS-REPORTING-ADAPTER 0.1.0` reconstructs eight files. It rejects unproven app/business/`ads_read`/test-account/cost/reconciliation state, account IDs outside numeric form, non-approved fields, unknown levels, invalid or >31-day ranges, and limits outside 1..1000. Credentials resolve only from named environment variables. It executes `AdAccount.get_insights`, preserves every returned row and emits atomic hash-linked evidence without the account ID or secrets. No business write exists.

## Verification and failures retained

- exact 18-wheel install, `pip check` and official SDK import/signature probes passed;
- six unit/contract/negative/atomicity tests passed in authored and reconstructed trees;
- Markdown reconstruction produced eight files with zero content hash differences;
- `META_ADS_REPORTING_PACK_PLAN.md` composes the 12-file official source core plus eight adapter files;
- `LIB-FAIL-047..056` preserve byte hashing, encoding/report, PowerShell syntax/binding, path discovery, cleanup-policy, stale fixture, inventory and patch-context failures;
- `UP-FAIL-021..022` preserve the missing wheel license and unsigned tagged commits.

## Admission decision

`REBUILD_VERIFIED / CONDITIONED / INTEGRATION_ONLY`. The code is immediately materializable and invokes Meta's official active SDK, but no real Meta account or API call was exercised. Promotion requires app/business verification, `ads_read`, least-privilege identity/token lifecycle, test ad account and ownership, approved fields/date semantics, comparison against UI/export, pagination/limits/cost, privacy/retention, security, observability, reconciliation and production pilot evidence. The license restriction remains permanent.
