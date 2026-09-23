# Elite Engineering Library

Biblioteca Markdown portable para que Codex, Claude Code, Grok Build, Cursor u otro agente diseñe, componga y materialice sistemas verificables **sin** imponer un backend TypeScript ni copiar código público sin admisión auditada.

Este repositorio es la **biblioteca de ingeniería**, no un producto desplegado. Contiene mapas de autoridad, contratos, implementation packs materializables, código de referencia local y evidencia de reconstrucción. Un PASS de fixtures locales no equivale a cloud, CI alojada, proveedor live, dispositivo físico, aprobación visual ni producción.

## Empezar aquí

Orden de puertas (una sola tabla):

| # | Puerta | Para qué |
|---|---|---|
| 1 | [`README.md`](README.md) (este archivo) | Contexto de biblioteca vs producto |
| 2 | [`AGENTS.md`](AGENTS.md) | Router compacto por agente |
| 3 | [`markdown_system/USE_LIBRARY_V403.md`](markdown_system/USE_LIBRARY_V403.md) | Uso, búsqueda y retoma |
| 4 | [`markdown_system/LIBRARY_HEALTH_CHECK.md`](markdown_system/LIBRARY_HEALTH_CHECK.md) | Salud, wiring y GAPs honestos (huérfanos tip: QR/`GO_QR_CORE`, leads; GTM/SDR pack FALTA) |
| 5 | [`START_FRANCHISE.md`](START_FRANCHISE.md) | Nueva franquicia / consumidor desde la copia elegida |
| 6 | [`reconstruction_evidence/LIBRARY_VS_PRODUCT_GATE_V402.md`](reconstruction_evidence/LIBRARY_VS_PRODUCT_GATE_V402.md) | Gate canónico biblioteca vs producto |
| 7 | [`qualification/FINAL_LIBRARY_READY_V402.json`](qualification/FINAL_LIBRARY_READY_V402.json) | Recibo de cierre V402 (solo evidencia) |
| - | [`docs/FRANCHISE_ARCHITECTURE.md`](docs/FRANCHISE_ARCHITECTURE.md) | **Franquicia / empresa** — dominios end-to-end, matriz HECHO/PARCIAL/NO, orden de build (tip pin) |
| - | [`docs/ROADMAP.md`](docs/ROADMAP.md) / [`docs/FRANCHISE_PLAYBOOK.md`](docs/FRANCHISE_PLAYBOOK.md) | Roadmap (`FRANCHISE / COMPANY DOMAINS` + `LIBRARY_INFRASTRUCTURE`) + playbook de armado (el CAN del playbook ≠ alcance producto franquicia) |
| — | [`START_REFERENCE_V402.md`](START_REFERENCE_V402.md) | **Después** de lo anterior: operadores de ZIP V402 / referencia local (no es la puerta de franquicia) |

No abras `implementation_packs/UNIFIED_REFERENCE_V403_R4.md` dentro del contexto del agente: transporta 2147 fuentes y se materializa con el comando de [`markdown_system/START_V403_LOCAL.md`](markdown_system/START_V403_LOCAL.md).

## Encontrar cualquier archivo

```bash
rg -n -F "término" markdown_system/LIBRARY_SEARCH_INDEX.md
```

Abrí solo la ruta que coincide; no cargues el índice entero ni el corpus.

## Índices principales

| Necesidad | Ruta |
|---|---|
| Persona (tres gráficos; no certifican implementación) | [`markdown_system/LIBRARY_HUMAN_GRAPH.md`](markdown_system/LIBRARY_HUMAN_GRAPH.md) |
| Un pack por claim acotado | [`markdown_system/PACK_PER_CLAIM_INDEX.md`](markdown_system/PACK_PER_CLAIM_INDEX.md) |
| Nueva franquicia desde una copia | [`START_FRANCHISE.md`](START_FRANCHISE.md) (puerta 5) |
| Arquitectura (6 `CANDIDATE_PACK`) | [`architecture_packs/`](architecture_packs/) |
| Implementación materializable | [`implementation_packs/`](implementation_packs/) |

## Integrar en un proyecto

Si la biblioteca no es la raíz del proyecto, instalá el bridge desde su raíz:

```powershell
pwsh -NoProfile -File .\INSTALL_AGENT_BRIDGE.ps1 -ProjectRoot .\ruta\al\proyecto -Agent Both
```

Abrí el agente en la **raíz exacta del proyecto consumidor**, no solo en la biblioteca vendorizada. Detalle completo: [`AGENTS.md`](AGENTS.md) y [`START_ANY_PROJECT.md`](START_ANY_PROJECT.md).

## Verificar una copia

```powershell
pwsh -NoProfile -File .\VERIFY_LIBRARY.ps1
pwsh -NoProfile -File .\VERIFY_EXECUTABLE_LIBRARY.ps1 -Mode Preflight
```

En Linux/macOS usá `./` en lugar de `.\`. `Preflight` nunca convierte toolchains o cuentas ausentes en PASS.

## Referencia local opcional (TypeScript web BFF)

El árbol `src/` incluye un BFF/frontend de referencia para el perfil Go/PostgreSQL. Su readme operativo está en [`docs/ENTERPRISE_WEB_BFF.md`](docs/ENTERPRISE_WEB_BFF.md); no es la entrada de la biblioteca.

## Estado vigente

Consultá [`markdown_system/MARKDOWN_SYSTEM_READINESS.md`](markdown_system/MARKDOWN_SYSTEM_READINESS.md). `READY_FOR_PROJECT_BOOTSTRAP` significa que la biblioteca puede iniciar y acelerar un proyecto; la admisión productiva depende del país, reglas de negocio, proveedor, IdP, infraestructura y artefacto concreto del destino.
