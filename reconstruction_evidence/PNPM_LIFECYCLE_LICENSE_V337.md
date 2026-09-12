# npm-lifecycle exact package licence and governed acquisition V337

2026-09-08; maintenance T2803/T2808. Resolves the npm-lifecycle source/notice
uncertainty recorded V335-V336 without substituting a moving branch. The package
is LICENSE_VERIFIED for the observed text/applicability; it is not admitted for
runtime, whole-bundle use, redistribution or production.

## Evidence and result

The current official npm registry confirms @pnpm/npm-lifecycle1100.1.0, published
2026-08-27T16:53:56.685Z, Artistic2.0, with no gitHead. Its integrity matches V321.
The tag listing has no1100.1.0 entry. Commit history identifies the exact release
commit da8e0e802ab2e50860fe237c705ca6fb3092e910, timestamp16:53:14Z and manifest
version1100.1.0. GitHub reports its signature verified/valid; this is an observed
GitHub claim, not local cryptographic validation of the Git commit signature.

Fresh official npm keys verify the ECDSA signature over package/version/integrity.
A changed-message control is rejected. A1byte HTTP Range response fixed the full
artifact size10780bytes; no guessed size or fictitious Git SHA was placed in lock.
Opaque acquisition from registry.npmjs.org verifies exact SHA512 and writes4files
(artifact, unchanged candidate licence sidecar, source receipt, profile receipt).
PRESENT replay preserves all4 byte-for-byte. Artifact SHA256:
ee4ae1558da7215a514f65adedc500af8c308ddc4689031c070f0d7722b72527.

Bounded read-only TAR inspection observes8regular files/30261bytes. Archive
LICENSE equals the fixed-commit LICENSE and decoded sidecar byte-for-byte.
Seven files match Git blob identities from that release tree: LICENSE, README,
index.js, lib/extendPath.js, lib/spawn.js and two node-gyp launchers. package.json
differs only by omission of packageManager (source pnpm@11.6.0); name/version/
license and the rest of its parsed fields match. This is an explicit publication
delta, not a claim that8files are byte-identical. Registry gitHead remains absent.

Only licence text is copied as inspection evidence; package JS and launchers are
never installed or executed. Acquisition receipt extracted=false describes the
earlier opaque operation; subsequent bounded member reads have their own receipt.
The library receives only the AUTHORED profile and ADAPTED licence envelope, not
the npm package implementation. Source read-only files remain in temporary evidence.

## Canonical correction and qualification

OFFICIAL-UPSTREAM-ACQUISITION-CORE0.4.79 contains33files/125sources. The existing
transport adds the reviewed scoped npm-lifecycle repository/filename/HTTPS rule;
the version parser now permits1-4major digits and preserves3digit minor/patch
bounds and exact stable syntax. Scope/name/URL spoofing, partial/prerelease/five-
digit major versions and missing commits are rejected. A new separate profile
maintenance-package-license-evidence selects only this source, excluding124others.
The original three-artifact maintenance profile/identities remain unchanged.

78opaque transport checks pass in candidate and independent reconstruction;
upstream acquisition9negative cases and source-profile18valid/7negative cases
pass in rebuild.33candidate/rebuilt files are byte-identical. G0-G8 and nine
assurance dimensions are documented in stage qualification.md. This narrows the
claim to the existing AUTHORED transport under observed Windows conditions, not
the upstream runtime. The actual profile was recorded from its template with all
three required inputs answered under existing maintenance authorization, validated
before acquisition and hash-bound to the approval/lock. No invented human approval.

29composition plans now reference0.4.79; initialization becomes2packs/39files.
Root inventory162packs/1461files/775Markdown/53composition profiles, provenance
AUTHORED1216/ADAPTED138/VERBATIM107.18source-acquisition profiles are a different
inventory from53composition profiles. Franchise67/746 and existing consumer locks
are unchanged. Integrated154-step Preflight follows checkpoint111.

Recovery: FAIL564/565 exposed the3digit-major assumption; actual1100fixture and
negative syntax controls now prove the fix. FAIL567 replaced a hardcoded profile
count with the executed matrix count. FAIL566 path lookup recovered. FAIL568
serializer initially missed2blocks because of blank lines (not differing fence
width); all31old block payloads were then checked against baseline before33file
reconstruction. No failed gate, old event, original receipt or licence rewritten.

## Remaining boundaries

The473component conservative pnpm graph still needs complete notice/use analysis;
the original full-native FAIL532 is not reopened as admitted. Current work resolves
one exact licence provenance gap. It does not certify compatibility of every
bundled work or correspond to a global completion percentage. V336cache policy
qualification remains dated; renew before claiming current freshness. Existing
consumer115tests/build, connected browsers/Lighthouse and43selector/routing tests
retain their original dates unless the integrated gate reruns them.

Readiness42observations,48contract tests/7pending and10macrofronts remain open.
No product expansion, global install, release/signing, provider writes or ARCA.

Primary endpoints retained in receipts: [release commit](https://github.com/pnpm/npm-lifecycle/commit/da8e0e802ab2e50860fe237c705ca6fb3092e910),
[fixed LICENSE](https://raw.githubusercontent.com/pnpm/npm-lifecycle/da8e0e802ab2e50860fe237c705ca6fb3092e910/LICENSE),
[npm metadata](https://registry.npmjs.org/@pnpm%2fnpm-lifecycle),
[npm signing keys](https://registry.npmjs.org/-/npm/v1/keys).
The fixed package artefact and release-tree comparisons govern these observations;
current branch popularity and registry license labels alone do not.

Stage: C:/Users/NL/AppData/Local/Temp/elite-v337-41c13ce2aa624c54b1d83784db789544.
Keep failed attempts; helpers use absent destinations and are not safe to rerun
over existing output. Rollback pack and29composition refs together from before
snapshots, preserving source acquisition receipts and checkpoint history.

| Evidence | SHA256 |
|---|---|
| package-license-inspection.json | `b66c52fe356d2e29993729b2bc44334eed0fefce25c595ecfdcdac5524be4d3b` |
| package-manifest-delta.json | `7b8e3dd3ac4579791125882d793eb48e3764e6466b11f8c8ecf9056c1b30d7f7` |
| registry-signature-receipt.json | `0cdc09c15ac635fb3af672d0717ff8d6f4c8670ffeee3ad987b05e5d72bbf30a` |
| registry-keys.json | `faf23d8753d5bb79df250f10391ac89b63ecf7743e48487a544a99c847f9c8df` |
| release-commit.json | `0220f69c857c7209d6126b481ff41155c82c635e11943b8eb3f1edd63bd6cf0f` |
| release-tree.json | `821aaa35e05fb18453208c5a9c03807227ff277d05b72b1dc5a001433938e0e5` |
| npm-lifecycle-registry.json | `5df0f1bbeb2942c52102a1f647f23cc5f3af46befa5e9a1b6da6290831082228` |
| release-license-receipts.json | `cc7bb2742d810467cc4b1853a73638fc675af118f021fdac5a90da66d8e0b777` |
| tarball-range-metadata.json | `02c9a1b7a620381316052651451de5160c49200b6639dfae2ad20e9294ab86cc` |
| rebuild-parity.json | `2a7fd52c02ccc7a580e71f9edbae03186afeb173e75cca3acd13e74c47c76e36` |
| profile-approval.json | `a9e1ed05b9b5ecc224e8ab22820d3952f50cc642e2573ce24b213fb0c919015e` |
| profile-validation.log | `5ba16d3fe1b16a47044b1dcdd59e27927e9086ac14519934d490471029a07461` |
| acquisition.log | `6cf39354e727e0fffc824a1a1a0ea7a0e7328b53ece6d1f263d553047dd2b38c` |
| acquisition-present.log | `796f2c68e3b9bcd714ba8cb696c874fb44c991fa97107e98b88678728a29540b` |
| acquisition-parity.json | `4bdfdc0e9a6ebc69c621879b2fb9f9adee892f35a59d496df450ff7e004cf730` |
| qualification.md | `c9171f3ee6212e24320ef966b839314500e4d48b5a06b558252df49db598ed4e` |
| test_opaque_artifact_transport-r2.log | `a111c1874239be3c47e734b46479528ea7c10c694b21bad8ee1eaad3b256fc21` |
| test_opaque_artifact_transport-rebuilt.log | `b2f7c6a8684a1a350fd2233441cd53cda30fde48e693b7a013e61043ec621cf8` |
| test_source_profiles-rebuilt.log | `607d424e7e9d25593b88778e709f1a5d01d475e0f77e95fd48d2a3f98c5c9769` |
| test_upstream_acquisition-rebuilt.log | `65fc96b0158e5caf8168af58001d0cd7f02579f0b97076865792942e68fa19b0` |
| governed-acquisition/receipts/maintenance-npm-lifecycle-1100.1.0.json | `539bbd0eb6210872a25ea0369b2dc754b5e9e1ca55826a15eb81a601ecf5164b` |
| governed-acquisition/profile-receipts/maintenance-package-license-evidence.json | `f56840b91f4364387885c594b5c4cf79bd4f48e5fa2094f563c453163686fa3d` |

## V337 composition reference correction

First integration updated29spaced JSON references. Preflight111 correctly rejected one compact JSON reference (FAIL569). The remaining AWS_SECURE_EMAIL_ATTACHMENT_PROCESSING_PACK_PLAN reference is corrected; parsed composition audit now confirms30plans selecting0.4.79. No source payload or acquisition receipt changed. Checkpoint112 precedes the repeated integrated gate; first failure log remains retained.

Reference audit SHA256 3b084a9ecdffc00c94d7ae1a4b6a30193104a410e75850cfd4a338e4d85d53a1; first Preflight log SHA256 fdacce7c74b77062b49abf69b90ca4269e3da8ff94cfc06b46062ad3d6dff7fd.

## V337 file-count correction

Preflight112 rejected a shipping profile count still at54 after two explicit new core files made56 (FAIL570). Audit against all53previous verified profile counts identified six hardcoded Amazon count assertions requiring+2, and ten profile introductions already stale before V337. The six independent expected assertions remain explicit and are corrected; the ten introductions now state current expected counts, preserving before/verified/after evidence. This is not a relaxed count gate. All30core consumers select `*`; the only added paths are the new acquisition profile and licence envelope. Full materialization must confirm every expected count.

Audit SHA256 a91465fbc67dc8467b29605bd86a4e3d4a71bef66614502a85ef8f0a70356c21; failed log SHA256 5d638058581ff295b1bc75750d416445b7c7fa5cf8e00a0bb9afea75bb1d5878. Checkpoint113 precedes another integrated gate.

## V337 integrated closure — checkpoint114

Preflight113 completed on 2026-09-08: all154 executable/structural steps PASS.
Overall availability remains BLOCKED solely because Docker is absent; availability
is not admission. Structural inventory is162packs/1461files/775Markdown and all53
materialized composition counts equal independently recorded expectations.
Initializer is2packs/39files; franchise stays67packs/746files. Core0.4.79 has33files,
125sources and18source profiles; these are not the53composition profiles.

FAIL569 compact JSON reference and FAIL570 count regression are REGRESSION_PROVEN
by the complete integrated replay and the53profile comparison. Failed111/112logs
and their checkpoints remain intact. Rollback encompasses all30composition plans
(the earlier29-reference wording above describes the incomplete first attempt).
No code was changed after this successful gate; checkpoint114 closes documentation.

Additional licence follow-up preserved full fixed-source texts for individual3.0.0
and qrcode-terminal0.12.0, with SHA256 and Gitblob identity checks. individual uses
LICENCE in its fixed tree although legacy metadata points to LICENSE. These two
records remain PINNED_SOURCE_TEXT, without npm archive comparison or redistribution
admission. Legacy declarations were not silently rewritten. OpenPGP LGPL-3.0+
requires use analysis despite texts retained in V336. semver-utils1.1.4 APACHEv2
still needs exact text: discovery could not open the declared source URL, which
does not establish repository unavailability or absence of a licence.

The exact npm-lifecycle1100.1.0 licence gap is resolved within the recorded claim.
Whole473component notice/use compatibility, original full-native FAIL532, Docker,
target access and production gates remain open. Readiness42observations,
48contract tests/7pending and10macrofronts remain unchanged; no global percentage
or production completion follows from154PASS steps.

| Closure evidence (V337 stage) | SHA256 |
|---|---|
| preflight113.json | `9ffe93fa25291d03dce922df8c0f273e50d6fe1b58eb38ebb66cee2a7523d7db` |
| preflight113.log | `abc12c9437b7f2fbf8b4cbaf13f64d9043e358928a913909dc366b77c6e3856f` |
| profile-count-verification.json | `5f73480b09b4868dbb5346cbf7853baab93dd2b26f9a8e56e4fbac914e2c2e37` |
| baseline-delta.json | `a5e49a8a37f0182d2e7c99bbec0567f7f66b8371d1abe044181eb66d73de66a7` |
| legacy-license-followup.json | `2147384e64c31a7771279d59fd2621c584cb42a2a98effbc5bdf8cd36199906a` |
| legacy-source-license-receipts.json | `bceeda622ac80abf1cf0fdb09ffdd5af92baa60e5bd389a8993f898b19d1af07` |
| individual-LICENCE.txt | `19338d17a974d1c747b3f6c618f08975b8908fdebccbc047c0e72b09daf6ddd3` |
| qrcode-terminal-LICENSE.txt | `b3c7a2fadb2515b8106eae58439a4b9c0581a4eaa88d6a265701f8d4dd7dadb8` |
| semver-source-observation.json | `c994d99189106c1ccc107eeb296952bc00bd88a7796d4fb680316a2141ccdf23` |
