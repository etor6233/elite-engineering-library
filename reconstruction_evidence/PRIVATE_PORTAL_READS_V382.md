# Private portal reads and exact amount presentation — V382

2026-09-11. Library maintenance / TEST02 / T2804 with T2803. Parent V381 checkpoint225,
resume/plan and final VERIFY_LIBRARY passed. Root product readiness remains DISCOVERY/BLOCKED.

## Verified result

V382 corrige dos defectos reales de lectura: admin:read ya no consulta leads sin lead:read y los portales dejan de mostrar unidades menores como importes mayores. Se comparte exactamente el formatter de cotizaciones, sin alterar precio/estado/aceptación.24regresiones focales,260tests web/1skip previo, typecheck/build y4navegadores contra Go/JWKS/PostgreSQL PASS;5negativos HTTP, consultas indebidas0, escrituras0 y snapshot de6tablas idéntico.45/48 permanece, sin promoción global.

The actual AdminPage admitted admin:read then unconditionally fetched /v1/franchise/leads,
whose unchanged Go handler requires lead:read. Two original-code regressions failed with
FORBIDDEN lead:read. The repair conditions only that request and section; overview/orders/cases
retain their exact admin gate. Case-sensitive grants and wildcard behavior remain unchanged.
The Go permission, tenant, organization and customer-subject enforcement was not altered.

Visual inspection of the first real connected browser run then found USD100 shown for a
stored total_minor_units=100. A separate original-display test failed. The same raw interpolation
existed in customer order/quote summaries. The dedicated quote page already had the correct
BigInt/Intl algorithm: it is now a shared BFF helper with only parameter-name substitutions.
Reversing that extraction/import reproduces the original quote source exactly, preserving its
validity, date, state and acceptance decisions. Admin/customer summaries call the same helper.
No binary float division, price/FX/tax/rounding rule or stored value was introduced or changed.

## Independent evidence

11private-view tests and13literal numeric/invalid-input cases pass. Oracles cover USD100→1,00,
JPY123→123, KWD1234→1,234, zero and the maximum safe integer; negative/fractional/non-finite/
unsafe integers and unknown/lowercase currencies retain the existing non-verifiable label.
The full web suite passes260tests plus one pre-existing explicit connected skip, typecheck/build.

The connected Go test reuses the existing real RS256/JWKS test issuer and53byte-identical SQL
migrations in a new protected loopback database. It registers the actual enterprise-query and
journey HTTP/services/PostgreSQL repositories. Five independent HTTP requests prove invalid
token401 and permission/organization403. Real encrypted Next sessions then drive four browser
projects (retries0): admin-only, admin+lead, wildcard, wrong-case grant, other organization/tenant,
non-admin denial, factory inventory read, customer-owned order/quote reads, and the dedicated
quote display. Permission removal on reload hides opportunities. Deliberately overclaimed
frontend fixture permissions with an insufficient bearer still produce500, disclose no domain
records, and recover under the corrected session. That500 is preserved fail-closed behavior,
not a claim of a polished recovery UI. Both browser and API write counters remain zero.

The admin-only bearer caused zero lead queries and the non-admin UI caused zero admin queries.
A before/after PostgreSQL content digest across orders, cases, factory units, leads, stock and
quotations is identical. Synthetic setup records are not business commands or production data.
Final mobile screenshot visually reviewed; labels show USD1,00/USD2,00 and no horizontal overflow.
The first permission-only green run and its misleading-money screenshot remain as history;
the final green run includes monetary and customer checks. No hidden retry or failed run erased.

## Failures, provenance and conditions

FAIL754: two real page regressions under admin-only/wrong-case grants, now green without ACL
expansion. FAIL755: initial test file suffix .test.tsx was outside the configured .test.ts
discovery; exit1 ran zero tests and was not counted as reproduction. Renaming the new non-JSX
test yielded2failed/7passed before the fix. FAIL756: visual and one explicit amount regression
proved the raw-unit defect; shared existing algorithm and final browser values now pass.
The earlier unrelated Firefox teardown FAIL423/742 remains OPEN. Local self-signed test TLS
logs include certificate/prefetch cancellations; actual tested navigations/assertions pass,
and no production TLS or external IdP admission is inferred. PostgreSQL/Next/API/issuer stopped.

No new dependency or upstream.8589exact Next/Sharp artifact files match V376 and the locked
pnpm graph is unchanged. Known pnpm CRLF manifest rewrites were normalized before hashing,
with exact JSON object equivalence. No scan/freshness claim beyond recorded artifact parity.
All new files are AUTHORED / LicenseRef-Workspace-Owner; existing formatter provenance remains
local code over the pinned Intl runtime, not an implementation attributed to another company.

## Canonical scope

BFF0.5.10/51files: two private pages, new page regression, shared amount helper/tests and
read-view documentation. Portals0.14.14/31files: exact quote-helper extraction. Query0.1.2/7files:
one optional connected test, zero production Go changes. Browser0.1.31/14files: existing test
file and package script. Role0.1.2/7files updates compatibility documentation only, preserving
all six source files. HTTP0.1.7 updates only the parent lock. Current metadata now names the
versions actually composed/tested; historical pins and original bytes remain in before receipts.
Five affected plans; six owners reconstruct identical file blocks in absent destinations.
Inventory165packs/1550files/830Markdown/56plans; franchise68/790 and HTTP69/798. Final composition
and Preflight226 passed with164executed steps and56plans; receipts appear below. The original unselected golden path is not promoted.

This closes the demonstrated read permission/display defects, not all portal workflows,
arbitrary dashboard KPIs, assessed onboarding, complete recovery UX, product IdP, native
security/licensing or signed release. All48acceptance definitions/actions/oracles/statuses
remain unchanged:45passed, TEST02/03/07blocked. No full-control closure from this partial boundary.

## Receipts

Local raw evidence resolves through elite-v382-current.txt; fixture databases/artifacts remain
outside distribution. Their local paths are not operational configuration for a consumer.

| Evidence | SHA-256 |
|---|---|
| baseline-red.log | 2c5458fb7ab8e15c27ecae1aecf1607b2b178ba6e7ca4add607c78bdfac3cc45 |
| baseline-discovered.log | ee0cee9bd0a39c4a010bc37f1fa086cc56b0b4d65ab16a5bb5428909bc9e34e7 |
| red-reproduction.json | b60e9d4f56c44ac6ba825f8464b3c1bc0de82b670af5230aca72b439089d3e3e |
| money-red.log | 261c80cf8b098e407462378017c9e68a94702beaaf5745e5fa7402f30c3d6b7b |
| money-red.json | 8ffd7ab719520d01631447272a834ff637e1c7d28579e12ad6b0e6acf25aae36 |
| money-extraction.json | ff3f53312599ca75b11c4427843de20082b42f57f87d5f30a223e7a7db58a663 |
| web-fixed-results.json | ef7312a5aff800d1d801c5f0a887656f60f23cec1e11ce806c9c94560eff45a9 |
| web-money-qualified-results.json | 1bcc3347ff6a835e3fadc453f7cdfb3858e8a4c1cabd02139bb9e7b7d7f6241b |
| admin-regression-money-qualified.log | fc73aa91f7cff793576a09785548d750cbf3b4a02cce371b9b4020966203332a |
| tests-money-qualified.log | abb763feecb107a53aad0b26cb404f90537b63673bb5734827a41e13e3532ede |
| typecheck-money-qualified.log | 99d7c2179cdd316467c51f37901025a0753453664d741214ba56fc631fa6c0bc |
| build-money-qualified.log | 2bb1b541b85b78a9021e7a87ec52f0a5bdcfc40e70b7b6fa3fb62ad52e37f7f0 |
| artifact-parity.json | 7563f34b2a058f7ce4c0f799cb29ad7b97d517326579e335c400497c8723cb43 |
| manifest-normalization.json | cee473e809d7efd3d3f3dbee690df60fafcc5d6e0c3547e7ad920a6f82cdda55 |
| pending-parity.json | b363f30d87f37b0583eb2d27e31ed89984cd80d4d070927abbbf2411bcd179e8 |
| canonical-pending-changes.json | 255311709578e909db6f3dea70dbe85cf2546a5bb9ba60c0e578c14c9e2e00c1 |
| final-connected/result.json | 695d6bf375be76ac00d36c62d1eb17c9e34b95d7966495072d07f3245e6f8f08 |
| final-connected/admin-reads.log | 5f368ca8ce8c080864955cb73f5eec70000775ed4438b8ee1bb87a1b0974b653 |

## Final reconstructed consumer

```json
{
  "status": "PASS",
  "packs": 69,
  "files": 798,
  "parent_files": 790,
  "parent_files_identical_to_connected_candidate": 789,
  "next_env_generated_imports": 2,
  "added": [
    "docs/private-portal-reads.md",
    "src/platform/i18n/money.test.ts",
    "src/platform/i18n/money.ts",
    "src/app/admin/page.test.ts",
    "internal/platform/httpapi/admin_read_browser_test.go"
  ],
  "modified": [
    "docs/role-workspace.md",
    "microsoft_playwright_browser_gate/package.json",
    "reference_http_metrics/source-lock.json",
    "src/app/admin/page.tsx",
    "src/app/customer/page.tsx",
    "src/app/customer/quotes/page.tsx",
    "microsoft_playwright_browser_gate/tests/role-workspace.spec.mjs"
  ],
  "http_binary_byte_identical_to_v375": true,
  "http_binary_sha256": "4643654625c71aaf5b45367bd3a198381257588d073ae0349fc90c1bc1b5ab85",
  "lock_sha256": "b77daa20003311c295d496ddbfda1569819121665c6b6d24e8421b90693da59f"
}
```

## Closure227 and unchanged acceptance

V382 cierre227: admin:read funciona sin consultar leads no autorizados; admin/customer muestran importes correctos mediante el formatter exacto de cotizaciones compartido.260tests web/1skip previo,24focales,typecheck/build y4navegadores con Go/JWKS/PostgreSQL PASS;5negativos HTTP,0consultas indebidas,0escrituras y snapshot6tablas idéntico.6packs reconstruidos,composición68/790 y HTTP69/798;8589artefactos y binario HTTP idénticos. Preflight226164pasos ejecutados/56planes PASS; Docker ausente. FAIL754/755/756 corregidos con originales preservados.45/48 conserva TEST02/03/07 pendientes; no cerrar integración,seguridad o release completos por estos fixes.

```json
{
  "checkpoint_under_test": 226,
  "executed_steps": 164,
  "all_executed_steps": "PASS",
  "plans": 56,
  "missing_toolchains": [
    "docker"
  ],
  "inventory": {
    "packs": 165,
    "files": 1550,
    "markdown": 830,
    "plans": 56
  },
  "franchise": {
    "packs": 68,
    "files": 790
  },
  "http_reference": {
    "packs": 69,
    "files": 798
  },
  "private_read_browser_projects": 4,
  "independent_http_negatives": 5,
  "admin_only_lead_queries": 0,
  "non_admin_admin_queries": 0,
  "browser_and_api_writes": 0,
  "unchanged_domain_snapshot_tables": 6,
  "browser_retries": 0,
  "unit_tests_passed": 260,
  "focused_private_view_and_amount_tests": 24,
  "existing_explicit_skip": 1,
  "controls_passed": 45,
  "controls_total": 48,
  "preflight_log_sha256": "62cc1c425554c424e5326326a43ee7084813341721e152ff75ad0accd2f9d6f2",
  "preflight_json_sha256": "deacc73752466726959c3105895b5bfe7bba45687625f1f3b97ccc16addaa82e",
  "final_parity_sha256": "eeff9747a36f7d1b098c5c741abf9a70a570ee96da4a1b75dca5edb5cef5a5b9",
  "acceptance_invariance_sha256": "d09615ac7cf9f189fd76a693ea5be031f05ca50a5d0152a225ca310b06798e41"
}
```
