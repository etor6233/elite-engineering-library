# Pinned pnpm consumer routing V335

2026-09-08; library maintenance, continuing checkpoint106. PNPM-ARTIFACT-
SELECTION-GATE0.2.0 adds2AUTHORED files to the existing owner (six total).
The original selector, selection policy and17projection tests are byte-unchanged.
New prepare/verify tooling fixes3consumer manifests/locks, Node24.20.0 bytes,
projection/acquisition binding, fresh consumer directories, separate store and
metadata cache, fixed offline/copy/no-scripts argv and an isolated replacement
environment. Plan/config files publish exclusively to an absent directory;
expected plan SHA256 is kept outside the bundle. No runtime is executed by the
pack; no runtime/license/redistribution admission is granted by its PASS.

43tests (17projection+26routing) pass in source and independent six-file rebuild.
New cases cover unknown profile, manifest/lock drift, inherited hooks/npmrc,
Node/acquisition drift, overlaps, existing installations, missing cache, tamper,
configuration changes, reparse, concurrent/late destinations and rename failure.

Actual evaluation from final rebuilt planner:3fresh consumer installations
(enterprise-web, Playwright, Lighthouse) PASS, frozen/offline/ignore-scripts,
zero downloads and unchanged manifests/locks. Source evaluation also passed3.
These are isolated maintenance probes using already acquired candidate bytes,
not canonical consumer promotion or source acquisition. Final web consumer:
115tests PASS/1existing skip, Next16.3.2 build PASS.745sources unchanged; only the
same V333Next-generated next-env.d.ts delta remains. Browser4/full Lighthouse5
from V334 are retained, not represented as rerun in this iteration.

FAIL555: pnpm rejects bare userconfig/globalconfig; config-prefixed syntax fixed.
FAIL557: store contents do not replace metadata cache. Dropping the exact known
enterprise workspace invalidated its cached supply-chain policy. Preserve that
workspace and use ignore-workspace only for standalone consumers. No fabricated
cache verification, trust-lockfile, policy relaxation or fallback download.
FAIL556 retry-helper naming and558 auxiliary path read recovered separately.

G0-G8 resolution USE_REUSABLE_PACK is only for this AUTHORED no-execution tooling
in observed maintenance conditions. Nine assurance dimensions and reproducible
commands are in stage qualification.md. Full license/use scope, cache trust and
freshness, selected runtime admission and later project gates remain separate.
Original native FAIL532 stays blocked. Readiness42observations,48contract tests
with7pending, BENCH01 and10macrofronts remain open; no exact global percentage.

Integrated154-step Preflight follows checkpoint107. Inventory expected162packs,
1459files/773Markdown/53profiles; standalone selection profile1/6, franchise67/746.
No new product code, global install, release/signature, provider effects or ARCA.

Stage: C:/Users/NL/AppData/Local/Temp/elite-v335-52739a31121c471a94ef0e420809407c.
Failed retries are preserved; do not rerun mutation helpers into existing targets.

| Evidence | SHA256 |
|---|---|
| all-tests.log | `91b483954d152b9d9ff3024953fd0d6f7fa7bac86fafc14b99ab8bd2e39abbf2` |
| rebuilt-tests.log | `dccdc830a43880a298db1f811b6ce536c01ceabb5abf2d35506f2b845e36cda8` |
| parity.json | `760fff9fab761c0e77e70c0ba9601ad9c78ce901855ee1ec3ea37a7bddf14132` |
| consumer-r4-results.json | `739deac06ac4b303ae8f0514f1bdf1ff71d6b405d7800055e1a14c635dadad5d` |
| consumer-r5-results.json | `102c45166ecddc023713175821961cdc33bd7996d8df8930495ba991ebafb7f1` |
| post-install-results.json | `706d831c65cbad8f2afc95dc83f55b27bb032e88382d62f5f69f0e98a5487851` |
| web-tests.log | `55d9f0de1dbd4822ea09821ffae177fe8f8e9c850db763823799409d1ef6fc6b` |
| web-build.log | `1201efe94118451a0f99e2de4847e29e37a129aa62d57e281836a00939eb8025` |
| enterprise-web-install.log | `5fe1a5e1f8fec1ba367c3afb940971783c4365552d6660287b578d4d78741f4d` |
| enterprise-web-r2-install.log | `846ef28bb033b2e5c8a087c2b6fbedd1bb54b6477d40e6060d3bc7cafa9b8be2` |
| enterprise-web-r3-install.log | `399d07645210c402c07eee9406fbc74e35d26c0917d2fed328a4df776615fede` |
| gap-record.json | `81335cf507aec94a66861ad361d07dd20e31312887b7bfcd912398f365edd8b3` |
| gap-receipt.json | `1abcce22399772a430f67a72e8b1b68ce4030a6abcd43a9798be64a44a4f19dd` |
| qualification.md | `6b8b0c120b761a73ac301fcd01d762ce0833e6da4a3ef7e7b11ac78559366f18` |
| enterprise-web-r5-install.log | `2c8a472c21b477f751fe0a0f3ae5eb83b73a2e0494d1680336842f642fe23386` |
| playwright-r5-install.log | `18f36937f4f41f4bba4538745c3a993e27f0fcf1ac37d20f457722fe87734a85` |
| lighthouse-r5-install.log | `55646549769d7978f2a3c9ced5e487df5dc0d83b9a67384d448c2605007b1703` |

## Integrated closure

Preflight107 completed154/154steps PASS, including17projection tests and the new
26routing regressions. Inventory162packs/1459files/773Markdown/53profiles and
franchise67/746 remain consistent. Availability status is BLOCKED only for Docker;
this does not waive project access/readiness. No canonical executable code changed
after this full gate. Checkpoint108 records the closure and preserves107events.

The recipe's executed=false describes its preparation operation. Separate
evaluation logs record the later isolated installs; there is no claim that
installation never occurred, or that those probes admit the runtime.

The license review now maps all473existing registry metadata records to their
declarations, with zero missing labels. The four special families are concretely
@pnpm/npm-lifecycle1100.1.0 (Artistic-2.0), next-path1.0.0 (MPL-2.0),
openpgp6.3.1 (LGPL-3.0+) and spdx-exceptions2.5.0 (CC-BY-3.0). Their paths
appear in the retained bundle. Three fixed repository revisions were observed;
npm-lifecycle metadata lacks gitHead, so the current branch is not substituted
for its package-applicable source/license revision. SPDX upstream redirects to
jslicense at the same commit; original metadata/history remain unchanged.
This is declaration/source discovery, not a selected-bundle SBOM or full license
clearance. Primary URLs and boundaries are in license-primary-observations.md.

Remaining next work: exact license texts/notices and internal-use versus
redistribution scope, cache trust/freshness and final selected-runtime admission.
Do not expand the canonical recipe to arbitrary projects or commands, reuse a
used home as pristine, bypass supply-chain policy, or present historical consumer
tests as new ones. Other independent roadmap owners may continue in parallel
with scoped human/external blockers; historicalD is not a blanket pause.

| Closure evidence | SHA256 |
|---|---|
| preflight107.json | `83ca0782ece00ae29b9aab095fedc922aabe44bfcbaf7a6a6e840bc04b74331e` |
| preflight107.log | `2e7b4af83802265530a24325bf9b09ac0cb26c8a1a1a62d59955a41fb06f76c3` |
| declared-license-inventory.json | `7ccd0f0148dc19f350986289dbb88520d5d826acc4698b9c93113815e7d55ff5` |
| license-primary-observations.md | `02b5dbbe2448317ec034e0ffc03f77a29fce75274b9a31ae685b53b583a89034` |
