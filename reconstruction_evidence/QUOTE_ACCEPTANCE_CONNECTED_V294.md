# V294 — aceptación de cotización conectada y recuperación

> **Errata de trazabilidad V317 (2026-09-08), prevalece sobre las menciones históricas de Go1.26.7 debajo.** Se retira la atribución histórica exacta Go1.26.7 no vinculada a un receipt de identidad por ejecución. No se sustituye por1.26.8 sin prueba del artefacto concreto. La ruta current-toolchain de V281 contiene1.26.8; official-toolchain contiene1.26.7. El cuerpo original queda conservado como historia. V317 revalida únicamente la composición actual y el scope allí enumerado con1.26.7 observado, no todos los snapshots previos. Ver [TOOLCHAIN_IDENTITY_AND_RUNTIME_SCA_V317.md](TOOLCHAIN_IDENTITY_AND_RUNTIME_SCA_V317.md).

Fecha: 2026-09-07. Owners: roadmap T2802/T2804, sin nuevo plan.
Estado: **PASS local del claim acotado**, packs CONDITIONED. No cierre integral de franquicia.

## Resultado ejecutado

Portal cliente real Next/BFF → cookie cifrada → OIDC/JWKS RS256 sintético verificado
por Go → repositorio existente → PostgreSQL 18.6 con 52 migraciones.
No se sustituyó el repositorio por memoria ni se modificó la transacción de negocio.

Se reprodujo un defecto real: el servidor creaba el pedido, pero una respuesta
perdida hacía mostrar «No se aceptó». La corrección AUTHORED mantiene resultado
incierto, deshabilita reenvío, recupera por GET y muestra la referencia durable.
La aceptación **no confirma pago, stock ni entrega**.

En la reconstrucción canónica pasaron 4 proyectos (Chromium escritorio/móvil,
Firefox, WebKit), cada uno con 3 fases: initial, read-failure, restart.
Total: 12 fases PASS; Go informó 38.906 s para el gate completo.
Cada tenant termina con 3 pedidos placed, 3 líneas, 3 aceptaciones y 6 eventos
outbox (3 quotation.accepted y 3 customer-order.placed), más dos cotizaciones
no aceptadas. El runner verifica los contadores después de cada fase.

- Aceptación normal por UI y activación por teclado.
- Respuesta perdida después del commit real; lectura muestra el pedido sin repetir POST.
- Dos aceptaciones simultáneas: 200/409, un solo pedido para esa cotización.
- Reenvío explícito: 409; no se promete respuesta 200 idempotente. La recuperación es GET.
- Rechazo de versión incorrecta, vencimiento, organización/tenant/cliente ajenos,
  permiso ausente y sesión ausente; sin pedido adicional.
- Fallo GET 503: mensaje comprensible, sin tarjetas engañosas, recuperación y ayuda.
- Reinicio real del proceso Next y recreación API/pool; lectura conserva referencias.
  No equivale a reiniciar el proceso Go completo, PostgreSQL, PITR o DR.
- Guía contextual quote-acceptance-view/1.0.0, práctica segura y escalamiento sin tokens.
- Capturas y ancho de viewport móvil comprobados. No auditoría WCAG completa,
  lector de pantalla, validación UX empresarial ni LMS/progreso/evaluación integral.

## Procedencia, reconstrucción y límites

Cinco archivos vuelven a los tres owners existentes:
GO-FRANCHISE-CUSTOMER-JOURNEY-API 0.10.1,
TS-FRANCHISE-JOURNEY-PORTALS 0.11.1,
MICROSOFT-PLAYWRIGHT-BROWSER-GATE 0.1.7.
67 packs/742 archivos/52 migraciones, sin nueva dependencia.
Recomposición en destino ausente: 5/5 archivos byte-idénticos.
Instalación offline/frozen-lockfile; el gate anidado necesita --ignore-workspace.
Build Next PASS; Go test ./..., vet y build PASS. Tests web: 102 PASS, 1 SKIP
opt-in, 103 total. Las suites generales sin TEST_DATABASE_URL no acreditan PostgreSQL.

Código local AUTHORED, no código copiado ni aprobado por Microsoft/Google.
Frameworks y licencias mantienen sus locks existentes. Autoridades del claim:
[Playwright API testing](https://playwright.dev/docs/api-testing) para combinar
acciones de navegador con pre/postcondiciones, y
[PostgreSQL 18 locking](https://www.postgresql.org/docs/18/explicit-locking.html)
para los límites de concurrencia. Citar estas fuentes no readmite el producto.
No se adquirió código upstream nuevo ni se ejecutaron proveedores, cobros o cuentas.

FAIL-420 cerrado para la falsa respuesta negativa. Fallos de fixtures/entorno
418/419/421/422/424 corregidos con logs conservados.
FAIL-423 permanece abierto: Firefox falló una vez durante cierre de contexto;
la repetición focal y la reconstrucción completa pasaron sin cambiar pin ni
ocultar/reintentar automáticamente el fallo. No se encontró confirmación oficial
de la causa exacta. Un PASS posterior no borra esa inestabilidad.

Revisión visual deja mejoras pendientes de presentación: estados técnicos,
importe en unidades menores y fecha ISO; el servidor rechaza vencidas aunque
la tarjeta issued aún ofrece aceptar. No declarar experiencia final simple.
T2802/T2804 no se cierran completos. Continuar pedido→stock/pago/entrega usando
sus owners; conservar capacitación/soporte y pruebas operativas pendientes.
Readiness41 y assurance integral siguen BLOCKED; no promoción ni ZIP final.

## Receipts

Staging local bajo el temporal del sistema:
elite-v294-8833fd547ff1493b9d5cb48f97234338. Datos sintéticos aislados.
Los logs rojos se preservan; un nombre green no sustituye el exit code.

| Archivo o receipt | SHA-256 |
|---|---|
| quote-baseline-ux.log | ea8e93e929f758267cb8f28e2d255c183e10592a03f5d750612ae08461d09720 |
| quote-rebuilt.log | bcff436ec94998e90b3fe5241fcc11da6bbfbdeca23c9848aa5028a321c79515 |
| quote-rebuilt-corrected.log | 6c724c489fd4db5780c26e34c8577006848b9d233fba3a981433c336a9861cab |
| quote-firefox-repeat.log | 9e6f870c858f24f2082fb9012b87d5c435bb52226ff8687ad2b7cff94192c15f |
| web-tests.log | e9992c547a7baee3bcf302928c4fecaf3a0d907ed759059f682e5233f0869d70 |
| go-unit.log | 958934d50d9252705992fe1e4f632449c56bc8419c289876909376e1ff236c18 |
| go-vet.log | e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 |
| go-build.log | e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 |
| rebuilt-web-build.log | f1f68fbb822f598a4a9046a12033dd22a2efd9c24f9d384aee8b368009a7366c |
| migrations.log | 62a1fd47c06cb823a5ef000e52eb60326fc7f696fe68f6300786b7acfa7cd6f7 |
| internal/platform/httpapi/franchisejourney_test.go | 23960103a415066fbe0ed546c969308520907cb24704deceb227526ff6807458 |
| src/components/customer-quote-actions.tsx | 3ea5223460bf368830fcdc827d776853e5a684cff5bcb8f80d196f5aca7b097f |
| src/app/customer/quotes/page.tsx | 2fc5171ea9156a95333ddf8703656d250425d1a0811bb89d8d0f666b6bd49f79 |
| microsoft_playwright_browser_gate/tests/enterprise-web.spec.mjs | ebea5b35a5e02ca284952240681bd06ecbe2fa489fac094d137d958c7897f476 |
| microsoft_playwright_browser_gate/playwright.config.mjs | 83f4ae8671dbd161c96068ea2b3e3fd1d8ebe12d3b6dc30ef580d56a6b535889 |

## Continuación después del checkpoint20 — misma vertical y siguiente dependencia

Revisión final: GO-FRANCHISE-CUSTOMER-JOURNEY-API **0.10.2**,
TS-FRANCHISE-JOURNEY-PORTALS **0.11.2**, browser gate **0.1.8**.
Los receipts anteriores son históricos del primer corte, no hashes de esta revisión.

Se continuó sin crear otro plan ni owner. Los pedidos generados por el navegador
se conectan directamente con el owner Commerce existente:

- Dos pedidos compiten por una misma unidad serial: exactamente una reserva.
- Tenant/org ajenos, versión incorrecta y asignación duplicada rechazados.
- Una solicitud payment_attempt en estado **created**, no authorized/captured.
- Importe/moneda/org incorrectos y clave duplicada no agregan pagos.
- Pool nuevo confirma reserva, stock, línea asignada, solicitud y sus eventos durables.
- Cero llamadas al proveedor y ningún worker de pago habilitado.
- Esto prueba contratos del repositorio sobre los mismos pedidos; **no prueba aún
  UI operativa/BFF/OIDC de stock/pagos ni una venta cobrada**.

La presentación inicialmente pendiente se corrigió: importe ARS 2.500,00,
estados en español, fecha legible UTC y cotización vencida sin botón.
Cálculo visual server-side con BigInt/Intl, sin convertir moneda ni definir impuestos;
importes no representables de forma segura o moneda no reconocida no permiten
aceptación. La captura móvil final fue inspeccionada. Pruebas genéricas de todas
las monedas, reader accessibility y aceptación empresarial siguen fuera del claim.
Autoridades de presentación: [ECMA-402](https://tc39.es/ecma402/#sec-currencydigits)
y [React hydration](https://react.dev/reference/react-dom/client/hydrateRoot).
Implementación AUTHORED; no se atribuye esta función a esas organizaciones.

Nueva composición vacía final-rebuilt: cinco archivos byte-idénticos.
Build Next y Go test/vet/build PASS; web 102 PASS/1 SKIP.
Gate final: cuatro navegadores, tres fases cada uno y cuatro continuaciones
Commerce PASS; 34.953 s. No actualización de dependencia ni reintentos ocultos.

Después, PostgreSQL se reinició de forma controlada. Snapshots ordenados de
pedidos, líneas, aceptaciones, stock, solicitudes de pago y outbox idénticos:
SHA-256 **cb2e03247bb35b24b105f77cf5f0d743e6c4cece239c2debe6fa9461f16cb80f**.
Logs pg-before-restart.log y pg-independent-after-restart.log conservados.
No crash test, backup, PITR ni DR. El cluster propio quedó detenido, datos conservados.

El primer orchestration de restart retuvo pipes; el segundo Start-Process -Wait
esperaba al árbol de descendientes. Regresión final: Hidden, redirecciones -l,
host/puerto explícitos y espera finita de pg_ctl (no de PostgreSQL), start→probe→stop PASS.
Fuentes: [Microsoft Start-Process](https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.management/start-process)
y [PostgreSQL pg_ctl](https://www.postgresql.org/docs/18/app-pg-ctl.html).
FAIL-429 registra los intentos fallidos; no se cuentan como pases.

Pendiente real: portal operativo de asignación/pago/entrega conectado, proveedor
sandbox/live, autorización/recuperación de esa continuación por HTTP, ayuda y
capacitación por rol, operación/seguridad/aceptación target. Firefox423 sigue
registrado como intermitente no explicado. Readiness/assurance no se promueven.

### Receipts finales (sustituyen hashes de código del primer corte)

| Archivo o receipt | SHA-256 |
|---|---|
| quote-commerce-rebuilt.log | efc29bd9e5fa13133687dac6454957cfb2a2edaf58c6a8a820f1fd853f4c6dee |
| quote-final-rebuilt.log | 5a1d15f86f2a97aed370d5a3ba3e2cdae05cf9b212b3357d8c50a649eefc1632 |
| final-build.log | 319cc7e2505ca96b1a1bdbf00d1d2ba1eb4e35caf0bcfc3104ee08a67942b603 |
| final-go-tests.log | 3e1513e88f6ceec8831f2b4f9bc63a7afc09e2776843d811adf86270682f3703 |
| final-go-vet.log | e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 |
| final-go-build.log | e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 |
| final-web-tests.log | 04826ffd43fb29b1a2e815f10bb2ba5d59fc5665cca341ab31449298d382aaa7 |
| quote-ux-build.log | 6aa4e850fe4622b97f4a63497d7b6644facb27cbce018aae581b7588c305153e |
| internal/platform/httpapi/franchisejourney_test.go | 920d420324555e17c23eeff97df0ef15f8bcb8a1668896314c9794cfe3d1206c |
| src/components/customer-quote-actions.tsx | ace08dab683e987f925382351e98c9b9bc363aaac01be3bf0075a1c0c5616e46 |
| src/app/customer/quotes/page.tsx | a776ade4087bb40d8db326dd411864f22f1acf98704ed29376d2ff500c8390aa |
| microsoft_playwright_browser_gate/tests/enterprise-web.spec.mjs | d1c65b7a0aac948594b43c809a77254e3f97d51a2edd16e72fdf2c6f55ea5a94 |
| microsoft_playwright_browser_gate/playwright.config.mjs | 83f4ae8671dbd161c96068ea2b3e3fd1d8ebe12d3b6dc30ef580d56a6b535889 |

