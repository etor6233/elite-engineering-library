# Sharp native build-input correspondence — V384

2026-09-11. Library maintenance / TEST03, parent V383 checkpoint229 validated, VERIFY_LIBRARY
passed. No new product requirement, private data, runtime substitution or production operation.

V384 demuestra correspondencia exacta del insumo nativo Sharp: libvips-42.dll18614784bytes, versions.json y avisos coinciden con @img/sharp-libvips-win32-x641.3.3 firmado. Perfil previo/2attestations/tamper, cuarentena sin ejecución y4suites del transporte PASS; core0.4.85/48files/194IDs/23perfiles. DLL C++ y build MXE/notices/SCA completos pendientes. Runtime y perfil68/795 intactos;45/48 sin promoción.

## Fixed chain and verified bytes

Sharp0.35.4 at verified7f1a0a22cc285fe180766f4935d50b55af6e8432 declares the Windows x64
development input @img/sharp-libvips-win32-x641.3.3. Its packaging script copies the generated
binary output plus versions.json and the third-party README suffix from that input. Source
inspection fixes package.json, build scripts and packaging script against their Git blob hashes;
the already verified CI workflow runs build, unit tests and package-from-local-build for Windows.
[Pinned Sharp packaging source](https://github.com/lovell/sharp/blob/7f1a0a22cc285fe180766f4935d50b55af6e8432/npm/from-local-build.js).

The1.3.3input npm metadata, registry signature and both npm-publish/SLSA attestations are verified
before acquisition. Altered registry message and DSSE payload are rejected. SLSA binds exact
package/version/SHA512 to lovell/sharp-libvips6e5971d333377743163edc3ad9e5d0b897abcbc9 and
.github/workflows/ci.yml@refs/tags/v1.3.3. This proves package provenance, not a complete native
build SBOM or identity of every downloaded build dependency.

An exact source profile and answers were materialized before acquisition. The authorized action
is delegated reversible library inspection; no human signature, external expenditure or product
authority was inferred. The transport adds only the exact reviewed package/version/commit/URL
roster entry, packaging-license envelope and profile; all four source/profile/opaque/wheel suites
pass after48-file reconstruction. The existing opaque tests exercise exact bytes and reject
changed version, commit, repository, filename, query/host/scheme and unreviewed package.
ValidateOnly and the actual acquisition both pass;8,263,820archive bytes remain in quarantine.

Bounded read-only tar inspection rechecks SHA512, rejects unsafe/duplicate/link paths, enforces
member/total sizes and inspects seven files without installation or binary execution. The
18,614,784-byte lib/libvips-42.dll equals the already authenticated installed Sharp DLL exactly:
SHA256 e6cc51bbc763e7deda536c6f56ce96b4c51ea769690ce4f4ed07607527b81dae.
versions.json and the complete third-party README suffix are also byte-identical. This closes
the upstream-input→installed-libvips-DLL/version/notice correspondence. The installed
libvips-cpp-8.18.6.dll is not supplied by this package; it belongs to the separate Sharp C++
wrapper build and remains explicit. No equivalence of that DLL or sharp.node is inferred.

## Windows recipe and actual remaining origin

The fixed sharp-libvips Windows script does not compile the native stack from versions.properties.
It downloads vips-dev-x64-web-8.18.6-static.zip from libvips/build-win64-mxe releasev8.18.6,
copies its DLLs and packages its JSON metadata. The global AOM3.15.0 declaration therefore does
not establish a Windows version mismatch against the observed3.14.1. The exact Windows zip,
its MXE recipe/toolchains/component sources and DLL correspondence are the next evidence chain.
The script also fetches third-party notices from a moving main URL. The signed archive fixes
the actual delivered notice bytes; a pinned recipe alone does not make that input reproducible.
[Pinned Windows input recipe](https://github.com/lovell/sharp-libvips/blob/6e5971d333377743163edc3ad9e5d0b897abcbc9/build/win.sh).

V383 already identified libnsgif's13vendored blobs at libvips426af3f44246fce9cfa8dd51a353aa4dfd48c553.
That source identity is reusable, but source→MXE DLL correspondence still requires proof. No
NetSurf semantic version was invented. The29license rows/28version entries remain explicit.
The source repository Apache-2.0 license concerns packaging scripts; it does not replace the
LGPL-3.0-or-later native payload or its component notices, source/relinking and redistribution
obligations. The sidecar's canonical provenance correctly labels only its Apache bytes. Complete
SCA/native/CRT/stdlib/build/notices remain gated, as do existing pnpm native FAIL532 and TEST03/07.

## Canonical changes and preservation

OFFICIAL-UPSTREAM-ACQUISITION-CORE0.4.85 adds two files (one AUTHORED profile, one ADAPTED
reversible packaging-license envelope) and updates five existing transport/lock/test files.
30actual source profiles select that owner; six hardcoded composition budgets increase by2.
The franchise and HTTP runtime profiles do not select acquisition tooling and remain exactly
68/795 and69/803, including their plan and source-lock hashes. No BFF/Go/SQL/runtime source or
dependency changed; V383 execution results are referenced through byte parity, not rerun counts.
Inventory165packs/1557files/832Markdown/56plans;1304AUTHORED/146ADAPTED/107VERBATIM.

FAIL760 records the original incorrect assumption that the source-tooling owner belonged to
the runtime parent. That failed pending preparation published no canonical changes. Original
pending/helper evidence remains, and the corrected preparation derives consumers from actual
plan selections. Repeated guessed-path tooling incidents are retained under FAIL719; resolve
current artifacts from the recorded inventory instead of historical stage-name assumptions.
No acceptance definitions, actions, oracles or statuses change:45passed, TEST02/03/07blocked.
This qualification closes one source/input boundary, not complete security or functional readiness.

## Receipts

Raw public inspection and quarantine evidence resolve through elite-v384-current.txt. No new
native binary is placed in the library or consumer runtime.

| Evidence | SHA-256 |
|---|---|
| inspection-receipts.json | 2e9c730c9ca26ef9be7b6259c03b94a2370379ef43f7602b73334ec4fa4a055d |
| artifact-signature-receipt.json | 5644981ac17a97634956ef9659de50f8be9c86a9ea996bd98d78fedd7545616b |
| artifact-metadata-inspection.json | a0aa2e826c4f5fe08e687df8f3a0f5760fb9f7ed9fe3779bbbba7715a7fe26c9 |
| source-candidate-changes.json | 7ee9637e8a1ed91594f8c9f2dd49e2a1aeeabb1c7cf18d6892b7e0ab66a99cad |
| source-core-test-results.json | 653cf0ed4951d1c5df86cf59b80cc88c6876cd25e5bf8e13451e83adb441212a |
| native-source-approval.json | 7285279b5de47518414161a6182809c2c299915a5a601f6342e98f95c6d3ee4f |
| native-source-validate-process.json | 6c7eeec9cb8df6fcfa07095adc253bb25f70c170009709e56e885b6dac9e2f5e |
| native-source-acquire-process.json | 2f596cf85c157fb004e85fc175ffcdb590590f009ad9c213714c18b90a04829a |
| native-input-correspondence.json | bb5cd216744ec06ef54e24d72d5e4fc033eac6337d22ff46e1c2d41611fb1139 |
| canonical-pending-changes.json | 338b294f02ea5ec6b690606c830d0e040d8b4aef9c3059194cd25de52530727f |
| quarantined-native-input/profile-receipts/maintenance-sharp-native-input-inspection.json | 43ac2e7612beb8bc3b534a63dbf83c5806b86b4f8fec4fe9d63299c74c845999 |

## Final canonical reconstruction

The canonical source core reconstructs48files byte-identical to the qualified candidate. A fresh HTTP consumer reconstructs69packs/803files, every file identical to V383; ordinary franchise remains68/795. No runtime rebuild/browser run is claimed for these unchanged bytes. Full Preflight230 follows against a stable root.

## Closure231

V384 cierre231: input nativo Sharp1.3.3 firmado y adquirido con perfil; libvips-42.dll18614784bytes, versiones y avisos byte-idénticos. Core0.4.85/48files/194IDs/23perfiles y4suites PASS;30planes actualizados. Preflight230164pasos/56composiciones PASS, Docker ausente. Runtime68/795 y69/803 idéntico aV383.45/48; integración completa, seguridad nativa y release pendientes.

```json
{
  "checkpoint_under_test": 230,
  "executed_steps": 164,
  "all_executed_steps": "PASS",
  "compositions": 56,
  "missing_toolchains": [
    "docker"
  ],
  "inventory": {
    "packs": 165,
    "files": 1557,
    "markdown": 832,
    "plans": 56
  },
  "source_core": {
    "version": "0.4.85",
    "files": 48,
    "source_ids": 194,
    "source_profiles": 23,
    "source_suites_passed": 4,
    "affected_compositions": 30
  },
  "franchise": {
    "packs": 68,
    "files": 795
  },
  "http_reference": {
    "packs": 69,
    "files": 803,
    "all_files_identical_to_v383": true
  },
  "matched_native_dll_bytes": 18614784,
  "matched_native_dll_sha256": "e6cc51bbc763e7deda536c6f56ce96b4c51ea769690ce4f4ed07607527b81dae",
  "runtime_substitution_or_execution": false,
  "controls": {
    "passed": 45,
    "total": 48,
    "acceptance_definitions_actions_oracles_statuses_unchanged": true
  },
  "preflight_log_sha256": "71260eaa89416201188f77b7011fb7b9a12c9b1f177f57793d4ecbc5f88e183a",
  "preflight_json_sha256": "f7ab38bfca9a013b121048c450fdfd4112dd2011da67a450b3c2e3ff50a08c8f",
  "final_parity_sha256": "6af0c5dd0fb40a17a4bbf913f6b087d87f5fc0c48c0ca00e9b7f937b420e8773",
  "input_correspondence_sha256": "bb5cd216744ec06ef54e24d72d5e4fc033eac6337d22ff46e1c2d41611fb1139"
}
```

## Next exact MXE input identified, not acquired

The official release metadata identifies vips-dev-x64-web-8.18.6-static.zip, 9322134bytes, sha256:42ceea0c2f53d1244f3cb54c059881487b37452ee46c1c98c0d4b47480b20588. The tag resolves to 09cfccf20b91b441fbe97fa7a7ed8a597e55e830; GitHub reports commit verification=false, which is preserved rather than called signed. Twelve listed release assets are binary ZIP variants, not a separately listed source bundle. No MXE binary was acquired. Before the next acquisition, materialize the exact source profile/lock and preserve the distinction between the platform release digest, source identity, package attestation and reproducible build. Metadata receipts: mxe-release-discovery.json SHA256 0e5804a98798bd4aae3872364b114ca3e77237924deeb95668b129e1f231b676 and mxe-recipe-discovery.json SHA256 01a07ed03e6a41b167f04a242787d26ae867936c5e0ac7f918e54c6db1db4556.

[Official MXE release](https://github.com/libvips/build-win64-mxe/releases/tag/v8.18.6).
