# V306 — emisión de cotización recuperable y autoría persistente

> **Errata de trazabilidad V317 (2026-09-08), prevalece sobre las menciones históricas de Go1.26.7 debajo.** Se retira la atribución histórica exacta Go1.26.7 no vinculada a un receipt de identidad por ejecución. No se sustituye por1.26.8 sin prueba del artefacto concreto. La ruta current-toolchain de V281 contiene1.26.8; official-toolchain contiene1.26.7. El cuerpo original queda conservado como historia. V317 revalida únicamente la composición actual y el scope allí enumerado con1.26.7 observado, no todos los snapshots previos. Ver [TOOLCHAIN_IDENTITY_AND_RUNTIME_SCA_V317.md](TOOLCHAIN_IDENTITY_AND_RUNTIME_SCA_V317.md).

2026-09-07. Checkpoint37/V305, mantenimiento correctivo T2802/T2803/T2804
del roadmap vigente. Continúa V294: primer defecto ejecutable, gate conectado
y siguiente bloqueo real del mismo recorrido. V295, decisiones comerciales y
ARCA diferida intactos. No nuevo proyecto, cuentas, gasto ni efectos externos.

## Corrección y alcance

Journey0.10.11, Portales0.14.4, OIDC0.2.4 y BrowserGate0.1.16. Once archivos
existentes, sin dependencias, migraciones ni módulos nuevos. Perfil67/745.
Reconstrucción a destino ausente:11/11 hashes iguales a candidate. El cambio
OIDC corresponde al cliente server-only, no modifica el protocolo de login.

FAIL466: create-quote del formulario/BFF omitía Idempotency-Key, requerida por
Go. El test BFF mockeado esperaba201 y no comprobaba esa frontera. La regresión
de transporte falla sobre los bytes V305: falta el header. El test HTTP real
confirma400 sin clave; no se presenta ese rojo como doble emisión demostrada.

QuoteCreation guarda antes del POST sólo requestKey aleatoria bajo namespace
hash de tenant/subject/organización y lead. Controles cerrados hasta hidratación,
fence síncrono y almacenamiento fallido con cero POST. No conserva contacto,
variante, precio, notas ni secretos. UI→BFF→protectedPost transmite la misma
Idempotency-Key; BFF/cliente rechazan claves ausentes, cortas o malformadas.
No cambia pricing: moneda/total se resuelven en PostgreSQL desde la lista activa.
El recibo se contrasta con identidad, organización, lead, variante/lista, fecha,
estado/versión y tipos. Timeout10s conservador, sin claim de latencia medida.

GET /v1/franchise/quotes/result y GET del BFF leen sólo la relación existente
idempotency_record→quotation→lead. Exigen tenant, organización, lead, clave,
status completed, tipo quotation, response_code201 y referencia coincidente.
La relación de cliente con el lead también se contrasta. quote:write y scope
organizacional se verifican en HTTP; el BFF no sustituye esa autorización.
GET no emite cotizaciones ni reintenta POST. Resultado ausente/incoherente o
consulta fallida conserva referencia y bloqueo; un409 no significa no aplicado.
No depende de buscar cotizaciones parecidas por fecha o de una lista paginada.

La consulta muestra ID/estado/moneda del registro actual, no inventa una escala
monetaria. Tras recuperar, el formulario permanece cerrado para esa referencia;
no existe aquí un flujo de reemisión comercial deliberada. sessionStorage no
garantiza continuidad entre dispositivos, borrado o cierre del navegador.
Ayuda quote-create-view/1.0.0 enlaza práctica sintética, consulta y escalamiento;
no equivale a capacitación evaluada ni soporte integral.

Después de comprobar transporte/recuperación se continuó la auditoría del mismo
owner: FAIL469 encontró6 eventos sintéticos y0 con actor_subject=writer.
CreateQuoteAs toma Subject autenticado, falla cerrado sin repositorio auditado
y persiste actor con cotización/idempotencia/outbox en la misma transacción.
Replay no reemplaza al autor original, incluso con otro operador. Los métodos
internos legacy conservan compatibilidad y no heredan esta garantía de autoría.

## Procedencia y autoridad

Delta AUTHORED / LicenseRef-Workspace-Owner. No código externo copiado ni
actualización de dependencias. Se conservan revisiones/locks/licencias existentes.
FRONTEND_PRODUCT_ENGINEERING_UX y SOFTWARE_BACKEND_API_ENGINEERING gobiernan
visibilidad, recuperación y semántica de reintentos; el owner PostgreSQL conserva
invariantes, precio, atomicidad y permisos. Consultadas2026-09-07:

- [AWS — retries e identidad de intención](https://aws.amazon.com/builders-library/making-retries-safe-with-idempotent-APIs/).
- [React — refs](https://react.dev/learn/referencing-values-with-refs).

Estas fuentes explican método; no atribuyen el código local ni una aprobación a
AWS/Meta. No nueva SCA, admisión de runtime ni actualización global del corpus.

## Evidencia ejecutada

Staging: `%LOCALAPPDATA%/Temp/elite-v306-fcec733ad8f147ba965bd6857df879f9`.
candidate contiene web compilada; rebuilt-audited se compone desde los Markdown
canónicos. Once archivos idénticos. gate.ps1 preserva invocaciones exactas.

- Validator1.3.0 reconstruido15 archivos: resume37 y cadena37 PASS antes de editar.
- Preflight ejecutado: BLOCKED por cuatro herramientas no resueltas por default.
  Go/psql existentes se usan por ruta explícita en este gate; no reinstalación ni
  cambio de PATH global. Docker/target y readiness integral permanecen pendientes.
- Web108 PASS/1 SKIP y build PASS. El skip conectado se acredita separadamente
  mediante agenda y PostgreSQL; no se cuenta como PASS del comando unitario.
- Go test/vet/build PASS desde rebuilt-audited; suites opt-in se acreditan por
  sus invocaciones separadas, no por omisiones del comando general.
- PostgreSQL: suite completa del owner y tres regresiones anteriores,3 veces.
  Nuevo lookup scoped/clave desconocida/tres registros idempotentes incoherentes
  rechazados; replay, carreras, precio server-side y atomicidad existentes pasan.
  Actor original conservado tras replay por otro operador; fallo del evento
  deja cero cotizaciones/keys de la solicitud fallida.
- browser-final.log:68 fases/4 proyectos antes de corregir autoría; acredita
  recuperación, no auditoría de actor. browser-audited-final.log:68 fases/4
  proyectos sobre backend final, con3 cotizaciones/keys/eventos/actores por
  proyecto. Chromium desktop/mobile, Firefox desktop y WebKit desktop.
- Nueva fase: respuesta perdida tras commit, JSON201 vacío, storage fallido
  con cero POST, recarga con marker, GET503 redactado y recuperación por GET,
  consulta ajena/lead distinto/clave desconocida, HTTP403 sin permiso, replay200
  y divergencia409. Dos requests de prueba concurrentes conservan una identidad;
  BFF conserva su201 histórico, HTTP diferencia creación201/replay200.
  Las peticiones extra son inyección del test, no retries del producto.
- SQL exige exactamente3 cotizaciones,3 keys completed y3 eventos quotation.issued
  con writer por cada proyecto. Los recorridos previos de pedido/stock/pago local,
  entrega de fixture, leads y cancelación de disponibilidad siguen pasando.
- Agenda cuatro proyectos PASS con el transporte compartido final.
- Reinicio controlado del cluster sintético: snapshot de cotizaciones, claves
  y eventos idéntico por hash antes/después. Cluster detenido y datos conservados.
  No PITR, DR, despliegue ni recuperación del producto en un target real.
- Go1.26.7/PostgreSQL18.6/Node24.14.1/pnpm11.19.0/Next16.3.2/Playwright1.62.1.
  Instalación offline/frozen-lockfile, gate anidado --ignore-workspace.

Reproducir: componer FRANCHISE_COMPLETE_PACK_PLAN a destino ausente; instalar
web y gate con locks; build web. Base loopback elite_confirmation_* descartable
con53 migraciones. TEST_DATABASE_URL y ELITE_WEB_ROOT absoluto; ELITE_QUOTE_E2E,
ELITE_QUOTE_CREATE_E2E y flagsV305=1, ELITE_QUOTE_PROJECT vacío. Ejecutar
TestQuoteAcceptanceBrowserPostgres -count=1 -v -timeout15m. Agenda usa
ELITE_CONFIRMATION_E2E=1 y TestAppointmentAgendaBrowserPostgres. Nunca usar datos
de negocio para estos fixtures. Restart sólo después de finalizar todos los gates.

## Fallos conservados y continuidad

FAIL467 conserva rutas inferidas ausentes, glob inválido, Go fuera de PATH y
salidas truncadas; se retomaron refs/inventario literal y herramientas existentes.
FAIL468 conserva Python sin UTF-8/anclas no únicas, declaración repo omitida,
mock sin tipos y colisión recovered Quote/CustomerJourney. Logs rojos permanecen;
corrección de fixtures, sin debilitar controles, seguida de reconstrucción y gates.

FAIL457 sigue OPEN para otros comandos. Se inspeccionó creación de intervalos:
CreateAvailability genera ID nuevo, su handler no recibe identidad de solicitud
y Availability lista por rango con límite1000. No hay lookup por operación.
Por eso copiar la recuperación de cancelación o comparar fechas no resuelve una
creación incierta. Próximo delta: identidad retenida y resultado persistente de
create-availability sobre su owner y el registro idempotente existente; preservar
lock organizacional/actor/restricciones de agenda, sin inventar horarios.
Recursos, checklist, recepción y disposición conservan sus propios pendientes.

No cierre de T2802/T2804, readiness o assurance integral. No pruebas de multi-tab,
browser storage borrado, todas las respuestas malformadas, todas las carreras,
carga, privacidad/SCA integral, proveedor/pago real, liberación inicial de entrega,
PITR/DR, aceptación del target ni ZIP/promoción. V295 y ARCA diferida conservados.

## Integridad

| Artefacto | SHA256 |
|---|---|
| quote-transport-red.log | 826eaf482ffc9e16e28eb0926c241f1ccc4c1986aa027478e0ed96125d0750e2 |
| quote-actor-red.log | cf6712e15d98c0a20a2deb6609b90bb4d6a1bda177b8ed22ec89026c49075ced |
| web-test-final.log | 649d200b583ba7d40657daf3db004be198a435caaba9e1a9f24e107da4f6900e |
| web-build-final.log | 31734389d12063e7d8cc1276014a1188edd40b15f14ae8ef69c56310fe454099 |
| go-test-audited-final.log | b0327f0489241802e135e081db91b142d4a86a9fe032f8c545cbd197a17c2321 |
| go-vet-audited-final.log | e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 |
| go-build-audited-final.log | e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 |
| postgres-audited-final.log | 1369c118b4a90485507baac9b1fd765774b9e13f7f14665722b396eaeb3c3b42 |
| browser-final.log | b923510adc1e1e06c73beb8252b17ca4224af6bf486c0ae48e79cb341bbf52a9 |
| browser-audited-final.log | 9eff288de40cd307feee65933c769d09bbefb9f8eece4740d02cd26309f8bf9b |
| agenda-audited-final.log | c4c51e496d891d4267f2f3599af45c0657c4eee4e43b4968c85c570689c70d6b |
| restart.log | 23b2859d8b2209783d1c8e2da4cdd6547b27344d40f1ffa8fcf2284c410686bb |
| internal/franchisejourney/service.go | dc89cb01e69663c5b8f61abf74966d057e794680d5313cf786142fa793ca2604 |
| internal/platform/httpapi/franchisejourney.go | 4484a116431265565fb90660ad05e1226616b6e8db1c3a95ab4c6322aa2dc3b9 |
| internal/platform/httpapi/franchisejourney_test.go | 5c5955b3a6b91773aaf2dd7ea9ff960b9a37362a40a46c856df1207bd1bf3ac8 |
| internal/platform/postgres/franchisejourney.go | 8628fd9de42a5371f0dcbfabcc77704fb27953dc45047554c4a706f2869da2d7 |
| internal/platform/postgres/franchisejourney_integration_test.go | 808aa50bb66265f9094fc67ebc9db81bbfada860b88d6457afa0c63d9abbcf00 |
| src/platform/backend/protected-client.ts | 607a0f1ced7953c45bfa5e24c27f662f4a39fd7680045dcda051af02e52d3e6a |
| src/platform/backend/protected-client.test.ts | 5da9a69bf52d697b1411cde59502b821e95d48a987c81df1ec46a933dc2fc9c5 |
| microsoft_playwright_browser_gate/tests/enterprise-web.spec.mjs | b8a1c45c0e5a93c1b709525db2bc55985e824fa6fdc901a000f0959323b4d883 |
| src/components/franchise-command-panel.tsx | 3840ccb1a8a457d725cdf2db853543c24c4eedb815168873d6741104f4dda278 |
| src/app/api/enterprise/franchise/commands/route.ts | f3f0cc1383af1807880aa0f57d051413f46ede09a3c63d2db8d3a23288999a95 |
| src/app/api/enterprise/franchise/commands/route.test.ts | a3328677b6d8284c433e991753db1b91b44f1ceb46dd9498d9489758d37a8614 |

## Verificación final de biblioteca y continuidad

VERIFY_LIBRARY_PASS

El verificador de biblioteca pasó con checkpoint39; el cierre documental se
enlaza después en checkpoint40 y se valida de nuevo en nivel resume. Contrato
de ingeniería plan PASS; readiness/assurance integral siguen BLOCKED.

FAIL470: next_action mayor a800 fue rechazado antes de append; se compactó a598,
conservando revisión38 y cadena previa. El detalle permanece en este expediente.
No se cambió el límite del validator ni se borraron eventos.

SHA256 library-final.log: `f841705a0495fa4c47358f1be6a6011bc1e0bdc19ad58a4c6efeca7ed024af7e`.

FAIL471: el primer verificador rechazó el inventario preflight738 tras añadir
V306; se corrigió a739 sin cambiar el gate, y se repitió el verificador completo.
SHA256 library-inventory-red.log: `a936721f415209f984a06574b4cd319d1d785d0148aab71ac7424ae00d23fa0a`.
