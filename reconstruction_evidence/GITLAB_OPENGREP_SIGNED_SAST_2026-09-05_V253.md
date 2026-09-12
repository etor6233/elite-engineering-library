# GitLab/OpenGrep Signed SAST — reconstruction evidence V253

## Decision

`GITLAB-OPENGREP-SIGNED-SAST-GATE 0.1.0` is admitted as `ELITE_REFERENCE / REBUILD_VERIFIED / CONDITIONED` for a narrow Windows x64 Go/TypeScript local SAST claim. Five materialized files are `AUTHORED` acquisition/execution glue. The installed OpenGrep binary and three GitLab rule files remain exact external bytes under their own licenses; no GitLab/OpenGrep authorship is claimed for the composition.

## Fixed authorities and artifacts

- OpenGrep `v1.29.0`, verified commit `344509d693c852eaac4fc1eeffaf2f655c531b5a`, tree `ccc75c2c86a22c3f1670a21c9c7fae1d6728259b`; Windows asset 53,536,256 bytes, SHA-256 `ee485b31912704dc6410bc43f04b5c6ad896697db56e360a98204abf95fa1025`.
- OpenGrep certificate/signature SHA-256 `eb7906ca…` / `64274240…`; official identity regex `https://github.com/opengrep/opengrep.+`, issuer `https://token.actions.githubusercontent.com`.
- Sigstore Cosign `v3.1.3`, signed tag object `2f3a85b…`, verified commit `11926fa5…`; Windows asset 198,819,314 bytes, SHA-256 `9fe59be0…`.
- GitLab SAST Rules `v2.9.3`, commit `39fc7de…`; official package 272,617 bytes, SHA-256 `abbff567…`.
- GitLab Semgrep analyzer authority commit `b8d1635…`, Dockerfile blob `ee66eca…`, fixes `SCANNER_VERSION=1.174.0` and `SAST_RULES_VERSION=2.9.3`. This proves the compatible official rules revision, not identity of this OpenGrep adaptation with GitLab's analyzer image.

## License admission

Selected exact rules: `eslint.yml` 11/MIT, `gosec.yml` 28/MIT, `lgpl/nodejs_scan.yml` 83/LGPL-3.0-only; total 122, with 117 at `WARNING|ERROR`. OpenGrep is LGPL-2.1-only; Cosign is Apache-2.0 build tooling. All exact license texts are retained.

Rejected from runtime:

- `dist/gitlab/**`: exact included license is GitLab Enterprise Edition and requires subscription for production.
- `dist/lgpl-cc/**`: exact included license adds Commons Clause and restricts sale.

The runtime contains neither rejected family. Public availability was not treated as commercial permission.

## Executed gates

- Release asset digests/lengths: PASS.
- OpenGrep version: `1.29.0` PASS.
- Cosign version: `v3.1.3` PASS.
- Official Cosign verification of the OpenGrep blob: `Verified OK`.
- GitHub OpenGrep commit and Cosign annotated tag/commit identities/signatures: PASS.
- GitLab rules tag and analyzer Dockerfile blob/content pins: PASS.
- ZIP inventory/traversal, rule bytes/counts and license evidence: PASS.
- Hostile TypeScript: `eslint.detect-eval-with-expression` detected.
- Hostile Go: `gosec.G204-1` detected.
- Clean TypeScript/Go: zero findings/errors.
- JSON/SARIF finding-count parity: PASS.
- Runtime manifest tamper: rejected before evidence publication.
- Materialization: 5/5, byte-identical round-trip.

Final focal output: `GITLAB_OPENGREP_SAST_PACK PASS static=1 negative=2 signature=1 authority=1 hostile=2 clean=0 json=1 sarif=1`.

## AWS ASH candidate disposition

AWS ASH `v3.7.1` was investigated, not silently adopted. Exact tag commit `af1206aff596f1288ad066828058b3011de2a998` is unsigned. Its documented dev sync and default suite on this Windows host produced 4,556 PASS, 156 skip and 26 fail: six symlink/privilege cases and twenty CDK-extra cases. With all extras, the focal CDK suite failed 13/53 under default xdist but passed 53/53 sequentially; the alternate run diagnoses but does not replace the official gate. Google OSV Scanner 2.5.1 found `GHSA-xvg9-69gf-fjrf` in dev-only `mkdocs-material 9.7.6`; the separately exported 83-package runtime graph was clean. Docker was absent, so no container claim exists.

ASH did successfully orchestrate the fixed OpenGrep/rules lane: hostile corpus exit `2` with two findings and clean corpus exit `0`. Because the full upstream source remains conditioned by the unsigned commit, dev advisory and default-suite failures, it was not published as a reusable ASH runtime. The smaller direct gate avoids that extra startup/dependency surface while preserving the exact external engine/rules and evidence.

## Remaining project conditions

This result does not certify a franchise or any target product. Each project must still run the gate on its real source and triage every finding; suppressions need owner, justification and expiry. SCA, secrets, DevSkim, fuzzing, DAST, threat modeling, browser/authorization journeys, load/resilience, offensive security, deployment/recovery and business acceptance remain independent gates.
