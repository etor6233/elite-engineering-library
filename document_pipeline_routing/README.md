# Document pipeline routing gate

This directory is deterministic local `AUTHORED` control code. It is not AWS, Microsoft or Google code and it performs no classification, extraction, provider call or storage.

Its purpose is to stop an agent from guessing a document class, provider or implementation plan. The 21 baseline rows must always remain explicit so common enterprise document families cannot be omitted. A project may append any additional canonical kebab-case class ID; every additional class must pass the same schema, field, MIME, size, security, corpus, access, evaluation, provider-profile and no-storage contract as a baseline class. The distributed template is intentionally invalid until the user completes it.

The exact provider-to-plan map selects only library profiles that already wrap official public SDKs and source locks:

- Azure → `AZURE_DOCUMENT_RUNTIME_PACK_PLAN.md`;
- Google → `GOOGLE_DOCUMENT_RUNTIME_PACK_PLAN.md`;
- AWS → `AWS_TEXTRACT_DOCUMENT_RUNTIME_PACK_PLAN.md`;
- Microsoft MarkItDown local → `MARKITDOWN_LOCAL_RUNTIME_PACK_PLAN.md`.

Run:

```text
python -m unittest -v test_validate_document_routing.py
python validate_document_routing.py --configuration <completed-routing.json> --output <new-evidence-directory>
```

The V3 receipt preserves the baseline inventory, explicit custom-class IDs and the complete class → provider → role → exact pack plan → configured provider-profile path map, so the agent never has to reconstruct or guess a routing decision after validation. It only proves that the configuration is complete and hash-linked. It never proves provider accuracy or authorizes persistence. Every selected profile still has to run real security, provider, strict evaluation, corpus, review, load and recovery gates.

## Successor 0.3.1 — explicit local Paddle lane

PADDLEOCR_LOCAL selects markdown_system/PADDLEOCR_LOCAL_RUNTIME_PACK_PLAN.md exactly, as PRIMARY or EVALUATION. Azure, Google, AWS and MarkItDown lanes remain available. This is AUTHORED routing glue; no upstream runtime or weights are incorporated.

A local lane still requires an actual runtime probe receipt, authorized corpus, security policy and strict evaluation. Use a non-secret local-runtime reference and region local in the existing access record; synthetic qualification records prove validator behavior only. Selection never authorizes automatic storage, field accuracy or production. The 21 classes remain explicit; routing PASS is not OCR21 PASS.
