# V278 — host controlado, historial WhatsApp y ayuda del operador

Fecha: 2026-09-06/07, mantenimiento de biblioteca. Baseline V277, sin reiniciar
arquitectura ni certificar un proyecto productivo. Staging conservado:
elite-v278-b5b67ce6d1084015be134d0f47f15061; árbol ejecutado: verified;
cierre documental reconstruido: release-review. Comparación de 742 archivos:
738 idénticos, tres README/PROVENANCE reordenados para poner V278 primero y
next-env.d.ts con las dos importaciones .next/types que genera Next durante
build (routes y root-params), comprobadas exactamente. Ninguna otra diferencia.
Los 22 hashes siguientes corresponden al cierre documental.
La revisión final alinea además compatible_with del gate Playwright y del
bridge con las versiones realmente probadas; no modifica bloques ejecutables.

## Resultado y procedencia

Se amplían cuatro owners existentes, sin nuevo pack, dependencia, migración,
cuenta, suscripción, Git, Actions ni envío externo:

- PYTHON-META-WHATSAPP-CLOUD-ADAPTER 0.12.0: host Run opt-in y colección de
  historial autorizada/read-only, reutilizando cola, router, observer y fence.
- TS-FRANCHISE-JOURNEY-PORTALS 0.11.0: BFF y consulta en agenda, recuperación
  GET, fechas distintas y ayuda vinculada a whatsapp-status-view/1.0.0.
- MICROSOFT-PLAYWRIGHT-BROWSER-GATE 0.1.6: fixture conectado y regresiones
  de coexistencia entre mensajes de agenda y notificación.
- TS-GO-API-WEB-BRIDGE 0.5.3: dos reglas CSS de la vista, mismo owner.

Los doce archivos nuevos son AUTHORED. Las veintidós rutas nuevas/modificadas
no son nuevo código copiado de Meta, Google o Microsoft. Se conservan los pins,
licencias, notices y referencias/adaptaciones oficiales previas. El runtime
Microsoft Playwright sigue siendo la dependencia oficial fijada, no el autor
del harness, la pantalla, las reglas locales ni las pruebas de dominio.
Estado de integración: REBUILD_VERIFIED / CONDITIONED para este claim acotado.

## Conexión demostrada

Turno y rol → agenda existente → consulta explícita → BFF/session existente
→ firma RS256/JWKS → autorización tenant/organización → snapshot PostgreSQL
→ observación anclada del worker → respuesta comprensible → ayuda versionada.

El fixture precrea una aprobación y un outbound accepted sintéticos; el worker
real procesa el webhook retenido. No demuestra que un navegador envió un mensaje
a Meta. La pantalla nunca transforma accepted/delivered/read en venta, cita
aceptada por cliente, autorización de reenvío ni resultado empresarial.

El host bloqueante no instala un daemon. Revalida identidad en cada poll;
auth cinco segundos, reporte tres, jitter/backoff máximo diez minutos y una
toma de job por iteración. Cancelación cooperativa: los callbacks deben respetar
context. Un reporter fallido detiene Run; el supervisor real debe alertarlo.
Los reports contienen outcomes acotados, no payload/contacto/token/ID de proveedor.
No se afirma destino de alertas, supervisor o identidad de servicio desplegados.

La colección reutiliza la consulta existente dentro de una transacción
repeatable-read/read-only; límite de veinte aprobaciones y conflicto explícito
si hay más, nunca historia parcial presentada como completa.
El navegador no realiza sondeo automático, cancela lecturas obsoletas y elimina
el éxito anterior ante fallo/scope divergente. Consulta y fecha del proveedor
se muestran por separado.

Activación: features.whatsapp_status_history=true en el owner de configuración
de negocio existente, después de montar y probar el módulo Go. Ausente/false
oculta la vista y bloquea su BFF antes de resolver sesión/backend. Configuración
inválida falla cerrada. No crea otro registro de features ni concede permisos,
consentimiento o acceso al proveedor. No fuerza TypeScript como backend/base;
es el adapter web opcional ya presente en la arquitectura.

## Gates ejecutados

| Gate | Resultado y evidencia conservada |
|---|---|
| Preflight | BLOCKED: dotnet/docker/psql no descubiertos; process 0 no equivale a READY. PostgreSQL se invocó por su ruta real en pruebas focales. preflight.log/json |
| Composición limpia | 67 packs / 742 archivos, 52 migraciones existentes; verified-composition.log |
| Identidad de artefacto | 22/22 rutas idénticas SHA-256 respecto del árbol corregido; tabla siguiente |
| Go/WhatsApp + PostgreSQL | 43 tests top-level PASS, cero SKIP; incluye cuatro proyectos navegador y cinco tests host. verified-whatsapp.log |
| Web | 102 PASS / 1 SKIP en 12 archivos; typecheck y build PASS. verified-web/typecheck/build.log |
| BFF conectado omitido en suite web general | Ejecutado aparte: PASS con RS256/Go/PG, duplicate=409, audit=1/outbox=1 y lectura customer-scoped. verified-agenda.log |
| Agenda existente | Cuatro proyectos navegador PASS y negativos de firma/permisos/org/tenant/recurso; verified-agenda.log |
| Historial de notificación | Cuatro proyectos PASS: pérdida de GET, recuperación, respuesta ajena rechazada, ayuda por versión, fecha histórica visible, cero POST; verified-whatsapp.log |
| Postcondiciones de notificación | Por fixture navegador: un job completado, una observación, attempt_count total uno |
| Python | 42 PASS; verified-python.log |
| Go vet/build | Ambos process 0 sobre ./...; verified-vet.log y verified-go-build.log |
| Playwright público histórico | Runtime 4 PASS, home/CSP 8 PASS; 12 SKIP opt-in conservados, no contados como PASS. verified-public-browser.log |
| Execution validator | 24 tests PASS con ELITE_AUTHORITY_ROOT real; execution-validator-green.log. No es checkpoint ni readiness de un proyecto productivo |
| Lighthouse público | Cinco mediciones: performance mediana/mínimo 0.99; accesibilidad automatizada, buenas prácticas y SEO 1.00 en todas; 155 audits/run, cero warnings, cinco tests de política PASS. lighthouse-target-green.log y lighthouse-results |

La prueba conectada de BFF usa filtro y deja otros 16 tests fuera de esa
invocación; ya se ejecutaron en la suite web general. No se suman subtests ni
pruebas repetidas para inflar un total. Los ocho recorridos de navegador son
cuatro configuraciones por journey, no ocho proveedores ni plataformas móviles
nativas. Capturas verificadas en desktop/mobile, conservadas en
notification-browser-artifacts; no equivalen a auditoría WCAG/AT completa.
Advertencias TLS del certificado sintético loopback y NO_COLOR/FORCE_COLOR
se preservan. No se relajó TLS productivo ni se ejecutó SCA actual en esta ronda.

## Fallos y recuperación

FAIL-20260906-382: inicio PostgreSQL sin puerto explícito; conexión rechazada
antes de migraciones. Reinicio exclusivo del data directory de prueba en
127.0.0.1:55959, readiness y 52 migraciones verificadas. Recurrencias de tooling
incluyen CLI Playwright aislado no instalado (LIB-FAIL-1119), SourceRoot con slash
final, rutas/contextos inferidos, patch con operaciones repetidas para un archivo
y capstones sin raíz autoridad. Se corrigió el entorno/llamada, no se relajaron
guards, allowlists ni se atribuyó éxito a intentos fallidos.

FAIL-20260906-383: tres navegadores de la agenda fallaron al encontrar dos
role=status sin nombre distinto. round-agenda.log conserva el rojo; corrección
con nombres accesibles específicos y selección role+name, no por índice.
final-agenda.log prueba cuatro navegadores con ambos paneles presentes; la
versión final además nace desactivada por configuración y verified-agenda.log
repite la agenda. verified-whatsapp.log prueba la vista activada junto a agenda.
Se preservan todas las ejecuciones previas; no se reescribe su resultado.

## Autoridad de método, no atribución de código

- [Go: pipelines y cancelación](https://go.dev/blog/pipelines): cancelación
  cooperativa y responsabilidad de consumidores/productores.
- [React: sincronización de effects](https://react.dev/learn/synchronizing-with-effects):
  cleanup y protección frente a resultados de una lectura anterior.
- [Google SRE: monitoreo](https://sre.google/sre-book/monitoring-distributed-systems/):
  señales operables; no justifica alertas inexistentes.
- [Google SRE: gestión de incidentes](https://sre.google/sre-book/managing-incidents/):
  coordinación y registro; la guía local no certifica personas capacitadas.
- [Microsoft Playwright: buenas prácticas](https://playwright.dev/docs/best-practices):
  verificaciones visibles para el usuario, aislamiento y datos controlados.

Consultadas 2026-09-06. Estas fuentes gobiernan claims estrechos de método;
no se afirma equivalencia universal con sus sistemas internos ni perfección.

## Hashes de las rutas cambiadas

| Ruta | SHA-256 |
|---|---|
| internal/whatsappbridge/status_host.go | bddcea35fd820d71e3fd62e49fc8f9f7e25a3947e344b374cde258e95baaf853 |
| internal/whatsappbridge/status_host_test.go | 9f8195112b929171068de6e7df41a2a0cc8f4421471059c15af20a8854e988c8 |
| internal/whatsappbridge/notification_history.go | 620e8d645eb9f711658afc1be3d61b4fce5ea626040dffb659c92d43de66b4ca |
| internal/whatsappbridge/notification_history_test.go | fdba06ff3e9224bd0504496aaed6c9fc3413a706fb182dc577bd90a8d6674cf8 |
| internal/whatsappbridge/notification_browser_test.go | 165507a3c0703de4d08beab732100d1f466fa9daa301ee23650b86548ba1dba0 |
| internal/whatsappbridge/appointment_notification.go | 939cafb14e13a689c415ad40e8e3e6c25336d732591fd5b90d702d876af9f8bb |
| internal/whatsappbridge/appointment_notification_status.go | 9bb6b30dda597ee96b784423de72eebd08921fec08f16f5c360a318e93d6d948 |
| whatsapp_cloud/README.md | 5e637e2b56e3c6ebc31d5d6de94f05698368e5f450896c7c74343f8dbd945a0b |
| whatsapp_cloud/PROVENANCE.md | 28399954d1a93da98983a2ad29c5001cdaef6cfa376100af9ffda046e161ac69 |
| src/platform/notifications/status-contract.ts | 501130c6f091a1e007d4c6af585cb8acb0ea2fab95f2eadc786329db41ca3fd7 |
| src/platform/notifications/status-contract.test.ts | 1c6b7eacbb277fe3c66e635852c9ffbd6c0b65a98e8892e0028d145339464f62 |
| src/app/api/enterprise/franchise/notifications/route.ts | ebfe5ab0ca8a45c3601666e2248e894bf1acd7bbf2a7a61e0c8cf813be438455 |
| src/app/api/enterprise/franchise/notifications/route.test.ts | b1cac0a7bc3605e88dab51b33d0f044495de55cd5de52079cdba1b7afdf25238 |
| src/components/appointment-notification-status.tsx | 8a1942f6968e4abbc4718686d212f5d0b1f6812c178e117484958ddfdf250678 |
| docs/whatsapp-status-operations.md | 5619c6d4d1791107592dccd99334d3758039d7e7200b5672fa2299036ece62a8 |
| src/components/franchise-command-panel.tsx | 4e3bdd1824f501757a3e92a57b24e25ae724ad2b16e8e984427361be9d115121 |
| src/app/franchise/page.tsx | dc80d1bc1a75b450003c8ae25b1558c477c32e9dddb77229e61b3a7908ed347e |
| microsoft_playwright_browser_gate/tests/notification-status.spec.mjs | 43c1dd61f0df5ddd3e6c1a06c36104170256da54e9d0d4243a0222b771f00c4a |
| microsoft_playwright_browser_gate/playwright.config.mjs | 1b92f88fc77bf577efe1589ab8ef64cadfa3e85ee63cde014df6d5beee277109 |
| microsoft_playwright_browser_gate/README.md | 9ddf554a06c090681a83bc91fc7438b252bcbb7d05adf47ec8f69db4c809e5e6 |
| microsoft_playwright_browser_gate/tests/enterprise-web.spec.mjs | 0c42ba1e94801537a0d37c2fc5684ca2d2d1adb1ee98cf45a5d7ddc775a0b9b0 |
| src/app/globals.css | b7964d15c39cd463b57f522409e07b4592a19fdcd900da05909eae2bb2c796aa |

## Condiciones abiertas y próximo paso

Aún requieren implementación/configuración y evidencia: montaje del módulo y
supervisor en el target, identidad de servicio legítima, reporter/alertas reales,
retención/PII, capacidad y redrive autorizado; inbound de clientes/autorespuestas,
otros journeys del roadmap, contratos/cuenta Meta/IdP reales y entrega live.
La ayuda cubre esta vista; no completa soporte, capacitación ni evaluación de
todos los roles. No se cierra publisher/outbox ni todo el sistema de franquicia.

El target conserva CDN/WAF, seguridad ofensiva, SCA actual, carga, PITR/recovery,
deploy/canary/rollback y aceptación empresarial. No se emite READY_TO_BUILD,
porcentaje global, plazo de franquicia ni ZIP final con esta evidencia.
Rollback de este delta: desactivar la feature y pausar host; preservar jobs,
approvals, observaciones, audit y fences; restaurar artefacto compatible probado.
Nunca borrar historia ni resetear attempts para simular recuperación.

La próxima integración debe montar el runtime y el reporte en el owner de
aplicación existente con accesos legítimos o mantener explícita su condición;
no copiar la identidad sintética de tests ni crear otra cola/registro.

## Control raíz y limpieza

VERIFY_LIBRARY_PASS, proceso 0: 160 packs / 1.433 archivos materializables /
711 Markdown / 51 perfiles; franquicia 67/742. AUTHORED=1195, ADAPTED=133,
VERBATIM=105. Ledger: 2.117 IDs locales + 209 upstream = 2.326. No Audit
ejecutable global ni SCA actual repetidos. Logs verify-library.log,
verify-library-final.log y verify-library-closure.log; este último incluye
la alineación final de metadata compatible_with sin cambiar ningún hash de código.

Se comprobó identidad exacta y cero sesiones ajenas antes de eliminar sólo
elite_whatsapp_v278 (478 tenants, 170 jobs, 149 webhooks sintéticos) y
elite_confirmation_v278 (15 tenants, cero jobs/webhooks). Después de comprobar
cero otros clientes se detuvo el PostgreSQL de pruebas iniciado en esta ronda.
Datos reconstruibles mediante fixtures/migraciones; árboles, logs y capturas
permanecen. No se eliminó información del usuario ni se activó servicio externo.
El servidor Next/Chrome de Lighthouse también terminó mediante su cleanup.

El primer helper Lighthouse falló por quoting de ArgumentList antes de iniciar
el target; lighthouse-target.log y lighthouse-next.stderr.log lo conservan.
Se corrigió la invocación, sin cambiar código web, pins, política ni floors;
lighthouse-target-green.log acredita cinco mediciones PASS. Son mediciones del
home público local, no del portal autenticado ni del despliegue productivo.
