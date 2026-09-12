# Project Readiness Record Template

> Copiar como `PROJECT_READINESS_RECORD.md` dentro del proyecto. Este registro es un control `AUTHORED` de Elite; no es código atribuido a una empresa externa. Se completa conversando por las rondas de `PROJECT_START_READINESS_GATE.md`, no marcando casillas por inferencia.

Materializar primero `PROJECT_READINESS_GATE_PACK_PLAN.md`, copiar `project-readiness.template.json` como `PROJECT_READINESS_GATE.json` y mantener ambos registros consistentes. Antes de cada ronda generar `PROJECT_ADVISORY_<A-H>.md` con el renderer, explicarlo y enlazarlo como `advisory_prompt_ref`. Este Markdown conserva conversación/decisiones; el JSON es la entrada determinista. Sólo exit 0 y un reporte nuevo `READY_TO_BUILD` abren implementación.

## 1. Estado global

```yaml
readiness_version: "1.0"
project_id: ""
project_name: ""
owner: ""
operational_owners: []
status: DISCOVERY
platform_mode: BLOCK
risk_tier: CRITICAL
data_classification: RESTRICTED
target_jurisdictions: []
languages: []
currencies: []
updated_at: ""
critical_unknowns: []
required_access_not_proven: []
accepted_constraints: []
explicitly_out_of_scope: []
failure_ledger_path: "PROJECT_FAILURE_LESSONS.md"
open_high_or_critical_failures: []
dependency_update_record_path: "PROJECT_DEPENDENCY_UPDATE_RECORD.md"
open_dependency_blockers: []
authority_freshness_record_path: "PROJECT_AUTHORITY_FRESHNESS_RECORD.md"
stale_or_conflicted_authorities: []
vulnerability_monitoring_record_path: "PROJECT_VULNERABILITY_MONITORING_RECORD.md"
vulnerability_monitoring_status: ZERO_COST_BASELINE_BLOCKED
production_runtime_monitoring: NOT_READY_FOR_PRODUCTION_RUNTIME_MONITORING
next_round: A
```

Valores permitidos:

- `status`: `DISCOVERY|AWAITING_USER|ACCESS_VALIDATION|READY_TO_PLAN|READY_TO_BUILD|BLOCKED`;
- `platform_mode`: `OFFICIAL_PLATFORM|BUSINESS_CENTRAL_PLATFORM|CUSTOM_PLATFORM|BLOCK`;
- ninguna lista crítica se vacía sin enlazar respuesta o evidencia;
- `READY_TO_BUILD` se calcula aplicando el gate, nunca se escribe por deseo.

## 2. Log de conversación por rondas

| Ronda | Estado | Prompt exacto | Preguntas presentadas | Respuestas confirmadas | Contradicciones | Evidencia | Siguiente pregunta material |
|---|---|---|---|---|---|---|---|
| A — propiedad/resultado/límites | `PENDING` | `PROJECT_ADVISORY_A.md` |  |  |  |  |  |
| B — actores/journeys | `PENDING` | `PROJECT_ADVISORY_B.md` |  |  |  |  |  |
| C — modelo empresarial | `PENDING` | `PROJECT_ADVISORY_C.md` |  |  |  |  |  |
| D — documentos/extracción | `PENDING` | `PROJECT_ADVISORY_D.md` |  |  |  |  |  |
| E — integraciones | `PENDING` | `PROJECT_ADVISORY_E.md` |  |  |  |  |  |
| F — UX/canales | `PENDING` | `PROJECT_ADVISORY_F.md` |  |  |  |  |  |
| G — datos/seguridad/cumplimiento | `PENDING` | `PROJECT_ADVISORY_G.md` |  |  |  |  |  |
| H — plataforma/operación/evidencia | `PENDING` | `PROJECT_ADVISORY_H.md` |  |  |  |  |  |

Estados de ronda: `PENDING|IN_PROGRESS|ANSWERED|PROVEN|BLOCKED`. El agente presenta una ronda manejable desde el prompt exacto, explica términos, registra la respuesta textual confirmada y continúa; no descarga un cuestionario gigante, inventa respuestas ni salta rondas.

## 3. Decisiones y unknowns

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
