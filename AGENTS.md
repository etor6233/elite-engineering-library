# Elite Engineering Library — agent index

Router compacto. No cargues el corpus entero ni leas `AGENT_SYSTEM_START.md` de punta a punta por defecto.

## Misión

Acelerar proyectos reales sin rebajar calidad, seguridad, rendimiento, operabilidad ni cumplimiento de licencias. La fama o popularidad no sustituyen admisión auditada.

## Entrada por agente

| Agente | Arranque |
|---|---|
| **Codex / OpenAI** | Este índice + `.agents/skills/elite-engineering-library/SKILL.md` |
| **Claude Code** | `CLAUDE.md` → este índice |
| **Grok Build** | `GROK_AGENT_ENTRY.md` → este índice + `.grok/skills/elite-engineering-library/SKILL.md` |

Si la biblioteca no es la raíz del proyecto: `INSTALL_AGENT_BRIDGE.ps1` (`-Agent Codex`, `Claude`, `Grok`, `Both` o `All`). Sin Git, abre el agente en la raíz exacta del proyecto.

## Divulgación progresiva (obligatoria)

1. Clasifica la tarea: mantenimiento de biblioteca · discovery · diseño · implementación · revisión · incidente.
2. **Busca primero** en la biblioteca; si falta o está bloqueado, escala — no inventes fuentes, packs ni estados de admisión.
3. Carga **un** mapa de autoridad y **un** pack por claim (`markdown_system/PACK_PER_CLAIM_INDEX.md`).
4. Abre secciones de gates sólo cuando el claim las exija.

**No cargar por defecto:** `AGENT_SYSTEM_START.md` entero, `PROJECT_PACK_PLAN.md`, `markdown_system/FRANCHISE_COMPLETE_PACK_PLAN.md`, `reconstruction_evidence/`, ni todos los `implementation_packs/`.

## Índice de autoridades (carga bajo demanda)

| Necesidad | Archivo |
|---|---|
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
