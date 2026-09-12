# V330 — devoluciones concurrentes entre pestañas

Mantenimiento y verificación de invariantes existentes de V312/T2804, 2026-09-08.
Resultado focal PASS en Chromium desktop/mobile, Firefox y WebKit. No cambia
código de producto ni reglas comerciales; sólo dos archivos AUTHORED de pruebas.
GO-FRANCHISE-CUSTOMER-JOURNEY-API0.10.18 y PLAYWRIGHT-BROWSER-GATE0.1.25
conservan admisión condicionada y los pins/licencias/dependencias anteriores.

## Oráculo y resultados

Cada navegador ejecuta nueve fases conectadas de fixture; la fase final hace
seis carreras entre dos páginas con la misma sesión autenticada. Una barrera
retiene ambos POST de UI hasta recibirlos y libera dos solicitudes al BFF/API/PG.
Tres carreras llevan igual intención y tres divergen en observaciones/acción.
Cada carrera exige exactamente 200/409, sin asumir cuál pestaña gana.

| Medida focal | Por navegador | Total cuatro proyectos |
|---|---:|---:|
| Carreras | 6 | 24 |
| POST iniciados por UI | 12 | 48 |
| Éxitos / conflictos | 6 / 6 | 24 / 24 |
| Respuestas perdidas después del commit | 2 | 8 |
| Copias iniciales de referencia por opener | 1 | 4 |
| Recibos / decisiones durables | 3 / 3 | 12 / 12 |
| Solicitudes downstream en estado requested | 11 | 44 |
| Eventos y actores comprobados | 6 / 6 | 24 / 24 |

Tras recargar, cada referencia conserva sólo action/authorizationId/receiptId/
digest; no serial ni observaciones. La consulta exacta confirma la ganadora.
Una perdedora idéntica también recupera; una divergente conserva bloqueo y no
presenta la evidencia ajena como confirmación de su intención. Limpiar una
referencia no borra la otra. La ventana abierta con opener recibe una copia
inicial, luego consulta y limpia su propia referencia independientemente.

GET autorizado comprueba observaciones/actor de la intención ganadora, acción
de inventario y solicitudes; el snapshot completo permanece igual después de
recuperar y continuar. Las recargas/consultas no agregan POST. El harness Go
exige marcador de ejecución del caso y comprueba recibos, decisiones, requests,
outbox y actores en PostgreSQL. No acepta un SKIP como prueba multi-tab.

## Alcance y límites

sessionStorage no coordina bloqueos entre pestañas: cada una puede intentar el
comando; la exclusión durable es del servidor y la recuperación compara digest.
No se añadió BroadcastChannel ni un mutex distribuido en frontend. La prueba
de popup verifica copia inicial y aislamiento posterior en estos cuatro
proyectos; no recuperación después de cerrar toda la sesión del navegador.
No sustituye pruebas de permisos entre usuarios/organizaciones, cierre comercial,
refund efectivo, integración de workers/proveedores, carga, seguridad ofensiva,
deploy/restore/rollback ni aceptación humana del artefacto final. TEST08 sigue
pendiente; FAIL423 histórico de Firefox no se cierra por este PASS distinto.

Fuentes consultadas para diseñar el fixture, sin copiar código upstream:
[MDN sessionStorage](https://developer.mozilla.org/en-US/docs/Web/API/Window/sessionStorage)
describe partición por pestaña y copia por opener;
[Playwright BrowserContext](https://playwright.dev/docs/api/class-browsercontext)
y [API testing](https://playwright.dev/docs/api-testing) documentan páginas y
solicitudes en el contexto autenticado. Observación local gobierna el PASS.

## Reconstrucción y ejecución

Composición limpia: 67packs/746archivos. Los 744archivos restantes coinciden con
la base; los dos tests reconstruidos coinciden byte a byte con los ejecutados.
La aplicación usa el build Next previo V321: se comprobaron 745fuentes iguales
antes del delta; next-env.d.ts difería sólo por dos imports generados y se
conservó esa diferencia aparte. No se afirma un nuevo build de producto.
Dependencias existentes enlazadas mediante junction; sin instalación ni upgrade.
Node24.20.0 y Go1.26.7 se verifican por SHA256 antes de cada navegador. PG18.6
loopback con tenants sintéticos separados. TLS autofirmado sólo en fixture.
go test ./internal/platform/httpapi -count=1 PASS en reconstrucción; esa invocación
separada omite suites opt-in de runtime, que se prueban en los cuatro logs.

Reproducir en una composición nueva con build del mismo código, dependencias
admitidas y PostgreSQL sintético disponible. Establecer TEST_DATABASE_URL,
ELITE_WEB_ROOT y ELITE_QUOTE_PROJECT a uno de los cuatro proyectos; habilitar
ELITE_QUOTE_E2E, ELITE_ORDER_E2E, ELITE_DELIVERY_READ_E2E,
ELITE_RETURN_RECOVERY_E2E y ELITE_RETURN_MULTITAB_E2E con valor1. Ejecutar:

```powershell
go test ./internal/platform/httpapi -run '^TestQuoteAcceptanceBrowserPostgres$' -count=1 -v -timeout 10m
```

No ejecutar contra datos reales. Las solicitudes downstream permanecen requested.
Stage local: elite-v330-10415e57a8a14d5594193c1ef51efbc3 en el temporal del sistema.
Logs, runner, before, baseline, parity y artefactos de fallos se conservan allí;
los tests canónicos y este informe son portables. No se genera ZIP/release.

## Fallos preservados

FAIL542: browser-first rechazó selector de botón ambiguo entre tres casos;
browser-second agotó timeout al usar texto exacto del label que incluye opciones
del select. Corrección: article por caso y rol combobox con nombre accesible.
Sin aumentar timeout, desactivar strict mode ni retirar oráculos. browser-third
PASS y reconstrucción + otros tres navegadores prueban la corrección. Los logs
anteriores y traces siguen disponibles. Lecturas/argumentos fallidos de FAIL473
no produjeron evidencia positiva. El mensaje Go multi-tab declara sólo el
oráculo de este escenario, preservando el mensaje de la rama V312 original.

## Evidencia hasheada

| Archivo del stage | SHA256 |
|---|---|
| browser-first.log | `ab2977eb08987eef48217bc9e2ac5d92c3dcc938b6da7f5f2a8fc2918524b4c1` |
| browser-second.log | `19ca57697bbc9ec96c396b31153b808105a6b9b12290cc7af614cc525274657c` |
| browser-third.log | `0599f20f731ef343e3c70154fdcf3a890b1d4a5c646876cf566839a745e5066a` |
| browser-mobile.log | `913f157154d5e997f251e73e179beda8149dbc59f6f61aa3f59d38f6248f8d58` |
| browser-firefox.log | `a0947a747aa94433a0f719debb50b157cdb8b80038ce46e708a2c5b0adb622bd` |
| browser-webkit.log | `69f7c76efc54627029fe5c47efa883616472a8138433422e7333cbf91d950fb1` |
| go-http.log | `0ce5d8a3a8cd22b7ab9c0abfc2d98580d145925b98bc3662c817edc07cc15a85` |
| compose.log | `90d2fc062ba8e80dd8814528d694b0e0036da16d6c4d5fa4bc1fac39afc7b605` |
| rebuild.log | `4c4fbd8afaed1b14f393664e3572cc309b74346413aeaa17aac6d50ac799617e` |
| source-baseline.json | `1be6da86456665a32a403add52772a066636bc717bdf6c0a2a6f172a6e7a2ba5` |
| parity.json | `9cbadfe436e8ac1ae2f7f056d732ea6485b7e652c94e0591721e2bc63aea7fab` |
| run.ps1 | `d32607f2d5b039c936883ede585ba7ac2c9fa54207657808e703c5a405d306cb` |
| next-env-built.d.ts | `1862ac4bbbc5192d4bf562161df66ea547ed3e67173100656ab606ae9797db2b` |

| Test canónico reconstruido | SHA256 |
|---|---|
| microsoft_playwright_browser_gate/tests/enterprise-web.spec.mjs | `f8455f0ecd3d971e217b7ff061b619af8de3186c986866fa50a144e34ceb6a16` |
| internal/platform/httpapi/franchisejourney_test.go | `931e6c2e6b8e84448ab1c8fc6ec44af59117789b5c7276ac12528955eb09730f` |

Gate integrado portable pendiente a continuación; esta evidencia focal no cierra readiness42 ni los10macrofrentes abiertos.

## Cierre integrado

Preflight94: 151pasos PASS; VERIFY_LIBRARY_PASS161packs/1453archivos/
766Markdown/52perfiles. Preflight general BLOCKED únicamente por Docker ausente;
los estados runtime SKIPPED/BLOCKED documentados por cada pack permanecen así.
Plan del contrato PASS; evidence integral BLOCKED por pendientes y evidencias
captured sin promoción integral. TEST48/EVID40 no cambian ese resultado.
Los93eventos previos están intactos; checkpoint94/resume94 PASS con111evidencias.
Checkpoint95 registra este cierre, sin nuevo delta de código desde la reconstrucción.

| Cierre del stage | SHA256 |
|---|---|
| preflight94.json | `5403787f5d29abdca42f6b31fa2f059891a6eaca52d198e72eb650dcb4f64f08` |
| preflight94.log | `8284b45f02a8c199ca01211cf13ae4f8ef800ec282492f21a6eaf499bf6c28b8` |
| contract-plan.log | `2fd65a7ad27c8b968c64f8079abf8610894d6e7732b0b0f514dba962cf93b565` |
| contract-evidence.log | `7283c42afd7da7c28a80322c7183d59bfe4d44f36624f5ee26f1ec3a8aed2f57` |
| checkpoint94.log | `aff833b92e51edf246b992550c0c8fc114f821a12f109951e5ed84aef5cba999` |
