# Public web metadata integration — V377

Date2026-09-10. Library maintenance, existing T2802/T2804 owners.

V377 integrates actual Next canonical/robots/sitemap/WebSite JSON-LD under the existing BFF owner0.5.6/43files:5new AUTHORED files and3page changes; no dependency update. Public indexing is opt-in, limited to the real home/model-list pages and public render assets; private pages default noindex.137web tests/1preserved explicit integration skip, strict typecheck/build and3actual HTTP configurations PASS. The same build proves disabled/enabled/changed origin,Host poisoning isolation,mandatory cache revalidation,private noindex and JSON-LD HTML escaping/CSP nonce. Connected qualification: three browser projects passed initially, Firefox failed teardown with known FAIL423/742; a separate exact Firefox23phase rerun plus4agenda projects passed. These results cover4projects/92distinct phases across two runs, not a flawless92phase initial run or a Firefox lifecycle fix. Known upstream incident stays open.45/48 unchanged. See reconstruction_evidence/PUBLIC_WEB_METADATA_INTEGRATION_V377.md; final canonical composition and Preflight pending.

## Scope and exact authority

The selected BFF uses Next16.3.4 already qualified in V376. Native/vendor/redistribution gates and FAIL532 remain; the unchanged152version package graph is not rescanned or relabeled as a new full-security proof.8589patched dependency files match the signed/inspected V376 inventory. No dependency, public-source artifact group, account, private corpus or paid service was added.

Authority snapshots:

- [robots](https://nextjs.org/docs/app/api-reference/file-conventions/metadata/robots), SHA256 cdd0782e111f747a0870e15a74e55cbfd9e481626a2b0ba481697971b309c363.
- [sitemap](https://nextjs.org/docs/app/api-reference/file-conventions/metadata/sitemap), SHA256 b35f459ee5d2243c11198f5567849b438061f962ddb956fdcd01aa8d6de419b0.
- [generate-metadata](https://nextjs.org/docs/app/api-reference/functions/generate-metadata), SHA256 5d32c82f826d6e6b7e712b18247b5bcc7428df67c334bbc0873f2f3144f356d1.
- [json-ld](https://nextjs.org/docs/app/guides/json-ld), SHA256 52a1f258533f8e880a5e81a274349360ad7831e19bac516fc66132f37eb07663.

[Google JavaScript SEO guidance](https://developers.google.com/search/docs/crawling-indexing/javascript/javascript-seo-basics) requires crawlable rendering resources; the allowlist retains /_next/static/ and the icon. Robots is advice, not authorization or index removal. This code does not claim crawler acceptance, ranking, translated journeys, per-model publishing, domain verification or multi-tenant origin resolution. The GO-SEO-CORE memory helper remains separate and CONDITIONED; its local60/160policy is not imposed as a Google rule. This closes the reference HTTP rendering boundary, not all SEO/i18n or TEST02/FAIL385 integration.

## Executed behavior

One production build, three isolated runtime configurations, six rendered public responses, three private responses and six robots/sitemap responses. Disabled means empty sitemap/noindex/noJSON-LD/disallow; enabled contains exactly two URLs. A second origin on the same build replaces the first, with exact mandatory-revalidation headers. Host/X-Forwarded-Host/Proto and input x-nonce cannot choose canonical origin or script nonce. A configured string containing a closing script tag is preserved as text and JSON but cannot terminate the JSON-LD element. The catalog backend in this focal probe is an explicit empty synthetic HTTP fixture, not a database-persistence claim.

Actual PostgreSQL/browser regression separately uses53unchanged migrations, isolated ownership and all18journey flags, including return multi-tab/session gates. Initial run retains72successful phase receipts (69from three complete projects plus3Firefox prefix) and a real Firefox teardown failure. Retry preserves all source/test/lock bytes, selects Firefox only, requires23phases and then4agenda projects, retries0 inside Playwright. Do not add the3prefix duplicates to92. Both owned PostgreSQL clusters stopped. Trace174records identifies three propagated errors in context/fixture/after-hook teardown; this does not prove the browser race fixed.

## Reproduction and retained failures

Compose FRANCHISE_COMPLETE_PACK_PLAN in an absent directory, use exact V376 Node/pnpm and frozen offline ignore-scripts install. Run test/typecheck/build. With a loopback catalog fixture, run the public HTTP probe against the production Next server in disabled and two enabled-origin configurations; preserve output HTML/XML/header hashes. Then use a new owned elite_confirmation_* database,53exact migrations and the current connected Go/Playwright harness. Require complete results per browser and preserve any failure. No automatic retry is hidden as a first-run PASS.

FAIL738 strict optional environment typing corrected without changing tsconfig; FAIL739 probe HTTP header case/indentation corrected; FAIL740 provisional no-store assumption replaced by the exact verified Next loader revalidation contract and runtime-origin proof. FAIL741 current dependency bill corrected to16.3.4, historical V11716.3.2 retained. FAIL742 remains an upstream lifecycle incident. A prepare helper also rejected an unobserved section heading before file writes (FAIL727 recurrence); actual Configuration surface heading was then used.

## Receipt index

Raw evidence stays outside distribution under the stage resolved locally by elite-v377-current.txt. Hashes bind raw files; display paths omit machine identifiers.

| Relative evidence | SHA256 |
|---|---|
| tests-final2.log | 57e20c9936ce6dc7a3e288999fd4bda8345db6195f9f33140bbe8d9c6fe428e8 |
| typecheck-final2.log | 99d7c2179cdd316467c51f37901025a0753453664d741214ba56fc631fa6c0bc |
| build-final2.log | f018627c67a7f78adc2b1b4ad753a08ab1c7d64d3e3b4c0f39fe694b58fdb308 |
| web-final2-results.json | 2ac0950c1cfecc5d424d166907894d9943a813d9c8e42342cb2aab66889f6cfc |
| artifact-parity.json | 2f676c4fae62514663d2b233198b64508629d214160b31f515251fe571ad792d |
| pending-block-parity.json | 67092721d123bd768f3a8632f5f11dc565ad52f2b162ddc5ba303b7e632c2f3a |
| canonical-pending-changes.json | 5c7c0dbb88db5c348d525a3e1c89531b41e87e479e987d3aa49dd47b07a39dd8 |
| probe-public-http377.py | 49cf2d1ea6e47c17027e799f670172cfadaf21ffea89193977a2264f8e4fc87e |
| connected-web377.py | f2780a1279ec26b6847df70cd1c89482ed341fdeb8abd0406799ec2a03ece715 |
| connected-firefox-retry377.py | 9647820f51cd3dc8fd0d55fab1390206fcaaa1c06ce40ca26f18c40d79e49b01 |
| firefox-failure-trace-inspection.json | 866e11264138664f3d8b099d97143576ccc6d66f4ce11ac4d066dee9005e39a2 |
| http-c476e71a85e74664bded20e9951ba868/public-http-result.json | 7c53613500fab09426b9ab1defac440343f07f409ad32178573627b1943956f3 |
| connected-0aee33742c6e42a29e67936bd0853244/browser-connected.log | 4c440d45188912d49444505e781feb09604c2be0cb46469b38b2f4427f721a45 |
| connected-0aee33742c6e42a29e67936bd0853244/result.json | 3ed8e700b2f8bceca14176023e18e530a3a27f5c490f2b738ae795345d67c616 |
| connected-f50a28280ed24896bd5af5570dc1e765/browser-connected.log | 20a43dd3a3af20a1b3c7415a2ab59beb2c6f7d39a6ad660b04497299824827fc |
| connected-f50a28280ed24896bd5af5570dc1e765/browser-agenda.log | 2bf05f6f22f7409279f913bca0701c5d01338a7a0a55c2a0e947d588bcfd1235 |
| connected-f50a28280ed24896bd5af5570dc1e765/result.json | cb1438c6a19aa18741664e527251636b1b54378f8ca2f6a61f65e94107e51e97 |

The43-file BFF reconstruction matches all8changed candidate files;35existing blocks remain identical. No claim or status is promoted from documentation alone. Conditions in docs/public-indexing.md remain mandatory for a consumer.

## Final reconstructed consumer

```json
{
  "status": "PASS",
  "packs": 68,
  "files": 768,
  "parent_files": 760,
  "parent_files_identical_to_connected_candidate": 759,
  "next_env_generated_imports": 2,
  "added": [
    "docs/public-indexing.md",
    "src/app/robots.ts",
    "src/app/sitemap.ts",
    "src/platform/seo/public-indexing.test.ts",
    "src/platform/seo/public-indexing.ts"
  ],
  "modified": [
    "reference_http_metrics/source-lock.json",
    "src/app/layout.tsx",
    "src/app/page.tsx",
    "src/app/models/page.tsx"
  ],
  "http_binary_byte_identical_to_v375": true,
  "http_binary_sha256": "4643654625c71aaf5b45367bd3a198381257588d073ae0349fc90c1bc1b5ab85",
  "lock_sha256": "b77daa20003311c295d496ddbfda1569819121665c6b6d24e8421b90693da59f",
  "reused_go_scope": "All Go,SQL,rule and dependency files unchanged; no repeat of unchanged runtime or unit/vet suites"
}
```

## Final closure and unchanged acceptance

V377 cierre217: integración pública canonical/robots/sitemap/JSON-LD en BFF0.5.6/43archivos comprobada con137tests,3configuraciones HTTP reales y composición final67/760. Regresión conectada:69fases de tres perfiles PASS en el primer run y23Firefox+4agendas PASS en una repetición aislada; primer fallo de teardown conservado y FAIL423/742 abierto. Preflight216:164pasos ejecutados PASS,56planes compuestos;165packs/1526files/825Markdown/56planes,1274AUTHORED/145ADAPTED/107VERBATIM. Docker sigue ausente. Binario HTTP idéntico a V375; fuentes Go/SQL/regla/dependencias sin cambios. Los48criterios/status son idénticos:45/48, TEST02/03/07 bloqueados. Esta frontera SEO de referencia está integrada; faltan i18n/localización y otras equivalencias funcionales, seguridad/admisión nativa y release integral. No contar el replay Firefox como reparación upstream ni repetir creación/recuperación ya probadas.

```json
{
  "checkpoint_under_test": 216,
  "executed_preflight_steps": 164,
  "all_executed_steps": "PASS",
  "composed_plans": 56,
  "availability_status": "BLOCKED",
  "missing": [
    "docker"
  ],
  "inventory": {
    "packs": 165,
    "files": 1526,
    "markdown": 825,
    "plans": 56
  },
  "provenance": {
    "AUTHORED": 1274,
    "ADAPTED": 145,
    "VERBATIM": 107
  },
  "franchise": {
    "packs": 67,
    "files": 760
  },
  "http_reference": {
    "packs": 68,
    "files": 768
  },
  "controls_passed": 45,
  "controls_total": 48,
  "preflight_log_sha256": "adac0735b698c00feae228fe72eb649b26262ec200f3db95a41c00a39e8af254",
  "preflight_json_sha256": "d3baddc9529438ec8d64bbf32ea5045e9c1bf8c7758bd6a30580de787aa6ba23",
  "final_parity_sha256": "5bbd9ddcd39eb1b5c1a03b0f22c258bcf7306f5bf4867f89d172570ec4f9f5de",
  "acceptance_invariance_sha256": "989931e4160b74414768c70294486202e8739676a1ec4088bf739abb82530fa0"
}
```

Full Preflight ran on checkpoint216. Successor217 records the completed receipts, without changing materialized source or claiming whole-control closure.
