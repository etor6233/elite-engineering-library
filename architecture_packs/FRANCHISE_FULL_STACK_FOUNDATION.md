# Architecture Pack — Franchise Full-Stack Foundation

> **Estado:** `CANDIDATE_PACK`.  
> **Autoridades:** `ENTERPRISE_FULL_STACK_BLUEPRINT.md`, `SOFTWARE_ARCHITECTURE_SYSTEM_DESIGN.md`, `FRONTEND_PRODUCT_ENGINEERING_UX.md`, `SOFTWARE_BACKEND_API_ENGINEERING.md`, `DATABASE_STORAGE_INTERNALS.md`, `SECURITY_SRE_CLOUD_INFRASTRUCTURE.md`.  
> **Promoción pendiente:** instanciación en stack concreto, harness J1–J5, threat model y evidencia G3–G8.

## 1. Claim estrecho

Proporciona el esqueleto de decisión y entrega para construir una plataforma multi-organización con web pública, portales, API, dominio transaccional, workers, integraciones y operación. No aporta reglas fiscales, contables, regulatorias ni comerciales de un país específico.

## 2. Inputs obligatorios

```yaml
markets: []
languages: []
currencies: []
organization_model:
  hq: true
  franchisees: true
  branches: true
  factories_suppliers: true
channels:
  public_web: true
  customer_portal: true
  admin_portal: true
  mobile: conditional
journeys:
  lead_to_delivery: required
  factory_to_stock: required
  catalog_publish: required
  service_warranty: required
  franchise_onboarding: required
quality:
  availability: TBD
  latency: TBD
  rpo: TBD
  rto: TBD
  privacy_jurisdictions: []
delivery:
  team: []
  cloud_or_on_prem: TBD
  deadline: TBD
```

Valores `TBD` impiden declarar producción, pero no impiden crear el vertical skeleton.

## 3. Salidas que Codex debe crear

```text
PROJECT_AUTHORITY_MAP.md
docs/journeys/J1_LEAD_TO_DELIVERY.md
docs/journeys/J2_FACTORY_TO_STOCK.md
docs/journeys/J3_CATALOG_PUBLISH.md
docs/journeys/J4_SERVICE_WARRANTY.md
docs/journeys/J5_FRANCHISE_ONBOARDING.md
docs/adr/0001-architecture-baseline.md
docs/adr/0002-identity-authorization.md
docs/threats/system-threat-model.md
docs/runbooks/
engineering_execution_kit/project.json
apps/ + modules/ + packages/ + platform/
```

## 4. Arquitectura baseline

```text
one repository
one modular application deployable + one worker deployable
three web surfaces sharing contracts/design primitives, not authorization state
one PostgreSQL cluster/database initially
one OIDC provider
one authorization authority
one object-store adapter
one outbox relay
one telemetry pipeline
```

No agregar microservices, distributed cache, search cluster, event log o workflow engine hasta que un ADR demuestre el driver.

## 5. Fronteras de frontend

| Superficie | Audiencia | Datos sensibles | Rendering/default |
|---|---|---|---|
| public-web | anónimo/prospect | mínimo; consented lead | SSR/static + progressive enhancement |
| customer-portal | customer/delegate | PII, orders, vehicles, service | authenticated web app; no secrets in browser |
| admin-portal | branch/HQ/factory/support | operational/PII/high-risk actions | authenticated app + server-side authorization every request |
| mobile | staff/customer if justified | camera/offline/secure tokens | native lifecycle and bounded sync |

Compartir design tokens y primitive components. No compartir stores globales, session assumptions ni rutas privileged entre public y admin.

## 6. Fronteras del core

Cada módulo expone:

```text
commands       mutate under transaction and authorization
queries        return authorized projections
events         facts after commit through outbox
contracts      versioned types, errors and compatibility rules
tests          invariants + negative authorization + concurrency
```

No se permite acceso cross-module directo a repositorios internos. Para lectura agregada, usar query composition/projections; para writes, command contract.

## 7. Orden de implementación

### Slice 0 — plataforma verificable

- health/readiness separados;
- config schema y secret references;
- DB migration runner con lock/timeout;
- request ID/trace context y structured logging;
- auth middleware stub fail-closed;
- CI unit/integration/license/secret/build;
- immutable artifact y local compose/dev environment.

### Slice 1 — organización + autorización

- organization/franchise/branch/membership;
- OIDC login y principal link;
- resource authorization con negative tests;
- audit event por membership/grant change;
- admin UI mínima para onboarding.

### Slice 2 — catálogo público + lead

- draft/publish catalog version;
- public listing/detail/SEO;
- lead form con consent version, rate limit y idempotency;
- branch assignment; notification outbox;
- Playwright: landing → lead → admin visible.

### Slice 3 — stock + order + payment sandbox

- serialized stock item/reservation;
- quote snapshot/order state machine;
- provider-neutral payment intent;
- signed webhook ingestion/dedup/reconciliation;
- Playwright/API concurrency: two buyers, one unit.

### Slice 4 — procurement + fulfillment

- PO/ASN/receiving/quality state machines;
- shipment/delivery/proof;
- factory/supplier scopes;
- out-of-order and duplicate contract tests.

### Slice 5 — service/warranty + integrations

- vehicle ownership/coverage/work order;
- parts reservation and claim evidence;
- marketplace/feed and ads adapters;
- provider outage/degraded-mode tests.

## 8. Feature flags y releases

- flags no reemplazan autorización;
- server evaluates high-risk flags; browser recibe sólo resultado necesario;
- cada flag tiene owner, expiry y removal issue;
- schema changes usan expand/contract antes de enable;
- external integrations habilitadas por connection/tenant/market;
- rollout: internal → pilot branch → percentage/market → full.

## 9. Datos de prueba

Fixtures deterministas:

```text
HQ_ACME
FRANCHISE_NORTH / BRANCH_NORTH_01
FRANCHISE_SOUTH / BRANCH_SOUTH_01
CUSTOMER_ALICE / CUSTOMER_BOB
SUPPLIER_CELL / FACTORY_ASSEMBLY
MODEL_E1 / VARIANT_E1_CITY
SERIAL_E1_0001 / BATTERY_B_0001
```

Tests de aislamiento siempre intentan leer/mutar recursos de otra franquicia. Ningún fixture contiene persona, VIN, token o documento real.

## 10. Gates de aceptación

- [ ] manifests pasan plan y luego implementation/release;
- [ ] cinco journeys trazados requirement → code → test → telemetry → runbook;
- [ ] cross-tenant negative suite;
- [ ] concurrency sobre reservation/order/payment;
- [ ] duplicate/out-of-order provider event suite;
- [ ] accessibility + keyboard + responsive + browser matrix;
- [ ] PII/secret telemetry scan;
- [ ] migration forward/backward compatibility;
- [ ] backup restore y rollback drill;
- [ ] SBOM, provenance y license evidence;
- [ ] SLOs y alerts basadas en usuario/business correctness.

## 11. Rechazos automáticos

- tenant ID aceptado del browser sin verificar membership;
- roles globales como única autorización;
- stock calculado desde cache/search;
- webhook que ejecuta side effect antes de persistir/deduplicar;
- pago marcado exitoso sólo por redirect del browser;
- marketplace como source of truth de catálogo/order interno;
- credentials por usuario guardadas en texto;
- admin y public app compartiendo privileged client bundle;
- `latest` o branch head como pin de producción;
- “production-ready” sin restore, rollback y negative security tests.

## 12. Evidencia para promoción

Para llegar a `REUSABLE_PACK` se necesita un starter concreto que ejecute Slice 0–3 como mínimo, licencia de cada dependencia, generated SBOM, tests reproducibles en CI, deployment local/staging, runbooks de fallos inducidos y una adopción piloto que no contradiga los manuales autoridad.
