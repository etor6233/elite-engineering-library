# V302 — permisos y recuperación de consultas del portal operativo

> **Errata de trazabilidad V317 (2026-09-08), prevalece sobre las menciones históricas de Go1.26.7 debajo.** Se retira la atribución histórica exacta Go1.26.7 no vinculada a un receipt de identidad por ejecución. No se sustituye por1.26.8 sin prueba del artefacto concreto. La ruta current-toolchain de V281 contiene1.26.8; official-toolchain contiene1.26.7. El cuerpo original queda conservado como historia. V317 revalida únicamente la composición actual y el scope allí enumerado con1.26.7 observado, no todos los snapshots previos. Ver [TOOLCHAIN_IDENTITY_AND_RUNTIME_SCA_V317.md](TOOLCHAIN_IDENTITY_AND_RUNTIME_SCA_V317.md).

2026-09-07. Continuación del checkpoint33/V301, T2803/T2804 del roadmap
existente. Mantenimiento correctivo de owners admitidos bajo sus condiciones;
no nueva capability, proveedor, regla comercial ni promoción de producto.
V295, V300, V301 y ARCA diferida permanecen con sus límites.

## Implementación y procedencia

Cuatro archivos existentes regresan a sus packs canónicos:

- TYPESCRIPT_FRANCHISE_JOURNEY_PORTALS0.14.0: page.tsx y
  franchise-command-panel.tsx. Secciones y formularios se seleccionan por los
  permisos ya consumidos por HTTP; acceso también para resource:manage y
  availability:read sin exigir permisos comerciales ajenos. Read-only no puede
  enviar el formulario de modificación. La cancelación permanece deshabilitada
  sin availability:manage. Sin permisos suministrados, controles cerrados.
- GO_FRANCHISE_CUSTOMER_JOURNEY_API0.10.7: sólo harness de prueba; no cambios
  al dominio, SQL de producto, permisos HTTP ni migraciones.
- MICROSOFT_PLAYWRIGHT_BROWSER_GATE0.1.12: regresión conectada de roles y fallo
  aislado de consulta; no nueva dependencia ni renovación de upstreams.

Las seis consultas de la página conservan ejecución concurrente, pero cada
fallo tiene resultado separado. Una consulta de discrepancias fallida no quita
los pedidos. Se muestra un error comprensible sin detalle interno y un GET
«Volver a consultar operación»; no reenvía comandos. Las secciones fallidas no
se presentan como listas vacías exitosas. La ayuda contextual versionada
operation-sections-view/1.0.0 explica recuperación, práctica y soporte autorizado.
No equivale a una plataforma de capacitación, progreso ni soporte integral.

Los cambios son AUTHORED / LicenseRef-Workspace-Owner, no código copiado de
Microsoft/Meta ni avalado por esas empresas. Se consultó el2026-09-07
[React: conditional rendering](https://react.dev/learn/conditional-rendering),
sólo para composición JSX. La autoridad de autorización es el contrato y el
handler HTTP existentes, no el renderizado. React y Playwright conservan los
locks/licencias fijados; reputación o documentación no promueven estos cambios.

## Reproducción y evidencia real

Staging: `%LOCALAPPDATA%/Temp/elite-v302-b1f9724d58bf48bfa8fecdd7cb79ef21`.
Composición integral67 packs/745 archivos, destino ausente final-source; los
cuatro archivos coinciden por SHA con candidate, cuyo Next se compiló. Perfil
de backend/web/integral/serverless sincronizado sin agregar rutas duplicadas.

1. sections-baseline.log: V300 construido muestra «Preparación versionada de
   entregas» a un lector de leads; Chromium esperaba0 y recibió1. Rojo real.
   El fallo propagado de Promise.all se observó en código, no fue un segundo
   rojo de navegador: la primera aserción interrumpió esa ejecución.
2. browser-closure.log: TestQuoteAcceptanceBrowserPostgres PASS,52 fases
   (13 por proyecto), Chromium desktop/mobile, Firefox desktop y WebKit desktop;
   148.755s paquete. Next/BFF→Go→OIDC sintético→PostgreSQL18.6 real en loopback.
   Conserva cotizaciones, pedidos, reservas, solicitudes de pago, entrega,
   reinicio, duplicados, concurrencia y negativos de las fases anteriores.
3. Cada fase nueva prueba cuatro roles: lector de leads, gestor de entrega,
   gestor de recursos y lector de disponibilidad. Primer GET de discrepancias
   devuelve503 inyectado; pedidos siguen visibles, detalle privado no aparece,
   GET posterior recupera la sección. Cero POST automáticos durante navegación
   y recuperación. No mide todas las combinaciones de permisos ni accesibilidad
   completa; comprueba overflow horizontal del viewport de esa fase.
4. Una sesión de prueba con permisos UI forzados no altera su access token:
   publicación de checklist por BFF y consulta directa no autorizadas reciben403.
   SQL en sales.delivery_checklist_template demuestra cero filas denied-probe.
   Cada fase conserva3 pedidos/3 aceptaciones/outbox6; operación conserva2
   reservas/2 asignaciones/outbox2. Fixtures locales, no pagos o entregas live.
5. agenda-final.log: TestAppointmentAgendaBrowserPostgres PASS en los mismos
   cuatro proyectos;13.427s paquete. Repetida por cambio de la página compartida.
6. Web:105 pruebas PASS,1 SKIP,12 archivos; build PASS. Go test ./..., vet y
   build ./... sobre final-source PASS, GOPROXY=off/GOTOOLCHAIN=local, Go1.26.7.
   Las suites opt-in PostgreSQL/browser se acreditan por2–5, no por tests omitidos.

Repetir desde composición limpia, toolchains fijados y PostgreSQL local con las
53 migraciones existentes aplicadas: instalar web/gate con lock y caché verificada,
ejecutar tests/build web y las suites Go. Para el gate conectado, definir
TEST_DATABASE_URL del fixture local, ELITE_WEB_ROOT absoluto a la web construida,
ELITE_QUOTE_E2E=1, ELITE_ORDER_E2E=1, ELITE_PAYMENT_E2E=1,
ELITE_DELIVERY_READ_E2E=1, ELITE_DELIVERY_ACTION_E2E=1 y
ELITE_OPERATOR_SECTIONS_E2E=1; ejecutar `go test ./internal/platform/httpapi
-run '^TestQuoteAcceptanceBrowserPostgres$' -count=1 -v -timeout 15m`.
El harness crea identidades sintéticas, inicia los servidores y comprueba SQL.
Nunca usar una base productiva para estos fixtures.

## Fallos de ejecución conservados

- FAIL454: primer arranque de pg_ctl usó5432; gate51913 no conectó. No es rojo
  del producto. Reinicio del cluster propio con host/puerto explícitos y
  pg_isready seguido de las suites conectadas PASS.
- FAIL456: primer SQL de evidencia tenía escape erróneo42601; segundo usó tabla
  supuesta42P01. Se consultó el owner PublishDeliveryChecklist, se parametrizó
  denied-probe y se corrigió a delivery_checklist_template. Los logs
  sections-candidate.log y browser-final.log permanecen FAIL. No se eliminó la
  aserción durable; browser-closure.log es el único cierre de la matriz52.
- FAIL455 queda REGRESSION_PROVEN con las pruebas anteriores.
- FAIL457 queda OPEN: FranchiseCommandPanel.send todavía afirma «No se aplicó»
  y permite reenviar tras respuesta perdida. Es siguiente reparación T2804,
  con post-commit, fence y consulta del resultado. V302 no la demuestra ni oculta.

## Hashes de los bytes probados

| Archivo/log | SHA256 |
|---|---|
| src/app/franchise/page.tsx | b16fee2fd692c81dcddcef69e8e4f185b665da04d3288d55ee23b2fd288f3360 |
| src/components/franchise-command-panel.tsx | c6e3a9db57c9deb4576e59067847894dac1a4b6d5b3ae11999319fba89a5645d |
| internal/platform/httpapi/franchisejourney_test.go | a735a9237b68982ce9cff83a4f35ee7a13176c9eac818298a34d66434c562620 |
| microsoft_playwright_browser_gate/tests/enterprise-web.spec.mjs | 4d773764f396c3c18ef8786cec4a54b4057d2d9a6e9ffbb06cb4f37fabe1794a |
| browser-closure.log | 81ddecdf5658b6e7293ad7d6712799b57f93ba620f278171217e7d1aa80f6bad |
| agenda-final.log | 00fa1ed82fcde6e861c4371e1caed60e9905a4afda557ee233a51ee2f1434095 |
| sections-baseline.log | 2562c555c4bc5843ba085a69b1076a2b1481771f95262eca2eca87a05e2f91cb |
| sections-red.log | f0f6879bcf6f3f11d096bfa0b79a2fefdd3b6de69bdf5b700163d50494cd4fcc |
| sections-candidate.log | fd3aefa65b89fe70e9429c45184499277ead8c2b0ff0d0d9753e0d5c4328ff95 |
| browser-final.log (FAIL preservado) | ecc14b6bd11f0160717c6b3eec31cea8f08894da0d3add5b8b37a43a86b1e39b |
| web-test.log | 352b586e197f79720a9b249d5c4c5f79cb2e7969c18e68652344682dd3fa8b47 |
| web-build.log | 24c8abeac760e13fa683826833161c92edd043fc6683cb7b66623cdfb910bb77 |
| final-go-test.log | b62ade0099ed79298e9da309a98b8ead6c5eb8fb0632fa3d4f524a0dfdebd033 |
| final-go-vet.log / final-go-build.log (vacíos, exit0) | e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 |

## Lo que no cierra

No nuevos módulos, migraciones, dependencias ni archivos de producto. No creador
inicial de entrega, liberación comercial, cobro real, transportista, credenciales
productivas, ARCA, observabilidad integral, SCA nuevo, carga, PITR, deploy,
rollback ni aceptación empresarial. No validación de toda respuesta malformada:
el fallo nuevo ejercitado es503. No cierre T2802/T2804 completo ni assurance.
No nueva estimación porcentual o promesa de plazo, ZIP ni REUSABLE_PACK por PASS.
