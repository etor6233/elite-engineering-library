# V266 — aprobación durable de confirmación para WhatsApp

Fecha: 2026-09-06. Mantenimiento de biblioteca EXISTING, delta sobre V265.
Estado: `REBUILD_VERIFIED / CONDITIONED`, no notificación live certificada.

## Implementación y procedencia

`PYTHON-META-WHATSAPP-CLOUD-ADAPTER 0.4.0` incorpora cinco archivos
AUTHORED: resolver Go, tests, migración 0050 up/down y test SQL. README y
PROVENANCE actualizados. El pack materializa 18 archivos y el perfil integral
67 packs/711 archivos con 50 migraciones; no añade dependencias ni runtime.

No son cinco archivos publicados por Meta/AWS. Las tres referencias Meta
VERBATIM, licencia Platform-API-only, snapshot y adaptación Python se conservan.
El nuevo código conecta los owners existentes; no añade otra agenda, CRM,
consentimiento, binding, cola o delivery ledger.

Autoridades consultadas 2026-09-06:

- [AWS transactional outbox](https://docs.aws.amazon.com/prescriptive-guidance/latest/cloud-design-patterns/transactional-outbox.html): sólo actuar sobre eventos comprometidos y mantener consumidores idempotentes.
- [WhatsApp Business Messaging Policy](https://whatsappbusiness.com/policy/): opt-in, respeto de opt-out y templates aprobados; el registro SQL no acredita por sí solo permiso real/legal ni aprobación Meta.

## Contrato demostrado

El caller aporta Principal verificado con appointment:manage y organización
permitida. Propósito y política se fijan en el servidor. El único INSERT…SELECT
exige turno confirmado vigente, transición inmutable y outbox con el mismo
event ID/version, lead/organización, binding HMAC activo/versionado/PII y
consentimiento exacto otorgado. Otra decisión posterior o empatada en tiempo
para ese propósito rechaza el grant; no se inventa desempate.

El grant inmutable vincula evento, versión, contacto, consentimiento, actor,
hashes del mensaje/perfil/evidencia y vencimiento. Una key determinista por
tenant/evento impide cambiarla para repetir la misma confirmación. No guarda
teléfono ni texto del template; sus hashes siguen siendo datos sensibles.

Resolve revalida turno, contacto, política, propósito, consentimiento y
vencimiento antes del send. Conserva prueba de confirmación aun después de la
limpieza normal del outbox. No renueva permisos ni recupera automáticamente
mensajes inciertos. Es snapshot previo, NO revocación atómica de un HTTP en vuelo.
El máximo de 24 horas del grant es presupuesto LOCAL, no norma de Meta.

## Pruebas ejecutadas

Windows, Go 1.26.7, CPython 3.14.4, PostgreSQL 18.6.
Staging `elite-v266-f1e4d3519b2f4e01b88f680368f84839`, bajo temporal del sistema.
Reconstrucción a destino nuevo `roundtrip`, 711 archivos/67 packs.

1. Tres tests superiores nuevos: 19 rechazos de prueba incompleta/divergente,
   seis cambios/revocaciones y caso positivo con ocho aprobaciones concurrentes
   idénticas que dejan un grant. Replay alterado y modificación/borrado del
   grant se rechazan; no se deshabilitan triggers.
2. Resolver PostgreSQL real → Sender → proceso Python sintético fijado →
   fence PostgreSQL real: dos intentos, una invocación y receipt accepted.
   Los seis casos de invalidación no invocan al hijo.
3. `go test -count=1 ./...`, `go vet ./...`, `go build ./...`: exit 0
   desde la reconstrucción; Go test registrado en `go-test.log` temporal.
   ELITE_WHATSAPP_PYTHON y ELITE_WHATSAPP_TEST_DATABASE_URL configurados:
   ejecutados los tests nuevos y los cinco casos PG anteriores V265.
   Esto no activa los demás gates de integración que requieren otras variables.
4. Python adapter: 22/22 PASS desde reconstrucción, transporte inyectado.
5. Base sintética `elite_whatsapp_v266` con 50 migraciones. Otra base vacía
   `elite_whatsapp_v266_migration` aplica 0001–0050 y prueba 0050 down,
   ausencia de tabla, up y test SQL. El test SQL prueba metadata; las
   invariantes funcionales anteriores se demuestran en Go/PostgreSQL.
6. Los siete archivos del delta son byte-idénticos por SHA-256 entre staging
   y reconstrucción. No se adquieren versiones nuevas ni se envía a Meta.
7. Verificador raíz: VERIFY_LIBRARY_PASS, exit 0; 160 packs, 1.402 archivos,
   699 Markdown y 51 perfiles. No se reejecutó Audit integral en V266: este
   cierre acredita los gates focales, full Go y verificación estructural.

La inspección final observa 60 tenants Synthetic approval y cinco Synthetic
bridge; 15 grants y siete deliveries (tres accepted, cuatro unknown) acumulados
entre fixtures y ejecuciones. No son mensajes reales ni un conteo de clientes.
La base de migración contiene cero tenants. Tras comprobar que no había tenants
ajenos a esos fixtures, se eliminaron sólo las dos bases V266 y se detuvo el
servidor que inició esta tarea. Los fixtures son reconstruibles; sus filas
temporales fueron descartadas. No se eliminaron otras bases ni sus archivos.

Las citas/outbox se siembran como fixtures SQL y el Principal es fixture:
NO se probó aquí el endpoint real de confirmación ni el recorrido navegador.
V263 conserva su evidencia distinta; concatenar ambos PASS no es E2E live.

## Hashes de reconstrucción

| Archivo | SHA-256 |
|---|---|
| internal/whatsappbridge/appointment_approval.go | 50f8f778be5cd5ced634037fe29c642bb3faf3f106d176e64239305e30644adc |
| internal/whatsappbridge/appointment_approval_test.go | 057c9fc0031b6566fe27bfe496f229282d69b88e44b2614c73990d2e720c34e2 |
| db/migrations/0050_whatsapp_appointment_approval.up.sql | e880a343999d9d79de15885a56ea6aab21776ee94ac53f489ce711f74418e8a8 |
| db/migrations/0050_whatsapp_appointment_approval.down.sql | 1aff4e7b751d1db0349074c0a9e36bcd55bd0e45b0fe5aa9da8d2fe9c91e2084 |
| db/tests/0050_whatsapp_appointment_approval.test.sql | db9e4dce1fb6c19b8d5590ed72f77cce5a2a012fdca53ebb62a8ce13f3a7fc41 |
| whatsapp_cloud/README.md | 7c11d4dd0efc35948f3c10fd17127186a25da36a78693191efe46ae394ff7b10 |
| whatsapp_cloud/PROVENANCE.md | f7ff5cd38abc71f189645b36e84b8c8c05a7e5611937a203414d8d86b48bb7f2 |

## Fallos y cierre pendiente

FAIL-20260906-366 / LIB-FAIL-2100 conserva recurrencia de ruta inferida.
FAIL-20260906-367 / LIB-FAIL-2101 registra el seed tenant_code inválido:
corregido con prefijo wa- sin relajar constraint; suite staging y reconstruida
pasan. No se borra el resultado inicial.

Pendiente: integración autenticada API/operador/worker con template aprobado,
secret store/receiver real, cuenta y consentimientos reales, inbox/status
correlacionados, recuperación/operación, carga, seguridad y gates del target.
No se afirma una franquicia terminada ni un porcentaje global de cierre.

Rollback operativo: pausar consumer y conservar grants/fence/evidencia.
No usar down para reintentar un envío; sólo rollback de schema con autorización,
export/retención y consumers detenidos. No cambiar DeliveryKey para superar
unknown. El README documenta construcción y condiciones de uso.
