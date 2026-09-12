# Complete existing versioned help coverage — V380

2026-09-10. Library maintenance, TEST-02/T2804. Entry checkpoint221 validated against
191evidence references and221events. This completes the existing-guide index boundary selected
by that checkpoint; root product readiness remains DISCOVERY/BLOCKED.

## Verified result

V380 completa el índice de las15/15guías versionadas existentes:13extracciones exactas, permisos por pantalla/sección y enlaces fijados.214tests,13regresiones de HTML inline, typecheck/build y8casos en4navegadores PASS;15perfiles ×15guías ×4navegadores=900decisiones de acceso, más60navegaciones exactas. Tres archivos de operación conservan todos sus handlers idénticos fuera de la ayuda. Sin nuevas dependencias ni reglas.45/48 permanece: no equivale a CMS, guías para tareas que antes carecían de ellas, capacitación integral ni target.

The prior index exposed quote acceptance and WhatsApp status only. The same owner now includes
all13other versioned inline guides found in the selected web source: order operations, lead
changes, quote creation, interval creation/cancellation, resource creation, slot publication,
checklist publication/completion, delivery discrepancy resolution, return operations, section
recovery and customer handover reads. Each retains its exact original id/version, instructions,
paragraph order and inline summary. Current versions include order-operations-view/1.1.0;
all other existing guides remain1.0.0. A pure shared component adds the exact-version link.

Three original source files were independently compared after reversing only the guide
extraction/imports. Every non-help source character is identical, including domain commands,
recovery markers, permission checks and handlers. No new domain rule or business action.

## Permission and functional evidence

The15profile matrix is independent of the production policy and records each expected article.
Every actual Next/JWE HTTPS browser checks all15article endpoints for each profile:900access
decisions across Chromium desktop/mobile, Firefox and WebKit. It also follows all15version links
in each browser:60actual UI navigations. These are checks, not900global controls. Both suites
passed in all four browsers (eight tests, retries0), retaining unauthorized404/403, expired and
adulterated401, exact-version selection,503/malformed response recovery and superseded-read checks.
No browser operation requests occurred; the explicit prior API POST negative still returns405.
The retained desktop image was visually reviewed; mobile overflow checks passed.

Policy mirrors the existing page/section visibility. customer:self receives quote/handover
guidance. Order guidance follows the operation panel's existing inventory/payment/handover/admin
permissions. Quote creation also requires lead:read because that page loads leads under this
permission. availability:read receives cancellation guidance without acquiring manage authority;
availability:manage receives creation guidance. Reading a guide never authorizes its command.
The WhatsApp feature flag removes only its article. Wildcard obeys that flag. Factory-only,
empty permissions and quote:write without lead:read have no existing applicable guides and403.

The bounded response expands from two to15known IDs; only id/version/title/paragraphs are
serialized, never inline metadata, tenant, organization, identity, token or transaction data.
All text is generic PUBLIC material already present in browser bundles. This is relevance
filtering, not confidential-document ACL. No external IdP, provider or database is required for
these static reads. Actual JWT decryption/session expiry runs in the production Next build;
this synthetic fixture does not prove IdP deployment or immediate external revocation.

## Qualification, retained failures and source provenance

214unit/API/render tests pass, including15role cases and13frozen before/after render oracles.
The13render oracles additionally ran as a focused regression. One pre-existing explicitly
separate connected-test skip remains; typecheck and production build pass. Domain Go/SQL,
dependency locks and operation handlers are unchanged, so no unrelated92phase operation replay.
Exact8589Next/Sharp artifact files match V376 and the frozen offline install succeeded with
scripts disabled. Owned Next/proxy processes stopped. No new upstream source or dependency.

FAIL747: the first candidate run had213passes/one failed inline parity check. The oracle used
JSX className="card" instead of its rendered HTML class="card" for the section-recovery guide.
Replacing only that syntax in the original captured source exactly matches the actual digest;
paragraphs, version and class value remain unchanged. No production fix was needed. The initial
test log and oracle are retained, alongside the narrow before/after diagnosis.

FAIL748: strict TypeScript rejected className:undefined in the new render test under
exactOptionalPropertyTypes. The test now omits the absent prop; compiler and component contracts
are unchanged. Focused/full tests, typecheck and build passed after correction. FAIL719 recurred
when a dependent build log was read before checking phase completion; the file was absent and
no mutation occurred. Inspect phase receipts before reading dependent artifacts.

No failure evidence was removed. The earlier Firefox teardown incident FAIL423/742 remains
open; this new green help suite is not its upstream repair. Native license/security/runtime,
Docker, release and consumer-target obligations stay separate.

## Canonical scope and closure boundary

Portal0.14.13 has31files, adding the pure guide renderer,13render regression tests in one file,
and coverage documentation. Browser0.1.29 retains13files and extends its existing help test.
HTTP reference0.1.5 updates only its parent composition lock. Three canonical packs reconstruct
in empty destinations and match the qualified source. Four affected plans select the new versions.
BFF0.5.8 is unchanged. Full67/777 and68/785compositions match the tested candidate; Preflight222 passed164executed steps/56plans, with Docker unavailable.

The existing-guide index is15/15complete. This does not invent help for every task that lacked
a versioned guide, factory instructions, CMS publication/archive, historical storage, confidential
documents, employee training/progress or support policy. These remain explicit gaps in the
general help/onboarding capability. The21-core map retains help as partial general equivalence.
All48control definitions, actions, oracles and statuses remain unchanged;45passed, TEST02/03/07blocked.

## Receipts

Local raw evidence resolves via elite-v380-current.txt and stays outside distribution.
| Evidence | SHA-256 |
|---|---|
| guide-extraction.json | 6779a784f96679f932fd8466af316f2e404a6f3f8c7eb93a493c32c0d22308b7 |
| expected-role-matrix.json | 0d61301167b2af93dbef03088061d18cf0caacb8d7079aee743ed0c0cc350671 |
| handler-parity.json | 6f3578b21e67881b477c25f5c78e15bf1d9dfb442f823232e1410ace661ad8a6 |
| oracle-diagnosis.json | 4ddc280b0bba9c5ec87c5abe87f968c0546d94a7261a88b1c301f827bacbb49d |
| tests.log | 8694de1958c91b981dbb729c78e4d70bbf3025e3367ffa5ac672569b3a58d2b0 |
| typecheck-fixed.log | 4712112eadcb344a9061437f70bac7d48b22febdc510caa2eccbf3cc9b8f3033 |
| web-qualified-results.json | b050d1ab22e6b142148a0864a4fbb5f9ce50ec71a11bf43e849321e970cf128c |
| help-browser/result.json | ece1f0c4c55c66b5f952673c4b65bae2e2aebb7d4c5b160677b7c3ea340eda3d |
| help-browser/help-four-browsers.log | 27ab1b6d8926f076c792935bc1945b00542136a6573aaffea6fa3d2f812e3800 |
| artifact-parity.json | 7563f34b2a058f7ce4c0f799cb29ad7b97d517326579e335c400497c8723cb43 |
| pending-parity.json | a14012a7c4e1cbd7f29b13cb37eb87935d4b4c3fe1af5cffc14b7026e59e2676 |
| canonical-pending-changes.json | 8fe5a571d17423caf4558e452c827681dcc210846d62bef1ebbbea0ac21beee8 |

## Final reconstructed consumer

```json
{
  "status": "PASS",
  "packs": 68,
  "files": 785,
  "parent_files": 777,
  "parent_files_identical_to_connected_candidate": 776,
  "next_env_generated_imports": 2,
  "added": [
    "docs/journey-help-coverage.md",
    "src/components/operational-guide.tsx",
    "src/platform/help/coverage.test.ts"
  ],
  "modified": [
    "docs/journey-help.md",
    "reference_http_metrics/source-lock.json",
    "src/components/franchise-command-panel.tsx",
    "src/platform/help/catalog.ts",
    "src/platform/help/content.ts",
    "src/platform/help/contract.ts",
    "src/app/franchise/page.tsx",
    "src/app/customer/handovers/page.tsx",
    "src/app/api/enterprise/help/route.test.ts",
    "microsoft_playwright_browser_gate/tests/journey-help.spec.mjs"
  ],
  "http_binary_byte_identical_to_v375": true,
  "http_binary_sha256": "4643654625c71aaf5b45367bd3a198381257588d073ae0349fc90c1bc1b5ab85",
  "lock_sha256": "b77daa20003311c295d496ddbfda1569819121665c6b6d24e8421b90693da59f"
}
```

## Closure223 and unchanged acceptance

V380 cierre223: índice15/15de las guías versionadas existentes, permisos por pantalla/sección y enlaces exactos;13extracciones preservan texto/versión/HTML y handlers operativos intactos.214tests, typecheck/build y8casos en4navegadores PASS:900decisiones de acceso y60enlaces UI, sin reintentos. Composición67/777 y HTTP68/785; Preflight222164pasos ejecutados/56planes PASS, Docker ausente. FAIL747/748 corregidos con evidencia;45/48 mantiene TEST02/03/07 bloqueados. El inventario de guías existentes queda cerrado; CMS, tareas sin guía previa, capacitación integral, seguridad/release y target permanecen separados.

```json
{
  "checkpoint_under_test": 222,
  "executed_steps": 164,
  "all_executed_steps": "PASS",
  "plans": 56,
  "missing_toolchains": [
    "docker"
  ],
  "inventory": {
    "packs": 165,
    "files": 1543,
    "markdown": 828,
    "plans": 56
  },
  "franchise": {
    "packs": 67,
    "files": 777
  },
  "http_reference": {
    "packs": 68,
    "files": 785
  },
  "existing_guides_indexed": 15,
  "existing_guides_total": 15,
  "permission_profiles": 15,
  "help_browser_cases": 8,
  "help_browser_projects": 4,
  "browser_direct_access_decisions": 900,
  "browser_exact_version_navigations": 60,
  "browser_retries": 0,
  "unit_tests_passed": 214,
  "inline_render_regressions": 13,
  "existing_explicit_skip": 1,
  "controls_passed": 45,
  "controls_total": 48,
  "preflight_log_sha256": "cd1d5ce5f5ccea533422a32fe507eac7f0c0326b4c9ac97cadf8452f81f89cdb",
  "preflight_json_sha256": "20aa79d81465a6a1dfa24863092313620293d8104cc3f23f6673aed1697739d8",
  "final_parity_sha256": "b4c240b7916c67afd401b47c9255c3e0e138255bd76a4369322a7b777ede9e37",
  "acceptance_invariance_sha256": "d09615ac7cf9f189fd76a693ea5be031f05ca50a5d0152a225ca310b06798e41"
}
```
