# Business function operating contracts — V403

## 1. Metadata

```yaml
pack_id: "BUSINESS-FUNCTION-OPERATING-V403"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: RECONSTRUCTIBLE
  admission: CONDITIONED
claim: "Eleven function contracts, 22 task mappings, no implicit grants, scoped method references, existing owner bindings, and local validation; no eleven-business-operation or Galaxy-event implementation claim"
stacks: ["Python 3.14.4 standard library", "PowerShell 7+ existing materializer"]
compatible_with: ["ENGINEERING_EXECUTION_VALIDATOR 1.3.1", "PROJECT-OPERATING-CONNECTION 0.1.0", "V402 library baseline337"]
incompatible_with: ["implicit business-title privileges", "inherited maintenance progress or approvals", "product readiness inferred from scaffold", "invented Galaxy event methods"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["No upstream code copied; eight official documentary references in roles/method-sources.v403.json"]
verified_at: "2026-09-13"
```

## 2. Applicability

Reusable library scaffolding for Grok Bot 101, Engineering, Product Managers, Founders, Sales Engineering, Sales, SDRs, Customer Support, Marketing Operations, Post-Sales and Marketing. These are task functions, not authorization roles. Nine functions reference existing packs with historical local evidence; two are documentary only. Every consumer resolves its own existing owners and product policies. No credentials or real business data are needed to validate this scope.

## 3. Architecture contract

AUTHORED contract, validation and owner-binding glue only. specification_status, method_admission, galaxy_admission and verification_status remain independent. Methods from current official documentation can inform scoped tasks; future Galaxy material remains WAITING_GALAXY_ADMISSION. The catalogue executes no business effects. The consumer's authorized product operations retain their own gates and enforcement.

Project owners remain external. MAINTENANCE_FIXTURE names the V403 example. Consumer binding requires explicit existing owner paths and rejects maintenance IDs/task paths; it resets verification evidence and grants without mutating state/events. The existing execution validator must independently validate resume. Receipt hashes check correspondence, not truth of semantic assertions.

## 4. Exact file manifest

```text
CREATE roles/README.md
CREATE roles/business-functions.v403.json
CREATE roles/method-sources.v403.json
CREATE roles/library-bindings.v403.json
CREATE roles/validate_business_functions.py
CREATE roles/bind_consumer_owners.py
CREATE roles/tests/test_business_functions.py
CREATE roles/evidence/source-observations.json
```

## 5. Materialization blocks

### FILE: `roles/README.md`

```yaml
block_id: "BUSINESS-FUNCTION-OPERATING-V403:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "local contract, provenance metadata, validator or owner-binding glue; no upstream method code"
license: "LicenseRef-Workspace-Owner"
sha256: "29bb58cb871fd7325006083bbcbc852147f9ebec4b82f0e18290114e760a9832"
variables: []
secrets_allowed: false
```

````markdown
# Contratos de funciones empresariales V403

Este directorio implementa un scaffold de contratos y validación para 11 funciones / 22 tareas. Es AUTHORED glue de biblioteca. No implementa once operaciones empresariales ni atribuye su código a Microsoft, DORA, HubSpot o xAI.

Las nueve referencias a packs existentes enlazan archivos y evidencia histórica de la biblioteca. Sus hashes se verifican contra el baseline 337. La consola no llama sus APIs; los límites de cada pack siguen vigentes.

## Lectura y frontend

business-functions.v403.json es el contrato consumido por la vista. Cada función contiene propósito, sistemas, tareas, inputs, outputs, permisos requeridos, aceptación, sources, owners, resume y binding de implementación. La ruta es #/functions/{id}. La cobertura de todas las secciones y tareas exige además prueba del frontend real; el validator de C sólo verifica sus vínculos declarados.

Grok Bot 101 es LEARNING_TRACK. Los demás nombres son BUSINESS_FUNCTION, nunca roles de autenticación. effective_grants permanece vacío y ninguna task concede permisos. Nombres propuestos se distinguen de marketing:read / marketing:request observados en el owner existente. El producto consumidor resuelve sus propios subjects, objetos, organizaciones y enforcement.

Los efectos están deshabilitados en este catálogo. Esa restricción no prohíbe operaciones de un producto que cuenten con autorización explícita y pruebas propias.

## Estados y DONE

- specification_status=SPEC_READY: el contrato local está definido.
- method_admission=OFFICIAL_METHOD_REFERENCED: método documental actual y claim limitado; no es admisión de código o arquitectura.
- galaxy_admission=WAITING_GALAXY_ADMISSION: sólo falta material futuro del evento. No impide usar las referencias documentales actuales.
- verification_status.target=NOT_TESTED: no hay prueba del negocio ni del target futuro.
- El JSON portable conserva scaffold=NOT_TESTED como declaración de entrada sin autoatestación. Los resultados ejecutados viven en receipts externos ligados al SHA exacto del contrato; la vista debe mostrar ambos alcances.
- PROVEN exige evidencia existente con result=PASS, claim_scope=FUNCTION_CONTRACT_SCAFFOLD, hash exacto y función incluida. Un hash acredita correspondencia, no la veracidad semántica de un informe escrito por un tercero.
- DONE del scaffold exige checks locales PASS, cobertura frontend independiente y cero defectos conocidos abiertos en ese alcance y revisión. CONDITIONED / NOT_TESTED no cuentan como PASS de una garantía requerida. Los owners canónicos conservan el triage; el validator no inventa un conteo de bugs.

## Owners y reanudación

El archivo entregado tiene owner_binding.scope=MAINTENANCE_FIXTURE. Referencia los seis owners existentes de este mantenimiento V403. Ningún proyecto consumidor debe heredar esas rutas, su identidad, progreso o aprobaciones.

Después de materializar el scaffold en el consumidor, preparar un owner-map JSON con claves progress, execution_state, execution_events, failures, dependency y freshness. Sus valores son rutas relativas de archivos YA existentes del proyecto. Resolver progress al tasks owner real (por ejemplo el PROJECT_TASKS_REF del consumidor), sin crear otro backlog.

Ejemplo de integración, con rutas elegidas por el consumidor:

    python roles/bind_consumer_owners.py --source roles/business-functions.v403.json --project-root . --owner-map owner-map.json --output roles/business-functions.consumer.json

El binder crea sólo un contrato nuevo, rechaza destinos existentes, rutas fuera de raíz, IDs/owners de mantenimiento y owners ausentes. No modifica state/events ni copia progreso. Resetea claims/evidencias de verificación y no concede permisos. Después ejecutar el validador de execution state del kit 1.3.1 del consumidor; OWNERS_BOUND_ONLY no autoriza resume.

El gate de contrato acepta --contract roles/business-functions.consumer.json. Sus referencias al kit y a los packs se resuelven contra --library-root. Una modificación material del contrato o de los refs invalida los receipts afectados.

## Comprobación local

Desde la raíz de un mantenimiento o consumidor con owners resueltos:

    python roles/validate_business_functions.py --project-root . --library-root "C:/ruta/biblioteca" --report roles/evidence/validation-new.json
    python roles/tests/test_business_functions.py --library-root "C:/ruta/biblioteca"

El report path debe ser nuevo. Los tests crean fixtures sintéticas aisladas bajo roles/evidence y no modifican los owners. Cubren privilegios implícitos, scope por objeto/organización, links, procedencia, Galaxy futuro, PROVEN no ligado, cambios de pack y herencia de owners. No ejercitan proveedores, BFF, usuarios humanos ni operaciones de negocio.

## Fuentes y procedencia

method-sources.v403.json contiene ocho fuentes oficiales actuales con consulta, claim y límite. Son enlaces y paráfrasis breves propias. No se adquirió código upstream ni se seleccionó Business Central por citar Dynamics 365.

evidence/source-observations.json conserva URL, fecha de consulta, HTTP200, tamaño y huella SHA-256 de la respuesta observada. No retiene cuerpos HTML: es evidencia histórica de observación, no snapshot reproducible ni firma editorial. Revalidar la fuente oficial si cambia una decisión material; no tratar el hash histórico como prueba de contenido actual.

library-bindings.v403.json sí permite comparar bytes locales de packs/evidencia existentes. Estas referencias no promueven el estado del producto consumidor.

## Reconstrucción

La guía canónica es markdown_system/BUSINESS_FUNCTION_OPERATING_CONTRACT_V403.md. La única fuente de código es implementation_packs/BUSINESS_FUNCTION_OPERATING_V403.md y usa los bloques y hashes estándar de la biblioteca.

Materializar con el script ya existente, sin crear otro compositor:

    pwsh -NoProfile -File materialize_markdown_pack.ps1 -PackFile implementation_packs/BUSINESS_FUNCTION_OPERATING_V403.md -Destination C:/ruta/ausente

El pack contiene ocho archivos. Conserva sus receipts y vuelve a ejecutar el gate con los owners del mantenimiento o consumidor. No copia state/events ni aprobaciones. Python usa sólo la biblioteca estándar; JSON, validator, binder y tests son AUTHORED glue. No se instala dependencia ni se atribuye código local a un proveedor.
````

### FILE: `roles/business-functions.v403.json`

```yaml
block_id: "BUSINESS-FUNCTION-OPERATING-V403:file2:v1"
operation: CREATE
provenance: AUTHORED
source: "local contract, provenance metadata, validator or owner-binding glue; no upstream method code"
license: "LicenseRef-Workspace-Owner"
sha256: "2598b2b8ac1259c7ad640c02330bfe92be212495438235394b6766adba80585e"
variables: []
secrets_allowed: false
```

````json
{
  "schema_version": "1.0.0",
  "contract_id": "BUSINESS-FUNCTIONS-V403",
  "scope": "LIBRARY_INFRASTRUCTURE",
  "claim_scope": "FUNCTION_CONTRACT_SCAFFOLD",
  "production_authorized": false,
  "provenance": {
    "classification": "AUTHORED",
    "kind": "CONTRACT_AND_VALIDATOR_GLUE",
    "upstream_code_copied": false,
    "official_method_code_claim": false
  },
  "frontend_contract": {
    "route_pattern": "#/functions/{id}",
    "mode": "READ_ONLY_LIBRARY_CONSOLE",
    "required_sections": [
      "purpose",
      "systems",
      "tasks",
      "permissions",
      "allowed",
      "prohibited",
      "inputs",
      "outputs",
      "acceptance",
      "implementation",
      "sources",
      "status",
      "owners",
      "resume"
    ],
    "external_effects_enabled": false
  },
  "owner_refs": {
    "progress": "specs/library-extension/tasks.md",
    "execution_state": "PROJECT_EXECUTION_STATE.json",
    "execution_events": "PROJECT_EXECUTION_EVENTS.jsonl",
    "failures": "PROJECT_FAILURE_LESSONS.md",
    "dependency": "PROJECT_DEPENDENCY_UPDATE_RECORD.md",
    "freshness": "PROJECT_AUTHORITY_FRESHNESS_RECORD.md"
  },
  "permission_policy": {
    "track_is_auth_role": false,
    "default": "DENY",
    "effective_grants": [],
    "requires": [
      "authenticated_subject",
      "explicit_product_permission",
      "organization_scope",
      "resource_scope",
      "server_side_enforcement",
      "audit_evidence"
    ],
    "unbound_permission_behavior": "BLOCK_PRODUCT_EFFECT",
    "inherited_privilege_allowed": false
  },
  "done_policy": {
    "scaffold": "All declared local checks PASS on fixed bytes; zero known open defects in this scaffold scope; exact frontend coverage independently checked.",
    "product": "Target gates required before runtime/production claims. CONDITIONED and NOT_TESTED do not count as PASS.",
    "galaxy": "Waiting for future event source does not block referenced current official methods; it blocks claims derived from that event."
  },
  "functions": [
    {
      "id": "grok-bot-101",
      "track_name": "Grok Bot 101",
      "function_kind": "LEARNING_TRACK",
      "auth_role": false,
      "purpose": "Aprender a delimitar una tarea, sus datos y la evidencia antes de delegarla a un agente.",
      "systems": [
        "Documentación oficial Grok",
        "Contratos y ejemplos locales sin datos reales"
      ],
      "allowed": [
        "Leer contratos, fuentes y evidencia local.",
        "Preparar propuestas sin efectos externos dentro del scope aprobado."
      ],
      "prohibited": [
        "Derivar privilegios del título o track.",
        "Usar secretos o datos empresariales reales en el scaffold.",
        "Enviar mensajes, publicar, gastar, modificar CRM o ejecutar efectos desde esta consola."
      ],
      "status": {
        "specification_status": "SPEC_READY",
        "method_admission": {
          "status": "OFFICIAL_METHOD_REFERENCED",
          "claim_scope": "DOCUMENTARY_METHOD_ONLY",
          "source_refs": [
            "XAI-GROK",
            "MS-RBAC"
          ],
          "code_admission": "NONE"
        },
        "verification_status": {
          "scaffold": "NOT_TESTED",
          "target": "NOT_TESTED",
          "evidence_refs": []
        }
      },
      "source_refs": [
        "XAI-GROK",
        "MS-RBAC"
      ],
      "galaxy_admission": {
        "status": "WAITING_GALAXY_ADMISSION",
        "source_ref": null,
        "applies_only_to": "FUTURE_EVENT_MATERIAL",
        "content_invented": false,
        "reopen_trigger": "Material oficial identificable del evento disponible y obtenido legalmente; fijar fecha/revisión/hash/términos/claim y ejecutar admisión aplicable antes de incorporarlo."
      },
      "owner_refs": {
        "progress": "specs/library-extension/tasks.md",
        "execution_state": "PROJECT_EXECUTION_STATE.json",
        "execution_events": "PROJECT_EXECUTION_EVENTS.jsonl",
        "failures": "PROJECT_FAILURE_LESSONS.md",
        "dependency": "PROJECT_DEPENDENCY_UPDATE_RECORD.md",
        "freshness": "PROJECT_AUTHORITY_FRESHNESS_RECORD.md"
      },
      "resume": {
        "read_refs": [
          "PROJECT_EXECUTION_STATE.json",
          "PROJECT_EXECUTION_EVENTS.jsonl",
          "specs/library-extension/tasks.md",
          "PROJECT_FAILURE_LESSONS.md",
          "roles/business-functions.v403.json"
        ],
        "rule": "Validar checkpoint y hashes; consultar los owners existentes; invalidar evidencia afectada por delta; registrar siguiente acción allí, sin backlog paralelo."
      },
      "implementation_binding": {
        "status": "DOCUMENTED_ONLY",
        "pack": null,
        "evidence": null,
        "action": null,
        "kind": "NO_PRODUCT_RUNTIME",
        "target_verified": false,
        "v403_executes_binding": false
      },
      "open_requirement": {
        "id": "TARGET-GROK-BOT-101",
        "scope": "FUTURE_CONSUMER_ONLY",
        "description": "Admitir el material oficial Galaxy cuando exista y se obtenga; una cuenta o runtime Grok queda fuera de este scaffold.",
        "owner_ref": "specs/library-extension/tasks.md"
      },
      "tasks": [
        {
          "id": "grok-bot-101.delimit-learning-exercise",
          "title": "Delimitar un ejercicio de aprendizaje",
          "inputs": [
            "Objetivo de práctica",
            "ejemplo sintético",
            "límites de datos y herramientas"
          ],
          "outputs": [
            "Ficha de tarea y preguntas abiertas"
          ],
          "product_permission": {
            "id": "learning:read",
            "definition_kind": "PROPOSED_TARGET_REQUIREMENT",
            "granted": false,
            "subject_binding": "REQUIRES_TARGET_CONFIGURATION",
            "organization_scope": "REQUIRED",
            "resource_scope": "REQUIRED",
            "enforcement": "EXISTING_PRODUCT_OWNER_OR_TARGET_GAP",
            "policy_ref": "SECURITY_SRE_CLOUD_INFRASTRUCTURE.md"
          },
          "frontend": {
            "route": "#/functions/grok-bot-101",
            "section": "tasks",
            "task_id": "grok-bot-101.delimit-learning-exercise",
            "interaction": "READ_CONTRACT",
            "effect_enabled": false
          },
          "acceptance": [
            {
              "id": "grok-bot-101.AC1",
              "claim_scope": "CONTRACT_SPECIFICATION",
              "criterion": "Ejercicio sin credenciales, identidad ni efectos; distinguir explicación, simulación y resultado observado.",
              "test_ref": "roles/tests/test_business_functions.py",
              "evidence_refs": [],
              "target_status": "NOT_TESTED"
            }
          ],
          "owner_ref": "specs/library-extension/tasks.md"
        },
        {
          "id": "grok-bot-101.review-agent-output",
          "title": "Revisar la salida propuesta por un agente",
          "inputs": [
            "Ficha de tarea",
            "salida sintética",
            "criterio de revisión"
          ],
          "outputs": [
            "Comparación con el objetivo y límite observado"
          ],
          "product_permission": {
            "id": "learning:review",
            "definition_kind": "PROPOSED_TARGET_REQUIREMENT",
            "granted": false,
            "subject_binding": "REQUIRES_TARGET_CONFIGURATION",
            "organization_scope": "REQUIRED",
            "resource_scope": "REQUIRED",
            "enforcement": "EXISTING_PRODUCT_OWNER_OR_TARGET_GAP",
            "policy_ref": "SECURITY_SRE_CLOUD_INFRASTRUCTURE.md"
          },
          "frontend": {
            "route": "#/functions/grok-bot-101",
            "section": "tasks",
            "task_id": "grok-bot-101.review-agent-output",
            "interaction": "READ_CONTRACT",
            "effect_enabled": false
          },
          "acceptance": [
            {
              "id": "grok-bot-101.AC2",
              "claim_scope": "CONTRACT_SPECIFICATION",
              "criterion": "Un texto convincente no se registra como efecto ejecutado ni como permiso.",
              "test_ref": "roles/tests/test_business_functions.py",
              "evidence_refs": [],
              "target_status": "NOT_TESTED"
            }
          ],
          "owner_ref": "specs/library-extension/tasks.md"
        }
      ],
      "owner_binding_scope": "MAINTENANCE_FIXTURE"
    },
    {
      "id": "engineering",
      "track_name": "Engineering",
      "function_kind": "BUSINESS_FUNCTION",
      "auth_role": false,
      "purpose": "Convertir un cambio acotado en implementación revisable y evidencia reproducible.",
      "systems": [
        "Spec Kit/owners de ejecución existentes",
        "Packs y gates de biblioteca"
      ],
      "allowed": [
        "Leer contratos, fuentes y evidencia local.",
        "Preparar propuestas sin efectos externos dentro del scope aprobado."
      ],
      "prohibited": [
        "Derivar privilegios del título o track.",
        "Usar secretos o datos empresariales reales en el scaffold.",
        "Enviar mensajes, publicar, gastar, modificar CRM o ejecutar efectos desde esta consola."
      ],
      "status": {
        "specification_status": "SPEC_READY",
        "method_admission": {
          "status": "OFFICIAL_METHOD_REFERENCED",
          "claim_scope": "DOCUMENTARY_METHOD_ONLY",
          "source_refs": [
            "DORA-BATCHES",
            "MS-RBAC"
          ],
          "code_admission": "NONE"
        },
        "verification_status": {
          "scaffold": "NOT_TESTED",
          "target": "NOT_TESTED",
          "evidence_refs": []
        }
      },
      "source_refs": [
        "DORA-BATCHES",
        "MS-RBAC"
      ],
      "galaxy_admission": {
        "status": "WAITING_GALAXY_ADMISSION",
        "source_ref": null,
        "applies_only_to": "FUTURE_EVENT_MATERIAL",
        "content_invented": false,
        "reopen_trigger": "Material oficial identificable del evento disponible y obtenido legalmente; fijar fecha/revisión/hash/términos/claim y ejecutar admisión aplicable antes de incorporarlo."
      },
      "owner_refs": {
        "progress": "specs/library-extension/tasks.md",
        "execution_state": "PROJECT_EXECUTION_STATE.json",
        "execution_events": "PROJECT_EXECUTION_EVENTS.jsonl",
        "failures": "PROJECT_FAILURE_LESSONS.md",
        "dependency": "PROJECT_DEPENDENCY_UPDATE_RECORD.md",
        "freshness": "PROJECT_AUTHORITY_FRESHNESS_RECORD.md"
      },
      "resume": {
        "read_refs": [
          "PROJECT_EXECUTION_STATE.json",
          "PROJECT_EXECUTION_EVENTS.jsonl",
          "specs/library-extension/tasks.md",
          "PROJECT_FAILURE_LESSONS.md",
          "roles/business-functions.v403.json"
        ],
        "rule": "Validar checkpoint y hashes; consultar los owners existentes; invalidar evidencia afectada por delta; registrar siguiente acción allí, sin backlog paralelo."
      },
      "implementation_binding": {
        "status": "EXISTING_LIBRARY_REFERENCE",
        "pack": "implementation_packs/ENGINEERING_EXECUTION_VALIDATOR.md",
        "evidence": "reconstruction_evidence/LIBRARY_INFRA_READY_V402.md",
        "action": "validate_execution_state.py --level resume; validate_project.py --level plan",
        "kind": "CLI_CONTRACT",
        "target_verified": false,
        "v403_executes_binding": false
      },
      "open_requirement": {
        "id": "TARGET-ENGINEERING",
        "scope": "FUTURE_CONSUMER_ONLY",
        "description": "El nuevo proyecto consumidor debe materializar su propio kit, manifest y checkpoint; esta vista no ejecuta gates de producto.",
        "owner_ref": "specs/library-extension/tasks.md"
      },
      "tasks": [
        {
          "id": "engineering.define-change",
          "title": "Delimitar cambio y aceptación",
          "inputs": [
            "Requisito",
            "baseline",
            "restricciones",
            "riesgo"
          ],
          "outputs": [
            "Delta trazable y criterios verificables"
          ],
          "product_permission": {
            "id": "engineering:plan",
            "definition_kind": "PROPOSED_TARGET_REQUIREMENT",
            "granted": false,
            "subject_binding": "REQUIRES_TARGET_CONFIGURATION",
            "organization_scope": "REQUIRED",
            "resource_scope": "REQUIRED",
            "enforcement": "EXISTING_PRODUCT_OWNER_OR_TARGET_GAP",
            "policy_ref": "SECURITY_SRE_CLOUD_INFRASTRUCTURE.md"
          },
          "frontend": {
            "route": "#/functions/engineering",
            "section": "tasks",
            "task_id": "engineering.define-change",
            "interaction": "READ_CONTRACT",
            "effect_enabled": false
          },
          "acceptance": [
            {
              "id": "engineering.AC1",
              "claim_scope": "CONTRACT_SPECIFICATION",
              "criterion": "Cada cambio conserva requisito, owner existente, pruebas afectadas y trigger de rollback.",
              "test_ref": "roles/tests/test_business_functions.py",
              "evidence_refs": [],
              "target_status": "NOT_TESTED"
            }
          ],
          "owner_ref": "specs/library-extension/tasks.md"
        },
        {
          "id": "engineering.verify-change",
          "title": "Verificar el cambio",
          "inputs": [
            "Delta",
            "artefactos y pruebas del alcance"
          ],
          "outputs": [
            "Receipt con revisión, comandos, resultados y límites"
          ],
          "product_permission": {
            "id": "engineering:verify",
            "definition_kind": "PROPOSED_TARGET_REQUIREMENT",
            "granted": false,
            "subject_binding": "REQUIRES_TARGET_CONFIGURATION",
            "organization_scope": "REQUIRED",
            "resource_scope": "REQUIRED",
            "enforcement": "EXISTING_PRODUCT_OWNER_OR_TARGET_GAP",
            "policy_ref": "SECURITY_SRE_CLOUD_INFRASTRUCTURE.md"
          },
          "frontend": {
            "route": "#/functions/engineering",
            "section": "tasks",
            "task_id": "engineering.verify-change",
            "interaction": "READ_CONTRACT",
            "effect_enabled": false
          },
          "acceptance": [
            {
              "id": "engineering.AC2",
              "claim_scope": "CONTRACT_SPECIFICATION",
              "criterion": "Un build no cierra autorización, integración, seguridad ni aceptación humana pendientes.",
              "test_ref": "roles/tests/test_business_functions.py",
              "evidence_refs": [],
              "target_status": "NOT_TESTED"
            }
          ],
          "owner_ref": "specs/library-extension/tasks.md"
        }
      ],
      "owner_binding_scope": "MAINTENANCE_FIXTURE"
    },
    {
      "id": "product-managers",
      "track_name": "Product Managers",
      "function_kind": "BUSINESS_FUNCTION",
      "auth_role": false,
      "purpose": "Relacionar problemas de usuarios con hipótesis, criterios de aceptación y feedback trazable.",
      "systems": [
        "Blueprint/spec/plan existentes",
        "Evidencia de journeys y feedback"
      ],
      "allowed": [
        "Leer contratos, fuentes y evidencia local.",
        "Preparar propuestas sin efectos externos dentro del scope aprobado."
      ],
      "prohibited": [
        "Derivar privilegios del título o track.",
        "Usar secretos o datos empresariales reales en el scaffold.",
        "Enviar mensajes, publicar, gastar, modificar CRM o ejecutar efectos desde esta consola."
      ],
      "status": {
        "specification_status": "SPEC_READY",
        "method_admission": {
          "status": "OFFICIAL_METHOD_REFERENCED",
          "claim_scope": "DOCUMENTARY_METHOD_ONLY",
          "source_refs": [
            "DORA-USER",
            "DORA-BATCHES"
          ],
          "code_admission": "NONE"
        },
        "verification_status": {
          "scaffold": "NOT_TESTED",
          "target": "NOT_TESTED",
          "evidence_refs": []
        }
      },
      "source_refs": [
        "DORA-USER",
        "DORA-BATCHES"
      ],
      "galaxy_admission": {
        "status": "WAITING_GALAXY_ADMISSION",
        "source_ref": null,
        "applies_only_to": "FUTURE_EVENT_MATERIAL",
        "content_invented": false,
        "reopen_trigger": "Material oficial identificable del evento disponible y obtenido legalmente; fijar fecha/revisión/hash/términos/claim y ejecutar admisión aplicable antes de incorporarlo."
      },
      "owner_refs": {
        "progress": "specs/library-extension/tasks.md",
        "execution_state": "PROJECT_EXECUTION_STATE.json",
        "execution_events": "PROJECT_EXECUTION_EVENTS.jsonl",
        "failures": "PROJECT_FAILURE_LESSONS.md",
        "dependency": "PROJECT_DEPENDENCY_UPDATE_RECORD.md",
        "freshness": "PROJECT_AUTHORITY_FRESHNESS_RECORD.md"
      },
      "resume": {
        "read_refs": [
          "PROJECT_EXECUTION_STATE.json",
          "PROJECT_EXECUTION_EVENTS.jsonl",
          "specs/library-extension/tasks.md",
          "PROJECT_FAILURE_LESSONS.md",
          "roles/business-functions.v403.json"
        ],
        "rule": "Validar checkpoint y hashes; consultar los owners existentes; invalidar evidencia afectada por delta; registrar siguiente acción allí, sin backlog paralelo."
      },
      "implementation_binding": {
        "status": "EXISTING_LIBRARY_REFERENCE",
        "pack": "implementation_packs/GO_CUSTOMER_SURVEY_API.md",
        "evidence": "reconstruction_evidence/CONNECTED_CUSTOMER_SURVEYS_V400.md",
        "action": "Referencia a encuesta y feedback del owner existente; no ejecución V403.",
        "kind": "PACK_REFERENCE",
        "target_verified": false,
        "v403_executes_binding": false
      },
      "open_requirement": {
        "id": "TARGET-PRODUCT-MANAGERS",
        "scope": "FUTURE_CONSUMER_ONLY",
        "description": "Validación con usuarios, métricas y decisión de prioridad del producto concreto; no hay motor universal de product management.",
        "owner_ref": "specs/library-extension/tasks.md"
      },
      "tasks": [
        {
          "id": "product-managers.frame-user-problem",
          "title": "Enmarcar problema y resultado",
          "inputs": [
            "Problema observado",
            "usuario",
            "contexto",
            "restricciones"
          ],
          "outputs": [
            "Hipótesis falsable y aceptación del journey"
          ],
          "product_permission": {
            "id": "product:plan",
            "definition_kind": "PROPOSED_TARGET_REQUIREMENT",
            "granted": false,
            "subject_binding": "REQUIRES_TARGET_CONFIGURATION",
            "organization_scope": "REQUIRED",
            "resource_scope": "REQUIRED",
            "enforcement": "EXISTING_PRODUCT_OWNER_OR_TARGET_GAP",
            "policy_ref": "SECURITY_SRE_CLOUD_INFRASTRUCTURE.md"
          },
          "frontend": {
            "route": "#/functions/product-managers",
            "section": "tasks",
            "task_id": "product-managers.frame-user-problem",
            "interaction": "READ_CONTRACT",
            "effect_enabled": false
          },
          "acceptance": [
            {
              "id": "product-managers.AC1",
              "claim_scope": "CONTRACT_SPECIFICATION",
              "criterion": "Separar evidencia observada de supuestos; el owner de negocio decide prioridades y umbrales materiales.",
              "test_ref": "roles/tests/test_business_functions.py",
              "evidence_refs": [],
              "target_status": "NOT_TESTED"
            }
          ],
          "owner_ref": "specs/library-extension/tasks.md"
        },
        {
          "id": "product-managers.review-feedback",
          "title": "Revisar feedback y decisión",
          "inputs": [
            "Evidencia de tarea",
            "versión del producto",
            "limitaciones"
          ],
          "outputs": [
            "Decisión ligada al requisito y a la evidencia"
          ],
          "product_permission": {
            "id": "product:review",
            "definition_kind": "PROPOSED_TARGET_REQUIREMENT",
            "granted": false,
            "subject_binding": "REQUIRES_TARGET_CONFIGURATION",
            "organization_scope": "REQUIRED",
            "resource_scope": "REQUIRED",
            "enforcement": "EXISTING_PRODUCT_OWNER_OR_TARGET_GAP",
            "policy_ref": "SECURITY_SRE_CLOUD_INFRASTRUCTURE.md"
          },
          "frontend": {
            "route": "#/functions/product-managers",
            "section": "tasks",
            "task_id": "product-managers.review-feedback",
            "interaction": "READ_CONTRACT",
            "effect_enabled": false
          },
          "acceptance": [
            {
              "id": "product-managers.AC2",
              "claim_scope": "CONTRACT_SPECIFICATION",
              "criterion": "La métrica propuesta no se declara medida; una encuesta no demuestra causalidad.",
              "test_ref": "roles/tests/test_business_functions.py",
              "evidence_refs": [],
              "target_status": "NOT_TESTED"
            }
          ],
          "owner_ref": "specs/library-extension/tasks.md"
        }
      ],
      "owner_binding_scope": "MAINTENANCE_FIXTURE"
    },
    {
      "id": "founders",
      "track_name": "Founders",
      "function_kind": "BUSINESS_FUNCTION",
      "auth_role": false,
      "purpose": "Explicitar resultado empresarial, restricciones y decisiones que requieren un responsable humano.",
      "systems": [
        "Blueprint y authority map existentes",
        "Intake y decisiones de alcance"
      ],
      "allowed": [
        "Leer contratos, fuentes y evidencia local.",
        "Preparar propuestas sin efectos externos dentro del scope aprobado."
      ],
      "prohibited": [
        "Derivar privilegios del título o track.",
        "Usar secretos o datos empresariales reales en el scaffold.",
        "Enviar mensajes, publicar, gastar, modificar CRM o ejecutar efectos desde esta consola.",
        "Suponer privilegios globales o aprobación de pagos por ser Founder."
      ],
      "status": {
        "specification_status": "SPEC_READY",
        "method_admission": {
          "status": "OFFICIAL_METHOD_REFERENCED",
          "claim_scope": "DOCUMENTARY_METHOD_ONLY",
          "source_refs": [
            "DORA-USER",
            "DORA-BATCHES"
          ],
          "code_admission": "NONE"
        },
        "verification_status": {
          "scaffold": "NOT_TESTED",
          "target": "NOT_TESTED",
          "evidence_refs": []
        }
      },
      "source_refs": [
        "DORA-USER",
        "DORA-BATCHES"
      ],
      "galaxy_admission": {
        "status": "WAITING_GALAXY_ADMISSION",
        "source_ref": null,
        "applies_only_to": "FUTURE_EVENT_MATERIAL",
        "content_invented": false,
        "reopen_trigger": "Material oficial identificable del evento disponible y obtenido legalmente; fijar fecha/revisión/hash/términos/claim y ejecutar admisión aplicable antes de incorporarlo."
      },
      "owner_refs": {
        "progress": "specs/library-extension/tasks.md",
        "execution_state": "PROJECT_EXECUTION_STATE.json",
        "execution_events": "PROJECT_EXECUTION_EVENTS.jsonl",
        "failures": "PROJECT_FAILURE_LESSONS.md",
        "dependency": "PROJECT_DEPENDENCY_UPDATE_RECORD.md",
        "freshness": "PROJECT_AUTHORITY_FRESHNESS_RECORD.md"
      },
      "resume": {
        "read_refs": [
          "PROJECT_EXECUTION_STATE.json",
          "PROJECT_EXECUTION_EVENTS.jsonl",
          "specs/library-extension/tasks.md",
          "PROJECT_FAILURE_LESSONS.md",
          "roles/business-functions.v403.json"
        ],
        "rule": "Validar checkpoint y hashes; consultar los owners existentes; invalidar evidencia afectada por delta; registrar siguiente acción allí, sin backlog paralelo."
      },
      "implementation_binding": {
        "status": "DOCUMENTED_ONLY",
        "pack": null,
        "evidence": null,
        "action": null,
        "kind": "NO_PRODUCT_RUNTIME",
        "target_verified": false,
        "v403_executes_binding": false
      },
      "open_requirement": {
        "id": "TARGET-FOUNDERS",
        "scope": "FUTURE_CONSUMER_ONLY",
        "description": "Decisiones de modelo de negocio, gasto y regulación pertenecen al proyecto consumidor y su owner real.",
        "owner_ref": "specs/library-extension/tasks.md"
      },
      "tasks": [
        {
          "id": "founders.record-business-hypothesis",
          "title": "Registrar hipótesis empresarial",
          "inputs": [
            "Problema",
            "segmento propuesto",
            "evidencia",
            "restricciones"
          ],
          "outputs": [
            "Hipótesis con incertidumbre y criterio de contraste"
          ],
          "product_permission": {
            "id": "business:plan",
            "definition_kind": "PROPOSED_TARGET_REQUIREMENT",
            "granted": false,
            "subject_binding": "REQUIRES_TARGET_CONFIGURATION",
            "organization_scope": "REQUIRED",
            "resource_scope": "REQUIRED",
            "enforcement": "EXISTING_PRODUCT_OWNER_OR_TARGET_GAP",
            "policy_ref": "SECURITY_SRE_CLOUD_INFRASTRUCTURE.md"
          },
          "frontend": {
            "route": "#/functions/founders",
            "section": "tasks",
            "task_id": "founders.record-business-hypothesis",
            "interaction": "READ_CONTRACT",
            "effect_enabled": false
          },
          "acceptance": [
            {
              "id": "founders.AC1",
              "claim_scope": "CONTRACT_SPECIFICATION",
              "criterion": "No inventar precios, obligaciones regulatorias, inversión ni autorización de gasto.",
              "test_ref": "roles/tests/test_business_functions.py",
              "evidence_refs": [],
              "target_status": "NOT_TESTED"
            }
          ],
          "owner_ref": "specs/library-extension/tasks.md"
        },
        {
          "id": "founders.review-scope-decision",
          "title": "Revisar decisión de alcance",
          "inputs": [
            "Opciones",
            "costos conocidos",
            "riesgos",
            "evidencia"
          ],
          "outputs": [
            "Decisión o pregunta material para owner real"
          ],
          "product_permission": {
            "id": "business:review",
            "definition_kind": "PROPOSED_TARGET_REQUIREMENT",
            "granted": false,
            "subject_binding": "REQUIRES_TARGET_CONFIGURATION",
            "organization_scope": "REQUIRED",
            "resource_scope": "REQUIRED",
            "enforcement": "EXISTING_PRODUCT_OWNER_OR_TARGET_GAP",
            "policy_ref": "SECURITY_SRE_CLOUD_INFRASTRUCTURE.md"
          },
          "frontend": {
            "route": "#/functions/founders",
            "section": "tasks",
            "task_id": "founders.review-scope-decision",
            "interaction": "READ_CONTRACT",
            "effect_enabled": false
          },
          "acceptance": [
            {
              "id": "founders.AC2",
              "claim_scope": "CONTRACT_SPECIFICATION",
              "criterion": "El título Founders no otorga admin, aprobación de pagos ni privilegios entre organizaciones.",
              "test_ref": "roles/tests/test_business_functions.py",
              "evidence_refs": [],
              "target_status": "NOT_TESTED"
            }
          ],
          "owner_ref": "specs/library-extension/tasks.md"
        }
      ],
      "owner_binding_scope": "MAINTENANCE_FIXTURE"
    },
    {
      "id": "sales-engineering",
      "track_name": "Sales Engineering",
      "function_kind": "BUSINESS_FUNCTION",
      "auth_role": false,
      "purpose": "Relacionar una necesidad comercial con capacidades técnicas demostradas y gaps explícitos.",
      "systems": [
        "Packs y evidencia local",
        "Contratos de integración y demo"
      ],
      "allowed": [
        "Leer contratos, fuentes y evidencia local.",
        "Preparar propuestas sin efectos externos dentro del scope aprobado."
      ],
      "prohibited": [
        "Derivar privilegios del título o track.",
        "Usar secretos o datos empresariales reales en el scaffold.",
        "Enviar mensajes, publicar, gastar, modificar CRM o ejecutar efectos desde esta consola."
      ],
      "status": {
        "specification_status": "SPEC_READY",
        "method_admission": {
          "status": "OFFICIAL_METHOD_REFERENCED",
          "claim_scope": "DOCUMENTARY_METHOD_ONLY",
          "source_refs": [
            "MS-SALES",
            "DORA-USER"
          ],
          "code_admission": "NONE"
        },
        "verification_status": {
          "scaffold": "NOT_TESTED",
          "target": "NOT_TESTED",
          "evidence_refs": []
        }
      },
      "source_refs": [
        "MS-SALES",
        "DORA-USER"
      ],
      "galaxy_admission": {
        "status": "WAITING_GALAXY_ADMISSION",
        "source_ref": null,
        "applies_only_to": "FUTURE_EVENT_MATERIAL",
        "content_invented": false,
        "reopen_trigger": "Material oficial identificable del evento disponible y obtenido legalmente; fijar fecha/revisión/hash/términos/claim y ejecutar admisión aplicable antes de incorporarlo."
      },
      "owner_refs": {
        "progress": "specs/library-extension/tasks.md",
        "execution_state": "PROJECT_EXECUTION_STATE.json",
        "execution_events": "PROJECT_EXECUTION_EVENTS.jsonl",
        "failures": "PROJECT_FAILURE_LESSONS.md",
        "dependency": "PROJECT_DEPENDENCY_UPDATE_RECORD.md",
        "freshness": "PROJECT_AUTHORITY_FRESHNESS_RECORD.md"
      },
      "resume": {
        "read_refs": [
          "PROJECT_EXECUTION_STATE.json",
          "PROJECT_EXECUTION_EVENTS.jsonl",
          "specs/library-extension/tasks.md",
          "PROJECT_FAILURE_LESSONS.md",
          "roles/business-functions.v403.json"
        ],
        "rule": "Validar checkpoint y hashes; consultar los owners existentes; invalidar evidencia afectada por delta; registrar siguiente acción allí, sin backlog paralelo."
      },
      "implementation_binding": {
        "status": "EXISTING_LIBRARY_REFERENCE",
        "pack": "implementation_packs/GO_FRANCHISE_CUSTOMER_JOURNEY_API.md",
        "evidence": "reconstruction_evidence/LIBRARY_SIGNED_REFERENCE_ACCEPTANCE_V402.md",
        "action": "Revisión de contratos del journey de referencia y su aceptación local firmada.",
        "kind": "PACK_REFERENCE",
        "target_verified": false,
        "v403_executes_binding": false
      },
      "open_requirement": {
        "id": "TARGET-SALES-ENGINEERING",
        "scope": "FUTURE_CONSUMER_ONLY",
        "description": "Mapeo técnico de un cliente concreto y validación de sus integraciones; no automatizar compromisos comerciales.",
        "owner_ref": "specs/library-extension/tasks.md"
      },
      "tasks": [
        {
          "id": "sales-engineering.map-technical-fit",
          "title": "Evaluar encaje técnico",
          "inputs": [
            "Necesidad del prospecto",
            "restricciones",
            "evidencia de capability"
          ],
          "outputs": [
            "Matriz requisito-capacidad-gap"
          ],
          "product_permission": {
            "id": "sales-engineering:assess",
            "definition_kind": "PROPOSED_TARGET_REQUIREMENT",
            "granted": false,
            "subject_binding": "REQUIRES_TARGET_CONFIGURATION",
            "organization_scope": "REQUIRED",
            "resource_scope": "REQUIRED",
            "enforcement": "EXISTING_PRODUCT_OWNER_OR_TARGET_GAP",
            "policy_ref": "SECURITY_SRE_CLOUD_INFRASTRUCTURE.md"
          },
          "frontend": {
            "route": "#/functions/sales-engineering",
            "section": "tasks",
            "task_id": "sales-engineering.map-technical-fit",
            "interaction": "READ_CONTRACT",
            "effect_enabled": false
          },
          "acceptance": [
            {
              "id": "sales-engineering.AC1",
              "claim_scope": "CONTRACT_SPECIFICATION",
              "criterion": "Cada capacidad afirmada enlaza una prueba con alcance; las capacidades futuras quedan pendientes.",
              "test_ref": "roles/tests/test_business_functions.py",
              "evidence_refs": [],
              "target_status": "NOT_TESTED"
            }
          ],
          "owner_ref": "specs/library-extension/tasks.md"
        },
        {
          "id": "sales-engineering.prepare-scoped-demo",
          "title": "Preparar una demostración acotada",
          "inputs": [
            "Matriz de encaje",
            "fixture permitido",
            "recorrido y límites"
          ],
          "outputs": [
            "Guion de demo reproducible y resultados esperados"
          ],
          "product_permission": {
            "id": "sales-engineering:demo",
            "definition_kind": "PROPOSED_TARGET_REQUIREMENT",
            "granted": false,
            "subject_binding": "REQUIRES_TARGET_CONFIGURATION",
            "organization_scope": "REQUIRED",
            "resource_scope": "REQUIRED",
            "enforcement": "EXISTING_PRODUCT_OWNER_OR_TARGET_GAP",
            "policy_ref": "SECURITY_SRE_CLOUD_INFRASTRUCTURE.md"
          },
          "frontend": {
            "route": "#/functions/sales-engineering",
            "section": "tasks",
            "task_id": "sales-engineering.prepare-scoped-demo",
            "interaction": "READ_CONTRACT",
            "effect_enabled": false
          },
          "acceptance": [
            {
              "id": "sales-engineering.AC2",
              "claim_scope": "CONTRACT_SPECIFICATION",
              "criterion": "Un fixture no se presenta como integración live, homologación, SLA ni aceptación del cliente.",
              "test_ref": "roles/tests/test_business_functions.py",
              "evidence_refs": [],
              "target_status": "NOT_TESTED"
            }
          ],
          "owner_ref": "specs/library-extension/tasks.md"
        }
      ],
      "owner_binding_scope": "MAINTENANCE_FIXTURE"
    },
    {
      "id": "sales",
      "track_name": "Sales",
      "function_kind": "BUSINESS_FUNCTION",
      "auth_role": false,
      "purpose": "Preparar una oportunidad y revisar una propuesta bajo reglas comerciales aprobadas.",
      "systems": [
        "CRM y cotización existentes",
        "Aprobaciones y dominio transaccional"
      ],
      "allowed": [
        "Leer contratos, fuentes y evidencia local.",
        "Preparar propuestas sin efectos externos dentro del scope aprobado."
      ],
      "prohibited": [
        "Derivar privilegios del título o track.",
        "Usar secretos o datos empresariales reales en el scaffold.",
        "Enviar mensajes, publicar, gastar, modificar CRM o ejecutar efectos desde esta consola."
      ],
      "status": {
        "specification_status": "SPEC_READY",
        "method_admission": {
          "status": "OFFICIAL_METHOD_REFERENCED",
          "claim_scope": "DOCUMENTARY_METHOD_ONLY",
          "source_refs": [
            "MS-SALES",
            "MS-RBAC"
          ],
          "code_admission": "NONE"
        },
        "verification_status": {
          "scaffold": "NOT_TESTED",
          "target": "NOT_TESTED",
          "evidence_refs": []
        }
      },
      "source_refs": [
        "MS-SALES",
        "MS-RBAC"
      ],
      "galaxy_admission": {
        "status": "WAITING_GALAXY_ADMISSION",
        "source_ref": null,
        "applies_only_to": "FUTURE_EVENT_MATERIAL",
        "content_invented": false,
        "reopen_trigger": "Material oficial identificable del evento disponible y obtenido legalmente; fijar fecha/revisión/hash/términos/claim y ejecutar admisión aplicable antes de incorporarlo."
      },
      "owner_refs": {
        "progress": "specs/library-extension/tasks.md",
        "execution_state": "PROJECT_EXECUTION_STATE.json",
        "execution_events": "PROJECT_EXECUTION_EVENTS.jsonl",
        "failures": "PROJECT_FAILURE_LESSONS.md",
        "dependency": "PROJECT_DEPENDENCY_UPDATE_RECORD.md",
        "freshness": "PROJECT_AUTHORITY_FRESHNESS_RECORD.md"
      },
      "resume": {
        "read_refs": [
          "PROJECT_EXECUTION_STATE.json",
          "PROJECT_EXECUTION_EVENTS.jsonl",
          "specs/library-extension/tasks.md",
          "PROJECT_FAILURE_LESSONS.md",
          "roles/business-functions.v403.json"
        ],
        "rule": "Validar checkpoint y hashes; consultar los owners existentes; invalidar evidencia afectada por delta; registrar siguiente acción allí, sin backlog paralelo."
      },
      "implementation_binding": {
        "status": "EXISTING_LIBRARY_REFERENCE",
        "pack": "implementation_packs/GO_FRANCHISE_CUSTOMER_JOURNEY_API.md",
        "evidence": "reconstruction_evidence/LIBRARY_SIGNED_REFERENCE_ACCEPTANCE_V402.md",
        "action": "Consultar el contrato existente de cotización/pedido y sus límites.",
        "kind": "PACK_REFERENCE",
        "target_verified": false,
        "v403_executes_binding": false
      },
      "open_requirement": {
        "id": "TARGET-SALES",
        "scope": "FUTURE_CONSUMER_ONLY",
        "description": "Permisos efectivos y reglas comerciales del target; los identificadores sales:* son requisitos propuestos, no una asignación existente.",
        "owner_ref": "specs/library-extension/tasks.md"
      },
      "tasks": [
        {
          "id": "sales.review-opportunity",
          "title": "Revisar oportunidad",
          "inputs": [
            "Necesidad y evidencia de lead",
            "identidad y organización",
            "política aprobada"
          ],
          "outputs": [
            "Oportunidad con trazabilidad y siguiente paso propuesto"
          ],
          "product_permission": {
            "id": "sales:read",
            "definition_kind": "PROPOSED_TARGET_REQUIREMENT",
            "granted": false,
            "subject_binding": "REQUIRES_TARGET_CONFIGURATION",
            "organization_scope": "REQUIRED",
            "resource_scope": "REQUIRED",
            "enforcement": "EXISTING_PRODUCT_OWNER_OR_TARGET_GAP",
            "policy_ref": "SECURITY_SRE_CLOUD_INFRASTRUCTURE.md"
          },
          "frontend": {
            "route": "#/functions/sales",
            "section": "tasks",
            "task_id": "sales.review-opportunity",
            "interaction": "READ_CONTRACT",
            "effect_enabled": false
          },
          "acceptance": [
            {
              "id": "sales.AC1",
              "claim_scope": "CONTRACT_SPECIFICATION",
              "criterion": "La identidad procede del owner existente; no inferir contacto ni organización del texto del agente.",
              "test_ref": "roles/tests/test_business_functions.py",
              "evidence_refs": [],
              "target_status": "NOT_TESTED"
            }
          ],
          "owner_ref": "specs/library-extension/tasks.md"
        },
        {
          "id": "sales.prepare-proposal",
          "title": "Preparar propuesta",
          "inputs": [
            "Oportunidad",
            "catálogo/versiones",
            "política comercial y approval refs"
          ],
          "outputs": [
            "Propuesta vinculada a su política y aceptación requerida"
          ],
          "product_permission": {
            "id": "sales:request",
            "definition_kind": "PROPOSED_TARGET_REQUIREMENT",
            "granted": false,
            "subject_binding": "REQUIRES_TARGET_CONFIGURATION",
            "organization_scope": "REQUIRED",
            "resource_scope": "REQUIRED",
            "enforcement": "EXISTING_PRODUCT_OWNER_OR_TARGET_GAP",
            "policy_ref": "SECURITY_SRE_CLOUD_INFRASTRUCTURE.md"
          },
          "frontend": {
            "route": "#/functions/sales",
            "section": "tasks",
            "task_id": "sales.prepare-proposal",
            "interaction": "READ_CONTRACT",
            "effect_enabled": false
          },
          "acceptance": [
            {
              "id": "sales.AC2",
              "claim_scope": "CONTRACT_SPECIFICATION",
              "criterion": "Sin política o autoridad, no inventar descuento, precio, disponibilidad ni compromiso contractual.",
              "test_ref": "roles/tests/test_business_functions.py",
              "evidence_refs": [],
              "target_status": "NOT_TESTED"
            }
          ],
          "owner_ref": "specs/library-extension/tasks.md"
        }
      ],
      "owner_binding_scope": "MAINTENANCE_FIXTURE"
    },
    {
      "id": "sdrs",
      "track_name": "SDRs",
      "function_kind": "BUSINESS_FUNCTION",
      "auth_role": false,
      "purpose": "Investigar y preparar el traspaso de un lead con procedencia y límites de contacto.",
      "systems": [
        "Ingress y promoción de leads existentes",
        "Identidad/contacto y aprobación"
      ],
      "allowed": [
        "Leer contratos, fuentes y evidencia local.",
        "Preparar propuestas sin efectos externos dentro del scope aprobado."
      ],
      "prohibited": [
        "Derivar privilegios del título o track.",
        "Usar secretos o datos empresariales reales en el scaffold.",
        "Enviar mensajes, publicar, gastar, modificar CRM o ejecutar efectos desde esta consola."
      ],
      "status": {
        "specification_status": "SPEC_READY",
        "method_admission": {
          "status": "OFFICIAL_METHOD_REFERENCED",
          "claim_scope": "DOCUMENTARY_METHOD_ONLY",
          "source_refs": [
            "MS-SALES",
            "MS-RBAC"
          ],
          "code_admission": "NONE"
        },
        "verification_status": {
          "scaffold": "NOT_TESTED",
          "target": "NOT_TESTED",
          "evidence_refs": []
        }
      },
      "source_refs": [
        "MS-SALES",
        "MS-RBAC"
      ],
      "galaxy_admission": {
        "status": "WAITING_GALAXY_ADMISSION",
        "source_ref": null,
        "applies_only_to": "FUTURE_EVENT_MATERIAL",
        "content_invented": false,
        "reopen_trigger": "Material oficial identificable del evento disponible y obtenido legalmente; fijar fecha/revisión/hash/términos/claim y ejecutar admisión aplicable antes de incorporarlo."
      },
      "owner_refs": {
        "progress": "specs/library-extension/tasks.md",
        "execution_state": "PROJECT_EXECUTION_STATE.json",
        "execution_events": "PROJECT_EXECUTION_EVENTS.jsonl",
        "failures": "PROJECT_FAILURE_LESSONS.md",
        "dependency": "PROJECT_DEPENDENCY_UPDATE_RECORD.md",
        "freshness": "PROJECT_AUTHORITY_FRESHNESS_RECORD.md"
      },
      "resume": {
        "read_refs": [
          "PROJECT_EXECUTION_STATE.json",
          "PROJECT_EXECUTION_EVENTS.jsonl",
          "specs/library-extension/tasks.md",
          "PROJECT_FAILURE_LESSONS.md",
          "roles/business-functions.v403.json"
        ],
        "rule": "Validar checkpoint y hashes; consultar los owners existentes; invalidar evidencia afectada por delta; registrar siguiente acción allí, sin backlog paralelo."
      },
      "implementation_binding": {
        "status": "EXISTING_LIBRARY_REFERENCE",
        "pack": "implementation_packs/GO_LEAD_CANDIDATE_PROMOTION.md",
        "evidence": "reconstruction_evidence/LIBRARY_INFRA_READY_V402.md",
        "action": "Consultar promoción gobernada de leads del pack existente.",
        "kind": "PACK_REFERENCE",
        "target_verified": false,
        "v403_executes_binding": false
      },
      "open_requirement": {
        "id": "TARGET-SDRS",
        "scope": "FUTURE_CONSUMER_ONLY",
        "description": "Admitir fuente/proveedor y criterios de calificación del target; vínculo de pack no prueba un CRM universal.",
        "owner_ref": "specs/library-extension/tasks.md"
      },
      "tasks": [
        {
          "id": "sdrs.review-lead-evidence",
          "title": "Revisar procedencia del lead",
          "inputs": [
            "Referencia de lead",
            "fuente",
            "consentimiento y propósito aplicables"
          ],
          "outputs": [
            "Ficha de procedencia y restricciones de uso"
          ],
          "product_permission": {
            "id": "lead:read",
            "definition_kind": "PROPOSED_TARGET_REQUIREMENT",
            "granted": false,
            "subject_binding": "REQUIRES_TARGET_CONFIGURATION",
            "organization_scope": "REQUIRED",
            "resource_scope": "REQUIRED",
            "enforcement": "EXISTING_PRODUCT_OWNER_OR_TARGET_GAP",
            "policy_ref": "SECURITY_SRE_CLOUD_INFRASTRUCTURE.md"
          },
          "frontend": {
            "route": "#/functions/sdrs",
            "section": "tasks",
            "task_id": "sdrs.review-lead-evidence",
            "interaction": "READ_CONTRACT",
            "effect_enabled": false
          },
          "acceptance": [
            {
              "id": "sdrs.AC1",
              "claim_scope": "CONTRACT_SPECIFICATION",
              "criterion": "La existencia de un dato público no demuestra consentimiento, exactitud ni permiso de contacto.",
              "test_ref": "roles/tests/test_business_functions.py",
              "evidence_refs": [],
              "target_status": "NOT_TESTED"
            }
          ],
          "owner_ref": "specs/library-extension/tasks.md"
        },
        {
          "id": "sdrs.prepare-qualified-handoff",
          "title": "Preparar handoff comercial",
          "inputs": [
            "Ficha",
            "criterios de calificación aprobados",
            "owner receptor"
          ],
          "outputs": [
            "Handoff propuesto con evidencia y gaps"
          ],
          "product_permission": {
            "id": "lead:request",
            "definition_kind": "PROPOSED_TARGET_REQUIREMENT",
            "granted": false,
            "subject_binding": "REQUIRES_TARGET_CONFIGURATION",
            "organization_scope": "REQUIRED",
            "resource_scope": "REQUIRED",
            "enforcement": "EXISTING_PRODUCT_OWNER_OR_TARGET_GAP",
            "policy_ref": "SECURITY_SRE_CLOUD_INFRASTRUCTURE.md"
          },
          "frontend": {
            "route": "#/functions/sdrs",
            "section": "tasks",
            "task_id": "sdrs.prepare-qualified-handoff",
            "interaction": "READ_CONTRACT",
            "effect_enabled": false
          },
          "acceptance": [
            {
              "id": "sdrs.AC2",
              "claim_scope": "CONTRACT_SPECIFICATION",
              "criterion": "No inventar scoring o criterios; no enviar mensajes ni promover automáticamente por nombre del track.",
              "test_ref": "roles/tests/test_business_functions.py",
              "evidence_refs": [],
              "target_status": "NOT_TESTED"
            }
          ],
          "owner_ref": "specs/library-extension/tasks.md"
        }
      ],
      "owner_binding_scope": "MAINTENANCE_FIXTURE"
    },
    {
      "id": "customer-support",
      "track_name": "Customer Support",
      "function_kind": "BUSINESS_FUNCTION",
      "auth_role": false,
      "purpose": "Clasificar una consulta y preparar una resolución o escalación comprobable.",
      "systems": [
        "Runtime conversacional existente",
        "Owner de atención/handoff y dominio"
      ],
      "allowed": [
        "Leer contratos, fuentes y evidencia local.",
        "Preparar propuestas sin efectos externos dentro del scope aprobado."
      ],
      "prohibited": [
        "Derivar privilegios del título o track.",
        "Usar secretos o datos empresariales reales en el scaffold.",
        "Enviar mensajes, publicar, gastar, modificar CRM o ejecutar efectos desde esta consola."
      ],
      "status": {
        "specification_status": "SPEC_READY",
        "method_admission": {
          "status": "OFFICIAL_METHOD_REFERENCED",
          "claim_scope": "DOCUMENTARY_METHOD_ONLY",
          "source_refs": [
            "MS-CASES",
            "MS-RBAC"
          ],
          "code_admission": "NONE"
        },
        "verification_status": {
          "scaffold": "NOT_TESTED",
          "target": "NOT_TESTED",
          "evidence_refs": []
        }
      },
      "source_refs": [
        "MS-CASES",
        "MS-RBAC"
      ],
      "galaxy_admission": {
        "status": "WAITING_GALAXY_ADMISSION",
        "source_ref": null,
        "applies_only_to": "FUTURE_EVENT_MATERIAL",
        "content_invented": false,
        "reopen_trigger": "Material oficial identificable del evento disponible y obtenido legalmente; fijar fecha/revisión/hash/términos/claim y ejecutar admisión aplicable antes de incorporarlo."
      },
      "owner_refs": {
        "progress": "specs/library-extension/tasks.md",
        "execution_state": "PROJECT_EXECUTION_STATE.json",
        "execution_events": "PROJECT_EXECUTION_EVENTS.jsonl",
        "failures": "PROJECT_FAILURE_LESSONS.md",
        "dependency": "PROJECT_DEPENDENCY_UPDATE_RECORD.md",
        "freshness": "PROJECT_AUTHORITY_FRESHNESS_RECORD.md"
      },
      "resume": {
        "read_refs": [
          "PROJECT_EXECUTION_STATE.json",
          "PROJECT_EXECUTION_EVENTS.jsonl",
          "specs/library-extension/tasks.md",
          "PROJECT_FAILURE_LESSONS.md",
          "roles/business-functions.v403.json"
        ],
        "rule": "Validar checkpoint y hashes; consultar los owners existentes; invalidar evidencia afectada por delta; registrar siguiente acción allí, sin backlog paralelo."
      },
      "implementation_binding": {
        "status": "EXISTING_LIBRARY_REFERENCE",
        "pack": "implementation_packs/GO_CONNECTED_CONVERSATION_RUNTIME.md",
        "evidence": "reconstruction_evidence/AI_CONNECTED_REFERENCE_RELEASE_V402.md",
        "action": "Run/Handoff del runtime existente; autorización, presupuesto y replay mantienen sus owners.",
        "kind": "PACK_REFERENCE",
        "target_verified": false,
        "v403_executes_binding": false
      },
      "open_requirement": {
        "id": "TARGET-CUSTOMER-SUPPORT",
        "scope": "FUTURE_CONSUMER_ONLY",
        "description": "Políticas/SLA y prueba del canal del consumidor; no representa ticketing universal ni aceptación del cliente.",
        "owner_ref": "specs/library-extension/tasks.md"
      },
      "tasks": [
        {
          "id": "customer-support.triage-case",
          "title": "Clasificar consulta",
          "inputs": [
            "Consulta permitida",
            "identidad scoped",
            "historial autorizado",
            "política de soporte"
          ],
          "outputs": [
            "Caso con categoría propuesta y owner de escalación"
          ],
          "product_permission": {
            "id": "support:read",
            "definition_kind": "PROPOSED_TARGET_REQUIREMENT",
            "granted": false,
            "subject_binding": "REQUIRES_TARGET_CONFIGURATION",
            "organization_scope": "REQUIRED",
            "resource_scope": "REQUIRED",
            "enforcement": "EXISTING_PRODUCT_OWNER_OR_TARGET_GAP",
            "policy_ref": "SECURITY_SRE_CLOUD_INFRASTRUCTURE.md"
          },
          "frontend": {
            "route": "#/functions/customer-support",
            "section": "tasks",
            "task_id": "customer-support.triage-case",
            "interaction": "READ_CONTRACT",
            "effect_enabled": false
          },
          "acceptance": [
            {
              "id": "customer-support.AC1",
              "claim_scope": "CONTRACT_SPECIFICATION",
              "criterion": "Sin identidad/contacto y acceso al objeto demostrados, no revelar historial ni datos de otras organizaciones.",
              "test_ref": "roles/tests/test_business_functions.py",
              "evidence_refs": [],
              "target_status": "NOT_TESTED"
            }
          ],
          "owner_ref": "specs/library-extension/tasks.md"
        },
        {
          "id": "customer-support.prepare-resolution",
          "title": "Preparar respuesta o escalación",
          "inputs": [
            "Caso",
            "conocimiento vigente",
            "política",
            "evidencia de resultado"
          ],
          "outputs": [
            "Respuesta propuesta o handoff con estado explícito"
          ],
          "product_permission": {
            "id": "support:request",
            "definition_kind": "PROPOSED_TARGET_REQUIREMENT",
            "granted": false,
            "subject_binding": "REQUIRES_TARGET_CONFIGURATION",
            "organization_scope": "REQUIRED",
            "resource_scope": "REQUIRED",
            "enforcement": "EXISTING_PRODUCT_OWNER_OR_TARGET_GAP",
            "policy_ref": "SECURITY_SRE_CLOUD_INFRASTRUCTURE.md"
          },
          "frontend": {
            "route": "#/functions/customer-support",
            "section": "tasks",
            "task_id": "customer-support.prepare-resolution",
            "interaction": "READ_CONTRACT",
            "effect_enabled": false
          },
          "acceptance": [
            {
              "id": "customer-support.AC2",
              "claim_scope": "CONTRACT_SPECIFICATION",
              "criterion": "Preparado/enviado/confirmado son estados distintos; no cerrar por texto generado ni prometer reembolso.",
              "test_ref": "roles/tests/test_business_functions.py",
              "evidence_refs": [],
              "target_status": "NOT_TESTED"
            }
          ],
          "owner_ref": "specs/library-extension/tasks.md"
        }
      ],
      "owner_binding_scope": "MAINTENANCE_FIXTURE"
    },
    {
      "id": "marketing-operations",
      "track_name": "Marketing Operations",
      "function_kind": "BUSINESS_FUNCTION",
      "auth_role": false,
      "purpose": "Preparar y revisar campañas con audiencia, consentimiento, aprobaciones y recuperación explícitos.",
      "systems": [
        "Campañas/agenda WhatsApp existentes",
        "Aprobaciones/jobs/fence outbound existentes"
      ],
      "allowed": [
        "Leer contratos, fuentes y evidencia local.",
        "Preparar propuestas sin efectos externos dentro del scope aprobado."
      ],
      "prohibited": [
        "Derivar privilegios del título o track.",
        "Usar secretos o datos empresariales reales en el scaffold.",
        "Enviar mensajes, publicar, gastar, modificar CRM o ejecutar efectos desde esta consola."
      ],
      "status": {
        "specification_status": "SPEC_READY",
        "method_admission": {
          "status": "OFFICIAL_METHOD_REFERENCED",
          "claim_scope": "DOCUMENTARY_METHOD_ONLY",
          "source_refs": [
            "HUBSPOT-CAMPAIGNS",
            "MS-RBAC"
          ],
          "code_admission": "NONE"
        },
        "verification_status": {
          "scaffold": "NOT_TESTED",
          "target": "NOT_TESTED",
          "evidence_refs": []
        }
      },
      "source_refs": [
        "HUBSPOT-CAMPAIGNS",
        "MS-RBAC"
      ],
      "galaxy_admission": {
        "status": "WAITING_GALAXY_ADMISSION",
        "source_ref": null,
        "applies_only_to": "FUTURE_EVENT_MATERIAL",
        "content_invented": false,
        "reopen_trigger": "Material oficial identificable del evento disponible y obtenido legalmente; fijar fecha/revisión/hash/términos/claim y ejecutar admisión aplicable antes de incorporarlo."
      },
      "owner_refs": {
        "progress": "specs/library-extension/tasks.md",
        "execution_state": "PROJECT_EXECUTION_STATE.json",
        "execution_events": "PROJECT_EXECUTION_EVENTS.jsonl",
        "failures": "PROJECT_FAILURE_LESSONS.md",
        "dependency": "PROJECT_DEPENDENCY_UPDATE_RECORD.md",
        "freshness": "PROJECT_AUTHORITY_FRESHNESS_RECORD.md"
      },
      "resume": {
        "read_refs": [
          "PROJECT_EXECUTION_STATE.json",
          "PROJECT_EXECUTION_EVENTS.jsonl",
          "specs/library-extension/tasks.md",
          "PROJECT_FAILURE_LESSONS.md",
          "roles/business-functions.v403.json"
        ],
        "rule": "Validar checkpoint y hashes; consultar los owners existentes; invalidar evidencia afectada por delta; registrar siguiente acción allí, sin backlog paralelo."
      },
      "implementation_binding": {
        "status": "EXISTING_LIBRARY_REFERENCE",
        "pack": "implementation_packs/GO_CONNECTED_WHATSAPP_CAMPAIGNS.md",
        "evidence": "reconstruction_evidence/CAMPAIGN_CONNECTED_RELEASE_V402.md",
        "action": "POST /v1/franchise/marketing/campaigns/prepare; GET /v1/franchise/marketing/campaigns/{id}",
        "kind": "OBSERVED_API_ROUTE",
        "target_verified": false,
        "v403_executes_binding": false
      },
      "open_requirement": {
        "id": "TARGET-MARKETING-OPERATIONS",
        "scope": "FUTURE_CONSUMER_ONLY",
        "description": "Provider/account/consent/template/costo del proyecto; la consola V403 sólo muestra el contrato y no llama esas rutas.",
        "owner_ref": "specs/library-extension/tasks.md"
      },
      "tasks": [
        {
          "id": "marketing-operations.prepare-campaign-operations",
          "title": "Preparar operación de campaña",
          "inputs": [
            "Audiencia permitida",
            "propósito",
            "template exacto",
            "ventanas y presupuesto aprobados"
          ],
          "outputs": [
            "Plan trazable de campaña y sus revisiones"
          ],
          "product_permission": {
            "id": "marketing:request",
            "definition_kind": "OBSERVED_EXISTING_OWNER_NAME",
            "granted": false,
            "subject_binding": "REQUIRES_TARGET_CONFIGURATION",
            "organization_scope": "REQUIRED",
            "resource_scope": "REQUIRED",
            "enforcement": "EXISTING_PRODUCT_OWNER_OR_TARGET_GAP",
            "policy_ref": "SECURITY_SRE_CLOUD_INFRASTRUCTURE.md"
          },
          "frontend": {
            "route": "#/functions/marketing-operations",
            "section": "tasks",
            "task_id": "marketing-operations.prepare-campaign-operations",
            "interaction": "READ_CONTRACT",
            "effect_enabled": false
          },
          "acceptance": [
            {
              "id": "marketing-operations.AC1",
              "claim_scope": "CONTRACT_SPECIFICATION",
              "criterion": "Conservar fuente/consentimiento; separar solicitud y aprobación; no activar por el track.",
              "test_ref": "roles/tests/test_business_functions.py",
              "evidence_refs": [],
              "target_status": "NOT_TESTED"
            }
          ],
          "owner_ref": "specs/library-extension/tasks.md"
        },
        {
          "id": "marketing-operations.review-campaign-state",
          "title": "Revisar estado y recuperación",
          "inputs": [
            "ID y revisión de campaña",
            "recibos",
            "pasos confirmados/inciertos"
          ],
          "outputs": [
            "Estado observado y siguiente paso seguro"
          ],
          "product_permission": {
            "id": "marketing:read",
            "definition_kind": "OBSERVED_EXISTING_OWNER_NAME",
            "granted": false,
            "subject_binding": "REQUIRES_TARGET_CONFIGURATION",
            "organization_scope": "REQUIRED",
            "resource_scope": "REQUIRED",
            "enforcement": "EXISTING_PRODUCT_OWNER_OR_TARGET_GAP",
            "policy_ref": "SECURITY_SRE_CLOUD_INFRASTRUCTURE.md"
          },
          "frontend": {
            "route": "#/functions/marketing-operations",
            "section": "tasks",
            "task_id": "marketing-operations.review-campaign-state",
            "interaction": "READ_CONTRACT",
            "effect_enabled": false
          },
          "acceptance": [
            {
              "id": "marketing-operations.AC2",
              "claim_scope": "CONTRACT_SPECIFICATION",
              "criterion": "No repetir efecto incierto; preparación parcial se reanuda desde owners existentes.",
              "test_ref": "roles/tests/test_business_functions.py",
              "evidence_refs": [],
              "target_status": "NOT_TESTED"
            }
          ],
          "owner_ref": "specs/library-extension/tasks.md"
        }
      ],
      "owner_binding_scope": "MAINTENANCE_FIXTURE"
    },
    {
      "id": "post-sales",
      "track_name": "Post-Sales",
      "function_kind": "BUSINESS_FUNCTION",
      "auth_role": false,
      "purpose": "Revisar entrega, onboarding y necesidades posteriores a la venta conservando el estado real.",
      "systems": [
        "Journey cliente/entrega/devolución existentes",
        "Agenda y aprobación de dominio existentes"
      ],
      "allowed": [
        "Leer contratos, fuentes y evidencia local.",
        "Preparar propuestas sin efectos externos dentro del scope aprobado."
      ],
      "prohibited": [
        "Derivar privilegios del título o track.",
        "Usar secretos o datos empresariales reales en el scaffold.",
        "Enviar mensajes, publicar, gastar, modificar CRM o ejecutar efectos desde esta consola."
      ],
      "status": {
        "specification_status": "SPEC_READY",
        "method_admission": {
          "status": "OFFICIAL_METHOD_REFERENCED",
          "claim_scope": "DOCUMENTARY_METHOD_ONLY",
          "source_refs": [
            "MS-SERVICE",
            "MS-CASES"
          ],
          "code_admission": "NONE"
        },
        "verification_status": {
          "scaffold": "NOT_TESTED",
          "target": "NOT_TESTED",
          "evidence_refs": []
        }
      },
      "source_refs": [
        "MS-SERVICE",
        "MS-CASES"
      ],
      "galaxy_admission": {
        "status": "WAITING_GALAXY_ADMISSION",
        "source_ref": null,
        "applies_only_to": "FUTURE_EVENT_MATERIAL",
        "content_invented": false,
        "reopen_trigger": "Material oficial identificable del evento disponible y obtenido legalmente; fijar fecha/revisión/hash/términos/claim y ejecutar admisión aplicable antes de incorporarlo."
      },
      "owner_refs": {
        "progress": "specs/library-extension/tasks.md",
        "execution_state": "PROJECT_EXECUTION_STATE.json",
        "execution_events": "PROJECT_EXECUTION_EVENTS.jsonl",
        "failures": "PROJECT_FAILURE_LESSONS.md",
        "dependency": "PROJECT_DEPENDENCY_UPDATE_RECORD.md",
        "freshness": "PROJECT_AUTHORITY_FRESHNESS_RECORD.md"
      },
      "resume": {
        "read_refs": [
          "PROJECT_EXECUTION_STATE.json",
          "PROJECT_EXECUTION_EVENTS.jsonl",
          "specs/library-extension/tasks.md",
          "PROJECT_FAILURE_LESSONS.md",
          "roles/business-functions.v403.json"
        ],
        "rule": "Validar checkpoint y hashes; consultar los owners existentes; invalidar evidencia afectada por delta; registrar siguiente acción allí, sin backlog paralelo."
      },
      "implementation_binding": {
        "status": "EXISTING_LIBRARY_REFERENCE",
        "pack": "implementation_packs/GO_FRANCHISE_CUSTOMER_JOURNEY_API.md",
        "evidence": "reconstruction_evidence/LIBRARY_SIGNED_REFERENCE_ACCEPTANCE_V402.md",
        "action": "Contratos existentes de checklist/entrega/devolución; límites por dominio.",
        "kind": "PACK_REFERENCE",
        "target_verified": false,
        "v403_executes_binding": false
      },
      "open_requirement": {
        "id": "TARGET-POST-SALES",
        "scope": "FUTURE_CONSUMER_ONLY",
        "description": "Acuerdo de servicio y aceptación target; Customer Success, renovación y garantías no son capacidades universales de este vínculo.",
        "owner_ref": "specs/library-extension/tasks.md"
      },
      "tasks": [
        {
          "id": "post-sales.review-delivery-onboarding",
          "title": "Revisar entrega y onboarding",
          "inputs": [
            "Pedido/servicio scoped",
            "acuerdo",
            "checklist y evidencia vigente"
          ],
          "outputs": [
            "Checklist y gaps vinculados al estado real"
          ],
          "product_permission": {
            "id": "post-sales:read",
            "definition_kind": "PROPOSED_TARGET_REQUIREMENT",
            "granted": false,
            "subject_binding": "REQUIRES_TARGET_CONFIGURATION",
            "organization_scope": "REQUIRED",
            "resource_scope": "REQUIRED",
            "enforcement": "EXISTING_PRODUCT_OWNER_OR_TARGET_GAP",
            "policy_ref": "SECURITY_SRE_CLOUD_INFRASTRUCTURE.md"
          },
          "frontend": {
            "route": "#/functions/post-sales",
            "section": "tasks",
            "task_id": "post-sales.review-delivery-onboarding",
            "interaction": "READ_CONTRACT",
            "effect_enabled": false
          },
          "acceptance": [
            {
              "id": "post-sales.AC1",
              "claim_scope": "CONTRACT_SPECIFICATION",
              "criterion": "Preparación local no equivale a entrega física, instalación ni aceptación del cliente.",
              "test_ref": "roles/tests/test_business_functions.py",
              "evidence_refs": [],
              "target_status": "NOT_TESTED"
            }
          ],
          "owner_ref": "specs/library-extension/tasks.md"
        },
        {
          "id": "post-sales.prepare-follow-up",
          "title": "Preparar seguimiento o escalación",
          "inputs": [
            "Estado de entrega",
            "consulta",
            "política de garantía/devolución"
          ],
          "outputs": [
            "Seguimiento propuesto con owner y aceptación pendiente"
          ],
          "product_permission": {
            "id": "post-sales:request",
            "definition_kind": "PROPOSED_TARGET_REQUIREMENT",
            "granted": false,
            "subject_binding": "REQUIRES_TARGET_CONFIGURATION",
            "organization_scope": "REQUIRED",
            "resource_scope": "REQUIRED",
            "enforcement": "EXISTING_PRODUCT_OWNER_OR_TARGET_GAP",
            "policy_ref": "SECURITY_SRE_CLOUD_INFRASTRUCTURE.md"
          },
          "frontend": {
            "route": "#/functions/post-sales",
            "section": "tasks",
            "task_id": "post-sales.prepare-follow-up",
            "interaction": "READ_CONTRACT",
            "effect_enabled": false
          },
          "acceptance": [
            {
              "id": "post-sales.AC2",
              "claim_scope": "CONTRACT_SPECIFICATION",
              "criterion": "No inventar garantía, renovación, fiscalidad ni devolución de dinero; usar el owner de dominio.",
              "test_ref": "roles/tests/test_business_functions.py",
              "evidence_refs": [],
              "target_status": "NOT_TESTED"
            }
          ],
          "owner_ref": "specs/library-extension/tasks.md"
        }
      ],
      "owner_binding_scope": "MAINTENANCE_FIXTURE"
    },
    {
      "id": "marketing",
      "track_name": "Marketing",
      "function_kind": "BUSINESS_FUNCTION",
      "auth_role": false,
      "purpose": "Relacionar mensaje, audiencia y activos con una hipótesis y revisión de campaña.",
      "systems": [
        "Contenido y campañas existentes",
        "Fuentes y reportes gobernados"
      ],
      "allowed": [
        "Leer contratos, fuentes y evidencia local.",
        "Preparar propuestas sin efectos externos dentro del scope aprobado."
      ],
      "prohibited": [
        "Derivar privilegios del título o track.",
        "Usar secretos o datos empresariales reales en el scaffold.",
        "Enviar mensajes, publicar, gastar, modificar CRM o ejecutar efectos desde esta consola."
      ],
      "status": {
        "specification_status": "SPEC_READY",
        "method_admission": {
          "status": "OFFICIAL_METHOD_REFERENCED",
          "claim_scope": "DOCUMENTARY_METHOD_ONLY",
          "source_refs": [
            "HUBSPOT-CAMPAIGNS",
            "DORA-USER"
          ],
          "code_admission": "NONE"
        },
        "verification_status": {
          "scaffold": "NOT_TESTED",
          "target": "NOT_TESTED",
          "evidence_refs": []
        }
      },
      "source_refs": [
        "HUBSPOT-CAMPAIGNS",
        "DORA-USER"
      ],
      "galaxy_admission": {
        "status": "WAITING_GALAXY_ADMISSION",
        "source_ref": null,
        "applies_only_to": "FUTURE_EVENT_MATERIAL",
        "content_invented": false,
        "reopen_trigger": "Material oficial identificable del evento disponible y obtenido legalmente; fijar fecha/revisión/hash/términos/claim y ejecutar admisión aplicable antes de incorporarlo."
      },
      "owner_refs": {
        "progress": "specs/library-extension/tasks.md",
        "execution_state": "PROJECT_EXECUTION_STATE.json",
        "execution_events": "PROJECT_EXECUTION_EVENTS.jsonl",
        "failures": "PROJECT_FAILURE_LESSONS.md",
        "dependency": "PROJECT_DEPENDENCY_UPDATE_RECORD.md",
        "freshness": "PROJECT_AUTHORITY_FRESHNESS_RECORD.md"
      },
      "resume": {
        "read_refs": [
          "PROJECT_EXECUTION_STATE.json",
          "PROJECT_EXECUTION_EVENTS.jsonl",
          "specs/library-extension/tasks.md",
          "PROJECT_FAILURE_LESSONS.md",
          "roles/business-functions.v403.json"
        ],
        "rule": "Validar checkpoint y hashes; consultar los owners existentes; invalidar evidencia afectada por delta; registrar siguiente acción allí, sin backlog paralelo."
      },
      "implementation_binding": {
        "status": "EXISTING_LIBRARY_REFERENCE",
        "pack": "implementation_packs/GO_CONNECTED_WHATSAPP_CAMPAIGNS.md",
        "evidence": "reconstruction_evidence/CAMPAIGN_CONNECTED_RELEASE_V402.md",
        "action": "Campaña existente y observación quote/order sin fórmula de atribución causal.",
        "kind": "PACK_REFERENCE",
        "target_verified": false,
        "v403_executes_binding": false
      },
      "open_requirement": {
        "id": "TARGET-MARKETING",
        "scope": "FUTURE_CONSUMER_ONLY",
        "description": "Publicación/canales/activos y medición del target; no se crea integración HubSpot ni autorización de publicación.",
        "owner_ref": "specs/library-extension/tasks.md"
      },
      "tasks": [
        {
          "id": "marketing.prepare-campaign-brief",
          "title": "Preparar brief de campaña",
          "inputs": [
            "Objetivo",
            "evidencia de audiencia",
            "restricciones",
            "claims con fuentes"
          ],
          "outputs": [
            "Brief y activos propuestos con revisión"
          ],
          "product_permission": {
            "id": "marketing:request",
            "definition_kind": "OBSERVED_EXISTING_OWNER_NAME",
            "granted": false,
            "subject_binding": "REQUIRES_TARGET_CONFIGURATION",
            "organization_scope": "REQUIRED",
            "resource_scope": "REQUIRED",
            "enforcement": "EXISTING_PRODUCT_OWNER_OR_TARGET_GAP",
            "policy_ref": "SECURITY_SRE_CLOUD_INFRASTRUCTURE.md"
          },
          "frontend": {
            "route": "#/functions/marketing",
            "section": "tasks",
            "task_id": "marketing.prepare-campaign-brief",
            "interaction": "READ_CONTRACT",
            "effect_enabled": false
          },
          "acceptance": [
            {
              "id": "marketing.AC1",
              "claim_scope": "CONTRACT_SPECIFICATION",
              "criterion": "No inventar beneficios, testimonios, claims regulatorios ni autoridad de publicación.",
              "test_ref": "roles/tests/test_business_functions.py",
              "evidence_refs": [],
              "target_status": "NOT_TESTED"
            }
          ],
          "owner_ref": "specs/library-extension/tasks.md"
        },
        {
          "id": "marketing.review-campaign-learning",
          "title": "Revisar aprendizaje de campaña",
          "inputs": [
            "Objetivo",
            "definición de métrica",
            "período/fuente",
            "resultados observados"
          ],
          "outputs": [
            "Informe con incertidumbres y siguiente hipótesis"
          ],
          "product_permission": {
            "id": "marketing:read",
            "definition_kind": "OBSERVED_EXISTING_OWNER_NAME",
            "granted": false,
            "subject_binding": "REQUIRES_TARGET_CONFIGURATION",
            "organization_scope": "REQUIRED",
            "resource_scope": "REQUIRED",
            "enforcement": "EXISTING_PRODUCT_OWNER_OR_TARGET_GAP",
            "policy_ref": "SECURITY_SRE_CLOUD_INFRASTRUCTURE.md"
          },
          "frontend": {
            "route": "#/functions/marketing",
            "section": "tasks",
            "task_id": "marketing.review-campaign-learning",
            "interaction": "READ_CONTRACT",
            "effect_enabled": false
          },
          "acceptance": [
            {
              "id": "marketing.AC2",
              "claim_scope": "CONTRACT_SPECIFICATION",
              "criterion": "Separar observación de atribución causal; no afirmar ROI sin fórmula/política y datos admitidos.",
              "test_ref": "roles/tests/test_business_functions.py",
              "evidence_refs": [],
              "target_status": "NOT_TESTED"
            }
          ],
          "owner_ref": "specs/library-extension/tasks.md"
        }
      ],
      "owner_binding_scope": "MAINTENANCE_FIXTURE"
    }
  ],
  "owner_binding": {
    "scope": "MAINTENANCE_FIXTURE",
    "project_id": "library-extension-v403",
    "require_existing_refs": true,
    "fixture_owners_must_not_be_inherited": true,
    "resume_gate": "implementation_packs/ENGINEERING_EXECUTION_VALIDATOR.md",
    "consumer_integration": "bind_consumer_owners.py resolves explicit existing owners inside the consumer root, rejects maintenance state/owners, resets claims. Run the consumer execution validator before resuming."
  },
  "actions_disabled_in_catalog": [
    "SEND_MESSAGE",
    "PUBLISH",
    "SPEND",
    "MUTATE_CRM",
    "RUN_PRODUCT_EFFECT"
  ],
  "catalog_restriction_scope": "READ_ONLY_LIBRARY_CONSOLE_ONLY",
  "consumer_authorization": "Consumer product applies its own explicitly authorized policy and existing enforcement; catalog restrictions do not prohibit separately authorized consumer operations."
}
````

### FILE: `roles/method-sources.v403.json`

```yaml
block_id: "BUSINESS-FUNCTION-OPERATING-V403:file3:v1"
operation: CREATE
provenance: AUTHORED
source: "local contract, provenance metadata, validator or owner-binding glue; no upstream method code"
license: "LicenseRef-Workspace-Owner"
sha256: "4bd1cc0ae8d211c84ce9bde611f3db9a70128ff40740b410465a4cc95b42f4e2"
variables: []
secrets_allowed: false
```

````json
{
  "schema_version": "1.0.0",
  "scope": "DOCUMENTARY_METHOD_REFERENCES",
  "retrieval_date": "2026-09-13",
  "upstream_code_acquired": false,
  "method_translation": "AUTHORED local task framing; scoped documentary references, not implementation admission or universal business rules.",
  "sources": [
    {
      "id": "MS-RBAC",
      "publisher": "Microsoft",
      "title": "What is Azure role-based access control (Azure RBAC)?",
      "url": "https://learn.microsoft.com/en-us/azure/role-based-access-control/overview",
      "query": "site.learn.microsoft.com azure role-based-access-control overview",
      "claim": "Separar principal, definición de permisos y scope al describir acceso.",
      "limit": "Semántica Azure; no asigna roles del producto ni instala Azure.",
      "method_status": "OFFICIAL_METHOD_REFERENCED",
      "architecture_code_admission": "NOT_REQUESTED",
      "upstream_code_copied": false,
      "redistribution": "LINK_AND_SHORT_LOCAL_PARAPHRASE_ONLY",
      "freshness": "Official page observed on retrieval_date; revalidate upon material change. Retrieval hash is not publisher signature.",
      "observation_receipt": "roles/evidence/source-observations.json"
    },
    {
      "id": "DORA-USER",
      "publisher": "Google / DORA",
      "title": "User-centric focus",
      "url": "https://dora.dev/capabilities/user-centric-focus/",
      "query": "site.dora.dev capabilities user-centricity",
      "claim": "Relacionar necesidades, feedback y prioridades con el valor para el usuario.",
      "limit": "Referencia de método; sin porcentajes de rendimiento ni causalidad transferidos.",
      "method_status": "OFFICIAL_METHOD_REFERENCED",
      "architecture_code_admission": "NOT_REQUESTED",
      "upstream_code_copied": false,
      "redistribution": "LINK_AND_SHORT_LOCAL_PARAPHRASE_ONLY",
      "freshness": "Official page observed on retrieval_date; revalidate upon material change. Retrieval hash is not publisher signature.",
      "observation_receipt": "roles/evidence/source-observations.json"
    },
    {
      "id": "DORA-BATCHES",
      "publisher": "Google / DORA",
      "title": "Working in small batches",
      "url": "https://dora.dev/capabilities/working-in-small-batches/",
      "query": "site.dora.dev capabilities working in small batches",
      "claim": "Dividir cambios en unidades pequeñas y verificables para obtener feedback temprano.",
      "limit": "No fija plazo prometido, cadencia o umbral universal; el proyecto conserva sus gates.",
      "method_status": "OFFICIAL_METHOD_REFERENCED",
      "architecture_code_admission": "NOT_REQUESTED",
      "upstream_code_copied": false,
      "redistribution": "LINK_AND_SHORT_LOCAL_PARAPHRASE_ONLY",
      "freshness": "Official page observed on retrieval_date; revalidate upon material change. Retrieval hash is not publisher signature.",
      "observation_receipt": "roles/evidence/source-observations.json"
    },
    {
      "id": "MS-SALES",
      "publisher": "Microsoft",
      "title": "Prospect to quote end-to-end business process flow overview",
      "url": "https://learn.microsoft.com/en-us/dynamics365/guidance/business-processes/prospect-to-quote-overview",
      "query": "site.learn.microsoft.com dynamics365 guidance business-processes prospect-to-quote overview",
      "claim": "Distinguir relaciones, leads, oportunidades y cotización en el flujo comercial.",
      "limit": "Proceso de referencia Dynamics 365; no selecciona Business Central ni define scoring/precios.",
      "method_status": "OFFICIAL_METHOD_REFERENCED",
      "architecture_code_admission": "NOT_REQUESTED",
      "upstream_code_copied": false,
      "redistribution": "LINK_AND_SHORT_LOCAL_PARAPHRASE_ONLY",
      "freshness": "Official page observed on retrieval_date; revalidate upon material change. Retrieval hash is not publisher signature.",
      "observation_receipt": "roles/evidence/source-observations.json"
    },
    {
      "id": "MS-CASES",
      "publisher": "Microsoft",
      "title": "Help organizations manage and optimize their case to resolution processes with Dynamics 365",
      "url": "https://learn.microsoft.com/en-us/dynamics365/guidance/business-processes/case-to-resolution-introduction",
      "query": "site.learn.microsoft.com dynamics365 guidance case to resolution business process overview",
      "claim": "Distinguir registro, asignación, investigación, resolución y cierre con responsables.",
      "limit": "El consumidor define variaciones y políticas; no prueba SLA ni ticketing universal.",
      "method_status": "OFFICIAL_METHOD_REFERENCED",
      "architecture_code_admission": "NOT_REQUESTED",
      "upstream_code_copied": false,
      "redistribution": "LINK_AND_SHORT_LOCAL_PARAPHRASE_ONLY",
      "freshness": "Official page observed on retrieval_date; revalidate upon material change. Retrieval hash is not publisher signature.",
      "observation_receipt": "roles/evidence/source-observations.json"
    },
    {
      "id": "MS-SERVICE",
      "publisher": "Microsoft",
      "title": "Streamline service delivery and invoicing processes with the service to deliver end-to-end scenario",
      "url": "https://learn.microsoft.com/en-us/dynamics365/guidance/business-processes/service-to-cash-introduction",
      "query": "service to deliver introduction business processes site:learn.microsoft.com",
      "claim": "Relacionar solicitud, recursos y estado de entrega al revisar seguimiento posterior a la venta.",
      "limit": "La URL conserva service-to-cash; el texto explica el cambio a service-to-deliver desde febrero 2025. No adopta facturación ni garantías.",
      "method_status": "OFFICIAL_METHOD_REFERENCED",
      "architecture_code_admission": "NOT_REQUESTED",
      "upstream_code_copied": false,
      "redistribution": "LINK_AND_SHORT_LOCAL_PARAPHRASE_ONLY",
      "freshness": "Official page observed on retrieval_date; revalidate upon material change. Retrieval hash is not publisher signature.",
      "observation_receipt": "roles/evidence/source-observations.json"
    },
    {
      "id": "HUBSPOT-CAMPAIGNS",
      "publisher": "HubSpot",
      "title": "Create campaigns",
      "url": "https://knowledge.hubspot.com/campaigns/create-campaigns",
      "query": "site.knowledge.hubspot.com campaigns create campaigns",
      "claim": "Documentar objetivo y activos de campaña distinguiendo permisos del producto.",
      "limit": "No integra API ni requiere suscripción HubSpot. No importa Super Admin como privilegio de Marketing.",
      "method_status": "OFFICIAL_METHOD_REFERENCED",
      "architecture_code_admission": "NOT_REQUESTED",
      "upstream_code_copied": false,
      "redistribution": "LINK_AND_SHORT_LOCAL_PARAPHRASE_ONLY",
      "freshness": "Official page observed on retrieval_date; revalidate upon material change. Retrieval hash is not publisher signature.",
      "observation_receipt": "roles/evidence/source-observations.json"
    },
    {
      "id": "XAI-GROK",
      "publisher": "xAI",
      "title": "Welcome to Grok",
      "url": "https://docs.x.ai/grok/overview",
      "query": "Grok API overview documentation site:docs.x.ai",
      "claim": "Orientar aprendizaje con documentación oficial que enlaza Grok Bot.",
      "limit": "No es material Galaxy, currículo confirmado, runtime admitido ni disponibilidad de una cuenta.",
      "method_status": "OFFICIAL_METHOD_REFERENCED",
      "architecture_code_admission": "NOT_REQUESTED",
      "upstream_code_copied": false,
      "redistribution": "LINK_AND_SHORT_LOCAL_PARAPHRASE_ONLY",
      "freshness": "Official page observed on retrieval_date; revalidate upon material change. Retrieval hash is not publisher signature.",
      "observation_receipt": "roles/evidence/source-observations.json"
    }
  ]
}
````

### FILE: `roles/library-bindings.v403.json`

```yaml
block_id: "BUSINESS-FUNCTION-OPERATING-V403:file4:v1"
operation: CREATE
provenance: AUTHORED
source: "local contract, provenance metadata, validator or owner-binding glue; no upstream method code"
license: "LicenseRef-Workspace-Owner"
sha256: "50b63aa7bff27c96d7260e08fd8f506e09481d16001a63926931aa5f48262005"
variables: []
secrets_allowed: false
```

````json
{
  "schema_version": "1.0.0",
  "library_baseline_revision": 337,
  "bindings": {
    "engineering": {
      "pack": "implementation_packs/ENGINEERING_EXECUTION_VALIDATOR.md",
      "evidence": "reconstruction_evidence/LIBRARY_INFRA_READY_V402.md",
      "pack_sha256": "368c4f9d0c328353634e3dc0738e0d300f225013b9c7a769628b971d0c854808",
      "evidence_sha256": "97c433e986c3beb29b183dc30c5ccdf7b0721156268bfcaab8a559554b46b883",
      "claim_scope": "EXISTING_LIBRARY_REFERENCE_ONLY",
      "target_verified": false
    },
    "product-managers": {
      "pack": "implementation_packs/GO_CUSTOMER_SURVEY_API.md",
      "evidence": "reconstruction_evidence/CONNECTED_CUSTOMER_SURVEYS_V400.md",
      "pack_sha256": "7638176e48ddd88f90a08c92108eb976e9e5bb770ebebcc77fc4ce78b619ade7",
      "evidence_sha256": "29970739d1d1e8d77a22c12619677f4fd0dcb491a6149b2c07ddcb394c4f2ca9",
      "claim_scope": "EXISTING_LIBRARY_REFERENCE_ONLY",
      "target_verified": false
    },
    "sales-engineering": {
      "pack": "implementation_packs/GO_FRANCHISE_CUSTOMER_JOURNEY_API.md",
      "evidence": "reconstruction_evidence/LIBRARY_SIGNED_REFERENCE_ACCEPTANCE_V402.md",
      "pack_sha256": "c4989e27c5ad7881362819ede39b8127019b119c4a5509e9c3d22d26b4d0cc78",
      "evidence_sha256": "fd98b23c028f7ee84da423460d6e51252196b4154282e9243371f466dc095bdf",
      "claim_scope": "EXISTING_LIBRARY_REFERENCE_ONLY",
      "target_verified": false
    },
    "sales": {
      "pack": "implementation_packs/GO_FRANCHISE_CUSTOMER_JOURNEY_API.md",
      "evidence": "reconstruction_evidence/LIBRARY_SIGNED_REFERENCE_ACCEPTANCE_V402.md",
      "pack_sha256": "c4989e27c5ad7881362819ede39b8127019b119c4a5509e9c3d22d26b4d0cc78",
      "evidence_sha256": "fd98b23c028f7ee84da423460d6e51252196b4154282e9243371f466dc095bdf",
      "claim_scope": "EXISTING_LIBRARY_REFERENCE_ONLY",
      "target_verified": false
    },
    "sdrs": {
      "pack": "implementation_packs/GO_LEAD_CANDIDATE_PROMOTION.md",
      "evidence": "reconstruction_evidence/LIBRARY_INFRA_READY_V402.md",
      "pack_sha256": "d04bbef352bcd931fbda523f5f139b6d7765a131f62c772dcc55dc753d83d92e",
      "evidence_sha256": "97c433e986c3beb29b183dc30c5ccdf7b0721156268bfcaab8a559554b46b883",
      "claim_scope": "EXISTING_LIBRARY_REFERENCE_ONLY",
      "target_verified": false
    },
    "customer-support": {
      "pack": "implementation_packs/GO_CONNECTED_CONVERSATION_RUNTIME.md",
      "evidence": "reconstruction_evidence/AI_CONNECTED_REFERENCE_RELEASE_V402.md",
      "pack_sha256": "8bbbe6f46e3f2af027ad9724d9abce4a2ecaf10fc7c3310c20f8daa36fddf288",
      "evidence_sha256": "5352c625aa0dffa700a559ba3b4f9d2795ab7d9816cebf8045747cd0954759ca",
      "claim_scope": "EXISTING_LIBRARY_REFERENCE_ONLY",
      "target_verified": false
    },
    "marketing-operations": {
      "pack": "implementation_packs/GO_CONNECTED_WHATSAPP_CAMPAIGNS.md",
      "evidence": "reconstruction_evidence/CAMPAIGN_CONNECTED_RELEASE_V402.md",
      "pack_sha256": "d3b0aef726f76fe44b0b74a990760bbbc57a635ba39e402ce8669cc652e28788",
      "evidence_sha256": "4fc7deec96b3ddc0442e04adffe5c00263e0c5597eb2398e83233728e61e611d",
      "claim_scope": "EXISTING_LIBRARY_REFERENCE_ONLY",
      "target_verified": false
    },
    "post-sales": {
      "pack": "implementation_packs/GO_FRANCHISE_CUSTOMER_JOURNEY_API.md",
      "evidence": "reconstruction_evidence/LIBRARY_SIGNED_REFERENCE_ACCEPTANCE_V402.md",
      "pack_sha256": "c4989e27c5ad7881362819ede39b8127019b119c4a5509e9c3d22d26b4d0cc78",
      "evidence_sha256": "fd98b23c028f7ee84da423460d6e51252196b4154282e9243371f466dc095bdf",
      "claim_scope": "EXISTING_LIBRARY_REFERENCE_ONLY",
      "target_verified": false
    },
    "marketing": {
      "pack": "implementation_packs/GO_CONNECTED_WHATSAPP_CAMPAIGNS.md",
      "evidence": "reconstruction_evidence/CAMPAIGN_CONNECTED_RELEASE_V402.md",
      "pack_sha256": "d3b0aef726f76fe44b0b74a990760bbbc57a635ba39e402ce8669cc652e28788",
      "evidence_sha256": "4fc7deec96b3ddc0442e04adffe5c00263e0c5597eb2398e83233728e61e611d",
      "claim_scope": "EXISTING_LIBRARY_REFERENCE_ONLY",
      "target_verified": false
    }
  },
  "new_v403_product_execution": false
}
````

### FILE: `roles/validate_business_functions.py`

```yaml
block_id: "BUSINESS-FUNCTION-OPERATING-V403:file5:v1"
operation: CREATE
provenance: AUTHORED
source: "local contract, provenance metadata, validator or owner-binding glue; no upstream method code"
license: "LicenseRef-Workspace-Owner"
sha256: "e8c730b913e18e75cadd7e3347711ffcce09aeacbc865a18fb5a5f9db0d2e201"
variables: []
secrets_allowed: false
```

````python
"""AUTHORED contract glue. Validates library scaffolding, never product permission or readiness."""
from __future__ import annotations

import argparse
import copy
import hashlib
import json
from pathlib import Path, PurePosixPath
import re
import sys
from urllib.parse import urlsplit

MAX_JSON_BYTES = 2 * 1024 * 1024
EXPECTED_TRACKS = {
    "grok-bot-101": "Grok Bot 101", "engineering": "Engineering",
    "product-managers": "Product Managers", "founders": "Founders",
    "sales-engineering": "Sales Engineering", "sales": "Sales", "sdrs": "SDRs",
    "customer-support": "Customer Support", "marketing-operations": "Marketing Operations",
    "post-sales": "Post-Sales", "marketing": "Marketing",
}
OWNER_KEYS = {"progress", "execution_state", "execution_events", "failures", "dependency", "freshness"}
OWNER_FIXTURE = "specs/library-extension/tasks.md"
MAINTENANCE_PREFIXES = ("elite-library", "library-extension-v403")
OFFICIAL_HOSTS = {"learn.microsoft.com", "dora.dev", "knowledge.hubspot.com", "docs.x.ai"}
SECTIONS = {"purpose", "systems", "tasks", "permissions", "allowed", "prohibited", "inputs",
            "outputs", "acceptance", "implementation", "sources", "status", "owners", "resume"}
ROOT_KEYS = {"schema_version", "contract_id", "scope", "claim_scope", "production_authorized",
             "provenance", "frontend_contract", "owner_refs", "permission_policy", "done_policy",
             "functions", "owner_binding", "actions_disabled_in_catalog",
             "catalog_restriction_scope", "consumer_authorization"}
FUNCTION_KEYS = {"id", "track_name", "function_kind", "auth_role", "purpose", "systems", "allowed",
                 "prohibited", "status", "source_refs", "galaxy_admission", "owner_refs", "resume",
                 "implementation_binding", "open_requirement", "tasks", "owner_binding_scope"}

def require(condition, code):
    if not condition:
        raise ValueError(code)

def closed(value, keys, code):
    require(type(value) is dict and set(value) == set(keys), code)

def strings(value, code, maximum=32):
    require(type(value) is list and 0 < len(value) <= maximum, code)
    require(all(type(v) is str and v.strip() and len(v) <= 4000 for v in value), code)

def sha256(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()

def pairs(items):
    result = {}
    for key, value in items:
        require(key not in result, "DUPLICATE_JSON_KEY")
        result[key] = value
    return result

def load_json(path):
    require(path.stat().st_size <= MAX_JSON_BYTES, "JSON_SIZE_LIMIT")
    return json.loads(path.read_text(encoding="utf-8"), object_pairs_hook=pairs,
                      parse_constant=lambda value: (_ for _ in ()).throw(ValueError("NONFINITE_JSON")))

def confined(root, reference, must_exist=True):
    require(type(reference) is str and reference.strip(), "MISSING_REFERENCE")
    require("\\" not in reference and ":" not in reference and "\x00" not in reference,
            "REFERENCE_SYNTAX")
    relative = PurePosixPath(reference)
    require(not relative.is_absolute() and ".." not in relative.parts and "." != reference,
            "REFERENCE_ESCAPE")
    base = root.resolve()
    path = (base / reference).resolve()
    require(path.is_relative_to(base), "REFERENCE_ESCAPE")
    if must_exist:
        require(path.is_file(), "MISSING_REFERENCE")
    return path

def check_owners(record, project_root, check_files):
    owners = record["owner_refs"]
    closed(owners, OWNER_KEYS, "OWNER_FIELDS")
    binding = record["owner_binding"]
    closed(binding, {"scope", "project_id", "require_existing_refs", "fixture_owners_must_not_be_inherited",
                     "resume_gate", "consumer_integration"}, "OWNER_BINDING_FIELDS")
    require(binding["scope"] in {"MAINTENANCE_FIXTURE", "CONSUMER"}, "OWNER_BINDING_SCOPE")
    require(binding["require_existing_refs"] is True and binding["fixture_owners_must_not_be_inherited"] is True,
            "OWNER_INHERITANCE_CONTROL")
    require(type(binding["project_id"]) is str and binding["project_id"].strip(), "OWNER_PROJECT_ID")
    for reference in owners.values():
        confined(project_root, reference, must_exist=check_files)
    require(len(set(owners.values())) == len(owners), "OWNER_COLLISION")
    if binding["scope"] == "MAINTENANCE_FIXTURE":
        require(binding["project_id"] == "library-extension-v403" and owners["progress"] == OWNER_FIXTURE,
                "FIXTURE_OWNER_DRIFT")
    else:
        forbidden = ("specs/library-extension/", "specs/library-maintenance/")
        require(not binding["project_id"].startswith(MAINTENANCE_PREFIXES), "MAINTENANCE_ID_IN_CONSUMER")
        require(not any(any(part in ref.lower() for part in forbidden) for ref in owners.values()),
                "MAINTENANCE_OWNER_IN_CONSUMER")
    if check_files:
        state = load_json(confined(project_root, owners["execution_state"]))
        require(state.get("project", {}).get("id") == binding["project_id"], "OWNER_STATE_ID_MISMATCH")
        if binding["scope"] == "CONSUMER":
            require(not state["project"]["id"].startswith(MAINTENANCE_PREFIXES), "MAINTENANCE_STATE_IN_CONSUMER")
    return owners

def validate(record, sources, project_root, library_root, check_files=True,
             contract_path=None, bindings=None, observations=None):
    closed(record, ROOT_KEYS, "ROOT_FIELDS")
    require(record["schema_version"] == "1.0.0" and record["contract_id"] == "BUSINESS-FUNCTIONS-V403",
            "SCHEMA_ID")
    require(record["scope"] == "LIBRARY_INFRASTRUCTURE" and
            record["claim_scope"] == "FUNCTION_CONTRACT_SCAFFOLD" and
            record["production_authorized"] is False, "SCOPE_ESCALATION")
    closed(record["provenance"], {"classification", "kind", "upstream_code_copied", "official_method_code_claim"},
           "PROVENANCE_FIELDS")
    require(record["provenance"] == {"classification": "AUTHORED", "kind": "CONTRACT_AND_VALIDATOR_GLUE",
                                    "upstream_code_copied": False, "official_method_code_claim": False},
            "FALSE_OFFICIAL_CODE")
    owners = check_owners(record, project_root, check_files)
    policy = record["permission_policy"]
    closed(policy, {"track_is_auth_role", "default", "effective_grants", "requires",
                    "unbound_permission_behavior", "inherited_privilege_allowed"}, "PERMISSION_POLICY_FIELDS")
    require(policy["track_is_auth_role"] is False and policy["default"] == "DENY" and
            policy["effective_grants"] == [] and policy["inherited_privilege_allowed"] is False,
            "IMPLICIT_PRIVILEGE")
    require(set(policy["requires"]) == {"authenticated_subject", "explicit_product_permission",
            "organization_scope", "resource_scope", "server_side_enforcement", "audit_evidence"},
            "AUTHORIZATION_DIMENSIONS")
    require(policy["unbound_permission_behavior"] == "BLOCK_PRODUCT_EFFECT", "UNBOUND_EFFECT")
    frontend = record["frontend_contract"]
    closed(frontend, {"route_pattern", "mode", "required_sections", "external_effects_enabled"}, "FRONTEND_FIELDS")
    require(frontend["route_pattern"] == "#/functions/{id}" and
            frontend["mode"] == "READ_ONLY_LIBRARY_CONSOLE" and frontend["external_effects_enabled"] is False,
            "CATALOG_EFFECT")
    require(set(frontend["required_sections"]) == SECTIONS, "FRONTEND_COVERAGE")
    require(record["catalog_restriction_scope"] == "READ_ONLY_LIBRARY_CONSOLE_ONLY", "CATALOG_SCOPE")
    require(set(record["actions_disabled_in_catalog"]) == {"SEND_MESSAGE", "PUBLISH", "SPEND", "MUTATE_CRM",
                                                         "RUN_PRODUCT_EFFECT"}, "CATALOG_EFFECT")
    require(type(record["consumer_authorization"]) is str and record["consumer_authorization"].strip(),
            "CONSUMER_AUTHORITY")
    closed(record["done_policy"], {"scaffold", "product", "galaxy"}, "DONE_POLICY")
    require(all(type(v) is str and v.strip() for v in record["done_policy"].values()), "DONE_POLICY")

    require(sources.get("scope") == "DOCUMENTARY_METHOD_REFERENCES" and sources.get("upstream_code_acquired") is False,
            "SOURCE_SCOPE")
    source_list = sources.get("sources", [])
    source_ids = [s["id"] for s in source_list]
    require(source_ids and len(source_ids) == len(set(source_ids)), "SOURCE_DUPLICATE")
    for source in source_list:
        parts = urlsplit(source["url"])
        require(parts.scheme == "https" and parts.hostname in OFFICIAL_HOSTS and not parts.username,
                "UNOFFICIAL_SOURCE")
        require(source.get("method_status") == "OFFICIAL_METHOD_REFERENCED" and
                source.get("upstream_code_copied") is False and
                source.get("architecture_code_admission") == "NOT_REQUESTED", "SOURCE_PROMOTION")
        require(all(type(source.get(k)) is str and source[k].strip() for k in
                    ("query", "claim", "limit", "publisher", "title", "freshness")), "SOURCE_PROVENANCE")
    if check_files:
        require(type(observations) is dict, "MISSING_SOURCE_OBSERVATIONS")
        obs = {o["id"]: o for o in observations.get("observations", [])}
        require(len(obs) == len(source_ids), "SOURCE_OBSERVATION_COUNT")
        for source in source_list:
            receipt = obs.get(source["id"], {})
            require(receipt.get("status") == "OBSERVED" and receipt.get("http_status") == 200 and
                    receipt.get("url") == source["url"] and receipt.get("body_stored") is False and
                    re.fullmatch("[0-9a-f]{64}", receipt.get("response_sha256", "")) is not None,
                    "SOURCE_OBSERVATION_MISMATCH")

    functions = record["functions"]
    require(type(functions) is list and len(functions) == 11, "FUNCTION_COUNT")
    require({f.get("id") for f in functions} == set(EXPECTED_TRACKS), "FUNCTION_SET")
    task_ids, acceptance_ids = set(), set()
    for function in functions:
        closed(function, FUNCTION_KEYS, "FUNCTION_FIELDS")
        fid = function["id"]
        require(function["track_name"] == EXPECTED_TRACKS[fid], "TRACK_NAME")
        kind = "LEARNING_TRACK" if fid == "grok-bot-101" else "BUSINESS_FUNCTION"
        require(function["function_kind"] == kind and function["auth_role"] is False, "TRACK_IS_NOT_ROLE")
        require(function["owner_binding_scope"] == record["owner_binding"]["scope"] and
                function["owner_refs"] == owners, "OWNER_BINDING_MISMATCH")
        require(type(function["purpose"]) is str and function["purpose"].strip(), "PURPOSE")
        for key in ("systems", "allowed", "prohibited", "source_refs"):
            strings(function[key], "FUNCTION_" + key.upper())
        require(set(function["source_refs"]).issubset(source_ids), "UNKNOWN_METHOD_SOURCE")
        closed(function["status"], {"specification_status", "method_admission", "verification_status"}, "STATUS_FIELDS")
        require(function["status"]["specification_status"] == "SPEC_READY", "SPECIFICATION_STATUS")
        method = function["status"]["method_admission"]
        require(method == {"status": "OFFICIAL_METHOD_REFERENCED", "claim_scope": "DOCUMENTARY_METHOD_ONLY",
                           "source_refs": function["source_refs"], "code_admission": "NONE"}, "METHOD_PROMOTION")
        verification = function["status"]["verification_status"]
        closed(verification, {"scaffold", "target", "evidence_refs"}, "VERIFICATION_FIELDS")
        require(verification["target"] == "NOT_TESTED", "TARGET_PROOF_OUT_OF_SCOPE")
        require(verification["scaffold"] in {"NOT_TESTED", "PROVEN"}, "VERIFICATION_STATUS")
        require(type(verification["evidence_refs"]) is list, "EVIDENCE_LIST")
        if verification["scaffold"] == "PROVEN":
            require(bool(verification["evidence_refs"]) and contract_path is not None, "PROVEN_WITHOUT_EVIDENCE")
            for reference in verification["evidence_refs"]:
                proof = load_json(confined(project_root, reference))
                require(proof.get("result") == "PASS" and proof.get("claim_scope") == "FUNCTION_CONTRACT_SCAFFOLD"
                        and proof.get("contract_sha256") == sha256(contract_path)
                        and fid in proof.get("function_ids", []), "UNBOUND_PROOF")
        galaxy = function["galaxy_admission"]
        closed(galaxy, {"status", "source_ref", "applies_only_to", "content_invented", "reopen_trigger"}, "GALAXY_FIELDS")
        require(galaxy["status"] == "WAITING_GALAXY_ADMISSION" and galaxy["source_ref"] is None and
                galaxy["applies_only_to"] == "FUTURE_EVENT_MATERIAL" and galaxy["content_invented"] is False
                and type(galaxy["reopen_trigger"]) is str and bool(galaxy["reopen_trigger"].strip()),
                "FUTURE_GALAXY_CLAIM")
        closed(function["resume"], {"read_refs", "rule"}, "RESUME_FIELDS")
        require(set(function["resume"]["read_refs"]) >= {owners["execution_state"], owners["execution_events"],
                owners["progress"], owners["failures"]}, "RESUME_OWNERS")
        for reference in function["resume"]["read_refs"]:
            confined(project_root, reference, must_exist=check_files)
        open_req = function["open_requirement"]
        closed(open_req, {"id", "scope", "description", "owner_ref"}, "OPEN_REQUIREMENT_FIELDS")
        require(open_req["scope"] == "FUTURE_CONSUMER_ONLY" and open_req["owner_ref"] == owners["progress"]
                and type(open_req["description"]) is str and open_req["description"].strip(), "OPEN_REQUIREMENT")
        impl = function["implementation_binding"]
        closed(impl, {"status", "pack", "evidence", "action", "kind", "target_verified", "v403_executes_binding"},
               "IMPLEMENTATION_FIELDS")
        require(impl["target_verified"] is False and impl["v403_executes_binding"] is False, "IMPLICIT_RUNTIME")
        require(impl["status"] in {"EXISTING_LIBRARY_REFERENCE", "DOCUMENTED_ONLY"}, "IMPLEMENTATION_STATUS")
        if impl["status"] == "DOCUMENTED_ONLY":
            require(impl["pack"] is None and impl["evidence"] is None and impl["action"] is None
                    and impl["kind"] == "NO_PRODUCT_RUNTIME", "UNDOCUMENTED_RUNTIME")
        else:
            require(type(impl["action"]) is str and impl["action"].strip() and
                    impl["kind"] in {"CLI_CONTRACT", "PACK_REFERENCE", "OBSERVED_API_ROUTE"}, "IMPLEMENTATION_ACTION")
            pack = confined(library_root, impl["pack"], must_exist=check_files)
            evidence = confined(library_root, impl["evidence"], must_exist=check_files)
            if check_files:
                require(type(bindings) is dict and fid in bindings, "MISSING_BINDING_RECEIPT")
                bound = bindings[fid]
                require(bound.get("pack") == impl["pack"] and bound.get("evidence") == impl["evidence"] and
                        bound.get("pack_sha256") == sha256(pack) and bound.get("evidence_sha256") == sha256(evidence),
                        "LIBRARY_BINDING_CHANGED")
                if impl["kind"] == "OBSERVED_API_ROUTE":
                    require("/v1/franchise/marketing/campaigns" in pack.read_text(encoding="utf-8"),
                            "ROUTE_NOT_IN_OWNER")
        require(type(function["tasks"]) is list and len(function["tasks"]) == 2, "TASK_COUNT")
        for task in function["tasks"]:
            closed(task, {"id", "title", "inputs", "outputs", "product_permission", "frontend", "acceptance",
                          "owner_ref"}, "TASK_FIELDS")
            tid = task["id"]
            require(type(tid) is str and tid.startswith(fid + ".") and tid not in task_ids, "TASK_ID")
            task_ids.add(tid)
            require(task["owner_ref"] == owners["progress"], "TASK_OWNER")
            require(type(task["title"]) is str and task["title"].strip(), "TASK_TITLE")
            strings(task["inputs"], "INPUT_CONTRACT")
            strings(task["outputs"], "OUTPUT_CONTRACT")
            permission = task["product_permission"]
            closed(permission, {"id", "definition_kind", "granted", "subject_binding", "organization_scope",
                                "resource_scope", "enforcement", "policy_ref"}, "PRODUCT_PERMISSION_FIELDS")
            require(re.fullmatch("[a-z][a-z-]*:[a-z][a-z-]*", permission["id"]) is not None and
                    permission["granted"] is False and
                    permission["subject_binding"] == "REQUIRES_TARGET_CONFIGURATION" and
                    permission["organization_scope"] == "REQUIRED" and permission["resource_scope"] == "REQUIRED",
                    "PRODUCT_PERMISSION_GRANT")
            require(permission["enforcement"] == "EXISTING_PRODUCT_OWNER_OR_TARGET_GAP", "ENFORCEMENT")
            require(permission["definition_kind"] in {"PROPOSED_TARGET_REQUIREMENT", "OBSERVED_EXISTING_OWNER_NAME"},
                    "PERMISSION_DEFINITION")
            confined(library_root, permission["policy_ref"], must_exist=check_files)
            if check_files and permission["definition_kind"] == "OBSERVED_EXISTING_OWNER_NAME":
                require(impl["status"] == "EXISTING_LIBRARY_REFERENCE" and permission["id"] in
                        confined(library_root, impl["pack"]).read_text(encoding="utf-8"), "PERMISSION_NOT_IN_OWNER")
            ui = task["frontend"]
            require(ui == {"route": "#/functions/" + fid, "section": "tasks", "task_id": tid,
                           "interaction": "READ_CONTRACT", "effect_enabled": False}, "TASK_FRONTEND_BINDING")
            require(type(task["acceptance"]) is list and task["acceptance"], "MISSING_ACCEPTANCE")
            for acceptance in task["acceptance"]:
                closed(acceptance, {"id", "claim_scope", "criterion", "test_ref", "evidence_refs", "target_status"},
                       "ACCEPTANCE_FIELDS")
                aid = acceptance["id"]
                require(type(aid) is str and aid.startswith(fid + ".") and aid not in acceptance_ids, "ACCEPTANCE_ID")
                acceptance_ids.add(aid)
                require(acceptance["claim_scope"] == "CONTRACT_SPECIFICATION" and
                        acceptance["target_status"] == "NOT_TESTED" and acceptance["evidence_refs"] == [],
                        "ACCEPTANCE_SCOPE")
                require(type(acceptance["criterion"]) is str and acceptance["criterion"].strip(), "ACCEPTANCE_CRITERION")
                confined(project_root, acceptance["test_ref"], must_exist=check_files)
    return {"result": "PASS", "claim_scope": "FUNCTION_CONTRACT_SCAFFOLD", "function_ids": sorted(EXPECTED_TRACKS),
            "function_count": len(functions), "task_count": len(task_ids), "acceptance_count": len(acceptance_ids),
            "owner_binding_scope": record["owner_binding"]["scope"], "production_authorized": False,
            "business_operations_verified": False, "frontend_render_verified": False,
            "galaxy_event_admitted": False, "known_open_scaffold_defects": "SEE_EXISTING_FAILURE_OWNER",
            "scope_limit": "Structural/provenance/link controls only; frontend execution and human/business acceptance are separate."}

def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--project-root", type=Path, required=True)
    parser.add_argument("--library-root", type=Path, required=True)
    parser.add_argument("--contract", default="roles/business-functions.v403.json")
    parser.add_argument("--report", required=True)
    args = parser.parse_args()
    output = confined(args.project_root, args.report, must_exist=False)
    require(not output.exists(), "REPORT_ALREADY_EXISTS")
    try:
        contract_path = confined(args.project_root, args.contract)
        source_path = confined(args.project_root, "roles/method-sources.v403.json")
        result = validate(load_json(contract_path), load_json(source_path), args.project_root, args.library_root,
                          contract_path=contract_path,
                          bindings=load_json(confined(args.project_root, "roles/library-bindings.v403.json"))["bindings"],
                          observations=load_json(confined(args.project_root, "roles/evidence/source-observations.json")))
        result.update(contract_sha256=sha256(contract_path), sources_sha256=sha256(source_path),
                      validator_sha256=sha256(Path(__file__)), python_version=sys.version)
        exit_code = 0
    except (ValueError, KeyError, TypeError, OSError, json.JSONDecodeError) as exc:
        result = {"result": "FAIL", "claim_scope": "FUNCTION_CONTRACT_SCAFFOLD", "error": str(exc),
                  "production_authorized": False}
        exit_code = 2
    output.parent.mkdir(parents=True, exist_ok=True)
    with output.open("x", encoding="utf-8", newline="\n") as stream:
        stream.write(json.dumps(result, ensure_ascii=False, indent=2) + "\n")
    print(json.dumps(result, ensure_ascii=False))
    return exit_code

if __name__ == "__main__":
    raise SystemExit(main())
````

### FILE: `roles/bind_consumer_owners.py`

```yaml
block_id: "BUSINESS-FUNCTION-OPERATING-V403:file6:v1"
operation: CREATE
provenance: AUTHORED
source: "local contract, provenance metadata, validator or owner-binding glue; no upstream method code"
license: "LicenseRef-Workspace-Owner"
sha256: "751e01b40aad812e30892d3d34e3ab97619aa582f00c768ae316923f5d9488f9"
variables: []
secrets_allowed: false
```

````python
"""AUTHORED owner-binding glue. Never copies progress, grants, approvals, or execution state."""
from __future__ import annotations
import argparse
import copy
import json
from pathlib import Path
from validate_business_functions import MAINTENANCE_PREFIXES, check_owners, confined, load_json, require

def bind_existing_owners(source, owners, project_root, output_ref):
    record = copy.deepcopy(source)
    state = load_json(confined(project_root, owners["execution_state"]))
    project_id = state.get("project", {}).get("id")
    require(type(project_id) is str and project_id and not project_id.startswith(MAINTENANCE_PREFIXES),
            "MAINTENANCE_STATE_IN_CONSUMER")
    record["owner_refs"] = copy.deepcopy(owners)
    record["owner_binding"]["scope"] = "CONSUMER"
    record["owner_binding"]["project_id"] = project_id
    check_owners(record, project_root, check_files=True)
    confined(project_root, output_ref, must_exist=False)
    record["permission_policy"]["effective_grants"] = []
    for function in record["functions"]:
        function["owner_refs"] = copy.deepcopy(owners)
        function["owner_binding_scope"] = "CONSUMER"
        function["open_requirement"]["owner_ref"] = owners["progress"]
        function["resume"]["read_refs"] = [
            owners["execution_state"], owners["execution_events"], owners["progress"],
            owners["failures"], output_ref,
        ]
        function["status"]["verification_status"] = {
            "scaffold": "NOT_TESTED", "target": "NOT_TESTED", "evidence_refs": [],
        }
        for task in function["tasks"]:
            task["owner_ref"] = owners["progress"]
            task["product_permission"]["granted"] = False
            for acceptance in task["acceptance"]:
                acceptance["target_status"] = "NOT_TESTED"
                acceptance["evidence_refs"] = []
    return record

def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--source", type=Path, required=True)
    parser.add_argument("--project-root", type=Path, required=True)
    parser.add_argument("--owner-map", type=Path, required=True)
    parser.add_argument("--output", required=True)
    args = parser.parse_args()
    output = confined(args.project_root, args.output, must_exist=False)
    require(not output.exists(), "OUTPUT_ALREADY_EXISTS")
    record = bind_existing_owners(load_json(args.source), load_json(args.owner_map),
                                  args.project_root, args.output)
    output.parent.mkdir(parents=True, exist_ok=True)
    with output.open("x", encoding="utf-8", newline="\n") as stream:
        stream.write(json.dumps(record, ensure_ascii=False, indent=2) + "\n")
    print(json.dumps({"result": "OWNERS_BOUND_ONLY", "output": str(output),
                      "resume_authorized": False, "production_authorized": False,
                      "next": "Run the consumer execution-state validator and this contract validator."}))
    return 0

if __name__ == "__main__":
    raise SystemExit(main())
````

### FILE: `roles/tests/test_business_functions.py`

```yaml
block_id: "BUSINESS-FUNCTION-OPERATING-V403:file7:v1"
operation: CREATE
provenance: AUTHORED
source: "local contract, provenance metadata, validator or owner-binding glue; no upstream method code"
license: "LicenseRef-Workspace-Owner"
sha256: "b64fa55715f6a046707abbea16ff2cc34a5c2118c2bfa7eb6f156827776ef099"
variables: []
secrets_allowed: false
```

````python
"""Focused adversarial tests of the AUTHORED contract validator and owner binder."""
import copy
import json
from pathlib import Path
import sys
import tempfile
import unittest

ROLES = Path(__file__).resolve().parents[1]
sys.path.insert(0, str(ROLES))
from validate_business_functions import validate, load_json, confined, sha256
from bind_consumer_owners import bind_existing_owners

LIBRARY = Path(sys.argv.pop(sys.argv.index("--library-root") + 1)) if "--library-root" in sys.argv else Path(r"C:/Users/NL/Desktop/Public Elite Codes")
if "--library-root" in sys.argv:
    sys.argv.remove("--library-root")

class ContractsTest(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.original = load_json(ROLES / "business-functions.v403.json")
        cls.sources = load_json(ROLES / "method-sources.v403.json")
        cls.bindings = load_json(ROLES / "library-bindings.v403.json")["bindings"]
        cls.observations = load_json(ROLES / "evidence/source-observations.json")

    def setUp(self):
        self.contract = copy.deepcopy(self.original)
        self.temp = tempfile.TemporaryDirectory(prefix="contract-fixture-", dir=ROLES / "evidence")
        self.project = Path(self.temp.name).resolve()
        self.assertTrue(self.project.is_relative_to((ROLES / "evidence").resolve()))
        for key, ref in self.contract["owner_refs"].items():
            path = self.project / ref
            path.parent.mkdir(parents=True, exist_ok=True)
            if key == "execution_state":
                path.write_text(json.dumps({"project": {"id": "library-extension-v403"}}), encoding="utf-8")
            else:
                path.write_text("SYNTHETIC CONTRACT VALIDATION FIXTURE\n", encoding="utf-8")
        (self.project / "roles/tests").mkdir(parents=True)
        (self.project / "roles/tests/test_business_functions.py").write_text("# synthetic existing test reference\n", encoding="utf-8")
        self.path = self.project / "roles/business-functions.v403.json"
        self.persist()

    def tearDown(self):
        self.temp.cleanup()

    def persist(self):
        self.path.write_text(json.dumps(self.contract, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")

    def run_validation(self, sources=None, bindings=None):
        return validate(self.contract, self.sources if sources is None else sources, self.project, LIBRARY,
                        contract_path=self.path, bindings=self.bindings if bindings is None else bindings,
                        observations=self.observations)

    def rejected(self, code):
        with self.assertRaisesRegex(ValueError, code):
            self.run_validation()

    def test_actual_eleven_contracts_and_library_bindings(self):
        result = self.run_validation()
        self.assertEqual((result["function_count"], result["task_count"], result["acceptance_count"]), (11, 22, 22))
        self.assertFalse(result["business_operations_verified"])
        self.assertFalse(result["frontend_render_verified"])

    def test_track_cannot_be_auth_role(self):
        self.contract["functions"][3]["auth_role"] = True
        self.rejected("TRACK_IS_NOT_ROLE")

    def test_grok_is_learning(self):
        self.contract["functions"][0]["function_kind"] = "BUSINESS_FUNCTION"
        self.rejected("TRACK_IS_NOT_ROLE")

    def test_global_grants_forbidden(self):
        self.contract["permission_policy"]["effective_grants"] = ["admin:*"]
        self.rejected("IMPLICIT_PRIVILEGE")

    def test_inherited_privilege_forbidden(self):
        self.contract["permission_policy"]["inherited_privilege_allowed"] = True
        self.rejected("IMPLICIT_PRIVILEGE")

    def test_unknown_privilege_field_rejected(self):
        self.contract["functions"][3]["is_admin"] = True
        self.rejected("FUNCTION_FIELDS")

    def test_product_grant_forbidden(self):
        self.contract["functions"][5]["tasks"][0]["product_permission"]["granted"] = True
        self.rejected("PRODUCT_PERMISSION_GRANT")

    def test_permission_scope_and_wildcards_rejected(self):
        for key, value in (("organization_scope", "ALL"), ("resource_scope", "*"), ("id", "admin:*")):
            with self.subTest(key=key):
                self.contract = copy.deepcopy(self.original)
                self.contract["functions"][5]["tasks"][0]["product_permission"][key] = value
                self.rejected("PRODUCT_PERMISSION_GRANT")

    def test_observed_permission_must_exist_in_pack(self):
        self.contract["functions"][8]["tasks"][0]["product_permission"]["id"] = "root:admin"
        self.rejected("PERMISSION_NOT_IN_OWNER")

    def test_missing_frontend_binding_rejected(self):
        self.contract["functions"][2]["tasks"][0]["frontend"]["route"] = ""
        self.rejected("TASK_FRONTEND_BINDING")

    def test_wrong_task_frontend_binding_rejected(self):
        self.contract["functions"][2]["tasks"][0]["frontend"]["task_id"] = "sales.fake"
        self.rejected("TASK_FRONTEND_BINDING")

    def test_frontend_effect_cannot_be_enabled(self):
        self.contract["functions"][8]["tasks"][0]["frontend"]["effect_enabled"] = True
        self.rejected("TASK_FRONTEND_BINDING")

    def test_all_declared_frontend_sections_required(self):
        self.contract["frontend_contract"]["required_sections"].remove("permissions")
        self.rejected("FRONTEND_COVERAGE")

    def test_missing_acceptance_rejected(self):
        self.contract["functions"][9]["tasks"][0]["acceptance"] = []
        self.rejected("MISSING_ACCEPTANCE")

    def test_duplicate_function_rejected(self):
        self.contract["functions"][1] = copy.deepcopy(self.contract["functions"][0])
        self.rejected("FUNCTION_SET")

    def test_duplicate_task_rejected(self):
        self.contract["functions"][0]["tasks"][1] = copy.deepcopy(self.contract["functions"][0]["tasks"][0])
        self.rejected("TASK_ID")

    def test_proven_requires_evidence(self):
        self.contract["functions"][0]["status"]["verification_status"]["scaffold"] = "PROVEN"
        self.rejected("PROVEN_WITHOUT_EVIDENCE")

    def test_proven_rejects_wrong_hash_or_scope(self):
        verification = self.contract["functions"][0]["status"]["verification_status"]
        verification.update(scaffold="PROVEN", evidence_refs=["roles/proof.json"])
        self.persist()
        proof_path = self.project / "roles/proof.json"
        for proof in [
            {"result": "PASS", "claim_scope": "FUNCTION_CONTRACT_SCAFFOLD", "contract_sha256": "0"*64, "function_ids": ["grok-bot-101"]},
            {"result": "PASS", "claim_scope": "PRODUCTION", "contract_sha256": sha256(self.path), "function_ids": ["grok-bot-101"]},
            {"result": "CONDITIONED", "claim_scope": "FUNCTION_CONTRACT_SCAFFOLD", "contract_sha256": sha256(self.path), "function_ids": ["grok-bot-101"]},
        ]:
            proof_path.write_text(json.dumps(proof), encoding="utf-8")
            self.rejected("UNBOUND_PROOF")

    def test_proven_evidence_is_narrowly_bound(self):
        verification = self.contract["functions"][0]["status"]["verification_status"]
        verification.update(scaffold="PROVEN", evidence_refs=["roles/proof.json"])
        self.persist()
        proof = {"result": "PASS", "claim_scope": "FUNCTION_CONTRACT_SCAFFOLD", "contract_sha256": sha256(self.path),
                 "function_ids": ["grok-bot-101"]}
        (self.project / "roles/proof.json").write_text(json.dumps(proof), encoding="utf-8")
        self.assertEqual(self.run_validation()["result"], "PASS")
        self.contract["functions"][0]["purpose"] += " changed"
        self.persist()
        self.rejected("UNBOUND_PROOF")

    def test_target_proven_is_out_of_scope(self):
        self.contract["functions"][0]["status"]["verification_status"]["target"] = "PROVEN"
        self.rejected("TARGET_PROOF_OUT_OF_SCOPE")

    def test_future_galaxy_content_rejected(self):
        self.contract["functions"][0]["galaxy_admission"]["source_ref"] = "imagined-event-method"
        self.rejected("FUTURE_GALAXY_CLAIM")

    def test_current_methods_do_not_wait_for_event(self):
        result = self.run_validation()
        self.assertFalse(result["galaxy_event_admitted"])
        self.assertEqual(self.contract["functions"][1]["status"]["method_admission"]["status"], "OFFICIAL_METHOD_REFERENCED")

    def test_official_code_reattribution_rejected(self):
        self.contract["provenance"]["official_method_code_claim"] = True
        self.rejected("FALSE_OFFICIAL_CODE")

    def test_unofficial_method_rejected(self):
        sources = copy.deepcopy(self.sources)
        sources["sources"][0]["url"] = "https://example.invalid/method"
        with self.assertRaisesRegex(ValueError, "UNOFFICIAL_SOURCE"):
            self.run_validation(sources=sources)

    def test_unknown_source_rejected(self):
        self.contract["functions"][0]["source_refs"].append("FUTURE")
        self.rejected("UNKNOWN_METHOD_SOURCE")

    def test_drift_in_historical_pack_rejected(self):
        bindings = copy.deepcopy(self.bindings)
        bindings["engineering"]["pack_sha256"] = "0"*64
        with self.assertRaisesRegex(ValueError, "LIBRARY_BINDING_CHANGED"):
            self.run_validation(bindings=bindings)

    def test_missing_owner_and_path_escape_rejected(self):
        for reference, error in (("missing.md", "MISSING_REFERENCE"), ("../other.md", "REFERENCE_ESCAPE"),
                                 ("C:/outside.md", "REFERENCE_SYNTAX")):
            self.contract = copy.deepcopy(self.original)
            self.contract["owner_refs"]["failures"] = reference
            self.rejected(error)

    def test_resume_must_reference_actual_owners(self):
        self.contract["functions"][0]["resume"]["read_refs"].remove("PROJECT_FAILURE_LESSONS.md")
        self.rejected("RESUME_OWNERS")

    def test_consumer_cannot_inherit_maintenance_state(self):
        with self.assertRaisesRegex(ValueError, "MAINTENANCE_STATE_IN_CONSUMER"):
            bind_existing_owners(self.contract, self.contract["owner_refs"], self.project, "roles/consumer.json")

    def test_consumer_cannot_inherit_maintenance_tasks(self):
        state_path = self.project / self.contract["owner_refs"]["execution_state"]
        state_path.write_text(json.dumps({"project": {"id": "consumer-actual"}}), encoding="utf-8")
        with self.assertRaisesRegex(ValueError, "MAINTENANCE_OWNER_IN_CONSUMER"):
            bind_existing_owners(self.contract, self.contract["owner_refs"], self.project, "roles/consumer.json")

    def test_consumer_rebinds_existing_owners_and_resets_claims(self):
        owner_map = copy.deepcopy(self.contract["owner_refs"])
        owner_map["progress"] = "specs/customer-journey/tasks.md"
        (self.project / owner_map["progress"]).parent.mkdir(parents=True)
        (self.project / owner_map["progress"]).write_text("# Existing consumer task owner\n", encoding="utf-8")
        state_path = self.project / owner_map["execution_state"]
        state_path.write_text(json.dumps({"project": {"id": "consumer-actual"}}), encoding="utf-8")
        before = {ref: (self.project / ref).read_bytes() for ref in owner_map.values()}
        self.contract["functions"][0]["status"]["verification_status"]["evidence_refs"] = ["old-proof.json"]
        rebound = bind_existing_owners(self.contract, owner_map, self.project, "roles/consumer.json")
        self.assertEqual(rebound["owner_binding"]["scope"], "CONSUMER")
        self.assertEqual(rebound["functions"][0]["status"]["verification_status"]["evidence_refs"], [])
        self.assertEqual(before, {ref: (self.project / ref).read_bytes() for ref in owner_map.values()})
        self.contract = rebound
        self.path = self.project / "roles/consumer.json"
        self.persist()
        self.assertEqual(self.run_validation()["result"], "PASS")

    def test_duplicate_json_keys_rejected(self):
        duplicate = self.project / "duplicate.json"
        duplicate.write_text('{"scope":"A","scope":"B"}', encoding="utf-8")
        with self.assertRaisesRegex(ValueError, "DUPLICATE_JSON_KEY"):
            load_json(duplicate)

    def test_nonfinite_json_rejected(self):
        invalid = self.project / "nonfinite.json"
        invalid.write_text('{"metric":NaN}', encoding="utf-8")
        with self.assertRaisesRegex(ValueError, "NONFINITE_JSON"):
            load_json(invalid)

if __name__ == "__main__":
    unittest.main(verbosity=2)
````

### FILE: `roles/evidence/source-observations.json`

```yaml
block_id: "BUSINESS-FUNCTION-OPERATING-V403:file8:v1"
operation: CREATE
provenance: AUTHORED
source: "local contract, provenance metadata, validator or owner-binding glue; no upstream method code"
license: "LicenseRef-Workspace-Owner"
sha256: "bb86045172dfdfa19a7851286c62406a7f77dc582ab0ee05bfcd1e67b2d8560e"
variables: []
secrets_allowed: false
```

````json
{
  "schema_version": "1.0.0",
  "scope": "DOCUMENTATION_RETRIEVAL_ONLY",
  "raw_documentation_distributed": false,
  "observations": [
    {
      "id": "MS-RBAC",
      "url": "https://learn.microsoft.com/en-us/azure/role-based-access-control/overview",
      "query": "site.learn.microsoft.com azure role-based-access-control overview",
      "observed_at": "2026-09-14T01:56:24.692667+00:00",
      "body_stored": false,
      "publisher_signature_verified": false,
      "status": "OBSERVED",
      "http_status": 200,
      "final_url": "https://learn.microsoft.com/en-us/azure/role-based-access-control/overview",
      "content_type": "text/html",
      "bytes": 58525,
      "response_sha256": "3ddc18157a0597901158053c80b32b3831510a4cf5099a5d3876de9a3845a6ed",
      "etag": "\"wa2CLDgk+7IuV35IjdZtStWqpQg=\"",
      "last_modified": "Sun, 13 Sep 2026 11:36:51 GMT"
    },
    {
      "id": "DORA-USER",
      "url": "https://dora.dev/capabilities/user-centric-focus/",
      "query": "site.dora.dev capabilities user-centricity",
      "observed_at": "2026-09-14T01:56:24.693428+00:00",
      "body_stored": false,
      "publisher_signature_verified": false,
      "status": "OBSERVED",
      "http_status": 200,
      "final_url": "https://dora.dev/capabilities/user-centric-focus/",
      "content_type": "text/html; charset=utf-8",
      "bytes": 16123,
      "response_sha256": "bc7fa93860500ebd0ac4a0b55a4412bfdb012cef9c1530915873f136833514f8",
      "etag": "\"57388c00434792f4263f8c1a3a47726ed61d0efccf243e4d81e03c9ed3fa7e4e\"",
      "last_modified": "Sun, 13 Sep 2026 12:39:25 GMT"
    },
    {
      "id": "DORA-BATCHES",
      "url": "https://dora.dev/capabilities/working-in-small-batches/",
      "query": "site.dora.dev capabilities working in small batches",
      "observed_at": "2026-09-14T01:56:24.693925+00:00",
      "body_stored": false,
      "publisher_signature_verified": false,
      "status": "OBSERVED",
      "http_status": 200,
      "final_url": "https://dora.dev/capabilities/working-in-small-batches/",
      "content_type": "text/html; charset=utf-8",
      "bytes": 20232,
      "response_sha256": "6b43fddc441be5fbf6f5c1e523e130b4e803d10c4d25725d7c7ea4835638af6e",
      "etag": "\"0fe23578d2798b1ac03dab5dfc94f1f159cb2f57ac7f870c71f1bb74ad45226c\"",
      "last_modified": "Sun, 13 Sep 2026 12:39:25 GMT"
    },
    {
      "id": "MS-SALES",
      "url": "https://learn.microsoft.com/en-us/dynamics365/guidance/business-processes/prospect-to-quote-overview",
      "query": "site.learn.microsoft.com dynamics365 guidance business-processes prospect-to-quote overview",
      "observed_at": "2026-09-14T01:56:24.694468+00:00",
      "body_stored": false,
      "publisher_signature_verified": false,
      "status": "OBSERVED",
      "http_status": 200,
      "final_url": "https://learn.microsoft.com/en-us/dynamics365/guidance/business-processes/prospect-to-quote-overview",
      "content_type": "text/html",
      "bytes": 58154,
      "response_sha256": "df88514da1be4809a1f54fdc6e6a577caf4cd98464659e850c6bd783ad043633",
      "etag": "\"T3poe7Tg6N9XATJoJNBoUhJBPFs=\"",
      "last_modified": "Tue, 25 Aug 2026 19:04:57 GMT"
    },
    {
      "id": "MS-CASES",
      "url": "https://learn.microsoft.com/en-us/dynamics365/guidance/business-processes/case-to-resolution-introduction",
      "query": "site.learn.microsoft.com dynamics365 guidance case to resolution business process overview",
      "observed_at": "2026-09-14T01:56:24.921659+00:00",
      "body_stored": false,
      "publisher_signature_verified": false,
      "status": "OBSERVED",
      "http_status": 200,
      "final_url": "https://learn.microsoft.com/en-us/dynamics365/guidance/business-processes/case-to-resolution-introduction",
      "content_type": "text/html",
      "bytes": 59787,
      "response_sha256": "57075bd3013bf8396fd41aac177dc2611290460811ba8df34aac6aed7bcfc922",
      "etag": "\"5f9ZhEeujbbiS9n9zfOe3Kb4jyA=\"",
      "last_modified": "Tue, 25 Aug 2026 19:04:58 GMT"
    },
    {
      "id": "MS-SERVICE",
      "url": "https://learn.microsoft.com/en-us/dynamics365/guidance/business-processes/service-to-cash-introduction",
      "query": "service to deliver introduction business processes site:learn.microsoft.com",
      "observed_at": "2026-09-14T01:56:25.208177+00:00",
      "body_stored": false,
      "publisher_signature_verified": false,
      "status": "OBSERVED",
      "http_status": 200,
      "final_url": "https://learn.microsoft.com/en-us/dynamics365/guidance/business-processes/service-to-cash-introduction",
      "content_type": "text/html",
      "bytes": 58829,
      "response_sha256": "39ad6499c9a5284a1d55992bd989873e87ad9bf8b5ef020cf2c6c01de8ebf6c1",
      "etag": "\"vLO31ghjQLp+1qpK9NlAy1R+Z88=\"",
      "last_modified": "Tue, 25 Aug 2026 19:04:57 GMT"
    },
    {
      "id": "HUBSPOT-CAMPAIGNS",
      "url": "https://knowledge.hubspot.com/campaigns/create-campaigns",
      "query": "site.knowledge.hubspot.com campaigns create campaigns",
      "observed_at": "2026-09-14T01:56:25.275407+00:00",
      "body_stored": false,
      "publisher_signature_verified": false,
      "status": "OBSERVED",
      "http_status": 200,
      "final_url": "https://knowledge.hubspot.com/campaigns/create-campaigns",
      "content_type": "text/html; charset=UTF-8",
      "bytes": 437750,
      "response_sha256": "d7ae474f1938d132eb1aedbaead72ccbbbcfbb44f3aff56ee2fff4711bd1b17f",
      "etag": "W/\"cdd6068a046227070075d1d05997fa48\"",
      "last_modified": "Mon, 14 Sep 2026 01:48:39 GMT"
    },
    {
      "id": "XAI-GROK",
      "url": "https://docs.x.ai/grok/overview",
      "query": "Grok API overview documentation site:docs.x.ai",
      "observed_at": "2026-09-14T01:56:25.409149+00:00",
      "body_stored": false,
      "publisher_signature_verified": false,
      "status": "OBSERVED",
      "http_status": 200,
      "final_url": "https://docs.x.ai/grok/overview",
      "content_type": "text/html; charset=utf-8",
      "bytes": 235385,
      "response_sha256": "e04a90678463e4fea4463412c5ed35cb0125927229888db7d82fc18c62a90f54",
      "etag": null,
      "last_modified": null
    }
  ]
}
````

## 6. Configuration surface

Explicit project root, library root, source contract, owner-map and new output/report path. No provider endpoints, keys, subscriptions, runtime identities or product settings. The read-only frontend shows every declared section/task and independently verifies rendered coverage.

## 7. Dependency bill

Python standard library only. Reuse materialize_markdown_pack.ps1 unchanged; no new parser, compositor, package or runtime acquisition. Microsoft/DORA/HubSpot/xAI sources are documentary method references with scoped short paraphrases, no copied implementation. Source-observation hashes record historical HTTP bodies that are not retained or distributed; they do not constitute reproducible snapshots or publisher signatures.

## 8. Apply order

Materialize this pack using the existing script into an absent destination. Resolve the consumer's real owners using bind_consumer_owners.py and a new output contract; do not inherit maintenance tasks/state. Validate the execution-state chain with the existing 1.3.1 kit before resume. Run validate_business_functions.py against the exact contract and library; run the focused tests. Integrate the frontend independently and retain its own UI/browser evidence.

## 9. Verification

validate_business_functions.py rejects implicit grants, unknown privilege fields, missing task/frontend/acceptance links, source promotion, future Galaxy claims, target PROVEN, unbound evidence, changed library refs and maintenance-owner inheritance. It emits only FUNCTION_CONTRACT_SCAFFOLD evidence with product and frontend execution false. tests/test_business_functions.py contains 33 focused cases including negative inheritance, stale hash and owner-preservation checks. A local PASS does not prove business, provider, native-agent or human acceptance.

## 10. Reconstruction evidence

Maintenance V403 receipts live outside this pack under roles/evidence in the durable extension sandbox: validation-standard.json, tests-standard.log and standard-reconstruction-result.json. Existing root progress/failure owners preserve test status and earlier tooling failures. Two independent standard materializations must match all eight files exactly; re-run the 33 cases from reconstructed source. Preserve the sealed V402 state337 and ZIPs. No production authorization.
