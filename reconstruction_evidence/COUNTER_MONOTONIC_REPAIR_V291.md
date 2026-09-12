# V291 — corrección real del contador, sin readmitir observabilidad

> **Errata de trazabilidad V317 (2026-09-08), prevalece sobre las menciones históricas de Go1.26.7 debajo.** Se retira la atribución histórica exacta Go1.26.7 no vinculada a un receipt de identidad por ejecución. No se sustituye por1.26.8 sin prueba del artefacto concreto. La ruta current-toolchain de V281 contiene1.26.8; official-toolchain contiene1.26.7. El cuerpo original queda conservado como historia. V317 revalida únicamente la composición actual y el scope allí enumerado con1.26.7 observado, no todos los snapshots previos. Ver [TOOLCHAIN_IDENTITY_AND_RUNTIME_SCA_V317.md](TOOLCHAIN_IDENTITY_AND_RUNTIME_SCA_V317.md).

2026-09-07. Mantenimiento correctivo aislado de T2809 / FAIL-387.
GO-OBSERVABILITY-CORE 0.1.0 → 0.1.1. Sigue CANDIDATE, fuera del perfil integral.

## Antes / después

Dos pruebas sobre implementación 0.1.0 reprodujeron:

- Inc(5), Inc(-1): baja a 4;
- Inc(MaxInt64), Inc(1): pasa a MinInt64.

Cinco tests históricos pasaban; ambos negativos nuevos fallaban. No se atribuye
esto a Go: la especificación permite overflow firmado, el defecto era nuestro
contrato de contador sin controles.

La fuente canónica ahora ofrece TryInc(delta) → (valor, aceptado). Bajo el mismo
mutex rechaza negativos y overflow antes de mutar; no entra en panic ni satura
parcialmente. Inc conserva firma y retorna el valor sin cambio cuando rechaza.
Para diagnosticar rechazo hay que usar TryInc y comprobar el booleano; este cambio
no agrega un logger recursivo ni dependencias. No copiar Counter tras usarlo.

## Verificación ejecutada

- 10 tests unitarios PASS: cinco históricos y cinco de negativos, overflow,
  estado/reanudación y dos escenarios concurrentes (32 goroutines).
- Un target fuzz con siete semillas explícitas y oráculo independiente math/big.
- Gate Go nativo materializado y verificado: positivo y dos negativos PASS.
- Fuzz candidate, 10 s y GOMAXPROCS=2: 1.791.391 ejecuciones PASS.
- Dos archivos vuelven a materializarse desde Markdown con bytes idénticos.
- Sobre rebuild: tests, vet y build PASS; fuzz 10 s: 1.774.273 ejecuciones PASS.
  No sumarlas como casos únicos ni afirmar cobertura exhaustiva.
- Tres contraejemplos V280 de privacidad siguen FAIL como se esperaba: email,
  password anidado y mensaje libre. Sólo datos sintéticos; FAIL-386 sigue abierto.
- Compositor rechaza 0.1.1 CANDIDATE incluso con acknowledgeConditions=true y
  conserva destino sin archivos. No se debilitó la admisión.

No se ejecutó race detector (compilador C no resuelto en el PATH inspeccionado),
ni integración de producto, carga, SCA del grafo completo o pruebas de privacidad
general. Las pruebas concurrentes normales no sustituyen -race. Histograma y
sink conservan límites pendientes; no se readmiten por el arreglo del contador.

## Identidad y evidencia

Go 1.26.7 oficial aislado ya fijado por V281, GOTOOLCHAIN=local, GOPROXY=off;
no actualización de runtime, adquisición de código, cuenta ni coste externo.
Procedencia de la corrección y tests: AUTHORED / LicenseRef-Workspace-Owner.
No es copia de un contador publicado por Google ni implementación de la API OTel.

| Artefacto canónico | SHA-256 |
|---|---|
| observability.go | c85d8295f1716bd7aafec3eca9df2536a5b6a341224dc0bf0f095c3089d8413a |
| observability_test.go | 0bfca2f1717957c5c9a85990bc89329d4d7a08172adc6f609b8d60e700241de6 |

Staging preservado: Temp/elite-v291-38e6f7c4dead460892e951e7415a7a58.
baseline/ conserva bytes 0.1.0; candidate/ y rebuilt/ contienen el experimento;
counter-before.log, counter-after.log, rebuilt-tests.log, rebuilt-vet.log,
rebuilt-build.log, counter-fuzz.log, rebuilt-fuzz.log, ambos receipts de fuzz,
perfil de 10s, prueba del gate, privacy-still-blocked.log y candidate-rejection.log.
El harness de privacidad se añadió al rebuild DESPUÉS de tests/vet/build/fuzz y
no se incorpora al pack; demuestra que esas tres fallas no fueron corregidas.
go.mod es sólo harness local sin dependencias, no otro runtime en la biblioteca.

FAIL-411 de invocación: concatenación de ruta pasó filename como Algorithm a
Get-FileHash. Corregido con Join-Path; comparación 2/2 y suites posteriores PASS.
FAIL-412: patch de índices repitió owner, fue rechazado sin cambios; se agruparon
los hunks por archivo y la sincronización pasó sin relajar los controles.

## Autoridades: claims estrechos, no atribución falsa

[OpenTelemetry Metrics API](https://opentelemetry.io/docs/specs/otel/metrics/api/)
describe Counter para incrementos no negativos; no prescribe nuestra API TryInc.
[OpenTelemetry error handling](https://opentelemetry.io/docs/specs/otel/error-handling/)
evita que errores de instrumentación interrumpan el negocio. La política local
de rechazo con booleano aplica esos criterios, no certifica conformidad OTel.
[Go, integer overflow](https://go.dev/ref/spec#Integer_overflow) explica que el
overflow firmado no genera panic: el guard debe comprobarse antes de sumar.
[Google SRE monitoring](https://sre.google/sre-book/monitoring-distributed-systems/)
trata observabilidad como señales operativas; un contador reparado no la completa.
Fuentes consultadas 2026-09-07; no se reaudita todo el corpus por esta consulta.

Se cierra el defecto local FAIL-387 después del retorno canónico y regresiones.
T2809 permanece abierto por privacidad, políticas de eventos/atributos, exportación,
sink, métricas, carga y recuperación. El inventario permanece 160/1433 y perfil67/742.
