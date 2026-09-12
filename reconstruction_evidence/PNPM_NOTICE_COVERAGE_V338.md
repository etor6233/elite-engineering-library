# pnpm notice coverage and nested-work evidence V338

2026-09-08. Library maintenance, T2803/T2808. Continues checkpoint114 and V337.
This is a reproducible inspection of the selected V334 projection and a draft
notice collection. It does not admit the runtime or a redistributable bundle.

## Measured scope

All473 V321 registry metadata files still match their recorded hashes. The exact
442-file projection contains2680 bundle source markers representing456 package
identities.21 external package manifests plus the root pnpm manifest produce468
distinct observed identities. All are in the conservative473-row inventory.
The remaining five are @reflink/reflink and its darwin-arm64, darwin-x64,
win32-arm64-msvc and win32-x64-msvc packages, all0.1.19. They are not found by this
probe; this does not independently prove runtime unreachability. The original
full-native FAIL532 remains blocked and its473-row inventory is not rewritten.

22 existing notice files were enumerated and hashed:17 direct files and five
nested notices.21 embedded attribution paths name19 packages; some describe
additional vendored works. Across all473 inventory rows,17 have a direct notice,
18 have only matched embedded attribution and438 have neither. Among the468
observed identities,451 have no direct notice file. These categories do not
establish451 infringements or complete licence coverage: other source headers,
permitted links and applicability still require review. A filename or package
declaration alone cannot certify the scope of a notice.

Independent reconciliation uses path splitting and manifest enumeration rather
than the probe's marker regex. It agrees on468 identities and2680 lines. Seven
parser boundary checks pass, including scoped npm-lifecycle1100.1.0, unscoped
identities and three mismatched module paths. All442 payload file hashes remain
unchanged. The draft collection preserves32 text copies byte-for-byte with a
manifest; no executable, original receipt or package has been changed.

## Concrete findings and prepared evidence

Five external packages — chownr3.0.0, isexe4.0.0, minipass7.1.3, tar7.5.22 and
yallist5.0.0 — declare BlueOak-1.0.0 and have no direct licence file. No exact
blueoakcouncil.org/license/1.0.0 URL occurs in the442 inspected files. The official
licence text is now retained in the draft. Its Notices clause requires recipients
to receive either the text or its URL; declaration labels are not automatically
treated as that proof. [Official Blue Oak text](https://blueoakcouncil.org/license/1.0.0).

qrcode-terminal0.12.0 has an Apache declaration and fixed root licence, but its
vendor/QRCode/index.js header declares MIT and attributes Kazuhiko Arase. The
fixed source file matches Git blob10eb8eb0a06aa50d0d3c508f886a7728f52e5d98 in
commit90f66cf5c6b10bcb4358df96a9580f9eb383307b. Ten vendor paths occur in the
selected bundle. The author string was not found in the442 files. The source
header is retained separately; this is source inspection, not proof that every
vendored byte equals the npm archive. Never normalize the entire package to
Apache while dropping its embedded MIT work. [Fixed source header](https://raw.githubusercontent.com/gtanner/qrcode-terminal/90f66cf5c6b10bcb4358df96a9580f9eb383307b/vendor/QRCode/index.js).

The embedded notice list also names noble-ciphers/noble-hashes/noble-curves inside
OpenPGP, formdata-polyfill/ws inside undici, Node.js contributors and Blake Embrey
inside Yarn PnP. node-gyp includes gyp and Python packaging licence files. These
are explicit nested-work follow-ups; the npm dependency graph alone cannot close
their versions, source obligations and attribution.

## Use and redistribution conditions

These are evidence-backed review requirements, not a legal or runtime approval:

| Component | Required continuation |
|---|---|
| OpenPGP6.3.1 / LGPL-3.0+ declaration | Full LGPL and GPL texts are retained. Assess the actual combined bundle, notices and a demonstrable source/relinking or suitable shared-library route before distribution; embedded one-line attribution is insufficient proof. Preserve the legacy declaration until a justified successor exists. |
| next-path1.0.0 / MPL-2.0 | Retain licence and applicable source; determine preferred source form and provide recipient access when distributing executable form. Internal use and delivery of a copy have different triggers. |
| npm-lifecycle1100.1.0 / Artistic-2.0 | V337 resolved exact licence identity, not the bundle's transformation/distribution conditions. Record standard/modified and source/compiled treatment, changes and source instructions applicable to the delivered artefact. |
| spdx-exceptions2.5.0 / CC-BY-3.0 | Preserve source README attribution, licence URI and changes where applicable. Its author's statement about mechanical extraction is retained without assuming a jurisdiction-specific exemption. |
| individual3.0.0, qrcode-terminal0.12.0 | V337 pinned texts exist; individual uses LICENCE despite legacy metadata naming LICENSE. Complete package/source binding and nested notices. |
| semver-utils1.1.4 / APACHEv2 | Exact text unresolved. Maintainer GitHub repo confirms relocation; fixed relocated source could not be opened by discovery. This is not proof of source absence. Use a governed exact archive or restored primary access; no mirror or current-branch substitution. |

Sources: [OpenPGP fixed LGPL](https://raw.githubusercontent.com/openpgpjs/openpgpjs/2ac0048404b74a3595d503125b53f3b3d0486bec/LICENSE),
[MPL sections3.1–3.4](https://www.mozilla.org/en-US/MPL/2.0/),
[Mozilla use/distribution FAQ](https://www.mozilla.org/en-US/MPL/2.0/FAQ/),
[Artistic2.0 sections1–7](https://perlfoundation.org/artistic-license-20.html),
[CC-BY3.0 section4](https://creativecommons.org/licenses/by/3.0/legalcode.en),
[semver-utils maintainer relocation](https://github.com/coolaj86/semver-utils).

## Continuity, recovery and boundaries

FAIL571 was a research-probe root-union omission: pnpm was incorrectly labelled
not found despite package.json/LICENSE. Initial script/output are preserved;
root-inclusive independent reconciliation proves the correction. FAIL572 was a
nonexistent guessed inventory filename; the actual owner is README.md from the
execution-state reference. No runtime or package was promoted on either result.

README's stale161/1453/52 summary is corrected to the V337 demonstrated162/1461/53;
its email profile142 becomes150, verified against all53 V337 composition counts.
Current documentary inventory becomes776Markdown because this report adds one
file; no materializable implementation file, provenance count or pack version
changes. Preflight113's154PASS steps retain their original timestamp. A structural
gate follows checkpoint115; production/readiness and10macrofronts stay open.

Next: governed exact semver-utils archive/text; verify remaining direct notices
and nested-work identities; complete source/relinking/attribution evidence for
the chosen use. The32-text collection remains DRAFT_EVIDENCE_ONLY. No external
distribution, new runtime execution, provider effect, global install or ARCA.

## Evidence

Stage: C:/Users/NL/AppData/Local/Temp/elite-v338-606f40d5d434468ba8ceb8d424340173.
The two research scripts are local probes, not admitted reusable product code.
Their independent checks verify these observations without shipping a new pack.

| File in stage | SHA256 |
|---|---|
| audit_notice_evidence.py | `07ef2ce395a61e82e1c9d2ebd4b7f6934dad67d28d45a9dafc51a3db754e4453` |
| verify_and_collect.py | `be7e192fed71fc5fdb576a7a1429f781ac130c6bcf5e0a340394cc439966b9e2` |
| notice-audit-summary.json | `be06b8fe9f59cfc84073563e4845eb8331d4e59151e98013ec7433dbe424dd72` |
| package-notice-matrix.json | `c5e1a953df4473169a51960bd77147e405ab850b22919f7f610728eed29d5a4a` |
| payload-manifest.json | `4cd6f11da178240bb7d8a925df68bd7c9c2874e7508b746586acd8bfecd52cc1` |
| notice-file-inventory.json | `4869f52822b1f7c3ad5db4c1100bc9ee1a5fde8a6d15eb9c889e38018f751211` |
| embedded-notice-inventory.json | `d9f345e1f7070c754d222c0caef30b925b80f35d7df1745f403127a8780c579c` |
| notice-verification.json | `c3b4fa9b7e69815bf8f9f6594244750a390215a4542b8e73628a2c027d003e46` |
| authority-receipts.json | `933a3cb6ca9347bae3e4299d7eeac3631b3cbd57fb48109aab0be31e322e8007` |
| qrcode-vendor-inspection.json | `90cbe3b87ba2ae1e5a6f11a35a88bc4a65010a0da70d43b069211a31ed7d27c3` |
| draft-notices/MANIFEST.json | `8e938692bf4d251cfc8ade5185c700508199e44cac2827cbb33de41d9571fd40` |
| draft-notices/README.md | `5ac2b212f96b6a6dbdd675aa46212e6d3013f4ed5fdc563d5d65189c719bc96c` |

## Structural gate correction after115

VERIFY_LIBRARY115 rejected duplicate ledger IDs (FAIL573): post-gate closure114 repeated two canonical identifiers, and V338 repeated one. The full Preflight113 ran before those closure edits, so its PASS never covered their final ledger formatting. Resume114/115 proved hashes and event continuity only. Three successor references now use plain text; original canonical rows/history are preserved. Exact gate regex reports 2499unique IDs. A second probe assumed162directory Markdown files (FAIL574), but the directory is162pack definitions plus README: all163hashes are unchanged. Structural gate will replay after checkpoint116.

Evidence: ledger-uniqueness.json SHA256 f6ce8700e03333de939c8ca4321f902666d12d21d06c68344cd01667d6606a43; pack-immutability.json 2d97a0496b9930a794f23ffca2fb5c31febaf6fc0856535607577f09aae214ca; failed verify-library115.log c4a39500eb38702b2d30e1b6f7b95eaac2d500841451412046b351ff98b0f96f.

## Closure117

VERIFY_LIBRARY116 PASS:162packs/1461materialized files/776Markdown. All53actual composition counts exactly match the V337 baseline. The162pack definitions plus their README (163files) remain byte-identical; all442selected payload files remain unchanged. Independent32copy draft verification PASS. FAIL573 ledger formatting is REGRESSION_PROVEN by the full structural replay and an exact uniqueness check repeated after closure; FAIL574 classification is corrected. Preflight113154steps retain their original date; no new executable gate or runtime admission is claimed.

| Closure evidence | SHA256 |
|---|---|
| structural-closure.json | `6fa70a64bdb11247129e1ec88c4079b4679b0dd8079e4bd4cc7a457987af6d16` |
| verify-library116.log | `2b2219b6655696ee463abdd9220b27dc210ad53a7a3336cee1d169e6ac640c9f` |
| pack-immutability.json | `2d97a0496b9930a794f23ffca2fb5c31febaf6fc0856535607577f09aae214ca` |
| draft-independent-check.json | `7653ad1d5134e960ed9492104f48cd29105ca9f8aa2c8e6a26f2dcd7881d771f` |
| source-discovery-observations.json | `f689f71abe2219642b0d04bf9c9f404ad6c5986d07a211d2574ce5d65a827da8` |
