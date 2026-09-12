# Microsoft DevSkim current-source reauditoria — V250

## Decision

`REJECTED_COMPONENT / UPSTREAM_OPEN` para uso como SAST reusable de la biblioteca. No se incorporan binarios, source ni reglas. La marca Microsoft, el commit firmado, el build y un test exit 0 no neutralizan una dependencia vulnerable alcanzable por el grafo del analizador.

## Identidad oficial observada

- repositorio: `microsoft/DevSkim`;
- commit actual de `main`: `a452aa506f80928b611c4c7bfe306b8115821aae`;
- fecha del commit: `2026-08-16T15:47:06Z`;
- verificacion GitHub: `verified=true`, `reason=valid`;
- tree: `f6be283e73f106afdb1345c4123539a7a2f9a136`;
- archive GitHub: 935.159 bytes, SHA-256 `2e75a048dcf3e3db97e4551207b0a912d3bc504e2d2689d21d82e48dfcd89410`;
- licencia MIT: 1.095 bytes, SHA-256 `ad0cf28f3381ca9bb0bf101d127402d44c17bfa0991e1a00bff7ae6679e9dada`;
- NOTICE: 49.773 bytes, SHA-256 `78535e301bab95ed0c86398020fbcc1654a4be55db2f48e4d1ba1eb8368b1947`.

La release GitHub y el paquete NuGet mas recientes siguen siendo `v1.0.90`/`1.0.90`, publicados el 2026-07-17; el changelog de `main` menciona cambios posteriores pero no constituye un artefacto publicado.

## Reconstruccion exacta

Se uso el SDK portable oficial .NET 10.0.400. El `nuget.config` del commit dirige a un feed Azure de Microsoft que respondio 401. Sin editar el source, se repitio el restore con NuGet.org explicito y luego se compilo `Microsoft.DevSkim.Tests` para `net10.0`. Build y ejecucion `--no-build --no-restore` con `DOTNET_ROOT[_X64]` explicito terminaron exit 0.

La restauracion emitio `NU1902` en los proyectos library, CLI y tests: `SharpCompress 0.40.0` esta afectado por `GHSA-6c8g-7p36-r338`. Los archivos `.deps.json` producidos confirman `SharpCompress/0.40.0` en el runtime CLI y tests. El commit actual conserva `Microsoft.CST.ApplicationInspector.RulesEngine 1.9.50` y `Microsoft.CST.ApplicationInspector.Logging 1.9.50`; no corrige el mismo blocker que rechazo DevSkim 1.0.90 en V106.

## Limite del claim

DevSkim se presenta oficialmente como security linting multi-lenguaje, no como demostracion completa de seguridad ni reemplazo universal de CodeQL. GitHub CodeQL si soporta Go/JavaScript y query suites de seguridad, pero para repositorios privados de una organizacion GitHub Team requiere ademas licencia GitHub Code Security; no se incorpora como baseline gratuito universal.

Reabrir DevSkim solo ante release/artefacto oficial fijado cuyo grafo elimine la version afectada, con restore reproducible, tests, SCA y regresiones de archivos hostiles. No realizar un override local de SharpCompress y atribuirlo a Microsoft.

Fuentes oficiales: <https://github.com/microsoft/DevSkim>, <https://github.com/microsoft/DevSkim/releases>, <https://docs.github.com/en/code-security/concepts/code-scanning/codeql/codeql-cli>.
