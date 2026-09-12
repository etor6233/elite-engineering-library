# V267 — operación HTTP autenticada de notificación de turno

Fecha: 2026-09-06. Mantenimiento de biblioteca EXISTING, delta sobre V266.
Estado: REBUILD_VERIFIED / CONDITIONED. No acredita una entrega Meta real.

## Cambio y conexión

PYTHON-META-WHATSAPP-CLOUD-ADAPTER 0.5.0 añade dos archivos AUTHORED y actualiza
README/PROVENANCE. El pack materializa 20 archivos; el perfil integral compone
67 packs/713 archivos, conserva 50 migraciones y no añade dependencia/runtime.

NewAppointmentNotificationModule recibe el resolver V266, Sender V265, store
PostgreSQL existente y receiver verificado. Copia sender/perfil y fuerza el mismo
resolver en la aprobación y el envío. Register instala una ruta en el mux existente:
POST /v1/franchise/appointments/{id}/whatsapp-confirmation.
No crea listener, worker, tarea programada, cola, consentimiento ni cuenta.

El cliente debe aportar los doce campos tipados de AppointmentNotificationRequest.
Tenant/actor se obtienen del verifier; propósito/política/perfil del servidor;
canal y key por evento se derivan sin intervención del LLM. Requiere permiso
appointment:manage, tenant configurado y organización autorizada. Sólo Bearer:
ninguna identidad desde cookie, body ni forwarded header.

El parser rechaza campos desconocidos, duplicados, aliases por case, nulls,
elementos null de parámetros, tipos incompatibles y JSON concatenado. Body
máximo 64 KiB y contexto 45 s son presupuestos locales, no límites de Meta.
La aprobación se comprueba antes de reclamar el fence; el Sender la vuelve
a verificar después. Un cambio concurrente aún puede invalidar esa segunda
verificación; el fence conserva el resultado incierto sin reenvío automático.

200 accepted significa aceptación del provider o replay de ese estado durable,
nunca entrega al destinatario. 409 distingue aprobación inválida/divergente,
operación en curso, terminal y reconciliación necesaria. Respuestas no-store
sin teléfono/template/token ni error interno. No se almacena una nueva copia
del mensaje: grant y fence conservan sus hashes y evidencias existentes.

## Procedencia y fuentes consultadas

Estos dos archivos son composición AUTHORED; no código publicado ni certificado
por Meta, AWS o Google. Siguen intactos los tres archivos Meta VERBATIM, licencia
API-only, snapshot y adaptación Python. El patrón Register y identity.Verifier
son los mismos contratos de los módulos existentes.

Fuentes primarias consultadas 2026-09-06:

- [AWS transactional outbox](https://docs.aws.amazon.com/prescriptive-guidance/latest/cloud-design-patterns/transactional-outbox.html): actuar sobre estado comprometido y consumidores idempotentes; no certifica este endpoint.
- [WhatsApp Business Messaging Policy](https://whatsappbusiness.com/policy/): consentimiento, opt-out y template admitido; un hash o token local no demuestra cumplimiento real.
- [Go encoding/json](https://pkg.go.dev/encoding/json): decoder, campos case-insensitive y tratamiento de null; el contrato estricto adicional es local.

## Ejecución observada

Windows, Go 1.26.7, CPython 3.14.4, PostgreSQL 18.6.
Temporal: elite-v267-ce668e893c794d97902de515c8dc5ab9.
Base dedicada elite_whatsapp_v267 sobre loopback:55959, 0001–0050 aplicadas
desde la composición. No migración nueva en V267.

- Tres tests superiores nuevos, 27 subcasos: 11 del parser, 13 negativos HTTP
  y tres recorridos de ejecución/recuperación. Pasan en staging y reconstrucción.
- HTTP real → discovery/JWKS local → verifier OIDC existente con firmas RS256
  reales → aprobación PostgreSQL → fence PostgreSQL → proceso Python sintético.
- Ocho POST simultáneos aceptan o informan conflicto/en curso, pero dejan un
  grant y una sola invocación. Repeticiones posteriores devuelven accepted.
- Pérdida de respuesta del proceso provider deja unknown; dos llamadas nuevas
  exigen reconciliación y no repiten el proceso.
- Pérdida de la respuesta HTTP después del commit se inyecta cerrando la conexión:
  el cliente observa error de red. Repetir exactamente el request obtiene accepted
  y sigue habiendo una sola invocación.
- Modificar parámetros bajo el mismo evento devuelve 409, sin otro envío.
- Token ausente/firma alterada, rol, organización, tenant, destinatario,
  consentimiento, revocación, ID de turno, tamaño, media type, identidad
  inyectada y Authorization duplicado dejan cero grants y cero deliveries.
- go test -count=1 ./..., go vet ./..., go build ./...: exit 0 en roundtrip.
  ELITE_WHATSAPP_PYTHON y ELITE_WHATSAPP_TEST_DATABASE_URL configurados;
  gates que dependen de otras variables no se presentan como ejecutados.
- Python: 22/22 tests PASS. Los cuatro archivos modificados/nuevos son idénticos
  por hash entre work y roundtrip. No se cambiaron tests anteriores ni triggers.

El issuer, tokens, cita y consentimiento son fixtures. Se reutiliza el fixture
de dominio V266, no el endpoint de confirmación V263. El child prueba el borde
de proceso, NO Meta. No hubo login real, navegador/BFF, tráfico Meta ni costos.

## Hashes del delta reconstruido

| Archivo | SHA-256 |
|---|---|
| internal/whatsappbridge/appointment_notification.go | 8e1410a0bba905671ec76257195e749bc53f6019937f4ef2a5075b8a70725348 |
| internal/whatsappbridge/appointment_notification_test.go | 32ae8990b26158e727bc990fca3109183c5d4d98c32dda33dcd99a5172283a51 |
| whatsapp_cloud/README.md | 00498f50b63bc68fcac4f04d4f9bbf8f0e5ea2ee38be8acd279b79c736178b05 |
| whatsapp_cloud/PROVENANCE.md | d9c33dfb763001764e0110d93150d3c52b03f7e3ddbdbfad445d7e060a24d337 |

## Uso, límites y rollback

El README materializado muestra construcción/registro, campos exactos, budgets
y semántica de errores. La integración es opt-in: el composition root existente
aún debe proporcionar los owners/configuración y registrar el módulo.
La API es explícita y síncrona, no un outbox consumer. Un crash entre grant y
send requiere reanudar con el mismo request mientras siga autorizado. Un grant
expirado/revocado devuelve 409; no cambiar evento/expiry para forzar un reenvío.
No existe todavía aquí una API de lectura histórica de delivery.

Pendiente en este recorrido: UI/BFF del operador, worker cuando el blueprint
requiera notificación eventual, lectura/correlación de statuses, inbox,
reconciliación/soporte, IdP y Meta reales, template/purpose/consentimiento,
rate/cost/egress, observabilidad, carga, seguridad y despliegue.
El flujo de confirmación de turnos anterior no llama automáticamente a este endpoint.

Rollback: retirar/pausar la ruta y consumidores, conservar grants/fence/evidencia.
No borrar registros ni cambiar keys para reenviar. Las dependencias y 0050 no
cambiaron; no corresponde ejecutar down como rollback de esta API.

FAIL-20260906-368 / LIB-FAIL-2102 conserva la recurrencia de inventario/ruta
de lectura incorrectos, corregida antes de implementar con owners observados.
No se detectó un fallo funcional nuevo en las suites de este delta.

## Gate raíz y limpieza

VERIFY_LIBRARY_PASS, exit 0: 160 packs, 1.404 archivos materializables,
700 Markdown y 51 perfiles; integral 67/713. No se reejecutó Audit integral
en V267 y no se atribuye su PASS histórico a este delta.

SQL final: 74 tenants Synthetic approval y cinco Synthetic bridge, 16 grants,
ocho deliveries accepted y siete unknown, acumulados por ejecuciones/fixtures.
Se comprobó cero tenants ajenos a esas identidades sintéticas y se eliminó
sólo elite_whatsapp_v267. El servidor iniciado por la tarea quedó detenido;
no se tocaron otras bases. Las filas temporales se descartaron y se reconstruyen
con los tests; los fuentes y esta evidencia se conservan.
