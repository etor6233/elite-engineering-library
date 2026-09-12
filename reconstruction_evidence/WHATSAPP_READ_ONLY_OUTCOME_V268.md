# V268 — consulta autenticada del resultado local de notificación

Fecha: 2026-09-06. Mantenimiento de biblioteca EXISTING, delta sobre V267.
Estado: REBUILD_VERIFIED / CONDITIONED; no observación de entrega Meta real.

## Cambio material y ownership

PYTHON-META-WHATSAPP-CLOUD-ADAPTER 0.6.0 incorpora dos archivos AUTHORED
(lector/tests), comparte la autenticación del módulo HTTP existente y actualiza
README/PROVENANCE. Son cinco archivos del delta, 22 en el pack y 715 desde
67 packs en el perfil integral. No añade migración, runtime, dependencia o tabla.

GET /v1/franchise/appointments/{id}/whatsapp-confirmation exige los parámetros
organization_id y confirmation_event_id, un Bearer verificado, tenant configurado,
appointment:manage y organización autorizada. Rechaza query duplicada/desconocida,
evento inválido y body. El endpoint usa el mismo mux/module que POST; no crea
servidor ni activa tráfico. Presupuesto local de consulta: cinco segundos.

ReadNotificationStatus hace un SELECT de aprobación, turno/lead actuales y
fence existente. Un grant expirado, contacto revocado o turno cancelado no
oculta el resultado histórico a un operador que mantiene alcance actual.
Reasignación de lead/organización invalida esa lectura. No es permiso de envío.

Comprueba igualdad de request hash y recipient HMAC entre grant y fence.
Si accepted carece de los campos de evidencia requeridos no se publica como
válido. Hashes, teléfono, texto, actor, consentimiento, provider ID y códigos
internos no forman parte de la respuesta. No-store y errores 404/409/503
distinguen ausencia, discrepancia y lectura indisponible sin detalles privados.

El snapshot separa fence_state y delivery_status. El primero refleja
not_started/sending/accepted/unknown/failed_terminal. El segundo permanece
not_observed_by_this_reader: no se correlacionaron webhooks Meta aquí.
Un sending con lease vencido señala reconciliation_required sin cambiar el
estado, crear eventos ni llamar Claim. No se usa POST para averiguar historia.

## Autoridades y procedencia

Todo el delta es AUTHORED, no código publicado o certificado por Meta.
Los tres archivos Meta VERBATIM, fuente/licencia, Python adaptado, SDKs y pins
permanecen intactos. Consultas oficiales del 2026-09-06:

- [Meta Message Status Update Notifications](https://www.postman.com/meta/whatsapp-business-platform/request/rgtfq23/message-status-update-notifications): sent, delivered y read son observaciones separadas con timestamp; su llegada puede no estar ordenada.
- [PostgreSQL 18 Transaction Isolation](https://www.postgresql.org/docs/18/transaction-iso.html): SELECT observa estado comprometido en un snapshot. No se atribuye atomicidad global entre la lectura y un envío posterior.

## Evidencia ejecutada

Windows, Go 1.26.7, CPython 3.14.4, PostgreSQL 18.6.
Staging elite-v268-711a4fa6f7e142ce9b5d4c578608cfc3 bajo temporal del sistema;
base dedicada elite_whatsapp_v268, loopback:55959; 50 migraciones aplicadas.

1. Tres tests nuevos: 11 casos de estado/integridad, 13 de scope/errores y un
   recorrido POST→GET→GET por HTTP real. Verifier OIDC real con issuer y claves
   RS256/JWKS sintéticos; no login/IdP productivo.
2. Los 11 casos de lectura usan un pool con default_transaction_read_only=on;
   SHOW transaction_read_only confirma on. Dos GET por caso conservan cantidad
   de eventos, no crean proceso y no mutan sending aunque el lease esté vencido.
3. Accepted sigue visible cuando vence la aprobación, se cancela el turno o se
   revoca el contacto. La lectura no vuelve a aprobar ni libera el fence.
4. Request/recipient hash discrepantes devuelven 409. Scope ausente/ajeno,
   reasignación y evento inexistente no filtran registros. Pool cerrado devuelve
   503, nunca falso 404/accepted ni información interna.
5. POST usa el proceso Python sintético existente; los dos GET posteriores
   recuperan accepted local y no repiten la única invocación.
6. Primer run: sólo falla la preparación del fixture de reasignación por
   lead_id NOT NULL. Corregido con segundo lead válido, conservando NOT NULL/FK.
   Focal staging y full Go en reconstrucción pasan. Registros FAIL-20260906-369
   / LIB-FAIL-2103; el patch de README rechazado queda como FAIL-370/LIB-2104.
7. go test -count=1 ./..., go vet ./..., go build ./...: exit 0 desde roundtrip.
   ELITE_WHATSAPP_PYTHON y ELITE_WHATSAPP_TEST_DATABASE_URL configuradas;
   otros gates live que requieren variables separadas no se presentan ejecutados.
8. Python 22/22 PASS. Cinco hashes idénticos entre staging y reconstrucción
   limpia de 715 archivos. Fuentes/test outputs temporales en go-test.log.

Los fixtures SQL no demuestran captura real del consentimiento ni ejecución de
la confirmación V263. El proveedor hijo es sintético. No hubo llamada Meta,
mensajes externos, coste nuevo ni actualización de dependencias.

## Hashes del delta

| Archivo | SHA-256 |
|---|---|
| internal/whatsappbridge/appointment_notification.go | 88257dccb1692f60f90a909f71ecc2800b84d289b61f2157c54b74dc852ea7ef |
| internal/whatsappbridge/appointment_notification_status.go | c3f410f947b3c727d540d4fab0e819e3b1eadcfd14fc5cfec4c7c003eaacf166 |
| internal/whatsappbridge/appointment_notification_status_test.go | 8158977f24dddc9cf2cd9725ae4f918347c3d7485cfd46b3edf62a45fb793d97 |
| whatsapp_cloud/README.md | 36d166a8460751e590577132ecad5a479ee1d9dacd6c1278a02467b244b92606 |
| whatsapp_cloud/PROVENANCE.md | 977f0401867085503608100c640983967abd1fee9fe9c0faa346e962cb8ad26c |

## Estado real y siguiente integración

El resultado local de envío ya es consultable sin reintentar POST y sin duplicar
ledger. Sigue faltando la experiencia UI/BFF del operador, automatización cuando
aplique, consumer/correlación de statuses, soporte/reconciliación operativos y
evidencia Meta/IdP/costos/consentimiento reales. Este GET no declara entrega ni
reemplaza esos gates.

El mapa general conserva además journeys completos de cotización/pedido/recovery,
ayuda/capacitación/soporte release-bound, telemetría sin PII, bootstrap compacto,
mutaciones de proveedores seleccionadas y promoción por claim. No se traduce
el conteo de archivos en un porcentaje de cierre ni en plazo de franquicia.

Rollback: pausar/retirar ruta conservando aprobación/fence/evidencia; no correr
down, borrar registros ni alterar keys. El README contiene contrato y límites.

## Verificación raíz y limpieza

VERIFY_LIBRARY_PASS exit 0: 160 packs, 1.406 archivos materializables,
701 Markdown y 51 perfiles; franquicia 67/715. No se reejecutó Audit integral
en V268: el PASS de aquí cubre gates focales, full Go y estructura.

Antes de limpiar, SQL observó 117 tenants Synthetic approval y cinco Synthetic
bridge; estados acumulados de fixtures: 25 accepted, ocho unknown, seis sending
y tres failed_terminal. Se comprobó cero tenants ajenos a esos fixtures,
se eliminó sólo elite_whatsapp_v268 y se detuvo el servidor iniciado por esta
tarea. No se tocaron otras bases. Esas filas temporales son reconstruibles desde
los tests; sus fuentes y evidencia se conservan.
