# V341 — Scoped source correspondence and remaining closure scope

Date: 2026-09-08. Library maintenance; read-only source qualification.

## Outcome

Five AST units from two bundled Undici7.29.0 modules match the fixed source at
9e38fc121d2eb26086d41c7d9379b47a6fada1c5 under explicit per-unit identifier renames.
Both source files were independently rechecked against the fixed Git tree blobs.
Two exact byte ranges were extracted from the existing selected pnpm11.25.0
bundle. No package source was executed; the existing local parser alone ran.

| Unit | Declared bundled → source identifier renames |
|---|---|
| formdataEscape arrow function | str2 → str |
| normalizeLinefeeds arrow function | none |
| WebsocketFrameSend.createFrame | buffer4 → buffer; i4 → i |
| WebsocketFrameSend.createFastTextFrame | buffer4 → buffer; i4 → i; head2 → head |
| generateMask | buffer3 → buffer |

AST comparison ignores source offsets/locations and the lexical spelling of
ordinary literals. It preserves actual literal values, regex patterns/flags,
operators, declaration kinds, argument order and template raw values. Declared
renames do not rewrite noncomputed member/property names. The comparison is
syntax correspondence within these selected units, not a proof that surrounding
bindings, dependencies, bundler configuration or complete modules are equivalent.

Nine initial negative controls include altered AST values/operators, omitted
renames, property-name retention, template raw retention and syntax rejection.
An independent pass additionally changes eight actual source snippets before
parsing: replacement character, regex flags/pattern, XOR operator, mask bound,
frame-size threshold, mask flag and return binding. Every alteration differs from
the original canonical AST. Two byte ranges independently equal their exact
slices in the untouched bundle. These controls are analysis-harness checks,
not new product tests or a promotion of the 48-row contract.

## Method and attribution boundaries

The existing Next16.3.2 analysis environment contains Acorn8.14.0 with its MIT
licence. Its executable bytes, local path, licence and Node24.20.0 identity are
recorded; nothing was installed or updated. Acorn parses JavaScript to ESTree;
the normalization/comparison harness is local AUTHORED analysis, not upstream
code. An initial optional TypeScript JS-parser probe failed because the assumed
legacy lib/typescript.js file is absent from the observed7.0.2 layout (FAIL576).
Inventory located the existing Acorn parser and its path/hash were verified.

Method source: [official Acorn parser interface](https://github.com/acornjs/acorn/blob/master/acorn/README.md).
Fixed source: [Undici body](https://github.com/nodejs/undici/blob/9e38fc121d2eb26086d41c7d9379b47a6fada1c5/lib/web/fetch/body.js),
[Undici frame](https://github.com/nodejs/undici/blob/9e38fc121d2eb26086d41c7d9379b47a6fada1c5/lib/web/websocket/frame.js).

V340 preserved the formdata-polyfill and ws attribution comments independently.
This work now relates selected code units in those Undici source files to their
bundled representation. It does not establish the original formdata-polyfill/ws
revisions, full upstream licence applicability, complete Undici7 reconstruction,
runtime reachability or redistribution admission. External Undici6's103-file
parity remains a different version/scope. Yarn registry/tag binding (FAIL575),
semver-utils/chownr evidence and LGPL source/relinking/delivery remain open.
Draft47 is unchanged; no additional licence text was promoted into it.

## Approximate remaining scope

The contract has41passed and7pending test rows: 14.5833% of unweighted rows remain,
rounded to approximately15%. This is the only percentage calculated here. It is
not an estimate of remaining effort or global readiness. All ten macrofronts
T2801–T2810 retain incomplete status, with substantial partial evidence.

| Pending row | Remaining closure | Dependency boundary |
|---|---|---|
| TEST-02 | cross-core contracts, coverage and claims | further local/source audit |
| TEST-03 | security and privacy over the final composition | exact graph, resolved findings and integrated evidence |
| TEST-05 | integrated telemetry, diagnosis and sink behavior | local integration plus selected operational target |
| TEST-06 | recovery from the final distributed artifact | candidate artifact after material closure |
| TEST-07 | independent reproducible release, SCA/notices/signature | final candidate, source/licence and release gates |
| TEST-08 | NEW/EXISTING usability of the exact artifact | final artifact and acceptance evidence |
| TEST-09 | governed historical-data import and live seam | authorized corpus and official channel contract absent |

No total-effort weighting exists, so an overall percentage would be invented.
The independent execution/preflight baseline is already substantial; these broad
closure tests cannot be weighted equally to each narrowly scoped passed test.
No user decision or approval is inferred from the numerical estimate. Historical
roundD and deferred ARCA remain as agreed; continue independent maintenance.

## Integration scope

Documentary-only update:162packs/1461materializable files/53profiles remain;
779Markdown expected after this report. All implementation pack hashes, selected
payload442 and draft47 are checked again at closure. V337 Preflight113's154steps
retain their original date; no unchanged executable suite is rerun or relabelled.
The structural gate follows checkpoint122, followed by final resume/plan checks.

## Evidence

Stage: %TEMP%/elite-v341-171c0a35821f4cf6aacf2ab09cd9d6c1.

| Evidence | SHA256 |
|---|---|
| ast-correspondence.json | `ffa5e50ed86c5303cb6646b60dd2eb0e9e5bcb87b132468b321afed37f030605` |
| independent-controls.json | `04e11f13e846951dd3fa6116aa4c285e2437b5784d14abe61e508c711290e44d` |
| map_undici.cjs | `e811a1d7cdfba00709f467651bc06d942e18e9f9028de00c1ff19e7de6d1db1e` |
| independent_controls.cjs | `30a09eff7b29511001339321e8a0d0b51eeda9be42ef5dac5181c63ea06ff125` |
| ast-first.log | `56fb08e4bfa6f8665e4023866f2b554538e1ba930aef6622fe798618d8a74034` |
| independent-controls.log | `91afcdfb9368ae2f5510e2aff9e4ae35bc444419cfd59448f219004aa38bc2d1` |
| body-bundle-region.evidence | `add1411bea4c2891ad5befa17b050b566dcdc4138590ccaca4aae31c9043764a` |
| frame-bundle-region.evidence | `bb06aebf91386947915aebf595523fce29f2d570e03694f3e3e68a28b84161b7` |
| research-tool-receipt.json | `c9d5d14714e76f3a5f8c71771dcb00a5801e0179133ae0ea582abd9267df9865` |
| research-parser-LICENSE.evidence | `76a876cf886ff9be2a8b5e2e86514fed06223c8c9f0c1e9ee9606e93841e00b7` |
| remaining-scope.json | `02284a60d44400970e0fabf22ade8b0fe065897de82201e31880345f9e6f7327` |
| parser-probe.txt | `4fd8ac852f8b8143c372b178e6f5288725f06667f82a0f0c0b488a19b6a57dd4` |

## Structural correction before replay123

VERIFY_LIBRARY122 rejected a newly added machine-local Windows home reference in PROJECT_FAILURE_LESSONS (FAIL577). Its failed log and before-fix122 snapshot remain intact. The canonical reference now uses %TEMP% plus the stage name; no verifier rule was weakened. Existing ledger lesson2081 informed the portable correction. Gate replay follows checkpoint123; no premature structural PASS.

Failed log SHA256: `08a2b7ae61f90549471691bd0422bf8286dbc3c287e7a80c8453907d8a79240a`.

## Closure124

VERIFY_LIBRARY123 PASS:162packs/1461materializable files/779Markdown. All53profile counts equal V340; all163implementation-directory Markdown hashes remain unchanged. Independent checks verify the unchanged47-copy notice draft and exact442-file selected payload. The ledger has2502unique canonical failure IDs. FAIL577 is regression-proven by structural replay; FAIL576 is recovered; Yarn FAIL575 remains unresolved. Final resume and contract-plan validation follow checkpoint124. No unchanged154-step Preflight is rerun or relabelled.

The five AST units are qualified only within their explicit mapping. Next work remains full source/build binding and licence delivery by owner; pending7/48test rows and ten macrofronts are not closed by this syntax evidence.

| Closure evidence | SHA256 |
|---|---|
| verify-library123.log | `2f0970d7cb4f1a468dd12cfe0fd15c85085bca1a4fec7c82d6977e7a0bc84a6c` |
| structural-closure.json | `b9dcd2148677fcb46399f0f77a552445a54781b0d084feb95787c050f08f2019` |
| draft-payload-verification.json | `801263f739075dc97edd45d36eb4a3c60b9473e3a35a520fb925a8aa744441ba` |
| final-uniqueness.json | `60420ff3e16e6db2a7acd98ccb84886963923bb3e658d960399449e73b8c59ff` |
