# Pack-per-claim index

Índice para elegir **un** pack por claim. Estado y condiciones: fila exacta en `CAPABILITY_CATALOG.md`. No sustituye admisión ni gates del proyecto.

## Cómo usar

1. Identifica el claim acotado (una capability o superficie).
2. Elige **una** fila de la tabla.
3. Lee **solo** ese pack y **un** mapa (`AI_ENGINEERING_MASTER_MAP.md` o `SYSTEMS_ENGINEERING_MASTER_MAP.md`).
4. Si no hay fila o el pack es `CONDITIONED`/`CANDIDATE`, escala o abre `implementation_packs/CAPABILITY_GAP_RESOLUTION_GATE.md`.

## Tabla claim → mapa → pack

| Claim | Mapa (uno) | Pack primario |
|---|---|---|
| advisories Node core/EOL | `SYSTEMS_ENGINEERING_MASTER_MAP.md` | `implementation_packs/NODE_OFFICIAL_RUNTIME_ADVISORY_GATE.md` |
| composición Markdown | `SYSTEMS_ENGINEERING_MASTER_MAP.md` | `implementation_packs/MARKDOWN_COMPOSITOR_CORE.md` |
| readiness previo al producto | `SYSTEMS_ENGINEERING_MASTER_MAP.md` | `implementation_packs/PROJECT_START_READINESS_VALIDATOR.md` |
| arquitectura empresarial (48 superficies) | `SYSTEMS_ENGINEERING_MASTER_MAP.md` | `ENTERPRISE_FULL_STACK_BLUEPRINT.md` |
| PostgreSQL transaccional | `SYSTEMS_ENGINEERING_MASTER_MAP.md` | `implementation_packs/POSTGRES_TRANSACTIONAL_FOUNDATION.md` |
| outbox / CDC / inbox | `SYSTEMS_ENGINEERING_MASTER_MAP.md` | `implementation_packs/DEBEZIUM_POSTGRES_OUTBOX_RUNTIME.md` |
| identidad y autorización | `SYSTEMS_ENGINEERING_MASTER_MAP.md` | `implementation_packs/GO_OIDC_SERVICE_TOKEN_BROKER.md` |
| web BFF / portales | `SYSTEMS_ENGINEERING_MASTER_MAP.md` | `implementation_packs/TYPESCRIPT_GO_API_WEB_BRIDGE.md` |
| SEO público (robots/sitemap/meta) | `SYSTEMS_ENGINEERING_MASTER_MAP.md` | `implementation_packs/GO_SEO_CORE.md` (+ `src/app/sitemap.ts`, `src/app/robots.ts`, `src/platform/seo/public-indexing.ts` en árbol V403) |
| calidad web automatizada (Lighthouse) | `SYSTEMS_ENGINEERING_MASTER_MAP.md` | `implementation_packs/GOOGLE_LIGHTHOUSE_WEB_QUALITY_GATE.md` |
| latencia / HTTP SLI / métricas | `SYSTEMS_ENGINEERING_MASTER_MAP.md` | `implementation_packs/GO_HTTP_METRICS_REFERENCE.md` |
| función SDR (handoff comercial) | `SYSTEMS_ENGINEERING_MASTER_MAP.md` | `markdown_system/BUSINESS_FUNCTION_OPERATING_CONTRACT_V403.md` → `implementation_packs/GO_LEAD_CANDIDATE_PROMOTION.md` (**PARCIAL** — sin pack `*SDR*`) |
| GTM / tag manager web | `SYSTEMS_ENGINEERING_MASTER_MAP.md` | — (**GAP** — sin pack; ver `markdown_system/LIBRARY_HEALTH_CHECK.md`) |
| captura QR / mobile browser | `AI_ENGINEERING_MASTER_MAP.md` | `qr_capture/` (**GAP** pack admitido; código de referencia sólo) |
| franquicia / electromovilidad backend | `SYSTEMS_ENGINEERING_MASTER_MAP.md` | `implementation_packs/GO_ENTERPRISE_BACKEND_CORE.md` |
| fiscal argentino / ARCA | `SYSTEMS_ENGINEERING_MASTER_MAP.md` | `implementation_packs/MICROSOFT_ARCA_WSFE_SOAP_ADAPTER.md` |
| pagos / marketplaces / ads | `SYSTEMS_ENGINEERING_MASTER_MAP.md` | `implementation_packs/PYTHON_AMAZON_SPAPI_CATALOG_ADAPTER.md` |
| ingesta email documental AWS | `SYSTEMS_ENGINEERING_MASTER_MAP.md` | `implementation_packs/SECURE_EMAIL_MIME_QUARANTINE_CORE.md` |
| clientes OpenAPI / WSDL | `SYSTEMS_ENGINEERING_MASTER_MAP.md` | `implementation_packs/MICROSOFT_KIOTA_OPENAPI_CLIENT_GATE.md` |
| observabilidad / CI / storage | `SYSTEMS_ENGINEERING_MASTER_MAP.md` | `implementation_packs/SECURE_OPERATIONS_DELIVERY_CORE.md` |
| deploy / canary / rollback | `SYSTEMS_ENGINEERING_MASTER_MAP.md` | `implementation_packs/SECURE_OPERATIONS_DELIVERY_CORE.md` |
| aceptación empresarial | `SYSTEMS_ENGINEERING_MASTER_MAP.md` | `implementation_packs/SECURE_OPERATIONS_DELIVERY_CORE.md` |
| privacidad / desidentificación (Presidio) | `AI_ENGINEERING_MASTER_MAP.md` | `implementation_packs/OFFICIAL_UPSTREAM_ACQUISITION_CORE.md` |
| runtime distribuido / multi-cloud | `SYSTEMS_ENGINEERING_MASTER_MAP.md` | `implementation_packs/DAPR_OFFICIAL_TRANSACTIONAL_OUTBOX.md` |
| edge / CDN / WAF | `SYSTEMS_ENGINEERING_MASTER_MAP.md` | `implementation_packs/SECURE_OPERATIONS_DELIVERY_CORE.md` |
| supply chain / licencias / OSV | `SYSTEMS_ENGINEERING_MASTER_MAP.md` | `implementation_packs/PORTABLE_SIGNED_RELEASE_EVIDENCE_GATE.md` |
| SAST / fuzz / DAST | `AI_ENGINEERING_MASTER_MAP.md` | `implementation_packs/MICROSOFT_DEVSKIM_ADAPTED_SAST_GATE.md` |
| IA / ML / RAG / agente conversacional | `AI_ENGINEERING_MASTER_MAP.md` | `implementation_packs/GO_CONVERSATIONAL_AGENT.md` |
| document intelligence empresarial | `AI_ENGINEERING_MASTER_MAP.md` | `implementation_packs/DOCUMENT_PIPELINE_ROUTING_GATE.md` |
| handoff humano durable (pg_durable) | `SYSTEMS_ENGINEERING_MASTER_MAP.md` | `implementation_packs/MICROSOFT_PG_DURABLE_HUMAN_HANDOFF.md` |
| sample Azure Document Intelligence | `AI_ENGINEERING_MASTER_MAP.md` | `implementation_packs/MICROSOFT_AZURE_DOCUMENT_INTELLIGENCE_OFFICIAL_INVOICE_SAMPLE.md` |
| Google Document AI | `AI_ENGINEERING_MASTER_MAP.md` | `implementation_packs/GOOGLE_CLOUD_DOCUMENT_AI_OFFICIAL_PROCESS_SAMPLE.md` |
| integridad evidencia (Tessera) | `SYSTEMS_ENGINEERING_MASTER_MAP.md` | `implementation_packs/TRANSPARENCY_DEV_TESSERA_POSIX_EVIDENCE_LOG.md` |
| ERP Microsoft Business Central | `SYSTEMS_ENGINEERING_MASTER_MAP.md` | `implementation_packs/MICROSOFT_BUSINESS_CENTRAL_EXECUTABLE_PLATFORM.md` |
| mobile / desktop | `SYSTEMS_ENGINEERING_MASTER_MAP.md` | — (`SPEC_ONLY`; gap gate si aplica) |
| data / analytics / search | `AI_ENGINEERING_MASTER_MAP.md` | `implementation_packs/GO_SEARCH_CORE.md` |
| cache | `SYSTEMS_ENGINEERING_MASTER_MAP.md` | `implementation_packs/GO_CACHE_CORE.md` |
| ciclo del agente / bridge | `AI_ENGINEERING_MASTER_MAP.md` | `INSTALL_AGENT_BRIDGE.ps1` + `GROK_AGENT_ENTRY.md` |
| AWS SQS idempotente | `SYSTEMS_ENGINEERING_MASTER_MAP.md` | `implementation_packs/AWS_POWERTOOLS_IDEMPOTENT_SQS_BATCH_COMPONENT.md` |
| AWS Durable Execution | `SYSTEMS_ENGINEERING_MASTER_MAP.md` | `implementation_packs/AWS_LAMBDA_DURABLE_EXECUTION_COMPONENT.md` |
| adquisición oficial sin Git | `SYSTEMS_ENGINEERING_MASTER_MAP.md` | `implementation_packs/OFFICIAL_UPSTREAM_ACQUISITION_CORE.md` |
| distribución portable | `SYSTEMS_ENGINEERING_MASTER_MAP.md` | `implementation_packs/PORTABLE_SIGNED_RELEASE_EVIDENCE_GATE.md` |
| ejecución / checkpoints | `SYSTEMS_ENGINEERING_MASTER_MAP.md` | `implementation_packs/ENGINEERING_EXECUTION_VALIDATOR.md` |
| fuzzing Go nativo | `AI_ENGINEERING_MASTER_MAP.md` | `implementation_packs/GO_NATIVE_FUZZ_GATE.md` |
| gap sin pack admitido | mapa del dominio | `implementation_packs/CAPABILITY_GAP_RESOLUTION_GATE.md` |

## Notas

- Varios claims pueden compartir un dominio; sigue siendo **un pack por claim** en cada iteración.
- Nombres de pack en catálogo pueden usar guiones; rutas aquí apuntan a archivos reales bajo `implementation_packs/`.
- Si un path no existe, busca en `implementation_packs/README.md` y escala — no inventes un pack.
