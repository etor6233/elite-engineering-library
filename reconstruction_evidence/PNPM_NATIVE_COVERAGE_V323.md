# pnpm native coverage V323

Result: focal source-lock/identity/provenance research passed; complete native admission remains BLOCKED. No new dependency version, product code, global installation, artifact execution, release or promotion. T2803 advances with explicit evidence; T2809 remains a separate owner.

## Exact scope and observations

pnpm11.25.0 contains four reflink0.1.19 addons (macOS arm64/x64, Windows arm64/x64) and fastlist0.3.0 x64/x86 executables. All six bytes match the V321 signed pnpm container inventory. The parent pnpm signature does not independently prove each native build's dependency graph or Microsoft CRT revision. This round did not download native tarballs/release executables or acquire an upstream source tree; it inspected existing bytes and public metadata/manifests.

The reflink fixed registry metadata identifies [pnpm/reflink commit2223f06f](https://github.com/pnpm/reflink/tree/2223f06f51cca2635846937281892673415f332b). GitHub reports its commit signature valid. Its Cargo.lock contains75package identities, including the local reflink0.0.0 source package; OSV2.5.1 scanned exactly that set with0known advisories, without resolution or call analysis. This is conservative source-lock coverage, not a binary SBOM or reachability result.

All four addon bytes contain the same rustc source revision f6e511eec7342f59a25f7c0534f1dbea00d01b14. The [official version file](https://raw.githubusercontent.com/rust-lang/rust/f6e511eec7342f59a25f7c0534f1dbea00d01b14/src/version) identifies1.82.0, matching the repository rust-toolchain file. Embedded strings support an identity inference; they do not replace a reproducible build. The fixed CI uses moving stable setup/actions and container tags, unfrozen pnpm installation and no displayed locked build flag. One Linux test path explicitly tolerates failure; no new build certification is inferred from its green workflow.

## Signature statements, not a fabricated binary receipt

The previously isolated Sigstore4.1.1 verifier (V321,50package SCA0) ran under exact Node24.20.0. It verified eight npm-publication/SLSA statements for all four native npm packages, with nine tampered-payload/wrong-workflow negatives. Trusted TUF roots and transparency/CT thresholds were enforced by the official verifier. The fixed certificate identity is pnpm/reflink CI.yml on refs/heads/main; the signed statement resolves that historical branch to commit2223f06f, never a moving source selection. Subjects equal the registry SHA512/SRI values.

The historical npm publication key expires2025-01-29; verified transparency entries precede that expiry. No claim of a current unexpired publication key is made. The native tarballs were not acquired in this round, so the signed statement is not reported as an independently verified tarball-to-binary correspondence. The current pnpm container still binds its own included files.

## Rust standard-library cross-check

RustSec advisory-db fixed commit8a1eb4f933fb5821add5b4e98601ebd90b8b3538 contains18rust/std records. All list patched boundaries no later than1.58.1 (the Error::type_id record uses >1.34.1);1.82.0 is outside those affected ranges. The current official Rust security endpoint supplies four std advisories: Windows batch fixes before1.77.2/1.81.0, remove_dir_all before1.58.1, and Cygwin separators affecting1.87.0–<1.89.0;1.82.0 is outside all four. This manually reviewed finite inventory supplements the crate scan; no npm/crates alias or general-purpose vulnerability matcher is invented.

[Rust's 2024 batch advisory](https://blog.rust-lang.org/2024/09/04/cve-2024-43402/) supplies the1.81.0 boundary. Cargo advisories in2026 concern the build/package manager and do not automatically label an already-built addon vulnerable; this round neither runs nor admits that old build toolchain. No compiler correctness, OS patch or complete future-advisory guarantee follows from this cross-check.

## Open native admission conditions — FAIL532

The exact reflink and fastlist trees contain no full LICENSE/COPYING file. The declarations available in metadata/README say MIT; a label does not manufacture copyright and permission text. Do not claim these projects are unlicensed; record the missing full-text evidence and seek an authoritative license artifact before independent redistribution/adoption. V322 preserves all supplied artifact license bytes; it did not close transitive redistribution obligations.

The [fastlist release source](https://github.com/MarkTiedemann/fastlist/tree/0da17e8f6654a1991e032faacb1494cfea1ece1f) is a small Windows process enumerator. Its [v0.3.0 release](https://github.com/MarkTiedemann/fastlist/releases/tag/v0.3.0) dates2020-06-15 and has no published asset digest in the current API. The release description in the current README says the runtime is included, while the fixed older README still describes packing: preserve that difference; do not infer CRT/build version from a filename. No exact Visual Studio project, compiler/CRT bill or reproducible build is provided by the inspected tree. Existing four Windows PE files are NotSigned and have no file/product version resource; that does not invalidate the verified enclosing pnpm signature.

Required reopening evidence: full upstream license notices; artifact-to-fixed-source/build and target-bound Rust/native dependency evidence; exact Fastlist static CRT/toolchain with official vulnerability evidence, or a compatible replacement qualified through its own profile/admission/regression gates. Keep pnpm's overall native/redistribution admission CONDITIONED/BLOCKED. No known vulnerability is asserted solely from age or absence of Authenticode.

## Reproduction and retained evidence

Staging `%LOCALAPPDATA%/Temp/elite-v323-b72a92c4fe7c4907947b1e919e2469d0`. OSV2.5.1 executable SHA25625e42f5ef6711fd8c0fb45390972205891dd44c6bd02ac93f0f63e8e98d9bfb6; scan source --no-resolve --no-call-analysis=go --no-call-analysis=rust --all-packages against the byte-identical Cargo.lock. Attestation verification uses V321's isolated official Sigstore tool; no native file is executed by the inspection scripts.

- `research-receipts.json` SHA256 `887c861f06c1c3be039ffc7271cae3cc2a4492b5e314df2594aa45d3484c929f`.
- `research-receipts-2.json` SHA256 `66a0eafec0cd9068fc07adc0b292044575243d7f3e68d147c85e3dbdfe33b412`.
- `reflink-osv.json` SHA256 `0992ae0a574c1a2ffb7df030b54b2f1233e076d7c33f22ed79d720231ee90754`.
- `native-identity-and-lock-coverage.json` SHA256 `0280fd2ee75e905e513ba0bdbb62a90c3cd9f0647d650738682853af7120f2ca`.
- `native-pe-metadata.json` SHA256 `1cdfa04c88c3792b684e96597c9c498ce0bd7a898826bd9f346f2bea8ded3484`.
- `rust-authority-receipts.json` SHA256 `0566d9f5a4109460d08dcd41572b5e045b023017cf0a4254d50c8a2241ac4b0f`.
- `rustsec-stdlib-records.json` SHA256 `bfdb371c8ba0fed2d9940cc4c96634038d7cbfdbdfbea8043d49e7c1fa404312`.
- `attestation-input-index.json` SHA256 `8f9116138e28f3de6aedf8cb0b7274c1ff9902029403baf957b621be5b4c8aaa`.
- `verified-native-statements.json` SHA256 `5b332c6174f2390c65808a1611c7b71375cc7817cd4a1161da154118d3fd2609`.
- `verify-native-statements.log` SHA256 `241eacdc5a2532ce636dd9246a478066f30c17f72cc9e23938ad700e947c257d`.

Continues: preserve the native blocker and inspect the existing T2809 status host/supervisor/identity/reporting contract; target/provider choices and readinessD–H remain pending. This is not100% completion.

## Integrated check and T2809 continuation audit

VERIFY_LIBRARY_PASS161packs/1453files/758Markdown/52profiles at checkpoint77; franchise67/746. Log SHA256 6b0e323085a05ff3eac18034e1211f3cc013046b308710c5e71756c7ef583885. Source inspection then found no NewStatusWorker constructor call outside tests in cmd/internal of the current reconstructed profile. StatusPrincipalSource and StatusReporter remain callbacks whose real identity/destination must be mounted; docs require supervisor alerting when reporting stops. Census SHA256 46f3e2aa80c3a4269d1b7500afca11c45fb3b7f705629096f03d11790ee3dd35. This is not a live host/supervisor/alert-delivery proof.

The current spec explicitly prohibits product expansion before semantic readiness; roundsD–H are still pending. Existing local maintenance corrections and research remain authorized. Real historical channel/format, commercial release condition and target/IdP/reporting choices already requested are not supplied by a synthetic test. No accounts, grants, provider or retention policy were invented to close T2809.
