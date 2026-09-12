# V274 — HTTP Meta-adaptado → inbox/job PostgreSQL → original verificable

Fecha: 2026-09-06. Mantenimiento de Elite Engineering Library. Pack WhatsApp 0.9.0, estado REBUILD_VERIFIED / CONDITIONED. Sin cuenta, envío ni despliegue externo.

## Objetivo al que aporta

Completar una biblioteca reutilizable para que un agente construya una franquicia o amplíe un sistema existente mediante componentes compatibles, fuentes/procedencia verificables y journeys de extremo a extremo demostrados. El roadmap sigue en FRANCHISE_PREFLIGHT_GAP.md: recorridos empresariales conectados, capacitación/soporte, telemetría, bootstrap compacto, integraciones elegidas y promoción por claim. Este expediente no reemplaza ese roadmap ni convierte la biblioteca en una franquicia productiva.

La recepción pendiente ahora tiene un tramo implementado y reconstruido: HTTP real -> proceso Python Meta-adaptado -> transacción del inbox/job existente. No se añadió una segunda cola. Todavía no es un worker autónomo de routing ni UI.

## Implementación entregada

- NewWebhookReceiver recibe configuración server-side ligada a tenant/conexión, copia el perfil, exige un hash de expediente de retención, secretos por interfaces y concurrencia explícita 1..16. No se registra ni ejecuta automáticamente al materializar.
- GET utiliza verify_subscription existente y devuelve un challenge numérico acotado; la comprobación del token permite configurar la suscripción antes de declarar el perfil POST demostrado. No prueba posesión de una cuenta Meta por sí sola.
- POST acepta JSON sin Content-Encoding, <=1 MiB, una sola firma X-Hub-Signature-256 y sin query que pueda escoger tenant. El proceso tiene un límite de cinco segundos, input/output acotados, -I/-B, ejecutable/script fijados por hash, ambiente mínimo y secretos por stdin; se mantienen necesarios los timeouts de servidor y protecciones del edge.
- ingress_bridge reutiliza el verificador/normalizador existente: firma sobre bytes exactos, WABA/teléfono y perfil aprobado, hasta 1.000 eventos. La salida del proceso se liga por hash al frame. No hay un nuevo algoritmo de firma Meta escrito en Go.
- Sólo después de AcceptWebhook del core 0.1.2 retorna correctamente se emite 200 EVENT_RECEIVED. Se conserva estado received y un job provider.webhook.received pendiente, no processed ni completado. Errores de almacenamiento/conexión deshabilitada o saturación devuelven 503 sin éxito falso; no se hacen envíos ni retries externos.
- Receipt determinista en integration.webhook_event.payload: schema elite-whatsapp-retained-webhook/v1, bytes originales base64, firma, hash del perfil/política y cantidad verificada. La identidad wa:SHA256(raw) se limita por tenant/conexión. El repo compara receipt completo y crea evento/job atómicamente.

**Cambio explícito de retención:** el writer offline sigue produciendo sólo evidencia pseudónima. El nuevo receptor opt-in conserva datos crudos potencialmente sensibles. Base64 NO es cifrado. El hash del expediente es un enlace que el agente debe demostrar en el proyecto, no una aprobación generada por este código. BD/backups cifrados, acceso mínimo, redacción de logs, plazos y borrado deben probarse antes de exponerlo. No se almacenan app secret ni verify token. Cambios de firma/perfil/política sobre los mismos bytes pueden generar conflicto de replay y necesitan reconciliación explícita, sin sobreescribir historial.

## Evidencia real ejecutada

Suite completa internal/whatsappbridge con ELITE_WHATSAPP_PYTHON y ELITE_WHATSAPP_TEST_DATABASE_URL, sin SKIP, en staging y desde reconstrucción limpia. Dos nuevos tests Go (uno con 14 subcasos), usando el proceso Python real y PostgreSQL real:

1. GET challenge HTTP correcto.
2. Se recibe un POST firmado, se compromete el receipt y se corta la conexión antes de entregar el ACK al cliente. Se comprueba una fila durable tras el error de transporte.
3. Cuatro reintentos concurrentes retornan 200; queda un evento y un job pendientes. No se acumulan cuatro trabajos.
4. Se relee de PostgreSQL el payload y se comprueba igualdad exacta de bytes/firma, hashes de perfil/retención y count. No contiene secretos. Estos datos siguen siendo sensibles aunque el transporte al proceso esté aislado.
5. Con routing explícito del fixture, se pasan los bytes releídos a StatusObserver existente; inserta una observación y la lectura de estado scoped devuelve observed_delivered. Ese routing no es un worker implementado ni una consulta de navegador.
6. La conexión disabled rechaza el replay HTTP sin añadir evento/job.
7. Catorce rechazos: firma inválida, WABA ajena firmada, firma duplicada, media type, encoding, sobrepresupuesto, query, JSON inválido firmado, hash de runtime, concurrencia ocupada, verify token incorrecto, challenge duplicado, query malformada y método. El fixture termina sin receipts de esos rechazos. Constructores sin aprobación de retención o presupuesto también se rechazan.

42 tests Python existentes PASS desde ambos árboles; las pruebas Go nuevas ejercitan el nuevo modo Python en proceso aislado. go vet ./... y go build ./... PASS. No hay nuevo PASS de navegador, SCA, carga, seguridad ofensiva, restore o prueba Meta live. El package Go reconstruido demoró 10,174 s en esa ejecución; no es un SLO.

## Reproducción y hashes

Perfil completo 67 packs / 724 archivos, WhatsApp 0.9.0 con 30 archivos. Dos archivos materializables nuevos AUTHORED. Sin nuevas dependencias ni migraciones: las 51 existentes se aplicaron a elite_whatsapp_v274. Go 1.26.7, CPython 3.14.4, PostgreSQL 18.6, Windows. No se repitió down/up porque no cambia schema.

Cinco archivos idénticos SHA-256 entre staging y reconstrucción:

| Archivo | SHA-256 |
|---|---|
| internal/whatsappbridge/webhook_receiver.go | 349ef32e52bcacdd216029cff466b976781ecf9ae83277a7b7fb016245a08724 |
| internal/whatsappbridge/webhook_receiver_test.go | c15b1a8352d9b75d8e35670d502b3e625ce945b4864759389e281bbfbb0e6c40 |
| whatsapp_cloud/whatsapp_cloud.py | 6310abef1497b31139aed605a8a3d6239ee8cfe23a12fe643cf22a681d673a28 |
| whatsapp_cloud/README.md | 1144ccc50194454f8a5b56a8eaeb61c3813de85d31925721d9deed8e6b1574c2 |
| whatsapp_cloud/PROVENANCE.md | b246231c69663a3b624e9aa3edd42330892471b8b1dd249dd929cad6ebe4c411 |

Temporal del sistema: elite-v274-a510f066303c44738bce9dcc13608bd0. Logs migrations.log, go-focused.log, go-whatsapp.log, python.log, roundtrip-go.log, roundtrip-python.log, roundtrip-vet.log y roundtrip-build.log. Repetir las instrucciones del README con base loopback descartable autorizada por approvalPool y Python explícito; ejecutar go test ./internal/whatsappbridge -v -count=1, Python unittest discover, vet y build.

## Autoridades y verdad de procedencia

El receiver, sus tests y el entry point ingress_bridge son AUTHORED, no nuevos archivos copiados de Meta. La función de validación que reutilizan es ADAPTED sobre los ejemplos exactos ya fijados de Meta; sus bytes VERBATIM, licencia y lock no cambiaron. No se afirma que esta integración sea código oficial, superior por reputación o productiva por pasar pruebas locales.

- [Meta/Postman — Message Status Update Notifications](https://www.postman.com/meta/whatsapp-business-platform/request/rgtfq23/message-status-update-notifications), consultado 2026-09-06: payload de statuses y uso de timestamps, no orden de llegada.
- [Meta — documentación del SDK Node archivado](https://whatsapp.github.io/WhatsApp-Nodejs-SDK/api-reference/webhooks/start/): contexto histórico del challenge GET y firma POST/ACK, explícitamente archivado; NO se adopta ese SDK ni se afirma una política vigente de reintentos a partir de él.

El overview actual de developers.facebook.com/documentation/business-messaging/whatsapp/webhooks/overview devolvió 429; su variante .md fue inaccesible. No se sustituyó autoridad por agregadores ni se afirmó verificación de ese contenido. La recurrencia y la lectura local truncada están registradas en FAIL-20260906-373. La suscripción, los límites y el contrato live actual deben demostrarse con el acceso del proyecto antes de exponer el endpoint.

## Qué falta para cerrar este journey

Worker que tome el job existente, revalide autoridad, resuelva el mensaje/receipt/tenant/organización correctos, procese todos los eventos admitidos y complete/reconcilie el job sólo tras efectos comprobados; retries/DLQ/retención operativos; UI y pruebas de navegador; cuenta, recepción y entrega reales. El fixture demuestra composición con routing conocido, no esos pasos automáticos. Ninguna compra, cita, consentimiento o envío se confirma a partir de un 200 de recepción.

## Verificación general y limpieza

VERIFY_LIBRARY_PASS, proceso exit 0: 160 packs / 1.415 archivos materializables / 707 Markdown / 51 perfiles. Franquicia 67/724. Procedencia actual AUTHORED=1177, ADAPTED=133, VERBATIM=105. No es un nuevo Audit ejecutable global ni certificación de producción.

Tras comprobar current_database=elite_whatsapp_v274 y cero sesiones activas, se eliminó exclusivamente esa base descartable (186 tenants sintéticos, tres webhooks y 17 observaciones de fixtures) y se detuvo el PostgreSQL iniciado para la prueba. Los fixtures y migraciones la recrean; staging/logs permanecen. No se activaron workers, cuentas, cron, suscripciones ni mensajes reales.
