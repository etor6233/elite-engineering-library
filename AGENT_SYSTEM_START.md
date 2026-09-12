# Agent Engineering System — Start Here

> **Producto:** sistema Markdown agnóstico para que un agente de código componga proyectos verificables a partir de conocimiento, contratos y packs de implementación admitidos.
> **No es:** una aplicación, un framework obligatorio, un volcado de repositorios famosos ni una promesa de que toda regla empresarial pueda convertirse en configuración.

## 1. Resultado esperado

```text
pedido del usuario
→ constitution + PROJECT_BLUEPRINT.md + spec.md
→ selección mínima de autoridades
→ arquitectura y stack justificados
→ plan.md + packs compatibles y licenciados
→ tasks.md + plan exacto de archivos
→ materialización de código/config/tests
→ build y vertical slices
→ gates de seguridad, corrección, performance y operación
→ análisis + convergencia + evidencia y readiness honesto
```

El sistema debe permitir iniciar rápido sin copiar componentes incompatibles ni fingir que un ejemplo educativo es producción.

Credenciales/configuración: usar la guía única
`markdown_system/PROJECT_SECRETS_TEMPLATE.md` al preparar accesos. Conciliarla con
la composición efectiva, explicar cuenta/config/secreto/clave interna/identidad
temporal y completar la receta del proveedor elegido antes de pedir accesos.
No credenciales en Markdown/chat ni alternativas simultáneas obligatorias;
custodia local mediante `markdown_system/project_local_secrets.ps1` sólo en Windows
DEV, target mediante gestor elegido. Diferir una integración requiere comprobar
que no se monta/ejecuta, no sólo un flag visual. No certificar adapter por token.

## 2. Activación multiagente

- Codex: el `AGENTS.md` raíz descubre este archivo mediante su orden obligatorio.
- Claude Code: `CLAUDE.md` importa `AGENTS.md`, conservando una autoridad compartida.
- Otro agente: se le indica explícitamente que lea `AGENTS.md` y este archivo antes de actuar.

Cuando la biblioteca no es la raíz del proyecto, instalar el bridge con `INSTALL_AGENT_BRIDGE.ps1`. Crea las Skills `.agents/skills/elite-engineering-library/SKILL.md` para Codex y `.claude/skills/elite-engineering-library/SKILL.md` para Claude: los agentes descubren metadata breve y cargan el workflow al activarlo. Sin raíz Git, iniciar Codex desde la raíz exacta del proyecto; la detección oficial de instrucciones sólo consulta el directorio actual cuando no encuentra una raíz de proyecto. Claude puede exigir aprobación la primera vez que el import apunta fuera del proyecto; hasta concederla, el agente no debe afirmar que recibió la autoridad Elite.

Los archivos de arranque son routers compactos. Los manuales de más de cientos de líneas se cargan sólo cuando una decisión material los necesita.

Historial para entrenar modelos: aplicar `markdown_system/PROJECT_HISTORY_MODEL_TRAINING_CONTRACT.md`. NEW no exige historial; EXISTING usa datos privados del proyecto consumidor. No confundir entrenamiento explícito con memoria/RAG/evals, ni pedir chats para completar la biblioteca. No importar, subir datos o entrenar sin código admitido y gates del proyecto.

Codex descubre instrucciones desde la raíz del proyecto hasta el directorio de trabajo y las más cercanas tienen precedencia. Por eso `AGENTS.md` permanece como router breve y los overrides específicos deben vivir cerca de su subsistema; no se vuelca el corpus entero en la cadena inicial. El límite combinado predeterminado es 32 KiB. Claude resuelve imports relativos desde el `CLAUDE.md` que los contiene y recomienda Skills de proyecto para procedimientos bajo demanda. Fuentes oficiales vigentes: <https://learn.chatgpt.com/docs/agent-configuration/agents-md>, <https://learn.chatgpt.com/docs/build-skills>, <https://code.claude.com/docs/en/memory> y <https://code.claude.com/docs/en/slash-commands>.

La autonomía operativa se rige por `markdown_system/AGENT_AUTONOMY_CONTRACT.md`. `START_ANY_PROJECT.md` contiene la instrucción portable para iniciar un proyecto desde la primera conversación.
Todo fallo durante materialización, build, test u operación sigue `markdown_system/AGENT_ERROR_RECOVERY_PROTOCOL.md` y `markdown_system/FAILURE_LEARNING_CONTRACT.md`: registrarlo en `PROJECT_FAILURE_LESSONS.md`, consultar huellas anteriores, corregir fuente canónica, añadir regresión, reconstruir limpio, repetir gates afectados, promover la lección reusable y continuar sin espera ceremonial.
Todo inicio o reanudación materializa `ENGINEERING_EXECUTION_VALIDATOR.md` 1.3.1 y valida `PROJECT_EXECUTION_STATE.json` más su cadena `PROJECT_EXECUTION_EVENTS.jsonl`. El engineering contract 1.1.0 exige `implementation_assurance`: procedencia, riesgo, nueve dimensiones, fuentes de método, locks cuando hay upstream y evidencia distinta del artefacto y release; no permite promover código propio/adaptado por compilación o reputación. El cursor enlaza por hash readiness, Spec Kit, tasks, fallos, dependencias, evidencia y release; carga sólo `context.must_read_refs`, conserva lo estable en `reuse_without_reload_refs` y no duplica esos owners. Un proyecto EXISTING exige inventario observado y delta; uno NEW exige baseline vacío o scaffold. Un paso sólo se checkpointa cuando su artifact y evidencia existen.
Toda dependencia o upstream se mantiene conforme a `markdown_system/DEPENDENCY_UPDATE_CONTRACT.md`; todo claim temporal o contradicción se gobierna por `markdown_system/AUTHORITY_FRESHNESS_AND_SELF_CORRECTION_CONTRACT.md`. El agente puede proponer y probar candidatos aislados, pero no adoptar “latest”, borrar un snapshot ni autoaprobar una corrección de negocio.
La vigilancia baseline sin gasto incremental se rige por `markdown_system/ZERO_COST_VULNERABILITY_MONITORING_PROFILE.md`: crear `PROJECT_VULNERABILITY_MONITORING_RECORD.md`, usar Dependabot incluido y Google OSV-Scanner local fijado, no gastar Actions en scans programados y no autorizar producción sin monitoring runtime probado.

Antes de escribir implementación, componer `markdown_system/PROJECT_READINESS_GATE_PACK_PLAN.md` y aplicar `markdown_system/PROJECT_START_READINESS_GATE.md`. El agente puede investigar y ejecutar probes read-only, pero sólo materializa producto cuando `PROJECT_READINESS_RECORD.md` y `PROJECT_READINESS_GATE.json` sean consistentes y el validador materializado emita exit 0/reporte `READY_TO_BUILD`. Para documentos requeridos, el JSON debe inventariar cada clase/variante; cualquier almacenamiento automático queda bloqueado sin modo de persistencia, mapping directamente compatible con el consumer, expresión y SHA coincidentes, schemas/output keys, capabilities, corpus, evaluación, aprobación y owner de esa misma clase. En modo `OUTBOX_CDC_CONSUMER`, el agente ejecuta después el renderer materializado para transferir sólo los diez campos exactos a un template consumer intacto; el perfil resultante permanece sin endpoints, secrets, owners, approvals, acknowledgement ni proofs. El gate también se aplica como auditoría delta al continuar un proyecto existente.

Antes de preguntar cada ronda pendiente, ejecutar `project_readiness_gate/render_project_advisory.py --project-root . --round <A-H> --output PROJECT_ADVISORY_<A-H>.md`, presentar su contenido en lenguaje simple y registrar la respuesta/evidencia real. Un ejemplo del prompt nunca es respuesta asumida. El validator rechaza una ronda `ANSWERED|PROVEN` si `advisory_prompt_ref` falta, fue modificado, corresponde a otra ronda o escapa la raíz.

En mantenimiento de esta biblioteca, `PROJECT_EXECUTION_STATE.json` y `PROJECT_EXECUTION_EVENTS.jsonl` viven en la raíz local y se excluyen del payload portable. `VERIFY_LIBRARY.ps1` exige el par íntegro cuando existe y ejecuta el validator reconstruido con Python; un checkpoint `DISCOVERY/BLOCKED` válido prueba continuidad, no `READY_TO_BUILD`. El kit generado permanece en un destino separado, sin crear otra fuente. La distribución sin esos dos archivos inicia su propio estado y no hereda aprobaciones ni progreso de este mantenimiento. Cualquier otro archivo raíz inesperado sigue rechazado.

## 3. Primer ciclo obligatorio

### Expedientes locales de mantenimiento

Con el par de checkpoint presente, el verificador y el empaquetador reconocen exclusivamente los veinte nombres de readiness/advisory/assurance definidos por `Test-LocalMaintenanceEntry`, más `.specify/memory/constitution.md` y `specs/library-maintenance/{spec,plan,tasks}.md`. Son expedientes locales excluidos del payload portable. No se permite un comodín PROJECT_* ni otros proyectos dentro de specs; se rechazan archivos desconocidos y reparse points antes de descender. Los kits ejecutables siguen materializándose fuera de la distribución.

Esta exclusión resuelve almacenamiento y empaquetado, no la validez semántica de los expedientes. Un registro vacío o BLOCKED nunca habilita implementación: ejecutar el readiness validator y enlazar su reporte exacto al checkpoint. Los paths spec/plan/tasks explican el delta y enlazan el roadmap vigente, no crean otro backlog. No convertir este mantenimiento en una franquicia hipotética ni asumir por él cuentas, costos, reglas fiscales o accesos de futuros proyectos.

1. Clasificar la solicitud: discovery, diseño, implementación, revisión, incidente o mantenimiento de biblioteca.
2. Ejecutar `markdown_system/PROJECT_START_READINESS_GATE.md`; generar y explicar el prompt exacto de cada ronda; crear desde template `PROJECT_READINESS_RECORD.md`, `PROJECT_FAILURE_LESSONS.md`, `PROJECT_DEPENDENCY_UPDATE_RECORD.md`, `PROJECT_AUTHORITY_FRESHNESS_RECORD.md`, `PROJECT_VULNERABILITY_MONITORING_RECORD.md` y `PROJECT_OFFICIAL_SOURCE_PROFILE_RECORD.md`; materializar el execution validator, crear el estado/event log y checkpointar el baseline observado; consultar el ledger compartido, completar el intake y mantener implementación bloqueada hasta `READY_TO_BUILD`.
3. Ejecutar `VERIFY_EXECUTABLE_LIBRARY.ps1 -Mode Preflight` desde la biblioteca; registrar faltantes sin interpretar un PASS estructural como acceso productivo.
4. Crear o actualizar `PROJECT_BLUEPRINT.md` con `markdown_system/PROJECT_BLUEPRINT_CONTRACT.md`.
5. Consultar `AI_ENGINEERING_MASTER_MAP.md`, `SYSTEMS_ENGINEERING_MASTER_MAP.md`, `markdown_system/CAPABILITY_CATALOG.md` y clasificar todas las filas de `markdown_system/TOTAL_SYSTEM_CAPABILITY_CONTRACT.md`; conectar cada capability `REQUIRED` mediante los journeys release-bound exigidos por el validator, sin aceptar pantallas, backends, capacitación o soporte aislados.
   Si una capability requerida no tiene pack compatible admitido, materializar `markdown_system/CAPABILITY_GAP_RESOLUTION_PACK_PLAN.md`, crear un expediente JSON por capability y ejecutar `CAPABILITY-GAP-RESOLUTION-GATE`: investigar primero autoridades, documentación, repositorios, releases y advisories oficiales actuales; preservar bytes/refs/licencia/revisión/hash/claim y G0–G8; continuar mientras el estado sea `RESEARCH_INCOMPLETE`; programar sólo con `USE_REUSABLE_PACK`. Un resultado condicionado o bloqueado se comunica con evidencia y trigger de reapertura y se devuelve a los owners canónicos; nunca se rellena con código o procedencia inventados.
   Para código privado Go/TypeScript sin GitHub Code Security, usar `MICROSOFT_DEVSKIM_ADAPTED_SAST_PACK_PLAN.md` como lint baseline. `GITLAB_OPENGREP_SIGNED_SAST_PACK_PLAN.md` queda temporalmente bloqueado: OpenGrep/rules siguen fijados, pero Cosign 3.1.3 reabrió admisión por SCA actual; no ejecutar ese verifier ni afirmar un segundo lane hasta fijar una alternativa limpia. Conservar JSON/SARIF/triage/fixes y no presentar lint/SAST como análisis interprocedural completo ni como prueba de seguridad.
6. Crear `PROJECT_AUTHORITY_MAP.md` y explicar inclusiones/exclusiones.
7. Crear `PROJECT_EXTERNAL_SOURCE_LOCK.md` desde `markdown_system/REVESTEX_ELITE_PUBLIC_CODE_ADMISSION_LEDGER.md` cuando aplique código externo.
8. Cuando se aprueben fuentes oficiales, completar `PROJECT_OFFICIAL_SOURCE_PROFILE_RECORD.md`, componer `markdown_system/PROJECT_INITIALIZATION_PACK_PLAN.md`, validar el JSON materializado y usar `elite_sources/apply_source_profile.ps1` para adquirirlas sin Git antes del perfil de producto.
   Para TikTok Lead v1.3, `markdown_system/TIKTOK_LEAD_ADAPTER_PACK_PLAN.md` es la lane separada de reporting: materializa retrieval autenticado, evidencia hash-linked e import durable al owner PostgreSQL/outbox compartido. El wrapper es `ADAPTED` sobre el transporte genérico del wheel oficial exacto; no atribuirlo a TikTok ni aceptar un webhook inbound como auténtico hasta demostrar el mecanismo oficial en el proyecto.
   Para documentos, materializar primero `DOCUMENT_PIPELINE_ROUTING_PACK_PLAN.md` y seleccionar perfiles sólo desde su receipt hash-linked. `DURABLE_DOCUMENT_PIPELINE_PACK_PLAN.md` es el plano de control vendor-neutral 5/51 para un destino separado; no se superpone con un perfil provider porque éstos ya incorporan su orquestador. Las opciones cerradas son `AZURE_DOCUMENT_RUNTIME_PACK_PLAN.md`, `GOOGLE_DOCUMENT_RUNTIME_PACK_PLAN.md`, `AWS_TEXTRACT_DOCUMENT_RUNTIME_PACK_PLAN.md`, `PADDLEOCR_LOCAL_RUNTIME_PACK_PLAN.md` y `MARKITDOWN_LOCAL_RUNTIME_PACK_PLAN.md`. Azure 10/85 incorpora Blob version-level WORM `Locked`, scope no reemplazable, descarga de la versión exacta y los samples Python oficiales Microsoft `prebuilt-invoice`, bytes locales `prebuilt-documentSearch` y same-resource analyzer copy; Google 9/102 incorpora el sample oficial Custom Document Extractor con `schema_override` variable por solicitud, lifecycle oficial processor/dataset/import/train/evaluate/deploy/default/undeploy, batch GCS y response handling OCR/form/table/entity/split/layout/custom y GCS create único/CRC32C/KMS/Object Retention y AWS 7/75 incorpora S3 Object Lock más Amazon Textractor oficial 1.10.0 hash-locked. `prebuilt-documentSearch` preserva contenido/layout pero no sustituye un schema de negocio. El copy sample materializa código pero nunca se autoejecuta: crea/actualiza/elimina recursos y exige autoridad, receipt e inventario/reconciliación. Las tres lanes cloud incluyen Microsoft Durable Task para security→retained original→extract→evaluate→review→idempotent persist→retained final, pero mantienen todo efecto bloqueado hasta autoridad/IAM/RBAC/costo/política/corpus y verificación target; el backend in-memory es sólo test. Si el blueprint exige un historial local tamper-evident, componer `TESSERA_POSIX_EVIDENCE_LOG_PACK_PLAN.md` en destino dedicado y mantenerlo bloqueado hasta Linux/POSIX, receipts sin datos, keys, checkpoint/witness externo, backup/restore, acceso y operación demostrados; nunca sustituir por él el storage WORM. Para carga portal/API hacia cuarentena S3 usar `AWS_SECURE_DOCUMENT_INTAKE_PACK_PLAN.md` 3/35 y exigir sesión server-side/CAS, versión exacta y posterior gate Magika/ClamAV/YARA-X; no confundirlo con recepción email/SFTP/scanner/eventos. Para recepción email, el source AWS SES Mail Manager 79314a9 permanece `SAMPLE_ONLY_REJECTED_FOR_IMMEDIATE_ADOPTION`; para SFTP, tanto AWS secure Transfer Family 474ca07 como AWS file-transfer-sync 5ee79ee permanecen rechazados: el primero por aislamiento compartido/idempotencia/delivery sin retry/DLQ y el segundo por instalación durante synth, comparación por timestamp y finalización asíncrona no reconciliada. Para cámara/mobile, AWS receipt 6eb72bf permanece rechazado: PWA/infra compilan, pero faltan tests y existen SCA 14/13, debug pre-auth, POST sin checksum/version binding y eventos fallidos absorbidos. Adquirir las cuatro sólo como referencias exactas y no desplegarlas. Para AWS S3/SES usar además `AWS_ENTERPRISE_ADAPTERS_PACK_PLAN.md` cuando corresponda; para Stripe/Mercado Pago usar `PAYMENT_WEBHOOK_ADAPTERS_PACK_PLAN.md`; para Google Ads read-only usar `GOOGLE_ADS_REPORTING_PACK_PLAN.md`; para Meta Ads Insights read-only usar `META_ADS_REPORTING_PACK_PLAN.md`; para TikTok Business reporting read-only usar `TIKTOK_ADS_REPORTING_PACK_PLAN.md`; para Meta WhatsApp template/webhook usar `META_WHATSAPP_CLOUD_PACK_PLAN.md`; para Firebase Cloud Messaging usar `FIREBASE_PUSH_PACK_PLAN.md`; para Mercado Libre item/order/validator/notification fetch usar `MERCADOLIBRE_MARKETPLACE_PACK_PLAN.md` y, si el journey responde preguntas, componer además `GO-MERCADOLIBRE-QUESTION-OUTBOUND` después de `GO-PG-OUTBOUND-DELIVERY-FENCE 0.2.x` —ya seleccionados juntos por el perfil integral—; para Amazon Catalog read-only usar `AMAZON_SPAPI_CATALOG_PACK_PLAN.md`; para Google Merchant ProductInput/status usar `GOOGLE_MERCHANT_PRODUCT_SYNC_PACK_PLAN.md`. Ninguno se activa sin la fila de integración/acceso/costo y sus conditions; TikTok requiere acquisition receipt del source oficial aprobado, WhatsApp requiere licencia/terms/consent/version/template, Firebase exige proyecto/ADC/API/dispositivo consentido/políticas, Mercado Libre exige app/seller/token/scopes/origin/policies/aprobación y Merchant requiere aprobación explícita del write externo.
   Si el blueprint incluye devoluciones de dinero dentro del backend empresarial, `ENTERPRISE_BACKEND_PACK_PLAN.md` incorpora `GO-OFFICIAL-RETURN-REFUND-WORKER` 0.1.0. Seleccionarlo sólo después de probar cuenta/sandbox, secretos por referencia, provider/currency, webhook y reconciliation. El worker deriva monto desde una única línea quantity-one ligada al stock y exige un único pago elegible; múltiples pagos, descuentos/impuestos/envío sin allocation aprobada o Mercado Pago fuera de ARS permanecen bloqueados. Ejecutar desde Markdown los contratos SDK, PostgreSQL 0001–0020/down-up e integración parcial→total; nunca presentar un refund solicitado/pending como dinero devuelto.
   Si el blueprint incluye cambio de producto, el mismo perfil incorpora `GO-RETURN-EXCHANGE-FULFILLMENT-WORKER` 0.1.0. Exigir devolución recibida, disposición exchange, efecto físico de inventario exitoso, pedido original de una línea quantity-one, misma variante/sucursal y request contable separado. Materializar y probar PostgreSQL 0001–0021; `prepared` sólo confirma pedido de reemplazo, reserva y handover locales. Diferencia de producto/precio, pedido multi-línea, cross-location, carrier, entrega, aceptación, accounting o fiscalidad no demostrados permanecen `BLOCKED`.
   Para el request contable usar `GO-RETURN-ACCOUNTING-REVERSAL-WORKER` 0.1.0: exige un único journal `SALE` original posteado y aprobado, inventory y remedy exitosos. Revierte exactamente sus líneas; nunca inventa cuentas/costos/impuestos. Probar el perfil PostgreSQL vigente 0001–0027 y recuperación post-commit. `accounting succeeded` no equivale por sí solo a nota de crédito ARCA autorizada ni cierre auditado.
   Para solicitar una nota de crédito desde una devolución usar `GO-RETURN-FISCAL-CREDIT-NOTE-WORKER` 0.1.0 junto con `GO-ARCA-FISCAL-ISSUANCE-API` 0.6.0, `MICROSOFT-ARCA-WSFE-SOAP-ADAPTER` 0.3.0 y `ARCA-WSFE-UDS-WORKER` 0.3.0. El worker exige inventory/refund/accounting exitosos y montos exactos; sólo convierte `1→3`, `6→8`, `11→13`, delega WSFE al owner existente y espera autorización. Probar generación WSFE, full Go y el perfil PostgreSQL vigente 0001–0027. No afirmar homologación, notas parciales/débito, rendering ni validez fiscal real sin evidencia del proyecto.
   La auditoría email vigente también fija `aws-samples/serverless-mail` 70ac931 y AWS Security IR email 08f11d4 como referencias rechazadas: el primero no tiene tests/lock ni delivery seguro; el segundo pasa 42 tests pero no extrae adjuntos, no valida verdicts SPF/DKIM/DMARC y deja la DLQ desconectada. `sample-bda-redaction` no publica licencia y no puede copiarse. La auditoría SFTP fija AWS secure Transfer Family 474ca07 y AWS file-transfer-sync 5ee79ee como referencias rechazadas por aislamiento/idempotencia/delivery y por no reconciliar finalización asíncrona, respectivamente. La auditoría mobile añade AWS receipt 6eb72bf como referencia rechazada por tests ausentes, vulnerabilidades, auth debug, ausencia de checksum/version y persistencia sin CAS/review. Ninguno autoriza recepción productiva ni almacenamiento automático.
   Si el proyecto selecciona Azure para recepción SFTP, materializar `MICROSOFT_AVM_SECURE_SFTP_INTAKE_PACK_PLAN.md` 1/8. Usa Microsoft AVM Storage 0.33.0 y Bicep 0.46.1 fijado por hash, SSH-only, permiso uploader `cw`, Shared Key/red pública deshabilitados y Private Endpoint. Mantener bloqueado hasta cuenta/red/DNS/Log Analytics/clave/costo, transferencia DEV y reconciliación; copiar bytes+SHA-256 a un retained-original separado antes de extracción porque Microsoft documenta incompatibilidades de HNS/SFTP con versioning/WORM.
   La auditoría object-event fija además AWS `amazon-s3-endedupe` 7209a02 como referencia exacta rechazada: 17 tests validan parte del conditional write, pero SCA 30/23, Python 3.9 deprecado, cobertura parcial, sequencers variables sin padding, unlock sin owner, lock sin lease y target sin DLQ/retry/reconciliation impiden copiar o desplegar el worker. Consultar el patrón no sustituye una implementación admitida ni sus gates target.
   AWS Powertools Python 3.34.0 queda disponible como componente runtime `PINNED_CANDIDATE_CONDITIONED`, con wheel/sdist/SLSA exactos, 314/314 parity, 189 focales PASS y SCA app 0. No lo describas como worker: Pydantic es dependencia separada, update/delete DynamoDB no tienen owner fencing y ordering/DLQ/replay/reconciliation/IaC/E2E AWS siguen obligatorios. Etiqueta toda integración local como `AUTHORED`.
   AWS DynamoDB Lock Client 1.5.0 queda disponible como componente Java `PINNED_CANDIDATE_CONDITIONED`, no como fencing externo ni worker. Usa sólo el JAR firmado/hash exacto; el grafo publicado tiene 11 advisories runtime + 4 test y no se admite. El candidate SDK 2.54.6/Log4j 2.25.5 pasa 112/OSV 0 pero es configuración `AUTHORED`, no release upstream. El UUID record version sólo sirve como CAS del lock; exige efecto fenced/transactional, DynamoDB live, order, DLQ/replay/reconciliation/IaC y conserva `holdLockOnServiceUnavailable=false` salvo decisión explícita de seguridad degradada.
   Si cualquier superficie requiere documentos/OCR/extracción, aplicar además `markdown_system/OFFICIAL_DOCUMENT_INTELLIGENCE_PROFILE.md`, crear `PROJECT_DOCUMENT_INTELLIGENCE_DECISION.md` desde su template y obtener un receipt válido del routing gate antes de seleccionar prebuilt, custom model o revisión.
   Si `platform_mode` es `OFFICIAL_PLATFORM`, exigir `VERBATIM|DEPENDENCY_PIN` para todo código de producto, adquirir una plataforma oficial coherente mediante su perfil exacto y conservar como `NO_SOURCE/BLOCKED` cada hueco que exigiría adaptación o código local. Si `platform_mode` es `BUSINESS_CENTRAL_PLATFORM`, aplicar `markdown_system/MICROSOFT_BUSINESS_CENTRAL_CAPABILITY_PROFILE.md`, componer `markdown_system/MICROSOFT_BUSINESS_CENTRAL_PLATFORM_PACK_PLAN.md`, crear `PROJECT_BUSINESS_CENTRAL_PLAN.md` y exigir config/approval hash-linked, EULA, administrador/Windows containers/Docker, demo o licencia, country view, build y test suites seleccionadas. La composición nunca ejecuta el contenedor.
9. Crear `PROJECT_PACK_PLAN.md` siguiendo `markdown_system/PACK_CONTRACT.md`.
   Para cerrar un candidate release Windows sin Git ni CI hospedado, componer `markdown_system/PORTABLE_SIGNED_RELEASE_EVIDENCE_PACK_PLAN.md`: dos builds byte-idénticos, ZIP determinista, OSV 2.5.1 exacto/0 hallazgos, SPDX 2.3, in-toto/SLSA y firma Ed25519 OpenSSH verificada. La clave real debe existir protegida fuera del proyecto; el PASS no sustituye deploy, restore, seguridad ofensiva ni aceptación.
   Para un backend empresarial compatible, partir de `markdown_system/ENTERPRISE_BACKEND_PACK_PLAN.md`, eliminar capacidades no requeridas y añadir las específicas del blueprint; nunca reconstruir manualmente una lista ya demostrada.
   Para la experiencia web opcional compatible, componer además `markdown_system/ENTERPRISE_WEB_PACK_PLAN.md`; no convertir ese adapter TypeScript en backend. El perfil 6/84 incluye la jornada pública/cliente/franquicia, aceptación quote→order y comandos operativos same-origin, Google SafeValues 1.2.0, CSP estricta con nonce fresco y gates Microsoft Playwright/Google Lighthouse: después del build, ejecutar 40 pruebas y reejecutar Playwright/Lighthouse sobre la revisión y el target exactos. Volver a probar origin/host/nonce/CSP en el edge/CDN real y añadir Trusted Types, revisión de sinks/raw HTML, journeys reales de roles/IdP/backend, accesibilidad con tecnología asistiva, RUM/carga, seguridad, deploy y rollback según el blueprint.
10. Aplicar `markdown_system/COMPOSITION_PROTOCOL.md`.
11. Para cualquier target operable, incluir `SECURE_OPERATIONS_DELIVERY_CORE.md` y `PORTABLE_CI_QUALITY_GATE_RUNNER.md`; si usa PostgreSQL, incluir también `POSTGRES_BACKUP_RESTORE_CORE.md`.
12. Instanciar las categorías CI seleccionadas por el blueprint y ejecutar los gates del pack y `ENGINEERING_EXECUTION_PLAYBOOK.md`.
13. Materializar `DEPENDENCY_LICENSE_EVIDENCE_CORE.md`, generar reportes sobre el artefacto exacto y detener la promoción si una expresión no está admitida por la política del proyecto.

El ciclo de artefactos sigue GitHub Spec Kit `v1.0.1` fijado en `9118ed15a0ba65053469a94c560ea5d233f75884`, extendido con el gate local de readiness y rigor adaptativo. Consultar `markdown_system/ELITE_PUBLIC_AGENT_METHODS.md`. La integración Git de Spec Kit es opt-in; no es necesaria para inicializar Codex o Claude. No importar plantillas o scripts upstream desde `main` ni sin licencia/notices.

## 4. Regla de código

Un agente puede materializar código desde un Markdown sólo cuando el bloque declara:

- ruta destino;
- lenguaje y versión o rango compatible;
- variables de configuración;
- dependencias directas y licencias;
- procedencia y clasificación (`AUTHORED`, `ADAPTED`, `VERBATIM`);
- invariantes y trust boundary;
- tests y comandos de verificación;
- estado de readiness del bloque.

Un ejemplo sin esos campos informa, pero no se copia automáticamente.

## 5. Estados separados

No mezclar calidad de conocimiento con facilidad de reutilización:

```text
authority: NOTE → SUPPORTED_REFERENCE → ELITE_REFERENCE
implementation: SPEC_ONLY → SNIPPET → RECONSTRUCTIBLE → REBUILD_VERIFIED
admission: DISCOVERED → LICENSE_VERIFIED → CANDIDATE → CONDITIONED → REUSABLE_PACK
```

Para generar archivos automáticamente se requiere `RECONSTRUCTIBLE`; para presentarlo como baseline comprobado se requiere `REBUILD_VERIFIED`. Sólo `REUSABLE_PACK` permite incorporación inmediata sin reabrir su expediente completo, aunque siempre se verifica compatibilidad con el proyecto.

`REBUILD_VERIFIED / CONDITIONED` no bloquea al agente: puede incorporarlo inmediatamente si verifica y registra las condiciones en el target. Los gates son aceleradores de corrección y evidencia, no pausas ceremoniales.

## 6. Configuración frente a código

Usar configuración para selección y política estable:

- mercados, monedas, locales y zonas horarias;
- jerarquías organizacionales;
- módulos y feature flags;
- roles y permisos declarativos;
- estados y transiciones;
- campos adicionales;
- selección de providers y modos.

Usar extensiones de código para semántica que exige corrección específica:

- impuestos y documentos fiscales;
- pricing, financiación y contabilidad;
- garantías, homologación y recalls;
- protocolos de fábrica, telemetría y logística;
- contratos reales de marketplaces, pagos y publicidad.

No crear ramas del core por cliente. Toda extensión entra detrás de contrato, tests y ADR.

## 7. Regla de stack

No existe un lenguaje universal de élite. El agente elige stack después del blueprint, comparando:

- ajuste al dominio y latencia;
- experiencia del equipo;
- ecosistema y soporte;
- superficie de ataque y supply chain;
- despliegue y costo operativo;
- portabilidad y salida;
- evidencia disponible en los packs.

TypeScript, Java, Go, Python, Rust, Kotlin/Swift u otros pueden ser correctos en contextos distintos. La fama del lenguaje o la empresa no decide.

La foundation de esta biblioteca es stack-neutral. TypeScript no es base ni default de backend. La vertical TypeScript existente es sólo un frontend/web adapter opcional. Para servicios empresariales nuevos, comparar primero Go y JVM; usar Rust/C++ donde seguridad de memoria, hardware o tails lo justifiquen, y registrar la decisión.

## 8. Prohibiciones

- no cargar toda la biblioteca indiscriminadamente;
- no generar desde `SPEC_ONLY` como si fuese código admitido;
- no activar integraciones sin credenciales aisladas y sandbox/contract tests;
- no copiar código público sin revisión fijada, licencia, notices y clasificación de procedencia;
- no sustituir autenticación, autorización o invariantes por instrucciones de prompt;
- no declarar “élite”, “seguro”, “baja latencia” o “production-ready” sin evidencia.

## 9. Definition of Done del sistema Markdown

La biblioteca estará lista cuando:

- los journeys de referencia se puedan componer usando sólo Markdown;
- cada archivo generado sea trazable a un bloque y una decisión;
- un workspace vacío pueda reconstruirse sin depender de `golden_starter/`;
- build, tests, migraciones y smoke tests pasen;
- licencias y notices sean reproducibles;
- una reconstrucción independiente produzca comportamiento equivalente;
- recién entonces se elimine el prototipo temporal y sus dependencias.

## 10. Recepción email documental vigente

Si el blueprint exige adjuntos por email en AWS, usar `markdown_system/AWS_SECURE_EMAIL_ATTACHMENT_PROCESSING_PACK_PLAN.md` 14/142. Materializa SES/raw/MIME, GuardDuty exact-version, Magika→AWS IDP, Powertools oficial, handoff/review, strict-field, decision worker v2, PG foundation, boundary intake+outbox, runtime Debezium y consumer inbox/proyección con CEL-Java integrado. Sólo `NO_THREATS_FOUND` exacto avanza; `Completed` exige cobertura/reviewer/tiempo; decision v2 liga snapshot exacto. El consumer exige mapping ID/clase/version/SHA, expresión single-line, schemas/corpus aprobados y claves exactas; rechaza cross-class y hace rollback de inbox+proyección si la clase o CEL fallan; sólo ACK después del commit. El gate fechado es 18/18 PostgreSQL 18.6, concurrencia 4×8 y SCA 170/0. Antes de efectos pedir y probar cuenta/MX/IAM/KMS/IDP/Aurora/Kafka/config/corpus/HITL/policy/mapping/migration/roles/RLS/costo/canaries/races/DLQ/outbox recovery/TLS/auth/rebalance/crash/restart/redelivery/order/load/restore/rollback. Nunca usar headers como tenant ni confundir staging revisado con mapping ERP final.

Los tres repositorios AWS email fijados en el source lock continúan rechazados como runtimes integrales. El receptor vigente es glue `AUTHORED` sobre contratos/SDK oficiales y no se atribuye a AWS; si el proyecto exige política `OFFICIAL_PLATFORM` estricta sin código authored, registrar `NO_SOURCE` y bloquear esta capability.

## 11. Ciclo Amazon Shipping vigente

Si el blueprint selecciona Amazon Shipping V2, usar `markdown_system/AMAZON_SPAPI_SHIPPING_TRACKING_PACK_PLAN.md` 2/48. Añade `create_claim(CreateClaimRequest)` oficial al ciclo+collection. Claim deriva tracking/package/insured value de purchase+quote+rate-request+selection, restringe proofs a hosts aprobados y persiste attempt; Amazon no expone idempotency key, por lo que ambigüedad bloquea retry. Exigir acceso/business/carrier/sandbox/costo/datos/legal/reconciliación. No prueba pickup, firma ni recepción.

Para una orden y marketplace Amazon Easy Ship demostrados, usar además `markdown_system/AMAZON_SPAPI_EASYSHIP_HANDOVER_PACK_PLAN.md` 2/34. El adapter llama `list_handover_slots`, `create_scheduled_package` y `get_scheduled_package` oficiales; el slot debe provenir del inventario hash-bound y de una aprobación separada. Schedule persiste attempt y bloquea retry ambiguo. `PickedUp` se conserva como reporte Amazon con `automatic_internal_handover_completion=false`. No inventar cancelación, firma/foto, cross-location ni carrier universal.

Para evidencia de entrega Multi-Channel Fulfillment, usar `markdown_system/AMAZON_SPAPI_FULFILLMENT_DELIVERY_EVIDENCE_PACK_PLAN.md` 2/33. Pedir cuenta/rol/scope, programa/marketplace, sandbox, cuota/costo, privacidad/retención/owner y listas observadas de status, tipos y hosts. Ligar orden, marketplace y paquetes exactos al SHA-256 del receipt empresarial; guardar sólo binarios validados por SHA-256 y respuesta filtrada. Nunca persistir URL firmada, tracking o identidad, inferir enum documental, ni convertir foto/firma/drop-off en aceptación empresarial.

Para supply sources Amazon por tienda/depósito, usar `markdown_system/AMAZON_SPAPI_SUPPLY_SOURCES_PACK_PLAN.md` 2/34. Pedir rol/scope, marketplace/programa, dynamic sandbox, cuota/costo, tratamiento de ubicación/contacto y owner. Inventariar antes de mutar; ligar approval al request, receipt y snapshot exactos; aceptar sólo modelos oficiales y reconciliar cada efecto. `UNKNOWN_EFFECT` bloquea retry y archive exige `Inactive`. Nunca convertir status Amazon en activación interna ni presentar esta API como multi-location universal.

Para publicar inventario FBM por ubicación Amazon, usar `markdown_system/AMAZON_SPAPI_MLI_INVENTORY_PACK_PLAN.md` 3/43. Exigir Product Listing, inscripción MLI, marketplace y Supply Sources activos ya inventariados. El adapter consulta `fulfillmentAvailability`, liga approval a request/receipt/snapshot exactos, reemplaza sólo con `ListingsItemPatchRequest(PRODUCT)` y `/attributes/fulfillment_availability`, escribe attempt y reconcilia por GET. Rechazo provider se sanitiza; excepción o falta de reconciliación queda `UNKNOWN_EFFECT` sin retry. No persistir seller/SKU/source IDs crudos ni hacer commit interno automático; no usar esta lane para FBA por ubicación.

Para inventario absoluto location-level de un programa Amazon External Fulfillment, usar `markdown_system/AMAZON_SPAPI_EXTERNAL_FULFILLMENT_INVENTORY_PACK_PLAN.md` 2/34. Exigir programa, rol/scope, canal, mapping location+SKU, sequence authority, sandbox, cuota/datos y owner. FETCH genera baseline filtrado; UPDATE debe usar sequence de la misma identidad, approval exacto, nuevo FETCH sin drift y batch máximo 10. Rechazo parcial, excepción o sequence no reconciliada queda `UNKNOWN_EFFECT` sin retry. Nunca persistir IDs crudos ni convertir acknowledgement en commit ERP/POS/WMS.

Para mantenimiento y reanudación, aplicar [Continuidad y atajos verificados](markdown_system/CONTINUITY_EXECUTION_SHORTCUTS.md): reutilizar evidencia compatible, verificar entradas completas y conservar gates y fallos.
