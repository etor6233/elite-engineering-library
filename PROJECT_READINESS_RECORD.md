# Project Readiness Record — mantenimiento local

> Copiar como `PROJECT_READINESS_RECORD.md` dentro del proyecto. Este registro es un control `AUTHORED` de Elite; no es código atribuido a una empresa externa. Se completa conversando por las rondas de `PROJECT_START_READINESS_GATE.md`, no marcando casillas por inferencia.

Materializar primero `PROJECT_READINESS_GATE_PACK_PLAN.md`, copiar `project-readiness.template.json` como `PROJECT_READINESS_GATE.json` y mantener ambos registros consistentes. Antes de cada ronda generar `PROJECT_ADVISORY_<A-H>.md` con el renderer, explicarlo y enlazarlo como `advisory_prompt_ref`. Este Markdown conserva conversación/decisiones; el JSON es la entrada determinista. Sólo exit 0 y un reporte nuevo `READY_TO_BUILD` abren implementación.

## 1. Estado global

```yaml
readiness_version: "1.0"
project_id: "elite-library-maintenance"
project_name: "Elite Engineering Library — maintenance"
owner: "usuario solicitante"
operational_owners: [agente_ejecutor_dentro_de_la_autoridad_de_tarea]
status: DISCOVERY
platform_mode: BLOCK
risk_tier: HIGH
data_classification: INTERNAL
target_jurisdictions: []
languages: [es]
currencies: []
updated_at: "2026-09-07T14:51:23Z"
critical_unknowns: [MAINTENANCE_READINESS_AND_ASSURANCE_NOT_PROVEN, CAPABILITY_EVIDENCE_RECONCILIATION_PENDING]
required_access_not_proven: []
accepted_constraints: [sin_servicios_pagos_adicionales_automaticos, sin_Git_obligatorio, sin_backend_TypeScript_impuesto, procedencia_y_evidencia_explicitas]
explicitly_out_of_scope: [efectos_productivos_de_franquicia_durante_intake_de_biblioteca]
failure_ledger_path: "PROJECT_FAILURE_LESSONS.md"
open_high_or_critical_failures: [FAIL-20260907-384, FAIL-20260907-385]
dependency_update_record_path: "PROJECT_DEPENDENCY_UPDATE_RECORD.md"
open_dependency_blockers: [FAIL-20260908-532, FAIL-20260911-804]
authority_freshness_record_path: "PROJECT_AUTHORITY_FRESHNESS_RECORD.md"
stale_or_conflicted_authorities: []
vulnerability_monitoring_record_path: "PROJECT_VULNERABILITY_MONITORING_RECORD.md"
vulnerability_monitoring_status: ZERO_COST_BASELINE_BLOCKED
production_runtime_monitoring: NOT_READY_FOR_PRODUCTION_RUNTIME_MONITORING
next_round: G
```

Valores permitidos:

- `status`: `DISCOVERY|AWAITING_USER|ACCESS_VALIDATION|READY_TO_PLAN|READY_TO_BUILD|BLOCKED`;
- `platform_mode`: `OFFICIAL_PLATFORM|BUSINESS_CENTRAL_PLATFORM|CUSTOM_PLATFORM|BLOCK`;
- ninguna lista crítica se vacía sin enlazar respuesta o evidencia;
- `READY_TO_BUILD` se calcula aplicando el gate, nunca se escribe por deseo.

## 2. Log de conversación por rondas

| Ronda | Estado | Prompt exacto | Preguntas presentadas | Respuestas confirmadas | Contradicciones | Evidencia | Siguiente pregunta material |
|---|---|---|---|---|---|---|---|
| A — propiedad/resultado/límites | `ANSWERED` | `PROJECT_ADVISORY_A.md` | Resultado, responsables y límites del mantenimiento, explicados en V282/V283 | Pedido previo del usuario recuperado en sección Respuestas A; no nueva aprobación productiva | La urgencia no prueba un plazo; mantenimiento no configura negocio futuro | PROJECT_BLUEPRINT.md y specs/library-maintenance/spec.md | Continuar B; consultar sólo nuevas decisiones materiales |
| B — actores/journeys | `ANSWERED` | `PROJECT_ADVISORY_B.md` | Actores, permisos y recorridos explicados para mantenimiento CLI | Mandato previo del usuario; matriz B abajo; no autorizaciones adicionales | Pruebas parciales V284/V285, falta cierre release-bound; ANSWERED no es PROVEN | specs/library-maintenance/spec.md; reconstruction_evidence/MAINTENANCE_ENTRY_AND_SCOPE_V285.md | Ronda C; conservar los pendientes de verificación |
| C — modelo empresarial | `ANSWERED` | `PROJECT_ADVISORY_C.md` | Scope de mantenimiento explicado en V287 | Respuestas C registradas abajo y scope existente | No aprobación fiscal ni cierre de T2802/T2805 | PROJECT_BLUEPRINT.md; specs/library-maintenance/spec.md | D permanece pendiente |
| D — documentos/extracción | `ANSWERED` | `PROJECT_ADVISORY_D.md` | Alcance documental del mantenimiento explicado V391 | Alcance ya registrado: documents.required=false, automatic_storage_requested=false; no corpus empresarial para preparar la biblioteca | V366 ya separó el entrenamiento del consumer; no se solicitan chats privados | PROJECT_BLUEPRINT.md; specs/library-maintenance/spec.md; PROJECT_READINESS_RECORD.md#reconciliacion-v391 | E; reabrir D si el mantenimiento ingiere documentos empresariales |
| E — integraciones | `ANSWERED` | `PROJECT_ADVISORY_E.md` | Alcance de mantenimiento explicado V391 | Alcance previo: CLI local, sin cobros/envíos/campañas ni cuentas de provider; adquisición oficial separada y gobernada | ANSWERED no sustituye pruebas de integración/seguridad/release | PROJECT_READINESS_RECORD.md#reconciliacion-v391 | F; reabrir ante cambios materiales de alcance |
| F — UX/canales | `ANSWERED` | `PROJECT_ADVISORY_F.md` | Alcance de mantenimiento explicado V391 | Alcance previo: agente/desarrollador y operador de mantenimiento por CLI; guía NEW/EXISTING, diagnósticos y recuperación versionados | ANSWERED no sustituye pruebas de integración/seguridad/release | PROJECT_READINESS_RECORD.md#reconciliacion-v391 | G; reabrir ante cambios materiales de alcance |
| G — datos/seguridad/cumplimiento | `BLOCKED` | `PROJECT_ADVISORY_G.md` | Seguridad del mantenimiento explicada V391 | Usuario pidió diferir la parte restringida y continuar la independiente | No permite declarar seguridad integral ni release | PROJECT_DEPENDENCY_UPDATE_RECORD.md; PROJECT_FAILURE_LESSONS.md | Reabrir investigación diferida sólo con cambio real de acceso/instrucción; conservar TEST03/07 bloqueados |
| H — plataforma/operación/evidencia | `BLOCKED` | `PROJECT_ADVISORY_H.md` | Entorno local y separación de release explicados V401 | Mandato de mantenimiento local existente; sin nueva aceptación productiva | No se infiere SLO, plataforma productiva ni aceptación final | PROJECT_READINESS_RECORD.md; reconstruction_evidence/PNPM_ZIP_CONTAINMENT_V401.md | Cerrar TEST02/03 antes de producir y verificar release TEST07 |

Estados de ronda: `PENDING|IN_PROGRESS|ANSWERED|PROVEN|BLOCKED`. El agente presenta una ronda manejable desde el prompt exacto, explica términos, registra la respuesta textual confirmada y continúa; no descarga un cuestionario gigante, inventa respuestas ni salta rondas.

## 3. Decisiones y unknowns

### Respuestas B — mandato existente, no aprobación productiva

| Actor | Puede hacer en el alcance autorizado | No puede autoaprobar |
|---|---|---|
| Usuario solicitante | definir objetivo y límites, entregar decisiones/datos por mecanismos seguros | este registro no declara que haya entregado cuentas, corpus o claves |
| Agente de mantenimiento | inspeccionar fuentes públicas y archivos en scope, materializar en staging, probar, corregir su fuente canónica, registrar fallos y checkpoint | gasto, campañas, cobros, fiscalidad, datos privados, decisiones irreversibles, cambio de licencia pública, claves/tareas persistentes sin autoridad específica |
| Agente consumidor/desarrollador | usar una copia admitida en NEW o EXISTING con baseline/delta, permisos y gates del proyecto real | heredar aprobaciones de esta biblioteca o borrar cambios ajenos |

Los recorridos pedidos son NEW, EXISTING, fallo→corrección→regresión y distribución
reutilizable. V284 verifica preservación/rollback del bridge; V285 compone el gate,
genera su primera explicación, mantiene BLOCKED sin respuestas y rechaza reentrada
sin alterar 13 archivos. No demuestra todavía un sistema de franquicia completo.
Ausencias laborales, no-show, horarios de empleados y entrega física no pertenecen
al mantenimiento CLI; permanecen en los requisitos de negocio por configurar.
Cancelar o fallar una escritura local sí aplica y tiene pruebas de recuperación V284.
No se pidió otra vez esta información: procede del mandato explícito conservado en spec.

V285 clasifica las 48 superficies del mantenimiento (21 REQUIRED, 27 NONE_WITH_REASON).
Las REQUIRED conservan condiciones abiertas. Los cuatro expedientes de dependencias,
freshness, monitoreo y fuentes existen con observaciones/pendientes reales, no como
aprobaciones. Las exclusiones locales no modifican el perfil integral de franquicia.

### Respuestas A recuperadas de esta conversación

- Resultado: biblioteca para que el agente construya sistemas completos y también trabaje dentro de otros proyectos; la franquicia es el primer uso previsto. LIB-R01–LIB-R09 convierten ese pedido en criterios verificables, todavía no aprobados como ejecutados.
- Autoridad: el usuario solicita y decide el alcance; el agente ejecuta cambios locales y verificaciones autorizadas. Gasto externo, pagos, fiscalidad, datos sensibles y acciones irreversibles no se infieren de «continúa».
- Presupuesto: el usuario pidió «sin gasto adicional» y evitar consumo de Actions limitado. No se crean servicios pagos ni automatizaciones hospedadas. Esto no afirma que la ejecución del agente no consuma su plan.
- Stack y procedencia: el usuario dijo «No quiero typescript como base» y exigió fuentes/código de empresas líderes y calidad comprobable. No se elige otro stack por capricho ni se atribuye glue propio al proveedor.
- Plazo: se conserva la urgencia reiterada y el objetivo de acelerar una franquicia; no existe evidencia para convertir «una semana» en garantía técnica.
- País/moneda/fiscalidad: este intake sólo prepara mantenimiento local de biblioteca, sin ventas ni documentos fiscales reales. No define los de ningún proyecto consumidor. Idioma de asesoramiento: español, observado en la conversación.
- Scope y exclusiones: completar lo pendiente del roadmap, sin rediseño indiscriminado, duplicación de módulos ni pérdida de historia. La frase del usuario «Cada fallo debe ser registrado» se mantiene como requisito de aceptación.

Fuente: mensajes previos del usuario en esta tarea, no una confirmación nueva inventada.
ANSWERED significa pedido registrado; no PROVEN ni READY_TO_BUILD. Las clasificaciones
HIGH/INTERNAL describen conservadoramente el mantenimiento con archivos locales;
no autorizan procesar datos restringidos ni sustituir el análisis de riesgo del target.

| ID | Pregunta/decisión | Estado | Respuesta o razón | Owner | Fuente | Evidencia | Impacto | Próxima acción |
|---|---|---|---|---|---|---|---|---|
| `DEC-001` |  | `UNKNOWN` |  |  |  |  |  |  |

Estados: `ANSWERED|PROVEN|NONE_WITH_REASON|UNKNOWN|BLOCKED`. Sólo `PROVEN` requiere y acepta evidencia verificable.

## 4. Inventario de documentos

Completar una fila por `{clase, variante, emisor, idioma, versión}`. Si no existen documentos, conservar una fila `NONE_WITH_REASON` confirmada.

| Class ID | Variant ID | Emisor/idioma/versión | Ejemplos reales | Campos/input schema ID+SHA | Ground-truth owner | Riesgo por error | Lane oficial | Cuenta/probe | Evaluación | Decisión | Storage automático | Persistence mode | Mapping exacto del consumer | Output schema ID+SHA | Approval/evidencia |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
|  |  |  |  |  |  |  |  |  |  | `CORPUS_BLOCKED` | `false` |  |  |  |  |

Estados de decisión documental: `CORPUS_BLOCKED|ACCESS_BLOCKED|POLICY_BLOCKED|EVALUATION_REQUIRED|REVIEW_ONLY|LIMITED_AUTOMATION|READY_FOR_AUTOMATIC_STORAGE|NONE_WITH_REASON`.

El detalle por campo vive en `PROJECT_DOCUMENT_INTELLIGENCE_DECISION.md`, creado desde `PROJECT_DOCUMENT_INTELLIGENCE_DECISION_TEMPLATE.md`. `Class ID` y `Variant ID` son identificadores canónicos minúsculos con guiones y el par no se repite. Si `Storage automático=true`, la decisión debe ser `READY_FOR_AUTOMATIC_STORAGE`; `Persistence mode` debe ser `TRANSACTIONAL_DIRECT` u `OUTBOX_CDC_CONSUMER`; y el mapping de `PROJECT_READINESS_GATE.json` usa exactamente `mappingId`, `mappingDomainType`, `mappingVersion`, `mappingExpression`, `mappingExpressionSha256`, `mappingExactOutputKeysCsv`, `mappingInputSchemaId`, `mappingInputSchemaSha256`, `mappingOutputSchemaId` y `mappingOutputSchemaSha256`. Los bytes de la expresión, hashes, schemas, output keys, corpus, evaluación, capabilities y aprobación deben coincidir; una aprobación para factura nunca habilita packing list.

## 5. Integraciones y accesos

Una fila por proveedor/entorno. Nunca incluir secretos, tokens, certificados, claves privadas ni PII.

| Capability | Proveedor/producto | Requerimiento | Entorno | Identity reference | Scopes | Probe read-only/sandbox | Resultado | Evidencia | Reconciliation/rollback |
|---|---|---|---|---|---|---|---|---|---|
|  |  | `REQUIRED` |  |  |  |  | `NOT_PROVIDED` |  |  |

Requerimiento: `REQUIRED|OPTIONAL|NONE_WITH_REASON`. Resultado: `PROVEN|FAILED|NOT_PROVIDED|NOT_APPLICABLE`.

Inventario mínimo a resolver explícitamente: identity provider, PostgreSQL/ERP, object storage, OCR/document provider, Mercado Libre, Amazon, Mercado Pago, Stripe, Google Ads, Merchant Center, Meta Ads, TikTok Ads, WhatsApp, email/SMS/push, carriers/tracking, bancos/fintech, contabilidad/facturación local, IA/Grok, cloud/hosting, DNS/TLS y observabilidad.

## 6. Elección de plataforma empresarial

```yaml
decision: BLOCK
decision_owner: ""
decided_at: ""
reason: ""
official_platform:
  selected_source_profile: ""
  selected_product_source_ids: []
  exact_source_lock: "PROJECT_EXTERNAL_SOURCE_LOCK.md"
  product_provenance_allowed: [VERBATIM, DEPENDENCY_PIN]
  authored_or_adapted_product_code_allowed: false
  uncovered_capabilities: []
business_central:
  plan_path: ""
  accepted_development_snapshot: false
  runtime_license_probe: NOT_PROVIDED
  target_country_view: ""
  localization_gaps: []
custom_platform:
  accepted_authored_domain_code: false
  accepted_stack: []
  required_upstream_sources: []
evidence: []
```

Si `decision` es `OFFICIAL_PLATFORM`, toda capability de producto debe enlazar código oficial exacto `VERBATIM|DEPENDENCY_PIN`; cualquier necesidad de adaptación queda `NO_SOURCE/BLOCKED` hasta una nueva decisión del propietario. Si es `BUSINESS_CENTRAL_PLATFORM`, aplicar `MICROSOFT_BUSINESS_CENTRAL_CAPABILITY_PROFILE.md`. Si es `CUSTOM_PLATFORM`, todo dominio/adaptation no publicado por una fuente oficial queda `AUTHORED`. Si el propietario no acepta ninguna vía, conservar `BLOCK`.

## 7. Las 48 superficies obligatorias

Cada fila debe enlazar un expediente con los campos completos de `TOTAL_SYSTEM_CAPABILITY_CONTRACT.md`: trigger, authority, pack/trabajo, configuration, source of truth, trust/failure boundaries, licencias, tests, deployment/rollback/recovery y open conditions.

| Capability ID | Clasificación | Owner | Expediente/evidencia | Condiciones abiertas |
|---|---|---|---|---|
| `PRD-INTAKE` | `BLOCKED` |  |  |  |
| `ARCH-DOMAIN` | `BLOCKED` |  |  |  |
| `REPO-SCM` | `BLOCKED` |  |  |  |
| `CONTRACTS` | `BLOCKED` |  |  |  |
| `RUNTIME` | `BLOCKED` |  |  |  |
| `WEB-PUBLIC` | `BLOCKED` |  |  |  |
| `WEB-PORTALS` | `BLOCKED` |  |  |  |
| `MOBILE` | `BLOCKED` |  |  |  |
| `DESKTOP` | `BLOCKED` |  |  |  |
| `EMBEDDED-IOT` | `BLOCKED` |  |  |  |
| `API-BACKEND` | `BLOCKED` |  |  |  |
| `DOMAIN-MODULES` | `BLOCKED` |  |  |  |
| `IDENTITY` | `BLOCKED` |  |  |  |
| `AUTHORIZATION` | `BLOCKED` |  |  |  |
| `SECRETS-PKI` | `BLOCKED` |  |  |  |
| `TX-DATABASE` | `BLOCKED` |  |  |  |
| `CACHE` | `BLOCKED` |  |  |  |
| `SEARCH` | `BLOCKED` |  |  |  |
| `OBJECT-STORAGE` | `BLOCKED` |  |  |  |
| `OUTBOX-INBOX` | `BLOCKED` |  |  |  |
| `JOBS-WORKFLOWS` | `BLOCKED` |  |  |  |
| `BROKER-STREAMING` | `BLOCKED` |  |  |  |
| `INTEGRATIONS` | `BLOCKED` |  |  |  |
| `PAYMENTS` | `BLOCKED` |  |  |  |
| `MARKETPLACES` | `BLOCKED` |  |  |  |
| `ADS-ATTRIBUTION` | `BLOCKED` |  |  |  |
| `NOTIFICATIONS` | `BLOCKED` |  |  |  |
| `DATA-INGEST` | `BLOCKED` |  |  |  |
| `ANALYTICS-BI` | `BLOCKED` |  |  |  |
| `ML-AI` | `BLOCKED` |  |  |  |
| `RAG-AGENTS` | `BLOCKED` |  |  |  |
| `GPU-ACCEL` | `BLOCKED` |  |  |  |
| `NETWORK-EDGE` | `BLOCKED` |  |  |  |
| `CONTAINERS` | `BLOCKED` |  |  |  |
| `ORCHESTRATION` | `BLOCKED` |  |  |  |
| `IAC-CLOUD` | `BLOCKED` |  |  |  |
| `CI` | `BLOCKED` |  |  |  |
| `CD-RELEASE` | `BLOCKED` |  |  |  |
| `SUPPLY-CHAIN` | `BLOCKED` |  |  |  |
| `OBSERVABILITY` | `BLOCKED` |  |  |  |
| `SLO-INCIDENT` | `BLOCKED` |  |  |  |
| `PERFORMANCE` | `BLOCKED` |  |  |  |
| `SECURITY-APPSEC` | `BLOCKED` |  |  |  |
| `PRIVACY-COMPLIANCE` | `BLOCKED` |  |  |  |
| `BACKUP-DR` | `BLOCKED` |  |  |  |
| `COST-FINOPS` | `BLOCKED` |  |  |  |
| `TEST-PLATFORM` | `BLOCKED` |  |  |  |
| `DOCS-OPS` | `BLOCKED` |  |  |  |

Clasificación permitida: `REQUIRED|OPTIONAL|NONE_WITH_REASON|BLOCKED`. El estado inicial `BLOCKED` obliga a decidir, no implica que todas las capacidades deban implementarse.

## 8. Sources, packs y procedencia

| Capability | Fuente oficial exacta | Revision/hash | Licencia/notices | Modo | Claim permitido | Non-claims | Tests requeridos | Estado |
|---|---|---|---|---|---|---|---|---|
|  |  |  |  | `NO_SOURCE` |  |  |  | `BLOCKED` |

Modo: `VERBATIM|DEPENDENCY_PIN|ADAPTED|AUTHORED|NO_SOURCE`. Una fila `NO_SOURCE` nunca se renombra como código de empresa líder. `OFFICIAL_PLATFORM` permite para producto únicamente `VERBATIM|DEPENDENCY_PIN`; `AUTHORED` sólo se admite con `CUSTOM_PLATFORM` o como extensión declarada de Business Central.

## 9. Evidencia de preparación

| Gate | Comando/probe | Expected | Observed | Evidence path | Estado |
|---|---|---|---|---|---|
| biblioteca | `VERIFY_EXECUTABLE_LIBRARY.ps1 -Mode Preflight` | PASS + toolchains |  |  | `PENDING` |
| upstream lock | adquisición seleccionada | hash/licencia/receipt PASS |  |  | `PENDING` |
| failure learning | ledger + huellas/lecciones aplicables | cero HIGH/CRITICAL abiertos para P1 |  |  | `PENDING` |
| dependency lifecycle | inventory/locks/EOL/advisories/update policy | cero dependency blocker aplicable a P1 |  |  | `PENDING` |
| authority freshness | sources/version/date/applicability/corrections | cero stale/conflict blocker aplicable a P1 |  |  | `PENDING` |
| plataforma | runtime/tenant/container/database | identidad y ambiente correctos |  |  | `PENDING` |
| documentos | corpus/ground truth/provider benchmark | decisión por clase/campo |  |  | `PENDING` |
| seguridad | threat/access/privacy | owners y controles cerrados |  |  | `PENDING` |
| operación | SLO/backup/restore/rollback | evidencia en target |  |  | `PENDING` |
| vertical slice | journey P1 end-to-end | acceptance tests PASS |  |  | `PENDING` |

## 10. Cálculo de readiness

Mantener `status: BLOCKED` o `AWAITING_USER` si ocurre cualquiera:

- ronda A–H no resuelta;
- capability sin clasificación o sin owner;
- `REQUIRED` sin source/pack/trabajo autorizado;
- `critical_unknowns` o `required_access_not_proven` no vacíos;
- `PROJECT_FAILURE_LESSONS.md` ausente, desactualizado o con fallos `HIGH`/`CRITICAL` abiertos que afecten el primer vertical slice;
- `PROJECT_DEPENDENCY_UPDATE_RECORD.md` ausente o con dependencia/runtime/upstream requerido sin owner, lock, policy, soporte o rollback;
- `PROJECT_AUTHORITY_FRESHNESS_RECORD.md` ausente o con authority aplicable `STALE_BLOCKED`/`CONFLICT_UNRESOLVED`;
- `PROJECT_VULNERABILITY_MONITORING_RECORD.md` ausente, OSV/Dependabot sin prueba vigente o finding/coverage blocker abierto;
- plataforma en `BLOCK` o sin probe;
- documentos requeridos sin expediente, inventario clase/variante y corpus; o almacenamiento automático sin autoridad exacta de mapping/schema/evaluación/aprobación por clase;
- integración requerida sin cuenta/sandbox/probe;
- reglas, jurisdicción, seguridad, SLO o rollback críticos desconocidos;
- licencia/procedencia no resuelta;
- contradicción entre blueprint, spec, plan, source lock o este registro.

`READY_TO_PLAN` autoriza arquitectura y tasks; no código de producto. `READY_TO_BUILD` exige además que todos los requisitos de la sección 9 de `PROJECT_START_READINESS_GATE.md` estén demostrados. El agente registra qué condición produjo el cambio de estado y conserva el historial.

## Baseline observado V282

### Ronda D — requisito aclarado, V289

PROJECT_ADVISORY_D.md fue generado desde catálogo 1.0.0 y explicado: originales,
campos, corpus y validación son cosas distintas. El usuario pide preservar crudos
y separar interpretación; queda incorporado a LIB-R10 y TEST-09 del contrato.
No indicó canal histórico, formatos, corpus autorizado, mapping, umbrales ni owner
de aprobación. D sigue PENDING; no inventar esas respuestas ni pedir secretos.
Esto no impide validar localmente el contrato de mantenimiento y sus referencias.
La creación de PROJECT_ENGINEERING_CONTRACT.json no cambia readiness a READY.

### Ronda C — mantenimiento, V287

Prompt exacto generado: PROJECT_ADVISORY_C.md. Explicado en conversación: stock,
dinero, compras/ventas y fiscalidad no se ejecutan en este mantenimiento; no se
asumen reglas ni autoridad contable de una franquicia todavía no definida.
La evidencia es el scope existente en specs/library-maintenance/spec.md y
PROJECT_BLUEPRINT.md y la petición vigente de completar la biblioteca, no una
nueva aprobación fiscal del usuario. Clasificación: ANSWERED para este scope,
NO_APLICA a efectos comerciales/fiscales en el primer slice.
Entidades reales de mantenimiento: pack, archivo, revisión, evidencia, checkpoint,
fallo y release; owners/estados son los contratos existentes, no nuevos agregados
financieros. Invariantes: selección/licencia/hash, no sobrescribir terceros,
no reescribir eventos, no promover por un test sintético.
Trigger de reapertura: cualquier futuro proyecto que cotice, reserve, compre,
cobre, devuelva, contabilice, emita comprobantes o configure franquicias deberá
completar C con sus responsables/datos/reglas reales. T2802/T2805 no se cierran.
V287 prueba NEW/EXISTING del tooling; no compra/importación/venta productivas.

Este expediente corresponde al mantenimiento de la biblioteca, no a una franquicia en producción. Se conserva DISCOVERY; los valores restrictivos no resueltos vienen del template y no son decisiones aprobadas. Ninguna ronda se marca contestada por crear el archivo. El roadmap único sigue siendo REUSABLE_CODE_READINESS_ROADMAP.md §10; el estado y sus eventos permanecen en la raíz. V278–V281 preservan pruebas y límites sin reemplazar requisitos ni implementation_assurance.

Datos ya expresados por el usuario que deben reconciliarse con cada ronda: reutilización en proyectos nuevos y existentes, código/fuentes oficiales trazables, adaptación explícita y probada sin atribución falsa, no imponer backend TypeScript, evitar gasto incremental y Git obligatorio, no perder historia ni heredar aprobaciones en otros proyectos. No significan cuenta cloud, documentos reales, reglas fiscales ni producción autorizados.

Próximo paso: completar el delta de preparación de mantenimiento con estas respuestas y los artefactos reales. No exigir cuentas de una franquicia hipotética para verificar herramientas locales; no excluir capacidades requeridas del perfil integral usando esa distinción.

### Conciliación V324 (sin nuevas respuestas)

La tabla C se alinea con JSON y respuesta V287 ya existentes. FAIL386 se retira de fallos abiertos por su cierre focal V292; la observabilidad integral sigue CANDIDATE/pendiente. FAIL532 se incorpora al resumen de dependencia por evidencia nativa/licencias aún insuficiente. D–H, FAIL384/385, plataforma y demás condiciones permanecen abiertas. Se preservaron los tres archivos before y se ejecutó de nuevo el validator0.6.1; el reporte actual continúa BLOCKED.

### Clarificación de historial/entrenamiento V366 — respuesta del usuario

El usuario distingue NEW sin historial y EXISTING con historial, y aclara que el propósito era entrenar un modelo con código de élite. No entregó corpus, canal, modelo, proveedor, método ni presupuesto. Esos inputs corresponden al proyecto consumidor que active el entrenamiento; no se le piden para preparar esta biblioteca. Se corrige la solicitud de V365 y la interpretación V288 mediante PROJECT_HISTORY_MODEL_TRAINING_CONTRACT. documents.required=false y automatic_storage_requested=false ya describen este mantenimiento y permanecen intactos. No se cierra toda la ronda D por esta aclaración puntual ni se heredan respuestas en futuros proyectos. TEST09 sigue bloqueado por capacidad reusable pendiente, no por no haber solicitado datos de una empresa hipotética.

## Reconciliacion V391

D: se recupera la respuesta de alcance existente, no una respuesta nueva inferida. El blueprint y el JSON ya excluyen ingesta y almacenamiento automático de documentos empresariales de este mantenimiento. V366 aclara que preparar la biblioteca no requiere chats privados de un consumer hipotético. No se procesó un corpus ni se aprobaron precisión, retención o mappings. T2806 y el contrato de entrenamiento conservan sus gates por proyecto. Reabrir esta ronda ante cualquier nueva ingesta empresarial en el mantenimiento.

E: NEW/EXISTING materializan y verifican archivos locales. No hay provider empresarial REQUIRED en ese primer slice: integrations=[] ya registrado, sin secreto, cuenta ni webhook de negocio. Las exclusiones de efectos externos proceden de la especificación aprobada; no se cancelan los adapters de T2805. La adquisición de herramientas/fuentes conserva PROJECT_OFFICIAL_SOURCE_PROFILE_RECORD y PROJECT_DEPENDENCY_UPDATE_RECORD, sus identidades exactas, costos y bloqueos. No se declara conexión live ni seguridad del grafo. Reabrir con cualquier journey del mantenimiento que dependa de un proveedor remoto; un consumer recibe su propio intake.

F: se recupera la definición de actores/journeys de las rondas A-B y del blueprint, sin elegir una interfaz nueva. El canal del mantenimiento es CLI; no requiere branding, SEO, analytics, pantallas o navegadores de negocio. La guía PROJECT_ENTRY_CLI_GUIDE0.1.2 cubre preparación, selección, diagnóstico, ayuda y recuperación para agente/desarrollador/operador con comandos y owners exactos. TEST01/46/47 aportan pruebas locales; TEST06/08 y V352-V353 prueban un candidato histórico exacto. Las mediciones V326 y los límites de los runners son evidencia técnica, no un SLO ni plazo comercial inventado. El user exige rapidez y conservación de cambios/fallos. ANSWERED registra ese alcance; no declara probado el payload actual, capacitación humana independiente, todos los portales, ni aceptación de una release nueva. Reabrir ante cambio de actores/canal/guía y revalidar el payload final bajo TEST07.

G: continúa BLOCKED. La instrucción explícita del usuario difiere la investigación nativa restringida; no equivale a aceptación de riesgo, prueba de seguridad ni eliminación del finding. No se solicita una cuenta de IdP para el CLI ni se heredan los fixtures de identidad de la franquicia. H permanece PENDING y no se declaran platform_mode, operación, implementación integral ni release resueltos.

## Ronda H — V401: explicación y frontera pendiente

PROJECT_ADVISORY_H.md fue generado desde catálogo1.0.0. La tarea autorizada sigue
siendo mantenimiento local de biblioteca, sin despliegue de una franquicia real.
Un ZIP local y sus pruebas no habilitan producción. No se inventan regiones,
on-call, SLO, RPO/RTO ni una aceptación de negocio. El usuario ya pidió preservar
la parte Daybreak diferida y continuar el trabajo independiente. H queda BLOCKED
por evidencia de release/seguridad/integración pendiente, no por falta del prompt.
La nueva contención de pnpm tiene pruebas acotadas y no elimina esas condiciones.

V402320: G/H ANSWERED para la decisión explícita LIBRARY_INFRASTRUCTURE, según LIBRARY_ADVISORY_GH_V402.md; nueve dimensiones reconciliadas en LIBRARY_ASSURANCE_REVIEW_V402.json. critical_unknowns vacío para este scope; trabajo conocido T2810/ARCA y dossier último no equivale a incertidumbre humana. El gate release no está aprobado todavía; el expediente legacy de proyecto no habilita producción.
