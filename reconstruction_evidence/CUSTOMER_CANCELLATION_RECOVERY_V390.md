# V390 — customer cancellation journey closure

Library maintenance / existing TEST02 owner integration. Entry checkpoint242
resume passed; checkpoint243 records the bounded cancellation work. V386 native
research remains explicitly deferred. No new business policy or production work.

## What is closed in the reference

Customer account → appointment management → cancellation → server result →
current durable state is now reachable and tested as one journey.

| Case | Evidence and outcome |
|---|---|
| Normal cancellation and double click | Synchronous guard emits one POST; actual backend returns cancelled/successor version; full reload shows durable cancellation. |
| Response lost after commit | Interceptor forwards the real command and receives200, then drops only the response; UI reports unknown outcome, disables commands and offers scoped GET recovery. |
| Unbound HTTP200 receipt | A response for a different appointment never confirms success; GET shows original requested state, after which a new intentional command succeeds. |
| Concurrent and stale commands | Two actual BFF requests with the same version return200/409; the stale page receives conflict, then GET recovers the unique cancellation. |
| Identity and organization | Per browser: foreign resource409, wrong organization403, missing permission403 and missing session401; existing five direct HTTP negatives also run. |
| Receipt binding | Ten unit cases require exact appointment/organization/cancelled state/successor numeric version; null, array and mismatched receipts fail. |
| Empty and navigation states | Account has a direct management link; empty appointments show a clear message; recovery is an ordinary GET. |

The client reuses the existing bounded10second request pattern; this run does
not claim a separate real-time timer-expiry campaign. No automatic POST retry or
new cancellation eligibility exists. The immutable appointment instant and its
configured display zone remain unchanged. Error text no longer asserts that a
committed cancellation did not happen or exposes raw internal error codes.

## Actual execution

294web tests pass with the same one inherited opt-in skip; typecheck and build
pass. The connected gate runs separately and is not substituted by that skip.
Actual Next/JWE → Go/RS256/JWKS → PostgreSQL, all53unchanged migrations, four
browser projects. Sixteen distinct synthetic appointments end cancelled at
version2, with exactly16transitions by customer-1/reason customer-request and
16matching appointment.cancelled outbox events. Twenty-eight backend POSTs
include the rejected duplicate/stale/foreign-resource attempts. No extra audit
or outbox rows exist. Three untouched fixture appointments remain requested/v1.
Seven-table snapshot comparison preserves every unrelated record plus original
appointment fields excluding the three intentionally mutable state/version/
updated_at columns. The isolated fixture process tree is stopped after the run.

Eight previously green read/time-zone browser runs (two business configurations,
four browser projects each) pass on the same final source. They preserve their
seven-table snapshots, zero writes, six-list cursor traversal, role scopes,
fixed read recovery and customer date consistency. Source guard and actual
DOM/new-document synchronization are checked without retries or error filtering.

Six canonical packs reconstruct. Final69packs/805files match the tested source,
parent68/797. Production Go/domain/SQL/locks unchanged; the HTTP reference binary
is byte-identical to V375.8589authenticated dependency files are unchanged.
New code is AUTHORED and retains CONDITIONED admission, with no claim that tests
replace complete implementation assurance, native qualification or target gates.

## Failures preserved

FAIL776: old client says No se canceló after response loss and accepts arbitrary
HTTP200JSON. FAIL777: all four baseline browsers demonstrate two outgoing POSTs
from a same-event double click; backend CAS prevented a duplicate effect.
A single-click diagnostic reaches the lost-response defect in Firefox/WebKit;
the Chromium diagnostic instead reveals the same reload race as below.
FAIL778: first corrected run passes Firefox/WebKit while Chromium inspects an
execution context being replaced by the successful reload. The final fixture
registers DOMContentLoaded before the action, waits for the new document and
then inspects it. Business/result/concurrency assertions remain intact.
Original failures, code candidates, traces and separate databases are retained.

## Closure boundary

This customer cancellation path is closed in the isolated reference. It is not
the closure of all customer journeys, all21core functional equivalences, full
localization, reminder delivery, provider/IdP live, production deployment or the
release. TEST02/03/07 and all48acceptance definitions/statuses remain unchanged:
45/48. Native research remains deferred by user instruction; release cannot be
promoted through that deferral. Do not reopen this cancellation defect without
a relevant delta; continue the next material owner gap from the existing roadmap.

## Receipts

Stage: Temp/elite-v390-current.txt; tools.json identifies the reused installed
candidate, while V389 final-consumer preserves the previous source. This stage
contains a separate final reconstruction and all baseline/failure pointers.

| Receipt | SHA-256 |
|---|---|
| web-qualified-results.json | 36d4ba6ca48e0f0c41e9ee04813c3efae9062e1de35b37edf9324ad2fb46ef0b |
| artifact-parity.json | 7563f34b2a058f7ce4c0f799cb29ad7b97d517326579e335c400497c8723cb43 |
| canonical-pending-changes.json | 8f6192bd5d39c24c3e839dd445dbcf06bb6957b9cb63049e7233059279a949da |
| final-parity.json | 66c7d36a50b97d1ce9d8af18930122d3c029a8875306a33a35f3e93206a9a0ce |
| baseline-double-click-stage.txt/result.json | 38c39aa9e8527a42c090f3964e265b6823b76af2b98f519806569cb82a47bc25 |
| baseline-double-click-stage.txt/admin-reads.log | 4aa45a6c5c22f4f6efd5d5ed2de272c91f6645caca006821e8ac9daccf4c361d |
| baseline-lost-response-stage.txt/result.json | 2ed4c43588400a26a769aef68145390a80bab75ab3d1d4d83f599ba907056be5 |
| baseline-lost-response-stage.txt/admin-reads.log | fa3553cfa549b3bf0db333ab139a167a6a2bfa123db311b6153c7456b61637b2 |
| first-fixed-failed-stage.txt/result.json | 8613675d6ad53477523d688c0be1eb4dde3566cd9cc2d4cfceb474c7580b1b67 |
| first-fixed-failed-stage.txt/admin-reads.log | 2bb541053f8ff97fa4ea86e6985d845a4b902a9e689b79970e8eebb9ed00f133 |
| connected-ar-stage.txt/result.json | f08c4dd3f3e7226d1bf008fb9e2ecf058fe23226b33c6c211b0cd94800702106 |
| connected-ar-stage.txt/admin-reads.log | 4c3700cb3a0e535d25f4ace8887ccdaddf8a6e15c93231ff036d0fca77928128 |
| connected-read-ar-stage.txt/result.json | cdc7d9677509ccd1a5403bc886c000eab7bf9499f53a3b9a650177737c40b015 |
| connected-read-ar-stage.txt/admin-reads.log | 371f7f7d2aecf517792fbd42b16d7d9f7d9fcb138f952fd28d0a4a4448c0f83b |
| connected-read-ny-stage.txt/result.json | a8163052743f97f3d5756b42622727f8abd255e4edba00b421f3617cffccf064 |
| connected-read-ny-stage.txt/admin-reads.log | 2440736ff6cbe99f3d365d616a6dd1b46cc08037e08d61b7700649b0edcb0b83 |

Full Preflight244 pending.

## Final checkpoint245

V390 cierre245: recorrido cliente cuenta→turnos→cancelacion→recuperacion GET cerrado en referencia.294tests/1skip previo,4runs cancelacion con16cambios/audits/outbox unicos y28POSTincluidos rechazos;8runs lectura/fechas siguen sin escrituras.6packs/805files reconstruidos,8589artefactos yHTTPbinario idénticos. Preflight244164pasos/56composiciones PASS; Docker ausente.165packs/1561files/838Markdown;68/797.45/48; TEST02integral pendiente, V386native diferido yTEST07dependiente.

The bounded cancellation reference closure above is supported by actual durable counts and final canonical reconstruction. Do not describe the28POST attempts as28successful database mutations:16succeed and12are rejected conflicts. Initial unknown-response/double-click and fixture-navigation failures remain preserved. Broader control and target conditions remain open.

```json
{
  "checkpoint_under_test": 244,
  "steps": 164,
  "step_status": "ALL_PASS",
  "compositions": 56,
  "inventory": {
    "packs": 165,
    "files": 1561,
    "markdown": 838
  },
  "franchise_files": 797,
  "http_reference_files": 805,
  "web_tests": 294,
  "inherited_opt_in_skip": 1,
  "cancellation_browser_runs": 4,
  "read_regression_browser_runs": 8,
  "synthetic_cancelled_appointments": 16,
  "matching_audit_rows": 16,
  "matching_outbox_rows": 16,
  "cancellation_api_post_attempts_including_rejections": 28,
  "scope_negatives_per_cancellation_browser": 4,
  "unexpected_durable_effects": 0,
  "read_regression_writes": 0,
  "snapshot_tables": 7,
  "missing_toolchains": [
    "docker"
  ],
  "controls_passed": 45,
  "controls_total": 48,
  "acceptance_unchanged": true,
  "preflight_report_sha256": "ebcf1f83d3b4fb20c65c6868a9e84653dd04128dc69131ba69f6c46af98f0d03",
  "final_parity_sha256": "66c7d36a50b97d1ce9d8af18930122d3c029a8875306a33a35f3e93206a9a0ce",
  "dependency_file_parity_sha256": "7563f34b2a058f7ce4c0f799cb29ad7b97d517326579e335c400497c8723cb43"
}
```
