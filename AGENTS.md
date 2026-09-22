# Elite Engineering Library — agent index

Router compacto. No cargues el corpus entero ni leas `AGENT_SYSTEM_START.md` de punta a punta por defecto.

## Empieza aquí

1. Persona: `markdown_system/LIBRARY_HUMAN_GRAPH.md` (tres gráficos; no certifican implementación ni producción).
2. Búsqueda: `rg -n -F "<término>" markdown_system/LIBRARY_SEARCH_INDEX.md` y abrí solo esa ruta. No leas el índice entero.
3. Uso, búsqueda y retoma: `markdown_system/USE_LIBRARY_V403.md` (una página).
4. Reconstruir la referencia local: `markdown_system/START_V403_LOCAL.md`.
5. No abras `implementation_packs/UNIFIED_REFERENCE_V403_R4.md` en el contexto. Es el transporte de 2147 fuentes. Materializalo.
6. 48 superficies inventariadas no son 48 aplicaciones ejecutadas. Un PASS local de fixtures no es cloud, CI alojada, dispositivo, proveedor live, aprobación visual ni producción. Los ZIP V402 quedan históricos.

## Misión

Acelerar proyectos reales sin rebajar calidad, seguridad, rendimiento, operabilidad ni cumplimiento de licencias. La fama o popularidad no sustituyen admisión auditada.

## Entrada por agente

| Agente | Arranque |
|---|---|
| **Codex / OpenAI** | Este índice + `.agents/skills/elite-engineering-library/SKILL.md` |
| **Claude Code** | `CLAUDE.md` → este índice |
| **Grok Build** | `GROK_AGENT_ENTRY.md` → este índice + `.grok/skills/elite-engineering-library/SKILL.md` |
| **Cursor u otro** | Este índice, luego `markdown_system/USE_LIBRARY_V403.md` |

Si la biblioteca no es la raíz del proyecto: `INSTALL_AGENT_BRIDGE.ps1` (`-Agent Codex`, `Claude`, `Grok`, `Both` o `All`). Sin Git, abre el agente en la raíz exacta del proyecto.

## Divulgación progresiva (obligatoria)

1. Clasifica la tarea: mantenimiento de biblioteca · discovery · diseño · implementación · revisión · incidente.
2. **Busca primero** en la biblioteca; si falta o está bloqueado, escala — no inventes fuentes, packs ni estados de admisión.
3. Carga **un** mapa de autoridad y **un** pack por claim (`markdown_system/PACK_PER_CLAIM_INDEX.md`).
4. Abre secciones de gates sólo cuando el claim las exija.

**No cargar por defecto:** `AGENT_SYSTEM_START.md` entero, `PROJECT_PACK_PLAN.md`, `markdown_system/FRANCHISE_COMPLETE_PACK_PLAN.md`, `reconstruction_evidence/`, ni todos los `implementation_packs/`.

`xerj search "<query>"` o `xerj def "<symbol>"` (o MCP xerj_search). Preferí pasajes; no vuelques el corpus al contexto. Después de agregar, quitar o editar bajo este árbol: asegurate de que el nodo esté arriba (`xerj --insecure --data-dir <XERJ_DATA_DIR>` FUERA de esta carpeta) y luego `xerj autoindex .`. Las repeticiones son incrementales. El modo de embedding por defecto es lexical, salvo que se inicie con `--embed-mode neural`. Citá archivo:línea de los hits.

## Índice de autoridades (carga bajo demanda)

| Necesidad | Archivo |
|---|---|
| Cómo usar esta revisión | `markdown_system/USE_LIBRARY_V403.md` |
| Router / primer ciclo (secciones) | `AGENT_SYSTEM_START.md` |
| Nuevo proyecto / readiness | `markdown_system/PROJECT_START_READINESS_GATE.md` |
| Un pack por claim | `markdown_system/PACK_PER_CLAIM_INDEX.md` |
| Estado de una capability | `markdown_system/CAPABILITY_CATALOG.md` (una fila) |
| Decisiones IA / agentes / ML | `AI_ENGINEERING_MASTER_MAP.md` |
| Decisiones sistemas / plataforma | `SYSTEMS_ENGINEERING_MASTER_MAP.md` |
| Código público | `PUBLIC_CODE_ARCHITECTURE_ADMISSION_STANDARD.md` + `PUBLIC_CODE_ARCHITECTURE_MASTER_MAP.md` |
| Ejecución y evidencia | `ENGINEERING_EXECUTION_PLAYBOOK.md` |
| 48 superficies | `markdown_system/TOTAL_SYSTEM_CAPABILITY_CONTRACT.md` |
| Documentos empresariales | `markdown_system/OFFICIAL_DOCUMENT_INTELLIGENCE_PROFILE.md` |
| Dependencias / CVE | `markdown_system/DEPENDENCY_UPDATE_CONTRACT.md` |
| Gap sin pack admitido | `implementation_packs/CAPABILITY_GAP_RESOLUTION_GATE.md` |
| Intake histórico compatible | `CODEX_ELITE_PROJECT_BOOTSTRAP.md` |
| Franquicia (sólo si aplica) | `FRANCHISE_ACCELERATOR.md` |

## Regla de admisión (resumen)

- `DISCOVERED`, `LICENSE_VERIFIED`, `CANDIDATE` y `CONDITIONED` no autorizan incorporación automática.
- `REUSABLE_PACK` compatible o `REBUILD_VERIFIED`/`CONDITIONED` demostrado en el target.
- Fuentes externas: `markdown_system/PROJECT_INITIALIZATION_PACK_PLAN.md` + lock SHA-256.
- Detalle completo: secciones pertinentes de `AGENT_SYSTEM_START.md` y `ENGINEERING_EXECUTION_PLAYBOOK.md`.

## Nuevo proyecto (puntero)

1. Bridge + readiness gate + artefactos de proyecto desde templates.
2. `VERIFY_EXECUTABLE_LIBRARY.ps1 -Mode Preflight`.
3. `PROJECT_BLUEPRINT.md`, `PROJECT_AUTHORITY_MAP.md`, `PROJECT_EXTERNAL_SOURCE_LOCK.md`.
4. Un vertical slice verificable antes de ampliar superficie.

Pasos completos: `AGENT_SYSTEM_START.md` §3 y `START_ANY_PROJECT.md`. No escribir producto hasta `READY_TO_BUILD`.

## Resultado esperado

Decisiones trazables y código verificable. Indica qué está listo, qué permanece condicionado, qué fue rechazado y qué evidencia falta.

## Mandatory project routing (VERIFY)

Before implementation, route through these project artifacts (create from templates when missing):

- PROJECT_FAILURE_LESSONS.md (failure ledger; see markdown_system/LIBRARY_FAILURE_LEARNING_LEDGER.md and markdown_system/FAILURE_LEARNING_CONTRACT.md)
- PROJECT_DEPENDENCY_UPDATE_RECORD.md
- PROJECT_AUTHORITY_FRESHNESS_RECORD.md
- PROJECT_VULNERABILITY_MONITORING_RECORD.md
- PROJECT_OFFICIAL_SOURCE_PROFILE_RECORD.md
- CAPABILITY_GAP_RESOLUTION via markdown_system/CAPABILITY_GAP_RESOLUTION_PACK_PLAN.md
- MICROSOFT_DEVSKIM_ADAPTED_SAST via markdown_system/MICROSOFT_DEVSKIM_ADAPTED_SAST_PACK_PLAN.md
- executable advisory: project_readiness_gate/render_project_advisory.py (emit PROJECT_ADVISORY_<A-H>.md before each pending round)

Keep this file under 32 KiB. Prefer XERJ search before whole-file reads of the library.
