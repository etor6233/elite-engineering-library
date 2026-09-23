# Markdown Engineering System

Este directorio define cómo convertir el corpus en código materializable por Codex, Claude Code u otro agente.

## Archivos autoridad

| Archivo | Función |
|---|---|
| `PROJECT_BLUEPRINT_CONTRACT.md` | contrato del negocio y requisitos antes de elegir stack |
| `PACK_CONTRACT.md` | formato obligatorio de todo pack con código |
| `COMPOSITION_PROTOCOL.md` | algoritmo de selección, conflicto y materialización |
| `CAPABILITY_CATALOG.md` | cobertura y readiness real del corpus |
| `TOTAL_SYSTEM_CAPABILITY_CONTRACT.md` | clasificación obligatoria de las 48 superficies de infraestructura |
| `MARKDOWN_SYSTEM_READINESS.md` | gate para eliminar el prototipo temporal |
| `AGENT_ERROR_RECOVERY_PROTOCOL.md` | bucle autocorrectivo: diagnóstico, fix canónico, regresión, reconstrucción limpia y continuidad |
| `FAILURE_LEARNING_CONTRACT.md` | memoria obligatoria de cada fallo, recurrencia, causa, corrección, regresión y promoción de lecciones |
| `PROJECT_FAILURE_LESSONS_TEMPLATE.md` | ledger copiable que acompaña a cada proyecto desde discovery |
| `LIBRARY_FAILURE_LEARNING_LEDGER.md` | índice de fallos propios corregidos y condiciones upstream retenidas |
| `DEPENDENCY_UPDATE_CONTRACT.md` | actualización segura de libraries/runtimes/SDKs/images/providers/models con candidate, SBOM y rollback |
| `PROJECT_DEPENDENCY_UPDATE_RECORD_TEMPLATE.md` | inventario y expediente copiable baseline↔candidate por proyecto |
| `AUTHORITY_FRESHNESS_AND_SELF_CORRECTION_CONTRACT.md` | detección de información obsoleta, contradicción, supersession y corrección trazable |
| `PROJECT_AUTHORITY_FRESHNESS_TEMPLATE.md` | registro copiable de fuentes, vigencia y correcciones del agente |
| `ELITE_PUBLIC_AGENT_METHODS.md` | comparación de metodologías públicas y base adoptada |
| `LIBRARY_HEALTH_CHECK.md` | wiring de entrypoints, matriz de dominios y GAPs honestos |
| `USE_LIBRARY_V403.md` | una página para usar o retomar la revisión |
| `LIBRARY_HUMAN_GRAPH.md` | tres gráficos para personas (no corpus de agente) |
| `LIBRARY_SEARCH_INDEX.md` | índice completo para `rg` |
| `AGENT_AUTONOMY_CONTRACT.md` | autonomía por defecto y únicos motivos materiales de pausa |
| `PROJECT_START_READINESS_GATE.md` | intake, accesos, corpus, plataforma y evidencia obligatorios antes de implementar |
| `PROJECT_READINESS_RECORD_TEMPLATE.md` | registro copiable de rondas, 48 superficies, accesos, sources y cálculo de readiness |
| `REVESTEX_ELITE_PUBLIC_CODE_ADMISSION_LEDGER.md` | admisión exacta de producto, SDK, sample y fuente oficial rechazada |
| `PROJECT_INITIALIZATION_PACK_PLAN.md` | composición mínima del source lock y bootstrap oficial sin Git |
| `PROJECT_OFFICIAL_SOURCE_PROFILE_RECORD_TEMPLATE.md` | respuestas, selección, receipts y blockers de los perfiles oficiales materializados |
| `OFFICIAL_DOCUMENT_INTELLIGENCE_PROFILE.md` | selección y gates oficiales por clase/campo documental |
| `STRICT_DOCUMENT_FIELD_EVALUATION_PACK_PLAN.md` | composición de adquisición Stickler exacta + gate cerrado de evaluación, sin storage authority |
| `MICROSOFT_BUSINESS_CENTRAL_CAPABILITY_PROFILE.md` | mapa verificable del producto BCApps y condiciones para elegir esa plataforma |

## Principio

Los manuales explican y gobiernan. Los implementation packs materializan. El agente nunca debe confundir ambos niveles.

```text
manual autoridad ≠ código copiable
architecture pack ≠ implementation pack
snippet ≠ archivo completo
archivo completo ≠ sistema verificado
```

Los packs de `architecture_packs/` existentes permanecen como arquitectura candidata hasta migrar cada capacidad al contrato de `PACK_CONTRACT.md` y superar reconstrucción.

El ciclo de artefactos será compatible conceptualmente con GitHub Spec Kit y adoptará rigor adaptativo de AWS AI-DLC. Esta compatibilidad no convierte automáticamente código upstream en parte de la biblioteca: cada archivo externo conserva su propio expediente de licencia y versión.
