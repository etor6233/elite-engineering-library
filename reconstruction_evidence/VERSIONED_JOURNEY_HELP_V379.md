# Versioned operational journey help — V379

2026-09-10. Library maintenance, TEST-02 / T2804. Entry checkpoint219 validated.
Current scope: search/read of the existing quote and WhatsApp operational guidance under
the existing BFF session and portal owners. No new product readiness, business policy or external service.

## Demonstrated boundary

V379 conecta búsqueda/lectura de ayuda operativa por permiso y versión exacta: cotizaciones y estado WhatsApp reutilizan sus textos1.0.0.186tests/1skip explícito heredado, typecheck/build y4navegadores PASS sin reintentos; sesión vencida/adulterada, cambio de permiso, versión ausente, lecturas tardías y recuperación503/malformed comprobados. Sin operaciones ni nuevas dependencias. Sólo cierra esta frontera de referencia;45/48 y TEST02/03/07 bloqueados, CMS/capacitación integral y target pendientes.

Before this change the two guides were inline-only. /help now provides bounded literal
NFC/case-insensitive search, current applicable guides and links pinned to id/version.
The exact previous paragraphs and1.0.0 versions are extracted into one shared source;
both inline components reuse it. A missing revision never silently selects the latest.
No domain command, price, payment, stock, delivery, notification/provider or database code changes.

Every GET reopens the existing encrypted BFF cookie. customer:self selects quote guidance;
appointment:manage plus whatsapp_status_history selects notification guidance. Wildcard
permission still obeys the feature flag. Inapplicable article and unavailable revision both404;
no applicable permission403; absent/expired/adulterated session401. Invalid/duplicate query
parameters fail400. Cross-site requests fail before session access. Unexpected errors are
redacted503. No tenant, organization, token, customer or transaction identifier is returned.
All help responses use no-store. Only GET is implemented; an actual POST returns405.

The guide text is generic PUBLIC operational content already present in client bundles.
Permission filtering selects relevance; it is explicitly NOT an ACL for confidential documents.
There is no tenant document store, CMS publication/archive, historical version repository,
external search, model learning, enrollment, completion or personnel assessment. Private help
remains Spanish and says so. The independent GO-HELP-CENTER-CORE is not imported or promoted.

The client clears previous results before new reads, aborts superseded requests, validates
response schema, clears when hidden and revalidates on focus. Errors offer explicit GET recovery.
The browser fixture holds a real authorized response while another search completes, then
releases it; stale guidance cannot replace the current search. Permission changes, malformed
responses and503 do not preserve old results. Existing cookie-expiry/revalidation semantics
remain: this does not prove immediate external IdP revocation of an already issued session.

## Actual qualification and reuse

186Vitest tests passed; one pre-existing connected-test skip remains explicit. Typecheck and
production Next build passed. Four Playwright projects passed on the first run with retries0:
Chromium desktop/mobile, Firefox, WebKit. Each exercises guest, customer, operator, unprivileged,
wildcard, expired and adulterated sessions, literal search, version pinning, unavailable versions,
role change,503/malformed response, a superseded read and recovery. It observes zero browser
operation requests; a separate explicit API negative proves POST405. Mobile overflow check
passes; its retained screenshot was visually reviewed. This does not close the older unrelated
Firefox teardown incident FAIL423/742.

The test uses the actual production Next/JWE implementation over an owned HTTPS loopback proxy
built with the already admitted Go stdlib. The new proxy is a test-only bounded fixture, ends on
stdin EOF or10minutes and must never be deployed. The owned Next and proxy both stopped.
No PostgreSQL instance or provider was started because these public release-bound guides have
no durable effects. The test does not certify external IdP provisioning or live support.

Frozen offline installation with scripts disabled; exact8589installed Next/Sharp artifact files
match the retained V376 inventory. The dependency lock is unchanged:
b77daa20003311c295d496ddbfda1569819121665c6b6d24e8421b90693da59f.
No third-party source acquisition or dependency promotion. Native redistribution, Docker,
security, deployment and target acceptance conditions remain. Domain Go/SQL sources unchanged;
do not replay unrelated92private phases or12public locale cases without a material delta.

## Canonical owners and freshness correction

BFF0.5.8 changes two navigation/catalog files. Portals0.14.12 adds eight files and changes the
two inline-help components,28files total. Browser gate0.1.28 adds the test and bounded TLS fixture,
13files total, with test:help and explicit ELITE_HELP_E2E. HTTP reference0.1.4 updates only its
parent composition lock. Four affected plans select those versions. Four packs reconstructed in
new empty directories and matched their candidate sources exactly.

The portal metadata still stated Node24.14.0/Next16.3.2 despite the already proven pinned
Node24.20.0/Next16.3.4 runtime from V321/V376. Its current metadata now agrees with those retained
artifacts; before bytes remain in canonical-before. This is a description correction, not an
additional dependency upgrade or retroactive rewrite of the older evidence.

The full67/774 and68/782compositions match the tested candidate; Preflight220 passed all164executed steps and56plans. Docker remains unavailable.
Global control definitions, actions, oracles and statuses are preserved; help stays a partial
overlap in the21-core map because this read boundary does not implement a general CMS or training.

## Receipts

Local raw artifacts resolve through elite-v379-current.txt; they remain outside distribution.
| Evidence | SHA-256 |
|---|---|
| web-results.json | c678dc33016f763f27965557da89919228adb2bc079b322b7414f70323053a7e |
| tests.log | 2b590e857fcd40544f8cae33a95bb68b3537f005de27c705e5d8550645cc2cd5 |
| typecheck.log | 99d7c2179cdd316467c51f37901025a0753453664d741214ba56fc631fa6c0bc |
| build.log | b41e9643c74190550945a6495e514fe1c03766dc514018f560be4875234591bb |
| help-browser/result.json | f22cc0ee85077bf6fb06960f46906b33f9cccc2cc6995e46e3be51728d327761 |
| help-browser/help-four-browsers.log | 6e59d15de6d35462c8342274009bdfaf7e4eaebe4cf9be8197e4404cda405fe1 |
| artifact-parity.json | 7563f34b2a058f7ce4c0f799cb29ad7b97d517326579e335c400497c8723cb43 |
| pending-parity.json | 8f944cb5f99d77098964d04d01abb8d398950c9c2de1454e1a05d1885a1380b8 |
| canonical-pending-changes.json | d7ff01b08861ab0dd0b524e34978dd34f603f4368db42a5e51c2ebe04bd6a47e |
| implement379.py | 5eef1ae3a8e2fe3082a26640186ae467324838c2eacfab8ffa8dd9e836445da3 |
| browser379.py | a45c271829cf1acb5be45024d305f8c56638089fdd95e0aea93b79135e8dcc7e |

## Final reconstructed consumer

```json
{
  "status": "PASS",
  "packs": 68,
  "files": 782,
  "parent_files": 774,
  "parent_files_identical_to_connected_candidate": 773,
  "next_env_generated_imports": 2,
  "added": [
    "docs/journey-help.md",
    "src/components/journey-help-panel.tsx",
    "src/platform/help/catalog.ts",
    "src/platform/help/content.ts",
    "src/platform/help/contract.ts",
    "src/app/help/page.tsx",
    "src/app/api/enterprise/help/route.test.ts",
    "src/app/api/enterprise/help/route.ts",
    "microsoft_playwright_browser_gate/fixtures/help-loopback-proxy.go",
    "microsoft_playwright_browser_gate/tests/journey-help.spec.mjs"
  ],
  "modified": [
    "microsoft_playwright_browser_gate/package.json",
    "microsoft_playwright_browser_gate/playwright.config.mjs",
    "reference_http_metrics/source-lock.json",
    "src/components/appointment-notification-status.tsx",
    "src/components/customer-quote-actions.tsx",
    "src/platform/config/registry.ts",
    "src/platform/i18n/public-catalog.ts"
  ],
  "http_binary_byte_identical_to_v375": true,
  "http_binary_sha256": "4643654625c71aaf5b45367bd3a198381257588d073ae0349fc90c1bc1b5ab85",
  "lock_sha256": "b77daa20003311c295d496ddbfda1569819121665c6b6d24e8421b90693da59f"
}
```

## Closure221 and unchanged acceptance

V379 cierre221: búsqueda/lectura de ayuda de cotizaciones y WhatsApp por permiso y versión exacta, compartiendo textos inline;186tests, typecheck/build y4/4navegadores PASS sin reintentos. Composición final67/774 y HTTP68/782; Preflight220164pasos ejecutados/56planes PASS, Docker ausente.8589archivos autenticados y ejecutable HTTP idénticos a evidencia previa. Esta frontera de lectura queda cerrada en referencia;45/48 permanece con TEST02/03/07 bloqueados. CMS, capacitación integral, otras guías y target conservan sus pendientes. No cambios de dependencia, reglas comerciales ni operaciones.

```json
{
  "checkpoint_under_test": 220,
  "executed_steps": 164,
  "all_executed_steps": "PASS",
  "plans": 56,
  "missing_toolchains": [
    "docker"
  ],
  "inventory": {
    "packs": 165,
    "files": 1540,
    "markdown": 827,
    "plans": 56
  },
  "franchise": {
    "packs": 67,
    "files": 774
  },
  "http_reference": {
    "packs": 68,
    "files": 782
  },
  "help_browser_projects": 4,
  "browser_retries": 0,
  "unit_tests_passed": 186,
  "existing_explicit_skip": 1,
  "controls_passed": 45,
  "controls_total": 48,
  "preflight_log_sha256": "e278e5bfce82bf2add09730ff7beb4f1fa23950e8f59e5283fcc1af62608fe41",
  "preflight_json_sha256": "572ff60068031c11872c655178ce3968c592b996f4cca93d8e45d2e52d39df4c",
  "final_parity_sha256": "062ccf8192926efe644c183725797e156c8f69dbc759b7a3f39dd15f6df14066",
  "acceptance_invariance_sha256": "d09615ac7cf9f189fd76a693ea5be031f05ca50a5d0152a225ca310b06798e41"
}
```
