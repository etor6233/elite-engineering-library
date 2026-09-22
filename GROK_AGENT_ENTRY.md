# Grok agent entry — Elite Engineering Library

Router breve para Grok Build. El corpus completo permanece en disco; carga por capas.

## Orden de descubrimiento

1. `.grok/skills/elite-engineering-library/SKILL.md` (tras `INSTALL_AGENT_BRIDGE.ps1 -Agent Grok`)
1b. Persona: `markdown_system/LIBRARY_HUMAN_GRAPH.md`. Búsqueda: `rg -n -F "<término>" markdown_system/LIBRARY_SEARCH_INDEX.md`. No leas el índice entero.
2. Este archivo
3. `AGENTS.md` (índice; no es lectura obligatoria del corpus)
4. **Un** mapa de autoridad + **un** pack (ver abajo)

## Buscar primero; escalar si falta

Antes de implementar, **busca** en la biblioteca la capability, pack, contrato, gate o fuente oficial. Si no existe, está bloqueado o figura `RESEARCH_INCOMPLETE`, **escala** con evidencia. No inventes fuentes élite, estados de admisión, resultados de tests ni versiones de packs.

## Pack-per-claim

Elige **un** pack por claim acotado en `markdown_system/PACK_PER_CLAIM_INDEX.md`. Confirma la fila correspondiente en `markdown_system/CAPABILITY_CATALOG.md`. Un claim → un pack → verificar → siguiente claim.

## Divulgación progresiva (límites)

| Paso | Cargar | No cargar |
|---|---|---|
| 1 | Skill + `AGENTS.md` | `AGENT_SYSTEM_START.md` entero |
| 2 | **Un** mapa (IA **o** sistemas) | ambos mapas |
| 3 | **Un** pack del índice | todos los `implementation_packs/` |
| 4 | Secciones de gate bajo demanda | `PROJECT_PACK_PLAN.md`, `reconstruction_evidence/` |

## Bridge de proyecto

Cuando la biblioteca no es la raíz:

```powershell
pwsh -NoProfile -File .\ruta\a\elite-engineering-library\INSTALL_AGENT_BRIDGE.ps1 -ProjectRoot . -Agent Grok
```

Abre Grok en la **raíz del proyecto**. Sin Git, el directorio de trabajo debe ser esa raíz.

## Admisión

Sólo `REUSABLE_PACK` o `REBUILD_VERIFIED`/`CONDITIONED` demostrado en el target con evidencia del proyecto. Ítems diferidos (p. ej. Daybreak/libxml2) pueden aparecer en documentación — no marcar completos sin evidencia.
