# Versioned Delivery Checklist — V138

Date: 2026-08-30  
Admission: `REBUILD_VERIFIED / CONDITIONED`

## Scope and provenance

V138 closes one concrete delivery weakness without creating a second order, inventory, identity or handover owner. The implementation is `AUTHORED`; it is not presented as Microsoft code. Microsoft Field Service inspections govern published/versioned inspection definitions and required responses, while Microsoft Business Central posting governs the durable delivery boundary:

- https://learn.microsoft.com/en-us/dynamics365/field-service/inspections-overview
- https://learn.microsoft.com/en-us/dynamics365/business-central/ui-post-sales
- https://www.postgresql.org/docs/18/ddl-constraints.html
- https://www.postgresql.org/docs/18/explicit-locking.html

The project must still approve the actual checklist content, evidence retention and legal acceptance/signature policy. This pack does not claim that a technical confirmation is a qualified electronic signature.

## Materialized contract

- `GO-FRANCHISE-CUSTOMER-JOURNEY-API 0.6.0`: 27 files; pack SHA-256 `22f62ca4b8d15a6b1fe769a62d4a888fa95a9af3a2634c61306622b476ebebab`.
- `TS-FRANCHISE-JOURNEY-PORTALS 0.7.0`: 14 files; pack SHA-256 `900cbc773175627229a9bfd8f97ac07de36037129213e9d4df91fc0b182ef6ef`.
- backend composition: 24 packs / 241 files; record SHA-256 `ddb33b575e0c942d2f0cf261fde33b45d2e2223678139dd38b5f818af7de481f`.
- web composition: 6 packs / 84 files; record SHA-256 `2b3cdcde88e2f8573f112f3038519ba724be1f5b0479af93f042dcb8c670141f`.

New PostgreSQL files:

| File | Bytes | SHA-256 |
|---|---:|---|
| `0016_versioned_delivery_checklist.up.sql` | 8,859 | `67fcc7ac5e5e11bc57f648f1bd5cf87d2e2da174e4be07eae63820fea953098e` |
| `0016_versioned_delivery_checklist.down.sql` | 1,294 | `c31137a931f72b765b1585863cb08257bf6f77a1c1e6c722e6a7a87aabfe3d15` |
| `0016_versioned_delivery_checklist.test.sql` | 1,117 | `13181c3d6eb1daef8d913c3cb178ee55d76bfaa1e20e699054c3eb40f6a4e550` |

## Demonstrated behavior

1. Publication inserts a draft definition plus 1–64 ordered items and publishes them in one transaction; published definition/items are immutable.
2. Responses bind tenant, organization, prepared handover, exact checklist ID/version, item, verified operator and timestamp; response rows are append-only.
3. PostgreSQL rejects presentation while a required item lacks a response and validates confirmation/evidence response forms.
4. Completion changes `prepared → presented`, fixes checklist ID/version/completion actor/time, increments the handover version and writes an outbox event in the same transaction.
5. Customer acceptance requires the verified customer, expected handover version, exact completed checklist ID/version and durable stock serial; the server hashes the exact command and atomically writes the acceptance event.
6. The customer portal displays the checklist version and controls and does not enable acceptance without server-reported completion.

## Reconstruction and gates

- Individual packs materialized 27/27 and 14/14; all 15 changed/new paths were byte-identical to their tested author trees.
- Clean backend profile composed 24/241. `go test ./...`, `go vet ./...` and four command builds passed.
- A fresh PostgreSQL 18.6 database applied migrations 0001–0016 with `ON_ERROR_STOP=1`; SQL structure test and the full repository integration passed from the Markdown composition.
- Negative integration paths passed: acceptance before completion, missing required response, wrong checklist version, wrong durable serial and cross-customer scope.
- Clean web profile composed 6/84; `pnpm install --frozen-lockfile --offline` downloaded zero packages, strict TypeScript passed, nine Vitest files/36 tests passed, Next.js 16.3.2 production build and license report passed.
- The temporary PostgreSQL server was explicitly stopped after both database runs.

## Learned failures

`LIB-FAIL-1372` through `LIB-FAIL-1379` record the read-path error, atomic patch rejections, two PostgreSQL supervisor behaviors, Windows glob error, split build window and the first undersized Audit window. None was hidden or treated as PASS. Their corrected repetitions are the gates above, including the extended `VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit` over 83 packs and 121 sources.

## Remaining condition

This block is reusable technical infrastructure. A project is not production-ready until its owner approves the real checklist/version content, operator roles, customer wording, retention/signature policy and proves real IdP, backend, browser/accessibility, edge, load, security, deployment and recovery evidence.
