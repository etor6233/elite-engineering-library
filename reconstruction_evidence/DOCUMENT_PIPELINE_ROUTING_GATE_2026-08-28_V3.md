# Document Pipeline Routing Gate V3

Date: 2026-08-28  
Pack: `DOCUMENT-PIPELINE-ROUTING-GATE` 0.3.0

## Problem closed

The official Google Custom Document Extractor sample admitted in V83 accepts `DocumentSchema` through `ProcessOptions.schema_override` on each request. The prior local routing validator nevertheless rejected every class ID outside a hard-coded inventory of 21 names. That would have forced a code edit whenever a project introduced a new legitimate document family.

V3 removes that artificial ceiling without weakening the baseline inventory or any project gate.

## Contract

- The 21 baseline enterprise classes remain mandatory and unique. Omitting one still fails closed.
- A project may append zero or more additional classes.
- Every additional ID must be canonical lowercase kebab-case and unique.
- An active additional class passes the exact same validator as a baseline class: schema version, required fields, supported MIME, positive byte/page limits, one PRIMARY provider lane, exact provider→pack-plan mapping, canonical provider-profile path, secure-file receipt schema, authorized ground truth, named owner, sandbox access `PROVEN`, strict evaluation and `automatic_storage=false`.
- No provider call, extraction, model choice or storage occurs in the routing gate.
- Receipt V3 records `baseline_class_ids`, `custom_class_ids`, required/optional classes and every class→provider→role→plan→profile route, bound to the configuration SHA-256.

The control code is declared `AUTHORED`; it is not presented as AWS, Microsoft or Google source. Extraction remains delegated to the exact official engines selected by the receipt.

## Regression

The positive additional-class test appends `letter-of-credit-amendment` without modifying validator code and supplies:

- schema `letter-of-credit-amendment/v1`;
- required fields `credit_number`, `amendment_number`, `effective_date`;
- Google PRIMARY lane and exact Google pack plan;
- security, corpus/ground truth, access and evaluation evidence paths;
- storage disabled.

The receipt preserves both `supplier-invoice` and the new class, lists the custom ID explicitly and selects both AWS and Google plans. A negative test rejects `../../new class` before any output.

Results from a fresh Markdown materialization:

- 4 implementation files reconstructed with manifest↔FILE and SHA gates;
- Python syntax PASS;
- 11/11 unit tests PASS;
- missing baseline, duplicate, incomplete decision/reason, unsafe path, mismatched provider plan, missing security/corpus/access/evaluation, storage true and occupied output remain rejected.

## Composition

`DOCUMENT_PIPELINE_ROUTING_PACK_PLAN.md` now selects pack 0.3.0 and remains 2 packs / 27 implementation files. The source acquisition pack is unchanged at 0.4.64. The agent can therefore add a document class from project configuration rather than authoring a new router branch or searching for a class-name-specific sample.

This closes a local orchestration limitation. It does not claim that an arbitrary provider response is correct; correctness evidence remains a property of the selected official engine, schema, corpus and evaluation record.
