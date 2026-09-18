# Versioned Yarn and Undici notice evidence V340

2026-09-08. Library maintenance T2803/T2808. Resume119 verified. Continues the
V339 per-version routing without changing source locks, runtime or selected code.

## Undici external files and bundled version

All103 files in the selected external undici6.28.0 directory match Git blobs in
nodejs/undici commit01a912e49a50c48009ed2639d2a457a6ec26752a. Root LICENSE,
lib/web/fetch/LICENSE and package.json are included in that exact comparison.
This establishes observed source/file/notice provenance, not runtime admission.

The112 bundle markers refer to undici7.29.0, separately fixed at
9e38fc121d2eb26086d41c7d9379b47a6fada1c5. Its root/fetch licences, package manifest,
body.js and frame.js were checked against fixed Git blobs. The source root MIT
notice names Matteo Collina and Undici contributors; the nested fetch licence
names Ethan Arrowood. body.js also preserves formdata-polyfill attribution to
Jimmy Wärting and frame.js preserves ws attribution to Einar Otto Stangvik.
The two attribution comments are byte-identical across the inspected6/7 source
versions. Preserve each scope; the root licence alone does not identify all
contributors. Their exact upstream component revisions remain unqualified.

Undici7 transformed byte regions were not reconstructed. Do not apply the103-file
Undici6 result to them, or call source-marker presence a reachability proof.

Sources: [Undici6 fixed tree](https://github.com/nodejs/undici/tree/01a912e49a50c48009ed2639d2a457a6ec26752a),
[Undici7 root licence](https://raw.githubusercontent.com/nodejs/undici/9e38fc121d2eb26086d41c7d9379b47a6fada1c5/LICENSE),
[Undici7 fetch licence](https://raw.githubusercontent.com/nodejs/undici/9e38fc121d2eb26086d41c7d9379b47a6fada1c5/lib/web/fetch/LICENSE).

## Yarn registry and release-source distinction

Current official npm metadata for @yarnpkg/pnp4.1.7 still declares gitHead
5761b03feb2146da8ce8cefeba6482f3cb6edb79. GitHub tree, package directory and direct
git/commits probes all returned404 (FAIL575). The direct commit probe follows the
existing ledger lesson that commit SHA and tree SHA must not be conflated.
These failures do not establish that the source never existed.

The official annotated tag @yarnpkg/pnp/4.1.7 resolves through tag object
bbe8151ee194b5e50d8350faad66e639d2ce3b81 to commit
4fe4d4bf45a13dca90181c5b7fee61376aa21794. Its package manifest is4.1.7/BSD-2-Clause.
The tag API reports signature verification false; no local signature verification
or signed-release claim is made. The tag and registry commits remain distinct.
Tagged source text is useful evidence, not proof of the published package bytes.

The fixed tagged root BSD licence and three complete MIT comment blocks are now
preserved independently: Node.js contributors in node/resolve.js, and Blake Embrey
plus Joyent/other Node contributors in loader/node-options.js. The source README
describes files adapted from Node18.9.0 with formatting and local modifications;
this is an origin statement, not byte parity with a Node release. Blake and Node
contributor strings occur in the selected bundle; the Joyent string does not.
That absence alone does not prove relevant code survives bundling or a breach.

Status is TAGGED_SOURCE_ONLY_REGISTRY_BINDING_UNRESOLVED. Reopen the published
binding on verified registry/source provenance or governed archive/build
comparison. Never replace the npm gitHead with the convenient release commit.

Sources: [official tag](https://api.github.com/repos/yarnpkg/berry/git/ref/tags/@yarnpkg/pnp/4.1.7),
[tagged root licence](https://raw.githubusercontent.com/yarnpkg/berry/4fe4d4bf45a13dca90181c5b7fee61376aa21794/LICENSE.md),
[Node origin README](https://raw.githubusercontent.com/yarnpkg/berry/4fe4d4bf45a13dca90181c5b7fee61376aa21794/packages/yarnpkg-pnp/sources/node/README.md),
[MIT comments](https://raw.githubusercontent.com/yarnpkg/berry/4fe4d4bf45a13dca90181c5b7fee61376aa21794/packages/yarnpkg-pnp/sources/loader/node-options.js).

## Preserved draft and verification

Draft47 preserves all39 V339 text copies and adds eight version-specific texts:
Undici7 root/fetch licences, two formdata/ws attributions, Yarn root BSD and three
MIT comments. Some text is duplicated across version-bound records.47 is a count
of evidence copies, not unique cleared dependencies or a completion percentage.
Every copy is hashed and the new manifest references its parent. Earlier drafts
and all442selected payload files remain byte-identical.

The probes inspect bounded documents/metadata and existing code without installing
or executing packages. No new source acquisition roster or licence sidecar was
invented. Canonical update is documentary:162packs/1461materializable files and53
profiles remain, with778Markdown expected after adding this report. Structural
verification follows checkpoint120; full Preflight113154steps keep their date.

Remaining work includes Yarn published binding, Undici7 transformed source mapping,
formdata/ws upstream identities, semver-utils/chownr text applicability, remaining
notices and LGPL/source/relinking/delivery conditions. Original473graph is retained
as historical scope; this is not a complete recursive SBOM. Original full-native
FAIL532, readiness42observations,48contract tests/7pending and10macrofronts stay
open. No runtime/redistribution admission, provider effects, global installation,
signing, publication or ARCA activity.

## Evidence

Stage: <LOCALAPPDATA>/Temp/elite-v340-9cc3d8c83aa4416cbdcd5f3f3b158fe7.

| Stage file | SHA256 |
|---|---|
| identities.json | `e4adc0fb86c51bc1b7b5074595d1ae3fd03402c85054b2c97cedb927db852038` |
| research-receipts.json | `c80c33bd789515ad60ebfc3f334d5f3b5127ce2a4d1344ba4b19b3ab9b149dab` |
| undici-external-parity.json | `4d096586f079d525ce6651d9ba266d91af213358147e22b7813a757442bf5b81` |
| yarn-current-registry.json | `a18398a7779c9c86b4cb26cebcb773056882740cf0d32287bbc405cf4bc4d6ee` |
| yarn-release-ref.json | `c2f5e8d676bec1c77cc757b846c1698e9674dd941fb510b4c0ddf9159750ae4e` |
| yarn-release-tag.json | `6b1fb854eb43817453065b7ff29f6d048bdca7832570cf936deac7ec409d71c4` |
| yarn-release-package.json | `ceee12be54146e2748ebf4ac3b5d3029d7a75360d9bbb3946c172f250d1983c7` |
| versioned-source-cases.json | `3129fe702d90b1889d84908121622127887767557fe5f082e15c126eea846ca0` |
| qualification-verification.json | `a252a856578083a2c3281c3c55a9dcaf97fad85cbe17da559433ebb70eae8626` |
| qualify_versioned_notices.py | `11ede755b58a3bff9787a562499ec8253929e8c6ed2fd06b7da94c32c1252b66` |
| draft-notices/MANIFEST.json | `27d666e86d769c368439b1cc20e0fa72a42f9366973216b42ced5434f003fa9d` |
| draft-notices/README.md | `641727a3379cd394b93a77c6214f1d1577e40fc3aaa8459063851fb401601245` |

## Closure121

VERIFY_LIBRARY120 PASS: 162 packs / 1461 materializable files / 778 Markdown. All 53 profile counts match V339, and all 163 implementation-directory Markdown files (162 pack definitions plus README) retain their baseline hashes. Independent verification confirms all 47 draft text copies and the exact 442-file selected payload. The ledger contains 2500 unique canonical failure IDs. This documentary closure does not claim a new full 154-step Preflight.

FAIL575 remains unresolved for Yarn published-source binding. Tagged source and preserved notices are conditioned evidence; Undici6 external parity is not Undici7 transformed parity. Next work remains exact component identities, published/bundled source mapping and licence delivery/relinking conditions; original readiness and production gates remain open.

| Closure evidence | SHA256 |
|---|---|
| verify-library120.log | `0d050d3074467077322c7aca5e2dedb319c6845c0acb85cd639f3ee1cf29d4a1` |
| structural-closure.json | `e6bd5bd1c204bc42e406ea7f5c528acf10c6978fa0225294c6ffd921171949d3` |
| draft-independent-verification.json | `42f619b14a5b01181d80c2856a0195343a25965a33cf46d093991a747225b068` |
| final-uniqueness.json | `eebeadc28614eb026b2d05e3157a73acffa62387c1ad86d6a377951e3b655236` |
