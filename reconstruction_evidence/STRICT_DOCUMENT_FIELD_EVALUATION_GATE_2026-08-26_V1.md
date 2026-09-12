# Strict Document Field Evaluation Gate — Reconstruction Evidence V1

Date: 2026-08-26. Scope: exact AWS Labs Stickler admission plus local fail-closed orchestration. This is not evidence of a REVESTEX corpus, an extraction provider, or production storage readiness.

## Authority fixed

- official repository: `awslabs/stickler`;
- release: `v0.6.0` / distribution `stickler-eval==0.6.0`;
- commit: `174ca9d3476c1ea2d0a36628c73d595604e0398d`, GitHub verification valid;
- source ZIP: 1,356,931 bytes, SHA-256 `f24ae4650b3ffc46d23e5b447a0e72a8f75a20e98df1b4ee0ca31490759707a3`;
- `LICENSE`: Apache-2.0, SHA-256 `09e8a9bcec8067104652c168685ab0931e7868f9c8284b66f5ae6edae5f1130b`;
- `NOTICE`: SHA-256 `d4290ed64c2edd0fce1d84e3f9dfb2881240fe534def76b8cd29ed6af683e287`;
- `uv.lock`: SHA-256 `325f27921c235b3a3620d3523e615fa3f63525a09deb0ed2fad96da901884dd6`;
- official source profile `strict-field-evaluation` selects exactly that one lock entry and lists five explicit production blockers.

AWS source was not copied or modified. Every materialized application file in `STRICT_DOCUMENT_FIELD_EVALUATION_GATE.md` declares `AUTHORED`. Imports use public APIs of the exact official engine; local policy/path/hash/schema/receipt code is not attributed to AWS.

## Upstream evidence retained

`uv sync --frozen` selected Python 3.12.13 and 75 exact development packages. The official suite executed 1,543 cases: 1,535 PASS, 2 skip and 6 FAIL because `signal.SIGALRM` is unavailable on Windows. A direct adversarial probe showed root extra FA/FP, missing duplicate row FN and wrong exact value FD/FP, but an extra inside `List[StructuredModel]` produced zero FA/FP. The same issue remains in the signed 0.7.0 candidate inspected on 2026-08-26. These are retained as `UP-FAIL-052` and `UP-FAIL-053`; no local patch is represented as upstream.

## Local containment and tests

The reconstructed gate validates both expected and actual payloads with Stickler's own Draft-07 validator before comparison, while an independent local policy requires:

- a closed root and every nested object (`additionalProperties=false`);
- every declared property required;
- local JSON Pointer references only;
- no combinator/dynamic schema surface admitted by this version;
- `ExactComparator` and threshold `1.0` for every scalar;
- exact source identity and installed distribution version;
- SHA-256-linked manifest inputs, confined non-symlink paths and create-only receipt;
- zero `fa`, `fd`, `fn`, `fp` and `overall_score == 1.0` for every case;
- `automatic_storage_authorized=false` unconditionally.

Clean reconstruction from the Markdown produced five byte-identical files. The twelve tests executed against the frozen official environment and passed:

```text
Ran 12 tests in 0.247s
OK
STRICT_FIELD_REBUILD_IDENTICAL files=5
```

Covered negatives: root extra, nested-list extra, required field missing, duplicated ground-truth row, incorrect exact value, open nested schema, external reference, altered input hash, installed version mismatch, schema-error raw-value leakage and output overwrite. The positive receipt contains no invoice number, SKU or raw approval ID.

The upstream acquisition pack reconstructed 18 files and its suites returned:

```text
UPSTREAM_ACQUISITION_TEST_PASS negatives=5
SOURCE_PROFILE_VALID profile=strict-field-evaluation selected=1
SOURCE_PROFILE_TEST_PASS valid=9 negatives=5 positives=1
```

## Failure learning

`LIB-FAIL-145` through `LIB-FAIL-168` retain every path, invocation, binding, placeholder, version-alignment, glob, patch, privacy, cleanup, archive-smoke and global-verification failure observed during this slice. None was deleted or reclassified as upstream. Their corrections are now represented by exact path enumeration, native array binding, canonical skeleton hashes, sanitized schema rejection, safe cleanup, all-plan composition and atomic patching.

## Decision

`STRICT-DOCUMENT-FIELD-EVALUATION-GATE 0.1.1` is `REBUILD_VERIFIED / CONDITIONED`. It is immediately materializable and its tests are executable once the exact official source environment has been acquired and frozen. It is not `REUSABLE_PACK`, because no project corpus, classes, service access, Linux upstream suite, load/drift/privacy/review/storage evidence exists. A gate PASS is evaluation evidence only and never a command to persist business facts.
