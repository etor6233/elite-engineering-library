# Document task map 0.1.0

This AUTHORED contract keeps all 21 document classes visible while separating proposed task names from selected business scope, source code, UI, runtime and field accuracy. It coordinates existing routing, security, strict field evaluation and connected review owners. It is not another document extraction engine.

Only supplier-invoice has the existing four-field reference (invoice_number, vendor, total, currency) and /experience/documents UI. This revision does not rerun that historical runtime. The other 20 classes do not inherit its UI, fields or PASS. A proposed task identifier is a planning link; it is not a implemented permission or business operation.

For each required class, the project must bind a schema and fields, authorized ground truth, provider lane and real probe, limits, secure ingestion receipt, field evaluation, review operations and domain commit/reconciliation owner. The routing template requires an explicit decision for every class and is intentionally incomplete before those bindings. Do not discard a class silently. Payment evidence is not confirmation that money settled; a technical specification requires units, revision and authoritative compatibility rules. Country/tax rules cannot be inferred from a sample.

Preserve review separation: the uploader proposes edited fields, then another authorized reviewer approves or rejects the exact payload SHA with a reason. The reviewer does not silently modify and approve the same proposal. Changes require a new proposal/version. Preserve uncertain-result recovery through the existing owners; no automatic repeated mutation.

Preserve safe original downloads: attachment, application/octet-stream, no-store, nosniff. Do not turn uploaded HTML/PDF into an inline active preview. Any future preview needs a separate safe rendering design and tests.

Run `python -m unittest -v test_task_mapping` here, or `python validate_task_mapping.py document-task-mapping.json`. PASS means this planning contract preserves inventory/state/review/download invariants. It does not mean OCR21, physical scanning, storage or production is accepted.
