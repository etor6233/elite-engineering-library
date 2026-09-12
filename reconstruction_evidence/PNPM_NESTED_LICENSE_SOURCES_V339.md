# pnpm nested licence sources V339

2026-09-08. Library maintenance T2803/T2808; resume117 verified before work.
This extends the V338 notice dossier with exact source metadata and licence
texts. No runtime, package archive or redistributable bundle is admitted.

## Results

29 bounded official registry/repository metadata and document responses are
preserved with hashes. Eight source package manifests match expected identity,
version and their fixed Git blobs. Seven full licence texts are retained with
fixed-tree blob checks. All442 selected pnpm payload hashes remain unchanged.

The new draft has39 byte-verified text copies: all32 V338 entries plus seven
source licences. It has its own manifest linked to the previous manifest and
passed a separate39-copy hash verification. No original dossier was overwritten.

### BlueOak exact source evidence

| Package | Fixed commit | Result |
|---|---|---|
| chownr3.0.0 | 8b9800ac5fe4da0b58bffc9c66dd618f0721472d | Exact source manifest declares BlueOak. Complete fixed tree has no named licence file; README is usage-only. Declaration/source-text distinction remains open. |
| isexe4.0.0 | 2e7df7dabc4f68e88cf6b32b9029225ba52c6b0c | LICENSE.md Git blob verified; bytes equal the official steward text retained V338. |
| minipass7.1.3 | ab4b3b05d0d557ac6bb178f38501b11d0c96454e | LICENSE.md Git blob verified; bytes equal the steward text. |
| tar7.5.22 | 2a22bfc5d3a432a606d9da0e2d87ba634aa3b1cb | LICENSE.md Git blob verified; bytes equal the steward text. |
| yallist5.0.0 | a082c2dd872bb3cd8fb92359ed66605064957cd4 | Full root LICENSE.md verified, including its scope preamble; do not strip that wording. |

yallist's preamble refers packages under src to their own licences. This fixed
tree contains src/index.ts, no nested package manifest and no other named licence.
The root package declares BlueOak. That observation must not be generalized into
permission to overwrite licences of hypothetical nested packages. Four source
texts are prepared; source-to-published implementation equality and delivery of
the final notices remain separate. The five selected external directories still
have their original bytes and do not gain notice files merely because a draft
exists elsewhere.

Primary texts: [isexe](https://raw.githubusercontent.com/isaacs/isexe/2e7df7dabc4f68e88cf6b32b9029225ba52c6b0c/LICENSE.md),
[minipass](https://raw.githubusercontent.com/isaacs/minipass/ab4b3b05d0d557ac6bb178f38501b11d0c96454e/LICENSE.md),
[tar](https://raw.githubusercontent.com/isaacs/node-tar/2a22bfc5d3a432a606d9da0e2d87ba634aa3b1cb/LICENSE.md),
[yallist](https://raw.githubusercontent.com/isaacs/yallist/a082c2dd872bb3cd8fb92359ed66605064957cd4/LICENSE.md).

### OpenPGP nested Noble build inputs

OpenPGP6.3.1's fixed package-lock identifies @noble/ciphers1.3.0,
@noble/curves1.9.7 and @noble/hashes1.8.0. Their registry integrity strings match
that lock, and each fixed source package manifest and LICENSE matches the
corresponding repository tree. These three identities are absent from the
original473 inventory. They are added as a nested-source overlay, not used to
rewrite that historical graph or claim a complete476-component runtime SBOM.

| Nested package | Registry gitHead / source commit | Status |
|---|---|---|
| @noble/ciphers1.3.0 | 953f9aabb4dd91f02bb64a8df424dd1912ddd62a | PINNED_BUILD_INPUT_AND_SOURCE_TEXT |
| @noble/curves1.9.7 | a0ac59846ee76c52f7c18886f4963e1211345d48 | PINNED_BUILD_INPUT_AND_SOURCE_TEXT |
| @noble/hashes1.8.0 | 32f700f38ec49d7e6b2ab687904d6b2d7d60d80a | PINNED_BUILD_INPUT_AND_SOURCE_TEXT |

The complete ciphers licence also attributes Thomas Pornin, beyond Paul Miller
named by the selected parent's abbreviated notice. Preserve the full source text.
Matching the parent lock does not demonstrate exact correspondence to every
embedded byte region: a reproducible parent build/source mapping remains pending.
No cryptographic quality, vulnerability status or runtime approval is inferred.

Sources: [fixed OpenPGP lock](https://raw.githubusercontent.com/openpgpjs/openpgpjs/2ac0048404b74a3595d503125b53f3b3d0486bec/package-lock.json),
[ciphers licence](https://raw.githubusercontent.com/paulmillr/noble-ciphers/953f9aabb4dd91f02bb64a8df424dd1912ddd62a/LICENSE),
[curves licence](https://raw.githubusercontent.com/paulmillr/noble-curves/a0ac59846ee76c52f7c18886f4963e1211345d48/LICENSE),
[hashes licence](https://raw.githubusercontent.com/paulmillr/noble-hashes/32f700f38ec49d7e6b2ab687904d6b2d7d60d80a/LICENSE).

### node-gyp nested source parity

All70 files under selected node-gyp12.3.0/gyp match Git blobs in the fixed
nodejs/node-gyp commit154a27ef069deafc8f680f3f77452e76343546f8. This includes gyp's
LICENSE and packaging's LICENSE, LICENSE.APACHE and LICENSE.BSD. Those exact
texts were already preserved V338 and now have direct fixed-tree provenance.

Both pyproject.toml and the release manifest declare gyp-next0.22.1. The bundled
packaging/__init__.py declares23.3.dev0 while the external dependency constraint
is packaging>=24.0. Record both: a dependency constraint is not the identity of
vendored bytes. No independently established PyPI release or security conclusion
is asserted for that development-labelled snapshot. Its evidence is the exact
node-gyp tree and file hashes, without executing Python imports or build tooling.

## semver-utils remains explicitly unresolved

Current discovery cannot open the exact/full npm metadata page, and the npm code
page returned403; prior maintainer relocation remains recorded. No alternate
channel bypassed those results. The existing acquisition roster does not cover
semver-utils and requires a real licence sidecar. A generic Apache text must not
be substituted for the unverified APACHEv2 declaration to mint an admitted
source. Status remains RESEARCH_INCOMPLETE, not proof that no licence exists.
Reopen on verified applicable primary text/source binding or a separately
qualified acquisition design explicitly limited to evidence without admission.
Other source/notice work continued without waiting for that unresolved item.

## Boundaries and continuation

No implementation pack, source roster, runtime version, original source receipt,
consumer lock, registry licence label or selected payload was modified. The local
research probes and document reads are not new acquisitions of executable SDKs.
Full Preflight113154PASS retains its original date; this documentary update needs
its own structural/resume/plan checks. Expected inventory is162packs/1461files,
53profiles and777Markdown, adding this one evidence report.

Next: remaining nested works in Yarn/undici and source/relinking obligations;
chownr applicable text and semver-utils exact evidence; then complete delivery
and qualification of notices for the selected use. Draft39 is not a release.
Original full-native FAIL532, runtime, redistribution, readiness42observations,
48contract tests/7pending and10macrofronts remain open. No provider effects,
global installation, signing, publication or ARCA activity.

## Evidence

Stage: C:/Users/NL/AppData/Local/Temp/elite-v339-9266379aff0044498a1d7004267002c8.

| Stage file | SHA256 |
|---|---|
| research-receipts.json | `da4978d52f01db8bafc3b1987682d4c624a0835c558b3525428ada403f887730` |
| source-identities.json | `f6ca56ae1184b6cbf3cd7d0b84a36150c1fb6d555b1be2ef07254f17683b6b5d` |
| source-license-cases.json | `dfb121e4d6e0b5ddcd1fa39847db4242c710d4ff857d8a6072c9f63f9df7d413` |
| openpgp-lock.json | `5a83da4c7d661b8893a211551e5f371bfe98fe01f814c5f82bc6ed3b9ba9bf59` |
| gyp-vendored-parity.json | `8925f8681e4878b53edc3265ebbdc82ed391c946da41c60bb650e6674200e1ce` |
| qualification-verification.json | `086c954c24f082aad45d4140b714343eaff477924fae8cdd89307dc6e0f287c5` |
| qualify_notice_sources.py | `2ec1fa629943bda8c5035981a392da2b61c72d72a7c4e45d428a657ef17a5b94` |
| semver-research-status.json | `c04c8c13095887c167197a8e16e3faa3be79a47ffaca73081e91c8da382d9638` |
| draft-notices/MANIFEST.json | `115ad73f3eef02e38acdc6509ce75d321ace4f85d320f5c8b647e022fb37a69b` |
| draft-notices/README.md | `aa3d442bd056b1dda6c5d3d63f74b0146cb10923bb5a0cb9d044ff3fbdc75884` |
| draft-independent-verification.json | `f4fcbd86dfc56be947e4f6f7224e8af1e2efa459eb62e64153ba52c2c8fedeb1` |

## Closure119

VERIFY_LIBRARY118 PASS:162packs/1461materializable files/777Markdown. All53profile counts match V338 and all163implementation-directoryMarkdown files (162packs plus README) are byte-identical to the baseline. Draft39 and442selected payload hashes are verified. No canonical code changed and no new full154step Preflight is claimed. Final closure checks ledger uniqueness and execution/plan consistency independently. No failures were introduced by this V339 source qualification.

Next-source routing is explicit: undici6.28.0 is external with zero bundle markers; undici7.29.0 supplies112bundle markers and both matched embedded attributions. Yarn PnP4.1.7 has one bundled marker with embedded MIT works under its parent BSD declaration. Do not mix evidence across these representations/versions. This routing reuses the V338 observations and does not admit new source code.

| Closure evidence | SHA256 |
|---|---|
| verify-library118.log | `b8acb0067320d45f10b9c53e5acbbb21759a7cceeb40c373232c8bf6d83a89d4` |
| structural-closure.json | `09dc1b85d74ddec9b61fc32f8fc075992c5d53b583aabfe319d8fd69451dcdf2` |
| next-nested-source-map.json | `81e662001028466618f3c5ac7a2fc15ee221d703e09a0b706a4f18bc3fbb05e0` |
| draft-independent-verification.json | `f4fcbd86dfc56be947e4f6f7224e8af1e2efa459eb62e64153ba52c2c8fedeb1` |
