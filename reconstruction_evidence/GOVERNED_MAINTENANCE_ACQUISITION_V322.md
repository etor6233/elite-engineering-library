# Governed maintenance acquisition V322

Current result: G0–G8 narrow transport USE_REUSABLE_PACK; three official artifacts ACQUIRED with new governed receipts and byte-exact licenses, ten output files immutable on reentry. All151executed preflight steps PASS; library161/1453/757/52 and franchise67/746. Overall project remains DISCOVERY/BLOCKED, not100% or production ready.

Scope: library maintenance, T2803. Owner OFFICIAL_UPSTREAM_ACQUISITION_CORE0.4.78 remains REBUILD_VERIFIED / CONDITIONED globally. The audited instance adds an AUTHORED opaque transport and profile, with three ADAPTED base64 license envelopes preserving the exact original decoded bytes. It does not install or execute Node, pnpm or Sigstore, acquire GitHub source trees, attest signatures, run SCA, deploy or admit production.

## Before and correction

FAIL525: the root official-source record described V285 while V319–V321 had manual tool receipts. Those receipts remain historical; no approval is backdated. Node/pnpm fixed GitHub trees contain6/9symlinks and the current Windows ZIP branch cannot materialize them safely. The PyPI document SDK transport is incompatible with these artifact types. pnpm11.x documentation redirected12.x; fixed11.25 CLI help was used. Official release/signature/SCA evidence remains in V319–V321.

FAIL527: an empty legacy source directory plus a metadata-only receipt was accepted as PRESENT and the profile emitted ACQUIRED. Its Expect-Failure helper also accepted a successful action because it caught its own sentinel. Both red fixtures are preserved. Legacy PRESENT now fails closed and requires a new destination; no existing source is deleted or overwritten. The corrected oracle rejects actions that do not throw. Opaque PRESENT independently hashes the actual artifact and license and checks receipt/profile/approval/lock bindings. Aggregate receipts are immutable and source receipt hashes must agree.

FAIL528: the Node license could not roundtrip through the LF-only materializer. Integration stopped before writing the canonical pack. Three explicitly ADAPTED JSON/base64 envelopes now decode to the original bytes under a1MiB budget; a CRLF/no-final-newline fixture verifies preservation. Original license SHA values remain unchanged. The envelopes are not called VERBATIM source files.

## Audited contract and boundaries

Only exact Node Windows x64 executables and pnpm/Sigstore tarballs on their official hosts are eligible. Strict roster, release version, URL, filename, revision, size and digest validation precedes use. Sources and profile/approval hashes bind intent to receipts. HTTPS redirects/cookies/default credentials are disabled. Local imports require exact filename and bytes. The transport uses64KiB streaming buffers,128MiB artifact maximum,60s cancellation and bounded profile/lock/license/receipt reads. Reparse ancestors are rejected, including aggregate receipt paths. Outputs use CreateNew; failed partial files are preserved and block reentry. Resume uses the same inputs and rehashes output; a different approval or changed receipt fails. Runtime is a single local maintenance process; no multitenant service or untrusted concurrent local writer defense is claimed.

An interrupted write can leave a bounded partial artifact or incomplete receipt; the operator preserves evidence and uses a fresh destination. No retries, fallback to floating versions, automatic cleanup, source tree verification for legacy caches, or production rollback are claimed. Safe rollback is reverting the canonical owner/plan from the preserved before snapshot and not using the new transport. A vulnerable prior tool version is forensic evidence, not a safe runtime fallback.

## Gates for the narrow local adapter

| Gate | Evidence and limit |
|---|---|
| G0 | AUTHORED transport, exact canonical pack hash and reconstruction manifest. Official artifacts retain their own revisions/digests and V319–V321 signature evidence. |
| G1 | Workspace-owner orchestration; original Node/pnpm/Sigstore full licenses retained in reversible envelopes. No license claim from repository popularity. |
| G2 | PUBLIC_CODE_ARCHITECTURE_ADMISSION_STANDARD G0–G8, DEPENDENCY_UPDATE_CONTRACT and official-source profile contract govern acquisition. Readiness BLOCKED is preserved. |
| G3 | States absent→partial→verified artifact+license→receipt; immutable PRESENT or fail closed. CLI/lock/profile/approval and orchestration boundaries are materialized. |
| G4 | 17 valid profiles,7 profile negatives,9 lock negatives,65 transport assertions,37/37 candidate/rebuild byte equality. Empty-source and false-oracle regressions fixed. |
| G5 | Host/path/digest/type/size checks; tamper, incomplete cache, invalid approval, junction, redirect and cancellation tests. No secrets/private data. No SAST or full host security claim. |
| G6 | Five sequential16MiB SHA256 streams under a10s local budget each. Bounded streaming, injected stalled HTTP cancellation, partial-output and reentry tests. Memory benchmark excludes network/disk and is not a production SLO. |
| G7 | Immutable receipts, explicit failure messages, preserved red evidence and fail-closed migration from legacy cache; ownership and reopening triggers recorded. |
| G8 | Owner0.4.78,30 dependent plans updated,31 core files/37 initialization files. New suite added to executable preflight; complete structural verification pending below. |

## Local evidence inventory

Staging: `%LOCALAPPDATA%/Temp/elite-v322-a5c006a784f04fc08764f4a3853f17a5`. Synthetic test trees use `elite-opaque-test-*` under the same Temp root and are retained without recursive deletion of junction-bearing trees.

- `negative-oracle-red.log` SHA256 `3b4a45788ef3a9055d67bb64f97bbf5f4e89458d49c6474bcf2b2f9f692f16c9`.
- `empty-source-red.log` SHA256 `153b9e38f4bb865966d7c3f24713547c702917cdb51bededbed677c85fa0b2f7`.
- `rebuild-parity.json` SHA256 `881dd54f55ebab50b6e22e6fc4068c3460ef9bb13089ecd23538a6eb246c0c78`.
- `rebuilt-test_upstream_acquisition.log` SHA256 `65fc96b0158e5caf8168af58001d0cd7f02579f0b97076865792942e68fa19b0`.
- `rebuilt-test_source_profiles.log` SHA256 `8b946d97fdb80664adbb9b58fa7980bf4d02a5e336c42ba6b2608e69df9f1189`.
- `rebuilt-test_opaque_artifact_transport.log` SHA256 `9b68796fdb3cccf036de15e1ae64b33e9d053b7621d6df34caca217d721f6f5d`.
- `transport-benchmark.json` SHA256 `818b052f9626e8bad077d67f4659d01959552303aee634da69d6262f91fd54a3`.
- `research/tree-receipts.json` SHA256 `307f0bc3ef2607bfa1542b0a69ec7c5bb92b9066525129bde458a2dc6865aa80`.

Benchmark p50 14.240ms, nearest-rank p95/p99 42.112ms; peak process working set 119279616bytes. Five samples are descriptive only.

## Acquisition status before qualification

Not yet executed through the new materialized profile. FAIL525 remains OPEN until qualification and new receipts. Product readiness and all profile production blockers remain BLOCKED.

## Final governed result

Core0.4.78 reconstructed into a fresh destination before acquisition. Current approval was recorded by the agent under the existing user authorization for reversible library maintenance; it is not a user signature and is not backdated. All three exact required inputs are answered, hashes bind the profile and source lock, and blockers are acknowledged verbatim. apply_source_profile ValidateOnly passed, then actual official HTTPS acquired all three artifacts. A second invocation rehashed them and preserved all ten output files byte-for-byte. Read-only tar member inspection corroborated pnpm/Sigstore package names, versions and complete licenses; zero files were extracted or executed.

FAIL525 closes for this new maintenance acquisition only. V285 and manual V319–V321 history remain unchanged. FAIL527/528 close through canonical fixes and regression tests. Integrated failures529/530/531 were corrected through checkpoint refresh, notices inventory and six exact profile-count deltas, without weakening any checker. Logs are retained. Global owner remains CONDITIONED and each product/provider/runtime gate remains separate.

| Source | Bytes | Observed SHA256 | License SHA256 |
|---|---:|---|---|
| maintenance-node-24.20.0-win-x64 | 93381448 | `5c976096e04e5c2c1f091938926234cc9fbebfe9787ddd149351b3b0ecc707b5` | `ed34dd8e3f0a78dbaf00d0444ce8e285b015b765379c2e17880455f70370f8e9` |
| maintenance-pnpm-11.25.0 | 5116259 | `33dd0748f27e7916c4f1c8b6943461983e3453b06bbda6312a6280130b4881e5` | `e0a867ff513ea7be2a0ddc339ac6a031e459a38668e077b8f0e649544062f9f2` |
| maintenance-sigstore-4.1.1 | 10329 | `4d7ecc73cd9559457209adab0d9a64c50145e5cb1286de92abc75f0a140928a0` | `364a130d2ca340bd56eb1e6d045fc6929bb0f9d0aa018f2c1949b29517e1cdd0` |

The source lock additionally binds Sigstore with its official SHA512/SRI. A digest is an integrity binding, not a digital signature. Signature/provenance verification remains the separate V319/V321 evidence.

- `preflight75.log` SHA256 `b2fc30327b4d14baaf3b645541ad9c1ddd43e07d29c35114ee69482fc51f4c1f`.
- `preflight75.json` SHA256 `2b1cec85815d1d0ba9b8a7d2c983d6eb5eb65c7e0ffeab685914410d87993f76`.
- `maintenance-acquisition-gap-qualified.json` SHA256 `37d2dcfeac10b40918966049efd7db099d2843f8286080974ab6b10e54ec7158`.
- `maintenance-acquisition-gap-receipt.json` SHA256 `483319ede8d0745bc928cb42df9a70d516ef159ff2d80b68f261571e3229cb8a`.
- `maintenance-approval.json` SHA256 `213f0383d868b3ff0f409d3cd6a0175c01660aa3d122c604e69bb4e62471daa7`.
- `governed-validate.log` SHA256 `a3ca1f3b9edbdc48adccb78f55c2fff978df47b0511fff067e7409c6d0950820`.
- `governed-acquire.log` SHA256 `f5e5f8c20cbdc6d464159663c8121109b52edeb87ba7aac04454e9bcf7206d55`.
- `governed-present.log` SHA256 `86597d43bdbf2ffa8df7152ab2d8351bee7b57f6a24b3b406db422ea0e3b447a`.
- `governed-byte-verification.json` SHA256 `e5a2f9f6ecf1f09a018dc2e2ea2e0fd8b2564661efb33bc783836b15c89e2cb3`.
- `profile-count-delta.json` SHA256 `e82a7d4d22c26de461abfca99f530fa66f0b9ccb1415dbbe9e2bcf039fcd5180`.

## Portable receipt and input snapshot

```json
{
  "profile": {
    "schema": "elite-official-source-profile/v1",
    "profile_id": "maintenance-runtime-artifacts",
    "description": "Exact reviewed Node Windows executable and npm tool tarballs, transported opaquely with separate full license sidecars.",
    "source_ids": [
      "maintenance-node-24.20.0-win-x64",
      "maintenance-pnpm-11.25.0",
      "maintenance-sigstore-4.1.1"
    ],
    "required_user_inputs": [
      "Confirm this scope is local maintenance acquisition only, without extracting, installing or executing the artifacts.",
      "Identify the exact reviewed runtime/tool versions, immutable revisions, artifact digests and license evidence.",
      "Confirm the public-source network or exact local cache path, zero incremental spend, absence of private data and preservation of all production blockers."
    ],
    "production_blockers": [
      "An acquisition receipt proves only exact artifact and sidecar bytes; it is not source extraction, installation, runtime execution, SCA, signature verification or production admission.",
      "Current signature/provenance, full runtime and transitive native vulnerability coverage, license/redistribution obligations and platform-specific gates must be demonstrated before use or promotion.",
      "Target identities, providers, monitoring, deployment, recovery and human project readiness decisions remain unproven."
    ]
  },
  "profile_sha256": "7d7e5ad5347ca1feedf8c430ce16b1970423d60dd7276bc31011e0232f513d0b",
  "lock_sha256": "42fd2a1b8ad54d0e2111074b666290c7d9d27b39c3dca97897c184aabcd9947a",
  "approval": {
    "schema": "elite-project-official-source-profile-approval/v1",
    "profile_id": "maintenance-runtime-artifacts",
    "profile_sha256": "7d7e5ad5347ca1feedf8c430ce16b1970423d60dd7276bc31011e0232f513d0b",
    "lock_sha256": "42fd2a1b8ad54d0e2111074b666290c7d9d27b39c3dca97897c184aabcd9947a",
    "approved_by": "codex-agent-under-existing-user-maintenance-authorization",
    "approved_at": "2026-09-08T07:57:45.559157+00:00",
    "answers": [
      {
        "input": "Confirm this scope is local maintenance acquisition only, without extracting, installing or executing the artifacts.",
        "status": "ANSWERED",
        "response_or_reason": "Existing user instruction authorizes autonomous reversible library maintenance to continue the roadmap. This new acquisition is local opaque storage only, with no extraction, install or execution. It is an agent-recorded scope decision under existing authorization, not a new human signature or backdated approval."
      },
      {
        "input": "Identify the exact reviewed runtime/tool versions, immutable revisions, artifact digests and license evidence.",
        "status": "ANSWERED",
        "response_or_reason": "Exact reviewed artifacts: [{\"id\":\"maintenance-node-24.20.0-win-x64\",\"revision\":\"71b8b174857e25106d39b61a9e6f30d927da8b01\",\"digest\":{\"algorithm\":\"sha256\",\"value\":\"5c976096e04e5c2c1f091938926234cc9fbebfe9787ddd149351b3b0ecc707b5\"},\"license_sha256\":\"ed34dd8e3f0a78dbaf00d0444ce8e285b015b765379c2e17880455f70370f8e9\"},{\"id\":\"maintenance-pnpm-11.25.0\",\"revision\":\"6d90c71efdffbc909b499490b64c66badc720327\",\"digest\":{\"algorithm\":\"sha256\",\"value\":\"33dd0748f27e7916c4f1c8b6943461983e3453b06bbda6312a6280130b4881e5\"},\"license_sha256\":\"e0a867ff513ea7be2a0ddc339ac6a031e459a38668e077b8f0e649544062f9f2\"},{\"id\":\"maintenance-sigstore-4.1.1\",\"revision\":\"c1dc7d4778a450787fc72b083f2490ad02b714c6\",\"digest\":{\"algorithm\":\"sha512\",\"value\":\"7a776a1022647e1a33ad730ae6782efd4000d3155c54415d9c724212519a131ca98d6f872b98baceedcdb5e2e8697fe215b4546c2dfe0e3b6d42cd06011afee7\"},\"license_sha256\":\"364a130d2ca340bd56eb1e6d045fc6929bb0f9d0aa018f2c1949b29517e1cdd0\"}]. Authority and signature/SCA records are V319-V321; decoded license bytes and current materialized lock were verified in V322."
      },
      {
        "input": "Confirm the public-source network or exact local cache path, zero incremental spend, absence of private data and preservation of all production blockers.",
        "status": "ANSWERED",
        "response_or_reason": "Use public HTTPS nodejs.org and registry.npmjs.org only in a new isolated temporary destination. No private data, paid API, credentials, installed runtime changes or incremental purchase. Preserve all three profile production blockers exactly; acquisition does not close them."
      }
    ],
    "production_blockers_acknowledged": [
      "An acquisition receipt proves only exact artifact and sidecar bytes; it is not source extraction, installation, runtime execution, SCA, signature verification or production admission.",
      "Current signature/provenance, full runtime and transitive native vulnerability coverage, license/redistribution obligations and platform-specific gates must be demonstrated before use or promotion.",
      "Target identities, providers, monitoring, deployment, recovery and human project readiness decisions remain unproven."
    ]
  },
  "approval_sha256": "213f0383d868b3ff0f409d3cd6a0175c01660aa3d122c604e69bb4e62471daa7",
  "aggregate": {
    "schema": "elite-official-source-profile-receipt/v1",
    "profile_id": "maintenance-runtime-artifacts",
    "profile_sha256": "7d7e5ad5347ca1feedf8c430ce16b1970423d60dd7276bc31011e0232f513d0b",
    "lock_sha256": "42fd2a1b8ad54d0e2111074b666290c7d9d27b39c3dca97897c184aabcd9947a",
    "approval_sha256": "213f0383d868b3ff0f409d3cd6a0175c01660aa3d122c604e69bb4e62471daa7",
    "acquired_at": "2026-09-08T07:58:01.9819967+00:00",
    "source_receipts": [
      {
        "id": "maintenance-node-24.20.0-win-x64",
        "receipt": "receipts/maintenance-node-24.20.0-win-x64.json",
        "receipt_sha256": "071b0605aef67ea446a3e60feb8a4582cf684106e16102e861127bf1df6d49ab"
      },
      {
        "id": "maintenance-pnpm-11.25.0",
        "receipt": "receipts/maintenance-pnpm-11.25.0.json",
        "receipt_sha256": "ab4a8972ef1e1cd124345c2a774dc16d3039fec0a55332902abdf5e08995f513"
      },
      {
        "id": "maintenance-sigstore-4.1.1",
        "receipt": "receipts/maintenance-sigstore-4.1.1.json",
        "receipt_sha256": "b7a399e6e46266b5eab773d35f2f2af8ebdf4f1c02dfdbb71d6c44d27509e6ed"
      }
    ],
    "production_status": "BLOCKED_PENDING_USER_INPUTS_AND_LIVE_GATES"
  },
  "aggregate_sha256": "e7117557eda48f4662b1d01a741b9f4521b4fa2f203a800614a20504ff0b370e",
  "source_receipts": [
    {
      "schema": "elite-official-source-receipt/v1",
      "id": "maintenance-node-24.20.0-win-x64",
      "repository": "nodejs/node",
      "release": "v24.20.0",
      "commit": "71b8b174857e25106d39b61a9e6f30d927da8b01",
      "transport_kind": "opaque-release-artifact/v1",
      "artifact_url": "https://nodejs.org/dist/v24.20.0/win-x64/node.exe",
      "artifact_filename": "node-v24.20.0-win-x64.exe",
      "artifact_bytes": 93381448,
      "artifact_digest": {
        "algorithm": "sha256",
        "value": "5c976096e04e5c2c1f091938926234cc9fbebfe9787ddd149351b3b0ecc707b5"
      },
      "artifact_sha256": "5c976096e04e5c2c1f091938926234cc9fbebfe9787ddd149351b3b0ecc707b5",
      "license_expression": "MIT AND LicenseRef-Node-Bundled-Notices",
      "license_sha256": "ed34dd8e3f0a78dbaf00d0444ce8e285b015b765379c2e17880455f70370f8e9",
      "profile_sha256": "7d7e5ad5347ca1feedf8c430ce16b1970423d60dd7276bc31011e0232f513d0b",
      "approval_sha256": "213f0383d868b3ff0f409d3cd6a0175c01660aa3d122c604e69bb4e62471daa7",
      "lock_sha256": "42fd2a1b8ad54d0e2111074b666290c7d9d27b39c3dca97897c184aabcd9947a",
      "extracted": false,
      "installed": false,
      "executed": false,
      "production_admitted": false,
      "platform": "windows",
      "input_mode": "official_https",
      "acquired_at": "2026-09-08T07:58:01.3569244+00:00"
    },
    {
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
      "profile_sha256": "7d7e5ad5347ca1feedf8c430ce16b1970423d60dd7276bc31011e0232f513d0b",
      "approval_sha256": "213f0383d868b3ff0f409d3cd6a0175c01660aa3d122c604e69bb4e62471daa7",
      "lock_sha256": "42fd2a1b8ad54d0e2111074b666290c7d9d27b39c3dca97897c184aabcd9947a",
      "extracted": false,
      "installed": false,
      "executed": false,
      "production_admitted": false,
      "platform": "windows",
      "input_mode": "official_https",
      "acquired_at": "2026-09-08T07:58:01.7582562+00:00"
    },
    {
      "schema": "elite-official-source-receipt/v1",
      "id": "maintenance-sigstore-4.1.1",
      "repository": "sigstore/sigstore-js",
      "release": "npm4.1.1",
      "commit": "c1dc7d4778a450787fc72b083f2490ad02b714c6",
      "transport_kind": "opaque-release-artifact/v1",
      "artifact_url": "https://registry.npmjs.org/sigstore/-/sigstore-4.1.1.tgz",
      "artifact_filename": "sigstore-4.1.1.tgz",
      "artifact_bytes": 10329,
      "artifact_digest": {
        "algorithm": "sha512",
        "value": "7a776a1022647e1a33ad730ae6782efd4000d3155c54415d9c724212519a131ca98d6f872b98baceedcdb5e2e8697fe215b4546c2dfe0e3b6d42cd06011afee7"
      },
      "artifact_sha256": "4d7ecc73cd9559457209adab0d9a64c50145e5cb1286de92abc75f0a140928a0",
      "license_expression": "Apache-2.0",
      "license_sha256": "364a130d2ca340bd56eb1e6d045fc6929bb0f9d0aa018f2c1949b29517e1cdd0",
      "profile_sha256": "7d7e5ad5347ca1feedf8c430ce16b1970423d60dd7276bc31011e0232f513d0b",
      "approval_sha256": "213f0383d868b3ff0f409d3cd6a0175c01660aa3d122c604e69bb4e62471daa7",
      "lock_sha256": "42fd2a1b8ad54d0e2111074b666290c7d9d27b39c3dca97897c184aabcd9947a",
      "extracted": false,
      "installed": false,
      "executed": false,
      "production_admitted": false,
      "platform": "windows",
      "input_mode": "official_https",
      "acquired_at": "2026-09-08T07:58:01.9510456+00:00"
    }
  ],
  "gap_receipt": {
    "authority_count": 2,
    "candidate_count": 1,
    "capability_id": "GOVERNED-MAINTENANCE-ARTIFACT-ACQUISITION",
    "decision_state": "USE_REUSABLE_PACK",
    "gate": "CAPABILITY-GAP-RESOLUTION-GATE",
    "implementation_ready": true,
    "observed_at": "2026-09-08T07:28:10.571919Z",
    "official_domain_count": 5,
    "outcome": "READY",
    "record_sha256": "37d2dcfeac10b40918966049efd7db099d2843f8286080974ab6b10e54ec7158",
    "requirement_ref": "research/PROJECT_BLUEPRINT.md",
    "schema_version": 1,
    "search_count": 1,
    "selected_candidate_id": "local-opaque-maintenance-transport-0.4.78"
  }
}
```

## Continuation

Continue T2803 native/transitive tool coverage and T2809 host/supervisor/identity/reporting/retention/alerts, using existing owners. pnpm native inventory has four reflink0.1.19 addons and two fastlist0.3.0 executables; six unchanged binaries are not covered by the473npm-package result. Read-only research stagingV323 is prepared. D/channel/commercial/target decisions remain pending; ARCA deferred. No ZIP, deployment, production promotion or scheduled monitor created.
