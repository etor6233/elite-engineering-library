# V270 — envío anclado → webhook firmado → observaciones correlacionadas

Fecha: 2026-09-06. Mantenimiento de biblioteca. Alcance `REBUILD_VERIFIED / CONDITIONED`: correlación de evidencia de WhatsApp, no entrega live, consumer PostgreSQL ni interfaz terminada.

## Avance conectado

El mismo send_bridge existente produce SEND_RECEIPT.json y devuelve su evidence_sha256; el fence Go ya guarda ese hash. La nueva función reconcile_status_webhooks exige los bytes del comprobante y ese anchor confiable, valida perfil/teléfono/versión y usa el mismo normalizador HMAC/scope existente para verificar cada webhook original. No acepta un normalized-events.json arbitrario como prueba de autenticidad.

Relaciona exclusivamente message_id y recipient del comprobante. Otros mensajes de un lote se excluyen; el mismo mensaje con otro destinatario rechaza. Las repeticiones exactas deduplican, pero un event key repetido con distinto payload rechaza en vez de descartar la contradicción. WABA/teléfono ajenos, firma inválida, anchor alterado o un lote posterior inválido no publican resultados parciales.

Genera STATUS_RECONCILIATION.json y status-observations.json juntos en un destino ausente. Ordena observaciones por timestamp del proveedor, no por llegada. Un empate del timestamp más reciente con estados distintos queda AMBIGUOUS_LATEST_TIMESTAMP; no inventa precedencia. Retiene todos los estados observados y su evidencia hash-linked. NOT_OBSERVED_IN_INPUT sólo habla de los lotes aportados, no del estado global. No realiza llamadas HTTP, no reenvía y no actualiza SQL ni reconoce un inbox como procesado.

## Procedencia y autoridad

- Nuevos status_reconciliation.py y test_status_reconciliation.py: **AUTHORED**, no código copiado de Meta. Presupuesto, anchor, snapshots, deduplicación y tratamiento de empates son decisiones locales declaradas, verificadas por tests.
- whatsapp_cloud.py permanece adaptación declarada; su frontera HMAC/scope/parser se extrae a normalize_verified_webhook y es compartida por la función anterior y la nueva. No se crea otro verificador de firmas. Evidence v2 mantiene su contrato.
- [Meta — Message Status Update Notifications](https://www.postman.com/meta/whatsapp-business-platform/request/rgtfq23/message-status-update-notifications), consultada 2026-09-06, documenta ID de mensaje/destinatario, status/timestamp y posible llegada fuera de orden. Esa fuente gobierna ese claim estrecho, no el módulo completo ni su readiness.
- Source lock Meta fbsamples de70ee908a67026e642aaee3703d20464e2a9466, upstreams, licencia/notices y dependencias no cambian. La verificación de bytes oficiales incluida en send_bridge se ejecuta en las pruebas nuevas. No se afirma un SCA general nuevo ni una revisión global de todos los upstreams.

Un enlace Postman secundario no pudo abrirse y no se utilizó como evidencia; la página oficial principal sí se verificó. Los fallos de inventario/cwd y el patch inicial sobre un alias de red incorrecto se documentan en FAIL-20260906-373 / LIB-FAIL-2107.

## Pruebas ejecutadas

**41 tests Python PASS** en staging y repetidos desde reconstrucción limpia: los 22 anteriores y 19 nuevos, con subcasos. Incluyen:

- send_bridge real → archivos de receipt → anchor retornado → webhook firmado → snapshot; transporte proveedor sintético, una llamada de envío del fixture y cero llamadas durante correlación;
- los seis órdenes posibles de sent/delivered/read producen el mismo hash de observaciones;
- duplicación de evento/lote, timestamp empatado, message ajeno, lotes mixtos y mismo message con recipient distinto;
- WABA/teléfono ajenos, firma alterada incluso en eventos no relacionados, anchor alterado, contrato/perfil inválidos;
- event key idéntico con evidencia divergente, failed/deleted conservados sin autorización de retry;
- lotes vacíos/excesivos/malformados, cuerpo mayor de un MiB, receipt mayor de 64 KiB;
- destino existente intacto, fallo de escritura sin snapshot parcial y lote inválido tardío sin resultados parciales;
- no persiste texto crudo, IDs, teléfono o secretos del fixture en los nuevos artifacts. Los hashes son pseudónimos, no anonimización.

No se usó cuenta Meta real. No se ejecutaron nuevos tests navegador/IdP/PostgreSQL ni se convirtió la evidencia histórica V263–V268 en PASS actual. La nueva función no depende de PostgreSQL y todavía no escribe su estado: eso sigue siendo un trabajo de integración explícito.

Desde la reconstrucción también pasan TestProcessBridgeFencedReplayAndAmbiguity, TestBridgeRejectsUnapprovedOrDriftedRequestBeforeChild y TestRequestContractIsExact con ELITE_WHATSAPP_PYTHON explícito; go vet ./... y go build ./... PASS, Go 1.26.7. Son tres tests focales con subcasos y proceso Python local, no una suite Go/PG completa. Logs roundtrip-go-bridge.log, roundtrip-vet.log y roundtrip-go-build.log.

## Materialización y reproducción

PYTHON-META-WHATSAPP-CLOUD-ADAPTER **0.7.0**, 24 archivos; dos archivos nuevos y tres existentes modificados. Perfil integral **67 packs / 718 archivos**. Sin migraciones, frameworks, dependencias ni cuentas nuevas. Los planes integral y WhatsApp seleccionan la misma versión.

Los cinco archivos coinciden por SHA-256 entre staging probado y reconstrucción nueva:

| Archivo | SHA-256 |
|---|---|
| whatsapp_cloud/whatsapp_cloud.py | 10cf0e71d514b1168a274243522903616b8c3d44cfe649c2dc4165fdd8b8a037 |
| whatsapp_cloud/status_reconciliation.py | 837291ce820e36c174073e99b288caadf620b7a5bcd04f6a5610c5d4188f28bd |
| whatsapp_cloud/test_status_reconciliation.py | 5fcd380318f20b29f41a4ebb2d6b0f0fbb434f76d3d04a2638a88c2a129df932 |
| whatsapp_cloud/README.md | c816307710c990578c5daa51a7a3707ada5f18a64b7913bf189c00bfbb3034a6 |
| whatsapp_cloud/PROVENANCE.md | 57143ce084006981583c67e89cf0e7f09b39e96d251995ad9b4c578d773d0d5f |

Ejecutar `python -m unittest discover -v` dentro de whatsapp_cloud materializado. El README incluye la llamada inmediata, entradas, outputs, límites y condiciones. Python observado: 3.14.4. Staging conservado bajo el temporal del sistema: elite-v270-da6ae0315faa4757a75bd50068c125db; logs python-tests-final.log y roundtrip-python.log. Los fixtures se limpian por TemporaryDirectory; no se borran datos de proyecto ni se activan procesos recurrentes.

## Distancia al objetivo completo

El siguiente tramo de este recorrido sigue siendo: entrada durable autorizada → correlación con el owner/fence → persistencia/lectura scoped de observaciones → interfaz de operador y recuperación. El nuevo snapshot reduce esa brecha, **no la cierra entera**. Un send unknown sin SEND_RECEIPT anclado no puede resolverse con esta función; hay que conservar su incertidumbre y adquirir evidencia admisible por la ruta correspondiente, nunca fabricar un anchor o reenviar.

Permanecen los seis frentes de FRANCHISE_PREFLIGHT_GAP.md: recorridos empresariales completos, ayuda/capacitación/soporte por release, telemetría sin PII, bootstrap compacto, integraciones/mutaciones pendientes y promoción por claim. Son frentes, no seis tareas iguales ni un porcentaje de esfuerzo. Cuentas, accesos, reglas/corpus y gates de operación del proyecto son otra clase de evidencia. No hay base para prometer 90/95/100 % ni una franquicia productiva en una semana.

## Resultado del control general

VERIFY_LIBRARY_PASS, proceso exit 0: 160 packs, 1.409 archivos materializables, 703 Markdown y 51 perfiles; composición integral 67/718. Inventarios actuales y memoria sincronizados: 2.107 lecciones locales + 209 condiciones upstream = 2.316 IDs. Este PASS de integridad/composición no sustituye Audit ejecutable integral, gates del entorno ni condiciones upstream abiertas. Log: verify-library.log en el staging indicado. Este párrafo registra el resultado después de ejecutarlo; no hubo cambios de código posteriores.
