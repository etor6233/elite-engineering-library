# Elite Agent Context Bridge — Reconstruction Evidence V1

Fecha: 2026-08-26  
Alcance: acceso inmediato y progresivo de Codex/Claude a Elite Engineering Library desde un proyecto nuevo o existente, sin exigir Git.

## Autoridad oficial comprobada

- OpenAI `AGENTS.md`: Codex construye la cadena de instrucciones antes de trabajar, busca desde la raíz del proyecto al directorio actual, sólo consulta el directorio actual cuando no encuentra raíz de proyecto y limita la cadena combinada a 32 KiB por defecto: <https://learn.chatgpt.com/docs/agent-configuration/agents-md>.
- OpenAI Skills: nombre/descripción se descubren primero y `SKILL.md` se carga sólo al activar la capacidad; una Skill puede incluir scripts/referencias y debe conservar un scope preciso: <https://learn.chatgpt.com/docs/build-skills>.
- Anthropic Claude Code memory: `CLAUDE.md` carga instrucciones de proyecto, puede importar `AGENTS.md` mediante `@path`, resuelve rutas desde el archivo importador y pide aprobación inicial para imports externos: <https://code.claude.com/docs/en/memory>.
- Anthropic Claude Code Skills: `.claude/skills/<name>/SKILL.md` es la ubicación de proyecto y el body se carga sólo al usar la Skill: <https://code.claude.com/docs/en/slash-commands>.

Estas fuentes gobiernan únicamente el mecanismo de contexto. No prueban la arquitectura ni el código de producto de Elite.

## Implementación local y frontera de procedencia

`INSTALL_AGENT_BRIDGE.ps1` y la Skill que emite son `AUTHORED` localmente. No contienen código de producto copiado de OpenAI ni se atribuyen a OpenAI. El instalador:

1. resuelve proyecto y biblioteca existentes;
2. comprueba los entrypoints de Elite;
3. calcula una ruta relativa portable;
4. integra bloques administrados en `AGENTS.md`/`CLAUDE.md` sin borrar contenido ajeno;
5. instala `.agents/skills/elite-engineering-library/SKILL.md` para Codex y `.claude/skills/elite-engineering-library/SKILL.md` para Claude;
6. rechaza marcadores dañados, Skill homónima no administrada, symlinks de destino y `AGENTS.md` mayor de 32 KiB;
7. preflighta todos los conflictos antes de escribir y revierte escrituras completadas si una escritura posterior falla;
8. no crea repositorio Git, no instala dependencias, no adquiere upstreams, no solicita secretos y no materializa producto.

La Skill enruta al readiness gate, preflight, catálogo y mapas; ordena carga selectiva y conserva la frontera entre orquestación `AUTHORED` y código público admitido con revisión/hash/licencia propios.

## Regresión ejecutada

`pwsh -NoProfile -File markdown_system/test_install_agent_bridge.ps1`:

- proyecto vacío sin `.git`: PASS;
- proyecto con instrucciones preexistentes: PASS, contenido preservado;
- instalación Codex+Claude, Codex-only y Claude-only: PASS;
- biblioteca vendorizada en `tools/elite-engineering-library`: PASS, ruta relativa exacta;
- segunda instalación: PASS, hashes idénticos;
- marcador roto: rechazo antes de escritura;
- Skill homónima ajena: rechazo sin escritura parcial;
- `AGENTS.md` mayor de 32 KiB: rechazo sin escritura parcial;
- `-WhatIf`: cero archivos/directorios creados;
- proyecto y biblioteca iguales: rechazo.

La Skill emitida fue validada además con `skill-creator/scripts/quick_validate.py` oficial: `Skill is valid!`. El fixture validado midió 2.184 bytes; su SHA cambia correctamente cuando cambia la ruta relativa de la biblioteca, mientras que una segunda instalación sobre el mismo target conserva bytes idénticos.

## Fallos retenidos

- `LIB-FAIL-101`: uso incorrecto de `$LASTEXITCODE` para un script PowerShell invocado en proceso; convertido en regresión.
- `LIB-FAIL-102`: cuatro rechazos de cleanup inline —incluido el ZIP temporal literal—; las regresiones limpian sus temporales y los cleanups extraordinarios usan helpers exactos/autocontenibles.
- `LIB-FAIL-103`: el primer Python y el runtime workspace no contenían `PyYAML`; no se instaló implícitamente. Se reutilizó un entorno aislado ya verificado con PyYAML 6.0.3 para ejecutar el validador oficial.

## Resultado y límites honestos

El bridge elimina la copia manual y demuestra acceso determinista al router sin Git. No vuelve `READY_TO_BUILD` a un proyecto: el usuario todavía debe aportar decisiones, corpus, cuentas, sandboxes, permisos y reglas materiales que el readiness gate identifique. Sin raíz Git, Codex debe iniciarse desde la raíz exacta del proyecto. Claude debe aprobar la primera importación externa cuando su interfaz lo solicite. Esas limitaciones pertenecen al descubrimiento oficial y no se ocultan.

No se añadió ni promovió código de producto upstream. Los conteos permanecen en 38 packs y 346 archivos materializables; esta evidencia cubre sólo el mecanismo de entrada del agente.

El smoke de distribución generó 236 entradas, verificó hashes internos y produjo SHA-256 lateral del ZIP; los dos artefactos temporales fueron eliminados después mediante helper allowlisted.

## Gates finales

- `VERIFY_LIBRARY_PASS`: 38 packs, 346 archivos materializados, 227 Markdown; los 17 perfiles declarados compusieron con sus conteos esperados.
- `VERIFY_EXECUTABLE_LIBRARY_PASS -Mode Audit`: 38 packs, 66 upstream sources, cinco artefactos SDK documentales, un artefacto Business Central y siete provider adapters; acquisition/source profiles/Business Central/WhatsApp negativos y positivos pasaron sin red.
- distribución portable temporal: 236 entradas; manifest SHA-256 interno y checksum lateral verificados; artefactos temporales retirados.
