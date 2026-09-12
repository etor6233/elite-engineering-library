# V261 — solicitud de turno conectada y recuperación de respuesta perdida

Fecha: 2026-09-06. Mantenimiento de biblioteca. Estado: `REBUILD_VERIFIED / CONDITIONED`, sin promoción global.

## Resultado demostrado

`catálogo → lead → enlace de turno → disponibilidad publicada → solicitud → commit PG → respuesta perdida → reintento → misma referencia visible`

Cuatro proyectos Microsoft Playwright 1.62.1: Chromium desktop/mobile, Firefox y WebKit. El nuevo gate pasó **4/4**, junto con el gate anterior de captura **4/4**. API Go, BFF construido y PostgreSQL son reales; no se sustituyen las respuestas válidas ni la persistencia por mocks. El fallo de red se inyecta expresamente con `route.fetch` (commit real) seguido de `route.abort` (sólo se pierde la respuesta al navegador). WebKit conserva el shim loopback HTTP/HTTPS del gate existente, no equivale a verificar HTTPS desplegado.

Después de perder la respuesta, el formulario no anuncia éxito; el usuario reintenta con **la misma clave**, recibe el mismo ID, ve una referencia validada y el botón queda deshabilitado para evitar otra solicitud accidental en ese formulario. Go propaga la señal de replay al BFF.

Negativos: payload divergente con misma clave → 409; otra clave sobre horario lleno → 409; falta de clave → 400; estado `confirmed` enviado por el navegador → 400. Al refrescar, el slot lleno ya no aparece.

SQL independiente verifica cuatro solicitudes `requested`, unidas a sus leads/modelos/organización/horarios, cuatro eventos `appointment.requested`, cuatro registros idempotentes completados y cero discrepancias de capacidad. Las cuatro capturas asociadas conservan sus consentimientos/eventos/receipts. Es una reserva de capacidad en estado **solicitado**, no confirmación de un operador.

## Correcciones y red/green

- LIB-FAIL-2082: fixture inicial rechazada por ausencia de jornada de organización (SQLSTATE 23514). Se agregó disponibilidad sintética acotada antes del slot; el trigger no se deshabilitó.
- LIB-FAIL-2083: formulario mostraba éxito sin consumir ni validar el receipt; Chromium alcanzó timeout leyendo response.json. Se valida ID UUID y estado solicitado antes de mostrar referencia. No se atribuye al defecto de validación una causa de transporte no demostrada.
- LIB-FAIL-2084: replay 200 sin `Idempotency-Replayed` (Firefox y luego Chromium). Se agrega el header sólo cuando el repositorio confirma replay.
- LIB-FAIL-2085: WebKit perdía la selección y no enviaba POST antes de completar hidratación. Se deshabilitan controles hasta el efecto de montaje, según Microsoft Playwright; no se añadieron sleeps ni se omitió ese navegador.

La primera matriz fallida se usó sólo para diagnóstico: hubo una reconstrucción web mientras esa ejecución seguía activa. Sus resultados no son evidencia de un artefacto inmutable. La matriz verde posterior se ejecutó con el build terminado, sin modificaciones en curso, y se repitió desde la reconstrucción canónica (lección LIB-FAIL-2086).

## Implementación canónica

- `GO-ELECTROMOBILITY-PUBLIC-CRM-API` **0.2.3**: harness reutilizable de navegador/Go/PG; no duplica backend. Fail-fast al primer fallo, cero retries de tests, base descartable loopback obligatoria, tenant aleatorio y limpieza scoped.
- `GO-FRANCHISE-CUSTOMER-JOURNEY-API` **0.9.1**: header de replay y fixture/checks del turno en el test existente.
- `TS-FRANCHISE-JOURNEY-PORTALS` **0.9.1**: validación del receipt, referencia visible, protección previa a hidratación y contra resubmit tras éxito.
- `MICROSOFT-PLAYWRIGHT-BROWSER-GATE` **0.1.4**: captura compartida y prueba del turno con fallo de respuesta posterior al commit; README ejecutable con prerequisitos y límites.

Seis archivos materializados modificados, ningún archivo de dominio duplicado, ninguna dependencia ni migración nueva. Los cambios son `AUTHORED`, no código copiado de Microsoft o AWS.

## Verificación ejecutada

- Composición limpia: **67 packs / 704 archivos**; comparación de los seis archivos ejecutados con los reconstruidos desde Markdown: **6/6 SHA-256 idénticos**.
- Web: instalación offline frozen, **40/40 tests**, typecheck y build de producción PASS. Node 24.14.1, pnpm efectivo 11.19.0.
- Go 1.26.7: `go test ./... -count=1`, `go vet ./...`, `go build ./...` PASS desde `roundtrip`. La ejecución general sin variables omite expresamente DB/browser; no reemplaza las siguientes.
- Conectados opt-in: `TestPublicLeadBrowserPostgres` y `TestPublicAppointmentBrowserPostgres`, **8/8 proyectos de navegador PASS**, comprobaciones SQL incluidas. Una ejecución focal del nuevo gate tomó 10,6 s de runner; no es un benchmark de rendimiento productivo.
- PostgreSQL 18.6: **49 migraciones** aplicadas a `elite_browser_v261`; `go test ./internal/platform/postgres -count=1` PASS separado.
- Smoke histórico Playwright: **4 runtime + 8 home/CSP PASS**. Los ocho casos conectados se omiten sólo en ese smoke sin opt-in; su PASS proviene del comando dedicado, no de sumar skips.

Staging local: `elite-v261-a39b9809b1b44532ba18486b929ad307` bajo el temporal del sistema. Los Markdown son la fuente de verdad; no depender de ese directorio para reutilizar.

## Reproducción

Hashes de owners/perfil al cierre V261; son snapshots históricos, no afirmaciones de inmutabilidad de futuras versiones:

```json
[
  {
    "path": "implementation_packs/GO_ELECTROMOBILITY_PUBLIC_CRM_API.md",
    "sha256": "69b0f715eab64ef831f510505ceb4b59d2a36444c7e2c557524eb44d424e6557"
  },
  {
    "path": "implementation_packs/GO_FRANCHISE_CUSTOMER_JOURNEY_API.md",
    "sha256": "f5020d042cfe09e5e46646f17a6da3adcf2622ab178b81a0e84adfb1f150ffc8"
  },
  {
    "path": "implementation_packs/TYPESCRIPT_FRANCHISE_JOURNEY_PORTALS.md",
    "sha256": "a08be6fe31881e68177e232d90e08462332e56901050a88b7410d63f65ebd0e7"
  },
  {
    "path": "implementation_packs/MICROSOFT_PLAYWRIGHT_BROWSER_GATE.md",
    "sha256": "34beb353a939915ee602a30aa5525b00170a2b2ddcf1e59c9c7919478b62b3c5"
  },
  {
    "path": "markdown_system/FRANCHISE_COMPLETE_PACK_PLAN.md",
    "sha256": "e41018f1128f8dd09986d207c669ec4202dd9279be278e8d5c28c866f55c64ef"
  }
]
```

Limpieza verificada: cero tenants de fixture `browser-*` residuales y PostgreSQL temporal detenido. Verificador raíz: `VERIFY_LIBRARY_PASS`, 160 packs / 1.395 archivos / 694 Markdown / 51 perfiles; franquicia 67/704. No se generó un ZIP final ni se declaró completitud global.

Componer backend y web actuales; instalar locks exactos, construir web y aplicar migraciones a una base **descartable** `elite_browser_*` en `127.0.0.1`. Desde el backend, con `TEST_DATABASE_URL` seguro y `ELITE_WEB_ROOT` absoluto al web construido:

```powershell
$env:ELITE_PUBLIC_LEAD_E2E = '1'
try {
  go test ./internal/platform/httpapi -run '^TestPublic(Lead|Appointment)BrowserPostgres$' -v -count=1 -timeout=5m
  if ($LASTEXITCODE -ne 0) { throw 'Connected journey gate failed' }
} finally {
  Remove-Item Env:ELITE_PUBLIC_LEAD_E2E -ErrorAction SilentlyContinue
}
```

No ejecutarlo sobre una base de proyecto. La autenticación del harness rechaza todos los intentos protegidos; no se simula un IdP exitoso. El test elimina sólo sus fixtures de tenant y no cambia triggers ni permisos.

## Autoridades oficiales consultadas

- [Microsoft Playwright — Hydration](https://playwright.dev/docs/navigations#hydration): controles interactivos deshabilitados hasta estar funcionales.
- [Microsoft Playwright — Best Practices](https://playwright.dev/docs/best-practices): comportamiento visible, aislamiento, datos controlados y locators semánticos.
- [Amazon Builders’ Library — Making retries safe with idempotent APIs](https://aws.amazon.com/builders-library/making-retries-safe-with-idempotent-APIs/): identidad de petición, efectos únicos ante reintentos y detección de intención divergente.
- [PostgreSQL 18 — Explicit Locking](https://www.postgresql.org/docs/18/explicit-locking.html): semántica de locks de fila; el owner de turnos existente usa `FOR UPDATE`. V261 no agrega una prueba nueva de carga concurrente.

Consulta: 2026-09-06. Runtime Microsoft mantenido en commit `26a9e470a7b3c7822084b09fb7f13902c5f37b51`, Apache-2.0, paquetes/browsers exactos sin cambios. No se ejecutó una nueva auditoría SCA global ni se atribuye a Microsoft la lógica local.

## Pendiente real

Confirmación/operador/recursos con identidad real, notificación externa, cotización/venta, recuperación tras reinicio del navegador, antiabuse, política legal, accesibilidad asistiva, carga/ofensiva/edge/deploy y aceptación. El resto de brechas permanece en `FRANCHISE_PREFLIGHT_GAP.md`; esta evidencia no cierra toda la biblioteca ni autoriza producción.
