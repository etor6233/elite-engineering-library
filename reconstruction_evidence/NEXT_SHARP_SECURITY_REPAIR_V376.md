# Next and Sharp security repair — V376

Date: 2026-09-10. Local library maintenance;45/48 controls remain closed. TEST02/03/07 and FAIL385 are still open.

V376: Next16.3.4/Sharp0.35.4 replaces the affected Next16.3.2/Sharp0.35.3 baseline in TS-GO-API-WEB-BRIDGE0.5.5. Five exact public artifacts, ten verified provenance attestations, 8589 installed files identical to inspected archives,152 declared npm versions/0OSV findings versus3baseline advisories,115web tests plus1explicit integration skip, typecheck/build,92connected browser phases and4agenda projects with53exact migrations PASS. Multitab/session gates enabled on all four browser projects. Native AVIF/PNG/WebP smoke confirms libheif1.23.2. This is scoped advisory remediation; full vendored/native security and redistribution obligations remain blocked. Acquisition core0.4.84/46files,193IDs/22profiles and all four transport suites PASS;30dependent source plans updated. HTTP reference0.1.1 records the changed web-parent plan; Go/SQL/runtime evidence is reused only with byte parity.45/48 unchanged. See reconstruction_evidence/NEXT_SHARP_SECURITY_REPAIR_V376.md. Final canonical composition, HTTP binary parity and full Preflight remain pending in this entry.

## Official authority and exact change

- [Next Windows RCE advisory](https://github.com/vercel/next.js/security/advisories/GHSA-p293-qw3h-jr36): affected16.x below16.3.3.
- [Next AVIF advisory](https://github.com/vercel/next.js/security/advisories/GHSA-2xp9-vwfh-vxw4).
- [Sharp libheif advisory](https://github.com/lovell/sharp/security/advisories/GHSA-rgj7-g3m4-5g8c).
- [Next16.3.4 release](https://github.com/vercel/next.js/releases/tag/v16.3.4) restores AVIF using patched Sharp0.35.4.

Next repository299180d3315c7ebd7b199d2b1a265b5986c5fc7d; Sharp7f1a0a22cc285fe180766f4935d50b55af6e8432. GitHub-verified commits, registry signatures and ten Sigstore publish/SLSA attestations bound to exact package/version/SHA512/repository/workflow. Altered signatures/payloads rejected. The signature receipt predates acquisition; subsequent quarantine and installed-byte receipts bind those authenticated digests to the actual bytes.

## Evidence boundaries

152exact npm version queries cover the declared lock graph. They do not cover every vendored package, native compiler, DLL or WASM component. The artifact inventory records146nested compiled package manifests (only10versioned) and native file hashes; versions.json is a component list, not a complete SBOM. Benign codec tests do not reproduce the exploits. Current known package advisories are repaired; do not claim complete native security.

@next/env and SWC omit archive-local license files: redistribute the exact Next root MIT companion license. Native Sharp declares Apache-2.0 AND LGPL-3.0-or-later; its root Apache file does not resolve every native dependency notice/source obligation. Keep that redistribution gate open. No upstream runtime binaries are included in the library Markdown source. Existing pnpm native license/build/SBOM/CRT blocker FAIL532 also remains.

The legacy unselected TypeScript golden-path candidate retains its historical quarantined dependencies; this update applies to the four active plans selecting the BFF. No product, provider, data, deployment or paid service authority is inherited.

## Actual connected execution

Fresh owned loopback PostgreSQL18.6; ownership verified by data_directory before database creation;53migration hashes identical to V372. Go1.26.8 and Node24.20.0 exact; four browser projects,92phases, four agendas, return multitab/session assertions enabled. PostgreSQL stopped after completion. The115/1web result preserves the explicitly skipped connected Vitest case; browser/PostgreSQL suites ran separately without skips.

## Reproduction and receipt index

Materialize FRANCHISE_COMPLETE_PACK_PLAN in an absent directory; install the exact frozen lock with the pinned Node/pnpm and ignore-scripts. Validate all installed patched artifacts against signed quarantined bytes. Run pnpm test, typecheck and build. Install the isolated browser gate with --ignore-workspace --offline --frozen-lockfile --ignore-scripts. Use a newly owned elite_confirmation_* loopback database with the exact migrations; enable all18flags in the preserved finite connected harness, execute TestQuoteAcceptanceBrowserPostgres and TestAppointmentAgendaBrowserPostgres, require92/4 and no skip/failure, then stop the owned database. Never substitute a project database or copy synthetic identity into production.

Raw files are retained under the external V376 evidence stage, resolved locally by the elite-v376-current.txt temporary pointer; machine paths are omitted from distribution. SHA256 below binds original raw bytes, not normalized display text.

| Relative evidence | SHA256 |
|---|---|
| artifact-signature-receipt.json | 634523a17c9c254bd045003a948581fb3fef045b88e8652a3ad79bf1449c5a0c |
| artifact-metadata-inspection.json | 672a5dc85b7e9aef152d37e8f347032b239feeb1213b4ad40ba3451470e61d40 |
| artifact-file-inventories.json | 42a27e8f69a7098b29cf53c1b043f607f3303f191d1f30ae5ba7b5d1e6def583 |
| artifact-inspection-summary.json | 4c9380dbc3374a740d7cccfe869ce6f9811fa156ceb6560286bef6ecadcdd58e |
| installed-artifact-parity.json | 31aa723a4c7a28ed9429f15f6fd615f0f806bee31b74ef6a788f7d47f619a8e6 |
| baseline-osv-request.json | 1f4595b71efaf8fcaceb5eb7e4e3141f25c522bbd99ba720967dd4e71fee3d64 |
| baseline-osv-response.json | 5854f26f16e092f84bd2777a06d96b82b50642b9d99cb3ddf1cefb543d68821c |
| candidate-osv-request.json | 2f512e198893adf5192b1ea42930ad1725221d3be05b1e98120d9999239efe27 |
| candidate-osv-response.json | f2c4ed7690153f9b7d6ac141b6dfcdebb80b8e045836d60a5485f37a09c947ed |
| candidate-web-tests.log | 640550ff4b5a15c3c696c3b0c691b10528d905c7209620a98d3118ea66458df6 |
| candidate-web-typecheck.log | 99d7c2179cdd316467c51f37901025a0753453664d741214ba56fc631fa6c0bc |
| candidate-web-build.log | e8134e2604357569363d5e02e6d278abf9810712bf4cd5381cd38ee56dbf66d1 |
| sharp-native-probe.json | 032e819a5bbb6b8c6e0a44c68863d1e43836bce66851ac7d1e5d8d33b58edf2b |
| source-core-test-results.json | 23dd11578b8e16b819ef2b732cb78915f1b0b8cbdde2c2166c1f37e6b396e3c1 |
| pending-block-parity.json | dfd664ec42e298141bd33c1b038460cade049370be405fdabfe9d36a24b95578 |
| canonical-pending-changes.json | 5643435777df279e6f3e43ba14368111d4daedf15c2fcfe8e54da9a0f17676f6 |
| connected-web376.py | 201166c2e218e0bd7a2582e60b3908ed842be1d20dbb51e5df0f9a8422ae34c3 |
| verify-next-attestations.cjs | 8319a5c4236fc4a69cf71fef595dfb458b69034c550289c1e57dc442fb370108 |
| inspect-web-tarballs.py | 2e93065616d01cf8d8ee2a290a353fdaefbc94cac614a32a51194d9c9206a4d7 |
| quarantined-web-artifacts/profile-receipts/maintenance-web-security-artifacts.json | b2fea3fbae6f70b40458164de21c67839ece4062c5023082dd894af0cc98c170 |
| connected-40263a330a1f48f0af541f7d8c187087/browser-connected.log | 870527a3d5e23f744f26d04f0bc79d34da9cd7ab36d19e0bd63d924c52005fe5 |
| connected-40263a330a1f48f0af541f7d8c187087/browser-agenda.log | 2e60908a9e37076b355e9c1d65138158eaac8b27834e9bc1de14cc8d7a207566 |
| connected-40263a330a1f48f0af541f7d8c187087/result.json | fddd21f308eaa4bdae20f889e1881747a537e62db1e91772bfc0906e204ce596 |

Rollback: preserve the baseline package and lock only for forensics/isolated comparisons; the known vulnerable baseline cannot be promoted as rollback. A functional regression keeps promotion blocked until a compatible fixed candidate passes the same gates. Preserve current failure logs and original attestations; do not rewrite source provenance.

## Final canonical closure

V376 closure: Preflight214 passed all164executed steps and all56compositions; only Docker unavailable. Canonical165packs/1521files/824Markdown/56plans,1269AUTHORED/145ADAPTED/107VERBATIM; ordinary67/755 and HTTP68/763. Exact Next/Sharp package-advisory repair, signed artifacts and92connected phases/4agendas complete. Current45/48; TEST02/03/07 and native/vendor license/security gates remain blocked. Final canonical consumer has754parent source files byte-identical to the connected candidate plus the two proven Next-generated type imports. HTTP executable matches V375 byte-for-byte, so its unchanged15minute runtime evidence is reused without repetition. Next work returns to TEST02 owner integration, including the public-web SEO/i18n boundary; do not repeat completed browser creation/recovery work or the fixed dependency qualification.

```json
{
  "checkpoint_under_test": 214,
  "full_preflight_executed_steps": 164,
  "all_executed_steps": "PASS",
  "all_compositions": 56,
  "availability_status": "BLOCKED",
  "missing_toolchains": [
    "docker"
  ],
  "inventory": {
    "packs": 165,
    "files": 1521,
    "markdown": 824,
    "plans": 56
  },
  "provenance": {
    "AUTHORED": 1269,
    "ADAPTED": 145,
    "VERBATIM": 107
  },
  "controls_passed": 45,
  "controls_total": 48,
  "whole_controls_unchanged": true,
  "preflight_json_sha256": "f282b689b569382c474c9085a0e89efef5f23c2d8b23cbce8a2f08cc36e91cab",
  "preflight_log_sha256": "238366e7e2d3870fdd51aac349c4393fb1ad0f060ccf3a1c4395e32d23af59df",
  "final_parity_sha256": "07e16946ff0263605db51c5685c6919e763901e1e1b95b2801de41ae35285c6a"
}
```

```json
{
  "status": "PASS",
  "files": 763,
  "packs": 68,
  "parent_files_identical_to_connected_candidate": 754,
  "generated_next_env_delta": "Exactly two type imports generated by the signed Next build; all other bytes equal",
  "next_type_generator_sha256": "3973ebf0093b08136059a790442951ba498421627a52bfcf06e1584360e034c9",
  "changes_from_v375": [
    "package.json",
    "pnpm-lock.yaml",
    "reference_http_metrics/source-lock.json"
  ],
  "http_binary_byte_identical_to_v375": true,
  "binary_sha256": "4643654625c71aaf5b45367bd3a198381257588d073ae0349fc90c1bc1b5ab85",
  "checks": [
    {
      "name": "final-http-tests",
      "exit_code": 0,
      "log_sha256": "2e3c54bca63f2844a59006b000d77bc4cdb2e4bf9524bbbd6decad8f1a415075"
    },
    {
      "name": "final-http-vet",
      "exit_code": 0,
      "log_sha256": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
    },
    {
      "name": "final-http-build",
      "exit_code": 0,
      "log_sha256": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
    }
  ],
  "reused_scope": "Only unchanged Go/rule/migration runtime evidence; no reuse of the old web dependency SCA"
}
```

FAIL737 preserves the first parity rejection and exact two-import correction. The final full Preflight read checkpoint214; successor215 records these completed receipts and the queued failure lesson. No materialized code changes follow that full run.
