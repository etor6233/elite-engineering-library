# Franchise Accelerator

Guía de aceleración para que **cualquier agente (Codex, Claude Code u otro)** componga una base de franquicia desde packs verificables y complete sus journeys con evidencia. El wiring local canal→Responses→tool→dominio→PostgreSQL→reply está reconstruido; TikTok Lead v1.3 dispone de retrieval autenticado, evidencia hash-linked e import durable al owner compartido, mientras cuentas/permisos live, autenticidad del webhook y evidencia productiva del proyecto continúan fail-closed. Nunca se promete producción ni un plazo menor a una semana sólo por materializar la biblioteca.

- Arranque desde cero (copiar carpeta → bridge → agente): **`START_FRANCHISE.md`**.
- Obtención/guardado de credenciales paso a paso: **`markdown_system/PROJECT_SECRETS_TEMPLATE.md`**.
- Estado y brechas verificadas: **`markdown_system/POST_DEEPSEEK_FRANCHISE_REAUDIT_2026-09-04.md`**.

## Qué hace y qué no

- **Hace**: reduce re-trabajo. El agente materializa componentes verificados, captura inputs y demuestra verticales propios de la franquicia.
- **No hace**: fabricar cuentas (LLM, WhatsApp, pasarela, ARCA), reglas fiscales/legales del país, ni aceptación de producción. Esos son inputs del proyecto y siguen `fail-closed`; inventarlos degrada la biblioteca de "élite" a "demo".

## Fast lane — 7 pasos

1. **Bridge** (raíz del proyecto, sin Git):
   ```powershell
   pwsh -NoProfile -File .\INSTALL_AGENT_BRIDGE.ps1 -ProjectRoot . -Agent Both
   ```
2. **Preflight** de toolchains (declara lo ausente, no lo finge):
   ```powershell
   pwsh -NoProfile -File .\VERIFY_EXECUTABLE_LIBRARY.ps1 -Mode Preflight
   ```
3. **Blueprint + readiness + cursor**: `PROJECT_BLUEPRINT.md` (negocio, journeys, invariantes), rondas A–H del `PROJECT-START-READINESS-VALIDATOR` y `PROJECT_EXECUTION_STATE.json`/`PROJECT_EXECUTION_EVENTS.jsonl` materializados por `EXECUTION-VALIDATOR 1.3.1`. El engineering contract 1.1.0 enlaza cada implementación con procedencia, riesgo, nueve dimensiones, tests/benchmarks y evidencia separada de artefacto/release. En brownfield son obligatorios inventario y delta; al retomar se leen sólo las referencias vencidas o necesarias.
4. **Composición de referencia**: usar `markdown_system/FRANCHISE_COMPLETE_PACK_PLAN.md` como punto de partida de 70 packs exactos/820 archivos; incluye un único runtime unido, su núcleo determinista, FinOps, resolver PostgreSQL de identidad/contacto, fence outbound 0.2.x, Google Lead Form durable, señal→reconciliation/import de Meta, retrieval→evidencia→import durable de TikTok y notificación→GET v4→candidato→respuesta aprobada→POST→confirmación/reconciliación GET-only de Mercado Libre Questions, DevSkim adaptado y release firmado portable, pero debe filtrarse por las 48 superficies del proyecto. OpenGrep/GitLab permanece fuera mientras su verifier esté rechazado. Reporting no sustituye captación, SAST no prueba seguridad y ninguna cuenta live se presume.
5. **Materializar + componer** (`materialize_markdown_pack.ps1` + `MARKDOWN_COMPOSITOR_CORE`) y **migrar** `0001`–`0049`.
6. **Vertical slice** end-to-end (captación → contacto → catálogo → cotización/pedido/pago o agenda → respuesta/handoff), checkpoint después de cada evidencia material y gates (`go test ./...`, `go vet`, `go build`, PostgreSQL real, navegador).
7. **Expandir módulos** por vertical, conservando el gate de aceptación empresarial hasta la evidencia real.

## Cobertura de 48 superficies para una franquicia de alto nivel

`R` = REQUIRED, `O` = OPTIONAL, `N` = NONE_WITH_REASON. "listo" = pack materializable y verificado; "input" = requiere un dato del proyecto.

| Superficie | Clasif. | Pack | Estado |
|---|---|---|---|
| outcome/users/journeys | R | readiness validator + blueprint | listo (inputs) |
| bounded contexts/ADR | R | blueprint + ADR | listo |
| repo/colaboración | R | políticas SCM | listo |
| contratos versionados | R | PACK_CONTRACT + clientes generados | listo |
| runtime/GC/memoria | R | Go 1.26.7 / PG 18.6 | listo |
| web pública/SEO | R | perfil web (BFF, SafeValues, CSP, Lighthouse) | listo |
| portales cliente/admin/fábrica | R | perfil web | listo |
| mobile | O | **PWA** (web responsive); nativa = `N` | recomendado PWA |
| desktop | N | — | no se justifica |
| embedded/IoT | O | telemetría de flota | `CONDITIONED` (solo electromovilidad) |
| API/backend transaccional | R | GO-ENTERPRISE-BACKEND | listo |
| dominio configurable | R | supply/commerce/fulfillment/royalty/accounting/fiscal | listo |
| identidad | R | OIDC + OpenID conformance | listo (IdP input) |
| autorización por objeto | R | scopes tenant/org + RLS | listo |
| secretos/PKI | R | vault/KMS + credential core ARCA | input (secret manager) |
| base transaccional | R | PG foundation + backup | listo |
| cache | R | GO-CACHE-CORE | listo |
| búsqueda | R | GO-SEARCH-CORE | listo |
| object storage | R | S3/GCS/Azure | input (cuenta) |
| outbox/inbox | R | Dapr/Debezium | listo |
| jobs/workflows | R | reliable workers + durable | listo |
| broker/streaming | O | Kafka (Debezium runtime) | `CONDITIONED` |
| integraciones | R | provider core + adapters | input (cuentas) |
| pagos | R | Stripe/Mercado Pago | input (sandbox) |
| marketplaces | O | Amazon/Google Merchant/Mercado Libre ingress + Questions outbound | listo (cuentas/políticas live condicionadas) |
| ads/atribución | O | Google Lead Form + Meta webhook/reconcile/import + TikTok retrieval/evidence/import + reporting Meta/Google/TikTok | tres ingress reconstruibles; cuentas, permisos, suscripciones y pruebas live siguen condicionados; TikTok inbound no se considera auténtico sin contrato provider demostrado |
| notificaciones | R | email/SMS/push/WhatsApp | `CONDITIONED`; adapters base, delivery/reconciliation live pendiente |
| data ingest | R | GO-OMNICHANNEL-LEAD-INGRESS + GO-MERCADOLIBRE-MARKETPLACE-ADAPTER + GO-LEAD-CANDIDATE-PROMOTION + adapters Meta/TikTok | Google/Meta/TikTok/Mercado Libre Questions convergen en evidencia/candidato/rechazo/outbox durables; la promoción CRM permanece separada y exige política/consentimiento |
| analytics/BI | R | GO-DATA-ANALYTICS-CORE | input (lakehouse) |
| ML/AI | R | GO-ML-AI-FOUNDATION | input (LLM provider) |
| RAG/agentes | R | GO-CONNECTED-CONVERSATION-RUNTIME + GO-APP-WIRING | wiring local V236 verificado; provider/evals/cuentas/target `CONDITIONED` |
| GPU/acceleración | O/N | serving self-hosted | según serving |
| edge/CDN/WAF | R | Cloudflare edge | input (zona) |
| contenedores | R | container packaging | listo |
| orquestación | R | k8s deploy/rollback | input (cluster) |
| IaC/cloud | R | AVM Bicep + módulos | `CONDITIONED` |
| CI | R | portable CI runner | listo |
| CD/release | R | deploy/rollback adapter | input (entornos) |
| supply chain | R | SBOM/OSV/provenance | listo |
| observabilidad | R | OTel/Prometheus | listo |
| SLO/incident | R | SLO contracts | `CONDITIONED` (producción) |
| performance | R | k6/Lighthouse | listo |
| security/appsec | R | SAST/DAST/fuzz | listo |
| privacidad/cumplimiento | R | Presidio + retención | input (jurisdicción) |
| backup/DR | R | PG backup/restore | listo |
| cost/FinOps | R | `GO-FINOPS-CORE` (tags/inventario/budget/teardown) | pack listo; conexión runtime/scan provider condicionada |
| test platform | R | matriz completa | listo |
| docs/ops | R | runbooks + evidencia | listo |

## Inyectar una capability en un proyecto existente

Ejemplo — **"agregar un chatbot a un stack que ya existe"**:

1. Vendorizar la biblioteca (`tools/elite-engineering-library/`) y ejecutar el bridge (Modo B).
2. Materializar `GO-ML-AI-FOUNDATION`, `GO-CONVERSATION-TOOLS`, `GO-LLM-ECONOMY-CORE`, `GO-HUMAN-APPROVAL-CORE`, `GO-LLM-OPENAI-ADAPTER`, `GO-OPENAI-RESPONSES-TOOL-ADAPTER`, `GO-CHANNELS-CORE`, `GO-AGENT-DOMAIN-BINDING`, `GO-CONNECTED-CONVERSATION-RUNTIME`, `GO-PG-CONTACT-CHANNEL-IDENTITY`, `GO-PG-OUTBOUND-DELIVERY-FENCE` y `GO-APP-WIRING` con sus versiones exactas del plan.
3. Implementar `ContactResolver` contra la identidad/CRM existente y configurar mappings explícitos hacia **los endpoints ya existentes** del proyecto; no crear un segundo CRM, agenda, pricing ni orders.
4. Registrar el adapter de canal live y proporcionar store PostgreSQL, token rotatorio, presupuesto y política de aprobación.
5. Demostrar canal→lead/contacto→diálogo→efecto de dominio→persistencia→reply/replay y recovery con provider real; el E2E local V236 no cierra el target.

El mismo patrón aplica a búsqueda (`GO-SEARCH-CORE`), cache (`GO-CACHE-CORE`), analytics (`GO-DATA-ANALYTICS-CORE`) o cualquier pack `REBUILD_VERIFIED`.

## Lo que debes aportar (fail-closed)

Cuentas/credenciales (LLM, WhatsApp Business, pasarela, ARCA/AFIP, cloud), reglas de negocio y fiscales del país, corpus/schemas del dominio, y la evidencia de aceptación productiva. Ninguno lo inventa el agente.

## Lo que nunca se inventa

- Reglas legales/fiscales, credenciales, proveedores o datos reales.
- `latest`/`auto` en versiones de modelo o dependencias.
- Efectos de dinero/inventario dentro del agente (siempre delega al backend determinista).
- Éxito de producción sin evidencia target.

Estado de la biblioteca: `markdown_system/MARKDOWN_SYSTEM_READINESS.md`. Reconstrucción completa: `VERIFY_LIBRARY.ps1` + `VERIFY_EXECUTABLE_LIBRARY.ps1 -Mode Audit`.
