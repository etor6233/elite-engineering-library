# Runtime IA: generación, presupuesto e identidad — V402 / execution314

PROVEN_LOCAL_DELTA dentro de T2807; el control completo sigue abierto para
evaluación conectada y consolidación de LIB-R10. No READY global ni producción.
Biblioteca205packs/2275blocks; franquicia113packs/1523files.13archivos AUTHORED
de glue/tests:8actualizados y5nuevos;1510previos conservan bytes. Sin dependencia
nueva, atribución empresarial ni SCA repetido sobre el mismo grafo.

La regresión original falló en seis aserciones: Complete y FailRetryable aceptaban
leases vencidos o generaciones anteriores; replay e historial devolvían textos
fuera de retención. El store ahora devuelve generación y exige el mismo intento,
estado y lease vivo según reloj PostgreSQL en las transiciones. READ COMMITTED
con lock por fila evita tratar un conflicto serializable transitorio como éxito.
Retención de lecturas está probada; purga operativa sigue en T2809.

El runtime reserva tokens en PostgreSQL antes del request. El cupo es acumulado
por tenant: reiniciar no lo repone, cambiarlo por configuración se rechaza y un
reintento incierto consume otra reserva. La reserva se recupera por generación.
Dos mensajes disputando el saldo final producen un ganador y un rechazo. Es una
cota conservadora basada en bytes y dos calls, no un precio/factura del proveedor.
El helper economy histórico no gobierna ya el presupuesto de este runtime.

La primera tool, sus argumentos canónicos, contacto y configuración se fijan con
un hash inmutable. Un reintento que cambie esa intención termina en handoff;
uno vencido no puede fijarla. El owner de dominio sigue imponiendo idempotencia
del efecto; la duración de una llamada remota no se presenta como una transacción
con el lease local. Respuesta final y argumentos tienen límites y schemas.

La continuación Responses reenvía instrucciones explícitas. La regresión de
wire falló antes del cambio y pasa después; la [referencia oficial OpenAI](https://developers.openai.com/api/reference/cli/resources/responses/methods/create)
confirma que previous_response_id no las hereda. Modelo y proveedor siguen siendo
selecciones del consumidor; fixtures no prueban calidad de un modelo live.

PG18.6 real, suites afectadas y compilación integral PASS.0086 down/up vacío PASS;
down con datos está protegido. Dos perfiles finales reconstruyen exactamente.
13archivos del delta recibieron lint adaptado:9avisos revisados con SHA/posición,
cero nuevos blockers o avisos sin revisar. No reemplaza evaluación de seguridad.
El finding fiscal1539 permanece para ARCA; no se investiga Daybreak/libxml2.

FAIL932 corregido para este delta; RED conservado. FAIL933: una sustitución de
firma duplicó un argumento sólo en el test de retry; compilación lo detectó y
la corrección puntual pasó. Próximo: cerrar evals/runtime de referencia y reutilizar
el pack histórico ya admitido V373; roadmap807 no era el requisito histórico,
la autoridad correcta es LIB-R10 / PROJECT_HISTORY_MODEL_TRAINING_CONTRACT.md.
