# Authority map automation — investigación de ajuste V329

Fecha: 2026-09-08. Clasificación: mantenimiento de biblioteca / discovery T2801.
Continuidad: checkpoint88 validado antes de investigar; no cambio de código,
dependencia, source lock, respuesta de readiness ni efecto externo.

## Resultado y límite

CAPABILITY-GAP-RESOLUTION-GATE0.1.0 materializado, cuatro archivos; su verificador
terminó PASS con3positivos/6negativos. El expediente real de cuatro búsquedas y
tres candidatos termina NO_ADMISSIBLE_SOURCE, implementation_ready=false.
El exit0 valida coherencia del expediente y no aprueba implementación ni prueba
inexistencia universal de soluciones. El alcance fue el registro de métodos
portables de la biblioteca y su catálogo, contrastados con fuentes oficiales.

No se instaló ningún método ni se adquirió/ejecutó código de los candidatos.
Las notas locales son AUTHORED; la clasificación OFFICIAL_ARCHITECTURE/ADAPTED
del expediente describe una evaluación adaptada del método documentado.
Los hashes fijan las notas preservadas, NO los bytes de código/licencia/archives
upstream. Las revisiones identifican la documentación consultada. No hay un
OFFICIAL_CODE ni un REUSABLE_PACK admitido como resultado de esta investigación.

## Evaluación

| Candidato | Evidencia acotada | Decisión para este claim |
|---|---|---|
| GitHub Spec Kit1.0.1 | README de agent-context y funciones _resolve_plan_path/_build_section observadas en la revisión fijada | Rechazado: referencia a un plan no demuestra selección de manuales |
| AWS AI-DLC1.0.1 | workflow fijado: carga progresiva, detección y reglas de ejecución; docs actuales separadas | Rechazado: sin contrato local de selección; no adoptar registro crudo/aprobaciones del método |
| Conductor99ba10e | README/revisión observada, artefactos de contexto y tracks; release0.4.1 separada | Rechazado: no demuestra selección por decisión/requisito de Elite |

G0 PASS sólo identidad documental. G1 y G4–G7 NOT_APPLICABLE a incorporación
porque no se incorpora el candidato; no fueron superados por leer una licencia,
un README o una página de advisories. G2/G3/G8 FAIL para el claim integral.
La decisión de ajuste es una inferencia explícita de las interfaces documentadas,
no una ejecución de suites upstream. Cada nota conserva URLs oficiales, revisión,
licencia observada, razonamiento y límites; fuentes no primarias se descartaron.

El registro también contempla Skills y archivos de instrucciones como formatos,
sin selector ejecutable propio; Kiro es un producto con términos, no el pack
portable faltante. La propuesta comunitaria Spec Kit issue2989 no equivale a
código oficial incorporado. No se confunde ausencia en una búsqueda con rechazo.

## Delta de contrato y autoridades

Antes: bootstrap §5 enumeraba ocho secciones, y codex.authority_docs conservaba
cuatro paths. validate_project.py comprueba extensión/existencia de esos paths;
el mapa manual contiene decisiones, pero no existe generador de su justificación.
Después: bootstrap §5.1 precisa seis criterios verificables: NEW, EXISTING,
negativos, determinismo/presupuesto, publicación/recuperación y pertinencia.
No se cambió el significado de READY_TO_BUILD ni se exigió aprobación ceremonial.

La selección razonada y la validación/renderización mecánica tienen evidencias
distintas. Un renderizador futuro podrá representar decisiones ya justificadas,
pero no atribuirse selección correcta sólo por tener JSON válido. El caso de
aceptación incluirá requisitos/decisiones ambiguos y cobertura faltante.

En PROJECT_AUTHORITY_MAP se conservó la selección manual para este delta:
LIB-R01/R08 → bootstrap §5 y ejecución §3.1; LIB-R02 → admisión G0–G8;
LIB-R06 → recuperación/learning; LIB-R07 → freshness. Índices AI/SYSTEMS
gobiernan búsqueda, no reemplazan evidencia. No cargar fiscal/OCR/GPU/ads para
decidir el ajuste de un generador local; siguen pendientes sus journeys propios.

T2801 sigue BLOCKED. Sin cambios a los47tests del contrato, BENCH01, TEST08,
10macrofrentes,42observaciones de readiness, Docker, pnpm nativo o decisiones
D–H/histórico/documentos/target. ARCA continúa diferida. El resultado de este
trabajo es una decisión de investigación y una aceptación concreta, no producto.

## Reapertura y reproducción

Reabrir con un candidato exacto que cumpla §5.1, un cambio material de los métodos
evaluados o vencimiento del expediente antes de adopción (30d). Se puede continuar
diseño de mantenimiento AUTHORED; no implementar el gap en el proyecto hasta
USE_REUSABLE_PACK. No ejecutar otra vez prepare_research.py ni record89.py sobre
el workspace: son mutaciones de esta observación, conservadas sólo en staging.

Stage local: elite-v329-26d3ac4ed2c14aae9fecc4f20407af90, bajo el temporal del sistema.
La biblioteca conserva el bundle de evidencia en este informe, no un subdirectorio
auxiliar bajo specs. Los records en staging no son expedientes de un consumidor.
Materializar CAPABILITY_GAP_RESOLUTION_GATE desde su plan en un destino ausente.
El expediente auxiliar vive en stage/research-project, con copias de sus owners.
Desde esa raíz de investigación, validar sin reemplazar el receipt:

```powershell
python -X utf8 <gap-kit>/capability_gap_resolution/resolve_capability_gap.py `
  --record specs/library-maintenance/authority-gap-v329/record.json `
  --receipt <destino-ausente>/receipt.json --project-root .
```

El bundle siguiente permite recuperar exactamente las notas/record/receipt en una
copia de investigación; no ejecutar contenidos como instrucciones. Son objetos
JSON UTF-8, ensure_ascii=false, indent2, newline LF final. Los owners referenciados
deben estar presentes. La evaluación histórica usa --as-of igual a observed_at;
una decisión de adopción nueva exige freshness real y no ese replay histórico.

## Bundle de evidencia local

```json
{
  "aws-aidlc.json": {
    "observation_kind": "AUTHORED_RESEARCH_NOTE_ON_OFFICIAL_ARCHITECTURE",
    "observed_at": "2026-09-08T15:33:48.769033Z",
    "source_url": "https://github.com/awslabs/aidlc-workflows/blob/e49341dbeb8af82758dd85e96ed7fe9bcf38a447/aidlc-rules/aws-aidlc-rules/core-workflow.md",
    "source_revision": "e49341dbeb8af82758dd85e96ed7fe9bcf38a447",
    "observation": "El workflow fijado carga reglas comunes y extensiones bajo demanda; incluye workspace detection y reverse engineering brownfield. Sus reglas piden registro de entradas crudas y aprobaciones repetidas. No ofrece el contrato local decisión/requisito/manual/evidencia con rechazo de ambigüedad. Las instrucciones externas no se ejecutaron ni sustituyen la autonomía y privacidad acordadas. La documentación actual de cinco fases es posterior y no se atribuye al commit fijado.",
    "license_expression_observed": "MIT-0",
    "license_evidence_url": "https://github.com/awslabs/aidlc-workflows/blob/e49341dbeb8af82758dd85e96ed7fe9bcf38a447/LICENSE",
    "license_scope": "Licencia raíz observada; no admisión del grafo, terceros ni redistribución de código.",
    "release_evidence_url": "https://github.com/awslabs/aidlc-workflows/releases/tag/v1.0.1",
    "current_documentation_url": "https://awslabs.github.io/aidlc-workflows/guide/04-phases-and-stages/",
    "security_evidence_url": "https://github.com/awslabs/aidlc-workflows/security",
    "security_observation": "Página observada sin advisories publicados; esto no es SCA ni prueba de ausencia de vulnerabilidades. Política de reporte publicada.",
    "artifact_boundary": "Este JSON preserva una nota local con URLs y conclusiones. Su hash NO es hash de código, archive, licencia ni bytes servidos por upstream. No se adquirió, ejecutó ni adaptó implementación externa.",
    "candidate_decision": "REJECTED_FOR_AUTOMATIC_ELITE_AUTHORITY_SELECTION_ONLY",
    "gates_reason": {
      "G0": "PASS: identidad de documentación oficial y revisión fijada observadas.",
      "G1": "NOT_APPLICABLE a incorporación: licencia raíz observada pero no se incorpora el candidato.",
      "G2": "FAIL: no satisface selección explicada por decisión/requisito conforme al bootstrap §5.",
      "G3": "FAIL: no se demuestra generador de ese contrato local.",
      "G4": "NOT_APPLICABLE: no se ejecutó suite del candidato rechazado por ajuste.",
      "G5": "NOT_APPLICABLE: no hubo admisión de seguridad; revisión de página no es SCA.",
      "G6": "NOT_APPLICABLE: sin benchmark de implementación para este claim.",
      "G7": "NOT_APPLICABLE: sin operación del target para este claim.",
      "G8": "FAIL: no hay pack compatible admitido para generación integral."
    }
  },
  "conductor.json": {
    "observation_kind": "AUTHORED_RESEARCH_NOTE_ON_OFFICIAL_ARCHITECTURE",
    "observed_at": "2026-09-08T15:33:48.769033Z",
    "source_url": "https://github.com/gemini-cli-extensions/conductor/blob/99ba10e1a11130fc159f681b7ba8803489239cbf/README.md",
    "source_revision": "99ba10e1a11130fc159f681b7ba8803489239cbf",
    "observation": "El README fijado describe generación de product, product-guidelines, tech-stack, workflow y tracks; soporta proyectos nuevos y existentes y advierte mayor consumo de contexto. Es metodología de contexto/spec/plan, sin el contrato de selección de autoridades Elite. El commit observado es una revisión de investigación, no una release admitida; la página releases lista conductor-v0.4.1 en otra revisión. No instalar ni inferir ahorro de tokens.",
    "license_expression_observed": "Apache-2.0",
    "license_evidence_url": "https://github.com/gemini-cli-extensions/conductor/blob/99ba10e1a11130fc159f681b7ba8803489239cbf/LICENSE",
    "license_scope": "Licencia raíz observada; no admisión del grafo, terceros ni redistribución de código.",
    "release_evidence_url": "https://github.com/gemini-cli-extensions/conductor/releases",
    "current_documentation_url": "https://github.com/gemini-cli-extensions/conductor/blob/99ba10e1a11130fc159f681b7ba8803489239cbf/README.md",
    "security_evidence_url": "https://github.com/gemini-cli-extensions/conductor/security",
    "security_observation": "Página observada sin advisories publicados; esto no es SCA ni prueba de ausencia de vulnerabilidades. Sin SECURITY.md detectado.",
    "artifact_boundary": "Este JSON preserva una nota local con URLs y conclusiones. Su hash NO es hash de código, archive, licencia ni bytes servidos por upstream. No se adquirió, ejecutó ni adaptó implementación externa.",
    "candidate_decision": "REJECTED_FOR_AUTOMATIC_ELITE_AUTHORITY_SELECTION_ONLY",
    "gates_reason": {
      "G0": "PASS: identidad de documentación oficial y revisión fijada observadas.",
      "G1": "NOT_APPLICABLE a incorporación: licencia raíz observada pero no se incorpora el candidato.",
      "G2": "FAIL: no satisface selección explicada por decisión/requisito conforme al bootstrap §5.",
      "G3": "FAIL: no se demuestra generador de ese contrato local.",
      "G4": "NOT_APPLICABLE: no se ejecutó suite del candidato rechazado por ajuste.",
      "G5": "NOT_APPLICABLE: no hubo admisión de seguridad; revisión de página no es SCA.",
      "G6": "NOT_APPLICABLE: sin benchmark de implementación para este claim.",
      "G7": "NOT_APPLICABLE: sin operación del target para este claim.",
      "G8": "FAIL: no hay pack compatible admitido para generación integral."
    }
  },
  "receipt.json": {
    "authority_count": 7,
    "candidate_count": 3,
    "capability_id": "ELITE-AUTHORITY-MAP-AUTOMATION",
    "decision_state": "NO_ADMISSIBLE_SOURCE",
    "gate": "CAPABILITY-GAP-RESOLUTION-GATE",
    "implementation_ready": false,
    "observed_at": "2026-09-08T15:33:48.769033Z",
    "official_domain_count": 3,
    "outcome": "BLOCKED",
    "record_sha256": "44204cf196b75760ebf3301fc14d375e4fb0ce0be36cfed4efec09fcb00afad4",
    "requirement_ref": "REUSABLE_CODE_READINESS_ROADMAP.md",
    "schema_version": 1,
    "search_count": 4,
    "selected_candidate_id": null
  },
  "record.json": {
    "schema_version": 1,
    "capability_id": "ELITE-AUTHORITY-MAP-AUTOMATION",
    "requirement_ref": "REUSABLE_CODE_READINESS_ROADMAP.md",
    "observed_at": "2026-09-08T15:33:48.769033Z",
    "max_age_days": 30,
    "authority_refs": [
      "CODEX_ELITE_PROJECT_BOOTSTRAP.md",
      "ENGINEERING_EXECUTION_PLAYBOOK.md",
      "PUBLIC_CODE_ARCHITECTURE_ADMISSION_STANDARD.md",
      "AI_ENGINEERING_MASTER_MAP.md",
      "SYSTEMS_ENGINEERING_MASTER_MAP.md",
      "markdown_system/ELITE_PUBLIC_AGENT_METHODS.md",
      "markdown_system/CAPABILITY_CATALOG.md"
    ],
    "source_searches": [
      {
        "search_id": "github-method",
        "query": "site:github.com/github/spec-kit update-agent-context plan.md constitution context",
        "official_domains": [
          "github.com",
          "github.github.com"
        ],
        "executed_at": "2026-09-08T15:33:48.769033Z",
        "evidence_refs": [
          "specs/library-maintenance/authority-gap-v329/spec-kit.json"
        ]
      },
      {
        "search_id": "aws-method",
        "query": "site:github.com/awslabs/aidlc-workflows brownfield application design reverse engineering",
        "official_domains": [
          "github.com",
          "awslabs.github.io"
        ],
        "executed_at": "2026-09-08T15:33:48.769033Z",
        "evidence_refs": [
          "specs/library-maintenance/authority-gap-v329/aws-aidlc.json"
        ]
      },
      {
        "search_id": "conductor-method",
        "query": "site:github.com/gemini-cli-extensions/conductor context Apache",
        "official_domains": [
          "github.com"
        ],
        "executed_at": "2026-09-08T15:33:48.769033Z",
        "evidence_refs": [
          "specs/library-maintenance/authority-gap-v329/conductor.json"
        ]
      },
      {
        "search_id": "authority-specific",
        "query": "site:github.com/github/spec-kit authority map",
        "official_domains": [
          "github.com"
        ],
        "executed_at": "2026-09-08T15:33:48.769033Z",
        "evidence_refs": [
          "specs/library-maintenance/authority-gap-v329/spec-kit.json"
        ]
      }
    ],
    "candidates": [
      {
        "candidate_id": "spec-kit",
        "origin_kind": "OFFICIAL_ARCHITECTURE",
        "provenance": "ADAPTED",
        "source_url": "https://github.com/github/spec-kit/blob/9118ed15a0ba65053469a94c560ea5d233f75884/extensions/agent-context/README.md",
        "claim": "Evaluar el método documentado como candidato a generación explicada de PROJECT_AUTHORITY_MAP.",
        "non_claims": [
          "No código incorporado ni ejecutado",
          "No admisión productiva",
          "Hash sólo de nota local adaptada, no bytes upstream"
        ],
        "status": "REJECTED",
        "immutable_revision": "9118ed15a0ba65053469a94c560ea5d233f75884",
        "artifact_sha256": "67ef4f14bd2bfd63bbd4cf9c9f41a4a8306a9086e66d1d36908e971a8cdbe7e8",
        "artifact_ref": "specs/library-maintenance/authority-gap-v329/spec-kit.json",
        "license_expression": "MIT",
        "license_evidence_ref": "specs/library-maintenance/authority-gap-v329/spec-kit.json",
        "gates": {
          "G0": "PASS",
          "G1": "NOT_APPLICABLE",
          "G2": "FAIL",
          "G3": "FAIL",
          "G4": "NOT_APPLICABLE",
          "G5": "NOT_APPLICABLE",
          "G6": "NOT_APPLICABLE",
          "G7": "NOT_APPLICABLE",
          "G8": "FAIL"
        },
        "rejection_reasons": [
          "Falta el contrato de selección explicado por decisión/requisito y autoridades Elite; no hay REUSABLE_PACK compatible."
        ]
      },
      {
        "candidate_id": "aws-aidlc",
        "origin_kind": "OFFICIAL_ARCHITECTURE",
        "provenance": "ADAPTED",
        "source_url": "https://github.com/awslabs/aidlc-workflows/blob/e49341dbeb8af82758dd85e96ed7fe9bcf38a447/aidlc-rules/aws-aidlc-rules/core-workflow.md",
        "claim": "Evaluar el método documentado como candidato a generación explicada de PROJECT_AUTHORITY_MAP.",
        "non_claims": [
          "No código incorporado ni ejecutado",
          "No admisión productiva",
          "Hash sólo de nota local adaptada, no bytes upstream"
        ],
        "status": "REJECTED",
        "immutable_revision": "e49341dbeb8af82758dd85e96ed7fe9bcf38a447",
        "artifact_sha256": "079539da8a82d324cc284f5b3aac26ae873a19f0bc708fc3302d4c657fb22d31",
        "artifact_ref": "specs/library-maintenance/authority-gap-v329/aws-aidlc.json",
        "license_expression": "MIT-0",
        "license_evidence_ref": "specs/library-maintenance/authority-gap-v329/aws-aidlc.json",
        "gates": {
          "G0": "PASS",
          "G1": "NOT_APPLICABLE",
          "G2": "FAIL",
          "G3": "FAIL",
          "G4": "NOT_APPLICABLE",
          "G5": "NOT_APPLICABLE",
          "G6": "NOT_APPLICABLE",
          "G7": "NOT_APPLICABLE",
          "G8": "FAIL"
        },
        "rejection_reasons": [
          "Falta el contrato de selección explicado por decisión/requisito y autoridades Elite; no hay REUSABLE_PACK compatible."
        ]
      },
      {
        "candidate_id": "conductor",
        "origin_kind": "OFFICIAL_ARCHITECTURE",
        "provenance": "ADAPTED",
        "source_url": "https://github.com/gemini-cli-extensions/conductor/blob/99ba10e1a11130fc159f681b7ba8803489239cbf/README.md",
        "claim": "Evaluar el método documentado como candidato a generación explicada de PROJECT_AUTHORITY_MAP.",
        "non_claims": [
          "No código incorporado ni ejecutado",
          "No admisión productiva",
          "Hash sólo de nota local adaptada, no bytes upstream"
        ],
        "status": "REJECTED",
        "immutable_revision": "99ba10e1a11130fc159f681b7ba8803489239cbf",
        "artifact_sha256": "1c35c55b7f78d105a634b45b7ad57a80a1c6311359b0cd37a990359045e8d7e9",
        "artifact_ref": "specs/library-maintenance/authority-gap-v329/conductor.json",
        "license_expression": "Apache-2.0",
        "license_evidence_ref": "specs/library-maintenance/authority-gap-v329/conductor.json",
        "gates": {
          "G0": "PASS",
          "G1": "NOT_APPLICABLE",
          "G2": "FAIL",
          "G3": "FAIL",
          "G4": "NOT_APPLICABLE",
          "G5": "NOT_APPLICABLE",
          "G6": "NOT_APPLICABLE",
          "G7": "NOT_APPLICABLE",
          "G8": "FAIL"
        },
        "rejection_reasons": [
          "Falta el contrato de selección explicado por decisión/requisito y autoridades Elite; no hay REUSABLE_PACK compatible."
        ]
      }
    ],
    "decision": {
      "state": "NO_ADMISSIBLE_SOURCE",
      "selected_candidate_id": null,
      "reason": "Búsqueda acotada a métodos públicos del registro interno. Tres candidatos oficiales no satisfacen el claim integral; el tooling local valida paths y no selección semántica. No demuestra inexistencia universal.",
      "canonical_updates": [
        "PROJECT_AUTHORITY_MAP.md",
        "CODEX_ELITE_PROJECT_BOOTSTRAP.md",
        "markdown_system/CAPABILITY_CATALOG.md",
        "markdown_system/FRANCHISE_GAP_MAP.md"
      ],
      "blockers": [
        "No existe pack compatible admitido en el catálogo evaluado para este claim; renderizar paths no lo cubre.",
        "T2801 conserva readiness y assurance pendientes; no comenzar implementación de este gap hasta USE_REUSABLE_PACK."
      ],
      "exhaustion": {
        "searched_authority_registry": true,
        "searched_official_docs": true,
        "searched_official_repositories": true,
        "checked_release_security": true,
        "scope_reason": "Registro ELITE_PUBLIC_AGENT_METHODS y CAPABILITY_CATALOG: GitHub Spec Kit, AWS AI-DLC y Conductor son los métodos portables candidatos. Skills y archivos de instrucciones son formatos, no selectores; Kiro es producto con términos y no generador portable admitido. Issue2989 blueprint es propuesta comunitaria no incorporada por reputación. Ninguna búsqueda fallida cuenta como ausencia. No se investigó todo el software público.",
        "next_review_trigger": "Reabrir ante pack/candidato exacto que implemente el contrato §5.1, revisión material de los tres métodos, o al vencer30d antes de adopción. Puede continuar diseño de mantenimiento AUTHORED, sin atribución de código al proveedor ni habilitación automática."
      }
    }
  },
  "spec-kit.json": {
    "observation_kind": "AUTHORED_RESEARCH_NOTE_ON_OFFICIAL_ARCHITECTURE",
    "observed_at": "2026-09-08T15:33:48.769033Z",
    "source_url": "https://github.com/github/spec-kit/blob/9118ed15a0ba65053469a94c560ea5d233f75884/extensions/agent-context/README.md",
    "source_revision": "9118ed15a0ba65053469a94c560ea5d233f75884",
    "observation": "La extensión opt-in mantiene bloques de contexto. El script Python _resolve_plan_path selecciona feature.json o el plan más reciente por mtime; _build_section crea una referencia a ese plan. No selecciona manuales Elite ni justifica requisitos, fronteras y exclusiones. Esta conclusión de ajuste deriva del README y las funciones observadas; no es un test ejecutado del upstream.",
    "license_expression_observed": "MIT",
    "license_evidence_url": "https://github.com/github/spec-kit/blob/9118ed15a0ba65053469a94c560ea5d233f75884/LICENSE",
    "license_scope": "Licencia raíz observada; no admisión del grafo, terceros ni redistribución de código.",
    "release_evidence_url": "https://github.com/github/spec-kit/releases/tag/v1.0.1",
    "current_documentation_url": "https://github.github.com/spec-kit/reference/agentic-sdd.html",
    "security_evidence_url": "https://github.com/github/spec-kit/security",
    "security_observation": "Página observada sin advisories publicados; esto no es SCA ni prueba de ausencia de vulnerabilidades. Política de reporte publicada.",
    "artifact_boundary": "Este JSON preserva una nota local con URLs y conclusiones. Su hash NO es hash de código, archive, licencia ni bytes servidos por upstream. No se adquirió, ejecutó ni adaptó implementación externa.",
    "candidate_decision": "REJECTED_FOR_AUTOMATIC_ELITE_AUTHORITY_SELECTION_ONLY",
    "gates_reason": {
      "G0": "PASS: identidad de documentación oficial y revisión fijada observadas.",
      "G1": "NOT_APPLICABLE a incorporación: licencia raíz observada pero no se incorpora el candidato.",
      "G2": "FAIL: no satisface selección explicada por decisión/requisito conforme al bootstrap §5.",
      "G3": "FAIL: no se demuestra generador de ese contrato local.",
      "G4": "NOT_APPLICABLE: no se ejecutó suite del candidato rechazado por ajuste.",
      "G5": "NOT_APPLICABLE: no hubo admisión de seguridad; revisión de página no es SCA.",
      "G6": "NOT_APPLICABLE: sin benchmark de implementación para este claim.",
      "G7": "NOT_APPLICABLE: sin operación del target para este claim.",
      "G8": "FAIL: no hay pack compatible admitido para generación integral."
    },
    "code_inspection_url": "https://github.com/github/spec-kit/blob/9118ed15a0ba65053469a94c560ea5d233f75884/extensions/agent-context/scripts/python/update_agent_context.py"
  }
}
```

## Corrección de empaquetado y revalidación

FAIL540: VERIFY_LIBRARY rechazó la carpeta authority-gap-v329 bajo specs,
porque no pertenece al layout local cerrado. Se trasladó únicamente esa carpeta
propia al stage/research-project con paths absolutos comprobados. Cuatro archivos
del manifest preservan hashes exactos y el quinto es el receipt original; el
receipt regenerado contra copias de los owners resulta idéntico al original.
El bundle JSON anterior conserva también las notas y los dos artefactos del gate.
No se modificaron VERIFY_LIBRARY, sus selectores ni la allowlist de distribución.
Se retiran del cursor los enlaces a auxiliares fuera de la raíz; el informe
portable es la evidencia durable y reproduce el expediente. Eventos89 intactos.

## Verificación integrada observada

VERIFY_LIBRARY_PASS tras la corrección de layout e inventario:161packs,
1453archivos materializables,765Markdown y52perfiles. El verificador, los
selectores y los implementation packs no cambian; no hubo adopción de runtime.
La suite completa de runtimes no se repite para este delta documental. El
Preflight87 de V328 conserva su fecha/alcance,151pasos PASS y Docker BLOCKED.

SHA-256 verify91.log: `0d7937e3a0f34e382dcdb1f3168477f22a419ddf5c89627425ed72c5a82cffd5`.

FAIL540 cerrado por preservación exacta del bundle/receipts y gate portable;
verify89/90 fallidos quedan preservados. Registro final añade sólo evidencia,
lección y continuidad, sin código, pins ni cambios de criterios posteriores al
gate. Checkpoint92 conserva91eventos previos y revalida contrato/cadena.

Siguiente verificación independiente: T2804, coordinación entre pestañas de
operaciones de devolución ya existentes. RETURN_OPERATIONS_RECOVERY_V312.md
explicita ese límite; no repetir sus casos verdes ni inventar nuevas reglas
comerciales. Readiness, generación automática, integración final y decisiones
D–H/corpus/target siguen pendientes.

Cierre de cursor: primer intento de checkpoint92 rechazó portal-owner/v312 duplicados entre must-read y reuse. FAIL541 conserva el rechazo sin append; se movieron de categoría, checkpoint92/resume92 PASS con110evidencias y92eventos. Checkpoint93 registra sólo esta corrección y su lección, sin delta de código/criterios respecto del gate portable observado.
