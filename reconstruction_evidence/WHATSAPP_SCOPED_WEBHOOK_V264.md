# V264 — WhatsApp: evidencia de webhook vinculada a cuenta y teléfono

Fecha: 2026-09-06. Tarea: mantenimiento de biblioteca existente; delta del
adaptador seleccionado, sin crear producto, canal paralelo ni servicio pago.
Estado: `REBUILD_VERIFIED / CONDITIONED`, claim exclusivamente local.

## Resultado demostrado

`PYTHON-META-WHATSAPP-CLOUD-ADAPTER` pasa de 0.1.0 a 0.2.0. Cambian cinco
archivos existentes; el pack sigue materializando once. La franquicia conserva
67 packs/704 archivos, sin colisiones ni nuevas dependencias/migraciones.

Antes, un webhook firmado con WABA/teléfono ajenos era aceptado. La regresión
`test_signed_foreign_scope_is_rejected` produjo `PermissionError not raised`:
1 FAIL y 8 PASS. No se presentó aquel resultado como aislamiento demostrado.

Ahora el perfil exige `business_account_id` y `phone_number_id` exactos; el
normalizador exige `profile=` desde configuración autorizada y compara ambos
antes de publicar. Rechaza lotes mixtos sin producir evidencia parcial.
Valida estructuras, identidad, destinatario, estado y timestamp. Rechaza keys
JSON duplicadas y limita cuerpo/arrays/eventos. Son presupuestos locales, no
límites atribuidos a Meta.

Eventos y receipt `/v2` conservan scope pseudonimizado, recipient/message IDs,
timestamp, estado, hashes de bytes/perfil y key de evento determinista. No
deduplican durablemente ni ordenan artificialmente: preservan read/sent/failed/
delivered/deleted sin inferir entrega desde un HTTP exitoso. No asignan tenant
ni escriben dominio. El inbox y fence existentes siguen siendo sus owners.

## Autoridad y procedencia

- [Meta: webhook payload reference](https://www.postman.com/meta/whatsapp-business-platform/folder/vzaxn16/webhook-payload-reference), consultada 2026-09-06: WABA en entry.id, teléfono en metadata y envoltorio messages/whatsapp.
- [Meta: message status update notifications](https://www.postman.com/meta/whatsapp-business-platform/request/rgtfq23/message-status-update-notifications), consultada 2026-09-06: status, id, recipient_id y timestamp; llegada potencialmente desordenada.

No se adquirió una nueva revisión. Se conserva
`fbsamples/whatsapp-api-examples@de70ee908a67026e642aaee3703d20464e2a9466`, su
source lock, tres archivos `VERBATIM_REFERENCE_ONLY`, LICENSE adaptada sólo en
newline y licencia `LicenseRef-Meta-Platform-API-Only`. Los cambios locales al
runtime son `ADAPTED`; tests/config/docs son `AUTHORED`, no código escrito por
Meta. Consultar su contrato no constituye certificación de producción.

Se corrige además el texto de privacidad: no se guarda raw webhook, pero
`provider-response.json` outbound sigue siendo dato crudo sujeto a permisos y
retención. SHA-256 es pseudonimización, no anonimización.

## Gates ejecutados

Entorno: Windows, CPython 3.14.4 stdlib. Staging desechable:
`elite-v264-feee98f9656346ea8d963f9ca456880b` bajo temporal del sistema.

1. Materialización original: 11 archivos.
2. Regresión roja: `python -m unittest -v test_whatsapp_cloud.py`, 1 FAIL/8 PASS.
3. Corrección y misma suite ampliada: 17/17 PASS.
4. Fuente canónica sincronizada con `update_pack_from_tree.ps1`.
5. Reconstrucción en destino vacío: 11 archivos; 17/17 PASS.
6. Composición del perfil integral: 67 packs/704 archivos.
7. Comparación SHA-256 de los cinco archivos: staging = roundtrip = perfil integral.
8. `VERIFY_LIBRARY.ps1`: exit 0, `VERIFY_LIBRARY_PASS`, 160 packs/1.395 archivos/
   697 Markdown/51 perfiles.
9. `VERIFY_EXECUTABLE_LIBRARY.ps1 -Mode Audit` con Python 3.14.4 y Go 1.26.7
   exactos, sin `AllowNetwork`: exit 0, `VERIFY_EXECUTABLE_LIBRARY_PASS`,
   160 packs/121 fuentes/16 adapters. Repite los 17 tests WhatsApp desde su
   propia materialización. DevSkim runtime, proveedores cloud/documentos y
   otros probes con cuenta/red mantienen skips explícitos; OpenGrep runtime
   sigue bloqueado por Cosign SCA. El conteo de gates no equivale a runtime PASS.

Los 17 métodos incluyen matrices negativas para scope ajeno/ausente/mixto,
identidades vacías/excesivas, timestamp no canónico, estados desconocidos,
containers incorrectos, perfil sin demostrar, duplicate keys y presupuestos.
Prueban conservación de estados fuera de orden, keys estables ante cambio de
formato/orden, integridad del receipt, pseudonimización y no sobreescritura.
Incluyen los ocho tests anteriores de fuente/template/send/HMAC/challenge.
Todos los datos, firmas y transports son sintéticos. Cero mensajes externos.

| Archivo | SHA-256 materializado |
|---|---|
| whatsapp_cloud/whatsapp_cloud.py | d18fcd16e459a188edd53bc89e6491169b34488a7ff994c50d1c68e8ac499613 |
| whatsapp_cloud/test_whatsapp_cloud.py | 31ab00c822d5e37a722f9f9369e23330b0587556cb5ebb30741fc816589c3d53 |
| whatsapp_cloud/provider-profile.template.json | 3029f068e1d67dcdc774b43c28f89705a7440247b7df02b8869c3813601bf841 |
| whatsapp_cloud/PROVENANCE.md | df877adf12dd73bba2a7b063d7b4edf9f9ffdbd3ace4d770ee71b7912e9be679 |
| whatsapp_cloud/README.md | e91492f09ad90d9b31a6155a9a5afa9aafc616ab7bc75f7a65a1654a6eacdd9a |

## Compatibilidad y siguiente tramo pendiente

Los dos planes consumidores fijan 0.2.0. La llamada anterior sin perfil falla
cerrado; eventos v1 no pueden reinterpretarse como scope demostrado. El README
incluye migración, controles de datos, errores, recuperación y rollback sin
reactivar el defecto. El fence Go ya no declara compatibilidad ejecutable con
Python por mera coexistencia: el adapter aún no implementa su interfaz Sender.

Falta conectar appointment.confirmed al template/destinatario autorizado,
envolver el send con el fence existente y consumir estos statuses mediante el
inbox durable con tenant/contact binding y reconciliación. También faltan
evidencias reales de cuenta, Graph version, consentimiento/opt-out, costo,
entrega, webhook ACK, permisos y operación. Este trabajo NO cierra esos gates,
no revalida SCA global ni acredita load/seguridad ofensiva/producción.

Lección cerrada: `FAIL-20260906-362 / LIB-FAIL-2096`. No se borran antecedentes.
