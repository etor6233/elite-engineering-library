# V271 — observaciones firmadas → PostgreSQL → consulta autenticada

Fecha: 2026-09-06. Mantenimiento de biblioteca. Estado `REBUILD_VERIFIED / CONDITIONED`, no certificación de Meta live ni de la franquicia completa.

## Resultado conectado

StatusObserver reutiliza el resolver PostgreSQL de aprobación, el fence outbound y el normalizador/correlador Python ya existentes. Obtiene del fence **accepted** el hash confiable de SEND_RECEIPT; no acepta del caller un anchor elegido libremente. Comprueba tenant, principal verificado, permiso appointment:manage, organización, asociación actual appointment/lead, perfil y hashes de la aprobación/fence.

Ejecuta Python -I/-B, sin shell, con hashes exactos del ejecutable y ambos scripts, ambiente mínimo, secreto por stdin, stderr descartado, input/output acotados y deadline de 20 segundos. El proceso usa la misma verificación Meta HMAC/scope y devuelve un resultado ligado por hash al frame exacto. Sus archivos temporales sólo contienen observaciones pseudónimas y se eliminan al terminar; la prueba del original se valida en memoria.

Después revalida scope en una transacción SQL: FOR UPDATE sobre el fence y FOR SHARE sobre appointment/lead. Guarda lotes de evidencia y eventos únicos por tenant/delivery/event key en dos tablas de observaciones, no en otro ledger de envíos. La evidencia es append-only. Un conflicto de payload con event key existente revierte el lote completo. No cambia el estado/attempt_count del envío ni autoriza retry.

La consulta GET existente une las observaciones en el mismo snapshot SQL. Agrega provider_event_count y provider_timestamp (Unix seconds). Devuelve observed_sent/delivered/read/failed/deleted según el timestamp más reciente registrado, o ambiguous_latest_timestamp si hay distintos estados empatados. Sin observaciones conserva not_observed_by_this_reader. No expone teléfono, provider ID, hashes ni app secret. Estas son observaciones registradas, no un estado global infalible o una autorización de reenvío.

## Evidencia ejecutada

- **42 tests Python PASS**, incluidos los 41 anteriores y una nueva prueba del proceso real aislado, binding de bytes, firma inválida, stderr acotado y cleanup. Repetidos desde reconstrucción limpia.
- **Suite completa del package Go internal/whatsappbridge PASS**, con ELITE_WHATSAPP_PYTHON y ELITE_WHATSAPP_TEST_DATABASE_URL explícitos: no se omiten sus pruebas PostgreSQL/proceso por falta de entorno. Go vet ./... y build se ejecutan por separado.
- La suite del package se repitió desde la reconstrucción limpia después de 0051 down/up; también go vet ./... y go build ./... PASS.
- Nuevo recorrido: fixture send_bridge con transporte de proveedor sintético → SEND_RECEIPT real/anchor → cuatro procesos de verificación concurrentes → tres observaciones únicas → GET HTTP real con RS256/JWKS local y sesión SQL read-only → observed_read. El total de inserciones concurrentes es tres, no doce.
- Un segundo estado con el mismo timestamp produce ambigüedad; quedan cuatro observaciones y el fence mantiene exactamente **un intento y dos eventos** (sending/accepted).
- Once subcasos negativos: firma, receipt, tenant, permiso, organización, destinatario, hash de adapter, hash de reconciler, fence unknown, lead reasignado y perfil. Todos dejan cero lotes parciales para ese fixture.
- Evidencia divergente entre lotes revierte también el evento nuevo que ya se había insertado dentro de la transacción; queda un solo evento y un solo lote anteriores.
- UPDATE de observaciones rechazado por trigger. Respuesta HTTP sin secretos ni identificadores privados. Los tests anteriores de lectura read-only, aprobación, proceso y API continúan pasando.
- **51 migraciones** aplicadas en base descartable elite_whatsapp_v271; 0051 down/up ejecutadas sobre datos sintéticos, seguidas por repetición de la suite sobre la reconstrucción. No es un restore drill productivo ni rollback sin pérdida de evidencia.

Las firmas Meta, las identidades y el transporte de envío son fixtures sintéticos. No se contactó Meta ni se probó una suscripción/cuenta real. No hay nuevo PASS de navegador ni suite Go/PG global: el alcance ejecutado es el package afectado, más vet/build globales.

## Procedencia y autoridades

StatusObserver, migración 0051, tests, loader y extensión del lector son **AUTHORED**, no código copiado de Meta, Microsoft o Google. Los ejemplos Meta exactos, licencia/notices, source lock y dependencias permanecen intactos. La frontera de firma sigue siendo la adaptación admitida; no se añade un segundo parser/HMAC con reglas divergentes.

Consultas oficiales del 2026-09-06:

- [Meta: Message Status Update Notifications](https://www.postman.com/meta/whatsapp-business-platform/request/rgtfq23/message-status-update-notifications): identificación de mensaje/destinatario, timestamps y posible desorden de callbacks.
- [PostgreSQL 18: Explicit Locking](https://www.postgresql.org/docs/18/explicit-locking.html): bloqueos de fila y duración transaccional. El diseño de las tablas/queries es local y no se atribuye a PostgreSQL ni a Meta.
- [Python: importlib, importing a source file directly](https://docs.python.org/3/library/importlib.html#importing-a-source-file-directly): carga del módulo local fijo bajo -I. La implementación no admite rutas de código elegidas desde un webhook.

No se admite un upstream nuevo por reputación ni se acredita SCA global actualizado mediante estas pruebas.

## Reconstrucción y límites

PYTHON-META-WHATSAPP-CLOUD-ADAPTER **0.8.0**, 28 archivos. Perfil integral **67 packs / 722 archivos**, cuatro archivos materializables nuevos. Diez archivos cambiados coinciden por SHA-256 entre staging probado y nueva reconstrucción:

| Archivo | SHA-256 |
|---|---|
| whatsapp_cloud/whatsapp_cloud.py | 30eb577f27454a654bd4546c6ad017cbac18bb2cd3fb612066824c32e589812f |
| whatsapp_cloud/status_reconciliation.py | ec31c2442ef176e89c4cd03de179345d1805d4f528b0710a9a4ed060ba7484ec |
| whatsapp_cloud/test_status_reconciliation.py | 1513c03f4e07e63118af547bdcee365262ec237e86a5624ee906b2aea6c18d15 |
| whatsapp_cloud/README.md | 8dbe133790d95b66123e2fac2a88cabd347be35225837f18cb4c24a36ca9cf2b |
| whatsapp_cloud/PROVENANCE.md | 79dc0d48f08c3078fd45860332fe52c1bdd90560aff059df7a76d991660d29aa |
| internal/whatsappbridge/status_observer.go | 8db8dd24e4101f89fd323954b9a8a9f648c41995b2a8c5e41a3e2668c061177a |
| internal/whatsappbridge/status_observer_test.go | 2b971c2f6e6b9c8d236daf240692ab02ff7041aea6c940d92d2a4cd816aef57b |
| internal/whatsappbridge/appointment_notification_status.go | 341b778d54752eeefef94d25d82c07083de2e296923b1c2b7d442cfbb8d7ccbf |
| db/migrations/0051_whatsapp_status_observations.up.sql | 00514d0f9c3285e38a10839d291edc966b96f235c1e97375823ec625319ac346 |
| db/migrations/0051_whatsapp_status_observations.down.sql | fce14de4dae9007ccc2dfa0139251161fce502ebaf8a9b39eb41bc41253c2ef5 |

Reproducción: README materializado, Python 3.14.4, Go 1.26.7 y PostgreSQL 18.6. Staging en el temporal del sistema: elite-v271-561dd5376d634665b8e5b938980dbf8a. Logs: python.log, go-whatsapp.log, migrations.log, migration-reversal.log, roundtrip-python.log, roundtrip-whatsapp.log, roundtrip-vet.log y roundtrip-build.log.

Fallo real de mantenimiento: un patch generado incluyó blancos como contexto, fue rechazado antes de editar y se corrigió prefijando todas las líneas nuevas. FAIL-20260906-374 / LIB-FAIL-2108. No se ocultó ni se relajó el compositor; después se reconstruyeron los diez archivos exactos.

## Qué queda en este recorrido

El servicio de observación durable y su lectura ya están conectados y probados localmente. Quedan el montaje con el ingress/inbox real y su política de acknowledgement/recovery, la pantalla del operador, sus pruebas de navegador, la suscripción/cuenta Meta y evidencia de entrega reales. No se activa un worker, cuenta, cron ni tráfico externo por materializar el pack. El caller debe autenticar al principal y usar una configuración tenant-bound inmutable; los permisos mínimos de BD, retención de pseudónimos, recuperación y protección del filesystem se prueban en el target. Un send unknown sin receipt anclado sigue bloqueado y no se resuelve fabricando evidencia.

## Control general y limpieza

VERIFY_LIBRARY_PASS, proceso exit 0: **160 packs / 1.413 archivos materializables / 704 Markdown / 51 perfiles**; composición integral 67/722. Memoria actual: 2.108 lecciones locales + 209 condiciones upstream = 2.317 IDs. La suite Go reconstruida registra cero SKIP. El PASS raíz acredita integridad/composición, no Audit ejecutable integral ni producción.

Se comprobó que elite_whatsapp_v271 era la base descartable creada en esta tarea y que ya no tenía sesiones abiertas. Contenía 183 tenants de fixtures acumulados por las pruebas y, tras el ensayo down/up, cinco observaciones y tres lotes. Se descartó exclusivamente esa base y se detuvo PostgreSQL temporal. Estos datos sintéticos son reconstruibles con los tests; no se eliminaron datos de usuario ni código. Logs y árboles de staging permanecen. Este cierre documental registra los resultados después de ejecutarlos; no se modificó código luego de los gates.
