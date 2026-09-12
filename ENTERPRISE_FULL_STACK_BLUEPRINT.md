# Enterprise Full-Stack Blueprint

> **Estado:** arquitectura de referencia v1; 2026-08-24.  
> **Caso conductor:** red de franquicias de motos, bicicletas y vehículos eléctricos, extensible a otros negocios multi-sucursal.  
> **Objetivo:** permitir que Codex pase de idea a vertical slice desplegable sin improvisar capas críticas.

## 1. Principios de diseño

1. Dominio y journeys antes que frameworks.
2. Monolito modular como baseline; separación física sólo con evidencia.
3. PostgreSQL es source of truth transaccional inicial.
4. Identidad no equivale a autorización de dominio.
5. Todo side effect externo es idempotente, observable y reconciliable.
6. Frontend público, portales y operación forman parte de la arquitectura.
7. Seguridad, privacidad, restore y rollback son funcionalidad del producto.
8. Los componentes públicos se adoptan por claim y gate, no por marca.

## 2. Vista de contexto

```text
prospect / customer / franchise staff / HQ / factory / supplier
                              │
             public web / customer portal / admin / mobile
                              │
                   edge + BFF/versioned API
                              │
                    modular application core
  ┌───────────────┬───────────┼───────────┬────────────────┐
identity/policy  PostgreSQL  object store  workers/outbox  search/cache optional
  │                │             │             │                 │
OIDC provider   backups/PITR  documents/media  integration adapters  derived views
                                               │
           payments / marketplaces / ads / CRM / ERP / carriers / vehicles

telemetry: application + infrastructure → OTel pipeline → metrics/traces/logs
delivery: source → CI evidence → immutable artifact → staged rollout → rollback
```

## 3. Superficies de producto

### 3.1 Web pública

- catálogo de modelos/variantes con disponibilidad por mercado;
- comparador, autonomía, batería, carga, homologación y documentos;
- landing pages por campaña/franquicia/territorio;
- buscador, SEO técnico, sitemap, structured data y canonical URLs;
- contacto, solicitud de prueba, cotización y reserva;
- consentimiento, preferencias y atribución de campaña;
- accesibilidad WCAG, Core Web Vitals y progressive enhancement;
- contenido legal, privacidad, cookies y recall/service notices.

### 3.2 Portal de cliente

- perfil, consentimientos, organizaciones y delegados;
- leads/cotizaciones/reservas/pedidos/pagos/facturas;
- vehículos/series/baterías/documentos;
- entrega, garantía, turnos, mantenimiento y casos;
- notificaciones y exportación/eliminación cuando corresponda;
- pairing/telemetría sólo mediante consentimiento y scope explícito.

### 3.3 Administración y operación

- HQ: red, territorios, franquicias, catálogos, price books y policies;
- sucursal: leads, stock, reservas, ventas, caja operacional y service;
- fábrica/proveedor: órdenes, hitos, lotes, ASN, calidad y documentos;
- logística: shipments, carriers, tracking, recepción y discrepancias;
- marketing: campañas, audiences permitidas, feeds y reconciliación;
- soporte: tickets, warranty claims, recalls y timeline auditable;
- seguridad: identidades, grants, sesiones, audit events y investigations.

### 3.4 Mobile/desktop

Sólo si el journey lo exige:

- escaneo QR/barcode/VIN/serie;
- fotos de recepción, daños o service con metadata controlada;
- firma/aceptación y prueba de entrega;
- offline bounded queue para operaciones de sucursal;
- push/deep links; secure storage; remote logout;
- distribución firmada, actualización y compatibilidad mínima.

## 4. Módulos de dominio y ownership

| Módulo | Entidades/raíces | Invariantes críticas | Publica |
|---|---|---|---|
| Organization Network | company, franchise, branch, territory, membership | una sucursal tiene ownership y vigencia; delegación es explícita | organization/branch lifecycle |
| Identity Link | principal, external identity, session metadata | no almacena password si OIDC externo; link único y auditable | identity linked/unlinked |
| Authorization | grants, relations/policies, model version | deny por defecto; checks por organización/objeto/acción | grant changed/model activated |
| Product Catalog | model, variant, option, specification, media | variante publicada es versionada; datos legales por mercado | catalog version published |
| Pricing | price book, tax class, promotion, quote | moneda/mercado/vigencia explícitos; quote es snapshot | price/quote changed |
| Lead/CRM | lead, consent, source, activity, assignment | consent purpose/version; owner y SLA visibles | lead created/qualified/converted |
| Customer | person/org, address, preference | PII minimizada y retención etiquetada | customer updated |
| Supplier/Factory | supplier, factory, agreement, capability | supplier status y terms versionados | supplier status changed |
| Procurement | purchase order, line, milestone, ASN | cantidades y tolerancias; state machine monotónica | PO issued/confirmed/shipped |
| Inventory | stock item, serial/VIN, battery, location, reservation | serial único; movimientos balanceados; reserva con expiry | stock moved/reserved/released |
| Sales/Order | cart/offer, order, line, allocation | order total reproducible; state transition validada | order placed/confirmed/cancelled |
| Payment | payment intent, attempt, refund, settlement link | dinero en minor units; webhook dedup; ledger externo conciliado | payment authorized/captured/refunded |
| Fulfillment | shipment, package, handover, proof | entrega sólo con allocation; tracking no es source of truth del order | shipment dispatched/delivered |
| Warranty/Service | policy, coverage, appointment, work order, claim | coverage evaluada sobre versión vendida; parts/labor trazables | case opened/work completed |
| Document | document, version, classification, retention | hash, owner, ACL, malware state y retention | document accepted/expired |
| Notification | template, request, delivery, preference | purpose y opt-out; no duplicar por retry | notification requested/delivered |
| Integration | connection, credential ref, cursor, mapping, sync run | secretos fuera de DB de dominio; cursor y mapping versionados | sync completed/failed |
| Audit/Compliance | audit event, legal hold, access review | append-only lógico; actor/resource/reason/correlation | evidence exported |

Cada módulo posee sus tablas y writes. Otros módulos consumen APIs internas o eventos; no escriben tablas ajenas.

## 5. Journeys verticales obligatorios

### J1 — Lead a entrega

```text
landing → consented lead → assignment → quote snapshot → reservation
→ order → payment authorization/capture → inventory allocation
→ delivery appointment → proof of delivery → warranty activation
```

Tests mínimos: atribución sin PII indebida, idempotencia de formulario, permiso por sucursal, quote reproducible, reserva concurrente, webhook duplicado, pago incierto, rollback de allocation y auditoría completa.

### J2 — Fábrica a stock disponible

```text
demand/plan → purchase order → supplier confirmation → production milestones
→ serial/lot manifest → shipment/ASN → receiving → quality hold/release
→ stock available by branch/channel
```

Tests mínimos: tolerancia de cantidades, eventos fuera de orden, ASN duplicado, serial repetido, recepción parcial, cuarentena y reconciliación.

### J3 — Publicar un modelo

```text
draft catalog → legal/technical review → media scan → market price book
→ approval → immutable catalog version → search/cache projection
→ storefront publish → marketplace/feed propagation
```

Tests mínimos: maker-checker, vigencia, rollback de versión, cache purge, SEO canonical, feed rejection y comparación de projection contra source of truth.

### J4 — Service y garantía

```text
customer/vehicle → appointment → diagnosis → work order → parts reservation
→ warranty decision → work/quality → customer acceptance → claim/reconciliation
```

Tests mínimos: ownership del vehículo, policy version, evidencia, parts balance, aprobación excepcional y privacidad.

### J5 — Onboarding de franquicia

```text
agreement → organization/territory → admins → branches → price/inventory policy
→ integrations → training/acceptance → go-live → access review
```

Tests mínimos: tenant isolation, bootstrap admin, least privilege, secret rotation, readiness gate y deprovisioning.

## 6. Arquitectura de aplicación

### 6.1 Repositorio recomendado

```text
apps/
  public-web/
  customer-portal/
  admin-portal/
  api/
  worker/
  mobile/                 # optional
packages/
  contracts/
  ui-system/
  observability/
  test-fixtures/
modules/
  organization/
  catalog/
  pricing/
  crm/
  procurement/
  inventory/
  orders/
  payments/
  fulfillment/
  service/
  documents/
  integrations/
platform/
  migrations/
  deployment/
  observability/
  policy/
docs/
  adr/
  threats/
  runbooks/
```

El lenguaje exacto se decide por equipo y restricciones. La estructura expresa ownership y contratos, no obliga a TypeScript en el core.

### 6.2 API

- HTTP/JSON versionado para superficies públicas; OpenAPI source-controlled;
- gRPC sólo para contratos internos que lo justifiquen;
- IDs opacos, pagination por cursor y filtros allowlisted;
- optimistic concurrency/version field para edición;
- `Idempotency-Key` en creates/commands con side effects;
- problem details/error codes estables; correlation/trace ID;
- deadlines, payload limits, rate limits y cancellation;
- BFF sólo si reduce acoplamiento de una superficie concreta.

### 6.3 Persistencia

- una base PostgreSQL puede alojar schemas por módulo inicialmente;
- transacción local mantiene invariantes; outbox se escribe en la misma transacción;
- migrations expand/contract y backward compatibility durante rollout;
- claves/constraints primero; índices derivados de queries y planes medidos;
- timestamps UTC + zona de negocio separada; dinero en minor units + currency;
- PII clasificada; encryption/KMS donde el threat model lo requiera;
- backup + PITR + restore drills; un backup no probado no es recovery.

### 6.4 Asincronía

- jobs en tabla/queue bounded para email, feeds y tareas cortas;
- outbox relay con claim/lease, retries limitados y dead-letter operable;
- inbox/dedup para consumers con side effects;
- Temporal sólo para procesos largos con timers/retries/compensation que superen el gate;
- Kafka/NATS sólo con throughput, fanout, replay u ownership demostrado.

### 6.5 Cache y búsqueda

- cache es prescindible y reconstruible; nunca confirma stock, pago o permiso;
- request coalescing, TTL+jitter, negative caching y invalidation metrics;
- PostgreSQL FTS/trigram primero;
- OpenSearch sólo con relevancia, volumen o faceting que lo justifique;
- projection lag y stale-state visibles al usuario/operador cuando importan.

## 7. Identidad, autorización y seguridad

### 7.1 Identidad

- OIDC Authorization Code + PKCE; MFA/passkeys según riesgo;
- sesiones cortas/rotables, refresh theft detection y logout propagado;
- service identities distintas de usuarios;
- break-glass controlado, temporal y auditado;
- account recovery tratado como flujo de alto riesgo.

### 7.2 Autorización

Modelo mínimo:

```text
principal ─membership→ organization/franchise/branch
principal ─relation→ customer/order/vehicle/document
permission = action(resource, context, model_version)
```

- checks server-side en cada resource access;
- list endpoints aplican autorización al query, no filtran después de devolver;
- model ID/version pinneado; migrations shadow-read antes de activar;
- cambios de grants mediante outbox y fail-closed;
- elegir OpenFGA, Cedar o policy local: uno como autoridad primaria.

### 7.3 Trust boundaries

- browser/mobile son no confiables;
- uploads pasan content/type/size validation, malware scan y quarantine;
- webhooks validan firma, timestamp, replay window y raw body;
- factories/suppliers tienen identities/scopes separados;
- marketplace/ads credentials viven en secret manager y por connection;
- telemetry vehicular se minimiza y separa de analytics de marketing;
- admin high-risk actions requieren step-up/reason/maker-checker cuando aplique.

## 8. Integraciones

Todo adapter implementa:

```text
canonical command/query
→ validation + authorization
→ provider mapping(versioned)
→ credential lookup
→ idempotent request with deadline
→ response/error classification
→ durable sync state
→ reconciliation
→ audit + metrics + trace
```

Nunca se modela el dominio con objetos de Stripe, Amazon, Mercado Libre, Google o Meta. Se conservan `provider_id`, version y raw evidence cifrada/retained sólo cuando sea necesario.

## 9. Observabilidad y SLO

### Señales mínimas

- RED por endpoint/worker/integration: rate, errors, duration;
- saturation: pools, queue age/depth, DB connections, CPU/memory;
- business correctness: reservation conflicts, payment uncertainty, stock mismatch, sync lag;
- security: denied checks, privilege changes, token anomalies, webhook failures;
- deploy markers y config/model versions en telemetry.

### SLO inicial a concretar

| Journey | SLI | Ejemplo de objetivo, no default contractual |
|---|---|---|
| catálogo público | successful page/API + LCP | 99.9%; p75 LCP por mercado |
| crear lead | accepted exactly once | 99.9%; p95 server < 500 ms |
| confirmar order | correct terminal/known state | 99.9%; no oversell por race |
| admin critical action | success + audit durability | 99.9%; audit correlation 100% |
| provider sync | freshness/reconciliation | lag según canal y cuota |

Los objetivos definitivos se derivan de failure cost, presupuesto y dependencia externa.

## 10. Delivery y operación

```text
pre-commit → unit/static/license/secret checks
→ integration with real PostgreSQL and provider sandboxes/mocks
→ contract + Playwright vertical journeys
→ build once + SBOM + provenance + signing
→ migration preflight
→ staging/synthetic checks
→ canary/rolling release
→ SLO and business invariant observation
→ promote or rollback
```

Runbooks mínimos:

- DB unavailable/slow, connection exhaustion y failed migration;
- queue backlog/poison message;
- payment webhook outage o state uncertainty;
- marketplace quota/credential revocation;
- identity provider outage y break-glass;
- bad catalog/price publish;
- security incident, secret rotation y evidence export;
- restore/PITR y regional/provider exit.

## 11. Build vs adopt

| Capability | Baseline | Adoptar cuando | Evitar cuando |
|---|---|---|---|
| identity | managed OIDC o Keycloak | requisitos de federation/MFA/admin lo justifican | construir password/session stack propio |
| authorization | local typed policy o OpenFGA/Cedar | relaciones/ABAC y administración lo requieren | duplicar engines |
| commerce | dominio propio modular o Saleor evaluado | fit de catálogo/order/payment alto | forzar ERP/commerce completo para pocos journeys |
| workflow | DB jobs/outbox | timers/compensation/historial complejo → Temporal | hot path o simple email |
| cache | none/local | medición muestra beneficio | source of truth/locks improvisados |
| search | PostgreSQL | relevancia/volumen/facets → OpenSearch | catálogo pequeño |
| messaging | DB queue/outbox | fanout/replay/scale → NATS/Kafka | moda arquitectónica |
| orchestration | simple container platform | multi-service/team/scale → Kubernetes | un deployable pequeño |
| edge runtime | CDN/WAF managed | isolation/custom logic → workerd | duplicar backend |

## 12. Definition of complete

Un sistema no está “completo” porque tenga pantallas y endpoints. Debe demostrar:

- [ ] J1–J5 o justificación de exclusión;
- [ ] tenant/object authorization y negative tests;
- [ ] invariantes de stock/order/payment bajo concurrencia;
- [ ] idempotencia y reconciliation de cada side effect;
- [ ] migrations, seed controlado y data retention;
- [ ] accesibilidad, responsive, SEO y budgets web;
- [ ] telemetry sin secretos/PII indebida;
- [ ] backup restore y rollback ejecutados;
- [ ] dependency/license/SBOM/provenance evidence;
- [ ] runbooks y ownership;
- [ ] manifest del Engineering Kit aprobado al nivel requerido.
