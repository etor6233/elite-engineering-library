# V275 — generación y vigencia del claim de jobs

Fecha: 2026-09-06. Mantenimiento de biblioteca, no implementación de una franquicia target. GO-RELIABLE-ASYNC-WORKERS 0.2.0 mantiene REBUILD_VERIFIED / CONDITIONED. No se añade upstream, cola, tabla, dependencia, servicio pago ni worker en segundo plano.

## Objetivo y delta

El objetivo sigue siendo una biblioteca reutilizable para construir una franquicia o ampliar otro sistema con componentes compatibles, procedencia verificable y recorridos de extremo a extremo. Este cambio protege el owner de jobs existente antes de conectar el procesamiento automático de avisos; no sustituye el roadmap FRANCHISE_PREFLIGHT_GAP.md.

Antes, Complete y Fail aceptaban un worker con lease vencido; el mismo nombre reutilizado podía cerrar otro intento. Ahora exigen Attempts exacto y vigencia comprobada con reloj PostgreSQL después del row lock. JobProcessor revalida cada elemento del batch y entrega al handler un contexto con el presupuesto restante menos el round trip de consulta. No confía en un check hecho al principio de todo el batch.

API incompatible intencional: Complete/Fail requieren generación, JobStore requiere RemainingLease. Los cuatro archivos y los tres perfiles consumidores se actualizan juntos. Drenar/parar versiones antiguas antes de desplegar; no reiniciar attempts, no reutilizar IDs ni volver al SQL 0.1.x como rollback seguro. No hay cambio de schema. Para retroceder, detener el consumidor y preservar la cola hasta disponer de un artefacto que mantenga estos controles.

## Autoridades y procedencia

Todo el delta es AUTHORED. No se presenta como código AWS ni como superior por fama.

- [AWS Builders' Library — elección de líder](https://aws.amazon.com/es/builders-library/leader-election-in-distributed-systems/), consultado 2026-09-06: revisar concesiones antes de efectos y considerar demoras/pausas. El enlace inglés redirigió a una página sin texto extraíble; se consultó esta edición oficial, no un agregador.
- [PostgreSQL 18 — SELECT](https://www.postgresql.org/docs/18/sql-select.html), consultado 2026-09-06: row locks y SKIP LOCKED para consumidores de colas. El SQL de integración y la selección de generación son responsabilidad local, no código de ejemplo atribuido a PostgreSQL o AWS.

## Pruebas ejecutadas

Windows, Go 1.26.7, PostgreSQL 18.6 loopback, base descartable elite_jobs_v275 con las 51 migraciones existentes. No se usaron datos, cuentas, pagos, envíos ni proveedores reales.

1. El primer red.log falló antes del escenario por SQLSTATE 42P08 del fixture; NO prueba el defecto. FAIL-20260906-379 conserva causa y corrección: parámetros separados para UUID y queue text.
2. red-semantic.log reproduce cuatro fallos reales: complete/fail después de vencimiento, con/sin takeover del mismo worker. Cada uno informa stale execution changed job state.
3. Los cuatro casos pasan después del fix; además comprueban que el intento vigente conserva su estado y puede completar/reprogramar.
4. Dos casos PostgreSQL confirman bloqueo real mediante pg_stat_activity, mantienen el row lock hasta vencimiento sin cambiar versión de fila y prueban rechazo de Complete/Fail al desbloquear. No se sustituye concurrencia por un mock ni por un sleep estimado del reloj del cliente.
5. Cinco subcasos del processor: claim perdido, vencido, segundo elemento del batch vencido, presupuesto transmitido y contexto cancelado. Un claim inválido no llama al handler ni confirma/reprograma.
6. Reconstrucción nueva desde Markdown: siete tests de nivel superior PASS, once subcasos nuevos PASS, cero SKIP en la selección; incluye regresiones existentes de inbox, outbox, retry/terminal y processors. El test histórico outbox no demuestra fencing por generación de ese subsistema.
7. go vet ./... y go build ./... PASS sobre el perfil reconstruido. No se ejecutó full go test ./..., SCA global, navegador, carga, seguridad ofensiva, recovery productivo ni cuenta Meta.

## Reproducción

Componer FRANCHISE_COMPLETE_PACK_PLAN.md a destino ausente. Aplicar las 51 migraciones a una base local descartable con prefijo elite_jobs_; definir TEST_DATABASE_URL. Ejecutar:

```text
go test ./internal/platform/postgres ./internal/platform/workers -run 'TestJobs|TestJobProcessor|TestProcessors|TestOutbox|TestInbox' -count=1 -v
go vet ./...
go build ./...
```

Temporal conservado: elite-v275-6aad65d689fd44269013e5ad785b5e21; logs migrations.log, red.log, red-semantic.log, green.log, green-all.log, roundtrip-tests.log, vet.log, build.log. Materialización: 67 packs/724 archivos. Cuatro archivos idénticos entre staging y reconstrucción:

| Archivo | SHA-256 |
|---|---|
| internal/platform/postgres/jobs.go | 23698fe4470e0b339eb2bff4045ee4aa8c60f3882a86eff90fa35d9802e7f2d5 |
| internal/platform/postgres/jobs_integration_test.go | 97ffdd220c22757d9ac86994eebea4870001f5ab5af9196892eaad50c9a207aa |
| internal/platform/workers/processors.go | a4eb2fca998da8c66e9906e70c74de98fa6998ff171e055ec02e7ced56ad282e |
| internal/platform/workers/processors_test.go | 834d943386aa4a25867a7e590684055582d73fba532c535c2abceb0e3f567985 |

## Límites y continuación

Un lease no es autorización empresarial ni puede cancelar un efecto remoto ya emitido. Los handlers deben respetar contexto, idempotencia, fence del efecto y reconciliación. No se añade heartbeat automático. JobProcessor es infraestructura de confianza; no selecciona tenant por petición pública.

Todavía falta el routing/worker WhatsApp que resuelva mensajes contra receipts autorizados y complete el inbox sólo tras procesar todos los eventos admitidos, además de retry/DLQ/retención, UI y pruebas live. Antes de activar otros consumidores, revisar también el fencing del publisher/outbox y la reconciliación de jobs cuyo último intento murió; esta prueba no declara resueltos esos casos. Los demás frentes del roadmap siguen vigentes. No se estima porcentaje de cierre por contar archivos.

## Verificación estructural y limpieza realizadas

VERIFY_LIBRARY_PASS y proceso exit 0: 160 packs / 1.415 archivos materializables / 708 Markdown / 51 perfiles; composición franquicia 67/724. Memoria: 2.113 IDs locales + 209 upstream. No se presenta este control estructural como nuevo Audit ejecutable global ni como READY_TO_BUILD del proyecto.

Después de verificar current_database=elite_jobs_v275, 24 tenants sintéticos, 20 jobs de fixture y cero sesiones ajenas, se eliminó exclusivamente esa base descartable y se detuvo el PostgreSQL iniciado para estas pruebas. Migraciones y fixtures permiten recrearla; staging, reconstrucción y logs se conservan. No se borraron datos de usuario ni se activaron workers/servicios externos.
