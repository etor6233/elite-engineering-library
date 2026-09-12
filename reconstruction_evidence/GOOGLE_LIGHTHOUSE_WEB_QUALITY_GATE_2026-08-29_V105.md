# Google Lighthouse Web Quality Gate — reconstruction V105

Date: `2026-08-29`

## Result

`PASS`, conditioned to the limits stated below. `GOOGLE-LIGHTHOUSE-WEB-QUALITY-GATE 0.1.4` reconstructs ten files, installs the exact official Google Lighthouse runtime offline from a frozen pnpm lock, and executes five governed audits against the enterprise web profile.

This evidence does not relabel authored policy or orchestration as Google code. Provenance is nine `AUTHORED`, zero `ADAPTED`, and one `VERBATIM` Apache-2.0 license. The executable Lighthouse and `chrome-launcher` packages are acquired by their exact official npm identities at installation time.

## Official identities

- Repository: `GoogleChrome/lighthouse`.
- Release/tag: `v13.4.1`; release id `356969559`; published `2026-07-20T20:10:47Z`.
- Commit: `1d58f5b06d28e3419b38817a6c7488ec4413c67d`.
- Tree: `405d16b2cdaac5c3ba3823b720031386a618165d`.
- Exact commit archive: `73,569,399` bytes; SHA-256 `f526c6719e496519fc85986226364ef4e6274b1b0e2d8f64a5b8e22fc51510e3`.
- Official license: Apache-2.0, `11,358` bytes; SHA-256 `cfc7749b96f63bd31c3c42b5c471bf756814053e847c10f3eb003417bc523d30`.
- Official package manifest SHA-256: `9a2acb5f3c0908842aacec66eb0bdd8307edd73469098e7eb186b5b79dfad1cf`.
- npm `lighthouse@13.4.1`: `3,562,319` bytes; SHA-256 `110759ba9e863c024e214e9b08ed2b0344d89b492286227235d4dfb990dc3e54`; integrity `sha512-fDu8lt3QLK/lTqIxtp1HkzQNJ32rsFHhbadYOepcMZFLgA8oINhxutMbMv8XXnpTOvZ0TXCo4JCk1LDTWaRLnA==`.
- npm attestations: `14,734` bytes; SHA-256 `b1bd42f081b45ae2c6929be8084452a6ddb6f8ac47832e869862fe6adc0641f7`; two statements bind the package digest.

The GitHub API reports the release commit as unsigned. It is therefore not presented as signed. `UP-FAIL-186` retains that condition; the pack compensates with exact commit/tree/archive/license/manifest identities plus npm integrity, signature metadata and attestations, but remains `CONDITIONED`.

## Reconstructed contract

- Ten manifest paths equal ten `CREATE` blocks and all embedded SHA-256 values verify.
- Dependency graph: Lighthouse `13.4.1`, chrome-launcher `1.2.1`, 119 resolved production packages.
- Runtime: Node `24.14.1`, pnpm `11.19.0`, PowerShell 7.
- Browser: Microsoft Playwright Chromium `151.0.7922.34`, revision `1234`; Windows executable SHA-256 `409805a16d6416087e6b2f778df1cf8f7bbb267d6b99f6b5bb0a618eace234f2`.
- `pnpm audit --prod --audit-level=low`: zero vulnerabilities on the dated graph.
- Five policy tests passed. Negative cases cover URL credentials, remote HTTP, fragments, missing categories, even run counts, version ranges/drift, warnings, origin escape, insufficient audits, weak individual runs and weak medians.

## Real target evidence

The enterprise web profile reconstructed 63 files, installed from frozen locks, passed TypeScript generation/typecheck, 14 Vitest tests and a Next.js production build. Microsoft Playwright then passed four isolated runtime tests and four target journeys across Chromium desktop/mobile, Firefox desktop and WebKit desktop.

Google Lighthouse ran five times against the built loopback public home. Final Foundation scores:

| Category | Median | Minimum | Policy |
|---|---:|---:|---:|
| Performance | `0.97` | at least `0.85` | median at least `0.90` |
| Accessibility (automated) | `1.00` | `1.00` | both `1.00` |
| Best practices | `1.00` | at least `0.90` | both at least `0.90` |
| SEO | `1.00` | `1.00` | both `1.00` |

All reports had zero runtime warnings, stayed on the requested origin and contained at least 140 audits. Raw reports and the summary were written exclusively to the requested evidence directory.

## Windows long-path regression

The composed module exposed a real Node/Windows boundary defect: at the long Foundation path, Lighthouse's nested `lib/cdt/package.json` was no longer honored and its CommonJS `SDK.js` was interpreted as ESM. Isolated bytes were identical; the relevant paths differed only in length.

Junction, Node symlink-preservation and `subst` attempts were each rejected after proving they either returned to the physical path or changed pnpm module resolution. The admitted correction copies only the exact package manifest, lock, policy, runner and validator into a uniquely named short temporary capsule, performs `pnpm install --frozen-lockfile --offline`, runs with standard Node semantics and deletes the validated capsule with bounded retries in `finally`.

The regression passed first from a 138-character composed scratch root and then from the canonical Foundation reconstruction. Five runs completed, the command exited zero, no `elite-lh-run-*` capsule, `elite-lh-link-*` junction, substituted drive, Chrome profile or owned server remained.

## Gates

- `VERIFY_LIBRARY_PASS`: 73 packs, 702 materializable files, enterprise web 5 packs/63 files.
- `VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit`: 73 packs, 113 official upstream source records, one Playwright gate and one Lighthouse gate.
- `VERIFY_EXECUTABLE_LIBRARY_PASS mode=Foundation backend=go-test web=typecheck-test-build browser=playwright-chromium-firefox-webkit quality=lighthouse-five-run`.
- Foundation used the exact official Go `1.26.7` temporary distribution; no global toolchain was installed.

## Failure learning retained

`LIB-FAIL-1136` through `LIB-FAIL-1160` and `UP-FAIL-186` remain in `markdown_system/LIBRARY_FAILURE_LEARNING_LEDGER.md`. Failed runs were never converted into success receipts. The ledger records search truncation, wrong archive identity, cleanup errors, missing toolchains/cache, path-length diagnosis, ineffective junction/symlink/substituted-drive attempts, cwd dependence, the final offline-capsule regression and idempotent temporal cleanup.

## Explicit non-claims

Lighthouse is automated lab evidence. It does not prove assistive-technology interoperability, real-user monitoring, load capacity, authenticated-role journeys, offensive security, provider production behavior, PostgreSQL integration/recovery, cloud deployment, canary/rollback or business acceptance. Those remain project/environment gates and cannot be inferred from these scores.
