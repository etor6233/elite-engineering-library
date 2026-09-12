# V276 — routing de estados retenidos hacia envíos autorizados

2026-09-06. Mantenimiento de biblioteca. PYTHON-META-WHATSAPP-CLOUD-ADAPTER 0.10.0 conserva REBUILD_VERIFIED / CONDITIONED. El objetivo es que el agente componga sistemas nuevos o existentes mediante componentes trazables y recorridos conectados; este delta cierra la selección automática del envío en el tramo de statuses, no todo el worker ni el roadmap de franquicia.

## Implementación real

`NewStatusRouter` configura tenant del observer, conexión, política de retención, clave HMAC del outbound y presupuesto de 1..8 rutas distintas. Copia perfil/clave, exige pool >=2 y admite una llamada concurrente por instancia. `ObserveRetained` recibe un principal realmente verificado y el ID del aviso; no recibe IDs de cita/organización/receipt elegidos por el webhook.

El recorrido ejecutado es:

```text
HTTP firmado → verificador Python existente → receipt/job PostgreSQL
→ relectura scoped y revalidación del original
→ message/recipient HMAC → envío accepted + aprobación + scope vigente
→ SEND_RECEIPT exacto → StatusObserver existente → GET scoped
```

La transacción del router conserva locks compartidos sobre conexión activa/receipt durante verificación y observación; otros pool users necesitan capacidad suficiente. Verifica hashes de cuerpo/perfil/política y cantidad de eventos. La proyección de IDs usa claves JSON exactas después del verificador Python compartido: no usa el matching case-insensitive de structs Go para reinterpretar el contrato del proveedor. Una firma válida por sí sola no decide tenant ni permisos.

La correlación exige mensaje y destinatario HMAC, perfil exacto, conexión/organización compatibles, un único envío accepted ligado a aprobación y scope actual de cita/lead. Resuelve todas las rutas y receipts antes del primer write. Lee el path determinista del Sender mediante os.Root, límite 64 KiB, hash anclado en DB e identidad mensaje/destinatario coincidente. La migración 0052 añade sólo un índice parcial de lookup: no nueva tabla/cola. No es UNIQUE porque un histórico ambiguo debe rechazarse, no borrarse o esconderse.

El único escritor de observaciones sigue siendo StatusObserver. Un fallo posterior puede dejar observaciones previas comprometidas: el caller debe mantener el job/inbox y reintentar idempotentemente. El número insertado no representa un ACK ni la entrega de todos los mensajes. El router NO completa jobs, NO marca inbox processed, NO responde al cliente, NO reenvía y NO crea aceptación comercial.

## Fuentes y procedencia

Los cuatro archivos nuevos (router, tests, índice up/down) son AUTHORED. No son un nuevo sample oficial ni código publicado por Meta/Google. No cambian el Python adaptado, los tres archivos Meta VERBATIM, licencia ni lock upstream.

- [Meta/Postman — Message Status Update Notifications](https://www.postman.com/meta/whatsapp-business-platform/request/rgtfq23/message-status-update-notifications), consultado 2026-09-06: id, recipient_id, estados y timestamp; llegada no equivale a orden temporal.
- [Go — Traversal-resistant file APIs](https://go.dev/blog/osroot), consultado 2026-09-06: os.Root confina aperturas y evita escapes por symlinks, con límites por plataforma. Se usa el runtime Go 1.26.7 ya fijado, no otra dependencia. Proteger el directorio y archivos frente a escritores no autorizados sigue siendo necesario; no se afirma aislamiento de kernel ni certificación multiplataforma.

## Verificación ejecutada

Windows, Go 1.26.7, Python 3.14.4 y PostgreSQL 18.6 loopback. Base sintética elite_whatsapp_v276, 52 migraciones aplicadas y 0052 down/up PASS. No cuentas, mensajes, clientes ni datos reales.

- HTTP real guarda el aviso; el router encuentra la cita automáticamente a partir del mensaje. Dos estados fuera de orden producen observed_delivered; dos replays insertan cero. Un campo Statuses con identidad ajena no reemplaza el statuses exacto verificado. Inbox permanece received, un job pendiente y fence accepted con un intento.
- 20 negativos PASS antes de cualquier observación: mensaje desconocido, destinatario distinto, mensajes inbound mezclados, lote conocido+desconocido, permiso, organización, tenant, conexión disabled, organización de conexión, política, perfil, secreto/firma, cuerpo retenido, count, receipt ausente/alterado, clave HMAC ajena, router ocupado, scope de lead cambiado y correlación ambigua.
- Presupuesto/count y claves JSON exactas probados. Constructor nil rechazado.
- Dos identidades distintas con receipts sintéticos anclados: fallo de secret source en el segundo observer deja exactamente una observación y devuelve error. Replay agrega sólo la faltante, siguiente replay cero; dos observaciones finales. Los receipts del fixture no prueban dos llamadas Meta reales.
- Reconstrucción vacía 67 packs/728 archivos, seis hashes idénticos. Suite completa internal/whatsappbridge: 25 tests de nivel superior PASS, cero SKIP; 42 Python PASS. go vet ./... y go build ./... PASS. No se ejecutó full go test ./..., Audit integrado global, SCA nuevo, navegador, carga, ofensiva, Meta live ni deploy/recovery productivo.

El primer intento de compilar los fixtures falló por sets de Principal declarados map[string]bool en lugar del tipo real map[string]struct{}. Se conserva router-tests.log, FAIL-20260906-380 y la regresión corregida; no se alteró autorización para hacer pasar el test.

## Archivos y reproducción

Componer FRANCHISE_COMPLETE_PACK_PLAN.md a destino ausente; aplicar 52 migraciones a una base descartable loopback elite_whatsapp_ en el puerto de fixture; definir ELITE_WHATSAPP_TEST_DATABASE_URL y ELITE_WHATSAPP_PYTHON. Ejecutar go test ./internal/whatsappbridge -count=1 -v; en whatsapp_cloud, Python unittest sobre test_whatsapp_cloud.py y test_status_reconciliation.py; luego vet/build globales.

Temporal conservado: elite-v276-63dca17866124a468db42d3a42b305b1. Logs migrations.log, router-tests.log, router-green.log, router-recovery.log, roundtrip-go.log, python.log, vet.log, build.log e index-down-up.log.

| Archivo | SHA-256 |
|---|---|
| internal/whatsappbridge/status_router.go | 44431bc8b59e3fad3396e028616a017975cd7ba7e0cecfeb445d590707f80898 |
| internal/whatsappbridge/status_router_test.go | cb230a423fdd290355a5753cbb6f94ac407515d18bfeaa7d57ad66a87f2f4dd5 |
| db/migrations/0052_whatsapp_status_routing.up.sql | 21a16e09525c2622c99f6df06c1f14cc58be4046d250e7afa41e6e86e9e28106 |
| db/migrations/0052_whatsapp_status_routing.down.sql | 1354387e5e45e34b199ba14a4392203d7e40e97a5e50f2beaa476950bcc9caae |
| whatsapp_cloud/README.md | e77dee489dd9c9424529231494bbd2ba125fc3140a1debed72a1b02a4d716568 |
| whatsapp_cloud/PROVENANCE.md | 17334ed3c07d6291bd39e17c6ce67f88e608f70ee77116bcdd2fb380bbcc43f1 |

## Lo que sigue abierto

Worker scoped sobre la cola compartida: claim y generación vigentes, finalización coordinada inbox/job, reintentos acotados, agotamiento/DLQ, retención y operación. No llamar a JobProcessor global sin un dispatcher compatible: puede tomar eventos de otros proveedores. Inbound customer messages necesitan otro handler admitido y no pueden descartarse como statuses. Frontend, cuenta/contratos actuales, recepción/entrega live y los demás frentes de FRANCHISE_PREFLIGHT_GAP.md siguen pendientes.

No hay activación automática. Ante fallo, parar consumidor, preservar originales/observaciones/fences y reconstruir; la reversión de 0052 sólo retira el índice, no datos. No se estima el cierre global contando archivos.

## Control general y limpieza

VERIFY_LIBRARY_PASS, proceso exit 0: 160 packs / 1.419 archivos materializables / 709 Markdown / 51 perfiles. Franquicia 67/728. Procedencia actual AUTHORED=1181, ADAPTED=133, VERBATIM=105. Memoria 2.114 IDs locales + 209 upstream = 2.323. No es un Audit ejecutable global repetido ni validación de READY_TO_BUILD de un proyecto.

Después del down/up del índice se repitieron satisfactoriamente los dos recorridos de routing/replay y recuperación parcial (post-migration.log). Se comprobó current_database=elite_whatsapp_v276, 138 tenants sintéticos, 47 webhooks, 20 observaciones y cero sesiones ajenas; se eliminó exclusivamente esa base descartable y se detuvo el PostgreSQL de pruebas. Fixtures/migraciones permiten recrearla, staging/logs permanecen y no se activó ninguna cuenta o envío real.
