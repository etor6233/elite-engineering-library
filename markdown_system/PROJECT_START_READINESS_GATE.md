# Project Start Readiness Gate

## 1. Propósito y procedencia

Este contrato impide comenzar implementación antes de que el usuario y el agente hayan cerrado las decisiones, accesos y evidencias que cambian materialmente el sistema.

El ciclo ejecutable se apoya en GitHub Spec Kit `v1.0.1` (`9118ed15a0ba65053469a94c560ea5d233f75884`) para constitution → specify → clarify → plan → checklist → tasks → analyze → implement → converge. La revisión arquitectónica puede apoyarse en AWS Well-Architected Skills `v6.3.1` (`4b58ba02670abcd67458eeafb44fb6118b6af2ef`). Este inventario de preparación es composición propia de Elite Engineering Library, no código ni garantía atribuida a GitHub, AWS, Microsoft u otra empresa.

## 2. Estados

| Estado | Qué puede hacer el agente |
|---|---|
| `DISCOVERY` | leer, investigar, inspeccionar repositorios y cuentas con probes no mutantes; no crear producto |
| `AWAITING_USER` | mantener el registro, explicar opciones y pedir la siguiente decisión material; no implementar |
| `ACCESS_VALIDATION` | verificar de forma no destructiva identidades, sandboxes, scopes, cuotas y toolchains |
| `READY_TO_PLAN` | producir arquitectura, authority map, source lock y pack plan; todavía no materializar producto |
| `READY_TO_BUILD` | materializar el primer vertical slice y sus pruebas |
| `BLOCKED` | existe una condición crítica que no puede resolverse sin decisión/acceso externo; no simularla |

`READY_TO_BUILD` exige cero campos críticos `UNKNOWN`, cero accesos requeridos sin probe y una elección de plataforma empresarial. Un `NONE_WITH_REASON` explícito es válido; omitir una superficie no lo es. El estado no se acepta por texto: `project_readiness_gate/validate_project_readiness.py --project-root . --report PROJECT_READINESS_REPORT.json` debe terminar 0 sobre un report path nuevo; exit 2 conserva `BLOCKED`.

## 2.1 Alcance explícito de infraestructura de biblioteca — V402

El usuario puede encargar infraestructura reusable sin aportar cuentas/secretos
ni autorizar producción. Ese alcance debe quedar decidido expresamente y no
se infiere de un proyecto bloqueado. Para el mantenimiento V402, la decisión
canónica es `reconstruction_evidence/LIBRARY_INFRA_SCOPE_V402.md`.

El pack PROJECT-START-READINESS-VALIDATOR 0.7.0 mantiene PROJECT por defecto y
añade un registro independiente `PROJECT_LIBRARY_READINESS_GATE.json`, creado
desde `library-readiness.template.json`. Ejecutar con --record de ese nombre,
--scope LIBRARY_INFRASTRUCTURE y --level preparation o release:

- PREPARATION_PROVEN / READY_FOR_LIBRARY_WORK: autoriza las correcciones locales
  del plan explícito. Exige decisión, plan, autoridades, assurance, tooling,
  lock, notices y un receipt de preparación verificables por SHA. Los defectos
  de implementación/procedencia/dependencias se conservan OPEN con controls y
  plan de resolución exactos. T2801 final y los demás trabajos no se cierran.
- PROVEN_LOCAL / READY_FOR_LIBRARY_USE: exige T2801–T2810 y ARCA_INFRA probados
  sobre una revisión, dos materializaciones independientes fuera de la raíz,
  inventario íntegro y admisión G0–G8 por procedencia. AUTHORED sólo glue inevitable
  con justificación; BUSINESS_LOGIC local no pasa por atribuirlo a una empresa.
  Todo blocker local OPEN rechaza el cierre.

Sólo se difieren USER_CREDENTIALS (código y prueba local completos) y el expediente
exacto DAYBREAK_LIBXML2 excluido por el usuario. No pedir credenciales para cerrar
infraestructura; no confundir ausencia de credenciales con políticas pendientes,
defectos, falta de código, SCA abierta o corpus no evaluado. La evidencia declara
lo probado con fixtures y deja la aceptación real por clase/proveedor/target al
proyecto futuro. No investigar Daybreak como parte del cierre local acordado.

Ninguno de estos estados autoriza producción ni sustituye READY_TO_BUILD de un
proyecto consumidor. Sus rondas A–H, cuentas, decisiones y validación target
conservan el gate PROJECT. El informe de biblioteca siempre lleva
production_authorized=false y nunca emite READY_TO_BUILD.

## 3. Artefactos obligatorios

Para el cierre de configuración/identidad/accesos, completar la guía canónica
`PROJECT_SECRETS_TEMPLATE.md` contra código/config/IaC y registrar únicamente
referencias/estados. Las instrucciones de consola deben corresponder al proveedor
elegido y a fuentes oficiales verificadas; lo no seleccionado no exige cuenta.
En mantenimiento V295 el usuario difirió ARCA: conservar dependencia fiscal y
procedimiento posterior sin solicitar sus credenciales ni habilitar efectos.
Esta decisión no se impone silenciosamente a futuros proyectos: su intake propio
debe confirmar selección/diferimiento. Ningún dato secreto entra en el record.

Antes de implementation, componer `PROJECT_READINESS_GATE_PACK_PLAN.md`; copiar su template como `PROJECT_READINESS_GATE.json`; y hacer que el proyecto contenga:

```text
.specify/memory/constitution.md
specs/<feature>/spec.md
specs/<feature>/plan.md
specs/<feature>/tasks.md
PROJECT_READINESS_RECORD.md
PROJECT_FAILURE_LESSONS.md
PROJECT_DEPENDENCY_UPDATE_RECORD.md
PROJECT_AUTHORITY_FRESHNESS_RECORD.md
PROJECT_VULNERABILITY_MONITORING_RECORD.md
PROJECT_OFFICIAL_SOURCE_PROFILE_RECORD.md
PROJECT_BLUEPRINT.md
PROJECT_AUTHORITY_MAP.md
PROJECT_EXTERNAL_SOURCE_LOCK.md
PROJECT_PACK_PLAN.md
PROJECT_READINESS_GATE.json
PROJECT_ADVISORY_A.md
PROJECT_ADVISORY_B.md
PROJECT_ADVISORY_C.md
PROJECT_ADVISORY_D.md
PROJECT_ADVISORY_E.md
PROJECT_ADVISORY_F.md
PROJECT_ADVISORY_G.md
PROJECT_ADVISORY_H.md
PROJECT_READINESS_REPORT.json
PROJECT_EXECUTION_STATE.json
PROJECT_EXECUTION_EVENTS.jsonl
```

Spec Kit gobierna constitution/spec/clarify/plan/tasks/analyze/implement/converge. Los artefactos `PROJECT_*` de Elite agregan cobertura empresarial, fuentes públicas, accesos, packs y memoria de fallos. `PROJECT_EXECUTION_STATE.json` es sólo el cursor compacto sobre esos owners: referencia sus hashes, el próximo paso y el contexto mínimo; `PROJECT_EXECUTION_EVENTS.jsonl` conserva checkpoints append-only y no sustituye ningún artifact.

## 4. Registro normativo

Crear desde sus templates `PROJECT_READINESS_RECORD.md`, `PROJECT_FAILURE_LESSONS.md`, `PROJECT_DEPENDENCY_UPDATE_RECORD.md`, `PROJECT_AUTHORITY_FRESHNESS_RECORD.md`, `PROJECT_VULNERABILITY_MONITORING_RECORD.md` y `PROJECT_OFFICIAL_SOURCE_PROFILE_RECORD.md`. Aplicar `FAILURE_LEARNING_CONTRACT.md`, `DEPENDENCY_UPDATE_CONTRACT.md`, `AUTHORITY_FRESHNESS_AND_SELF_CORRECTION_CONTRACT.md` y `ZERO_COST_VULNERABILITY_MONITORING_PROFILE.md` desde el primer probe, no después de comenzar a programar. La plantilla de readiness contiene las ocho rondas, las 48 superficies, documentos, integraciones, accesos, plataforma, procedencia y gates. Su encabezado debe comenzar con:

```yaml
readiness_version: "1.0"
project_id: ""
owner: ""
status: DISCOVERY|AWAITING_USER|ACCESS_VALIDATION|READY_TO_PLAN|READY_TO_BUILD|BLOCKED
platform_mode: OFFICIAL_PLATFORM|BUSINESS_CENTRAL_PLATFORM|CUSTOM_PLATFORM|BLOCK
risk_tier: LIGHT|STANDARD|HIGH|CRITICAL
data_classification: PUBLIC|INTERNAL|CONFIDENTIAL|RESTRICTED
target_jurisdictions: []
updated_at: "YYYY-MM-DDTHH:MM:SSZ"
critical_unknowns: []
required_access_not_proven: []
accepted_constraints: []
```

Luego debe contener una fila por superficie:

```yaml
- id: "BUSINESS-OUTCOME"
  status: ANSWERED|PROVEN|NONE_WITH_REASON|UNKNOWN|BLOCKED
  answer_or_reason: ""
  owner: ""
  source: "user|document|contract|probe|official-upstream"
  evidence: []
  affects: []
  next_action: ""
```

El agente nunca marca `PROVEN` por inferencia. `PROVEN` enlaza un documento, contrato o probe real.

Además del inventario de 48 superficies, `PROJECT_READINESS_GATE.json` mantiene `connected_journeys`. Cada capability `REQUIRED` aparece en al menos uno; el primer vertical slice selecciona IDs existentes. Cada journey usa una versión semántica única para release, ayuda, capacitación y soporte, y enlaza evidencia local de persona/rol, interfaz, autorización, contrato API/evento, dominio, dato/efecto externo, auditoría/observabilidad, respuesta/error, ayuda contextual, capacitación, soporte, impacto de actualización, rollback/recovery, E2E y negativos. El validador rechaza eslabones vacíos, capabilities no requeridas, IDs desconocidos, campos extra y drift de versiones.

## 5. Rondas obligatorias de intake

El agente no arroja un formulario gigante sin conversación. Pregunta por rondas, actualiza el registro tras cada respuesta, detecta contradicciones y muestra brevemente qué falta. Puede agrupar preguntas estrechamente relacionadas, pero no saltar una ronda requerida.

Antes de cada ronda pendiente debe ejecutar `project_readiness_gate/render_project_advisory.py --project-root . --round <A-H> --output PROJECT_ADVISORY_<A-H>.md`, leer el prompt exacto y explicarlo al usuario. Ese prompt define términos, importancia, información que el usuario puede aportar, ejemplo no asumido, verificación y límite de `NO_APLICA`. La respuesta y su evidencia se registran separadamente; cada objeto `rounds.<A-H>` del JSON enlaza el archivo mediante `advisory_prompt_ref`. El validator compara bytes con el catálogo distribuido y rechaza ausencia, alteración, cruce de ronda, symlink, traversal u overwrite.

### Ronda A — propiedad, resultado y límites

- nombre del proyecto, owner decisor y responsables operativos;
- problema económico/operativo, resultado medible y costo de fallo;
- países/jurisdicciones, idiomas, monedas, impuestos y fecha objetivo;
- presupuesto inicial/operativo, restricciones de proveedor y lock-in aceptable;
- qué queda expresamente fuera y qué decisiones no puede tomar el agente.

### Ronda B — actores y journeys completos

- público anónimo, prospecto, cliente B2C, cliente B2B y representante;
- administrador, ventas, compras, tesorería, depósito, fábrica, logística, soporte y auditor;
- franquiciado, sucursal, territorio, proveedor, carrier, despachante y marketplace;
- alta/baja, permisos por organización/objeto/acción, delegación y break-glass;
- journeys P1 con Given/When/Then, negativos, cancelación, reversa y recuperación.

### Ronda C — modelo empresarial

- catálogo, productos, variantes, atributos, documentos, lotes/series y lifecycle;
- proveedores, RFQ/cotización, PO, PI, PL, producción, inspección y recepción;
- inventario, depósitos, ubicaciones, reservas, transferencias y costo aterrizado;
- listas de precio, monedas, promociones, presupuestos, pedidos, impuestos y facturación;
- pagos, comisiones, comprobantes, devoluciones, contracargos, conciliación y cierre;
- franquicias, territorios, sucursales, políticas, regalías y aislamiento;
- containers/pallets/shipments, tracking, aduana, discrepancias y claims;
- clientes, obras/proyectos, instalación, garantías, recalls, reclamos y postventa;
- source of truth, owner, estados, invariantes y reglas configurables de cada agregado.

### Ronda D — documentos y extracción

- canales de ingreso completos: portal/API/email/SFTP/object event/scanner/mobile/marketplace/ERP; identidad, trust boundary, checksum, partial cleanup, límites, rate y retry/idempotencia por canal;
- original inmutable y SHA-256 antes de procesar; tenant, actor, source, timestamps, nombre/tamaño/tipo declarado/detectado, engine/signature/ruleset/parser/model/schema version y decisión como provenance auditable;
- estados explícitos recepción→cuarentena→tipo/antimalware→extracción→validación/revisión→aprobación/persistencia→superseded/delete, con owner y transición autorizada;
- Magika para tipo por contenido; ClamAV o GuardDuty para antimalware; YARA-X sólo como regla adicional; firma/ruleset freshness, aislamiento, límites, outage fail-closed, release/recall e incident owner;
- archives prohibidos por defecto; si Linux es aprobado, safearchive exacto más límites de compressed/expanded bytes, ratio, entries, depth/time y rechazo traversal/symlink/hardlink/special/duplicate/Unicode collision;
- tipos documentales reales, emisores, idiomas, layouts, scans/fotos/PDFs y volumen;
- campos requeridos, tipos/unidades/monedas, relaciones, tolerancias y reglas cruzadas;
- `class_id` y `variant_id` canónicos, únicos y estables por emisor/idioma/layout/version; una clase nunca hereda la aprobación de otra;
- ground truth aprobado y conjunto representativo sin mezclar documentos sensibles sin autorización;
- costo de falso positivo/falso negativo por campo y acción permitida ante ambigüedad;
- retención, cifrado, residencia, acceso, borrado, malware/content gate y auditoría;
- proveedor permitido: local Docling/Table Transformer o servicio oficial Azure Content Understanding/Document Intelligence, Google Document AI, AWS Textract/Bedrock IDP, SAP DOX, Oracle Document Understanding u OpenAI Responses como lane agentic condicionado;
- umbrales de aceptación medidos sobre el corpus; ninguna garantía se inventa.
- si se solicita almacenamiento automático: modo `TRANSACTIONAL_DIRECT|OUTBOX_CDC_CONSUMER`; mapping con nombres exactos del consumer, expresión single-line y SHA de sus bytes UTF-8, output keys canónicas, input/output schema ID+SHA, corpus/evaluación/aprobación y owner exactos por `{class_id, variant_id}`; el JSON ejecutable rechaza duplicados, cross-class, hashes/IDs/versiones inválidos, capabilities arquitectónicas ausentes o autoridad incompleta. Para outbox, tras el PASS ejecutar `render_consumer_profile.py` con clase/variante explícitas contra el template materializado del consumer; exigir output nuevo y todavía cerrado, sin completar manualmente los diez campos ni activar pruebas.

Antes de implementar, validar `secure-file-ingestion-leaders` contra el lock exacto y completar su approval. El perfil sólo adquiere Google Magika, Cisco ClamAV, VirusTotal YARA-X, Presidio, Docling y Microsoft MarkItDown: no conecta el pipeline ni sustituye los gates del proyecto. `secure-archive-linux` es una aprobación separada y nunca se selecciona automáticamente en Windows.

Si el blueprint requiere pipeline documental, el agente materializa `DOCUMENT_PIPELINE_ROUTING_PACK_PLAN.md`; el usuario completa sus 21 clases base y agrega toda clase adicional, y el gate emite qué perfiles exactos evaluar sin clasificar ni elegir por él. Para cloud, los candidatos de source son `aws-secure-document-pipeline` (25 fuentes), `google-secure-document-pipeline` (7) o `microsoft-secure-document-pipeline` (17). El agente valida cada perfil seleccionado sin red, pide todas sus entradas —cuenta/región/identidad/storage/quarantine/costos/corpus/schemas/revisión/recovery/deploy— y sólo adquiere después de recibir un approval ligado al SHA del perfil y del lock. Varias lanes sólo se admiten como comparación explícita; no crea recursos ni transforma un sample/accelerator en producción por existir el código.

Sin corpus y ground truth, `REVESTEX_DOCUMENT_INTELLIGENCE` queda `BLOCKED` para almacenamiento automático de datos de negocio. El agente puede preparar el harness, nunca aprobar exactitud inexistente.

Una vez elegido el proveedor, el plan puede reutilizar el runtime oficial ya reconstruido: `AZURE_DOCUMENT_RUNTIME_PACK_PLAN.md`, `GOOGLE_DOCUMENT_RUNTIME_PACK_PLAN.md` o `AWS_TEXTRACT_DOCUMENT_RUNTIME_PACK_PLAN.md`. Google incluye un adapter GCS condicionado y AWS uno S3 Object Lock condicionado; su presencia no responde región, cuenta, costo, corpus, processor/analyzer, IAM/KMS, bucket ni política de retención. Esos campos continúan en `required_access_not_proven` y el agente no llama ni configura el provider hasta probarlos.

### Ronda E — integraciones

Para cada integración clasificar `REQUIRED`, `OPTIONAL` o `NONE_WITH_REASON` y cerrar:

- proveedor/producto, país, cuenta owner, sandbox/test account y producción;
- API/SDK/version, scopes, OAuth/service identity, webhook, firma y allowlists;
- quotas/rate limits, costos/caps, retries, idempotencia, polling y reconciliation;
- data shared, región/retención, términos/DPA, deletion y exit plan;
- test fixture, contract probe, failure mode, runbook y contacto de soporte.

Cuando una captación de lead deba continuar por WhatsApp, SMS, email, chat u otro canal, cerrar además: identificador estable de ambos lados; método verificable que autoriza el binding; subject/lead/organización esperados; policy y hash de evidencia; finalidad/PII/opt-out; secreto HMAC y plan de rotación; revocación/rebinding; receipt, estado incierto y reconciliación outbound. No asociar identidades por similitud de nombre, teléfono, email ni texto del modelo. NIST SP 800-63 Rev. 4 gobierna el principio de cuenta/sujeto/binding y registro; `GO-PG-CONTACT-CHANNEL-IDENTITY` materializa el borde local condicionado.

Inventario mínimo: Mercado Libre, Amazon SP-API, Mercado Pago, Stripe, Google Ads, Merchant Center, Meta Ads, TikTok Ads, WhatsApp, email/SMS/push, carriers/tracking, bancos/fintech, contabilidad/facturación argentina y proveedor IA/OCR incluido Grok.

Perfiles ejecutables disponibles sólo después de cerrar esa fila: `PAYMENT_WEBHOOK_ADAPTERS_PACK_PLAN.md` para Stripe/Mercado Pago, `MERCADOLIBRE_MARKETPLACE_PACK_PLAN.md` para item/order read, validator y notification fetch sin publicación, `GO-MERCADOLIBRE-QUESTION-OUTBOUND` sobre `GO-PG-OUTBOUND-DELIVERY-FENCE 0.2.x` para Questions con aprobación hash-bound, un POST y confirmación/reconciliación GET-only, `AWS_ENTERPRISE_ADAPTERS_PACK_PLAN.md` para S3/SES, `AWS_SECURE_DOCUMENT_INTAKE_PACK_PLAN.md` para portal/API→cuarentena S3→gate local condicionado, `GOOGLE_ADS_REPORTING_PACK_PLAN.md`, `META_ADS_REPORTING_PACK_PLAN.md` y `TIKTOK_ADS_REPORTING_PACK_PLAN.md` para reporting read-only, `META_WHATSAPP_CLOUD_PACK_PLAN.md` para template send/webhook, `FIREBASE_PUSH_PACK_PLAN.md` para FCM dry-run/send condicionado, `AMAZON_SPAPI_CATALOG_PACK_PLAN.md` para Amazon Catalog read-only y `GOOGLE_MERCHANT_PRODUCT_SYNC_PACK_PLAN.md` para ProductInput/status con write externo explícito. Mercado Libre exige app/owner/seller/token/scopes/origin/policies/aprobación y excluye sus SDKs archivados; TikTok requiere adquisición aprobada del exacto source commit; WhatsApp requiere licencia/terms/consent/version/template exactos; Firebase requiere proyecto/ADC/API/dispositivo consentido/políticas y aprobación separada del send real. Si la capacidad requerida excede esos claims —Mercado Libre publish/update/stock/precio/OAuth refresh/postventa; Amazon Orders/RDT/Listings; mutaciones Ads; WhatsApp conversations/media/Flows no cubiertos; push topics/APNs/WebPush no cubiertos; carriers, bancos o fiscalidad— la fila conserva `BLOCKED` o un trabajo `AUTHORED`/`ADAPTED` explícitamente autorizado; nunca se sustituye por otro provider en silencio.

Si el proyecto requiere refunds en el backend empresarial, la fila de pago debe seleccionar proveedor/cuenta/sandbox/moneda, secretos por referencia, webhook/reconciliation y política de allocation. `ENTERPRISE_BACKEND_PACK_PLAN.md` incluye `GO-OFFICIAL-RETURN-REFUND-WORKER` 0.1.0: sólo admite automáticamente una línea quantity-one ligada al stock y un pago elegible; split payments, descuentos/impuestos/envío sin regla aprobada o Mercado Pago fuera de ARS conservan `BLOCKED`. Exigir contratos SDK y PostgreSQL 0001–0020 desde Markdown; ningún estado no final confirma devolución de dinero.

Si requiere exchange, registrar variante admitida, organization/local delivery, settlement policy y accounting owner. El backend incluye `GO-RETURN-EXCHANGE-FULFILLMENT-WORKER` 0.1.0, pero sólo automatiza un pedido original entregado de una línea quantity-one, misma variante/sucursal, inventory effect exitoso y accounting request existente. Exigir PostgreSQL 0001–0021, reserva concurrente, retry sin stock, inmutabilidad y down/up desde Markdown. Multi-line, diferencia de precio, cross-organization, carrier o fiscalidad no demostrados conservan `BLOCKED`; `prepared` no es `delivered`.

### Ronda F — UX, contenido y canales

- web pública, catálogo, buscador, contacto/leads, SEO, analytics y consent;
- portales cliente/B2B/franquicia/admin/fábrica/proveedor y sus journeys;
- responsive, navegadores/dispositivos, accesibilidad y performance budgets;
- branding, contenido legal/comercial, assets con licencia y responsables editoriales;
- mobile/desktop/offline/hardware sólo si el journey lo exige.

### Ronda G — datos, seguridad y cumplimiento

- clasificación por campo, PII/financiero/secreto, propósito y base legal;
- consentimiento, retención, exportación, rectificación, borrado y legal hold;
- tenant/organization isolation, policy owner, MFA/federation y session lifecycle;
- secretos/KMS/PKI, rotación, incident access y separación DEV/TEST/PROD;
- threat/abuse cases, fraude, scraping, takeover, webhook forgery y supply chain;
- logging/audit sin datos sensibles, seguridad ofensiva y acceptance owner.

### Ronda H — plataforma, operación y evidencia

- decisión `BUSINESS_CENTRAL_PLATFORM`, `CUSTOM_PLATFORM` o `BLOCK`;
- cloud/on-prem/híbrido, regiones, DNS/TLS, identities, network boundaries y IaC;
- entornos local/CI/DEV/STAGE/PROD y promoción;
- SLI/SLO, carga, p50/p95/p99, capacidad, goodput y costo unitario;
- observabilidad, on-call, alerts, runbooks e incident owner;
- backup/PITR, restore, RPO/RTO, DR, canary, rollback y rollback owner;
- test matrix: unit/property/contract/integration/E2E/a11y/load/security/recovery;
- definición de piloto y evidencia que autoriza producción.
- presupuesto de seguridad/CI, GitHub Team/Dependabot, OSV local, periodicidad autorizada y monitor runtime requerido antes de producción.

## 6. Validación real de accesos

Para un primer slice únicamente CLI/AUTOMATION, PROJECT-START-READINESS-VALIDATOR 0.6.1 acepta esos journeys seleccionados sin imponer web/API/mobile ficticios. Exige vincular RUNTIME, CONTRACTS, DOMAIN-MODULES y AUTHORIZATION, además de todos los refs, roles, evidencia y versiones de connected_journeys. Una superficie web/portal/mobile/desktop/API declarada siempre exige su capability correspondiente. No es una excepción a seguridad, aceptación, ayuda o recuperación, ni una aprobación de este mantenimiento por existir sus archivos.

El usuario no pega secretos, claves privadas, tokens o certificados dentro de Markdown ni del chat. Debe proporcionar acceso por un mecanismo seguro del entorno. El agente registra sólo:

```yaml
provider: ""
environment: sandbox|test|production
identity_reference: "nombre lógico, nunca secreto"
required_scopes: []
probe_command_or_request: "redactado"
probe_expected: ""
probe_observed: ""
proven_at: "YYYY-MM-DDTHH:MM:SSZ"
evidence_path: ""
status: PROVEN|FAILED|NOT_PROVIDED|NOT_APPLICABLE
```

Un probe:

- es read-only o usa un recurso sandbox desechable autorizado;
- no imprime secretos ni PII;
- confirma identidad/cuenta, ambiente, scopes, región y contrato API;
- captura request ID/status/schema redacted;
- falla cerrado si responde el ambiente equivocado;
- no convierte acceso sandbox en evidencia de producción.

## 7. Elección obligatoria de plataforma empresarial

### `OFFICIAL_PLATFORM`

Seleccionar cuando el propietario exige que el código de producto provenga únicamente de implementaciones publicadas por empresas líderes y admitidas por revisión, licencia, hash y pruebas. Cada capability debe usar `VERBATIM` o `DEPENDENCY_PIN` contra una revisión exacta de `PROJECT_EXTERNAL_SOURCE_LOCK.md`. `ADAPTED` y `AUTHORED` quedan prohibidos para producto; los registros, approvals, materializadores y gates de Elite siguen siendo controles locales y nunca se atribuyen al upstream.

El agente elige una plataforma oficial coherente —no mezcla fragmentos de proveedores por fama—, adquiere su perfil exacto y ejecuta los comandos upstream en el entorno soportado. Un sample, accelerator o branch condicionado conserva sus límites y no se convierte en producción por elegir este modo. Si falta una capability requerida o sería necesario modificar semántica, registrar `NO_SOURCE` y mantenerla `BLOCKED`; sólo una decisión explícita posterior del propietario puede cambiar a `CUSTOM_PLATFORM`.

### `BUSINESS_CENTRAL_PLATFORM`

Seleccionar sólo si el usuario acepta Microsoft Business Central/AL y puede aportar sandbox/responsables. Para desarrollo local, componer `MICROSOFT_BUSINESS_CENTRAL_PLATFORM_PACK_PLAN.md`: fija BcContainerHelper 6.1.16 y artifact 28.4.53241.0, pero exige EULA, administrador, Windows containers/Docker, capacidad, credenciales demo y approval hash-linked antes de crear nada. El snapshot integral auditado es el commit Microsoft firmado `microsoft/BCApps@2eae56d704a1fd035d104f333602aea7091b7749`, bajo MIT y condiciones por path: fija 36.673 AL, Inventory 958, Warehouse 355 y 22 tests AL SCM-Reservation. Proviene de `main`; el tag firmado `releases/28.4/StrictMode@cb07eef12935e07dd4258c122d0e7b1920fc1223` no incluye `src/Layers/W1/BaseApp`, APIV2 ni tests completos. Exigir una release integral o aceptar explícitamente el snapshot de desarrollo antes de compilar source. No portar silenciosamente su dominio a Go, no afirmar que el runtime es MIT y no asumir localización argentina.

Aplicar obligatoriamente `MICROSOFT_BUSINESS_CENTRAL_CAPABILITY_PROFILE.md` y crear `PROJECT_BUSINESS_CENTRAL_PLAN.md`. Una aceptación del snapshot no reemplaza los probes de licencia/runtime/container, la country view, las suites oficiales por área ni los adapters/localizaciones faltantes.

### `CUSTOM_PLATFORM`

Seleccionar si se construirá sobre la foundation portable/Go/PostgreSQL de Elite. Los SDKs oficiales del ledger pueden ser `VERBATIM` o dependencias fijadas; reglas REVESTEX, adapters y glue nuevos son `AUTHORED`. No atribuirlos a las empresas fuente.

### `BLOCK`

Seleccionar si el propietario prohíbe tanto una plataforma empresarial condicionada como cualquier código específico nuevo. El agente continúa investigación y admisión, pero no implementa dominio.

## 8. Gate de fuente pública

Por cada capability `REQUIRED`, `PROJECT_EXTERNAL_SOURCE_LOCK.md` registra:

```yaml
capability: ""
source_owner_repo: ""
revision: ""
paths_or_package: []
license_expression: ""
notices: []
provenance_mode: VERBATIM|DEPENDENCY_PIN|ADAPTED|AUTHORED|NO_SOURCE
claim: ""
non_claims: []
tests_to_run: []
access_conditions: []
admission: ""
```

`NO_SOURCE` bloquea cualquier afirmación de código público de élite. En `OFFICIAL_PLATFORM`, el producto admite sólo `VERBATIM|DEPENDENCY_PIN`; `ADAPTED|AUTHORED` permanecen bloqueados. `AUTHORED` es válido sólo si el usuario eligió `CUSTOM_PLATFORM`; nunca se renombra como upstream.

Cuando una capability `REQUIRED` quede sin una fuente compatible admitida, no basta escribir `NO_SOURCE`: materializar `CAPABILITY_GAP_RESOLUTION_PACK_PLAN.md`, crear un expediente por capability y ejecutar `CAPABILITY-GAP-RESOLUTION-GATE`. El agente investiga autoridades, docs, repos, releases y advisories oficiales actuales; fija evidencia/revisión/hash/licencia/claim y G0–G8; continúa si resulta `RESEARCH_INCOMPLETE`. Sólo `USE_REUSABLE_PACK` habilita implementación inmediata. Cualquier estado condicionado o bloqueado conserva evidencia, blockers y trigger de reapertura y actualiza owners canónicos sin atribuir código local al proveedor.

Cuando el target sea privado Go/TypeScript y no exista entitlement GitHub Code Security, el baseline sin costo es `MICROSOFT_DEVSKIM_ADAPTED_SAST_PACK_PLAN.md`: debe reconstruirse con su commit/árbol, adaptación, tests y SCA exactos, ejecutar sobre el scope real y registrar triage. Su estado `ADAPTED / CONDITIONED` no cierra análisis interprocedural, fuzzing, DAST, threat model ni ofensiva.

Cuando una capability use uno de los perfiles materializados en `elite_sources/source-profiles/`, completar `PROJECT_OFFICIAL_SOURCE_PROFILE_RECORD.md` antes de ejecutar adquisición no validante. Todos sus `required_user_inputs` deben tener respuesta y cada `production_blockers` debe permanecer abierto o enlazar evidencia de cierre. `apply_source_profile.ps1` es el único runner admitido para el perfil completo.

## 9. Condición exacta de `READY_TO_BUILD`

El agente demuestra:

- constitution/spec y readiness record consistentes;
- ocho prompts A–H generados desde el catálogo exacto, explicados por ronda y enlazados mediante `advisory_prompt_ref`;
- 48 superficies de `TOTAL_SYSTEM_CAPABILITY_CONTRACT.md` clasificadas;
- journeys P1, invariantes, owners y acceptance tests cerrados;
- plataforma elegida y accesible;
- integraciones requeridas con sandbox/contract probe o plan explícitamente diferido;
- corpus/ground truth e inventario clase/variante presentes cuando hay document intelligence; almacenamiento automático además con modo de persistencia, mapping compatible con el consumer, expresión/schema hash-locked, output keys, capabilities, evaluación y aprobación por clase;
- source lock con revisión/licencia/notices y sin branch móvil;
- inventario de dependencias/toolchains/upstreams con policy, EOL, owner, locks y rollback; todo candidato aplicable está evaluado o bloqueado explícitamente;
- authorities materiales con versión/fecha/applicability y cero `STALE_BLOCKED`/`CONFLICT_UNRESOLVED` que afecten el primer vertical slice;
- datos, seguridad, privacidad, SLO, backup/restore y rollout definidos según riesgo;
- pack plan compatible, sin colisiones y con rollback;
- ledger de fallos creado, lecciones previas aplicables consultadas y cero fallos `HIGH`/`CRITICAL` abiertos que afecten el primer vertical slice;
- cero critical unknowns y cero required access not proven.

Sólo entonces cambia `status: READY_TO_BUILD` y comienza un vertical slice P1. “Completar el formulario” sin evidencia no abre el gate.

## 10. Continuación de un proyecto existente

El protocolo se adapta: el agente primero descubre código, IaC, contratos, datos, proveedores y evidencia existentes; no obliga a reiniciar. Toda respuesta demostrable se marca `PROVEN`, las contradicciones se abren como blockers y el plan se limita al delta. El execution state debe declarar `mode: EXISTING`, enlazar `baseline.inventory_ref` y `baseline.delta_scope_ref`, clasificar cada task del `tasks.md` exacto y checkpointar sólo artifacts observados. El mismo gate impide ampliar una zona crítica cuyo owner, contrato o acceso siga desconocido.
