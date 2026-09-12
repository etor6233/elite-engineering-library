# V262 — confirmación autenticada BFF/Go/PostgreSQL

Fecha: 2026-09-06. Mantenimiento de biblioteca. `REBUILD_VERIFIED / CONDITIONED`; no promoción global ni certificación productiva.

## Confirmación demostrada

`solicitud pública → recurso con jornada → token RS256 verificado → cliente BFF protegido → API Go → PostgreSQL → confirmed → lectura del cliente autorizado`

El nuevo gate usa el verificador OIDC existente, discovery/JWKS loopback y claves RSA efímeras. No usa un verifier que acepte cualquier token. Las identidades y los datos son sintéticos. Crea disponibilidad de organización, slot, solicitud, recurso y jornada del recurso mediante la API real; el lead y su asociación al cliente son fixtures SQL explícitas. Esa asociación no demuestra resolución de identidad real de un visitante.

Dos llamadas concurrentes del cliente BFF a la confirmación con la misma versión producen **una respuesta 200 y un rechazo 409/JOURNEY_CONFLICT**. SQL independiente confirma `state=confirmed`, `version=3`, un recurso asignado, una transición inmutable con `actor_subject=operator` y un evento `appointment.confirmed`. El cliente propietario obtiene ese turno en su timeline; otro subject obtiene cero turnos.

Rechazos ejecutados: firma alterada 401, permiso insuficiente 403, organización ajena 403, tenant ajeno 409, confirmación sin recurso 409, asignación sin jornada del recurso 409 y asignación con versión obsoleta 409. Un intento SQL de borrar el audit falla. No se deshabilitan triggers para este gate.

## Corrección concreta

El BFF protegido perdía los errores Problem Details de Go: 401/403/409 se convertían en 502. Además, aceptaba el prefijo inválido application/jsonp como éxito. Red: **8 FAIL / 6 PASS**. Fix: comparación del media type normalizado exacto y application/problem+json únicamente en error. Green: **14/14** tests focales; el gate conectado comprueba además el 409 real, sin mock de fetch.

Owner TS-OIDC-PORTAL-ADAPTER **0.2.2**: protected-client.ts y su test. Owner GO-FRANCHISE-CUSTOMER-JOURNEY-API **0.9.2**: test HTTP existente. Tres archivos modificados, cero dependencias, migraciones o archivos materializables nuevos. Cambios `AUTHORED`, no código copiado de Microsoft ni nuevo upstream admitido. Licencias y locks existentes se conservan.

## Gates ejecutados

- PostgreSQL 18.6: 49 migraciones aplicadas a bases descartables locales.
- Confirmación conectada: PASS inicial y repetición desde backend reconstruido; cliente BFF ejecutado desde árbol probado, con bytes idénticos comprobados.
- Web: **52 tests PASS**, un test conectado omitido explícitamente en la suite sin opt-in; typecheck y build de producción PASS. El test omitido pasó por separado dentro del runner conectado.
- Go 1.26.7: test ./..., vet ./... y build ./... PASS. La suite general no acredita tests opt-in omitidos.
- Suite PostgreSQL completa: PASS separado en elite_browser_v262.
- Composición limpia: **67 packs / 704 archivos**. Tres archivos ejecutados/reconstruidos idénticos:

| Archivo | SHA-256 |
| --- | --- |
| internal/platform/httpapi/franchisejourney_test.go | 7fc644cd3d92679c329539967e46e31720fe4af5207874d8e214cf1978c897b9 |
| src/platform/backend/protected-client.ts | 949ad3a4d1ed98804146cc5ab068e1741b299c52c8069ebb5dd38d7af686a18b |
| src/platform/backend/protected-client.test.ts | 34e1ac8053ff6c89585cb1dee3f50fe4f6f0582bf1c652dd73b972c752e1daf3 |

Staging: `elite-v262-e9f4da5b09e0441da223fbc3e922b301` bajo el temporal del sistema; reconstrucción válida `roundtrip-r1`. El intento anterior se detuvo por una edición incorrecta de versión del plan, no se contó como PASS. Fallos y correcciones: LIB-FAIL-2088–2090 / FAIL-20260906-354–356.

## Reproducción y limpieza

Componer el perfil vigente en un destino ausente. Instalar su lock web exacto con pnpm offline/frozen. Aplicar las 49 migraciones a una base **descartable dedicada** cuyo nombre comience por `elite_confirmation_`, en 127.0.0.1; nunca usar una base de proyecto. Con TEST_DATABASE_URL y ELITE_WEB_ROOT absoluto al web materializado:

```powershell
$env:ELITE_CONFIRMATION_E2E = '1'
try {
  go test ./internal/platform/httpapi -run '^TestAppointmentConfirmationBFFPostgres$' -v -count=1 -timeout=3m
  if ($LASTEXITCODE -ne 0) { throw 'Confirmation gate failed' }
} finally {
  Remove-Item Env:ELITE_CONFIRMATION_E2E -ErrorAction SilentlyContinue
}
```

El test conserva sus filas sintéticas para inspección en esa base dedicada: no elimina registros inmutables ni modifica protecciones de auditoría. Descartar exclusivamente la base temporal identificada al terminar. No guardar tokens de fixture ni tratarlos como credenciales del proyecto.

## Autoridades oficiales consultadas

- [Microsoft Azure — API design](https://learn.microsoft.com/en-us/azure/architecture/best-practices/api-design): contratos HTTP, errores y representaciones de recursos; no aprobación del código local ni de toda la arquitectura.
- [Microsoft Playwright — Authentication](https://playwright.dev/docs/auth): aislamiento y precauciones del estado autenticado de pruebas. V262 no ejecuta un login Playwright ni afirma cubrir sus cookies.
- [PostgreSQL 18 — Explicit Locking](https://www.postgresql.org/docs/18/explicit-locking.html): serialización mediante locks de fila; dos intentos concurrentes aquí no equivalen a prueba de carga.

Consulta 2026-09-06. Sin actualización de dependencias, nuevo SCA global ni adquisición de código externo. La lógica de negocio del owner preexistente mantiene sus referencias Microsoft BCApps fijadas; esta tarea agrega evidencia y corrige el transporte del error.

## Pendiente, sin ocultarlo

No es un recorrido de navegador autenticado: el gate ejecuta el cliente server-only del BFF, no el formulario, cookie, callback ni ruta Next completa. El portal de operador todavía exige IDs/versión manuales y no tiene una agenda seleccionable con estado recuperable suficiente; debe cerrarse antes de dar esta experiencia por óptima. También quedan notificación externa, cotización/venta, recovery tras pérdida/reinicio de sesión y todos los gates target de IdP real, seguridad, carga, edge/deploy y aceptación. Un evento outbox no demuestra notificación entregada.

El avance global sigue NO_ESTIMADO; no se generó ZIP final. El resto del cierre permanece en FRANCHISE_PREFLIGHT_GAP.md.

## Cierre verificado de esta iteración

`VERIFY_LIBRARY_PASS`: 160 packs / 1.395 archivos materializables / 695 Markdown / 51 perfiles. Franquicia 67 packs / 704 archivos, sin duplicaciones de paths. Memoria: 2.090 lecciones locales + 209 condiciones upstream = 2.299 IDs; este PASS estructural no cierra condiciones upstream ni producción.

Inspección SQL antes de limpiar: dos ejecuciones independientes dejaron exactamente dos tenants sintéticos, dos confirmaciones, dos audits y dos eventos. Se comprobó ausencia de tenants ajenos y se descartó únicamente la base efímera `elite_confirmation_v262` creada en esta tarea; sus fixtures no son datos recuperables de negocio. PostgreSQL temporal quedó detenido. No se eliminaron datos de proyecto ni evidencias canónicas.
