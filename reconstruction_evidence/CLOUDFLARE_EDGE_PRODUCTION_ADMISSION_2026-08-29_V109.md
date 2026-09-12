# Cloudflare Edge Production Admission — V109

Date: 2026-08-29

## Scope and non-claim

This checkpoint adds one exact official Cloudflare SDK source and a local fail-closed semantic adapter for the existing `EDGE_CDN_WAF` production control. It does not claim that Cloudflare authored the adapter, that an API success proves security, or that any real account, zone, hostname, origin, WAF policy, cache policy or production traffic was exercised.

## Exact official source identity

- repository/release: `cloudflare/cloudflare-go` `v7.9.0`;
- release ID/published: `373826187` / `2026-08-20T15:05:12Z`;
- lightweight tag commit: `4542c2cb87007f7de371e118b3cf3d7cc7f71daf`;
- GitHub commit verification: `verified=true`, tree `87112de22edbac2311945faad8b1f0484fdf87b2`;
- exact commit archive: 8,024,155 bytes, SHA-256 `67fdde9324bfe519e1c39bf0fe214798de2bb00c758200f59b5c1f19912ab380`;
- Apache-2.0 `LICENSE`: SHA-256 `8659207f29b327845dd32bcc5bcbde42f7b9ed221d8596c38c13a9ac00cdb07e`;
- `go.mod`: SHA-256 `c00b553a32436f3ee320478a8121f3d53a961798a0097e4dcee60aeff0f70671`;
- generated `rulesets/ruleset.go`: SHA-256 `bb012b8dc28039608d6533b8ca0df839149e97af0d30e27d5db02bc6e658bef5`.

The source lock now contains 119 entries. Both `enterprise-platform-leaders` and `cloudflare-network-edge` select this SDK; the latter selects it together with exact Cloudflare Pingora 0.8.1. Lock validation reports 119/119, fourteen source profiles pass, enterprise selects 20 and Cloudflare edge selects 2.

Official references: <https://github.com/cloudflare/cloudflare-go/releases/tag/v7.9.0>, <https://github.com/cloudflare/cloudflare-go/tree/4542c2cb87007f7de371e118b3cf3d7cc7f71daf> and <https://developers.cloudflare.com/api/resources/rulesets/methods/list/>.

## Executable semantic adapter

`SECURE-OPS-DELIVERY-CORE` 0.7.0 adds three `AUTHORED` files without embedding Cloudflare source:

1. `cloudflare-edge-admission.template.json`;
2. `validate_cloudflare_edge_admission.py`;
3. `test_validate_cloudflare_edge_admission.py`.

The adapter requires four ordered executions through the common no-shell, explicit-authorization, binary/argv/output hash-bound boundary:

- a Cloudflare API configuration snapshot for the exact zone and hostname;
- a direct-origin bypass negative;
- two-request public-cache behavior;
- authenticated private-response cache bypass.

It admits only proxied A/AAAA/CNAME records, approved active managed and custom WAF rulesets, managed `execute`, a defensive custom action, Full (strict) TLS, TLS 1.2 or 1.3, active certificate, HTTPS redirect, approved origin protection, a repeat public `HIT` with identical body hash, and a private `BYPASS`/`DYNAMIC` response marked private/no-store and not shared. Policy requires two distinct approvers. All runs must bind the same exact project probe binary, official cloudflare-go source identity, canonical argv, target and observation hash.

Nine regressions pass. They prove the positive contract and reject unproxied DNS, missing managed WAF, weak TLS, successful origin bypass, repeated public cache misses, private cache hits, a non-zero probe exit and observation tampering.

## Reconstruction and global gates

- acquisition pack 0.4.73 reconstructs 23/23 files; its two suites pass;
- secure operations 0.7.0 reconstructs 25/25 files;
- production-admission suite: 47 tests, 46 PASS and one platform symlink skip;
- `VERIFY_LIBRARY_PASS`: 73 packs, 727 materializable files, 437 Markdown before adding this evidence, backend 16/114 and web 5/69;
- `VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit`: 73 packs and 119 upstream sources.
- official Go 1.26.7 Windows archive: 74,955,002 bytes, SHA-256 `f4f534a486e4bc3387fa18f08208f2f854b7aaea8a08f2a2d829a914a05abb11`;
- `VERIFY_EXECUTABLE_LIBRARY_PASS mode=Foundation`: backend Go tests, 47 production-admission cases, web typecheck/21 tests/build, Playwright runtime 4 plus target 8 and Lighthouse five-run quality gate all pass. The final status remains `PROJECT_CONDITIONED`.

Canonical checkpoint hashes before adding this evidence:

- `OFFICIAL_UPSTREAM_ACQUISITION_CORE.md`: `da640b07bfb3ade1b6902b19fa88461b1ba640072e07820a69cfc239e493acf3`;
- `SECURE_OPERATIONS_DELIVERY_CORE.md`: `18df02dc4e40990584ce63e1f821128eb8526e511e35b823f3a2e292b7a5ee85`;
- `ENTERPRISE_BACKEND_PACK_PLAN.md`: `11e5f57d085666c55275714bf043e024ee8c6398cbb3f130e112c188e365b401`.

## Failure memory and remaining work

`LIB-FAIL-1201` through `LIB-FAIL-1204` retain the blocked destructive staging command, recurrent array CLI invocation, initially incorrect focal hash and blocked recursive cleanup. `UP-FAIL-198` retains the Cloudflare boundary: official SDK code cannot provide a project's real edge evidence. The verified leaf-first cleanup ended with `V109_TEMP_REMAINING=0`.

The library therefore has 1,204 local failure identities plus 198 upstream conditions: 1,402 unique IDs. `EDGE_CDN_WAF` is now an executable adapter, but no concrete project is `PRODUCTION_ADMITTED`. Four semantic adapters remain absent: `IDENTITY_AUTHORIZATION`, `PROVIDERS`, `DEPLOY_ROLLBACK` and `BUSINESS_ACCEPTANCE`; all eight real receipts must still pass together for one immutable release digest.
