# V284 — límites de escritura y recuperación del bridge

Fecha: 2026-09-07. Owner: T2801 del roadmap vigente; mantenimiento de biblioteca,
no implementación de una franquicia. Procedencia de cambios y pruebas: AUTHORED.

## Resultado demostrado

Tres defectos del instalador se reprodujeron y corrigieron en su fuente raíz
`INSTALL_AGENT_BRIDGE.ps1`. La suite adicional ejecuta 23 aserciones con el
instalador real y filesystem Windows, no mocks de sus funciones.

| Frontera | Antes | Después |
|---|---|---|
| `.agents`, `.agents/skills`, directorio de skill Claude y ancestro del proyecto redirigidos | cuatro junctions permiten escritura fuera del destino | rechazo antes de la primera escritura; destino externo intacto |
| END anterior a BEGIN | éxito sin reemplazar el bloque administrado | rechazo con bytes originales intactos y sin skill creada |
| segundo archivo bloqueado tras reemplazar AGENTS | rollback llama un parámetro inexistente y deja cambios | error original propagado; AGENTS restaurado byte por byte, incluido BOM, o eliminado si antes no existía |

Se comprueban también ausencia de temporales residuales y se conserva la suite
anterior: instalación NEW sin Git, EXISTING con instrucciones propias,
idempotencia por SHA-256, selección Codex/Claude, biblioteca incluida en el
proyecto, conflicto con skill no administrada, límite de instrucciones y WhatIf.
Estos tests prueban los archivos instalados, no que un agente distinto ya los
haya cargado en su sesión ni que un producto haya pasado readiness.

## Evidencia y reproducción

Entorno: Windows, PowerShell 7.6.5, .NET 10.0.11. Sin red de producto, credenciales,
Git, instalaciones globales o gasto externo.

```powershell
./markdown_system/test_agent_bridge_safety.ps1
./markdown_system/test_install_agent_bridge.ps1
./markdown_system/test_library_maintenance_state.ps1
```

Primera ejecución: `BRIDGE_SAFETY_FAILED failures=18 checks=23`. Huellas:
12 aserciones de redirección, 2 de marcadores invertidos, 4 de rollback.
Después: `BRIDGE_SAFETY_PASS checks=23` y suite anterior PASS.
Una copia nueva de nueve archivos canónicos, comparados por SHA-256, repite la
suite y los 23 checks. Los 32 checks de selectores de distribución también PASS.
El instalador es un script fuente de raíz: aquí se verifica copia limpia, no se
presenta como reconstrucción de un pack Markdown. Logs red/green/clean-tests en
staging local `elite-v284-b229f5481daf4ff687a50f2176e4c99d`; las pruebas eliminan
sólo sus propios fixtures y junctions temporales, nunca archivos de proyectos.

| Archivo canónico | SHA-256 |
|---|---|
| INSTALL_AGENT_BRIDGE.ps1 | fe32bf216ea48e9f218dbcf3075161bef445a82bd498bea3f14429d172a12a4b |
| markdown_system/test_agent_bridge_safety.ps1 | 127eb65145b5aee53b12c68349d0378261e1e51378e751a6bdfc13b682e48387 |
| markdown_system/test_install_agent_bridge.ps1 | 265995c5efc49d98374fdb645b51ef49f859629e79d414f7cdf3ce42fd963d1f |

La nueva suite forma parte del verificador existente y de las dos allowlists
exactas de distribución. No se agrega un comodín de scripts ni otro subsystem.
Los packs y el perfil integral conservan 160/1.433 y 67/742 respectivamente.

## Autoridad estrecha, no atribución ficticia

Microsoft documenta el bit `ReparsePoint` en
[FileAttributes](https://learn.microsoft.com/en-us/dotnet/api/system.io.fileattributes?view=net-10.0).
La firma oficial de
[Select-Object](https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.utility/select-object?view=powershell-7.5)
no incluye `-Reverse`; el fallo además se reprodujo en PowerShell 7.6.5.
Consulta: 2026-09-07. Esas páginas gobiernan la semántica de las APIs, no certifican
nuestro instalador. La validación de ancestros, orden de marcadores, iteración
inversa y pruebas son implementación propia declarada, no código copiado de
Microsoft, OpenAI o Anthropic. No se alteró el contenido de las Skills generadas.

## Límites y continuidad

- Se rechazan enlaces ya presentes y se revalida antes de escribir/restaurar.
  No equivale a una defensa contra TOCTOU con un escritor hostil concurrente:
  ejecutar en workspace controlado por el propietario, sin edición concurrente.
- Rollback probado para fallo de reemplazo con filesystem disponible. No es
  transacción resistente a corte eléctrico, disco lleno durante restauración o
  cambios hostiles de ACL. Pueden quedar directorios vacíos creados por la operación.
- Los escenarios de bloqueo por FileShare son Windows; otro sistema no hereda
  ese PASS y la suite declara la condición de plataforma.
- FAIL-390/391/392 se cierran por regresión. T2801 no se cierra: aún faltan
  expedientes de dependencias/freshness/monitoring/fuentes, clasificación de las
  48 capacidades y assurance conectado. T2802–T2810 permanecen en el roadmap.
- No se promovió un pack ni se creó el ZIP final; este avance no certifica
  integraciones live, documentos reales, fiscalidad, carga o producción.
