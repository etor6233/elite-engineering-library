# Evidencia de reconstrucción V240 — control de ejecución y reanudación

Fecha: 2026-09-04  
Pack: `EXECUTION-VALIDATOR`  
Versión: `1.2.0`

## Claim estrecho demostrado

El pack materializa 15 archivos exactos. Siete archivos `AUTHORED` nuevos agregan un control persistente para proyectos `NEW` y `EXISTING`: baseline y delta explícitos, rigor proporcional, referencias internas ligadas por SHA-256, paridad exacta con las tareas de `tasks.md`, cursor de contexto mínimo, slice activo y checkpoints append-only encadenados. El estado no reemplaza los owners de readiness, Spec Kit, fallos, dependencias, evidencia o release; sólo los indexa para reanudar sin inventar ni releer indiscriminadamente.

## Fuentes públicas fijadas y uso permitido

| Fuente | Identidad verificada | Claim adoptado |
|---|---|---|
| GitHub Spec Kit | `v1.0.1`, commit `9118ed15a0ba65053469a94c560ea5d233f75884`, MIT | cadena constitution/spec/plan/tasks/implement/converge y tareas verificables |
| AWS AI-DLC | `v1.0.1`, commit GitHub-verificado `e49341dbeb8af82758dd85e96ed7fe9bcf38a447`, MIT-0 | detección greenfield/brownfield, profundidad adaptativa, continuidad de sesión y auditoría |
| AWS AI-DLC source ZIP | SHA-256 `d7c2029a5957a4fc16c43e688a3638f41ea833e4a97a372195ee4772215d3958` | receipt de la revisión inspeccionada |
| AWS AI-DLC release asset | 109267 bytes, SHA-256 `ac0601544b6c7ba41b541a7a96259d6ea3995ea0b94b2a58a333ae7898e22b39` | receipt del artefacto publicado |
| NASA Systems Engineering Handbook | apéndice oficial consultado 2026-09-04 | separación de requirements, verification, validation e integración de sistema completo |
| Google SRE Release Engineering | página oficial consultada 2026-09-04 | identidad exacta del release, builds repetibles y evidencia archivada |
| Google SRE Reliable Product Launches | página oficial consultada 2026-09-04 | gates proporcionales, robustos, adaptables y rápidos en el camino común |
| xAI Grok Build | commit `19d42e35c07a9c9244f03f6df0c4c353f970d4f9`, Apache-2.0 y notices | sesiones, herramientas, permisos, sandbox y compactación; no se presenta como método empresarial completo |

No se incrustó código de estas fuentes. Los siete archivos nuevos son `AUTHORED`, con procedencia honesta y reglas gobernadas sólo por los claims estrechos anteriores. De AWS AI-DLC se excluyeron deliberadamente su área Operations placeholder, el logging de prompts crudos y cualquier workaround que permita continuar sin tests.

## Artefactos agregados

```text
engineering_execution_kit/execution_method_lock.json
engineering_execution_kit/execution_state.template.json
engineering_execution_kit/execution_state.schema.json
engineering_execution_kit/EXECUTION_STATE_PROTOCOL.md
engineering_execution_kit/validate_execution_state.py
engineering_execution_kit/checkpoint_execution_state.py
engineering_execution_kit/test_execution_state.py
```

`execution_method_lock.json` conserva las identidades, licencias, hashes y exclusiones exactas. Cada bloque materializable conserva además su SHA-256 en `implementation_packs/ENGINEERING_EXECUTION_VALIDATOR.md`.

## Verificación focal ejecutada

```powershell
python -m compileall -q .\engineering_execution_kit
python .\engineering_execution_kit\test_execution_state.py
python .\engineering_execution_kit\test_validate_project.py
```

Resultado focal: **20/20 tests PASS** — 13 del control de estado y 7 del validador histórico — después de materializar desde Markdown en una raíz limpia reconocida con sus autoridades exactas. También pasó la sincronización pack↔árbol y la verificación de confinamiento: ningún estado, evento o evidencia puede escapar de la raíz del proyecto ni atravesar symlinks admitidos.

Resultado global final:

```text
VERIFY_LIBRARY_PASS packs=151 materialized_files=1338 markdown_files=661
VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit packs=151 upstream_sources=121 execution_state_controls=1
```

El runner global materializa ahora el pack, compila sus 15 archivos y exige las 20 regresiones en cada Audit; la verificación focal dejó de ser una operación manual separada.

Los casos cubren estado válido `NEW`/`EXISTING`, baseline/delta brownfield, paridad de tasks, evidencia ausente/alterada/fuera de raíz, secretos, transiciones, hash chain, secuencia, cierre y producción fail-closed.

## Fallos aprendidos durante la reconstrucción

- `LIB-FAIL-1867` / `FAIL-20260904-133`: directorios temporales UUID sin limpieza destructiva previa.
- `LIB-FAIL-1868` / `FAIL-20260904-134`: suite limpia con sentinel y autoridades exactas.
- `LIB-FAIL-1869` / `FAIL-20260904-135`: invocación PowerShell directa para preservar `-Files`.
- `LIB-FAIL-1870` / `FAIL-20260904-136`: patches documentales contra contexto observado, no supuesto.

Ningún gate se debilitó.

## Límites que permanecen verdaderos

Este control demuestra continuidad y disciplina de ejecución, no el funcionamiento de una franquicia concreta. CDN/WAF, IdP, cuentas y permisos de proveedores, PostgreSQL/recovery, carga, seguridad ofensiva, despliegue, reglas regulatorias y aceptación empresarial deben demostrarse en el proyecto real. El validador impide convertir su ausencia en un falso `COMPLETE` o `PRODUCTION`.
