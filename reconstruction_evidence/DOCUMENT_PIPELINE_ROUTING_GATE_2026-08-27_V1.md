# Document Pipeline Routing Gate — Reconstruction Evidence V1

Date: 2026-08-27  
Result: `REBUILD_VERIFIED / CONDITIONED`  
Production/admission result: **not promoted**; this gate validates project decisions and calls no classifier, extractor, model, provider or storage.

## Official leader basis and exact boundary

Current official sources were rechecked before implementation:

- AWS `accelerated-intelligent-document-processing-on-aws`: latest release remains `v0.6.5`/commit `1b5fd74454e593de233342a02ee911af8ee38359`; exact source is already in the 72-source lock under MIT-0 + NOTICE. AWS documents page/package classification, segmentation, JSON Schema document classes, versioned configuration, evaluation and HITL. Current main head `720f3052a57c22ae0d1e50886965c065ee86d60b` was unsigned and was not substituted for the release.
- Microsoft `content-processing-solution-accelerator`: stable `v2.1.2`/commit `b47cec48475cfedd7109debf2f407ef0c5530311` is already locked under MIT with its `uv.lock`. Microsoft documents the four-stage per-document pipeline Extract → Map → Evaluate → Save. Current signed main remains `659eaa1f503dd08b1e1aea1c72eab11c7c191d00`; its recorded failures prevent silent replacement.
- Google `document-ai-samples`: current signed main remains locked commit `001ba391ab4a2f40d001cc0387618cb3c3699523`, Apache-2.0. Google demonstrates specialized processors plus custom classification/extraction and explicitly says the repository is demonstrative, not an officially supported product.

Primary sources: https://github.com/aws-solutions-library-samples/accelerated-intelligent-document-processing-on-aws, https://github.com/microsoft/content-processing-solution-accelerator and https://github.com/GoogleCloudPlatform/document-ai-samples.

No upstream publishes a universal router that knows a new company's document classes, schemas, access, corpus and storage policy. The local validator is therefore declared `AUTHORED`; it copies no vendor engine and is never labeled AWS, Microsoft or Google code. It operationalizes the existing project decision contract and only maps an explicit provider to an already verified library plan.

## Materialized capability

`DOCUMENT-PIPELINE-ROUTING-GATE 0.1.0` produces four files:

- a deliberately incomplete 21-class JSON template;
- a CPython stdlib validator/atomic receipt writer;
- eight positive/negative tests;
- a README with provenance and non-authority boundaries.

The 21 mandatory inventory IDs cover supplier invoices, receipts, proformas, commercial invoices, packing lists, purchase orders/confirmations, BOL/AWB, delivery notes, origin/quality certificates, customs, freight/insurance, quotations, product specs, contracts, bank/payment statements, warranty/service claims, identity/regulatory documents, multi-document packages and a named extension slot.

Every row must be exactly `REQUIRED`, `OPTIONAL`, `NONE_WITH_REASON` or `BLOCKED`; `AWAITING_USER` is rejected. Active rows require schema version, fields, MIME/byte/page limits, authorized ground truth, owner, real access probe evidence, security policy/receipt schema, strict evaluation profile and storage false. Candidate providers are closed to Azure, Google, AWS or MarkItDown local. Their plan paths are fixed, relative and traversal-safe; every active class has exactly one primary lane.

The receipt contains the configuration SHA-256, required/optional classes and selected pack plans. Its next gate is materialization followed by real security/provider/evaluation. It never claims that routing, provider accuracy or storage is proven.

## Tests and reconstruction

- source-tree `py_compile`: PASS;
- eight tests: PASS;
- Markdown pack materialization: 4/4 files;
- source-tree versus canonical materialization SHA-256: 4/4 equal;
- canonical `py_compile` and eight tests: PASS;
- composition `DOCUMENT_PIPELINE_ROUTING_PACK_PLAN.md`: 2 packs/22 files PASS;
- global structural verifier after addition: 42 packs/370 files and 21 profiles PASS.

Negative coverage rejects empty template/owners, `AWAITING_USER`, missing/duplicate inventory, no reason, provider/plan mismatch, path traversal, missing security/corpus/access/evaluation, unproved access, unauthorized corpus, evaluation attempting storage, storage true at either level and occupied output. The one initial test expected the later `AWAITING_USER` message but the template correctly stopped earlier at missing project owner; retained as `LIB-FAIL-251`, the expectation was aligned and a separate test continues proving `AWAITING_USER` rejection.

## Honest boundary

This gate makes the agent ask and prove the prerequisites before selecting code. It does not make business documents accurate, provide accounts, create schemas, supply ground truth, run malware engines, call cloud services or authorize persistence. The four selected document profiles remain conditioned by their exact official SDK/source, real project inputs and provider-specific gates.
