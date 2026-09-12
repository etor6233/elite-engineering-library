# V277 — worker scoped y cierre durable de statuses WhatsApp

2026-09-06. Mantenimiento de biblioteca. WhatsApp 0.11.0 y async workers 0.3.0,
REBUILD_VERIFIED / CONDITIONED; sin promoción a REUSABLE_PACK ni producción.
Objetivo: componentes materializables y conectados para crear una franquicia o
extender otro proyecto, con procedencia y evidencia verificables.

## Delta implementado

```text
HTTP firmado → inbox/job existente → ClaimScoped
→ router/receipt autorizado → observer idempotente
→ transacción [job completado + inbox procesado + audit]
                       fallo → retry acotado o terminal preservado
```

NewStatusWorker configura router, worker ID, lease y retry; ProcessOnce recibe
un principal realmente autenticado con permisos y organización actuales. No
inventa identidad desde payload. Selecciona como máximo un job de tenant,
conexión, proveedor, queue, tipo, schema y event_type configurados. El contrato
del payload exige exactamente cuatro campos string y un ID wa:SHA256 válido.
No toma trabajo de otra integración ni inicia un daemon automáticamente.

ClaimScoped y CompleteJobInTx/FailJobInTx amplían el owner jobs.go; no hay otra
cola, store o implementación paralela de las transiciones. Los helpers reusan
row lock, generación attempts y comprobación de reloj existentes, añadiendo
comparación exacta del snapshot de payload/queue/type/schema. El commit es del
caller: completion, inbox processed y audit son atómicos. El callback privado
del router permite ese cierre; ObserveRetained público sigue sin hacer ACK.

StatusObserver conserva commits independientes. Si el cierre posterior falla,
las observaciones ya durables se reproducen sin duplicación; no se reenvía un
mensaje. Fallos recuperables conservan el inbox recibido y reprograman hasta
max_attempts. Fallos terminales conservan job/inbox/audit; un último intento
caído y vencido se marca LEASE_EXHAUSTED sin ejecutar un efecto nuevo. No se
resetean generaciones ni existe redrive automático de terminales.

Audit utiliza audit.event append-only existente: actor real, job scoped,
attempt, terminal, reason acotado y hash de provider_event_id. No guarda raw
mensajes, contactos, tokens ni secretos en el nuevo registro. FailureRecorded
false significa que no se demostró registro durable (claim perdido, payload
alterado o DB indisponible): el host debe comunicar error seguro y reconciliar,
nunca interpretar inserciones parciales o cero como éxito. No se afirma que
cualquier fallo puede guardarse cuando la base no está disponible.

## Autoridades y procedencia

Los dos archivos nuevos y los cambios locales son AUTHORED. No son código
publicado por Meta/AWS/PostgreSQL. Reutilizan la adaptación Meta ya existente,
router, observer, approval/fence, cola y auditoría. No cambian dependencias,
migraciones, protocolo, source lock, licencia ni archivos upstream.

- [Amazon Builders' Library — elección de líder](https://aws.amazon.com/es/builders-library/leader-election-in-distributed-systems/), consultado 2026-09-06: revisar leases y considerar pausas/fallos. Es autoridad de método, no autoría ni certificación de este worker.
- [PostgreSQL 18 — transacciones](https://www.postgresql.org/docs/18/tutorial-transactions.html), consultado 2026-09-06: atomicidad local y visibilidad del commit. No demuestra exactly-once remoto.

## Verificación real

Windows, Go 1.26.7, Python 3.14.4, PostgreSQL 18.6, loopback 55959. Dos bases
dedicadas elite_whatsapp_v277 y elite_jobs_v277, 52 migraciones aplicadas a cada
una. Sólo fixtures sintéticos; sin cuenta/proveedor/envío real.

- Nueve tests nuevos del worker, ejecutados dos veces en staging: aislamiento
  de siete clases de jobs ajenos; una finalización y audit; retry recuperable y
  terminal; rollback entre escrituras del cierre; recuperación de intento caído
  y cuarentena del último; consumidores concurrentes; claim tomado por otro
  worker después de observar; cinco rechazos preclaim; drift del payload; router
  sin inicializar rechazado.
- El fallo inducido entre UPDATE job y UPDATE inbox revierte ambos; queda una
  observación previamente committed. Replay cierra con cero nuevas observaciones.
- Claim takeover y drift después de observar impiden completion y failure write
  antiguos: FailureRecorded=false, snapshot vigente intacto. El takeover se
  recupera con nuevo claim sin duplicar la observación.
- Reconstrucción en destino ausente: 67 packs / 730 archivos, seis hashes
  idénticos respecto al staging probado.
- Artefacto reconstruido: 34 tests superiores Go/WhatsApp PASS sin SKIP; siete
  tests superiores de jobs/inbox/outbox/processors PASS sin SKIP; 42 Python PASS.
  go vet ./... y go build ./... PASS. No full go test ./..., Audit ejecutable
  integral, SCA nuevo, navegador, carga, ofensiva, provider live ni deploy.

FAIL-20260906-381 conserva worker-green.log con SQLSTATE 23514: un UUID de
fixture podía comenzar con número y violar tenant_code. Prefijo scope- corrige
el fixture, no la constraint. Regresión repetida y reconstruida PASS. La primera
composición también rechazó versión Journey modificada por un hunk sin packId;
se restauró 0.10.0 y se actualizó WhatsApp 0.11.0 con contexto explícito. Fallo
guardado en composition.log; composición posterior verde sin relajar gates.

## Reproducción y archivos

Componer FRANCHISE_COMPLETE_PACK_PLAN.md a destino ausente, aplicar las 52
migraciones a bases descartables con los nombres/prefijos y puerto de fixture.
Definir ELITE_WHATSAPP_TEST_DATABASE_URL, ELITE_WHATSAPP_PYTHON y
TEST_DATABASE_URL. Ejecutar go test ./internal/whatsappbridge -count=1 -v;
go test ./internal/platform/postgres ./internal/platform/workers -run
'TestJobs|TestJobProcessor|TestProcessors|TestOutbox|TestInbox' -count=1 -v;
en whatsapp_cloud ejecutar Python unittest test_whatsapp_cloud.py y
test_status_reconciliation.py; luego go vet ./... y go build ./....

Temporal conservado: elite-v277-8918729c80eb485daae1b08a3fd4863f, árboles work y
roundtrip. Logs worker-first.log, worker-green.log, worker-regression.log,
jobs-green.log, composition.log, composition-green.log, roundtrip-go.log,
roundtrip-jobs.log, python.log, vet.log, build.log y migraciones por base.
compile-unit.log fue sólo compilación para whatsappbridge (no tests to run),
no evidencia de escenarios. Las pruebas completas posteriores sí los ejecutan.

| Archivo | SHA-256 |
|---|---|
| internal/platform/postgres/jobs.go | 417bff0d640ea13146a873e1e556e61cc428d507ab714ef10ef4d073ec8f8345 |
| internal/whatsappbridge/status_router.go | 75d79c1cf44133276a4cdf3d8fa067db54d8b8bd6b1c4600b61e2b34d8c107e5 |
| internal/whatsappbridge/status_worker.go | a94051c187923cb19a10807369c28f78dea7c3e676d2da5bd2659ec85e7b5ffa |
| internal/whatsappbridge/status_worker_test.go | 917692f3144e005b1e3c5968cbc1f9ef5c0657963ef3d2b5edb58cb13d3c2128 |
| whatsapp_cloud/README.md | fd73cb1ff105ab4d0bbefda90504dadb42f5dff364d20978c840b60ea90ed8f7 |
| whatsapp_cloud/PROVENANCE.md | 742e0f1d4d32a584138feea26bed30c4ed5d05cf43bef5b4bd036bdc7f86de1c |

## Alcance que permanece abierto

Este bloque implementa procesamiento de statuses, no mensajes de cliente ni
autorespuestas universales. El host debe aportar scheduling con backoff,
shutdown, identidad legítima, capacidad, alertas, política de retención y
reconciliación/redrive autorizado. Hay terminales en tablas existentes, no una
DLQ externa desplegada. No se probó matar el proceso/servidor/host: se inyectaron
estados vencidos y fallos transaccionales en PostgreSQL real.

Frontend conectado, ayuda/capacitación/soporte versionados, telemetría segura,
otros journeys y mutaciones de proveedores siguen en FRANCHISE_PREFLIGHT_GAP.md.
Cuenta/contratos/consentimiento/entrega live, IdP/CDN/WAF, carga, seguridad,
recovery/deploy y aceptación siguen siendo evidencia del proyecto. No se
certifica READY_TO_BUILD ni checkpoint de un proyecto productivo con este PASS.

Rollback: pausar host, preservar jobs/inbox/observaciones/audit/fences y corregir
el owner canónico antes de reconstruir. Nunca activar workers pre-0.2.0 ni
borrar/reiniciar terminales para esconder fallos. No porcentaje global estimado.

## Control general y limpieza

VERIFY_LIBRARY_PASS y proceso exit 0. Inventario vigente: 160 packs / 1.421
archivos materializables / 710 Markdown / 51 perfiles; franquicia 67/730.
AUTHORED=1183, ADAPTED=133, VERBATIM=105. Memoria 2.115 IDs locales + 209
upstream = 2.324. Log verify-library.log; no Audit ejecutable global repetido.

Se verificaron current_database y cero sesiones ajenas en ambas bases:
elite_whatsapp_v277 tenía 182 tenants sintéticos, 115 jobs y 87 webhooks;
elite_jobs_v277 tenía 12 tenants, 12 jobs y cero webhooks. Se eliminaron sólo
esas dos bases descartables y se detuvo el PostgreSQL iniciado para las pruebas.
Los datos sintéticos pueden recrearse con fixtures/migraciones; árboles y logs
se conservan. No se borraron datos del usuario ni se activaron proveedores.
