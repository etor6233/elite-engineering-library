# Architecture Packs Index

> **Regla:** este índice informa estado; no promueve ningún pack. Sólo `REUSABLE_PACK` puede proponerse para incorporación inmediata después de comprobar compatibilidad con el proyecto.

Estos documentos son architecture packs, no implementation packs. La migración a código reconstruible se gobierna mediante `../markdown_system/PACK_CONTRACT.md`; el estado real está en `../markdown_system/CAPABILITY_CATALOG.md`.

| Pack | Estado | Cubre | Bloqueo principal |
|---|---|---|---|
| `FRANCHISE_FULL_STACK_FOUNDATION.md` | `CANDIDATE_PACK` | composición, slices J1–J5, gates full-stack | starter concreto + harness + piloto |
| `POSTGRES_TRANSACTIONAL_MODULAR_MONOLITH.md` | `CANDIDATE_PACK` | schemas, idempotency, outbox, inbox, jobs, migrations, restore | implementación por stack + concurrency/recovery evidence |
| `PUBLIC_WEB_ADMIN_PORTALS.md` | `CANDIDATE_PACK` | public/customer/admin, accessibility, security, performance, Playwright | framework/deployment selection + executable UI harness |
| `IDENTITY_AUTHORIZATION_KEYCLOAK_OPENFGA.md` | `CANDIDATE_PACK` | OIDC + fine-grained multi-organization authorization | upstream code-path audit + model migration/load/security harness |
| `INTEGRATION_PAYMENTS_MARKETPLACES_ADS.md` | `CANDIDATE_PACK` | payments, Amazon/Mercado Libre, Google/Meta/TikTok Ads, webhooks/reconciliation | provider/market terms + sandbox/contract evidence |
| `OBSERVABILITY_SECURE_DELIVERY.md` | `CANDIDATE_PACK` | OTel/metrics/traces/logs, CI, SBOM, rollout, rollback, restore | instantiated pipeline + induced-failure/recovery evidence |

## Orden recomendado para un proyecto

```text
FRANCHISE_FULL_STACK_FOUNDATION
→ POSTGRES_TRANSACTIONAL_MODULAR_MONOLITH
→ IDENTITY_AUTHORIZATION_KEYCLOAK_OPENFGA or a justified alternative
→ PUBLIC_WEB_ADMIN_PORTALS
→ INTEGRATION_PAYMENTS_MARKETPLACES_ADS as providers enter scope
→ OBSERVABILITY_SECURE_DELIVERY from Slice 0, not at the end
```

## Regla de selección

Un proyecto crea `PROJECT_PACK_SELECTION.md` con:

```yaml
pack: ""
version_or_commit: ""
status_at_selection: ""
applicable_claims: []
non_applicable_sections: []
project_overrides: []
license_evidence: []
tests_and_harness: []
open_gates: []
owner: ""
review_date: ""
```

No se copia un pack completo cuando sólo aplica una capability. Toda desviación material se registra como ADR y no se modifica el significado del estado de admisión.
