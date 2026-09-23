# Local runtime 0.9.0 — current library scope

The selector and legacy planner keep their historical nonexecution receipts.
The separate local_runtime.py entry point now admits one exact command: a
frozen offline copy install for the three hash-locked consumers. It uses only
Node24.20.0 (exact executable SHA) with --jitless --no-addons and the exact
ZIP-disabled445-file projection. In this pinned process, WebAssembly is absent
in main and worker contexts and native addon loading returns ERR_DLOPEN_DISABLED.
Those restrictions apply to the package installer. Product build/application
runtimes have separate dependency/native gates.

The source inventory retains467published package archives, eight nested
OpenPGP archives,42exact source-map/member matches, six embedded WASM identities,
and477npm/Python identities with dated OSV2.5.1 results (472base+5new, zero findings).
This is not an independent native rebuild or a claim that the inactive parsers
are safe for network use. Three actual restricted installs and11negative router
cases pass;8focal preservation/recovery tests cover payload/notices, exit failure
and timeout. No live secrets, pnpm exec or arbitrary CLI arguments are admitted.

## Start from an already verified official archive

1. Acquire the locked official archive and receipt with the existing official
   source profile. Materialize it with select_artifact.py into an absent path.
2. Run contain_zip.py create --projection <verified-selection> --target
   <absent-contained>. The original full455-file extraction is also supported
   by --source. Both routes produce the same445fixed files; selection validates
   all original archive entries before any excluded payload is omitted.
3. Run local_runtime.py prepare --consumer enterprise-web|playwright|lighthouse
   with --project, --projection, --contained, --node, --store, --cache, --target,
   --acquisition-receipt and --acquisition-sha256. Each is an explicit local
   path except the externally retained acquisition receipt SHA. The target and
   consumer node_modules must be absent. The official Node path is an executable,
   never a PATH-resolved shim. Archive/source stores remain outside product output.
4. Retain the printed plan SHA separately. Run local_runtime.py verify --target
   <plan> --sha256 <printed-SHA>, then execute with those same arguments. Execute
   verifies the entire recipe and all inputs again, runs once and writes a
   durable start marker, log and result. It preserves notices and immutable
   consumer inputs. Environment is replaced, not merged; scripts/hooks disabled.

The planner supplies local-tool-notices.json and readable THIRD_PARTY_NOTICES.md
in addition to the original47texts, BlueOak/QRCode/semver and MPL source notices.
Text copies and published declarations preserve their actual sources and SHA.
Missing upstream copyright text is never filled with an invented author/year.
This admission covers private local use, not redistribution of the altered pnpm
payload. Keep that payload out of both library and product ZIPs; acquire/adapt
it on each consuming machine. GPL/LGPL/MPL private-use and redistribution scopes
remain distinct: https://www.gnu.org/licenses/gpl.en.html and
https://www.mozilla.org/en-US/MPL/2.0/FAQ/ .

A failed/aborted attempt retains its marker, log and any partial consumer files.
Do not erase or reuse that plan. Rebuild the consumer in a new absent staging
directory from its unchanged manifest/lock, then prepare a new plan. Source,
Node, notices, workspace configuration or consumer hash changes fail closed.
--offline is not an operating-system network sandbox; trusted local inputs and
the fixed recipe are conditions of this claim. A new upstream advisory, runtime
or input revision reopens admission. Archive redistribution, production release,
security testing of network parsers and Daybreak/libxml2 are not granted here.

Reference: reconstruction_evidence/PNPM_LOCAL_RUNTIME_V402.md/json.

## Historical selector and planner documentation (preserved)

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

