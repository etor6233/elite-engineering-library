# Private portal read recovery — V383

2026-09-11. Library maintenance / TEST02 / T2804, parent checkpoint227 validated before resuming.
Root remains DISCOVERY/BLOCKED. No new business, identity, financial or permission policy.

V383 completa recuperación de lecturas admin/customer/factory: error sin detalles internos, rechazo persistente sin permisos y nueva consulta GET tras corregir sesión; ninguna repetición de operaciones.263tests web/1skip previo,27focales,typecheck/build y4navegadores Go/JWKS/PostgreSQL PASS;5negativos HTTP,0escrituras y snapshot6tablas idéntico. Cinco archivos AUTHORED nuevos; producción Go/SQL/dependencias intactas.45/48 mantiene TEST02/03/07 pendientes.

## Implementation and qualification

The V382 production build returned a generic Next error on an insufficient backend bearer.
The four-browser baseline reached the original500 and no domain records, then failed because
no recovery alert existed. Raw failure, traces and screenshots remain outside distribution.
Three segment error.tsx boundaries now use one fixed-text component; they never read the
thrown error, stack, digest, reset or retry callback. A normal fixed-path anchor performs GET
only on explicit activation. It does not replay an operation or infer the result of an earlier
write. The segment may catch a descendant error; recovery intentionally returns to its portal
root. Failures in root layout and external IdP lifecycle require separate handling.

The component matches the admitted local Next16.3.4 ErrorComponent contract. Its three tests
pass synthetic sensitive error properties and verify absent disclosure, no callback invocation,
and exact link/operation-uncertainty wording. Existing permission and amount regressions remain:
27focused tests are a subset of263full tests, not additive. One existing explicit connected
skip remains documented; the connected suite here is executed independently.

Actual production Next/JWE → Go/RS256/JWKS → PostgreSQL,53unchanged migrations, four browser
projects and retries0: each admin/customer/factory route receives an overclaimed encrypted
fixture session and a bearer lacking the endpoint grant. The page shows the fixed error and
no domain records. Clicking with unchanged credentials still fails. Changing to the existing
valid identity, then clicking the same link, issues an observed navigation GET and restores
scoped records. The existing independent-permission, tenant/org/customer-subject, money,
wildcard and removal-on-reload assertions remain intact. Five direct HTTP401/403 negatives,
zero unauthorized admin-only lead queries, zero non-admin admin queries, zero browser/API
writes, and before/after content digest of six durable tables equal. These are permission
rejection scenarios, not injected database outages or general security certification.

Next may stream its segment error under a200 transport response; admission is established by
the actual backend denial and absence of domain data, not by requiring the historical generic
Next500. No backend status or ACL was modified. Mobile recovery screenshot visually reviewed:
message and action remain legible with no layout overflow. All owned fixture processes stopped.
Test TLS/prefetch diagnostics are preserved and do not establish production TLS admission.

FAIL757 is the missing recovery defect. FAIL758 is a test-oracle issue: Next's independent
route announcer also has role=alert. The first fixed run failed in all four browsers; scoping
all alert assertions to the existing main landmark removes ambiguity while retaining exact
text, persistent denial and zero alert after recovery. Both failed stages remain preserved.
FAIL719 recurred on a guessed .md contract filename; the actual .json owner was then used.
Earlier Firefox teardown FAIL423/742 remains open; this run is not its repair.

## Canonical scope and limits

BFF0.5.11/56files adds one component, its test and three error boundaries; existing read docs
are extended. Browser0.1.32/14files extends the existing connected test only. Role0.1.3/7files
updates compatibility docs, with six role source files unchanged. Portals0.14.15/31files changes
current compatibility metadata only: every product block is identical. HTTP0.1.8 updates the
parent composition lock. Four affected plans and five owners rebuilt in absent destinations.
No dependency changes:8589Next/Sharp artifact files exactly match V376 and the pnpm lock is
unchanged. Expected browser package CRLF normalization preserves parsed JSON identity.
All five additions are AUTHORED / LicenseRef-Workspace-Owner; no attribution to official code.

Current inventory165packs/1555files/831Markdown/56plans; franchise68/795 and HTTP69/803.
Full final composition and Preflight receipts follow below. The executable HTTP reference is
required byte-identical to V375 before unchanged runtime evidence can be reused.
This closes this read-recovery boundary. TEST02 still requires broader functional integration;
TEST03 native/vendor security/licensing and TEST07 final signed release remain blocked.
All48acceptance definitions/actions/oracles/statuses are preserved. No complete-control promotion.

## Receipts

Local raw stages resolve through elite-v383-current.txt. Synthetic identities/config/DBs and
diagnostic logs remain outside the distributed packs.

| Evidence | SHA-256 |
|---|---|
| browser-before.mjs | ae1e1c64c70329c7f77d5f5e5873abf56bdeb589f704f35502143bc78cecc5c6 |
| browser-before-alert-scope.mjs | 6c62c5e4ce8a78196c3b3c5cf1ff9f5cd81522f33d0ffa3a7675f3deaeeaff57 |
| web-recovery-qualified-results.json | 2ecd0cfa2f93d0152484c050e2acc84f09da6a9e573639d6657305e6eb366d94 |
| recovery-regression-recovery-qualified.log | 786c587a6636844355f672a8c74d6aac56ae25577cb3fe8a2fe410a1a04808a8 |
| tests-recovery-qualified.log | 10087b9aadee11b00693b5d83aba97571f264b2d3c623c708d6eee1b91733f00 |
| typecheck-recovery-qualified.log | 99d7c2179cdd316467c51f37901025a0753453664d741214ba56fc631fa6c0bc |
| build-recovery-qualified.log | f4693a4c844dc9f79b2edcf1f26a8b9373c0a194981c9916e5de183f752a24a9 |
| artifact-parity.json | 7563f34b2a058f7ce4c0f799cb29ad7b97d517326579e335c400497c8723cb43 |
| manifest-normalization.json | cee473e809d7efd3d3f3dbee690df60fafcc5d6e0c3547e7ad920a6f82cdda55 |
| pending-parity.json | ffacb4b173fd2059a09fd1aab3a8eca2af42ee7b6f5a28b5dbf14a1bedff1f07 |
| canonical-pending-changes.json | e5ca383fc8dd0ecb826a3889df0e7a9edf22f016f3c88b561532fd20d3a0ab3d |
| baseline/result.json | 4c631205af6702f20fb684b6d11240d75d987a3cb06c902b1fd39f97b79c4c29 |
| baseline/admin-reads.log | 39a7e30dbf86c0dc1bf74f405a5a33032a0ea53ff0c5b3cddc950e889b2377e3 |
| alert-oracle/result.json | 7e3956e0164b64dbfa9c6e13f0420dfc467cb0756e190830fd2377236a4cba50 |
| alert-oracle/admin-reads.log | fba0c53d55157ea0e49f10fcfb499f0ba51211def80cd449ab248f4f5ccb0005 |
| final/result.json | 4319b99474fd6c39d9bf8ee915b9f5b1d6dd777b7252ae34eb2cab2cf2b36c64 |
| final/admin-reads.log | dd5c8f7c644a052e2df7d43bd51650da6d99e1d3c4c93b47e408174320a378a3 |

## Final reconstructed consumer

```json
{
  "status": "PASS",
  "packs": 69,
  "files": 803,
  "parent_files": 795,
  "parent_files_identical_to_connected_candidate": 794,
  "next_env_generated_imports": 2,
  "added": [
    "src/platform/backend/private-read-failure.test.ts",
    "src/platform/backend/private-read-failure.tsx",
    "src/app/admin/error.tsx",
    "src/app/customer/error.tsx",
    "src/app/factory/error.tsx"
  ],
  "modified": [
    "docs/private-portal-reads.md",
    "docs/role-workspace.md",
    "reference_http_metrics/source-lock.json",
    "microsoft_playwright_browser_gate/tests/role-workspace.spec.mjs"
  ],
  "http_binary_byte_identical_to_v375": true,
  "http_binary_sha256": "4643654625c71aaf5b45367bd3a198381257588d073ae0349fc90c1bc1b5ab85",
  "lock_sha256": "b77daa20003311c295d496ddbfda1569819121665c6b6d24e8421b90693da59f"
}
```

V383 FAIL759 REGRESSION_PROVEN: composition803and parent795byte checks had passed; a global count substitution corrupted the baseline UUID. Original helper/log preserved, baseline now resolved from recorded tools.json with required materialization record. All final added/modified/byte/hash assertions pass and the rebuilt HTTP executable is identical to V375. No source, permission or acceptance relaxation.

## Closure229 and unchanged acceptance

V383 cierre229: recuperación admin/customer/factory integrada, error fijo sin datos internos y GET explícito; bearer insuficiente permanece denegado y sesión válida recupera datos.263tests/1skip previo,27focales,typecheck/build y4navegadores Go/JWKS/PostgreSQL PASS;0escrituras/snapshot6tablas idéntico.5packs reconstruidos,68/795 y HTTP69/803,8589artefactos/binario HTTP idénticos. Preflight228164pasos/56planes PASS, Docker ausente.45/48; TEST02/03/07 pendientes, sin cambiar oráculos.

```json
{
  "checkpoint_under_test": 228,
  "executed_steps": 164,
  "all_executed_steps": "PASS",
  "plans": 56,
  "missing_toolchains": [
    "docker"
  ],
  "inventory": {
    "packs": 165,
    "files": 1555,
    "markdown": 831,
    "plans": 56
  },
  "franchise": {
    "packs": 68,
    "files": 795
  },
  "http_reference": {
    "packs": 69,
    "files": 803
  },
  "private_read_browser_projects": 4,
  "independent_http_negatives": 5,
  "admin_only_lead_queries": 0,
  "non_admin_admin_queries": 0,
  "browser_and_api_writes": 0,
  "unchanged_domain_snapshot_tables": 6,
  "browser_retries": 0,
  "unit_tests_passed": 263,
  "focused_private_view_amount_and_recovery_tests": 27,
  "existing_explicit_skip": 1,
  "controls_passed": 45,
  "controls_total": 48,
  "preflight_log_sha256": "69b104b563fa5e1912580546a944667ad048c092213feaa4dee7f5fdfa5b1496",
  "preflight_json_sha256": "0faa8f56d414db6f3a03fb4131a0a42b72120b92ee331714f85a9b555b6189af",
  "final_parity_sha256": "37d9935e4e5a51b3dbf856778bd419624b1704ff5091b0307a49acf461af4186",
  "acceptance_invariance_sha256": "d09615ac7cf9f189fd76a693ea5be031f05ca50a5d0152a225ca310b06798e41"
}
```

## Additional TEST03 progress: exact notices and vendored identity

While Preflight228 held the root stable, read-only dependency investigation continued in the
owned evidence stage. No further product source or dependency changes follow that Preflight.
The exact Next root MIT license was captured from the already authenticated16.3.4 artifact
and equals the official license bytes at299180d3315c7ebd7b199d2b1a265b5986c5fc7d. The existing
verified SLSA receipts bind Next, @next/env and @next/swc-win32-x64-msvc to that same commit.
The companion artifact and per-package receipt are now concrete. They have not yet been
integrated into a consumer distribution; they do not resolve vendored Rust/JS/native notices.
Authority: [pinned Next license](https://raw.githubusercontent.com/vercel/next.js/299180d3315c7ebd7b199d2b1a265b5986c5fc7d/license.md).

The authenticated Sharp native README declares29license rows versus28versions.json entries.
All28version names reconcile through an explicit alias map; libnsgif is the only versionless
row. This is an inventory gap, not evidence of a vulnerable or absent runtime component.
Official source inspection resolves its source owner: libvips8.18.6 vendors libnsgif under
libvips/foreign/libnsgif and links its static library into libvips_components. The tag resolves
to GitHub-verified commit426af3f44246fce9cfa8dd51a353aa4dfd48c553; the complete tree is not
truncated and fixes13vendored blob IDs. COPYING, README, meson.build and update.sh were read
at the fixed commit and each byte sequence matches its Git blob hash. Its update script
clones the then-current NetSurf source and applies patches; no independent semantic version
can be inferred from that script or a date. The13file identity is recorded without inventing
a version or claiming correspondence to the installed DLL. No upstream update script ran.
Authority: [pinned vendored tree](https://github.com/libvips/libvips/tree/426af3f44246fce9cfa8dd51a353aa4dfd48c553/libvips/foreign/libnsgif).

The remaining concrete step is to bind Sharp's exact native build recipe and acquired input
artifacts to this vendored tree and the installed DLL, then finish the other component notices,
source/relinking obligations and security mapping. A new upstream source/binary acquisition
group must first get its materialized profile/lock/required-input record. This audit downloaded
public metadata and four small source documents for inspection; it did not acquire/build/adopt
a new runtime group. Existing pnpm native FAIL532, companion FAIL735 distribution conditions,
full TEST03 and TEST07 stay open. The failed web API page lookup was replaced by a successful
direct official metadata read; no unrelated search-result package list was used as authority.

| Additional receipt | SHA-256 |
|---|---|
| notice-analysis/next-license.md | ee765244e2d59f5234d474f62e0766fa0c8b99af967fdd4c0cb8dcb0c76ea224 |
| notice-analysis/next-companion-receipt.json | 676bbdf864517e9261bee6f92bf891036d6b8ac7a3e5661189e63f814f7871ea |
| notice-analysis/sharp-native-reconciliation.json | c41affa4ba9584bfa2939caeef27761a9b235d14279f697b7a162e2178264157 |
| notice-analysis/libnsgif-source-binding.json | a4d1105330c0311e0897d082f52dcbbefc9876d280df8699555158a1ad2bb5bd |
