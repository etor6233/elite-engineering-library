# Enterprise Backend Profile — Reconstruction Evidence V1

## Scope

- Plan: `markdown_system/ENTERPRISE_BACKEND_PACK_PLAN.md`
- Composition: 12 compatible packs selected by one plan
- Output: 68 materialized files plus `MATERIALIZATION_RECORD.md`
- Runtime: Go 1.26.7, Python 3.14.4, PostgreSQL 18.6

## Results

| Gate | Result |
|---|---:|
| one-command preflight and composition into an empty directory | PASS |
| migrations 0001/0002/0003 present in the output | PASS |
| fresh database creation and all three migrations | PASS |
| all SQL invariant tests | PASS |
| all Go domain/HTTP/PostgreSQL/worker tests | PASS |
| `go vet` and both executable builds | PASS |
| portable CI runner tests (6) | PASS |
| operational-readiness validator tests (6) | PASS |
| logical backup, manifest SHA-256, isolated restore and invariants | PASS |

The first profile revision omitted migrations 0001 and 0003. Inspection of the materialized output caught the omission before it was accepted; the canonical plan was corrected to include the PostgreSQL foundation and vertical schema packs, then reconstructed in a second empty directory.

## Conditions

This is the default backend profile for compatible enterprise blueprints, not a universal stack mandate. Production still depends on the selected identity, frontend, providers, legal rules and deployment platform, with their specific gates.
