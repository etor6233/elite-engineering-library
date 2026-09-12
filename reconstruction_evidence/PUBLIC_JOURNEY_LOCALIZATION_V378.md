# Public journey localization — V378

2026-09-10. Library maintenance under existing T2802/T2804 owners. Entry checkpoint217 validated. No new product, platform, external service, private data or financial policy.

## Verified result

V378 integra el recorrido público home→modelos→consulta→ubicaciones→solicitud de turno en español/inglés, con fallback español explícito y zona horaria del mercado.12/12recorridos navegador/configuración PASS,12leads y12turnos durables sin duplicados/sobrecupo; pérdida de respuesta y replay probados.164tests/1skip explícito conservado, typecheck/build PASS. El probe HTTP reproduce y corrige labels españoles bajo en-US y offset fijo de3horas bajo UTC. Cuatro archivos nuevos y ocho existentes modificados bajo owners BFF/portales/browser, sin nuevas dependencias ni reglas de dominio.45/48 permanece; no cerrar TEST02 por una sola frontera.

## Actual before and after

The previous exact V377 Next production build accepted en-US/UTC configuration, labeled HTML en-US,
but rendered Spanish messages and displayed12:00 for an appointment at15:00Z. This was an actual HTTP
probe against a bounded synthetic catalog/slot fixture. The new build displays English and15:00UTC.
Raw HTML and source revisions are preserved; the fixture proves presentation only.

Separately, the unchanged Go public-lead/public-appointment harness runs the actual BFF/Go/PostgreSQL
flow under a newly owned loopback elite_browser_* database with53exact migrations. Three configurations
(es-AR/Buenos Aires, en-US/UTC, fr-FR/UTC→Spanish fallback) each run Chromium desktop/mobile,
Firefox and WebKit with zero Playwright retries. Per configuration SQL independently verifies four leads,
four consent records, four lead events/receipts, four requested appointments and four appointment
events/receipts, plus no capacity mismatch. The browser loses a response after the real commit,
retries with the same key and receives the same reference. Divergent intent, missing keys,
missing consent, forged state and full capacity remain rejected. The owned PostgreSQL process stopped.

The first es-AR run failed only the newly added Node-versus-WebKit Intl string oracle: U+0020 versus
U+00A0 within Spanish p. m. Its trace/video/screenshot/log are preserved. The corrected comparison
equates only the two non-breaking space variants with space; all digits, punctuation, kind,
configured time zone and submitted UTC instant remain asserted. All12cases then passed in a new run.
This is a test-oracle correction, not a Firefox lifecycle repair. FAIL423/742 stays open.

## Design and boundaries

The BFF now owns the fixed AUTHORED es/en message catalog, canonical supported-locale selection,
number/plural formatting, and loading the default market's explicit timeZone. Invalid/unsupported
language selects and labels Spanish; invalid timeZone rejects configuration. No locale comes from
request headers, query, cookies or untrusted translation files. Unknown message keys remain literal
keys, with own-property lookup. Static complete catalogs and TypeScript check prevent an incomplete
English dictionary. Node/Intl supplies formatting semantics; no downloaded CLDR tables or new dependency.

Home, both catalog routes, lead form, location page and appointment form use the same resolved locale.
Navigation/footer/metadata and WebSite JSON-LD follow that actual language. Private portals remain
Spanish under main lang=es; each translated public section overrides lang. Existing dynamic backend
model names/classification are not auto-translated. No CMS, user locale negotiation, localized URL
variants/hreflang, general CLDR interpolation, translated private journeys or legal consent approval.
The standalone GO-I18N-CORE remains an isolated binary-plural helper, not substituted into product.

All request payload fields, endpoint paths, UTC instants, idempotency keys and domain/SQL implementations
remain unchanged. Unknown transport exceptions now produce a fixed localized failure instead of raw
engine text; invalid receipts preserve a distinct fixed message and retry guidance. No third-party
text is interpolated as HTML. Target policies and business terminology require consumer review.

## Reuse and validation

164Vitest tests PASS, one existing explicitly separated connected test skip retained; typecheck and
Next production build PASS. The affected connected public journey ran separately,12cases/3configurations.
No unrelated92phase private-operation replay: all Go, migrations, backend contracts and private components
are unchanged, and the original default-Spanish navigation labels remain exact. Prior V376/V377 private
runtime evidence stays bounded to its original scope. Full library Preflight218 passed all164executed steps and56compositions; Docker availability remains blocked.

8589installed Next/Sharp files match the exact authenticated V376 artifact inventories; the lock remains
b77daa20003311c295d496ddbfda1569819121665c6b6d24e8421b90693da59f.
This reuses scoped advisory remediation, not a fresh whole-native/vendor scan or license clearance.
All existing native redistribution, Docker, security and release blockers remain.

Official semantic references consulted: [ECMA-402](https://tc39.es/ecma402/),
[Node Intl support](https://nodejs.org/api/intl.html),
[W3C language declarations](https://www.w3.org/International/questions/qa-html-language-declarations).
The current Node documentation is newer than the pinned runtime; behavior here is proved on the exact
Node24.20.0 binary and browser builds already retained, not inferred from that newer documentation.

## Retained failures and canonical reconstruction

FAIL743 is the demonstrated configuration defect. FAIL744 was a candidate-only shadowed state name
caught by strict typecheck. FAIL745 is the Intl space oracle. FAIL746 was the evidence reader's
legacy Windows path length; extended absolute paths verified every8589file, no waiver.
FAIL719 recurred in read-only guessed paths; failures and corrections are retained.
Four canonical packs reconstruct with exact source parity: BFF47, portals20, browser11, HTTPreference8.
The final integrated consumer and Preflight receipts are appended below after actual execution.

## Receipts

External local staging is resolved by elite-v378-current.txt; raw data remains outside distribution.
| Evidence | SHA-256 |
|---|---|
| baseline.json | e8bacff9777670a708d64db9bf75f850f9a95063b8d2a8e78d8d3312c46cb1c0 |
| web-fixed-results.json | dd212a9b523535136b5b58ead38fa30552f933a2f5fa65fb1fa836bd03644fb4 |
| http-red-green.json | b7e69417800c52efdc43b3284dea79f3794706e47ce55f50f71146131c82e23e |
| artifact-parity.json | 7563f34b2a058f7ce4c0f799cb29ad7b97d517326579e335c400497c8723cb43 |
| intl-runtime.json | ff1cf8f1f0e68f680ddbc548714d49a8d1bb13181574352ef0c96cd017a9d6ec |
| pending-parity.json | 2edfdcb8b00f4a25dd692743cc8173f784053db67f6183102262e629ca13dce1 |
| canonical-pending-changes.json | 42377bf3d2abe5ab47886b6b60f594f69aeb6cc8c5e3ca69ff2b10ac89d22a08 |
| connected-public-fixed378.py | 0f783620df952728cede047f0d44553a62aa9c885071d85bbe1134d21ec851c0 |
| probe-locale378.py | a8767f6bfd3e6fd9c80afcfb606902adb92ce87affa2a05c9680099d0a6b2bbf |
| connected-fdaf246beee548059cd3f03674cfcfd8/result.json | babfe5319cb1015bda0d86c571fcd7c688863f6cf3118044f6a8f4e0cd5911a6 |
| connected-fdaf246beee548059cd3f03674cfcfd8/public-es-AR.log | 9c625f4b37eb5d72e50dbee4761c48fa1e0d7ab25438fa4fbd901f13c1bfd84a |
| connected-fdaf246beee548059cd3f03674cfcfd8/public-en-US.log | e1bda6b5153f975cded01a03912c219f9282476ff86d2de617dd6acbd81d04cf |
| connected-fdaf246beee548059cd3f03674cfcfd8/public-fr-FR.log | 9f384341495348c05662d2efec0ead4da7fa9cc2e140fc8c2e331396423c345e |

## Final reconstructed consumer

```json
{
  "status": "PASS",
  "packs": 68,
  "files": 772,
  "parent_files": 764,
  "parent_files_identical_to_connected_candidate": 763,
  "next_env_generated_imports": 2,
  "added": [
    "docs/public-localization.md",
    "src/platform/i18n/load-public-locale.ts",
    "src/platform/i18n/public-catalog.test.ts",
    "src/platform/i18n/public-catalog.ts"
  ],
  "modified": [
    "reference_http_metrics/source-lock.json",
    "src/app/layout.tsx",
    "src/app/page.tsx",
    "src/app/connected/lead-form.tsx",
    "src/app/connected/page.tsx",
    "src/app/locations/appointment-form.tsx",
    "src/app/locations/page.tsx",
    "src/app/models/page.tsx",
    "microsoft_playwright_browser_gate/tests/enterprise-web.spec.mjs"
  ],
  "http_binary_byte_identical_to_v375": true,
  "http_binary_sha256": "4643654625c71aaf5b45367bd3a198381257588d073ae0349fc90c1bc1b5ab85",
  "lock_sha256": "b77daa20003311c295d496ddbfda1569819121665c6b6d24e8421b90693da59f"
}
```

## Closure219 and unchanged acceptance

V378 cierre219: recorrido público localizado es/en/fallback y zona horaria configurada,12/12casos navegador/configuración,12leads y12turnos durables sin duplicados/sobrecupo;164tests y typecheck/build PASS. Composición final67/764 y HTTP68/772; Preflight218164pasos ejecutados/56planes PASS, Docker ausente. Sin nuevas dependencias;8589archivos autenticados idénticos y ejecutable HTTP idéntico a V375.45/48, TEST02/03/07 siguen bloqueados. La brecha pública de idiomas queda demostrada; portales privados y demás integraciones mantienen su alcance. Se conservan los fallos iniciales y FAIL423/742 abierto.

```json
{
  "checkpoint_under_test": 218,
  "executed_steps": 164,
  "all_executed_steps": "PASS",
  "plans": 56,
  "missing_toolchains": [
    "docker"
  ],
  "inventory": {
    "packs": 165,
    "files": 1530,
    "markdown": 826,
    "plans": 56
  },
  "franchise": {
    "packs": 67,
    "files": 764
  },
  "http_reference": {
    "packs": 68,
    "files": 772
  },
  "public_browser_configurations": 12,
  "durable_leads": 12,
  "durable_appointments": 12,
  "unit_tests_passed": 164,
  "existing_explicit_skip": 1,
  "controls_passed": 45,
  "controls_total": 48,
  "preflight_log_sha256": "f8cdfaae8f1f1e8f61ab8329ce23fa154e538c82c4d090f255dc177708a561fe",
  "preflight_json_sha256": "c3d6356e0f1d0cf9fd3998bb66a64a4d7713b6ffe09ce04353aa36ffa0e96a39",
  "final_parity_sha256": "9fcef6ba909d1cb8f640780ea2dc3b58b73ddff6d6093cd7e599bb934a3966e9"
}
```
