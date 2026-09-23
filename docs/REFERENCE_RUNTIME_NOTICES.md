# Exact local reference runtime notices

The portable builder checks the actual traced runtime against
`licenses/reference/runtime-notice-lock.json` before producing its catalogue.
All129observed package manifests are bound by SHA256, including framework bundles
without published version fields and unnamed export metadata. Versions are not
invented.71installed package identities are listed separately; an installed build
input is not automatically redistributed. The exact reference contains no native
addon, DLL or WASM payload under its web runtime.

The collector retains direct and nested original license/notice/LEGAL files and
the fixed source supplements listed in SOURCE_NOTICES.json. It validates all
inputs before writing the catalogue. Missing or changed runtime manifests,
unknown dependencies, missing/changed/empty notices, duplicate JSON fields,
unsafe paths, reparses and an excluded build tool in the runtime reject. Existing
catalogues cannot be overwritten. Notice collection neither installs dependencies
nor executes their scripts, and does not replace the separate SCA/source gates.

Next16.3.4 official source commit299180d3315c7ebd7b199d2b1a265b5986c5fc7d binds
the package manifests and bundling recipe. Full original README files are retained
where the publisher places its license there. The embedded dotenv16.3.1 and
dotenv-expand10.0.0 carry their original BSD-2-Clause texts. The three edge-runtime
packages carry their original npm LICENSE.md files alongside embedded notices.
String-hash retains its publisher's original CC0 waiver. Local glue is AUTHORED;
original notice files are VERBATIM, not an attribution of library code.

Three notice files are explicitly ADAPTED, not claimed as original upstream files:

- Unistore3.4.1 offers MIT with copyright Jason Miller in its original README.
  Its separate notice reproduces that stated owner, without inventing a year,
  and the standard MIT permission paragraphs. The README grant accompanies it.
- Client-only0.0.1 declares MIT in its original package.json without a copyright
  notice. The notice reproduces standard MIT terms without assigning an owner.
  Its three original files match the commit-fixed Next distributed source.
- The Babel bundle notice contains complete original copyright/license comments
  extracted from the exact bundle in their original order. The original Babel
  MIT notice and fixed Next build/import metadata accompany this correspondence.

The standard MIT permission paragraphs above match the fixed original Next MIT
text; its Vercel copyright heading is not assigned to Unistore or client-only.
These are declared notice renderings, not fabricated author/source receipts.

The runtime excludes @next/swc-win32-x64-msvc16.3.4,
@rolldown/binding-win32-x64-msvc1.2.5, server-only0.0.1 and stackback0.0.2.
Their build-input metadata remains visible. This exclusion does not certify a
native binary or authorize redistribution of the restricted package-manager
bundle. The lock must be regenerated and reviewed with any package/trace change.

The tests cover actual historical reference runtime files read-only and12adverse
publication cases. The final post-ARCA artifact must run this same collector as
part of both independent builds. No production, provider or jurisdiction claim
is derived from this local packaging proof.
