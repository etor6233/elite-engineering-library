# Total System Capability Contract

Este contrato obliga al agente a evaluar toda la infraestructura de programación sin instalar piezas innecesarias. Cada capability termina como `REQUIRED`, `OPTIONAL`, `NONE_WITH_REASON` o `BLOCKED`; sólo `REQUIRED` entra al plan de materialización.

## Registro mínimo por capability

```yaml
capability_id: ""
classification: REQUIRED|OPTIONAL|NONE_WITH_REASON|BLOCKED
trigger_or_reason: ""
authority: []
pack_or_authored_work: []
configuration_keys: []
source_of_truth_and_owner: ""
trust_and_failure_boundaries: []
licenses_and_provenance: []
tests_and_evidence: []
deployment_rollback_recovery: []
connected_journey_ids: []
open_conditions: []
```

Una capability no se considera cubierta porque exista un framework, SaaS o manual. Debe tener owner, archivos materializables, configuración, pruebas, operación, licencia y salida/reemplazo.

Si una capability `REQUIRED` no tiene un pack compatible admitido, aplicar `CAPABILITY_GAP_RESOLUTION_PACK_PLAN.md`. El expediente ejecutable debe enlazar requirement y autoridad, registrar búsquedas oficiales actuales, fijar candidatos/bytes/revisión/licencia/hash/claim, aplicar G0–G8 y devolver la decisión al owner canónico. `RESEARCH_INCOMPLETE` obliga a continuar; sólo `USE_REUSABLE_PACK` habilita implementación. Un bloqueo honesto no cierra la capability y nunca autoriza inventar código, exactitud o procedencia.

## Cierre conectado obligatorio

Toda capability `REQUIRED` debe pertenecer a por lo menos un journey release-bound validado por `PROJECT-START-READINESS-VALIDATOR`. El journey puede ser humano, API, CLI o automatización, pero siempre identifica persona/actor y rol responsable. Conserva evidencia separada de interfaz, autorización, contrato, dominio, dato o efecto externo, auditoría/observabilidad, respuesta y errores comprensibles, ayuda contextual, capacitación, soporte, impacto de actualización, rollback/recovery, E2E y negativos. Ayuda, capacitación y soporte usan la misma versión que la release; una modificación del recorrido reabre esas evidencias.

Una pantalla sin efecto durable, un backend sin interfaz operable, una capacitación desactualizada o un soporte sin contexto autorizado no cierran la capability. Un archivo puede servir como evidencia de varios eslabones sólo cuando contiene explícitamente cada contrato; compartir un path no permite inferir contenido ausente.

## Inventario obligatorio

| ID | Superficie | Se activa cuando | Cierre material mínimo |
|---|---|---|---|
| `PRD-INTAKE` | outcome, usuarios, journeys, aceptación | siempre | blueprint, requirements↔tests, supuestos y blockers |
| `ARCH-DOMAIN` | bounded contexts, invariantes, ownership | siempre | vistas, ADR, source of truth, evolución y rollback |
| `REPO-SCM` | repositorio y colaboración | siempre | estructura, ignore/attributes, branch/review policy y ownership |
| `CONTRACTS` | API/events/config/schemas | toda frontera | contratos versionados, compatibilidad, generated artifacts y contract tests |
| `RUNTIME` | lenguaje/runtime/GC/memoria | toda ejecución | versión fijada, support window, build y profiling gates |
| `WEB-PUBLIC` | adquisición/catálogo/contacto | interfaz pública | SSR/SEO/consent/a11y/performance/security y E2E |
| `WEB-PORTALS` | cliente/admin/partner/fábrica | operación humana autenticada | journeys completos, sesiones, authz por recurso y E2E negativos |
| `MOBILE` | lifecycle/offline/push/hardware | app móvil requerida | storage seguro, sync, firma, updates, permisos y device tests |
| `DESKTOP` | integración OS/offline/periféricos | desktop requerido | sandbox/installer/update/signing/crash recovery |
| `EMBEDDED-IOT` | vehículo/dispositivo/firmware | hardware gestionado | boot/RTOS, device identity, secure OTA, rollback, fleet telemetry y safety |
| `API-BACKEND` | comandos/queries externos | sistema transaccional | strict HTTP/RPC, limits, idempotencia, errors, deadlines y tests |
| `DOMAIN-MODULES` | negocio configurable | siempre que haya reglas | state machines, invariantes, permissions, migrations y property tests |
| `IDENTITY` | usuarios o servicios autenticados | casi toda producción | issuer/audience/keys, session/token lifecycle, MFA/federation y negatives |
| `AUTHORIZATION` | acceso por tenant/objeto | datos o acciones protegidos | policy versionada, deny-by-default, tenant isolation y migration tests |
| `SECRETS-PKI` | credenciales/cifrado/TLS | toda producción | vault/KMS/PKI adapter, rotación, envelope encryption y break-glass auditado |
| `TX-DATABASE` | source of truth transaccional | writes consistentes | schema/constraints/indexes/transactions/migrations/concurrency/restore |
| `CACHE` | hot data/latencia/descarga | workload lo demuestra | keys/TTL/invalidation/eviction/stampede/fallback y consistency tests |
| `SEARCH` | full text/vector/facets | journeys lo requieren | indexing ownership, freshness, ACL filtering, rebuild y relevance eval |
| `OBJECT-STORAGE` | archivos/media/modelos/backups | blobs requeridos | naming, checksum, encryption, lifecycle, malware/content gate y restore |
| `OUTBOX-INBOX` | side effects/eventos | commit + publish separados | claim/lease, dedup, schema, replay, poison handling y crash tests |
| `JOBS-WORKFLOWS` | trabajo durable/largo | async/retry/human steps | lease, retry budget, timeout, compensation, history y operator tooling |
| `BROKER-STREAMING` | fanout/throughput/replay | Postgres no alcanza | partition/key/order, retention, backpressure, consumer lag y DR |
| `INTEGRATIONS` | APIs/webhooks de terceros | provider externo | adapter, secret isolation, deadlines, idempotencia, signatures, reconciliation |
| `PAYMENTS` | movimiento de dinero | cobros/reembolsos | provider ledger, webhook verification, reconciliation, disputes y sandbox |
| `MARKETPLACES` | catálogo/orden externo | canal requerido | mapping/version skew, quotas, polling/webhooks, reconciliation y ownership |
| `ADS-ATTRIBUTION` | campañas/medición | marketing requerido | consent, offline conversions, attribution limits, deletion y spend controls |
| `NOTIFICATIONS` | email/SMS/push/chat | comunicación requerida | templates/versiones, preferences, suppression, provider fallback y receipts |
| `DATA-INGEST` | CDC/batch/stream | analytics/ML | contracts, quality, lineage, replay, late data y source ownership |
| `ANALYTICS-BI` | métricas/decisiones | reporting requerido | semantic metrics, warehouse/lakehouse, freshness, access, cost y reconciliation |
| `ML-AI` | predicción/generación | valor medible | data/model versions, evals, safety, serving SLO, rollback y monitoring |
| `RAG-AGENTS` | retrieval/tools/autonomía | evidencia lo justifica | corpus ACL, provenance, injection defense, tool policy, eval y human controls |
| `GPU-ACCEL` | throughput/tails/model size | benchmark lo exige | hardware matrix, kernels, precision/correctness, memory y fallback |
| `NETWORK-EDGE` | DNS/TLS/proxy/CDN/WAF | todo deployment remoto | trust proxy, certificates, timeouts, limits, health y regional failover |
| `CONTAINERS` | artifact portable | runtime containerizado | minimal image, non-root, read-only, health, limits, digest y scan |
| `ORCHESTRATION` | scheduling/HA/scale | topology lo exige | rollout, probes, budgets, autoscaling, disruption, quotas y rollback |
| `IAC-CLOUD` | infraestructura reproducible | entornos compartidos | state/locks, modules, policy, drift, imports, destroy guard y recovery |
| `CI` | integración automática | siempre | formatting, unit/contract/security/license/build gates con cache segura |
| `CD-RELEASE` | promoción/deploy | entrega fuera de local | immutable artifact, environments, approvals por riesgo, canary y rollback |
| `SUPPLY-CHAIN` | dependencias/artifacts | siempre | pins, SBOM, provenance, signing, secrets scan, vulnerability/license policy |
| `OBSERVABILITY` | detectar/diagnosticar | toda producción | OTel/logs/metrics/traces, correlation, redaction, dashboards y retention |
| `SLO-INCIDENT` | confiabilidad operada | production/regulated | SLIs/SLOs, alerts accionables, runbooks, on-call, postmortem y error budget |
| `PERFORMANCE` | latency/capacity/cost | siempre según riesgo | workload, correctness gate, p50–p99.9, goodput, profiling, overload y budget |
| `SECURITY-APPSEC` | threat/abuse/supply chain | siempre | threat model, secure defaults, SAST/DAST/fuzz, dependency/secrets y regressions |
| `PRIVACY-COMPLIANCE` | PII/regulación/jurisdicción | datos regulados | inventory, purpose/consent, retention/deletion/export, audit y DPIA controls |
| `BACKUP-DR` | datos/servicio recuperables | toda producción con estado | backup, PITR, restore verification, RPO/RTO, failover y DR drill |
| `COST-FINOPS` | gasto externo | cloud/SaaS/AI/providers | tags/budgets/unit cost/alerts/caps, capacity plan y exit cost |
| `TEST-PLATFORM` | evidencia transversal | siempre | unit/property/contract/integration/E2E/a11y/load/security/recovery matrix |
| `DOCS-OPS` | transferencia/operación | siempre | generated docs, ADR, threat/failure models, runbooks y versioned evidence |

## Reglas de dependencia

- `WEB-PORTALS`, `MOBILE`, `DESKTOP`, `API-BACKEND` o integrations protegidas implican `IDENTITY` y `AUTHORIZATION`.
- todo write con side effects remotos implica `OUTBOX-INBOX` o una alternativa durable demostrada.
- `TX-DATABASE` en producción implica `BACKUP-DR`, `OBSERVABILITY`, `SECURITY-APPSEC` y `PRIVACY-COMPLIANCE` evaluada.
- `ML-AI` implica `DATA-INGEST`, model/data/eval governance, seguridad IA, observabilidad y rollback.
- `RAG-AGENTS` implica `ML-AI`, ACL end-to-end, tool authorization, prompt-injection tests y budgets.
- `CONTAINERS`, `IAC-CLOUD` o `CD-RELEASE` implican `SUPPLY-CHAIN`, signing/provenance y rollback.
- `PAYMENTS`, `MARKETPLACES` y `ADS-ATTRIBUTION` siempre implican reconciliation y provider sandbox evidence.
- `EMBEDDED-IOT` implica device identity, secure OTA, rollback, fleet observability y threat/safety model.

## Definition of closure

El proyecto puede comenzar a expandir features cuando todas las capabilities están clasificadas, cada `REQUIRED` está cubierta por journeys conectados y el primer vertical slice selecciona journeys demostrados end-to-end. Puede llamarse listo para producción sólo cuando no hay `BLOCKED` crítico y cada capability requerida enlaza artifacts, experiencia operable, autorización, efectos, auditoría, ayuda/capacitación/soporte versionados, pruebas ejecutadas, licencia/procedencia, deployment, rollback y recovery.
