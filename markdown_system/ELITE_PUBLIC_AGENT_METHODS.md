# Métodos públicos de ingeniería con agentes

Fecha de verificación: 2026-09-04.

## Respuesta ejecutiva

Sí existen metodologías públicas, consistentes y reutilizables creadas por organizaciones de primer nivel. Ninguna hace que un agente sea perfecto ni entrega por sí sola un sistema empresarial completo con frontend, backend, datos, seguridad, infraestructura y operaciones.

La base pública más cercana al objetivo de esta biblioteca es **GitHub Spec Kit**. La combinación recomendada es:

```text
GitHub Spec Kit
  ciclo y artefactos: constitution → spec → plan → tasks → implement → converge

+ AWS AI-DLC
  adaptación por riesgo, profundidad variable, aprobaciones y trazabilidad

+ Agent Skills
  empaquetado portable y carga progresiva de conocimiento/procedimientos

+ este corpus
  autoridad técnica, admisión de licencias, implementation packs y evidencia
```

No se incorporará upstream ciegamente. Toda reutilización debe fijar versión o commit, conservar licencia/notices, revisar el contenido y superar un piloto de reconstrucción.

## Comparación

| Método | Organización | Qué aporta | Licencia/reutilización observada | Qué no resuelve | Decisión |
|---|---|---|---|---|---|
| GitHub Spec Kit | GitHub | harness agnóstico; Markdown encadenado; constitution, specify, clarify, plan, checklist, tasks, analyze, implement y converge; integraciones con múltiples agentes | v1.0.1 / commit `9118ed15a0ba65053469a94c560ea5d233f75884`, MIT | no contiene nuestra biblioteca full-stack ni garantiza seguridad, licencias de dependencias o calidad del código generado | `PINNED_CANDIDATE`; ciclo y aceptación ya ligados a la identidad exacta |
| AI-DLC workflows | AWS Labs | workspace detection, greenfield/brownfield, workflow adaptativo, estado de sesión y auditoría | release v1.0.1 / commit verificado `e49341dbeb8af82758dd85e96ed7fe9bcf38a447`, tag object unsigned `e40f6a93...`, MIT-0; source ZIP SHA-256 `d7c2029a...`, release asset `ac060154...` | Operations es placeholder; no se adoptan logging de prompts crudos ni workarounds sin tests; no entrega módulos empresariales completos | `SUPPORTED_REFERENCE_NARROW`; receipt exacto en execution validator 1.3.1 |
| Agent Skills | estándar abierto soportado por varios agentes | `SKILL.md`, metadata, instrucciones y recursos bajo demanda; reduce saturación de contexto | especificación pública; la licencia corresponde a cada skill, no al formato por sí solo | no define SDLC, arquitectura ni calidad; un skill sigue siendo una instrucción probabilística | `SUPPORTED_REFERENCE` para empaquetado portable |
| AGENTS.md + CLAUDE.md | OpenAI / Anthropic | descubrimiento inicial de instrucciones del repositorio | formatos públicos de configuración; no implican permiso para copiar repositorios ajenos | contexto y routing, no enforcement ni verificación | `SUPPORTED_REFERENCE` para arranque fino |
| Conductor | ecosistema Gemini CLI | context → spec/plan → implement; artefactos de producto, stack, workflow y tracks | Apache-2.0 observado en upstream; requiere admisión fijada | menor cobertura de gates/evidencia que la combinación elegida | `DISCOVERED`; referencia secundaria |
| Kiro steering/specs | AWS | steering por producto/tecnología/estructura y specs con requirements/design/tasks | condiciones del producto, no una licencia permisiva general para copiar su implementación | portabilidad y reutilización legal limitadas | aprender conceptos; no incorporar código por defecto |

## Hallazgo sobre empresas de máximo valor

En esta investigación no se encontró una metodología pública integral equivalente proveniente de Tesla, SpaceX, xAI u Oracle que, con licencia permisiva verificada, gobierne de punta a punta a un agente de código y además incluya una biblioteca empresarial lista para materializar.

Esto no significa que esas empresas no publiquen software, papers o prácticas valiosas. Significa que no se admite por reputación ni se atribuye una metodología que no fue encontrada y licenciada. Los activos concretos de esas organizaciones se evalúan individualmente en `LEADING_COMPANY_PUBLIC_CODE_MATRIX.md` y por el estándar de admisión.

## Por qué no existe el “agente perfecto” por instrucciones

Los Markdown orientan el razonamiento; no son controles deterministas. Un agente puede perder contexto, interpretar mal una regla o generar una implementación defectuosa. La consistencia real exige que las reglas importantes se transformen en controles verificables:

| Intención Markdown | Control determinista requerido |
|---|---|
| contrato de API | schema, code generation compatible y contract tests |
| autorización | policy engine/middleware y negative tests |
| invariantes de datos | constraints, transacciones y property tests |
| seguridad | análisis de dependencias, secretos, SAST/DAST y abuse tests |
| baja latencia | workload, presupuesto, profiling y regresión automática |
| recuperación | backup, restore drill y RPO/RTO medidos |
| licencia | lockfile, SBOM, notices y expediente de procedencia |
| completitud | trazabilidad spec ↔ plan ↔ tareas ↔ archivos ↔ tests y convergencia |

Por eso el objetivo correcto no es “un prompt perfecto”, sino un **sistema autocorrectivo y auditable** en el que el agente no pueda declarar éxito sin evidencia.

## Adopción propuesta

1. Mantener `AGENTS.md`, `CLAUDE.md` y `AGENT_SYSTEM_START.md` como routers agnósticos y breves.
2. Alinear los artefactos de proyecto con Spec Kit sin copiar todavía sus archivos: constitution, spec, plan, tasks, análisis y convergencia.
3. Mantener el rigor adaptativo `LIGHT`, `STANDARD`, `HIGH` y `CRITICAL` y la detección NEW/EXISTING en `PROJECT_EXECUTION_STATE.json`; AWS AI-DLC sólo gobierna esos claims estrechos y no reemplaza los owners Elite.
4. Empaquetar capacidades grandes con divulgación progresiva compatible con Agent Skills.
5. Mantener los implementation packs como nuestra capa diferencial: código completo, configuración, pruebas, licencias y reconstrucción.
6. Ejecutar un piloto desde workspace vacío y otro sobre un sistema existente; ambos deben validar resume cursor, tasks, evidencia y checkpoint chain además de los gates del producto.

## Fuentes primarias

- GitHub Spec Kit: <https://github.github.com/spec-kit/>
- Referencia agentic SDD: <https://github.github.com/spec-kit/reference/agentic-sdd.html>
- Licencia MIT de Spec Kit: <https://github.com/github/spec-kit/blob/main/LICENSE>
- AWS AI-DLC: <https://github.com/awslabs/aidlc-workflows>
- Explicación oficial de AWS: <https://aws.amazon.com/blogs/devops/open-sourcing-adaptive-workflows-for-ai-driven-development-life-cycle-ai-dlc/>
- Agent Skills specification: <https://agentskills.io/specification>
- OpenAI Build Skills: <https://learn.chatgpt.com/docs/build-skills>
- Anthropic project memory: <https://code.claude.com/docs/en/memory>

## Observación de ajuste V329 — 2026-09-08

Las decisiones generales anteriores conservan su claim estrecho. Para generar
PROJECT_AUTHORITY_MAP con justificación por decisión/requisito, la revisión
V329 no encontró pack compatible entre Spec Kit, AI-DLC y Conductor. Resultado
NO_ADMISSIBLE_SOURCE acotado, no rechazo de sus otros usos. Revisiones fijadas,
G0–G8, fuentes y límites en reconstruction_evidence/AUTHORITY_MAP_GAP_V329.md.
Las docs actuales de AI-DLC no se atribuyen retroactivamente a v1.0.1; ninguna
revisión nueva, extensión o regla de logging/aprobación se incorpora aquí.
