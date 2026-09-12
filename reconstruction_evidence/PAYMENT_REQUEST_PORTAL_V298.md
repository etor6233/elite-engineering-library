# V298 — solicitud de pago desde el portal operativo

> **Errata de trazabilidad V317 (2026-09-08), prevalece sobre las menciones históricas de Go1.26.7 debajo.** Se retira la atribución histórica exacta Go1.26.7 no vinculada a un receipt de identidad por ejecución. No se sustituye por1.26.8 sin prueba del artefacto concreto. La ruta current-toolchain de V281 contiene1.26.8; official-toolchain contiene1.26.7. El cuerpo original queda conservado como historia. V317 revalida únicamente la composición actual y el scope allí enumerado con1.26.7 observado, no todos los snapshots previos. Ver [TOOLCHAIN_IDENTITY_AND_RUNTIME_SCA_V317.md](TOOLCHAIN_IDENTITY_AND_RUNTIME_SCA_V317.md).

2026-09-07. Continuación del checkpoint29 y T2802/T2804, sin nuevo roadmap.
V295, V296, V297 y ARCA diferida conservados. Cierre **local condicionado**;
no es un cobro, una entrega ni una aprobación productiva.

## Implementación

Doce archivos existentes, cero archivos de producto nuevos, cero dependencias
nuevas. Perfil67/745. Owners: Commerce0.6.0, aplicación1.9.2, Journey0.10.4
(sólo harness), portales0.13.0 y browser gate0.1.10. Todos los cambios son
AUTHORED; no se presentan como código de Stripe, AWS, Meta o Microsoft.

- La misma pantalla de pedidos permite registrar una solicitud inicial del total.
  El servidor deriva moneda/importe bajo lock y audita el actor autenticado.
  BFF strict/same-origin rechaza importe/proveedor/actor aportados por el cliente.
- La clave UUID se conserva en sessionStorage bajo scope hash de tenant/usuario/
  organización y pedido. Sin storage no hay envío. No guarda tokens ni secretos.
  Una respuesta incierta conserva la clave y exige consultar el estado durable.
- Tanto la ruta nueva como la anterior usan exclusión transaccional del intento
  inicial: otra clave/proveedor para el mismo pedido da409, incluida concurrencia;
  replay válido recupera el mismo recibo. No se elimina historia para reintentar.
- PAYMENT_REQUEST_PROVIDER ausente/vacío desactiva requests/transiciones HTTP y
  acciones UI; stripe/mercadopago son la allowlist. La UI usa la selección del
  backend, no otra configuración. Valor inválido falla cerrado. El fixture usa
  stripe sin cuenta, SDK ni llamadas de cobro. La guía de credenciales explica
  exactamente ese consumer y distingue configuración de acceso real.
- Ayuda, práctica segura y escalamiento quedan en order-operations-view/1.1.0.
  No se infiere dinero recibido ni autorización de entrega de un estado created.

Fundamentos oficiales consultados, no aprobación del código local:
[AWS Builders’ Library](https://aws.amazon.com/builders-library/making-retries-safe-with-idempotent-APIs/)
para identificar intención/reintentos y
[React](https://react.dev/learn/you-might-not-need-an-effect) para ejecutar la acción
del usuario desde su evento, no desde un efecto de render. Se reutilizan los
locks y licencias existentes; no hay una adquisición de upstream nueva.

## Evidencia ejecutada

Staging: `%LOCALAPPDATA%/Temp/elite-v298-7384649c848f4383b34d7f1017b9ecd7`.
`final-source` recompuesto desde Markdown;12/12 hashes iguales a candidate y
runtime probado. Se reutilizó PostgreSQL18.6/schema53 demostrado en V297;
sin nuevas migraciones, datos reales, cuentas ni efectos externos.

- Go raíz: test/vet/build PASS; PG focal explícito PASS. Tests opt-in no se
  presentan como ejecutados por la suite raíz cuando no tienen entorno.
- Web:105 PASS y1 SKIP opt-in; build Next PASS. Artefacto web se reutiliza después
  del último fix sólo Go/test porque sus archivos no cambiaron.
- Playwright final:36 fases, nueve en cada Chromium desktop/mobile, Firefox y
  WebKit.24 revalidan cotización/stock afectados y12 cubren pago/desactivado/
  procesos nuevos. Duración observada98.35s, no estimación del proyecto.
- Mismos pedidos creados por cliente→reserva→solicitud; identidad RS256/JWKS local,
  sesión cifrada, BFF, HTTP, owner y PostgreSQL. Dos intents/auditorías por tenant,
  importe250000ARS derivado del pedido y actor payment verificados en SQL.
- Respuesta abortada después del commit real, GET recupera ID; replay misma clave
  conserva ID. Ocho claves concurrentes dejan un éxito y siete409. Negativos:
  actor/provider/monto inyectados, organización, permiso, sesión ausente y claims
  del BFF incompatibles con el token Go. Storage denegado produce cero requests.
- Tras desactivar provider, no botones ni escrituras; GET conserva los registros.
  Nuevos pools/servidores Go/Next recuperan estados. Vista móvil inspeccionada y
  sin overflow en los cuatro proyectos; no certifica accesibilidad completa.
- PostgreSQL propio detenido/reiniciado: snapshot de pedidos/pagos/outbox idéntico
  SHA256 `c12639dcb9751b38d0d82570c0cf744c7a6f46889d0733ed35a65f3284bc6713`.
  Quedó detenido. No es crash recovery, PITR ni restore productivo.

| Receipt local | SHA256 |
|---|---|
| final-browser.log | 216c776f36c34522b73882014ffc27e600ef28d27c516a3506bd606a89e00dad |
| final-pg.log | 09eeaf8befae3c0b6ba0f213d4819aaf60092a8745c00d47720b1573cb1a087a |
| final-go-test.log | 725a9e45facd1e8e65d0cb31dd55e2d0a9c29ea2235c83430a1b54a61cf8f427 |
| rebuilt-web-tests.log | 3aec76ee95237467488b3d74e95d598d106d6eb04f65544cda915f4d1b32fae5 |
| legacy-bypass-red.log | 501761513d5b9a47909243990f2a52ffffc196d7e65c5fa9f0d5d09e70b66dac |

Reproducir con composición integral vigente, sus locks/builds, DB sintética
loopback elite_confirmation_*, TEST_DATABASE_URL y ELITE_WEB_ROOT absolutos.
ELITE_QUOTE_E2E=1, ELITE_ORDER_E2E=1 y ELITE_PAYMENT_E2E=1 activan el harness:

```text
go test ./internal/platform/httpapi -run '^TestQuoteAcceptanceBrowserPostgres$' -count=1 -timeout 22m -v
go test ./internal/platform/httpapi ./internal/platform/postgres ./cmd/electromobility-api -run 'TestPayment|TestCommerce|TestOperationRead' -count=1 -v
go test ./...
go vet ./...
go build ./...
pnpm test
pnpm build
```

FAIL444 selección corregida; FAIL445 espera de interacción WebKit corregida sin
retirar teclado/timeout/oráculo; FAIL446 reproducido202 y corregido409 en ruta
anterior. Regresiones y fuente canónica preservadas. FAIL423 anterior no se cierra
por coincidencia de un PASS; mantiene su expediente.

## Pendientes del mismo cierre

SDK/adapter real de solicitud/cobro, cuenta/sandbox, callback autenticado,
reconciliación y reglas de pago fallido/parcial/segundo intento siguen pendientes.
No arrancar consumidores externos de payment.requested por este PASS. La fila
de selección futura debe decidir proveedor y titular antes de credenciales.
Preparación inicial de entrega del pedido común y regla de liberación empresarial
siguen pendientes: los flows de devolución/reemplazo no los sustituyen.
Siguiente trabajo independiente: cerrar esa preparación sobre Fulfillment/Commerce
y contratos vigentes sin inferir crédito, financiación ni entrega autorizada.
Credenciales restantes y gates target permanecen en los mismos owners.

El harness monta módulos reales con issuer local; no prueba login con IdP real,
cloud, CDN/WAF, main desplegado, seguridad ofensiva, carga ni aceptación comercial.
T2802/T2804 y assurance integral no se promueven. No ZIP final ni porcentaje.
