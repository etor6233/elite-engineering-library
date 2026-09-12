# Document Pipeline Routing Gate — 2026-08-27 V2

Status: `REBUILD_VERIFIED / CONDITIONED`

## Scope

`DOCUMENT-PIPELINE-ROUTING-GATE` was advanced from 0.1.0 to 0.2.0 without adding a classifier, extractor, provider call or storage effect. The pack remains local `AUTHORED` control code and is not represented as AWS, Microsoft or Google source.

The V1 configuration contract still requires an explicit 21-class inventory, named business/security/data/review owners, one primary provider lane per active class, exact library pack-plan mapping, a canonical configured provider-profile path, security receipt contract, authorized corpus/ground truth, proven access evidence, strict evaluation and `automatic_storage=false`.

## V2 receipt correction

V1 preserved only the deduplicated selected pack plans. That was insufficient for a fast handoff because an agent could have to reconstruct which class selected which provider role and configured profile. Receipt schema `elite-document-pipeline-routing-receipt/v2` now preserves, for every active lane:

- class ID and `REQUIRED`/`OPTIONAL` decision;
- provider and `PRIMARY`/`EVALUATION` role;
- exact closed library pack plan;
- canonical configured provider-profile path.

The receipt remains hash-linked to the complete routing configuration, atomic, non-overwriting and explicitly denies automatic storage.

## Reconstruction and tests

The four manifest files were regenerated through the canonical updater and their declared SHA-256 values match reconstructed bytes. On CPython available in the audit host:

```text
python -m py_compile validate_document_routing.py test_validate_document_routing.py
python -m unittest -v test_validate_document_routing.py
Ran 9 tests
OK
```

The new regression proves a two-lane AWS-primary/Google-evaluation class preserves both exact pack plans, roles and provider-profile paths deterministically. The prior negatives for incomplete template, missing/duplicate class, `AWAITING_USER`, reason, provider-plan/path mapping, security/corpus/access/evaluation, storage authorization and overwrite remain green.

## Boundary

This change eliminates re-guessing after routing; it does not prove provider availability, extraction accuracy, document authenticity, production performance or legal authority. Those remain real project gates. The governing public sources remain the exact locked official AWS accelerated IDP, Microsoft Content Processing accelerator and Google Document AI samples; none of their conditioned code is copied into this local gate.
