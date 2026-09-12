# Microsoft BCApps Inventory Source Freshness — V158

## Resultado focal

`REBUILD_VERIFIED / CONDITIONED` para la adquisición oficial Microsoft BCApps y su lane Business Central.

- `OFFICIAL-UPSTREAM-ACQUISITION-CORE` 0.4.76: 25 archivos reconstruidos; lock 121/121, 9 negativos y 16 perfiles/5 negativos/1 positivo PASS.
- `MICROSOFT-BUSINESS-CENTRAL-EXECUTABLE-PLATFORM` 0.1.1: 6/6 archivos reconstruidos; config, dos negativos y regresión del source lock PASS.
- Perfil Business Central: 1 pack / 6 archivos.
- `VERIFY_LIBRARY_PASS`: 94 packs / 1015 archivos materializables / 514 Markdown / 42 perfiles.
- `VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit`: 94 packs / 121 fuentes oficiales / 1 artifact Business Central / 14 adapters de proveedor.

## Autoridad, licencia y procedencia

- Repositorio oficial: `microsoft/BCApps`.
- Commit fijado: `2eae56d704a1fd035d104f333602aea7091b7749`; tree `f9846fb1254c6311985c6131bb2af11c7f169e1b`; firma GitHub verificada; fecha del commit `2026-08-31T17:57:04Z`.
- ZIP exacto: 226.293.141 bytes; SHA-256 `f7e984f2a1e9784a351068f42f0a0cefdc8317704b309f2a94fce5f0f91bc5b6`.
- Licencia raíz MIT: SHA-256 `c2cfccb812fe482101a8f04597dfc5a9991a6b2748266c47ac91b6a5aae15383`; licencias/terceros siguen siendo path-specific.
- Microsoft afirma en el README oficial que el repositorio contiene el código, pipelines y tooling usados para construir Business Central. Elite no atribuye a Microsoft ningún script local de adquisición, validación o aprobación.

## Superficie oficial verificada

- Archive: 56.479 entradas y 36.673 archivos AL.
- BaseApp W1: 8.154 archivos AL.
- Inventory: 958 AL; Warehouse: 355 AL; SCM-Reservation: 22 AL.
- Availability: manifest determinista de 49 archivos `cf0b1c27267f341a912cf2bd69b0d3e8cc4f529cc667777d43d6e601b250189b`.
- Tracking/reservas: 108 archivos `a1d432cb8b9a995712b57d1b4e2bc6839ae8a2b2d20a7bf5be4a32c4eba0f9f0`.
- Warehouse: 372 archivos `95fab623e5d4c28935469fddbde128ce7c7cc0f6431113b7bfba8a16d1f08e62`.
- SCM-Reservation con manifest de app: 23 archivos `86de5f6d08f4014c1905c371fc93ea3ac5ab47f8e138d5754a6555673ed1c0b2`.

Microsoft documenta reservas como vínculos firmes supply↔demand, ATP a partir de inventario/recepciones/requerimientos y disponibilidad warehouse descontando picks, movimientos, reservas y bins outbound. Los paths oficiales contienen `AvailabletoPromise`, `ReservationEngineMgt`, `CreateReservEntry`, allocation policies, item tracking y transfer-reservation tests.

## Corrección de frescura

La release firmada `releases/28.4/StrictMode@cb07eef12935e07dd4258c122d0e7b1920fc1223` fue descargada y verificada por separado: 40.855.028 bytes, SHA-256 `48410c37872ea1a974902282a86f662ede8282ae8c5b6adb58fc56cf8e7a9a4a`, 4.824 AL. No contiene BaseApp Inventory/Warehouse ni SCM-Reservation. Por eso no sustituye el source funcional y se conserva únicamente como identidad de release del artifact 28.4.

## Fallos convertidos en memoria

- `LIB-FAIL-1486`: una release StrictMode no implica la superficie BaseApp del snapshot integral.
- `LIB-FAIL-1487`: segunda recurrencia de pipe directo tras `foreach`; colección explícita obligatoria.
- `LIB-FAIL-1488`: globs Markdown literales no son paths válidos para `rg` en Windows.
- `LIB-FAIL-1489`: un hunk con hash transcripto incorrectamente fue rechazado atómicamente.
- `LIB-FAIL-1490`: planes JSON compactos y multilínea deben sincronizarse por representación literal.
- `LIB-FAIL-1491`: los controles del ledger deben contar únicamente IDs en la primera celda de filas canónicas, no referencias históricas legítimas dentro de otras celdas.

## Límites honestos

El código empresarial es Microsoft y está listo para adquisición exacta. No es portable fuera de Business Central/AL y no se ejecutó aquí porque faltan Windows containers, runtime, EULA/entitlement, licencia o demo, country view y autoridad para crear el contenedor. Las suites oficiales SCM-Reservation todavía no son PASS: sólo se probó su identidad exacta. No completa por sí mismo el inventario Go/PostgreSQL existente y no autoriza portar silenciosamente reglas Microsoft. Producción conserva CDN/WAF, IdP, proveedores, PostgreSQL/recovery cuando aplique, carga, seguridad ofensiva, deploy/rollback y aceptación empresarial reales.
