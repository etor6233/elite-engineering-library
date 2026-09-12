# pnpm Artifact Selection Gate

## 1. Metadata

```yaml
pack_id: "PNPM-ARTIFACT-SELECTION-GATE"
pack_version: "0.8.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa y verifica una proyección candidata de pnpm11.25.0 desde artefacto/receipt fijados:442archivos intactos,13excluidos,22notices preservados. Bytes, publicación local exclusiva y recetas offline de tres consumers fijados; nunca instalación, ejecución ni admisión de runtime o redistribución."
stacks: ["CPython 3.12+", "Windows local filesystem"]
compatible_with: ["OFFICIAL-UPSTREAM-ACQUISITION-CORE 0.4.78", "MARKDOWN-COMPOSITOR-CORE 0.2.x"]
incompatible_with: ["artefactos móviles", "política sustituida", "ejecución automática", "directorio hostil", "promoción implícita de pnpm"]
license_expression: "LicenseRef-Workspace-Owner AND MIT AND LicenseRef-Pnpm-Bundled-Licenses AND MPL-2.0"
upstream_sources: ["https://github.com/sholladay/next-path/tree/ce2c1386836339bc473b76c9888947899ead8d56", "https://registry.npmjs.org/semver-utils/-/semver-utils-1.1.4.tgz", "https://github.com/gtanner/qrcode-terminal/tree/90f66cf5c6b10bcb4358df96a9580f9eb383307b/vendor/QRCode", "https://opensource.org/license/mit", "https://github.com/pnpm/pnpm/releases/tag/v11.25.0", "https://docs.python.org/3/library/tarfile.html"]
verified_at: "2026-09-11"
```

## 2. Applicability

Tooling de mantenimiento T2803/T2808, separado de adquisición opaca y producto.
Composición AUTHORED conforme a COMPOSITION_PROTOCOL y DEPENDENCY_UPDATE_CONTRACT;
los bytes proyectados son selección ADAPTED y conservan su licencia propia.
Usar sólo después de validar adquisición y sus condiciones; no habilita pnpm.

## 3. Architecture contract

Policy hash-bound, archive SHA256 y455miembros exactos. Recibo de adquisición
referenciado por SHA256 explícito y campos exactos; sus grants/firmas se verifican
en el workflow anterior. Todas las entradas se validan, incluso las excluidas.
Salida staging hermana con442archivos y receipt determinista; rename Windows
rechaza destinos ocupados. Verificación independiente rechaza drift, reparse,
faltantes y extras. No red, extractall, install, ejecución ni escritura de fuentes.
Directorio confiable obligatorio; sin claim de atomicidad tras power loss.

## 4. Exact file manifest

```text
CREATE pnpm_artifact_selection/README.md
CREATE pnpm_artifact_selection/selection-policy.json
CREATE pnpm_artifact_selection/select_artifact.py
CREATE pnpm_artifact_selection/test_selection.py
CREATE pnpm_artifact_selection/plan_install.py
CREATE pnpm_artifact_selection/retained-notice-catalog.json
CREATE pnpm_artifact_selection/next-path-source-evidence.json
CREATE pnpm_artifact_selection/test_routing.py
CREATE pnpm_artifact_selection/contain_zip.py
```

## 5. Materialization blocks

### FILE: `pnpm_artifact_selection/README.md`
```yaml
block_id: "PNPM-ARTIFACT-SELECTION-GATE:README.md:v395"
operation: CREATE
provenance: AUTHORED
source: "local, bounded projection and consumer planning governed by canonical acquisition and execution contracts"
license: "LicenseRef-Workspace-Owner"
sha256: "569a728e21e78eb168b32193310d6724f3ad180c10fc130f4c8dae2baabf98a6"
variables: []
secrets_allowed: false
```
````markdown
# pnpm artifact selection — narrow tooling gate

This AUTHORED utility creates an ADAPTED candidate from the exact pnpm11.25.0
archive already acquired with the official source profile. It has no network,
install, execution, signature, runtime admission or redistribution operation.
The original455-file artifact remains immutable. The projection retains442files,
all22supplied notices, and excludes the entire physical @reflink family plus
two Fastlist files (13files, including6native binaries). It is not the official
unmodified pnpm release. Policy/source changes require a new version and gates.

## Conditions and invocation

Windows, CPython3.12+ with the project's admitted exact runtime, trusted local
directories and a previously validated acquisition profile/receipt are required.
Keep the raw artifact and acquisition receipts. Obtain the receipt SHA256 from
that validated workflow, not from an untrusted sender's claim. The selector
matches exact artifact/receipt fields and preserves their hash references; it
does not reverify signatures, grants or approvals of the acquisition workflow.
The acquisition receipt remains `extracted=false`: it describes that earlier
operation. This separate receipt describes the local projection.

```powershell
python select_artifact.py materialize --artifact <fixed-pnpm-11.25.0.tgz> --acquisition-receipt <receipt.json> --acquisition-receipt-sha256 <trusted-SHA256> --target <absent-directory>
python select_artifact.py verify --target <selection-directory>
python -m unittest -v test_selection.py
```

Output contains `payload/` and `selection-receipt.json`, whose runtime and
redistribution admission fields are always false. Success verifies exact bytes;
it does not authorize pointing a consumer at `payload/bin/pnpm.mjs`. The CLI
accepts no substitute policy, extraction mode, version, install or execute flag.

## Failure and recovery

Wrong hashes, duplicated JSON, unsafe TAR entries, links/reparse points, altered
or missing output, extra files/directories and occupied targets are rejected.
Payload SHA256 is checked before parsing; every member, even an excluded member,
is verified. No `extractall` or external program is used. Staging is a newly
created sibling, verified before Windows no-overwrite directory rename. A late
competitor survives. Ordinary exceptions clean only this invocation's staging;
source bytes, receipts and occupied targets are never removed or replaced.

After interruption, preserve any remaining `.pnpm-selection-*` candidate as
evidence, verify the final target if present, and retry into another absent
target. No fsync/power-loss atomicity or hostile-parent-directory guarantee is
claimed. Directory ACLs and concurrent hostile modification remain outside
this utility's guarantee. Never delete an unrelated target to make a retry pass.

## Before actual use

V333 proved Windows copy/hardlink installs and bounded compatibility; it did not
admit all pnpm CLI configurations. Scope the consumer routing to those methods,
test unsupported-mode rejection, all affected browser/Lighthouse/operational
paths, complete scoped license review and runtime/dependency monitoring before
adoption. Application native dependencies are distinct and are not removed by
this projection. Retained LGPL/MPL/Artistic/CC-BY and other notices keep their
obligations; the parent MIT label does not replace them. No redistribution right
or security acceptance is inferred solely from a clean scan or this PASS.

Tests use explicitly synthetic tiny archives and receipts. A separate integration
must materialize the actual pinned artifact twice and compare every output with
the locked selection, including receipts and retained notices. The public CLI
always reads the fixed hash-bound455-entry policy; synthetic policies are only
internal unit-test data, never a command-line bypass.

## Pinned consumer routing (0.2.0)

plan_install.py prepares and verifies a single-use offline-install recipe bundle;
it never runs it. Only enterprise-web, playwright and lighthouse manifest/lock
hashes recorded in the code are supported, with Node24.20.0 Windows x64 exact.
The project must be a fresh isolated copy without node_modules. Store and metadata
cache are explicit existing, separate directories. Content-store completeness
alone does not prove offline policy metadata or cached verification is available.

prepare requires --consumer, --project, --projection, --node, --store, --cache,
--target, --acquisition-receipt and --acquisition-sha256. The output target must
be absent; inputs cannot overlap it or each other and cannot use reparse points.
Retain the printed SHA256 externally. Recheck immediately before any authorized
use: python plan_install.py verify --target <bundle> --sha256 <retained-SHA256>.

The recipe fixes argv and cwd, shell=false, a REPLACE_NOT_MERGE environment,
empty user/global npmrc files and isolated home/config/temp directories. It omits
host credentials, NODE_OPTIONS and NPM_CONFIG options. It requires copy imports,
offline, frozen lockfile, ignore-scripts, ignore-pnpmfile and store integrity.
The exact enterprise workspace is retained (including its supply-chain policy);
standalone browser/quality consumers use ignore-workspace. Unknown consumer
inputs, local/ancestor npmrc/hooks, workspace variants and existing installs
are rejected. No arbitrary arguments, scripts, cloning or global-install lane.

A recipe is not runtime/license admission or an OS sandbox. Metadata caches,
operating system and input directories must remain trusted and unchanged until
use. A separate authorized consumer must replace the environment completely,
not merge it, and satisfy remaining runtime/license gates. Never use trust-lockfile
or weaken policy merely to make an offline install pass. A used home/consumer or
any tampered receipt/config is no longer a pristine verified plan. Concurrent or
late target writers never get overwritten. A failure after publication preserves
the diagnostic candidate; delete no existing project to recover. No power-loss
or hostile-directory atomicity claim. Recreate into a new absent destination.

V335:26routing regressions plus the17projection regressions, independent six-file
rebuild and actual isolated consumer probes. The17/4-file V334 evidence remains
historical; this version adds2AUTHORED files and updates this guide.

## Current consumer and supplemental notice revision (0.3.0)

The enterprise-web inputs now match TS-GO-API-WEB-BRIDGE0.5.16, including the
reviewed 122-version lock and ignored optional Sharp workspace policy. Playwright
matches0.1.37; its lock and Lighthouse inputs are unchanged. Earlier consumer
identities remain in V335 evidence, not as accepted fallback inputs. Revalidate
this exact routing profile after any consumer manifest, lock or workspace update.

Every new v2 plan includes BLUEOAK-NOTICE.md, bound by SHA-256 in install-plan.json.
It supplies the official Blue Oak Model License1.0.0 link for five exact package
root declarations: chownr3.0.0, isexe4.0.0, minipass7.1.3, tar7.5.22, yallist5.0.0.
The generator verifies each manifest hash and name/version/license before use;
fixed author-declaration links are retained. Blue Oak's Notices clause permits
delivery of the license text OR its official link:
https://blueoakcouncil.org/license/1.0.0

This generated supplement is outside the immutable442-file projection. Keep it
with any subsequent copy of those works. The planner does not copy or distribute
the payload and cannot prove delivery to downstream recipients. It does not
invent an original chownr LICENSE file, drop the yallist nested-package caveat,
or satisfy the remaining pnpm licenses and source/publishing obligations.
runtime_admitted, redistribution_admitted and executed remain false.

verify rejects missing or changed notices, even when the attacker rehashes the
plan and replacement notice together. The old v1 plan is historical and cannot
be verified by this revision; create a new plan into an absent destination.
No inputs or previously prepared plans are migrated or overwritten.
Run both suites: python -m unittest -v test_selection.py test_routing.py
Current suites contain17projection and34routing tests (51total).

## QRCode complete attribution supplement (0.4.0)

New v3 recipes additionally contain QRCODE-NOTICE.md. The notice preserves the
QRCode vendor index header (copyright2009 Kazuhiko Arase, MIT declaration,
trademark and local modification notice) from qrcode-terminal0.12.0 at commit
90f66cf5c6b10bcb4358df96a9580f9eb383307b. It supplies the complete MIT permission
text from https://opensource.org/license/mit with that copyright identity.
This is an explicitly assembled supplement, not an upstream archive LICENSE file.

The pinned pnpm bundle and its10exact vendor module regions are checked before
generation. Fixed-source blobs and selected bundle region hashes are recorded
separately; no source/build equivalence is inferred from a module path.
The containing package's Apache license and all original notices remain intact.
Keep this supplement and BLUEOAK-NOTICE.md with copies of the relevant works.
The planner delivers both files locally; it does not distribute the payload.

Removal of the author, permission text, notice, or jointly rehashed notice/receipt
is rejected. Publication faults preserve occupied destinations and clean only
the newly owned stage. Prior v1/v2 plans remain historical; recreate into an
absent destination. Runtime, redistribution and execution fields remain false.
Current suites:17projection +42routing tests =59. No new runtime dependency.

V396 / 0.5.0: v4 recipes also retain SEMVER-UTILS-NOTICE.md, including the original dual-license file and its complete MIT option. Exact artifact SHA512 and original member hashes are recorded separately from the selected bundle region. APACHEv2 metadata discrepancy and expired registry key remain explicit. This closes license-text discovery and local notice delivery, not whole pnpm admission. Verify refuses missing/altered/rehashed notices. Existing BlueOak/QRCode supplements and immutable payload are preserved.

V397 / 0.6.0: recipe v5 delivers PNPM-RETAINED-NOTICES.md with all47 exact retained text copies (22 original payload notices and25 research texts). The catalogue is pinned, bounded and per-entry checked; original notices must also equal selected payload bytes. This closes local evidence collection/delivery only. It is not complete license coverage, source offer, relinking fulfillment or runtime/redistribution admission. Existing notices and442payloadfiles remain unchanged. No external package code or scripts execute.

V398 / 0.7.0: v6 recipe delivers NEXT-PATH-MPL-SOURCE.md with complete fixed index.js, package.json and MPL-2.0 license, URLs/Git blobs/digests and explicit bundler transformation record. Source files are evidence data, never installed/executed. Four module statements compare under named adapters with10negative probes; __require host resolution, npm artifact identity and whole-pnpm build remain unproven. Prior47texts and three supplements stay byte-identical. No blanket runtime/redistribution admission.

## V401: published pnpm security block and isolated ZIP removal

The unmodified pnpm11.25.0 payload is affected by GHSA-vwc7-r8mq-g2x9.
Its historical scan does not authorize new execution. pnpm11.26.0 is also
affected. Version12.4.1 is a native candidate, not an admitted replacement.

contain_zip.py creates a separate ADAPTED candidate from all455 exact published
files. It retains441 selected files unchanged, changes only dist/pnpm.mjs,
removes all16 adm-zip modules and rejects binary ZIP acquisition before download
or extraction-directory creation. TAR support is unchanged. It retains original
notices and writes ADAPTATION.txt, adaptation.diff and a deterministic receipt.
Original input and existing targets are never overwritten. Trusted Windows
directories are required; concurrent hostile directory mutation is out of scope.

python contain_zip.py create --source <verified-published-pnpm-directory> --target <absent-candidate-directory>
python contain_zip.py verify --source <same-published-directory> --target <candidate-directory>

This utility performs no execution, installation, network call or admission.
The old plan_install.py recipes still bind the unmodified projection and must
not be redirected to the changed bundle. The candidate's bounded qualification
uses a separate environment and explicit digest. Its identity is the receipt
and bundle SHA, not the unchanged upstream version banner. Original licenses,
source/relinking obligations and general runtime/redistribution admission remain.
Binary ZIP runtime installation is intentionally unavailable. Reopen on an
admissible official fixed release, a new consumer or any change to the bundle.
Evidence: reconstruction_evidence/PNPM_ZIP_CONTAINMENT_V401.md.

````

### FILE: `pnpm_artifact_selection/selection-policy.json`
```yaml
block_id: "PNPM-ARTIFACT-SELECTION-GATE:selection-policy.json:v1"
operation: CREATE
provenance: AUTHORED
source: "local, bounded projection and consumer planning governed by canonical acquisition and execution contracts"
license: "LicenseRef-Workspace-Owner"
sha256: "93c27e601ccc04bbadb755483186e96d3aca0bad6b37063c668e934e55bc43f6"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-pnpm-selection-policy/v1",
  "archive_sha256": "33dd0748f27e7916c4f1c8b6943461983e3453b06bbda6312a6280130b4881e5",
  "archive_bytes": 5116259,
  "receipt_bindings": {
    "schema": "elite-official-source-receipt/v1",
    "id": "maintenance-pnpm-11.25.0",
    "repository": "pnpm/pnpm",
    "release": "v11.25.0",
    "commit": "6d90c71efdffbc909b499490b64c66badc720327",
    "transport_kind": "opaque-release-artifact/v1",
    "artifact_url": "https://registry.npmjs.org/pnpm/-/pnpm-11.25.0.tgz",
    "artifact_filename": "pnpm-11.25.0.tgz",
    "artifact_bytes": 5116259,
    "artifact_digest": {
      "algorithm": "sha256",
      "value": "33dd0748f27e7916c4f1c8b6943461983e3453b06bbda6312a6280130b4881e5"
    },
    "artifact_sha256": "33dd0748f27e7916c4f1c8b6943461983e3453b06bbda6312a6280130b4881e5",
    "license_expression": "MIT AND LicenseRef-Pnpm-Bundled-Licenses",
    "license_sha256": "e0a867ff513ea7be2a0ddc339ac6a031e459a38668e077b8f0e649544062f9f2",
    "extracted": false,
    "installed": false,
    "executed": false,
    "production_admitted": false,
    "platform": "windows"
  },
  "files": [
    {
      "path": "bin/pnpm.cjs",
      "sha256": "67b035e322203961795e8e34ca63a08c37a4386eda94107fb3d28f3246d882ad",
      "bytes": 105,
      "selection": "unchanged"
    },
    {
      "path": "bin/pnpm.mjs",
      "sha256": "ff3224d46b47fbb24a7e9fe15fededef7e00892d07d4e376b6762d4899906bfd",
      "bytes": 1464,
      "selection": "unchanged"
    },
    {
      "path": "bin/pnpx.cjs",
      "sha256": "3cbdc69ac2a2bc3d88e38840b73c80d40c553ab9d8c8ab832639b69f1a1baf5d",
      "bytes": 105,
      "selection": "unchanged"
    },
    {
      "path": "bin/pnpx.mjs",
      "sha256": "7e2a61f1636e6d85fbd894b7062b1837e175c673c3dc9ba9bc7538fe1ce66502",
      "bytes": 127,
      "selection": "unchanged"
    },
    {
      "path": "CHANGELOG.md",
      "sha256": "43632aa0d0766c58c5d6e7efe2ed89485fb4759096bec2fbbeaedb6e14b0f14b",
      "bytes": 291496,
      "selection": "unchanged"
    },
    {
      "path": "dist/node-gyp-bin/node-gyp",
      "sha256": "7e140dbab20e3a8bce969eea8c4ae0103638ee1837cf7e08e027f19e648b4a6d",
      "bytes": 169,
      "selection": "unchanged"
    },
    {
      "path": "dist/node-gyp-bin/node-gyp.cmd",
      "sha256": "3f68a71f7cc3a097061e66169bddd70f7661ccf9c0e992dad5c2dbe972e1383c",
      "bytes": 141,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/.bin/node-gyp",
      "sha256": "b16c955c37f0562c2130f85f92c79e11b59177f6de644d617699fbb74a2619ee",
      "bytes": 1538,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/.bin/node-which",
      "sha256": "7bb4d292d2a84d4a412e2771a89cdac4cc4ec6a7dfcb2996b96d05d48f6b8b4c",
      "bytes": 1502,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/.bin/nopt",
      "sha256": "ea746e06e46ebe585c848e17e4ebe9707c4a5eb2edf0f2de15e64a549ac7b2e1",
      "bytes": 1488,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/.bin/semver",
      "sha256": "e936f771c42a3305617205c75792531292ff31cd81e1293b00a89bb941afd922",
      "bytes": 1512,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/.package-map.json",
      "sha256": "a13ca782d47aabcab1c575da8805e331aed27d4cdbe6811dbd80d7b92a412e3c",
      "bytes": 3176,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/@isaacs/fs-minipass/dist/commonjs/index.js",
      "sha256": "2c934649bd1f847982fc06b02721244cb3271bf254afd71f23f9d20c14f9b59a",
      "bytes": 12791,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/@isaacs/fs-minipass/dist/commonjs/package.json",
      "sha256": "8005a3491db7d92f36ac66369861589f9c47123d3a7c71e643fc2c06168cd45a",
      "bytes": 25,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/@isaacs/fs-minipass/dist/esm/index.js",
      "sha256": "6feab1ba6210a3dd4c2bf1eb9cadfca5674c82972d0d015a162bbeb097a886da",
      "bytes": 12129,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/@isaacs/fs-minipass/dist/esm/package.json",
      "sha256": "3ca9d4afd21425087cf31893b8f9f63c81b0b8408db5e343ca76e5f8aa26ab9a",
      "bytes": 23,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/@isaacs/fs-minipass/LICENSE",
      "sha256": "4ec3d4c66cd87f5c8d8ad911b10f99bf27cb00cdfcff82621956e379186b016b",
      "bytes": 765,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/@isaacs/fs-minipass/package.json",
      "sha256": "ab8c1c2bce664e4c76b1937fa1a28f00ca5231a70cf0c1f78dd473101a9d2a88",
      "bytes": 1675,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/@reflink/reflink/binding.js",
      "sha256": "5b71b22594a8131250ff2a5d9bd003a7fc7a5dbb632e03e5770f2e8bf346ffab",
      "bytes": 7516,
      "selection": "excluded_native_family"
    },
    {
      "path": "dist/node_modules/@reflink/reflink/index.js",
      "sha256": "e07548eb450428530430ec9d2c43daccccffe1305c754ff558b21a525f515cfe",
      "bytes": 1326,
      "selection": "excluded_native_family"
    },
    {
      "path": "dist/node_modules/@reflink/reflink/package.json",
      "sha256": "e99ee5ff222a9356c129c180a6ae8cbd9aa2af62e40ab4bd4c1f1c0212f72521",
      "bytes": 1826,
      "selection": "excluded_native_family"
    },
    {
      "path": "dist/node_modules/@reflink/reflink-darwin-arm64/package.json",
      "sha256": "fbbf219b021f4d1175864cd5ea0eec9ff8aaee56f710fa2c5ff4cbae099f08e1",
      "bytes": 397,
      "selection": "excluded_native_family"
    },
    {
      "path": "dist/node_modules/@reflink/reflink-darwin-arm64/reflink.darwin-arm64.node",
      "sha256": "0cf6caf1b125009161c98e7fb39d996da26548b09e7b17d7688357f9c3b97a02",
      "bytes": 404184,
      "selection": "excluded_native_family"
    },
    {
      "path": "dist/node_modules/@reflink/reflink-darwin-x64/package.json",
      "sha256": "ef663980cb91dd3070a8ca5945a087d0ffd6e3bbc842d35697e143b18614ac05",
      "bytes": 387,
      "selection": "excluded_native_family"
    },
    {
      "path": "dist/node_modules/@reflink/reflink-darwin-x64/reflink.darwin-x64.node",
      "sha256": "e107629c8000ad509954bdc84e4bb8fec38453c210d8a4370a9fadb28aff5b97",
      "bytes": 417392,
      "selection": "excluded_native_family"
    },
    {
      "path": "dist/node_modules/@reflink/reflink-win32-arm64-msvc/package.json",
      "sha256": "b6684d5ab4d22565fc50083c86a088ad3b9f13d31f09048fabc79a2a37e570be",
      "bytes": 412,
      "selection": "excluded_native_family"
    },
    {
      "path": "dist/node_modules/@reflink/reflink-win32-arm64-msvc/reflink.win32-arm64-msvc.node",
      "sha256": "b7c3fe5620ead8872593acd949fac7bb4d860ec015e320f7e117260fa3b12d96",
      "bytes": 332800,
      "selection": "excluded_native_family"
    },
    {
      "path": "dist/node_modules/@reflink/reflink-win32-x64-msvc/package.json",
      "sha256": "17234edc0cff64b71ad78d6194731b5932998cea4e878665078e3e9e2f8ccd96",
      "bytes": 402,
      "selection": "excluded_native_family"
    },
    {
      "path": "dist/node_modules/@reflink/reflink-win32-x64-msvc/reflink.win32-x64-msvc.node",
      "sha256": "370cdee825881d1f0c3e87a85df5589ee51f1a7eec3bca3c252255b6521224dc",
      "bytes": 361472,
      "selection": "excluded_native_family"
    },
    {
      "path": "dist/node_modules/abbrev/lib/index.js",
      "sha256": "92f1fd161bdab2f33a468e4849c58d3a1dfe9bba4988e5b3686cdd4c81706caa",
      "bytes": 1379,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/abbrev/LICENSE",
      "sha256": "9e0d5c7989f7e9f07d7c4b158aceff270f235eb7464ace41c5e7b200834a43e0",
      "bytes": 2011,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/abbrev/package.json",
      "sha256": "86c0a87ee85682437c2617ea516604ac19dcc7b8e1fc5d34a2220f320347c627",
      "bytes": 1232,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/chownr/dist/commonjs/index.js",
      "sha256": "02b4ca7ab42bdb4ab00bfe7f17d4dafd215e66fd96105e949b4f918b37101bc4",
      "bytes": 3059,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/chownr/dist/commonjs/package.json",
      "sha256": "8005a3491db7d92f36ac66369861589f9c47123d3a7c71e643fc2c06168cd45a",
      "bytes": 25,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/chownr/dist/esm/index.js",
      "sha256": "09b582ed90c2acc9c457a621fdbd7e6930a099e5089958917c1ad4d823f832e9",
      "bytes": 2534,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/chownr/dist/esm/package.json",
      "sha256": "3ca9d4afd21425087cf31893b8f9f63c81b0b8408db5e343ca76e5f8aa26ab9a",
      "bytes": 23,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/chownr/package.json",
      "sha256": "4300e90fdd91ec7035047473c60f880251a9801bd786302729d4277751d3b948",
      "bytes": 1621,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/env-paths/index.js",
      "sha256": "84351667051b005f9856319267407b06affad12888355462c7c8740f22ca3999",
      "bytes": 2155,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/env-paths/license",
      "sha256": "48da2f39e100d4085767e94966b43f4fa95ff6a0698fba57ed460914e35f94a0",
      "bytes": 1109,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/env-paths/package.json",
      "sha256": "b92833e5851ec53bd4cd8093f6099d0c6e6818c4374ec1d09aa25d4f9ba91ec4",
      "bytes": 698,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/exponential-backoff/dist/backoff.js",
      "sha256": "2b5e92bc4531c746ede193e095e3564902083334df9b53a66dc91ecb4e0a05cd",
      "bytes": 5903,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/exponential-backoff/dist/delay/always/always.delay.js",
      "sha256": "82ed246fdf73d87c8219a57cf2818a0eaa2b2eb890155e3b4fa933fc81f26f3b",
      "bytes": 1028,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/exponential-backoff/dist/delay/delay.base.js",
      "sha256": "ff428341472df66189bc55552bdc1fdcc82414fe67f6ac7a6d085a90cb87e39c",
      "bytes": 1530,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/exponential-backoff/dist/delay/delay.factory.js",
      "sha256": "fd0b733e12e019be65d70f4edb78e48ce9b8ea1ed19a653a4839e11a995e88a0",
      "bytes": 613,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/exponential-backoff/dist/delay/delay.interface.js",
      "sha256": "1c59b7a3a0c9fa05dfba47568210865c8e16a1ff60183b3da5782b6785e9dc28",
      "bytes": 120,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/exponential-backoff/dist/delay/skip-first/skip-first.delay.js",
      "sha256": "0e45504b53ec2a110f642939b18dd1b963728e8a3ed209208a72860c55639c0a",
      "bytes": 4249,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/exponential-backoff/dist/jitter/full/full.jitter.js",
      "sha256": "547b50cb458c943ba82ad3602e5129289430bac64c483e3db99427aa2703c7c1",
      "bytes": 265,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/exponential-backoff/dist/jitter/jitter.factory.js",
      "sha256": "e51d94ddc719c2b952a50d7a295596cf6bd4280b53e4dcab1abd15784ce83ca4",
      "bytes": 471,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/exponential-backoff/dist/jitter/no/no.jitter.js",
      "sha256": "d94bf184e944abd02643a95a37bfed32757d8763b6fb1e52f5b593d37f8bf903",
      "bytes": 190,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/exponential-backoff/dist/options.js",
      "sha256": "8d2c5b2d2b69e7f8f18c32fe255e83413a0013cad2216ecd34d07fe566ea3b17",
      "bytes": 969,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/exponential-backoff/LICENSE",
      "sha256": "f5afc95c69e116ea6b1755f91a8763f53eadb78b4f1582a08d2d18586a03a117",
      "bytes": 11351,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/exponential-backoff/package.json",
      "sha256": "c3076cc23aa8d9f5123900305b9c174af612bfef91211c658ab914611dd86fa8",
      "bytes": 1329,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/fdir/dist/index.cjs",
      "sha256": "ec92c1fdd78905a7b6d1ea3391aeca2cf19d10e6bbd04db0de158fbfc8b6cff4",
      "bytes": 18210,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/fdir/dist/index.d.cts",
      "sha256": "84e9a7a30adfbe64e51de98c71cc5b8403f6c1af16301c6f40fc2febc4be6fd3",
      "bytes": 5129,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/fdir/dist/index.d.mts",
      "sha256": "84e9a7a30adfbe64e51de98c71cc5b8403f6c1af16301c6f40fc2febc4be6fd3",
      "bytes": 5129,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/fdir/dist/index.mjs",
      "sha256": "4a0266bf60373bc08bcac35a813587ead767864b3a80b5dca10bff999e05e32d",
      "bytes": 17225,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/fdir/LICENSE",
      "sha256": "9a39f2aadab11a3697edd668ff2d8ad885b649737b7ab4d3bf12b34e5ada0c86",
      "bytes": 1053,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/fdir/package.json",
      "sha256": "aca49395d61ef383a09f74555c66ed45817ca2b5724e710a4dc9e0471b83d8f5",
      "bytes": 2694,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/graceful-fs/clone.js",
      "sha256": "7258eca52e65d69845759503f9fdd66c252f40e5eafb76db5d481172e31ac9ed",
      "bytes": 496,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/graceful-fs/graceful-fs.js",
      "sha256": "b47953721e408e70e23753c83ccdd06baaf40e3b942aa5e0e1d56db71287cc54",
      "bytes": 12704,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/graceful-fs/legacy-streams.js",
      "sha256": "60a6a7ecf7c3e55a3ffaae13433b6cff388b7205bba6daf393c863f77a949e36",
      "bytes": 2655,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/graceful-fs/LICENSE",
      "sha256": "f65c5d9f22a317b2a10803bd1868461ce6499c2ed7217bc80c0cc772a748789c",
      "bytes": 791,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/graceful-fs/package.json",
      "sha256": "5747d4ba6b17165c6ecac30ab3a331715f41c7ad546e1f1574dab1bdcb116181",
      "bytes": 1031,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/graceful-fs/polyfills.js",
      "sha256": "66ea1687ed5edf39d67296d26edccc8da695d9a869303a78d0e580cd770aca27",
      "bytes": 10141,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/isexe/dist/commonjs/index.js",
      "sha256": "fbaffc9133ade6b2df4a2f2824f4a84ac92bfd8046f380ea110c86852f88492a",
      "bytes": 2283,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/isexe/dist/commonjs/index.min.js",
      "sha256": "abf1fa35f4710b8f74351b63e61945b59994cdbceaf78cb9e309fb8b582a3a35",
      "bytes": 3030,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/isexe/dist/commonjs/options.js",
      "sha256": "474f49f173fdd0657ded670ab46bb209a05b08f9f52efb5e4da2045ba7d097a3",
      "bytes": 112,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/isexe/dist/commonjs/package.json",
      "sha256": "8005a3491db7d92f36ac66369861589f9c47123d3a7c71e643fc2c06168cd45a",
      "bytes": 25,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/isexe/dist/commonjs/posix.js",
      "sha256": "97ad32a2c5583b0c88ee3c79f52fa0e6dea3433eb4e63b46e880dd85635bad82",
      "bytes": 2092,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/isexe/dist/commonjs/win32.js",
      "sha256": "bfde01a62a183af4936a3dec114bd24faa57535f8de0fb2ccbc0de737741af8f",
      "bytes": 1936,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/isexe/dist/esm/index.js",
      "sha256": "c123f7ac6fde508df0623580c8b0d7596cee6bf3955842cd6a0d18ede42ca6dc",
      "bytes": 516,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/isexe/dist/esm/index.min.js",
      "sha256": "c1d43c88f2b5eff554cc40b5702c83efd0f33c7ff92277675d344a047cc72984",
      "bytes": 1670,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/isexe/dist/esm/options.js",
      "sha256": "fea76c2f7b85cf0fae3fc883565127911873222c84c8ee41bbf8f3a6ac3881ea",
      "bytes": 46,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/isexe/dist/esm/package.json",
      "sha256": "3ca9d4afd21425087cf31893b8f9f63c81b0b8408db5e343ca76e5f8aa26ab9a",
      "bytes": 23,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/isexe/dist/esm/posix.js",
      "sha256": "db612c862f36ae60256e45d98f883d8db649a4928431c2d7255415f885bcb94a",
      "bytes": 1906,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/isexe/dist/esm/win32.js",
      "sha256": "072cc87f6d0b77ea25766de04623295ad45ee82d42da60af9cdab6be1f994cdc",
      "bytes": 1735,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/isexe/package.json",
      "sha256": "1889746f84d1e2a524353a2e31464995b808f24457f0dd55d6fd9da638f978ed",
      "bytes": 1977,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/minipass/dist/commonjs/index.js",
      "sha256": "085dc9747f0f4b28a196353df3702e2ff8cf0bc883a02378df062525a7cf9a38",
      "bytes": 34050,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/minipass/dist/commonjs/package.json",
      "sha256": "8005a3491db7d92f36ac66369861589f9c47123d3a7c71e643fc2c06168cd45a",
      "bytes": 25,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/minipass/dist/esm/index.js",
      "sha256": "2690b4752df99b1e7959802070037ba2898908675e75f1e32c358578a4456dd3",
      "bytes": 33328,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/minipass/dist/esm/package.json",
      "sha256": "3ca9d4afd21425087cf31893b8f9f63c81b0b8408db5e343ca76e5f8aa26ab9a",
      "bytes": 23,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/minipass/package.json",
      "sha256": "b8ab116187d5e9375d494f17437eb862cbee4b329c21f15412c6a5170ec2f55c",
      "bytes": 1904,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/minizlib/dist/commonjs/constants.js",
      "sha256": "1574ad656e72d59598f1fd44bf2d242ff57d05ebd27cdf177b504e24b3e64fdf",
      "bytes": 4299,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/minizlib/dist/commonjs/index.js",
      "sha256": "1545b857e5a597a4fb9294e5040f8b8e3f207e53901f935de2b135a4200abbbf",
      "bytes": 14870,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/minizlib/dist/commonjs/package.json",
      "sha256": "8005a3491db7d92f36ac66369861589f9c47123d3a7c71e643fc2c06168cd45a",
      "bytes": 25,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/minizlib/dist/esm/constants.js",
      "sha256": "7c219d10414779e805f027504b8e1b20105dfa727cd01b6a5ef32158725f6033",
      "bytes": 4034,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/minizlib/dist/esm/index.js",
      "sha256": "d5416c33c5948b2cd9eb196c1a7f789811767a77a616129dd9f782fa5781c64a",
      "bytes": 12247,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/minizlib/dist/esm/package.json",
      "sha256": "3ca9d4afd21425087cf31893b8f9f63c81b0b8408db5e343ca76e5f8aa26ab9a",
      "bytes": 23,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/minizlib/LICENSE",
      "sha256": "df05697a04b8997c79d423930534e59cc18828f6e22d0c8c2d62fcbea1c4f50b",
      "bytes": 1339,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/minizlib/package.json",
      "sha256": "3e3b206f94f3bd444f18fed847d8caeacb84925c628ae3d87805bc1dd6c4c470",
      "bytes": 1842,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/.release-please-manifest.json",
      "sha256": "f88c6931d20e9023fe94321dadf4234f56b49bd4d48ff4fe9f13ccf5f7a4b5c4",
      "bytes": 22,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/addon.gypi",
      "sha256": "ab6904affe2b17de17b0450bada73243e62c841295a3efd69403e269f6eddbee",
      "bytes": 5952,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/bin/node-gyp.js",
      "sha256": "733bdbf2080c31a67d94cb8f7123110f44b0a3e2be24b8e575754af550eee6bf",
      "bytes": 3620,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/eslint.config.js",
      "sha256": "44f2ae8aadd8d676a4ce2ac98a6d779aca5f1f1c507b9965b534ea0004a81283",
      "bytes": 58,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/.release-please-manifest.json",
      "sha256": "070c8bff5eef46612632151175413f1a9bd032ef55dfbb41c157446f71c8e893",
      "bytes": 22,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/data/ninja/build.ninja",
      "sha256": "d67fbbbf43d27816323c80fef7dad205536c831e521e650579074d82499b5b1b",
      "bytes": 56,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/data/win/large-pdb-shim.cc",
      "sha256": "e1c759dda39fa50264575092a136e33f28211139b332d88e933a1b953d564f90",
      "bytes": 653,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/gyp",
      "sha256": "38e8b886cf06fe7cec4d89634fd2850891706308e2bbbc0556b3d299bd6a7993",
      "bytes": 240,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/gyp.bat",
      "sha256": "77bf45e8c077df03d65e6c076920f24bee04752e29bcb21b63d3622fffe84f10",
      "bytes": 201,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/gyp_main.py",
      "sha256": "80eb0a1948230c832f17bf99f31a8729d54044d80469cfe0fecace69e7ab66a6",
      "bytes": 1250,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/LICENSE",
      "sha256": "ca90abb6ed71de0774461ef9f928de33e748b617aeb79f9e52415cf08d69230e",
      "bytes": 1537,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/pylib/gyp/__init__.py",
      "sha256": "0bc04e8561199ae81da51e7f45fdf271fd2d6025c3bda411d622d59f54ee3429",
      "bytes": 24614,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/pylib/gyp/common.py",
      "sha256": "cbbc59947402ff878a121a6d95034e8b04016486154502beb9f9b84c4b2d8db8",
      "bytes": 25313,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/pylib/gyp/common_test.py",
      "sha256": "d9831158cb9393a3d3b8d626731298e883602fffaaf2394ee6db658d01b92c70",
      "bytes": 6252,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/pylib/gyp/easy_xml.py",
      "sha256": "4604946abac2a8cb6dd24ad2f35c0644f6c0f2342967ce916ab3091fd41477ad",
      "bytes": 5487,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/pylib/gyp/easy_xml_test.py",
      "sha256": "d4ca2f5b01ab57596bcee58881a75759f00cd978c413cb7d814c174f22d91229",
      "bytes": 3941,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/pylib/gyp/flock_tool.py",
      "sha256": "36a3f8725c44fbaa555d57c39c3896d170283e164c53ff3ebe59cb43db393c1a",
      "bytes": 1886,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/pylib/gyp/generator/__init__.py",
      "sha256": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
      "bytes": 0,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/pylib/gyp/generator/analyzer.py",
      "sha256": "8acc9c46eb01b763453bf72174243e451df1e6f992c634b0bdac4d44d68a9d3e",
      "bytes": 31639,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/pylib/gyp/generator/android.py",
      "sha256": "fc3591b7740e964bccaf7b7e54d5b08652030c6b503d95b661c3edb29fb7f995",
      "bytes": 49773,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/pylib/gyp/generator/cmake.py",
      "sha256": "3a2b1c63330afa01368b38e08340eac30850110e75959da07a823c185d22a792",
      "bytes": 49299,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/pylib/gyp/generator/compile_commands_json.py",
      "sha256": "4b4be68b5e9c8d2392ad81ff7ec3d16204876cf6e553c83e003a37b307c19a5f",
      "bytes": 4851,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/pylib/gyp/generator/dump_dependency_json.py",
      "sha256": "be29b3c4cccaed8a29f78f933f0d864e95e5223010f70b7d9593e4a756ae5a55",
      "bytes": 3104,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/pylib/gyp/generator/eclipse.py",
      "sha256": "309aaac9b59683cbfaa72f4d5906c9d9da88cc50a55164c043a240adfd53129d",
      "bytes": 17519,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/pylib/gyp/generator/gypd.py",
      "sha256": "327da16f225f0b25ad0468b3b1e5cad5e5b5f7943adae943345647796835a2ca",
      "bytes": 3505,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/pylib/gyp/generator/gypsh.py",
      "sha256": "712cd6f9d744c7680250f20e8fd6b87d74408264f3dafd5e0b28e9fc084f6938",
      "bytes": 1692,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/pylib/gyp/generator/make.py",
      "sha256": "0ce6c09023c2cf4fa48f2c5d17d0f6c96ff6debc44df5d8477bbee6eb33dc0ca",
      "bytes": 111764,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/pylib/gyp/generator/msvs.py",
      "sha256": "adade980df0828897a83917f982545042480471eeee18c4851e64f5ef997cbb6",
      "bytes": 151278,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/pylib/gyp/generator/msvs_test.py",
      "sha256": "fff6c3e2608ad9d70dc3a8d73d55546bdefb6a9c01af3feeed608b03adc20f95",
      "bytes": 1261,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/pylib/gyp/generator/ninja.py",
      "sha256": "a812b9dab9168bebd93ce45a7175466a749e4824d5f894aaee886b326dbc3a61",
      "bytes": 119530,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/pylib/gyp/generator/ninja_test.py",
      "sha256": "021d757e992e6ae651fe2dcc77019ed6dde7d0ccfdbd983402beea7398e95208",
      "bytes": 2601,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/pylib/gyp/generator/xcode.py",
      "sha256": "bcdb30428678fff56b512e598fa949485d4451d2c9797af9adda2188ae269947",
      "bytes": 66005,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/pylib/gyp/generator/xcode_test.py",
      "sha256": "000187305f0f180a71e6540de9e5665bf6fbf4f082148674b4b2a69aa55b644e",
      "bytes": 667,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/pylib/gyp/input.py",
      "sha256": "4030bc4af92103bd4fa31412185e0a647871042a5f0b0705f30e9499f33f20c9",
      "bytes": 126437,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/pylib/gyp/input_test.py",
      "sha256": "9c2bfb6b356f492faf17d12e8b5d6abbf4554f6bcf1befc4ae3fa852681aceab",
      "bytes": 3426,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/pylib/gyp/mac_tool.py",
      "sha256": "4a5c954b90a64b9da94e4988992ed90b0adbfef08dc97298b39afd786e7bb3f6",
      "bytes": 30469,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/pylib/gyp/msvs_emulation.py",
      "sha256": "dcedafe9e280dd6de3ddc11632a5e690fcac3c68f9020a203a60c3dd4c7970ca",
      "bytes": 54030,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/pylib/gyp/MSVSNew.py",
      "sha256": "30ab60341a68288aff8c8f8ae82709d55c631a3041b16a621cf677b7acbd99c1",
      "bytes": 13202,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/pylib/gyp/MSVSProject.py",
      "sha256": "c7e611ab14306def936c6dc47304b0f44d053c1e1f701081edff16659108313a",
      "bytes": 6940,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/pylib/gyp/MSVSSettings.py",
      "sha256": "13c9006c1af1bf40389eadbb1f303221544a37bcb4615be8e3d227856c40280c",
      "bytes": 45747,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/pylib/gyp/MSVSSettings_test.py",
      "sha256": "8c30d346ac54fb1d2a5f15a2f91be068336f490e1e5b648de7a5dba9eaade742",
      "bytes": 74434,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/pylib/gyp/MSVSToolFile.py",
      "sha256": "bddfddd604777a0f1b9e3fa344b8bef2a90a18fd1ed14e6b9e7142ab4785b995",
      "bytes": 1830,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/pylib/gyp/MSVSUserFile.py",
      "sha256": "97bb28014bc84494b9e0397481c9749c94f0c04f73a03a632b2cace4d0f38b29",
      "bytes": 5376,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/pylib/gyp/MSVSUtil.py",
      "sha256": "c9e2bd297d2ea44b32e6b6d7697c1aeb0709ef3728dc06b8e3d29f8e3844f1fb",
      "bytes": 10326,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/pylib/gyp/MSVSVersion.py",
      "sha256": "7e3cebf3372de53182cf59c23978d6aac33f9c06b52e1338170f684804020019",
      "bytes": 20667,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/pylib/gyp/ninja_syntax.py",
      "sha256": "ccc1e407c9743bba1192e605a1a37d768b676c04c11d9804ad0cedc04fd6eb8f",
      "bytes": 5640,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/pylib/gyp/simple_copy.py",
      "sha256": "ba81d910a604bc32b3b26a05a79a3d413bb9b1e6ac6adcaefc17b5d8c3d43ce1",
      "bytes": 1295,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/pylib/gyp/win_tool.py",
      "sha256": "601401e4e26c3f6211943b4062f9e72dbe8b21870c8b8e0ecdc66ac4152a7c27",
      "bytes": 15178,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/pylib/gyp/xcode_emulation.py",
      "sha256": "153b5cdec70c1875efbfb391df371d9650e55effaf572fc500dba14289f7cf01",
      "bytes": 82454,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/pylib/gyp/xcode_emulation_test.py",
      "sha256": "784a6188ed50125e6a74d72e5b57bba2076d0d5693c523776542a5530091cab7",
      "bytes": 1489,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/pylib/gyp/xcode_ninja.py",
      "sha256": "e2a5068a09bb4b832fe29532d9b2519be768898fbd89f72390a9c30e1020f527",
      "bytes": 12132,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/pylib/gyp/xcodeproj_file.py",
      "sha256": "2c00554661b8b6677644ae64117f0dbed21fe0c0873c3da778ea63e2ee1976c5",
      "bytes": 136274,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/pylib/gyp/xml_fix.py",
      "sha256": "f688dfce0efe63cc5c28ddd7d85a5fb9d3ecb14840ac7925f1f814d04164ca52",
      "bytes": 2244,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/pylib/packaging/__init__.py",
      "sha256": "2a00427c7cfdbe084df43b72682531c9ec1f7b956b8d50ed89ab40d66e8287b4",
      "bytes": 501,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/pylib/packaging/_elffile.py",
      "sha256": "e0cd8d28e5072aa90b5e2bc1009a7ff7d88846e1ccd500163965612520872e7b",
      "bytes": 3255,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/pylib/packaging/_manylinux.py",
      "sha256": "46aea9a570311fc5c5b4d7fab42f81fb548abaf08e0cf06f702a12ba532d6ed9",
      "bytes": 9526,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/pylib/packaging/_musllinux.py",
      "sha256": "92098118b1726e9cbceb4f7e293bf39addb30a108f598be1a790563f8257edd1",
      "bytes": 2676,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/pylib/packaging/_parser.py",
      "sha256": "e4384aff3609138538cb34a4804053e05eed4f6c59e8f931e204912b2bd79de4",
      "bytes": 10382,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/pylib/packaging/_structures.py",
      "sha256": "ab77953666d62461bf4b40e2b7f4b7028f2a42acffe4f6135c500a0597b9cabe",
      "bytes": 1431,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/pylib/packaging/_tokenizer.py",
      "sha256": "6a50ad6f05e138502614667a050fb0093485a11009db3fb2b087fbfff31327f9",
      "bytes": 5292,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/pylib/packaging/LICENSE",
      "sha256": "cad1ef5bd340d73e074ba614d26f7deaca5c7940c3d8c34852e65c4909686c48",
      "bytes": 197,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/pylib/packaging/LICENSE.APACHE",
      "sha256": "0d542e0c8804e39aa7f37eb00da5a762149dc682d7829451287e11b938e94594",
      "bytes": 10174,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/pylib/packaging/LICENSE.BSD",
      "sha256": "b70e7e9b742f1cc6f948b34c16aa39ffece94196364bc88ff0d2180f0028fac5",
      "bytes": 1344,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/pylib/packaging/markers.py",
      "sha256": "72c7d49e0186aa124fd8c46e32caeb048c0360f1e7bf5515b4b8d26c2f467c33",
      "bytes": 8202,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/pylib/packaging/metadata.py",
      "sha256": "90ac1cfb7819d967cddc11866e9d6733838a05c8744f8591a2270002be4e02bb",
      "bytes": 32535,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/pylib/packaging/py.typed",
      "sha256": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
      "bytes": 0,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/pylib/packaging/requirements.py",
      "sha256": "c2cc06e265c74a013dc38363367952be02c6a3ac98bdf1d51059d686e1265f18",
      "bytes": 2952,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/pylib/packaging/specifiers.py",
      "sha256": "7e6392c8cebfd4ab851cb5d0b117c883a2a2a806b937244f797bb1ca8a0f2f90",
      "bytes": 39969,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/pylib/packaging/tags.py",
      "sha256": "643bc6a797f1f60b432ac1ad516feadb700c213accf1efc2dd8fd048c499f96b",
      "bytes": 17850,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/pylib/packaging/utils.py",
      "sha256": "5e07663f7cb1f7ec101058ceecebcc8fd46311fe49951e4714547af6fed243d1",
      "bytes": 5268,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/pylib/packaging/version.py",
      "sha256": "5e34412cd2b5ed430380b78ff141e7ab0898dd37528b4df1150511b5e736d750",
      "bytes": 16236,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/pyproject.toml",
      "sha256": "56a752cd7a5d394d219537c041447a775b44e7cdacd2456220bb3fea6243c9dd",
      "bytes": 2981,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/release-please-config.json",
      "sha256": "e60dec6429d83a7578464d52669f17053481f1ec408671b1a56b470ee454837c",
      "bytes": 279,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/gyp/test_gyp.py",
      "sha256": "c47a7aade761cd8092d6169ac3c478ce45e2efa13bd3e8406bba715ede229264",
      "bytes": 7724,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/lib/build.js",
      "sha256": "d11b1d799a2522bd38a364bc39d36175d9e68225e0e5a36d9fc3a64e1710fbc0",
      "bytes": 6724,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/lib/clean.js",
      "sha256": "5a4f5a39a10523fd559e262bd0e4b00816b3274912de6d8637904180484b0df0",
      "bytes": 414,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/lib/configure.js",
      "sha256": "640d14e5719d34a4b5f1b7728cd94f383b0992a6089f3a133ecb9a349625d861",
      "bytes": 11915,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/lib/create-config-gypi.js",
      "sha256": "b634df313bd97a588d7e235e4043128d71571d1ea7e4282906170d58d9f6e504",
      "bytes": 4798,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/lib/download.js",
      "sha256": "d6d9f6adf3e73ca72e6879620eb37ab5603e431796d0375adcfad10ac2e826a3",
      "bytes": 2406,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/lib/find-node-directory.js",
      "sha256": "aa15c361f9703bf37d99d5d1d6b85491987bfa170d8403f3b5d50696b409d4a9",
      "bytes": 2378,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/lib/find-python.js",
      "sha256": "43ebfc0fb75165c46218d16d441de1da7898ca5339452a65aca0cb291aaa160e",
      "bytes": 10700,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/lib/Find-VisualStudio.cs",
      "sha256": "bdb19763c5d23fab534ddc945f64c4cd956584eac934f0bed96c536d3fc53502",
      "bytes": 7931,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/lib/find-visualstudio.js",
      "sha256": "9356f1e312812091a796b11b729704f91247adb71c0238223640674edafb4731",
      "bytes": 20382,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/lib/install.js",
      "sha256": "f0a0017fec48692a8eb8922b6eaff8d12f4a4eb1af4e42f23efa3e8797dff2c3",
      "bytes": 13975,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/lib/list.js",
      "sha256": "6d727b4490ec61d795c7e0ac01af7d00fb185a56c27437dfbc5ab511fe9e919c",
      "bytes": 580,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/lib/log.js",
      "sha256": "8a8be6cc02b9803fbdff0e7e7b25a50ea2e221b53ba59e146fadc6448578ba5d",
      "bytes": 3462,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/lib/node-gyp.js",
      "sha256": "6221f4df47b863bcff62ae052d53f11b7487c9ea015f787771323c4c232a2139",
      "bytes": 5651,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/lib/process-release.js",
      "sha256": "a5f87aaeb425a0e3e36aced7635a7269a6950e4ab3b38e71828bb1eeb9f4ba1c",
      "bytes": 5762,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/lib/rebuild.js",
      "sha256": "8e061a061f56fcc329eba702d5e2b08107ea1d2d0f64a29f93291e2d3fae43a9",
      "bytes": 281,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/lib/remove.js",
      "sha256": "887a83966f9739067be94d9e0aa71197fcadafc55d7e3b4cfebe62e6a4ba260f",
      "bytes": 1233,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/lib/util.js",
      "sha256": "6a806d26943f783094826cb23ff23b85b112e98f86e3d2508f16d1d5c13586d8",
      "bytes": 2228,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/LICENSE",
      "sha256": "662a1b0115251cfb29c6aed0f221f8847bc49c6365d1c53a62c9f4bccc2489c3",
      "bytes": 1102,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/macOS_Catalina_acid_test.sh",
      "sha256": "c4ff0028080b20ec21f77b06dc84f27110e3c925e0cb65553a64d08f0989ef6a",
      "bytes": 495,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/package.json",
      "sha256": "d19b594bcbf958f80ec6863a44be2c5bd3ad84077c1de5f57b18168014fa97bd",
      "bytes": 1261,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/release-please-config.json",
      "sha256": "5e3fcb59fead830b74d15d69261a11347aa205bc63fa1fae6fc957838f3cf165",
      "bytes": 2405,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/node-gyp/src/win_delay_load_hook.cc",
      "sha256": "ec2357ffdf512151c21a52326ad3396aaa650b83e5c4a31153d216a155f68ecc",
      "bytes": 1003,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/nopt/bin/nopt.js",
      "sha256": "95527c67ac7a1e294f7fcb09e648d1e454f6cfef06346a18a297173389b97d21",
      "bytes": 644,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/nopt/lib/debug.js",
      "sha256": "b05ff7b220d143feb139105e9a37d1e42705f952e862d476786bb18a801c013a",
      "bytes": 181,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/nopt/lib/nopt-lib.js",
      "sha256": "73cf6ad54ac9ea1a29b59954fc0473e28693069aede1586ac4acf64a92899d2e",
      "bytes": 14184,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/nopt/lib/nopt.js",
      "sha256": "0cd507599b0b67afaaf6aa1019dfb2ea9dc958efd1d79f5e74cd997205d65469",
      "bytes": 1127,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/nopt/lib/type-defs.js",
      "sha256": "230fee3a48e92b863c5d2d9d62e5c8de020cdb636037cd589730ddebe221c902",
      "bytes": 2030,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/nopt/LICENSE",
      "sha256": "4ec3d4c66cd87f5c8d8ad911b10f99bf27cb00cdfcff82621956e379186b016b",
      "bytes": 765,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/nopt/package.json",
      "sha256": "9e9aa0342cc3232da295048d60276edff062faf18b317b600c6ccb6bfc0d212f",
      "bytes": 1214,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/picomatch/index.js",
      "sha256": "02ae26406508e28f9ba0865be16cf77f0b33de6c4d4d49c212c4176df1d5625a",
      "bytes": 479,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/picomatch/lib/constants.js",
      "sha256": "c28f952d2758261fe5428a121add485cc25b03e0193633f649de903dc0e5e222",
      "bytes": 4552,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/picomatch/lib/parse.js",
      "sha256": "4e1e05447668aa1e94366c5516a2bcd8f5fff800cc92b1d0cd7d77703b8618c4",
      "bytes": 35144,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/picomatch/lib/picomatch.js",
      "sha256": "fdd249d7ffff853f3bb1c77eafc8baa1a437bea85a94c35877b83c6a40d9ca94",
      "bytes": 10722,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/picomatch/lib/scan.js",
      "sha256": "6457a1764e5c2e543cc3a0c00ce12ebb3764097628868198370b64d4e642aaae",
      "bytes": 9483,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/picomatch/lib/utils.js",
      "sha256": "f04604f70cd1a3f6ace80b84766468674ff38b4bb83818d2ef1d81fc6355f935",
      "bytes": 1994,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/picomatch/LICENSE",
      "sha256": "d0cd141b0c322fded5dfad1d4645bb2fedfc05b7321fe1009469638190d59ef9",
      "bytes": 1091,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/picomatch/package.json",
      "sha256": "479b70b87aa7b44344d62c8a459ae7e53dd2bca6c61f1396bd55e63e9de73fc6",
      "bytes": 1915,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/picomatch/posix.js",
      "sha256": "5695ff5dfcd5a338a40630bd506056f2950bcb08fe1cc068519cbbfe60add480",
      "bytes": 60,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/proc-log/lib/index.js",
      "sha256": "2251aa02ae472e0174b94f37bcbf433f7f2decc77edb5c7b6c56a4a3d90f7c72",
      "bytes": 3597,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/proc-log/LICENSE",
      "sha256": "dc32a0dee275e0a9aeffbc974dbf4899a30dcdc2e5ffa8934aecb69261065864",
      "bytes": 742,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/proc-log/package.json",
      "sha256": "5da081088a404333c924e99523d09d00d90f0e2c2bff011bc14f8899dff8897d",
      "bytes": 1132,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/semver/bin/semver.js",
      "sha256": "bd6c871026985937dc945011fc54b74b47f5998154216d7b8d40d5ed782e4402",
      "bytes": 4957,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/semver/classes/comparator.js",
      "sha256": "054202956430d63d5ff4599fae09760ce465b489e4f0b5ef5ce7cc7ac21157ac",
      "bytes": 3631,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/semver/classes/index.js",
      "sha256": "3bb69280c2a788d0eb16f915bb9df4dbe812075182024c753dca2283bcea1b17",
      "bytes": 143,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/semver/classes/range.js",
      "sha256": "a6544485aa8575aa9854c67e2caa67726815de18a482d0ceb8a1003244de3bc1",
      "bytes": 15647,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/semver/classes/semver.js",
      "sha256": "813b2c185512d30c9b931c2fed8140889d9d7490ceeeedc370da30337dc8ca57",
      "bytes": 9891,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/semver/functions/clean.js",
      "sha256": "4eadb0892844cf3ae295121a86163a66c73f89acd1b7f0b114ec115b4539512f",
      "bytes": 205,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/semver/functions/cmp.js",
      "sha256": "a63d74e87b73788e78e9ce0a4892b5333d6b809c0de88b31e4ed76cbf17f94b3",
      "bytes": 961,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/semver/functions/coerce.js",
      "sha256": "28a251c5ab210ddf9e97551b9f37a53329fcce91f0c3943dcbc02de1a1de915a",
      "bytes": 2004,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/semver/functions/compare-build.js",
      "sha256": "5ab651d5b40af289bd85c645a92b6d8cfe1a986dc413c797cfcc8d623d7c844c",
      "bytes": 281,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/semver/functions/compare-loose.js",
      "sha256": "07b6a3a1db0a5210ceb784c1708bd4679f3a94fc73a9c9eb349349e7070a6f78",
      "bytes": 132,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/semver/functions/compare.js",
      "sha256": "d404b5aa48aaddc8a654c5da8fb7d4443404b7948589b21ac4b045d1cee4e34c",
      "bytes": 170,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/semver/functions/diff.js",
      "sha256": "7a11fd39b987cdf06c65e928cb1aca49bec583feb86fb5c8fe47a6fd61d7de31",
      "bytes": 1423,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/semver/functions/eq.js",
      "sha256": "b6e30a7168e52723216fc163d300e2bbabf92ec0251f9ac5438bb6ccf57c8936",
      "bytes": 126,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/semver/functions/gt.js",
      "sha256": "135523704aa48cd98834dd170ee9f74f0e68043b379f32d021db11e6304c5c93",
      "bytes": 124,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/semver/functions/gte.js",
      "sha256": "991c5bbe48ecb210a562646872f05862ad9fc0d42186d85aa60bdc6fa323eb9d",
      "bytes": 127,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/semver/functions/inc.js",
      "sha256": "952069fc8690b7d3af0fe9d55f7c54fe2ac067b48c5e74f6a54f9ce19a334493",
      "bytes": 478,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/semver/functions/lt.js",
      "sha256": "1c897a9bc849320e2e9dc0f6c09555c01ee3ddf30734515d717b88ad7740ea25",
      "bytes": 124,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/semver/functions/lte.js",
      "sha256": "3a8d0b1d00423f60cc7cb810b36ab77b61330831ad237dbe73eb5cecfa412800",
      "bytes": 127,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/semver/functions/major.js",
      "sha256": "5c678480d882f511200fed2c16ec3847dfedb08a1d70328dc2f031d35d825276",
      "bytes": 136,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/semver/functions/minor.js",
      "sha256": "1b051794f1713adec2a236517196691687c35f82e0b596ff3316a78b3cc10ae6",
      "bytes": 136,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/semver/functions/neq.js",
      "sha256": "a662883751918822c162183b46b9e20d09489132f82686c92ab78bee67f3a127",
      "bytes": 128,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/semver/functions/parse.js",
      "sha256": "29a69e15b6d02fe381d573f861881a89590e9d0f0f0ca740c5f85eaf0234c4ad",
      "bytes": 331,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/semver/functions/patch.js",
      "sha256": "f02d3c1b059fe3d96ce124886f7eef321d381a95638fd3c4a8d5ccd8e76ffadd",
      "bytes": 136,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/semver/functions/prerelease.js",
      "sha256": "baedbf503d5610ad041bfb56071efa48feb331ca278295c399537d35d3ffb593",
      "bytes": 234,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/semver/functions/rcompare.js",
      "sha256": "84fa5e88adf08d15c993cac8cdec6d1a65045b7e95a9c55184230a7f807f4dc0",
      "bytes": 132,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/semver/functions/rsort.js",
      "sha256": "6ec659ce3b6c2b173c719286caab04409adba046c0917874ec3b5e36ddfbc7e3",
      "bytes": 163,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/semver/functions/satisfies.js",
      "sha256": "8cf5e122b757251671ed6c9d9680904b71cd375845853f05312e608cf2cc2946",
      "bytes": 247,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/semver/functions/sort.js",
      "sha256": "c2fe2d3ed0be8a4e9de8f02abdfd5d9c0d3bcc510d88e87af84185592882c4b8",
      "bytes": 161,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/semver/functions/truncate.js",
      "sha256": "e145833a250311927c92f424f56498c52f1b7bb1271ce6afbf006db87773ff43",
      "bytes": 1021,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/semver/functions/valid.js",
      "sha256": "0de7ea736cb7807179d46dfac09c830a9338e02b6a07db12d7040cde2def6025",
      "bytes": 176,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/semver/index.js",
      "sha256": "4b3e57d3d40e29e0706002eba113d09f35aea593578376bbeec83b777b9912ab",
      "bytes": 2691,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/semver/internal/constants.js",
      "sha256": "38a112baf27ceca0260082ff26ac2fd7a9861cab1af12dd65e720277f68e6ce9",
      "bytes": 873,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/semver/internal/debug.js",
      "sha256": "8a9f420572260f3cf944463b5090d62a60f0730589dc23a7ec4ca25e2ee41bb3",
      "bytes": 240,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/semver/internal/identifiers.js",
      "sha256": "b4916b09dc7869ae0eb05e71c855c942a3a7365f7e6e89185d59a9e45a2451d7",
      "bytes": 525,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/semver/internal/lrucache.js",
      "sha256": "14d087c87da87b6f5c36fc4cdd7d2d14077874b14a68e20fce5b6138fa2ca34f",
      "bytes": 802,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/semver/internal/parse-options.js",
      "sha256": "fdf51d0de8d5442c35a997ef58cd530d239ce206f961d14c5121354451b01d01",
      "bytes": 338,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/semver/internal/re.js",
      "sha256": "5833262888e2b5d843a69193f83c05e374818dbe55379b497819b5bf58e48cd8",
      "bytes": 8139,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/semver/LICENSE",
      "sha256": "4ec3d4c66cd87f5c8d8ad911b10f99bf27cb00cdfcff82621956e379186b016b",
      "bytes": 765,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/semver/package.json",
      "sha256": "7c94cb7f2a53c27b20d76386ec144c062894dbcc909cfabd0f728c37874b1776",
      "bytes": 1661,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/semver/preload.js",
      "sha256": "edb6808911bebcb324b2df57e5c9935149e56984ff083b74c6cfe215f5b710ba",
      "bytes": 83,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/semver/range.bnf",
      "sha256": "15d0baf9b7b98e6d862c0cb9e822d2533d2bc23e136bb75314d802ca1fcb0392",
      "bytes": 702,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/semver/ranges/gtr.js",
      "sha256": "8fadad28e36d28e93d498ac7ac20badba2a407312845eabc18e82e90a0732b19",
      "bytes": 231,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/semver/ranges/intersects.js",
      "sha256": "fe87ac5d3020010ad3ec00636dadbf0c669ff07d0f57e0a8165a8264f79a676f",
      "bytes": 224,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/semver/ranges/ltr.js",
      "sha256": "e5186fcc03018acf9be6d968755d4c49727aeb0d981d179eb568ae5fbe983038",
      "bytes": 227,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/semver/ranges/max-satisfying.js",
      "sha256": "e1a2c0d6144cc772cd20bbc8ecb9e8a3a4074e9172a3d8e794838b591cdeb416",
      "bytes": 593,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/semver/ranges/min-satisfying.js",
      "sha256": "2681abf54098aa670f12826b76a6ec77a2441186ae4243afde3be8ae4908f7ca",
      "bytes": 591,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/semver/ranges/min-version.js",
      "sha256": "fefba0a88c2bf74d5cede504b5ae50a8dc3edfd69cf0174a491e2cf3e442614b",
      "bytes": 1514,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/semver/ranges/outside.js",
      "sha256": "3a2b0b23593d2f49419c06af2af75450cab103b0c25d665d48fe5bca495a21ca",
      "bytes": 2204,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/semver/ranges/simplify.js",
      "sha256": "7b78581c13322bc68ece2088685386b2a9b51c15b94d0a2063bdf2546bd41934",
      "bytes": 1355,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/semver/ranges/subset.js",
      "sha256": "ceb9eb6ed5bfc6c7ac7af7f47e5c3535444a4e66a0e2b58cdad4bd02a35d454d",
      "bytes": 7478,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/semver/ranges/to-comparators.js",
      "sha256": "6c5e966210cff270fa2850668aaf8460fac7759f8d99f282521ef7a78f4564e9",
      "bytes": 282,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/semver/ranges/valid.js",
      "sha256": "5ef6f995af801868925940cf8df5735d565ebabb090b068695cae65218bcd3ac",
      "bytes": 326,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tar/dist/commonjs/create.js",
      "sha256": "eff71718fc8424bc6c7cda76d5d2831da98ee1457da1df19bd31b92e563d2336",
      "bytes": 2588,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tar/dist/commonjs/cwd-error.js",
      "sha256": "faae7f047a57227c0452abb13a36a6ae2bc560d64f01e300440afe392e3d1500",
      "bytes": 436,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tar/dist/commonjs/extract.js",
      "sha256": "a1aa096322772eb0b341c39fc732857f8a3efbbd4934a0c5376a9b3b18129acd",
      "bytes": 3447,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tar/dist/commonjs/get-write-flag.js",
      "sha256": "55302560e2fc123275d0b77549caa6be3a9098bbb2c3622c43c6901626c3ac6c",
      "bytes": 1491,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tar/dist/commonjs/header.js",
      "sha256": "2840b352ffdee3ed31cb30b87e51d9b96c6255ea484ced0fe2437ad46f49ad93",
      "bytes": 14094,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tar/dist/commonjs/index.js",
      "sha256": "b5c4c31f353da9fe11ed352a95c6c247946549fec4f1326e871a3cd32158f9bf",
      "bytes": 3137,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tar/dist/commonjs/index.min.js",
      "sha256": "b9464172723801f6ad013a1a819b5f31838753559d56c879202d7fb3c1326167",
      "bytes": 102753,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tar/dist/commonjs/large-numbers.js",
      "sha256": "6de856e25f82701b0ab91e1d1e39efbffc4e73530abcf0a66caba0239b4c9472",
      "bytes": 2733,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tar/dist/commonjs/list.js",
      "sha256": "9e8c59276c643ceeba7ecf786045dc611150aa440fa98b16f457b9c5aa211093",
      "bytes": 5599,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tar/dist/commonjs/make-command.js",
      "sha256": "c81990a4be4d005c5dcfc30ea1067c790c8eaefe7eb5869474022b28d1899438",
      "bytes": 1862,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tar/dist/commonjs/mkdir.js",
      "sha256": "e1664d23922cf2db56a8ccd448ca6a576c3514a4b86e55df9587cc64820597bf",
      "bytes": 6758,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tar/dist/commonjs/mode-fix.js",
      "sha256": "1a52161925617f0e0798ddd0df4fb7dc2cc8da52b3b2c60fd7b8756506d5deae",
      "bytes": 876,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tar/dist/commonjs/normalize-unicode.js",
      "sha256": "012057324bfbac08b14ff36e29fa6d27fc3377c747efc5d43f7861ecb43131e2",
      "bytes": 1159,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tar/dist/commonjs/normalize-windows-path.js",
      "sha256": "549d841aa9e43dbc1eae3ca716944b8f7ac9cdfca1638a4c9cb6370beeac97a4",
      "bytes": 615,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tar/dist/commonjs/options.js",
      "sha256": "0f17e6496cb4b75ec3a4bb25c8fe540cf08b810b1a21f86e27b62d349e9651e2",
      "bytes": 2117,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tar/dist/commonjs/pack.js",
      "sha256": "2f68f54f8a61c33d26b353bf9311f961b614c85e52b4b0b0521b2a3fc314be6c",
      "bytes": 17624,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tar/dist/commonjs/package.json",
      "sha256": "8005a3491db7d92f36ac66369861589f9c47123d3a7c71e643fc2c06168cd45a",
      "bytes": 25,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tar/dist/commonjs/parse.js",
      "sha256": "c90f3771fa6c839b81c1d24961786463a195dfb8a59a8f190da30d33e7e1f769",
      "bytes": 24607,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tar/dist/commonjs/path-reservations.js",
      "sha256": "06fe7819ab48dec1c884d0aaeba9ae75bc804e2b1b6cf1bd16436d48ab3af92b",
      "bytes": 5699,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tar/dist/commonjs/pax.js",
      "sha256": "8b3fc3092056e28a0e137c4b7ddc230465d62d27df7e9d1665df3f3acbb7c0b2",
      "bytes": 5397,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tar/dist/commonjs/process-umask.js",
      "sha256": "4d39e1ae623afaae403476b97495d2ee12d5d9aace509a88b8c92ba8fb9a7886",
      "bytes": 273,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tar/dist/commonjs/read-entry.js",
      "sha256": "3152ca6809643af42d413214fbb268736c5f60285941746a9b9d78086e666062",
      "bytes": 4394,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tar/dist/commonjs/replace.js",
      "sha256": "4aaf077d176067a5ca2e1e79a9f9394fa3cbd1cf7f1d399c86ae0ef35ce29358",
      "bytes": 7737,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tar/dist/commonjs/strip-absolute-path.js",
      "sha256": "9b92e2de39dcddfc6ad4d4755ecf308013eace0bc5d0f71fdaebf6b56bdaa2aa",
      "bytes": 1232,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tar/dist/commonjs/strip-trailing-slashes.js",
      "sha256": "9fba838d409865ec9a63def2e6fa6bb8439689ab95e46df36e3ef38cdc52324b",
      "bytes": 651,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tar/dist/commonjs/symlink-error.js",
      "sha256": "2d8e680987edb0a1a372391776e6514054f5575663823a5ef335c258d2b3eea1",
      "bytes": 528,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tar/dist/commonjs/types.js",
      "sha256": "c07cef6564bcef1f93b7a18a3e22c97b57f18229e31a3f66d15127c43b182b5c",
      "bytes": 1889,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tar/dist/commonjs/unpack.js",
      "sha256": "03515ed98f5df6ec0f23177b7016274f8cdc91b038d0dcd036f438a70e614eed",
      "bytes": 36414,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tar/dist/commonjs/update.js",
      "sha256": "9e01f728164cea547b955b6e78d4f544ae480a6a8ffeda8bfaf63b9dbbf07d12",
      "bytes": 1229,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tar/dist/commonjs/warn-method.js",
      "sha256": "eace06acccb8ab8a6372b3b7635da64e8dc7ff4577937a5d5714b61abf1c6589",
      "bytes": 927,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tar/dist/commonjs/winchars.js",
      "sha256": "05a3e6a1dee09aa7ed7040cdf84be24c624bb935b751f4a7fc5760818e5697d3",
      "bytes": 714,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tar/dist/commonjs/write-entry.js",
      "sha256": "4c62f3a5b860d40be8f075b3e9ed9aea4c274f73af5ea5accfa9e04565139d37",
      "bytes": 25354,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tar/dist/esm/create.js",
      "sha256": "ede714930e3e35d67374b6b18d98669a11f5d98cbabd2684e01e33e773ff7226",
      "bytes": 2177,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tar/dist/esm/cwd-error.js",
      "sha256": "34acd57feffcf46ad860d29f7ae94f89db149e3eea90ee7ca01c52be17fa5787",
      "bytes": 310,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tar/dist/esm/extract.js",
      "sha256": "c2781e3e025fb52495812fbaba10a35a06f19a6cb55874373b1dd14d163446ed",
      "bytes": 1692,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tar/dist/esm/get-write-flag.js",
      "sha256": "24b3460458f841db9d1a47becf59940a112060f4208b049fa5cff350837e8e8d",
      "bytes": 1205,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tar/dist/esm/header.js",
      "sha256": "8af8946db14f867fa1f1d7bb887554e31002721907b167d8bcfde5f38aaba3c7",
      "bytes": 12552,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tar/dist/esm/index.js",
      "sha256": "3f6eee754a43f2154124af1db2251ec0edd427367dd2a6edb1e4a3dc76efe3e2",
      "bytes": 647,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tar/dist/esm/index.min.js",
      "sha256": "1902a2a0b8cc2771f7e2aa0a855e9215a2d1023b1345fd504d28eb9c420c958c",
      "bytes": 84473,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tar/dist/esm/large-numbers.js",
      "sha256": "b50ed568260de6ba28af885adf8bb8e8d68d459e18576941fb6cc6ffdc54fc83",
      "bytes": 2581,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tar/dist/esm/list.js",
      "sha256": "386be389ed1d1765c37252013cfcaf8f1c06d2af5b9e9960492d555e8310382e",
      "bytes": 3607,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tar/dist/esm/make-command.js",
      "sha256": "87a0facf60b127ec94a07fe86d6e7115129057186598c9d07fe25c5519bbc9c7",
      "bytes": 1686,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tar/dist/esm/mkdir.js",
      "sha256": "2233ec5f5dfee8edd2e67190ce0d8dc87c5a6e16b6841bf38f077b4808432caa",
      "bytes": 5695,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tar/dist/esm/mode-fix.js",
      "sha256": "bb3e54bdf1b4c8a023ba2881baa3383c60100400cc5fa381f198bee00c969002",
      "bytes": 753,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tar/dist/esm/normalize-unicode.js",
      "sha256": "ebdbddce1e999d86d7790e2e6d912e7688ecc64c2a229e93058336283f68e44f",
      "bytes": 1009,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tar/dist/esm/normalize-windows-path.js",
      "sha256": "5d0e4a5e9da9f04bce100205d19ba2a45ccf699d2ab5377d467f43a5a86d1453",
      "bytes": 504,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tar/dist/esm/options.js",
      "sha256": "f44b5a92c02fe6a9f751427c243e9221a33b1c224d0135a72594897b695971c2",
      "bytes": 1639,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tar/dist/esm/pack.js",
      "sha256": "fdab2a58e536e2541edc212c6faa67406d62fd0764585448b5b5f85495f4da3a",
      "bytes": 15622,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tar/dist/esm/package.json",
      "sha256": "3ca9d4afd21425087cf31893b8f9f63c81b0b8408db5e343ca76e5f8aa26ab9a",
      "bytes": 23,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tar/dist/esm/parse.js",
      "sha256": "6528a15bef7787adca4920eacc3131263c3d4114448a60522923fdfe37a584f6",
      "bytes": 24366,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tar/dist/esm/path-reservations.js",
      "sha256": "40171e54eb27f9f0136e1a660257c7e2047d71ee04d8ba111f50892c30a5bf9d",
      "bytes": 5431,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tar/dist/esm/pax.js",
      "sha256": "de39d126f4641b18ab3558921817a0a52ec0dfbfc429a2148a6dc9d53926b0c8",
      "bytes": 5247,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tar/dist/esm/process-umask.js",
      "sha256": "585ba6ff23a5943e563d260e34968ee6869a448b892a3c153e99ccb26ee160e5",
      "bytes": 156,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tar/dist/esm/read-entry.js",
      "sha256": "f6633626a5ea499bd83701aeab9b5f6b813f9ed0e4b26027bf49060ab6b59edd",
      "bytes": 4111,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tar/dist/esm/replace.js",
      "sha256": "4f005a92f6883e3841e7e60a563387ccb5744b61c96b0e43a6e80bf668b2b602",
      "bytes": 7066,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tar/dist/esm/strip-absolute-path.js",
      "sha256": "23d2cf99644ab42b4802391fe0a9315c1ffe55ae81cb3dcd1cbf63494a47db3e",
      "bytes": 1060,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tar/dist/esm/strip-trailing-slashes.js",
      "sha256": "b834c1e75829a195ab94fb902d5c2d71e8b1dd899e83953be7873e756bdf0f3c",
      "bytes": 489,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tar/dist/esm/symlink-error.js",
      "sha256": "397c77c7391c0e4cda53bec71f860ac96806d0c70fbe3e743ae4848d0866b7f7",
      "bytes": 390,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tar/dist/esm/types.js",
      "sha256": "91a4d5ea812aae1440a53e32fee635aa8627da9d8d7b3df84fbd3d9b39f51320",
      "bytes": 1671,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tar/dist/esm/unpack.js",
      "sha256": "d95795464b13cfa4b5a9da8497d189ca946749951b73b22f4d017c4eace7dfbe",
      "bytes": 33094,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tar/dist/esm/update.js",
      "sha256": "e7379bb9d34d223208ba2f6540d007b4820c761adb13dfe317cdca34a10caecd",
      "bytes": 1006,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tar/dist/esm/warn-method.js",
      "sha256": "c3246912795352d8610dc66a94bc2802b56bd90c0638a089c053ffa9eede0b56",
      "bytes": 795,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tar/dist/esm/winchars.js",
      "sha256": "424214f5fd77dc70ecd64bf0c1c9b5480e0a9bfb0ed20abe0f62640182de384f",
      "bytes": 559,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tar/dist/esm/write-entry.js",
      "sha256": "97e5699a8ac9afdf58d1863622397acb444c424675a4b2c63d586f5afe844315",
      "bytes": 22763,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tar/package.json",
      "sha256": "31b18f4fb83ba7183f54fa9e2cdca610897108c37ba80570c2c6fa7ccf112d97",
      "bytes": 7839,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tinyglobby/dist/index.cjs",
      "sha256": "ee408bbc9665cba5800a6b2a0c89dd3df9820baa30d4ff3ba2f2435f16ce9f94",
      "bytes": 13706,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tinyglobby/dist/index.d.cts",
      "sha256": "737da79b3269b992ed4aa9a9f799ec85ae0eb5236a3f63125382e52e02e0c832",
      "bytes": 4688,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tinyglobby/dist/index.d.mts",
      "sha256": "737da79b3269b992ed4aa9a9f799ec85ae0eb5236a3f63125382e52e02e0c832",
      "bytes": 4688,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tinyglobby/dist/index.mjs",
      "sha256": "fc1bed0d5da56d4638002af7deaf6d179e183a85de648c180839ca36b21b3b50",
      "bytes": 12438,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tinyglobby/LICENSE",
      "sha256": "22c68811e174cbbfb3813d4135918df4f540959c14d872e601d9abe83d3cde8f",
      "bytes": 1076,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/tinyglobby/package.json",
      "sha256": "bd2b2af3c9ba0580f1acfb22c55d0b6d29abf683af5cac8a2efe82fa7c943f23",
      "bytes": 1761,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/index-fetch.js",
      "sha256": "3b99e505f9dc3ac97b44fb6bd77b67987b2e9c4ec936270971326de28993de8c",
      "bytes": 1558,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/index.js",
      "sha256": "c42e75dd4f2dfa00b8bafa7b85b8594ab367b2d8d22dcf72a67585f1f5ee41ee",
      "bytes": 6057,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/api/abort-signal.js",
      "sha256": "940054d80e9a82a4dda26fbe0b26ec3fb8ef2d283ab1b74febc4c698dfb40c5f",
      "bytes": 1057,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/api/api-connect.js",
      "sha256": "33229ace71010bfb1c2df21c0751d96042ec7c9d1e6b8ae3d5067fec4936e86f",
      "bytes": 2565,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/api/api-pipeline.js",
      "sha256": "1fff148ec69841c6a31c12350f9afe788a72fcb9c42a4082462987c31b345eb1",
      "bytes": 5445,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/api/api-request.js",
      "sha256": "dba2cdb7c724d7898cc11c4f2635e7cfbbdeadbf141b8b07d063c41782b126b8",
      "bytes": 5734,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/api/api-stream.js",
      "sha256": "b2687896c2a99be3ef9e6d9d4befde55158a4336eb39b5081b600aa360215145",
      "bytes": 5320,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/api/api-upgrade.js",
      "sha256": "efa17f0b83ad86f8135b1576dab81fe6b7983b29dd1223d092c0f734bd6498bd",
      "bytes": 2576,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/api/index.js",
      "sha256": "a7b4b74ea6b81b2401f8659c138df54fd3f2c45450741682c026abb4cf133096",
      "bytes": 264,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/api/readable.js",
      "sha256": "7c1ab798eb492b3b5fb55425d84d9ee115bfc39e763b4c64befeaf30937c8434",
      "bytes": 9085,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/api/util.js",
      "sha256": "8fb5cca52bdf03160ca7b4a41eab2088c6d4213240aa7f70425b3c71102df9c3",
      "bytes": 2374,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/core/connect.js",
      "sha256": "24ed02e5a7080d7873ffc0b02870393d87529f7a8b0a6783f59a00a14bd8aabb",
      "bytes": 7286,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/core/constants.js",
      "sha256": "962e2278404f3c0b0d45bf67eb1966c39dc33c04b89e2e9f4a955ac057272b82",
      "bytes": 2610,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/core/diagnostics.js",
      "sha256": "b07763ed42b4463f95fad4bc042e04c894587ce2ed52fb3c7ebd22d77f8d005c",
      "bytes": 5687,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/core/errors.js",
      "sha256": "aaca1e0952616cb68db77153f24e1c2a7da3ee54d699b507f40d88eecc62b061",
      "bytes": 12171,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/core/request.js",
      "sha256": "dd77a5199e35cdd6fc2f1962e6b18f00b11c24db43ac8025d45b2ff0113ab106",
      "bytes": 11222,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/core/symbols.js",
      "sha256": "c27f3539f763111179c821a91e391025ffc5a19bdadaf33dda85c2b6ccd796f4",
      "bytes": 2547,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/core/tree.js",
      "sha256": "7f2625127a44664c9fd531c0938b00a8fbac467e2ad15e65dbef48cb5029a385",
      "bytes": 3455,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/core/util.js",
      "sha256": "ebbaa1263b06489bc467b6b62d95b247f38006f94d17131482ea03685ef18da9",
      "bytes": 19042,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/dispatcher/agent.js",
      "sha256": "6bf76d0a07820033b65fc831839a17794a6582874d7d4573f900a5a00393a6f7",
      "bytes": 4006,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/dispatcher/balanced-pool.js",
      "sha256": "6beac1850ec97449ab72f87a623547c0152acf27b65fe3faa432a3f00782209d",
      "bytes": 5531,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/dispatcher/client-h1.js",
      "sha256": "41f689e08a3b3c9ece2b181cd77d5a4919c6db2c930dc3f8432028b500dbe6ec",
      "bytes": 40409,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/dispatcher/client-h2.js",
      "sha256": "bc902c3b3999ec669a0a5dd05204e7efba6d07f98e5af766483bb74d02487220",
      "bytes": 18623,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/dispatcher/client.js",
      "sha256": "bf715aab5839b1f88ff42eedf69eff8badb43e671c0fd7e785cfc28102e69ab6",
      "bytes": 17057,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/dispatcher/dispatcher-base.js",
      "sha256": "9a627a0e5a9e3a17d4fd3c0df1c52f5d421e995b193ae3d196984aa43eec37aa",
      "bytes": 4847,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/dispatcher/dispatcher.js",
      "sha256": "e76ad1e0766f839303d2208843df639a642edd0c218cf3bb468af6d740f70d78",
      "bytes": 1437,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/dispatcher/env-http-proxy-agent.js",
      "sha256": "96a0cbb0dc557f76e2d11c9969d8b5e840a5565bb6501dcbed2f7406be37fa5e",
      "bytes": 4479,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/dispatcher/fixed-queue.js",
      "sha256": "244a0b8b734395490f13430449e06754ec6ea7b102447d6b77e8eef3cb20c979",
      "bytes": 4222,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/dispatcher/pool-base.js",
      "sha256": "5263c1112ec9be6d7cf1548f916c6ed903379db97fb0eb6d4a7a7dd397134fc6",
      "bytes": 4558,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/dispatcher/pool-stats.js",
      "sha256": "b77bfdc27744e464ec58099982b0b1c4c2d5926e5c8e3929bf1321772b29d593",
      "bytes": 553,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/dispatcher/pool.js",
      "sha256": "5ccc278d8e34fad79d6ba2b656a73da0c6cc2225715522021c1fb4586633990b",
      "bytes": 3013,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/dispatcher/proxy-agent.js",
      "sha256": "3bf00deb1018c397adb7049a85d6a6d9a116f639441a4ac0b6b762976f5d13a9",
      "bytes": 8135,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/dispatcher/retry-agent.js",
      "sha256": "465e9429750441852cb0ed19bd669dbe2eca2b531ce6966040551a9b6550ae30",
      "bytes": 684,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/global.js",
      "sha256": "2743c8128f9004d7d82dfae10cdcae10b27806b7c5d378b94f2c50651f3b6dbc",
      "bytes": 871,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/handler/decorator-handler.js",
      "sha256": "56e55b0a2683669007c43454030b2d488f5577aeb693ef7d4b74a3b03b89fb12",
      "bytes": 858,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/handler/redirect-handler.js",
      "sha256": "35945253f190fd6149839d9b4c4d6d99d966fca3da64fc84d803aa071fb96fac",
      "bytes": 7700,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/handler/retry-handler.js",
      "sha256": "7efa2ca22dc394d7b9383d600e7cc01e403e0dfb7192ba04100299cf7423b544",
      "bytes": 10860,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/interceptor/dns.js",
      "sha256": "bca388666ef90bc4a89f125c786af07760e55c2fc8f7510bbe2f647d10b59c86",
      "bytes": 9338,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/interceptor/dump.js",
      "sha256": "3bfd2bfbd074fd23e1867b5d14962b065bdfb43ab8c83d30a4108939be651cf8",
      "bytes": 2514,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/interceptor/redirect-interceptor.js",
      "sha256": "dc4f2739771211d5dd3940d6b41abcced797de6cce9431665037ace2149ab3d9",
      "bytes": 661,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/interceptor/redirect.js",
      "sha256": "377a3f459d5652d2668d49b0668cf259e537745abb9174ad0677ecae7fb31427",
      "bytes": 588,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/interceptor/response-error.js",
      "sha256": "aece2c6eee6314919df653a0c5bcd89f9d5243241c11186546455af26722e57d",
      "bytes": 2148,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/interceptor/retry.js",
      "sha256": "488cf270ae143cc13a6ecf82ba523314fdbb17fd9ee08677c4488d4410f047c3",
      "bytes": 419,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/llhttp/.gitkeep",
      "sha256": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
      "bytes": 0,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/llhttp/constants.js",
      "sha256": "3cba0763867e5e3753f3efbf22853e42b3da26e2353ff77d06156d923eb4d8bd",
      "bytes": 11002,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/llhttp/llhttp-wasm.js",
      "sha256": "f6e22d23d57d86c1e613ce26ef1eb4e997e81663987a8127eef758ce5f20a655",
      "bytes": 64920,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/llhttp/llhttp_simd-wasm.js",
      "sha256": "1b8fa2e4e50e4f5774dde59855f0d34d112ef9900caf97aea4b2e6ec4ca9e3f1",
      "bytes": 64960,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/llhttp/utils.js",
      "sha256": "0c72c339551b7e963cada1733f0d95e059d340e0f77c5d704c3e2e5989dc86a3",
      "bytes": 394,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/mock/mock-agent.js",
      "sha256": "746d3abcf9f8e352549833e1a18e8880989e4f6541b23f90491adfa148a49663",
      "bytes": 4556,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/mock/mock-client.js",
      "sha256": "e5fa1ccc3f1025276166affcca04f04e6659631b9f106019fdd95f264061d792",
      "bytes": 1513,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/mock/mock-errors.js",
      "sha256": "30c997cff400a8cf3f0a9b2eadeb1923c6506a9743ba892f14fae741c49b9840",
      "bytes": 740,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/mock/mock-interceptor.js",
      "sha256": "0794f09f20e76f5b4d74d603e6557fe5462d09e8a7c307306bfc1efcfbbc460e",
      "bytes": 6824,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/mock/mock-pool.js",
      "sha256": "9526cb9c0a989987cdde04e15647e8fab21166b711a6b6dba4ffcc781d8045f0",
      "bytes": 1499,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/mock/mock-symbols.js",
      "sha256": "2bb0167aca64371dc350fa9ee1077018517d3fd43f0d48ed068cdfc3e0530950",
      "bytes": 769,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/mock/mock-utils.js",
      "sha256": "7979b322a2c7a47cf864eee799febae6c33cfc923d270e1d19e1e5714f698511",
      "bytes": 10783,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/mock/pending-interceptors-formatter.js",
      "sha256": "9d3737f40ba601fd7a9c6dd28a328e6797686ff08da712dcd0a0f011d5639f8d",
      "bytes": 1186,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/mock/pluralizer.js",
      "sha256": "d5d71bf11a309f2368693f40737a011ae59734ea1b45c27678c6d49f9888ef33",
      "bytes": 495,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/util/timers.js",
      "sha256": "bea70a68b7d31ea58611892c54795399d6bd3bcc4aca635081d40a5b4a564f7a",
      "bytes": 11737,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/web/cache/cache.js",
      "sha256": "4f6c667a0566e1e916dbfa23a3ede274d151433c073c131e058df42bc6e5e504",
      "bytes": 21080,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/web/cache/cachestorage.js",
      "sha256": "bf271e7d5bd3b6e3adfd7b0a092490bf17ae10dcb396f35efd3c28f7fd496567",
      "bytes": 3663,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/web/cache/symbols.js",
      "sha256": "499a6cdb7a77eb92655c929145116c0724e81caa6c351f23232f55dff05e073d",
      "bytes": 90,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/web/cache/util.js",
      "sha256": "638d5c2424a730e6d2400ac518765b3fc9b50e2b3816b65fb7b46f93c66a5930",
      "bytes": 1020,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/web/cookies/constants.js",
      "sha256": "85a24ae49eb5954e15363049bbc68967f01aae7260246daa3f49ad375429acb4",
      "bytes": 306,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/web/cookies/index.js",
      "sha256": "55f0e191558a019507374c918a08060ca5b9d5810b1bc4901cf5dd2dd5293ff5",
      "bytes": 4247,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/web/cookies/parse.js",
      "sha256": "550f060bdae6ac7d0215be1e21a9e95a439486f58883c49c38ed18b8c4aebbe4",
      "bytes": 12481,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/web/cookies/util.js",
      "sha256": "51941b2fb0a6f29e389b91236050b1fc4bdaf57051ae53145767219a9d3c5f09",
      "bytes": 9261,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/web/eventsource/eventsource-stream.js",
      "sha256": "baef83e10219926475f6d8ad2d28016561cf68e1af13f669d11d5541cc209e8d",
      "bytes": 11586,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/web/eventsource/eventsource.js",
      "sha256": "12709d9420c06126f1886d16e4c74b38a9d81023ee8e907c49cabf46bb098245",
      "bytes": 14143,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/web/eventsource/util.js",
      "sha256": "23db924c5eed1b4421811da40e78ed6dc754cdf165cc8eb6e9f1b9cf1dc424d7",
      "bytes": 788,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/web/fetch/body.js",
      "sha256": "edc925847d383b21cc286c27972075aea26cc52d2c87c0bc5bc9d0d3aa1a3543",
      "bytes": 17339,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/web/fetch/constants.js",
      "sha256": "d6c377a0ce91b947a1299845516886be525f474aedad65c7dde35f2c9301b97e",
      "bytes": 3395,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/web/fetch/data-url.js",
      "sha256": "9593cb61a2844fabc1849f6343c0fbb132d4ab885195ee1a87eeee54579dcaaa",
      "bytes": 21989,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/web/fetch/dispatcher-weakref.js",
      "sha256": "9f34f1dda099666ef510e41399c00c66da21d0482780aa82aaf6478acd687479",
      "bytes": 1083,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/web/fetch/file.js",
      "sha256": "29a79309900579667fb5df3a0629f508bde34464f2c08a29990de2fea1b21078",
      "bytes": 3300,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/web/fetch/formdata-parser.js",
      "sha256": "a99daed7968ae5031a50d020a881a97afbeb613fa1db13a98561ce8e6b8f4384",
      "bytes": 15227,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/web/fetch/formdata.js",
      "sha256": "2755aa6653a131084c3c9d8123a8556a5bb4de4102cbe291b7134a1975faa9d8",
      "bytes": 7756,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/web/fetch/global.js",
      "sha256": "6b38813b82ec8d5002009640515d35f83ab96d4f8e0066d71fae89324079a289",
      "bytes": 890,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/web/fetch/headers.js",
      "sha256": "8492ae2455a524b46017cb74e6a7353d63ea705914a4fb28e718b7b6d2de4470",
      "bytes": 20603,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/web/fetch/index.js",
      "sha256": "ec5fa587004c3d852f4e29ca4726eff00bede5cc4793dc6431b0c13ca0576c9f",
      "bytes": 82477,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/web/fetch/LICENSE",
      "sha256": "a21d6c8d3fc631198b97e57584697f7f9c37805c71c6cf9a12e4442b40394b88",
      "bytes": 1071,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/web/fetch/request.js",
      "sha256": "eed7953967abb30acad4c7ae54417e61d2c53bd3393a6f1929270865724345ad",
      "bytes": 34565,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/web/fetch/response.js",
      "sha256": "103a92a5e44ee47822ad5892b6134b8b237fad505a4ccf9fb4d2963af1c49fb6",
      "bytes": 19205,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/web/fetch/symbols.js",
      "sha256": "dd88b1385665e69fca2465f33dc99d4f69919441555e5504ce65b55019e077a7",
      "bytes": 181,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/web/fetch/util.js",
      "sha256": "85ee9fbb5c6083b521ac108959dddc4029e3549350da5e4ec10bfc2077846039",
      "bytes": 50492,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/web/fetch/webidl.js",
      "sha256": "c1197b314e62c8dc9e1cce920b59679474ccb2e74dc4779a7d33ee1ecd0c91d0",
      "bytes": 20571,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/web/fileapi/encoding.js",
      "sha256": "1104cabac55bb11fb9ce5127017f7c1b5dd816e63e2ec1a1c8eb40d7dfb4ad20",
      "bytes": 6618,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/web/fileapi/filereader.js",
      "sha256": "758f05d98581fa73a36df3efde7db45b209df05f9fbbdcd21ab90cd527c2feb2",
      "bytes": 8515,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/web/fileapi/progressevent.js",
      "sha256": "a4efded37a5b1bf78ae1a3c9a09c90615e24f1b8fcaed5310a317f43bd40ffc6",
      "bytes": 1651,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/web/fileapi/symbols.js",
      "sha256": "ac1dd1e588f039bee4917ff017e32af5eb727d68683d869acf1479828480c371",
      "bytes": 317,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/web/fileapi/util.js",
      "sha256": "7fe9d4d9e553ff05bbd5fe63e9c2d1d75acf4e47291936a513bcab712a1e907b",
      "bytes": 11526,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/web/websocket/connection.js",
      "sha256": "399bfffecf05b8c1c43db4ca7b8ed11ad2176eb0e08d89e917c643e2c38c2761",
      "bytes": 14206,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/web/websocket/constants.js",
      "sha256": "55527aeddaa0d2ca2a9b57080be9ee3113f674dce45069cf6d494d5149f8767c",
      "bytes": 1073,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/web/websocket/events.js",
      "sha256": "156e4faea9fae2665344d5af1d08efb0cb6dc1deb96a787879e411fd80d55d2d",
      "bytes": 7353,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/web/websocket/frame.js",
      "sha256": "a60454364e7f2ec5db26c6ff03d20262d35fb2f1a0b2f6d625fd644be21cb118",
      "bytes": 2307,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/web/websocket/permessage-deflate.js",
      "sha256": "808e4fc2eb52ac7e0140cc37a0efe0ed2b13f4bd0119133c569030d3fab05f25",
      "bytes": 2757,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/web/websocket/receiver.js",
      "sha256": "279aa204e739cf532ad06edd575b563c5b73ee0864ccacfd9c21ce9826378a47",
      "bytes": 16127,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/web/websocket/sender.js",
      "sha256": "147f118843cdbbb884b39b2f4644cbec6671a530018e3baedafadf04339e4dcd",
      "bytes": 2291,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/web/websocket/symbols.js",
      "sha256": "8b0fefd6a165e147061da314db2f1d5f7e870b6720835e540d0be3376c6132fa",
      "bytes": 330,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/web/websocket/util.js",
      "sha256": "0a391dc293d65810a39c9c336904e98956fba902eb4402e55556e3dbe151ac46",
      "bytes": 9355,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/lib/web/websocket/websocket.js",
      "sha256": "9b3173f0a96feda58e36293002f2341ca3c1c3340e9e1e070f149a4dcfe28f8c",
      "bytes": 18741,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/LICENSE",
      "sha256": "a6db8096b2707bc0102d256917d4d33f298ba36d8c3f25de067a2b5bb379db27",
      "bytes": 1090,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/package.json",
      "sha256": "f529aa8b145f6e270a72700ceb6c8dd3b87f742313078ea9f774a9d2ce8c865d",
      "bytes": 6290,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/undici/scripts/strip-comments.js",
      "sha256": "dda77b813446eda232f84105181302ab39e3e58c8d92d3350ac1ce3eace5f4c4",
      "bytes": 260,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/v8-compile-cache/LICENSE",
      "sha256": "c77674258a3fdf3036a5d13d2aecd30d7a25aa6191cb0a9a7dd45b975dc7fe69",
      "bytes": 1080,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/v8-compile-cache/package.json",
      "sha256": "0fa03730c3968110f6f8c088ebfba37a43190690a964fb8bcc6e485adbe5f1cc",
      "bytes": 832,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/v8-compile-cache/v8-compile-cache.js",
      "sha256": "3d8a6789f59d1deb9d58bbad5dd4150c69b13c15fda4b991bd0704599d7ddd43",
      "bytes": 10837,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/which/bin/which.js",
      "sha256": "7c37493fcac8af6526f51f4d83606a733ca7ca3c7b943d1e9239bb30346b019d",
      "bytes": 960,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/which/lib/index.js",
      "sha256": "d9d32d45a01826692a1a0c70d40c6182aa0405713166ecd4a3bd3df68d3b90bf",
      "bytes": 3129,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/which/LICENSE",
      "sha256": "4ec3d4c66cd87f5c8d8ad911b10f99bf27cb00cdfcff82621956e379186b016b",
      "bytes": 765,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/which/package.json",
      "sha256": "a587811d0c730bfc8be6a78adc62c8ca4d55f80a253c88d5d827fd5b3e97703d",
      "bytes": 1245,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/yallist/dist/commonjs/index.js",
      "sha256": "ba0199785e5dd531ee9fe5b8a8f9903f85d3787de85e9c8ba358ea38fb367e44",
      "bytes": 9914,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/yallist/dist/commonjs/package.json",
      "sha256": "8005a3491db7d92f36ac66369861589f9c47123d3a7c71e643fc2c06168cd45a",
      "bytes": 25,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/yallist/dist/esm/index.js",
      "sha256": "ddc4c2d4f599190677c553dfa9e4db28338283645508cb3db35aa2c6282b7af7",
      "bytes": 9762,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/yallist/dist/esm/package.json",
      "sha256": "3ca9d4afd21425087cf31893b8f9f63c81b0b8408db5e343ca76e5f8aa26ab9a",
      "bytes": 23,
      "selection": "unchanged"
    },
    {
      "path": "dist/node_modules/yallist/package.json",
      "sha256": "1b9d47057ce39814531ff93f668823b4fa03e7d23945449c274a1ff6d4cc297f",
      "bytes": 1618,
      "selection": "unchanged"
    },
    {
      "path": "dist/pnpm.mjs",
      "sha256": "ddc64218bc85fb88d28b5def06eb01fafb39cb67f7a9465ce736456658cc7f11",
      "bytes": 13891041,
      "selection": "unchanged"
    },
    {
      "path": "dist/pnpmrc",
      "sha256": "6985e6e532eed7e6fd69a7d560ac5a5753d74a9f9709627092d5be02afddd84e",
      "bytes": 161,
      "selection": "unchanged"
    },
    {
      "path": "dist/templates/completion.bash",
      "sha256": "20beebf887ab851f23e2d950913f5a78d2c27ddac782108d81febe1e3027b155",
      "bytes": 994,
      "selection": "unchanged"
    },
    {
      "path": "dist/templates/completion.fish",
      "sha256": "6aaa8f0ef03f6aee183ac70cd7a0d822f9dee96f7bf8f3b4fd498040b3a4958a",
      "bytes": 672,
      "selection": "unchanged"
    },
    {
      "path": "dist/templates/completion.ps1",
      "sha256": "34e3cc6406959fdc32ec8f41ed8d86723f151238aae64a38d25ba0d52b5463b8",
      "bytes": 7685,
      "selection": "unchanged"
    },
    {
      "path": "dist/templates/completion.zsh",
      "sha256": "0553017483862f22121af97b6cf0a3333cdfdd7d3a94c683317d09cacf2f02cd",
      "bytes": 839,
      "selection": "unchanged"
    },
    {
      "path": "dist/vendor/fastlist-0.3.0-x64.exe",
      "sha256": "d1a71f9ac1728082c1b276392725c3e010b98714888579b99152e401abedbf11",
      "bytes": 271872,
      "selection": "excluded_native_family"
    },
    {
      "path": "dist/vendor/fastlist-0.3.0-x86.exe",
      "sha256": "017411f3b0b5c0402cc3b2cb87c32c6fc71abd82e5b17ea6108990096c75a65d",
      "bytes": 215040,
      "selection": "excluded_native_family"
    },
    {
      "path": "dist/worker.js",
      "sha256": "81648f68f288deeace8ca799f1d1021bda7493cdc8afa2582600de9129fd9a82",
      "bytes": 443263,
      "selection": "unchanged"
    },
    {
      "path": "LICENSE",
      "sha256": "e0a867ff513ea7be2a0ddc339ac6a031e459a38668e077b8f0e649544062f9f2",
      "bytes": 1170,
      "selection": "unchanged"
    },
    {
      "path": "package.json",
      "sha256": "81ad221a244294ce48eb2b6d7185e316e66de890b82aeb4ec0c3524ae9e4cddd",
      "bytes": 2349,
      "selection": "unchanged"
    },
    {
      "path": "README.md",
      "sha256": "c3c7b6016fb2a2d55e20785dbdad200c9bcf604f08c8c2dd819e206bbecdcafe",
      "bytes": 14233,
      "selection": "unchanged"
    }
  ]
}
````

### FILE: `pnpm_artifact_selection/select_artifact.py`
```yaml
block_id: "PNPM-ARTIFACT-SELECTION-GATE:select_artifact.py:v1"
operation: CREATE
provenance: AUTHORED
source: "local, bounded projection and consumer planning governed by canonical acquisition and execution contracts"
license: "LicenseRef-Workspace-Owner"
sha256: "c3357136987d8cde019aba907f39bc23f6e1ef371169f9e9b98122a097b0f640"
variables: []
secrets_allowed: false
```
````python
#!/usr/bin/env python3
"""AUTHORED local projection of a pinned pnpm archive; never executes its bytes."""
from __future__ import annotations
import argparse
import hashlib
import io
import json
import os
from pathlib import Path, PurePosixPath
import re
import shutil
import stat
import sys
import tarfile
import tempfile

POLICY_SHA256 = "93c27e601ccc04bbadb755483186e96d3aca0bad6b37063c668e934e55bc43f6"
RECEIPT = "selection-receipt.json"
REPARSE = 0x400

class SelectionError(ValueError):
    pass

def require(condition, message):
    if not condition:
        raise SelectionError(message)

def digest(data):
    return hashlib.sha256(data).hexdigest()

def json_bytes(value):
    return (json.dumps(value, sort_keys=True, ensure_ascii=False, indent=2) + "\n").encode("utf-8")

def unique(pairs):
    result = {}
    for key, value in pairs:
        require(key not in result, "duplicate JSON key")
        result[key] = value
    return result

def decode(data):
    return json.loads(data.decode("utf-8"), object_pairs_hook=unique)

def no_reparse(path):
    path = Path(os.path.abspath(path))
    for item in [*reversed(path.parents), path]:
        try:
            info = item.lstat()
        except FileNotFoundError:
            continue
        require(not stat.S_ISLNK(info.st_mode) and not (getattr(info, "st_file_attributes", 0) & REPARSE), "reparse point rejected")
    return path

def read_bounded(path, limit):
    path = no_reparse(path)
    with path.open("rb") as stream:
        info = os.fstat(stream.fileno())
        require(stat.S_ISREG(info.st_mode) and info.st_size <= limit, "input size/type rejected")
        data = stream.read(limit + 1)
    require(len(data) <= limit, "input read budget exceeded")
    return data

def safe_name(name):
    require(isinstance(name, str) and name and "\\" not in name, "unsafe member name")
    parts = name.split("/")
    require(not name.startswith("/") and all(p not in ("", ".", "..") for p in parts), "unsafe member path")
    for part in parts:
        require(not re.search(r'[\x00-\x1f<>:"|?*]', part) and not part.endswith((".", " ")), "unsafe Windows member name")
        require(part.split(".")[0].upper() not in {"CON", "PRN", "AUX", "NUL", *[f"COM{i}" for i in range(1, 10)], *[f"LPT{i}" for i in range(1, 10)]}, "reserved Windows member name")
    require(str(PurePosixPath(name)) == name, "noncanonical member name")
    return name

def load_policy():
    raw = read_bounded(Path(__file__).with_name("selection-policy.json"), 262144)
    require(digest(raw) == POLICY_SHA256, "selection policy identity mismatch")
    policy = decode(raw)
    require(policy["schema"] == "elite-pnpm-selection-policy/v1", "policy schema rejected")
    entries = policy["files"]
    require(len(entries) == 455, "policy member count rejected")
    folded = set()
    for item in entries:
        name = safe_name(item["path"])
        require(name.casefold() not in folded, "policy case collision")
        folded.add(name.casefold())
        require(item["selection"] in ("unchanged", "excluded_native_family"), "policy selection rejected")
        require(type(item["bytes"]) is int and 0 <= item["bytes"] <= 22000000, "policy size rejected")
        require(re.fullmatch("[0-9a-f]{64}", item["sha256"]), "policy digest rejected")
    require(sum(x["selection"] == "unchanged" for x in entries) == 442, "selected count rejected")
    return policy

def acquisition_binding(raw, policy):
    value = decode(raw)
    require(isinstance(value, dict), "acquisition receipt must be object")
    for key, expected in policy["receipt_bindings"].items():
        require(type(value.get(key)) is type(expected) and value[key] == expected, "acquisition binding rejected: " + key)
    for key in ("profile_sha256", "approval_sha256", "lock_sha256"):
        require(isinstance(value.get(key), str) and re.fullmatch("[0-9a-f]{64}", value[key]), "acquisition reference missing: " + key)
    return digest(raw)

def make_receipt(policy, acquisition_sha):
    return {"schema": "elite-pnpm-selection-receipt/v1", "policy_sha256": POLICY_SHA256,
            "artifact_sha256": policy["archive_sha256"], "acquisition_receipt_sha256": acquisition_sha,
            "provenance": "ADAPTED_SELECTION_UNCHANGED_FILES",
            "selected_files": sum(x["selection"] == "unchanged" for x in policy["files"]),
            "excluded_files": sum(x["selection"] != "unchanged" for x in policy["files"]), "installed": False, "executed": False,
            "runtime_admitted": False, "redistribution_admitted": False}

def project_payload(payload, entries, destination):
    """Internal parser, tested with synthetic archives; public CLI uses only locked bytes."""
    expected = {"package/" + item["path"]: item for item in entries}
    seen = set()
    with tarfile.open(fileobj=io.BytesIO(payload), mode="r:gz") as archive:
        for member in archive:
            require(member.name in expected and member.name not in seen, "unexpected/duplicate archive member")
            seen.add(member.name)
            item = expected[member.name]
            require(member.isreg() and not member.sparse and member.size == item["bytes"], "member type/size rejected")
            name = safe_name(item["path"])
            stream = archive.extractfile(member)
            require(stream is not None, "member content absent")
            with stream:
                data = stream.read(item["bytes"] + 1)
            require(len(data) == item["bytes"] and digest(data) == item["sha256"], "member digest rejected")
            if item["selection"] == "unchanged":
                out = destination.joinpath(*name.split("/"))
                out.parent.mkdir(parents=True, exist_ok=True)
                with out.open("xb") as writer:
                    writer.write(data)
    require(seen == set(expected), "archive is missing locked members")

def verify(target, policy):
    target = no_reparse(target)
    expected = {"payload/" + x["path"]: x for x in policy["files"] if x["selection"] == "unchanged"}
    allowed_dirs = {"payload"}
    for rel in expected:
        allowed_dirs.update(str(x) for x in PurePosixPath(rel).parents if str(x) != ".")
    found = set()
    def walk(directory):
        with os.scandir(directory) as items:
            for item in items:
                path = Path(item.path)
                rel = path.relative_to(target).as_posix()
                no_reparse(path)
                if item.is_dir(follow_symlinks=False):
                    require(rel in allowed_dirs, "unexpected output directory")
                    walk(path)
                else:
                    require(rel in expected or rel == RECEIPT, "unexpected output file")
                    found.add(rel)
                    if rel in expected:
                        entry = expected[rel]
                        data = read_bounded(path, entry["bytes"])
                        require(len(data) == entry["bytes"] and digest(data) == entry["sha256"], "output member changed")
    walk(target)
    require(found == set(expected) | {RECEIPT}, "output files missing")
    raw = read_bounded(target / RECEIPT, 4096)
    receipt = decode(raw)
    require(isinstance(receipt, dict), "selection receipt must be object")
    source_sha = receipt.get("acquisition_receipt_sha256")
    require(isinstance(source_sha, str) and re.fullmatch("[0-9a-f]{64}", source_sha), "source receipt hash rejected")
    require(raw == json_bytes(make_receipt(policy, source_sha)), "selection receipt changed")
    return receipt

def materialize(artifact, acquisition_receipt, acquisition_sha256, target, policy):
    require(os.name == "nt", "publication supported only on Windows")
    target = no_reparse(target)
    require(target.parent.is_dir() and not os.path.lexists(target), "target must be absent with existing parent")
    raw_receipt = read_bounded(acquisition_receipt, 1048576)
    require(digest(raw_receipt) == acquisition_sha256, "acquisition receipt identity mismatch")
    acquisition_sha = acquisition_binding(raw_receipt, policy)
    payload = read_bounded(artifact, policy["archive_bytes"])
    require(len(payload) == policy["archive_bytes"] and digest(payload) == policy["archive_sha256"], "archive identity mismatch")
    staging = Path(tempfile.mkdtemp(prefix=".pnpm-selection-", dir=target.parent))
    try:
        (staging / "payload").mkdir()
        project_payload(payload, policy["files"], staging / "payload")
        (staging / RECEIPT).write_bytes(json_bytes(make_receipt(policy, acquisition_sha)))
        verify(staging, policy)
        no_reparse(target.parent)
        # Windows rename refuses any occupied destination, including an empty directory.
        os.rename(staging, target)
    finally:
        if staging.exists():
            # Only this invocation's exact newly-created sibling may be removed.
            require(staging.parent == target.parent and staging.name.startswith(".pnpm-selection-"), "cleanup scope rejected")
            shutil.rmtree(staging)
    return verify(target, policy)

def main(argv=None):
    parser = argparse.ArgumentParser(description=__doc__)
    sub = parser.add_subparsers(dest="action", required=True)
    create = sub.add_parser("materialize", help="Create a candidate selection, never install or execute")
    for arg in ("artifact", "acquisition-receipt", "acquisition-receipt-sha256", "target"):
        create.add_argument("--" + arg, required=True)
    check = sub.add_parser("verify", help="Verify exact candidate bytes, not runtime admission")
    check.add_argument("--target", required=True)
    args = parser.parse_args(argv)
    try:
        policy = load_policy()
        if args.action == "materialize":
            materialize(args.artifact, args.acquisition_receipt, args.acquisition_receipt_sha256, args.target, policy)
        else:
            verify(args.target, policy)
        print("PNPM_SELECTION_BYTES_PASS runtime_admitted=false redistribution_admitted=false")
        return 0
    except (SelectionError, OSError, UnicodeError, ValueError, KeyError, TypeError, RecursionError, tarfile.TarError) as exc:
        print("PNPM_SELECTION_BLOCKED: " + str(exc), file=sys.stderr)
        return 2

if __name__ == "__main__":
    raise SystemExit(main())
````

### FILE: `pnpm_artifact_selection/test_selection.py`
```yaml
block_id: "PNPM-ARTIFACT-SELECTION-GATE:test_selection.py:v1"
operation: CREATE
provenance: AUTHORED
source: "local, bounded projection and consumer planning governed by canonical acquisition and execution contracts"
license: "LicenseRef-Workspace-Owner"
sha256: "561ea3cc314a7a81779c04541b2fd4e201227750dd690c0d0bf87a7834008a5d"
variables: []
secrets_allowed: false
```
````python
"""AUTHORED safety regressions; tiny synthetic inputs never represent pnpm admission."""
import concurrent.futures
import contextlib
import io
import json
import os
from pathlib import Path
import subprocess
import tarfile
import tempfile
import unittest
from unittest import mock
import select_artifact as m

def tar_bytes(items):
    out = io.BytesIO()
    with tarfile.open(fileobj=out, mode="w:gz") as tar:
        for name, data, kind in items:
            info = tarfile.TarInfo(name); info.type = kind; info.size = len(data)
            if kind != tarfile.REGTYPE:
                info.linkname = "../../outside"
            tar.addfile(info, io.BytesIO(data))
    return out.getvalue()

class SelectionTests(unittest.TestCase):
    def setUp(self):
        self.tmp = tempfile.TemporaryDirectory(prefix="elite-selection-test-")
        self.root = Path(self.tmp.name)
        self.items = [("package/bin/tool.js", b"fixture only", tarfile.REGTYPE),
                      ("package/LICENSE", b"synthetic notice", tarfile.REGTYPE),
                      ("package/native.exe", b"not executable", tarfile.REGTYPE)]
        self.payload = tar_bytes(self.items)
        self.entries = [{"path": n[8:], "bytes": len(data), "sha256": m.digest(data),
                         "selection": "excluded_native_family" if n.endswith(".exe") else "unchanged"}
                        for n, data, _ in self.items]
        self.policy = {"archive_sha256": m.digest(self.payload), "archive_bytes": len(self.payload),
                       "files": self.entries, "receipt_bindings": {"id": "synthetic", "executed": False}}
        self.acquisition = {"id": "synthetic", "executed": False,
                            "profile_sha256": "a" * 64, "approval_sha256": "b" * 64, "lock_sha256": "c" * 64}
        self.artifact = self.root / "fixture.tgz"; self.artifact.write_bytes(self.payload)
        self.receipt = self.root / "acquisition.json"; self.receipt.write_bytes(m.json_bytes(self.acquisition))
        self.target = self.root / "selected"

    def tearDown(self):
        self.tmp.cleanup()

    def build(self, target=None):
        return m.materialize(self.artifact, self.receipt, m.digest(self.receipt.read_bytes()), target or self.target, self.policy)

    def test_real_locked_policy(self):
        policy = m.load_policy()
        self.assertEqual(len(policy["files"]), 455)
        self.assertEqual(sum(x["selection"] == "unchanged" for x in policy["files"]), 442)

    def test_projection_deterministic_and_exclusion(self):
        self.build(); other = self.root / "other"; self.build(other)
        left = {p.relative_to(self.target).as_posix(): p.read_bytes() for p in self.target.rglob("*") if p.is_file()}
        right = {p.relative_to(other).as_posix(): p.read_bytes() for p in other.rglob("*") if p.is_file()}
        self.assertEqual(left, right)
        self.assertNotIn("payload/native.exe", left)
        self.assertFalse(m.verify(self.target, self.policy)["runtime_admitted"])

    def test_occupied_targets_untouched(self):
        for kind in ("file", "directory"):
            with self.subTest(kind=kind):
                target = self.root / kind
                if kind == "file": target.write_bytes(b"owner")
                else: target.mkdir()
                with self.assertRaises(m.SelectionError): self.build(target)
                self.assertTrue(target.exists())
                if kind == "file": self.assertEqual(target.read_bytes(), b"owner")

    def test_archive_and_receipt_tamper_before_output(self):
        for file in (self.artifact, self.receipt):
            before = file.read_bytes()
            file.write_bytes(before + b"changed")
            with self.assertRaises((m.SelectionError, ValueError)): self.build()
            self.assertFalse(self.target.exists())
            file.write_bytes(before)

    def test_receipt_hash_is_required(self):
        with self.assertRaises(m.SelectionError):
            m.materialize(self.artifact, self.receipt, "0" * 64, self.target, self.policy)
        self.assertFalse(self.target.exists())

    def test_acquisition_bindings_and_no_admission(self):
        for key, value in (("id", "wrong"), ("executed", True), ("executed", 0), ("approval_sha256", "")):
            receipt = dict(self.acquisition); receipt[key] = value
            self.receipt.write_bytes(m.json_bytes(receipt))
            with self.subTest(key=key, value=value), self.assertRaises(m.SelectionError): self.build()
            self.assertFalse(self.target.exists())

    def test_duplicate_json_rejected(self):
        with self.assertRaises(m.SelectionError): m.decode(b'{"id":1,"id":2}')

    def test_tar_duplicate_unexpected_missing(self):
        for index, items in enumerate((self.items + [self.items[0]], self.items + [("package/extra", b"x", tarfile.REGTYPE)], self.items[:-1])):
            with self.subTest(count=len(items)), self.assertRaises(m.SelectionError):
                m.project_payload(tar_bytes(items), self.entries, self.root / ("parser-" + str(index)))

    def test_tar_links_and_digest_mismatch(self):
        for kind in (tarfile.SYMTYPE, tarfile.LNKTYPE, tarfile.FIFOTYPE):
            items = [(self.items[0][0], self.items[0][1], kind), *self.items[1:]]
            with self.subTest(kind=kind), self.assertRaises(m.SelectionError):
                m.project_payload(tar_bytes(items), self.entries, self.root / "parser")
        entries = [dict(e) for e in self.entries]; entries[2]["sha256"] = "0" * 64
        with self.assertRaises(m.SelectionError): m.project_payload(self.payload, entries, self.root / "digest")

    def test_unsafe_names(self):
        for name in ("../outside", "/absolute", "x\\y", "a//b", "a/./b", "a/b.", "a/b ", "a:stream", "NUL.txt", "com1", "x\x00y"):
            with self.subTest(name=name), self.assertRaises(m.SelectionError): m.safe_name(name)

    def test_extraction_fault_cleans_only_staging(self):
        before = self.artifact.read_bytes()
        def fault(_payload, _entries, target):
            (target / "partial").write_bytes(b"partial")
            raise OSError("synthetic disk failure")
        with mock.patch.object(m, "project_payload", side_effect=fault), self.assertRaises(OSError): self.build()
        self.assertFalse(self.target.exists())
        self.assertEqual(list(self.root.glob(".pnpm-selection-*")), [])
        self.assertEqual(self.artifact.read_bytes(), before)

    def test_late_competitor_survives(self):
        rename = m.os.rename
        def compete(source, target):
            target.mkdir(); (target / "owner").write_bytes(b"other publisher")
            return rename(source, target)
        with mock.patch.object(m.os, "rename", side_effect=compete), self.assertRaises(OSError): self.build()
        self.assertEqual((self.target / "owner").read_bytes(), b"other publisher")
        self.assertEqual(list(self.root.glob(".pnpm-selection-*")), [])

    def test_concurrent_publication_one_winner(self):
        def attempt():
            try: self.build(); return "ok"
            except (OSError, m.SelectionError): return "blocked"
        with concurrent.futures.ThreadPoolExecutor(max_workers=2) as pool:
            results = list(pool.map(lambda _: attempt(), range(2)))
        self.assertEqual(sorted(results), ["blocked", "ok"])
        m.verify(self.target, self.policy)

    def test_verify_detects_missing_extra_and_changed(self):
        self.build()
        file = self.target / "payload/LICENSE"; original = file.read_bytes()
        for mutate in (lambda: file.write_bytes(b"changed"), lambda: file.unlink()):
            mutate()
            with self.assertRaises(m.SelectionError): m.verify(self.target, self.policy)
            file.write_bytes(original)
        (self.target / "extra").write_bytes(b"unexpected")
        with self.assertRaises(m.SelectionError): m.verify(self.target, self.policy)

    def test_false_admission_receipt_rejected(self):
        self.build(); file = self.target / m.RECEIPT; data = m.decode(file.read_bytes()); data["runtime_admitted"] = True
        file.write_bytes(m.json_bytes(data))
        with self.assertRaises(m.SelectionError): m.verify(self.target, self.policy)

    def test_junction_input_and_target_parent_rejected(self):
        actual = self.root / "real"; actual.mkdir(); link = self.root / "junction"
        result = subprocess.run(["cmd.exe", "/d", "/c", "mklink", "/J", str(link), str(actual)], capture_output=True, creationflags=0x08000000)
        self.assertEqual(result.returncode, 0, result.stderr.decode(errors="replace"))
        try:
            with self.assertRaises(m.SelectionError): m.no_reparse(link / "missing")
            with self.assertRaises(m.SelectionError): self.build(link / "target")
            self.assertEqual(list(actual.iterdir()), [])
        finally:
            link.rmdir()

    def test_cli_policy_tamper_and_missing_paths_no_traceback(self):
        with mock.patch.object(m, "POLICY_SHA256", "0" * 64), contextlib.redirect_stderr(io.StringIO()) as err:
            self.assertEqual(m.main(["verify", "--target", str(self.target)]), 2)
        self.assertNotIn("Traceback", err.getvalue())
        with contextlib.redirect_stderr(io.StringIO()) as err:
            self.assertEqual(m.main(["verify", "--target", str(self.target)]), 2)
        self.assertNotIn("Traceback", err.getvalue())

if __name__ == "__main__":
    unittest.main()
````

### FILE: `pnpm_artifact_selection/plan_install.py`
```yaml
block_id: "PNPM-ARTIFACT-SELECTION-GATE:plan_install.py:v395"
operation: CREATE
provenance: ADAPTED
source: "AUTHORED planner; VERBATIM QRCode header from gtanner/qrcode-terminal 90f66cf5c6b10bcb4358df96a9580f9eb383307b vendor/QRCode/index.js (blob 10eb8eb0a06aa50d0d3c508f886a7728f52e5d98); MIT permission text from https://opensource.org/license/mit, copyright populated from retained author header; no upstream algorithm executed or copied as runtime; VERBATIM semver-utils1.1.4 original package/LICENSE SHA25680b98c1b20edfc51abd2c802ee7a1d3c5561151d108a24f19c207b6935beaed8 from exact npm artifact fa6980458971864f25f2afc7d7f6632b015d8b767296d48421e6ff13221bd2b3; full MIT option chosen, entire dual-license text retained"
license: "LicenseRef-Workspace-Owner AND MIT"
sha256: "5e13ccfa6228278ff66bd992d0db6415750fcf900157e25467e03d4aee535d67"
variables: []
secrets_allowed: false
```
````python
#!/usr/bin/env python3
"""Prepare/verify a pinned offline-install recipe. Never execute or admit pnpm."""
from __future__ import annotations
import argparse
import base64
import binascii
import os
from pathlib import Path
import shutil
import sys
import tempfile
from select_artifact import (SelectionError, require, digest, json_bytes, decode,
    no_reparse, read_bounded, load_policy, acquisition_binding, verify, safe_name)

NODE_SHA256 = '5c976096e04e5c2c1f091938926234cc9fbebfe9787ddd149351b3b0ecc707b5'
PROFILES = {'enterprise-web': {'package.json': 'acc4e4bd62c13cf48e071e933266d82d067cab849f92474d71a533ebf6508951', 'pnpm-lock.yaml': 'd6c73eb82a56c7e8d21ee89ffb402d0adb2f0c15db9c33bb8ec7035775119928', 'pnpm-workspace.yaml': 'f52887318fe556dbdb3574622f5c6248483e843edc3b61ad98b6ae3b351571ba'}, 'playwright': {'package.json': 'eabb3e705b59759748724b2be3e8720184d26b570ab51e6c8a2ee8a0bdf41b90', 'pnpm-lock.yaml': '63eec1e3d5bac29965e850bf751269c7a8b317929b570205a1e223544f80a470'}, 'lighthouse': {'package.json': 'dcdb7fd7c346f73543643070f14c15578a3f425245b24a5398b314d26366b7bb', 'pnpm-lock.yaml': 'c690da3a314e5275d68c4ea4a88e700dbd9a58b0da6fcbc25aec2bdc486dd9aa'}}
RECIPE = 'install-plan.json'
NOTICE = 'BLUEOAK-NOTICE.md'
QRCODE_NOTICE = 'QRCODE-NOTICE.md'

SEMVER_NOTICE = 'SEMVER-UTILS-NOTICE.md'

RETAINED_NOTICE = 'PNPM-RETAINED-NOTICES.md'

NEXT_PATH_NOTICE = 'NEXT-PATH-MPL-SOURCE.md'
NEXT_PATH_DATA_SHA256 = '581c15a6e63ecc37d3c4f4fde4e701ed25532ffc141c519dd88919baf8685dd7'
NEXT_PATH_REGION = {'start': 7619794, 'end': 7620528, 'sha256': 'd8b58764b8a66de031a13ad8b76cccbb936e49f69c7df6d6bda306ced71fab83'}

CATALOG_SHA256 = 'e4a906897aeae3dec44d63beb40514329353f41859c8d380c126df8b2d599c1e'

SEMVER_LICENSE = 'Copyright 2013 AJ ONeal\n\nThis is open source software; you can redistribute it and/or modify it under the\nterms of either:\n\n   a) the "MIT License"\n   b) the "Apache-2.0 License"\n\nMIT License\n\n   Permission is hereby granted, free of charge, to any person obtaining a copy\n   of this software and associated documentation files (the "Software"), to deal\n   in the Software without restriction, including without limitation the rights\n   to use, copy, modify, merge, publish, distribute, sublicense, and/or sell\n   copies of the Software, and to permit persons to whom the Software is\n   furnished to do so, subject to the following conditions:\n\n   The above copyright notice and this permission notice shall be included in all\n   copies or substantial portions of the Software.\n\n   THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR\n   IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,\n   FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE\n   AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER\n   LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,\n   OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE\n   SOFTWARE.\n\nApache-2.0 License Summary\n\n   Licensed under the Apache License, Version 2.0 (the "License");\n   you may not use this file except in compliance with the License.\n   You may obtain a copy of the License at\n\n     http://www.apache.org/licenses/LICENSE-2.0\n\n   Unless required by applicable law or agreed to in writing, software\n   distributed under the License is distributed on an "AS IS" BASIS,\n   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.\n   See the License for the specific language governing permissions and\n   limitations under the License.\n'
SEMVER_LICENSE_SHA256 = '80b98c1b20edfc51abd2c802ee7a1d3c5561151d108a24f19c207b6935beaed8'
SEMVER_REGION = {'start': 1347913, 'end': 1350984, 'region_sha256': 'aed20138b14658d55690a793c31b50d4f97e3c0d5260e323b7a74fb8279d9cf2', 'marker': '// ../../../../../setup-pnpm/store/v11/links/@/semver-utils/1.1.4/bbbf6bcd6e2cad096d8c5e4cd94c8815535d76e88f7c54b6d35929d441dd7625/node_modules/semver-utils/semver-utils.js', 'source_member_sha256': '1f85bceb6f2e6cdbefadd5bf9969db8637c754c9d6553f3d7663dca6ade7c293'}

QRCODE_BUNDLE = 'dist/pnpm.mjs'
QRCODE_BUNDLE_SHA256 = 'ddc64218bc85fb88d28b5def06eb01fafb39cb67f7a9465ce736456658cc7f11'
QRCODE_MODULES = [{'path': 'vendor/QRCode/QRMode.js', 'source_sha256': '6b8ec04257a2d23b01e8189815292dca3651b38ea0a8f9c975b3c1d18dfb1b01', 'source_git_blob': '050c8a3037f60d1e034764035394c55bca34f3c3', 'start': 2065080, 'end': 2065663, 'region_sha256': '008dd02912868c88714e36bce44f2b472b7812d417183053b7fb3ae47d7126e5'}, {'path': 'vendor/QRCode/QR8bitByte.js', 'source_sha256': 'a67f0b2239db81b1fc1dfd8e169a879d7075dd79d0ae00dc155e9c3bac595891', 'source_git_blob': '94bf74f0e897de3d7c432d0d6f6422dcee8e471a', 'start': 2065664, 'end': 2066555, 'region_sha256': '3ec44883fe4fab2eb98f4f4ae841b836d19fb7fa47523fb4f019bbca23a77e60'}, {'path': 'vendor/QRCode/QRMath.js', 'source_sha256': '481fe65cd1a049a3cdd659ff20c45eb4e0cb2db285fa63a42478727e1b051667', 'source_git_blob': '8f4a0370ebb32a2e93ae053fb475732215a2cd02', 'start': 2066556, 'end': 2067839, 'region_sha256': 'b617520021a3eb20389c9483ab7fb517eaa5c5720e1058e60c5841302783a997'}, {'path': 'vendor/QRCode/QRPolynomial.js', 'source_sha256': '76eb786a451ceee003cb4279b7bc559e8a77321dad19ce11825a4d98d470b422', 'source_git_blob': '0c05f38ef324680c1cf53e57159567639177d546', 'start': 2067840, 'end': 2069843, 'region_sha256': '7789cf38de59770fe3377f63d3449aefc08ecea0ba9d3c65e47b037f5918a75c'}, {'path': 'vendor/QRCode/QRMaskPattern.js', 'source_sha256': 'b1f5a99876a31fccbfba89b973e11a4eb295f47b4b00e923814215309c0a725e', 'source_git_blob': 'f6fdeb53730c8b5aecec0182543955f112408df9', 'start': 2069844, 'end': 2070503, 'region_sha256': '7e34365450e4bf9a0a35f7075df7dc7ba336183767469919d695d14201fbf29c'}, {'path': 'vendor/QRCode/QRUtil.js', 'source_sha256': '4191ce852ba66124c4f1fc3bb1a507f667193c0466731339d8c1e66a19aa6bc5', 'source_git_blob': 'e5b7d5b3cc542b49373420ee2f3d3b1d14f2863b', 'start': 2070504, 'end': 2078555, 'region_sha256': '939ca3276f9e9c17531ea48823952ee064d383f64ab1c4e305d131b6622f67f9'}, {'path': 'vendor/QRCode/QRErrorCorrectLevel.js', 'source_sha256': 'd11a145632cea07057084190e86243b3054f30fc77256dc5ea0dc0e0cae54608', 'source_git_blob': '9b4b30099d033344b7b4da95f6774a36d132511c', 'start': 2078556, 'end': 2079113, 'region_sha256': 'ac3dbc0fa27e54a4f201e4541db90dc3630af98683beb04a757fef56b39c7631'}, {'path': 'vendor/QRCode/QRRSBlock.js', 'source_sha256': '78281d6a39b575a1078f1f70e7311e4a3c8b67e15e5468c25521b64d6ff6b931', 'source_git_blob': 'd150af174607997d03d501312c47a6969440a32a', 'start': 2079114, 'end': 2086288, 'region_sha256': 'db806a292b99a099c6ad6012716e4a06f36c3665a1c55f0c4cc1e12ee43b1a1e'}, {'path': 'vendor/QRCode/QRBitBuffer.js', 'source_sha256': '0b5de11b341f5dd92caf3e3a26469f86fa3eb9b3795db6a489e4d53d91ecb67d', 'source_git_blob': 'e2861f68d1b38a93e6630c3403c0769ba41af376', 'start': 2086289, 'end': 2087573, 'region_sha256': '80048565dce106ebc0d7d47d5a1238c43fe4cb5ea34c6eac572c2cb96fea7232'}, {'path': 'vendor/QRCode/index.js', 'source_sha256': '7377be90fc61a40268acf7f30d5bd89c2fca99c57ef5391623de8c151b8da7df', 'source_git_blob': '10eb8eb0a06aa50d0d3c508f886a7728f52e5d98', 'start': 2087574, 'end': 2099556, 'region_sha256': 'cd16b1735b8762154fb6e8eca5a9fb3d3578c12f15e8cbfccc8bcfacc241f6a8'}]
QRCODE_COMMIT = '90f66cf5c6b10bcb4358df96a9580f9eb383307b'
QRCODE_HEADER = '//---------------------------------------------------------------------\n// QRCode for JavaScript\n//\n// Copyright (c) 2009 Kazuhiko Arase\n//\n// URL: http://www.d-project.com/\n//\n// Licensed under the MIT license:\n//   http://www.opensource.org/licenses/mit-license.php\n//\n// The word "QR Code" is registered trademark of \n// DENSO WAVE INCORPORATED\n//   http://www.denso-wave.com/qrcode/faqpatent-e.html\n//\n//---------------------------------------------------------------------\n// Modified to work in node for this project (and some refactoring)\n//---------------------------------------------------------------------\n'
QRCODE_LICENSE = 'MIT License\n\nCopyright (c) 2009 Kazuhiko Arase\n\nPermission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:\n\nThe above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.\n\nTHE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.\n'
LICENSE_URL = 'https://blueoakcouncil.org/license/1.0.0'
BLUEOAK = {'chownr': {'version': '3.0.0', 'license': 'BlueOak-1.0.0', 'sha256': '4300e90fdd91ec7035047473c60f880251a9801bd786302729d4277751d3b948', 'path': 'dist/node_modules/chownr/package.json', 'source': 'https://github.com/isaacs/chownr/blob/8b9800ac5fe4da0b58bffc9c66dd618f0721472d/package.json'}, 'isexe': {'version': '4.0.0', 'license': 'BlueOak-1.0.0', 'sha256': '1889746f84d1e2a524353a2e31464995b808f24457f0dd55d6fd9da638f978ed', 'path': 'dist/node_modules/isexe/package.json', 'source': 'https://github.com/isaacs/isexe/blob/2e7df7dabc4f68e88cf6b32b9029225ba52c6b0c/package.json'}, 'minipass': {'version': '7.1.3', 'license': 'BlueOak-1.0.0', 'sha256': 'b8ab116187d5e9375d494f17437eb862cbee4b329c21f15412c6a5170ec2f55c', 'path': 'dist/node_modules/minipass/package.json', 'source': 'https://github.com/isaacs/minipass/blob/ab4b3b05d0d557ac6bb178f38501b11d0c96454e/package.json'}, 'tar': {'version': '7.5.22', 'license': 'BlueOak-1.0.0', 'sha256': '31b18f4fb83ba7183f54fa9e2cdca610897108c37ba80570c2c6fa7ccf112d97', 'path': 'dist/node_modules/tar/package.json', 'source': 'https://github.com/isaacs/node-tar/blob/2a22bfc5d3a432a606d9da0e2d87ba634aa3b1cb/package.json'}, 'yallist': {'version': '5.0.0', 'license': 'BlueOak-1.0.0', 'sha256': '1b9d47057ce39814531ff93f668823b4fa03e7d23945449c274a1ff6d4cc297f', 'path': 'dist/node_modules/yallist/package.json', 'source': 'https://github.com/isaacs/yallist/blob/a082c2dd872bb3cd8fb92359ed66605064957cd4/package.json'}}
EMPTY_FILES = ('config/user.npmrc', 'config/global.npmrc')
EMPTY_DIRS = ('config', 'home', 'home/appdata', 'home/localappdata', 'home/config',
              'home/data', 'home/pnpm', 'tmp')



def blueoak_notice(projection):
    """Link-form notice for five exact root declarations; no invented archive licence."""
    rows = []
    for name, binding in BLUEOAK.items():
        raw = read_bounded(Path(projection) / 'payload' / binding['path'], 1048576)
        require(digest(raw) == binding['sha256'], 'BlueOak manifest changed: ' + name)
        value = decode(raw)
        require(value.get('name') == name and value.get('version') == binding['version']
                and value.get('license') == 'BlueOak-1.0.0', 'BlueOak declaration changed: ' + name)
        rows.append('| ' + name + '@' + binding['version'] + ' | ' + binding['path']
                    + ' | ' + binding['sha256'] + ' | ' + binding['source'] + ' |')
    text = ('# Supplemental Blue Oak license notice\n\n'
            'License: Blue Oak Model License 1.0.0 (BlueOak-1.0.0).\n'
            'License link: ' + LICENSE_URL + '\n\n'
            'This link-form notice accompanies the following five exact package-root declarations '
            'in the selected pnpm 11.25.0 candidate. Keep this notice with any copy of these works.\n\n'
            '| Package | Manifest in payload | Manifest SHA-256 | Fixed author declaration |\n'
            '|---|---|---|---|\n' + '\n'.join(rows) + '\n\n'
            'This supplemental file is generated by local tooling; it is not an original file '
            'from the pnpm archive. chownr has a fixed author declaration but no named source '
            'license file in the retained tree. The yallist source license scopes packages under '
            'src to their own licenses; this notice does not relicense any nested work.\n\n'
            'This notice covers the listed declarations only. All original notices remain in '
            'the unchanged payload. Other pnpm licenses, source obligations, publishing provenance '
            'and runtime/redistribution admission remain separate and unresolved.\n')
    return text.encode('utf-8')


def qrcode_notice(projection):
    """Deliver the retained vendor attribution and full MIT permission notice."""
    bundle = read_bounded(Path(projection) / 'payload' / QRCODE_BUNDLE, 20000000)
    require(digest(bundle) == QRCODE_BUNDLE_SHA256, 'QRCode bundle changed')
    require(len(QRCODE_MODULES) == 10, 'QRCode module inventory changed')
    rows = []
    for item in QRCODE_MODULES:
        require(0 <= item['start'] < item['end'] <= len(bundle), 'QRCode module range invalid')
        require(digest(bundle[item['start']:item['end']]) == item['region_sha256'], 'QRCode module changed')
        rows.append('| ' + item['path'] + ' | ' + item['source_sha256'] + ' | '
                    + item['region_sha256'] + ' |')
    text = ('# QRCode vendored attribution and MIT permission notice\n\n'
            'Component: qrcode-terminal 0.12.0, vendor/QRCode.\n'
            'Selected pnpm 11.25.0 bundle SHA-256: ' + QRCODE_BUNDLE_SHA256 + '\n'
            'Fixed source: https://github.com/gtanner/qrcode-terminal/tree/' + QRCODE_COMMIT + '/vendor/QRCode\n'
            'Permission text authority: https://opensource.org/license/mit\n\n'
            '## Retained author and modification notice\n\n' + QRCODE_HEADER + '\n'
            '## MIT permission notice\n\n' + QRCODE_LICENSE + '\n'
            '## Exact source and bundle-region inventory\n\n'
            '| Vendor path | Fixed source SHA-256 | Selected bundle region SHA-256 |\n'
            '|---|---|---|\n' + '\n'.join(rows) + '\n\n'
            'This is a locally assembled supplement. The source header is retained verbatim; '
            'the complete MIT permission text is supplied with the copyright identity from that header. '
            'It is not an original pnpm archive file. It does not replace the containing package\'s '
            'Apache license or any other notice. The path/hash inventory is not proof of a complete '
            'reproducible build or equivalent module bindings. Keep this notice with copies of the work. '
            'This planner does not distribute the payload or admit runtime/redistribution.\n')
    return text.encode('utf-8')


def semver_notice(projection):
    """Retain the exact dual-license file; select its full MIT option for this supplement."""
    bundle = read_bounded(Path(projection) / 'payload' / QRCODE_BUNDLE, 20000000)
    require(digest(bundle) == QRCODE_BUNDLE_SHA256, 'semver-utils bundle changed')
    item = SEMVER_REGION
    require(0 <= item['start'] < item['end'] <= len(bundle), 'semver-utils region invalid')
    require(digest(bundle[item['start']:item['end']]) == item['region_sha256'], 'semver-utils region changed')
    require(digest(SEMVER_LICENSE.encode('utf-8')) == SEMVER_LICENSE_SHA256, 'semver-utils original license changed')
    text = ('# semver-utils 1.1.4 original license notice\n\n'
            'Original npm artifact: https://registry.npmjs.org/semver-utils/-/semver-utils-1.1.4.tgz\n'
            'Artifact SHA-256: fa6980458971864f25f2afc7d7f6632b015d8b767296d48421e6ff13221bd2b3\n'
            'Original package/LICENSE SHA-256: ' + SEMVER_LICENSE_SHA256 + '\n'
            'Original package/semver-utils.js SHA-256: ' + item['source_member_sha256'] + '\n'
            'Selected pnpm bundle SHA-256: ' + QRCODE_BUNDLE_SHA256 + '\n'
            'Selected bundle region SHA-256: ' + item['region_sha256'] + '\n\n'
            'The original LICENSE explicitly permits either MIT or Apache-2.0. This supplement '
            'selects the complete MIT option and retains the entire original file below. '
            'The package manifest and registry declaration APACHEv2 remain unchanged as historical '
            'metadata; this is a locally assembled supplement, not a replacement pnpm archive file.\n\n'
            '## Verbatim original license\n\n' + SEMVER_LICENSE + '\n'
            '## Evidence scope\n\n'
            'The artifact was acquired in quarantine against exact SHA512 registry integrity. '
            'Its registry signature cryptographically verifies with a key expired on 2025-01-29; '
            'no current valid signature or trusted signing time is claimed. '
            'The source member and bundle region are separately fixed observations, not proof of '
            'complete source/build equivalence. Keep the copyright and permission notice with copies '
            'of the work. Other pnpm obligations and runtime/redistribution admission remain open.\n')
    return text.encode('utf-8')


def retained_notices(projection):
    """Deliver the reviewed evidence collection without converting it into license admission."""
    raw = read_bounded(Path(__file__).with_name('retained-notice-catalog.json'), 524288)
    require(digest(raw) == CATALOG_SHA256, 'retained notice catalog changed')
    catalog = decode(raw)
    require(set(catalog) == {'schema', 'status', 'release_ready', 'parent_manifest_sha256', 'entries'},
            'retained notice catalog fields changed')
    require(catalog['schema'] == 'elite-pnpm-retained-notice-catalog/v1'
            and catalog['status'] == 'RETAINED_TEXTS_ONLY' and catalog['release_ready'] is False,
            'retained notice scope inflated')
    require(catalog['parent_manifest_sha256'] == '27d666e86d769c368439b1cc20e0fa72a42f9366973216b42ced5434f003fa9d',
            'retained notice parent changed')
    require(isinstance(catalog['entries'], list) and len(catalog['entries']) == 47,
            'retained notice inventory changed')
    seen = set(); total = 0; originals = 0
    out = [b'# Retained pnpm notice and license-text evidence\n\n'
           b'47 exact text copies: 22 original selected-payload notices and 25 retained research texts.\n'
           b'This collection preserves evidence; it is not a complete recursive SBOM, license clearance, '
           b'source offer or source/relinking fulfillment. Runtime and redistribution admission remain open.\n'
           b'Texts can repeat across versions; 47 is not a count of uniquely cleared dependencies. '
           b'Original pnpm payload, BlueOak, QRCode and semver-utils supplements remain separate and unchanged.\n\n'
           b'Authority and qualification records: PNPM_NOTICE_COVERAGE_V338.md, '
           b'PNPM_NESTED_LICENSE_SOURCES_V339.md, PNPM_YARN_UNDICI_NOTICES_V340.md, '
           b'SEMVER_ORIGINAL_LICENSE_DELIVERY_V396.md. Tagged source and published-byte identity '
           b'must not be conflated. Follow the exact scope attached to each retained text.\n\n']
    for entry in catalog['entries']:
        require(set(entry) == {'path', 'source_locator', 'scope', 'origin', 'payload_path', 'bytes', 'sha256', 'content_base64'},
                'retained notice entry fields changed')
        path = safe_name(entry['path'])
        require(path.casefold() not in seen, 'duplicate retained notice path'); seen.add(path.casefold())
        require(type(entry['bytes']) is int and 0 < entry['bytes'] <= 65536, 'retained notice byte budget')
        require(isinstance(entry['content_base64'], str), 'retained notice encoding required')
        try:
            content = base64.b64decode(entry['content_base64'], validate=True)
        except (ValueError, binascii.Error) as exc:
            raise SelectionError('retained notice encoding invalid') from exc
        require(base64.b64encode(content).decode('ascii') == entry['content_base64'], 'noncanonical notice encoding')
        require(len(content) == entry['bytes'] and digest(content) == entry['sha256'], 'retained notice text changed')
        total += len(content); require(total <= 200000, 'retained collection byte budget')
        require(all(isinstance(entry[k], str) and entry[k] and '\n' not in entry[k] and '\r' not in entry[k]
                    for k in ('source_locator', 'scope')), 'retained notice provenance required')
        if entry['origin'] == 'ORIGINAL_SELECTED_PAYLOAD':
            require(path.startswith('original/') and entry['payload_path'] == path[len('original/'):],
                    'original notice path mismatch')
            actual = read_bounded(Path(projection) / 'payload' / safe_name(entry['payload_path']), 65536)
            require(actual == content, 'original payload notice changed'); originals += 1
        else:
            require(entry['origin'] == 'RETAINED_RESEARCH_EVIDENCE' and entry['payload_path'] is None
                    and not path.startswith('original/'), 'retained notice origin inflated')
        header = ('## ' + path + '\n\nOrigin: ' + entry['origin'] + '\nScope: ' + entry['scope']
                  + '\nSource evidence locator: ' + entry['source_locator'] + '\nSHA-256: '
                  + entry['sha256'] + '\nBytes: ' + str(entry['bytes']) + '\n\nBEGIN VERBATIM TEXT\n\n').encode('utf-8')
        out.extend([header, content, b'\n\nEND VERBATIM TEXT\n\n'])
    require(originals == 22 and total == 151253, 'retained notice coverage changed')
    result = b''.join(out); require(len(result) <= 262144, 'retained notice output budget')
    return result


def next_path_notice(projection):
    """Supply original covered source, manifest, license and bounded comparison qualifications."""
    raw = read_bounded(Path(__file__).with_name('next-path-source-evidence.json'), 65536)
    require(digest(raw) == NEXT_PATH_DATA_SHA256, 'next-path source evidence changed')
    data = decode(raw)
    require(data['schema'] == 'elite-next-path-source-evidence/v1'
            and data['package'] == 'next-path' and data['version'] == '1.0.0'
            and data['license'] == 'MPL-2.0' and data['runtime_admitted'] is False
            and data['status'] == 'SOURCE_AND_LICENSE_DELIVERY_ONLY', 'next-path evidence scope inflated')
    require([entry['path'] for entry in data['files']] == ['index.js', 'package.json', 'LICENSE'],
            'next-path source inventory changed')
    bundle = read_bounded(Path(projection) / 'payload' / QRCODE_BUNDLE, 20000000)
    require(digest(bundle) == QRCODE_BUNDLE_SHA256, 'next-path selected bundle changed')
    region = NEXT_PATH_REGION
    require(0 <= region['start'] < region['end'] <= len(bundle)
            and digest(bundle[region['start']:region['end']]) == region['sha256'],
            'next-path bundle region changed')
    proof = data['correspondence']
    require(proof['statements_compared'] == 4 and len(proof['negative_probes']) == 10
            and proof['source_code_executed'] is False and proof['runtime_equivalence_claimed'] is False
            and proof['whole_pnpm_build_proven'] is False and proof['npm_tarball_byte_identity_proven'] is False,
            'next-path comparison scope inflated')
    out = [('# next-path 1.0.0 — MPL source and license\n\n'
            'Covered component license: Mozilla Public License 2.0 (MPL-2.0).\n'
            'The original source files below are made available under that license; the full license is included. '
            'No additional restriction is imposed here on the rights granted for that covered source.\n'
            'Official fixed source: https://github.com/sholladay/next-path/tree/' + data['commit'] + '\n'
            'Selected pnpm bundle SHA-256: ' + QRCODE_BUNDLE_SHA256 + '\n'
            'Selected module region SHA-256: ' + region['sha256'] + '\n\n'
            '## Source and transformation scope\n\n'
            'The fixed package manifest publishes index.js as its main and sole files entry. '
            'All four module statements were compared structurally after only the explicitly listed adapters. '
            'Ten changed-source controls were rejected. This is source correspondence under those adapters, '
            'not runtime equivalence, an npm tarball byte match or a reproducible whole-pnpm build.\n'
            'Generated packaging changes: one __commonJS wrapper, two named top-level const declarations emitted as var, '
            'and the following identifier/loader mappings (emitted name to original name): '
            + ', '.join(k + ' -> ' + v for k, v in proof['identifier_renames'].items()) + '.\n'
            '__require(path) -> require(path) is a conditional loader adapter; actual host resolution remains unproven. '
            'The function body, property names, literals, argument order and export were preserved by the comparison. '
            'No additional business or algorithm change is claimed. Other pnpm components keep their own licenses.\n\n'
            '## Original files — complete, unchanged bytes\n\n').encode('utf-8')]
    for entry in data['files']:
        require(type(entry['bytes']) is int and 0 < entry['bytes'] <= 32768, 'next-path source byte budget')
        try:
            content = base64.b64decode(entry['content_base64'], validate=True)
        except (ValueError, binascii.Error) as exc:
            raise SelectionError('next-path source encoding invalid') from exc
        require(len(content) == entry['bytes'] and digest(content) == entry['sha256'], 'next-path original source changed')
        out.extend([('### ' + entry['path'] + '\n\nSource: ' + entry['source_url'] + '\nGit blob: '
                     + entry['git_blob'] + '\nSHA-256: ' + entry['sha256'] + '\nBytes: '
                     + str(entry['bytes']) + '\n\nBEGIN ORIGINAL FILE\n\n').encode('utf-8'),
                    content, b'\n\nEND ORIGINAL FILE\n\n'])
    out.append(b'This local delivery does not authorize redistribution of the entire pnpm bundle, '
               b'waive other source/relinking obligations, or admit any runtime. '
               b'No supplied source file is installed or executed by the planner.\n')
    result = b''.join(out); require(len(result) <= 65536, 'next-path notice output budget')
    return result

def paths(project, projection, node, store, cache, target):
    result = [no_reparse(x) for x in (project, projection, node, store, cache, target)]
    project, projection, node, store, cache, target = result
    require(all(p.is_dir() for p in (project, projection, store, cache)), 'input directory missing')
    require(node.is_file(), 'Node executable missing')
    for i, a in enumerate(result):
        for b in result[i + 1:]:
            require(a != b and a not in b.parents and b not in a.parents, 'overlapping inputs/output rejected')
    return result


def consumer_files(consumer, project):
    require(consumer in PROFILES, 'unsupported consumer')
    policy = PROFILES[consumer]
    require(not os.path.lexists(project / 'node_modules'), 'fresh consumer without node_modules required')
    for name, expected in policy.items():
        require(digest(read_bounded(project / name, 2097152)) == expected, 'consumer input changed: ' + name)
    for ancestor in (project, *project.parents):
        for name in ('.npmrc', '.pnpmfile.cjs', 'pnpmfile.cjs', '.pnpmfile.mjs', 'pnpmfile.mjs'):
            require(not os.path.lexists(ancestor / name), 'ambient project configuration rejected')
    if 'pnpm-workspace.yaml' not in policy:
        require(not os.path.lexists(project / 'pnpm-workspace.yaml'), 'unsupported workspace configuration')
    require(not os.path.lexists(project / 'pnpm-workspace.yml'), 'unsupported workspace extension')
    return dict(policy)


def recipe(consumer, project, projection, node, store, cache, target, acquisition_receipt, acquisition_sha256):
    require(os.name == 'nt', 'routing supported only on Windows')
    project, projection, node, store, cache, target = paths(project, projection, node, store, cache, target)
    acquisition_receipt = no_reparse(acquisition_receipt)
    require(target != acquisition_receipt and target not in acquisition_receipt.parents, 'receipt/output overlap')
    raw = read_bounded(acquisition_receipt, 1048576)
    require(digest(raw) == acquisition_sha256, 'acquisition identity changed')
    policy = load_policy()
    acquisition_binding(raw, policy)
    selection = verify(projection, policy)
    require(selection['acquisition_receipt_sha256'] == acquisition_sha256, 'projection acquisition mismatch')
    require(digest(read_bounded(node, 150000000)) == NODE_SHA256, 'Node identity changed')
    files = consumer_files(consumer, project)
    system = no_reparse(os.environ.get('SystemRoot', ''))
    require(system.is_dir() and system.name.casefold() == 'windows', 'trusted Windows SystemRoot required')
    environment = {
        'SystemRoot': str(system), 'WINDIR': str(system),
        'PATH': str(node.parent) + os.pathsep + str(system / 'System32'),
        'TEMP': str(target / 'tmp'), 'TMP': str(target / 'tmp'),
        'HOME': str(target / 'home'), 'USERPROFILE': str(target / 'home'),
        'APPDATA': str(target / 'home/appdata'), 'LOCALAPPDATA': str(target / 'home/localappdata'),
        'XDG_CONFIG_HOME': str(target / 'home/config'), 'XDG_DATA_HOME': str(target / 'home/data'),
        'PNPM_HOME': str(target / 'home/pnpm'), 'CI': 'true', 'NEXT_TELEMETRY_DISABLED': '1',
    }
    argv = [str(node), str(projection / 'payload/bin/pnpm.mjs'), 'install',
            '--offline', '--frozen-lockfile', '--ignore-scripts', '--ignore-pnpmfile',
            '--package-import-method=copy', '--verify-store-integrity',
            '--store-dir=' + str(store), '--config.cache-dir=' + str(cache), '--config.userconfig=' + str(target / EMPTY_FILES[0]),
            '--config.globalconfig=' + str(target / EMPTY_FILES[1])]
    if 'pnpm-workspace.yaml' not in files:
        argv.append('--ignore-workspace')
    return {'schema': 'elite-pnpm-install-plan/v6', 'consumer': consumer,
            'inputs': {'project': str(project), 'projection': str(projection), 'node': str(node),
                       'store': str(store), 'cache': str(cache), 'acquisition_receipt': str(acquisition_receipt),
                       'acquisition_receipt_sha256': acquisition_sha256},
            'consumer_files': files, 'node_sha256': NODE_SHA256,
            'next_path_source': {'path': NEXT_PATH_NOTICE, 'sha256': digest(next_path_notice(projection)),
                                 'data_sha256': NEXT_PATH_DATA_SHA256, 'scope': 'SOURCE_AND_LICENSE_DELIVERY_ONLY'},
            'retained_notices': {'path': RETAINED_NOTICE, 'sha256': digest(retained_notices(projection)),
                                 'catalog_sha256': CATALOG_SHA256, 'copies': 47, 'scope': 'RETAINED_TEXTS_ONLY'},
            'semver_notice': {'path': SEMVER_NOTICE, 'sha256': digest(semver_notice(projection)),
                              'scope': 'ORIGINAL_SEMVER_UTILS_LICENSE_MIT_OPTION'},
            'qrcode_notice': {'path': QRCODE_NOTICE, 'sha256': digest(qrcode_notice(projection)),
                              'scope': 'QRCODE_VENDOR_ATTRIBUTION_AND_MIT_PERMISSION'},
            'supplemental_notice': {'path': NOTICE, 'sha256': digest(blueoak_notice(projection)),
                                    'license_url': LICENSE_URL, 'scope': 'FIVE_PINNED_ROOT_DECLARATIONS'},
            'selection_receipt_sha256': digest(read_bounded(projection / 'selection-receipt.json', 4096)),
            'cwd': str(project), 'argv': argv, 'environment': environment,
            'environment_mode': 'REPLACE_NOT_MERGE', 'shell': False,
            'runtime_admitted': False, 'redistribution_admitted': False, 'executed': False,
            'conditions': ['Verify this plan and its externally retained SHA256 immediately before use.',
                           'Separate runtime/license admission required before execution.',
                           'Keep NEXT-PATH-MPL-SOURCE.md, PNPM-RETAINED-NOTICES.md, BLUEOAK-NOTICE.md, QRCODE-NOTICE.md and SEMVER-UTILS-NOTICE.md with copies of their listed works; this plan does not distribute the payload.',
                           'Offline and ignore-scripts flags are not an operating-system sandbox.',
                           'Trusted local inputs must remain immutable through any later use.',
                           'No scripts, remote resolution, cloning, global installs or arbitrary CLI arguments in this recipe.']}


def verify_plan(target, expected_sha256):
    target = no_reparse(target)
    raw = read_bounded(target / RECIPE, 65536)
    require(digest(raw) == expected_sha256, 'install plan identity changed')
    value = decode(raw)
    inp = value['inputs']
    expected = recipe(value['consumer'], inp['project'], inp['projection'], inp['node'], inp['store'], inp['cache'],
                      target, inp['acquisition_receipt'], inp['acquisition_receipt_sha256'])
    require(raw == json_bytes(expected), 'install plan changed')
    found_files, found_dirs = set(), set()
    for root, dirs, files in os.walk(target, followlinks=False):
        for name in dirs + files:
            p = no_reparse(Path(root) / name)
            rel = p.relative_to(target).as_posix()
            if p.is_dir():
                require(rel in EMPTY_DIRS, 'unexpected plan directory')
                found_dirs.add(rel)
            else:
                require(rel in (*EMPTY_FILES, RECIPE, NOTICE, QRCODE_NOTICE, SEMVER_NOTICE, RETAINED_NOTICE, NEXT_PATH_NOTICE), 'unexpected plan file')
                found_files.add(rel)
                if rel == NEXT_PATH_NOTICE:
                    require(digest(read_bounded(p, 65536)) == expected['next_path_source']['sha256'], 'next-path source delivery changed')
                if rel == RETAINED_NOTICE:
                    require(digest(read_bounded(p, 262144)) == expected['retained_notices']['sha256'], 'retained notice delivery changed')
                if rel == SEMVER_NOTICE:
                    require(digest(read_bounded(p, 16384)) == expected['semver_notice']['sha256'], 'semver-utils notice changed')
                if rel == QRCODE_NOTICE:
                    require(digest(read_bounded(p, 16384)) == expected['qrcode_notice']['sha256'], 'QRCode notice changed')
                if rel == NOTICE:
                    require(digest(read_bounded(p, 16384)) == expected['supplemental_notice']['sha256'], 'supplemental notice changed')
                if rel in EMPTY_FILES:
                    require(read_bounded(p, 0) == b'', 'plan configuration is not empty')
    require(found_files == {*EMPTY_FILES, RECIPE, NOTICE, QRCODE_NOTICE, SEMVER_NOTICE, RETAINED_NOTICE, NEXT_PATH_NOTICE} and found_dirs == set(EMPTY_DIRS), 'plan files/directories missing')
    return value


def create_plan(consumer, project, projection, node, store, cache, target, acquisition_receipt, acquisition_sha256):
    target = no_reparse(target)
    require(target.parent.is_dir() and not os.path.lexists(target), 'target must be absent with existing parent')
    value = recipe(consumer, project, projection, node, store, cache, target, acquisition_receipt, acquisition_sha256)
    raw = json_bytes(value)
    stage = Path(tempfile.mkdtemp(prefix='.pnpm-plan-', dir=target.parent))
    try:
        for name in EMPTY_DIRS:
            (stage / name).mkdir(exist_ok=True)
        for name in EMPTY_FILES:
            (stage / name).write_bytes(b'')
        (stage / NOTICE).write_bytes(blueoak_notice(projection))
        (stage / QRCODE_NOTICE).write_bytes(qrcode_notice(projection))
        (stage / SEMVER_NOTICE).write_bytes(semver_notice(projection))
        (stage / RETAINED_NOTICE).write_bytes(retained_notices(projection))
        (stage / NEXT_PATH_NOTICE).write_bytes(next_path_notice(projection))
        (stage / RECIPE).write_bytes(raw)
        # Recipe paths refer to the final target. Publication never replaces a competitor.
        no_reparse(target.parent)
        os.rename(stage, target)
    finally:
        if stage.exists():
            require(stage.parent == target.parent and stage.name.startswith('.pnpm-plan-'), 'cleanup scope rejected')
            shutil.rmtree(stage)
    verify_plan(target, digest(raw))
    return digest(raw)


def main(argv=None):
    parser = argparse.ArgumentParser(description=__doc__)
    sub = parser.add_subparsers(dest='command', required=True)
    make = sub.add_parser('prepare')
    make.add_argument('--consumer', choices=sorted(PROFILES), required=True)
    for name in ('project', 'projection', 'node', 'store', 'cache', 'target', 'acquisition-receipt', 'acquisition-sha256'):
        make.add_argument('--' + name, required=True)
    check = sub.add_parser('verify')
    check.add_argument('--target', required=True)
    check.add_argument('--sha256', required=True)
    args = parser.parse_args(argv)
    try:
        if args.command == 'prepare':
            sha = create_plan(args.consumer, args.project, args.projection, args.node, args.store, args.cache,
                              args.target, args.acquisition_receipt, args.acquisition_sha256)
        else:
            verify_plan(args.target, args.sha256)
            sha = args.sha256
        print('PNPM_INSTALL_PLAN_PASS sha256=' + sha + ' executed=false runtime_admitted=false')
        return 0
    except (SelectionError, OSError, ValueError, UnicodeError, KeyError, TypeError, RecursionError) as exc:
        print('PNPM_INSTALL_PLAN_BLOCKED: ' + str(exc), file=sys.stderr)
        return 2

if __name__ == '__main__':
    raise SystemExit(main())
````

### FILE: `pnpm_artifact_selection/test_routing.py`
```yaml
block_id: "PNPM-ARTIFACT-SELECTION-GATE:test_routing.py:v395"
operation: CREATE
provenance: AUTHORED
source: "local, bounded projection and consumer planning governed by canonical acquisition and execution contracts"
license: "LicenseRef-Workspace-Owner"
sha256: "7d35af69567c2293349e226bc0dacde8aa04769a11817288375d73097c5f2680"
variables: []
secrets_allowed: false
```
````python
import concurrent.futures
import base64
import hashlib
import json
import os
from pathlib import Path
import subprocess
import tempfile
import unittest
from unittest.mock import patch
import plan_install as gate

@unittest.skipUnless(os.name == 'nt', 'Windows exclusive-publication contract')
class RoutingTests(unittest.TestCase):
    def setUp(self):
        self.tmp = tempfile.TemporaryDirectory(prefix='elite-route-test-')
        self.addCleanup(self.tmp.cleanup)
        self.root = Path(self.tmp.name)
        self.project, self.projection, self.store, self.cache = [self.root / x for x in ('project', 'projection', 'store', 'cache')]
        for p in (self.project, self.projection, self.store, self.cache): p.mkdir()
        self.node = self.root / 'node.exe'; self.node.write_bytes(b'fixture never executed')
        self.acq = self.root / 'acquisition.json'; self.acq.write_bytes(b'{}\n')
        self.acq_sha = gate.digest(self.acq.read_bytes())
        (self.projection / 'selection-receipt.json').write_bytes(b'fixture receipt')
        self.filehash = {}
        for name in ('package.json', 'pnpm-lock.yaml'):
            (self.project / name).write_bytes((name + '\n').encode())
            self.filehash[name] = gate.digest((self.project / name).read_bytes())
        self.output = self.root / 'plan'
        self.catalog = json.loads(Path(gate.__file__).with_name('retained-notice-catalog.json').read_bytes())
        for entry in self.catalog['entries']:
            if entry['origin'] == 'ORIGINAL_SELECTED_PAYLOAD':
                dest = self.projection / 'payload' / entry['payload_path']
                dest.parent.mkdir(parents=True, exist_ok=True)
                dest.write_bytes(base64.b64decode(entry['content_base64']))
        self.notice_bindings = {}
        for name, original in gate.BLUEOAK.items():
            binding = dict(original)
            raw = gate.json_bytes({'name': name, 'version': binding['version'], 'license': 'BlueOak-1.0.0'})
            dest = self.projection / 'payload' / binding['path']; dest.parent.mkdir(parents=True, exist_ok=True); dest.write_bytes(raw)
            binding['sha256'] = gate.digest(raw)
            self.notice_bindings[name] = binding
        notice_patch = patch.dict(gate.BLUEOAK, self.notice_bindings, clear=True)
        notice_patch.start(); self.addCleanup(notice_patch.stop)
        bundle = b''; modules = []
        for index, original in enumerate(gate.QRCODE_MODULES):
            raw = ('synthetic module ' + str(index) + '\n').encode()
            item = dict(original); item.update(start=len(bundle), end=len(bundle)+len(raw), region_sha256=gate.digest(raw))
            modules.append(item); bundle += raw
        self.qrcode_bundle = self.projection / 'payload' / gate.QRCODE_BUNDLE
        self.qrcode_bundle.parent.mkdir(parents=True, exist_ok=True); self.qrcode_bundle.write_bytes(bundle)
        for patched in [patch.object(gate, 'QRCODE_BUNDLE_SHA256', gate.digest(bundle)),
                        patch.object(gate, 'QRCODE_MODULES', modules),
                        patch.object(gate, 'NEXT_PATH_REGION', {'start': 0, 'end': len(bundle), 'sha256': gate.digest(bundle)}),
                        patch.object(gate, 'SEMVER_REGION', dict(gate.SEMVER_REGION, start=0, end=len(bundle), region_sha256=gate.digest(bundle)))]:
            patched.start(); self.addCleanup(patched.stop)
        for p in [patch.dict(gate.PROFILES, {'fixture': self.filehash}, clear=True),
                  patch.object(gate, 'NODE_SHA256', gate.digest(self.node.read_bytes())),
                  patch.object(gate, 'load_policy', return_value={}),
                  patch.object(gate, 'acquisition_binding', return_value=self.acq_sha),
                  patch.object(gate, 'verify', return_value={'acquisition_receipt_sha256': self.acq_sha})]:
            p.start(); self.addCleanup(p.stop)

    def create(self, **changes):
        args = dict(consumer='fixture', project=self.project, projection=self.projection, node=self.node,
                    store=self.store, cache=self.cache, target=self.output, acquisition_receipt=self.acq, acquisition_sha256=self.acq_sha)
        args.update(changes)
        return gate.create_plan(**args)

    def rejected(self, **changes):
        with self.assertRaises((gate.SelectionError, OSError, ValueError)):
            self.create(**changes)
        self.assertFalse(self.output.exists())
        self.assertEqual([], list(self.root.glob('.pnpm-plan-*')))

    def test_recipe_is_bound_and_never_admitted(self):
        sha = self.create(); value = gate.verify_plan(self.output, sha)
        self.assertFalse(value['executed']); self.assertFalse(value['runtime_admitted'])
        self.assertFalse(value['redistribution_admitted']); self.assertFalse(value['shell'])
        for arg in ['--offline','--frozen-lockfile','--ignore-scripts','--ignore-pnpmfile','--ignore-workspace','--package-import-method=copy']:
            self.assertIn(arg, value['argv'])
        self.assertEqual(str(self.project), value['cwd'])
        self.assertTrue(any(x.startswith('--config.userconfig=') for x in value['argv']))
        self.assertTrue(any(x.startswith('--config.globalconfig=') for x in value['argv']))
        self.assertFalse(any(x.startswith(('--userconfig=', '--globalconfig=')) for x in value['argv']))

    def test_host_secrets_and_options_are_not_inherited(self):
        with patch.dict(os.environ, {'NODE_OPTIONS':'--require=evil.cjs','NPM_CONFIG_REGISTRY':'secret-value','PNPM_HOME':'secret-home','AWS_SECRET_ACCESS_KEY':'secret-key'}):
            sha=self.create(); value=gate.verify_plan(self.output,sha)
        self.assertEqual('REPLACE_NOT_MERGE',value['environment_mode'])
        self.assertNotIn('NODE_OPTIONS',value['environment']); self.assertNotIn('NPM_CONFIG_REGISTRY',value['environment'])
        self.assertNotIn('secret-',json.dumps(value))

    def test_unknown_consumer(self): self.rejected(consumer='arbitrary')
    def test_manifest_drift(self):
        (self.project/'package.json').write_text('changed'); self.rejected()
    def test_lock_drift(self):
        (self.project/'pnpm-lock.yaml').write_text('changed'); self.rejected()
    def test_local_npmrc(self):
        (self.project/'.npmrc').write_text('registry=changed'); self.rejected()
    def test_ancestor_hook(self):
        (self.root/'.pnpmfile.cjs').write_text('throw new Error()'); self.rejected()
    def test_workspace_not_in_profile(self):
        (self.project/'pnpm-workspace.yaml').write_text('packages: [evil]'); self.rejected()
    def test_wrong_node(self):
        self.node.write_bytes(b'changed'); self.rejected()
    def test_wrong_acquisition_sha(self): self.rejected(acquisition_sha256='0'*64)
    def test_projection_from_other_receipt(self):
        with patch.object(gate,'verify',return_value={'acquisition_receipt_sha256':'1'*64}): self.rejected()
    def test_existing_install_is_rejected(self):
        (self.project/'node_modules').mkdir(); self.rejected()
    def test_known_workspace_policy_is_preserved(self):
        p=self.project/'pnpm-workspace.yaml';p.write_bytes(b'known-workspace')
        self.filehash[p.name]=gate.digest(p.read_bytes())
        sha=self.create(); value=gate.verify_plan(self.output,sha)
        self.assertNotIn('--ignore-workspace',value['argv'])
        self.assertIn(p.name,value['consumer_files'])
    def test_missing_metadata_cache_is_rejected(self):
        self.rejected(cache=self.root/'absent-cache')
    def test_overlap(self): self.rejected(target=self.project/'plan')
    def test_existing_destination_is_preserved(self):
        self.output.mkdir(); (self.output/'marker').write_bytes(b'old')
        with self.assertRaises(gate.SelectionError): self.create()
        self.assertEqual(b'old',(self.output/'marker').read_bytes())
    def test_late_competitor_is_preserved(self):
        original=gate.os.rename
        def compete(a,b):
            Path(b).mkdir(); (Path(b)/'marker').write_bytes(b'winner'); original(a,b)
        with patch.object(gate.os,'rename',side_effect=compete):
            with self.assertRaises(OSError): self.create()
        self.assertEqual(b'winner',(self.output/'marker').read_bytes()); self.assertEqual([],list(self.root.glob('.pnpm-plan-*')))
    def test_normal_publication_failure_cleans_owned_stage(self):
        with patch.object(gate.os,'rename',side_effect=PermissionError('injected')): self.rejected()
    def test_two_writers_have_one_winner(self):
        def attempt(_):
            try: self.create(); return 'win'
            except (gate.SelectionError,OSError): return 'rejected'
        with concurrent.futures.ThreadPoolExecutor(max_workers=2) as pool: results=list(pool.map(attempt,range(2)))
        self.assertEqual(['rejected','win'],sorted(results)); self.assertEqual([],list(self.root.glob('.pnpm-plan-*')))
    def test_tampered_plan_hash(self):
        sha=self.create(); p=self.output/gate.RECIPE; value=json.loads(p.read_bytes()); value['runtime_admitted']=True; p.write_bytes(gate.json_bytes(value))
        with self.assertRaises(gate.SelectionError): gate.verify_plan(self.output,sha)
        with self.assertRaises(gate.SelectionError): gate.verify_plan(self.output,gate.digest(p.read_bytes()))
    def test_config_change_is_rejected(self):
        sha=self.create(); (self.output/gate.EMPTY_FILES[0]).write_bytes(b'unsafe=true')
        with self.assertRaises(gate.SelectionError): gate.verify_plan(self.output,sha)
    def test_extra_file_is_rejected(self):
        sha=self.create(); (self.output/'home/evil.cjs').write_bytes(b'unsafe')
        with self.assertRaises(gate.SelectionError): gate.verify_plan(self.output,sha)
    def test_missing_directory_is_rejected(self):
        sha=self.create(); (self.output/'home/data').rmdir()
        with self.assertRaises(gate.SelectionError): gate.verify_plan(self.output,sha)
    def test_drift_after_plan_is_rejected(self):
        sha=self.create(); (self.project/'pnpm-lock.yaml').write_bytes(b'drift')
        with self.assertRaises(gate.SelectionError): gate.verify_plan(self.output,sha)
    def test_junction_input_rejected(self):
        link=self.root/'junction'; r=subprocess.run(['cmd','/d','/c','mklink','/J',str(link),str(self.project)],capture_output=True)
        self.assertEqual(0,r.returncode,r.stderr.decode(errors='replace'))
        try:self.rejected(project=link)
        finally:link.rmdir()
    def test_invalid_json_cli_is_actionable(self):
        sha=self.create(); p=self.output/gate.RECIPE;p.write_bytes(b'{')
        self.assertEqual(2,gate.main(['verify','--target',str(self.output),'--sha256',gate.digest(p.read_bytes())]))


    def test_notice_delivered_for_all_five_exact_declarations(self):
        sha = self.create(); value = gate.verify_plan(self.output, sha)
        raw = (self.output / gate.NOTICE).read_bytes()
        self.assertIn(b'https://blueoakcouncil.org/license/1.0.0', raw)
        for name, row in self.notice_bindings.items():
            self.assertIn((name + '@' + row['version']).encode(), raw)
            self.assertIn(row['sha256'].encode(), raw)
        self.assertEqual(hashlib.sha256(raw).hexdigest(), value['supplemental_notice']['sha256'])
        self.assertEqual('elite-pnpm-install-plan/v6', value['schema'])

    def test_missing_notice_rejected(self):
        sha = self.create(); (self.output / gate.NOTICE).unlink()
        with self.assertRaises(gate.SelectionError): gate.verify_plan(self.output, sha)

    def test_changed_notice_rejected(self):
        sha = self.create(); p = self.output / gate.NOTICE
        p.write_bytes(p.read_bytes().replace(b'https://blueoakcouncil.org/license/1.0.0', b'https://invalid.example/license'))
        with self.assertRaises(gate.SelectionError): gate.verify_plan(self.output, sha)

    def test_forged_notice_and_rehashed_receipt_rejected(self):
        self.create(); p = self.output / gate.NOTICE; p.write_bytes(b'false license')
        receipt = self.output / gate.RECIPE; value = json.loads(receipt.read_bytes())
        value['supplemental_notice']['sha256'] = gate.digest(p.read_bytes())
        receipt.write_bytes(gate.json_bytes(value))
        with self.assertRaises(gate.SelectionError): gate.verify_plan(self.output, gate.digest(receipt.read_bytes()))

    def test_notice_manifest_drift_rejected(self):
        p = self.projection / 'payload' / self.notice_bindings['chownr']['path']
        p.write_bytes(p.read_bytes() + b' ')
        self.rejected()

    def test_notice_declaration_mismatch_rejected_even_with_matching_hash(self):
        # Independent semantic guard; projection verifier is mocked only for the fixture.
        for field, changed in [('name', 'another-work'), ('version', '9.0.0'), ('license', 'MIT')]:
            with self.subTest(field=field):
                p = self.projection / 'payload' / self.notice_bindings['chownr']['path']
                original = p.read_bytes(); raw = json.loads(original); raw[field] = changed
                p.write_bytes(gate.json_bytes(raw))
                before = self.notice_bindings['chownr']['sha256']
                self.notice_bindings['chownr']['sha256'] = gate.digest(p.read_bytes())
                try: self.rejected()
                finally: p.write_bytes(original); self.notice_bindings['chownr']['sha256'] = before

    def test_missing_notice_manifest_rejected(self):
        (self.projection / 'payload' / self.notice_bindings['chownr']['path']).unlink()
        self.rejected()

    def test_notice_write_failure_cleans_unpublished_stage(self):
        original = Path.write_bytes
        def fail_notice(path, data):
            if path.name == gate.NOTICE: raise PermissionError('injected notice write fault')
            return original(path, data)
        with patch.object(Path, 'write_bytes', fail_notice): self.rejected()


    def test_qrcode_notice_contains_author_permission_and_all_modules(self):
        sha = self.create(); value = gate.verify_plan(self.output, sha)
        raw = (self.output / gate.QRCODE_NOTICE).read_bytes()
        self.assertIn(b'Copyright (c) 2009 Kazuhiko Arase', raw)
        self.assertIn(b'Permission is hereby granted', raw)
        self.assertIn(b'The above copyright notice and this permission notice', raw)
        self.assertIn(b'THE SOFTWARE IS PROVIDED', raw)
        self.assertIn(gate.QRCODE_HEADER.encode('utf-8'), raw)
        for item in gate.QRCODE_MODULES: self.assertIn(item['path'].encode(), raw)
        self.assertEqual(gate.digest(raw), value['qrcode_notice']['sha256'])
        self.assertTrue((self.output / gate.NOTICE).is_file())

    def test_missing_qrcode_notice_rejected(self):
        sha = self.create(); (self.output / gate.QRCODE_NOTICE).unlink()
        with self.assertRaises(gate.SelectionError): gate.verify_plan(self.output, sha)

    def test_qrcode_author_removed_rejected(self):
        sha = self.create(); p = self.output / gate.QRCODE_NOTICE
        p.write_bytes(p.read_bytes().replace(b'Kazuhiko Arase', b'unknown author'))
        with self.assertRaises(gate.SelectionError): gate.verify_plan(self.output, sha)

    def test_qrcode_permission_removed_rejected(self):
        sha = self.create(); p = self.output / gate.QRCODE_NOTICE
        p.write_bytes(p.read_bytes().replace(gate.QRCODE_LICENSE.encode('utf-8'), b'MIT'))
        with self.assertRaises(gate.SelectionError): gate.verify_plan(self.output, sha)

    def test_qrcode_notice_and_receipt_rehashed_together_rejected(self):
        self.create(); p = self.output / gate.QRCODE_NOTICE; p.write_bytes(b'Copyright omitted')
        receipt = self.output / gate.RECIPE; value = json.loads(receipt.read_bytes())
        value['qrcode_notice']['sha256'] = gate.digest(p.read_bytes()); receipt.write_bytes(gate.json_bytes(value))
        with self.assertRaises(gate.SelectionError): gate.verify_plan(self.output, gate.digest(receipt.read_bytes()))

    def test_qrcode_bundle_missing_rejected(self):
        self.qrcode_bundle.unlink(); self.rejected()

    def test_qrcode_bundle_drift_rejected(self):
        self.qrcode_bundle.write_bytes(self.qrcode_bundle.read_bytes() + b'changed'); self.rejected()

    def test_qrcode_notice_write_failure_cleans_unpublished_stage(self):
        original = Path.write_bytes
        def fail_notice(path, data):
            if path.name == gate.QRCODE_NOTICE: raise PermissionError('injected QRCode notice write fault')
            return original(path, data)
        with patch.object(Path, 'write_bytes', fail_notice): self.rejected()


    def test_semver_original_license_delivered_and_scoped(self):
        sha = self.create(); value = gate.verify_plan(self.output, sha)
        raw = (self.output / gate.SEMVER_NOTICE).read_bytes()
        self.assertIn(gate.SEMVER_LICENSE.encode('utf-8'), raw)
        for text in [b'Copyright 2013 AJ ONeal', b'Permission is hereby granted', b'copies or substantial portions', b'THE SOFTWARE IS PROVIDED', b'APACHEv2', b'expired on 2025-01-29']:
            self.assertIn(text, raw)
        self.assertEqual(gate.digest(raw), value['semver_notice']['sha256'])
        self.assertFalse(value['runtime_admitted']); self.assertFalse(value['redistribution_admitted'])
    def test_missing_semver_notice_rejected(self):
        sha = self.create(); (self.output / gate.SEMVER_NOTICE).unlink()
        with self.assertRaises(gate.SelectionError): gate.verify_plan(self.output, sha)
    def test_semver_copyright_removed_rejected(self):
        sha = self.create(); p = self.output / gate.SEMVER_NOTICE
        p.write_bytes(p.read_bytes().replace(b'Copyright 2013 AJ ONeal', b''))
        with self.assertRaises(gate.SelectionError): gate.verify_plan(self.output, sha)
    def test_semver_permission_removed_rejected(self):
        sha = self.create(); p = self.output / gate.SEMVER_NOTICE
        p.write_bytes(p.read_bytes().replace(gate.SEMVER_LICENSE.encode('utf-8'), b'MIT'))
        with self.assertRaises(gate.SelectionError): gate.verify_plan(self.output, sha)
    def test_semver_notice_and_receipt_rehashed_together_rejected(self):
        self.create(); p = self.output / gate.SEMVER_NOTICE; p.write_bytes(b'MIT')
        receipt = self.output / gate.RECIPE; value = json.loads(receipt.read_bytes())
        value['semver_notice']['sha256'] = gate.digest(p.read_bytes()); receipt.write_bytes(gate.json_bytes(value))
        with self.assertRaises(gate.SelectionError): gate.verify_plan(self.output, gate.digest(receipt.read_bytes()))
    def test_semver_region_drift_rejected(self):
        with patch.dict(gate.SEMVER_REGION, {'region_sha256': '0' * 64}): self.rejected()
    def test_semver_original_license_drift_rejected(self):
        with patch.object(gate, 'SEMVER_LICENSE', 'MIT'): self.rejected()
    def test_semver_notice_write_failure_cleans_unpublished_stage(self):
        original = Path.write_bytes
        def fail_notice(path, data):
            if path.name == gate.SEMVER_NOTICE: raise PermissionError('injected semver notice write fault')
            return original(path, data)
        with patch.object(Path, 'write_bytes', fail_notice): self.rejected()


    def test_retained_collection_delivers_all_exact_texts(self):
        sha = self.create(); value = gate.verify_plan(self.output, sha)
        raw = (self.output / gate.RETAINED_NOTICE).read_bytes()
        for entry in self.catalog['entries']:
            self.assertIn(base64.b64decode(entry['content_base64']), raw)
            self.assertIn(entry['sha256'].encode(), raw)
        self.assertEqual(47, value['retained_notices']['copies'])
        self.assertEqual('RETAINED_TEXTS_ONLY', value['retained_notices']['scope'])
        self.assertFalse(value['redistribution_admitted'])
    def test_retained_delivery_missing_rejected(self):
        sha = self.create(); (self.output / gate.RETAINED_NOTICE).unlink()
        with self.assertRaises(gate.SelectionError): gate.verify_plan(self.output, sha)
    def test_retained_permission_removed_rejected(self):
        sha = self.create(); p = self.output / gate.RETAINED_NOTICE
        p.write_bytes(p.read_bytes().replace(b'Permission', b'Omitted'))
        with self.assertRaises(gate.SelectionError): gate.verify_plan(self.output, sha)
    def test_retained_delivery_and_receipt_rehashed_rejected(self):
        self.create(); p = self.output / gate.RETAINED_NOTICE; p.write_bytes(b'All licensed')
        receipt = self.output / gate.RECIPE; value = json.loads(receipt.read_bytes())
        value['retained_notices']['sha256'] = gate.digest(p.read_bytes()); receipt.write_bytes(gate.json_bytes(value))
        with self.assertRaises(gate.SelectionError): gate.verify_plan(self.output, gate.digest(receipt.read_bytes()))
    def test_original_retained_notice_drift_rejected(self):
        entry = next(e for e in self.catalog['entries'] if e['origin'] == 'ORIGINAL_SELECTED_PAYLOAD')
        (self.projection / 'payload' / entry['payload_path']).write_bytes(b'changed')
        self.rejected()
    def test_retained_write_failure_cleans_unpublished_stage(self):
        original = Path.write_bytes
        def fail(path, data):
            if path.name == gate.RETAINED_NOTICE: raise PermissionError('injected retained notice write fault')
            return original(path, data)
        with patch.object(Path, 'write_bytes', fail): self.rejected()
    def catalog_rejected(self, mutate):
        value = json.loads(json.dumps(self.catalog)); mutate(value); raw = gate.json_bytes(value)
        original = gate.read_bounded
        def read(path, limit):
            if Path(path).name == 'retained-notice-catalog.json': return raw
            return original(path, limit)
        with patch.object(gate, 'read_bounded', read), patch.object(gate, 'CATALOG_SHA256', gate.digest(raw)):
            self.rejected()
    def test_catalog_duplicate_path_rejected(self):
        self.catalog_rejected(lambda v: v['entries'].__setitem__(1, v['entries'][0]))
    def test_catalog_traversal_rejected(self):
        self.catalog_rejected(lambda v: v['entries'][0].update(path='../outside'))
    def test_catalog_inflated_status_rejected(self):
        self.catalog_rejected(lambda v: v.update(release_ready=True))
    def test_catalog_content_tampered_rejected(self):
        self.catalog_rejected(lambda v: v['entries'][0].update(content_base64='YQ=='))
    def test_catalog_invalid_encoding_rejected(self):
        self.catalog_rejected(lambda v: v['entries'][0].update(content_base64='!!!!'))
    def test_catalog_origin_inflated_rejected(self):
        self.catalog_rejected(lambda v: v['entries'][22].update(origin='ORIGINAL_SELECTED_PAYLOAD'))
    def test_catalog_count_drift_rejected(self):
        self.catalog_rejected(lambda v: v['entries'].pop())
    def test_catalog_extra_field_rejected(self):
        self.catalog_rejected(lambda v: v.update(licensed=True))
    def test_catalog_wrong_hash_rejected(self):
        with patch.object(gate, 'CATALOG_SHA256', '0'*64): self.rejected()


    def test_next_path_original_source_manifest_license_delivered(self):
        sha = self.create(); value = gate.verify_plan(self.output, sha)
        raw = (self.output / gate.NEXT_PATH_NOTICE).read_bytes()
        data = json.loads(Path(gate.__file__).with_name('next-path-source-evidence.json').read_bytes())
        for entry in data['files']: self.assertIn(base64.b64decode(entry['content_base64']), raw)
        self.assertEqual('SOURCE_AND_LICENSE_DELIVERY_ONLY', value['next_path_source']['scope'])
        self.assertIn(b'conditional loader adapter', raw)
        self.assertFalse(value['redistribution_admitted'])
    def test_next_path_source_notice_missing_rejected(self):
        sha = self.create(); (self.output / gate.NEXT_PATH_NOTICE).unlink()
        with self.assertRaises(gate.SelectionError): gate.verify_plan(self.output, sha)
    def test_next_path_source_body_mutation_rejected(self):
        sha = self.create(); p = self.output / gate.NEXT_PATH_NOTICE
        p.write_bytes(p.read_bytes().replace(b'path.relative(from, to)', b'path.relative(to, from)'))
        with self.assertRaises(gate.SelectionError): gate.verify_plan(self.output, sha)
    def test_next_path_full_license_omission_rejected(self):
        sha = self.create(); p = self.output / gate.NEXT_PATH_NOTICE
        data = json.loads(Path(gate.__file__).with_name('next-path-source-evidence.json').read_bytes())
        original = base64.b64decode(data['files'][2]['content_base64'])
        p.write_bytes(p.read_bytes().replace(original, b'MPL-2.0'))
        with self.assertRaises(gate.SelectionError): gate.verify_plan(self.output, sha)
    def test_next_path_notice_and_receipt_rehashed_rejected(self):
        self.create(); p = self.output / gate.NEXT_PATH_NOTICE; p.write_bytes(b'Source omitted')
        receipt = self.output / gate.RECIPE; value = json.loads(receipt.read_bytes())
        value['next_path_source']['sha256'] = gate.digest(p.read_bytes()); receipt.write_bytes(gate.json_bytes(value))
        with self.assertRaises(gate.SelectionError): gate.verify_plan(self.output, gate.digest(receipt.read_bytes()))
    def test_next_path_bundle_region_changed_rejected(self):
        with patch.dict(gate.NEXT_PATH_REGION, {'sha256': '0'*64}): self.rejected()
    def test_next_path_evidence_changed_rejected(self):
        with patch.object(gate, 'NEXT_PATH_DATA_SHA256', '0'*64): self.rejected()
    def test_next_path_write_failure_cleans_stage(self):
        original = Path.write_bytes
        def fail(path, raw):
            if path.name == gate.NEXT_PATH_NOTICE: raise PermissionError('injected next-path source write fault')
            return original(path, raw)
        with patch.object(Path, 'write_bytes', fail): self.rejected()

if __name__ == '__main__': unittest.main()
````

### FILE: `pnpm_artifact_selection/retained-notice-catalog.json`
```yaml
block_id: "PNPM-ARTIFACT-SELECTION-GATE:retained-notice-catalog:v397"
operation: CREATE
provenance: ADAPTED
source: "AUTHORED portable JSON envelope; VERBATIM47license and attribution texts from fixed V340manifest 27d666e86d769c368439b1cc20e0fa72a42f9366973216b42ced5434f003fa9d;22original pnpm payload notices plus25research copies, each retains scope/source locator/hash; no executable upstream algorithm"
license: "LicenseRef-Pnpm-Bundled-Licenses"
sha256: "e4a906897aeae3dec44d63beb40514329353f41859c8d380c126df8b2d599c1e"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-pnpm-retained-notice-catalog/v1",
  "status": "RETAINED_TEXTS_ONLY",
  "release_ready": false,
  "parent_manifest_sha256": "27d666e86d769c368439b1cc20e0fa72a42f9366973216b42ced5434f003fa9d",
  "entries": [
    {
      "path": "original/dist/node_modules/@isaacs/fs-minipass/LICENSE",
      "source_locator": "Temp/elite-v334-a443d28e2128479b87937c040aa8b6a8/projection-rebuilt/payload/dist/node_modules/@isaacs/fs-minipass/LICENSE",
      "scope": "VERBATIM existing selected payload notice; applicability remains per path",
      "origin": "ORIGINAL_SELECTED_PAYLOAD",
      "payload_path": "dist/node_modules/@isaacs/fs-minipass/LICENSE",
      "bytes": 765,
      "sha256": "4ec3d4c66cd87f5c8d8ad911b10f99bf27cb00cdfcff82621956e379186b016b",
      "content_base64": "VGhlIElTQyBMaWNlbnNlCgpDb3B5cmlnaHQgKGMpIElzYWFjIFouIFNjaGx1ZXRlciBhbmQgQ29udHJpYnV0b3JzCgpQZXJtaXNzaW9uIHRvIHVzZSwgY29weSwgbW9kaWZ5LCBhbmQvb3IgZGlzdHJpYnV0ZSB0aGlzIHNvZnR3YXJlIGZvciBhbnkKcHVycG9zZSB3aXRoIG9yIHdpdGhvdXQgZmVlIGlzIGhlcmVieSBncmFudGVkLCBwcm92aWRlZCB0aGF0IHRoZSBhYm92ZQpjb3B5cmlnaHQgbm90aWNlIGFuZCB0aGlzIHBlcm1pc3Npb24gbm90aWNlIGFwcGVhciBpbiBhbGwgY29waWVzLgoKVEhFIFNPRlRXQVJFIElTIFBST1ZJREVEICJBUyBJUyIgQU5EIFRIRSBBVVRIT1IgRElTQ0xBSU1TIEFMTCBXQVJSQU5USUVTCldJVEggUkVHQVJEIFRPIFRISVMgU09GVFdBUkUgSU5DTFVESU5HIEFMTCBJTVBMSUVEIFdBUlJBTlRJRVMgT0YKTUVSQ0hBTlRBQklMSVRZIEFORCBGSVRORVNTLiBJTiBOTyBFVkVOVCBTSEFMTCBUSEUgQVVUSE9SIEJFIExJQUJMRSBGT1IKQU5ZIFNQRUNJQUwsIERJUkVDVCwgSU5ESVJFQ1QsIE9SIENPTlNFUVVFTlRJQUwgREFNQUdFUyBPUiBBTlkgREFNQUdFUwpXSEFUU09FVkVSIFJFU1VMVElORyBGUk9NIExPU1MgT0YgVVNFLCBEQVRBIE9SIFBST0ZJVFMsIFdIRVRIRVIgSU4gQU4KQUNUSU9OIE9GIENPTlRSQUNULCBORUdMSUdFTkNFIE9SIE9USEVSIFRPUlRJT1VTIEFDVElPTiwgQVJJU0lORyBPVVQgT0YgT1IKSU4gQ09OTkVDVElPTiBXSVRIIFRIRSBVU0UgT1IgUEVSRk9STUFOQ0UgT0YgVEhJUyBTT0ZUV0FSRS4K"
    },
    {
      "path": "original/dist/node_modules/abbrev/LICENSE",
      "source_locator": "Temp/elite-v334-a443d28e2128479b87937c040aa8b6a8/projection-rebuilt/payload/dist/node_modules/abbrev/LICENSE",
      "scope": "VERBATIM existing selected payload notice; applicability remains per path",
      "origin": "ORIGINAL_SELECTED_PAYLOAD",
      "payload_path": "dist/node_modules/abbrev/LICENSE",
      "bytes": 2011,
      "sha256": "9e0d5c7989f7e9f07d7c4b158aceff270f235eb7464ace41c5e7b200834a43e0",
      "content_base64": "VGhpcyBzb2Z0d2FyZSBpcyBkdWFsLWxpY2Vuc2VkIHVuZGVyIHRoZSBJU0MgYW5kIE1JVCBsaWNlbnNlcy4KWW91IG1heSB1c2UgdGhpcyBzb2Z0d2FyZSB1bmRlciBFSVRIRVIgb2YgdGhlIGZvbGxvd2luZyBsaWNlbnNlcy4KCi0tLS0tLS0tLS0KClRoZSBJU0MgTGljZW5zZQoKQ29weXJpZ2h0IChjKSBJc2FhYyBaLiBTY2hsdWV0ZXIgYW5kIENvbnRyaWJ1dG9ycwoKUGVybWlzc2lvbiB0byB1c2UsIGNvcHksIG1vZGlmeSwgYW5kL29yIGRpc3RyaWJ1dGUgdGhpcyBzb2Z0d2FyZSBmb3IgYW55CnB1cnBvc2Ugd2l0aCBvciB3aXRob3V0IGZlZSBpcyBoZXJlYnkgZ3JhbnRlZCwgcHJvdmlkZWQgdGhhdCB0aGUgYWJvdmUKY29weXJpZ2h0IG5vdGljZSBhbmQgdGhpcyBwZXJtaXNzaW9uIG5vdGljZSBhcHBlYXIgaW4gYWxsIGNvcGllcy4KClRIRSBTT0ZUV0FSRSBJUyBQUk9WSURFRCAiQVMgSVMiIEFORCBUSEUgQVVUSE9SIERJU0NMQUlNUyBBTEwgV0FSUkFOVElFUwpXSVRIIFJFR0FSRCBUTyBUSElTIFNPRlRXQVJFIElOQ0xVRElORyBBTEwgSU1QTElFRCBXQVJSQU5USUVTIE9GCk1FUkNIQU5UQUJJTElUWSBBTkQgRklUTkVTUy4gSU4gTk8gRVZFTlQgU0hBTEwgVEhFIEFVVEhPUiBCRSBMSUFCTEUgRk9SCkFOWSBTUEVDSUFMLCBESVJFQ1QsIElORElSRUNULCBPUiBDT05TRVFVRU5USUFMIERBTUFHRVMgT1IgQU5ZIERBTUFHRVMKV0hBVFNPRVZFUiBSRVNVTFRJTkcgRlJPTSBMT1NTIE9GIFVTRSwgREFUQSBPUiBQUk9GSVRTLCBXSEVUSEVSIElOIEFOCkFDVElPTiBPRiBDT05UUkFDVCwgTkVHTElHRU5DRSBPUiBPVEhFUiBUT1JUSU9VUyBBQ1RJT04sIEFSSVNJTkcgT1VUIE9GIE9SCklOIENPTk5FQ1RJT04gV0lUSCBUSEUgVVNFIE9SIFBFUkZPUk1BTkNFIE9GIFRISVMgU09GVFdBUkUuCgotLS0tLS0tLS0tCgpDb3B5cmlnaHQgSXNhYWMgWi4gU2NobHVldGVyIGFuZCBDb250cmlidXRvcnMKQWxsIHJpZ2h0cyByZXNlcnZlZC4KClBlcm1pc3Npb24gaXMgaGVyZWJ5IGdyYW50ZWQsIGZyZWUgb2YgY2hhcmdlLCB0byBhbnkgcGVyc29uCm9idGFpbmluZyBhIGNvcHkgb2YgdGhpcyBzb2Z0d2FyZSBhbmQgYXNzb2NpYXRlZCBkb2N1bWVudGF0aW9uCmZpbGVzICh0aGUgIlNvZnR3YXJlIiksIHRvIGRlYWwgaW4gdGhlIFNvZnR3YXJlIHdpdGhvdXQKcmVzdHJpY3Rpb24sIGluY2x1ZGluZyB3aXRob3V0IGxpbWl0YXRpb24gdGhlIHJpZ2h0cyB0byB1c2UsCmNvcHksIG1vZGlmeSwgbWVyZ2UsIHB1Ymxpc2gsIGRpc3RyaWJ1dGUsIHN1YmxpY2Vuc2UsIGFuZC9vciBzZWxsCmNvcGllcyBvZiB0aGUgU29mdHdhcmUsIGFuZCB0byBwZXJtaXQgcGVyc29ucyB0byB3aG9tIHRoZQpTb2Z0d2FyZSBpcyBmdXJuaXNoZWQgdG8gZG8gc28sIHN1YmplY3QgdG8gdGhlIGZvbGxvd2luZwpjb25kaXRpb25zOgoKVGhlIGFib3ZlIGNvcHlyaWdodCBub3RpY2UgYW5kIHRoaXMgcGVybWlzc2lvbiBub3RpY2Ugc2hhbGwgYmUKaW5jbHVkZWQgaW4gYWxsIGNvcGllcyBvciBzdWJzdGFudGlhbCBwb3J0aW9ucyBvZiB0aGUgU29mdHdhcmUuCgpUSEUgU09GVFdBUkUgSVMgUFJPVklERUQgIkFTIElTIiwgV0lUSE9VVCBXQVJSQU5UWSBPRiBBTlkgS0lORCwKRVhQUkVTUyBPUiBJTVBMSUVELCBJTkNMVURJTkcgQlVUIE5PVCBMSU1JVEVEIFRPIFRIRSBXQVJSQU5USUVTCk9GIE1FUkNIQU5UQUJJTElUWSwgRklUTkVTUyBGT1IgQSBQQVJUSUNVTEFSIFBVUlBPU0UgQU5ECk5PTklORlJJTkdFTUVOVC4gSU4gTk8gRVZFTlQgU0hBTEwgVEhFIEFVVEhPUlMgT1IgQ09QWVJJR0hUCkhPTERFUlMgQkUgTElBQkxFIEZPUiBBTlkgQ0xBSU0sIERBTUFHRVMgT1IgT1RIRVIgTElBQklMSVRZLApXSEVUSEVSIElOIEFOIEFDVElPTiBPRiBDT05UUkFDVCwgVE9SVCBPUiBPVEhFUldJU0UsIEFSSVNJTkcKRlJPTSwgT1VUIE9GIE9SIElOIENPTk5FQ1RJT04gV0lUSCBUSEUgU09GVFdBUkUgT1IgVEhFIFVTRSBPUgpPVEhFUiBERUFMSU5HUyBJTiBUSEUgU09GVFdBUkUuCg=="
    },
    {
      "path": "original/dist/node_modules/env-paths/license",
      "source_locator": "Temp/elite-v334-a443d28e2128479b87937c040aa8b6a8/projection-rebuilt/payload/dist/node_modules/env-paths/license",
      "scope": "VERBATIM existing selected payload notice; applicability remains per path",
      "origin": "ORIGINAL_SELECTED_PAYLOAD",
      "payload_path": "dist/node_modules/env-paths/license",
      "bytes": 1109,
      "sha256": "48da2f39e100d4085767e94966b43f4fa95ff6a0698fba57ed460914e35f94a0",
      "content_base64": "TUlUIExpY2Vuc2UKCkNvcHlyaWdodCAoYykgU2luZHJlIFNvcmh1cyA8c2luZHJlc29yaHVzQGdtYWlsLmNvbT4gKHNpbmRyZXNvcmh1cy5jb20pCgpQZXJtaXNzaW9uIGlzIGhlcmVieSBncmFudGVkLCBmcmVlIG9mIGNoYXJnZSwgdG8gYW55IHBlcnNvbiBvYnRhaW5pbmcgYSBjb3B5IG9mIHRoaXMgc29mdHdhcmUgYW5kIGFzc29jaWF0ZWQgZG9jdW1lbnRhdGlvbiBmaWxlcyAodGhlICJTb2Z0d2FyZSIpLCB0byBkZWFsIGluIHRoZSBTb2Z0d2FyZSB3aXRob3V0IHJlc3RyaWN0aW9uLCBpbmNsdWRpbmcgd2l0aG91dCBsaW1pdGF0aW9uIHRoZSByaWdodHMgdG8gdXNlLCBjb3B5LCBtb2RpZnksIG1lcmdlLCBwdWJsaXNoLCBkaXN0cmlidXRlLCBzdWJsaWNlbnNlLCBhbmQvb3Igc2VsbCBjb3BpZXMgb2YgdGhlIFNvZnR3YXJlLCBhbmQgdG8gcGVybWl0IHBlcnNvbnMgdG8gd2hvbSB0aGUgU29mdHdhcmUgaXMgZnVybmlzaGVkIHRvIGRvIHNvLCBzdWJqZWN0IHRvIHRoZSBmb2xsb3dpbmcgY29uZGl0aW9uczoKClRoZSBhYm92ZSBjb3B5cmlnaHQgbm90aWNlIGFuZCB0aGlzIHBlcm1pc3Npb24gbm90aWNlIHNoYWxsIGJlIGluY2x1ZGVkIGluIGFsbCBjb3BpZXMgb3Igc3Vic3RhbnRpYWwgcG9ydGlvbnMgb2YgdGhlIFNvZnR3YXJlLgoKVEhFIFNPRlRXQVJFIElTIFBST1ZJREVEICJBUyBJUyIsIFdJVEhPVVQgV0FSUkFOVFkgT0YgQU5ZIEtJTkQsIEVYUFJFU1MgT1IgSU1QTElFRCwgSU5DTFVESU5HIEJVVCBOT1QgTElNSVRFRCBUTyBUSEUgV0FSUkFOVElFUyBPRiBNRVJDSEFOVEFCSUxJVFksIEZJVE5FU1MgRk9SIEEgUEFSVElDVUxBUiBQVVJQT1NFIEFORCBOT05JTkZSSU5HRU1FTlQuIElOIE5PIEVWRU5UIFNIQUxMIFRIRSBBVVRIT1JTIE9SIENPUFlSSUdIVCBIT0xERVJTIEJFIExJQUJMRSBGT1IgQU5ZIENMQUlNLCBEQU1BR0VTIE9SIE9USEVSIExJQUJJTElUWSwgV0hFVEhFUiBJTiBBTiBBQ1RJT04gT0YgQ09OVFJBQ1QsIFRPUlQgT1IgT1RIRVJXSVNFLCBBUklTSU5HIEZST00sIE9VVCBPRiBPUiBJTiBDT05ORUNUSU9OIFdJVEggVEhFIFNPRlRXQVJFIE9SIFRIRSBVU0UgT1IgT1RIRVIgREVBTElOR1MgSU4gVEhFIFNPRlRXQVJFLgo="
    },
    {
      "path": "original/dist/node_modules/exponential-backoff/LICENSE",
      "source_locator": "Temp/elite-v334-a443d28e2128479b87937c040aa8b6a8/projection-rebuilt/payload/dist/node_modules/exponential-backoff/LICENSE",
      "scope": "VERBATIM existing selected payload notice; applicability remains per path",
      "origin": "ORIGINAL_SELECTED_PAYLOAD",
      "payload_path": "dist/node_modules/exponential-backoff/LICENSE",
      "bytes": 11351,
      "sha256": "f5afc95c69e116ea6b1755f91a8763f53eadb78b4f1582a08d2d18586a03a117",
      "content_base64": "CiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIEFwYWNoZSBMaWNlbnNlCiAgICAgICAgICAgICAgICAgICAgICAgICAgIFZlcnNpb24gMi4wLCBKYW51YXJ5IDIwMDQKICAgICAgICAgICAgICAgICAgICAgICAgaHR0cDovL3d3dy5hcGFjaGUub3JnL2xpY2Vuc2VzLwoKICAgVEVSTVMgQU5EIENPTkRJVElPTlMgRk9SIFVTRSwgUkVQUk9EVUNUSU9OLCBBTkQgRElTVFJJQlVUSU9OCgogICAxLiBEZWZpbml0aW9ucy4KCiAgICAgICJMaWNlbnNlIiBzaGFsbCBtZWFuIHRoZSB0ZXJtcyBhbmQgY29uZGl0aW9ucyBmb3IgdXNlLCByZXByb2R1Y3Rpb24sCiAgICAgIGFuZCBkaXN0cmlidXRpb24gYXMgZGVmaW5lZCBieSBTZWN0aW9ucyAxIHRocm91Z2ggOSBvZiB0aGlzIGRvY3VtZW50LgoKICAgICAgIkxpY2Vuc29yIiBzaGFsbCBtZWFuIHRoZSBjb3B5cmlnaHQgb3duZXIgb3IgZW50aXR5IGF1dGhvcml6ZWQgYnkKICAgICAgdGhlIGNvcHlyaWdodCBvd25lciB0aGF0IGlzIGdyYW50aW5nIHRoZSBMaWNlbnNlLgoKICAgICAgIkxlZ2FsIEVudGl0eSIgc2hhbGwgbWVhbiB0aGUgdW5pb24gb2YgdGhlIGFjdGluZyBlbnRpdHkgYW5kIGFsbAogICAgICBvdGhlciBlbnRpdGllcyB0aGF0IGNvbnRyb2wsIGFyZSBjb250cm9sbGVkIGJ5LCBvciBhcmUgdW5kZXIgY29tbW9uCiAgICAgIGNvbnRyb2wgd2l0aCB0aGF0IGVudGl0eS4gRm9yIHRoZSBwdXJwb3NlcyBvZiB0aGlzIGRlZmluaXRpb24sCiAgICAgICJjb250cm9sIiBtZWFucyAoaSkgdGhlIHBvd2VyLCBkaXJlY3Qgb3IgaW5kaXJlY3QsIHRvIGNhdXNlIHRoZQogICAgICBkaXJlY3Rpb24gb3IgbWFuYWdlbWVudCBvZiBzdWNoIGVudGl0eSwgd2hldGhlciBieSBjb250cmFjdCBvcgogICAgICBvdGhlcndpc2UsIG9yIChpaSkgb3duZXJzaGlwIG9mIGZpZnR5IHBlcmNlbnQgKDUwJSkgb3IgbW9yZSBvZiB0aGUKICAgICAgb3V0c3RhbmRpbmcgc2hhcmVzLCBvciAoaWlpKSBiZW5lZmljaWFsIG93bmVyc2hpcCBvZiBzdWNoIGVudGl0eS4KCiAgICAgICJZb3UiIChvciAiWW91ciIpIHNoYWxsIG1lYW4gYW4gaW5kaXZpZHVhbCBvciBMZWdhbCBFbnRpdHkKICAgICAgZXhlcmNpc2luZyBwZXJtaXNzaW9ucyBncmFudGVkIGJ5IHRoaXMgTGljZW5zZS4KCiAgICAgICJTb3VyY2UiIGZvcm0gc2hhbGwgbWVhbiB0aGUgcHJlZmVycmVkIGZvcm0gZm9yIG1ha2luZyBtb2RpZmljYXRpb25zLAogICAgICBpbmNsdWRpbmcgYnV0IG5vdCBsaW1pdGVkIHRvIHNvZnR3YXJlIHNvdXJjZSBjb2RlLCBkb2N1bWVudGF0aW9uCiAgICAgIHNvdXJjZSwgYW5kIGNvbmZpZ3VyYXRpb24gZmlsZXMuCgogICAgICAiT2JqZWN0IiBmb3JtIHNoYWxsIG1lYW4gYW55IGZvcm0gcmVzdWx0aW5nIGZyb20gbWVjaGFuaWNhbAogICAgICB0cmFuc2Zvcm1hdGlvbiBvciB0cmFuc2xhdGlvbiBvZiBhIFNvdXJjZSBmb3JtLCBpbmNsdWRpbmcgYnV0CiAgICAgIG5vdCBsaW1pdGVkIHRvIGNvbXBpbGVkIG9iamVjdCBjb2RlLCBnZW5lcmF0ZWQgZG9jdW1lbnRhdGlvbiwKICAgICAgYW5kIGNvbnZlcnNpb25zIHRvIG90aGVyIG1lZGlhIHR5cGVzLgoKICAgICAgIldvcmsiIHNoYWxsIG1lYW4gdGhlIHdvcmsgb2YgYXV0aG9yc2hpcCwgd2hldGhlciBpbiBTb3VyY2Ugb3IKICAgICAgT2JqZWN0IGZvcm0sIG1hZGUgYXZhaWxhYmxlIHVuZGVyIHRoZSBMaWNlbnNlLCBhcyBpbmRpY2F0ZWQgYnkgYQogICAgICBjb3B5cmlnaHQgbm90aWNlIHRoYXQgaXMgaW5jbHVkZWQgaW4gb3IgYXR0YWNoZWQgdG8gdGhlIHdvcmsKICAgICAgKGFuIGV4YW1wbGUgaXMgcHJvdmlkZWQgaW4gdGhlIEFwcGVuZGl4IGJlbG93KS4KCiAgICAgICJEZXJpdmF0aXZlIFdvcmtzIiBzaGFsbCBtZWFuIGFueSB3b3JrLCB3aGV0aGVyIGluIFNvdXJjZSBvciBPYmplY3QKICAgICAgZm9ybSwgdGhhdCBpcyBiYXNlZCBvbiAob3IgZGVyaXZlZCBmcm9tKSB0aGUgV29yayBhbmQgZm9yIHdoaWNoIHRoZQogICAgICBlZGl0b3JpYWwgcmV2aXNpb25zLCBhbm5vdGF0aW9ucywgZWxhYm9yYXRpb25zLCBvciBvdGhlciBtb2RpZmljYXRpb25zCiAgICAgIHJlcHJlc2VudCwgYXMgYSB3aG9sZSwgYW4gb3JpZ2luYWwgd29yayBvZiBhdXRob3JzaGlwLiBGb3IgdGhlIHB1cnBvc2VzCiAgICAgIG9mIHRoaXMgTGljZW5zZSwgRGVyaXZhdGl2ZSBXb3JrcyBzaGFsbCBub3QgaW5jbHVkZSB3b3JrcyB0aGF0IHJlbWFpbgogICAgICBzZXBhcmFibGUgZnJvbSwgb3IgbWVyZWx5IGxpbmsgKG9yIGJpbmQgYnkgbmFtZSkgdG8gdGhlIGludGVyZmFjZXMgb2YsCiAgICAgIHRoZSBXb3JrIGFuZCBEZXJpdmF0aXZlIFdvcmtzIHRoZXJlb2YuCgogICAgICAiQ29udHJpYnV0aW9uIiBzaGFsbCBtZWFuIGFueSB3b3JrIG9mIGF1dGhvcnNoaXAsIGluY2x1ZGluZwogICAgICB0aGUgb3JpZ2luYWwgdmVyc2lvbiBvZiB0aGUgV29yayBhbmQgYW55IG1vZGlmaWNhdGlvbnMgb3IgYWRkaXRpb25zCiAgICAgIHRvIHRoYXQgV29yayBvciBEZXJpdmF0aXZlIFdvcmtzIHRoZXJlb2YsIHRoYXQgaXMgaW50ZW50aW9uYWxseQogICAgICBzdWJtaXR0ZWQgdG8gTGljZW5zb3IgZm9yIGluY2x1c2lvbiBpbiB0aGUgV29yayBieSB0aGUgY29weXJpZ2h0IG93bmVyCiAgICAgIG9yIGJ5IGFuIGluZGl2aWR1YWwgb3IgTGVnYWwgRW50aXR5IGF1dGhvcml6ZWQgdG8gc3VibWl0IG9uIGJlaGFsZiBvZgogICAgICB0aGUgY29weXJpZ2h0IG93bmVyLiBGb3IgdGhlIHB1cnBvc2VzIG9mIHRoaXMgZGVmaW5pdGlvbiwgInN1Ym1pdHRlZCIKICAgICAgbWVhbnMgYW55IGZvcm0gb2YgZWxlY3Ryb25pYywgdmVyYmFsLCBvciB3cml0dGVuIGNvbW11bmljYXRpb24gc2VudAogICAgICB0byB0aGUgTGljZW5zb3Igb3IgaXRzIHJlcHJlc2VudGF0aXZlcywgaW5jbHVkaW5nIGJ1dCBub3QgbGltaXRlZCB0bwogICAgICBjb21tdW5pY2F0aW9uIG9uIGVsZWN0cm9uaWMgbWFpbGluZyBsaXN0cywgc291cmNlIGNvZGUgY29udHJvbCBzeXN0ZW1zLAogICAgICBhbmQgaXNzdWUgdHJhY2tpbmcgc3lzdGVtcyB0aGF0IGFyZSBtYW5hZ2VkIGJ5LCBvciBvbiBiZWhhbGYgb2YsIHRoZQogICAgICBMaWNlbnNvciBmb3IgdGhlIHB1cnBvc2Ugb2YgZGlzY3Vzc2luZyBhbmQgaW1wcm92aW5nIHRoZSBXb3JrLCBidXQKICAgICAgZXhjbHVkaW5nIGNvbW11bmljYXRpb24gdGhhdCBpcyBjb25zcGljdW91c2x5IG1hcmtlZCBvciBvdGhlcndpc2UKICAgICAgZGVzaWduYXRlZCBpbiB3cml0aW5nIGJ5IHRoZSBjb3B5cmlnaHQgb3duZXIgYXMgIk5vdCBhIENvbnRyaWJ1dGlvbi4iCgogICAgICAiQ29udHJpYnV0b3IiIHNoYWxsIG1lYW4gTGljZW5zb3IgYW5kIGFueSBpbmRpdmlkdWFsIG9yIExlZ2FsIEVudGl0eQogICAgICBvbiBiZWhhbGYgb2Ygd2hvbSBhIENvbnRyaWJ1dGlvbiBoYXMgYmVlbiByZWNlaXZlZCBieSBMaWNlbnNvciBhbmQKICAgICAgc3Vic2VxdWVudGx5IGluY29ycG9yYXRlZCB3aXRoaW4gdGhlIFdvcmsuCgogICAyLiBHcmFudCBvZiBDb3B5cmlnaHQgTGljZW5zZS4gU3ViamVjdCB0byB0aGUgdGVybXMgYW5kIGNvbmRpdGlvbnMgb2YKICAgICAgdGhpcyBMaWNlbnNlLCBlYWNoIENvbnRyaWJ1dG9yIGhlcmVieSBncmFudHMgdG8gWW91IGEgcGVycGV0dWFsLAogICAgICB3b3JsZHdpZGUsIG5vbi1leGNsdXNpdmUsIG5vLWNoYXJnZSwgcm95YWx0eS1mcmVlLCBpcnJldm9jYWJsZQogICAgICBjb3B5cmlnaHQgbGljZW5zZSB0byByZXByb2R1Y2UsIHByZXBhcmUgRGVyaXZhdGl2ZSBXb3JrcyBvZiwKICAgICAgcHVibGljbHkgZGlzcGxheSwgcHVibGljbHkgcGVyZm9ybSwgc3VibGljZW5zZSwgYW5kIGRpc3RyaWJ1dGUgdGhlCiAgICAgIFdvcmsgYW5kIHN1Y2ggRGVyaXZhdGl2ZSBXb3JrcyBpbiBTb3VyY2Ugb3IgT2JqZWN0IGZvcm0uCgogICAzLiBHcmFudCBvZiBQYXRlbnQgTGljZW5zZS4gU3ViamVjdCB0byB0aGUgdGVybXMgYW5kIGNvbmRpdGlvbnMgb2YKICAgICAgdGhpcyBMaWNlbnNlLCBlYWNoIENvbnRyaWJ1dG9yIGhlcmVieSBncmFudHMgdG8gWW91IGEgcGVycGV0dWFsLAogICAgICB3b3JsZHdpZGUsIG5vbi1leGNsdXNpdmUsIG5vLWNoYXJnZSwgcm95YWx0eS1mcmVlLCBpcnJldm9jYWJsZQogICAgICAoZXhjZXB0IGFzIHN0YXRlZCBpbiB0aGlzIHNlY3Rpb24pIHBhdGVudCBsaWNlbnNlIHRvIG1ha2UsIGhhdmUgbWFkZSwKICAgICAgdXNlLCBvZmZlciB0byBzZWxsLCBzZWxsLCBpbXBvcnQsIGFuZCBvdGhlcndpc2UgdHJhbnNmZXIgdGhlIFdvcmssCiAgICAgIHdoZXJlIHN1Y2ggbGljZW5zZSBhcHBsaWVzIG9ubHkgdG8gdGhvc2UgcGF0ZW50IGNsYWltcyBsaWNlbnNhYmxlCiAgICAgIGJ5IHN1Y2ggQ29udHJpYnV0b3IgdGhhdCBhcmUgbmVjZXNzYXJpbHkgaW5mcmluZ2VkIGJ5IHRoZWlyCiAgICAgIENvbnRyaWJ1dGlvbihzKSBhbG9uZSBvciBieSBjb21iaW5hdGlvbiBvZiB0aGVpciBDb250cmlidXRpb24ocykKICAgICAgd2l0aCB0aGUgV29yayB0byB3aGljaCBzdWNoIENvbnRyaWJ1dGlvbihzKSB3YXMgc3VibWl0dGVkLiBJZiBZb3UKICAgICAgaW5zdGl0dXRlIHBhdGVudCBsaXRpZ2F0aW9uIGFnYWluc3QgYW55IGVudGl0eSAoaW5jbHVkaW5nIGEKICAgICAgY3Jvc3MtY2xhaW0gb3IgY291bnRlcmNsYWltIGluIGEgbGF3c3VpdCkgYWxsZWdpbmcgdGhhdCB0aGUgV29yawogICAgICBvciBhIENvbnRyaWJ1dGlvbiBpbmNvcnBvcmF0ZWQgd2l0aGluIHRoZSBXb3JrIGNvbnN0aXR1dGVzIGRpcmVjdAogICAgICBvciBjb250cmlidXRvcnkgcGF0ZW50IGluZnJpbmdlbWVudCwgdGhlbiBhbnkgcGF0ZW50IGxpY2Vuc2VzCiAgICAgIGdyYW50ZWQgdG8gWW91IHVuZGVyIHRoaXMgTGljZW5zZSBmb3IgdGhhdCBXb3JrIHNoYWxsIHRlcm1pbmF0ZQogICAgICBhcyBvZiB0aGUgZGF0ZSBzdWNoIGxpdGlnYXRpb24gaXMgZmlsZWQuCgogICA0LiBSZWRpc3RyaWJ1dGlvbi4gWW91IG1heSByZXByb2R1Y2UgYW5kIGRpc3RyaWJ1dGUgY29waWVzIG9mIHRoZQogICAgICBXb3JrIG9yIERlcml2YXRpdmUgV29ya3MgdGhlcmVvZiBpbiBhbnkgbWVkaXVtLCB3aXRoIG9yIHdpdGhvdXQKICAgICAgbW9kaWZpY2F0aW9ucywgYW5kIGluIFNvdXJjZSBvciBPYmplY3QgZm9ybSwgcHJvdmlkZWQgdGhhdCBZb3UKICAgICAgbWVldCB0aGUgZm9sbG93aW5nIGNvbmRpdGlvbnM6CgogICAgICAoYSkgWW91IG11c3QgZ2l2ZSBhbnkgb3RoZXIgcmVjaXBpZW50cyBvZiB0aGUgV29yayBvcgogICAgICAgICAgRGVyaXZhdGl2ZSBXb3JrcyBhIGNvcHkgb2YgdGhpcyBMaWNlbnNlOyBhbmQKCiAgICAgIChiKSBZb3UgbXVzdCBjYXVzZSBhbnkgbW9kaWZpZWQgZmlsZXMgdG8gY2FycnkgcHJvbWluZW50IG5vdGljZXMKICAgICAgICAgIHN0YXRpbmcgdGhhdCBZb3UgY2hhbmdlZCB0aGUgZmlsZXM7IGFuZAoKICAgICAgKGMpIFlvdSBtdXN0IHJldGFpbiwgaW4gdGhlIFNvdXJjZSBmb3JtIG9mIGFueSBEZXJpdmF0aXZlIFdvcmtzCiAgICAgICAgICB0aGF0IFlvdSBkaXN0cmlidXRlLCBhbGwgY29weXJpZ2h0LCBwYXRlbnQsIHRyYWRlbWFyaywgYW5kCiAgICAgICAgICBhdHRyaWJ1dGlvbiBub3RpY2VzIGZyb20gdGhlIFNvdXJjZSBmb3JtIG9mIHRoZSBXb3JrLAogICAgICAgICAgZXhjbHVkaW5nIHRob3NlIG5vdGljZXMgdGhhdCBkbyBub3QgcGVydGFpbiB0byBhbnkgcGFydCBvZgogICAgICAgICAgdGhlIERlcml2YXRpdmUgV29ya3M7IGFuZAoKICAgICAgKGQpIElmIHRoZSBXb3JrIGluY2x1ZGVzIGEgIk5PVElDRSIgdGV4dCBmaWxlIGFzIHBhcnQgb2YgaXRzCiAgICAgICAgICBkaXN0cmlidXRpb24sIHRoZW4gYW55IERlcml2YXRpdmUgV29ya3MgdGhhdCBZb3UgZGlzdHJpYnV0ZSBtdXN0CiAgICAgICAgICBpbmNsdWRlIGEgcmVhZGFibGUgY29weSBvZiB0aGUgYXR0cmlidXRpb24gbm90aWNlcyBjb250YWluZWQKICAgICAgICAgIHdpdGhpbiBzdWNoIE5PVElDRSBmaWxlLCBleGNsdWRpbmcgdGhvc2Ugbm90aWNlcyB0aGF0IGRvIG5vdAogICAgICAgICAgcGVydGFpbiB0byBhbnkgcGFydCBvZiB0aGUgRGVyaXZhdGl2ZSBXb3JrcywgaW4gYXQgbGVhc3Qgb25lCiAgICAgICAgICBvZiB0aGUgZm9sbG93aW5nIHBsYWNlczogd2l0aGluIGEgTk9USUNFIHRleHQgZmlsZSBkaXN0cmlidXRlZAogICAgICAgICAgYXMgcGFydCBvZiB0aGUgRGVyaXZhdGl2ZSBXb3Jrczsgd2l0aGluIHRoZSBTb3VyY2UgZm9ybSBvcgogICAgICAgICAgZG9jdW1lbnRhdGlvbiwgaWYgcHJvdmlkZWQgYWxvbmcgd2l0aCB0aGUgRGVyaXZhdGl2ZSBXb3Jrczsgb3IsCiAgICAgICAgICB3aXRoaW4gYSBkaXNwbGF5IGdlbmVyYXRlZCBieSB0aGUgRGVyaXZhdGl2ZSBXb3JrcywgaWYgYW5kCiAgICAgICAgICB3aGVyZXZlciBzdWNoIHRoaXJkLXBhcnR5IG5vdGljZXMgbm9ybWFsbHkgYXBwZWFyLiBUaGUgY29udGVudHMKICAgICAgICAgIG9mIHRoZSBOT1RJQ0UgZmlsZSBhcmUgZm9yIGluZm9ybWF0aW9uYWwgcHVycG9zZXMgb25seSBhbmQKICAgICAgICAgIGRvIG5vdCBtb2RpZnkgdGhlIExpY2Vuc2UuIFlvdSBtYXkgYWRkIFlvdXIgb3duIGF0dHJpYnV0aW9uCiAgICAgICAgICBub3RpY2VzIHdpdGhpbiBEZXJpdmF0aXZlIFdvcmtzIHRoYXQgWW91IGRpc3RyaWJ1dGUsIGFsb25nc2lkZQogICAgICAgICAgb3IgYXMgYW4gYWRkZW5kdW0gdG8gdGhlIE5PVElDRSB0ZXh0IGZyb20gdGhlIFdvcmssIHByb3ZpZGVkCiAgICAgICAgICB0aGF0IHN1Y2ggYWRkaXRpb25hbCBhdHRyaWJ1dGlvbiBub3RpY2VzIGNhbm5vdCBiZSBjb25zdHJ1ZWQKICAgICAgICAgIGFzIG1vZGlmeWluZyB0aGUgTGljZW5zZS4KCiAgICAgIFlvdSBtYXkgYWRkIFlvdXIgb3duIGNvcHlyaWdodCBzdGF0ZW1lbnQgdG8gWW91ciBtb2RpZmljYXRpb25zIGFuZAogICAgICBtYXkgcHJvdmlkZSBhZGRpdGlvbmFsIG9yIGRpZmZlcmVudCBsaWNlbnNlIHRlcm1zIGFuZCBjb25kaXRpb25zCiAgICAgIGZvciB1c2UsIHJlcHJvZHVjdGlvbiwgb3IgZGlzdHJpYnV0aW9uIG9mIFlvdXIgbW9kaWZpY2F0aW9ucywgb3IKICAgICAgZm9yIGFueSBzdWNoIERlcml2YXRpdmUgV29ya3MgYXMgYSB3aG9sZSwgcHJvdmlkZWQgWW91ciB1c2UsCiAgICAgIHJlcHJvZHVjdGlvbiwgYW5kIGRpc3RyaWJ1dGlvbiBvZiB0aGUgV29yayBvdGhlcndpc2UgY29tcGxpZXMgd2l0aAogICAgICB0aGUgY29uZGl0aW9ucyBzdGF0ZWQgaW4gdGhpcyBMaWNlbnNlLgoKICAgNS4gU3VibWlzc2lvbiBvZiBDb250cmlidXRpb25zLiBVbmxlc3MgWW91IGV4cGxpY2l0bHkgc3RhdGUgb3RoZXJ3aXNlLAogICAgICBhbnkgQ29udHJpYnV0aW9uIGludGVudGlvbmFsbHkgc3VibWl0dGVkIGZvciBpbmNsdXNpb24gaW4gdGhlIFdvcmsKICAgICAgYnkgWW91IHRvIHRoZSBMaWNlbnNvciBzaGFsbCBiZSB1bmRlciB0aGUgdGVybXMgYW5kIGNvbmRpdGlvbnMgb2YKICAgICAgdGhpcyBMaWNlbnNlLCB3aXRob3V0IGFueSBhZGRpdGlvbmFsIHRlcm1zIG9yIGNvbmRpdGlvbnMuCiAgICAgIE5vdHdpdGhzdGFuZGluZyB0aGUgYWJvdmUsIG5vdGhpbmcgaGVyZWluIHNoYWxsIHN1cGVyc2VkZSBvciBtb2RpZnkKICAgICAgdGhlIHRlcm1zIG9mIGFueSBzZXBhcmF0ZSBsaWNlbnNlIGFncmVlbWVudCB5b3UgbWF5IGhhdmUgZXhlY3V0ZWQKICAgICAgd2l0aCBMaWNlbnNvciByZWdhcmRpbmcgc3VjaCBDb250cmlidXRpb25zLgoKICAgNi4gVHJhZGVtYXJrcy4gVGhpcyBMaWNlbnNlIGRvZXMgbm90IGdyYW50IHBlcm1pc3Npb24gdG8gdXNlIHRoZSB0cmFkZQogICAgICBuYW1lcywgdHJhZGVtYXJrcywgc2VydmljZSBtYXJrcywgb3IgcHJvZHVjdCBuYW1lcyBvZiB0aGUgTGljZW5zb3IsCiAgICAgIGV4Y2VwdCBhcyByZXF1aXJlZCBmb3IgcmVhc29uYWJsZSBhbmQgY3VzdG9tYXJ5IHVzZSBpbiBkZXNjcmliaW5nIHRoZQogICAgICBvcmlnaW4gb2YgdGhlIFdvcmsgYW5kIHJlcHJvZHVjaW5nIHRoZSBjb250ZW50IG9mIHRoZSBOT1RJQ0UgZmlsZS4KCiAgIDcuIERpc2NsYWltZXIgb2YgV2FycmFudHkuIFVubGVzcyByZXF1aXJlZCBieSBhcHBsaWNhYmxlIGxhdyBvcgogICAgICBhZ3JlZWQgdG8gaW4gd3JpdGluZywgTGljZW5zb3IgcHJvdmlkZXMgdGhlIFdvcmsgKGFuZCBlYWNoCiAgICAgIENvbnRyaWJ1dG9yIHByb3ZpZGVzIGl0cyBDb250cmlidXRpb25zKSBvbiBhbiAiQVMgSVMiIEJBU0lTLAogICAgICBXSVRIT1VUIFdBUlJBTlRJRVMgT1IgQ09ORElUSU9OUyBPRiBBTlkgS0lORCwgZWl0aGVyIGV4cHJlc3Mgb3IKICAgICAgaW1wbGllZCwgaW5jbHVkaW5nLCB3aXRob3V0IGxpbWl0YXRpb24sIGFueSB3YXJyYW50aWVzIG9yIGNvbmRpdGlvbnMKICAgICAgb2YgVElUTEUsIE5PTi1JTkZSSU5HRU1FTlQsIE1FUkNIQU5UQUJJTElUWSwgb3IgRklUTkVTUyBGT1IgQQogICAgICBQQVJUSUNVTEFSIFBVUlBPU0UuIFlvdSBhcmUgc29sZWx5IHJlc3BvbnNpYmxlIGZvciBkZXRlcm1pbmluZyB0aGUKICAgICAgYXBwcm9wcmlhdGVuZXNzIG9mIHVzaW5nIG9yIHJlZGlzdHJpYnV0aW5nIHRoZSBXb3JrIGFuZCBhc3N1bWUgYW55CiAgICAgIHJpc2tzIGFzc29jaWF0ZWQgd2l0aCBZb3VyIGV4ZXJjaXNlIG9mIHBlcm1pc3Npb25zIHVuZGVyIHRoaXMgTGljZW5zZS4KCiAgIDguIExpbWl0YXRpb24gb2YgTGlhYmlsaXR5LiBJbiBubyBldmVudCBhbmQgdW5kZXIgbm8gbGVnYWwgdGhlb3J5LAogICAgICB3aGV0aGVyIGluIHRvcnQgKGluY2x1ZGluZyBuZWdsaWdlbmNlKSwgY29udHJhY3QsIG9yIG90aGVyd2lzZSwKICAgICAgdW5sZXNzIHJlcXVpcmVkIGJ5IGFwcGxpY2FibGUgbGF3IChzdWNoIGFzIGRlbGliZXJhdGUgYW5kIGdyb3NzbHkKICAgICAgbmVnbGlnZW50IGFjdHMpIG9yIGFncmVlZCB0byBpbiB3cml0aW5nLCBzaGFsbCBhbnkgQ29udHJpYnV0b3IgYmUKICAgICAgbGlhYmxlIHRvIFlvdSBmb3IgZGFtYWdlcywgaW5jbHVkaW5nIGFueSBkaXJlY3QsIGluZGlyZWN0LCBzcGVjaWFsLAogICAgICBpbmNpZGVudGFsLCBvciBjb25zZXF1ZW50aWFsIGRhbWFnZXMgb2YgYW55IGNoYXJhY3RlciBhcmlzaW5nIGFzIGEKICAgICAgcmVzdWx0IG9mIHRoaXMgTGljZW5zZSBvciBvdXQgb2YgdGhlIHVzZSBvciBpbmFiaWxpdHkgdG8gdXNlIHRoZQogICAgICBXb3JrIChpbmNsdWRpbmcgYnV0IG5vdCBsaW1pdGVkIHRvIGRhbWFnZXMgZm9yIGxvc3Mgb2YgZ29vZHdpbGwsCiAgICAgIHdvcmsgc3RvcHBhZ2UsIGNvbXB1dGVyIGZhaWx1cmUgb3IgbWFsZnVuY3Rpb24sIG9yIGFueSBhbmQgYWxsCiAgICAgIG90aGVyIGNvbW1lcmNpYWwgZGFtYWdlcyBvciBsb3NzZXMpLCBldmVuIGlmIHN1Y2ggQ29udHJpYnV0b3IKICAgICAgaGFzIGJlZW4gYWR2aXNlZCBvZiB0aGUgcG9zc2liaWxpdHkgb2Ygc3VjaCBkYW1hZ2VzLgoKICAgOS4gQWNjZXB0aW5nIFdhcnJhbnR5IG9yIEFkZGl0aW9uYWwgTGlhYmlsaXR5LiBXaGlsZSByZWRpc3RyaWJ1dGluZwogICAgICB0aGUgV29yayBvciBEZXJpdmF0aXZlIFdvcmtzIHRoZXJlb2YsIFlvdSBtYXkgY2hvb3NlIHRvIG9mZmVyLAogICAgICBhbmQgY2hhcmdlIGEgZmVlIGZvciwgYWNjZXB0YW5jZSBvZiBzdXBwb3J0LCB3YXJyYW50eSwgaW5kZW1uaXR5LAogICAgICBvciBvdGhlciBsaWFiaWxpdHkgb2JsaWdhdGlvbnMgYW5kL29yIHJpZ2h0cyBjb25zaXN0ZW50IHdpdGggdGhpcwogICAgICBMaWNlbnNlLiBIb3dldmVyLCBpbiBhY2NlcHRpbmcgc3VjaCBvYmxpZ2F0aW9ucywgWW91IG1heSBhY3Qgb25seQogICAgICBvbiBZb3VyIG93biBiZWhhbGYgYW5kIG9uIFlvdXIgc29sZSByZXNwb25zaWJpbGl0eSwgbm90IG9uIGJlaGFsZgogICAgICBvZiBhbnkgb3RoZXIgQ29udHJpYnV0b3IsIGFuZCBvbmx5IGlmIFlvdSBhZ3JlZSB0byBpbmRlbW5pZnksCiAgICAgIGRlZmVuZCwgYW5kIGhvbGQgZWFjaCBDb250cmlidXRvciBoYXJtbGVzcyBmb3IgYW55IGxpYWJpbGl0eQogICAgICBpbmN1cnJlZCBieSwgb3IgY2xhaW1zIGFzc2VydGVkIGFnYWluc3QsIHN1Y2ggQ29udHJpYnV0b3IgYnkgcmVhc29uCiAgICAgIG9mIHlvdXIgYWNjZXB0aW5nIGFueSBzdWNoIHdhcnJhbnR5IG9yIGFkZGl0aW9uYWwgbGlhYmlsaXR5LgoKICAgRU5EIE9GIFRFUk1TIEFORCBDT05ESVRJT05TCgogICBBUFBFTkRJWDogSG93IHRvIGFwcGx5IHRoZSBBcGFjaGUgTGljZW5zZSB0byB5b3VyIHdvcmsuCgogICAgICBUbyBhcHBseSB0aGUgQXBhY2hlIExpY2Vuc2UgdG8geW91ciB3b3JrLCBhdHRhY2ggdGhlIGZvbGxvd2luZwogICAgICBib2lsZXJwbGF0ZSBub3RpY2UsIHdpdGggdGhlIGZpZWxkcyBlbmNsb3NlZCBieSBicmFja2V0cyAiW10iCiAgICAgIHJlcGxhY2VkIHdpdGggeW91ciBvd24gaWRlbnRpZnlpbmcgaW5mb3JtYXRpb24uIChEb24ndCBpbmNsdWRlCiAgICAgIHRoZSBicmFja2V0cyEpICBUaGUgdGV4dCBzaG91bGQgYmUgZW5jbG9zZWQgaW4gdGhlIGFwcHJvcHJpYXRlCiAgICAgIGNvbW1lbnQgc3ludGF4IGZvciB0aGUgZmlsZSBmb3JtYXQuIFdlIGFsc28gcmVjb21tZW5kIHRoYXQgYQogICAgICBmaWxlIG9yIGNsYXNzIG5hbWUgYW5kIGRlc2NyaXB0aW9uIG9mIHB1cnBvc2UgYmUgaW5jbHVkZWQgb24gdGhlCiAgICAgIHNhbWUgInByaW50ZWQgcGFnZSIgYXMgdGhlIGNvcHlyaWdodCBub3RpY2UgZm9yIGVhc2llcgogICAgICBpZGVudGlmaWNhdGlvbiB3aXRoaW4gdGhpcmQtcGFydHkgYXJjaGl2ZXMuCgogICBDb3B5cmlnaHQgMjAxOSBDb3ZlbyBTb2x1dGlvbnMgSW5jLgoKICAgTGljZW5zZWQgdW5kZXIgdGhlIEFwYWNoZSBMaWNlbnNlLCBWZXJzaW9uIDIuMCAodGhlICJMaWNlbnNlIik7CiAgIHlvdSBtYXkgbm90IHVzZSB0aGlzIGZpbGUgZXhjZXB0IGluIGNvbXBsaWFuY2Ugd2l0aCB0aGUgTGljZW5zZS4KICAgWW91IG1heSBvYnRhaW4gYSBjb3B5IG9mIHRoZSBMaWNlbnNlIGF0CgogICAgICAgaHR0cDovL3d3dy5hcGFjaGUub3JnL2xpY2Vuc2VzL0xJQ0VOU0UtMi4wCgogICBVbmxlc3MgcmVxdWlyZWQgYnkgYXBwbGljYWJsZSBsYXcgb3IgYWdyZWVkIHRvIGluIHdyaXRpbmcsIHNvZnR3YXJlCiAgIGRpc3RyaWJ1dGVkIHVuZGVyIHRoZSBMaWNlbnNlIGlzIGRpc3RyaWJ1dGVkIG9uIGFuICJBUyBJUyIgQkFTSVMsCiAgIFdJVEhPVVQgV0FSUkFOVElFUyBPUiBDT05ESVRJT05TIE9GIEFOWSBLSU5ELCBlaXRoZXIgZXhwcmVzcyBvciBpbXBsaWVkLgogICBTZWUgdGhlIExpY2Vuc2UgZm9yIHRoZSBzcGVjaWZpYyBsYW5ndWFnZSBnb3Zlcm5pbmcgcGVybWlzc2lvbnMgYW5kCiAgIGxpbWl0YXRpb25zIHVuZGVyIHRoZSBMaWNlbnNlLgo="
    },
    {
      "path": "original/dist/node_modules/fdir/LICENSE",
      "source_locator": "Temp/elite-v334-a443d28e2128479b87937c040aa8b6a8/projection-rebuilt/payload/dist/node_modules/fdir/LICENSE",
      "scope": "VERBATIM existing selected payload notice; applicability remains per path",
      "origin": "ORIGINAL_SELECTED_PAYLOAD",
      "payload_path": "dist/node_modules/fdir/LICENSE",
      "bytes": 1053,
      "sha256": "9a39f2aadab11a3697edd668ff2d8ad885b649737b7ab4d3bf12b34e5ada0c86",
      "content_base64": "Q29weXJpZ2h0IDIwMjMgQWJkdWxsYWggQXR0YQoKUGVybWlzc2lvbiBpcyBoZXJlYnkgZ3JhbnRlZCwgZnJlZSBvZiBjaGFyZ2UsIHRvIGFueSBwZXJzb24gb2J0YWluaW5nIGEgY29weSBvZiB0aGlzIHNvZnR3YXJlIGFuZCBhc3NvY2lhdGVkIGRvY3VtZW50YXRpb24gZmlsZXMgKHRoZSAiU29mdHdhcmUiKSwgdG8gZGVhbCBpbiB0aGUgU29mdHdhcmUgd2l0aG91dCByZXN0cmljdGlvbiwgaW5jbHVkaW5nIHdpdGhvdXQgbGltaXRhdGlvbiB0aGUgcmlnaHRzIHRvIHVzZSwgY29weSwgbW9kaWZ5LCBtZXJnZSwgcHVibGlzaCwgZGlzdHJpYnV0ZSwgc3VibGljZW5zZSwgYW5kL29yIHNlbGwgY29waWVzIG9mIHRoZSBTb2Z0d2FyZSwgYW5kIHRvIHBlcm1pdCBwZXJzb25zIHRvIHdob20gdGhlIFNvZnR3YXJlIGlzIGZ1cm5pc2hlZCB0byBkbyBzbywgc3ViamVjdCB0byB0aGUgZm9sbG93aW5nIGNvbmRpdGlvbnM6CgpUaGUgYWJvdmUgY29weXJpZ2h0IG5vdGljZSBhbmQgdGhpcyBwZXJtaXNzaW9uIG5vdGljZSBzaGFsbCBiZSBpbmNsdWRlZCBpbiBhbGwgY29waWVzIG9yIHN1YnN0YW50aWFsIHBvcnRpb25zIG9mIHRoZSBTb2Z0d2FyZS4KClRIRSBTT0ZUV0FSRSBJUyBQUk9WSURFRCAiQVMgSVMiLCBXSVRIT1VUIFdBUlJBTlRZIE9GIEFOWSBLSU5ELCBFWFBSRVNTIE9SIElNUExJRUQsIElOQ0xVRElORyBCVVQgTk9UIExJTUlURUQgVE8gVEhFIFdBUlJBTlRJRVMgT0YgTUVSQ0hBTlRBQklMSVRZLCBGSVRORVNTIEZPUiBBIFBBUlRJQ1VMQVIgUFVSUE9TRSBBTkQgTk9OSU5GUklOR0VNRU5ULiBJTiBOTyBFVkVOVCBTSEFMTCBUSEUgQVVUSE9SUyBPUiBDT1BZUklHSFQgSE9MREVSUyBCRSBMSUFCTEUgRk9SIEFOWSBDTEFJTSwgREFNQUdFUyBPUiBPVEhFUiBMSUFCSUxJVFksIFdIRVRIRVIgSU4gQU4gQUNUSU9OIE9GIENPTlRSQUNULCBUT1JUIE9SIE9USEVSV0lTRSwgQVJJU0lORyBGUk9NLCBPVVQgT0YgT1IgSU4gQ09OTkVDVElPTiBXSVRIIFRIRSBTT0ZUV0FSRSBPUiBUSEUgVVNFIE9SIE9USEVSIERFQUxJTkdTIElOIFRIRSBTT0ZUV0FSRS4K"
    },
    {
      "path": "original/dist/node_modules/graceful-fs/LICENSE",
      "source_locator": "Temp/elite-v334-a443d28e2128479b87937c040aa8b6a8/projection-rebuilt/payload/dist/node_modules/graceful-fs/LICENSE",
      "scope": "VERBATIM existing selected payload notice; applicability remains per path",
      "origin": "ORIGINAL_SELECTED_PAYLOAD",
      "payload_path": "dist/node_modules/graceful-fs/LICENSE",
      "bytes": 791,
      "sha256": "f65c5d9f22a317b2a10803bd1868461ce6499c2ed7217bc80c0cc772a748789c",
      "content_base64": "VGhlIElTQyBMaWNlbnNlCgpDb3B5cmlnaHQgKGMpIDIwMTEtMjAyMiBJc2FhYyBaLiBTY2hsdWV0ZXIsIEJlbiBOb29yZGh1aXMsIGFuZCBDb250cmlidXRvcnMKClBlcm1pc3Npb24gdG8gdXNlLCBjb3B5LCBtb2RpZnksIGFuZC9vciBkaXN0cmlidXRlIHRoaXMgc29mdHdhcmUgZm9yIGFueQpwdXJwb3NlIHdpdGggb3Igd2l0aG91dCBmZWUgaXMgaGVyZWJ5IGdyYW50ZWQsIHByb3ZpZGVkIHRoYXQgdGhlIGFib3ZlCmNvcHlyaWdodCBub3RpY2UgYW5kIHRoaXMgcGVybWlzc2lvbiBub3RpY2UgYXBwZWFyIGluIGFsbCBjb3BpZXMuCgpUSEUgU09GVFdBUkUgSVMgUFJPVklERUQgIkFTIElTIiBBTkQgVEhFIEFVVEhPUiBESVNDTEFJTVMgQUxMIFdBUlJBTlRJRVMKV0lUSCBSRUdBUkQgVE8gVEhJUyBTT0ZUV0FSRSBJTkNMVURJTkcgQUxMIElNUExJRUQgV0FSUkFOVElFUyBPRgpNRVJDSEFOVEFCSUxJVFkgQU5EIEZJVE5FU1MuIElOIE5PIEVWRU5UIFNIQUxMIFRIRSBBVVRIT1IgQkUgTElBQkxFIEZPUgpBTlkgU1BFQ0lBTCwgRElSRUNULCBJTkRJUkVDVCwgT1IgQ09OU0VRVUVOVElBTCBEQU1BR0VTIE9SIEFOWSBEQU1BR0VTCldIQVRTT0VWRVIgUkVTVUxUSU5HIEZST00gTE9TUyBPRiBVU0UsIERBVEEgT1IgUFJPRklUUywgV0hFVEhFUiBJTiBBTgpBQ1RJT04gT0YgQ09OVFJBQ1QsIE5FR0xJR0VOQ0UgT1IgT1RIRVIgVE9SVElPVVMgQUNUSU9OLCBBUklTSU5HIE9VVCBPRiBPUgpJTiBDT05ORUNUSU9OIFdJVEggVEhFIFVTRSBPUiBQRVJGT1JNQU5DRSBPRiBUSElTIFNPRlRXQVJFLgo="
    },
    {
      "path": "original/dist/node_modules/minizlib/LICENSE",
      "source_locator": "Temp/elite-v334-a443d28e2128479b87937c040aa8b6a8/projection-rebuilt/payload/dist/node_modules/minizlib/LICENSE",
      "scope": "VERBATIM existing selected payload notice; applicability remains per path",
      "origin": "ORIGINAL_SELECTED_PAYLOAD",
      "payload_path": "dist/node_modules/minizlib/LICENSE",
      "bytes": 1339,
      "sha256": "df05697a04b8997c79d423930534e59cc18828f6e22d0c8c2d62fcbea1c4f50b",
      "content_base64": "TWluaXpsaWIgd2FzIGNyZWF0ZWQgYnkgSXNhYWMgWi4gU2NobHVldGVyLgpJdCBpcyBhIGRlcml2YXRpdmUgd29yayBvZiB0aGUgTm9kZS5qcyBwcm9qZWN0LgoKIiIiCkNvcHlyaWdodCAoYykgMjAxNy0yMDIzIElzYWFjIFouIFNjaGx1ZXRlciBhbmQgQ29udHJpYnV0b3JzCkNvcHlyaWdodCAoYykgMjAxNy0yMDIzIE5vZGUuanMgY29udHJpYnV0b3JzLiBBbGwgcmlnaHRzIHJlc2VydmVkLgpDb3B5cmlnaHQgKGMpIDIwMTctMjAyMyBKb3llbnQsIEluYy4gYW5kIG90aGVyIE5vZGUgY29udHJpYnV0b3JzLiBBbGwgcmlnaHRzIHJlc2VydmVkLgoKUGVybWlzc2lvbiBpcyBoZXJlYnkgZ3JhbnRlZCwgZnJlZSBvZiBjaGFyZ2UsIHRvIGFueSBwZXJzb24gb2J0YWluaW5nIGEKY29weSBvZiB0aGlzIHNvZnR3YXJlIGFuZCBhc3NvY2lhdGVkIGRvY3VtZW50YXRpb24gZmlsZXMgKHRoZSAiU29mdHdhcmUiKSwKdG8gZGVhbCBpbiB0aGUgU29mdHdhcmUgd2l0aG91dCByZXN0cmljdGlvbiwgaW5jbHVkaW5nIHdpdGhvdXQgbGltaXRhdGlvbgp0aGUgcmlnaHRzIHRvIHVzZSwgY29weSwgbW9kaWZ5LCBtZXJnZSwgcHVibGlzaCwgZGlzdHJpYnV0ZSwgc3VibGljZW5zZSwKYW5kL29yIHNlbGwgY29waWVzIG9mIHRoZSBTb2Z0d2FyZSwgYW5kIHRvIHBlcm1pdCBwZXJzb25zIHRvIHdob20gdGhlClNvZnR3YXJlIGlzIGZ1cm5pc2hlZCB0byBkbyBzbywgc3ViamVjdCB0byB0aGUgZm9sbG93aW5nIGNvbmRpdGlvbnM6CgpUaGUgYWJvdmUgY29weXJpZ2h0IG5vdGljZSBhbmQgdGhpcyBwZXJtaXNzaW9uIG5vdGljZSBzaGFsbCBiZSBpbmNsdWRlZCBpbgphbGwgY29waWVzIG9yIHN1YnN0YW50aWFsIHBvcnRpb25zIG9mIHRoZSBTb2Z0d2FyZS4KClRIRSBTT0ZUV0FSRSBJUyBQUk9WSURFRCAiQVMgSVMiLCBXSVRIT1VUIFdBUlJBTlRZIE9GIEFOWSBLSU5ELCBFWFBSRVNTCk9SIElNUExJRUQsIElOQ0xVRElORyBCVVQgTk9UIExJTUlURUQgVE8gVEhFIFdBUlJBTlRJRVMgT0YKTUVSQ0hBTlRBQklMSVRZLCBGSVRORVNTIEZPUiBBIFBBUlRJQ1VMQVIgUFVSUE9TRSBBTkQgTk9OSU5GUklOR0VNRU5ULgpJTiBOTyBFVkVOVCBTSEFMTCBUSEUgQVVUSE9SUyBPUiBDT1BZUklHSFQgSE9MREVSUyBCRSBMSUFCTEUgRk9SIEFOWQpDTEFJTSwgREFNQUdFUyBPUiBPVEhFUiBMSUFCSUxJVFksIFdIRVRIRVIgSU4gQU4gQUNUSU9OIE9GIENPTlRSQUNULApUT1JUIE9SIE9USEVSV0lTRSwgQVJJU0lORyBGUk9NLCBPVVQgT0YgT1IgSU4gQ09OTkVDVElPTiBXSVRIIFRIRQpTT0ZUV0FSRSBPUiBUSEUgVVNFIE9SIE9USEVSIERFQUxJTkdTIElOIFRIRSBTT0ZUV0FSRS4KIiIiCg=="
    },
    {
      "path": "original/dist/node_modules/node-gyp/gyp/LICENSE",
      "source_locator": "Temp/elite-v334-a443d28e2128479b87937c040aa8b6a8/projection-rebuilt/payload/dist/node_modules/node-gyp/gyp/LICENSE",
      "scope": "VERBATIM existing selected payload notice; applicability remains per path",
      "origin": "ORIGINAL_SELECTED_PAYLOAD",
      "payload_path": "dist/node_modules/node-gyp/gyp/LICENSE",
      "bytes": 1537,
      "sha256": "ca90abb6ed71de0774461ef9f928de33e748b617aeb79f9e52415cf08d69230e",
      "content_base64": "Q29weXJpZ2h0IChjKSAyMDIwIE5vZGUuanMgY29udHJpYnV0b3JzLiBBbGwgcmlnaHRzIHJlc2VydmVkLgpDb3B5cmlnaHQgKGMpIDIwMDkgR29vZ2xlIEluYy4gQWxsIHJpZ2h0cyByZXNlcnZlZC4KClJlZGlzdHJpYnV0aW9uIGFuZCB1c2UgaW4gc291cmNlIGFuZCBiaW5hcnkgZm9ybXMsIHdpdGggb3Igd2l0aG91dAptb2RpZmljYXRpb24sIGFyZSBwZXJtaXR0ZWQgcHJvdmlkZWQgdGhhdCB0aGUgZm9sbG93aW5nIGNvbmRpdGlvbnMgYXJlCm1ldDoKCiAgICogUmVkaXN0cmlidXRpb25zIG9mIHNvdXJjZSBjb2RlIG11c3QgcmV0YWluIHRoZSBhYm92ZSBjb3B5cmlnaHQKbm90aWNlLCB0aGlzIGxpc3Qgb2YgY29uZGl0aW9ucyBhbmQgdGhlIGZvbGxvd2luZyBkaXNjbGFpbWVyLgogICAqIFJlZGlzdHJpYnV0aW9ucyBpbiBiaW5hcnkgZm9ybSBtdXN0IHJlcHJvZHVjZSB0aGUgYWJvdmUKY29weXJpZ2h0IG5vdGljZSwgdGhpcyBsaXN0IG9mIGNvbmRpdGlvbnMgYW5kIHRoZSBmb2xsb3dpbmcgZGlzY2xhaW1lcgppbiB0aGUgZG9jdW1lbnRhdGlvbiBhbmQvb3Igb3RoZXIgbWF0ZXJpYWxzIHByb3ZpZGVkIHdpdGggdGhlCmRpc3RyaWJ1dGlvbi4KICAgKiBOZWl0aGVyIHRoZSBuYW1lIG9mIEdvb2dsZSBJbmMuIG5vciB0aGUgbmFtZXMgb2YgaXRzCmNvbnRyaWJ1dG9ycyBtYXkgYmUgdXNlZCB0byBlbmRvcnNlIG9yIHByb21vdGUgcHJvZHVjdHMgZGVyaXZlZCBmcm9tCnRoaXMgc29mdHdhcmUgd2l0aG91dCBzcGVjaWZpYyBwcmlvciB3cml0dGVuIHBlcm1pc3Npb24uCgpUSElTIFNPRlRXQVJFIElTIFBST1ZJREVEIEJZIFRIRSBDT1BZUklHSFQgSE9MREVSUyBBTkQgQ09OVFJJQlVUT1JTCiJBUyBJUyIgQU5EIEFOWSBFWFBSRVNTIE9SIElNUExJRUQgV0FSUkFOVElFUywgSU5DTFVESU5HLCBCVVQgTk9UCkxJTUlURUQgVE8sIFRIRSBJTVBMSUVEIFdBUlJBTlRJRVMgT0YgTUVSQ0hBTlRBQklMSVRZIEFORCBGSVRORVNTIEZPUgpBIFBBUlRJQ1VMQVIgUFVSUE9TRSBBUkUgRElTQ0xBSU1FRC4gSU4gTk8gRVZFTlQgU0hBTEwgVEhFIENPUFlSSUdIVApPV05FUiBPUiBDT05UUklCVVRPUlMgQkUgTElBQkxFIEZPUiBBTlkgRElSRUNULCBJTkRJUkVDVCwgSU5DSURFTlRBTCwKU1BFQ0lBTCwgRVhFTVBMQVJZLCBPUiBDT05TRVFVRU5USUFMIERBTUFHRVMgKElOQ0xVRElORywgQlVUIE5PVApMSU1JVEVEIFRPLCBQUk9DVVJFTUVOVCBPRiBTVUJTVElUVVRFIEdPT0RTIE9SIFNFUlZJQ0VTOyBMT1NTIE9GIFVTRSwKREFUQSwgT1IgUFJPRklUUzsgT1IgQlVTSU5FU1MgSU5URVJSVVBUSU9OKSBIT1dFVkVSIENBVVNFRCBBTkQgT04gQU5ZClRIRU9SWSBPRiBMSUFCSUxJVFksIFdIRVRIRVIgSU4gQ09OVFJBQ1QsIFNUUklDVCBMSUFCSUxJVFksIE9SIFRPUlQKKElOQ0xVRElORyBORUdMSUdFTkNFIE9SIE9USEVSV0lTRSkgQVJJU0lORyBJTiBBTlkgV0FZIE9VVCBPRiBUSEUgVVNFCk9GIFRISVMgU09GVFdBUkUsIEVWRU4gSUYgQURWSVNFRCBPRiBUSEUgUE9TU0lCSUxJVFkgT0YgU1VDSCBEQU1BR0UuCg=="
    },
    {
      "path": "original/dist/node_modules/node-gyp/gyp/pylib/packaging/LICENSE",
      "source_locator": "Temp/elite-v334-a443d28e2128479b87937c040aa8b6a8/projection-rebuilt/payload/dist/node_modules/node-gyp/gyp/pylib/packaging/LICENSE",
      "scope": "VERBATIM existing selected payload notice; applicability remains per path",
      "origin": "ORIGINAL_SELECTED_PAYLOAD",
      "payload_path": "dist/node_modules/node-gyp/gyp/pylib/packaging/LICENSE",
      "bytes": 197,
      "sha256": "cad1ef5bd340d73e074ba614d26f7deaca5c7940c3d8c34852e65c4909686c48",
      "content_base64": "VGhpcyBzb2Z0d2FyZSBpcyBtYWRlIGF2YWlsYWJsZSB1bmRlciB0aGUgdGVybXMgb2YgKmVpdGhlciogb2YgdGhlIGxpY2Vuc2VzCmZvdW5kIGluIExJQ0VOU0UuQVBBQ0hFIG9yIExJQ0VOU0UuQlNELiBDb250cmlidXRpb25zIHRvIHRoaXMgc29mdHdhcmUgaXMgbWFkZQp1bmRlciB0aGUgdGVybXMgb2YgKmJvdGgqIHRoZXNlIGxpY2Vuc2VzLgo="
    },
    {
      "path": "original/dist/node_modules/node-gyp/gyp/pylib/packaging/LICENSE.APACHE",
      "source_locator": "Temp/elite-v334-a443d28e2128479b87937c040aa8b6a8/projection-rebuilt/payload/dist/node_modules/node-gyp/gyp/pylib/packaging/LICENSE.APACHE",
      "scope": "VERBATIM existing selected payload notice; applicability remains per path",
      "origin": "ORIGINAL_SELECTED_PAYLOAD",
      "payload_path": "dist/node_modules/node-gyp/gyp/pylib/packaging/LICENSE.APACHE",
      "bytes": 10174,
      "sha256": "0d542e0c8804e39aa7f37eb00da5a762149dc682d7829451287e11b938e94594",
      "content_base64": "CiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIEFwYWNoZSBMaWNlbnNlCiAgICAgICAgICAgICAgICAgICAgICAgICAgIFZlcnNpb24gMi4wLCBKYW51YXJ5IDIwMDQKICAgICAgICAgICAgICAgICAgICAgICAgaHR0cDovL3d3dy5hcGFjaGUub3JnL2xpY2Vuc2VzLwoKICAgVEVSTVMgQU5EIENPTkRJVElPTlMgRk9SIFVTRSwgUkVQUk9EVUNUSU9OLCBBTkQgRElTVFJJQlVUSU9OCgogICAxLiBEZWZpbml0aW9ucy4KCiAgICAgICJMaWNlbnNlIiBzaGFsbCBtZWFuIHRoZSB0ZXJtcyBhbmQgY29uZGl0aW9ucyBmb3IgdXNlLCByZXByb2R1Y3Rpb24sCiAgICAgIGFuZCBkaXN0cmlidXRpb24gYXMgZGVmaW5lZCBieSBTZWN0aW9ucyAxIHRocm91Z2ggOSBvZiB0aGlzIGRvY3VtZW50LgoKICAgICAgIkxpY2Vuc29yIiBzaGFsbCBtZWFuIHRoZSBjb3B5cmlnaHQgb3duZXIgb3IgZW50aXR5IGF1dGhvcml6ZWQgYnkKICAgICAgdGhlIGNvcHlyaWdodCBvd25lciB0aGF0IGlzIGdyYW50aW5nIHRoZSBMaWNlbnNlLgoKICAgICAgIkxlZ2FsIEVudGl0eSIgc2hhbGwgbWVhbiB0aGUgdW5pb24gb2YgdGhlIGFjdGluZyBlbnRpdHkgYW5kIGFsbAogICAgICBvdGhlciBlbnRpdGllcyB0aGF0IGNvbnRyb2wsIGFyZSBjb250cm9sbGVkIGJ5LCBvciBhcmUgdW5kZXIgY29tbW9uCiAgICAgIGNvbnRyb2wgd2l0aCB0aGF0IGVudGl0eS4gRm9yIHRoZSBwdXJwb3NlcyBvZiB0aGlzIGRlZmluaXRpb24sCiAgICAgICJjb250cm9sIiBtZWFucyAoaSkgdGhlIHBvd2VyLCBkaXJlY3Qgb3IgaW5kaXJlY3QsIHRvIGNhdXNlIHRoZQogICAgICBkaXJlY3Rpb24gb3IgbWFuYWdlbWVudCBvZiBzdWNoIGVudGl0eSwgd2hldGhlciBieSBjb250cmFjdCBvcgogICAgICBvdGhlcndpc2UsIG9yIChpaSkgb3duZXJzaGlwIG9mIGZpZnR5IHBlcmNlbnQgKDUwJSkgb3IgbW9yZSBvZiB0aGUKICAgICAgb3V0c3RhbmRpbmcgc2hhcmVzLCBvciAoaWlpKSBiZW5lZmljaWFsIG93bmVyc2hpcCBvZiBzdWNoIGVudGl0eS4KCiAgICAgICJZb3UiIChvciAiWW91ciIpIHNoYWxsIG1lYW4gYW4gaW5kaXZpZHVhbCBvciBMZWdhbCBFbnRpdHkKICAgICAgZXhlcmNpc2luZyBwZXJtaXNzaW9ucyBncmFudGVkIGJ5IHRoaXMgTGljZW5zZS4KCiAgICAgICJTb3VyY2UiIGZvcm0gc2hhbGwgbWVhbiB0aGUgcHJlZmVycmVkIGZvcm0gZm9yIG1ha2luZyBtb2RpZmljYXRpb25zLAogICAgICBpbmNsdWRpbmcgYnV0IG5vdCBsaW1pdGVkIHRvIHNvZnR3YXJlIHNvdXJjZSBjb2RlLCBkb2N1bWVudGF0aW9uCiAgICAgIHNvdXJjZSwgYW5kIGNvbmZpZ3VyYXRpb24gZmlsZXMuCgogICAgICAiT2JqZWN0IiBmb3JtIHNoYWxsIG1lYW4gYW55IGZvcm0gcmVzdWx0aW5nIGZyb20gbWVjaGFuaWNhbAogICAgICB0cmFuc2Zvcm1hdGlvbiBvciB0cmFuc2xhdGlvbiBvZiBhIFNvdXJjZSBmb3JtLCBpbmNsdWRpbmcgYnV0CiAgICAgIG5vdCBsaW1pdGVkIHRvIGNvbXBpbGVkIG9iamVjdCBjb2RlLCBnZW5lcmF0ZWQgZG9jdW1lbnRhdGlvbiwKICAgICAgYW5kIGNvbnZlcnNpb25zIHRvIG90aGVyIG1lZGlhIHR5cGVzLgoKICAgICAgIldvcmsiIHNoYWxsIG1lYW4gdGhlIHdvcmsgb2YgYXV0aG9yc2hpcCwgd2hldGhlciBpbiBTb3VyY2Ugb3IKICAgICAgT2JqZWN0IGZvcm0sIG1hZGUgYXZhaWxhYmxlIHVuZGVyIHRoZSBMaWNlbnNlLCBhcyBpbmRpY2F0ZWQgYnkgYQogICAgICBjb3B5cmlnaHQgbm90aWNlIHRoYXQgaXMgaW5jbHVkZWQgaW4gb3IgYXR0YWNoZWQgdG8gdGhlIHdvcmsKICAgICAgKGFuIGV4YW1wbGUgaXMgcHJvdmlkZWQgaW4gdGhlIEFwcGVuZGl4IGJlbG93KS4KCiAgICAgICJEZXJpdmF0aXZlIFdvcmtzIiBzaGFsbCBtZWFuIGFueSB3b3JrLCB3aGV0aGVyIGluIFNvdXJjZSBvciBPYmplY3QKICAgICAgZm9ybSwgdGhhdCBpcyBiYXNlZCBvbiAob3IgZGVyaXZlZCBmcm9tKSB0aGUgV29yayBhbmQgZm9yIHdoaWNoIHRoZQogICAgICBlZGl0b3JpYWwgcmV2aXNpb25zLCBhbm5vdGF0aW9ucywgZWxhYm9yYXRpb25zLCBvciBvdGhlciBtb2RpZmljYXRpb25zCiAgICAgIHJlcHJlc2VudCwgYXMgYSB3aG9sZSwgYW4gb3JpZ2luYWwgd29yayBvZiBhdXRob3JzaGlwLiBGb3IgdGhlIHB1cnBvc2VzCiAgICAgIG9mIHRoaXMgTGljZW5zZSwgRGVyaXZhdGl2ZSBXb3JrcyBzaGFsbCBub3QgaW5jbHVkZSB3b3JrcyB0aGF0IHJlbWFpbgogICAgICBzZXBhcmFibGUgZnJvbSwgb3IgbWVyZWx5IGxpbmsgKG9yIGJpbmQgYnkgbmFtZSkgdG8gdGhlIGludGVyZmFjZXMgb2YsCiAgICAgIHRoZSBXb3JrIGFuZCBEZXJpdmF0aXZlIFdvcmtzIHRoZXJlb2YuCgogICAgICAiQ29udHJpYnV0aW9uIiBzaGFsbCBtZWFuIGFueSB3b3JrIG9mIGF1dGhvcnNoaXAsIGluY2x1ZGluZwogICAgICB0aGUgb3JpZ2luYWwgdmVyc2lvbiBvZiB0aGUgV29yayBhbmQgYW55IG1vZGlmaWNhdGlvbnMgb3IgYWRkaXRpb25zCiAgICAgIHRvIHRoYXQgV29yayBvciBEZXJpdmF0aXZlIFdvcmtzIHRoZXJlb2YsIHRoYXQgaXMgaW50ZW50aW9uYWxseQogICAgICBzdWJtaXR0ZWQgdG8gTGljZW5zb3IgZm9yIGluY2x1c2lvbiBpbiB0aGUgV29yayBieSB0aGUgY29weXJpZ2h0IG93bmVyCiAgICAgIG9yIGJ5IGFuIGluZGl2aWR1YWwgb3IgTGVnYWwgRW50aXR5IGF1dGhvcml6ZWQgdG8gc3VibWl0IG9uIGJlaGFsZiBvZgogICAgICB0aGUgY29weXJpZ2h0IG93bmVyLiBGb3IgdGhlIHB1cnBvc2VzIG9mIHRoaXMgZGVmaW5pdGlvbiwgInN1Ym1pdHRlZCIKICAgICAgbWVhbnMgYW55IGZvcm0gb2YgZWxlY3Ryb25pYywgdmVyYmFsLCBvciB3cml0dGVuIGNvbW11bmljYXRpb24gc2VudAogICAgICB0byB0aGUgTGljZW5zb3Igb3IgaXRzIHJlcHJlc2VudGF0aXZlcywgaW5jbHVkaW5nIGJ1dCBub3QgbGltaXRlZCB0bwogICAgICBjb21tdW5pY2F0aW9uIG9uIGVsZWN0cm9uaWMgbWFpbGluZyBsaXN0cywgc291cmNlIGNvZGUgY29udHJvbCBzeXN0ZW1zLAogICAgICBhbmQgaXNzdWUgdHJhY2tpbmcgc3lzdGVtcyB0aGF0IGFyZSBtYW5hZ2VkIGJ5LCBvciBvbiBiZWhhbGYgb2YsIHRoZQogICAgICBMaWNlbnNvciBmb3IgdGhlIHB1cnBvc2Ugb2YgZGlzY3Vzc2luZyBhbmQgaW1wcm92aW5nIHRoZSBXb3JrLCBidXQKICAgICAgZXhjbHVkaW5nIGNvbW11bmljYXRpb24gdGhhdCBpcyBjb25zcGljdW91c2x5IG1hcmtlZCBvciBvdGhlcndpc2UKICAgICAgZGVzaWduYXRlZCBpbiB3cml0aW5nIGJ5IHRoZSBjb3B5cmlnaHQgb3duZXIgYXMgIk5vdCBhIENvbnRyaWJ1dGlvbi4iCgogICAgICAiQ29udHJpYnV0b3IiIHNoYWxsIG1lYW4gTGljZW5zb3IgYW5kIGFueSBpbmRpdmlkdWFsIG9yIExlZ2FsIEVudGl0eQogICAgICBvbiBiZWhhbGYgb2Ygd2hvbSBhIENvbnRyaWJ1dGlvbiBoYXMgYmVlbiByZWNlaXZlZCBieSBMaWNlbnNvciBhbmQKICAgICAgc3Vic2VxdWVudGx5IGluY29ycG9yYXRlZCB3aXRoaW4gdGhlIFdvcmsuCgogICAyLiBHcmFudCBvZiBDb3B5cmlnaHQgTGljZW5zZS4gU3ViamVjdCB0byB0aGUgdGVybXMgYW5kIGNvbmRpdGlvbnMgb2YKICAgICAgdGhpcyBMaWNlbnNlLCBlYWNoIENvbnRyaWJ1dG9yIGhlcmVieSBncmFudHMgdG8gWW91IGEgcGVycGV0dWFsLAogICAgICB3b3JsZHdpZGUsIG5vbi1leGNsdXNpdmUsIG5vLWNoYXJnZSwgcm95YWx0eS1mcmVlLCBpcnJldm9jYWJsZQogICAgICBjb3B5cmlnaHQgbGljZW5zZSB0byByZXByb2R1Y2UsIHByZXBhcmUgRGVyaXZhdGl2ZSBXb3JrcyBvZiwKICAgICAgcHVibGljbHkgZGlzcGxheSwgcHVibGljbHkgcGVyZm9ybSwgc3VibGljZW5zZSwgYW5kIGRpc3RyaWJ1dGUgdGhlCiAgICAgIFdvcmsgYW5kIHN1Y2ggRGVyaXZhdGl2ZSBXb3JrcyBpbiBTb3VyY2Ugb3IgT2JqZWN0IGZvcm0uCgogICAzLiBHcmFudCBvZiBQYXRlbnQgTGljZW5zZS4gU3ViamVjdCB0byB0aGUgdGVybXMgYW5kIGNvbmRpdGlvbnMgb2YKICAgICAgdGhpcyBMaWNlbnNlLCBlYWNoIENvbnRyaWJ1dG9yIGhlcmVieSBncmFudHMgdG8gWW91IGEgcGVycGV0dWFsLAogICAgICB3b3JsZHdpZGUsIG5vbi1leGNsdXNpdmUsIG5vLWNoYXJnZSwgcm95YWx0eS1mcmVlLCBpcnJldm9jYWJsZQogICAgICAoZXhjZXB0IGFzIHN0YXRlZCBpbiB0aGlzIHNlY3Rpb24pIHBhdGVudCBsaWNlbnNlIHRvIG1ha2UsIGhhdmUgbWFkZSwKICAgICAgdXNlLCBvZmZlciB0byBzZWxsLCBzZWxsLCBpbXBvcnQsIGFuZCBvdGhlcndpc2UgdHJhbnNmZXIgdGhlIFdvcmssCiAgICAgIHdoZXJlIHN1Y2ggbGljZW5zZSBhcHBsaWVzIG9ubHkgdG8gdGhvc2UgcGF0ZW50IGNsYWltcyBsaWNlbnNhYmxlCiAgICAgIGJ5IHN1Y2ggQ29udHJpYnV0b3IgdGhhdCBhcmUgbmVjZXNzYXJpbHkgaW5mcmluZ2VkIGJ5IHRoZWlyCiAgICAgIENvbnRyaWJ1dGlvbihzKSBhbG9uZSBvciBieSBjb21iaW5hdGlvbiBvZiB0aGVpciBDb250cmlidXRpb24ocykKICAgICAgd2l0aCB0aGUgV29yayB0byB3aGljaCBzdWNoIENvbnRyaWJ1dGlvbihzKSB3YXMgc3VibWl0dGVkLiBJZiBZb3UKICAgICAgaW5zdGl0dXRlIHBhdGVudCBsaXRpZ2F0aW9uIGFnYWluc3QgYW55IGVudGl0eSAoaW5jbHVkaW5nIGEKICAgICAgY3Jvc3MtY2xhaW0gb3IgY291bnRlcmNsYWltIGluIGEgbGF3c3VpdCkgYWxsZWdpbmcgdGhhdCB0aGUgV29yawogICAgICBvciBhIENvbnRyaWJ1dGlvbiBpbmNvcnBvcmF0ZWQgd2l0aGluIHRoZSBXb3JrIGNvbnN0aXR1dGVzIGRpcmVjdAogICAgICBvciBjb250cmlidXRvcnkgcGF0ZW50IGluZnJpbmdlbWVudCwgdGhlbiBhbnkgcGF0ZW50IGxpY2Vuc2VzCiAgICAgIGdyYW50ZWQgdG8gWW91IHVuZGVyIHRoaXMgTGljZW5zZSBmb3IgdGhhdCBXb3JrIHNoYWxsIHRlcm1pbmF0ZQogICAgICBhcyBvZiB0aGUgZGF0ZSBzdWNoIGxpdGlnYXRpb24gaXMgZmlsZWQuCgogICA0LiBSZWRpc3RyaWJ1dGlvbi4gWW91IG1heSByZXByb2R1Y2UgYW5kIGRpc3RyaWJ1dGUgY29waWVzIG9mIHRoZQogICAgICBXb3JrIG9yIERlcml2YXRpdmUgV29ya3MgdGhlcmVvZiBpbiBhbnkgbWVkaXVtLCB3aXRoIG9yIHdpdGhvdXQKICAgICAgbW9kaWZpY2F0aW9ucywgYW5kIGluIFNvdXJjZSBvciBPYmplY3QgZm9ybSwgcHJvdmlkZWQgdGhhdCBZb3UKICAgICAgbWVldCB0aGUgZm9sbG93aW5nIGNvbmRpdGlvbnM6CgogICAgICAoYSkgWW91IG11c3QgZ2l2ZSBhbnkgb3RoZXIgcmVjaXBpZW50cyBvZiB0aGUgV29yayBvcgogICAgICAgICAgRGVyaXZhdGl2ZSBXb3JrcyBhIGNvcHkgb2YgdGhpcyBMaWNlbnNlOyBhbmQKCiAgICAgIChiKSBZb3UgbXVzdCBjYXVzZSBhbnkgbW9kaWZpZWQgZmlsZXMgdG8gY2FycnkgcHJvbWluZW50IG5vdGljZXMKICAgICAgICAgIHN0YXRpbmcgdGhhdCBZb3UgY2hhbmdlZCB0aGUgZmlsZXM7IGFuZAoKICAgICAgKGMpIFlvdSBtdXN0IHJldGFpbiwgaW4gdGhlIFNvdXJjZSBmb3JtIG9mIGFueSBEZXJpdmF0aXZlIFdvcmtzCiAgICAgICAgICB0aGF0IFlvdSBkaXN0cmlidXRlLCBhbGwgY29weXJpZ2h0LCBwYXRlbnQsIHRyYWRlbWFyaywgYW5kCiAgICAgICAgICBhdHRyaWJ1dGlvbiBub3RpY2VzIGZyb20gdGhlIFNvdXJjZSBmb3JtIG9mIHRoZSBXb3JrLAogICAgICAgICAgZXhjbHVkaW5nIHRob3NlIG5vdGljZXMgdGhhdCBkbyBub3QgcGVydGFpbiB0byBhbnkgcGFydCBvZgogICAgICAgICAgdGhlIERlcml2YXRpdmUgV29ya3M7IGFuZAoKICAgICAgKGQpIElmIHRoZSBXb3JrIGluY2x1ZGVzIGEgIk5PVElDRSIgdGV4dCBmaWxlIGFzIHBhcnQgb2YgaXRzCiAgICAgICAgICBkaXN0cmlidXRpb24sIHRoZW4gYW55IERlcml2YXRpdmUgV29ya3MgdGhhdCBZb3UgZGlzdHJpYnV0ZSBtdXN0CiAgICAgICAgICBpbmNsdWRlIGEgcmVhZGFibGUgY29weSBvZiB0aGUgYXR0cmlidXRpb24gbm90aWNlcyBjb250YWluZWQKICAgICAgICAgIHdpdGhpbiBzdWNoIE5PVElDRSBmaWxlLCBleGNsdWRpbmcgdGhvc2Ugbm90aWNlcyB0aGF0IGRvIG5vdAogICAgICAgICAgcGVydGFpbiB0byBhbnkgcGFydCBvZiB0aGUgRGVyaXZhdGl2ZSBXb3JrcywgaW4gYXQgbGVhc3Qgb25lCiAgICAgICAgICBvZiB0aGUgZm9sbG93aW5nIHBsYWNlczogd2l0aGluIGEgTk9USUNFIHRleHQgZmlsZSBkaXN0cmlidXRlZAogICAgICAgICAgYXMgcGFydCBvZiB0aGUgRGVyaXZhdGl2ZSBXb3Jrczsgd2l0aGluIHRoZSBTb3VyY2UgZm9ybSBvcgogICAgICAgICAgZG9jdW1lbnRhdGlvbiwgaWYgcHJvdmlkZWQgYWxvbmcgd2l0aCB0aGUgRGVyaXZhdGl2ZSBXb3Jrczsgb3IsCiAgICAgICAgICB3aXRoaW4gYSBkaXNwbGF5IGdlbmVyYXRlZCBieSB0aGUgRGVyaXZhdGl2ZSBXb3JrcywgaWYgYW5kCiAgICAgICAgICB3aGVyZXZlciBzdWNoIHRoaXJkLXBhcnR5IG5vdGljZXMgbm9ybWFsbHkgYXBwZWFyLiBUaGUgY29udGVudHMKICAgICAgICAgIG9mIHRoZSBOT1RJQ0UgZmlsZSBhcmUgZm9yIGluZm9ybWF0aW9uYWwgcHVycG9zZXMgb25seSBhbmQKICAgICAgICAgIGRvIG5vdCBtb2RpZnkgdGhlIExpY2Vuc2UuIFlvdSBtYXkgYWRkIFlvdXIgb3duIGF0dHJpYnV0aW9uCiAgICAgICAgICBub3RpY2VzIHdpdGhpbiBEZXJpdmF0aXZlIFdvcmtzIHRoYXQgWW91IGRpc3RyaWJ1dGUsIGFsb25nc2lkZQogICAgICAgICAgb3IgYXMgYW4gYWRkZW5kdW0gdG8gdGhlIE5PVElDRSB0ZXh0IGZyb20gdGhlIFdvcmssIHByb3ZpZGVkCiAgICAgICAgICB0aGF0IHN1Y2ggYWRkaXRpb25hbCBhdHRyaWJ1dGlvbiBub3RpY2VzIGNhbm5vdCBiZSBjb25zdHJ1ZWQKICAgICAgICAgIGFzIG1vZGlmeWluZyB0aGUgTGljZW5zZS4KCiAgICAgIFlvdSBtYXkgYWRkIFlvdXIgb3duIGNvcHlyaWdodCBzdGF0ZW1lbnQgdG8gWW91ciBtb2RpZmljYXRpb25zIGFuZAogICAgICBtYXkgcHJvdmlkZSBhZGRpdGlvbmFsIG9yIGRpZmZlcmVudCBsaWNlbnNlIHRlcm1zIGFuZCBjb25kaXRpb25zCiAgICAgIGZvciB1c2UsIHJlcHJvZHVjdGlvbiwgb3IgZGlzdHJpYnV0aW9uIG9mIFlvdXIgbW9kaWZpY2F0aW9ucywgb3IKICAgICAgZm9yIGFueSBzdWNoIERlcml2YXRpdmUgV29ya3MgYXMgYSB3aG9sZSwgcHJvdmlkZWQgWW91ciB1c2UsCiAgICAgIHJlcHJvZHVjdGlvbiwgYW5kIGRpc3RyaWJ1dGlvbiBvZiB0aGUgV29yayBvdGhlcndpc2UgY29tcGxpZXMgd2l0aAogICAgICB0aGUgY29uZGl0aW9ucyBzdGF0ZWQgaW4gdGhpcyBMaWNlbnNlLgoKICAgNS4gU3VibWlzc2lvbiBvZiBDb250cmlidXRpb25zLiBVbmxlc3MgWW91IGV4cGxpY2l0bHkgc3RhdGUgb3RoZXJ3aXNlLAogICAgICBhbnkgQ29udHJpYnV0aW9uIGludGVudGlvbmFsbHkgc3VibWl0dGVkIGZvciBpbmNsdXNpb24gaW4gdGhlIFdvcmsKICAgICAgYnkgWW91IHRvIHRoZSBMaWNlbnNvciBzaGFsbCBiZSB1bmRlciB0aGUgdGVybXMgYW5kIGNvbmRpdGlvbnMgb2YKICAgICAgdGhpcyBMaWNlbnNlLCB3aXRob3V0IGFueSBhZGRpdGlvbmFsIHRlcm1zIG9yIGNvbmRpdGlvbnMuCiAgICAgIE5vdHdpdGhzdGFuZGluZyB0aGUgYWJvdmUsIG5vdGhpbmcgaGVyZWluIHNoYWxsIHN1cGVyc2VkZSBvciBtb2RpZnkKICAgICAgdGhlIHRlcm1zIG9mIGFueSBzZXBhcmF0ZSBsaWNlbnNlIGFncmVlbWVudCB5b3UgbWF5IGhhdmUgZXhlY3V0ZWQKICAgICAgd2l0aCBMaWNlbnNvciByZWdhcmRpbmcgc3VjaCBDb250cmlidXRpb25zLgoKICAgNi4gVHJhZGVtYXJrcy4gVGhpcyBMaWNlbnNlIGRvZXMgbm90IGdyYW50IHBlcm1pc3Npb24gdG8gdXNlIHRoZSB0cmFkZQogICAgICBuYW1lcywgdHJhZGVtYXJrcywgc2VydmljZSBtYXJrcywgb3IgcHJvZHVjdCBuYW1lcyBvZiB0aGUgTGljZW5zb3IsCiAgICAgIGV4Y2VwdCBhcyByZXF1aXJlZCBmb3IgcmVhc29uYWJsZSBhbmQgY3VzdG9tYXJ5IHVzZSBpbiBkZXNjcmliaW5nIHRoZQogICAgICBvcmlnaW4gb2YgdGhlIFdvcmsgYW5kIHJlcHJvZHVjaW5nIHRoZSBjb250ZW50IG9mIHRoZSBOT1RJQ0UgZmlsZS4KCiAgIDcuIERpc2NsYWltZXIgb2YgV2FycmFudHkuIFVubGVzcyByZXF1aXJlZCBieSBhcHBsaWNhYmxlIGxhdyBvcgogICAgICBhZ3JlZWQgdG8gaW4gd3JpdGluZywgTGljZW5zb3IgcHJvdmlkZXMgdGhlIFdvcmsgKGFuZCBlYWNoCiAgICAgIENvbnRyaWJ1dG9yIHByb3ZpZGVzIGl0cyBDb250cmlidXRpb25zKSBvbiBhbiAiQVMgSVMiIEJBU0lTLAogICAgICBXSVRIT1VUIFdBUlJBTlRJRVMgT1IgQ09ORElUSU9OUyBPRiBBTlkgS0lORCwgZWl0aGVyIGV4cHJlc3Mgb3IKICAgICAgaW1wbGllZCwgaW5jbHVkaW5nLCB3aXRob3V0IGxpbWl0YXRpb24sIGFueSB3YXJyYW50aWVzIG9yIGNvbmRpdGlvbnMKICAgICAgb2YgVElUTEUsIE5PTi1JTkZSSU5HRU1FTlQsIE1FUkNIQU5UQUJJTElUWSwgb3IgRklUTkVTUyBGT1IgQQogICAgICBQQVJUSUNVTEFSIFBVUlBPU0UuIFlvdSBhcmUgc29sZWx5IHJlc3BvbnNpYmxlIGZvciBkZXRlcm1pbmluZyB0aGUKICAgICAgYXBwcm9wcmlhdGVuZXNzIG9mIHVzaW5nIG9yIHJlZGlzdHJpYnV0aW5nIHRoZSBXb3JrIGFuZCBhc3N1bWUgYW55CiAgICAgIHJpc2tzIGFzc29jaWF0ZWQgd2l0aCBZb3VyIGV4ZXJjaXNlIG9mIHBlcm1pc3Npb25zIHVuZGVyIHRoaXMgTGljZW5zZS4KCiAgIDguIExpbWl0YXRpb24gb2YgTGlhYmlsaXR5LiBJbiBubyBldmVudCBhbmQgdW5kZXIgbm8gbGVnYWwgdGhlb3J5LAogICAgICB3aGV0aGVyIGluIHRvcnQgKGluY2x1ZGluZyBuZWdsaWdlbmNlKSwgY29udHJhY3QsIG9yIG90aGVyd2lzZSwKICAgICAgdW5sZXNzIHJlcXVpcmVkIGJ5IGFwcGxpY2FibGUgbGF3IChzdWNoIGFzIGRlbGliZXJhdGUgYW5kIGdyb3NzbHkKICAgICAgbmVnbGlnZW50IGFjdHMpIG9yIGFncmVlZCB0byBpbiB3cml0aW5nLCBzaGFsbCBhbnkgQ29udHJpYnV0b3IgYmUKICAgICAgbGlhYmxlIHRvIFlvdSBmb3IgZGFtYWdlcywgaW5jbHVkaW5nIGFueSBkaXJlY3QsIGluZGlyZWN0LCBzcGVjaWFsLAogICAgICBpbmNpZGVudGFsLCBvciBjb25zZXF1ZW50aWFsIGRhbWFnZXMgb2YgYW55IGNoYXJhY3RlciBhcmlzaW5nIGFzIGEKICAgICAgcmVzdWx0IG9mIHRoaXMgTGljZW5zZSBvciBvdXQgb2YgdGhlIHVzZSBvciBpbmFiaWxpdHkgdG8gdXNlIHRoZQogICAgICBXb3JrIChpbmNsdWRpbmcgYnV0IG5vdCBsaW1pdGVkIHRvIGRhbWFnZXMgZm9yIGxvc3Mgb2YgZ29vZHdpbGwsCiAgICAgIHdvcmsgc3RvcHBhZ2UsIGNvbXB1dGVyIGZhaWx1cmUgb3IgbWFsZnVuY3Rpb24sIG9yIGFueSBhbmQgYWxsCiAgICAgIG90aGVyIGNvbW1lcmNpYWwgZGFtYWdlcyBvciBsb3NzZXMpLCBldmVuIGlmIHN1Y2ggQ29udHJpYnV0b3IKICAgICAgaGFzIGJlZW4gYWR2aXNlZCBvZiB0aGUgcG9zc2liaWxpdHkgb2Ygc3VjaCBkYW1hZ2VzLgoKICAgOS4gQWNjZXB0aW5nIFdhcnJhbnR5IG9yIEFkZGl0aW9uYWwgTGlhYmlsaXR5LiBXaGlsZSByZWRpc3RyaWJ1dGluZwogICAgICB0aGUgV29yayBvciBEZXJpdmF0aXZlIFdvcmtzIHRoZXJlb2YsIFlvdSBtYXkgY2hvb3NlIHRvIG9mZmVyLAogICAgICBhbmQgY2hhcmdlIGEgZmVlIGZvciwgYWNjZXB0YW5jZSBvZiBzdXBwb3J0LCB3YXJyYW50eSwgaW5kZW1uaXR5LAogICAgICBvciBvdGhlciBsaWFiaWxpdHkgb2JsaWdhdGlvbnMgYW5kL29yIHJpZ2h0cyBjb25zaXN0ZW50IHdpdGggdGhpcwogICAgICBMaWNlbnNlLiBIb3dldmVyLCBpbiBhY2NlcHRpbmcgc3VjaCBvYmxpZ2F0aW9ucywgWW91IG1heSBhY3Qgb25seQogICAgICBvbiBZb3VyIG93biBiZWhhbGYgYW5kIG9uIFlvdXIgc29sZSByZXNwb25zaWJpbGl0eSwgbm90IG9uIGJlaGFsZgogICAgICBvZiBhbnkgb3RoZXIgQ29udHJpYnV0b3IsIGFuZCBvbmx5IGlmIFlvdSBhZ3JlZSB0byBpbmRlbW5pZnksCiAgICAgIGRlZmVuZCwgYW5kIGhvbGQgZWFjaCBDb250cmlidXRvciBoYXJtbGVzcyBmb3IgYW55IGxpYWJpbGl0eQogICAgICBpbmN1cnJlZCBieSwgb3IgY2xhaW1zIGFzc2VydGVkIGFnYWluc3QsIHN1Y2ggQ29udHJpYnV0b3IgYnkgcmVhc29uCiAgICAgIG9mIHlvdXIgYWNjZXB0aW5nIGFueSBzdWNoIHdhcnJhbnR5IG9yIGFkZGl0aW9uYWwgbGlhYmlsaXR5LgoKICAgRU5EIE9GIFRFUk1TIEFORCBDT05ESVRJT05TCg=="
    },
    {
      "path": "original/dist/node_modules/node-gyp/gyp/pylib/packaging/LICENSE.BSD",
      "source_locator": "Temp/elite-v334-a443d28e2128479b87937c040aa8b6a8/projection-rebuilt/payload/dist/node_modules/node-gyp/gyp/pylib/packaging/LICENSE.BSD",
      "scope": "VERBATIM existing selected payload notice; applicability remains per path",
      "origin": "ORIGINAL_SELECTED_PAYLOAD",
      "payload_path": "dist/node_modules/node-gyp/gyp/pylib/packaging/LICENSE.BSD",
      "bytes": 1344,
      "sha256": "b70e7e9b742f1cc6f948b34c16aa39ffece94196364bc88ff0d2180f0028fac5",
      "content_base64": "Q29weXJpZ2h0IChjKSBEb25hbGQgU3R1ZmZ0IGFuZCBpbmRpdmlkdWFsIGNvbnRyaWJ1dG9ycy4KQWxsIHJpZ2h0cyByZXNlcnZlZC4KClJlZGlzdHJpYnV0aW9uIGFuZCB1c2UgaW4gc291cmNlIGFuZCBiaW5hcnkgZm9ybXMsIHdpdGggb3Igd2l0aG91dAptb2RpZmljYXRpb24sIGFyZSBwZXJtaXR0ZWQgcHJvdmlkZWQgdGhhdCB0aGUgZm9sbG93aW5nIGNvbmRpdGlvbnMgYXJlIG1ldDoKCiAgICAxLiBSZWRpc3RyaWJ1dGlvbnMgb2Ygc291cmNlIGNvZGUgbXVzdCByZXRhaW4gdGhlIGFib3ZlIGNvcHlyaWdodCBub3RpY2UsCiAgICAgICB0aGlzIGxpc3Qgb2YgY29uZGl0aW9ucyBhbmQgdGhlIGZvbGxvd2luZyBkaXNjbGFpbWVyLgoKICAgIDIuIFJlZGlzdHJpYnV0aW9ucyBpbiBiaW5hcnkgZm9ybSBtdXN0IHJlcHJvZHVjZSB0aGUgYWJvdmUgY29weXJpZ2h0CiAgICAgICBub3RpY2UsIHRoaXMgbGlzdCBvZiBjb25kaXRpb25zIGFuZCB0aGUgZm9sbG93aW5nIGRpc2NsYWltZXIgaW4gdGhlCiAgICAgICBkb2N1bWVudGF0aW9uIGFuZC9vciBvdGhlciBtYXRlcmlhbHMgcHJvdmlkZWQgd2l0aCB0aGUgZGlzdHJpYnV0aW9uLgoKVEhJUyBTT0ZUV0FSRSBJUyBQUk9WSURFRCBCWSBUSEUgQ09QWVJJR0hUIEhPTERFUlMgQU5EIENPTlRSSUJVVE9SUyAiQVMgSVMiIEFORApBTlkgRVhQUkVTUyBPUiBJTVBMSUVEIFdBUlJBTlRJRVMsIElOQ0xVRElORywgQlVUIE5PVCBMSU1JVEVEIFRPLCBUSEUgSU1QTElFRApXQVJSQU5USUVTIE9GIE1FUkNIQU5UQUJJTElUWSBBTkQgRklUTkVTUyBGT1IgQSBQQVJUSUNVTEFSIFBVUlBPU0UgQVJFCkRJU0NMQUlNRUQuIElOIE5PIEVWRU5UIFNIQUxMIFRIRSBDT1BZUklHSFQgSE9MREVSIE9SIENPTlRSSUJVVE9SUyBCRSBMSUFCTEUKRk9SIEFOWSBESVJFQ1QsIElORElSRUNULCBJTkNJREVOVEFMLCBTUEVDSUFMLCBFWEVNUExBUlksIE9SIENPTlNFUVVFTlRJQUwKREFNQUdFUyAoSU5DTFVESU5HLCBCVVQgTk9UIExJTUlURUQgVE8sIFBST0NVUkVNRU5UIE9GIFNVQlNUSVRVVEUgR09PRFMgT1IKU0VSVklDRVM7IExPU1MgT0YgVVNFLCBEQVRBLCBPUiBQUk9GSVRTOyBPUiBCVVNJTkVTUyBJTlRFUlJVUFRJT04pIEhPV0VWRVIKQ0FVU0VEIEFORCBPTiBBTlkgVEhFT1JZIE9GIExJQUJJTElUWSwgV0hFVEhFUiBJTiBDT05UUkFDVCwgU1RSSUNUIExJQUJJTElUWSwKT1IgVE9SVCAoSU5DTFVESU5HIE5FR0xJR0VOQ0UgT1IgT1RIRVJXSVNFKSBBUklTSU5HIElOIEFOWSBXQVkgT1VUIE9GIFRIRSBVU0UKT0YgVEhJUyBTT0ZUV0FSRSwgRVZFTiBJRiBBRFZJU0VEIE9GIFRIRSBQT1NTSUJJTElUWSBPRiBTVUNIIERBTUFHRS4K"
    },
    {
      "path": "original/dist/node_modules/node-gyp/LICENSE",
      "source_locator": "Temp/elite-v334-a443d28e2128479b87937c040aa8b6a8/projection-rebuilt/payload/dist/node_modules/node-gyp/LICENSE",
      "scope": "VERBATIM existing selected payload notice; applicability remains per path",
      "origin": "ORIGINAL_SELECTED_PAYLOAD",
      "payload_path": "dist/node_modules/node-gyp/LICENSE",
      "bytes": 1102,
      "sha256": "662a1b0115251cfb29c6aed0f221f8847bc49c6365d1c53a62c9f4bccc2489c3",
      "content_base64": "KFRoZSBNSVQgTGljZW5zZSkKCkNvcHlyaWdodCAoYykgMjAxMiBOYXRoYW4gUmFqbGljaCA8bmF0aGFuQHRvb3RhbGxuYXRlLm5ldD4KClBlcm1pc3Npb24gaXMgaGVyZWJ5IGdyYW50ZWQsIGZyZWUgb2YgY2hhcmdlLCB0byBhbnkgcGVyc29uCm9idGFpbmluZyBhIGNvcHkgb2YgdGhpcyBzb2Z0d2FyZSBhbmQgYXNzb2NpYXRlZCBkb2N1bWVudGF0aW9uCmZpbGVzICh0aGUgIlNvZnR3YXJlIiksIHRvIGRlYWwgaW4gdGhlIFNvZnR3YXJlIHdpdGhvdXQKcmVzdHJpY3Rpb24sIGluY2x1ZGluZyB3aXRob3V0IGxpbWl0YXRpb24gdGhlIHJpZ2h0cyB0byB1c2UsCmNvcHksIG1vZGlmeSwgbWVyZ2UsIHB1Ymxpc2gsIGRpc3RyaWJ1dGUsIHN1YmxpY2Vuc2UsIGFuZC9vciBzZWxsCmNvcGllcyBvZiB0aGUgU29mdHdhcmUsIGFuZCB0byBwZXJtaXQgcGVyc29ucyB0byB3aG9tIHRoZQpTb2Z0d2FyZSBpcyBmdXJuaXNoZWQgdG8gZG8gc28sIHN1YmplY3QgdG8gdGhlIGZvbGxvd2luZwpjb25kaXRpb25zOgoKVGhlIGFib3ZlIGNvcHlyaWdodCBub3RpY2UgYW5kIHRoaXMgcGVybWlzc2lvbiBub3RpY2Ugc2hhbGwgYmUKaW5jbHVkZWQgaW4gYWxsIGNvcGllcyBvciBzdWJzdGFudGlhbCBwb3J0aW9ucyBvZiB0aGUgU29mdHdhcmUuCgpUSEUgU09GVFdBUkUgSVMgUFJPVklERUQgIkFTIElTIiwgV0lUSE9VVCBXQVJSQU5UWSBPRiBBTlkgS0lORCwKRVhQUkVTUyBPUiBJTVBMSUVELCBJTkNMVURJTkcgQlVUIE5PVCBMSU1JVEVEIFRPIFRIRSBXQVJSQU5USUVTCk9GIE1FUkNIQU5UQUJJTElUWSwgRklUTkVTUyBGT1IgQSBQQVJUSUNVTEFSIFBVUlBPU0UgQU5ECk5PTklORlJJTkdFTUVOVC4gSU4gTk8gRVZFTlQgU0hBTEwgVEhFIEFVVEhPUlMgT1IgQ09QWVJJR0hUCkhPTERFUlMgQkUgTElBQkxFIEZPUiBBTlkgQ0xBSU0sIERBTUFHRVMgT1IgT1RIRVIgTElBQklMSVRZLApXSEVUSEVSIElOIEFOIEFDVElPTiBPRiBDT05UUkFDVCwgVE9SVCBPUiBPVEhFUldJU0UsIEFSSVNJTkcKRlJPTSwgT1VUIE9GIE9SIElOIENPTk5FQ1RJT04gV0lUSCBUSEUgU09GVFdBUkUgT1IgVEhFIFVTRSBPUgpPVEhFUiBERUFMSU5HUyBJTiBUSEUgU09GVFdBUkUuCg=="
    },
    {
      "path": "original/dist/node_modules/nopt/LICENSE",
      "source_locator": "Temp/elite-v334-a443d28e2128479b87937c040aa8b6a8/projection-rebuilt/payload/dist/node_modules/nopt/LICENSE",
      "scope": "VERBATIM existing selected payload notice; applicability remains per path",
      "origin": "ORIGINAL_SELECTED_PAYLOAD",
      "payload_path": "dist/node_modules/nopt/LICENSE",
      "bytes": 765,
      "sha256": "4ec3d4c66cd87f5c8d8ad911b10f99bf27cb00cdfcff82621956e379186b016b",
      "content_base64": "VGhlIElTQyBMaWNlbnNlCgpDb3B5cmlnaHQgKGMpIElzYWFjIFouIFNjaGx1ZXRlciBhbmQgQ29udHJpYnV0b3JzCgpQZXJtaXNzaW9uIHRvIHVzZSwgY29weSwgbW9kaWZ5LCBhbmQvb3IgZGlzdHJpYnV0ZSB0aGlzIHNvZnR3YXJlIGZvciBhbnkKcHVycG9zZSB3aXRoIG9yIHdpdGhvdXQgZmVlIGlzIGhlcmVieSBncmFudGVkLCBwcm92aWRlZCB0aGF0IHRoZSBhYm92ZQpjb3B5cmlnaHQgbm90aWNlIGFuZCB0aGlzIHBlcm1pc3Npb24gbm90aWNlIGFwcGVhciBpbiBhbGwgY29waWVzLgoKVEhFIFNPRlRXQVJFIElTIFBST1ZJREVEICJBUyBJUyIgQU5EIFRIRSBBVVRIT1IgRElTQ0xBSU1TIEFMTCBXQVJSQU5USUVTCldJVEggUkVHQVJEIFRPIFRISVMgU09GVFdBUkUgSU5DTFVESU5HIEFMTCBJTVBMSUVEIFdBUlJBTlRJRVMgT0YKTUVSQ0hBTlRBQklMSVRZIEFORCBGSVRORVNTLiBJTiBOTyBFVkVOVCBTSEFMTCBUSEUgQVVUSE9SIEJFIExJQUJMRSBGT1IKQU5ZIFNQRUNJQUwsIERJUkVDVCwgSU5ESVJFQ1QsIE9SIENPTlNFUVVFTlRJQUwgREFNQUdFUyBPUiBBTlkgREFNQUdFUwpXSEFUU09FVkVSIFJFU1VMVElORyBGUk9NIExPU1MgT0YgVVNFLCBEQVRBIE9SIFBST0ZJVFMsIFdIRVRIRVIgSU4gQU4KQUNUSU9OIE9GIENPTlRSQUNULCBORUdMSUdFTkNFIE9SIE9USEVSIFRPUlRJT1VTIEFDVElPTiwgQVJJU0lORyBPVVQgT0YgT1IKSU4gQ09OTkVDVElPTiBXSVRIIFRIRSBVU0UgT1IgUEVSRk9STUFOQ0UgT0YgVEhJUyBTT0ZUV0FSRS4K"
    },
    {
      "path": "original/dist/node_modules/picomatch/LICENSE",
      "source_locator": "Temp/elite-v334-a443d28e2128479b87937c040aa8b6a8/projection-rebuilt/payload/dist/node_modules/picomatch/LICENSE",
      "scope": "VERBATIM existing selected payload notice; applicability remains per path",
      "origin": "ORIGINAL_SELECTED_PAYLOAD",
      "payload_path": "dist/node_modules/picomatch/LICENSE",
      "bytes": 1091,
      "sha256": "d0cd141b0c322fded5dfad1d4645bb2fedfc05b7321fe1009469638190d59ef9",
      "content_base64": "VGhlIE1JVCBMaWNlbnNlIChNSVQpCgpDb3B5cmlnaHQgKGMpIDIwMTctcHJlc2VudCwgSm9uIFNjaGxpbmtlcnQuCgpQZXJtaXNzaW9uIGlzIGhlcmVieSBncmFudGVkLCBmcmVlIG9mIGNoYXJnZSwgdG8gYW55IHBlcnNvbiBvYnRhaW5pbmcgYSBjb3B5Cm9mIHRoaXMgc29mdHdhcmUgYW5kIGFzc29jaWF0ZWQgZG9jdW1lbnRhdGlvbiBmaWxlcyAodGhlICJTb2Z0d2FyZSIpLCB0byBkZWFsCmluIHRoZSBTb2Z0d2FyZSB3aXRob3V0IHJlc3RyaWN0aW9uLCBpbmNsdWRpbmcgd2l0aG91dCBsaW1pdGF0aW9uIHRoZSByaWdodHMKdG8gdXNlLCBjb3B5LCBtb2RpZnksIG1lcmdlLCBwdWJsaXNoLCBkaXN0cmlidXRlLCBzdWJsaWNlbnNlLCBhbmQvb3Igc2VsbApjb3BpZXMgb2YgdGhlIFNvZnR3YXJlLCBhbmQgdG8gcGVybWl0IHBlcnNvbnMgdG8gd2hvbSB0aGUgU29mdHdhcmUgaXMKZnVybmlzaGVkIHRvIGRvIHNvLCBzdWJqZWN0IHRvIHRoZSBmb2xsb3dpbmcgY29uZGl0aW9uczoKClRoZSBhYm92ZSBjb3B5cmlnaHQgbm90aWNlIGFuZCB0aGlzIHBlcm1pc3Npb24gbm90aWNlIHNoYWxsIGJlIGluY2x1ZGVkIGluCmFsbCBjb3BpZXMgb3Igc3Vic3RhbnRpYWwgcG9ydGlvbnMgb2YgdGhlIFNvZnR3YXJlLgoKVEhFIFNPRlRXQVJFIElTIFBST1ZJREVEICJBUyBJUyIsIFdJVEhPVVQgV0FSUkFOVFkgT0YgQU5ZIEtJTkQsIEVYUFJFU1MgT1IKSU1QTElFRCwgSU5DTFVESU5HIEJVVCBOT1QgTElNSVRFRCBUTyBUSEUgV0FSUkFOVElFUyBPRiBNRVJDSEFOVEFCSUxJVFksCkZJVE5FU1MgRk9SIEEgUEFSVElDVUxBUiBQVVJQT1NFIEFORCBOT05JTkZSSU5HRU1FTlQuIElOIE5PIEVWRU5UIFNIQUxMIFRIRQpBVVRIT1JTIE9SIENPUFlSSUdIVCBIT0xERVJTIEJFIExJQUJMRSBGT1IgQU5ZIENMQUlNLCBEQU1BR0VTIE9SIE9USEVSCkxJQUJJTElUWSwgV0hFVEhFUiBJTiBBTiBBQ1RJT04gT0YgQ09OVFJBQ1QsIFRPUlQgT1IgT1RIRVJXSVNFLCBBUklTSU5HIEZST00sCk9VVCBPRiBPUiBJTiBDT05ORUNUSU9OIFdJVEggVEhFIFNPRlRXQVJFIE9SIFRIRSBVU0UgT1IgT1RIRVIgREVBTElOR1MgSU4KVEhFIFNPRlRXQVJFLgo="
    },
    {
      "path": "original/dist/node_modules/proc-log/LICENSE",
      "source_locator": "Temp/elite-v334-a443d28e2128479b87937c040aa8b6a8/projection-rebuilt/payload/dist/node_modules/proc-log/LICENSE",
      "scope": "VERBATIM existing selected payload notice; applicability remains per path",
      "origin": "ORIGINAL_SELECTED_PAYLOAD",
      "payload_path": "dist/node_modules/proc-log/LICENSE",
      "bytes": 742,
      "sha256": "dc32a0dee275e0a9aeffbc974dbf4899a30dcdc2e5ffa8934aecb69261065864",
      "content_base64": "VGhlIElTQyBMaWNlbnNlCgpDb3B5cmlnaHQgKGMpIEdpdEh1YiwgSW5jLgoKUGVybWlzc2lvbiB0byB1c2UsIGNvcHksIG1vZGlmeSwgYW5kL29yIGRpc3RyaWJ1dGUgdGhpcyBzb2Z0d2FyZSBmb3IgYW55CnB1cnBvc2Ugd2l0aCBvciB3aXRob3V0IGZlZSBpcyBoZXJlYnkgZ3JhbnRlZCwgcHJvdmlkZWQgdGhhdCB0aGUgYWJvdmUKY29weXJpZ2h0IG5vdGljZSBhbmQgdGhpcyBwZXJtaXNzaW9uIG5vdGljZSBhcHBlYXIgaW4gYWxsIGNvcGllcy4KClRIRSBTT0ZUV0FSRSBJUyBQUk9WSURFRCAiQVMgSVMiIEFORCBUSEUgQVVUSE9SIERJU0NMQUlNUyBBTEwgV0FSUkFOVElFUwpXSVRIIFJFR0FSRCBUTyBUSElTIFNPRlRXQVJFIElOQ0xVRElORyBBTEwgSU1QTElFRCBXQVJSQU5USUVTIE9GCk1FUkNIQU5UQUJJTElUWSBBTkQgRklUTkVTUy4gSU4gTk8gRVZFTlQgU0hBTEwgVEhFIEFVVEhPUiBCRSBMSUFCTEUgRk9SCkFOWSBTUEVDSUFMLCBESVJFQ1QsIElORElSRUNULCBPUiBDT05TRVFVRU5USUFMIERBTUFHRVMgT1IgQU5ZIERBTUFHRVMKV0hBVFNPRVZFUiBSRVNVTFRJTkcgRlJPTSBMT1NTIE9GIFVTRSwgREFUQSBPUiBQUk9GSVRTLCBXSEVUSEVSIElOIEFOCkFDVElPTiBPRiBDT05UUkFDVCwgTkVHTElHRU5DRSBPUiBPVEhFUiBUT1JUSU9VUyBBQ1RJT04sIEFSSVNJTkcgT1VUIE9GIE9SCklOIENPTk5FQ1RJT04gV0lUSCBUSEUgVVNFIE9SIFBFUkZPUk1BTkNFIE9GIFRISVMgU09GVFdBUkUuCg=="
    },
    {
      "path": "original/dist/node_modules/semver/LICENSE",
      "source_locator": "Temp/elite-v334-a443d28e2128479b87937c040aa8b6a8/projection-rebuilt/payload/dist/node_modules/semver/LICENSE",
      "scope": "VERBATIM existing selected payload notice; applicability remains per path",
      "origin": "ORIGINAL_SELECTED_PAYLOAD",
      "payload_path": "dist/node_modules/semver/LICENSE",
      "bytes": 765,
      "sha256": "4ec3d4c66cd87f5c8d8ad911b10f99bf27cb00cdfcff82621956e379186b016b",
      "content_base64": "VGhlIElTQyBMaWNlbnNlCgpDb3B5cmlnaHQgKGMpIElzYWFjIFouIFNjaGx1ZXRlciBhbmQgQ29udHJpYnV0b3JzCgpQZXJtaXNzaW9uIHRvIHVzZSwgY29weSwgbW9kaWZ5LCBhbmQvb3IgZGlzdHJpYnV0ZSB0aGlzIHNvZnR3YXJlIGZvciBhbnkKcHVycG9zZSB3aXRoIG9yIHdpdGhvdXQgZmVlIGlzIGhlcmVieSBncmFudGVkLCBwcm92aWRlZCB0aGF0IHRoZSBhYm92ZQpjb3B5cmlnaHQgbm90aWNlIGFuZCB0aGlzIHBlcm1pc3Npb24gbm90aWNlIGFwcGVhciBpbiBhbGwgY29waWVzLgoKVEhFIFNPRlRXQVJFIElTIFBST1ZJREVEICJBUyBJUyIgQU5EIFRIRSBBVVRIT1IgRElTQ0xBSU1TIEFMTCBXQVJSQU5USUVTCldJVEggUkVHQVJEIFRPIFRISVMgU09GVFdBUkUgSU5DTFVESU5HIEFMTCBJTVBMSUVEIFdBUlJBTlRJRVMgT0YKTUVSQ0hBTlRBQklMSVRZIEFORCBGSVRORVNTLiBJTiBOTyBFVkVOVCBTSEFMTCBUSEUgQVVUSE9SIEJFIExJQUJMRSBGT1IKQU5ZIFNQRUNJQUwsIERJUkVDVCwgSU5ESVJFQ1QsIE9SIENPTlNFUVVFTlRJQUwgREFNQUdFUyBPUiBBTlkgREFNQUdFUwpXSEFUU09FVkVSIFJFU1VMVElORyBGUk9NIExPU1MgT0YgVVNFLCBEQVRBIE9SIFBST0ZJVFMsIFdIRVRIRVIgSU4gQU4KQUNUSU9OIE9GIENPTlRSQUNULCBORUdMSUdFTkNFIE9SIE9USEVSIFRPUlRJT1VTIEFDVElPTiwgQVJJU0lORyBPVVQgT0YgT1IKSU4gQ09OTkVDVElPTiBXSVRIIFRIRSBVU0UgT1IgUEVSRk9STUFOQ0UgT0YgVEhJUyBTT0ZUV0FSRS4K"
    },
    {
      "path": "original/dist/node_modules/tinyglobby/LICENSE",
      "source_locator": "Temp/elite-v334-a443d28e2128479b87937c040aa8b6a8/projection-rebuilt/payload/dist/node_modules/tinyglobby/LICENSE",
      "scope": "VERBATIM existing selected payload notice; applicability remains per path",
      "origin": "ORIGINAL_SELECTED_PAYLOAD",
      "payload_path": "dist/node_modules/tinyglobby/LICENSE",
      "bytes": 1076,
      "sha256": "22c68811e174cbbfb3813d4135918df4f540959c14d872e601d9abe83d3cde8f",
      "content_base64": "TUlUIExpY2Vuc2UKCkNvcHlyaWdodCAoYykgMjAyNCBNYWRlbGluZSBHdXJyaWFyw6FuCgpQZXJtaXNzaW9uIGlzIGhlcmVieSBncmFudGVkLCBmcmVlIG9mIGNoYXJnZSwgdG8gYW55IHBlcnNvbiBvYnRhaW5pbmcgYSBjb3B5Cm9mIHRoaXMgc29mdHdhcmUgYW5kIGFzc29jaWF0ZWQgZG9jdW1lbnRhdGlvbiBmaWxlcyAodGhlICJTb2Z0d2FyZSIpLCB0byBkZWFsCmluIHRoZSBTb2Z0d2FyZSB3aXRob3V0IHJlc3RyaWN0aW9uLCBpbmNsdWRpbmcgd2l0aG91dCBsaW1pdGF0aW9uIHRoZSByaWdodHMKdG8gdXNlLCBjb3B5LCBtb2RpZnksIG1lcmdlLCBwdWJsaXNoLCBkaXN0cmlidXRlLCBzdWJsaWNlbnNlLCBhbmQvb3Igc2VsbApjb3BpZXMgb2YgdGhlIFNvZnR3YXJlLCBhbmQgdG8gcGVybWl0IHBlcnNvbnMgdG8gd2hvbSB0aGUgU29mdHdhcmUgaXMKZnVybmlzaGVkIHRvIGRvIHNvLCBzdWJqZWN0IHRvIHRoZSBmb2xsb3dpbmcgY29uZGl0aW9uczoKClRoZSBhYm92ZSBjb3B5cmlnaHQgbm90aWNlIGFuZCB0aGlzIHBlcm1pc3Npb24gbm90aWNlIHNoYWxsIGJlIGluY2x1ZGVkIGluIGFsbApjb3BpZXMgb3Igc3Vic3RhbnRpYWwgcG9ydGlvbnMgb2YgdGhlIFNvZnR3YXJlLgoKVEhFIFNPRlRXQVJFIElTIFBST1ZJREVEICJBUyBJUyIsIFdJVEhPVVQgV0FSUkFOVFkgT0YgQU5ZIEtJTkQsIEVYUFJFU1MgT1IKSU1QTElFRCwgSU5DTFVESU5HIEJVVCBOT1QgTElNSVRFRCBUTyBUSEUgV0FSUkFOVElFUyBPRiBNRVJDSEFOVEFCSUxJVFksCkZJVE5FU1MgRk9SIEEgUEFSVElDVUxBUiBQVVJQT1NFIEFORCBOT05JTkZSSU5HRU1FTlQuIElOIE5PIEVWRU5UIFNIQUxMIFRIRQpBVVRIT1JTIE9SIENPUFlSSUdIVCBIT0xERVJTIEJFIExJQUJMRSBGT1IgQU5ZIENMQUlNLCBEQU1BR0VTIE9SIE9USEVSCkxJQUJJTElUWSwgV0hFVEhFUiBJTiBBTiBBQ1RJT04gT0YgQ09OVFJBQ1QsIFRPUlQgT1IgT1RIRVJXSVNFLCBBUklTSU5HIEZST00sCk9VVCBPRiBPUiBJTiBDT05ORUNUSU9OIFdJVEggVEhFIFNPRlRXQVJFIE9SIFRIRSBVU0UgT1IgT1RIRVIgREVBTElOR1MgSU4gVEhFClNPRlRXQVJFLgo="
    },
    {
      "path": "original/dist/node_modules/undici/lib/web/fetch/LICENSE",
      "source_locator": "Temp/elite-v334-a443d28e2128479b87937c040aa8b6a8/projection-rebuilt/payload/dist/node_modules/undici/lib/web/fetch/LICENSE",
      "scope": "VERBATIM existing selected payload notice; applicability remains per path",
      "origin": "ORIGINAL_SELECTED_PAYLOAD",
      "payload_path": "dist/node_modules/undici/lib/web/fetch/LICENSE",
      "bytes": 1071,
      "sha256": "a21d6c8d3fc631198b97e57584697f7f9c37805c71c6cf9a12e4442b40394b88",
      "content_base64": "TUlUIExpY2Vuc2UKCkNvcHlyaWdodCAoYykgMjAyMCBFdGhhbiBBcnJvd29vZAoKUGVybWlzc2lvbiBpcyBoZXJlYnkgZ3JhbnRlZCwgZnJlZSBvZiBjaGFyZ2UsIHRvIGFueSBwZXJzb24gb2J0YWluaW5nIGEgY29weQpvZiB0aGlzIHNvZnR3YXJlIGFuZCBhc3NvY2lhdGVkIGRvY3VtZW50YXRpb24gZmlsZXMgKHRoZSAiU29mdHdhcmUiKSwgdG8gZGVhbAppbiB0aGUgU29mdHdhcmUgd2l0aG91dCByZXN0cmljdGlvbiwgaW5jbHVkaW5nIHdpdGhvdXQgbGltaXRhdGlvbiB0aGUgcmlnaHRzCnRvIHVzZSwgY29weSwgbW9kaWZ5LCBtZXJnZSwgcHVibGlzaCwgZGlzdHJpYnV0ZSwgc3VibGljZW5zZSwgYW5kL29yIHNlbGwKY29waWVzIG9mIHRoZSBTb2Z0d2FyZSwgYW5kIHRvIHBlcm1pdCBwZXJzb25zIHRvIHdob20gdGhlIFNvZnR3YXJlIGlzCmZ1cm5pc2hlZCB0byBkbyBzbywgc3ViamVjdCB0byB0aGUgZm9sbG93aW5nIGNvbmRpdGlvbnM6CgpUaGUgYWJvdmUgY29weXJpZ2h0IG5vdGljZSBhbmQgdGhpcyBwZXJtaXNzaW9uIG5vdGljZSBzaGFsbCBiZSBpbmNsdWRlZCBpbiBhbGwKY29waWVzIG9yIHN1YnN0YW50aWFsIHBvcnRpb25zIG9mIHRoZSBTb2Z0d2FyZS4KClRIRSBTT0ZUV0FSRSBJUyBQUk9WSURFRCAiQVMgSVMiLCBXSVRIT1VUIFdBUlJBTlRZIE9GIEFOWSBLSU5ELCBFWFBSRVNTIE9SCklNUExJRUQsIElOQ0xVRElORyBCVVQgTk9UIExJTUlURUQgVE8gVEhFIFdBUlJBTlRJRVMgT0YgTUVSQ0hBTlRBQklMSVRZLApGSVRORVNTIEZPUiBBIFBBUlRJQ1VMQVIgUFVSUE9TRSBBTkQgTk9OSU5GUklOR0VNRU5ULiBJTiBOTyBFVkVOVCBTSEFMTCBUSEUKQVVUSE9SUyBPUiBDT1BZUklHSFQgSE9MREVSUyBCRSBMSUFCTEUgRk9SIEFOWSBDTEFJTSwgREFNQUdFUyBPUiBPVEhFUgpMSUFCSUxJVFksIFdIRVRIRVIgSU4gQU4gQUNUSU9OIE9GIENPTlRSQUNULCBUT1JUIE9SIE9USEVSV0lTRSwgQVJJU0lORyBGUk9NLApPVVQgT0YgT1IgSU4gQ09OTkVDVElPTiBXSVRIIFRIRSBTT0ZUV0FSRSBPUiBUSEUgVVNFIE9SIE9USEVSIERFQUxJTkdTIElOIFRIRQpTT0ZUV0FSRS4K"
    },
    {
      "path": "original/dist/node_modules/undici/LICENSE",
      "source_locator": "Temp/elite-v334-a443d28e2128479b87937c040aa8b6a8/projection-rebuilt/payload/dist/node_modules/undici/LICENSE",
      "scope": "VERBATIM existing selected payload notice; applicability remains per path",
      "origin": "ORIGINAL_SELECTED_PAYLOAD",
      "payload_path": "dist/node_modules/undici/LICENSE",
      "bytes": 1090,
      "sha256": "a6db8096b2707bc0102d256917d4d33f298ba36d8c3f25de067a2b5bb379db27",
      "content_base64": "TUlUIExpY2Vuc2UKCkNvcHlyaWdodCAoYykgTWF0dGVvIENvbGxpbmEgYW5kIFVuZGljaSBjb250cmlidXRvcnMKClBlcm1pc3Npb24gaXMgaGVyZWJ5IGdyYW50ZWQsIGZyZWUgb2YgY2hhcmdlLCB0byBhbnkgcGVyc29uIG9idGFpbmluZyBhIGNvcHkKb2YgdGhpcyBzb2Z0d2FyZSBhbmQgYXNzb2NpYXRlZCBkb2N1bWVudGF0aW9uIGZpbGVzICh0aGUgIlNvZnR3YXJlIiksIHRvIGRlYWwKaW4gdGhlIFNvZnR3YXJlIHdpdGhvdXQgcmVzdHJpY3Rpb24sIGluY2x1ZGluZyB3aXRob3V0IGxpbWl0YXRpb24gdGhlIHJpZ2h0cwp0byB1c2UsIGNvcHksIG1vZGlmeSwgbWVyZ2UsIHB1Ymxpc2gsIGRpc3RyaWJ1dGUsIHN1YmxpY2Vuc2UsIGFuZC9vciBzZWxsCmNvcGllcyBvZiB0aGUgU29mdHdhcmUsIGFuZCB0byBwZXJtaXQgcGVyc29ucyB0byB3aG9tIHRoZSBTb2Z0d2FyZSBpcwpmdXJuaXNoZWQgdG8gZG8gc28sIHN1YmplY3QgdG8gdGhlIGZvbGxvd2luZyBjb25kaXRpb25zOgoKVGhlIGFib3ZlIGNvcHlyaWdodCBub3RpY2UgYW5kIHRoaXMgcGVybWlzc2lvbiBub3RpY2Ugc2hhbGwgYmUgaW5jbHVkZWQgaW4gYWxsCmNvcGllcyBvciBzdWJzdGFudGlhbCBwb3J0aW9ucyBvZiB0aGUgU29mdHdhcmUuCgpUSEUgU09GVFdBUkUgSVMgUFJPVklERUQgIkFTIElTIiwgV0lUSE9VVCBXQVJSQU5UWSBPRiBBTlkgS0lORCwgRVhQUkVTUyBPUgpJTVBMSUVELCBJTkNMVURJTkcgQlVUIE5PVCBMSU1JVEVEIFRPIFRIRSBXQVJSQU5USUVTIE9GIE1FUkNIQU5UQUJJTElUWSwKRklUTkVTUyBGT1IgQSBQQVJUSUNVTEFSIFBVUlBPU0UgQU5EIE5PTklORlJJTkdFTUVOVC4gSU4gTk8gRVZFTlQgU0hBTEwgVEhFCkFVVEhPUlMgT1IgQ09QWVJJR0hUIEhPTERFUlMgQkUgTElBQkxFIEZPUiBBTlkgQ0xBSU0sIERBTUFHRVMgT1IgT1RIRVIKTElBQklMSVRZLCBXSEVUSEVSIElOIEFOIEFDVElPTiBPRiBDT05UUkFDVCwgVE9SVCBPUiBPVEhFUldJU0UsIEFSSVNJTkcgRlJPTSwKT1VUIE9GIE9SIElOIENPTk5FQ1RJT04gV0lUSCBUSEUgU09GVFdBUkUgT1IgVEhFIFVTRSBPUiBPVEhFUiBERUFMSU5HUyBJTiBUSEUKU09GVFdBUkUuCg=="
    },
    {
      "path": "original/dist/node_modules/v8-compile-cache/LICENSE",
      "source_locator": "Temp/elite-v334-a443d28e2128479b87937c040aa8b6a8/projection-rebuilt/payload/dist/node_modules/v8-compile-cache/LICENSE",
      "scope": "VERBATIM existing selected payload notice; applicability remains per path",
      "origin": "ORIGINAL_SELECTED_PAYLOAD",
      "payload_path": "dist/node_modules/v8-compile-cache/LICENSE",
      "bytes": 1080,
      "sha256": "c77674258a3fdf3036a5d13d2aecd30d7a25aa6191cb0a9a7dd45b975dc7fe69",
      "content_base64": "VGhlIE1JVCBMaWNlbnNlIChNSVQpCgpDb3B5cmlnaHQgKGMpIDIwMTkgQW5kcmVzIFN1YXJlegoKUGVybWlzc2lvbiBpcyBoZXJlYnkgZ3JhbnRlZCwgZnJlZSBvZiBjaGFyZ2UsIHRvIGFueSBwZXJzb24gb2J0YWluaW5nIGEgY29weQpvZiB0aGlzIHNvZnR3YXJlIGFuZCBhc3NvY2lhdGVkIGRvY3VtZW50YXRpb24gZmlsZXMgKHRoZSAiU29mdHdhcmUiKSwgdG8gZGVhbAppbiB0aGUgU29mdHdhcmUgd2l0aG91dCByZXN0cmljdGlvbiwgaW5jbHVkaW5nIHdpdGhvdXQgbGltaXRhdGlvbiB0aGUgcmlnaHRzCnRvIHVzZSwgY29weSwgbW9kaWZ5LCBtZXJnZSwgcHVibGlzaCwgZGlzdHJpYnV0ZSwgc3VibGljZW5zZSwgYW5kL29yIHNlbGwKY29waWVzIG9mIHRoZSBTb2Z0d2FyZSwgYW5kIHRvIHBlcm1pdCBwZXJzb25zIHRvIHdob20gdGhlIFNvZnR3YXJlIGlzCmZ1cm5pc2hlZCB0byBkbyBzbywgc3ViamVjdCB0byB0aGUgZm9sbG93aW5nIGNvbmRpdGlvbnM6CgpUaGUgYWJvdmUgY29weXJpZ2h0IG5vdGljZSBhbmQgdGhpcyBwZXJtaXNzaW9uIG5vdGljZSBzaGFsbCBiZSBpbmNsdWRlZCBpbiBhbGwKY29waWVzIG9yIHN1YnN0YW50aWFsIHBvcnRpb25zIG9mIHRoZSBTb2Z0d2FyZS4KClRIRSBTT0ZUV0FSRSBJUyBQUk9WSURFRCAiQVMgSVMiLCBXSVRIT1VUIFdBUlJBTlRZIE9GIEFOWSBLSU5ELCBFWFBSRVNTIE9SCklNUExJRUQsIElOQ0xVRElORyBCVVQgTk9UIExJTUlURUQgVE8gVEhFIFdBUlJBTlRJRVMgT0YgTUVSQ0hBTlRBQklMSVRZLApGSVRORVNTIEZPUiBBIFBBUlRJQ1VMQVIgUFVSUE9TRSBBTkQgTk9OSU5GUklOR0VNRU5ULiBJTiBOTyBFVkVOVCBTSEFMTCBUSEUKQVVUSE9SUyBPUiBDT1BZUklHSFQgSE9MREVSUyBCRSBMSUFCTEUgRk9SIEFOWSBDTEFJTSwgREFNQUdFUyBPUiBPVEhFUgpMSUFCSUxJVFksIFdIRVRIRVIgSU4gQU4gQUNUSU9OIE9GIENPTlRSQUNULCBUT1JUIE9SIE9USEVSV0lTRSwgQVJJU0lORyBGUk9NLApPVVQgT0YgT1IgSU4gQ09OTkVDVElPTiBXSVRIIFRIRSBTT0ZUV0FSRSBPUiBUSEUgVVNFIE9SIE9USEVSIERFQUxJTkdTIElOIFRIRQpTT0ZUV0FSRS4K"
    },
    {
      "path": "original/dist/node_modules/which/LICENSE",
      "source_locator": "Temp/elite-v334-a443d28e2128479b87937c040aa8b6a8/projection-rebuilt/payload/dist/node_modules/which/LICENSE",
      "scope": "VERBATIM existing selected payload notice; applicability remains per path",
      "origin": "ORIGINAL_SELECTED_PAYLOAD",
      "payload_path": "dist/node_modules/which/LICENSE",
      "bytes": 765,
      "sha256": "4ec3d4c66cd87f5c8d8ad911b10f99bf27cb00cdfcff82621956e379186b016b",
      "content_base64": "VGhlIElTQyBMaWNlbnNlCgpDb3B5cmlnaHQgKGMpIElzYWFjIFouIFNjaGx1ZXRlciBhbmQgQ29udHJpYnV0b3JzCgpQZXJtaXNzaW9uIHRvIHVzZSwgY29weSwgbW9kaWZ5LCBhbmQvb3IgZGlzdHJpYnV0ZSB0aGlzIHNvZnR3YXJlIGZvciBhbnkKcHVycG9zZSB3aXRoIG9yIHdpdGhvdXQgZmVlIGlzIGhlcmVieSBncmFudGVkLCBwcm92aWRlZCB0aGF0IHRoZSBhYm92ZQpjb3B5cmlnaHQgbm90aWNlIGFuZCB0aGlzIHBlcm1pc3Npb24gbm90aWNlIGFwcGVhciBpbiBhbGwgY29waWVzLgoKVEhFIFNPRlRXQVJFIElTIFBST1ZJREVEICJBUyBJUyIgQU5EIFRIRSBBVVRIT1IgRElTQ0xBSU1TIEFMTCBXQVJSQU5USUVTCldJVEggUkVHQVJEIFRPIFRISVMgU09GVFdBUkUgSU5DTFVESU5HIEFMTCBJTVBMSUVEIFdBUlJBTlRJRVMgT0YKTUVSQ0hBTlRBQklMSVRZIEFORCBGSVRORVNTLiBJTiBOTyBFVkVOVCBTSEFMTCBUSEUgQVVUSE9SIEJFIExJQUJMRSBGT1IKQU5ZIFNQRUNJQUwsIERJUkVDVCwgSU5ESVJFQ1QsIE9SIENPTlNFUVVFTlRJQUwgREFNQUdFUyBPUiBBTlkgREFNQUdFUwpXSEFUU09FVkVSIFJFU1VMVElORyBGUk9NIExPU1MgT0YgVVNFLCBEQVRBIE9SIFBST0ZJVFMsIFdIRVRIRVIgSU4gQU4KQUNUSU9OIE9GIENPTlRSQUNULCBORUdMSUdFTkNFIE9SIE9USEVSIFRPUlRJT1VTIEFDVElPTiwgQVJJU0lORyBPVVQgT0YgT1IKSU4gQ09OTkVDVElPTiBXSVRIIFRIRSBVU0UgT1IgUEVSRk9STUFOQ0UgT0YgVEhJUyBTT0ZUV0FSRS4K"
    },
    {
      "path": "original/LICENSE",
      "source_locator": "Temp/elite-v334-a443d28e2128479b87937c040aa8b6a8/projection-rebuilt/payload/LICENSE",
      "scope": "VERBATIM existing selected payload notice; applicability remains per path",
      "origin": "ORIGINAL_SELECTED_PAYLOAD",
      "payload_path": "LICENSE",
      "bytes": 1170,
      "sha256": "e0a867ff513ea7be2a0ddc339ac6a031e459a38668e077b8f0e649544062f9f2",
      "content_base64": "VGhlIE1JVCBMaWNlbnNlIChNSVQpCgpDb3B5cmlnaHQgKGMpIDIwMTUtMjAxNiBSaWNvIFN0YS4gQ3J1eiBhbmQgb3RoZXIgY29udHJpYnV0b3JzCkNvcHlyaWdodCAoYykgMjAxNi0yMDI1IFpvbHRhbiBLb2NoYW4gYW5kIG90aGVyIGNvbnRyaWJ1dG9ycwoKUGVybWlzc2lvbiBpcyBoZXJlYnkgZ3JhbnRlZCwgZnJlZSBvZiBjaGFyZ2UsIHRvIGFueSBwZXJzb24gb2J0YWluaW5nIGEgY29weQpvZiB0aGlzIHNvZnR3YXJlIGFuZCBhc3NvY2lhdGVkIGRvY3VtZW50YXRpb24gZmlsZXMgKHRoZSAiU29mdHdhcmUiKSwgdG8gZGVhbAppbiB0aGUgU29mdHdhcmUgd2l0aG91dCByZXN0cmljdGlvbiwgaW5jbHVkaW5nIHdpdGhvdXQgbGltaXRhdGlvbiB0aGUgcmlnaHRzCnRvIHVzZSwgY29weSwgbW9kaWZ5LCBtZXJnZSwgcHVibGlzaCwgZGlzdHJpYnV0ZSwgc3VibGljZW5zZSwgYW5kL29yIHNlbGwKY29waWVzIG9mIHRoZSBTb2Z0d2FyZSwgYW5kIHRvIHBlcm1pdCBwZXJzb25zIHRvIHdob20gdGhlIFNvZnR3YXJlIGlzCmZ1cm5pc2hlZCB0byBkbyBzbywgc3ViamVjdCB0byB0aGUgZm9sbG93aW5nIGNvbmRpdGlvbnM6CgpUaGUgYWJvdmUgY29weXJpZ2h0IG5vdGljZSBhbmQgdGhpcyBwZXJtaXNzaW9uIG5vdGljZSBzaGFsbCBiZSBpbmNsdWRlZCBpbiBhbGwKY29waWVzIG9yIHN1YnN0YW50aWFsIHBvcnRpb25zIG9mIHRoZSBTb2Z0d2FyZS4KClRIRSBTT0ZUV0FSRSBJUyBQUk9WSURFRCAiQVMgSVMiLCBXSVRIT1VUIFdBUlJBTlRZIE9GIEFOWSBLSU5ELCBFWFBSRVNTIE9SCklNUExJRUQsIElOQ0xVRElORyBCVVQgTk9UIExJTUlURUQgVE8gVEhFIFdBUlJBTlRJRVMgT0YgTUVSQ0hBTlRBQklMSVRZLApGSVRORVNTIEZPUiBBIFBBUlRJQ1VMQVIgUFVSUE9TRSBBTkQgTk9OSU5GUklOR0VNRU5ULiBJTiBOTyBFVkVOVCBTSEFMTCBUSEUKQVVUSE9SUyBPUiBDT1BZUklHSFQgSE9MREVSUyBCRSBMSUFCTEUgRk9SIEFOWSBDTEFJTSwgREFNQUdFUyBPUiBPVEhFUgpMSUFCSUxJVFksIFdIRVRIRVIgSU4gQU4gQUNUSU9OIE9GIENPTlRSQUNULCBUT1JUIE9SIE9USEVSV0lTRSwgQVJJU0lORyBGUk9NLApPVVQgT0YgT1IgSU4gQ09OTkVDVElPTiBXSVRIIFRIRSBTT0ZUV0FSRSBPUiBUSEUgVVNFIE9SIE9USEVSIERFQUxJTkdTIElOIFRIRQpTT0ZUV0FSRS4K"
    },
    {
      "path": "supplemental/GNU-GPL-3.0.txt",
      "source_locator": "Temp/elite-v336-32e87d99be6448e98b47fb0febcdf2be/license-texts/GNU-GPL-3.0.txt",
      "scope": "fixed-source/authority text retained V336; not an archive match",
      "origin": "RETAINED_RESEARCH_EVIDENCE",
      "payload_path": null,
      "bytes": 35149,
      "sha256": "3972dc9744f6499f0f9b2dbf76696f2ae7ad8af9b23dde66d6af86c9dfb36986",
      "content_base64": "ICAgICAgICAgICAgICAgICAgICBHTlUgR0VORVJBTCBQVUJMSUMgTElDRU5TRQogICAgICAgICAgICAgICAgICAgICAgIFZlcnNpb24gMywgMjkgSnVuZSAyMDA3CgogQ29weXJpZ2h0IChDKSAyMDA3IEZyZWUgU29mdHdhcmUgRm91bmRhdGlvbiwgSW5jLiA8aHR0cHM6Ly9mc2Yub3JnLz4KIEV2ZXJ5b25lIGlzIHBlcm1pdHRlZCB0byBjb3B5IGFuZCBkaXN0cmlidXRlIHZlcmJhdGltIGNvcGllcwogb2YgdGhpcyBsaWNlbnNlIGRvY3VtZW50LCBidXQgY2hhbmdpbmcgaXQgaXMgbm90IGFsbG93ZWQuCgogICAgICAgICAgICAgICAgICAgICAgICAgICAgUHJlYW1ibGUKCiAgVGhlIEdOVSBHZW5lcmFsIFB1YmxpYyBMaWNlbnNlIGlzIGEgZnJlZSwgY29weWxlZnQgbGljZW5zZSBmb3IKc29mdHdhcmUgYW5kIG90aGVyIGtpbmRzIG9mIHdvcmtzLgoKICBUaGUgbGljZW5zZXMgZm9yIG1vc3Qgc29mdHdhcmUgYW5kIG90aGVyIHByYWN0aWNhbCB3b3JrcyBhcmUgZGVzaWduZWQKdG8gdGFrZSBhd2F5IHlvdXIgZnJlZWRvbSB0byBzaGFyZSBhbmQgY2hhbmdlIHRoZSB3b3Jrcy4gIEJ5IGNvbnRyYXN0LAp0aGUgR05VIEdlbmVyYWwgUHVibGljIExpY2Vuc2UgaXMgaW50ZW5kZWQgdG8gZ3VhcmFudGVlIHlvdXIgZnJlZWRvbSB0bwpzaGFyZSBhbmQgY2hhbmdlIGFsbCB2ZXJzaW9ucyBvZiBhIHByb2dyYW0tLXRvIG1ha2Ugc3VyZSBpdCByZW1haW5zIGZyZWUKc29mdHdhcmUgZm9yIGFsbCBpdHMgdXNlcnMuICBXZSwgdGhlIEZyZWUgU29mdHdhcmUgRm91bmRhdGlvbiwgdXNlIHRoZQpHTlUgR2VuZXJhbCBQdWJsaWMgTGljZW5zZSBmb3IgbW9zdCBvZiBvdXIgc29mdHdhcmU7IGl0IGFwcGxpZXMgYWxzbyB0bwphbnkgb3RoZXIgd29yayByZWxlYXNlZCB0aGlzIHdheSBieSBpdHMgYXV0aG9ycy4gIFlvdSBjYW4gYXBwbHkgaXQgdG8KeW91ciBwcm9ncmFtcywgdG9vLgoKICBXaGVuIHdlIHNwZWFrIG9mIGZyZWUgc29mdHdhcmUsIHdlIGFyZSByZWZlcnJpbmcgdG8gZnJlZWRvbSwgbm90CnByaWNlLiAgT3VyIEdlbmVyYWwgUHVibGljIExpY2Vuc2VzIGFyZSBkZXNpZ25lZCB0byBtYWtlIHN1cmUgdGhhdCB5b3UKaGF2ZSB0aGUgZnJlZWRvbSB0byBkaXN0cmlidXRlIGNvcGllcyBvZiBmcmVlIHNvZnR3YXJlIChhbmQgY2hhcmdlIGZvcgp0aGVtIGlmIHlvdSB3aXNoKSwgdGhhdCB5b3UgcmVjZWl2ZSBzb3VyY2UgY29kZSBvciBjYW4gZ2V0IGl0IGlmIHlvdQp3YW50IGl0LCB0aGF0IHlvdSBjYW4gY2hhbmdlIHRoZSBzb2Z0d2FyZSBvciB1c2UgcGllY2VzIG9mIGl0IGluIG5ldwpmcmVlIHByb2dyYW1zLCBhbmQgdGhhdCB5b3Uga25vdyB5b3UgY2FuIGRvIHRoZXNlIHRoaW5ncy4KCiAgVG8gcHJvdGVjdCB5b3VyIHJpZ2h0cywgd2UgbmVlZCB0byBwcmV2ZW50IG90aGVycyBmcm9tIGRlbnlpbmcgeW91CnRoZXNlIHJpZ2h0cyBvciBhc2tpbmcgeW91IHRvIHN1cnJlbmRlciB0aGUgcmlnaHRzLiAgVGhlcmVmb3JlLCB5b3UgaGF2ZQpjZXJ0YWluIHJlc3BvbnNpYmlsaXRpZXMgaWYgeW91IGRpc3RyaWJ1dGUgY29waWVzIG9mIHRoZSBzb2Z0d2FyZSwgb3IgaWYKeW91IG1vZGlmeSBpdDogcmVzcG9uc2liaWxpdGllcyB0byByZXNwZWN0IHRoZSBmcmVlZG9tIG9mIG90aGVycy4KCiAgRm9yIGV4YW1wbGUsIGlmIHlvdSBkaXN0cmlidXRlIGNvcGllcyBvZiBzdWNoIGEgcHJvZ3JhbSwgd2hldGhlcgpncmF0aXMgb3IgZm9yIGEgZmVlLCB5b3UgbXVzdCBwYXNzIG9uIHRvIHRoZSByZWNpcGllbnRzIHRoZSBzYW1lCmZyZWVkb21zIHRoYXQgeW91IHJlY2VpdmVkLiAgWW91IG11c3QgbWFrZSBzdXJlIHRoYXQgdGhleSwgdG9vLCByZWNlaXZlCm9yIGNhbiBnZXQgdGhlIHNvdXJjZSBjb2RlLiAgQW5kIHlvdSBtdXN0IHNob3cgdGhlbSB0aGVzZSB0ZXJtcyBzbyB0aGV5Cmtub3cgdGhlaXIgcmlnaHRzLgoKICBEZXZlbG9wZXJzIHRoYXQgdXNlIHRoZSBHTlUgR1BMIHByb3RlY3QgeW91ciByaWdodHMgd2l0aCB0d28gc3RlcHM6CigxKSBhc3NlcnQgY29weXJpZ2h0IG9uIHRoZSBzb2Z0d2FyZSwgYW5kICgyKSBvZmZlciB5b3UgdGhpcyBMaWNlbnNlCmdpdmluZyB5b3UgbGVnYWwgcGVybWlzc2lvbiB0byBjb3B5LCBkaXN0cmlidXRlIGFuZC9vciBtb2RpZnkgaXQuCgogIEZvciB0aGUgZGV2ZWxvcGVycycgYW5kIGF1dGhvcnMnIHByb3RlY3Rpb24sIHRoZSBHUEwgY2xlYXJseSBleHBsYWlucwp0aGF0IHRoZXJlIGlzIG5vIHdhcnJhbnR5IGZvciB0aGlzIGZyZWUgc29mdHdhcmUuICBGb3IgYm90aCB1c2VycycgYW5kCmF1dGhvcnMnIHNha2UsIHRoZSBHUEwgcmVxdWlyZXMgdGhhdCBtb2RpZmllZCB2ZXJzaW9ucyBiZSBtYXJrZWQgYXMKY2hhbmdlZCwgc28gdGhhdCB0aGVpciBwcm9ibGVtcyB3aWxsIG5vdCBiZSBhdHRyaWJ1dGVkIGVycm9uZW91c2x5IHRvCmF1dGhvcnMgb2YgcHJldmlvdXMgdmVyc2lvbnMuCgogIFNvbWUgZGV2aWNlcyBhcmUgZGVzaWduZWQgdG8gZGVueSB1c2VycyBhY2Nlc3MgdG8gaW5zdGFsbCBvciBydW4KbW9kaWZpZWQgdmVyc2lvbnMgb2YgdGhlIHNvZnR3YXJlIGluc2lkZSB0aGVtLCBhbHRob3VnaCB0aGUgbWFudWZhY3R1cmVyCmNhbiBkbyBzby4gIFRoaXMgaXMgZnVuZGFtZW50YWxseSBpbmNvbXBhdGlibGUgd2l0aCB0aGUgYWltIG9mCnByb3RlY3RpbmcgdXNlcnMnIGZyZWVkb20gdG8gY2hhbmdlIHRoZSBzb2Z0d2FyZS4gIFRoZSBzeXN0ZW1hdGljCnBhdHRlcm4gb2Ygc3VjaCBhYnVzZSBvY2N1cnMgaW4gdGhlIGFyZWEgb2YgcHJvZHVjdHMgZm9yIGluZGl2aWR1YWxzIHRvCnVzZSwgd2hpY2ggaXMgcHJlY2lzZWx5IHdoZXJlIGl0IGlzIG1vc3QgdW5hY2NlcHRhYmxlLiAgVGhlcmVmb3JlLCB3ZQpoYXZlIGRlc2lnbmVkIHRoaXMgdmVyc2lvbiBvZiB0aGUgR1BMIHRvIHByb2hpYml0IHRoZSBwcmFjdGljZSBmb3IgdGhvc2UKcHJvZHVjdHMuICBJZiBzdWNoIHByb2JsZW1zIGFyaXNlIHN1YnN0YW50aWFsbHkgaW4gb3RoZXIgZG9tYWlucywgd2UKc3RhbmQgcmVhZHkgdG8gZXh0ZW5kIHRoaXMgcHJvdmlzaW9uIHRvIHRob3NlIGRvbWFpbnMgaW4gZnV0dXJlIHZlcnNpb25zCm9mIHRoZSBHUEwsIGFzIG5lZWRlZCB0byBwcm90ZWN0IHRoZSBmcmVlZG9tIG9mIHVzZXJzLgoKICBGaW5hbGx5LCBldmVyeSBwcm9ncmFtIGlzIHRocmVhdGVuZWQgY29uc3RhbnRseSBieSBzb2Z0d2FyZSBwYXRlbnRzLgpTdGF0ZXMgc2hvdWxkIG5vdCBhbGxvdyBwYXRlbnRzIHRvIHJlc3RyaWN0IGRldmVsb3BtZW50IGFuZCB1c2Ugb2YKc29mdHdhcmUgb24gZ2VuZXJhbC1wdXJwb3NlIGNvbXB1dGVycywgYnV0IGluIHRob3NlIHRoYXQgZG8sIHdlIHdpc2ggdG8KYXZvaWQgdGhlIHNwZWNpYWwgZGFuZ2VyIHRoYXQgcGF0ZW50cyBhcHBsaWVkIHRvIGEgZnJlZSBwcm9ncmFtIGNvdWxkCm1ha2UgaXQgZWZmZWN0aXZlbHkgcHJvcHJpZXRhcnkuICBUbyBwcmV2ZW50IHRoaXMsIHRoZSBHUEwgYXNzdXJlcyB0aGF0CnBhdGVudHMgY2Fubm90IGJlIHVzZWQgdG8gcmVuZGVyIHRoZSBwcm9ncmFtIG5vbi1mcmVlLgoKICBUaGUgcHJlY2lzZSB0ZXJtcyBhbmQgY29uZGl0aW9ucyBmb3IgY29weWluZywgZGlzdHJpYnV0aW9uIGFuZAptb2RpZmljYXRpb24gZm9sbG93LgoKICAgICAgICAgICAgICAgICAgICAgICBURVJNUyBBTkQgQ09ORElUSU9OUwoKICAwLiBEZWZpbml0aW9ucy4KCiAgIlRoaXMgTGljZW5zZSIgcmVmZXJzIHRvIHZlcnNpb24gMyBvZiB0aGUgR05VIEdlbmVyYWwgUHVibGljIExpY2Vuc2UuCgogICJDb3B5cmlnaHQiIGFsc28gbWVhbnMgY29weXJpZ2h0LWxpa2UgbGF3cyB0aGF0IGFwcGx5IHRvIG90aGVyIGtpbmRzIG9mCndvcmtzLCBzdWNoIGFzIHNlbWljb25kdWN0b3IgbWFza3MuCgogICJUaGUgUHJvZ3JhbSIgcmVmZXJzIHRvIGFueSBjb3B5cmlnaHRhYmxlIHdvcmsgbGljZW5zZWQgdW5kZXIgdGhpcwpMaWNlbnNlLiAgRWFjaCBsaWNlbnNlZSBpcyBhZGRyZXNzZWQgYXMgInlvdSIuICAiTGljZW5zZWVzIiBhbmQKInJlY2lwaWVudHMiIG1heSBiZSBpbmRpdmlkdWFscyBvciBvcmdhbml6YXRpb25zLgoKICBUbyAibW9kaWZ5IiBhIHdvcmsgbWVhbnMgdG8gY29weSBmcm9tIG9yIGFkYXB0IGFsbCBvciBwYXJ0IG9mIHRoZSB3b3JrCmluIGEgZmFzaGlvbiByZXF1aXJpbmcgY29weXJpZ2h0IHBlcm1pc3Npb24sIG90aGVyIHRoYW4gdGhlIG1ha2luZyBvZiBhbgpleGFjdCBjb3B5LiAgVGhlIHJlc3VsdGluZyB3b3JrIGlzIGNhbGxlZCBhICJtb2RpZmllZCB2ZXJzaW9uIiBvZiB0aGUKZWFybGllciB3b3JrIG9yIGEgd29yayAiYmFzZWQgb24iIHRoZSBlYXJsaWVyIHdvcmsuCgogIEEgImNvdmVyZWQgd29yayIgbWVhbnMgZWl0aGVyIHRoZSB1bm1vZGlmaWVkIFByb2dyYW0gb3IgYSB3b3JrIGJhc2VkCm9uIHRoZSBQcm9ncmFtLgoKICBUbyAicHJvcGFnYXRlIiBhIHdvcmsgbWVhbnMgdG8gZG8gYW55dGhpbmcgd2l0aCBpdCB0aGF0LCB3aXRob3V0CnBlcm1pc3Npb24sIHdvdWxkIG1ha2UgeW91IGRpcmVjdGx5IG9yIHNlY29uZGFyaWx5IGxpYWJsZSBmb3IKaW5mcmluZ2VtZW50IHVuZGVyIGFwcGxpY2FibGUgY29weXJpZ2h0IGxhdywgZXhjZXB0IGV4ZWN1dGluZyBpdCBvbiBhCmNvbXB1dGVyIG9yIG1vZGlmeWluZyBhIHByaXZhdGUgY29weS4gIFByb3BhZ2F0aW9uIGluY2x1ZGVzIGNvcHlpbmcsCmRpc3RyaWJ1dGlvbiAod2l0aCBvciB3aXRob3V0IG1vZGlmaWNhdGlvbiksIG1ha2luZyBhdmFpbGFibGUgdG8gdGhlCnB1YmxpYywgYW5kIGluIHNvbWUgY291bnRyaWVzIG90aGVyIGFjdGl2aXRpZXMgYXMgd2VsbC4KCiAgVG8gImNvbnZleSIgYSB3b3JrIG1lYW5zIGFueSBraW5kIG9mIHByb3BhZ2F0aW9uIHRoYXQgZW5hYmxlcyBvdGhlcgpwYXJ0aWVzIHRvIG1ha2Ugb3IgcmVjZWl2ZSBjb3BpZXMuICBNZXJlIGludGVyYWN0aW9uIHdpdGggYSB1c2VyIHRocm91Z2gKYSBjb21wdXRlciBuZXR3b3JrLCB3aXRoIG5vIHRyYW5zZmVyIG9mIGEgY29weSwgaXMgbm90IGNvbnZleWluZy4KCiAgQW4gaW50ZXJhY3RpdmUgdXNlciBpbnRlcmZhY2UgZGlzcGxheXMgIkFwcHJvcHJpYXRlIExlZ2FsIE5vdGljZXMiCnRvIHRoZSBleHRlbnQgdGhhdCBpdCBpbmNsdWRlcyBhIGNvbnZlbmllbnQgYW5kIHByb21pbmVudGx5IHZpc2libGUKZmVhdHVyZSB0aGF0ICgxKSBkaXNwbGF5cyBhbiBhcHByb3ByaWF0ZSBjb3B5cmlnaHQgbm90aWNlLCBhbmQgKDIpCnRlbGxzIHRoZSB1c2VyIHRoYXQgdGhlcmUgaXMgbm8gd2FycmFudHkgZm9yIHRoZSB3b3JrIChleGNlcHQgdG8gdGhlCmV4dGVudCB0aGF0IHdhcnJhbnRpZXMgYXJlIHByb3ZpZGVkKSwgdGhhdCBsaWNlbnNlZXMgbWF5IGNvbnZleSB0aGUKd29yayB1bmRlciB0aGlzIExpY2Vuc2UsIGFuZCBob3cgdG8gdmlldyBhIGNvcHkgb2YgdGhpcyBMaWNlbnNlLiAgSWYKdGhlIGludGVyZmFjZSBwcmVzZW50cyBhIGxpc3Qgb2YgdXNlciBjb21tYW5kcyBvciBvcHRpb25zLCBzdWNoIGFzIGEKbWVudSwgYSBwcm9taW5lbnQgaXRlbSBpbiB0aGUgbGlzdCBtZWV0cyB0aGlzIGNyaXRlcmlvbi4KCiAgMS4gU291cmNlIENvZGUuCgogIFRoZSAic291cmNlIGNvZGUiIGZvciBhIHdvcmsgbWVhbnMgdGhlIHByZWZlcnJlZCBmb3JtIG9mIHRoZSB3b3JrCmZvciBtYWtpbmcgbW9kaWZpY2F0aW9ucyB0byBpdC4gICJPYmplY3QgY29kZSIgbWVhbnMgYW55IG5vbi1zb3VyY2UKZm9ybSBvZiBhIHdvcmsuCgogIEEgIlN0YW5kYXJkIEludGVyZmFjZSIgbWVhbnMgYW4gaW50ZXJmYWNlIHRoYXQgZWl0aGVyIGlzIGFuIG9mZmljaWFsCnN0YW5kYXJkIGRlZmluZWQgYnkgYSByZWNvZ25pemVkIHN0YW5kYXJkcyBib2R5LCBvciwgaW4gdGhlIGNhc2Ugb2YKaW50ZXJmYWNlcyBzcGVjaWZpZWQgZm9yIGEgcGFydGljdWxhciBwcm9ncmFtbWluZyBsYW5ndWFnZSwgb25lIHRoYXQKaXMgd2lkZWx5IHVzZWQgYW1vbmcgZGV2ZWxvcGVycyB3b3JraW5nIGluIHRoYXQgbGFuZ3VhZ2UuCgogIFRoZSAiU3lzdGVtIExpYnJhcmllcyIgb2YgYW4gZXhlY3V0YWJsZSB3b3JrIGluY2x1ZGUgYW55dGhpbmcsIG90aGVyCnRoYW4gdGhlIHdvcmsgYXMgYSB3aG9sZSwgdGhhdCAoYSkgaXMgaW5jbHVkZWQgaW4gdGhlIG5vcm1hbCBmb3JtIG9mCnBhY2thZ2luZyBhIE1ham9yIENvbXBvbmVudCwgYnV0IHdoaWNoIGlzIG5vdCBwYXJ0IG9mIHRoYXQgTWFqb3IKQ29tcG9uZW50LCBhbmQgKGIpIHNlcnZlcyBvbmx5IHRvIGVuYWJsZSB1c2Ugb2YgdGhlIHdvcmsgd2l0aCB0aGF0Ck1ham9yIENvbXBvbmVudCwgb3IgdG8gaW1wbGVtZW50IGEgU3RhbmRhcmQgSW50ZXJmYWNlIGZvciB3aGljaCBhbgppbXBsZW1lbnRhdGlvbiBpcyBhdmFpbGFibGUgdG8gdGhlIHB1YmxpYyBpbiBzb3VyY2UgY29kZSBmb3JtLiAgQQoiTWFqb3IgQ29tcG9uZW50IiwgaW4gdGhpcyBjb250ZXh0LCBtZWFucyBhIG1ham9yIGVzc2VudGlhbCBjb21wb25lbnQKKGtlcm5lbCwgd2luZG93IHN5c3RlbSwgYW5kIHNvIG9uKSBvZiB0aGUgc3BlY2lmaWMgb3BlcmF0aW5nIHN5c3RlbQooaWYgYW55KSBvbiB3aGljaCB0aGUgZXhlY3V0YWJsZSB3b3JrIHJ1bnMsIG9yIGEgY29tcGlsZXIgdXNlZCB0bwpwcm9kdWNlIHRoZSB3b3JrLCBvciBhbiBvYmplY3QgY29kZSBpbnRlcnByZXRlciB1c2VkIHRvIHJ1biBpdC4KCiAgVGhlICJDb3JyZXNwb25kaW5nIFNvdXJjZSIgZm9yIGEgd29yayBpbiBvYmplY3QgY29kZSBmb3JtIG1lYW5zIGFsbAp0aGUgc291cmNlIGNvZGUgbmVlZGVkIHRvIGdlbmVyYXRlLCBpbnN0YWxsLCBhbmQgKGZvciBhbiBleGVjdXRhYmxlCndvcmspIHJ1biB0aGUgb2JqZWN0IGNvZGUgYW5kIHRvIG1vZGlmeSB0aGUgd29yaywgaW5jbHVkaW5nIHNjcmlwdHMgdG8KY29udHJvbCB0aG9zZSBhY3Rpdml0aWVzLiAgSG93ZXZlciwgaXQgZG9lcyBub3QgaW5jbHVkZSB0aGUgd29yaydzClN5c3RlbSBMaWJyYXJpZXMsIG9yIGdlbmVyYWwtcHVycG9zZSB0b29scyBvciBnZW5lcmFsbHkgYXZhaWxhYmxlIGZyZWUKcHJvZ3JhbXMgd2hpY2ggYXJlIHVzZWQgdW5tb2RpZmllZCBpbiBwZXJmb3JtaW5nIHRob3NlIGFjdGl2aXRpZXMgYnV0CndoaWNoIGFyZSBub3QgcGFydCBvZiB0aGUgd29yay4gIEZvciBleGFtcGxlLCBDb3JyZXNwb25kaW5nIFNvdXJjZQppbmNsdWRlcyBpbnRlcmZhY2UgZGVmaW5pdGlvbiBmaWxlcyBhc3NvY2lhdGVkIHdpdGggc291cmNlIGZpbGVzIGZvcgp0aGUgd29yaywgYW5kIHRoZSBzb3VyY2UgY29kZSBmb3Igc2hhcmVkIGxpYnJhcmllcyBhbmQgZHluYW1pY2FsbHkKbGlua2VkIHN1YnByb2dyYW1zIHRoYXQgdGhlIHdvcmsgaXMgc3BlY2lmaWNhbGx5IGRlc2lnbmVkIHRvIHJlcXVpcmUsCnN1Y2ggYXMgYnkgaW50aW1hdGUgZGF0YSBjb21tdW5pY2F0aW9uIG9yIGNvbnRyb2wgZmxvdyBiZXR3ZWVuIHRob3NlCnN1YnByb2dyYW1zIGFuZCBvdGhlciBwYXJ0cyBvZiB0aGUgd29yay4KCiAgVGhlIENvcnJlc3BvbmRpbmcgU291cmNlIG5lZWQgbm90IGluY2x1ZGUgYW55dGhpbmcgdGhhdCB1c2VycwpjYW4gcmVnZW5lcmF0ZSBhdXRvbWF0aWNhbGx5IGZyb20gb3RoZXIgcGFydHMgb2YgdGhlIENvcnJlc3BvbmRpbmcKU291cmNlLgoKICBUaGUgQ29ycmVzcG9uZGluZyBTb3VyY2UgZm9yIGEgd29yayBpbiBzb3VyY2UgY29kZSBmb3JtIGlzIHRoYXQKc2FtZSB3b3JrLgoKICAyLiBCYXNpYyBQZXJtaXNzaW9ucy4KCiAgQWxsIHJpZ2h0cyBncmFudGVkIHVuZGVyIHRoaXMgTGljZW5zZSBhcmUgZ3JhbnRlZCBmb3IgdGhlIHRlcm0gb2YKY29weXJpZ2h0IG9uIHRoZSBQcm9ncmFtLCBhbmQgYXJlIGlycmV2b2NhYmxlIHByb3ZpZGVkIHRoZSBzdGF0ZWQKY29uZGl0aW9ucyBhcmUgbWV0LiAgVGhpcyBMaWNlbnNlIGV4cGxpY2l0bHkgYWZmaXJtcyB5b3VyIHVubGltaXRlZApwZXJtaXNzaW9uIHRvIHJ1biB0aGUgdW5tb2RpZmllZCBQcm9ncmFtLiAgVGhlIG91dHB1dCBmcm9tIHJ1bm5pbmcgYQpjb3ZlcmVkIHdvcmsgaXMgY292ZXJlZCBieSB0aGlzIExpY2Vuc2Ugb25seSBpZiB0aGUgb3V0cHV0LCBnaXZlbiBpdHMKY29udGVudCwgY29uc3RpdHV0ZXMgYSBjb3ZlcmVkIHdvcmsuICBUaGlzIExpY2Vuc2UgYWNrbm93bGVkZ2VzIHlvdXIKcmlnaHRzIG9mIGZhaXIgdXNlIG9yIG90aGVyIGVxdWl2YWxlbnQsIGFzIHByb3ZpZGVkIGJ5IGNvcHlyaWdodCBsYXcuCgogIFlvdSBtYXkgbWFrZSwgcnVuIGFuZCBwcm9wYWdhdGUgY292ZXJlZCB3b3JrcyB0aGF0IHlvdSBkbyBub3QKY29udmV5LCB3aXRob3V0IGNvbmRpdGlvbnMgc28gbG9uZyBhcyB5b3VyIGxpY2Vuc2Ugb3RoZXJ3aXNlIHJlbWFpbnMKaW4gZm9yY2UuICBZb3UgbWF5IGNvbnZleSBjb3ZlcmVkIHdvcmtzIHRvIG90aGVycyBmb3IgdGhlIHNvbGUgcHVycG9zZQpvZiBoYXZpbmcgdGhlbSBtYWtlIG1vZGlmaWNhdGlvbnMgZXhjbHVzaXZlbHkgZm9yIHlvdSwgb3IgcHJvdmlkZSB5b3UKd2l0aCBmYWNpbGl0aWVzIGZvciBydW5uaW5nIHRob3NlIHdvcmtzLCBwcm92aWRlZCB0aGF0IHlvdSBjb21wbHkgd2l0aAp0aGUgdGVybXMgb2YgdGhpcyBMaWNlbnNlIGluIGNvbnZleWluZyBhbGwgbWF0ZXJpYWwgZm9yIHdoaWNoIHlvdSBkbwpub3QgY29udHJvbCBjb3B5cmlnaHQuICBUaG9zZSB0aHVzIG1ha2luZyBvciBydW5uaW5nIHRoZSBjb3ZlcmVkIHdvcmtzCmZvciB5b3UgbXVzdCBkbyBzbyBleGNsdXNpdmVseSBvbiB5b3VyIGJlaGFsZiwgdW5kZXIgeW91ciBkaXJlY3Rpb24KYW5kIGNvbnRyb2wsIG9uIHRlcm1zIHRoYXQgcHJvaGliaXQgdGhlbSBmcm9tIG1ha2luZyBhbnkgY29waWVzIG9mCnlvdXIgY29weXJpZ2h0ZWQgbWF0ZXJpYWwgb3V0c2lkZSB0aGVpciByZWxhdGlvbnNoaXAgd2l0aCB5b3UuCgogIENvbnZleWluZyB1bmRlciBhbnkgb3RoZXIgY2lyY3Vtc3RhbmNlcyBpcyBwZXJtaXR0ZWQgc29sZWx5IHVuZGVyCnRoZSBjb25kaXRpb25zIHN0YXRlZCBiZWxvdy4gIFN1YmxpY2Vuc2luZyBpcyBub3QgYWxsb3dlZDsgc2VjdGlvbiAxMAptYWtlcyBpdCB1bm5lY2Vzc2FyeS4KCiAgMy4gUHJvdGVjdGluZyBVc2VycycgTGVnYWwgUmlnaHRzIEZyb20gQW50aS1DaXJjdW12ZW50aW9uIExhdy4KCiAgTm8gY292ZXJlZCB3b3JrIHNoYWxsIGJlIGRlZW1lZCBwYXJ0IG9mIGFuIGVmZmVjdGl2ZSB0ZWNobm9sb2dpY2FsCm1lYXN1cmUgdW5kZXIgYW55IGFwcGxpY2FibGUgbGF3IGZ1bGZpbGxpbmcgb2JsaWdhdGlvbnMgdW5kZXIgYXJ0aWNsZQoxMSBvZiB0aGUgV0lQTyBjb3B5cmlnaHQgdHJlYXR5IGFkb3B0ZWQgb24gMjAgRGVjZW1iZXIgMTk5Niwgb3IKc2ltaWxhciBsYXdzIHByb2hpYml0aW5nIG9yIHJlc3RyaWN0aW5nIGNpcmN1bXZlbnRpb24gb2Ygc3VjaAptZWFzdXJlcy4KCiAgV2hlbiB5b3UgY29udmV5IGEgY292ZXJlZCB3b3JrLCB5b3Ugd2FpdmUgYW55IGxlZ2FsIHBvd2VyIHRvIGZvcmJpZApjaXJjdW12ZW50aW9uIG9mIHRlY2hub2xvZ2ljYWwgbWVhc3VyZXMgdG8gdGhlIGV4dGVudCBzdWNoIGNpcmN1bXZlbnRpb24KaXMgZWZmZWN0ZWQgYnkgZXhlcmNpc2luZyByaWdodHMgdW5kZXIgdGhpcyBMaWNlbnNlIHdpdGggcmVzcGVjdCB0bwp0aGUgY292ZXJlZCB3b3JrLCBhbmQgeW91IGRpc2NsYWltIGFueSBpbnRlbnRpb24gdG8gbGltaXQgb3BlcmF0aW9uIG9yCm1vZGlmaWNhdGlvbiBvZiB0aGUgd29yayBhcyBhIG1lYW5zIG9mIGVuZm9yY2luZywgYWdhaW5zdCB0aGUgd29yaydzCnVzZXJzLCB5b3VyIG9yIHRoaXJkIHBhcnRpZXMnIGxlZ2FsIHJpZ2h0cyB0byBmb3JiaWQgY2lyY3VtdmVudGlvbiBvZgp0ZWNobm9sb2dpY2FsIG1lYXN1cmVzLgoKICA0LiBDb252ZXlpbmcgVmVyYmF0aW0gQ29waWVzLgoKICBZb3UgbWF5IGNvbnZleSB2ZXJiYXRpbSBjb3BpZXMgb2YgdGhlIFByb2dyYW0ncyBzb3VyY2UgY29kZSBhcyB5b3UKcmVjZWl2ZSBpdCwgaW4gYW55IG1lZGl1bSwgcHJvdmlkZWQgdGhhdCB5b3UgY29uc3BpY3VvdXNseSBhbmQKYXBwcm9wcmlhdGVseSBwdWJsaXNoIG9uIGVhY2ggY29weSBhbiBhcHByb3ByaWF0ZSBjb3B5cmlnaHQgbm90aWNlOwprZWVwIGludGFjdCBhbGwgbm90aWNlcyBzdGF0aW5nIHRoYXQgdGhpcyBMaWNlbnNlIGFuZCBhbnkKbm9uLXBlcm1pc3NpdmUgdGVybXMgYWRkZWQgaW4gYWNjb3JkIHdpdGggc2VjdGlvbiA3IGFwcGx5IHRvIHRoZSBjb2RlOwprZWVwIGludGFjdCBhbGwgbm90aWNlcyBvZiB0aGUgYWJzZW5jZSBvZiBhbnkgd2FycmFudHk7IGFuZCBnaXZlIGFsbApyZWNpcGllbnRzIGEgY29weSBvZiB0aGlzIExpY2Vuc2UgYWxvbmcgd2l0aCB0aGUgUHJvZ3JhbS4KCiAgWW91IG1heSBjaGFyZ2UgYW55IHByaWNlIG9yIG5vIHByaWNlIGZvciBlYWNoIGNvcHkgdGhhdCB5b3UgY29udmV5LAphbmQgeW91IG1heSBvZmZlciBzdXBwb3J0IG9yIHdhcnJhbnR5IHByb3RlY3Rpb24gZm9yIGEgZmVlLgoKICA1LiBDb252ZXlpbmcgTW9kaWZpZWQgU291cmNlIFZlcnNpb25zLgoKICBZb3UgbWF5IGNvbnZleSBhIHdvcmsgYmFzZWQgb24gdGhlIFByb2dyYW0sIG9yIHRoZSBtb2RpZmljYXRpb25zIHRvCnByb2R1Y2UgaXQgZnJvbSB0aGUgUHJvZ3JhbSwgaW4gdGhlIGZvcm0gb2Ygc291cmNlIGNvZGUgdW5kZXIgdGhlCnRlcm1zIG9mIHNlY3Rpb24gNCwgcHJvdmlkZWQgdGhhdCB5b3UgYWxzbyBtZWV0IGFsbCBvZiB0aGVzZSBjb25kaXRpb25zOgoKICAgIGEpIFRoZSB3b3JrIG11c3QgY2FycnkgcHJvbWluZW50IG5vdGljZXMgc3RhdGluZyB0aGF0IHlvdSBtb2RpZmllZAogICAgaXQsIGFuZCBnaXZpbmcgYSByZWxldmFudCBkYXRlLgoKICAgIGIpIFRoZSB3b3JrIG11c3QgY2FycnkgcHJvbWluZW50IG5vdGljZXMgc3RhdGluZyB0aGF0IGl0IGlzCiAgICByZWxlYXNlZCB1bmRlciB0aGlzIExpY2Vuc2UgYW5kIGFueSBjb25kaXRpb25zIGFkZGVkIHVuZGVyIHNlY3Rpb24KICAgIDcuICBUaGlzIHJlcXVpcmVtZW50IG1vZGlmaWVzIHRoZSByZXF1aXJlbWVudCBpbiBzZWN0aW9uIDQgdG8KICAgICJrZWVwIGludGFjdCBhbGwgbm90aWNlcyIuCgogICAgYykgWW91IG11c3QgbGljZW5zZSB0aGUgZW50aXJlIHdvcmssIGFzIGEgd2hvbGUsIHVuZGVyIHRoaXMKICAgIExpY2Vuc2UgdG8gYW55b25lIHdobyBjb21lcyBpbnRvIHBvc3Nlc3Npb24gb2YgYSBjb3B5LiAgVGhpcwogICAgTGljZW5zZSB3aWxsIHRoZXJlZm9yZSBhcHBseSwgYWxvbmcgd2l0aCBhbnkgYXBwbGljYWJsZSBzZWN0aW9uIDcKICAgIGFkZGl0aW9uYWwgdGVybXMsIHRvIHRoZSB3aG9sZSBvZiB0aGUgd29yaywgYW5kIGFsbCBpdHMgcGFydHMsCiAgICByZWdhcmRsZXNzIG9mIGhvdyB0aGV5IGFyZSBwYWNrYWdlZC4gIFRoaXMgTGljZW5zZSBnaXZlcyBubwogICAgcGVybWlzc2lvbiB0byBsaWNlbnNlIHRoZSB3b3JrIGluIGFueSBvdGhlciB3YXksIGJ1dCBpdCBkb2VzIG5vdAogICAgaW52YWxpZGF0ZSBzdWNoIHBlcm1pc3Npb24gaWYgeW91IGhhdmUgc2VwYXJhdGVseSByZWNlaXZlZCBpdC4KCiAgICBkKSBJZiB0aGUgd29yayBoYXMgaW50ZXJhY3RpdmUgdXNlciBpbnRlcmZhY2VzLCBlYWNoIG11c3QgZGlzcGxheQogICAgQXBwcm9wcmlhdGUgTGVnYWwgTm90aWNlczsgaG93ZXZlciwgaWYgdGhlIFByb2dyYW0gaGFzIGludGVyYWN0aXZlCiAgICBpbnRlcmZhY2VzIHRoYXQgZG8gbm90IGRpc3BsYXkgQXBwcm9wcmlhdGUgTGVnYWwgTm90aWNlcywgeW91cgogICAgd29yayBuZWVkIG5vdCBtYWtlIHRoZW0gZG8gc28uCgogIEEgY29tcGlsYXRpb24gb2YgYSBjb3ZlcmVkIHdvcmsgd2l0aCBvdGhlciBzZXBhcmF0ZSBhbmQgaW5kZXBlbmRlbnQKd29ya3MsIHdoaWNoIGFyZSBub3QgYnkgdGhlaXIgbmF0dXJlIGV4dGVuc2lvbnMgb2YgdGhlIGNvdmVyZWQgd29yaywKYW5kIHdoaWNoIGFyZSBub3QgY29tYmluZWQgd2l0aCBpdCBzdWNoIGFzIHRvIGZvcm0gYSBsYXJnZXIgcHJvZ3JhbSwKaW4gb3Igb24gYSB2b2x1bWUgb2YgYSBzdG9yYWdlIG9yIGRpc3RyaWJ1dGlvbiBtZWRpdW0sIGlzIGNhbGxlZCBhbgoiYWdncmVnYXRlIiBpZiB0aGUgY29tcGlsYXRpb24gYW5kIGl0cyByZXN1bHRpbmcgY29weXJpZ2h0IGFyZSBub3QKdXNlZCB0byBsaW1pdCB0aGUgYWNjZXNzIG9yIGxlZ2FsIHJpZ2h0cyBvZiB0aGUgY29tcGlsYXRpb24ncyB1c2VycwpiZXlvbmQgd2hhdCB0aGUgaW5kaXZpZHVhbCB3b3JrcyBwZXJtaXQuICBJbmNsdXNpb24gb2YgYSBjb3ZlcmVkIHdvcmsKaW4gYW4gYWdncmVnYXRlIGRvZXMgbm90IGNhdXNlIHRoaXMgTGljZW5zZSB0byBhcHBseSB0byB0aGUgb3RoZXIKcGFydHMgb2YgdGhlIGFnZ3JlZ2F0ZS4KCiAgNi4gQ29udmV5aW5nIE5vbi1Tb3VyY2UgRm9ybXMuCgogIFlvdSBtYXkgY29udmV5IGEgY292ZXJlZCB3b3JrIGluIG9iamVjdCBjb2RlIGZvcm0gdW5kZXIgdGhlIHRlcm1zCm9mIHNlY3Rpb25zIDQgYW5kIDUsIHByb3ZpZGVkIHRoYXQgeW91IGFsc28gY29udmV5IHRoZQptYWNoaW5lLXJlYWRhYmxlIENvcnJlc3BvbmRpbmcgU291cmNlIHVuZGVyIHRoZSB0ZXJtcyBvZiB0aGlzIExpY2Vuc2UsCmluIG9uZSBvZiB0aGVzZSB3YXlzOgoKICAgIGEpIENvbnZleSB0aGUgb2JqZWN0IGNvZGUgaW4sIG9yIGVtYm9kaWVkIGluLCBhIHBoeXNpY2FsIHByb2R1Y3QKICAgIChpbmNsdWRpbmcgYSBwaHlzaWNhbCBkaXN0cmlidXRpb24gbWVkaXVtKSwgYWNjb21wYW5pZWQgYnkgdGhlCiAgICBDb3JyZXNwb25kaW5nIFNvdXJjZSBmaXhlZCBvbiBhIGR1cmFibGUgcGh5c2ljYWwgbWVkaXVtCiAgICBjdXN0b21hcmlseSB1c2VkIGZvciBzb2Z0d2FyZSBpbnRlcmNoYW5nZS4KCiAgICBiKSBDb252ZXkgdGhlIG9iamVjdCBjb2RlIGluLCBvciBlbWJvZGllZCBpbiwgYSBwaHlzaWNhbCBwcm9kdWN0CiAgICAoaW5jbHVkaW5nIGEgcGh5c2ljYWwgZGlzdHJpYnV0aW9uIG1lZGl1bSksIGFjY29tcGFuaWVkIGJ5IGEKICAgIHdyaXR0ZW4gb2ZmZXIsIHZhbGlkIGZvciBhdCBsZWFzdCB0aHJlZSB5ZWFycyBhbmQgdmFsaWQgZm9yIGFzCiAgICBsb25nIGFzIHlvdSBvZmZlciBzcGFyZSBwYXJ0cyBvciBjdXN0b21lciBzdXBwb3J0IGZvciB0aGF0IHByb2R1Y3QKICAgIG1vZGVsLCB0byBnaXZlIGFueW9uZSB3aG8gcG9zc2Vzc2VzIHRoZSBvYmplY3QgY29kZSBlaXRoZXIgKDEpIGEKICAgIGNvcHkgb2YgdGhlIENvcnJlc3BvbmRpbmcgU291cmNlIGZvciBhbGwgdGhlIHNvZnR3YXJlIGluIHRoZQogICAgcHJvZHVjdCB0aGF0IGlzIGNvdmVyZWQgYnkgdGhpcyBMaWNlbnNlLCBvbiBhIGR1cmFibGUgcGh5c2ljYWwKICAgIG1lZGl1bSBjdXN0b21hcmlseSB1c2VkIGZvciBzb2Z0d2FyZSBpbnRlcmNoYW5nZSwgZm9yIGEgcHJpY2Ugbm8KICAgIG1vcmUgdGhhbiB5b3VyIHJlYXNvbmFibGUgY29zdCBvZiBwaHlzaWNhbGx5IHBlcmZvcm1pbmcgdGhpcwogICAgY29udmV5aW5nIG9mIHNvdXJjZSwgb3IgKDIpIGFjY2VzcyB0byBjb3B5IHRoZQogICAgQ29ycmVzcG9uZGluZyBTb3VyY2UgZnJvbSBhIG5ldHdvcmsgc2VydmVyIGF0IG5vIGNoYXJnZS4KCiAgICBjKSBDb252ZXkgaW5kaXZpZHVhbCBjb3BpZXMgb2YgdGhlIG9iamVjdCBjb2RlIHdpdGggYSBjb3B5IG9mIHRoZQogICAgd3JpdHRlbiBvZmZlciB0byBwcm92aWRlIHRoZSBDb3JyZXNwb25kaW5nIFNvdXJjZS4gIFRoaXMKICAgIGFsdGVybmF0aXZlIGlzIGFsbG93ZWQgb25seSBvY2Nhc2lvbmFsbHkgYW5kIG5vbmNvbW1lcmNpYWxseSwgYW5kCiAgICBvbmx5IGlmIHlvdSByZWNlaXZlZCB0aGUgb2JqZWN0IGNvZGUgd2l0aCBzdWNoIGFuIG9mZmVyLCBpbiBhY2NvcmQKICAgIHdpdGggc3Vic2VjdGlvbiA2Yi4KCiAgICBkKSBDb252ZXkgdGhlIG9iamVjdCBjb2RlIGJ5IG9mZmVyaW5nIGFjY2VzcyBmcm9tIGEgZGVzaWduYXRlZAogICAgcGxhY2UgKGdyYXRpcyBvciBmb3IgYSBjaGFyZ2UpLCBhbmQgb2ZmZXIgZXF1aXZhbGVudCBhY2Nlc3MgdG8gdGhlCiAgICBDb3JyZXNwb25kaW5nIFNvdXJjZSBpbiB0aGUgc2FtZSB3YXkgdGhyb3VnaCB0aGUgc2FtZSBwbGFjZSBhdCBubwogICAgZnVydGhlciBjaGFyZ2UuICBZb3UgbmVlZCBub3QgcmVxdWlyZSByZWNpcGllbnRzIHRvIGNvcHkgdGhlCiAgICBDb3JyZXNwb25kaW5nIFNvdXJjZSBhbG9uZyB3aXRoIHRoZSBvYmplY3QgY29kZS4gIElmIHRoZSBwbGFjZSB0bwogICAgY29weSB0aGUgb2JqZWN0IGNvZGUgaXMgYSBuZXR3b3JrIHNlcnZlciwgdGhlIENvcnJlc3BvbmRpbmcgU291cmNlCiAgICBtYXkgYmUgb24gYSBkaWZmZXJlbnQgc2VydmVyIChvcGVyYXRlZCBieSB5b3Ugb3IgYSB0aGlyZCBwYXJ0eSkKICAgIHRoYXQgc3VwcG9ydHMgZXF1aXZhbGVudCBjb3B5aW5nIGZhY2lsaXRpZXMsIHByb3ZpZGVkIHlvdSBtYWludGFpbgogICAgY2xlYXIgZGlyZWN0aW9ucyBuZXh0IHRvIHRoZSBvYmplY3QgY29kZSBzYXlpbmcgd2hlcmUgdG8gZmluZCB0aGUKICAgIENvcnJlc3BvbmRpbmcgU291cmNlLiAgUmVnYXJkbGVzcyBvZiB3aGF0IHNlcnZlciBob3N0cyB0aGUKICAgIENvcnJlc3BvbmRpbmcgU291cmNlLCB5b3UgcmVtYWluIG9ibGlnYXRlZCB0byBlbnN1cmUgdGhhdCBpdCBpcwogICAgYXZhaWxhYmxlIGZvciBhcyBsb25nIGFzIG5lZWRlZCB0byBzYXRpc2Z5IHRoZXNlIHJlcXVpcmVtZW50cy4KCiAgICBlKSBDb252ZXkgdGhlIG9iamVjdCBjb2RlIHVzaW5nIHBlZXItdG8tcGVlciB0cmFuc21pc3Npb24sIHByb3ZpZGVkCiAgICB5b3UgaW5mb3JtIG90aGVyIHBlZXJzIHdoZXJlIHRoZSBvYmplY3QgY29kZSBhbmQgQ29ycmVzcG9uZGluZwogICAgU291cmNlIG9mIHRoZSB3b3JrIGFyZSBiZWluZyBvZmZlcmVkIHRvIHRoZSBnZW5lcmFsIHB1YmxpYyBhdCBubwogICAgY2hhcmdlIHVuZGVyIHN1YnNlY3Rpb24gNmQuCgogIEEgc2VwYXJhYmxlIHBvcnRpb24gb2YgdGhlIG9iamVjdCBjb2RlLCB3aG9zZSBzb3VyY2UgY29kZSBpcyBleGNsdWRlZApmcm9tIHRoZSBDb3JyZXNwb25kaW5nIFNvdXJjZSBhcyBhIFN5c3RlbSBMaWJyYXJ5LCBuZWVkIG5vdCBiZQppbmNsdWRlZCBpbiBjb252ZXlpbmcgdGhlIG9iamVjdCBjb2RlIHdvcmsuCgogIEEgIlVzZXIgUHJvZHVjdCIgaXMgZWl0aGVyICgxKSBhICJjb25zdW1lciBwcm9kdWN0Iiwgd2hpY2ggbWVhbnMgYW55CnRhbmdpYmxlIHBlcnNvbmFsIHByb3BlcnR5IHdoaWNoIGlzIG5vcm1hbGx5IHVzZWQgZm9yIHBlcnNvbmFsLCBmYW1pbHksCm9yIGhvdXNlaG9sZCBwdXJwb3Nlcywgb3IgKDIpIGFueXRoaW5nIGRlc2lnbmVkIG9yIHNvbGQgZm9yIGluY29ycG9yYXRpb24KaW50byBhIGR3ZWxsaW5nLiAgSW4gZGV0ZXJtaW5pbmcgd2hldGhlciBhIHByb2R1Y3QgaXMgYSBjb25zdW1lciBwcm9kdWN0LApkb3VidGZ1bCBjYXNlcyBzaGFsbCBiZSByZXNvbHZlZCBpbiBmYXZvciBvZiBjb3ZlcmFnZS4gIEZvciBhIHBhcnRpY3VsYXIKcHJvZHVjdCByZWNlaXZlZCBieSBhIHBhcnRpY3VsYXIgdXNlciwgIm5vcm1hbGx5IHVzZWQiIHJlZmVycyB0byBhCnR5cGljYWwgb3IgY29tbW9uIHVzZSBvZiB0aGF0IGNsYXNzIG9mIHByb2R1Y3QsIHJlZ2FyZGxlc3Mgb2YgdGhlIHN0YXR1cwpvZiB0aGUgcGFydGljdWxhciB1c2VyIG9yIG9mIHRoZSB3YXkgaW4gd2hpY2ggdGhlIHBhcnRpY3VsYXIgdXNlcgphY3R1YWxseSB1c2VzLCBvciBleHBlY3RzIG9yIGlzIGV4cGVjdGVkIHRvIHVzZSwgdGhlIHByb2R1Y3QuICBBIHByb2R1Y3QKaXMgYSBjb25zdW1lciBwcm9kdWN0IHJlZ2FyZGxlc3Mgb2Ygd2hldGhlciB0aGUgcHJvZHVjdCBoYXMgc3Vic3RhbnRpYWwKY29tbWVyY2lhbCwgaW5kdXN0cmlhbCBvciBub24tY29uc3VtZXIgdXNlcywgdW5sZXNzIHN1Y2ggdXNlcyByZXByZXNlbnQKdGhlIG9ubHkgc2lnbmlmaWNhbnQgbW9kZSBvZiB1c2Ugb2YgdGhlIHByb2R1Y3QuCgogICJJbnN0YWxsYXRpb24gSW5mb3JtYXRpb24iIGZvciBhIFVzZXIgUHJvZHVjdCBtZWFucyBhbnkgbWV0aG9kcywKcHJvY2VkdXJlcywgYXV0aG9yaXphdGlvbiBrZXlzLCBvciBvdGhlciBpbmZvcm1hdGlvbiByZXF1aXJlZCB0byBpbnN0YWxsCmFuZCBleGVjdXRlIG1vZGlmaWVkIHZlcnNpb25zIG9mIGEgY292ZXJlZCB3b3JrIGluIHRoYXQgVXNlciBQcm9kdWN0IGZyb20KYSBtb2RpZmllZCB2ZXJzaW9uIG9mIGl0cyBDb3JyZXNwb25kaW5nIFNvdXJjZS4gIFRoZSBpbmZvcm1hdGlvbiBtdXN0CnN1ZmZpY2UgdG8gZW5zdXJlIHRoYXQgdGhlIGNvbnRpbnVlZCBmdW5jdGlvbmluZyBvZiB0aGUgbW9kaWZpZWQgb2JqZWN0CmNvZGUgaXMgaW4gbm8gY2FzZSBwcmV2ZW50ZWQgb3IgaW50ZXJmZXJlZCB3aXRoIHNvbGVseSBiZWNhdXNlCm1vZGlmaWNhdGlvbiBoYXMgYmVlbiBtYWRlLgoKICBJZiB5b3UgY29udmV5IGFuIG9iamVjdCBjb2RlIHdvcmsgdW5kZXIgdGhpcyBzZWN0aW9uIGluLCBvciB3aXRoLCBvcgpzcGVjaWZpY2FsbHkgZm9yIHVzZSBpbiwgYSBVc2VyIFByb2R1Y3QsIGFuZCB0aGUgY29udmV5aW5nIG9jY3VycyBhcwpwYXJ0IG9mIGEgdHJhbnNhY3Rpb24gaW4gd2hpY2ggdGhlIHJpZ2h0IG9mIHBvc3Nlc3Npb24gYW5kIHVzZSBvZiB0aGUKVXNlciBQcm9kdWN0IGlzIHRyYW5zZmVycmVkIHRvIHRoZSByZWNpcGllbnQgaW4gcGVycGV0dWl0eSBvciBmb3IgYQpmaXhlZCB0ZXJtIChyZWdhcmRsZXNzIG9mIGhvdyB0aGUgdHJhbnNhY3Rpb24gaXMgY2hhcmFjdGVyaXplZCksIHRoZQpDb3JyZXNwb25kaW5nIFNvdXJjZSBjb252ZXllZCB1bmRlciB0aGlzIHNlY3Rpb24gbXVzdCBiZSBhY2NvbXBhbmllZApieSB0aGUgSW5zdGFsbGF0aW9uIEluZm9ybWF0aW9uLiAgQnV0IHRoaXMgcmVxdWlyZW1lbnQgZG9lcyBub3QgYXBwbHkKaWYgbmVpdGhlciB5b3Ugbm9yIGFueSB0aGlyZCBwYXJ0eSByZXRhaW5zIHRoZSBhYmlsaXR5IHRvIGluc3RhbGwKbW9kaWZpZWQgb2JqZWN0IGNvZGUgb24gdGhlIFVzZXIgUHJvZHVjdCAoZm9yIGV4YW1wbGUsIHRoZSB3b3JrIGhhcwpiZWVuIGluc3RhbGxlZCBpbiBST00pLgoKICBUaGUgcmVxdWlyZW1lbnQgdG8gcHJvdmlkZSBJbnN0YWxsYXRpb24gSW5mb3JtYXRpb24gZG9lcyBub3QgaW5jbHVkZSBhCnJlcXVpcmVtZW50IHRvIGNvbnRpbnVlIHRvIHByb3ZpZGUgc3VwcG9ydCBzZXJ2aWNlLCB3YXJyYW50eSwgb3IgdXBkYXRlcwpmb3IgYSB3b3JrIHRoYXQgaGFzIGJlZW4gbW9kaWZpZWQgb3IgaW5zdGFsbGVkIGJ5IHRoZSByZWNpcGllbnQsIG9yIGZvcgp0aGUgVXNlciBQcm9kdWN0IGluIHdoaWNoIGl0IGhhcyBiZWVuIG1vZGlmaWVkIG9yIGluc3RhbGxlZC4gIEFjY2VzcyB0byBhCm5ldHdvcmsgbWF5IGJlIGRlbmllZCB3aGVuIHRoZSBtb2RpZmljYXRpb24gaXRzZWxmIG1hdGVyaWFsbHkgYW5kCmFkdmVyc2VseSBhZmZlY3RzIHRoZSBvcGVyYXRpb24gb2YgdGhlIG5ldHdvcmsgb3IgdmlvbGF0ZXMgdGhlIHJ1bGVzIGFuZApwcm90b2NvbHMgZm9yIGNvbW11bmljYXRpb24gYWNyb3NzIHRoZSBuZXR3b3JrLgoKICBDb3JyZXNwb25kaW5nIFNvdXJjZSBjb252ZXllZCwgYW5kIEluc3RhbGxhdGlvbiBJbmZvcm1hdGlvbiBwcm92aWRlZCwKaW4gYWNjb3JkIHdpdGggdGhpcyBzZWN0aW9uIG11c3QgYmUgaW4gYSBmb3JtYXQgdGhhdCBpcyBwdWJsaWNseQpkb2N1bWVudGVkIChhbmQgd2l0aCBhbiBpbXBsZW1lbnRhdGlvbiBhdmFpbGFibGUgdG8gdGhlIHB1YmxpYyBpbgpzb3VyY2UgY29kZSBmb3JtKSwgYW5kIG11c3QgcmVxdWlyZSBubyBzcGVjaWFsIHBhc3N3b3JkIG9yIGtleSBmb3IKdW5wYWNraW5nLCByZWFkaW5nIG9yIGNvcHlpbmcuCgogIDcuIEFkZGl0aW9uYWwgVGVybXMuCgogICJBZGRpdGlvbmFsIHBlcm1pc3Npb25zIiBhcmUgdGVybXMgdGhhdCBzdXBwbGVtZW50IHRoZSB0ZXJtcyBvZiB0aGlzCkxpY2Vuc2UgYnkgbWFraW5nIGV4Y2VwdGlvbnMgZnJvbSBvbmUgb3IgbW9yZSBvZiBpdHMgY29uZGl0aW9ucy4KQWRkaXRpb25hbCBwZXJtaXNzaW9ucyB0aGF0IGFyZSBhcHBsaWNhYmxlIHRvIHRoZSBlbnRpcmUgUHJvZ3JhbSBzaGFsbApiZSB0cmVhdGVkIGFzIHRob3VnaCB0aGV5IHdlcmUgaW5jbHVkZWQgaW4gdGhpcyBMaWNlbnNlLCB0byB0aGUgZXh0ZW50CnRoYXQgdGhleSBhcmUgdmFsaWQgdW5kZXIgYXBwbGljYWJsZSBsYXcuICBJZiBhZGRpdGlvbmFsIHBlcm1pc3Npb25zCmFwcGx5IG9ubHkgdG8gcGFydCBvZiB0aGUgUHJvZ3JhbSwgdGhhdCBwYXJ0IG1heSBiZSB1c2VkIHNlcGFyYXRlbHkKdW5kZXIgdGhvc2UgcGVybWlzc2lvbnMsIGJ1dCB0aGUgZW50aXJlIFByb2dyYW0gcmVtYWlucyBnb3Zlcm5lZCBieQp0aGlzIExpY2Vuc2Ugd2l0aG91dCByZWdhcmQgdG8gdGhlIGFkZGl0aW9uYWwgcGVybWlzc2lvbnMuCgogIFdoZW4geW91IGNvbnZleSBhIGNvcHkgb2YgYSBjb3ZlcmVkIHdvcmssIHlvdSBtYXkgYXQgeW91ciBvcHRpb24KcmVtb3ZlIGFueSBhZGRpdGlvbmFsIHBlcm1pc3Npb25zIGZyb20gdGhhdCBjb3B5LCBvciBmcm9tIGFueSBwYXJ0IG9mCml0LiAgKEFkZGl0aW9uYWwgcGVybWlzc2lvbnMgbWF5IGJlIHdyaXR0ZW4gdG8gcmVxdWlyZSB0aGVpciBvd24KcmVtb3ZhbCBpbiBjZXJ0YWluIGNhc2VzIHdoZW4geW91IG1vZGlmeSB0aGUgd29yay4pICBZb3UgbWF5IHBsYWNlCmFkZGl0aW9uYWwgcGVybWlzc2lvbnMgb24gbWF0ZXJpYWwsIGFkZGVkIGJ5IHlvdSB0byBhIGNvdmVyZWQgd29yaywKZm9yIHdoaWNoIHlvdSBoYXZlIG9yIGNhbiBnaXZlIGFwcHJvcHJpYXRlIGNvcHlyaWdodCBwZXJtaXNzaW9uLgoKICBOb3R3aXRoc3RhbmRpbmcgYW55IG90aGVyIHByb3Zpc2lvbiBvZiB0aGlzIExpY2Vuc2UsIGZvciBtYXRlcmlhbCB5b3UKYWRkIHRvIGEgY292ZXJlZCB3b3JrLCB5b3UgbWF5IChpZiBhdXRob3JpemVkIGJ5IHRoZSBjb3B5cmlnaHQgaG9sZGVycyBvZgp0aGF0IG1hdGVyaWFsKSBzdXBwbGVtZW50IHRoZSB0ZXJtcyBvZiB0aGlzIExpY2Vuc2Ugd2l0aCB0ZXJtczoKCiAgICBhKSBEaXNjbGFpbWluZyB3YXJyYW50eSBvciBsaW1pdGluZyBsaWFiaWxpdHkgZGlmZmVyZW50bHkgZnJvbSB0aGUKICAgIHRlcm1zIG9mIHNlY3Rpb25zIDE1IGFuZCAxNiBvZiB0aGlzIExpY2Vuc2U7IG9yCgogICAgYikgUmVxdWlyaW5nIHByZXNlcnZhdGlvbiBvZiBzcGVjaWZpZWQgcmVhc29uYWJsZSBsZWdhbCBub3RpY2VzIG9yCiAgICBhdXRob3IgYXR0cmlidXRpb25zIGluIHRoYXQgbWF0ZXJpYWwgb3IgaW4gdGhlIEFwcHJvcHJpYXRlIExlZ2FsCiAgICBOb3RpY2VzIGRpc3BsYXllZCBieSB3b3JrcyBjb250YWluaW5nIGl0OyBvcgoKICAgIGMpIFByb2hpYml0aW5nIG1pc3JlcHJlc2VudGF0aW9uIG9mIHRoZSBvcmlnaW4gb2YgdGhhdCBtYXRlcmlhbCwgb3IKICAgIHJlcXVpcmluZyB0aGF0IG1vZGlmaWVkIHZlcnNpb25zIG9mIHN1Y2ggbWF0ZXJpYWwgYmUgbWFya2VkIGluCiAgICByZWFzb25hYmxlIHdheXMgYXMgZGlmZmVyZW50IGZyb20gdGhlIG9yaWdpbmFsIHZlcnNpb247IG9yCgogICAgZCkgTGltaXRpbmcgdGhlIHVzZSBmb3IgcHVibGljaXR5IHB1cnBvc2VzIG9mIG5hbWVzIG9mIGxpY2Vuc29ycyBvcgogICAgYXV0aG9ycyBvZiB0aGUgbWF0ZXJpYWw7IG9yCgogICAgZSkgRGVjbGluaW5nIHRvIGdyYW50IHJpZ2h0cyB1bmRlciB0cmFkZW1hcmsgbGF3IGZvciB1c2Ugb2Ygc29tZQogICAgdHJhZGUgbmFtZXMsIHRyYWRlbWFya3MsIG9yIHNlcnZpY2UgbWFya3M7IG9yCgogICAgZikgUmVxdWlyaW5nIGluZGVtbmlmaWNhdGlvbiBvZiBsaWNlbnNvcnMgYW5kIGF1dGhvcnMgb2YgdGhhdAogICAgbWF0ZXJpYWwgYnkgYW55b25lIHdobyBjb252ZXlzIHRoZSBtYXRlcmlhbCAob3IgbW9kaWZpZWQgdmVyc2lvbnMgb2YKICAgIGl0KSB3aXRoIGNvbnRyYWN0dWFsIGFzc3VtcHRpb25zIG9mIGxpYWJpbGl0eSB0byB0aGUgcmVjaXBpZW50LCBmb3IKICAgIGFueSBsaWFiaWxpdHkgdGhhdCB0aGVzZSBjb250cmFjdHVhbCBhc3N1bXB0aW9ucyBkaXJlY3RseSBpbXBvc2Ugb24KICAgIHRob3NlIGxpY2Vuc29ycyBhbmQgYXV0aG9ycy4KCiAgQWxsIG90aGVyIG5vbi1wZXJtaXNzaXZlIGFkZGl0aW9uYWwgdGVybXMgYXJlIGNvbnNpZGVyZWQgImZ1cnRoZXIKcmVzdHJpY3Rpb25zIiB3aXRoaW4gdGhlIG1lYW5pbmcgb2Ygc2VjdGlvbiAxMC4gIElmIHRoZSBQcm9ncmFtIGFzIHlvdQpyZWNlaXZlZCBpdCwgb3IgYW55IHBhcnQgb2YgaXQsIGNvbnRhaW5zIGEgbm90aWNlIHN0YXRpbmcgdGhhdCBpdCBpcwpnb3Zlcm5lZCBieSB0aGlzIExpY2Vuc2UgYWxvbmcgd2l0aCBhIHRlcm0gdGhhdCBpcyBhIGZ1cnRoZXIKcmVzdHJpY3Rpb24sIHlvdSBtYXkgcmVtb3ZlIHRoYXQgdGVybS4gIElmIGEgbGljZW5zZSBkb2N1bWVudCBjb250YWlucwphIGZ1cnRoZXIgcmVzdHJpY3Rpb24gYnV0IHBlcm1pdHMgcmVsaWNlbnNpbmcgb3IgY29udmV5aW5nIHVuZGVyIHRoaXMKTGljZW5zZSwgeW91IG1heSBhZGQgdG8gYSBjb3ZlcmVkIHdvcmsgbWF0ZXJpYWwgZ292ZXJuZWQgYnkgdGhlIHRlcm1zCm9mIHRoYXQgbGljZW5zZSBkb2N1bWVudCwgcHJvdmlkZWQgdGhhdCB0aGUgZnVydGhlciByZXN0cmljdGlvbiBkb2VzCm5vdCBzdXJ2aXZlIHN1Y2ggcmVsaWNlbnNpbmcgb3IgY29udmV5aW5nLgoKICBJZiB5b3UgYWRkIHRlcm1zIHRvIGEgY292ZXJlZCB3b3JrIGluIGFjY29yZCB3aXRoIHRoaXMgc2VjdGlvbiwgeW91Cm11c3QgcGxhY2UsIGluIHRoZSByZWxldmFudCBzb3VyY2UgZmlsZXMsIGEgc3RhdGVtZW50IG9mIHRoZQphZGRpdGlvbmFsIHRlcm1zIHRoYXQgYXBwbHkgdG8gdGhvc2UgZmlsZXMsIG9yIGEgbm90aWNlIGluZGljYXRpbmcKd2hlcmUgdG8gZmluZCB0aGUgYXBwbGljYWJsZSB0ZXJtcy4KCiAgQWRkaXRpb25hbCB0ZXJtcywgcGVybWlzc2l2ZSBvciBub24tcGVybWlzc2l2ZSwgbWF5IGJlIHN0YXRlZCBpbiB0aGUKZm9ybSBvZiBhIHNlcGFyYXRlbHkgd3JpdHRlbiBsaWNlbnNlLCBvciBzdGF0ZWQgYXMgZXhjZXB0aW9uczsKdGhlIGFib3ZlIHJlcXVpcmVtZW50cyBhcHBseSBlaXRoZXIgd2F5LgoKICA4LiBUZXJtaW5hdGlvbi4KCiAgWW91IG1heSBub3QgcHJvcGFnYXRlIG9yIG1vZGlmeSBhIGNvdmVyZWQgd29yayBleGNlcHQgYXMgZXhwcmVzc2x5CnByb3ZpZGVkIHVuZGVyIHRoaXMgTGljZW5zZS4gIEFueSBhdHRlbXB0IG90aGVyd2lzZSB0byBwcm9wYWdhdGUgb3IKbW9kaWZ5IGl0IGlzIHZvaWQsIGFuZCB3aWxsIGF1dG9tYXRpY2FsbHkgdGVybWluYXRlIHlvdXIgcmlnaHRzIHVuZGVyCnRoaXMgTGljZW5zZSAoaW5jbHVkaW5nIGFueSBwYXRlbnQgbGljZW5zZXMgZ3JhbnRlZCB1bmRlciB0aGUgdGhpcmQKcGFyYWdyYXBoIG9mIHNlY3Rpb24gMTEpLgoKICBIb3dldmVyLCBpZiB5b3UgY2Vhc2UgYWxsIHZpb2xhdGlvbiBvZiB0aGlzIExpY2Vuc2UsIHRoZW4geW91cgpsaWNlbnNlIGZyb20gYSBwYXJ0aWN1bGFyIGNvcHlyaWdodCBob2xkZXIgaXMgcmVpbnN0YXRlZCAoYSkKcHJvdmlzaW9uYWxseSwgdW5sZXNzIGFuZCB1bnRpbCB0aGUgY29weXJpZ2h0IGhvbGRlciBleHBsaWNpdGx5IGFuZApmaW5hbGx5IHRlcm1pbmF0ZXMgeW91ciBsaWNlbnNlLCBhbmQgKGIpIHBlcm1hbmVudGx5LCBpZiB0aGUgY29weXJpZ2h0CmhvbGRlciBmYWlscyB0byBub3RpZnkgeW91IG9mIHRoZSB2aW9sYXRpb24gYnkgc29tZSByZWFzb25hYmxlIG1lYW5zCnByaW9yIHRvIDYwIGRheXMgYWZ0ZXIgdGhlIGNlc3NhdGlvbi4KCiAgTW9yZW92ZXIsIHlvdXIgbGljZW5zZSBmcm9tIGEgcGFydGljdWxhciBjb3B5cmlnaHQgaG9sZGVyIGlzCnJlaW5zdGF0ZWQgcGVybWFuZW50bHkgaWYgdGhlIGNvcHlyaWdodCBob2xkZXIgbm90aWZpZXMgeW91IG9mIHRoZQp2aW9sYXRpb24gYnkgc29tZSByZWFzb25hYmxlIG1lYW5zLCB0aGlzIGlzIHRoZSBmaXJzdCB0aW1lIHlvdSBoYXZlCnJlY2VpdmVkIG5vdGljZSBvZiB2aW9sYXRpb24gb2YgdGhpcyBMaWNlbnNlIChmb3IgYW55IHdvcmspIGZyb20gdGhhdApjb3B5cmlnaHQgaG9sZGVyLCBhbmQgeW91IGN1cmUgdGhlIHZpb2xhdGlvbiBwcmlvciB0byAzMCBkYXlzIGFmdGVyCnlvdXIgcmVjZWlwdCBvZiB0aGUgbm90aWNlLgoKICBUZXJtaW5hdGlvbiBvZiB5b3VyIHJpZ2h0cyB1bmRlciB0aGlzIHNlY3Rpb24gZG9lcyBub3QgdGVybWluYXRlIHRoZQpsaWNlbnNlcyBvZiBwYXJ0aWVzIHdobyBoYXZlIHJlY2VpdmVkIGNvcGllcyBvciByaWdodHMgZnJvbSB5b3UgdW5kZXIKdGhpcyBMaWNlbnNlLiAgSWYgeW91ciByaWdodHMgaGF2ZSBiZWVuIHRlcm1pbmF0ZWQgYW5kIG5vdCBwZXJtYW5lbnRseQpyZWluc3RhdGVkLCB5b3UgZG8gbm90IHF1YWxpZnkgdG8gcmVjZWl2ZSBuZXcgbGljZW5zZXMgZm9yIHRoZSBzYW1lCm1hdGVyaWFsIHVuZGVyIHNlY3Rpb24gMTAuCgogIDkuIEFjY2VwdGFuY2UgTm90IFJlcXVpcmVkIGZvciBIYXZpbmcgQ29waWVzLgoKICBZb3UgYXJlIG5vdCByZXF1aXJlZCB0byBhY2NlcHQgdGhpcyBMaWNlbnNlIGluIG9yZGVyIHRvIHJlY2VpdmUgb3IKcnVuIGEgY29weSBvZiB0aGUgUHJvZ3JhbS4gIEFuY2lsbGFyeSBwcm9wYWdhdGlvbiBvZiBhIGNvdmVyZWQgd29yawpvY2N1cnJpbmcgc29sZWx5IGFzIGEgY29uc2VxdWVuY2Ugb2YgdXNpbmcgcGVlci10by1wZWVyIHRyYW5zbWlzc2lvbgp0byByZWNlaXZlIGEgY29weSBsaWtld2lzZSBkb2VzIG5vdCByZXF1aXJlIGFjY2VwdGFuY2UuICBIb3dldmVyLApub3RoaW5nIG90aGVyIHRoYW4gdGhpcyBMaWNlbnNlIGdyYW50cyB5b3UgcGVybWlzc2lvbiB0byBwcm9wYWdhdGUgb3IKbW9kaWZ5IGFueSBjb3ZlcmVkIHdvcmsuICBUaGVzZSBhY3Rpb25zIGluZnJpbmdlIGNvcHlyaWdodCBpZiB5b3UgZG8Kbm90IGFjY2VwdCB0aGlzIExpY2Vuc2UuICBUaGVyZWZvcmUsIGJ5IG1vZGlmeWluZyBvciBwcm9wYWdhdGluZyBhCmNvdmVyZWQgd29yaywgeW91IGluZGljYXRlIHlvdXIgYWNjZXB0YW5jZSBvZiB0aGlzIExpY2Vuc2UgdG8gZG8gc28uCgogIDEwLiBBdXRvbWF0aWMgTGljZW5zaW5nIG9mIERvd25zdHJlYW0gUmVjaXBpZW50cy4KCiAgRWFjaCB0aW1lIHlvdSBjb252ZXkgYSBjb3ZlcmVkIHdvcmssIHRoZSByZWNpcGllbnQgYXV0b21hdGljYWxseQpyZWNlaXZlcyBhIGxpY2Vuc2UgZnJvbSB0aGUgb3JpZ2luYWwgbGljZW5zb3JzLCB0byBydW4sIG1vZGlmeSBhbmQKcHJvcGFnYXRlIHRoYXQgd29yaywgc3ViamVjdCB0byB0aGlzIExpY2Vuc2UuICBZb3UgYXJlIG5vdCByZXNwb25zaWJsZQpmb3IgZW5mb3JjaW5nIGNvbXBsaWFuY2UgYnkgdGhpcmQgcGFydGllcyB3aXRoIHRoaXMgTGljZW5zZS4KCiAgQW4gImVudGl0eSB0cmFuc2FjdGlvbiIgaXMgYSB0cmFuc2FjdGlvbiB0cmFuc2ZlcnJpbmcgY29udHJvbCBvZiBhbgpvcmdhbml6YXRpb24sIG9yIHN1YnN0YW50aWFsbHkgYWxsIGFzc2V0cyBvZiBvbmUsIG9yIHN1YmRpdmlkaW5nIGFuCm9yZ2FuaXphdGlvbiwgb3IgbWVyZ2luZyBvcmdhbml6YXRpb25zLiAgSWYgcHJvcGFnYXRpb24gb2YgYSBjb3ZlcmVkCndvcmsgcmVzdWx0cyBmcm9tIGFuIGVudGl0eSB0cmFuc2FjdGlvbiwgZWFjaCBwYXJ0eSB0byB0aGF0CnRyYW5zYWN0aW9uIHdobyByZWNlaXZlcyBhIGNvcHkgb2YgdGhlIHdvcmsgYWxzbyByZWNlaXZlcyB3aGF0ZXZlcgpsaWNlbnNlcyB0byB0aGUgd29yayB0aGUgcGFydHkncyBwcmVkZWNlc3NvciBpbiBpbnRlcmVzdCBoYWQgb3IgY291bGQKZ2l2ZSB1bmRlciB0aGUgcHJldmlvdXMgcGFyYWdyYXBoLCBwbHVzIGEgcmlnaHQgdG8gcG9zc2Vzc2lvbiBvZiB0aGUKQ29ycmVzcG9uZGluZyBTb3VyY2Ugb2YgdGhlIHdvcmsgZnJvbSB0aGUgcHJlZGVjZXNzb3IgaW4gaW50ZXJlc3QsIGlmCnRoZSBwcmVkZWNlc3NvciBoYXMgaXQgb3IgY2FuIGdldCBpdCB3aXRoIHJlYXNvbmFibGUgZWZmb3J0cy4KCiAgWW91IG1heSBub3QgaW1wb3NlIGFueSBmdXJ0aGVyIHJlc3RyaWN0aW9ucyBvbiB0aGUgZXhlcmNpc2Ugb2YgdGhlCnJpZ2h0cyBncmFudGVkIG9yIGFmZmlybWVkIHVuZGVyIHRoaXMgTGljZW5zZS4gIEZvciBleGFtcGxlLCB5b3UgbWF5Cm5vdCBpbXBvc2UgYSBsaWNlbnNlIGZlZSwgcm95YWx0eSwgb3Igb3RoZXIgY2hhcmdlIGZvciBleGVyY2lzZSBvZgpyaWdodHMgZ3JhbnRlZCB1bmRlciB0aGlzIExpY2Vuc2UsIGFuZCB5b3UgbWF5IG5vdCBpbml0aWF0ZSBsaXRpZ2F0aW9uCihpbmNsdWRpbmcgYSBjcm9zcy1jbGFpbSBvciBjb3VudGVyY2xhaW0gaW4gYSBsYXdzdWl0KSBhbGxlZ2luZyB0aGF0CmFueSBwYXRlbnQgY2xhaW0gaXMgaW5mcmluZ2VkIGJ5IG1ha2luZywgdXNpbmcsIHNlbGxpbmcsIG9mZmVyaW5nIGZvcgpzYWxlLCBvciBpbXBvcnRpbmcgdGhlIFByb2dyYW0gb3IgYW55IHBvcnRpb24gb2YgaXQuCgogIDExLiBQYXRlbnRzLgoKICBBICJjb250cmlidXRvciIgaXMgYSBjb3B5cmlnaHQgaG9sZGVyIHdobyBhdXRob3JpemVzIHVzZSB1bmRlciB0aGlzCkxpY2Vuc2Ugb2YgdGhlIFByb2dyYW0gb3IgYSB3b3JrIG9uIHdoaWNoIHRoZSBQcm9ncmFtIGlzIGJhc2VkLiAgVGhlCndvcmsgdGh1cyBsaWNlbnNlZCBpcyBjYWxsZWQgdGhlIGNvbnRyaWJ1dG9yJ3MgImNvbnRyaWJ1dG9yIHZlcnNpb24iLgoKICBBIGNvbnRyaWJ1dG9yJ3MgImVzc2VudGlhbCBwYXRlbnQgY2xhaW1zIiBhcmUgYWxsIHBhdGVudCBjbGFpbXMKb3duZWQgb3IgY29udHJvbGxlZCBieSB0aGUgY29udHJpYnV0b3IsIHdoZXRoZXIgYWxyZWFkeSBhY3F1aXJlZCBvcgpoZXJlYWZ0ZXIgYWNxdWlyZWQsIHRoYXQgd291bGQgYmUgaW5mcmluZ2VkIGJ5IHNvbWUgbWFubmVyLCBwZXJtaXR0ZWQKYnkgdGhpcyBMaWNlbnNlLCBvZiBtYWtpbmcsIHVzaW5nLCBvciBzZWxsaW5nIGl0cyBjb250cmlidXRvciB2ZXJzaW9uLApidXQgZG8gbm90IGluY2x1ZGUgY2xhaW1zIHRoYXQgd291bGQgYmUgaW5mcmluZ2VkIG9ubHkgYXMgYQpjb25zZXF1ZW5jZSBvZiBmdXJ0aGVyIG1vZGlmaWNhdGlvbiBvZiB0aGUgY29udHJpYnV0b3IgdmVyc2lvbi4gIEZvcgpwdXJwb3NlcyBvZiB0aGlzIGRlZmluaXRpb24sICJjb250cm9sIiBpbmNsdWRlcyB0aGUgcmlnaHQgdG8gZ3JhbnQKcGF0ZW50IHN1YmxpY2Vuc2VzIGluIGEgbWFubmVyIGNvbnNpc3RlbnQgd2l0aCB0aGUgcmVxdWlyZW1lbnRzIG9mCnRoaXMgTGljZW5zZS4KCiAgRWFjaCBjb250cmlidXRvciBncmFudHMgeW91IGEgbm9uLWV4Y2x1c2l2ZSwgd29ybGR3aWRlLCByb3lhbHR5LWZyZWUKcGF0ZW50IGxpY2Vuc2UgdW5kZXIgdGhlIGNvbnRyaWJ1dG9yJ3MgZXNzZW50aWFsIHBhdGVudCBjbGFpbXMsIHRvCm1ha2UsIHVzZSwgc2VsbCwgb2ZmZXIgZm9yIHNhbGUsIGltcG9ydCBhbmQgb3RoZXJ3aXNlIHJ1biwgbW9kaWZ5IGFuZApwcm9wYWdhdGUgdGhlIGNvbnRlbnRzIG9mIGl0cyBjb250cmlidXRvciB2ZXJzaW9uLgoKICBJbiB0aGUgZm9sbG93aW5nIHRocmVlIHBhcmFncmFwaHMsIGEgInBhdGVudCBsaWNlbnNlIiBpcyBhbnkgZXhwcmVzcwphZ3JlZW1lbnQgb3IgY29tbWl0bWVudCwgaG93ZXZlciBkZW5vbWluYXRlZCwgbm90IHRvIGVuZm9yY2UgYSBwYXRlbnQKKHN1Y2ggYXMgYW4gZXhwcmVzcyBwZXJtaXNzaW9uIHRvIHByYWN0aWNlIGEgcGF0ZW50IG9yIGNvdmVuYW50IG5vdCB0bwpzdWUgZm9yIHBhdGVudCBpbmZyaW5nZW1lbnQpLiAgVG8gImdyYW50IiBzdWNoIGEgcGF0ZW50IGxpY2Vuc2UgdG8gYQpwYXJ0eSBtZWFucyB0byBtYWtlIHN1Y2ggYW4gYWdyZWVtZW50IG9yIGNvbW1pdG1lbnQgbm90IHRvIGVuZm9yY2UgYQpwYXRlbnQgYWdhaW5zdCB0aGUgcGFydHkuCgogIElmIHlvdSBjb252ZXkgYSBjb3ZlcmVkIHdvcmssIGtub3dpbmdseSByZWx5aW5nIG9uIGEgcGF0ZW50IGxpY2Vuc2UsCmFuZCB0aGUgQ29ycmVzcG9uZGluZyBTb3VyY2Ugb2YgdGhlIHdvcmsgaXMgbm90IGF2YWlsYWJsZSBmb3IgYW55b25lCnRvIGNvcHksIGZyZWUgb2YgY2hhcmdlIGFuZCB1bmRlciB0aGUgdGVybXMgb2YgdGhpcyBMaWNlbnNlLCB0aHJvdWdoIGEKcHVibGljbHkgYXZhaWxhYmxlIG5ldHdvcmsgc2VydmVyIG9yIG90aGVyIHJlYWRpbHkgYWNjZXNzaWJsZSBtZWFucywKdGhlbiB5b3UgbXVzdCBlaXRoZXIgKDEpIGNhdXNlIHRoZSBDb3JyZXNwb25kaW5nIFNvdXJjZSB0byBiZSBzbwphdmFpbGFibGUsIG9yICgyKSBhcnJhbmdlIHRvIGRlcHJpdmUgeW91cnNlbGYgb2YgdGhlIGJlbmVmaXQgb2YgdGhlCnBhdGVudCBsaWNlbnNlIGZvciB0aGlzIHBhcnRpY3VsYXIgd29yaywgb3IgKDMpIGFycmFuZ2UsIGluIGEgbWFubmVyCmNvbnNpc3RlbnQgd2l0aCB0aGUgcmVxdWlyZW1lbnRzIG9mIHRoaXMgTGljZW5zZSwgdG8gZXh0ZW5kIHRoZSBwYXRlbnQKbGljZW5zZSB0byBkb3duc3RyZWFtIHJlY2lwaWVudHMuICAiS25vd2luZ2x5IHJlbHlpbmciIG1lYW5zIHlvdSBoYXZlCmFjdHVhbCBrbm93bGVkZ2UgdGhhdCwgYnV0IGZvciB0aGUgcGF0ZW50IGxpY2Vuc2UsIHlvdXIgY29udmV5aW5nIHRoZQpjb3ZlcmVkIHdvcmsgaW4gYSBjb3VudHJ5LCBvciB5b3VyIHJlY2lwaWVudCdzIHVzZSBvZiB0aGUgY292ZXJlZCB3b3JrCmluIGEgY291bnRyeSwgd291bGQgaW5mcmluZ2Ugb25lIG9yIG1vcmUgaWRlbnRpZmlhYmxlIHBhdGVudHMgaW4gdGhhdApjb3VudHJ5IHRoYXQgeW91IGhhdmUgcmVhc29uIHRvIGJlbGlldmUgYXJlIHZhbGlkLgoKICBJZiwgcHVyc3VhbnQgdG8gb3IgaW4gY29ubmVjdGlvbiB3aXRoIGEgc2luZ2xlIHRyYW5zYWN0aW9uIG9yCmFycmFuZ2VtZW50LCB5b3UgY29udmV5LCBvciBwcm9wYWdhdGUgYnkgcHJvY3VyaW5nIGNvbnZleWFuY2Ugb2YsIGEKY292ZXJlZCB3b3JrLCBhbmQgZ3JhbnQgYSBwYXRlbnQgbGljZW5zZSB0byBzb21lIG9mIHRoZSBwYXJ0aWVzCnJlY2VpdmluZyB0aGUgY292ZXJlZCB3b3JrIGF1dGhvcml6aW5nIHRoZW0gdG8gdXNlLCBwcm9wYWdhdGUsIG1vZGlmeQpvciBjb252ZXkgYSBzcGVjaWZpYyBjb3B5IG9mIHRoZSBjb3ZlcmVkIHdvcmssIHRoZW4gdGhlIHBhdGVudCBsaWNlbnNlCnlvdSBncmFudCBpcyBhdXRvbWF0aWNhbGx5IGV4dGVuZGVkIHRvIGFsbCByZWNpcGllbnRzIG9mIHRoZSBjb3ZlcmVkCndvcmsgYW5kIHdvcmtzIGJhc2VkIG9uIGl0LgoKICBBIHBhdGVudCBsaWNlbnNlIGlzICJkaXNjcmltaW5hdG9yeSIgaWYgaXQgZG9lcyBub3QgaW5jbHVkZSB3aXRoaW4KdGhlIHNjb3BlIG9mIGl0cyBjb3ZlcmFnZSwgcHJvaGliaXRzIHRoZSBleGVyY2lzZSBvZiwgb3IgaXMKY29uZGl0aW9uZWQgb24gdGhlIG5vbi1leGVyY2lzZSBvZiBvbmUgb3IgbW9yZSBvZiB0aGUgcmlnaHRzIHRoYXQgYXJlCnNwZWNpZmljYWxseSBncmFudGVkIHVuZGVyIHRoaXMgTGljZW5zZS4gIFlvdSBtYXkgbm90IGNvbnZleSBhIGNvdmVyZWQKd29yayBpZiB5b3UgYXJlIGEgcGFydHkgdG8gYW4gYXJyYW5nZW1lbnQgd2l0aCBhIHRoaXJkIHBhcnR5IHRoYXQgaXMKaW4gdGhlIGJ1c2luZXNzIG9mIGRpc3RyaWJ1dGluZyBzb2Z0d2FyZSwgdW5kZXIgd2hpY2ggeW91IG1ha2UgcGF5bWVudAp0byB0aGUgdGhpcmQgcGFydHkgYmFzZWQgb24gdGhlIGV4dGVudCBvZiB5b3VyIGFjdGl2aXR5IG9mIGNvbnZleWluZwp0aGUgd29yaywgYW5kIHVuZGVyIHdoaWNoIHRoZSB0aGlyZCBwYXJ0eSBncmFudHMsIHRvIGFueSBvZiB0aGUKcGFydGllcyB3aG8gd291bGQgcmVjZWl2ZSB0aGUgY292ZXJlZCB3b3JrIGZyb20geW91LCBhIGRpc2NyaW1pbmF0b3J5CnBhdGVudCBsaWNlbnNlIChhKSBpbiBjb25uZWN0aW9uIHdpdGggY29waWVzIG9mIHRoZSBjb3ZlcmVkIHdvcmsKY29udmV5ZWQgYnkgeW91IChvciBjb3BpZXMgbWFkZSBmcm9tIHRob3NlIGNvcGllcyksIG9yIChiKSBwcmltYXJpbHkKZm9yIGFuZCBpbiBjb25uZWN0aW9uIHdpdGggc3BlY2lmaWMgcHJvZHVjdHMgb3IgY29tcGlsYXRpb25zIHRoYXQKY29udGFpbiB0aGUgY292ZXJlZCB3b3JrLCB1bmxlc3MgeW91IGVudGVyZWQgaW50byB0aGF0IGFycmFuZ2VtZW50LApvciB0aGF0IHBhdGVudCBsaWNlbnNlIHdhcyBncmFudGVkLCBwcmlvciB0byAyOCBNYXJjaCAyMDA3LgoKICBOb3RoaW5nIGluIHRoaXMgTGljZW5zZSBzaGFsbCBiZSBjb25zdHJ1ZWQgYXMgZXhjbHVkaW5nIG9yIGxpbWl0aW5nCmFueSBpbXBsaWVkIGxpY2Vuc2Ugb3Igb3RoZXIgZGVmZW5zZXMgdG8gaW5mcmluZ2VtZW50IHRoYXQgbWF5Cm90aGVyd2lzZSBiZSBhdmFpbGFibGUgdG8geW91IHVuZGVyIGFwcGxpY2FibGUgcGF0ZW50IGxhdy4KCiAgMTIuIE5vIFN1cnJlbmRlciBvZiBPdGhlcnMnIEZyZWVkb20uCgogIElmIGNvbmRpdGlvbnMgYXJlIGltcG9zZWQgb24geW91ICh3aGV0aGVyIGJ5IGNvdXJ0IG9yZGVyLCBhZ3JlZW1lbnQgb3IKb3RoZXJ3aXNlKSB0aGF0IGNvbnRyYWRpY3QgdGhlIGNvbmRpdGlvbnMgb2YgdGhpcyBMaWNlbnNlLCB0aGV5IGRvIG5vdApleGN1c2UgeW91IGZyb20gdGhlIGNvbmRpdGlvbnMgb2YgdGhpcyBMaWNlbnNlLiAgSWYgeW91IGNhbm5vdCBjb252ZXkgYQpjb3ZlcmVkIHdvcmsgc28gYXMgdG8gc2F0aXNmeSBzaW11bHRhbmVvdXNseSB5b3VyIG9ibGlnYXRpb25zIHVuZGVyIHRoaXMKTGljZW5zZSBhbmQgYW55IG90aGVyIHBlcnRpbmVudCBvYmxpZ2F0aW9ucywgdGhlbiBhcyBhIGNvbnNlcXVlbmNlIHlvdSBtYXkKbm90IGNvbnZleSBpdCBhdCBhbGwuICBGb3IgZXhhbXBsZSwgaWYgeW91IGFncmVlIHRvIHRlcm1zIHRoYXQgb2JsaWdhdGUgeW91CnRvIGNvbGxlY3QgYSByb3lhbHR5IGZvciBmdXJ0aGVyIGNvbnZleWluZyBmcm9tIHRob3NlIHRvIHdob20geW91IGNvbnZleQp0aGUgUHJvZ3JhbSwgdGhlIG9ubHkgd2F5IHlvdSBjb3VsZCBzYXRpc2Z5IGJvdGggdGhvc2UgdGVybXMgYW5kIHRoaXMKTGljZW5zZSB3b3VsZCBiZSB0byByZWZyYWluIGVudGlyZWx5IGZyb20gY29udmV5aW5nIHRoZSBQcm9ncmFtLgoKICAxMy4gVXNlIHdpdGggdGhlIEdOVSBBZmZlcm8gR2VuZXJhbCBQdWJsaWMgTGljZW5zZS4KCiAgTm90d2l0aHN0YW5kaW5nIGFueSBvdGhlciBwcm92aXNpb24gb2YgdGhpcyBMaWNlbnNlLCB5b3UgaGF2ZQpwZXJtaXNzaW9uIHRvIGxpbmsgb3IgY29tYmluZSBhbnkgY292ZXJlZCB3b3JrIHdpdGggYSB3b3JrIGxpY2Vuc2VkCnVuZGVyIHZlcnNpb24gMyBvZiB0aGUgR05VIEFmZmVybyBHZW5lcmFsIFB1YmxpYyBMaWNlbnNlIGludG8gYSBzaW5nbGUKY29tYmluZWQgd29yaywgYW5kIHRvIGNvbnZleSB0aGUgcmVzdWx0aW5nIHdvcmsuICBUaGUgdGVybXMgb2YgdGhpcwpMaWNlbnNlIHdpbGwgY29udGludWUgdG8gYXBwbHkgdG8gdGhlIHBhcnQgd2hpY2ggaXMgdGhlIGNvdmVyZWQgd29yaywKYnV0IHRoZSBzcGVjaWFsIHJlcXVpcmVtZW50cyBvZiB0aGUgR05VIEFmZmVybyBHZW5lcmFsIFB1YmxpYyBMaWNlbnNlLApzZWN0aW9uIDEzLCBjb25jZXJuaW5nIGludGVyYWN0aW9uIHRocm91Z2ggYSBuZXR3b3JrIHdpbGwgYXBwbHkgdG8gdGhlCmNvbWJpbmF0aW9uIGFzIHN1Y2guCgogIDE0LiBSZXZpc2VkIFZlcnNpb25zIG9mIHRoaXMgTGljZW5zZS4KCiAgVGhlIEZyZWUgU29mdHdhcmUgRm91bmRhdGlvbiBtYXkgcHVibGlzaCByZXZpc2VkIGFuZC9vciBuZXcgdmVyc2lvbnMgb2YKdGhlIEdOVSBHZW5lcmFsIFB1YmxpYyBMaWNlbnNlIGZyb20gdGltZSB0byB0aW1lLiAgU3VjaCBuZXcgdmVyc2lvbnMgd2lsbApiZSBzaW1pbGFyIGluIHNwaXJpdCB0byB0aGUgcHJlc2VudCB2ZXJzaW9uLCBidXQgbWF5IGRpZmZlciBpbiBkZXRhaWwgdG8KYWRkcmVzcyBuZXcgcHJvYmxlbXMgb3IgY29uY2VybnMuCgogIEVhY2ggdmVyc2lvbiBpcyBnaXZlbiBhIGRpc3Rpbmd1aXNoaW5nIHZlcnNpb24gbnVtYmVyLiAgSWYgdGhlClByb2dyYW0gc3BlY2lmaWVzIHRoYXQgYSBjZXJ0YWluIG51bWJlcmVkIHZlcnNpb24gb2YgdGhlIEdOVSBHZW5lcmFsClB1YmxpYyBMaWNlbnNlICJvciBhbnkgbGF0ZXIgdmVyc2lvbiIgYXBwbGllcyB0byBpdCwgeW91IGhhdmUgdGhlCm9wdGlvbiBvZiBmb2xsb3dpbmcgdGhlIHRlcm1zIGFuZCBjb25kaXRpb25zIGVpdGhlciBvZiB0aGF0IG51bWJlcmVkCnZlcnNpb24gb3Igb2YgYW55IGxhdGVyIHZlcnNpb24gcHVibGlzaGVkIGJ5IHRoZSBGcmVlIFNvZnR3YXJlCkZvdW5kYXRpb24uICBJZiB0aGUgUHJvZ3JhbSBkb2VzIG5vdCBzcGVjaWZ5IGEgdmVyc2lvbiBudW1iZXIgb2YgdGhlCkdOVSBHZW5lcmFsIFB1YmxpYyBMaWNlbnNlLCB5b3UgbWF5IGNob29zZSBhbnkgdmVyc2lvbiBldmVyIHB1Ymxpc2hlZApieSB0aGUgRnJlZSBTb2Z0d2FyZSBGb3VuZGF0aW9uLgoKICBJZiB0aGUgUHJvZ3JhbSBzcGVjaWZpZXMgdGhhdCBhIHByb3h5IGNhbiBkZWNpZGUgd2hpY2ggZnV0dXJlCnZlcnNpb25zIG9mIHRoZSBHTlUgR2VuZXJhbCBQdWJsaWMgTGljZW5zZSBjYW4gYmUgdXNlZCwgdGhhdCBwcm94eSdzCnB1YmxpYyBzdGF0ZW1lbnQgb2YgYWNjZXB0YW5jZSBvZiBhIHZlcnNpb24gcGVybWFuZW50bHkgYXV0aG9yaXplcyB5b3UKdG8gY2hvb3NlIHRoYXQgdmVyc2lvbiBmb3IgdGhlIFByb2dyYW0uCgogIExhdGVyIGxpY2Vuc2UgdmVyc2lvbnMgbWF5IGdpdmUgeW91IGFkZGl0aW9uYWwgb3IgZGlmZmVyZW50CnBlcm1pc3Npb25zLiAgSG93ZXZlciwgbm8gYWRkaXRpb25hbCBvYmxpZ2F0aW9ucyBhcmUgaW1wb3NlZCBvbiBhbnkKYXV0aG9yIG9yIGNvcHlyaWdodCBob2xkZXIgYXMgYSByZXN1bHQgb2YgeW91ciBjaG9vc2luZyB0byBmb2xsb3cgYQpsYXRlciB2ZXJzaW9uLgoKICAxNS4gRGlzY2xhaW1lciBvZiBXYXJyYW50eS4KCiAgVEhFUkUgSVMgTk8gV0FSUkFOVFkgRk9SIFRIRSBQUk9HUkFNLCBUTyBUSEUgRVhURU5UIFBFUk1JVFRFRCBCWQpBUFBMSUNBQkxFIExBVy4gIEVYQ0VQVCBXSEVOIE9USEVSV0lTRSBTVEFURUQgSU4gV1JJVElORyBUSEUgQ09QWVJJR0hUCkhPTERFUlMgQU5EL09SIE9USEVSIFBBUlRJRVMgUFJPVklERSBUSEUgUFJPR1JBTSAiQVMgSVMiIFdJVEhPVVQgV0FSUkFOVFkKT0YgQU5ZIEtJTkQsIEVJVEhFUiBFWFBSRVNTRUQgT1IgSU1QTElFRCwgSU5DTFVESU5HLCBCVVQgTk9UIExJTUlURUQgVE8sClRIRSBJTVBMSUVEIFdBUlJBTlRJRVMgT0YgTUVSQ0hBTlRBQklMSVRZIEFORCBGSVRORVNTIEZPUiBBIFBBUlRJQ1VMQVIKUFVSUE9TRS4gIFRIRSBFTlRJUkUgUklTSyBBUyBUTyBUSEUgUVVBTElUWSBBTkQgUEVSRk9STUFOQ0UgT0YgVEhFIFBST0dSQU0KSVMgV0lUSCBZT1UuICBTSE9VTEQgVEhFIFBST0dSQU0gUFJPVkUgREVGRUNUSVZFLCBZT1UgQVNTVU1FIFRIRSBDT1NUIE9GCkFMTCBORUNFU1NBUlkgU0VSVklDSU5HLCBSRVBBSVIgT1IgQ09SUkVDVElPTi4KCiAgMTYuIExpbWl0YXRpb24gb2YgTGlhYmlsaXR5LgoKICBJTiBOTyBFVkVOVCBVTkxFU1MgUkVRVUlSRUQgQlkgQVBQTElDQUJMRSBMQVcgT1IgQUdSRUVEIFRPIElOIFdSSVRJTkcKV0lMTCBBTlkgQ09QWVJJR0hUIEhPTERFUiwgT1IgQU5ZIE9USEVSIFBBUlRZIFdITyBNT0RJRklFUyBBTkQvT1IgQ09OVkVZUwpUSEUgUFJPR1JBTSBBUyBQRVJNSVRURUQgQUJPVkUsIEJFIExJQUJMRSBUTyBZT1UgRk9SIERBTUFHRVMsIElOQ0xVRElORyBBTlkKR0VORVJBTCwgU1BFQ0lBTCwgSU5DSURFTlRBTCBPUiBDT05TRVFVRU5USUFMIERBTUFHRVMgQVJJU0lORyBPVVQgT0YgVEhFClVTRSBPUiBJTkFCSUxJVFkgVE8gVVNFIFRIRSBQUk9HUkFNIChJTkNMVURJTkcgQlVUIE5PVCBMSU1JVEVEIFRPIExPU1MgT0YKREFUQSBPUiBEQVRBIEJFSU5HIFJFTkRFUkVEIElOQUNDVVJBVEUgT1IgTE9TU0VTIFNVU1RBSU5FRCBCWSBZT1UgT1IgVEhJUkQKUEFSVElFUyBPUiBBIEZBSUxVUkUgT0YgVEhFIFBST0dSQU0gVE8gT1BFUkFURSBXSVRIIEFOWSBPVEhFUiBQUk9HUkFNUyksCkVWRU4gSUYgU1VDSCBIT0xERVIgT1IgT1RIRVIgUEFSVFkgSEFTIEJFRU4gQURWSVNFRCBPRiBUSEUgUE9TU0lCSUxJVFkgT0YKU1VDSCBEQU1BR0VTLgoKICAxNy4gSW50ZXJwcmV0YXRpb24gb2YgU2VjdGlvbnMgMTUgYW5kIDE2LgoKICBJZiB0aGUgZGlzY2xhaW1lciBvZiB3YXJyYW50eSBhbmQgbGltaXRhdGlvbiBvZiBsaWFiaWxpdHkgcHJvdmlkZWQKYWJvdmUgY2Fubm90IGJlIGdpdmVuIGxvY2FsIGxlZ2FsIGVmZmVjdCBhY2NvcmRpbmcgdG8gdGhlaXIgdGVybXMsCnJldmlld2luZyBjb3VydHMgc2hhbGwgYXBwbHkgbG9jYWwgbGF3IHRoYXQgbW9zdCBjbG9zZWx5IGFwcHJveGltYXRlcwphbiBhYnNvbHV0ZSB3YWl2ZXIgb2YgYWxsIGNpdmlsIGxpYWJpbGl0eSBpbiBjb25uZWN0aW9uIHdpdGggdGhlClByb2dyYW0sIHVubGVzcyBhIHdhcnJhbnR5IG9yIGFzc3VtcHRpb24gb2YgbGlhYmlsaXR5IGFjY29tcGFuaWVzIGEKY29weSBvZiB0aGUgUHJvZ3JhbSBpbiByZXR1cm4gZm9yIGEgZmVlLgoKICAgICAgICAgICAgICAgICAgICAgRU5EIE9GIFRFUk1TIEFORCBDT05ESVRJT05TCgogICAgICAgICAgICBIb3cgdG8gQXBwbHkgVGhlc2UgVGVybXMgdG8gWW91ciBOZXcgUHJvZ3JhbXMKCiAgSWYgeW91IGRldmVsb3AgYSBuZXcgcHJvZ3JhbSwgYW5kIHlvdSB3YW50IGl0IHRvIGJlIG9mIHRoZSBncmVhdGVzdApwb3NzaWJsZSB1c2UgdG8gdGhlIHB1YmxpYywgdGhlIGJlc3Qgd2F5IHRvIGFjaGlldmUgdGhpcyBpcyB0byBtYWtlIGl0CmZyZWUgc29mdHdhcmUgd2hpY2ggZXZlcnlvbmUgY2FuIHJlZGlzdHJpYnV0ZSBhbmQgY2hhbmdlIHVuZGVyIHRoZXNlIHRlcm1zLgoKICBUbyBkbyBzbywgYXR0YWNoIHRoZSBmb2xsb3dpbmcgbm90aWNlcyB0byB0aGUgcHJvZ3JhbS4gIEl0IGlzIHNhZmVzdAp0byBhdHRhY2ggdGhlbSB0byB0aGUgc3RhcnQgb2YgZWFjaCBzb3VyY2UgZmlsZSB0byBtb3N0IGVmZmVjdGl2ZWx5CnN0YXRlIHRoZSBleGNsdXNpb24gb2Ygd2FycmFudHk7IGFuZCBlYWNoIGZpbGUgc2hvdWxkIGhhdmUgYXQgbGVhc3QKdGhlICJjb3B5cmlnaHQiIGxpbmUgYW5kIGEgcG9pbnRlciB0byB3aGVyZSB0aGUgZnVsbCBub3RpY2UgaXMgZm91bmQuCgogICAgPG9uZSBsaW5lIHRvIGdpdmUgdGhlIHByb2dyYW0ncyBuYW1lIGFuZCBhIGJyaWVmIGlkZWEgb2Ygd2hhdCBpdCBkb2VzLj4KICAgIENvcHlyaWdodCAoQykgPHllYXI+ICA8bmFtZSBvZiBhdXRob3I+CgogICAgVGhpcyBwcm9ncmFtIGlzIGZyZWUgc29mdHdhcmU6IHlvdSBjYW4gcmVkaXN0cmlidXRlIGl0IGFuZC9vciBtb2RpZnkKICAgIGl0IHVuZGVyIHRoZSB0ZXJtcyBvZiB0aGUgR05VIEdlbmVyYWwgUHVibGljIExpY2Vuc2UgYXMgcHVibGlzaGVkIGJ5CiAgICB0aGUgRnJlZSBTb2Z0d2FyZSBGb3VuZGF0aW9uLCBlaXRoZXIgdmVyc2lvbiAzIG9mIHRoZSBMaWNlbnNlLCBvcgogICAgKGF0IHlvdXIgb3B0aW9uKSBhbnkgbGF0ZXIgdmVyc2lvbi4KCiAgICBUaGlzIHByb2dyYW0gaXMgZGlzdHJpYnV0ZWQgaW4gdGhlIGhvcGUgdGhhdCBpdCB3aWxsIGJlIHVzZWZ1bCwKICAgIGJ1dCBXSVRIT1VUIEFOWSBXQVJSQU5UWTsgd2l0aG91dCBldmVuIHRoZSBpbXBsaWVkIHdhcnJhbnR5IG9mCiAgICBNRVJDSEFOVEFCSUxJVFkgb3IgRklUTkVTUyBGT1IgQSBQQVJUSUNVTEFSIFBVUlBPU0UuICBTZWUgdGhlCiAgICBHTlUgR2VuZXJhbCBQdWJsaWMgTGljZW5zZSBmb3IgbW9yZSBkZXRhaWxzLgoKICAgIFlvdSBzaG91bGQgaGF2ZSByZWNlaXZlZCBhIGNvcHkgb2YgdGhlIEdOVSBHZW5lcmFsIFB1YmxpYyBMaWNlbnNlCiAgICBhbG9uZyB3aXRoIHRoaXMgcHJvZ3JhbS4gIElmIG5vdCwgc2VlIDxodHRwczovL3d3dy5nbnUub3JnL2xpY2Vuc2VzLz4uCgpBbHNvIGFkZCBpbmZvcm1hdGlvbiBvbiBob3cgdG8gY29udGFjdCB5b3UgYnkgZWxlY3Ryb25pYyBhbmQgcGFwZXIgbWFpbC4KCiAgSWYgdGhlIHByb2dyYW0gZG9lcyB0ZXJtaW5hbCBpbnRlcmFjdGlvbiwgbWFrZSBpdCBvdXRwdXQgYSBzaG9ydApub3RpY2UgbGlrZSB0aGlzIHdoZW4gaXQgc3RhcnRzIGluIGFuIGludGVyYWN0aXZlIG1vZGU6CgogICAgPHByb2dyYW0+ICBDb3B5cmlnaHQgKEMpIDx5ZWFyPiAgPG5hbWUgb2YgYXV0aG9yPgogICAgVGhpcyBwcm9ncmFtIGNvbWVzIHdpdGggQUJTT0xVVEVMWSBOTyBXQVJSQU5UWTsgZm9yIGRldGFpbHMgdHlwZSBgc2hvdyB3Jy4KICAgIFRoaXMgaXMgZnJlZSBzb2Z0d2FyZSwgYW5kIHlvdSBhcmUgd2VsY29tZSB0byByZWRpc3RyaWJ1dGUgaXQKICAgIHVuZGVyIGNlcnRhaW4gY29uZGl0aW9uczsgdHlwZSBgc2hvdyBjJyBmb3IgZGV0YWlscy4KClRoZSBoeXBvdGhldGljYWwgY29tbWFuZHMgYHNob3cgdycgYW5kIGBzaG93IGMnIHNob3VsZCBzaG93IHRoZSBhcHByb3ByaWF0ZQpwYXJ0cyBvZiB0aGUgR2VuZXJhbCBQdWJsaWMgTGljZW5zZS4gIE9mIGNvdXJzZSwgeW91ciBwcm9ncmFtJ3MgY29tbWFuZHMKbWlnaHQgYmUgZGlmZmVyZW50OyBmb3IgYSBHVUkgaW50ZXJmYWNlLCB5b3Ugd291bGQgdXNlIGFuICJhYm91dCBib3giLgoKICBZb3Ugc2hvdWxkIGFsc28gZ2V0IHlvdXIgZW1wbG95ZXIgKGlmIHlvdSB3b3JrIGFzIGEgcHJvZ3JhbW1lcikgb3Igc2Nob29sLAppZiBhbnksIHRvIHNpZ24gYSAiY29weXJpZ2h0IGRpc2NsYWltZXIiIGZvciB0aGUgcHJvZ3JhbSwgaWYgbmVjZXNzYXJ5LgpGb3IgbW9yZSBpbmZvcm1hdGlvbiBvbiB0aGlzLCBhbmQgaG93IHRvIGFwcGx5IGFuZCBmb2xsb3cgdGhlIEdOVSBHUEwsIHNlZQo8aHR0cHM6Ly93d3cuZ251Lm9yZy9saWNlbnNlcy8+LgoKICBUaGUgR05VIEdlbmVyYWwgUHVibGljIExpY2Vuc2UgZG9lcyBub3QgcGVybWl0IGluY29ycG9yYXRpbmcgeW91ciBwcm9ncmFtCmludG8gcHJvcHJpZXRhcnkgcHJvZ3JhbXMuICBJZiB5b3VyIHByb2dyYW0gaXMgYSBzdWJyb3V0aW5lIGxpYnJhcnksIHlvdQptYXkgY29uc2lkZXIgaXQgbW9yZSB1c2VmdWwgdG8gcGVybWl0IGxpbmtpbmcgcHJvcHJpZXRhcnkgYXBwbGljYXRpb25zIHdpdGgKdGhlIGxpYnJhcnkuICBJZiB0aGlzIGlzIHdoYXQgeW91IHdhbnQgdG8gZG8sIHVzZSB0aGUgR05VIExlc3NlciBHZW5lcmFsClB1YmxpYyBMaWNlbnNlIGluc3RlYWQgb2YgdGhpcyBMaWNlbnNlLiAgQnV0IGZpcnN0LCBwbGVhc2UgcmVhZAo8aHR0cHM6Ly93d3cuZ251Lm9yZy9saWNlbnNlcy93aHktbm90LWxncGwuaHRtbD4uCg=="
    },
    {
      "path": "supplemental/next-path-LICENSE.txt",
      "source_locator": "Temp/elite-v336-32e87d99be6448e98b47fb0febcdf2be/license-texts/next-path-LICENSE.txt",
      "scope": "fixed-source/authority text retained V336; not an archive match",
      "origin": "RETAINED_RESEARCH_EVIDENCE",
      "payload_path": null,
      "bytes": 15978,
      "sha256": "786ef75c24eb986a2ffe31eb878b51a16826af1f5c9bd5cbdb5ad9a4223b5ab7",
      "content_base64": "Q29weXJpZ2h0IChDKSAyMDE2IFNldGggSG9sbGFkYXksIG1lQHNldGgtaG9sbGFkYXkuY29tCgpNb3ppbGxhIFB1YmxpYyBMaWNlbnNlLCB2ZXJzaW9uIDIuMAoKMS4gRGVmaW5pdGlvbnMKCjEuMS4gIkNvbnRyaWJ1dG9yIgoKICAgICBtZWFucyBlYWNoIGluZGl2aWR1YWwgb3IgbGVnYWwgZW50aXR5IHRoYXQgY3JlYXRlcywgY29udHJpYnV0ZXMgdG8gdGhlCiAgICAgY3JlYXRpb24gb2YsIG9yIG93bnMgQ292ZXJlZCBTb2Z0d2FyZS4KCjEuMi4gIkNvbnRyaWJ1dG9yIFZlcnNpb24iCgogICAgIG1lYW5zIHRoZSBjb21iaW5hdGlvbiBvZiB0aGUgQ29udHJpYnV0aW9ucyBvZiBvdGhlcnMgKGlmIGFueSkgdXNlZCBieSBhCiAgICAgQ29udHJpYnV0b3IgYW5kIHRoYXQgcGFydGljdWxhciBDb250cmlidXRvcidzIENvbnRyaWJ1dGlvbi4KCjEuMy4gIkNvbnRyaWJ1dGlvbiIKCiAgICAgbWVhbnMgQ292ZXJlZCBTb2Z0d2FyZSBvZiBhIHBhcnRpY3VsYXIgQ29udHJpYnV0b3IuCgoxLjQuICJDb3ZlcmVkIFNvZnR3YXJlIgoKICAgICBtZWFucyBTb3VyY2UgQ29kZSBGb3JtIHRvIHdoaWNoIHRoZSBpbml0aWFsIENvbnRyaWJ1dG9yIGhhcyBhdHRhY2hlZCB0aGUKICAgICBub3RpY2UgaW4gRXhoaWJpdCBBLCB0aGUgRXhlY3V0YWJsZSBGb3JtIG9mIHN1Y2ggU291cmNlIENvZGUgRm9ybSwgYW5kCiAgICAgTW9kaWZpY2F0aW9ucyBvZiBzdWNoIFNvdXJjZSBDb2RlIEZvcm0sIGluIGVhY2ggY2FzZSBpbmNsdWRpbmcgcG9ydGlvbnMKICAgICB0aGVyZW9mLgoKMS41LiAiSW5jb21wYXRpYmxlIFdpdGggU2Vjb25kYXJ5IExpY2Vuc2VzIgogICAgIG1lYW5zCgogICAgIGEuIHRoYXQgdGhlIGluaXRpYWwgQ29udHJpYnV0b3IgaGFzIGF0dGFjaGVkIHRoZSBub3RpY2UgZGVzY3JpYmVkIGluCiAgICAgICAgRXhoaWJpdCBCIHRvIHRoZSBDb3ZlcmVkIFNvZnR3YXJlOyBvcgoKICAgICBiLiB0aGF0IHRoZSBDb3ZlcmVkIFNvZnR3YXJlIHdhcyBtYWRlIGF2YWlsYWJsZSB1bmRlciB0aGUgdGVybXMgb2YKICAgICAgICB2ZXJzaW9uIDEuMSBvciBlYXJsaWVyIG9mIHRoZSBMaWNlbnNlLCBidXQgbm90IGFsc28gdW5kZXIgdGhlIHRlcm1zIG9mCiAgICAgICAgYSBTZWNvbmRhcnkgTGljZW5zZS4KCjEuNi4gIkV4ZWN1dGFibGUgRm9ybSIKCiAgICAgbWVhbnMgYW55IGZvcm0gb2YgdGhlIHdvcmsgb3RoZXIgdGhhbiBTb3VyY2UgQ29kZSBGb3JtLgoKMS43LiAiTGFyZ2VyIFdvcmsiCgogICAgIG1lYW5zIGEgd29yayB0aGF0IGNvbWJpbmVzIENvdmVyZWQgU29mdHdhcmUgd2l0aCBvdGhlciBtYXRlcmlhbCwgaW4gYQogICAgIHNlcGFyYXRlIGZpbGUgb3IgZmlsZXMsIHRoYXQgaXMgbm90IENvdmVyZWQgU29mdHdhcmUuCgoxLjguICJMaWNlbnNlIgoKICAgICBtZWFucyB0aGlzIGRvY3VtZW50LgoKMS45LiAiTGljZW5zYWJsZSIKCiAgICAgbWVhbnMgaGF2aW5nIHRoZSByaWdodCB0byBncmFudCwgdG8gdGhlIG1heGltdW0gZXh0ZW50IHBvc3NpYmxlLCB3aGV0aGVyCiAgICAgYXQgdGhlIHRpbWUgb2YgdGhlIGluaXRpYWwgZ3JhbnQgb3Igc3Vic2VxdWVudGx5LCBhbnkgYW5kIGFsbCBvZiB0aGUKICAgICByaWdodHMgY29udmV5ZWQgYnkgdGhpcyBMaWNlbnNlLgoKMS4xMC4gIk1vZGlmaWNhdGlvbnMiCgogICAgIG1lYW5zIGFueSBvZiB0aGUgZm9sbG93aW5nOgoKICAgICBhLiBhbnkgZmlsZSBpbiBTb3VyY2UgQ29kZSBGb3JtIHRoYXQgcmVzdWx0cyBmcm9tIGFuIGFkZGl0aW9uIHRvLAogICAgICAgIGRlbGV0aW9uIGZyb20sIG9yIG1vZGlmaWNhdGlvbiBvZiB0aGUgY29udGVudHMgb2YgQ292ZXJlZCBTb2Z0d2FyZTsgb3IKCiAgICAgYi4gYW55IG5ldyBmaWxlIGluIFNvdXJjZSBDb2RlIEZvcm0gdGhhdCBjb250YWlucyBhbnkgQ292ZXJlZCBTb2Z0d2FyZS4KCjEuMTEuICJQYXRlbnQgQ2xhaW1zIiBvZiBhIENvbnRyaWJ1dG9yCgogICAgICBtZWFucyBhbnkgcGF0ZW50IGNsYWltKHMpLCBpbmNsdWRpbmcgd2l0aG91dCBsaW1pdGF0aW9uLCBtZXRob2QsCiAgICAgIHByb2Nlc3MsIGFuZCBhcHBhcmF0dXMgY2xhaW1zLCBpbiBhbnkgcGF0ZW50IExpY2Vuc2FibGUgYnkgc3VjaAogICAgICBDb250cmlidXRvciB0aGF0IHdvdWxkIGJlIGluZnJpbmdlZCwgYnV0IGZvciB0aGUgZ3JhbnQgb2YgdGhlIExpY2Vuc2UsCiAgICAgIGJ5IHRoZSBtYWtpbmcsIHVzaW5nLCBzZWxsaW5nLCBvZmZlcmluZyBmb3Igc2FsZSwgaGF2aW5nIG1hZGUsIGltcG9ydCwKICAgICAgb3IgdHJhbnNmZXIgb2YgZWl0aGVyIGl0cyBDb250cmlidXRpb25zIG9yIGl0cyBDb250cmlidXRvciBWZXJzaW9uLgoKMS4xMi4gIlNlY29uZGFyeSBMaWNlbnNlIgoKICAgICAgbWVhbnMgZWl0aGVyIHRoZSBHTlUgR2VuZXJhbCBQdWJsaWMgTGljZW5zZSwgVmVyc2lvbiAyLjAsIHRoZSBHTlUgTGVzc2VyCiAgICAgIEdlbmVyYWwgUHVibGljIExpY2Vuc2UsIFZlcnNpb24gMi4xLCB0aGUgR05VIEFmZmVybyBHZW5lcmFsIFB1YmxpYwogICAgICBMaWNlbnNlLCBWZXJzaW9uIDMuMCwgb3IgYW55IGxhdGVyIHZlcnNpb25zIG9mIHRob3NlIGxpY2Vuc2VzLgoKMS4xMy4gIlNvdXJjZSBDb2RlIEZvcm0iCgogICAgICBtZWFucyB0aGUgZm9ybSBvZiB0aGUgd29yayBwcmVmZXJyZWQgZm9yIG1ha2luZyBtb2RpZmljYXRpb25zLgoKMS4xNC4gIllvdSIgKG9yICJZb3VyIikKCiAgICAgIG1lYW5zIGFuIGluZGl2aWR1YWwgb3IgYSBsZWdhbCBlbnRpdHkgZXhlcmNpc2luZyByaWdodHMgdW5kZXIgdGhpcwogICAgICBMaWNlbnNlLiBGb3IgbGVnYWwgZW50aXRpZXMsICJZb3UiIGluY2x1ZGVzIGFueSBlbnRpdHkgdGhhdCBjb250cm9scywgaXMKICAgICAgY29udHJvbGxlZCBieSwgb3IgaXMgdW5kZXIgY29tbW9uIGNvbnRyb2wgd2l0aCBZb3UuIEZvciBwdXJwb3NlcyBvZiB0aGlzCiAgICAgIGRlZmluaXRpb24sICJjb250cm9sIiBtZWFucyAoYSkgdGhlIHBvd2VyLCBkaXJlY3Qgb3IgaW5kaXJlY3QsIHRvIGNhdXNlCiAgICAgIHRoZSBkaXJlY3Rpb24gb3IgbWFuYWdlbWVudCBvZiBzdWNoIGVudGl0eSwgd2hldGhlciBieSBjb250cmFjdCBvcgogICAgICBvdGhlcndpc2UsIG9yIChiKSBvd25lcnNoaXAgb2YgbW9yZSB0aGFuIGZpZnR5IHBlcmNlbnQgKDUwJSkgb2YgdGhlCiAgICAgIG91dHN0YW5kaW5nIHNoYXJlcyBvciBiZW5lZmljaWFsIG93bmVyc2hpcCBvZiBzdWNoIGVudGl0eS4KCgoyLiBMaWNlbnNlIEdyYW50cyBhbmQgQ29uZGl0aW9ucwoKMi4xLiBHcmFudHMKCiAgICAgRWFjaCBDb250cmlidXRvciBoZXJlYnkgZ3JhbnRzIFlvdSBhIHdvcmxkLXdpZGUsIHJveWFsdHktZnJlZSwKICAgICBub24tZXhjbHVzaXZlIGxpY2Vuc2U6CgogICAgIGEuIHVuZGVyIGludGVsbGVjdHVhbCBwcm9wZXJ0eSByaWdodHMgKG90aGVyIHRoYW4gcGF0ZW50IG9yIHRyYWRlbWFyaykKICAgICAgICBMaWNlbnNhYmxlIGJ5IHN1Y2ggQ29udHJpYnV0b3IgdG8gdXNlLCByZXByb2R1Y2UsIG1ha2UgYXZhaWxhYmxlLAogICAgICAgIG1vZGlmeSwgZGlzcGxheSwgcGVyZm9ybSwgZGlzdHJpYnV0ZSwgYW5kIG90aGVyd2lzZSBleHBsb2l0IGl0cwogICAgICAgIENvbnRyaWJ1dGlvbnMsIGVpdGhlciBvbiBhbiB1bm1vZGlmaWVkIGJhc2lzLCB3aXRoIE1vZGlmaWNhdGlvbnMsIG9yCiAgICAgICAgYXMgcGFydCBvZiBhIExhcmdlciBXb3JrOyBhbmQKCiAgICAgYi4gdW5kZXIgUGF0ZW50IENsYWltcyBvZiBzdWNoIENvbnRyaWJ1dG9yIHRvIG1ha2UsIHVzZSwgc2VsbCwgb2ZmZXIgZm9yCiAgICAgICAgc2FsZSwgaGF2ZSBtYWRlLCBpbXBvcnQsIGFuZCBvdGhlcndpc2UgdHJhbnNmZXIgZWl0aGVyIGl0cwogICAgICAgIENvbnRyaWJ1dGlvbnMgb3IgaXRzIENvbnRyaWJ1dG9yIFZlcnNpb24uCgoyLjIuIEVmZmVjdGl2ZSBEYXRlCgogICAgIFRoZSBsaWNlbnNlcyBncmFudGVkIGluIFNlY3Rpb24gMi4xIHdpdGggcmVzcGVjdCB0byBhbnkgQ29udHJpYnV0aW9uCiAgICAgYmVjb21lIGVmZmVjdGl2ZSBmb3IgZWFjaCBDb250cmlidXRpb24gb24gdGhlIGRhdGUgdGhlIENvbnRyaWJ1dG9yIGZpcnN0CiAgICAgZGlzdHJpYnV0ZXMgc3VjaCBDb250cmlidXRpb24uCgoyLjMuIExpbWl0YXRpb25zIG9uIEdyYW50IFNjb3BlCgogICAgIFRoZSBsaWNlbnNlcyBncmFudGVkIGluIHRoaXMgU2VjdGlvbiAyIGFyZSB0aGUgb25seSByaWdodHMgZ3JhbnRlZCB1bmRlcgogICAgIHRoaXMgTGljZW5zZS4gTm8gYWRkaXRpb25hbCByaWdodHMgb3IgbGljZW5zZXMgd2lsbCBiZSBpbXBsaWVkIGZyb20gdGhlCiAgICAgZGlzdHJpYnV0aW9uIG9yIGxpY2Vuc2luZyBvZiBDb3ZlcmVkIFNvZnR3YXJlIHVuZGVyIHRoaXMgTGljZW5zZS4KICAgICBOb3R3aXRoc3RhbmRpbmcgU2VjdGlvbiAyLjEoYikgYWJvdmUsIG5vIHBhdGVudCBsaWNlbnNlIGlzIGdyYW50ZWQgYnkgYQogICAgIENvbnRyaWJ1dG9yOgoKICAgICBhLiBmb3IgYW55IGNvZGUgdGhhdCBhIENvbnRyaWJ1dG9yIGhhcyByZW1vdmVkIGZyb20gQ292ZXJlZCBTb2Z0d2FyZTsgb3IKCiAgICAgYi4gZm9yIGluZnJpbmdlbWVudHMgY2F1c2VkIGJ5OiAoaSkgWW91ciBhbmQgYW55IG90aGVyIHRoaXJkIHBhcnR5J3MKICAgICAgICBtb2RpZmljYXRpb25zIG9mIENvdmVyZWQgU29mdHdhcmUsIG9yIChpaSkgdGhlIGNvbWJpbmF0aW9uIG9mIGl0cwogICAgICAgIENvbnRyaWJ1dGlvbnMgd2l0aCBvdGhlciBzb2Z0d2FyZSAoZXhjZXB0IGFzIHBhcnQgb2YgaXRzIENvbnRyaWJ1dG9yCiAgICAgICAgVmVyc2lvbik7IG9yCgogICAgIGMuIHVuZGVyIFBhdGVudCBDbGFpbXMgaW5mcmluZ2VkIGJ5IENvdmVyZWQgU29mdHdhcmUgaW4gdGhlIGFic2VuY2Ugb2YKICAgICAgICBpdHMgQ29udHJpYnV0aW9ucy4KCiAgICAgVGhpcyBMaWNlbnNlIGRvZXMgbm90IGdyYW50IGFueSByaWdodHMgaW4gdGhlIHRyYWRlbWFya3MsIHNlcnZpY2UgbWFya3MsCiAgICAgb3IgbG9nb3Mgb2YgYW55IENvbnRyaWJ1dG9yIChleGNlcHQgYXMgbWF5IGJlIG5lY2Vzc2FyeSB0byBjb21wbHkgd2l0aAogICAgIHRoZSBub3RpY2UgcmVxdWlyZW1lbnRzIGluIFNlY3Rpb24gMy40KS4KCjIuNC4gU3Vic2VxdWVudCBMaWNlbnNlcwoKICAgICBObyBDb250cmlidXRvciBtYWtlcyBhZGRpdGlvbmFsIGdyYW50cyBhcyBhIHJlc3VsdCBvZiBZb3VyIGNob2ljZSB0bwogICAgIGRpc3RyaWJ1dGUgdGhlIENvdmVyZWQgU29mdHdhcmUgdW5kZXIgYSBzdWJzZXF1ZW50IHZlcnNpb24gb2YgdGhpcwogICAgIExpY2Vuc2UgKHNlZSBTZWN0aW9uIDEwLjIpIG9yIHVuZGVyIHRoZSB0ZXJtcyBvZiBhIFNlY29uZGFyeSBMaWNlbnNlIChpZgogICAgIHBlcm1pdHRlZCB1bmRlciB0aGUgdGVybXMgb2YgU2VjdGlvbiAzLjMpLgoKMi41LiBSZXByZXNlbnRhdGlvbgoKICAgICBFYWNoIENvbnRyaWJ1dG9yIHJlcHJlc2VudHMgdGhhdCB0aGUgQ29udHJpYnV0b3IgYmVsaWV2ZXMgaXRzCiAgICAgQ29udHJpYnV0aW9ucyBhcmUgaXRzIG9yaWdpbmFsIGNyZWF0aW9uKHMpIG9yIGl0IGhhcyBzdWZmaWNpZW50IHJpZ2h0cyB0bwogICAgIGdyYW50IHRoZSByaWdodHMgdG8gaXRzIENvbnRyaWJ1dGlvbnMgY29udmV5ZWQgYnkgdGhpcyBMaWNlbnNlLgoKMi42LiBGYWlyIFVzZQoKICAgICBUaGlzIExpY2Vuc2UgaXMgbm90IGludGVuZGVkIHRvIGxpbWl0IGFueSByaWdodHMgWW91IGhhdmUgdW5kZXIKICAgICBhcHBsaWNhYmxlIGNvcHlyaWdodCBkb2N0cmluZXMgb2YgZmFpciB1c2UsIGZhaXIgZGVhbGluZywgb3Igb3RoZXIKICAgICBlcXVpdmFsZW50cy4KCjIuNy4gQ29uZGl0aW9ucwoKICAgICBTZWN0aW9ucyAzLjEsIDMuMiwgMy4zLCBhbmQgMy40IGFyZSBjb25kaXRpb25zIG9mIHRoZSBsaWNlbnNlcyBncmFudGVkIGluCiAgICAgU2VjdGlvbiAyLjEuCgoKMy4gUmVzcG9uc2liaWxpdGllcwoKMy4xLiBEaXN0cmlidXRpb24gb2YgU291cmNlIEZvcm0KCiAgICAgQWxsIGRpc3RyaWJ1dGlvbiBvZiBDb3ZlcmVkIFNvZnR3YXJlIGluIFNvdXJjZSBDb2RlIEZvcm0sIGluY2x1ZGluZyBhbnkKICAgICBNb2RpZmljYXRpb25zIHRoYXQgWW91IGNyZWF0ZSBvciB0byB3aGljaCBZb3UgY29udHJpYnV0ZSwgbXVzdCBiZSB1bmRlcgogICAgIHRoZSB0ZXJtcyBvZiB0aGlzIExpY2Vuc2UuIFlvdSBtdXN0IGluZm9ybSByZWNpcGllbnRzIHRoYXQgdGhlIFNvdXJjZQogICAgIENvZGUgRm9ybSBvZiB0aGUgQ292ZXJlZCBTb2Z0d2FyZSBpcyBnb3Zlcm5lZCBieSB0aGUgdGVybXMgb2YgdGhpcwogICAgIExpY2Vuc2UsIGFuZCBob3cgdGhleSBjYW4gb2J0YWluIGEgY29weSBvZiB0aGlzIExpY2Vuc2UuIFlvdSBtYXkgbm90CiAgICAgYXR0ZW1wdCB0byBhbHRlciBvciByZXN0cmljdCB0aGUgcmVjaXBpZW50cycgcmlnaHRzIGluIHRoZSBTb3VyY2UgQ29kZQogICAgIEZvcm0uCgozLjIuIERpc3RyaWJ1dGlvbiBvZiBFeGVjdXRhYmxlIEZvcm0KCiAgICAgSWYgWW91IGRpc3RyaWJ1dGUgQ292ZXJlZCBTb2Z0d2FyZSBpbiBFeGVjdXRhYmxlIEZvcm0gdGhlbjoKCiAgICAgYS4gc3VjaCBDb3ZlcmVkIFNvZnR3YXJlIG11c3QgYWxzbyBiZSBtYWRlIGF2YWlsYWJsZSBpbiBTb3VyY2UgQ29kZSBGb3JtLAogICAgICAgIGFzIGRlc2NyaWJlZCBpbiBTZWN0aW9uIDMuMSwgYW5kIFlvdSBtdXN0IGluZm9ybSByZWNpcGllbnRzIG9mIHRoZQogICAgICAgIEV4ZWN1dGFibGUgRm9ybSBob3cgdGhleSBjYW4gb2J0YWluIGEgY29weSBvZiBzdWNoIFNvdXJjZSBDb2RlIEZvcm0gYnkKICAgICAgICByZWFzb25hYmxlIG1lYW5zIGluIGEgdGltZWx5IG1hbm5lciwgYXQgYSBjaGFyZ2Ugbm8gbW9yZSB0aGFuIHRoZSBjb3N0CiAgICAgICAgb2YgZGlzdHJpYnV0aW9uIHRvIHRoZSByZWNpcGllbnQ7IGFuZAoKICAgICBiLiBZb3UgbWF5IGRpc3RyaWJ1dGUgc3VjaCBFeGVjdXRhYmxlIEZvcm0gdW5kZXIgdGhlIHRlcm1zIG9mIHRoaXMKICAgICAgICBMaWNlbnNlLCBvciBzdWJsaWNlbnNlIGl0IHVuZGVyIGRpZmZlcmVudCB0ZXJtcywgcHJvdmlkZWQgdGhhdCB0aGUKICAgICAgICBsaWNlbnNlIGZvciB0aGUgRXhlY3V0YWJsZSBGb3JtIGRvZXMgbm90IGF0dGVtcHQgdG8gbGltaXQgb3IgYWx0ZXIgdGhlCiAgICAgICAgcmVjaXBpZW50cycgcmlnaHRzIGluIHRoZSBTb3VyY2UgQ29kZSBGb3JtIHVuZGVyIHRoaXMgTGljZW5zZS4KCjMuMy4gRGlzdHJpYnV0aW9uIG9mIGEgTGFyZ2VyIFdvcmsKCiAgICAgWW91IG1heSBjcmVhdGUgYW5kIGRpc3RyaWJ1dGUgYSBMYXJnZXIgV29yayB1bmRlciB0ZXJtcyBvZiBZb3VyIGNob2ljZSwKICAgICBwcm92aWRlZCB0aGF0IFlvdSBhbHNvIGNvbXBseSB3aXRoIHRoZSByZXF1aXJlbWVudHMgb2YgdGhpcyBMaWNlbnNlIGZvcgogICAgIHRoZSBDb3ZlcmVkIFNvZnR3YXJlLiBJZiB0aGUgTGFyZ2VyIFdvcmsgaXMgYSBjb21iaW5hdGlvbiBvZiBDb3ZlcmVkCiAgICAgU29mdHdhcmUgd2l0aCBhIHdvcmsgZ292ZXJuZWQgYnkgb25lIG9yIG1vcmUgU2Vjb25kYXJ5IExpY2Vuc2VzLCBhbmQgdGhlCiAgICAgQ292ZXJlZCBTb2Z0d2FyZSBpcyBub3QgSW5jb21wYXRpYmxlIFdpdGggU2Vjb25kYXJ5IExpY2Vuc2VzLCB0aGlzCiAgICAgTGljZW5zZSBwZXJtaXRzIFlvdSB0byBhZGRpdGlvbmFsbHkgZGlzdHJpYnV0ZSBzdWNoIENvdmVyZWQgU29mdHdhcmUKICAgICB1bmRlciB0aGUgdGVybXMgb2Ygc3VjaCBTZWNvbmRhcnkgTGljZW5zZShzKSwgc28gdGhhdCB0aGUgcmVjaXBpZW50IG9mCiAgICAgdGhlIExhcmdlciBXb3JrIG1heSwgYXQgdGhlaXIgb3B0aW9uLCBmdXJ0aGVyIGRpc3RyaWJ1dGUgdGhlIENvdmVyZWQKICAgICBTb2Z0d2FyZSB1bmRlciB0aGUgdGVybXMgb2YgZWl0aGVyIHRoaXMgTGljZW5zZSBvciBzdWNoIFNlY29uZGFyeQogICAgIExpY2Vuc2UocykuCgozLjQuIE5vdGljZXMKCiAgICAgWW91IG1heSBub3QgcmVtb3ZlIG9yIGFsdGVyIHRoZSBzdWJzdGFuY2Ugb2YgYW55IGxpY2Vuc2Ugbm90aWNlcwogICAgIChpbmNsdWRpbmcgY29weXJpZ2h0IG5vdGljZXMsIHBhdGVudCBub3RpY2VzLCBkaXNjbGFpbWVycyBvZiB3YXJyYW50eSwgb3IKICAgICBsaW1pdGF0aW9ucyBvZiBsaWFiaWxpdHkpIGNvbnRhaW5lZCB3aXRoaW4gdGhlIFNvdXJjZSBDb2RlIEZvcm0gb2YgdGhlCiAgICAgQ292ZXJlZCBTb2Z0d2FyZSwgZXhjZXB0IHRoYXQgWW91IG1heSBhbHRlciBhbnkgbGljZW5zZSBub3RpY2VzIHRvIHRoZQogICAgIGV4dGVudCByZXF1aXJlZCB0byByZW1lZHkga25vd24gZmFjdHVhbCBpbmFjY3VyYWNpZXMuCgozLjUuIEFwcGxpY2F0aW9uIG9mIEFkZGl0aW9uYWwgVGVybXMKCiAgICAgWW91IG1heSBjaG9vc2UgdG8gb2ZmZXIsIGFuZCB0byBjaGFyZ2UgYSBmZWUgZm9yLCB3YXJyYW50eSwgc3VwcG9ydCwKICAgICBpbmRlbW5pdHkgb3IgbGlhYmlsaXR5IG9ibGlnYXRpb25zIHRvIG9uZSBvciBtb3JlIHJlY2lwaWVudHMgb2YgQ292ZXJlZAogICAgIFNvZnR3YXJlLiBIb3dldmVyLCBZb3UgbWF5IGRvIHNvIG9ubHkgb24gWW91ciBvd24gYmVoYWxmLCBhbmQgbm90IG9uCiAgICAgYmVoYWxmIG9mIGFueSBDb250cmlidXRvci4gWW91IG11c3QgbWFrZSBpdCBhYnNvbHV0ZWx5IGNsZWFyIHRoYXQgYW55CiAgICAgc3VjaCB3YXJyYW50eSwgc3VwcG9ydCwgaW5kZW1uaXR5LCBvciBsaWFiaWxpdHkgb2JsaWdhdGlvbiBpcyBvZmZlcmVkIGJ5CiAgICAgWW91IGFsb25lLCBhbmQgWW91IGhlcmVieSBhZ3JlZSB0byBpbmRlbW5pZnkgZXZlcnkgQ29udHJpYnV0b3IgZm9yIGFueQogICAgIGxpYWJpbGl0eSBpbmN1cnJlZCBieSBzdWNoIENvbnRyaWJ1dG9yIGFzIGEgcmVzdWx0IG9mIHdhcnJhbnR5LCBzdXBwb3J0LAogICAgIGluZGVtbml0eSBvciBsaWFiaWxpdHkgdGVybXMgWW91IG9mZmVyLiBZb3UgbWF5IGluY2x1ZGUgYWRkaXRpb25hbAogICAgIGRpc2NsYWltZXJzIG9mIHdhcnJhbnR5IGFuZCBsaW1pdGF0aW9ucyBvZiBsaWFiaWxpdHkgc3BlY2lmaWMgdG8gYW55CiAgICAganVyaXNkaWN0aW9uLgoKNC4gSW5hYmlsaXR5IHRvIENvbXBseSBEdWUgdG8gU3RhdHV0ZSBvciBSZWd1bGF0aW9uCgogICBJZiBpdCBpcyBpbXBvc3NpYmxlIGZvciBZb3UgdG8gY29tcGx5IHdpdGggYW55IG9mIHRoZSB0ZXJtcyBvZiB0aGlzIExpY2Vuc2UKICAgd2l0aCByZXNwZWN0IHRvIHNvbWUgb3IgYWxsIG9mIHRoZSBDb3ZlcmVkIFNvZnR3YXJlIGR1ZSB0byBzdGF0dXRlLAogICBqdWRpY2lhbCBvcmRlciwgb3IgcmVndWxhdGlvbiB0aGVuIFlvdSBtdXN0OiAoYSkgY29tcGx5IHdpdGggdGhlIHRlcm1zIG9mCiAgIHRoaXMgTGljZW5zZSB0byB0aGUgbWF4aW11bSBleHRlbnQgcG9zc2libGU7IGFuZCAoYikgZGVzY3JpYmUgdGhlCiAgIGxpbWl0YXRpb25zIGFuZCB0aGUgY29kZSB0aGV5IGFmZmVjdC4gU3VjaCBkZXNjcmlwdGlvbiBtdXN0IGJlIHBsYWNlZCBpbiBhCiAgIHRleHQgZmlsZSBpbmNsdWRlZCB3aXRoIGFsbCBkaXN0cmlidXRpb25zIG9mIHRoZSBDb3ZlcmVkIFNvZnR3YXJlIHVuZGVyCiAgIHRoaXMgTGljZW5zZS4gRXhjZXB0IHRvIHRoZSBleHRlbnQgcHJvaGliaXRlZCBieSBzdGF0dXRlIG9yIHJlZ3VsYXRpb24sCiAgIHN1Y2ggZGVzY3JpcHRpb24gbXVzdCBiZSBzdWZmaWNpZW50bHkgZGV0YWlsZWQgZm9yIGEgcmVjaXBpZW50IG9mIG9yZGluYXJ5CiAgIHNraWxsIHRvIGJlIGFibGUgdG8gdW5kZXJzdGFuZCBpdC4KCjUuIFRlcm1pbmF0aW9uCgo1LjEuIFRoZSByaWdodHMgZ3JhbnRlZCB1bmRlciB0aGlzIExpY2Vuc2Ugd2lsbCB0ZXJtaW5hdGUgYXV0b21hdGljYWxseSBpZiBZb3UKICAgICBmYWlsIHRvIGNvbXBseSB3aXRoIGFueSBvZiBpdHMgdGVybXMuIEhvd2V2ZXIsIGlmIFlvdSBiZWNvbWUgY29tcGxpYW50LAogICAgIHRoZW4gdGhlIHJpZ2h0cyBncmFudGVkIHVuZGVyIHRoaXMgTGljZW5zZSBmcm9tIGEgcGFydGljdWxhciBDb250cmlidXRvcgogICAgIGFyZSByZWluc3RhdGVkIChhKSBwcm92aXNpb25hbGx5LCB1bmxlc3MgYW5kIHVudGlsIHN1Y2ggQ29udHJpYnV0b3IKICAgICBleHBsaWNpdGx5IGFuZCBmaW5hbGx5IHRlcm1pbmF0ZXMgWW91ciBncmFudHMsIGFuZCAoYikgb24gYW4gb25nb2luZwogICAgIGJhc2lzLCBpZiBzdWNoIENvbnRyaWJ1dG9yIGZhaWxzIHRvIG5vdGlmeSBZb3Ugb2YgdGhlIG5vbi1jb21wbGlhbmNlIGJ5CiAgICAgc29tZSByZWFzb25hYmxlIG1lYW5zIHByaW9yIHRvIDYwIGRheXMgYWZ0ZXIgWW91IGhhdmUgY29tZSBiYWNrIGludG8KICAgICBjb21wbGlhbmNlLiBNb3Jlb3ZlciwgWW91ciBncmFudHMgZnJvbSBhIHBhcnRpY3VsYXIgQ29udHJpYnV0b3IgYXJlCiAgICAgcmVpbnN0YXRlZCBvbiBhbiBvbmdvaW5nIGJhc2lzIGlmIHN1Y2ggQ29udHJpYnV0b3Igbm90aWZpZXMgWW91IG9mIHRoZQogICAgIG5vbi1jb21wbGlhbmNlIGJ5IHNvbWUgcmVhc29uYWJsZSBtZWFucywgdGhpcyBpcyB0aGUgZmlyc3QgdGltZSBZb3UgaGF2ZQogICAgIHJlY2VpdmVkIG5vdGljZSBvZiBub24tY29tcGxpYW5jZSB3aXRoIHRoaXMgTGljZW5zZSBmcm9tIHN1Y2gKICAgICBDb250cmlidXRvciwgYW5kIFlvdSBiZWNvbWUgY29tcGxpYW50IHByaW9yIHRvIDMwIGRheXMgYWZ0ZXIgWW91ciByZWNlaXB0CiAgICAgb2YgdGhlIG5vdGljZS4KCjUuMi4gSWYgWW91IGluaXRpYXRlIGxpdGlnYXRpb24gYWdhaW5zdCBhbnkgZW50aXR5IGJ5IGFzc2VydGluZyBhIHBhdGVudAogICAgIGluZnJpbmdlbWVudCBjbGFpbSAoZXhjbHVkaW5nIGRlY2xhcmF0b3J5IGp1ZGdtZW50IGFjdGlvbnMsCiAgICAgY291bnRlci1jbGFpbXMsIGFuZCBjcm9zcy1jbGFpbXMpIGFsbGVnaW5nIHRoYXQgYSBDb250cmlidXRvciBWZXJzaW9uCiAgICAgZGlyZWN0bHkgb3IgaW5kaXJlY3RseSBpbmZyaW5nZXMgYW55IHBhdGVudCwgdGhlbiB0aGUgcmlnaHRzIGdyYW50ZWQgdG8KICAgICBZb3UgYnkgYW55IGFuZCBhbGwgQ29udHJpYnV0b3JzIGZvciB0aGUgQ292ZXJlZCBTb2Z0d2FyZSB1bmRlciBTZWN0aW9uCiAgICAgMi4xIG9mIHRoaXMgTGljZW5zZSBzaGFsbCB0ZXJtaW5hdGUuCgo1LjMuIEluIHRoZSBldmVudCBvZiB0ZXJtaW5hdGlvbiB1bmRlciBTZWN0aW9ucyA1LjEgb3IgNS4yIGFib3ZlLCBhbGwgZW5kIHVzZXIKICAgICBsaWNlbnNlIGFncmVlbWVudHMgKGV4Y2x1ZGluZyBkaXN0cmlidXRvcnMgYW5kIHJlc2VsbGVycykgd2hpY2ggaGF2ZSBiZWVuCiAgICAgdmFsaWRseSBncmFudGVkIGJ5IFlvdSBvciBZb3VyIGRpc3RyaWJ1dG9ycyB1bmRlciB0aGlzIExpY2Vuc2UgcHJpb3IgdG8KICAgICB0ZXJtaW5hdGlvbiBzaGFsbCBzdXJ2aXZlIHRlcm1pbmF0aW9uLgoKNi4gRGlzY2xhaW1lciBvZiBXYXJyYW50eQoKICAgQ292ZXJlZCBTb2Z0d2FyZSBpcyBwcm92aWRlZCB1bmRlciB0aGlzIExpY2Vuc2Ugb24gYW4gImFzIGlzIiBiYXNpcywKICAgd2l0aG91dCB3YXJyYW50eSBvZiBhbnkga2luZCwgZWl0aGVyIGV4cHJlc3NlZCwgaW1wbGllZCwgb3Igc3RhdHV0b3J5LAogICBpbmNsdWRpbmcsIHdpdGhvdXQgbGltaXRhdGlvbiwgd2FycmFudGllcyB0aGF0IHRoZSBDb3ZlcmVkIFNvZnR3YXJlIGlzIGZyZWUKICAgb2YgZGVmZWN0cywgbWVyY2hhbnRhYmxlLCBmaXQgZm9yIGEgcGFydGljdWxhciBwdXJwb3NlIG9yIG5vbi1pbmZyaW5naW5nLgogICBUaGUgZW50aXJlIHJpc2sgYXMgdG8gdGhlIHF1YWxpdHkgYW5kIHBlcmZvcm1hbmNlIG9mIHRoZSBDb3ZlcmVkIFNvZnR3YXJlCiAgIGlzIHdpdGggWW91LiBTaG91bGQgYW55IENvdmVyZWQgU29mdHdhcmUgcHJvdmUgZGVmZWN0aXZlIGluIGFueSByZXNwZWN0LAogICBZb3UgKG5vdCBhbnkgQ29udHJpYnV0b3IpIGFzc3VtZSB0aGUgY29zdCBvZiBhbnkgbmVjZXNzYXJ5IHNlcnZpY2luZywKICAgcmVwYWlyLCBvciBjb3JyZWN0aW9uLiBUaGlzIGRpc2NsYWltZXIgb2Ygd2FycmFudHkgY29uc3RpdHV0ZXMgYW4gZXNzZW50aWFsCiAgIHBhcnQgb2YgdGhpcyBMaWNlbnNlLiBObyB1c2Ugb2YgIGFueSBDb3ZlcmVkIFNvZnR3YXJlIGlzIGF1dGhvcml6ZWQgdW5kZXIKICAgdGhpcyBMaWNlbnNlIGV4Y2VwdCB1bmRlciB0aGlzIGRpc2NsYWltZXIuCgo3LiBMaW1pdGF0aW9uIG9mIExpYWJpbGl0eQoKICAgVW5kZXIgbm8gY2lyY3Vtc3RhbmNlcyBhbmQgdW5kZXIgbm8gbGVnYWwgdGhlb3J5LCB3aGV0aGVyIHRvcnQgKGluY2x1ZGluZwogICBuZWdsaWdlbmNlKSwgY29udHJhY3QsIG9yIG90aGVyd2lzZSwgc2hhbGwgYW55IENvbnRyaWJ1dG9yLCBvciBhbnlvbmUgd2hvCiAgIGRpc3RyaWJ1dGVzIENvdmVyZWQgU29mdHdhcmUgYXMgcGVybWl0dGVkIGFib3ZlLCBiZSBsaWFibGUgdG8gWW91IGZvciBhbnkKICAgZGlyZWN0LCBpbmRpcmVjdCwgc3BlY2lhbCwgaW5jaWRlbnRhbCwgb3IgY29uc2VxdWVudGlhbCBkYW1hZ2VzIG9mIGFueQogICBjaGFyYWN0ZXIgaW5jbHVkaW5nLCB3aXRob3V0IGxpbWl0YXRpb24sIGRhbWFnZXMgZm9yIGxvc3QgcHJvZml0cywgbG9zcyBvZgogICBnb29kd2lsbCwgd29yayBzdG9wcGFnZSwgY29tcHV0ZXIgZmFpbHVyZSBvciBtYWxmdW5jdGlvbiwgb3IgYW55IGFuZCBhbGwKICAgb3RoZXIgY29tbWVyY2lhbCBkYW1hZ2VzIG9yIGxvc3NlcywgZXZlbiBpZiBzdWNoIHBhcnR5IHNoYWxsIGhhdmUgYmVlbgogICBpbmZvcm1lZCBvZiB0aGUgcG9zc2liaWxpdHkgb2Ygc3VjaCBkYW1hZ2VzLiBUaGlzIGxpbWl0YXRpb24gb2YgbGlhYmlsaXR5CiAgIHNoYWxsIG5vdCBhcHBseSB0byBsaWFiaWxpdHkgZm9yIGRlYXRoIG9yIHBlcnNvbmFsIGluanVyeSByZXN1bHRpbmcgZnJvbQogICBzdWNoIHBhcnR5J3MgbmVnbGlnZW5jZSB0byB0aGUgZXh0ZW50IGFwcGxpY2FibGUgbGF3IHByb2hpYml0cyBzdWNoCiAgIGxpbWl0YXRpb24uIFNvbWUganVyaXNkaWN0aW9ucyBkbyBub3QgYWxsb3cgdGhlIGV4Y2x1c2lvbiBvciBsaW1pdGF0aW9uIG9mCiAgIGluY2lkZW50YWwgb3IgY29uc2VxdWVudGlhbCBkYW1hZ2VzLCBzbyB0aGlzIGV4Y2x1c2lvbiBhbmQgbGltaXRhdGlvbiBtYXkKICAgbm90IGFwcGx5IHRvIFlvdS4KCjguIExpdGlnYXRpb24KCiAgIEFueSBsaXRpZ2F0aW9uIHJlbGF0aW5nIHRvIHRoaXMgTGljZW5zZSBtYXkgYmUgYnJvdWdodCBvbmx5IGluIHRoZSBjb3VydHMKICAgb2YgYSBqdXJpc2RpY3Rpb24gd2hlcmUgdGhlIGRlZmVuZGFudCBtYWludGFpbnMgaXRzIHByaW5jaXBhbCBwbGFjZSBvZgogICBidXNpbmVzcyBhbmQgc3VjaCBsaXRpZ2F0aW9uIHNoYWxsIGJlIGdvdmVybmVkIGJ5IGxhd3Mgb2YgdGhhdAogICBqdXJpc2RpY3Rpb24sIHdpdGhvdXQgcmVmZXJlbmNlIHRvIGl0cyBjb25mbGljdC1vZi1sYXcgcHJvdmlzaW9ucy4gTm90aGluZwogICBpbiB0aGlzIFNlY3Rpb24gc2hhbGwgcHJldmVudCBhIHBhcnR5J3MgYWJpbGl0eSB0byBicmluZyBjcm9zcy1jbGFpbXMgb3IKICAgY291bnRlci1jbGFpbXMuCgo5LiBNaXNjZWxsYW5lb3VzCgogICBUaGlzIExpY2Vuc2UgcmVwcmVzZW50cyB0aGUgY29tcGxldGUgYWdyZWVtZW50IGNvbmNlcm5pbmcgdGhlIHN1YmplY3QKICAgbWF0dGVyIGhlcmVvZi4gSWYgYW55IHByb3Zpc2lvbiBvZiB0aGlzIExpY2Vuc2UgaXMgaGVsZCB0byBiZQogICB1bmVuZm9yY2VhYmxlLCBzdWNoIHByb3Zpc2lvbiBzaGFsbCBiZSByZWZvcm1lZCBvbmx5IHRvIHRoZSBleHRlbnQKICAgbmVjZXNzYXJ5IHRvIG1ha2UgaXQgZW5mb3JjZWFibGUuIEFueSBsYXcgb3IgcmVndWxhdGlvbiB3aGljaCBwcm92aWRlcyB0aGF0CiAgIHRoZSBsYW5ndWFnZSBvZiBhIGNvbnRyYWN0IHNoYWxsIGJlIGNvbnN0cnVlZCBhZ2FpbnN0IHRoZSBkcmFmdGVyIHNoYWxsIG5vdAogICBiZSB1c2VkIHRvIGNvbnN0cnVlIHRoaXMgTGljZW5zZSBhZ2FpbnN0IGEgQ29udHJpYnV0b3IuCgoKMTAuIFZlcnNpb25zIG9mIHRoZSBMaWNlbnNlCgoxMC4xLiBOZXcgVmVyc2lvbnMKCiAgICAgIE1vemlsbGEgRm91bmRhdGlvbiBpcyB0aGUgbGljZW5zZSBzdGV3YXJkLiBFeGNlcHQgYXMgcHJvdmlkZWQgaW4gU2VjdGlvbgogICAgICAxMC4zLCBubyBvbmUgb3RoZXIgdGhhbiB0aGUgbGljZW5zZSBzdGV3YXJkIGhhcyB0aGUgcmlnaHQgdG8gbW9kaWZ5IG9yCiAgICAgIHB1Ymxpc2ggbmV3IHZlcnNpb25zIG9mIHRoaXMgTGljZW5zZS4gRWFjaCB2ZXJzaW9uIHdpbGwgYmUgZ2l2ZW4gYQogICAgICBkaXN0aW5ndWlzaGluZyB2ZXJzaW9uIG51bWJlci4KCjEwLjIuIEVmZmVjdCBvZiBOZXcgVmVyc2lvbnMKCiAgICAgIFlvdSBtYXkgZGlzdHJpYnV0ZSB0aGUgQ292ZXJlZCBTb2Z0d2FyZSB1bmRlciB0aGUgdGVybXMgb2YgdGhlIHZlcnNpb24KICAgICAgb2YgdGhlIExpY2Vuc2UgdW5kZXIgd2hpY2ggWW91IG9yaWdpbmFsbHkgcmVjZWl2ZWQgdGhlIENvdmVyZWQgU29mdHdhcmUsCiAgICAgIG9yIHVuZGVyIHRoZSB0ZXJtcyBvZiBhbnkgc3Vic2VxdWVudCB2ZXJzaW9uIHB1Ymxpc2hlZCBieSB0aGUgbGljZW5zZQogICAgICBzdGV3YXJkLgoKMTAuMy4gTW9kaWZpZWQgVmVyc2lvbnMKCiAgICAgIElmIHlvdSBjcmVhdGUgc29mdHdhcmUgbm90IGdvdmVybmVkIGJ5IHRoaXMgTGljZW5zZSwgYW5kIHlvdSB3YW50IHRvCiAgICAgIGNyZWF0ZSBhIG5ldyBsaWNlbnNlIGZvciBzdWNoIHNvZnR3YXJlLCB5b3UgbWF5IGNyZWF0ZSBhbmQgdXNlIGEKICAgICAgbW9kaWZpZWQgdmVyc2lvbiBvZiB0aGlzIExpY2Vuc2UgaWYgeW91IHJlbmFtZSB0aGUgbGljZW5zZSBhbmQgcmVtb3ZlCiAgICAgIGFueSByZWZlcmVuY2VzIHRvIHRoZSBuYW1lIG9mIHRoZSBsaWNlbnNlIHN0ZXdhcmQgKGV4Y2VwdCB0byBub3RlIHRoYXQKICAgICAgc3VjaCBtb2RpZmllZCBsaWNlbnNlIGRpZmZlcnMgZnJvbSB0aGlzIExpY2Vuc2UpLgoKMTAuNC4gRGlzdHJpYnV0aW5nIFNvdXJjZSBDb2RlIEZvcm0gdGhhdCBpcyBJbmNvbXBhdGlibGUgV2l0aCBTZWNvbmRhcnkKICAgICAgTGljZW5zZXMgSWYgWW91IGNob29zZSB0byBkaXN0cmlidXRlIFNvdXJjZSBDb2RlIEZvcm0gdGhhdCBpcwogICAgICBJbmNvbXBhdGlibGUgV2l0aCBTZWNvbmRhcnkgTGljZW5zZXMgdW5kZXIgdGhlIHRlcm1zIG9mIHRoaXMgdmVyc2lvbiBvZgogICAgICB0aGUgTGljZW5zZSwgdGhlIG5vdGljZSBkZXNjcmliZWQgaW4gRXhoaWJpdCBCIG9mIHRoaXMgTGljZW5zZSBtdXN0IGJlCiAgICAgIGF0dGFjaGVkLgoKRXhoaWJpdCBBIC0gU291cmNlIENvZGUgRm9ybSBMaWNlbnNlIE5vdGljZQoKICAgICAgVGhpcyBTb3VyY2UgQ29kZSBGb3JtIGlzIHN1YmplY3QgdG8gdGhlCiAgICAgIHRlcm1zIG9mIHRoZSBNb3ppbGxhIFB1YmxpYyBMaWNlbnNlLCB2LgogICAgICAyLjAuIElmIGEgY29weSBvZiB0aGUgTVBMIHdhcyBub3QKICAgICAgZGlzdHJpYnV0ZWQgd2l0aCB0aGlzIGZpbGUsIFlvdSBjYW4KICAgICAgb2J0YWluIG9uZSBhdAogICAgICBodHRwOi8vbW96aWxsYS5vcmcvTVBMLzIuMC8uCgpJZiBpdCBpcyBub3QgcG9zc2libGUgb3IgZGVzaXJhYmxlIHRvIHB1dCB0aGUgbm90aWNlIGluIGEgcGFydGljdWxhciBmaWxlLAp0aGVuIFlvdSBtYXkgaW5jbHVkZSB0aGUgbm90aWNlIGluIGEgbG9jYXRpb24gKHN1Y2ggYXMgYSBMSUNFTlNFIGZpbGUgaW4gYQpyZWxldmFudCBkaXJlY3RvcnkpIHdoZXJlIGEgcmVjaXBpZW50IHdvdWxkIGJlIGxpa2VseSB0byBsb29rIGZvciBzdWNoIGEKbm90aWNlLgoKWW91IG1heSBhZGQgYWRkaXRpb25hbCBhY2N1cmF0ZSBub3RpY2VzIG9mIGNvcHlyaWdodCBvd25lcnNoaXAuCgpFeGhpYml0IEIgLSAiSW5jb21wYXRpYmxlIFdpdGggU2Vjb25kYXJ5IExpY2Vuc2VzIiBOb3RpY2UKCiAgICAgIFRoaXMgU291cmNlIENvZGUgRm9ybSBpcyAiSW5jb21wYXRpYmxlCiAgICAgIFdpdGggU2Vjb25kYXJ5IExpY2Vuc2VzIiwgYXMgZGVmaW5lZCBieQogICAgICB0aGUgTW96aWxsYSBQdWJsaWMgTGljZW5zZSwgdi4gMi4wLgoK"
    },
    {
      "path": "supplemental/openpgp-LICENSE.txt",
      "source_locator": "Temp/elite-v336-32e87d99be6448e98b47fb0febcdf2be/license-texts/openpgp-LICENSE.txt",
      "scope": "fixed-source/authority text retained V336; not an archive match",
      "origin": "RETAINED_RESEARCH_EVIDENCE",
      "payload_path": null,
      "bytes": 7652,
      "sha256": "e3a994d82e644b03a792a930f574002658412f62407f5fee083f2555c5f23118",
      "content_base64": "ICAgICAgICAgICAgICAgICAgIEdOVSBMRVNTRVIgR0VORVJBTCBQVUJMSUMgTElDRU5TRQogICAgICAgICAgICAgICAgICAgICAgIFZlcnNpb24gMywgMjkgSnVuZSAyMDA3CgogQ29weXJpZ2h0IChDKSAyMDA3IEZyZWUgU29mdHdhcmUgRm91bmRhdGlvbiwgSW5jLiA8aHR0cHM6Ly9mc2Yub3JnLz4KIEV2ZXJ5b25lIGlzIHBlcm1pdHRlZCB0byBjb3B5IGFuZCBkaXN0cmlidXRlIHZlcmJhdGltIGNvcGllcwogb2YgdGhpcyBsaWNlbnNlIGRvY3VtZW50LCBidXQgY2hhbmdpbmcgaXQgaXMgbm90IGFsbG93ZWQuCgoKICBUaGlzIHZlcnNpb24gb2YgdGhlIEdOVSBMZXNzZXIgR2VuZXJhbCBQdWJsaWMgTGljZW5zZSBpbmNvcnBvcmF0ZXMKdGhlIHRlcm1zIGFuZCBjb25kaXRpb25zIG9mIHZlcnNpb24gMyBvZiB0aGUgR05VIEdlbmVyYWwgUHVibGljCkxpY2Vuc2UsIHN1cHBsZW1lbnRlZCBieSB0aGUgYWRkaXRpb25hbCBwZXJtaXNzaW9ucyBsaXN0ZWQgYmVsb3cuCgogIDAuIEFkZGl0aW9uYWwgRGVmaW5pdGlvbnMuCgogIEFzIHVzZWQgaGVyZWluLCAidGhpcyBMaWNlbnNlIiByZWZlcnMgdG8gdmVyc2lvbiAzIG9mIHRoZSBHTlUgTGVzc2VyCkdlbmVyYWwgUHVibGljIExpY2Vuc2UsIGFuZCB0aGUgIkdOVSBHUEwiIHJlZmVycyB0byB2ZXJzaW9uIDMgb2YgdGhlIEdOVQpHZW5lcmFsIFB1YmxpYyBMaWNlbnNlLgoKICAiVGhlIExpYnJhcnkiIHJlZmVycyB0byBhIGNvdmVyZWQgd29yayBnb3Zlcm5lZCBieSB0aGlzIExpY2Vuc2UsCm90aGVyIHRoYW4gYW4gQXBwbGljYXRpb24gb3IgYSBDb21iaW5lZCBXb3JrIGFzIGRlZmluZWQgYmVsb3cuCgogIEFuICJBcHBsaWNhdGlvbiIgaXMgYW55IHdvcmsgdGhhdCBtYWtlcyB1c2Ugb2YgYW4gaW50ZXJmYWNlIHByb3ZpZGVkCmJ5IHRoZSBMaWJyYXJ5LCBidXQgd2hpY2ggaXMgbm90IG90aGVyd2lzZSBiYXNlZCBvbiB0aGUgTGlicmFyeS4KRGVmaW5pbmcgYSBzdWJjbGFzcyBvZiBhIGNsYXNzIGRlZmluZWQgYnkgdGhlIExpYnJhcnkgaXMgZGVlbWVkIGEgbW9kZQpvZiB1c2luZyBhbiBpbnRlcmZhY2UgcHJvdmlkZWQgYnkgdGhlIExpYnJhcnkuCgogIEEgIkNvbWJpbmVkIFdvcmsiIGlzIGEgd29yayBwcm9kdWNlZCBieSBjb21iaW5pbmcgb3IgbGlua2luZyBhbgpBcHBsaWNhdGlvbiB3aXRoIHRoZSBMaWJyYXJ5LiAgVGhlIHBhcnRpY3VsYXIgdmVyc2lvbiBvZiB0aGUgTGlicmFyeQp3aXRoIHdoaWNoIHRoZSBDb21iaW5lZCBXb3JrIHdhcyBtYWRlIGlzIGFsc28gY2FsbGVkIHRoZSAiTGlua2VkClZlcnNpb24iLgoKICBUaGUgIk1pbmltYWwgQ29ycmVzcG9uZGluZyBTb3VyY2UiIGZvciBhIENvbWJpbmVkIFdvcmsgbWVhbnMgdGhlCkNvcnJlc3BvbmRpbmcgU291cmNlIGZvciB0aGUgQ29tYmluZWQgV29yaywgZXhjbHVkaW5nIGFueSBzb3VyY2UgY29kZQpmb3IgcG9ydGlvbnMgb2YgdGhlIENvbWJpbmVkIFdvcmsgdGhhdCwgY29uc2lkZXJlZCBpbiBpc29sYXRpb24sIGFyZQpiYXNlZCBvbiB0aGUgQXBwbGljYXRpb24sIGFuZCBub3Qgb24gdGhlIExpbmtlZCBWZXJzaW9uLgoKICBUaGUgIkNvcnJlc3BvbmRpbmcgQXBwbGljYXRpb24gQ29kZSIgZm9yIGEgQ29tYmluZWQgV29yayBtZWFucyB0aGUKb2JqZWN0IGNvZGUgYW5kL29yIHNvdXJjZSBjb2RlIGZvciB0aGUgQXBwbGljYXRpb24sIGluY2x1ZGluZyBhbnkgZGF0YQphbmQgdXRpbGl0eSBwcm9ncmFtcyBuZWVkZWQgZm9yIHJlcHJvZHVjaW5nIHRoZSBDb21iaW5lZCBXb3JrIGZyb20gdGhlCkFwcGxpY2F0aW9uLCBidXQgZXhjbHVkaW5nIHRoZSBTeXN0ZW0gTGlicmFyaWVzIG9mIHRoZSBDb21iaW5lZCBXb3JrLgoKICAxLiBFeGNlcHRpb24gdG8gU2VjdGlvbiAzIG9mIHRoZSBHTlUgR1BMLgoKICBZb3UgbWF5IGNvbnZleSBhIGNvdmVyZWQgd29yayB1bmRlciBzZWN0aW9ucyAzIGFuZCA0IG9mIHRoaXMgTGljZW5zZQp3aXRob3V0IGJlaW5nIGJvdW5kIGJ5IHNlY3Rpb24gMyBvZiB0aGUgR05VIEdQTC4KCiAgMi4gQ29udmV5aW5nIE1vZGlmaWVkIFZlcnNpb25zLgoKICBJZiB5b3UgbW9kaWZ5IGEgY29weSBvZiB0aGUgTGlicmFyeSwgYW5kLCBpbiB5b3VyIG1vZGlmaWNhdGlvbnMsIGEKZmFjaWxpdHkgcmVmZXJzIHRvIGEgZnVuY3Rpb24gb3IgZGF0YSB0byBiZSBzdXBwbGllZCBieSBhbiBBcHBsaWNhdGlvbgp0aGF0IHVzZXMgdGhlIGZhY2lsaXR5IChvdGhlciB0aGFuIGFzIGFuIGFyZ3VtZW50IHBhc3NlZCB3aGVuIHRoZQpmYWNpbGl0eSBpcyBpbnZva2VkKSwgdGhlbiB5b3UgbWF5IGNvbnZleSBhIGNvcHkgb2YgdGhlIG1vZGlmaWVkCnZlcnNpb246CgogICBhKSB1bmRlciB0aGlzIExpY2Vuc2UsIHByb3ZpZGVkIHRoYXQgeW91IG1ha2UgYSBnb29kIGZhaXRoIGVmZm9ydCB0bwogICBlbnN1cmUgdGhhdCwgaW4gdGhlIGV2ZW50IGFuIEFwcGxpY2F0aW9uIGRvZXMgbm90IHN1cHBseSB0aGUKICAgZnVuY3Rpb24gb3IgZGF0YSwgdGhlIGZhY2lsaXR5IHN0aWxsIG9wZXJhdGVzLCBhbmQgcGVyZm9ybXMKICAgd2hhdGV2ZXIgcGFydCBvZiBpdHMgcHVycG9zZSByZW1haW5zIG1lYW5pbmdmdWwsIG9yCgogICBiKSB1bmRlciB0aGUgR05VIEdQTCwgd2l0aCBub25lIG9mIHRoZSBhZGRpdGlvbmFsIHBlcm1pc3Npb25zIG9mCiAgIHRoaXMgTGljZW5zZSBhcHBsaWNhYmxlIHRvIHRoYXQgY29weS4KCiAgMy4gT2JqZWN0IENvZGUgSW5jb3Jwb3JhdGluZyBNYXRlcmlhbCBmcm9tIExpYnJhcnkgSGVhZGVyIEZpbGVzLgoKICBUaGUgb2JqZWN0IGNvZGUgZm9ybSBvZiBhbiBBcHBsaWNhdGlvbiBtYXkgaW5jb3Jwb3JhdGUgbWF0ZXJpYWwgZnJvbQphIGhlYWRlciBmaWxlIHRoYXQgaXMgcGFydCBvZiB0aGUgTGlicmFyeS4gIFlvdSBtYXkgY29udmV5IHN1Y2ggb2JqZWN0CmNvZGUgdW5kZXIgdGVybXMgb2YgeW91ciBjaG9pY2UsIHByb3ZpZGVkIHRoYXQsIGlmIHRoZSBpbmNvcnBvcmF0ZWQKbWF0ZXJpYWwgaXMgbm90IGxpbWl0ZWQgdG8gbnVtZXJpY2FsIHBhcmFtZXRlcnMsIGRhdGEgc3RydWN0dXJlCmxheW91dHMgYW5kIGFjY2Vzc29ycywgb3Igc21hbGwgbWFjcm9zLCBpbmxpbmUgZnVuY3Rpb25zIGFuZCB0ZW1wbGF0ZXMKKHRlbiBvciBmZXdlciBsaW5lcyBpbiBsZW5ndGgpLCB5b3UgZG8gYm90aCBvZiB0aGUgZm9sbG93aW5nOgoKICAgYSkgR2l2ZSBwcm9taW5lbnQgbm90aWNlIHdpdGggZWFjaCBjb3B5IG9mIHRoZSBvYmplY3QgY29kZSB0aGF0IHRoZQogICBMaWJyYXJ5IGlzIHVzZWQgaW4gaXQgYW5kIHRoYXQgdGhlIExpYnJhcnkgYW5kIGl0cyB1c2UgYXJlCiAgIGNvdmVyZWQgYnkgdGhpcyBMaWNlbnNlLgoKICAgYikgQWNjb21wYW55IHRoZSBvYmplY3QgY29kZSB3aXRoIGEgY29weSBvZiB0aGUgR05VIEdQTCBhbmQgdGhpcyBsaWNlbnNlCiAgIGRvY3VtZW50LgoKICA0LiBDb21iaW5lZCBXb3Jrcy4KCiAgWW91IG1heSBjb252ZXkgYSBDb21iaW5lZCBXb3JrIHVuZGVyIHRlcm1zIG9mIHlvdXIgY2hvaWNlIHRoYXQsCnRha2VuIHRvZ2V0aGVyLCBlZmZlY3RpdmVseSBkbyBub3QgcmVzdHJpY3QgbW9kaWZpY2F0aW9uIG9mIHRoZQpwb3J0aW9ucyBvZiB0aGUgTGlicmFyeSBjb250YWluZWQgaW4gdGhlIENvbWJpbmVkIFdvcmsgYW5kIHJldmVyc2UKZW5naW5lZXJpbmcgZm9yIGRlYnVnZ2luZyBzdWNoIG1vZGlmaWNhdGlvbnMsIGlmIHlvdSBhbHNvIGRvIGVhY2ggb2YKdGhlIGZvbGxvd2luZzoKCiAgIGEpIEdpdmUgcHJvbWluZW50IG5vdGljZSB3aXRoIGVhY2ggY29weSBvZiB0aGUgQ29tYmluZWQgV29yayB0aGF0CiAgIHRoZSBMaWJyYXJ5IGlzIHVzZWQgaW4gaXQgYW5kIHRoYXQgdGhlIExpYnJhcnkgYW5kIGl0cyB1c2UgYXJlCiAgIGNvdmVyZWQgYnkgdGhpcyBMaWNlbnNlLgoKICAgYikgQWNjb21wYW55IHRoZSBDb21iaW5lZCBXb3JrIHdpdGggYSBjb3B5IG9mIHRoZSBHTlUgR1BMIGFuZCB0aGlzIGxpY2Vuc2UKICAgZG9jdW1lbnQuCgogICBjKSBGb3IgYSBDb21iaW5lZCBXb3JrIHRoYXQgZGlzcGxheXMgY29weXJpZ2h0IG5vdGljZXMgZHVyaW5nCiAgIGV4ZWN1dGlvbiwgaW5jbHVkZSB0aGUgY29weXJpZ2h0IG5vdGljZSBmb3IgdGhlIExpYnJhcnkgYW1vbmcKICAgdGhlc2Ugbm90aWNlcywgYXMgd2VsbCBhcyBhIHJlZmVyZW5jZSBkaXJlY3RpbmcgdGhlIHVzZXIgdG8gdGhlCiAgIGNvcGllcyBvZiB0aGUgR05VIEdQTCBhbmQgdGhpcyBsaWNlbnNlIGRvY3VtZW50LgoKICAgZCkgRG8gb25lIG9mIHRoZSBmb2xsb3dpbmc6CgogICAgICAgMCkgQ29udmV5IHRoZSBNaW5pbWFsIENvcnJlc3BvbmRpbmcgU291cmNlIHVuZGVyIHRoZSB0ZXJtcyBvZiB0aGlzCiAgICAgICBMaWNlbnNlLCBhbmQgdGhlIENvcnJlc3BvbmRpbmcgQXBwbGljYXRpb24gQ29kZSBpbiBhIGZvcm0KICAgICAgIHN1aXRhYmxlIGZvciwgYW5kIHVuZGVyIHRlcm1zIHRoYXQgcGVybWl0LCB0aGUgdXNlciB0bwogICAgICAgcmVjb21iaW5lIG9yIHJlbGluayB0aGUgQXBwbGljYXRpb24gd2l0aCBhIG1vZGlmaWVkIHZlcnNpb24gb2YKICAgICAgIHRoZSBMaW5rZWQgVmVyc2lvbiB0byBwcm9kdWNlIGEgbW9kaWZpZWQgQ29tYmluZWQgV29yaywgaW4gdGhlCiAgICAgICBtYW5uZXIgc3BlY2lmaWVkIGJ5IHNlY3Rpb24gNiBvZiB0aGUgR05VIEdQTCBmb3IgY29udmV5aW5nCiAgICAgICBDb3JyZXNwb25kaW5nIFNvdXJjZS4KCiAgICAgICAxKSBVc2UgYSBzdWl0YWJsZSBzaGFyZWQgbGlicmFyeSBtZWNoYW5pc20gZm9yIGxpbmtpbmcgd2l0aCB0aGUKICAgICAgIExpYnJhcnkuICBBIHN1aXRhYmxlIG1lY2hhbmlzbSBpcyBvbmUgdGhhdCAoYSkgdXNlcyBhdCBydW4gdGltZQogICAgICAgYSBjb3B5IG9mIHRoZSBMaWJyYXJ5IGFscmVhZHkgcHJlc2VudCBvbiB0aGUgdXNlcidzIGNvbXB1dGVyCiAgICAgICBzeXN0ZW0sIGFuZCAoYikgd2lsbCBvcGVyYXRlIHByb3Blcmx5IHdpdGggYSBtb2RpZmllZCB2ZXJzaW9uCiAgICAgICBvZiB0aGUgTGlicmFyeSB0aGF0IGlzIGludGVyZmFjZS1jb21wYXRpYmxlIHdpdGggdGhlIExpbmtlZAogICAgICAgVmVyc2lvbi4KCiAgIGUpIFByb3ZpZGUgSW5zdGFsbGF0aW9uIEluZm9ybWF0aW9uLCBidXQgb25seSBpZiB5b3Ugd291bGQgb3RoZXJ3aXNlCiAgIGJlIHJlcXVpcmVkIHRvIHByb3ZpZGUgc3VjaCBpbmZvcm1hdGlvbiB1bmRlciBzZWN0aW9uIDYgb2YgdGhlCiAgIEdOVSBHUEwsIGFuZCBvbmx5IHRvIHRoZSBleHRlbnQgdGhhdCBzdWNoIGluZm9ybWF0aW9uIGlzCiAgIG5lY2Vzc2FyeSB0byBpbnN0YWxsIGFuZCBleGVjdXRlIGEgbW9kaWZpZWQgdmVyc2lvbiBvZiB0aGUKICAgQ29tYmluZWQgV29yayBwcm9kdWNlZCBieSByZWNvbWJpbmluZyBvciByZWxpbmtpbmcgdGhlCiAgIEFwcGxpY2F0aW9uIHdpdGggYSBtb2RpZmllZCB2ZXJzaW9uIG9mIHRoZSBMaW5rZWQgVmVyc2lvbi4gKElmCiAgIHlvdSB1c2Ugb3B0aW9uIDRkMCwgdGhlIEluc3RhbGxhdGlvbiBJbmZvcm1hdGlvbiBtdXN0IGFjY29tcGFueQogICB0aGUgTWluaW1hbCBDb3JyZXNwb25kaW5nIFNvdXJjZSBhbmQgQ29ycmVzcG9uZGluZyBBcHBsaWNhdGlvbgogICBDb2RlLiBJZiB5b3UgdXNlIG9wdGlvbiA0ZDEsIHlvdSBtdXN0IHByb3ZpZGUgdGhlIEluc3RhbGxhdGlvbgogICBJbmZvcm1hdGlvbiBpbiB0aGUgbWFubmVyIHNwZWNpZmllZCBieSBzZWN0aW9uIDYgb2YgdGhlIEdOVSBHUEwKICAgZm9yIGNvbnZleWluZyBDb3JyZXNwb25kaW5nIFNvdXJjZS4pCgogIDUuIENvbWJpbmVkIExpYnJhcmllcy4KCiAgWW91IG1heSBwbGFjZSBsaWJyYXJ5IGZhY2lsaXRpZXMgdGhhdCBhcmUgYSB3b3JrIGJhc2VkIG9uIHRoZQpMaWJyYXJ5IHNpZGUgYnkgc2lkZSBpbiBhIHNpbmdsZSBsaWJyYXJ5IHRvZ2V0aGVyIHdpdGggb3RoZXIgbGlicmFyeQpmYWNpbGl0aWVzIHRoYXQgYXJlIG5vdCBBcHBsaWNhdGlvbnMgYW5kIGFyZSBub3QgY292ZXJlZCBieSB0aGlzCkxpY2Vuc2UsIGFuZCBjb252ZXkgc3VjaCBhIGNvbWJpbmVkIGxpYnJhcnkgdW5kZXIgdGVybXMgb2YgeW91cgpjaG9pY2UsIGlmIHlvdSBkbyBib3RoIG9mIHRoZSBmb2xsb3dpbmc6CgogICBhKSBBY2NvbXBhbnkgdGhlIGNvbWJpbmVkIGxpYnJhcnkgd2l0aCBhIGNvcHkgb2YgdGhlIHNhbWUgd29yayBiYXNlZAogICBvbiB0aGUgTGlicmFyeSwgdW5jb21iaW5lZCB3aXRoIGFueSBvdGhlciBsaWJyYXJ5IGZhY2lsaXRpZXMsCiAgIGNvbnZleWVkIHVuZGVyIHRoZSB0ZXJtcyBvZiB0aGlzIExpY2Vuc2UuCgogICBiKSBHaXZlIHByb21pbmVudCBub3RpY2Ugd2l0aCB0aGUgY29tYmluZWQgbGlicmFyeSB0aGF0IHBhcnQgb2YgaXQKICAgaXMgYSB3b3JrIGJhc2VkIG9uIHRoZSBMaWJyYXJ5LCBhbmQgZXhwbGFpbmluZyB3aGVyZSB0byBmaW5kIHRoZQogICBhY2NvbXBhbnlpbmcgdW5jb21iaW5lZCBmb3JtIG9mIHRoZSBzYW1lIHdvcmsuCgogIDYuIFJldmlzZWQgVmVyc2lvbnMgb2YgdGhlIEdOVSBMZXNzZXIgR2VuZXJhbCBQdWJsaWMgTGljZW5zZS4KCiAgVGhlIEZyZWUgU29mdHdhcmUgRm91bmRhdGlvbiBtYXkgcHVibGlzaCByZXZpc2VkIGFuZC9vciBuZXcgdmVyc2lvbnMKb2YgdGhlIEdOVSBMZXNzZXIgR2VuZXJhbCBQdWJsaWMgTGljZW5zZSBmcm9tIHRpbWUgdG8gdGltZS4gU3VjaCBuZXcKdmVyc2lvbnMgd2lsbCBiZSBzaW1pbGFyIGluIHNwaXJpdCB0byB0aGUgcHJlc2VudCB2ZXJzaW9uLCBidXQgbWF5CmRpZmZlciBpbiBkZXRhaWwgdG8gYWRkcmVzcyBuZXcgcHJvYmxlbXMgb3IgY29uY2VybnMuCgogIEVhY2ggdmVyc2lvbiBpcyBnaXZlbiBhIGRpc3Rpbmd1aXNoaW5nIHZlcnNpb24gbnVtYmVyLiBJZiB0aGUKTGlicmFyeSBhcyB5b3UgcmVjZWl2ZWQgaXQgc3BlY2lmaWVzIHRoYXQgYSBjZXJ0YWluIG51bWJlcmVkIHZlcnNpb24Kb2YgdGhlIEdOVSBMZXNzZXIgR2VuZXJhbCBQdWJsaWMgTGljZW5zZSAib3IgYW55IGxhdGVyIHZlcnNpb24iCmFwcGxpZXMgdG8gaXQsIHlvdSBoYXZlIHRoZSBvcHRpb24gb2YgZm9sbG93aW5nIHRoZSB0ZXJtcyBhbmQKY29uZGl0aW9ucyBlaXRoZXIgb2YgdGhhdCBwdWJsaXNoZWQgdmVyc2lvbiBvciBvZiBhbnkgbGF0ZXIgdmVyc2lvbgpwdWJsaXNoZWQgYnkgdGhlIEZyZWUgU29mdHdhcmUgRm91bmRhdGlvbi4gSWYgdGhlIExpYnJhcnkgYXMgeW91CnJlY2VpdmVkIGl0IGRvZXMgbm90IHNwZWNpZnkgYSB2ZXJzaW9uIG51bWJlciBvZiB0aGUgR05VIExlc3NlcgpHZW5lcmFsIFB1YmxpYyBMaWNlbnNlLCB5b3UgbWF5IGNob29zZSBhbnkgdmVyc2lvbiBvZiB0aGUgR05VIExlc3NlcgpHZW5lcmFsIFB1YmxpYyBMaWNlbnNlIGV2ZXIgcHVibGlzaGVkIGJ5IHRoZSBGcmVlIFNvZnR3YXJlIEZvdW5kYXRpb24uCgogIElmIHRoZSBMaWJyYXJ5IGFzIHlvdSByZWNlaXZlZCBpdCBzcGVjaWZpZXMgdGhhdCBhIHByb3h5IGNhbiBkZWNpZGUKd2hldGhlciBmdXR1cmUgdmVyc2lvbnMgb2YgdGhlIEdOVSBMZXNzZXIgR2VuZXJhbCBQdWJsaWMgTGljZW5zZSBzaGFsbAphcHBseSwgdGhhdCBwcm94eSdzIHB1YmxpYyBzdGF0ZW1lbnQgb2YgYWNjZXB0YW5jZSBvZiBhbnkgdmVyc2lvbiBpcwpwZXJtYW5lbnQgYXV0aG9yaXphdGlvbiBmb3IgeW91IHRvIGNob29zZSB0aGF0IHZlcnNpb24gZm9yIHRoZQpMaWJyYXJ5Lgo="
    },
    {
      "path": "supplemental/spdx-exceptions-README.txt",
      "source_locator": "Temp/elite-v336-32e87d99be6448e98b47fb0febcdf2be/license-texts/spdx-exceptions-README.txt",
      "scope": "fixed-source/authority text retained V336; not an archive match",
      "origin": "RETAINED_RESEARCH_EVIDENCE",
      "payload_path": null,
      "bytes": 1243,
      "sha256": "554b19eee11d2964e9f7b244e47944c08d52ca75539260a04f3227e6c0144513",
      "content_base64": "VGhlIHBhY2thZ2UgZXhwb3J0cyBhbiBhcnJheSBvZiBzdHJpbmdzLiBFYWNoIHN0cmluZyBpcyBhbiBpZGVudGlmaWVyCmZvciBhIGxpY2Vuc2UgZXhjZXB0aW9uIHVuZGVyIHRoZSBbU29mdHdhcmUgUGFja2FnZSBEYXRhIEV4Y2hhbmdlCihTUERYKV1bU1BEWF0gc29mdHdhcmUgbGljZW5zZSBtZXRhZGF0YSBzdGFuZGFyZC4KCltTUERYXTogaHR0cHM6Ly9zcGR4Lm9yZwoKIyMgQ29weXJpZ2h0IGFuZCBMaWNlbnNpbmcKCiMjIyBTUERYCgoiU1BEWCIgaXMgYSBmZWRlcmFsbHkgcmVnaXN0ZXJlZCBVbml0ZWQgU3RhdGVzIHRyYWRlbWFyayBvZiBUaGUgTGludXgKRm91bmRhdGlvbiBDb3Jwb3JhdGlvbi4KCkZyb20gdmVyc2lvbiAyLjAgb2YgdGhlIFtTUERYXSBzcGVjaWZpY2F0aW9uOgoKPiBDb3B5cmlnaHQgwqkgMjAxMC0yMDE1IExpbnV4IEZvdW5kYXRpb24gYW5kIGl0cyBDb250cmlidXRvcnMuIExpY2Vuc2VkCj4gdW5kZXIgdGhlIENyZWF0aXZlIENvbW1vbnMgQXR0cmlidXRpb24gTGljZW5zZSAzLjAgVW5wb3J0ZWQuIEFsbCBvdGhlcgo+IHJpZ2h0cyBhcmUgZXhwcmVzc2x5IHJlc2VydmVkLgoKVGhlIExpbnV4IEZvdW5kYXRpb24gYW5kIHRoZSBTUERYIHdvcmtpbmcgZ3JvdXBzIGFyZSBnb29kIHBlb3BsZS4gT25seQp0aGV5IGRlY2lkZSB3aGF0ICJTUERYIiBtZWFucywgYXMgYSBzdGFuZGFyZCBhbmQgb3RoZXJ3aXNlLiBJIHJlc3BlY3QKdGhlaXIgd29yayBhbmQgdGhlaXIgcmlnaHRzLiBZb3Ugc2hvdWxkLCB0b28uCgojIyMgVGhpcyBQYWNrYWdlCgo+IEkgY3JlYXRlZCB0aGlzIHBhY2thZ2UgYnkgY29weWluZyBleGNlcHRpb24gaWRlbnRpZmllcnMgb3V0IG9mIHRoZQo+IFNQRFggc3BlY2lmaWNhdGlvbi4gVGhhdCB3b3JrIHdhcyBtZWNoYW5pY2FsLCByb3V0aW5lLCBhbmQgcmVxdWlyZWQgbm8KPiBjcmVhdGl2aXR5IHdoYXRzb2V2ZXIuIC0gS3lsZSBNaXRjaGVsbCwgcGFja2FnZSBhdXRob3IKClVuaXRlZCBTdGF0ZXMgdXNlcnMgY29uY2VybmVkIGFib3V0IGludGVsbGVjdHVhbCBwcm9wZXJ0eSBtYXkgd2lzaCB0bwpkaXNjdXNzIHRoZSBmb2xsb3dpbmcgU3VwcmVtZSBDb3VydCBkZWNpc2lvbnMgd2l0aCB0aGVpciBhdHRvcm5leXM6CgotIF9CYWtlciB2LiBTZWxkZW5fLCAxMDEgVS5TLiA5OSAoMTg3OSkKCi0gX0ZlaXN0IFB1YmxpY2F0aW9ucywgSW5jLiwgdi4gUnVyYWwgVGVsZXBob25lIFNlcnZpY2UgQ28uXywKICA0OTkgVS5TLiAzNDAgKDE5OTEpCg=="
    },
    {
      "path": "supplemental/individual-LICENCE.txt",
      "source_locator": "Temp/elite-v337-41c13ce2aa624c54b1d83784db789544/individual-LICENCE.txt",
      "scope": "V337 exact lifecycle archive licence or fixed-source legacy text; see original receipts",
      "origin": "RETAINED_RESEARCH_EVIDENCE",
      "payload_path": null,
      "bytes": 1050,
      "sha256": "19338d17a974d1c747b3f6c618f08975b8908fdebccbc047c0e72b09daf6ddd3",
      "content_base64": "Q29weXJpZ2h0IChjKSAyMDEyIFJheW5vcy4KClBlcm1pc3Npb24gaXMgaGVyZWJ5IGdyYW50ZWQsIGZyZWUgb2YgY2hhcmdlLCB0byBhbnkgcGVyc29uIG9idGFpbmluZyBhIGNvcHkKb2YgdGhpcyBzb2Z0d2FyZSBhbmQgYXNzb2NpYXRlZCBkb2N1bWVudGF0aW9uIGZpbGVzICh0aGUgIlNvZnR3YXJlIiksIHRvIGRlYWwKaW4gdGhlIFNvZnR3YXJlIHdpdGhvdXQgcmVzdHJpY3Rpb24sIGluY2x1ZGluZyB3aXRob3V0IGxpbWl0YXRpb24gdGhlIHJpZ2h0cwp0byB1c2UsIGNvcHksIG1vZGlmeSwgbWVyZ2UsIHB1Ymxpc2gsIGRpc3RyaWJ1dGUsIHN1YmxpY2Vuc2UsIGFuZC9vciBzZWxsCmNvcGllcyBvZiB0aGUgU29mdHdhcmUsIGFuZCB0byBwZXJtaXQgcGVyc29ucyB0byB3aG9tIHRoZSBTb2Z0d2FyZSBpcwpmdXJuaXNoZWQgdG8gZG8gc28sIHN1YmplY3QgdG8gdGhlIGZvbGxvd2luZyBjb25kaXRpb25zOgoKVGhlIGFib3ZlIGNvcHlyaWdodCBub3RpY2UgYW5kIHRoaXMgcGVybWlzc2lvbiBub3RpY2Ugc2hhbGwgYmUgaW5jbHVkZWQgaW4KYWxsIGNvcGllcyBvciBzdWJzdGFudGlhbCBwb3J0aW9ucyBvZiB0aGUgU29mdHdhcmUuCgpUSEUgU09GVFdBUkUgSVMgUFJPVklERUQgIkFTIElTIiwgV0lUSE9VVCBXQVJSQU5UWSBPRiBBTlkgS0lORCwgRVhQUkVTUyBPUgpJTVBMSUVELCBJTkNMVURJTkcgQlVUIE5PVCBMSU1JVEVEIFRPIFRIRSBXQVJSQU5USUVTIE9GIE1FUkNIQU5UQUJJTElUWSwKRklUTkVTUyBGT1IgQSBQQVJUSUNVTEFSIFBVUlBPU0UgQU5EIE5PTklORlJJTkdFTUVOVC4gSU4gTk8gRVZFTlQgU0hBTEwgVEhFCkFVVEhPUlMgT1IgQ09QWVJJR0hUIEhPTERFUlMgQkUgTElBQkxFIEZPUiBBTlkgQ0xBSU0sIERBTUFHRVMgT1IgT1RIRVIKTElBQklMSVRZLCBXSEVUSEVSIElOIEFOIEFDVElPTiBPRiBDT05UUkFDVCwgVE9SVCBPUiBPVEhFUldJU0UsIEFSSVNJTkcgRlJPTSwKT1VUIE9GIE9SIElOIENPTk5FQ1RJT04gV0lUSCBUSEUgU09GVFdBUkUgT1IgVEhFIFVTRSBPUiBPVEhFUiBERUFMSU5HUyBJTgpUSEUgU09GVFdBUkUu"
    },
    {
      "path": "supplemental/qrcode-terminal-LICENSE.txt",
      "source_locator": "Temp/elite-v337-41c13ce2aa624c54b1d83784db789544/qrcode-terminal-LICENSE.txt",
      "scope": "V337 exact lifecycle archive licence or fixed-source legacy text; see original receipts",
      "origin": "RETAINED_RESEARCH_EVIDENCE",
      "payload_path": null,
      "bytes": 11962,
      "sha256": "b3c7a2fadb2515b8106eae58439a4b9c0581a4eaa88d6a265701f8d4dd7dadb8",
      "content_base64": "CiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIEFwYWNoZSBMaWNlbnNlCiAgICAgICAgICAgICAgICAgICAgICAgICAgIFZlcnNpb24gMi4wLCBKYW51YXJ5IDIwMDQKICAgICAgICAgICAgICAgICAgICAgICAgaHR0cDovL3d3dy5hcGFjaGUub3JnL2xpY2Vuc2VzLwoKICAgVEVSTVMgQU5EIENPTkRJVElPTlMgRk9SIFVTRSwgUkVQUk9EVUNUSU9OLCBBTkQgRElTVFJJQlVUSU9OCgogICAxLiBEZWZpbml0aW9ucy4KCiAgICAgICJMaWNlbnNlIiBzaGFsbCBtZWFuIHRoZSB0ZXJtcyBhbmQgY29uZGl0aW9ucyBmb3IgdXNlLCByZXByb2R1Y3Rpb24sCiAgICAgIGFuZCBkaXN0cmlidXRpb24gYXMgZGVmaW5lZCBieSBTZWN0aW9ucyAxIHRocm91Z2ggOSBvZiB0aGlzIGRvY3VtZW50LgoKICAgICAgIkxpY2Vuc29yIiBzaGFsbCBtZWFuIHRoZSBjb3B5cmlnaHQgb3duZXIgb3IgZW50aXR5IGF1dGhvcml6ZWQgYnkKICAgICAgdGhlIGNvcHlyaWdodCBvd25lciB0aGF0IGlzIGdyYW50aW5nIHRoZSBMaWNlbnNlLgoKICAgICAgIkxlZ2FsIEVudGl0eSIgc2hhbGwgbWVhbiB0aGUgdW5pb24gb2YgdGhlIGFjdGluZyBlbnRpdHkgYW5kIGFsbAogICAgICBvdGhlciBlbnRpdGllcyB0aGF0IGNvbnRyb2wsIGFyZSBjb250cm9sbGVkIGJ5LCBvciBhcmUgdW5kZXIgY29tbW9uCiAgICAgIGNvbnRyb2wgd2l0aCB0aGF0IGVudGl0eS4gRm9yIHRoZSBwdXJwb3NlcyBvZiB0aGlzIGRlZmluaXRpb24sCiAgICAgICJjb250cm9sIiBtZWFucyAoaSkgdGhlIHBvd2VyLCBkaXJlY3Qgb3IgaW5kaXJlY3QsIHRvIGNhdXNlIHRoZQogICAgICBkaXJlY3Rpb24gb3IgbWFuYWdlbWVudCBvZiBzdWNoIGVudGl0eSwgd2hldGhlciBieSBjb250cmFjdCBvcgogICAgICBvdGhlcndpc2UsIG9yIChpaSkgb3duZXJzaGlwIG9mIGZpZnR5IHBlcmNlbnQgKDUwJSkgb3IgbW9yZSBvZiB0aGUKICAgICAgb3V0c3RhbmRpbmcgc2hhcmVzLCBvciAoaWlpKSBiZW5lZmljaWFsIG93bmVyc2hpcCBvZiBzdWNoIGVudGl0eS4KCiAgICAgICJZb3UiIChvciAiWW91ciIpIHNoYWxsIG1lYW4gYW4gaW5kaXZpZHVhbCBvciBMZWdhbCBFbnRpdHkKICAgICAgZXhlcmNpc2luZyBwZXJtaXNzaW9ucyBncmFudGVkIGJ5IHRoaXMgTGljZW5zZS4KCiAgICAgICJTb3VyY2UiIGZvcm0gc2hhbGwgbWVhbiB0aGUgcHJlZmVycmVkIGZvcm0gZm9yIG1ha2luZyBtb2RpZmljYXRpb25zLAogICAgICBpbmNsdWRpbmcgYnV0IG5vdCBsaW1pdGVkIHRvIHNvZnR3YXJlIHNvdXJjZSBjb2RlLCBkb2N1bWVudGF0aW9uCiAgICAgIHNvdXJjZSwgYW5kIGNvbmZpZ3VyYXRpb24gZmlsZXMuCgogICAgICAiT2JqZWN0IiBmb3JtIHNoYWxsIG1lYW4gYW55IGZvcm0gcmVzdWx0aW5nIGZyb20gbWVjaGFuaWNhbAogICAgICB0cmFuc2Zvcm1hdGlvbiBvciB0cmFuc2xhdGlvbiBvZiBhIFNvdXJjZSBmb3JtLCBpbmNsdWRpbmcgYnV0CiAgICAgIG5vdCBsaW1pdGVkIHRvIGNvbXBpbGVkIG9iamVjdCBjb2RlLCBnZW5lcmF0ZWQgZG9jdW1lbnRhdGlvbiwKICAgICAgYW5kIGNvbnZlcnNpb25zIHRvIG90aGVyIG1lZGlhIHR5cGVzLgoKICAgICAgIldvcmsiIHNoYWxsIG1lYW4gdGhlIHdvcmsgb2YgYXV0aG9yc2hpcCwgd2hldGhlciBpbiBTb3VyY2Ugb3IKICAgICAgT2JqZWN0IGZvcm0sIG1hZGUgYXZhaWxhYmxlIHVuZGVyIHRoZSBMaWNlbnNlLCBhcyBpbmRpY2F0ZWQgYnkgYQogICAgICBjb3B5cmlnaHQgbm90aWNlIHRoYXQgaXMgaW5jbHVkZWQgaW4gb3IgYXR0YWNoZWQgdG8gdGhlIHdvcmsKICAgICAgKGFuIGV4YW1wbGUgaXMgcHJvdmlkZWQgaW4gdGhlIEFwcGVuZGl4IGJlbG93KS4KCiAgICAgICJEZXJpdmF0aXZlIFdvcmtzIiBzaGFsbCBtZWFuIGFueSB3b3JrLCB3aGV0aGVyIGluIFNvdXJjZSBvciBPYmplY3QKICAgICAgZm9ybSwgdGhhdCBpcyBiYXNlZCBvbiAob3IgZGVyaXZlZCBmcm9tKSB0aGUgV29yayBhbmQgZm9yIHdoaWNoIHRoZQogICAgICBlZGl0b3JpYWwgcmV2aXNpb25zLCBhbm5vdGF0aW9ucywgZWxhYm9yYXRpb25zLCBvciBvdGhlciBtb2RpZmljYXRpb25zCiAgICAgIHJlcHJlc2VudCwgYXMgYSB3aG9sZSwgYW4gb3JpZ2luYWwgd29yayBvZiBhdXRob3JzaGlwLiBGb3IgdGhlIHB1cnBvc2VzCiAgICAgIG9mIHRoaXMgTGljZW5zZSwgRGVyaXZhdGl2ZSBXb3JrcyBzaGFsbCBub3QgaW5jbHVkZSB3b3JrcyB0aGF0IHJlbWFpbgogICAgICBzZXBhcmFibGUgZnJvbSwgb3IgbWVyZWx5IGxpbmsgKG9yIGJpbmQgYnkgbmFtZSkgdG8gdGhlIGludGVyZmFjZXMgb2YsCiAgICAgIHRoZSBXb3JrIGFuZCBEZXJpdmF0aXZlIFdvcmtzIHRoZXJlb2YuCgogICAgICAiQ29udHJpYnV0aW9uIiBzaGFsbCBtZWFuIGFueSB3b3JrIG9mIGF1dGhvcnNoaXAsIGluY2x1ZGluZwogICAgICB0aGUgb3JpZ2luYWwgdmVyc2lvbiBvZiB0aGUgV29yayBhbmQgYW55IG1vZGlmaWNhdGlvbnMgb3IgYWRkaXRpb25zCiAgICAgIHRvIHRoYXQgV29yayBvciBEZXJpdmF0aXZlIFdvcmtzIHRoZXJlb2YsIHRoYXQgaXMgaW50ZW50aW9uYWxseQogICAgICBzdWJtaXR0ZWQgdG8gTGljZW5zb3IgZm9yIGluY2x1c2lvbiBpbiB0aGUgV29yayBieSB0aGUgY29weXJpZ2h0IG93bmVyCiAgICAgIG9yIGJ5IGFuIGluZGl2aWR1YWwgb3IgTGVnYWwgRW50aXR5IGF1dGhvcml6ZWQgdG8gc3VibWl0IG9uIGJlaGFsZiBvZgogICAgICB0aGUgY29weXJpZ2h0IG93bmVyLiBGb3IgdGhlIHB1cnBvc2VzIG9mIHRoaXMgZGVmaW5pdGlvbiwgInN1Ym1pdHRlZCIKICAgICAgbWVhbnMgYW55IGZvcm0gb2YgZWxlY3Ryb25pYywgdmVyYmFsLCBvciB3cml0dGVuIGNvbW11bmljYXRpb24gc2VudAogICAgICB0byB0aGUgTGljZW5zb3Igb3IgaXRzIHJlcHJlc2VudGF0aXZlcywgaW5jbHVkaW5nIGJ1dCBub3QgbGltaXRlZCB0bwogICAgICBjb21tdW5pY2F0aW9uIG9uIGVsZWN0cm9uaWMgbWFpbGluZyBsaXN0cywgc291cmNlIGNvZGUgY29udHJvbCBzeXN0ZW1zLAogICAgICBhbmQgaXNzdWUgdHJhY2tpbmcgc3lzdGVtcyB0aGF0IGFyZSBtYW5hZ2VkIGJ5LCBvciBvbiBiZWhhbGYgb2YsIHRoZQogICAgICBMaWNlbnNvciBmb3IgdGhlIHB1cnBvc2Ugb2YgZGlzY3Vzc2luZyBhbmQgaW1wcm92aW5nIHRoZSBXb3JrLCBidXQKICAgICAgZXhjbHVkaW5nIGNvbW11bmljYXRpb24gdGhhdCBpcyBjb25zcGljdW91c2x5IG1hcmtlZCBvciBvdGhlcndpc2UKICAgICAgZGVzaWduYXRlZCBpbiB3cml0aW5nIGJ5IHRoZSBjb3B5cmlnaHQgb3duZXIgYXMgIk5vdCBhIENvbnRyaWJ1dGlvbi4iCgogICAgICAiQ29udHJpYnV0b3IiIHNoYWxsIG1lYW4gTGljZW5zb3IgYW5kIGFueSBpbmRpdmlkdWFsIG9yIExlZ2FsIEVudGl0eQogICAgICBvbiBiZWhhbGYgb2Ygd2hvbSBhIENvbnRyaWJ1dGlvbiBoYXMgYmVlbiByZWNlaXZlZCBieSBMaWNlbnNvciBhbmQKICAgICAgc3Vic2VxdWVudGx5IGluY29ycG9yYXRlZCB3aXRoaW4gdGhlIFdvcmsuCgogICAyLiBHcmFudCBvZiBDb3B5cmlnaHQgTGljZW5zZS4gU3ViamVjdCB0byB0aGUgdGVybXMgYW5kIGNvbmRpdGlvbnMgb2YKICAgICAgdGhpcyBMaWNlbnNlLCBlYWNoIENvbnRyaWJ1dG9yIGhlcmVieSBncmFudHMgdG8gWW91IGEgcGVycGV0dWFsLAogICAgICB3b3JsZHdpZGUsIG5vbi1leGNsdXNpdmUsIG5vLWNoYXJnZSwgcm95YWx0eS1mcmVlLCBpcnJldm9jYWJsZQogICAgICBjb3B5cmlnaHQgbGljZW5zZSB0byByZXByb2R1Y2UsIHByZXBhcmUgRGVyaXZhdGl2ZSBXb3JrcyBvZiwKICAgICAgcHVibGljbHkgZGlzcGxheSwgcHVibGljbHkgcGVyZm9ybSwgc3VibGljZW5zZSwgYW5kIGRpc3RyaWJ1dGUgdGhlCiAgICAgIFdvcmsgYW5kIHN1Y2ggRGVyaXZhdGl2ZSBXb3JrcyBpbiBTb3VyY2Ugb3IgT2JqZWN0IGZvcm0uCgogICAzLiBHcmFudCBvZiBQYXRlbnQgTGljZW5zZS4gU3ViamVjdCB0byB0aGUgdGVybXMgYW5kIGNvbmRpdGlvbnMgb2YKICAgICAgdGhpcyBMaWNlbnNlLCBlYWNoIENvbnRyaWJ1dG9yIGhlcmVieSBncmFudHMgdG8gWW91IGEgcGVycGV0dWFsLAogICAgICB3b3JsZHdpZGUsIG5vbi1leGNsdXNpdmUsIG5vLWNoYXJnZSwgcm95YWx0eS1mcmVlLCBpcnJldm9jYWJsZQogICAgICAoZXhjZXB0IGFzIHN0YXRlZCBpbiB0aGlzIHNlY3Rpb24pIHBhdGVudCBsaWNlbnNlIHRvIG1ha2UsIGhhdmUgbWFkZSwKICAgICAgdXNlLCBvZmZlciB0byBzZWxsLCBzZWxsLCBpbXBvcnQsIGFuZCBvdGhlcndpc2UgdHJhbnNmZXIgdGhlIFdvcmssCiAgICAgIHdoZXJlIHN1Y2ggbGljZW5zZSBhcHBsaWVzIG9ubHkgdG8gdGhvc2UgcGF0ZW50IGNsYWltcyBsaWNlbnNhYmxlCiAgICAgIGJ5IHN1Y2ggQ29udHJpYnV0b3IgdGhhdCBhcmUgbmVjZXNzYXJpbHkgaW5mcmluZ2VkIGJ5IHRoZWlyCiAgICAgIENvbnRyaWJ1dGlvbihzKSBhbG9uZSBvciBieSBjb21iaW5hdGlvbiBvZiB0aGVpciBDb250cmlidXRpb24ocykKICAgICAgd2l0aCB0aGUgV29yayB0byB3aGljaCBzdWNoIENvbnRyaWJ1dGlvbihzKSB3YXMgc3VibWl0dGVkLiBJZiBZb3UKICAgICAgaW5zdGl0dXRlIHBhdGVudCBsaXRpZ2F0aW9uIGFnYWluc3QgYW55IGVudGl0eSAoaW5jbHVkaW5nIGEKICAgICAgY3Jvc3MtY2xhaW0gb3IgY291bnRlcmNsYWltIGluIGEgbGF3c3VpdCkgYWxsZWdpbmcgdGhhdCB0aGUgV29yawogICAgICBvciBhIENvbnRyaWJ1dGlvbiBpbmNvcnBvcmF0ZWQgd2l0aGluIHRoZSBXb3JrIGNvbnN0aXR1dGVzIGRpcmVjdAogICAgICBvciBjb250cmlidXRvcnkgcGF0ZW50IGluZnJpbmdlbWVudCwgdGhlbiBhbnkgcGF0ZW50IGxpY2Vuc2VzCiAgICAgIGdyYW50ZWQgdG8gWW91IHVuZGVyIHRoaXMgTGljZW5zZSBmb3IgdGhhdCBXb3JrIHNoYWxsIHRlcm1pbmF0ZQogICAgICBhcyBvZiB0aGUgZGF0ZSBzdWNoIGxpdGlnYXRpb24gaXMgZmlsZWQuCgogICA0LiBSZWRpc3RyaWJ1dGlvbi4gWW91IG1heSByZXByb2R1Y2UgYW5kIGRpc3RyaWJ1dGUgY29waWVzIG9mIHRoZQogICAgICBXb3JrIG9yIERlcml2YXRpdmUgV29ya3MgdGhlcmVvZiBpbiBhbnkgbWVkaXVtLCB3aXRoIG9yIHdpdGhvdXQKICAgICAgbW9kaWZpY2F0aW9ucywgYW5kIGluIFNvdXJjZSBvciBPYmplY3QgZm9ybSwgcHJvdmlkZWQgdGhhdCBZb3UKICAgICAgbWVldCB0aGUgZm9sbG93aW5nIGNvbmRpdGlvbnM6CgogICAgICAoYSkgWW91IG11c3QgZ2l2ZSBhbnkgb3RoZXIgcmVjaXBpZW50cyBvZiB0aGUgV29yayBvcgogICAgICAgICAgRGVyaXZhdGl2ZSBXb3JrcyBhIGNvcHkgb2YgdGhpcyBMaWNlbnNlOyBhbmQKCiAgICAgIChiKSBZb3UgbXVzdCBjYXVzZSBhbnkgbW9kaWZpZWQgZmlsZXMgdG8gY2FycnkgcHJvbWluZW50IG5vdGljZXMKICAgICAgICAgIHN0YXRpbmcgdGhhdCBZb3UgY2hhbmdlZCB0aGUgZmlsZXM7IGFuZAoKICAgICAgKGMpIFlvdSBtdXN0IHJldGFpbiwgaW4gdGhlIFNvdXJjZSBmb3JtIG9mIGFueSBEZXJpdmF0aXZlIFdvcmtzCiAgICAgICAgICB0aGF0IFlvdSBkaXN0cmlidXRlLCBhbGwgY29weXJpZ2h0LCBwYXRlbnQsIHRyYWRlbWFyaywgYW5kCiAgICAgICAgICBhdHRyaWJ1dGlvbiBub3RpY2VzIGZyb20gdGhlIFNvdXJjZSBmb3JtIG9mIHRoZSBXb3JrLAogICAgICAgICAgZXhjbHVkaW5nIHRob3NlIG5vdGljZXMgdGhhdCBkbyBub3QgcGVydGFpbiB0byBhbnkgcGFydCBvZgogICAgICAgICAgdGhlIERlcml2YXRpdmUgV29ya3M7IGFuZAoKICAgICAgKGQpIElmIHRoZSBXb3JrIGluY2x1ZGVzIGEgIk5PVElDRSIgdGV4dCBmaWxlIGFzIHBhcnQgb2YgaXRzCiAgICAgICAgICBkaXN0cmlidXRpb24sIHRoZW4gYW55IERlcml2YXRpdmUgV29ya3MgdGhhdCBZb3UgZGlzdHJpYnV0ZSBtdXN0CiAgICAgICAgICBpbmNsdWRlIGEgcmVhZGFibGUgY29weSBvZiB0aGUgYXR0cmlidXRpb24gbm90aWNlcyBjb250YWluZWQKICAgICAgICAgIHdpdGhpbiBzdWNoIE5PVElDRSBmaWxlLCBleGNsdWRpbmcgdGhvc2Ugbm90aWNlcyB0aGF0IGRvIG5vdAogICAgICAgICAgcGVydGFpbiB0byBhbnkgcGFydCBvZiB0aGUgRGVyaXZhdGl2ZSBXb3JrcywgaW4gYXQgbGVhc3Qgb25lCiAgICAgICAgICBvZiB0aGUgZm9sbG93aW5nIHBsYWNlczogd2l0aGluIGEgTk9USUNFIHRleHQgZmlsZSBkaXN0cmlidXRlZAogICAgICAgICAgYXMgcGFydCBvZiB0aGUgRGVyaXZhdGl2ZSBXb3Jrczsgd2l0aGluIHRoZSBTb3VyY2UgZm9ybSBvcgogICAgICAgICAgZG9jdW1lbnRhdGlvbiwgaWYgcHJvdmlkZWQgYWxvbmcgd2l0aCB0aGUgRGVyaXZhdGl2ZSBXb3Jrczsgb3IsCiAgICAgICAgICB3aXRoaW4gYSBkaXNwbGF5IGdlbmVyYXRlZCBieSB0aGUgRGVyaXZhdGl2ZSBXb3JrcywgaWYgYW5kCiAgICAgICAgICB3aGVyZXZlciBzdWNoIHRoaXJkLXBhcnR5IG5vdGljZXMgbm9ybWFsbHkgYXBwZWFyLiBUaGUgY29udGVudHMKICAgICAgICAgIG9mIHRoZSBOT1RJQ0UgZmlsZSBhcmUgZm9yIGluZm9ybWF0aW9uYWwgcHVycG9zZXMgb25seSBhbmQKICAgICAgICAgIGRvIG5vdCBtb2RpZnkgdGhlIExpY2Vuc2UuIFlvdSBtYXkgYWRkIFlvdXIgb3duIGF0dHJpYnV0aW9uCiAgICAgICAgICBub3RpY2VzIHdpdGhpbiBEZXJpdmF0aXZlIFdvcmtzIHRoYXQgWW91IGRpc3RyaWJ1dGUsIGFsb25nc2lkZQogICAgICAgICAgb3IgYXMgYW4gYWRkZW5kdW0gdG8gdGhlIE5PVElDRSB0ZXh0IGZyb20gdGhlIFdvcmssIHByb3ZpZGVkCiAgICAgICAgICB0aGF0IHN1Y2ggYWRkaXRpb25hbCBhdHRyaWJ1dGlvbiBub3RpY2VzIGNhbm5vdCBiZSBjb25zdHJ1ZWQKICAgICAgICAgIGFzIG1vZGlmeWluZyB0aGUgTGljZW5zZS4KCiAgICAgIFlvdSBtYXkgYWRkIFlvdXIgb3duIGNvcHlyaWdodCBzdGF0ZW1lbnQgdG8gWW91ciBtb2RpZmljYXRpb25zIGFuZAogICAgICBtYXkgcHJvdmlkZSBhZGRpdGlvbmFsIG9yIGRpZmZlcmVudCBsaWNlbnNlIHRlcm1zIGFuZCBjb25kaXRpb25zCiAgICAgIGZvciB1c2UsIHJlcHJvZHVjdGlvbiwgb3IgZGlzdHJpYnV0aW9uIG9mIFlvdXIgbW9kaWZpY2F0aW9ucywgb3IKICAgICAgZm9yIGFueSBzdWNoIERlcml2YXRpdmUgV29ya3MgYXMgYSB3aG9sZSwgcHJvdmlkZWQgWW91ciB1c2UsCiAgICAgIHJlcHJvZHVjdGlvbiwgYW5kIGRpc3RyaWJ1dGlvbiBvZiB0aGUgV29yayBvdGhlcndpc2UgY29tcGxpZXMgd2l0aAogICAgICB0aGUgY29uZGl0aW9ucyBzdGF0ZWQgaW4gdGhpcyBMaWNlbnNlLgoKICAgNS4gU3VibWlzc2lvbiBvZiBDb250cmlidXRpb25zLiBVbmxlc3MgWW91IGV4cGxpY2l0bHkgc3RhdGUgb3RoZXJ3aXNlLAogICAgICBhbnkgQ29udHJpYnV0aW9uIGludGVudGlvbmFsbHkgc3VibWl0dGVkIGZvciBpbmNsdXNpb24gaW4gdGhlIFdvcmsKICAgICAgYnkgWW91IHRvIHRoZSBMaWNlbnNvciBzaGFsbCBiZSB1bmRlciB0aGUgdGVybXMgYW5kIGNvbmRpdGlvbnMgb2YKICAgICAgdGhpcyBMaWNlbnNlLCB3aXRob3V0IGFueSBhZGRpdGlvbmFsIHRlcm1zIG9yIGNvbmRpdGlvbnMuCiAgICAgIE5vdHdpdGhzdGFuZGluZyB0aGUgYWJvdmUsIG5vdGhpbmcgaGVyZWluIHNoYWxsIHN1cGVyc2VkZSBvciBtb2RpZnkKICAgICAgdGhlIHRlcm1zIG9mIGFueSBzZXBhcmF0ZSBsaWNlbnNlIGFncmVlbWVudCB5b3UgbWF5IGhhdmUgZXhlY3V0ZWQKICAgICAgd2l0aCBMaWNlbnNvciByZWdhcmRpbmcgc3VjaCBDb250cmlidXRpb25zLgoKICAgNi4gVHJhZGVtYXJrcy4gVGhpcyBMaWNlbnNlIGRvZXMgbm90IGdyYW50IHBlcm1pc3Npb24gdG8gdXNlIHRoZSB0cmFkZQogICAgICBuYW1lcywgdHJhZGVtYXJrcywgc2VydmljZSBtYXJrcywgb3IgcHJvZHVjdCBuYW1lcyBvZiB0aGUgTGljZW5zb3IsCiAgICAgIGV4Y2VwdCBhcyByZXF1aXJlZCBmb3IgcmVhc29uYWJsZSBhbmQgY3VzdG9tYXJ5IHVzZSBpbiBkZXNjcmliaW5nIHRoZQogICAgICBvcmlnaW4gb2YgdGhlIFdvcmsgYW5kIHJlcHJvZHVjaW5nIHRoZSBjb250ZW50IG9mIHRoZSBOT1RJQ0UgZmlsZS4KCiAgIDcuIERpc2NsYWltZXIgb2YgV2FycmFudHkuIFVubGVzcyByZXF1aXJlZCBieSBhcHBsaWNhYmxlIGxhdyBvcgogICAgICBhZ3JlZWQgdG8gaW4gd3JpdGluZywgTGljZW5zb3IgcHJvdmlkZXMgdGhlIFdvcmsgKGFuZCBlYWNoCiAgICAgIENvbnRyaWJ1dG9yIHByb3ZpZGVzIGl0cyBDb250cmlidXRpb25zKSBvbiBhbiAiQVMgSVMiIEJBU0lTLAogICAgICBXSVRIT1VUIFdBUlJBTlRJRVMgT1IgQ09ORElUSU9OUyBPRiBBTlkgS0lORCwgZWl0aGVyIGV4cHJlc3Mgb3IKICAgICAgaW1wbGllZCwgaW5jbHVkaW5nLCB3aXRob3V0IGxpbWl0YXRpb24sIGFueSB3YXJyYW50aWVzIG9yIGNvbmRpdGlvbnMKICAgICAgb2YgVElUTEUsIE5PTi1JTkZSSU5HRU1FTlQsIE1FUkNIQU5UQUJJTElUWSwgb3IgRklUTkVTUyBGT1IgQQogICAgICBQQVJUSUNVTEFSIFBVUlBPU0UuIFlvdSBhcmUgc29sZWx5IHJlc3BvbnNpYmxlIGZvciBkZXRlcm1pbmluZyB0aGUKICAgICAgYXBwcm9wcmlhdGVuZXNzIG9mIHVzaW5nIG9yIHJlZGlzdHJpYnV0aW5nIHRoZSBXb3JrIGFuZCBhc3N1bWUgYW55CiAgICAgIHJpc2tzIGFzc29jaWF0ZWQgd2l0aCBZb3VyIGV4ZXJjaXNlIG9mIHBlcm1pc3Npb25zIHVuZGVyIHRoaXMgTGljZW5zZS4KCiAgIDguIExpbWl0YXRpb24gb2YgTGlhYmlsaXR5LiBJbiBubyBldmVudCBhbmQgdW5kZXIgbm8gbGVnYWwgdGhlb3J5LAogICAgICB3aGV0aGVyIGluIHRvcnQgKGluY2x1ZGluZyBuZWdsaWdlbmNlKSwgY29udHJhY3QsIG9yIG90aGVyd2lzZSwKICAgICAgdW5sZXNzIHJlcXVpcmVkIGJ5IGFwcGxpY2FibGUgbGF3IChzdWNoIGFzIGRlbGliZXJhdGUgYW5kIGdyb3NzbHkKICAgICAgbmVnbGlnZW50IGFjdHMpIG9yIGFncmVlZCB0byBpbiB3cml0aW5nLCBzaGFsbCBhbnkgQ29udHJpYnV0b3IgYmUKICAgICAgbGlhYmxlIHRvIFlvdSBmb3IgZGFtYWdlcywgaW5jbHVkaW5nIGFueSBkaXJlY3QsIGluZGlyZWN0LCBzcGVjaWFsLAogICAgICBpbmNpZGVudGFsLCBvciBjb25zZXF1ZW50aWFsIGRhbWFnZXMgb2YgYW55IGNoYXJhY3RlciBhcmlzaW5nIGFzIGEKICAgICAgcmVzdWx0IG9mIHRoaXMgTGljZW5zZSBvciBvdXQgb2YgdGhlIHVzZSBvciBpbmFiaWxpdHkgdG8gdXNlIHRoZQogICAgICBXb3JrIChpbmNsdWRpbmcgYnV0IG5vdCBsaW1pdGVkIHRvIGRhbWFnZXMgZm9yIGxvc3Mgb2YgZ29vZHdpbGwsCiAgICAgIHdvcmsgc3RvcHBhZ2UsIGNvbXB1dGVyIGZhaWx1cmUgb3IgbWFsZnVuY3Rpb24sIG9yIGFueSBhbmQgYWxsCiAgICAgIG90aGVyIGNvbW1lcmNpYWwgZGFtYWdlcyBvciBsb3NzZXMpLCBldmVuIGlmIHN1Y2ggQ29udHJpYnV0b3IKICAgICAgaGFzIGJlZW4gYWR2aXNlZCBvZiB0aGUgcG9zc2liaWxpdHkgb2Ygc3VjaCBkYW1hZ2VzLgoKICAgOS4gQWNjZXB0aW5nIFdhcnJhbnR5IG9yIEFkZGl0aW9uYWwgTGlhYmlsaXR5LiBXaGlsZSByZWRpc3RyaWJ1dGluZwogICAgICB0aGUgV29yayBvciBEZXJpdmF0aXZlIFdvcmtzIHRoZXJlb2YsIFlvdSBtYXkgY2hvb3NlIHRvIG9mZmVyLAogICAgICBhbmQgY2hhcmdlIGEgZmVlIGZvciwgYWNjZXB0YW5jZSBvZiBzdXBwb3J0LCB3YXJyYW50eSwgaW5kZW1uaXR5LAogICAgICBvciBvdGhlciBsaWFiaWxpdHkgb2JsaWdhdGlvbnMgYW5kL29yIHJpZ2h0cyBjb25zaXN0ZW50IHdpdGggdGhpcwogICAgICBMaWNlbnNlLiBIb3dldmVyLCBpbiBhY2NlcHRpbmcgc3VjaCBvYmxpZ2F0aW9ucywgWW91IG1heSBhY3Qgb25seQogICAgICBvbiBZb3VyIG93biBiZWhhbGYgYW5kIG9uIFlvdXIgc29sZSByZXNwb25zaWJpbGl0eSwgbm90IG9uIGJlaGFsZgogICAgICBvZiBhbnkgb3RoZXIgQ29udHJpYnV0b3IsIGFuZCBvbmx5IGlmIFlvdSBhZ3JlZSB0byBpbmRlbW5pZnksCiAgICAgIGRlZmVuZCwgYW5kIGhvbGQgZWFjaCBDb250cmlidXRvciBoYXJtbGVzcyBmb3IgYW55IGxpYWJpbGl0eQogICAgICBpbmN1cnJlZCBieSwgb3IgY2xhaW1zIGFzc2VydGVkIGFnYWluc3QsIHN1Y2ggQ29udHJpYnV0b3IgYnkgcmVhc29uCiAgICAgIG9mIHlvdXIgYWNjZXB0aW5nIGFueSBzdWNoIHdhcnJhbnR5IG9yIGFkZGl0aW9uYWwgbGlhYmlsaXR5LgoKICAgRU5EIE9GIFRFUk1TIEFORCBDT05ESVRJT05TCgogICBBUFBFTkRJWDogSG93IHRvIGFwcGx5IHRoZSBBcGFjaGUgTGljZW5zZSB0byB5b3VyIHdvcmsuCgogICAgICBUbyBhcHBseSB0aGUgQXBhY2hlIExpY2Vuc2UgdG8geW91ciB3b3JrLCBhdHRhY2ggdGhlIGZvbGxvd2luZwogICAgICBib2lsZXJwbGF0ZSBub3RpY2UsIHdpdGggdGhlIGZpZWxkcyBlbmNsb3NlZCBieSBicmFja2V0cyAiW10iCiAgICAgIHJlcGxhY2VkIHdpdGggeW91ciBvd24gaWRlbnRpZnlpbmcgaW5mb3JtYXRpb24uIChEb24ndCBpbmNsdWRlCiAgICAgIHRoZSBicmFja2V0cyEpICBUaGUgdGV4dCBzaG91bGQgYmUgZW5jbG9zZWQgaW4gdGhlIGFwcHJvcHJpYXRlCiAgICAgIGNvbW1lbnQgc3ludGF4IGZvciB0aGUgZmlsZSBmb3JtYXQuIFdlIGFsc28gcmVjb21tZW5kIHRoYXQgYQogICAgICBmaWxlIG9yIGNsYXNzIG5hbWUgYW5kIGRlc2NyaXB0aW9uIG9mIHB1cnBvc2UgYmUgaW5jbHVkZWQgb24gdGhlCiAgICAgIHNhbWUgInByaW50ZWQgcGFnZSIgYXMgdGhlIGNvcHlyaWdodCBub3RpY2UgZm9yIGVhc2llcgogICAgICBpZGVudGlmaWNhdGlvbiB3aXRoaW4gdGhpcmQtcGFydHkgYXJjaGl2ZXMuCgogICBDb3B5cmlnaHQgW3l5eXldIFtuYW1lIG9mIGNvcHlyaWdodCBvd25lcl0KCiAgIExpY2Vuc2VkIHVuZGVyIHRoZSBBcGFjaGUgTGljZW5zZSwgVmVyc2lvbiAyLjAgKHRoZSAiTGljZW5zZSIpOwogICB5b3UgbWF5IG5vdCB1c2UgdGhpcyBmaWxlIGV4Y2VwdCBpbiBjb21wbGlhbmNlIHdpdGggdGhlIExpY2Vuc2UuCiAgIFlvdSBtYXkgb2J0YWluIGEgY29weSBvZiB0aGUgTGljZW5zZSBhdAoKICAgICAgIGh0dHA6Ly93d3cuYXBhY2hlLm9yZy9saWNlbnNlcy9MSUNFTlNFLTIuMAoKICAgVW5sZXNzIHJlcXVpcmVkIGJ5IGFwcGxpY2FibGUgbGF3IG9yIGFncmVlZCB0byBpbiB3cml0aW5nLCBzb2Z0d2FyZQogICBkaXN0cmlidXRlZCB1bmRlciB0aGUgTGljZW5zZSBpcyBkaXN0cmlidXRlZCBvbiBhbiAiQVMgSVMiIEJBU0lTLAogICBXSVRIT1VUIFdBUlJBTlRJRVMgT1IgQ09ORElUSU9OUyBPRiBBTlkgS0lORCwgZWl0aGVyIGV4cHJlc3Mgb3IgaW1wbGllZC4KICAgU2VlIHRoZSBMaWNlbnNlIGZvciB0aGUgc3BlY2lmaWMgbGFuZ3VhZ2UgZ292ZXJuaW5nIHBlcm1pc3Npb25zIGFuZAogICBsaW1pdGF0aW9ucyB1bmRlciB0aGUgTGljZW5zZS4KCj09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09ClRoaXMgcHJvZHVjdCBhbHNvIGluY2x1ZGUgdGhlIGZvbGxvd2luZyBzb2Z0d2FyZToKPT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT0KCiBRUkNvZGUgZm9yIEphdmFTY3JpcHQKCiBDb3B5cmlnaHQgKGMpIDIwMDkgS2F6dWhpa28gQXJhc2UKCiBVUkw6IGh0dHA6Ly93d3cuZC1wcm9qZWN0LmNvbS8KCiBMaWNlbnNlZCB1bmRlciB0aGUgTUlUIGxpY2Vuc2U6CiAgIGh0dHA6Ly93d3cub3BlbnNvdXJjZS5vcmcvbGljZW5zZXMvbWl0LWxpY2Vuc2UucGhwCgogVGhlIHdvcmQgIlFSIENvZGUiIGlzIHJlZ2lzdGVyZWQgdHJhZGVtYXJrIG9mIAogREVOU08gV0FWRSBJTkNPUlBPUkFURUQKICAgaHR0cDovL3d3dy5kZW5zby13YXZlLmNvbS9xcmNvZGUvZmFxcGF0ZW50LWUuaHRtbAoKTG9jYXRlZCBpbiAuL3ZlbmRvci9RUkNvZGUKLSBwcm9qZWN0IGhhcyBiZWVuIG1vZGlmaWVkIHRvIHdvcmsgaW4gTm9kZSBhbmQgc29tZSByZWZhY3RvcmluZyB3YXMgZG9uZSBmb3IgY29kZSBjbGVhbnVwCg=="
    },
    {
      "path": "supplemental/npm-lifecycle-LICENSE.txt",
      "source_locator": "Temp/elite-v337-41c13ce2aa624c54b1d83784db789544/archive-LICENSE.txt",
      "scope": "V337 exact lifecycle archive licence or fixed-source legacy text; see original receipts",
      "origin": "RETAINED_RESEARCH_EVIDENCE",
      "payload_path": null,
      "bytes": 9742,
      "sha256": "7610d223851f421d315df5e77974f1c68a04b97e02060e5bbbcf13d95e3ca257",
      "content_base64": "VGhlIG5wbSBhcHBsaWNhdGlvbgpDb3B5cmlnaHQgKGMpIG5wbSwgSW5jLiBhbmQgQ29udHJpYnV0b3JzCkxpY2Vuc2VkIG9uIHRoZSB0ZXJtcyBvZiBUaGUgQXJ0aXN0aWMgTGljZW5zZSAyLjAKCk5vZGUgcGFja2FnZSBkZXBlbmRlbmNpZXMgb2YgdGhlIG5wbSBhcHBsaWNhdGlvbgpDb3B5cmlnaHQgKGMpIHRoZWlyIHJlc3BlY3RpdmUgY29weXJpZ2h0IG93bmVycwpMaWNlbnNlZCBvbiB0aGVpciByZXNwZWN0aXZlIGxpY2Vuc2UgdGVybXMKClRoZSBucG0gcHVibGljIHJlZ2lzdHJ5IGF0IGh0dHBzOi8vcmVnaXN0cnkubnBtanMub3JnCmFuZCB0aGUgbnBtIHdlYnNpdGUgYXQgaHR0cHM6Ly93d3cubnBtanMuY29tCk9wZXJhdGVkIGJ5IG5wbSwgSW5jLgpVc2UgZ292ZXJuZWQgYnkgdGVybXMgcHVibGlzaGVkIG9uIGh0dHBzOi8vd3d3Lm5wbWpzLmNvbQoKIk5vZGUuanMiClRyYWRlbWFyayBKb3llbnQsIEluYy4sIGh0dHBzOi8vam95ZW50LmNvbQpOZWl0aGVyIG5wbSBub3IgbnBtLCBJbmMuIGFyZSBhZmZpbGlhdGVkIHdpdGggSm95ZW50LCBJbmMuCgpUaGUgTm9kZS5qcyBhcHBsaWNhdGlvbgpQcm9qZWN0IG9mIE5vZGUgRm91bmRhdGlvbiwgaHR0cHM6Ly9ub2RlanMub3JnCgpUaGUgbnBtIExvZ28KQ29weXJpZ2h0IChjKSBNYXRoaWFzIFBldHRlcnNzb24gYW5kIEJyaWFuIEhhbW1vbmQKCiJHdWJibGVidW0gQmxvY2t5IiB0eXBlZmFjZQpDb3B5cmlnaHQgKGMpIFRqYXJkYSBLb3N0ZXIsIGh0dHBzOi8vamVsbG93ZWVuLmRldmlhbnRhcnQuY29tClVzZWQgd2l0aCBwZXJtaXNzaW9uCgoKLS0tLS0tLS0KCgpUaGUgQXJ0aXN0aWMgTGljZW5zZSAyLjAKCkNvcHlyaWdodCAoYykgMjAwMC0yMDA2LCBUaGUgUGVybCBGb3VuZGF0aW9uLgoKRXZlcnlvbmUgaXMgcGVybWl0dGVkIHRvIGNvcHkgYW5kIGRpc3RyaWJ1dGUgdmVyYmF0aW0gY29waWVzCm9mIHRoaXMgbGljZW5zZSBkb2N1bWVudCwgYnV0IGNoYW5naW5nIGl0IGlzIG5vdCBhbGxvd2VkLgoKUHJlYW1ibGUKClRoaXMgbGljZW5zZSBlc3RhYmxpc2hlcyB0aGUgdGVybXMgdW5kZXIgd2hpY2ggYSBnaXZlbiBmcmVlIHNvZnR3YXJlClBhY2thZ2UgbWF5IGJlIGNvcGllZCwgbW9kaWZpZWQsIGRpc3RyaWJ1dGVkLCBhbmQvb3IgcmVkaXN0cmlidXRlZC4KVGhlIGludGVudCBpcyB0aGF0IHRoZSBDb3B5cmlnaHQgSG9sZGVyIG1haW50YWlucyBzb21lIGFydGlzdGljCmNvbnRyb2wgb3ZlciB0aGUgZGV2ZWxvcG1lbnQgb2YgdGhhdCBQYWNrYWdlIHdoaWxlIHN0aWxsIGtlZXBpbmcgdGhlClBhY2thZ2UgYXZhaWxhYmxlIGFzIG9wZW4gc291cmNlIGFuZCBmcmVlIHNvZnR3YXJlLgoKWW91IGFyZSBhbHdheXMgcGVybWl0dGVkIHRvIG1ha2UgYXJyYW5nZW1lbnRzIHdob2xseSBvdXRzaWRlIG9mIHRoaXMKbGljZW5zZSBkaXJlY3RseSB3aXRoIHRoZSBDb3B5cmlnaHQgSG9sZGVyIG9mIGEgZ2l2ZW4gUGFja2FnZS4gIElmIHRoZQp0ZXJtcyBvZiB0aGlzIGxpY2Vuc2UgZG8gbm90IHBlcm1pdCB0aGUgZnVsbCB1c2UgdGhhdCB5b3UgcHJvcG9zZSB0bwptYWtlIG9mIHRoZSBQYWNrYWdlLCB5b3Ugc2hvdWxkIGNvbnRhY3QgdGhlIENvcHlyaWdodCBIb2xkZXIgYW5kIHNlZWsKYSBkaWZmZXJlbnQgbGljZW5zaW5nIGFycmFuZ2VtZW50LgoKRGVmaW5pdGlvbnMKCiAgICAiQ29weXJpZ2h0IEhvbGRlciIgbWVhbnMgdGhlIGluZGl2aWR1YWwocykgb3Igb3JnYW5pemF0aW9uKHMpCiAgICBuYW1lZCBpbiB0aGUgY29weXJpZ2h0IG5vdGljZSBmb3IgdGhlIGVudGlyZSBQYWNrYWdlLgoKICAgICJDb250cmlidXRvciIgbWVhbnMgYW55IHBhcnR5IHRoYXQgaGFzIGNvbnRyaWJ1dGVkIGNvZGUgb3Igb3RoZXIKICAgIG1hdGVyaWFsIHRvIHRoZSBQYWNrYWdlLCBpbiBhY2NvcmRhbmNlIHdpdGggdGhlIENvcHlyaWdodCBIb2xkZXIncwogICAgcHJvY2VkdXJlcy4KCiAgICAiWW91IiBhbmQgInlvdXIiIG1lYW5zIGFueSBwZXJzb24gd2hvIHdvdWxkIGxpa2UgdG8gY29weSwKICAgIGRpc3RyaWJ1dGUsIG9yIG1vZGlmeSB0aGUgUGFja2FnZS4KCiAgICAiUGFja2FnZSIgbWVhbnMgdGhlIGNvbGxlY3Rpb24gb2YgZmlsZXMgZGlzdHJpYnV0ZWQgYnkgdGhlCiAgICBDb3B5cmlnaHQgSG9sZGVyLCBhbmQgZGVyaXZhdGl2ZXMgb2YgdGhhdCBjb2xsZWN0aW9uIGFuZC9vciBvZgogICAgdGhvc2UgZmlsZXMuIEEgZ2l2ZW4gUGFja2FnZSBtYXkgY29uc2lzdCBvZiBlaXRoZXIgdGhlIFN0YW5kYXJkCiAgICBWZXJzaW9uLCBvciBhIE1vZGlmaWVkIFZlcnNpb24uCgogICAgIkRpc3RyaWJ1dGUiIG1lYW5zIHByb3ZpZGluZyBhIGNvcHkgb2YgdGhlIFBhY2thZ2Ugb3IgbWFraW5nIGl0CiAgICBhY2Nlc3NpYmxlIHRvIGFueW9uZSBlbHNlLCBvciBpbiB0aGUgY2FzZSBvZiBhIGNvbXBhbnkgb3IKICAgIG9yZ2FuaXphdGlvbiwgdG8gb3RoZXJzIG91dHNpZGUgb2YgeW91ciBjb21wYW55IG9yIG9yZ2FuaXphdGlvbi4KCiAgICAiRGlzdHJpYnV0b3IgRmVlIiBtZWFucyBhbnkgZmVlIHRoYXQgeW91IGNoYXJnZSBmb3IgRGlzdHJpYnV0aW5nCiAgICB0aGlzIFBhY2thZ2Ugb3IgcHJvdmlkaW5nIHN1cHBvcnQgZm9yIHRoaXMgUGFja2FnZSB0byBhbm90aGVyCiAgICBwYXJ0eS4gIEl0IGRvZXMgbm90IG1lYW4gbGljZW5zaW5nIGZlZXMuCgogICAgIlN0YW5kYXJkIFZlcnNpb24iIHJlZmVycyB0byB0aGUgUGFja2FnZSBpZiBpdCBoYXMgbm90IGJlZW4KICAgIG1vZGlmaWVkLCBvciBoYXMgYmVlbiBtb2RpZmllZCBvbmx5IGluIHdheXMgZXhwbGljaXRseSByZXF1ZXN0ZWQKICAgIGJ5IHRoZSBDb3B5cmlnaHQgSG9sZGVyLgoKICAgICJNb2RpZmllZCBWZXJzaW9uIiBtZWFucyB0aGUgUGFja2FnZSwgaWYgaXQgaGFzIGJlZW4gY2hhbmdlZCwgYW5kCiAgICBzdWNoIGNoYW5nZXMgd2VyZSBub3QgZXhwbGljaXRseSByZXF1ZXN0ZWQgYnkgdGhlIENvcHlyaWdodAogICAgSG9sZGVyLgoKICAgICJPcmlnaW5hbCBMaWNlbnNlIiBtZWFucyB0aGlzIEFydGlzdGljIExpY2Vuc2UgYXMgRGlzdHJpYnV0ZWQgd2l0aAogICAgdGhlIFN0YW5kYXJkIFZlcnNpb24gb2YgdGhlIFBhY2thZ2UsIGluIGl0cyBjdXJyZW50IHZlcnNpb24gb3IgYXMKICAgIGl0IG1heSBiZSBtb2RpZmllZCBieSBUaGUgUGVybCBGb3VuZGF0aW9uIGluIHRoZSBmdXR1cmUuCgogICAgIlNvdXJjZSIgZm9ybSBtZWFucyB0aGUgc291cmNlIGNvZGUsIGRvY3VtZW50YXRpb24gc291cmNlLCBhbmQKICAgIGNvbmZpZ3VyYXRpb24gZmlsZXMgZm9yIHRoZSBQYWNrYWdlLgoKICAgICJDb21waWxlZCIgZm9ybSBtZWFucyB0aGUgY29tcGlsZWQgYnl0ZWNvZGUsIG9iamVjdCBjb2RlLCBiaW5hcnksCiAgICBvciBhbnkgb3RoZXIgZm9ybSByZXN1bHRpbmcgZnJvbSBtZWNoYW5pY2FsIHRyYW5zZm9ybWF0aW9uIG9yCiAgICB0cmFuc2xhdGlvbiBvZiB0aGUgU291cmNlIGZvcm0uCgoKUGVybWlzc2lvbiBmb3IgVXNlIGFuZCBNb2RpZmljYXRpb24gV2l0aG91dCBEaXN0cmlidXRpb24KCigxKSAgWW91IGFyZSBwZXJtaXR0ZWQgdG8gdXNlIHRoZSBTdGFuZGFyZCBWZXJzaW9uIGFuZCBjcmVhdGUgYW5kIHVzZQpNb2RpZmllZCBWZXJzaW9ucyBmb3IgYW55IHB1cnBvc2Ugd2l0aG91dCByZXN0cmljdGlvbiwgcHJvdmlkZWQgdGhhdAp5b3UgZG8gbm90IERpc3RyaWJ1dGUgdGhlIE1vZGlmaWVkIFZlcnNpb24uCgoKUGVybWlzc2lvbnMgZm9yIFJlZGlzdHJpYnV0aW9uIG9mIHRoZSBTdGFuZGFyZCBWZXJzaW9uCgooMikgIFlvdSBtYXkgRGlzdHJpYnV0ZSB2ZXJiYXRpbSBjb3BpZXMgb2YgdGhlIFNvdXJjZSBmb3JtIG9mIHRoZQpTdGFuZGFyZCBWZXJzaW9uIG9mIHRoaXMgUGFja2FnZSBpbiBhbnkgbWVkaXVtIHdpdGhvdXQgcmVzdHJpY3Rpb24sCmVpdGhlciBncmF0aXMgb3IgZm9yIGEgRGlzdHJpYnV0b3IgRmVlLCBwcm92aWRlZCB0aGF0IHlvdSBkdXBsaWNhdGUKYWxsIG9mIHRoZSBvcmlnaW5hbCBjb3B5cmlnaHQgbm90aWNlcyBhbmQgYXNzb2NpYXRlZCBkaXNjbGFpbWVycy4gIEF0CnlvdXIgZGlzY3JldGlvbiwgc3VjaCB2ZXJiYXRpbSBjb3BpZXMgbWF5IG9yIG1heSBub3QgaW5jbHVkZSBhCkNvbXBpbGVkIGZvcm0gb2YgdGhlIFBhY2thZ2UuCgooMykgIFlvdSBtYXkgYXBwbHkgYW55IGJ1ZyBmaXhlcywgcG9ydGFiaWxpdHkgY2hhbmdlcywgYW5kIG90aGVyCm1vZGlmaWNhdGlvbnMgbWFkZSBhdmFpbGFibGUgZnJvbSB0aGUgQ29weXJpZ2h0IEhvbGRlci4gIFRoZSByZXN1bHRpbmcKUGFja2FnZSB3aWxsIHN0aWxsIGJlIGNvbnNpZGVyZWQgdGhlIFN0YW5kYXJkIFZlcnNpb24sIGFuZCBhcyBzdWNoCndpbGwgYmUgc3ViamVjdCB0byB0aGUgT3JpZ2luYWwgTGljZW5zZS4KCgpEaXN0cmlidXRpb24gb2YgTW9kaWZpZWQgVmVyc2lvbnMgb2YgdGhlIFBhY2thZ2UgYXMgU291cmNlCgooNCkgIFlvdSBtYXkgRGlzdHJpYnV0ZSB5b3VyIE1vZGlmaWVkIFZlcnNpb24gYXMgU291cmNlIChlaXRoZXIgZ3JhdGlzCm9yIGZvciBhIERpc3RyaWJ1dG9yIEZlZSwgYW5kIHdpdGggb3Igd2l0aG91dCBhIENvbXBpbGVkIGZvcm0gb2YgdGhlCk1vZGlmaWVkIFZlcnNpb24pIHByb3ZpZGVkIHRoYXQgeW91IGNsZWFybHkgZG9jdW1lbnQgaG93IGl0IGRpZmZlcnMKZnJvbSB0aGUgU3RhbmRhcmQgVmVyc2lvbiwgaW5jbHVkaW5nLCBidXQgbm90IGxpbWl0ZWQgdG8sIGRvY3VtZW50aW5nCmFueSBub24tc3RhbmRhcmQgZmVhdHVyZXMsIGV4ZWN1dGFibGVzLCBvciBtb2R1bGVzLCBhbmQgcHJvdmlkZWQgdGhhdAp5b3UgZG8gYXQgbGVhc3QgT05FIG9mIHRoZSBmb2xsb3dpbmc6CgogICAgKGEpICBtYWtlIHRoZSBNb2RpZmllZCBWZXJzaW9uIGF2YWlsYWJsZSB0byB0aGUgQ29weXJpZ2h0IEhvbGRlcgogICAgb2YgdGhlIFN0YW5kYXJkIFZlcnNpb24sIHVuZGVyIHRoZSBPcmlnaW5hbCBMaWNlbnNlLCBzbyB0aGF0IHRoZQogICAgQ29weXJpZ2h0IEhvbGRlciBtYXkgaW5jbHVkZSB5b3VyIG1vZGlmaWNhdGlvbnMgaW4gdGhlIFN0YW5kYXJkCiAgICBWZXJzaW9uLgoKICAgIChiKSAgZW5zdXJlIHRoYXQgaW5zdGFsbGF0aW9uIG9mIHlvdXIgTW9kaWZpZWQgVmVyc2lvbiBkb2VzIG5vdAogICAgcHJldmVudCB0aGUgdXNlciBpbnN0YWxsaW5nIG9yIHJ1bm5pbmcgdGhlIFN0YW5kYXJkIFZlcnNpb24uIEluCiAgICBhZGRpdGlvbiwgdGhlIE1vZGlmaWVkIFZlcnNpb24gbXVzdCBiZWFyIGEgbmFtZSB0aGF0IGlzIGRpZmZlcmVudAogICAgZnJvbSB0aGUgbmFtZSBvZiB0aGUgU3RhbmRhcmQgVmVyc2lvbi4KCiAgICAoYykgIGFsbG93IGFueW9uZSB3aG8gcmVjZWl2ZXMgYSBjb3B5IG9mIHRoZSBNb2RpZmllZCBWZXJzaW9uIHRvCiAgICBtYWtlIHRoZSBTb3VyY2UgZm9ybSBvZiB0aGUgTW9kaWZpZWQgVmVyc2lvbiBhdmFpbGFibGUgdG8gb3RoZXJzCiAgICB1bmRlcgoKICAgICAgICAoaSkgIHRoZSBPcmlnaW5hbCBMaWNlbnNlIG9yCgogICAgICAgIChpaSkgIGEgbGljZW5zZSB0aGF0IHBlcm1pdHMgdGhlIGxpY2Vuc2VlIHRvIGZyZWVseSBjb3B5LAogICAgICAgIG1vZGlmeSBhbmQgcmVkaXN0cmlidXRlIHRoZSBNb2RpZmllZCBWZXJzaW9uIHVzaW5nIHRoZSBzYW1lCiAgICAgICAgbGljZW5zaW5nIHRlcm1zIHRoYXQgYXBwbHkgdG8gdGhlIGNvcHkgdGhhdCB0aGUgbGljZW5zZWUKICAgICAgICByZWNlaXZlZCwgYW5kIHJlcXVpcmVzIHRoYXQgdGhlIFNvdXJjZSBmb3JtIG9mIHRoZSBNb2RpZmllZAogICAgICAgIFZlcnNpb24sIGFuZCBvZiBhbnkgd29ya3MgZGVyaXZlZCBmcm9tIGl0LCBiZSBtYWRlIGZyZWVseQogICAgICAgIGF2YWlsYWJsZSBpbiB0aGF0IGxpY2Vuc2UgZmVlcyBhcmUgcHJvaGliaXRlZCBidXQgRGlzdHJpYnV0b3IKICAgICAgICBGZWVzIGFyZSBhbGxvd2VkLgoKCkRpc3RyaWJ1dGlvbiBvZiBDb21waWxlZCBGb3JtcyBvZiB0aGUgU3RhbmRhcmQgVmVyc2lvbgpvciBNb2RpZmllZCBWZXJzaW9ucyB3aXRob3V0IHRoZSBTb3VyY2UKCig1KSAgWW91IG1heSBEaXN0cmlidXRlIENvbXBpbGVkIGZvcm1zIG9mIHRoZSBTdGFuZGFyZCBWZXJzaW9uIHdpdGhvdXQKdGhlIFNvdXJjZSwgcHJvdmlkZWQgdGhhdCB5b3UgaW5jbHVkZSBjb21wbGV0ZSBpbnN0cnVjdGlvbnMgb24gaG93IHRvCmdldCB0aGUgU291cmNlIG9mIHRoZSBTdGFuZGFyZCBWZXJzaW9uLiAgU3VjaCBpbnN0cnVjdGlvbnMgbXVzdCBiZQp2YWxpZCBhdCB0aGUgdGltZSBvZiB5b3VyIGRpc3RyaWJ1dGlvbi4gIElmIHRoZXNlIGluc3RydWN0aW9ucywgYXQgYW55CnRpbWUgd2hpbGUgeW91IGFyZSBjYXJyeWluZyBvdXQgc3VjaCBkaXN0cmlidXRpb24sIGJlY29tZSBpbnZhbGlkLCB5b3UKbXVzdCBwcm92aWRlIG5ldyBpbnN0cnVjdGlvbnMgb24gZGVtYW5kIG9yIGNlYXNlIGZ1cnRoZXIgZGlzdHJpYnV0aW9uLgpJZiB5b3UgcHJvdmlkZSB2YWxpZCBpbnN0cnVjdGlvbnMgb3IgY2Vhc2UgZGlzdHJpYnV0aW9uIHdpdGhpbiB0aGlydHkKZGF5cyBhZnRlciB5b3UgYmVjb21lIGF3YXJlIHRoYXQgdGhlIGluc3RydWN0aW9ucyBhcmUgaW52YWxpZCwgdGhlbgp5b3UgZG8gbm90IGZvcmZlaXQgYW55IG9mIHlvdXIgcmlnaHRzIHVuZGVyIHRoaXMgbGljZW5zZS4KCig2KSAgWW91IG1heSBEaXN0cmlidXRlIGEgTW9kaWZpZWQgVmVyc2lvbiBpbiBDb21waWxlZCBmb3JtIHdpdGhvdXQKdGhlIFNvdXJjZSwgcHJvdmlkZWQgdGhhdCB5b3UgY29tcGx5IHdpdGggU2VjdGlvbiA0IHdpdGggcmVzcGVjdCB0bwp0aGUgU291cmNlIG9mIHRoZSBNb2RpZmllZCBWZXJzaW9uLgoKCkFnZ3JlZ2F0aW5nIG9yIExpbmtpbmcgdGhlIFBhY2thZ2UKCig3KSAgWW91IG1heSBhZ2dyZWdhdGUgdGhlIFBhY2thZ2UgKGVpdGhlciB0aGUgU3RhbmRhcmQgVmVyc2lvbiBvcgpNb2RpZmllZCBWZXJzaW9uKSB3aXRoIG90aGVyIHBhY2thZ2VzIGFuZCBEaXN0cmlidXRlIHRoZSByZXN1bHRpbmcKYWdncmVnYXRpb24gcHJvdmlkZWQgdGhhdCB5b3UgZG8gbm90IGNoYXJnZSBhIGxpY2Vuc2luZyBmZWUgZm9yIHRoZQpQYWNrYWdlLiAgRGlzdHJpYnV0b3IgRmVlcyBhcmUgcGVybWl0dGVkLCBhbmQgbGljZW5zaW5nIGZlZXMgZm9yIG90aGVyCmNvbXBvbmVudHMgaW4gdGhlIGFnZ3JlZ2F0aW9uIGFyZSBwZXJtaXR0ZWQuIFRoZSB0ZXJtcyBvZiB0aGlzIGxpY2Vuc2UKYXBwbHkgdG8gdGhlIHVzZSBhbmQgRGlzdHJpYnV0aW9uIG9mIHRoZSBTdGFuZGFyZCBvciBNb2RpZmllZCBWZXJzaW9ucwphcyBpbmNsdWRlZCBpbiB0aGUgYWdncmVnYXRpb24uCgooOCkgWW91IGFyZSBwZXJtaXR0ZWQgdG8gbGluayBNb2RpZmllZCBhbmQgU3RhbmRhcmQgVmVyc2lvbnMgd2l0aApvdGhlciB3b3JrcywgdG8gZW1iZWQgdGhlIFBhY2thZ2UgaW4gYSBsYXJnZXIgd29yayBvZiB5b3VyIG93biwgb3IgdG8KYnVpbGQgc3RhbmQtYWxvbmUgYmluYXJ5IG9yIGJ5dGVjb2RlIHZlcnNpb25zIG9mIGFwcGxpY2F0aW9ucyB0aGF0CmluY2x1ZGUgdGhlIFBhY2thZ2UsIGFuZCBEaXN0cmlidXRlIHRoZSByZXN1bHQgd2l0aG91dCByZXN0cmljdGlvbiwKcHJvdmlkZWQgdGhlIHJlc3VsdCBkb2VzIG5vdCBleHBvc2UgYSBkaXJlY3QgaW50ZXJmYWNlIHRvIHRoZSBQYWNrYWdlLgoKCkl0ZW1zIFRoYXQgYXJlIE5vdCBDb25zaWRlcmVkIFBhcnQgb2YgYSBNb2RpZmllZCBWZXJzaW9uCgooOSkgV29ya3MgKGluY2x1ZGluZywgYnV0IG5vdCBsaW1pdGVkIHRvLCBtb2R1bGVzIGFuZCBzY3JpcHRzKSB0aGF0Cm1lcmVseSBleHRlbmQgb3IgbWFrZSB1c2Ugb2YgdGhlIFBhY2thZ2UsIGRvIG5vdCwgYnkgdGhlbXNlbHZlcywgY2F1c2UKdGhlIFBhY2thZ2UgdG8gYmUgYSBNb2RpZmllZCBWZXJzaW9uLiAgSW4gYWRkaXRpb24sIHN1Y2ggd29ya3MgYXJlIG5vdApjb25zaWRlcmVkIHBhcnRzIG9mIHRoZSBQYWNrYWdlIGl0c2VsZiwgYW5kIGFyZSBub3Qgc3ViamVjdCB0byB0aGUKdGVybXMgb2YgdGhpcyBsaWNlbnNlLgoKCkdlbmVyYWwgUHJvdmlzaW9ucwoKKDEwKSAgQW55IHVzZSwgbW9kaWZpY2F0aW9uLCBhbmQgZGlzdHJpYnV0aW9uIG9mIHRoZSBTdGFuZGFyZCBvcgpNb2RpZmllZCBWZXJzaW9ucyBpcyBnb3Zlcm5lZCBieSB0aGlzIEFydGlzdGljIExpY2Vuc2UuIEJ5IHVzaW5nLAptb2RpZnlpbmcgb3IgZGlzdHJpYnV0aW5nIHRoZSBQYWNrYWdlLCB5b3UgYWNjZXB0IHRoaXMgbGljZW5zZS4gRG8gbm90CnVzZSwgbW9kaWZ5LCBvciBkaXN0cmlidXRlIHRoZSBQYWNrYWdlLCBpZiB5b3UgZG8gbm90IGFjY2VwdCB0aGlzCmxpY2Vuc2UuCgooMTEpICBJZiB5b3VyIE1vZGlmaWVkIFZlcnNpb24gaGFzIGJlZW4gZGVyaXZlZCBmcm9tIGEgTW9kaWZpZWQKVmVyc2lvbiBtYWRlIGJ5IHNvbWVvbmUgb3RoZXIgdGhhbiB5b3UsIHlvdSBhcmUgbmV2ZXJ0aGVsZXNzIHJlcXVpcmVkCnRvIGVuc3VyZSB0aGF0IHlvdXIgTW9kaWZpZWQgVmVyc2lvbiBjb21wbGllcyB3aXRoIHRoZSByZXF1aXJlbWVudHMgb2YKdGhpcyBsaWNlbnNlLgoKKDEyKSAgVGhpcyBsaWNlbnNlIGRvZXMgbm90IGdyYW50IHlvdSB0aGUgcmlnaHQgdG8gdXNlIGFueSB0cmFkZW1hcmssCnNlcnZpY2UgbWFyaywgdHJhZGVuYW1lLCBvciBsb2dvIG9mIHRoZSBDb3B5cmlnaHQgSG9sZGVyLgoKKDEzKSAgVGhpcyBsaWNlbnNlIGluY2x1ZGVzIHRoZSBub24tZXhjbHVzaXZlLCB3b3JsZHdpZGUsCmZyZWUtb2YtY2hhcmdlIHBhdGVudCBsaWNlbnNlIHRvIG1ha2UsIGhhdmUgbWFkZSwgdXNlLCBvZmZlciB0byBzZWxsLApzZWxsLCBpbXBvcnQgYW5kIG90aGVyd2lzZSB0cmFuc2ZlciB0aGUgUGFja2FnZSB3aXRoIHJlc3BlY3QgdG8gYW55CnBhdGVudCBjbGFpbXMgbGljZW5zYWJsZSBieSB0aGUgQ29weXJpZ2h0IEhvbGRlciB0aGF0IGFyZSBuZWNlc3NhcmlseQppbmZyaW5nZWQgYnkgdGhlIFBhY2thZ2UuIElmIHlvdSBpbnN0aXR1dGUgcGF0ZW50IGxpdGlnYXRpb24KKGluY2x1ZGluZyBhIGNyb3NzLWNsYWltIG9yIGNvdW50ZXJjbGFpbSkgYWdhaW5zdCBhbnkgcGFydHkgYWxsZWdpbmcKdGhhdCB0aGUgUGFja2FnZSBjb25zdGl0dXRlcyBkaXJlY3Qgb3IgY29udHJpYnV0b3J5IHBhdGVudAppbmZyaW5nZW1lbnQsIHRoZW4gdGhpcyBBcnRpc3RpYyBMaWNlbnNlIHRvIHlvdSBzaGFsbCB0ZXJtaW5hdGUgb24gdGhlCmRhdGUgdGhhdCBzdWNoIGxpdGlnYXRpb24gaXMgZmlsZWQuCgooMTQpICBEaXNjbGFpbWVyIG9mIFdhcnJhbnR5OgpUSEUgUEFDS0FHRSBJUyBQUk9WSURFRCBCWSBUSEUgQ09QWVJJR0hUIEhPTERFUiBBTkQgQ09OVFJJQlVUT1JTICJBUwpJUycgQU5EIFdJVEhPVVQgQU5ZIEVYUFJFU1MgT1IgSU1QTElFRCBXQVJSQU5USUVTLiBUSEUgSU1QTElFRApXQVJSQU5USUVTIE9GIE1FUkNIQU5UQUJJTElUWSwgRklUTkVTUyBGT1IgQSBQQVJUSUNVTEFSIFBVUlBPU0UsIE9SCk5PTi1JTkZSSU5HRU1FTlQgQVJFIERJU0NMQUlNRUQgVE8gVEhFIEVYVEVOVCBQRVJNSVRURUQgQlkgWU9VUiBMT0NBTApMQVcuIFVOTEVTUyBSRVFVSVJFRCBCWSBMQVcsIE5PIENPUFlSSUdIVCBIT0xERVIgT1IgQ09OVFJJQlVUT1IgV0lMTApCRSBMSUFCTEUgRk9SIEFOWSBESVJFQ1QsIElORElSRUNULCBJTkNJREVOVEFMLCBPUiBDT05TRVFVRU5USUFMCkRBTUFHRVMgQVJJU0lORyBJTiBBTlkgV0FZIE9VVCBPRiBUSEUgVVNFIE9GIFRIRSBQQUNLQUdFLCBFVkVOIElGCkFEVklTRUQgT0YgVEhFIFBPU1NJQklMSVRZIE9GIFNVQ0ggREFNQUdFLgoKCi0tLS0tLS0tCg=="
    },
    {
      "path": "supplemental/BlueOak-1.0.0.txt",
      "source_locator": "Temp/elite-v338-606f40d5d434468ba8ceb8d424340173/BlueOak-1.0.0.txt",
      "scope": "official text, pinned vendor header, or exact selected-bundle attribution; not blanket clearance",
      "origin": "RETAINED_RESEARCH_EVIDENCE",
      "payload_path": null,
      "bytes": 1552,
      "sha256": "8a1af140fdfbf5afd3df27f7e662f989c5b963a300020dfafce42033cae9e004",
      "content_base64": "IyBCbHVlIE9hayBNb2RlbCBMaWNlbnNlCgpWZXJzaW9uIDEuMC4wCgojIyBQdXJwb3NlCgpUaGlzIGxpY2Vuc2UgZ2l2ZXMgZXZlcnlvbmUgYXMgbXVjaCBwZXJtaXNzaW9uIHRvIHdvcmsgd2l0aAp0aGlzIHNvZnR3YXJlIGFzIHBvc3NpYmxlLCB3aGlsZSBwcm90ZWN0aW5nIGNvbnRyaWJ1dG9ycwpmcm9tIGxpYWJpbGl0eS4KCiMjIEFjY2VwdGFuY2UKCkluIG9yZGVyIHRvIHJlY2VpdmUgdGhpcyBsaWNlbnNlLCB5b3UgbXVzdCBhZ3JlZSB0byBpdHMKcnVsZXMuICBUaGUgcnVsZXMgb2YgdGhpcyBsaWNlbnNlIGFyZSBib3RoIG9ibGlnYXRpb25zCnVuZGVyIHRoYXQgYWdyZWVtZW50IGFuZCBjb25kaXRpb25zIHRvIHlvdXIgbGljZW5zZS4KWW91IG11c3Qgbm90IGRvIGFueXRoaW5nIHdpdGggdGhpcyBzb2Z0d2FyZSB0aGF0IHRyaWdnZXJzCmEgcnVsZSB0aGF0IHlvdSBjYW5ub3Qgb3Igd2lsbCBub3QgZm9sbG93LgoKIyMgQ29weXJpZ2h0CgpFYWNoIGNvbnRyaWJ1dG9yIGxpY2Vuc2VzIHlvdSB0byBkbyBldmVyeXRoaW5nIHdpdGggdGhpcwpzb2Z0d2FyZSB0aGF0IHdvdWxkIG90aGVyd2lzZSBpbmZyaW5nZSB0aGF0IGNvbnRyaWJ1dG9yJ3MKY29weXJpZ2h0IGluIGl0LgoKIyMgTm90aWNlcwoKWW91IG11c3QgZW5zdXJlIHRoYXQgZXZlcnlvbmUgd2hvIGdldHMgYSBjb3B5IG9mCmFueSBwYXJ0IG9mIHRoaXMgc29mdHdhcmUgZnJvbSB5b3UsIHdpdGggb3Igd2l0aG91dApjaGFuZ2VzLCBhbHNvIGdldHMgdGhlIHRleHQgb2YgdGhpcyBsaWNlbnNlIG9yIGEgbGluayB0bwo8aHR0cHM6Ly9ibHVlb2FrY291bmNpbC5vcmcvbGljZW5zZS8xLjAuMD4uCgojIyBFeGN1c2UKCklmIGFueW9uZSBub3RpZmllcyB5b3UgaW4gd3JpdGluZyB0aGF0IHlvdSBoYXZlIG5vdApjb21wbGllZCB3aXRoIFtOb3RpY2VzXSgjbm90aWNlcyksIHlvdSBjYW4ga2VlcCB5b3VyCmxpY2Vuc2UgYnkgdGFraW5nIGFsbCBwcmFjdGljYWwgc3RlcHMgdG8gY29tcGx5IHdpdGhpbiAzMApkYXlzIGFmdGVyIHRoZSBub3RpY2UuICBJZiB5b3UgZG8gbm90IGRvIHNvLCB5b3VyIGxpY2Vuc2UKZW5kcyBpbW1lZGlhdGVseS4KCiMjIFBhdGVudAoKRWFjaCBjb250cmlidXRvciBsaWNlbnNlcyB5b3UgdG8gZG8gZXZlcnl0aGluZyB3aXRoIHRoaXMKc29mdHdhcmUgdGhhdCB3b3VsZCBvdGhlcndpc2UgaW5mcmluZ2UgYW55IHBhdGVudCBjbGFpbXMKdGhleSBjYW4gbGljZW5zZSBvciBiZWNvbWUgYWJsZSB0byBsaWNlbnNlLgoKIyMgUmVsaWFiaWxpdHkKCk5vIGNvbnRyaWJ1dG9yIGNhbiByZXZva2UgdGhpcyBsaWNlbnNlLgoKIyMgTm8gTGlhYmlsaXR5CgoqKipBcyBmYXIgYXMgdGhlIGxhdyBhbGxvd3MsIHRoaXMgc29mdHdhcmUgY29tZXMgYXMgaXMsCndpdGhvdXQgYW55IHdhcnJhbnR5IG9yIGNvbmRpdGlvbiwgYW5kIG5vIGNvbnRyaWJ1dG9yCndpbGwgYmUgbGlhYmxlIHRvIGFueW9uZSBmb3IgYW55IGRhbWFnZXMgcmVsYXRlZCB0byB0aGlzCnNvZnR3YXJlIG9yIHRoaXMgbGljZW5zZSwgdW5kZXIgYW55IGtpbmQgb2YgbGVnYWwgY2xhaW0uKioqCg=="
    },
    {
      "path": "supplemental/qrcode-vendor-header.txt",
      "source_locator": "Temp/elite-v338-606f40d5d434468ba8ceb8d424340173/qrcode-vendor-header.txt",
      "scope": "official text, pinned vendor header, or exact selected-bundle attribution; not blanket clearance",
      "origin": "RETAINED_RESEARCH_EVIDENCE",
      "payload_path": null,
      "bytes": 618,
      "sha256": "f265b9225bb2a1a209d60f81be487f23257c3f8637498282212d04af7822c2b3",
      "content_base64": "Ly8tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0KLy8gUVJDb2RlIGZvciBKYXZhU2NyaXB0Ci8vCi8vIENvcHlyaWdodCAoYykgMjAwOSBLYXp1aGlrbyBBcmFzZQovLwovLyBVUkw6IGh0dHA6Ly93d3cuZC1wcm9qZWN0LmNvbS8KLy8KLy8gTGljZW5zZWQgdW5kZXIgdGhlIE1JVCBsaWNlbnNlOgovLyAgIGh0dHA6Ly93d3cub3BlbnNvdXJjZS5vcmcvbGljZW5zZXMvbWl0LWxpY2Vuc2UucGhwCi8vCi8vIFRoZSB3b3JkICJRUiBDb2RlIiBpcyByZWdpc3RlcmVkIHRyYWRlbWFyayBvZiAKLy8gREVOU08gV0FWRSBJTkNPUlBPUkFURUQKLy8gICBodHRwOi8vd3d3LmRlbnNvLXdhdmUuY29tL3FyY29kZS9mYXFwYXRlbnQtZS5odG1sCi8vCi8vLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tCi8vIE1vZGlmaWVkIHRvIHdvcmsgaW4gbm9kZSBmb3IgdGhpcyBwcm9qZWN0IChhbmQgc29tZSByZWZhY3RvcmluZykKLy8tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0K"
    },
    {
      "path": "supplemental/bundled-attributions.txt",
      "source_locator": "Temp/elite-v338-606f40d5d434468ba8ceb8d424340173/bundled-attributions.txt",
      "scope": "official text, pinned vendor header, or exact selected-bundle attribution; not blanket clearance",
      "origin": "RETAINED_RESEARCH_EVIDENCE",
      "payload_path": null,
      "bytes": 7062,
      "sha256": "d958cbc03540bf9fb23ef6e096d94ed2567b4bd0c378c96b997ae63e1ef1ff44",
      "content_base64": "LyohIEJ1bmRsZWQgbGljZW5zZSBpbmZvcm1hdGlvbjoKCmlzLXdpbmRvd3MvaW5kZXguanM6CiAgKCohCiAgICogaXMtd2luZG93cyA8aHR0cHM6Ly9naXRodWIuY29tL2pvbnNjaGxpbmtlcnQvaXMtd2luZG93cz4KICAgKgogICAqIENvcHlyaWdodCDCqSAyMDE1LTIwMTgsIEpvbiBTY2hsaW5rZXJ0LgogICAqIFJlbGVhc2VkIHVuZGVyIHRoZSBNSVQgTGljZW5zZS4KICAgKikKCmlzLWV4dGdsb2IvaW5kZXguanM6CiAgKCohCiAgICogaXMtZXh0Z2xvYiA8aHR0cHM6Ly9naXRodWIuY29tL2pvbnNjaGxpbmtlcnQvaXMtZXh0Z2xvYj4KICAgKgogICAqIENvcHlyaWdodCAoYykgMjAxNC0yMDE2LCBKb24gU2NobGlua2VydC4KICAgKiBMaWNlbnNlZCB1bmRlciB0aGUgTUlUIExpY2Vuc2UuCiAgICopCgppcy1nbG9iL2luZGV4LmpzOgogICgqIQogICAqIGlzLWdsb2IgPGh0dHBzOi8vZ2l0aHViLmNvbS9qb25zY2hsaW5rZXJ0L2lzLWdsb2I+CiAgICoKICAgKiBDb3B5cmlnaHQgKGMpIDIwMTQtMjAxNywgSm9uIFNjaGxpbmtlcnQuCiAgICogUmVsZWFzZWQgdW5kZXIgdGhlIE1JVCBMaWNlbnNlLgogICAqKQoKaXMtbnVtYmVyL2luZGV4LmpzOgogICgqIQogICAqIGlzLW51bWJlciA8aHR0cHM6Ly9naXRodWIuY29tL2pvbnNjaGxpbmtlcnQvaXMtbnVtYmVyPgogICAqCiAgICogQ29weXJpZ2h0IChjKSAyMDE0LXByZXNlbnQsIEpvbiBTY2hsaW5rZXJ0LgogICAqIFJlbGVhc2VkIHVuZGVyIHRoZSBNSVQgTGljZW5zZS4KICAgKikKCnRvLXJlZ2V4LXJhbmdlL2luZGV4LmpzOgogICgqIQogICAqIHRvLXJlZ2V4LXJhbmdlIDxodHRwczovL2dpdGh1Yi5jb20vbWljcm9tYXRjaC90by1yZWdleC1yYW5nZT4KICAgKgogICAqIENvcHlyaWdodCAoYykgMjAxNS1wcmVzZW50LCBKb24gU2NobGlua2VydC4KICAgKiBSZWxlYXNlZCB1bmRlciB0aGUgTUlUIExpY2Vuc2UuCiAgICopCgpmaWxsLXJhbmdlL2luZGV4LmpzOgogICgqIQogICAqIGZpbGwtcmFuZ2UgPGh0dHBzOi8vZ2l0aHViLmNvbS9qb25zY2hsaW5rZXJ0L2ZpbGwtcmFuZ2U+CiAgICoKICAgKiBDb3B5cmlnaHQgKGMpIDIwMTQtcHJlc2VudCwgSm9uIFNjaGxpbmtlcnQuCiAgICogTGljZW5zZWQgdW5kZXIgdGhlIE1JVCBMaWNlbnNlLgogICAqKQoKcXVldWUtbWljcm90YXNrL2luZGV4LmpzOgogICgqISBxdWV1ZS1taWNyb3Rhc2suIE1JVCBMaWNlbnNlLiBGZXJvc3MgQWJvdWtoYWRpamVoIDxodHRwczovL2Zlcm9zcy5vcmcvb3BlbnNvdXJjZT4gKikKCnJ1bi1wYXJhbGxlbC9pbmRleC5qczoKICAoKiEgcnVuLXBhcmFsbGVsLiBNSVQgTGljZW5zZS4gRmVyb3NzIEFib3VraGFkaWplaCA8aHR0cHM6Ly9mZXJvc3Mub3JnL29wZW5zb3VyY2U+ICopCgpAemtvY2hhbi9qcy15YW1sL2Rpc3QvanMteWFtbC5tanM6CiAgKCohIEB6a29jaGFuL2pzLXlhbWwgMC4wLjExIGh0dHBzOi8vZ2l0aHViLmNvbS9ub2RlY2EvanMteWFtbCBAbGljZW5zZSBNSVQgKikKCm5vcm1hbGl6ZS1wYXRoL2luZGV4LmpzOgogICgqIQogICAqIG5vcm1hbGl6ZS1wYXRoIDxodHRwczovL2dpdGh1Yi5jb20vam9uc2NobGlua2VydC9ub3JtYWxpemUtcGF0aD4KICAgKgogICAqIENvcHlyaWdodCAoYykgMjAxNC0yMDE4LCBKb24gU2NobGlua2VydC4KICAgKiBSZWxlYXNlZCB1bmRlciB0aGUgTUlUIExpY2Vuc2UuCiAgICopCgp1bmRpY2kvbGliL3dlYi9mZXRjaC9ib2R5LmpzOgogICgqISBmb3JtZGF0YS1wb2x5ZmlsbC4gTUlUIExpY2Vuc2UuIEppbW15IFfDpHJ0aW5nIDxodHRwczovL2ppbW15LndhcnRpbmcuc2Uvb3BlbnNvdXJjZT4gKikKCnVuZGljaS9saWIvd2ViL3dlYnNvY2tldC9mcmFtZS5qczoKICAoKiEgd3MuIE1JVCBMaWNlbnNlLiBFaW5hciBPdHRvIFN0YW5ndmlrIDxlaW5hcm9zQGdtYWlsLmNvbT4gKikKCm9wZW5wZ3AvZGlzdC9ub2RlL29wZW5wZ3AubWpzOgogICgqISBPcGVuUEdQLmpzIHY2LjMuMSAtIDIwMjYtMDYtMDQgLSB0aGlzIGlzIExHUEwgbGljZW5zZWQgY29kZSwgc2VlIExJQ0VOU0Uvb3VyIHdlYnNpdGUgaHR0cHM6Ly9vcGVucGdwanMub3JnLyBmb3IgbW9yZSBpbmZvcm1hdGlvbi4gKikKICAoKiEgbm9ibGUtY2lwaGVycyAtIE1JVCBMaWNlbnNlIChjKSAyMDIzIFBhdWwgTWlsbGVyIChwYXVsbWlsbHIuY29tKSAqKQogICgqISBub2JsZS1oYXNoZXMgLSBNSVQgTGljZW5zZSAoYykgMjAyMiBQYXVsIE1pbGxlciAocGF1bG1pbGxyLmNvbSkgKikKICAoKiEgbm9ibGUtY3VydmVzIC0gTUlUIExpY2Vuc2UgKGMpIDIwMjIgUGF1bCBNaWxsZXIgKHBhdWxtaWxsci5jb20pICopCgppbXVybXVyaGFzaC9pbXVybXVyaGFzaC5qczoKICAoKioKICAgKiBAcHJlc2VydmUKICAgKiBKUyBJbXBsZW1lbnRhdGlvbiBvZiBpbmNyZW1lbnRhbCBNdXJtdXJIYXNoMyAocjE1MCkgKGFzIG9mIE1heSAxMCwgMjAxMykKICAgKgogICAqIEBhdXRob3IgPGEgaHJlZj0ibWFpbHRvOmplbnN5dEBnbWFpbC5jb20iPkplbnMgVGF5bG9yPC9hPgogICAqIEBzZWUgaHR0cDovL2dpdGh1Yi5jb20vaG9tZWJyZXdpbmcvYnJhdWhhdXMtZGlmZgogICAqIEBhdXRob3IgPGEgaHJlZj0ibWFpbHRvOmdhcnkuY291cnRAZ21haWwuY29tIj5HYXJ5IENvdXJ0PC9hPgogICAqIEBzZWUgaHR0cDovL2dpdGh1Yi5jb20vZ2FyeWNvdXJ0L211cm11cmhhc2gtanMKICAgKiBAYXV0aG9yIDxhIGhyZWY9Im1haWx0bzphYXBwbGVieUBnbWFpbC5jb20iPkF1c3RpbiBBcHBsZWJ5PC9hPgogICAqIEBzZWUgaHR0cDovL3NpdGVzLmdvb2dsZS5jb20vc2l0ZS9tdXJtdXJoYXNoLwogICAqKQoKQHlhcm5wa2cvcG5wL2xpYi9pbmRleC5qczoKICAoKioKICAgIEBsaWNlbnNlCiAgICBDb3B5cmlnaHQgTm9kZS5qcyBjb250cmlidXRvcnMuIEFsbCByaWdodHMgcmVzZXJ2ZWQuCiAgCiAgICBQZXJtaXNzaW9uIGlzIGhlcmVieSBncmFudGVkLCBmcmVlIG9mIGNoYXJnZSwgdG8gYW55IHBlcnNvbiBvYnRhaW5pbmcgYSBjb3B5CiAgICBvZiB0aGlzIHNvZnR3YXJlIGFuZCBhc3NvY2lhdGVkIGRvY3VtZW50YXRpb24gZmlsZXMgKHRoZSAiU29mdHdhcmUiKSwgdG8KICAgIGRlYWwgaW4gdGhlIFNvZnR3YXJlIHdpdGhvdXQgcmVzdHJpY3Rpb24sIGluY2x1ZGluZyB3aXRob3V0IGxpbWl0YXRpb24gdGhlCiAgICByaWdodHMgdG8gdXNlLCBjb3B5LCBtb2RpZnksIG1lcmdlLCBwdWJsaXNoLCBkaXN0cmlidXRlLCBzdWJsaWNlbnNlLCBhbmQvb3IKICAgIHNlbGwgY29waWVzIG9mIHRoZSBTb2Z0d2FyZSwgYW5kIHRvIHBlcm1pdCBwZXJzb25zIHRvIHdob20gdGhlIFNvZnR3YXJlIGlzCiAgICBmdXJuaXNoZWQgdG8gZG8gc28sIHN1YmplY3QgdG8gdGhlIGZvbGxvd2luZyBjb25kaXRpb25zOgogIAogICAgVGhlIGFib3ZlIGNvcHlyaWdodCBub3RpY2UgYW5kIHRoaXMgcGVybWlzc2lvbiBub3RpY2Ugc2hhbGwgYmUgaW5jbHVkZWQgaW4KICAgIGFsbCBjb3BpZXMgb3Igc3Vic3RhbnRpYWwgcG9ydGlvbnMgb2YgdGhlIFNvZnR3YXJlLgogIAogICAgVEhFIFNPRlRXQVJFIElTIFBST1ZJREVEICJBUyBJUyIsIFdJVEhPVVQgV0FSUkFOVFkgT0YgQU5ZIEtJTkQsIEVYUFJFU1MgT1IKICAgIElNUExJRUQsIElOQ0xVRElORyBCVVQgTk9UIExJTUlURUQgVE8gVEhFIFdBUlJBTlRJRVMgT0YgTUVSQ0hBTlRBQklMSVRZLAogICAgRklUTkVTUyBGT1IgQSBQQVJUSUNVTEFSIFBVUlBPU0UgQU5EIE5PTklORlJJTkdFTUVOVC4gSU4gTk8gRVZFTlQgU0hBTEwgVEhFCiAgICBBVVRIT1JTIE9SIENPUFlSSUdIVCBIT0xERVJTIEJFIExJQUJMRSBGT1IgQU5ZIENMQUlNLCBEQU1BR0VTIE9SIE9USEVSCiAgICBMSUFCSUxJVFksIFdIRVRIRVIgSU4gQU4gQUNUSU9OIE9GIENPTlRSQUNULCBUT1JUIE9SIE9USEVSV0lTRSwgQVJJU0lORwogICAgRlJPTSwgT1VUIE9GIE9SIElOIENPTk5FQ1RJT04gV0lUSCBUSEUgU09GVFdBUkUgT1IgVEhFIFVTRSBPUiBPVEhFUiBERUFMSU5HUwogICAgSU4gVEhFIFNPRlRXQVJFLgogICopCiAgKCoqCiAgICBAbGljZW5zZQogICAgVGhlIE1JVCBMaWNlbnNlIChNSVQpCiAgCiAgICBDb3B5cmlnaHQgKGMpIDIwMTQgQmxha2UgRW1icmV5IChoZWxsb0BibGFrZWVtYnJleS5jb20pCiAgCiAgICBQZXJtaXNzaW9uIGlzIGhlcmVieSBncmFudGVkLCBmcmVlIG9mIGNoYXJnZSwgdG8gYW55IHBlcnNvbiBvYnRhaW5pbmcgYSBjb3B5CiAgICBvZiB0aGlzIHNvZnR3YXJlIGFuZCBhc3NvY2lhdGVkIGRvY3VtZW50YXRpb24gZmlsZXMgKHRoZSAiU29mdHdhcmUiKSwgdG8gZGVhbAogICAgaW4gdGhlIFNvZnR3YXJlIHdpdGhvdXQgcmVzdHJpY3Rpb24sIGluY2x1ZGluZyB3aXRob3V0IGxpbWl0YXRpb24gdGhlIHJpZ2h0cwogICAgdG8gdXNlLCBjb3B5LCBtb2RpZnksIG1lcmdlLCBwdWJsaXNoLCBkaXN0cmlidXRlLCBzdWJsaWNlbnNlLCBhbmQvb3Igc2VsbAogICAgY29waWVzIG9mIHRoZSBTb2Z0d2FyZSwgYW5kIHRvIHBlcm1pdCBwZXJzb25zIHRvIHdob20gdGhlIFNvZnR3YXJlIGlzCiAgICBmdXJuaXNoZWQgdG8gZG8gc28sIHN1YmplY3QgdG8gdGhlIGZvbGxvd2luZyBjb25kaXRpb25zOgogIAogICAgVGhlIGFib3ZlIGNvcHlyaWdodCBub3RpY2UgYW5kIHRoaXMgcGVybWlzc2lvbiBub3RpY2Ugc2hhbGwgYmUgaW5jbHVkZWQgaW4KICAgIGFsbCBjb3BpZXMgb3Igc3Vic3RhbnRpYWwgcG9ydGlvbnMgb2YgdGhlIFNvZnR3YXJlLgogIAogICAgVEhFIFNPRlRXQVJFIElTIFBST1ZJREVEICJBUyBJUyIsIFdJVEhPVVQgV0FSUkFOVFkgT0YgQU5ZIEtJTkQsIEVYUFJFU1MgT1IKICAgIElNUExJRUQsIElOQ0xVRElORyBCVVQgTk9UIExJTUlURUQgVE8gVEhFIFdBUlJBTlRJRVMgT0YgTUVSQ0hBTlRBQklMSVRZLAogICAgRklUTkVTUyBGT1IgQSBQQVJUSUNVTEFSIFBVUlBPU0UgQU5EIE5PTklORlJJTkdFTUVOVC4gSU4gTk8gRVZFTlQgU0hBTEwgVEhFCiAgICBBVVRIT1JTIE9SIENPUFlSSUdIVCBIT0xERVJTIEJFIExJQUJMRSBGT1IgQU5ZIENMQUlNLCBEQU1BR0VTIE9SIE9USEVSCiAgICBMSUFCSUxJVFksIFdIRVRIRVIgSU4gQU4gQUNUSU9OIE9GIENPTlRSQUNULCBUT1JUIE9SIE9USEVSV0lTRSwgQVJJU0lORyBGUk9NLAogICAgT1VUIE9GIE9SIElOIENPTk5FQ1RJT04gV0lUSCBUSEUgU09GVFdBUkUgT1IgVEhFIFVTRSBPUiBPVEhFUiBERUFMSU5HUyBJTgogICAgVEhFIFNPRlRXQVJFLgogICopCgpAY3ljbG9uZWR4L2N5Y2xvbmVkeC1saWJyYXJ5L2Rpc3Qubm9kZS9zcGR4LmpzOgogICgqIQogIFRoaXMgZmlsZSBpcyBwYXJ0IG9mIEN5Y2xvbmVEWCBKYXZhU2NyaXB0IExpYnJhcnkuCiAgCiAgTGljZW5zZWQgdW5kZXIgdGhlIEFwYWNoZSBMaWNlbnNlLCBWZXJzaW9uIDIuMCAodGhlICJMaWNlbnNlIik7CiAgeW91IG1heSBub3QgdXNlIHRoaXMgZmlsZSBleGNlcHQgaW4gY29tcGxpYW5jZSB3aXRoIHRoZSBMaWNlbnNlLgogIFlvdSBtYXkgb2J0YWluIGEgY29weSBvZiB0aGUgTGljZW5zZSBhdAogIAogICAgIGh0dHA6Ly93d3cuYXBhY2hlLm9yZy9saWNlbnNlcy9MSUNFTlNFLTIuMAogIAogIFVubGVzcyByZXF1aXJlZCBieSBhcHBsaWNhYmxlIGxhdyBvciBhZ3JlZWQgdG8gaW4gd3JpdGluZywgc29mdHdhcmUKICBkaXN0cmlidXRlZCB1bmRlciB0aGUgTGljZW5zZSBpcyBkaXN0cmlidXRlZCBvbiBhbiAiQVMgSVMiIEJBU0lTLAogIFdJVEhPVVQgV0FSUkFOVElFUyBPUiBDT05ESVRJT05TIE9GIEFOWSBLSU5ELCBlaXRoZXIgZXhwcmVzcyBvciBpbXBsaWVkLgogIFNlZSB0aGUgTGljZW5zZSBmb3IgdGhlIHNwZWNpZmljIGxhbmd1YWdlIGdvdmVybmluZyBwZXJtaXNzaW9ucyBhbmQKICBsaW1pdGF0aW9ucyB1bmRlciB0aGUgTGljZW5zZS4KICAKICBTUERYLUxpY2Vuc2UtSWRlbnRpZmllcjogQXBhY2hlLTIuMAogIENvcHlyaWdodCAoYykgT1dBU1AgRm91bmRhdGlvbi4gQWxsIFJpZ2h0cyBSZXNlcnZlZC4KICAqKQoKY29udGVudC10eXBlL2Rpc3QvaW5kZXguanM6CiAgKCohCiAgICogY29udGVudC10eXBlCiAgICogQ29weXJpZ2h0KGMpIDIwMTUgRG91Z2xhcyBDaHJpc3RvcGhlciBXaWxzb24KICAgKiBNSVQgTGljZW5zZWQKICAgKikKCm5lZ290aWF0b3IvbGliL2FjY2VwdC5qczoKICAoKiEKICAgKiBuZWdvdGlhdG9yCiAgICogQ29weXJpZ2h0KGMpIDIwMjYgQmxha2UgRW1icmV5CiAgICogTUlUIExpY2Vuc2VkCiAgICopCgpuZWdvdGlhdG9yL2luZGV4LmpzOgogICgqIQogICAqIG5lZ290aWF0b3IKICAgKiBDb3B5cmlnaHQoYykgMjAxMiBGZWRlcmljbyBSb21lcm8KICAgKiBDb3B5cmlnaHQoYykgMjAxMi0yMDE0IElzYWFjIFouIFNjaGx1ZXRlcgogICAqIENvcHlyaWdodChjKSAyMDE1IERvdWdsYXMgQ2hyaXN0b3BoZXIgV2lsc29uCiAgICogTUlUIExpY2Vuc2VkCiAgICopCgptYWtlLWZldGNoLWhhcHBlbi9saWIvZmV0Y2guanM6CiAgKCoqCiAgICogQGxpY2Vuc2UKICAgKiBDb3B5cmlnaHQgKGMpIDIwMTAtMjAxMiBNaWtlYWwgUm9nZXJzCiAgICogTGljZW5zZWQgdW5kZXIgdGhlIEFwYWNoZSBMaWNlbnNlLCBWZXJzaW9uIDIuMCAodGhlICJMaWNlbnNlIik7CiAgICogeW91IG1heSBub3QgdXNlIHRoaXMgZmlsZSBleGNlcHQgaW4gY29tcGxpYW5jZSB3aXRoIHRoZSBMaWNlbnNlLgogICAqIFlvdSBtYXkgb2J0YWluIGEgY29weSBvZiB0aGUgTGljZW5zZSBhdAogICAqIGh0dHA6Ly93d3cuYXBhY2hlLm9yZy9saWNlbnNlcy9MSUNFTlNFLTIuMAogICAqIFVubGVzcyByZXF1aXJlZCBieSBhcHBsaWNhYmxlIGxhdyBvciBhZ3JlZWQgdG8gaW4gd3JpdGluZywKICAgKiBzb2Z0d2FyZSBkaXN0cmlidXRlZCB1bmRlciB0aGUgTGljZW5zZSBpcyBkaXN0cmlidXRlZCBvbiBhbiAiQVMKICAgKiBJUyIgQkFTSVMsIFdJVEhPVVQgV0FSUkFOVElFUyBPUiBDT05ESVRJT05TIE9GIEFOWSBLSU5ELCBlaXRoZXIKICAgKiBleHByZXNzIG9yIGltcGxpZWQuIFNlZSB0aGUgTGljZW5zZSBmb3IgdGhlIHNwZWNpZmljIGxhbmd1YWdlCiAgICogZ292ZXJuaW5nIHBlcm1pc3Npb25zIGFuZCBsaW1pdGF0aW9ucyB1bmRlciB0aGUgTGljZW5zZS4KICAgKikKKi8K"
    },
    {
      "path": "supplemental-v339/isexe-LICENSE.txt",
      "source_locator": "Temp/elite-v339-9266379aff0044498a1d7004267002c8/isexe-LICENSE.txt",
      "scope": "Pinned source licence; see V339 source-license-cases, no whole-package admission",
      "origin": "RETAINED_RESEARCH_EVIDENCE",
      "payload_path": null,
      "bytes": 1552,
      "sha256": "8a1af140fdfbf5afd3df27f7e662f989c5b963a300020dfafce42033cae9e004",
      "content_base64": "IyBCbHVlIE9hayBNb2RlbCBMaWNlbnNlCgpWZXJzaW9uIDEuMC4wCgojIyBQdXJwb3NlCgpUaGlzIGxpY2Vuc2UgZ2l2ZXMgZXZlcnlvbmUgYXMgbXVjaCBwZXJtaXNzaW9uIHRvIHdvcmsgd2l0aAp0aGlzIHNvZnR3YXJlIGFzIHBvc3NpYmxlLCB3aGlsZSBwcm90ZWN0aW5nIGNvbnRyaWJ1dG9ycwpmcm9tIGxpYWJpbGl0eS4KCiMjIEFjY2VwdGFuY2UKCkluIG9yZGVyIHRvIHJlY2VpdmUgdGhpcyBsaWNlbnNlLCB5b3UgbXVzdCBhZ3JlZSB0byBpdHMKcnVsZXMuICBUaGUgcnVsZXMgb2YgdGhpcyBsaWNlbnNlIGFyZSBib3RoIG9ibGlnYXRpb25zCnVuZGVyIHRoYXQgYWdyZWVtZW50IGFuZCBjb25kaXRpb25zIHRvIHlvdXIgbGljZW5zZS4KWW91IG11c3Qgbm90IGRvIGFueXRoaW5nIHdpdGggdGhpcyBzb2Z0d2FyZSB0aGF0IHRyaWdnZXJzCmEgcnVsZSB0aGF0IHlvdSBjYW5ub3Qgb3Igd2lsbCBub3QgZm9sbG93LgoKIyMgQ29weXJpZ2h0CgpFYWNoIGNvbnRyaWJ1dG9yIGxpY2Vuc2VzIHlvdSB0byBkbyBldmVyeXRoaW5nIHdpdGggdGhpcwpzb2Z0d2FyZSB0aGF0IHdvdWxkIG90aGVyd2lzZSBpbmZyaW5nZSB0aGF0IGNvbnRyaWJ1dG9yJ3MKY29weXJpZ2h0IGluIGl0LgoKIyMgTm90aWNlcwoKWW91IG11c3QgZW5zdXJlIHRoYXQgZXZlcnlvbmUgd2hvIGdldHMgYSBjb3B5IG9mCmFueSBwYXJ0IG9mIHRoaXMgc29mdHdhcmUgZnJvbSB5b3UsIHdpdGggb3Igd2l0aG91dApjaGFuZ2VzLCBhbHNvIGdldHMgdGhlIHRleHQgb2YgdGhpcyBsaWNlbnNlIG9yIGEgbGluayB0bwo8aHR0cHM6Ly9ibHVlb2FrY291bmNpbC5vcmcvbGljZW5zZS8xLjAuMD4uCgojIyBFeGN1c2UKCklmIGFueW9uZSBub3RpZmllcyB5b3UgaW4gd3JpdGluZyB0aGF0IHlvdSBoYXZlIG5vdApjb21wbGllZCB3aXRoIFtOb3RpY2VzXSgjbm90aWNlcyksIHlvdSBjYW4ga2VlcCB5b3VyCmxpY2Vuc2UgYnkgdGFraW5nIGFsbCBwcmFjdGljYWwgc3RlcHMgdG8gY29tcGx5IHdpdGhpbiAzMApkYXlzIGFmdGVyIHRoZSBub3RpY2UuICBJZiB5b3UgZG8gbm90IGRvIHNvLCB5b3VyIGxpY2Vuc2UKZW5kcyBpbW1lZGlhdGVseS4KCiMjIFBhdGVudAoKRWFjaCBjb250cmlidXRvciBsaWNlbnNlcyB5b3UgdG8gZG8gZXZlcnl0aGluZyB3aXRoIHRoaXMKc29mdHdhcmUgdGhhdCB3b3VsZCBvdGhlcndpc2UgaW5mcmluZ2UgYW55IHBhdGVudCBjbGFpbXMKdGhleSBjYW4gbGljZW5zZSBvciBiZWNvbWUgYWJsZSB0byBsaWNlbnNlLgoKIyMgUmVsaWFiaWxpdHkKCk5vIGNvbnRyaWJ1dG9yIGNhbiByZXZva2UgdGhpcyBsaWNlbnNlLgoKIyMgTm8gTGlhYmlsaXR5CgoqKipBcyBmYXIgYXMgdGhlIGxhdyBhbGxvd3MsIHRoaXMgc29mdHdhcmUgY29tZXMgYXMgaXMsCndpdGhvdXQgYW55IHdhcnJhbnR5IG9yIGNvbmRpdGlvbiwgYW5kIG5vIGNvbnRyaWJ1dG9yCndpbGwgYmUgbGlhYmxlIHRvIGFueW9uZSBmb3IgYW55IGRhbWFnZXMgcmVsYXRlZCB0byB0aGlzCnNvZnR3YXJlIG9yIHRoaXMgbGljZW5zZSwgdW5kZXIgYW55IGtpbmQgb2YgbGVnYWwgY2xhaW0uKioqCg=="
    },
    {
      "path": "supplemental-v339/minipass-LICENSE.txt",
      "source_locator": "Temp/elite-v339-9266379aff0044498a1d7004267002c8/minipass-LICENSE.txt",
      "scope": "Pinned source licence; see V339 source-license-cases, no whole-package admission",
      "origin": "RETAINED_RESEARCH_EVIDENCE",
      "payload_path": null,
      "bytes": 1552,
      "sha256": "8a1af140fdfbf5afd3df27f7e662f989c5b963a300020dfafce42033cae9e004",
      "content_base64": "IyBCbHVlIE9hayBNb2RlbCBMaWNlbnNlCgpWZXJzaW9uIDEuMC4wCgojIyBQdXJwb3NlCgpUaGlzIGxpY2Vuc2UgZ2l2ZXMgZXZlcnlvbmUgYXMgbXVjaCBwZXJtaXNzaW9uIHRvIHdvcmsgd2l0aAp0aGlzIHNvZnR3YXJlIGFzIHBvc3NpYmxlLCB3aGlsZSBwcm90ZWN0aW5nIGNvbnRyaWJ1dG9ycwpmcm9tIGxpYWJpbGl0eS4KCiMjIEFjY2VwdGFuY2UKCkluIG9yZGVyIHRvIHJlY2VpdmUgdGhpcyBsaWNlbnNlLCB5b3UgbXVzdCBhZ3JlZSB0byBpdHMKcnVsZXMuICBUaGUgcnVsZXMgb2YgdGhpcyBsaWNlbnNlIGFyZSBib3RoIG9ibGlnYXRpb25zCnVuZGVyIHRoYXQgYWdyZWVtZW50IGFuZCBjb25kaXRpb25zIHRvIHlvdXIgbGljZW5zZS4KWW91IG11c3Qgbm90IGRvIGFueXRoaW5nIHdpdGggdGhpcyBzb2Z0d2FyZSB0aGF0IHRyaWdnZXJzCmEgcnVsZSB0aGF0IHlvdSBjYW5ub3Qgb3Igd2lsbCBub3QgZm9sbG93LgoKIyMgQ29weXJpZ2h0CgpFYWNoIGNvbnRyaWJ1dG9yIGxpY2Vuc2VzIHlvdSB0byBkbyBldmVyeXRoaW5nIHdpdGggdGhpcwpzb2Z0d2FyZSB0aGF0IHdvdWxkIG90aGVyd2lzZSBpbmZyaW5nZSB0aGF0IGNvbnRyaWJ1dG9yJ3MKY29weXJpZ2h0IGluIGl0LgoKIyMgTm90aWNlcwoKWW91IG11c3QgZW5zdXJlIHRoYXQgZXZlcnlvbmUgd2hvIGdldHMgYSBjb3B5IG9mCmFueSBwYXJ0IG9mIHRoaXMgc29mdHdhcmUgZnJvbSB5b3UsIHdpdGggb3Igd2l0aG91dApjaGFuZ2VzLCBhbHNvIGdldHMgdGhlIHRleHQgb2YgdGhpcyBsaWNlbnNlIG9yIGEgbGluayB0bwo8aHR0cHM6Ly9ibHVlb2FrY291bmNpbC5vcmcvbGljZW5zZS8xLjAuMD4uCgojIyBFeGN1c2UKCklmIGFueW9uZSBub3RpZmllcyB5b3UgaW4gd3JpdGluZyB0aGF0IHlvdSBoYXZlIG5vdApjb21wbGllZCB3aXRoIFtOb3RpY2VzXSgjbm90aWNlcyksIHlvdSBjYW4ga2VlcCB5b3VyCmxpY2Vuc2UgYnkgdGFraW5nIGFsbCBwcmFjdGljYWwgc3RlcHMgdG8gY29tcGx5IHdpdGhpbiAzMApkYXlzIGFmdGVyIHRoZSBub3RpY2UuICBJZiB5b3UgZG8gbm90IGRvIHNvLCB5b3VyIGxpY2Vuc2UKZW5kcyBpbW1lZGlhdGVseS4KCiMjIFBhdGVudAoKRWFjaCBjb250cmlidXRvciBsaWNlbnNlcyB5b3UgdG8gZG8gZXZlcnl0aGluZyB3aXRoIHRoaXMKc29mdHdhcmUgdGhhdCB3b3VsZCBvdGhlcndpc2UgaW5mcmluZ2UgYW55IHBhdGVudCBjbGFpbXMKdGhleSBjYW4gbGljZW5zZSBvciBiZWNvbWUgYWJsZSB0byBsaWNlbnNlLgoKIyMgUmVsaWFiaWxpdHkKCk5vIGNvbnRyaWJ1dG9yIGNhbiByZXZva2UgdGhpcyBsaWNlbnNlLgoKIyMgTm8gTGlhYmlsaXR5CgoqKipBcyBmYXIgYXMgdGhlIGxhdyBhbGxvd3MsIHRoaXMgc29mdHdhcmUgY29tZXMgYXMgaXMsCndpdGhvdXQgYW55IHdhcnJhbnR5IG9yIGNvbmRpdGlvbiwgYW5kIG5vIGNvbnRyaWJ1dG9yCndpbGwgYmUgbGlhYmxlIHRvIGFueW9uZSBmb3IgYW55IGRhbWFnZXMgcmVsYXRlZCB0byB0aGlzCnNvZnR3YXJlIG9yIHRoaXMgbGljZW5zZSwgdW5kZXIgYW55IGtpbmQgb2YgbGVnYWwgY2xhaW0uKioqCg=="
    },
    {
      "path": "supplemental-v339/tar-LICENSE.txt",
      "source_locator": "Temp/elite-v339-9266379aff0044498a1d7004267002c8/tar-LICENSE.txt",
      "scope": "Pinned source licence; see V339 source-license-cases, no whole-package admission",
      "origin": "RETAINED_RESEARCH_EVIDENCE",
      "payload_path": null,
      "bytes": 1552,
      "sha256": "8a1af140fdfbf5afd3df27f7e662f989c5b963a300020dfafce42033cae9e004",
      "content_base64": "IyBCbHVlIE9hayBNb2RlbCBMaWNlbnNlCgpWZXJzaW9uIDEuMC4wCgojIyBQdXJwb3NlCgpUaGlzIGxpY2Vuc2UgZ2l2ZXMgZXZlcnlvbmUgYXMgbXVjaCBwZXJtaXNzaW9uIHRvIHdvcmsgd2l0aAp0aGlzIHNvZnR3YXJlIGFzIHBvc3NpYmxlLCB3aGlsZSBwcm90ZWN0aW5nIGNvbnRyaWJ1dG9ycwpmcm9tIGxpYWJpbGl0eS4KCiMjIEFjY2VwdGFuY2UKCkluIG9yZGVyIHRvIHJlY2VpdmUgdGhpcyBsaWNlbnNlLCB5b3UgbXVzdCBhZ3JlZSB0byBpdHMKcnVsZXMuICBUaGUgcnVsZXMgb2YgdGhpcyBsaWNlbnNlIGFyZSBib3RoIG9ibGlnYXRpb25zCnVuZGVyIHRoYXQgYWdyZWVtZW50IGFuZCBjb25kaXRpb25zIHRvIHlvdXIgbGljZW5zZS4KWW91IG11c3Qgbm90IGRvIGFueXRoaW5nIHdpdGggdGhpcyBzb2Z0d2FyZSB0aGF0IHRyaWdnZXJzCmEgcnVsZSB0aGF0IHlvdSBjYW5ub3Qgb3Igd2lsbCBub3QgZm9sbG93LgoKIyMgQ29weXJpZ2h0CgpFYWNoIGNvbnRyaWJ1dG9yIGxpY2Vuc2VzIHlvdSB0byBkbyBldmVyeXRoaW5nIHdpdGggdGhpcwpzb2Z0d2FyZSB0aGF0IHdvdWxkIG90aGVyd2lzZSBpbmZyaW5nZSB0aGF0IGNvbnRyaWJ1dG9yJ3MKY29weXJpZ2h0IGluIGl0LgoKIyMgTm90aWNlcwoKWW91IG11c3QgZW5zdXJlIHRoYXQgZXZlcnlvbmUgd2hvIGdldHMgYSBjb3B5IG9mCmFueSBwYXJ0IG9mIHRoaXMgc29mdHdhcmUgZnJvbSB5b3UsIHdpdGggb3Igd2l0aG91dApjaGFuZ2VzLCBhbHNvIGdldHMgdGhlIHRleHQgb2YgdGhpcyBsaWNlbnNlIG9yIGEgbGluayB0bwo8aHR0cHM6Ly9ibHVlb2FrY291bmNpbC5vcmcvbGljZW5zZS8xLjAuMD4uCgojIyBFeGN1c2UKCklmIGFueW9uZSBub3RpZmllcyB5b3UgaW4gd3JpdGluZyB0aGF0IHlvdSBoYXZlIG5vdApjb21wbGllZCB3aXRoIFtOb3RpY2VzXSgjbm90aWNlcyksIHlvdSBjYW4ga2VlcCB5b3VyCmxpY2Vuc2UgYnkgdGFraW5nIGFsbCBwcmFjdGljYWwgc3RlcHMgdG8gY29tcGx5IHdpdGhpbiAzMApkYXlzIGFmdGVyIHRoZSBub3RpY2UuICBJZiB5b3UgZG8gbm90IGRvIHNvLCB5b3VyIGxpY2Vuc2UKZW5kcyBpbW1lZGlhdGVseS4KCiMjIFBhdGVudAoKRWFjaCBjb250cmlidXRvciBsaWNlbnNlcyB5b3UgdG8gZG8gZXZlcnl0aGluZyB3aXRoIHRoaXMKc29mdHdhcmUgdGhhdCB3b3VsZCBvdGhlcndpc2UgaW5mcmluZ2UgYW55IHBhdGVudCBjbGFpbXMKdGhleSBjYW4gbGljZW5zZSBvciBiZWNvbWUgYWJsZSB0byBsaWNlbnNlLgoKIyMgUmVsaWFiaWxpdHkKCk5vIGNvbnRyaWJ1dG9yIGNhbiByZXZva2UgdGhpcyBsaWNlbnNlLgoKIyMgTm8gTGlhYmlsaXR5CgoqKipBcyBmYXIgYXMgdGhlIGxhdyBhbGxvd3MsIHRoaXMgc29mdHdhcmUgY29tZXMgYXMgaXMsCndpdGhvdXQgYW55IHdhcnJhbnR5IG9yIGNvbmRpdGlvbiwgYW5kIG5vIGNvbnRyaWJ1dG9yCndpbGwgYmUgbGlhYmxlIHRvIGFueW9uZSBmb3IgYW55IGRhbWFnZXMgcmVsYXRlZCB0byB0aGlzCnNvZnR3YXJlIG9yIHRoaXMgbGljZW5zZSwgdW5kZXIgYW55IGtpbmQgb2YgbGVnYWwgY2xhaW0uKioqCg=="
    },
    {
      "path": "supplemental-v339/yallist-LICENSE.txt",
      "source_locator": "Temp/elite-v339-9266379aff0044498a1d7004267002c8/yallist-LICENSE.txt",
      "scope": "Pinned source licence; see V339 source-license-cases, no whole-package admission",
      "origin": "RETAINED_RESEARCH_EVIDENCE",
      "payload_path": null,
      "bytes": 1764,
      "sha256": "a49c9ba464796f65b59fca3f1e6ca40912df1e859f575383223f7ec6c5baae09",
      "content_base64": "QWxsIHBhY2thZ2VzIHVuZGVyIGBzcmMvYCBhcmUgbGljZW5zZWQgYWNjb3JkaW5nIHRvIHRoZSB0ZXJtcyBpbgp0aGVpciByZXNwZWN0aXZlIGBMSUNFTlNFYCBvciBgTElDRU5TRS5tZGAgZmlsZXMuCgpUaGUgcmVtYWluZGVyIG9mIHRoaXMgcHJvamVjdCBpcyBsaWNlbnNlZCB1bmRlciB0aGUgQmx1ZSBPYWsKTW9kZWwgTGljZW5zZSwgYXMgZm9sbG93czoKCi0tLS0tCgojIEJsdWUgT2FrIE1vZGVsIExpY2Vuc2UKClZlcnNpb24gMS4wLjAKCiMjIFB1cnBvc2UKClRoaXMgbGljZW5zZSBnaXZlcyBldmVyeW9uZSBhcyBtdWNoIHBlcm1pc3Npb24gdG8gd29yayB3aXRoCnRoaXMgc29mdHdhcmUgYXMgcG9zc2libGUsIHdoaWxlIHByb3RlY3RpbmcgY29udHJpYnV0b3JzCmZyb20gbGlhYmlsaXR5LgoKIyMgQWNjZXB0YW5jZQoKSW4gb3JkZXIgdG8gcmVjZWl2ZSB0aGlzIGxpY2Vuc2UsIHlvdSBtdXN0IGFncmVlIHRvIGl0cwpydWxlcy4gIFRoZSBydWxlcyBvZiB0aGlzIGxpY2Vuc2UgYXJlIGJvdGggb2JsaWdhdGlvbnMKdW5kZXIgdGhhdCBhZ3JlZW1lbnQgYW5kIGNvbmRpdGlvbnMgdG8geW91ciBsaWNlbnNlLgpZb3UgbXVzdCBub3QgZG8gYW55dGhpbmcgd2l0aCB0aGlzIHNvZnR3YXJlIHRoYXQgdHJpZ2dlcnMKYSBydWxlIHRoYXQgeW91IGNhbm5vdCBvciB3aWxsIG5vdCBmb2xsb3cuCgojIyBDb3B5cmlnaHQKCkVhY2ggY29udHJpYnV0b3IgbGljZW5zZXMgeW91IHRvIGRvIGV2ZXJ5dGhpbmcgd2l0aCB0aGlzCnNvZnR3YXJlIHRoYXQgd291bGQgb3RoZXJ3aXNlIGluZnJpbmdlIHRoYXQgY29udHJpYnV0b3Incwpjb3B5cmlnaHQgaW4gaXQuCgojIyBOb3RpY2VzCgpZb3UgbXVzdCBlbnN1cmUgdGhhdCBldmVyeW9uZSB3aG8gZ2V0cyBhIGNvcHkgb2YKYW55IHBhcnQgb2YgdGhpcyBzb2Z0d2FyZSBmcm9tIHlvdSwgd2l0aCBvciB3aXRob3V0CmNoYW5nZXMsIGFsc28gZ2V0cyB0aGUgdGV4dCBvZiB0aGlzIGxpY2Vuc2Ugb3IgYSBsaW5rIHRvCjxodHRwczovL2JsdWVvYWtjb3VuY2lsLm9yZy9saWNlbnNlLzEuMC4wPi4KCiMjIEV4Y3VzZQoKSWYgYW55b25lIG5vdGlmaWVzIHlvdSBpbiB3cml0aW5nIHRoYXQgeW91IGhhdmUgbm90CmNvbXBsaWVkIHdpdGggW05vdGljZXNdKCNub3RpY2VzKSwgeW91IGNhbiBrZWVwIHlvdXIKbGljZW5zZSBieSB0YWtpbmcgYWxsIHByYWN0aWNhbCBzdGVwcyB0byBjb21wbHkgd2l0aGluIDMwCmRheXMgYWZ0ZXIgdGhlIG5vdGljZS4gIElmIHlvdSBkbyBub3QgZG8gc28sIHlvdXIgbGljZW5zZQplbmRzIGltbWVkaWF0ZWx5LgoKIyMgUGF0ZW50CgpFYWNoIGNvbnRyaWJ1dG9yIGxpY2Vuc2VzIHlvdSB0byBkbyBldmVyeXRoaW5nIHdpdGggdGhpcwpzb2Z0d2FyZSB0aGF0IHdvdWxkIG90aGVyd2lzZSBpbmZyaW5nZSBhbnkgcGF0ZW50IGNsYWltcwp0aGV5IGNhbiBsaWNlbnNlIG9yIGJlY29tZSBhYmxlIHRvIGxpY2Vuc2UuCgojIyBSZWxpYWJpbGl0eQoKTm8gY29udHJpYnV0b3IgY2FuIHJldm9rZSB0aGlzIGxpY2Vuc2UuCgojIyBObyBMaWFiaWxpdHkKCioqKkFzIGZhciBhcyB0aGUgbGF3IGFsbG93cywgdGhpcyBzb2Z0d2FyZSBjb21lcyBhcyBpcywKd2l0aG91dCBhbnkgd2FycmFudHkgb3IgY29uZGl0aW9uLCBhbmQgbm8gY29udHJpYnV0b3IKd2lsbCBiZSBsaWFibGUgdG8gYW55b25lIGZvciBhbnkgZGFtYWdlcyByZWxhdGVkIHRvIHRoaXMKc29mdHdhcmUgb3IgdGhpcyBsaWNlbnNlLCB1bmRlciBhbnkga2luZCBvZiBsZWdhbCBjbGFpbS4qKioK"
    },
    {
      "path": "supplemental-v339/noble-ciphers-LICENSE.txt",
      "source_locator": "Temp/elite-v339-9266379aff0044498a1d7004267002c8/noble-ciphers-LICENSE.txt",
      "scope": "Pinned source licence; see V339 source-license-cases, no whole-package admission",
      "origin": "RETAINED_RESEARCH_EVIDENCE",
      "payload_path": null,
      "bytes": 1161,
      "sha256": "f36671a5487c9c5050efacb58011c37c24c55a889803cb036cf9d9a6347c1e2d",
      "content_base64": "VGhlIE1JVCBMaWNlbnNlIChNSVQpCgpDb3B5cmlnaHQgKGMpIDIwMjIgUGF1bCBNaWxsZXIgKGh0dHBzOi8vcGF1bG1pbGxyLmNvbSkKQ29weXJpZ2h0IChjKSAyMDE2IFRob21hcyBQb3JuaW4gPHBvcm5pbkBib2xldC5vcmc+CgpQZXJtaXNzaW9uIGlzIGhlcmVieSBncmFudGVkLCBmcmVlIG9mIGNoYXJnZSwgdG8gYW55IHBlcnNvbiBvYnRhaW5pbmcgYSBjb3B5Cm9mIHRoaXMgc29mdHdhcmUgYW5kIGFzc29jaWF0ZWQgZG9jdW1lbnRhdGlvbiBmaWxlcyAodGhlIOKAnFNvZnR3YXJl4oCdKSwgdG8gZGVhbAppbiB0aGUgU29mdHdhcmUgd2l0aG91dCByZXN0cmljdGlvbiwgaW5jbHVkaW5nIHdpdGhvdXQgbGltaXRhdGlvbiB0aGUgcmlnaHRzCnRvIHVzZSwgY29weSwgbW9kaWZ5LCBtZXJnZSwgcHVibGlzaCwgZGlzdHJpYnV0ZSwgc3VibGljZW5zZSwgYW5kL29yIHNlbGwKY29waWVzIG9mIHRoZSBTb2Z0d2FyZSwgYW5kIHRvIHBlcm1pdCBwZXJzb25zIHRvIHdob20gdGhlIFNvZnR3YXJlIGlzCmZ1cm5pc2hlZCB0byBkbyBzbywgc3ViamVjdCB0byB0aGUgZm9sbG93aW5nIGNvbmRpdGlvbnM6CgpUaGUgYWJvdmUgY29weXJpZ2h0IG5vdGljZSBhbmQgdGhpcyBwZXJtaXNzaW9uIG5vdGljZSBzaGFsbCBiZSBpbmNsdWRlZCBpbgphbGwgY29waWVzIG9yIHN1YnN0YW50aWFsIHBvcnRpb25zIG9mIHRoZSBTb2Z0d2FyZS4KClRIRSBTT0ZUV0FSRSBJUyBQUk9WSURFRCDigJxBUyBJU+KAnSwgV0lUSE9VVCBXQVJSQU5UWSBPRiBBTlkgS0lORCwgRVhQUkVTUyBPUgpJTVBMSUVELCBJTkNMVURJTkcgQlVUIE5PVCBMSU1JVEVEIFRPIFRIRSBXQVJSQU5USUVTIE9GIE1FUkNIQU5UQUJJTElUWSwKRklUTkVTUyBGT1IgQSBQQVJUSUNVTEFSIFBVUlBPU0UgQU5EIE5PTklORlJJTkdFTUVOVC4gSU4gTk8gRVZFTlQgU0hBTEwgVEhFCkFVVEhPUlMgT1IgQ09QWVJJR0hUIEhPTERFUlMgQkUgTElBQkxFIEZPUiBBTlkgQ0xBSU0sIERBTUFHRVMgT1IgT1RIRVIKTElBQklMSVRZLCBXSEVUSEVSIElOIEFOIEFDVElPTiBPRiBDT05UUkFDVCwgVE9SVCBPUiBPVEhFUldJU0UsIEFSSVNJTkcgRlJPTSwKT1VUIE9GIE9SIElOIENPTk5FQ1RJT04gV0lUSCBUSEUgU09GVFdBUkUgT1IgVEhFIFVTRSBPUiBPVEhFUiBERUFMSU5HUyBJTgpUSEUgU09GVFdBUkUu"
    },
    {
      "path": "supplemental-v339/noble-curves-LICENSE.txt",
      "source_locator": "Temp/elite-v339-9266379aff0044498a1d7004267002c8/noble-curves-LICENSE.txt",
      "scope": "Pinned source licence; see V339 source-license-cases, no whole-package admission",
      "origin": "RETAINED_RESEARCH_EVIDENCE",
      "payload_path": null,
      "bytes": 1109,
      "sha256": "4f221aee6e072336700c408c68ab3b96a3fc09f6aebe6f48f1bd99e5ef13faec",
      "content_base64": "VGhlIE1JVCBMaWNlbnNlIChNSVQpCgpDb3B5cmlnaHQgKGMpIDIwMjIgUGF1bCBNaWxsZXIgKGh0dHBzOi8vcGF1bG1pbGxyLmNvbSkKClBlcm1pc3Npb24gaXMgaGVyZWJ5IGdyYW50ZWQsIGZyZWUgb2YgY2hhcmdlLCB0byBhbnkgcGVyc29uIG9idGFpbmluZyBhIGNvcHkKb2YgdGhpcyBzb2Z0d2FyZSBhbmQgYXNzb2NpYXRlZCBkb2N1bWVudGF0aW9uIGZpbGVzICh0aGUg4oCcU29mdHdhcmXigJ0pLCB0byBkZWFsCmluIHRoZSBTb2Z0d2FyZSB3aXRob3V0IHJlc3RyaWN0aW9uLCBpbmNsdWRpbmcgd2l0aG91dCBsaW1pdGF0aW9uIHRoZSByaWdodHMKdG8gdXNlLCBjb3B5LCBtb2RpZnksIG1lcmdlLCBwdWJsaXNoLCBkaXN0cmlidXRlLCBzdWJsaWNlbnNlLCBhbmQvb3Igc2VsbApjb3BpZXMgb2YgdGhlIFNvZnR3YXJlLCBhbmQgdG8gcGVybWl0IHBlcnNvbnMgdG8gd2hvbSB0aGUgU29mdHdhcmUgaXMKZnVybmlzaGVkIHRvIGRvIHNvLCBzdWJqZWN0IHRvIHRoZSBmb2xsb3dpbmcgY29uZGl0aW9uczoKClRoZSBhYm92ZSBjb3B5cmlnaHQgbm90aWNlIGFuZCB0aGlzIHBlcm1pc3Npb24gbm90aWNlIHNoYWxsIGJlIGluY2x1ZGVkIGluCmFsbCBjb3BpZXMgb3Igc3Vic3RhbnRpYWwgcG9ydGlvbnMgb2YgdGhlIFNvZnR3YXJlLgoKVEhFIFNPRlRXQVJFIElTIFBST1ZJREVEIOKAnEFTIElT4oCdLCBXSVRIT1VUIFdBUlJBTlRZIE9GIEFOWSBLSU5ELCBFWFBSRVNTIE9SCklNUExJRUQsIElOQ0xVRElORyBCVVQgTk9UIExJTUlURUQgVE8gVEhFIFdBUlJBTlRJRVMgT0YgTUVSQ0hBTlRBQklMSVRZLApGSVRORVNTIEZPUiBBIFBBUlRJQ1VMQVIgUFVSUE9TRSBBTkQgTk9OSU5GUklOR0VNRU5ULiBJTiBOTyBFVkVOVCBTSEFMTCBUSEUKQVVUSE9SUyBPUiBDT1BZUklHSFQgSE9MREVSUyBCRSBMSUFCTEUgRk9SIEFOWSBDTEFJTSwgREFNQUdFUyBPUiBPVEhFUgpMSUFCSUxJVFksIFdIRVRIRVIgSU4gQU4gQUNUSU9OIE9GIENPTlRSQUNULCBUT1JUIE9SIE9USEVSV0lTRSwgQVJJU0lORyBGUk9NLApPVVQgT0YgT1IgSU4gQ09OTkVDVElPTiBXSVRIIFRIRSBTT0ZUV0FSRSBPUiBUSEUgVVNFIE9SIE9USEVSIERFQUxJTkdTIElOClRIRSBTT0ZUV0FSRS4="
    },
    {
      "path": "supplemental-v339/noble-hashes-LICENSE.txt",
      "source_locator": "Temp/elite-v339-9266379aff0044498a1d7004267002c8/noble-hashes-LICENSE.txt",
      "scope": "Pinned source licence; see V339 source-license-cases, no whole-package admission",
      "origin": "RETAINED_RESEARCH_EVIDENCE",
      "payload_path": null,
      "bytes": 1109,
      "sha256": "4f221aee6e072336700c408c68ab3b96a3fc09f6aebe6f48f1bd99e5ef13faec",
      "content_base64": "VGhlIE1JVCBMaWNlbnNlIChNSVQpCgpDb3B5cmlnaHQgKGMpIDIwMjIgUGF1bCBNaWxsZXIgKGh0dHBzOi8vcGF1bG1pbGxyLmNvbSkKClBlcm1pc3Npb24gaXMgaGVyZWJ5IGdyYW50ZWQsIGZyZWUgb2YgY2hhcmdlLCB0byBhbnkgcGVyc29uIG9idGFpbmluZyBhIGNvcHkKb2YgdGhpcyBzb2Z0d2FyZSBhbmQgYXNzb2NpYXRlZCBkb2N1bWVudGF0aW9uIGZpbGVzICh0aGUg4oCcU29mdHdhcmXigJ0pLCB0byBkZWFsCmluIHRoZSBTb2Z0d2FyZSB3aXRob3V0IHJlc3RyaWN0aW9uLCBpbmNsdWRpbmcgd2l0aG91dCBsaW1pdGF0aW9uIHRoZSByaWdodHMKdG8gdXNlLCBjb3B5LCBtb2RpZnksIG1lcmdlLCBwdWJsaXNoLCBkaXN0cmlidXRlLCBzdWJsaWNlbnNlLCBhbmQvb3Igc2VsbApjb3BpZXMgb2YgdGhlIFNvZnR3YXJlLCBhbmQgdG8gcGVybWl0IHBlcnNvbnMgdG8gd2hvbSB0aGUgU29mdHdhcmUgaXMKZnVybmlzaGVkIHRvIGRvIHNvLCBzdWJqZWN0IHRvIHRoZSBmb2xsb3dpbmcgY29uZGl0aW9uczoKClRoZSBhYm92ZSBjb3B5cmlnaHQgbm90aWNlIGFuZCB0aGlzIHBlcm1pc3Npb24gbm90aWNlIHNoYWxsIGJlIGluY2x1ZGVkIGluCmFsbCBjb3BpZXMgb3Igc3Vic3RhbnRpYWwgcG9ydGlvbnMgb2YgdGhlIFNvZnR3YXJlLgoKVEhFIFNPRlRXQVJFIElTIFBST1ZJREVEIOKAnEFTIElT4oCdLCBXSVRIT1VUIFdBUlJBTlRZIE9GIEFOWSBLSU5ELCBFWFBSRVNTIE9SCklNUExJRUQsIElOQ0xVRElORyBCVVQgTk9UIExJTUlURUQgVE8gVEhFIFdBUlJBTlRJRVMgT0YgTUVSQ0hBTlRBQklMSVRZLApGSVRORVNTIEZPUiBBIFBBUlRJQ1VMQVIgUFVSUE9TRSBBTkQgTk9OSU5GUklOR0VNRU5ULiBJTiBOTyBFVkVOVCBTSEFMTCBUSEUKQVVUSE9SUyBPUiBDT1BZUklHSFQgSE9MREVSUyBCRSBMSUFCTEUgRk9SIEFOWSBDTEFJTSwgREFNQUdFUyBPUiBPVEhFUgpMSUFCSUxJVFksIFdIRVRIRVIgSU4gQU4gQUNUSU9OIE9GIENPTlRSQUNULCBUT1JUIE9SIE9USEVSV0lTRSwgQVJJU0lORyBGUk9NLApPVVQgT0YgT1IgSU4gQ09OTkVDVElPTiBXSVRIIFRIRSBTT0ZUV0FSRSBPUiBUSEUgVVNFIE9SIE9USEVSIERFQUxJTkdTIElOClRIRSBTT0ZUV0FSRS4="
    },
    {
      "path": "supplemental-v340/undici-formdata-polyfill-attribution.txt",
      "source_locator": "Temp/elite-v340-9cc3d8c83aa4416cbdcd5f3f3b158fe7/undici-formdata-polyfill-attribution.txt",
      "scope": "source licence or attribution; qualification in V340 versioned-source-cases; not automatic clearance",
      "origin": "RETAINED_RESEARCH_EVIDENCE",
      "payload_path": null,
      "bytes": 91,
      "sha256": "2337029828ec5a04414ff9716cd9eb46a19a042d6a2c4d88c649d20df07d6a42",
      "content_base64": "LyohIGZvcm1kYXRhLXBvbHlmaWxsLiBNSVQgTGljZW5zZS4gSmltbXkgV8OkcnRpbmcgPGh0dHBzOi8vamltbXkud2FydGluZy5zZS9vcGVuc291cmNlPiAqLw=="
    },
    {
      "path": "supplemental-v340/undici-ws-attribution.txt",
      "source_locator": "Temp/elite-v340-9cc3d8c83aa4416cbdcd5f3f3b158fe7/undici-ws-attribution.txt",
      "scope": "source licence or attribution; qualification in V340 versioned-source-cases; not automatic clearance",
      "origin": "RETAINED_RESEARCH_EVIDENCE",
      "payload_path": null,
      "bytes": 63,
      "sha256": "3e981c50590c1237334f5399b4747edcb8dc23f509df80ad67b32f29caec69d4",
      "content_base64": "LyohIHdzLiBNSVQgTGljZW5zZS4gRWluYXIgT3R0byBTdGFuZ3ZpayA8ZWluYXJvc0BnbWFpbC5jb20+ICov"
    },
    {
      "path": "supplemental-v340/yarn-node-contributors-MIT.txt",
      "source_locator": "Temp/elite-v340-9cc3d8c83aa4416cbdcd5f3f3b158fe7/yarn-node-contributors-MIT.txt",
      "scope": "source licence or attribution; qualification in V340 versioned-source-cases; not automatic clearance",
      "origin": "RETAINED_RESEARCH_EVIDENCE",
      "payload_path": null,
      "bytes": 1126,
      "sha256": "389721ac6dbc529aa431f9d2a5c4f215383fb7a9a7208a13258706704a84ac43",
      "content_base64": "LyoqCiAgQGxpY2Vuc2UKICBDb3B5cmlnaHQgTm9kZS5qcyBjb250cmlidXRvcnMuIEFsbCByaWdodHMgcmVzZXJ2ZWQuCgogIFBlcm1pc3Npb24gaXMgaGVyZWJ5IGdyYW50ZWQsIGZyZWUgb2YgY2hhcmdlLCB0byBhbnkgcGVyc29uIG9idGFpbmluZyBhIGNvcHkKICBvZiB0aGlzIHNvZnR3YXJlIGFuZCBhc3NvY2lhdGVkIGRvY3VtZW50YXRpb24gZmlsZXMgKHRoZSAiU29mdHdhcmUiKSwgdG8KICBkZWFsIGluIHRoZSBTb2Z0d2FyZSB3aXRob3V0IHJlc3RyaWN0aW9uLCBpbmNsdWRpbmcgd2l0aG91dCBsaW1pdGF0aW9uIHRoZQogIHJpZ2h0cyB0byB1c2UsIGNvcHksIG1vZGlmeSwgbWVyZ2UsIHB1Ymxpc2gsIGRpc3RyaWJ1dGUsIHN1YmxpY2Vuc2UsIGFuZC9vcgogIHNlbGwgY29waWVzIG9mIHRoZSBTb2Z0d2FyZSwgYW5kIHRvIHBlcm1pdCBwZXJzb25zIHRvIHdob20gdGhlIFNvZnR3YXJlIGlzCiAgZnVybmlzaGVkIHRvIGRvIHNvLCBzdWJqZWN0IHRvIHRoZSBmb2xsb3dpbmcgY29uZGl0aW9uczoKCiAgVGhlIGFib3ZlIGNvcHlyaWdodCBub3RpY2UgYW5kIHRoaXMgcGVybWlzc2lvbiBub3RpY2Ugc2hhbGwgYmUgaW5jbHVkZWQgaW4KICBhbGwgY29waWVzIG9yIHN1YnN0YW50aWFsIHBvcnRpb25zIG9mIHRoZSBTb2Z0d2FyZS4KCiAgVEhFIFNPRlRXQVJFIElTIFBST1ZJREVEICJBUyBJUyIsIFdJVEhPVVQgV0FSUkFOVFkgT0YgQU5ZIEtJTkQsIEVYUFJFU1MgT1IKICBJTVBMSUVELCBJTkNMVURJTkcgQlVUIE5PVCBMSU1JVEVEIFRPIFRIRSBXQVJSQU5USUVTIE9GIE1FUkNIQU5UQUJJTElUWSwKICBGSVRORVNTIEZPUiBBIFBBUlRJQ1VMQVIgUFVSUE9TRSBBTkQgTk9OSU5GUklOR0VNRU5ULiBJTiBOTyBFVkVOVCBTSEFMTCBUSEUKICBBVVRIT1JTIE9SIENPUFlSSUdIVCBIT0xERVJTIEJFIExJQUJMRSBGT1IgQU5ZIENMQUlNLCBEQU1BR0VTIE9SIE9USEVSCiAgTElBQklMSVRZLCBXSEVUSEVSIElOIEFOIEFDVElPTiBPRiBDT05UUkFDVCwgVE9SVCBPUiBPVEhFUldJU0UsIEFSSVNJTkcKICBGUk9NLCBPVVQgT0YgT1IgSU4gQ09OTkVDVElPTiBXSVRIIFRIRSBTT0ZUV0FSRSBPUiBUSEUgVVNFIE9SIE9USEVSIERFQUxJTkdTCiAgSU4gVEhFIFNPRlRXQVJFLgoqLw=="
    },
    {
      "path": "supplemental-v340/yarn-blake-embrey-MIT.txt",
      "source_locator": "Temp/elite-v340-9cc3d8c83aa4416cbdcd5f3f3b158fe7/yarn-blake-embrey-MIT.txt",
      "scope": "source licence or attribution; qualification in V340 versioned-source-cases; not automatic clearance",
      "origin": "RETAINED_RESEARCH_EVIDENCE",
      "payload_path": null,
      "bytes": 1154,
      "sha256": "505e12262fd3973c023c2225a98b143bd926c3149a632a8a50fcb34efc48c268",
      "content_base64": "LyoqCiAgQGxpY2Vuc2UKICBUaGUgTUlUIExpY2Vuc2UgKE1JVCkKCiAgQ29weXJpZ2h0IChjKSAyMDE0IEJsYWtlIEVtYnJleSAoaGVsbG9AYmxha2VlbWJyZXkuY29tKQoKICBQZXJtaXNzaW9uIGlzIGhlcmVieSBncmFudGVkLCBmcmVlIG9mIGNoYXJnZSwgdG8gYW55IHBlcnNvbiBvYnRhaW5pbmcgYSBjb3B5CiAgb2YgdGhpcyBzb2Z0d2FyZSBhbmQgYXNzb2NpYXRlZCBkb2N1bWVudGF0aW9uIGZpbGVzICh0aGUgIlNvZnR3YXJlIiksIHRvIGRlYWwKICBpbiB0aGUgU29mdHdhcmUgd2l0aG91dCByZXN0cmljdGlvbiwgaW5jbHVkaW5nIHdpdGhvdXQgbGltaXRhdGlvbiB0aGUgcmlnaHRzCiAgdG8gdXNlLCBjb3B5LCBtb2RpZnksIG1lcmdlLCBwdWJsaXNoLCBkaXN0cmlidXRlLCBzdWJsaWNlbnNlLCBhbmQvb3Igc2VsbAogIGNvcGllcyBvZiB0aGUgU29mdHdhcmUsIGFuZCB0byBwZXJtaXQgcGVyc29ucyB0byB3aG9tIHRoZSBTb2Z0d2FyZSBpcwogIGZ1cm5pc2hlZCB0byBkbyBzbywgc3ViamVjdCB0byB0aGUgZm9sbG93aW5nIGNvbmRpdGlvbnM6CgogIFRoZSBhYm92ZSBjb3B5cmlnaHQgbm90aWNlIGFuZCB0aGlzIHBlcm1pc3Npb24gbm90aWNlIHNoYWxsIGJlIGluY2x1ZGVkIGluCiAgYWxsIGNvcGllcyBvciBzdWJzdGFudGlhbCBwb3J0aW9ucyBvZiB0aGUgU29mdHdhcmUuCgogIFRIRSBTT0ZUV0FSRSBJUyBQUk9WSURFRCAiQVMgSVMiLCBXSVRIT1VUIFdBUlJBTlRZIE9GIEFOWSBLSU5ELCBFWFBSRVNTIE9SCiAgSU1QTElFRCwgSU5DTFVESU5HIEJVVCBOT1QgTElNSVRFRCBUTyBUSEUgV0FSUkFOVElFUyBPRiBNRVJDSEFOVEFCSUxJVFksCiAgRklUTkVTUyBGT1IgQSBQQVJUSUNVTEFSIFBVUlBPU0UgQU5EIE5PTklORlJJTkdFTUVOVC4gSU4gTk8gRVZFTlQgU0hBTEwgVEhFCiAgQVVUSE9SUyBPUiBDT1BZUklHSFQgSE9MREVSUyBCRSBMSUFCTEUgRk9SIEFOWSBDTEFJTSwgREFNQUdFUyBPUiBPVEhFUgogIExJQUJJTElUWSwgV0hFVEhFUiBJTiBBTiBBQ1RJT04gT0YgQ09OVFJBQ1QsIFRPUlQgT1IgT1RIRVJXSVNFLCBBUklTSU5HIEZST00sCiAgT1VUIE9GIE9SIElOIENPTk5FQ1RJT04gV0lUSCBUSEUgU09GVFdBUkUgT1IgVEhFIFVTRSBPUiBPVEhFUiBERUFMSU5HUyBJTgogIFRIRSBTT0ZUV0FSRS4KKi8="
    },
    {
      "path": "supplemental-v340/yarn-joyent-node-MIT.txt",
      "source_locator": "Temp/elite-v340-9cc3d8c83aa4416cbdcd5f3f3b158fe7/yarn-joyent-node-MIT.txt",
      "scope": "source licence or attribution; qualification in V340 versioned-source-cases; not automatic clearance",
      "origin": "RETAINED_RESEARCH_EVIDENCE",
      "payload_path": null,
      "bytes": 1127,
      "sha256": "14895e2d944ac3300cc2daf6208ed55982e1690a6a34090078199a0b6ac07110",
      "content_base64": "LyoqCiAgQGxpY2Vuc2UKICBDb3B5cmlnaHQgSm95ZW50LCBJbmMuIGFuZCBvdGhlciBOb2RlIGNvbnRyaWJ1dG9ycy4KCiAgUGVybWlzc2lvbiBpcyBoZXJlYnkgZ3JhbnRlZCwgZnJlZSBvZiBjaGFyZ2UsIHRvIGFueSBwZXJzb24gb2J0YWluaW5nIGEKICBjb3B5IG9mIHRoaXMgc29mdHdhcmUgYW5kIGFzc29jaWF0ZWQgZG9jdW1lbnRhdGlvbiBmaWxlcyAodGhlCiAgIlNvZnR3YXJlIiksIHRvIGRlYWwgaW4gdGhlIFNvZnR3YXJlIHdpdGhvdXQgcmVzdHJpY3Rpb24sIGluY2x1ZGluZwogIHdpdGhvdXQgbGltaXRhdGlvbiB0aGUgcmlnaHRzIHRvIHVzZSwgY29weSwgbW9kaWZ5LCBtZXJnZSwgcHVibGlzaCwKICBkaXN0cmlidXRlLCBzdWJsaWNlbnNlLCBhbmQvb3Igc2VsbCBjb3BpZXMgb2YgdGhlIFNvZnR3YXJlLCBhbmQgdG8gcGVybWl0CiAgcGVyc29ucyB0byB3aG9tIHRoZSBTb2Z0d2FyZSBpcyBmdXJuaXNoZWQgdG8gZG8gc28sIHN1YmplY3QgdG8gdGhlCiAgZm9sbG93aW5nIGNvbmRpdGlvbnM6CgogIFRoZSBhYm92ZSBjb3B5cmlnaHQgbm90aWNlIGFuZCB0aGlzIHBlcm1pc3Npb24gbm90aWNlIHNoYWxsIGJlIGluY2x1ZGVkCiAgaW4gYWxsIGNvcGllcyBvciBzdWJzdGFudGlhbCBwb3J0aW9ucyBvZiB0aGUgU29mdHdhcmUuCgogIFRIRSBTT0ZUV0FSRSBJUyBQUk9WSURFRCAiQVMgSVMiLCBXSVRIT1VUIFdBUlJBTlRZIE9GIEFOWSBLSU5ELCBFWFBSRVNTCiAgT1IgSU1QTElFRCwgSU5DTFVESU5HIEJVVCBOT1QgTElNSVRFRCBUTyBUSEUgV0FSUkFOVElFUyBPRgogIE1FUkNIQU5UQUJJTElUWSwgRklUTkVTUyBGT1IgQSBQQVJUSUNVTEFSIFBVUlBPU0UgQU5EIE5PTklORlJJTkdFTUVOVC4gSU4KICBOTyBFVkVOVCBTSEFMTCBUSEUgQVVUSE9SUyBPUiBDT1BZUklHSFQgSE9MREVSUyBCRSBMSUFCTEUgRk9SIEFOWSBDTEFJTSwKICBEQU1BR0VTIE9SIE9USEVSIExJQUJJTElUWSwgV0hFVEhFUiBJTiBBTiBBQ1RJT04gT0YgQ09OVFJBQ1QsIFRPUlQgT1IKICBPVEhFUldJU0UsIEFSSVNJTkcgRlJPTSwgT1VUIE9GIE9SIElOIENPTk5FQ1RJT04gV0lUSCBUSEUgU09GVFdBUkUgT1IgVEhFCiAgVVNFIE9SIE9USEVSIERFQUxJTkdTIElOIFRIRSBTT0ZUV0FSRS4KKi8="
    },
    {
      "path": "supplemental-v340/undici-7.29.0-LICENSE.evidence",
      "source_locator": "Temp/elite-v340-9cc3d8c83aa4416cbdcd5f3f3b158fe7/undici-7.29.0-LICENSE.evidence",
      "scope": "source licence or attribution; qualification in V340 versioned-source-cases; not automatic clearance",
      "origin": "RETAINED_RESEARCH_EVIDENCE",
      "payload_path": null,
      "bytes": 1090,
      "sha256": "a6db8096b2707bc0102d256917d4d33f298ba36d8c3f25de067a2b5bb379db27",
      "content_base64": "TUlUIExpY2Vuc2UKCkNvcHlyaWdodCAoYykgTWF0dGVvIENvbGxpbmEgYW5kIFVuZGljaSBjb250cmlidXRvcnMKClBlcm1pc3Npb24gaXMgaGVyZWJ5IGdyYW50ZWQsIGZyZWUgb2YgY2hhcmdlLCB0byBhbnkgcGVyc29uIG9idGFpbmluZyBhIGNvcHkKb2YgdGhpcyBzb2Z0d2FyZSBhbmQgYXNzb2NpYXRlZCBkb2N1bWVudGF0aW9uIGZpbGVzICh0aGUgIlNvZnR3YXJlIiksIHRvIGRlYWwKaW4gdGhlIFNvZnR3YXJlIHdpdGhvdXQgcmVzdHJpY3Rpb24sIGluY2x1ZGluZyB3aXRob3V0IGxpbWl0YXRpb24gdGhlIHJpZ2h0cwp0byB1c2UsIGNvcHksIG1vZGlmeSwgbWVyZ2UsIHB1Ymxpc2gsIGRpc3RyaWJ1dGUsIHN1YmxpY2Vuc2UsIGFuZC9vciBzZWxsCmNvcGllcyBvZiB0aGUgU29mdHdhcmUsIGFuZCB0byBwZXJtaXQgcGVyc29ucyB0byB3aG9tIHRoZSBTb2Z0d2FyZSBpcwpmdXJuaXNoZWQgdG8gZG8gc28sIHN1YmplY3QgdG8gdGhlIGZvbGxvd2luZyBjb25kaXRpb25zOgoKVGhlIGFib3ZlIGNvcHlyaWdodCBub3RpY2UgYW5kIHRoaXMgcGVybWlzc2lvbiBub3RpY2Ugc2hhbGwgYmUgaW5jbHVkZWQgaW4gYWxsCmNvcGllcyBvciBzdWJzdGFudGlhbCBwb3J0aW9ucyBvZiB0aGUgU29mdHdhcmUuCgpUSEUgU09GVFdBUkUgSVMgUFJPVklERUQgIkFTIElTIiwgV0lUSE9VVCBXQVJSQU5UWSBPRiBBTlkgS0lORCwgRVhQUkVTUyBPUgpJTVBMSUVELCBJTkNMVURJTkcgQlVUIE5PVCBMSU1JVEVEIFRPIFRIRSBXQVJSQU5USUVTIE9GIE1FUkNIQU5UQUJJTElUWSwKRklUTkVTUyBGT1IgQSBQQVJUSUNVTEFSIFBVUlBPU0UgQU5EIE5PTklORlJJTkdFTUVOVC4gSU4gTk8gRVZFTlQgU0hBTEwgVEhFCkFVVEhPUlMgT1IgQ09QWVJJR0hUIEhPTERFUlMgQkUgTElBQkxFIEZPUiBBTlkgQ0xBSU0sIERBTUFHRVMgT1IgT1RIRVIKTElBQklMSVRZLCBXSEVUSEVSIElOIEFOIEFDVElPTiBPRiBDT05UUkFDVCwgVE9SVCBPUiBPVEhFUldJU0UsIEFSSVNJTkcgRlJPTSwKT1VUIE9GIE9SIElOIENPTk5FQ1RJT04gV0lUSCBUSEUgU09GVFdBUkUgT1IgVEhFIFVTRSBPUiBPVEhFUiBERUFMSU5HUyBJTiBUSEUKU09GVFdBUkUuCg=="
    },
    {
      "path": "supplemental-v340/undici-7.29.0-lib__web__fetch__LICENSE.evidence",
      "source_locator": "Temp/elite-v340-9cc3d8c83aa4416cbdcd5f3f3b158fe7/undici-7.29.0-lib__web__fetch__LICENSE.evidence",
      "scope": "source licence or attribution; qualification in V340 versioned-source-cases; not automatic clearance",
      "origin": "RETAINED_RESEARCH_EVIDENCE",
      "payload_path": null,
      "bytes": 1071,
      "sha256": "a21d6c8d3fc631198b97e57584697f7f9c37805c71c6cf9a12e4442b40394b88",
      "content_base64": "TUlUIExpY2Vuc2UKCkNvcHlyaWdodCAoYykgMjAyMCBFdGhhbiBBcnJvd29vZAoKUGVybWlzc2lvbiBpcyBoZXJlYnkgZ3JhbnRlZCwgZnJlZSBvZiBjaGFyZ2UsIHRvIGFueSBwZXJzb24gb2J0YWluaW5nIGEgY29weQpvZiB0aGlzIHNvZnR3YXJlIGFuZCBhc3NvY2lhdGVkIGRvY3VtZW50YXRpb24gZmlsZXMgKHRoZSAiU29mdHdhcmUiKSwgdG8gZGVhbAppbiB0aGUgU29mdHdhcmUgd2l0aG91dCByZXN0cmljdGlvbiwgaW5jbHVkaW5nIHdpdGhvdXQgbGltaXRhdGlvbiB0aGUgcmlnaHRzCnRvIHVzZSwgY29weSwgbW9kaWZ5LCBtZXJnZSwgcHVibGlzaCwgZGlzdHJpYnV0ZSwgc3VibGljZW5zZSwgYW5kL29yIHNlbGwKY29waWVzIG9mIHRoZSBTb2Z0d2FyZSwgYW5kIHRvIHBlcm1pdCBwZXJzb25zIHRvIHdob20gdGhlIFNvZnR3YXJlIGlzCmZ1cm5pc2hlZCB0byBkbyBzbywgc3ViamVjdCB0byB0aGUgZm9sbG93aW5nIGNvbmRpdGlvbnM6CgpUaGUgYWJvdmUgY29weXJpZ2h0IG5vdGljZSBhbmQgdGhpcyBwZXJtaXNzaW9uIG5vdGljZSBzaGFsbCBiZSBpbmNsdWRlZCBpbiBhbGwKY29waWVzIG9yIHN1YnN0YW50aWFsIHBvcnRpb25zIG9mIHRoZSBTb2Z0d2FyZS4KClRIRSBTT0ZUV0FSRSBJUyBQUk9WSURFRCAiQVMgSVMiLCBXSVRIT1VUIFdBUlJBTlRZIE9GIEFOWSBLSU5ELCBFWFBSRVNTIE9SCklNUExJRUQsIElOQ0xVRElORyBCVVQgTk9UIExJTUlURUQgVE8gVEhFIFdBUlJBTlRJRVMgT0YgTUVSQ0hBTlRBQklMSVRZLApGSVRORVNTIEZPUiBBIFBBUlRJQ1VMQVIgUFVSUE9TRSBBTkQgTk9OSU5GUklOR0VNRU5ULiBJTiBOTyBFVkVOVCBTSEFMTCBUSEUKQVVUSE9SUyBPUiBDT1BZUklHSFQgSE9MREVSUyBCRSBMSUFCTEUgRk9SIEFOWSBDTEFJTSwgREFNQUdFUyBPUiBPVEhFUgpMSUFCSUxJVFksIFdIRVRIRVIgSU4gQU4gQUNUSU9OIE9GIENPTlRSQUNULCBUT1JUIE9SIE9USEVSV0lTRSwgQVJJU0lORyBGUk9NLApPVVQgT0YgT1IgSU4gQ09OTkVDVElPTiBXSVRIIFRIRSBTT0ZUV0FSRSBPUiBUSEUgVVNFIE9SIE9USEVSIERFQUxJTkdTIElOIFRIRQpTT0ZUV0FSRS4K"
    },
    {
      "path": "supplemental-v340/yarn-root-LICENSE.txt",
      "source_locator": "Temp/elite-v340-9cc3d8c83aa4416cbdcd5f3f3b158fe7/yarn-root-LICENSE.txt",
      "scope": "source licence or attribution; qualification in V340 versioned-source-cases; not automatic clearance",
      "origin": "RETAINED_RESEARCH_EVIDENCE",
      "payload_path": null,
      "bytes": 1336,
      "sha256": "238d933f5c226cc197bd1dae2ad0c468e157b4cba8ed844f81549ba6db777dc4",
      "content_base64": "QlNEIDItQ2xhdXNlIExpY2Vuc2UKCkNvcHlyaWdodCAoYykgMjAxNi1wcmVzZW50LCBZYXJuIENvbnRyaWJ1dG9ycy4KQWxsIHJpZ2h0cyByZXNlcnZlZC4KClJlZGlzdHJpYnV0aW9uIGFuZCB1c2UgaW4gc291cmNlIGFuZCBiaW5hcnkgZm9ybXMsIHdpdGggb3Igd2l0aG91dAptb2RpZmljYXRpb24sIGFyZSBwZXJtaXR0ZWQgcHJvdmlkZWQgdGhhdCB0aGUgZm9sbG93aW5nIGNvbmRpdGlvbnMgYXJlIG1ldDoKCjEuIFJlZGlzdHJpYnV0aW9ucyBvZiBzb3VyY2UgY29kZSBtdXN0IHJldGFpbiB0aGUgYWJvdmUgY29weXJpZ2h0IG5vdGljZSwgdGhpcwogICBsaXN0IG9mIGNvbmRpdGlvbnMgYW5kIHRoZSBmb2xsb3dpbmcgZGlzY2xhaW1lci4KCjIuIFJlZGlzdHJpYnV0aW9ucyBpbiBiaW5hcnkgZm9ybSBtdXN0IHJlcHJvZHVjZSB0aGUgYWJvdmUgY29weXJpZ2h0IG5vdGljZSwKICAgdGhpcyBsaXN0IG9mIGNvbmRpdGlvbnMgYW5kIHRoZSBmb2xsb3dpbmcgZGlzY2xhaW1lciBpbiB0aGUgZG9jdW1lbnRhdGlvbgogICBhbmQvb3Igb3RoZXIgbWF0ZXJpYWxzIHByb3ZpZGVkIHdpdGggdGhlIGRpc3RyaWJ1dGlvbi4KClRISVMgU09GVFdBUkUgSVMgUFJPVklERUQgQlkgVEhFIENPUFlSSUdIVCBIT0xERVJTIEFORCBDT05UUklCVVRPUlMgIkFTIElTIgpBTkQgQU5ZIEVYUFJFU1MgT1IgSU1QTElFRCBXQVJSQU5USUVTLCBJTkNMVURJTkcsIEJVVCBOT1QgTElNSVRFRCBUTywgVEhFCklNUExJRUQgV0FSUkFOVElFUyBPRiBNRVJDSEFOVEFCSUxJVFkgQU5EIEZJVE5FU1MgRk9SIEEgUEFSVElDVUxBUiBQVVJQT1NFIEFSRQpESVNDTEFJTUVELiBJTiBOTyBFVkVOVCBTSEFMTCBUSEUgQ09QWVJJR0hUIEhPTERFUiBPUiBDT05UUklCVVRPUlMgQkUgTElBQkxFCkZPUiBBTlkgRElSRUNULCBJTkRJUkVDVCwgSU5DSURFTlRBTCwgU1BFQ0lBTCwgRVhFTVBMQVJZLCBPUiBDT05TRVFVRU5USUFMCkRBTUFHRVMgKElOQ0xVRElORywgQlVUIE5PVCBMSU1JVEVEIFRPLCBQUk9DVVJFTUVOVCBPRiBTVUJTVElUVVRFIEdPT0RTIE9SClNFUlZJQ0VTOyBMT1NTIE9GIFVTRSwgREFUQSwgT1IgUFJPRklUUzsgT1IgQlVTSU5FU1MgSU5URVJSVVBUSU9OKSBIT1dFVkVSCkNBVVNFRCBBTkQgT04gQU5ZIFRIRU9SWSBPRiBMSUFCSUxJVFksIFdIRVRIRVIgSU4gQ09OVFJBQ1QsIFNUUklDVCBMSUFCSUxJVFksCk9SIFRPUlQgKElOQ0xVRElORyBORUdMSUdFTkNFIE9SIE9USEVSV0lTRSkgQVJJU0lORyBJTiBBTlkgV0FZIE9VVCBPRiBUSEUgVVNFCk9GIFRISVMgU09GVFdBUkUsIEVWRU4gSUYgQURWSVNFRCBPRiBUSEUgUE9TU0lCSUxJVFkgT0YgU1VDSCBEQU1BR0UuCg=="
    }
  ]
}
````

### FILE: `pnpm_artifact_selection/next-path-source-evidence.json`
```yaml
block_id: "PNPM-ARTIFACT-SELECTION-GATE:next-path-source-evidence:v398"
operation: CREATE
provenance: ADAPTED
source: "AUTHORED evidence envelope with VERBATIM next-path1.0.0 index.js/package.json/LICENSE from official commit ce2c1386836339bc473b76c9888947899ead8d56; Git blobs and original SHA256 retained, encoded as data and never installed/executed; locally authored scoped module correspondence receipt"
license: "LicenseRef-Workspace-Owner AND MPL-2.0"
sha256: "581c15a6e63ecc37d3c4f4fde4e701ed25532ffc141c519dd88919baf8685dd7"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-next-path-source-evidence/v1",
  "package": "next-path",
  "version": "1.0.0",
  "license": "MPL-2.0",
  "commit": "ce2c1386836339bc473b76c9888947899ead8d56",
  "status": "SOURCE_AND_LICENSE_DELIVERY_ONLY",
  "runtime_admitted": false,
  "files": [
    {
      "path": "index.js",
      "bytes": 317,
      "git_blob": "18367e387bf296e217ae554b0c3b5799d47536a9",
      "sha256": "ef96862a543004139a10f8f52f0a45b4c972d0c4aa26e85396fd7aa2575d290e",
      "source_url": "https://raw.githubusercontent.com/sholladay/next-path/ce2c1386836339bc473b76c9888947899ead8d56/index.js",
      "content_base64": "J3VzZSBzdHJpY3QnOwoKY29uc3QgcGF0aCA9IHJlcXVpcmUoJ3BhdGgnKTsKCmNvbnN0IG5leHRQYXRoID0gKGZyb20sIHRvKSA9PiB7CiAgICBjb25zdCBkaWZmID0gcGF0aC5yZWxhdGl2ZShmcm9tLCB0byk7CiAgICBjb25zdCBzZXBJbmRleCA9IGRpZmYuaW5kZXhPZihwYXRoLnNlcCk7CiAgICBjb25zdCBuZXh0ID0gc2VwSW5kZXggPj0gMCA/CiAgICAgICAgZGlmZi5zdWJzdHJpbmcoMCwgc2VwSW5kZXgpIDoKICAgICAgICBkaWZmOwoKICAgIHJldHVybiBwYXRoLmpvaW4oZnJvbSwgbmV4dCk7Cn07Cgptb2R1bGUuZXhwb3J0cyA9IG5leHRQYXRoOwo="
    },
    {
      "path": "package.json",
      "bytes": 1010,
      "git_blob": "839a4dec0b49dcc4e799bcb39dd3e4f9996ebdf1",
      "sha256": "07d01d2d44846a76138c89985a17eda58358db997b5a25cdc7e87163aadcfdf4",
      "source_url": "https://raw.githubusercontent.com/sholladay/next-path/ce2c1386836339bc473b76c9888947899ead8d56/package.json",
      "content_base64": "ewogICJuYW1lIjogIm5leHQtcGF0aCIsCiAgInZlcnNpb24iOiAiMS4wLjAiLAogICJkZXNjcmlwdGlvbiI6ICJPbmUgc3RlcCBjbG9zZXIgdG8geW91ciBkZXN0aW5hdGlvbi4iLAogICJob21lcGFnZSI6ICJodHRwczovL2dpdGh1Yi5jb20vc2hvbGxhZGF5L25leHQtcGF0aCIsCiAgIm1haW4iOiAiaW5kZXguanMiLAogICJhdXRob3IiOiB7CiAgICAibmFtZSI6ICJTZXRoIEhvbGxhZGF5IiwKICAgICJ1cmwiOiAiaHR0cDovL3NldGgtaG9sbGFkYXkuY29tIiwKICAgICJlbWFpbCI6ICJtZUBzZXRoLWhvbGxhZGF5LmNvbSIKICB9LAogICJzY3JpcHRzIjogewogICAgInRlc3QiOiAieG8gJiYgYXZhIgogIH0sCiAgInJlcG9zaXRvcnkiOiB7CiAgICAidHlwZSI6ICJnaXQiLAogICAgInVybCI6ICJnaXRAZ2l0aHViLmNvbTpzaG9sbGFkYXkvbmV4dC1wYXRoLmdpdCIKICB9LAogICJidWdzIjogewogICAgInVybCI6ICJodHRwczovL2dpdGh1Yi5jb20vc2hvbGxhZGF5L25leHQtcGF0aC9pc3N1ZXMiLAogICAgImVtYWlsIjogIm1lQHNldGgtaG9sbGFkYXkuY29tIgogIH0sCiAgImVuZ2luZXMiOiB7CiAgICAibm9kZSI6ICI+PTYiCiAgfSwKICAibGljZW5zZSI6ICJNUEwtMi4wIiwKICAiZmlsZXMiOiBbCiAgICAiaW5kZXguanMiCiAgXSwKICAiZGV2RGVwZW5kZW5jaWVzIjogewogICAgImF2YSI6ICJeMC4xOC4xIiwKICAgICJlc2xpbnQtY29uZmlnLXRpZHkiOiAiXjAuNC4xIiwKICAgICJ4byI6ICJeMC4xNy4xIgogIH0sCiAgImtleXdvcmRzIjogWwogICAgImFkZCIsCiAgICAiY29uY2F0IiwKICAgICJuZXh0IiwKICAgICJwYXRoIiwKICAgICJwYXRocyIsCiAgICAic2VnbWVudCIsCiAgICAiY29tcG9uZW50IiwKICAgICJyZXNvbHZlIiwKICAgICJkaXN0YW5jZSIsCiAgICAiYmV0d2VlbiIsCiAgICAicmVsYXRpdmUiLAogICAgImRpciIsCiAgICAiZGlyZWN0b3J5IiwKICAgICJmaWxlIgogIF0sCiAgInhvIjogewogICAgImV4dGVuZCI6ICJ0aWR5IgogIH0KfQo="
    },
    {
      "path": "LICENSE",
      "bytes": 15978,
      "git_blob": "3daaa01ab4d1e6713c7e741bdd5b33e5e851eb87",
      "sha256": "786ef75c24eb986a2ffe31eb878b51a16826af1f5c9bd5cbdb5ad9a4223b5ab7",
      "source_url": "https://raw.githubusercontent.com/sholladay/next-path/ce2c1386836339bc473b76c9888947899ead8d56/LICENSE",
      "content_base64": "Q29weXJpZ2h0IChDKSAyMDE2IFNldGggSG9sbGFkYXksIG1lQHNldGgtaG9sbGFkYXkuY29tCgpNb3ppbGxhIFB1YmxpYyBMaWNlbnNlLCB2ZXJzaW9uIDIuMAoKMS4gRGVmaW5pdGlvbnMKCjEuMS4gIkNvbnRyaWJ1dG9yIgoKICAgICBtZWFucyBlYWNoIGluZGl2aWR1YWwgb3IgbGVnYWwgZW50aXR5IHRoYXQgY3JlYXRlcywgY29udHJpYnV0ZXMgdG8gdGhlCiAgICAgY3JlYXRpb24gb2YsIG9yIG93bnMgQ292ZXJlZCBTb2Z0d2FyZS4KCjEuMi4gIkNvbnRyaWJ1dG9yIFZlcnNpb24iCgogICAgIG1lYW5zIHRoZSBjb21iaW5hdGlvbiBvZiB0aGUgQ29udHJpYnV0aW9ucyBvZiBvdGhlcnMgKGlmIGFueSkgdXNlZCBieSBhCiAgICAgQ29udHJpYnV0b3IgYW5kIHRoYXQgcGFydGljdWxhciBDb250cmlidXRvcidzIENvbnRyaWJ1dGlvbi4KCjEuMy4gIkNvbnRyaWJ1dGlvbiIKCiAgICAgbWVhbnMgQ292ZXJlZCBTb2Z0d2FyZSBvZiBhIHBhcnRpY3VsYXIgQ29udHJpYnV0b3IuCgoxLjQuICJDb3ZlcmVkIFNvZnR3YXJlIgoKICAgICBtZWFucyBTb3VyY2UgQ29kZSBGb3JtIHRvIHdoaWNoIHRoZSBpbml0aWFsIENvbnRyaWJ1dG9yIGhhcyBhdHRhY2hlZCB0aGUKICAgICBub3RpY2UgaW4gRXhoaWJpdCBBLCB0aGUgRXhlY3V0YWJsZSBGb3JtIG9mIHN1Y2ggU291cmNlIENvZGUgRm9ybSwgYW5kCiAgICAgTW9kaWZpY2F0aW9ucyBvZiBzdWNoIFNvdXJjZSBDb2RlIEZvcm0sIGluIGVhY2ggY2FzZSBpbmNsdWRpbmcgcG9ydGlvbnMKICAgICB0aGVyZW9mLgoKMS41LiAiSW5jb21wYXRpYmxlIFdpdGggU2Vjb25kYXJ5IExpY2Vuc2VzIgogICAgIG1lYW5zCgogICAgIGEuIHRoYXQgdGhlIGluaXRpYWwgQ29udHJpYnV0b3IgaGFzIGF0dGFjaGVkIHRoZSBub3RpY2UgZGVzY3JpYmVkIGluCiAgICAgICAgRXhoaWJpdCBCIHRvIHRoZSBDb3ZlcmVkIFNvZnR3YXJlOyBvcgoKICAgICBiLiB0aGF0IHRoZSBDb3ZlcmVkIFNvZnR3YXJlIHdhcyBtYWRlIGF2YWlsYWJsZSB1bmRlciB0aGUgdGVybXMgb2YKICAgICAgICB2ZXJzaW9uIDEuMSBvciBlYXJsaWVyIG9mIHRoZSBMaWNlbnNlLCBidXQgbm90IGFsc28gdW5kZXIgdGhlIHRlcm1zIG9mCiAgICAgICAgYSBTZWNvbmRhcnkgTGljZW5zZS4KCjEuNi4gIkV4ZWN1dGFibGUgRm9ybSIKCiAgICAgbWVhbnMgYW55IGZvcm0gb2YgdGhlIHdvcmsgb3RoZXIgdGhhbiBTb3VyY2UgQ29kZSBGb3JtLgoKMS43LiAiTGFyZ2VyIFdvcmsiCgogICAgIG1lYW5zIGEgd29yayB0aGF0IGNvbWJpbmVzIENvdmVyZWQgU29mdHdhcmUgd2l0aCBvdGhlciBtYXRlcmlhbCwgaW4gYQogICAgIHNlcGFyYXRlIGZpbGUgb3IgZmlsZXMsIHRoYXQgaXMgbm90IENvdmVyZWQgU29mdHdhcmUuCgoxLjguICJMaWNlbnNlIgoKICAgICBtZWFucyB0aGlzIGRvY3VtZW50LgoKMS45LiAiTGljZW5zYWJsZSIKCiAgICAgbWVhbnMgaGF2aW5nIHRoZSByaWdodCB0byBncmFudCwgdG8gdGhlIG1heGltdW0gZXh0ZW50IHBvc3NpYmxlLCB3aGV0aGVyCiAgICAgYXQgdGhlIHRpbWUgb2YgdGhlIGluaXRpYWwgZ3JhbnQgb3Igc3Vic2VxdWVudGx5LCBhbnkgYW5kIGFsbCBvZiB0aGUKICAgICByaWdodHMgY29udmV5ZWQgYnkgdGhpcyBMaWNlbnNlLgoKMS4xMC4gIk1vZGlmaWNhdGlvbnMiCgogICAgIG1lYW5zIGFueSBvZiB0aGUgZm9sbG93aW5nOgoKICAgICBhLiBhbnkgZmlsZSBpbiBTb3VyY2UgQ29kZSBGb3JtIHRoYXQgcmVzdWx0cyBmcm9tIGFuIGFkZGl0aW9uIHRvLAogICAgICAgIGRlbGV0aW9uIGZyb20sIG9yIG1vZGlmaWNhdGlvbiBvZiB0aGUgY29udGVudHMgb2YgQ292ZXJlZCBTb2Z0d2FyZTsgb3IKCiAgICAgYi4gYW55IG5ldyBmaWxlIGluIFNvdXJjZSBDb2RlIEZvcm0gdGhhdCBjb250YWlucyBhbnkgQ292ZXJlZCBTb2Z0d2FyZS4KCjEuMTEuICJQYXRlbnQgQ2xhaW1zIiBvZiBhIENvbnRyaWJ1dG9yCgogICAgICBtZWFucyBhbnkgcGF0ZW50IGNsYWltKHMpLCBpbmNsdWRpbmcgd2l0aG91dCBsaW1pdGF0aW9uLCBtZXRob2QsCiAgICAgIHByb2Nlc3MsIGFuZCBhcHBhcmF0dXMgY2xhaW1zLCBpbiBhbnkgcGF0ZW50IExpY2Vuc2FibGUgYnkgc3VjaAogICAgICBDb250cmlidXRvciB0aGF0IHdvdWxkIGJlIGluZnJpbmdlZCwgYnV0IGZvciB0aGUgZ3JhbnQgb2YgdGhlIExpY2Vuc2UsCiAgICAgIGJ5IHRoZSBtYWtpbmcsIHVzaW5nLCBzZWxsaW5nLCBvZmZlcmluZyBmb3Igc2FsZSwgaGF2aW5nIG1hZGUsIGltcG9ydCwKICAgICAgb3IgdHJhbnNmZXIgb2YgZWl0aGVyIGl0cyBDb250cmlidXRpb25zIG9yIGl0cyBDb250cmlidXRvciBWZXJzaW9uLgoKMS4xMi4gIlNlY29uZGFyeSBMaWNlbnNlIgoKICAgICAgbWVhbnMgZWl0aGVyIHRoZSBHTlUgR2VuZXJhbCBQdWJsaWMgTGljZW5zZSwgVmVyc2lvbiAyLjAsIHRoZSBHTlUgTGVzc2VyCiAgICAgIEdlbmVyYWwgUHVibGljIExpY2Vuc2UsIFZlcnNpb24gMi4xLCB0aGUgR05VIEFmZmVybyBHZW5lcmFsIFB1YmxpYwogICAgICBMaWNlbnNlLCBWZXJzaW9uIDMuMCwgb3IgYW55IGxhdGVyIHZlcnNpb25zIG9mIHRob3NlIGxpY2Vuc2VzLgoKMS4xMy4gIlNvdXJjZSBDb2RlIEZvcm0iCgogICAgICBtZWFucyB0aGUgZm9ybSBvZiB0aGUgd29yayBwcmVmZXJyZWQgZm9yIG1ha2luZyBtb2RpZmljYXRpb25zLgoKMS4xNC4gIllvdSIgKG9yICJZb3VyIikKCiAgICAgIG1lYW5zIGFuIGluZGl2aWR1YWwgb3IgYSBsZWdhbCBlbnRpdHkgZXhlcmNpc2luZyByaWdodHMgdW5kZXIgdGhpcwogICAgICBMaWNlbnNlLiBGb3IgbGVnYWwgZW50aXRpZXMsICJZb3UiIGluY2x1ZGVzIGFueSBlbnRpdHkgdGhhdCBjb250cm9scywgaXMKICAgICAgY29udHJvbGxlZCBieSwgb3IgaXMgdW5kZXIgY29tbW9uIGNvbnRyb2wgd2l0aCBZb3UuIEZvciBwdXJwb3NlcyBvZiB0aGlzCiAgICAgIGRlZmluaXRpb24sICJjb250cm9sIiBtZWFucyAoYSkgdGhlIHBvd2VyLCBkaXJlY3Qgb3IgaW5kaXJlY3QsIHRvIGNhdXNlCiAgICAgIHRoZSBkaXJlY3Rpb24gb3IgbWFuYWdlbWVudCBvZiBzdWNoIGVudGl0eSwgd2hldGhlciBieSBjb250cmFjdCBvcgogICAgICBvdGhlcndpc2UsIG9yIChiKSBvd25lcnNoaXAgb2YgbW9yZSB0aGFuIGZpZnR5IHBlcmNlbnQgKDUwJSkgb2YgdGhlCiAgICAgIG91dHN0YW5kaW5nIHNoYXJlcyBvciBiZW5lZmljaWFsIG93bmVyc2hpcCBvZiBzdWNoIGVudGl0eS4KCgoyLiBMaWNlbnNlIEdyYW50cyBhbmQgQ29uZGl0aW9ucwoKMi4xLiBHcmFudHMKCiAgICAgRWFjaCBDb250cmlidXRvciBoZXJlYnkgZ3JhbnRzIFlvdSBhIHdvcmxkLXdpZGUsIHJveWFsdHktZnJlZSwKICAgICBub24tZXhjbHVzaXZlIGxpY2Vuc2U6CgogICAgIGEuIHVuZGVyIGludGVsbGVjdHVhbCBwcm9wZXJ0eSByaWdodHMgKG90aGVyIHRoYW4gcGF0ZW50IG9yIHRyYWRlbWFyaykKICAgICAgICBMaWNlbnNhYmxlIGJ5IHN1Y2ggQ29udHJpYnV0b3IgdG8gdXNlLCByZXByb2R1Y2UsIG1ha2UgYXZhaWxhYmxlLAogICAgICAgIG1vZGlmeSwgZGlzcGxheSwgcGVyZm9ybSwgZGlzdHJpYnV0ZSwgYW5kIG90aGVyd2lzZSBleHBsb2l0IGl0cwogICAgICAgIENvbnRyaWJ1dGlvbnMsIGVpdGhlciBvbiBhbiB1bm1vZGlmaWVkIGJhc2lzLCB3aXRoIE1vZGlmaWNhdGlvbnMsIG9yCiAgICAgICAgYXMgcGFydCBvZiBhIExhcmdlciBXb3JrOyBhbmQKCiAgICAgYi4gdW5kZXIgUGF0ZW50IENsYWltcyBvZiBzdWNoIENvbnRyaWJ1dG9yIHRvIG1ha2UsIHVzZSwgc2VsbCwgb2ZmZXIgZm9yCiAgICAgICAgc2FsZSwgaGF2ZSBtYWRlLCBpbXBvcnQsIGFuZCBvdGhlcndpc2UgdHJhbnNmZXIgZWl0aGVyIGl0cwogICAgICAgIENvbnRyaWJ1dGlvbnMgb3IgaXRzIENvbnRyaWJ1dG9yIFZlcnNpb24uCgoyLjIuIEVmZmVjdGl2ZSBEYXRlCgogICAgIFRoZSBsaWNlbnNlcyBncmFudGVkIGluIFNlY3Rpb24gMi4xIHdpdGggcmVzcGVjdCB0byBhbnkgQ29udHJpYnV0aW9uCiAgICAgYmVjb21lIGVmZmVjdGl2ZSBmb3IgZWFjaCBDb250cmlidXRpb24gb24gdGhlIGRhdGUgdGhlIENvbnRyaWJ1dG9yIGZpcnN0CiAgICAgZGlzdHJpYnV0ZXMgc3VjaCBDb250cmlidXRpb24uCgoyLjMuIExpbWl0YXRpb25zIG9uIEdyYW50IFNjb3BlCgogICAgIFRoZSBsaWNlbnNlcyBncmFudGVkIGluIHRoaXMgU2VjdGlvbiAyIGFyZSB0aGUgb25seSByaWdodHMgZ3JhbnRlZCB1bmRlcgogICAgIHRoaXMgTGljZW5zZS4gTm8gYWRkaXRpb25hbCByaWdodHMgb3IgbGljZW5zZXMgd2lsbCBiZSBpbXBsaWVkIGZyb20gdGhlCiAgICAgZGlzdHJpYnV0aW9uIG9yIGxpY2Vuc2luZyBvZiBDb3ZlcmVkIFNvZnR3YXJlIHVuZGVyIHRoaXMgTGljZW5zZS4KICAgICBOb3R3aXRoc3RhbmRpbmcgU2VjdGlvbiAyLjEoYikgYWJvdmUsIG5vIHBhdGVudCBsaWNlbnNlIGlzIGdyYW50ZWQgYnkgYQogICAgIENvbnRyaWJ1dG9yOgoKICAgICBhLiBmb3IgYW55IGNvZGUgdGhhdCBhIENvbnRyaWJ1dG9yIGhhcyByZW1vdmVkIGZyb20gQ292ZXJlZCBTb2Z0d2FyZTsgb3IKCiAgICAgYi4gZm9yIGluZnJpbmdlbWVudHMgY2F1c2VkIGJ5OiAoaSkgWW91ciBhbmQgYW55IG90aGVyIHRoaXJkIHBhcnR5J3MKICAgICAgICBtb2RpZmljYXRpb25zIG9mIENvdmVyZWQgU29mdHdhcmUsIG9yIChpaSkgdGhlIGNvbWJpbmF0aW9uIG9mIGl0cwogICAgICAgIENvbnRyaWJ1dGlvbnMgd2l0aCBvdGhlciBzb2Z0d2FyZSAoZXhjZXB0IGFzIHBhcnQgb2YgaXRzIENvbnRyaWJ1dG9yCiAgICAgICAgVmVyc2lvbik7IG9yCgogICAgIGMuIHVuZGVyIFBhdGVudCBDbGFpbXMgaW5mcmluZ2VkIGJ5IENvdmVyZWQgU29mdHdhcmUgaW4gdGhlIGFic2VuY2Ugb2YKICAgICAgICBpdHMgQ29udHJpYnV0aW9ucy4KCiAgICAgVGhpcyBMaWNlbnNlIGRvZXMgbm90IGdyYW50IGFueSByaWdodHMgaW4gdGhlIHRyYWRlbWFya3MsIHNlcnZpY2UgbWFya3MsCiAgICAgb3IgbG9nb3Mgb2YgYW55IENvbnRyaWJ1dG9yIChleGNlcHQgYXMgbWF5IGJlIG5lY2Vzc2FyeSB0byBjb21wbHkgd2l0aAogICAgIHRoZSBub3RpY2UgcmVxdWlyZW1lbnRzIGluIFNlY3Rpb24gMy40KS4KCjIuNC4gU3Vic2VxdWVudCBMaWNlbnNlcwoKICAgICBObyBDb250cmlidXRvciBtYWtlcyBhZGRpdGlvbmFsIGdyYW50cyBhcyBhIHJlc3VsdCBvZiBZb3VyIGNob2ljZSB0bwogICAgIGRpc3RyaWJ1dGUgdGhlIENvdmVyZWQgU29mdHdhcmUgdW5kZXIgYSBzdWJzZXF1ZW50IHZlcnNpb24gb2YgdGhpcwogICAgIExpY2Vuc2UgKHNlZSBTZWN0aW9uIDEwLjIpIG9yIHVuZGVyIHRoZSB0ZXJtcyBvZiBhIFNlY29uZGFyeSBMaWNlbnNlIChpZgogICAgIHBlcm1pdHRlZCB1bmRlciB0aGUgdGVybXMgb2YgU2VjdGlvbiAzLjMpLgoKMi41LiBSZXByZXNlbnRhdGlvbgoKICAgICBFYWNoIENvbnRyaWJ1dG9yIHJlcHJlc2VudHMgdGhhdCB0aGUgQ29udHJpYnV0b3IgYmVsaWV2ZXMgaXRzCiAgICAgQ29udHJpYnV0aW9ucyBhcmUgaXRzIG9yaWdpbmFsIGNyZWF0aW9uKHMpIG9yIGl0IGhhcyBzdWZmaWNpZW50IHJpZ2h0cyB0bwogICAgIGdyYW50IHRoZSByaWdodHMgdG8gaXRzIENvbnRyaWJ1dGlvbnMgY29udmV5ZWQgYnkgdGhpcyBMaWNlbnNlLgoKMi42LiBGYWlyIFVzZQoKICAgICBUaGlzIExpY2Vuc2UgaXMgbm90IGludGVuZGVkIHRvIGxpbWl0IGFueSByaWdodHMgWW91IGhhdmUgdW5kZXIKICAgICBhcHBsaWNhYmxlIGNvcHlyaWdodCBkb2N0cmluZXMgb2YgZmFpciB1c2UsIGZhaXIgZGVhbGluZywgb3Igb3RoZXIKICAgICBlcXVpdmFsZW50cy4KCjIuNy4gQ29uZGl0aW9ucwoKICAgICBTZWN0aW9ucyAzLjEsIDMuMiwgMy4zLCBhbmQgMy40IGFyZSBjb25kaXRpb25zIG9mIHRoZSBsaWNlbnNlcyBncmFudGVkIGluCiAgICAgU2VjdGlvbiAyLjEuCgoKMy4gUmVzcG9uc2liaWxpdGllcwoKMy4xLiBEaXN0cmlidXRpb24gb2YgU291cmNlIEZvcm0KCiAgICAgQWxsIGRpc3RyaWJ1dGlvbiBvZiBDb3ZlcmVkIFNvZnR3YXJlIGluIFNvdXJjZSBDb2RlIEZvcm0sIGluY2x1ZGluZyBhbnkKICAgICBNb2RpZmljYXRpb25zIHRoYXQgWW91IGNyZWF0ZSBvciB0byB3aGljaCBZb3UgY29udHJpYnV0ZSwgbXVzdCBiZSB1bmRlcgogICAgIHRoZSB0ZXJtcyBvZiB0aGlzIExpY2Vuc2UuIFlvdSBtdXN0IGluZm9ybSByZWNpcGllbnRzIHRoYXQgdGhlIFNvdXJjZQogICAgIENvZGUgRm9ybSBvZiB0aGUgQ292ZXJlZCBTb2Z0d2FyZSBpcyBnb3Zlcm5lZCBieSB0aGUgdGVybXMgb2YgdGhpcwogICAgIExpY2Vuc2UsIGFuZCBob3cgdGhleSBjYW4gb2J0YWluIGEgY29weSBvZiB0aGlzIExpY2Vuc2UuIFlvdSBtYXkgbm90CiAgICAgYXR0ZW1wdCB0byBhbHRlciBvciByZXN0cmljdCB0aGUgcmVjaXBpZW50cycgcmlnaHRzIGluIHRoZSBTb3VyY2UgQ29kZQogICAgIEZvcm0uCgozLjIuIERpc3RyaWJ1dGlvbiBvZiBFeGVjdXRhYmxlIEZvcm0KCiAgICAgSWYgWW91IGRpc3RyaWJ1dGUgQ292ZXJlZCBTb2Z0d2FyZSBpbiBFeGVjdXRhYmxlIEZvcm0gdGhlbjoKCiAgICAgYS4gc3VjaCBDb3ZlcmVkIFNvZnR3YXJlIG11c3QgYWxzbyBiZSBtYWRlIGF2YWlsYWJsZSBpbiBTb3VyY2UgQ29kZSBGb3JtLAogICAgICAgIGFzIGRlc2NyaWJlZCBpbiBTZWN0aW9uIDMuMSwgYW5kIFlvdSBtdXN0IGluZm9ybSByZWNpcGllbnRzIG9mIHRoZQogICAgICAgIEV4ZWN1dGFibGUgRm9ybSBob3cgdGhleSBjYW4gb2J0YWluIGEgY29weSBvZiBzdWNoIFNvdXJjZSBDb2RlIEZvcm0gYnkKICAgICAgICByZWFzb25hYmxlIG1lYW5zIGluIGEgdGltZWx5IG1hbm5lciwgYXQgYSBjaGFyZ2Ugbm8gbW9yZSB0aGFuIHRoZSBjb3N0CiAgICAgICAgb2YgZGlzdHJpYnV0aW9uIHRvIHRoZSByZWNpcGllbnQ7IGFuZAoKICAgICBiLiBZb3UgbWF5IGRpc3RyaWJ1dGUgc3VjaCBFeGVjdXRhYmxlIEZvcm0gdW5kZXIgdGhlIHRlcm1zIG9mIHRoaXMKICAgICAgICBMaWNlbnNlLCBvciBzdWJsaWNlbnNlIGl0IHVuZGVyIGRpZmZlcmVudCB0ZXJtcywgcHJvdmlkZWQgdGhhdCB0aGUKICAgICAgICBsaWNlbnNlIGZvciB0aGUgRXhlY3V0YWJsZSBGb3JtIGRvZXMgbm90IGF0dGVtcHQgdG8gbGltaXQgb3IgYWx0ZXIgdGhlCiAgICAgICAgcmVjaXBpZW50cycgcmlnaHRzIGluIHRoZSBTb3VyY2UgQ29kZSBGb3JtIHVuZGVyIHRoaXMgTGljZW5zZS4KCjMuMy4gRGlzdHJpYnV0aW9uIG9mIGEgTGFyZ2VyIFdvcmsKCiAgICAgWW91IG1heSBjcmVhdGUgYW5kIGRpc3RyaWJ1dGUgYSBMYXJnZXIgV29yayB1bmRlciB0ZXJtcyBvZiBZb3VyIGNob2ljZSwKICAgICBwcm92aWRlZCB0aGF0IFlvdSBhbHNvIGNvbXBseSB3aXRoIHRoZSByZXF1aXJlbWVudHMgb2YgdGhpcyBMaWNlbnNlIGZvcgogICAgIHRoZSBDb3ZlcmVkIFNvZnR3YXJlLiBJZiB0aGUgTGFyZ2VyIFdvcmsgaXMgYSBjb21iaW5hdGlvbiBvZiBDb3ZlcmVkCiAgICAgU29mdHdhcmUgd2l0aCBhIHdvcmsgZ292ZXJuZWQgYnkgb25lIG9yIG1vcmUgU2Vjb25kYXJ5IExpY2Vuc2VzLCBhbmQgdGhlCiAgICAgQ292ZXJlZCBTb2Z0d2FyZSBpcyBub3QgSW5jb21wYXRpYmxlIFdpdGggU2Vjb25kYXJ5IExpY2Vuc2VzLCB0aGlzCiAgICAgTGljZW5zZSBwZXJtaXRzIFlvdSB0byBhZGRpdGlvbmFsbHkgZGlzdHJpYnV0ZSBzdWNoIENvdmVyZWQgU29mdHdhcmUKICAgICB1bmRlciB0aGUgdGVybXMgb2Ygc3VjaCBTZWNvbmRhcnkgTGljZW5zZShzKSwgc28gdGhhdCB0aGUgcmVjaXBpZW50IG9mCiAgICAgdGhlIExhcmdlciBXb3JrIG1heSwgYXQgdGhlaXIgb3B0aW9uLCBmdXJ0aGVyIGRpc3RyaWJ1dGUgdGhlIENvdmVyZWQKICAgICBTb2Z0d2FyZSB1bmRlciB0aGUgdGVybXMgb2YgZWl0aGVyIHRoaXMgTGljZW5zZSBvciBzdWNoIFNlY29uZGFyeQogICAgIExpY2Vuc2UocykuCgozLjQuIE5vdGljZXMKCiAgICAgWW91IG1heSBub3QgcmVtb3ZlIG9yIGFsdGVyIHRoZSBzdWJzdGFuY2Ugb2YgYW55IGxpY2Vuc2Ugbm90aWNlcwogICAgIChpbmNsdWRpbmcgY29weXJpZ2h0IG5vdGljZXMsIHBhdGVudCBub3RpY2VzLCBkaXNjbGFpbWVycyBvZiB3YXJyYW50eSwgb3IKICAgICBsaW1pdGF0aW9ucyBvZiBsaWFiaWxpdHkpIGNvbnRhaW5lZCB3aXRoaW4gdGhlIFNvdXJjZSBDb2RlIEZvcm0gb2YgdGhlCiAgICAgQ292ZXJlZCBTb2Z0d2FyZSwgZXhjZXB0IHRoYXQgWW91IG1heSBhbHRlciBhbnkgbGljZW5zZSBub3RpY2VzIHRvIHRoZQogICAgIGV4dGVudCByZXF1aXJlZCB0byByZW1lZHkga25vd24gZmFjdHVhbCBpbmFjY3VyYWNpZXMuCgozLjUuIEFwcGxpY2F0aW9uIG9mIEFkZGl0aW9uYWwgVGVybXMKCiAgICAgWW91IG1heSBjaG9vc2UgdG8gb2ZmZXIsIGFuZCB0byBjaGFyZ2UgYSBmZWUgZm9yLCB3YXJyYW50eSwgc3VwcG9ydCwKICAgICBpbmRlbW5pdHkgb3IgbGlhYmlsaXR5IG9ibGlnYXRpb25zIHRvIG9uZSBvciBtb3JlIHJlY2lwaWVudHMgb2YgQ292ZXJlZAogICAgIFNvZnR3YXJlLiBIb3dldmVyLCBZb3UgbWF5IGRvIHNvIG9ubHkgb24gWW91ciBvd24gYmVoYWxmLCBhbmQgbm90IG9uCiAgICAgYmVoYWxmIG9mIGFueSBDb250cmlidXRvci4gWW91IG11c3QgbWFrZSBpdCBhYnNvbHV0ZWx5IGNsZWFyIHRoYXQgYW55CiAgICAgc3VjaCB3YXJyYW50eSwgc3VwcG9ydCwgaW5kZW1uaXR5LCBvciBsaWFiaWxpdHkgb2JsaWdhdGlvbiBpcyBvZmZlcmVkIGJ5CiAgICAgWW91IGFsb25lLCBhbmQgWW91IGhlcmVieSBhZ3JlZSB0byBpbmRlbW5pZnkgZXZlcnkgQ29udHJpYnV0b3IgZm9yIGFueQogICAgIGxpYWJpbGl0eSBpbmN1cnJlZCBieSBzdWNoIENvbnRyaWJ1dG9yIGFzIGEgcmVzdWx0IG9mIHdhcnJhbnR5LCBzdXBwb3J0LAogICAgIGluZGVtbml0eSBvciBsaWFiaWxpdHkgdGVybXMgWW91IG9mZmVyLiBZb3UgbWF5IGluY2x1ZGUgYWRkaXRpb25hbAogICAgIGRpc2NsYWltZXJzIG9mIHdhcnJhbnR5IGFuZCBsaW1pdGF0aW9ucyBvZiBsaWFiaWxpdHkgc3BlY2lmaWMgdG8gYW55CiAgICAganVyaXNkaWN0aW9uLgoKNC4gSW5hYmlsaXR5IHRvIENvbXBseSBEdWUgdG8gU3RhdHV0ZSBvciBSZWd1bGF0aW9uCgogICBJZiBpdCBpcyBpbXBvc3NpYmxlIGZvciBZb3UgdG8gY29tcGx5IHdpdGggYW55IG9mIHRoZSB0ZXJtcyBvZiB0aGlzIExpY2Vuc2UKICAgd2l0aCByZXNwZWN0IHRvIHNvbWUgb3IgYWxsIG9mIHRoZSBDb3ZlcmVkIFNvZnR3YXJlIGR1ZSB0byBzdGF0dXRlLAogICBqdWRpY2lhbCBvcmRlciwgb3IgcmVndWxhdGlvbiB0aGVuIFlvdSBtdXN0OiAoYSkgY29tcGx5IHdpdGggdGhlIHRlcm1zIG9mCiAgIHRoaXMgTGljZW5zZSB0byB0aGUgbWF4aW11bSBleHRlbnQgcG9zc2libGU7IGFuZCAoYikgZGVzY3JpYmUgdGhlCiAgIGxpbWl0YXRpb25zIGFuZCB0aGUgY29kZSB0aGV5IGFmZmVjdC4gU3VjaCBkZXNjcmlwdGlvbiBtdXN0IGJlIHBsYWNlZCBpbiBhCiAgIHRleHQgZmlsZSBpbmNsdWRlZCB3aXRoIGFsbCBkaXN0cmlidXRpb25zIG9mIHRoZSBDb3ZlcmVkIFNvZnR3YXJlIHVuZGVyCiAgIHRoaXMgTGljZW5zZS4gRXhjZXB0IHRvIHRoZSBleHRlbnQgcHJvaGliaXRlZCBieSBzdGF0dXRlIG9yIHJlZ3VsYXRpb24sCiAgIHN1Y2ggZGVzY3JpcHRpb24gbXVzdCBiZSBzdWZmaWNpZW50bHkgZGV0YWlsZWQgZm9yIGEgcmVjaXBpZW50IG9mIG9yZGluYXJ5CiAgIHNraWxsIHRvIGJlIGFibGUgdG8gdW5kZXJzdGFuZCBpdC4KCjUuIFRlcm1pbmF0aW9uCgo1LjEuIFRoZSByaWdodHMgZ3JhbnRlZCB1bmRlciB0aGlzIExpY2Vuc2Ugd2lsbCB0ZXJtaW5hdGUgYXV0b21hdGljYWxseSBpZiBZb3UKICAgICBmYWlsIHRvIGNvbXBseSB3aXRoIGFueSBvZiBpdHMgdGVybXMuIEhvd2V2ZXIsIGlmIFlvdSBiZWNvbWUgY29tcGxpYW50LAogICAgIHRoZW4gdGhlIHJpZ2h0cyBncmFudGVkIHVuZGVyIHRoaXMgTGljZW5zZSBmcm9tIGEgcGFydGljdWxhciBDb250cmlidXRvcgogICAgIGFyZSByZWluc3RhdGVkIChhKSBwcm92aXNpb25hbGx5LCB1bmxlc3MgYW5kIHVudGlsIHN1Y2ggQ29udHJpYnV0b3IKICAgICBleHBsaWNpdGx5IGFuZCBmaW5hbGx5IHRlcm1pbmF0ZXMgWW91ciBncmFudHMsIGFuZCAoYikgb24gYW4gb25nb2luZwogICAgIGJhc2lzLCBpZiBzdWNoIENvbnRyaWJ1dG9yIGZhaWxzIHRvIG5vdGlmeSBZb3Ugb2YgdGhlIG5vbi1jb21wbGlhbmNlIGJ5CiAgICAgc29tZSByZWFzb25hYmxlIG1lYW5zIHByaW9yIHRvIDYwIGRheXMgYWZ0ZXIgWW91IGhhdmUgY29tZSBiYWNrIGludG8KICAgICBjb21wbGlhbmNlLiBNb3Jlb3ZlciwgWW91ciBncmFudHMgZnJvbSBhIHBhcnRpY3VsYXIgQ29udHJpYnV0b3IgYXJlCiAgICAgcmVpbnN0YXRlZCBvbiBhbiBvbmdvaW5nIGJhc2lzIGlmIHN1Y2ggQ29udHJpYnV0b3Igbm90aWZpZXMgWW91IG9mIHRoZQogICAgIG5vbi1jb21wbGlhbmNlIGJ5IHNvbWUgcmVhc29uYWJsZSBtZWFucywgdGhpcyBpcyB0aGUgZmlyc3QgdGltZSBZb3UgaGF2ZQogICAgIHJlY2VpdmVkIG5vdGljZSBvZiBub24tY29tcGxpYW5jZSB3aXRoIHRoaXMgTGljZW5zZSBmcm9tIHN1Y2gKICAgICBDb250cmlidXRvciwgYW5kIFlvdSBiZWNvbWUgY29tcGxpYW50IHByaW9yIHRvIDMwIGRheXMgYWZ0ZXIgWW91ciByZWNlaXB0CiAgICAgb2YgdGhlIG5vdGljZS4KCjUuMi4gSWYgWW91IGluaXRpYXRlIGxpdGlnYXRpb24gYWdhaW5zdCBhbnkgZW50aXR5IGJ5IGFzc2VydGluZyBhIHBhdGVudAogICAgIGluZnJpbmdlbWVudCBjbGFpbSAoZXhjbHVkaW5nIGRlY2xhcmF0b3J5IGp1ZGdtZW50IGFjdGlvbnMsCiAgICAgY291bnRlci1jbGFpbXMsIGFuZCBjcm9zcy1jbGFpbXMpIGFsbGVnaW5nIHRoYXQgYSBDb250cmlidXRvciBWZXJzaW9uCiAgICAgZGlyZWN0bHkgb3IgaW5kaXJlY3RseSBpbmZyaW5nZXMgYW55IHBhdGVudCwgdGhlbiB0aGUgcmlnaHRzIGdyYW50ZWQgdG8KICAgICBZb3UgYnkgYW55IGFuZCBhbGwgQ29udHJpYnV0b3JzIGZvciB0aGUgQ292ZXJlZCBTb2Z0d2FyZSB1bmRlciBTZWN0aW9uCiAgICAgMi4xIG9mIHRoaXMgTGljZW5zZSBzaGFsbCB0ZXJtaW5hdGUuCgo1LjMuIEluIHRoZSBldmVudCBvZiB0ZXJtaW5hdGlvbiB1bmRlciBTZWN0aW9ucyA1LjEgb3IgNS4yIGFib3ZlLCBhbGwgZW5kIHVzZXIKICAgICBsaWNlbnNlIGFncmVlbWVudHMgKGV4Y2x1ZGluZyBkaXN0cmlidXRvcnMgYW5kIHJlc2VsbGVycykgd2hpY2ggaGF2ZSBiZWVuCiAgICAgdmFsaWRseSBncmFudGVkIGJ5IFlvdSBvciBZb3VyIGRpc3RyaWJ1dG9ycyB1bmRlciB0aGlzIExpY2Vuc2UgcHJpb3IgdG8KICAgICB0ZXJtaW5hdGlvbiBzaGFsbCBzdXJ2aXZlIHRlcm1pbmF0aW9uLgoKNi4gRGlzY2xhaW1lciBvZiBXYXJyYW50eQoKICAgQ292ZXJlZCBTb2Z0d2FyZSBpcyBwcm92aWRlZCB1bmRlciB0aGlzIExpY2Vuc2Ugb24gYW4gImFzIGlzIiBiYXNpcywKICAgd2l0aG91dCB3YXJyYW50eSBvZiBhbnkga2luZCwgZWl0aGVyIGV4cHJlc3NlZCwgaW1wbGllZCwgb3Igc3RhdHV0b3J5LAogICBpbmNsdWRpbmcsIHdpdGhvdXQgbGltaXRhdGlvbiwgd2FycmFudGllcyB0aGF0IHRoZSBDb3ZlcmVkIFNvZnR3YXJlIGlzIGZyZWUKICAgb2YgZGVmZWN0cywgbWVyY2hhbnRhYmxlLCBmaXQgZm9yIGEgcGFydGljdWxhciBwdXJwb3NlIG9yIG5vbi1pbmZyaW5naW5nLgogICBUaGUgZW50aXJlIHJpc2sgYXMgdG8gdGhlIHF1YWxpdHkgYW5kIHBlcmZvcm1hbmNlIG9mIHRoZSBDb3ZlcmVkIFNvZnR3YXJlCiAgIGlzIHdpdGggWW91LiBTaG91bGQgYW55IENvdmVyZWQgU29mdHdhcmUgcHJvdmUgZGVmZWN0aXZlIGluIGFueSByZXNwZWN0LAogICBZb3UgKG5vdCBhbnkgQ29udHJpYnV0b3IpIGFzc3VtZSB0aGUgY29zdCBvZiBhbnkgbmVjZXNzYXJ5IHNlcnZpY2luZywKICAgcmVwYWlyLCBvciBjb3JyZWN0aW9uLiBUaGlzIGRpc2NsYWltZXIgb2Ygd2FycmFudHkgY29uc3RpdHV0ZXMgYW4gZXNzZW50aWFsCiAgIHBhcnQgb2YgdGhpcyBMaWNlbnNlLiBObyB1c2Ugb2YgIGFueSBDb3ZlcmVkIFNvZnR3YXJlIGlzIGF1dGhvcml6ZWQgdW5kZXIKICAgdGhpcyBMaWNlbnNlIGV4Y2VwdCB1bmRlciB0aGlzIGRpc2NsYWltZXIuCgo3LiBMaW1pdGF0aW9uIG9mIExpYWJpbGl0eQoKICAgVW5kZXIgbm8gY2lyY3Vtc3RhbmNlcyBhbmQgdW5kZXIgbm8gbGVnYWwgdGhlb3J5LCB3aGV0aGVyIHRvcnQgKGluY2x1ZGluZwogICBuZWdsaWdlbmNlKSwgY29udHJhY3QsIG9yIG90aGVyd2lzZSwgc2hhbGwgYW55IENvbnRyaWJ1dG9yLCBvciBhbnlvbmUgd2hvCiAgIGRpc3RyaWJ1dGVzIENvdmVyZWQgU29mdHdhcmUgYXMgcGVybWl0dGVkIGFib3ZlLCBiZSBsaWFibGUgdG8gWW91IGZvciBhbnkKICAgZGlyZWN0LCBpbmRpcmVjdCwgc3BlY2lhbCwgaW5jaWRlbnRhbCwgb3IgY29uc2VxdWVudGlhbCBkYW1hZ2VzIG9mIGFueQogICBjaGFyYWN0ZXIgaW5jbHVkaW5nLCB3aXRob3V0IGxpbWl0YXRpb24sIGRhbWFnZXMgZm9yIGxvc3QgcHJvZml0cywgbG9zcyBvZgogICBnb29kd2lsbCwgd29yayBzdG9wcGFnZSwgY29tcHV0ZXIgZmFpbHVyZSBvciBtYWxmdW5jdGlvbiwgb3IgYW55IGFuZCBhbGwKICAgb3RoZXIgY29tbWVyY2lhbCBkYW1hZ2VzIG9yIGxvc3NlcywgZXZlbiBpZiBzdWNoIHBhcnR5IHNoYWxsIGhhdmUgYmVlbgogICBpbmZvcm1lZCBvZiB0aGUgcG9zc2liaWxpdHkgb2Ygc3VjaCBkYW1hZ2VzLiBUaGlzIGxpbWl0YXRpb24gb2YgbGlhYmlsaXR5CiAgIHNoYWxsIG5vdCBhcHBseSB0byBsaWFiaWxpdHkgZm9yIGRlYXRoIG9yIHBlcnNvbmFsIGluanVyeSByZXN1bHRpbmcgZnJvbQogICBzdWNoIHBhcnR5J3MgbmVnbGlnZW5jZSB0byB0aGUgZXh0ZW50IGFwcGxpY2FibGUgbGF3IHByb2hpYml0cyBzdWNoCiAgIGxpbWl0YXRpb24uIFNvbWUganVyaXNkaWN0aW9ucyBkbyBub3QgYWxsb3cgdGhlIGV4Y2x1c2lvbiBvciBsaW1pdGF0aW9uIG9mCiAgIGluY2lkZW50YWwgb3IgY29uc2VxdWVudGlhbCBkYW1hZ2VzLCBzbyB0aGlzIGV4Y2x1c2lvbiBhbmQgbGltaXRhdGlvbiBtYXkKICAgbm90IGFwcGx5IHRvIFlvdS4KCjguIExpdGlnYXRpb24KCiAgIEFueSBsaXRpZ2F0aW9uIHJlbGF0aW5nIHRvIHRoaXMgTGljZW5zZSBtYXkgYmUgYnJvdWdodCBvbmx5IGluIHRoZSBjb3VydHMKICAgb2YgYSBqdXJpc2RpY3Rpb24gd2hlcmUgdGhlIGRlZmVuZGFudCBtYWludGFpbnMgaXRzIHByaW5jaXBhbCBwbGFjZSBvZgogICBidXNpbmVzcyBhbmQgc3VjaCBsaXRpZ2F0aW9uIHNoYWxsIGJlIGdvdmVybmVkIGJ5IGxhd3Mgb2YgdGhhdAogICBqdXJpc2RpY3Rpb24sIHdpdGhvdXQgcmVmZXJlbmNlIHRvIGl0cyBjb25mbGljdC1vZi1sYXcgcHJvdmlzaW9ucy4gTm90aGluZwogICBpbiB0aGlzIFNlY3Rpb24gc2hhbGwgcHJldmVudCBhIHBhcnR5J3MgYWJpbGl0eSB0byBicmluZyBjcm9zcy1jbGFpbXMgb3IKICAgY291bnRlci1jbGFpbXMuCgo5LiBNaXNjZWxsYW5lb3VzCgogICBUaGlzIExpY2Vuc2UgcmVwcmVzZW50cyB0aGUgY29tcGxldGUgYWdyZWVtZW50IGNvbmNlcm5pbmcgdGhlIHN1YmplY3QKICAgbWF0dGVyIGhlcmVvZi4gSWYgYW55IHByb3Zpc2lvbiBvZiB0aGlzIExpY2Vuc2UgaXMgaGVsZCB0byBiZQogICB1bmVuZm9yY2VhYmxlLCBzdWNoIHByb3Zpc2lvbiBzaGFsbCBiZSByZWZvcm1lZCBvbmx5IHRvIHRoZSBleHRlbnQKICAgbmVjZXNzYXJ5IHRvIG1ha2UgaXQgZW5mb3JjZWFibGUuIEFueSBsYXcgb3IgcmVndWxhdGlvbiB3aGljaCBwcm92aWRlcyB0aGF0CiAgIHRoZSBsYW5ndWFnZSBvZiBhIGNvbnRyYWN0IHNoYWxsIGJlIGNvbnN0cnVlZCBhZ2FpbnN0IHRoZSBkcmFmdGVyIHNoYWxsIG5vdAogICBiZSB1c2VkIHRvIGNvbnN0cnVlIHRoaXMgTGljZW5zZSBhZ2FpbnN0IGEgQ29udHJpYnV0b3IuCgoKMTAuIFZlcnNpb25zIG9mIHRoZSBMaWNlbnNlCgoxMC4xLiBOZXcgVmVyc2lvbnMKCiAgICAgIE1vemlsbGEgRm91bmRhdGlvbiBpcyB0aGUgbGljZW5zZSBzdGV3YXJkLiBFeGNlcHQgYXMgcHJvdmlkZWQgaW4gU2VjdGlvbgogICAgICAxMC4zLCBubyBvbmUgb3RoZXIgdGhhbiB0aGUgbGljZW5zZSBzdGV3YXJkIGhhcyB0aGUgcmlnaHQgdG8gbW9kaWZ5IG9yCiAgICAgIHB1Ymxpc2ggbmV3IHZlcnNpb25zIG9mIHRoaXMgTGljZW5zZS4gRWFjaCB2ZXJzaW9uIHdpbGwgYmUgZ2l2ZW4gYQogICAgICBkaXN0aW5ndWlzaGluZyB2ZXJzaW9uIG51bWJlci4KCjEwLjIuIEVmZmVjdCBvZiBOZXcgVmVyc2lvbnMKCiAgICAgIFlvdSBtYXkgZGlzdHJpYnV0ZSB0aGUgQ292ZXJlZCBTb2Z0d2FyZSB1bmRlciB0aGUgdGVybXMgb2YgdGhlIHZlcnNpb24KICAgICAgb2YgdGhlIExpY2Vuc2UgdW5kZXIgd2hpY2ggWW91IG9yaWdpbmFsbHkgcmVjZWl2ZWQgdGhlIENvdmVyZWQgU29mdHdhcmUsCiAgICAgIG9yIHVuZGVyIHRoZSB0ZXJtcyBvZiBhbnkgc3Vic2VxdWVudCB2ZXJzaW9uIHB1Ymxpc2hlZCBieSB0aGUgbGljZW5zZQogICAgICBzdGV3YXJkLgoKMTAuMy4gTW9kaWZpZWQgVmVyc2lvbnMKCiAgICAgIElmIHlvdSBjcmVhdGUgc29mdHdhcmUgbm90IGdvdmVybmVkIGJ5IHRoaXMgTGljZW5zZSwgYW5kIHlvdSB3YW50IHRvCiAgICAgIGNyZWF0ZSBhIG5ldyBsaWNlbnNlIGZvciBzdWNoIHNvZnR3YXJlLCB5b3UgbWF5IGNyZWF0ZSBhbmQgdXNlIGEKICAgICAgbW9kaWZpZWQgdmVyc2lvbiBvZiB0aGlzIExpY2Vuc2UgaWYgeW91IHJlbmFtZSB0aGUgbGljZW5zZSBhbmQgcmVtb3ZlCiAgICAgIGFueSByZWZlcmVuY2VzIHRvIHRoZSBuYW1lIG9mIHRoZSBsaWNlbnNlIHN0ZXdhcmQgKGV4Y2VwdCB0byBub3RlIHRoYXQKICAgICAgc3VjaCBtb2RpZmllZCBsaWNlbnNlIGRpZmZlcnMgZnJvbSB0aGlzIExpY2Vuc2UpLgoKMTAuNC4gRGlzdHJpYnV0aW5nIFNvdXJjZSBDb2RlIEZvcm0gdGhhdCBpcyBJbmNvbXBhdGlibGUgV2l0aCBTZWNvbmRhcnkKICAgICAgTGljZW5zZXMgSWYgWW91IGNob29zZSB0byBkaXN0cmlidXRlIFNvdXJjZSBDb2RlIEZvcm0gdGhhdCBpcwogICAgICBJbmNvbXBhdGlibGUgV2l0aCBTZWNvbmRhcnkgTGljZW5zZXMgdW5kZXIgdGhlIHRlcm1zIG9mIHRoaXMgdmVyc2lvbiBvZgogICAgICB0aGUgTGljZW5zZSwgdGhlIG5vdGljZSBkZXNjcmliZWQgaW4gRXhoaWJpdCBCIG9mIHRoaXMgTGljZW5zZSBtdXN0IGJlCiAgICAgIGF0dGFjaGVkLgoKRXhoaWJpdCBBIC0gU291cmNlIENvZGUgRm9ybSBMaWNlbnNlIE5vdGljZQoKICAgICAgVGhpcyBTb3VyY2UgQ29kZSBGb3JtIGlzIHN1YmplY3QgdG8gdGhlCiAgICAgIHRlcm1zIG9mIHRoZSBNb3ppbGxhIFB1YmxpYyBMaWNlbnNlLCB2LgogICAgICAyLjAuIElmIGEgY29weSBvZiB0aGUgTVBMIHdhcyBub3QKICAgICAgZGlzdHJpYnV0ZWQgd2l0aCB0aGlzIGZpbGUsIFlvdSBjYW4KICAgICAgb2J0YWluIG9uZSBhdAogICAgICBodHRwOi8vbW96aWxsYS5vcmcvTVBMLzIuMC8uCgpJZiBpdCBpcyBub3QgcG9zc2libGUgb3IgZGVzaXJhYmxlIHRvIHB1dCB0aGUgbm90aWNlIGluIGEgcGFydGljdWxhciBmaWxlLAp0aGVuIFlvdSBtYXkgaW5jbHVkZSB0aGUgbm90aWNlIGluIGEgbG9jYXRpb24gKHN1Y2ggYXMgYSBMSUNFTlNFIGZpbGUgaW4gYQpyZWxldmFudCBkaXJlY3RvcnkpIHdoZXJlIGEgcmVjaXBpZW50IHdvdWxkIGJlIGxpa2VseSB0byBsb29rIGZvciBzdWNoIGEKbm90aWNlLgoKWW91IG1heSBhZGQgYWRkaXRpb25hbCBhY2N1cmF0ZSBub3RpY2VzIG9mIGNvcHlyaWdodCBvd25lcnNoaXAuCgpFeGhpYml0IEIgLSAiSW5jb21wYXRpYmxlIFdpdGggU2Vjb25kYXJ5IExpY2Vuc2VzIiBOb3RpY2UKCiAgICAgIFRoaXMgU291cmNlIENvZGUgRm9ybSBpcyAiSW5jb21wYXRpYmxlCiAgICAgIFdpdGggU2Vjb25kYXJ5IExpY2Vuc2VzIiwgYXMgZGVmaW5lZCBieQogICAgICB0aGUgTW96aWxsYSBQdWJsaWMgTGljZW5zZSwgdi4gMi4wLgoK"
    }
  ],
  "correspondence": {
    "schema": "elite-next-path-source-correspondence/v1",
    "status": "COMPLETE_MODULE_AST_UNDER_EXPLICIT_ADAPTERS",
    "source_sha256": "ef96862a543004139a10f8f52f0a45b4c972d0c4aa26e85396fd7aa2575d290e",
    "bundle_sha256": "ddc64218bc85fb88d28b5def06eb01fafb39cb67f7a9465ce736456658cc7f11",
    "bundle_region_sha256": "d8b58764b8a66de031a13ad8b76cccbb936e49f69c7df6d6bda306ced71fab83",
    "bundle_byte_start": 7619794,
    "bundle_byte_end": 7620528,
    "source_commit": "ce2c1386836339bc473b76c9888947899ead8d56",
    "parser": {
      "version": "8.14.0",
      "sha256": "758cead0e9764f94320f938ac169fb95c7eeea30dc675a6e5b1133fae239fa19"
    },
    "statements_compared": 4,
    "identifier_renames": {
      "path246": "path",
      "nextPath2": "nextPath",
      "from5": "from",
      "diff2": "diff",
      "next2": "next",
      "module2": "module",
      "__require": "require"
    },
    "other_adapters": [
      "unwrap one __commonJS module property; preserve four statements",
      "normalize only two named top-level var declarations to source const",
      "__require(path) to require(path) is explicit conditional loader adapter; host resolution not proven"
    ],
    "ast_sha256": "14befa9a0127d30d1d160aec651ab7ed32c7a6a669537c370a45f7cba13da2d7",
    "negative_probes": [
      {
        "label": "path builtin substitution",
        "rejected": true
      },
      {
        "label": "relative argument order",
        "rejected": true
      },
      {
        "label": "separator property",
        "rejected": true
      },
      {
        "label": "inclusive comparison",
        "rejected": true
      },
      {
        "label": "substring start",
        "rejected": true
      },
      {
        "label": "join argument order",
        "rejected": true
      },
      {
        "label": "return value",
        "rejected": true
      },
      {
        "label": "export binding",
        "rejected": true
      },
      {
        "label": "strict directive",
        "rejected": true
      },
      {
        "label": "top declaration binding",
        "rejected": true
      }
    ],
    "source_code_executed": false,
    "runtime_equivalence_claimed": false,
    "whole_pnpm_build_proven": false,
    "npm_tarball_byte_identity_proven": false
  }
}
````



## 6. Configuration surface

| Campo | Default | Validación | Secreto | Efecto |
|---|---|---|---|---|
| artifact | ninguno | tamaño/SHA256 exactos | no | lectura local |
| acquisition-receipt | ninguno | JSON estricto/bindings | no | referencia trazable |
| acquisition-receipt-sha256 | ninguno | hash esperado del workflow validado | no | rechaza drift |
| target | ausente obligatorio | Windows/sin reparse/sin sobrescritura | no | proyección candidata |
| selection-policy | fija, no argumento CLI | SHA256 embebido | no | única selección |

## 7. Dependency bill

CPython stdlib3.12+, PSF-2.0, runtime exacto y admisión del proyecto. pnpm11.25.0
es input opaco ya adquirido, no dependencia ejecutada ni código incluido en pack.
El JSON sólo contiene inventario y hashes. No nueva dependencia ni resolución.

## 8. Apply order

Validar perfil/recibo de adquisición; materializar el tooling en destino separado;
pasar17tests; proyectar archivo exacto a destino ausente; verificar y conservar
receipt. Before/after y rollback no alteran los originales. Cierre de runtime y
redistribución requieren sus gates propios, incluyendo consumers y licencias.

## 9. Verification

python -m unittest -v test_selection.py:17tests Windows, incluyendo negativas,
concurrencia y junction. Integración sobre archive oficial y receipt previo:
dos proyecciones idénticas,442payload files+receipt,22notices intactos, exclusión
exacta13files/6native; prueba de cambios, colisiones y política fija. Reconstruir
en directorio nuevo antes de incorporar tooling a un proyecto.

## 10. Reconstruction evidence

V334:17tests canónicos y17reconstruidos PASS;4/4fuentes idénticas. Dos
proyecciones del archivo oficial son byte-idénticas:442payload+receipt; entradas
originales intactas. Creación10.68/5.71s; verify0.39s. Colisión, execute y tamper
real rechazan. Evidencia reconstruction_evidence/PNPM_SELECTION_GATE_V334.md.
No producto, runtime, redistribución, release ni cierre global.

## 11. V335 consumer routing extension

Historical0.2.0 has6files and43regressions (17projection+26routing). plan_install.py
adds a separate no-execution prepare/verify command for3exact consumer profiles.
The6input paths and acquisition receipt are checked, Node bytes pinned, inherited
configuration excluded, and each recipe has an external expected SHA256 and
empty isolated configuration. Store and metadata cache are explicit. Existing
installs, ambient npmrc/hooks, drift and unsupported profiles are rejected.

All new code is AUTHORED. No runtime dependency, pnpm source byte, credential or
production effect is included. Windows trusted-directory conditions and runtime/
redistribution non-admission remain. Rebuild, adversarial and real consumer
results belong to reconstruction_evidence/PNPM_CONSUMER_ROUTING_V335.md.
The standalone profile changes1/4 to1/6; product profiles are unchanged.

## V336 — condición operacional de la cache

El protocolo de renovación y revalidación offline para los3consumers fijados está en reconstruction_evidence/PNPM_CACHE_FRESHNESS_V336.md. Una ruta de cache válida no demuestra frescura: no interpretar verifiedAt como TTL. Empezar en destino nuevo, preservar metadata oficial y comprobar sin verdicts copiados cuando se necesite replay offline. Mantener directorios confiables y hashes externos; sin garantía ante mutación concurrente hostil. Este addendum no cambia los6bloques0.2.0 ni admite runtime/redistribución.

## 12. V394 — current inputs and five-package supplemental notice

0.3.0 retains6materialized files.51tests (17projection +34routing). Current
enterprise-web0.5.16 and Playwright0.1.37 identities replace stale routing pins;
Lighthouse remains byte-identical. Immutable projection/source policy and
original22notices are unchanged. A v2 recipe adds BLUEOAK-NOTICE.md outside the
442-file payload, binds its bytes and five exact package manifests, and rejects
missing, modified or jointly rehashed notice/receipt. The Blue Oak official
link-form notice is delivered locally with the plan; downstream delivery and all
other licenses/source/publishing requirements remain separate. No fake upstream
license file, runtime admission, execution or redistribution authorization.
Evidence: reconstruction_evidence/PNPM_CURRENT_ROUTING_NOTICES_V394.md.

## 13. V395 — complete QRCode author and permission supplement

0.4.0 preserves6files;59tests (17projection +42routing). New v3 recipe delivers
QRCODE-NOTICE.md alongside the unchanged BlueOak notice, with exact author/
modification header and full MIT permission text. The plan_install.py block is ADAPTED: it combines AUTHORED control logic with
VERBATIM notice text. Its conceptual mixed origin is declared in the source field.
No QRCode algorithm or new runtime dependency is included in the pack.

Fixed-source index SHA2567377be90fc61a40268acf7f30d5bd89c2fca99c57ef5391623de8c151b8da7df,
header SHA256f265b9225bb2a1a209d60f81be487f23257c3f8637498282212d04af7822c2b3,
ten fixed source blobs and ten immutable selected bundle regions remain separate
evidence. Paths and byte ranges do not prove whole source/build equivalence.
No blanket Apache-to-MIT relabeling, execution or runtime/redistribution admission.
Evidence: reconstruction_evidence/QRCODE_VENDOR_NOTICE_DELIVERY_V395.md.

## V396 — original semver-utils license and local delivery

0.5.0/6files;67tests (17selection+50routing). v4 recipe delivers original LICENSE1839bytes under its complete MIT option in SEMVER-UTILS-NOTICE.md. Artifact SHA512, six-member inventory and exact LICENSE hash are preserved; APACHEv2 manifest/registry declaration remains unchanged and the signing key is expired. Original source and bundle region are fixed separately, not asserted build-equivalent. Retain BlueOak/QRCode notices and all442payloadfiles. No pnpm execution, dependency update or runtime/redistribution admission follows from these recipes. See reconstruction_evidence/SEMVER_ORIGINAL_LICENSE_DELIVERY_V396.md.

## V397 — retained notice collection delivery

0.6.0/7files;82tests (17selection+65routing). Recipev5 delivers PNPM-RETAINED-NOTICES.md containing47exact text copies/151253original bytes.22original notices are compared directly against selected payload;25research texts retain original scoped qualifications. Catalogue SHA256, individual digests/bytes, origin/schema/path/duplicate/read budgets and delivery hashes are checked. Canonical catalogue is ADAPTED because a portable JSON envelope carries unchanged third-party notices; LicenseRef-Pnpm-Bundled-Licenses preserves each work’s original terms and never re-licenses them. No runtime or redistribution admission, complete recursive SBOM, source offer or source/relinking fulfillment is inferred. Prior three supplements,442payload files and original receipts remain unchanged. Evidence reconstruction_evidence/PNPM_RETAINED_NOTICE_DELIVERY_V397.md.

## V398 — fixed next-path source and MPL delivery

0.7.0/8files;90tests (17selection+73routing). Recipev6 delivers NEXT-PATH-MPL-SOURCE.md: original index.js317bytes, package.json1010bytes and full MPL15978bytes, with source URLs/Git blobs/hashes and explicit bundler transformations. Four module statements compare structurally under named adapters;10mutations reject. Loader resolution, npm tarball identity and whole pnpm build are not proven. Source is carried as evidence data under MPL-2.0, never installed/executed. Previous47texts/three supplements and442payloadfiles remain unchanged. No blanket license/runtime/redistribution admission. Evidence reconstruction_evidence/NEXT_PATH_MPL_SOURCE_DELIVERY_V398.md.

## V401 — isolated containment implementation

## V401: published pnpm security block and isolated ZIP removal

The unmodified pnpm11.25.0 payload is affected by GHSA-vwc7-r8mq-g2x9.
Its historical scan does not authorize new execution. pnpm11.26.0 is also
affected. Version12.4.1 is a native candidate, not an admitted replacement.

contain_zip.py creates a separate ADAPTED candidate from all455 exact published
files. It retains441 selected files unchanged, changes only dist/pnpm.mjs,
removes all16 adm-zip modules and rejects binary ZIP acquisition before download
or extraction-directory creation. TAR support is unchanged. It retains original
notices and writes ADAPTATION.txt, adaptation.diff and a deterministic receipt.
Original input and existing targets are never overwritten. Trusted Windows
directories are required; concurrent hostile directory mutation is out of scope.

python contain_zip.py create --source <verified-published-pnpm-directory> --target <absent-candidate-directory>
python contain_zip.py verify --source <same-published-directory> --target <candidate-directory>

This utility performs no execution, installation, network call or admission.
The old plan_install.py recipes still bind the unmodified projection and must
not be redirected to the changed bundle. The candidate's bounded qualification
uses a separate environment and explicit digest. Its identity is the receipt
and bundle SHA, not the unchanged upstream version banner. Original licenses,
source/relinking obligations and general runtime/redistribution admission remain.
Binary ZIP runtime installation is intentionally unavailable. Reopen on an
admissible official fixed release, a new consumer or any change to the bundle.
Evidence: reconstruction_evidence/PNPM_ZIP_CONTAINMENT_V401.md.

### FILE: `pnpm_artifact_selection/contain_zip.py`
```yaml
block_id: "PNPM-ARTIFACT-SELECTION-GATE:contain_zip.py:v401"
operation: CREATE
provenance: AUTHORED
source: "local exact-byte containment under DEPENDENCY_UPDATE_CONTRACT; upstream source is not embedded in this tool"
license: "LicenseRef-Workspace-Owner"
sha256: "864d0c7d9931ac9e8b2d8c0d2b79d095002802b01bddf7c752cafe768e56c527"
variables: []
secrets_allowed: false
```
````python
"""Build an isolated pnpm11.25.0 candidate with binary ZIP support removed.

AUTHORED containment, not an official pnpm fix or general runtime admission.
Only exact published input bytes accepted. No network or execution operation.
"""
import argparse
import difflib
import os
from pathlib import Path
import re
import shutil
import sys
import tempfile
import select_artifact as base

BUNDLE = 'dist/pnpm.mjs'
ORIGINAL_SHA = 'ddc64218bc85fb88d28b5def06eb01fafb39cb67f7a9465ce736456658cc7f11'
NOTICE = b'''Local adaptation: Elite pnpm11.25.0 ZIP-disabled candidate, revision1.
The published pnpm artifact is retained separately. This is not an official release.
All16 bundled adm-zip0.6.0 modules have been removed. Binary ZIP fetch/extraction
fails closed before download or creation of the extraction directory. TAR package
installation is unchanged. Automatic runtime acquisition through ZIP is unavailable.
Reason: GHSA-vwc7-r8mq-g2x9. No upstream fixed release was identified2026-09-11.
Original license and notice files remain. Additional source/license/relinking and
runtime admission conditions remain open. No production or redistribution grant.
'''
DENIAL = 'throw new PnpmError("ZIP_DISABLED_BY_ELITE_CONTAINMENT", "Binary ZIP acquisition is disabled in this local pnpm adaptation; use a separately admitted runtime.");'

def transform(raw):
    base.require(base.digest(raw) == ORIGINAL_SHA, 'published bundle identity mismatch')
    text = raw.decode('utf-8')
    chunks = list(re.finditer(r'^// .*?/node_modules/adm-zip/[^\n]+\n', text, re.M))
    base.require(len(chunks) == 16, 'adm-zip module inventory changed')
    start = chunks[0].start()
    after = text.index('\n// ', chunks[-1].end()) + 1
    removed = text[start:after]
    symbols = re.findall(r'^var (\w+) = __commonJS\(', removed, re.M)
    base.require(len(symbols) == 16, 'adm-zip wrapper inventory changed')
    text = text[:start] + '// Elite containment: adm-zip implementation removed.\n' + text[after:]
    for name, next_name in [('downloadAndUnpackZip', 'downloadWithIntegrityCheck'), ('extractZipToTarget', 'toStatelessTester')]:
        pattern = r'async function ' + name + r'\([^\n]+\) \{\n[\s\S]*?\n\}\n(?=(?:async )?function ' + next_name + r'\()'
        signature = re.search(pattern, text)
        base.require(signature is not None, 'extraction function boundary changed')
        replacement = signature.group().split('\n', 1)[0] + '\n  ' + DENIAL + '\n}\n'
        text, count = re.subn(pattern, lambda _: replacement, text)
        base.require(count == 1, 'ambiguous extraction function')
    old = '      case "zip": {\n        const tempLocation = await cafs.tempDir();'
    base.require(text.count(old) == 1, 'binary ZIP dispatch boundary changed')
    # Reject before even creating the extraction directory.
    end = text.index('      default: {', text.index(old))
    text = text[:text.index(old)] + '      case "zip": {\n        ' + DENIAL + '\n      }\n' + text[end:]
    old_init = '    import_adm_zip = __toESM(require_adm_zip(), 1);\n'
    base.require(text.count(old_init) == 1, 'ZIP initialization changed')
    text = text.replace(old_init, '')
    base.require(text.count('var import_adm_zip, import_ssri3;') == 1, 'ZIP import binding changed')
    text = text.replace('var import_adm_zip, import_ssri3;', 'var import_ssri3;')
    for symbol in symbols:
        base.require(not re.search(r'\b' + re.escape(symbol) + r'\b', text), 'removed wrapper still referenced: ' + symbol)
    base.require('/node_modules/adm-zip/' not in text and 'import_adm_zip' not in text, 'ZIP implementation still present')
    return text.encode('utf-8')

def expected_files(source):
    source = base.no_reparse(source)
    policy = base.load_policy()
    entries = {x['path']: x for x in policy['files']}
    seen = set()
    for directory, dirs, files in os.walk(source, followlinks=False):
        for name in dirs: base.no_reparse(Path(directory) / name)
        for name in files:
            path = Path(directory) / name
            rel = path.relative_to(source).as_posix()
            base.require(rel in entries, 'unexpected source file')
            entry = entries[rel]
            raw = base.read_bounded(path, entry['bytes'])
            base.require(len(raw) == entry['bytes'] and base.digest(raw) == entry['sha256'], 'source file identity mismatch: ' + rel)
            seen.add(rel)
    base.require(seen == set(entries), 'published source incomplete')
    original = base.read_bounded(source / BUNDLE, entries[BUNDLE]['bytes'])
    changed = transform(original)
    result = {}
    for rel, entry in entries.items():
        if entry['selection'] == 'unchanged':
            result['payload/' + rel] = changed if rel == BUNDLE else base.read_bounded(source / rel, entry['bytes'])
    diff = ''.join(difflib.unified_diff(original.decode().splitlines(True), changed.decode().splitlines(True), fromfile='official/dist/pnpm.mjs', tofile='elite-zip-disabled/dist/pnpm.mjs'))
    result['adaptation.diff'] = diff.encode('utf-8')
    result['ADAPTATION.txt'] = NOTICE
    result['containment-receipt.json'] = base.json_bytes({
        'schema':'elite-pnpm-zip-disabled/v1', 'provenance':'ADAPTED_AUTHORED_CONTAINMENT',
        'base_policy_sha256':base.POLICY_SHA256, 'base_bundle_sha256':ORIGINAL_SHA,
        'candidate_bundle_sha256':base.digest(changed), 'removed_modules':16,
        'changed_payload_files':[BUNDLE], 'unchanged_payload_files':441,
        'excluded_physical_files':13, 'binary_zip_enabled':False,
        'upstream_advisory':'GHSA-vwc7-r8mq-g2x9', 'runtime_admitted':False,
        'redistribution_admitted':False, 'executed_by_materializer':False,
        'files':{p:{'sha256':base.digest(b),'bytes':len(b)} for p,b in sorted(result.items())}})
    return result

def verify(target, expected):
    target = base.no_reparse(target)
    found = set()
    allowed_dirs = {p.as_posix() for rel in expected for p in Path(rel).parents if p.as_posix() != '.'}
    for directory, dirs, files in os.walk(target, followlinks=False):
        for name in dirs:
            path = base.no_reparse(Path(directory) / name)
            base.require(path.relative_to(target).as_posix() in allowed_dirs, 'unexpected candidate directory')
        for name in files:
            path = Path(directory) / name
            rel = path.relative_to(target).as_posix()
            base.require(rel in expected, 'unexpected candidate file')
            base.require(base.read_bounded(path, len(expected[rel])) == expected[rel], 'candidate bytes changed: ' + rel)
            found.add(rel)
    base.require(found == set(expected), 'candidate files missing')

def create(source, target):
    base.require(os.name == 'nt', 'Windows publication required')
    target = base.no_reparse(target)
    base.require(target.parent.is_dir() and not os.path.lexists(target), 'target must be absent with existing parent')
    expected = expected_files(source)
    staging = Path(tempfile.mkdtemp(prefix='.pnpm-zip-disabled-', dir=target.parent))
    try:
        for rel, raw in expected.items():
            path = staging / rel
            path.parent.mkdir(parents=True, exist_ok=True)
            with path.open('xb') as stream: stream.write(raw)
        verify(staging, expected)
        base.no_reparse(target.parent)
        os.rename(staging, target)
    finally:
        if staging.exists():
            base.require(staging.parent == target.parent and staging.name.startswith('.pnpm-zip-disabled-'), 'cleanup scope rejected')
            shutil.rmtree(staging)
    verify(target, expected)

def main(argv=None):
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('action', choices=['create','verify'])
    parser.add_argument('--source', required=True)
    parser.add_argument('--target', required=True)
    args = parser.parse_args(argv)
    try:
        if args.action == 'create': create(args.source, args.target)
        else: verify(args.target, expected_files(args.source))
        print('PNPM_ZIP_CONTAINMENT_BYTES_PASS runtime_admitted=false redistribution_admitted=false')
        return 0
    except (base.SelectionError, OSError, ValueError, TypeError, KeyError) as exc:
        print('PNPM_ZIP_CONTAINMENT_BLOCKED: ' + str(exc), file=sys.stderr)
        return 2

if __name__ == '__main__': raise SystemExit(main())
````
