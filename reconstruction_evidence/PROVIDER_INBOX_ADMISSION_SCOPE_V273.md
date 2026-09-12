# V273 — admisión durable de proveedores con conexión vigente y replay exacto

Fecha: 2026-09-06. Mantenimiento de biblioteca, no despliegue de franquicia.

## Hallazgo real y corrección

Al revisar el receptor compartido para conectar WhatsApp se encontró que AcceptWebhook no consultaba el estado actual de integration.provider_connection. Una conexión deshabilitada seguía insertando eventos/jobs. Además, ON CONFLICT seguido de comparación exclusiva de body_sha256_hex podía reconocer como replay un receipt con otro proveedor, tipo o payload. El registry de arranque no resolvía la revocación durable. No se registró un incidente productivo: el hallazgo se reprodujo con fixtures locales.

La prueba inicial TestProviderWebhookConnectionAndReplayScope falló en seis subcasos; el fixture terminó con dos eventos y dos jobs en vez de uno. FAIL-20260906-376 / LIB-FAIL-2110 conserva el fallo. La corrección AUTHORED del pack GO-PROVIDER-INTEGRATION-CORE 0.1.2:

- comprueba tenant/conexión/proveedor exactos y estado active dentro de la misma transacción;
- mantiene FOR SHARE sobre la conexión hasta commit, serializando una desactivación concurrente;
- exige que replay conserve provider_code, event_type, body hash y payload JSONB; acepta diferencias de formato JSON sin cambiar semántica;
- rechaza nil repository/pool sin panic;
- conserva la única transacción evento + job existentes, sin nueva cola, tabla o worker.

El INSERT no es confirmación de procesamiento: un job ya comprometido requiere su propia revalidación antes de efectos. Deshabilitar conexión no cancela retroactivamente trabajos ni envíos existentes. La autenticación de bytes debe seguir haciéndola el adapter oficial del proveedor; esta corrección no convierte el HMAC de referencia Elite en firma Meta.

## Evidencia ejecutada desde reconstrucción

- TestProviderWebhookConnectionAndReplayScope: siete subcasos PASS; disabled new/replay, tipo, payload, proveedor, conexión inexistente y desactivación concurrente. En este último se observó la espera de lock en pg_stat_activity, se comprometió disabled y el receptor rechazó sin persistir otro evento/job. Cierre del fixture: un evento y un job originales.
- TestProviderWebhookDeduplicationAndJobAreAtomic anterior PASS; el conflicto de body hash y reconciliación existentes siguen funcionando.
- Cuatro tests del package internal/providerintegration PASS.
- TestProviderWebhookHTTPAdmissionAndReplay PASS. Este test HTTP usa el fixture de repositorio del adapter; no se presenta como recorrido HTTP→PostgreSQL ni webhook Meta live.
- go vet ./... y go build ./... PASS.
- 51 migraciones aplicadas a una base descartable nueva elite_provider_v273. No cambia schema; no se repite ni se atribuye nuevo down/up, backup o restore.
- Perfil integral reconstruido en carpeta ausente: 67 packs / 722 archivos. Dos archivos modificados idénticos por SHA-256 entre staging y reconstrucción.

| Archivo | SHA-256 |
|---|---|
| internal/platform/postgres/providerintegration.go | 82bfc01c9615f9c2b96115c1dcf7f261c75036ce67fd819d314f7832ad0a41bd |
| internal/platform/postgres/providerintegration_integration_test.go | cb81f9228899be5b7f83b7d24c05b09a07dbf3a9562ae7ec252818a48d28d997 |

Entorno: Go 1.26.7, PostgreSQL 18.6, Windows. Temporal conservado elite-v273-46628c4593d141b09e129b62a32ef904; logs migrations.log, red.log, green.log, green-concurrent.log, roundtrip-pg.log, roundtrip-domain.log, roundtrip-http.log, roundtrip-vet.log y roundtrip-build.log.

Comandos reproducibles: materializar el perfil; aplicar migraciones a una BD local descartable elite_provider_ en 127.0.0.1:55959; configurar DATABASE_URL; ejecutar go test ./internal/platform/postgres -run '^TestProviderWebhook' -v -count=1, go test ./internal/providerintegration -v -count=1 y go test ./internal/platform/httpapi -run '^TestProviderWebhookHTTPAdmissionAndReplay$' -v -count=1. Las tres selecciones focales no tuvieron SKIP.

## Límites, errores de mantenimiento y procedencia

El intento previo de ejecutar todo el package HTTP terminó exit 0, pero el guard local rechazó siete SKIP correspondientes a browser/PG opt-in: tres tests y cuatro variantes de agenda. No eran evidencia de navegador ejecutado. Se corrigió el scope mediante inventario de tests y se ejecutaron dominio/HTTP del adapter y vet/build por separado. FAIL-20260906-377 / LIB-FAIL-2111 conserva el hecho. roundtrip-transport.log mantiene los SKIP originales. No hay nuevo PASS global de pruebas Go/PG, navegador, carga o SCA.

La lectura de un plan minificado fue truncada y su guard rechazó la propuesta antes de editar; se capturó después la línea exacta y sólo se cambió el packId/version correspondiente. Un patch de lecciones con hunks fuera de orden también se rechazó y se corrigió sin cambios parciales. Recurrencia en FAIL-20260906-372. El enlace Meta developers de componentes no fue accesible; la página oficial Meta/Postman de statuses sí se pudo consultar (recurrencia FAIL-20260906-373). No se derivaron reglas de fuentes secundarias.

Fuentes oficiales consultadas el 2026-09-06, con claims limitados:

- [PostgreSQL 18 — Explicit Locking](https://www.postgresql.org/docs/18/explicit-locking.html): FOR SHARE y serialización de actualizaciones de la fila; no publica este adapter.
- [Microsoft Azure — API design](https://learn.microsoft.com/en-us/azure/architecture/best-practices/api-design): semántica consistente de operaciones e idempotencia; no certifica esta implementación.
- [Meta — Message Status Update Notifications](https://www.postman.com/meta/whatsapp-business-platform/request/rgtfq23/message-status-update-notifications): scope y timestamps de callbacks; consulta para el siguiente montaje, no evidencia de un receiver integrado.

No se copió nuevo código upstream ni cambió una dependencia, licencia, runtime o tabla. Se conserva clasificación AUTHORED/CONDITIONED. La versión anterior carecía de estos controles: no usarla como alternativa de rollback de seguridad; detener admisión afectada y reconstruir la corrección si falla el despliegue. Los tres planes consumidores apuntan ahora a 0.1.2. Credenciales, retención, costos, cuenta, operación live y demás gates target siguen abiertos.

## Siguiente conexión pendiente

La bandeja y el job genéricos están disponibles con el control corregido. Todavía falta conectar el verificador Meta existente a una recepción HTTP concreta, fijar retención protegida de bytes originales, routing tenant/mensaje, worker y acknowledgement/recovery; después, la UI del operador y prueba con cuenta real. No crear un ledger duplicado ni marcar procesado un aviso sólo porque su recepción quedó guardada.

## Cierre del incremento

VERIFY_LIBRARY_PASS, proceso exit 0: 160 packs / 1.413 archivos materializables / 706 Markdown; 51 perfiles y franquicia 67/722. Log verify-library.log. No es un nuevo Audit ejecutable global ni certificación productiva. Tras comprobar la base exacta elite_provider_v273 vacía de tenants/events/jobs y cero sesiones activas, se eliminó sólo esa BD descartable y se detuvo el PostgreSQL local iniciado para las pruebas. Los fixtures/migraciones la recrean; evidencia y staging permanecen. No se dejaron servicios, suscripciones, cron ni envíos activos.
