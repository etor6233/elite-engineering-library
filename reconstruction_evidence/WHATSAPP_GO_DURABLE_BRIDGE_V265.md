# V265 — enlace Go/Python al fence durable existente

Fecha: 2026-09-06. Mantenimiento de biblioteca EXISTING, delta sobre V264.
Estado: `REBUILD_VERIFIED / CONDITIONED`; no notificación productiva certificada.

## Implementación y ownership

`PYTHON-META-WHATSAPP-CLOUD-ADAPTER 0.3.0` añade dos archivos Go y modifica
cuatro archivos Python/documentación. Materializa 13 archivos; el perfil integral
compone 67 packs/706 archivos, sin nueva migración, runtime o dependencia.

`whatsappbridge.Sender` implementa el contrato `outbounddelivery.Sender` del
owner existente. Usa un perfil/tenant configurado, ApprovalResolver obligatorio
con hash del mensaje y perfil, actor/evidencia/ventana temporal, y TokenSource
externo. La request de template se valida y el destinatario debe coincidir con
el envelope fenced. No genera aprobación ni consentimiento.

El proceso hijo es Python exacto con adapter exacto, paths absolutos, SHA-256,
`-I -B`, sin shell ni variables heredadas salvo SystemRoot. Token sólo por stdin,
frame/response acotados, stderr descartado, deadline finito. El entry point de
Python reutiliza `send_template_to_evidence`, verifica la referencia empaquetada
y devuelve identidad única y hash del receipt conservado; el receipt vincula los
bytes de respuesta. No agrega un segundo cliente HTTP, inbox o cola.

El caller SIEMPRE debe envolver el Sender en el Channel/Store PostgreSQL
existente. Resultados inciertos no autorizan otra invocación. Un ACK aceptado
no confirma entrega al cliente, asistencia al turno ni una venta.

## Procedencia y autoridades

Los archivos Go/tests/configuración son `AUTHORED`; la entrada de proceso amplía
la adaptación Python `ADAPTED`. No se presentan como código Go publicado por
Meta. Se conserva el source lock, commit Meta firmado
`de70ee908a67026e642aaee3703d20464e2a9466`, los tres archivos VERBATIM de
referencia y la licencia Platform-API-only con newline declarado.

- [Meta: status notifications](https://www.postman.com/meta/whatsapp-business-platform/request/rgtfq23/message-status-update-notifications), consultada 2026-09-06: estados posteriores separados y llegada potencialmente fuera de orden.
- [AWS: making retries safe](https://aws.amazon.com/builders-library/making-retries-safe-with-idempotent-APIs/), consultada 2026-09-06: request identity y efectos ambiguos; no demuestra idempotency key nativa de Meta.

## Ejecución observada

Windows; Go oficial 1.26.7, CPython 3.14.4 y PostgreSQL 18.6 locales.
Staging `elite-v265-d22ed88f5013470cb0b6b46ac1669996` bajo temporal del sistema.
Base exclusivamente sintética `elite_whatsapp_v265`, 127.0.0.1:55959, creada
para esta tarea con las 49 migraciones existentes verificadas antes de aplicar.

1. Python: 22/22 PASS, incluyendo sender reutilizado con transport inyectado,
   receipt ligado a evidencia, frame malformado, incertidumbre, identidad múltiple
   y CLI real con perfil bloqueado/sin fuga de secrets.
2. Go: cuatro tests superiores y 19 subcasos. Proceso real con script de proveedor
   expresamente sintético y hash-fijado; no se sustituyó por fixture el CLI
   productivo de Python ni se llamó a Meta.
3. Cinco casos PostgreSQL: aceptación, pérdida de respuesta, binding incorrecto,
   ID ausente y stdout excesivo. Cada caso prueba dos intentos con nuevas
   instancias Store, una sola invocación, conflicto al alterar el request,
   attempt_count=1 y dos eventos durables. Estados: accepted o unknown, nunca
   delivered inferido. No se deshabilita el trigger de inmutabilidad.
4. Después de dos ejecuciones (staging y reconstrucción), consulta SQL observó
   2 accepted, 8 unknown y 20 eventos; cada fila conserva un único intento.
5. `go test -count=1 ./...`, `go vet ./...`, `go build ./...`: exit 0 en el perfil
   reconstruido. ELITE_WHATSAPP_PYTHON y ELITE_WHATSAPP_TEST_DATABASE_URL estaban
   configurados; los otros gates live/DB del producto no se activaron por ello.
6. Materialización limpia: 706 archivos/67 packs; los seis archivos del delta
   son idénticos por SHA-256. Tests Python y Go/PG repetidos desde esa copia.
7. Verificador raíz: `VERIFY_LIBRARY_PASS`, exit 0, 160 packs/1.397 archivos/
   698 Markdown/51 perfiles. No se presenta el Audit V264 como reejecutado V265;
   aquí se ejecutaron los gates focales, full Go y verificación estructural.

No se enviaron mensajes externos, adquirieron dependencias ni habilitaron costos.
Los tests de procesos usan un provider sintético y los Python un transport
inyectado: no son evidencia de una cuenta o entrega Meta real.

## Hashes reconstruidos

| Archivo | SHA-256 |
|---|---|
| whatsapp_cloud/whatsapp_cloud.py | a5835e1fe346f115ac920119f5e1d6c0a31396340602d90c4cd8b97b8c67361e |
| whatsapp_cloud/test_whatsapp_cloud.py | c45a97b3584b6f308316d9760610fca4a5a4354c3d9c2146af5180402f42f421 |
| whatsapp_cloud/PROVENANCE.md | 43cd3db789424c82432615198e46619c7013f0f7b9ad68ce5a38277e791a183d |
| whatsapp_cloud/README.md | 18d230405819cefc246c2f04f18ff275fe9596a708fb03e1d5ff50133a830107 |
| internal/whatsappbridge/sender.go | 094db800ac72af136449260d607b216f07741c9ed19715f0b219eb7cde3f3735 |
| internal/whatsappbridge/sender_test.go | 4f7d4b04e4a0e8811b14f3234873f4b0308a7c164f5609e8a315389ffb3e045d |

## Fallos conservados

`FAIL-20260906-363 / LIB-FAIL-2097`: glob de rg no portable, corregido con filtro
nativo antes de implementar. `FAIL-20260906-364 / LIB-FAIL-2098`: cinco casos Go
y uno Python fallaron por LF→CRLF de streams de texto Windows. Se corrigieron
protocolo/marcador con bytes, sin relajar aserciones. La reconstrucción repite
los negativos y positivos; no se borra el primer resultado.

`FAIL-20260906-365 / LIB-FAIL-2099`: el primer gate raíz rechazó el resumen
de notices por omitir los dos archivos AUTHORED nuevos. Se sincronizó a
AUTHORED=1159, ADAPTED=133, VERBATIM=105, TOTAL=1397 y el gate pasó sin cambios
en su política. Luego se eliminó sólo elite_whatsapp_v265, comprobando antes
que todos sus tenants eran Synthetic bridge, y se detuvo el servidor local que
esta tarea había iniciado. El puerto 55959 quedó sin listener; no se eliminaron
datos de proyectos ni otras bases históricas.

## Lo que todavía no cierra

ApprovalResolver durable ligado a appointment.confirmed/contacto/policy,
selección del template y su contenido de negocio, receiver real, correlación
del status con el outbound e inbox y recuperación operativa. También cuenta,
permisos, consentimiento/opt-out, costo, TLS/egress, capacidad, carga, seguridad,
SCA actual del despliegue y aceptación. El hash del ejecutable Python no prueba
toda su distribución; el target requiere bundle admitido y árbol de solo lectura.

Rollback: conservar evidencia y unknown; pausar el consumer, no llamar al CLI
directamente ni cambiar DeliveryKey para forzar un reenvío. El README incluye
configuración, comandos focales y límites para la siguiente implementación.
