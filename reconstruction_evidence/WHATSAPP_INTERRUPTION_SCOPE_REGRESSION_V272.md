# V272 — revalidación de ámbito, interrupción SQL y replay durable

Fecha: 2026-09-06. Mantenimiento de biblioteca. Estado `REBUILD_VERIFIED / CONDITIONED`.

## Delta real

WhatsApp 0.8.1 incorpora tres tests Go nuevos (uno con tres subcasos), README y procedencia actualizados. El código de producción StatusObserver queda idéntico byte por byte al de V271. No se incorpora código nuevo de Meta, dependencia, migración, endpoint, cola ni segundo ledger. La tabla de lotes/observaciones ya ofrece deduplicación y atomicidad; añadir un receipt de inbox con la misma finalidad no cerraría la recepción pública pendiente.

## Pruebas ejecutadas

1. `TestStatusObserverPostgresRechecksAfterVerification`: el callback server-side de secretos se ejecuta después de la primera lectura de anchor. Cambia realmente en PostgreSQL la organización del lead, el estado del fence o el hash del receipt. Se exige que la mutación haya tenido éxito y que el segundo control rechace el procesamiento con cero lotes/eventos. Se conserva el verificador Python real aislado; no hay sleep ni sustituto que finja verificar firmas.
2. `TestStatusObserverPostgresCanceledWriteAndRecovery`: una transacción de fixture bloquea la tabla de observaciones. Se demuestra mediante pg_stat_activity que el INSERT del evento está esperando un lock, después de haber ejecutado el INSERT del lote. Se cancela la operación, se comprueba el error y que ni lote ni evento hayan persistido; liberado el lock, el mismo aviso se procesa y queda exactamente un lote/evento.
3. `TestStatusObserverPostgresReplayAcrossConnectionRestart`: tras commit se cierran las conexiones y se crean un pool y valores de servicio nuevos. El replay inserta cero eventos y conserva un lote, un evento y un solo intento outbound. No depende de una cache de deduplicación en memoria.

Control de sensibilidad: sólo en staging se eliminó temporalmente el segundo control de anchor. Los tres subcasos de ámbito fallaron porque el mutante aceptó una observación en cada uno. Se restauró el archivo original y se comprobó su identidad SHA-256 contra V271 antes de continuar. Ese mutante nunca se copió al pack. Es un negativo deliberado, no un fallo del componente publicado.

La suite completa de `internal/whatsappbridge` pasó en staging y otra vez desde el Markdown reconstruido, con Python/BD explícitos y sin SKIP. En el segundo recorrido demoró 9,822 segundos según go test; no es un benchmark ni un SLO. También pasaron `go vet ./...`, `go build ./...` y los 42 tests Python. No se ejecutaron navegador, SCA ni una suite Go/PostgreSQL global en esta revisión.

## Reconstrucción y entorno

67 packs / 722 archivos en destino nuevo; sin cambio del número de archivos materializables. Pack `PYTHON-META-WHATSAPP-CLOUD-ADAPTER` 0.8.1, 28 archivos. Se aplicaron las 51 migraciones a la base local descartable elite_whatsapp_v272. Python 3.14.4, Go 1.26.7 y PostgreSQL 18.6. No se volvió a ejecutar down/up porque no cambió schema; V271 conserva esa evidencia histórica.

Los tres archivos modificados coinciden por SHA-256 entre staging y reconstrucción:

| Archivo | SHA-256 |
|---|---|
| internal/whatsappbridge/status_observer_test.go | 6ed7e3c4c777baf6ca4c031a0d94b9adfe9d712634898ffbe55dbcd971e0db35 |
| whatsapp_cloud/README.md | 598a1b8e213ed73a9382d23afdf2bd0e5a3a99d4c3579eee3db8cf71b616ad00 |
| whatsapp_cloud/PROVENANCE.md | 8f14294d3ca499b844895b120fb78e41e298e8626036b145d4ea7571e39f312d |

Staging bajo temporal del sistema: elite-v272-973c54902270423ca9fb73f514fa62f1. Logs: migrations.log, go-focused.log, mutant.log, go-whatsapp.log, roundtrip-whatsapp.log, roundtrip-python.log, roundtrip-vet.log y roundtrip-build.log. Comando reproducible: README del pack, con ELITE_WHATSAPP_PYTHON y ELITE_WHATSAPP_TEST_DATABASE_URL configurados explícitamente a una BD descartable admitida por approvalPool, seguido de `go test ./internal/whatsappbridge -v -count=1`.

## Procedencia y límites de las autoridades

Los tests son **AUTHORED**, no copiados de una empresa ni certificados por ella. Las fuentes oficiales consultadas el 2026-09-06 gobiernan sólo el método o la semántica indicada:

- [Google SRE — Testing for Reliability](https://sre.google/sre-book/testing-reliability/): pruebas de integración/sistema, regresión y fallos; un PASS no prueba toda la fiabilidad operativa.
- [PostgreSQL 18 — Explicit Locking](https://www.postgresql.org/docs/18/explicit-locking.html): locks de filas/tablas y su duración transaccional.
- [PostgreSQL 18 — Transactions](https://www.postgresql.org/docs/18/tutorial-transactions.html): atomicidad del conjunto de escrituras.

Esto no prueba revocación live de permisos IdP, muerte del proceso/host/servidor PostgreSQL, restore, redelivery Meta, entrega real ni rendimiento bajo carga. No se modificaron fuentes upstream, licencias ni locks. No hay promoción a REUSABLE_PACK por añadir tests.

## Pendientes que no se cierran artificialmente

Montaje de ingress/worker con routing autenticado de tenant/mensaje, retención y acknowledgement/recovery; pantalla del operador y navegador; suscripción/cuenta y callbacks reales. La biblioteca completa conserva los demás frentes del roadmap vigente. Este incremento mejora evidencia de controles existentes, no agrega esas capacidades ni acredita preparación productiva integral.

La recurrencia de inventario/rutas y salidas truncadas de mantenimiento está registrada en FAIL-20260906-368; se corrigió el scope de lectura sin atribuir revisión a contenido omitido.

El primer VERIFY_LIBRARY detectó el snapshot preflight todavía en 704 Markdown; FAIL-20260906-375 / LIB-FAIL-2109 conserva el fallo y la corrección a 705 sin debilitar el check. La repetición terminó con **VERIFY_LIBRARY_PASS, exit 0** (verify-library-recheck.log). Inventario: 160 packs / 1.413 archivos materializables / 705 Markdown / 51 perfiles; memoria 2.109 lecciones locales + 209 condiciones upstream. Es verificación de biblioteca, no un nuevo Audit ejecutable global ni evidencia de accesos productivos.

Limpieza: tras comprobar current_database=elite_whatsapp_v272 y cero sesiones activas, se eliminó exclusivamente esa base descartable (188 tenants sintéticos, 15 lotes, 19 observaciones, incluidos fixtures del mutante) y se detuvo el servidor local iniciado para la prueba. Los fixtures permiten recrearla; logs y staging se conservan. No se eliminaron datos de usuario ni se dejó un worker ejecutándose.
