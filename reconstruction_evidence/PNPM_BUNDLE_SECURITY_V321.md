# V321 — pnpm compiled bundle security and reconstruction

Date:2026-09-08. Scope:T2803 maintenance, existing library. Local candidate
gates passed; global admission remains CONDITIONED/BLOCKED. No deployment,
release ZIP, global installation, private data, provider or scheduled scanner.
Staging:%LOCALAPPDATA%/Temp/elite-v321-184b4d5926e44dc5a3cd1ecbb05ff077. Before/logs remain immutable evidence.

## Finding and official resolution

V319's891/891 official pnpm11.19.0 files proved identity, not bundled SCA.
V321 extracts5586 published source headers plus48 physical manifests,471unique
npm packages,0unresolved headers. OSV2.5.1/all-packages finds1package/1advisory:
brace-expansion5.0.8/GHSA-rgw5-rvv9-x895,CVE-2026-69152. Reachability UNKNOWN.
The official advisory fixes5.0.9; destructive memory/CPU PoCs were not run.
Authority:https://github.com/juliangruber/brace-expansion/security/advisories/GHSA-rgw5-rvv9-x895

OFFICIAL_FIXED_RELEASE candidate pnpm11.25.0,29-August release, commit
6d90c71efdffbc909b499490b64c66badc720327; prior536b7a2c8a5b3db8703a81f5d06db7841956d055.
GitHub API reports both commits verified/valid; candidate annotated tag also
verified/valid. Baseline is a lightweight ref, not a verified annotated tag.
Registry gitHead is absent; the commit binding is independently verified via
the SLSA attestation. Release:https://github.com/pnpm/pnpm/releases/tag/v11.25.0
Tarball:https://registry.npmjs.org/pnpm/-/pnpm-11.25.0.tgz
SHA256:33dd0748f27e7916c4f1c8b6943461983e3453b06bbda6312a6280130b4881e5;455files/20297324unpacked bytes.
Integrity:sha512-XN6SW08HX3Jetx+64YpC/+eEUkeJ8ZthxzHLhyHsKKruFg4BqNWvT+2ypCzb8wDv4j2zVrDUoXtNY+EfirfJVg==. Safe path/type/duplicate/size extraction;
exact file manifest retained. MIT LICENSE SHA256:e0a867ff513ea7be2a0ddc339ac6a031e459a38668e077b8f0e649544062f9f2.

Candidate2825headers+27physical manifests yield473packages/0unresolved/0known
OSV advisories, including brace-expansion5.0.9. All472transitive identities match
literal package keys/line receipts in the exact upstream lock; baseline470also
match. The excluded root pnpm identity is proven by its published package.json
and signed attestation, not invented from the workspace lock. Inventory parser
does not implement vulnerability matching. Official OSV performs that check.
Dependency delta:38added identities/36removed; some official bundled versions
decrease, so release number was not treated as evidence of a uniform upgrade.

## Authentication and verification tool scope

Registry ECDSA verifies name@version:integrity with the official npm public key
SHA256:DhQ8wR5APBvFHLF/+Tc+AYvPOdTpcIDqOhxsBHRwC7U. A changed message fails.
Node crypto uses npm's documented format; it does not infer authenticity from
SHA256 alone. https://docs.npmjs.com/about-registry-signatures/

Official Sigstore4.1.1 verifies both npm publish and SLSA DSSE bundles with
transparency thresholds1/1 and TUF refreshed trust material. Exact GitHub issuer,
anchored release-workflow identity, repository/tag/commit, GitHub-hosted builder
and artifact SHA512 are required. Two altered DSSE payloads and a wrong workflow
identity are rejected, plus the ECDSA altered-message control:4negatives.
No private key, signing operation or token was used. TUF cache/seed/tool files
are hash-inventoried; identity policy is not a claim of reproducible upstream build.

The existing Node ZIP Sigstore closure45packages contains vulnerable
brace-expansion5.0.7/ip-address10.2.0 and was not executed as a clean verifier.
An isolated official sigstore4.1.1 tool with50locked packages was installed with
scripts disabled and scanned50/0 before verification. Root commit from registry:
c1dc7d4778a450787fc72b083f2490ad02b714c6; Apache-2.0. This is audit tooling,
not a new product dependency or admitted reusable crypto implementation.

## Consumers, compatibility and licenses

Active owners:TS-GO-API-WEB-BRIDGE0.5.4, Microsoft Playwright gate0.1.24,
Google Lighthouse gate0.1.5. Eight materialized files changed:3package manifests,
Playwright source-lock tool version,2verifiers and2README cache instructions.
All3dependency locks remain byte-identical. Plans:ENTERPRISE_WEB,
FRANCHISE_COMPLETE and FRANCHISE_SERVERLESS. Root preflight availability floor
and routing fixture use11.25.0; availability is still explicitly not admission.
Historical TS-ENTERPRISE-WEB0.2.1 remains unselected/CANDIDATE with its old pin
retained for forensics and an explicit reuse prohibition pending its own gates.

All473fixed registry metadata records have license declarations; all22published
license/notice files and inline bundle notices remain byte-exact. openpgp6.3.1
LGPL-3.0+,next-path1.0.0MPL-2.0,spdx-exceptions2.5.0CC-BY-3.0 are unchanged
baseline identities. @pnpm/npm-lifecycle1100.1.0 is Artistic-2.0. Tool artifacts
are not redistributed by these8pin/config changes. Future tool redistribution
must satisfy these obligations; the parent pnpm MIT label is not their license.
uid-number0.0.6 is unsupported upstream and remains a visible condition.
Published pnpm has no install lifecycle script; msgpackr-extract install script
exists in registry metadata of a bundled dependency, not a newly executed install.
Six native binaries/addons are unchanged byte-for-byte; their transitive native
SCA is not covered by the473npm package scan.

## Executed gates and failure recovery

Node24.20.0 standalone SHA256:5c976096e04e5c2c1f091938926234cc9fbebfe9787ddd149351b3b0ecc707b5.
Go1.26.7 from official-toolchainV281 observed; GOTOOLCHAIN=local.
Candidate and rebuilt: frozen/offline web install,115PASS/1pre-existing SKIP,
Next build. Candidate:92connected browser phases,23per4projects; agenda4PASS;
Playwright runtime4PASS; Lighthouse policy5PASS. No new Lighthouse target-quality
claim, external identity or real provider acceptance. No schema migration:
NONE_WITH_REASON, package manager only and application locks unchanged.

First isolated Playwright offline install failed ERR_PNPM_NO_OFFLINE_META:
3packages linked did not close the policy gate. Frozen online registry metadata
prime, then unchanged-policy offline retry passed. Lighthouse same sequence
passed. README source now documents both bytes and policy metadata in the cache.
No trust, release-age, ignore or risk-acceptance bypass. FAIL524 regression fixed.

Initial root preflight rejected stale routing ID in its fixture (FAIL526);
fixture corrected and46routing checks passed. Final integrated preflight/library
result is recorded after checkpoint validation, not inferred here.
Canonical reconstruction746/746 equals the tested tree after Next generates
the same next-env.d.ts. All8changed files exact; no second browser run was needed
to claim source identity. Go-only business suites retain their previous evidence.

Eight clean project directories, alternating4baseline/4candidate warm-cache
frozen installs: baseline p50=7.193s,
p95/p99=8.635s; candidate p50=7.133s,
p95/p99=8.034s. Maximum RSS KiB:
623340→630596. Predeclared3xelapsed/2xmemory guards pass;
4samples are descriptive, not statistical SLO/load evidence. Same lock after
each install; last alternating pair is local rollback/forward. No target/data
restore, rollout or production canary claim. Old vulnerable tool is retained
only for forensics and isolated known-input comparison, never safe rollback.

## Remaining normative gate

FAIL525: official source profile record still represented V285. Manual V319–V321
acquisitions have hashes/signatures/SCA but no matching materialized source-profile
receipt. Before retained, discrepancy recorded. Do not claim ACQUIRED or invent
past approvals. Next:fix canonical acquisition owner/profile, materialize, answer
exact reversible inputs from instructions/probes and execute governed acquisition.
No100% completion: this gate, native coverage, T2809 operation and target readiness
remain blocked/unfinished. Human D/channel/commercial release/target decisions
are still pending; ARCA remains deferred.

## Receipt hashes

| Receipt under staging | SHA256 |
|---|---|
| baseline-findings.json | 4fe2aa0c268b632f67c594173f240f268d8f3b224f5c9f04cb4b3aaa7f538c99 |
| baseline-osv.json | 4c5d495bc725c062f266c46ba8fc3facf2b1cc3b448de471299a7586faf419a0 |
| bundle-header-proofs.json | f8a18d5e0ef69d3d49522ef4baa3d67847acdf098f2112f38128e35ae16c954b |
| inventory-summary.json | d95419f8bce59c8428f19fda16cfb142a7642ac4923bcbfb84d2af953910b3a4 |
| candidate-artifact-receipt.json | 50681af2ace3c9bb854425b05e09ad9ce028db67a9bb44a0eded51e97491b687 |
| candidate-files.json | 8eeda0dc0e3ec668ebf3cc1de1e77c4005f29154dd5b87e5fa7576b03733b0c8 |
| candidate-findings.json | 9c89933261c4356f3d67a936ed5f3a491cd828aa0d0e3640233bee0c9b1597ae |
| candidate-osv.json | f20f9eb888d7a66d27f54b8341a94f94ab0def70aeb5061a697c441711201449 |
| candidate-header-proofs.json | 5f534aaf6a44e26c7b87340c6b8ec5272f967e8d04f4594bb5d3c2ecd9b217cb |
| source-lock-corroboration.json | 786954bd5fac477ebf3ef5de7487e6cf3f0139d8c0b2ee9d167da3c00268993b |
| authorities/receipts.json | 76e032cd8b27a6d6f0b9cb8a58a4036e7d0f17afd31a24eec0332353abf4e4a5 |
| package-license-receipts.json | d4c151c878b36c1b0c62db5fe754629cf4cafeee4cb07f80d564ad654925b647 |
| signature-authorities/receipts.json | 81c6d973b72a6eee255d8e3a46582047648220691b6e824c645395e7e7da99a9 |
| signature-tool/pnpm-lock.yaml | 05e5728d523256bdd36a60973915b173f6bb9d750705b51cd6f178ebf325a872 |
| signature-osv.json | 449ca6a3f9059dda44da4b40cd1f37cbf9fb0b4c448318d5c691839fa6af51cf |
| signature-verification-receipt.json | 8cffa2275ade9396abd59369e9c327e02d12d7317f041a688943f91e9e68a532 |
| signature-verification.log | 611c8a01b4bdf131b244f8ef40cc4afaefb1861717a4367446de981a937c4dfe |
| final-supplement.json | 0b765241662370c9d47c4953f1808be142b9fa97c3246dbf550fed7528ba2f24 |
| candidate-web-toolchain.json | 23662c4a251c193e8493c4a2334e2c63ad8f1fb971092363239a8aa29be680c9 |
| candidate-web-tests.log | d80de1658c708f9fb56ac56a211b3dea70a754ce9ad0e4ac43e7080f12da3bb6 |
| candidate-web-build.log | f18b494408d2aee055262033461e0d553f07dc0b316b66ecaac072a61f10c623 |
| candidate-browser-install.log | 62aacc384d8694deaddb32537e6b658be53f4267d10ac1699ed9895f30dae9d1 |
| browser-prime.log | 5a4aae7547a6b3216eca8cb3513842d1eb7b58e07044cd5cb2e63c901b9d0b98 |
| browser-offline-recheck.log | 736fbd30c747692ec866db09d40de0e870bede0776d213ce0f318f936aae8f11 |
| playwright-runtime.log | 678ef748376318e91f115b690c4f30c0d6da09302632e2472f53748ec0616881 |
| lighthouse-contract.log | 282ee9aac41e9db4a38394b2632743e9edcb2d966a49da0489961a2ead6e71a8 |
| candidate-browser.log | 781796bb7eafe8c27067ecf71e5a64abfcd0b003f5f6d8e8eb7a5b0c12610b63 |
| candidate-agenda.log | 4bdc76138dc059a27556475c3ae23047ed3d802475320c9391f625d3f3a45b4d |
| rebuilt-web-tests.log | 5341cf1cda810a7973c590df4aff9e0a56e24626e06e9a58c34c42fcd9f46fa7 |
| rebuilt-web-build.log | 2a96c6a5cd48e15459eea1b33f6ff8255c995060bc37e5390bd185f94e9969bd |
| rebuild-parity.json | aabb001712d073f9921959274fab8643042d40edaade3bf541217632b75687bf |
| pm-benchmark-results.json | 08112dc47dd772d839c11babc6ea2537357b75a56e472b4df607ac18b2c48b92 |
| executable-preflight.log | a1f9b16597c5cf9be1ee8f3821fc49605c5caa19bf8a21988c99e1dde78c51df |
| toolchain-routing-green.log | 1ce317d57e585594d93584a0f21cc3bcde024762679e441fc04a319c2cdc679c |

## Integrated result after checkpoint71

VERIFY_LIBRARY_PASS161packs/1447files/756Markdown/52profiles;franchise67/746. All150executed preflight steps PASS. Preflight status is BLOCKED, exit0 informational:dotnet-10+,docker,psql unresolved in this configuration; not evidence of absent software everywhere or project readiness. Exact observed core tools:PowerShell7.6.5,Python3.14.4,Go1.26.7,Node24.20.0,pnpm11.25.0. Full log SHA256:89c9d1a7b62b2ed71449dffce1155d979383620c3a7dac1fbf0fd790acd731b1; executable-preflight.json SHA256:39a46078caee88ada2b9ab5ce819cc59e0d15f6fe61044f7a1a1b82ebf4a37cc. First red log retained. Checkpoint72 closes this evidence round and continues FAIL525.
