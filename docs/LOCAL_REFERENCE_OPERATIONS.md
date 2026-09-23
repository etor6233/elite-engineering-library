# Operación de la referencia local

AUTHORED glue alrededor de Go/OTel/Prometheus y del supervisor Windows existentes.
Scope LIBRARY_INFRASTRUCTURE_LOCAL_FIXTURES, un operador local de confianza.

## Arranque y ensayo

Primero construir la referencia y aplicar ci/migrate_local_reference.py en una
base vacía propia. El contrato del ensayo espera las 84 migraciones de esta
composición, no una base productiva. Su configuración JSON contiene solamente
database_url del patrón elite_payment_connected_<32hex>, loopback sin contraseña,
y data_path de ese clúster detenido. No acepta cuentas ni proveedores live.

Ejecutar con Python fijado:

```text
python -X utf8 -B ci/run_observed_reference.py --artifact <artefacto> --inventory <inventario> --inventory-sha256 <sha> --runtime-lock <lock> --runtime-lock-sha256 <sha> --database-config <fixture.json> --database-config-sha256 <sha> --rules ops/prometheus/platform.rules.yml --rules-sha256 <sha> --work <directorio-ausente>
```

El lock externo fija postgres/pg_ctl/psql, prometheus/promtool y
telemetry-reference. Sus identidades/licencias se conservan en los expedientes
V372/V375; nunca se descargan versiones móviles desde este comando. El artefacto
fija API, Node y Next. La identidad RSA efímera sólo vive en el fixture loopback.

El ensayo arranca API+Next bajo el Job Object ya admitido: doce procesos como
máximo, 1 GiB de commit, salida capturada de 2 MiB, 1200 segundos con métricas y
25 segundos de cierre. Fuera de este ensayo, el controlador conserva 600 segundos
y el entrypoint de uso local sigue limitado a 480. Un puerto de métricas opcional
debe ser entero, loopback y distinto de los puertos API/web. No hay SCM ni HA.

## Señal, alerta y recuperación

Sólo se exporta el histograma HTTP de método/status, con cardinalidad limitada y
sin query, ruta, tenant, tokens, cookies, trazas ni exemplars. Se prueba contra
el handler real, su OIDC, PostgreSQL y un pedido durable. Prometheus consulta con
mTLS efímero y retiene WAL/series entre reinicios. La regla HTTP conserva ventana
5m, umbral 1%, hold 10m y evaluación 30s. No se inyectan muestras ni se acelera
el reloj. El propio ensayo renombra y restaura exclusivamente la tabla de su
base sintética; no copiar esa operación a una base real.

Ante PlatformHighErrorRate: consultar progress.json/result.json, comprobar salud
del API y estado del clúster propietario; restaurar la tabla sólo si el receipt
declara fault_restore_failure. El finally intenta restaurar antes del apagado.
En operación fuera del ensayo, no modificar tablas: recuperar la última release
confirmada con Deployment.recover y conservar el fallo. Los runbook_url históricos
replace.invalid no son enlaces de operación: esta guía gobierna el caso local.
OutboxBacklog y BackupStale no tienen señales conectadas por este cambio y no se
presentan como alertas probadas del host.

Después del firing real: 1000 replay requests, concurrencia ocho, identidad y
pedido constantes; se observan p95/máximo y se exige la misma fila durable.
La reparación limpia la alerta; el reinicio de Prometheus conserva una consulta
con timestamp fijo y el recovery del controlador conserva el pedido. Cada
instancia termina con recibo del árbol vacío, Go0/NextSIGTERM143/wrapper0.

## Retención y límites del claim

Prometheus: 1h y 64MB para bloques persistentes. Eso NO es una cuota dura de disco:
WAL, head y compacción pueden excederla. La documentación oficial lo especifica:
https://prometheus.io/docs/prometheus/latest/storage/ . Se registra tamaño real y
replay de WAL; no se infiere limpieza de bloques expirados de una ejecución de
menos de una hora. La captura del proceso sí tiene un límite duro de 2 MiB.

El owner WINDOWS_REFERENCE_TELEMETRY_RUNTIME conserva el ensayo separado V372:
rotación de archivo del Collector a 1MiB, dos backups, un día, 2000 observaciones
y reinicio. Se conservan código/hashes/binarios de esa evidencia, sin atribuirle
logs del API nuevo. El despliegue productivo permanente, routing de notificaciones,
respondedor humano, SLOs comerciales, cuotas target, PITR y fallos de host quedan
CANDIDATE_TARGET, no CONDITIONED sólo por una credencial. El límite de T2809
admitido aquí es el host y ensayo local finito explícito.
