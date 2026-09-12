# Reauditoría post-DeepSeek de franquicia y captación omnicanal — 2026-09-04

## Decisión vigente

```yaml
library_structural_integrity: PASS
generic_reuse_for_other_projects: YES_CONDITIONED
franchise_foundation_materialization: PASS
franchise_connected_runtime: PASS_LOCAL_CONDITIONED
production_ready_project: NOT_CLAIMED
open_critical_failures: []
```

La biblioteca puede acelerar proyectos distintos porque los packs separan contratos, dominio, adapters, datos, web, seguridad, operación y evidencia. No autoriza copiar todos los packs a cualquier proyecto: las 48 superficies de `TOTAL_SYSTEM_CAPABILITY_CONTRACT.md` se clasifican por proyecto y sólo las `REQUIRED` entran al plan.

Actualización vigente V258 (la evidencia V251 de las secciones siguientes se conserva como historia): la composición selecciona 67 packs/704 archivos y materializa un único runtime conversacional conectado con su núcleo determinista, FinOps, resolver PostgreSQL de identidad/contacto, fence outbound 0.2.x, ingreso Google/Meta/TikTok y Mercado Libre Questions v4 al owner durable común, más respuesta Mercado Libre aprobada con POST/confirmación/reconciliación sin retry y release firmado portable; **todavía no certifica un proyecto productivo**. La biblioteca raíz gobierna 160 packs/1.395 archivos/51 perfiles y agrega un resolver opt-in para cualquier capability requerida sin cobertura; fuzz Go continúa opt-in porque sus targets e invariantes se derivan del código real. OpenGrep/GitLab queda preservado como evidencia, bloqueado y fuera del perfil ejecutable por el SCA actual de Cosign; DevSkim adaptado es el lane seleccionado. El cursor verificable del agente permite comenzar o retomar sin duplicar planes ni releer todo el corpus. Los nombres `FRANCHISE_COMPLETE_PACK_PLAN` y `FRANCHISE_SERVERLESS_PACK_PLAN` no deben interpretarse como aceptación live: cuentas, providers, reglas, infraestructura y seguridad se demuestran en cada target.

## Evidencia reproducida

| Comprobación | Resultado | Alcance real |
|---|---|---|
| inventario raíz post-V251 | 156 packs, 1.371 archivos materializables | suma resolver de gaps 4/4, conserva Kiota OpenAPI→Go, Go fuzz, execution validator, TikTok Lead y lifecycle durable; la auditoría raíz gobierna hashes, no accesos live del proyecto |
| control del agente V240 | `NEW`/`EXISTING`, baseline/delta, tasks, contexto, evidencia y checkpoints | cursor materializable y verificable; no duplica readiness/Spec Kit/fallos ni convierte blockers live en PASS |
| perfil completo vigente | 63 packs únicos / 675 archivos | incluye un runtime, núcleo determinista, FinOps, resolver, fence outbound, endpoint Meta y retrieval/import TikTok; ninguna cuenta live se inventa |
| perfil serverless | 45 packs, 537 archivos | agrega LLM/economía/canales/SMS/FinOps; no agrega tools, domain binding, aprobación ni app wiring |
| fixture serverless + cuatro packs faltantes | 17 archivos añadidos sin colisión; `go test ./...` PASS | demuestra compatibilidad de compilación, no recorrido conductual ni durabilidad |
| `GO-APP-WIRING 0.2.0` + `GO-CONNECTED-CONVERSATION-RUNTIME 0.1.0` | `REBUILD_VERIFIED / CONDITIONED` | canal→presupuesto→Responses→aprobación→scope dinámico→dominio→PostgreSQL→reply/replay probado |
| `GO-PG-CONTACT-CHANNEL-IDENTITY 0.1.0` | `REBUILD_VERIFIED / CONDITIONED` | HMAC/evidencia/CAS/revocación y resolver PostgreSQL real; E2E lead→conversación→cotización/replay probado |
| `GO-PG-OUTBOUND-DELIVERY-FENCE 0.2.0` + `GO-MERCADOLIBRE-QUESTION-OUTBOUND 0.1.0` | `REBUILD_VERIFIED / CONDITIONED` | un intento provider automático, rechazo terminal evidenciado, ambigüedad inmovilizada, POST aprobado y reconciliación GET-only |
| `GO-META-LEAD-WEBHOOK-SIGNAL 0.1.0` | `REBUILD_VERIFIED / CONDITIONED` | challenge/firma raw-body, batch completo, PostgreSQL antes del ACK y outbox hacia recuperación; live condicionado |
| `PYTHON-TIKTOK-LEAD-ADAPTER 0.2.0` + `GO-TIKTOK-LEAD-DURABLE-IMPORT 0.1.0` | `REBUILD_VERIFIED / CONDITIONED` | wheel oficial hash-pinned, retrieval v1.3, artefactos hash-linked y verificación/import al owner PostgreSQL/outbox; webhook no autenticado y live condicionados |
| Audit integral post-V251 online | `PASS`: 156 packs/121 fuentes/16 adapters/gap resolver 1/Kiota gate 1/Go fuzz gate 1 | incluyó resolver 3 positivos/6 negativos, adquisición/generación/compilación/SCA Kiota, baseline/semilla/fuzz Go real y runtimes documentales; no convierte gates live del target en PASS |

Un PASS de compilación no sustituye un journey. El gate pendiente debe demostrar:

```text
provider verificado
  -> raw recibido y confirmado durablemente
  -> deduplicación/replay/divergencia
  -> lead/contacto/atribución tenant-scoped
  -> binding exacto y evidenciado de identidad de canal
  -> conversación tenant-scoped
  -> política de consentimiento y respuesta
  -> respuesta económica gobernada
  -> tool autorizada
  -> cita/reserva/cotización/pedido durable real
  -> auditoría + outcome + reconciliación/postback
  -> error comprensible, handoff y recuperación
```

## Estado de hallazgos

### 1. Wiring del agente

V236 reemplazó el agente paralelo por un único `ConversationRuntime`, exige store/resolver/canal/presupuesto y conecta `Dispatcher` con Responses, safety, aprobación, gateway autenticado, scopes por contacto y PostgreSQL. La prueba E2E repitió el mismo message ID y demostró dos replies con una sola llamada de dominio y sin segunda inferencia. `LIB-FAIL-1743` queda corregido; FinOps conserva su claim separado de inventario/costo de infraestructura y no se finge como ledger distribuido de tokens.

V237 cerró la implementación faltante del resolver: `GO-PG-CONTACT-CHANNEL-IDENTITY` vincula una identidad exacta a un lead existente mediante HMAC tenant/canal, policy/evidencia, versión esperada y ledger append-only. El lead/subject no puede reasignarse silenciosamente y una revocación hace que el runtime falle a handoff. El E2E ya no usa un resolver ficticio para el nuevo recorrido lead→cotización; propiedad de cuenta y método de verificación siguen siendo evidencia del proyecto.

V238 cerró el retry directo de la respuesta: `GO-PG-OUTBOUND-DELIVERY-FENCE` registra el intento antes del proveedor, valida receipt, retorna replay sin segundo send y convierte timeout/lease/commit incierto en `unknown`. Sólo una reconciliación específica puede cerrar el estado; no se afirma exactly-once delivery. El E2E local demuestra una sola llamada sender ante dos dispatch.

### 2. Captación omnicanal

V224 aporta el owner PostgreSQL para raw/candidato/outbox de Google y V225 promoción explícita a CRM. V226 agrega recuperación/reconciliación Meta por SDK oficial exacto; V227 verifica sus tres artefactos y persiste candidato o rechazo en el mismo owner con replay. V239 agrega el endpoint Meta firmado: guarda lote/signal y emite outbox antes del ACK. V242–V245 extienden esa semántica a TikTok mediante retrieval autenticado, respuesta SDK preservada, lote candidato/rechazado, receipt hash-linked y verificación/import estricto antes de escribir en el mismo owner. Falta demostrar ambos recorridos con cuentas live y, para TikTok inbound, un mecanismo oficial de autenticidad. El cierre requiere:

- identidad `(tenant, provider, provider_event_id)` y hash de payload;
- recepción cruda inmutable antes de confirmar éxito;
- duplicado exacto idempotente y duplicado divergente fail-closed;
- normalización versionada a lead/contact/conversation/attribution;
- inbox/outbox durable, replay, poison queue, reconciliación y gaps explícitos;
- consentimiento, finalidad, retención, borrado/exportación y ventanas de contacto;
- freshness/correctness SLO por fuente;
- resultado comercial y feedback de conversión separado y aprobado.

### 3. Adapters por proveedor

| Superficie | Lo oficial demostrado | Estado permitido |
|---|---|---|
| Google Ads Lead Form | webhook JSON, `google_key`, `lead_id` para dedupe, campos desconocidos ignorados, 200/4xx/5xx y recurso de submissions | parser/handler puede implementarse y probarse contra el contrato oficial |
| Meta Lead Ads | dos repos oficiales fijados gobiernan challenge, HMAC SHA-256 y envelope leadgen; SDK fijado expone `LeadgenForm.get_leads`/`get_test_leads`; V239 persiste señal/outbox, V226 recupera y V227 importa | cadena offline reconstruible; cuenta, permisos, Page/form, suscripción, test lead y entrega live siguen condicionados |
| TikTok Lead Generation | referencia oficial vigente v1.3 para GET `/lead/get/` y POST `/subscription/subscribe/`; wheel oficial `tiktok-business-api-sdk-official==1.1.3` SHA `663b…33b7` aporta transporte genérico y declara internamente 1.2.1; entrega al menos una vez | wrapper `ADAPTED`, 16 tests, tres artefactos hash-linked e import Go estricto al owner durable pasan offline y en PostgreSQL 18.6; cuenta/permisos/suscripción/test lead/reconciliación live y autenticidad inbound siguen condicionados, sin inventar firma |
| Google Merchant | sincronización de productos/catálogo | no clasificar como captación de leads |
| Meta/Google/TikTok reporting actuales | lectura de métricas | no presentarlos como lead ingestion ni conversion feedback |

### 4. Packs recientes fuera de la composición

De los 49 packs con verificación del 2026-09-02, el perfil serverless selecciona 12 y deja 37 fuera. No todos son obligatorios en cada negocio, pero deben ser clasificados, no olvidados:

```text
GO-AGENT-DOMAIN-BINDING, GO-APP-WIRING, GO-APPLICANT-ONBOARDING-CORE,
GO-CONVERSATION-TOOLS, GO-DASHBOARDS-CORE, GO-DATA-RESIDENCY-CORE,
GO-FX-CORE, GO-GDPR-CONSENT-ERASURE-CORE, GO-GIFT-CARDS-CORE,
GO-HELP-CENTER-CORE, GO-HUMAN-APPROVAL-CORE, GO-I18N-CORE,
GO-LOYALTY-CORE, GO-MARKETING-CORE, GO-OBSERVABILITY-CORE,
GO-ONBOARDING-CORE, GO-OTP-VERIFICATION-CORE, GO-PAYROLL-CORE,
GO-PCI-DSS-SCOPE-CORE, GO-POS-CORE, GO-PROMOTIONS-CORE, GO-QR-CORE,
GO-RATE-LIMIT-CORE, GO-REFERRALS-CORE, GO-REMINDERS-CORE,
GO-RESILIENCE-CORE, GO-REVIEWS-CORE, GO-SEO-CORE, GO-SLO-CORE,
GO-SOCIAL-POSTING-CORE, GO-SURVEYS-CORE, GO-WAITLIST-CORE,
GO-WARRANTY-CLAIMS-CORE, SOC2-TSC-CONTROL-MATRIX,
SUSTAINABILITY-PILLAR-GUIDANCE, TS-DESIGN-SYSTEM,
TS-MULTIROLE-ONBOARDING
```

La solución correcta es un selector de capabilities que genere un manifest sin duplicados desde la clasificación de 48 superficies. Un “superset” ciego aumenta superficie de ataque, costo y tiempo.

## Transferencia válida desde NEW BINANCE

El proyecto local autorizado `NEW BINANCE/Binance` fue inspeccionado nuevamente en modo read-only. Aporta evidencia de un sistema propio de ingesta en tiempo real: `src/binance_lob/raw_log.py` separa append de ACK durable y usa flush/fsync; `src/binance_lob/segment_chain.py` verifica hash-chain, prefijo durable, progreso monotónico, gaps y seals; `tools/market_replay_oracle/complete_replay.py` aporta replay independiente ligado a una ejecución completa. No se atribuye a Binance ni se copia automáticamente. Sus invariantes reutilizables son:

- recibir no equivale a persistir;
- confirmar sólo después del commit durable;
- raw append-only y lineage/hash verificable;
- replay idempotente y recuperación no destructiva;
- watchdog por progreso semántico, no sólo proceso vivo;
- gaps/degradación explícitos;
- backoff acotado con jitter;
- lane independiente de reconciliación/backfill.

No se trasladan semánticas financieras de sequence/order-book, formatos BNACK ni arbitraje de WebSockets. Antes de copiar código se requiere declarar propiedad/licencia/procedencia; mientras tanto, se adoptan sólo invariantes verificadas como requisitos. V239 demuestra que la adaptación correcta para leads es commit durable antes del ACK, hash/replay/divergencia y outbox; no copia el runtime financiero.

## Autoridades que gobiernan el diseño

- CloudEvents 1.0.2 fija envelope interoperable e identidad duplicada por `source + id`: <https://github.com/cloudevents/spec/blob/ce@v1.0.2/cloudevents/spec.md>
- Google Ads Lead Form Webhook fija schema, validación, dedupe, compatibilidad y respuestas: <https://developers.google.com/google-ads/webhook/docs/implementation>
- Google publica el sample Apache-2.0 para crear el asset y configurar webhook delivery: <https://developers.google.com/google-ads/api/samples/add-lead-form-asset>
- Meta mantiene el Business SDK y su edge oficial de leads: <https://github.com/facebook/facebook-python-business-sdk>
- TikTok confirma Custom API/Webhooks en tiempo real: <https://ads.tiktok.com/resources/help/article/available-crm-integrations-tiktok-lead-generation?lang=en-GB>
- TikTok confirma postback CRM y etapas deep-funnel: <https://ads.tiktok.com/help/article/about-signal-optimization-for-lead-generation-campaigns?lang=en>
- AWS Lambda/SQS exige idempotencia ante entrega al menos una vez y partial batch: <https://docs.aws.amazon.com/lambda/latest/dg/with-sqs.html>
- Google Cloud Pub/Sub documenta redelivery, backoff y dead-letter: <https://docs.cloud.google.com/pubsub/docs/subscription-retry-policy> y <https://docs.cloud.google.com/pubsub/docs/dead-letter-topics>
- Google SRE usa SLO de freshness/correctness de pipelines: <https://sre.google/workbook/data-processing/>
- OpenAI exige herramientas declaradas mediante el campo API y Responses expone `function_call`/`function_call_output`: <https://developers.openai.com/api/reference/cli/resources/responses/methods/create> y <https://developers.openai.com/api/docs/guides/latest-model?model=gpt-4.1>

Estas fuentes gobiernan claims estrechos. El glue que las conecte será `ADAPTED` o `AUTHORED`, con tests propios; nunca se presentará como código textual de esas empresas.

## Qué puede reutilizar otro proyecto

| Capa | Reutilización |
|---|---|
| protocolos, seguridad, inbox/outbox, observabilidad, CI, pruebas | alta, si versiones y condiciones siguen válidas |
| dominio configurable, web/portales, agente y workflows | alta con blueprint, mappings y journeys propios |
| adapters de proveedor | condicionada a cuenta, permisos, región, términos, versión, sandbox y reconciliación |
| fiscal, privacidad, pagos, catálogo y documentos | condicionada a país, empresa, corpus y políticas aprobadas |
| producción | nunca heredada; se demuestra en el target concreto |

## Definition of done para reabrir el inicio de franquicia

- [x] pack omnicanal durable materializable y verificado (V224; condicionado a cuentas/target);
- [x] Google lead webhook conectado con contrato oficial, PostgreSQL y replay (V224; live condicionado);
- [x] Meta fetch/reconcile materializable con SDK/revisión exactos (V226; cuenta/live condicionados);
- [x] Meta artifact import→owner PostgreSQL conectado y probado (V227);
- [x] Meta webhook/señal inmediata reconstruible, firmado y durable con contratos oficiales exactos (V239); su ejecución live conserva cuenta/permisos/Page/form/suscripción/test lead;
- [x] TikTok demuestra contratos retrieval/subscription v1.3, artefacto oficial exacto, wrapper atribuido como `ADAPTED`, evidencia hash-linked e import durable; webhook inbound y live permanecen explícitamente condicionados;
- [x] adapter OpenAI Responses con tools estrictas y validación server-side materializable (V228; cuenta/evals condicionados);
- [x] `GO-APP-WIRING` conecta conducta real, no sólo objetos (V236);
- [x] binding identidad de canal→lead y resolver PostgreSQL real con evidencia/revocación/concurrencia (V237);
- [x] fence outbound durable; replay no reenvía y ambigüedad exige reconciliación (V238);
- [x] plan de referencia incluye los 67 packs admitidos seleccionados exactamente una vez y fija versiones vigentes; OpenGrep bloqueado queda fuera;
- [ ] browser→lead→respuesta→cita/cotización/pedido y fallos pasan E2E;
- [x] Go, PostgreSQL, frontend, seguridad, carga, restore y rollback de la composición de biblioteca pasan desde carpeta vacía; el target productivo debe repetirlos con su infraestructura real;
- [ ] documentación, ayuda, capacitación y soporte apuntan a la misma release;
- [ ] cuentas/credenciales/políticas/jurisdicción/target reales pasan sus gates.

La biblioteca ya permite iniciar discovery y construir el vertical conectado sin investigar ni escribir este plumbing desde cero. No permite prometer una franquicia productiva en una semana: los puntos live pendientes y los últimos gates son evidencia del proyecto real.
